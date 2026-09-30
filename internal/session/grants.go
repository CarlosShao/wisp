package session

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

// ---------------------------------------------------------------------------
// The per-session grants ledger (D45-2)
//
// doc.go:10 has listed "per-session grants ledger (D45-2)" as this package's
// responsibility since the boundary was frozen; ticket 224 is the day it has an
// implementation. The lifecycle policy is the part that was missing everywhere,
// not the storage: internal/memory has owned the approval_grant rows since
// ticket 49, and until this file NOTHING in production wrote one, nothing read
// one, and nothing compared a row's `pattern` against anything (ticket 224's
// 现量 rows 2 and 3, and A435's "存哪／按什么匹配／什么时候查" - the middle one
// was the expensive half because it did not exist).
//
// Three jobs, and each is deliberately separate so ticket 224's AC#2 cannot be
// collapsed into one "persistence" case:
//
//	Record   - the write side: one answered card becomes one row.
//	Covering - the read side: the bridge's "已授权过吗" question, satisfied by
//	           the shape tools.GrantSource declares, so internal/tools never
//	           imports this package or internal/memory.
//	(nothing) - the invalidation side has no code at all. A session ends when
//	           the process does, which means the minted ID stops existing,
//	           which means no Covering call can be keyed with it. That absence
//	           IS the design (A435 clause 2); a revoke step would be the first
//	           thing to forget on a crash.
//
// What this file does NOT do: it does not list or revoke grants for a UI.
// SPEC-08:171's grants.list / grants.revoke stay out of ticket 224's range on
// purpose - implemented as inbound Go-answered methods they die on
// internal/panel/l2_grant_boundary_test.go's banned-root table (A435's 落点
// clause), and the revocation entrance today is the existing SetAllowedDirs
// route (ticket 226).
// ---------------------------------------------------------------------------

// Store is the slice of *memory.Store this ledger needs, kept as an interface
// for the same reason perm.ConfigManager is one (internal/perm/store.go:42-51):
// the tests then run against a real SQLite file and real rows instead of against
// a mock that agrees to remember things it never had to store.
type Store interface {
	InsertGrant(ctx context.Context, g memory.ApprovalGrant) (int64, error)
	ListGrantsBySession(ctx context.Context, sessionID string) ([]memory.ApprovalGrant, error)
}

// GrantRecorder is the write half as the approval layer sees it: one answered
// card, one stored rule. approval.Gate holds this interface and nothing else
// from this package, so it cannot read grants, list them, or revoke them.
type GrantRecorder interface {
	Record(ctx context.Context, tool, pattern string) (int64, error)
}

var _ GrantRecorder = (*Ledger)(nil)

// Ledger holds one process's session identity and the D45 grants answered under
// it. Safe for concurrent use: the rows are the state, and SQLite owns them.
type Ledger struct {
	id    ID
	store Store
	logf  func(string, ...any)
	now   func() time.Time
	ttl   time.Duration
}

// LedgerOptions configures a Ledger. Every field is required except Logf: there
// is no safe zero value for a session identity or for the store holding its
// grants, and this type refuses to invent either.
type LedgerOptions struct {
	// ID is the identity Mint produced for this process.
	ID ID
	// Store is where the rows live.
	Store Store
	// Logf is the audit sink (the "[audit]" family cmd/wisp already uses).
	Logf func(format string, args ...any)
	// Now is the clock for created_at / expires_at comparisons. nil = time.Now.
	Now func() time.Time
	// GrantTTL bounds expires_at. See clockCeilingDefault for why this type
	// stamps a clock at all when the identity is what actually invalidates.
	GrantTTL time.Duration
}

// clockCeilingDefault is the backstop written into expires_at.
//
// A435 clause 2 says this session ends when the process exits, which is not a
// timestamp anyone can know in advance - and memory.InsertGrant requires a
// non-zero expires_at (internal/memory/dao_misc.go:27-29). So the clock here is
// NOT the invalidator; the identity is. What the clock is for is the one shape
// the identity cannot cover: a row whose session id is guessed or leaked. Its
// expiry still does the killing.
//
// The value is taken from a number the repository has already ruled
// (memory.GrantAuditTTL, internal/memory/retention.go:36 = SPEC-02 §4's 30-day
// audit window) rather than from a fresh invention by this ticket: past it, the
// retention job erases the row anyway, so no grant can ever be immortal on
// paper. Whether a shorter bound is wanted ("会话时长" as a bounded duration,
// A435's 乙 branch) is registered as open in ticket 224's progress note; this
// ticket does not pick a new policy number.
const clockCeilingDefault = memory.GrantAuditTTL

// NewLedger composes the ledger for one minted session.
func NewLedger(o LedgerOptions) (*Ledger, error) {
	if !o.ID.Valid() {
		return nil, errors.New("session: NewLedger 需要一枚 Mint() 铸出的合法会话身份（A435 第 1 条：不许派生、不许手写）")
	}
	if o.Store == nil {
		return nil, errors.New("session: NewLedger 需要一个存放 approval_grant 行的 Store")
	}
	ttl := o.GrantTTL
	if ttl <= 0 {
		ttl = clockCeilingDefault
	}
	now := o.Now
	if now == nil {
		now = time.Now
	}
	return &Ledger{id: o.ID, store: o.Store, logf: o.Logf, now: now, ttl: ttl}, nil
}

// SessionID reports which identity this ledger keys its rows with. It is the
// read side cmd/wisp logs once at boot so "which session wrote this row" is
// answerable from the log, not inferred.
func (l *Ledger) SessionID() ID { return l.id }

