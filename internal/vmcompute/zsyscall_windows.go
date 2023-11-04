//go:build windows

// Code generated by 'go generate' using "github.com/Microsoft/go-winio/tools/mkwinsyscall"; DO NOT EDIT.

package vmcompute

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modvmcompute = windows.NewLazySystemDLL("vmcompute.dll")

	procGrantVmAccess                      = modvmcompute.NewProc("GrantVmAccess")
	procHcsCloseComputeSystem              = modvmcompute.NewProc("HcsCloseComputeSystem")
	procHcsCloseProcess                    = modvmcompute.NewProc("HcsCloseProcess")
	procHcsCreateComputeSystem             = modvmcompute.NewProc("HcsCreateComputeSystem")
	procHcsCreateProcess                   = modvmcompute.NewProc("HcsCreateProcess")
	procHcsEnumerateComputeSystems         = modvmcompute.NewProc("HcsEnumerateComputeSystems")
	procHcsGetComputeSystemProperties      = modvmcompute.NewProc("HcsGetComputeSystemProperties")
	procHcsGetProcessInfo                  = modvmcompute.NewProc("HcsGetProcessInfo")
	procHcsGetProcessProperties            = modvmcompute.NewProc("HcsGetProcessProperties")
	procHcsGetServiceProperties            = modvmcompute.NewProc("HcsGetServiceProperties")
	procHcsModifyComputeSystem             = modvmcompute.NewProc("HcsModifyComputeSystem")
	procHcsModifyProcess                   = modvmcompute.NewProc("HcsModifyProcess")
	procHcsModifyServiceSettings           = modvmcompute.NewProc("HcsModifyServiceSettings")
	procHcsOpenComputeSystem               = modvmcompute.NewProc("HcsOpenComputeSystem")
	procHcsOpenProcess                     = modvmcompute.NewProc("HcsOpenProcess")
	procHcsPauseComputeSystem              = modvmcompute.NewProc("HcsPauseComputeSystem")
	procHcsRegisterComputeSystemCallback   = modvmcompute.NewProc("HcsRegisterComputeSystemCallback")
	procHcsRegisterProcessCallback         = modvmcompute.NewProc("HcsRegisterProcessCallback")
	procHcsResumeComputeSystem             = modvmcompute.NewProc("HcsResumeComputeSystem")
	procHcsSaveComputeSystem               = modvmcompute.NewProc("HcsSaveComputeSystem")
	procHcsShutdownComputeSystem           = modvmcompute.NewProc("HcsShutdownComputeSystem")
	procHcsSignalProcess                   = modvmcompute.NewProc("HcsSignalProcess")
	procHcsStartComputeSystem              = modvmcompute.NewProc("HcsStartComputeSystem")
	procHcsTerminateComputeSystem          = modvmcompute.NewProc("HcsTerminateComputeSystem")
	procHcsTerminateProcess                = modvmcompute.NewProc("HcsTerminateProcess")
	procHcsUnregisterComputeSystemCallback = modvmcompute.NewProc("HcsUnregisterComputeSystemCallback")
	procHcsUnregisterProcessCallback       = modvmcompute.NewProc("HcsUnregisterProcessCallback")
)

func GrantVmAccess(vmid string, filepath string) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(vmid)
	if hr != nil {
		return
	}
	var _p1 *uint16
	_p1, hr = syscall.UTF16PtrFromString(filepath)
	if hr != nil {
		return
	}
	return _GrantVmAccess(_p0, _p1)
}

