package speech

import (
	"context"
	"fmt"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
)

// TranscriptKind says what a Transcript is.
type TranscriptKind int

const (
	// TranscriptPartial is a hypothesis for the utterance still in flight. It
	// exists for latency: the ball can start thinking before the speaker stops
	// (AC#1 / D43#11) because text arrives while the audio does.
	TranscriptPartial TranscriptKind = iota + 1

	// TranscriptFinal closes one utterance and carries the endpoint that
	// closed it, including the utterance audio the ASR was fed.
	TranscriptFinal
)

// String is the log/diagnostic spelling.
func (k TranscriptKind) String() string {
	switch k {
	case TranscriptPartial:
		return "partial"
	case TranscriptFinal:
		return "final"
	default:
		return fmt.Sprintf("transcript-kind(%d)", int(k))
	}
}

// Transcript is one thing the voice chain has to say about an utterance.
type Transcript struct {
	Kind TranscriptKind
	// Text is the hypothesis. On a final it is the utterance's transcript, and
	// it is empty when the ASR heard nothing sayable -- that empty final is the
	// signal AC#5's retry-then-"did not catch that" policy keys on (that policy
	// lives with the caller, not here).
	Text string
	// Latency is the distance from the utterance's first sample to the window
	// that produced this event, measured by counting fed samples rather than by
	// reading a clock, so the same fixture always reports the same number
	// (D42#9). On a final it is the utterance's length.
	Latency time.Duration
	// Endpoint is the transition that closed the utterance; the zero value on
	// partials. Its Samples are the audio this ASR streamed.
	Endpoint Endpoint
}

// Transcriber is the seam between the two halves of ticket 15: it drives the
// VAD endpointer (15-b) and the streaming ASR (this slice) from one audio
// stream, and turns "an utterance opened / an utterance closed" into "partials
// while it is open / one final when it closes".
//
// It owns neither part: the caller constructs the Detector and the AsrEngine
// (C9; *StreamingAsr is the local sherpa one) and closes them (C11). That keeps
// the expensive recognizer reusable across utterances and keeps the utterance
// policy in one place. The engine is the seam, not a concrete session: nothing
// in this driver knows which provider it is driving.
//
// It starts no goroutines and is not safe for concurrent use: one owner
// goroutine drives it (D38b), exactly like Detector.
type Transcriber struct {
	det *Detector
	asr AsrEngine

	rest []float32 // input not yet filling one seam window
	ring []float32 // tail of the audio, so an opening utterance can be replayed
	text string    // last emitted partial, to suppress repeats
}

// NewTranscriber wires an endpointer to any ASR engine (C9).
func NewTranscriber(det *Detector, asr AsrEngine) *Transcriber {
	return &Transcriber{det: det, asr: asr}
}

// Accept feeds samples (any length, any chunking) and appends the transcripts
// they produce. It buffers to whole windows and drives the detector one window
// at a time, because the window is the unit the ASR replay has to line up with.
func (t *Transcriber) Accept(samples []float32, out []Transcript) []Transcript {
	for len(samples) > 0 {
		need := WindowSamples - len(t.rest)
		if need > len(samples) {
			need = len(samples)
		}
		t.rest = append(t.rest, samples[:need]...)
		samples = samples[need:]
		if len(t.rest) < WindowSamples {
			break
		}
		out = t.window(t.rest, out)
		t.rest = t.rest[:0]
	}
	return out
}

