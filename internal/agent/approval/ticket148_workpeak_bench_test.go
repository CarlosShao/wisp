package approval

// Ticket 148 AC#2's before/after readings for "copying the eight reference
// slots" - the arm ticket 146 AC#3 ran out-of-tree and never committed, made
// reproducible in-tree here.
//
// THE THREE ARMS AND WHAT EACH ONE IS:
//
//	A1 bareValueCopyProjection148  the PRE-ticket-146 shape: `Decision: it.Dec`,
//	   every reference backing shared with the queue record. This is the
//	   "before" of "多拷那 8 枚槽位".
//	A2 fullCloneProjection148      the ticket-146..pre-148 shape: one
//	   cloneDecision per row (map + seven backings). This is what ticket 146's
//	   evidence measured as +4736 B/op / +72 allocs/op at depth 8, on
//	   out-of-tree snapshots.
//	A3 production LiveApprovals    the ticket-148 shape: panelVerdict's
//	   whitelist projection (two backings per row: Args, RulesHit; Params and
//	   the rest no longer leave the queue at all).
//
// The queue is filled to DefaultMaxPending with the full eight-slot fixture,
// because that is the depth ticket 146's readings were taken at and the bound
// the pump's cost is stated against (pending_read.go's COST note).
//
// WHAT THESE NUMBERS ARE AND ARE NOT. They are the mechanism account
// (B/op, allocs/op, ns/op) of the copy itself, measured on the subject ticket
// 148 AC#2 proved really runs tool calls (ticket148_workpeak_subject_test.go).
// They are NOT the production WorkPeak band reading: the `wisp slo` instrument
// still boots the sit-still skeleton (cmd/wisp), which is off this leg's write
// surface, so that gap stays registered - the evidence file quotes these arms
// under that explicit caveat, per ticket 146 §4.3's honesty rule.

import (
	"fmt"
	"testing"
)

// bench148Queue fills the FIFO to DefaultMaxPending with the full fixture.
func bench148Queue(b *testing.B) *Queue {
	b.Helper()
	q := NewQueue(0, 0, 0, nil)
	for i := 0; i < DefaultMaxPending; i++ {
		d := sharedDecision()
		d.CorrelationID = fmt.Sprintf("corr-bench-%d", i)
		d.TaskID = fmt.Sprintf("task-bench-%d", i)
		if _, err := q.push(d); err != nil {
			b.Fatalf("push(%s) 失败，bench 前提就没了：%v", d.CorrelationID, err)
		}
	}
	return q
}

// bareValueCopyProjection148 replicates the pre-ticket-146 shape of
// LiveApprovals, one arm's worth of honesty: the copy that shared every
// backing.
func bareValueCopyProjection148(q *Queue) []LiveApproval {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]LiveApproval, 0, len(q.pending))
	for _, it := range q.pending {
		if it == nil || it.state != statePending {
			continue
		}
		out = append(out, LiveApproval{
			CorrelationID: it.Corr,
			Decision:      it.Dec,
			Position:      q.position(it),
		})
	}
	return out
}

// fullCloneProjection148 replicates the ticket-146 shape of LiveApprovals: one
// cloneDecision per reported row.
func fullCloneProjection148(q *Queue) []LiveApproval {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]LiveApproval, 0, len(q.pending))
	for _, it := range q.pending {
		if it == nil || it.state != statePending {
			continue
		}
		out = append(out, LiveApproval{
			CorrelationID: it.Corr,
			Decision:      cloneDecision(it.Dec),
			Position:      q.position(it),
		})
	}
	return out
}

func Benchmark148LiveApprovalsArms(b *testing.B) {
	b.Run("A1BareValueCopy_preTicket146", func(b *testing.B) {
		q := bench148Queue(b)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = bareValueCopyProjection148(q)
		}
	})
	b.Run("A2FullClone_ticket146", func(b *testing.B) {
		q := bench148Queue(b)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fullCloneProjection148(q)
		}
	})
	b.Run("A3PanelWhitelist_ticket148", func(b *testing.B) {
		q := bench148Queue(b)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = q.LiveApprovals()
		}
	})
}
