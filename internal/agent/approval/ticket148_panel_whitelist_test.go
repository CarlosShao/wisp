package approval

// Ticket 148 Ⓑ's reconciliation ruler for the outbound field whitelist
// (pending_read.go's panelVerdict).
//
// WHAT THE OWNER APPROVED. Ticket 148 Ⓑ narrows the outward field range to
// "the fields the panel needs", by name. The name of that set is not a guess:
// the one production reader of LiveApprovals is cmd/wisp/panel_pump.go's
// liveVerdicts, and it reads exactly six fields off the row's Decision - Tool,
// Args (through panelArgs), Level, RulesHit, Reason, SessionOverrideBlocked -
// plus CorrelationID and Position, which travel on the LiveApproval wrapper.
// panelVerdict carries those six and nothing else; CorrelationID stays on the
// wrapper, so a row whose Decision still carries an id string is itself a
// leak.
//
// THE TWO WAYS THIS FILE GOES RED (the ticket's "字段增删即红"):
//
//	ruler 1  the ROSTER check. It reflects over tools.Decision's declared
//	         fields and requires them to be exactly allowed ∪ dropped, both
//	         named above. A field ADDED to tools.Decision is in neither list ->
//	         red, and whoever adds it must decide in the same review whether
//	         the panel needs it (allowed) or the queue keeps it (dropped). A
//	         field REMOVED, or a whitelist line renamed, is a mismatch too. No
//	         drift in either direction can pass silently.
//	ruler 2  the ENFORCEMENT check, against the real queue. One admitted
//	         verdict with every field filled (so a leak cannot hide behind an
//	         empty source), one LiveApprovals call: every allowed field must
//	         arrive with its value and its own backing (the two reference
//	         slots are cloneBacking copies, so an in-place write through them
//	         still cannot reach the queue record); every dropped field must be
//	         the zero value on the returned row. Swap panelVerdict back to
//	         cloneDecision - or delete the whitelist line for any one field -
//	         and ruler 2 goes red; that is the mutation reading recorded in
//	         docs/evidence/s1/148-r3.md.
//
// THE Params SECOND LEVEL, stated as disposed rather than forgotten. The
// containers a JSON decode leaves inside Params' values ([]any, nested
// map[string]any) were still shared after ticket 146's one-level clone.
// Ticket 148 Ⓑ closes that hole by exclusion: Params is not on the whitelist,
// so the panel payload does not carry it and nothing outside the queue can
// reach those containers through LiveApprovals any more. Queue.replay's
// returned Decision still carries a one-level Params clone for the in-process
// approval flow; its second level remains cloneDecision's documented limit and
// ledger A261's registered item, not this file's claim.

import (
	"reflect"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/tools"
)

// panelVerdictAllowedFields is the whitelist itself, named. Keep it in sync
// with panelVerdict's literal and with the panel reader it is a census of.
var panelVerdictAllowedFields = []string{
	"Args",
	"Level",
	"Reason",
	"RulesHit",
	"SessionOverrideBlocked",
	"Tool",
}

// panelVerdictDroppedFields names every other declared field of
// tools.Decision: the fields the queue keeps. The reference-typed ones among
// them (Params, Paths, Capabilities, Blacklist) are exactly the slots the
// ticket-146 clone guarded; on the panel payload they are gone outright.
var panelVerdictDroppedFields = []string{
	"Blacklist",
	"CallID",
	"Capabilities",
	"CorrelationID",
	"DecisionColumn",
	"Mode",
	"ModeKept",
	"ModeSilenced",
	"Params",
	"Paths",
	"PluginID",
	"Provider",
	"Refusal",
	"TaskID",
	"Timeout",
}

// fullDecision148 fills EVERY declared field of tools.Decision, so that ruler
// 2's "dropped fields arrive zero" assertion has teeth: a whitelist line that
// leaks a dropped field cannot hide behind an empty source value.
func fullDecision148() tools.Decision {
	d := sharedDecision()
	d.PluginID = "plugin-148"
	d.Refusal = tools.RefusalOutOfAllowlist
	d.ModeSilenced = true
	d.ModeKept = "mode-kept-148"
	d.DecisionColumn = "allow"
	d.CallID = "call-148"
	d.Timeout = 5 * time.Second
	return d
}

// decisionFieldRoster reflects the declared field names of tools.Decision.
func decisionFieldRoster(t *testing.T) []string {
	t.Helper()
	ty := reflect.TypeOf(tools.Decision{})
	out := make([]string, 0, ty.NumField())
	for i := 0; i < ty.NumField(); i++ {
		out = append(out, ty.Field(i).Name)
	}
	return sorted(out)
}

func unionSorted(a, b []string) []string {
	return sorted(append(append([]string(nil), a...), b...))
}

