package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/streamkey"
)

// The subagent layer (ticket 197 leg A: 实体层).
//
// What this file is: the ONE place in this repository that derives a subagent -
// a record with its own task id in the TaskRoster, its own agent.Loop built from
// the host's already-assembled agent.Options, and its conclusion handed back to
// the parent task with C25 provenance on it.
//
// What it deliberately is not:
//   - not a second runtime. The child loop comes from SubagentDeps.BaseOptions,
//     which the composition root answers with the SAME agent.Options the root
//     loop runs on (same C5 provider, same bridge, same approval layer).
//   - not a new state vocabulary. The status dimension is statemachine.State and
//     only ever carries one of D43's 20 names (internal/statemachine/states.go:11-31),
//     checked by statemachine.Valid; memory.TaskLog.State's four legacy words
//     (done/cancelled/running/succeeded) are ticket 196's and never enter here.
//   - not a new taint source. The child's conclusion is stamped with the
//     existing C25 name risk.SrcTaskOutput ("task.output",
//     internal/risk/provenance.go:92-94). No "subagent.output" name is invented:
//     an off-register source would be exactly the unowned channel D30 warns about.
//   - not a second stop path. The only cancel that exists is the child's own
//     context.CancelFunc, stored in the roster (TaskRoster.AttachCancel) and read
//     back by TaskRoster.Cancel.

// Subagent identity, pool and depth constants. They live in this one file on
// purpose: leg B (internal/panel) reads the KEY SHAPE, not these numbers.
const (
	// SubagentStreamKeyPrefix is ticket 197 §0's frozen stream key shape, and as
	// of the 载体层 leg it is an ALIAS of internal/streamkey's one literal, not a
	// second copy of it. Before that, this package and internal/panel each wrote
	// the string "subagent:" themselves, and only the composition root - the one
	// place that imports both - could notice a drift (the pair of pins in
	// cmd/wisp/subagent_stream_key_197_test.go, kept in force). Now a re-fork of
	// the literal is what goes red, mechanically:
	// TestSubagentStreamKeyHasOneMintSite scans the tree and allows exactly one
	// non-test file to contain it.
	SubagentStreamKeyPrefix = streamkey.SubagentPrefix

	// TaskKindRoot and TaskKindSubagent are the roster's two kinds. A root row
	// has no parent; a subagent row always names the task that derived it.
	TaskKindRoot     = "root"
	TaskKindSubagent = "subagent"

	// MaxConcurrentSubagents is the hard cap on subagents in flight. The root
	// and every one of its descendants share this ONE pool (it lives on the
	// roster, which is process-local), and exceeding it is a hard refusal with a
	// readable reason - never a silent queue into somewhere nobody can see.
	//
	// The number is the bridge's D38d ceiling (MaxToolConcurrency in
	// internal/tools/bridge.go, which is a frozen contract and not tunable
	// upward), NOT a number picked from the roster's side. Ticket 222 changed
	// WHY the two numbers mean the same thing, not what they are: a spawn that is
	// waiting for its child no longer holds a bridge slot (see
	// giveBackWhileWaiting), so the ceiling is no longer a limit on how many
	// children may EXIST at once - it is the limit on how many of them may be
	// executing a tool in the same instant. A pool LARGER than that ceiling would
	// therefore still admit rows the machine cannot run at once: the extra
	// children now take turns (轮转) through the four execution slots instead of
	// starving behind their own parents' holds, and the roster would still print
	// all of them as 「在跑」, which is the lie ticket 211 was filed for. Whether
	// the pool should be raised to 8 on top of a ceiling of 4 is NOT this file's
	// call (ticket 222 AC#5 hands that reading to the owner); until it is made,
	// the constant stays equal to the ceiling and
	// internal/tools/subagent_197_test.go's
	// Test197SubagentPoolNeverExceedsBridgeCeiling is the resident nail: raise
	// this constant above the ceiling and that leg goes red. It is written as the
	// number 4, not as an alias of MaxToolConcurrency, so that leg compares two
	// independent readings instead of one constant against itself.
	MaxConcurrentSubagents = 4

	// MaxSubagentDepth is how deep the tree may go: 1, i.e. a subagent may not
	// derive a subagent. Enforced structurally (the child's tool directory has no
	// task.spawn in it) and again by the roster check in Execute, because
	// "reachable but refused" is a weaker guarantee than "not offered".
	MaxSubagentDepth = 1
)

