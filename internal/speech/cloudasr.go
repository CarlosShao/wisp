package speech

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// This file is the C9 cloud ASR engine of ticket 61: recognition against
// an OpenAI-compatible-style audio endpoint, written on the SAME seam 15-c
// froze (Begin/Feed/Partial/Finish/InFlight/Fed/Close) so Transcriber
// cannot tell it from *StreamingAsr -- that indistinguishability is C9's
// entire point, and engines_test's stub said "replacing this stub with a
// real provider is ticket 61's work and changes nothing above it".
//
// Three shapes differ from the local engine, each for a stated reason:
//
//   - Partials are produced by PERIODIC SNAPSHOT REQUESTS from one spawned
//     loop per utterance (observe.Registry.Spawn; the roster temp name
//     asr-infer, owner "speech"), not by a duplex stream. Measured on this
//     host, not theorized: Go's net/http HTTP/1.1 transport does NOT
//     reliably deliver a chunked request body to a server that responds
//     before the upload completes -- the early response ends the round trip
//     while body bytes are still buffered (the design probe that proved it
//     is in this ticket's commit notes). A "one request, stream audio up,
//     stream partials down" engine would work against some vendors and
//     silently starve others. The cascade's honest REST/SSE-only shape is:
//     every snapshot cadence of the AUDIO CLOCK, POST what has been heard
//     so far and take the hypothesis from the reply; Finish POSTs the whole
//     utterance as the authoritative final. Each such call is exactly one
//     OpenAI-compatible /audio/transcriptions request (the StepFun
//     step-1-asr family), the VAD->ASR pipeline gets real streaming
//     partials (SPEC-04 sec 4), and everything rides the
//     response-after-complete-body pattern the transport is verified
//     against.
//   - The owner thread never touches the socket: Feed appends bytes under a
//     lock and signals a capacity-1 wake channel, so a slow endpoint cannot
//     stall the 32ms capture cadence, and merged wake-ups mean a busy loop
//     skips stale snapshots instead of queueing them.
//   - Finish blocks until the final reply lands -- bounded by monotonic
//     budgets (context deadlines plus observe.NewTimeout; the D22
//     wall-clock ban). A blown budget is a network failure with a class,
//     not a hang. The failure is recorded beside the empty final
//     (LastFailure): the seam stays two-valued -- cloud failure reads like
//     a mishear to its caller, exactly what AC#5's retry policy keys on --
//     while the cascade (cloudchain.go) reads the class and steps.
//
// Wire dialect (first candidates: StepAudio 2.5 ASR via StepFun's
// OpenAI-compatible audio endpoints). Every snapshot and the final are:
//
//	POST {base}/audio/transcriptions?model=<m>[&language=<l>]
//	Content-Type: application/octet-stream   (body: s16le mono 16 kHz PCM)
//	Accept: application/json, text/event-stream
//
// The reply is {"text": <hypothesis>, "usage": {"seconds","input_tokens"}};
// a provider that answers with an SSE stream instead is read to its end
// and its LAST data event wins (is_final honoured when present). The
// request builder is one function (buildASRRequest) so a later vendor's
// dialect drops in behind it -- the ticket's "provider abstraction keeps
// any later vendor drop-in".

// cloudASROwner is the observe.Registry owner label for this engine's
// transport goroutine. The roster NAME is "asr-infer", the frozen temporary
// name (internal/observe/goroutine.go) for ASR inference work inside a voice
// scope -- this goroutine is exactly that work, moved off-box.
const cloudASROwner = "speech"

// utteranceFinishBudget bounds how long Finish waits for the final reply
// after the utterance audio is fully in. It is a transport watchdog (a hung
// connection must become a class-network failure the chain steps on), NOT a
// re-stated SLO: SPEC-04 sec 8's 800ms streaming budget is the measurement
// target the ball holds the chain to, and the ticket forbids inventing a
// gate here for a cold-network cascade.
const utteranceFinishBudget = 15 * time.Second

