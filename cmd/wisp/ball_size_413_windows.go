//go:build windows

package main

// Ticket 413's boot-side half of the [ball].size feed line: the seed that reads
// the configured diameter and hands it to the ball host, closing the gap the
// reader roster measured (cmd/wisp/config_readers_255.go's [ball] row,
// ticket 419's adjudication): internal/ball/ball_windows.go:64 carries the
// SizePx option ("configured orb size 44..160 (0 = default 56)") and nothing
// production fed it, so the slider's number stopped at config.toml and the
// window was always built at the package default.
//
// WHY A FRESH READ INSTEAD OF THE MANAGER. The ball is built before the task
// pipeline's config.Manager exists (the 258-a1 census measured that order, and
// hotCfg258 lives on the same constraint), so the seed reads config.toml the
// same per-use way ballStyleSeed65 reads [ball].style and hotCfg258 reads
// [hotkey]: config.LoadFile over the data dir.
//
// NO LITERAL AND NO CLAMP OF ITS OWN. The range lives in the two places that
// already own it: config.Validate rejects ball.size outside 44..160
// (internal/config/validate.go, the schema's own pin), and ball.New clamps
// again (0 = the package default 56, below 44 = 44, above 160 = 160). This
// file therefore carries the file's own answer, or 0 for "the file could not
// be read and the ball package's own default stands" - the Warn names the
// fallback so "fell back" and "the file said 56" stay two different facts,
// the same honesty shape ballStyleSeed65 keeps for the style family.

import (
	"log/slog"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/config"
)

// ballSizeSeed413 reads [ball].size the way ballStyleSeed65 reads [ball].style
// and hotCfg258 reads [hotkey]. Returns (size, true) when the file loaded - the
// size is the file's own answer, 56 included, never a local default. (0, false)
// says the file could not be read and ball.New's own default (56) stays; the
// Warn names it so the fallback is said out loud, not inferred.
func ballSizeSeed413(dataDir string) (int, bool) {
	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil || c == nil {
		slog.Warn("ball: [ball].size source unreadable at boot; the ball package default (56) stays",
			"path", cfgPath, "err", err)
		return 0, false
	}
	return c.Ball.Size, true
}
