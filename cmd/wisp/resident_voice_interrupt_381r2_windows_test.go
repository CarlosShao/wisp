//go:build windows

package main

// Ticket 381-r2's readings: the interrupt channel 381-r1 reserved becomes a
// real stop. Three structures, each with its named mutant (the mutant shapes
// are what the failure messages spell out, 26r3's convention):
//
//	M1  TriggerInterrupt aborts the speaker in flight - during playback the
//	    owner goroutine is parked inside Drain where the chunk-boundary select
//	    never runs, so the abort is the only half that can stop audible audio
//	    (the 6,653ms play-out 381-v2 measured dies here).
//	M2  speak() drains a stale interrupt token at the entry - a token left by
//	    an idle Esc must not silently kill the next announcement (381-v2 §5).
//	M3  the 381-r1 veto-word substring branch is gone - benign sentences with
//	    停/取消/别 substrings reach the agent (381-v2 §4 measured 12/12 hits).
//
// Plus the live-device case, which is where the ≤400ms stop reading comes
// from: five rounds of Esc-during-playback against the real waveOut device and
// the real matcha synthesis, through the production TriggerInterrupt path.
//
// Staging is 15e's: the sherpa DLLs beside the test binary, WISP_* for the
// models; the live-device case skips where there is no output device.

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// vetoCorpus381r2 is this leg's benign corpus: twelve sentences that contain
// 停, 取消 or 别 as a substring and are not stop commands (381-v2 §4's shape,
// re-collected here because the original probe corpus was scratch-side). Under
// the 381-r1 branch every one of them was swallowed with a D43 rejection;
// after 381-r2's deletion every one of them reaches the agent.
var vetoCorpus381r2 = []string{
	"导航去最近的停车场",
	"特别提醒我明早八点起床",
	"查一下那趟高铁取消没有",
	"别忘了帮我带伞",
	"附近哪里可以停汽车",
	"这个快递能不能取消然后再下单",
	"提醒我别错过末班车",
	"帮我查一下停车场收费多少",
	"那个会议取消了吗",
	"特别的日子要提前准备",
	"我的会员能不能取消自动续费",
	"把车停到地库再去上班",
}

// voice381r2Gate is the capture gate's recorder: it keeps the on/off sequence
// so "the gate was closed for the announcement and reopened on the way out"
// is a reading for the interrupted path too.
type voice381r2Gate struct {
	mu sync.Mutex
	on []bool
}

func (g *voice381r2Gate) SetSpeaking(speaking bool) {
	g.mu.Lock()
	g.on = append(g.on, speaking)
	g.mu.Unlock()
}

func (g *voice381r2Gate) count() (closed, reopened int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, speaking := range g.on {
		if speaking {
			closed++
		} else {
			reopened++
		}
	}
	return closed, reopened
}

// voice381r2ProbeSpeaker is the render seam with a device that PARKS: Drain
// holds the caller until the queued samples have "played" on the audio clock
// or an abort lands, whichever comes first - the shape the winmm speaker has
// during playback and the shape the chunk-boundary select cannot reach. The
// abort answers errSpeakerAborted, the same sentinel the real speaker answers;
// that contract is what speak()'s abort branch reads.
type voice381r2ProbeSpeaker struct {
	mu         sync.Mutex
	rate       int
	chunks     int
	samples    int
	aborts     int
	parked     chan struct{}
	parkedOnce sync.Once
	aborted    chan struct{}
	abortOnce  sync.Once
}

func newVoice381r2ProbeSpeaker(rate int) *voice381r2ProbeSpeaker {
	return &voice381r2ProbeSpeaker{rate: rate, parked: make(chan struct{}), aborted: make(chan struct{})}
}

func (s *voice381r2ProbeSpeaker) SampleRate() int { return s.rate }

func (s *voice381r2ProbeSpeaker) Play(samples []float32, engineRate int) error {
	if len(samples) == 0 {
		return nil
	}
	s.mu.Lock()
	s.chunks++
	s.samples += len(samples)
	s.mu.Unlock()
	return nil
}

