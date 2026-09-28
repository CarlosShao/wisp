package tools

import (
	"context"

	"github.com/CarlosShao/wisp/internal/agent"
)

// This file is the WRITE side of ticket 164's process-local task table, the
// half ticket 176 says was never built: internal/agent has an async entry point
// (Loop.RunAsync, loop.go:321) with zero production call sites, so
// TaskRoster.Record had zero production writers and task.output could only ever
// answer "查不到这个任务". The composition root calls RunAsync and, once the
// task has fully joined, hands the finished answer to Backfill below.
//
// What this file deliberately does NOT do:
//
//   - it does not change Record. TaskRoster.Record stays last-writer-wins
//     (task.go's own comment states that verbatim, and D15(3) documents the same
//     semantics for the artifacts it points at). "Do not write a cancelled
//     task's output" is achieved the dumb way the ruling named: the caller only
//     reaches this after the task has joined, and this function looks at the
//     cancellation state once before it writes.
//   - it adds no method to agent.Loop (ticket 164 定案④ + the G3 leg of
//     probes/154/gate-clauses.sh). Everything here runs on already-exported
//     surface: Loop.RunAsync, RunningTask.Wait / Root / Pending, Loop.Budgets
//     and agent.NewSpiller.

// taskArtifactPrefix is the logical prefix a background task's artifact key
// carries, and it is not a new invention: it is the exact name the loop gives
// its own spawned task in the D38 roster (loop.go:324,
// reg.Spawn("agent-task-"+id, "agent", ...)).
//
// The reason the artifact key cannot be the bare task id: the spill file name
// is derived from whatever key it is handed (agent.spill.go artifactName), and
// the loop's ordinary spill path hands it the tool-CALL id, which is
// model-supplied. agent.newTaskID() is a plain lowercase hex uuid, so a bare
// task id is a perfectly reachable call id - and two keys that encode to the
// same name land on the same file, where spill.go:117-131 documents
// last-writer-wins. A collision there is not an error, it is a silent byte
// swap: the model re-reads a pointer the host printed and gets somebody else's
// bytes, with C25 provenance broken and nothing raised. Prefixing the key (not
// changing the overwrite semantics - that is contract-adjacent) is what keeps
// the two namespaces apart.
const taskArtifactPrefix = "agent-task-"

// TaskBackfill is the host's write-side dependency set for one finished
// background task: the roster the record lands in, and the real D15(3) spill
// layer that produces the artifact the pointer names.
//
// Spills is the SAME agent.Spiller the loop uses for tool results, built by the
// composition root from the loop's own scaled budgets and the store's artifacts
// directory. It is a *agent.Spiller and not an interface on purpose: a mock
// here would let this file "pass" while nothing was ever written to disk, which
// is exactly the shape AGENTS §1.3 forbids.
type TaskBackfill struct {
	Roster *TaskRoster
	Spills *agent.Spiller
}

// Backfill files one background task's answer so task.output can find it.
//
// ctx is the finished task's OWN cancellation context (the running task's root
// ctx), never the loop's internal terminal-write carrier, and never a ctx the
// model chose. The rule it enforces is ticket 176 AC#3: once a task has been
// cancelled, nobody may still be writing its output - so this looks at
// Stopped(ctx) once, before it touches the table, and reports why it stayed
// quiet. Stopped is the existing tools helper (cancel.go:74): the task's ctx is
// gone, or a veto reached this call. A ctx that deliberately survives its
// parent's cancellation - context.WithoutCancel, which is how the loop's own
// terminal journal write at loop.go:956-958 is built, and the one legitimate
// "still writing after cancel" in this repository - reads as NOT stopped here,
// which is what keeps this criterion from reading the host's own teardown as a
// violation.
//
// The second return value is "" when the record is in the table, and otherwise
// the reason it is not. It is a string and not an error because every branch
// here is a decision the host reports, not a failure the run returns to: an
// uncalled roster, an empty answer, a cancelled task, a spill that could not
// land. The caller says any of those out loud (cmd/wisp's audit line) rather
// than letting the model read "查不到这个任务" as "the task printed nothing".
//
// ArtifactPath follows ticket 174's rule of telling the truth: it is only ever
// the path of bytes that really landed. Short answers have no artifact at all
// and get "", and task.output then answers with the whole text; a spill that
// failed also gets "", and task.output announces 「截断且无副本文件」 instead of
// pointing at a file that is not there.
func (b TaskBackfill) Backfill(ctx context.Context, taskID, text string) (TaskOutput, string) {
	switch {
	case b.Roster == nil:
		return TaskOutput{}, "任务名册未接线（fail-closed：不替没人建的名册编一条记录）"
	case taskID == "":
		return TaskOutput{}, "这个任务没有宿主铸的 id（id 只由环路生成，组合根不许自造一枚）"
	case text == "":
		return TaskOutput{}, "这个任务没有留下正文（空正文不进名册，免得「查不到」和「没打印」被读成同一件事）"
	}
	if stopped, why := Stopped(ctx); stopped {
		return TaskOutput{}, "任务已取消，不再写它的输出：" + why
	}

	rec := TaskOutput{Text: text}
	var spillWhy string
	if b.Spills != nil {
		sp, err := b.Spills.Prepare(taskID, text)
		switch {
		case err != nil:
			// No path, and the reason travels with it: the record still holds
			// the text the task produced, and task.output will say the rest is
			// unrecoverable rather than inventing a file.
			spillWhy = "超长输出没能落盘，名册里不许诺任何路径：" + err.Error()
		case sp.Path != "":
			rec.ArtifactPath = sp.Path
		}
	}
	b.Roster.Record(taskID, rec)
	return rec, spillWhy
}
