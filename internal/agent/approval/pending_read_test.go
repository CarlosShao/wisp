package approval

// White-box read-side cases for LiveApprovals (internal/agent/approval/
// pending_read.go, ticket 35's snapshot pump). F-PUMP-2 of
// docs/evidence/s1/35-panel-snapshot-pump-r1-accept-r1.md §5.2 asked for one
// case in this package, because before these two functions the whole file could
// be deleted and `go test ./internal/agent/approval/` still said ok: the only
// caller in the tree is cmd/wisp, one package out.
//
// Why this file is `package approval` and not `package approval_test` like the
// other eight: the queue's own settle paths (deliver -> dropLocked, and the
// grantNonce error path) remove an item from q.pending in the same locked
// section where they change its state, so NO public path can leave a
// non-pending row inside q.pending. The row filter in pending_read.go is
// therefore only observable from inside the package - see
// TestLiveApprovalsSkipsRowsThatAreNotPending for the case that plants one,
// and the honest note there about what that does and does not prove.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// readSideDecision builds one admitted verdict with every field the L2 card
// reads off Decision set to a value that differs per item, so a returned row
// that lost or swapped a field cannot pass by accident.
func readSideDecision(corr, tool string, level risk.Level, rules []risk.RuleID, blocked bool, mode risk.Mode) tools.Decision {
	return tools.Decision{
		CorrelationID:          corr,
		TaskID:                 "task-" + corr,
		Tool:                   tool,
		Level:                  level,
		RulesHit:               rules,
		Reason:                 "reason-" + corr,
		Paths:                  []string{"C:/dir/" + corr},
		SessionOverrideBlocked: blocked,
		Mode:                   mode,
		ModeSilenced:           !blocked,
	}
}

// mustPush admits one decision through the queue's own admission path.
func mustPush(t *testing.T, q *Queue, d tools.Decision) *qitem {
	t.Helper()
	it, err := q.push(d)
	if err != nil {
		t.Fatalf("push(%s) 失败，用例前提就没了：%v", d.CorrelationID, err)
	}
	return it
}

func assertRow(t *testing.T, got LiveApproval, want tools.Decision, position int) {
	t.Helper()
	if got.CorrelationID != want.CorrelationID {
		t.Errorf("CorrelationID=%q，期望 %q", got.CorrelationID, want.CorrelationID)
	}
	if got.Position != position {
		t.Errorf("%s 的 Position=%d，期望 FIFO 第 %d 位（不是登记时的序号）",
			want.CorrelationID, got.Position, position)
	}
	d := got.Decision
	if d.Tool != want.Tool {
		t.Errorf("%s 的 Decision.Tool=%q，期望原样 %q", want.CorrelationID, d.Tool, want.Tool)
	}
	if d.CorrelationID != want.CorrelationID || d.TaskID != want.TaskID {
		t.Errorf("%s 的行把 id 串到了别枚上：corr=%q task=%q",
			want.CorrelationID, d.CorrelationID, d.TaskID)
	}
	if d.Level != want.Level {
		t.Errorf("%s 的 Level=%v，期望原样 %v（等级是卡片上那枚徽标）",
			want.CorrelationID, d.Level, want.Level)
	}
	if d.LevelString() != want.LevelString() {
		t.Errorf("%s 的 LevelString=%q，期望原样 %q",
			want.CorrelationID, d.LevelString(), want.LevelString())
	}
	if len(d.RulesHit) != len(want.RulesHit) {
		t.Fatalf("%s 的 RulesHit=%v，期望原样 %v（命中规则是 L2 卡要逐条念出来的）",
			want.CorrelationID, d.RulesHit, want.RulesHit)
	}
	for i := range want.RulesHit {
		if d.RulesHit[i] != want.RulesHit[i] {
			t.Errorf("%s 的 RulesHit[%d]=%q，期望原样 %q",
				want.CorrelationID, i, d.RulesHit[i], want.RulesHit[i])
		}
	}
	if d.Reason != want.Reason {
		t.Errorf("%s 的 Reason=%q，期望原样 %q", want.CorrelationID, d.Reason, want.Reason)
	}
	if len(d.Paths) != len(want.Paths) || d.Paths[0] != want.Paths[0] {
		t.Errorf("%s 的 Paths=%v，期望原样 %v", want.CorrelationID, d.Paths, want.Paths)
	}
	if d.SessionOverrideBlocked != want.SessionOverrideBlocked {
		t.Errorf("%s 的 SessionOverrideBlocked=%v，期望原样 %v（R4 那枚标记，卡片据此说话）",
			want.CorrelationID, d.SessionOverrideBlocked, want.SessionOverrideBlocked)
	}
	if d.Mode != want.Mode || d.ModeSilenced != want.ModeSilenced {
		t.Errorf("%s 的 Mode=%v ModeSilenced=%v，期望原样 %v/%v",
			want.CorrelationID, d.Mode, d.ModeSilenced, want.Mode, want.ModeSilenced)
	}
}