// subagent states, all of them D43 names (statemachine.Valid is the judge).
// The mapping is ticket 197's own reading of D43 for a task record, and it is
// named here rather than spread across call sites:
//
//	Thinking  the child loop is running
//	Settling  the child ended with an answer
//	Muted     the host stopped it (D43 has no "cancelled"; Muted is D43's name
//	          for "silenced on purpose", which is what a host cancel is)
//	Error     the child ended without an answer
//
// Reconciling this dimension with memory's four legacy words stays ticket 196's
// (ruling A394): this ticket neither reads that column nor renames anything.
const (
	subagentStateRunning = statemachine.StateThinking
	subagentStateSettled = statemachine.StateSettling
	subagentStateStopped = statemachine.StateMuted
	subagentStateFailed  = statemachine.StateError
)

// subagentInFlightState is the one D43 name that means "still holding a slot".
const subagentInFlightState = subagentStateRunning

// SubagentStreamKey is the ONE way this package names a subagent's stream, and
// it is streamkey's function, not a second spelling: the writing side and the
// reading side now cannot disagree by construction. A blank task id yields "",
// which this package used to answer with the bare prefix - a row nobody could
// attribute to. Divergence named and closed by 197-r1b §⑤2.
func SubagentStreamKey(taskID string) string {
	return streamkey.Subagent(taskID)
}

// SubagentDeps is what task.spawn needs from the host. Every field is a seam the
// composition root wires; nil on any of them is a fail-closed refusal, never a
// degraded success.
type SubagentDeps struct {
	// Roster is the shared process-local table: the pool counter, the identity
	// rows and the cancel handles all live there, so host and tool see one table.
	Roster *TaskRoster
	// BaseOptions answers the agent.Options the ROOT loop runs on. The child gets
	// a copy with three fields changed (Tools / Registry / Sink) and nothing else,
	// which is what "同一套宿主装配" means concretely. ok=false refuses the spawn.
	BaseOptions func() (agent.Options, bool)
	// ParentTools is the host's C4 surface (the bridge). The child receives it
	// wrapped so task.spawn is not merely refused but absent from its directory.
	ParentTools agent.ToolProvider
	// Provenance is the C25 engine used to stamp the child's conclusion into the
	// PARENT task's taint scope as risk.SrcTaskOutput. nil means the stamp cannot
	// happen, and the reply then says so out loud instead of handing the parent an
	// untraced outside text.
	Provenance *risk.Provenance
	// Stream is the subagent's text channel, and it is the panel's stream log
	// handed over as the two methods this package needs - nothing more, so the
	// tool layer holds no pointer to the view layer and no new field enters the
	// AC#6 "is this an approval outlet?" enumeration.
	//
	// 197-r3 widened it from a bare append to Append+Close, and that is the whole
	// carrier fix: opt.Sink used to be nil outright, so a child's streamed deltas
	// were discarded and its work page held exactly the two lifecycle lines below
	// (197-r2's own §6.1 is the reading that named it), and Close was the missing
	// half of "this page is finished" - without it every subagent row streamed
	// forever on a settled task, which is the failure consoleSink's EvDone branch
	// exists to prevent for a root.
	//
	// nil keeps the pre-r3 shape (lifecycle lines only, no page) and is a
	// statement about the assembly, never a degraded claim: nothing below pretends
	// a child streamed when it was given nowhere to stream to.
	Stream SubagentStreamSink
}

// SubagentStreamSink is the whole text surface one subagent has. *panel.StreamLog
// satisfies it as it stands, which is why the composition root can hand the log
// itself over instead of wrapping it, and why this package needs neither the
// panel's types nor a second channel.
type SubagentStreamSink interface {
	Append(key, text string)
	Close(key string)
}

type subagentSpawn struct{ d SubagentDeps }

