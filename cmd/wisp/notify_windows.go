//go:build windows

package main

// The CLI-side notification surface (D34 notify / D10 result presentation).
//
// It is a tray balloon tip, not a toast: a balloon cannot steal focus, cannot
// be left open by a click, and disappears on its own, which is what a
// non-GUI console run owes a user who is working in another window. The
// message-only window is created, used and destroyed inside one call on one
// OS thread, so the run leaves no window, no tray icon and no goroutine
// behind it (the resident process's own tray icon is ticket 07/62's, in
// internal/ball, and nothing here touches it).

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	nimAdd    = 0x00000000
	nimModify = 0x00000001
	nimDelete = 0x00000002

	nifMessage = 0x00000001
	nifIcon    = 0x00000002
	nifTip     = 0x00000004
	nifInfo    = 0x00000010

	niifInfo        = 0x00000001
	hwndMessage     = ^windows.HWND(2) // HWND_MESSAGE == -3
	wmDestroy       = 0x0002
	wmClose         = 0x0010
	notifyBalloonID = 0x5712

	// balloonTipClass is deliberately NOT the ball's window class: an
	// enumeration-based test that looks for a Wisp ball must never match a
	// console run's message-only window.
	balloonTipClass = "WispCLINotifyWindow"
)

// notifyBalloonTimeout is how long the shell shows the balloon, and therefore
// how long this call holds the icon before removing it.
var notifyBalloonTimeout = 1200 * time.Millisecond

type notifyIconDataW struct {
	cbSize           uint32
	hWnd             windows.Handle
	uID              uint32
	uFlags           uint32
	uCallbackMessage uint16
	hIcon            windows.Handle
	szTip            [64]uint16
	dwInfoFlags      uint32
	szInfo           [256]uint16
	szInfoTitle      [64]uint16
}

var (
	modShell32        = windows.NewLazySystemDLL("shell32.dll")
	pShellNotifyIconW = modShell32.NewProc("Shell_NotifyIconW")
	pRegisterClassExW = modUser32ForNotify.NewProc("RegisterClassExW")
	pCreateWindowExW  = modUser32ForNotify.NewProc("CreateWindowExW")
	pDestroyWindow    = modUser32ForNotify.NewProc("DestroyWindow")
	pDefWindowProcW   = modUser32ForNotify.NewProc("DefWindowProcW")
	pLoadIconW        = modUser32ForNotify.NewProc("LoadIconW")
	pGetModuleHandleW = modKernel32ForNotify.NewProc("GetModuleHandleW")

	balloonClassOnce sync.Once
	balloonClassErr  error
	balloonWndProc   = windows.NewCallback(func(hwnd windows.HWND, msg, wParam, lParam uintptr) uintptr {
		switch msg {
		case wmDestroy, wmClose:
			return 0
		}
		r, _, _ := pDefWindowProcW.Call(uintptr(hwnd), msg, wParam, lParam)
		return r
	})
)

var (
	modUser32ForNotify    = windows.NewLazySystemDLL("user32.dll")
	modKernel32ForNotify  = windows.NewLazySystemDLL("kernel32.dll")
	idiApplicationForInfo = uintptr(3251) // IDI_INFORMATION
)

// registerBalloonClass registers the message-only window class once per
// process. It is guarded by a Once because a run may notify twice (the tool
// call and the result), and re-registering the same class name would leak a
// second class.
func registerBalloonClass() error {
	balloonClassOnce.Do(func() {
		inst, _, _ := pGetModuleHandleW.Call(0)
		var wc struct {
			cbSize        uint32
			style         uint32
			lpfnWndProc   uintptr
			cnClsExtra    int32
			cnWndExtra    int32
			hInstance     windows.Handle
			hIcon         windows.Handle
			hCursor       windows.Handle
			hbrBackground windows.Handle
			lpszMenuName  *uint16
			lpszClassName *uint16
			hIconSm       windows.Handle
		}
		wc.cbSize = uint32(unsafe.Sizeof(wc))
		wc.lpfnWndProc = balloonWndProc
		wc.hInstance = windows.Handle(inst)
		name, err := windows.UTF16PtrFromString(balloonTipClass)
		if err != nil {
			balloonClassErr = err
			return
		}
		wc.lpszClassName = name
		icon, _, _ := pLoadIconW.Call(0, idiApplicationForInfo)
		wc.hIcon = windows.Handle(icon)
		r, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		if r == 0 {
			// ERROR_CLASS_ALREADY_EXISTS (1413) is what a second registration
			// attempt returns; under Once it cannot happen, so anything else is
			// a real failure and the caller must see it.
			balloonClassErr = fmt.Errorf("RegisterClassExW(%s): %v", balloonTipClass, err)
		}
	})
	return balloonClassErr
}

// postSystemNotification posts one balloon tip and returns after the shell has
// had it. An error means the notification did not reach the shell.
func postSystemNotification(title, body string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := registerBalloonClass(); err != nil {
		return err
	}
	cls, err := windows.UTF16PtrFromString(balloonTipClass)
	if err != nil {
		return err
	}
	appName, _ := windows.UTF16PtrFromString("wisp run")
	hwnd, _, err := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(cls)),
		uintptr(unsafe.Pointer(appName)), 0, 0, 0, 0, 0,
		uintptr(hwndMessage), 0, 0, 0)
	if hwnd == 0 {
		return fmt.Errorf("CreateWindowExW(message-only): %v", err)
	}
	defer pDestroyWindow.Call(hwnd)

	icon, _, _ := pLoadIconW.Call(0, idiApplicationForInfo)
	d := notifyIconDataW{
		cbSize: uint32(unsafe.Sizeof(notifyIconDataW{})),
		hWnd:   windows.Handle(hwnd),
		uID:    notifyBalloonID,
		// NIF_INFO is the balloon; without NIF_ICON the shell rejects the
		// balloon outright, and without NIF_TIP the tray hover text is empty.
		uFlags:      nifMessage | nifIcon | nifTip | nifInfo,
		hIcon:       windows.Handle(icon),
		dwInfoFlags: niifInfo,
	}
	copyU16(d.szTip[:], "wisp")
	copyU16(d.szInfoTitle[:], title)
	copyU16(d.szInfo[:], body)

	if r, _, e := pShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&d))); r == 0 {
		return fmt.Errorf("Shell_NotifyIconW(NIM_ADD): %v", e)
	}
	// NIM_MODIFY with NIF_INFO is what actually shows the balloon.
	if r, _, e := pShellNotifyIconW.Call(nimModify, uintptr(unsafe.Pointer(&d))); r == 0 {
		pShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&d)))
		return fmt.Errorf("Shell_NotifyIconW(NIM_MODIFY/NIF_INFO): %v", e)
	}
	// Hold the icon for the balloon's own display time, then remove it: a
	// console run must not leave a ghost tray icon behind (the icon is the
	// balloon's owner, so deleting early dismisses it).
	time.Sleep(notifyBalloonTimeout)
	pShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&d)))
	return nil
}

// copyU16 writes s into dst as UTF-16, truncating (never overflowing) rather
// than panicking on a long model reply.
func copyU16(dst []uint16, s string) {
	u, err := windows.UTF16FromString(s)
	if err != nil {
		return
	}
	n := copy(dst, u)
	if n < len(dst) {
		dst[n] = 0
	}
}

var _ = windows.HWND(0)
