package tools

// Ticket 92 AC#3: choosing a workspace must change what the RISK ASSESSOR sees,
// refuse anything that crossed a link or moved trees, and stay narrowable only.
//
// The (i) case is the load-bearing one: it drives a real risk.RiskAssessor with
// the real canonicalizer wired in (the same wiring cmd/wisp does), so "the new
// workspace is used" is a verdict change on screen, not a field that was set.

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// mkDir makes a real directory under t.TempDir and returns its canonical form.
func mkDir(t *testing.T, parent, name string) string {
	t.Helper()
	p := filepath.Join(parent, name)
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", name, err)
	}
	return p
}

func assessWrite(t *testing.T, canon *PathCanonicalizer, target string) risk.Decision {
	t.Helper()
	a := risk.NewRiskAssessor().WithCanonicalizer(canon)
	return a.Assess("fs.write", map[string]any{"path": target}, risk.Facts{
		Declared: risk.L1,
		Paths:    []string{target},
	})
}

func TestWorkspaceSwitchNarrowsWhatTheAssessorJudges(t *testing.T) {
	base := sealableTempDir124(t)
	wsA := mkDir(t, base, "alpha")
	wsB := mkDir(t, base, "beta")
	// Both trees are authorized by config; the file under beta is a plain L1
	// write while no workspace is chosen.
	canon := NewPathCanonicalizer([]string{wsA, wsB}, nil)
	if u := canon.UnusableRoots(); len(u) > 0 {
		t.Fatalf("fixture premise broken, unusable roots: %v", u)
	}
	before := assessWrite(t, canon, filepath.Join(wsB, "note.txt"))
	if before.Level != risk.L1 {
		t.Fatalf("premise: an in-allowlist write with no workspace chosen assessed %q (%v), want L1",
			before.Level, before.RulesHit)
	}
	if canon.WorkspaceRoot() != "" {
		t.Fatal("premise: a fresh canonicalizer already reports a workspace")
	}

	res, err := canon.ResolveWorkspace(wsA)
	if err != nil {
		t.Fatalf("ResolveWorkspace(%s) = %v, want success", wsA, err)
	}
	if err := canon.SetWorkspaceRoot(res.Canonical); err != nil {
		t.Fatalf("SetWorkspaceRoot: %v", err)
	}

	after := assessWrite(t, canon, filepath.Join(wsB, "note.txt"))
	if after.Level != risk.L2 {
		t.Errorf("AC#3(i): after switching the workspace to %s, a write under %s still assessed %q (rules %v) - "+
			"the decision chain is not reading the new workspace", wsA, wsB, after.Level, after.RulesHit)
	}
	if !containsRule(after.RulesHit, "R2") {
		t.Errorf("AC#3(i): expected R2 (out of scope) to fire after the switch, rules=%v reason=%q",
			after.RulesHit, after.Reason)
	}
	inside := assessWrite(t, canon, filepath.Join(wsA, "note.txt"))
	if inside.Level != risk.L1 {
		t.Errorf("AC#3(i): a write INSIDE the new workspace assessed %q, want L1 - narrowing must not be a blanket deny",
			inside.Level)
	}
	canon.ClearWorkspace()
	if back := assessWrite(t, canon, filepath.Join(wsB, "note.txt")); back.Level != risk.L1 {
		t.Errorf("clearing the workspace did not restore the config-root judgement: %q", back.Level)
	}
	t.Logf("AC#3(i): %s write L1 -> L2 after switching to %s; inside the workspace still L1", wsB, wsA)
}

func containsRule(rules []risk.RuleID, want string) bool {
	for _, r := range rules {
		if string(r) == want {
			return true
		}
	}
	return false
}

func TestWorkspaceSwitchRefusesOutOfScopeAndMissingPaths(t *testing.T) {
	base := t.TempDir()
	inside := mkDir(t, base, "inside")
	outside := mkDir(t, t.TempDir(), "outside")
	canon := NewPathCanonicalizer([]string{inside}, nil)

	if _, err := canon.ResolveWorkspace(outside); err == nil {
		t.Errorf("a workspace outside [fs] allowed_dirs (%s) was accepted - a panel picker would have just widened authority", outside)
	} else if !strings.Contains(err.Error(), "allowed_dirs") {
		t.Errorf("out-of-scope refusal said %q, want it to name the allowlist", err)
	}
	if canon.WorkspaceRoot() != "" {
		t.Errorf("a refused switch left a workspace set: %q", canon.WorkspaceRoot())
	}

	if _, err := canon.ResolveWorkspace(filepath.Join(inside, "does-not-exist")); err == nil {
		t.Error("a nonexistent directory was accepted as a workspace")
	}
	if _, err := canon.ResolveWorkspace("   "); err == nil {
		t.Error("a blank workspace path was accepted")
	}
	// An empty allowlist authorizes nothing, so no workspace may be set either.
	none := NewPathCanonicalizer(nil, nil)
	if _, err := none.ResolveWorkspace(inside); err == nil {
		t.Error("a workspace was accepted with an empty allowlist")
	}
	if err := none.SetWorkspaceRoot(inside); err == nil {
		t.Error("SetWorkspaceRoot accepted an unauthorized root with an empty allowlist")
	}
	// Handing SetWorkspaceRoot a raw path that was never resolved is the other
	// way a caller can get it wrong.
	if err := canon.SetWorkspaceRoot(outside); err == nil {
		t.Error("SetWorkspaceRoot accepted a path outside the roots without resolving it")
	}
}

