package speech

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
)

// This file is the ticket's measurement harness: the readings it produces are
// what docs/evidence/s1/15c-r1-impl.md reports. Everything is synthesised --
// ticket 15d owns the recorded baseline sets, and 15c must not block on a
// corpus (same rule 15b-r1 set for endpointing).

// synthBandNoise is white noise band-limited by a difference of two one-pole
// filters. It gives the recognizer something with speech-band energy and no
// phonetic content.
func synthBandNoise(n int, lo, hi, amp float64, seed int64) []int16 {
	rng := rand.New(rand.NewSource(seed))
	aHi := math.Exp(-2 * math.Pi * hi / float64(audio.TargetRate))
	aLo := math.Exp(-2 * math.Pi * lo / float64(audio.TargetRate))
	lpHi, lpLo := 0.0, 0.0
	out := make([]int16, n)
	for i := range out {
		x := rng.Float64()*2 - 1
		lpHi = aHi*lpHi + (1-aHi)*x
		lpLo = aLo*lpLo + (1-aLo)*x
		out[i] = int16(math.Max(-1, math.Min(1, amp*(lpHi-lpLo)*4)) * 32767)
	}
	return out
}

// synthFormant is a glottal harmonic stack shaped by three formant resonances:
// a vowel, at the same 3.5Hz syllable rate as synthSpeech. Its nominal phonetic
// script is known (one syllable per envelope period), which is what makes a
// word-error reading possible without recorded speech.
func synthFormant(n int, f0 float64, formants []float64, amp float64) []int16 {
	out := make([]int16, n)
	for i := range out {
		tt := float64(i) / float64(audio.TargetRate)
		var acc float64
		for h := 1; h <= 40; h++ {
			f := f0 * float64(h)
			if f > 7000 {
				break
			}
			var g float64
			for _, fc := range formants {
				g += 1 / (1 + math.Pow((f-fc)/90, 2))
			}
			acc += g * math.Sin(2*math.Pi*f*tt) / math.Sqrt(float64(h))
		}
		env := 0.5 + 0.5*math.Sin(2*math.Pi*surrogateSyllableHz*tt)
		out[i] = int16(math.Max(-1, math.Min(1, amp*env*acc)) * 32767)
	}
	return out
}

// nominalSyllables is how many syllables the envelope asked for in n samples.
func nominalSyllables(n int) int {
	return int(surrogateSyllableHz * float64(n) / float64(audio.TargetRate))
}

// TestPartialsStreamWhileTheUtteranceIsOpen is the ticket's central reading:
// partials have to arrive DURING the utterance, with a first-partial latency
// inside the streaming budget (SPEC-04 sec 8: 800ms; the offline engines get
// 1500ms, which is why streaming is the engine that ships).
//
// The latency is read off the audio clock (samples fed), not off a wall clock,
// so the same fixture reports the same number on every host (D42#9).
func TestPartialsStreamWhileTheUtteranceIsOpen(t *testing.T) {
	tr, asr := newTestTranscriber(t)

	pcm := synthSpeech(msSamples(3000))
	frames := len(pcm) / WindowSamples
	var partials []Transcript
	var out []Transcript
	out = tr.Accept(nil, out)

	for f := 0; f < frames; f++ {
		win := pcm[f*WindowSamples : (f+1)*WindowSamples]
		before := len(out)
		out = tr.Accept(pcmToFloat(win), out)
		for _, got := range out[before:] {
			if got.Kind != TranscriptPartial {
				continue
			}
			partials = append(partials, got)
			t.Logf("frame=%d (fixture t=%dms) -> partial #%d %q (utterance t=%dms, fed=%d samples)",
				f+1, (f+1)*int(audio.FrameDuration/time.Millisecond), len(partials), got.Text,
				int(got.Latency/time.Millisecond), asr.Fed())
		}
	}

	if len(partials) < 2 {
		t.Fatalf("%d partial(s) over 3s of speech, want at least 2: the hypothesis has to "+
			"grow while the audio is arriving, or the ball cannot start thinking early (AC#1)",
			len(partials))
	}
	first := partials[0]
	if first.Latency > StreamingLatencyBudget {
		t.Fatalf("first partial arrived %v after the utterance opened, over the %v streaming "+
			"budget (SPEC-04 sec 8)", first.Latency, StreamingLatencyBudget)
	}
	if first.Latency <= 0 {
		t.Fatalf("first partial reports a non-positive latency (%v)", first.Latency)
	}
	for i := 1; i < len(partials); i++ {
		if len([]rune(partials[i].Text)) < len([]rune(partials[i-1].Text)) {
			t.Fatalf("partial #%d (%q) is shorter than partial #%d (%q): the hypothesis "+
				"rewound inside one utterance", i+1, partials[i].Text, i, partials[i-1].Text)
		}
	}

	out = tr.Finalize(out)
	last := out[len(out)-1]
	if last.Kind != TranscriptFinal {
		t.Fatalf("the run ended on a %v, want the final", last.Kind)
	}
	if last.Text == "" {
		t.Fatalf("the final transcript is empty")
	}
	t.Logf("summary: %d partials, first at %v (budget %v, margin %v), final %q over %d samples",
		len(partials), first.Latency, StreamingLatencyBudget,
		StreamingLatencyBudget-first.Latency, last.Text, last.Endpoint.Total)
}

