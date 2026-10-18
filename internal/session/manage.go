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
// The Go-side session management layer (ticket 166)
//
// doc.go has listed "the wake -> explicit end / 90s-idle lifecycle" as this
// package's responsibility since the boundary was frozen; that is the WITHIN
// -BOOT scope (C31). Ticket 166 adds the thing nothing in production had: an
// outlet for the CROSS-RESTART conversation — list (by time / by workspace),
// search (title + body), rename, archive WITH an undo, and a fork that does not
// swallow the original.
//
// This file is deliberately the only place that turns "the four outlets" into
// calls on storage; internal/memory/dao_session.go stores rows and decides no
// policy. The storage slice is an interface for the same reason Ledger.Store is
// one (grants.go:47-54): the tests then run against a real SQLite file and real
// rows, not a mock that agrees to remember things it never had to store.
//
// What this file does NOT do: it does not mint the ON-DISK identity as a
// security key. That is Mint()'s job for the D45 ledger, where a derived id
// would be a hole. Here an id is a durable primary key, and Mint is reused only
// because it is already the repository's one "crypto/rand + hex" mint.
// ---------------------------------------------------------------------------

// ManageStore is the storage slice this layer needs (satisfied by
// *memory.Store).
type ManageStore interface {
	InsertSession(ctx context.Context, s memory.Session) error
	SessionByID(ctx context.Context, id string) (memory.Session, error)
	ListSessions(ctx context.Context, opts memory.SessionListOptions) ([]memory.Session, error)
	SearchSessions(ctx context.Context, q string, opts memory.SessionListOptions) ([]memory.Session, error)
	RenameSession(ctx context.Context, id, title string) error
	SetSessionArchived(ctx context.Context, id string, archived bool) error
	ForkSession(ctx context.Context, parentID, newID string) (memory.Session, error)
	AppendSessionTurn(ctx context.Context, t memory.SessionTurn) (int64, error)
	ListSessionTurns(ctx context.Context, sessionID string) ([]memory.SessionTurn, error)
}

// ManageAssemblyGuard is the compile-time proof that the host shape assembles:
// the value a production root gets from memory.Open is a concrete *memory.Store
// (cmd/wisp/run.go:463 today), and this line asserts it satisfies ManageStore
// with NO adapter — so the one production line ticket 166 AC#2 still waits for,
// the host-qualified NewManager(ManagerOptions{Store: mem}) call from cmd/wisp,
// compiles the moment the shared compile surface gets its own wave.
// It lives in production code, not a test, because 166-v1's ruling is that
// "tests green + functions complete" is NOT wiring (生产调用者 0/0/0/0/0): only
// a byte that the compiler enforces on every build belongs on the assembly
// axis. If *memory.Store ever drifts from ManageStore, this is the line that
// fails the build before any host can be written against the drift.
var _ ManageStore = (*memory.Store)(nil)

// Manager is the session-management outlet set. Every method is safe for
// concurrent use: the rows are the state and SQLite owns them.
type Manager struct {
	store ManageStore
	mint  func() (ID, error)
	now   func() time.Time
}

// ManagerOptions configures a Manager. Store is required; Mint and Now have
// safe defaults (Mint, time.Now) but are overridable so tests can pin ids and
// clocks.
type ManagerOptions struct {
	Store ManageStore
	Mint  func() (ID, error)
	Now   func() time.Time
}

// NewManager composes the outlet set. A nil store is refused rather than
// papered over: without rows there is nothing to list, so fail at composition.
func NewManager(o ManagerOptions) (*Manager, error) {
	if o.Store == nil {
		return nil, errors.New("session: NewManager 需要一枚存放会话行的 ManageStore")
	}
	m := &Manager{store: o.Store, mint: o.Mint, now: o.Now}
	if m.mint == nil {
		m.mint = Mint
	}
	if m.now == nil {
		m.now = time.Now
	}
	return m, nil
}

// Create starts a new conversation and returns its row. Title is required
// (a nameless list row cannot be renamed back to anything); workspace may be "".
func (m *Manager) Create(ctx context.Context, title, workspace string) (memory.Session, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return memory.Session{}, errors.New("session: 会话标题不能为空")
	}
	id, err := m.mint()
	if err != nil {
		return memory.Session{}, fmt.Errorf("session: 铸造会话 id 失败: %w", err)
	}
	now := m.now().Unix()
	sess := memory.Session{ID: id.String(), Title: title, Workspace: strings.TrimSpace(workspace), CreatedAt: now, UpdatedAt: now}
	if err := m.store.InsertSession(ctx, sess); err != nil {
		return memory.Session{}, err
	}
	return sess, nil
}

// List is outlet #1: conversations newest-first, optionally by workspace.
// Archived ones stay hidden unless opts asks for them.
func (m *Manager) List(ctx context.Context, opts memory.SessionListOptions) ([]memory.Session, error) {
	return m.store.ListSessions(ctx, opts)
}

// Search is outlet #2: title + body ("正文") match. An empty query returns no
// rows rather than everything — "search for nothing" is a request for nothing.
func (m *Manager) Search(ctx context.Context, q string, opts memory.SessionListOptions) ([]memory.Session, error) {
	if strings.TrimSpace(q) == "" {
		return nil, nil
	}
	return m.store.SearchSessions(ctx, q, opts)
}

// Rename is outlet #3. It rewrites exactly one title column; no id, no turn and
// no createdAt of the original history changes (ticket 166's append-only rule).
func (m *Manager) Rename(ctx context.Context, id, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("session: 会话标题不能为空")
	}
	return m.store.RenameSession(ctx, id, title)
}

// Archive is outlet #4 (the forward half).
func (m *Manager) Archive(ctx context.Context, id string) error {
	return m.store.SetSessionArchived(ctx, id, true)
}

// Unarchive is outlet #4's undo. It is a first-class outlet, not an afterthought
// on Archive: archiving must be reachable back out of, and the test that pairs
// them is what pins "归档带撤销" as reversible rather than one-way.
func (m *Manager) Unarchive(ctx context.Context, id string) error {
	return m.store.SetSessionArchived(ctx, id, false)
}

// Fork branches a conversation. The new row records forked_from=id and COPIES
// the parent's turns; the parent's row and turns are untouched (AC#4).
func (m *Manager) Fork(ctx context.Context, id string) (memory.Session, error) {
	newID, err := m.mint()
	if err != nil {
		return memory.Session{}, fmt.Errorf("session: 铸造分叉 id 失败: %w", err)
	}
	return m.store.ForkSession(ctx, id, newID.String())
}

// AppendTurn adds one message to a conversation (role: 'user'/'assistant'/…).
func (m *Manager) AppendTurn(ctx context.Context, sessionID, role, content string) (int64, error) {
	return m.store.AppendSessionTurn(ctx, memory.SessionTurn{
		SessionID: sessionID, Role: strings.TrimSpace(role), Content: content,
		CreatedAt: m.now().Unix(),
	})
}

// Turns returns one conversation's messages in order.
func (m *Manager) Turns(ctx context.Context, sessionID string) ([]memory.SessionTurn, error) {
	return m.store.ListSessionTurns(ctx, sessionID)
}
