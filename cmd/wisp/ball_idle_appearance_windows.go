//go:build windows

package main

// The appearance wave's boot-side half for the two [ball] keys the config
// screen left locked: the [ball] opacity_idle / click_through seed the ball
// host reads at assembly, and the runtime-apply seam the settings write path
// uses - the same two halves ball_style_65_windows.go established for the
// style family and ball_size_413_windows.go for the size.
//
// WHY A FRESH READ INSTEAD OF THE MANAGER. The ball is built before the task
// pipeline's config.Manager exists (the 258-a1 census measured that order,
// and the style/size seeds live on the same constraint), so the seed reads
// config.toml the same per-use way ballSizeSeed413 reads [ball].size:
// config.LoadFile over the data dir.
//
// THE CARRIER, STATED. The seed answers (opacity, clickThrough, ok). With
// ok=false (unreadable file) the opacity handed over is nil - ball.Options'
// "no configured value reached this option" carrier, which keeps the ball
// package's own idle rendering standing while still passing the clickThrough
// half as false (the pre-wave behavior). "The file said 0.35" and "nobody
// could read a file" stay two different facts: the seed's Warn names the
// fallback. The POINTER is the schema band's own doing: 0 is a real
// configured opacity_idle (an invisible idle orb), so a plain float64 with a
// sentinel would swallow it - the lying-option shape - and would dim every
// Options literal built without the field; the pointer carries the file's
// own answer, 0 included.
//
// NO BAND AND NO DEFAULT OF ITS OWN. The range lives where it already owns:
// config.Validate refuses opacity_idle outside 0..1 (validate.go), and
// ball's idleVisualOpacity clamps defensively. This file carries the file's
// own answer or nil, never a third copy of either.

import (
	"log/slog"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
)

// ballIdleAppearanceSeed reads [ball].opacity_idle and [ball].click_through
// the way ballSizeSeed413 reads [ball].size. ok=true says the file loaded -
// the values are the file's own answer, defaults included, never a local
// default. ok=false says the file could not be read: opacity comes back as
// nil (the ball package's own idle rendering stands) and the Warn names it,
// so the fallback is said out loud, not inferred.
func ballIdleAppearanceSeed(dataDir string) (opacity *float64, clickThrough bool, ok bool) {
	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil || c == nil {
		slog.Warn("ball: [ball].opacity_idle/click_through source unreadable at boot; the ball package's own idle rendering and clickable idle ball stay",
			"path", cfgPath, "err", err)
		return nil, false, false
	}
	return &c.Ball.OpacityIdle, c.Ball.ClickThrough, true
}

// installBallIdleAppearanceAppliers wires both seams over the host's ball. A
// nil ball (the no-window boot) installs nothing.
func installBallIdleAppearanceAppliers(b *ball.Ball) {
	if b == nil {
		return
	}
	ballIdleOpacityApplier = func(v float64) bool { return b.SetIdleOpacity(v) }
	ballClickThroughApplier = func(on bool) bool { return b.SetClickThrough(on) }
}