// utteranceRequestBudget bounds one transcription request end to end.
// 16k mono s16le is 32 KB per second of speech, so even a two-minute
// utterance sends ~4 MB well inside this; a budget hit is a network-side
// failure by class, and the chain then falls to the local cascade -- the
// ticket's terminal fallback is a correctness guarantee, not a preference.
const utteranceRequestBudget = 60 * time.Second

// asrSnapshotCadenceSeconds is the audio-clock period of partial snapshot
// requests. Two seconds: long enough that a four-seconds-per-second-of-
// audio fast endpoint (the StepAudio 2.5 ASR claim) never overlaps itself
// at realistic utterance lengths, short enough that the ball has a growing
// hypothesis while the speaker talks. Utterances shorter than the cadence
// simply get the final -- the seam documents Partial as "empty until the
// engine has enough audio to say anything".
const asrSnapshotCadenceSeconds = 2

// asrSnapshotCadenceSamples is the same cadence in samples.
const asrSnapshotCadenceSamples = asrSnapshotCadenceSeconds * audio.TargetRate

// asrUtteranceCapSamples bounds one utterance's buffered audio. At 16k mono
// this is ten minutes -- far past any endpointer utterance; overflow fails
// the utterance visibly (class resource) rather than dropping a window
// silently, because a dropped window would corrupt the transcript and hide
// the corruption.
const asrUtteranceCapSamples = 600 * audio.TargetRate

// CloudASR is one cloud recognition session behind the C9 seam.
//
// It is driven by a single owner goroutine (D38b): Begin/Feed/Partial/
// Finish/Close are called from the voice loop; the one goroutine it spawns
// per utterance owns that utterance's HTTP requests. It is not safe for
// concurrent use, exactly like StreamingAsr.
type CloudASR struct {
	ep     CloudEndpoint
	client HTTPDoer
	reg    *observe.Registry
	root   *observe.Root
	cb     ChainCallbacks

	// ctx is the engine lifetime: derived from opts.Root.Ctx when a root is
	// wired (so the voice loop's cancellation reaches the sockets) and
	// canceled by Close.
	ctx    context.Context
	cancel context.CancelFunc

	mu       sync.Mutex
	sess     *cloudASRSession // non-nil while an utterance is open
	lastErr  *observe.Error   // failure behind the most recent Finish
	closed   bool
	fedCarry int // Fed() of the last closed utterance, until the next Begin
}

// cloudASRSession is one utterance: the sample buffer the owner appends to,
// the snapshot loop's wake channel, and the published outcome.
type cloudASRSession struct {
	inMu     sync.Mutex
	in       []byte // utterance audio so far, s16le, grow-only
	fed      int    // samples fed -- the audio clock (D42#9)
	sent     int    // snapshots already dispatched
	finalReq bool   // Finish has asked for the authoritative final

	wake chan struct{} // capacity-1 signal from Feed/Finish to the loop

	// Outcome side, owned by pubMu; written once (first publisher wins).
	pubMu          sync.Mutex
	partial        string
	final          string
	finalErr       *observe.Error
	publishedFinal bool
	usageSeen      bool
	usageInTokens  int
	usageSeconds   float64
	done           chan struct{}

	handle *observe.Handle
}

// The compile-time conformance C9 asks for: a cloud engine that stops
// matching the seam stops the build, exactly like the local assertion in
// engines.go.
var _ AsrEngine = (*CloudASR)(nil)