// Drain parks the caller the way the device parks in playback. The wait is
// bounded by the queued samples on the AUDIO clock (samples to duration, the
// same arithmetic the winmm speaker uses), never by a wall-clock measurement;
// an abort ends it with errSpeakerAborted.
func (s *voice381r2ProbeSpeaker) Drain() error {
	s.mu.Lock()
	playback := time.Duration(s.samples) * time.Second / time.Duration(max(s.rate, 1))
	s.mu.Unlock()
	s.parkedOnce.Do(func() { close(s.parked) })
	if playback <= 0 {
		return nil
	}
	timer := time.NewTimer(playback)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-s.aborted:
		return errSpeakerAborted
	}
}

func (s *voice381r2ProbeSpeaker) Abort() int {
	s.mu.Lock()
	s.aborts++
	s.mu.Unlock()
	s.abortOnce.Do(func() { close(s.aborted) })
	return 1
}

func (s *voice381r2ProbeSpeaker) Close() error { return nil }

// voice381r2DeviceSpy wraps the real waveOut speaker: it reports when the
// owner parked in Drain and what Abort dropped, without changing the device's
// behavior (every call delegates).
type voice381r2DeviceSpy struct {
	voiceSpeaker
	parked     chan struct{}
	parkedOnce sync.Once
	abortDone  chan struct{}
	abortOnce  sync.Once
	dropped    int
}

func (s *voice381r2DeviceSpy) Drain() error {
	s.parkedOnce.Do(func() { close(s.parked) })
	return s.voiceSpeaker.Drain()
}

func (s *voice381r2DeviceSpy) Abort() int {
	n := s.voiceSpeaker.Abort()
	s.abortOnce.Do(func() {
		s.dropped = n
		close(s.abortDone)
	})
	return n
}

// walkToThinking381r2 moves the machine to Thinking through the table's own
// rows, so a test-side announce (which dispatches #15 from Thinking) is legal
// without feeding a real utterance: #4/#29/#33 into Listening, then #11.
func walkToThinking381r2(t *testing.T, rv *residentVoice) {
	t.Helper()
	switch rv.machine.State() {
	case statemachine.StateThinking:
		return
	case statemachine.StateSleeping:
		rv.dispatch(statemachine.EvSummon, nil, "test: D43#4 Sleeping 到 Listening")
	case statemachine.StateWarm:
		rv.dispatch(statemachine.EvSummon, nil, "test: D43#29 Warm 到 Listening")
	case statemachine.StateSettling:
		rv.dispatch(statemachine.EvSummon, nil, "test: D43#33 Settling 回 Listening")
	}
	if got := rv.machine.State(); got != statemachine.StateListening {
		t.Fatalf("state = %s, want Listening before the #11 walk", string(got))
	}
	rv.dispatch(statemachine.EvVadStop, &statemachine.Facts{SpeechMS: 400}, "test: D43#11 Listening 到 Thinking")
	if got := rv.machine.State(); got != statemachine.StateThinking {
		t.Fatalf("state = %s, want Thinking after the #11 walk", string(got))
	}
}

// Test381r2IdleInterruptIsInert is the positive control for the non-announce
// half of the interrupt channel: an Esc while nothing speaks buffers a token,
// aborts nothing (no speaker in flight), dispatches nothing, records no
// problem, and moves no state. The token it leaves is speak()'s entry-drain
// business, not the idle machine's.
func Test381r2IdleInterruptIsInert(t *testing.T) {
	rv := newResidentVoice(make(chan []byte), voicePlan{}, &residentBall{}, nil)
	rv.TriggerInterrupt()
	if got := rv.machine.State(); got != statemachine.StateSleeping {
		t.Fatalf("an idle interrupt moved the machine to %s", string(got))
	}
	if len(rv.timelineSnapshot()) != 0 {
		t.Fatalf("an idle interrupt walked %v", rv.timelineSnapshot())
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("an idle interrupt polluted the problem list: %v", got)
	}
	if rv.ttsAborted.Load() != 0 || rv.ttsSpoken.Load() != 0 {
		t.Fatalf("an idle interrupt touched the announce counters: aborted=%d spoken=%d",
			rv.ttsAborted.Load(), rv.ttsSpoken.Load())
	}
	// nil-safety: the trigger on a nil leg and on a leg without the channel.
	var nilVoice *residentVoice
	nilVoice.TriggerInterrupt()
	noChannel := newResidentVoice(make(chan []byte), voicePlan{}, nil, nil)
	noChannel.interruptCh = nil
	noChannel.TriggerInterrupt()
}

