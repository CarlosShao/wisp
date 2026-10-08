package approval

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/tools"
)

// Ticket 242 AC#1 (form A): the binding layer between a grant nonce and the
// ONE pending item it was minted for had zero test assertions. These cases
// pin both ends of that layer:
//
//   - the spend end: a forged binding (right nonce, wrong digest) must be
//     rejected AND must still consume the nonce (a rejected caller cannot
//     retry);
//   - the mint end: every queued item's grant is bound to a digest covering
//     the item's own identity fields, and the sequence number is folded into
//     that digest.
//
// Wording note, corrected against ticket 242 v2 (mutation M-D3): the seq term
// IS folded in, but no queue-level case can notice its removal, because
// Queue.push hands every push a fresh correlation id (queue.go:163-165) and the
// correlation id alone already separates two items of the same identity. The
// seq term is therefore measured by
// TestTicket242BindDigestSeparatesItemsBySequenceNumber, which holds every
// other field fixed and moves only seq; the queue-level case below says what it
// actually detects instead.
//
// 259-r1 mechanical note: spend used to return a bool and now returns a
// grantDenial (ticket 259 AC#2). Every judgement below is the SAME failure
// condition written against the new shape - denialNone is exactly the old true,
// anything else exactly the old false - and no message text moved. The old/new
// equivalence pair is transcribed in
// .scratch/wisp/probes/259/r1/evidence.md §2.
func TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce(t *testing.T) {
	s := newGrantStore()
	s.issue("nonce-live", "digest-of-item-one")

	if d := s.spend("nonce-live", "digest-of-item-two"); d == denialNone {
		t.Fatal("AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding")
	}
	if s.live() != 0 {
		t.Fatalf("AC#1 RED: after a rejected spend the nonce survived (%d live); a rejected answer must still consume the grant or the forged caller can retry it", s.live())
	}
	if d := s.spend("nonce-live", "digest-of-item-one"); d == denialNone {
		t.Fatal("AC#1 RED: the rejected nonce was still spendable afterwards - single-use was broken by the forged attempt")
	}
}

func TestTicket242SpendRequiresTheExactBinding(t *testing.T) {
	// The empty-binding probe gets its own store: spend consumes the nonce
	// whether or not the binding matched, so sharing a store with the happy
	// path would burn the grant before the exact-digest assertion ran.
	s := newGrantStore()
	s.issue("nonce-e", bindDigest("corr", "task", "tool", "L2", 1, []byte("args")))
	if d := s.spend("nonce-e", ""); d == denialNone {
		t.Fatal("AC#1 RED: empty binding spent a live grant")
	}
	if s.live() != 0 {
		t.Fatalf("AC#1 RED: the rejected empty-binding spend left %d live nonces; consumption-on-reject is the retry killer", s.live())
	}

	happy := newGrantStore()
	want := bindDigest("corr", "task", "tool", "L2", 1, []byte("args"))
	happy.issue("nonce-a", want)
	if d := happy.spend("nonce-a", want); d != denialNone {
		t.Fatal("AC#1 RED: the exact minted digest failed to spend its own grant - the mint/spend round trip is broken")
	}
	if happy.live() != 0 {
		t.Fatalf("AC#1 RED: a successful spend left %d live nonces; spending deletes", happy.live())
	}
}

func TestTicket242QueuedItemsBindGrantsToTheirOwnDigest(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	d := tools.Decision{
		Tool: "shell.run", TaskID: "task-1", Args: []byte(`{"cmd":"a"}`),
		Level: 2, Paths: []string{"C:/tmp/a"},
	}
	itA, err := q.push(d)
	if err != nil {
		t.Fatalf("enqueue A: %v", err)
	}
	if itA == nil || itA.bind == "" {
		t.Fatal("AC#1 RED: the queued item's grant carries an empty binding - the mint side stopped binding")
	}
	want := bindDigest(itA.Corr, d.TaskID, d.Tool, d.LevelString(), itA.Seq, d.Args)
	if itA.bind != want {
		t.Fatalf("AC#1 RED: the queued item's binding does not cover its own identity fields (seq/task/tool/level/args)")
	}

	// The replay half, stated honestly: pushing the SAME identity a second time
	// must still land on a different digest. What separates the two here is the
	// correlation id Queue assigns per push (d.CorrelationID is empty, so
	// queue.go:163-165 mints "approval-1" then "approval-2") - NOT the sequence
	// number. Ticket 242 v2's M-D3 folded seq out of bindDigest and this line
	// stayed green, which is why the seq claim now lives in its own case:
	// TestTicket242BindDigestSeparatesItemsBySequenceNumber.
	itB, err := q.push(d)
	if err != nil {
		t.Fatalf("enqueue B (replayed identity): %v", err)
	}
	if itB != nil && itB.bind == itA.bind {
		t.Fatal("AC#1 RED: a replayed identity produced the same binding digest; the digest stopped covering the fields that make the two items distinct (on this path: the correlation id the queue issued)")
	}
}

