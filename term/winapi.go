// +build windows

package term

import (
	"syscall"
)

var (
	kernel32DLL        = syscall.NewLazyDLL("kernel32.dll")
	getStdHandleProc   = kernel32DLL.NewProc("GetStdHandle")
	setConsoleModeProc = kernel32DLL.NewProc("SetConsoleMode")
)

const STD_OUTPUT_HANDLE = uintptr(1) + ^uintptr(11)
const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004

func GetStdHandle(handle syscall.Handle) (syscall.Handle, error) {
	r1, _, err := getStdHandleProc.Call(uintptr(handle), 0)
	if r1 != 0 && syscall.Handle(r1) != syscall.InvalidHandle {
		return syscall.Handle(r1), nil
	} else {
		return syscall.InvalidHandle, winError(err)
	}
}

func GetConsoleMode(handle syscall.Handle) (uint32, error) {
	var mode uint32 = 0
	err := syscall.GetConsoleMode(handle, &mode)
	return mode, err
}

func SetConsoleMode(handle syscall.Handle, mode uint32) error {
	r1, _, err := setConsoleModeProc.Call(uintptr(handle), uintptr(mode), 0)
	if r1 != 0 {
		return nil
	} else {
		return winError(err)
	}
}

func winError(err error) error {
	if err != nil {
		return syscall.GetLastError()
	} else {
		return syscall.EINVAL
	}
}
