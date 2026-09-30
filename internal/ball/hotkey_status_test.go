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
//
// attemptedAcc (ticket 245) is the parallel accelerator list: "which slot did
// it try" alone cannot answer "did anything bind the bare Esc", which is the
// question this ticket is about.
type fakeRegistry struct {
	attempted    []uint32
	attemptedAcc []Accelerator // parallel to attempted
	unattempt    []uint32      // ids whose UnregisterHotKey was called
	fail         map[uint32]error
	held         map[string]bool // "mods|vk" -> already owned
}

func (f *fakeRegistry) register(id uint32, acc Accelerator) error {
	f.attempted = append(f.attempted, id)
	f.attemptedAcc = append(f.attemptedAcc, acc)
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

// attemptsOf counts RegisterHotKey calls made for one id (the return leg of the
// borrow must be exactly the borrow's one attempt - a release that re-binds
// would show up here as a second one).
func (f *fakeRegistry) attemptsOf(id uint32) int {
	n := 0
	for _, got := range f.attempted {
		if got == id {
			n++
		}
	}
	return n
}

// unregistersOf counts UnregisterHotKey calls made for one id.
func (f *fakeRegistry) unregistersOf(id uint32) int {
	n := 0
	for _, got := range f.unattempt {
		if got == id {
			n++
		}
	}
	return n
}

// bindsEsc reports whether any attempt registered the bare VK_ESCAPE with no
// modifier but MOD_NOREPEAT - the spelling RegisterHotKey takes desktop-wide.
// ParseAccelerator("Esc") produces exactly this accelerator, so the helper
// catches both shapes of the defect: the borrow, and an idle pass that binds a
// configured bare Esc.
func (f *fakeRegistry) bindsEsc() bool {
	for _, acc := range f.attemptedAcc {
		if acc.VK == vkEscape && acc.Mods == modNoRepeat {
			return true
		}
	}
	return false
}

func accelKey(a Accelerator) string { return fmt.Sprintf("%d|%d", a.Mods, a.VK) }

func fullConfig() HotkeyConfig {
	return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}
}

// TestRegisterAllLiveSet asserts the registration SET itself - the thing no
// test in this repo ever checked (ticket 64 A1: "四项默认热键全注册失败、测试仍报 PASS").
//
// Ticket 245 rewrote the expected set from four to three, and that edit is a
// TIGHTENING in two directions at once, so both halves are asserted here:
//   - an idle ball may hold EXACTLY summon / mute / panel - one more, or one
//     less, is red (the old "exactly four" would have called the pre-245 bug
//     green, which is why it is not enough any more);
//   - and the thing it may never hold while idle is the cancel slot, which the
//     production default spells as a bare Esc. That is asserted as an
//     accelerator, not as a slot name: bindsEsc() must be false on the idle
//     pass, and true the moment the borrow runs (the positive control below,
//     in the same fake registry - a ruler that only ever sees "no Esc" because
//     nothing registers ever would be a blind ruler).
func TestRegisterAllLiveSet(t *testing.T) {
	reg := &fakeRegistry{}
	rep := registerAllWith(0, fullConfig(), reg.register)

	want := []uint32{hkSummon, hkMute, hkPanel}
	if len(reg.attempted) != len(want) {
		t.Fatalf("attempted %v, want exactly the three idle slots %v (cancel is standby, ticket 245)",
			reg.attempted, want)
	}
	for _, id := range reg.attempted {
		if id == hkCancel {
			t.Fatalf("the idle registration pass bound the cancel slot: %v", reg.attempted)
		}
	}
	live := rep.Live()
	if len(live) != 3 {
		t.Fatalf("live set = %d entries %v, want exactly 3 (summon/mute/panel)", len(live), live)
	}
	for _, id := range want {
		if !rep.IsLive(id) {
			t.Errorf("id %d not live; report %v", id, rep.Bindings())
		}
	}
	if _, held := live[hkCancel]; held {
		t.Errorf("cancel is in the idle live set: %+v", rep.Bindings())
	}
	if reg.bindsEsc() {
		t.Errorf("the idle pass registered a bare Esc (mods|vk seen: %v); every other app on this desktop "+
			"just lost its Esc key (ticket 245)", reg.attemptedAcc)
	}
	cancel, ok := rep.Binding(hkCancel)
	if !ok {
		t.Fatalf("no cancel line in %+v", rep.Bindings())
	}
	if cancel.Status != HotkeyStandby {
		t.Errorf("cancel status = %q, want standby (configured, parsable, not bound while idle)", cancel.Status)
	}
	if cancel.Binding != "Esc" {
		t.Errorf("cancel binding string = %q, want the configured Esc (the report keeps the user's spelling)", cancel.Binding)
	}
	if cancel.Status.Attempted() {
		t.Error("standby must be a never-attempted status")
	}
	if cancel.Acc != (Accelerator{}) {
		t.Errorf("standby carries accelerator %+v, want the zero value (Acc is zero when not attempted)", cancel.Acc)
	}
	if !rep.AllLive() {
		t.Error("AllLive false with every idle-allowed binding accepted")
	}
	if p := rep.Problems(); len(p) != 0 {
		t.Errorf("problems on a clean run: %v", p)
	}
	if live[hkSummon].VK != 'Q' || live[hkMute].VK != 'M' {
		t.Errorf("live accelerators wrong: %+v", live)
	}
}

