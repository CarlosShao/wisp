//go:build windows

package main

// Ticket 413 AC#4's full round trip, end to end over the production surfaces:
// a guarded settings write of 90 lands in config.toml, the read face answers
// 90 back, the boot seed (ballSizeSeed413 - what the resident ball is actually
// built from) reads 90 off the same file, and a live boot with the production
// hook holds 90 in ball.Options.SizePx (the same reflection instrument the
// AC#1 leg pinned). The refusal half: 43 and 161 are refused with validate.go's
// own sentence verbatim, the file stays byte-identical, and the seed still
// reads the value that was there before.
//
// The Go-side write/read faces live in internal/config (settings_413.go:
// SetBallSize / BallSizeOptions); this file proves they, the seed and the ball
// option agree on one number.

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// ballSize413RoundTrip is AC#4's in-band, non-default, non-endpoint size: 90
// sits inside 44..160, differs from the default 56, so neither a default nor a
// clamp can answer the same reading.
const ballSize413RoundTrip = 90

// writeBallSizeRoundTrip413 plants a config.toml whose [ball] section carries
// one size - the shape every 413 fixture uses.
func writeBallSizeRoundTrip413(t *testing.T, dir string, size int) {
	t.Helper()
	body := "[ball]\nsize = " + strconv.Itoa(size) + "\n"
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("write [ball]-only config: %v", err)
	}
}

// Test413RoundTripSettingsWrite90ReachesTheBallWindow: write 90 through the
// guarded setter, read 90 back through the reader, seed 90 off the file, and
// the live ball holds 90.
func Test413RoundTripSettingsWrite90ReachesTheBallWindow(t *testing.T) {
	dir := t.TempDir()
	writeBallSizeRoundTrip413(t, dir, 56)
	m, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	written, err := m.SetBallSize(ballSize413RoundTrip)
	if err != nil {
		t.Fatalf("SetBallSize(%d): %v", ballSize413RoundTrip, err)
	}
	if len(written) == 0 || written[0] != "ball.size" {
		t.Fatalf("receipt = %v, want [ball.size]", written)
	}
	if v := m.BallSizeOptions(); v.Size != ballSize413RoundTrip {
		t.Fatalf("BallSizeOptions().Size = %d, want %d (the read-back half)", v.Size, ballSize413RoundTrip)
	}

	px, ok := ballSizeSeed413(dir)
	if !ok || px != ballSize413RoundTrip {
		t.Fatalf("ballSizeSeed413 = (%d, %v), want (%d, true): the seed and the setter disagree about the same file", px, ok, ballSize413RoundTrip)
	}

	rb := startResidentBallBounded(t, observe.NewRegistry(), nil, nil, nil, withBallSizePx(px))
	if rb == nil || rb.b == nil {
		t.Skipf("SKIP-LOUD: this host cannot bring up a ball window (headless service session), so the fed size has no Win32 to land in: %q", rb.verdictStatusLineFor258())
	}
	defer rb.stop()

	got := ballConfiguredSizePx413(t, rb.b)
	if got != px {
		t.Errorf("ball.Options.SizePx = %d, want the written %d: the round trip broke between the file and the window", got, px)
	}
}

// Test413RoundTripRefuses43And161WithValidateSentence: the out-of-band writes
// are refused with internal/config/validate.go's own sentence verbatim, the
// file stays byte-identical, and the seed still reads the pre-refusal value.
func Test413RoundTripRefuses43And161WithValidateSentence(t *testing.T) {
	dir := t.TempDir()
	writeBallSizeRoundTrip413(t, dir, ballSize413RoundTrip)
	m, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	before, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatal(err)
	}

	for _, bad := range []int{43, 161} {
		written, err := m.SetBallSize(bad)
		if err == nil {
			t.Fatalf("SetBallSize(%d) must be refused", bad)
		}
		// AC#4 names the sentence: the load gate's own words
		// (internal/config/validate.go, ballSizeRangeError - the same helper
		// the write guard shares), verbatim, value spelled inside it.
		want := "config.toml: ball.size " + strconv.Itoa(bad) + " out of range [44, 160]"
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("SetBallSize(%d) refusal must carry validate.go's sentence %q, got: %v", bad, want, err)
		}
		if written != nil {
			t.Fatalf("a refused write must return no receipt, got %v", written)
		}
	}

	after, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("a refused write must leave config.toml byte-identical")
	}

	px, ok := ballSizeSeed413(dir)
	if !ok || px != ballSize413RoundTrip {
		t.Fatalf("after the refusals the seed = (%d, %v), want the untouched %d", px, ok, ballSize413RoundTrip)
	}
}
