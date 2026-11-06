package ball

// Ticket 65 (A1581 direction, docs/evidence/s1/65-style-direction.md):
// style switching, the Wensity band token/geometry mapping, the
// pseudo-spectrum derivation and the D32 static-frame pins.
//
// Every test here is pure (no window, no D2D, no timer, no desktop): the
// renderer consumes Visual + tokens only, so the whole surface is testable
// from the mapping functions - which is also what keeps the desktop gates
// (diff/timers/CPU) owner-side without a green-washing gap here.

import (
	"math"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TestStyleSelectionPinsTheConfigVocabulary: SetStyle routes exactly the
// [ball] style vocabulary, rejects the rest case-insensitively without
// losing the previous family, and the default is glass.
func TestStyleSelectionPinsTheConfigVocabulary(t *testing.T) {
	// The style family is a process-wide global (like activeLook): save and
	// restore it so the test order can never leak a jarvis selection into a
	// later run of this file (count=2 runs it twice).
	prev := ActiveStyle()
	t.Cleanup(func() {
		if prev == StyleGlass {
			SetStyle("glass")
		} else {
			SetStyle("jarvis")
		}
	})
	if got := ActiveStyle(); got != StyleGlass {
		t.Fatalf("the default style = %v, want StyleGlass", got)
	}
	names := StyleNames()
	if len(names) != 2 || names[0] != "glass" || names[1] != "jarvis" {
		t.Fatalf("StyleNames() = %v, want [glass jarvis]", names)
	}
	if !SetStyle("jarvis") || ActiveStyle() != StyleJarvis {
		t.Fatalf("SetStyle(jarvis) must select the jarvis family, got %v", ActiveStyle())
	}
	if !SetStyle("GLASS") || ActiveStyle() != StyleGlass {
		t.Fatalf("SetStyle is case-insensitive on the vocabulary, got %v", ActiveStyle())
	}
	if SetStyle("aurora") || SetStyle("") || SetStyle("glass ") {
		t.Fatalf("a name outside the vocabulary must be refused (aurora is a LOOK, not a STYLE)")
	}
	if ActiveStyle() != StyleGlass {
		t.Fatalf("a refused SetStyle must keep the previous family, got %v", ActiveStyle())
	}
}

// TestStyleSwitchNeverWakesSleeping is the D32 pin on the new seam: a style
// change is one repaint of a static frame - it must not move the animation
// policy of any state, and Sleeping stays zero-timer and untouched by the
// audio feed through the whole switch.
func TestStyleSwitchNeverWakesSleeping(t *testing.T) {
	prev := ActiveStyle()
	t.Cleanup(func() {
		if prev == StyleGlass {
			SetStyle("glass")
		} else {
			SetStyle("jarvis")
		}
	})
	snapshot := map[statemachine.State]AnimPolicy{}
	for _, s := range allKnownStates {
		snapshot[s] = AnimationPolicy(s)
	}
	for _, name := range StyleNames() {
		if !SetStyle(name) {
			t.Fatalf("SetStyle(%q) must succeed", name)
		}
		for _, s := range allKnownStates {
			if AnimationPolicy(s) != snapshot[s] {
				t.Fatalf("style %q moved AnimationPolicy(%s): %+v -> %+v", name, s, snapshot[s], AnimationPolicy(s))
			}
		}
	}
	if p := AnimationPolicy(statemachine.StateSleeping); p.Kind != AnimNone || p.PeriodMs != 0 {
		t.Fatalf("Sleeping policy after style switching = %+v, want zero (D32)", p)
	}
	if liquidDriven(statemachine.StateSleeping) || transitionDriven(statemachine.StateSleeping) {
		t.Fatal("Sleeping must stay outside both motion whitelists regardless of style")
	}
}

// TestWensityLookIsTheGlassDefault pins the A1581 direction: the wensity
// band treatment is DefaultLook (looks[0]), the old four families stay
// selectable as glass sub-treatments, and the band token set is fully
// populated (a zero colour would draw an invisible band silently).
func TestWensityLookIsTheGlassDefault(t *testing.T) {
	def := DefaultLook()
	if def.Name != "wensity" {
		t.Fatalf("DefaultLook() = %q, want wensity (ticket 65)", def.Name)
	}
	if !wensityActive() {
		t.Fatal("the default look must activate the wensity band branch")
	}
	if len(looks) != 5 {
		t.Fatalf("looks carries %d treatments, want 5 (wensity + the four ticket-62 families)", len(looks))
	}
	if l, ok := LookByName("aurora"); !ok || l.Name != "aurora" {
		t.Fatal("the four ticket-62 sub-treatments must stay reachable via -look")
	}
	if !SetLook("aurora") || wensityActive() {
		t.Fatal("selecting a blob look must deactivate the wensity band branch")
	}
	if !SetLook("wensity") || !wensityActive() {
		t.Fatal("selecting wensity by name must re-activate the band branch")
	}
	wb := WensityBandTokens()
	for _, c := range []Color{wb.StopRed, wb.StopMagenta, wb.StopWhite, wb.StopBlue, wb.StopCyan, wb.StopOrange, wb.SpecArc, wb.RimLip, wb.BottomGlow} {
		if c.A <= 0 || c.A > 1 || (c.R == 0 && c.G == 0 && c.B == 0) {
			t.Fatalf("the band token set carries a blank colour: %+v", wb)
		}
	}
}

// TestWensityBandGeometryPinsTheQuantifiedSpec: the spine peaks 0.13R above
// the centre, its ends differ by the 0.33R cold-end rise, the span reaches
// 0.84R/0.68R, and every segment paints one of the six pure stop colours
// (no interpolation inventing colours the table never declared).
func TestWensityBandGeometryPinsTheQuantifiedSpec(t *testing.T) {
	if got := wensityBandAxisY(0.5); math.Abs(float64(got-WensityBandCenterFrac)) > 1e-6 {
		t.Fatalf("band spine at the centre = %v, want %v above the centre", got, WensityBandCenterFrac)
	}
	rise := wensityBandAxisY(1) - wensityBandAxisY(0)
	if math.Abs(float64(rise-WensityBandTiltFrac)) > 1e-6 {
		t.Fatalf("cold-end rise = %v, want the measured %v", rise, WensityBandTiltFrac)
	}
	if got := wensityBandAxisX(0); math.Abs(float64(got+WensityBandSpanLeftFrac)) > 1e-6 {
		t.Fatalf("warm end reach = %v, want -%v", got, WensityBandSpanLeftFrac)
	}
	if got := wensityBandAxisX(1); math.Abs(float64(got-WensityBandSpanRightFrac)) > 1e-6 {
		t.Fatalf("cold end reach = %v, want %v", got, WensityBandSpanRightFrac)
	}
	stops := map[Color]bool{
		WensityBandTokens().StopRed: true, WensityBandTokens().StopMagenta: true,
		WensityBandTokens().StopWhite: true, WensityBandTokens().StopBlue: true,
		WensityBandTokens().StopCyan: true, WensityBandTokens().StopOrange: true,
	}
	seen := map[Color]bool{}
	for i := 0; i <= 32; i++ {
		col := wensityBandSegmentColor(float32(i) / 32)
		if !stops[col] {
			t.Fatalf("t=%.3f produced a colour no table row declares: %v", float32(i)/32, col)
		}
		seen[col] = true
	}
	if len(seen) != 6 {
		t.Fatalf("only %d of the 6 stops are reachable along the span: %v", len(seen), seen)
	}
	// The stop order reads red -> magenta -> white -> blue -> cyan -> orange.
	order := []Color{}
	var last Color
	for i := 0; i <= 64; i++ {
		if col := wensityBandSegmentColor(float32(i) / 64); col != last {
			order = append(order, col)
			last = col
		}
	}
	want := []Color{
		WensityBandTokens().StopRed, WensityBandTokens().StopMagenta,
		WensityBandTokens().StopWhite, WensityBandTokens().StopBlue,
		WensityBandTokens().StopCyan, WensityBandTokens().StopOrange,
	}
	if len(order) != len(want) {
		t.Fatalf("the band reads %d colour groups, want the 6 measured stops in order", len(order))
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("band stop order drifted at %d: got %v, want %v", i, order[i], want[i])
		}
	}
}

