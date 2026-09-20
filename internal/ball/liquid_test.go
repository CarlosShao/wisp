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
