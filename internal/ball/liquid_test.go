package ball

import (
	"math"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// The ticket 62 audio gate: rotation amplitude must be monotonic in the level
// signal (AC#4), the liquid must converge and the border must arrive when the
// voice stops, and none of it may run outside the session states.

func stepMotion(m *liquidMotion, level float32, n int, dt time.Duration) {
	for i := 0; i < n; i++ {
		m.pushLevel(level, dt)
	}
}

func TestLiquidRotationIsMonotonicInLevel(t *testing.T) {
	const frames = 30
	const dt = 33 * time.Millisecond
	var prev float64
	for _, level := range []float32{0, 0.15, 0.35, 0.6, 0.85, 1} {
		m := &liquidMotion{active: true}
		stepMotion(m, level, frames, dt)
		travel := m.angle
		if travel <= prev {
			t.Fatalf("rotation travel must grow with the level: level %.2f gave %.4f rad, previous level gave %.4f",
				level, travel, prev)
		}
		prev = travel
	}
}

func TestLiquidConvergesAndBorderFadesInOnSilence(t *testing.T) {
	m := &liquidMotion{active: true}
	m.beginSession()
	if m.border != 0 {
		t.Fatalf("a fresh summon must have no border yet, got %v", m.border)
	}
	// Speak loudly: the liquid gathers, the border stays away.
	stepMotion(m, 0.9, 20, 33*time.Millisecond)
	if m.border != 0 {
		t.Fatalf("while speaking the border must stay out, got %v", m.border)
	}
	if m.level < 0.5 {
		t.Fatalf("a loud feed must lift the envelope, got %v", m.level)
	}
	// Go silent: the border must arrive through the transition, in
	// BorderOpenMs of simulated time - not instantly, not never.
	stepMotion(m, 0, BorderOpenMs+SilenceHoldMs+20, time.Millisecond)
	if m.border != 1 {
		t.Fatalf("the border must be fully in BorderOpenMs after the silence hold, got %v", m.border)
	}
	// The liquid itself converges more slowly on purpose (release tau): the
	// blob glide carries over the gap between syllables.
	stepMotion(m, 0, int(LevelReleaseTauMs)*4, time.Millisecond)
	if m.level > SilenceLevelGate {
		t.Fatalf("sustained silence must decay the envelope below the gate, got %v", m.level)
	}
	// Speaking again pulls the border back out.
	stepMotion(m, 1, 30, 33*time.Millisecond)
	if m.border > 0.1 {
		t.Fatalf("the voice must hand the orb back from the border, got %v", m.border)
	}
}

func TestSummonBurstIsBoundedAndRetires(t *testing.T) {
	m := &liquidMotion{active: true}
	m.beginSession()
	if !m.busy(statemachine.StateListening) {
		t.Fatal("a just-summoned orb must be producing frames")
	}
	if m.summon != 1 {
		t.Fatalf("summon burst must start full, got %v", m.summon)
	}
	stepMotion(m, 0, SummonFlowMs+400, time.Millisecond)
	if m.summon != 0 {
		t.Fatalf("the summon burst must end by itself, got %v", m.summon)
	}
	if m.busy(statemachine.StateListening) {
		t.Fatal("motion must go idle once the burst landed and the border is in")
	}
	// A loud level keeps it busy only while something still ramps.
	m2 := &liquidMotion{active: true}
	m2.beginSession()
	stepMotion(m2, 1, SummonFlowMs+400, time.Millisecond)
	if m2.busy(statemachine.StateListening) {
		t.Fatal("a steady loud level with no ramp left must not hold a timer open")
	}
}

func TestMotionOnlyAppliesInPrototypeMode(t *testing.T) {
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	EnablePrototypeVisuals(false)
	m := &liquidMotion{active: true}
	m.beginSession()
	stepMotion(m, 0.8, 10, 33*time.Millisecond)
	v := Visual{Glass: true, BorderAlpha: 1}
	m.applyTo(&v)
	if v.LiquidLevel != 0 || v.LiquidAngle != 0 || v.SummonFlow != 0 || v.BorderAlpha != 1 {
		t.Fatalf("frozen mode must not receive prototype motion: %+v", v)
	}
	EnablePrototypeVisuals(true)
	v2 := Visual{Glass: true, BorderAlpha: 1}
	m.applyTo(&v2)
	if v2.LiquidLevel == 0 || v2.SummonFlow == 0 {
		t.Fatalf("prototype mode must carry the live motion: %+v", v2)
	}
	if v2.BorderAlpha >= 1 {
		t.Fatalf("the border must be mid-transition, not already full: %v", v2.BorderAlpha)
	}
}

// allKnownStates is the D43 roster (20 states), used by the gates below so a
// new state cannot silently join the timer set or the border set.
var allKnownStates = []statemachine.State{
	statemachine.StateFirstRun, statemachine.StateSleeping, statemachine.StateArmed,
	statemachine.StateMuted, statemachine.StateListening, statemachine.StateThinking,
	statemachine.StateActing, statemachine.StateSpeaking, statemachine.StateWarm,
	statemachine.StateConversation, statemachine.StateConfirming,
	statemachine.StateAwaitingApproval, statemachine.StateSettling,
	statemachine.StateDownloading, statemachine.StateError, statemachine.StateNoNetwork,
	statemachine.StateUnconfigured, statemachine.StateWatchdogAlert,
	statemachine.StateQueued, statemachine.StateStuck,
}

// TestBorderIsNeverAppliedInOneFrame is ticket 62 item 1 in one line: arriving
// in a state must only set the border TARGET, never the value - the owner's
// words are "不说话的时候…通过一个过度动画丝滑引入边框", and a border that shows
// up in a single frame is exactly what was rejected.
func TestBorderIsNeverAppliedInOneFrame(t *testing.T) {
	for _, s := range allKnownStates {
		m := &liquidMotion{}
		m.enterState(s)
		if m.border != 0 {
			t.Fatalf("enterState(%s) put the border at %v in one frame", s, m.border)
		}
		if m.wantBorder() != borderAtRest(s) {
			t.Fatalf("%s: target %v must be the state's resting level %v",
				s, m.wantBorder(), borderAtRest(s))
		}
	}
}

// TestBorderGlidesInAfterTheSessionEnds drives the state-driven half of the
// transition with no audio at all: an idle orb that the assistant has stopped
// talking about takes the border over BorderOpenMs of frames.
func TestBorderGlidesInAfterTheSessionEnds(t *testing.T) {
	const dt = 33 * time.Millisecond
	m := &liquidMotion{}
	m.enterState(statemachine.StateSpeaking)
	for i := 0; i < 40; i++ { // a long pause in a speaking state
		m.stepBorder(dt)
	}
	if m.border != 0 {
		t.Fatalf("Speaking must keep the border away, got %v", m.border)
	}
	m.endSession()
	m.enterState(statemachine.StateWarm)

	frames := rampBorderFrames(m, statemachine.StateWarm, dt)
	if len(frames) < 5 {
		t.Fatalf("the border must arrive through several frames, got %v", frames)
	}
	for _, f := range frames[:len(frames)-1] {
		if f <= 0 || f >= 1 {
			t.Fatalf("a travelling frame must sit strictly between, got %v", frames)
		}
	}
	if frames[len(frames)-1] != 1 {
		t.Fatalf("the border must end fully in, got %v", frames)
	}
	for i := 1; i < len(frames); i++ {
		if frames[i] <= frames[i-1] {
			t.Fatalf("the border must travel monotonically in, got %v", frames)
		}
	}
	// The idle orb must not swirl just because a ramp is travelling.
	if m.level != 0 || m.angle != 0 {
		t.Fatalf("stepBorder moved the liquid: level=%v angle=%v", m.level, m.angle)
	}
}

// TestSpeakingGlidesTheBorderBackOut covers the other direction: entering a
// speaking state takes the border OUT through a ramp, not by dropping it.
func TestSpeakingGlidesTheBorderBackOut(t *testing.T) {
	const dt = 33 * time.Millisecond
	m := &liquidMotion{}
	m.enterState(statemachine.StateListening)
	if got := rampBorderFrames(m, statemachine.StateListening, dt); len(got) == 0 || got[len(got)-1] != 1 {
		t.Fatalf("setup: Listening must carry a full border, got %v", got)
	}
	m.enterState(statemachine.StateSpeaking)
	if !m.busy(statemachine.StateSpeaking) {
		t.Fatal("a full border entering Speaking must ramp out, not vanish")
	}
	frames := rampBorderFrames(m, statemachine.StateSpeaking, dt)
	if len(frames) < 3 {
		t.Fatalf("the border must leave through several frames, got %v", frames)
	}
	if frames[0] <= 0 || frames[0] >= 1 {
		t.Fatalf("the first frame out must still be mid-travel, got %v", frames)
	}
	if frames[len(frames)-1] != 0 {
		t.Fatalf("Speaking must end with no border, got %v", frames)
	}
}

// rampBorderFrames runs the bounded transition to its end the way the UI-thread
// burst does and returns the border value of every frame it drew.
func rampBorderFrames(m *liquidMotion, s statemachine.State, dt time.Duration) []float32 {
	var out []float32
	for i := 0; m.busy(s) && i < 60; i++ {
		if !m.stepBorder(dt) {
			return out // the model says nothing moved: no frame would be drawn
		}
		out = append(out, m.border)
	}
	return out
}

// TestTransitionTimerSetIsTheFrozenTimerSet is the D32 belt: the bounded
// transition timer may only ever exist where ticket 07's frozen policy already
// granted a timer, and Sleeping is provably not one of them no matter what the
// motion fields say.
func TestTransitionTimerSetIsTheFrozenTimerSet(t *testing.T) {
	for _, s := range allKnownStates {
		want := Animated(s) || AnimationPolicy(s).PeriodMs > 0
		if got := transitionDriven(s); got != want {
			t.Errorf("%s: transitionDriven=%v but the frozen policy grants a timer=%v", s, got, want)
		}
		// The worst-case motion state (a burst in flight and a border still
		// travelling) must still refuse a timer wherever the frozen set does.
		m := &liquidMotion{active: true, summon: 1, border: 0.5, borderState: 1, ownsBorder: true}
		if m.busy(s) && !want {
			t.Errorf("%s: the liquid burst wants a timer in a state the frozen policy keeps static", s)
		}
	}
	m := &liquidMotion{active: true, summon: 1, border: 0.5, borderState: 1, ownsBorder: true}
	if m.busy(statemachine.StateSleeping) {
		t.Fatal("Sleeping must never host the transition timer (zero-timer discipline)")
	}
	if transitionDriven(statemachine.StateSleeping) {
		t.Fatal("Sleeping must never be transition-driven")
	}
	for _, s := range []statemachine.State{statemachine.StateConversation, statemachine.StateArmed,
		statemachine.StateMuted, statemachine.StateError, statemachine.StateDownloading} {
		if transitionDriven(s) {
			t.Errorf("%s must stay a static frame with no timer", s)
		}
	}
}

// TestSleepingFrameLosesTheBorderAuthority guards the committed glass frame:
// entering Sleeping must hand the border back to the frozen mapping (which
// carries none) so the Sleeping pixels stay exactly the shipped ones.
func TestSleepingFrameLosesTheBorderAuthority(t *testing.T) {
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	EnablePrototypeVisuals(true)

	m := &liquidMotion{}
	m.beginSession()
	m.enterState(statemachine.StateListening)
	stepMotion(m, 0, 40, 33*time.Millisecond) // idles into the border
	// The order Ball.motionStateChanged uses for a state that leaves the
	// session family: park the motion, then hand the border back.
	m.endSession()
	m.enterState(statemachine.StateSleeping)
	if m.ownsBorder || m.border != 0 || m.active {
		t.Fatalf("Sleeping must park the motion: owns=%v border=%v active=%v",
			m.ownsBorder, m.border, m.active)
	}
	v := VisualFor(pal, BallSizeDefaultPx, statemachine.StateSleeping, 0)
	if v.BorderAlpha != 0 || !v.Glass {
		t.Fatalf("the Sleeping frame must carry no border: %+v", v)
	}
	m.applyTo(&v)
	if v.BorderAlpha != 0 || v.LiquidLevel != 0 || v.SummonFlow != 0 || v.LiquidAngle != 0 {
		t.Fatalf("the parked motion must not touch the Sleeping frame: %+v", v)
	}
}

func TestLiquidDrivenStateSet(t *testing.T) {
	for _, s := range []statemachine.State{
		statemachine.StateListening, statemachine.StateSpeaking,
		statemachine.StateThinking, statemachine.StateConversation,
	} {
		if !liquidDriven(s) {
			t.Errorf("session state %s must be liquid-driven", s)
		}
	}
	// The zero-activity promise: nothing outside the session family moves.
	for _, s := range []statemachine.State{
		statemachine.StateSleeping, statemachine.StateArmed, statemachine.StateMuted,
		statemachine.StateAwaitingApproval, statemachine.StateError,
		statemachine.StateNoNetwork, statemachine.StateUnconfigured,
		statemachine.StateFirstRun, statemachine.StateWatchdogAlert,
		statemachine.StateQueued, statemachine.StateStuck, statemachine.StateDownloading,
	} {
		if liquidDriven(s) {
			t.Errorf("state %s must stay a static frame", s)
		}
		if AnimationPolicy(s).PeriodMs > 0 && s != statemachine.StateWarm {
			t.Errorf("state %s must not arm a looping timer", s)
		}
	}
}

func TestPushLevelIsInputHardened(t *testing.T) {
	m := &liquidMotion{active: true}
	m.pushLevel(-5, 33*time.Millisecond)
	m.pushLevel(9, -time.Second)
	if m.level < 0 || m.level > 1 {
		t.Fatalf("envelope must stay in 0..1, got %v", m.level)
	}
	if math.IsNaN(float64(m.angle)) || m.angle < 0 {
		t.Fatalf("angle must stay sane, got %v", m.angle)
	}
	// A long stall (a suspended process, a paused device) must not fling the
	// liquid around: dt is clamped inside the model to 250ms.
	before := m.angle
	m.pushLevel(1, time.Hour)
	maxTravel := (SpinBaseRadPerS + SpinLevelRadPerS + SpinSummonRadPerS) * 0.25
	if m.angle-before > maxTravel {
		t.Fatalf("a stalled feed must be clamped to 250ms: travelled %v rad, cap %v", m.angle-before, maxTravel)
	}
}

func TestRepaintIsSkippedWhenNothingChanged(t *testing.T) {
	m := &liquidMotion{active: true}
	m.pushLevel(0, 33*time.Millisecond)
	if m.pushLevel(0, time.Millisecond) {
		t.Fatal("a 1ms step with a settled envelope must not ask for a frame")
	}
	if !m.pushLevel(1, 33*time.Millisecond) {
		t.Fatal("a level jump must ask for a frame")
	}
}
