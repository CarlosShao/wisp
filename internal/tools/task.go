package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// The D34 task family (ticket 164, plus ticket 221 甲形 per ruling A434, plus
// ticket 47 AC#6). All three names the roster line carries now have an
// implementation here:
//
//	task.output  L0  read what a background task printed   (this file)
//	task.cancel  L1  stop one subagent the CALLER itself derived (this file, 221)
//	task.list    L0  read the scheduler's task rows         (this file, 47)
//
// task.output is AC#3's leg, task.cancel is ticket 221's. task.list lands
// with ticket 47: its done-criterion (PLAN.md:1537's registry row) is
// "注册表里 task.list／task.cancel 各有一枚真实现并过契约测试", ticket 163's
// prerequisite already closed the 163/164 half, and the D34 row
// (PLAN.md:2564, L0 for task.list) was in the frozen table from the start -
// so this is the "标记与接线同一次核销" shape ruling A434 set for task.cancel,
// applied to the second half of the same row. The DEFERRED line in the
// header is lifted in THIS commit, with the registration and the tests that
// pin the wiring (tasklist_deferred_236r3_teeth_test.go's roster nail now
// demands the row) in the same landing.
//
// task.cancel's authority shape (only the row's own parent may stop it) is
// UNCHANGED: task.cancel stops what the caller derived. The scheduler's rows
// are reachable through task.list (read) and - for the scheduler's own rows
// only, whose ParentTaskID the host files as the CALLER's task id - through
// task.cancel's scheduler branch below. task.cancel also carries the scheduler
// branch (task.list read + scheduler row stop, ticket 47 AC#6).

