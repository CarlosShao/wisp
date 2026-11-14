//go:build windows

package main

// Ticket 390-r1's readings: the ball click - the interrupt channel's second
// trigger source, named in the channel's own comment since 381 and left
// unhosted by 381-r1 - now fires the same production TriggerInterrupt the Esc
// path fires. Four structures, all staged the way 381-r2 staged theirs (the
// DLLs beside the test binary, the WISP_* models; the live case skips where
// there is no output device):
//
//	Play   during playback a click goes through the ball host's gesture entry
//	       (clickInterruptGesture, the one OnClickBall runs) and stops the
//	       announcement: abort in flight, D43#27 walk, gate reopened, next
//	       announcement intact. The trigger line is the ATTACHED HOST's, not a
//	       bare rv.TriggerInterrupt - what is under test is the hop the
//	       assembly builds (attachVoiceInterrupt), not the hop 381-r2 already
//	       proved.
//	Device the live reading: one control round plus five rounds of click-
//	       during-playback against the real waveOut device and the real matcha
//	       synthesis; every round's unstick and reset must land within 400ms.
//	Idle   the positive control: a click while nothing speaks dispatches
//	       nothing, moves no state, pollutes no list (381-r2's idle Esc shape,
//	       now through the ball gesture), and a ball without an attached
//	       executor records the gesture as unhosted without touching the leg.
//	Unhook the unattached shape: with no executor attached, a click during
//	       playback leaves the announcement alone - it plays out. This is what
//	       makes the mutation reading meaningful: the ONLY delta between Play
//	       and Unhook is the attachVoiceInterrupt hop.
//
// Ticket 445's repair of this file (CI run 43575891409, the red at old line 107):
// the Play case judged "the edge that landed last" at the instant its completion
// marker said the walk was over, and the marker is satisfied one edge early - the
// state is published by statemachine.Dispatch while rv.dispatch appends the edge
// only after Dispatch returns. The judgement now waits for the walk to land first
// (waitHopAfter445), so the same assertion reads the edge it is about. The
// instrument's tooth is unchanged and re-measured in the export tree: with the
// click's executor really not called the case still reds, naming the MUTANT line;
// with the timeline's append really removed the wait itself reds with its numbers.

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio/aecspike"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Test390r1BallClickPlayingInterruptAbortsAndWalks27 is the staged half of
// AC#1+AC#2: the owner is parked inside Drain, a BALL CLICK (through the ball
// host's gesture entry) stops the probe device, the walk is D43#27, the gate
// reopens, and the next announcement still speaks.
func Test390r1BallClickPlayingInterruptAbortsAndWalks27(t *testing.T) {
	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "390-r1 test",
		Tts:          voiceWireTtsFiles(t),
	}
	gate := &voice381r2Gate{}
	var speakerMu sync.Mutex
	var probes []*voice381r2ProbeSpeaker
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
	// The hop under test: the ball host's click executor is the voice leg's
	// own TriggerInterrupt, attached after the leg exists - the same order the
	// assembly root's startResidentVoice uses.
	ball.attachVoiceInterrupt(rv.TriggerInterrupt)

	walkToThinking381r2(t, rv)
	if !rv.announce("这是第一句话，播报需要一点时间。") {
		t.Fatalf("the announce path refused the sentence; problems=%v", rv.problemsSnapshot())
	}
	voice15eWaitFor(t, "the owner to park in Drain", func() bool {
		speakerMu.Lock()
		defer speakerMu.Unlock()
		return len(probes) > 0 && func() bool {
			select {
			case <-probes[0].parked:
				return true
			default:
				return false
			}
		}()
	})
	// The baseline for "the click appended exactly one edge" is taken before the
	// gesture: the announce's own #15 edge is already on record by now.
	hopsBefore := len(rv.timelineSnapshot())
	ball.clickInterruptGesture()
	voice15eWaitFor(t, "the click to end the announcement (aborted or played out)", func() bool {
		return (rv.ttsAborted.Load() >= 1 && rv.machine.State() == statemachine.StateListening) ||
			rv.ttsSpoken.Load() >= 1
	})
	if got := rv.ttsAborted.Load(); got != 1 {
		t.Fatalf("MUTANT 390r1-M1 SHAPE: the ball click did not stop the announcement in flight "+
			"(spoken=%d, aborted=%d): the click gesture did not fire the attached TriggerInterrupt, so the "+
			"sentence played out - the second trigger source is not wired (ticket 390: 球单击→TriggerInterrupt 真挂)",
			rv.ttsSpoken.Load(), got)
	}
	if got := rv.ttsSpoken.Load(); got != 0 {
		t.Fatalf("an interrupted announcement was counted as spoken (%d)", got)
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("a click interrupt left problems behind: %v", got)
	}
	// TICKET 445 AC#2: wait for the walk to LAND, then judge it. The marker above
	// is satisfied by (aborted>=1 && state==Listening), and that state is visible
	// BEFORE this leg appends the edge - statemachine.Dispatch publishes the state
	// inside its own critical section and rv.dispatch appends after Dispatch
	// returns - so a last-edge read at that instant can see the #15 edge. CI run
	// 43575891409 read exactly there. The wait is bounded and named, and the
	// assertion below is unchanged: #27 must be the edge that landed LAST.
	hops := waitHopAfter445(t, rv, hopsBefore, "the click's stop (D43#27 voice.interrupt Speaking -> Listening)")
	last := hops[len(hops)-1]
	if last.Event != statemachine.EvInterrupt || last.To != statemachine.StateListening {
		t.Fatalf("the click's stop walked %+v, want #27 voice.interrupt -> Listening", last)
	}
	if n := len(hops) - hopsBefore; n != 1 {
		t.Fatalf("the click appended %d D43 edge(s), want exactly 1: %+v", n, hops[hopsBefore:])
	}
	voice15eWaitFor(t, "the capture gate to reopen after the click's stop", func() bool {
		closed, reopened := gate.count()
		return closed == 1 && reopened == 1
	})
	speakerMu.Lock()
	aborts1 := probes[0].aborts
	speakerMu.Unlock()
	if aborts1 != 1 {
		t.Fatalf("the click's stop fired %d times on the speaker, want 1", aborts1)
	}

	// The next announcement enters, speaks, and its gate releases - the token
	// the click left is speak()'s entry-drain business, not this sentence's.
	walkToThinking381r2(t, rv)
	if !rv.announce("这是第二句话。") {
		t.Fatalf("the announce path refused the next sentence; problems=%v", rv.problemsSnapshot())
	}
	voice15eWaitFor(t, "the next announcement to be spoken and its gate released", func() bool {
		closed, reopened := gate.count()
		return rv.ttsSpoken.Load() == 1 && closed == 2 && reopened == 2
	})
	speakerMu.Lock()
	if len(probes) < 2 {
		speakerMu.Unlock()
		t.Fatalf("the next announcement opened no speaker of its own")
	}
	aborts2 := probes[1].aborts
	speakerMu.Unlock()
	if aborts2 != 0 {
		t.Fatalf("a stale abort reached the next announcement's speaker (%d aborts)", aborts2)
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("the next announcement left problems behind: %v", got)
	}
	closed, reopened := gate.count()
	t.Logf("CLICK-27 aborted=%d spoken=%d next_spoken=%d probe1_aborts=%d probe2_aborts=%d gate_closed=%d gate_reopened=%d",
		rv.ttsAborted.Load(), 0, rv.ttsSpoken.Load(), aborts1, aborts2, closed, reopened)
}

