package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Session storage (ticket 166, schema v3). One `session` row is one
// cross-restart conversation; its `session_turn` rows are the "正文" the
// search reads. Lifecycle POLICY (when a session is created, when it is
// archived, whether forks copy turns) lives in internal/session — this layer
// only stores rows, exactly like the task_log DAO it sits next to.

// Session is one conversation row (session table).
type Session struct {
	ID         string
	Title      string
	Workspace  string // "" = unassigned
	CreatedAt  int64
	UpdatedAt  int64
	ArchivedAt *int64 // nil = active; set = archived (undo clears it)
	ForkedFrom string // "" = NULL (an original, not a fork)
}

// SessionTurn is one message row (session_turn table).
type SessionTurn struct {
	ID        int64
	SessionID string
	Seq       int64
	Role      string // 'user' | 'assistant' | …
	Content   string
	CreatedAt int64
}

// SessionListOptions filters ListSessions. The zero value lists every ACTIVE
// session across all workspaces (archived ones stay hidden until asked for).
type SessionListOptions struct {
	Workspace       string // "" = every workspace
	IncludeArchived bool
}

func validateSession(s Session) error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("memory: session.id is required")
	}
	if strings.TrimSpace(s.Title) == "" {
		return fmt.Errorf("memory: session.title is required")
	}
	return nil
}

func validateTurn(t SessionTurn) error {
	if strings.TrimSpace(t.SessionID) == "" {
		return fmt.Errorf("memory: session_turn.session_id is required")
	}
	if strings.TrimSpace(t.Role) == "" {
		return fmt.Errorf("memory: session_turn.role is required")
	}
	if t.Content == "" {
		return fmt.Errorf("memory: session_turn.content is required")
	}
	return nil
}

const sessionCols = `id, title, workspace, created_at, updated_at, archived_at, COALESCE(forked_from,'')`

// InsertSession creates one conversation row. The id is caller-minted (the
// session layer mints it so storage never invents identity).
func (s *Store) InsertSession(ctx context.Context, sess Session) error {
	sess.CreatedAt = orNow(sess.CreatedAt)
	if sess.UpdatedAt == 0 {
		sess.UpdatedAt = sess.CreatedAt
	}
	if err := validateSession(sess); err != nil {
		return err
	}
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO session(id, title, workspace, created_at, updated_at, archived_at, forked_from)
			 VALUES(?,?,?,?,?,?,?)`,
			sess.ID, sess.Title, sess.Workspace, sess.CreatedAt, sess.UpdatedAt,
			nullableInt64(sess.ArchivedAt), nullIfEmpty(sess.ForkedFrom))
		return err
	})
	if err != nil {
		return fmt.Errorf("memory: insert session %s: %w", sess.ID, err)
	}
	return nil
}

// SessionByID returns one conversation row.
func (s *Store) SessionByID(ctx context.Context, id string) (Session, error) {
	var sess Session
	var archived sql.NullInt64
	err := s.reader.QueryRowContext(ctx,
		`SELECT `+sessionCols+` FROM session WHERE id=?`, id).
		Scan(&sess.ID, &sess.Title, &sess.Workspace, &sess.CreatedAt, &sess.UpdatedAt, &archived, &sess.ForkedFrom)
	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, fmt.Errorf("session %s: %w", id, ErrNotFound)
	}
	if err != nil {
		return Session{}, fmt.Errorf("memory: session by id: %w", err)
	}
	sess.ArchivedAt = nullInt64Ptr(archived)
	return sess, nil
}

// ListSessions returns conversations newest-first, optionally narrowed to one
// workspace and optionally including archived ones.
func (s *Store) ListSessions(ctx context.Context, opts SessionListOptions) ([]Session, error) {
	q := `SELECT ` + sessionCols + ` FROM session WHERE 1=1`
	var args []any
	if !opts.IncludeArchived {
		q += ` AND archived_at IS NULL`
	}
	if opts.Workspace != "" {
		q += ` AND workspace=?`
		args = append(args, opts.Workspace)
	}
	q += ` ORDER BY updated_at DESC, id DESC`
	return s.querySessions(ctx, q, args...)
}

// SearchSessions returns conversations whose TITLE or any TURN CONTENT contains
// q (LIKE, case-insensitive for ASCII via SQLite's default NOCASE-less LIKE;
// wildcards in q are escaped, so a literal `%` searches for `%`).
func (s *Store) SearchSessions(ctx context.Context, q string, opts SessionListOptions) ([]Session, error) {
	like := "%" + escapeLike(q) + "%"
	sqlText := `SELECT ` + sessionCols + ` FROM session WHERE (title LIKE ? ESCAPE '\' OR EXISTS (
		SELECT 1 FROM session_turn t WHERE t.session_id = session.id AND t.content LIKE ? ESCAPE '\'))`
	args := []any{like, like}
	if !opts.IncludeArchived {
		sqlText += ` AND archived_at IS NULL`
	}
	if opts.Workspace != "" {
		sqlText += ` AND workspace=?`
		args = append(args, opts.Workspace)
	}
	sqlText += ` ORDER BY updated_at DESC, id DESC`
	return s.querySessions(ctx, sqlText, args...)
}

func (s *Store) querySessions(ctx context.Context, q string, args ...any) ([]Session, error) {
	rows, err := s.reader.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("memory: query sessions: %w", err)
	}
	defer rows.Close()
	var out []Session
	for rows.Next() {
		var sess Session
		var archived sql.NullInt64
		if err := rows.Scan(&sess.ID, &sess.Title, &sess.Workspace, &sess.CreatedAt,
			&sess.UpdatedAt, &archived, &sess.ForkedFrom); err != nil {
			return nil, err
		}
		sess.ArchivedAt = nullInt64Ptr(archived)
		out = append(out, sess)
	}
	return out, rows.Err()
}

