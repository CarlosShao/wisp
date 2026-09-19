package ball

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TestTokenGoldenValues pins key C21 values against design/assets/tokens.css
// (dark + light). If tokens.css changes, this test and tokens.go change with
// it - the cross-check table is docs/evidence/s1/c21-native-tokens.md.
func TestTokenGoldenValues(t *testing.T) {
	dark := DarkPalette()
	light := LightPalette()

	type chk struct {
		name string
		got  Color
		want Color
	}
	checks := []chk{
		{"dark accent", dark.Accent, hex(0x86C2B9, 1)},
		{"dark danger", dark.Danger, hex(0xE07A70, 1)},
		{"dark warn", dark.Warn, hex(0xD9B26A, 1)},
		{"dark success", dark.Success, hex(0x7DB896, 1)},
		{"dark info", dark.Info, hex(0x8AAAD2, 1)},
		{"dark warm", dark.Warm, hex(0xD6B184, 1)},
		{"dark accent-line", dark.AccentLine, rgba(134, 194, 185, 0.36)},
		{"dark orb-hi", dark.OrbHi, rgba(255, 255, 255, 0.90)},
		{"dark orb-body", dark.OrbBody, rgba(255, 255, 255, 0.30)},
		{"dark ball-halo", dark.BallHalo, rgba(255, 255, 255, 0.10)},
		{"dark danger-line", dark.DangerLine, rgba(224, 122, 112, 0.38)},
		{"dark tint-danger-hi", dark.TintDangerHi, hex(0xF3C0BA, 1)},
		{"dark tint-neutral-mid", dark.TintNeutralMid, hex(0x6A6D75, 1)},
		{"dark on-solid", dark.OnSolid, hex(0xFFFFFF, 1)},
		{"light accent", light.Accent, hex(0x3E837A, 1)},
		{"light danger", light.Danger, hex(0xBC544C, 1)},
		{"light orb-rim", light.OrbRim, rgba(20, 24, 28, 0.16)},
		{"light glass-ring", light.GlassRing, rgba(255, 255, 255, 0.66)},
		// Light does not redefine on-solid/ball-halo: cascade keeps :root.
		{"light on-solid cascade", light.OnSolid, hex(0xFFFFFF, 1)},
		{"light ball-halo cascade", light.BallHalo, rgba(255, 255, 255, 0.10)},
	}
	for _, c := range checks {
		if !colorEq(c.got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, c.got, c.want)
		}
	}
}

// colorEq compares with a small tolerance (float32 token math).
func colorEq(a, b Color) bool {
	const eps = 1e-6
	return abs32(a.R-b.R) < eps && abs32(a.G-b.G) < eps && abs32(a.B-b.B) < eps && abs32(a.A-b.A) < eps
}

func abs32(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}

// TestPremultipliedDerivation guards the ULW bitmap color math.
func TestPremultipliedDerivation(t *testing.T) {
	c := rgba(100, 150, 200, 0.5)
	p := c.Premultiplied()
	if p.R != 100.0/255*0.5 || p.G != 150.0/255*0.5 || p.B != 200.0/255*0.5 || p.A != 0.5 {
		t.Errorf("premultiplied wrong: %v", p)
	}
	if w := c.WithAlpha(0.25); w.R != c.R || w.A != 0.25 {
		t.Errorf("WithAlpha mutated RGB: %v", w)
	}
}

// TestNoHardcodedColorsInBallPackage enforces the C21 rule: only tokens.go
// may contain color literals. Any rgba(/6-digit hex elsewhere in the package
// is a review-rejecting violation (SPEC-08 §2, ticket 07 constraint).
func TestNoHardcodedColorsInBallPackage(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil || len(files) == 0 {
		t.Fatalf("glob ball package: %v (%d files)", err, len(files))
	}
	badColor := regexp.MustCompile(`#[0-9A-Fa-f]{6}\b`)
	badRGBA := regexp.MustCompile(`\brgba\(`)
	for _, f := range files {
		base := filepath.Base(f)
		if base == "tokens.go" || strings.HasSuffix(base, "_test.go") {
			continue
		}
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			if badColor.MatchString(line) {
				t.Errorf("%s:%d: hardcoded hex color (C21 violation): %s", base, i+1, strings.TrimSpace(line))
			}
			if badRGBA.MatchString(line) {
				t.Errorf("%s:%d: hardcoded rgba literal (C21 violation): %s", base, i+1, strings.TrimSpace(line))
			}
		}
	}
}

