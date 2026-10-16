package session

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// Ticket 166 AC#2/AC#4: the four outlets and the fork invariant, run against a
// real wisp.db (openStore lives in grants_test.go, same package).

// newTestManager builds a Manager over a real store. A non-nil clock advances
// one minute per call, so rows get DISTINCT second-resolution stamps; nil means
// the real time.Now, which is fine when a case compares sets, but makes any case
// that compares ORDER a coin flip (same second -> tie -> random hex id).
func newTestManager(t *testing.T, store *memory.Store, clock *time.Time) *Manager {
	t.Helper()
	opts := ManagerOptions{Store: store}
	if clock != nil {
		tick := *clock
		opts.Now = func() time.Time { tick = tick.Add(time.Minute); return tick }
	}
	m, err := NewManager(opts)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func mustCreate(t *testing.T, m *Manager, title, ws string) memory.Session {
	t.Helper()
	s, err := m.Create(context.Background(), title, ws)
	if err != nil {
		t.Fatalf("Create(%q): %v", title, err)
	}
	return s
}

func sessionIDs(ss []memory.Session) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = s.ID
	}
	return out
}

// Outlet #1 (by time): newest first.
func TestSessionListByTime(t *testing.T) {
	m := newTestManager(t, openStore(t, filepath.Join(t.TempDir(), "data")), ptrTime(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)))
	a := mustCreate(t, m, "alpha", "")
	b := mustCreate(t, m, "beta", "")
	c := mustCreate(t, m, "gamma", "")
	got, err := m.List(context.Background(), memory.SessionListOptions{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	want := []string{c.ID, b.ID, a.ID}
	if len(got) != 3 {
		t.Fatalf("List = %v, want 3 rows", sessionIDs(got))
	}
	for i, w := range want {
		if got[i].ID != w {
			t.Errorf("List[%d] = %s, want %s (order %v)", i, got[i].ID, w, sessionIDs(got))
		}
	}
}

// Outlet #1 (by workspace).
// The two wk1 rows are compared BY ORDER, so this case rides the stepping clock
// (the shape TestSessionListByTime already uses): on a wall clock both rows land
// in the same second (manage.go:93 keeps a Unix SECOND), "ORDER BY updated_at
// DESC, id DESC" (dao_session.go:125) then breaks the tie on a random hex id
// (session.go:96), and the case turns into a coin flip. Ticket 166 AC#6.
func TestSessionListByWorkspace(t *testing.T) {
	m := newTestManager(t, openStore(t, filepath.Join(t.TempDir(), "data")), ptrTime(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)))
	a := mustCreate(t, m, "alpha", "wk1")
	_ = mustCreate(t, m, "beta", "wk2")
	c := mustCreate(t, m, "gamma", "wk1")
	got, err := m.List(context.Background(), memory.SessionListOptions{Workspace: "wk1"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 || got[0].ID != c.ID || got[1].ID != a.ID {
		t.Errorf("workspace filter = %v, want [%s %s]", sessionIDs(got), c.ID, a.ID)
	}
}

// Outlet #2: search over title AND body.
// Same instrument rule as the two list cases: no wall clock here, so the rows
// the store hands back are stamped in the order the case created them and a
// future order-sensitive assertion cannot quietly become a coin flip.
func TestSessionSearchTitleAndBody(t *testing.T) {
	m := newTestManager(t, openStore(t, filepath.Join(t.TempDir(), "data")), ptrTime(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)))
	ctx := context.Background()
	notes := mustCreate(t, m, "Meeting notes", "")
	other := mustCreate(t, m, "Shopping list", "")
	for _, turn := range []struct{ id, role, body string }{
		{notes.ID, "user", "we discussed the pineapple roadmap"},
		{other.ID, "user", "buy milk"},
	} {
		if _, err := m.AppendTurn(ctx, turn.id, turn.role, turn.body); err != nil {
			t.Fatalf("AppendTurn: %v", err)
		}
	}
	byTitle, err := m.Search(ctx, "Meeting", memory.SessionListOptions{})
	if err != nil {
		t.Fatalf("Search(title): %v", err)
	}
	if len(byTitle) != 1 || byTitle[0].ID != notes.ID {
		t.Errorf("title search = %v, want [%s]", sessionIDs(byTitle), notes.ID)
	}
	byBody, err := m.Search(ctx, "pineapple", memory.SessionListOptions{})
	if err != nil {
		t.Fatalf("Search(body): %v", err)
	}
	if len(byBody) != 1 || byBody[0].ID != notes.ID {
		t.Errorf("body search = %v, want [%s]", sessionIDs(byBody), notes.ID)
	}
	if miss, err := m.Search(ctx, "no-such-token", memory.SessionListOptions{}); err != nil || len(miss) != 0 {
		t.Errorf("miss search = %v (%v), want empty", sessionIDs(miss), err)
	}
	if empty, err := m.Search(ctx, "   ", memory.SessionListOptions{}); err != nil || len(empty) != 0 {
		t.Errorf("blank search = %v (%v), want empty", sessionIDs(empty), err)
	}
}

// Outlet #3 (rename) + outlet #4 (archive with undo).
func TestSessionRenameAndArchiveUndo(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := newTestManager(t, store, nil)
	ctx := context.Background()
	s := mustCreate(t, m, "old title", "")
	if _, err := m.AppendTurn(ctx, s.ID, "user", "hello"); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}

	if err := m.Rename(ctx, s.ID, "  new title  "); err != nil {
		t.Fatalf("Rename: %v", err)
	}
	row, err := store.SessionByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("SessionByID: %v", err)
	}
	if row.Title != "new title" {
		t.Errorf("title = %q, want %q (trimmed)", row.Title, "new title")
	}
	if turns, _ := m.Turns(ctx, s.ID); len(turns) != 1 || turns[0].Content != "hello" {
		t.Errorf("rename disturbed the body: %+v", turns)
	}

	if err := m.Archive(ctx, s.ID); err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if active, _ := m.List(ctx, memory.SessionListOptions{}); len(active) != 0 {
		t.Errorf("archived session still active: %v", sessionIDs(active))
	}
	all, _ := m.List(ctx, memory.SessionListOptions{IncludeArchived: true})
	if len(all) != 1 || all[0].ArchivedAt == nil {
		t.Errorf("IncludeArchived should show the archived row: %+v", all)
	}

	// The undo.
	if err := m.Unarchive(ctx, s.ID); err != nil {
		t.Fatalf("Unarchive: %v", err)
	}
	active, _ := m.List(ctx, memory.SessionListOptions{})
	if len(active) != 1 || active[0].ArchivedAt != nil {
		t.Errorf("unarchive did not restore: %+v", active)
	}
	if _, err := store.SessionByID(ctx, s.ID); err != nil {
		t.Errorf("archive/undo deleted the row: %v", err)
	}
}