// NewCloudASR builds a session for one resolved endpoint. It loads nothing:
// the endpoint's whole cost is HTTP requests per utterance, which is why
// cloud engines are cheap to hold resident while sherpa sessions are not --
// the residency slot (slot.go) stays a local-family concern.
func NewCloudASR(ep CloudEndpoint, opts ChainOptions) (*CloudASR, error) {
	if ep.BaseURL == "" {
		return nil, observe.New(observe.ClassConfig,
			"speech: cloud ASR endpoint "+ep.String()+" has no base URL")
	}
	base := context.Background()
	if opts.Root != nil {
		base = opts.Root.Ctx
	}
	ctx, cancel := context.WithCancel(base)
	return &CloudASR{
		ep:     ep,
		client: opts.resolveClient(),
		reg:    opts.registry(),
		root:   opts.Root,
		cb:     opts.Callbacks,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Begin opens an utterance: a fresh decoder state, so nothing from the
// previous utterance can leak into this one (same clause as StreamingAsr).
// A still-open previous utterance is failed first -- Begin is the owner
// saying "the last one is over".
func (a *CloudASR) Begin() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.closed {
		return
	}
	if a.sess != nil {
		a.fedCarry = a.sess.fedCount()
		a.sess.publish("", observe.New(observe.ClassCancelled,
			"speech: cloud ASR utterance abandoned by a new Begin"))
		a.sess = nil
	}
	s := &cloudASRSession{
		wake: make(chan struct{}, 1),
		done: make(chan struct{}),
	}
	uctx, cancel := context.WithCancel(a.ctx)
	a.sess = s
	a.lastErr = nil

	s.handle = a.reg.Spawn("asr-infer", cloudASROwner, a.root, func(context.Context) {
		a.runUtterance(uctx, s)
		cancel()
	})
}

// Feed appends utterance audio: 16 kHz mono float32, any length. It never
// touches the socket: float32 -> s16le under the input lock, then a
// non-blocking wake. Merged wake-ups are the backpressure policy: a loop
// still busy posting an older snapshot skips stale cadence boundaries
// rather than queueing them, and the next due-check catches up.
func (a *CloudASR) Feed(samples []float32) {
	if len(samples) == 0 {
		return
	}
	a.mu.Lock()
	s := a.sess
	a.mu.Unlock()
	if s == nil {
		return // like StreamingAsr: no utterance open, no-op
	}
	pcm := floatToPCM16(samples)

	s.inMu.Lock()
	if s.fed+len(samples) > asrUtteranceCapSamples {
		s.inMu.Unlock()
		// Overflow is a failure, not a drop: the transcript would be
		// quietly wrong and the seam carries no error channel to say so.
		// Publishing makes Finish return "" with a class, and the chain
		// steps -- loud beats silently wrong, always.
		s.publish("", observe.New(observe.ClassResource, fmt.Sprintf(
			"speech: cloud ASR utterance exceeded the %d-sample cap (%d seconds of 16k audio)",
			asrUtteranceCapSamples, asrUtteranceCapSamples/audio.TargetRate)))
		return
	}
	s.in = append(s.in, pcm...)
	s.fed += len(samples)
	due := s.fed >= (s.sent+1)*asrSnapshotCadenceSamples
	s.inMu.Unlock()

	if due {
		select {
		case s.wake <- struct{}{}:
		default: // loop already awake
		}
	}
}

// Partial is the hypothesis for the open utterance: the last snapshot's
// answer. Empty until the first snapshot replies -- the network's share of
// the latency the ticket says to measure and NOT gate. Inside one utterance
// it only ever grows while the provider's hypotheses do.
func (a *CloudASR) Partial() string {
	a.mu.Lock()
	s := a.sess
	a.mu.Unlock()
	if s == nil {
		return ""
	}
	s.pubMu.Lock()
	defer s.pubMu.Unlock()
	return s.partial
}

// Finish closes the utterance: the whole buffered audio becomes the
// authoritative transcription request, and Finish waits (monotonically) for
// its answer. It never blocks forever: a blown watchdog becomes a
// class-network failure in LastFailure, and the final reads like a mishear
// (empty) to the seam's caller -- the single signal AC#5's retry policy
// keys on, unchanged between engines.
func (a *CloudASR) Finish() string {
	a.mu.Lock()
	s := a.sess
	a.mu.Unlock()
	if s == nil {
		return "" // safe to call with no utterance open (returns ""), as Finish does locally
	}

	s.inMu.Lock()
	s.finalReq = true
	s.inMu.Unlock()
	select {
	case s.wake <- struct{}{}:
	default:
	}

	// The watchdog covers the WORST honest sequence: one in-flight snapshot
	// at its full request budget, then the final at its own, plus the
	// finish slack. The per-request budgets self-cap long before this fires
	// in practice; this is the last line against a wedged transport, and
	// the budget math is monotonic throughout (D42#9).
	budget := observe.NewTimeout(utteranceFinishBudget + 2*utteranceRequestBudget)
	timer := time.NewTimer(budget.Remaining())
	select {
	case <-s.done:
		timer.Stop()
	case <-timer.C:
		s.publish("", llm.WithCode(
			observe.New(observe.ClassNetwork,
				"speech: cloud ASR final reply did not arrive within the finish budget"),
			llm.CodeTimeout))
	}

	s.pubMu.Lock()
	text := s.final
	cause := s.finalErr
	use := asrUsageSnapshot{seen: s.usageSeen, inputTokens: s.usageInTokens, seconds: s.usageSeconds}
	s.pubMu.Unlock()

	fed := s.fedCount()
	a.mu.Lock()
	a.fedCarry = fed
	a.lastErr = cause
	a.sess = nil
	a.mu.Unlock()
	s.stop() // idempotent: publish already stopped the loop when it ran

	if cause != nil {
		return ""
	}
	a.recordUsage(fed, use)
	return text
}

// InFlight reports whether an utterance is open.
func (a *CloudASR) InFlight() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.sess != nil
}

