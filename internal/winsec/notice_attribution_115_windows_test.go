//go:build windows

package winsec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// Ticket 115's own cases. They exist because run 35599458439 is the first CI run
// in which the C26 pipeline really got installed into the sealing seam (verbatim
// from that job's step 4, the winsec gate:
//
//	INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2
//
// where the previous run, 35595651898, printed "refusing to install" at the same
// point and every seal therefore acted on the caller's own spelling). With the
// resolver installed, a notice names the object in the OS's answer, and the four
// existing notice cases identify "their" notice by comparing the caller's spelling
// with strings.EqualFold or ==. On GitHub's runner the caller's spelling comes out
// of USERPROFILE/TEMP in 8.3 form (C:\Users\RUNNER~1), the answer comes out long
// (C:\Users\runneradmin), and those two are not a case difference:
//
//	inherited_narrow_notice_104_windows_test.go:134: reported 0 notice(s), want exactly 1;
//	  all notices: [{Path:C:\Users\runneradmin\AppData\Local\Temp\...readable-by-inheritance.txt
//	  Principals:[] Inherited:[S-1-1-0(A;ID;0x1200a9;;;WD)]}]
//
// The notice was there, correct, and unattributable. That is a comparison defect,
// and the ruling this file pins (ticket 115 AC#2) is "attribute by tree": the
// caller's spelling goes through the one resolver this package may use and the two
// sides are then compared by component (noticeNamesTree in winsec_windows.go).
// R-115-2 holds this file's own legs to that same rule, via answerNamesTree115:
// after the caller-side migration landed, the only two cases run 35608530583 still
// lost were two of these self-proof legs, each red because it held the caller's
// spelling and compared it with the answer by hand.
//
// Why these cases plant the alias instead of waiting for it: this box's own temp
// root needs no short name, so living on a machine cannot reproduce the runner's
// shape - ticket 112 established that and planted the same way (see
// tree_ownership_112_windows_test.go). Unlike that file, a plant that fails here
// is a Fatal rather than a Skip: these two cases are the only evidence AC#3 and
// AC#4 have, so a machine that cannot plant must say "this instrument measured
// nothing" out loud instead of reporting a green that tested nothing at all.

// captureNotices115 swaps the notice seam for the duration of one test.
func captureNotices115(t *testing.T) *[]narrowNotice {
	t.Helper()
	got := &[]narrowNotice{}
	orig := noticeNarrowed
	noticeNarrowed = func(n narrowNotice) { *got = append(*got, n) }
	t.Cleanup(func() { noticeNarrowed = orig })
	return got
}

// shortFormOf115 asks the OS for the 8.3 rendering of an existing object. No
// hand-built alias: the whole defect is about a spelling this package does not
// control, so the fixture uses the one the filesystem chose (ticket 112's rule).
func shortFormOf115(t *testing.T, path string) string {
	t.Helper()
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatal(err)
	}
	n, err := windows.GetShortPathName(p, nil, 0)
	if err != nil && n == 0 {
		t.Fatalf("GetShortPathName(%s) size probe: %v", path, err)
	}
	buf := make([]uint16, n)
	if _, err := windows.GetShortPathName(p, &buf[0], n); err != nil {
		t.Fatalf("GetShortPathName(%s): %v", path, err)
	}
	return windows.UTF16ToString(buf)
}

// widen115 puts an explicit out-of-band grant on an object, so the next seal of
// that object has something to report. Every call re-widens: the seal that
// follows clears it again and emits exactly one notice about exactly that object.
func widen115(t *testing.T, path string) {
	t.Helper()
	mustExec(t, "icacls", path, "/grant", "*"+everyoneSID+":(RX)")
}

// answerNamesTree115 asks this file's own question - "does the answer this notice
// carries name the same tree as the spelling the caller typed?" - and answers it
// by tree, which is the ruling ticket 115 AC#2 wrote down. The attribution itself
// is delegated to the package's single noticeNamesTree rather than re-derived
// here, so this file does not grow a comparison face beside the ones commit
// 527d303 migrated on the caller side.
//
// What the wrapper adds over calling that helper directly is the guard.
// noticeNamesTree answers false both for "this notice is not about that tree" and
// for "ResolvePath would not vouch for that spelling". At the production seam the
// collapse is the right direction - a caller asking "was my tree reported on?"
// hears "no" and goes looking. Inside a case it is not: this file's back half
// reads a false as "these are two different trees", so an unanswerable question
// could satisfy a rejection leg without having compared anything. Here a spelling
// nothing can name is a Fatal, which keeps both directions answerable-or-loud.
//
// It replaces the five faces in this file that held the caller's own spelling and
// compared it with the notice's answer: the strings.EqualFold in the forward half
// and the four bare sameTree calls in the back half (R-115-2 - on run
// 35608530583 those were the only two cases still red, and both are this file's
// self-proof legs, red because the caller's "long" spelling carries an 8.3
// segment there).
//
// The one literal comparison left in the file is the instrument leg that demands
// the notice NOT be fold-equal to the short spelling it was sealed through. It
// stays literal on purpose: its claim is "an 8.3 expansion happened at all", and
// answered by tree that claim is a tautology - the resolver's answer is always in
// the tree the short spelling names - so the leg would stop proving anything.
func answerNamesTree115(t *testing.T, n narrowNotice, spelling string) bool {
	t.Helper()
	if _, err := ResolvePath(spelling); err != nil {
		t.Fatalf("this leg cannot be asked at all: ResolvePath(%q) refused to vouch for the spelling it was handed: %v", spelling, err)
	}
	return noticeNamesTree(n, spelling)
}

