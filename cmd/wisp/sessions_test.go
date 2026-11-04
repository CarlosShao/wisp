package main

// Ticket 364-r1's end-to-end cases for the `wisp sessions` console leg.
//
// The five outlets in sessions.go were pinned at their layer below by
// internal/session (TestSessionListByTime, TestSessionListByWorkspace,
// TestSessionSearchTitleAndBody, TestSessionRenameAndArchiveUndo) and the
// no-data-root refusal by the dataroot_128 leg table. What none of those can say
// is whether THIS file's dispatch — the five `case` arms of cmdSessions and the
// one sessionsMutate all three writers share — is wired to those outlets at all.
// A commented-out case arm leaves every lower-layer case green: the store still
// renames, the console just stops reaching it. These two cases are that reading.
//
// Both run the production entry in-process (cmdSessions, the same call main.go
// makes) with a real data dir injected through sessionsIO.dataDir, so the store
// the command opens is a real SQLite wisp.db opened by memory.Open — the same
// file, the same migration chain a user's machine runs. No mock store: the
// session, memory and sqlite layers are exercised for real, and a rename that
// never reached SQLite would come back from the store unchanged.
//
// The clock is the only injected thing, and it is injected at the SEED, not in
// the command: session.ManagerOptions.Now steps one minute per call, so the rows
// the seed creates carry distinct updated_at stamps instead of two rows sharing
// one wall-clock second (the tie that internal/session's own cases document).
// The commands themselves read the real clock, exactly as a user's run does.

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/session"
)

// sessionSeed364 is one conversation to plant in the real store: a title and the
// turn bodies that make up its 正文.
type sessionSeed364 struct {
	title string
	turns []string
}

// sessionsFixture364 is a throwaway data root plus the console surface the
// command prints through. Every call resets the two buffers, so the assertions
// below always read the output of the call they just made.
type sessionsFixture364 struct {
	t   *testing.T
	dir string
	out *bytes.Buffer
	err *bytes.Buffer
}

func newSessionsFixture364(t *testing.T) *sessionsFixture364 {
	t.Helper()
	// sealableTempDir124 resolves the symlinks in the temp prefix: memory.Open
	// refuses an unresolved data root (ticket 124's discipline), and a case that
	// named one directly would die in setup on an OS that spells TMPDIR through a
	// link. "data" is a fresh child so each fixture owns its own wisp.db.
	return &sessionsFixture364{
		t:   t,
		dir: filepath.Join(sealableTempDir124(t), "data"),
		out: &bytes.Buffer{},
		err: &bytes.Buffer{},
	}
}

// call runs one `wisp sessions ...` line through the production dispatch.
func (f *sessionsFixture364) call(argv ...string) int {
	f.t.Helper()
	f.out.Reset()
	f.err.Reset()
	return cmdSessions(argv, sessionsIO{stdout: f.out, stderr: f.err, dataDir: f.dir})
}

// search runs `sessions search <query>` and fails the case if the command
// refuses; the output is what the match assertions read.
func (f *sessionsFixture364) search(query string) string {
	f.t.Helper()
	if rc := f.call("search", query); rc != 0 {
		f.t.Fatalf("364-r1 RED: `sessions search %q` exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			query, rc, f.out.String(), f.err.String())
	}
	return f.out.String()
}

// seedSessions364 plants the conversations through the production store and
// Manager — the same two calls the command makes — and returns the rows the
// Manager created, so the cases can name the ids the console has to print.
// The store is closed before the first command runs: the command has to open the
// file itself, which is the wiring under test.
func seedSessions364(t *testing.T, dir string, seeds []sessionSeed364) []memory.Session {
	t.Helper()
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("seed: memory.Open(%s): %v", dir, err)
	}
	defer func() { _ = store.Close() }()

	tick := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	mgr, err := session.NewManager(session.ManagerOptions{
		Store: store,
		Now:   func() time.Time { tick = tick.Add(time.Minute); return tick },
	})
	if err != nil {
		t.Fatalf("seed: session.NewManager: %v", err)
	}
	ctx := context.Background()
	rows := make([]memory.Session, 0, len(seeds))
	for _, seed := range seeds {
		row, err := mgr.Create(ctx, seed.title, "")
		if err != nil {
			t.Fatalf("seed: Create(%q): %v", seed.title, err)
		}
		for _, body := range seed.turns {
			if _, err := mgr.AppendTurn(ctx, row.ID, "user", body); err != nil {
				t.Fatalf("seed: AppendTurn(%s, %q): %v", row.ID, body, err)
			}
		}
		rows = append(rows, row)
	}
	return rows
}