// Fed is how many samples the open -- or, right after Finish, the most
// recently closed -- utterance streamed. It is the audio clock every
// transcript latency in this package is reported on (D42#9), and it counts
// what the CALLER fed, not what the socket managed to send.
func (a *CloudASR) Fed() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.sess != nil {
		return a.sess.fedCount()
	}
	return a.fedCarry
}

// LastFailure is the cascade's private window -- NOT part of the seam: the
// chain stepper asks why a Finish came back empty before it decides to
// step. Callers holding the AsrEngine interface never see this method and
// never branch on it (C9: at the seam a cloud failure is indistinguishable
// from a local mishear, by contract).
func (a *CloudASR) LastFailure() *observe.Error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.lastErr
}

// Endpoint names the target this session recognizes against.
func (a *CloudASR) Endpoint() CloudEndpoint { return a.ep }

// Close releases the engine (C11): it fails any open utterance, cancels the
// sockets, and joins the snapshot loop. A closed engine Begin-refuses
// instead of spawning onto a dead context. Idempotent.
func (a *CloudASR) Close() error {
	a.mu.Lock()
	if a.closed {
		a.mu.Unlock()
		return nil
	}
	a.closed = true
	s := a.sess
	a.sess = nil
	if s != nil {
		a.fedCarry = s.fedCount()
	}
	a.mu.Unlock()

	if s != nil {
		s.publish("", observe.New(observe.ClassCancelled,
			"speech: cloud ASR closed with an open utterance"))
		s.stop()
		if s.handle != nil {
			<-s.handle.Done()
		}
	}
	a.cancel()
	return nil
}

