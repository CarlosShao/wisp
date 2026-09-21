//go:build windows

package ball

import "testing"

// TestParseAccelerator covers the [hotkey] binding grammar.
//
// This file is //go:build windows because its whole subject is Win32: the
// MOD_/VK_ constants, modNoRepeat, DefaultHotkeys, AltSummonSpace and
// ParseAccelerator are all declared in hotkey_windows.go, so an untagged test
// here cannot type-check on any other GOOS (ticket 78 - it was the second
// error hidden behind undefined: mulA in the same package).
func TestParseAccelerator(t *testing.T) {
	cases := []struct {
		in   string
		mods uint32
		vk   uint32
		fail bool
	}{
		{"Ctrl+Alt+W", modControl | modAlt, 'W', false},
		{"Ctrl+Alt+w", modControl | modAlt, 'W', false}, // case-insensitive key
		{"Win+Shift+A", modWin | modShift, 'A', false},
		{"Esc", 0, vkEscape, false},
		{"Escape", 0, vkEscape, false},
		{"F9", 0, 0x78, false},
		{"Ctrl+F12", modControl, 0x7B, false},
		{"Space", 0, vkSpace, false},
		{"7", 0, '7', false},
		{"", 0, 0, true},
		{"Ctrl+", 0, 0, true},
		{"Ctrl+Bogus", 0, 0, true},
	}
	for _, c := range cases {
		acc, err := ParseAccelerator(c.in)
		if c.fail {
			if err == nil {
				t.Errorf("ParseAccelerator(%q): expected error", c.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseAccelerator(%q): %v", c.in, err)
			continue
		}
		if acc.VK != c.vk {
			t.Errorf("ParseAccelerator(%q): vk=0x%X want 0x%X", c.in, acc.VK, c.vk)
		}
		if got := acc.Mods &^ modNoRepeat; got != c.mods {
			t.Errorf("ParseAccelerator(%q): mods=0x%X want 0x%X", c.in, got, c.mods)
		}
		if acc.Mods&modNoRepeat == 0 {
			t.Errorf("ParseAccelerator(%q): no-repeat flag missing", c.in)
		}
	}
}

// TestDefaultHotkeys pins SPEC-03 §4 (cancel = Esc, B1) and owner ruling R10
// (2026-09-20): the default summon key is Ctrl+Alt+Q, because Ctrl+Alt+W was
// measured occupied on the owner's machine (the only one of 84 candidates).
// Ctrl+Alt+Space is the documented, configurable alternative and stays NOT the
// default (some Chinese IMEs grab ctrl/shift+space before the app can).
func TestDefaultHotkeys(t *testing.T) {
	d := DefaultHotkeys()
	if d.Cancel != "Esc" {
		t.Errorf("cancel default = %q, want Esc", d.Cancel)
	}
	if d.Summon != "Ctrl+Alt+Q" {
		t.Errorf("summon default = %q, want Ctrl+Alt+Q (R10)", d.Summon)
	}
	if d.Summon == "Ctrl+Alt+W" || d.Summon == AltSummonSpace {
		t.Errorf("summon default regressed to %q", d.Summon)
	}
	if _, err := ParseAccelerator(AltSummonSpace); err != nil {
		t.Errorf("the documented alternative %q must stay parseable: %v", AltSummonSpace, err)
	}
	if d.Mute == "" || d.Panel == "" {
		t.Errorf("summon/mute/panel defaults must be non-empty: %+v", d)
	}
	for _, s := range []string{d.Summon, d.Mute, d.Cancel, d.Panel} {
		if _, err := ParseAccelerator(s); err != nil {
			t.Errorf("default %q must parse: %v", s, err)
		}
	}
}
