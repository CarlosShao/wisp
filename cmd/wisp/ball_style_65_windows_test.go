//go:build windows

package main

// Ticket 65's wiring teeth: the boot seed reads [ball].style (the default is
// the SCHEMA's glass, never a local literal), and the settings write path
// flips the process-wide family through the installed applier.
//
// The BALL half of Ball.ApplyStyle (the STA orb page swap) is proven live by
// the ball package's winlive suite (TestBallLiveOrbJarvisStyleRendersParticles
// runs the ApplyStyle glass/jarvis round-trip against a real window); these
// tests prove the WIRING a unit can reach: seed -> SetStyle before the ball,
// write -> applyBallStyle65 -> ActiveStyle, and the honest tier per applier
// state. The applier installed here is ball.SetStyle - the real product
// function (the process-wide half of ApplyStyle, tokens.go:391), not a stub.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// ballStyleTestConfig65 writes a minimal valid config.toml (style line optional)
// and returns its directory. The section set mirrors internal/config's own
// matrixBase harness so validation sees a loadable file, not a half one.
func ballStyleTestConfig65(t *testing.T, styleLine string) string {
	t.Helper()
	dir := t.TempDir()
	body := "schema_version = 2\n\n[ball]\nsize = 56\n" + styleLine +
		"\n[app]\ntheme = \"dark\"\nlanguage = \"zh-CN\"\n" +
		"\n[voice]\nenabled = true\n\n[voice.asr]\nprovider = \"local-sherpa\"\nmodel = \"parakeet-zh\"\n" +
		"\n[fs]\nallowed_dirs = []\n"
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TestBallStyleSeed65DefaultGlass: a config without a style line seeds glass,
// and (name, true) is the FILE's answer - the mutation tooth is the second
// return value: delete the LoadFile read in ballStyleSeed65 and this test goes
// red on ok=false, because no local default exists to fake it.
func Test65SeedDefaultGlass(t *testing.T) {
	dir := ballStyleTestConfig65(t, "")
	name, ok := ballStyleSeed65(dir)
	if !ok {
		t.Fatalf("ballStyleSeed65 = ok=false for a loadable config - the read is load-bearing")
	}
	if name != "glass" {
		t.Fatalf("seed = %q, want glass (the schema default, not a literal here)", name)
	}
}

// TestBallStyleSeed65Jarvis: the file's own family reaches the seed verbatim.
func Test65SeedJarvis(t *testing.T) {
	dir := ballStyleTestConfig65(t, "style = \"jarvis\"")
	name, ok := ballStyleSeed65(dir)
	if !ok || name != "jarvis" {
		t.Fatalf("seed = (%q, %v), want (jarvis, true)", name, ok)
	}
}

// TestBallStyleSeed65Unreadable: an unreadable file answers ok=false (and the
// caller leaves the package default standing) - "fell back" is said, not
// silently swallowed.
func Test65SeedUnreadable(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte("this is not toml [[["), 0o600); err != nil {
		t.Fatal(err)
	}
	name, ok := ballStyleSeed65(dir)
	if ok {
		t.Fatalf("an unreadable config must answer ok=false, got (%q, true)", name)
	}
	if name != "" {
		t.Fatalf("an unreadable config must not answer a family, got %q", name)
	}
}

// TestApplySettingBallStyleFlipsActiveStyle: the write path's teeth - a real
// Manager write through the guarded setter, then the installed applier flips
// the process-wide family. jarvis lands, glass lands back, the tier is said
// per applier state, and a word outside the vocabulary is refused before any
// byte moves.
func Test65ApplySettingBallStyleFlipsActiveStyle(t *testing.T) {
	t.Cleanup(func() {
		ball.SetStyle("glass") // the package default; tests here must not leave the family flipped
		ballStyleApplier65 = nil
	})
	mgr, err := config.NewManager(filepath.Join(ballStyleTestConfig65(t, ""), configFileName), nil)
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	s := newConfigStore(mgr, nil, nil, "test-65")

	// No applier installed (a process hosting no ball): the write lands on disk
	// but the honest tier is restart, not now.
	res, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallStyle, Value: "jarvis", Actor: "test-65"})
	if err != nil {
		t.Fatalf("ApplySetting(jarvis) with no applier: %v", err)
	}
	if res.Tier != panel.EffectiveRestart {
		t.Fatalf("tier with no ball = %v, want restart (nothing to apply it to)", res.Tier)
	}
	if len(res.Written) == 0 || res.Written[0] != "ball.style" {
		t.Fatalf("receipt = %v, want [ball.style]", res.Written)
	}
	if ball.ActiveStyle() != ball.StyleGlass {
		t.Fatalf("no applier must mean no family change, got %v", ball.ActiveStyle())
	}

	// The real product function as the applier (ball.SetStyle, the process-wide
	// half of Ball.ApplyStyle): the write now takes effect NOW, and the family
	// really moved.
	ballStyleApplier65 = ball.SetStyle
	res, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallStyle, Value: "jarvis", Actor: "test-65"})
	if err != nil {
		t.Fatalf("ApplySetting(jarvis) with applier: %v", err)
	}
	if res.Tier != panel.EffectiveNow {
		t.Fatalf("tier with applier = %v, want now", res.Tier)
	}
	if ball.ActiveStyle() != ball.StyleJarvis {
		t.Fatalf("ActiveStyle after jarvis write = %v, want jarvis", ball.ActiveStyle())
	}

	// And back: glass lands too (the default is a selectable value, not just
	// the absence of one).
	res, err = s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallStyle, Value: "glass", Actor: "test-65"})
	if err != nil {
		t.Fatalf("ApplySetting(glass): %v", err)
	}
	if res.Tier != panel.EffectiveNow || ball.ActiveStyle() != ball.StyleGlass {
		t.Fatalf("glass round-trip: tier=%v active=%v, want (now, glass)", res.Tier, ball.ActiveStyle())
	}

	// The refusal: a word the schema does not carry never reaches the applier.
	if _, err := s.ApplySetting(t.Context(), panel.SettingWrite{Field: panel.FieldBallStyle, Value: "aurora", Actor: "test-65"}); err == nil {
		t.Fatalf("an unknown family must be refused")
	} else if !strings.Contains(err.Error(), "ball.style") {
		t.Fatalf("the refusal must name the key, got: %v", err)
	}
	if ball.ActiveStyle() != ball.StyleGlass {
		t.Fatalf("a refused write must leave the family standing, got %v", ball.ActiveStyle())
	}
}
