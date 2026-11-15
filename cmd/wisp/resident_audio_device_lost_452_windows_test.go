//go:build windows

package main

// Ticket 452, AC#1: the two shipped disconnects become judgeable facts BEFORE
// either production line is touched.
//
//   (a) the shipped build carries no AEC3 binding: scripts/build.ps1's
//       `go build` passes no -tags webrtcap, so internal/audio/aec compiles
//       aec_stub.go, whose NewApm answers ErrUnavailable for ever. The
//       consequence the product never stated: the shipped binary has NO
//       echo cancellation, so Path-C barge-in cannot happen at all. AC#2's
//       form ② is that the product SAY it; this file pins the statement.
//
//   (b) a microphone unplugged AFTER the device was handed over is read by
//       nobody. D43 row #14 (internal/statemachine/table.go:97,
//       Listening + EvAudioDeviceLost -> Error) exists, EvAudioDeviceLost
//       had zero dispatchers in the production tree, and the capture leg's
//       terminal error only reached the BOOT posture switch. The reading
//       below drives the production consumer loop (the owner goroutine that
//       walks D43 edges), hands it a live device, then makes it die, and
//       demands the machine and the ball leave Listening.
//
// Both readings are taken off production symbols only: the assembly
// (assembleCapture), the seam (captureSource.Err), the monitor
// (voiceBargeinMonitor.reason), the consumer (residentVoice.run) and the D43
// machine. Nothing here is a mock standing in for the fact under test.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// aecBindingInThisBuild452 asks the package the same question the monitor asks
// at arm time (aec.NewApm): the stub answers ErrUnavailable, the native binding
// answers an engine. This is AC#1(a)'s measurement taken with the production
// call itself, so no constant can drift away from what the binary really does.
func aecBindingInThisBuild452(t *testing.T) error {
	t.Helper()
	apm, err := aec.NewApm(audio.TargetRate)
	if err == nil && apm != nil {
		_ = apm.Close()
	}
	return err
}

// lossSource452 is the capture side of AC#1(b): a device that is handed over
// healthy and dies later. The three telemetry reads are the same three the
// production WASAPIMicrophone exposes (Start/Stop, Stats().LastError, Err()),
// and the error object is the package's own canonical audio_device mapping -
// the exact value a failed hotplug reopen posts (wasapimic_windows.go:313).
type lossSource452 struct {
	mu  sync.Mutex
	err error
}

func (s *lossSource452) Start(context.Context, chan<- []byte) error { return nil }
func (s *lossSource452) Stop() error                                { return nil }

func (s *lossSource452) Stats() audio.Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err == nil {
		return audio.Stats{}
	}
	return audio.Stats{LastError: s.err.Error()}
}

func (s *lossSource452) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

// lose is the unplug: from this call on the source answers Err() with err,
// exactly as WASAPIMicrophone.postErr does after its one reopen attempt failed.
func (s *lossSource452) lose(err error) {
	s.mu.Lock()
	s.err = err
	s.mu.Unlock()
}

// TestAC452ShippedBuildSaysItHasNoEchoCancellation is AC#1(a) plus AC#2's
// form-② statement, in one place so the two cannot drift apart: the shipped
// build has no binding (pinned off the script and the build tag), and the
// product says so instead of letting the configuration surface imply a
// capability the binary does not have.
func TestAC452ShippedBuildSaysItHasNoEchoCancellation(t *testing.T) {
	// (a1) the build tag the shipped build really carries: none.
	scriptPath := filepath.Join("..", "..", "scripts", "build.ps1")
	raw, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("read %s: %v", scriptPath, err)
	}
	script := string(raw)
	if strings.Contains(script, "webrtcap") {
		t.Fatalf("452 AC#1(a) RED (ruler): scripts/build.ps1 now mentions webrtcap. The shipped shape moved; this reading and the product statement below must move with it, deliberately (AC#2 is a ruling, not an accident)")
	}
	var buildLine string
	for _, line := range strings.Split(script, "\n") {
		if strings.Contains(line, "build -trimpath") {
			buildLine = strings.TrimSpace(line)
		}
	}
	if buildLine == "" {
		t.Fatalf("452 AC#1(a) RED (ruler): no `go build -trimpath` line found in %s - the shipped build command moved and this reading now measures nothing", scriptPath)
	}
	if strings.Contains(buildLine, "-tags") {
		t.Fatalf("452 AC#1(a) RED: the shipped build line carries tags now: %q. If the factory build really gained the AEC3 binding, AC#2's statement below is stale and must be replaced by that fact - not left saying the opposite", buildLine)
	}

	// (a2) the fact measured through the production call: this binary's NewApm
	// is the stub, which is what the tagged-off shipped build resolves.
	if err := aecBindingInThisBuild452(t); !errors.Is(err, aec.ErrUnavailable) {
		t.Fatalf("452 AC#1(a) RED (ruler): this test binary's aec.NewApm answered %v, not ErrUnavailable - so this is NOT the shipped shape and the statement below is not owed here. Run the package without -tags webrtcap to read the factory posture", err)
	}

	// (b) the product statement AC#2 form ② requires: the absence is said, not
	// left for the user to discover.
	rv := newResidentVoice(nil, voicePlan{}, &residentBall{}, nil)
	if rv == nil || rv.bargein == nil {
		t.Fatal("452 AC#1(a) RED: newResidentVoice built no barge-in monitor, so there is no posture surface to state the absence on")
	}
	if !strings.Contains(rv.bargein.reason, "本构建无回声消除") {
		t.Fatalf("452 AC#1(a)/AC#2 RED: this build has no AEC3 binding (aec.NewApm answers ErrUnavailable, scripts/build.ps1 passes no -tags webrtcap), and the barge-in monitor's posture at construction is %q: the product states neither that this binary cannot cancel echo nor that Path-C barge-in is therefore unavailable. A build without the capability may not stay silent about it (AC#2 form ②), and the 400ms-class readings belong to the dedicated -tags webrtcap staging, not to what ships",
			rv.bargein.reason)
	}
}