// runToFinal feeds whole windows until the first final transcript appears and
// returns it with the wall-clock cost of the call that produced it. That cost
// IS the VAD-stop -> final leg SPEC-04 sec 8 budgets.
//
// The clock on the dev host quantises to roughly 0.57ms (measured: consecutive
// time.Now() deltas are 0 or 568500ns), so a single-call reading of 0 means
// "under one tick", not "no work". The instrument's own control is the
// whole-segment mutation, which reads 321ms on the same line of code.
func runToFinal(tr *Transcriber, pcm []int16) (Transcript, time.Duration) {
	frames := len(pcm) / WindowSamples
	for f := 0; f < frames; f++ {
		start := time.Now()
		out := tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), nil)
		cost := time.Since(start)
		for _, got := range out {
			if got.Kind == TranscriptFinal {
				return got, cost
			}
		}
	}
	return Transcript{}, 0
}

// TestFinalArrivesWithinTheBudgetOfTheVadStop is the budget as SPEC-04 sec 8
// actually words it: from the moment the VAD closes the utterance to the final
// transcript. That leg is real compute (flush the feature extractor, decode
// what is left), so it is the one number measured on a wall clock.
func TestFinalArrivesWithinTheBudgetOfTheVadStop(t *testing.T) {
	tr, _ := newTestTranscriber(t)

	// 1.2s of speech, then the 1500ms pause 15-b calibrated as the shortest
	// nominal silence that actually closes an utterance (silero over-reports
	// speech by roughly 400ms).
	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	final, stopCost := runToFinal(tr, pcm)
	if final.Kind != TranscriptFinal {
		t.Fatalf("the utterance never closed")
	}
	if final.Text == "" {
		t.Fatalf("the final transcript is empty")
	}
	note := ""
	if stopCost < time.Millisecond {
		note = " (below this host's ~0.57ms clock tick)"
	}
	t.Logf("VAD-stop -> final: %d ns%s, budget %v; text %q over %d samples",
		stopCost.Nanoseconds(), note, StreamingLatencyBudget, final.Text, final.Endpoint.Total)
	if stopCost > StreamingLatencyBudget {
		t.Fatalf("VAD-stop to final took %v, over the %v streaming budget (SPEC-04 sec 8)",
			stopCost, StreamingLatencyBudget)
	}
}

// TestSilentInputProducesNoTranscripts is the positive control: three seconds of
// digital silence through the whole chain must produce no endpoint, no partial
// and no final. A recognizer that invents a word here would be a
// hallucinating product, and AC#5's "did not catch that" path depends on
// silence being silent.
func TestSilentInputProducesNoTranscripts(t *testing.T) {
	tr, _ := newTestTranscriber(t)

	pcm := silence(msSamples(3000))
	frames := len(pcm) / WindowSamples
	var out []Transcript
	for f := 0; f < frames; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	partials, finals := countTranscripts(out)
	t.Logf("3s of digital silence: %d partial(s), %d final(s), %d transcript(s) total",
		partials, finals, len(out))
	if len(out) != 0 {
		t.Fatalf("silence produced %d transcript(s) (%+v), want none", len(out), out)
	}
}