// childTree115 is a private directory holding one artifact, both in the long
// spelling the filesystem gives them and in the spellings a caller can legitimately
// type for the same objects.
type childTree115 struct {
	dir      string   // long spelling of the directory
	dirShort string   // the same directory, in the OS's own 8.3 form
	child    string   // long spelling of the artifact
	spelling []string // every spelling of child that names the same object
}

// plantChildTree115 builds one such tree and proves the alias is real before
// anything is asserted over it (AC#3's self-proof leg).
func plantChildTree115(t *testing.T, name string) childTree115 {
	t.Helper()
	if r := PathResolverInstalled(); r == nil {
		t.Fatal("no resolver is installed in this test binary and the built-in floor rewrites nothing, " +
			"so a notice can only ever repeat the caller's spelling and this instrument would measure nothing")
	}
	base := t.TempDir()
	dir := filepath.Join(base, name)
	if err := PrivateDirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(dir, "artifact-under-test.txt")
	if err := os.WriteFile(child, []byte("artifact"), 0o600); err != nil {
		t.Fatal(err)
	}
	tr := childTree115{dir: dir, child: child}

	dirShort := shortFormOf115(t, dir)
	if strings.EqualFold(dirShort, dir) {
		t.Fatalf("the instrument planted nothing: %q has no 8.3 short form on this volume "+
			"(GetShortPathName answered %q), so the runner's RUNNER~1 asymmetry does not exist here to test against",
			dir, dirShort)
	}
	tr.dirShort = dirShort
	tr.spelling = []string{
		child,                    // the filesystem's own long form
		shortFormOf115(t, child), // every component shortened
		filepath.Join(dirShort, filepath.Base(child)), // the runner's mixed form: short ancestors, long leaf
		strings.ToUpper(child),                        // a case-only variation
	}
	for _, s := range tr.spelling {
		if s == "" {
			t.Fatalf("planted an empty spelling for %s", child)
		}
	}
	if strings.EqualFold(tr.spelling[1], child) {
		t.Fatalf("the instrument planted nothing: the artifact's short form %q is already its long form", tr.spelling[1])
	}
	// The mixed form must be a genuine alias too, and it is the one the runner
	// actually drives: a shortened profile directory under a long test directory.
	if strings.EqualFold(tr.spelling[2], child) {
		t.Fatalf("the instrument planted nothing: the mixed spelling %q equals the long spelling", tr.spelling[2])
	}
	return tr
}

// TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree is AC#3's forward half: one
// folder, several legal spellings (long name / 8.3 short name / mixed / case-only),
// and the notice's attribution must not depend on which one the caller typed.
//
// It also reproduces the four CI reds' mechanism on this box (AC#1): the notice
// carries the resolver's answer, which is fold-equal to this box's long spelling
// and never to the short one, so the attribution leg asks the tree rule (a string
// comparison there is what lost the runner two cases over one spelling) while the
// instrument leg stays a literal comparison, because what it has to prove is that
// the answer is not the spelling that was typed.
func TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree(t *testing.T) {
	got := captureNotices115(t)
	tr := plantChildTree115(t, "wisp115longdirname")

	for i, typed := range tr.spelling {
		widen115(t, tr.child)
		*got = nil
		if err := SealFile(typed); err != nil {
			t.Fatalf("seal #%d through spelling %q: %v", i, typed, err)
		}
		// Nothing was suppressed and nothing lost its path (AC#5's forbidden
		// shapes): one seal of one widened object is one notice, carrying a path.
		if len(*got) != 1 {
			t.Fatalf("seal #%d through %q reported %d notice(s), want exactly 1: %+v", i, typed, len(*got), *got)
		}
		n := (*got)[0]
		if n.Path == "" {
			t.Fatalf("seal #%d through %q emitted a notice with no path: %+v", i, typed, n)
		}
		if len(n.Principals) == 0 {
			t.Fatalf("seal #%d through %q cleared nothing it was shown: %+v", i, typed, n)
		}
		for _, asked := range tr.spelling {
			hits := noticesAboutTree(*got, asked)
			if len(hits) != 1 {
				t.Fatalf("AC#3: attribution lost - notice from seal #%d (typed %q) is not attributable through spelling %q: %d hit(s), notice=%+v",
					i, typed, asked, len(hits), n)
			}
		}
		// The mechanism of the four reds, read off the same object: the short
		// spelling is what breaks a string comparison, and only that.
		if strings.EqualFold(n.Path, tr.spelling[1]) {
			t.Fatalf("the instrument measured nothing: the notice repeated the caller's short spelling (%q), "+
				"so no 8.3 expansion happened on this path and the CI reds are not reproduced here", n.Path)
		}
		if !answerNamesTree115(t, n, tr.child) {
			t.Errorf("the notice no longer carries the resolver's answer for the object that was sealed: path=%q, long spelling=%q",
				n.Path, tr.child)
		}
		t.Logf("seal #%d typed %q -> notice path %q (EqualFold with the short spelling: %v)",
			i, typed, n.Path, strings.EqualFold(n.Path, tr.spelling[1]))
	}
}

