package agent

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Shared test fixture: a ToolProvider that HOLDS every call until the test
// releases it. It exists because the behaviours the adversarial review demanded
// proof of are unobservable with the instant-returning EchoProvider:
//
//   - real tool concurrency (D38d): EchoProvider returns so fast that
//     MaxConcurrent() stays 1, so a `<= 4` ceiling assertion passes even for a
//     fully serial implementation.
//   - a cancelled task's in-flight tool_call rows: the calls have to still be
//     running when the cancellation lands.
//
// It is also CONTRACT-HONEST about host failure, which EchoProvider is not
// (tools.go returns nil error on ctx.Done, violating its own documented
// contract at tools.go:82-85): when the ctx dies this provider returns
// ctx.Err() alongside a zero outcome, which is exactly the shape a real host
// bridge (ticket 20+) will produce, and exactly the shape that used to swallow
// the structured TOOL_TIMEOUT (MAJOR-2).

// sleepArgs is the argument shape the golden fixtures use for the sleep tool.
type sleepArgs struct {
	Text string `json:"text"`
	MS   int    `json:"ms"`
}

// blockingProvider holds tool calls until released. Safe for concurrent use.
type blockingProvider struct {
	mu sync.Mutex

	hold     chan struct{} // closed by release() to let held calls return
	released bool

	started    atomic.Int32 // calls that entered Execute
	inflight   atomic.Int32
	obsMax     atomic.Int32
	callCount  atomic.Int32
	errOnAbort bool // return ctx.Err() (host failure) instead of a soft outcome

	done []string // recorded "<name>:<abort>" lines, in completion order
}

// newBlockingProvider builds a provider whose calls block until release().
func newBlockingProvider() *blockingProvider {
	return &blockingProvider{hold: make(chan struct{}), errOnAbort: true}
}

// Tools implements ToolProvider with an L0 directory covering the tool names
// the golden fixtures call (echo, sleep).
func (p *blockingProvider) Tools(context.Context) ([]ToolInfo, error) {
	return []ToolInfo{
		{
			Name: "echo", Description: "Echo back the given text.", Resident: true, RiskLevel: RiskL0,
			Parameters: json.RawMessage(`{"type":"object"}`),
		},
		{
			Name: "sleep", Description: "Sleep for the given number of milliseconds.",
			Resident: true, RiskLevel: RiskL0,
			Parameters: json.RawMessage(`{"type":"object"}`),
		},
	}, nil
}

// Execute implements ToolProvider: it announces its own start, holds until the
// test releases it or the ctx dies, and reports a ctx death as a host error
// (the documented contract).
func (p *blockingProvider) Execute(ctx context.Context, req ToolRequest) (ToolOutcome, error) {
	p.started.Add(1)
	cur := p.inflight.Add(1)
	for {
		prev := p.obsMax.Load()
		if cur <= prev || p.obsMax.CompareAndSwap(prev, cur) {
			break
		}
	}
	defer p.inflight.Add(-1)
	defer p.callCount.Add(1)

	select {
	case <-p.hold:
	case <-ctx.Done():
	}
	// Cancellation is authoritative over the release wake-up: a select with both
	// cases ready picks at random, and a real host call whose ctx died reports
	// the death rather than the answer it happened to finish writing.
	if cerr := ctx.Err(); cerr != nil && p.errOnAbort {
		return ToolOutcome{}, cerr
	}
	if cerr := ctx.Err(); cerr != nil {
		return ToolOutcome{
			Text: "cancelled", IsError: true,
			ErrorClass: string(observe.ClassCancelled),
		}, nil
	}
	return ToolOutcome{Text: "released"}, nil
}

// release unblocks every held call.
func (p *blockingProvider) release() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.released {
		p.released = true
		close(p.hold)
	}
}

// Started is how many calls entered Execute.
func (p *blockingProvider) Started() int32 { return p.started.Load() }

// CallCount is how many calls returned from Execute.
func (p *blockingProvider) CallCount() int32 { return p.callCount.Load() }

// MaxConcurrent is the highest simultaneously observed execution count.
func (p *blockingProvider) MaxConcurrent() int32 { return p.obsMax.Load() }

// waitForInflight blocks until n calls are in flight at once (i.e. the calls
// are really running concurrently, not queued). The bound is a monotonic
// observe.Timeout budget, never a wall-clock difference (D42#9).
func (p *blockingProvider) waitForInflight(n int32, within time.Duration) bool {
	deadline := observe.NewTimeout(within)
	for !deadline.Expired() {
		if p.inflight.Load() >= n {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return p.inflight.Load() >= n
}