// TestVisualForCoversAllTwentyStates: every frozen state has a visual and the
// SPEC-08 §2.1 anchor values hold.
func TestVisualForCoversAllTwentyStates(t *testing.T) {
	states := []statemachine.State{
		statemachine.StateFirstRun, statemachine.StateSleeping, statemachine.StateArmed,
		statemachine.StateMuted, statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateActing, statemachine.StateSpeaking, statemachine.StateWarm,
		statemachine.StateConversation, statemachine.StateConfirming,
		statemachine.StateAwaitingApproval, statemachine.StateSettling,
		statemachine.StateDownloading, statemachine.StateUnconfigured,
		statemachine.StateNoNetwork, statemachine.StateError, statemachine.StateQueued,
		statemachine.StateStuck, statemachine.StateWatchdogAlert,
	}
	if len(states) != 20 {
		t.Fatalf("vocabulary drifted: %d", len(states))
	}
	p := DarkPalette()
	for _, s := range states {
		v := VisualFor(p, BallSizeDefaultPx, s, 0)
		if v.SizePx <= 0 || v.Opacity <= 0 || v.Opacity > 1 {
			t.Errorf("%s: bad visual %+v", s, v)
		}
	}

	// SPEC-08 §2.1 anchors.
	if v := VisualFor(p, 56, statemachine.StateSleeping, 0); v.SizePx != 12 || v.Opacity != 0.35 {
		t.Errorf("Sleeping: want 12px @ 0.35, got %v px @ %v", v.SizePx, v.Opacity)
	}
	// Sleeping ignores the configured size (ball.html: 放大就失去意义).
	if v := VisualFor(p, 72, statemachine.StateSleeping, 0); v.SizePx != 12 {
		t.Errorf("Sleeping must stay 12px at any configured size, got %v", v.SizePx)
	}
	if v := VisualFor(p, 56, statemachine.StateArmed, 0); v.Opacity != 0.6 {
		t.Errorf("Armed opacity = %v, want 0.6", v.Opacity)
	}
	if v := VisualFor(p, 56, statemachine.StateMuted, 0); v.Opacity != 0.4 || v.Icon != IconSlash {
		t.Errorf("Muted: want 0.4 + slash, got %v / %d", v.Opacity, v.Icon)
	}
	// Conversation ring must not be weaker than Confirming's.
	conv := VisualFor(p, 56, statemachine.StateConversation, 0)
	conf := VisualFor(p, 56, statemachine.StateConfirming, 0)
	if conv.RingWidth <= conf.RingWidth {
		t.Errorf("Conversation ring (%v) must beat Confirming (%v)", conv.RingWidth, conf.RingWidth)
	}
	// Settling fade endpoints: 1 at entry, 0.35 fully faded (float32 eps).
	if v := VisualFor(p, 56, statemachine.StateSettling, 0); v.Opacity != 1 {
		t.Errorf("Settling entry opacity = %v, want 1", v.Opacity)
	}
	if v := VisualFor(p, 56, statemachine.StateSettling, 1); abs32(v.Opacity-SleepOpacity) > 1e-6 {
		t.Errorf("Settling end opacity = %v, want %v", v.Opacity, SleepOpacity)
	}
	// Size clamping to the token range.
	if v := VisualFor(p, 999, statemachine.StateWarm, 0); v.SizePx != BallSizeMaxPx {
		t.Errorf("size clamp high: %v", v.SizePx)
	}
	if v := VisualFor(p, 10, statemachine.StateWarm, 0); v.SizePx != BallSizeMinPx {
		t.Errorf("size clamp low: %v", v.SizePx)
	}
}

// TestAnimatedWhitelist pins the five loop-animated states.
func TestAnimatedWhitelist(t *testing.T) {
	yes := []statemachine.State{
		statemachine.StateListening, statemachine.StateThinking, statemachine.StateActing,
		statemachine.StateSpeaking, statemachine.StateConfirming,
	}
	for _, s := range yes {
		if !Animated(s) {
			t.Errorf("%s must be in the animated whitelist", s)
		}
	}
	no := []statemachine.State{
		statemachine.StateSleeping, statemachine.StateArmed, statemachine.StateMuted,
		statemachine.StateWarm, statemachine.StateSettling, statemachine.StateConversation,
		statemachine.StateAwaitingApproval, statemachine.StateQueued,
	}
	for _, s := range no {
		if Animated(s) {
			t.Errorf("%s must NOT be in the animated whitelist", s)
		}
	}
}

