//go:build windows && winlive

package main

// Ticket 228 AC#1, live tier: the shipped no-args process really puts a ball
// window on the desktop, and really takes it away again.
//
// The default-tier cases in this family read what the process SAYS (console
// line) and what it BOOKS (persistent log). Both are one indirection away from
// the claim AC#1 makes, which is about a window the user can see:
// "让它真在桌面上出现". Only Win32 can answer that, so this file asks
// FindWindowExW for windows of the ball's own window class and matches each
// one's owning pid against the child it started - the same two facts
// internal/ball's live suite insists on (a window a measurement did not create
// is pollution, and liveness is proven by pid/title, never by class alone).
//
// Run it alone, on a machine with a desktop and no other Wisp ball:
//
//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
//	  go test -tags winlive ./cmd/wisp -run TestLive228 -v
//
// Why "alone" is a rule and not a preference: internal/ball's winlive guard
// (live_guard_windows_test.go, requireQuietBallDesktop, ticket 64 rule 1) fails
// on ANY WispBallWindow it did not create, and this case creates exactly one on
// purpose, in another process. `go test -tags winlive ./...` runs the two
// packages concurrently, so the two live suites may not share a desktop; the
// same constraint already sits between balldebug's interactive modes and that
// guard. The default tier has no such coupling - it opens no window of its own
// when the host cannot (and this file is not built without the tag).
//
// Nothing here skips. A desktop that will not host the ball is a FAIL with the
// refusal the process itself reported, because ticket 228's claim is that the
// ball arrives in the leg the owner launches.

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	live228User32              = windows.NewLazySystemDLL("user32.dll")
	pLive228FindWindowExW      = live228User32.NewProc("FindWindowExW")
	pLive228IsWindow           = live228User32.NewProc("IsWindow")
	pLive228GetWindowThreadPID = live228User32.NewProc("GetWindowThreadProcessId")
)

// ballWindowClass228 is internal/ball's own registration (ball_windows.go
// createOnSTA / registerBallClass). Copied as a literal: when the class changes,
// this file has to be told, because it is the only handle the live tier has on
// "the window is the ball's".
const ballWindowClass228 = "WispBallWindow"

// findBallWindows228 enumerates every top-level window of the ball class and
// returns hwnd -> owning pid. It takes no shortcut through the title: the
// resident leg creates the default title, so pid is the only identity here.
func findBallWindows228(t *testing.T) map[windows.HWND]uint32 {
	t.Helper()
	cls, err := windows.UTF16PtrFromString(ballWindowClass228)
	if err != nil {
		t.Fatalf("UTF16PtrFromString(%q): %v", ballWindowClass228, err)
	}
	out := map[windows.HWND]uint32{}
	var hwnd windows.HWND
	for {
		r, _, _ := pLive228FindWindowExW.Call(0, uintptr(hwnd), uintptr(unsafe.Pointer(cls)), 0)
		if r == 0 {
			break
		}
		hwnd = windows.HWND(r)
		var pid uint32
		pLive228GetWindowThreadPID.Call(r, uintptr(unsafe.Pointer(&pid)))
		out[windows.HWND(r)] = pid
	}
	return out
}

func isWindow228(h windows.HWND) bool {
	r, _, _ := pLive228IsWindow.Call(uintptr(h))
	return r != 0
}