// TaskDeps is what a task.* tool needs from the host.
type TaskDeps struct {
	// Roster is the process-local background-task table. It is a pointer with
	// a constructor rather than a map field so the host and the tool always
	// look at the SAME table; nil means the host never wired one, which
	// taskOutput answers fail-closed instead of pretending the task list is
	// merely empty.
	Roster *TaskRoster
	// Scheduler is ticket 47's C20 scheduler face (nil = not wired, the
	// fail-closed shape below): task.list reads the rows through View, and
	// task.cancel's scheduler branch stops a scheduler row through Control.
	// Both faces are narrow by construction (SchedulerView carries no
	// allow/answer verb, SchedulerControl carries exactly Cancel), which is
	// what the 197 allow-door walker demands of anything TaskDeps holds.
	Scheduler SchedulerDeps
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

// SchedulerDeps is ticket 47's two narrow faces onto the scheduler, carried
// as interfaces so internal/tools never imports agent/scheduler (the arrow
// would drag every fs tool into the scheduler's graph, and tools is
// consume-only for packages below it). The host (cmd/wisp's assembly, or a
// test) fills them with adapters over *scheduler.Scheduler. A zero value
// means "no scheduler wired": task.list then answers fail-closed.
type SchedulerDeps struct {
	// View is the read face: rows only, no verbs that answer anything.
	View TaskListView
	// Control is the write face: exactly one verb, Cancel.
	Control TaskCanceler
}

// TaskListView is the read face task.list consumes. Deliberately minimal:
// rows in, nothing else - so the 197 allow-door walker finds no answerer on
// it, and a test double can be three rows wide.
type TaskListView interface {
	// Rows returns the scheduler's task rows (id, label, D43 state, note,
	// held paths) in FIFO accept order.
	Rows() []TaskListRow
}

// TaskCanceler is the write face task.cancel's scheduler branch consumes:
// one verb, the same narrow shape TaskRoster.Cancel has.
type TaskCanceler interface {
	// Cancel stops one scheduler task (queued: dequeue+settle; running:
	// root ctx cancel). An unknown id answers an error - the branch below
	// surfaces it, never a silent no-op.
	Cancel(id string) error
}

// TaskListRow is the scheduler row as the tool family sees it: the D43 state
// vocabulary only (statemachine.State - the same discipline ticket 188 put
// on the roster's status dimension), plus the free-text note where the
// 「在等 <task> 释放 <path>」 reason lives, and the paths the row holds.
// ParentTaskID is the authority join task.cancel's scheduler branch sits on:
// the host files the scheduler row's parent as the CALLER's task id, so the
// branch's "只有父任务能停" check is data, not a convention. Empty on a root
// row (nobody may stop a root through this tool - the user's own stop path
// is the panel, not a model verb).
type TaskListRow struct {
	ID           string
	Label        string
	State        statemachine.State
	Note         string
	ParentTaskID string
	HeldPaths    []string
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
	// State is this task's status dimension, as the host filed it (ticket 188
	// AC#2, ruling A394 (i): the carrier is this struct, not memory.TaskLog).
	//
	// Its type is statemachine.State, whose value set is the verbatim in-repo
	// copy of D43's frozen transition table: internal/statemachine/states.go:11-31
	// - a file whose own header line :8 reads "names exactly as in D43 / SPEC-08" -
	// sourced from docs/PLAN.md:3055 plus the 40 transitions at :3061-3100.
	// Typing the field that way is what makes "no invented state names" a
	// structural fact instead of a convention, and statemachine.Valid (states.go:39)
	// is the judge every writer here goes through.
	//
	// What this field is NOT: it is not memory.TaskLog.State
	// (internal/memory/models.go:62, and its schema at :57 carries no CHECK).
	// That column is already being fed done / cancelled / running / succeeded by
	// internal/agent/loop.go:983-998 and cmd/wisp/run.go:651-670 - four words,
	// none of them on the D43 table. Reconciling the two vocabularies is ticket
	// 196's, and per A394 this ticket does not unify them, does not read that
	// column and does not touch either producer.
	//
	// Empty means the host filed none. That is a real answer and StateAnswer
	// announces it, rather than this record quietly reading as "finished".
	State statemachine.State
	// ParentTaskID is the task that derived this one; empty means this row is a
	// root (ticket 197 §0: one identity dimension, no second status field).
	ParentTaskID string
	// Label is the short title of the work - the task's NAME only, never its
	// body, which stays in Text. It is what a roster row reads as in a list.
	Label string
	// Kind is TaskKindRoot or TaskKindSubagent (see subagent_197.go), and the
	// spawner's depth check reads it: a subagent may not derive another.
	// Empty means the host filed no kind, which is how every pre-197 row reads.
	Kind string
}

// StateAnswer reads the status dimension of one roster record the same way
// pointerNotice reads its pointer below: it never vouches for what it cannot
// check.
//
//	(state, "")      the host filed one of D43's 20 names
//	("", reason)     nothing was filed, or what was filed is not a D43 name
//
// The reason string is what a display surface (the panel snapshot's task
// section, which is ticket 145's carrier and is deliberately NOT built here)
// is supposed to show instead of a verdict. A reader that skipped this
// function and printed its own word would be hardcoding the status - the
// exact shape ticket 188 AC#2's "值来自真名册，不许硬编码" forbids - and a writer
// that bypassed it would be how a fifth legacy word (ticket 196's defect)
// slips into this dimension, which is why both branches below are named and
// loud.
func (o TaskOutput) StateAnswer() (statemachine.State, string) {
	if o.State == "" {
		return "", "宿主没有登记这一维（fail-closed：不替任务编一个状态）"
	}
	if !statemachine.Valid(o.State) {
		return "", fmt.Sprintf(
			"登记的「%s」不是 D43 状态机表里的 20 个名字之一（票 188 AC#2：不许自造；"+
				"memory 侧那套旧词表属票 196，本维不接它）", string(o.State))
	}
	return o.State, ""
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
	// cancel is the ONE stop path a subagent has: the child loop's own
	// context.CancelFunc, handed over at spawn time and read back by Cancel.
	// No second cancel mechanism exists for a roster row by design (197 §0).
	cancel map[string]func()
	// reservedSubagents counts the spawn calls holding a pool slot, i.e. the
	// subagents in flight. It is a counter and not a scan because a slot is
	// reserved BEFORE the child exists (its id is not) and released only when the
	// child has joined, and the rows on both sides of that window would otherwise
	// either double-count or count nothing.
	reservedSubagents int
	// watchers are the rows a reader asked to keep an eye on (ticket 197: the
	// display layer needs "a row appeared / moved" without re-polling a table it
	// cannot snapshot atomically). Sends are non-blocking by design: a slow reader
	// loses an intermediate state, never the row itself, and never holds the lock.
	watchers map[string][]chan TaskOutput
}

// NewTaskRoster returns an empty roster, ready to share between host and tools.
func NewTaskRoster() *TaskRoster {
	return &TaskRoster{
		byTask:   map[string]TaskOutput{},
		cancel:   map[string]func(){},
		watchers: map[string][]chan TaskOutput{},
	}
}

// Record files (or replaces) one task's output. Last writer wins: a retried
// task's later output is the answer, and D15(3) already documents the same
// semantics for the spill artifacts it points at.
//
// The one exception is identity (ticket 197): ParentTaskID / Label / Kind are
// preserved from the row already on file when the incoming record leaves them
// empty, because the host's write side (TaskBackfill) files {Text, State} and
// would otherwise erase the parent link of a row the spawner just published.
// Text, ArtifactPath and State are NOT preserved - an empty State is a real
// answer that StateAnswer announces (ticket 188), not a hole to paper over.
func (r *TaskRoster) Record(taskID string, o TaskOutput) {
	if r == nil || taskID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recordLocked(taskID, o)
}

// recordLocked is Record's body; the caller already holds r.mu. It exists
// because ticket 197's in-flight publish has to write the row and the stop
// handle in ONE critical section (see PublishSubagentInFlight), so it must reuse
// this exact identity-preserving write rather than keep a second copy of it.
func (r *TaskRoster) recordLocked(taskID string, o TaskOutput) {
	if r.byTask == nil {
		r.byTask = map[string]TaskOutput{}
	}
	if prev, ok := r.byTask[taskID]; ok {
		if o.ParentTaskID == "" {
			o.ParentTaskID = prev.ParentTaskID
		}
		if o.Label == "" {
			o.Label = prev.Label
		}
		if o.Kind == "" {
			o.Kind = prev.Kind
		}
	}
	r.byTask[taskID] = o
	for _, ch := range r.watchers[taskID] {
		select {
		case ch <- o:
		default:
		}
	}
}

// WatchRow returns a channel that receives every later write to one row, from
// the moment it is asked for. It is how a reader sees a subagent's row appear
// and move without polling, and it is deliberately a live view only: a row
// written before the call is not replayed (Look is the way to read that).
func (r *TaskRoster) WatchRow(taskID string) <-chan TaskOutput {
	if r == nil || taskID == "" {
		ch := make(chan TaskOutput)
		close(ch)
		return ch
	}
	ch := make(chan TaskOutput, 8)
	r.mu.Lock()
	if r.watchers == nil {
		r.watchers = map[string][]chan TaskOutput{}
	}
	r.watchers[taskID] = append(r.watchers[taskID], ch)
	r.mu.Unlock()
	return ch
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
// diagnostics. The model-facing list is task.list (ticket 47): it reads the
// scheduler's rows, never this roster, so giving task.output a "list
// everything" mode stays out (that would be a second list surface the D34
// row never froze).
func (r *TaskRoster) Count() int {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byTask)
}

// TaskRecord is one roster row with the key it is filed under. The display layer
// needs both (ticket 197 leg B reads the stream by task id), and Look's bool is
// deliberately not smuggled in here: a record that is not in the table simply
// never appears in a list built from it.
type TaskRecord struct {
	TaskID string
	Out    TaskOutput
}

// Descendants returns every row this task derived, directly or through a chain,
// nearest generation first. That read port is what lets the panel open a subagent
// by clicking into its parent (ticket 197 §0: 名册里一行有状态、可查、可停的记录).
// It walks ParentTaskID only; a row with no parent is a root and is never
// anybody's descendant, so a cycle in hand-filed data cannot hang this loop: the
// seen-set below stops it.
func (r *TaskRoster) Descendants(taskID string) []TaskRecord {
	if r == nil || taskID == "" {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []TaskRecord
	frontier := []string{taskID}
	seen := map[string]bool{taskID: true}
	for len(frontier) > 0 {
		next := make([]string, 0, len(frontier))
		for _, parent := range frontier {
			for id, o := range r.byTask {
				if o.ParentTaskID != parent || seen[id] {
					continue
				}
				seen[id] = true
				out = append(out, TaskRecord{TaskID: id, Out: o})
				next = append(next, id)
			}
		}
		frontier = next
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].TaskID == out[j].TaskID {
			return out[i].TaskID < out[j].TaskID
		}
		return out[i].TaskID < out[j].TaskID
	})
	return out
}

