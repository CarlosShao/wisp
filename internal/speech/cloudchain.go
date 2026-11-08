package speech

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
)

// This file is the voice cascade itself: CloudASRChain and CloudTTSChain,
// the two AsrEngine/TtsEngine implementations that hide the chain walk from
// the caller. The voice loop is wired against exactly the same seams as
// before ticket 61 -- NewTranscriber(det, chain) is the whole integration
// surface, and the announce path calls chain.Speak(text, emit) like any
// other TtsEngine. Interchangeability is not a slogan here: the chain is an
// engine, so "cloud or local, who knows" is enforced by the type.
//
// The walk rules are SPEC-03 sec 3.1 plus the ticket's Key constraints, and
// each one is load-bearing for a named reason:
//
//   - The local sherpa cascade is the TERMINAL fallback and never a chain
//     member (the config validator rejects "local-sherpa" elements). When
//     the chain exhausts, it falls to local with a visible event -- never
//     silent: a silent cloud-to-local switch would make the privacy
//     indicator lie and the failover counters un-auditable.
//   - Quota is checked BEFORE dispatch (ticket 44's anti-negative rule):
//     an over-quota entry is skipped without ever sending audio, and the
//     skip is a visible event of its own.
//   - ASR failover REPLAYS the buffered utterance into the next element.
//     The endpointer already closed the utterance; an engine that never
//     heard the audio cannot produce a transcript, so the chain owns one
//     replay buffer (capped, and a truncation is said out loud).
//   - TTS failover follows text_chain's "delivered" rule (llm/chain.go):
//     once audio reached the consumer, a second element would append a
//     fresh sentence to a half-spoken one -- the stream ends with the
//     error instead. Before the first chunk, however anything goes.
//   - Cancellation is never a failover cause: it is the user (or the
//     state machine) stopping the voice loop, and the walk must not
//     "retry" that across providers.

// headNone and headLocal are the chain's pseudo-positions: "nothing active"
// (every element quota-skipped with no fallback wired) and "the terminal
// local cascade is active". Negative so element indices stay comparable.
const (
	headNone  = -1
	headLocal = -2
)

// replayCapSeconds bounds the ASR chain's utterance replay buffer. The
// endpointer closes utterances far below this in practice; a pathological
// non-closing speaker is truncated loudly rather than unbounded.
const replayCapSeconds = 120

// chainReplayCapSamples is replayCapSeconds at the seam rate -- the sample
// ceiling of one buffered utterance (~7 MB as float32; transient).
const chainReplayCapSamples = replayCapSeconds * audio.TargetRate

// CloudASRChain is the AsrEngine that walks voice.cloud_asr_chain and falls
// through to the local sherpa cascade. Like every engine in this package it
// is driven by one owner goroutine (D38b) and is not safe for concurrent
// use; the only goroutines anywhere in the cascade are the per-utterance
// transports CloudASR itself spawns.
type CloudASRChain struct {
	engines  []*CloudASR
	names    []string // "provider/model" per element (events/diagnostics)
	fallback AsrEngine

	cb ChainCallbacks

	head       int  // current element index, headNone, or headLocal
	open       bool // between Begin and Finish
	fed        int  // audio clock of the open/last utterance (chain-owned)
	buf        []float32
	bufDropped bool // replay buffer hit its cap this utterance
	lastErr    *observe.Error
}

// CloudTTSChain is the TtsEngine twin: it walks voice.cloud_tts_chain, falls
// through to the local cascade, and reports the same event vocabulary.
type CloudTTSChain struct {
	engines  []*CloudTTS
	names    []string
	fallback TtsEngine

	cb ChainCallbacks

	active int  // head that served the most recent Speak (SampleRate follows it)
	cloud  bool // whether that head takes audio off the machine
}

// The compile-time C9 assertions for the chains themselves: they ARE
// engines, so the voice loop needs no chain-awareness at all.
var (
	_ AsrEngine = (*CloudASRChain)(nil)
	_ TtsEngine = (*CloudTTSChain)(nil)
)

