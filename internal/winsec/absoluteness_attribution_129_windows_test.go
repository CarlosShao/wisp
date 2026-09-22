//go:build windows

package winsec

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Ticket 129's cases. The defect: sameTree and answerInsideTree compare two
// spellings after dropping, per component, everything filepath.VolumeName
// claims - and the absoluteness of a spelling lives in exactly that prefix on
// Windows. "C:wisp\p" (drive-relative, i.e. relative to whatever directory the
// process happens to be standing in *on C:*) and "C:\wisp\p" therefore compare
// equal although they name two objects. Ticket 126 put the volume segment into
// the rule; this is the segment of the same prefix that decides *which root the
// tail hangs off*, and it predates 126 and is untouched by it: sameVolume gives
// both spellings "C:", so 126's new leg reads them as one volume before and
// after alike.
//
// Three boundaries this file is written inside, from ticket 126's acceptance 7:
//
//	(1) no leg here claims 126 introduced or worsened anything - legs below are
//	    annotated with the verdict the pre-126 tree gives too, by construction
//	    (sameVolume("C:x","C:\\x") is true in both trees);
//	(2) the harm claim is measured in absoluteness_seam_landing_129_windows_test.go,
//	    not assumed here; this file judges the comparisons and the seam verdict;
//	(3) everything here is only allowed to get stricter, and none of it moves
//	    into pathComponents (ticket 108's TestAC2ComponentsAndTraversalPerSeparatorShape
//	    is the nail that pins that, and 126's acceptance measured it going red).
//
// The fixture fake is ticket 126's crossVolumeWitness126, reused rather than
// copied: it answers only the spellings the guard asks and records what it was
// asked, which is what turns "the guard never reached this leg" into a red
// rather than a quiet green.