// TestStreamedSamplesMatchTheClosedUtterance pins the invariant that makes the
// seam trustworthy: the recognizer streamed exactly the audio the endpoint
// closed. Endpointer trims the samples that arrived after the closing threshold
// and truncates at MaxUtterance, so the allowed slack is one window.
func TestStreamedSamplesMatchTheClosedUtterance(t *testing.T) {
	tr, asr := newTestTranscriber(t)

	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	final, _ := runToFinal(tr, pcm)
	if final.Kind != TranscriptFinal {
		t.Fatalf("the utterance never closed")
	}
	fed, total := asr.Fed(), final.Endpoint.Total
	t.Logf("streamed %d samples, endpoint reports %d (slack %d, window %d)",
		fed, total, fed-total, WindowSamples)
	if total <= 0 {
		t.Fatalf("the closed utterance carries no audio")
	}
	if fed < total || fed-total >= WindowSamples {
		t.Fatalf("the ASR streamed %d samples but the endpoint closed an utterance of %d: "+
			"the recognizer and the endpointer disagree about what the utterance was", fed, total)
	}
}

// TestSpeakPauseSpeakYieldsTwoFinals is the 15-b seam integration: the
// endpointer's split, not the recognizer's, decides where utterances end, and
// each utterance gets its own final with its own partials.
func TestSpeakPauseSpeakYieldsTwoFinals(t *testing.T) {
	tr, _ := newTestTranscriber(t)

	// 15-b calibrated this fixture by measurement: silero over-reports speech by
	// roughly 400ms, so a nominal 1100ms pause is the first one that closes an
	// utterance. 1500ms leaves 400ms of margin.
	pcm := append(synthSpeech(msSamples(1200)), silence(msSamples(1500))...)
	pcm = append(pcm, synthSpeech(msSamples(1200))...)

	frames := len(pcm) / WindowSamples
	var out []Transcript
	for f := 0; f < frames; f++ {
		out = tr.Accept(pcmToFloat(pcm[f*WindowSamples:(f+1)*WindowSamples]), out)
	}
	out = tr.Finalize(out)

	partials, finals := countTranscripts(out)
	t.Logf("speak-pause-speak: %d partial(s), %d final(s)", partials, finals)
	for i, got := range out {
		if got.Kind != TranscriptFinal {
			continue
		}
		t.Logf("  final %d: kind=%v %q over %d samples (%v of audio)",
			i, got.Endpoint.Kind, got.Text, got.Endpoint.Total, got.Latency)
	}
	if finals != 2 {
		t.Fatalf("%d final(s) over speak-pause-speak, want 2: the endpointer's split is not "+
			"reaching the recognizer", finals)
	}
	if partials == 0 {
		t.Fatalf("no partials were produced while the utterances were open")
	}
	if out[len(out)-1].Text == "" {
		t.Fatalf("the second utterance produced an empty transcript")
	}
}

// wordErrors is the edit script of one hypothesis against one reference.
type wordErrors struct {
	sub, del, ins int
	refLen        int
}

// rate is the word error rate, undefined-as-zero when there is no reference.
func (w wordErrors) rate() float64 {
	if w.refLen == 0 {
		if w.ins == 0 {
			return 0
		}
		return float64(w.ins)
	}
	return float64(w.sub+w.del+w.ins) / float64(w.refLen)
}

// alignWords is a plain Levenshtein alignment over runes, with a traceback so
// the substitution/deletion/insertion split is exact. It exists to make the
// table below arithmetically checkable; the CER harness itself is ticket 15d's,
// and this is not it.
func alignWords(ref, hyp []rune) wordErrors {
	w := wordErrors{refLen: len(ref)}
	dp := make([][]int, len(ref)+1)
	for i := range dp {
		dp[i] = make([]int, len(hyp)+1)
		dp[i][0] = i
	}
	for j := 0; j <= len(hyp); j++ {
		dp[0][j] = j
	}
	for i := 1; i <= len(ref); i++ {
		for j := 1; j <= len(hyp); j++ {
			cost := 1
			if ref[i-1] == hyp[j-1] {
				cost = 0
			}
			dp[i][j] = min3(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+cost)
		}
	}
	for i, j := len(ref), len(hyp); i > 0 || j > 0; {
		switch {
		case i > 0 && j > 0 && ref[i-1] == hyp[j-1] && dp[i][j] == dp[i-1][j-1]:
			i--
			j--
		case i > 0 && j > 0 && dp[i][j] == dp[i-1][j-1]+1:
			w.sub++
			i--
			j--
		case i > 0 && dp[i][j] == dp[i-1][j]+1:
			w.del++
			i--
		default:
			w.ins++
			j--
		}
	}
	return w
}

