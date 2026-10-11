//go:build windows

package main

// Ticket 296 AC#1: the world ticket 258's suite never built - a config.toml
// whose [hotkey] section answers THREE EMPTY SLOTS, which is what "the user
// never wrote those keys" looks like in this file format.
//
// WHY THIS FILE EXISTS. The owner's own file (read verbatim at 2026-10-09
// 17:5x, recorded in
// .scratch/wisp/probes/orch/2026-10-09-c1c2-resident.raw.md) is
// summon = '' / mute = '' / cancel = 'Esc' / panel = ''. One second after a
// boot that printed hotkeys_live=3, the reload bridge rebound onto exactly
// those four strings and printed live=0: summon, mute and panel died on the
// machine, and the only thing between them and a dead key was the one hop this
// chain skipped - ball.ApplyHotkeyDefaults.
//
// WHAT THE CASES DRIVE. The real assembly shape, not a stand-in:
// startResidentBall(reg, nil, hotCfg, hotReload) over a live ball, with
//   - hotCfg    = the construction half (the naked config view; its defaults
//                 are merged by the consumer, residentBallHotkeyChain258), and
//   - hotReload = hotkeyReloadSource296, the production hot-tier source, hoisted
//                 out of runResident so this file executes the body the
//                 assembly root executes.
// The two are SEPARATE closures on purpose. Ticket 258's positive control passes
// the SAME closure to both halves (startResidentBall(..., src, src)) and plants
// four explicit values, so the bridge could never answer anything but a full
// set - that shape is why this regression shipped green.
//
// TWO FORMS, both from the orchestrator's 2026-10-09 18:2x ruling, because the
// bug has two entrances into the bridge and fixing only the naked mapping would
// leave the second one live:
//   - form 1: the file IS readable and carries three empty slots (the owner's
//     shape, AC#2's 甲 target).
//   - form 2: the file is NOT readable at tick time - the branch answering
//     ball.HotkeyConfig{} - four empty slots straight into the bridge.
//
// Each case asserts two things after ONE bridge hop (the Check() the 1s poll
// performs): (a) the report still holds three live keys - summon, mute, panel -
// and (b) ConfiguredHotkeys() reads back the merged bindings, not empty strings.
// cancel is never counted live at idle (ticket 245: borrowed only during
// Confirming), so live=3 is the healthy reading, not a shortfall.
//
// NOT TOUCHED HERE: ticket 258's assertions (loadHotkey258's "the source may
// answer all-empty" premise at :153 is TRUE - the defaults are not the source's
// business; what was missing is the combination "all-empty source => merged set
// must still be live=3"), internal/config/schema.go's four [hotkey] defaults,
// and anything in internal/ball. The none/off family is Q-82 (owner's list,
// default action: not built) - no normalize for it here.

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// ownerHotkeyBody296 is the owner's [hotkey] section, copied line for line from
// the reading quoted above (three empty slots, cancel = "Esc"). It is a fixture
// inside a temp dir - the real file under %APPDATA% is never opened by these
// cases.
const ownerHotkeyBody296 = "[hotkey]\nsummon = \"\"\nmute = \"\"\ncancel = \"Esc\"\npanel = \"\"\n"

// liveNames296 lists the slots an idle ball is expected to hold, in the
// registration order internal/ball uses, so a failure names what happened to
// each key instead of only a count.
func liveNames296(rep ball.HotkeyReport) []string {
	out := make([]string, 0, len(rep.Bindings()))
	for _, b := range rep.Bindings() {
		if b.Status == ball.HotkeyLive {
			out = append(out, b.Name)
		}
	}
	return out
}

// hotkeyStatuses296 renders the per-slot picture for a failure message
// ("summon=disabled (unset in config)=\"\" mute=... cancel=... panel=...").
func hotkeyStatuses296(rep ball.HotkeyReport) string {
	parts := make([]string, 0, len(rep.Bindings()))
	for _, b := range rep.Bindings() {
		parts = append(parts, b.Name+"="+b.Status.String()+"="+quoteBinding296(b.Binding))
	}
	return strings.Join(parts, " ")
}

// quoteBinding296 shows the slot's raw value, because an empty pair of quotes IS
// the shape under test.
func quoteBinding296(binding string) string {
	return `"` + binding + `"`
}

// missingLive296 names the idle keys among summon/mute/panel that are not live.
func missingLive296(live []string) []string {
	out := make([]string, 0, 3)
	for _, want := range []string{"summon", "mute", "panel"} {
		if !slices.Contains(live, want) {
			out = append(out, want)
		}
	}
	return out
}

