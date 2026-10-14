package main

// zz315a1probe_test.go — ticket 315's live re-measurement (AC#0, probe leg
// 315-a1) and, since leg 315-r1 landed the fix, the STANDING REGRESSION NAIL for
// the production shape it reproduced.
//
// WHAT IS MEASURED. A REAL `wisp run` composition root — real config.toml ->
// real secret.Store -> real approval.Gate (run.go's approval.New) -> real
// tools.Bridge with the real C25 engine -> real agent.Loop — driven by a
// scripted golden-SSE provider. Golden is the C5 seam AGENTS.md §1.3 lists, so
// the MODEL (not a test) emits the tool calls: the root's turn asks for
// task.spawn, and each child's OWN loop then asks for fs.read through the tool
// face its spawner handed down (internal/tools/subagent_197.go:296, the parent
// bridge), dispatched with the child's own task id (internal/agent/loop.go
// dispatch). The bridge judges every call exactly as it judges a real model's.
//
// WHY GOLDEN AND NOT mockllm. mockllm cannot script a tool_call — its own
// harness comment says no tool_choice is sent, so it never answers with one, and
// a child spawned on that fixture never reaches the bridge at all. The golden
// replayer serves scripted SSE bytes over a real HTTP listener through the real
// adapter, which is the seam the repo's own rules name for exactly this.
//
// THE CHAIN THIS FILE PUTS ON THE WIRE (ticket 315 §现量 #3/#4):
//
//	root turn 1  -> task.spawn(甲)          (root's own bridge call, no marks anywhere yet)
//	child 甲  1  -> fs.read(allowlisted)    (甲's FIRST bridge call: control leg)
//	child 甲  ends  -> stampConclusion marks the PARENT's C25 scope (subagent_197.go finalize/stampConclusion)
//	root turn 2  -> task.spawn(乙)          (root's own scope now registered+marked)
//	child 乙  1  -> fs.read(allowlisted)    (THE measured call: sibling child's first batch)
//	child 乙  ends, root turn 3 -> run closes
//
// The assertions are the EXPECTED-CORRECT readings (an allowlisted read inside
// the allowlist is L0; no confirmation card is owed).
//
// RED BEFORE, GREEN AFTER, AND WHY IT IS NO LONGER ENV-GATED. While ticket 315's
// bug was live, case B was RED by construction (乙's first bridge call judged
// L2/R4 "unbound-scope", a real card, auto-reject after the 31s production
// default), so 315-a1 landed it behind WISP_315A1=1. Two things ended that:
//   1. leg 315-r1 fixed the root cause (cmd/wisp/run.go admitTask now OPENS the
//      task's own C25 scope at its start, pairing the CloseTask ticket 151 put
//      in the same hook), and the same run is green and fast — 0.21s, no card,
//      no 31s wait (pre-fix: rc=1, 31.22s; readings in
//      .scratch/wisp/probes/315/{a1,r1}/).
//   2. An env-gated t.Skip is NOT neutral in this repo: the cli scope
//      (scripts/wisp-cli-tests.sh -> portable-tests.sh --scope=cli) runs every
//      cmd/wisp test through tools/d22scan/runtests.sh, where ANY top-level
//      `--- SKIP` is fatal (ticket 71 AC#3), and no ledger row accounts for this
//      name. Measured at HEAD 5ed56c2f: the strict runner exits 1 with SKIP=1 on
//      this very test (probes/315/r1/logs/12-prefix-gated-skip-under-strict-runner.txt).
//      De-gating removes the skip, registers nothing in scripts/portable-tests.sh
//      (another leg's face) and gives the fix a permanent witness in the default
//      package run, which is this package's own tradition for composition-root
//      tests (run_test.go, task_scope_close_151_test.go).
//
// Case A (TestZZ315A1RootLoopFirstBridgeCall) is the no-other-marks control: no
// seed scope holds a mark, so the root's own first bridge call must execute as
// L0 — that is the baseline which makes case B's reading mean something.
//
// 只读复量：this file is a zz-prefixed carrier; the only production byte ticket
// 315 changed is cmd/wisp/run.go's admitTask. Kept (临时件只建不删).

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// locked console buffers: the run writes stdout/stderr from more than one
// goroutine (its own sink/pump), and this probe polls them while it runs, so a
// bytes.Buffer (what run_test.go's fixture uses, read only after the run) would
// be a data race here.
// ---------------------------------------------------------------------------

