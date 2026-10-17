//go:build windows && winlive

package main

// Ticket 228 AC#11 ⓑ-2, live tier: the SHIPPED no-args process, a real ball
// window and a real tray icon on a real desktop, and the tray branch's request
// carrying it out with exit code 0.
//
// This is the only case in the family allowed to speak about a request reaching a
// live process, and the only thing it claims is the observable half:
//
//   - that one request ends the process with exit code 0 (not 0xc000013a, which
//     is what "nothing was registered for the signal" looks like on this machine,
//     and not a kill);
//   - the trail on disk still holds the sink's own install booking, the tray
//     booking, the ball teardown and the D38(e) records after them;
//   - the window that belonged to that pid is gone.
//
// ⛔ It does NOT read a 1..10 step sequence off the file, and it may not:
// internal/proc/shutdown.go books a line for skipped, failed and fast steps only
// (shutdown.go's four logging outlets), so a step that ran cleanly is structurally
// absent from the disk, and adding per-step logs would move D38(e)'s observation
// surface - ledger A495 names exactly this, and the in-process reading of all ten
// records is TestAC11TrayExitTriggerWalksTheTenSteps in
// resident_tray_stop_228_windows_test.go.
//
// HOW the tray branch is driven here, stated plainly because the honest version is
// narrower than "a user clicked": internal/ball's tray branch runs its menu on the
// ui-sta thread and takes the selection ONLY from TrackPopupMenu's synchronous
// return (internal/ball/tray_windows.go's showMenu; the same file's header says
// WM_COMMAND is NOT used, and ledger A495 closes the "post WM_COMMAND with the
// menu id" shortcut by name - no case here may rebuild it). So this case hands the
// child's own ball window the tray callback message the shell would hand it, then
// completes the menu with REAL injected mouse input on the notification-area popup
// that the child itself opened, and the selection is the value TrackPopupMenu
// returns. What is NOT exercised is the hit-test of the icon in the notification
// area itself: that belongs to the shell, has no handle this process owns, and the
// gap stays named here rather than claimed.
//
// Run it alone on a machine with a desktop, and never beside internal/ball's own
// live suite (its requireQuietBall guard fails on any WispBallWindow it did not
// create, and this case creates one on purpose in another process - the same
// constraint resident_ball_live_228_windows_test.go's header states):
//
//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
//	  go test -tags winlive ./cmd/wisp -run 'TestLive228R3' -v
//
// The owner's own live instance is never touched: every handle here is claimed by
// the pid of the child this case started.

import (
	"strings"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// live228R3TrayMsg is WM_APP + 0x202, the callback message the tray icon was
	// registered with (internal/ball/win32_windows.go's wmAppTray). Copied as a
	// literal because internal/ball keeps it unexported; a rename there has to
	// tell this file, which is the same rule that keeps ballWindowClass228 a
	// literal.
	live228R3TrayMsg  = 0x0400 + 0x202
	live228R3TrayUID  = 0x5701 // internal/ball/tray_windows.go's trayUID
	live228R3RBUp     = 0x0205 // WM_RBUTTONUP, the low word of lParam the branch reads
	live228R3MenuCls  = "#32768"
	live228R3LastItem = 12 * time.Millisecond // one item row, SM_CYSMSIZE's neighbourhood
)

var (
	live228R3User32         = windows.NewLazySystemDLL("user32.dll")
	pLive228R3PostMessageW  = live228R3User32.NewProc("PostMessageW")
	pLive228R3GetWindowRect = live228R3User32.NewProc("GetWindowRect")
	pLive228R3SetCursorPos  = live228R3User32.NewProc("SetCursorPos")
	pLive228R3GetCursorPos  = live228R3User32.NewProc("GetCursorPos")
	pLive228R3MouseEvent    = live228R3User32.NewProc("mouse_event")
)

const (
	live228R3LeftDown = 0x0002
	live228R3LeftUp   = 0x0004
)