// TestLiveApprovalsReportsPendingRowsAndDoesNotConsumeThem is the real-path
// case: three decisions admitted through the queue's own admission path, one
// of them then settled by a real refusal, and the read asserted against what is
// left. It also pins that the read is a read: Depth() is the number the ball's
// badge shows, so a LiveApprovals call that consumed or reordered anything
// would be a user-visible lie.
func TestLiveApprovalsReportsPendingRowsAndDoesNotConsumeThem(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	dA := readSideDecision("corr-a", "fs.write", risk.L2, []risk.RuleID{risk.R2, risk.R3}, true, risk.ModeAskEveryStep)
	dB := readSideDecision("corr-b", "shell.run", risk.L1, []risk.RuleID{risk.R1}, false, risk.ModeAskHighRisk)
	dC := readSideDecision("corr-c", "net.fetch", risk.L2, []risk.RuleID{risk.R7}, true, risk.ModeAutoApprove)
	mustPush(t, q, dA)
	mustPush(t, q, dB)
	mustPush(t, q, dC)

	before := q.Depth()
	got := q.LiveApprovals()
	if len(got) != 3 {
		t.Fatalf("三枚 pending 时返回 %d 行，期望 3 行：%+v", len(got), got)
	}
	// Head first, FIFO, with the position the depth badge uses.
	assertRow(t, got[0], dA, 1)
	assertRow(t, got[1], dB, 2)
	assertRow(t, got[2], dC, 3)
	if after := q.Depth(); after != before {
		t.Errorf("读一次就把 Depth() 从 %d 变成了 %d，这不是只读", before, after)
	}
	// Twice-in-a-row is what the pump actually does: nothing may drain.
	if again := q.LiveApprovals(); len(again) != 3 {
		t.Errorf("第二次读返回 %d 行，期望仍 3 行", len(again))
	}

	// Settle the head through the real refusal funnel (no forged state).
	if err := q.reject("corr-a", "用户拒绝了本次操作"); err != nil {
		t.Fatalf("reject(corr-a) 失败：%v", err)
	}
	got = q.LiveApprovals()
	if len(got) != 2 {
		t.Fatalf("结掉队头后返回 %d 行，期望 2 行：%+v", len(got), got)
	}
	// The settled row must be gone AND the survivors must be renumbered:
	// Position is the place in the FIFO right now, not the sequence number the
	// queue handed out at admission.
	assertRow(t, got[0], dB, 1)
	assertRow(t, got[1], dC, 2)
	for _, row := range got {
		if row.CorrelationID == "corr-a" {
			t.Fatal("已决断的 corr-a 仍出现在返回里，面板会把它重画成还在等")
		}
	}
	if q.Depth() != 2 {
		t.Errorf("Depth()=%d，期望 2（返回的行数应与深度一致）", q.Depth())
	}
	if len(q.LiveApprovals()) != q.Depth() {
		t.Errorf("行数 %d 与 Depth() %d 不一致，同一枚队列两套口径",
			len(q.LiveApprovals()), q.Depth())
	}

	// Empty queue reads as an empty list, not as nil-with-a-length.
	if err := q.reject("corr-b", ""); err != nil {
		t.Fatalf("reject(corr-b) 失败：%v", err)
	}
	if err := q.reject("corr-c", ""); err != nil {
		t.Fatalf("reject(corr-c) 失败：%v", err)
	}
	if got = q.LiveApprovals(); len(got) != 0 {
		t.Errorf("全部结掉后返回 %d 行，期望 0 行：%+v", len(got), got)
	}
}

