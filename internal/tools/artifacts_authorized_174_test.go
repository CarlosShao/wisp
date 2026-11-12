package tools

// Ticket 174 AC#4 - the artifacts tree joins the [fs] authorization roots.
//
// The defect this pins: the agent's own spill writes into <dataDir>\artifacts,
// and task.output then points the model at that file. Until this change nothing
// authorized that tree, so fs.read answered "outside" for a file this very
// process had just created - D14's "read on demand" was an empty promise. The
// owner approved adding the artifacts tree to the authorization roots
// (census-199 §4 row 14 asks exactly this; A1571's B bucket covers 174 and
// approves the contract-wording step), and the wording needs no edit: PLAN.md's
// promise becomes true because of the assembly line, not because of a text
// change.
//
// Two directions, because one alone is satisfied by a change that authorizes
// everything:
//   - forward: a spilled artifact IS readable through the same C26 instance,
//     at L0;
//   - reverse: the artifacts directory's SIBLINGS (secrets\, config.toml, the
//     db) stay outside - a root that quietly widened to the whole data root
//     would pass the forward direction alone.
//
// The root set here is built the way the production assembly builds it
// (mem.ArtifactsDir() appended to the configured roots), because a seam test
// that hands itself a different spelling would not be measuring the change.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// t174DataRoot lays out the data root the production assembly has by the time
// it builds the canonicalizer: the sealed data dir, the artifacts tree inside
// it (created and sealed by memory.Open at open.go:283), and the siblings that
// must stay outside any artifact grant.
func t174DataRoot(t *testing.T) (dataDir, artifactsDir string) {
	t.Helper()
	dataDir = sealableTempDir124(t)
	artifactsDir = filepath.Join(dataDir, "artifacts")
	if err := os.MkdirAll(artifactsDir, 0o700); err != nil {
		t.Fatalf("create artifacts dir: %v", err)
	}
	// The siblings: a secrets tree and the config file. Neither may become
	// readable just because the artifacts tree did.
	if err := os.MkdirAll(filepath.Join(dataDir, "secrets"), 0o700); err != nil {
		t.Fatalf("create secrets dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "secrets", "token.txt"), []byte("s3cr3t"), 0o600); err != nil {
		t.Fatalf("write secret: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "config.toml"), []byte("[fs]\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return dataDir, artifactsDir
}

// t174Roots is the production root set: the configured roots plus the artifacts
// tree, in exactly the shape cmd/wisp/run.go builds it.
func t174Roots(configured []string, artifactsDir string) []string {
	allowed := make([]string, 0, len(configured)+1)
	allowed = append(allowed, configured...)
	allowed = append(allowed, artifactsDir)
	return allowed
}

func Test174SpilledArtifactIsReadableThroughTheSameCanonicalizer(t *testing.T) {
	dataDir, artifactsDir := t174DataRoot(t)

	// The production shape: a workspace root the operator configured, plus the
	// artifacts tree the assembly appends.
	workspace := filepath.Join(dataDir, "workspace")
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	paths := NewPathCanonicalizer(t174Roots([]string{workspace}, artifactsDir), nil)

	// The root really is in the set, in the same spelling the spill target uses.
	roots := paths.Roots()
	found := false
	for _, r := range roots {
		if strings.EqualFold(filepath.Clean(r), filepath.Clean(artifactsDir)) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("artifacts root missing from Roots()=%v (wanted %q); an unusable root is silently dropped, which would leave the ticket unfixed while every test passed", roots, artifactsDir)
	}

	// Forward: a real spill into that tree, then fs.read over the same C26.
	spills := agent.NewSpiller(artifactsDir, agent.BudgetsFor(128000))
	spill, err := spills.Prepare("call-174-a", strings.Repeat("artifact-body-174 ", 1200))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if !spill.Spilled || spill.Path == "" {
		t.Fatalf("Prepare did not spill (Spilled=%v Path=%q); the forward direction needs a real artifact", spill.Spilled, spill.Path)
	}

	b := t174Bridge(t, paths)
	out := r183Call(t, b, "task-174-a", "fs.read", r183Args(t, spill.Path, 0))
	if out.IsError {
		t.Fatalf("fs.read of the artifact this process just spilled was refused: %s", out.Text)
	}
	if out.RiskLevel != "L0" {
		t.Fatalf("fs.read of an authorized artifact landed at %q, want L0", out.RiskLevel)
	}
	if !strings.Contains(out.Text, "artifact-body-174") {
		t.Fatalf("fs.read did not return the artifact body; got %q", truncateFor174(out.Text, 200))
	}
}

func Test174DataRootSiblingsStayOutsideTheArtifactGrant(t *testing.T) {
	dataDir, artifactsDir := t174DataRoot(t)
	workspace := filepath.Join(dataDir, "workspace")
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	paths := NewPathCanonicalizer(t174Roots([]string{workspace}, artifactsDir), nil)
	b := t174Bridge(t, paths)

	// Reverse: every sibling of the artifacts tree is still outside. A change
	// that appended the whole data root instead of the artifacts tree would
	// make all three readable, and only this direction notices.
	for _, sibling := range []string{
		filepath.Join(dataDir, "secrets", "token.txt"),
		filepath.Join(dataDir, "config.toml"),
		filepath.Join(dataDir, "wisp.db"),
	} {
		out := r183Call(t, b, "task-174-b", "fs.read", r183Args(t, sibling, 0))
		if !out.IsError {
			t.Fatalf("fs.read of %q was allowed; the artifact grant must not widen to the data root", sibling)
		}
		if out.RiskLevel != "L2" {
			t.Fatalf("fs.read of %q landed at %q, want L2 (R2)", sibling, out.RiskLevel)
		}
	}
}

// Test174NoArtifactRootMeansTheArtifactIsUnreadable is the mutation control: the
// same spilled artifact, read through a canonicalizer built WITHOUT the
// artifacts root - the pre-change assembly. It must be refused, or the forward
// direction above proves nothing about this change.
func Test174NoArtifactRootMeansTheArtifactIsUnreadable(t *testing.T) {
	dataDir, artifactsDir := t174DataRoot(t)
	workspace := filepath.Join(dataDir, "workspace")
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	spills := agent.NewSpiller(artifactsDir, agent.BudgetsFor(128000))
	spill, err := spills.Prepare("call-174-c", strings.Repeat("artifact-body-174 ", 1200))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if !spill.Spilled {
		t.Fatalf("Prepare did not spill; the control needs a real artifact")
	}

	// The pre-change assembly: the configured roots only.
	paths := NewPathCanonicalizer([]string{workspace}, nil)
	b := t174Bridge(t, paths)
	out := r183Call(t, b, "task-174-c", "fs.read", r183Args(t, spill.Path, 0))
	if !out.IsError {
		t.Fatalf("without the artifacts root the artifact read was ALLOWED - the forward direction above is not measuring the change; got %q", truncateFor174(out.Text, 200))
	}
	if out.RiskLevel != "L2" {
		t.Fatalf("without the artifacts root the read landed at %q, want L2", out.RiskLevel)
	}
}

// t174Bridge runs fs.read over the real bridge with the given canonicalizer.
func t174Bridge(t *testing.T, paths *PathCanonicalizer) *Bridge {
	t.Helper()
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil})
	return New(Options{
		Registry:   reg,
		Paths:      paths,
		Provenance: prov,
		Gate:       NoGate{},
		Logf:       func(string, ...any) {},
	})
}

func truncateFor174(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