// runUtterance is the snapshot loop: one goroutine owning every HTTP
// request of one utterance. It posts due snapshots as the audio clock runs
// past cadence boundaries, switches to the authoritative final when Finish
// asks, and publishes exactly once. The panic boundary is
// observe.Registry.Spawn's, not this loop's.
func (a *CloudASR) runUtterance(ctx context.Context, s *cloudASRSession) {
	for {
		select {
		case <-ctx.Done():
			s.publish("", observe.New(observe.ClassCancelled, "speech: cloud ASR utterance cancelled"))
			return
		case <-s.wake:
		}

		// The final outranks any pending snapshot: once the utterance
		// closed, its answer IS the transcript.
		s.inMu.Lock()
		wantFinal := s.finalReq
		s.inMu.Unlock()
		if wantFinal {
			pcm := s.snapshotAll()
			text, use, cause := a.transcribe(ctx, pcm)
			if cause != nil {
				s.publish("", cause)
			} else {
				s.publishWithUsage(text, nil, use)
			}
			return
		}

		// Catch up on every due boundary; stale ones collapse (one fresh
		// snapshot of the whole buffer replaces N skipped ones), which is
		// the point of the capacity-1 wake. `sent` records the boundaries
		// the last posted snapshot already covers.
		for {
			s.inMu.Lock()
			due := s.fed >= (s.sent+1)*asrSnapshotCadenceSamples
			fedNow := s.fed
			wantFinal := s.finalReq
			s.inMu.Unlock()
			if !due || wantFinal {
				break
			}
			pcm := s.snapshotAll()
			text, _, cause := a.transcribe(ctx, pcm)
			s.inMu.Lock()
			s.sent = fedNow / asrSnapshotCadenceSamples
			s.inMu.Unlock()
			if cause != nil {
				// A snapshot failure does not fail the utterance: the final
				// attempt is the one whose class decides.
				continue
			}
			if text != "" {
				s.pubMu.Lock()
				s.partial = text
				s.pubMu.Unlock()
			}
		}
	}
}

// transcribe is THE request builder and reader for this engine: one POST of
// the utterance-so-far, one hypothesis out. Every snapshot and the final
// run through it, which is what makes the "later vendor drops in behind
// one function" claim true rather than aspirational.
func (a *CloudASR) transcribe(ctx context.Context, pcm []byte) (string, *asrWireUse, *observe.Error) {
	req, cancel, err := a.buildASRRequest(ctx, pcm)
	if err != nil {
		return "", nil, observe.Wrap(observe.ClassInternal, err, "speech: build cloud ASR request")
	}
	defer cancel()

	resp, err := a.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return "", nil, observe.New(observe.ClassCancelled, "speech: cloud ASR request cancelled")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return "", nil, llm.WithCode(
				observe.New(observe.ClassNetwork, "speech: cloud ASR request exceeded the per-request budget"),
				llm.CodeTimeout)
		}
		return "", nil, llm.ClassifyTransportError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, classifyCascadeHTTP(resp)
	}

	ct := resp.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "application/json"):
		return readASRJSON(resp.Body)
	case strings.HasPrefix(ct, "text/event-stream"):
		return readASRSSE(ctx, resp.Body)
	default:
		return "", nil, observe.New(observe.ClassProvider,
			"speech: cloud ASR endpoint answered with content type "+ct)
	}
}

