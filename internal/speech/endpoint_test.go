package speech

import (
	"fmt"
	"strings"
	"testing"
)

// The state machine owns the three AC#3 numbers, so it is tested on an audio
// clock with sample-exact fixtures: no cgo, no model, no corpus. Every
// boundary is probed on BOTH sides (699 vs 701 etc.) because a single-sided
// reading cannot tell "the constant is 700" from "the constant is 50".

// machine drives an Endpointer over an audio clock, recording when each
// transition fired.
type machine struct {
	ep  *Endpointer
	out []Endpoint
	at  []int // audio-clock ms at which out[i] was emitted
	ms  int
}

func newMachine() *machine { return &machine{ep: NewEndpointer()} }

// feed advances the machine by steps x stepMS milliseconds of one class.
func (m *machine) feed(stepMS, steps int, speech bool) {
	win := make([]float32, msSamples(stepMS))
	if speech {
		for i := range win {
			win[i] = 0.1
		}
	}
	for i := 0; i < steps; i++ {
		before := len(m.out)
		m.out = m.ep.Accept(win, speech, m.out)
		m.ms += stepMS
		for range m.out[before:] {
			m.at = append(m.at, m.ms)
		}
	}
}

// ofKind returns the endpoints of one kind.
func (m *machine) ofKind(k EndpointKind) []Endpoint {
	var out []Endpoint
	for _, e := range m.out {
		if e.Kind == k {
			out = append(out, e)
		}
	}
	return out
}

// trace renders the transition log for the evidence table.
func (m *machine) trace() string {
	var b strings.Builder
	for i, e := range m.out {
		if i > 0 {
			b.WriteString(" ")
		}
		fmt.Fprintf(&b, "%s@%dms", e.Kind, m.at[i])
		if e.Kind != EndpointStart {
			fmt.Fprintf(&b, "(total=%dms speech=%dms)", e.Duration().Milliseconds(), e.SpeechDuration().Milliseconds())
		}
	}
	if b.Len() == 0 {
		return "(no transitions)"
	}
	return b.String()
}

// TestEndpointSilenceBoundary is AC#3's "speak-pause-speak splits correctly"
// boundary: a pause of EndpointSilence (700ms) closes the utterance, and the
// two sides differ.
func TestEndpointSilenceBoundary(t *testing.T) {
	for _, tc := range []struct {
		silenceMS int
		wantEnds  int
	}{
		{699, 0}, // just under: still one utterance
		{700, 1}, // exactly at the constant: inclusive by design
		{701, 1}, // just over
	} {
		m := newMachine()
		m.feed(1, 1000, true)          // 1.0s of speech
		m.feed(1, tc.silenceMS, false) // the pause
		ends := m.ofKind(EndpointEnd)
		if len(ends) != tc.wantEnds {
			t.Fatalf("pause %dms: got %d end(s), want %d; trace=%s",
				tc.silenceMS, len(ends), tc.wantEnds, m.trace())
		}
		if tc.wantEnds == 1 {
			// It must close at exactly 1000+700ms, i.e. the constant is 700.
			if got := m.at[len(m.at)-1]; got != 1700 {
				t.Fatalf("pause %dms: closed at %dms, want 1700ms; trace=%s",
					tc.silenceMS, got, m.trace())
			}
			// 1.0s speech + 0.7s trailing silence.
			if d := ends[0].Duration().Milliseconds(); d != 1700 {
				t.Fatalf("pause %dms: utterance %dms, want 1700ms", tc.silenceMS, d)
			}
			if d := ends[0].SpeechDuration().Milliseconds(); d != 1000 {
				t.Fatalf("pause %dms: speech %dms, want 1000ms", tc.silenceMS, d)
			}
		}
		if tc.wantEnds == 0 && !m.ep.InSpeech() {
			t.Fatalf("pause %dms: machine left the speaking state, want still speaking", tc.silenceMS)
		}
	}
}

// TestMinSpeechBoundary is AC#3's "<300ms speech ignored" boundary.
func TestMinSpeechBoundary(t *testing.T) {
	for _, tc := range []struct {
		burstMS   int
		wantStart int
	}{
		{299, 0},
		{300, 1},
		{301, 1},
	} {
		m := newMachine()
		m.feed(1, 200, false)       // 200ms of lead-in silence
		m.feed(1, tc.burstMS, true) // the burst
		m.feed(1, 1500, false)      // plenty of silence to close anything opened
		if got := len(m.ofKind(EndpointStart)); got != tc.wantStart {
			t.Fatalf("burst %dms: %d utterance(s) opened, want %d; trace=%s",
				tc.burstMS, got, tc.wantStart, m.trace())
		}
		if tc.wantStart == 0 && len(m.out) != 0 {
			t.Fatalf("burst %dms: want no transition at all, got %s", tc.burstMS, m.trace())
		}
	}
}

