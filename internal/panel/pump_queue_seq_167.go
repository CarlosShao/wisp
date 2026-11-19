package panel

import "sync"

// The panel side's own display-sequence ledger (ticket 167 AC#3's stable queue
// number, re-homed here by AC#7's contract axis, 2026-11-19; reshaped the same
// day by the audit-repair leg "167-d").
//
// TWO NUMBERING SYSTEMS EXIST ON HEAD, AND THEY ARE NOT THE SAME NUMBER (the
// 2026-11-19 audit leg made this explicit; this paragraph is the semantic
// note the dispatch ordered). The QUEUE's own counter lives in
// internal/agent/approval/queue.go: `push` increments q.seq (queue.go:165),
// casts qitem.Seq (queue.go:174) and feeds bindDigest (queue.go:180). It is
// an ADMISSION number: it is consumed by binding/correlation plumbing inside
// the approval package and never reaches any user-visible face. THE NUMBER
// THE USER SEES is cast HERE, on the panel side: pump.go's Snapshot calls
// castAll + number (pump.go:478/:481) and stamps card.Seq, and that is the
// value the wire carries (`seq` on ApprovalCardView). The two counters happen
// to agree on a fresh process (both start at 1 and both are first-seen
// monotone) but they are two independent counters with no code syncing them,
// so nothing here may assume they stay equal. Any change that needs a card
// number must come HERE, not from qitem.Seq:
// internal/agent/approval/** is the contract axis ticket 167 AC#7 keeps at
// zero changed bytes, and reading qitem.Seq across the boundary would couple
// the display face to an admission counter that was never meant for it.
//
// WHY THIS EXISTS AND WHY IT LIVES ON THE PANEL SIDE. Ticket 167 AC#3 needs
// each approval card to carry a number that does NOT renumber when an earlier
// card is cancelled. The first leg carried it by adding a Seq field to
// internal/agent/approval's LiveApproval row - which violated ticket 167
// AC#7's contract axis verbatim: internal/agent/approval/** must stay at ZERO
// changed bytes for this ticket. The number the page shows is a display fact,
// and the display side is this package, so the mint moved here: the panel now
// casts the number itself, on first sight of a pending item, from a stable
// correlation id -> seq map. Nothing in internal/agent/approval changes - the
// queue keeps its own admission counter exactly as it was, and this ledger is
// NOT a second queue: it holds no items, grants nothing, answers nothing, and
// never decides anything. It is one map of display numbers beside the pump,
// the same class of bookkeeping StreamLog already does for stream rows.
//
// WHY THE API IS TWO CALLS AND NOT ONE. The first version exposed a single
// Assign(correlationID) that the pump called once per pending card, in the
// card order of THAT snapshot. That shape was constructively green: the
// stability proposition (取消不重排) and the traversal were fused into one
// call, so a ledger that renumbered on every snapshot - casting 1..n by the
// current order - would have been indistinguishable from the real one through
// any pump-driven test. The 2026-11-19 morning audit named exactly that hole,
// and the repair separates the two halves: castAll hands the ledger the WHOLE
// pending set of one snapshot (the pass boundary is the call itself), and
// number recalls a card's number. The proposition is now exercisable on its
// own: two castAll calls with different sets - the second missing a cancelled
// id - must leave every surviving id's number where the first call cast it.
// Test167QueueSeqLedgerNumberSurvivesSetShrinkage (this package) pins that
// shape, and the mutation "renumber 1..n by the current order on every
// castAll" reds it with the named reading (card 3 reads 2); both readings are
// booked in docs/evidence/s1/167-tooth-repair-d.md.
//
// THE ONE RULE (取消不重排): the first set that presents a correlation id
// casts the next monotone number for it; every later sighting of the SAME id
// returns the SAME number, no matter where the id sits in the later set - so
// when an earlier card is cancelled, the cards behind it keep their numbers
// (the exact proposition the shrinkage tooth above pins: cancel card 2, card
// 3 stays 3). A correlation id this process has never seen gets the next
// number, so numbers stay monotone in first-seen order and are never reused.
//
// POSITION STAYS WHERE IT WAS. This ledger casts only the stable number. The
// card's 当前第几条 (Position, the live FIFO place) is still read verbatim off
// the queue's own record by the host's mapping (nativeVerdictsFromGate), the
// same source it always had - the ledger touches no queue state and no
// position computation.
//
// BOUND, STATED NOT ASSUMED: one entry per correlation id this pump ever
// showed, a string key plus a uint64 - tens of bytes each. The queue itself
// holds at most DefaultMaxPending (8) pending rows at once, but settled ids'
// entries stay so a re-shown id (it cannot recur today - correlation ids are
// minted per tool call) would still read its old number. A very long-lived
// resident process that shows ten thousand approvals holds well under a
// megabyte. No eviction on purpose: evicting a live item's number IS the
// renumber this ledger exists to refuse.
type QueueSeqLedger struct {
	mu   sync.Mutex
	last map[string]uint64
	next uint64
}

// NewQueueSeqLedger builds an empty ledger. The result is never nil; a nil
// *QueueSeqLedger is safe through both calls - castAll does nothing and
// number returns 0 - so a pump assembled without one keeps the card's seq key
// off the wire (omitempty) instead of inventing a number.
func NewQueueSeqLedger() *QueueSeqLedger {
	return &QueueSeqLedger{last: map[string]uint64{}}
}

// castAll hands the ledger the correlation ids of ONE snapshot's pending
// section, in that section's own order. Ids the ledger has never seen are
// cast the next monotone number, in the order this set presents them; ids it
// has seen keep the number they were cast, no matter where they sit in THIS
// set. That last clause is the whole proposition: presenting a set with an
// earlier card cancelled must not move any surviving id's number. An empty id
// casts nothing - a row with no correlation id is not an item the page can
// address, and numbering it would be a fabricated fact on a card. A nil
// ledger does nothing.
func (l *QueueSeqLedger) castAll(ids []string) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := l.last[id]; ok {
			continue
		}
		l.next++
		l.last[id] = l.next
	}
}

// number recalls the number an id was cast: the same number on every sighting,
// 0 when the id was never cast or is empty - the card then travels without the
// seq key (omitempty), the same reader-gated absence rule the rest of this
// package's keys follow. A nil ledger recalls nothing.
func (l *QueueSeqLedger) number(correlationID string) uint64 {
	if l == nil || correlationID == "" {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.last[correlationID]
}