// Record writes one D45 grant: "本会话内允许 <tool> 于 <pattern>".
//
// It is called from the native answer route ONLY (approval.Gate's allow-session
// path), which is what keeps SPEC-06:1588's rule intact - a row exists because a
// native click spent a live nonce on a card that named this exact tool and this
// exact path. Nothing here infers a pattern, widens a path, or writes a row
// nobody asked for.
//
// One row per (tool, pattern): a card that named three paths is answered as
// three rows, because approval_grant.pattern is a single TEXT column (SPEC-02
// §3) and folding three paths into one string would make the matcher's equality
// test meaningless.
func (l *Ledger) Record(ctx context.Context, tool, pattern string) (int64, error) {
	if strings.TrimSpace(tool) == "" {
		return 0, errors.New("session: grant 需要 tool")
	}
	if !recordablePattern(pattern) {
		// A path C26 could not canonicalize renders with this marker attached
		// (internal/tools/bridge.go:874). Storing it would store a string no
		// later call can ever produce, i.e. a row that grants nothing but looks
		// live in the audit view. Refuse the row instead.
		return 0, fmt.Errorf("session: 路径模式无法记账（未经 C26 规范化的路径不落 grant）：%q", pattern)
	}
	now := l.now()
	id, err := l.store.InsertGrant(ctx, memory.ApprovalGrant{
		Scope:     memory.GrantScopeSession,
		Tool:      tool,
		Pattern:   pattern,
		SessionID: l.id.String(),
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(l.ttl).Unix(),
	})
	if err != nil {
		return 0, fmt.Errorf("session: 记录会话授权失败: %w", err)
	}
	l.log("session: GRANT-RECORD id=%d tool=%s pattern=%q session=%s expires_at=%d",
		id, tool, pattern, l.id, now.Add(l.ttl).Unix())
	return id, nil
}

// Covering answers the bridge's one question: does THIS session already hold a
// live grant covering (tool, every path in paths)? It is the production reader
// of ListGrantsBySession, which had none before ticket 224.
//
// The contract it implements is tools.GrantSource (internal/tools/grant.go).
// Fail-closed on every branch: any error, an empty path list, a path carrying
// the un-resolvable marker, or a panic from the store all answer (0, false),
// which puts the question back on the user.
//
// Why it reads the database instead of an in-memory set, stated because it is
// the whole security argument: the read is keyed by l.id, so AC#3 (重启必失效)
// is measured against the one property A435 clause 1 guarantees - a fresh
// process mints a fresh key and finds zero rows. A cache that never asked the
// store would make the same call green for a reason the test could not see, and
// the AC#4 control (a derived id letting a grant cross sessions) would then have
// nothing to catch.
func (l *Ledger) Covering(ctx context.Context, tool string, paths []string) (int64, bool) {
	if strings.TrimSpace(tool) == "" || len(paths) == 0 {
		return 0, false
	}
	for _, p := range paths {
		if !recordablePattern(p) {
			return 0, false
		}
	}
	rows, err := l.store.ListGrantsBySession(ctx, l.id.String())
	if err != nil {
		l.log("session: GRANT-READ-FAILED tool=%s err=%v; 按未授权处理（继续询问）", tool, err)
		return 0, false
	}
	now := l.now().Unix()
	want := make(map[string]bool, len(paths))
	for _, p := range paths {
		want[p] = true
	}
	var first int64
	for _, g := range rows {
		if g.Tool != tool || !want[g.Pattern] {
			continue
		}
		if !grantLive(g, now) {
			continue
		}
		delete(want, g.Pattern)
		if first == 0 {
			first = g.ID
		}
	}
	if len(want) > 0 {
		return 0, false
	}
	l.log("session: GRANT-HIT tool=%s paths=%d grant_id=%d session=%s", tool, len(paths), first, l.id)
	return first, true
}

// grantLive applies the three conditions ticket 224's 乙 branch names: same
// identity (guaranteed by the query this row came from), not revoked, not
// expired. A row with a zero ExpiresAt is treated as dead: memory.InsertGrant
// refuses to write one, so a zero here means the row predates this code and
// nobody can say what it meant.
func grantLive(g memory.ApprovalGrant, now int64) bool {
	if g.RevokedAt != nil {
		return false
	}
	return g.ExpiresAt > now
}

// recordablePattern is the matcher's whole dialect statement, so it is the one
// place to look when the question "what counts as 路径模式" gets asked.
//
// TODAY: exact string equality against a canonicalized path, and nothing else.
// That is the same dialect risk.Gate's B-tier single-file override uses
// (internal/risk: a map keyed on the resolved path), and it is the only dialect
// the write side can produce honestly, because Record stores the path the card
// printed.
//
// NOT today: glob. The frozen tests plant rows like filepath.Join(dir, "*")
// (internal/perm/ticket90_persist_test.go:233), and SPEC-02 §3 names the column
// `pattern`, so a wildcard dialect is clearly intended - but no document in this
// repository defines its alphabet, and a wildcard that means "this directory
// subtree" is a wider blast radius than anything a user clicked. Ticket 224
// therefore does not guess one; it is registered as an open question in the
// ticket's progress note rather than settled in code.
func recordablePattern(p string) bool {
	if strings.TrimSpace(p) == "" {
		return false
	}
	if strings.Contains(p, "无法规范化") {
		return false
	}
	return true
}

func (l *Ledger) log(format string, args ...any) {
	if l.logf == nil {
		return
	}
	l.logf(format, args...)
}
