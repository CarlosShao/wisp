package approval

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/tools"
)

// CancellationReport is D31's cancel non-atomicity surface: what the user is
// told when a veto arrived AFTER execution had already started.
//
// The contract has one law - never render a cancelled call as a clean undo.
// AppliedSteps is the only evidence of what the call already did, and it is
// filled by the tool that did it (internal/tools.Result.AppliedSteps), which
// for the fs family is ticket 20 segment 2's temp+atomic-rename writer. This
// segment owns the TYPE, the aggregation and the wording; segment 2 owns the
// data. A report with zero applied steps therefore says so out loud instead of
// implying "nothing happened".
type CancellationReport struct {
	Tool          string
	CorrelationID string
	TaskID        string
	Level         string

	// Vetoed is true when a veto landed on this call at all.
	Vetoed bool
	// Channel names which veto channel fired, so the wording can attribute it.
	Channel Channel
	// StartedBeforeVeto is the D31 fact: execution had begun, so the veto is a
	// stop request, not a rewind.
	StartedBeforeVeto bool
	// AppliedSteps is verbatim from tools.Result.AppliedSteps.
	AppliedSteps []string
	// UnreportedSteps is true when a veto fired but the tool listed nothing.
	// That is the "may have produced side effects without saying so" case and
	// the reason this field exists rather than an empty slice.
	UnreportedSteps bool
	// Outcome is the tool's own verdict string for the audit line.
	Outcome string
	// Text is the user/LLM-visible report. Built by Text, never hand-set.
	Text string
}

// Bus is the seam a RUNNING tool uses to ask whether its own call was vetoed
// mid-flight, and to turn what it did into D31's report. It exists because
// tools.Gate's two callbacks are pre-execution only: once the bridge has
// handed off, nothing in the gate can reach into the tool.
//
// The contract for consumers (ticket 20 segment 2's fs.write/fs.trash/fs.move,
// and any later side-effecting builtin):
//
//  1. poll Vetoed(correlationID) at every step boundary that can be stopped;
//  2. on a hit, stop, and return Result{AppliedSteps: <what actually landed>};
//  3. hand that Result plus the Decision to Report and surface the Text.
//
// A tool that cannot stop between steps must still return its applied steps;
// the report then reads "已产生副作用，无法回退", which is the honest answer
// and the whole point of B1's correction away from "undo".
type Bus interface {
	// Vetoed reports the veto recorded for a started call, if any.
	Vetoed(correlationID string) (Veto, bool)
	// Started reports whether the L1 window for this id handed off to
	// execution (i.e. the call is live).
	Started(correlationID string) bool
	// Report builds the applied-steps report for a finished call.
	Report(d tools.Decision, res tools.Result) CancellationReport
	// Complete releases one call's tracking; the host defers it.
	Complete(correlationID string)
}

// Bus returns the D31 seam for the running tool. It is bound to this gate's
// veto bookkeeping and adds no authority: nothing on it can allow anything.
func (g *Gate) Bus() Bus { return bus{g: g} }

type bus struct{ g *Gate }

func (b bus) Vetoed(corr string) (Veto, bool) { return b.g.LateVeto(corr) }

func (b bus) Started(corr string) bool {
	b.g.mu.Lock()
	defer b.g.mu.Unlock()
	_, ok := b.g.running[corr]
	return ok
}

func (b bus) Report(d tools.Decision, res tools.Result) CancellationReport {
	return b.g.Report(d, res)
}

func (b bus) Complete(corr string) { b.g.Complete(corr) }

// Report assembles one call's applied-steps report. See Bus for the contract.
func (g *Gate) Report(d tools.Decision, res tools.Result) CancellationReport {
	corr := orDefaultText(d.CorrelationID, d.TaskID)
	rep := CancellationReport{
		Tool: d.Tool, CorrelationID: corr, TaskID: d.TaskID, Level: d.LevelString(),
		AppliedSteps: append([]string(nil), res.AppliedSteps...),
	}
	g.mu.Lock()
	if r := g.running[corr]; r != nil {
		rep.StartedBeforeVeto = true
		if r.veto != nil {
			rep.Vetoed = true
			rep.Channel = r.veto.Channel
		}
	}
	g.mu.Unlock()
	if len(rep.AppliedSteps) > 0 {
		// A tool that lists what it did has started, whatever the post-handoff
		// table says (it is bounded by MaxTracked and can evict). Trusting the
		// ledger here is what stops a call that wrote bytes from being rendered
		// as "never entered execution" - D31's named failure.
		rep.StartedBeforeVeto = true
	}

	rep.UnreportedSteps = rep.Vetoed && len(rep.AppliedSteps) == 0
	rep.Text = rep.TextFor()
	return rep
}

// TextFor renders the user-visible wording. Exported on the value so a host
// can re-render a stored report without the gate's state.
func (r CancellationReport) TextFor() string {
	var b strings.Builder
	switch {
	case !r.StartedBeforeVeto && !r.Vetoed:
		b.WriteString("本次调用未进入执行阶段")
	case r.Vetoed && len(r.AppliedSteps) > 0:
		b.WriteString("取消不是原子的（D31）：否决到达时执行已开始，以下步骤已生效，未自动回退")
	case len(r.AppliedSteps) > 0:
		// Steps landed with no veto on record: a task cancellation, a tool error
		// or a timeout. Still not a clean undo, and saying 「未进入执行阶段」
		// over a list of side effects would be the fake-clean-cancel wording D31
		// exists to prevent, so this case gets its own sentence.
		b.WriteString("执行已开始，以下步骤已生效，未自动回退（本次未收到否决记录）")
	default:
		b.WriteString("取消不是原子的（D31）：否决到达时执行已开始，但工具未上报其步骤，以下步骤可能已产生副作用")
	}
	if r.Vetoed && r.Channel != "" {
		b.WriteString("（否决通道：")
		b.WriteString(channelNames[r.Channel])
		b.WriteString("）")
	}
	for _, s := range r.AppliedSteps {
		b.WriteString("\n- 已执行：")
		b.WriteString(s)
	}
	if r.UnreportedSteps {
		b.WriteString("\n- 已执行：（无上报，按最坏情况处理）")
	}
	return b.String()
}

// String implements fmt.Stringer over TextFor, which is how this report reaches
// internal/tools: the bridge carries a fmt.Stringer (see tools.CancelBus)
// instead of defining a second applied-steps shape of its own.
func (r CancellationReport) String() string {
	if r.Text != "" {
		return r.Text
	}
	return r.TextFor()
}

// ---------------------------------------------------------------------------
// the tools.CancelBus adapter
// ---------------------------------------------------------------------------

// ToolsCancelBus wires this gate into the bridge's D31 seam (tools.Options
// .Cancel). The adapter lives here because internal/tools may not import this
// package - this one imports its Decision/Result types.
//
// Vetoed maps onto LateVeto deliberately: a tool only ever asks after the
// window handed off, which is precisely the case where a veto is a stop request
// and not a rewind.
func (g *Gate) ToolsCancelBus() tools.CancelBus { return toolsBus{g: g} }

type toolsBus struct{ g *Gate }

func (b toolsBus) Vetoed(corr string) bool { _, ok := b.g.LateVeto(corr); return ok }

func (b toolsBus) Started(corr string) bool { return b.g.Bus().Started(corr) }

func (b toolsBus) Report(d tools.Decision, res tools.Result) fmt.Stringer {
	return b.g.Bus().Report(d, res)
}

func (b toolsBus) Complete(corr string) { b.g.Bus().Complete(corr) }

// compile-time proof the adapter is the seam the bridge asks for.
var _ tools.CancelBus = toolsBus{}
