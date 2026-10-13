//go:build windows

package main

// Ticket 294 AC#2/AC#3: the two readings that need real Win32 or a real gate.
//
// AC#2: on the real ball window (the startResidentBall shape ticket 258's suite
// already drives) the MUTE hot key's own registration row must be live. That
// reading never existed: hkMute is internal/ball's unexported constant and
// cmd/wisp only ever mirrored summon(1)/panel(4), so nobody read IsLive for the
// mute id from this side.
//
// AC#3: the other half of toggleMute's (string, bool) contract. Ticket 290
// pinned the executed=false shape (a leg with no gate may not claim it turned
// one); the happy path's executed=true was asserted by nobody, so mutating the
// gate-was-turned returns to ", false" stayed green (ticket 294's M2 reading).
// The judge here is the gate's own readback plus this bool - never the outcome
// sentence (ticket 290's M1 proved a sentence can say success while the bool
// says otherwise).

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/observe"
)

// hkMute294 mirrors internal/ball/hotkey_windows.go's hkMute. It is the second
// hand-copied truth source in this package (resident_hotkey_258_windows_test.go
// carries hkSummon258=1 and hkPanel258=4; mute=2 had no mirror at all).
//
// If this mirror is copied wrong, the IsLive read below answers about the wrong
// slot: 1 (summon) would read "live" for a mute registration that never
// happened, 3 (cancel) is the idle-set standby and would go red for the wrong
// reason, and 4 (panel) points at a different key entirely. So the mirror is
// not trusted on its name: TestAC294HkMuteMirrorMatchesInternalBallSource
// re-reads the constant from internal/ball's source and requires equality.
const hkMute294 uint32 = 2

// TestAC294HkMuteMirrorMatchesInternalBallSource is the same-source ruler for
// the mirror above: parse internal/ball/hotkey_windows.go, find the integer
// const hkMute, require it to equal hkMute294. A hand-copied constant that
// drifts from its source is a silent wrong reading in every test below.
func TestAC294HkMuteMirrorMatchesInternalBallSource(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join("..", "..", "internal", "ball", "hotkey_windows.go"), nil, 0)
	if err != nil {
		t.Fatalf("parsing internal/ball/hotkey_windows.go (the source of truth, not the code): %v", err)
	}
	srcVal := int64(-1)
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				if name.Name != "hkMute" || i >= len(vs.Values) {
					continue
				}
				if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.INT {
					if v, err := strconv.ParseInt(lit.Value, 0, 64); err == nil {
						srcVal = v
					}
				}
			}
		}
	}
	if srcVal < 0 {
		t.Fatalf("no integer const hkMute found in internal/ball/hotkey_windows.go: the id table moved, re-adjudicate this row")
	}
	if uint32(srcVal) != hkMute294 {
		t.Fatalf("mirror drift: cmd/wisp hkMute294 = %d but internal/ball hkMute = %d; every IsLive(hkMute294) read in this file would be about the wrong slot", hkMute294, srcVal)
	}
}

// TestAC294MuteHotkeyRegistersLiveOnARealBallWindow is AC#2: bring up the real
// ball window exactly the way the 258 suite does, plant a NON-default mute
// binding, and read the mute id's own registration row off HotkeyReport.
// Ticket 290 proved the flip once the gesture asks; this is the other
// direction - the user's key must actually reach production registration.
//
// SKIP-LOUD is for headless hosts only; on a host with a desktop this case
// must produce the reading, not the skip (ticket 290 AC#3's precedent: a skip
// is not a reading).
func TestAC294MuteHotkeyRegistersLiveOnARealBallWindow(t *testing.T) {
	dir := t.TempDir()
	writeConfig258(t, dir, `summon = "Ctrl+Alt+Q"
mute = "Ctrl+Alt+X"
cancel = "Esc"
panel = "Ctrl+Alt+P"`)
	src := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }
	rb := startResidentBall(observe.NewRegistry(), nil, src, src)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the mute registration has no real Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()

	boot := rb.b.HotkeyReport()
	if !boot.IsLive(hkMute294) {
		t.Fatalf("AC#2 RED: the mute hot key is not live on the real window: %+v", boot.Bindings())
	}
	if got := boot.Live()[hkMute294].VK; got != 'X' {
		t.Fatalf("AC#2 RED: mute VK = 0x%X, want 'X' (the planted non-default; 'M' here would mean the mute slot fell back to the compiled default)", got)
	}
	// The reading itself, printed so the evidence carries the real row and not
	// just the fact this case passed.
	row, ok := boot.Binding(hkMute294)
	if !ok {
		t.Fatalf("AC#2 RED: no mute row in the boot report: %+v", boot.Bindings())
	}
	t.Logf("AC#2 reading: id=%d name=%q binding=%q status=%q vk=0x%X", row.ID, row.Name, row.Binding, row.Status, row.Acc.VK)
}

// TestAC294ToggleMuteExecutedTrueOnTheHappyPath is AC#3: the gate-was-turned
// path must answer executed=true. The rig is ticket 290's own (a real
// HalfDuplexGate over the capture seam), and the assertions are the gate's
// readback plus the bool. The mutating control (the gate.Open() branch's
// return flipped from ", true" to ", false") is in the leg's probes: it
// reddens this case while leaving the outcome sentence intact - which is
// exactly why the sentence may not be the judge.
func TestAC294ToggleMuteExecutedTrueOnTheHappyPath(t *testing.T) {
	src := &openableSource{}
	ra, rb, _ := muteRig(t, src)

	// Premise said out loud: the shipped posture boots muted with no device
	// handed over, so the first press is an OPEN the fixture really performs.
	if !ra.gate.Muted() || ra.gate.Open() {
		t.Fatalf("premise moved: the rig did not boot in the shipped muted posture (muted=%v open=%v)", ra.gate.Muted(), ra.gate.Open())
	}

	outcome, executed := ra.toggleMute()
	if !executed {
		t.Fatalf("AC#3 RED: the happy path answered executed=false while the gate DID turn (muted=%v open=%v): %q", ra.gate.Muted(), ra.gate.Open(), outcome)
	}
	if !ra.gate.Open() || ra.gate.Muted() {
		t.Fatalf("the gate did not really open: muted=%v open=%v (%q)", ra.gate.Muted(), ra.gate.Open(), outcome)
	}

	// And back: the same executor mutes again, still executed=true.
	outcome, executed = ra.toggleMute()
	if !executed {
		t.Fatalf("AC#3 RED: the mute-back press answered executed=false while the gate DID turn (muted=%v open=%v): %q", ra.gate.Muted(), ra.gate.Open(), outcome)
	}
	if !ra.gate.Muted() || ra.gate.Open() {
		t.Fatalf("the gate did not really close: muted=%v open=%v (%q)", ra.gate.Muted(), ra.gate.Open(), outcome)
	}

	// The host entry point both user gestures use lands on the same executor
	// (ticket 294 AC#1 pins where production attaches it; this is the turn
	// itself). Judge: the gate's readback, not the sentence muteGesture
	// returns.
	rb.muteGesture("mute-hotkey")
	if !ra.gate.Open() || ra.gate.Muted() {
		t.Fatalf("the gesture did not reach the gate: muted=%v open=%v", ra.gate.Muted(), ra.gate.Open())
	}
}