// window runs one seam window through the detector and keeps the ASR in step.
//
// The three cases are the whole seam:
//   - the window opened an utterance: start a decoder state and replay the
//     pre-roll, so the onset MinSpeech waited for is not clipped out of the
//     transcript;
//   - the window closed one: feed it, finish, emit the final;
//   - neither: feed it and publish a partial if the hypothesis grew.
func (t *Transcriber) window(win []float32, out []Transcript) []Transcript {
	t.ring = ringAppend(t.ring, win)
	var eps []Endpoint
	eps = t.det.Accept(win, eps)
	for _, ep := range eps {
		if ep.Kind == EndpointStart {
			t.asr.Begin()
			// An EndpointStart carries no audio (15-b): it carries the length
			// of the speech run that opened the utterance instead. The ring's
			// tail of exactly that length IS that run, so the replay is exact
			// rather than approximate.
			t.asr.Feed(t.ringTail(ep.Total))
			out = t.partial(out)
			continue
		}
		// A closing window is part of the utterance (Endpoint.Samples includes
		// it), so it is fed before the flush.
		t.asr.Feed(win)
		out = t.final(ep, out)
	}
	if len(eps) == 0 && t.asr.InFlight() {
		t.asr.Feed(win)
		out = t.partial(out)
	}
	return out
}

// ringTail is the last n samples of the ring, or everything it holds when the
// run is longer than the ring (which cannot happen for a MinSpeech opening, but
// is not worth a panic).
func (t *Transcriber) ringTail(n int) []float32 {
	if n >= len(t.ring) {
		return t.ring
	}
	if n <= 0 {
		return nil
	}
	return t.ring[len(t.ring)-n:]
}

// partial publishes the hypothesis if it grew. Empty and unchanged hypotheses
// are not events: a partial that says nothing is not news, and repeating the
// previous one would make the latency table meaningless.
func (t *Transcriber) partial(out []Transcript) []Transcript {
	text := t.asr.Partial()
	if text == "" || text == t.text {
		return out
	}
	t.text = text
	return append(out, Transcript{
		Kind:    TranscriptPartial,
		Text:    text,
		Latency: samplesDuration(t.asr.Fed()),
	})
}

// final closes the utterance.
func (t *Transcriber) final(ep Endpoint, out []Transcript) []Transcript {
	text := t.asr.Finish()
	t.text = ""
	return append(out, Transcript{
		Kind:     TranscriptFinal,
		Text:     text,
		Latency:  samplesDuration(t.asr.Fed()),
		Endpoint: ep,
	})
}

// Finalize closes the stream: it zero-pads a partial trailing window, then
// finalizes the endpointer, so an utterance still open when the audio stops
// still produces a transcript (EndpointEndOfStream).
func (t *Transcriber) Finalize(out []Transcript) []Transcript {
	if n := len(t.rest); n > 0 {
		pad := make([]float32, WindowSamples)
		copy(pad, t.rest)
		t.rest = t.rest[:0]
		out = t.window(pad, out)
	}
	var eps []Endpoint
	eps = t.det.Finalize(eps)
	for _, ep := range eps {
		out = t.final(ep, out)
	}
	return out
}

// Consume drains seam frames until the channel closes or ctx is cancelled,
// decoding each one to samples and feeding the transcriber. It runs in the
// caller's goroutine by design: the voice chain starts no detached goroutines
// (D38b), so cancellation is the caller's ctx and shutdown is the caller's
// wait. A closed channel is end of stream.
func (t *Transcriber) Consume(ctx context.Context, frames <-chan []byte, out []Transcript) ([]Transcript, error) {
	for {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		case frame, ok := <-frames:
			if !ok {
				return t.Finalize(out), nil
			}
			pcm, err := audio.DecodeFrame(frame)
			if err != nil {
				return out, fmt.Errorf("speech: %w", err)
			}
			out = t.Accept(pcmToFloat(pcm), out)
		}
	}
}

// ringAppend appends a window to the replay ring and keeps it bounded to one
// MinSpeech opening plus one window (the longest run a single Accept can turn
// into an EndpointStart).
func ringAppend(ring, win []float32) []float32 {
	ring = append(ring, win...)
	limit := minSpeechSamples + WindowSamples
	if len(ring) > limit {
		n := copy(ring, ring[len(ring)-limit:])
		ring = ring[:n]
	}
	return ring
}
