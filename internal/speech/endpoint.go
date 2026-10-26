package speech

import (
	"fmt"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
)

// Endpointing is the VAD half of the voice chain (SPEC-04 sec 4): it decides
// where an utterance starts and stops. The neural part (silero) only answers
// "is this 32ms window speech?"; every duration-based policy lives HERE, in
// Go, on an audio clock - never on a wall clock (D42#9). That split is what
// makes the three AC#3 numbers testable without a corpus and without cgo.
//
// The policy numbers are package constants, not configuration keys: the SPEC
// writes them as defaults and D36's config model has no silence/vad key yet
// (ticket 15b-r1 refused to invent one; adding [voice].asr.silence_ms is a
// separate ticket because it moves the frozen config schema).
const (
	// EndpointSilence is the trailing-silence run that closes an utterance.
	// SPEC-04 sec 4 states the default as ~700ms and no config key exists.
	EndpointSilence = 700 * time.Millisecond

	// MinSpeech is the shortest speech run that opens an utterance; a shorter
	// burst is ignored outright and never reaches the ASR (AC#3 second clause).
	MinSpeech = 300 * time.Millisecond

	// MaxUtterance force-finalizes a run-on utterance so the ASR stream cannot
	// grow without bound (SPEC-04 sec 4: 60s guard, auto-finalize + notify).
	MaxUtterance = 60 * time.Second
)

// WindowSamples is the silero input window: one C8 seam frame (512 @ 16k =
// 32ms, internal/audio). The endpointing resolution is this window, so the
// thresholds above are honoured to the sample but OBSERVED to the window.
const WindowSamples = audio.FrameSamples

// samplesOf converts a duration to a sample count at the seam rate.
func samplesOf(d time.Duration) int {
	return int(int64(d) * int64(audio.TargetRate) / int64(time.Second))
}

// The policy in samples, computed once. Comparisons are done on these so the
// state machine never converts per frame.
var (
	endpointSilenceSamples = samplesOf(EndpointSilence)
	minSpeechSamples       = samplesOf(MinSpeech)
	maxUtteranceSamples    = samplesOf(MaxUtterance)
)

// EndpointKind says why an Endpoint was emitted.
type EndpointKind int

const (
	// EndpointStart means an utterance just opened: at least MinSpeech of
	// speech arrived. The ball's Listening state keys off this (D43#11); the
	// audio is NOT carried on a start event.
	EndpointStart EndpointKind = iota + 1

	// EndpointEnd means the utterance closed because EndpointSilence of
	// trailing silence elapsed.
	EndpointEnd

	// EndpointMaxDuration means the utterance closed because MaxUtterance
	// elapsed first (the 60s guard).
	EndpointMaxDuration

	// EndpointEndOfStream means the utterance was still open when the stream
	// ended, so the caller's Finalize closed it.
	EndpointEndOfStream
)

// String is the log/diagnostic spelling (D38d: counts and reasons reach logs).
func (k EndpointKind) String() string {
	switch k {
	case EndpointStart:
		return "start"
	case EndpointEnd:
		return "end"
	case EndpointMaxDuration:
		return "max-duration"
	case EndpointEndOfStream:
		return "end-of-stream"
	default:
		return fmt.Sprintf("endpoint-kind(%d)", int(k))
	}
}

// Endpoint is one state-machine transition. Start events carry no audio and
// zero durations; the three closing kinds carry the utterance.
type Endpoint struct {
	Kind EndpointKind
	// Samples is the utterance audio, 16k mono float32, including the leading
	// MinSpeech pre-roll and the trailing silence that closed it. Nil on
	// EndpointStart.
	Samples []float32
	// Speech is how many of the utterance's samples the VAD called speech.
	Speech int
	// Total is the utterance length in samples (Speech plus trailing silence).
	Total int
}

// SpeechDuration is the VAD-classified speech inside the utterance.
func (e Endpoint) SpeechDuration() time.Duration { return samplesDuration(e.Speech) }

// Duration is the whole utterance, trailing silence included.
func (e Endpoint) Duration() time.Duration { return samplesDuration(e.Total) }

