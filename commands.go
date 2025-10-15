package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevpar/hcstool/internal/computecore"
	"github.com/kevpar/hcstool/internal/hcsschema"
	"github.com/kevpar/repl-go"
	"golang.org/x/sys/windows"
)

func allCommands() []repl.Command[*state] {
	return []repl.Command[*state]{
		&createCommand{},
		&startCommand{},
		&closeCommand{},
		&suspendCommand{},
		&resumeCommand{},
		&saveCommand{},
		&propsCommand{},
		&grantCommand{},
		&defaultCommand{},
		&listCommand{},
		&openCommand{},
		&svcPropsCommand{},
		&modifyCommand{},
		&lmSourceInitializeCommand{},
		&lmSourceStartCommand{},
		&lmTransferCommand{},
		&lmFinalizeCommand{},
	}
}

type state struct {
	def     string
	systems map[string]*cs
}

type cs struct {
	handle computecore.HCS_SYSTEM
}

func setupCommonFlags(cf *commonFlags, fs *flag.FlagSet) {
	cf.cs = fs.String("cs", "", "Specifies the compute system to operate on.")
}

type commonFlags struct {
	cs *string
}

func getCS(state *state, cf *commonFlags) (string, *cs, error) {
	key := state.def
	if cs := *cf.cs; cs != "" {
		key = cs
	}
	if key == "" {
		return "", nil, fmt.Errorf("must specify a default compute system or use -cs flag")
	}
	cs, ok := state.systems[key]
	if !ok {
		return "", nil, fmt.Errorf("compute system not opened: %s", key)
	}
	return key, cs, nil
}

type createCommand struct{ setDefault *bool }

func (c *createCommand) Name() string        { return "create" }
func (c *createCommand) Description() string { return "Creates a compute system." }
func (c *createCommand) ArgHelp() string     { return "ID PATH" }
func (c *createCommand) SetupFlags(fs *flag.FlagSet) {
	c.setDefault = fs.Bool("def", false, "Set the new compute system as the default.")
}

func (c *createCommand) Execute(state *state, fs *flag.FlagSet) error {
	id := fs.Arg(0)
	if _, ok := state.systems[id]; ok {
		return fmt.Errorf("compute system already open: %s", id)
	}
	var cs cs
	doc, err := os.ReadFile(fs.Arg(1))
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsCreateComputeSystem(fs.Arg(0), string(doc), op, nil, &cs.handle); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	state.systems[id] = &cs
	if *c.setDefault {
		state.def = id
	}
	return nil
}

type startCommand struct {
	cf        commonFlags
	migsocket *string
}

func (c *startCommand) Name() string        { return "start" }
func (c *startCommand) Description() string { return "Starts a compute system." }
func (c *startCommand) ArgHelp() string     { return "" }
func (c *startCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
	c.migsocket = fs.String("migsocket", "", "TCP address to dial for live migration connection.")
}

func (c *startCommand) Execute(state *state, fs *flag.FlagSet) error {
	var sock uintptr
	if *c.migsocket != "" {
		addr, err := netip.ParseAddrPort(*c.migsocket)
		if err != nil {
			return err
		}
		s, err := dial(addr)
		if err != nil {
			return err
		}
		sock = uintptr(s)
	}

	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	var optionsRaw []byte
	if sock != 0 {
		if err := computecore.HcsAddResourceToOperation(op, computecore.HcsResourceTypeSocket, "hcs:/VirtualMachine/LiveMigrationSocket", sock); err != nil {
			return err
		}
		options := hcsschema.StartOptions{
			DestinationMigrationOptions: &hcsschema.MigrationStartOptions{
				NetworkSettings: &hcsschema.MigrationNetworkSettings{
					SessionID: 1,
				},
			},
		}
		optionsRaw, err = json.Marshal(options)
		if err != nil {
			return err
		}
	}
	if err := computecore.HcsStartComputeSystem(cs.handle, op, string(optionsRaw)); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

func dial(addr netip.AddrPort) (_ windows.Handle, err error) {
	conn, err := windows.Socket(windows.AF_INET, windows.SOCK_STREAM, windows.IPPROTO_TCP)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			windows.Closesocket(conn)
		}
	}()
	fmt.Printf("connecting...\n")
	if err := windows.Connect(conn, &windows.SockaddrInet4{Port: int(addr.Port()), Addr: addr.Addr().As4()}); err != nil {
		return 0, err
	}
	fmt.Printf("connected\n")
	return conn, nil
}

type closeCommand struct{ cf commonFlags }

