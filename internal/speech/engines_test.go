package speech

import (
	"errors"
	"reflect"
	"testing"
)

// This file is the receipt for the C9 layer. It does not measure the
// recognizer -- the 15-b/15-c tests do that, on the real models. What is pinned
// here is the seam itself: the method roster of each engine, that the three
// reserved slots can actually be written against, and that the landed 15-c
// driver runs over an engine that is not sherpa at all.

// c9Rosters is the C9 method roster of each seam, as data.
//
// It exists to make one specific edit loud. Removing a method is already fatal
// (the implementations stop matching and the compile-time assertion in
// engines.go fires first), but ADDING one is invisible to Go's structural
// typing: every existing implementation still conforms, callers still pass, and
// nothing says a word. Without this pin RealtimeEngine could grow a method that
// performs something, and C32's "zero tool permissions" would stop being a
// property of the type while every test stayed green.
var c9Rosters = []struct {
	name  string
	iface reflect.Type
	want  []string
}{
	{"AsrEngine", reflect.TypeOf((*AsrEngine)(nil)).Elem(), []string{
		"Begin", "Close", "Fed", "Feed", "Finish", "InFlight", "Partial",
	}},
	{"TtsEngine", reflect.TypeOf((*TtsEngine)(nil)).Elem(), []string{
		"Close", "SampleRate", "Speak",
	}},
	{"WakeWordEngine", reflect.TypeOf((*WakeWordEngine)(nil)).Elem(), []string{
		"Accept", "Close",
	}},
	{"RealtimeEngine", reflect.TypeOf((*RealtimeEngine)(nil)).Elem(), []string{
		"Close", "Resume", "SampleRate", "Send", "Suspend",
	}},
}

// TestC9EngineRostersArePinned is the roster nail: four seams, exactly these
// methods, no more. reflect returns exported methods sorted by name, so the
// pin is a set, not a source order.
func TestC9EngineRostersArePinned(t *testing.T) {
	for _, tc := range c9Rosters {
		got := make([]string, 0, tc.iface.NumMethod())
		for i := 0; i < tc.iface.NumMethod(); i++ {
			got = append(got, tc.iface.Method(i).Name)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%s roster = %v, want exactly %v (C9): a removed method breaks the "+
				"implementations, and an added one has to be argued here first -- "+
				"RealtimeEngine in particular may never grow a method that could act "+
				"(C32: zero tool permissions)", tc.name, got, tc.want)
		}
		t.Logf("%s: %d methods, %v", tc.name, len(got), got)
	}
}

// The three reserved slots, implemented by the smallest types that satisfy
// them. They are not stand-ins for the real engines -- those are tickets 26, 41
// and 60 -- and nothing here is cited as evidence about a model. They are the
// proof that each shape can be written against at all: a seam with no
// implementation AND no conceivable one is not a seam.
type stubTts struct{ rate int }

func (s *stubTts) Speak(text string, emit func(samples []float32) error) error {
	for range text {
		if err := emit([]float32{1}); err != nil {
			return err
		}
	}
	return nil
}

func (s *stubTts) SampleRate() int { return s.rate }
func (s *stubTts) Close() error    { return nil }

type stubKws struct {
	fed int
	hit bool
}

func (s *stubKws) Accept(samples []float32, out []KeywordHit) []KeywordHit {
	s.fed += len(samples)
	if !s.hit && s.fed >= minSpeechSamples {
		s.hit = true
		out = append(out, KeywordHit{Word: "wisp", At: s.fed})
	}
	return out
}

func (s *stubKws) Close() error { return nil }

type stubRealtime struct{ suspended bool }

func (s *stubRealtime) Send(samples []float32, out []RealtimeEvent) []RealtimeEvent {
	if s.suspended {
		return out
	}
	return append(out,
		RealtimeEvent{Kind: RealtimeEndpoint},
		RealtimeEvent{Kind: RealtimeHandoff, Text: "set a timer"},
	)
}

func (s *stubRealtime) SampleRate() int { return 24000 }
func (s *stubRealtime) Suspend() error  { s.suspended = true; return nil }
func (s *stubRealtime) Resume() error   { s.suspended = false; return nil }
func (s *stubRealtime) Close() error    { return nil }

var (
	_ TtsEngine      = (*stubTts)(nil)
	_ WakeWordEngine = (*stubKws)(nil)
	_ RealtimeEngine = (*stubRealtime)(nil)
)

