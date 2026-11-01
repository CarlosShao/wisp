package speech

import "testing"

// Ticket 398 probes: the silero RNN hangover and the session reset face.
//
// The mechanism (392-r2 mechanism probe, three-process reproduced in 392-v3):
// the RNN context carries across silent gaps because an idle session gets no
// input and nothing decays, so after Finalize closes the open utterance the
// hangover can still feed the endpointer's speech run through post-speech
// silence - the next start lands on the SILENT gap rather than on the talker.
// 392-r2 read that control arm as (silence=1 talker=0). Until this ticket the
// only owner-side cure was feeding 1.2s of cold silence and discarding what
// fired in it; Reset is the harder clear, and these probes pin it.
//
// The readings keep 392-r2's named shape (silence=N talker=M). The fixtures
// are the package's measured surrogate: silero hangs over for ~720ms after it
// stops (vad_test.go), so a 400ms gap sits wholly inside the hangover window
// while staying under the 700ms close threshold. The assertions are pinned to
// measured shape, and the reset arm doubles as the ticket's stop-the-line
// check: a talker=0 there would mean the reset broke the classifier itself,
// which the ticket forbids (reset clears state, never judgement).

// count398Starts feeds whole windows and counts the starts that fire.
func count398Starts(d *Detector, samples []float32) int {
	n := 0
	var eps []Endpoint
	for off := 0; off+WindowSamples <= len(samples); off += WindowSamples {
		eps = d.Accept(samples[off:off+WindowSamples], eps[:0])
		for _, ep := range eps {
			if ep.Kind == EndpointStart {
				n++
			}
		}
	}
	return n
}

// Test398SileroHangoverReset pins the ticket's two readings on one measured
// shape. Pre-reset (the 392-r2 control): Finalize returns the endpointer to
// idle but the RNN still carries the talker's context, so the start lands on
// the silent gap (silence=1 talker=0). Post-reset: the gap is clean and the
// start lands on the talker it belongs to (silence=0 talker=1).
func Test398SileroHangoverReset(t *testing.T) {
	talker := pcmToFloat(synthSpeech(msSamples(1500)))
	gap := pcmToFloat(silence(msSamples(400)))

	// Pre-reset: the hangover reading.
	d := newTestDetector(t)
	if n := count398Starts(d, talker); n == 0 {
		t.Fatalf("the probe talker did not open an utterance: the probe has no subject")
	}
	_ = d.Finalize(nil)
	nSil := count398Starts(d, gap)
	nTalk := count398Starts(d, talker)
	if nSil != 1 || nTalk != 0 {
		t.Fatalf("MECHANISM MOVED: pre-reset reading is (silence=%d talker=%d), want (1 0); "+
			"the hangover no longer reproduces on this fixture, so the reset readings have no control",
			nSil, nTalk)
	}
	t.Logf("398-PRE-RESET (silence=%d talker=%d): the start fired on post-speech silence", nSil, nTalk)

	// Post-reset: the same sequence with Reset between the windows.
	d2 := newTestDetector(t)
	if n := count398Starts(d2, talker); n == 0 {
		t.Fatalf("the reset-arm talker did not open an utterance")
	}
	_ = d2.Finalize(nil)
	d2.Reset()
	nSil = count398Starts(d2, gap)
	nTalk = count398Starts(d2, talker)
	if nSil != 0 || nTalk != 1 {
		t.Fatalf("MECHANISM MOVED: post-reset reading is (silence=%d talker=%d), want (0 1); "+
			"silence>0 means the hangover survived the reset, talker=0 means the reset reached the "+
			"classifier, which ticket 398 forbids (stop the line)",
			nSil, nTalk)
	}
	t.Logf("398-POST-RESET (silence=%d talker=%d): the gap is clean and the start landed on the talker", nSil, nTalk)
}

// Test398ResetMatchesAFreshDetector states the reset contract directly: after
// a Reset a detector that has just finished an utterance reads exactly like a
// just-constructed detector on the same three-material run. The readings are
// per-material start counts, which are landing-sensitive - a kind-only
// sequence would hide a hangover start that merely moved onto the gap.
func Test398ResetMatchesAFreshDetector(t *testing.T) {
	warmUp := pcmToFloat(synthSpeech(msSamples(1500)))
	talker := pcmToFloat(synthSpeech(msSamples(1500)))
	gap := pcmToFloat(silence(msSamples(400)))

	counts := func(d *Detector) (int, int, int) {
		return count398Starts(d, gap), count398Starts(d, talker), count398Starts(d, gap)
	}

	fresh := newTestDetector(t)
	wantGapA, wantTalk, wantGapB := counts(fresh)
	if wantTalk == 0 {
		t.Fatalf("the fresh detector never started on the talker: the probe has no subject")
	}
	if wantGapA != 0 || wantGapB != 0 {
		t.Fatalf("MECHANISM MOVED: the fresh detector fired on silence (gapA=%d gapB=%d), want 0/0", wantGapA, wantGapB)
	}

	dirty := newTestDetector(t)
	if n := count398Starts(dirty, warmUp); n == 0 {
		t.Fatalf("the warm-up talker did not open an utterance")
	}
	_ = dirty.Finalize(nil)
	dirty.Reset()
	gotGapA, gotTalk, gotGapB := counts(dirty)
	if gotGapA != wantGapA || gotTalk != wantTalk || gotGapB != wantGapB {
		t.Fatalf("reset detector reading is (gapA=%d talker=%d gapB=%d), fresh detector is (%d %d %d): the reset did not return the detector to its fresh state",
			gotGapA, gotTalk, gotGapB, wantGapA, wantTalk, wantGapB)
	}
	t.Logf("398-RESET-EQUALS-FRESH (gapA=%d talker=%d gapB=%d) on both arms", gotGapA, gotTalk, gotGapB)
}
