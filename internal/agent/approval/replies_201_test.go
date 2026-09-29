package approval

// White-box cases for the host-facing reply seam (replies.go, ticket 201 AC#1/
// AC#4/AC#6).
//
// These live in `package approval` because the two readings they pin are not
// reachable from outside: the rule text is derived from a card the gate mints, and
// the waiting state has to be judged with a ledger that no console is attached to.
// The end-to-end half of the same ticket (a card answered through the seam inside
// an assembled `wisp run`) is in cmd/wisp/approval_seam_201_test.go; this file is
// the part that must stay red when a rule is invented rather than read off the
// card.

import (
	"context"
	"errors"
	"testing"
	"time"
)

// card187Prompt is one displayed confirmation, shaped the way Gate.promptFor
// shapes it: Level as the verdict's own string, Paths as the canonical list the
// risk layer produced.
func card187Prompt(corr, level, tool string, paths []string, grant string) Prompt {
	return Prompt{
		CorrelationID: corr,
		TaskID:        "task-" + corr,
		Tool:          tool,
		Level:         level,
		Paths:         paths,
		Grant:         grant,
		Window:        2 * time.Second,
		Deadline:      300 * time.Second,
	}
}

// TestReplySeamWaitingStateNamesTheTwoD43Rows is AC#6's producer, pinned at the
// level the seam owns: an L2 card is AwaitingApproval, an L1 window is Confirming,
// a settled ledger is neither, and the two rows never collapse into one name.
func TestReplySeamWaitingStateNamesTheTwoD43Rows(t *testing.T) {
	r := NewReplies()
	if st, awaiting := r.WaitingState(); awaiting || st != "" {
		t.Fatalf("an empty ledger must not name a state, got (%q,%v)", st, awaiting)
	}

	r.Record(card187Prompt("l1-1", "L1", "fs.write", []string{"C:/proj/a.txt"}, ""))
	if st, awaiting := r.WaitingState(); !awaiting || string(st) != "Confirming" {
		t.Errorf("an open L1 window must read Confirming (D43 row 17), got (%q,%v)", st, awaiting)
	}

	r.Record(card187Prompt("l2-1", "L2", "fs.delete", []string{"C:/proj/b.txt"}, "nonce-1"))
	if st, awaiting := r.WaitingState(); !awaiting || string(st) != "AwaitingApproval" {
		t.Errorf("with an L2 card pending the reading must be AwaitingApproval (D43 row 17's L2 target), got (%q,%v)", st, awaiting)
	}
	if card, ok := r.AwaitingHuman(); !ok || card.CorrelationID != "l2-1" {
		t.Errorf("AwaitingHuman must prefer the C18 card that can die, got (%+v,%v)", card, ok)
	}

	// Positive control on the producer itself: retiring both entries must clear
	// the reading, or 「等人」 becomes a word this process never stops saying.
	r.Forget("l2-1")
	if st, awaiting := r.WaitingState(); !awaiting || string(st) != "Confirming" {
		t.Errorf("after the L2 card left, the L1 window must still read Confirming, got (%q,%v)", st, awaiting)
	}
	r.Forget("l1-1")
	if st, awaiting := r.WaitingState(); awaiting || st != "" {
		t.Errorf("an emptied ledger must stop naming a waiting state, got (%q,%v)", st, awaiting)
	}
}

// TestReplySeamCarriesNoAllowWithoutAGrant is AC#3's structural half, one package
// closer to the queue than the console: the seam has no way to allow that does not
// start from a grant this host was handed, and no way to allow at all before a
// Gate is bound.
func TestReplySeamCarriesNoAllowWithoutAGrant(t *testing.T) {
	ctx := context.Background()

	// Unbound: nothing can be answered, and the refusal is an error, not a yes.
	unbound := NewReplies()
	unbound.Record(card187Prompt("c-1", "L2", "fs.write", []string{"C:/proj/a.txt"}, "nonce-1"))
	if err := unbound.Allow(ctx, "c-1"); !errors.Is(err, ErrNoGateAttached) {
		t.Errorf("Allow before Attach = %v, want ErrNoGateAttached", err)
	}
	if err := unbound.Reject(ctx, "c-1", "x"); !errors.Is(err, ErrNoGateAttached) {
		t.Errorf("Reject before Attach = %v, want ErrNoGateAttached", err)
	}
	if _, err := unbound.PanelAllow(ctx, "c-1"); !errors.Is(err, ErrNoGateAttached) {
		t.Errorf("PanelAllow before Attach = %v, want ErrNoGateAttached", err)
	}

	r := NewReplies()
	r.Attach(HostBinding{Gate: nil, VetoChannel: ChannelBall})
	if err := r.Allow(ctx, "never-displayed"); !errors.Is(err, ErrNoTrackedCard) {
		t.Errorf("Allow for a card this host never displayed = %v, want ErrNoTrackedCard", err)
	}
	if err := r.Veto("never-displayed"); !errors.Is(err, ErrNoTrackedCard) {
		t.Errorf("Veto for an unknown card = %v, want ErrNoTrackedCard", err)
	}
	// An L1 window has no allow verb to answer for.
	r.Record(card187Prompt("w-1", "L1", "fs.write", []string{"C:/proj/a.txt"}, ""))
	if err := r.Allow(ctx, "w-1"); !errors.Is(err, ErrRouteHasNoAllow) {
		t.Errorf("Allow on an L1 window = %v, want ErrRouteHasNoAllow", err)
	}
	// …and the ledger entry survives that refusal, because the window is still
	// open and the operator can still veto it.
	if _, ok := r.Look("w-1"); !ok {
		t.Errorf("a refused allow must not retire the card it was asked about")
	}
}

// TestWideningRuleReadsOneLineOffTheCard is AC#4's "what rule would this store"
// text, pinned where it is computed: one directory yields the line the card shows,
// anything else yields a refusal rather than a guess.
func TestWideningRuleReadsOneLineOffTheCard(t *testing.T) {
	one := ReplyCard{Paths: []string{"C:\\Users\\me\\proj\\note.md"}}
	dir, rule, ok := one.WideningRule()
	if !ok {
		t.Fatalf("a single-directory card must offer a rule, got (%q,%q,%v)", dir, rule, ok)
	}
	if dir != "C:/Users/me/proj" {
		t.Errorf("dir = %q, want the card's own parent segment with forward slashes", dir)
	}
	if rule != `[fs] allowed_dirs += "C:/Users/me/proj"` {
		t.Errorf("rule text = %q, want the line exactly as config.toml will hold it", rule)
	}

	two := ReplyCard{Paths: []string{"C:/proj/a/x.md", "C:/proj/b/y.md"}}
	if _, _, ok := two.WideningRule(); ok {
		t.Errorf("two directories must not produce one stored rule")
	}
	dup := ReplyCard{Paths: []string{"C:/proj/a/x.md", "C:/proj/a/y.md"}}
	if _, _, ok := dup.WideningRule(); !ok {
		t.Errorf("the same directory twice is one rule, not an ambiguity")
	}
	root := ReplyCard{Paths: []string{"D:\\x.md"}}
	if _, _, ok := root.WideningRule(); ok {
		t.Errorf("a bare drive root must never be offered as an allowlist entry")
	}
	none := ReplyCard{}
	if _, _, ok := none.WideningRule(); ok {
		t.Errorf("a card that named no path can store no rule")
	}
}
