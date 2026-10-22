package tools

// Ticket 164 AC#4: 取消要真取消 —— 取消之后不许有"还在写"的尾巴。
//
// The object this case builds is the one ticket 164's ruling 5 said did not
// exist on 2026-09-27 ("RunAsync 生产零调用点 ⇒ 连被测对象都没有"). It exists
// now: cmd/wisp/run.go:1176 starts a task through Loop.RunAsync and :1184-1187
// files its answer through TaskBackfill once Wait() returns, while the subagent
// door has filed terminal rows through subagentSpawn.finalize since ticket 197.
// So a "tail" here is not a hypothetical: it is the host's own write side,
// handed a finished task's context, trying to land output after the stop.
//
// "No tail" has to be stated precisely, because ONE write after a cancel is
// legitimate: subagentSpawn.finalize (subagent_197.go:450-472) files the
// terminal state that task.cancel's reply quotes, and ticket 176-r1's L-4 leg
// (TestTerminalWriteCarrierStaysOutOfScope176r1) names that carrier as out of
// the criterion's range. The tail this case forbids is the NEXT write, the one
// that moves a row which has already settled.
//
// The two seams this case composes are the two production ones:
//
//	TaskRoster.Cancel      the ONE stop path (task.go:447), which task.cancel
//	                       calls; its whole job is to invoke the handle the
//	                       child was spawned with, so "cancel" is not a word
//	                       the roster says, it is a context that ends.
//	TaskBackfill.Backfill  the output-write port (task_backfill.go:105); it
//	                       reads the finished task's own ctx once, via Stopped
//	                       (cancel.go:86), before it touches the table.
//
// Neither is mocked: the terminal filing runs through finalize itself, and the
// tail runs through Backfill itself. What this case adds is the ORDER - cancel,
// settle, then attempt the tail - which is AC#4's sentence and which neither
// ticket's own legs state.
//
// Teeth (measured 2026-10-22, both directions, production restored byte for
// byte afterwards; the readings live in the delivery evidence):
//
//	M1  TaskRoster.Cancel files the stop without calling the handle
//	    (task.go:462) -> the child's ctx survives, the tail lands, and the row
//	    moves: red on the tail line, on the row line, and on the ctx invariant.
//	M2  Backfill's Stopped check is neutered (task_backfill.go:123) -> the
//	    stop is real but the write door no longer reads it -> red on the tail
//	    line and on the "row moved" line.
//
// As-is the case is green (the guard landed with ticket 176 AC#3). It is a
// guard, not a historical reading: the ticket's AC#4 note says which readings
// predate it and which one this one is.

import (
	"context"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

const (
	ac4Parent = "ac4-parent-0000000000000000"
	ac4Child  = "ac4-child-00000000000000000"
	ac4Live   = "ac4-live-000000000000000000"
)

// ac4CancelledChild builds one child row, stops it through the roster's single
// stop path, and settles the row the way production does. It returns the roster
// and the child's own context - the one the host's write side is handed - so the
// caller can attempt the tail on exactly the context the cancel killed.
func ac4CancelledChild(t *testing.T) (*TaskRoster, context.Context) {
	t.Helper()
	roster := NewTaskRoster()
	roster.MarkRoot(ac4Parent, "根任务")
	roster.PublishSubagent(ac4Child, ac4Parent, "尾巴")

	childCtx, cancelChild := context.WithCancel(context.Background())
	roster.AttachCancel(ac4Child, cancelChild)

	if stopped, why := roster.Cancel(ac4Child); !stopped {
		t.Fatalf("名册没停掉自己的孩子：%s", why)
	}

	// The legitimate post-cancel write, through the real finalize rather than a
	// stand-in: this is the row the tail would land on, produced by the same
	// function production uses (subagent_197.go:462-466, then DetachCancel).
	spawn := subagentSpawn{d: SubagentDeps{Roster: roster}}
	spawn.finalize(ac4Parent, ac4Child, "尾巴", agent.Result{
		Status: agent.StatusCancelled, Text: "被停掉前的半句正文",
	})
	return roster, childCtx
}

func Test164AC4CancelLeavesNoWritingTail(t *testing.T) {
	roster, childCtx := ac4CancelledChild(t)

	settled, ok := roster.Look(ac4Child)
	if !ok {
		t.Fatalf("终态那一行不见了：取消之后名册不许丢行")
	}
	if settled.State != statemachine.StateMuted {
		t.Fatalf("终态不是取消那一支该落的 D43 名：state=%q, want %q", settled.State, statemachine.StateMuted)
	}

	// The tail attempt: output that arrived after the stop, handed to the write
	// port with the SAME context the task was stopped through - exactly what
	// cmd/wisp/run.go:1187 hands Backfill after bg.Wait().
	_, why := (TaskBackfill{Roster: roster}).Backfill(childCtx, ac4Child, "取消之后尾巴还在写的正文")
	if why == "" {
		t.Errorf("取消之后那股输出还是被写进名册了：Backfill 没有拒绝——取消之后不许有还在写的尾巴")
	} else if !strings.Contains(why, "取消") {
		t.Errorf("拒绝了，但理由没有说是取消：%q（模型/宿主读到的必须是这一支）", why)
	}

	after, _ := roster.Look(ac4Child)
	if after.Text != settled.Text || after.State != settled.State {
		t.Errorf("取消后名册那一行被改动了：Text %q -> %q，State %q -> %q"+
			"（取消之后落到这一行的写就是那条尾巴）",
			settled.Text, after.Text, settled.State, after.State)
	}

	// Why the two lines above can be trusted to mean "the cancel was real": a
	// successful Cancel must have RUN the handle, so the writer's own context is
	// gone. Read last on purpose - when a stop announces itself without
	// delivering, the tail lands first and this line names the cause instead of
	// pre-empting the observation.
	if err := childCtx.Err(); err == nil {
		t.Errorf("Cancel 返回成功，但它挂上的句柄没有真停掉子任务：ctx.Err() == nil（取消只是把话说出去了）")
	}

	// Control half, same port and same row shape: a context nobody cancelled
	// must still write, or the refusal above is green on silence.
	live := NewTaskRoster()
	live.PublishSubagent(ac4Live, ac4Parent, "没被停的")
	if _, liveWhy := (TaskBackfill{Roster: live}).Backfill(
		context.Background(), ac4Live, "没被取消的正文"); liveWhy != "" {
		t.Fatalf("对照组没能写进名册（why=%q）：上面那条拒绝就是空响", liveWhy)
	}
	if rec, ok := live.Look(ac4Live); !ok || rec.Text != "没被取消的正文" {
		t.Fatalf("对照组那一行读不回刚写的正文：present=%v rec=%+v", ok, rec)
	}
}