// TestWensityBandDeformFullScalePP: the audio deform reaches exactly the
// token's peak-to-peak travel at a full envelope, is monotonic in the
// spectrum, and clamps - so a runaway feed can never fling the band.
func TestWensityBandDeformFullScalePP(t *testing.T) {
	if got := wensityBandDeform(0); got != 0 {
		t.Fatalf("a silent band must not deform, got %v", got)
	}
	if got := wensityBandDeform(1); math.Abs(float64(got-WensityBandDeformFrac)) > 1e-6 {
		t.Fatalf("full-scale deform = %v, want the token %v (peak-to-peak)", got, WensityBandDeformFrac)
	}
	if got := wensityBandDeform(-3); got != 0 {
		t.Fatalf("negative spectrum must clamp to zero, got %v", got)
	}
	if got := wensityBandDeform(4); math.Abs(float64(got-WensityBandDeformFrac)) > 1e-6 {
		t.Fatalf("spectrum above 1 must clamp to the token, got %v", got)
	}
	prev := float32(-1)
	for i := 0; i <= 20; i++ {
		got := wensityBandDeform(float32(i) / 20)
		if got < prev-1e-9 {
			t.Fatalf("deform must not decrease: level %.2f gave %v after %v", float32(i)/20, got, prev)
		}
		prev = got
	}
}

