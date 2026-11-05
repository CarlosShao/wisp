package main

// Ticket 378 AC#2 (form 甲; the ruling is docs/evidence/s1/378-a0.md 「AC#1 裁形」):
// every task.spawn-derived CHILD TASK gets ONE task_log row of its own, so the
// tool_call rows the production bridge books under the child's task id have a
// parent that really exists.
//
// THE ACCEPTANCE RULER IS NOT "a row appeared". The ticket says so in as many
// words (「只测『多了一行 task_log』」): what breaks while a child has no row is
// the FORENSIC chain - a tool_call row whose parent was never written - and the
// reading for that is SQLite's PRAGMA foreign_key_check, taken after one real
// task.spawn. That is what this case reads, on the database this run itself wrote.
//
// WHY THE WHOLE COMPOSITION ROOT. The decision under test is neither
// internal/agent's (the loop writes whatever Journal it is handed) nor
// internal/tools' (the spawner copies the Options snapshot it is given): it is
// cmd/wisp's - which payload execute() publishes onto rt.loopOpt, i.e. what a
// child inherits from its parent's loop. A case that assembled the child's
// agent.Options by hand would prove a REPRODUCTION of that decision and would
// stay green with the production line deleted (the sibling case
// tasklog_subloop_378_test.go owns the other half of the same journal - that the
// inherited payload books no tool_call row of its own; it pins that with a
// journal-less bridge, so its child's calls are never booked at all). Here the
// model emits the calls through the C5 seam (golden SSE, AGENTS.md §1.3), the
// child's own loop dispatches a real fs.read through the SAME bridge the root
// uses (run.go's `Journal: mem`), and every row below is read back from the real
// store.
//
// WHAT MAKES A 0 MEAN SOMETHING (all four are asserted, not assumed):
//
//	[A] the child's tool_call row EXISTS (one, `fs.read`)    -> not "nothing ran"
//	[B] the child's task_log row EXISTS                      -> the parent the
//	    ruler needs; TaskLogByID reads it back
//	[C] the root's own rows exist AND stay clean             -> the positive
//	    control (ticket 373's half; this ticket must not break it)
//	[D] total tool_call rows == the model's tool calls (2)   -> the inherited
//	    journal did NOT add a second row per call (hard reminder #1: the record
//	    page must not be flooded)
//
// THE 改前 READING is not a second arm of this function. Deleting the one
// publication line in run.go (`rt.loopOpt = loopOpts`) IS "改前", and this case
// goes red on it, naming the child row and its tool. A self-built arm would only
// prove that the test can write an orphan - not that the wiring prevents one.

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/tools"
)

const (
	// fk378Label is the spawn's short title, i.e. the roster row's label and the
	// handle this case uses to find the child.
	fk378Label = "378-r1 子任务读笔记"
	// fk378NoteName is the file the CHILD reads. It lives in the run's data dir,
	// which the fixture allowlists as its only fs root.
	fk378NoteName = "fk378-note.txt"
	// fk378Task is the parent's task text.
	fk378Task = "378：父任务派生一枚子代理，让它把笔记读完再回来"
)

// fk378TaskIDRe is this file's own copy of the status-line reader (run_test.go
// keeps one for its own assertions; sharing it would tie this case to that
// fixture's helper).
var fk378TaskIDRe = regexp.MustCompile(`任务 (\S+) 结束`)

// fk378Buf is a console buffer the run writes from more than one goroutine.
type fk378Buf struct {
	mu sync.Mutex
	b  strings.Builder
}

func (z *fk378Buf) Write(p []byte) (int, error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	return z.b.Write(p)
}

func (z *fk378Buf) String() string {
	z.mu.Lock()
	defer z.mu.Unlock()
	return z.b.String()
}

// ---------------------------------------------------------------------------
// the scripted model: one fk378Turn per reply, in the order the requests arrive.
// A task.spawn call BLOCKS its caller until the child joins, so the order is
// deterministic: root turn 1 (spawn) -> child turn 1 (fs.read) -> child turn 2
// (its conclusion) -> root turn 2 (close).
// ---------------------------------------------------------------------------

type fk378Turn struct {
	text string            // assistant text of this reply
	tool string            // tool name when this reply calls one ("" = text-only)
	args map[string]string // the tool's string arguments
	id   string            // scripted tool_call id
}

// fk378Data renders one SSE data line with the blank dispatcher line after it.
func fk378Data(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal sse chunk: %v", err)
	}
	return "data: " + string(b) + "\n\n"
}