type subagentSpawnArgs struct {
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

// subagentSpawnSchema keeps description and prompt as TWO fields: a short title
// for the roster label and a self-contained task book for the child. Merging
// them would make the roster row unreadable and the child under-specified.
var subagentSpawnSchema = JSONSchema(`{"type":"object","properties":{` +
	`"description":{"type":"string","description":"短标题，只放任务名，进名册当 label，不作为子代理的输入"},` +
	`"prompt":{"type":"string","description":"交给子代理的自包含任务书：子代理看不到父任务的对话，需要的背景都得写进来"}` +
	`},"required":["description","prompt"],"additionalProperties":false}`)

func (subagentSpawn) Name() string { return "task.spawn" }

// Description carries the two facts the model cannot be allowed to guess:
// cancellation does not cascade, and a subagent cannot derive a subagent.
// The two numbers are read from the constants, never typed by hand: a model-facing
// sentence that says "8" while the pool holds 4 is the same lie ticket 211 was
// filed for, just told to the model instead of to the user.
//
// Ticket 221 (乙形, landed in the SAME commit as 甲形 per A434 item 7): the clause
// that used to read 「可以用 task.cancel 单独停它」 promised an unqualified power
// while nothing registered that tool and TaskRoster.Cancel had zero production
// callers. What stands here now is what is true on both sides of that commit:
// task.cancel exists, and it only lets the row's OWN parent stop it. ⚠ Two
// literals are load-bearing and must survive any future reword: 「停掉父任务不会级联」
// (pinned by subagent_197_test.go:797) and the two numbers (pinned at :419).
func (subagentSpawn) Description() string {
	return fmt.Sprintf("派生一枚子代理去独立完成一个子任务，等它跑完并把结论带回本任务；"+
		"它会在任务名册里留下一行有父子关系与状态的记录；"+
		"注意：停掉父任务不会级联停掉子代理；task.cancel 只有派生它的那枚父任务能用它单独停孩子——"+
		"子代理停兄弟、停自己都一律被拒；由用户停某一枚子代理还是另一条通道（票 181／票 220），今天没有落点；"+
		"子代理不能再派生子代理（深度 %d），同时在跑的子代理上限 %d 枚",
		MaxSubagentDepth, MaxConcurrentSubagents)
}

func (subagentSpawn) Parameters() JSONSchema { return subagentSpawnSchema }

// TaskSpawnDecl declares the row. L0 floor and NO capability tokens, the same
// shape task.output documents above: the spawn itself touches no C3 capability
// face - everything the child does goes through the same bridge, the same C19
// assessor and the same approval queue as the parent's, so declaring a capability
// here would be a second, redundant grant surface. C19 still owns the verdict.
func TaskSpawnDecl() Decl {
	return Decl{
		Capabilities: nil,
		Needs:        nil,
		Declared:     risk.L0,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinSubagentEntries returns the ticket 197 spawn entry. It is a separate
// function from BuiltinTaskEntries because that one's registration shape is
// pinned by ticket 164's tests as the task.output-only family; nothing here
// re-pins or loosens those legs.
func BuiltinSubagentEntries(d SubagentDeps) []Entry {
	return []Entry{{Tool: subagentSpawn{d: d}, Decl: TaskSpawnDecl()}}
}

// Execute implements Tool: derive, record, run, record back, answer the parent.
func (t subagentSpawn) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a subagentSpawnArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	label := strings.TrimSpace(a.Description)
	prompt := strings.TrimSpace(a.Prompt)
	if label == "" || prompt == "" {
		return Result{Text: "description（短标题）与 prompt（自包含任务书）都得给，两枚字段不能互相替代", IsError: true}, nil
	}
	label = capRunes(label, 120)

	if t.d.Roster == nil {
		return Result{Text: "任务名册未接线（fail-closed：拒绝派生子代理，派生了也没有地方登记它）", IsError: true}, nil
	}
	if t.d.BaseOptions == nil || t.d.ParentTools == nil {
		return Result{Text: "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，fail-closed：不起第二套运行时）", IsError: true}, nil
	}
	// The task-level id of the caller: since ticket 242 the correlation id is
	// minted per call (C18 routes approval replies by it), so it can no longer
	// name a task; the roster is keyed by task ids.
	parentID := TaskID(ctx)
	if parentID == "" {
		return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，" +
			"就没法登记父子关系，也没法把结论盖戳进父任务的 C25 作用域）", IsError: true}, nil
	}
	if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {
		return Result{Text: fmt.Sprintf(
			"拒绝派生：深度上限是 %d，而派生者 %s 自己就是一枚子代理（父 kind=%s）。"+
				"子代理不许再派子代理。", MaxSubagentDepth, parentID, rec.Kind), IsError: true}, nil
	}

	release, ok := t.d.Roster.TryAcquireSubagentSlot(MaxConcurrentSubagents)
	if !ok {
		return Result{Text: fmt.Sprintf(
			"拒绝派生：同时在跑的子代理已达上限 %d 枚（在跑的：%s）。"+
				"这是硬拒，不是排队——等哪一枚结束了再派，或者把任务并成一枚子代理。",
			MaxConcurrentSubagents, strings.Join(t.d.Roster.RunningSubagentIDs(), ", ")), IsError: true}, nil
	}

	base, wired := t.d.BaseOptions()
	if !wired {
		release()
		return Result{Text: "宿主现在给不出装配快照（fail-closed：不在没有父装配的情况下起子环路）", IsError: true}, nil
	}
	if base.AdmitTask == nil {
		// A child with no admission hook has no route to the host's approval
		// layer at all. Deriving one would be handing it a private decision
		// surface, which is the "子代理自批" shape this ticket must not ship.
		release()
		return Result{Text: "拒绝派生：宿主装配里没有审批准入钩子（AdmitTask 未接线），" +
			"子代理的请求就没有任何队列可去（fail-closed：子代理永不自批）", IsError: true}, nil
	}

	opt := base
	opt.Tools = newSubagentToolProvider(t.d.ParentTools)
	opt.Registry = observe.NewRegistry()
	// The child's text must NOT ride the parent's console sink: that is how two
	// tasks' output gets merged into one visible stream. It gets its own sink,
	// built from the one text channel the host injected, and that channel keys the
	// child by its own task id (see SubagentDeps.Stream for why r3 had to widen it
	// - before this leg the child had NO sink at all, so its streaming work was
	// thrown away and the page was two lifecycle lines).
	opt.Sink = nil
	if t.d.Stream != nil {
		opt.Sink = subagentTextSink{d: t.d}
	}
	// The child is admitted through the PARENT's admission hook, i.e. the host's
	// one approval queue - which is also the second publisher of its roster row
	// (the loop calls this before its first model call, internal/agent/loop.go:360).
	baseAdmit := base.AdmitTask
	opt.AdmitTask = func(taskID string) func() {
		t.d.Roster.PublishSubagent(taskID, parentID, label)
		revoke := baseAdmit(taskID)
		return func() {
			if revoke != nil {
				revoke()
			}
		}
	}

	child, err := agent.New(opt)
	if err != nil {
		release()
		return Result{Text: "子代理环路装配失败：" + err.Error(), IsError: true}, nil
	}

	// Parent cancellation must not cascade (ticket 197 §0). WithoutCancel keeps
	// the call's values (including its cancel handle and C22 deadline carrier)
	// while dropping its cancellation/deadline, so the child is a task the host
	// can still find, stop and audit - not a detached goroutine.
	// The child loop's task id is minted by the loop itself (agent.newTaskID,
	// internal/agent/loop.go:321-327, which mints the id BEFORE reg.Spawn), and
	// Loop has no exported entry point that takes a caller-supplied id - that
	// gap is the open question Q-56 recorded in bridge.go's CloseTask comment,
	// and inventing one is not this ticket's call. So the row is published from
	// THIS goroutine as the first act after the id exists, and the child's own
	// admission hook (below, always fired before its first model call) publishes
	// the identical row a second time: there is no window in which a running
	// child has no roster row, which is the DSH/Step-Code failure mode this step
	// exists to prevent.
	// Parent cancellation must not cascade (ticket 197 §0). WithoutCancel keeps
	// the call's values (its cancel handle carrier, the C25 boxes) while dropping
	// its cancellation and its deadline, so the child is a task the host can still
	// find, stop and audit - not a detached goroutine.
	childCtx, cancelChild := context.WithCancel(context.WithoutCancel(ctx))
	bg := child.RunAsync(childCtx, prompt)
	t.d.Roster.PublishSubagent(bg.ID, parentID, label)
	t.d.Roster.AttachCancel(bg.ID, cancelChild)
	// Ticket 222 (丙): the wait below is not an execution, so it must not hold a
	// D38d bridge token. Until here the spawn held one across the whole wait for
	// its child, and the child asks for one of those SAME tokens to run any tool
	// (cmd/wisp hands the bridge over as SubagentDeps.ParentTools) - with the pool
	// full, every parent collected its own per-tool deadline and answered
	// "不等了"/"工具 task.spawn 超时" while a healthy child sat queued behind the slot
	// its own parent was holding. Giving the token back costs the frozen ceiling
	// nothing: at most four tool calls still RUN at once, because a parked parent
	// runs nothing. The token is not taken back afterwards (see inFlightSlot for
	// why reacquiring would recreate the block), and everything this call still
	// does after the wait - finalize, the C25 stamp, the journal row - is
	// bookkeeping, not capability execution.
	//
	// The call sits HERE, ahead of the two "已派生" lines below, on purpose: a
	// reader that has seen this call's spawn delta therefore knows the slot is
	// already back, which is what makes ticket 222's occupancy leg decide without
	// a deadline. The pool slot is a different thing and is NOT handed back here -
	// the child still owns it until it really joins.
	giveBackWhileWaiting(ctx)
	t.feed(bg.ID, "已派生："+label)
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("task.spawn 已派生子代理 %s（%s）", bg.ID, label))
	}

	// The finish watcher runs on the registry (D38b: no bare goroutine), and it
	// owns the slot and the row: the child's roster record is finalised when it
	// has really joined, INCLUDING on the path where the parent stopped listening
	// - otherwise a detached child would sit in the roster reading "Thinking"
	// forever, which is the lie this layer exists to prevent.
	done := make(chan agent.Result, 1)
	var releaseOnce sync.Once
	finish := func() { releaseOnce.Do(release) }
	var stampWhy string
	watchRoot := observe.NewRootFrom(context.WithoutCancel(ctx), "subagent-finish-"+bg.ID)
	observe.Default.Spawn("subagent-finish-"+bg.ID, "tools", watchRoot, func(c context.Context) {
		defer finish()
		res := bg.Wait()
		// The slot goes back BEFORE the row is finalised, so a reader that sees
		// the settled row can never see the pool still holding its place.
		finish()
		stampWhy = t.finalize(parentID, bg.ID, label, res)
		done <- res
	})

	select {
	case <-ctx.Done():
		// The parent's call went away. The child keeps running and stays in the
		// roster; nothing here cancels it.
		return Result{
			Text: fmt.Sprintf(
				"父任务这一侧已经不等了（%v）。子代理 %s 没有被级联取消，它仍在名册里，"+
					"它的父任务 %s 还能用 task.cancel 单独停它（只有派生它的那一枚能停）；"+
					"它自己的流键是 %s。",
				ctx.Err(), bg.ID, parentID, SubagentStreamKey(bg.ID)),
			IsError: true,
		}, nil
	case res := <-done:
		return t.answer(bg.ID, label, res, stampWhy), nil
	}
}

