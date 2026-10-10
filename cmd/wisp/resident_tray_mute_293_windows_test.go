//go:build windows

package main

// Ticket 293 AC#2: the tray's "静音" checkmark mirrors the gate.
//
// What is under test is the projection hop, on top of ticket 290's rig: the gate
// is the REAL object production holds (assembleCapture mounts it, the injected
// captureSource only stands in for the device, exactly as in
// resident_mute_290_windows_test.go), the read half is the real
// residentAudio.trayMuteState, and the flips are driven through the one entry
// point both user gestures use (residentBall.muteGesture) plus the boot line the
// assembly root runs.
//
// The single thing these tests stand in for is the display surface: cmd/wisp's
// harness has no ball window (no test in this package has ever built a real one,
// and Ball.SetTrayChecks queues onto an STA thread a zero-value Ball does not
// have), so the push half is a recorder. Two things pin that recorder to the real
// setter rather than letting it drift: the compile-time assertion below, which
// makes (*ball.Ball).SetTrayChecks's exact signature a checked fact, and the
// untouched ball-side chain the value travels down (internal/ball/ball_windows.go:
// SetTrayChecks writes trayMuted, the right-click handler hands trayMuted to
// showMenu, tray_windows.go's appendItem turns it into MF_CHECKED). Whether the
// user SEES the tick is AC#5's owner-present reading, and no test here claims it.

import (
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
)

// The display surface this hop writes to is the ball's own exported setter. If
// that method is renamed, changes arity or changes types, this line stops
// compiling - which is the point: the recorder below is a stand-in for a target
// whose shape is checked against the real thing, not against a wish.
var _ func(*ball.Ball, bool, bool) = (*ball.Ball).SetTrayChecks

// trayCheck is one write to the tray's two checkmarks.
type trayCheck struct {
	muted      bool
	pausedWake bool
}

// trayCheckRec records what the projection pushed. It is the display side only:
// every assertion about the state being projected compares against ra.gate.Muted(),
// never against this recorder, so a projection that ignored the gate could not
// agree with itself and pass.
type trayCheckRec struct {
	mu   sync.Mutex
	seen []trayCheck
}

func (r *trayCheckRec) push(muted, pausedWake bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seen = append(r.seen, trayCheck{muted: muted, pausedWake: pausedWake})
}

func (r *trayCheckRec) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.seen)
}

func (r *trayCheckRec) last(t *testing.T) trayCheck {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.seen) == 0 {
		t.Fatal("the tray checkmark was never written - the projection did not fire")
	}
	return r.seen[len(r.seen)-1]
}

// trayMuteRig is ticket 290's mute rig plus ticket 293's projection pair, attached
// the way the assembly root attaches them (read off the capture leg, push into the
// tray display surface).
func trayMuteRig(t *testing.T, src captureSource) (*residentAudio, *residentBall, *trayCheckRec) {
	t.Helper()
	ra, rb, _ := muteRig(t, src)
	rec := &trayCheckRec{}
	rb.attachTrayMuteProjection(ra.trayMuteState, rec.push)
	return ra, rb, rec
}

// assertMirrorsGate is the ruler every start point is measured with: the bool the
// menu will be built from must equal the gate's own answer at that moment, and the
// pause-wake half must still say false (禁区② - see mirrorTrayMute's comment).
func assertMirrorsGate(t *testing.T, ra *residentAudio, rec *trayCheckRec, where string) {
	t.Helper()
	got := rec.last(t)
	if got.muted != ra.gate.Muted() {
		t.Fatalf("%s: the tray's mute checkmark is %v while the gate reports %v - the checkmark must mirror gate.Muted()",
			where, got.muted, ra.gate.Muted())
	}
	if got.pausedWake {
		t.Fatalf("%s: the projection wrote a pause-wake checkmark this process has no truth source for", where)
	}
}

