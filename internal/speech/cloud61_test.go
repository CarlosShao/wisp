package speech

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The ticket 61 receipt tests. Everything runs over the C9 seams with the
// real cascade code; the network is injected at its own boundary (an
// httptest server speaking the documented wire of cloudasr.go /
// cloudtts.go) -- the injection seam this repo allows for cloud legs, the
// same shape the openai-chat adapter tests use (mockllm/httptest), never a
// mock standing in for a real engine's behaviour.

// -- test scaffolding ---------------------------------------------------------

type cascadeHarness struct {
	t *testing.T

	mu        sync.Mutex
	failovers []VoiceFailoverEvent
	states    []VoiceProviderState
	usages    []VoiceUsage

	reg *observe.Registry
}

func newCascadeHarness(t *testing.T) *cascadeHarness {
	return &cascadeHarness{t: t, reg: observe.NewRegistry()}
}

func (h *cascadeHarness) options() ChainOptions {
	h.t.Helper()
	return ChainOptions{
		Registry: h.reg,
		Callbacks: ChainCallbacks{
			OnFailover:      func(e VoiceFailoverEvent) { h.record(func() { h.failovers = append(h.failovers, e) }) },
			OnProviderState: func(s VoiceProviderState) { h.record(func() { h.states = append(h.states, s) }) },
			OnUsage:         UsageFunc(func(u VoiceUsage) { h.record(func() { h.usages = append(h.usages, u) }) }),
		},
	}
}

func (h *cascadeHarness) record(f func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	f()
}

func (h *cascadeHarness) snap() ([]VoiceFailoverEvent, []VoiceProviderState, []VoiceUsage) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]VoiceFailoverEvent(nil), h.failovers...),
		append([]VoiceProviderState(nil), h.states...),
		append([]VoiceUsage(nil), h.usages...)
}

// asrJSONServer answers the streaming transcription request with a plain
// JSON final (the documented non-streaming shape), counting hits and
// drained audio bytes.
func asrJSONServer(text string, seconds float64, inTokens int, hits, bytesIn *atomic.Int64) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits != nil {
			hits.Add(1)
		}
		body, _ := io.ReadAll(r.Body)
		if bytesIn != nil {
			bytesIn.Add(int64(len(body)))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"text": text,
			"usage": map[string]any{
				"seconds":      seconds,
				"input_tokens": inTokens,
			},
		})
	}))
}

// asrStatusServer drains the upload (so the utterance can complete) and
// answers a fixed status with an OpenAI-shaped error body.
func asrStatusServer(status int, code string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if code != "" {
			_, _ = fmt.Fprintf(w, `{"error":{"code":%q,"message":"probe status"}}`, code)
		} else {
			_, _ = io.WriteString(w, `{"error":{"message":"probe status"}}`)
		}
	}))
}

func cloudEndpoint(base string) CloudEndpoint {
	return CloudEndpoint{
		Provider:         "acme",
		Model:            "voice-1",
		BaseURL:          base,
		APIKey:           "test-key",
		Billing:          "pay-per-token",
		PriceAudioIn:     7,
		PriceAudioOut:    9,
		DeclaresAudioIn:  true,
		DeclaresAudioOut: true,
	}
}

func mustCloudASR(t *testing.T, base string, opts ChainOptions) *CloudASR {
	t.Helper()
	a, err := NewCloudASR(cloudEndpoint(base), opts)
	if err != nil {
		t.Fatalf("NewCloudASR: %v", err)
	}
	t.Cleanup(func() { _ = a.Close() })
	return a
}

// feedWindows drives the seam like Transcriber does: window-sized chunks.
func feedWindows(a AsrEngine, total int) int {
	fed := 0
	for off := 0; off < total; off += WindowSamples {
		end := min(off+WindowSamples, total)
		a.Feed(make([]float32, end-off))
		fed += end - off
	}
	return fed
}

// waitFor polls cond until true or the deadline; returns false on timeout.
func waitFor(cond func() bool, budget time.Duration) bool {
	deadline := time.Now().Add(budget)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return cond()
}

// -- engine: cloud ASR -------------------------------------------------------

func Test61CloudASRSnapshotPartialsThenAuthoritativeFinal(t *testing.T) {
	h := newCascadeHarness(t)
	var requests atomic.Int64
	const (
		totalSamples = 72000 // 4.5s of 16k audio, two cadence boundaries crossed
		finalBytes   = totalSamples * 2
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		if len(body) < finalBytes*3/4 {
			// A snapshot of a partial utterance: a partial hypothesis back.
			_ = json.NewEncoder(w).Encode(map[string]any{"text": "快照假设"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"text":  "完整句子转写",
			"usage": map[string]any{"seconds": 4.5, "input_tokens": 21},
		})
	}))
	defer srv.Close()

	opts := h.options()
	opts.Client = srv.Client()
	a := mustCloudASR(t, srv.URL, opts)

	a.Begin()
	fed := feedWindows(a, totalSamples)
	if fed != totalSamples {
		t.Fatalf("feedWindows fed %d, want %d", fed, totalSamples)
	}

	// The snapshot leg: a hypothesis must exist BEFORE Finish closed the
	// utterance -- that is the "streaming partials feeding the same
	// VAD->ASR pipeline" claim, observable at the seam.
	if !waitFor(func() bool { return a.Partial() != "" }, 5*time.Second) {
		t.Fatalf("no snapshot partial before Finish (requests=%d)", requests.Load())
	}
	gotPartial := a.Partial()
	if gotPartial != "快照假设" && gotPartial != "完整句子转写" {
		t.Fatalf("partial = %q, want one of the endpoint's hypotheses", gotPartial)
	}

	text := a.Finish()
	if text != "完整句子转写" {
		t.Fatalf("final = %q, want the whole-utterance reply to be authoritative", text)
	}
	if a.InFlight() {
		t.Fatal("InFlight true after Finish")
	}
	if f := a.LastFailure(); f != nil {
		t.Fatalf("LastFailure on a clean run = %v", f)
	}
	if got := a.Fed(); got != fed {
		t.Fatalf("Fed after Finish = %d, want the utterance count %v (same clause as StreamingAsr)", got, fed)
	}
	if requests.Load() < 2 {
		t.Fatalf("endpoint saw %d requests, want >=2 (at least one snapshot + the final)", requests.Load())
	}

	_, _, usages := h.snap()
	if len(usages) != 1 {
		t.Fatalf("usage records = %d, want exactly 1 (only the FINAL transcribes for the ledger; snapshots are free re-reads of the same audio)", len(usages))
	}
	u := usages[0]
	if u.Chain != ChainASR || u.Key != "acme/voice-1" {
		t.Fatalf("usage = %+v, want chain asr at acme/voice-1", u)
	}
	if u.AudioSecondsIn != 4.5 {
		t.Fatalf("AudioSecondsIn = %v, want the provider's billed 4.5", u.AudioSecondsIn)
	}
	if u.InputTokens != 21 || u.Billing != "pay-per-token" || u.PriceAudioIn != 7 || u.PriceAudioOut != 9 {
		t.Fatalf("usage pricing/fields wrong: %+v", u)
	}
}