// Test381r2StaleTokenDoesNotKillTheNextAnnouncement is M2's surface: a token
// left by an earlier idle Esc sits in interruptCh, and the NEXT announcement
// must still be spoken - speak() drains the stale token once at the entry.
// The announce runs synchronously on this goroutine (the same call the owner
// makes), so after queue.run() the sentence is fully spoken and the counters
// are final - no waiting, no flake.
func Test381r2StaleTokenDoesNotKillTheNextAnnouncement(t *testing.T) {
	speaker := &voiceWireSpeaker{rate: 22050}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voiceVadNone(),
		Source:       "381-r2 test",
		Tts:          voiceWireTtsFiles(t),
	}
	frames := make(chan []byte)
	rv := newResidentVoice(frames, plan, ball, nil)
	rv.residency = newVoiceResidency(plan, rv.queue, nil)
	rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return speaker, nil }
	t.Cleanup(rv.stop)

	// The stale token: an Esc while the machine was idle.
	rv.TriggerInterrupt()

	// The legal walk into Thinking so the announce's #15 hop is real.
	walkToThinking381r2(t, rv)

	if !rv.announce("任务已完成，明天下午三点的会改到四点。") {
		t.Fatalf("the announce path refused the sentence; problems=%v", rv.problemsSnapshot())
	}
	rv.queue.run()

	if got := rv.ttsSpoken.Load(); got != 1 {
		t.Fatalf("MUTANT 381r2-M2 SHAPE: the announcement was killed by a stale interrupt token "+
			"(spoken=%d, chunks=%d): speak() does not drain interruptCh at the entry, so a token left by "+
			"an earlier idle Esc silently kills the next sentence (381-v2 §5: spoken=0, no problem row)",
			got, speaker.chunks)
	}
	if got := rv.ttsAborted.Load(); got != 0 {
		t.Fatalf("the stale token was counted as a mid-playback abort (%d)", got)
	}
	if speaker.chunks == 0 {
		t.Fatalf("the announcement emitted no chunks; problems=%v", rv.problemsSnapshot())
	}
	for _, tr := range rv.timelineSnapshot() {
		if tr.Event == statemachine.EvInterrupt {
			t.Fatalf("the stale token walked an interrupt edge: %+v", tr)
		}
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("the stale-token announce left problems behind: %v", got)
	}
	t.Logf("STALE-TOKEN spoken=%d chunks=%d samples=%d timeline=%d problems=0",
		rv.ttsSpoken.Load(), speaker.chunks, speaker.written(), len(rv.timelineSnapshot()))
}