func (c *closeCommand) Name() string                { return "close" }
func (c *closeCommand) Description() string         { return "Closes a compute system." }
func (c *closeCommand) ArgHelp() string             { return "" }
func (c *closeCommand) SetupFlags(fs *flag.FlagSet) { setupCommonFlags(&c.cf, fs) }

func (c *closeCommand) Execute(state *state, fs *flag.FlagSet) error {
	id, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	computecore.HcsCloseComputeSystem(cs.handle)
	cs.handle = 0
	delete(state.systems, id)
	if state.def == id {
		state.def = ""
	}
	return nil
}

type suspendCommand struct{ cf commonFlags }

func (c *suspendCommand) Name() string                { return "pause" }
func (c *suspendCommand) Description() string         { return "Pauses a compute system." }
func (c *suspendCommand) ArgHelp() string             { return "" }
func (c *suspendCommand) SetupFlags(fs *flag.FlagSet) { setupCommonFlags(&c.cf, fs) }

func (c *suspendCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsPauseComputeSystem(cs.handle, op, ""); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

type resumeCommand struct{ cf commonFlags }

func (c *resumeCommand) Name() string                { return "resume" }
func (c *resumeCommand) Description() string         { return "Resumes a compute system." }
func (c *resumeCommand) ArgHelp() string             { return "" }
func (c *resumeCommand) SetupFlags(fs *flag.FlagSet) { setupCommonFlags(&c.cf, fs) }

func (c *resumeCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsResumeComputeSystem(cs.handle, op, ""); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

type saveCommand struct{ cf commonFlags }

func (c *saveCommand) Name() string                { return "save" }
func (c *saveCommand) Description() string         { return "Saves the compute system to disk." }
func (c *saveCommand) ArgHelp() string             { return "PATH" }
func (c *saveCommand) SetupFlags(fs *flag.FlagSet) { setupCommonFlags(&c.cf, fs) }

func (c *saveCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	path := fs.Arg(0)
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	so := hcsschema.SaveOptions{
		SaveType:          "ToFile",
		SaveStateFilePath: path,
	}
	j, err := json.Marshal(so)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsSaveComputeSystem(cs.handle, op, string(j)); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

type propsCommand struct {
	cf         commonFlags
	rawQuery   *string
	vmVersion  *bool
	compatInfo *bool
	procReqs   *bool
}

func (c *propsCommand) Name() string        { return "props" }
func (c *propsCommand) Description() string { return "Lists compute system properties." }
func (c *propsCommand) ArgHelp() string     { return "" }
func (c *propsCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
	c.rawQuery = fs.String("rawquery", "", "Exact query string to use.")
	c.vmVersion = fs.Bool("vmversion", false, "Query for VmVersion property as well.")
	c.compatInfo = fs.Bool("compatinfo", false, "Query for CompatibilityInfo property as well.")
	c.procReqs = fs.Bool("procreqs", false, "Query for VmProcessorRequirements as well.")
}

func (c *propsCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	var query string
	if c.rawQuery != nil {
		query = *c.rawQuery
	} else {
		pq := hcsschema.PropertyQuery{
			Queries: map[string]interface{}{
				"Basic": nil,
			},
		}
		if *c.vmVersion {
			pq.Queries["VmVersion"] = nil
		}
		if *c.compatInfo {
			pq.Queries["CompatibilityInfo"] = nil
		}
		if *c.procReqs {
			pq.Queries["VmProcessorRequirements"] = nil
		}
		j, err := json.Marshal(pq)
		if err != nil {
			return err
		}
		query = string(j)
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsGetComputeSystemProperties(cs.handle, op, query); err != nil {
		return err
	}
	properties, err := op.WaitResult(windows.INFINITE)
	if err != nil {
		return err
	}
	var results any
	if err := json.Unmarshal([]byte(properties), &results); err != nil {
		return err
	}
	resultsStr, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(resultsStr))

	return nil
}

type grantCommand struct{}

func (c *grantCommand) Name() string { return "grant" }
func (c *grantCommand) Description() string {
	return "Grants a given compute system ID access to a file."
}
func (c *grantCommand) ArgHelp() string             { return "ID PATH" }
func (c *grantCommand) SetupFlags(fs *flag.FlagSet) {}

func (c *grantCommand) Execute(state *state, fs *flag.FlagSet) error {
	path := fs.Arg(1)
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := computecore.HcsGrantVmAccess(fs.Arg(0), path); err != nil {
		return err
	}
	return nil
}

type defaultCommand struct{ unset *bool }

func (c *defaultCommand) Name() string { return "default" }
func (c *defaultCommand) Description() string {
	return "Sets the default compute system to operate on."
}
func (c *defaultCommand) ArgHelp() string { return "ID" }
func (c *defaultCommand) SetupFlags(fs *flag.FlagSet) {
	c.unset = fs.Bool("unset", false, "Set the default to nothing.")
}

func (c *defaultCommand) Execute(state *state, fs *flag.FlagSet) error {
	if *c.unset {
		state.def = ""
		return nil
	}
	id := fs.Arg(0)
	if _, ok := state.systems[id]; !ok {
		return fmt.Errorf("compute system not found: %s", id)
	}
	state.def = id
	return nil
}

type listCommand struct{ all *bool }

func (c *listCommand) Name() string        { return "list" }
func (c *listCommand) Description() string { return "Lists compute systems." }
func (c *listCommand) ArgHelp() string     { return "" }
func (c *listCommand) SetupFlags(fs *flag.FlagSet) {
	c.all = fs.Bool("all", false, "Show all systems instead of only those you have open.")
}

func (c *listCommand) Execute(state *state, fs *flag.FlagSet) error {
	if *c.all {
		op := computecore.NewOperation(0)
		defer op.Close()
		if err := computecore.HcsEnumerateComputeSystems("", op); err != nil {
			return err
		}
		systemsRaw, err := op.WaitResult(windows.INFINITE)
		if err != nil {
			return err
		}
		type systemData struct {
			ID         string `json:"Id"`
			Name       string
			SystemType string
			Owner      string
			State      string
		}
		var systems []systemData
		if err := json.Unmarshal([]byte(systemsRaw), &systems); err != nil {
			return err
		}
		if err := printTable(
			[]colInfo{{"ID", "%s"}, {"NAME", "%s"}, {"TYPE", "%s"}, {"OWNER", "%s"}, {"STATE", "%s"}},
			systems,
			func(rd systemData) []any { return []any{rd.ID, rd.Name, rd.SystemType, rd.Owner, rd.State} },
		); err != nil {
			return err
		}
	} else {
		for id := range state.systems {
			fmt.Printf("%s\n", id)
		}
	}
	return nil
}

type openCommand struct{}

func (c *openCommand) Name() string                { return "open" }
func (c *openCommand) Description() string         { return "Opens a compute systems." }
func (c *openCommand) ArgHelp() string             { return "ID" }
func (c *openCommand) SetupFlags(fs *flag.FlagSet) {}

func (c *openCommand) Execute(state *state, fs *flag.FlagSet) error {
	id := fs.Arg(0)
	if _, ok := state.systems[id]; ok {
		return fmt.Errorf("compute system already open: %s", id)
	}
	var cs cs
	if err := computecore.HcsOpenComputeSystem(id, windows.GENERIC_ALL, &cs.handle); err != nil {
		return err
	}
	state.systems[id] = &cs
	return nil
}

type svcPropsCommand struct {
	rawQuery *string
}

func (c *svcPropsCommand) Name() string        { return "svcprops" }
func (c *svcPropsCommand) Description() string { return "Lists HCS service properties." }
func (c *svcPropsCommand) ArgHelp() string     { return "" }
func (c *svcPropsCommand) SetupFlags(fs *flag.FlagSet) {
	c.rawQuery = fs.String("rawquery", "", "Exact query string to use.")
}

func (c *svcPropsCommand) Execute(state *state, fs *flag.FlagSet) error {
	var query string
	if c.rawQuery != nil {
		query = *c.rawQuery
	} else {
		pq := struct {
			PropertyQueries map[string]any
		}{
			PropertyQueries: map[string]any{
				"Basic":                 nil,
				"ProcessorCapabilities": nil,
			},
		}
		j, err := json.Marshal(pq)
		if err != nil {
			return err
		}
		query = string(j)
	}
	var properties *uint16
	if err := computecore.HcsGetServiceProperties(query, &properties); err != nil {
		return err
	}
	j, err := json.MarshalIndent(json.RawMessage(windows.UTF16PtrToString(properties)), "\t", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("Properties:\n\t%s\n", string(j))
	return nil
}

type modifyCommand struct{ cf commonFlags }

func (c *modifyCommand) Name() string        { return "modify" }
func (c *modifyCommand) Description() string { return "Modifies a compute system." }
func (c *modifyCommand) ArgHelp() string     { return "add|remove|update PATH SETTINGS" }
func (c *modifyCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
}

func (c *modifyCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	typ, ok := map[string]string{
		"add":    "Add",
		"remove": "Remove",
		"update": "Update",
	}[strings.ToLower(fs.Arg(0))]
	if !ok {
		return fmt.Errorf("unrecognized operation: %s", fs.Arg(0))
	}
	req := hcsschema.ModifySettingRequest{
		RequestType:  typ,
		ResourcePath: fs.Arg(1),
		Settings:     json.RawMessage(fs.Arg(2)),
	}
	j, err := json.Marshal(req)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsModifyComputeSystem(cs.handle, op, string(j), 0); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

type lmSourceInitializeCommand struct{ cf commonFlags }

func (c *lmSourceInitializeCommand) Name() string { return "lmsrcinit" }
func (c *lmSourceInitializeCommand) Description() string {
	return "Initializes the source for live migration."
}
func (c *lmSourceInitializeCommand) ArgHelp() string { return "" }
func (c *lmSourceInitializeCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
}

func (c *lmSourceInitializeCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	options := hcsschema.MigrationInitializeOptions{}
	optionsRaw, err := json.Marshal(options)
	if err != nil {
		return err
	}
	if err := computecore.HcsInitializeLiveMigrationOnSource(cs.handle, op, string(optionsRaw)); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

type lmSourceStartCommand struct{ cf commonFlags }

func (c *lmSourceStartCommand) Name() string { return "lmsrcstart" }
func (c *lmSourceStartCommand) Description() string {
	return "Starts the source for live migration."
}
func (c *lmSourceStartCommand) ArgHelp() string { return "SOCKET" }
func (c *lmSourceStartCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
}

func (c *lmSourceStartCommand) Execute(state *state, fs *flag.FlagSet) error {
	addr, err := netip.ParseAddrPort(fs.Arg(0))
	if err != nil {
		return err
	}
	sock, err := listen(addr)
	if err != nil {
		return err
	}

	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsAddResourceToOperation(op, computecore.HcsResourceTypeSocket, "hcs:/VirtualMachine/LiveMigrationSocket", uintptr(sock)); err != nil {
		return err
	}
	options := hcsschema.MigrationStartOptions{
		NetworkSettings: &hcsschema.MigrationNetworkSettings{
			SessionID: 1,
		},
	}
	optionsRaw, err := json.Marshal(options)
	if err != nil {
		return err
	}
	if err := computecore.HcsStartLiveMigrationOnSource(cs.handle, op, string(optionsRaw)); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

func listen(addr netip.AddrPort) (_ windows.Handle, err error) {
	l, err := windows.Socket(windows.AF_INET, windows.SOCK_STREAM, windows.IPPROTO_TCP)
	if err != nil {
		return 0, err
	}
	defer windows.Closesocket(l)
	if err := windows.Bind(l, &windows.SockaddrInet4{Port: int(addr.Port()), Addr: addr.Addr().As4()}); err != nil {
		return 0, err
	}
	conn, err := windows.Socket(windows.AF_INET, windows.SOCK_STREAM, windows.IPPROTO_TCP)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			windows.Closesocket(conn)
		}
	}()
	if err := windows.Listen(l, 1); err != nil {
		return 0, err
	}
	var buf [64]byte
	var recvd uint32
	event, err := windows.CreateEvent(nil, 1, 0, nil)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(event)
	overlapped := windows.Overlapped{HEvent: event}
	if err := windows.AcceptEx(l, conn, &buf[0], 0, 32, 32, &recvd, &overlapped); err != nil && err != windows.ERROR_IO_PENDING {
		return 0, err
	}
	fmt.Printf("connecting...\n")
	if _, err := windows.WaitForSingleObject(event, windows.INFINITE); err != nil {
		return 0, err
	}
	fmt.Printf("connected\n")
	return conn, nil
}

type lmTransferCommand struct{ cf commonFlags }

func (c *lmTransferCommand) Name() string { return "lmtransfer" }
func (c *lmTransferCommand) Description() string {
	return "Initiates live migration transfer."
}
func (c *lmTransferCommand) ArgHelp() string { return "" }
func (c *lmTransferCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
}

func (c *lmTransferCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	options := hcsschema.MigrationTransferOptions{}
	optionsRaw, err := json.Marshal(options)
	if err != nil {
		return err
	}
	if err := computecore.HcsStartLiveMigrationTransfer(cs.handle, op, string(optionsRaw)); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}

type lmFinalizeCommand struct{ cf commonFlags }

func (c *lmFinalizeCommand) Name() string { return "lmfinalize" }
func (c *lmFinalizeCommand) Description() string {
	return "Finalizes the live migration."
}
func (c *lmFinalizeCommand) ArgHelp() string { return "" }
func (c *lmFinalizeCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
}

func (c *lmFinalizeCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	options := hcsschema.MigrationFinalizedOptions{}
	optionsRaw, err := json.Marshal(options)
	if err != nil {
		return err
	}
	if err := computecore.HcsFinalizeLiveMigration(cs.handle, op, string(optionsRaw)); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
}
