//go:build windows

package main

// Ticket 258 AC#1/AC#2/AC#3, default tier: the [hotkey] chain at the three
// seams this host actually has.
//
// WHAT THIS FILE MEASURES. Three rulers, one per AC, all read off production
// shapes rather than re-typed copies:
//
//  1. THE CONSTRUCTION CHAIN (AC#1). residentBallHotkeyChain258 is the function
//     resident_windows.go feeds and resident_ball_windows.go consumes. The case
//     drives it the way the assembly root does - a config.LoadFile read off a
//     real config.toml, including the MISSING-file shape - and asserts the four
//     bindings and the provenance word. AC#1's sentence ("退回 DefaultHotkeys 并
//     说得出这句话") is checked as a sentence the process prints: the winlive
//     family reads it from the shipped binary's stdout (see
//     resident_hotkey_live_258_windows_test.go); here the word itself is pinned
//     as data, so a rename of the provenance goes red instead of drifting.
//     Since ticket 258-r2 the word follows the SET, not the file's mere
//     readability: a file whose [hotkey] section is absent merges to exactly
//     ball.DefaultHotkeys() (only the schema's cancel tag arrives) and must
//     print defaults. Test258SectionMissingTierWordIsDefaults pins that with
//     no ball window; Test258ConstructionChainMissingFileFallsBackAndSaysIt
//     keeps the other two fallback shapes honest.
//
//  2. THE BRIDGE (AC#2's in-process half). The hot tier's promise is "no
//     restart": the case installs the ticket-64 bridge over a REAL ball
//     (startResidentBall with a live src closure, exactly the production
//     wiring), edits the config file on disk inside the same process, drives
//     Check() the way the 1s poll would, and requires the NEW combination to be
//     what Win32 holds. The mutation control for this case lives in
//     Test258BridgeMutationNoSrcKeepsOldBinding: with the rebind hop removed,
//     the same edit must leave the OLD binding live and the case red.
//
//  3. THE HONEST FAILURE SENTENCE (AC#3). A squatted combination (registered on
//     the ball's own window under a spare id, the same in-process squat
//     internal/ball/hotkey_live_test.go uses) must surface as HotkeyTaken with
//     the NEW value named, in Problems(), in the verdict's live count, and in
//     the per-slot summary - never as a silent keep-the-old-value success.
//
// WHAT THIS FILE DELIBERATELY DOES NOT ASSERT: the `wisp run` half of the
// ticket's "两条入口" wording is an empty set (258-a1 census Q1c: run.go's
// production code references no ball), so there is no run-leg ball to bind
// anything; the resident leg is the only subject. And the verdict-string pins of
// ticket 228 (resident_ball_228_windows_test.go:45,
// resident_ball_live_228_windows_test.go:144/160/163) are NOT touched here:
// their literals still hold - the verdict gained a provenance word, it did not
// lose "the floating ball window is up in this process" or change the
// "hotkeys live %d/4" shape.

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"golang.org/x/sys/windows"
)

// spareHKID258 is the id the squat registers under. Any id that is not one of
// the four real ones works; 9 keeps clear of hkSummon..hkPanel (1..4).
const spareHKID258 = 9

// hkSummon258 / hkPanel258 are the wParam ids internal/ball registers the four
// slots under (internal/ball/hotkey_windows.go's hkSummon=1, hkPanel=4).
const (
	hkSummon258 uint32 = 1
	hkPanel258  uint32 = 4
)

// writeConfig258 writes a config.toml body and returns the path. The canonical
// branch uses SaveFile - what the product itself writes (firstrun.go) - so the
// fixture is the production shape and not a hand-rolled TOML dialect.
func writeConfig258(t *testing.T, dir string, hotkeyBody string) string {
	t.Helper()
	path := filepath.Join(dir, configFileName)
	if hotkeyBody == "" {
		if err := config.SaveFile(path, config.NewDefaults()); err != nil {
			t.Fatalf("SaveFile canonical config: %v", err)
		}
		return path
	}
	body := "[hotkey]\n" + hotkeyBody + "\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write [hotkey]-only config: %v", err)
	}
	return path
}

