package ball

import "testing"

// TestParseAccelerator covers the [hotkey] binding grammar.
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

// TestDefaultHotkeys pins the SPEC-03 §4 default set (cancel = Esc, B1).
func TestDefaultHotkeys(t *testing.T) {
	d := DefaultHotkeys()
	if d.Cancel != "Esc" {
		t.Errorf("cancel default = %q, want Esc", d.Cancel)
	}
	if d.Summon == "" || d.Mute == "" || d.Panel == "" {
		t.Errorf("summon/mute/panel defaults must be non-empty: %+v", d)
	}
}
