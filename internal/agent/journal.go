package agent

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/CarlosShao/wisp/internal/memory"
)

// The 04-tables writer. task_log (L3) and tool_call rows are written through
// this interface, which *memory.Store satisfies directly, so the loop never
// re-implements storage and never invents column values: the enum-shaped
// fields (risk_level / decision / outcome / error_class) are validated by the
// DAO and the loop uses memory's exported constants.

// Journal is the subset of the memory store the loop writes through.
type Journal interface {
	StartTaskLog(ctx context.Context, tl memory.TaskLog) error
	FinishTaskLog(ctx context.Context, id, state, summary string, tokensIn, tokensOut, costMicro int64, errorClass string) error
	InsertToolCall(ctx context.Context, tc memory.ToolCall) (int64, error)
	DecideToolCall(ctx context.Context, id int64, decision string, grantID *int64) error
	FinishToolCall(ctx context.Context, id int64, outcome, errorClass string) error
}

var (
	_ Journal = (*memory.Store)(nil)
)

// Decision values (tool_call.decision). Ticket 21 replaces the pass-through
// decision with real gate outcomes; the vocabulary is already frozen here.
const (
	DecisionAllow      = "allow"
	DecisionAllowGrant = "allow_session_grant"
	DecisionReject     = "reject"
	DecisionTimeout    = "timeout"
	DecisionBatchAggr  = "batch_aggregated"
)

// Outcome values (tool_call.outcome).
const (
	OutcomeSuccess   = "success"
	OutcomeError     = "error"
	OutcomeCancelled = "cancelled"
	OutcomeTruncated = "truncated"
)

// taskJournal books one task's rows: monotonic seq, correlation id, and the
// store call sites. A nil Journal disables all writes (unit tests of the pure
// loop run without a database).
type taskJournal struct {
	mu     sync.Mutex
	j      Journal
	taskID string
	corrID string
	seq    int64
	// rows maps tool_call id -> DB row id for finish/update calls.
	rows map[string]int64
}

func newTaskJournal(j Journal, taskID, corrID string) *taskJournal {
	if j == nil {
		return &taskJournal{taskID: taskID, corrID: corrID, rows: map[string]int64{}}
	}
	return &taskJournal{j: j, taskID: taskID, corrID: corrID, rows: map[string]int64{}}
}

// startCall books a pending tool_call row (decision recorded separately once
// the gate/pass-through resolves it).
func (t *taskJournal) startCall(ctx context.Context, callID, tool string, args json.RawMessage, risk string) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.j == nil {
		return 0
	}
	t.seq++
	rowID, err := t.j.InsertToolCall(ctx, memory.ToolCall{
		TaskID:        t.taskID,
		Seq:           t.seq,
		Tool:          tool,
		ArgsJSON:      string(args),
		RiskLevel:     riskColumn(risk),
		CorrelationID: t.corrID,
	})
	if err != nil {
		return 0
	}
	if t.rows == nil {
		t.rows = map[string]int64{}
	}
	t.rows[callID] = rowID
	return rowID
}

// decide books the gate decision for a row.
func (t *taskJournal) decide(ctx context.Context, rowID int64, decision string) {
	t.mu.Lock()
	j := t.j
	t.mu.Unlock()
	if j == nil || rowID == 0 {
		return
	}
	_ = j.DecideToolCall(ctx, rowID, decision, nil)
}

// finish books the outcome of a row.
func (t *taskJournal) finish(ctx context.Context, callID string, rowID int64, outcome, errorClass string) {
	t.mu.Lock()
	j := t.j
	if rowID == 0 && t.rows != nil {
		rowID = t.rows[callID]
	}
	t.mu.Unlock()
	if j == nil || rowID == 0 {
		return
	}
	if outcome == "" {
		outcome = OutcomeError
	}
	_ = j.FinishToolCall(ctx, rowID, outcome, errorClass)
}

// startTask opens the task_log row.
func (t *taskJournal) startTask(ctx context.Context, query string) error {
	t.mu.Lock()
	j := t.j
	t.mu.Unlock()
	if j == nil {
		return nil
	}
	return j.StartTaskLog(ctx, memory.TaskLog{
		ID:        t.taskID,
		State:     "running",
		QueryText: query,
	})
}

// finishTask closes the task_log row with the C23 numbers.
func (t *taskJournal) finishTask(ctx context.Context, state, summary string, tokensIn, tokensOut, costMicro int64, errorClass string) error {
	t.mu.Lock()
	j := t.j
	t.mu.Unlock()
	if j == nil {
		return nil
	}
	return j.FinishTaskLog(ctx, t.taskID, state, summary, tokensIn, tokensOut, costMicro, errorClass)
}

// riskColumn maps the pre-gate risk onto the NOT NULL risk_level column:
// declared L1/L2 stay as declared, an unclassified call is booked as L0 while
// the pass-through flag is on (ticket 21 writes the assessed level instead).
func riskColumn(risk string) string {
	switch risk {
	case memory.RiskL1, memory.RiskL2:
		return risk
	default:
		return memory.RiskL0
	}
}