// loadHotkey258 maps one LoadFile result into ball.HotkeyConfig the way both
// production closures (resident_windows.go) do.
func loadHotkey258(t *testing.T, dir string) ball.HotkeyConfig {
	t.Helper()
	c, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
	if err != nil {
		return ball.HotkeyConfig{}
	}
	return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
}

// Test258ConstructionChainTakesConfigValues is AC#1's positive control at the
// chain seam: a config.toml whose [hotkey] carries non-default values produces
// exactly those bindings, provenance "config".
func Test258ConstructionChainTakesConfigValues(t *testing.T) {
	dir := t.TempDir()
	writeConfig258(t, dir, `summon = "Ctrl+Alt+Z"
mute = "Ctrl+Alt+X"
cancel = "Esc"
panel = "Ctrl+Alt+V"`)
	if got := loadHotkey258(t, dir); got.Summon != "Ctrl+Alt+Z" {
		t.Fatalf("fixture load answered %+v, want the planted summon", got)
	}
	hotCfg := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }

	got, provenance := residentBallHotkeyChain258(hotCfg)
	if provenance != hotkeyProvenanceConfig {
		t.Errorf("provenance = %q, want %q: a readable config with [hotkey] values is config-sourced", provenance, hotkeyProvenanceConfig)
	}
	want := ball.HotkeyConfig{Summon: "Ctrl+Alt+Z", Mute: "Ctrl+Alt+X", Cancel: "Esc", Panel: "Ctrl+Alt+V"}
	if got != want {
		t.Errorf("chain = %+v, want the planted values %+v: the four slots must follow config.toml (AC#1)", got, want)
	}
	if got == ball.DefaultHotkeys() {
		t.Errorf("chain answered the compiled defaults for a config that named other values: the [hotkey] section was ignored")
	}
}

// Test258ConstructionChainMissingFileFallsBackAndSaysIt is AC#1's fallback
// half: a MISSING config file (and a host view that answers empty) yields the
// compiled defaults, and the provenance word is the "说得出这句话" half.
func Test258ConstructionChainMissingFileFallsBackAndSaysIt(t *testing.T) {
	// The nil-view shape: no host config at all (the pre-258 shape, still legal).
	got, provenance := residentBallHotkeyChain258(nil)
	if provenance != hotkeyProvenanceNone {
		t.Errorf("nil view provenance = %q, want %q", provenance, hotkeyProvenanceNone)
	}
	if got != ball.DefaultHotkeys() {
		t.Errorf("nil view chain = %+v, want DefaultHotkeys()", got)
	}

	// The missing-file shape, driven exactly the way resident_windows.go drives
	// it: a LoadFile over a directory with no config.toml answers an error, the
	// closure returns the empty view, and the chain names the fallback.
	dir := t.TempDir()
	hotCfg := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }
	if c := loadHotkey258(t, dir); c != (ball.HotkeyConfig{}) {
		t.Fatalf("premise moved: an empty directory produced the non-empty view %+v", c)
	}
	got, provenance = residentBallHotkeyChain258(hotCfg)
	if provenance != hotkeyProvenanceDefaults {
		t.Errorf("missing-file provenance = %q, want %q: AC#1 requires the fallback to be named, not silent", provenance, hotkeyProvenanceDefaults)
	}
	if got != ball.DefaultHotkeys() {
		t.Errorf("missing-file chain = %+v, want DefaultHotkeys()", got)
	}

	// The section-missing shape: a valid file with no [hotkey] table at all.
	// The schema only defaults cancel, so the view arrives with three empty
	// slots plus that tag; ApplyHotkeyDefaults merges exactly the compiled
	// defaults out of it, and since ticket 258-r2 the tier word is defaults -
	// 258-v1 §3's wording ruling: a file that contributed no binding may not
	// be named as the source of the set. The VALUE is unchanged by that fix
	// (AC#1's 不许静默换成另一套 is about the set, and no shape here returns a
	// fourth one).
	writeConfig258(t, dir, "")
	got, provenance = residentBallHotkeyChain258(hotCfg)
	if provenance != hotkeyProvenanceDefaults {
		t.Errorf("section-missing provenance = %q, want %q: the [hotkey] table is absent, the set is the compiled defaults", provenance, hotkeyProvenanceDefaults)
	}
	if got != ball.DefaultHotkeys() {
		t.Errorf("canonical-file chain = %+v, want the defaults the empty slots fall back to", got)
	}
}

