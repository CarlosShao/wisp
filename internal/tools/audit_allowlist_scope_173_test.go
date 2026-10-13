package tools

// Ticket 173 — the audit line's `in_allowlist_scope` field is a THREE-state
// answer, and this file is the ruler for all three states.
//
// The defect (ticket 162's AC#5 acceptance run measured it, this ticket owns
// it): when NO path reached the judge - a blank path is dropped by pathArgs
// before C19 ever sees it, a tool like task.output has no path parameter at
// all, or the call was refused before C19 ran - the line still printed
// `in_allowlist_scope=true`. "Nothing was judged" was rendered as "judged, and
// in range", the exact sentence a forensics reader filters on to find
// out-of-bounds calls. The old comment called it "the audit line's cheap
// answer"; the cheapness was never registered.
//
// The chosen third value is the literal `unknown`, not an omitted column:
//   - every `tools: call` line keeps ONE shape, so a reader (or a filter) never
//     has to distinguish "field absent because nothing was judged" from "field
//     absent because this line predates the fix / a log got truncated";
//   - the same line already prints empty-valued columns rather than dropping
//     them (an allowed call carries `reason=""`), so a token is the local
//     convention;
//   - `unknown` is a positive, greppable claim: "no verdict about scope
//     exists", which is exactly true on these calls.
//
// `false` was refused on purpose (票面原话): it says "judged, and out of
// range", a second false sentence - and it would paint ping/task.output-style
// pathless calls as "out of bounds", drowning the field's real signal.
//
// Contract axis (AC#3, checked before writing a byte): the literal
// `in_allowlist_scope` appears in NO frozen artifact - not the SPEC-02 §3
// `tool_call` DDL (its columns are listed at docs/specs/SPEC-02-data-storage.md
// §3 and this field is not one of them), not docs/PLAN.md, not docs/specs/**.
// The string is built by this package's own book() log call through
// Options.Logf, and the "rendered verbatim" promise in Decision's doc applies
// to RulesHit/Reason on the confirmation card, not to this line. So the field
// is code-level, not contract-level: fixable here without touching any spec.

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// TestAuditAllowlistScopeDoesNotInventAVerdict carries all four readings of
// ticket 173's AC#1 + AC#2. Each subtest owns its bridge and gate polarity so
// the audit line it inspects is the only `tools: call` line of its kind.
//
//	unjudged, blank path  -> must NOT be true; today (unfixed) it IS true (AC#1)
//	unjudged, no path tool-> must NOT be true (the ping/task.output family)
//	judged, in range      -> must stay true   (AC#2①, one word unchanged)
//	judged, out of range  -> must stay false  (AC#2②, one word unchanged)
func TestAuditAllowlistScopeDoesNotInventAVerdict(t *testing.T) {
	t.Run("unjudged: a blank path never reaches the judge", func(t *testing.T) {
		root := sealableTempCanonical124(t)
		var audit []string
		// The gate says yes so the call reaches fs.edit's own refusal - the
		// exact shape ticket 162's acceptance run measured (the judge never saw
		// a path; rules_hit=[R1] from the declared floor only).
		b, _ := fsEditBridgeAuditedPair(t, root, &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}, &audit)

		if out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, "   ",
			map[string]any{"old": "a", "new": "b"}))); err != nil {
			t.Fatalf("Execute: %v", err)
		} else if !out.IsError {
			t.Fatalf("a path C26 cannot answer must never report success: %+v", out)
		}
		line := r4AuditLine(t, audit, "fs.edit", "risk=L2")
		t.Logf("AC#1 读数（空白 path）：%s", line)
		if strings.Contains(line, "in_allowlist_scope=true") {
			t.Errorf("nothing was judged about a path, yet the line claims scope=true: %s", line)
		}
		if !strings.Contains(line, "in_allowlist_scope=unknown") {
			t.Errorf("the no-judgment state must render as unknown, got: %s", line)
		}
	})

	t.Run("unjudged: a tool with no path parameter at all", func(t *testing.T) {
		root := sealableTempCanonical124(t)
		var audit []string
		b, _ := fsEditBridgeAuditedPair(t, root, &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}, &audit)
		if err := b.reg.Register(entry("probe.nopath", risk.L0, nil, nil)); err != nil {
			t.Fatal(err)
		}

		out, err := b.Execute(t.Context(), req("probe.nopath", `{}`))
		if err != nil || out.IsError {
			t.Fatalf("pathless L0 fixture must execute: out=%+v err=%v", out, err)
		}
		line := r4AuditLine(t, audit, "probe.nopath", "outcome=success")
		t.Logf("AC#1 读数（无 path 参数工具）：%s", line)
		if strings.Contains(line, "in_allowlist_scope=true") {
			t.Errorf("a tool whose declaration has no path parameter cannot be in scope: %s", line)
		}
		if !strings.Contains(line, "in_allowlist_scope=unknown") {
			t.Errorf("the no-judgment state must render as unknown, got: %s", line)
		}
	})

	t.Run("judged in range stays true", func(t *testing.T) {
		root := sealableTempCanonical124(t)
		var audit []string
		b, _ := fsEditBridgeAuditedPair(t, root, &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}, &audit)

		f := writeUnder(t, root, "in.txt", "hello")
		if out, err := b.Execute(t.Context(), req("fs.read", argsFor(f))); err != nil || out.IsError {
			t.Fatalf("in-allowlist read: out=%+v err=%v", out, err)
		}
		line := r4AuditLine(t, audit, "fs.read", "outcome=success")
		t.Logf("AC#2① 读数（真在范围内）：%s", line)
		if !strings.Contains(line, "in_allowlist_scope=true") {
			t.Errorf("a path the judge vouched for must keep saying true: %s", line)
		}
	})

	t.Run("judged out of range stays false", func(t *testing.T) {
		root := sealableTempCanonical124(t)
		outside := sealableTempCanonical124(t)
		var audit []string
		// A human says no on the card: the R2 family's refusal shape, the same
		// one fs_edit_ac5_gate_r4_test.go pins side by side.
		b, _ := fsEditBridgeAuditedPair(t, root, &gateSpy{approveAns: AnswerReject, windowAns: AnswerVeto}, &audit)

		out, err := b.Execute(t.Context(), req("fs.read", argsFor(writeUnder(t, outside, "out.txt", "x"))))
		if err != nil {
			t.Fatal(err)
		}
		if !out.IsError {
			t.Fatalf("an out-of-scope read the human refused must be an error outcome: %+v", out)
		}
		line := r4AuditLine(t, audit, "fs.read", "risk=L2", "R2")
		t.Logf("AC#2② 读数（真越界，R2 族）：%s", line)
		if !strings.Contains(line, "in_allowlist_scope=false") {
			t.Errorf("a judged out-of-range path must keep saying false: %s", line)
		}
	})
}

