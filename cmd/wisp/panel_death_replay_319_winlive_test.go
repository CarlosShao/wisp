//go:build windows && winlive

package main

// Ticket 319 AC#1/AC#3 (N3), the LIVE half: the census's death chain replayed on
// a REAL window, then the next gesture.
//
// What 319-a1 established from the owner's log (census section 1, and the
// BEFORE half of this file's first case replays it): at 15:41:45 the panel
// window received WM_SYSCOMMAND(0x112) + WM_CLOSE(0x10) in the same
// millisecond, went through WM_DESTROY(0x2), the library pump returned, and the
// thread ended. After that every gesture was refused with "panel gesture could
// not be handed over". This file plants that exact message pair on a real
// window this test process created (PostMessageW(WM_SYSCOMMAND, SC_CLOSE), the
// same pair the system sends for the X button), waits for finished to close,
// and then reads the two answers that matter:
//
//  1. BEFORE shape, verbatim: the next gesture handed to the dead thread is
//     refused and no window comes back. This assertion is deliberate and stays
//     true after the fix - post() must keep refusing a finished thread, which
//     is what makes the refusal honest; the fix is at the holder, not here.
//  2. AFTER shape (second case, the AC#1 judgement read live): the SAME
//     gesture through the assembly root's panelHolder asks for a fresh thread
//     on the same manager and a NEW window really comes up. The new HWND is
//     reported; a rebuild that reused the old handle would read as a failure,
//     not as a pass.
//
// The window is this test process's own (its own WebView2 profile under
// /tmp), and every event name any case in this tree opens is unique per run -
// the owner's running wisp.exe is never touched. winlive has no CI job: this
// runs on the desktop box only, one case at a time, by name.

import (
	"testing"
)

// The two message numbers the close chain starts with (winuser.h), named so the
// plant reads as itself.
const (
	pan319WmSysCommand = 0x112
	pan319ScClose      = 0xF060
)

// TestPanelDeathReplay319UnfixedReading replays the 15:41:45 chain on a real
// window and records the shape the ticket was filed about: dead thread, refused
// gesture, no window. It is the live "今天不响" half of N3's two-way judgement.
func TestPanelDeathReplay319UnfixedReading(t *testing.T) {
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)
	old := mgr.windowHandle()
	if old == 0 {
		t.Fatalf("the fixture reported a shown window with HWND 0")
	}

	// The system's own close, programmatically: the same (0x112, SC_CLOSE) pair
	// the 15:41:45 log line carries.
	r, _, callErr := t33r9NailPostMessage.Call(old, pan319WmSysCommand, pan319ScClose, 0)
	if r == 0 {
		t.Fatalf("PostMessageW(hwnd 0x%X, WM_SYSCOMMAND, SC_CLOSE) was refused: %v", old, callErr)
	}
	waitPanelTrue(t, "the panel thread to end after its window was closed", rp.isFinished)

	// BEFORE shape, verbatim: the same executor the ball, the hot key and the
	// tray item drive, handed to the dead thread.
	toggle := rp.RequestToggle("319-death-replay")
	if toggle {
		t.Errorf("a gesture was accepted by a panel thread whose pump had returned - the refusal is the honest half of the defect, and this case pins it")
	}
	if mgr.IsCreated() || mgr.windowHandle() != 0 {
		t.Errorf("a window is recorded after the close but before any rebuild: IsCreated=%v hwnd=0x%X", mgr.IsCreated(), mgr.windowHandle())
	}
	t.Logf("319 unfixed live reading: hwnd=0x%X closed via WM_SYSCOMMAND/SC_CLOSE; thread finished=%v; next gesture RequestToggle=%v; mgr.IsCreated=%v; mgr hwnd=0x%X",
		old, rp.isFinished(), toggle, mgr.IsCreated(), mgr.windowHandle())
}
