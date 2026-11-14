//go:build windows

package ball

// Ticket 260 AC#1, 形ⓐ (the Confirming borrow reads [hotkey] cancel instead of a
// hard-coded key). These are the four resident rulers the orchestrator's ruling
// names as its conditions, plus the two boundary guards the ruling forbids
// crossing:
//
//	① the default档 did not move one bit - an unset or default slot still borrows
//	  {MOD_NOREPEAT, VK_ESCAPE}, asserted as an Accelerator value and never as a
//	  string that happens to contain "Esc";
//	② the receipt line and the accelerator handed to RegisterHotKey come from ONE
//	  resolution (resolveCancelBorrow), so the card can never name a key that was
//	  not the one borrowed;
//	③ a configured binding the parser refuses falls back to the default AND says
//	  so on the receipt (P6) - no swallowed parse failure;
//	④ plus the "借的时机与存在性" guards from ticket 245, in the shape the 2026-11-14
//	  AC#6② ruling settled: the idle pass still never registers a BARE cancel
//	  binding, a modifier-combination one is held while idle, and the borrow still
//	  happens exactly once per card for whatever binding the config says.
//
// Deterministic layer: the fake registry is the one ticket 64/245 built in
// hotkey_status_test.go, so these run in the default suite and take no desktop.

import (
	"errors"
	"strings"
	"syscall"
	"testing"
)

// accelDefault260 is condition ①'s expected value spelled as raw Win32 bits, not
// through the constants the production code uses - otherwise "the default moved"
// and "the constant moved" would be indistinguishable, which is the drift this
// condition exists to prevent.
var accelDefault260 = Accelerator{Mods: 0x4000, VK: 0x1B} // MOD_NOREPEAT, VK_ESCAPE

// lastAttempt260 returns the accelerator the fake registry saw last. In every
// fixture here the borrow is the final RegisterHotKey attempt.
func lastAttempt260(t *testing.T, reg *fakeRegistry) Accelerator {
	t.Helper()
	if len(reg.attemptedAcc) == 0 {
		t.Fatal("the fake registry saw no RegisterHotKey attempt at all")
	}
	return reg.attemptedAcc[len(reg.attemptedAcc)-1]
}

// cleanReport260 is an idle report over cfg with its cancel line replaced, which
// is the shape TakeEscForCancel leaves behind on the real ball.
func cleanReport260(cfg HotkeyConfig, line HotkeyBinding) HotkeyReport {
	return registerAllWith(0, cfg, (&fakeRegistry{}).register).withCancel(line)
}

