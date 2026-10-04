//go:build windows && winlive

package ball

// Ticket 260 AC#1 on the REAL desktop: does the configured [hotkey] cancel
// binding actually become the hot key Win32 holds for the length of the card,
// and does the default still land on bare Esc?
//
//	WRITTEN, NOT RUN by 260-r1. The instruction for this round forbids taking the
//	owner's desktop and real hotkeys; these two cases do exactly that for a few
//	seconds each, so they are owed to the winlive window (named in
//	.scratch/wisp/probes/260/r1/landing.md as 〔未量有没有牙〕归 260-v1). Run with:
//
//	go test -tags winlive ./internal/ball/ -run "260" -v
//
// requireEscBorrowed (hotkey_live_test.go:90) is deliberately NOT reused by the
// configured case: it asserts the borrowed slot is VK_ESCAPE, which is true of the
// DEFAULT config and false of any other value - the premise ticket 260 moved. The
// "real Win32 holds it" check is done here against the configured combination
// itself, on a foreign id, the same way escBorrowProbe checks bare Esc.

import (
	"errors"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
	"golang.org/x/sys/windows"
)

// TestLiveCancelBorrowDefaultUnchanged260 is condition ① at the Win32 layer: a
// ball booted on the product defaults borrows the same bare Esc it borrowed
// before this ticket, and hands it back.
func TestLiveCancelBorrowDefaultUnchanged260(t *testing.T) {
	b := newLiveBall(t, Options{Initial: statemachine.StateSleeping, Hotkeys: DefaultHotkeys()})
	requireIdleRoster(t, b)

	b.TakeEscForCancel()
	if !b.EscTakenOver() {
		t.Fatal("the default borrow did not mark the slot taken over")
	}
	got := b.RegisteredHotkeys()[hkCancel]
	if got != accelDefault260 {
		t.Fatalf("Win32 holds %+v for the cancel slot, want the default %+v", got, accelDefault260)
	}
	line, _ := b.HotkeyReport().Binding(hkCancel)
	if line.Binding != "Esc" || line.Acc != accelDefault260 || line.Note != "" {
		t.Fatalf("the receipt line for a default borrow = %+v, want live \"Esc\" / %+v / no note", line, accelDefault260)
	}
	// The desktop, not our bookkeeping, says the key is ours.
	requireEscBorrowed(t, b)

	b.ReleaseEscAfterSession()
	requireEscReturned(t, b)
}

// TestLiveCancelBorrowFollowsConfig260 is AC#1's positive control at the Win32
// layer: the key the config names is the key that really registers.
func TestLiveCancelBorrowFollowsConfig260(t *testing.T) {
	cfg := ApplyHotkeyDefaults(HotkeyConfig{
		Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Ctrl+Alt+Y", Panel: "Ctrl+Alt+B",
	})
	b := newLiveBall(t, Options{Initial: statemachine.StateSleeping, Hotkeys: cfg})
	requireIdleRoster(t, b)

	borrow := resolveCancelBorrow(cfg.Cancel)
	if borrow.Fallback != nil {
		t.Fatalf("the fixture binding is not parsable, which makes this case measure nothing: %v", borrow.Fallback)
	}

	b.TakeEscForCancel()
	if !b.EscTakenOver() {
		t.Fatal("a configured cancel binding never took the slot over")
	}
	got := b.RegisteredHotkeys()[hkCancel]
	if got != borrow.Acc {
		t.Fatalf("the report claims %+v while the resolved borrow is %+v - config did not reach the registration", got, borrow.Acc)
	}
	if got.VK == vkEscape {
		t.Fatalf("the borrow took the hard-coded bare Esc instead of the configured %q (%+v)", cfg.Cancel, got)
	}
	line, _ := b.HotkeyReport().Binding(hkCancel)
	if line.Binding != cfg.Cancel || line.Acc != borrow.Acc || line.Status != HotkeyLive || line.Note != "" {
		t.Fatalf("the receipt line = %+v, want live %q with %+v and no note", line, cfg.Cancel, borrow.Acc)
	}

	// Ask Win32 whether that exact combination is taken. On a foreign id, the same
	// window, a successful registration would mean our own borrow is not what the
	// report says it is; ERROR_HOTKEY_ALREADY_REGISTERED means it is.
	err := uiRegisterHotkey(t, b, spareHKID, borrow.Acc)
	if err == nil {
		uiUnregisterHotkey(t, b, spareHKID)
		t.Fatalf("Win32 accepted a second registration of %+v: the ball's cancel slot never bound the configured key", borrow.Acc)
	}
	if !errors.Is(err, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
		t.Fatalf("the probe for %+v failed for a reason other than 1409: %v", borrow.Acc, err)
	}

	// While the card is up, a bare-Esc probe must find Esc FREE again - the whole
	// point of 形ⓐ is that a user who moved off Esc does not also lose Esc.
	free, perr := escBorrowProbe(t, b)
	if perr != nil {
		t.Fatalf("the bare-Esc probe could not run: %v", perr)
	}
	if !free {
		t.Errorf("bare Esc is still claimed while [hotkey] cancel says %q - the borrow is holding two keys", cfg.Cancel)
	}

	b.ReleaseEscAfterSession()
	requireEscReturned(t, b)
}