// AC#4: a fork must not swallow the original. A "tidy up the parent while
// forking" mutation (moving or deleting its turns) breaks the byte-for-byte
// comparison below — that is the mutation this test is built to catch.
func TestSessionForkDoesNotSwallowOriginal(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := newTestManager(t, store, nil)
	ctx := context.Background()
	parent := mustCreate(t, m, "orig", "wk1")
	if _, err := m.AppendTurn(ctx, parent.ID, "user", "first"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.AppendTurn(ctx, parent.ID, "assistant", "second"); err != nil {
		t.Fatal(err)
	}
	before, err := m.Turns(ctx, parent.ID)
	if err != nil || len(before) != 2 {
		t.Fatalf("parent turns before fork = %+v (%v)", before, err)
	}

	fork, err := m.Fork(ctx, parent.ID)
	if err != nil {
		t.Fatalf("Fork: %v", err)
	}
	if fork.ID == parent.ID {
		t.Fatal("fork reused the parent id")
	}
	if fork.ForkedFrom != parent.ID {
		t.Errorf("forked_from = %q, want %q", fork.ForkedFrom, parent.ID)
	}

	after, err := m.Turns(ctx, parent.ID)
	if err != nil {
		t.Fatalf("Turns(parent): %v", err)
	}
	if len(after) != len(before) {
		t.Fatalf("original lost turns (%d -> %d): the fork swallowed history", len(before), len(after))
	}
	for i := range before {
		if before[i].ID != after[i].ID || before[i].Seq != after[i].Seq || before[i].Content != after[i].Content {
			t.Errorf("original turn[%d] changed: %+v -> %+v", i, before[i], after[i])
		}
	}
	if prow, err := store.SessionByID(ctx, parent.ID); err != nil || prow.Title != "orig" || prow.ForkedFrom != "" {
		t.Errorf("parent row mutated: %+v (%v)", prow, err)
	}

	forkTurns, err := m.Turns(ctx, fork.ID)
	if err != nil || len(forkTurns) != 2 {
		t.Fatalf("fork turns = %+v (%v), want 2", forkTurns, err)
	}
	for i := range before {
		if forkTurns[i].Content != before[i].Content {
			t.Errorf("fork turn[%d] content = %q, want %q", i, forkTurns[i].Content, before[i].Content)
		}
		if forkTurns[i].ID == before[i].ID {
			t.Error("fork reused the parent's turn rowid (moved, not copied)")
		}
	}
}

func ptrTime(t time.Time) *time.Time { return &t }
