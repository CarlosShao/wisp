//go:build windows

package ball

// Ticket-65 Jarvis particle-orb renderer (docs/evidence/s1/65-style-direction.md
// §1.2/§2 - the second A1581-approved style family; the mirror of the glass
// renderer's wensity branch). Pure D2D: a dark base glow, 2-3 elliptical
// particle orbit rings plus a sparse field, and the core highlight. Baked once
// per state/size/DPI/style change in setAssets exactly like every other brush
// family: an animated frame only FILLS with existing COM objects (the D32 CPU
// discipline). The particle layout is computed once per bake from a fixed seed
// (jarvis.go) and never re-rolled per frame; the resting frame is one
// deterministic static frame with zero timers.

import (
	"math"
	"unsafe"
)

// jarvisAssets is the pre-baked jarvis material for the CURRENT state visual
// (rebuilt by setAssets on state/size/dpi/theme/style changes).
type jarvisAssets struct {
	glow unsafe.Pointer // base tile glow (radial)
	core unsafe.Pointer // core halo (radial)
	dots []jarvisDot    // fixed-seed particle layout for this state
}

// setJarvis bakes the jarvis brushes + layout for one state visual.
func (r *renderer) setJarvis(v Visual) {
	R := v.SizePx / 2
	jt := jarvisTokens
	r.jarvis.dots = jarvisLayout(v.SizePx, jarvisMainRingFrac(v))
	r.jarvis.glow = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(mulA(jt.BgGlow, v.Opacity))},
		{0.55, colorF(mulA(jt.BgGlow, 0.55*v.Opacity))},
		{1.00, colorF(mulA(jt.BgGlow, 0))},
	}, r.center(), d2d1Point2F{0, 0}, R*jarvisGlowRadiusFrac, R*jarvisGlowRadiusFrac)
	coreCol := jarvisCoreColor(v)
	r.jarvis.core = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(mulA(coreCol, v.Opacity))},
		{0.40, colorF(mulA(coreCol, 0.55*v.Opacity))},
		{1.00, colorF(mulA(coreCol, 0))},
	}, r.center(), d2d1Point2F{0, 0}, R*jarvisCoreGlowFrac, R*jarvisCoreGlowFrac)
}

// dropJarvis releases the jarvis material (dropAssets calls it with the rest).
func (r *renderer) dropJarvis() {
	comRelease(&r.jarvis.glow)
	comRelease(&r.jarvis.core)
	r.jarvis.dots = nil
}

// jarvisDotColors picks the particle pair for the frame: the active family's
// (dim, bright) tokens, the thinking p90 override while Thinking, and the
// danger-derived C21 pair for the alert structure (§3.2: the alert palette is
// danger-derived). Zero literals - C21 tokens only.
func jarvisDotColors(v Visual) (dim, bright Color) {
	switch {
	case v.RingPulse || v.Icon == IconX:
		return pal.TintDangerLo, pal.TintDangerHi
	case v.Sweep:
		jt := jarvisTokens
		return jt.ParticleDim, jt.ParticleThinking
	default:
		dim, bright, _ = jarvisFamilyColors()
		return dim, bright
	}
}

// jarvisCoreColor picks the core highlight for the frame: danger hi on alert,
// the C21 success tint on success/Speaking (the ember tile's sampling was
// recorded unreliable, so success derives from C21 per §1.2), the thinking
// white, else the active family's core token.
func jarvisCoreColor(v Visual) Color {
	switch {
	case v.RingPulse || v.Icon == IconX:
		return pal.TintDangerHi
	case v.Icon == IconVolume:
		return pal.TintSuccessHi
	case v.Sweep:
		return jarvisTokens.CoreThinking
	default:
		_, _, core := jarvisFamilyColors()
		return core
	}
}

// drawJarvis paints the particle-orb frame: base glow, particle dots, core.
// Everything is a cached brush fill - no COM object is created per frame, so
// an animated frame costs the same as the static one.
func (r *renderer) drawJarvis(v Visual, c d2d1Point2F, R, s float32, phase float64) {
	// The edge dock squeezes the whole body along the dock axis (the same law
	// drawGlass applies: the tab is the squashed orb tucked against the edge).
	dockX, dockY := float32(1), float32(1)
	switch v.Dock {
	case EdgeLeft, EdgeRight:
		dockX = DockSquash(v.DockProgress)
	case EdgeTop, EdgeBottom:
		dockY = DockSquash(v.DockProgress)
	case EdgeNone:
	}

	// 0. Base tile glow: the dark backdrop that lets the particles read on a
	// light desktop (the reference tile base, JarvisPalette.BgGlow).
	if r.jarvis.glow != nil {
		r.fillEllipse(c, R*jarvisGlowRadiusFrac*dockX, R*jarvisGlowRadiusFrac*dockY, r.jarvis.glow)
	}

	// 1. The particle field. Per frame: the main orbit breathes with the
	// envelope (Listening), its phase advances with the animation frame
	// (Thinking), and the audio scatter wobbles every dot (level-proportional,
	// so a silent or resting frame is the exact baked geometry - D32).
	dim, bright := jarvisDotColors(v)
	spin := jarvisRingSpin(v, phase)
	pulse := jarvisListenPulse(v.LiquidLevel)
	scatter := jarvisScatterAmp(v.LiquidLevel)
	dotR := jarvisDotRadiusPx * s
	brightR := jarvisDotBrightPx * s
	for i := range r.jarvis.dots {
		d := &r.jarvis.dots[i]
		rad := d.radius
		ang := d.angle
		if d.ring == jarvisRingMain {
			rad *= pulse
			ang += spin
		}
		rad += d.wobble * scatter
		tilt := float32(1)
		if d.ring != jarvisRingField {
			tilt = jarvisRingTiltFrac // the orbits read as tilted ellipses
		}
		col, rr := dim, dotR
		if d.bright {
			col, rr = bright, brightR
		}
		r.setSolid(mulA(col, v.Opacity))
		r.fillEllipse(d2d1Point2F{
			x: c.x + rad*R*float32(math.Cos(ang))*dockX,
			y: c.y + rad*R*tilt*float32(math.Sin(ang))*dockY,
		}, rr, rr, r.solid)
	}

	// 2. The core: soft halo + the bright centre point (round, not tilted).
	if r.jarvis.core != nil {
		r.fillEllipse(c, R*jarvisCoreGlowFrac*dockX, R*jarvisCoreGlowFrac*dockY, r.jarvis.core)
	}
	r.setSolid(mulA(jarvisCoreColor(v), v.Opacity))
	r.fillEllipse(c, R*jarvisCoreDotFrac*dockX, R*jarvisCoreDotFrac*dockY, r.solid)
}
