package session

// Ticket 166 AC#2 · the three cases leg 166-r1 adds on 2026-10-22, on top of the
// 2026-10-18 assembly set (ticket166r1_assembly_test.go). Each one is written to
// the same rule as that file: build the Manager from a CONCRETE *memory.Store
// handed over as ManageStore (the edge a production root performs), run against
// a real wisp.db, and assert only what the outlet itself can be seen doing.
//
// Case map, per item of AC#2 ("列 / 搜 / 改名 / 归档带撤销"):
//
//	item 1 列（按时间）      TestAC2UndoPutsTheConversationBackWhereItWas (the ORDER half)
//	item 2 搜（标题与正文）  TestAC2SearchHonoursArchiveMaskAndWorkspace
//	item 3 改名             TestAC2OutletsRefuseUnknownSession (fail-loud half)
//	item 4 归档带撤销        TestAC2UndoPutsTheConversationBackWhereItWas (the undo half)
//
// Why these three and not four: the four outlets themselves already have one
// assembly case each (TestAC2Outlet*Assembled, 10-18). What had NO case at any
// level when this leg measured — the readings are in
// docs/evidence/s1/166-r1-impl.md — is the INTERLOCK between the items:
//   - nothing asserted that archive hides a conversation from SEARCH (only from
//     List), nor that search honours the workspace scope it is offered;
//   - nothing asserted that an undo is an UNDO: the archived_at column was read
//     back by tests, but the row's ordering stamp was not, so "archive then undo"
//     could (and did) put an old conversation on top of "newest first";
//   - nothing asserted that a miss on rename/archive is an error rather than a
//     silent no-op.
//
// An outlet that answers with the wrong rows, or that succeeds while moving
// nothing, is not an outlet: the ledger of the four items would read "complete"
// while the user still cannot trust what they see. That is what these pin.

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// item 4 (归档带撤销) + item 1 (按时间列): the undo must restore the row it
// archived, and "restore" includes the position the row holds in the time-ordered
// list. The failure this case is built to catch is real and was on the tree when
// it was written: SetSessionArchived wrote updated_at=now on BOTH branches, so
//   - archiving counted as activity (an archived row would sort as the freshest
//     thing the user touched), and
//   - the undo moved an old conversation to the top of "newest first", i.e. the
//     user asked to take the archive back and got a reordered list back too.
//
// The clock is the stepping one (one minute per read) for the same reason the
// two list cases use it: same-second rows tie on "ORDER BY updated_at DESC,
// id DESC" (dao_session.go:125) and fall through to a random hex id (ticket 166
// AC#6). Nothing here may lean on that coin.
func TestAC2UndoPutsTheConversationBackWhereItWas(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := assembleManagerThroughInterface(t, store, ptrTime(time.Date(2026, 10, 22, 9, 0, 0, 0, time.UTC)))
	ctx := context.Background()

	old := mustCreate(t, m, "what we settled yesterday", "wk1")
	if _, err := m.AppendTurn(ctx, old.ID, "user", "kept through the undo"); err != nil {
		t.Fatalf("AppendTurn(old): %v", err)
	}
	newer := mustCreate(t, m, "today's thread", "wk1")

	before, err := store.SessionByID(ctx, old.ID)
	if err != nil {
		t.Fatalf("SessionByID(before): %v", err)
	}
	// Premise, asserted rather than assumed: the old conversation is NOT on top.
	if got, err := m.List(ctx, memory.SessionListOptions{Workspace: "wk1"}); err != nil ||
		len(got) != 2 || got[0].ID != newer.ID || got[1].ID != old.ID {
		t.Fatalf("premise broken: list = %v (%v), want [%s %s]", sessionIDs(got), err, newer.ID, old.ID)
	}

	if err := m.Archive(ctx, old.ID); err != nil {
		t.Fatalf("Archive: %v", err)
	}
	mid, err := store.SessionByID(ctx, old.ID)
	if err != nil {
		t.Fatalf("SessionByID(after archive): %v", err)
	}
	if mid.ArchivedAt == nil {
		t.Error("archive did not stamp archived_at")
	}
	if mid.UpdatedAt != before.UpdatedAt {
		t.Errorf("archive rewrote the last-activity stamp: updated_at %d -> %d (archiving is not activity)",
			before.UpdatedAt, mid.UpdatedAt)
	}
	if got, err := m.List(ctx, memory.SessionListOptions{}); err != nil || len(got) != 1 || got[0].ID != newer.ID {
		t.Errorf("archived row still listed by default: %v (%v)", sessionIDs(got), err)
	}

	// The undo, through the outlet the panel/CLI would call.
	if err := m.Unarchive(ctx, old.ID); err != nil {
		t.Fatalf("Unarchive: %v", err)
	}
	back, err := store.SessionByID(ctx, old.ID)
	if err != nil {
		t.Fatalf("SessionByID(after undo): %v", err)
	}
	if back.ArchivedAt != nil {
		t.Errorf("undo did not clear archived_at: %v", back.ArchivedAt)
	}
	if back.UpdatedAt != before.UpdatedAt {
		t.Errorf("undo moved the row instead of putting it back: updated_at %d, want %d (the value it had before the archive)",
			back.UpdatedAt, before.UpdatedAt)
	}
	got, err := m.List(ctx, memory.SessionListOptions{Workspace: "wk1"})
	if err != nil {
		t.Fatalf("List(after undo): %v", err)
	}
	if len(got) != 2 || got[0].ID != newer.ID || got[1].ID != old.ID {
		t.Errorf("undo did not restore the position in the time-ordered list: %v, want [%s %s]",
			sessionIDs(got), newer.ID, old.ID)
	}
	// The undo is observable on the SECOND read path too: a conversation the user
	// took back must be findable again, not merely un-flagged in a column.
	found, err := m.Search(ctx, "kept through the undo", memory.SessionListOptions{})
	if err != nil || len(found) != 1 || found[0].ID != old.ID {
		t.Errorf("restored conversation is not searchable again: %v (%v)", sessionIDs(found), err)
	}
	if turns, err := m.Turns(ctx, old.ID); err != nil || len(turns) != 1 || turns[0].Content != "kept through the undo" {
		t.Errorf("archive+undo disturbed the body: %+v (%v)", turns, err)
	}
}

