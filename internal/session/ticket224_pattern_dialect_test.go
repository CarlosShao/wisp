package session

// Ticket 224-r2 cell N#5 - the 路径模式 dialect, decided, and nailed.
//
// What was missing: rows could be written in the wildcard form the column's own
// name promises (SPEC-02 §3 calls it `pattern`, and the pre-224 guards plant
// exactly `filepath.Join(dir, "*")`), while Covering compared strings with ==. The
// result was a rule that reads like a rule and matches nothing, with no
// instrument to notice. 224-v1 §6-c row 4 registered it as 「待人拍」; ticket 224's
// 续做 N#5 assigns the choice to this leg: pick one of path.Match / hand-rolled /
// regex and add the standing nail against "looks like a rule, never matches".
//
// The decision is path.Match applied only to patterns carrying `*`, argued in
// grants.go's patternCovers comment. These cases are the four consequences that
// decision is answerable for, each pinned separately so a later change cannot
// quietly move one of them:
//
//	HITS      - a wildcard row covers a concrete path (the anti-"等于没写" nail)
//	BOUNDS    - and nothing wider than direct children of that one directory
//	DIALECT   - it is path.Match's grammar, not regex's (the two disagree on
//		  observable cases, so this is checkable rather than declarative)
//	EXACTNESS - a row without `*` is never reinterpreted, which is what keeps a
//		  path a user clicked from ever covering a path they did not.

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// seedRow224r2 writes one row straight through the DAO, so the read side can be
// measured against rows the production write path would never produce - which is
// exactly the population this dialect question is about.
func seedRow224r2(t *testing.T, s Store, sessionID, tool, pattern string) int64 {
	t.Helper()
	now := time.Now().Unix()
	id, err := s.InsertGrant(context.Background(), memory.ApprovalGrant{
		Scope: memory.GrantScopeSession, Tool: tool, Pattern: pattern,
		SessionID: sessionID, CreatedAt: now, ExpiresAt: now + 3600,
	})
	if err != nil {
		t.Fatalf("InsertGrant(%q): %v", pattern, err)
	}
	return id
}

// ---------------------------------------------------------------------------
// HITS - the nail N#5 names: "写了通配行等于没写" must now be impossible
// ---------------------------------------------------------------------------

// TestTicket224WildcardRowCoversAConcretePath is that nail. The row is written in
// the same idiom the frozen guards already plant (filepath.Join(dir,"*")), and the
// assertion is that a real path under it IS covered. Under the == matcher this
// case fails with "covered=false" and an auditor looking at the table would see a
// rule that grants nothing.
func TestTicket224WildcardRowCoversAConcretePath(t *testing.T) {
	dir := filepath.ToSlash(filepath.Join(t.TempDir(), "notes"))
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	l := newTestLedger(t, store)
	ctx := context.Background()
	pattern := dir + "/*"

	id, err := l.Record(ctx, "fs.write", pattern)
	if err != nil {
		t.Fatalf("Record(%q): %v - a wildcard row that cannot be stored is the other half of "+
			"「写了通配行等于没写」", pattern, err)
	}
	child := dir + "/todo.txt"

	gotID, ok := l.Covering(ctx, "fs.write", []string{child})
	if !ok {
		t.Fatalf("Covering(%q) = false across the live row %q (id %d): the pattern dialect is "+
			"back to exact equality, so every wildcard row on disk grants nothing while the "+
			"audit view shows a rule", child, pattern, id)
	}
	if gotID != id {
		t.Errorf("Covering returned grant_id %d, want the wildcard row %d", gotID, id)
	}

	// One row may discharge several of a call's paths, because the row is a rule
	// about a directory rather than about one file.
	if _, ok := l.Covering(ctx, "fs.write", []string{child, dir + "/other.txt"}); !ok {
		t.Errorf("two children of the same granted directory are not covered by that directory's " +
			"row: the rule was read as one file, which is narrower than what was written")
	}

	// BOUNDS - and no wider than that one directory level. path.Match's `*` stops
	// at a separator, which is the whole reason this dialect was chosen over a
	// regex: the blast radius is readable off the row.
	for _, p := range []string{
		dir + "/sub/deep.txt", // `*` never crosses a separator
		dir + "2/x.txt",       // not a prefix sibling
		"/elsewhere/x.txt",    // a different tree entirely
		dir,                   // the directory itself is not one of its children
	} {
		if _, ok := l.Covering(ctx, "fs.write", []string{p}); ok {
			t.Errorf("Covering(%q) = true across the row %q: wider than the row says, and wider "+
				"than anything a user could have meant", p, pattern)
		}
	}

	// The other three terms stay exact: a wildcard pattern never turns into a
	// wildcard tool.
	if _, ok := l.Covering(ctx, "fs.delete", []string{child}); ok {
		t.Error("the wildcard row covers another tool: the dialect applies to the pattern term " +
			"only, and a tool the user never named is the widest widening available")
	}

	// Partial coverage still asks: one child covered, one outside the directory.
	if _, ok := l.Covering(ctx, "fs.write", []string{child, "/somewhere/else.txt"}); ok {
		t.Error("a call whose paths are only PARTLY covered reports covered: that would let one " +
			"remembered directory authorize a write outside it")
	}

	// And the row that hit is the row that exists: read it back through the DAO.
	rows, err := store.ListGrantsBySession(ctx, l.SessionID().String())
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Pattern != pattern || rows[0].ID != id {
		t.Errorf("ledger rows = %+v, want the one wildcard row %d", rows, id)
	}
}