// TestLiveApprovalsSkipsRowsThatAreNotPending pins the filter in
// pending_read.go, and it is honest about how it gets there.
//
// What is real here: an item can sit in q.pending with stateDropped only if
// some future path changes state without going through dropLocked - which is
// exactly what the guard in pending_read.go exists for, and exactly the shape
// a mutation of that guard produces. Today no public path builds it (deliver
// sets the state and drops the row inside one locked section), so this case
// plants both shapes from inside the package rather than pretending the queue
// made them. Read it as "the guard is load-bearing": delete the `continue`, or
// the `it == nil` half of it, and this case goes red (the nil row turns the
// loop into a nil-pointer dereference, which is a panic, which is red).
//
// What it does NOT prove: that the queue ever reaches this state on its own.
// That is a different claim, and no case in this file makes it.
func TestLiveApprovalsSkipsRowsThatAreNotPending(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	dLive := readSideDecision("corr-live", "fs.write", risk.L2, []risk.RuleID{risk.R2}, false, risk.ModeAskHighRisk)
	mustPush(t, q, dLive)

	ghost := &qitem{Corr: "ghost-answered", Dec: dLive, state: stateAnswered}
	dropped := &qitem{Corr: "ghost-dropped", Dec: dLive, state: stateDropped}
	q.mu.Lock()
	q.pending = append(q.pending, ghost, dropped, nil)
	q.mu.Unlock()

	got := q.LiveApprovals()
	if len(got) != 1 {
		t.Fatalf("pending 里混进两枚非 pending 行与一枚 nil 行时返回 %d 行，期望只有那 1 枚活行：%+v",
			len(got), got)
	}
	assertRow(t, got[0], dLive, 1)
	if got[0].Position != q.position(mustFindLive(q, "corr-live")) {
		t.Errorf("Position=%d，与 position() 算出的位次 %d 不一致",
			got[0].Position, q.position(mustFindLive(q, "corr-live")))
	}
	// The planted rows are read-side invisible but they are still queue
	// members: the guard must skip them, not delete them.
	q.mu.Lock()
	sliceLen := len(q.pending)
	q.mu.Unlock()
	if sliceLen != 4 {
		t.Errorf("一次读之后 q.pending 长度为 %d，期望仍是 4（只读不得改动队列）", sliceLen)
	}
}

// TestLiveApprovalsOnNilQueueIsNil covers the guard's other half. Without the
// `q == nil` branch this is a nil-pointer dereference inside q.mu.Lock(), i.e.
// a panic, i.e. red - which is the point of writing it at all.
func TestLiveApprovalsOnNilQueueIsNil(t *testing.T) {
	var q *Queue
	if got := q.LiveApprovals(); got != nil {
		t.Errorf("nil 队列应返回 nil，得 %+v", got)
	}
}

// mustFindLive is a test-side lookup: it returns the row it was asked to find
// so the assertion above compares LiveApprovals against position() itself
// rather than against a literal the same file wrote.
func mustFindLive(q *Queue, corr string) *qitem {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, it := range q.pending {
		if it != nil && it.Corr == corr {
			return it
		}
	}
	panic("test setup: no pending row named " + corr)
}