// TestBorrowDefaultIsBitIdentical260 is condition ①: the shipped default borrow is
// the same accelerator it was before ticket 260, bit for bit, for every way the
// default can arrive (unset in config, the product default, the schema default run
// through ApplyHotkeyDefaults).
func TestBorrowDefaultIsBitIdentical260(t *testing.T) {
	// The numeric literals above must still be what the constants say, or the
	// whole test is measuring a pair nobody registers.
	if modNoRepeat != 0x4000 || vkEscape != 0x1B {
		t.Fatalf("modNoRepeat=0x%X vkEscape=0x%X, want 0x4000/0x1B - condition ①'s literals went stale",
			modNoRepeat, vkEscape)
	}
	if acc, err := ParseAccelerator("Esc"); err != nil || acc != accelDefault260 {
		t.Fatalf(`ParseAccelerator("Esc") = %+v (%v), want %+v - the fact 260-a2 measured moved`, acc, err, accelDefault260)
	}

	for _, bind := range []string{"", DefaultHotkeys().Cancel, ApplyHotkeyDefaults(HotkeyConfig{}).Cancel} {
		borrow := resolveCancelBorrow(bind)
		if borrow.Acc != accelDefault260 {
			t.Errorf("cancel = %q borrows %+v, want %+v (mods|VK), the default档 drifted", bind, borrow.Acc, accelDefault260)
		}
		if borrow.Fallback != nil {
			t.Errorf("cancel = %q reports a fallback (%v); the default is not a fallback", bind, borrow.Fallback)
		}

		reg := &fakeRegistry{}
		if err := takeEscWithAcc(reg.unregister, reg.register, borrow); err != nil {
			t.Fatalf("cancel = %q: the default borrow was refused by the fake registry: %v", bind, err)
		}
		if got := lastAttempt260(t, reg); got != accelDefault260 {
			t.Errorf("RegisterHotKey was handed %+v, want %+v", got, accelDefault260)
		}
		if !reg.bindsEsc() {
			t.Fatal("the ruler is blind: the fake registry never saw a bare Esc, so nothing here proves the default binds one")
		}

		line := cancelBorrowedLineFor(borrow)
		if line.Acc != accelDefault260 || line.Binding != "Esc" || line.Status != HotkeyLive {
			t.Errorf("default receipt line = %+v, want live \"Esc\" with %+v", line, accelDefault260)
		}
		if line.Note != "" {
			t.Errorf("the default receipt carries a fallback sentence (%q); condition ① says nothing about the default may change, wording included", line.Note)
		}
		if p := cleanReport260(ApplyHotkeyDefaults(HotkeyConfig{}), line).Problems(); len(p) != 0 {
			t.Errorf("a default borrow produced user-visible problems: %v", p)
		}
	}

	// The pre-260 shape itself, so the claim "unchanged" is checked against the
	// value the old code returned rather than against this file's own wording.
	if escBorrowAcc() != accelDefault260 {
		t.Errorf("escBorrowAcc() = %+v, want %+v", escBorrowAcc(), accelDefault260)
	}
	if got := resolveCancelBorrow(""); got.Acc != escBorrowAcc() {
		t.Errorf("unset cancel borrows %+v, want escBorrowAcc() %+v", got.Acc, escBorrowAcc())
	}
}

// TestBorrowFollowsConfiguredBinding260 is condition ①'s counterpart and AC#1's
// actual proof of life: a configured cancel binding is the key that gets
// borrowed. Before ticket 260 every case below borrowed the bare Esc and the
// configuration never reached RegisterHotKey.
func TestBorrowFollowsConfiguredBinding260(t *testing.T) {
	cases := []struct {
		name string
		bind string
		want Accelerator
	}{
		{"modifiers plus a letter", "Ctrl+Alt+K", Accelerator{Mods: 0x4000 | 0x0002 | 0x0001, VK: 'K'}},
		{"the documented alternative summon shape", "Ctrl+Alt+Space", Accelerator{Mods: 0x4000 | 0x0002 | 0x0001, VK: 0x20}},
		{"a bare function key", "F9", Accelerator{Mods: 0x4000, VK: 0x78}},
	}
	for _, c := range cases {
		borrow := resolveCancelBorrow(ApplyHotkeyDefaults(HotkeyConfig{Cancel: c.bind}).Cancel)
		if borrow.Acc != c.want {
			t.Errorf("%s: cancel = %q resolves to %+v, want %+v", c.name, c.bind, borrow.Acc, c.want)
			continue
		}
		if borrow.Acc == escBorrowAcc() {
			t.Errorf("%s: cancel = %q still resolves to the hard-coded bare Esc - the config did not reach the borrow", c.name, c.bind)
		}
		if borrow.Fallback != nil {
			t.Errorf("%s: a valid binding reported a fallback: %v", c.name, borrow.Fallback)
		}

		reg := &fakeRegistry{}
		if err := takeEscWithAcc(reg.unregister, reg.register, borrow); err != nil {
			t.Fatalf("%s: the fake registry refused the borrow: %v", c.name, err)
		}
		if got := lastAttempt260(t, reg); got != c.want {
			t.Errorf("%s: RegisterHotKey was handed %+v, want %+v", c.name, got, c.want)
		}
		if reg.bindsEsc() {
			t.Errorf("%s: the borrow bound a bare Esc although [hotkey] cancel says %q", c.name, c.bind)
		}

		cfg := ApplyHotkeyDefaults(HotkeyConfig{Cancel: c.bind})
		line := cancelBorrowedLineFor(borrow)
		if line.Binding != c.bind || line.Acc != c.want || line.Status != HotkeyLive || line.Note != "" {
			t.Errorf("%s: receipt line = %+v, want live %q with %+v and no note", c.name, line, c.bind, c.want)
		}
		if p := cleanReport260(cfg, line).Problems(); len(p) != 0 {
			t.Errorf("%s: a configured borrow that worked produced user-visible problems: %v", c.name, p)
		}
	}
}

