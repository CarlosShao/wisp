package memory

// Ticket 326 — adjudication instrument added by the NON-implementer leg 326-v1.
//
// Why this file exists (write-face 4 of the 326-v1 dispatch): AC#1 does not only
// require "session.updated_at gets refreshed", it requires the refresh to live in
// THE SAME TRANSACTION as the turn INSERT. Nothing in the delivered suite catches
// the opposite shape deterministically: the three AC#1 cases read only the FINAL
// state, so "INSERT committed, then a separate UPDATE" passes them, and the real
// kill bench (dao_session_append_326_test.go) only catches that shape by racing the
// kill against the second transaction (326-r1 measured 18/32 kill points). This
// file turns that probabilistic reading into a deterministic one, by making the
// refresh fail on purpose.
//
// The shape it pins is the one ticket 326 forbids in its own landing criteria
// ("写失败要 propagate，⛔ 静默降级成'只插子表'"): if the two statements share one
// transaction, a failing refresh MUST take the turn row with it. A refresh moved
// into its own transaction after the turn commit leaves the child row on disk while
// still returning an error — that is exactly a half state, and here it is caught
// without killing a process and without waiting on a wall clock.
//
// Injection seam: a SQLite BEFORE UPDATE trigger that RAISE(ABORT) on the session
// row, installed through the store's own writer handle and dropped again in
// cleanup. No production code is touched, no assertion elsewhere is loosened, no
// real SQL is mocked — the same driver, the same writeQueue, the same transaction.

import (
	"context"
	"testing"
	"time"
)

const (
	trig326v1Install = `CREATE TRIGGER zz326v1_fail_session_touch
BEFORE UPDATE OF updated_at ON session
BEGIN
  SELECT RAISE(ABORT, '326v1 injected refresh failure');
END;`
	trig326v1Drop = `DROP TRIGGER zz326v1_fail_session_touch;`
)

// TestAppendSessionTurn326v1RefreshSharesTheTurnTransaction: with the refresh made
// to fail, the turn must not survive. Red here means the touch and the turn are in
// different transactions (or the refresh error is swallowed).
func TestAppendSessionTurn326v1RefreshSharesTheTurnTransaction(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	at := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC).Unix()
	if err := s.InsertSession(ctx, Session{
		ID: "s326v1", Title: "same transaction or it is a half state", CreatedAt: at, UpdatedAt: at,
	}); err != nil {
		t.Fatalf("InsertSession: %v", err)
	}
	if _, err := s.writer.ExecContext(ctx, trig326v1Install); err != nil {
		t.Fatalf("install refresh-failing trigger: %v", err)
	}
	t.Cleanup(func() { _, _ = s.writer.ExecContext(context.Background(), trig326v1Drop) })

	_, appendErr := s.AppendSessionTurn(ctx, SessionTurn{
		SessionID: "s326v1", Role: "user", Content: "the turn that must not outlive its touch", CreatedAt: at + 60,
	})
	if appendErr == nil {
		t.Fatalf("AppendSessionTurn succeeded while the session touch was forced to fail — the " +
			"refresh error is being swallowed (ticket 326 landing criterion ③)")
	}

	turns, err := s.ListSessionTurns(ctx, "s326v1")
	if err != nil {
		t.Fatalf("ListSessionTurns after the failed append: %v", err)
	}
	row, err := s.SessionByID(ctx, "s326v1")
	if err != nil {
		t.Fatalf("SessionByID after the failed append: %v", err)
	}
	t.Logf("326-v1 determinism reading: append err=%v | session_turn rows=%d | updated_at=%d (born %d)",
		appendErr, len(turns), row.UpdatedAt, at)
	if len(turns) != 0 {
		t.Errorf("the touch and the turn do not share one transaction: the refresh failed (%v) yet %d "+
			"session_turn row(s) stayed committed while session.updated_at = %d — a committed turn with an "+
			"untouched parent row is the HALF-STATE ticket 326 AC#3 kills a subprocess to look for, and the "+
			"shape AC#1's final-state cases cannot see",
			appendErr, len(turns), row.UpdatedAt)
	}
	if row.UpdatedAt != at {
		t.Errorf("a failed touch still moved session.updated_at: %d -> %d", at, row.UpdatedAt)
	}
}

// TestAppendSessionTurn326v1TriggerBlocksNothingElse is the premise control for the
// case above: once the trigger is gone the very same call succeeds and refreshes.
// Without it the first case could be green because the trigger never fired at all.
func TestAppendSessionTurn326v1TriggerBlocksNothingElse(t *testing.T) {
	ctx := context.Background()
	s := openTestStore(t)
	at := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC).Unix()
	if err := s.InsertSession(ctx, Session{ID: "s326v1b", Title: "control", CreatedAt: at, UpdatedAt: at}); err != nil {
		t.Fatalf("InsertSession: %v", err)
	}
	if _, err := s.writer.ExecContext(ctx, trig326v1Install); err != nil {
		t.Fatalf("install trigger: %v", err)
	}
	blocked := true
	if _, err := s.AppendSessionTurn(ctx, SessionTurn{SessionID: "s326v1b", Role: "user", Content: "blocked", CreatedAt: at + 30}); err == nil {
		blocked = false
	}
	if _, err := s.writer.ExecContext(ctx, trig326v1Drop); err != nil {
		t.Fatalf("drop trigger: %v", err)
	}
	if !blocked {
		t.Fatalf("premise broken: the trigger did not make the touch fail, so the determinism case above proves nothing")
	}
	if _, err := s.AppendSessionTurn(ctx, SessionTurn{SessionID: "s326v1b", Role: "user", Content: "allowed again", CreatedAt: at + 60}); err != nil {
		t.Fatalf("append after dropping the trigger: %v", err)
	}
	turns, err := s.ListSessionTurns(ctx, "s326v1b")
	if err != nil || len(turns) != 1 {
		t.Fatalf("ListSessionTurns = %+v (%v), want the one allowed turn", turns, err)
	}
	row, err := s.SessionByID(ctx, "s326v1b")
	if err != nil || row.UpdatedAt != at+60 {
		t.Errorf("after the trigger is gone the touch must refresh again: updated_at = %d (%v), want %d",
			row.UpdatedAt, err, at+60)
	}
	t.Logf("326-v1 premise control: blocked with trigger, refreshed without it (updated_at=%d)", row.UpdatedAt)
}