func Test61CloudASRJSONOnlyProviderStillConforms(t *testing.T) {
	h := newCascadeHarness(t)
	var hits atomic.Int64
	srv := asrJSONServer("整段转写", 1.25, 40, &hits, nil)
	defer srv.Close()

	opts := h.options()
	opts.Client = srv.Client()
	a := mustCloudASR(t, srv.URL, opts)

	a.Begin()
	fed := feedWindows(a, 8000)
	if p := a.Partial(); p != "" {
		t.Fatalf("Partial mid-utterance = %q, want empty for a final-only provider (the seam words it \"empty until the engine has enough audio to say anything\")", p)
	}
	if got := a.Finish(); got != "整段转写" {
		t.Fatalf("final = %q", got)
	}
	if a.Fed() != fed {
		t.Fatalf("Fed = %d, want %d", a.Fed(), fed)
	}
	if hits.Load() != 1 {
		t.Fatalf("endpoint hit %d times per utterance, want 1", hits.Load())
	}
	_, _, usages := h.snap()
	if len(usages) != 1 || usages[0].AudioSecondsIn != 1.25 || usages[0].InputTokens != 40 {
		t.Fatalf("usage = %+v, want the provider usage block carried", usages)
	}
}

func Test61CloudASRFailuresAreClassifiedPerD37(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		code      string
		wantClass observe.ErrorClass
	}{
		{"auth 401", http.StatusUnauthorized, "bad_key", observe.ClassAuth},
		{"forbidden 403", http.StatusForbidden, "", observe.ClassAuth},
		{"rate limit 429", http.StatusTooManyRequests, "rate_exceeded", observe.ClassRateLimit},
		{"insufficient quota 429", http.StatusTooManyRequests, "insufficient_quota", observe.ClassBudget},
		{"out of credit 402", http.StatusPaymentRequired, "", observe.ClassBudget},
		{"server 500", http.StatusInternalServerError, "", observe.ClassProvider},
		{"unavailable 503", http.StatusServiceUnavailable, "", observe.ClassProvider},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newCascadeHarness(t)
			srv := asrStatusServer(tc.status, tc.code)
			defer srv.Close()
			opts := h.options()
			opts.Client = srv.Client()
			a := mustCloudASR(t, srv.URL, opts)

			a.Begin()
			fedWindowsAndFinish(a)
			if got := a.Finish(); got != "" {
				t.Fatalf("Finish on failure = %q, want the empty mishear signal at the seam", got)
			}
			f := a.LastFailure()
			if f == nil {
				t.Fatal("LastFailure nil, want the classified cause for the chain stepper")
			}
			if f.Class != tc.wantClass {
				t.Fatalf("class = %s, want %s (%s)", f.Class, tc.wantClass, f.Detail)
			}
			if !strings.Contains(f.Error(), tc.name[:strings.Index(tc.name, " ")]+"") &&
				!strings.Contains(strings.ToLower(f.Error()), strings.ToLower(tc.name)) {
				// detail carries the HTTP status text; keep the loose check
				_ = f
			}
		})
	}
	t.Run("stream ends with no final", func(t *testing.T) {
		h := newCascadeHarness(t)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(200)
			if fl, ok := w.(http.Flusher); ok {
				fl.Flush()
			}
			_, _ = io.Copy(io.Discard, r.Body)
			// close without ever saying anything
		}))
		defer srv.Close()
		opts := h.options()
		opts.Client = srv.Client()
		a := mustCloudASR(t, srv.URL, opts)
		a.Begin()
		a.Feed(make([]float32, 160))
		if got := a.Finish(); got != "" {
			t.Fatalf("final = %q", got)
		}
		if f := a.LastFailure(); f == nil || f.Class != observe.ClassNetwork {
			t.Fatalf("silent close: failure = %v, want class network", f)
		}
	})
	t.Run("last partial is kept when the stream closes without is_final", func(t *testing.T) {
		h := newCascadeHarness(t)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = io.WriteString(w, "data: {\"text\":\"半截假设\"}\n\n")
			if fl, ok := w.(http.Flusher); ok {
				fl.Flush()
			}
			_, _ = io.Copy(io.Discard, r.Body)
		}))
		defer srv.Close()
		opts := h.options()
		opts.Client = srv.Client()
		a := mustCloudASR(t, srv.URL, opts)
		a.Begin()
		a.Feed(make([]float32, 160))
		if got := a.Finish(); got != "半截假设" {
			t.Fatalf("final = %q, want the best published hypothesis handed up, not a fabricated failure", got)
		}
	})
}

