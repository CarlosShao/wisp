// Probe fixture for the read-only census leg 174-c2 (cost of wiring the path
// judge into production; this leg changes NO tracked file).
//
// It compiles into package tools through `go test -overlay` (see overlay.json
// next to it), so internal/tools/ stays byte-clean.
//
// Purpose:票 174 AC#2b's two-way reading, measured on ONE tree. The only thing
// AC#2b asks for is `Paths:` on the `tools.TaskDeps{...}` literal in
// cmd/wisp/run.go. This probe therefore builds BOTH arms side by side:
//
//	before = tools.TaskDeps{Roster: roster}            (what run.go wires today)
//	after  = tools.TaskDeps{Roster: roster, Paths: J}  (the one added field)
//
// fs.* is registered in both arms over the SAME canonicalizer J, so the reply
// text and the fs.read verdict for one path come out of one and the same C26
// instance - which is what lets this probe say whether wiring moves authority
// or only moves words.
//
// It reads nothing into the tree: all output goes to stdout, which the runner
// redirects into this probe's own log directory.
package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

const c2BodyBytes = 30000

// c2Pair builds one arm. givePaths is the single difference between the arms.
func c2Pair(t *testing.T, roster *TaskRoster, judge *PathCanonicalizer, givePaths bool) *Bridge {
	t.Helper()
	reg := NewRegistry()
	td := TaskDeps{Roster: roster}
	if givePaths {
		td.Paths = judge
	}
	for _, e := range append(
		BuiltinFSEntries(FSDeps{Paths: judge}),
		BuiltinTaskEntries(td)...) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry:   reg,
		Paths:      judge,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		Gate:       NoGate{},
		Logf:       func(string, ...any) {},
	})
}

func c2Call(t *testing.T, b *Bridge, name string, args map[string]any) agent.ToolOutcome {
	t.Helper()
	raw, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}
	res, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-c2", CallID: "c2-" + name, Name: name, Args: raw,
	})
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return res
}

// c2Announce returns only the bracketed pointer line, so a log line stays short.
func c2Announce(text string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "[…") {
			return line
		}
	}
	return text
}

func c2Fixtures(t *testing.T) (artifacts, spill, ghost, dirShape, outside string) {
	t.Helper()
	dataDir := t.TempDir()
	artifacts = filepath.Join(dataDir, "artifacts") // cmd/wisp/run.go:612 shape
	sp := agent.NewSpiller(artifacts, agent.BudgetsFor(0))
	st, err := sp.Prepare("tool-output-c2", strings.Repeat("a", c2BodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	if !st.Spilled || st.Path == "" {
		t.Fatalf("no real spill: Spilled=%v Path=%q", st.Spilled, st.Path)
	}
	spill = st.Path
	ghost = filepath.Join(artifacts, "tool-output-c2-ghost.txt")
	dirShape = artifacts
	outsideDir := t.TempDir()
	outside = filepath.Join(outsideDir, "c2-outside-real-file.txt")
	if err := os.WriteFile(outside, []byte(strings.Repeat("b", c2BodyBytes)), 0o600); err != nil {
		t.Fatal(err)
	}
	return
}

func TestProbe174C2WiringBeforeAfterPair(t *testing.T) {
	artifacts, spill, ghost, dirShape, outside := c2Fixtures(t)
	body := strings.Repeat("a", c2BodyBytes)
	t.Logf("FIXTURES artifacts=%s spill=%s", filepath.ToSlash(artifacts), filepath.ToSlash(spill))

	for _, cfg := range []struct {
		label string
		roots []string
	}{
		{"default-allowed-dirs-empty", nil}, // [fs] allowed_dirs = [] (SPEC-03:35)
		{"root-is-artifacts-dir", []string{artifacts}},
	} {
		judge := NewPathCanonicalizer(cfg.roots, nil)
		t.Logf("==== arm-set %s (roots=%v, canonical roots=%v) ====",
			cfg.label, cfg.roots, judge.Roots())

		for _, shape := range []struct{ id, path string }{
			{"real-spill", spill},
			{"ghost-missing", ghost},
			{"dir-not-a-file", dirShape},
			{"outside-root-real-file", outside},
		} {
			roster := NewTaskRoster()
			roster.Record(shape.id, TaskOutput{Text: body, ArtifactPath: shape.path})

			before := c2Pair(t, roster, judge, false)
			after := c2Pair(t, roster, judge, true)

			ob := c2Call(t, before, "task.output", map[string]any{"task_id": shape.id})
			oa := c2Call(t, after, "task.output", map[string]any{"task_id": shape.id})
			probeRead := func(b *Bridge, tag string) {
				r := c2Call(t, b, "fs.read", map[string]any{
					"path": filepath.ToSlash(shape.path)})
				t.Logf("  [%s/%s] fs.read  -> isError=%v level=%v class=%q bytes=%d",
					tag, shape.id, r.IsError, r.RiskLevel, r.ErrorClass, len(r.Text))
			}
			probeRead(before, "before")
			probeRead(after, "after")

			t.Logf("  [%s/%s] BEFORE task.output 注意：x%d  %q",
				cfg.label, shape.id, strings.Count(ob.Text, "注意："), c2Announce(ob.Text))
			t.Logf("  [%s/%s] AFTER  task.output 注意：x%d  %q",
				cfg.label, shape.id, strings.Count(oa.Text, "注意："), c2Announce(oa.Text))
			t.Logf("  [%s/%s] DID THE POINTER LINE LOSE 全文见? before=%v after=%v",
				cfg.label, shape.id, strings.Contains(ob.Text, "全文见 "), strings.Contains(oa.Text, "全文见 "))
		}
	}
}
