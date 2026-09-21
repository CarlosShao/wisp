package winsec

// Ticket 108 AC#2, the lexical half: which characters split a path, and what the
// ancestor list looks like for each spelling shape. Pure string facts, so this
// file compiles and runs on both platforms and the POSIX rows are read on a real
// Linux kernel rather than inferred (the same judgement the container run in the
// ticket's AC#5 row makes for the filesystem-shaped cases).
//
// The invariant every row checks first: each prefix must be an exact substring of
// the input. A prefix that is not a substring is a rebuilt path, and rebuilding is
// where a folded name (`a\b` read as `a` + `b`) turns into a different object -
// the fail-open direction this repository booked as A74(3).

import (
	"runtime"
	"strings"
	"testing"
)

func TestAC2AncestorPrefixesForEverySeparatorShape(t *testing.T) {
	onWindows := runtime.GOOS == "windows"
	type row struct {
		name        string
		input       string
		wantWindows []string
		wantPOSIX   []string
	}
	rows := []row{
		{
			name:        "native-backslash",
			input:       `C:\a\b\keep.txt`,
			wantWindows: []string{`C:\a`, `C:\a\b`, `C:\a\b\keep.txt`},
			wantPOSIX:   []string{`C:\a\b\keep.txt`},
		},
		{
			name:        "all-forward-slash",
			input:       `C:/a/b/keep.txt`,
			wantWindows: []string{`C:/a`, `C:/a/b`, `C:/a/b/keep.txt`},
			// On POSIX this whole string is one directory name with slashes in it,
			// so the first piece is the segment "C:" - a component, not a volume.
			wantPOSIX: []string{`C:`, `C:/a`, `C:/a/b`, `C:/a/b/keep.txt`},
		},
		{
			name:        "mixed-separators",
			input:       `C:\a/b\keep.txt`,
			wantWindows: []string{`C:\a`, `C:\a/b`, `C:\a/b\keep.txt`},
			wantPOSIX:   []string{`C:\a`, `C:\a/b\keep.txt`},
		},
		{
			name:        "trailing-separator",
			input:       `/a/b/keep.txt/`,
			wantWindows: []string{`/a`, `/a/b`, `/a/b/keep.txt`},
			wantPOSIX:   []string{`/a`, `/a/b`, `/a/b/keep.txt`},
		},
		{
			name:        "doubled-separator",
			input:       `/a/b//keep.txt`,
			wantWindows: []string{`/a`, `/a/b`, `/a/b//keep.txt`},
			wantPOSIX:   []string{`/a`, `/a/b`, `/a/b//keep.txt`},
		},
		{
			name:        "posix-backslash-is-a-name-character",
			input:       `/a/b\c/keep.txt`,
			wantWindows: []string{`/a`, `/a/b`, `/a/b\c`, `/a/b\c/keep.txt`},
			wantPOSIX:   []string{`/a`, `/a/b\c`, `/a/b\c/keep.txt`},
		},
		{
			name: "extended-length-prefix",
			// filepath.VolumeName takes \\?\C: as the volume on Windows, so the
			// device-namespace root is not offered as an ancestor to Lstat; on
			// POSIX the whole string is one component, because none of it is a
			// separator there. Either way the leaf stays out of the ancestor set.
			input:       `\\?\C:\a\b\keep.txt`,
			wantWindows: []string{`\\?\C:\a`, `\\?\C:\a\b`, `\\?\C:\a\b\keep.txt`},
			wantPOSIX:   []string{`\\?\C:\a\b\keep.txt`},
		},
	}
	for _, r := range rows {
		t.Run(r.name, func(t *testing.T) {
			want := r.wantPOSIX
			if onWindows {
				want = r.wantWindows
			}
			got := pathPieces(r.input)
			t.Logf("%s: pathPieces(%q) = %q", r.name, r.input, got)
			if strings.Join(got, "|") != strings.Join(want, "|") {
				t.Errorf("AC#2: %s on %s: got %q, want %q", r.name, runtime.GOOS, got, want)
			}
			for _, piece := range got {
				if !strings.HasPrefix(r.input, piece) {
					t.Errorf("AC#2 RED: prefix %q of %q is not an exact substring of the input, so a name got folded or rebuilt", piece, r.input)
				}
			}
		})
	}
}

func TestAC2ComponentsAndTraversalPerSeparatorShape(t *testing.T) {
	onWindows := runtime.GOOS == "windows"
	// The dot/space scan and the ".." scan must use the same separator set as the
	// ancestor walk, or one of the two walks disagrees about what a component is.
	comps := pathComponents(`C:\a/b\keep.txt`)
	wantComps := []string{`C:\a`, `b\keep.txt`}
	if onWindows {
		wantComps = []string{"a", "b", "keep.txt"}
	}
	if strings.Join(comps, "|") != strings.Join(wantComps, "|") {
		t.Errorf("AC#2: pathComponents(`C:\\a/b\\keep.txt`) = %q, want %q on %s", comps, wantComps, runtime.GOOS)
	}

	type traversal struct {
		input       string
		wantOnWin   bool
		wantOnPOSIX bool
		why         string
	}
	for _, tc := range []traversal{
		{`C:/a/../b`, true, true, "a dot-dot the Windows parser would honour"},
		{`/a/../b`, true, true, "a portable dot-dot"},
		{`/a/b\..\c`, true, false, "backslash-dot-dot: one real component name on POSIX, two separators on Windows"},
		{`/a/b//c`, true, true, "an empty segment"},
		{`/a/./b`, true, true, "a self reference"},
	} {
		want := tc.wantOnPOSIX
		if onWindows {
			want = tc.wantOnWin
		}
		got, found := lexicalTraversal(tc.input)
		t.Logf("%s: lexicalTraversal(%q) -> found=%v (%s)", tc.why, tc.input, found, got)
		if found != want {
			t.Errorf("AC#2: lexicalTraversal(%q) found=%v, want %v on %s (%s)", tc.input, found, want, runtime.GOOS, tc.why)
		}
	}
}