// buildASRRequest pins the wire: URL, method, headers, auth. A vendor whose
// dialect differs changes THIS FUNCTION and nothing else in the engine.
func (a *CloudASR) buildASRRequest(ctx context.Context, pcm []byte) (*http.Request, context.CancelFunc, error) {
	q := "model=" + a.ep.Model
	if a.ep.Language != "" {
		q += "&language=" + a.ep.Language
	}
	url := strings.TrimRight(a.ep.BaseURL, "/") + "/audio/transcriptions?" + q
	rctx, cancel := context.WithTimeout(ctx, utteranceRequestBudget)
	req, err := http.NewRequestWithContext(rctx, http.MethodPost, url, bytes.NewReader(pcm))
	if err != nil {
		cancel()
		return nil, cancel, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json, text/event-stream")
	a.ep.applyAuth(req)
	return req, cancel, nil
}

// readASRJSON parses the documented single-object reply.
func readASRJSON(body io.Reader) (string, *asrWireUse, *observe.Error) {
	raw, err := io.ReadAll(io.LimitReader(body, 1<<20))
	if err != nil {
		return "", nil, llm.ClassifyTransportError(err)
	}
	var ev asrWireEvent
	if err := json.Unmarshal(raw, &ev); err != nil {
		return "", nil, observe.Wrap(observe.ClassProvider, err,
			"speech: cloud ASR response body is not valid JSON")
	}
	if ev.Error != nil {
		return "", nil, ev.Error.classify()
	}
	return ev.Text, ev.Usage, nil
}

// readASRSSE reads a streaming provider to its end: the LAST data event
// wins (is_final honoured when it appears, an in-band error aborts). A
// stream that produced nothing parseable is a network failure -- an
// endpoint that opened and said nothing did not transcribe.
func readASRSSE(ctx context.Context, body io.Reader) (string, *asrWireUse, *observe.Error) {
	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64<<10), 4<<20)
	var (
		lastText string
		lastUse  *asrWireUse
		sawAny   bool
		finalIn  bool
	)
	for sc.Scan() {
		if ctx.Err() != nil {
			return "", nil, observe.New(observe.ClassCancelled, "speech: cloud ASR stream cancelled mid-read")
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, ":") || !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		var ev asrWireEvent
		if err := json.Unmarshal([]byte(data), &ev); err != nil {
			return "", nil, observe.Wrap(observe.ClassProvider, err,
				"speech: cloud ASR stream event is not valid JSON")
		}
		if ev.Error != nil {
			return "", nil, ev.Error.classify()
		}
		if ev.Text != "" {
			lastText = ev.Text
			sawAny = true
		}
		if ev.Usage != nil {
			lastUse = ev.Usage
			sawAny = true
		}
		if ev.IsFinal {
			finalIn = true
			break
		}
	}
	if err := sc.Err(); err != nil {
		return "", nil, llm.ClassifyTransportError(err)
	}
	if !sawAny && !finalIn {
		return "", nil, observe.New(observe.ClassNetwork,
			"speech: cloud ASR stream ended without a parseable hypothesis")
	}
	return lastText, lastUse, nil
}

// classifyCascadeHTTP turns a non-2xx on either audio endpoint into a D37
// observe.Error using the SAME status/code table the C5 adapters run
// (llm/errors.go). One taxonomy for voice and text failures is what keeps
// the chain-stepping rules ("auth never retried, rate_limit honors
// retry-after, budget means an account out of money") true across both
// chains without a second, drifting copy of the table -- ticket 196's
// vocabulary lesson, applied to error classes.
func classifyCascadeHTTP(resp *http.Response) *observe.Error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	detail := llm.HTTPErrorDetail{Status: resp.StatusCode}
	var probe struct {
		Error asrWireError `json:"error"`
	}
	if json.Unmarshal(body, &probe) == nil && probe.Error.Message != "" {
		detail.Code = probe.Error.codeString()
		detail.Type = probe.Error.Type
		detail.Message = probe.Error.Message
	} else {
		detail.Message = strings.TrimSpace(string(body))
	}
	return llm.NewHTTPError(detail, llm.ParseRetryAfter(resp.Header.Get("Retry-After"), time.Now()))
}

// asrUsageSnapshot is the provider usage block copied out under pubMu for
// metering outside it.
type asrUsageSnapshot struct {
	seen        bool
	inputTokens int
	seconds     float64
}

// recordUsage books one successful recognition against C23 when the host
// wired a sink. The audio seconds come off the sample clock (D42#9); a
// provider-reported seconds number replaces it when present, because that
// is the quantity the vendor's invoice uses.
func (a *CloudASR) recordUsage(fedSamples int, use asrUsageSnapshot) {
	if a.cb.OnUsage == nil {
		return
	}
	u := VoiceUsage{
		Chain:          ChainASR,
		Key:            ChainKey(a.ep.Provider, a.ep.Model),
		Billing:        a.ep.Billing,
		AudioSecondsIn: float64(fedSamples) / float64(audio.TargetRate),
		PriceAudioIn:   a.ep.PriceAudioIn,
		PriceAudioOut:  a.ep.PriceAudioOut,
	}
	if use.seconds > 0 {
		u.AudioSecondsIn = use.seconds
	}
	if use.seen {
		u.InputTokens = use.inputTokens
	}
	a.cb.OnUsage.RecordVoiceUsage(u)
}

