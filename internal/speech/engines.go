package speech

import "fmt"

// This file is the C9 layer (PLAN.md C9, SPEC-04 section 4): the engine seams
// the voice chain is written against, and the one property C9 says makes them
// worth having -- "local sherpa and cloud provider interchangeable".
//
// Until this file existed, that property was unfalsifiable: 15-a1's census
// found all four names living only in doc.go's prose, with zero Go definitions,
// so there was nothing to conform to, nothing to swap and nothing to check.
//
// The roster, and where each implementation lands:
//
//	AsrEngine       local: *StreamingAsr (15-c, landed). Cascade cloud: ticket 61.
//	TtsEngine       reserved slot -- ticket 26 local, ticket 61 cloud.
//	WakeWordEngine  reserved slot -- ticket 41, local only by contract.
//	RealtimeEngine  cloud only by contract -- ticket 60, specified by C32.
//
// "Reserved slot" is the shape C8 already uses for its AEC seam: declared and
// unimplemented. Ticket 15 puts TTS, KWS and cloud providers out of scope, so
// the alternative to a reserved slot is not an implementation, it is silence --
// and a seam added after the pipeline that would have used it is a rewrite, not
// a seam.
//
// What this file deliberately does not do: it adds no configuration key (the
// engine choice is read from the existing voice section, and every number
// behind these seams stays a package constant, as 15-b and 15-c decided), and
// it pulls in no cgo (this is plain Go and builds without the sherpa DLLs; only
// the implementations behind AsrEngine cross the C boundary).

// AsrEngine is the C9 ASR seam: one streaming recognition session, driven by
// its single owner goroutine (D38b).
//
// The method set is 15-c's landed session, extracted rather than designed. The
// direction matters: the local implementation is the one carrying measurements
// (first partial 640ms, flush 1ms), so bending it to fit a hypothetical cloud
// provider would trade a real number for a guess. An engine that cannot express
// this shape is the thing that is wrong, not the shape.
//
// Three clauses carry reasons, not taste:
//
//   - No context, no error return. A session is driven by its owner, so
//     cancellation is the caller's context and shutdown is the caller's wait;
//     and a decode that fails surfaces as a short or empty final transcript,
//     which is the single signal AC#5's retry-then-"did not catch that" policy
//     keys on. At this seam a cloud provider that errors is therefore
//     indistinguishable from a local model that misheard, and the caller never
//     branches on which one it is holding.
//   - Fed is the audio clock, in samples, never a wall clock (D42#9). That is
//     what makes Transcript.Latency the same number on every host.
//   - Close releases (C11). Sessions are released, not unloaded: freeing
//     onnxruntime is not reliable, and D5 keeps one runtime for the whole chain.
type AsrEngine interface {
	// Begin opens an utterance with a fresh decoder state.
	Begin()
	// Feed appends utterance audio: 16kHz mono float32, any length.
	Feed(samples []float32)
	// Partial is the hypothesis for the open utterance so far. It is empty
	// until the engine has enough audio to say anything, and inside one
	// utterance it only ever grows.
	Partial() string
	// Finish closes the utterance and returns its transcript. Empty means
	// "nothing sayable", not a failure code.
	Finish() string
	// InFlight reports whether an utterance is open.
	InFlight() bool
	// Fed is how many samples the open -- or, right after Finish, the most
	// recently closed -- utterance streamed.
	Fed() int
	// Close releases the engine. Idempotent.
	Close() error
}

// The local sherpa implementation of AsrEngine, asserted at compile time: if
// *StreamingAsr stops matching this seam the package stops building, which is
// the only kind of conformance check that cannot rot.
var _ AsrEngine = (*StreamingAsr)(nil)

// TtsEngine is the C9 TTS seam: text in, speech out, in playback order.
//
// The shape is streaming because the budget is. SPEC-04 section 8 gives TTS
// 800ms from the first token to the first audio frame, while the sherpa offline
// path returns a whole clip in one call -- a seam of the form
// "synthesize(text) []float32" would accept that wrapper happily and blow the
// budget with nothing to point at. Here the chunk is the unit, so a whole-clip
// implementation cannot express itself at all; the compiler rejects it before
// anyone has to measure anything (the tooth: 15a-r1-impl.md section 4).
//
// emit is called once per chunk, in order, and its error aborts synthesis and
// is returned verbatim -- the cooperative-abort contract C5's LlmProvider
// already has. That is how Path C's barge-in stops a sentence that is
// mid-flight, with no goroutine, no timer and no wall clock.
type TtsEngine interface {
	// Speak synthesizes text and emits its audio, starting with the first
	// chunk as soon as one exists.
	Speak(text string, emit func(samples []float32) error) error
	// SampleRate is the rate of the samples Speak emits, for the render side to
	// resample from.
	SampleRate() int
	// Close releases the engine (C11). Idempotent.
	Close() error
}

