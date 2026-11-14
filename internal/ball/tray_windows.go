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

	// Ticket 228 AC#2: the three ANSWER items. They are appended, never
	// renumbered - ids 1..4 and the four labels below them are byte-identical
	// to what shipped before this ticket, because a renumbering would silently
	// re-bind a menu entry a consumer already dispatches on.
	//
	// What these three are, and what they are not: they carry the human's answer
	// to the ONE card the process is being waited on (Replies.Allow /
	// Replies.Reject / replySurface.always on the host side). internal/ball knows
	// nothing about approvals and must keep knowing nothing - the callbacks are
	// func(), fire() only tests them for nil (ball_windows.go's fire), and every
	// sentence, refusal and audit line is the host's. 「本次会话内允许」 is
	// deliberately NOT here: the Answer verb table has four words and the
	// session-scoped grant is ticket 224's, not this ticket's (AC#2's own text).
	menuAllowOnce   = 5
	menuReject      = 6
	menuAllowAlways = 7
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
	// Ticket 228 AC#2: the answer group. It is ALWAYS built, never hidden and
	// never greyed - there is no "is a card waiting" input in this function, and
	// adding one would put the ledger's depth into a package that has no business
	// knowing it. The host's ruling on the落空 shape is "always show, say so when
	// the press finds nothing" (ticket 228 AC#2 supplement, orchestrator F4): the
	// item is present, and a press with no pending card produces the host's named
	// no-pending-card line instead of a silent no-op or a fake success.
	appendItem(menuAllowOnce, "允许一次", false)
	appendItem(menuReject, "拒绝", false)
	appendItem(menuAllowAlways, "长期允许", false)
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