// Test390r1BallClickRealDeviceWithin400ms is the live reading AC#2 asks for:
// one control round plus five rounds of click-during-playback, each against the
// REAL waveOut device and the REAL matcha synthesis, the trigger delivered
// through the ball host's gesture entry. Skipped where there is no output
// device (a CI box has none).
func Test390r1BallClickRealDeviceWithin400ms(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	// Ticket 388's前置: the sentence below is real matcha speech on the real
	// device, so the process render session is attenuated BEFORE the first
	// announce (control round included) and the owner's original volume is
	// restored in t.Cleanup. The cgo float path only (59-r1's finding #2);
	// the waveOut stream itself is untouched, so the measured abort-latency
	// path keeps its exact semantics - only loudness changes.
	attenuateAnnounce388(t, aecspike.DefaultAnnounceAttenuatedVolume)
	const sentence = "任务已经完成，明天下午三点的会议改到四点，请准时参加，谢谢。"

	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "390-r1 test",
		Tts:          voiceWireTtsFiles(t),
	}
	gate := &voice381r2Gate{}
	var speakerMu sync.Mutex
	var spies []*voice381r2DeviceSpy
	rv, _ := voiceWireStart(t, plan, ball, sink, func(rv *residentVoice) {
		rv.newSpeaker = func(rate int) (voiceSpeaker, error) {
			sp, err := openWaveOutSpeaker(rate)
			if err != nil {
				return nil, err
			}
			spy := &voice381r2DeviceSpy{voiceSpeaker: sp, parked: make(chan struct{}), abortDone: make(chan struct{})}
			speakerMu.Lock()
			spies = append(spies, spy)
			speakerMu.Unlock()
			return spy, nil
		}
		rv.attachCaptureGate(gate)
	})
	ball.attachVoiceInterrupt(rv.TriggerInterrupt)

	spyAt := func(i int) *voice381r2DeviceSpy {
		voice15eWaitFor(t, fmt.Sprintf("speaker %d to be opened", i), func() bool {
			speakerMu.Lock()
			defer speakerMu.Unlock()
			return len(spies) > i
		})
		speakerMu.Lock()
		defer speakerMu.Unlock()
		return spies[i]
	}

	// Control round: the same sentence with no click, so "played out" has a
	// number on this machine.
	controlStart := time.Now()
	walkToThinking381r2(t, rv)
	if !rv.announce(sentence) {
		t.Fatalf("the announce path refused the sentence; problems=%v", rv.problemsSnapshot())
	}
	control := spyAt(0)
	<-control.parked
	parkedAt := time.Now()
	voice15eWaitFor(t, "the control sentence to play out (spoken and back in Warm)", func() bool {
		return rv.ttsSpoken.Load() >= 1 && rv.machine.State() == statemachine.StateWarm
	})
	t.Logf("CLICK-CONTROL parked_to_done=%s announce_to_done=%s",
		time.Since(parkedAt).Round(time.Millisecond), time.Since(controlStart).Round(time.Millisecond))

	// Five record rounds: the click lands 300ms into the playback; the reading
	// is the distance from the gesture to the device actually stopped.
	const rounds = 5
	var unsticks388, resets388 []time.Duration
	for r := 1; r <= rounds; r++ {
		walkToThinking381r2(t, rv)
		if !rv.announce(sentence) {
			t.Fatalf("round %d: the announce path refused the sentence; problems=%v", r, rv.problemsSnapshot())
		}
		spy := spyAt(r)
		<-spy.parked
		spokenBefore := rv.ttsSpoken.Load()
		time.Sleep(300 * time.Millisecond) // the click lands 300ms into the playback
		t0 := time.Now()
		ball.clickInterruptGesture()
		// The completion marker is the abort branch FULLY done - the counter
		// AND the #27 walk (381-r2's lesson: the counter alone can race the
		// state read, it is incremented one statement before the dispatch).
		voice15eWaitFor(t, fmt.Sprintf("round %d: the stop (aborted or played out)", r), func() bool {
			return (rv.ttsAborted.Load() >= uint64(r) && rv.machine.State() == statemachine.StateListening) ||
				rv.ttsSpoken.Load() > spokenBefore
		})
		if got := rv.ttsAborted.Load(); got < uint64(r) {
			t.Fatalf("MUTANT 390r1-M1 SHAPE: round %d played out instead of stopping (spoken=%d, aborted=%d): "+
				"the ball click did not fire the attached TriggerInterrupt, so the sentence played out",
				r, rv.ttsSpoken.Load(), got)
		}
		unstick := time.Since(t0)
		if unstick > 400*time.Millisecond {
			t.Fatalf("round %d: the device answered the click in %s, want within 400ms", r, unstick)
		}
		<-spy.abortDone
		reset := time.Since(t0)
		if reset > 400*time.Millisecond {
			t.Fatalf("round %d: waveOutReset landed %s after the click, want within 400ms", r, reset)
		}
		unsticks388 = append(unsticks388, unstick)
		resets388 = append(resets388, reset)
		if got := rv.ttsSpoken.Load(); got != spokenBefore {
			t.Fatalf("round %d: the interrupted sentence was counted as spoken (%d -> %d)", r, spokenBefore, got)
		}
		// Abort's return value stays a READING, not an assert (381-r2's named
		// race: the owner's Close and the async Abort both land a reset; either
		// way the audio stops).
		t.Logf("CLICK-ROUND-%d unstick=%s reset=%s abort_returned=%d",
			r, unstick.Round(time.Millisecond), reset.Round(time.Millisecond), spy.dropped)
	}
	// Ticket 388 AC#4's distribution reading: the attenuation rides the
	// session volume, not the waveOut stream, so the stop-latency
	// distribution must be unchanged - every sample inside the same 400ms
	// bound the per-round asserts above already enforce.
	assertInterruptDistribution388(t, "unstick", unsticks388)
	assertInterruptDistribution388(t, "waveOutReset", resets388)
	voice15eWaitFor(t, "every round's gate release to land", func() bool {
		closed, reopened := gate.count()
		return closed == rounds+1 && reopened == rounds+1
	})
	if got := rv.ttsSpoken.Load(); got != 1 {
		t.Fatalf("spoken=%d after the control and %d interrupted rounds, want 1", got, rounds)
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("the live rounds left problems behind: %v", got)
	}
}