// TestPseudoSpectrumDerivation drives pushLevel with a fixed feed and pins
// the three bands: band-mid IS the smoothed level; band-low lags it on the
// slow tau (400ms vs the 60ms attack) and holds the felt volume through the
// release; band-high spikes on the rise and decays between syllables. The
// same feed must produce the same bands twice (pure, no clock).
func TestPseudoSpectrumDerivation(t *testing.T) {
	feed := func() *liquidMotion {
		m := &liquidMotion{}
		m.beginSession()
		// One attack step at the 30fps cadence: raw 0->0.9.
		m.pushLevel(0.9, 33*time.Millisecond)
		return m
	}
	m := feed()
	if m.specLow <= 0 || m.specLow > m.level {
		t.Fatalf("band-low must lag the level on attack: specLow=%v level=%v", m.specLow, m.level)
	}
	if m.specHigh <= 0 {
		t.Fatal("a rising envelope must lift band-high (the sibilance proxy)")
	}
	// band-mid is level, by construction, not by copy drift.
	v := Visual{Glass: true}
	m.applyTo(&v, true)
	if v.SpectrumMid != m.level || v.SpectrumLow != m.specLow || v.SpectrumHigh != m.specHigh {
		t.Fatalf("applyTo must stamp the live bands: bands=(%v,%v,%v) motion=(%v,%v,%v)",
			v.SpectrumLow, v.SpectrumMid, v.SpectrumHigh, m.specLow, m.level, m.specHigh)
	}
	// Release: hold at zero; band-low decays slower than nothing at all, and
	// band-high decays fast toward the silence floor.
	m2 := feed()
	for i := 0; i < 20; i++ {
		m2.pushLevel(0, 33*time.Millisecond)
	}
	if m2.specHigh >= m.specHigh {
		t.Fatalf("band-high must decay when the rise stops: %v vs %v", m2.specHigh, m.specHigh)
	}
	if m2.specLow <= 0 {
		t.Fatal("band-low is the slow integrator: it must hold the felt volume through the pause")
	}
	// Determinism: the identical feed rebuilds the identical bands.
	m3 := feed()
	if m3.specLow != m.specLow || m3.specHigh != m.specHigh || m3.level != m.level {
		t.Fatalf("the spectrum derivation must be pure: (%v,%v,%v) vs (%v,%v,%v)",
			m3.specLow, m3.level, m3.specHigh, m.specLow, m.level, m.specHigh)
	}
}