// fk378TurnSSE renders one openai-chat reply: role chunk, optional text chunk,
// optional complete tool_call, finish chunk, usage, [DONE]. The byte shape
// mirrors internal/llm/testdata/golden/tool-call.sse (fragmented arguments are
// not needed here: the adapter accumulates one complete arguments string either
// way).
func fk378TurnSSE(t *testing.T, turn fk378Turn) string {
	t.Helper()
	var b strings.Builder
	b.WriteString(fk378Data(t, map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{"role": "assistant", "content": ""}, "finish_reason": nil,
	}}}))
	if turn.text != "" {
		b.WriteString(fk378Data(t, map[string]any{"choices": []any{map[string]any{
			"index": 0, "delta": map[string]any{"content": turn.text}, "finish_reason": nil,
		}}}))
	}
	finish := "stop"
	if turn.tool != "" {
		args, err := json.Marshal(turn.args)
		if err != nil {
			t.Fatalf("marshal tool args: %v", err)
		}
		b.WriteString(fk378Data(t, map[string]any{"choices": []any{map[string]any{
			"index": 0,
			"delta": map[string]any{"tool_calls": []any{map[string]any{
				"index": 0, "id": turn.id, "type": "function",
				"function": map[string]any{"name": turn.tool, "arguments": string(args)},
			}}},
			"finish_reason": nil,
		}}}))
		finish = "tool_calls"
	}
	b.WriteString(fk378Data(t, map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{}, "finish_reason": finish,
	}}}))
	b.WriteString(`data: {"choices":[],"usage":{"prompt_tokens":11,"completion_tokens":3,"total_tokens":14}}` + "\n\n")
	b.WriteString("data: [DONE]\n")
	return b.String()
}

