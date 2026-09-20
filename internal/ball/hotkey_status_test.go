//go:build windows

package ball

// Deterministic (no-window) tests for the hotkey registration semantics and
// the config-reload bridge (ticket 64 A1/A1b). These run in the DEFAULT suite:
// the status classification and the reloader diff need no desktop, and the
// live layer (hotkey_live_test.go, -tags winlive) proves the same codes come
// back from the real RegisterHotKey.

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

// fakeRegistry is an in-memory RegisterHotKey: it records what was attempted
// and can be told which combinations fail (and with which error).
type fakeRegistry struct {
	attempted []uint32
	unattempt []uint32 // ids whose UnregisterHotKey was called
	fail      map[uint32]error
	held      map[string]bool // "mods|vk" -> already owned
}

func (f *fakeRegistry) register(id uint32, acc Accelerator) error {
	f.attempted = append(f.attempted, id)
	if err, ok := f.fail[id]; ok {
		return err
	}
	key := accelKey(acc)
	if f.held[key] {
		return windows.ERROR_HOTKEY_ALREADY_REGISTERED
	}
	if f.held == nil {
		f.held = map[string]bool{}
	}
	f.held[key] = true
	return nil
}

func (f *fakeRegistry) unregister(id uint32) {
	f.unattempt = append(f.unattempt, id)
}

func accelKey(a Accelerator) string { return fmt.Sprintf("%d|%d", a.Mods, a.VK) }

func fullConfig() HotkeyConfig {
	return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}
}

// TestRegisterAllLiveSet asserts the registration SET itself - the thing no
// test in this repo ever checked (ticket 64 A1: "四项默认热键全注册失败、测试仍报 PASS").
func TestRegisterAllLiveSet(t *testing.T) {
	reg := &fakeRegistry{}
	rep := registerAllWith(0, fullConfig(), reg.register)

	want := []uint32{hkSummon, hkMute, hkCancel, hkPanel}
	if len(reg.attempted) != len(want) {
		t.Fatalf("attempted %v, want all four %v", reg.attempted, want)
	}
	live := rep.Live()
	if len(live) != 4 {
		t.Fatalf("live set %v, want 4 entries", live)
	}
	for _, id := range want {
		if !rep.IsLive(id) {
			t.Errorf("id %d not live; report %v", id, rep.Bindings())
		}
	}
	if !rep.AllLive() {
		t.Error("AllLive false with every binding accepted")
	}
	if p := rep.Problems(); len(p) != 0 {
		t.Errorf("problems on a clean run: %v", p)
	}
	if live[hkSummon].VK != 'Q' || live[hkMute].VK != 'M' {
		t.Errorf("live accelerators wrong: %+v", live)
	}
}

// TestRegisterAllSplitsFailureFamilies is the A1b contract: "occupied by
// somebody else" and "we never tried" must not share a message, and one
// failure must not drop the others.
func TestRegisterAllSplitsFailureFamilies(t *testing.T) {
	reg := &fakeRegistry{fail: map[uint32]error{
		hkMute:  windows.ERROR_HOTKEY_ALREADY_REGISTERED,
		hkPanel: syscall.EINVAL,
	}}
	cfg := HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+Bogus+"}
	rep := registerAllWith(0, cfg, reg.register)

	// summon + cancel survive a neighbour's failure.
	if !rep.IsLive(hkSummon) || !rep.IsLive(hkCancel) {
		t.Fatalf("one failure dropped the others: %+v", rep.Bindings())
	}
	// mute: ATTEMPTED and refused with 1409 => occupied.
	mute, _ := rep.Binding(hkMute)
	if mute.Status != HotkeyTaken || !mute.Status.Attempted() {
		t.Errorf("mute status = %v (attempted=%v), want occupied-and-attempted", mute.Status, mute.Status.Attempted())
	}
	if !errors.Is(mute.Err, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
		t.Errorf("mute err = %v, want ERROR_HOTKEY_ALREADY_REGISTERED", mute.Err)
	}
	// panel: the binding string is junk => NOT attempted at all.
	panel, _ := rep.Binding(hkPanel)
	if panel.Status != HotkeyUnparsable || panel.Status.Attempted() {
		t.Errorf("panel status = %v (attempted=%v), want not-attempted/unparsable", panel.Status, panel.Status.Attempted())
	}
	for _, id := range reg.attempted {
		if id == hkPanel {
			t.Fatal("an unparsable binding reached RegisterHotKey")
		}
	}
	// An empty binding is "disabled", never "failed".
	off := registerAllWith(0, HotkeyConfig{Summon: "Ctrl+Alt+Q"}, (&fakeRegistry{}).register)
	if d, _ := off.Binding(hkMute); d.Status != HotkeyDisabled || d.Status.Attempted() {
		t.Errorf("unset mute = %v, want disabled-and-not-attempted", d.Status)
	}
	if off.IsLive(hkMute) || off.IsLive(hkCancel) || off.IsLive(hkPanel) {
		t.Errorf("unset bindings reported live: %+v", off.Bindings())
	}

	// The user-visible lines must name the family, and must never re-introduce
	// the old guessy "already taken?" wording (that question mark WAS the bug).
	probs := strings.Join(rep.Problems(), "\n")
	if !strings.Contains(probs, "occupied by another program") {
		t.Errorf("occupied case missing from the user-visible problems: %q", probs)
	}
	if !strings.Contains(probs, "not a valid key combination") {
		t.Errorf("never-attempted case missing: %q", probs)
	}
	if strings.Contains(probs, "already taken?") {
		t.Errorf("the indeterminate wording is back: %q", probs)
	}
	if rep.AllLive() {
		t.Error("AllLive true with two failed bindings")
	}
}