// TestC9ReservedSlotsAreImplementable drives each stub once, so the seam is
// known to be usable and not only declarable. Three properties are load-bearing
// and each has a clause in the interface's doc: TTS chunks arrive in order and
// an emit error aborts synthesis verbatim (that is Path C's barge-in); a
// keyword hit carries the veto bit and the audio clock (B1's only voice
// channel); a realtime handoff is an event, and suspending stops the stream
// (C32 (1)).
func TestC9ReservedSlotsAreImplementable(t *testing.T) {
	tts := &stubTts{rate: 22050}
	stop := errors.New("barge-in")
	chunks := 0
	err := tts.Speak("abc", func(samples []float32) error {
		chunks++
		if chunks == 2 {
			return stop
		}
		return nil
	})
	if !errors.Is(err, stop) {
		t.Fatalf("Speak returned %v, want the emit error verbatim: an aborted sentence has to "+
			"be distinguishable from a failed one", err)
	}
	if chunks != 2 {
		t.Fatalf("Speak emitted %d chunk(s) before the abort, want 2", chunks)
	}
	t.Logf("TtsEngine: 2 ordered chunks, rate %d, abort verbatim %q", tts.SampleRate(), err)

	kws := &stubKws{}
	var hits []KeywordHit
	// Exactly as many windows as the MinSpeech gate needs, rounded up: the stub
	// fires on the window that crosses it and never again.
	for i := 0; i < (minSpeechSamples+WindowSamples-1)/WindowSamples; i++ {
		hits = kws.Accept(make([]float32, WindowSamples), hits)
	}
	if len(hits) != 1 {
		t.Fatalf("WakeWordEngine produced %d hit(s), want 1", len(hits))
	}
	if hits[0].At < minSpeechSamples {
		t.Fatalf("hit position %d samples is before the audio it was computed from (%d)",
			hits[0].At, kws.fed)
	}
	t.Logf("WakeWordEngine: hit %q veto=%v at %d samples", hits[0].Word, hits[0].Veto, hits[0].At)

	rt := &stubRealtime{}
	var handoffs, endpoints int
	for _, ev := range rt.Send(make([]float32, WindowSamples), nil) {
		switch ev.Kind {
		case RealtimeHandoff:
			handoffs++
		case RealtimeEndpoint:
			endpoints++
		}
	}
	if handoffs != 1 || endpoints != 1 {
		t.Fatalf("RealtimeEngine produced %d handoff(s) and %d endpoint(s), want 1 each",
			handoffs, endpoints)
	}
	if err := rt.Suspend(); err != nil {
		t.Fatalf("Suspend: %v", err)
	}
	if got := rt.Send(make([]float32, WindowSamples), nil); len(got) != 0 {
		t.Fatalf("a suspended session emitted %d event(s), want none: C32 (1) parks the "+
			"session while the text loop works", len(got))
	}
	if err := rt.Resume(); err != nil {
		t.Fatalf("Resume: %v", err)
	}
	t.Logf("RealtimeEngine: handoff=%d endpoint=%d, suspended stream silent, resumed",
		handoffs, endpoints)
}

// cloudStubAsr is an AsrEngine that shares no code with StreamingAsr: it counts
// samples and spells out a fixed text. It is deliberately not a model, and it is
// not cited as evidence about recognition. What it proves is the SEAM: the 15-c
// driver runs over an engine it has never seen, so C9's "local sherpa and cloud
// provider interchangeable" is a property of the code rather than a sentence in
// a plan. Replacing this stub with a real provider is ticket 61's work and
// changes nothing above it.
type cloudStubAsr struct {
	open bool
	fed  int
	text []rune
}

// cloudStubText is the whole hypothesis this stub can ever produce.
const cloudStubText = "far-side transcript"

func (c *cloudStubAsr) Begin() { c.open, c.fed, c.text = true, 0, nil }

func (c *cloudStubAsr) Feed(samples []float32) {
	if !c.open {
		return
	}
	c.fed += len(samples)
	if n := len(c.text); n < len([]rune(cloudStubText)) {
		c.text = append(c.text, []rune(cloudStubText)[n])
	}
}

func (c *cloudStubAsr) Partial() string {
	if !c.open {
		return ""
	}
	return string(c.text)
}

func (c *cloudStubAsr) Finish() string {
	if !c.open {
		return ""
	}
	c.open = false
	return string(c.text)
}

func (c *cloudStubAsr) InFlight() bool { return c.open }
func (c *cloudStubAsr) Fed() int       { return c.fed }
func (c *cloudStubAsr) Close() error   { c.open = false; return nil }

var _ AsrEngine = (*cloudStubAsr)(nil)

// TestTheTranscriberRunsOverAnyAsrEngine is the interchangeability reading: the
// real endpointer (silero, staging) plus a non-sherpa engine produce the same
// shape of transcript stream that the sherpa session produces -- partials while
// the utterance is open, one final when it closes, and a latency that comes off
// the audio clock rather than from whichever engine is plugged in.
func TestTheTranscriberRunsOverAnyAsrEngine(t *testing.T) {
	det, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
	if err != nil {
		t.Skipf("VAD session unavailable on this host: %v", err)
	}
	t.Cleanup(func() { _ = det.Close() })

	eng := &cloudStubAsr{}
	tr := NewTranscriber(det, eng)

	// The 15-b calibrated fixture: 1.2s of speech, then the 1500ms nominal
	// silence that actually closes an utterance (silero over-reports speech by
	// roughly 400ms).
	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	var out []Transcript
	for f := 0; f < len(pcm)/WindowSamples; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	partials, finals := countTranscripts(out)
	if len(out) == 0 {
		t.Fatalf("the non-sherpa engine produced no transcripts at all")
	}
	last := out[len(out)-1]
	t.Logf("non-sherpa engine over the real endpointer: %d partial(s), %d final(s), "+
		"final %q closed by %v over %d samples (%v)",
		partials, finals, last.Text, last.Endpoint.Kind, last.Endpoint.Total, last.Latency)

	if finals != 1 {
		t.Fatalf("%d final(s), want 1: the endpointer's split is not reaching this engine", finals)
	}
	if partials == 0 {
		t.Fatalf("no partials while the utterance was open: the driver is not streaming")
	}
	if last.Kind != TranscriptFinal || last.Text != cloudStubText {
		t.Fatalf("last transcript = %v %q, want the engine's own text %q",
			last.Kind, last.Text, cloudStubText)
	}
	if last.Endpoint.Kind != EndpointEnd || last.Endpoint.Total <= 0 {
		t.Fatalf("the final carries %v over %d samples, want a real closed utterance",
			last.Endpoint.Kind, last.Endpoint.Total)
	}
	if want := samplesDuration(eng.Fed()); last.Latency != want {
		t.Fatalf("latency %v is not the audio clock reading (%v over %d samples): a latency "+
			"that depends on the engine cannot be compared across engines (D42#9)",
			last.Latency, want, eng.Fed())
	}
}
