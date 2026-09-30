package tools

// Ticket 224-r2 cell N#1, second half: the REAL tool_call.grant_id column.
//
// What was missing is not the booking logic and not the DAO. TestTicket224Granted
// CallBooksAllowSessionGrantWithItsRow (grant_test.go:445) asserts the row against
// t224Journal, a stub that keeps whatever it is handed in a slice - so the claim
// "tool_call.grant_id names the approval_grant row that paid for this call"
// travelled through a double that never had a column to fill. The 224-v1 table
// §2-③ named exactly that: 「真 memory.Store 的 tool_call.grant_id 列今天零用例」.
// internal/memory's own dao_test.go:355 round-trips the column in isolation, so
// neither half was wrong - the two halves were simply never joined, and a
// parameter that gets dropped between bridge.book and the UPDATE statement would
// have been invisible to every test in the repository.
//
// These cases join them: a real *memory.Store over a real wisp.db, a real
// approval_grant row inserted through the real DAO, the covering id handed to the
// bridge by a grant source, and then the tool_call row READ BACK OFF DISK by
// correlation id. The assertion is the join, not the pointer: the value in the
// column must be the id of a row that actually exists in approval_grant.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// openStore224r2 is the real journal these cases need. A second temp dir keeps the
// database out of the allowlist directory, so no grant pattern can accidentally
// name a path the store itself owns.
func openStore224r2(t *testing.T) *memory.Store {
	t.Helper()
	s, err := memory.Open(filepath.Join(t.TempDir(), "data"))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// seedGrant224r2 puts one live approval_grant row on disk through the real DAO and
// returns its id, so the id asserted below is a row a reader can follow.
func seedGrant224r2(t *testing.T, s *memory.Store, tool, pattern string) int64 {
	t.Helper()
	now := time.Now().Unix()
	id, err := s.InsertGrant(context.Background(), memory.ApprovalGrant{
		Scope:     memory.GrantScopeSession,
		Tool:      tool,
		Pattern:   pattern,
		SessionID: "sess_0123456789abcdef0123456789abcdef",
		CreatedAt: now,
		ExpiresAt: now + 3600,
	})
	if err != nil {
		t.Fatalf("InsertGrant: %v", err)
	}
	if id == 0 {
		t.Fatal("the DAO handed back id 0: nothing below could point at a row")
	}
	return id
}

// req224r2 is t90Req with a caller-chosen correlation id, because reading the row
// back off disk needs a key that identifies ONE dispatch.
func req224r2(corr, target string) agent.ToolRequest {
	return agent.ToolRequest{
		TaskID: "t224-r2", CorrelationID: corr, CallID: "call-" + corr,
		Name: "t224.write", Args: json.RawMessage(fmt.Sprintf(`{"path":%q}`, target)),
	}
}

// rowByCorr224r2 reads tool_call back from the real store rather than from the
// journal the bridge was handed.
func rowByCorr224r2(t *testing.T, s *memory.Store, corr string) memory.ToolCall {
	t.Helper()
	rows, err := s.ListToolCallsByTask(context.Background(), "t224-r2")
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	for _, r := range rows {
		if r.CorrelationID == corr {
			return r
		}
	}
	t.Fatalf("no tool_call row on disk carries correlation %q; rows=%+v", corr, rows)
	return memory.ToolCall{}
}

// TestTicket224GrantIDColumnRoundTripsThroughTheRealStore is the join: the column
// holds the id of the approval_grant row that covered the call, written by the
// bridge and read back by a second handle on the same database.
func TestTicket224GrantIDColumnRoundTripsThroughTheRealStore(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	store := openStore224r2(t)

	grants := newT224Grants(0)
	b, gate, tool, paths := t224Bridge(t, risk.L1, grants, store, dir)
	canon := t224Seed(t, paths, grants, target)
	realID := seedGrant224r2(t, store, "t224.write", canon)
	// The source hands back the SAME id the DAO produced: this is the shape a real
	// *session.Ledger returns, and it is the only way the join can be asserted.
	grants.grantID = realID

	if _, err := b.Execute(context.Background(), req224r2("grant-on-disk", target)); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if w, _ := gate.counts(); w != 0 {
		t.Errorf("a live grant left the L1 window on screen (%d asks)", w)
	}
	if tool.runs != 1 {
		t.Errorf("the granted tool ran %d times, want 1", tool.runs)
	}

	row := rowByCorr224r2(t, store, "grant-on-disk")
	if row.Decision != agent.DecisionAllowGrant {
		t.Errorf("tool_call.decision on disk = %q, want %q", row.Decision, agent.DecisionAllowGrant)
	}
	if row.GrantID == nil {
		t.Fatalf("tool_call.grant_id on disk is NULL for a call a live grant covered: the column "+
			"the frozen DDL has carried since SPEC-02 §3 is still never filled in practice (%+v)", row)
	}
	if *row.GrantID != realID {
		t.Errorf("tool_call.grant_id = %d, want the covering row %d", *row.GrantID, realID)
	}
	// THE JOIN: the value must name a row that exists, with the tool and the path
	// this call was judged on. An id that points at nothing is the same lie the
	// audit line would be.
	all, err := store.ListGrants(context.Background())
	if err != nil {
		t.Fatalf("ListGrants: %v", err)
	}
	var found *memory.ApprovalGrant
	for i := range all {
		if all[i].ID == *row.GrantID {
			found = &all[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("grant_id %d points at no approval_grant row (rows=%+v): the column stores a "+
			"number, not evidence", *row.GrantID, all)
	}
	if found.Tool != "t224.write" || found.Pattern != canon {
		t.Errorf("the row grant_id points at is (%q,%q), want (%q,%q)",
			found.Tool, found.Pattern, "t224.write", canon)
	}
	if found.RevokedAt != nil {
		t.Errorf("the covering row is revoked on disk (revoked_at=%v) and still authorized a call", *found.RevokedAt)
	}

	// CONTROL, same store and same boot: an uncovered path books no grant_id. Without
	// this the column above could be a constant every call gets.
	b2, gate2, _, _ := t224Bridge(t, risk.L1, nil, store, dir)
	other := filepath.Join(dir, "never-granted.txt")
	if err := os.WriteFile(other, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := b2.Execute(context.Background(), req224r2("grant-absent", other)); err != nil {
		t.Fatalf("Execute (control): %v", err)
	}
	if w, _ := gate2.counts(); w != 1 {
		t.Fatalf("the control call asked %d times, want 1: a chain that never asks would make "+
			"the NULL below vacuous", w)
	}
	ctrl := rowByCorr224r2(t, store, "grant-absent")
	if ctrl.Decision == agent.DecisionAllowGrant {
		t.Error("the control call booked allow_session_grant with no grant source: the value " +
			"would be a constant, not a reading")
	}
	if ctrl.GrantID != nil {
		t.Errorf("control row carries grant_id %d, want SQL NULL", *ctrl.GrantID)
	}
}
