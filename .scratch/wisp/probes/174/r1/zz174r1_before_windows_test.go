//go:build windows

package main

// PROBE 174-r1, BEFORE arm (measurement only; zero production code touched).
//
// It wires EXACTLY what the pre-fix tree offers - tools.TaskDeps{Roster: ...}
// and nothing else, i.e. what cmd/wisp/run.go:362 builds today - so the
// readings below are reproducible on the unmodified commit. Legs (seam level:
// real tools.Bridge + real task.output + real C26 canonicalizer + real
// fs.read inside `go test`):
//
//	A-outside : a real agent.Spiller artifact under <dataDir>\artifacts while
//	          the allowlist is the default empty list (run.go:327-331 shape).
//	          The pointer names a path fs.read cannot open.
//	B-ghost   : ArtifactPath names a file that does not exist, INSIDE the
//	          authorized root, so shape (b) is isolated from shape (a).
//	B-dir     : ArtifactPath names a directory - the shape ticket 164's own
//	          TestTruncationShapeIsTheD15Triple sits on.
//	POS-ctrl  : a path inside the authorized root that IS a real regular file.
//	          This leg must stay silent both before and after the fix.
//
// Readings append to ./logs/ next to this file (runtime.Caller, never CWD).
// Tag a run with WISP174R1_LABEL=before|after. The AFTER arm lives in
// zz174r1_after_windows_test.go, which additionally wires a path judge.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

const spillBody = 51633 // the byte count probe 174-c1 measured on a real spill

func label174r1() string {
	if l := os.Getenv("WISP174R1_LABEL"); l != "" {
		return l
	}
	return "unlabeled"
}

