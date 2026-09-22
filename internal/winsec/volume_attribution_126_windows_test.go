//go:build windows

package winsec

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Ticket 126's cases. The defect: pathComponents starts after
// filepath.VolumeName, so every tree comparison in this package compared tails
// only, and C:\store-44440\artifact.txt was one tree with
// D:\store-44440\artifact.txt. Two production faces read that verdict and both
// are pinned here: noticeNamesTree (which notice is about which tree) and the
// install-time tree-ownership leg of the sealing seam guard, which decides
// whether a candidate resolver may be installed at all.
//
// Three of the four legs below are expressible with planted volume letters and
// so run on a single-volume machine, which is what makes this defect class
// visible on the CI runner at all (AC#4). The fourth needs two REAL volumes,
// because that is the only shape in which both trees exist and neither one is a
// lie; it is a named subtest that prints its own skip reason, ticket 112's rule
// for a capability gate, since a top-level SKIP is fatal to the winsec gate and
// quietly dropping the leg is worse than naming the machine that cannot build
// it. The fixtures are ticket 115's (captureNotices115, widen115), not copies:
// this file asks the same question that file settled how to ask.

// crossVolumePair126 returns an existing object's resolved answer and a second
// spelling of a byte-identical tail on a volume letter this machine does not
// have. The tail is read off the answer rather than hand-written, so the pair is
// derived from what the resolver actually says.
func crossVolumePair126(t *testing.T, existing string) (real, planted string) {
	t.Helper()
	resolved, err := ResolvePath(existing)
	if err != nil {
		t.Fatalf("this instrument cannot be asked at all: ResolvePath(%q) refused to vouch for it: %v", existing, err)
	}
	answer := resolved.String()
	tail := answer[len(filepath.VolumeName(answer)):]
	if strings.Trim(tail, string(filepath.Separator)) == "" {
		t.Fatalf("the planted object %q resolved to a volume root %q, so it has no tail to copy", existing, answer)
	}
	for letter := 'Z'; letter >= 'A'; letter-- {
		root := string(letter) + `:\`
		if _, statErr := os.Lstat(root); statErr == nil {
			continue // a volume this machine has; the real-volume case covers those
		}
		candidate := string(letter) + `:` + tail
		if strings.EqualFold(filepath.VolumeName(candidate), filepath.VolumeName(answer)) {
			continue
		}
		return answer, candidate
	}
	t.Fatalf("no free volume letter was found to plant a second volume on")
	return "", ""
}

// TestCrossVolumeSpellingsAreNotOneTree is the unit face: the two comparisons
// themselves, over spellings this package must decide about without asking the
// filesystem. It carries the legs that keep the fix from becoming "nothing is
// ever one tree", which is the direction a nail this narrow rots in.
func TestCrossVolumeSpellingsAreNotOneTree(t *testing.T) {
	const tail = `wisp126-trees\store-44440\artifact.txt`
	a := `C:\` + tail
	b := `D:\` + tail
	if sameTree(a, b) {
		t.Errorf("AC#3 RED: sameTree(%q, %q) = true, but these are two trees on two volumes; the comparison dropped the volume segment", a, b)
	}
	if answerInsideTree(b, a) {
		t.Errorf("AC#3 RED: answerInsideTree(%q, %q) = true, so a witness naming %q would vouch for a seal on the tree %q", b, a, b, a)
	}
	// The same pair with the tails unequal in length, so neither leg can be
	// satisfied by the component-count check alone.
	if answerInsideTree(`D:\wisp126-trees\store-44440\artifact.txt\deeper`, `C:\wisp126-trees\store-44440`) {
		t.Errorf("AC#3 RED: answerInsideTree admitted a deeper path on D: as being inside a tree on C:")
	}
	// Two UNC shares are the same defect one layer over: pathComponents strips
	// \\server\share as the volume, so different shares with one tail were also
	// one tree. platformVerifyPlacement refuses a UNC spelling before anything
	// can be sealed through it, so this leg pins the rule rather than a route a
	// seal takes - and it goes red under the mutation that drops the volume
	// comparison, which is what makes it a rule and not a decoration.
	if sameTree(`\\fileserver-1\data\wisp126\artifact.txt`, `\\fileserver-2\data\wisp126\artifact.txt`) {
		t.Errorf("AC#3 RED: sameTree admitted two different UNC shares as one tree")
	}
	if !sameTree(`\\fileserver-1\data\wisp126\artifact.txt`, `\\fileserver-1\data\wisp126\artifact.txt`) {
		t.Errorf("CONTROL RED: sameTree denied two identical UNC spellings, so the new check is refusing rather than comparing")
	}

	// Forward legs. One volume is spelled in either case, and an object is on its
	// volume however it is named: these are what keep sameVolume from being
	// "return false" in disguise.
	if !sameTree(a, a) {
		t.Errorf("CONTROL RED: sameTree(%q, itself) = false", a)
	}
	if !sameTree(`c:\`+tail, `C:\`+tail) {
		t.Errorf("CONTROL RED: sameTree denied a mixed-case spelling of the same volume: %q vs %q", `c:\`+tail, `C:\`+tail)
	}
	if !answerInsideTree(`c:\wisp126-trees\store-44440\artifact.txt`, `C:\wisp126-trees\store-44440`) {
		t.Errorf("CONTROL RED: answerInsideTree denied a mixed-case volume for a path that is inside the same tree")
	}
	// And the pre-existing component rule survives untouched: a forward slash is
	// still the platform's own separator, and that is not a volume question.
	if !sameTree(`C:\wisp126-trees\store-44440\artifact.txt`, `C:/wisp126-trees/store-44440/artifact.txt`) {
		t.Errorf("CONTROL RED: sameTree denied two separator spellings of one path, a rule this change may not alter")
	}
}

// TestNoticeFromOneVolumeIsNotAttributedToAnotherVolumeSpelling is the
// attribution face driven through the production entry point: seal a real object,
// then ask whether its notice is about a spelling that differs only in volume.
// The other spelling is planted, so this leg runs on a single-volume machine.
func TestNoticeFromOneVolumeIsNotAttributedToAnotherVolumeSpelling(t *testing.T) {
	got := captureNotices115(t)
	dir := filepath.Join(t.TempDir(), "store-44440")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(dir, "artifact.txt")
	if err := os.WriteFile(child, []byte("wisp 126\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	widen115(t, child)
	if err := SealFile(child); err != nil {
		t.Fatalf("SealFile(%s): %v", child, err)
	}
	if len(*got) == 0 {
		t.Fatal("this instrument measured nothing: the seal emitted no notice to attribute")
	}
	sealed, planted := crossVolumePair126(t, child)
	sealedTail := sealed[len(filepath.VolumeName(sealed)):]
	plantedTail := planted[len(filepath.VolumeName(planted)):]
	if sealedTail != plantedTail {
		t.Fatalf("the instrument planted nothing: the two tails differ, %q vs %q", sealedTail, plantedTail)
	}
	if strings.EqualFold(filepath.VolumeName(sealed), filepath.VolumeName(planted)) {
		t.Fatalf("the instrument planted nothing: both spellings are on volume %q", filepath.VolumeName(sealed))
	}
	t.Logf("AC#4 shapes: sealed=%s planted=%s (volumes %q vs %q, tails byte-identical)",
		sealed, planted, filepath.VolumeName(sealed), filepath.VolumeName(planted))
	mine := noticesAboutTree(*got, child)
	theirs := noticesAboutTree(*got, planted)
	t.Logf("notices attributed to the tree that was sealed: %d of %d; to the planted other-volume spelling: %d", len(mine), len(*got), len(theirs))
	if len(mine) != len(*got) {
		t.Errorf("AC#3 RED: %d of %d notices from a seal on %q are no longer attributed to that tree, so the comparison now refuses instead of distinguishing",
			len(*got)-len(mine), len(*got), sealed)
	}
	if len(theirs) != 0 {
		t.Errorf("AC#3 RED: %d notice(s) from a seal that ran on %q are attributed to %q, a spelling on the volume %q that was never sealed. pathComponents dropped the volume segment and the two trees became one.",
			len(theirs), sealed, planted, filepath.VolumeName(planted))
	}
	for _, n := range *got {
		if noticeNamesTree(n, planted) {
			t.Errorf("AC#3 RED: noticeNamesTree(notice for %q, %q) = true", n.Path, planted)
		}
	}
}

// TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume is the shape only a
// real second volume can express: both trees exist, both carry the same
// reportable grant, and exactly one of them is sealed.
func TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume(t *testing.T) {
	t.Run("two real volumes", func(t *testing.T) {
		roots, refused := writableVolumeRoots126(t)
		if len(roots) < 2 {
			// Named, not quiet: the denominator and the reason both print, and
			// AC#4 registers that a single-volume runner takes this branch.
			t.Skipf("this machine has %d volume root(s) that accept a directory (%v; %d refused), so two real trees with one tail cannot be planted",
				len(roots), roots, len(refused))
		}
		fixture := fmt.Sprintf("wisp126-xvol-%d", os.Getpid())
		aDir := filepath.Join(roots[0], fixture, "store-44440")
		bDir := filepath.Join(roots[1], fixture, "store-44440")
		a := filepath.Join(aDir, "artifact.txt")
		b := filepath.Join(bDir, "artifact.txt")
		t.Cleanup(func() {
			for _, root := range []string{roots[0], roots[1]} {
				target := filepath.Join(root, fixture)
				if err := os.RemoveAll(target); err != nil {
					t.Errorf("AC#6 cleanup RED: RemoveAll(%s): %v", target, err)
					continue
				}
				if _, err := os.Lstat(target); !os.IsNotExist(err) {
					t.Errorf("AC#6 cleanup RED: %s still names something after RemoveAll (err=%v)", target, err)
				}
			}
		})
		for _, d := range []string{aDir, bDir} {
			if err := os.MkdirAll(d, 0o700); err != nil {
				t.Fatalf("MkdirAll(%s): %v", d, err)
			}
		}
		for _, f := range []string{a, b} {
			if err := os.WriteFile(f, []byte("wisp 126\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			// Both trees carry the same out-of-band grant, so "0 notices about B"
			// cannot be satisfied by "B has nothing to report".
			widen115(t, f)
		}
		got := captureNotices115(t)
		if err := SealFile(a); err != nil {
			t.Fatalf("SealFile(%s): %v", a, err)
		}
		mine, theirs := noticesAboutTree(*got, a), noticesAboutTree(*got, b)
		t.Logf("AC#8 shapes: A=%s B=%s (volumes %q vs %q)", a, b, filepath.VolumeName(a), filepath.VolumeName(b))
		t.Logf("notices after sealing A only: %d; attributed to A: %d; attributed to never-sealed B: %d", len(*got), len(mine), len(theirs))
		if len(*got) == 0 || len(mine) != len(*got) {
			t.Errorf("AC#3 RED: the tree that was sealed gets %d of its own %d notices", len(mine), len(*got))
		}
		if len(theirs) != 0 {
			t.Errorf("AC#3 RED: %d notice(s) from a seal that ran on %q are attributed to %q, a tree on the volume %q that was never sealed",
				len(theirs), a, b, filepath.VolumeName(b))
		}
	})
}

// crossVolumeWitness126 is the candidate ticket 118's acceptance recorded as
// never built (R-118-9): a resolver whose answers at the seam guard name a tree
// on ANOTHER volume whose tail is byte-identical to a tree it named elsewhere.
// It answers only the spellings the guard actually asks about, and records what
// it was asked, so "the guard never reached the leg this case is about" is a
// Fatal rather than a quiet green.
type crossVolumeWitness126 struct {
	answers map[string]string
	asked   map[string]bool
}

func (f *crossVolumeWitness126) Resolve(input string) (string, error) {
	p, _, err := f.ResolveAccounted(input)
	return p, err
}

func (f *crossVolumeWitness126) ResolveAccounted(input string) (string, bool, error) {
	if f.asked == nil {
		f.asked = map[string]bool{}
	}
	f.asked[input] = true
	answer, ok := f.answers[input]
	if !ok {
		return "", false, fmt.Errorf("%w: instrument was not expecting %q", ErrUnresolvedPath, input)
	}
	return answer, false, nil
}

// TestSeamGuardRefusesACandidateWhoseSecondWitnessNamesAnotherVolume is AC#2's
// leg, now nailed. It needs no second volume: the candidate lies about paths and
// the guard compares only what it is handed, which is exactly why this leg runs
// on a single-volume runner.
//
// Each refusal leg is paired with a control that differs from it ONLY in the
// volume segment, so the verdict cannot come from the component rule, the
// equal-answers leg or a refusal. Which of the guard's three comparisons carries
// each leg is attributed by the mutation matrix on the ticket, not by reading
// prose off a log line.
func TestSeamGuardRefusesACandidateWhoseSecondWitnessNamesAnotherVolume(t *testing.T) {
	sep := string(filepath.Separator)
	parent := os.TempDir()
	child := parent + sep + "wisp-126-tree-ownership-probe"
	const (
		parentOnD = `D:\wisp126-seam\probe-tree`
		movedOnD  = `D:\wisp126-seam\moved-seal\leaf`
		movedDir  = `D:\wisp126-seam\moved-seal`
	)
	// The instrument's own shapes must not collapse: if the child's answer were
	// inside the tree named for the probe parent, the guard would stop at its
	// first containment leg and the two witness legs below would never be asked.
	if answerInsideTree(movedOnD, parentOnD) {
		t.Fatal("the instrument planted nothing: the child's answer is inside the probe parent's tree")
	}
	const childOnOtherVolumeDir = `C:\wisp126-seam\probe-tree`
	legs := []struct {
		name        string
		answers     map[string]string
		wantRefused bool
		// anchor is the spelling whose being asked about proves the guard reached
		// the second witness at resolve.go:325 rather than stopping earlier.
		anchor        string
		wantAnchorAsk bool
	}{
		{
			name:        "second witness names the same tree on another volume",
			answers:     map[string]string{parent: parentOnD, child: movedOnD, movedDir: `C:\wisp126-seam\probe-tree`},
			wantRefused: true, anchor: movedDir, wantAnchorAsk: true,
		},
		{
			name:        "second witness names a deeper tree on another volume",
			answers:     map[string]string{parent: parentOnD, child: movedOnD, movedDir: `C:\wisp126-seam\probe-tree\deeper`},
			wantRefused: true, anchor: movedDir, wantAnchorAsk: true,
		},
		{
			name: "the child's own answer names a deeper tree on another volume",
			answers: map[string]string{
				parent:                parentOnD,
				child:                 childOnOtherVolumeDir + `\leaf`,
				childOnOtherVolumeDir: `C:\wisp126-seam\unrelated-tree`,
			},
			wantRefused: true, anchor: childOnOtherVolumeDir, wantAnchorAsk: true,
		},
		{
			name:        "control: the second witness names that very tree, same volume",
			answers:     map[string]string{parent: parentOnD, child: movedOnD, movedDir: parentOnD},
			wantRefused: false, anchor: movedDir, wantAnchorAsk: true,
		},
		{
			name:        "control: the second witness names a containing tree, same volume",
			answers:     map[string]string{parent: parentOnD, child: movedOnD, movedDir: parentOnD + `\deeper`},
			wantRefused: false, anchor: movedDir, wantAnchorAsk: true,
		},
		{
			name:        "control: the child's answer is inside the probe parent's tree",
			answers:     map[string]string{parent: parentOnD, child: parentOnD + `\leaf`},
			wantRefused: false, anchor: movedDir, wantAnchorAsk: false,
		},
	}
	for _, leg := range legs {
		fake := &crossVolumeWitness126{answers: leg.answers}
		reason := treeOwnershipFailureForPair(fake, parent, child)
		refused := reason != ""
		asked := make([]string, 0, len(fake.asked))
		for input := range fake.asked {
			asked = append(asked, input)
		}
		sort.Strings(asked)
		t.Logf("AC#2 leg %q: answers=%v asked=%v refusal=%q", leg.name, leg.answers, asked, reason)
		if !fake.asked[parent] || !fake.asked[child] {
			t.Errorf("AC#2 instrument measured nothing on leg %q: the guard never asked about both probe shapes (asked %v)", leg.name, asked)
			continue
		}
		if fake.asked[leg.anchor] != leg.wantAnchorAsk {
			t.Errorf("AC#2 instrument measured the wrong leg %q: second witness %q asked=%v, wanted %v",
				leg.name, leg.anchor, fake.asked[leg.anchor], leg.wantAnchorAsk)
		}
		if refused != leg.wantRefused {
			if leg.wantRefused {
				t.Errorf("AC#2/AC#3 RED on leg %q: the seam guard admitted a candidate that moves a seal across volumes (answers %v); it owed a refusal and said %q",
					leg.name, leg.answers, reason)
			} else {
				t.Errorf("CONTROL RED on leg %q: the seam guard now refuses a candidate whose answers name one tree on one volume: %s. A guard that over-refuses takes the whole sealing pipeline down to the built-in floor, which is the failure run 35595651898 measured.",
					leg.name, reason)
			}
		}
	}
}

// writableVolumeRoots126 enumerates volumes by creating a DIRECTORY at each root,
// which is the rule R-118-8 registers after two acceptance rounds reported
// different denominators for the same machine: this box's C: root refuses
// ordinary-user file creation and allows directory creation, so a
// WriteFile-based enumeration loses C: and never reports the pair the shape
// needs. Refusals are kept with their reason so the denominator is readable
// rather than merely counted, and every probe this leaves on a volume root is
// removed here rather than by the caller.
func writableVolumeRoots126(t *testing.T) ([]string, map[string]string) {
	t.Helper()
	var roots []string
	refused := map[string]string{}
	for letter := 'A'; letter <= 'Z'; letter++ {
		root := string(letter) + `:\`
		probe := root + `wisp126-xvol-enumeration-probe`
		if err := os.Mkdir(probe, 0o700); err != nil {
			refused[root] = err.Error()
			continue
		}
		roots = append(roots, root)
	}
	t.Cleanup(func() {
		for letter := 'A'; letter <= 'Z'; letter++ {
			probe := string(letter) + `:\wisp126-xvol-enumeration-probe`
			if _, err := os.Lstat(probe); os.IsNotExist(err) {
				continue
			}
			if err := os.RemoveAll(probe); err != nil {
				t.Errorf("AC#6 cleanup RED: the enumeration probe %s still exists: %v", probe, err)
			}
		}
	})
	return roots, refused
}
