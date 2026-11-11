//go:build windows

package ball

// The idle-appearance pair (the appearance wave's native half): the two
// [ball] keys the config screen left locked because "那颗球是原生绘制的
// （D29），页面改不到" - now fed as options and applied where only native
// code can apply them:
//
//   - OpacityIdle: ball.opacity_idle 0..1, the Sleeping visual's opacity.
//     The resting look's SIZE and material stay exactly what the owner signed
//     (ticket 68/R15); only the opacity value the config names rides the
//     existing Visual.Opacity lane the renderer already multiplies every
//     drawn color by (renderer_windows.go's mulA(..., v.Opacity)) - no second
//     alpha mechanism, no new timer, no window-alpha write into the lane the
//     Warm breathing and Settling fade own (b.constAlpha stays theirs).
//   - ClickThrough: ball.click_through, WS_EX_TRANSPARENT on the window while
//     Sleeping. WS_EX_TRANSPARENT on a layered window is the one bit that
//     makes the WHOLE window (body included) pass mouse events through - the
//     per-pixel HTTRANSPARENT corners stay as they were, and the active
//     states remove the bit again so OnClickBall and the drag keep working
//     exactly as before.
//
// PURE CORES, ONE WRITER. The two mappings below are pure functions of
// (options, state, current value) so the tests can judge them without a
// window, and the two Ball methods (SetIdleOpacity / SetClickThrough) are the
// only runtime writers, posted to the ui-sta thread like every other surface
// in this package. The config-side feed lives in cmd/wisp (the seed and the
// settings-write applier seam, mirroring the 65/413 shape); internal/ball
// carries no config knowledge (SPEC-01 §3).

import (
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// idleVisualOpacity folds the configured idle opacity into a state's visual.
// The override applies to Sleeping ONLY and ONLY when the option carries a
// real value (nil = "no reading reached this option", and the current
// visual's own opacity stands; the pointer exists because 0 is a real
// configured value a zero-value field would swallow). The schema band is
// [0, 1] and both gates that judge it (the load gate and the write gate) hand
// out the same sentence, so this function clamps defensively instead of
// owning a second band.
func idleVisualOpacity(opts Options, s statemachine.State, current float32) float32 {
	if s != statemachine.StateSleeping {
		return current
	}
	if opts.OpacityIdle == nil {
		return current
	}
	v := float32(*opts.OpacityIdle)
	if v > 1 {
		v = 1
	}
	return v
}

// clickThroughExStyle answers the extended-style word the window should carry
// in this state: the click-through bit present in Sleeping when the option
// says so, absent in every other state and whenever the option says no - so
// the active ball stays clickable and a click_through=false boot is byte-for-
// byte the ex-style this package has always built.
func clickThroughExStyle(current uintptr, opts Options, s statemachine.State) uintptr {
	want := s == statemachine.StateSleeping && opts.ClickThrough
	has := current&wsExTransparent != 0
	switch {
	case want && !has:
		return current | wsExTransparent
	case !want && has:
		return current &^ wsExTransparent
	default:
		return current
	}
}

// visualForLocked maps a state to its visual with the idle-opacity override
// folded in - the ONE mapping site, so the boot frame, every state change and
// every repaint can never disagree about where the override ran.
func (b *Ball) visualForLocked(s statemachine.State) Visual {
	v := VisualFor(pal, b.opts.SizePx, s, 0)
	v.Opacity = idleVisualOpacity(b.opts, s, v.Opacity)
	return v
}

// applyClickThroughStyleLocked re-applies the click-through bit for the given
// state, reading the window's current extended style first and writing only
// on a real difference - a state change that does not move the answer must
// not churn Win32. STA thread only (GetWindowLongPtrW/SetWindowLongPtrW on
// the window this ball owns).
func (b *Ball) applyClickThroughStyleLocked(s statemachine.State) {
	if b == nil || b.hwnd == 0 {
		return
	}
	idxVal := int32(gwlExStyle)
	idx := uintptr(idxVal)
	cur, _, _ := pGetWindowLongPtrW.Call(uintptr(b.hwnd), idx)
	next := clickThroughExStyle(cur, b.opts, s)
	if next != cur {
		pSetWindowLongPtrW.Call(uintptr(b.hwnd), idx, next)
	}
}

// SetIdleOpacity is the runtime half of ball.opacity_idle: the settings write
// leg calls it when this process hosts the ball, and the idle look follows on
// the ui-sta thread without a restart. Safe from any goroutine (the write is
// posted like every other surface task); false = this ball hosts no window,
// and the caller answers restart honestly instead of claiming a live switch.
func (b *Ball) SetIdleOpacity(v float64) bool {
	if b == nil || b.hwnd == 0 {
		return false
	}
	b.sta.PostTask(func() {
		b.opts.OpacityIdle = &v
		if b.curState == statemachine.StateSleeping {
			b.repaintCurrent()
		}
	})
	return true
}

// SetClickThrough is the runtime half of ball.click_through, the same shape
// as SetIdleOpacity: the option is stored and the bit is re-applied for the
// CURRENT state, so a write while the ball is active lands the moment the
// ball next rests (and a write while it rests is visible immediately).
func (b *Ball) SetClickThrough(on bool) bool {
	if b == nil || b.hwnd == 0 {
		return false
	}
	b.sta.PostTask(func() {
		b.opts.ClickThrough = on
		b.applyClickThroughStyleLocked(b.curState)
	})
	return true
}