// Test390r1IdleBallClickIsInert is the positive control: a click while nothing
// speaks dispatches nothing, moves no state, records no problem, touches no
// counter - the same inert shape 381-r2 proved for the idle Esc, now through
// the ball gesture. A nil ball and a ball without an attached executor are the
// nil-safety half: neither may panic or reach a leg.
func Test390r1IdleBallClickIsInert(t *testing.T) {
	ball := &residentBall{}
	rv := newResidentVoice(make(chan []byte), voicePlan{}, ball, nil)
	ball.attachVoiceInterrupt(rv.TriggerInterrupt)
	ball.clickInterruptGesture()
	if got := rv.machine.State(); got != statemachine.StateSleeping {
		t.Fatalf("an idle click moved the machine to %s", string(got))
	}
	if len(rv.timelineSnapshot()) != 0 {
		t.Fatalf("an idle click walked %v", rv.timelineSnapshot())
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("an idle click polluted the problem list: %v", got)
	}
	if rv.ttsAborted.Load() != 0 || rv.ttsSpoken.Load() != 0 {
		t.Fatalf("an idle click touched the announce counters: aborted=%d spoken=%d",
			rv.ttsAborted.Load(), rv.ttsSpoken.Load())
	}
	// nil-safety: the gesture on a nil host and on a host without an executor.
	var nilBall *residentBall
	nilBall.clickInterruptGesture()
	unhosted := &residentBall{}
	unhosted.clickInterruptGesture()
}