const (
	// tail129 is byte-identical below the volume prefix in both spellings, so
	// pathComponents cannot tell the pair apart - which is the defect.
	tail129 = `wisp129-trees\store-44440\artifact.txt`
	// abs129 is the absolute spelling: tail hangs off the volume root.
	abs129 = `C:\` + tail129
	// rel129 is the drive-relative spelling: tail hangs off the process's own
	// current directory on C:.
	rel129 = `C:` + tail129
)

// TestAbsolutenessSpellingsAreNotOneTree is the unit face: the two comparisons
// themselves, over spellings this package has to decide about without asking the
// filesystem. Like ticket 126's file it carries the CONTROL legs that keep a nail
// this narrow from rotting into "nothing is ever one tree".
func TestAbsolutenessSpellingsAreNotOneTree(t *testing.T) {
	// Instrument sanity first: if the two spellings ever stop differing in the
	// one property this case is about, every leg below measures nothing.
	if filepath.IsAbs(rel129) {
		t.Fatalf("the instrument planted nothing: filepath.IsAbs(%q) = true, so the pair is not a drive-relative vs absolute pair", rel129)
	}
	if !filepath.IsAbs(abs129) {
		t.Fatalf("the instrument planted nothing: filepath.IsAbs(%q) = false, so the pair is not a drive-relative vs absolute pair", abs129)
	}
	if v1, v2 := filepath.VolumeName(rel129), filepath.VolumeName(abs129); v1 != v2 {
		t.Fatalf("the instrument planted a cross-volume pair (%q vs %q), which is ticket 126's shape and not this one", v1, v2)
	}
	// The same claim from the other direction: this package's own floor already
	// refuses the drive-relative spelling, so the judgement it needs is already
	// written down in this package - the comparison leg just never consults it.
	if _, err := (builtinVerifier{}).Resolve(rel129); err == nil || !strings.Contains(err.Error(), "is not absolute") {
		t.Errorf("AC#2 RED: builtinVerifier.Resolve(%q) = %v, but the floor refusing a drive-relative answer is the whole premise of this ticket", rel129, err)
	}
	if _, err := (builtinVerifier{}).Resolve(abs129); err != nil {
		t.Errorf("CONTROL RED: builtinVerifier refused the absolute spelling %q: %v", abs129, err)
	}

	// The two directions of the same pair, on both comparisons.
	if sameTree(rel129, abs129) {
		t.Errorf("AC#1/AC#3 RED: sameTree(%q, %q) = true, but a drive-relative spelling hangs its tail off the process's current directory on C: and an absolute one off C:\\ - two objects, one tree in this comparison's eyes", rel129, abs129)
	}
	if sameTree(abs129, rel129) {
		t.Errorf("AC#1/AC#3 RED: sameTree(%q, %q) = true (mirror order of the leg above)", abs129, rel129)
	}
	if answerInsideTree(rel129, `C:\wisp129-trees\store-44440`) {
		t.Errorf("AC#1/AC#3 RED: answerInsideTree(%q, %q) = true, so a witness naming the process's own C: directory would vouch for a seal on the tree at C:\\", rel129, `C:\wisp129-trees\store-44440`)
	}
	if answerInsideTree(abs129, `C:wisp129-trees\store-44440`) {
		t.Errorf("AC#1/AC#3 RED: answerInsideTree(%q, %q) = true, the mirror direction: an absolute answer admitted as sitting inside a drive-relative parent", abs129, `C:wisp129-trees\store-44440`)
	}
	// Unequal component counts, so neither leg can be carried by the length check
	// alone (ticket 126's same shape, one segment over).
	if answerInsideTree(`C:wisp129-trees\store-44440\artifact.txt\deeper`, `C:\wisp129-trees\store-44440`) {
		t.Errorf("AC#1/AC#3 RED: answerInsideTree admitted a deeper drive-relative path as being inside an absolute tree")
	}

	// CONTROL legs, both directions of the rule. Two spellings of one kind that
	// differ only in case or separator still name one tree: this is what keeps
	// the fix a pair rule rather than "refuse everything non-absolute", which
	// would be a different (and unjustified) verdict.
	if !sameTree(rel129, rel129) {
		t.Errorf("CONTROL RED: sameTree(%q, itself) = false, so the new leg is refusing rather than comparing", rel129)
	}
	if !sameTree(`c:`+tail129, rel129) {
		t.Errorf("CONTROL RED: sameTree denied a mixed-case volume letter on two drive-relative spellings of one shape")
	}
	if !answerInsideTree(`C:wisp129-trees\store-44440\artifact.txt`, `C:wisp129-trees\store-44440`) {
		t.Errorf("CONTROL RED: answerInsideTree denied a drive-relative path inside its own drive-relative tree")
	}
	if !sameTree(abs129, `C:/wisp129-trees/store-44440/artifact.txt`) {
		t.Errorf("CONTROL RED: sameTree denied two separator spellings of one absolute path, a rule this change may not alter")
	}
}

// TestAttributionFaceNeverSeesAMixedAbsolutenessPair is the reading behind this
// ticket's AC#2 claim that the attribution face cannot be reached by this defect
// today. noticeNamesTree resolves the caller's spelling through ResolvePath and
// compares it against the notice's own path, so both sides of that comparison are
// answers this package already vouched for. If every one of them is absolute, the
// new leg can only ever agree there, and the face it protects is the guard's, not
// the audit channel's. Measured on a real seal rather than argued from the source.
func TestAttributionFaceNeverSeesAMixedAbsolutenessPair(t *testing.T) {
	got := captureNotices115(t)
	dir := filepath.Join(t.TempDir(), "store-44440")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(dir, "artifact.txt")
	if err := os.WriteFile(child, []byte("wisp 129\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	widen115(t, child)
	if err := SealFile(child); err != nil {
		t.Fatalf("SealFile(%s): %v", child, err)
	}
	if len(*got) == 0 {
		t.Fatal("this instrument measured nothing: the seal emitted no notice to inspect")
	}
	resolved, err := ResolvePath(child)
	if err != nil {
		t.Fatalf("ResolvePath(%q) refused the spelling that was just sealed: %v", child, err)
	}
	if !filepath.IsAbs(resolved.String()) {
		t.Errorf("AC#2 RED: ResolvePath answered %q with %q, which is not absolute, so the attribution face can be handed a mixed pair after all",
			child, resolved.String())
	}
	for _, n := range *got {
		if !filepath.IsAbs(n.Path) {
			t.Errorf("AC#2 RED: a notice carries %q, which is not absolute, so sameTree can be asked about a pair that disagrees in absoluteness here", n.Path)
		}
		if !noticeNamesTree(n, child) {
			t.Errorf("CONTROL RED: the notice for %q is no longer attributed to the tree that was just sealed", n.Path)
		}
	}
	t.Logf("AC#2 attribution-face reading: %d notice(s), all absolute=%v, caller's resolved answer absolute=%v (%s)",
		len(*got), allAbsolute129(*got), filepath.IsAbs(resolved.String()), resolved.String())
}

func allAbsolute129(notices []narrowNotice) bool {
	for _, n := range notices {
		if !filepath.IsAbs(n.Path) {
			return false
		}
	}
	return true
}

// absolutenessPair129 names the seam guard's probe shapes: two candidate answers
// that name two objects while agreeing in every segment pathComponents looks at.
const (
	probeTreeAbs129 = `D:\wisp129-seam\probe-tree`
	probeTreeRel129 = `D:wisp129-seam\probe-tree`
	movedAbs129     = `D:\wisp129-seam\moved-seal\leaf`
	movedDir129     = `D:\wisp129-seam\moved-seal`
	otherTree129    = `D:\wisp129-seam\unrelated-tree`
)

// TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness is
// AC#1's verdict face and the same measurement ticket 126's acceptance recorded
// as 7 for this hole: the pair goes into treeOwnershipFailureForPair and either
// comes out with refusal "" (the seam would install the candidate) or with a
// refusal. Refusal "" on a leg that wants one is a gate being walked through,
// not a bookkeeping entry: both comparisons in this function are read only as
// grounds to PASS a candidate (resolve.go:391 and :405).
//
// Each refusal leg is paired with a control that differs from it ONLY in the
// absoluteness of one spelling, so the verdict cannot come from the component
// rule, the equal-answers leg or a refusal - and so that the fix cannot be
// scored by a leg that is green only because this package now refuses every
// non-absolute answer.
func TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness(t *testing.T) {
	sep := string(filepath.Separator)
	parent := os.TempDir()
	child := parent + sep + "wisp-129-tree-ownership-probe"
	// The instrument's own shapes must not collapse: if the child's answer were
	// inside the tree named for the probe parent by the component rule alone, the
	// guard would stop at its first containment leg and the witness legs below
	// would never be asked.
	if answerInsideTree(movedAbs129, probeTreeAbs129) {
		t.Fatal("the instrument planted nothing: the child's answer is inside the probe parent's tree")
	}
	legs := []struct {
		name        string
		answers     map[string]string
		wantRefused bool
		// anchor is the spelling whose being asked about proves which of the
		// guard's two comparisons carries this leg.
		anchor        string
		wantAnchorAsk bool
		why           string
	}{
		{
			// 126's leg 1 shape with the volume difference replaced by an
			// absoluteness difference: the second witness names the probe parent's
			// tree only in a spelling that hangs off the process's C: directory.
			name:        "second witness vouches for the probe parent's tree in a drive-relative spelling",
			answers:     map[string]string{parent: probeTreeAbs129, child: movedAbs129, movedDir129: probeTreeRel129},
			wantRefused: true, anchor: movedDir129, wantAnchorAsk: true,
			why: "sameTree read a drive-relative answer as the absolute tree it names only after absolutization",
		},
		{
			// The first containment leg, :391, has the same blind spot: a child
			// answer inside a *drive-relative* parent answer is not inside the
			// absolute tree the candidate named for the probe parent.
			name:        "the child's own answer sits in the probe parent's tree only when read drive-relative",
			answers:     map[string]string{parent: probeTreeAbs129, child: probeTreeRel129 + `\leaf`, probeTreeRel129: otherTree129},
			wantRefused: true, anchor: probeTreeRel129, wantAnchorAsk: true,
			why: "answerInsideTree at :391 took a drive-relative child as inside an absolute parent",
		},
		{
			// Isolates the gap: an absolute second witness naming another tree is
			// refused today, so the admission above is not "the guard never looks".
			name:        "control for the gap: an absolute second witness naming another tree",
			answers:     map[string]string{parent: probeTreeAbs129, child: movedAbs129, movedDir129: otherTree129},
			wantRefused: true, anchor: movedDir129, wantAnchorAsk: true,
			why: "this leg is green before this ticket too; it is here to show what a refusal looks like",
		},
		{
			name:        "CONTROL: the second witness names that very tree, both spellings absolute",
			answers:     map[string]string{parent: probeTreeAbs129, child: movedAbs129, movedDir129: probeTreeAbs129},
			wantRefused: false, anchor: movedDir129, wantAnchorAsk: true,
			why: "the honest self-vouching ticket 112 and 126 already protect",
		},
		{
			// The anti-tautology leg: the fix is a pair rule. Two drive-relative
			// answers naming one shape are still one tree, so "refuse everything
			// non-absolute" is not what this nails down.
			name:        "CONTROL: both the probe parent's answer and the second witness are drive-relative and identical",
			answers:     map[string]string{parent: probeTreeRel129, child: movedAbs129, movedDir129: probeTreeRel129},
			wantRefused: false, anchor: movedDir129, wantAnchorAsk: true,
			why: "agreeing absoluteness is what the rule requires, not absolute-ness itself",
		},
		{
			name:        "CONTROL: the child's answer is honestly inside the probe parent's tree",
			answers:     map[string]string{parent: probeTreeAbs129, child: probeTreeAbs129 + `\leaf`},
			wantRefused: false, anchor: movedDir129, wantAnchorAsk: false,
			why: "the containment leg must still stop there, exactly as before the change",
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
		t.Logf("AC#1 leg %q: answers=%v asked=%v refusal=%q (%s)", leg.name, leg.answers, asked, reason, leg.why)
		if !fake.asked[parent] || !fake.asked[child] {
			t.Errorf("AC#1 instrument measured nothing on leg %q: the guard never asked about both probe shapes (asked %v)", leg.name, asked)
			continue
		}
		if fake.asked[leg.anchor] != leg.wantAnchorAsk {
			t.Errorf("AC#1 instrument measured the wrong leg %q: witness %q asked=%v, wanted %v",
				leg.name, leg.anchor, fake.asked[leg.anchor], leg.wantAnchorAsk)
		}
		if refused != leg.wantRefused {
			if leg.wantRefused {
				t.Errorf("AC#1/AC#3 RED on leg %q: the seam guard admitted a candidate whose answers name two objects (%s); it owed a refusal and said %q. answers=%v",
					leg.name, leg.why, reason, leg.answers)
			} else {
				t.Errorf("CONTROL RED on leg %q: the seam guard now refuses a candidate whose answers name one tree: %s. A guard that over-refuses takes the whole sealing pipeline down to the built-in floor, which is the failure run 35595651898 measured.",
					leg.name, reason)
			}
		}
	}
}

// TestHonestPipelineKeepsPassingItsOwnProbe is the cost side of the rule, on the
// pipeline that is really installed rather than on a scripted candidate. Ticket
// 112's failure (run 35595651898) was the guard reading two spellings of ONE
// object as two trees and dropping C26 to the floor for a whole machine, so a
// change to this comparison has to show that the honest shapes still pass - and
// that the pair it now refuses is a pair no honest answer can produce.
func TestHonestPipelineKeepsPassingItsOwnProbe(t *testing.T) {
	// Leg 1: the production tree-ownership probe, run against the built-in floor.
	if reason := resolverTreeOwnershipFailure(builtinVerifier{}); reason != "" {
		t.Errorf("CONTROL RED: the built-in floor answers the probe parent and child with the caller's own absolute spellings, one inside the other, and is now refused by its own tree-ownership probe: %s", reason)
	}
	// Leg 2: the same probe on the machine's own 8.3 asymmetry, planted (ticket
	// 112's rule: this box's temp root needs no short name, so living here does
	// not reproduce the runner). Both spellings are absolute, so this ticket's
	// rule has no say here and the verdict must not move.
	t.Run("on a planted 8.3 short spelling of an existing parent", func(t *testing.T) {
		fixture := filepath.Join(t.TempDir(), "wisp129-long-profile-name")
		if err := os.MkdirAll(fixture, 0o700); err != nil {
			t.Fatal(err)
		}
		short := shortFormOf115(t, fixture)
		if short == fixture {
			t.Skipf("this volume gives %q no short name (GetShortPathName answered %q), so the runner's RUNNER~1 shape cannot be planted here", fixture, short)
		}
		sep := string(filepath.Separator)
		if reason := treeOwnershipFailureForPair(builtinVerifier{}, short, short+sep+"wisp-129-honest-probe"); reason != "" {
			t.Errorf("CONTROL RED: an honest answer pair for the 8.3 spelling %q is now refused: %s", short, reason)
		} else {
			t.Logf("CONTROL: honest 8.3 pair still passes; long=%q short=%q both absolute=%v/%v",
				fixture, short, filepath.IsAbs(fixture), filepath.IsAbs(short))
		}
	})
	// Leg 3 is the cost reading, and it belongs to the external file, where the
	// pipeline that is actually installed can be asked through the exported
	// ResolvePath: absoluteness_seam_landing_129_windows_test.go's
	// TestAC2NoSealEverActsOnAGloballyDifferentAbsoluteness measures that every
	// spelling ResolvePath answers without an error answers absolute, which is
	// what makes this leg refuse no seal rather than merely no refusal.
}