// TestBorrowUnparsableFallsBackLoudly260 is condition ③ (part P6): the branch that
// did not exist before this ticket. A cancel binding the parser refuses must
// (a) still borrow - ticket 245's ruling is about WHEN, and dropping the borrow
// would change it, and (b) say the substitution on the receipt in the
// user-visible report. Nothing here skips.
func TestBorrowUnparsableFallsBackLoudly260(t *testing.T) {
	const junk = "Ctrl+Alt+NotAKey+"
	cfg := ApplyHotkeyDefaults(HotkeyConfig{Cancel: junk})
	if cfg.Cancel != junk {
		t.Fatalf("premise moved: ApplyHotkeyDefaults turned %q into %q", junk, cfg.Cancel)
	}
	borrow := resolveCancelBorrow(cfg.Cancel)
	if borrow.Fallback == nil {
		t.Fatal("P6 RED: an unparsable cancel binding resolved with no fallback notice - the parse failure was swallowed")
	}
	if borrow.Acc != accelDefault260 || borrow.Binding != "Esc" {
		t.Errorf("the fallback is not the documented default Esc: %+v / %q", borrow.Acc, borrow.Binding)
	}

	reg := &fakeRegistry{}
	if err := takeEscWithAcc(reg.unregister, reg.register, borrow); err != nil {
		t.Fatalf("the fallback borrow was refused by the fake registry: %v", err)
	}
	if n := reg.attemptsOf(hkCancel); n != 1 {
		t.Errorf("cancel registration attempts = %d, want exactly 1: a broken binding must still borrow the default key (ticket 245's timing, unchanged)", n)
	}
	if got := lastAttempt260(t, reg); got != accelDefault260 {
		t.Errorf("the fallback bound %+v, want the default %+v", got, accelDefault260)
	}

	line := cancelBorrowedLineFor(borrow)
	if line.Status != HotkeyLive {
		t.Errorf("the fallback line is %q, want live - the borrow did happen", line.Status)
	}
	if line.Note == "" {
		t.Fatal("P6 RED: the fallback borrow left no sentence on the receipt line")
	}
	probs := cleanReport260(cfg, line).Problems()
	if len(probs) != 1 {
		t.Fatalf("user-visible problems after a fallback borrow = %d lines, want the 1 that names it: %q", len(probs), probs)
	}
	said := strings.Join(probs, "\n")
	for _, want := range []string{"cancel", junk, "not a valid key combination", "default Esc", "borrowed for this card"} {
		if !strings.Contains(said, want) {
			t.Errorf("the fallback sentence does not say %q; it reads: %s", want, said)
		}
	}

	// Positive control on the same ruler: a default borrow says nothing, so the
	// line above is a reaction to the junk binding and not a standing announcement.
	clean := cleanReport260(DefaultHotkeys(), cancelBorrowedLineFor(resolveCancelBorrow(DefaultHotkeys().Cancel)))
	if p := clean.Problems(); len(p) != 0 {
		t.Errorf("the default borrow printed a fallback sentence: %v", p)
	}
	// And the idle half of the same junk config keeps reporting its own bug the
	// way ticket 245 left it (this change did not hush the idle line).
	if p := registerAllWith(0, cfg, (&fakeRegistry{}).register).Problems(); len(p) == 0 ||
		!strings.Contains(p[0], junk) {
		t.Errorf("the idle pass no longer names the junk binding: %q", p)
	}
}

