package ball

// Ticket 65 jarvis particle-orb tests (A1581 direction,
// docs/evidence/s1/65-style-direction.md §1.2/§2.5): token rows present,
// the per-state ring-radius law, the particle count band, layout
// determinism and the D32 static-frame pin.
//
// Same posture as style_65_test.go: every test here is pure (no window, no
// D2D, no timer, no desktop) - the jarvis mapping functions in jarvis.go are
// platform-free, so the whole surface is testable without a GUI.

import (
	"math"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TestJarvisTokenRowsArePresent: every colour row of the JarvisPalette token
// set is populated (a blank colour would draw invisible particles silently) -
// the cyan family (hue ~187deg), the thinking overrides and the violet
// variant (#9F63FF measured). Also pins the family selector vocabulary.
func TestJarvisTokenRowsArePresent(t *testing.T) {
	t.Cleanup(func() { SetJarvisFamily("cyan") })
	jt := JarvisTokens()
	rows := map[string]Color{
		"BgGlow":           jt.BgGlow,
		"Core":             jt.Core,
		"ParticleDim":      jt.ParticleDim,
		"ParticleBright":   jt.ParticleBright,
		"ParticleThinking": jt.ParticleThinking,
		"CoreThinking":     jt.CoreThinking,
		"VioletDim":        jt.VioletDim,
		"VioletBright":     jt.VioletBright,
		"VioletCore":       jt.VioletCore,
	}
	for name, c := range rows {
		if c.A <= 0 || c.A > 1 {
			t.Fatalf("JarvisPalette.%s carries a blank alpha: %+v", name, c)
		}
		if c.R == 0 && c.G == 0 && c.B == 0 {
			t.Fatalf("JarvisPalette.%s carries a blank colour", name)
		}
	}
	// The two families really are two hue families: the cyan bright particle
	// is green-dominant, the violet one is blue-dominant.
	if jt.ParticleBright.R >= jt.ParticleBright.G || jt.VioletBright.G >= jt.VioletBright.B {
		t.Fatalf("the cyan/violet families drifted: cyan %+v violet %+v", jt.ParticleBright, jt.VioletBright)
	}
	if !SetJarvisFamily("violet") || JarvisFamily() != "violet" {
		t.Fatal("SetJarvisFamily(violet) must select the violet variant")
	}
	if !SetJarvisFamily("CYAN") || JarvisFamily() != "cyan" {
		t.Fatal("SetJarvisFamily must be case-insensitive on the vocabulary")
	}
	if SetJarvisFamily("") || SetJarvisFamily("aurora") {
		t.Fatal("a name outside the cyan/violet vocabulary must be refused (aurora is a LOOK, not a jarvis family)")
	}
	if JarvisFamily() != "cyan" {
		t.Fatal("a refused SetJarvisFamily must keep the previous family")
	}
}

// TestJarvisPerStateRingRadii: the state structure law - thinking 0.67R,
// alert 0.39R, idle/success 0.50R - read through VisualFor so the mapping is
// pinned to the real state visuals, not hand-built structs.
func TestJarvisPerStateRingRadii(t *testing.T) {
	p := DarkPalette()
	cases := []struct {
		s    statemachine.State
		want float32
	}{
		{statemachine.StateThinking, jarvisRingMainThinkFrac},
		{statemachine.StateConfirming, jarvisRingMainAlertFrac},
		{statemachine.StateError, jarvisRingMainAlertFrac},
		{statemachine.StateSleeping, jarvisRingMainIdleFrac},
		{statemachine.StateListening, jarvisRingMainIdleFrac},
		{statemachine.StateSpeaking, jarvisRingMainIdleFrac}, // success keeps the idle structure (ember tile unreliable, §1.2)
		{statemachine.StateConversation, jarvisRingMainIdleFrac},
		{statemachine.StateWarm, jarvisRingMainIdleFrac},
		{statemachine.StateArmed, jarvisRingMainIdleFrac},
	}
	for _, c := range cases {
		if got := jarvisMainRingFrac(VisualFor(p, BallSizeDefaultPx, c.s, 0)); math.Abs(float64(got-c.want)) > 1e-6 {
			t.Fatalf("jarvisMainRingFrac(%s) = %v, want %v", c.s, got, c.want)
		}
	}
}

// TestJarvisParticleCountBand: 60-120 dots at the 44px floor, scaling with
// the orb size, and pure (same size, same count).
func TestJarvisParticleCountBand(t *testing.T) {
	n := jarvisParticleCount(BallSizeMinPx)
	if n < 60 || n > 120 {
		t.Fatalf("jarvisParticleCount(%vpx) = %d, want inside the measured 60-120 band", BallSizeMinPx, n)
	}
	if got := jarvisParticleCount(BallSizeMinPx); got != n {
		t.Fatalf("the particle budget must be pure: %d vs %d", got, n)
	}
	if jarvisParticleCount(BallSizeMaxPx) <= n {
		t.Fatalf("the budget must grow with size: %vpx gives %d, %vpx gives %d",
			BallSizeMaxPx, jarvisParticleCount(BallSizeMaxPx), BallSizeMinPx, n)
	}
	if jarvisParticleCount(BallSizeMaxPx) > 128 {
		t.Fatalf("the 72px budget = %d, above the reference grid's 128 ceiling", jarvisParticleCount(BallSizeMaxPx))
	}
}

// TestJarvisLayoutIsDeterministicAndBounded: the fixed-seed layout agrees bit
// for bit across bakes and every dot stays inside the orb with a unit wobble.
func TestJarvisLayoutIsDeterministicAndBounded(t *testing.T) {
	a := jarvisLayout(BallSizeDefaultPx, jarvisRingMainIdleFrac)
	b := jarvisLayout(BallSizeDefaultPx, jarvisRingMainIdleFrac)
	if len(a) != len(b) || len(a) == 0 {
		t.Fatalf("layout size drifted/empty: %d vs %d", len(a), len(b))
	}
	bound := jarvisFieldMaxFrac + jarvisRingJitterMain + float32(1e-6)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("dot %d drifted: %+v vs %+v (the seed must fully determine the field)", i, a[i], b[i])
		}
		if a[i].radius <= 0 || a[i].radius > bound {
			t.Fatalf("dot %d radius %v escapes the orb (bound %v)", i, a[i].radius, bound)
		}
		if a[i].wobble < -1 || a[i].wobble > 1 {
			t.Fatalf("dot %d wobble %v escapes the unit band", i, a[i].wobble)
		}
	}
	// A different main radius must move the main-orbit dots: the state
	// difference is structure, not palette.
	c := jarvisLayout(BallSizeDefaultPx, jarvisRingMainThinkFrac)
	moved := 0
	for i := range a {
		if a[i].ring == jarvisRingMain && a[i].radius != c[i].radius {
			moved++
		}
	}
	if moved == 0 {
		t.Fatal("the thinking main radius must move the main-orbit dots")
	}
}

