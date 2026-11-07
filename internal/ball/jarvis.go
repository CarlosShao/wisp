package ball

// Ticket-65 Jarvis particle-orb style: pure geometry/motion mapping shared by
// the renderer (jarvis_windows.go does the D2D drawing) and the tests. Same
// posture as wensity.go: deliberately free of Win32/D2D so the mapping stays
// unit-testable on any GOOS. Every colour is a token (JarvisPalette in
// tokens.go / the C21 palette); every number below is a measurement of the
// reference tile grid (docs/evidence/s1/65-style-direction.md §1.2) - the
// renderer draws from these numbers only.

import "math"

// Ring identity of one particle. The measured radial brightness profile peaks
// at 0.07R / 0.16R / main (65-style-direction §1.2); the sparse in-orb field
// lifts the tile coverage into the measured 10-15% band.
const (
	jarvisRingInner uint8 = iota // the 0.16R inner orbit
	jarvisRingMain               // the per-state main orbit
	jarvisRingCore               // the 0.07R core cluster
	jarvisRingField              // sparse scatter inside the orb
)

// Jarvis geometry/motion tokens (fractions of the orb radius R unless stated;
// measured on jarvis-orb.png, 65-style-direction §1.2). Unexported on purpose:
// a within-style layout, not a knob - the same law as wensityStopTs.
const (
	jarvisRingInnerFrac     = 0.16               // inner orbit radius
	jarvisRingMainIdleFrac  = 0.50               // main orbit, idle/rest/success
	jarvisRingMainThinkFrac = 0.67               // main orbit, thinking (measured)
	jarvisRingMainAlertFrac = 0.39               // main orbit, alert (measured)
	jarvisRingCoreFrac      = 0.07               // core cluster radius
	jarvisCoreDotFrac       = 0.05               // the bright centre point (core diameter ~0.10R)
	jarvisRingTiltFrac      = 0.62               // orbit ellipse minor axis (the tilted-orbit read)
	jarvisFieldMinFrac      = 0.10               // field scatter inner bound
	jarvisFieldMaxFrac      = 0.85               // field scatter outer bound
	jarvisGlowRadiusFrac    = 1.05               // base tile-glow brush radius
	jarvisCoreGlowFrac      = 0.16               // core halo brush radius
	jarvisDotRadiusPx       = 1.0                // particle body radius at 96 DPI (1-2px dots)
	jarvisDotBrightPx       = 1.5                // bright (p90) particle radius
	jarvisBrightShare       = 0.18               // share of bright (p90) particles
	jarvisRingJitterMain    = 0.06               // main-orbit radial jitter (a band, not a circle)
	jarvisRingJitterInner   = 0.15               // inner-orbit radial jitter
	jarvisParticleCountBase = 96                 // particle budget at BallSizeDefaultPx (56px)
	jarvisParticleCountMin  = 60                 // measured visible-particle band floor
	jarvisParticleCountMax  = 120                // measured visible-particle band ceiling
	jarvisListenPulseFrac   = 0.08               // main-orbit swell per unit envelope (§3.2: 0.50R±0.08R)
	jarvisScatterPerLevel   = 0.20               // particle scatter per unit envelope (§3.2 +0.15R plus the floor, kept level-proportional)
	jarvisParticleSeed      = 0x9E3779B97F4A7C15 // fixed seed: the layout bakes once, never re-rolls per frame
)

// jarvisRng is a fixed-seed xorshift64* generator: the particle positions are
// computed ONCE per bake from this seed (there is no per-frame RNG anywhere),
// so the resting frame is byte-for-byte reproducible (D32) and two bakes of
// the same (state,size) agree bit for bit.
type jarvisRng struct{ s uint64 }

func newJarvisRng() jarvisRng { return jarvisRng{s: jarvisParticleSeed} }

func (g *jarvisRng) f64() float64 {
	x := g.s
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	g.s = x
	return float64((x*0x2545F4914F6CDD1D)>>11) / float64(1<<53)
}

// jarvisDot is one particle: its orbit, fixed angle and radius (fractions of
// R, jitter baked in), a deterministic wobble direction for the audio scatter,
// and whether it is a bright (p90) particle.
type jarvisDot struct {
	ring   uint8
	bright bool
	angle  float64
	radius float32
	wobble float32
}

// jarvisParticleCount scales the particle budget with the orb size (linear in
// the diameter): 96 at the default 56px, 75 at the 44px floor - inside the
// measured 60-120 visible-particle band at every configured size.
func jarvisParticleCount(sizePx float32) int {
	n := int(jarvisParticleCountBase*sizePx/float32(BallSizeDefaultPx) + 0.5)
	if n < jarvisParticleCountMin {
		n = jarvisParticleCountMin
	}
	if n > jarvisParticleCountMax {
		n = jarvisParticleCountMax
	}
	return n
}