// TestBorrowReceiptSharesOneSourceWithRegistration260 is condition ②: the receipt
// and the RegisterHotKey call read the SAME resolved value. The two assertions
// together close the "two spellings,拼两次" shape: the line's accelerator must be
// the one Win32 was handed, and the line's printed spelling must parse back to
// that same accelerator.
func TestBorrowReceiptSharesOneSourceWithRegistration260(t *testing.T) {
	for _, bind := range []string{"", "Esc", "Ctrl+Alt+K", "Ctrl+Alt+Space", "F9", "Ctrl+Alt+NotAKey+"} {
		borrow := resolveCancelBorrow(bind)
		reg := &fakeRegistry{}
		if err := takeEscWithAcc(reg.unregister, reg.register, borrow); err != nil {
			t.Fatalf("cancel = %q: borrow refused: %v", bind, err)
		}
		line := cancelBorrowedLineFor(borrow)

		if got := lastAttempt260(t, reg); line.Acc != got {
			t.Fatalf("cancel = %q: the receipt carries %+v while RegisterHotKey was handed %+v - the two are not the same resolution",
				bind, line.Acc, got)
		}
		parsed, err := ParseAccelerator(line.Binding)
		if err != nil {
			t.Fatalf("cancel = %q: the receipt prints %q, which does not even parse: %v", bind, line.Binding, err)
		}
		if parsed != line.Acc {
			t.Fatalf("cancel = %q: the receipt prints %q = %+v but registers %+v - the page and the key disagree",
				bind, line.Binding, parsed, line.Acc)
		}
		wantSpelling := bind
		if borrow.Fallback != nil || bind == "" {
			wantSpelling = "Esc"
		}
		if line.Binding != wantSpelling {
			t.Errorf("cancel = %q: receipt spells %q, want %q", bind, line.Binding, wantSpelling)
		}
	}
}

// TestBorrowFailureNamesTheKeyItTried260 keeps condition ② honest on the failure
// half: when Win32 refuses a configured borrow, the line that says so must name
// the key that was really attempted, not the default one.
func TestBorrowFailureNamesTheKeyItTried260(t *testing.T) {
	const want = "Ctrl+Alt+K"
	borrow := resolveCancelBorrow(want)
	reg := &fakeRegistry{fail: map[uint32]error{hkCancel: syscall.EINVAL}}
	err := takeEscWithAcc(reg.unregister, reg.register, borrow)
	if err == nil {
		t.Fatal("a failing registrar reported a successful borrow")
	}
	if got := lastAttempt260(t, reg); got != borrow.Acc {
		t.Fatalf("the refused attempt was for %+v, want the configured %+v", got, borrow.Acc)
	}

	line := cancelFailedLineFor(borrow, err)
	if line.Binding != want || line.Acc != borrow.Acc || line.Status != HotkeyError || !line.Status.Attempted() {
		t.Fatalf("the failed-borrow line = %+v, want an attempted error naming %q", line, want)
	}
	probs := strings.Join(cleanReport260(ApplyHotkeyDefaults(HotkeyConfig{Cancel: want}), line).Problems(), "\n")
	if !strings.Contains(probs, want) || !strings.Contains(probs, "was not registered") {
		t.Errorf("the refused borrow does not name the key it tried: %q", probs)
	}

	// Positive control: the default-档 failure line still names Esc (the
	// pre-260 wording ticket 245's rulers read).
	defLine := cancelFailedLine(syscall.EINVAL)
	if defLine.Binding != "Esc" || defLine.Acc != accelDefault260 || defLine.Status != HotkeyError {
		t.Errorf("the default failed-borrow line moved: %+v", defLine)
	}
}

