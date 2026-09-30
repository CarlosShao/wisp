package session

// Ticket 224 tests: the minted identity and the per-session grants ledger.
//
// The store here is a REAL *memory.Store over a temp directory, not a double:
// the property under test is "a row keyed by this session cannot be read by any
// other session", and a mock that agrees to partition rows it never had to
// partition would prove nothing. (Same reasoning internal/perm/store.go:42-51
// gives for keeping ConfigManager an interface.)
//
// Case map against the ticket's AC list:
//
//	AC#1  TestTicket224MintIsRandomAndValid          (the minter, and the shape lock)
//	AC#2写 TestTicket224RecordWritesOneRowPerAnswer   (the row exists, with its four terms)
//	AC#2读 TestTicket224CoveringMatchesItsOwnSession   (the same identity hits)
//	AC#2失效/AC#3
//	      TestTicket224CoveringDoesNotCrossToANewSession (⑩'s control ②: second minted
//	                                                     id reads 0 rows)
//	AC#4  the three controls named in the two tests above plus cmd/wisp's
//	      TestTicket224ProductionSessionDoesNotSurviveRestart; a derived id makes
//	      AC#3's case red by construction, and the row-lock test says why
//
// ⚠ That cmd/wisp case is a pointer to a REAL case since ticket 224-r2
// (cmd/wisp/ticket224_assembly_test.go); before that it named a test nobody had
// written, which 224-v1 §6-a judged as this ticket's heaviest defect. Read its own
// header for what it can and cannot claim: its "restart" is a SECOND ASSEMBLY
// inside the same test process, so it proves the minted key changes per assembly -
// not per OS process, which no instrument in this repository distinguishes.
//
// AC#2's 写/读/失效 are three functions because the AC says they must be, and the
// reason is in internal/perm/ticket90_persist_test.go's header: a merged
// persistence case can only ever be satisfied by the looser half.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// openStore gives the test a real wisp.db and a real handle on it.
func openStore(t *testing.T, dir string) *memory.Store {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	s, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func newTestLedger(t *testing.T, store Store) *Ledger {
	t.Helper()
	id, err := Mint()
	if err != nil {
		t.Fatalf("Mint: %v", err)
	}
	l, err := NewLedger(LedgerOptions{ID: id, Store: store})
	if err != nil {
		t.Fatalf("NewLedger: %v", err)
	}
	return l
}

// brokenStore answers every call with an error: the fail-closed control.
type brokenStore struct{ err error }

func (b brokenStore) InsertGrant(context.Context, memory.ApprovalGrant) (int64, error) {
	return 0, b.err
}

func (b brokenStore) ListGrantsBySession(context.Context, string) ([]memory.ApprovalGrant, error) {
	return nil, b.err
}

// ---------------------------------------------------------------------------
// AC#1 - the minter
// ---------------------------------------------------------------------------

// TestTicket224MintIsRandomAndValid is AC#1's construction half: one minting
// point, and what it produces is not recomputable.
//
// The "not derived" property is stated in A435 clause 1 as a prohibition on
// deriving from the process number, the clock, the working directory or anything
// else recomputable. Two mints INSIDE one process share all of those inputs, so
// any derivation from them would collide here - which is the strongest form of
// this claim a single-process test can make. The second-boot form
// ("两次启动断两枚 id 不等") is cmd/wisp's
// TestTicket224ProductionSessionDoesNotSurviveRestart (ticket224_assembly_test.go,
// written by 224-r2; this comment pointed at that name for a case that did not yet
// exist). ⛔ Its ceiling has to travel with any quote of it: that case's second
// boot is a second assembly in the SAME test process, so neither this file's
// same-process pair nor that case distinguishes "restart" as a new OS process -
// which is why this function asserts 200 distinct mints rather than trusting
// construction.
func TestTicket224MintIsRandomAndValid(t *testing.T) {
	seen := map[ID]bool{}
	for i := 0; i < 200; i++ {
		id, err := Mint()
		if err != nil {
			t.Fatalf("Mint #%d: %v", i, err)
		}
		if !id.Valid() {
			t.Fatalf("Mint #%d produced %q, which Valid rejects: the minter and the "+
				"shape lock disagree, so the lock cannot be used to exclude test literals",
				i, id)
		}
		if seen[id] {
			t.Fatalf("Mint produced %q twice in one process: that is a derived identity "+
				"(A435 clause 1 forbids it), and a derived identity turns 「本会话内允许」 "+
				"into a permanent pass because the next boot recomputes this boot's key", id)
		}
		seen[id] = true
	}

	// The zero value is not a session: it must not validate, and a ledger built
	// with it must refuse to exist.
	if ID("").Valid() {
		t.Error("the empty identity validates: a host that never minted would still have " +
			"a keyable session")
	}
	if _, err := NewLedger(LedgerOptions{ID: "", Store: nil}); err == nil {
		t.Error("NewLedger accepted an unminted identity")
	}
}

// TestTicket224TestLiteralsAreOutsideTheMintedShape is the structural half of the
// lock ticket 224's nail precheck asks for. The pre-224 guards in cmd/wisp and
// internal/perm key their fixture rows with hand-written strings
// ("session-before-restart", "session-after-restart"); if any production mint
// could ever equal one, AC#4's control would go green while granting nothing.
//
// This is the cheap direction. The expensive direction - "a session the REAL
// assembly minted is never one of those literals" - is asserted by the production
// minting point in cmd/wisp's TestTicket224ProductionSessionDoesNotSurviveRestart
// (cmd/wisp/ticket224_assembly_test.go, written by ticket 224-r2; until then this
// comment pointed at a name no file defined). That case compares BOTH of its
// boots' minted ids against both literals, so the lock now has both directions,
// with the same ceiling as everywhere else here: two assemblies, one test process.
func TestTicket224TestLiteralsAreOutsideTheMintedShape(t *testing.T) {
	for _, lit := range []string{
		"session-before-restart",
		"session-after-restart",
		"sess_",
		"sess_deadbeef",
		strings.TrimSuffix(idPrefix, "_"),
	} {
		if ID(lit).Valid() {
			t.Errorf("the hand-written literal %q validates as a minted session: the shape "+
				"lock is no longer excluding fixtures", lit)
		}
	}
	// And a minted id is never equal to any of them, on this run or any future one.
	for i := 0; i < 50; i++ {
		id, err := Mint()
		if err != nil {
			t.Fatal(err)
		}
		for _, lit := range []string{"session-before-restart", "session-after-restart"} {
			if id.String() == lit {
				t.Fatalf("a production mint produced the test literal %q", lit)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// AC#2 写 - the answer produces a row
// ---------------------------------------------------------------------------

// TestTicket224RecordWritesOneRowPerAnswer is AC#2's write cell: answering
// 「本会话内允许」 leaves a real approval_grant row carrying the terms D45-2 names
// - (工具, 路径模式, 会话) - plus a creation time and an expiry the DAO demands.
func TestTicket224RecordWritesOneRowPerAnswer(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	l := newTestLedger(t, store)
	ctx := context.Background()

	const tool = "fs.write"
	pattern := filepath.ToSlash(filepath.Join(t.TempDir(), "notes.txt"))

	before := time.Now().Unix()
	id, err := l.Record(ctx, tool, pattern)
	if err != nil {
		t.Fatalf("Record: %v", err)
	}
	if id == 0 {
		t.Fatal("Record returned no row id: the answer wrote nothing")
	}
	after := time.Now().Unix()

	rows, err := store.ListGrantsBySession(ctx, l.SessionID().String())
	if err != nil {
		t.Fatalf("ListGrantsBySession: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("%d rows for this session, want exactly 1", len(rows))
	}
	g := rows[0]
	if g.ID != id {
		t.Errorf("row id = %d, want %d", g.ID, id)
	}
	if g.Scope != memory.GrantScopeSession {
		t.Errorf("scope = %q, want %q: SPEC-02 §3 has exactly one scope value",
			g.Scope, memory.GrantScopeSession)
	}
	if g.Tool != tool {
		t.Errorf("tool = %q, want %q", g.Tool, tool)
	}
	if g.Pattern != pattern {
		t.Errorf("pattern = %q, want %q", g.Pattern, pattern)
	}
	if g.SessionID != l.SessionID().String() {
		t.Errorf("session_id = %q, want the minted %q", g.SessionID, l.SessionID())
	}
	if g.CreatedAt < before || g.CreatedAt > after {
		t.Errorf("created_at = %d, want inside [%d,%d]", g.CreatedAt, before, after)
	}
	if g.ExpiresAt <= g.CreatedAt {
		t.Errorf("expires_at = %d <= created_at = %d: the row is born dead, so Covering "+
			"will never honour it", g.ExpiresAt, g.CreatedAt)
	}
	if g.RevokedAt != nil {
		t.Errorf("revoked_at = %v, want NULL: a fresh answer is not a revocation", *g.RevokedAt)
	}

	// Two answers are two rows, and one answer for two paths is two rows: the
	// column is a single TEXT (SPEC-02 §3), so folding a multi-path card into one
	// string would make the matcher's equality test meaningless.
	if _, err := l.Record(ctx, tool, filepath.ToSlash(filepath.Join(t.TempDir(), "b.txt"))); err != nil {
		t.Fatalf("second Record: %v", err)
	}
	rows, _ = store.ListGrantsBySession(ctx, l.SessionID().String())
	if len(rows) != 2 {
		t.Errorf("%d rows after two answers, want 2", len(rows))
	}
}

// TestTicket224RecordRefusesUnusableSubjects keeps a row from being written for
// something no later call can ever name: an empty tool, an empty pattern, or a
// path C26 could not canonicalize (those arrive carrying the 无法规范化 marker,
// bridge.go:874). A row that can never match is not inert - it is a live-looking
// entry in the audit view that grants nothing.
func TestTicket224RecordRefusesUnusableSubjects(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	l := newTestLedger(t, store)
	ctx := context.Background()

	for _, tc := range []struct {
		name, tool, pattern string
	}{
		{"no tool", "", "/x/y"},
		{"blank tool", "   ", "/x/y"},
		{"no pattern", "fs.write", ""},
		{"unresolved path", "fs.write", `C:\weird (无法规范化: symlink loop)`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := l.Record(ctx, tc.tool, tc.pattern); err == nil {
				t.Errorf("Record(%q, %q) succeeded, want a refusal", tc.tool, tc.pattern)
			}
		})
	}
	rows, err := store.ListGrants(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Errorf("%d rows survived the refusals, want 0", len(rows))
	}

	// Control: a well-formed subject on this same ledger does land, so the four
	// refusals above are not "Record always fails".
	if _, err := l.Record(ctx, "fs.write", "/x/y"); err != nil {
		t.Fatalf("control Record failed: %v", err)
	}
	if rows, _ = store.ListGrants(ctx); len(rows) != 1 {
		t.Errorf("%d rows after the control answer, want 1", len(rows))
	}
}

// ---------------------------------------------------------------------------
// AC#2 读 - the same identity hits
// ---------------------------------------------------------------------------

func TestTicket224CoveringMatchesItsOwnSession(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	l := newTestLedger(t, store)
	ctx := context.Background()
	const tool = "fs.write"
	granted := "/work/notes.txt"
	if _, err := l.Record(ctx, tool, granted); err != nil {
		t.Fatal(err)
	}

	if id, ok := l.Covering(ctx, tool, []string{granted}); !ok || id == 0 {
		t.Errorf("Covering(own tool, own path) = (%d,%v), want (row id,true): AC#2's read "+
			"cell is not happening at the ledger level", id, ok)
	}
	// Controls, each one a different way the match could be too wide.
	for _, tc := range []struct {
		name, tool string
		paths      []string
	}{
		{"another tool", "fs.delete", []string{granted}},
		{"another path", tool, []string{"/work/other.txt"}},
		{"no paths", tool, nil},
		{"blank path", tool, []string{""}},
		{"unresolved path", tool, []string{`/work/x (无法规范化: nope)`}},
		{"partial coverage of two", tool, []string{granted, "/work/other.txt"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if id, ok := l.Covering(ctx, tc.tool, tc.paths); ok {
				t.Errorf("Covering(%q, %v) = (%d,true), want false: %s would be covered "+
					"by a rule the user never clicked", tc.tool, tc.paths, id, tc.name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AC#2 失效 + AC#3 + ⑩'s control ② - a new session reads nothing
// ---------------------------------------------------------------------------

// TestTicket224CoveringDoesNotCrossToANewSession is the permanent-pass test.
//
// It is the shape ticket 224's ⑩ rewrites AC#4 around, control ②: write with one
// production-minted id, then look up with a SECOND production-minted id and
// require zero rows and zero coverage. Same-process, different-mint is exactly
// what a restart is from the ledger's point of view, because A435 clause 2 makes
// the process the whole lifetime of an identity.
//
// Why this is also AC#4's teeth: if Mint() ever became derivable (pid, clock,
// cwd), the second mint would equal the first, ListGrantsBySession would return
// the row, Covering would return true, and both this case and cmd/wisp's restart
// case go red. That is the property the ticket demands - "a derived session id
// letting a grant cross sessions must make AC#2/AC#3 red" - and it is why the
// read goes to the database instead of an in-memory set.
func TestTicket224CoveringDoesNotCrossToANewSession(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	store := openStore(t, dir)
	writer := newTestLedger(t, store)
	ctx := context.Background()
	const tool = "fs.write"
	pattern := "/work/notes.txt"

	if _, err := writer.Record(ctx, tool, pattern); err != nil {
		t.Fatal(err)
	}
	if _, ok := writer.Covering(ctx, tool, []string{pattern}); !ok {
		t.Fatal("the writing session cannot cover its own answer: the fixture is broken, " +
			"so a false reading below would be vacuously true")
	}

	// "A new process": a second minted identity over the same bytes on disk.
	nextID, err := Mint()
	if err != nil {
		t.Fatal(err)
	}
	if nextID == writer.SessionID() {
		t.Fatal("two mints produced one identity: session end cannot be process exit if " +
			"the next process recomputes this one's key")
	}
	next, err := NewLedger(LedgerOptions{ID: nextID, Store: store})
	if err != nil {
		t.Fatal(err)
	}
	rows, err := store.ListGrantsBySession(ctx, nextID.String())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Errorf("the new session can read %d grant rows left by the old one, want 0 "+
			"(PLAN.md:1642 「会话结束后授权必须失效」)", len(rows))
	}
	if _, ok := next.Covering(ctx, tool, []string{pattern}); ok {
		t.Error("AC#3/AC#4 FAILS: the new session is covered by the old session's grant. " +
			"That is the permanent免审通行证 internal/perm/store.go:19-27 names")
	}

	// The row itself is still there, keyed to the dead session, for the 30-day
	// audit window (SPEC-02 §4). Invalidating must not mean deleting.
	all, err := store.ListGrants(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || all[0].SessionID != writer.SessionID().String() {
		t.Errorf("grant rows = %+v, want the one dead-session row kept for audit", all)
	}
}

// TestTicket224ExpiredOrRevokedRowIsNotLive checks the other two of the three
// conditions ticket 224's 乙 branch names (same identity AND not revoked AND not
// expired) with the identity held constant - so a green AC#2 读 cell cannot be
// explained by "the only thing that ever vetoes is the session key".
func TestTicket224ExpiredOrRevokedRowIsNotLive(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	store := openStore(t, dir)
	ctx := context.Background()
	const tool = "fs.write"

	live := "/work/live.txt"
	expired := "/work/expired.txt"
	revoked := "/work/revoked.txt"

	ledger := newTestLedger(t, store)
	// `live` goes through the ledger (the honest write path); the two dead rows
	// are seeded straight through the DAO with their own clocks. They must NOT be
	// Record()ed first: that would leave a live row for the same pattern, and the
	// test would then be measuring "the second row lost the race" instead of
	// "liveness is checked".
	if _, err := ledger.Record(ctx, tool, live); err != nil {
		t.Fatal(err)
	}
	// Age two of the three rows out through the DAO, so the ledger is not
	// validating against its own bookkeeping.
	seed := func(pattern string, expires int64, revoked bool) {
		t.Helper()
		var rid *int64
		if revoked {
			n := time.Now().Unix()
			rid = &n
		}
		_, err := store.InsertGrant(ctx, memory.ApprovalGrant{
			Scope: memory.GrantScopeSession, Tool: tool, Pattern: pattern,
			SessionID: ledger.SessionID().String(),
			CreatedAt: time.Now().Add(-48 * time.Hour).Unix(),
			ExpiresAt: expires, RevokedAt: rid,
		})
		if err != nil {
			t.Fatalf("seed %s: %v", pattern, err)
		}
	}
	past := time.Now().Add(-time.Hour).Unix()
	future := time.Now().Add(time.Hour).Unix()
	seed(expired, past, false)
	seed(revoked, future, true)

	if _, ok := ledger.Covering(ctx, tool, []string{live}); !ok {
		t.Error("the live row is not covering: the control below cannot be trusted")
	}
	for _, p := range []string{expired, revoked} {
		if _, ok := ledger.Covering(ctx, tool, []string{p}); ok {
			t.Errorf("Covering(%q) = true, want false: an expired or revoked row is being "+
				"honoured inside its own session", p)
		}
	}
}

// ---------------------------------------------------------------------------
// fail-closed
// ---------------------------------------------------------------------------

// TestTicket224LedgerFailsClosedOnStoreErrors: a grant lookup that cannot be
// answered is a question that must still be asked.
func TestTicket224LedgerFailsClosedOnStoreErrors(t *testing.T) {
	boom := errors.New("session fixture: store unavailable")
	id, err := Mint()
	if err != nil {
		t.Fatal(err)
	}
	l, err := NewLedger(LedgerOptions{ID: id, Store: brokenStore{err: boom}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := l.Covering(context.Background(), "fs.write", []string{"/a/b"}); ok {
		t.Error("Covering reported a grant across a broken store, want false")
	}
	if _, err := l.Record(context.Background(), "fs.write", "/a/b"); err == nil {
		t.Error("Record succeeded across a broken store, want an error")
	}

	// A nil sink must not be the reason anything panics: both paths log.
	l2, err := NewLedger(LedgerOptions{ID: id, Store: brokenStore{err: boom}})
	if err != nil {
		t.Fatal(err)
	}
	l2.logf = nil
	if _, ok := l2.Covering(context.Background(), "fs.write", []string{"/a/b"}); ok {
		t.Error("Logf-less ledger reported coverage across a broken store")
	}
}

// TestTicket224NewLedgerRefusesHandWrittenIDs is A435 clause 1 enforced at the
// only door that matters: the constructor. A host cannot hand in a session string
// it computed, guessed or copied from an earlier boot's log line.
func TestTicket224NewLedgerRefusesHandWrittenIDs(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	for _, lit := range []string{"session-before-restart", "sess_short", "sess_deadbeef"} {
		if _, err := NewLedger(LedgerOptions{ID: ID(lit), Store: store}); err == nil {
			t.Errorf("NewLedger accepted the hand-written identity %q: a derived or reused "+
				"session key is exactly what makes 「本会话内允许」 survive its session", lit)
		}
	}
	if _, err := NewLedger(LedgerOptions{ID: mustMint(t), Store: nil}); err == nil {
		t.Error("NewLedger accepted a nil Store")
	}
}

func mustMint(t *testing.T) ID {
	t.Helper()
	id, err := Mint()
	if err != nil {
		t.Fatal(err)
	}
	return id
}