// Test390r1UnattachedClickDoesNotInterrupt is the unattached shape the
// mutation reading stands on: with NO executor attached to the ball host, a
// click during playback records the gesture as unhosted and the announcement
// plays out. The only delta to Test390r1BallClickPlayingInterruptAbortsAndWalks27
// is the attachVoiceInterrupt call - which is why the mutation (removing the
// click's trigger hop in clickInterruptGesture) turns that test red and leaves
// this one green.
func Test390r1UnattachedClickDoesNotInterrupt(t *testing.T) {
	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "390-r1 test",
		Tts:          voiceWireTtsFiles(t),
	}
	gate := &voice381r2Gate{}
	var speakerMu sync.Mutex
	var probes []*voice381r2ProbeSpeaker
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
	// Deliberately NO ball.attachVoiceInterrupt here: the pre-390 shape.

	walkToThinking381r2(t, rv)
	if !rv.announce("这是第二句话。") {
		t.Fatalf("the announce path refused the sentence; problems=%v", rv.problemsSnapshot())
	}
	voice15eWaitFor(t, "the owner to park in Drain", func() bool {
		speakerMu.Lock()
		defer speakerMu.Unlock()
		return len(probes) > 0 && func() bool {
			select {
			case <-probes[0].parked:
				return true
			default:
				return false
			}
		}()
	})
	// A click arrives with no executor behind it: nothing may reach the
	// interrupt channel, so the announcement plays out.
	ball.clickInterruptGesture()
	voice15eWaitFor(t, "the announcement to play out untouched", func() bool {
		return rv.ttsSpoken.Load() >= 1
	})
	if got := rv.ttsAborted.Load(); got != 0 {
		t.Fatalf("an unattached click stopped the announcement (aborted=%d): the gesture reached the "+
			"interrupt channel without an executor", got)
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("the unattached click left problems behind: %v", got)
	}
	t.Logf("UNATTACHED spoken=%d aborted=0 problems=0 (the pre-390 shape, kept honest)",
		rv.ttsSpoken.Load())
}

