//go:build windows

package ball

// Ticket 260 · 260-r2: the EXPECTATION half of ticket 245's live "the cancel key
// really was borrowed" ruler.
//
// The winlive helper requireEscBorrowed (hotkey_live_test.go:90) used to spell its
// expectation as a constant: the borrowed slot must hold VK_ESCAPE. Ticket 260 形ⓐ
// moved which key the borrow registers to [hotkey] cancel (resolveCancelBorrow,
// hotkey_windows.go:588), so for any config that is not the shipped default that
// constant is no longer the behaviour - it is a guarantee about an old build.
//
// The ruling is "the expectation follows the seed", which needs two cells and not
// one sentence, and both cells must be readable WITHOUT taking the owner's desktop
// (the winlive tier is not run this round):
//
//	cell ① (condition ①, zero drift): seed unset / "Esc"  =>  the borrow is bit-for
//	         bit {MOD_NOREPEAT, VK_ESCAPE} = the raw pair below, and the bare-Esc
//	         desktop probe still has to say "ours";
//	cell ② (AC#1): a non-default seed => the borrow is THAT combination, the probe
//	         runs against that combination, and the bare Esc has to be free again.
//
// So the helper below derives the expectation from the seed, and these tests pin
// that derivation to raw Win32 bit patterns (not to the constants production
// passes, which would make "the default moved" and "the expectation moved"
// indistinguishable - same reason accelDefault260 exists in ticket 260-r1's file)
// and to what the deterministic fake registry actually saw. That is the positive
// control the repo requires of every negative ruler: seed a different combination
// and the named step has to move with it.

import "testing"

// wantDefaultBorrow260r2 is the shipped default cancel borrow spelled as raw
// Win32 bits (MOD_NOREPEAT 0x4000, VK_ESCAPE 0x1B). It is deliberately NOT built
// from the named constants: a ruler that reads its expected value out of the same
// constants the code under test uses cannot see those constants moving.
var wantDefaultBorrow260r2 = Accelerator{Mods: 0x4000, VK: 0x1B}

// borrowWish260r2 is one seed's expectation: the accelerator the borrowed cancel
// slot must be seen holding while Confirming owns the key, and whether that seed
// is the default case (which decides which desktop probe the live ruler runs).
type borrowWish260r2 struct {
	Seed    string
	Acc     Accelerator
	Default bool  // the seed resolves to the shipped default bare Esc
	Refused error // non-nil: the seed is not a parsable combination, so the default is what gets borrowed
}

// wantCancelBorrow260r2 turns a configured [hotkey] cancel binding into the
// expectation the live borrow rulers compare against. It is the ruler-side mirror
// of production's resolveCancelBorrow and is kept deliberately independent of it:
// the value comes from ParseAccelerator (the binding-spelling contract, itself
// pinned by hotkey_test.go) over the seed string, never from the borrow helper.
// If production ever goes back to ignoring the config, the seed stays "Ctrl+Alt+V"
// while the borrowed slot reads Esc, and that is a red here, not an agreement.
func wantCancelBorrow260r2(bind string) borrowWish260r2 {
	if bind == "" {
		return borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: true}
	}
	acc, err := ParseAccelerator(bind)
	if err != nil {
		return borrowWish260r2{Seed: bind, Acc: wantDefaultBorrow260r2, Default: true, Refused: err}
	}
	return borrowWish260r2{Seed: bind, Acc: acc, Default: acc == wantDefaultBorrow260r2}
}

// TestBorrowWishFollowsSeed260r2 is the positive control for the synced live
// ruler, runnable in the default suite: two seeds, two different expectations,
// each spelled as its own bit pattern, and each matching what the deterministic
// borrow actually hands RegisterHotKey. "The hotkey is wired up" is not a verdict
// that can carry this; the named value has to move when the seed moves.
func TestBorrowWishFollowsSeed260r2(t *testing.T) {
	cases := []struct {
		name  string
		seed  string
		want  Accelerator
		isDef bool
	}{
		{"unset in config", "", Accelerator{Mods: 0x4000, VK: 0x1B}, true},
		{"the shipped default spelling", "Esc", Accelerator{Mods: 0x4000, VK: 0x1B}, true},
		{"the winlive fixture seed at hotkey_live_test.go:42", "Ctrl+Alt+V", Accelerator{Mods: 0x4003, VK: 0x56}, false},
		{"a different configured combination", "Ctrl+Alt+Y", Accelerator{Mods: 0x4003, VK: 0x59}, false},
		{"a bare function key", "F9", Accelerator{Mods: 0x4000, VK: 0x78}, false},
	}
	for _, c := range cases {
		wish := wantCancelBorrow260r2(c.seed)
		if wish.Acc != c.want {
			t.Errorf("%s: seed %q expects the borrowed slot to read %+v, got %+v", c.name, c.seed, c.want, wish.Acc)
			continue
		}
		if wish.Default != c.isDef {
			t.Errorf("%s: seed %q Default = %v, want %v", c.name, c.seed, wish.Default, c.isDef)
		}
		if wish.Refused != nil {
			t.Errorf("%s: a valid seed reported a refused binding: %v", c.name, wish.Refused)
		}

		// The expectation must be the thing production actually registers, or the
		// live ruler is measuring a formula nobody uses.
		reg := &fakeRegistry{}
		if err := takeEscWithAcc(reg.unregister, reg.register, resolveCancelBorrow(c.seed)); err != nil {
			t.Fatalf("%s: the fake registry refused the borrow: %v", c.name, err)
		}
		if got := lastAttempt260(t, reg); got != wish.Acc {
			t.Errorf("%s: RegisterHotKey was handed %+v while the ruler expects %+v", c.name, got, wish.Acc)
		}
		if got := reg.bindsEsc(); got != c.isDef {
			t.Errorf("%s: the borrow bound a bare Esc = %v, want %v for seed %q", c.name, got, c.isDef, c.seed)
		}
	}

	// The two cells are two cells: the default expectation and the configured one
	// must not collapse into one value (that collapse is exactly what the pre-260
	// constant was, and it is what this round is replacing).
	if wantCancelBorrow260r2("Esc").Acc == wantCancelBorrow260r2("Ctrl+Alt+Y").Acc {
		t.Fatal("the expectation does not follow the seed: \"Esc\" and \"Ctrl+Alt+Y\" resolve to the same accelerator")
	}
}

