//go:build windows

package ball

// Window-side glue for the audio-driven liquid (ticket 62 step 3). This is
// the file the killed agent left stranded: it wrote b.liq onto a *Ball that
// had no such field. Here it is against the real struct - the model from
// liquid.go is used verbatim (it owns no clock), the last raw envelope lives
// on the Ball (liqRaw), and the liquid is stamped onto the visual at draw
// time (frameVisual) so the base visual never accumulates a multiplication.
//
// Frame sources, and only these two:
//
//  1. an incoming audio level (Ball.SetAudioLevel) - the capture thread pushes
//     one envelope per 512-sample frame (~31/s at 16k), which is already the
//     30fps cap, so speaking costs zero extra timers;
//  2. a transition in flight (the summon burst / the border ramp), served by a
//     BOUNDED timer that kills itself the moment liquidMotion.busy() goes false
//     and that may only run in a state the frozen SPEC-08 animation whitelist
//     already grants a timer to.
//
// Sleeping appears in neither list: liquidDriven(StateSleeping) is false, so
// SetAudioLevel returns before it touches anything and no timer exists to arm.
// While the frozen SPEC-08 visuals are in force (prototypeVisuals off) the
// whole seam is inert: no motion, no timer, byte-identical to the pre-ticket
// render paths.

import (
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// SetAudioLevel feeds one capture-frame envelope (RMS of the last 512 samples,
// normalised 0..1) to the liquid. Safe from any goroutine: it posts to the UI
// thread. Levels arriving in a static state are dropped on the floor - they
// must not wake a Sleeping orb.
//
// This is the only render-side audio input in the project: no samples and no
// transcript text cross into the ball (C25 contamination surface).
func (b *Ball) SetAudioLevel(level float32) {
	if b.closed.Load() {
		return
	}
	b.sta.PostTask(func() { b.applyLevel(level) })
}

// applyLevel advances the liquid with one sample and repaints if a pixel can
// move. Runs on the UI thread.
func (b *Ball) applyLevel(level float32) {
	if !prototypeVisuals {
		return
	}
	b.liqRaw = level
	if !liquidDriven(b.curState) {
		b.liq.endSession()
		b.stopLiquidTimer()
		return
	}
	if !b.liq.active {
		b.liq.beginSession()
		b.syncLiquidTimer()
	}
	if b.liq.pushLevel(level, b.motionDt()) {
		b.renderFrame()
	}
	b.syncLiquidTimer()
}

// motionDt returns the real time since the previous motion step, which is the
// only clock the model ever sees (it owns no timer of its own). This is a
// sample-interval measurement for the physics step, not timeout logic.
func (b *Ball) motionDt() time.Duration {
	now := time.Now()
	dt := now.Sub(b.lastMotion)
	b.lastMotion = now
	if dt < 0 {
		return 0
	}
	return dt
}

// onLiquidTick is the bounded transition burst: it exists only to carry the
// summon flow and the border ramp to their end, then it retires itself.
func (b *Ball) onLiquidTick() {
	if !b.liq.busy(b.curState) {
		b.stopLiquidTimer()
		return
	}
	if b.liq.pushLevel(b.liqRaw, b.motionDt()) {
		b.renderFrame()
	}
	if !b.liq.busy(b.curState) {
		b.stopLiquidTimer()
	}
}

// syncLiquidTimer arms or retires the burst timer to match the motion. Called
// after every state change and every level push. busy() already refuses
// outside the SPEC-08 animated-state whitelist, so no static state (Sleeping
// above all) can ever see this timer.
func (b *Ball) syncLiquidTimer() {
	if b.liq.busy(b.curState) {
		if !b.liquidTimerActive {
			b.liquidTimerActive = true
			pSetTimer.Call(uintptr(b.hwnd), timerLiquidID, uintptr(FrameIntervalMs), 0)
		}
		return
	}
	b.stopLiquidTimer()
}

func (b *Ball) stopLiquidTimer() {
	if b.liquidTimerActive {
		pKillTimer.Call(uintptr(b.hwnd), timerLiquidID)
		b.liquidTimerActive = false
	}
}

// frameVisual returns the current visual with the live liquid values stamped
// on. The stamp is applied to a copy at draw time and only while the liquid
// actually participates (an active session or a ramp still travelling), so
// (a) applyTo's BorderAlpha multiplication never compounds into the base
// visual, and (b) static states render exactly the committed glass frame.
func (b *Ball) frameVisual() Visual {
	v := b.curVisual
	if b.liq.active || b.liq.border > 0 || b.liq.summon > 0 {
		b.liq.applyTo(&v)
	}
	return v
}

// motionStateChanged folds a state change into the liquid: a session entry
// fires the summon flow burst (the owner's "调出悬浮球…液态流动"), leaving one
// parks the liquid back to a static body.
func (b *Ball) motionStateChanged(from, to statemachine.State) {
	if !liquidDriven(to) || !prototypeVisuals {
		b.liq.endSession()
		b.stopLiquidTimer()
		return
	}
	if !liquidDriven(from) || from != to {
		b.liq.beginSession()
	}
	b.lastMotion = time.Now()
}
