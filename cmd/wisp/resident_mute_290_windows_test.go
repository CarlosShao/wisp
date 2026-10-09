//go:build windows

package main

// Ticket 290 AC#2/AC#4: the mute gestures reach THIS process's gate.
//
// What is under test is the hop, not the microphone. Before this leg the chain
// ended at recordBallGesture: the orb's mute hot key and the tray's mute item
// were booked by name and touched nothing, while nothing in production called
// HalfDuplexGate.SetMuted at all - the only way anything ever opened the gate was
// a test that wrote its own mic_muted_default=false. So these tests run the real
// assembly (assembleCapture -> the real HalfDuplexGate -> the real D38(e)
// registration) and the real gesture entry point (residentBall.muteGesture),
// against the capture seam ticket 247 already provides for the device side. The
// gate is not faked: it is the object production holds.
//
// What they do NOT cover, named so nobody reads them as the live reading: the
// number that varies with a real voice belongs to AC#3's owner-present window
// (form 2: a real microphone and a real hand on the key), and whether the orb
// visibly breathes belongs to ticket 68. Neither is claimed here.

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// openableSource is the capture side of the seam: it answers Start with success
// (or with the error a subtest asked for) and counts both calls. The device is
// deliberately not opened - these tests ask which object the gesture reached, and
// the live-device question is the one named above as owed to its owner.
type openableSource struct {
	mu       sync.Mutex
	starts   int
	stops    int
	startErr error
	lastErr  string
}

func (s *openableSource) Start(_ context.Context, _ chan<- []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.startErr != nil {
		return s.startErr
	}
	s.starts++
	return nil
}

func (s *openableSource) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stops++
	return nil
}

func (s *openableSource) Stats() audio.Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return audio.Stats{LastError: s.lastErr}
}

func (s *openableSource) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.startErr
}

func (s *openableSource) calls() (starts, stops int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.starts, s.stops
}

// muteRig lays the shipped default config down (mic_muted_default untouched),
// assembles the capture leg over the injected source and hands the ball host the
// executor the assembly root hands it in production.
func muteRig(t *testing.T, src captureSource) (*residentAudio, *residentBall, *observe.Registry) {
	t.Helper()
	rt, _ := bootAudioRuntime(t)
	dir := writeAudioConfig(t, nil)
	tap := &levelTap{}
	ra := assembleCapture(rt, dir, tap.call, func(...audio.CaptureOption) captureSource { return src })
	if ra.gate == nil {
		t.Fatalf("the default assembled no gate: %q", ra.verdict)
	}
	rb := &residentBall{}
	rb.attachMuteGate(ra.toggleMute)
	return ra, rb, rt.Registry
}