// TestConfiguredCancelIdleBindingFollowsModifierRuling245 is ticket 245 AC#6②'s
// criterion on the idle registration pass, in the shape the 2026-11-14 ruling
// settled. The old guard here said "a configured cancel is NEVER bound while
// idle"; the ruling narrowed that to BARE keys, because the defect the ticket is
// about is a bare key taken from the whole desktop, not a steady-state slot count
// of three. So both halves are judged, each with its own reading:
//
//   - a cancel binding with a real modifier (Ctrl+Alt+K) is HELD while idle:
//     attempted once, live in the very report the host prints
//     (HotkeyReport), and free of problem lines. A user who configures a
//     combination is asking for a permanent binding, and it steals nothing from
//     another application - which is also the direction ApplyHotkeyDefaults
//     documents for this slot;
//   - a BARE key (the shipped "Esc") is still standby: never attempted, still
//     carried on the report, and the borrow still binds it exactly once for
//     Confirming while the return still drops it without re-binding. That half is
//     today's behaviour, one word unchanged.
//
// Both halves read the seam the host consumes (HotkeyReport / EscTakenOver), not
// only the unexported slice, because that is where AC#6② is judged.
func TestConfiguredCancelIdleBindingFollowsModifierRuling245(t *testing.T) {
	// (1) The ruled-in case: a modifier combination is held while idle.
	const combo = "Ctrl+Alt+K"
	cfg := ApplyHotkeyDefaults(HotkeyConfig{Cancel: combo})
	reg := &fakeRegistry{}
	rep := registerAllWith(0, cfg, reg.register)
	if n := reg.attemptsOf(hkCancel); n != 1 {
		t.Fatalf("AC#6② RED: the idle pass made %d RegisterHotKey attempts for a modifier-combination cancel "+
			"(want exactly 1: a combination is held, not borrowed): %v", n, reg.attempted)
	}
	wantAcc, err := ParseAccelerator(combo)
	if err != nil {
		t.Fatalf("ParseAccelerator(%q): %v", combo, err)
	}
	c, ok := rep.Binding(hkCancel)
	if !ok || c.Status != HotkeyLive || c.Binding != combo || c.Acc != wantAcc {
		t.Fatalf("AC#6② RED: idle cancel line = %+v (present=%t), want live %q / %+v", c, ok, combo, wantAcc)
	}
	if !rep.IsLive(hkCancel) {
		t.Fatalf("AC#6② RED: the cancel slot is not in the idle live set: %+v", rep.Bindings())
	}
	if n := len(rep.Live()); n != 4 {
		t.Errorf("idle live set = %d, want 4 while a cancel combination is held: %+v", n, rep.Bindings())
	}
	if p := rep.Problems(); len(p) != 0 {
		t.Errorf("a held cancel combination reports problems: %q", p)
	}
	if _, held := idleBareKeyHeld245(rep); held {
		t.Errorf("the held combination was read as a bare key: %+v", rep.Bindings())
	}
	if reds := idleRosterReds245(rep); len(reds) != 0 {
		t.Errorf("the AC#6② nail reds on a held modifier combination: %v", reds)
	}
	live, taken := cancelIdleSeam245(t, rep, combo)
	if !live {
		t.Error("AC#6② RED: HotkeyReport().IsLive(cancel) is false for a configured combination")
	}
	if taken {
		t.Error("EscTakenOver() reads true at steady state: nothing was borrowed")
	}

	// (2) The unchanged case: the bare Esc is standby, never attempted.
	bare := ApplyHotkeyDefaults(HotkeyConfig{Cancel: "Esc"})
	reg2 := &fakeRegistry{}
	rep2 := registerAllWith(0, bare, reg2.register)
	if n := reg2.attemptsOf(hkCancel); n != 0 {
		t.Fatalf("ticket 245 RED: the idle pass attempted a BARE cancel key (%d attempts): %v", n, reg2.attempted)
	}
	if reg2.bindsEsc() {
		t.Fatalf("ticket 245 RED: the idle pass bound the bare Esc: %v", reg2.attemptedAcc)
	}
	bc, _ := rep2.Binding(hkCancel)
	if bc.Status != HotkeyStandby || bc.Binding != "Esc" || bc.Acc != (Accelerator{}) {
		t.Fatalf("bare cancel idle line = %+v, want standby \"Esc\" with a zero accelerator", bc)
	}
	if n := len(rep2.Live()); n != 3 {
		t.Errorf("idle live set with a bare cancel = %d, want 3: %+v", n, rep2.Bindings())
	}
	liveBare, takenBare := cancelIdleSeam245(t, rep2, "Esc")
	if liveBare {
		t.Error("ticket 245 RED: HotkeyReport().IsLive(cancel) is true for the bare Esc at steady state")
	}
	if takenBare {
		t.Error("EscTakenOver() reads true at steady state for the bare Esc")
	}
	if err := takeEscWithAcc(reg2.unregister, reg2.register, resolveCancelBorrow(bare.Cancel)); err != nil {
		t.Fatalf("the bare borrow was refused: %v", err)
	}
	if n := reg2.attemptsOf(hkCancel); n != 1 {
		t.Fatalf("the card borrowed the bare Esc %d times, want exactly 1", n)
	}
	releaseEscWith(reg2.unregister)
	if n := reg2.attemptsOf(hkCancel); n != 1 {
		t.Fatalf("after the return cancel attempts = %d, want still 1 - the return must drop the slot, "+
			"not re-bind the configured key (ticket 245)", n)
	}
	if n := reg2.unregistersOf(hkCancel); n < 1 {
		t.Fatal("the return never called UnregisterHotKey for the cancel slot")
	}
}

