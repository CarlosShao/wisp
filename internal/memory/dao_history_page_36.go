package memory

// Readface wave 2 (2026-11-17), ticket 36 AC#3's Go half: the paged HISTORY
// read over the two tables the D35 schema already holds. The real table names
// are read off this package's own schema (schema.go's ddlV1): `task_log` is
// the task history (id, started_at, state, summary_text, cost columns) and
// `tool_call` is the per-call result/forensics history (tool, outcome,
// correlation id). Both have existed since schema v1; what never existed was a
// paged read face over them - the panel's history page had no reader, which is
// exactly the gap ticket 36 AC#3 names.
//
// WHY THIS IS A READ AND NOT A SECOND PRIVACY FACE. delete-one / purge /
// export already have their ONLY route (privacy.go's DeletePrivacyItem /
// PurgePrivacy / ExportPrivacy, answered on the panel by the wave-1 privacy.*
// doors); ticket 373's DeleteTaskLog / PurgeTaskLogs stay the destructive
// verbs for the task_log domain. This file adds no delete, no purge and no
// export - one read method, paged, nothing else - so "⛔新造第二套" holds by
// construction: the destructive vocabulary of this face is zero.
//
// WHY SQL PAGING AND NOT A GO-SIDE SLICE. The face has to support 1000-row
// scrolling; ListTaskLogs already answers the whole task_log table newest
// first, and slicing that in the caller loads every row to show one page.
// LIMIT/OFFSET over one deterministic total order is the mainstream paging
// shape, and the index that backs it (idx_task_log_started) already exists.
//
// THE ORDER IS TOTAL, NOT CONVENIENT. Rows come from two tables, so the
// stable order is (ts DESC, kind DESC, id DESC): every column in the order is
// a real column, no NULL sorts ahead of a value (COALESCE keeps ts an
// integer), and page N of the same database is page N again. A page boundary
// that can reorder between requests would scroll rows past the reader, which
// is the failure virtualization exists to hide and this query must not
// produce.

import (
	"context"
	"fmt"
)

// HistoryPageRow is one history row as the panel's page reads it. Kind says
// which table the row came from ("task" | "tool") and which columns mean
// what: a task row's Outcome is the task_log state, a tool row's Outcome is
// tool_call.outcome ('success'|'error'|'cancelled'|'truncated'). Summary is
// the task's own summary text when there is one; a tool row carries no
// summary here - args_json is the §5.1-sanitized forensics copy and the
// history page does not re-expose it under a second name.
type HistoryPageRow struct {
	ID      string `json:"id"`
	TS      int64  `json:"ts"`
	Kind    string `json:"kind"`
	Tool    string `json:"tool,omitempty"`
	Outcome string `json:"outcome,omitempty"`
	Summary string `json:"summary,omitempty"`
}

// HistoryPageTotal is the count half of the paging answer: how many rows the
// whole history holds, so a scroller can bound itself without a second
// request per page. It travels with the rows it was counted with.
type HistoryPageTotal struct {
	Rows  []HistoryPageRow `json:"rows"`
	Total int64            `json:"total"`
	// Offset / Limit are the request's own answers back, clamped to what was
	// actually served - the reader's scroll state stays honest about which
	// window it is looking at.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// HistoryPageLimitCeiling is the largest page this face serves. 500 rows is
// past any single screen and keeps a runaway client from asking for the
// whole table in one page; scrolling past it is what offset is for.
const HistoryPageLimitCeiling = 500

// historyPageQuery is the union read. The columns are the two tables' own;
// nothing is derived, renamed twice or defaulted here: a task row's outcome
// IS task_log.state (the store's own task vocabulary), a tool row's outcome
// IS tool_call.outcome. tool_call.started_at is nullable by schema (a call
// rejected before dispatch books no start), so the ts half COALESCEs it with
// the ended stamp and finally 0 - a row with no time sorts last, it does not
// sort as "now".
const historyPageQuery = `SELECT id, ts, kind, tool, outcome, summary FROM (
  SELECT id, started_at AS ts, 'task' AS kind, '' AS tool,
         state AS outcome, COALESCE(summary_text,'') AS summary
    FROM task_log
  UNION ALL
  SELECT correlation_id AS id, COALESCE(ended_at, started_at, 0) AS ts, 'tool' AS kind,
         tool AS tool, COALESCE(outcome,'') AS outcome, '' AS summary
    FROM tool_call
) ORDER BY ts DESC, kind DESC, id DESC LIMIT ? OFFSET ?`

// historyPageCountQuery is the same union without the page window: the total
// the scroller bounds itself with. It counts the same rows the page query
// reads - one definition of "the history", not two.
const historyPageCountQuery = `SELECT COUNT(*) FROM (
  SELECT id FROM task_log
  UNION ALL
  SELECT correlation_id FROM tool_call
)`

// HistoryPage answers one page of the history, newest first, with the total.
// Negative offsets and non-positive or over-ceiling limits are refused by
// name before any query runs - a caller that asks for a page this face does
// not serve hears which argument was wrong, never gets a silent clamp that
// looks like data.
func (s *Store) HistoryPage(ctx context.Context, offset, limit int) (HistoryPageTotal, error) {
	if offset < 0 {
		return HistoryPageTotal{}, fmt.Errorf("memory: history.page offset=%d 不成立：页窗口不为负", offset)
	}
	if limit <= 0 {
		return HistoryPageTotal{}, fmt.Errorf("memory: history.page limit=%d 不成立：一页至少 1 行", limit)
	}
	if limit > HistoryPageLimitCeiling {
		return HistoryPageTotal{}, fmt.Errorf("memory: history.page limit=%d 超过上限 %d：翻页用 offset，不用放大页",
			limit, HistoryPageLimitCeiling)
	}
	rows, err := s.reader.QueryContext(ctx, historyPageQuery, limit, offset)
	if err != nil {
		return HistoryPageTotal{}, fmt.Errorf("memory: history.page read: %w", err)
	}
	defer rows.Close()
	out := HistoryPageTotal{Offset: offset, Limit: limit, Rows: []HistoryPageRow{}}
	for rows.Next() {
		var r HistoryPageRow
		if err := rows.Scan(&r.ID, &r.TS, &r.Kind, &r.Tool, &r.Outcome, &r.Summary); err != nil {
			return HistoryPageTotal{}, fmt.Errorf("memory: history.page scan: %w", err)
		}
		out.Rows = append(out.Rows, r)
	}
	if err := rows.Err(); err != nil {
		return HistoryPageTotal{}, fmt.Errorf("memory: history.page rows: %w", err)
	}
	if err := s.reader.QueryRowContext(ctx, historyPageCountQuery).Scan(&out.Total); err != nil {
		return HistoryPageTotal{}, fmt.Errorf("memory: history.page count: %w", err)
	}
	return out, nil
}