// Test258SectionMissingTierWordIsDefaults is AC#1 condition ① (ticket 258-r2):
// the honest word for a file whose [hotkey] section is missing. It needs no
// ball window - the chain is a pure function over the view the assembly root
// reads - so it runs in CI every day, and it has teeth against exactly the
// rot it replaced: restore the pre-r2 branch (any non-empty view answers
// provenance "config") and the first clause goes red, because that view
// carries the schema's cancel = "Esc" tag and is not all-empty.
func Test258SectionMissingTierWordIsDefaults(t *testing.T) {
	dir := t.TempDir()
	// A file a hand editor would actually leave: version line only, no
	// [hotkey] table anywhere in the bytes.
	body := "schema_version = " + strconv.Itoa(config.SchemaVersionCurrent) + "\n"
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("write section-missing file: %v", err)
	}
	hotCfg := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }

	// Premise, said out loud: LoadFile succeeds and the view is NOT all-empty
	// (the schema tags cancel), which is why the pre-r2 code fell through to
	// the config branch here.
	c := loadHotkey258(t, dir)
	if c == (ball.HotkeyConfig{}) {
		t.Fatalf("premise moved: a file without [hotkey] answered the all-empty view; this case pins the NON-empty schema-tag shape")
	}
	if c.Cancel != "Esc" || c.Summon != "" || c.Mute != "" || c.Panel != "" {
		t.Fatalf("premise moved: the section-missing view is %+v, want only the schema's cancel = \"Esc\"", c)
	}

	got, provenance := residentBallHotkeyChain258(hotCfg)
	if provenance != hotkeyProvenanceDefaults {
		t.Errorf("AC#1 condition-1 RED: section-missing provenance = %q, want %q (the merged set IS ball.DefaultHotkeys(); \"config\" names a file that contributed no binding)", provenance, hotkeyProvenanceDefaults)
	}
	if got != ball.DefaultHotkeys() {
		t.Errorf("the fallback swapped the set instead of only the word: chain = %+v, want DefaultHotkeys()", got)
	}

	// The word still tells the shapes apart: the same directory, one section
	// that moves a single slot, answers "config". Without this control the
	// "defaults" side could be bought by a judge that always says defaults.
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body+"[hotkey]\nsummon = \"Ctrl+Alt+Z\"\n"), 0o600); err != nil {
		t.Fatalf("plant one-slot section: %v", err)
	}
	got, provenance = residentBallHotkeyChain258(hotCfg)
	if provenance != hotkeyProvenanceConfig {
		t.Errorf("positive control RED: a file naming a non-default summon answered %q, want %q", provenance, hotkeyProvenanceConfig)
	}
	if got.Summon != "Ctrl+Alt+Z" || got != ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: "Ctrl+Alt+Z"}) {
		t.Errorf("positive control set RED: chain = %+v, want the planted summon over the defaulted rest", got)
	}

	// The rule's boundary, said in a pin instead of a comment: a file that
	// DOES carry [hotkey] but spells out the compiled defaults verbatim also
	// prints defaults - the live keys ARE the defaults there, and that is the
	// honest reading the word owes (258-v1 §3: the set, not the file's mood).
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body+"[hotkey]\nsummon = \"Ctrl+Alt+Q\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"Ctrl+Alt+P\"\n"), 0o600); err != nil {
		t.Fatalf("plant all-defaults section: %v", err)
	}
	if _, provenance := residentBallHotkeyChain258(hotCfg); provenance != hotkeyProvenanceDefaults {
		t.Errorf("a verbatim-defaults [hotkey] section answered %q: the set equals ball.DefaultHotkeys(), so only defaults is true of it", provenance)
	}
}