// NewCloudASRChain builds the voice.cloud_asr_chain cascade: one CloudASR
// per resolved endpoint (in chain order), the local sherpa engine as the
// terminal fallback. fallback may be nil -- then chain exhaustion is the
// end of the line (the terminal event names an empty To), and the
// composition root is choosing cloud-or-nothing.
//
// The chain owns the engines it built (Close releases them); it does NOT
// own the fallback: that engine belongs to whoever loaded it (the slot in
// slot.go / the composition root), and closing someone else's resident
// recognizer from inside a chain would be a surprise that outlives the
// chain.
func NewCloudASRChain(eps []CloudEndpoint, fallback AsrEngine, opts ChainOptions) (*CloudASRChain, error) {
	if len(eps) == 0 && fallback == nil {
		return nil, observe.New(observe.ClassConfig,
			"speech: cloud ASR chain with no endpoints and no local fallback would swallow every utterance")
	}
	engines := make([]*CloudASR, 0, len(eps))
	names := make([]string, 0, len(eps))
	for _, ep := range eps {
		e, err := NewCloudASR(ep, opts)
		if err != nil {
			// Release what was already built: a half-built chain must not
			// leak sockets.
			for _, built := range engines {
				_ = built.Close()
			}
			return nil, err
		}
		engines = append(engines, e)
		names = append(names, ChainKey(ep.Provider, ep.Model))
	}
	c := &CloudASRChain{
		engines:  engines,
		names:    names,
		fallback: fallback,
		cb:       opts.Callbacks,
	}
	c.head = initialHead(len(engines), fallback != nil)
	fireProviderState(c.cb.OnProviderState, ChainASR, c.currentName(), c.isCloudHead())
	return c, nil
}

// NewCloudTTSChain builds the voice.cloud_tts_chain cascade (same ownership
// rules as NewCloudASRChain).
func NewCloudTTSChain(eps []CloudEndpoint, fallback TtsEngine, opts ChainOptions) (*CloudTTSChain, error) {
	if len(eps) == 0 && fallback == nil {
		return nil, observe.New(observe.ClassConfig,
			"speech: cloud TTS chain with no endpoints and no local fallback would swallow every announcement")
	}
	engines := make([]*CloudTTS, 0, len(eps))
	names := make([]string, 0, len(eps))
	for _, ep := range eps {
		e, err := NewCloudTTS(ep, opts)
		if err != nil {
			for _, built := range engines {
				_ = built.Close()
			}
			return nil, err
		}
		engines = append(engines, e)
		names = append(names, ChainKey(ep.Provider, ep.Model))
	}
	c := &CloudTTSChain{
		engines:  engines,
		names:    names,
		fallback: fallback,
		cb:       opts.Callbacks,
	}
	c.active, c.cloud = initialTTSHead(len(engines), fallback != nil)
	fireProviderState(c.cb.OnProviderState, ChainTTS, c.Active(), c.cloud)
	return c, nil
}

// -- ASR chain: the seam, walked ----------------------------------------------

// Begin opens an utterance on the first quota-allowed element. Quota skips
// fire visible events BEFORE any audio dispatches (AC#3's "chain advances
// BEFORE dispatch, no negative, event visible").
func (c *CloudASRChain) Begin() {
	c.open = true
	c.fed = 0
	c.buf = c.buf[:0]
	c.bufDropped = false
	c.lastErr = nil

	idx := c.selectHead(0)
	c.setHead(idx)
	c.beginAt(idx)
	if idx == headNone && c.fallback != nil {
		// Every element was over quota. The loud half already happened: the
		// per-element skip events fired inside selectHead, and if the last
		// element could see the local cascade it named it as the To. The
		// SEAM must still carry the utterance -- falling to silence here
		// would drop the caller's audio without a transcript, which is the
		// "silently" the terminal-fallback rule forbids. Landing on local
		// is announced exactly when no skip event has already named it.
		c.setHead(headLocal)
		c.beginAt(headLocal)
	}
}