// TestAC293TrayItemCheckmarkFollowsTheGate is start point one: the tray's own mute
// item, which since ticket 290 really turns the gate, must leave the checkmark
// equal to the gate it just changed - both on the way out of mute and back.
func TestAC293TrayItemCheckmarkFollowsTheGate(t *testing.T) {
	src := &openableSource{}
	ra, rb, rec := trayMuteRig(t, src)
	if !ra.gate.Muted() {
		t.Fatalf("the rig did not boot in the shipped muted posture: %q", ra.posture())
	}

	if got := rb.muteGesture("tray-mute"); !strings.Contains(got, "已取消静音") {
		t.Fatalf("the tray item did not open the gate: %q", got)
	}
	assertMirrorsGate(t, ra, rec, "tray item, after unmuting")
	if rec.count() != 1 {
		t.Fatalf("one accepted flip wrote %d checkmarks, want exactly 1", rec.count())
	}

	if got := rb.muteGesture("tray-mute"); !strings.Contains(got, "已静音") {
		t.Fatalf("the tray item did not put the gate back: %q", got)
	}
	assertMirrorsGate(t, ra, rec, "tray item, after muting again")
	if rec.count() != 2 {
		t.Fatalf("two accepted flips wrote %d checkmarks, want exactly 2", rec.count())
	}
}

// TestAC293MuteHotkeyCheckmarkFollowsTheGate is start point two: the orb's global
// mute hot key reaches the same gate through the same entry point, and the checkmark
// the tray would show afterwards has to agree with it. A user who mutes with the
// key and then right-clicks the icon is exactly the shape this covers.
func TestAC293MuteHotkeyCheckmarkFollowsTheGate(t *testing.T) {
	src := &openableSource{}
	ra, rb, rec := trayMuteRig(t, src)

	if got := rb.muteGesture("mute-hotkey"); !strings.Contains(got, "已取消静音") {
		t.Fatalf("the hot key did not open the gate: %q", got)
	}
	assertMirrorsGate(t, ra, rec, "mute hot key, after unmuting")

	if got := rb.muteGesture("mute-hotkey"); !strings.Contains(got, "已静音") {
		t.Fatalf("the hot key did not close the gate: %q", got)
	}
	assertMirrorsGate(t, ra, rec, "mute hot key, after muting")
}

// TestAC293BootProjectionShowsTheFactoryMutePosture is start point three: the
// factory mute posture never passes through a gesture, so the boot hop is the only
// thing that can put the checkmark right before the user's first right-click. It
// must read the gate the shipped default mounted, and it must do so without
// touching the device (mirroring a state is not changing one).
func TestAC293BootProjectionShowsTheFactoryMutePosture(t *testing.T) {
	if defaults := config.NewDefaults(); !defaults.Audio.MicMutedDefault {
		t.Fatalf("audio.mic_muted_default moved off the shipped true, and this test's premise with it: %v",
			defaults.Audio.MicMutedDefault)
	}
	src := &openableSource{}
	ra, rb, rec := trayMuteRig(t, src)

	if !rb.mirrorTrayMute() {
		t.Fatal("the boot projection wrote nothing even though this leg owns a gate")
	}
	assertMirrorsGate(t, ra, rec, "boot, factory posture")
	if !rec.last(t).muted {
		t.Fatal("the booted checkmark says not muted while the gate booted muted - a double click would show a lie")
	}
	if st, sp := src.calls(); st != 0 || sp != 0 {
		t.Fatalf("projecting the checkmark touched the capture source (%d starts, %d stops); a mirror reads, it does not turn", st, sp)
	}
	if !ra.gate.Muted() {
		t.Fatal("the gate is no longer muted after the boot projection")
	}
}

// TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome is the truth-source
// half with teeth: when the user asks to unmute and the device refuses the open,
// the gate still reports its own flag, and the checkmark must copy THAT and not the
// optimistic reading of the request. Reading the checkmark off the outcome sentence
// (or off a hoped-for flip) passes nothing here.
func TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome(t *testing.T) {
	want := audio.DeviceError(0x80070005, audio.DeviceDescriptor{ID: "id-293", Name: "Blocked Mic"}, "capture open")
	src := &openableSource{startErr: want, lastErr: want.Error()}
	ra, rb, rec := trayMuteRig(t, src)

	got := rb.muteGesture("tray-mute")
	if !strings.Contains(got, "设备未交接") {
		t.Fatalf("the refused open was not reported as refused: %q", got)
	}
	assertMirrorsGate(t, ra, rec, "refused device open")
	if rec.last(t).muted != ra.gate.Muted() || ra.gate.Muted() {
		t.Fatalf("gate.Muted()=%v, checkmark=%v: the projection must copy the gate, not the failure and not the request",
			ra.gate.Muted(), rec.last(t).muted)
	}
}