type zz315Buf struct {
	mu sync.Mutex
	b  strings.Builder
}

func (z *zz315Buf) Write(p []byte) (int, error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	return z.b.Write(p)
}

func (z *zz315Buf) String() string {
	z.mu.Lock()
	defer z.mu.Unlock()
	return z.b.String()
}

// ---------------------------------------------------------------------------
// the golden script: one zz315Turn per model reply, in the exact order the
// requests arrive (the replayer serves sequentially and a task.spawn call
// BLOCKS its caller until the child joins, so the order below is deterministic:
// root -> 甲 (two turns) -> root -> 乙 (two turns) -> root).
// ---------------------------------------------------------------------------

type zz315Turn struct {
	text string            // assistant text of this reply
	tool string            // tool name when this reply calls one ("" = text-only close)
	args map[string]string // the tool's string arguments
	id   string            // the scripted tool_call id
}

// zz315Data renders one SSE data line with a blank dispatcher line after it.
func zz315Data(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal sse chunk: %v", err)
	}
	return "data: " + string(b) + "\n\n"
}

// zz315TurnSSE renders one openai-chat reply: role chunk, optional text chunk,
// optional complete tool_call, finish chunk, usage, [DONE]. The byte shape
// mirrors internal/llm/testdata/golden/tool-call.sse (fragmented args are not
// needed here; the adapter accumulates a single complete arguments string the
// same way).
func zz315TurnSSE(t *testing.T, turn zz315Turn) string {
	t.Helper()
	var b strings.Builder
	b.WriteString(zz315Data(t, map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{"role": "assistant", "content": ""}, "finish_reason": nil,
	}}}))
	if turn.text != "" {
		b.WriteString(zz315Data(t, map[string]any{"choices": []any{map[string]any{
			"index": 0, "delta": map[string]any{"content": turn.text}, "finish_reason": nil,
		}}}))
	}
	finish := "stop"
	if turn.tool != "" {
		args, err := json.Marshal(turn.args)
		if err != nil {
			t.Fatalf("marshal tool args: %v", err)
		}
		b.WriteString(zz315Data(t, map[string]any{"choices": []any{map[string]any{
			"index": 0,
			"delta": map[string]any{"tool_calls": []any{map[string]any{
				"index": 0, "id": turn.id, "type": "function",
				"function": map[string]any{"name": turn.tool, "arguments": string(args)},
			}}},
			"finish_reason": nil,
		}}}))
		finish = "tool_calls"
	}
	b.WriteString(zz315Data(t, map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{}, "finish_reason": finish,
	}}}))
	b.WriteString(`data: {"choices":[],"usage":{"prompt_tokens":11,"completion_tokens":3,"total_tokens":14}}` + "\n\n")
	b.WriteString("data: [DONE]\n")
	return b.String()
}

