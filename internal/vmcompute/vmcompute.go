//go:build windows

package vmcompute

import (
	"syscall"
)

//go:generate go run github.com/Microsoft/go-winio/tools/mkwinsyscall -output zsyscall_windows.go vmcompute.go

//sys HcsEnumerateComputeSystems(query string, computeSystems **uint16, result **uint16) (hr error) = vmcompute.HcsEnumerateComputeSystems?
//sys HcsCreateComputeSystem(id string, configuration string, identity syscall.Handle, computeSystem *HcsSystem, result **uint16) (hr error) = vmcompute.HcsCreateComputeSystem?
//sys HcsOpenComputeSystem(id string, computeSystem *HcsSystem, result **uint16) (hr error) = vmcompute.HcsOpenComputeSystem?
//sys HcsCloseComputeSystem(computeSystem HcsSystem) (hr error) = vmcompute.HcsCloseComputeSystem?
//sys HcsStartComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) = vmcompute.HcsStartComputeSystem?
//sys HcsShutdownComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) = vmcompute.HcsShutdownComputeSystem?
//sys HcsTerminateComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) = vmcompute.HcsTerminateComputeSystem?
//sys HcsPauseComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) = vmcompute.HcsPauseComputeSystem?
//sys HcsResumeComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) = vmcompute.HcsResumeComputeSystem?
//sys HcsGetComputeSystemProperties(computeSystem HcsSystem, propertyQuery string, properties **uint16, result **uint16) (hr error) = vmcompute.HcsGetComputeSystemProperties?
//sys HcsModifyComputeSystem(computeSystem HcsSystem, configuration string, result **uint16) (hr error) = vmcompute.HcsModifyComputeSystem?
//sys HcsModifyServiceSettings(settings string, result **uint16) (hr error) = vmcompute.HcsModifyServiceSettings?
//sys HcsRegisterComputeSystemCallback(computeSystem HcsSystem, callback uintptr, context uintptr, callbackHandle *HcsCallback) (hr error) = vmcompute.HcsRegisterComputeSystemCallback?
//sys HcsUnregisterComputeSystemCallback(callbackHandle HcsCallback) (hr error) = vmcompute.HcsUnregisterComputeSystemCallback?
//sys HcsSaveComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) = vmcompute.HcsSaveComputeSystem?

//sys HcsCreateProcess(computeSystem HcsSystem, processParameters string, processInformation *HcsProcessInformation, process *HcsProcess, result **uint16) (hr error) = vmcompute.HcsCreateProcess?
//sys HcsOpenProcess(computeSystem HcsSystem, pid uint32, process *HcsProcess, result **uint16) (hr error) = vmcompute.HcsOpenProcess?
//sys HcsCloseProcess(process HcsProcess) (hr error) = vmcompute.HcsCloseProcess?
//sys HcsTerminateProcess(process HcsProcess, result **uint16) (hr error) = vmcompute.HcsTerminateProcess?
//sys HcsSignalProcess(process HcsProcess, options string, result **uint16) (hr error) = vmcompute.HcsSignalProcess?
//sys HcsGetProcessInfo(process HcsProcess, processInformation *HcsProcessInformation, result **uint16) (hr error) = vmcompute.HcsGetProcessInfo?
//sys HcsGetProcessProperties(process HcsProcess, processProperties **uint16, result **uint16) (hr error) = vmcompute.HcsGetProcessProperties?
//sys HcsModifyProcess(process HcsProcess, settings string, result **uint16) (hr error) = vmcompute.HcsModifyProcess?
//sys HcsGetServiceProperties(propertyQuery string, properties **uint16, result **uint16) (hr error) = vmcompute.HcsGetServiceProperties?
//sys HcsRegisterProcessCallback(process HcsProcess, callback uintptr, context uintptr, callbackHandle *HcsCallback) (hr error) = vmcompute.HcsRegisterProcessCallback?
//sys HcsUnregisterProcessCallback(callbackHandle HcsCallback) (hr error) = vmcompute.HcsUnregisterProcessCallback?

//sys GrantVmAccess(vmid string, filepath string) (hr error) = vmcompute.GrantVmAccess?

// errVmcomputeOperationPending is an error encountered when the operation is being completed asynchronously
const ErrVmcomputeOperationPending = syscall.Errno(0xC0370103)

// HcsSystem is the handle associated with a created compute system.
type HcsSystem syscall.Handle

// HcsProcess is the handle associated with a created process in a compute
// system.
type HcsProcess syscall.Handle

// HcsCallback is the handle associated with the function to call when events
// occur.
type HcsCallback syscall.Handle

// HcsProcessInformation is the structure used when creating or getting process
// info.
type HcsProcessInformation struct {
	// ProcessId is the pid of the created process.
	ProcessId uint32
	_         uint32 // reserved padding
	// StdInput is the handle associated with the stdin of the process.
	StdInput syscall.Handle
	// StdOutput is the handle associated with the stdout of the process.
	StdOutput syscall.Handle
	// StdError is the handle associated with the stderr of the process.
	StdError syscall.Handle
}