// Test381r2PlayingPhaseInterruptAbortsAndWalks27 is M1's surface: the
// interrupt lands while the owner is parked inside Drain (the playing phase),
// the abort stops the probe device, the walk is D43#27 (not #26), the sentence
// is never counted as spoken, the capture gate reopens, and the NEXT
// announcement still speaks (the leftover token is the entry-drain's business).
func Test381r2PlayingPhaseInterruptAbortsAndWalks27(t *testing.T) {
	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "381-r2 test",
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

	walkToThinking381r2(t, rv)
	if !rv.announce("这是第一句话，播报需要一点时间。") {
		t.Fatalf("the announce path refused the sentence; problems=%v", rv.problemsSnapshot())
	}
	// The owner is parked in Drain with the sentence on the device: the phase
	// where the chunk-boundary select runs zero times.
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
	rv.TriggerInterrupt()
	// Same completion marker as the live case: the counter AND the #27 walk,
	// since ttsAborted is incremented one statement before the dispatch.
	voice15eWaitFor(t, "the interrupt to end the announcement (aborted or played out)", func() bool {
		return (rv.ttsAborted.Load() >= 1 && rv.machine.State() == statemachine.StateListening) ||
			rv.ttsSpoken.Load() >= 1
	})
	if got := rv.ttsAborted.Load(); got != 1 {
		t.Fatalf("MUTANT 381r2-M1 SHAPE: TriggerInterrupt did not abort the speaker in flight "+
			"(spoken=%d, aborted=%d): during playback the owner is parked in Drain where the "+
			"chunk-boundary select never runs, so without the abort hop nothing stops the audio and "+
			"the sentence plays out (381-v2's 6,653ms shape)", rv.ttsSpoken.Load(), got)
	}
	if got := rv.ttsSpoken.Load(); got != 0 {
		t.Fatalf("an aborted announcement was counted as spoken (%d)", got)
	}
	if got := rv.problemsSnapshot(); len(got) != 0 {
		t.Fatalf("a user interrupt left problems behind: %v", got)
	}
	timeline := rv.timelineSnapshot()
	last := timeline[len(timeline)-1]
	if last.Event != statemachine.EvInterrupt || last.To != statemachine.StateListening {
		t.Fatalf("the abort walked %+v, want #27 voice.interrupt -> Listening", last)
	}
	for _, tr := range timeline {
		if tr.Event == statemachine.EvSpeakDone {
			t.Fatalf("an aborted announcement also walked #26: %+v", tr)
		}
	}
	// The gate reopens in speak()'s deferred release - wait for it before
	// reading it, so the assert cannot race the owner's way out.
	voice15eWaitFor(t, "the capture gate to reopen after the abort", func() bool {
		closed, reopened := gate.count()
		return closed == 1 && reopened == 1
	})
	speakerMu.Lock()
	aborts1 := probes[0].aborts
	speakerMu.Unlock()
	if aborts1 != 1 {
		t.Fatalf("the abort hop fired %d times on the speaker, want 1", aborts1)
	}

	// AC#4's second half: the next announcement enters, speaks, and lands in
	// Warm - the leftover token from the interrupt is gone (entry drain) and the
	// cleared abort registration does not reach this sentence's speaker.
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
	t.Logf("ABORT-27 aborted=%d spoken=%d next_spoken=%d probe1_aborts=%d probe2_aborts=%d gate_closed=%d gate_reopened=%d",
		rv.ttsAborted.Load(), 0, rv.ttsSpoken.Load(), aborts1, aborts2, closed, reopened)
}

// Test381r2VetoSubstringFinalsReachTheAgent is M3's surface: the 381-r1
// string-match veto branch is deleted, so finals that merely CONTAIN 停/取消/别
// reach the agent like any other utterance - finals counted, handed to the text
// path, no D43 rejection, no interrupt edge. The real stop channels (Esc, KWS
// veto) do not ride the transcript path; a spoken stop sentence is now a task
// for the agent, which is the named tradeoff of deleting the branch.
func Test381r2VetoSubstringFinalsReachTheAgent(t *testing.T) {
	sink := &voice15eSink{}
	for i, text := range vetoCorpus381r2 {
		rv := newResidentVoice(make(chan []byte), voicePlan{}, &residentBall{}, sink)
		rv.dispatch(statemachine.EvSummon, nil, "test: D43#4 Sleeping 到 Listening")
		rv.consume([]speech.Transcript{{
			Kind: speech.TranscriptFinal,
			Text: text,
			Endpoint: speech.Endpoint{
				Kind:   speech.EndpointEnd,
				Speech: 400 * audio.TargetRate / 1000,
				Total:  400 * audio.TargetRate / 1000,
			},
		}})
		if got := rv.finals.Load(); got != 1 {
			t.Fatalf("MUTANT 381r2-M3 SHAPE: the veto substring branch came back and swallowed the "+
				"user's sentence %q (case %d: finals=%d, handed=%d): strings.Contains 停/取消/别 is not a "+
				"veto channel - it cannot see the Speaking state, cannot stop the TTS, and its continue "+
				"eats the final behind a D43 rejection (381-v2 §4: 12/12 false positives)",
				text, i, got, rv.handed.Load())
		}
		if got := rv.handed.Load(); got != 1 {
			t.Fatalf("case %d (%q): the final was counted but not handed to the text path (%d)", i, text, got)
		}
		if got := rv.problemsSnapshot(); len(got) != 0 {
			t.Fatalf("case %d (%q) polluted the problem list: %v", i, text, got)
		}
		for _, tr := range rv.timelineSnapshot() {
			if tr.Event == statemachine.EvInterrupt {
				t.Fatalf("case %d (%q) walked an interrupt edge: %+v", i, text, tr)
			}
		}
	}
	if got := len(sink.texts()); got != len(vetoCorpus381r2) {
		t.Fatalf("the text path received %d sentences, want %d", got, len(vetoCorpus381r2))
	}
	t.Logf("VETO-GONE benign=%d/%d handed=true problems=0 (false-positive rate 0/12; the 381-r1 branch measured 12/12)",
		len(sink.texts()), len(vetoCorpus381r2))
}

