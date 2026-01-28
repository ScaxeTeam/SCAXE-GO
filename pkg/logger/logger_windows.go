//go:build windows

package logger

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode = kernel32.NewProc("SetConsoleMode")
	getConsoleMode = kernel32.NewProc("GetConsoleMode")
)

const (
	enableVirtualTerminalProcessing = 0x0004
)

func enableWindowsVT() {

	stdout := os.Stdout.Fd()

	var mode uint32
	r, _, _ := getConsoleMode.Call(stdout, uintptr(unsafe.Pointer(&mode)))
	if r == 0 {

		colorEnabled = false
		return
	}

	mode |= enableVirtualTerminalProcessing
	r, _, _ = setConsoleMode.Call(stdout, uintptr(mode))
	if r == 0 {

		colorEnabled = false
	}
}
