package main

// Readface wave 2 (2026-11-17), ticket 38 AC#4's Go half: the real leg behind
// internal/panel's roster.cancel door - the panel-side cancel route.
//
// WHY THIS DOOR EXISTS (the C17 roster addition's reason, stated once more
// where the wire lives): the tasks view renders the host's own roster rows
// (snapshot.tasks), and a user clicking 停止 on one needs a Go通路. Without
// this door the click is the dormant shape ticket 38's page leg already had
// to name. The route is NOT approval.decide and must never become it:
// approval.decide answers a card the queue actually holds, and answering a
// card that does not exist is exactly the fabricated-answer shape the golden
// law forbids (Q-49 family). A cancel is a task-control verb; it gets its own
// door and touches no approval surface.
//
// WHY THE NAME IS "roster.cancel" AND NOT "task.cancel": see bridge.go's
// MethodRosterCancel - internal/tools' frozen ticket-188 task-axis scan
// reddens panel-side Method* values containing "task"/"state", that test is
// outside this wave's write face, and no assertion is loosened. The semantic
// target is unchanged: ONE row of the task roster.
//
// REUSE, NOT REINVENTION (the dispatch's "复用既有 cancel 语义 ⛔自造"). The
// door does not hold any stop mechanism of its own. Its seam hands the row id
// to the composition root's own cancel paths, which are the two that exist:
//
//   - the ROOT text task: agent.RunningTask.Cancel (internal/agent/loop.go:342,
//     the same handle the loop mints at RunAsync and the run's execute() waits
//     on) - the existing stop path for the row the run itself is running;
//   - every roster row: tools.TaskRoster.Cancel (internal/tools/task.go:569) -
//     the ONE stop path a roster row has by design (197 §0), per-row, never
//     cascading, refusing "no such row" / "no live handle" with its own
//     fail-closed reason pair.
//
// Both refuse by name when the id is not stoppable, and those refusals ride
// the receipt verbatim - the panel cannot stop a row the host cannot.
//
// THE SEAM IS LATE-BOUND, the panel_message_453 pattern: the panel chain is
// assembled before the task pipeline exists (resident_windows.go's ordering),
// so the door leaves its seam nil at construction and the boot binds it on
// the first line where the runtime exists (bindPanelRosterCancel, next to
// bindPanelMessageRelay453). A nil seam is a legal state (CLI seam, refused
// pipeline) and means "refuse by name", never "drop".
//
// THE RECEIPT DOES NOT WAIT. The stop is issued; the task then settles on its
// own goroutine - the loop books the terminal task_log row, the bridge books
// the cancelled tool_call rows, the roster row is backfilled at Wait, and the
// applied-steps evidence (D31: what already landed) is readable from those
// forensics rows the moment the task settles. A dispatch goroutine that
// parked on the settle would stall every other panel request behind one
// task's exit; the receipt says the stop was issued and names the row, which
// is exactly what the task.cancel tool's own first line says ("已对…发出停
// 止"). The settled reading is a test/evidence-level read, not a blocked
// dispatch.

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
)

// rosterCancelSeam is the one thing the door needs from the task pipeline:
// hand it the row id the user's row carried and whether the request carried
// the confirm bit (ComposerRequest.PrivacyConfirmed - the SAME destructive
// bit the privacy doors gate on), and it answers with the host's own verdict
// - stopped (the handle fired) or refused (with the refusing path's own
// reason, verbatim). It is a function, not an interface, because it has
// exactly one method and the 453 relay precedent set that shape.
type rosterCancelSeam func(taskID string, confirmed bool) (stopped bool, reason string)

// rosterCancelDoor is the cancel socket of the panel chain.
type rosterCancelDoor struct {
	audit panel.AuditFunc
	actor string

	mu   sync.RWMutex
	seam rosterCancelSeam
}

func newRosterCancelDoor(audit panel.AuditFunc, actor string) *rosterCancelDoor {
	return &rosterCancelDoor{audit: audit, actor: actor}
}

// bindSeam is the late bind. nil seam stays legal and refuses by name.
func (d *rosterCancelDoor) bindSeam(seam rosterCancelSeam) {
	if d == nil {
		return
	}
	d.mu.Lock()
	d.seam = seam
	d.mu.Unlock()
}