// TestAC452LostMicWalksD43FromListeningIntoError is AC#1(b)/AC#3: the
// production consumer loop is started (the same body the boot books on the
// D38b roster), the assembled capture leg hands over a healthy device, and
// then the device dies. The reading demands two things at once: the machine
// walks D43#14 to Error, and the ball is projected onto that state.
//
// The healthy phase is half the reading: a probe that fires on a live device
// would throw the ball into Error on every boot, so "Listening is still
// Listening one probe interval into a healthy run" is asserted before the
// unplug.
func TestAC452LostMicWalksD43FromListeningIntoError(t *testing.T) {
	rt, _ := bootAudioRuntime(t)
	// An unmuted default is deliberate and it is not the shipped table: the
	// device has to be handed over for a RUNTIME loss to exist at all.
	dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })

	src := &lossSource452{}
	ra := assembleCapture(rt, dir, (&levelTap{}).call, func(...audio.CaptureOption) captureSource { return src })
	if ra.gate == nil || ra.mic == nil {
		t.Fatalf("the capture leg did not assemble: verdict=%q", ra.verdict)
	}
	if !ra.started.Load() || !ra.gate.Open() {
		t.Fatalf("the capture leg never handed the device over (started=%v open=%v verdict=%q): a runtime loss cannot be measured on a device that was never running",
			ra.started.Load(), ra.gate.Open(), ra.verdict)
	}

	ball := &residentBall{}
	rv := newResidentVoice(nil, voicePlan{}, ball, nil)
	// Listening is the row's From state and the state the ticket names: a ball
	// listening to a microphone that has been unplugged.
	rv.machine = statemachine.New(statemachine.Options{Initial: statemachine.StateListening})

	// The seam this reading hands the loop is the same line startResidentVoice
	// runs before it books the consumer; a hand-wired test on top of a
	// production line nobody calls would be the vacuous-green shape this repo
	// counts. So both are demanded: the wiring line exists in the production
	// file, and the loop really walks the edge.
	prodSrc, err := os.ReadFile("resident_voice_windows.go")
	if err != nil {
		t.Fatalf("read resident_voice_windows.go: %v", err)
	}
	if !strings.Contains(string(prodSrc), "rv.deviceErr = ra.captureErr") {
		t.Fatalf("452 AC#3 RED: resident_voice_windows.go has no `rv.deviceErr = ra.captureErr` line, so no boot actually hands the health probe a capture leg: the machine below would only walk because the TEST wired the seam")
	}
	rv.deviceErr = ra.captureErr
	ra.attachVoice(rv)

	// The production consumer: the owner goroutine every D43 walk in this leg
	// runs on. No VAD is handed in because this reading feeds no frames - the
	// health read is the only work this loop is asked to do.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		rv.run(ctx, nil)
	}()
	defer func() {
		cancel()
		<-done
	}()

	// The healthy phase: one probe interval with a live device must not move
	// the machine.
	time.Sleep(400 * time.Millisecond)
	if got := rv.machine.State(); got != statemachine.StateListening {
		t.Fatalf("452 AC#1(b) RED: a HEALTHY device moved the machine out of Listening into %s within one probe interval: the probe invents a loss, which would walk every boot into Error", got)
	}

	dev := audio.DeviceDescriptor{ID: "id-452", Name: "USB Mic 452"}
	lossErr := audio.DeviceError(0x88890004, dev, "capture device lost")
	src.lose(lossErr)

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) && rv.machine.State() != statemachine.StateError {
		time.Sleep(10 * time.Millisecond)
	}
	if got := rv.machine.State(); got != statemachine.StateError {
		t.Fatalf("452 AC#1(b) RED: the microphone was unplugged after the device was handed over (mic.Err() now answers %q) while the machine sat in Listening, and 3s later the state is still %s. Nothing in the production tree reads the capture leg's terminal error, so D43#14 (Listening + EvAudioDeviceLost -> Error, table.go:97) has zero dispatchers and the ball stays in Listening with a microphone that no longer exists. Ball states: %v; problems: %v",
			lossErr.Error(), got, ball.voiceStatesSnapshot(), rv.problemsSnapshot())
	}

	// The ball is the D43 projection surface; the walk is only half the fact.
	states := ball.voiceStatesSnapshot()
	if len(states) == 0 || states[len(states)-1] != statemachine.StateError {
		t.Fatalf("452 AC#3 RED: the machine reached %s but the ball was projected %v: the walk happened without the projection hop, so the orb still shows what it showed", statemachine.StateError, states)
	}

	// And the walk is the TABLE's edge, not a fabricated one: the timeline
	// names the event and the row's From/To pair.
	var walked bool
	for _, tr := range rv.timelineSnapshot() {
		if tr.Event == statemachine.EvAudioDeviceLost && tr.From == statemachine.StateListening && tr.To == statemachine.StateError {
			walked = true
		}
	}
	if !walked {
		t.Fatalf("452 AC#3 RED: the machine reached Error but no timeline entry records Listening + EvAudioDeviceLost -> Error: %v", rv.timelineSnapshot())
	}
}