// TestDefaultHotkeysIdlePassHoldsNoEsc is AC#1 over the PRODUCT DEFAULTS (not
// a test fixture): DefaultHotkeys() is what cmd/wisp's resident leg hands
// ball.New, so this is the deterministic half of "the shipped process does not
// take Esc from the desktop".
func TestDefaultHotkeysIdlePassHoldsNoEsc(t *testing.T) {
	reg := &fakeRegistry{}
	rep := registerAllWith(0, DefaultHotkeys(), reg.register)
	if reg.bindsEsc() {
		t.Fatalf("the production defaults bound a bare Esc globally: %v", reg.attemptedAcc)
	}
	if len(rep.Live()) != 3 {
		t.Fatalf("idle live set with the production defaults = %d, want 3: %+v", len(rep.Live()), rep.Bindings())
	}
	if rep.IsLive(hkCancel) {
		t.Fatalf("cancel live with the production defaults: %+v", rep.Bindings())
	}
	// Positive control for the same fixture one line further: DefaultHotkeys
	// really does say Esc, so the assertion above is not passing because the
	// default changed (it must not change - PLAN.md D43 row 22 is frozen).
	if d := DefaultHotkeys(); d.Cancel != "Esc" {
		t.Fatalf("the cancel default moved to %q; ticket 245 form B keeps the key name", d.Cancel)
	}
	if c, _ := rep.Binding(hkCancel); c.Binding != "Esc" || c.Status != HotkeyStandby {
		t.Fatalf("cancel line = %+v, want the Esc default left standby", c)
	}
}

