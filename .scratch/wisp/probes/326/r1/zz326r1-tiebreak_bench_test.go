package session

// 326-r1 tie-break exposure bench (ticket 326 AC#2 — READ ONLY, lives in the
// out-of-repo exported tree, archived under .scratch/wisp/probes/326/r1/).
//
// Ticket 326 AC#2 puts the same-second tie OUT of scope: `ORDER BY updated_at
// DESC, id DESC` (internal/memory/dao_session.go:125/:144) is byte-for-byte
// untouched by this leg. This file does not change behaviour and does not ask to
// — it measures ONE thing the fix necessarily changes: after a turn counts as a
// touch, how often does the touch actually move the row out of the second its
// neighbour sits in, and how often do the two still tie (in which case the order
// a user sees is decided by the random hex id, ticket 166 AC#6's coin).
//
// Production clock on purpose (newTestManager(..., nil) => time.Now, second
// resolution), because that is the shape the coin lives in.

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

func TestZZ326R1TouchTiesAfterFix(t *testing.T) {
	ctx := context.Background()
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := newTestManager(t, store, nil)

	const pairs = 200
	var advanced, tied, reversedWhenTied, wrongWhenAdvanced int
	for i := 0; i < pairs; i++ {
		ws := fmt.Sprintf("wk-%03d", i)
		first := mustCreate(t, m, fmt.Sprintf("alpha-%03d", i), ws)
		second := mustCreate(t, m, fmt.Sprintf("beta-%03d", i), ws)
		if _, err := m.AppendTurn(ctx, first.ID, "user", "still talking in the older one"); err != nil {
			t.Fatalf("AppendTurn: %v", err)
		}
		row, err := store.SessionByID(ctx, first.ID)
		if err != nil {
			t.Fatalf("SessionByID: %v", err)
		}
		got, err := m.List(ctx, memory.SessionListOptions{Workspace: ws})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("pair %d: list = %v, want 2 rows", i, sessionIDs(got))
		}
		switch {
		case row.UpdatedAt > second.UpdatedAt:
			advanced++
			if got[0].ID != first.ID {
				wrongWhenAdvanced++
			}
		case row.UpdatedAt == second.UpdatedAt:
			tied++
			if got[0].ID != first.ID {
				reversedWhenTied++
			}
		default:
			t.Fatalf("pair %d: the touch moved updated_at BACKWARDS (%d < %d)", i, row.UpdatedAt, second.UpdatedAt)
		}
	}
	fmt.Printf("TIEBREAK-AFTERFIX pairs=%d | touch advanced the second=%d (order wrong in %d of those) | "+
		"still tied with the neighbour=%d (coin landed against the touched one in %d)\n",
		pairs, advanced, wrongWhenAdvanced, tied, reversedWhenTied)
	if advanced > 0 && wrongWhenAdvanced != 0 {
		t.Errorf("premise broken: %d/%d advanced rows came back in the wrong order", wrongWhenAdvanced, advanced)
	}
}

// TestZZ326R1TouchAcrossSecondIsDeterministic is the other half of the same
// question: if the touch lands in a LATER second than the neighbour's stamp, does
// the order come back right every time? Ten pairs, one deliberate sleep per pair
// so the append crosses the second boundary. (This is the shape a user is actually
// in — "I came back to yesterday's chat and said something".)
func TestZZ326R1TouchAcrossSecondIsDeterministic(t *testing.T) {
	ctx := context.Background()
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := newTestManager(t, store, nil)

	const pairs = 10
	advanced, tied, wrong := 0, 0, 0
	for i := 0; i < pairs; i++ {
		ws := fmt.Sprintf("wd-%03d", i)
		first := mustCreate(t, m, fmt.Sprintf("alpha-%03d", i), ws)
		second := mustCreate(t, m, fmt.Sprintf("beta-%03d", i), ws)
		time.Sleep(1100 * time.Millisecond)
		if _, err := m.AppendTurn(ctx, first.ID, "user", "back to the older one a second later"); err != nil {
			t.Fatalf("AppendTurn: %v", err)
		}
		row, err := store.SessionByID(ctx, first.ID)
		if err != nil {
			t.Fatalf("SessionByID: %v", err)
		}
		got, err := m.List(ctx, memory.SessionListOptions{Workspace: ws})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("pair %d: list = %v, want 2 rows", i, sessionIDs(got))
		}
		if row.UpdatedAt > second.UpdatedAt {
			advanced++
			if got[0].ID != first.ID {
				wrong++
			}
		} else {
			tied++
		}
	}
	fmt.Printf("TIEBREAK-AFTERFIX-SLOW pairs=%d | touch advanced the second=%d (order wrong in %d) | still tied=%d\n",
		pairs, advanced, wrong, tied)
	if wrong != 0 {
		t.Errorf("the touch crossed a second boundary and the list still contradicted it: %d/%d", wrong, advanced)
	}
}