// item 2 (搜: title and body) crossed with item 4 (归档): three clauses the
// search outlet is offered and nothing pinned:
//   - a hit in the TITLE and a hit in the BODY both count (the two halves the AC
//     names), over two different conversations;
//   - the workspace scope narrows the result — a body hit in another workspace
//     must not leak into this one's search;
//   - an archived conversation is hidden from search by default and comes back
//     with IncludeArchived, without the undo being needed to see that the row
//     still exists.
//
// Assertions are membership-only (never order): search was never a coin, and
// ticket 166 AC#6 recorded that it must not become one by accident.
func TestAC2SearchHonoursArchiveMaskAndWorkspace(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := assembleManagerThroughInterface(t, store, nil)
	ctx := context.Background()

	mine := mustCreate(t, m, "Pineapple notes", "wk1")
	theirs := mustCreate(t, m, "Pineapple plans", "wk2")
	if _, err := m.AppendTurn(ctx, mine.ID, "user", "the roadmap mentions zucchini"); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}

	if got, err := m.Search(ctx, "Pineapple", memory.SessionListOptions{}); err != nil || len(got) != 2 {
		t.Errorf("title hit across both workspaces = %v (%v), want 2 rows", sessionIDs(got), err)
	}
	if got, err := m.Search(ctx, "Pineapple", memory.SessionListOptions{Workspace: "wk1"}); err != nil ||
		len(got) != 1 || got[0].ID != mine.ID {
		t.Errorf("workspace-scoped title hit = %v (%v), want [%s]", sessionIDs(got), err, mine.ID)
	}
	if got, err := m.Search(ctx, "zucchini", memory.SessionListOptions{Workspace: "wk2"}); err != nil ||
		len(got) != 0 {
		t.Errorf("body hit leaked across the workspace filter: %v (%v), want none", sessionIDs(got), err)
	}

	if err := m.Archive(ctx, mine.ID); err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if got, err := m.Search(ctx, "Pineapple", memory.SessionListOptions{}); err != nil ||
		len(got) != 1 || got[0].ID != theirs.ID {
		t.Errorf("archived conversation still surfaces in search: %v (%v), want [%s]", sessionIDs(got), err, theirs.ID)
	}
	if got, err := m.Search(ctx, "Pineapple", memory.SessionListOptions{IncludeArchived: true}); err != nil ||
		len(got) != 2 {
		t.Errorf("IncludeArchived search = %v (%v), want 2 rows (archive is a mask, not a delete)", sessionIDs(got), err)
	}
	if err := m.Unarchive(ctx, mine.ID); err != nil {
		t.Fatalf("Unarchive: %v", err)
	}
	if got, err := m.Search(ctx, "Pineapple", memory.SessionListOptions{}); err != nil || len(got) != 2 {
		t.Errorf("undo did not return the conversation to search: %v (%v)", sessionIDs(got), err)
	}
}

// items 3 and 4 (改名 / 归档+撤销): a miss must be an ERROR, not a silent success.
// These outlets are what a panel draws "renamed"/"archived" from; a call that
// reports success while moving no row is worse than one that fails, because the
// UI then shows a state storage never entered (166's hard constraint 3: never
// silently lose or invent history).
//
// The positive control lives in the same case: the same calls on a REAL row do
// return nil, so "all three error" cannot be an artefact of the assertions being
// satisfied by an outlet that always fails.
func TestAC2OutletsRefuseUnknownSession(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	m := assembleManagerThroughInterface(t, store, nil)
	ctx := context.Background()

	const ghost = "does-not-exist"
	if err := m.Rename(ctx, ghost, "a name for nothing"); err == nil {
		t.Error("Rename(unknown id) reported success; the title moved nowhere")
	}
	if err := m.Archive(ctx, ghost); err == nil {
		t.Error("Archive(unknown id) reported success; nothing was archived")
	}
	if err := m.Unarchive(ctx, ghost); err == nil {
		t.Error("Unarchive(unknown id) reported success; nothing was restored")
	}
	// A failed outlet must not have created anything on its way out.
	if got, err := m.List(ctx, memory.SessionListOptions{IncludeArchived: true}); err != nil || len(got) != 0 {
		t.Errorf("failed outlets left rows behind: %v (%v)", sessionIDs(got), err)
	}
	if got, err := m.Search(ctx, "nothing", memory.SessionListOptions{}); err != nil || len(got) != 0 {
		t.Errorf("search over an empty store = %v (%v), want none and no error", sessionIDs(got), err)
	}

	// Positive control: the same three calls against a row that IS there.
	real := mustCreate(t, m, "a real conversation", "")
	if err := m.Rename(ctx, real.ID, "renamed for real"); err != nil {
		t.Errorf("Rename(real id): %v", err)
	}
	if err := m.Archive(ctx, real.ID); err != nil {
		t.Errorf("Archive(real id): %v", err)
	}
	if err := m.Unarchive(ctx, real.ID); err != nil {
		t.Errorf("Unarchive(real id): %v", err)
	}
}
