//go:build windows

// Ticket 108 AC#1 / AC#3 / AC#4 probes: the three ways ticket 103's acceptance
// agent walked through ticket 103's own seam guard (P1b, P3 and the R-103-1
// tree-ownership half), turned into rerunnable in-repo cases. Written and run
// BEFORE the fix, so the red is on record; the shapes are copied from
// docs/evidence/s1/103-adversarial-acceptance.md sections AC1-E and AC2-E.
//
// Like seam_guard_windows_test.go this file is in the *external* test package:
// the threat model is code outside internal/winsec that can reach the exported
// setter, which is exactly what P1b exploited.
package winsec_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// treeMovingResolver108 is P1b's fake: it answers every spelling, whatever it is,
// with one fixed clean absolute path in somebody else's tree. It passes today's
// conformance probe because the fixed answer is itself floor-clean (the probe
// only rejects "refuse" or "a spelling the floor refuses", so a rewrite into a
// different clean tree falls through the second leg) - which is R-103-1.
type treeMovingResolver108 struct{ fixed string }

func (r treeMovingResolver108) Resolve(string) (string, error) { return r.fixed, nil }

// seamAt108 puts r into the sealing seam for the duration of one test and puts
// the incumbent back afterwards.
//
// RED-COMMIT SHAPE: the only exported way to reach the floor today is
// SetPathResolver(nil), and that nil branch is precisely the door ticket 108 is
// closing (it is what lets P1b free the seam and re-install a fake). The fix
// commit replaces this body with the test-only hook in export_test.go, so that
// the post-fix reading of these probes does not depend on the door being open.
func seamAt108(t *testing.T, r winsec.C26Resolver) {
	t.Helper()
	incumbent := winsec.PathResolverInstalled()
	winsec.SetPathResolver(nil)
	if r != nil {
		winsec.SetPathResolver(r)
	}
	t.Cleanup(func() {
		winsec.SetPathResolver(nil)
		winsec.SetPathResolver(incumbent)
	})
}

// wideFileInOurTree is our own file carrying an explicit Everyone grant: the
// marker proving that a refused seal left it alone and an accepted one stripped
// it. Ticket 108's judgement is at SID level, not "the call returned an error".
func wideFileInOurTree(t *testing.T, dir, name string) string {
	t.Helper()
	f := filepath.Join(dir, name)
	if err := os.WriteFile(f, []byte("ours, and wide on purpose"), 0o600); err != nil {
		t.Fatal(err)
	}
	run(t, "icacls", f, "/inheritance:r", "/grant:r",
		"*"+currentSID(t)+":(F)", "*"+systemSID+":(F)", "*"+adminsSID+":(F)", "*"+everyoneSID+":(RX)")
	sids, names := aclSIDs(t, f)
	t.Logf("our wide file %s: sids=%v names=%v", name, sids, names)
	if !containsSID(sids, everyoneSID) {
		t.Fatalf("fixture is not wide: %s absent from %v", everyoneSID, sids)
	}
	return f
}

// ----------------------------------------------------------------- AC#1 / P1b --

