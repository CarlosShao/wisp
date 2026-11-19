package panel

import "testing"

// Test167QueueSeqLedgerNumberSurvivesSetShrinkage pins the ledger half of the
// 取消不重排 proposition ON ITS OWN, at the seam the 2026-11-19 167-d repair
// leg created: castAll is the pass boundary, so "an earlier card was
// cancelled" arrives as a SMALLER set than the one that cast the numbers, and
// the ledger must leave every surviving id's number where the first set cast
// it. Under the old one-call Assign shape this proposition was fused into the
// pump's traversal order and was constructively green (the morning audit's
// finding); here it is exercisable directly, and the repair leg's mutation
// (castAll renumbering 1..n by the current order on every call, injected via
// GOFLAGS -overlay) reds this case at the card-3 assertion below.
func Test167QueueSeqLedgerNumberSurvivesSetShrinkage(t *testing.T) {
	l := NewQueueSeqLedger()

	// First presentation: three cards, cast in first-seen order 1/2/3.
	l.castAll([]string{"corr-1", "corr-2", "corr-3"})
	if got := l.number("corr-1"); got != 1 {
		t.Fatalf("corr-1 first cast = %d, want 1", got)
	}
	if got := l.number("corr-3"); got != 3 {
		t.Fatalf("corr-3 first cast = %d, want 3", got)
	}

	// THE PROPOSITION: the middle card is cancelled, so the next presentation
	// is a SMALLER set missing corr-2. Both survivors keep the numbers the
	// first set cast - corr-3 stays 3 even though it now sits in slot 2.
	l.castAll([]string{"corr-1", "corr-3"})
	if got := l.number("corr-1"); got != 1 {
		t.Errorf("corr-1 after the shrink = %d, want 1 (a survivor keeps its cast number)", got)
	}
	if got := l.number("corr-3"); got != 3 {
		t.Errorf("corr-3 after the shrink = %d, want 3 (NOT renumbered to 2 by its new place - the renumber this ledger exists to refuse; under the repair leg's renumber-every-snapshot mutation this reads 2)", got)
	}

	// The freed number is not reused: a brand-new id joins BEHIND the cast
	// frontier, monotone in first-seen order.
	l.castAll([]string{"corr-4"})
	if got := l.number("corr-4"); got != 4 {
		t.Errorf("corr-4 first cast = %d, want 4 (the cancelled card's number 2 stays retired, never reused)", got)
	}

	// An empty id casts nothing and recalls nothing: a row with no correlation
	// id is not an item the page can address.
	l.castAll([]string{""})
	if got := l.number(""); got != 0 {
		t.Errorf("empty id = %d, want 0 (no fabricated number on an unaddressable row)", got)
	}

	// The nil ledger is safe through both calls and invents nothing.
	var nl *QueueSeqLedger
	nl.castAll([]string{"corr-1"})
	if got := nl.number("corr-1"); got != 0 {
		t.Errorf("nil ledger recall = %d, want 0 (a pump without a ledger keeps the seq key off the wire)", got)
	}
}