// RenameSession writes a new title. It touches only the title column — no id,
// no turn, no timestamp of the ORIGINAL history is rewritten (ticket 166's
// "真相源只追加" rule; renaming is not a rewrite of the conversation).
func (s *Store) RenameSession(ctx context.Context, id, title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("memory: session title is required")
	}
	now := observe.NowWallUTC().Unix()
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `UPDATE session SET title=?, updated_at=? WHERE id=?`, title, now, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("session %s: %w", id, ErrNotFound)
		}
		return nil
	})
}

// SetSessionArchived archives (archived=true) or restores (false) a session.
// Archiving is reversible by construction: the undo is the same call with
// archived=false, and neither branch deletes a row or a turn.
//
// The two branches write archived_at and NOTHING ELSE (ticket 166 AC#2, item 4:
// "归档带撤销"). updated_at is the row's last-activity stamp — it is what
// ListSessions/SearchSessions order by — so a stamp written here would buy two
// lies at once: archiving would count as activity, and the UNDO would move the
// row to the top of "newest first", i.e. undoing would not put the conversation
// back where it was. An undo that reorders the list is not an undo; the ticket's
// whole point is that archive is a mask the user can take back.
func (s *Store) SetSessionArchived(ctx context.Context, id string, archived bool) error {
	var stamp any
	if archived {
		stamp = observe.NowWallUTC().Unix()
	}
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `UPDATE session SET archived_at=? WHERE id=?`, stamp, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("session %s: %w", id, ErrNotFound)
		}
		return nil
	})
}

