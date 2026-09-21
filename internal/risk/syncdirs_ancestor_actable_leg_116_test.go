package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 116 - the BEHAVIOR case whose absence AC#3 of ticket 105 was filed over.
//
// Scope of what this file pins, stated narrowly (read this before generalizing
// it): each test below pins that DELETING the `ares.Actable()` leg at the
// ancestor re-resolve of syncdirs.go resolveTarget (the SECOND Resolve, on the
// nearest existing ancestor) turns a test red. This file does NOT pin that
// Actable()'s decision semantics are correct - that is ticket 102's
// producer-side accounting, already closed and covered by
// pathresolver_expansion_test.go and the static criterion in
// pathresolver_rewrite_account_test.go. A wrong-but-called Actable() would keep
// these tests green; a never-called one does not. Say the coverage at this
// layer and no further.
//
// Why ticket 105's "cannot be triggered on a real machine" claim is false. Its
// argument was that the ancestor string is cut out of C26's own output and so
// cannot hold an unexpanded construct. It is true that
// deepestExistingAncestor does not re-spell anything, and that the ancestor is
// a literal prefix of C26's canonical form. What the argument misses is that
// the two Resolve calls do not expand the same STRING: the first runs
// expandAccounted over the RAW spelling, while the ancestor is cut from the
// form AFTER lexCanonical - and Clean folds a doubled separator. A `%...%` pair
// is delimited by characters, not by path components, so folding the
// separators between them changes WHICH two `%` pair up, and with it the env
// var name that gets looked up:
//
//	RAW       <tree>\<p1><s><s><p2>%...  pair names `WISP116ANC<S><S>B` -> unset
//	ancestor  <tree>\%WISP116ANC<S>B%    pair names `WISP116ANC<S>B`     -> set
//
// so the first Actable() passes (nothing substituted, Rewritten=false, and the
// target does not exist so the walk runs), while the ancestor's Actable()
// genuinely errors. That is the leg, live, reached through IsSyncPath.
//
// Mutations the acceptance side should run. Anchor: the line
// `anceCanon, err := ares.Actable()` in syncdirs.go, measured at
// internal/risk/syncdirs.go:226 on 9141d4c - re-grep it, line numbers drift.
// Keep the re-resolve and delete only the leg, in either of its two spellings:
//
//	M1  `anceCanon := ares.Canonical`
//	    the substituted tree becomes the comparison candidate, it is outside
//	    the registered root, and the write is RELEASED UNSUSPECTED -> the first
//	    assertion of TestSyncAncestorActableLegFailsClosedOnMovedAncestor116
//	    goes red on Sync:false.
//	M2  `anceCanon, _ := ares.Actable()`
//	    the empty canonical falls into the `anceCanon == ""` branch below, so
//	    the verdict stays fail-closed but loses the expansion account -> the
//	    same test goes red on its Why assertion.
//
// Only that one test goes red under either mutation: the two AC#2 halves below
// stay green, which is what shows the red names the leg rather than a
// coincidence of a broken harness.
//
// Reachability, without inflation: all three preconditions are supplyable -
// a directory name containing `%`, a write target spelled with a doubled
// separator, and an env var whose NAME contains a separator (measured settable
// via os.Setenv on both platforms here). The third is only moderately-to-weakly
// reachable in a real deployment. What this file establishes is that the
// absolute claim "triggerable: never" is false. It is not a claim that an
// external attacker can drive this into a hole.

const (
	// ancPairName is the env var the ANCESTOR spelling resolves to: one
	// separator inside the %...% pair. This file sets it.
	ancPairName = "WISP116ANC" + sepStr + "B"
	// rawPairName is the env var the RAW spelling resolves to: two separators,
	// because Clean has not folded the doubled one yet. It must stay unset,
	// which is exactly what lets the first Actable() pass.
	rawPairName = "WISP116ANC" + sepStr + sepStr + "B"
	// firstPairName is an ordinary, single-separator-free env var, used by the
	// AC#2 half that keeps the loud refusal of a legit rewrite pinned.
	firstPairName = "WISP116FIRST"

	// movedDir is the tree the substitution moves the ancestor onto: a sibling
	// of the sync tree, so it is genuinely outside the registered guard.
	movedDir = "wisp116-escape"
	// plainDir is a third tree, neither sync root nor under the test home, used
	// as the release control.
	plainDir = "wisp116-plain"
	// hiddenLeaf is the not-yet-existing write target under the % pair, which
	// is what forces resolveTarget to take the ancestor walk.
	hiddenLeaf = "new.md"
)