// finalize files the finished conclusion back onto the SAME row (state checked by
// statemachine.Valid through TaskOutput.StateAnswer's own judge) and stamps it
// into the parent's C25 scope. It returns the reason the stamp did not land, or
// "" when it did, so the reply can say so out loud.
func (t subagentSpawn) finalize(parentID, childID, label string, res agent.Result) string {
	state := subagentStateSettled
	switch res.Status {
	case agent.StatusCancelled:
		state = subagentStateStopped
	case agent.StatusFailed:
		state = subagentStateFailed
	}
	text := res.Text
	if text == "" && res.Message != "" {
		text = res.Message
	}
	t.d.Roster.Record(childID, TaskOutput{
		Text: text, State: state,
		ParentTaskID: parentID, Label: label, Kind: TaskKindSubagent,
	})
	t.d.Roster.DetachCancel(childID)
	t.feed(childID, fmt.Sprintf("已结束（%s）", string(state)))
	// The row the page renders is closed here, not left open: "已结束" as text
	// without done=true would still read as a stream in flight.
	t.feedDone(childID)
	return t.stampConclusion(parentID, childID, text)
}

// answer writes the text the parent's model reads.
func (t subagentSpawn) answer(childID, label string, res agent.Result, stampWhy string) Result {
	text := res.Text
	if text == "" && res.Message != "" {
		text = res.Message
	}
	state := subagentStateSettled
	switch res.Status {
	case agent.StatusCancelled:
		state = subagentStateStopped
	case agent.StatusFailed:
		state = subagentStateFailed
	}
	body := fmt.Sprintf("子代理 %s（%s）已结束，状态 %s。\n结论：\n%s",
		childID, label, string(state), text)
	if strings.TrimSpace(text) == "" {
		body = fmt.Sprintf("子代理 %s（%s）已结束，状态 %s，但它没有留下正文（这是「没有正文」，不是「没跑」）。",
			childID, label, string(state))
	}
	if stampWhy != "" {
		body += "\n[注意：" + stampWhy + "]"
	}
	return Result{Text: body}
}