// jarvisLayout computes the deterministic particle field for one size and
// main-orbit fraction: the three orbit rings carry ~90% of the dots, the
// leftover scatters inside the orb so the tile reads as a particle globe
// (reference coverage 10-15%, 65-style-direction §1.2). Pure: same inputs,
// same dots, bit for bit.
func jarvisLayout(sizePx, mainFrac float32) []jarvisDot {
	n := jarvisParticleCount(sizePx)
	g := newJarvisRng()
	dots := make([]jarvisDot, 0, n)
	// One orbit ring: evenly spaced angles plus a deterministic half-step
	// jitter, radius with the ring's band jitter, one rng draw each for the
	// bright tail, the wobble direction. The draws run in explicit statement
	// order so the field is fully seed-determined.
	ring := func(count int, frac, jitter float32, id uint8) {
		for i := 0; i < count; i++ {
			ang := float64(i)/float64(count)*2*math.Pi + (g.f64()-0.5)*(2*math.Pi/float64(count))
			bright := g.f64() < jarvisBrightShare
			rad := frac * (1 + float32(g.f64()-0.5)*2*jitter)
			wob := float32(g.f64())*2 - 1
			dots = append(dots, jarvisDot{ring: id, bright: bright, angle: ang, radius: rad, wobble: wob})
		}
	}
	// Orbit budget: main 55% / inner 25% / core 10%; the remainder fields.
	nMain := n * 55 / 100
	nInner := n * 25 / 100
	nCore := n * 10 / 100
	ring(nMain, mainFrac, jarvisRingJitterMain, jarvisRingMain)
	ring(nInner, jarvisRingInnerFrac, jarvisRingJitterInner, jarvisRingInner)
	ring(nCore, jarvisRingCoreFrac, jarvisRingCoreFrac*0.5, jarvisRingCore)
	for i := n - len(dots); i > 0; i-- {
		bright := g.f64() < jarvisBrightShare
		ang := g.f64() * 2 * math.Pi
		rad := jarvisFieldMinFrac + float32(g.f64())*(jarvisFieldMaxFrac-jarvisFieldMinFrac)
		wob := float32(g.f64())*2 - 1
		dots = append(dots, jarvisDot{ring: jarvisRingField, bright: bright, angle: ang, radius: rad, wobble: wob})
	}
	return dots
}

// jarvisMainRingFrac is the per-state structure law (65-style-direction §1.2:
// state differences are STRUCTURE - the main-orbit radius - not palette):
// thinking stretches the main orbit to the measured 0.67R, alert (Confirming's
// danger pulse / Error) contracts it to 0.39R, every other state - including
// success/Speaking, whose ember tile sampling was honestly recorded as
// unreliable - keeps the idle 0.50R. The renderer is a pure function of the
// Visual, so the state is read off its structural markers (Sweep = Thinking,
// RingPulse = Confirming, IconX = Error), exactly like the glass path reads
// its motion off LiquidLevel/Spectrum.
func jarvisMainRingFrac(v Visual) float32 {
	switch {
	case v.Sweep:
		return jarvisRingMainThinkFrac
	case v.RingPulse || v.Icon == IconX:
		return jarvisRingMainAlertFrac
	default:
		return jarvisRingMainIdleFrac
	}
}

// jarvisRingSpin advances the main-orbit phase ONLY while an animation frame
// feeds one (Thinking): every static state returns 0, so at rest the
// particles never spin (D32 - one static frame, zero timers).
func jarvisRingSpin(v Visual, phase float64) float64 {
	if !v.Sweep {
		return 0
	}
	return phase * 2 * math.Pi
}

// jarvisListenPulse is the main-orbit swell law (§3.2: 0.50R±0.08R×level);
// a silent frame carries pulse 1 - the identity.
func jarvisListenPulse(level float32) float32 {
	return 1 + jarvisListenPulseFrac*jarvisClamp01(level)
}

// jarvisScatterAmp is the audio scatter amplitude in R units. Kept strictly
// level-proportional (the spec's +0.05R base floor would move a silent active
// frame; the resting frame must stay the exact identity); a full envelope
// adds the measured 0.15R plus that floor.
func jarvisScatterAmp(level float32) float32 {
	return jarvisScatterPerLevel * jarvisClamp01(level)
}

func jarvisClamp01(f float32) float32 {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}

// Family selection (v1 seam; the [ball] jarvis_palette_hue config key is v2 -
// 65-style-direction §4.1/§7 ⑥). The cyan family (hue ~187deg) is the default;
// violet (p90 0x9F63FF measured) is the second reference family and proves the
// family swap is a pure token re-read with zero geometry change.
const (
	jarvisFamilyCyanName   = "cyan"
	jarvisFamilyVioletName = "violet"
)

type jarvisFamilyKind uint8

const (
	jarvisFamilyCyan jarvisFamilyKind = iota
	jarvisFamilyViolet
)

var activeJarvisFamily = jarvisFamilyCyan

// SetJarvisFamily selects the jarvis colour family; a repaint must follow.
// Returns false for an unknown name (mirrors SetStyle/SetLook).
func SetJarvisFamily(name string) bool {
	if sameName(name, jarvisFamilyCyanName) {
		activeJarvisFamily = jarvisFamilyCyan
		return true
	}
	if sameName(name, jarvisFamilyVioletName) {
		activeJarvisFamily = jarvisFamilyViolet
		return true
	}
	return false
}

// JarvisFamily reports the active jarvis colour family.
func JarvisFamily() string {
	if activeJarvisFamily == jarvisFamilyViolet {
		return jarvisFamilyVioletName
	}
	return jarvisFamilyCyanName
}

// jarvisFamilyColors returns the (dim, bright, core) particle colours of the
// active family - all three are token reads, zero literals.
func jarvisFamilyColors() (dim, bright, core Color) {
	jt := jarvisTokens
	if activeJarvisFamily == jarvisFamilyViolet {
		return jt.VioletDim, jt.VioletBright, jt.VioletCore
	}
	return jt.ParticleDim, jt.ParticleBright, jt.Core
}