// TestJarvisRestingFrameIsStatic mirrors TestWensityRestingFrameIsStatic for
// the jarvis style: the resting Sleeping frame carries zero spectrum, the
// identity pulse/scatter, zero ring spin and the idle ring radius - ONE
// deterministic static frame, zero timers, the particles do NOT spin at rest
// (D32).
func TestJarvisRestingFrameIsStatic(t *testing.T) {
	v := VisualFor(DarkPalette(), BallSizeDefaultPx, statemachine.StateSleeping, 0)
	if !v.Glass {
		t.Fatal("Sleeping must keep the glass body flag (the jarvis branch hangs off it)")
	}
	if v.SpectrumLow != 0 || v.SpectrumMid != 0 || v.SpectrumHigh != 0 {
		t.Fatalf("the resting frame must carry no spectrum: %+v", v)
	}
	if jarvisScatterAmp(v.LiquidLevel) != 0 {
		t.Fatal("the resting particles must not scatter")
	}
	if jarvisListenPulse(v.LiquidLevel) != 1 {
		t.Fatal("the resting main ring must hold its baked radius")
	}
	if jarvisRingSpin(v, 0.42) != 0 {
		t.Fatal("the resting particles must not spin (D32: zero timers at rest)")
	}
	if jarvisMainRingFrac(v) != jarvisRingMainIdleFrac {
		t.Fatal("the resting frame must use the idle ring radius")
	}
	// The static mapping never sets a spectrum for ANY state: bands arrive
	// only through the liquid seam (applyTo), never through VisualFor.
	for _, s := range allKnownStates {
		if vv := VisualFor(DarkPalette(), BallSizeDefaultPx, s, 0); vv.SpectrumLow != 0 || vv.SpectrumMid != 0 || vv.SpectrumHigh != 0 {
			t.Fatalf("VisualFor(%s) set a spectrum - bands must only come from the liquid seam", s)
		}
	}
}

// TestJarvisMotionOnlyInActiveStates: the scatter/pulse/spin laws return the
// identity for every static state's visual, and only the audio-driven /
// frame-driven states move.
func TestJarvisMotionOnlyInActiveStates(t *testing.T) {
	p := DarkPalette()
	for _, s := range allKnownStates {
		v := VisualFor(p, BallSizeDefaultPx, s, 0)
		if !Animated(s) {
			spin := jarvisRingSpin(v, 0.5)
			scatter := jarvisScatterAmp(v.LiquidLevel)
			pulse := jarvisListenPulse(v.LiquidLevel)
			if spin != 0 || scatter != 0 || pulse != 1 {
				t.Fatalf("%s is a static state but the jarvis mapping moves: spin=%v scatter=%v pulse=%v", s, spin, scatter, pulse)
			}
		}
	}
	// Listening: the envelope drives pulse + scatter.
	lv := VisualFor(p, BallSizeDefaultPx, statemachine.StateListening, 0)
	lv.LiquidLevel = 1
	if jarvisListenPulse(lv.LiquidLevel) == 1 || jarvisScatterAmp(lv.LiquidLevel) == 0 {
		t.Fatal("a full envelope must pulse and scatter the Listening ring")
	}
	// Thinking: the animation frame drives the ring spin.
	th := VisualFor(p, BallSizeDefaultPx, statemachine.StateThinking, 0)
	if jarvisRingSpin(th, 0.5) == 0 {
		t.Fatal("Thinking must advance the ring phase with the frame")
	}
}
