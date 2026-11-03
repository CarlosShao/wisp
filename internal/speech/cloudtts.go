package speech

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// This file is the C9 cloud TTS engine of ticket 61 -- the twin of
// cloudasr.go on the synthesis side, and the shipping default per ledger R5
// ("默认 TTS 走云 TTS，供应商 StepFun"): the local matcha engine in tts.go
// is a measurement-grade reference the owner's ruling keeps OFF the default
// path; this engine is the path.
//
// Unlike the ASR leg, Speak needs NO goroutine of its own. The seam already
// streams: chunks leave through emit in playback order, and the caller's
// cooperative abort (emit returning an error, barge-in's channel) is the
// same contract the local engine honors through sherpa's progress callback.
// So Speak does one HTTP round trip on the OWNER goroutine and hands out
// chunks as bytes arrive -- a goroutine here would only add a queue between
// the socket and the callback for no latency and no concurrency.
//
// Wire dialect (OpenAI-compatible audio/speech, which the StepFun stepaudio
// endpoints speak -- same base URL family as the C5 text leg):
//
//	POST {base}/audio/speech
//	{"model": <m>, "input": <text>, "voice": <v>, "response_format": "pcm",
//	 "speed": <s>, "stream": true}
//	Accept: application/x-ndjson | text/event-stream | audio/L16
//
// Response, three accepted shapes, tried by content type:
//   - text/plain / application/x-ndjson: one JSON object per line, audio in
//     base64 under "audio_b64" (the shape that carries a provider-side
//     first-audio event for a low-latency cascade);
//   - text/event-stream: SSE "data:" lines, same object shape;
//   - audio/*: raw little-endian s16le PCM for the whole body (the plain
//     OpenAI audio/speech answer with response_format=pcm).
//
// PCM is the required wire format on purpose: mp3/wav would need a decoder,
// and the dependency whitelist (HTTP streaming only, no new native deps)
// rules one out. The first-frame reading the ticket demands is measured
// from request start to first emitted chunk, monotonic, carried on the
// usage record as "provider-dependent" -- there is deliberately NO SLO
// gate in this file.

// firstFrame is measured from request start to first emitted chunk,
// monotonic, carried for diagnostics only: the ticket says "measure, mark
// provider-dependent" and explicitly forbids inventing an SLO gate here,
// so nothing in this file compares against SPEC-04's 800/1600ms budgets.

// ttsRequestBudget bounds one synthesis request end to end. Synthesis output
// arrives as a stream, so the budget covers connect + drain; blowing it
// classifies as network and the chain steps (or falls to the local cascade).
const ttsRequestBudget = 60 * time.Second

// ttsChunkSamples is the emit unit: 2048 samples at the endpoint's rate, so
// a 24 kHz voice hands ~85ms of audio per chunk -- small enough for the
// render side's latency, large enough that per-chunk base64 JSON does not
// dominate the pipe. The seam documents chunking as the engine's choice.
const ttsChunkSamples = 2048

// CloudTTS is one cloud synthesis engine behind the C9 TtsEngine seam.
// Like MatchaTts it starts no goroutines and is driven by one owner
// goroutine (D38b); Speak may be called from any single goroutine in
// sequence but not concurrently.
type CloudTTS struct {
	ep     CloudEndpoint
	client HTTPDoer
	cb     ChainCallbacks

	mu      sync.Mutex
	closed  bool
	samples int // emitted samples of the most recent Speak (audio clock)
	first   time.Duration
}

// The compile-time conformance assertion, twin of the local one in tts.go:
// if this engine ever stops matching the seam, the package stops building.
var _ TtsEngine = (*CloudTTS)(nil)

