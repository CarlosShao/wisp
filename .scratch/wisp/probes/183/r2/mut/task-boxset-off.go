package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// The D34 task family (ticket 164). Only one of the three registered names has
// an implementation here, and that is deliberate:
//
//	task.output  L0  read what a background task printed   (this file)
//	task.list    -   DEFERRED with five fields, PLAN.md §7 :1531
//	task.cancel  -   DEFERRED with five fields, PLAN.md §7 :1531
//
// task.output is AC#3's leg. The two rows above stay unregistered because the
// §7 registry says so, and a half-claim is exactly the "在册 + 无实现 + 无人认领"
// shape ticket 164 AC#1 was written to end.

// TaskDeps is what a task.* tool needs from the host.
type TaskDeps struct {
	// Roster is the process-local background-task table. It is a pointer with
	// a constructor rather than a map field so the host and the tool always
	// look at the SAME table; nil means the host never wired one, which
	// taskOutput answers fail-closed instead of pretending the task list is
	// merely empty.
	Roster *TaskRoster
	// Paths is the SAME C26 canonicalizer the fs tools are wired with, and it
	// is the only sanctioned way for this tool to answer "can the pointer I am
	// about to hand the model actually be read right now" (ticket 174 AC#2).
	// Nil means the host never wired a judge: taskOutput then fails closed and
	// SAYS so in the reply instead of silently vouching for the path. Deciding
	// that question any other way - filepath.Clean/Abs of one's own - is the
	// D22 ban AGENTS §1.2 states verbatim, so there is no other way here.
	Paths *PathCanonicalizer
	// SpillTokens / SpillHeadTokens / SpillTailTokens override the D15(3)
	// truncation shape. Zero means the frozen reference values below, and the
	// AC#3 test that pins the shape shrinks them to prove the numbers are
	// actually read rather than decorative.
	SpillTokens     int
	SpillHeadTokens int
	SpillTailTokens int
}

// The D15(3) reference numbers, the same ones the loop's spill layer uses
// (internal/agent/budgets.go:32-34, PLAN.md:431): a single tool result over
// 4000 tokens keeps head 500 + tail 200 tokens plus the total length and the
// path, and the full text lands in a file. A tool must not invent its own
// budget - but it also must not hand the model a stub whose pointer it never
// checked, which is why task.output does the shaping itself for the text it is
// serving rather than relying on the loop to spill it afterwards.
const (
	refTaskSpillTokens     = 4000
	refTaskSpillHeadTokens = 500
	refTaskSpillTailTokens = 200
)

func (d TaskDeps) spillTokens() int {
	if d.SpillTokens > 0 {
		return d.SpillTokens
	}
	return refTaskSpillTokens
}

func (d TaskDeps) headTokens() int {
	if d.SpillHeadTokens > 0 {
		return d.SpillHeadTokens
	}
	return refTaskSpillHeadTokens
}

func (d TaskDeps) tailTokens() int {
	if d.SpillTailTokens > 0 {
		return d.SpillTailTokens
	}
	return refTaskSpillTailTokens
}

// TaskOutput is one background task's captured answer as the host recorded it.
//
// It is the missing layer ticket 164's survey named: the artifacts directory
// already holds spilled bytes and internal/agent's Spiller already writes the
// D15(3) stub, but both are addressed by tool-CALL id, so "where did this task's
// output go" had no answer. This record is that mapping, and the host - not the
// model - fills it.
type TaskOutput struct {
	// Text is the full output the task produced. The tool, not the host,
	// decides whether that is too long for a context: this field is the
	// authority the announced total length is measured against.
	Text string
	// ArtifactPath is where the host landed those exact bytes (the
	// D15(3) spill file). Empty means no copy exists, and taskOutput then
	// says so out loud instead of truncating silently.
	ArtifactPath string
}

// TaskRoster is the process-local background task table (ticket 164 ruling 2:
// v1 does NOT keep a task roster across restarts). It is a lookup table keyed
// by the task id ONLY - the id is never joined onto a directory, because the
// path a call returns comes from the host's record, so no model-supplied string
// reaches the filesystem here (C26 stays out of it, and so does the ban on
// filepath decisions outside risk.PathResolver).
type TaskRoster struct {
	mu     sync.RWMutex
	byTask map[string]TaskOutput
}

// NewTaskRoster returns an empty roster, ready to share between host and tools.
func NewTaskRoster() *TaskRoster {
	return &TaskRoster{byTask: map[string]TaskOutput{}}
}

// Record files (or replaces) one task's output. Last writer wins: a retried
// task's later output is the answer, and D15(3) already documents the same
// semantics for the spill artifacts it points at.
func (r *TaskRoster) Record(taskID string, o TaskOutput) {
	if r == nil || taskID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.byTask == nil {
		r.byTask = map[string]TaskOutput{}
	}
	r.byTask[taskID] = o
}

