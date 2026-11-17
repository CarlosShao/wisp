package main

// Readface wave 2 (2026-11-17), ticket 36 AC#3's cmd/wisp teeth: the real
// history.page leg over a LIVE SQLite database (memory.Open on a real data
// root, real tables task_log and tool_call), and the wiring nail for the
// production chain.
//
// The e2e reads are against rows seeded through the store's own DAOs
// (StartTaskLog / InsertToolCall / FinishToolCall) - the same writers
// production uses - never against hand-cut SQL, so what the door answers is
// what the schema + DAO pair actually produces.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// seedHistoryTask inserts one task_log row through the store's own writer.
func seedHistoryTask(t *testing.T, st *memory.Store, id string, ts int64, state, summary string) {
	t.Helper()
	var sum *string
	if summary != "" {
		s := summary
		sum = &s
	}
	if err := st.StartTaskLog(context.Background(), memory.TaskLog{
		ID: id, StartedAt: ts, State: state, QueryText: "票36 历史读面种子 " + id,
		SummaryText: sum,
	}); err != nil {
		t.Fatalf("seed task_log %s: %v", id, err)
	}
	if state != string(string(memory.TaskStateRunning)) && state != "" {
		if err := st.FinishTaskLog(context.Background(), id, state, summary, 0, 0, 0, ""); err != nil {
			t.Fatalf("finish task_log %s: %v", id, err)
		}
	}
}

func seedHistoryTool(t *testing.T, st *memory.Store, taskID, corr, tool string, ts int64, outcome string) {
	t.Helper()
	tc := memory.ToolCall{
		TaskID: taskID, Seq: ts, Tool: tool, ArgsJSON: "{}", RiskLevel: "L0",
		StartedAt: &ts, CorrelationID: corr,
	}
	if outcome != "" {
		// Booked after the fact, the DAO's own vocabulary (stampToolCallInsert):
		// caller-supplied instants win, so the history ts is the row's own
		// ended stamp and the ordering stays deterministic - no wall clock in
		// the expected order.
		o := outcome
		tc.Outcome = o
		tc.EndedAt = &ts
	}
	id, err := st.InsertToolCall(context.Background(), tc)
	if err != nil {
		t.Fatalf("seed tool_call %s: %v", corr, err)
	}
	_ = id
}

// TestHistoryPageDoorAnswersTheRealTables drives the door end to end over a
// live store: newest first across BOTH tables, the window the envelope asked
// for, and the payload line the page parses.
func TestHistoryPageDoorAnswersTheRealTables(t *testing.T) {
	dir := t.TempDir()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	seedHistoryTask(t, st, "task-a", 1000, string(memory.TaskStateDone), "已完成的任务甲")
	seedHistoryTask(t, st, "task-b", 2000, string(memory.TaskStateRunning), "")
	seedHistoryTool(t, st, "task-a", "task-a#call-1", "fs.read", 1500, "success")
	seedHistoryTool(t, st, "task-a", "task-a#call-2", "fs.write", 3000, "")
	_ = st.Close()

	h := newHistoryHandlerWave2(dir, nil, "test")
	req := panel.ComposerRequest{
		Method: panel.MethodHistoryPage, RequestID: "rw2-h", Source: panel.ComposerRequestSource,
		HistoryOffset: 0, HistoryLimit: 3,
	}
	reply, err := h.HandleHistoryPage(context.Background(), req)
	if err != nil {
		t.Fatalf("history.page refused: %v", err)
	}
	line, payload, ok := strings.Cut(reply, "\n"+historyPagePayloadMarker)
	if !ok {
		t.Fatalf("receipt carries no payload line behind the marker:\n%s", reply)
	}
	var got memory.HistoryPageTotal
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("payload is not the paging shape: %v (%s)", err, payload)
	}
	if got.Total != 4 {
		t.Errorf("total = %d, want 4 (two task rows + two tool rows): %+v", got.Total, got)
	}
	if len(got.Rows) != 3 || got.Offset != 0 || got.Limit != 3 {
		t.Errorf("window = off %d / limit %d / %d rows, want 0/3/3", got.Offset, got.Limit, len(got.Rows))
	}
	// Newest first across both tables, ts the store's own vocabulary: the
	// pending fs.write carries its booked start (3000), then the running task
	// (2000), then the booked-after-the-fact fs.read (its own end, 1500).
	wantLead := []struct{ kind, tool, outcome string }{
		{"tool", "fs.write", ""},
		{"task", "", string(memory.TaskStateRunning)},
		{"tool", "fs.read", "success"},
	}
	for i, w := range wantLead {
		r := got.Rows[i]
		if r.Kind != w.kind || r.Tool != w.tool || r.Outcome != w.outcome {
			t.Errorf("row %d = {%s %s %s}, want {%s %s %s}; page=%+v", i, r.Kind, r.Tool, r.Outcome, w.kind, w.tool, w.outcome, got.Rows)
		}
	}
	if !strings.Contains(line, "4") {
		t.Errorf("receipt sentence does not state the total:\n%s", line)
	}
	t.Logf("history.page 回执逐字：\n%s", reply)
}