// TestAC293NoGateWritesNoCheckmark covers the legs that own no gate: the projection
// has nothing to mirror, so it writes nothing rather than defaulting a checkmark
// into existence. This is the shape that keeps 禁区① honest - the answer comes from
// the gate's absence, not from a second flag like residentAudio's dead mutedAtBoot.
func TestAC293NoGateWritesNoCheckmark(t *testing.T) {
	rt, _ := bootAudioRuntime(t)
	dir := writeAudioConfig(t, func(c *config.Config) { c.Voice.Enabled = false })
	ra := assembleCapture(rt, dir, (&levelTap{}).call, newRealCaptureSource)
	if ra.gate != nil {
		t.Fatal("voice.enabled=false still built a gate for the projection to read")
	}
	rec := &trayCheckRec{}
	rb := &residentBall{}
	rb.attachMuteGate(ra.toggleMute)
	rb.attachTrayMuteProjection(ra.trayMuteState, rec.push)

	if rb.mirrorTrayMute() {
		t.Fatal("a process with no gate reported having projected one")
	}
	if rec.count() != 0 {
		t.Fatalf("the no-gate leg wrote %d checkmarks, want 0", rec.count())
	}
	if muted, ownsGate := ra.trayMuteState(); muted || ownsGate {
		t.Fatalf("the read half answered (%v, %v) for a leg that owns no gate", muted, ownsGate)
	}
	// The same answer from a nil handle, because the boot may never have assembled
	// a leg at all.
	var nilLeg *residentAudio
	if muted, ownsGate := nilLeg.trayMuteState(); muted || ownsGate {
		t.Fatalf("a nil capture handle answered the projection with (%v, %v)", muted, ownsGate)
	}
	// And the gesture path says the true thing without writing a checkmark either.
	if got := rb.muteGesture("tray-mute"); !strings.Contains(got, "未构造") {
		t.Fatalf("the no-gate gesture lost its reason: %q", got)
	}
	if rec.count() != 0 {
		t.Fatalf("a rejected gesture wrote %d checkmarks, want 0", rec.count())
	}
}

// TestAC293ProjectionIsNilSafeWithoutABallWindow is 禁区③'s shape: the host may be
// asked to project before it holds a window, and the assembly root's guard is what
// keeps the method value from ever being taken on a nil ball. Nothing here panics,
// and nothing is written.
func TestAC293ProjectionIsNilSafeWithoutABallWindow(t *testing.T) {
	src := &openableSource{}
	ra, _, _ := trayMuteRig(t, src)

	// Nothing attached yet: the boot window in which the hot key is already live.
	rb := &residentBall{}
	if rb.mirrorTrayMute() {
		t.Fatal("a host with nothing attached reported a projection")
	}
	// Read half attached, display half missing (the no-window boot branch never
	// attaches either; this is the same answer from the other side).
	rb.attachTrayMuteProjection(ra.trayMuteState, nil)
	if rb.mirrorTrayMute() {
		t.Fatal("a host with no display surface wrote a checkmark")
	}
	if rb.b != nil {
		t.Fatal("the rig grew a ball window")
	}
	// A nil host is safe to ask, like every other residentBall read.
	var nilBall *residentBall
	nilBall.attachTrayMuteProjection(ra.trayMuteState, (&trayCheckRec{}).push)
	if nilBall.mirrorTrayMute() {
		t.Fatal("a nil host reported a projection")
	}
	// The gate really did not move while all of that was being asked.
	if !ra.gate.Muted() {
		t.Fatal("asking about the checkmark turned the gate")
	}
}

// TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth names the bug this ticket replaced:
// Ball.SetTrayChecks had no production caller, so the ball-side flag stayed at its
// zero value no matter what the gate did. The projection is now that caller, and the
// asymmetry it keeps (pause-wake stays false) is asserted on every write above.
func TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth(t *testing.T) {
	src := &openableSource{}
	ra, rb, rec := trayMuteRig(t, src)

	// Three flips, three writes, each one equal to the gate at that instant.
	for i := 0; i < 3; i++ {
		rb.muteGesture("mute-hotkey")
		assertMirrorsGate(t, ra, rec, "flip")
	}
	if rec.count() != 3 {
		t.Fatalf("%d checkmarks written for 3 accepted flips", rec.count())
	}
	// A flip the gate refused to accept writes nothing (the gesture answered
	// executed=false), so the checkmark can never drift ahead of the gate.
	nilLeg := &residentAudio{}
	rb.attachMuteGate(nilLeg.toggleMute)
	before := rec.count()
	rb.muteGesture("tray-mute")
	if rec.count() != before {
		t.Fatalf("a gesture the leg did not execute still moved the checkmark (%d -> %d writes)", before, rec.count())
	}
}