// KeywordHit is one keyword detection (SPEC-04 section 5).
type KeywordHit struct {
	// Word is the keyword, spelled as the keyword table spells it.
	Word string
	// Veto marks a veto word ("cancel", "stop") rather than a wake word. Both
	// live in one table, and the veto channel is the only voice way out of
	// Confirming (B1) -- an engine that cannot say which kind it hit cannot
	// serve that state.
	Veto bool
	// At is the hit's position on the audio clock: how many samples had been
	// fed when it fired. Never a wall clock (D42#9).
	At int
}

// WakeWordEngine is the C9 wake-word seam: audio in, keyword hits out, for as
// long as the engine is alive. Reserved slot: the local kws-zipformer
// implementation is ticket 41.
//
// Accept has Detector's shape on purpose -- samples in, append to out -- because
// KWS is fed by the same seam: 512 samples of 16kHz mono out of the same C8
// frames, on the same owner goroutine (D38b). No method arms or disarms
// anything: arming is the state machine's business (see doc.go's
// non-responsibilities) and the engine's lifetime is the Armed period.
type WakeWordEngine interface {
	// Accept feeds samples (any length, any chunking) and appends the keywords
	// they triggered.
	Accept(samples []float32, out []KeywordHit) []KeywordHit
	// Close releases the engine (C11). Idempotent.
	Close() error
}

// RealtimeEventKind says what a RealtimeEvent is.
type RealtimeEventKind int

const (
	// RealtimeAudio is model speech to play, at RealtimeEngine.SampleRate.
	RealtimeAudio RealtimeEventKind = iota + 1

	// RealtimeTranscript is the far side's transcript of user speech. It
	// reaches task_log like any other transcript, under the same rule: length
	// and confidence may be logged, the text may not.
	RealtimeTranscript

	// RealtimeEndpoint is the far side's own end of utterance. Unlike the local
	// path, this engine endpoints on the server (C32), so the Go state machine
	// hears about it instead of deciding it.
	RealtimeEndpoint

	// RealtimeBargeIn is the user speaking over playback (D47 Path C).
	RealtimeBargeIn

	// RealtimeHandoff is C32 (1)'s structured task-intent signal: the host is
	// asked to run its own text loop, and is not handed a tool call.
	RealtimeHandoff

	// RealtimeUsage carries the billable audio-token count for C23.
	RealtimeUsage
)

// String is the log/diagnostic spelling (D38d: reasons reach logs).
func (k RealtimeEventKind) String() string {
	switch k {
	case RealtimeAudio:
		return "audio"
	case RealtimeTranscript:
		return "transcript"
	case RealtimeEndpoint:
		return "endpoint"
	case RealtimeBargeIn:
		return "barge-in"
	case RealtimeHandoff:
		return "handoff"
	case RealtimeUsage:
		return "usage"
	default:
		return fmt.Sprintf("realtime-event(%d)", int(k))
	}
}

// RealtimeEvent is one thing a realtime session tells the host.
type RealtimeEvent struct {
	Kind RealtimeEventKind
	// Audio is model speech; RealtimeAudio only.
	Audio []float32
	// Text is a transcript; RealtimeTranscript and RealtimeHandoff only.
	Text string
	// AudioTokens is the billable count; RealtimeUsage only (C23 bills per
	// audio token).
	AudioTokens int
}

// RealtimeEngine is the C9 realtime seam: the cloud speech-to-speech brain of
// D47's Path C, specified in full by C32. Cloud only, by contract. Reserved
// slot: ticket 60.
//
// Two of C32's clauses are structural here rather than documented:
//
//   - Zero tool permissions. The interface has no method that could run
//     anything -- no tool list, no approval surface -- so code holding a
//     RealtimeEngine cannot make it act: it can only talk. A task intent
//     arrives as RealtimeHandoff, a host-side signal, and the host then runs its
//     own fully gated text loop (C5). That is the whole of C32 (1)'s handoff
//     protocol as far as this seam is concerned, and it is why the handoff is an
//     event and not a call.
//   - Suspend, not destroy. C32 (1) parks the session while the text loop works
//     and resumes it once the result is delivered; that is why a session is a
//     session and not one call per turn.
//
// Send has Accept's shape for the same reason WakeWordEngine does: the owner
// goroutine drives the session (D38b), audio goes in as seam samples, and
// whatever the far side produced since the last call comes back on out.
type RealtimeEngine interface {
	// Send pushes caller audio (16kHz mono) and appends the events it produced.
	Send(samples []float32, out []RealtimeEvent) []RealtimeEvent
	// SampleRate is the rate of RealtimeAudio samples.
	SampleRate() int
	// Suspend parks the session for a task handoff (C32 (1)).
	Suspend() error
	// Resume returns to the conversation once the task's result is delivered.
	Resume() error
	// Close ends the session (C11). Idempotent.
	Close() error
}