// waitHopAfter445 is ticket 445 AC#2's shape: wait for the walk to LAND, then
// judge what landed. It returns the leg's whole D43 timeline once it has grown
// past `from`, and fails with a named line and numbers if it never does.
//
// WHY A WAIT AND NOT A LONGER SLEEP: the state and the edge are published in two
// steps. statemachine.Dispatch sets the new state inside its own critical
// section (internal/statemachine/machine.go:144) and returns only after firing
// the row's side effects (machine.go:162-164); cmd/wisp's rv.dispatch appends
// the edge after that call returns (resident_voice_windows.go:1441-1453). Every
// completion marker in this family waits for (a counter, machine.State()), so it
// can be satisfied one edge early - which is exactly what CI run 43575891409
// read at this file's old line 107 (the #15 edge instead of #27).
//
// The probe under .scratch/wisp/probes/445 measures the window on this host:
// with a continuous sampler the pre-#27 read is 3/20 shots, at the shipped 5ms
// poll it is 0/20 - the defect is real and its visual rate is a sampling-rate
// artefact. So the fix cannot be "wait longer": the wait's predicate is the
// assertion itself (an edge past the baseline exists), and the assertion still
// says WHICH edge must have landed last.
func waitHopAfter445(t *testing.T, rv *residentVoice, from int, what string) []voiceTransition {
	t.Helper()
	const bound = 5 * time.Second
	start := time.Now()
	for {
		if hops := rv.timelineSnapshot(); len(hops) > from {
			return hops
		}
		if time.Since(start) > bound {
			hops := rv.timelineSnapshot()
			t.Fatalf("the D43 walk never landed within %s: %s: timeline has %d edge(s) (baseline %d), state=%s, "+
				"aborted=%d spoken=%d, last=%s", bound, what, len(hops), from,
				string(rv.machine.State()), rv.ttsAborted.Load(), rv.ttsSpoken.Load(), voiceLastHop445(hops))
		}
		time.Sleep(time.Millisecond)
	}
}

// waitWalked445 is the same discipline for the cases whose assertion is "this row
// is on record" rather than "this row is last": it waits, bounded, for the named
// edge, so a reading taken at a counter's flip cannot judge a walk the announce
// has not appended yet. A row that never lands fails here with the numbers -
// never passes quietly.
func waitWalked445(t *testing.T, rv *residentVoice, from statemachine.State, ev statemachine.Event, to statemachine.State, what string) {
	t.Helper()
	const bound = 5 * time.Second
	start := time.Now()
	for {
		for _, tr := range rv.timelineSnapshot() {
			if tr.From == from && tr.Event == ev && tr.To == to {
				return
			}
		}
		if time.Since(start) > bound {
			t.Fatalf("the D43 walk never landed within %s: %s (%s + %s -> %s): state=%s, aborted=%d spoken=%d, timeline=%+v",
				bound, what, string(from), string(ev), string(to), string(rv.machine.State()),
				rv.ttsAborted.Load(), rv.ttsSpoken.Load(), rv.timelineSnapshot())
		}
		time.Sleep(time.Millisecond)
	}
}

// voiceLastHop445 is the "what did we actually read" half of a failure line: an
// empty timeline is named instead of indexed.
func voiceLastHop445(hops []voiceTransition) string {
	if len(hops) == 0 {
		return "<none>"
	}
	return fmt.Sprintf("%+v", hops[len(hops)-1])
}