// ticket116MovedTree lays out `<scratch>/001` as the registered sync tree -
// including the literal directory names `%WISP116ANC` and `B%` inside it, which
// NTFS and ext4 both accept - plus `<scratch>/<movedDir>` as a real tree that is
// NOT under the sync root, and installs the env var that only the cleaned
// spelling resolves to. Everything created lives in the test's own scratch
// directory; t.Setenv restores the environment.
func ticket116MovedTree(t *testing.T) (tree, moved, raw, ancestor string) {
	t.Helper()
	tree = t.TempDir()
	if os.Getenv(rawPairName) != "" {
		t.Fatalf("cannot build the shape: the two-separator name %q is already set in this process", rawPairName)
	}
	if err := os.MkdirAll(filepath.Join(tree, "%WISP116ANC", "B%"), 0o755); err != nil {
		t.Fatalf("mkdir the percent-spelled directory pair: %v", err)
	}
	moved = filepath.Join(filepath.Dir(tree), movedDir)
	if err := os.MkdirAll(moved, 0o755); err != nil {
		t.Fatalf("mkdir the tree expansion moves the ancestor onto: %v", err)
	}
	// Relative on purpose: the substitution is spliced in where the pair stood,
	// so `..` is folded by the same Clean that folded the doubled separator,
	// and the ancestor lands on a directory that really exists. That matters -
	// if the moved-onto tree did not exist, the `!ares.Resolved` branch below
	// the leg could be what saves the verdict, and a mutation red would not
	// prove the leg was doing anything.
	t.Setenv(ancPairName, ".."+sepStr+movedDir)

	raw = tree + sepStr + "%WISP116ANC" + sepStr + sepStr + "B%" + sepStr + hiddenLeaf
	ancestor = tree + sepStr + "%WISP116ANC" + sepStr + "B%"
	return tree, moved, raw, ancestor
}

// ticket116Engine returns a Provenance whose only sync root is `tree`, with a
// deliberately unrelated home so the under-profile suspect net cannot be what
// answers (the same precondition ticket 105's junction cases assert).
//
// SyncDetectionComplete() is NOT asserted here, unlike those tests: on POSIX no
// injected root can be C26-canonical, because resolveHandle is unsupported
// there (see the DEFERRED marker in pathresolver_other.go), so complete is
// legitimately false on that platform and asserting it would make this file a
// Windows-only test. The burden the assertion carries - "a Sync:false verdict is
// reachable on this engine, so a mutation red means the leg let the write
// through and not merely that the net failed to fire" - is carried instead by
// the release control inside the ancestor test below, which holds on both
// platforms.
func ticket116Engine(t *testing.T, tree string) *Provenance {
	t.Helper()
	p := NewProvenance(ProvOptions{
		NoProbe:   true,
		HomeDir:   t.TempDir(),
		SyncRoots: []SyncRoot{{Provider: "Ticket116Sync", Path: tree, Source: "env"}},
	})
	if len(p.SyncRoots()) != 1 {
		t.Fatalf("precondition: exactly one injected sync root expected, got %v", p.SyncRoots())
	}
	return p
}

