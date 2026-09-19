package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// orNow returns v, or the current wall-clock unix seconds when v == 0
// (persisted timestamps only — D42#9).
func orNow(v int64) int64 {
	if v != 0 {
		return v
	}
	return observe.NowWallUTC().Unix()
}

// L3 task-log DAO (D20/D35). QueryText is stored as given — redaction per
// §14.4 happens upstream before this layer is reached.

// StartTaskLog inserts a running task row (state 'running' by caller choice).
func (s *Store) StartTaskLog(ctx context.Context, tl TaskLog) error {
	tl.Currency = defaultCurrency(tl.Currency)
	tl.StartedAt = orNow(tl.StartedAt)
	if err := validateTaskLog(tl); err != nil {
		return err
	}
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO task_log(id, started_at, ended_at, state, query_text, summary_text,
			   cost_tokens_in, cost_tokens_out, cost_amount_micro, currency, error_class)
			 VALUES(?,?,NULL,?,?,?,?,?,?,?,?)`,
			tl.ID, tl.StartedAt, tl.State, tl.QueryText, tl.SummaryText,
			tl.CostTokensIn, tl.CostTokensOut, tl.CostAmountMicro, tl.Currency, nullIfEmpty(tl.ErrorClass))
		return err
	})
	if err != nil {
		return fmt.Errorf("memory: start task_log %s: %w", tl.ID, err)
	}
	return nil
}

// FinishTaskLog completes a task row: terminal state, optional summary,
// token/cost accounting and the D37 error class. ValidateErrorClass guards
// the error_class enum (ticket 03's reserved hook).
func (s *Store) FinishTaskLog(ctx context.Context, id, state, summary string, tokensIn, tokensOut, costMicro int64, errorClass string) error {
	if errorClass != "" {
		if err := validateTaskLog(TaskLog{ID: id, State: state, ErrorClass: errorClass}); err != nil {
			return err
		}
	}
	now := observe.NowWallUTC().Unix()
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`UPDATE task_log SET ended_at=?, state=?, summary_text=COALESCE(?, summary_text),
			   cost_tokens_in=?, cost_tokens_out=?, cost_amount_micro=?, error_class=NULLIF(?, '')
			 WHERE id=?`,
			now, state, nullIfEmpty(summary), tokensIn, tokensOut, costMicro, errorClass, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("task_log %s: %w", id, ErrNotFound)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("memory: finish task_log %s: %w", id, err)
	}
	return nil
}

// TaskLogByID returns one task row.
func (s *Store) TaskLogByID(ctx context.Context, id string) (TaskLog, error) {
	var tl TaskLog
	row := s.reader.QueryRowContext(ctx,
		`SELECT id, started_at, ended_at, state, query_text, summary_text,
		        cost_tokens_in, cost_tokens_out, cost_amount_micro, currency, COALESCE(error_class,'')
		 FROM task_log WHERE id=?`, id)
	if err := scanTaskLog(row.Scan, &tl); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TaskLog{}, fmt.Errorf("task_log %s: %w", id, ErrNotFound)
		}
		return TaskLog{}, fmt.Errorf("memory: task_log by id: %w", err)
	}
	return tl, nil
}

// ListTaskLogs returns task rows, newest first.
func (s *Store) ListTaskLogs(ctx context.Context) ([]TaskLog, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, started_at, ended_at, state, query_text, summary_text,
		        cost_tokens_in, cost_tokens_out, cost_amount_micro, currency, COALESCE(error_class,'')
		 FROM task_log ORDER BY started_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("memory: list task_log: %w", err)
	}
	defer rows.Close()
	var out []TaskLog
	for rows.Next() {
		var tl TaskLog
		if err := scanTaskLog(rows.Scan, &tl); err != nil {
			return nil, err
		}
		out = append(out, tl)
	}
	return out, rows.Err()
}

// DeleteTaskLog removes one task row by id.
func (s *Store) DeleteTaskLog(ctx context.Context, id string) error {
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM task_log WHERE id=?`, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("task_log %s: %w", id, ErrNotFound)
		}
		return nil
	})
}

// PurgeTaskLogs removes every task row (privacy: one-click clear). Forensic
// tool_call rows are governed by their own domain and retention window.
func (s *Store) PurgeTaskLogs(ctx context.Context) (int64, error) {
	var n int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM task_log`)
		if err != nil {
			return err
		}
		n, _ = res.RowsAffected()
		return nil
	})
	return n, err
}

type scanFunc func(dest ...any) error

func scanTaskLog(scan scanFunc, tl *TaskLog) error {
	var endedAt sql.NullInt64
	var summary sql.NullString
	if err := scan(&tl.ID, &tl.StartedAt, &endedAt, &tl.State, &tl.QueryText, &summary,
		&tl.CostTokensIn, &tl.CostTokensOut, &tl.CostAmountMicro, &tl.Currency, &tl.ErrorClass); err != nil {
		return err
	}
	if endedAt.Valid {
		v := endedAt.Int64
		tl.EndedAt = &v
	}
	if summary.Valid {
		v := summary.String
		tl.SummaryText = &v
	}
	return nil
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func defaultCurrency(c string) string {
	if c == "" {
		return "CNY"
	}
	return c
}
