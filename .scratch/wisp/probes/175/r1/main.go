// Command zz175r1 is ticket 175 AC#2's rig: it measures, on UNMODIFIED
// production code, the one question the ticket asks — when foreign content
// comes back through task.output, is prov.Mark even reached, does a mark land,
// and is that mark visible on the R4/Inspect side?
//
// The sub-questions are separated deliberately, because "Mark was never
// called" and "Mark was called but recorded nothing" produce the same empty
// ScopeTaints reading and are NOT the same defect:
//
//	LEG A  task.output across the real bridge      -> did the call site run?
//	LEG B  same content, same bridge, IN-ROSTER name -> the rig is live
//	LEG C  Provenance.Mark called directly with the OFF-roster name "task.output"
//	LEG D  Inspect (the R4 side) for each task id
//	LEG E  the Result task.output itself returns: does it carry an origin?
//
// Every leg runs against the real risk.Provenance and the real tools.Bridge —
// the injection surfaces AGENTS §1.3 allows. Nothing here patches production
// code (AC#2 forbids it): the only instrumentation is the two log sinks
// (tools.Options.Logf and the exported risk.Logf), which is exactly how "was
// Mark called" is separated from "did a mark land".
//
// Readings are written into this rig's own directory (see writeReport), never
// into whatever CWD the caller happened to inherit.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// probeText is the foreign content. It exists ONLY inside this process: no
// allowlisted file holds it, so the only way it can reach a taint index is the
// leg being measured.
const probeText = "外来内容探针 ZZ175r1-foreign-9c1f7d：这段文本只存在于后台任务的输出里"

// rosterTool is an IN-ROSTER source name (risk.SrcWebFetch) used as the
// same-bridge control. Its Decl mirrors task.output's (L0, no caps, no path
// params), so the only difference between LEG A and LEG B is the tool name.
const rosterTool = risk.SrcWebFetch

// offRosterTool is the name under investigation.
const offRosterTool = "task.output"

type foreignTool struct{ name string }

func (t foreignTool) Name() string                 { return t.name }
func (t foreignTool) Description() string          { return "returns foreign content (175-r1 rig control)" }
func (t foreignTool) Parameters() tools.JSONSchema { return tools.JSONSchema(`{"type":"object"}`) }

func (t foreignTool) Execute(ctx context.Context, _ json.RawMessage, _ func(string)) (tools.Result, error) {
	return tools.Result{Text: probeText, Origin: "https://example/foreign-body"}, nil
}

func foreignDecl(name string) tools.Decl {
	return tools.Decl{
		Capabilities: nil, Needs: nil,
		Declared: risk.L0, PathParams: nil,
		Resident: true, Provider: tools.KindBuiltin,
	}
}

type rig struct {
	b         *tools.Bridge
	prov      *risk.Provenance
	bridgeLog []string
	riskLog   []string
}

