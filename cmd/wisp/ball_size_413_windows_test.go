//go:build windows

package main

// Ticket 413's rulers, both directions. The gap the ticket closes: ball.New
// takes a SizePx option (internal/ball/ball_windows.go:64, "configured orb
// size 44..160 (0 = default 56)") that production never fed, so the slider's
// number stopped at config.toml and every resident boot built the ball at the
// package default. The feed line is SizePx: hs.sizePx in startResidentBall,
// sourced by ballSizeSeed413 and handed over with withBallSizePx.
//
// The reverse criterion is the REFLECTION one (ball.Options.SizePx ==
// cfg.Ball.Size), and the reason it is legal: production carries no reader of
// the configured size back out of the ball - that absence is the gap, closed
// at the feed side instead of by growing an accessor - so the only honest
// instrument is reflect on the unexported opts field, which reads but never
// writes. The two live boots below follow the house shape
// (startResidentBallBounded, skip-loud on a headless host, rb.stop cleanup).

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// ballSize413Plant is a non-default, in-band size: 88 sits inside 44..160,
// differs from the package default 56, and is not an endpoint (44/160), so a
// clamp hiding in the feed cannot answer the same reading.
const ballSize413Plant = 88

// writeBallSize413 writes a config.toml whose [ball] section carries one size,
// the production shape (config.LoadFile fills every other section's schema
// defaults, exactly as it does for writeConfig258's [hotkey]-only body).
func writeBallSize413(t *testing.T, dir string) {
	t.Helper()
	body := "[ball]\nsize = 88\n"
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("write [ball]-only config: %v", err)
	}
}

// ballConfiguredSizePx413 reads the SizePx the ball actually holds off its
// unexported Options. Reading an unexported field through reflect is allowed;
// only Interface() is not - and nothing here writes, so the instrument cannot
// have fed the answer it reports.
func ballConfiguredSizePx413(t *testing.T, b *ball.Ball) int {
	t.Helper()
	if b == nil {
		t.Fatal("no ball handle: the reflection ruler has no subject")
	}
	v := reflect.ValueOf(b).Elem().FieldByName("opts").FieldByName("SizePx")
	if !v.IsValid() || v.Kind() != reflect.Int {
		t.Fatalf("ball.Ball's opts.SizePx is not an int field (valid=%v kind=%v): the option the feed rides has moved", v.IsValid(), v.Kind())
	}
	return int(v.Int())
}

// Test413SeedReadsTheFileIsTheCfgBallSize is the seed half: a config.toml with
// [ball] size=88 answers 88 from ballSizeSeed413, and that 88 IS the
// cfg.Ball.Size a fresh LoadFile of the same file answers - the seed adds
// nothing of its own.
func Test413SeedReadsTheFileIsTheCfgBallSize(t *testing.T) {
	dir := t.TempDir()
	writeBallSize413(t, dir)
	px, ok := ballSizeSeed413(dir)
	if !ok {
		t.Fatalf("ballSizeSeed413 said the file was unreadable: the fixture did not load")
	}
	c, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
	if err != nil || c == nil {
		t.Fatalf("control LoadFile failed: %v", err)
	}
	if px != c.Ball.Size {
		t.Errorf("seed = %d, cfg.Ball.Size = %d: the seed answered something the file did not say", px, c.Ball.Size)
	}
	if px != ballSize413Plant {
		t.Errorf("seed = %d, want the planted 88", px)
	}

	// The fallback half: a directory with no config.toml answers (0, false),
	// and 0 means ball.New's own default 56 stands - said by the seed's Warn,
	// not inferred.
	empty := t.TempDir()
	px, ok = ballSizeSeed413(empty)
	if ok || px != 0 {
		t.Errorf("seed over a config-less dir = (%d, %v), want (0, false)", px, ok)
	}
}

// Test413FeedReachesBallOptionsSizePx is the positive criterion: the full
// production wiring (startResidentBall over a live ball, the size handed in
// with withBallSizePx the way resident_windows.go hands it) puts the
// configured number into ball.Options.SizePx, and a reflection read of the
// ball's own opts answers cfg.Ball.Size.
func Test413FeedReachesBallOptionsSizePx(t *testing.T) {
	dir := t.TempDir()
	writeBallSize413(t, dir)
	px, ok := ballSizeSeed413(dir)
	if !ok || px != ballSize413Plant {
		t.Fatalf("seed = (%d, %v), want the planted 88: the fixture broke, not the feed", px, ok)
	}

	rb := startResidentBallBounded(t, observe.NewRegistry(), nil, nil, nil, withBallSizePx(px))
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the feed has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()

	got := ballConfiguredSizePx413(t, rb.b)
	if got != px {
		t.Errorf("ball.Options.SizePx = %d, want cfg.Ball.Size = %d: the configured size did not reach the option it names (the pre-413 gap is back)", got, px)
	}
	if got == ball.BallSizeDefaultPx {
		t.Errorf("a fed ball still holds the package default %d: the feed line is a no-op", ball.BallSizeDefaultPx)
	}
}

// Test413UnfedBallKeepsThePackageDefault is the reverse criterion: with no
// withBallSizePx hook (the pre-413 wiring), the same construction answers the
// package default 56 - proving the DIFFERENTIATION the two live boots make:
// not fed -> 56, fed -> the configured value, so the reading cannot be the
// same number by accident.
func Test413UnfedBallKeepsThePackageDefault(t *testing.T) {
	rb := startResidentBallBounded(t, observe.NewRegistry(), nil, nil, nil)
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the unfed construction has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()

	got := ballConfiguredSizePx413(t, rb.b)
	if got != ball.BallSizeDefaultPx {
		t.Errorf("unfed ball.Options.SizePx = %d, want the package default %d", got, ball.BallSizeDefaultPx)
	}
	if ballSize413Plant == ball.BallSizeDefaultPx {
		t.Fatalf("the plant equals the default: this test could not tell fed from unfed")
	}
}