// samplesDuration converts a sample count back to an audio duration.
func samplesDuration(n int) time.Duration {
	return time.Duration(int64(n) * int64(time.Second) / int64(audio.TargetRate))
}

// Endpointer is the endpoint state machine. It is fed (window, isSpeech) in
// audio order and owns every duration decision. It is not safe for concurrent
// use; one owner goroutine drives it (D38b).
type Endpointer struct {
	inSpeech   bool
	speechRun  int       // consecutive speech samples while idle
	silenceRun int       // consecutive silence samples while in speech
	total      int       // utterance length in samples
	speech     int       // speech samples inside the utterance
	pre        []float32 // idle-time speech pre-roll (bounded by MinSpeech)
	utter      []float32 // the open utterance's audio
}

// NewEndpointer returns a machine in the idle state.
func NewEndpointer() *Endpointer { return &Endpointer{} }

// InSpeech reports whether an utterance is currently open.
func (e *Endpointer) InSpeech() bool { return e.inSpeech }

// Accept feeds one window of audio and its VAD verdict, appending any
// transitions to out. window may be any length: production feeds
// WindowSamples, tests feed exact boundary lengths. out is reused by the
// caller, so the returned slice aliases it (append idiom, no per-frame
// allocation on the hot path).
func (e *Endpointer) Accept(window []float32, isSpeech bool, out []Endpoint) []Endpoint {
	n := len(window)
	if n == 0 {
		return out
	}
	if !e.inSpeech {
		if !isSpeech {
			e.speechRun = 0
			e.pre = e.pre[:0]
			return out
		}
		e.speechRun += n
		e.pre = append(e.pre, window...)
		if e.speechRun < minSpeechSamples {
			return out
		}
		// MinSpeech reached: open, keeping the speech that triggered it so the
		// onset is not clipped out of the transcript's audio.
		e.inSpeech = true
		e.utter = append(e.utter[:0], e.pre...)
		e.total = e.speechRun
		e.speech = e.speechRun
		e.silenceRun = 0
		e.pre = e.pre[:0]
		e.speechRun = 0
		out = append(out, Endpoint{Kind: EndpointStart, Speech: e.speech, Total: e.total})
		return out
	}

	e.utter = append(e.utter, window...)
	e.total += n
	if isSpeech {
		e.silenceRun = 0
		e.speech += n
	} else {
		e.silenceRun += n
	}

	if e.silenceRun >= endpointSilenceSamples {
		// Trim the samples that arrived after the threshold was crossed, so an
		// utterance is exactly MinSpeech..EndpointSilence long no matter how the
		// caller chunked the audio (see TestStateMachineIgnoresChunking).
		out = append(out, e.close(EndpointEnd, e.silenceRun-endpointSilenceSamples))
	} else if e.total >= maxUtteranceSamples {
		out = append(out, e.close(EndpointMaxDuration, e.total-maxUtteranceSamples))
	}
	return out
}

// Finalize closes an open utterance as EndpointEndOfStream. It is the
// end-of-stream counterpart of Accept (capture stopped, file ended).
func (e *Endpointer) Finalize(out []Endpoint) []Endpoint {
	if !e.inSpeech {
		return out
	}
	return append(out, e.close(EndpointEndOfStream, 0))
}

// close emits the closing transition and returns the machine to idle. drop is
// how many tail samples arrived after the transition and must not be attributed
// to the utterance; the 60s guard additionally truncates the audio, so a
// run-on utterance never reports more than MaxUtterance.
func (e *Endpointer) close(kind EndpointKind, drop int) Endpoint {
	n := e.total - drop
	if n < 0 {
		n = 0
	}
	if n > maxUtteranceSamples {
		n = maxUtteranceSamples
	}
	ep := Endpoint{Kind: kind, Speech: e.speech, Total: n}
	if ep.Speech > n {
		ep.Speech = n
	}
	ep.Samples = append([]float32(nil), e.utter[:n]...)
	e.inSpeech = false
	e.speechRun = 0
	e.silenceRun = 0
	e.speech = 0
	e.total = 0
	e.utter = e.utter[:0]
	e.pre = e.pre[:0]
	return ep
}