// TestAC290BothMuteGesturesTurnTheGateAndBack is AC#2's capability reading, taken
// through the host entry point both gestures use: the gesture opens the gate that
// booted muted, the state is read back off the gate (never off the ball host), and
// pressing again closes it. Both names the ball surface has - the global hot key
// and the tray item - go through the same hop.
func TestAC290BothMuteGesturesTurnTheGateAndBack(t *testing.T) {
	src := &openableSource{}
	ra, rb, reg := muteRig(t, src)

	// The default posture is the one ticket 247 shipped and this leg may not
	// move: muted at boot, device never handed over, no capture thread booked.
	if !ra.gate.Muted() || ra.gate.Open() {
		t.Fatalf("the rig did not start in the shipped posture: muted=%v open=%v", ra.gate.Muted(), ra.gate.Open())
	}
	if st, sp := src.calls(); st != 0 || sp != 0 {
		t.Fatalf("a muted boot touched the source (%d starts, %d stops)", st, sp)
	}

	for i, gesture := range []string{"mute-hotkey", "tray-mute"} {
		wantOpens := i + 1 // one device handoff per open, and no extra one
		outcome := rb.muteGesture(gesture)
		if ra.gate.Muted() {
			t.Fatalf("%s: the gate still reports muted after the gesture: %q", gesture, outcome)
		}
		if !ra.gate.Open() {
			t.Fatalf("%s: the gate is unmuted but never handed the device over: %q", gesture, outcome)
		}
		if !ra.started {
			t.Fatalf("%s: the leg's derived mirror says nothing is running while the gate is open", gesture)
		}
		if !strings.Contains(outcome, "已取消静音") || !strings.Contains(outcome, "设备已交接") {
			t.Fatalf("%s: the outcome does not say what the gate now is: %q", gesture, outcome)
		}
		st, _ := src.calls()
		if st != wantOpens {
			t.Fatalf("%s: the source was started %d times, want %d (one device handoff per open)", gesture, st, wantOpens)
		}

		// And back the other way: the same key is the way out again
		// (docs/PLAN.md's Muted row lists pressing the mute key again as its exit).
		outcome = rb.muteGesture(gesture)
		if !ra.gate.Muted() {
			t.Fatalf("%s: the second press did not put the gate back to muted: %q", gesture, outcome)
		}
		if ra.gate.Open() || ra.started {
			t.Fatalf("%s: capture stayed handed over after muting (open=%v started=%v)", gesture, ra.gate.Open(), ra.started)
		}
		if !strings.Contains(outcome, "已静音") || !strings.Contains(outcome, "设备未打开") {
			t.Fatalf("%s: the muted outcome is not the gate's state: %q", gesture, outcome)
		}
		_, sp := src.calls()
		if sp != wantOpens {
			t.Fatalf("%s: muting closed the source %d times, want %d (one close per open)", gesture, sp, wantOpens)
		}
	}

	// Zero new goroutines is AC#4's shape, and the ruler with teeth here is the
	// registry: this leg books no thread of its own and the injected source spawns
	// none, so a non-zero count would mean the gesture path invented one.
	if n := reg.CountByName("audio-capture"); n != 0 {
		t.Fatalf("the mute path booked %d audio-capture threads beyond the capture leg's own roster name", n)
	}
}

// TestAC290OutcomeIsReadOffTheGateNotOffTheRequest is the truth-source half: an
// open the device refused must not come back as a success, because the sentence
// the user gets is the only evidence they have that the microphone is or is not
// live.
func TestAC290OutcomeIsReadOffTheGateNotOffTheRequest(t *testing.T) {
	want := audio.DeviceError(0x80070005, audio.DeviceDescriptor{ID: "id-290", Name: "Blocked Mic"}, "capture open")
	src := &openableSource{startErr: want, lastErr: want.Error()}
	ra, rb, _ := muteRig(t, src)

	outcome := rb.muteGesture("mute-hotkey")
	if ra.gate.Open() {
		t.Fatal("the gate claims it handed a refused device over")
	}
	if ra.started {
		t.Fatal("the leg's derived mirror says capture is running over a refused device")
	}
	if !strings.Contains(outcome, "设备未交接") {
		t.Fatalf("a failed open was reported without saying it failed: %q", outcome)
	}
	if !strings.Contains(outcome, string(observe.ClassAudioDevice)) || !strings.Contains(outcome, "Privacy & security") {
		t.Fatalf("the failure lost the error class or the guidance text: %q", outcome)
	}
}

// TestAC290NoGateSaysWhichShapeThisProcessIsIn covers the legs that own no gate:
// voice disabled (no collector built at all) and a nil handle. The gesture must
// say the reason, not perform a mute nobody asked for.
func TestAC290NoGateSaysWhichShapeThisProcessIsIn(t *testing.T) {
	rt, _ := bootAudioRuntime(t)
	dir := writeAudioConfig(t, func(c *config.Config) { c.Voice.Enabled = false })
	ra := assembleCapture(rt, dir, (&levelTap{}).call, newRealCaptureSource)
	if ra.gate != nil {
		t.Fatal("voice.enabled=false still built a gate for the gesture to turn")
	}

	outcome, executed := ra.toggleMute()
	if executed {
		t.Fatalf("a leg with no collector claimed to have turned a gate: %q", outcome)
	}
	if !strings.Contains(outcome, "未构造") {
		t.Fatalf("the gesture lost the reason this process has no microphone leg: %q", outcome)
	}
	// The gesture entry point takes that answer without inventing a state of its
	// own, and it hands the same words back to whoever pressed the key.
	rb := &residentBall{}
	rb.attachMuteGate(ra.toggleMute)
	if got := rb.muteGesture("mute-hotkey"); got != outcome {
		t.Fatalf("the host changed the sentence it was given: %q vs %q", got, outcome)
	}
	var nilLeg *residentAudio
	if noGate, ok := nilLeg.toggleMute(); ok || noGate == "" {
		t.Fatalf("a nil capture handle answered the gesture with (%q, %v)", noGate, ok)
	}
}

