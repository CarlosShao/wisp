package approval

// Ticket 148 AC#1's ruler for Queue.replay (queue.go).
//
// WHY THIS CASE EXISTS. Ticket 146 made cloneDecision the copy primitive of the
// LiveApprovals door and left the sibling outlet - Queue.replay - doing
// `d := it.Dec`: a bare value copy of a struct that carries eight
// reference-typed slots, so the re-opened verdict shared one map and seven
// backings with the queue's own history record. The symptom of forgetting that
// is a caller of Gate.Replay rewriting an already-settled verdict in place -
// the record a replayed card is re-displayed from, and the only copy of what
// the user originally answered. Ticket 148 Ⓐ (owner-approved, that one line of
// that one function only) switched the line to cloneDecision; this file is the
// instrument that keeps it switched.
//
// HOW IT RINGS. Same three-way shape as ticket 146's probes, driven through
// q.replay itself: value fidelity (the returned Decision must equal the
// history record field for field, except the correlation id the queue
// deliberately re-stamps), backing isolation (no declared reference slot of the
// returned Decision may share memory with the history record - walked by
// DECLARED TYPE with the same fail-closed census the ticket-146 file owns), and
// the literal in-place-write assertion (clobber every slot through the returned
// value only, then read the history record back and require it unchanged).
// Revert queue.go's line to `d := it.Dec` and every section below goes red -
// that is the mutation reading recorded in docs/evidence/s1/148-r3.md.
//
// What is NOT claimed: the copy is one level deep, so the containers reachable
// through Params' values are still shared with the history record. That limit
// is cloneDecision's documented shape, it is registered in ledger A261
// (ticket 146's "one more level" recommendation), and its only outbound panel
// door was closed by ticket 148 Ⓑ: LiveApprovals does not carry Params at all.

