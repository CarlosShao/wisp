package session

// Ticket 166 AC#2 · the assembly-reachable outlet cases (leg 166-r1, 2026-10-18).
//
// 166-v1's ruling (docs/evidence/s1/166-session-verdict.md): the four outlets
// have functions AND unit cases, yet 生产调用者＝0/0/0/0/0 — the missing half is
// an ASSEMBLY, not logic. manage.go's guard line (var _ ManageStore =
// (*memory.Store)(nil)) pins compile-time that the concrete store a production
// root holds satisfies the slice with no adapter; the cases below pin the
// runtime half with the same shape: a Manager assembled from a concrete
// *memory.Store (declared as the interface, exactly what
// session.NewManager(session.ManagerOptions{Store: mem}) would receive at
// cmd/wisp/run.go's assembly root — mem := memory.Open(...) is already held
// there as *memory.Store at run.go:468) answers each of the four outlets
// correctly through the interface, over a real wisp.db.
//
// Why these cases are the right instrument for "接线那一格":
//   - They take ManageStore (the interface), never *ManagerOptions-with-a-mock.
//     A wiring mutation — ManageStore drifts from *memory.Store so the guard
//     line stops holding — is caught by the compiler (guard) and here by the
//     interface-typed construction refusing to compile the concrete store.
//   - Each outlet gets its own case so the dispatch's "每件都要答这一发在未修码
//     上响不响" is answerable case by case: on the baseline these four names do
//     not exist (go test -run <name> ⇒ PATTERN=rc1 / RUN=0 = 响), after they run
//     green. That baseline⇄HEAD pair is the per-outlet reading, logged in
//     .scratch/wisp/probes/166/r1/logs/.
//   - AC#6's closed coin (same-second ordering, ticket face ⛔ 重开) and the
//     326 AC#1 touch are not reopened here: every case below either compares
//     membership (Search/Archive, order-insensitive assertions) or rides the
//     stepping clock like TestSessionListByTime does.
//
// What this file does NOT claim: that production has a caller. It cannot —
// cmd/wisp/** and internal/panel/** are shared compile surfaces with zero write
// legs this wave (dispatch 166r1-dispatch.txt §射程边界), and manufacturing a
// fake caller (a new exported name, a panel JSON key, a second Store) is
// exactly what the dispatch forbids. The residual jump is reported verbatim in
// .scratch/wisp/probes/166/r1/logs/06-gap-report.txt with file:line.

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// assembleManagerThroughInterface is the wiring shape: the concrete store is
// handed to NewManager through the ManageStore interface — the same conversion
// a production root performs at the call the AC waits for. *memory.Store must
// satisfy ManageStore at THIS compile site too, not only at manage.go's guard.
func assembleManagerThroughInterface(t *testing.T, store *memory.Store, clock *time.Time) *Manager {
	t.Helper()
	var ms ManageStore = store // compile-time assembly check, same edge as run.go would hit
	m, err := NewManager(ManagerOptions{Store: ms, Now: clockFunc(clock)})
	if err != nil {
		t.Fatalf("NewManager through ManageStore: %v", err)
	}
	return m
}

func clockFunc(clock *time.Time) func() time.Time {
	if clock == nil {
		return nil
	}
	tick := *clock
	return func() time.Time { tick = tick.Add(time.Minute); return tick }
}

// Outlet #1 (列): assembled through the interface, lists newest-first with the
// workspace filter and the archived mask — membership asserted; order asserted
// only inside the one workspace, on the stepping clock (AC#6 shape, ⛔ 新硬币).
func TestAC2OutletListAssembled(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := assembleManagerThroughInterface(t, store, ptrTime(time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)))
	ctx := context.Background()

	a := mustCreate(t, m, "alpha", "wk1")
	b := mustCreate(t, m, "beta", "wk2")
	c := mustCreate(t, m, "gamma", "wk1")

	got, err := m.List(ctx, memory.SessionListOptions{Workspace: "wk1"})
	if err != nil {
		t.Fatalf("List(assembled): %v", err)
	}
	if len(got) != 2 || got[0].ID != c.ID || got[1].ID != a.ID {
		t.Errorf("workspace list = %v, want [%s %s]", sessionIDs(got), c.ID, a.ID)
	}
	if all, err := m.List(ctx, memory.SessionListOptions{}); err != nil || len(all) != 3 {
		t.Errorf("full list = %v (%v), want 3 rows", sessionIDs(all), err)
	}
	// Archive masking is part of "列" semantics (dao_session.go:118): an archived
	// row disappears from the default list and reappears on demand.
	if err := m.Archive(ctx, b.ID); err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if act, err := m.List(ctx, memory.SessionListOptions{}); err != nil || len(act) != 2 {
		t.Errorf("active list after archive = %v (%v), want 2 rows", sessionIDs(act), err)
	}
	if incl, err := m.List(ctx, memory.SessionListOptions{IncludeArchived: true}); err != nil || len(incl) != 3 {
		t.Errorf("IncludeArchived = %v (%v), want 3 rows", sessionIDs(incl), err)
	}
}

