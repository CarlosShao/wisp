//go:build windows

// Ticket 129 AC#1: the leg ticket 126's acceptance recorded as reasoning rather
// than reading ("过缝之后能把 seal 落到别人那棵树是推理不是读数"). This file
// measures it end to end on real objects: plant the two spellings for real, show
// the OS treats them as two objects, then let the install-time guard decide about
// a candidate whose second witness differs only in absoluteness, and read where
// the next real seal lands with icacls before and after.
//
// Nothing here re-judges ticket 126. The candidate pairs below are read the same
// way before and after 126's volume leg (sameVolume gives "C:" and "C:" for the
// two spellings), so every red in this file is on the absoluteness segment only.
//
// Layout of the files in this package: absoluteness_attribution_129_windows_test.go
// (internal) judges the comparisons and the guard's verdict; this file judges what
// the verdict costs, which needs the icacls fixtures that live on this side.
package winsec_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// seamProbePair129 is the probe pair the install-time guard actually asks, read
// off the guard's own questions rather than hand-written: pass one installs a
// recorder that answers everything with one fixed absolute path (ticket 108's
// treeMovingResolver108 shape, refused today by the equal-answers leg), and the
// inputs it was asked name the probe parent and child exactly as resolverProbeRoot
// spells them on this machine. That matters: a scripted candidate keyed on a
// guessed parent path would turn a miss into a resolver refusal, and the guard
// reads a refusal as "narrow" - which would look like a pass for the wrong reason.
func seamProbePair129(t *testing.T, fixed string) (parent, child string) {
	t.Helper()
	rec := &witness129{fixed: fixed}
	restore := winsec.SetSeamForTest(nil)
	defer restore()
	winsec.SetPathResolver(rec)
	if got := resolverName(winsec.PathResolverInstalled()); got != "<floor>" {
		t.Fatalf("the instrument's own pass measured nothing: the recorder got installed (seam is %s), so the guard did not refuse it", got)
	}
	var childAsked []string
	for input := range rec.asked {
		if strings.HasSuffix(input, `wisp-108-tree-ownership-probe`) {
			childAsked = append(childAsked, input)
		}
	}
	if len(childAsked) != 1 {
		t.Fatalf("the guard never asked the tree-ownership probe this ticket is about (%d candidate inputs ending in the 112 marker; asked %v)",
			len(childAsked), askedNames(rec))
	}
	child = childAsked[0]
	parent = strings.TrimSuffix(child, string(filepath.Separator)+"wisp-108-tree-ownership-probe")
	if parent == child || parent == "" {
		t.Fatalf("the instrument planted nothing: stripping the probe marker from %q gave %q", child, parent)
	}
	t.Logf("AC#1 probe pair read off the guard's own asks: parent=%q child=%q (absolute parent: %v)", parent, child, filepath.IsAbs(parent))
	return parent, child
}

func askedNames(rec *witness129) []string {
	out := make([]string, 0, len(rec.asked))
	for input := range rec.asked {
		out = append(out, input)
	}
	sort.Strings(out)
	return out
}

// witness129 is the candidate. It answers the three spellings the guard's
// tree-ownership leg asks about from a script, and everything else - which is to
// say any real caller's path - with one fixed absolute file in somebody else's
// tree. The scripted answers are what the seam judges; the fixed answer is the
// move it is supposed to catch, and it is the same answer in every leg below, so
// the only thing that can change a verdict is the comparison.
type witness129 struct {
	answers map[string]string
	fixed   string
	asked   map[string]bool
}

func (w *witness129) Resolve(input string) (string, error) {
	p, _, err := w.ResolveAccounted(input)
	return p, err
}

func (w *witness129) ResolveAccounted(input string) (string, bool, error) {
	if w.asked == nil {
		w.asked = map[string]bool{}
	}
	w.asked[input] = true
	if answer, ok := w.answers[input]; ok {
		return answer, false, nil
	}
	return w.fixed, false, nil
}

// ------------------------------------------------------------------ AC#1 a ----

