//go:build windows

package main

// Ticket 445 AC#1 probe. CLONE-TREE ONLY: this file is staged into the
// out-of-repo export tree (D:/tmp/wisp/20261114/445-tree/cmd/wisp/
// zz445probe_windows_test.go) and is never committed to the package, so the
// main repo keeps a zero-new-case surface and the probe cannot add a CI row.
//
// THE QUESTION THE TICKET ASKS FIRST, AND WHY IT DECIDES THE TICKET: in the red
// shot (CI run 43575891409, cmd/wisp/resident_voice_interrupt_390r1_windows_test.go:107)
// the last edge of the D43 timeline was still llm.first-token -> Speaking when
// the case judged it. Either the click's TriggerInterrupt never reached the leg
// (a REAL product hole between the ball and the voice leg, and the ticket
// becomes a defect) or it did and the instrument read the walk before it
// landed (an instrument defect). This probe answers that with readings:
//
//	trigger_calls - ball.attachVoiceInterrupt is wrapped in a counter, so the
//	                executor the case itself attaches is counted where the
//	                assembly hands it over. Not a log line, not an inference.
//	redline_hits  - the shipped completion marker `(aborted>=1 && state==Listening)`
//	                is waited for EXACTLY as the shipped case waits for it, then
//	                the last edge is read exactly as the shipped case reads it.
//	                A round counts as a hit when that reading is the pre-#27 edge
//	                - i.e. the shipped assertion would have gone red.
//	lag           - the measured gap between "the walk's state is visible" and
//	                "the walk's edge is on record". This is the window the old
//	                instrument lives in, and it is a property of production code:
//	                statemachine.Dispatch publishes the state inside its own
//	                critical section (machine.go:144) and rv.dispatch appends the
//	                edge only after Dispatch returns (resident_voice_windows.go:1441-1453),
//	                with the row's side effects fired in between (machine.go:162-164).
//
// WHY A SPIN AND NOT THE SHIPPED 5MS POLL: the ticket forbids "ran a hundred
// times and never saw red" as a reproduction. The shipped poll samples the
// window once per 5ms, so the hit rate there is window/5ms - which is exactly
// why CI saw one red instead of a deterministic failure. This probe samples the
// SAME predicate continuously, so the window is measured instead of hoped for,
// and the hit rate it prints is an upper bound on the shipped instrument's
// (the honest reading of the difference is in the evidence file).
//
// Every round also re-reads the round's own hop AFTER it lands and requires the
// pre-#27 edge to be gone: the probe demonstrates the defect AND the shape of
// the fix in the same run.

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

