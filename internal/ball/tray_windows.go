//go:build windows

package ball

// Tray icon + menu (SPEC-08 §7): left click = open panel (no-op stub until
// ticket 33), right click menu = 打开面板 / 静音 / 暂停唤醒 / 退出.

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Tray menu command ids (wParam of WM_COMMAND after TPM_RETURNCMD is NOT
// used; we use TrackPopupMenu's return value inline in the wndproc instead).
const (
	trayUID = 0x5701

	menuOpenPanel = 1
	menuMute      = 2
	menuPauseWake = 3
	menuExit      = 4
)

type tray struct {
	data notifyIconData
}

func addTrayIcon(hwnd windows.HWND, tip string) (*tray, error) {
	icon, _, _ := pLoadIconW.Call(0, idiApplication)
	d := notifyIconData{
		cbSize:           uint32(unsafe.Sizeof(notifyIconData{})),
		hWnd:             hwnd,
		uID:              trayUID,
		uFlags:           nifMessage | nifIcon | nifTip,
		uCallbackMessage: wmAppTray,
		hIcon:            windows.Handle(icon),
	}
	copy(d.szTip[:], toU16(tip))
	t := &tray{data: d}
	r, _, err := pShellNotifyIconW.Call(nimAdd, unsafePtr(&t.data))
	if r == 0 {
		return nil, fmt.Errorf("ball: Shell_NotifyIconW(NIM_ADD): %v", err)
	}
	return t, nil
}

func toU16(s string) []uint16 {
	u, _ := windows.UTF16FromString(s)
	return u
}

// setTip updates the tray tooltip (used to surface the current state).
func (t *tray) setTip(tip string) {
	u := toU16(tip)
	copy(t.data.szTip[:], u)
	t.data.uFlags |= nifTip
	pShellNotifyIconW.Call(nimModify, unsafePtr(&t.data))
}

func (t *tray) remove() {
	if t == nil {
		return
	}
	pShellNotifyIconW.Call(nimDelete, unsafePtr(&t.data))
}

// showMenu builds and tracks the right-click menu. Runs on the STA thread
// inside the tray message handling. Returns the selected command id (0 =
// dismissed).
func showMenu(hwnd windows.HWND, muted, pausedWake bool) uint32 {
	menu, _, _ := pCreatePopupMenu.Call()
	if menu == 0 {
		return 0
	}
	defer pDestroyMenu.Call(menu)

	appendItem := func(id uintptr, label string, checked bool) {
		flags := uintptr(mfString)
		if checked {
			flags |= mfCheckd
		}
		pAppendMenuW.Call(menu, flags, id, unsafePtr(utf16(label)))
	}
	appendItem(menuOpenPanel, "打开面板", false)
	pAppendMenuW.Call(menu, mfSepart, 0, 0)
	appendItem(menuMute, "静音", muted)
	appendItem(menuPauseWake, "暂停唤醒", pausedWake)
	pAppendMenuW.Call(menu, mfSepart, 0, 0)
	appendItem(menuExit, "退出", false)

	// Tray menu focus quirk: without SetForegroundWindow the menu refuses
	// to dismiss on outside clicks. The user explicitly interacted with the
	// tray, so this activation is theirs (the ball window itself still
	// never activates on its own).
	pSetForegroundWindow.Call(uintptr(hwnd))

	var pt point
	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	sel, _, _ := pTrackPopupMenu.Call(menu, tpmRightButton|tpmReturnCMD|tpmNoActivate,
		uintptr(pt.x), uintptr(pt.y), 0, uintptr(hwnd), 0)
	// Dismiss-restore quirk: a benign WM_NULL makes the next outside click
	// close the menu cleanly.
	pPostMessageW.Call(uintptr(hwnd), wmNull, 0, 0)
	return uint32(sel)
}

var pGetCursorPos = modUser32.NewProc("GetCursorPos")