// Test258ProvenanceWordsAreThePrintedOnes pins the three provenance words as
// data, because the verdict embeds them ("hotkeys from %s") and the winlive
// family greps them off the shipped binary's stdout. Renaming one of these
// constants is a verdict-sentence change and has to be said out loud.
func Test258ProvenanceWordsAreThePrintedOnes(t *testing.T) {
	for _, w := range []string{hotkeyProvenanceConfig, hotkeyProvenanceDefaults, hotkeyProvenanceNone} {
		if w == "" {
			t.Fatal("a provenance word is empty: the verdict would print 'hotkeys from '")
		}
		if strings.ContainsAny(w, "%\n") {
			t.Errorf("provenance word %q carries format/newline bytes", w)
		}
	}
	if hotkeyProvenanceDefaults != "defaults" {
		t.Errorf("the AC#1 fallback word is %q, want %q (the live case greps 'hotkeys from defaults')", hotkeyProvenanceDefaults, "defaults")
	}
}

// Test258BridgeRebindsLiveKeysFromConfigEdit is AC#2's in-process positive
// control: the full production wiring (startResidentBall over a live ball with
// a src closure), an on-disk [hotkey] edit, and the poll hop - no restart, new
// combination live.
func Test258BridgeRebindsLiveKeysFromConfigEdit(t *testing.T) {
	dir := t.TempDir()
	// The boot plant is a NON-default summon since 258-r2: with Ctrl+Alt+Q the
	// boot VK could not tell "config carried the value" from "the fallback
	// printed the same word", so this line only became a reading once it
	// differed from ball.DefaultHotkeys().
	writeConfig258(t, dir, `summon = "Ctrl+Alt+Z"
mute = "Ctrl+Alt+M"
cancel = "Esc"
panel = "Ctrl+Alt+P"`)
	src := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }

	rb := startResidentBall(observe.NewRegistry(), nil, src, src)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the live rebind has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()
	if !rb.hotkeyBridgeArmed {
		t.Fatalf("the bridge did not arm with a live src: the hot tier has no rebind path in this process")
	}
	if rb.hotkeySourceName() != hotkeyProvenanceConfig {
		t.Errorf("provenance = %q, want %q", rb.hotkeySourceName(), hotkeyProvenanceConfig)
	}
	boot := rb.b.HotkeyReport()
	if !boot.IsLive(hkSummon258) {
		t.Fatalf("boot registration incomplete: %+v", boot.Bindings())
	}
	if got := boot.Live()[hkSummon258].VK; got != 'Z' {
		t.Fatalf("boot summon VK = 0x%X, want 'Z' (the planted non-default; a Q here means the chain answered defaults)", got)
	}

	// THE EDIT: summon -> Ctrl+Alt+R, the same plant internal/ball's live suite
	// uses, written to the file the src closure reads.
	writeConfig258(t, dir, `summon = "Ctrl+Alt+R"
mute = "Ctrl+Alt+M"
cancel = "Esc"
panel = "Ctrl+Alt+P"`)

	// Drive Check() the way the bridge's own 1s tick does (r.Run's body is
	// checkSafe -> Check; calling Check here is the same hop without the extra
	// goroutine). Bounded by attempts, never by wall-clock arithmetic (D42#9).
	seen := false
	for i := 0; i < 50 && !seen; i++ {
		rb.hotkeyBridgeCheck258()
		seen = rb.b.RegisteredHotkeys()[hkSummon258].VK == 'R'
		if !seen {
			time.Sleep(20 * time.Millisecond)
		}
	}
	if !seen {
		live := rb.b.RegisteredHotkeys()
		t.Fatalf("the [hotkey] edit never reached Win32: summon VK = 0x%X (%+v), want 'R'. The hot tier's promise (no-restart rebind) did not hold.", live[hkSummon258].VK, live)
	}
	rep := rb.b.HotkeyReport()
	if got := rep.Live()[hkSummon258].VK; got != 'R' {
		t.Errorf("post-rebind report summon VK = 0x%X, want 'R'", got)
	}
	if rb.b.ConfiguredHotkeys().Summon != "Ctrl+Alt+R" {
		t.Errorf("ConfiguredHotkeys().Summon = %q, want the planted Ctrl+Alt+R", rb.b.ConfiguredHotkeys().Summon)
	}
}

