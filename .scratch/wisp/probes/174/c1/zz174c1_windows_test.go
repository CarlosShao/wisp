//go:build windows

package main

// PROBE 174-c1 (READ-ONLY measurement, not a delivered test; zero production
// code touched). It measures, at the SEAM layer only (real tools.Bridge +
// real fs.read + real risk assessor + real agent.Spiller inside `go test`),
// the two directions of ticket 174 AC#1 and the fake-pointer shape of AC#2(b):
//
//	dir 1: default config (allowed_dirs = [], exactly the run.go:327-331
//	       shape) - a REAL spill lands in <dataDir>\artifacts (run.go:609
//	       shape), then fs.read is called on the path from the stub.
//	dir 2: the same artifact dir configured INTO the allowlist - the
//	       control leg proving dir 1's refusal is the allowlist's doing,
//	       not a misbuilt harness.
//	fake : task.output roster records whose ArtifactPath is a nonexistent
//	       file / a directory - does the "全文见 …" pointer branch fire and
//	       stay silent (branch = task.go:227, `rec.ArtifactPath != ""`,
//	       no stat)?
//
// Readings land in ./logs/ next to this file (runtime.Caller, never CWD).

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
	"github.com/CarlosShao/wisp/internal/tools"
)

var logOnce bool

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

// captured holds the L2 cards the bridge actually showed and every final
// Decision it emitted, so tier/rules/reason are read verbatim from the same
// objects the confirmation card would render (OnDecision contract: the card
// shows RulesHit/Reason VERBATIM).
type captured struct {
	decisions []tools.Decision
	l2Asks    int
}

func buildBridge(allowed []string, cap *captured, roster *tools.TaskRoster,
	l2Answer tools.Answer, l2Why string,
) *tools.Bridge {
	pc := tools.NewPathCanonicalizer(allowed, nil) // same call as run.go:331
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{Paths: pc, DeleteEnabled: false}) {
		if err := reg.Register(e); err != nil {
			panic(err)
		}
	}
	if roster != nil {
		for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: roster}) {
			if err := reg.Register(e); err != nil {
				panic(err)
			}
		}
	}
	gate := tools.GateFuncs{
		Window: func(context.Context, tools.Decision) (tools.Answer, string) {
			return tools.AnswerAllow, ""
		},
		Approval: func(_ context.Context, d tools.Decision) (tools.Answer, string) {
			cap.l2Asks++
			cp := d
			cap.decisions = append(cap.decisions, cp)
			return l2Answer, l2Why // harness plays the human/timeout leg
		},
	}
	return tools.New(tools.Options{
		Registry: reg,
		Paths:    pc,
		Gate:     gate,
	})
}

func exec(t *testing.T, b *tools.Bridge, name string, params map[string]any) agent.ToolOutcome {
	t.Helper()
	args, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "p174-c1", CallID: "call-1", Name: name, Args: args,
	})
	if err != nil {
		t.Fatalf("Execute(%s) returned Go error: %v", name, err)
	}
	return out
}

// bigOutput is ~40KB, i.e. ~10000 approx-tokens: over both frozen 4000-token
// spill thresholds and under the 1MiB raw cap and the 256KiB fs.read cap.
func bigOutput() string {
	var sb strings.Builder
	sb.WriteString("W174C1-HEAD-MARK")
	for i := 0; i < 1200; i++ {
		fmt.Fprintf(&sb, "\nline %04d: padding padding padding padding", i)
	}
	sb.WriteString("\nW174C1-TAIL-MARK")
	return sb.String()
}

