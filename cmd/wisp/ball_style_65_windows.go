//go:build windows

package main

// Ticket 65's boot-side half of the appearance style switch: the [ball].style
// seed the ball host reads at assembly, and the runtime-apply seam the settings
// write path uses (docs/evidence/s1/65-style-switch-config.md §2, ruling 甲).
//
// WHY A FRESH READ INSTEAD OF THE MANAGER. The ball is built before the task
// pipeline's config.Manager exists (the 258-a1 census measured that order, and
// hotCfg258 lives on the same constraint), so the seed reads config.toml the
// same per-use way hotCfg258 reads [hotkey]: config.LoadFile over the data dir.
// The DEFAULT IS NOT WRITTEN HERE: a config without a style line loads as
// glass by the schema default (internal/config's own pin, ball_style_65_test.go),
// and the ball package's process-wide default is glass too - this file carries
// no family literal at all, so "默认就是玻璃球" flows from the two defaults that
// already existed instead of a third copy.
//
// THE ORDER, and why the seed runs twice-shaped. Before startResidentBall the
// seed's answer feeds ball.SetStyle: orbWindow seeds its served style from
// ActiveStyle at CREATION (orb_windows.go:280), so a jarvis config never paints
// a glass first frame - the same order cmd/balldebug -style runs (main.go:130).
// After the ball exists, ApplyStyle re-syncs (idempotent when SetStyle already
// landed) and installs the runtime applier: the seam the settings write path
// calls so a panel style switch takes effect without a restart. The applier is
// process-global because the ball is (one orb per process, activeOrb's guard);
// tests restore it.

import (
	"log/slog"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
)

// ballStyleSeed65 reads [ball].style the way hotCfg258 reads [hotkey]. Returns
// (name, true) when the file loaded - the name is the file's own answer, glass
// included, never a local default. ("", false) says the file could not be
// read and the ball package's own glass default stays; the Warn names it so
// "fell back" and "the file said glass" stay two different facts.
func ballStyleSeed65(dataDir string) (string, bool) {
	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil || c == nil {
		slog.Warn("ball: [ball].style source unreadable at boot; the package default (glass) stays",
			"path", cfgPath, "err", err)
		return "", false
	}
	return c.Ball.Style, true
}

// installBallStyleApplier65 wires the seam over the host's ball. A nil ball
// (the no-window boot) installs nothing.
func installBallStyleApplier65(b *ball.Ball) {
	if b == nil {
		return
	}
	ballStyleApplier65 = func(name string) bool { return b.ApplyStyle(name) }
}