// TestMaxUtteranceBoundary is AC#3's "60s utterance auto-finalizes" boundary.
func TestMaxUtteranceBoundary(t *testing.T) {
	const speechS = 60_000 // ms

	// Side 1: 59.9s of speech is not truncated.
	m := newMachine()
	m.feed(1, 59_900, true)
	if n := len(m.ofKind(EndpointMaxDuration)); n != 0 {
		t.Fatalf("59.9s: got %d max-duration close(s), want 0; trace=%s", n, m.trace())
	}
	if !m.ep.InSpeech() {
		t.Fatalf("59.9s: machine closed early; trace=%s", m.trace())
	}
	// Side 2: 100ms more crosses the cap, at exactly 60.000s.
	m.feed(1, 200, true)
	forced := m.ofKind(EndpointMaxDuration)
	if len(forced) != 1 {
		t.Fatalf("60.0s+: got %d max-duration close(s), want 1; trace=%s", len(forced), m.trace())
	}
	if got := m.at[len(m.at)-1]; got != speechS {
		t.Fatalf("60.0s+: closed at %dms, want %dms; trace=%s", got, speechS, m.trace())
	}
	if d := forced[0].Duration().Milliseconds(); d != speechS {
		t.Fatalf("60.0s+: utterance reported %dms, want the %dms cap", d, speechS)
	}
	if n := len(forced[0].Samples); n != maxUtteranceSamples {
		t.Fatalf("60.0s+: carried %d samples, want the %d sample cap", n, maxUtteranceSamples)
	}

	// Uninterrupted 60.1s: one forced close, and the 100ms tail is too short to
	// open a second utterance.
	m2 := newMachine()
	m2.feed(1, 60_100, true)
	if n := len(m2.ofKind(EndpointMaxDuration)); n != 1 {
		t.Fatalf("60.1s: got %d max-duration close(s), want 1; trace=%s", n, m2.trace())
	}
	if n := len(m2.ofKind(EndpointStart)); n != 1 {
		t.Fatalf("60.1s: got %d start(s), want 1; trace=%s", n, m2.trace())
	}
}

// TestPureSilenceProducesNoEndpoint is the positive control: silence is not a
// sentence. It is also the assertion that catches a zeroed MinSpeech, because
// a zero threshold would let the first silent window open an utterance.
func TestPureSilenceProducesNoEndpoint(t *testing.T) {
	m := newMachine()
	m.feed(1, 5000, false)
	m.out = m.ep.Finalize(m.out)
	if len(m.out) != 0 {
		t.Fatalf("5s of pure silence produced %d transition(s), want 0; trace=%s", len(m.out), m.trace())
	}
	if m.ep.InSpeech() {
		t.Fatalf("pure silence left the machine in the speaking state")
	}
}

// TestSpeakPauseSpeakSplits is AC#3's first clause at the state-machine level:
// two speech runs separated by more than EndpointSilence become two utterances.
func TestSpeakPauseSpeakSplits(t *testing.T) {
	m := newMachine()
	m.feed(1, 400, false)
	m.feed(1, 500, true)
	m.feed(1, 800, false)
	m.feed(1, 500, true)
	m.feed(1, 800, false)

	if got := len(m.ofKind(EndpointStart)); got != 2 {
		t.Fatalf("got %d start(s), want 2; trace=%s", got, m.trace())
	}
	ends := m.ofKind(EndpointEnd)
	if len(ends) != 2 {
		t.Fatalf("got %d end(s), want 2; trace=%s", len(ends), m.trace())
	}
	for i, e := range ends {
		if d := e.SpeechDuration().Milliseconds(); d != 500 {
			t.Fatalf("utterance %d: speech %dms, want 500ms; trace=%s", i, d, m.trace())
		}
		if d := e.Duration().Milliseconds(); d != 1200 {
			t.Fatalf("utterance %d: total %dms, want 500+700=1200ms; trace=%s", i, d, m.trace())
		}
		if n := len(e.Samples); n != msSamples(1200) {
			t.Fatalf("utterance %d: carried %d samples, want %d", i, n, msSamples(1200))
		}
	}
	// The 800ms pause must not be swallowed whole into either utterance, and no
	// silence may be claimed twice: utterance 1 runs 300..1600ms and utterance
	// 2 runs 1700..2900ms, so the two closures fire at 1600ms and 2900ms.
	var endAt []int
	for i, e := range m.out {
		if e.Kind == EndpointEnd {
			endAt = append(endAt, m.at[i])
		}
	}
	if len(endAt) != 2 || endAt[0] != 1600 || endAt[1] != 2900 {
		t.Fatalf("closures fired at %v, want [1600 2900]; trace=%s", endAt, m.trace())
	}
}

// TestStateMachineIgnoresChunking: the machine counts samples, so the same
// audio yields the same closures whether it arrives in 1ms steps or in one
// chunk per region.
func TestStateMachineIgnoresChunking(t *testing.T) {
	regions := []struct {
		ms     int
		speech bool
	}{{400, false}, {500, true}, {800, false}, {500, true}, {800, false}}

	m1 := newMachine()
	for _, r := range regions {
		m1.feed(1, r.ms, r.speech)
	}

	m2 := newMachine()
	for _, r := range regions {
		m2.feed(r.ms, 1, r.speech)
	}

	a, b := m1.ofKind(EndpointEnd), m2.ofKind(EndpointEnd)
	if len(a) != 2 || len(b) != 2 {
		t.Fatalf("closures differ: 1ms-steps=%d chunked=%d; traces=%s | %s",
			len(a), len(b), m1.trace(), m2.trace())
	}
	for i := range a {
		if a[i].Duration() != b[i].Duration() || a[i].SpeechDuration() != b[i].SpeechDuration() {
			t.Fatalf("closure %d differs: 1ms-steps=%s chunked=%s", i, m1.trace(), m2.trace())
		}
	}
}