// TestTicket242ForgedBindingCannotSpendAnotherItemsGrant is AC#1's cross-card
// case, and as of this rewrite the proof it presents is really live: the nonce
// is minted for item A through Queue.grantNonce - the same call the native
// prompt makes - and A and B sit pending at the same time (asserted via
// q.Depth, plus each store's own live count).
//
// The previous version of this case passed a literal nonce that had never been
// issued into any store, so grantStore.spend walked its scan, found no row and
// returned denialSpentNonce WITHOUT ever reaching the binding comparison
// (approval.go's spend: the delete/compare pair is inside the loop body). Ticket
// 242 v2's mutation M-A emptied that comparison and the case stayed green; see
// .scratch/wisp/probes/242/v2/verdict.md Q2. That is why every refusal below
// names WHICH denial it expects instead of just "not denialNone", and why the
// live-proof precondition is asserted rather than assumed.
//
// The three shapes are told apart on purpose:
//   - A's live proof against B's digest    -> denialMisbound (AC#1's outcome)
//   - A's live proof against A's own digest -> denialNone     (positive control;
//     without it any cause at all could have produced the refusal above)
//   - a value never issued anywhere        -> denialSpentNonce (the branch the
//     toothless case actually landed on; kept asserted so this one cannot slide
//     back onto it unnoticed)
func TestTicket242ForgedBindingCannotSpendAnotherItemsGrant(t *testing.T) {
	q := NewQueue(0, 0, 0, nil)
	itA, err := q.push(tools.Decision{
		Tool: "shell.run", TaskID: "task-1", Args: []byte(`{"cmd":"a"}`),
		Level: 2, Paths: []string{"C:/tmp/a"},
	})
	if err != nil {
		t.Fatalf("enqueue A: %v", err)
	}
	itB, err := q.push(tools.Decision{
		Tool: "fs.write", TaskID: "task-2", Args: []byte(`{"path":"b"}`),
		Level: 2, Paths: []string{"C:/tmp/b"},
	})
	if err != nil {
		t.Fatalf("enqueue B: %v", err)
	}
	// Two cards alive at once, and two digests that actually differ - otherwise
	// "bound to the wrong card" is not a claim about anything.
	if q.Depth() != 2 {
		t.Fatalf("AC#1 RED: expected two simultaneously pending items, queue depth is %d", q.Depth())
	}
	if itA.bind == "" || itB.bind == "" {
		t.Fatal("AC#1 RED: one of the two pending items carries an empty binding digest, so the cross-card comparison below has nothing to mismatch")
	}
	if itA.bind == itB.bind {
		t.Fatal("AC#1 RED: two distinct pending items share one binding digest - bindDigest stopped separating cards")
	}

	nonceA, err := q.grantNonce(itA)
	if err != nil {
		t.Fatalf("mint A's grant: %v", err)
	}
	if itA.grants.live() != 1 {
		t.Fatalf("AC#1 RED: the grant minted for item A is not live on A's own store (live=%d, want 1) - the case below would then be spending a value that was never issued, which is the toothless shape ticket 242 v2 rejected", itA.grants.live())
	}
	if itB.grants.live() != 0 {
		t.Fatalf("AC#1 RED: item B's store already holds %d grant rows after A was minted - the per-item split that really refuses cross-card spends is gone", itB.grants.live())
	}

	d := itA.grants.spend(nonceA, itB.bind)
	if d != denialMisbound {
		t.Fatalf(`AC#1 RED: item A's live grant, spent against item B's binding digest, returned %s and not misbound. spent-or-never-live-nonce means the comparison never ran because the row is not in A's store (the shape this case had before, zero power over bindings); missing-nonce means nothing was presented at all; none means the cross-card grant was accepted.`, d.label())
	}
	if itA.grants.live() != 0 {
		t.Fatalf("AC#1 RED: the refused cross-card spend left %d live grant(s) on A - a rejected answer must still consume the proof, or its holder can retry it", itA.grants.live())
	}

	// Positive control: the same store, the same minting path, the item's own
	// digest - this one must open, or the refusal above proves nothing.
	nonceA2, err := q.grantNonce(itA)
	if err != nil {
		t.Fatalf("mint A's second grant: %v", err)
	}
	if d := itA.grants.spend(nonceA2, itA.bind); d != denialNone {
		t.Fatalf("AC#1 RED: the positive control was refused too (%s) - item A's own live grant against item A's own binding digest must spend, otherwise the cross-card refusal above could come from any cause", d.label())
	}
	if itA.grants.live() != 0 {
		t.Fatalf("AC#1 RED: a successful spend left %d live nonces; spending deletes", itA.grants.live())
	}

	// The discriminator: keep this case from ever becoming the empty-store shape
	// again by pinning what an unissued value answers.
	if d := itA.grants.spend("never-issued-value", itA.bind); d != denialSpentNonce {
		t.Fatalf(`AC#1 RED: a value that was never issued on A returned %s instead of spent-or-never-live-nonce - this case can no longer tell "bound to the wrong card" from "no such proof", so the denial it names is not evidence`, d.label())
	}
}

// TestTicket242BindDigestSeparatesItemsBySequenceNumber is the ruler ticket 242
// v2 found missing: its M-D3 replaced strconv.FormatUint(seq, 10) inside
// bindDigest with a constant and the whole package stayed green (verdict.md Q2
// note 4), because every queue-level pair differs by correlation id first. Here
// corr / task / tool / level / args are held equal and only seq moves, so seq is
// the sole difference-maker and removing it turns this case red.
func TestTicket242BindDigestSeparatesItemsBySequenceNumber(t *testing.T) {
	args := []byte(`{"cmd":"a"}`)
	one := bindDigest("approval-7", "task-1", "shell.run", "L2", 1, args)
	two := bindDigest("approval-7", "task-1", "shell.run", "L2", 2, args)
	if one == two {
		t.Fatal("AC#1 RED: two bind digests identical in every field except the sequence number came out equal - seq is not folded into bindDigest, so a replayed correlation id can land on the earlier card's digest")
	}
}