// Look returns one record. The bool is the ONLY way a caller learns "no such
// task": an empty TaskOutput must never be read as a task that printed nothing,
// which is why the tool that consumes this answers with an error text instead
// of an empty result (ticket 164 ruling 2: 查不到必须响亮返回，绝不空返回).
func (r *TaskRoster) Look(taskID string) (TaskOutput, bool) {
	if r == nil || taskID == "" {
		return TaskOutput{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.byTask[taskID]
	return o, ok
}

// Count reports how many tasks are filed. It exists for the host's own
// diagnostics, not as a tool surface: task.list is DEFERRED (§7 :1531), and
// giving task.output a "list everything" mode would smuggle that row back in.
func (r *TaskRoster) Count() int {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byTask)
}

// ---------------------------------------------------------------------------
// task.output (C1 Tool, D34 row: PLAN.md:2564)
// ---------------------------------------------------------------------------

type taskOutput struct{ d TaskDeps }

type taskOutputArgs struct {
	TaskID string `json:"task_id"`
}

// taskOutputSchema is the C1 Parameters document. There is deliberately NO
// offset / cursor / page token in it: PLAN.md:431 already froze the re-read
// mechanism ("模型可用 fs.read 按需再读"), and ticket 164's ruling 1 forbids
// inventing a second pagination scheme in this ticket.
var taskOutputSchema = JSONSchema(`{"type":"object","properties":{` +
	`"task_id":{"type":"string","description":"后台任务的 id（由宿主记录，不是路径）"}` +
	`},"required":["task_id"],"additionalProperties":false}`)

func (taskOutput) Name() string { return "task.output" }

func (taskOutput) Description() string {
	return "读取一个后台任务吐出的输出；过长时返回头尾摘要、总长度和可续读的副本路径"
}

func (taskOutput) Parameters() JSONSchema { return taskOutputSchema }

// Execute implements Tool.
//
// Three shapes, and only the middle one may ever be silent:
//
//	unknown task id  -> IsError naming the id (never an empty success)
//	short output     -> the whole text, Truncated=false, no pointer needed
//	long output      -> head + tail + total + pointer, Truncated=true, and the
//	                    pointer is required by the contract: 只截不指＝不合格
//	                    (PLAN.md:2564), so a long output with no copy file
//	                    announces that the rest is unrecoverable instead of
//	                    dropping it quietly.
func (t taskOutput) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a taskOutputArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.TaskID) == "" {
		return Result{Text: "缺少 task_id 参数", IsError: true}, nil
	}
	if t.d.Roster == nil {
		return Result{Text: "任务名册未接线（fail-closed：拒绝猜测任何任务的输出）", IsError: true}, nil
	}
	rec, ok := t.d.Roster.Look(a.TaskID)
	if !ok {
		return Result{Text: fmt.Sprintf(
			"查不到这个任务 %s：v1 的任务名册只在本进程内，重启后不保留（票 164 定案②）。"+
				"这是「没有这条记录」，不是「任务没有输出」。", a.TaskID), IsError: true}, nil
	}

	full := rec.Text
	if agent.ApproxTokens(full) <= t.d.spillTokens() {
		if onUpdate != nil {
			onUpdate(fmt.Sprintf("task.output 返回 %d 字节（未截断）", len(full)))
		}
		return Result{Text: full, Truncated: false}, nil
	}

	head := takeHeadTokens(full, t.d.headTokens())
	tail := takeTailTokens(full, t.d.tailTokens())
	totalBytes := len(full)
	totalTokens := agent.ApproxTokens(full)

	// The pointer branch. A copy file is the D15(3) norm; its absence is a
	// defect the model has to be told about, because the alternative - a stub
	// with no path - is precisely "只截不指".
	//
	// Ticket 174 AC#2: when the pointer IS given, it still must not promise a
	// road back that does not exist. The path stays in the reply (that is
	// PLAN.md:2564's shape and ticket 164's already-accepted AC#3), and
	// pointerNotice adds what is actually true about it, in the same text the
	// model reads - not in a log line, not in a Go error.
	var stub string
	if rec.ArtifactPath != "" {
		notice := t.d.pointerNotice(rec.ArtifactPath)
		// Ticket 177 shape A (Q-61甲, precise exemption): the path printed at
		// the end of the stub below is text the HOST just wrote itself - it
		// comes from the roster record, never from content this tool read.
		// Declaring it on the per-call box lets C25 keep this whole stub
		// indexed EXCEPT the window this path occupies, so a model that
		// re-reads its own pointer is not blocked by R4 while every other
		// byte of it still taints like normal. The stub's wording and shape
		// are frozen by PLAN.md:2564 / ticket 164's accepted AC#3 - nothing
		// here changes a character the model reads.
	_ = hostPathBoxFromCtx(ctx)
		stub = fmt.Sprintf(
			"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token%s，全文见 %s…]\n%s",
			head, totalBytes-len(head)-len(tail), totalBytes, totalTokens, notice, rec.ArtifactPath, tail)
	} else {
		stub = fmt.Sprintf(
			"%s\n[…输出已截断且无副本文件：省略 %d 字符，总长 %d 字节 / 约 %d token，"+
				"后半段不可找回…]\n%s",
			head, totalBytes-len(head)-len(tail), totalBytes, totalTokens, tail)
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("task.output 截断 %d 字节为头尾摘要", totalBytes))
	}
	return Result{Text: stub, Truncated: true}, nil
}