func probeLog(t *testing.T, format string, args ...any) {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Join(filepath.Dir(file), "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	line := fmt.Sprintf(format, args...)
	t.Log(line)
	f, err := os.OpenFile(filepath.Join(dir, "readings.txt"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fmt.Fprintln(f, line)
}

// bigOutput is deterministic, so the before/after runs compare bytes.
func bigOutput() string {
	var b strings.Builder
	b.WriteString("W174R1-HEAD-MARK\n")
	for i := 0; b.Len() < spillBody; i++ {
		fmt.Fprintf(&b, "line %05d: padding padding padding padding padding\n", i)
	}
	return b.String()[:spillBody]
}

// bridgeTaskFS wires fs.* (judged by judge) plus task.* the way the host wires
// them today: the task tool gets the roster and NOTHING else.
func bridgeTaskFS(t *testing.T, roster *tools.TaskRoster, judge *tools.PathCanonicalizer) *tools.Bridge {
	t.Helper()
	reg := tools.NewRegistry()
	for _, e := range append(tools.BuiltinFSEntries(tools.FSDeps{Paths: judge}),
		tools.BuiltinTaskEntries(tools.TaskDeps{Roster: roster})...) {
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

// taskOut calls task.output through the real bridge and returns the model-facing text.
func taskOut(t *testing.T, b *tools.Bridge, id string) agent.ToolOutcome {
	t.Helper()
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-r1", CallID: "c-" + id, Name: "task.output",
		Args: mustJSON174(t, map[string]any{"task_id": id}),
	})
	if err != nil {
		t.Fatalf("task.output Go error (a reply is data, not a fault): %v", err)
	}
	return out
}

// speakFlags reports whether the model-visible text says anything about the
// pointer being unreadable. "不响" is all three false.
func speakFlags(stub string) string {
	return fmt.Sprintf("提到读不到=%v 提到不存在=%v 提到不是一般文件=%v 提到未接线=%v 全文见=%v 不可找回=%v",
		strings.Contains(stub, "读不到"), strings.Contains(stub, "不存在"),
		strings.Contains(stub, "不是一般文件"), strings.Contains(stub, "未接线"),
		strings.Contains(stub, "全文见"), strings.Contains(stub, "不可找回"))
}

// announceLine pulls the bracketed pointer sentence out of a stub - the part
// the model reads about the pointer.
func announceLine(stub string) string {
	i := strings.Index(stub, "[…")
	if i < 0 {
		return "(no announcement)"
	}
	j := strings.Index(stub[i:], "…]")
	if j < 0 {
		return stub[i:]
	}
	return stub[i : i+j+len("…]")]
}

func mustJSON174(t *testing.T, m map[string]any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestProbe174R1PointerHonestyBefore(t *testing.T) {
	lab := label174r1()
	probeLog(t, "==== 174-r1 BEFORE arm, label=%s ====", lab)

	full := bigOutput()

	// A REAL spill under <dataDir>\artifacts (run.go:609 shape).
	dataDir := t.TempDir()
	artifacts := filepath.Join(dataDir, "artifacts")
	sp := agent.NewSpiller(artifacts, agent.BudgetsFor(0))
	st, err := sp.Prepare("call-1", full)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Spilled || st.Path == "" {
		t.Fatalf("expected a real spill, got Spilled=%v Path=%q", st.Spilled, st.Path)
	}

	ghost := filepath.Join(artifacts, "tool-output-does-not-exist.txt")
	inside := filepath.Join(artifacts, "inside-root-real-file.txt")
	if err := os.WriteFile(inside, []byte(full), 0o600); err != nil {
		t.Fatal(err)
	}

	// A-outside: default allowlist (empty) - the pointer leaves every root.
	emptyJudge := tools.NewPathCanonicalizer(nil, nil)
	rosterA := tools.NewTaskRoster()
	rosterA.Record("A-outside", tools.TaskOutput{Text: full, ArtifactPath: st.Path})
	bA := bridgeTaskFS(t, rosterA, emptyJudge)
	outA := taskOut(t, bA, "A-outside")
	probeLog(t, "[A-outside] spill path=%s", st.Path)
	probeLog(t, "[A-outside] isError=%v truncated=%v announcement=%q", outA.IsError, outA.Truncated, announceLine(outA.Text))
	probeLog(t, "[A-outside] %s", speakFlags(outA.Text))
	rd, err := bA.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-r1", CallID: "read-A", Name: "fs.read",
		Args: mustJSON174(t, map[string]any{"path": filepath.ToSlash(st.Path)}),
	})
	if err != nil {
		t.Fatal(err)
	}
	probeLog(t, "[A-outside] the fact behind the words, fs.read on that path: isError=%v level=%v class=%q text=%q",
		rd.IsError, rd.RiskLevel, rd.ErrorClass, rd.Text)

	// B-ghost / B-dir / POS-ctrl: authorized root = the artifacts dir, so only
	// existence differs between the three.
	rootJudge := tools.NewPathCanonicalizer([]string{artifacts}, nil)
	rosterB := tools.NewTaskRoster()
	rosterB.Record("B-ghost", tools.TaskOutput{Text: full, ArtifactPath: ghost})
	rosterB.Record("B-dir", tools.TaskOutput{Text: full, ArtifactPath: artifacts})
	rosterB.Record("POS-ctrl", tools.TaskOutput{Text: full, ArtifactPath: inside})
	bB := bridgeTaskFS(t, rosterB, rootJudge)
	for _, id := range []string{"B-ghost", "B-dir", "POS-ctrl"} {
		out := taskOut(t, bB, id)
		probeLog(t, "[%s] path-said=%q", id, announceLine(out.Text))
		probeLog(t, "[%s] %s", id, speakFlags(out.Text))
	}
	rdB, err := bB.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-r1", CallID: "read-POS", Name: "fs.read",
		Args: mustJSON174(t, map[string]any{"path": filepath.ToSlash(inside)}),
	})
	if err != nil {
		t.Fatal(err)
	}
	probeLog(t, "[POS-ctrl] fs.read of the control path: isError=%v level=%v bytes=%d (whole file is %d)",
		rdB.IsError, rdB.RiskLevel, len(rdB.Text), len(full))
	probeLog(t, "==== 174-r1 BEFORE arm, label=%s done ====", lab)
}