// Test258BridgeMutationNoSrcKeepsOldBinding is the mutation control: the same
// edit, but the bridge hop removed (hotReload nil - the pre-258 wiring). The
// OLD binding must stay live and the NEW one must never arrive; if this case
// ever goes green with the new value, the case above has no teeth.
func Test258BridgeMutationNoSrcKeepsOldBinding(t *testing.T) {
	dir := t.TempDir()
	writeConfig258(t, dir, `summon = "Ctrl+Alt+Q"
mute = "Ctrl+Alt+M"
cancel = "Esc"
panel = "Ctrl+Alt+P"`)
	src := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }
	// hotCfg = src (the construction still reads config), hotReload = nil: this
	// is exactly the wiring the mutation removes.
	rb := startResidentBall(observe.NewRegistry(), nil, src, nil)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: no ball window on this host (%q)", rb.verdictStatusLineFor258())
	}
	defer rb.stop()
	if rb.hotkeyBridgeArmed {
		t.Fatalf("the bridge armed with a nil hotReload: the mutation did not remove the hop, so this control proves nothing")
	}
	writeConfig258(t, dir, `summon = "Ctrl+Alt+R"
mute = "Ctrl+Alt+M"
cancel = "Esc"
panel = "Ctrl+Alt+P"`)
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
	}
	live := rb.b.RegisteredHotkeys()
	if live[hkSummon258].VK == 'R' {
		t.Fatalf("MUTATION CONTROL RED: with the bridge hop removed the NEW summon still went live (VK=0x%X), so Test258BridgeRebindsLiveKeysFromConfigEdit is measuring nothing", live[hkSummon258].VK)
	}
	if live[hkSummon258].VK != 'Q' {
		t.Fatalf("the old binding did not survive a no-rebind host (VK=0x%X): the premise moved", live[hkSummon258].VK)
	}
	// The construction-time view also still answers the boot-time value for the
	// ball (the ball was TOLD Ctrl+Alt+Q at boot; nothing re-told it).
	if got := rb.b.ConfiguredHotkeys().Summon; got != "Ctrl+Alt+Q" {
		t.Errorf("ConfiguredHotkeys().Summon = %q after the edit with no bridge, want the boot-time Ctrl+Alt+Q", got)
	}
}

