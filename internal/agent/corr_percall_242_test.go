package agent

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

// corrWaitProvider holds every tool call until the test releases THAT call by
// the correlation id the loop minted for it. It is the thin fixture ticket 242
// needs: with one held call per correlation id, "the corr addresses exactly one
// in-flight ask" is observable without a database or an approval card.
type corrWaitProvider struct {
	mu     sync.Mutex
	hold   map[string]chan struct{} // corr -> release channel
	byCall map[string]string        // call id -> corr
	start  chan struct{}            // one token per call that entered Execute
	finish chan string              // call ids, in the order they returned
}

func newCorrWaitProvider() *corrWaitProvider {
	return &corrWaitProvider{
		hold:   map[string]chan struct{}{},
		byCall: map[string]string{},
		start:  make(chan struct{}, 8),
		finish: make(chan string, 8),
	}
}

// Tools implements ToolProvider with the L0 directory the parallel-tools
// fixture calls (echo), mirroring blockingProvider's entries.
func (p *corrWaitProvider) Tools(context.Context) ([]ToolInfo, error) {
	return []ToolInfo{{
		Name: "echo", Description: "Echo back the given text.", Resident: true,
		RiskLevel: RiskL0, Parameters: json.RawMessage(`{"type":"object"}`),
	}}, nil
}

// Execute holds until its own correlation id is released (or the ctx dies).
func (p *corrWaitProvider) Execute(ctx context.Context, req ToolRequest) (ToolOutcome, error) {
	ch := make(chan struct{})
	p.mu.Lock()
	p.hold[req.CorrelationID] = ch
	p.byCall[req.CallID] = req.CorrelationID
	p.mu.Unlock()
	p.start <- struct{}{}

	select {
	case <-ch:
	case <-ctx.Done():
		return ToolOutcome{}, ctx.Err()
	}
	p.finish <- req.CallID
	return ToolOutcome{Text: "answered"}, nil
}

// release lifts exactly the call its corr names.
func (p *corrWaitProvider) release(corr string) {
	p.mu.Lock()
	ch := p.hold[corr]
	p.mu.Unlock()
	if ch != nil {
		close(ch)
	}
}

// corr reads the corr the loop minted for one call id.
func (p *corrWaitProvider) corr(callID string) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.byCall[callID]
}

// TestCorrPerCallTwoAsksSameTask242 pins ticket 242's ruling: the loop mints an
// independent correlation id for EVERY tool call, so two concurrent asks of one
// task carry two ids - each id routes back to exactly its own call and to the
// task (C18 keys the approval reply by correlationId; the task-level trace
// stays taskID).
func TestCorrPerCallTwoAsksSameTask242(t *testing.T) {
	p := newCorrWaitProvider()
	h := newHarness(t, "parallel-tools", withTools(p),
		withConfig(func(c *Config) { c.PerToolTimeout = 30 * time.Second }))

	done := make(chan Result, 1)
	go func() { done <- h.loop.Run(context.Background(), "两件事一起办") }()

	// Both calls of the one turn must be in flight together: the collapse
	// ticket 242 guards against needs the asks to be concurrent.
	for i := 0; i < 2; i++ {
		select {
		case <-p.start:
		case <-time.After(10 * time.Second):
			t.Fatal("the second call never entered the tool: the two asks are not concurrent")
		}
	}
	c1, c2 := p.corr("call_p1"), p.corr("call_p2")
	if c1 == "" || c2 == "" {
		t.Fatalf("correlation ids = %q/%q, want both non-empty (C18 routes the reply by them)", c1, c2)
	}
	if c1 == c2 {
		t.Fatalf("both concurrent calls of one task carry the same corr %q: the two asks collapsed onto one routing key", c1)
	}

	// Answer the first ask by ITS corr: exactly call_p1 must wake, and the
	// task cannot finish while call_p2 is still held.
	p.release(c1)
	select {
	case callID := <-p.finish:
		if callID != "call_p1" {
			t.Fatalf("releasing corr %q woke %s, want call_p1", c1, callID)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("the call addressed by corr %q never returned", c1)
	}
	select {
	case <-done:
		t.Fatal("the task completed although the second ask was never answered: the two asks share one wake-up")
	default:
	}

	// Answer the second ask by its own corr; now the task may finish.
	p.release(c2)
	select {
	case callID := <-p.finish:
		if callID != "call_p2" {
			t.Fatalf("releasing corr %q woke %s, want call_p2", c2, callID)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("the call addressed by corr %q never returned", c2)
	}
	res := <-done
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	// Each corr routes back to the task that asked: the task id is the prefix.
	if !strings.HasPrefix(c1, res.TaskID) || !strings.HasPrefix(c2, res.TaskID) {
		t.Errorf("corrs %q/%q do not route back to task %q", c1, c2, res.TaskID)
	}
	if c1 == res.TaskID || c2 == res.TaskID {
		t.Errorf("corrs %q/%q still collapse onto the task id", c1, c2)
	}
}
