package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

type startCommand struct{ cf commonFlags }

func (c *startCommand) Name() string                { return "start" }
func (c *startCommand) Description() string         { return "Starts a compute system." }
func (c *startCommand) ArgHelp() string             { return "" }
func (c *startCommand) SetupFlags(fs *flag.FlagSet) { setupCommonFlags(&c.cf, fs) }

func (c *startCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsStartComputeSystem(cs.handle, op, ""); err != nil {
		return err
	}
	if _, err := op.WaitResult(windows.INFINITE); err != nil {
		return err
	}
	return nil
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
	cf        commonFlags
	vmVersion *bool
}

func (c *propsCommand) Name() string        { return "props" }
func (c *propsCommand) Description() string { return "Lists compute system properties." }
func (c *propsCommand) ArgHelp() string     { return "" }
func (c *propsCommand) SetupFlags(fs *flag.FlagSet) {
	setupCommonFlags(&c.cf, fs)
	c.vmVersion = fs.Bool("vmversion", false, "Query for VmVersion property as well.")
}

func (c *propsCommand) Execute(state *state, fs *flag.FlagSet) error {
	_, cs, err := getCS(state, &c.cf)
	if err != nil {
		return err
	}
	pq := hcsschema.PropertyQuery{
		Queries: map[string]interface{}{
			"Basic": nil,
		},
	}
	if *c.vmVersion {
		pq.Queries["VmVersion"] = nil
	}
	j, err := json.Marshal(pq)
	if err != nil {
		return err
	}
	op := computecore.NewOperation(0)
	defer op.Close()
	if err := computecore.HcsGetComputeSystemProperties(cs.handle, op, string(j)); err != nil {
		return err
	}
	properties, err := op.WaitResult(windows.INFINITE)
	if err != nil {
		return err
	}
	var props struct {
		PropertyResponses struct {
			Basic struct {
				Response struct {
					State     string
					RuntimeId string
				}
			}
			VmVersion struct {
				Response struct {
					Major uint
					Minor uint
				}
			}
		}
	}
	if err := json.Unmarshal([]byte(properties), &props); err != nil {
		return err
	}
	fmt.Printf("State: %s\n", props.PropertyResponses.Basic.Response.State)
	fmt.Printf("RuntimeID: %s\n", props.PropertyResponses.Basic.Response.RuntimeId)
	fmt.Printf("VmVersion: %d.%d\n", props.PropertyResponses.VmVersion.Response.Major, props.PropertyResponses.VmVersion.Response.Minor)
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
	if err := computecore.HcsOpenComputeSystem(id, 0, &cs.handle); err != nil {
		return err
	}
	state.systems[id] = &cs
	return nil
}

type svcPropsCommand struct{}

func (c *svcPropsCommand) Name() string                { return "svcprops" }
func (c *svcPropsCommand) Description() string         { return "Lists HCS service properties." }
func (c *svcPropsCommand) ArgHelp() string             { return "" }
func (c *svcPropsCommand) SetupFlags(fs *flag.FlagSet) {}

func (c *svcPropsCommand) Execute(state *state, fs *flag.FlagSet) error {
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
	var properties *uint16
	if err := computecore.HcsGetServiceProperties(string(j), &properties); err != nil {
		return err
	}
	var props struct {
		PropertyResponses struct {
			Basic struct {
				Response struct {
					SupportedSchemaVersions []struct {
						Major uint
						Minor uint
					}
				}
			}
			ProcessorCapabilities struct {
				Response json.RawMessage
			}
		}
	}
	if err := json.Unmarshal([]byte(windows.UTF16PtrToString(properties)), &props); err != nil {
		return err
	}
	fmt.Printf("Supported schema versions: ")
	for i, sv := range props.PropertyResponses.Basic.Response.SupportedSchemaVersions {
		if i != 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%d.%d", sv.Major, sv.Minor)
	}
	fmt.Printf("\n")
	fmt.Printf("ProcessorCapabilities:\n")
	j, err = json.MarshalIndent(props.PropertyResponses.ProcessorCapabilities.Response, "\t", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("\t%s\n", string(j))
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