// TestBorrowWishDefaultCellIsBitIdentical260r2 is condition ① at the expectation
// layer: every way the default arrives must still be {0x4000, 0x1B}, and that pair
// must still be what the DEFAULT-档 wrapper takeEscWith registers - so the
// synced live ruler's default branch says the same thing the pre-260 ruler said.
func TestBorrowWishDefaultCellIsBitIdentical260r2(t *testing.T) {
	if wantDefaultBorrow260r2 != (Accelerator{Mods: 0x4000, VK: 0x1B}) {
		t.Fatalf("the raw default pair moved to %+v; condition ①'s ruler is measuring a value nobody registers",
			wantDefaultBorrow260r2)
	}
	if escBorrowAcc() != wantDefaultBorrow260r2 {
		t.Errorf("escBorrowAcc() = %+v, want the raw default pair %+v - the default档 drifted",
			escBorrowAcc(), wantDefaultBorrow260r2)
	}
	seeds := []string{"", "Esc", DefaultHotkeys().Cancel, ApplyHotkeyDefaults(HotkeyConfig{}).Cancel}
	for _, seed := range seeds {
		wish := wantCancelBorrow260r2(seed)
		if !wish.Default || wish.Acc != wantDefaultBorrow260r2 {
			t.Fatalf("seed %q resolves to %+v (Default=%v), want the bare-Esc default pair %+v",
				seed, wish.Acc, wish.Default, wantDefaultBorrow260r2)
		}
	}
	// The pre-260 shape, through the wrapper production no longer uses, so
	// "unchanged" is checked against the value the old code returned.
	reg := &fakeRegistry{}
	if err := takeEscWith(reg.unregister, reg.register); err != nil {
		t.Fatalf("the default borrow was refused by the fake registry: %v", err)
	}
	if got := lastAttempt260(t, reg); got != wantCancelBorrow260r2(DefaultHotkeys().Cancel).Acc {
		t.Errorf("the default-档 borrow registered %+v while the synced ruler expects %+v",
			got, wantCancelBorrow260r2(DefaultHotkeys().Cancel).Acc)
	}
}

// TestBorrowWishRefusedSeedIsNotSilent260r2 keeps the third reading straight: a
// seed the parser refuses is not "unset" - the ruler expects the default pair AND
// carries the refusal, so the live helper can require the receipt to say the
// substitution out loud (condition ③) instead of nodding at a silent fallback.
func TestBorrowWishRefusedSeedIsNotSilent260r2(t *testing.T) {
	const junk = "Ctrl+Alt+NotAKey+"
	wish := wantCancelBorrow260r2(junk)
	if wish.Refused == nil {
		t.Fatalf("seed %q parsed cleanly; the ruler would then expect the user's key while production borrows %+v",
			junk, wish.Acc)
	}
	if wish.Acc != wantDefaultBorrow260r2 || !wish.Default {
		t.Errorf("the refused seed expects %+v (Default=%v), want the default pair %+v", wish.Acc, wish.Default, wantDefaultBorrow260r2)
	}
	// Positive control on the same ruler: the seeds the winlive fixtures use parse,
	// so the line above is a reaction to junk and not a standing announcement.
	for _, seed := range []string{"", "Esc", "Ctrl+Alt+V"} {
		if got := wantCancelBorrow260r2(seed); got.Refused != nil {
			t.Errorf("seed %q reported a refusal: %v", seed, got.Refused)
		}
	}
	// And the expectation still matches what production borrows for a junk seed.
	reg := &fakeRegistry{}
	if err := takeEscWithAcc(reg.unregister, reg.register, resolveCancelBorrow(ApplyHotkeyDefaults(HotkeyConfig{Cancel: junk}).Cancel)); err != nil {
		t.Fatalf("the fallback borrow was refused by the fake registry: %v", err)
	}
	if got := lastAttempt260(t, reg); got != wish.Acc {
		t.Errorf("a junk seed: RegisterHotKey was handed %+v, want the default pair %+v", got, wish.Acc)
	}
}