// PublishSubagent files (or re-files) the identity row of one subagent as
// in-flight, WITHOUT a stop handle. That is honest only for a row that has no
// handle by construction: a fixture, or the "sibling that nobody may stop" a leg
// files to prove Cancel touches exactly one row.
//
// The production spawn path must NOT use it. A row that reads in flight and
// cannot be stopped is precisely the lie ticket 70 step 7 read out of CI
// (state=Thinking yet 「没有在跑 … 没有可停的句柄」), so the spawner files through
// PublishSubagentInFlight, which writes the row and the handle in one critical
// section. Test197SpawnPathPublishesRowAndStopHandleTogether is the ruler that
// keeps that split; this function's own callers are tests only.
func (r *TaskRoster) PublishSubagent(taskID, parentTaskID, label string) {
	r.Record(taskID, TaskOutput{
		ParentTaskID: parentTaskID, Label: label, Kind: TaskKindSubagent,
		State: statemachine.State(subagentStateRunning),
	})
}

// PublishSubagentInFlight files (or re-files) one in-flight subagent's roster
// row AND attaches its one stop handle in the SAME critical section. This is the
// spawn path's only publisher (ticket 197, rewired by ticket 70 step 7).
//
// Why one call and not two calls: Cancel answers 「没有在跑」 for a known row whose
// handle is nil (Cancel below). With the row published first and the handle
// attached second there is an instant where both facts are true of the same row
// - the row says this child IS running, the refusal says it is NOT - and the only
// way to close it is to make 「行可见 ⇒ 句柄已挂」 a property of the table instead
// of an ordering convention between two call sites. CI caught exactly that
// instant in 2 of 20 runs of Test221SubagentCannotStopSiblingOrItself /
// Test221ParentStopsItsOwnChildRowAndStreamSettle; the local instrument
// (Test197InFlightRowNeverLacksItsStopHandle) reproduces it under scheduler
// contention and reads 0 hits afterwards.
//
// Re-publishing is idempotent and safe: both publishers (the child's admission
// hook, which is the one that matters because it runs before the child's first
// model call, and the spawn goroutine, which keeps "a child that has started has
// a row" true) pass identical values and the SAME cancel handle, so
// last-writer-wins on either map cannot erase the other half.
//
// A nil cancel publishes NOTHING: this call exists to make a row stoppable, and
// a row that cannot be stopped must not become visible as in flight. The two
// call sites in subagent_197.go both pass the child's real cancel handle.
func (r *TaskRoster) PublishSubagentInFlight(taskID, parentTaskID, label string, cancel func()) {
	if r == nil || taskID == "" || cancel == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recordLocked(taskID, TaskOutput{
		ParentTaskID: parentTaskID, Label: label, Kind: TaskKindSubagent,
		State: statemachine.State(subagentStateRunning),
	})
	if r.cancel == nil {
		r.cancel = map[string]func(){}
	}
	r.cancel[taskID] = cancel
}

