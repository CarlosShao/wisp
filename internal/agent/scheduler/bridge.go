package scheduler

import (
	"context"
	"strings"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The composition adapter (ticket 47 AC#6): the seam that lets a REAL agent
// loop run INSIDE the scheduler, so task.list / task.cancel can reach the
// scheduler's state through the agent's own C4 tool surface. tools may not
// import agent/scheduler (internal/tools is consume-only for the scheduler;
// the dependency arrow would also make every fs tool pay for it), so this
// file hands tools the two narrow adapters its SchedulerDeps bag asks for.
//
// Disposition note: the scheduler code (this package) lands with ticket 47
// with its own tests over the REAL loop through this adapter; the
// composition wiring into cmd/wisp's run.go assembly is the assembly
// ticket's write face (the scheduler is created there, the adapter is
// wrapped around the bridge, the entries register through the existing
// TaskDeps seams).
//
// No allow-door travels here: the adapters implement exactly the two
// interfaces tools.TaskDeps declares (tools.TaskListView: Rows only;
// tools.TaskCanceler: Cancel only) - neither carries an allow/answer verb,
// which is what the 197 allow-door walker demands of anything TaskDeps
// holds (it walks every field's type graph).
var (
	_ tools.TaskListView = schedulerView{}
	_ tools.TaskCanceler = schedulerControl{}
)

// NewBridgeDeps assembles the deps bag the task family extends with. A nil
// scheduler yields the zero SchedulerDeps (the tools' fail-closed shape),
// never a panicking adapter.
func NewBridgeDeps(s *Scheduler) tools.SchedulerDeps {
	if s == nil {
		return tools.SchedulerDeps{}
	}
	return tools.SchedulerDeps{
		View:    schedulerView{s},
		Control: schedulerControl{s},
	}
}

// schedulerView adapts *Scheduler onto tools.TaskListView (no new verbs).
type schedulerView struct{ s *Scheduler }

// Rows maps the scheduler's rows onto the tools' TaskListRow shape. The
// authority join travels here: ParentTaskID is the task id whose task.cancel
// may stop this row (the host files it through SetParent when it submits a
// task on behalf of a parent task). A row without one is a root - and the
// tools' stop branch refuses roots, exactly like the roster branch does.
func (v schedulerView) Rows() []tools.TaskListRow {
	rows := v.s.Rows()
	out := make([]tools.TaskListRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, tools.TaskListRow{
			ID:           r.ID,
			Label:        r.Label,
			State:        r.State,
			Note:         r.Note,
			ParentTaskID: v.s.ParentOf(r.ID),
			HeldPaths:    r.HeldPaths,
		})
	}
	return out
}

// schedulerControl adapts *Scheduler onto tools.TaskCanceler.
type schedulerControl struct{ s *Scheduler }

func (c schedulerControl) Cancel(id string) error { return c.s.Cancel(id) }

// ExecAdapter wraps an agent.Loop runner into a scheduler TaskSpec Exec: it
// runs one loop task under the scheduler's root ctx, maps the loop's Status
// onto the TaskResult vocabulary, and carries the applied-steps ledger of
// every tool call the task logged (D31: a cancelled task's evidence, never
// rendered as a clean undo).
type ExecAdapter struct {
	// Run runs one loop task synchronously (agent.Loop.Run or an equivalent
	// with the same Result contract).
	Run func(ctx context.Context, input string) agent.Result
	// ParentTaskID is the task id the row files as its parent: the CALLER
	// whose task.cancel may stop this row. Empty = root row (nobody stops
	// it through the tool).
	ParentTaskID string
}

// TaskSpec builds the scheduler TaskSpec for one input. The parent filing
// happens at Exec start (the row exists from Submit, the parent link is
// host-owned data, so the Submit caller can also SetParent directly).
func (a ExecAdapter) TaskSpec(input string) TaskSpec {
	return TaskSpec{
		Label: labelOf(input),
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			return mapResult(a.Run(ctx, input))
		},
	}
}

// mapResult translates the loop's Result onto the scheduler's TaskResult,
// carrying the applied-steps ledger for a cancelled task (D31: the evidence
// travels with the row, never rendered as a clean undo).
func mapResult(res agent.Result) TaskResult {
	out := TaskResult{Status: statusOf(res.Status), Text: res.Text}
	switch out.Status {
	case "failed":
		out.Err = res.Message
	case "cancelled":
		// Every tool call that REPORTS success already landed. The loop
		// books a cancelled in-flight call as outcome "cancelled" with the
		// 「任务已取消，调用被中止」 wording, so the two never blur.
		var steps []string
		for _, tl := range res.ToolLog {
			if tl.Outcome == "success" {
				steps = append(steps, toolStep(tl))
			}
		}
		out.AppliedSteps = steps
	}
	return out
}

// statusOf maps the loop's Status vocabulary onto the TaskResult's. Control
// never reaches the scheduler (a control word is answered before a task is
// spawned); everything unmapped is a failure, not a success.
func statusOf(s agent.Status) string {
	switch s {
	case agent.StatusCompleted:
		return "completed"
	case agent.StatusCancelled:
		return "cancelled"
	default:
		return "failed"
	}
}

// toolStep renders one applied tool call for the ledger.
func toolStep(tl agent.ToolResultLog) string {
	return tl.Name + "：" + firstLine(tl.Text)
}

// firstLine is the ledger's one-line summary of what a call did.
func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 120 {
		s = s[:120] + "…"
	}
	return s
}

// labelOf derives the row label from the input text (the task's NAME only,
// never its body - the same rule the roster's Label field documents).
func labelOf(input string) string {
	s := strings.TrimSpace(input)
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i]
	}
	if len(s) > 40 {
		s = s[:40] + "…"
	}
	if s == "" {
		s = "（无标题任务）"
	}
	return s
}