// TestAC1DriveRelativeAndAbsoluteSpellingsNameTwoObjectsOnThisBox is the
// object-level half of the harm claim. Without it the ticket's premise - that the
// two spellings are two objects rather than two ways of typing one - is a claim
// about filepath.IsAbs, and filepath.IsAbs is not what a seal acts on.
func TestAC1DriveRelativeAndAbsoluteSpellingsNameTwoObjectsOnThisBox(t *testing.T) {
	t.Run("on the volume the process is standing on", func(t *testing.T) {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		vol := filepath.VolumeName(wd)
		if vol == "" {
			t.Skipf("os.Getwd() = %q carries no volume segment, so this box has no drive-relative shape to plant", wd)
		}
		fixture := fmt.Sprintf("wisp129-dr-%d", os.Getpid())
		tail := fixture + `\tree\object.txt`
		// The drive-relative spelling resolves under the process's own current
		// directory on this volume, so planting at that absolute path makes the
		// object reachable by the spelling under test - no normalization needed,
		// the two spellings are built from os.Getwd by construction.
		relDirSide := filepath.Join(wd, fixture, "tree")
		absSideRoot := driveAbsolute129(vol, fixture)
		rel := driveRelative129(vol, tail)
		abs := driveAbsolute129(vol, tail)
		// Registered before anything is planted: this file's first run skipped
		// between the two MkdirAll calls, and a skip that leaves a directory
		// behind is a residue bug ticket 118 and 126 each paid once for.
		t.Cleanup(func() {
			for _, target := range []string{filepath.Join(wd, fixture), absSideRoot} {
				if rmErr := os.RemoveAll(target); rmErr != nil {
					t.Errorf("AC#5 cleanup RED: RemoveAll(%s): %v", target, rmErr)
					continue
				}
				if _, lErr := os.Lstat(target); !os.IsNotExist(lErr) {
					t.Errorf("AC#5 cleanup RED: %s still names something after RemoveAll (err=%v)", target, lErr)
				}
			}
		})
		if wd == absRootOf(vol) {
			t.Skipf("this process stands at the volume root %q, so %q and %q name one object here and the pair cannot be planted", wd, rel, abs)
		}
		if err := os.MkdirAll(relDirSide, 0o700); err != nil {
			t.Skipf("the drive-relative side cannot be planted under the process's own directory %q: %v", relDirSide, err)
		}
		if err := os.MkdirAll(filepath.Dir(abs), 0o700); err != nil {
			t.Skipf("the absolute side cannot be planted at the volume root %q: %v", filepath.Dir(abs), err)
		}
		if err := os.WriteFile(rel, []byte("drive-relative"), 0o600); err != nil {
			t.Fatalf("WriteFile through the drive-relative spelling %q: %v", rel, err)
		}
		if err := os.WriteFile(abs, []byte("absolute"), 0o600); err != nil {
			t.Fatalf("WriteFile through the absolute spelling %q: %v", abs, err)
		}
		// One object would show one content after two writes. Two objects show
		// each write standing on its own, read back through its own spelling.
		relGot, err := os.ReadFile(rel)
		if err != nil {
			t.Fatalf("ReadFile(%q): %v", rel, err)
		}
		absGot, err := os.ReadFile(abs)
		if err != nil {
			t.Fatalf("ReadFile(%q): %v", abs, err)
		}
		t.Logf("AC#1 object reading: wd=%s volume=%s", wd, vol)
		t.Logf("AC#1 object reading: drive-relative spelling %q reads %q; absolute spelling %q reads %q", rel, relGot, abs, absGot)
		if string(relGot) == string(absGot) {
			t.Fatalf("the instrument planted nothing: both spellings read back %q, so on this box they are one object", relGot)
		}
		// The ACL dimension, since this is the thing a seal actually rewrites:
		// widen one object only and read both.
		run(t, "icacls", rel, "/grant", "*"+everyoneSID+":(RX)")
		relSIDs, relNames := aclSIDs(t, rel)
		absSIDs, absNames := aclSIDs(t, abs)
		t.Logf("AC#1 ACL reading: %q sids=%v names=%v", rel, relSIDs, relNames)
		t.Logf("AC#1 ACL reading: %q sids=%v names=%v", abs, absSIDs, absNames)
		if !containsSID(relSIDs, everyoneSID) {
			t.Fatalf("the fixture is not wide on the drive-relative side: %s absent from %v", everyoneSID, relSIDs)
		}
		if containsSID(absSIDs, everyoneSID) {
			t.Errorf("AC#1 RED: widening %q also widened %q, so the two spellings are one object and this ticket's premise is wrong here", rel, abs)
		}
		// What the installed pipeline answers for the drive-relative spelling, so
		// the claim "the resolver absolutizes it" is a reading rather than a
		// quote from another package's comment. Judged as a log line only:
		// internal/risk owns that answer. The comparison verdict on these two
		// real spellings is judged by the internal file, which can reach the
		// comparisons directly, and by the seam leg below, which is the face that
		// decides anything in production.
		if got, rErr := winsec.ResolvePath(rel); rErr != nil {
			t.Logf("AC#1 pipeline reading: ResolvePath(%q) refused: %v", rel, rErr)
		} else {
			t.Logf("AC#1 pipeline reading: ResolvePath(%q) = %q (drive-relative spelling names an object under %q)", rel, got.String(), wd)
		}
	})
}

// TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness is the same
// verdict on material that exists on this box: the candidate's two answers are
// the two directories planted above, byte-identical below their volume, one
// absolute and one drive-relative, holding different bytes and different DACLs.
// The guard compares them through the containment leg and has to say whether one
// tree is inside the other.
func TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	vol := filepath.VolumeName(wd)
	fixture := fmt.Sprintf("wisp129-vouch-%d", os.Getpid())
	if vol == "" || wd == absRootOf(vol) {
		t.Skipf("this box gives no drive-relative vs absolute pair to plant (wd=%q volume=%q)", wd, vol)
	}
	relDir := driveRelative129(vol, fixture+`\tree`)
	absDir := driveAbsolute129(vol, filepath.Join(fixture, "tree"))
	t.Cleanup(func() {
		for _, target := range []string{filepath.Join(wd, fixture), absRootOf(vol) + fixture} {
			if rmErr := os.RemoveAll(target); rmErr != nil {
				t.Errorf("AC#5 cleanup RED: RemoveAll(%s): %v", target, rmErr)
				continue
			}
			if _, lErr := os.Lstat(target); !os.IsNotExist(lErr) {
				t.Errorf("AC#5 cleanup RED: %s still names something after RemoveAll (err=%v)", target, lErr)
			}
		}
	})
	if err := os.MkdirAll(filepath.Join(wd, fixture, "tree"), 0o700); err != nil {
		t.Skipf("the drive-relative side cannot be planted under %q: %v", wd, err)
	}
	if err := os.MkdirAll(absDir, 0o700); err != nil {
		t.Skipf("the absolute side cannot be planted at the volume root %q: %v", vol, err)
	}
	probeParent, probeChild := seamProbePair129(t, absDir)
	otherTree := driveAbsolute129(vol, fixture+`-unrelated`)
	candidate := &witness129{
		answers: map[string]string{
			probeParent: relDir,
			probeChild:  absDir + `\leaf`,
			absDir:      otherTree,
		},
		fixed: absDir,
	}
	reason := winsec.TreeOwnershipProbeForTest(candidate, probeParent, probeChild)
	t.Logf("AC#1 real-object verdict: refusal=%q asked=%v (parent answer %q, child answer %q)",
		reason, askedNames(candidate), relDir, absDir+`\leaf`)
	if !candidate.asked[probeParent] || !candidate.asked[probeChild] {
		t.Fatalf("AC#1 instrument measured nothing: the probe pair was not asked (asked %v)", askedNames(candidate))
	}
	if reason == "" {
		t.Errorf("AC#1 RED: the seam's containment witness read %q as sitting inside the tree %q names, and those are two directories on this box - one under the process's own current directory %q, one at the volume root - planted for this leg with different contents",
			absDir+`\leaf`, relDir, wd)
	}
}

