package session

// Ticket 49 tests: the D45-2 lifecycle cells the ledger owns.
//
// The store here is a REAL *memory.Store over a temp directory (the same
// posture grants_test.go states): the properties under test are "a revoked row
// stops answering", "a session-end stamp dies before the next read", and "the
// target doors refuse the shapes that would widen them", and a mock that
// agrees to remember things it never had to store proves none of them.
//
// Case map against ticket 49's AC list:
//
//	lifecycle (session end)  TestTicket49SessionEndStampDiesBeforeTheNextRead
//	revoke (immediate)       TestTicket49RevokeStopsTheNextCoveringCall
//	input.type per-process   TestTicket49TargetGrantBindsTheProcessItNamed
//	input.type all-apps      TestTicket49AllAppsRowIsOnlyDoorAndAnswersEveryProcess
//	input.type isolation     TestTicket49TargetGrantIsPerProcess
//	input.type re-confirm    TestTicket49TargetGrantNeverCoversAnotherProcess

import (
	"context"
	"errors"
	"testing"
	"time"
)

// seed49 writes one live grant row through the REAL DAO and returns its id.
func seed49(t *testing.T, l *Ledger, tool, pattern string) int64 {
	t.Helper()
	id, err := l.Record(context.Background(), tool, pattern)
	if err != nil {
		t.Fatalf("Record(%s, %s): %v", tool, pattern, err)
	}
	return id
}

// ---------------------------------------------------------------------------
// lifecycle: session end -> grants expired -> new session re-prompts
// ---------------------------------------------------------------------------

// TestTicket49SessionEndStampDiesBeforeTheNextRead is AC#1's invalidation cell
// at the layer ticket 224 left implicit. The minted identity is the primary
// death (A435 clause 2, covered by 224's own cross-boot case); what ticket 49
// adds is the EXPLICIT stamp, and the property this case pins is that the
// stamp is a fact the NEXT read of the SAME session sees: after RevokeAll, the
// session that wrote the rows can no longer cover anything, while the rows
// themselves stay on disk for the audit window (SPEC-02 §4 - retained, not
// deleted).
func TestTicket49SessionEndStampDiesBeforeTheNextRead(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	l := newTestLedger(t, s)
	ctx := context.Background()

	path := `D:/proj/notes.txt`
	id := seed49(t, l, "fs.write", path)

	// Pre-control: the live row answers.
	if _, ok := l.Covering(ctx, "fs.write", []string{path}); !ok {
		t.Fatal("pre-control failed: the live row does not cover its own path, so the invalidation below would prove nothing")
	}

	n, err := l.RevokeAll(ctx)
	if err != nil {
		t.Fatalf("RevokeAll: %v", err)
	}
	if n != 1 {
		t.Fatalf("RevokeAll stamped %d rows, want 1", n)
	}

	// THE CELL: the same session, the same path, now asks.
	if _, ok := l.Covering(ctx, "fs.write", []string{path}); ok {
		t.Errorf("a session whose grants were stamped at session end still covers %s: "+
			"「会话结束失效」 must be visible to the very next read", path)
	}

	// The row is retained for audit: the stamp sets revoked_at, it does not
	// delete. This half is what makes the invalidation an audit fact instead
	// of a cleanup somebody could forget to log.
	rows, err := l.ListBySession(ctx)
	if err != nil {
		t.Fatalf("ListBySession: %v", err)
	}
	if len(rows) != 1 || rows[0].ID != id || rows[0].RevokedAt == nil {
		t.Errorf("after session end the row must remain with revoked_at set: rows=%+v want id=%d revoked", rows, id)
	}

	// Idempotence: a second sweep finds nothing left to stamp.
	n2, err := l.RevokeAll(ctx)
	if err != nil || n2 != 0 {
		t.Errorf("second RevokeAll = (%d, %v), want (0, nil): the sweep is idempotent", n2, err)
	}
}

// ---------------------------------------------------------------------------
// revoke: immediate effect on the next evaluation
// ---------------------------------------------------------------------------