// TestSyncAncestorActableLegFailsClosedOnMovedAncestor116 is the case AC#1 asks
// for: the ancestor leg's Actable() is the only thing standing between this
// write verdict and a "not a sync dir" answer.
func TestSyncAncestorActableLegFailsClosedOnMovedAncestor116(t *testing.T) {
	tree, moved, raw, ancestor := ticket116MovedTree(t)
	p := ticket116Engine(t, tree)

	// --- preconditions, asserted rather than assumed -----------------------
	// 1. The RAW spelling resolves clean and is not handle-resolvable, so the
	//    first Actable() passes and resolveTarget falls through to the walk.
	first, err := Resolve(raw, nil)
	if err != nil {
		t.Fatalf("precondition: the first Resolve denied the target: %v", err)
	}
	if first.Rewritten {
		t.Fatalf("precondition: the raw spelling must NOT be rewritten (rewrites=%v): if it were, the first "+
			"Actable() would answer here and the ancestor leg would never run", first.Rewrites)
	}
	if first.Resolved {
		t.Fatalf("precondition: %q must not be handle-resolvable as a whole, otherwise there is no ancestor walk", raw)
	}
	canon, err := first.Actable()
	if err != nil {
		t.Fatalf("precondition: the first Actable() must pass, got %v", err)
	}

	// 2. The ancestor is cut out of C26's OWN output - the fact ticket 105
	//    rested its claim on - and it still carries the unexpanded %...% pair,
	//    because the doubled separator that kept that pair unresolvable was
	//    folded away two pipeline steps earlier.
	gotAnc, rest := deepestExistingAncestor(canon)
	if gotAnc != ancestor {
		t.Fatalf("precondition: deepestExistingAncestor(%q) = %q, want %q", canon, gotAnc, ancestor)
	}
	if len(rest) != 1 || rest[0] != hiddenLeaf {
		t.Fatalf("precondition: the remainder must be exactly the one missing leaf, got %v", rest)
	}
	if !strings.Contains(gotAnc, "%") {
		t.Fatalf("precondition: the ancestor taken from C26's output must still hold a percent construct, got %q", gotAnc)
	}

	// 3. Re-resolving THAT ancestor is where the account fires. Nothing above
	//    this line can answer for the leg.
	ares, err := Resolve(gotAnc, nil)
	if err != nil {
		t.Fatalf("precondition: the ancestor Resolve denied instead of reporting: %v", err)
	}
	if !ares.Rewritten {
		t.Fatalf("precondition: the ancestor spelling must be the one C26 substituted (rewrites=%v)", ares.Rewrites)
	}
	if _, e := ares.Actable(); e == nil {
		t.Fatal("precondition: the ancestor's Actable() must be the thing that errors here")
	}

	// 4. Release control: this engine does answer Sync:false for a write that
	//    is genuinely outside the guard, including for the very tree the
	//    substitution moves onto. Without this, a Sync:false red under a
	//    mutation could be an artifact of the harness, and a Sync:true
	//    no-mutation-green could be the suspect net papering over everything.
	if st := p.IsSyncPath(filepath.Join(moved, hiddenLeaf)); st.Sync {
		t.Fatalf("precondition: %q is not a sync tree and must read as such, got %+v", moved, st)
	}
	if st := p.IsSyncPath(filepath.Join(filepath.Dir(tree), plainDir, hiddenLeaf)); st.Sync {
		t.Fatalf("precondition: an unrelated non-sync write target must read Sync:false, got %+v", st)
	}

	// --- the verdict -------------------------------------------------------
	st := p.IsSyncPath(raw)
	t.Logf("raw=%q\n anc=%q\n moved=%q\nverdict: Sync=%v Root=%+v Why=%q", raw, gotAnc, moved, st.Sync, st.Root, st.Why)
	if !st.Sync {
		t.Errorf("IsSyncPath(%q) = Sync:false (why=%q): the ancestor's spelling resolved onto %q, which is "+
			"outside the registered sync root %q. Deleting resolveTarget's ares.Actable() leg releases this write "+
			"as NOT a sync dir, so the fs.write channel would let it through un-suspected.",
			raw, st.Why, moved, tree)
	} else if st.Root.Source != "suspect-fallback" {
		t.Errorf("verdict root = %+v: the only legitimate Sync:true here is the fail-closed suspect one; a "+
			"matched root would mean the candidate got anchored on a tree C26 never verified", st.Root)
	}
	if !strings.Contains(st.Why, "expansion moved") {
		t.Errorf("Why = %q: an ancestor C26 moved onto another tree must be refused WITH the expansion account in "+
			"the operator-visible reason, and that text only exists while the leg is called - it is how a dropped "+
			"Actable() hides behind the generic unverified-ancestor message.", st.Why)
	}
}