// TestAnimationPolicyZeroTimerInSleeping is the ticket 07 acceptance: the
// policy grants Sleeping no timer at all; every granted period respects the
// 30fps cap; Warm breathes via the alpha path; Settling fades one-shot.
func TestAnimationPolicyZeroTimerInSleeping(t *testing.T) {
	if p := AnimationPolicy(statemachine.StateSleeping); p.Kind != AnimNone || p.PeriodMs != 0 {
		t.Fatalf("Sleeping policy = %+v; must be zero (zero-timer discipline)", p)
	}

	all := []statemachine.State{
		statemachine.StateFirstRun, statemachine.StateSleeping, statemachine.StateArmed,
		statemachine.StateMuted, statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateActing, statemachine.StateSpeaking, statemachine.StateWarm,
		statemachine.StateConversation, statemachine.StateConfirming,
		statemachine.StateAwaitingApproval, statemachine.StateSettling,
		statemachine.StateDownloading, statemachine.StateUnconfigured,
		statemachine.StateNoNetwork, statemachine.StateError, statemachine.StateQueued,
		statemachine.StateStuck, statemachine.StateWatchdogAlert,
	}
	for _, s := range all {
		p := AnimationPolicy(s)
		if p.PeriodMs != 0 && p.PeriodMs < MinFrameMs {
			t.Errorf("%s: period %dms violates the 30fps cap (min %dms)", s, p.PeriodMs, MinFrameMs)
		}
	}
	if p := AnimationPolicy(statemachine.StateWarm); p.Kind != AnimAlpha || p.PeriodMs != 100 {
		t.Errorf("Warm policy = %+v; want alpha path at 10fps", p)
	}
	if p := AnimationPolicy(statemachine.StateSettling); p.Kind != AnimFade {
		t.Errorf("Settling policy = %+v; want the one-shot fade", p)
	}
	for _, s := range []statemachine.State{statemachine.StateListening,
		statemachine.StateThinking, statemachine.StateConfirming, statemachine.StateSpeaking} {
		if p := AnimationPolicy(s); p.Kind != AnimFrame || p.PeriodMs <= 0 {
			t.Errorf("%s policy = %+v; want an animation frame timer", s, p)
		}
	}
}

// TestHitTestAndDPIInjection covers click-through geometry and the DPI
// injection seam at the three common per-monitor V2 DPIs.
func TestHitTestAndDPIInjection(t *testing.T) {
	for _, dpi := range []uint32{96, 120, 144, 192} {
		edge := int32(WindowEdgePx(BallSizeDefaultPx, dpi))
		if edge <= 0 {
			t.Fatalf("dpi %d: window edge %d", dpi, edge)
		}
		// Center and mid-radius: clickable.
		c := edge / 2
		if got := HitTest(c, c, edge, BallSizeDefaultPx, dpi); got != HTClient {
			t.Errorf("dpi %d: center hit = %d, want HTClient", dpi, got)
		}
		if got := HitTest(c+int32(float64(BallSizeDefaultPx)*float64(dpi)/96/4), c, edge, BallSizeDefaultPx, dpi); got != HTClient {
			t.Errorf("dpi %d: mid-radius hit = %d, want HTClient", dpi, got)
		}
		// Window corner: transparent (click-through).
		if got := HitTest(0, 0, edge, BallSizeDefaultPx, dpi); got != HTTransparent {
			t.Errorf("dpi %d: corner hit = %d, want HTTransparent", dpi, got)
		}
		// Just outside the orb radius: transparent.
		if got := HitTest(c, 1, edge, BallSizeDefaultPx, dpi); got != HTTransparent {
			t.Errorf("dpi %d: top-edge hit = %d, want HTTransparent", dpi, got)
		}
	}

	// Edges grow with DPI and with the orb size.
	if WindowEdgePx(BallSizeDefaultPx, 144) <= WindowEdgePx(BallSizeDefaultPx, 96) {
		t.Error("window edge must grow with DPI")
	}
	if WindowEdgePx(BallSizeLargePx, 96) <= WindowEdgePx(BallSizeSmallPx, 96) {
		t.Error("window edge must grow with orb size")
	}
	if SleepWindowEdgePx(96) >= WindowEdgePx(BallSizeSmallPx, 96) {
		t.Error("Sleeping window must be smaller than the Armed orb window")
	}
}