// TestHistoryPageScrollsOneThousandRowsWithoutOverlap is AC#3's 1000-row
// support, measured as paging fact rather than an adjective: 1005 seeded
// rows, three consecutive windows, no row twice, no row skipped, every
// window still newest-first.
func TestHistoryPageScrollsOneThousandRowsWithoutOverlap(t *testing.T) {
	dir := t.TempDir()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	const seed = 1005
	for i := 0; i < seed; i++ {
		seedHistoryTask(t, st, fmt.Sprintf("scroll-%04d", i), int64(10000+i), string(memory.TaskStateDone), "")
	}
	_ = st.Close()

	h := newHistoryHandlerWave2(dir, nil, "test")
	ctx := context.Background()
	seen := map[string]bool{}
	var firstOfPage []string
	for off := 0; off < seed; off += 500 {
		req := panel.ComposerRequest{
			Method: panel.MethodHistoryPage, RequestID: fmt.Sprintf("rw2-s%d", off),
			Source: panel.ComposerRequestSource, HistoryOffset: off, HistoryLimit: 500,
		}
		_, payload, ok := strings.Cut(mustHistoryPage(t, h, ctx, req), "\n"+historyPagePayloadMarker)
		if !ok {
			t.Fatalf("page at offset %d carries no payload line", off)
		}
		var page memory.HistoryPageTotal
		if err := json.Unmarshal([]byte(payload), &page); err != nil {
			t.Fatalf("payload at offset %d unreadable: %v", off, err)
		}
		if page.Total != seed {
			t.Errorf("total at offset %d = %d, want %d", off, page.Total, seed)
		}
		if len(page.Rows) == 0 {
			t.Fatalf("offset %d answered an empty page; scrolling would stop short of the real end", off)
		}
		firstOfPage = append(firstOfPage, page.Rows[0].ID)
		for _, r := range page.Rows {
			if seen[r.ID] {
				t.Errorf("row %s answered twice across the scroll: paging is not a partition", r.ID)
			}
			seen[r.ID] = true
		}
	}
	if len(seen) != seed {
		t.Errorf("the scroll covered %d distinct rows, want %d - %d rows were never served", len(seen), seed, seed-len(seen))
	}
	// Rank math, spelled out: 1005 rows newest-first puts scroll-1004 at rank
	// 1, so the windows at offsets 0/500/1000 lead at ranks 1/501/1001 =
	// scroll-1004 / scroll-0504 / scroll-0004.
	if firstOfPage[0] != "scroll-1004" || firstOfPage[1] != "scroll-0504" || firstOfPage[2] != "scroll-0004" {
		t.Errorf("page leads drifted off the newest-first order: %v", firstOfPage)
	}
}

// TestHistoryPageRefusesAnUnservableWindowByName: negative offset and a page
// over the ceiling are refused before the store is opened, with the wrong
// argument named.
func TestHistoryPageRefusesAnUnservableWindowByName(t *testing.T) {
	dir := t.TempDir()
	h := newHistoryHandlerWave2(dir, nil, "test")
	ctx := context.Background()
	for _, tc := range []struct {
		name       string
		off, limit int
		wantInErr  string
	}{
		{"negative offset", -1, 10, "Offset"},
		{"zero limit", 0, -3, "limit"},
		{"over ceiling", 0, memory.HistoryPageLimitCeiling + 1, "上限"},
	} {
		req := panel.ComposerRequest{
			Method: panel.MethodHistoryPage, RequestID: "rw2-bad", Source: panel.ComposerRequestSource,
			HistoryOffset: tc.off, HistoryLimit: tc.limit,
		}
		_, err := h.HandleHistoryPage(ctx, req)
		if err == nil {
			t.Errorf("%s: the door answered an unservable window", tc.name)
			continue
		}
		if !strings.Contains(err.Error(), tc.wantInErr) {
			t.Errorf("%s: refusal %v does not name the wrong argument (%q)", tc.name, err, tc.wantInErr)
		}
	}
}

// TestHistoryDoorIsWiredIntoTheProductionChain is the wiring nail: unplug the
// History socket in newComposerDispatchChain and this reddens - the door a
// page can call must not ship unattached (the 194 precedent's nail shape).
func TestHistoryDoorIsWiredIntoTheProductionChain(t *testing.T) {
	disp, _, err := newComposerDispatchChain(newInboundDataDir33(t, "ask_every_step"), auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("chain assembly: %v", err)
	}
	if disp.History == nil {
		t.Fatal("newComposerDispatchChain left the History socket nil: history.page is on the whitelist, so a machine with no leg must at least refuse by name, not fall to rosterMismatch")
	}
	req := panel.ComposerRequest{
		Method: panel.MethodHistoryPage, RequestID: "rw2-wire", Source: panel.ComposerRequestSource,
	}
	if _, err := disp.Handle(context.Background(), marshalComposerForTest(req)); err != nil {
		t.Errorf("the wired history door refused its own parsed request: %v", err)
	}
}

func mustHistoryPage(t *testing.T, h *historyHandlerWave2, ctx context.Context, req panel.ComposerRequest) string {
	t.Helper()
	reply, err := h.HandleHistoryPage(ctx, req)
	if err != nil {
		t.Fatalf("history.page refused (offset %d): %v", req.HistoryOffset, err)
	}
	return reply
}

// marshalComposerForTest renders the typed request as the page's JSON line.
func marshalComposerForTest(req panel.ComposerRequest) string {
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return string(b)
}
