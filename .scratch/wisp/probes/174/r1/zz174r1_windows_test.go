//go:build windows && wisp174r1after

package main

// PROBE 174-r1, AFTER arm. Build-tagged (`-tags wisp174r1after`) so the BEFORE
// readings in zz174r1_before_windows_test.go stay reproducible on the pre-fix
// tree, where the field this file passes - tools.TaskDeps.Paths - does not
// exist. Self-contained apart from the BEFORE file's logging/formatting
// helpers, so the two arms share one fixture shape and one log file.
//
// Same legs as the BEFORE arm:
//
//	A-outside  authorization leg: the pointer leaves every authorized root
//	           (empty allowlist, i.e. [fs] allowed_dirs = [] today).
//	B-ghost    existence leg: the named file is not there.
//	B-dir      existence leg, stronger: the name is a directory.
//	POS-ctrl   the anti-noise control: inside the root AND a real regular
//	           file - it must stay silent after the fix exactly as before.
//
// Re-running the BEFORE arm after the fix is itself the fail-closed reading
// (TaskDeps with no judge = what cmd/wisp/run.go:362 wires today).

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// afterFixtures lands one REAL spill (agent.Spiller, the D15(3) code path)
// under <dataDir>\artifacts and returns the four pointer spellings the legs
// need: the real spill, a ghost inside the same root, the root itself, and a
// real regular file inside the root.
type afterFixtures struct {
	dataDir, artifacts, spill, ghost, rootDir, realFile string
	body                                                string
}

func buildAfterFixtures(t *testing.T) afterFixtures {
	t.Helper()
	body := bigOutput()
	dataDir := t.TempDir()
	artifacts := filepath.Join(dataDir, "artifacts") // run.go:609 shape
	sp := agent.NewSpiller(artifacts, agent.BudgetsFor(0))
	st, err := sp.Prepare("call-1", body)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Spilled || st.Path == "" {
		t.Fatalf("expected a real spill, got Spilled=%v Path=%q", st.Spilled, st.Path)
	}
	realFile := filepath.Join(artifacts, "inside-root-real-file.txt")
	if err := os.WriteFile(realFile, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return afterFixtures{
		dataDir: dataDir, artifacts: artifacts, spill: st.Path,
		ghost:   filepath.Join(artifacts, "tool-output-does-not-exist.txt"),
		rootDir: artifacts, realFile: realFile, body: body,
	}
}

// taskBridgeWithJudge wires fs.* AND task.* over the SAME C26 canonicalizer,
// which is the point of this ticket's fix: the tool judges the pointer with
// the very code fs.read will later be judged by.
func taskBridgeWithJudge(t *testing.T, roster *tools.TaskRoster, judge *tools.PathCanonicalizer) *tools.Bridge {
	t.Helper()
	reg := tools.NewRegistry()
	for _, e := range append(tools.BuiltinFSEntries(tools.FSDeps{Paths: judge}),
		tools.BuiltinTaskEntries(tools.TaskDeps{Roster: roster, Paths: judge})...) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return tools.New(tools.Options{
		Registry: reg, Paths: judge,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		Gate:       tools.NoGate{}, Logf: func(string, ...any) {},
	})
}

func TestProbe174R1PointerHonestyAfter(t *testing.T) {
	lab := label174r1()
	probeLog(t, "==== 174-r1 AFTER arm (judge wired), label=%s ====", lab)
	fx := buildAfterFixtures(t)

	// A-outside: default allowlist (empty), exactly run.go:327-331.
	emptyJudge := tools.NewPathCanonicalizer(nil, nil)
	rosterA := tools.NewTaskRoster()
	rosterA.Record("A-outside", tools.TaskOutput{Text: fx.body, ArtifactPath: fx.spill})
	bA := taskBridgeWithJudge(t, rosterA, emptyJudge)
	outA := taskOut(t, bA, "A-outside")
	probeLog(t, "[after A-outside] dataDir=%s artifacts=%s spill=%s", fx.dataDir, fx.artifacts, fx.spill)
	probeLog(t, "[after A-outside] isError=%v truncated=%v announcement=%q",
		outA.IsError, outA.Truncated, announceLine(outA.Text))
	probeLog(t, "[after A-outside] %s", speakFlags(outA.Text))
	rdA, err := bA.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-r1", CallID: "read-A", Name: "fs.read",
		Args: mustJSON174(t, map[string]any{"path": filepath.ToSlash(fx.spill)}),
	})
	if err != nil {
		t.Fatal(err)
	}
	probeLog(t, "[after A-outside] fs.read of that same path: isError=%v level=%v class=%q",
		rdA.IsError, rdA.RiskLevel, rdA.ErrorClass)

	// B-ghost / B-dir / POS-ctrl: authorized root = the artifacts dir, so only
	// existence/shape differs between the three.
	rootJudge := tools.NewPathCanonicalizer([]string{fx.artifacts}, nil)
	rosterB := tools.NewTaskRoster()
	rosterB.Record("B-ghost", tools.TaskOutput{Text: fx.body, ArtifactPath: fx.ghost})
	rosterB.Record("B-dir", tools.TaskOutput{Text: fx.body, ArtifactPath: fx.rootDir})
	rosterB.Record("POS-ctrl", tools.TaskOutput{Text: fx.body, ArtifactPath: fx.realFile})
	bB := taskBridgeWithJudge(t, rosterB, rootJudge)
	for _, id := range []string{"B-ghost", "B-dir", "POS-ctrl"} {
		out := taskOut(t, bB, id)
		probeLog(t, "[after %s] announcement=%q", id, announceLine(out.Text))
		probeLog(t, "[after %s] %s", id, speakFlags(out.Text))
	}
	rdB, err := bB.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-r1", CallID: "read-POS", Name: "fs.read",
		Args: mustJSON174(t, map[string]any{"path": filepath.ToSlash(fx.realFile)}),
	})
	if err != nil {
		t.Fatal(err)
	}
	probeLog(t, "[after POS-ctrl] fs.read of the control path: isError=%v level=%v bytes=%d (whole file is %d)",
		rdB.IsError, rdB.RiskLevel, len(rdB.Text), len(fx.body))
	probeLog(t, "==== 174-r1 AFTER arm, label=%s done ====", lab)
}