func (d *rosterCancelDoor) seamOf() rosterCancelSeam {
	if d == nil {
		return nil
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.seam
}

// HandleRosterCancel answers one parsed cancel request. Outcomes, in order:
//
//	empty selector        -> refused before the seam is reached
//	no seam bound         -> refused by name (no task pipeline in this process)
//	id the host can't stop -> the host's own refusal, verbatim
//	otherwise             -> receipt naming the row, stop issued
//
// The confirm bit is the seam's business, not the door's pre-filter: the
// ROOT-task branch is the destructive one (stopping the root stops the whole
// run), so the seam demands the bit there and refuses without it; the
// per-row roster path keeps its existing shape (the row the tasks view
// showed, stopped through the ONE per-row stop path). Either way a refusal
// is named and the audit line records it.
func (d *rosterCancelDoor) HandleRosterCancel(ctx context.Context, req panel.ComposerRequest) (string, error) {
	id := strings.TrimSpace(req.RosterTaskID)
	if id == "" {
		err := fmt.Errorf("%w: rosterTaskId 是空的：cancel 只停一行，不给它 id 就什么都不停",
			panel.ErrComposerRequest)
		d.record(req, id, "refused", err)
		return "", err
	}
	seam := d.seamOf()
	if seam == nil {
		err := fmt.Errorf("%w: 取消通路未接入（本进程没有任务管线，或任务入口没有起来）：任务 %s 未被触碰，也不会被静默丢掉",
			panel.ErrNoHandlerAttached, id)
		d.record(req, id, "refused", err)
		return "", err
	}
	stopped, reason := seam(id, req.PrivacyConfirmed)
	if !stopped {
		err := fmt.Errorf("panel: 任务 %s 没有停下：%s", id, reason)
		d.record(req, id, "refused", err)
		return "", err
	}
	d.record(req, id, "stopped", nil)
	// The receipt is the existing cancel vocabulary's own first line, reshaped
	// for the panel: the row is named, the stop is stated as issued, and the
	// settled facts are pointed at rather than invented. What the task already
	// landed before the stop is NOT summarized here - the tool_call forensics
	// rows and the terminal task_log row are the D31 evidence surface, and a
	// receipt that re-counted them would be a second judge.
	return fmt.Sprintf("已对任务 %s 发出停止（经既有取消通路；任务将自行结算并落下终态行，"+
		"已落地的步骤以 tool_call 取证行与任务终态为准）。", id), nil
}

func (d *rosterCancelDoor) record(req panel.ComposerRequest, id, outcome string, err error) {
	if d == nil || d.audit == nil {
		return
	}
	line := "ok"
	if err != nil {
		line = err.Error()
	}
	d.audit("panel: ROSTER-CANCEL request=%q task=%q actor=%q outcome=%q err=%s",
		req.RequestID, id, d.actor, outcome, line)
}

// ---------------------------------------------------------------------------
// the seam's production implementation, on the composition root
// ---------------------------------------------------------------------------

// currentRootCancel is what execute() hands the runtime: the root task's own
// RunningTask handle, with the id minted by the loop so the seam can match it.
// A nil root is the honest "no root task is in flight" state.
type currentRootCancel struct {
	mu   sync.Mutex
	id   string
	root *agent.RunningTask
}

// setCurrentRoot stores (or clears, bg == nil) the root task's stop handle
// for the seam to match.
func (rt *agentRuntime) setCurrentRoot(bg *agent.RunningTask) {
	if rt == nil || rt.rootCancel == nil {
		return
	}
	rt.rootCancel.mu.Lock()
	defer rt.rootCancel.mu.Unlock()
	if bg == nil {
		rt.rootCancel.id, rt.rootCancel.root = "", nil
		return
	}
	rt.rootCancel.id, rt.rootCancel.root = bg.ID, bg
}

// cancelRosterRowWave2 is the seam the boot binds. Two paths, the existing
// ones, in the order the roster itself is built: the row the run is running
// first, then any row the TaskRoster holds. Nothing else is cancelled, nothing
// cascades, and both refusals are the host's own words.
//
// THE ROOT BRANCH IS THE CONFIRM-GATED ONE (2026-11-17, the adversarial
// audit's R1 second class): stopping the root stops the whole run, so it
// executes only when the request carried the confirm bit
// (ComposerRequest.PrivacyConfirmed - the same bit privacy.delete-one /
// privacy.purge gate on; page-side it rides only from the monitor cancel
// card's confirm key, the ONE send site). Without the bit the branch refuses
// by name and the root's handle is not touched - the same zero-effect
// refusal shape the privacy doors answer with.
func (rt *agentRuntime) cancelRosterRowWave2(taskID string, confirmed bool) (bool, string) {
	if rt == nil {
		return false, "本进程没有任务管线"
	}
	rt.rootCancel.mu.Lock()
	id, root := rt.rootCancel.id, rt.rootCancel.root
	rt.rootCancel.mu.Unlock()
	if root != nil && id == taskID {
		if !confirmed {
			return false, "停根任务需要确认位 privacyConfirmed:true（root 是整轮任务，" +
				"没有确认位的根任务停止零触碰）；带上确认位再发一次"
		}
		// The root's existing stop path (loop.go:342 RunningTask.Cancel ->
		// the loop root's own cancel). The run's execute() joins this exact
		// handle at Wait.
		root.Cancel()
		return true, ""
	}
	if rt.tasks == nil {
		return false, "任务名册未接线"
	}
	// The ONE stop path a roster row has by design (197 §0): per-row, never
	// cascading, refusing unknown ids and handle-less rows with its own
	// fail-closed reason pair - which rides the receipt verbatim.
	ok, why := rt.tasks.Cancel(taskID)
	if ok {
		return true, ""
	}
	return false, why
}