// TestCancelBorrowRoundTrip is AC#2's two ends at the deterministic layer: the
// borrow binds the bare Esc and the report follows it; the return DROPS it and
// the report follows back - and the return must not re-bind anything, because
// re-binding the configured default is the defect this ticket is about.
func TestCancelBorrowRoundTrip(t *testing.T) {
	reg := &fakeRegistry{}
	rep := registerAllWith(0, fullConfig(), reg.register)
	if reg.attemptsOf(hkCancel) != 0 {
		t.Fatalf("idle pass attempted cancel: %v", reg.attempted)
	}

	// --- borrow (what entering Confirming does)
	if err := takeEscWith(reg.unregister, reg.register); err != nil {
		t.Fatalf("borrow refused by the fake registry: %v", err)
	}
	if reg.attemptsOf(hkCancel) != 1 {
		t.Fatalf("cancel attempts after the borrow = %d, want exactly 1: %v", reg.attemptsOf(hkCancel), reg.attempted)
	}
	if !reg.bindsEsc() {
		t.Fatalf("the borrow did not bind the bare Esc: %v", reg.attemptedAcc)
	}
	borrowed := rep.withCancel(cancelBorrowedLine())
	if len(borrowed.Live()) != 4 {
		t.Fatalf("live set while borrowed = %d, want 4: %+v", len(borrowed.Live()), borrowed.Bindings())
	}
	if got := borrowed.Live()[hkCancel]; got.VK != vkEscape || got.Mods != modNoRepeat {
		t.Errorf("cancel registered as %+v, want the bare VK_ESCAPE with MOD_NOREPEAT", got)
	}
	if !borrowed.AllLive() {
		t.Error("AllLive false while the borrow is live")
	}
	if p := borrowed.Problems(); len(p) != 0 {
		t.Errorf("a successful borrow produced problems: %v", p)
	}

	// --- return (what leaving Confirming does). THIS is the "还" reading:
	// unregister, and nothing re-registered afterwards.
	releaseEscWith(reg.unregister)
	if reg.unregistersOf(hkCancel) < 1 {
		t.Fatal("the return never called UnregisterHotKey for the cancel slot")
	}
	if reg.attemptsOf(hkCancel) != 1 {
		t.Fatalf("cancel registration attempts = %d after the return, want still 1 - the return must drop "+
			"the slot, never re-bind the configured binding (that is ticket 245's defect coming back)",
			reg.attemptsOf(hkCancel))
	}
	returned := borrowed.withCancel(cancelIdleLine(fullConfig().Cancel))
	if len(returned.Live()) != 3 {
		t.Fatalf("live set after the return = %d, want 3: %+v", len(returned.Live()), returned.Bindings())
	}
	if c, _ := returned.Binding(hkCancel); c.Status != HotkeyStandby || c.Binding != "Esc" {
		t.Errorf("cancel after the return = %+v, want standby with the configured spelling kept", c)
	}
	if !returned.AllLive() {
		t.Error("AllLive false after a clean return")
	}

	// Borrow twice without returning is the shape of a re-sync while a card is
	// still up: the second attempt must be refused by the registry that already
	// holds it, and the caller must be told (TakeEscForCancel is idempotent on
	// top of that, but the primitive itself does not lie).
	if err := takeEscWith(reg.unregister, reg.register); err == nil {
		t.Error("double borrow succeeded without a return: the fake registry is not modelling 1409")
	}
}

// TestCancelBorrowFailureIsAProblemLine pins the one failure path the ticket
// rules out silently: if Win32 refuses the borrow (another app owns bare Esc),
// Confirming has no cancel key and the report must SAY so, not leave the slot
// reading as if everything were fine.
func TestCancelBorrowFailureIsAProblemLine(t *testing.T) {
	reg := &fakeRegistry{fail: map[uint32]error{hkCancel: syscall.EINVAL}}
	if err := takeEscWith(reg.unregister, reg.register); err == nil {
		t.Fatal("a failing registrar reported a successful borrow")
	}
	rep := registerAllWith(0, fullConfig(), (&fakeRegistry{}).register).
		withCancel(cancelFailedLine(syscall.EINVAL))
	if rep.AllLive() {
		t.Error("AllLive true while the borrow failed")
	}
	line := strings.Join(rep.Problems(), "\n")
	if !strings.Contains(line, "cancel") || !strings.Contains(line, "was not registered") {
		t.Errorf("the failed borrow produced no user-visible line: %q", line)
	}
	if c, _ := rep.Binding(hkCancel); c.Status.Attempted() != true {
		t.Error("a refused borrow must be an attempted status")
	}
}

