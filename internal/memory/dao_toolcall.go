package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// tool_call DAO (D35): forensics for "what did it actually do to my machine".
// Enum-shaped columns (risk_level/decision/outcome/error_class) are validated
// at this boundary; the DDL stays CHECK-free by contract.

// InsertToolCall records a tool invocation (pending decision/outcome).
func (s *Store) InsertToolCall(ctx context.Context, tc ToolCall) (int64, error) {
	if err := validateToolCall(tc); err != nil {
		return 0, err
	}
	var id int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, decision, decided_at,
			   started_at, ended_at, outcome, error_class, correlation_id, grant_id)
			 VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			tc.TaskID, tc.Seq, tc.Tool, tc.ArgsJSON, tc.RiskLevel,
			nullIfEmpty(tc.Decision), nullableInt64(tc.DecidedAt),
			nullableInt64(tc.StartedAt), nullableInt64(tc.EndedAt),
			nullIfEmpty(tc.Outcome), nullIfEmpty(tc.ErrorClass), tc.CorrelationID,
			nullableInt64(tc.GrantID))
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("memory: insert tool_call: %w", err)
	}
	return id, nil
}

// DecideToolCall records the approval decision (and optional grant linkage).
func (s *Store) DecideToolCall(ctx context.Context, id int64, decision string, grantID *int64) error {
	if !toolCallDecisions[decision] {
		return fmt.Errorf("memory: invalid tool_call.decision %q", decision)
	}
	now := observe.NowWallUTC().Unix()
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`UPDATE tool_call SET decision=?, decided_at=?, grant_id=COALESCE(?, grant_id) WHERE id=?`,
			decision, now, grantID, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("tool_call %d: %w", id, ErrNotFound)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("memory: decide tool_call %d: %w", id, err)
	}
	return nil
}

// FinishToolCall records outcome and optional D37 error class.
func (s *Store) FinishToolCall(ctx context.Context, id int64, outcome, errorClass string) error {
	if outcome != "" && !toolCallOutcomes[outcome] {
		return fmt.Errorf("memory: invalid tool_call.outcome %q", outcome)
	}
	if errorClass != "" {
		if err := observe.ValidateErrorClass(errorClass); err != nil {
			return err
		}
	}
	now := observe.NowWallUTC().Unix()
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`UPDATE tool_call SET ended_at=?, outcome=?, error_class=NULLIF(?, '') WHERE id=?`,
			now, nullIfEmpty(outcome), errorClass, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("tool_call %d: %w", id, ErrNotFound)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("memory: finish tool_call %d: %w", id, err)
	}
	return nil
}

// ToolCallByID returns one forensics row.
func (s *Store) ToolCallByID(ctx context.Context, id int64) (ToolCall, error) {
	var tc ToolCall
	row := s.reader.QueryRowContext(ctx, toolCallSelect+` WHERE id=?`, id)
	if err := scanToolCall(row.Scan, &tc); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ToolCall{}, fmt.Errorf("tool_call %d: %w", id, ErrNotFound)
		}
		return ToolCall{}, fmt.Errorf("memory: tool_call by id: %w", err)
	}
	return tc, nil
}

// ListToolCallsByTask returns a task's forensics in seq order.
func (s *Store) ListToolCallsByTask(ctx context.Context, taskID string) ([]ToolCall, error) {
	rows, err := s.reader.QueryContext(ctx, toolCallSelect+` WHERE task_id=? ORDER BY seq`, taskID)
	if err != nil {
		return nil, fmt.Errorf("memory: list tool_call by task: %w", err)
	}
	defer rows.Close()
	return scanToolCalls(rows)
}

// ListToolCalls returns forensics rows newest-first (privacy page listing).
func (s *Store) ListToolCalls(ctx context.Context) ([]ToolCall, error) {
	rows, err := s.reader.QueryContext(ctx,
		toolCallSelect+` ORDER BY COALESCE(started_at, decided_at, 0) DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("memory: list tool_call: %w", err)
	}
	defer rows.Close()
	return scanToolCalls(rows)
}

// DeleteToolCall removes one forensics row by id.
func (s *Store) DeleteToolCall(ctx context.Context, id int64) error {
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM tool_call WHERE id=?`, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("tool_call %d: %w", id, ErrNotFound)
		}
		return nil
	})
}

// PurgeToolCalls removes every forensics row (privacy: one-click clear).
func (s *Store) PurgeToolCalls(ctx context.Context) (int64, error) {
	var n int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM tool_call`)
		if err != nil {
			return err
		}
		n, _ = res.RowsAffected()
		return nil
	})
	return n, err
}

const toolCallSelect = `SELECT id, task_id, seq, tool, args_json, risk_level, decision, decided_at,
	started_at, ended_at, outcome, COALESCE(error_class,''), correlation_id, grant_id
	FROM tool_call`

func scanToolCall(scan scanFunc, tc *ToolCall) error {
	var decision, outcome sql.NullString
	var decidedAt, startedAt, endedAt, grantID sql.NullInt64
	if err := scan(&tc.ID, &tc.TaskID, &tc.Seq, &tc.Tool, &tc.ArgsJSON, &tc.RiskLevel,
		&decision, &decidedAt, &startedAt, &endedAt, &outcome, &tc.ErrorClass,
		&tc.CorrelationID, &grantID); err != nil {
		return err
	}
	tc.Decision = decision.String
	tc.Outcome = outcome.String
	tc.DecidedAt = nullInt64Ptr(decidedAt)
	tc.StartedAt = nullInt64Ptr(startedAt)
	tc.EndedAt = nullInt64Ptr(endedAt)
	tc.GrantID = nullInt64Ptr(grantID)
	return nil
}

func scanToolCalls(rows *sql.Rows) ([]ToolCall, error) {
	var out []ToolCall
	for rows.Next() {
		var tc ToolCall
		if err := scanToolCall(rows.Scan, &tc); err != nil {
			return nil, err
		}
		out = append(out, tc)
	}
	return out, rows.Err()
}

func nullableInt64(v *int64) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullInt64Ptr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	v := n.Int64
	return &v
}
