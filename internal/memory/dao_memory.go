package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// L2 explicit-memory DAO (D20/D35). Retrieval is deliberately simple per
// SPEC-02 §3: keywords LIKE matching + most-recently-hit-first ordering.
// FTS5 waits until the table exceeds ~500 rows (a future schema migration).

// AddMemory stores one explicit memory. Keywords are normalized (lowercase,
// whitespace-split, deduped, space-joined); empty keywords after
// normalization are rejected because the row would be unsearchable.
func (s *Store) AddMemory(ctx context.Context, content string, keywords ...string) (int64, error) {
	if strings.TrimSpace(content) == "" {
		return 0, fmt.Errorf("memory: memory.content is required")
	}
	kw := normalizeKeywords(strings.Join(keywords, " "))
	if kw == "" {
		return 0, fmt.Errorf("memory: memory.keywords is required after normalization")
	}
	now := observe.NowWallUTC().Unix()
	var id int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO memory(kind, content, keywords, created_at, last_hit_at, hit_count)
			 VALUES(?,?,?,?,0,0)`,
			MemoryKindExplicit, content, kw, now)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("memory: add memory: %w", err)
	}
	return id, nil
}

// SearchMemory matches rows whose keywords contain ANY of the query terms
// (LIKE, wildcard-escaped), ordered most-recently-hit first, then by hit
// count. Every returned row's last_hit_at/hit_count is updated through the
// db-writer ("recently used first" needs real hit statistics; the update is
// part of the search contract).
func (s *Store) SearchMemory(ctx context.Context, query string, limit int) ([]Memory, error) {
	terms := strings.Fields(strings.ToLower(query))
	if len(terms) == 0 {
		return nil, fmt.Errorf("memory: search query is empty")
	}
	if limit <= 0 {
		limit = 10
	}

	clauses := make([]string, len(terms))
	args := make([]any, 0, len(terms))
	for i, t := range terms {
		clauses[i] = `keywords LIKE ? ESCAPE '\'`
		args = append(args, "%"+escapeLike(t)+"%")
	}
	q := `SELECT id, kind, content, keywords, created_at, last_hit_at, hit_count
	      FROM memory WHERE ` + strings.Join(clauses, " OR ") + `
	      ORDER BY last_hit_at DESC, hit_count DESC, id DESC LIMIT ?`
	args = append(args, limit)

	rows, err := s.reader.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("memory: search memory: %w", err)
	}
	defer rows.Close()
	hits, err := scanMemories(rows)
	if err != nil {
		return nil, err
	}

	if len(hits) > 0 {
		// Hit-statistics update goes through the single writer like every
		// other write.
		now := observe.NowWallUTC().Unix()
		ids := make([]int64, len(hits))
		for i, m := range hits {
			ids[i] = m.ID
		}
		if err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
			return updateMemoryHits(ctx, tx, ids, now)
		}); err != nil {
			return nil, fmt.Errorf("memory: update hit stats: %w", err)
		}
	}
	return hits, nil
}

func updateMemoryHits(ctx context.Context, tx *sql.Tx, ids []int64, now int64) error {
	var b strings.Builder
	b.WriteString(`UPDATE memory SET last_hit_at=?, hit_count=hit_count+1 WHERE id IN (`)
	for i, id := range ids {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(fmt.Sprintf("%d", id))
	}
	b.WriteString(`)`)
	_, err := tx.ExecContext(ctx, b.String(), now)
	return err
}

// MemoryByID returns one row (ErrNotFound when absent).
func (s *Store) MemoryByID(ctx context.Context, id int64) (Memory, error) {
	var m Memory
	err := s.reader.QueryRowContext(ctx,
		`SELECT id, kind, content, keywords, created_at, last_hit_at, hit_count FROM memory WHERE id=?`, id).
		Scan(&m.ID, &m.Kind, &m.Content, &m.Keywords, &m.CreatedAt, &m.LastHitAt, &m.HitCount)
	if errors.Is(err, sql.ErrNoRows) {
		return Memory{}, fmt.Errorf("memory %d: %w", id, ErrNotFound)
	}
	if err != nil {
		return Memory{}, fmt.Errorf("memory: memory by id: %w", err)
	}
	return m, nil
}

// ListMemories returns rows newest-created first (privacy page listing).
func (s *Store) ListMemories(ctx context.Context) ([]Memory, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, kind, content, keywords, created_at, last_hit_at, hit_count
		 FROM memory ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("memory: list memories: %w", err)
	}
	defer rows.Close()
	return scanMemories(rows)
}

// DeleteMemory removes one memory row by id.
func (s *Store) DeleteMemory(ctx context.Context, id int64) error {
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM memory WHERE id=?`, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("memory %d: %w", id, ErrNotFound)
		}
		return nil
	})
}

// PurgeMemories removes every L2 row (privacy: one-click clear).
func (s *Store) PurgeMemories(ctx context.Context) (int64, error) {
	var n int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM memory`)
		if err != nil {
			return err
		}
		n, _ = res.RowsAffected()
		return nil
	})
	return n, err
}

func scanMemories(rows *sql.Rows) ([]Memory, error) {
	var out []Memory
	for rows.Next() {
		var m Memory
		if err := rows.Scan(&m.ID, &m.Kind, &m.Content, &m.Keywords, &m.CreatedAt, &m.LastHitAt, &m.HitCount); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