// zz315Script joins the turns into one golden file.
func zz315Script(t *testing.T, turns []zz315Turn) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("# wisp golden sse v1\n")
	b.WriteString("# @scenario zz315a1 ticket 315 AC#0 live probe\n")
	for _, turn := range turns {
		b.WriteString("# @response 200\n")
		b.WriteString(zz315TurnSSE(t, turn))
		b.WriteString("\n")
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// the flow: config + DPAPI blob + golden server, driven through runTextTask
// exactly as `wisp run` drives it.
// ---------------------------------------------------------------------------

type zz315Flow struct {
	dir  string
	srv  *httptest.Server
	rep  *golden.Replayer
	out  *zz315Buf
	errb *zz315Buf

	rtMu sync.Mutex
	rt   *agentRuntime

	hookRT func(*agentRuntime)
}

// zz315WriteConfig writes the fixture config. Two deliberate choices:
//   - confirm_timeout_sec = 31 (the validator's own floor, config/validate.go):
//     the shortest legal L2 card, so the production auto-reject path runs in
//     ~31s instead of 300 — a real card either way;
//   - per_tool_timeout_ms = 60000 so the C22 budget outlives the card (the
//     bridge's default is 30s, which would abandon the call one second BEFORE
//     the card expires and turn the reading into a ctx-deadline reject).
func zz315WriteConfig(t *testing.T, dir, baseURL string) {
	t.Helper()
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]

[agent]
per_tool_timeout_ms = 60000

[risk]
confirm_timeout_sec = 31
`, baseURL+"/v1", filepath.ToSlash(dir))
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}
}

func zz315Setup(t *testing.T, build func(dir string) string) *zz315Flow {
	t.Helper()
	f := &zz315Flow{dir: t.TempDir(), out: &zz315Buf{}, errb: &zz315Buf{}}
	srv, rep, err := golden.ServerFromBytes([]byte(build(f.dir)))
	if err != nil {
		t.Fatalf("golden.ServerFromBytes: %v", err)
	}
	t.Cleanup(srv.Close)
	f.srv, f.rep = srv, rep
	zz315WriteConfig(t, f.dir, srv.URL)
	return f
}

func (f *zz315Flow) run(task string) int {
	return runTextTask(runSpec{
		argv:    []string{task},
		stdout:  f.out,
		stderr:  f.errb,
		dataDir: f.dir,
		onRuntime: func(rt *agentRuntime) {
			f.rtMu.Lock()
			f.rt = rt
			h := f.hookRT
			f.rtMu.Unlock()
			if h != nil {
				h(rt)
			}
		},
		notify: func(title, body string) error { return nil },
	})
}

func (f *zz315Flow) runtime() *agentRuntime {
	f.rtMu.Lock()
	defer f.rtMu.Unlock()
	return f.rt
}

// rows reopen the run's own store (the composition root closes it on exit) and
// list one task's authoritative tool_call rows.
func (f *zz315Flow) rows(t *testing.T, taskID string) []memory.ToolCall {
	t.Helper()
	st, err := memory.Open(f.dir)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := st.ListToolCallsByTask(ctx, taskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask(%q): %v", taskID, err)
	}
	return rows
}

// rootTaskID reads the id this run reported on its status line.
func (f *zz315Flow) rootTaskID(t *testing.T) string {
	t.Helper()
	m := taskIDRe.FindStringSubmatch(f.out.String())
	if m == nil {
		t.Fatalf("no status line in stdout:\n%s", f.out.String())
	}
	return m[1]
}

// children197 returns (label -> taskID) for the subagent rows under rootID.
func (f *zz315Flow) children(t *testing.T, rootID string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, row := range f.runtime().taskRosterState().Rows {
		if row.Kind == tools.TaskKindSubagent && row.ParentTaskID == rootID {
			out[row.Label] = row.TaskID
		}
	}
	return out
}

// scopeTable renders the C25 scope snapshot the probe needs: which scope ids
// are OPEN and hold marks at the sampling instant, and which were never opened.
// ScopeTaints is the engine's own introspection view (its content is never
// printed, only tool/origin/length).
func (f *zz315Flow) scopeTable(ids map[string]string) string {
	rt := f.runtime()
	if rt == nil {
		return "(runtime unavailable)"
	}
	var b strings.Builder
	for label, id := range ids {
		tt := rt.c25.ScopeTaints(id)
		if len(tt) == 0 {
			fmt.Fprintf(&b, "  %-28s %s: marks=0 (open/closed not distinguishable via ScopeTaints; the run's own \"C25 scope closed was_open=...\" audit names it)\n", label, id)
			continue
		}
		fmt.Fprintf(&b, "  %-28s %s:", label, id)
		for _, ti := range tt {
			fmt.Fprintf(&b, " [tool=%s origin=%q runes=%d]", ti.Tool, ti.Origin, ti.RuneLen)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// cardSnap is the mid-run reading taken the instant the L2 card reaches stdout.
type zz315CardSnap struct {
	seen   bool
	line   string
	scopes string
	roster string
}

// watchCard polls stdout (safe: locked buffer) until the first L2 confirmation
// card appears, then snapshots the roster and the C25 scopes while the card is
// still up. Returns a reader and a stop func.
func (f *zz315Flow) watchCard() (func() zz315CardSnap, func()) {
	var mu sync.Mutex
	var snap zz315CardSnap
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
			}
			out := f.out.String()
			if i := strings.Index(out, "[确认 L2"); i >= 0 && !strings.Contains(out[:i], "[确认 L2") {
				line := out[i:]
				if j := strings.IndexByte(line, '\n'); j >= 0 {
					line = line[:j]
				}
				mu.Lock()
				if !snap.seen {
					snap.seen = true
					snap.line = line
					rootID := ""
					kids := map[string]string{}
					if rt := f.runtime(); rt != nil {
						for _, row := range rt.taskRosterState().Rows {
							switch row.Kind {
							case tools.TaskKindRoot:
								rootID = row.TaskID
							case tools.TaskKindSubagent:
								kids[row.Label] = row.TaskID
							}
						}
						ids := map[string]string{"root": rootID, "run:project-instructions": "run:project-instructions"}
						for label, id := range kids {
							ids["child "+label] = id
						}
						snap.scopes = f.scopeTable(ids)
						var rb strings.Builder
						for _, row := range rt.taskRosterState().Rows {
							fmt.Fprintf(&rb, "  %s kind=%s parent=%s label=%q status=%q known=%v\n",
								row.TaskID, row.Kind, row.ParentTaskID, row.Label, row.Status, row.StatusKnown)
						}
						snap.roster = rb.String()
					}
				}
				mu.Unlock()
				return
			}
			time.Sleep(2 * time.Millisecond)
		}
	}()
	return func() zz315CardSnap {
			mu.Lock()
			defer mu.Unlock()
			return snap
		}, func() {
			select {
			case <-stop:
			default:
				close(stop)
			}
			select {
			case <-done:
			case <-time.After(3 * time.Second):
			}
		}
}

// cardLines pulls the card + dismissal lines out of stdout verbatim.
func zz315CardLines(out string) string {
	i := strings.Index(out, "[确认 L2")
	if i < 0 {
		return "(no L2 card on stdout)"
	}
	j := strings.Index(out[i:], "\n\n")
	var seg string
	if j < 0 {
		seg = out[i:]
	} else {
		seg = out[i : i+j]
	}
	// attach the next dismissal line, if any
	if k := strings.Index(out[i:], "[dismissed]"); k >= 0 {
		end := strings.IndexByte(out[i+k:], '\n')
		if end < 0 {
			end = len(out) - i - k
		}
		seg += "\n" + out[i+k:i+k+end]
	}
	return seg
}

// auditLine pulls the first stderr audit line containing needle, verbatim.
func zz315AuditLine(errText, needle string) string {
	for _, line := range strings.Split(errText, "\n") {
		if strings.Contains(line, needle) {
			return line
		}
	}
	return "(no line containing " + needle + ")"
}

// ---------------------------------------------------------------------------
// Case A: the no-other-marks control. No seed scope holds a mark, so the ROOT
// loop's own first bridge call (an allowlisted fs.read, model-emitted) must
// execute L0: that is the baseline making case B's unbound reading mean
// something (scopeMarks fail-closes ONLY "while other scopes hold taints").
// ---------------------------------------------------------------------------

func TestZZ315A1RootLoopFirstBridgeCall(t *testing.T) {
	note := ""
	f := zz315Setup(t, func(dir string) string {
		note = filepath.Join(dir, "zz315a1-note.txt")
		return zz315Script(t, []zz315Turn{
			{tool: "fs.read", args: map[string]string{"path": filepath.ToSlash(note)}, id: "call_zz315_a1_root_read"},
			{text: "根任务读完了，收口。"},
		})
	})
	if err := os.WriteFile(note, []byte("zz315a1 note body：窗口编号例如 315-A1。"), 0o600); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	code := f.run("315a1 探针：根任务先读一份允许目录内的笔记")
	took := time.Since(start)
	t.Logf("case A: rc=%d in %s; golden requests=%v", code, took, f.rep.Requests)
	t.Logf("case A raw stdout:\n%s", f.out.String())
	t.Logf("case A raw stderr:\n%s", f.errb.String())
	if code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.errb.String())
	}
	rootID := f.rootTaskID(t)
	rows := f.rows(t, rootID)
	for _, r := range rows {
		t.Logf("case A row: tool=%s level=%s decision=%s outcome=%s class=%q corr=%q", r.Tool, r.RiskLevel, r.Decision, r.Outcome, r.ErrorClass, r.CorrelationID)
	}
	if !strings.Contains(f.out.String(), "echo") && !strings.Contains(f.out.String(), "根任务读完") {
		t.Logf("case A stdout (no streamed text match):\n%s", f.out.String())
	}
	if strings.Contains(f.out.String(), "[确认") {
		t.Errorf("case A: a confirmation card appeared with no seed mark anywhere — stdout:\n%s", f.out.String())
	}
	var read *memory.ToolCall
	for i := range rows {
		if rows[i].Tool == "fs.read" {
			read = &rows[i]
		}
	}
	if read == nil {
		t.Fatalf("case A: no fs.read row for the root task %s: %+v", rootID, rows)
	}
	if read.RiskLevel != memory.RiskL0 || read.Decision != agent.DecisionAllow || read.Outcome != agent.OutcomeSuccess {
		t.Errorf("case A: the no-other-marks control read is %s/%s/%s, want L0/allow/success: %+v",
			read.RiskLevel, read.Decision, read.Outcome, *read)
	}
}

// ---------------------------------------------------------------------------
// Case B: the ticket's shape — a sibling child's first bridge call after the
// first child's conclusion has stamped the PARENT's C25 scope. This is the
// STANDING NAIL for ticket 315: it is red on the pre-fix tree (probes/315/a1)
// and green on the fix (admitTask opens the task's own scope at its start).
// ---------------------------------------------------------------------------

const (
	zz315LabelJia = "315a1 甲组只读核对"
	zz315LabelYi  = "315a1 乙组独立复查"
)

func TestZZ315A1SiblingChildFirstBridgeCall(t *testing.T) {
	note := ""
	f := zz315Setup(t, func(dir string) string {
		note = filepath.Join(dir, "zz315a1-note.txt")
		return zz315Script(t, []zz315Turn{
			// ① root turn 1: spawn 甲 (root's own bridge call; no marks anywhere yet)
			{tool: "task.spawn", id: "call_zz315_a1_spawn_jia", args: map[string]string{
				"description": zz315LabelJia,
				"prompt":      "甲组先把笔记读完，再回来汇报结论。",
			}},
			// ② 甲 turn 1: its FIRST bridge call — an allowlisted read
			{tool: "fs.read", id: "call_zz315_a1_jia_read", args: map[string]string{"path": filepath.ToSlash(note)}},
			// ③ 甲 turn 2: the conclusion; on join it is stamped into the PARENT's scope
			{text: "甲组读完笔记，没有发现需要复核的条目。"},
			// ④ root turn 2: spawn 乙 (root's scope now holds 甲's conclusion mark)
			{tool: "task.spawn", id: "call_zz315_a1_spawn_yi", args: map[string]string{
				"description": zz315LabelYi,
				"prompt":      "乙组独立复核编号规则，回来报到。",
			}},
			// ⑤ 乙 turn 1: THE measured call — a sibling child's first bridge call
			{tool: "fs.read", id: "call_zz315_a1_yi_read", args: map[string]string{"path": filepath.ToSlash(note)}},
			// ⑥ 乙 turn 2
			{text: "乙组这边已经收尾。"},
			// ⑦ root turn 3
			{text: "父任务收口：两组都已回来。"},
		})
	})
	if err := os.WriteFile(note, []byte("zz315a1 note body：窗口编号例如 315-A1。"), 0o600); err != nil {
		t.Fatal(err)
	}

	snapOf, stopWatch := f.watchCard()
	start := time.Now()
	code := f.run("315a1 探针：父任务先后派生两枚子代理各自读一份允许目录内的笔记")
	took := time.Since(start)
	stopWatch()
	snap := snapOf()

	rootID := ""
	if rt := f.runtime(); rt != nil {
		for _, row := range rt.taskRosterState().Rows {
			if row.Kind == tools.TaskKindRoot {
				rootID = row.TaskID
			}
		}
	}
	kids := map[string]string{}
	if rootID != "" {
		kids = f.children(t, rootID)
	}
	jia, yi := kids[zz315LabelJia], kids[zz315LabelYi]

	t.Logf("case B: rc=%d in %s; golden requests=%d %v", code, took, len(f.rep.Requests), f.rep.Requests)
	t.Logf("case B: root=%s 甲=%s 乙=%s", rootID, jia, yi)
	t.Logf("case B raw stdout:\n%s", f.out.String())
	t.Logf("case B raw stderr:\n%s", f.errb.String())
	if code != 0 {
		t.Fatalf("exit %d — the scripted run did not complete; readings below would be misaligned.\nstdout:\n%s\nstderr:\n%s",
			code, f.out.String(), f.errb.String())
	}
	if jia == "" || yi == "" {
		t.Fatalf("the roster gave no two-child pair (甲=%q 乙=%q); rows: %+v", jia, yi, f.runtime().taskRosterState())
	}

	t.Logf("case B card snapshot (seen=%v): %s", snap.seen, snap.line)
	t.Logf("case B mid-card roster:\n%s", snap.roster)
	t.Logf("case B mid-card C25 scopes:\n%s", snap.scopes)
	t.Logf("case B stdout card segment:\n%s", zz315CardLines(f.out.String()))
	t.Logf("case B stderr expiry audit: %s", zz315AuditLine(f.errb.String(), "ANSWER-EXPIRED"))

	// The three calls this probe judges, in order.
	jiaRead := zz315FindRow(t, f.rows(t, jia), "fs.read")
	yiRead := zz315FindRow(t, f.rows(t, yi), "fs.read")
	for _, r := range f.rows(t, rootID) {
		t.Logf("case B root row: tool=%s level=%s decision=%s outcome=%s class=%q", r.Tool, r.RiskLevel, r.Decision, r.Outcome, r.ErrorClass)
		if r.Tool == "task.spawn" && r.Decision != agent.DecisionAllow {
			t.Errorf("the parent's own task.spawn call was judged %s/%s — the child-conclusion mark should not match those args: %+v",
				r.RiskLevel, r.Decision, r)
		}
	}
	t.Logf("case B 甲 read row: %+v", *jiaRead)
	t.Logf("case B 乙 read row: %+v", *yiRead)

	// Control leg: with no marks anywhere yet, 甲's own first bridge call runs.
	if jiaRead.RiskLevel != memory.RiskL0 || jiaRead.Decision != agent.DecisionAllow || jiaRead.Outcome != agent.OutcomeSuccess {
		t.Errorf("CONTROL: 甲's first bridge call (no marks anywhere yet) is %s/%s/%s, want L0/allow/success: %+v",
			jiaRead.RiskLevel, jiaRead.Decision, jiaRead.Outcome, *jiaRead)
	}

	// THE measured leg: 甲's conclusion now holds a mark in the PARENT's scope;
	// 乙 is a sibling child whose own scope was never opened. An allowlisted
	// read must still be L0 — a false L2 here is ticket 315 reproduced.
	if yiRead.RiskLevel != memory.RiskL0 || yiRead.Decision != agent.DecisionAllow || yiRead.Outcome != agent.OutcomeSuccess {
		t.Errorf("TICKET 315 REPRODUCED: 乙's first bridge call (allowlisted fs.read, sibling child's own task id) was judged %s/%s class=%q while 甲's identical call was L0.\n"+
			"  card line (verbatim): %s\n"+
			"  stdout card segment:\n%s\n"+
			"  stderr audit (verbatim): %s\n"+
			"  row: %+v\n"+
			"  mid-card scopes:\n%s",
			yiRead.RiskLevel, yiRead.Decision, yiRead.ErrorClass,
			snap.line, zz315CardLines(f.out.String()),
			zz315AuditLine(f.errb.String(), "ANSWER-EXPIRED"), *yiRead, snap.scopes)
	}
}

// ---------------------------------------------------------------------------
// Case C: the OTHER shape ticket 315 named as unmeasured — a seed mark that
// exists in a scope that is NOT the task's own, BEFORE that task's first bridge
// call. The run's data dir carries AGENTS.md, so the projctx loader stamps
// run:project-instructions (risk.SrcFSRead) while building the first request.
// On the pre-fix tree even the FIRST task's own allowlisted fs.read then fails
// closed (its own scope unregistered + a seeded mark elsewhere => L2/R4); the
// admitTask open covers this for the same reason it covers case B: the task's
// own scope is registered before its first bridge call. Not a new fix - a
// reading of what the one line already covers (dispatch 315-r1 item 4).
//
// The precondition is ASSERTED, not assumed: the load-manifest line naming the
// global-tier file must be on stderr, so "the seed never landed" cannot read as
// a green.
// ---------------------------------------------------------------------------

func TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask(t *testing.T) {
	note := ""
	f := zz315Setup(t, func(dir string) string {
		note = filepath.Join(dir, "zz315a1-note.txt")
		return zz315Script(t, []zz315Turn{
			{tool: "fs.read", args: map[string]string{"path": filepath.ToSlash(note)}, id: "call_zz315_a1c_root_read"},
			{text: "根任务读过笔记，收口。"},
		})
	})
	if err := os.WriteFile(note, []byte("zz315a1 note body：窗口编号例如 315-A1。"), 0o600); err != nil {
		t.Fatal(err)
	}
	// The global tier is the run's data dir (run.go hands it to the loader as
	// GlobalDir), and AGENTS.md is first in projctx's FileNamePriority.
	instr := filepath.Join(f.dir, "AGENTS.md")
	if err := os.WriteFile(instr, []byte("项目说明种子正文：这条说明必须在任何工具调用之前被 C25 盖戳。"), 0o600); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	code := f.run("315a1 探针：数据目录带项目说明时，根任务的第一发走桥调用")
	took := time.Since(start)
	t.Logf("case C: rc=%d in %s; golden requests=%v", code, took, f.rep.Requests)
	t.Logf("case C raw stdout:\n%s", f.out.String())
	t.Logf("case C raw stderr:\n%s", f.errb.String())
	if code != 0 {
		t.Fatalf("exit %d — the scripted run did not complete; the reading below would be misaligned.\nstdout:\n%s\nstderr:\n%s",
			code, f.out.String(), f.errb.String())
	}
	// Precondition first: the seed mark really landed. The manifest line the
	// loader logs names the file, its tier and its fate.
	if !strings.Contains(f.errb.String(), instr) || !strings.Contains(f.errb.String(), "tier=global") {
		t.Fatalf("前置条件缺失：数据目录里那枚 AGENTS.md 没有被装载（stderr 里找不到 tier=global 那一行），本形量不到，不算绿\nstderr:\n%s",
			f.errb.String())
	}
	if strings.Contains(f.out.String(), "[确认") {
		t.Errorf("有种子 mark（run:project-instructions）时，根任务的第一发走桥调用出了确认卡：stdout:\n%s", f.out.String())
	}
	read := zz315FindRow(t, f.rows(t, f.rootTaskID(t)), "fs.read")
	t.Logf("case C root read row: %+v", *read)
	if read.RiskLevel != memory.RiskL0 || read.Decision != agent.DecisionAllow || read.Outcome != agent.OutcomeSuccess {
		t.Errorf("TICKET 315 REPRODUCED (project-instructions seed shape): the root's own first bridge call (allowlisted fs.read) was judged %s/%s/%s class=%q while a seed mark sits only in run:project-instructions.\n"+
			"  stderr audit (verbatim): %s\n"+
			"  stdout card segment:\n%s\n"+
			"  row: %+v",
			read.RiskLevel, read.Decision, read.Outcome, read.ErrorClass,
			zz315AuditLine(f.errb.String(), "R4"), zz315CardLines(f.out.String()), *read)
	}
}

func zz315FindRow(t *testing.T, rows []memory.ToolCall, tool string) *memory.ToolCall {
	t.Helper()
	for i := range rows {
		if rows[i].Tool == tool {
			return &rows[i]
		}
	}
	t.Fatalf("no %s row in %+v", tool, rows)
	return nil
}
