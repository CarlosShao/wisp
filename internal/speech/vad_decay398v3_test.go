package speech

import "testing"

// Ticket 398 residual (c) probes: how the silero RNN hangover decays across
// fed silence, measured as a sweep over the number of empty gaps between the
// closed utterance and the 398-r1 probe. 392-v3 sec 7.4 left "how much decay
// is enough" at a single 1.2s cold-segment reading; 398-r1/v2 landed the
// Reset face. This file is pure measurement - it adds no product code and
// pins no new curve shape - and its readings feed the keepalive-default
// decision (392-v3 sec 6 step 3): whether the owner can drop the cold
// segment in favour of Finalize+Reset, and how much fed silence the natural
// decay actually needs.
//
// The sweep reuses 398-r1's measured fixtures verbatim: the 1500ms surrogate
// talker (900ms does not clear onset+MinSpeech) and the 400ms gap (under the
// 700ms close threshold, inside the ~720ms hangover window noted in
// vad_test.go). One trial at depth n:
//
//	talker 1500ms -> Finalize -> [arm 2: Reset] gap 400ms x n -> Finalize
//	  -> probe gap 400ms -> probe talker 1500ms
//
// The readings keep 392-r2's named (silence=N talker=M) shape: nSil counts
// starts on the probe gap (a start there is the hangover still alive after n
// gaps), nTalk counts starts on the probe talker (the next speech opens
// cleanly). The pre-probe Finalize closes whatever the decay feed left open
// and nothing else - Finalize never reaches the RNN (vad.go) - so the probe
// reads exactly the RNN state the gaps left behind. Arm 1 (natural) feeds the
// gaps with no Reset; arm 2 (zeroed control) Resets before every gap, the
// 398-r1 POST-RESET shape repeated per gap.

// sweep398v3Trial runs one sweep trial at decay depth n and returns what
// fired. resetEachGap selects the arm: false feeds the gaps as they come
// (natural decay); true Resets before every gap (the zeroed control).
// gapStarts[k]/gapEnds[k] count the transitions that fired during decay gap
// k (where the hangover false start lands, and where a false utterance
// closes on its own); finCloses counts the utterances the pre-probe Finalize
// still had to close (a false utterance alive at depth n); nSil/nTalk are
// the probe readings in the 392-r2 shape.
func sweep398v3Trial(t *testing.T, n int, resetEachGap bool) (gapStarts, gapEnds []int, finCloses, nSil, nTalk int) {
	t.Helper()
	d := newTestDetector(t)
	talker := pcmToFloat(synthSpeech(msSamples(1500)))
	gap := pcmToFloat(silence(msSamples(400)))

	if count398Starts(d, talker) == 0 {
		t.Fatalf("the sweep talker did not open an utterance: the probe has no subject")
	}
	_ = d.Finalize(nil)

	// countSeg counts the transitions one whole-window feed produces.
	countSeg := func(samples []float32) (starts, closings int) {
		var eps []Endpoint
		for off := 0; off+WindowSamples <= len(samples); off += WindowSamples {
			eps = d.Accept(samples[off:off+WindowSamples], eps[:0])
			for _, ep := range eps {
				if ep.Kind == EndpointStart {
					starts++
				} else {
					closings++
				}
			}
		}
		return starts, closings
	}

	gapStarts = make([]int, n)
	gapEnds = make([]int, n)
	for k := 0; k < n; k++ {
		if resetEachGap {
			d.Reset()
		}
		gapStarts[k], gapEnds[k] = countSeg(gap)
	}
	if eps := d.Finalize(nil); len(eps) > 0 {
		finCloses = len(eps)
	}
	nSil, _ = countSeg(gap)
	nTalk = count398Starts(d, talker)
	return gapStarts, gapEnds, finCloses, nSil, nTalk
}

// Test398v3HangoverDecaySweep measures the natural decay curve: after the
// utterance closes, from which fed silent gap onward does the hangover stop
// firing. Arm 1 of the residual (c) sweep - no Reset anywhere. The n=0 point
// is the 392-r2/398-r1 control shape (talker -> Finalize -> gap -> talker)
// and reads (silence=1 talker=0); the log flags any other reading as a shape
// mismatch against 398-r1.
//
// This is a measurement, not a pinned assertion: the curve is logged, the
// first fully-clean depth is derived in the log, and nothing about the curve
// fails the test. A flat-zero curve (no firing at any depth) would be
// reported as-is - for the owner it is stronger evidence that skipping the
// cold segment is safe.
func Test398v3HangoverDecaySweep(t *testing.T) {
	decayed := -1
	for n := 0; n <= 16; n++ {
		gapStarts, gapEnds, finCloses, nSil, nTalk := sweep398v3Trial(t, n, false)
		t.Logf("398v3-NATURAL n=%2d decayGapStarts=%v decayGapEnds=%v leftOpen=%d probe (silence=%d talker=%d)",
			n, gapStarts, gapEnds, finCloses, nSil, nTalk)
		if n == 0 && (nSil != 1 || nTalk != 0) {
			t.Logf("398v3-SHAPE-MISMATCH: the n=0 control read (silence=%d talker=%d); "+
				"the 398-r1 control shape is (1 0), so this sweep is not on the 398-r1 probe shape",
				nSil, nTalk)
		}
		if decayed < 0 && nSil == 0 && nTalk == 1 {
			decayed = n
		}
	}
	if decayed < 0 {
		t.Logf("398v3-DECAY-POINT none within 16 gaps: no depth read (silence=0 talker=1)")
	} else {
		t.Logf("398v3-DECAY-POINT n=%d (%dms of fed silence): the first depth whose probe read (silence=0 talker=1)",
			decayed, decayed*400)
	}
}

// Test398v3DecaySweepZeroedControl runs the same sweep with a Reset before
// every gap - the 398-r1 POST-RESET shape repeated per gap, the zeroed
// reference the natural curve is read against. Unlike the natural arm this
// one is guarded: the control must be flat at every depth. Firing on a reset
// gap means the per-gap Reset did not clear the hangover and the natural
// readings are confounded; a silent probe talker means the reset reached the
// classifier, which ticket 398 forbids (stop the line).
func Test398v3DecaySweepZeroedControl(t *testing.T) {
	for n := 1; n <= 16; n++ {
		gapStarts, _, finCloses, nSil, nTalk := sweep398v3Trial(t, n, true)
		t.Logf("398v3-ZEROED n=%2d decayGapStarts=%v leftOpen=%d probe (silence=%d talker=%d)",
			n, gapStarts, finCloses, nSil, nTalk)
		fired := nSil != 0 || finCloses != 0
		for _, c := range gapStarts {
			fired = fired || c != 0
		}
		if fired {
			t.Fatalf("MECHANISM MOVED: the zeroed control fired on silence at n=%d (decayGapStarts=%v leftOpen=%d nSil=%d); "+
				"the per-gap Reset did not clear the hangover, so the natural-decay readings have no zeroed reference",
				n, gapStarts, finCloses, nSil)
		}
		if nTalk != 1 {
			t.Fatalf("MECHANISM MOVED: the zeroed control's probe talker did not open at n=%d (talker=%d); "+
				"talker=0 means the reset reached the classifier, which ticket 398 forbids (stop the line)",
				n, nTalk)
		}
	}
}
