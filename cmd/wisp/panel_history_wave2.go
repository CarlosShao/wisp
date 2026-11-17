package main

// Readface wave 2 (2026-11-17), ticket 36 AC#3's Go half: the real leg behind
// internal/panel's history.page door. internal/panel names the roster door
// (bridge.go's MethodHistoryPage, routed by composer_dispatch.go's
// HistoryPageHandler) and this file attaches the handler that owns the
// objects: internal/memory's paged history read (dao_history_page_36.go's
// HistoryPage) over the two tables the D35 schema already holds - task_log
// (the task history) and tool_call (the per-call result forensics). The real
// table names and their row counts are the evidence file's to state; the
// read itself asks the store, never a second count.
//
// WHAT THIS DOOR IS NOT. No delete, no purge, no export: the destructive
// verbs of the same data already have their ONLY route (panel_privacy_40.go's
// three privacy doors over memory's privacy layer, plus ticket 373's
// DeleteTaskLog/PurgeTaskLogs), and the wave dispatch's "⛔新造第二套" holds
// by construction - this handler is one read method, nothing else.
//
// Q-50, answered: panel.view.request was DELETED by the owner's ruling (甲,
// A238) and the C17 whitelist never carried it; this face does not depend on
// it and does not resurrect it. The read rides its own named roster door, so
// the ticket-36 whitespace ("缺 SQLite 历史 reader＋panel.view.request 被
// Q-50 封") closes without touching that ruling.
//
// SHAPE. The settings-route shape (ticket 248/194/40 precedent): one prose
// sentence plus one payload line behind the marker, the payload the page
// parses for its virtualized list. The paging window comes off the envelope
// (HistoryOffset/HistoryLimit), the store refuses a window it does not serve
// by name, and the receipt states the row count and the total the store
// returned - never a recount.
//
// WHY OPEN PER READ. The same shape every read outlet in this package already
// uses (sessions.go's outlets, panel_cost_194.go, panel_privacy_40.go):
// memory.Open, one operation, defer Close. A history request is a
// user-triggered, rare operation, and a handle held for the process's whole
// life would pin wisp.db open for a machine that may never ask.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// historyPagePayloadMarker starts the history read's payload line. Same
// naming properties as every marker in this family: no `panel.` prefix and no
// dot - it is neither a method name nor a route.
const historyPagePayloadMarker = "wisp-history-page:"

// historyPageDefaultLimit is the page the door serves when the envelope
// carries no limit. A scroller that names no window still gets a bounded
// answer, never the whole table and never a refusal for the honest
// "give me the first page".
const historyPageDefaultLimit = 100

// historyHandlerWave2 is the history read leg: the data root it opens wisp.db
// over, the audit sink its refusals go to, and the actor for the audit line.
type historyHandlerWave2 struct {
	dir       string
	auditf    panel.AuditFunc
	actorName string
}

// newHistoryHandlerWave2 builds the leg over the data root the chain was
// given. No database is opened here - see the file comment.
func newHistoryHandlerWave2(dir string, auditf panel.AuditFunc, actor string) *historyHandlerWave2 {
	return &historyHandlerWave2{dir: dir, auditf: auditf, actorName: actor}
}

// HandleHistoryPage answers one parsed history request: the paging window off
// the envelope (limit 0 = the door's own default page), the store's paged
// read, the receipt sentence plus the payload line. A request with an
// unreadable window is refused before the store is opened; a store that will
// not open refuses THAT request loudly and audited.
func (h *historyHandlerWave2) HandleHistoryPage(ctx context.Context, req panel.ComposerRequest) (string, error) {
	if h == nil {
		return "", fmt.Errorf("%w: 历史读取门不存在", panel.ErrNoHandlerAttached)
	}
	if req.HistoryOffset < 0 {
		err := h.refuse(req, fmt.Errorf("historyOffset=%d 不成立：页窗口不为负", req.HistoryOffset))
		return "", err
	}
	limit := req.HistoryLimit
	if limit == 0 {
		limit = historyPageDefaultLimit
	}
	st, err := memory.Open(h.dir)
	if err != nil {
		err := h.refuse(req, fmt.Errorf("历史存储不可用（%s）：%w", h.dir, err))
		return "", err
	}
	defer func() { _ = st.Close() }()
	page, err := st.HistoryPage(ctx, req.HistoryOffset, limit)
	if err != nil {
		err := h.refuse(req, err)
		return "", err
	}
	h.record(req, "answered")
	var b strings.Builder
	fmt.Fprintf(&b, "历史共 %d 行（新的事件在前），本页第 %d 行起的 %d 行。",
		page.Total, page.Offset+1, len(page.Rows))
	b.WriteString("\n" + historyPagePayloadMarker + marshalHistoryPayloadWave2(page))
	return b.String(), nil
}

// historyPagePayload is the door's wire shape: the page's rows plus the total
// the page was read against, so a scroller can bound itself without a second
// request.
type historyPagePayload struct {
	Rows   []memory.HistoryPageRow `json:"rows"`
	Total  int64                   `json:"total"`
	Offset int                     `json:"offset"`
	Limit  int                     `json:"limit"`
}

func marshalHistoryPayloadWave2(p memory.HistoryPageTotal) string {
	data, err := json.Marshal(historyPagePayload{
		Rows: p.Rows, Total: p.Total, Offset: p.Offset, Limit: p.Limit,
	})
	if err != nil {
		return fmt.Sprintf("%s{\"payloadError\":%q}", historyPagePayloadMarker, err.Error())
	}
	return string(data)
}

func (h *historyHandlerWave2) refuse(req panel.ComposerRequest, cause error) error {
	err := fmt.Errorf("panel: history.page 被拒绝（requestId=%q）：%w", req.RequestID, cause)
	h.record(req, "refused")
	return err
}

func (h *historyHandlerWave2) actor() string {
	if h == nil || h.actorName == "" {
		return "unknown"
	}
	return h.actorName
}

// record writes one audit line per outcome; the payload is the user's own
// history and never rides the audit line - the window and the outcome are
// what the line needs.
func (h *historyHandlerWave2) record(req panel.ComposerRequest, outcome string) {
	if h == nil || h.auditf == nil {
		return
	}
	h.auditf("panel: HISTORY-ROUTE request=%q method=%q offset=%d limit=%d actor=%q outcome=%q",
		req.RequestID, req.Method, req.HistoryOffset, req.HistoryLimit, h.actor(), outcome)
}