// Outlet #2 (搜): title AND body, assembled through the interface. Assertions
// stay membership-only (166-v1 §必答 1: this outlet never was a coin).
func TestAC2OutletSearchAssembled(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := assembleManagerThroughInterface(t, store, ptrTime(time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)))
	ctx := context.Background()

	notes := mustCreate(t, m, "Pineapple notes", "")
	other := mustCreate(t, m, "Shopping list", "")
	if _, err := m.AppendTurn(ctx, notes.ID, "user", "the roadmap mentions zucchini"); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}

	byTitle, err := m.Search(ctx, "Pineapple", memory.SessionListOptions{})
	if err != nil {
		t.Fatalf("Search(title): %v", err)
	}
	if len(byTitle) != 1 || byTitle[0].ID != notes.ID {
		t.Errorf("title hit = %v, want [%s]", sessionIDs(byTitle), notes.ID)
	}
	byBody, err := m.Search(ctx, "zucchini", memory.SessionListOptions{})
	if err != nil {
		t.Fatalf("Search(body): %v", err)
	}
	if len(byBody) != 1 || byBody[0].ID != notes.ID {
		t.Errorf("body hit = %v, want [%s]", sessionIDs(byBody), notes.ID)
	}
	if miss, err := m.Search(ctx, "mango", memory.SessionListOptions{}); err != nil || len(miss) != 0 {
		t.Errorf("miss = %v (%v), want none", sessionIDs(miss), err)
	}
	// The other conversation matches nothing of this query either way: assert it
	// exists as a row (2 rows total), so "search returned 1" cannot be an
	// artefact of the body never having been stored at all.
	if all, err := m.List(ctx, memory.SessionListOptions{}); err != nil || len(all) != 2 {
		t.Errorf("fixture row count = %v (%v), want 2 (notes + %s)", sessionIDs(all), err, other.ID)
	}
}

// Outlet #3 (改名): one title column moves, the body and the timestamps of the
// original stay; rename then reads back through the SAME assembled outlet.
func TestAC2OutletRenameAssembled(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := assembleManagerThroughInterface(t, store, nil)
	ctx := context.Background()

	s := mustCreate(t, m, "before name", "wk")
	if _, err := m.AppendTurn(ctx, s.ID, "user", "body stays put"); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}
	if err := m.Rename(ctx, s.ID, "  after name  "); err != nil {
		t.Fatalf("Rename(assembled): %v", err)
	}
	row, err := store.SessionByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("SessionByID: %v", err)
	}
	if row.Title != "after name" {
		t.Errorf("title = %q, want %q", row.Title, "after name")
	}
	if row.CreatedAt != s.CreatedAt {
		t.Errorf("rename rewrote created_at: %d -> %d", s.CreatedAt, row.CreatedAt)
	}
	if turns, err := m.Turns(ctx, s.ID); err != nil || len(turns) != 1 || turns[0].Content != "body stays put" {
		t.Errorf("rename disturbed the body: %+v (%v)", turns, err)
	}
	hit, err := m.Search(ctx, "after name", memory.SessionListOptions{})
	if err != nil || len(hit) != 1 || hit[0].ID != s.ID {
		t.Errorf("renamed row not searchable: %v (%v)", sessionIDs(hit), err)
	}
}

// Outlet #4 (归档带撤销): reversible through the assembled outlet — archive
// hides, unarchive restores, and neither step deletes the row or its turns.
func TestAC2OutletArchiveUndoAssembled(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := assembleManagerThroughInterface(t, store, nil)
	ctx := context.Background()

	s := mustCreate(t, m, "archive me", "wk")
	if _, err := m.AppendTurn(ctx, s.ID, "user", "kept through archive"); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}

	if err := m.Archive(ctx, s.ID); err != nil {
		t.Fatalf("Archive(assembled): %v", err)
	}
	if act, err := m.List(ctx, memory.SessionListOptions{}); err != nil || len(act) != 0 {
		t.Errorf("archived row still listed: %v (%v)", sessionIDs(act), err)
	}
	mid, err := store.SessionByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("SessionByID after archive: %v", err)
	}
	if mid.ArchivedAt == nil {
		t.Error("archive did not stamp archived_at")
	}
	if err := m.Unarchive(ctx, s.ID); err != nil {
		t.Fatalf("Unarchive(assembled): %v", err)
	}
	back, err := store.SessionByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("SessionByID after undo: %v", err)
	}
	if back.ArchivedAt != nil {
		t.Errorf("undo did not clear archived_at: %v", back.ArchivedAt)
	}
	if act, err := m.List(ctx, memory.SessionListOptions{}); err != nil || len(act) != 1 || act[0].ID != s.ID {
		t.Errorf("undo did not restore listing: %v (%v)", sessionIDs(act), err)
	}
	if turns, err := m.Turns(ctx, s.ID); err != nil || len(turns) != 1 || turns[0].Content != "kept through archive" {
		t.Errorf("archive/undo disturbed the body: %+v (%v)", turns, err)
	}
}
