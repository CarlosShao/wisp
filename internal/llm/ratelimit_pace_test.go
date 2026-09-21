package llm_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat" // protocol registration
	"github.com/CarlosShao/wisp/internal/observe"
)

// AC#7: per-provider rpm/tpm self-restraint. With rpm=10 the ELEVENTH request
// must wait LOCALLY - the fixture server never sees it, so no 429 is ever
// exchanged (that is the whole point: "behave before the provider 429s us").
//
// The proof is the SERVER'S OWN received-request count, never elapsed wall
// time: every local-wait step re-reads the server counter and requires it to
// still be 10. Pacing arithmetic runs on an injected monotonic source
// (D42#9); the 250ms re-check granularity of tokenBucket.take (there for
// context cancellation) is asserted, not hidden.

// pacedServer counts requests; its own >10 rule 429s anything the client
// failed to hold back (a live tripwire, not decoration: if Wisp ever sends an
// 11th request in-window, status429 goes up and the test fails).
type pacedServer struct {
	srv *httptest.Server

	mu        sync.Mutex
	hits      int
	status429 int
}

func newPacedServer(t *testing.T) *pacedServer {
	t.Helper()
	ps := &pacedServer{}
	ps.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		ps.mu.Lock()
		ps.hits++
		hits := ps.hits
		if hits > 10 {
			ps.status429++
		}
		ps.mu.Unlock()
		if hits > 10 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"too many","type":"rate_limit_error"}}`))
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		body := fmt.Sprintf("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"echo %d\"},\"finish_reason\":null}]}\n\n", hits) +
			"data: {\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n" +
			"data: {\"choices\":[],\"usage\":{\"prompt_tokens\":5,\"completion_tokens\":2,\"total_tokens\":7}}\n\n" +
			"data: [DONE]\n\n"
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(ps.srv.Close)
	return ps
}

func (ps *pacedServer) hitsCount() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.hits
}

func (ps *pacedServer) tooMany() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.status429
}

// fakeMono is a MONOTONIC time source a test advances by hand (D42#9: the
// bucket may only ever difference readings taken from its own source).
type fakeMono struct {
	mu sync.Mutex
	t  time.Time
}

func (f *fakeMono) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.t
}

// Advance moves the fake clock forward (or backward, for the clock-jump case).
func (f *fakeMono) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.t = f.t.Add(d)
}

func chatProviderAt(t *testing.T, base string) llm.LlmProvider {
	t.Helper()
	var o llm.EndpointOptions
	o.Provider = "mock"
	o.Model = "mock-small"
	o.Protocol = "openai-chat"
	o.BaseURL = base
	o.APIKey = "sk-test-not-a-real-key"
	o.HTTPClient = &http.Client{Transport: &http.Transport{Proxy: nil}}
	p, err := llm.NewProvider(o)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// TestRPM10EleventhRequestWaitsLocally is the AC#7 case end to end: a real
// adapter, a real HTTP server and the pacing decorator. The 11th request is
// proven to be held BEFORE the connection is opened: while it waits, every
// local-wait step re-reads the server's own counter and finds 10, and the
// flow ends in cancellation (the test stops it after a few steps) with the
// server counter frozen and its 429 tripwire never tripped.
func TestRPM10EleventhRequestWaitsLocally(t *testing.T) {
	ps := newPacedServer(t)
	inner := chatProviderAt(t, ps.srv.URL+"/v1")

	clock := &fakeMono{t: time.Unix(1700000000, 0)}
	var mu sync.Mutex
	var waits []time.Duration
	var hitsDuringWait []int
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var steps int
	// The injected sleep stands in for time.Sleep only: it advances the
	// monotonic source (so the refill arithmetic runs unchanged) and records
	// what the SERVER had received at that instant. After a few re-check
	// steps it cancels - the point is to prove the request never left the
	// process, not to spend the full simulated 6s of steps (that total is
	// pinned by TestRPM10RefillQuantumIsOneTokenPerSixSeconds below).
	lim := llm.NewBucketLimiterWithClock(llm.RateLimits{RPM: 10}, clock.Now,
		func(_ context.Context, d time.Duration) error {
			mu.Lock()
			waits = append(waits, d)
			hitsDuringWait = append(hitsDuringWait, ps.hitsCount())
			steps++
			stop := steps >= 3
			mu.Unlock()
			clock.Advance(d) // the wait consumes simulated time
			if stop {
				cancel()
			}
			return nil
		})
	paced := llm.NewLimiterProvider(inner, lim)

	req := &llm.Request{Model: "mock-small", Messages: []llm.Message{{
		Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "hi"}},
	}}}

	// Ten requests fit the burst capacity: each reaches the server
	// immediately, with no local wait at all.
	for i := 0; i < 10; i++ {
		if err := paced.Stream(context.Background(), req, func(llm.StreamEvent) error { return nil }); err != nil {
			t.Fatalf("request %d: %v", i+1, err)
		}
		if got := ps.hitsCount(); got != i+1 {
			t.Fatalf("server hits = %d, want %d", got, i+1)
		}
	}
	mu.Lock()
	burstWaits := len(waits)
	mu.Unlock()
	if burstWaits != 0 {
		t.Fatalf("waits during the allowed burst = %d, want none", burstWaits)
	}

	// The eleventh must NOT reach the server: it waits for the refill.
	err := paced.Stream(ctx, req, func(llm.StreamEvent) error { return nil })
	if err == nil {
		t.Fatal("11th request completed without being held - pacing is not gating")
	}
	if cls, _ := observe.ClassOf(err); cls != observe.ClassCancelled {
		t.Fatalf("11th request err class = %s (%v), want cancelled: it must have been "+
			"stopped while waiting locally, never sent", cls, err)
	}
	if got := ps.hitsCount(); got != 10 {
		t.Errorf("server hits = %d, want 10: the 11th request must never reach "+
			"the fixture while its window is exhausted", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(waits) < 3 {
		t.Fatalf("local wait steps = %d, want >= 3 before the test stopped it (%v)", len(waits), waits)
	}
	for i, w := range waits {
		if w <= 0 || w > 250*time.Millisecond {
			t.Errorf("wait step %d = %v, want a positive step capped at the 250ms re-check granularity", i, w)
		}
		if hitsDuringWait[i] != 10 {
			t.Errorf("at wait step %d the server had seen %d requests, want 10 "+
				"(the 11th must be held LOCALLY, before the connection is opened)", i, hitsDuringWait[i])
		}
	}
	if ps.tooMany() != 0 {
		t.Errorf("server answered %d 429s, want 0: Wisp paced itself first", ps.tooMany())
	}
}

// TestRPM10RefillQuantumIsOneTokenPerSixSeconds proves the wait the cancelled
// case above was cut short of: a full one-token deficit at 10 rpm is a ~6s
// local wait assembled from 250ms re-check steps, measured on the injected
// monotonic source. No HTTP involved, so no server-window semantics can hide
// a miscount: the only gate is the bucket.
func TestRPM10RefillQuantumIsOneTokenPerSixSeconds(t *testing.T) {
	clock := &fakeMono{t: time.Unix(1700000000, 0)}
	var waits []time.Duration
	lim := llm.NewBucketLimiterWithClock(llm.RateLimits{RPM: 10}, clock.Now,
		func(_ context.Context, d time.Duration) error {
			waits = append(waits, d)
			clock.Advance(d)
			return nil
		})
	for i := 0; i < 10; i++ {
		if err := lim.Acquire(context.Background(), 0); err != nil {
			t.Fatalf("acquire %d: %v", i+1, err)
		}
	}
	if len(waits) != 0 {
		t.Fatalf("bursts waited %v, want no waits", waits)
	}
	if err := lim.Acquire(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
	var total time.Duration
	for i, w := range waits {
		if w <= 0 || w > 250*time.Millisecond {
			t.Fatalf("wait step %d = %v, want (0, 250ms]", i, w)
		}
		total += w
	}
	// 10 rpm = exactly one token per 6s; the 250ms re-checks must sum to the
	// deficit and nothing else.
	if total < 5900*time.Millisecond || total > 6100*time.Millisecond {
		t.Errorf("total local wait = %v, want ~6s (one token at 10 rpm)", total)
	}
	if len(waits) < 20 {
		t.Errorf("wait steps = %d, want >= 20 (periodic re-checks, not one long sleep)", len(waits))
	}
}

// TestClockJumpCannotMintTokens proves the bucket holds no wall-clock math: a
// source that jumps backward must neither grant extra tokens nor wedge.
func TestClockJumpCannotMintTokens(t *testing.T) {
	clock := &fakeMono{t: time.Unix(1700000000, 0)}
	var slept time.Duration
	lim := llm.NewBucketLimiterWithClock(llm.RateLimits{RPM: 60}, clock.Now,
		func(_ context.Context, d time.Duration) error {
			slept += d
			clock.Advance(d)
			return nil
		})
	// 60 rpm: burst of 60, then one token per second.
	for i := 0; i < 60; i++ {
		if err := lim.Acquire(context.Background(), 0); err != nil {
			t.Fatal(err)
		}
	}
	// A backward jump of an hour must not refill anything: the next acquire
	// still waits, and the elapsed math stays positive.
	clock.Advance(-time.Hour)
	if err := lim.Acquire(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
	if slept <= 0 {
		t.Errorf("slept = %v, want a positive wait after a backward clock jump "+
			"(tokens must not be minted out of a wall-clock rewind)", slept)
	}
}

// TestUnpacedLimitsNeverBlock: rpm/tpm of 0 mean "no limit configured" (the
// bucket is nil, not infinite-loop).
func TestUnpacedLimitsNeverBlock(t *testing.T) {
	lim := llm.NewBucketLimiter(llm.RateLimits{})
	for i := 0; i < 500; i++ {
		if err := lim.Acquire(context.Background(), 1000); err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
	}
}

// TestCancellationWhilePacingIsNotAProviderFailure: waiting must honor the
// caller's context and classify as cancelled, never as an error to retry.
// This one runs on the REAL production wiring (time.Now + time.Sleep).
func TestCancellationWhilePacingIsNotAProviderFailure(t *testing.T) {
	ps := newPacedServer(t)
	inner := chatProviderAt(t, ps.srv.URL+"/v1")
	lim := llm.NewBucketLimiter(llm.RateLimits{RPM: 1})
	paced := llm.NewLimiterProvider(inner, lim)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // owner: this test; cancel only, no work escapes
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	// The first request takes the only token, so the second one waits.
	if err := paced.Stream(context.Background(), adapterReq(), func(llm.StreamEvent) error { return nil }); err != nil {
		t.Fatal(err)
	}
	err := paced.Stream(ctx, adapterReq(), func(llm.StreamEvent) error { return nil })
	if err == nil {
		t.Fatal("expected the paced request to be cut short by cancellation")
	}
	if cls, _ := observe.ClassOf(err); cls != observe.ClassCancelled {
		t.Errorf("class = %s, want cancelled (%v)", cls, err)
	}
}

func adapterReq() *llm.Request {
	return &llm.Request{Model: "mock-small", Messages: []llm.Message{{
		Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "hi"}},
	}}}
}