// stampConclusion is the taint leg: the child's conclusion enters the PARENT's
// context, so it is marked in the parent's scope as the existing C25 source
// risk.SrcTaskOutput. It returns "" when the stamp landed, and the reason it did
// not - which the reply then says out loud rather than quietly dropping.
func (t subagentSpawn) stampConclusion(parentID, childID, text string) string {
	switch {
	case strings.TrimSpace(text) == "":
		return ""
	case t.d.Provenance == nil:
		return "C25 溯源引擎没接线，这段子代理结论没有盖 task.output 源名戳"
	case parentID == "":
		return "父任务 id 未知，子代理结论没有盖 task.output 源名戳"
	}
	if !t.d.Provenance.Mark(parentID, risk.SrcTaskOutput, SubagentStreamKey(childID), text) {
		return "子代理结论没能进父任务的 task.output 索引（文本太短，C25 取不到可匹配的片段）"
	}
	return ""
}

func (t subagentSpawn) feed(childID, text string) {
	if t.d.Stream == nil {
		return
	}
	key := SubagentStreamKey(childID)
	if key == "" {
		// A child with no id has no stream to write to, and appending under ""
		// would open a row the panel cannot attribute to anybody.
		return
	}
	t.d.Stream.Append(key, text)
}

// feedDone closes one child's own stream, i.e. puts done=true on the row the page
// renders. Without it a finished subagent keeps streaming on screen forever,
// which is the failure internal/panel's Close exists for (pump.go:443-461) and
// what consoleSink's EvDone branch says out loud for a root task.
func (t subagentSpawn) feedDone(childID string) {
	if t.d.Stream == nil {
		return
	}
	key := SubagentStreamKey(childID)
	if key == "" {
		return
	}
	t.d.Stream.Close(key)
}

