package session

// 166-v1 tie-break material bench (ticket 166 AC#6 half-b, READ ONLY).
// It answers, with numbers, the two shapes the implementation legs handed to the
// orchestrator:
//   (1) AppendSessionTurn does NOT refresh session.updated_at, so "list by time"
//       orders by the last create/rename/archive/fork - never by the last thing
//       somebody said;
//   (2) two sessions touched in the same second tie on updated_at and the
//       ORDER BY falls through to "id DESC", where id is a random hex - so the
//       order a user sees is a coin.
// Nothing here changes behaviour; it only measures. Lives in the out-of-repo
// exported tree and is archived under .scratch/wisp/probes/166/v1/.

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// TestZZ166V1AppendTurnDoesNotSurfaceSession measures (1).
func TestZZ166V1AppendTurnDoesNotSurfaceSession(t *testing.T) {
	ctx := context.Background()
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := newTestManager(t, store, ptrTime(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)))

	old := mustCreate(t, m, "the one I keep talking in", "wk")
	fresh := mustCreate(t, m, "a session I created later", "wk")

	before, err := m.List(ctx, memory.SessionListOptions{Workspace: "wk"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("TIEBREAK-BASELINE list order = %v (old.updated_at=%d fresh.updated_at=%d)\n",
		sessionIDs(before), old.UpdatedAt, fresh.UpdatedAt)

	rowBefore, err := store.SessionByID(ctx, old.ID)
	if err != nil {
		t.Fatal(err)
	}
	// Ten more turns in the OLDER conversation - the thing a user calls
	// "continuing yesterday's chat".
	for i := 0; i < 10; i++ {
		if _, err := m.AppendTurn(ctx, old.ID, "user", fmt.Sprintf("turn %d", i)); err != nil {
			t.Fatal(err)
		}
	}
	rowAfter, err := store.SessionByID(ctx, old.ID)
	if err != nil {
		t.Fatal(err)
	}
	after, err := m.List(ctx, memory.SessionListOptions{Workspace: "wk"})
	if err != nil {
		t.Fatal(err)
	}
	turns, err := m.Turns(ctx, old.ID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("TIEBREAK-AFTER appends=10 turns-in-db=%d | old.updated_at %d -> %d (delta=%d) | list order = %v\n",
		len(turns), rowBefore.UpdatedAt, rowAfter.UpdatedAt, rowAfter.UpdatedAt-rowBefore.UpdatedAt, sessionIDs(after))
	if after[0].ID != fresh.ID {
		t.Errorf("premise moved: the session with 10 new turns jumped to the top (order %v)", sessionIDs(after))
	}
	fmt.Printf("TIEBREAK-VERDICT continuing-a-chat does NOT surface it: list[0] is still %q (%q), while the\n"+
		"conversation that just got 10 turns sits at list[1]\n", after[0].ID, after[0].Title)
}

// TestZZ166V1SameSecondCoin measures (2) on the production shape: real wall
// clock (nil). For N pairs of sessions created in the SAME second, each in its
// own workspace (so the list holds exactly those two rows), count how often the
// order that comes back contradicts the order they were created in. The tie is
// broken by "id DESC", and id is a random hex (session.go:96), so the expected
// answer is a coin: the older of the two shows up on top about half the time.
func TestZZ166V1SameSecondCoin(t *testing.T) {
	ctx := context.Background()
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := newTestManager(t, store, nil) // production clock: time.Now, second resolution
	const pairs = 200
	sameSecond, reversed, differentSecond := 0, 0, 0
	for i := 0; i < pairs; i++ {
		ws := fmt.Sprintf("wk-%03d", i)
		first := mustCreate(t, m, fmt.Sprintf("alpha-%03d", i), ws)
		second := mustCreate(t, m, fmt.Sprintf("beta-%03d", i), ws)
		if first.UpdatedAt != second.UpdatedAt {
			differentSecond++
			continue
		}
		sameSecond++
		got, err := m.List(ctx, memory.SessionListOptions{Workspace: ws})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Fatalf("pair %d: list returned %d rows", i, len(got))
		}
		// Creation order says "first" was created earlier; both share updated_at,
		// so which one a user sees on top is decided by the random id.
		if got[0].ID != first.ID {
			reversed++
		}
	}
	denominator := sameSecond
	if denominator == 0 {
		denominator = 1
	}
	fmt.Printf("TIEBREAK-COIN pairs=%d same-second=%d different-second=%d | order contradicts creation order in %d/%d = %.1f%%\n",
		pairs, sameSecond, differentSecond, reversed, denominator, 100*float64(reversed)/float64(denominator))
}
