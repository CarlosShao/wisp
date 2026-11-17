package main

// Readface wave 2 (2026-11-17), ticket 38 AC#4's wire reading: the three
// task-card fields the roster reader fills reach the packet's tasks section
// verbatim, and a task the ledger never saw keeps all three ABSENT (omitempty
// on the wire) - "宿主没有这一维的读数", not epoch/zero/empty.

import (
	"encoding/json"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/tools"
)

// taskLiveSectionBytes renders the roster section the pump would wire, off a
// hand-built runtime (the same readers the production pump composes: the
// roster, the stream log, the live ledger).
func taskLiveSectionBytes(t *testing.T, rt *agentRuntime) []byte {
	t.Helper()
	state := rt.taskRosterState()
	sect := panel.TaskRosterSectionFrom(state, nil)
	data, err := json.Marshal(sect)
	if err != nil {
		t.Fatalf("marshal roster section: %v", err)
	}
	return data
}

func taskLiveRow(t *testing.T, data []byte, id string) map[string]json.RawMessage {
	t.Helper()
	var sect struct {
		Rows []map[string]json.RawMessage `json:"rows"`
	}
	if err := json.Unmarshal(data, &sect); err != nil {
		t.Fatalf("tasks section unreadable: %v (%s)", err, data)
	}
	for _, r := range sect.Rows {
		var taskID string
		_ = json.Unmarshal(r["taskId"], &taskID)
		if taskID == id {
			return r
		}
	}
	t.Fatalf("no roster row for %q in %s", id, data)
	return nil
}

func TestTaskCardFieldsRideTheRosterWire(t *testing.T) {
	rt := &agentRuntime{
		tasks:      tools.NewTaskRoster(),
		taskLive:   newTaskLiveLedger(),
		stream:     newPanelStreamLogForWave2(),
		rootCancel: &currentRootCancel{},
	}
	// The values a live task carries: the admission stamp (noteTask is the
	// production stamping hook), one EvUsage cumulative pair, one in-flight
	// tool call. The reader enumerates seenTasks, so the id goes through the
	// same admission hook the loop uses.
	const id = "card-w2-1"
	rt.noteTask(id)
	rt.tasks.Record(id, tools.TaskOutput{
		Text: "in flight", Label: "票38 卡片三字段正控", Kind: tools.TaskKindRoot,
	})
	rt.taskLive.usageRecorded(id, 150, 40, 20, 26000, "CNY")
	rt.taskLive.stepStarted(id, "fs.read")

	data := taskLiveSectionBytes(t, rt)
	row := taskLiveRow(t, data, id)
	var got struct {
		StartedAt   int64           `json:"startedAt"`
		CostToDate  json.RawMessage `json:"costToDate"`
		CurrentStep string          `json:"currentStep"`
	}
	for _, key := range []string{"startedAt", "costToDate", "currentStep"} {
		if _, ok := row[key]; !ok {
			t.Fatalf("wire row carries no %s key: %s", key, data)
		}
	}
	if err := json.Unmarshal(rowBytes(row), &got); err != nil {
		t.Fatalf("card fields unreadable: %v", err)
	}
	if got.StartedAt == 0 {
		t.Errorf("startedAt = 0: noteTask's admission stamp never reached the wire")
	}
	var cost panel.RunUsageRow
	if err := json.Unmarshal(got.CostToDate, &cost); err != nil {
		t.Fatalf("costToDate unreadable: %v (%s)", err, got.CostToDate)
	}
	if cost.Micros != 26000 || cost.TokensIn != 150 || cost.TokensOut != 40 || cost.Currency != "CNY" {
		t.Errorf("costToDate = %+v, want the EvUsage pair verbatim", cost)
	}
	if got.CurrentStep != "fs.read" {
		t.Errorf("currentStep = %q, want the in-flight tool name", got.CurrentStep)
	}
}

func TestTaskCardFieldsStayAbsentWhenTheHostStampedNothing(t *testing.T) {
	rt := &agentRuntime{
		tasks:      tools.NewTaskRoster(),
		taskLive:   newTaskLiveLedger(),
		stream:     newPanelStreamLogForWave2(),
		rootCancel: &currentRootCancel{},
	}
	const id = "card-w2-2"
	// Seen by the reader WITHOUT a ledger stamp: the direct seenTasks write is
	// the shape of a row filed by a pre-wave host (the ledger exists, the
	// stamping hook never ran for this id).
	rt.taskMu.Lock()
	rt.seenTasks = append(rt.seenTasks, id)
	rt.taskMu.Unlock()
	rt.tasks.Record(id, tools.TaskOutput{
		Text: "no ledger saw this", Label: "票38 缺读数正控", Kind: tools.TaskKindRoot,
	})
	data := taskLiveSectionBytes(t, rt)
	row := taskLiveRow(t, data, id)
	for _, key := range []string{"startedAt", "costToDate", "currentStep"} {
		if _, ok := row[key]; ok {
			t.Errorf("wire row carries %s for a task the ledger never saw (%s): absence is the honest shape, not epoch/zero/empty", key, data)
		}
	}
}

func rowBytes(row map[string]json.RawMessage) []byte {
	b, err := json.Marshal(row)
	if err != nil {
		panic(err)
	}
	return b
}

// newPanelStreamLogForWave2 builds the stream log the roster reader reads
// truncation facts from; the card fields never touch it, but the reader
// dereferences it unconditionally.
func newPanelStreamLogForWave2() *panel.StreamLog {
	return panel.NewStreamLog(panel.DefaultStreamKeys)
}
