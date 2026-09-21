package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The C4 tool surface as the agent loop consumes it. Real tools land with
// ticket 20+ (host bridge + builtin slice); this ticket owns only the
// consumer-side interface, the loop's dispatch policy (concurrency ceiling,
// per-tool timeout, result spill) and a trivial echo provider used by the
// loop tests.
//
// Gating is NOT implemented here (ticket 21): a provider that reports no risk
// level ("", i.e. unclassified) is routed as L0 pass-through while
// Config.PassThroughUnclassifiedRisk is true. With the flag off, an
// unclassified call fails closed instead of executing.
//
// Spill is a HOST-INTERNAL artifact write (D15(3) + D34 note(2)): it is
// deliberately not modelled as a gated tool.

// Risk levels re-exported from the storage enum so the loop and the
// tool_call rows cannot drift apart.
const (
	RiskL0 = memory.RiskL0
	RiskL1 = memory.RiskL1
	RiskL2 = memory.RiskL2
	// RiskUnclassified is what a provider returns before ticket 21 wires the
	// C19 assessor into the host bridge.
	RiskUnclassified = ""
)

// ToolInfo is one directory entry (name, description, JSON Schema, residency).
// Resident tools are builtin: they are injected into every request (D15(1));
// non-resident (third-party/long-tail) tools are injected only when the local
// BM25/keyword retrieval selects them (D15(2)).
type ToolInfo struct {
	Name        string
	Description string
	Parameters  json.RawMessage
	Resident    bool
	RiskLevel   string
}

// ToolRequest is one dispatch. Timeout is the resolved per-tool budget
// (C22); Args is the assembled argument object.
type ToolRequest struct {
	TaskID        string
	CorrelationID string
	CallID        string
	Name          string
	Args          json.RawMessage
	Timeout       time.Duration
}

// ToolOutcome is a tool's answer. Text is the (already size-capped) result;
// the loop applies spill on top of it. ErrorClass is a D37 class and is only
// meaningful when IsError is true: internal tool failures go back to the LLM
// as class "tool" for self-correction (D37), they never surface as task
// failures.
type ToolOutcome struct {
	Text       string
	IsError    bool
	RiskLevel  string
	ErrorClass string
	Truncated  bool
	// AppliedSteps is the D31 ledger the tool reported for itself (see
	// tools.Result.AppliedSteps): what already landed when the call stopped
	// early. It is DATA, not a second report type - the report wording stays
	// approval.CancellationReport's - because a host that shows the strip, the
	// panel or the transcript needs the list, not a re-parse of Text.
	AppliedSteps []string
}

// ToolProvider is the loop's view of the tool host (C4).
type ToolProvider interface {
	// Tools returns the full directory every turn (the loop does the
	// selection locally, D15(2): zero extra LLM round-trips).
	Tools(ctx context.Context) ([]ToolInfo, error)

	// Execute runs one call and must honor ctx cancellation and the request
	// timeout. The returned error is a host failure (class internal/tool);
	// a tool's own failure is reported in ToolOutcome.
	Execute(ctx context.Context, req ToolRequest) (ToolOutcome, error)
}

// ---------------------------------------------------------------------------
// Echo provider: the trivial ToolProvider used by the loop tests and by
// `wisp run` before ticket 20 lands real tools.

// EchoCall records one executed call (test/assertion surface).
type EchoCall struct {
	Req    ToolRequest
	Out    ToolOutcome
	Offset time.Duration // monotonic since the provider was created
}

// EchoProvider is a trivial ToolProvider: it echoes its arguments, can sleep
// on request (per-tool timeout tests), counts calls and tracks observed
// concurrency. It is safe for concurrent use.
type EchoProvider struct {
	tools []ToolInfo

	mu    sync.Mutex
	calls []EchoCall
	// inflight/outMax drive the concurrency-ceiling assertion (D38d).
	inflight atomic.Int32
	obsMax   atomic.Int32
	created  time.Time
}

// NewEchoProvider builds a provider over the given directory; with no
// arguments it exposes the two default dev tools (echo, sleep).
func NewEchoProvider(tools ...ToolInfo) *EchoProvider {
	if len(tools) == 0 {
		tools = []ToolInfo{
			{
				Name: "echo", Description: "Echo back the given text.", Resident: true,
				Parameters: json.RawMessage(`{"type":"object","properties":{"text":{"type":"string"}},"required":["text"]}`),
			},
			{
				Name: "sleep", Description: "Sleep for the given number of milliseconds.", Resident: true,
				Parameters: json.RawMessage(`{"type":"object","properties":{"ms":{"type":"integer"}},"required":["ms"]}`),
			},
		}
	}
	return &EchoProvider{tools: tools, created: time.Now()}
}

// Tools implements ToolProvider.
func (p *EchoProvider) Tools(context.Context) ([]ToolInfo, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]ToolInfo, len(p.tools))
	copy(out, p.tools)
	return out, nil
}

type echoArgs struct {
	Text string `json:"text"`
	MS   int    `json:"ms"`
}

// Execute implements ToolProvider.
func (p *EchoProvider) Execute(ctx context.Context, req ToolRequest) (ToolOutcome, error) {
	cur := p.inflight.Add(1)
	for {
		prev := p.obsMax.Load()
		if cur <= prev || p.obsMax.CompareAndSwap(prev, cur) {
			break
		}
	}
	defer p.inflight.Add(-1)

	var a echoArgs
	out := ToolOutcome{}
	if len(req.Args) > 0 {
		if err := json.Unmarshal(req.Args, &a); err != nil {
			out = ToolOutcome{
				Text: "invalid arguments: " + err.Error(), IsError: true,
				ErrorClass: string(observe.ClassTool),
			}
			p.record(req, out)
			return out, nil
		}
	}

	switch req.Name {
	case "echo":
		out = ToolOutcome{Text: a.Text}
	case "sleep":
		d := time.Duration(a.MS) * time.Millisecond
		tm := observe.NewTimeout(d)
		select {
		case <-ctx.Done():
			out = ToolOutcome{Text: "cancelled", IsError: true, ErrorClass: string(observe.ClassCancelled)}
		case <-time.After(tm.Remaining()):
			out = ToolOutcome{Text: fmt.Sprintf("slept %dms", a.MS)}
		}
	default:
		out = ToolOutcome{
			Text: "unknown tool " + req.Name, IsError: true,
			ErrorClass: string(observe.ClassTool),
		}
	}
	p.record(req, out)
	return out, nil
}

// record appends one completed call to the assertion log.
func (p *EchoProvider) record(req ToolRequest, out ToolOutcome) {
	p.mu.Lock()
	p.calls = append(p.calls, EchoCall{Req: req, Out: out, Offset: time.Since(p.created)})
	p.mu.Unlock()
}

// Calls returns the recorded calls in execution order.
func (p *EchoProvider) Calls() []EchoCall {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]EchoCall, len(p.calls))
	copy(out, p.calls)
	return out
}

// CallCount returns how many calls were executed.
func (p *EchoProvider) CallCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.calls)
}

// MaxConcurrent returns the highest simultaneously observed execution count.
func (p *EchoProvider) MaxConcurrent() int32 { return p.obsMax.Load() }
