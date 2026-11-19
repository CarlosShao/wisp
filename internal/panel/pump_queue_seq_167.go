package panel

import "sync"

// The panel side's own display-sequence ledger (ticket 167 AC#3's stable queue
// number, re-homed here by AC#7's contract axis, 2026-11-19).
//
// WHY THIS EXISTS AND WHY IT LIVES ON THE PANEL SIDE. Ticket 167 AC#3 needs
// each approval card to carry a number that does NOT renumber when an earlier
// card is cancelled. The first leg carried it by adding a Seq field to
// internal/agent/approval's LiveApproval row - which violated ticket 167
// AC#7's contract axis verbatim: internal/agent/approval/** must stay at ZERO
// changed bytes for this ticket. The number the page shows is a display fact,
// and the display side is this package, so the mint moved here: the panel now
// casts the number itself, on first sight of a pending item, from a stable
// correlationId -> seq map. Nothing in internal/agent/approval changes - the
// queue keeps its own admission counter exactly as it was, and this ledger is
// NOT a second queue: it holds no items, grants nothing, answers nothing, and
// never decides anything. It is one map of display numbers beside the pump,
// the same class of bookkeeping StreamLog already does for stream rows.
//
// THE ONE RULE (取消不重排): the first time the mapping sees a correlation id
// it casts the next monotone number for it; every later sighting of the SAME
// id returns the SAME number - so when an earlier card is cancelled, the
// cards behind it keep their numbers (the exact proposition
// Test167QueueSeqStableAcrossCancel pins: cancel card 2, card 3 stays 3).
// A correlation id this process has never seen gets the next number, so
// numbers stay monotone in first-seen order and are never reused.
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
// *QueueSeqLedger used through Assign is a no-op that returns 0, so a pump
// assembled without one keeps the card's seq key off the wire (omitempty)
// instead of inventing a number.
func NewQueueSeqLedger() *QueueSeqLedger {
	return &QueueSeqLedger{last: map[string]uint64{}}
}

// Assign returns the stable display sequence for one correlation id: the
// number it was cast on first sight, casting the next monotone number the
// first time. An empty id casts nothing - a row with no correlation id is not
// an item the page can address, and numbering it would be a fabricated fact
// on a card; the card then travels without the seq key (the same
// reader-gated absence rule the rest of this package's keys follow).
func (l *QueueSeqLedger) Assign(correlationID string) uint64 {
	if l == nil || correlationID == "" {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if seq, ok := l.last[correlationID]; ok {
		return seq
	}
	l.next++
	l.last[correlationID] = l.next
	return l.next
}
