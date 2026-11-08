package ball

import (
	"math"
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
		// light.GlassRing went with the retired v2 glass family (ticket 344
		// shape A, ledger A1395): no value is pinned here any more because the
		// field no longer exists; its ABSENCE is pinned by c21V2RetiredGlass in
		// tokens_table_test.go, which fails if anybody re-adds the field.
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

	// The flip itself (ticket 68 AC#2, 2026-11-01): the library default is the
	// look the owner signed, so every anchor below reads the DEFAULT mode and
	// asserts the prototype reading. Under each anchor, what it asserted
	// before the flip and where that number lives now:
	//   - Sleeping 12px @ 0.35 -> the resting glass body 34.72px @ 1
	//     (12px @ 0.35 is re-asserted verbatim in the frozen block below);
	//   - "12px at any configured size" -> "grows with the configured size"
	//     (SPEC-08 §2 INTERIM: 用户尺寸 × 0.62, floored - same re-homing);
	//   - Armed 0.6 / Muted 0.4 -> the prototype ladder 0.92 / 0.85 (R15#1
	//     Q-1: the glass material images exactly by being this opaque), with
	//     the 0.6/0.4 numbers re-asserted under the frozen flag;
	//   - Settling end 0.35 -> RestSettledOpacity 0.95 (R15#3 accepted the
	//     composite ≈0.33 as implemented: 0.95 here times the ULW 0.35),
	//     frozen 0.35 re-asserted below. Entry 1 held in both modes, kept.
	if !PrototypeVisualsEnabled() {
		t.Fatal("prototypeVisuals must default to ON after ticket 68 AC#2")
	}
	if v := VisualFor(p, 56, statemachine.StateSleeping, 0); abs32(v.SizePx-BallSizeDefaultPx*SleepRestRatio) > 1e-4 || v.Opacity != 1 || !v.Glass {
		t.Errorf("Sleeping: want the %vpx resting glass body @ 1, got %vpx @ %v glass=%v",
			BallSizeDefaultPx*SleepRestRatio, v.SizePx, v.Opacity, v.Glass)
	}
	// The resting body rides the configured size (floored; see the truth table).
	if v := VisualFor(p, 72, statemachine.StateSleeping, 0); v.SizePx <= BallSizeDefaultPx*SleepRestRatio {
		t.Errorf("Sleeping must grow past its default rest size at 72px configured, got %v", v.SizePx)
	}
	if v := VisualFor(p, 56, statemachine.StateArmed, 0); v.Opacity != 0.92 || !v.Glass {
		t.Errorf("Armed opacity = %v, want the prototype 0.92", v.Opacity)
	}
	if v := VisualFor(p, 56, statemachine.StateMuted, 0); v.Opacity != 0.85 || v.Icon != IconSlash {
		t.Errorf("Muted: want 0.85 + slash, got %v / %d", v.Opacity, v.Icon)
	}
	// Conversation ring must not be weaker than Confirming's.
	conv := VisualFor(p, 56, statemachine.StateConversation, 0)
	conf := VisualFor(p, 56, statemachine.StateConfirming, 0)
	if conv.RingWidth <= conf.RingWidth {
		t.Errorf("Conversation ring (%v) must beat Confirming (%v)", conv.RingWidth, conf.RingWidth)
	}
	// Settling fade endpoints: 1 at entry, RestSettledOpacity fully faded.
	if v := VisualFor(p, 56, statemachine.StateSettling, 0); v.Opacity != 1 {
		t.Errorf("Settling entry opacity = %v, want 1", v.Opacity)
	}
	if v := VisualFor(p, 56, statemachine.StateSettling, 1); abs32(v.Opacity-RestSettledOpacity) > 1e-6 {
		t.Errorf("Settling end opacity = %v, want %v", v.Opacity, RestSettledOpacity)
	}
	// The composite the user sees must land on the same number: the ULW
	// constant alpha is the fade's only driver (the frame underneath is the
	// fadeProgress=0 resting body), so settlingAlpha ends at
	// RestSettledOpacity in prototype mode (R15#3) - the SleepOpacity it
	// ended at before stacked a second fade on top and composited to ~0.33.
	// Monotonic, clamped, both ends pinned.
	if got := settlingAlpha(0); got != 1 {
		t.Errorf("settlingAlpha(0) = %v, want 1", got)
	}
	if got := settlingAlpha(1); abs32(got-RestSettledOpacity) > 1e-6 {
		t.Errorf("prototype settlingAlpha(1) = %v, want %v", got, RestSettledOpacity)
	}
	prev := float32(2)
	for i := 0; i <= 10; i++ {
		got := settlingAlpha(float64(i) / 10)
		if got > prev+1e-6 {
			t.Fatalf("settlingAlpha must not increase: p=%.1f gave %v after %v", float64(i)/10, got, prev)
		}
		prev = got
	}
	// Size clamping to the token range.
	if v := VisualFor(p, 999, statemachine.StateWarm, 0); v.SizePx != BallSizeMaxPx {
		t.Errorf("size clamp high: %v", v.SizePx)
	}
	if v := VisualFor(p, 10, statemachine.StateWarm, 0); v.SizePx != BallSizeMinPx {
		t.Errorf("size clamp low: %v", v.SizePx)
	}

	// The frozen §2.1 table is the escape hatch, and these are its pins: the
	// same anchors as before the flip, re-homed here verbatim (ticket 68 AC#2:
	// migrate every assertion, delete none, relax none - a test whose whole
	// purpose is the frozen table stays, made flag-explicit).
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	EnablePrototypeVisuals(false)
	if v := VisualFor(p, 56, statemachine.StateSleeping, 0); v.SizePx != 12 || v.Opacity != 0.35 {
		t.Errorf("frozen Sleeping: want 12px @ 0.35, got %v px @ %v", v.SizePx, v.Opacity)
	}
	// Sleeping ignores the configured size (ball.html: 放大就失去意义).
	if v := VisualFor(p, 72, statemachine.StateSleeping, 0); v.SizePx != 12 {
		t.Errorf("frozen Sleeping must stay 12px at any configured size, got %v", v.SizePx)
	}
	if v := VisualFor(p, 56, statemachine.StateArmed, 0); v.Opacity != 0.6 {
		t.Errorf("frozen Armed opacity = %v, want 0.6", v.Opacity)
	}
	if v := VisualFor(p, 56, statemachine.StateMuted, 0); v.Opacity != 0.4 || v.Icon != IconSlash {
		t.Errorf("frozen Muted: want 0.4 + slash, got %v / %d", v.Opacity, v.Icon)
	}
	if v := VisualFor(p, 56, statemachine.StateSettling, 1); abs32(v.Opacity-SleepOpacity) > 1e-6 {
		t.Errorf("frozen Settling end opacity = %v, want %v", v.Opacity, SleepOpacity)
	}
	// The frozen fade keeps the §2.1 composite endpoint verbatim (the escape
	// hatch must stay byte-identical to the shipped table).
	if got := settlingAlpha(1); abs32(got-SleepOpacity) > 1e-6 {
		t.Errorf("frozen settlingAlpha(1) = %v, want %v", got, SleepOpacity)
	}
}

