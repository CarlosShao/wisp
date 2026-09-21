//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

// AttachConsole is not exported by golang.org/x/sys/windows; declare it.
var (
	modKernel32       = windows.NewLazySystemDLL("kernel32.dll")
	procAttachConsole = modKernel32.NewProc("AttachConsole")
)

// attachParentProcess is the Win32 ATTACH_PARENT_PROCESS pseudo handle: (DWORD)-1.
const attachParentProcess = 0xFFFFFFFF

const (
	consoleAccessMask = windows.GENERIC_READ | windows.GENERIC_WRITE
	consoleShareMode  = windows.FILE_SHARE_READ | windows.FILE_SHARE_WRITE
)

// attachParentConsole makes console output work when this binary is built with
// the windowsgui subsystem (the final GUI build, ticket 07): such processes
// start detached from any console, so `wisp run` / `wisp doctor` must attach to
// the parent console and rebind the standard handles before printing
// (SPEC-11 §2.2, SPEC-01 §3).
//
// No-op when a usable stdout already exists: console-subsystem builds and
// redirected output (`wisp doctor > file`) must stay untouched.
func attachParentConsole() {
	out, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err == nil && out != 0 && out != windows.InvalidHandle {
		return
	}
	// Fails harmlessly when already attached or when no parent console exists.
	_, _, _ = procAttachConsole.Call(uintptr(attachParentProcess))
	rebindStdHandle(windows.STD_OUTPUT_HANDLE, "CONOUT$")
	rebindStdHandle(windows.STD_ERROR_HANDLE, "CONOUT$")
	rebindStdHandle(windows.STD_INPUT_HANDLE, "CONIN$")
}

// rebindStdHandle opens a console device and installs it as one of the three
// standard handles, including the os.File wrappers Go's fmt package writes to.
func rebindStdHandle(stdHandleID uint32, device string) {
	h, err := windows.CreateFile(windows.StringToUTF16Ptr(device), consoleAccessMask, consoleShareMode, nil, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		return
	}
	_ = windows.SetStdHandle(stdHandleID, h)
	switch stdHandleID {
	case windows.STD_OUTPUT_HANDLE:
		os.Stdout = os.NewFile(uintptr(h), device)
	case windows.STD_ERROR_HANDLE:
		os.Stderr = os.NewFile(uintptr(h), device)
	case windows.STD_INPUT_HANDLE:
		os.Stdin = os.NewFile(uintptr(h), device)
	}
}