// MarkRoot files the root identity for a task the host is running itself.
//
// State is deliberately NOT written here: the root task's producer still files
// no D43 name (ruling A394 keeps that rewiring with ticket 196), so
// StateAnswer keeps saying 「宿主没有登记这一维」 for a root row. What lands is
// the parent link (none) and the kind, which is what makes the tree readable.
func (r *TaskRoster) MarkRoot(taskID, label string) {
	r.Record(taskID, TaskOutput{Label: label, Kind: TaskKindRoot})
}

// TryAcquireSubagentSlot reserves one of the shared pool's slots. The pool is
// process-wide: the root and all of its descendants count against the same
// limit, because they share this table. release is idempotent, and every caller
// must reach it exactly once - on the refusal paths and after the child joins.
func (r *TaskRoster) TryAcquireSubagentSlot(limit int) (release func(), ok bool) {
	if r == nil {
		return func() {}, false
	}
	if limit <= 0 {
		limit = 1
	}
	r.mu.Lock()
	if r.reservedSubagents >= limit {
		r.mu.Unlock()
		return func() {}, false
	}
	r.reservedSubagents++
	var once sync.Once
	r.mu.Unlock()
	return func() {
		once.Do(func() {
			r.mu.Lock()
			if r.reservedSubagents > 0 {
				r.reservedSubagents--
			}
			r.mu.Unlock()
		})
	}, true
}

// InFlightSubagents reports how many slots the pool is holding.
func (r *TaskRoster) InFlightSubagents() int {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.reservedSubagents
}