// Test173ApprovalTimeoutKeepsTheTimeoutFlavor is AC#4's appended second item
// (same family, different field): the L2 approval timeout and the manual
// rejection share ONE D37 class - `user_rejected`. That sharing is the frozen
// model's own bucket, not a defect this ticket may fix: PLAN.md:2790 defines
// the class as "审批拒绝/超时" and D37(a)'s taxonomy holds 17 classes with no
// timeout class, so splitting it would be a contract change (human approval).
//
// The ticket offered two forms: split the classes, or carry the "timeout"
// flavor in the receipt. The second form is the one chosen, and this test is
// its criterion - the cause is readable from the record even though the class
// is shared: the receipt Text names the timeout, and the audit line separates
// the two dispositions in decision= (timeout vs reject) and kind=
// (approval-timeout vs refused). No production byte changed for this item; the
// class assertion below is exactly what a future approved class split would
// have to update.
func Test173ApprovalTimeoutKeepsTheTimeoutFlavor(t *testing.T) {
	readCall := func(t *testing.T, ans Answer) (agent.ToolOutcome, string) {
		t.Helper()
		root := sealableTempCanonical124(t)
		var audit []string
		b, _ := fsEditBridgeAuditedPair(t, root, &gateSpy{approveAns: ans, windowAns: AnswerVeto}, &audit)
		if err := b.reg.Register(entry("probe.irrev", risk.L2,
			[]Capability{CapFSWrite}, []Capability{CapFSWrite})); err != nil {
			t.Fatal(err)
		}
		out, err := b.Execute(t.Context(), req("probe.irrev", `{}`))
		if err != nil {
			t.Fatalf("a refusal is not a host fault: %v", err)
		}
		return out, r4AuditLine(t, audit, "probe.irrev", "outcome=error")
	}

	timeoutOut, timeoutLine := readCall(t, AnswerTimeout)
	rejectOut, rejectLine := readCall(t, AnswerReject)
	t.Logf("AC#4 第二枚读数（超时）：class=%q 回执=%q 审计行=%s", timeoutOut.ErrorClass, timeoutOut.Text, timeoutLine)
	t.Logf("AC#4 第二枚读数（人工拒）：class=%q 回执=%q 审计行=%s", rejectOut.ErrorClass, rejectOut.Text, rejectLine)

	if timeoutOut.ErrorClass != "user_rejected" || rejectOut.ErrorClass != "user_rejected" {
		t.Fatalf("D37 buckets both under user_rejected (PLAN.md:2790); got timeout=%q reject=%q",
			timeoutOut.ErrorClass, rejectOut.ErrorClass)
	}
	if !strings.Contains(timeoutOut.Text, "超时") {
		t.Errorf("the timeout receipt must name the timeout: %q", timeoutOut.Text)
	}
	if !strings.Contains(timeoutLine, "kind=approval-timeout") || !strings.Contains(timeoutLine, "decision=timeout") {
		t.Errorf("the timeout audit line must separate itself from a human no: %s", timeoutLine)
	}
	if !strings.Contains(timeoutLine, "decision="+agent.DecisionTimeout) {
		t.Errorf("the booked decision column must keep the timeout value: %s", timeoutLine)
	}
	if !strings.Contains(rejectLine, "kind=refused") || !strings.Contains(rejectLine, "decision="+agent.DecisionReject) {
		t.Errorf("a human no must stay kind=refused decision=reject: %s", rejectLine)
	}
}