// realSpill drives the PRODUCTION agent.Spiller (the D15(3) code path) so the
// artifact on disk is a real spill, not hand-stubbed text.
func realSpill(t *testing.T) (dataDir, artifacts string, st agent.Spill, body string) {
	t.Helper()
	dataDir = t.TempDir()
	artifacts = filepath.Join(dataDir, "artifacts") // run.go:609 shape, verbatim join
	body = bigOutput()
	sp := agent.NewSpiller(artifacts, agent.BudgetsFor(0))
	st, err := sp.Prepare("call-1", body)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Spilled || st.Path == "" {
		t.Fatalf("expected a real spill, got Spilled=%v Path=%q", st.Spilled, st.Path)
	}
	raw, err := os.ReadFile(st.Path)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != len(body) {
		t.Fatalf("artifact on disk is %d bytes, body is %d", len(raw), len(body))
	}
	return
}

func dumpDecision(t *testing.T, tag string, d tools.Decision) {
	t.Helper()
	probeLog(t, "[%s] decision: tool=%s level=%v rules=%v column=%v reason=%q",
		tag, d.Tool, d.Level, d.RulesHit, d.DecisionColumn, d.Reason)
}

func TestDir1DefaultAllowlist(t *testing.T) {
	_, _, st, _ := realSpill(t)

	for _, tc := range []struct {
		name string
		ans  tools.Answer
		why  string
	}{
		{"operator-rejects", tools.AnswerReject, ""},
		{"card-times-out", tools.AnswerTimeout, ""},
	} {
		cap := &captured{}
		b := buildBridge(nil, cap, nil, tc.ans, tc.why)
		out := exec(t, b, "fs.read", map[string]any{"path": st.Path})
		probeLog(t, "[dir1/%s] spill path = %s", tc.name, st.Path)
		probeLog(t, "[dir1/%s] gate L2 asks = %d", tc.name, cap.l2Asks)
		for _, d := range cap.decisions {
			dumpDecision(t, "dir1/"+tc.name, d)
		}
		probeLog(t, "[dir1/%s] outcome: isError=%v riskLevel=%s class=%s text=%q",
			tc.name, out.IsError, out.RiskLevel, out.ErrorClass, out.Text)
	}
}

func TestDir2ConfiguredAllowlist(t *testing.T) {
	dataDir, _, st, body := realSpill(t)
	cap := &captured{}
	b := buildBridge([]string{dataDir}, cap, nil, tools.AnswerReject, "")
	out := exec(t, b, "fs.read", map[string]any{"path": st.Path})
	for _, d := range cap.decisions {
		dumpDecision(t, "dir2", d)
	}
	probeLog(t, "[dir2] L2 asks = %d isError=%v riskLevel=%s",
		cap.l2Asks, out.IsError, out.RiskLevel)
	if out.IsError {
		probeLog(t, "[dir2] READBACK FAILED: %q", out.Text)
		t.Fatalf("control leg failed: %s", out.Text)
	}
	if !strings.HasPrefix(out.Text, "W174C1-HEAD-MARK") || !strings.Contains(out.Text, "line 0001:") {
		t.Fatalf("readback text is not the artifact content: %q", out.Text[:min(120, len(out.Text))])
	}
	probeLog(t, "[dir2] readback first 80 bytes = %q", out.Text[:min(80, len(out.Text))])
	probeLog(t, "[dir2] readback total bytes = %d (artifact body %d)", len(out.Text), len(body))
}

func TestFakePointerBranch(t *testing.T) {
	_, artifacts, _, body := realSpill(t)
	roster := tools.NewTaskRoster()
	ghost := filepath.Join(artifacts, "tool-output-does-not-exist.txt")
	roster.Record("t-ghost", tools.TaskOutput{Text: body, ArtifactPath: ghost})
	roster.Record("t-dir", tools.TaskOutput{Text: body, ArtifactPath: artifacts})
	cap := &captured{}
	b := buildBridge(nil, cap, roster, tools.AnswerReject, "")
	for _, id := range []string{"t-ghost", "t-dir"} {
		out := exec(t, b, "task.output", map[string]any{"task_id": id})
		probeLog(t, "[fake/%s] isError=%v level=%s", id, out.IsError, out.RiskLevel)
		probeLog(t, "[fake/%s] full text (%d bytes) follows:\n%s", id, len(out.Text), out.Text)
	}
}