// Test258OccupiedCombinationNamesTheNewValue is AC#3: a combination another
// owner holds must surface as taken, with the NEW binding named, in Problems()
// and in the summary - the honest sentence, not a silent keep-the-old-value.
func Test258OccupiedCombinationNamesTheNewValue(t *testing.T) {
	dir := t.TempDir()
	// Boot summon is a NON-default Z since 258-r2: the verdict this case reads
	// must legitimately say "hotkeys from config", which under the fixed tier
	// rule only holds when the file moved at least one slot off the compiled
	// defaults.
	writeConfig258(t, dir, `summon = "Ctrl+Alt+Z"
mute = "Ctrl+Alt+M"
cancel = "Esc"
panel = "Ctrl+Alt+P"`)
	src := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }
	rb := startResidentBall(observe.NewRegistry(), nil, src, src)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: no ball window on this host (%q)", rb.verdictStatusLineFor258())
	}
	defer rb.stop()

	// Squat Ctrl+Alt+V on the ball's own window under a spare id: RegisterHotKey
	// is desktop-wide and id-blind, so from panel's point of view "another
	// program" owns it. The squat lands on the OWNING thread (uiRun inside the
	// seam), like internal/ball's own helper does.
	acc, err := ball.ParseAccelerator("Ctrl+Alt+V")
	if err != nil {
		t.Fatalf("parse Ctrl+Alt+V: %v", err)
	}
	if sqErr := rb.b.DebugRegisterHotkeySquat(spareHKID258, acc); sqErr != nil {
		if errors.Is(sqErr, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
			t.Skipf("SKIP-LOUD: Ctrl+Alt+V is already owned by another program on this machine, so the in-process squat cannot be built and the real occupied path was not exercised (%v)", sqErr)
		}
		t.Fatalf("squat RegisterHotKey failed: %v", sqErr)
	}
	defer rb.b.DebugUnregisterHotkeySquat(spareHKID258)

	// THE EDIT: panel moves onto the squatted combination, bridge hops it in.
	// One slot moves - the boot Z stays put so the diff is exactly the panel.
	writeConfig258(t, dir, `summon = "Ctrl+Alt+Z"
mute = "Ctrl+Alt+M"
cancel = "Esc"
panel = "Ctrl+Alt+V"`)
	arrived := false
	for i := 0; i < 50 && !arrived; i++ {
		rb.hotkeyBridgeCheck258()
		arrived = rb.b.ConfiguredHotkeys().Panel == "Ctrl+Alt+V"
		if !arrived {
			time.Sleep(20 * time.Millisecond)
		}
	}
	if !arrived {
		t.Fatalf("the panel edit never reached the bridge (ConfiguredHotkeys=%+v)", rb.b.ConfiguredHotkeys())
	}
	rep := rb.b.HotkeyReport()
	pn, ok := rep.Binding(hkPanel258)
	if !ok {
		t.Fatalf("no panel line in %+v", rep.Bindings())
	}
	if pn.Status != ball.HotkeyTaken {
		t.Fatalf("panel status = %q, want occupied (the squat holds Ctrl+Alt+V)", pn.Status)
	}
	if pn.Binding != "Ctrl+Alt+V" {
		t.Errorf("the taken line names binding %q, want the NEW value Ctrl+Alt+V (not allowed: silently keeping the old value and calling it success)", pn.Binding)
	}
	probs := strings.Join(rep.Problems(), "\n")
	if !strings.Contains(probs, "occupied by another program") || !strings.Contains(probs, "Ctrl+Alt+V") {
		t.Fatalf("Problems() does not name the occupied NEW combination:\n%s", probs)
	}
	// The verdict is the BOOT sentence: rb.verdict is written once by
	// startResidentBall and is NOT rewritten by the rebind (the rebind's honest
	// readings are the reloader's own slog lines plus the report). Assert the
	// boot shape held (provenance + live count) and that the summary, read NOW
	// through hotkeySummary, answers with the NEW binding.
	verdict := rb.verdict
	if !strings.Contains(verdict, "hotkeys from config") {
		t.Errorf("the boot verdict does not name the config provenance: %q", verdict)
	}
	if !strings.Contains(verdict, "hotkeys live 3/4") {
		t.Errorf("the boot verdict lost the live-count shape: %q", verdict)
	}
	// The rebind's outcome is said by the bridge's own record, not by rewriting
	// the boot verdict; Problems() and the summary are the per-slot truth.
	if liveCount := len(rep.Live()); liveCount != 2 {
		t.Errorf("live slots after the rebind = %d, want 2 (summon+mute; cancel is standby, panel is taken)", liveCount)
	}
	if strings.Contains(hotkeySummary(rb.b), "panel=Ctrl+Alt+P") {
		t.Errorf("the summary still names the OLD panel binding: the per-slot line is lying about what Win32 holds: %q", hotkeySummary(rb.b))
	}
	if !strings.Contains(hotkeySummary(rb.b), "panel=Ctrl+Alt+V") {
		t.Errorf("the summary does not name the NEW panel binding: %q", hotkeySummary(rb.b))
	}
}