// TestSleepingSizeTruthTable is ticket 68 AC#1's three-column size table, read
// off the pure mapping only (no window, no desktop, no -tags winlive):
//
//	① prototypeVisuals=false           stateSize -> SleepingDotPx        12px
//	② =true, free (undocked)           configuredPx*SleepRestRatio   34.72px
//	③ =true, docked, ramp landed       same as ②: the dock never touches size
//
// Column ③ is where the "it shrinks when docked" story actually lives, so the
// table also pins the two pure laws the dock draws with - DockSquash and
// DockPos - and records what the position law is FED: dockOrbRadius hands it
// the state's drawn body (R15#4/Q-4); the window radius it was fed before
// leaked ~82% of the rest body past the edge at the landed ramp.
func TestSleepingSizeTruthTable(t *testing.T) {
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	p := DarkPalette()

	// ① the frozen micro dot: 12px regardless of the configured size.
	EnablePrototypeVisuals(false)
	for _, cfg := range []int{BallSizeMinPx, BallSizeDefaultPx, BallSizeMaxPx} {
		if got := stateSize(cfg, statemachine.StateSleeping); got != SleepingDotPx {
			t.Errorf("① stateSize(%d, Sleeping) = %v, want the frozen %v", cfg, got, SleepingDotPx)
		}
		if v := VisualFor(p, cfg, statemachine.StateSleeping, 0); v.SizePx != SleepingDotPx || v.Opacity != SleepOpacity {
			t.Errorf("① VisualFor(%d, Sleeping) = %vpx @ %v, want %vpx @ %v",
				cfg, v.SizePx, v.Opacity, SleepingDotPx, SleepOpacity)
		}
	}

	// ② the resting glass body at the default size, and the floor that binds
	// for every configured size the token range allows below 49px.
	EnablePrototypeVisuals(true)
	const rest56 = float32(BallSizeDefaultPx) * SleepRestRatio // 34.72
	if got := stateSize(BallSizeDefaultPx, statemachine.StateSleeping); abs32(got-rest56) > 1e-4 {
		t.Errorf("② stateSize(56, Sleeping) = %v, want %v", got, rest56)
	}
	for _, cfg := range []int{44, 48} {
		if got := stateSize(cfg, statemachine.StateSleeping); abs32(got-SleepingRestMinPx) > 1e-4 {
			t.Errorf("② stateSize(%d, Sleeping) = %v, want the %vpx floor", cfg, got, SleepingRestMinPx)
		}
	}
	if got := stateSize(49, statemachine.StateSleeping); got <= SleepingRestMinPx {
		t.Errorf("② the floor must stop binding above 48px, got %v", got)
	}
	if v := VisualFor(p, BallSizeDefaultPx, statemachine.StateSleeping, 0); abs32(v.SizePx-rest56) > 1e-4 ||
		v.Opacity != 1 || !v.Glass {
		t.Errorf("② VisualFor(56, Sleeping) = %+v, want %vpx @ opacity 1, glass", v, rest56)
	}

	// ③ docking changes no size: stateSize has no dock input at all, so the
	// same call still answers 34.72px after the ramp lands.
	if got := stateSize(BallSizeDefaultPx, statemachine.StateSleeping); abs32(got-rest56) > 1e-4 {
		t.Errorf("③ the dock must not move stateSize, got %v", got)
	}
	if DockSquash(0) != 1 || abs32(DockSquash(1)-DockOverlapFrac) > 1e-6 {
		t.Errorf("③ DockSquash endpoints = %v/%v, want 1 and %v", DockSquash(0), DockSquash(1), DockOverlapFrac)
	}
	// What a landed tab leaves on screen, at the 96-DPI default. Since
	// R15#4/Q-4 the position law is fed the DRAWN body (dockOrbRadius), not
	// the window's radius: rest body 34.72px -> orbR 17, off =
	// int(17*DockOverlapFrac) = 7, so with the same 72px window injected the
	// tab keeps 36+7 = 43px of window on screen; and drawGlass squashes the
	// body along the axis with the same factor, so the tab shows
	// 2*17.36*0.42 ~ 14.6px of the 34.72px body = DockOverlapFrac - the 42%
	// the token name promises (the window radius kept a ~82%-visible tab).
	// The 47px row below is the pure position law pinned at an INJECTED
	// radius (28: the geometry the 62-diff-signoff "-dock-right" rows were
	// shot at); the call-site radius is pinned right after it.
	work := Rect{L: 0, T: 0, R: 1920, B: 1080}
	if p1 := DockPos(work, EdgeRight, 900, 500, 72, 28, 1); work.R-p1.X != 47 {
		t.Errorf("③ DockPos law at an injected radius: tab keeps %d px, want 47 (36 + int(28*%v))",
			work.R-p1.X, DockOverlapFrac)
	}
	callR := dockOrbRadius(rest56) // 17: what dockGeometry now hands DockPos
	if callR != 17 {
		t.Errorf("③ dockOrbRadius(34.72) = %d, want the drawn body's 17 (not the window's 28)", callR)
	}
	if p1 := DockPos(work, EdgeRight, 900, 500, 72, callR, 1); work.R-p1.X != 36+int(float32(callR)*DockOverlapFrac) {
		t.Errorf("③ a landed tab fed the state body keeps %d px on screen, want %d (36 + int(%d*%v))",
			work.R-p1.X, 36+int(float32(callR)*DockOverlapFrac), callR, DockOverlapFrac)
	}
}