func Test445ProbeClickStopTriggerCalledAndWindow(t *testing.T) {
	rounds := 100
	if v := os.Getenv("WISP_445_ROUNDS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			t.Fatalf("WISP_445_ROUNDS=%q is not a positive integer", v)
		}
		rounds = n
	}

	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "445 probe",
		Tts:          voiceWireTtsFiles(t),
	}
	gate := &voice381r2Gate{}
	var speakerMu sync.Mutex
	var probes []*voice381r2ProbeSpeaker
	var triggerCalls atomic.Int64
	rv, _ := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) {
			p := newVoice381r2ProbeSpeaker(rate)
			speakerMu.Lock()
			probes = append(probes, p)
			speakerMu.Unlock()
			return p, nil
		}
		rv.attachCaptureGate(gate)
	})
	// The hop under test, wrapped: same executor the shipped case attaches, with
	// a counter in front of it so "was TriggerInterrupt called" is a reading.
	ball.attachVoiceInterrupt(func() {
		triggerCalls.Add(1)
		rv.TriggerInterrupt()
	})

	// A short sentence: the click lands at the parking point, so the sentence's
	// own length is not part of any measurement here.
	const sentence = "第一句话。"
	hits := 0
	fastHits, pollHits := 0, 0
	fastRounds, pollRounds := 0, 0
	var lags []time.Duration
	redlines := make([]string, 0, 4)

	for round := 1; round <= rounds; round++ {
		roundStart := time.Now()
		walkToThinking381r2(t, rv)
		abortedBefore := rv.ttsAborted.Load()
		spokenBefore := rv.ttsSpoken.Load()
		if !rv.announce(sentence) {
			t.Fatalf("round %d: the announce path refused the sentence; problems=%v", round, rv.problemsSnapshot())
		}
		voice15eWaitFor(t, fmt.Sprintf("round %d: the owner to park in Drain", round), func() bool {
			speakerMu.Lock()
			defer speakerMu.Unlock()
			if len(probes) < round {
				return false
			}
			select {
			case <-probes[round-1].parked:
				return true
			default:
				return false
			}
		})
		// The baseline for "the click appended exactly one edge" is taken AFTER
		// the announce's own #15 edge has landed: the parked wait above is one
		// statement after synthesis, so the announce edge is on record here.
		hopsBefore := len(rv.timelineSnapshot())
		ball.clickInterruptGesture()

		// ---- the shipped predicate, at two sampling rates (alternating by round).
		// The window is a property of the production publication order; the hit
		// rate is a property of how often the reader looks. Running both in one
		// pass turns "is this a flake" into a pair of numbers:
		//   fast    - a spin: the window is sampled continuously, so the hit rate
		//             is (close to) the window's own probability of being caught;
		//   shipped - exactly the shipped 5ms poll (voice15eWaitFor), i.e. the
		//             instrument CI ran.
		fast := round%2 == 1
		if fast {
			fastRounds++
		} else {
			pollRounds++
		}
		if fast {
			deadline := time.Now().Add(20 * time.Second)
			for {
				if rv.ttsAborted.Load() > abortedBefore && rv.machine.State() == statemachine.StateListening {
					break
				}
				if rv.ttsSpoken.Load() > spokenBefore {
					t.Fatalf("round %d: the sentence played out instead of stopping (spoken=%d aborted=%d trigger_calls=%d): "+
						"the click's executor did not stop the playback", round, rv.ttsSpoken.Load(), rv.ttsAborted.Load(), triggerCalls.Load())
				}
				if time.Now().After(deadline) {
					t.Fatalf("round %d: timed out 20s after the click waiting for (aborted>%d && state==Listening): aborted=%d state=%s "+
						"trigger_calls=%d timeline=%+v", round, abortedBefore, rv.ttsAborted.Load(),
						string(rv.machine.State()), triggerCalls.Load(), rv.timelineSnapshot())
				}
			}
		} else {
			voice15eWaitFor(t, fmt.Sprintf("round %d: the click to end the announcement (aborted or played out)", round), func() bool {
				return (rv.ttsAborted.Load() > abortedBefore && rv.machine.State() == statemachine.StateListening) ||
					rv.ttsSpoken.Load() > spokenBefore
			})
			if rv.ttsSpoken.Load() > spokenBefore {
				t.Fatalf("round %d: the sentence played out instead of stopping (spoken=%d aborted=%d trigger_calls=%d)",
					round, rv.ttsSpoken.Load(), rv.ttsAborted.Load(), triggerCalls.Load())
			}
		}
		flipAt := time.Now()
		naive := rv.timelineSnapshot()
		last := naive[len(naive)-1]
		if last.Event != statemachine.EvInterrupt || last.To != statemachine.StateListening {
			hits++
			mode := "shipped"
			if fast {
				mode = "fast"
				fastHits++
			} else {
				pollHits++
			}
			if len(redlines) < 4 {
				redlines = append(redlines, fmt.Sprintf(
					"round %d (%s sampling): the click's stop walked %+v, want #27 voice.interrupt -> Listening "+
						"(trigger_calls=%d aborted=%d state=%s edges=%d)",
					round, mode, last, triggerCalls.Load(), rv.ttsAborted.Load(), string(rv.machine.State()), len(naive)))
			}
		}
		// ---- the window: how long after the state became visible did the edge land.
		hopDeadline := time.Now().Add(20 * time.Second)
		for len(rv.timelineSnapshot()) <= hopsBefore {
			if time.Now().After(hopDeadline) {
				t.Fatalf("round %d: the #27 edge never reached the timeline 20s after the state flipped: state=%s aborted=%d edges=%+v",
					round, string(rv.machine.State()), rv.ttsAborted.Load(), rv.timelineSnapshot())
			}
		}
		lags = append(lags, time.Since(flipAt))

		// ---- the round's own readings, all from the same shot as the hit above.
		landed := rv.timelineSnapshot()
		got := landed[len(landed)-1]
		if got.Event != statemachine.EvInterrupt || got.To != statemachine.StateListening {
			t.Fatalf("round %d: the walk that landed was %+v, want #27 voice.interrupt -> Listening", round, got)
		}
		if n := len(landed) - hopsBefore; n != 1 {
			t.Fatalf("round %d: the click appended %d edge(s), want exactly 1: %+v", round, n, landed[hopsBefore:])
		}
		if c := triggerCalls.Load(); c != int64(round) {
			t.Fatalf("round %d: the click's executor ran %d time(s) in total, want %d (one per round)", round, c, round)
		}
		if a := rv.ttsAborted.Load(); a != uint64(round) {
			t.Fatalf("round %d: aborted=%d, want %d", round, a, round)
		}
		if s := rv.ttsSpoken.Load(); s != 0 {
			t.Fatalf("round %d: an interrupted sentence was counted as spoken (spoken=%d)", round, s)
		}
		if p := rv.problemsSnapshot(); len(p) != 0 {
			t.Fatalf("round %d: a click interrupt left problems behind: %v", round, p)
		}
		speakerMu.Lock()
		aborts := probes[round-1].aborts
		speakerMu.Unlock()
		if aborts != 1 {
			t.Fatalf("round %d: the abort reached the round's speaker %d time(s), want 1", round, aborts)
		}
		if round <= 3 || round%20 == 0 {
			t.Logf("445-ROUND %d took=%s lag=%s", round, time.Since(roundStart).Round(time.Millisecond),
				lags[len(lags)-1].Round(time.Nanosecond))
		}
	}

	sort.Slice(lags, func(i, j int) bool { return lags[i] < lags[j] })
	median := lags[len(lags)/2]
	for _, line := range redlines {
		t.Logf("445-REDLINE %s", line)
	}
	t.Logf("445-WINDOW rounds=%d redline_hits=%d hit_rate=%.1f%% fast=%d/%d shipped_5ms=%d/%d trigger_calls=%d aborted=%d spoken=0 settle_min=%s settle_median=%s settle_max=%s",
		rounds, hits, 100*float64(hits)/float64(rounds), fastHits, fastRounds, pollHits, pollRounds,
		triggerCalls.Load(), rv.ttsAborted.Load(),
		lags[0].Round(time.Nanosecond), median.Round(time.Nanosecond), lags[len(lags)-1].Round(time.Nanosecond))
}