// sessionRow364 reopens the store the console wrote to and reads one row back.
// Reading through a fresh Open is deliberate: it proves the change landed on
// disk, not in a connection the command happened to keep.
func sessionRow364(t *testing.T, dir, id string) memory.Session {
	t.Helper()
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("read back: memory.Open(%s): %v", dir, err)
	}
	defer func() { _ = store.Close() }()
	row, err := store.SessionByID(context.Background(), id)
	if err != nil {
		t.Fatalf("read back: SessionByID(%s): %v", id, err)
	}
	return row
}

// TestSessionsListSearch is outlet #1 and #2 through the console, against two
// real rows: `sessions list` prints both conversations with their ids and
// titles, `sessions search <word>` matches the title of one and the turn body of
// the other, and a word that appears in neither prints the empty list.
//
// The two searches are aimed at different columns on purpose. A search that only
// read `session.title` would still pass the first one; only the body match says
// the command reached SearchSessions' turn half. The counterpart absence checks
// (the title match must not print the other row) are what keeps a widened query
// — one that quietly returns every row — from reading as a match.
func TestSessionsListSearch(t *testing.T) {
	f := newSessionsFixture364(t)
	seeded := seedSessions364(t, f.dir, []sessionSeed364{
		{title: "Pineapple roadmap", turns: []string{"we walked the budget line"}},
		{title: "Grocery run", turns: []string{"buy mango and milk"}},
	})
	roadmap, grocery := seeded[0], seeded[1]

	if rc := f.call("list"); rc != 0 {
		t.Fatalf("364-r1 RED: `sessions list` exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			rc, f.out.String(), f.err.String())
	}
	listOut := f.out.String()
	for _, want := range []string{roadmap.ID, roadmap.Title, grocery.ID, grocery.Title} {
		if !strings.Contains(listOut, want) {
			t.Errorf("364-r1 RED: `sessions list` never printed %q, so the console is not reading the row the Manager stored:\n%s",
				want, listOut)
		}
	}

	byTitle := f.search("Pineapple")
	if !strings.Contains(byTitle, roadmap.ID) || !strings.Contains(byTitle, roadmap.Title) {
		t.Errorf("364-r1 RED: search on a TITLE word did not print %q (%s):\n%s",
			roadmap.Title, roadmap.ID, byTitle)
	}
	if strings.Contains(byTitle, grocery.ID) || strings.Contains(byTitle, grocery.Title) {
		t.Errorf("364-r1 RED: search on %q also printed the conversation that has no such title or body (a widened query reading as a match):\n%s",
			"Pineapple", byTitle)
	}

	byBody := f.search("mango")
	if !strings.Contains(byBody, grocery.ID) || !strings.Contains(byBody, grocery.Title) {
		t.Errorf("364-r1 RED: search on a turn-body word did not print %q (%s); Search must match the 正文, not only the title:\n%s",
			grocery.Title, grocery.ID, byBody)
	}
	if strings.Contains(byBody, roadmap.ID) || strings.Contains(byBody, roadmap.Title) {
		t.Errorf("364-r1 RED: search on %q also printed %q, whose title and body never mention it:\n%s",
			"mango", roadmap.Title, byBody)
	}

	if miss := f.search("no-such-token-364"); !strings.Contains(miss, "没有匹配的会话") {
		t.Errorf("364-r1 RED: a search that matches nothing printed neither the empty-list line nor a refusal:\n%s",
			miss)
	}
}

// TestSessionsRenameArchiveUndo is outlets #3 and #4 through the console, on two
// real rows so "hidden" can be told apart from "the list broke": rename rewrites
// the title and the list shows it; archive removes that one row from list while
// the bystander stays; unarchive brings it back, and the updated_at the row had
// before the archive is still the one it has after the undo.
//
// The stamp assertion is ticket 166's "归档带撤销" read through the console: the
// DAO's SetSessionArchived writes archived_at and nothing else, so an undo that
// reordered the list would show up here as a changed updated_at, not as a
// cosmetic difference.
func TestSessionsRenameArchiveUndo(t *testing.T) {
	f := newSessionsFixture364(t)
	seeded := seedSessions364(t, f.dir, []sessionSeed364{
		{title: "old title"},
		{title: "bystander conversation"},
	})
	target, bystander := seeded[0], seeded[1]

	// Premise: both are listed before anything is changed. Without this, "archive
	// hid it" and "the store was empty all along" would print the same way.
	if rc := f.call("list"); rc != 0 {
		t.Fatalf("364-r1 RED: premise: `sessions list` exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			rc, f.out.String(), f.err.String())
	}
	if premise := f.out.String(); !strings.Contains(premise, target.ID) || !strings.Contains(premise, bystander.ID) {
		t.Fatalf("364-r1 RED: premise broke: the two seeded rows are not both listed, so nothing below can attribute a change to a command:\n%s",
			premise)
	}

	// Rename.
	if rc := f.call("rename", target.ID, "new title"); rc != 0 {
		t.Fatalf("364-r1 RED: `sessions rename %s \"new title\"` exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			target.ID, rc, f.out.String(), f.err.String())
	}
	if confirm := f.out.String(); !strings.Contains(confirm, "已改名 "+target.ID) {
		t.Errorf("364-r1 RED: rename printed no confirmation naming %s:\n%s", target.ID, confirm)
	}
	renamed := sessionRow364(t, f.dir, target.ID)
	if renamed.Title != "new title" {
		t.Errorf("364-r1 RED: after rename the stored title is %q, want %q", renamed.Title, "new title")
	}
	if renamed.CreatedAt != target.CreatedAt {
		t.Errorf("364-r1 RED: rename rewrote created_at (%d -> %d); only the title column may change",
			target.CreatedAt, renamed.CreatedAt)
	}
	if rc := f.call("list"); rc != 0 {
		t.Fatalf("364-r1 RED: `sessions list` after rename exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			rc, f.out.String(), f.err.String())
	}
	renamedList := f.out.String()
	if !strings.Contains(renamedList, "new title") {
		t.Errorf("364-r1 RED: the list after rename does not show the new title:\n%s", renamedList)
	}
	if strings.Contains(renamedList, "old title") {
		t.Errorf("364-r1 RED: the list after rename still shows the old title:\n%s", renamedList)
	}
	if !strings.Contains(renamedList, bystander.ID) {
		t.Errorf("364-r1 RED: rename removed the other conversation from the list:\n%s", renamedList)
	}

	// The stamp the archive and its undo must leave alone: whatever updated_at the
	// row carries after the rename (renaming IS a touch, so it was just stamped).
	stamp := renamed.UpdatedAt

	// Archive.
	if rc := f.call("archive", target.ID); rc != 0 {
		t.Fatalf("364-r1 RED: `sessions archive %s` exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			target.ID, rc, f.out.String(), f.err.String())
	}
	if confirm := f.out.String(); !strings.Contains(confirm, "已归档 "+target.ID) {
		t.Errorf("364-r1 RED: archive printed no confirmation naming %s:\n%s", target.ID, confirm)
	}
	if rc := f.call("list"); rc != 0 {
		t.Fatalf("364-r1 RED: `sessions list` after archive exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			rc, f.out.String(), f.err.String())
	}
	archivedList := f.out.String()
	if strings.Contains(archivedList, target.ID) || strings.Contains(archivedList, "new title") {
		t.Errorf("364-r1 RED: the archived conversation is still listed:\n%s", archivedList)
	}
	if !strings.Contains(archivedList, bystander.ID) {
		t.Errorf("364-r1 RED: archiving one conversation hid the other one too:\n%s", archivedList)
	}
	archivedRow := sessionRow364(t, f.dir, target.ID)
	if archivedRow.ArchivedAt == nil {
		t.Error("364-r1 RED: the row is hidden from the list but archived_at was never written, so it is hidden by something other than the archive")
	}
	if archivedRow.UpdatedAt != stamp {
		t.Errorf("364-r1 RED: archive touched updated_at (%d -> %d); archiving is a mask, not activity",
			stamp, archivedRow.UpdatedAt)
	}

	// Unarchive.
	if rc := f.call("unarchive", target.ID); rc != 0 {
		t.Fatalf("364-r1 RED: `sessions unarchive %s` exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			target.ID, rc, f.out.String(), f.err.String())
	}
	if confirm := f.out.String(); !strings.Contains(confirm, "已取消归档 "+target.ID) {
		t.Errorf("364-r1 RED: unarchive printed no confirmation naming %s:\n%s", target.ID, confirm)
	}
	if rc := f.call("list"); rc != 0 {
		t.Fatalf("364-r1 RED: `sessions list` after unarchive exited %d, want 0\nstdout:\n%s\nstderr:\n%s",
			rc, f.out.String(), f.err.String())
	}
	restoredList := f.out.String()
	if !strings.Contains(restoredList, target.ID) || !strings.Contains(restoredList, "new title") {
		t.Errorf("364-r1 RED: the unarchived conversation did not come back to the list with its renamed title:\n%s",
			restoredList)
	}
	if !strings.Contains(restoredList, bystander.ID) {
		t.Errorf("364-r1 RED: unarchive disturbed the other conversation:\n%s", restoredList)
	}
	restored := sessionRow364(t, f.dir, target.ID)
	if restored.ArchivedAt != nil {
		t.Error("364-r1 RED: archived_at is still set after unarchive, so the undo did not take")
	}
	if restored.UpdatedAt != stamp {
		t.Errorf("364-r1 RED: the undo changed updated_at (%d -> %d), so it did not put the conversation back where it was",
			stamp, restored.UpdatedAt)
	}
}