// subagentTextSink is the child loop's own sink (ticket 197 载体层). Two rules
// make it a sibling of cmd/wisp's consoleSink rather than a copy of it: nothing
// here writes to stdout, because two tasks' text on one console is a merge; and
// every line goes under the CHILD's own stream key, so the page opened by
// clicking that row is that agent's work and nobody else's.
//
// Like the console sink it must never block - the loop publishes synchronously on
// the task's goroutine - and reasoning deltas stay out of the results channel for
// the reason consoleSink gives: panel.ResultChunk's contract is an assistant
// RESULT (internal/panel/composer.go's ResultChunk header).
type subagentTextSink struct{ d SubagentDeps }

func (s subagentTextSink) Publish(e agent.Event) {
	if s.d.Stream == nil {
		return
	}
	switch e.Kind {
	case agent.EvTextDelta:
		s.d.Stream.Append(SubagentStreamKey(e.TaskID), e.Text)
	case agent.EvDone:
		if key := SubagentStreamKey(e.TaskID); key != "" {
			s.d.Stream.Close(key)
		}
	}
}

func capRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	return string([]rune(s)[:max])
}

// subagentToolProvider is the structural half of MaxSubagentDepth: the child's
// C4 directory is the parent's, minus task.spawn, and a call for the name it can
// not see is refused with the reason rather than dispatched.
type subagentToolProvider struct{ inner agent.ToolProvider }

func newSubagentToolProvider(inner agent.ToolProvider) agent.ToolProvider {
	return &subagentToolProvider{inner: inner}
}

const subagentHiddenTool = "task.spawn"

func (p *subagentToolProvider) Tools(ctx context.Context) ([]agent.ToolInfo, error) {
	if p.inner == nil {
		return nil, nil
	}
	list, err := p.inner.Tools(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]agent.ToolInfo, 0, len(list))
	for _, info := range list {
		if info.Name == subagentHiddenTool {
			continue
		}
		out = append(out, info)
	}
	return out, nil
}

func (p *subagentToolProvider) Execute(ctx context.Context, req agent.ToolRequest) (agent.ToolOutcome, error) {
	if req.Name == subagentHiddenTool {
		return agent.ToolOutcome{
			Text: fmt.Sprintf("拒绝执行：子代理不能再派生子代理（深度上限 %d），"+
				"task.spawn 也不在你的工具目录里。", MaxSubagentDepth),
			IsError: true,
		}, nil
	}
	return p.inner.Execute(ctx, req)
}