// TestNoticeAttributionKeepsTwoTreesApart is AC#3's back half, and the reason the
// forward half is not simply "match everything": the tree rule must stay narrower
// than a spelling and wider than an alias, so two objects are never one object and
// a directory is never one of its children. Ticket 112's guard has the same two
// halves for the same reason, and widening this one is what AC#5 forbids.
func TestNoticeAttributionKeepsTwoTreesApart(t *testing.T) {
	got := captureNotices115(t)
	a := plantChildTree115(t, "wisp115treeAAAA")
	b := plantChildTree115(t, "wisp115treeBBBB")

	widen115(t, a.child)
	widen115(t, b.child)
	*got = nil
	if err := SealFile(a.spelling[1]); err != nil { // A's child, typed short
		t.Fatalf("sealing A's artifact through its short spelling: %v", err)
	}
	if err := SealFile(b.child); err != nil { // B's child, typed long
		t.Fatalf("sealing B's artifact: %v", err)
	}
	if len(*got) != 2 {
		t.Fatalf("two widened objects sealed, want 2 notices, got %d: %+v", len(*got), *got)
	}

	for _, tc := range []struct {
		name  string
		tree  childTree115
		other childTree115
	}{
		{"A", a, b},
		{"B", b, a},
	} {
		for _, typed := range tc.tree.spelling {
			hits := noticesAboutTree(*got, typed)
			if len(hits) != 1 {
				t.Fatalf("AC#3 back half: %d notice(s) attributed to %s through %q, want the 1 that names it: %+v",
					len(hits), tc.name, typed, *got)
			}
			if !answerNamesTree115(t, hits[0], tc.tree.child) {
				t.Fatalf("AC#3 back half: the notice attributed to %q is not about %s: path=%q", typed, tc.name, hits[0].Path)
			}
		}
		for _, foreign := range tc.other.spelling {
			if hits := noticesAboutTree(*got, foreign); len(hits) != 1 {
				// Both trees are in the captured set, so a foreign spelling must
				// still find exactly its own notice - and never the other one.
				t.Fatalf("spelling %q of the other tree matched %d notice(s) of this pair, want its own 1: %+v",
					foreign, len(hits), *got)
			} else if answerNamesTree115(t, hits[0], tc.tree.child) {
				t.Fatalf("AC#3 back half: %q (the other tree) was attributed to this tree's notice: %+v", foreign, hits[0])
			}
		}
	}

	// A directory is not its child: the containment shape of comparison (any
	// prefix or Contains test, which is what a lazy fix reaches for) would read
	// the child's notice as "about the directory too".
	dirSpellings := []string{a.dir, a.dirShort, strings.ToUpper(a.dir)}
	for _, typed := range dirSpellings {
		if hits := noticesAboutTree(*got, typed); len(hits) != 0 {
			t.Errorf("AC#3 back half: the child's notice was attributed to its own parent %q: %+v", typed, hits)
		}
	}
	// And the same objects, once the directory itself is widened and sealed: now
	// there is one notice per object, and each spelling still finds exactly one.
	mustExec(t, "icacls", a.dir, "/grant", "*"+everyoneSID+":(RX)")
	if err := SealDir(a.dirShort); err != nil { // the directory, typed short
		t.Fatalf("sealing A's directory through its short spelling: %v", err)
	}
	for _, typed := range dirSpellings {
		hits := noticesAboutTree(*got, typed)
		if len(hits) != 1 {
			t.Fatalf("the directory seal's own notice is not attributable through %q: %d hit(s), notices=%+v", typed, len(hits), *got)
		}
		if !answerNamesTree115(t, hits[0], a.dir) {
			t.Fatalf("a notice about %q was attributed to the directory %q: %+v", hits[0].Path, typed, hits[0])
		}
	}
	for _, typed := range a.spelling {
		if hits := noticesAboutTree(*got, typed); len(hits) != 1 || !answerNamesTree115(t, hits[0], a.child) {
			t.Fatalf("the child's notice must stay its own: spelling %q matched %+v", typed, hits)
		}
	}
}