// TestAC1SeamCannotBeFreedThenGivenATreeMovingFake is probe P1b: free the seam
// with SetPathResolver(nil), then install a fake that rewrites every answer into
// another tree. Today the guard accepts it, SealFile returns nil and the foreign
// S-1-1-0 grant is stripped - the exact fail-open ticket 103 AC#1 claimed to
// prevent, hence the rejected-needs-fix that opened this ticket.
func TestAC1SeamCannotBeFreedThenGivenATreeMovingFake(t *testing.T) {
	incumbent := winsec.PathResolverInstalled()
	// RED-COMMIT SHAPE: this restore leans on SetPathResolver(nil) being allowed,
	// which is the door ticket 108 is closing. The fix commit moves it onto the
	// test-only hook in export_test.go. Note what the leak costs if it is missing:
	// the accepted fake then answers every later PrivateDirAll in this binary,
	// which reports success and creates nothing where the caller pointed.
	t.Cleanup(func() {
		winsec.SetPathResolver(nil)
		winsec.SetPathResolver(incumbent)
	})
	t.Logf("incumbent before the freeing step: %s", resolverName(incumbent))

	victim, innocentFile, victimSidsBefore := foreignVictimTree(t)
	foreignFile := filepath.Join(victim, "sub", "keep-me.txt")
	foreignSidsBefore, foreignNamesBefore := aclSIDs(t, foreignFile)
	t.Logf("foreign file DACL before: sids=%v names=%v", foreignSidsBefore, foreignNamesBefore)

	own := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(own, 0o700); err != nil {
		t.Fatal(err)
	}
	blob := wideFileInOurTree(t, own, "blob.bin")
	blobSidsBefore, _ := aclSIDs(t, blob)

	// Leg 1: the seam must not be detachable, so the incumbent (or the floor) is
	// still what ResolvePath sees after the nil call.
	logged := captureSeamLog(t)
	winsec.SetPathResolver(nil)
	if resolverName(winsec.PathResolverInstalled()) != resolverName(incumbent) {
		t.Errorf("AC#1 RED: SetPathResolver(nil) detached the seam, which is what makes it one-use rather than one-way (was %s, now %s)",
			resolverName(incumbent), resolverName(winsec.PathResolverInstalled()))
	}

	// Leg 2: with or without the freeing step, the tree-moving fake must not get in.
	winsec.SetPathResolver(treeMovingResolver108{fixed: foreignFile})
	if got := resolverName(winsec.PathResolverInstalled()); got == resolverName(treeMovingResolver108{}) {
		t.Errorf("AC#1 RED: after the nil reset the guard ACCEPTED a fake that rewrites every answer to %s (seam is now %s)",
			foreignFile, got)
	}
	if !strings.Contains(logged.String(), "treeMovingResolver108") {
		t.Errorf("AC#1 RED: the refused install left no audit record naming the fake; log was: %s", logged.String())
	}

	// Leg 3: the outcome judgement - sealing our own blob may not end up narrowing
	// the foreign tree, and neither DACL may move.
	if err := winsec.SealFile(blob); err == nil {
		t.Errorf("AC#1 RED: SealFile(%s) returned nil while a tree-moving fake was in the seam", blob)
	} else {
		t.Logf("SealFile through the refused/kept seam: %v", err)
	}
	foreignSidsAfter, foreignNamesAfter := aclSIDs(t, foreignFile)
	t.Logf("foreign file DACL after:  sids=%v names=%v", foreignSidsAfter, foreignNamesAfter)
	if !containsSID(foreignSidsAfter, everyoneSID) {
		t.Errorf("AC#1 RED: %s (Everyone) was stripped from the foreign file, i.e. its DACL was rewritten: %v",
			everyoneSID, foreignSidsAfter)
	}
	if !equalSIDs(foreignSidsAfter, foreignSidsBefore) {
		t.Errorf("AC#1 RED: foreign file DACL changed: before=%v after=%v", foreignSidsBefore, foreignSidsAfter)
	}
	if _, err := os.Stat(innocentFile); err != nil {
		t.Errorf("AC#1: the innocent file behind the foreign tree must be untouched: %v", err)
	}
	blobSidsAfter, _ := aclSIDs(t, blob)
	victimSidsAfter, _ := aclSIDs(t, victim)
	if !equalSIDs(victimSidsAfter, victimSidsBefore) || !containsSID(victimSidsAfter, everyoneSID) {
		t.Errorf("AC#1 RED: the foreign tree directory's DACL moved: before=%v after=%v", victimSidsBefore, victimSidsAfter)
	}
	if !containsSID(blobSidsAfter, everyoneSID) {
		t.Errorf("AC#1 RED: %s was stripped from our own wide blob too, so the fake's answer is being acted on: %v",
			everyoneSID, blobSidsAfter)
	}
	if !equalSIDs(blobSidsAfter, blobSidsBefore) {
		t.Errorf("AC#1 RED: our blob DACL changed: before=%v after=%v", blobSidsBefore, blobSidsAfter)
	}
	t.Logf("icacls(foreign) after:\n%s", icaclsRaw(t, foreignFile))
	t.Logf("icacls(blob) after:\n%s", icaclsRaw(t, blob))
}

// -------------------------------------------------------------------- AC#4 ----

// TestAC4TreeOwnershipIsPartOfTheConformanceContract is R-103-1's second half:
// even on a seam that is genuinely free (the built-in floor in place, nothing
// installed), installing a resolver must require more than "its answers are
// floor-clean". The answer also has to stay in the tree the caller named, and a
// resolver that cannot answer that way may not be installed at all.
func TestAC4TreeOwnershipIsPartOfTheConformanceContract(t *testing.T) {
	victim, _, _ := foreignVictimTree(t)
	foreignFile := filepath.Join(victim, "sub", "keep-me.txt")
	foreignSidsBefore, _ := aclSIDs(t, foreignFile)

	seamAt108(t, nil)
	if got := winsec.PathResolverInstalled(); got != nil {
		t.Fatalf("this probe needs the floor in place, seam is %s", resolverName(got))
	}

	logged := captureSeamLog(t)
	winsec.SetPathResolver(treeMovingResolver108{fixed: foreignFile})
	if got := resolverName(winsec.PathResolverInstalled()); got == resolverName(treeMovingResolver108{}) {
		t.Errorf("AC#4 RED: the seam accepted a resolver whose answers name a tree nobody called (%s): the conformance probe checks spelling shape only, never tree ownership", foreignFile)
	} else {
		t.Logf("the tree-moving fake was refused; seam stayed %s", got)
	}
	reason := strings.ToLower(logged.String())
	if !strings.Contains(reason, "tree") && !strings.Contains(reason, "ownership") {
		t.Errorf("AC#4 RED: the refusal names no tree-ownership reason; log was: %s", logged.String())
	}
	sidsAfter, _ := aclSIDs(t, foreignFile)
	if !containsSID(sidsAfter, everyoneSID) || !equalSIDs(sidsAfter, foreignSidsBefore) {
		t.Errorf("AC#4 RED: the foreign tree's DACL moved anyway: before=%v after=%v", foreignSidsBefore, sidsAfter)
	}
}