// TestSyncFirstActableLegStillRefusesEnvRewriteAtVerdict116 is AC#2's first
// half: the leg this file's shape skips is NOT skipped for a spelling whose own
// RAW form expands onto another tree. Same engine, same moved-onto tree, an
// ordinary single-construct env var - the answer stays a loud fail-closed.
func TestSyncFirstActableLegStillRefusesEnvRewriteAtVerdict116(t *testing.T) {
	tree, _, _, _ := ticket116MovedTree(t)
	t.Setenv(firstPairName, ".."+sepStr+movedDir)

	raw := tree + sepStr + "%" + firstPairName + "%" + sepStr + hiddenLeaf
	res, err := Resolve(raw, nil)
	if err != nil {
		t.Fatalf("precondition: Resolve denied instead of reporting: %v", err)
	}
	if !res.Rewritten {
		t.Fatalf("precondition: the raw spelling of a set %s must be accounted as a rewrite, got %+v", firstPairName, res)
	}
	if _, e := res.Actable(); e == nil {
		t.Fatal("the first Actable() leg must refuse a spelling whose own RAW form moved trees")
	}

	st := ticket116Engine(t, tree).IsSyncPath(raw)
	if !st.Sync || st.Root.Source != "suspect-fallback" {
		t.Errorf("IsSyncPath(%q) = %+v: a legitimate percent-env rewrite onto another tree must stay a loudly "+
			"refused, fail-closed sync-suspect verdict; ticket 116's ancestor case may not have relaxed it", raw, st)
	}
	if !strings.Contains(st.Why, "expansion moved") {
		t.Errorf("Why = %q: want the refused expansion named in the write verdict too", st.Why)
	}
}

// TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116 is the same half
// for the other construct the contract expands: a leading `~`. The home is
// redirected into the test's own scratch directory (os.UserHomeDir reads
// USERPROFILE on Windows and HOME elsewhere, both via the process environment),
// so nothing is created under a real profile and there is no skip path.
func TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116(t *testing.T) {
	tree := t.TempDir()
	fakeHome := t.TempDir()
	p := ticket116Engine(t, tree) // built while the environment is still pristine
	t.Setenv("HOME", fakeHome)
	t.Setenv("USERPROFILE", fakeHome)

	raw := "~" + sepStr + "wisp116-tilde-probe" + sepStr + hiddenLeaf
	res, err := Resolve(raw, nil)
	if err != nil {
		t.Fatalf("precondition: Resolve denied instead of reporting: %v", err)
	}
	if !res.Rewritten {
		t.Fatalf("precondition: a leading ~ must be accounted as a rewrite, got %+v", res)
	}
	if !strings.HasPrefix(strings.ToLower(res.Canonical), strings.ToLower(fakeHome)) {
		t.Fatalf("precondition: the ~ spelling must resolve under the relocated home %q, got %q", fakeHome, res.Canonical)
	}
	if _, e := res.Actable(); e == nil {
		t.Fatal("the first Actable() leg must refuse a ~-spelled target that names another tree")
	}

	st := p.IsSyncPath(raw)
	if !st.Sync || st.Root.Source != "suspect-fallback" {
		t.Errorf("IsSyncPath(%q) = %+v: a ~-spelled write target must stay a refused sync-suspect verdict", raw, st)
	}
	if !strings.Contains(st.Why, "expansion moved") {
		t.Errorf("Why = %q: want the home expansion named as the reason this write was refused", st.Why)
	}
}

// TestSyncAncestorLegStillMatchesPlainSameTreeWrite116 is AC#2's second half:
// the normal case, ancestor and target in the same tree, nothing to expand.
// The leg's ordinary job - anchor the candidate on the existing ancestor so the
// comparison can match the registered root - must keep working, and must match
// as a DETECTED root rather than fall back to a blanket refusal. This one stays
// green under both M1 and M2, which is why it belongs next to the red.
func TestSyncAncestorLegStillMatchesPlainSameTreeWrite116(t *testing.T) {
	tree := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tree, "notes"), 0o755); err != nil {
		t.Fatalf("mkdir the existing ancestor directory: %v", err)
	}
	raw := filepath.Join(tree, "notes", hiddenLeaf)

	res, err := Resolve(raw, nil)
	if err != nil {
		t.Fatalf("precondition: Resolve denied the plain target: %v", err)
	}
	if res.Rewritten {
		t.Fatalf("precondition: nothing here may expand, got %+v", res)
	}
	if res.Resolved {
		t.Fatalf("precondition: %q must NOT be handle-resolvable as a whole, otherwise the ancestor leg is not "+
			"the thing under test here", raw)
	}

	st := ticket116Engine(t, tree).IsSyncPath(raw)
	if !st.Sync {
		t.Errorf("IsSyncPath(%q) = Sync:false (why=%q): a new file inside the registered sync tree must still be "+
			"caught through the ancestor leg", raw, st.Why)
	} else if st.Root.Provider != "Ticket116Sync" || st.Root.Source != "env" {
		t.Errorf("matched root = %+v, want the registered env root: the suspect fallback answering here would hide "+
			"a broken ancestor leg behind a blanket refusal", st.Root)
	}
}
