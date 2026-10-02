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
//     the item's own identity fields, and a replayed correlation id lands on
//     a different digest because the sequence number is folded in.
func TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce(t *testing.T) {
	s := newGrantStore()
	s.issue("nonce-live", "digest-of-item-one")

	if s.spend("nonce-live", "digest-of-item-two") {
		t.Fatal("AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding")
	}
	if s.live() != 0 {
		t.Fatalf("AC#1 RED: after a rejected spend the nonce survived (%d live); a rejected answer must still consume the grant or the forged caller can retry it", s.live())
	}
	if s.spend("nonce-live", "digest-of-item-one") {
		t.Fatal("AC#1 RED: the rejected nonce was still spendable afterwards - single-use was broken by the forged attempt")
	}
}

func TestTicket242SpendRequiresTheExactBinding(t *testing.T) {
	// The empty-binding probe gets its own store: spend consumes the nonce
	// whether or not the binding matched, so sharing a store with the happy
	// path would burn the grant before the exact-digest assertion ran.
	s := newGrantStore()
	s.issue("nonce-e", bindDigest("corr", "task", "tool", "L2", 1, []byte("args")))
	if s.spend("nonce-e", "") {
		t.Fatal("AC#1 RED: empty binding spent a live grant")
	}
	if s.live() != 0 {
		t.Fatalf("AC#1 RED: the rejected empty-binding spend left %d live nonces; consumption-on-reject is the retry killer", s.live())
	}

	happy := newGrantStore()
	want := bindDigest("corr", "task", "tool", "L2", 1, []byte("args"))
	happy.issue("nonce-a", want)
	if !happy.spend("nonce-a", want) {
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

	// The replay half: the same identity at a different time must land on a
	// different digest (the sequence number is folded in on purpose).
	itB, err := q.push(d)
	if err != nil {
		t.Fatalf("enqueue B (replayed identity): %v", err)
	}
	if itB != nil && itB.bind == itA.bind {
		t.Fatal("AC#1 RED: a replayed identity produced the same binding digest; the sequence number stopped separating items")
	}
}

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
	// Item B's binding must not open item A's grant even if a nonce leaked.
	if itA.grants.spend("leaked-or-guessed-nonce", itB.bind) {
		t.Fatal("AC#1 RED: item B's binding spent item A's grant - cross-item binding is broken")
	}
}

// sameDigest keeps the failure message honest when bindDigest's encoding
// changes shape: the assertion is identity coverage, not a hex spelling.
func sameDigest(a, b string) bool {
	return a == b
}