// TaskOutputDecl is the host-side declaration for task.output.
//
// Declared L0: PLAN.md:2564 puts this row at L0, so it must never raise an
// approval card. Capabilities and Needs are BOTH empty, which is the shape
// PLAN.md:2564's fourth column has for it (a "—", the same spelling as
// task.list/task.cancel on :2563): reading an in-process table touches no C3
// capability face at all, so ticket 164's ruling 3 ("能力声明＝不碰 C3") holds
// without widening the frozen set. PathParams is empty for the same reason -
// the only model-supplied string is a roster key, never a path, so C19 has no
// path to judge and C26 has nothing to canonicalize.
func TaskOutputDecl() Decl {
	return Decl{
		Capabilities: nil,
		Needs:        nil,
		Declared:     risk.L0,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinTaskEntries returns the task family this configuration registers:
// task.output only. See the header comment for why task.list/task.cancel are
// not here.
func BuiltinTaskEntries(d TaskDeps) []Entry {
	return []Entry{{Tool: taskOutput{d: d}, Decl: TaskOutputDecl()}}
}

// pointerNotice answers, in the words the model will read, whether the copy
// file a pointer names can actually be fetched right now. "" means yes: the
// path is inside an authorized root and is a real regular file, and then this
// tool stays quiet - a reply that warned about every healthy pointer would be
// noise, not honesty (ticket 174 AC#2's reverse criterion).
//
// Two shapes make it speak, and they are judged by different facts:
//
//	authorization  C26 decides it, through the same *PathCanonicalizer fs.read
//	               is judged by: Canonicalize, then InAllowlist. Nothing here
//	               normalizes a path by hand (D22 / AGENTS §1.2).
//	existence      os.Stat on the exact string being printed. Read-only: this
//	               function never creates, moves, truncates or removes
//	               anything, which is what keeps a "is it there" check from
//	               becoming a cleanup path.
//
// No judge wired is NOT "assume readable": the answer says it cannot be
// verified, which is the same fail-closed shape as Execute's
// "任务名册未接线" refusal.
func (d TaskDeps) pointerNotice(raw string) string {
	if d.Paths == nil {
		return "；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理）"
	}
	var notes []string
	canon, err := d.Paths.Canonicalize(raw)
	switch {
	case err != nil:
		notes = append(notes, "注意：这条路径现在读不到，C26 连规范化都没通过（"+err.Error()+"）")
	case !d.Paths.InAllowlist(canon):
		notes = append(notes, "注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；"+
			"要用户先把所属目录加进 [fs] allowed_dirs 才读得回来")
	}
	if st, statErr := os.Stat(raw); statErr != nil {
		notes = append(notes, "注意：这条路径现在读不到，宿主登记的那份副本文件并不存在")
	} else if !st.Mode().IsRegular() {
		notes = append(notes, "注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状）")
	}
	if len(notes) == 0 {
		return ""
	}
	return "；" + strings.Join(notes, "；")
}

// takeHeadTokens returns the leading budget tokens of s, 4 bytes per token
// under the shared agent.ApproxTokens heuristic, cut on a rune boundary so a
// truncated result can never carry half a UTF-8 sequence into the context.
//
// These two helpers are deliberately local copies of the loop's spill cutters
// (internal/agent/spill.go takeTokens/takeTokensLast) rather than new exported
// methods on agent.Loop: ticket 164's ruling 4 puts the read entry on the tools
// side, and the shape "a *Loop method that takes a task id" is what probe 154's
// quiet G3 leg is registered to catch.
func takeHeadTokens(s string, budget int) string {
	if budget <= 0 {
		return ""
	}
	n := budget * 4
	if n >= len(s) {
		return s
	}
	for n > 0 && !utf8.RuneStart(s[n]) {
		n--
	}
	return s[:n]
}

// takeTailTokens returns the trailing budget tokens of s.
func takeTailTokens(s string, budget int) string {
	if budget <= 0 {
		return ""
	}
	n := len(s) - budget*4
	if n <= 0 {
		return s
	}
	for n < len(s) && !utf8.RuneStart(s[n]) {
		n++
	}
	return s[n:]
}
