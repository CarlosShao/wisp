//go:build windows

package main

// Ticket 422's write-path teeth: config.set's ball_size field reaches the
// REAL guarded setter (internal/config's SetBallSize) through this store's
// ApplySetting, lands in memory hot (readable by BallSizeOptions in the same
// process), answers the honest restart tier (the orb window geometry is built
// once at boot from the seed read - there is no in-process re-size to claim),
// and every refusal carries the sentence the 413 contract pinned:
// out-of-band values answer validate.go's own ballSizeRangeError verbatim
// ("config.toml: ball.size 43 out of range [44, 160]"), a non-integer answers
// the parse sentence, and a refused write leaves the config byte-identical.
//
// The harness is ball_style_65_windows_test.go's ballStyleTestConfig65 (a
// minimal valid config.toml) and newConfigStore - the same real objects the
// 65 write leg proved its route with, no stub on the write path.

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// ballSize422InBand is AC#2's in-band, non-default, non-endpoint size: 90
// sits inside 44..160, is not the schema default 56, and is what the 413
// round-trip leg already proved travels from the file to the ball window.
const ballSize422InBand = 90

func Test422ApplySettingBallSizeWritesAndRefuses(t *testing.T) {
	mgr, err := config.NewManager(filepath.Join(ballStyleTestConfig65(t, ""), configFileName), nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	s := newConfigStore(mgr, nil, nil, "test-422")

	// The in-band write: lands, answers restart (no in-process applier
	// exists for size - the honest tier per the 413 contract §1), and the
	// receipt names the key path the FILE changed by.
	res, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallSize, Value: "90", Actor: "test-422"})
	if err != nil {
		t.Fatalf("ApplySetting(ball_size=90): %v", err)
	}
	if res.Tier != panel.EffectiveRestart {
		t.Fatalf("tier = %v, want restart (no in-process re-size exists; claiming now would be the lie the contract forbids)", res.Tier)
	}
	if len(res.Written) != 1 || res.Written[0] != "ball.size" {
		t.Fatalf("receipt = %v, want [ball.size]", res.Written)
	}
	if v := mgr.BallSizeOptions(); v.Size != ballSize422InBand {
		t.Fatalf("BallSizeOptions().Size = %d, want %d - the write must be readable hot in this same process", v.Size, ballSize422InBand)
	}

	// The value is trimmed before the parse: " 120 " is 120, not a refusal.
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallSize, Value: " 120 ", Actor: "test-422"}); err != nil {
		t.Fatalf("ApplySetting(ball_size=\" 120 \"): %v", err)
	}
	if v := mgr.BallSizeOptions(); v.Size != 120 {
		t.Fatalf("BallSizeOptions().Size = %d, want 120 after the trimmed write", v.Size)
	}

	// Out of band, low: refused with the validate.go sentence VERBATIM (the
	// 413 contract's error table), and nothing moved.
	_, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallSize, Value: "43", Actor: "test-422"})
	if err == nil {
		t.Fatalf("43 must be refused by SetBallSize's guard")
	}
	if !strings.Contains(err.Error(), "config.toml: ball.size 43 out of range [44, 160]") {
		t.Fatalf("the 43 refusal is not the single-sourced validate sentence, got: %v", err)
	}
	if v := mgr.BallSizeOptions(); v.Size != 120 {
		t.Fatalf("after the refused write Size = %d, want the untouched 120 - a refusal moves no byte", v.Size)
	}

	// Out of band, high: same sentence shape.
	_, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallSize, Value: "161", Actor: "test-422"})
	if err == nil {
		t.Fatalf("161 must be refused by SetBallSize's guard")
	}
	if !strings.Contains(err.Error(), "config.toml: ball.size 161 out of range [44, 160]") {
		t.Fatalf("the 161 refusal is not the single-sourced validate sentence, got: %v", err)
	}

	// Not an integer at all: the parse sentence names the field and the band,
	// and the setter is never reached.
	_, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallSize, Value: "big", Actor: "test-422"})
	if err == nil {
		t.Fatalf("a non-integer value must be refused before the setter")
	}
	if !strings.Contains(err.Error(), `字段 "ball_size" 要的是 44..160 的整数（px），收到的值解析不了`) {
		t.Fatalf("the parse refusal is not the contract's sentence, got: %v", err)
	}
	if v := mgr.BallSizeOptions(); v.Size != 120 {
		t.Fatalf("after the parse refusal Size = %d, want the untouched 120", v.Size)
	}
}