// TestUnregisterAllDropsKnownSlots pins why rebind unregisters all four ids:
// a slot we lost track of would otherwise return 1409 on the next rebind and
// be misreported as "occupied by another program".
func TestUnregisterAllDropsKnownSlots(t *testing.T) {
	reg := &fakeRegistry{}
	unregisterAllWith(reg.unregister)
	if len(reg.unattempt) != 4 {
		t.Fatalf("unregistered %v, want the four known slots", reg.unattempt)
	}
	for _, id := range []uint32{hkSummon, hkMute, hkCancel, hkPanel} {
		found := false
		for _, got := range reg.unattempt {
			if got == id {
				found = true
			}
		}
		if !found {
			t.Errorf("slot %d not unregistered", id)
		}
	}
}

// TestHotkeyStatusString checks the log/UI vocabulary is stable and specific.
func TestHotkeyStatusString(t *testing.T) {
	cases := map[HotkeyStatus]string{
		HotkeyLive:       "live",
		HotkeyDisabled:   "disabled (unset in config)",
		HotkeyUnparsable: "not attempted (binding cannot be parsed)",
		HotkeyTaken:      "occupied by another app",
		HotkeyError:      "registration failed",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("status %d = %q, want %q", s, got, want)
		}
	}
}

// recorder is the HotkeyBinder side of the bridge tests.
type recorder struct {
	calls []HotkeyConfig
	rep   HotkeyReport
}

func (r *recorder) RebindHotkeys(cfg HotkeyConfig) HotkeyReport {
	r.calls = append(r.calls, cfg)
	reg := &fakeRegistry{}
	r.rep = registerAllWith(0, cfg, reg.register)
	return r.rep
}

// TestHotkeyReloaderRebindsOnConfigChange is the A1 wiring proof at the
// deterministic layer: a changed [hotkey] section drives exactly one rebind,
// an unchanged one drives none (no Win32 churn per tick).
func TestHotkeyReloaderRebindsOnConfigChange(t *testing.T) {
	b := &recorder{}
	current := fullConfig()
	src := current // the "config file" the source reads
	r := NewHotkeyReloader(b, current, func() HotkeyConfig { return src })
	refreshed := 0
	r.Refresh = func() error { refreshed++; return nil }

	if changed, _ := r.Check(); changed {
		t.Fatal("Check rebound an unchanged config")
	}
	if r.Rebinds() != 0 || len(b.calls) != 0 {
		t.Fatalf("unchanged config triggered %d rebinds", r.Rebinds())
	}
	if refreshed != 1 {
		t.Fatalf("Refresh called %d times, want 1 (the bridge drives the config poll)", refreshed)
	}

	src = HotkeyConfig{Summon: "Ctrl+Alt+R", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}
	changed, rep := r.Check()
	if !changed {
		t.Fatal("a changed summon binding did not rebind")
	}
	if len(b.calls) != 1 || b.calls[0].Summon != "Ctrl+Alt+R" {
		t.Fatalf("binder calls %+v", b.calls)
	}
	if live, _ := rep.Binding(hkSummon); live.Acc.VK != 'R' {
		t.Fatalf("live set still holds the old key: %+v", rep.Bindings())
	}
	if r.Applied().Summon != "Ctrl+Alt+R" {
		t.Errorf("Applied = %v, want the new set", r.Applied())
	}

	// Same content again is still "unchanged".
	r.Check()
	if r.Rebinds() != 1 {
		t.Errorf("the bridge re-registered an unchanged set (rebinds=%d)", r.Rebinds())
	}

	// Back to the old key: the registration set must follow the file again.
	src = current
	if _, rep := r.Check(); rep.IsLive(hkSummon) {
		if got, _ := rep.Binding(hkSummon); got.Acc.VK != 'Q' {
			t.Errorf("reverting to the default summon did not land: %+v", got)
		}
	}
	if r.Rebinds() != 2 {
		t.Errorf("rebinds = %d, want 2", r.Rebinds())
	}
}

// TestHotkeyReloaderSectionsHookIgnoresOtherSections pins the OnReload path:
// only a [hotkey] notification may re-register.
func TestHotkeyReloaderSectionsHookIgnoresOtherSections(t *testing.T) {
	b := &recorder{}
	src := fullConfig()
	src.Panel = "Ctrl+Alt+J"
	r := NewHotkeyReloader(b, fullConfig(), func() HotkeyConfig { return src })
	r.Sections([]string{"voice", "session"})
	if len(b.calls) != 0 {
		t.Fatalf("unrelated sections triggered a rebind: %+v", b.calls)
	}
	r.OnReload()([]string{"hotkey"})
	if len(b.calls) != 1 || b.calls[0].Panel != "Ctrl+Alt+J" {
		t.Fatalf("the hotkey notification did not rebind: %+v", b.calls)
	}
}

// TestApplyHotkeyDefaults keeps a config file that never mentions a key from
// silently disabling it (the schema has no `default:` tags for summon/mute/panel).
func TestApplyHotkeyDefaults(t *testing.T) {
	got := ApplyHotkeyDefaults(HotkeyConfig{Cancel: "Ctrl+."})
	want := DefaultHotkeys()
	want.Cancel = "Ctrl+."
	if got != want {
		t.Errorf("ApplyHotkeyDefaults(empty-ish) = %+v, want %+v", got, want)
	}
	if AltSummonSpace != "Ctrl+Alt+Space" {
		t.Errorf("the documented alternative moved: %q", AltSummonSpace)
	}
	if _, err := ParseAccelerator(AltSummonSpace); err != nil {
		t.Errorf("the documented alternative %q must parse: %v", AltSummonSpace, err)
	}
}