// Feed forwards utterance audio to the active head and keeps the chain's own
// replay copy. Fed counts on the audio clock at the chain level, so a
// provider switch mid-utterance cannot move the latency baseline.
func (c *CloudASRChain) Feed(samples []float32) {
	if !c.open || len(samples) == 0 {
		return
	}
	c.fed += len(samples)
	if len(c.buf)+len(samples) <= chainReplayCapSamples {
		c.buf = append(c.buf, samples...)
	} else {
		// Truncation is said out loud at the next step; silently dropping
		// audio from the replay would corrupt the next provider's transcript
		// with nobody able to attribute why.
		c.bufDropped = true
	}
	if e := c.activeEngine(); e != nil {
		e.Feed(samples)
	}
}

// Partial is the active head's hypothesis -- whatever the currently winning
// provider has said so far.
func (c *CloudASRChain) Partial() string {
	if e := c.activeEngine(); e != nil {
		return e.Partial()
	}
	return ""
}

// Finish closes the utterance and walks the chain if the active head failed.
// The walk replays the buffered utterance into each next element (Begin +
// Feed of the whole buffer), so every provider that gets a turn hears
// exactly what the caller fed, and each hop fires a visible failover event
// -- including the terminal hop to local sherpa.
func (c *CloudASRChain) Finish() string {
	if !c.open {
		return ""
	}
	c.open = false

	if c.head == headNone {
		// Every element was quota-skipped and no local fallback is wired:
		// the skips already spoke; the empty final is the caller's
		// "nothing sayable" signal, and the ledger has the reason.
		return ""
	}

	text, cause := c.finishAt(c.head)
	for cause != nil && cause.Class != observe.ClassCancelled {
		next := headNone
		if c.head != headLocal {
			next = c.selectHead(c.head + 1)
		}
		if next == headNone && c.fallback == nil {
			c.reportFailover(c.currentName(), "", string(cause.Class))
			break
		}
		to := c.nameOf(next)
		if next == headNone {
			to = localSherpaName
			next = headLocal
		}
		if c.bufDropped {
			to += fmt.Sprintf(" (replay truncated at %ds)", replayCapSeconds)
		}
		c.reportFailover(c.currentName(), to, string(cause.Class))
		c.setHead(next)
		c.beginAt(next)
		for _, chunk := range splitReplay(c.buf) {
			if e := c.activeEngine(); e != nil {
				e.Feed(chunk)
			}
		}
		text, cause = c.finishAt(next)
		if next == headLocal {
			// The terminal fallback's own reading is final by contract: no
			// replay loop past local, and no inspecting sherpa for a
			// "failure" the seam cannot express.
			break
		}
	}
	c.lastErr = cause
	return text
}

// InFlight reports whether an utterance is open.
func (c *CloudASRChain) InFlight() bool { return c.open }

// Fed is the chain's audio clock -- what the CALLER fed, not what any
// provider managed to receive (transcript latencies stay comparable across
// engines and across mid-utterance switches, which is the point D42#9 makes
// for the single engine; the chain inherits it verbatim).
func (c *CloudASRChain) Fed() int { return c.fed }