func main() {
	out := &strings.Builder{}
	pf := func(format string, args ...any) { fmt.Fprintf(out, format+"\n", args...) }

	r := &rig{}
	risk.Logf = func(format string, args ...any) {
		r.riskLog = append(r.riskLog, fmt.Sprintf(format, args...))
	}

	roster := tools.NewTaskRoster()
	roster.Record("bg-1", tools.TaskOutput{Text: probeText})

	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: roster}) {
		must(reg.Register(e))
	}
	must(reg.Register(tools.Entry{Tool: foreignTool{name: rosterTool}, Decl: foreignDecl(rosterTool)}))

	r.prov = risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil})
	r.b = tools.New(tools.Options{
		Registry:   reg,
		Paths:      tools.NewPathCanonicalizer([]string{tempDir()}, nil),
		Provenance: r.prov,
		Gate:       tools.NoGate{},
		Logf:       func(format string, args ...any) { r.bridgeLog = append(r.bridgeLog, fmt.Sprintf(format, args...)) },
	})

	pf("== rosters under test ==")
	pf("IsSensitiveSource(%q) = %v   IsSensitiveSource(%q) = %v",
		offRosterTool, risk.IsSensitiveSource(offRosterTool), rosterTool, risk.IsSensitiveSource(rosterTool))

	// ---- LEG A: foreign content returning through task.output ----------------
	// Runs FIRST, while the engine holds no taint at all: with any other
	// scope tainted, Inspect on this (never-opened) id would answer the
	// unbound-scope fail-closed instead of the question asked here.
	const taskA = "175r1-A-task-output"
	outcomeA, errA := r.execute(taskA, offRosterTool, `{"task_id":"bg-1"}`)
	pf("A1 execute            err=%v iserror=%v risk=%s text_prefix=%q",
		errA, outcomeA.IsError, outcomeA.RiskLevel, first(outcomeA.Text, 20))
	pf("A2 mark_return_early  bridge saw no marking path: bridge lines mentioning %q = %q",
		offRosterTool, containing(r.bridgeLog, "taint"))
	pf("A3 risk_log C25 lines %q", containing(r.riskLog, "C25"))
	r.reading(pf, "A4", taskA)

	// ---- LEG B: same bridge, same content, an IN-ROSTER source name ----------
	const taskB = "175r1-B-web-fetch"
	outcomeB, errB := r.execute(taskB, rosterTool, `{}`)
	pf("B1 execute            err=%v iserror=%v risk=%s", errB, outcomeB.IsError, outcomeB.RiskLevel)
	pf("B2 risk_log C25 lines %q", containing(r.riskLog, "C25"))
	r.reading(pf, "B3", taskB)

	// ---- LEG C: Mark() called DIRECTLY with the off-roster name --------------
	const taskC = "175r1-C-direct-mark"
	r.prov.OpenScope(taskC)
	markOK := r.prov.Mark(taskC, offRosterTool, "", probeText)
	pf("C1 direct Mark()      returned=%v (engine accepts an off-roster name)", markOK)
	pf("C2 risk_log C25 lines %q", containing(r.riskLog, "C25"))
	r.reading(pf, "C3", taskC)

	// ---- LEG E: does task.output's own Result carry any origin? --------------
	res, err := toolExecute(roster, `{"task_id":"bg-1"}`)
	pf("E1 tool_result_origin %q err=%v  (empty => task.output names no source of its own)", res.Origin, err)

	must(writeReport(out))
	fmt.Print(out.String())
}

// reading prints the three observable faces of one task id: the taint index,
// the R4/Inspect side, and the bridge's own close audit (was_open / dropped).
func (r *rig) reading(pf func(string, ...any), tag, taskID string) {
	taints := r.prov.ScopeTaints(taskID)
	pf("%s.1 scope_taints      n=%d %v", tag, len(taints), taints)
	hit, ok := r.prov.Inspect(taskID, "notify", map[string]any{"text": probeText})
	pf("%s.2 r4_inspect_hit    %v src=%q channel=%q", tag, ok, hit.Source(), hit.Channel)
	r.b.CloseTask(taskID)
	lines := containing(r.bridgeLog, "C25 scope closed")
	pf("%s.3 close_audit       %q", tag, lines)
}

func (r *rig) execute(taskID, tool, args string) (agent.ToolOutcome, error) {
	return r.b.Execute(context.Background(), agent.ToolRequest{
		TaskID: taskID, CorrelationID: taskID, CallID: "call-" + taskID,
		Name: tool, Args: json.RawMessage(args),
	})
}

func containing(lines []string, sub string) []string {
	var out []string
	for _, l := range lines {
		if strings.Contains(l, sub) {
			out = append(out, l)
		}
	}
	return out
}

func first(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		r = r[:n]
	}
	return string(r)
}

func toolExecute(roster *tools.TaskRoster, args string) (tools.Result, error) {
	e := tools.BuiltinTaskEntries(tools.TaskDeps{Roster: roster})[0]
	return e.Tool.Execute(context.Background(), json.RawMessage(args), nil)
}

func tempDir() string {
	d, err := os.MkdirTemp("", "175r1")
	must(err)
	return d
}

// writeReport drops the reading in this rig's own directory. The directory is
// taken from the source file's location via ZZ175R1_DIR (the runner exports it);
// the inherited CWD is never used, because a relative name there would land the
// reading outside the repo or in the repo root.
func writeReport(s *strings.Builder) error {
	dir := os.Getenv("ZZ175R1_DIR")
	if dir == "" {
		return fmt.Errorf("ZZ175R1_DIR is not set: refusing to write readings into the inherited CWD")
	}
	name := "rig-reading.txt"
	if v := os.Getenv("ZZ175R1_OUT"); v != "" {
		name = v
	}
	return os.WriteFile(filepath.Join(dir, name), []byte(s.String()), 0o644)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "rig failure:", err)
		os.Exit(1)
	}
}