// TestLive228ResidentLegOwnsABallWindowOnTheDesktop is AC#1's live claim, and
// the teardown half of it: after the process leaves through its own shutdown
// path, no window of the ball class is owned by that pid any more.
func TestLive228ResidentLegOwnsABallWindowOnTheDesktop(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()

	// Pollution first: this case identifies its own window by pid, so a stray
	// is not fatal to the reading, but it does say who else is on the desktop -
	// the same disclosure ticket 64 rule 1 added for the ball suite.
	before := findBallWindows228(t)
	for h, pid := range before {
		t.Logf("pre-existing %s window hwnd=%v pid=%d (not created by this case)", ballWindowClass228, h, pid)
	}

	leg := bootResidentLeg(t, exe, dataDir)
	t.Cleanup(leg.stop)

	// Two facts, waited for in the order the boot produces them: the process
	// first says which posture it has (printed after the host returns), and the
	// window exists from inside the host. Polling the window alone would read
	// the console too early and then blame the sentence for a race in here.
	var verdict string
	sawPosture := pollUntil127(200, func() bool {
		verdict = leg.stdout.String()
		return strings.Contains(verdict, ballUpClaim) || strings.Contains(verdict, ballAbsentClaim)
	})
	if !sawPosture {
		t.Fatalf("AC#1 LIVE RED: the process never printed its ball posture (%q / %q).\nstdout:\n%s",
			ballUpClaim, ballAbsentClaim, verdict)
	}
	if strings.Contains(verdict, ballAbsentClaim) {
		// Say which branch happened rather than leaving it to the window list:
		// the process reports its own ball posture, and a refusal is a real
		// reading on a machine with no desktop.
		t.Fatalf("AC#1 LIVE RED: the resident process reported NO ball (%q in its own console). This tier needs "+
			"a desktop; the default tier covers the sentences.\nstdout:\n%s", ballAbsentClaim, verdict)
	}

	var mine []windows.HWND
	saw := pollUntil127(200, func() bool {
		mine = nil
		for h, pid := range findBallWindows228(t) {
			if pid == leg.pid() && isWindow228(h) {
				mine = append(mine, h)
			}
		}
		return len(mine) > 0
	})
	if !saw {
		t.Fatalf("AC#1 LIVE RED: the process claims %q while no %s window is owned by pid %d.\nall ball windows "+
			"seen now: %v\nstdout:\n%s", ballUpClaim, ballWindowClass228, leg.pid(), findBallWindows228(t), verdict)
	}
	if len(mine) != 1 {
		t.Errorf("AC#1 RED: the resident process owns %d ball windows (%v), want exactly 1 (one Ball per process).",
			len(mine), mine)
	}
	hot := ""
	if i := strings.Index(verdict, "hotkeys live"); i >= 0 {
		hot = verdict[i:min(len(verdict), i+80)]
	}
	t.Logf("AC#1 LIVE: hwnd=%v pid=%d owns the ball window; %s", mine[0], leg.pid(), hot)
	// Two counts are findings here, and they are findings about DIFFERENT bugs
	// (ticket 64 A1b and ticket 245, in that order):
	//   0/4 - a ball that registered nothing still looks like a working window,
	//         which is exactly how A1b survived a sign-off run;
	//   4/4 - the idle ball bound the cancel slot too. The production cancel
	//         binding is a bare Esc and RegisterHotKey is desktop-wide, so 4/4
	//         at idle means this process takes Esc from every other application
	//         on the machine while no card is waiting. The steady roster is
	//         three (summon/mute/panel); the fourth slot appears only while
	//         Confirming borrows it, and this leg has no card path at all.
	// The count is read from the number the process prints, which is
	// len(HotkeyReport().Live()) - the registration set, not a sentence about it.
	if strings.Contains(verdict, "hotkeys live 0/4") {
		t.Errorf("AC#1 RED: the ball came up with none of its global hot keys registered.\n%s", verdict)
	}
	if strings.Contains(verdict, "hotkeys live 4/4") {
		t.Errorf("AC#1 RED (ticket 245): the idle resident ball registered all four hot key slots, so the "+
			"cancel slot is a desktop-wide hot key while nothing is being confirmed. The idle roster is three "+
			"(summon/mute/panel); cancel is borrowed only during Confirming.\n%s", verdict)
	}

	// Leave the way an operator does, then require the window to be gone: the
	// D38(e) step 2 work this leg now owns is "hotkey + wake-word listening
	// stops", and a destroyed window is what makes that sentence checkable.
	if err := leg.breakToLoop(); err != nil {
		t.Fatalf("GenerateConsoleCtrlEvent(CTRL_BREAK_EVENT, %d): %v\n%s", leg.pid(), err, leg.console())
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("the child did not exit through its own shutdown path with a live ball window.\n%s", leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("the child exited with %v, want 0.\n%s", leg.waitErr, leg.console())
	}
	for _, h := range mine {
		if isWindow228(h) {
			t.Errorf("AC#1 RED: hwnd %v is still a window after the process that owned it left through D38(e).", h)
		}
	}
	for h, pid := range findBallWindows228(t) {
		if pid == leg.pid() {
			t.Errorf("AC#1 RED: %s window %v is still attributed to the exited pid %d", ballWindowClass228, h, pid)
		}
	}
	fmt.Printf("live228: pid %d created and destroyed its %s window through the D38(e) order\n",
		leg.pid(), ballWindowClass228)
}
