package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// plan.present (ticket 165): the "先出方案、你点头再动手" brake the voice entry
// needs. 甲路, owner-approved 2026-09-27 (A322/A323), frozen as the D34 row at
// PLAN.md:2567 and mirrored at SPEC-07 §3:73:
//
//	plan.present | 交回一份"我要干什么"的清单并停下等点 | L2 | S3
//
// 甲路 means exactly this: NO new state, NO new approval kind, NO new gate
// method. The tool declares L2, so C19 R1 hands the call the L2 verdict it
// declares and the bridge's EXISTING L2 route raises the EXISTING C18 card -
// D43 transitions #17 (Acting -> AwaitingApproval), #23/#24 (native allow /
// reject + the 300s timeout) and #25 are the four rows the frozen text
// reuses. Nothing in this file, and nothing the tool can return, adds a state,
// a card kind or an answer channel.
//
// What this tool deliberately does NOT do, and why each omission is the
// security half:
//
//   - it writes nothing, holds no scope and mints no grant. SPEC-06 §8.3's
//     first bullet says an L2 never enters any persistent authorization (D45-3
//     restates it), so a plan approval can never become "all later writes are
//     approved" - the decision column for THIS call is what the click bought.
//     A write that the plan did not list keeps asking, one card at a time;
//     PLAN.md:2567 states the same consequence as a warning: this row does not
//     reduce the number of confirmations.
//   - it has no capability face at all (Capabilities and Needs both empty, the
//     way the frozen row's fourth column reads): it reads no host table, so it
//     cannot widen C3.
//   - it has no PathParams: the plan's steps are prose for a human, not paths,
//     so C19 has no path of this call to judge and C26 nothing to canonicalize.
//     A step that names a file is text the human reads on the card; the write
//     that later touches that file is judged as its own call, on its own path.
//
// The card the human sees carries the arguments verbatim (Decision.Params is
// what the C18 queue renders), which is why the steps are plain sentences and
// why the scope clause below is IN the returned text and not only in a comment:
// ticket 165 AC#3② names that sentence as the ingredient whose removal must
// make an out-of-scope write impossible to redden.
// ---------------------------------------------------------------------------

type planPresent struct{}

type planPresentArgs struct {
	// Summary is optional; it becomes the one-line "what this batch is for".
	Summary string `json:"summary"`
	// Steps is the list itself. It is prose for a human, one step per item.
	Steps []string `json:"steps"`
}

// planPresentSchema is the C1 Parameters document. There is deliberately no
// field a caller could use to mark anything approved: no "approved", no
// "apply", no "scope". The only way this call runs is the existing L2 card,
// and the only thing the click covers is this call (AC#4).
var planPresentSchema = JSONSchema(`{"type":"object","properties":{` +
	`"summary":{"type":"string","description":"一句话：这批改动要达成什么（给人读的）"},` +
	`"steps":{"type":"array","minItems":1,"items":{"type":"string"},` +
	`"description":"你打算做的每一步，一句一步；确认卡上给人读的就是这份清单"}` +
	`},"required":["steps"],"additionalProperties":false}`)

func (planPresent) Name() string { return "plan.present" }

// Description says what the tool is and, in the same breath, what the click
// does not buy: the boundary sentence is the security half, and a description
// that only sold the capability would hand the model a power it does not have.
func (planPresent) Description() string {
	return "先交一份「我要干什么」的清单并停下等点头：" +
		"这份清单作为一张 L2 确认卡给人读，点头只批清单里的这几步、只批这一批；" +
		"清单之外的写仍会各自弹卡，不会顺带被放行"
}

func (planPresent) Parameters() JSONSchema { return planPresentSchema }

// Execute implements Tool. It produces the checklist text and reports that it
// applied nothing: by the time it runs, the human already answered the L2 card
// the bridge raised for this very call, and the text is what the model reads
// back as the approved batch.
//
// A caller that reaches Execute WITHOUT the bridge (a unit test, a future host
// path) gets the same text - which is why the text never claims "you already
// nodded". It states the plan and the boundary; the asking is the gate's job.
func (planPresent) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a planPresentArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	steps, err := planPresentSteps(a)
	if err != nil {
		return Result{Text: err.Error(), IsError: true}, nil
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("清单 %d 步", len(steps)))
	}
	return Result{Text: planPresentText(strings.TrimSpace(a.Summary), steps)}, nil
}

// PlanPresentDecl is the host-side declaration for plan.present.
//
// Declared L2 because that is the frozen row's own number: PLAN.md:2567 prints
// "L2" and SPEC-07 §3 mirrors it. The consequence is the ticket's whole point -
// R1 escalates the call to L2, and the L2 route is the card, so the batch is
// never run off a 2-3 second unopposed L1 window (a window is exactly what a
// voice entry misheard cannot take back).
//
// Capabilities, Needs and PathParams are all empty for the reasons in this
// file's header. Resident is true: the model has to know this door exists to
// knock on it, and a plan tool that is only reachable by retrieval is a door
// nobody feeds (the failure ticket 161 was written about).
func PlanPresentDecl() Decl {
	return Decl{
		Capabilities: nil,
		Needs:        nil,
		Declared:     risk.L2,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinPlanEntries returns the plan family (today exactly one tool). It is a
// separate builder for the same reason BuiltinSubagentEntries is: it needs no
// host dependency - no roster, no canonicalizer, no file handle - so the
// composition root can register it without arranging anything first.
func BuiltinPlanEntries() []Entry {
	return []Entry{{Tool: planPresent{}, Decl: PlanPresentDecl()}}
}

// planPresentSteps normalizes the list and refuses the empty one. An empty
// checklist is not a plan: it would put a card in front of a human whose only
// content is the boundary sentence, and "I will do something" with no steps is
// precisely the ambiguous request this tool exists to turn into a list. The
// refusal is model-actable (IsError), not a host failure.
func planPresentSteps(a planPresentArgs) ([]string, error) {
	out := make([]string, 0, len(a.Steps))
	for _, s := range a.Steps {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("清单是空的：plan.present 要一份至少一步的清单（空清单不是一份方案）")
	}
	return out, nil
}

// planPresentText renders the checklist the model reads back. The steps are
// numbered so a human can point at one ("第 3 步不算"), and the batch size is
// echoed because the scope clause names it.
func planPresentText(summary string, steps []string) string {
	var b strings.Builder
	if summary != "" {
		b.WriteString("这批改动要达成：" + summary + "\n")
	}
	fmt.Fprintf(&b, "我要做的事（共 %d 步）：\n", len(steps))
	for i, s := range steps {
		fmt.Fprintf(&b, "%d. %s\n", i+1, s)
	}
	b.WriteString(planPresentScopeClause(len(steps)))
	return b.String()
}

// planPresentScopeClause is the sentence that carries the BOUNDARY of an
// approval: what was approved is this list, this size, this once - and every
// write outside it still asks on its own.
//
// It is a function of the batch size on purpose (AC#3): a generic sentence
// ("本次只批这一批") can stay in the text while the plan it belongs to changes
// shape, and a ruler reading it would still pass. Naming the count ties the
// sentence to the list it bounds, so the instrument in
// plan_present_165_test.go reddens when either the sentence or the count stops
// matching the list the human was shown.
func planPresentScopeClause(n int) string {
	return fmt.Sprintf("批准范围只限上面这 %d 步、只限这一批："+
		"清单之外的任何写（包括随后想到的新步骤、以及之后的每一条命令）都仍会各自弹卡，"+
		"不会因为这份清单被点头而顺带放行。", n)
}