// NewCloudTTS builds a synthesis engine for one resolved endpoint. It loads
// nothing (there is nothing to load), so construction cannot fail except on
// a missing base URL, which is a config error by name.
func NewCloudTTS(ep CloudEndpoint, opts ChainOptions) (*CloudTTS, error) {
	if ep.BaseURL == "" {
		return nil, observe.New(observe.ClassConfig,
			"speech: cloud TTS endpoint "+ep.String()+" has no base URL")
	}
	return &CloudTTS{
		ep:     ep,
		client: opts.resolveClient(),
		cb:     opts.Callbacks,
	}, nil
}

// Speak synthesizes text and emits its audio in playback order, starting
// with the first chunk as soon as one exists.
//
// The cooperative-abort contract is the seam's: an emit error stops the
// stream immediately and is returned VERBATIM (that is how Path C's
// barge-in stops a sentence mid-flight), with the HTTP body closed on the
// way out. Failure classes reach the caller as *observe.Error exactly like
// the C5 adapters (SPEC-05 sec 3: errors classified at the seam), and the
// chain stepper reads them the same way the text chain does.
func (t *CloudTTS) Speak(text string, emit func(samples []float32) error) error {
	if emit == nil {
		return errors.New("speech: Speak needs an emit callback")
	}
	t.mu.Lock()
	closed := t.closed
	t.mu.Unlock()
	if closed {
		return errors.New("speech: Speak on a closed cloud TTS engine")
	}
	if !speakable(text) {
		// The same guard the local engine runs, same error: an empty or
		// punctuation-only utterance must not become a silent blip billed
		// as speech (ErrNothingToSpeak's whole argument applies to the cloud
		// leg verbatim).
		return ErrNothingToSpeak
	}

	payload, err := json.Marshal(ttsWireRequest{
		Model:          t.ep.Model,
		Input:          text,
		Voice:          t.ep.Voice,
		ResponseFormat: "pcm",
		Speed:          t.ep.ttsSpeed(),
		Stream:         true,
	})
	if err != nil {
		return observe.Wrap(observe.ClassInternal, err, "speech: encode cloud TTS request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ttsRequestBudget)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(t.ep.BaseURL, "/")+"/audio/speech", bytes.NewReader(payload))
	if err != nil {
		return observe.Wrap(observe.ClassInternal, err, "speech: build cloud TTS request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream, application/x-ndjson, audio/L16")
	t.ep.applyAuth(req)

	start := time.Now() // monotonic reading; used only for elapsed math (D42#9)
	resp, err := t.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return llm.WithCode(
				observe.New(observe.ClassNetwork, "speech: cloud TTS request timed out"),
				llm.CodeTimeout)
		}
		return llm.ClassifyTransportError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return classifyCascadeHTTP(resp)
	}

	rate := t.ep.ttsSampleRate()

	// Delivered-count tracking lives on the engine and updates UNDER the
	// consume closure, not after it: the chain's "delivered" rule (no
	// failover once audio reached the consumer, llm/chain.go's rule for
	// text) asks emittedSamples() precisely when Speak is RETURNING AN
	// ERROR, and a counter that only moved on success would let a
	// half-spoken sentence restart from the next provider.
	t.mu.Lock()
	t.samples = 0
	t.first = 0
	t.mu.Unlock()

	var (
		samples int
		// Coarse local estimate the C5 seam uses for the same purpose
		// (llm/provider.go EstimateTokens: utf-8 bytes / 4); the provider's
		// usage block replaces it when the wire carries one.
		inTok    = (len(text) + 3) / 4
		outTok   int
		sawUsage bool
	)
	consume := func(pcm []float32) error {
		samples += len(pcm)
		t.mu.Lock()
		if t.samples == 0 && samples > 0 {
			// A measured reading, never a zero. A localhost or in-process
			// provider can answer faster than the clock's own resolution on
			// a loaded machine, and "0" would read downstream as "no
			// reading at all" -- the same distinction provider_health's
			// latency column makes for the probes (see probeLatencyMS).
			t.first = time.Since(start)
			if t.first <= 0 {
				t.first = time.Nanosecond
			}
		}
		t.samples = samples
		t.mu.Unlock()
		return emit(pcm)
	}

	ct := resp.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "text/event-stream") || strings.HasPrefix(ct, "application/x-ndjson") ||
		strings.HasPrefix(ct, "text/plain"):
		inTok, outTok, sawUsage, err = t.consumeTTSDeltas(resp.Body, ct, consume)
	default:
		// Raw PCM body (audio/* or anything opaque): stream it through in
		// fixed-size chunks. This is the plain OpenAI audio/speech answer.
		err = t.consumeRawPCM(resp.Body, consume)
	}
	if err != nil {
		return err
	}

	if t.cb.OnUsage != nil {
		u := VoiceUsage{
			Chain:           ChainTTS,
			Key:             ChainKey(t.ep.Provider, t.ep.Model),
			Billing:         t.ep.Billing,
			AudioSecondsOut: float64(samples) / float64(rate),
			PriceAudioIn:    t.ep.PriceAudioIn,
			PriceAudioOut:   t.ep.PriceAudioOut,
		}
		if sawUsage {
			u.InputTokens = inTok
			u.OutputTokens = outTok
		}
		t.cb.OnUsage.RecordVoiceUsage(u)
	}
	return nil
}