func _GrantVmAccess(vmid *uint16, filepath *uint16) (hr error) {
	hr = procGrantVmAccess.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procGrantVmAccess.Addr(), 2, uintptr(unsafe.Pointer(vmid)), uintptr(unsafe.Pointer(filepath)), 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsCloseComputeSystem(computeSystem HcsSystem) (hr error) {
	hr = procHcsCloseComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsCloseComputeSystem.Addr(), 1, uintptr(computeSystem), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsCloseProcess(process HcsProcess) (hr error) {
	hr = procHcsCloseProcess.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsCloseProcess.Addr(), 1, uintptr(process), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsCreateComputeSystem(id string, configuration string, identity syscall.Handle, computeSystem *HcsSystem, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(id)
	if hr != nil {
		return
	}
	var _p1 *uint16
	_p1, hr = syscall.UTF16PtrFromString(configuration)
	if hr != nil {
		return
	}
	return _HcsCreateComputeSystem(_p0, _p1, identity, computeSystem, result)
}

func _HcsCreateComputeSystem(id *uint16, configuration *uint16, identity syscall.Handle, computeSystem *HcsSystem, result **uint16) (hr error) {
	hr = procHcsCreateComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall6(procHcsCreateComputeSystem.Addr(), 5, uintptr(unsafe.Pointer(id)), uintptr(unsafe.Pointer(configuration)), uintptr(identity), uintptr(unsafe.Pointer(computeSystem)), uintptr(unsafe.Pointer(result)), 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsCreateProcess(computeSystem HcsSystem, processParameters string, processInformation *HcsProcessInformation, process *HcsProcess, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(processParameters)
	if hr != nil {
		return
	}
	return _HcsCreateProcess(computeSystem, _p0, processInformation, process, result)
}

func _HcsCreateProcess(computeSystem HcsSystem, processParameters *uint16, processInformation *HcsProcessInformation, process *HcsProcess, result **uint16) (hr error) {
	hr = procHcsCreateProcess.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall6(procHcsCreateProcess.Addr(), 5, uintptr(computeSystem), uintptr(unsafe.Pointer(processParameters)), uintptr(unsafe.Pointer(processInformation)), uintptr(unsafe.Pointer(process)), uintptr(unsafe.Pointer(result)), 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsEnumerateComputeSystems(query string, computeSystems **uint16, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(query)
	if hr != nil {
		return
	}
	return _HcsEnumerateComputeSystems(_p0, computeSystems, result)
}

func _HcsEnumerateComputeSystems(query *uint16, computeSystems **uint16, result **uint16) (hr error) {
	hr = procHcsEnumerateComputeSystems.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsEnumerateComputeSystems.Addr(), 3, uintptr(unsafe.Pointer(query)), uintptr(unsafe.Pointer(computeSystems)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsGetComputeSystemProperties(computeSystem HcsSystem, propertyQuery string, properties **uint16, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(propertyQuery)
	if hr != nil {
		return
	}
	return _HcsGetComputeSystemProperties(computeSystem, _p0, properties, result)
}

func _HcsGetComputeSystemProperties(computeSystem HcsSystem, propertyQuery *uint16, properties **uint16, result **uint16) (hr error) {
	hr = procHcsGetComputeSystemProperties.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall6(procHcsGetComputeSystemProperties.Addr(), 4, uintptr(computeSystem), uintptr(unsafe.Pointer(propertyQuery)), uintptr(unsafe.Pointer(properties)), uintptr(unsafe.Pointer(result)), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsGetProcessInfo(process HcsProcess, processInformation *HcsProcessInformation, result **uint16) (hr error) {
	hr = procHcsGetProcessInfo.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsGetProcessInfo.Addr(), 3, uintptr(process), uintptr(unsafe.Pointer(processInformation)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsGetProcessProperties(process HcsProcess, processProperties **uint16, result **uint16) (hr error) {
	hr = procHcsGetProcessProperties.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsGetProcessProperties.Addr(), 3, uintptr(process), uintptr(unsafe.Pointer(processProperties)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsGetServiceProperties(propertyQuery string, properties **uint16, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(propertyQuery)
	if hr != nil {
		return
	}
	return _HcsGetServiceProperties(_p0, properties, result)
}

func _HcsGetServiceProperties(propertyQuery *uint16, properties **uint16, result **uint16) (hr error) {
	hr = procHcsGetServiceProperties.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsGetServiceProperties.Addr(), 3, uintptr(unsafe.Pointer(propertyQuery)), uintptr(unsafe.Pointer(properties)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsModifyComputeSystem(computeSystem HcsSystem, configuration string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(configuration)
	if hr != nil {
		return
	}
	return _HcsModifyComputeSystem(computeSystem, _p0, result)
}

func _HcsModifyComputeSystem(computeSystem HcsSystem, configuration *uint16, result **uint16) (hr error) {
	hr = procHcsModifyComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsModifyComputeSystem.Addr(), 3, uintptr(computeSystem), uintptr(unsafe.Pointer(configuration)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsModifyProcess(process HcsProcess, settings string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(settings)
	if hr != nil {
		return
	}
	return _HcsModifyProcess(process, _p0, result)
}

func _HcsModifyProcess(process HcsProcess, settings *uint16, result **uint16) (hr error) {
	hr = procHcsModifyProcess.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsModifyProcess.Addr(), 3, uintptr(process), uintptr(unsafe.Pointer(settings)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsModifyServiceSettings(settings string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(settings)
	if hr != nil {
		return
	}
	return _HcsModifyServiceSettings(_p0, result)
}

func _HcsModifyServiceSettings(settings *uint16, result **uint16) (hr error) {
	hr = procHcsModifyServiceSettings.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsModifyServiceSettings.Addr(), 2, uintptr(unsafe.Pointer(settings)), uintptr(unsafe.Pointer(result)), 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsOpenComputeSystem(id string, computeSystem *HcsSystem, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(id)
	if hr != nil {
		return
	}
	return _HcsOpenComputeSystem(_p0, computeSystem, result)
}

func _HcsOpenComputeSystem(id *uint16, computeSystem *HcsSystem, result **uint16) (hr error) {
	hr = procHcsOpenComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsOpenComputeSystem.Addr(), 3, uintptr(unsafe.Pointer(id)), uintptr(unsafe.Pointer(computeSystem)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsOpenProcess(computeSystem HcsSystem, pid uint32, process *HcsProcess, result **uint16) (hr error) {
	hr = procHcsOpenProcess.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall6(procHcsOpenProcess.Addr(), 4, uintptr(computeSystem), uintptr(pid), uintptr(unsafe.Pointer(process)), uintptr(unsafe.Pointer(result)), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsPauseComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(options)
	if hr != nil {
		return
	}
	return _HcsPauseComputeSystem(computeSystem, _p0, result)
}

func _HcsPauseComputeSystem(computeSystem HcsSystem, options *uint16, result **uint16) (hr error) {
	hr = procHcsPauseComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsPauseComputeSystem.Addr(), 3, uintptr(computeSystem), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsRegisterComputeSystemCallback(computeSystem HcsSystem, callback uintptr, context uintptr, callbackHandle *HcsCallback) (hr error) {
	hr = procHcsRegisterComputeSystemCallback.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall6(procHcsRegisterComputeSystemCallback.Addr(), 4, uintptr(computeSystem), uintptr(callback), uintptr(context), uintptr(unsafe.Pointer(callbackHandle)), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsRegisterProcessCallback(process HcsProcess, callback uintptr, context uintptr, callbackHandle *HcsCallback) (hr error) {
	hr = procHcsRegisterProcessCallback.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall6(procHcsRegisterProcessCallback.Addr(), 4, uintptr(process), uintptr(callback), uintptr(context), uintptr(unsafe.Pointer(callbackHandle)), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsResumeComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(options)
	if hr != nil {
		return
	}
	return _HcsResumeComputeSystem(computeSystem, _p0, result)
}

func _HcsResumeComputeSystem(computeSystem HcsSystem, options *uint16, result **uint16) (hr error) {
	hr = procHcsResumeComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsResumeComputeSystem.Addr(), 3, uintptr(computeSystem), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsSaveComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(options)
	if hr != nil {
		return
	}
	return _HcsSaveComputeSystem(computeSystem, _p0, result)
}

func _HcsSaveComputeSystem(computeSystem HcsSystem, options *uint16, result **uint16) (hr error) {
	hr = procHcsSaveComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsSaveComputeSystem.Addr(), 3, uintptr(computeSystem), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsShutdownComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(options)
	if hr != nil {
		return
	}
	return _HcsShutdownComputeSystem(computeSystem, _p0, result)
}

func _HcsShutdownComputeSystem(computeSystem HcsSystem, options *uint16, result **uint16) (hr error) {
	hr = procHcsShutdownComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsShutdownComputeSystem.Addr(), 3, uintptr(computeSystem), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsSignalProcess(process HcsProcess, options string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(options)
	if hr != nil {
		return
	}
	return _HcsSignalProcess(process, _p0, result)
}

func _HcsSignalProcess(process HcsProcess, options *uint16, result **uint16) (hr error) {
	hr = procHcsSignalProcess.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsSignalProcess.Addr(), 3, uintptr(process), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsStartComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(options)
	if hr != nil {
		return
	}
	return _HcsStartComputeSystem(computeSystem, _p0, result)
}

func _HcsStartComputeSystem(computeSystem HcsSystem, options *uint16, result **uint16) (hr error) {
	hr = procHcsStartComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsStartComputeSystem.Addr(), 3, uintptr(computeSystem), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsTerminateComputeSystem(computeSystem HcsSystem, options string, result **uint16) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(options)
	if hr != nil {
		return
	}
	return _HcsTerminateComputeSystem(computeSystem, _p0, result)
}

func _HcsTerminateComputeSystem(computeSystem HcsSystem, options *uint16, result **uint16) (hr error) {
	hr = procHcsTerminateComputeSystem.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsTerminateComputeSystem.Addr(), 3, uintptr(computeSystem), uintptr(unsafe.Pointer(options)), uintptr(unsafe.Pointer(result)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsTerminateProcess(process HcsProcess, result **uint16) (hr error) {
	hr = procHcsTerminateProcess.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsTerminateProcess.Addr(), 2, uintptr(process), uintptr(unsafe.Pointer(result)), 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsUnregisterComputeSystemCallback(callbackHandle HcsCallback) (hr error) {
	hr = procHcsUnregisterComputeSystemCallback.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsUnregisterComputeSystemCallback.Addr(), 1, uintptr(callbackHandle), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func HcsUnregisterProcessCallback(callbackHandle HcsCallback) (hr error) {
	hr = procHcsUnregisterProcessCallback.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.Syscall(procHcsUnregisterProcessCallback.Addr(), 1, uintptr(callbackHandle), 0, 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}