func min3(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}

// TestSyntheticWordErrorRate is the accuracy shape on synthetic fixtures.
//
// It is deliberately NOT a lexical accuracy claim: a recorded corpus is 15d's,
// and ticket 15's own text says the committed baseline sets are the thing that
// makes a CER number mean anything. What a synthetic fixture CAN pin is the
// link between what the synthesizer was asked to say and what came out:
//
//   - the two formant vowels have a known script (one vowel per envelope
//     period), so insertions/deletions are real errors;
//   - the non-phonetic fixtures (silence, band-limited noise, the harmonic
//     stack the VAD accepts) have no script at all, so every token emitted is
//     an insertion -- which is the hallucination rate the product cares about.
func TestSyntheticWordErrorRate(t *testing.T) {
	asr := newTestAsr(t)

	const seconds = 3
	n := msSamples(seconds * 1000)
	type fixture struct {
		name string
		pcm  []int16
		ref  []rune // nil means "no nominal script": any token is an insertion
	}
	fixtures := []fixture{
		{"A-harmonic-stack", synthSpeech(n), nil},
		{"B-silence", silence(n), nil},
		{"C-white-noise", synthBandNoise(n, 20, 7000, 0.3, 1), nil},
		{"D-speechband-noise", synthBandNoise(n, 250, 3600, 0.9, 2), nil},
		{
			"E-vowel-a", synthFormant(n, 120, []float64{700, 1220, 2600}, 0.12),
			repeatRune('啊', nominalSyllables(n)),
		},
		{
			"F-vowel-i", synthFormant(n, 130, []float64{300, 2200, 3000}, 0.12),
			repeatRune('一', nominalSyllables(n)),
		},
	}

	hyps := make(map[string]string, len(fixtures))
	var scripted, unscripted wordErrors
	for _, f := range fixtures {
		asr.Begin()
		for i := 0; i+WindowSamples <= len(f.pcm); i += WindowSamples {
			asr.Feed(pcmToFloat(f.pcm[i : i+WindowSamples]))
		}
		hyp := asr.Finish()
		hyps[f.name] = hyp
		got := alignWords(f.ref, []rune(hyp))
		if f.ref == nil {
			unscripted.ins += got.ins
			t.Logf("%-19s ref=%s hyp=%-12q tokens=%d (all insertions: the fixture has no script)",
				f.name, "(none)", hyp, got.ins)
			continue
		}
		scripted.sub += got.sub
		scripted.del += got.del
		scripted.ins += got.ins
		scripted.refLen += got.refLen
		t.Logf("%-19s ref=%q hyp=%-12q S=%d D=%d I=%d WER=%.2f",
			f.name, string(f.ref), hyp, got.sub, got.del, got.ins, got.rate())
	}
	t.Logf("scripted (E,F): S=%d D=%d I=%d over %d nominal words -> WER=%.2f",
		scripted.sub, scripted.del, scripted.ins, scripted.refLen, scripted.rate())
	t.Logf("unscripted (A-D): %d tokens emitted against 0 reference words", unscripted.ins)

	// The reading that has to hold no matter what the model does: fixtures with
	// no phonetic content must decode to nothing.
	for _, name := range []string{"B-silence", "C-white-noise", "D-speechband-noise"} {
		if hyp := hyps[name]; hyp != "" {
			t.Fatalf("%s decoded to %q: a fixture with no phonetic content must not produce "+
				"words", name, hyp)
		}
	}
}

func repeatRune(r rune, n int) []rune {
	out := make([]rune, n)
	for i := range out {
		out[i] = r
	}
	return out
}