// Test381r2RealDevicePlayingPhaseAbortWithin400ms is the live reading AC#3
// asks for: five rounds of Esc-during-playback, each against the REAL waveOut
// device and the REAL matcha synthesis through the production TriggerInterrupt
// path. One control round first, so the play-out contrast is measured on this
// machine and not quoted from another tree. Skipped where there is no output
// device (a CI box has none).
func Test381r2RealDevicePlayingPhaseAbortWithin400ms(t *testing.T) {
	if probe, err := openWaveOutSpeaker(22050); err != nil {
		t.Skipf("no output device on this machine: %v", err)
	} else {
		_ = probe.Close()
	}
	const sentence = "任务已经完成，明天下午三点的会议改到四点，请准时参加，谢谢。"

	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "381-r2 test",
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

	// Control round: the same sentence with no interrupt, so "played out" has a
	// number on this machine (381-v2 measured 6,653ms remaining on theirs).
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
	playOut := time.Since(parkedAt)
	t.Logf("PLAYOUT-CONTROL parked_to_done=%s announce_to_done=%s",
		playOut.Round(time.Millisecond), time.Since(controlStart).Round(time.Millisecond))

	// Five abort rounds: Esc lands 300ms into the playback; the reading is the
	// distance from the trigger to the device actually stopped.
	const rounds = 5
	for r := 1; r <= rounds; r++ {
		walkToThinking381r2(t, rv)
		if !rv.announce(sentence) {
			t.Fatalf("round %d: the announce path refused the sentence; problems=%v", r, rv.problemsSnapshot())
		}
		spy := spyAt(r)
		<-spy.parked
		parkedRound := time.Now()
		spokenBefore := rv.ttsSpoken.Load()
		time.Sleep(300 * time.Millisecond) // the Esc lands 300ms into the playback
		t0 := time.Now()
		rv.TriggerInterrupt()
		// The round's completion marker is the abort branch FULLY done - the
		// counter AND the #27 walk (ttsAborted is incremented one statement
		// before the dispatch, so the counter alone can race the state read).
		voice15eWaitFor(t, fmt.Sprintf("round %d: the stop (aborted or played out)", r), func() bool {
			return (rv.ttsAborted.Load() >= uint64(r) && rv.machine.State() == statemachine.StateListening) ||
				rv.ttsSpoken.Load() > spokenBefore
		})
		if got := rv.ttsAborted.Load(); got < uint64(r) {
			t.Fatalf("MUTANT 381r2-M1 SHAPE: round %d played out instead of stopping (spoken=%d, aborted=%d): "+
				"TriggerInterrupt did not abort the speaker in flight, so the sentence played out "+
				"(381-v2's 6,653ms shape)", r, rv.ttsSpoken.Load(), got)
		}
		unstick := time.Since(t0)
		if unstick > 400*time.Millisecond {
			t.Fatalf("round %d: the device answered the interrupt in %s, want within 400ms", r, unstick)
		}
		<-spy.abortDone
		reset := time.Since(t0)
		if reset > 400*time.Millisecond {
			t.Fatalf("round %d: waveOutReset landed %s after the trigger, want within 400ms", r, reset)
		}
		if got := rv.ttsSpoken.Load(); got != spokenBefore {
			t.Fatalf("round %d: the aborted sentence was counted as spoken (%d -> %d)", r, spokenBefore, got)
		}
		// Abort's return value is a READING, not an assert: once the parked Drain
		// answers the sentinel, the owner's Close and the async Abort race for the
		// device lock, and either of them lands the reset that stops the audio
		// (both call waveOutReset). Which one won is not observable from here and
		// does not matter for the stop; the timings above are the proof.
		t.Logf("ABORT-ROUND-%d parked_to_stop=%s unstick=%s reset=%s abort_returned=%d",
			r, time.Since(parkedRound).Round(time.Millisecond), unstick.Round(time.Millisecond),
			reset.Round(time.Millisecond), spy.dropped)
	}
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

// compile-time: both test speakers satisfy the seam the announce path is
// written against.
var (
	_ voiceSpeaker = (*voice381r2ProbeSpeaker)(nil)
	_ voiceSpeaker = (*voice381r2DeviceSpy)(nil)
)