// TestAC2NoSealEverActsOnAGloballyDifferentAbsoluteness is the cost side of the
// AC#2 ruling, measured rather than argued: the rule this ticket asks for refuses
// a comparison whose two sides disagree about absoluteness, so the question is
// whether any spelling that survives ResolvePath - the only mint for a path a seal
// acts on - can ever be in the other class. Each row below prints its own verdict,
// because a table that only asserts "some error" cannot tell a refused spelling
// from an accidentally failing one (ticket 108's AC#2 rule).
func TestAC2NoSealEverActsOnAGloballyDifferentAbsoluteness(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	vol := filepath.VolumeName(wd)
	existing := filepath.Join(t.TempDir(), "artifact.txt")
	if err := os.WriteFile(existing, []byte("wisp 129\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	rows := []struct {
		name     string
		spelling string
	}{
		{"existing absolute file", existing},
		{"existing file, slash separators", filepath.ToSlash(existing)},
		{"nonexistent absolute tail", driveAbsolute129(vol, `wisp129-shape\tree\object.txt`)},
		{"nonexistent drive-relative tail", driveRelative129(vol, `wisp129-shape\tree\object.txt`)},
		{"relative spelling, no volume at all", `wisp129-shape\tree\object.txt`},
	}
	for _, r := range rows {
		if r.spelling == "" {
			t.Logf("AC#2 reading: %s is not available on this box (no volume segment in %q)", r.name, wd)
			continue
		}
		got, rErr := winsec.ResolvePath(r.spelling)
		switch {
		case rErr != nil:
			t.Logf("AC#2 reading: %s %q refused: %v", r.name, r.spelling, rErr)
		case !filepath.IsAbs(got.String()):
			t.Errorf("AC#2 RED: ResolvePath(%q) = %q with no error, and the answer is not absolute, so a seal can act on a spelling that has no tree of its own - which is the pair this ticket's rule would refuse",
				r.spelling, got.String())
		default:
			t.Logf("AC#2 reading: %s %q answered %q (absolute)", r.name, r.spelling, got.String())
		}
	}
}

// ------------------------------------------------------------------ AC#1 b ----

// TestAC1SealLandsInAForeignTreeWhenTheSeamAdmittedAnAbsolutenessBlindCandidate
// is the harm reading this ticket was opened for.
//
// The candidate answers the guard's second witness with a drive-relative spelling
// of the tree it already named for the probe parent. Both spellings hang the same
// tail off different roots, so the candidate has vouched for nothing - and the
// guard's own refusal text is the sentence that says why that matters ("the seam
// may not be used to move a seal into another tree"). Once admitted, every seal
// the process performs lands wherever the candidate says: measured here on a real
// victim tree carrying a real inherited Everyone read grant.
//
// Before the fix this goes red three times over: the candidate is installed, the
// caller's own file keeps its Everyone grant, and the victim's is stripped. After
// the fix the guard refuses the same candidate for the same reason it already
// gives an absolute witness naming another tree, and the seal lands where the
// caller pointed.
func TestAC1SealLandsInAForeignTreeWhenTheSeamAdmittedAnAbsolutenessBlindCandidate(t *testing.T) {
	incumbent := winsec.PathResolverInstalled()
	t.Logf("incumbent resolver on this machine: %s", resolverName(incumbent))
	// The seam has to be genuinely free before the install attempt below, or the
	// guard's own single-use latch refuses the candidate for the wrong reason and
	// the reading would be about the latch rather than about the comparison.
	// SetSeamForTest is the test-only door production code does not have (ticket
	// 108's AC#1 welded that shut); everything judged from here on is the real
	// SetPathResolver guard.
	seamAt108(t, nil)

	victim, innocentFile, victimSidsBefore := foreignVictimTree(t)
	foreignFile := filepath.Join(victim, "sub", "keep-me.txt")
	foreignSidsBefore, foreignNamesBefore := aclSIDs(t, foreignFile)
	t.Logf("AC#1 foreign file DACL before: sids=%v names=%v", foreignSidsBefore, foreignNamesBefore)
	if !containsSID(foreignSidsBefore, everyoneSID) {
		t.Fatalf("the victim's own file is not wide, so this instrument could not measure a stripped grant: %v", foreignSidsBefore)
	}

	own := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(own, 0o700); err != nil {
		t.Fatal(err)
	}
	blob := wideFileInOurTree(t, own, "blob.bin")
	blobSidsBefore, blobNamesBefore := aclSIDs(t, blob)
	t.Logf("AC#1 caller's own file DACL before: sids=%v names=%v", blobSidsBefore, blobNamesBefore)

	probeParent, probeChild := seamProbePair129(t, foreignFile)
	const (
		probeTreeAbs = `D:\wisp129-seam\probe-tree`
		probeTreeRel = `D:wisp129-seam\probe-tree`
		movedAbs     = `D:\wisp129-seam\moved-seal\leaf`
		movedDir     = `D:\wisp129-seam\moved-seal`
	)
	candidate := &witness129{
		answers: map[string]string{
			probeParent: probeTreeAbs,
			probeChild:  movedAbs,
			// The second witness: the tree it named for the probe parent, in the
			// other absoluteness class. This is the whole lie. If the guard stopped
			// at its first containment leg it would never ask about movedDir, and
			// the three-way sentinel below reports that as "measured nothing"
			// rather than as a pass.
			movedDir: probeTreeRel,
		},
		fixed: foreignFile,
	}

	logged := captureSeamLog(t)
	winsec.SetPathResolver(candidate)
	installedAs := resolverName(winsec.PathResolverInstalled())
	asked := askedNames(candidate)
	t.Logf("AC#1 install attempt: seam now %s; asked=%v", installedAs, asked)
	if !candidate.asked[probeParent] || !candidate.asked[probeChild] || !candidate.asked[movedDir] {
		t.Errorf("AC#1 instrument measured nothing: the guard did not ask all three shapes (parent=%v child=%v second witness=%v, asked %v)",
			candidate.asked[probeParent], candidate.asked[probeChild], candidate.asked[movedDir], asked)
	}
	admitted := installedAs == resolverName(candidate)
	if admitted {
		t.Errorf("AC#1 RED: the install-time seam ADMITTED a candidate whose second witness vouches for the probe parent's tree only in a drive-relative spelling (%v): it named %q for the probe parent and %q for that answer's own parent, two objects, and the comparison read them as one",
			candidate.answers, probeTreeAbs, probeTreeRel)
	} else {
		t.Logf("AC#1 the candidate was refused; reason on record: %s", strings.TrimSpace(logged.String()))
	}

	// The outcome judgement, at icacls level rather than "the call returned an
	// error": where did the next seal actually land?
	err := winsec.SealFile(blob)
	t.Logf("AC#1 SealFile(%q) -> %v (candidate admitted: %v; seam %s)", blob, err, admitted, installedAs)
	foreignSidsAfter, foreignNamesAfter := aclSIDs(t, foreignFile)
	blobSidsAfter, blobNamesAfter := aclSIDs(t, blob)
	victimSidsAfter, _ := aclSIDs(t, victim)
	t.Logf("AC#1 foreign file DACL after:  sids=%v names=%v", foreignSidsAfter, foreignNamesAfter)
	t.Logf("AC#1 caller's own file DACL after:  sids=%v names=%v", blobSidsAfter, blobNamesAfter)
	if !containsSID(foreignSidsAfter, everyoneSID) {
		t.Errorf("AC#1 RED: %s was stripped from %s, i.e. a seal asked for %s landed inside a tree nobody called and rewrote its DACL. Grants that disappeared: %v. The install-time seam had just admitted the candidate that answers with that path (admitted=%v).",
			everyoneSID, foreignFile, blob, dropped129(foreignSidsBefore, foreignSidsAfter), admitted)
	}
	if containsSID(blobSidsAfter, everyoneSID) {
		t.Errorf("AC#1 RED: the file the caller actually named still grants %s, so the seal never landed where it was pointed: sids=%v names=%v",
			everyoneSID, blobSidsAfter, blobNamesAfter)
	}
	if !equalSIDs(victimSidsAfter, victimSidsBefore) {
		t.Errorf("AC#1 RED: the victim tree directory's DACL moved too: before=%v after=%v", victimSidsBefore, victimSidsAfter)
	}
	if _, statErr := os.Stat(innocentFile); statErr != nil {
		t.Errorf("AC#1 RED: the innocent file behind the victim tree is gone: %v", statErr)
	}
	t.Logf("icacls(foreign) after:\n%s", icaclsRaw(t, foreignFile))
	t.Logf("icacls(blob) after:\n%s", icaclsRaw(t, blob))
}

// dropped129 names the grants that disappeared, which is the per-grant reading
// rather than a count.
func dropped129(before, after []string) []string {
	have := map[string]bool{}
	for _, s := range after {
		have[s] = true
	}
	var out []string
	for _, s := range before {
		if !have[s] {
			out = append(out, s)
		}
	}
	return out
}

// The three spelling builders below exist because filepath.VolumeName answers
// "D:" with its own colon: this file's first run built "D::wisp129-..." by adding
// a second one, which is not a spelling of anything (the OS refused it, and
// filepath.IsAbs said false for the absolute side too, so an instrument that
// looked like the pair under test was measuring nothing). Keeping the two shapes
// in one place each is the fix for that mistake, not a comment about it.

// absRootOf returns the volume root for a volume segment as VolumeName spells it.
func absRootOf(vol string) string { return vol + `\` }

// driveAbsolute129 hangs tail off the volume root: the absolute spelling, the one
// the built-in floor accepts.
func driveAbsolute129(vol, tail string) string { return absRootOf(vol) + tail }

// driveRelative129 hangs the same tail off whatever directory the process
// happens to be standing in on that volume: the drive-relative spelling, the one
// the floor refuses.
func driveRelative129(vol, tail string) string {
	if vol == "" {
		return "" // no volume segment to be relative to
	}
	return vol + tail
}