// ----------------------------------------------------------------- AC#3 / P3 --

// TestAC3PlacementFloorHoldsForEverySeparatorSpelling is probe P3: the built-in
// floor's own reparse walk (platformVerifyPlacement) cuts only on the backslash,
// so the same input spelled with forward slashes gets past it and rewrites a
// foreign DACL. Predates ticket 103 (placement_windows.go is not in 0717bf2) but
// is the same mistake in the same judgement, so ticket 108 collects it.
func TestAC3PlacementFloorHoldsForEverySeparatorSpelling(t *testing.T) {
	seamAt108(t, nil)
	for _, shape := range sealSpellings(t) {
		t.Run(shape.name, func(t *testing.T) {
			victim, innocentFile, _ := foreignVictimTree(t)
			own := filepath.Join(t.TempDir(), "data")
			if err := winsec.PrivateDirAll(own, 0o700); err != nil {
				t.Fatal(err)
			}
			link := filepath.Join(own, "link")
			if _, statErr := os.Stat(own); statErr != nil {
				t.Fatalf("fixture: PrivateDirAll(%s) reported success but it is not there: %v", own, statErr)
			}
			t.Logf("fixture: own=%s seam=%s", own, resolverName(winsec.PathResolverInstalled()))
			makeJunction103(t, link, victim)
			through := filepath.Join(link, "sub", "keep-me.txt")
			spell := shape.spell(through, own, victim)
			before, namesBefore := aclSIDs(t, through)
			t.Logf("spelling %q: foreign DACL before sids=%v names=%v", spell, before, namesBefore)
			if !containsSID(before, everyoneSID) {
				t.Fatalf("fixture is not wide: %s absent from %v", everyoneSID, before)
			}
			err := winsec.SealFile(spell)
			if err == nil {
				t.Errorf("AC#3 RED: SealFile(%q) through a junction returned nil", spell)
			} else if !strings.Contains(strings.ToLower(err.Error()), "reparse") &&
				!errors.Is(err, winsec.ErrUnresolvedPath) {
				t.Logf("refused, reason: %v", err)
			}
			after, namesAfter := aclSIDs(t, through)
			t.Logf("spelling %q: foreign DACL after  sids=%v names=%v", spell, after, namesAfter)
			if !containsSID(after, everyoneSID) {
				t.Errorf("AC#3 RED: %s stripped from the foreign file by spelling %q - the floor's reparse walk did not see the junction: %v",
					everyoneSID, spell, after)
			}
			if !equalSIDs(after, before) {
				t.Errorf("AC#3 RED: foreign DACL changed for spelling %q: before=%v after=%v", spell, before, after)
			}
			if _, statErr := os.Stat(innocentFile); statErr != nil {
				t.Errorf("AC#3 RED: the innocent file is gone after spelling %q: %v", spell, statErr)
			}
			if _, lErr := os.Lstat(link); lErr != nil {
				t.Errorf("the junction itself must survive a refused seal: %v", lErr)
			}
		})
	}
}

// spelling108 names one input shape and how to build it from a native
// (backslash) path. The list is AC#2/AC#3's "per-shape conclusion" requirement:
// every row reports its own verdict, because a table that only asserts "some
// error" cannot tell a refused spelling from an accidentally failing one.
type spelling108 struct {
	name  string
	spell func(native, own, victim string) string
}