// cancelIdleSeam245 reads a finished idle report through the two accessors the
// host and the acceptance use: Ball.HotkeyReport() and Ball.EscTakenOver(). The
// ball is windowless and its ui-sta thread IS this goroutine, the same door the
// AC#7 rig measured (hotkey_refused_borrow_release_245_test.go), so the reads run
// inline and cost no desktop.
func cancelIdleSeam245(t *testing.T, seed HotkeyReport, cancelBinding string) (live, taken bool) {
	t.Helper()
	b := &Ball{hotkeyReport: seed, cancelBinding: cancelBinding}
	b.sta = newSTAThread(nil)
	b.sta.start(func(*staThread) error {
		live = b.HotkeyReport().IsLive(hkCancel)
		taken = b.EscTakenOver()
		return errT245NoWindowDoor
	})
	if err := b.sta.waitStarted(); !errors.Is(err, errT245NoWindowDoor) {
		t.Fatalf("the windowless ui-sta door did not take its create-error exit (err=%v): the reads above "+
			"would describe a pump that is not running", err)
	}
	return live, taken
}

// TestLiveBorrowHelperPremiseIsConfigDependent260 records a reading, not a fix: the
// winlive helper requireEscBorrowed (internal/ball/hotkey_live_test.go:90, whose
// VK assertion sits at :107) is fed by liveHotkeys() (:41), whose cancel binding is
// "Ctrl+Alt+V" - a value copied here because that file is the winlive layer and this
// round was told not to run it. Deterministically, that input now borrows VK 'V', so
// the helper's "the borrowed slot is VK_ESCAPE" premise is a property of the CONFIG,
// not of the code, and those four winlive callers would read red until the orchestrator
// rules on it. 260-r1 did not edit that file and did not run the winlive tier.
func TestLiveBorrowHelperPremiseIsConfigDependent260(t *testing.T) {
	// The fixture value is copied from liveHotkeys(), which this build cannot see.
	fixture := HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Ctrl+Alt+V", Panel: "Ctrl+Alt+B"}

	borrow := resolveCancelBorrow(fixture.Cancel)
	if borrow.Acc.VK == vkEscape {
		t.Fatalf("the borrow of %q resolved to VK_ESCAPE again - condition ①'s default-only rule was over-applied", fixture.Cancel)
	}
	if borrow.Acc.VK != 'V' || borrow.Acc.Mods != (0x4000|0x0002|0x0001) {
		t.Fatalf("the borrow of %q resolved to %+v, want {mods 0x4003, VK 0x56}", fixture.Cancel, borrow.Acc)
	}
	line := cancelBorrowedLineFor(borrow)
	if line.Acc.VK == vkEscape {
		t.Fatal("RegisterHotKey would get the configured key while the receipt still prints Esc")
	}
	t.Logf("named conflict for the orchestrator: hotkey_live_test.go:107 reads rep.Live()[hkCancel].VK == vkEscape "+
		"for a ball booted with liveHotkeys().Cancel = %q; under 形ⓐ that slot holds VK 0x%X, so the helper's "+
		"bare-Esc premise is now config-dependent (winlive tier not run this round)", fixture.Cancel, line.Acc.VK)
}