// TestStandbyIsNotDisabledAndNotAProblem keeps the three never-attempted
// statuses apart: the user turning a key off, a binding that is junk, and the
// cancel slot being idle-unbound by design must not collapse into one line -
// they mean three different things to the person reading the log.
func TestStandbyIsNotDisabledAndNotAProblem(t *testing.T) {
	rep := registerAllWith(0, fullConfig(), (&fakeRegistry{}).register)
	cancel, _ := rep.Binding(hkCancel)
	mute, _ := rep.Binding(hkMute)
	if cancel.Status == mute.Status {
		t.Fatalf("cancel standby and mute-disabled share a status (%v)", cancel.Status)
	}
	if strings.Contains(cancel.Status.String(), "unset in config") {
		t.Errorf("the standby line claims the user unset the key: %q", cancel.Status.String())
	}
	// An empty cancel binding is still the user's choice, not standby.
	off := registerAllWith(0, HotkeyConfig{Summon: "Ctrl+Alt+Q"}, (&fakeRegistry{}).register)
	if c, _ := off.Binding(hkCancel); c.Status != HotkeyDisabled {
		t.Errorf("unset cancel = %q, want disabled (the user asked for no key)", c.Status)
	}
	// A junk cancel binding still reports the config bug, standby or not.
	junk := registerAllWith(0, HotkeyConfig{Summon: "Ctrl+Alt+Q", Cancel: "Ctrl+Alt+NotAKey+"},
		(&fakeRegistry{}).register)
	if c, _ := junk.Binding(hkCancel); c.Status != HotkeyUnparsable || c.Status.Attempted() {
		t.Errorf("junk cancel = %+v, want unparsable-and-never-attempted", c)
	}
	if p := strings.Join(junk.Problems(), "\n"); !strings.Contains(p, "cancel") {
		t.Errorf("the junk cancel binding vanished from the user-visible problems: %q", p)
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

	// summon survives a neighbour's failure. Cancel used to be asserted here as
	// "also live"; since ticket 245 an idle ball must NOT hold cancel, so the
	// independence-of-the-cancel-path property is asserted where it now lives:
	// after the same broken set, the borrow still gets its registration (the
	// cancel slot does not inherit its neighbours' failures).
	if !rep.IsLive(hkSummon) {
		t.Fatalf("one failure dropped the others: %+v", rep.Bindings())
	}
	if c, _ := rep.Binding(hkCancel); c.Status != HotkeyStandby {
		t.Errorf("cancel while neighbours failed = %q, want standby (idle balls never bind it)", c.Status)
	}
	borrowReg := &fakeRegistry{}
	if err := takeEscWith(borrowReg.unregister, borrowReg.register); err != nil {
		t.Errorf("the cancel path is not independent: borrow failed after neighbours failed (%v)", err)
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
//
// Ticket 245 adds a second reason for the cancel id specifically: an Esc the
// ball borrowed during Confirming is a slot this process holds even though the
// idle pass never registered it, so shutdown and rebind must still drop it -
// otherwise the borrowed key outlives the window that borrowed it.
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
// Every status the code can produce is listed here on purpose: a new status
// with no line in this table is a status nobody can decode from a log.
func TestHotkeyStatusString(t *testing.T) {
	cases := map[HotkeyStatus]string{
		HotkeyLive:       "live",
		HotkeyDisabled:   "disabled (unset in config)",
		HotkeyUnparsable: "not attempted (binding cannot be parsed)",
		HotkeyTaken:      "occupied by another app",
		HotkeyError:      "registration failed",
		HotkeyStandby:    "not bound while idle (cancel is borrowed only during Confirming)",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("status %d = %q, want %q", s, got, want)
		}
	}
	// The two never-attempted-but-different cases must not read the same way:
	// one is the user's choice, the other is this process's design rule.
	if HotkeyDisabled.String() == HotkeyStandby.String() {
		t.Error("disabled and standby share a sentence")
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