// TestAC290GestureBeforeTheAttachSaysSo is the boot window the setter shape
// creates: the window and its hot key are live before the capture leg is
// assembled, and a key pressed in that gap must be reported as what it was.
func TestAC290GestureBeforeTheAttachSaysSo(t *testing.T) {
	rb := &residentBall{}
	if fn := rb.currentMuteGate(); fn != nil {
		t.Fatal("a host that was handed nothing reports an executor")
	}
	got := rb.muteGesture("mute-hotkey")
	if !strings.Contains(got, "no capture leg had been assembled") {
		t.Fatalf("an early key press was reported as something else: %q", got)
	}
	if rb.muteGate != nil {
		t.Fatal("reading the executor attached one")
	}

	// A host with no window at all is still safe to ask, because nil-safety is
	// the posture every residentBall read keeps.
	var nilBall *residentBall
	nilBall.attachMuteGate(func() (string, bool) { return "", false })
	nilBall.muteGesture("tray-mute")
}

// TestAC290MutedDefaultStillDecidesTheBoot is AC#4's "not one notch" half read
// from this ticket's side: adding a gesture that can open the gate must not have
// moved the shipped default, and the boot must still be the muted posture the
// privacy ruling picked.
func TestAC290MutedDefaultStillDecidesTheBoot(t *testing.T) {
	c := config.NewDefaults()
	if !c.Audio.MicMutedDefault || !c.Voice.Enabled || c.Voice.WakeWord.Enabled {
		t.Fatalf("a default this ticket was told not to touch moved: mic_muted_default=%v voice.enabled=%v wake_word.enabled=%v",
			c.Audio.MicMutedDefault, c.Voice.Enabled, c.Voice.WakeWord.Enabled)
	}

	src := &openableSource{}
	ra, rb, _ := muteRig(t, src)
	if !ra.gate.Muted() {
		t.Fatal("the gate booted unmuted: a double click would open the microphone")
	}
	st, _ := src.calls()
	if st != 0 {
		t.Fatalf("the device was opened %d times at boot without the user asking", st)
	}
	if !strings.Contains(ra.posture(), "设备未打开") {
		t.Fatalf("the boot line stopped saying the device is closed: %q", ra.posture())
	}
	// The gesture is the only thing that may leave that posture, and it is the
	// same toggle the tests above drive.
	if got := rb.muteGesture("mute-hotkey"); !strings.Contains(got, "已取消静音") {
		t.Fatalf("the user's gesture cannot leave the boot posture: %q", got)
	}
}

// TestAC290UnhostedGestureWordingStillSaysTrueThing guards the sentence this leg
// had to re-derive: it may not keep claiming this process has no microphone, and
// it may not claim the mute gestures are unhosted now that they are not.
func TestAC290UnhostedGestureWordingStillSaysTrueThing(t *testing.T) {
	if strings.Contains(ballGestureWhy, "no microphone") {
		t.Fatalf("ballGestureWhy still says this process has no microphone: %q", ballGestureWhy)
	}
	for _, want := range []string{"ticket 290", "capture leg"} {
		if !strings.Contains(ballGestureWhy, want) {
			t.Fatalf("ballGestureWhy does not name %q among the gestures with executors: %q", want, ballGestureWhy)
		}
	}
	// The reason an unhosted gesture gives is one string said in two places
	// (console line and log record), so it cannot drift into a claim the other
	// half does not make; that single const is what this test reads.
	if !strings.Contains(ballGestureWhy, "recorded by name") {
		t.Fatalf("ballGestureWhy lost the sentence both places say: %q", ballGestureWhy)
	}
}
