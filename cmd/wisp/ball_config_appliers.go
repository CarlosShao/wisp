package main

// The runtime-apply seams the settings write path uses for the three [ball]
// appearance keys (style 65, opacity_idle, click_through). The wrappers and
// their applier variables are pure Go with no Windows dependency, but
// panel_config_store.go (no build constraint) calls the wrappers, so they
// moved here out of ball_style_65_windows.go and
// ball_idle_appearance_windows.go (424-fix-b) to keep the Linux type-check
// whole. On a process that hosts no ball the variables stay nil (the windows
// install* functions are the only writers) and the wrappers answer false -
// the same honest "no live switch" the write path already handled.
//
// The install* half (installBallStyleApplier65, installBallIdleAppearanceAppliers)
// stays in the windows files: it takes *ball.Ball and assigns the real
// appliers. The seed readers (ballStyleSeed65, ballIdleAppearanceSeed) stay
// with them.

// ballStyleApplier65 is the runtime-apply seam: the ball host installs it over
// Ball.ApplyStyle once the ball exists; the settings write path (ApplySetting's
// ball_style case) calls applyBallStyle65 with the written family. nil means
// this process hosts no ball and the write answers restart honestly.
var ballStyleApplier65 func(name string) bool

// ballIdleOpacityApplier is the runtime-apply seam for the idle opacity: the
// ball host installs it over Ball.SetIdleOpacity once the ball exists; the
// settings write path (ApplySetting's ball_opacity_idle case) calls
// applyBallIdleOpacity with the written value. nil means this process hosts
// no ball and the write answers restart honestly.
var ballIdleOpacityApplier func(v float64) bool

// ballClickThroughApplier is the runtime-apply seam for the idle click-through
// switch, the same shape as the opacity seam above.
var ballClickThroughApplier func(on bool) bool

// applyBallStyle65 applies one family to the live ball. False = no applier is
// installed or the name is outside the vocabulary; the caller (ApplySetting)
// answers EffectiveRestart in that case instead of claiming a live switch.
func applyBallStyle65(name string) bool {
	if ballStyleApplier65 == nil {
		return false
	}
	return ballStyleApplier65(name)
}

// applyBallIdleOpacity applies one idle-opacity value to the live ball.
// False = no applier is installed (this process hosts no ball); the caller
// (ApplySetting) answers EffectiveRestart in that case instead of claiming a
// live switch.
func applyBallIdleOpacity(v float64) bool {
	if ballIdleOpacityApplier == nil {
		return false
	}
	return ballIdleOpacityApplier(v)
}

// applyBallClickThrough applies the idle click-through switch to the live
// ball, the same measured shape as applyBallIdleOpacity.
func applyBallClickThrough(on bool) bool {
	if ballClickThroughApplier == nil {
		return false
	}
	return ballClickThroughApplier(on)
}
