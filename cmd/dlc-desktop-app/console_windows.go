package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleTitle = kernel32.NewProc("SetConsoleTitleW")
	procSetConsoleMode  = kernel32.NewProc("SetConsoleMode")
)

// setupConsole prepares the console window the app runs in when it is
// started by double-clicking. Failures are ignored: without a console (for
// example when output is redirected) there is nothing to set up.
func setupConsole(title string) {
	if p, err := syscall.UTF16PtrFromString(title); err == nil {
		procSetConsoleTitle.Call(uintptr(unsafe.Pointer(p)))
	}

	// Turn off QuickEdit mode. With it on, clicking inside the console starts
	// a text selection that pauses the program's output; the next log line
	// then blocks, and the app stops answering every browser until someone
	// presses a key in the console.
	const (
		enableQuickEditMode = 0x0040
		enableExtendedFlags = 0x0080
	)
	in, err := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		return
	}
	var mode uint32
	if syscall.GetConsoleMode(in, &mode) != nil {
		return
	}
	procSetConsoleMode.Call(uintptr(in), uintptr(mode&^enableQuickEditMode|enableExtendedFlags))
}