// ---------------------------------------------------------------------------
// DIALECT - path.Match's grammar, checkable rather than declared
// ---------------------------------------------------------------------------

// TestTicket224DialectIsPathMatchNotRegex pins WHICH dialect, using the cases the
// two candidates disagree on. Swap patternCovers to regexp and at least one
// assertion below goes red, which is the point: "we chose path.Match" has to be a
// reading, not a comment.
func TestTicket224DialectIsPathMatchNotRegex(t *testing.T) {
	cases := []struct {
		name, pattern string
		covers        []string
		refuses       []string
	}{
		{
			// `?` is one non-separator character in path.Match; in regex it is a
			// quantifier with nothing to quantify, and a regex matcher would either
			// error or read the pattern as "wor" + optional k.
			name:    "question mark is a single character",
			pattern: "/w?rk/*",
			covers:  []string{"/work/a.txt", "/wark/b.txt"},
			refuses: []string{"/woork/a.txt", "/work/sub/a.txt", "/w/a.txt"},
		},
		{
			// `[ab]` is a byte class in path.Match and an alternation-free group in
			// regex - the observable difference is the hyphen and the escaping.
			name:    "bracket is a class over one character",
			pattern: "/work/[ab]*.txt",
			covers:  []string{"/work/a.txt", "/work/bb.txt"},
			refuses: []string{"/work/c.txt", "/work/a/deep.txt"},
		},
		{
			// The regex idiom. Under path.Match `.` is a literal dot, so a row written
			// as a regex matches nothing - and that is now a pinned consequence of the
			// choice instead of an undetected accident. ⛔ Anyone who wants regex has to
			// change this case on purpose, which makes it a contract conversation, not a
			// refactor.
			name:    "a regex-shaped row is honoured as a literal pattern",
			pattern: "/work/.*",
			covers:  []string{"/work/.hidden"},
			refuses: []string{"/work/a.txt", "/work/x"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// One ledger per case: an earlier case's wildcard row legitimately covers a
			// later case's refused path, and sharing one table would make those hits
			// look like dialect failures.
			l := newTestLedger(t, openStore(t, filepath.Join(t.TempDir(), "data")))
			ctx := context.Background()
			if _, err := l.Record(ctx, "fs.write", tc.pattern); err != nil {
				t.Fatalf("Record(%q): %v", tc.pattern, err)
			}
			for _, p := range tc.covers {
				if _, ok := l.Covering(ctx, "fs.write", []string{p}); !ok {
					t.Errorf("%q does not cover %q under this dialect", tc.pattern, p)
				}
			}
			for _, p := range tc.refuses {
				if _, ok := l.Covering(ctx, "fs.write", []string{p}); ok {
					t.Errorf("%q covers %q, which the chosen dialect does not license", tc.pattern, p)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EXACTNESS - the security half: only `*` can ever make a row wider than itself
// ---------------------------------------------------------------------------

// TestTicket224PatternWithoutStarIsNeverReinterpreted is the reason the dialect is
// gated on `*` instead of being applied to every row. A card prints a canonicalized
// CONCRETE path, and the classes path.Match would interpret are reachable in real
// file names (`[` and `]` are legal on Windows; `?` is not, but the case is cheap to
// pin anyway). Interpreting them would let a click on `note[s].txt` authorize
// `notes.txt`, which is wider than anything shown on the card.
func TestTicket224PatternWithoutStarIsNeverReinterpreted(t *testing.T) {
	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	l := newTestLedger(t, store)
	ctx := context.Background()

	for _, tc := range []struct {
		row, nearMiss, exact string
	}{
		{"/work/note[s].txt", "/work/notes.txt", "/work/note[s].txt"},
		{"/work/a?c.txt", "/work/abc.txt", "/work/a?c.txt"},
		{"/work/[a-c].txt", "/work/b.txt", "/work/[a-c].txt"},
		{"/work/plain.txt", "/work/plainxtxt", "/work/plain.txt"},
	} {
		t.Run(tc.row, func(t *testing.T) {
			if _, err := l.Record(ctx, "fs.write", tc.row); err != nil {
				t.Fatalf("Record(%q): %v", tc.row, err)
			}
			if _, ok := l.Covering(ctx, "fs.write", []string{tc.nearMiss}); ok {
				t.Errorf("the row %q covers %q: a class was interpreted in a row that never said "+
					"wildcard, so a path nobody clicked is authorized by a path somebody did",
					tc.row, tc.nearMiss)
			}
			if _, ok := l.Covering(ctx, "fs.write", []string{tc.exact}); !ok {
				t.Errorf("the row %q does not even cover its own string: exact equality is gone, "+
					"which would silently un-authorize the answer the user actually gave", tc.row)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// the write side refuses a rule that could never match
// ---------------------------------------------------------------------------

// TestTicket224MalformedWildcardIsRefusedAtTheDoorAndIgnoredOnDisk closes the
// "看着有规则其实从不匹配" hole from the other side: a `*`-bearing pattern that
// path.Match cannot compile is refused before it becomes a row, and if it is on
// disk anyway (seeded through the DAO, or written by an older build) it is honoured
// only as its own literal.
func TestTicket224MalformedWildcardIsRefusedAtTheDoorAndIgnoredOnDisk(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	store := openStore(t, dir)
	l := newTestLedger(t, store)
	ctx := context.Background()

	// The one shape the write side refuses: a pattern that carries the wildcard and
	// that path.Match cannot compile. Probe over the standard library at this anchor:
	// "/work/[a*" is ErrBadPattern, while "/work/[a-z" and "/work/]" are NOT errors
	// once the class is left unterminated/unopened - and since they carry no `*` this
	// ledger never compiles them anyway, which the block below pins on its own.
	for _, bad := range []string{"/work/[a*", "/work/[a-9*", "/work/x[ab*"} {
		if _, err := l.Record(ctx, "fs.write", bad); err == nil {
			t.Errorf("Record(%q) succeeded, want a refusal: this row would sit in the audit view "+
				"looking like a rule while matching nothing at all", bad)
		}
	}
	rows, err := store.ListGrants(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("%d malformed wildcard rows survived the refusals: %+v", len(rows), rows)
	}

	// The star-free bracket shapes stay storeable, because they stay literal: a real
	// file may be named `note[s].txt`, and refusing that row would be refusing an
	// answer the user really gave. They match exactly one string and nothing else.
	for _, literal := range []string{"/work/[a-z", "/work/]", "/work/note[s].txt"} {
		if _, err := l.Record(ctx, "fs.write", literal); err != nil {
			t.Errorf("Record(%q): %v - a star-free row is its own literal and must stay storeable",
				literal, err)
			continue
		}
		if _, ok := l.Covering(ctx, "fs.write", []string{literal}); !ok {
			t.Errorf("the literal row %q does not cover itself", literal)
		}
	}

	// Control: a well-formed wildcard on the same ledger does land, so the two
	// refusals above are not "Record always fails".
	if _, err := l.Record(ctx, "fs.write", "/work/*"); err != nil {
		t.Fatalf("control Record of a well-formed wildcard: %v", err)
	}
	if _, ok := l.Covering(ctx, "fs.write", []string{"/work/x.txt"}); !ok {
		t.Error("the control wildcard row does not cover: the refusals above cannot be attributed " +
			"to syntax, because nothing matches at all")
	}

	// Now the same broken shape through the DAO, on a ledger that holds ONLY it, so
	// the control row above cannot be the reason anything matches. Fail closed,
	// never a hit - except as its own exact string.
	ghostStore := openStore(t, filepath.Join(t.TempDir(), "ghost-data"))
	ghost := newTestLedger(t, ghostStore)
	ghostID := seedRow224r2(t, ghostStore, ghost.SessionID().String(), "fs.write", "/work/[a-z*")
	for _, p := range []string{"/work/a*", "/work/x.txt", "/work/-"} {
		if _, ok := ghost.Covering(ctx, "fs.write", []string{p}); ok {
			t.Errorf("a row path.Match cannot compile (id %d) covers %q: this is the live-looking "+
				"nothing this case exists to refuse", ghostID, p)
		}
	}
	if _, ok := ghost.Covering(ctx, "fs.write", []string{"/work/[a-z*"}); !ok {
		t.Errorf("the malformed row does not cover its own literal: the fallback is exact equality, " +
			"not silence - a caller that names the row verbatim is still answered")
	}

	// And a bad row for another session must stay invisible regardless of dialect.
	other := newTestLedger(t, ghostStore)
	if _, ok := other.Covering(ctx, "fs.write", []string{"/work/[a-z*"}); ok {
		t.Error("the malformed row crossed sessions: the session term was never in scope here, " +
			"and losing it would be the permanent免审通行证")
	}
}

// ---------------------------------------------------------------------------
// the dialect is one function's business, not several
// ---------------------------------------------------------------------------

// TestTicket224PatternCoversIsTheOnlyMatcher keeps the decision in one place: if a
// second interpretation ever appears (a normalizer here, a TrimSuffix there), this
// is the case that notices. It reads the truth table of patternCovers directly, so
// the dialect statement and the code cannot drift apart silently.
func TestTicket224PatternCoversIsTheOnlyMatcher(t *testing.T) {
	cases := []struct {
		pattern, p string
		want       bool
		why        string
	}{
		{"/work/a.txt", "/work/a.txt", true, "identical"},
		{"/work/a.txt", "/work/b.txt", false, "different concrete paths"},
		{"/work/*", "/work/a.txt", true, "direct child"},
		{"/work/*", "/work/", true, "the empty tail is still one non-separator segment"},
		{"/work/*", "/work/sub/a.txt", false, "star stops at the separator"},
		{"/work/*", "/workx/a.txt", false, "not a prefix match"},
		{"/work/a*", "/work/a.txt", true, "anchored star"},
		{"*", "/anything", false, "a bare star cannot reach a leading separator's first segment"},
		{"/work/[a*", "/work/anything", false, "ErrBadPattern fails closed"},
		{"/work/a?c.txt", "/work/abc.txt", false, "no star, so no interpretation at all"},
	}
	for _, tc := range cases {
		if got := patternCovers(tc.pattern, tc.p); got != tc.want {
			t.Errorf("patternCovers(%q, %q) = %v, want %v (%s)", tc.pattern, tc.p, got, tc.want, tc.why)
		}
	}
	// The dialect is slash-form only: a row and a path that differ by separator
	// style must not match, because the write side stores ToSlash and the bridge
	// hands canonicalized slash paths. If a future caller skips that, the pattern
	// silently stops matching - which is fail-closed, but worth a named assertion.
	if patternCovers(`/work/*`, `\work\a.txt`) {
		t.Error("a backslashed path is covered by a slashed pattern: the two spellings are different " +
			"sessions to this matcher, and the failure mode is a rule that quietly never fires")
	}
}
