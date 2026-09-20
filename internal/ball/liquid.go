package ball

// Audio-driven liquid motion (ticket 62 step 3).
//
// The owner's requirement is that the liquid rotates *in rhythm with the
// voice*, not on its own clock: "说话的时候液体会随着声音有节奏旋转". So the
// driver of every frame here is a level update or a state transition - never a
// resident render loop. This file is deliberately pure logic with no clock and
// no Win32 calls: the caller supplies the elapsed time, which keeps the whole
// envelope/rotation model unit-testable and keeps Sleeping free of it.
//
// Signal path: the capture side (ticket 13) hands the ball an RMS/packet
// envelope via Ball.SetAudioLevel. Only that scalar crosses into the render
// layer - no audio samples, no transcript text (C25 contamination surface).

import (
	"math"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Spin rates (radians/s) at the ends of the envelope. A quiet session still
// creeps (the liquid is alive); a loud one whirls. Strictly monotonic in
// level, which is what AC#4 measures.
const (
	SpinBaseRadPerS    = 0.45 // level 0: slow drift while the session idles
	SpinLevelRadPerS   = 4.20 // added at level 1.0
	SpinSummonRadPerS  = 2.20 // added while the summon burst is in flight
	SilenceLevelGate   = 0.06 // envelope below this counts as "nobody speaking"
	SilenceHoldMs      = 150  // how long the voice must pause before the border starts
	LevelAttackTauMs   = 60   // envelope attack time constant (fast)
	LevelReleaseTauMs  = 420  // envelope release time constant (slow, musical)
	MotionEpsilon      = 1e-4 // below this, a quantity is considered settled
	FrameIntervalMs    = 1000 / MaxAnimFPS
	summonFullRotScale = 1.6
)

// liquidMotion is the whole simulation state of the liquid. It is owned by the
// ball window's UI thread (no locking) and is advanced exclusively by level
// updates and by the bounded burst timer.
type liquidMotion struct {
	level  float32       // smoothed envelope 0..1: swells/shrinks the liquid
	angle  float64       // accumulated rotation in radians (the renderer reads it)
	summon float32       // remaining summon flow burst 0..1 (1 = just summoned)
	border float32       // "not speaking" border fade 0..1 (1 = fully introduced)
	silent time.Duration // how long the raw feed has been under the gate
	active bool          // a session is running (Sleeping never sets this)
}

// liquidDriven reports whether a state receives audio-driven motion. Only
// session states do; Sleeping/Armed/Muted and the failure helpers are static
// glass frames by D32 discipline.
func liquidDriven(s statemachine.State) bool {
	switch s {
	case statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateSpeaking, statemachine.StateConversation,
		statemachine.StateWarm, statemachine.StateActing:
		return true
	}
	return false
}

// beginSession arms the summon flow burst: this is the "调出悬浮球" moment. The
// border starts away and fades in as soon as the voice stops (owner: "不说话
// 的时候…通过一个过度动画丝滑引入边框").
func (m *liquidMotion) beginSession() {
	m.active = true
	m.summon = 1
	m.border = 0
}

// endSession drops the liquid back to a static body (the state left the
// session family). Nothing here needs a frame: the caller repaints anyway.
func (m *liquidMotion) endSession() {
	m.active = false
	m.summon = 0
	m.angle = 0
	m.border = 0
	m.level = 0
	m.silent = 0
}

// pushLevel advances the simulation with one audio envelope sample. It returns
// true when the rendered image changed enough to be worth a repaint, which is
// how the ball avoids painting 30 identical frames a second while the mic
// hisses. dt is the caller's elapsed time since the previous sample.
func (m *liquidMotion) pushLevel(raw float32, dt time.Duration) bool {
	if raw < 0 {
		raw = 0
	}
	if raw > 1 {
		raw = 1
	}
	if dt < 0 {
		dt = 0
	}
	if dt > 250*time.Millisecond {
		dt = 250 * time.Millisecond // a stalled feed cannot teleport the liquid
	}
	before := m.snapshot()

	// Envelope follower: fast attack, slow release, so syllables read as
	// rhythm instead of as a square wave.
	tau := LevelReleaseTauMs
	if raw >= m.level {
		tau = LevelAttackTauMs
	}
	k := 1 - math.Exp(-float64(dt.Milliseconds())/float64(tau))
	m.level += float32(float64(raw-m.level) * k)
	if math.Abs(float64(m.level)) < MotionEpsilon {
		m.level = 0
	}

	// Has the voice actually stopped? The raw feed answers that (the smoothed
	// envelope is deliberately slow so syllables read as rhythm), with a short
	// hold so an in-word pause never pops the border in.
	if raw < SilenceLevelGate {
		m.silent += dt
	} else {
		m.silent = 0
	}

	// Rotation: rate is strictly increasing in (level, summon burst). The
	// angle accumulates, so a louder voice twists the liquid further per frame
	// and the swirl keeps its momentum between syllables.
	omega := SpinBaseRadPerS + SpinLevelRadPerS*float64(m.level) +
		SpinSummonRadPerS*float64(m.summon)
	m.angle += omega * dt.Seconds()
	if m.angle > 2*math.Pi*summonFullRotScale {
		m.angle = math.Mod(m.angle, 2*math.Pi)
	}

	// The summon burst is a one-shot decay over SummonFlowMs: it ends by
	// itself, which is what lets the burst timer retire.
	if m.summon > 0 {
		m.summon -= float32(dt) / float32(SummonFlowMs*time.Millisecond)
		if m.summon < 0 {
			m.summon = 0
		}
	}

	// Border transition, only once the voice has clearly stopped: the owner
	// wants the border "introduced smoothly" when nobody speaks, and wants the
	// liquid to own the orb while somebody does.
	target := float32(0)
	if m.silent >= SilenceHoldMs*time.Millisecond {
		target = 1
	}
	ramp := BorderCloseMs * time.Millisecond
	if target > m.border {
		ramp = BorderOpenMs * time.Millisecond
	}
	m.border = approach(m.border, target, float32(dt)/float32(ramp))
	if m.border > 1 {
		m.border = 1
	}
	if m.border < 0 {
		m.border = 0
	}

	return m.changed(before)
}

// approach moves cur toward tgt by at most step of the full travel (step is
// dt/rampDuration, so each transition lasts exactly its token duration).
func approach(cur, tgt, step float32) float32 {
	if step <= 0 {
		return cur
	}
	if step > 1 {
		step = 1
	}
	if tgt >= cur {
		if step >= tgt-cur {
			return tgt // land exactly: no float drift on a finished ramp
		}
		return cur + step
	}
	if step >= cur-tgt {
		return tgt
	}
	return cur - step
}

// motionSnapshot is the render-relevant slice of the state, used to decide
// whether a repaint would actually change any pixel.
type motionSnapshot struct {
	level, summon, border float32
	angle                 float64
}

func (m *liquidMotion) snapshot() motionSnapshot {
	return motionSnapshot{level: m.level, summon: m.summon, border: m.border, angle: m.angle}
}

func (m *liquidMotion) changed(before motionSnapshot) bool {
	if quantChanged(before.level, m.level) || quantChanged(before.summon, m.summon) ||
		quantChanged(before.border, m.border) {
		return true
	}
	// 0.004 rad is under a tenth of a pixel of blob travel at this radius.
	return math.Abs(before.angle-m.angle) > 4e-3
}

// quantChanged reports a change that survives the 1/255 alpha quantisation the
// renderer applies, i.e. one that can move a pixel.
func quantChanged(a, b float32) bool { return math.Abs(float64(a-b)) > 1.0/512 }

// busy reports whether the motion still has ramping left, which is the only
// reason a frame may be drawn without an incoming level sample. The state gate
// is the D32 belt: the burst only ticks where SPEC-08 §2 already grants an
// animation timer, so ticket 62 never arms a timer in a state that was static
// before it - Sleeping above all.
func (m *liquidMotion) busy(s statemachine.State) bool {
	if !m.active || !Animated(s) {
		return false
	}
	return m.summon > 0 || (m.border > MotionEpsilon && m.border < 1-MotionEpsilon)
}

// applyTo stamps the live motion onto a freshly mapped visual. It is a no-op
// while the frozen SPEC-08 §2.1 visuals are in force, so the contract tests
// keep asserting the frozen table.
func (m *liquidMotion) applyTo(v *Visual) {
	if !prototypeVisuals || v == nil {
		return
	}
	v.LiquidLevel = m.level
	v.LiquidAngle = m.angle
	v.SummonFlow = m.summon
	if v.Glass {
		// The border is a transition, not a constant: whatever the state
		// mapping asked for gets multiplied by how far the fade has travelled.
		v.BorderAlpha *= m.border
	}
}