type (
	live228R3Point struct{ x, y int32 }
	live228R3Rect  struct{ left, top, right, bottom int32 }
)

// live228R3MenuWindow finds a popup menu window owned by pid. A menu is a
// top-level window of the system class #32768; several may exist on a desktop, so
// the owning pid is the identity - the same rule the ball window claim follows.
func live228R3MenuWindow(t *testing.T, pid uint32) windows.HWND {
	t.Helper()
	cls, err := windows.UTF16PtrFromString(live228R3MenuCls)
	if err != nil {
		t.Fatalf("UTF16PtrFromString(%q): %v", live228R3MenuCls, err)
	}
	var hwnd windows.HWND
	for {
		r, _, _ := pLive228FindWindowExW.Call(0, uintptr(hwnd), uintptr(unsafe.Pointer(cls)), 0)
		if r == 0 {
			return 0
		}
		hwnd = windows.HWND(r)
		var owner uint32
		pLive228GetWindowThreadPID.Call(r, uintptr(unsafe.Pointer(&owner)))
		if owner == pid && isWindow228(hwnd) {
			return hwnd
		}
	}
}

func live228R3RectOf(h windows.HWND) (live228R3Rect, bool) {
	var rc live228R3Rect
	r, _, _ := pLive228R3GetWindowRect.Call(uintptr(h), uintptr(unsafe.Pointer(&rc)))
	return rc, r != 0
}

func live228R3Click(x, y int32) {
	pLive228R3SetCursorPos.Call(uintptr(x), uintptr(y))
	time.Sleep(live228R3LastItem)
	pLive228R3MouseEvent.Call(live228R3LeftDown, 0, 0, 0, 0)
	pLive228R3MouseEvent.Call(live228R3LeftUp, 0, 0, 0, 0)
}