// TestWensityDeformGate is the anti-inductive guard: the pseudo-spectrum
// reaches the renderer ONLY in Listening/Thinking/Speaking. Everything else
// - including the session states Conversation/Warm/Acting that still receive
// the audio feed - renders a zero-spectrum, non-deforming band, and a
// parked motion can never smuggle bands into a static state's frame.
func TestWensityDeformGate(t *testing.T) {
	yes := []statemachine.State{statemachine.StateListening, statemachine.StateThinking, statemachine.StateSpeaking}
	for _, s := range allKnownStates {
		want := false
		for _, y := range yes {
			if s == y {
				want = true
			}
		}
		if bandDeformState(s) != want {
			t.Fatalf("bandDeformState(%s) = %v, want %v", s, !want, want)
		}
	}
	// A live motion with a full spectrum: the gate decides what lands.
	m := &liquidMotion{}
	m.beginSession()
	m.pushLevel(1, 33*time.Millisecond)
	m.pushLevel(1, 33*time.Millisecond)
	if m.level <= 0 || m.specLow <= 0 {
		t.Fatalf("setup: the feed must produce a live spectrum (level=%v specLow=%v)", m.level, m.specLow)
	}
	v := Visual{Glass: true}
	m.applyTo(&v, false)
	if v.SpectrumLow != 0 || v.SpectrumMid != 0 || v.SpectrumHigh != 0 {
		t.Fatalf("a non-deform state must carry zero bands, got (%v,%v,%v)",
			v.SpectrumLow, v.SpectrumMid, v.SpectrumHigh)
	}
	m.applyTo(&v, true)
	if v.SpectrumMid == 0 {
		t.Fatal("the deform states must receive the spectrum")
	}
	// And the parked motion after endSession carries nothing at all.
	m.endSession()
	v2 := Visual{Glass: true}
	m.applyTo(&v2, true)
	if v2.SpectrumLow != 0 || v2.SpectrumMid != 0 || v2.SpectrumHigh != 0 ||
		v2.LiquidLevel != 0 || v2.SummonFlow != 0 {
		t.Fatalf("endSession must zero the spectrum: %+v", v2)
	}
}

// TestWensityRestingFrameIsStatic is the D32 static-frame pin for the band:
// the resting Sleeping mapping carries a glass body with a zero spectrum,
// so the band draws its spine at the token geometry with no deform and no
// drift - one deterministic frame, zero timers, exactly like the pre-65
// resting orb.
func TestWensityRestingFrameIsStatic(t *testing.T) {
	v := VisualFor(DarkPalette(), BallSizeDefaultPx, statemachine.StateSleeping, 0)
	if !v.Glass {
		t.Fatal("Sleeping must keep the glass body")
	}
	if v.SpectrumLow != 0 || v.SpectrumMid != 0 || v.SpectrumHigh != 0 {
		t.Fatalf("the resting frame must carry no spectrum: %+v", v)
	}
	if wensityBandDeform(v.SpectrumMid) != 0 {
		t.Fatal("the resting band must not deform")
	}
	// The static mapping never sets a spectrum for ANY state: bands arrive
	// only through the liquid seam (applyTo), never through VisualFor.
	for _, s := range allKnownStates {
		if vv := VisualFor(DarkPalette(), BallSizeDefaultPx, s, 0); vv.SpectrumLow != 0 || vv.SpectrumMid != 0 || vv.SpectrumHigh != 0 {
			t.Fatalf("VisualFor(%s) set a spectrum - bands must only come from the liquid seam", s)
		}
	}
}