// ForkSession creates newID as a branch of parentID: the parent's title and
// workspace are copied and every parent turn is COPIED (new rows) so the fork
// starts with the history up to this point. The parent's row AND turns are left
// byte-for-byte untouched — a fork that "tidies up" the original (moving or
// deleting its turns) is exactly the mutation ticket 166 AC#4 exists to catch.
func (s *Store) ForkSession(ctx context.Context, parentID, newID string) (Session, error) {
	if strings.TrimSpace(newID) == "" {
		return Session{}, fmt.Errorf("memory: fork needs a new session id")
	}
	now := observe.NowWallUTC().Unix()
	var forked Session
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `SELECT `+sessionCols+` FROM session WHERE id=?`, parentID)
		var archived sql.NullInt64
		if err := row.Scan(&forked.ID, &forked.Title, &forked.Workspace, &forked.CreatedAt,
			&forked.UpdatedAt, &archived, &forked.ForkedFrom); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("session %s: %w", parentID, ErrNotFound)
			}
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO session(id, title, workspace, created_at, updated_at, archived_at, forked_from)
			 VALUES(?,?,?,?,?,NULL,?)`,
			newID, forked.Title, forked.Workspace, now, now, parentID); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx,
			`INSERT INTO session_turn(session_id, seq, role, content, created_at)
			 SELECT ?, seq, role, content, created_at FROM session_turn WHERE session_id=? ORDER BY seq, id`,
			newID, parentID)
		return err
	})
	if err != nil {
		return Session{}, fmt.Errorf("memory: fork session %s -> %s: %w", parentID, newID, err)
	}
	forked.ID = newID
	forked.CreatedAt, forked.UpdatedAt = now, now
	forked.ArchivedAt = nil
	forked.ForkedFrom = parentID
	return forked, nil
}

// AppendSessionTurn adds one message to a conversation (seq auto-incremented).
//
// Appending a turn IS a touch of the conversation, so the same transaction
// refreshes session.updated_at to the stamp of the turn it just wrote (ticket 326
// AC#1): without it, "list newest-first" sorts by the last time the user renamed
// or forked the conversation and never by the last thing they said, so a live
// chat cannot rise to the top of the list. (Archiving is deliberately NOT in that
// list of touches any more — see SetSessionArchived: undo has to be exact.)
//
// Two constraints shape the how:
//   - ONE transaction. The INSERT and the UPDATE share s.write, so a kill between
//     them cannot leave a turn whose conversation still carries an older
//     updated_at — the HALF-STATE ticket 326 AC#3 kills a subprocess to prove
//     unreachable.
//   - ONE clock. The value written is t.CreatedAt, already normalized by orNow
//     (== observe.NowWallUTC().Unix() when the caller passed none), which is the
//     family the row itself was stamped with; no second clock is read inside the
//     transaction (D32/D38: never derive a timeout or a stamp from wall-clock
//     differences).
//
// A failing refresh returns an error and rolls the whole transaction back: this
// never degrades silently into "child row only".
func (s *Store) AppendSessionTurn(ctx context.Context, t SessionTurn) (int64, error) {
	t.CreatedAt = orNow(t.CreatedAt)
	if err := validateTurn(t); err != nil {
		return 0, err
	}
	var id int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if t.Seq == 0 {
			if err := tx.QueryRowContext(ctx,
				`SELECT COALESCE(MAX(seq),0)+1 FROM session_turn WHERE session_id=?`, t.SessionID).Scan(&t.Seq); err != nil {
				return err
			}
		}
		res, err := tx.ExecContext(ctx,
			`INSERT INTO session_turn(session_id, seq, role, content, created_at) VALUES(?,?,?,?,?)`,
			t.SessionID, t.Seq, t.Role, t.Content, t.CreatedAt)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		// The touch half of "每一轮都算动过" (ticket 326). A conversation row that
		// is not in the table is not this call's business to police — storage
		// stores rows, the lifecycle policy lives in internal/session — so the
		// statement's error is what propagates, nothing else.
		_, err = tx.ExecContext(ctx, `UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt, t.SessionID)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("memory: append session turn: %w", err)
	}
	return id, nil
}

// ListSessionTurns returns one conversation's turns in order.
func (s *Store) ListSessionTurns(ctx context.Context, sessionID string) ([]SessionTurn, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, session_id, seq, role, content, created_at
		 FROM session_turn WHERE session_id=? ORDER BY seq, id`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("memory: list session turns: %w", err)
	}
	defer rows.Close()
	var out []SessionTurn
	for rows.Next() {
		var t SessionTurn
		if err := rows.Scan(&t.ID, &t.SessionID, &t.Seq, &t.Role, &t.Content, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