import (
	"reflect"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// mustFindHistory returns the settled record q.replay reads from, white-box:
// no public path exposes q.history, and the whole point of this case is that
// the replay outlet used to hand out an alias into it.
func mustFindHistory(t *testing.T, q *Queue, corr string) *qitem {
	t.Helper()
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, it := range q.history {
		if it != nil && it.Corr == corr {
			return it
		}
	}
	t.Fatalf("test setup: history 里没有 %s（settle 没有走 keepForReplay，用例前提就没了）", corr)
	return nil
}

func Test148ReplayReturnsACloneNotTheHistoryRecord(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	want := sharedDecision()
	mustPush(t, q, want)

	// Settle through the real refusal funnel - the fail-closed direction, so
	// the item lands in q.history the way every settled item does (deliver ->
	// dropLocked keepForReplay), not through a planted state.
	if err := q.reject("corr-backing", "148 AC#1 用例拒绝"); err != nil {
		t.Fatalf("reject(corr-backing) 失败，用例前提就没了：%v", err)
	}

	d, fresh, err := q.replay("corr-backing")
	if err != nil {
		t.Fatalf("replay(corr-backing) 失败：%v", err)
	}
	if fresh != "corr-backing#replay" {
		t.Errorf("重放 corr=%q，期望 %q（一键重放要换新 correlation id，旧令牌不能带过来）", fresh, "corr-backing#replay")
	}

	hist := mustFindHistory(t, q, "corr-backing")

	// VALUE FIDELITY first: the clone must carry the verdict, minus the one
	// field replay is supposed to re-stamp. A copy that drops values (a
	// cloneBacking that allocates empty, a cloneParamsMap that nils values)
	// goes red HERE, before any clobber can be mistaken for a lost value.
	wantDec := hist.Dec
	wantDec.CorrelationID = d.CorrelationID
	if !reflect.DeepEqual(d, wantDec) {
		t.Fatalf("AC#1 RED (值保真): replay() 交出的 Decision 与 history 记录不逐字段相等（corr 重贴除外）。\n"+
			"返回值:    %+v\nhistory:  %+v\n"+
			"这就是 \"copies\" 的第二半没人钉的形状：不共用底层，却把值拷丢了。",
			d, wantDec)
	}

	// BACKING ISOLATION, walked by declared type with the same fail-closed
	// census ticket 146's file owns. A ninth reference field on tools.Decision
	// that cloneDecision (and therefore replay) forgets is named here, not
	// discovered by an exploit.
	found := typeCensus(t, reflect.TypeOf(d))
	expect := sorted(queueSlotPaths)
	if len(found) != len(expect) {
		t.Fatalf("replay 返回值里的引用槽位数 %d (%v)，期望 %d (%v)：tools.Decision 的引用字段表变了，"+
			"queueSlotPaths / cloneDecision 必须同一轮评审跟上。", len(found), found, len(expect), expect)
	}
	for i := range expect {
		if found[i] != expect[i] {
			t.Fatalf("引用槽位名册不符：第 %d 枚 got=%q want=%q（全量 got=%v）", i, found[i], expect[i], found)
		}
	}
	histV := reflect.ValueOf(hist.Dec)
	var shared []string
	for _, path := range found {
		a := slotValue(t, reflect.ValueOf(d), path)
		b := slotValue(t, histV, path)
		if b.IsNil() || b.Len() == 0 {
			t.Fatalf("fixture 未填引用槽 %q（nil 或空）：backing 指针对它失明，此检会退化成恒真", path)
		}
		if pa, pb := backingPointer(a), backingPointer(b); pa != 0 && pa == pb {
			shared = append(shared, path)
		}
	}
	if len(shared) != 0 {
		t.Errorf("AC#1 RED: replay() 返回值里这些槽位与 history 记录共用底层：%v。\n"+
			"queue.go 的 replay 行必须走 cloneDecision（票 148 Ⓐ 批准的那一行），不许退回 d := it.Dec。",
			shared)
	}

	// THE LITERAL WRITE, ticket 146 AC#2's shape re-driven through the replay
	// outlet: clobber every slot IN PLACE through the returned value only, then
	// read the history record back and require it unchanged.
	d.Params["argv"] = "clobbered"
	d.Params["injected-by-replayer"] = true
	copy(d.Args, []byte("XXXXXXXX"))
	d.RulesHit[0] = risk.RuleID("R9")
	d.Paths[0] = "C:/clobbered"
	d.Capabilities[0] = tools.CapNotify
	d.Blacklist.Absolute[0] = "C:/clobbered"
	d.Blacklist.Unlockable[0] = "C:/clobbered"
	d.Blacklist.AlreadyUnlocked[0] = "C:/clobbered"

	have := hist.Dec
	if _, injected := have.Params["injected-by-replayer"]; injected {
		t.Errorf("AC#1 RED: history 记录的 Params 里出现了重放方塞进的键，map 是共用的（argv 现值 %v）",
			have.Params["argv"])
	}
	if argv, ok := have.Params["argv"].(string); ok && argv == "clobbered" {
		t.Errorf("AC#1 RED: history 记录的 Params[\"argv\"] 被就地改成 %q", argv)
	}
	if string(have.Args[:8]) == "XXXXXXXX" {
		t.Errorf("AC#1 RED: history 记录的 Args 前 8 字节被就地改写")
	}
	if have.RulesHit[0] == risk.RuleID("R9") {
		t.Errorf("AC#1 RED: history 记录的 RulesHit[0] 被就地改写（L2 卡逐字念这枚，票 17 冻结注）")
	}
	if have.Paths[0] == "C:/clobbered" {
		t.Errorf("AC#1 RED: history 记录的 Paths[0] 被就地改写")
	}
	if have.Capabilities[0] == tools.CapNotify {
		t.Errorf("AC#1 RED: history 记录的 Capabilities[0] 被就地改写（能力位被改写就是权限面被改写）")
	}
	for name, slot := range map[string]string{
		"Blacklist.Absolute":        have.Blacklist.Absolute[0],
		"Blacklist.Unlockable":      have.Blacklist.Unlockable[0],
		"Blacklist.AlreadyUnlocked": have.Blacklist.AlreadyUnlocked[0],
	} {
		if slot == "C:/clobbered" {
			t.Errorf("AC#1 RED: history 记录的 %s 被就地改写 —— A/B 档证据被改，重放卡据此说话", name)
		}
	}

	// And the outlet stays a read: replaying must not consume the history row
	// or move the pending set.
	if len(mustListHistory(t, q)) != 1 {
		t.Errorf("AC#1 RED: 一次 replay 后 history 长度变了 —— 出口必须是只读")
	}
	if q.Depth() != 0 {
		t.Errorf("Depth()=%d，期望 0：replay 不重开队列，重开是调用方下一程的事", q.Depth())
	}
}

// mustListHistory is the read-back helper for the history-length assertion.
func mustListHistory(t *testing.T, q *Queue) []*qitem {
	t.Helper()
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]*qitem(nil), q.history...)
}