// RunningSubagentIDs names the rows currently in flight. It reads the D43
// dimension, so the ids it returns are the ones a reader can also see in the
// panel; the refusal text uses it to say WHICH subagents are holding the pool.
func (r *TaskRoster) RunningSubagentIDs() []string {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var ids []string
	for id, o := range r.byTask {
		if o.Kind == TaskKindSubagent && o.State == statemachine.State(subagentInFlightState) {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

// AttachCancel stores the child's own context.CancelFunc as the ONLY stop path
// for that row. Re-attaching replaces it (last writer wins, same as Record).
func (r *TaskRoster) AttachCancel(taskID string, cancel func()) {
	if r == nil || taskID == "" || cancel == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel == nil {
		r.cancel = map[string]func(){}
	}
	r.cancel[taskID] = cancel
}

// DetachCancel forgets one row's stop handle once the child has joined: a cancel
// after that would be a no-op pretending to be a stop.
func (r *TaskRoster) DetachCancel(taskID string) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cancel, taskID)
}

// Cancel stops exactly one roster row through the handle it was spawned with, and
// never touches anybody else's - not the parent, not a sibling. The bool plus
// reason pair is the same fail-closed shape as Look's: "no such handle" is a
// different fact from "already stopped", and Cancel says which one it found.
func (r *TaskRoster) Cancel(taskID string) (bool, string) {
	if r == nil || taskID == "" {
		return false, "没有这条任务 id（宿主没登记过它）"
	}
	r.mu.RLock()
	cancel := r.cancel[taskID]
	o, known := r.byTask[taskID]
	r.mu.RUnlock()
	if !known {
		return false, fmt.Sprintf("查不到任务 %s：v1 的任务名册只在本进程内，重启后不保留（票 164 定案②）", taskID)
	}
	if cancel == nil {
		return false, fmt.Sprintf("任务 %s 没有在跑（kind=%s，state=%s），没有可停的句柄",
			taskID, o.Kind, string(o.State))
	}
	cancel()
	return true, ""
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
		if box := hostPathBoxFromCtx(ctx); box != nil {
			box.set(rec.ArtifactPath)
		}
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
// task.output, task.cancel (ticket 221 甲形, ruling A434) and - since ticket
// 47, whose done-criterion is PLAN.md:1537's registry row - task.list. See
// the header comment for the marker/wiring same-landing rule.
func BuiltinTaskEntries(d TaskDeps) []Entry {
	return []Entry{
		{Tool: taskOutput{d: d}, Decl: TaskOutputDecl()},
		{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},
		{Tool: taskList{d: d}, Decl: taskListDecl()},
	}
}

// ---------------------------------------------------------------------------
// task.cancel (C1 Tool, D34 row: PLAN.md:2564, declared L1, ticket 221 甲形)
// ---------------------------------------------------------------------------

type taskCancel struct{ d TaskDeps }

type taskCancelArgs struct {
	TaskID string `json:"task_id"`
}

// taskCancelSchema carries the same single roster key task.output reads, and for
// the same reason no path appears here: the only model-supplied string is a
// roster id, never a filesystem path, so C19 has no path to judge and C26
// nothing to canonicalize.
var taskCancelSchema = JSONSchema(`{"type":"object","properties":{` +
	`"task_id":{"type":"string","description":"要停掉的那枚子代理的任务 id（名册那一行的 id，不是路径）"}` +
	`},"required":["task_id"],"additionalProperties":false}`)

func (taskCancel) Name() string { return "task.cancel" }

// Description says what this tool is allowed to do and, in the same breath, what
// it is not: the boundary is the security half of ticket 221 (A434 批准范围 item
// 2), and a description that only named the capability would hand the model a
// power it cannot use. It also keeps telling the model the half ticket 197 AC#6
// made it owe: 父取消不级联 - and the reverse, stopping a child never reaches
// anybody else either.
func (taskCancel) Description() string {
	return fmt.Sprintf("停掉你自己派出的那一枚子代理：只有派生它的父任务能停它，"+
		"子代理停兄弟、停自己都一律被拒（深度上限 %d 的树里，别人才不是你的孩子）；"+
		"停掉父任务不会级联停掉子代理，停掉一枚子代理也只动它那一行——"+
		"那一行不会消失，它会落到 D43 已有的状态名上",
		MaxSubagentDepth)
}

func (taskCancel) Parameters() JSONSchema { return taskCancelSchema }

// taskCancelDecl declares the row.
//
// Declared L1 because that is D34's own number and not this ticket's to move:
// PLAN.md:2564 reads "task.list / task.cancel | 查看/取消任务 | L0 / L1", the
// pair splitting into task.list L0 and task.cancel L1, and SPEC-07 §3's mirror
// row repeats it. ⚠ Consequence worth naming: with a gate that cannot answer an
// L1 window (tools.NoGate, i.e. 票 21 未接入), a task.cancel call is REFUSED
// before it reaches Execute - which is the honest posture of "取消是可逆写但要有
// 一枚执行前窗口" (SPEC-06 §2), not a bug this ticket papers over by declaring
// L0. cmd/wisp wires approval.Gate, so the production path gets the window.
//
// Capabilities and Needs are empty the way task.output's are (the frozen row's
// fourth column is a "—"): stopping a row in the host's own in-process table
// touches no C3 capability face. PathParams is empty because the only
// model-supplied string is a roster key.
func taskCancelDecl() Decl {
	return Decl{
		Capabilities: nil,
		Needs:        nil,
		Declared:     risk.L1,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// Execute implements Tool: judge the caller, then stop exactly that one row.
//
// Five answers, and the refusal shapes are the point of this tool rather than
// its edges (A434 批准范围 item 2, same family as ticket 197 AC#5's 子级永不自批):
//
//	who am I?          no host-minted task id on this call -> refuse, fail closed
//	whose row?         caller is not the row's parent       -> refuse (a subagent
//	                                                    stopping a SIBLING lives
//	                                                    on this branch)
//	my own row?        caller == target                      -> refuse (SELF-STOP)
//	root row?          the row has no parent                 -> refuse: nobody in
//	                    this tree may stop a root but its own caller, and the
//	                    "user stops it" half has no landing yet (票 181／票 220)
//	one of mine?       stop it through the handle IT was spawned with
//
// The stop itself is TaskRoster.Cancel - the function this repository already
// had and nobody called (ticket 221's whole complaint). This tool is its first
// production caller, and it is the ONLY one: no second stop path gets invented.
func (t taskCancel) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a taskCancelArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	target := strings.TrimSpace(a.TaskID)
	if target == "" {
		return Result{Text: "缺少 task_id 参数", IsError: true}, nil
	}
	if t.d.Roster == nil {
		return Result{Text: "任务名册未接线（fail-closed：拒绝停掉任何任务——名册才是唯一的停法）", IsError: true}, nil
	}
	// The caller's identity is the host's, not an argument: the bridge stamps
	// the dispatching TASK's id into the call context (TaskID; the correlation
	// id went per-call in ticket 242 and no longer names a task), so a model
	// cannot name whose child it is by editing parameters. Same source
	// task.spawn uses to fill in ParentTaskID, which is what makes the two
	// sides agree by construction.
	caller := TaskID(ctx)
	if caller == "" {
		return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）", IsError: true}, nil
	}
	// --- ticket 47 AC#6: the scheduler branch --------------------------------
	//
	// A row filed by the C20 scheduler is NOT a roster row: it has no
	// ParentTaskID in this roster (Look misses it), so without this branch a
	// model could see the row through task.list and have task.cancel answer
	// 「查不到这个任务」 - the promise/ capability fork ticket 221 was filed to
	// kill, one row over. The branch sits BEFORE the roster read and keeps
	// the same authority shape: only the row's own parent (the host files the
	// scheduler row's parent as the CALLER's task id) may stop it, self-stop
	// and foreign-stop are refused, and an unknown id is answered loudly.
	if t.d.Scheduler.Control != nil && t.d.Scheduler.View != nil {
		if row, known := schedulerLookup(t.d.Scheduler.View, target); known {
			if target == caller {
				return Result{
					Text: fmt.Sprintf(
						"拒绝停止：%s 就是这次调用所属的任务自己（调度器行）。"+
							"子代理不许停自己（票 221 AC#3 的正控）；要收摊就把手上的活做完。", target),
					IsError: true,
				}, nil
			}
			if row.ParentTaskID != caller {
				return Result{Text: fmt.Sprintf(
					"拒绝停止：%s 是调度器任务 %s 派生的孩子，不是调用者 %s 的——只有父任务能停自己的那行"+
						"（票 221 AC#3 的正控，调度器行同规）。用 task.list 看的是整张队列，能停的只有你自己派的那行。",
					target, row.ParentTaskID, caller), IsError: true}, nil
			}
			if err := t.d.Scheduler.Control.Cancel(target); err != nil {
				return Result{Text: fmt.Sprintf(
					"没能停掉调度器任务 %s：%v（这一行没有被改动，别的行也没有被牵连）", target, err), IsError: true}, nil
			}
			if onUpdate != nil {
				onUpdate(fmt.Sprintf("task.cancel 已对调度器任务 %s 发出停止", target))
			}
			return Result{Text: fmt.Sprintf(
				"已对调度器任务 %s（%s）发出停止：排队中的那一行就地结算，在跑的那一行走任务根 ctx 取消（D38c）。"+
					"它落地后的 applied-steps 台账留在它自己的行上，用 task.list 再看状态；"+
					"这次停止只动了这一行，队列里别的任务没有受影响。", target, row.Label)}, nil
		}
		// The target is not a scheduler row either: fall through to the
		// roster branch (a subagent row), which answers 「查不到」 for ids
		// neither table filed.
	}
	rec, known := t.d.Roster.Look(target)
	if !known {
		return Result{Text: fmt.Sprintf(
			"查不到这个任务 %s：v1 的任务名册只在本进程内，重启后不保留（票 164 定案②）。"+
				"这是「没有这条记录」，不是「它已经停了」。", target), IsError: true}, nil
	}
	if target == caller {
		return Result{
			Text: fmt.Sprintf(
				"拒绝停止：%s 就是这次调用所属的任务自己。子代理不许停自己（票 221 AC#3 的正控，"+
					"票 197 AC#5「子级永不自批」同族）；要收摊就把手上的活做完，或者把这一行留给父任务。", target),
			IsError: true,
		}, nil
	}
	if rec.ParentTaskID == "" {
		return Result{Text: fmt.Sprintf(
			"拒绝停止：%s 那一行没有父任务（kind=%s），它不是这次调用派生的孩子，本工具动不了它。"+
				"由用户停一棵正在跑的树是另一条入向通道（票 181／票 220），今天没有落点。",
			target, orDefaultKind(rec.Kind)), IsError: true}, nil
	}
	if rec.ParentTaskID != caller {
		return Result{Text: fmt.Sprintf(
			"拒绝停止：%s 是任务 %s 派生的孩子，不是调用者 %s 的孩子——只有父任务能停自己的孩子，"+
				"兄弟之间、旁支之间都停不了（票 221 AC#3 的正控）。",
			target, rec.ParentTaskID, caller), IsError: true}, nil
	}

	// Subscribe before signalling: the row's terminal write is what this call
	// reports, and a watcher attached after the cancel could miss it.
	watch := t.d.Roster.WatchRow(target)
	stopped, why := t.d.Roster.Cancel(target)
	if !stopped {
		return Result{
			Text:    fmt.Sprintf("没能停掉 %s：%s（这一行没有被改动，别的孩子也没被牵连）", target, why),
			IsError: true,
		}, nil
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("task.cancel 已对子代理 %s 发出停止", target))
	}

	final, settled := awaitSettledRow(ctx, watch)
	if settled {
		if st, readWhy := final.StateAnswer(); st == "" {
			return Result{Text: fmt.Sprintf(
				"已停掉子代理 %s（%s），但那一行的状态维读不回来：%s", target, final.Label, readWhy), IsError: true}, nil
		} else {
			if onUpdate != nil {
				onUpdate(fmt.Sprintf("task.cancel 收尾：子代理 %s 落到 %s", target, string(st)))
			}
			return Result{
				Text: fmt.Sprintf(
					"已停掉子代理 %s（%s）。只动了这一行：%s 的父任务仍是 %s，那一行没有消失，"+
						"它的流已经收到终态，状态落在 D43 已有的名「%s」（被宿主停掉这一支借的是 Muted，本票不新造态名）。"+
						"停掉父任务不会级联停掉子代理，反过来也一样：这次停止没有碰它的任何兄弟，也没有碰 %s。"+
						"这次停止由 %s 发起，目标是 %s。",
					target, final.Label, target, rec.ParentTaskID, string(st), rec.ParentTaskID, caller, target),
			}, nil
		}
	}
	// The signal went out and the row simply has not reported back inside this
	// call's own budget. That is a real, sayable state - what must NOT happen is
	// this reply claiming a terminal state the roster has not filed, or leaving
	// the model to assume the child is still running.
	return Result{Text: fmt.Sprintf(
		"停止信号已经通过 %s 自己的句柄发出（没有第二套停法），但这一行在这次调用的预算内还没有回报终态（%v）。"+
			"派生它的那次 task.spawn 的收尾监视器仍会替它落账，那一行不会凭空消失，也不会被级联波及别的行。"+
			"现在能核实的是：它仍挂在父任务 %s 名下，停它的调用者是 %s。",
		target, ctx.Err(), rec.ParentTaskID, caller), IsError: true}, nil
}

// awaitSettledRow waits for one roster row to leave the in-flight D43 name and
// returns the record that moved. The bound is the CALL's own context (the C22
// deadline the bridge already carries), never a wall-clock difference - D22's
// ban on Sub(time.Now()) timeouts is exactly the shortcut this avoids. A closed
// channel and a dead context both answer false, which the caller says out loud
// instead of inventing a state.
func awaitSettledRow(ctx context.Context, watch <-chan TaskOutput) (TaskOutput, bool) {
	for {
		select {
		case o, open := <-watch:
			if !open {
				return TaskOutput{}, false
			}
			if o.State != "" && o.State != statemachine.State(subagentStateRunning) {
				return o, true
			}
		case <-ctx.Done():
			return TaskOutput{}, false
		}
	}
}

func orDefaultKind(kind string) string {
	if kind == "" {
		return "未登记"
	}
	return kind
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
		// The step is named and the sentence stops. The upstream text is NOT
		// spliced in: this arm was the one place err.Error() reached the
		// model-visible reply, and the relay is unbounded - any future C26
		// error text would ride it into the context (174-v1/A375; 174-a3
		// form 乙, struck by A587 as one line with zero criteria moved). The
		// sentence below is the one internal/agent/spill.go already speaks for
		// this same failure, so both legs now say the same thing.
		notes = append(notes, "注意：这条路径现在读不到，C26 连规范化都没通过")
	case d.Paths.DenyA(canon):
		notes = append(notes, "注意：这条路径现在绝对读不到——它命中了 A 档敏感路径（绝对禁止、任何授权都不可豁免，不会有卡，加 [fs] allowed_dirs 或批卡都无效）")
	case !d.Paths.InAllowlist(canon):
		// Both roads back are named, because both exist. [fs] allowed_dirs is a
		// risk-TIER input, not an execution wall: FSDeps.open canonicalizes and
		// never re-checks the allowlist, so one L2 card answered yes really does
		// fetch these bytes (the criterion in
		// task_output_pointer_two_roads_174r1_test.go measures that on the same
		// bridge). The previous sentence said flatly "fs.read 会被拒 ... 才读得
		// 回来", which is true only of the nobody-approves case and made this leg
		// contradict internal/agent/spill.go, which has told the model both roads
		// for this identical verdict since 174-c2 measured the asymmetry. One
		// promise, two legs: they now say the same thing.
		notes = append(notes, "注意：这条路径现在读不到，它不在你被授权的目录范围内："+
			"fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——"+
			"要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡")
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
	if n < 0 {
		n = 0
	}
	if n > len(s) {
		return s
	}
	for n < len(s) && !utf8.RuneStart(s[n]) {
		n++
	}
	return s[n:]
}

// ---------------------------------------------------------------------------
// task.list (C1 Tool, D34 row: PLAN.md:2564, declared L0, ticket 47 AC#6)
// ---------------------------------------------------------------------------

type taskList struct{ d TaskDeps }

// taskListSchema carries no argument at all: the list is the whole state the
// model may see, filtered by nothing (the D31 property is "被阻塞的东西必须
// 可见" - a model that must guess a filter first would not see it).
var taskListSchema = JSONSchema(`{"type":"object","properties":{},"additionalProperties":false}`)

func (taskList) Name() string { return "task.list" }

// Description says what the list is and what it is not: the scheduler's own
// rows (queue order, D43 state, wait reason), not the subagent roster -
// that one is reachable through task.output's task_id read.
func (taskList) Description() string {
	return "列出调度器当前的任务行（排队顺序）：每一行带任务 id、标题、D43 状态名" +
		"（Queued=在排队，含「在等谁释放什么路径」的原因、Thinking/Acting/AwaitingApproval=在跑、" +
		"Settling/Error=已结束）和它占着的路径。查某一行吐过什么用 task.output；" +
		"停你自己派生的子代理用 task.cancel。"
}

func (taskList) Parameters() JSONSchema { return taskListSchema }

// taskListDecl declares the row. L0 comes from the frozen D34 pair
// "task.list / task.cancel | 查看/取消任务 | L0 / L1": reading an in-process
// table touches no C3 capability face (Capabilities/Needs empty, same shape
// as task.output's) and has no path parameter (the rows carry ids and
// canonical paths the host already judged - a list is not a path argument).
func taskListDecl() Decl {
	return Decl{
		Capabilities: nil,
		Needs:        nil,
		Declared:     risk.L0,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// Execute implements Tool. One answer per wiring state, and the unwired one
// fails closed loudly: an empty list that read as "nothing runs" would be
// the exact lie ticket 164 定案② bans for task.output.
//
// Answer shapes:
//
//	unwired   -> error: "调度器未接线"（fail-closed，绝不空成功）
//	empty     -> success text saying the queue is EMPTY (a real answer)
//	rows      -> one line per row: id, label, D43 state, note, held paths
func (t taskList) Execute(ctx context.Context, _ json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if t.d.Scheduler.View == nil {
		return Result{Text: "调度器未接线（fail-closed：读不到任务行，不代表没有任务在跑）", IsError: true}, nil
	}
	rows := t.d.Scheduler.View.Rows()
	if len(rows) == 0 {
		return Result{Text: "调度器队列为空：当前没有任何任务。"}, nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "调度器任务 %d 行（排队顺序）：\n", len(rows))
	for i, r := range rows {
		fmt.Fprintf(&b, "%d. %s「%s」 状态=%s", i+1, r.ID, r.Label, string(r.State))
		if r.Note != "" {
			fmt.Fprintf(&b, "；%s", r.Note)
		}
		if len(r.HeldPaths) > 0 {
			fmt.Fprintf(&b, "；占用路径：%s", strings.Join(r.HeldPaths, "、"))
		}
		b.WriteString("\n")
	}
	return Result{Text: b.String()}, nil
}

// schedulerLookup finds one row by id through the wired view.
func schedulerLookup(v TaskListView, taskID string) (TaskListRow, bool) {
	for _, r := range v.Rows() {
		if r.ID == taskID {
			return r, true
		}
	}
	return TaskListRow{}, false
}