// Close releases every cloud engine the chain built. The fallback is NOT
// closed -- its owner loads and releases it (ownership rule documented on
// NewCloudASRChain). Idempotent: engine Close is idempotent (C11).
func (c *CloudASRChain) Close() error {
	var errs []string
	for _, e := range c.engines {
		if err := e.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	c.head = headNone
	c.open = false
	if len(errs) > 0 {
		return errors.New("speech: cloud ASR chain close: " + strings.Join(errs, "; "))
	}
	return nil
}

// LastFailure is the cascade-private window onto the most recent Finish,
// mirroring CloudASR.LastFailure: it exists so the chain's caller (and
// tests) can attribute a stepped walk after the fact. Seam callers never
// see it.
func (c *CloudASRChain) LastFailure() *observe.Error { return c.lastErr }

// Active names the head currently serving (a "provider/model" or
// localSherpaName) -- the same string OnProviderState carries, exposed for
// pull-style diagnostics (the panel renders from events; tests assert on
// this without a callback).
func (c *CloudASRChain) Active() string { return c.currentName() }

// -- ASR chain internals ----------------------------------------------------

// selectHead walks elements from fromIdx and returns the first the quota
// gate allows, firing a visible skip event for each denial. headNone comes
// back when no element is allowed.
func (c *CloudASRChain) selectHead(fromIdx int) int {
	for i := fromIdx; i < len(c.engines); i++ {
		ep := c.engines[i].Endpoint()
		ok, reason := c.quotaAllows(ep)
		if ok {
			return i
		}
		to := ""
		if i+1 < len(c.engines) {
			to = c.names[i+1]
		} else if c.fallback != nil {
			to = localSherpaName
		}
		c.reportSkip(i, to, reason)
	}
	return headNone
}

func (c *CloudASRChain) quotaAllows(ep CloudEndpoint) (bool, string) {
	if c.cb.Quota == nil {
		return true, ""
	}
	return c.cb.Quota.Allow(ep.Provider, ep.Model)
}

// beginAt opens the utterance on the head position (no-op for headNone).
func (c *CloudASRChain) beginAt(idx int) {
	if e := c.engineAt(idx); e != nil {
		e.Begin()
	}
}

// finishAt closes the utterance at the head position and reports (text,
// classified cause or nil). For cloud heads the cause comes from
// LastFailure -- the seam itself stays two-valued (text / empty) exactly
// like the local engine, and only the cascade reads the third value.
func (c *CloudASRChain) finishAt(idx int) (string, *observe.Error) {
	e := c.engineAt(idx)
	if e == nil {
		return "", nil
	}
	text := e.Finish()
	if asr, ok := e.(*CloudASR); ok {
		return text, asr.LastFailure()
	}
	// The local terminal fallback: whatever came back is the product's
	// final reading. A sherpa failure is invisible at the seam (an empty
	// final) and stepping "past" local is forbidden by contract.
	return text, nil
}

// engineAt resolves a head position to an engine (nil for headNone).
func (c *CloudASRChain) engineAt(idx int) AsrEngine {
	switch idx {
	case headLocal:
		return c.fallback
	case headNone:
		return nil
	default:
		if idx >= 0 && idx < len(c.engines) {
			return c.engines[idx]
		}
		return nil
	}
}

func (c *CloudASRChain) activeEngine() AsrEngine { return c.engineAt(c.head) }

// setHead records the head and fires the provider-state transition when the
// answer to "is audio leaving the machine right now" actually changed (the
// privacy indicator is a state, not a log line: repeated identical states
// must not flicker the UI).
func (c *CloudASRChain) setHead(idx int) {
	oldName, oldCloud := c.currentName(), c.isCloudHead()
	c.head = idx
	newName, newCloud := c.currentName(), c.isCloudHead()
	if oldName == newName && oldCloud == newCloud {
		return
	}
	fireProviderState(c.cb.OnProviderState, ChainASR, newName, newCloud)
}

func (c *CloudASRChain) isCloudHead() bool { return c.head != headNone && c.head != headLocal }

func (c *CloudASRChain) currentName() string { return c.nameOf(c.head) }

func (c *CloudASRChain) nameOf(idx int) string {
	switch idx {
	case headLocal:
		return localSherpaName
	case headNone:
		return ""
	default:
		if idx >= 0 && idx < len(c.names) {
			return c.names[idx]
		}
		return ""
	}
}

func (c *CloudASRChain) reportSkip(i int, to, reason string) {
	if c.cb.OnFailover == nil {
		return
	}
	c.cb.OnFailover(VoiceFailoverEvent{
		Chain: ChainASR, Index: i, From: c.nameOf(i), To: to,
		Cause: quotaCauseSkip, Detail: reason, Attempt: i + 1, Total: len(c.names),
	})
}

func (c *CloudASRChain) reportFailover(from, to, cause string) {
	if c.cb.OnFailover == nil {
		return
	}
	c.cb.OnFailover(VoiceFailoverEvent{
		Chain: ChainASR, Index: c.indexOfName(from), From: from, To: to,
		Cause: cause, Attempt: c.indexOfName(from) + 1, Total: len(c.names),
	})
}

func (c *CloudASRChain) indexOfName(name string) int {
	for i, n := range c.names {
		if n == name {
			return i
		}
	}
	return 0
}

// -- TTS chain: the seam, walked ---------------------------------------------

// Speak walks voice.cloud_tts_chain for one utterance. Rules, each with its
// reason (see the file header): speakable-guard first (one check, never a
// per-provider re-attempt), quota BEFORE dispatch, no failover once audio
// was delivered, cancellation never a step cause, local cascade terminal
// and never silent.
func (c *CloudTTSChain) Speak(text string, emit func(samples []float32) error) error {
	if emit == nil {
		return errors.New("speech: Speak needs an emit callback")
	}
	if !speakable(text) {
		return ErrNothingToSpeak
	}

	var lastErr *observe.Error
	// skips counts the elements quota refused before they could try. The
	// skip event has already been spoken for each; the terminal hop into
	// local is announced here exactly once, and never twice.
	skips := 0
	for i := range c.engines {
		ep := c.engines[i].ep // CloudTTS stores its endpoint; no getter needed inside the package
		ok, reason := c.quotaAllows(ep)
		if !ok {
			skips++
			to := localSherpaName
			if i+1 < len(c.engines) {
				to = c.names[i+1]
			}
			c.reportSkip(i, to, reason)
			continue
		}
		c.setActive(i, true)
		err := c.engines[i].Speak(text, emit)
		if err == nil {
			return nil
		}
		oe := classifySpeakError(err)
		lastErr = oe
		if oe.Class == observe.ClassCancelled {
			return err
		}
		delivered := emittedSamples(c.engines[i]) > 0
		if delivered {
			// Half the sentence already played; a second provider would
			// append fresh audio to a broken utterance. The announce path
			// owns what to do with the remainder (ticket 26's text-only
			// fallback still delivers the result).
			return err
		}
		to := ""
		if i+1 < len(c.engines) {
			to = c.names[i+1]
		} else if c.fallback != nil {
			to = localSherpaName
		}
		c.reportFailover(c.names[i], to, string(oe.Class))
	}
	if c.fallback == nil {
		if lastErr != nil {
			return lastErr
		}
		return observe.New(observe.ClassProvider,
			"speech: cloud TTS chain exhausted without a classified failure")
	}
	// The hop into the terminal local cascade is ANNOUNCED exactly once.
	// The last element's failure already spoke it (to == localSherpaName)
	// inside the loop above; re-reporting here would log the same hop twice
	// and make the failover counter claims un-auditable -- the same reason
	// the ASR chain's walk emits one event per hop. A walk where every
	// element was quota-skipped never reached that line, and its skips all
	// pointed at each other, so the landing on local is said here.
	if len(c.engines) > 0 && skips == len(c.engines) {
		c.reportFailover(c.nameOfAt(len(c.names)-1), localSherpaName, quotaCauseSkip)
	}
	c.setActive(headLocal, false)
	return c.fallback.Speak(text, emit)
}

// SampleRate follows the head that last spoke (or will first speak): the
// render side resamples from the rate of the audio IT is about to receive,
// and after a mid-sentence switch there is no such thing as one global
// cloud rate -- there is the current engine's.
//
// A hop into local is the case the render side must survive: it is holding
// the previous cloud head's rate, the samples arriving now are the sherpa
// vocoder's, and answering with the stale rate would resample a 22.05kHz
// clip as if it were 24kHz (a 9% pitch and speed error over the whole
// terminal fallback). So while the local cascade owns audio out, the
// answer is the FALLBACK's rate -- exactly the invariant this getter has
// always promised (SampleRate is "the rate of the samples Speak emits").
func (c *CloudTTSChain) SampleRate() int {
	if c.active == headLocal && c.fallback != nil {
		return c.fallback.SampleRate()
	}
	if c.active >= 0 && c.active < len(c.engines) {
		return c.engines[c.active].SampleRate()
	}
	return 0
}

// Close releases the cloud engines; the local fallback belongs to its
// loader (same ownership rule as the ASR chain). Idempotent.
func (c *CloudTTSChain) Close() error {
	var errs []string
	for _, e := range c.engines {
		if err := e.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	c.active = headNone
	if len(errs) > 0 {
		return errors.New("speech: cloud TTS chain close: " + strings.Join(errs, "; "))
	}
	return nil
}

// Active names the engine currently responsible for audio out.
func (c *CloudTTSChain) Active() string {
	if c.active == headLocal {
		return localSherpaName
	}
	if c.active >= 0 && c.active < len(c.names) {
		return c.names[c.active]
	}
	return ""
}

// -- TTS chain internals -----------------------------------------------------

func (c *CloudTTSChain) quotaAllows(ep CloudEndpoint) (bool, string) {
	if c.cb.Quota == nil {
		return true, ""
	}
	return c.cb.Quota.Allow(ep.Provider, ep.Model)
}

func (c *CloudTTSChain) setActive(idx int, cloud bool) {
	if c.active == idx && c.cloud == cloud {
		return
	}
	c.active = idx
	c.cloud = cloud
	fireProviderState(c.cb.OnProviderState, ChainTTS, c.Active(), cloud)
}

func (c *CloudTTSChain) reportSkip(i int, to, reason string) {
	if c.cb.OnFailover == nil {
		return
	}
	c.cb.OnFailover(VoiceFailoverEvent{
		Chain: ChainTTS, Index: i, From: c.nameOfAt(i), To: to,
		Cause: quotaCauseSkip, Detail: reason, Attempt: i + 1, Total: len(c.names),
	})
}

func (c *CloudTTSChain) reportFailover(from, to, cause string) {
	if c.cb.OnFailover == nil {
		return
	}
	idx := c.indexOfName(from)
	c.cb.OnFailover(VoiceFailoverEvent{
		Chain: ChainTTS, Index: idx, From: from, To: to,
		Cause: cause, Attempt: idx + 1, Total: len(c.names),
	})
}

func (c *CloudTTSChain) nameOfAt(i int) string {
	if i >= 0 && i < len(c.names) {
		return c.names[i]
	}
	return localSherpaName
}

func (c *CloudTTSChain) indexOfName(name string) int {
	for i, n := range c.names {
		if n == name {
			return i
		}
	}
	return 0
}

// -- shared helpers -----------------------------------------------------------

// classifySpeakError normalizes whatever a TTS engine returned into a
// D37 *observe.Error (engines here already return classified errors; this
// covers third-party TtsEngines that return plain errors -- the local
// matcha's plain errors are read as internal, which steps nothing extra:
// the chain never steps PAST local anyway).
func classifySpeakError(err error) *observe.Error {
	var oe *observe.Error
	if errors.As(err, &oe) {
		return oe
	}
	return observe.Wrap(observe.ClassOfOrInternal(err), err, "speech: cloud TTS chain step")
}

// splitReplay feeds the replay buffer in seam-window-sized chunks. Chunking
// matters: cloud engines convert float32->s16le per call and the local
// engine runs its per-window decode loop; one 2-minute Feed call would be
// legal per the interface but would spike the owner thread and stall
// endpointing math that assumes window-scale calls.
func splitReplay(buf []float32) [][]float32 {
	if len(buf) == 0 {
		return nil
	}
	const window = 512
	var out [][]float32
	for off := 0; off < len(buf); off += window {
		end := off + window
		if end > len(buf) {
			end = len(buf)
		}
		out = append(out, buf[off:end])
	}
	return out
}

// emittedSamples reports whether an engine already handed audio to the
// consumer (the "delivered" rule's observable, read under the engine lock).
func emittedSamples(t *CloudTTS) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.samples
}

func fireProviderState(fn func(VoiceProviderState), kind CloudChainKind, active string, cloud bool) {
	if fn != nil {
		fn(VoiceProviderState{Chain: kind, Active: active, Cloud: cloud})
	}
}

func initialHead(numEngines int, hasFallback bool) int {
	if numEngines > 0 {
		return 0
	}
	if hasFallback {
		return headLocal
	}
	return headNone
}

func initialTTSHead(numEngines int, hasFallback bool) (int, bool) {
	if numEngines > 0 {
		return 0, true
	}
	if hasFallback {
		return headLocal, false
	}
	return headNone, false
}
