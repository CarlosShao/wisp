//go:build windows

package main

// Ticket 245 AC#4: "用户要能改" - the [hotkey] cancel line the user writes in
// config.toml is the key the resident leg hands the ball.
//
// WHY THIS PIN EXISTS SEPARATELY. Ticket 258 wired the whole [hotkey] section
// into the resident leg (residentBallHotkeyChain258) and its positive control
// (Test258ConstructionChainTakesConfigValues) plants non-default summon / mute /
// panel values while leaving cancel at "Esc" - deliberately, because that file's
// job was the chain, not this slot. So the one slot ticket 245 is about had no
// pin of its own: a chain that quietly dropped the user's cancel line and kept
// the shipped default would have passed every case in the tree.
//
// WHAT THIS PIN DOES NOT COVER, said so nobody reads it wider: it stops at the
// ball's config (the value the idle standby line carries and the borrow
// registers). The borrow itself - "the key that really reaches Win32 is the
// configured one" - is pinned in internal/ball by ticket 260's
// TestConfiguredCancelStillNeverBoundWhileIdle260 / resolveCancelBorrow cases,
// because resolveCancelBorrow is not reachable from this package. And the real
// desktop reading of the shipped binary belongs to the winlive tier
// (resident_hotkey_live_258_windows_test.go), which this round did not run.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
)

// Test245AC4ResidentChainReadsAConfiguredCancel is AC#4 at the seam this
// package owns: config.toml -> config.LoadFile -> the host's [hotkey] view ->
// residentBallHotkeyChain258 -> the ball.HotkeyConfig that ball.New receives.
func Test245AC4ResidentChainReadsAConfiguredCancel(t *testing.T) {
	dir := t.TempDir()
	writeConfig258(t, dir, `summon = "Ctrl+Alt+Q"
mute = "Ctrl+Alt+M"
cancel = "Ctrl+Alt+V"
panel = "Ctrl+Alt+P"`)
	hotCfg := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }

	// Fixture reading, so the case says what it measured: the file really does
	// carry a cancel value that is NOT the shipped default. Without this line a
	// green could mean "the user did not change anything".
	seen := hotCfg()
	if seen.Cancel != "Ctrl+Alt+V" {
		t.Fatalf("fixture reading: the [hotkey] view answered cancel = %q, want the planted value", seen.Cancel)
	}
	if seen.Cancel == ball.DefaultHotkeys().Cancel {
		t.Fatalf("fixture is measuring nothing: the planted cancel equals the shipped default %q", seen.Cancel)
	}

	got, provenance := residentBallHotkeyChain258(hotCfg)
	if got.Cancel != "Ctrl+Alt+V" {
		t.Fatalf("the resident chain handed the ball cancel = %q, want the user's planted value: a hand edit "+
			"of [hotkey] cancel must reach the slot ticket 245 borrows (AC#4)", got.Cancel)
	}
	if provenance != hotkeyProvenanceConfig {
		t.Fatalf("provenance = %q, want %q: the set the ball boots with carries the user's cancel line",
			provenance, hotkeyProvenanceConfig)
	}
	// The other three slots are the file's too - a chain that answered the
	// compiled table for them would be a different defect, and cheap to name here.
	if got.Summon != "Ctrl+Alt+Q" || got.Mute != "Ctrl+Alt+M" || got.Panel != "Ctrl+Alt+P" {
		t.Fatalf("chain answered %+v, want the planted summon/mute/panel too", got)
	}

	// Positive control in the same shape: delete ONLY the user's cancel line and
	// the same chain must answer the shipped default again. A chain that ignored
	// its input (or a fixture that always says "Ctrl+Alt+V") fails one of the two
	// halves.
	dirDefault := t.TempDir()
	writeConfig258(t, dirDefault, `summon = "Ctrl+Alt+Q"
mute = "Ctrl+Alt+M"
panel = "Ctrl+Alt+P"`)
	def, defProvenance := residentBallHotkeyChain258(func() ball.HotkeyConfig {
		return loadHotkey258(t, dirDefault)
	})
	if want := ball.DefaultHotkeys(); def != want {
		t.Fatalf("control: without a cancel line the chain answered %+v, want the shipped table %+v", def, want)
	}
	if defProvenance != hotkeyProvenanceDefaults {
		t.Fatalf("control: provenance = %q, want %q (no slot was the file's)", defProvenance, hotkeyProvenanceDefaults)
	}
}