// TestWorkspaceSwitchRefusesAnExpandedSpelling is ticket 102's account on the
// authorization leg, and it is portable: Rewritten is computed by the shared
// part of C26 (expandAccounted), not by the Windows-only handle half.
//
// Each case puts the EXPANDED tree into [fs] allowed_dirs first, so the only
// reason a refusal can come from is the account: the spelling the panel sent is
// not the tree its own string names.
func TestWorkspaceSwitchRefusesAnExpandedSpelling(t *testing.T) {
	t.Run("leading tilde", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("no home directory to expand against on this host: %v", err)
		}
		canon := NewPathCanonicalizer([]string{home}, nil)
		if u := canon.UnusableRoots(); len(u) > 0 {
			t.Skipf("the probe root itself did not canonicalize here: %v", u)
		}
		// The spelling is the SAME tree the config authorized, so a refusal can
		// only come from the account: the panel named "~", the tree it authorizes
		// is somebody's home directory.
		res, err := canon.ResolveWorkspace(`~`)
		requireRewritten(t, res, err, canon)
	})
	t.Run("environment variable", func(t *testing.T) {
		for _, name := range []string{"TEMP", "TMPDIR", "HOME"} {
			v := os.Getenv(name)
			if v == "" {
				continue
			}
			canon := NewPathCanonicalizer([]string{v}, nil)
			if u := canon.UnusableRoots(); len(u) > 0 {
				continue
			}
			res, err := canon.ResolveWorkspace("%" + name + "%")
			if err == nil && !res.Rewritten {
				// POSIX ignores %NAME%; only a substitution counts as this case.
				continue
			}
			requireRewritten(t, res, err, canon)
			return
		}
		t.Skip("no expandable environment variable naming an existing directory on this host")
	})
}

func requireRewritten(t *testing.T, res risk.Result, err error, canon *PathCanonicalizer) {
	t.Helper()
	if err == nil {
		t.Fatalf("an expanded spelling was accepted as a workspace: %+v - the panel would have authorized %s while telling the user it authorized the string they typed", res, res.Canonical)
	}
	if !errors.Is(err, risk.ErrRewrittenPath) {
		t.Errorf("refusal %q does not carry risk.ErrRewrittenPath (res.Rewritten=%v) - the account was not read", err, res.Rewritten)
	}
	if canon.WorkspaceRoot() != "" {
		t.Errorf("a refused switch left a workspace set: %q", canon.WorkspaceRoot())
	}
	t.Logf("expanded spelling refused as required: %v", err)
}

// TestWorkspaceSwitchRefusesAJunctionToOutside is AC#3(ii) with a REAL junction.
// Windows-only because the detection itself is (risk.pathresolver_other.go's own
// DEFERRED stub): the skip is reported, never averaged away.
func TestWorkspaceSwitchRefusesAJunctionToOutside(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("C26's reparse detection is a Windows implementation (risk.pathresolver_other.go reparseComponents returns nil elsewhere); nothing to deny on %s", runtime.GOOS)
	}
	outside := mkDir(t, t.TempDir(), "elsewhere")
	base := t.TempDir()
	link := filepath.Join(base, "link")
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, outside).CombinedOutput()
	if err != nil {
		t.Skipf("mklink /J unavailable here: %v: %s", err, out)
	}
	t.Cleanup(func() { _ = os.Remove(link) })

	canon := NewPathCanonicalizer([]string{base}, nil)
	_, err = canon.ResolveWorkspace(link)
	if err == nil {
		t.Fatalf("a junction into %s was accepted as a workspace - C26's reparse denial did not reach the composer path", outside)
	}
	if !errors.Is(err, risk.ErrReparseDenied) {
		t.Errorf("refusal %q is not risk.ErrReparseDenied; the reason shown to the user must be C26's own", err)
	}
	if canon.WorkspaceRoot() != "" {
		t.Errorf("a refused switch left a workspace set: %q", canon.WorkspaceRoot())
	}
	t.Logf("junction workspace refused with C26's reason: %v", err)
}