// requireThreeLive296 is the pair of assertions both forms owe: the report must
// still hold summon/mute/panel, and the ball must have been TOLD the merged
// bindings. Provenance is the failure message's own subject, so it is said by
// the callers, which know which world they built.
func requireThreeLive296(t *testing.T, rb *residentBall, form string) {
	t.Helper()

	rep := rb.b.HotkeyReport()
	live := liveNames296(rep)
	cfg := rb.b.ConfiguredHotkeys()

	// (b) the merged set, not the empty strings, is what the ball holds. Both
	// forms merge to exactly the compiled defaults (the owner's cancel tag and
	// the unreadable branch's empty cancel both read "Esc"). Said first and
	// non-fatally so one red run reports both halves.
	want := ball.DefaultHotkeys()
	if cfg != want {
		t.Errorf("AC#1 %s RED: ConfiguredHotkeys() = %+v, want the merged set %+v - the bridge rebound onto the raw empty slots", form, cfg, want)
	}

	// (a) three keys still live. The headline number is the one the machine
	// printed when this broke (live=0), so the red line quotes it.
	if missing := missingLive296(live); len(live) != 3 || len(missing) > 0 {
		t.Fatalf("AC#1 %s RED: after the bridge's first Check live=%d %v, want live=3 (summon, mute, panel; missing %v), configured=%+v, per-slot=%s. The reload source handed the bridge un-merged empty slots and the keys were unregistered - ticket 296.",
			form, len(live), live, missing, cfg, hotkeyStatuses296(rep))
	}
}

// Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots is AC#1 form 1: the owner's
// file verbatim - readable, three empty slots, cancel = "Esc" - driven through
// the real bridge over a real ball. Red before ticket 296's fix (live=0 one
// second after a live=3 boot), green after.
func Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(ownerHotkeyBody296), 0o600); err != nil {
		t.Fatalf("write the owner's [hotkey] fixture: %v", err)
	}

	// Premise, said out loud: the raw view really is three empties plus the
	// owner's cancel tag. If a schema change ever fills those slots at parse
	// time, this case must stop measuring - it is not a fallback for the
	// defaults, it is the hole between the raw view and the merged set.
	if raw := loadHotkey258(t, dir); raw != (ball.HotkeyConfig{Cancel: "Esc"}) {
		t.Fatalf("premise moved: the [hotkey] fixture answers %+v, want three empty slots plus cancel = %q", raw, "Esc")
	}

	hotCfg := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }
	hotReload := func() ball.HotkeyConfig { return hotkeyReloadSource296(dir) }

	rb := startResidentBall(observe.NewRegistry(), nil, hotCfg, hotReload)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window, so the rebind has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()
	if !rb.hotkeyBridgeArmed {
		t.Fatalf("the bridge did not arm with a live src: this case would measure nothing")
	}

	// The boot half is the control: residentBallHotkeyChain258 merges, so the
	// three keys ARE live before the bridge speaks. If the host cannot register
	// them at all, neither the bug nor the fix is measurable here.
	boot := rb.b.HotkeyReport()
	if len(liveNames296(boot)) != 3 {
		t.Skipf("SKIP-LOUD: boot could not register the three default combinations on this host (%s), so the live count cannot be read here", hotkeyStatuses296(boot))
	}

	rb.hotkeyBridgeCheck258()
	requireThreeLive296(t, rb, "form-1 (owner's three empty slots)")
}

// Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable is AC#1 form 2, the branch
// the orchestrator's ruling added to 甲's range: the tick-time read of
// config.toml FAILS (no file in the directory the source reads), and that
// branch answers four empty slots - which reach the bridge the same way and kill
// the same three keys. Wrapping only the naked mapping would leave this one
// live.
func Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable(t *testing.T) {
	bootDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(bootDir, configFileName), []byte(ownerHotkeyBody296), 0o600); err != nil {
		t.Fatalf("write the boot config: %v", err)
	}
	// The reload half reads a directory with no config.toml at all - the shape
	// of a file that vanished, was never created, or cannot be opened mid-run.
	reloadDir := t.TempDir()

	// Premise: that directory is the UNREADABLE branch, not a file whose
	// [hotkey] happens to be empty. Two different failures, one same fix.
	if c, _, err := config.LoadFile(filepath.Join(reloadDir, configFileName), nil); err == nil || c != nil {
		t.Fatalf("premise moved: the reload dir was readable (cfg=%v err=%v); form 2 must exercise the unreadable branch", c, err)
	}

	hotCfg := func() ball.HotkeyConfig { return loadHotkey258(t, bootDir) }
	hotReload := func() ball.HotkeyConfig { return hotkeyReloadSource296(reloadDir) }

	rb := startResidentBall(observe.NewRegistry(), nil, hotCfg, hotReload)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: no ball window on this host, so the rebind has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()
	if !rb.hotkeyBridgeArmed {
		t.Fatalf("the bridge did not arm with a live src: this case would measure nothing")
	}
	boot := rb.b.HotkeyReport()
	if len(liveNames296(boot)) != 3 {
		t.Skipf("SKIP-LOUD: boot could not register the three default combinations on this host (%s), so the live count cannot be read here", hotkeyStatuses296(boot))
	}

	rb.hotkeyBridgeCheck258()
	requireThreeLive296(t, rb, "form-2 (config.toml unreadable at tick time)")
}