// SampleRate is the rate of the samples Speak emits. It is a property of
// the endpoint's declared output format (the catalog/wire convention),
// which is why it can be answered before the first request -- the render
// side resamples from it exactly as it does for the local engine's 22050.
func (t *CloudTTS) SampleRate() int { return t.ep.ttsSampleRate() }

// Endpoint names the target this engine synthesizes against (chain
// diagnostics; mirrors CloudASR.Endpoint).
func (t *CloudTTS) Endpoint() CloudEndpoint { return t.ep }

// FirstFrame reports how long the most recent Speak waited for its first
// chunk (monotonic, measured, never gated -- the ticket's "measure, mark
// provider-dependent" reading). Zero before the first successful Speak.
func (t *CloudTTS) FirstFrame() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.first
}

// Close releases the engine (C11). There is nothing resident to release --
// the promise is the interface, and it is kept by being safe to call twice.
func (t *CloudTTS) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	return nil
}

// ttsWireRequest is the audio/speech body. Speed 0 is normalized by the
// endpoint method, and Voice empty means "provider default" (the field is
// omitted so the server's own default applies).
type ttsWireRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice,omitempty"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed,omitempty"`
	Stream         bool    `json:"stream"`
}

// ttsWireChunk is one delta line/object of the streamed response: audio as
// base64 little-endian s16le, plus the optional usage/finish fields of the
// final event.
type ttsWireChunk struct {
	AudioB64 string `json:"audio_b64"`
	IsFinal  bool   `json:"is_final"`
	Usage    *struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *asrWireError `json:"error"`
}

// consumeTTSDeltas reads the streaming shapes line by line -- SSE ("data:"
// prefixed) and NDJSON (bare JSON) share the parser -- converting each base64
// audio delta through the same emit path the raw-PCM shape uses. A chunk
// error aborts; an emit error aborts and is returned verbatim (the seam's
// cooperative-abort contract).
func (t *CloudTTS) consumeTTSDeltas(body io.Reader, contentType string, consume func([]float32) error) (inTok, outTok int, sawUsage bool, err error) {
	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64<<10), 8<<20)
	sse := strings.HasPrefix(contentType, "text/event-stream")
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if sse {
			if !strings.HasPrefix(line, "data:") {
				continue // event:/id:/retry: lines carry nothing for us
			}
			line = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if line == "[DONE]" {
				break
			}
		}
		var ch ttsWireChunk
		if uerr := json.Unmarshal([]byte(line), &ch); uerr != nil {
			return 0, 0, false, observe.Wrap(observe.ClassProvider, uerr,
				"speech: cloud TTS stream chunk is not valid JSON")
		}
		if ch.Error != nil {
			return 0, 0, false, ch.Error.classify()
		}
		if ch.AudioB64 != "" {
			raw, derr := base64.StdEncoding.DecodeString(ch.AudioB64)
			if derr != nil {
				return 0, 0, false, observe.Wrap(observe.ClassProvider, derr,
					"speech: cloud TTS audio delta is not valid base64")
			}
			if cerr := consume(pcm16ToFloat(raw)); cerr != nil {
				return inTok, outTok, sawUsage, cerr // verbatim: barge-in is not a failure
			}
		}
		if ch.Usage != nil {
			inTok, outTok, sawUsage = ch.Usage.InputTokens, ch.Usage.OutputTokens, true
		}
		if ch.IsFinal {
			break
		}
	}
	if serr := sc.Err(); serr != nil {
		return 0, 0, false, llm.ClassifyTransportError(serr)
	}
	return inTok, outTok, sawUsage, nil
}

// consumeRawPCM streams a raw little-endian s16le body through consume in
// fixed-size chunks, so a whole-clip provider still satisfies the seam's
// chunked contract (the same argument engines.go makes for the seam's
// shape: a whole-clip implementation must not be able to hand one giant
// buffer to the render side).
//
// A body that stopped early is NOT a clean end: a truncated synthesis is the
// D37 transport failure, and that is what the chain's "delivered" rule reads
// to decide between stepping and stopping.
//
// The read loop hands net/http a CONTINUOUS reader, which is what makes the
// discrimination the stdlib already knows how to make usable here: on a
// whole-body read, io.ErrUnexpectedEOF means exactly "fewer bytes than the
// Content-Length promised" (including a body that promised 8KB and wrote
// nothing), while a chunked honest end simply returns io.EOF and is copied
// out without error. io.ReadFull must NOT wrap this loop: it reports a short
// FINAL chunk of an honest body as io.ErrUnexpectedEOF too, which would
// condemn every healthy provider. The chunk boundaries are imposed by the
// emit size instead, which is what the seam's chunked contract asks for.
//
// What neither signal can repair is an exchange that delivered NO audio at
// all: a 200 with an empty body is a provider that answered without
// speaking, and calling that a success would let a silent endpoint pass the
// chain's delivered test with zero samples played.
func (t *CloudTTS) consumeRawPCM(body io.Reader, consume func([]float32) error) error {
	const chunkBytes = ttsChunkSamples * 2
	buf := make([]byte, chunkBytes)
	pending := make([]byte, 0, chunkBytes*2)
	delivered := false

	emitChunks := func(final bool) error {
		for len(pending) >= chunkBytes || (final && len(pending) > 0) {
			n := min(len(pending), chunkBytes)
			if err := consume(pcm16ToFloat(pending[:n])); err != nil {
				return err
			}
			delivered = true
			pending = append(pending[:0], pending[n:]...)
		}
		return nil
	}

	for {
		n, rerr := body.Read(buf)
		if n > 0 {
			pending = append(pending, buf[:n]...)
			if err := emitChunks(false); err != nil {
				return err
			}
		}
		if rerr != nil {
			if err := emitChunks(true); err != nil {
				return err
			}
			switch {
			case rerr == io.EOF:
				if !delivered {
					return observe.New(observe.ClassProvider,
						"speech: cloud TTS endpoint answered 200 with an empty audio body")
				}
				return nil
			case rerr == io.ErrUnexpectedEOF:
				return llm.ClassifyTransportError(rerr)
			default:
				return llm.ClassifyTransportError(rerr)
			}
		}
	}
}

// pcm16ToFloat converts little-endian s16le bytes to float32 in -1..1.
// Odd trailing bytes are dropped -- a half sample is not audio, and
// inventing one from the high byte would speak a click.
func pcm16ToFloat(raw []byte) []float32 {
	n := len(raw) / 2
	out := make([]float32, n)
	for i := 0; i < n; i++ {
		out[i] = float32(int16(binary.LittleEndian.Uint16(raw[i*2:]))) / 32768.0
	}
	return out
}