// TestRecordedSleepingDiffBoxIsNotA44pxBody is the AC#1 arithmetic that decides
// which of the three columns docs/SLO.md A.2 actually measured, without
// re-measuring anything. The diff parent counts a pixel as imaged at >=8/255,
// and the widest thing a Sleeping frame draws is its halo: drawGlass fills an
// ellipse of radius 1.5*R while glowBrush only reaches zero alpha at 1.9*R, so
// the >=8 count is the area of a disc of radius 1.5*SizePx/2.
//
// Recorded for this exact configuration (default 56, 96 DPI, window edge 72,
// prototype on, undocked) by three independent runs: 2098 px in
// docs/evidence/s1/62-diff-glass/diff-table.txt, 2120 px in
// docs/evidence/s1/62-diff-signoff/diff-table.txt and 2103 px in docs/SLO.md
// A.2. Column ② (34.72px) predicts 2130. A 44px body predicts 3421, unclipped
// by that same window - so SPEC-08 §2's "直径 44px" cannot be the body of the
// frame that produced ANY of those rows. The 40x38, 44x44 and 46x46 the three
// runs report as the imaging box are the >=24/255 cut INSIDE that same halo,
// and that column is the unstable one: it moved 6px across runs that drew the
// same frame. The frozen dot (column ①, docs/evidence/s1/62-diff-baseline) has
// no box at all, because not one pixel moved >=8/255 (max delta 4/255).
func TestRecordedSleepingDiffBoxIsNotA44pxBody(t *testing.T) {
	// 1.5*R is drawGlass's halo fill; 34.72/2 is column ②'s radius.
	haloPx := func(sizePx float64) float64 {
		r := 1.5 * sizePx / 2
		return math.Pi * r * r
	}
	// Every recorded px_delta_ge8 for a prototype-mode, undocked Sleeping orb.
	for _, measured := range []float64{2098, 2103, 2120} {
		if got := haloPx(BallSizeDefaultPx * SleepRestRatio); math.Abs(got-measured)/measured > 0.02 {
			t.Errorf("column ② predicts %.0f imaged px, want within 2%% of the recorded %.0f", got, measured)
		}
		if got := haloPx(44); got < measured*1.5 {
			t.Errorf("a 44px body must predict far MORE imaged px than the evidence recorded, got %.0f vs %.0f", got, measured)
		}
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
	for _, s := range []statemachine.State{
		statemachine.StateListening,
		statemachine.StateThinking, statemachine.StateConfirming, statemachine.StateSpeaking,
	} {
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

	// R15#5 (owner 2026-11-01): the clickable circle reaches the visible outer
	// edge. The Sleeping rest body draws 34.72px (halo to 26px) but was
	// clickable only to body/2+4 = 21px; the WM_NCHITTEST call site hands
	// HitTest the int-truncated 34, so the circle now lands on ceil(34*0.75)
	// = 26 with the old body+tolerance circle as the floor.
	edge96 := int32(WindowEdgePx(BallSizeDefaultPx, 96))
	c96 := edge96 / 2
	restBodyPx := float32(BallSizeDefaultPx) * SleepRestRatio // 34.72
	restBody := int(restBodyPx)                               // 34, as the call site truncates it
	restHit := int(math.Ceil(float64(restBody) * 0.75))       // 26
	if got := HitTest(c96+int32(restHit)-1, c96, edge96, restBody, 96); got != HTClient {
		t.Errorf("the visible fringe of the resting Sleeping body (r=%d) must be clickable, got %d", restHit-1, got)
	}
	if got := HitTest(c96+int32(restHit), c96, edge96, restBody, 96); got != HTClient {
		t.Errorf("the visible edge itself (r=%d) counts as the ball, got %d", restHit, got)
	}
	if got := HitTest(c96+int32(restHit)+1, c96, edge96, restBody, 96); got != HTTransparent {
		t.Errorf("outside the visible edge (r=%d) must stay click-through, got %d", restHit+1, got)
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

// TestBallSizeMaxWidenedTo160 pins the A1636 widening (owner authorization,
// 2026-11-08): the config ceiling is 160 while the --ball-lg tier stays at 72 -
// the tier and the ceiling are different things, and a future edit that moves
// them together is wrong. The window edge follows the new ceiling exactly:
// orb + 2*RingMarginPx = 160 + 44 = 204 physical px at 96 DPI.
func TestBallSizeMaxWidenedTo160(t *testing.T) {
	if BallSizeMaxPx != 160 {
		t.Errorf("BallSizeMaxPx = %d, want the A1636 ceiling 160", BallSizeMaxPx)
	}
	if BallSizeLargePx != 72 {
		t.Errorf("BallSizeLargePx = %d, want the --ball-lg tier 72 unchanged (the tier is not the ceiling)", BallSizeLargePx)
	}
	if BallSizeDefaultPx != 56 || BallSizeSmallPx != 44 || BallSizeMinPx != 44 {
		t.Errorf("default/small/min drifted: %d/%d/%d, want 56/44/44", BallSizeDefaultPx, BallSizeSmallPx, BallSizeMinPx)
	}
	if got := WindowEdgePx(BallSizeMaxPx, 96); got != 204 {
		t.Errorf("WindowEdgePx(160, 96) = %d, want 204 (160 + 2*22)", got)
	}
	// The clamp tracks the ceiling: a configured 160px orb must survive both
	// clamp sites instead of being cut back to the old 72.
	if got := stateSize(BallSizeMaxPx, statemachine.StateWarm); got != BallSizeMaxPx {
		t.Errorf("stateSize(160, Warm) = %v, want the ceiling itself", got)
	}
}
