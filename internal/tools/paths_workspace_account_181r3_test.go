package tools

// Ticket 181 AC#7: the canonicalizer is the PRODUCER of the workspace rewrite
// account the panel renders. These tests drive the real C26 resolver over real
// directories - the scripted seam lives in internal/panel, and what is nailed
// here is that the values come out of this package's own book rather than out of
// a test's struct literal.
//
// What each one is for:
//
//   - TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing: the account
//     travels. Before this ticket WorkspaceRoot returned a bare string and
//     internal/panel could only ever render rewritten=false, so the field existed
//     and spoke nothing.
//   - TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26: the field is a
//     READING of C26's answer, not a constant. The rewritten account here is the
//     one risk.Resolve actually returned for a %VAR% spelling of a directory that
//     exists, the same construction
//     TestWorkspaceSwitchRefusesAnExpandedSpelling uses - and the refusal that
//     test requires still fires, unchanged, in the first sub-case.
//   - TestSetWorkspaceRootRefusesAnAccountAboutAnotherTree /
//     TestClearWorkspaceDropsTheAccountWithTheRoot: the two ways the new record
//     could rot (an account attached to a coordinate it does not describe, and an
//     account outliving the narrowing it was recorded for).

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

func TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing(t *testing.T) {
	base := sealableTempDir124(t)
	wsA := mkDir(t, base, "alpha")
	canon := NewPathCanonicalizer([]string{wsA}, nil)
	if u := canon.UnusableRoots(); len(u) > 0 {
		t.Fatalf("fixture premise broken, unusable roots: %v", u)
	}

	// Spell the tree the way a composer's picker does - separators the caller
	// used, not the ones C26 answers with - so Spelling is a value that can only
	// arrive by being carried, never by being re-derived from the canonical.
	spelled := filepath.ToSlash(wsA)
	res, err := canon.ResolveWorkspace(spelled)
	if err != nil {
		t.Fatalf("ResolveWorkspace(%s) = %v, want success", spelled, err)
	}
	actable, aErr := res.Actable()
	if aErr != nil {
		t.Fatalf("Actable() on the fixture's own spelling: %v", aErr)
	}
	if err := canon.SetWorkspaceRoot(actable, res); err != nil {
		t.Fatalf("SetWorkspaceRoot: %v", err)
	}

	got := canon.WorkspaceRoot()
	if got.Canonical != actable {
		t.Errorf("WorkspaceRoot().Canonical = %q, want the tree C26 authorised %q", got.Canonical, actable)
	}
	if got.Spelling != spelled {
		t.Errorf("WorkspaceRoot().Spelling = %q, want the caller's own string %q verbatim - "+
			"a producer that drops the spelling drops ticket 102's traceability: a verdict taken on "+
			"Canonical could no longer be traced back to the form it was asked about", got.Spelling, spelled)
	}
	if got.Rewritten || len(got.Rewrites) != 0 {
		t.Errorf("rewritten=%v rewrites=%v for a spelling C26 substituted nothing in", got.Rewritten, got.Rewrites)
	}
	if got.Resolved != res.Resolved || got.Reparse != res.Reparse {
		t.Errorf("the answer disagrees with the Result it was recorded from: got %+v, recorded %+v", got, res)
	}

	// The narrowing itself must still be the fold key this package has always
	// compared against: the account is added to the door, not swapped for it.
	other := mkDir(t, base, "beta")
	if canon.InAllowlist(mustCanonical(t, filepath.Join(other, "note.txt"))) {
		t.Error("recording an account widened the narrowing: a write outside the workspace is still in allowlist")
	}
	if !canon.InAllowlist(mustCanonical(t, filepath.Join(wsA, "note.txt"))) {
		t.Error("the narrowed workspace stopped containing its own tree")
	}
}

func TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26(t *testing.T) {
	const name = "WISP181R3ACC"
	base := sealableTempDir124(t)
	expanded := mkDir(t, base, "expanded")
	canon := NewPathCanonicalizer([]string{expanded}, nil)
	if u := canon.UnusableRoots(); len(u) > 0 {
		t.Fatalf("fixture premise broken, unusable roots: %v", u)
	}
	t.Setenv(name, expanded)
	spelled := "%" + name + "%"

	t.Run("the_refusal_leg_is_untouched", func(t *testing.T) {
		// The gate ticket 102 installed still fires first: a %VAR% spelling is
		// refused as a workspace, and nothing is recorded. This sub-case is the
		// reason this ticket is not a loosening.
		res, err := canon.ResolveWorkspace(spelled)
		if err == nil {
			t.Fatalf("an expanded spelling was accepted as a workspace: %+v", res)
		}
		if !errors.Is(err, risk.ErrRewrittenPath) {
			t.Errorf("refusal %q does not carry risk.ErrRewrittenPath (res.Rewritten=%v)", err, res.Rewritten)
		}
		if !res.Rewritten {
			t.Errorf("C26 returned rewritten=false for %q, which expanded onto %q: the account this package exists to carry is missing at the source",
				spelled, res.Canonical)
		}
		if got := canon.WorkspaceRoot(); got.Canonical != "" {
			t.Errorf("a refused switch left a workspace set: %+v", got)
		}
	})

	t.Run("a_recorded_account_comes_back_as_written", func(t *testing.T) {
		// The leg ErrRewrittenPath's own sentence permits: act on the expanded
		// tree when that tree is named for itself. The host here does exactly
		// that, and hands the account along with it - which is the only way
		// panel.WorkspaceView.Rewritten can ever be a reading instead of a zero.
		res, _ := canon.ResolveWorkspace(spelled)
		if !res.Rewritten {
			t.Fatalf("premise broken: C26 substituted nothing for %q (canonical %q), so no rewritten account exists to record",
				spelled, res.Canonical)
		}
		if err := canon.SetWorkspaceRoot(res.Canonical, res); err != nil {
			t.Fatalf("SetWorkspaceRoot(%s): %v", res.Canonical, err)
		}
		got := canon.WorkspaceRoot()
		if !got.Rewritten {
			t.Errorf("rewritten = false after recording C26's own rewritten Result (%+v): "+
				"the producer is filling the field with a constant again, and 181-r2's consumer branch is dead code", got)
		}
		if got.Spelling != spelled {
			t.Errorf("spelling = %q, want the substituted spelling %q", got.Spelling, spelled)
		}
		if got.Canonical != res.Canonical {
			t.Errorf("canonical = %q, want the expanded tree %q", got.Canonical, res.Canonical)
		}
		if !strings.Contains(strings.Join(got.Rewrites, "+"), "env") {
			t.Errorf("rewrites = %v, want it to name the construct that substituted (env)", got.Rewrites)
		}

		// The answer is a copy: handing out the stored slice would let a reader
		// edit this package's own book.
		got.Rewrites[0] = "tampered"
		if again := canon.WorkspaceRoot(); strings.Join(again.Rewrites, "+") == "tampered" {
			t.Error("WorkspaceRoot handed out the stored rewrites slice; a caller can rewrite the canonicalizer's own account")
		}
	})
}

func TestSetWorkspaceRootRefusesAnAccountAboutAnotherTree(t *testing.T) {
	base := sealableTempDir124(t)
	wsA := mkDir(t, base, "alpha")
	wsB := mkDir(t, base, "beta")
	canon := NewPathCanonicalizer([]string{wsA, wsB}, nil)
	if u := canon.UnusableRoots(); len(u) > 0 {
		t.Fatalf("fixture premise broken, unusable roots: %v", u)
	}

	resA, err := canon.ResolveWorkspace(wsA)
	if err != nil {
		t.Fatalf("ResolveWorkspace(%s): %v", wsA, err)
	}
	// Both trees are authorised, so the allowlist check alone cannot catch this:
	// the account attached to alpha must not be allowed to vouch for beta.
	if err := canon.SetWorkspaceRoot(resA.Canonical, risk.Result{Canonical: wsB, Spelling: wsB}); err == nil {
		t.Fatal("SetWorkspaceRoot recorded an account that names a different tree - the packet would carry beta's book for alpha's coordinate")
	}
	if got := canon.WorkspaceRoot(); got.Canonical != "" {
		t.Errorf("the mismatched account still narrowed the scope: %+v", got)
	}
	// The honest pairing still works.
	if err := canon.SetWorkspaceRoot(resA.Canonical, resA); err != nil {
		t.Fatalf("SetWorkspaceRoot with its own account: %v", err)
	}
	if got := canon.WorkspaceRoot(); got.Canonical != resA.Canonical {
		t.Errorf("WorkspaceRoot = %+v, want alpha %q", got, resA.Canonical)
	}
}

func TestClearWorkspaceDropsTheAccountWithTheRoot(t *testing.T) {
	base := sealableTempDir124(t)
	wsA := mkDir(t, base, "alpha")
	canon := NewPathCanonicalizer([]string{wsA}, nil)
	res, err := canon.ResolveWorkspace(wsA)
	if err != nil {
		t.Fatalf("ResolveWorkspace(%s): %v", wsA, err)
	}
	if err := canon.SetWorkspaceRoot(res.Canonical, res); err != nil {
		t.Fatalf("SetWorkspaceRoot: %v", err)
	}
	canon.ClearWorkspace()
	if got := canon.WorkspaceRoot(); !reflect.DeepEqual(got, risk.Result{}) {
		t.Errorf("WorkspaceRoot after ClearWorkspace = %+v, want the zero Result: "+
			"an account about a tree that is no longer in force is a phantom coordinate in the packet", got)
	}
}