// sealSpellings is the shape list: native is the control, the rest are the ones
// P2/P3 measured as bypasses.
func sealSpellings(t *testing.T) []spelling108 {
	t.Helper()
	return []spelling108{
		{"native-backslash", func(n, _, _ string) string { return n }},
		{"all-forward-slash", func(n, _, _ string) string { return strings.ReplaceAll(n, `\`, "/") }},
		{"mixed-separators", func(n, own, victim string) string {
			return strings.ReplaceAll(own, `\`, "/") + `\link/sub/keep-me.txt`
		}},
		{"trailing-separator", func(n, _, _ string) string { return n + "/" }},
		{"doubled-separator", func(n, _, _ string) string { return strings.Replace(n, `\sub`, `//sub`, 1) }},
		{"dot-segment", func(n, _, _ string) string { return strings.Replace(n, `\link`, `\link\.`, 1) }},
		{"extended-length-prefix", func(n, _, _ string) string { return `\\?\` + n }},
		{"volume-only-relative-tail", func(n, own, _ string) string {
			return filepath.Base(own) + `\link/sub/keep-me.txt`
		}},
	}
}

// ----------------------------------------------------------------- AC#2 / P2 --

// TestAC2AncestorGuardHoldsForEverySeparatorSpelling is probe P2: RemoveUnlinked
// does not go through the resolver at all, so its only defence is
// firstLinkAncestor - which today cuts on filepath.Separator only, so a
// junction-crossing spelling built with forward slashes has exactly one
// component in the guard's eyes and no ancestor is Lstat'ed.
func TestAC2AncestorGuardHoldsForEverySeparatorSpelling(t *testing.T) {
	for _, shape := range sealSpellings(t) {
		t.Run(shape.name, func(t *testing.T) {
			victim, innocentFile, victimSids := foreignVictimTree(t)
			own := filepath.Join(t.TempDir(), "data")
			if err := winsec.PrivateDirAll(own, 0o700); err != nil {
				t.Fatal(err)
			}
			link := filepath.Join(own, "link")
			if _, statErr := os.Stat(own); statErr != nil {
				t.Fatalf("fixture: PrivateDirAll(%s) reported success but it is not there: %v", own, statErr)
			}
			t.Logf("fixture: own=%s seam=%s", own, resolverName(winsec.PathResolverInstalled()))
			makeJunction103(t, link, victim)
			through := filepath.Join(link, "sub", "keep-me.txt")
			spell := shape.spell(through, own, victim)

			err := winsec.RemoveUnlinked(spell)
			t.Logf("RemoveUnlinked(%q) -> err=%v", spell, err)
			if err == nil {
				t.Errorf("AC#2 RED: RemoveUnlinked(%q) returned nil for a spelling that reaches the leaf through a junction", spell)
			} else if !errors.Is(err, winsec.ErrIsReparsePoint) {
				t.Logf("refused, but not with ErrIsReparsePoint: %v", err)
			}
			if _, statErr := os.Stat(innocentFile); statErr != nil {
				t.Errorf("AC#2 RED: the foreign file behind the junction is GONE after spelling %q: %v", spell, statErr)
			}
			if _, statErr := os.Stat(filepath.Join(victim, "sub")); statErr != nil {
				t.Errorf("AC#2 RED: the foreign directory was damaged by spelling %q: %v", spell, statErr)
			}
			sidsAfter, _ := aclSIDs(t, victim)
			if !equalSIDs(sidsAfter, victimSids) {
				t.Errorf("AC#2 RED: foreign tree DACL changed for spelling %q: before=%v after=%v", spell, victimSids, sidsAfter)
			}
			if _, lErr := os.Lstat(link); lErr != nil && !errors.Is(lErr, fs.ErrNotExist) {
				t.Errorf("the junction could not be re-read: %v", lErr)
			}
		})
	}
}

// TestAC2AncestorGuardStillUnlinksAPlainFile is the positive control that keeps
// a separator fix from being scored by a guard that refuses the world: an
// ordinary spelling of a real file with no link in its ancestors is still
// removed, and a link spelled with forward slashes is still unlinked itself.
func TestAC2AncestorGuardStillUnlinksAPlainFile(t *testing.T) {
	victim, innocentFile, _ := foreignVictimTree(t)
	root := t.TempDir()
	f := filepath.Join(root, "plain.txt")
	if err := os.WriteFile(f, []byte("stray"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := winsec.RemoveUnlinked(filepath.ToSlash(f)); err != nil {
		t.Fatalf("RemoveUnlinked of a slash-spelled plain file must work: %v", err)
	}
	if _, err := os.Stat(f); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the plain file was not removed: %v", err)
	}
	link := filepath.Join(root, "stray-link")
	makeJunction103(t, link, victim)
	if err := winsec.RemoveUnlinked(filepath.ToSlash(link)); err != nil {
		t.Fatalf("RemoveUnlinked of a slash-spelled standalone link must work: %v", err)
	}
	if _, err := os.Lstat(link); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the link survived its own unlink: %v", err)
	}
	if _, err := os.Stat(innocentFile); err != nil {
		t.Errorf("unlinking the link must not touch its target: %v", err)
	}
}
