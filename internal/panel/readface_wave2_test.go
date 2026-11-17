package panel

// Readface wave 2 (2026-11-17) - the two new doors' internal/panel teeth.
//
// Ticket 36 AC#3's history read (history.page) and ticket 38 AC#4's cancel
// door (roster.cancel) join the C17 whitelist this wave. The teeth here are
// the package's own: the envelope parses the two names, the router reaches
// the attached handler and refuses an unattached one by name, and the
// roster-cancel spelling is measured against the one instrument that shaped
// its name. The legs themselves (store paging, real cancel semantics) are
// judged on the cmd/wisp side; the closed-roster and page-gap rulers moved in
// the same change (inbound_roster_253_test.go's 11 -> 13, bridge_test.go's
// rosterPageGapLedger194).

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func readfaceWave2Request(method, field, value string) ComposerRequest {
	req := ComposerRequest{Method: method, RequestID: "rw2-1", Source: ComposerRequestSource}
	switch field {
	case "rosterTaskId":
		req.RosterTaskID = value
	}
	return req
}

func readfaceWave2Raw(req ComposerRequest) string {
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// TestReadfaceWave2DoorsAreOnTheWhitelistAndRouted drives both doors through
// the production envelope parser and the production router.
func TestReadfaceWave2DoorsAreOnTheWhitelistAndRouted(t *testing.T) {
	history := &recordingHistoryHandlerWave2{}
	cancel := &recordingRosterCancelHandlerWave2{}
	d := &ComposerDispatch{History: history, RosterCancel: cancel}

	// history.page with a paging window: parsed, routed, receipt passthrough.
	raw := `{"method":"history.page","requestId":"rw2-h1","source":"panel-composer","historyOffset":20,"historyLimit":50}`
	reply, err := d.Handle(context.Background(), raw)
	if err != nil {
		t.Fatalf("history.page was refused by the production hop: %v", err)
	}
	if reply != history.receipt {
		t.Errorf("history.page reply = %q, want the handler's own receipt %q", reply, history.receipt)
	}
	if !history.called || history.last.Offset != 20 || history.last.Limit != 50 {
		t.Errorf("history page handler saw %+v (called=%v): the paging window did not ride the envelope", history.last, history.called)
	}

	// roster.cancel with its row selector: parsed, routed, receipt passthrough.
	raw = `{"method":"roster.cancel","requestId":"rw2-c1","source":"panel-composer","rosterTaskId":"child-42"}`
	reply, err = d.Handle(context.Background(), raw)
	if err != nil {
		t.Fatalf("roster.cancel was refused by the production hop: %v", err)
	}
	if reply != cancel.receipt {
		t.Errorf("roster.cancel reply = %q, want the handler's own receipt %q", reply, cancel.receipt)
	}
	if !cancel.called || cancel.lastTaskID != "child-42" {
		t.Errorf("cancel handler saw called=%v taskId=%q: the row selector did not ride the envelope", cancel.called, cancel.lastTaskID)
	}
}

// TestReadfaceWave2UnattachedDoorsRefuseByName is the nil-socket half: a
// whitelisted name on a machine with no leg attached refuses by name - never
// a silence, never a fake success.
func TestReadfaceWave2UnattachedDoorsRefuseByName(t *testing.T) {
	d := &ComposerDispatch{}
	for _, method := range []string{MethodHistoryPage, MethodRosterCancel} {
		raw := readfaceWave2Raw(readfaceWave2Request(method, "", ""))
		_, err := d.Handle(context.Background(), raw)
		if !errors.Is(err, ErrNoHandlerAttached) {
			t.Errorf("%s with no handler attached returned %v, want ErrNoHandlerAttached", method, err)
		}
	}
}

// TestReadfaceWave2NamesStayClearOfTheTaskAxisInstrument is the naming ruling
// measured, not asserted in prose: internal/tools' frozen ticket-188 scan
// reddens every panel-side Method* constant whose VALUE contains "task" or
// "state" (task_state_188_test.go:385-401). That test is outside this wave's
// write face and its assertion is not loosened; the two names land clear of
// it BY VALUE, and this case pins the fact so a later rename re-reads the
// instrument instead of rediscovering it red.
func TestReadfaceWave2NamesStayClearOfTheTaskAxisInstrument(t *testing.T) {
	for _, c := range []struct {
		name  string
		value string
	}{
		{"MethodHistoryPage", MethodHistoryPage},
		{"MethodRosterCancel", MethodRosterCancel},
	} {
		lower := strings.ToLower(c.value)
		if strings.Contains(lower, "task") || strings.Contains(lower, "state") {
			t.Errorf("%s = %q: the value claims the task axis, and internal/tools' ticket-188 scan (task_state_188_test.go:395) reddens exactly that shape - rename the door or re-read the instrument's scope before landing", c.name, c.value)
		}
	}
}

type recordingHistoryHandlerWave2 struct {
	called  bool
	last    HistoryWindowWave2
	receipt string
}

type HistoryWindowWave2 struct {
	Offset int
	Limit  int
}

func (h *recordingHistoryHandlerWave2) HandleHistoryPage(ctx context.Context, req ComposerRequest) (string, error) {
	h.called = true
	h.last = HistoryWindowWave2{Offset: req.HistoryOffset, Limit: req.HistoryLimit}
	h.receipt = "成本按天共 3 行。" // 形状占位由 cmd/wisp 的真腿判；这里只证路由与封套
	return h.receipt, nil
}

type recordingRosterCancelHandlerWave2 struct {
	called      bool
	lastTaskID  string
	lastRequest ComposerRequest
	receipt     string
}

func (h *recordingRosterCancelHandlerWave2) HandleRosterCancel(ctx context.Context, req ComposerRequest) (string, error) {
	h.called = true
	h.lastTaskID = req.RosterTaskID
	h.lastRequest = req
	h.receipt = "已对任务 child-42 发出停止"
	return h.receipt, nil
}