// asrWireEvent is one hypothesis object: the JSON reply's shape, and (per
// data line) the SSE variant's.
type asrWireEvent struct {
	Text    string        `json:"text"`
	IsFinal bool          `json:"is_final"`
	Usage   *asrWireUse   `json:"usage"`
	Error   *asrWireError `json:"error"`
}

// asrWireUse is the provider's usage block when it sends one: billed audio
// seconds and the OpenAI-style audio-token counter.
type asrWireUse struct {
	Seconds     float64 `json:"seconds"`
	InputTokens int     `json:"input_tokens"`
}

// asrWireError is the wire's error object (OpenAI-style code/type/message;
// code may be numeric on some compatible servers, same tolerance as the
// openai-chat adapter's).
type asrWireError struct {
	Code    any    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *asrWireError) codeString() string {
	switch v := e.Code.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprint(int64(v))
	default:
		return ""
	}
}

// classify turns an in-band error object into a D37 observe.Error. An
// in-band error carries no HTTP status, so the class defaults to provider
// with the wire code preserved for the chain's decision; the chain steps
// on every class except cancelled, so the walk is the same either way.
func (e *asrWireError) classify() *observe.Error {
	oe := observe.New(observe.ClassProvider, "speech: cloud ASR endpoint error: "+e.Message)
	oe.ProviderCode = e.codeString()
	return oe
}

// -- session helpers ----------------------------------------------------------

// fedCount reads the audio clock under the input lock.
func (s *cloudASRSession) fedCount() int {
	s.inMu.Lock()
	defer s.inMu.Unlock()
	return s.fed
}

// snapshotAll copies the buffered utterance so the loop can POST it without
// holding the input lock across the network.
func (s *cloudASRSession) snapshotAll() []byte {
	s.inMu.Lock()
	defer s.inMu.Unlock()
	out := make([]byte, len(s.in))
	copy(out, s.in)
	return out
}

// publish closes the utterance with a transcript ("" on failure) and
// optionally the classified cause. First writer wins: a deadline-fired
// failure must not be overwritten by a late reply, and a server answer
// must not be overwritten by a Close that raced it.
func (s *cloudASRSession) publish(text string, cause *observe.Error) {
	s.publishWithUsage(text, cause, nil)
}

// publishWithUsage is publish plus the provider's usage numbers, written
// under the same lock so Finish's snapshot is atomic.
func (s *cloudASRSession) publishWithUsage(text string, cause *observe.Error, use *asrWireUse) {
	s.pubMu.Lock()
	if !s.publishedFinal {
		if cause != nil {
			s.finalErr = cause
		} else {
			s.final = text
			if use != nil {
				s.usageSeen = true
				s.usageInTokens = use.InputTokens
				s.usageSeconds = use.Seconds
			}
		}
		s.publishedFinal = true
		close(s.done)
	}
	s.pubMu.Unlock()
}

// stop wakes a parked loop and lets it retire; the loop's deferred cancel
// kills any in-flight request. Called by Finish/Close after the outcome was
// published (or raced to it), and safe more than once.
func (s *cloudASRSession) stop() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

// floatToPCM16 converts 16k mono float32 (-1..1) to little-endian s16le.
// Values outside the range clamp instead of wrapping -- a wrapped sample
// reads as a different word; a clamped one at least reads as a louder
// syllable of the right word.
func floatToPCM16(samples []float32) []byte {
	out := make([]byte, len(samples)*2)
	for i, v := range samples {
		f := v * 32768.0
		switch {
		case f > math.MaxInt16:
			f = math.MaxInt16
		case f < math.MinInt16:
			f = math.MinInt16
		}
		binary.LittleEndian.PutUint16(out[i*2:], uint16(int16(f)))
	}
	return out
}