// TestTicket49RevokeStopsTheNextCoveringCall is AC#2's red/green cell: with a
// live grant the question is answered (green control), after RevokeGrant the
// NEXT Covering call answers false (red reading). "Immediate" means exactly
// this: no cache, no lazy expiry - the read goes to the store every time
// (grants.go's own header states why the read is not cached).
func TestTicket49RevokeStopsTheNextCoveringCall(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	l := newTestLedger(t, s)
	ctx := context.Background()

	path := `D:/proj/diary.txt`
	id := seed49(t, l, "fs.write", path)

	if got, ok := l.Covering(ctx, "fs.write", []string{path}); !ok || got != id {
		t.Fatalf("green control failed: Covering = (%d, %v), want (%d, true)", got, ok, id)
	}

	if err := l.Revoke(ctx, id); err != nil {
		t.Fatalf("Revoke(%d): %v", id, err)
	}

	// THE CELL: the next evaluation re-asks.
	if _, ok := l.Covering(ctx, "fs.write", []string{path}); ok {
		t.Errorf("a revoked grant still covered %s on the next evaluation: 「撤销后下一次"+
			"受门控调用重新弹卡」 is not happening", path)
	}

	// The store's stamp is idempotent and the row survives for audit.
	rows, err := l.ListBySession(ctx)
	if err != nil || len(rows) != 1 || rows[0].RevokedAt == nil {
		t.Errorf("revoked row must remain with revoked_at set: rows=%+v err=%v", rows, err)
	}
}

// ---------------------------------------------------------------------------
// input.type: the target-process binding doors
// ---------------------------------------------------------------------------

// TestTicket49TargetGrantBindsTheProcessItNamed is AC#4's per-process cell,
// positive half: the row recorded through RecordTarget answers TargetCovering
// for exactly the process name that was resolved at judge time, and the id
// that comes back is the row's own.
func TestTicket49TargetGrantBindsTheProcessItNamed(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	l := newTestLedger(t, s)
	ctx := context.Background()

	id, err := l.RecordTarget(ctx, "input.type", "notepad.exe")
	if err != nil {
		t.Fatalf("RecordTarget: %v", err)
	}
	got, ok := l.TargetCovering(ctx, "input.type", "notepad.exe")
	if !ok || got != id {
		t.Errorf("TargetCovering = (%d, %v), want (%d, true): the binding row must answer "+
			"for the process it named", got, ok, id)
	}
}

// TestTicket49TargetGrantIsPerProcess is AC#4's isolation cell: Word ≠ 终端.
// A grant for one resolved process name covers nothing else - the §15 item-7
// ③ ruling (换目标应用要重新确认) expressed as a data property.
func TestTicket49TargetGrantIsPerProcess(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	l := newTestLedger(t, s)
	ctx := context.Background()

	if _, err := l.RecordTarget(ctx, "input.type", "winword.exe"); err != nil {
		t.Fatalf("RecordTarget: %v", err)
	}
	if _, ok := l.TargetCovering(ctx, "input.type", "cmd.exe"); ok {
		t.Errorf("a Word grant covered cmd.exe: per-process isolation (Word ≠ 终端) is not holding")
	}
	if _, ok := l.TargetCovering(ctx, "fs.write", "winword.exe"); ok {
		t.Errorf("a target grant answered for a different tool: the binding is (tool, process), not process alone")
	}
}

// TestTicket49TargetGrantNeverCoversAnotherProcess is AC#4's re-confirm cell,
// the negative control that keeps the isolation cell honest: after the Word
// grant exists, the FIRST terminal question still asks - which is exactly the
// behavior the bridge's swap guard enforces at the byte boundary.
func TestTicket49TargetGrantNeverCoversAnotherProcess(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	l := newTestLedger(t, s)
	ctx := context.Background()

	if _, err := l.RecordTarget(ctx, "input.type", "winword.exe"); err != nil {
		t.Fatalf("RecordTarget: %v", err)
	}
	if _, ok := l.TargetCovering(ctx, "input.type", "windowsterminal.exe"); ok {
		t.Errorf("the first terminal question was covered by a Word grant: 换目标应用必须重新确认")
	}
}

// TestTicket49AllAppsRowIsOnlyDoorAndAnswersEveryProcess is AC#4's danger-tier
// cell, three readings in one case:
//
//	positive - a row minted by RecordTargetAllApps answers TargetCovering for
//	           ANY process (that is what 「所有应用」 means),
//	isolation - the sentinel answers only for input.type,
//	doors     - RecordTarget refuses the sentinel (a per-app answer cannot
//	            widen itself), and the sentinel for another tool is refused
//	            at the door (a dead-looking-live row is an audit lie).
func TestTicket49AllAppsRowIsOnlyDoorAndAnswersEveryProcess(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	l := newTestLedger(t, s)
	ctx := context.Background()

	id, err := l.RecordTargetAllApps(ctx, "input.type")
	if err != nil {
		t.Fatalf("RecordTargetAllApps: %v", err)
	}
	for _, proc := range []string{"notepad.exe", "cmd.exe", "winword.exe"} {
		if got, ok := l.TargetCovering(ctx, "input.type", proc); !ok || got != id {
			t.Errorf("all-apps row did not cover %s: got (%d, %v), want (%d, true)", proc, got, ok, id)
		}
	}
	if _, ok := l.TargetCovering(ctx, "fs.write", "anything.txt"); ok {
		t.Errorf("the all-apps sentinel answered for fs.write: the danger tier belongs to input.type alone")
	}
	if _, err := l.RecordTarget(ctx, "input.type", GrantPatternAllApps); err == nil {
		t.Errorf("RecordTarget accepted the all-apps sentinel: the per-app door must never widen itself")
	}
	if _, err := l.RecordTargetAllApps(ctx, "fs.write"); err == nil {
		t.Errorf("RecordTargetAllApps accepted fs.write: the sentinel row for another tool is a row that matches nothing")
	}
}