// TestLive228R3TrayExitStopsTheShippedProcessWithCodeZero is AC#11 ⓑ-2.
func TestLive228R3TrayExitStopsTheShippedProcessWithCodeZero(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	sinkDir := logSinkDir(dataDir)

	// The cursor is the owner's, not this case's: saved first, restored first out.
	var restore live228R3Point
	pLive228R3GetCursorPos.Call(uintptr(unsafe.Pointer(&restore)))
	t.Cleanup(func() { pLive228R3SetCursorPos.Call(uintptr(restore.x), uintptr(restore.y)) })

	// 1. the leg must be the one that hosts a ball, or there is no tray to ask.
	if !pollUntil127(400, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
	}) {
		t.Fatalf("AC#11 ⓑ-2 RED: the shipped process never printed its ball posture.\n%s", leg.console())
	}
	if strings.Contains(leg.stdout.String(), ballAbsentClaim) {
		t.Fatalf("AC#11 ⓑ-2 RED: the resident process reported NO ball on this desktop, so it owns no tray icon "+
			"to ask anything of.\n%s", leg.console())
	}

	// 2. claim its window by pid - the only handle this case may address.
	var hwnd windows.HWND
	if !pollUntil127(200, func() bool {
		for h, pid := range findBallWindows228(t) {
			if pid == leg.pid() {
				hwnd = h
				return true
			}
		}
		return false
	}) {
		t.Fatalf("AC#11 ⓑ-2 RED: no %s window is owned by pid %d although the process claims a ball.\n%s",
			ballWindowClass228, leg.pid(), leg.console())
	}

	// 3. park the cursor somewhere a popup fits, then hand the child the tray
	//    callback the shell would hand it on a right button up.
	pLive228R3SetCursorPos.Call(400, 300)
	time.Sleep(50 * time.Millisecond)
	r, _, err := pLive228R3PostMessageW.Call(uintptr(hwnd), live228R3TrayMsg, live228R3TrayUID, live228R3RBUp)
	if r == 0 {
		t.Fatalf("AC#11 ⓑ-2 RED: PostMessageW(hwnd %v, WM_APP+0x202, %d, WM_RBUTTONUP) failed: %v",
			hwnd, live228R3TrayUID, err)
	}

	// 4. the menu the CHILD opened, on the desktop, addressed by its pid.
	var menu windows.HWND
	if !pollUntil127(200, func() bool {
		menu = live228R3MenuWindow(t, leg.pid())
		return menu != 0
	}) {
		t.Fatalf("AC#11 ⓑ-2 RED: no popup menu window is owned by pid %d after the tray callback message; the "+
			"break is between the message and internal/ball's showMenu on the ui-sta thread.\nhwnd=%v\n%s",
			leg.pid(), hwnd, leg.console())
	}
	rc, ok := live228R3RectOf(menu)
	if !ok {
		t.Fatalf("AC#11 ⓑ-2 RED: GetWindowRect on the menu window %v failed", menu)
	}
	t.Logf("AC#11 ⓑ-2: child pid=%d hwnd=%v menu=%v rect=(%d,%d)-(%d,%d)",
		leg.pid(), hwnd, menu, rc.left, rc.top, rc.right, rc.bottom)

	// 5. pick the LAST row, which is 退出 (menuExit is the sixth entry of
	//    showMenu's own build order). TrackPopupMenu returns that id synchronously
	//    and internal/ball's branch fires Events.OnTrayExit - the callback this
	//    ticket gave an executor.
	live228R3Click((rc.left+rc.right)/2, rc.bottom-6)

	// 6. the process must leave by itself, through its own ordered path.
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("AC#11 ⓑ-2 RED: the tray request did not end the process; it is still running, so the stop bit "+
			"or the loop's check of it is not reached from the tray branch.\nmenu=%v\n%s", menu, leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("AC#11 ⓑ-2 RED: the child ended as %s, want exit code 0. 0xc000013a here would mean the request "+
			"never reached the loop at all.\n%s", leg.exitStatus(), leg.console())
	}

	out := leg.stdout.String()
	for _, want := range []string{
		"tray exit requested: recorded",
		"event loop ending (stop request (tray Exit item selected in this process))",
		residentShutdownLine,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("AC#11 ⓑ-2 RED: the console of the process that really left through the tray branch lacks %q.\n%s", want, out)
		}
	}

	// 7. the trail, read from the file that process closed. install + tray booking +
	//    ball teardown + at least one D38(e) record after them. ⛔ no 1..10 claim.
	recs := readResidentSink(t, sinkDir)
	installIdx := installRecordIndex127(recs)
	if installIdx < 0 {
		t.Fatalf("AC#11 ⓑ-2 RED: no %q record in the file the tray-stopped process wrote (msgs: %v)",
			residentInstallMsg, msgsOf127(recs))
	}
	requested := firstIndexOfContains228(recs, "tray exit requested")
	stopped := indexOfMsg228(recs, ballStoppedMsg)
	shutdown := firstIndexOfContains228(recs, residentShutdownRecord)
	for _, pair := range []struct {
		what string
		idx  int
	}{
		{"tray exit booking", requested},
		{"ball teardown booking", stopped},
		{"D38(e) record", shutdown},
	} {
		if pair.idx < 0 {
			t.Errorf("AC#11 ⓑ-2 RED: the audit trail of a process stopped by the tray branch holds no %s (msgs: %v)", pair.what, msgsOf127(recs))
			continue
		}
		if pair.idx < installIdx {
			t.Errorf("AC#11 ⓑ-2 RED: the %s at record %d precedes the sink's own install booking at %d", pair.what, pair.idx, installIdx)
		}
	}
	if isWindow228(hwnd) {
		t.Errorf("AC#11 ⓑ-2 RED: hwnd %v is still a window after the process that owned it left through the tray branch.", hwnd)
	}
	t.Logf("AC#11 ⓑ-2 READING: exit code 0; install=%d tray=%d stopped=%d shutdown=%d of %d record(s)",
		installIdx, requested, stopped, shutdown, len(recs))
}