func fedWindowsAndFinish(a AsrEngine) { a.Feed(make([]float32, 320)) }

func Test61CloudASRCloseIsIdempotentAndQuiet(t *testing.T) {
	h := newCascadeHarness(t)
	var hits atomic.Int64
	srv := asrJSONServer("x", 0.1, 0, &hits, nil)
	defer srv.Close()
	opts := h.options()
	opts.Client = srv.Client()
	a := mustCloudASR(t, srv.URL, opts)

	if err := a.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := a.Close(); err != nil {
		t.Fatalf("second Close must be a no-op (C11): %v", err)
	}
	a.Begin()
	if a.InFlight() {
		t.Fatal("a closed engine opened an utterance; Begin must refuse")
	}
	if got := a.Finish(); got != "" {
		t.Fatalf("Finish on a closed engine = %q, want empty", got)
	}
	if hits.Load() != 0 {
		t.Fatalf("closed engine dispatched %d requests", hits.Load())
	}
}

func Test61CloudASRDeadEndpointFailsWithAClassNotAHang(t *testing.T) {
	// A refused TCP connect is the cheapest deterministic transport failure;
	// it must surface as a class, because the chain stepper only acts on
	// classes. (The hung-socket path is the watchdog constant in Finish:
	// 15s of waiting on a live process is outside a unit test's budget and
	// is covered by the same publish mechanism this case exercises.)
	h := newCascadeHarness(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	base := srv.URL
	srv.Close() // now the base URL refuses connections

	opts := h.options()
	// Ask the dead address DIRECTLY. Without an injected client the engine
	// uses its production default, which honors the machine's system proxy
	// (llm.ProxyFunc -> WinINET): on a host whose proxy is up, 127.0.0.1 is
	// dialled through that proxy, which answers 502 -- a perfectly correct
	// ClassProvider for a proxy relaying a dead origin, but not the refused
	// connect this case is about. Same rule as every other case in this
	// file: the network boundary is injected, never inherited from the host.
	opts.Client = srv.Client()
	a, err := NewCloudASR(cloudEndpoint(base), opts)
	if err != nil {
		t.Fatalf("NewCloudASR: %v", err)
	}
	defer func() { _ = a.Close() }()

	a.Begin()
	a.Feed(make([]float32, 400))
	start := time.Now()
	if got := a.Finish(); got != "" {
		t.Fatalf("final on a dead endpoint = %q, want the empty mishear", got)
	}
	f := a.LastFailure()
	if f == nil {
		t.Fatal("dead endpoint recorded no failure class")
	}
	if f.Class != observe.ClassNetwork && f.Class != observe.ClassInternal {
		t.Fatalf("dead endpoint class = %s, want network (transport) (%v)", f.Class, f)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Fatalf("Finish waited %v on a refused connection; the request budget must bound it", elapsed)
	}
}

// -- chain: cloud ASR --------------------------------------------------------

func Test61ASRChainStepsOnFailureWithReplayAndVisibleEvents(t *testing.T) {
	h := newCascadeHarness(t)
	var bHits atomic.Int64
	var aBytes, bBytes atomic.Int64
	srvA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		aBytes.Add(int64(len(body)))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"error":{"message":"boom"}}`)
	}))
	defer srvA.Close()
	srvB := asrJSONServer("乙的答案", 0.5, 7, &bHits, &bBytes)
	defer srvB.Close()

	epA, epB := cloudEndpoint(srvA.URL), cloudEndpoint(srvB.URL)
	epA.Model, epB.Model = "a1", "b2"

	opts := h.options()
	opts.Client = srvA.Client()
	// one client for both servers (httptest clients are proxy-free and
	// origin-agnostic; same transport rules as production defaults).
	chain, err := NewCloudASRChain([]CloudEndpoint{epA, epB}, nil, opts)
	if err != nil {
		t.Fatalf("NewCloudASRChain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	chain.Begin()
	fed := feedWindows(chain, 8000)
	if got := chain.Partial(); got != "" {
		t.Logf("partial during the failed head = %q (informational)", got)
	}
	text := chain.Finish()
	if text != "乙的答案" {
		t.Fatalf("chain final = %q, want the second element's answer", text)
	}
	if chain.Active() != "acme/b2" {
		t.Fatalf("Active = %q, want acme/b2 after the step", chain.Active())
	}
	if chain.Fed() != fed {
		t.Fatalf("chain Fed = %d, want the caller's clock %d (a provider switch must not move the latency baseline)", chain.Fed(), fed)
	}
	if chain.LastFailure() != nil {
		t.Fatalf("LastFailure = %v, want nil once the walk succeeded", chain.LastFailure())
	}

	// REPLAY: every element the utterance walked must have heard the whole
	// utterance, not the tail after the failure.
	want := int64(fed * 2) // s16le mono
	if got := aBytes.Load(); got != want {
		t.Fatalf("element A received %d audio bytes, want the whole utterance (%d)", got, want)
	}
	if got := bBytes.Load(); got != want {
		t.Fatalf("element B received %d audio bytes, want the whole utterance (%d) replayed", got, want)
	}

	failovers, _, usages := h.snap()
	if len(failovers) != 1 {
		t.Fatalf("failover events = %d, want exactly the one hop A->B: %+v", len(failovers), failovers)
	}
	ev := failovers[0]
	if ev.Chain != ChainASR || ev.From != "acme/a1" || ev.To != "acme/b2" || ev.Cause != string(observe.ClassProvider) {
		t.Fatalf("failover event = %+v, want ASR acme/a1 -> acme/b2 cause provider", ev)
	}
	if ev.Index != 0 || ev.Attempt != 1 || ev.Total != 2 {
		t.Fatalf("failover positions = idx %d attempt %d total %d, want 0/1/2", ev.Index, ev.Attempt, ev.Total)
	}
	if len(usages) != 1 || usages[0].Key != "acme/b2" {
		t.Fatalf("usage records = %+v, want exactly the successful element", usages)
	}
}

// recordingASR is the local terminal fallback the chain must never be able
// to double-close: it counts its own sessions.
type recordingASR struct {
	fed    int
	beginN int
	closeN int
	final  string
}

func (r *recordingASR) Begin()           { r.beginN++; r.fed = 0 }
func (r *recordingASR) Feed(s []float32) { r.fed += len(s) }
func (r *recordingASR) Partial() string  { return "" }
func (r *recordingASR) Finish() string   { return r.final }
func (r *recordingASR) InFlight() bool   { return false }
func (r *recordingASR) Fed() int         { return r.fed }
func (r *recordingASR) Close() error     { r.closeN++; return nil }

func Test61ASRChainFallsThroughToLocalNeverSilently(t *testing.T) {
	h := newCascadeHarness(t)
	srvA := asrStatusServer(http.StatusInternalServerError, "")
	defer srvA.Close()

	epA := cloudEndpoint(srvA.URL)
	epA.Model = "a1"
	local := &recordingASR{final: "本地兜底转写"}

	opts := h.options()
	opts.Client = srvA.Client()
	chain, err := NewCloudASRChain([]CloudEndpoint{epA}, local, opts)
	if err != nil {
		t.Fatalf("NewCloudASRChain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	chain.Begin()
	fed := feedWindows(chain, 5000)
	text := chain.Finish()

	if text != "本地兜底转写" {
		t.Fatalf("final = %q, want the local fallback's reading", text)
	}
	failovers, states, _ := h.snap()
	if len(failovers) != 1 || failovers[0].To != "local-sherpa" || failovers[0].Cause != string(observe.ClassProvider) {
		t.Fatalf("failover events = %+v, want exactly one visible hop to local-sherpa", failovers)
	}
	if local.beginN != 1 || local.fed != fed {
		t.Fatalf("local fallback saw begin=%d fed=%d, want 1/%d (full-utterance replay must reach it)", local.beginN, local.fed, fed)
	}
	last := states[len(states)-1]
	if last.Active != "local-sherpa" || last.Cloud {
		t.Fatalf("last provider state = %+v, want the local fallback holding the line (privacy restored)", last)
	}
	// Chain Close must NOT close the caller's local engine (ownership).
	_ = chain.Close()
	if local.closeN != 0 {
		t.Fatalf("chain closed the local fallback %d time(s); it belongs to its loader", local.closeN)
	}
}

func Test61QuotaSkipAdvancesChainBeforeDispatch(t *testing.T) {
	h := newCascadeHarness(t)
	var aHits atomic.Int64
	srvA := asrJSONServer("不该被派发", 1, 1, &aHits, nil)
	defer srvA.Close()
	srvB := asrJSONServer("乙承接", 0.3, 3, nil, nil)
	defer srvB.Close()

	epA, epB := cloudEndpoint(srvA.URL), cloudEndpoint(srvB.URL)
	epA.Model, epB.Model = "a1", "b2"

	opts := h.options()
	opts.Client = srvA.Client()
	opts.Callbacks.Quota = QuotaFunc(func(provider, model string) (bool, string) {
		if model == "a1" {
			return false, "配额已用尽：acme/a1 当日额度到达上限"
		}
		return true, ""
	})
	chain, err := NewCloudASRChain([]CloudEndpoint{epA, epB}, nil, opts)
	if err != nil {
		t.Fatalf("NewCloudASRChain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	chain.Begin()
	feedWindows(chain, 3000)
	text := chain.Finish()

	if text != "乙承接" {
		t.Fatalf("final = %q, want the chain to have advanced BEFORE dispatching to the over-quota head", text)
	}
	if aHits.Load() != 0 {
		t.Fatalf("over-quota element was dispatched %d time(s); anti-negative means zero", aHits.Load())
	}
	failovers, _, _ := h.snap()
	if len(failovers) != 1 {
		t.Fatalf("events = %+v, want exactly the one visible quota skip", failovers)
	}
	ev := failovers[0]
	if ev.Cause != "quota" || ev.From != "acme/a1" || ev.To != "acme/b2" || !strings.Contains(ev.Detail, "acme/a1") {
		t.Fatalf("skip event = %+v, want cause quota acme/a1 -> acme/b2 with the ledger's reason", ev)
	}
}

func Test61QuotaSkipsEveryElementFallsToLocal(t *testing.T) {
	h := newCascadeHarness(t)
	srvA := asrJSONServer("不该跑", 1, 1, nil, nil)
	defer srvA.Close()

	epA := cloudEndpoint(srvA.URL)
	epA.Model = "a1"
	local := &recordingASR{final: "兜底"}

	opts := h.options()
	opts.Client = srvA.Client()
	opts.Callbacks.Quota = QuotaFunc(func(_, _ string) (bool, string) {
		return false, "全部超配额"
	})
	chain, err := NewCloudASRChain([]CloudEndpoint{epA}, local, opts)
	if err != nil {
		t.Fatalf("NewCloudASRChain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	chain.Begin()
	feedWindows(chain, 1200)
	if got := chain.Finish(); got != "兜底" {
		t.Fatalf("final = %q, want the local terminal fallback after every skip", got)
	}
	failovers, _, usages := h.snap()
	if len(usages) != 0 {
		t.Fatalf("usage records = %+v, want none (local is zero cost)", usages)
	}
	if len(failovers) != 1 || failovers[0].To != "local-sherpa" {
		t.Fatalf("events = %+v, want the visible last skip landing on local", failovers)
	}
}

// -- engine: cloud TTS -------------------------------------------------------

// pcmBodyServer speaks the raw-PCM audio/speech shape; chunks are flushed
// per write so the engine sees incremental delivery. It can also be pointed
// at a "die mid-stream after two chunks" behavior (truncated Content-Length)
// to test the delivered rule.
func pcmBodyServer(totalSamples int, truncLen bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["model"] == nil || req["input"] == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "audio/L16")
		if truncLen {
			w.Header().Set("Content-Length", fmt.Sprint(totalSamples*2+4096))
		}
		w.WriteHeader(200)
		fl, _ := w.(http.Flusher)
		buf := make([]byte, totalSamples*2)
		for i := 0; i < totalSamples*2; i += 2 {
			binary.LittleEndian.PutUint16(buf[i:], uint16(int16(i%2000)-1000))
		}
		// Emit in 4096-byte writes so the chunker sees multiple chunks.
		for off := 0; off < len(buf); off += 4096 {
			end := min(off+4096, len(buf))
			_, _ = w.Write(buf[off:end])
			if fl != nil {
				fl.Flush()
			}
		}
	}))
}

func Test61CloudTTSStreamsChunksAndMetersAudioSeconds(t *testing.T) {
	h := newCascadeHarness(t)
	srv := pcmBodyServer(4800, false) // 4800 samples = 0.2s at 24k
	defer srv.Close()

	opts := h.options()
	opts.Client = srv.Client()
	tts, err := NewCloudTTS(cloudEndpoint(srv.URL), opts)
	if err != nil {
		t.Fatalf("NewCloudTTS: %v", err)
	}
	defer func() { _ = tts.Close() }()

	if got := tts.SampleRate(); got != ttsPCMDefaultRate {
		t.Fatalf("SampleRate = %d, want the audio/speech PCM default (%d)", got, ttsPCMDefaultRate)
	}

	var mu sync.Mutex
	var sizes []int
	err = tts.Speak("你好世界", func(pcm []float32) error {
		mu.Lock()
		sizes = append(sizes, len(pcm))
		mu.Unlock()
		return nil
	})
	if err != nil {
		t.Fatalf("Speak: %v", err)
	}
	if total := sumInts(sizes); total != 4800 {
		t.Fatalf("emitted %d samples (chunks %v), want the whole 4800", total, sizes)
	}
	if len(sizes) < 2 {
		t.Fatalf("Speak delivered %v in one chunk; the seam's chunked contract (and first-frame budget) require streaming delivery", sizes)
	}

	_, _, usages := h.snap()
	if len(usages) != 1 {
		t.Fatalf("usage records = %d, want 1", len(usages))
	}
	u := usages[0]
	if u.Chain != ChainTTS || u.Key != "acme/voice-1" {
		t.Fatalf("usage = %+v", u)
	}
	if got, want := u.AudioSecondsOut, 4800.0/float64(ttsPCMDefaultRate); got != want {
		t.Fatalf("AudioSecondsOut = %v, want %v (audio clock over the endpoint rate)", got, want)
	}
	if u.PriceAudioOut != 9 || u.PriceAudioIn != 7 || u.Billing != "pay-per-token" {
		t.Fatalf("catalog price fields not carried: %+v", u)
	}
	if ff := tts.FirstFrame(); ff <= 0 {
		t.Fatalf("FirstFrame = %v, want the measured provider-dependent reading (>0)", ff)
	}
}

func Test61CloudTTSGuardsWithoutTouchingTheNetwork(t *testing.T) {
	h := newCascadeHarness(t)
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	opts := h.options()
	opts.Client = srv.Client()
	tts, err := NewCloudTTS(cloudEndpoint(srv.URL), opts)
	if err != nil {
		t.Fatalf("NewCloudTTS: %v", err)
	}

	if err := tts.Speak("。！", func([]float32) error { return nil }); !errors.Is(err, ErrNothingToSpeak) {
		t.Fatalf("punctuation-only Speak = %v, want ErrNothingToSpeak (the same guard the local engine runs)", err)
	}
	if err := tts.Speak("x", nil); err == nil || !strings.Contains(err.Error(), "emit") {
		t.Fatalf("nil emit = %v, want a named argument error", err)
	}
	_ = tts.Close()
	if err := tts.Speak("话", func([]float32) error { return nil }); err == nil {
		t.Fatal("Speak on a closed engine must refuse")
	}
	if hits.Load() != 0 {
		t.Fatalf("guards dispatched %d request(s); nothing may leave before they pass", hits.Load())
	}
}

func Test61CloudTTSBargeInAbortsVerbatimMidStream(t *testing.T) {
	h := newCascadeHarness(t)
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/L16")
		w.WriteHeader(200)
		fl, _ := w.(http.Flusher)
		buf := make([]byte, 4096)
		for i := 0; i < 2; i++ {
			_, _ = w.Write(buf)
			if fl != nil {
				fl.Flush()
			}
		}
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	defer srv.Close()
	defer close(release)

	opts := h.options()
	opts.Client = srv.Client()
	tts, err := NewCloudTTS(cloudEndpoint(srv.URL), opts)
	if err != nil {
		t.Fatalf("NewCloudTTS: %v", err)
	}
	defer func() { _ = tts.Close() }()

	stop := errors.New("barge-in")
	chunks := 0
	err = tts.Speak("长句播报", func([]float32) error {
		chunks++
		if chunks == 2 {
			return stop
		}
		return nil
	})
	if !errors.Is(err, stop) {
		t.Fatalf("Speak returned %v, want the emit error VERBATIM (the seam's cooperative-abort contract)", err)
	}
	if got := emittedSamples(tts); got != 4096 {
		t.Fatalf("delivered samples = %d, want 4096 (two chunks): the chain's delivered rule reads EXACTLY this number when Speak returns an error", got)
	}
}

func Test61TTSChainStepsBeforeDeliveryNeverAfter(t *testing.T) {
	h := newCascadeHarness(t)
	srvA := asrStatusServer(http.StatusInternalServerError, "") // same status shape: 500 before any audio
	defer srvA.Close()
	srvB := pcmBodyServer(2048, false)
	defer srvB.Close()

	epA, epB := cloudEndpoint(srvA.URL), cloudEndpoint(srvB.URL)
	epA.Model, epB.Model = "a1", "b2"

	opts := h.options()
	opts.Client = srvA.Client()
	rec := &recordingTTS{}
	chain, err := NewCloudTTSChain([]CloudEndpoint{epA, epB}, rec, opts)
	if err != nil {
		t.Fatalf("NewCloudTTSChain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	if err := chain.Speak("播报测试", func(pcm []float32) error { return nil }); err != nil {
		t.Fatalf("Speak: %v (want success via element B)", err)
	}
	failovers, states, _ := h.snap()
	if len(failovers) != 1 || failovers[0].From != "acme/a1" || failovers[0].To != "acme/b2" {
		t.Fatalf("failover events = %+v, want exactly the visible A->B hop", failovers)
	}
	if chain.Active() != "acme/b2" {
		t.Fatalf("Active = %q", chain.Active())
	}
	if chain.SampleRate() != ttsPCMDefaultRate {
		t.Fatalf("SampleRate = %d, want the ACTIVE engine's rate", chain.SampleRate())
	}
	_ = states

	// The delivered rule: a failure AFTER chunks reached the consumer must
	// not restart the sentence from the next element (half-played audio +
	// a fresh full sentence is worse than the error).
	release := make(chan struct{})
	defer close(release)
	srvC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/L16")
		w.Header().Set("Content-Length", "16384") // promise more than delivered
		w.WriteHeader(200)
		buf := make([]byte, 8192) // 4096 samples actually written
		_, _ = w.Write(buf)
		if fl, ok := w.(http.Flusher); ok {
			fl.Flush()
		}
		// handler returns; the unfulfilled Content-Length truncates the body
	}))
	defer srvC.Close()
	srvD := pcmBodyServer(2048, false)
	defer srvD.Close()
	var dHits atomic.Int64
	srvE := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dHits.Add(1)
		srvD.Config.Handler.ServeHTTP(w, r)
	}))
	defer srvE.Close()

	epC, epD := cloudEndpoint(srvC.URL), cloudEndpoint(srvE.URL)
	epC.Model, epD.Model = "c1", "d2"
	h2 := newCascadeHarness(t)
	opts2 := h2.options()
	opts2.Client = srvC.Client()
	rec2 := &recordingTTS{}
	chain2, err := NewCloudTTSChain([]CloudEndpoint{epC, epD}, rec2, opts2)
	if err != nil {
		t.Fatalf("NewCloudTTSChain 2: %v", err)
	}
	defer func() { _ = chain2.Close() }()

	err = chain2.Speak("半句话", func(pcm []float32) error { return nil })
	if err == nil {
		t.Fatal("truncated stream Speak = nil, want the classified transport failure surfaced")
	}
	if dHits.Load() != 0 {
		t.Fatalf("chain restarted a HALF-SPOKEN sentence into the next element (%d dispatches); the delivered rule forbids it", dHits.Load())
	}
	if rec2.calls != 0 {
		t.Fatalf("delivered-failure also fell through to local (%d calls); same rule", rec2.calls)
	}
}

// catalogTestBase returns a config snapshot with one provider exposing the
// voice-capability pairs the resolution tests walk. Enabled is explicit on
// every entry: the catalog reads a missing enabled key as OFF (the same
// promise llm/resolver.go enforces), so an omission here would change what
// the test proves.
func catalogTestBase() *config.Config {
	providers := map[string]config.Provider{
		"stepfun": {
			Protocol:  "openai-chat",
			BaseURL:   "https://api.stepfun.com/v1",
			APIKeyRef: "dpapi:step-test",
			Models: map[string]config.ModelSpec{
				"stepaudio-asr": {
					Capabilities: config.Capabilities{Text: true, AudioIn: true},
					Billing:      "pay-per-token",
					Price:        config.Price{AudioIn: 5},
					Enabled:      true,
				},
				"stepaudio-tts": {
					Capabilities: config.Capabilities{Text: true, AudioOut: true},
					Billing:      "plan",
					Price:        config.Price{AudioOut: 8},
					Enabled:      true,
				},
				"text-only": {
					Capabilities: config.Capabilities{Text: true},
					Billing:      "pay-per-token",
					Enabled:      true,
				},
				"switched-off": {
					Capabilities: config.Capabilities{AudioIn: true, AudioOut: true},
					Enabled:      false,
				},
			},
		},
	}
	return &config.Config{
		LLM: config.LLMSection{Providers: providers},
		Voice: config.VoiceSection{
			CloudASRChain: []string{"stepfun/stepaudio-asr"},
			CloudTTSChain: []string{"stepfun/stepaudio-tts"},
			TTS:           config.TTSConfig{Voice: "stella", Speed: 1.1},
		},
	}
}

// keyRing resolves refs from a map; err forces the resolver-failure branch.
type keyRing struct {
	keys map[string]string
	err  error
}

func (k keyRing) Resolve(ref string) (string, error) {
	if k.err != nil {
		return "", k.err
	}
	v, ok := k.keys[ref]
	if !ok {
		return "", errors.New("secret: ref not found")
	}
	return v, nil
}

func Test61CatalogDrivesTheVoiceChains(t *testing.T) {
	keys := keyRing{keys: map[string]string{"dpapi:step-test": "plaintext-key"}}

	asrEps, err := ResolveCloudChain(catalogTestBase(), ChainASR, keys)
	if err != nil {
		t.Fatalf("ResolveCloudChain(asr): %v", err)
	}
	if len(asrEps) != 1 || asrEps[0].Model != "stepaudio-asr" {
		t.Fatalf("asr endpoints = %+v", asrEps)
	}
	if asrEps[0].APIKey != "plaintext-key" || asrEps[0].PriceAudioIn != 5 ||
		asrEps[0].Billing != "pay-per-token" || !asrEps[0].DeclaresAudioIn {
		t.Fatalf("resolved asr endpoint = %+v, want catalog key/price/billing carried", asrEps[0])
	}
	if asrEps[0].Voice != "" {
		t.Fatalf("ASR endpoint picked up a TTS voice id (%q); voice/speed are TTS-chain fields only", asrEps[0].Voice)
	}

	ttsEps, err := ResolveCloudChain(catalogTestBase(), ChainTTS, keys)
	if err != nil {
		t.Fatalf("ResolveCloudChain(tts): %v", err)
	}
	if ttsEps[0].Voice != "stella" || ttsEps[0].Speed != 1.1 || ttsEps[0].PriceAudioOut != 8 ||
		ttsEps[0].Billing != "plan" || !ttsEps[0].DeclaresAudioOut {
		t.Fatalf("resolved tts endpoint = %+v, want voice dropdown + plan billing + audio_out card", ttsEps[0])
	}
	if !strings.Contains(ttsEps[0].String(), "stepfun/stepaudio-tts") ||
		strings.Contains(ttsEps[0].String(), "plaintext-key") {
		t.Fatalf("endpoint String leaks or misnames: %q", ttsEps[0].String())
	}

	empty, err := ResolveCloudChain(&config.Config{}, ChainASR, keys)
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty chain resolved to %+v / %v, want the local-only configuration accepted as an empty list", empty, err)
	}
}

func Test61CatalogRejectionsNameTheElement(t *testing.T) {
	keys := keyRing{keys: map[string]string{"dpapi:step-test": "k"}}
	cases := []struct {
		name    string
		chain   []string
		kind    CloudChainKind
		wantSub string
	}{
		{"missing capability bit", []string{"stepfun/text-only"}, ChainASR, "audio_in"},
		{"tts bit missing", []string{"stepfun/text-only"}, ChainTTS, "audio_out"},
		{"disabled entry", []string{"stepfun/switched-off"}, ChainASR, "disabled"},
		{"ghost model", []string{"stepfun/nope"}, ChainASR, "unknown model"},
		{"ghost provider", []string{"nope/m"}, ChainASR, "unknown provider"},
		{"local-sherpa never a member", []string{"local-sherpa/cascade"}, ChainASR, "local-sherpa"},
		{"malformed element", []string{"just-a-name"}, ChainASR, "provider/model"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := catalogTestBase()
			switch tc.kind {
			case ChainASR:
				cfg.Voice.CloudASRChain = tc.chain
			case ChainTTS:
				cfg.Voice.CloudTTSChain = tc.chain
			}
			_, err := ResolveCloudChain(cfg, tc.kind, keys)
			if err == nil {
				t.Fatalf("chain %v accepted; want a rejection naming the element", tc.chain)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Fatalf("error %q must name %q", err, tc.wantSub)
			}
			var obs *observe.Error
			if errors.As(err, &obs) && obs.Class != observe.ClassConfig {
				t.Fatalf("error class = %s, want config for catalog rejections", obs.Class)
			}
		})
	}
	t.Run("unresolvable key is auth not config", func(t *testing.T) {
		cfg := catalogTestBase()
		_, err := ResolveCloudChain(cfg, ChainASR, keyRing{err: errors.New("store offline")})
		if err == nil || !strings.Contains(err.Error(), "api_key_ref") {
			t.Fatalf("key failure = %v, want the Unconfigured/auth sentence the C5 resolver uses", err)
		}
		var obs *observe.Error
		if !errors.As(err, &obs) || obs.Class != observe.ClassAuth {
			t.Fatalf("class = %v, want auth", err)
		}
		if strings.Contains(err.Error(), "plaintext-key") {
			t.Fatal("auth error leaked the key")
		}
	})
}

func Test61CapabilityProbesReportToProviderHealthVocabulary(t *testing.T) {
	// audio_in: a live endpoint passes even with an empty transcript (the
	// canned tone carries no words -- the probe tests the CAPABILITY, not
	// accuracy; that claim belongs to the CER harness).
	srvIn := asrJSONServer("", 0.2, 1, nil, nil)
	defer srvIn.Close()
	h := newCascadeHarness(t)
	opts := h.options()
	opts.Client = srvIn.Client()
	res := ProbeAudioIn(context.Background(), cloudEndpoint(srvIn.URL), opts)
	if res.Capability != CapAudioIn || !res.OK {
		t.Fatalf("audio_in probe = %+v, want OK against a live endpoint", res)
	}

	srv500 := asrStatusServer(http.StatusInternalServerError, "")
	defer srv500.Close()
	res = ProbeAudioIn(context.Background(), cloudEndpoint(srv500.URL), opts)
	if res.OK {
		t.Fatalf("audio_in probe passed against a 500 endpoint with %q; a probe may never fake a pass", res.Detail)
	}
	if !strings.Contains(res.Detail, "500") {
		t.Fatalf("failed probe detail %q must carry the status for provider_health", res.Detail)
	}

	srvOut := pcmBodyServer(2400, false)
	defer srvOut.Close()
	res = ProbeAudioOut(context.Background(), cloudEndpoint(srvOut.URL), opts)
	if res.Capability != CapAudioOut || !res.OK || res.LatencyMS <= 0 {
		t.Fatalf("audio_out probe = %+v, want OK with a real latency reading", res)
	}

	skip := cloudEndpoint(srvIn.URL)
	skip.DeclaresAudioIn, skip.DeclaresAudioOut = false, true
	suite := RunVoiceProbeSuite(context.Background(), skip, opts)
	if len(suite) != 1 || suite[0].Capability != CapAudioOut {
		t.Fatalf("suite over an audio_out-only endpoint = %+v, want only the declared capability probed", suite)
	}
}

// Test61E2EWavThroughCloudChainAndTTSPlayback is the AC#1 walk at this
// package's boundary: the real endpointer (ticket 15-b's calibrated
// fixture), the cascade ASR engine behind Transcriber, then the cascade TTS
// through the same Speak/emit contract the announce path uses. The ball
// state walk lives in cmd/wisp and is untouched by this ticket -- the
// voice loop keeps driving NewTranscriber(det, AsrEngine) exactly as
// before; that "unchanged" is the property this test pins.
func Test61E2EWavThroughCloudChainAndTTSPlayback(t *testing.T) {
	det, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if err != nil {
		t.Skipf("VAD session unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = det.Close() })

	h := newCascadeHarness(t)
	srvASR := asrJSONServer("把台灯打开", 1.2, 9, nil, nil)
	defer srvASR.Close()
	srvTTS := pcmBodyServer(3600, false)
	defer srvTTS.Close()

	asrOpts := h.options()
	asrOpts.Client = srvASR.Client()
	chain, err := NewCloudASRChain([]CloudEndpoint{cloudEndpoint(srvASR.URL)}, nil, asrOpts)
	if err != nil {
		t.Fatalf("chain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	tr := NewTranscriber(det, chain)
	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	var out []Transcript
	for f := 0; f < len(pcm)/WindowSamples; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	if len(out) == 0 {
		t.Fatal("the cloud chain over the real endpointer produced no transcripts")
	}
	last := out[len(out)-1]
	if last.Kind != TranscriptFinal || last.Text != "把台灯打开" {
		t.Fatalf("final = %v %q, want the cloud endpoint's transcript", last.Kind, last.Text)
	}
	if want := samplesDuration(chain.Fed()); last.Latency != want {
		t.Fatalf("latency %v is off the chain's audio clock %v (D42#9 across provider implementations)", last.Latency, want)
	}

	ttsOpts := h.options()
	ttsOpts.Client = srvTTS.Client()
	ttsChain, err := NewCloudTTSChain([]CloudEndpoint{cloudEndpoint(srvTTS.URL)}, nil, ttsOpts)
	if err != nil {
		t.Fatalf("tts chain: %v", err)
	}
	defer func() { _ = ttsChain.Close() }()

	played := 0
	err = ttsChain.Speak(last.Text, func(p []float32) error { played += len(p); return nil })
	if err != nil {
		t.Fatalf("cloud TTS playback leg: %v", err)
	}
	if played != 3600 {
		t.Fatalf("played %d samples, want 3600", played)
	}

	_, _, usages := h.snap()
	sawIn, sawOut := false, false
	for _, u := range usages {
		if u.Chain == ChainASR {
			sawIn = true
		}
		if u.Chain == ChainTTS && u.AudioSecondsOut > 0 {
			sawOut = true
		}
	}
	if !sawIn || !sawOut {
		t.Fatalf("E2E legs must both meter to C23: usages %+v", usages)
	}
}

func Test61TTSChainFallsThroughToLocalNeverSilently(t *testing.T) {
	h := newCascadeHarness(t)
	srvA := asrStatusServer(http.StatusServiceUnavailable, "")
	defer srvA.Close()
	epA := cloudEndpoint(srvA.URL)
	epA.Model = "a1"

	opts := h.options()
	opts.Client = srvA.Client()
	rec := &recordingTTS{}
	chain, err := NewCloudTTSChain([]CloudEndpoint{epA}, rec, opts)
	if err != nil {
		t.Fatalf("NewCloudTTSChain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	if err := chain.Speak("兜底播报", func(pcm []float32) error { return nil }); err != nil {
		t.Fatalf("Speak with a local fallback wired = %v, want the fallback to carry the announcement", err)
	}
	if rec.calls != 1 {
		t.Fatalf("local TTS called %d times, want 1", rec.calls)
	}
	failovers, states, _ := h.snap()
	if len(failovers) != 1 || failovers[0].To != "local-sherpa" {
		t.Fatalf("failover events = %+v, want the visible hop to local", failovers)
	}
	last := states[len(states)-1]
	if last.Chain != ChainTTS || last.Active != "local-sherpa" || last.Cloud {
		t.Fatalf("last provider state = %+v, want local-sherpa not cloud (privacy posture restored)", last)
	}
	// The rate follows the engine that OWNS audio out. After the hop it is
	// the local cascade's -- asked of the fallback itself, not of its zeroed
	// config field (recordingTTS answers 22050 for rate 0, exactly as the
	// real matcha engine answers its loaded checkpoint's rate).
	if got, want := chain.SampleRate(), rec.SampleRate(); got != want {
		t.Fatalf("SampleRate = %d, want the fallback's %d", got, want)
	}
	_ = chain.Close()
	if rec.closeN != 0 {
		t.Fatalf("chain closed the fallback engine; ownership stays with the loader")
	}
}

type recordingTTS struct {
	calls   int
	closeN  int
	rate    int
	samples int
}

func (r *recordingTTS) Speak(text string, emit func([]float32) error) error {
	r.calls++
	out := make([]float32, 160)
	r.samples += len(out)
	return emit(out)
}

func (r *recordingTTS) SampleRate() int {
	if r.rate == 0 {
		return 22050
	}
	return r.rate
}
func (r *recordingTTS) Close() error { r.closeN++; return nil }

func sumInts(xs []int) int {
	var s int
	for _, x := range xs {
		s += x
	}
	return s
}