// fk378Script joins the turns into one golden file.
func fk378Script(t *testing.T, turns []fk378Turn) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("# wisp golden sse v1\n")
	b.WriteString("# @scenario fk378 ticket 378 AC#2 foreign_key_check after one task.spawn\n")
	for _, turn := range turns {
		b.WriteString("# @response 200\n")
		b.WriteString(fk378TurnSSE(t, turn))
		b.WriteString("\n")
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// the run fixture: real config.toml -> real secret.Store -> real approval.Gate
// -> real tools.Bridge -> real agent.Loop, driven exactly as `wisp run` drives
// it. Only the model's bytes and the notification poster are injected.
// ---------------------------------------------------------------------------

type fk378Flow struct {
	dir  string
	srv  *httptest.Server
	rep  *golden.Replayer
	out  *fk378Buf
	errb *fk378Buf

	rtMu sync.Mutex
	rt   *agentRuntime
}

// fk378WriteConfig writes the fixture's config.toml and the fake key.
//
// Two knobs are load-bearing, both copied from the zz315 fixture that proved
// this shape fast: confirm_timeout_sec = 31 is the validator's own floor (no
// card is expected here - fs.read inside the allowlist is L0 - but a shape that
// regressed into asking would then end in 31s instead of 300s), and
// per_tool_timeout_ms = 60000 keeps the C22 budget clear of that window.
func fk378WriteConfig(t *testing.T, dir, baseURL string) {
	t.Helper()
	keyRef := fixtureKeyRef181(t, dir, "dpapi:acme", fakeStoreKey)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = %q

[llm.providers.acme.models.m1]
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]

[agent]
per_tool_timeout_ms = 60000

[risk]
confirm_timeout_sec = 31
`, baseURL+"/v1", keyRef, filepath.ToSlash(dir))
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// fk378Setup builds the fixture, calling build(dir) to render the golden script.
func fk378Setup(t *testing.T, build func(dir string) string) *fk378Flow {
	t.Helper()
	f := &fk378Flow{dir: t.TempDir(), out: &fk378Buf{}, errb: &fk378Buf{}}
	srv, rep, err := golden.ServerFromBytes([]byte(build(f.dir)))
	if err != nil {
		t.Fatalf("golden.ServerFromBytes: %v", err)
	}
	t.Cleanup(srv.Close)
	f.srv, f.rep = srv, rep
	fk378WriteConfig(t, f.dir, srv.URL)
	return f
}

// run drives the composition root exactly as `wisp run` does, with the real
// store, the real gate and the real bridge; only the notifier is a no-op.
func (f *fk378Flow) run(task string) int {
	return runTextTask(runSpec{
		argv:    []string{task},
		stdout:  f.out,
		stderr:  f.errb,
		dataDir: f.dir,
		onRuntime: func(rt *agentRuntime) {
			f.rtMu.Lock()
			f.rt = rt
			f.rtMu.Unlock()
		},
		notify: func(string, string) error { return nil },
	})
}

func (f *fk378Flow) runtime() *agentRuntime {
	f.rtMu.Lock()
	defer f.rtMu.Unlock()
	return f.rt
}

// rootTaskID reads the id this run reported on its own status line.
func (f *fk378Flow) rootTaskID(t *testing.T) string {
	t.Helper()
	m := fk378TaskIDRe.FindStringSubmatch(f.out.String())
	if m == nil {
		t.Fatalf("这一发没有打印任务状态行：\nstdout:\n%s\nstderr:\n%s", f.out.String(), f.errb.String())
	}
	return m[1]
}

// childIDs lists the subagent rows the roster recorded under rootID, in the
// roster's own (task-id) order.
func (f *fk378Flow) childIDs(t *testing.T, rootID string) []string {
	t.Helper()
	rt := f.runtime()
	if rt == nil {
		t.Fatal("onRuntime 没有交出这发的 runtime，名册读不到")
	}
	var out []string
	for _, row := range rt.taskRosterState().Rows {
		if row.Kind == tools.TaskKindSubagent && row.ParentTaskID == rootID {
			out = append(out, row.TaskID)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// the acceptance ruler
// ---------------------------------------------------------------------------

// fk378Violation is one row of PRAGMA foreign_key_check. For a tool_call row the
// owning task and tool are read back too, so a failure names the call and not
// just a rowid (the ticket's reverse judgment asks for rows named one by one).
type fk378Violation struct {
	Table  string
	RowID  int64
	Parent string
	FKID   int
	TaskID string
	Tool   string
}

func (v fk378Violation) String() string {
	if v.Table == "tool_call" && v.TaskID != "" {
		return fmt.Sprintf("%s rowid=%d -> %s.id fk#%d：task_id=%s tool=%s（这发工具调用没有父任务行）",
			v.Table, v.RowID, v.Parent, v.FKID, v.TaskID, v.Tool)
	}
	return fmt.Sprintf("%s rowid=%d -> %s.id fk#%d", v.Table, v.RowID, v.Parent, v.FKID)
}

// fk378ForeignKeys runs the ruler on a SECOND connection to the run's own
// database: PRAGMA foreign_key_check answers with enforcement OFF as well
// (measured by ticket 366-a1 §1; memory's business connections never turn it on),
// so the reading is about the rows, not about a switch. The Store itself exposes
// no query surface, and reading through a connection the writer never touched is
// the point - the ruler must not inherit the writer's assumptions.
func fk378ForeignKeys(t *testing.T, dir string) []fk378Violation {
	t.Helper()
	// D35 fixes the file name; memory keeps it unexported, so this is the one
	// place the name is repeated in this package (the same direct open appears in
	// resident_task_source_246_windows_test.go:355).
	dbPath := filepath.Join(dir, "wisp.db")
	db, err := sql.Open(memory.DriverName, "file:"+filepath.ToSlash(dbPath)+"?_pragma=busy_timeout(3000)")
	if err != nil {
		t.Fatalf("打开 %s 读 foreign_key_check：%v", dbPath, err)
	}
	defer func() { _ = db.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "PRAGMA foreign_key_check")
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_check: %v", err)
	}
	defer func() { _ = rows.Close() }()
	var out []fk378Violation
	for rows.Next() {
		var v fk378Violation
		if err := rows.Scan(&v.Table, &v.RowID, &v.Parent, &v.FKID); err != nil {
			t.Fatalf("读 foreign_key_check 行：%v", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("遍历 foreign_key_check：%v", err)
	}
	// Second pass (the result set must be closed before a nested query on the
	// same connection pool): name the task and tool behind every tool_call row,
	// and fail loudly if a violation cannot be attributed - an unattributable
	// orphan is not "probably fine".
	for i := range out {
		if out[i].Table != "tool_call" {
			continue
		}
		var taskID, tool sql.NullString
		err := db.QueryRowContext(ctx, `SELECT task_id, tool FROM tool_call WHERE rowid = ?`, out[i].RowID).
			Scan(&taskID, &tool)
		if err != nil {
			t.Fatalf("孤儿 tool_call rowid=%d 查不到 task_id/tool：%v", out[i].RowID, err)
		}
		out[i].TaskID, out[i].Tool = taskID.String, tool.String
	}
	return out
}

// ---------------------------------------------------------------------------
// the case
// ---------------------------------------------------------------------------

// Test378SpawnedChildToolCallsHaveAParentRow is the ticket's acceptance reading.
// See the file header for what each assertion discriminates.
func Test378SpawnedChildToolCallsHaveAParentRow(t *testing.T) {
	note := ""
	f := fk378Setup(t, func(dir string) string {
		note = filepath.Join(dir, fk378NoteName)
		return fk378Script(t, []fk378Turn{
			// ① root turn 1: the model asks for the spawn (a real bridge call
			//    under the root's own task id).
			{tool: "task.spawn", id: "call_fk378_spawn", args: map[string]string{
				"description": fk378Label,
				"prompt":      "把数据目录里那份笔记读一遍，然后一句话回报结论。",
			}},
			// ② child turn 1: THE measured call - a real fs.read dispatched by
			//    the child's own loop through the parent's bridge, so the
			//    tool_call row carries the CHILD's task id.
			{tool: "fs.read", id: "call_fk378_child_read", args: map[string]string{
				"path": filepath.ToSlash(note),
			}},
			// ③ child turn 2: its conclusion, i.e. the child task ends here.
			{text: "子任务读过笔记了，没有别的发现。"},
			// ④ root turn 2: the parent closes.
			{text: "父任务收口：子任务回来了。"},
		})
	})
	if err := os.WriteFile(note, []byte("fk378 note body（票 378 用）"), 0o600); err != nil {
		t.Fatal(err)
	}

	code := f.run(fk378Task)
	t.Logf("[378] rc=%d；golden requests=%v", code, f.rep.Requests)
	t.Logf("[378] raw stdout:\n%s", f.out.String())
	t.Logf("[378] raw stderr:\n%s", f.errb.String())
	if code != 0 {
		t.Fatalf("exit %d - 脚本化的一发没跑完，下面的读数会全部错位\nstdout:\n%s\nstderr:\n%s",
			code, f.out.String(), f.errb.String())
	}
	// A pre-condition, not a decoration: this case must be about two rows and
	// not about a 31s confirmation window that happened to end in a reject.
	if strings.Contains(f.out.String(), "[确认") {
		t.Fatalf("这一发出了确认卡（fs.read 在允许目录内应是 L0）：本例量到的行就会是拒绝形状\nstdout:\n%s", f.out.String())
	}

	rootID := f.rootTaskID(t)
	childIDs := f.childIDs(t, rootID)
	if len(childIDs) != 1 {
		t.Fatalf("名册里父任务 %s 下面的子代理行数 = %d（%v），want 1：本例的形状是「一枚子任务、子任务调一次工具」",
			rootID, len(childIDs), childIDs)
	}
	childID := childIDs[0]

	// THE RULER FIRST, before any row assertion can turn into a Fatal: on the
	// reverse run (the publication line deleted) the reading that must be printed
	// is exactly this one, and it must name the orphan rows one by one.
	violations := fk378ForeignKeys(t, f.dir)
	for _, v := range violations {
		t.Logf("[378] fk 违规行：%s", v)
	}
	t.Logf("[378] PRAGMA foreign_key_check = %d 行（root=%s child=%s）", len(violations), rootID, childID)
	if len(violations) != 0 {
		named := make([]string, 0, len(violations))
		for _, v := range violations {
			named = append(named, v.String())
		}
		t.Errorf("PRAGMA foreign_key_check = %d 行，want 0：一发 task.spawn 之后子任务的 task_log 父行必须真的存在\n  %s",
			len(violations), strings.Join(named, "\n  "))
	}
	// [C] the positive control, stated on its own so "0" cannot hide behind an
	// empty database: a violation naming the ROOT task would be ticket 373
	// regressing, which this ticket must not do.
	for _, v := range violations {
		if v.TaskID == rootID {
			t.Errorf("[C] 正控破了：根任务 %s 自己的 tool_call 行成了孤儿 - %s（root 那一支是票 373 治好的）", rootID, v)
		}
	}

	st, err := memory.Open(f.dir)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// [B] the child's own row: the parent every one of its tool_call rows names.
	childRow, err := st.TaskLogByID(ctx, childID)
	if err != nil {
		t.Errorf("[B] 子任务 %s 在 task_log 里没有行（%v）：那样上面那句 fkc=0 只是「没有子调用」的 0", childID, err)
	} else {
		// Ticket 405: the child's row MUST reach a terminal state. It used to stay
		// "running" forever (finishTask was inert in taskLifecycleJournal and the
		// CLI only closed the ROOT), which crash_recovery later mislabelled
		// "interrupted" even on success. Now the loop's own finishTask closes each
		// loop's row, so a completed child is "done".
		t.Logf("[378/405] 子任务行：id=%s state=%s query=%q", childRow.ID, childRow.State, childRow.QueryText)
		if childRow.QueryText == "" {
			t.Errorf("[B] 子任务 %s 的 task_log 行没有 query_text：这行不是「子任务自己的任务书」", childID)
		}
		if childRow.State != string(memory.TaskStateDone) {
			t.Errorf("[405] 子任务 %s 的 task_log 行 state=%q，want %q：子任务行必须被关成终态（否则恒 running、crash_recovery 会把成功的子任务误标 interrupted）",
				childID, childRow.State, string(memory.TaskStateDone))
		}
	}
	// Exactly one row for that id: task_log's id is its primary key, so a second
	// row would mean something rewrote the child's identity.
	taskRows, err := st.ListTaskLogs(ctx)
	if err != nil {
		t.Fatalf("ListTaskLogs: %v", err)
	}
	childRowCount := 0
	for _, tl := range taskRows {
		if tl.ID == childID {
			childRowCount++
		}
	}
	if childRowCount != 1 {
		t.Errorf("[B] 子任务 %s 在 task_log 里的行数 = %d，want 1：每个子任务恰好一行（不许每发工具调用另写一行）",
			childID, childRowCount)
	}
	t.Logf("[378] task_log 共 %d 行（root=%d child=%d）", len(taskRows), countTaskRows378(taskRows, rootID), childRowCount)

	// [A] the child's tool_call rows: they must exist, so [B]+fkc=0 is about a
	// real call and not about a child that never touched a tool.
	childCalls, err := st.ListToolCallsByTask(ctx, childID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask(%s): %v", childID, err)
	}
	if len(childCalls) != 1 {
		names := make([]string, 0, len(childCalls))
		for _, c := range childCalls {
			names = append(names, c.Tool+"/"+c.Outcome)
		}
		t.Errorf("[A] 子任务 %s 的 tool_call 行 = %d（%v），want 1（fs.read）：子任务没调工具的话，fkc=0 读不出任何东西",
			childID, len(childCalls), names)
	} else {
		c := childCalls[0]
		t.Logf("[378] 子任务 tool_call 行：tool=%s level=%s decision=%s outcome=%s", c.Tool, c.RiskLevel, c.Decision, c.Outcome)
		if c.Tool != "fs.read" || c.Outcome != "success" {
			t.Errorf("[A] 子任务的这一行是 %s/%s，want fs.read/success：被拒或没跑成的调用，fkc 照样是 0 但什么都不证明",
				c.Tool, c.Outcome)
		}
	}

	// [C] the root's own two rows (its task.spawn call and its task_log row) must
	// be there too - the control that the ruler looked at a populated database.
	rootCalls, err := st.ListToolCallsByTask(ctx, rootID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask(%s): %v", rootID, err)
	}
	if len(rootCalls) != 1 || rootCalls[0].Tool != "task.spawn" {
		names := make([]string, 0, len(rootCalls))
		for _, c := range rootCalls {
			names = append(names, c.Tool)
		}
		t.Errorf("[C] 根任务 %s 的 tool_call 行 = %v，want 恰好一枚 task.spawn：正控是空的话，fkc=0 不含根那一支",
			rootID, names)
	}
	if _, err := st.TaskLogByID(ctx, rootID); err != nil {
		t.Errorf("[C] 根任务 %s 自己的 task_log 行也读不到（%v）：票 373 的那半不许被本票弄坏", rootID, err)
	}

	// [D] the quantity: one row per tool call, no second row from the inherited
	// journal (ticket 378 hard reminder #1).
	allCalls, err := st.ListToolCalls(ctx)
	if err != nil {
		t.Fatalf("ListToolCalls: %v", err)
	}
	named := make([]string, 0, len(allCalls))
	for _, c := range allCalls {
		named = append(named, c.TaskID+"/"+c.Tool)
	}
	if len(allCalls) != 2 {
		t.Errorf("[D] tool_call 总行数 = %d（%v），want 2（父的 task.spawn + 子任务的 fs.read）：每发工具调用只该有一行",
			len(allCalls), named)
	}
}

// countTaskRows378 counts the task_log rows carrying one id.
func countTaskRows378(rows []memory.TaskLog, id string) int {
	n := 0
	for _, tl := range rows {
		if tl.ID == id {
			n++
		}
	}
	return n
}