// Test148PanelVerdictWhitelistRosterIsClosed is ruler 1.
func Test148PanelVerdictWhitelistRosterIsClosed(t *testing.T) {
	have := decisionFieldRoster(t)
	want := unionSorted(panelVerdictAllowedFields, panelVerdictDroppedFields)
	if len(have) != len(want) {
		t.Fatalf("tools.Decision 现有 %d 个字段 (%v)，白名单两侧名册合计 %d 个 (%v)。\n"+
			"字段数不符说明 Decision 增删了字段而本票的白名单没跟上 —— 这条用例正是为此而红：\n"+
			"新字段面板要用的，进 panelVerdictAllowedFields 并在 panelVerdict 里携出；\n"+
			"面板不需要的，进 panelVerdictDroppedFields 留在队列里。两件事必须同一轮评审做掉，不许放宽。",
			len(have), have, len(want), want)
	}
	for i := range want {
		if have[i] != want[i] {
			t.Fatalf("字段名册不符：第 %d 枚 got=%q want=%q（全量 have=%v）", i, have[i], want[i], have)
		}
	}
}

// Test148LiveApprovalsCarriesOnlyWhitelistedFields is ruler 2.
func Test148LiveApprovalsCarriesOnlyWhitelistedFields(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	want := fullDecision148()
	mustPush(t, q, want)
	stored := mustFindLive(q, "corr-backing").Dec

	rows := q.LiveApprovals()
	if len(rows) != 1 {
		t.Fatalf("期望 1 行，得 %d：%+v", len(rows), rows)
	}
	row := rows[0]
	if row.CorrelationID != "corr-backing" {
		t.Errorf("包装层 CorrelationID=%q，期望队列签发的 %q", row.CorrelationID, "corr-backing")
	}

	storedV := reflect.ValueOf(stored)
	rowV := reflect.ValueOf(row.Decision)
	allowed := map[string]bool{}
	for _, name := range panelVerdictAllowedFields {
		allowed[name] = true
	}

	for _, name := range decisionFieldRoster(t) {
		src := storedV.FieldByName(name)
		if !src.IsValid() {
			t.Fatalf("字段 %q 在存储记录上取不到，尺子自己坏了", name)
		}
		got := rowV.FieldByName(name)
		if allowed[name] {
			// Fidelity on the carried set: the whitelist is a projection, not a
			// reshaper. A field the panel needs that arrives zero (or whose
			// backing is a fresh empty) is a deformed verdict.
			if src.IsZero() {
				t.Fatalf("fixture 未填白名单字段 %q：携带集断言对它是失明的，fullDecision148 必须填满", name)
			}
			if got.IsZero() {
				t.Errorf("Ⓑ RED: 白名单字段 %q 在返回行上是零值，面板会读到一张缺字段的卡", name)
				continue
			}
			if !reflect.DeepEqual(got.Interface(), src.Interface()) {
				t.Errorf("Ⓑ RED: 白名单字段 %q 的值在投影途中变了：got=%v want=%v",
					name, got.Interface(), src.Interface())
			}
		} else {
			// Enforcement on the dropped set: zero on the row, whatever the
			// source carried. This is the assertion a leaked field trips.
			if !got.IsZero() {
				t.Errorf("Ⓑ RED: 非白名单字段 %q 出现在返回行上（值 %v）——"+
					"panelVerdict 携出了白名单外的字段，或白名单行被人删了", name, got.Interface())
			}
		}
	}

	// BACKING ISOLATION for the two reference slots the whitelist carries: the
	// reader may write in place, the queue's record may not change.
	if pa, pb := backingPointer(reflect.ValueOf(row.Decision.Args)), backingPointer(reflect.ValueOf(stored.Args)); pa != 0 && pa == pb {
		t.Errorf("Ⓑ RED: 返回行的 Args 与队列记录共用底层（panelVerdict 必须经 cloneBacking 拷出）")
	}
	if pa, pb := backingPointer(reflect.ValueOf(row.Decision.RulesHit)), backingPointer(reflect.ValueOf(stored.RulesHit)); pa != 0 && pa == pb {
		t.Errorf("Ⓑ RED: 返回行的 RulesHit 与队列记录共用底层（panelVerdict 必须经 cloneBacking 拷出）")
	}

	// THE LITERAL WRITE, through the carried slots only.
	row.Decision.Args[0] = 'X'
	if len(row.Decision.RulesHit) > 0 {
		row.Decision.RulesHit[0] = "R9-clobber"
	}
	if string(stored.Args[:1]) == "X" {
		t.Errorf("Ⓑ RED: 队列记录的 Args 被面板读者就地改写")
	}
	if len(stored.RulesHit) > 0 && stored.RulesHit[0] == "R9-clobber" {
		t.Errorf("Ⓑ RED: 队列记录的 RulesHit 被面板读者就地改写（L2 卡逐字念这枚）")
	}

	// A second read must describe the queue, not the earlier reader's write.
	again := q.LiveApprovals()
	if len(again) != 1 {
		t.Fatalf("第二次读期望 1 行，得 %d", len(again))
	}
	if string(again[0].Decision.Args) != string(sharedArgs148) {
		t.Errorf("Ⓑ RED: 第二次读到的 Args 带着上一枚读者的改写")
	}
}

// sharedArgs148 pins the fixture's Args bytes for the second-read assertion.
var sharedArgs148 = fullDecision148().Args