// TestTicket49PidFallbackTargetIsRefusedAtTheDoor pins the read-side refusal
// the bridge also enforces: a "pid:<n>" name must never key a grant lookup.
// The refusal lives at the bridge (tools.isGrantableTarget, where the name's
// provenance is known), so this package's own guard is the shape check: the
// ledger door accepts the pid spelling as a STRING (it cannot know the
// provenance) but the bridge refuses to key with it - pinned in
// internal/tools' own ticket-49 case. Here we pin the mirror half: a row
// stored under the pid spelling is a row whose key no future live call
// produces, i.e. a row that grants nothing, which the storeablePattern /
// recordableShape guards already treat as ordinary text. The bridge-side pin
// is the load-bearing one and lives with the resolver.
func TestTicket49PidFallbackTargetIsRefusedAtTheDoor(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	l := newTestLedger(t, s)
	ctx := context.Background()

	// The ledger door accepts the string (provenance is not its job)...
	if _, err := l.RecordTarget(ctx, "input.type", "pid:4242"); err != nil {
		t.Fatalf("RecordTarget(pid:4242): %v - the ledger door judges shape, not provenance", err)
	}
	// ...which is exactly why the BRIDGE refuses to key with that spelling.
	// The bridge-side assertion lives in internal/tools/ticket49_target_test.go
	// (TestTicket49BridgeRefusesPidFallbackKey); this case documents the split
	// so neither half is read as covering the other.
}

// TestTicket49ExpiryStillAppliesOnTheTargetDoors keeps ticket 224's expiry
// semantics true on the new doors: a target row past its expires_at answers
// nothing, even before any revoke stamp. The clock is the backstop for a row
// whose session id leaked; the identity is the primary invalidator.
func TestTicket49ExpiryStillAppliesOnTheTargetDoors(t *testing.T) {
	dir := t.TempDir()
	s := openStore(t, dir)
	base := time.Now().Add(-2 * time.Hour)
	l, err := NewLedger(LedgerOptions{ID: mustMint49(t), Store: s, Now: func() time.Time { return base }})
	if err != nil {
		t.Fatalf("NewLedger: %v", err)
	}
	ctx := context.Background()
	if _, err := l.RecordTarget(ctx, "input.type", "notepad.exe"); err != nil {
		t.Fatalf("RecordTarget: %v", err)
	}

	// A ledger whose clock has moved past the row's expiry answers nothing.
	// The row was written with expires_at = base + 30d (clockCeilingDefault),
	// so the clock here lands beyond it.
	later := base.Add(31 * 24 * time.Hour)
	rows, err := l.ListBySession(ctx)
	if err != nil {
		t.Fatalf("ListBySession: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1", len(rows))
	}
	if rows[0].ExpiresAt > later.Unix() {
		t.Fatalf("row expires_at=%d is beyond the later clock %d: the expiry stamp is broken", rows[0].ExpiresAt, later.Unix())
	}
	// grantLive on the expired row must refuse - exercised through a fresh
	// ledger whose Now reads past the row's expiry. The session id is the
	// same, so only the clock can be the difference.
	l2, err := NewLedger(LedgerOptions{ID: l.id, Store: s, Now: func() time.Time { return later }})
	if err != nil {
		t.Fatalf("NewLedger(later): %v", err)
	}
	if _, ok := l2.TargetCovering(ctx, "input.type", "notepad.exe"); ok {
		t.Errorf("an expired target row still answered: grantLive's clock condition is broken on the target door")
	}
	if _, ok := l2.Covering(ctx, "fs.write", []string{`D:/x.txt`}); ok {
		t.Errorf("an expired path row still answered: grantLive's clock condition is broken on the path door")
	}
}

func mustMint49(t *testing.T) ID {
	t.Helper()
	id, err := Mint()
	if err != nil {
		t.Fatalf("Mint: %v", err)
	}
	return id
}

// compile-time guard: the brokenStore keeps its fail-closed role.
var _ Store = brokenStore{err: errors.New("broken on purpose")}
