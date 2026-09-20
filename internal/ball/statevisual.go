package ball

import "github.com/CarlosShao/wisp/internal/statemachine"

// IconKind selects the vector glyph drawn inside the orb (SPEC-08 §6 rule 10:
// zero emoji, Lucide-style line icons, stroke 1.5). All glyphs are drawn with
// D2D primitives in renderer_windows.go - no bitmaps, no font glyphs.
type IconKind uint8

const (
	IconNone          IconKind = iota
	IconAudioLines             // Listening (audio-lines)
	IconVolume                 // Speaking (volume-2)
	IconMic                    // Conversation (mic)
	IconSlash                  // Muted (1.5px slash, SPEC-08 §2.1)
	IconX                      // Error (x)
	IconWifiOff                // NoNetwork (wifi-off)
	IconKey                    // Unconfigured (key-round)
	IconArrowUpRight           // FirstRun (arrow-up-right)
	IconAlertTriangle          // WatchdogAlert (alert-triangle)
	IconRefresh                // Stuck (refresh-cw)
)

// Visual is the complete render recipe of one state at one instant. The
// renderer consumes only this struct - state knowledge lives here, pixels
// live there.
type Visual struct {
	SizePx     float32 // orb diameter in px, drawn as-is (DPI note after this struct)
	Opacity    float32 // whole-orb opacity (Sleeping 0.35, Armed 0.6, Muted 0.4...)
	CoreColor  Color   // inner glow color (tint hi/lo pair of the state)
	CoreAlpha  float32 // inner glow strength (ball.html --tint-opacity)
	GlowColor  Color   // outer glow (ball.html --orb-glow)
	RingColor  Color   // outer ring color; A==0 means "no ring"
	RingWidth  float32 // ring stroke width (px @96)
	RingPulse  bool    // danger pulse animation active (Confirming)
	Icon       IconKind
	BadgeCount int     // >0 draws the top-right approval badge with this number
	QueueDot   bool    // bottom-left info dot (Queued overlay)
	Progress   float32 // 0..1: downloading ring progress; <0 = no progress ring
	BadgeText  string  // non-empty: center text override (countdown / percent)
	Sweep      bool    // Thinking: flowing 2px band at the orb bottom

	// --- ticket 62 liquid-glass form (prototype mode; see
	// EnablePrototypeVisuals). With the mode off every field below keeps the
	// frozen SPEC-08 §2.1 reading, which is the escape hatch (`balldebug
	// -frozen`) now that the owner has signed the new look provisionally
	// (SPEC-08 §2 INTERIM, 2026-09-20; the texture is still 赝品 - ticket 65).
	Glass bool // draw the glass body + liquid
	// BorderAlpha is how far the "not speaking" border has faded in (0..1).
	BorderAlpha float32
	// LiquidLevel is the audio envelope 0..1 driving blob swell and rotation
	// rate; 0 when nobody is speaking (the liquid converges, border stays).
	LiquidLevel float32
	// LiquidAngle is the accumulated rotation of the liquid (radians). The
	// window owns the accumulator (audio sets the rate), so the renderer stays
	// a pure function of the visual - no clock inside the draw.
	LiquidAngle float64
	// SummonFlow is the one-shot flow burst 0..1 after the ball is summoned.
	SummonFlow float32
	// Dock is the screen edge the orb squashes against (EdgeNone when free);
	// DockProgress is how far that squash has completed (0..1).
	Dock         Edge
	DockProgress float32
}

// SizePx's DPI treatment is asymmetric, and ticket 68 AC#1 records it as found
// rather than fixing it here (every diff number in docs/SLO.md A.2 and
// docs/evidence/s1/62-* was measured at 96 DPI, where the asymmetry is a
// no-op): the DC render target is created with dpiX/dpiY = 96, so its units ARE
// physical pixels, and drawFrame uses R = SizePx/2 unscaled - while the stroke
// widths and ring margins get the dpi/96 factor and the window edge gets
// WindowEdgePx(). Above 96 DPI the body therefore keeps its pixel size inside a
// window that grew (a 56 configured orb draws 56px of body in a 108px window at
// 144 DPI). Scaling the body with the window is a rendering change, so it
// belongs with AC#2's desktop run, not with a comment.

// Edge names a screen edge the ball docks against.
type Edge uint8

const (
	EdgeNone Edge = iota
	EdgeLeft
	EdgeRight
	EdgeTop
	EdgeBottom
)

// prototypeVisuals selects the ticket 62 liquid-glass rendering. It is still
// OFF by default, but the reason is no longer "the owner has not signed the
// new look": on 2026-09-20 the owner signed it PROVISIONALLY (SPEC-08 §2
// INTERIM - "先勉强用吧"), and SPEC-08 §2's Sleeping row now describes THIS
// mode. Flipping the default is ticket 68 AC#2, held back only because its
// re-measurement needs a desktop that is not busy; EnablePrototypeVisuals(false)
// and `balldebug -frozen` stay as the escape hatch back to the old frozen table.
//
// Two facts AC#1 established from code, both unresolved by this ticket:
//   - the resting body is stateSize()'s 56*SleepRestRatio = 34.72px, which is
//     NOT the "直径 44px" SPEC-08 §2 now claims (see TestSleepingSizeTruthTable
//     for the derivation of the recorded diff boxes) - 不符, awaiting ruling;
//   - cmd/balldebug is today the only caller that turns this on, so a library
//     consumer (the panel, a future wisp GUI host) still gets the frozen look.
var prototypeVisuals bool

// EnablePrototypeVisuals turns on the ticket 62 liquid-glass mapping. The
// caller must trigger a repaint (and a setAssets rebuild) afterwards.
func EnablePrototypeVisuals(on bool) { prototypeVisuals = on }

// PrototypeVisualsEnabled reports the mode (test/evidence seam).
func PrototypeVisualsEnabled() bool { return prototypeVisuals }

// stateSize returns the configured orb size clamped to the token range.
// While the mode is off this is the frozen SPEC-08 §2.1 reading: Sleeping is
// the fixed 12px micro dot, independent of the user's size (ball.html: "不随用
// 户配置的尺寸放大"). Prototype mode replaces it with the resting glass orb,
// because the 12px dot measured 0 visible pixels on the owner's desktop:
// configuredPx*SleepRestRatio, floored to SleepingRestMinPx - i.e. 34.72px at
// the default 56, and 30px (the floor) for any configured size from 44 to 48.
// The dock never enters here: a docked Sleeping orb keeps this same body size,
// and only drawGlass's DockSquash factor plus the window's new origin change
// what lands on screen (ticket 68 AC#1 column 3).
func stateSize(configuredPx int, s statemachine.State) float32 {
	if configuredPx < BallSizeMinPx {
		configuredPx = BallSizeMinPx
	}
	if configuredPx > BallSizeMaxPx {
		configuredPx = BallSizeMaxPx
	}
	if s == statemachine.StateSleeping {
		if !prototypeVisuals {
			return SleepingDotPx
		}
		size := float32(configuredPx) * SleepRestRatio
		if size < SleepingRestMinPx {
			size = SleepingRestMinPx
		}
		return size
	}
	return float32(configuredPx)
}

// VisualFor maps a state to its SPEC-08 §2.1 visual using the palette.
// Settling additionally takes fadeProgress (0=enter, 1=fully faded) to drive
// the 260ms opacity fade; other states ignore it.
func VisualFor(p Palette, configuredPx int, s statemachine.State, fadeProgress float32) Visual {
	v := Visual{
		SizePx:  stateSize(configuredPx, s),
		Opacity: 1,
		// Default glass: no tint, halo glow (ball.html base .ball).
		CoreColor: p.BallHalo,
		CoreAlpha: 0,
		GlowColor: p.BallHalo,
		Progress:  -1,
	}
	switch s {
	case statemachine.StateSleeping:
		v.Opacity = SleepOpacity
		v.GlowColor = p.BallHalo
	case statemachine.StateArmed:
		v.Opacity = 0.6
	case statemachine.StateMuted:
		v.Opacity = 0.4
		v.Icon = IconSlash
	case statemachine.StateListening:
		v.RingColor = p.Accent
		v.RingWidth = RingStrokePx
		v.Icon = IconAudioLines
	case statemachine.StateThinking:
		// 2px accent light band flowing at the orb bottom (renderer animates
		// the band position; color fixed).
		v.CoreColor = p.Accent
		v.CoreAlpha = 0.35
		v.Sweep = true
	case statemachine.StateActing:
		v.CoreColor = p.Accent
		v.CoreAlpha = 0.95
		v.GlowColor = p.AccentLine
	case statemachine.StateConfirming:
		v.RingColor = p.Danger
		v.RingWidth = RingStrokePx
		v.RingPulse = true
	case statemachine.StateAwaitingApproval:
		v.RingColor = p.Danger
		v.RingWidth = RingStrokePx
		v.BadgeCount = 1 // depth rendered by the caller via SetBadge
	case statemachine.StateSpeaking:
		v.CoreColor = p.Success
		v.CoreAlpha = 0.95
		v.GlowColor = p.SuccessLine
		v.Icon = IconVolume
	case statemachine.StateSettling:
		// SPEC-08 §2.1: opacity 1 -> 0.35 over 260ms.
		v.Opacity = 1 - fadeProgress*(1-SleepOpacity)
	case statemachine.StateWarm:
		v.CoreColor = p.Warm
		v.CoreAlpha = 0.9
		v.GlowColor = p.WarmLine
	case statemachine.StateConversation:
		// "danger 常亮环，不得渐隐、不得弱于 Confirming" - 2px ring, solid core.
		v.RingColor = p.Danger
		v.RingWidth = ConvRingStrokePx
		v.CoreColor = p.Danger
		v.CoreAlpha = 0.95
		v.GlowColor = p.DangerLine
		v.Icon = IconMic
	case statemachine.StateDownloading:
		v.CoreColor = p.BGOverlay
		v.CoreAlpha = 0.7
		v.Progress = 0.0 // caller updates via SetProgress
	case statemachine.StateError:
		v.CoreColor = p.Danger
		v.CoreAlpha = 0.95
		v.GlowColor = p.DangerLine
		v.Icon = IconX
	case statemachine.StateNoNetwork:
		v.CoreColor = p.TintNeutralMid
		v.CoreAlpha = 0.8
		v.Icon = IconWifiOff
	case statemachine.StateUnconfigured:
		v.CoreColor = p.Warn
		v.CoreAlpha = 0.95
		v.GlowColor = p.WarnLine
		v.Icon = IconKey
	case statemachine.StateFirstRun:
		v.CoreColor = p.Accent
		v.CoreAlpha = 0.95
		v.GlowColor = p.AccentLine
		v.Icon = IconArrowUpRight
	case statemachine.StateWatchdogAlert:
		v.CoreColor = p.Warn
		v.CoreAlpha = 0.95
		v.GlowColor = p.WarnLine
		v.Icon = IconAlertTriangle
	case statemachine.StateQueued:
		// Queued is an overlay on the main state: S1 renders the accent core
		// plus the info dot (the "main state" compositing lands with the
		// session wiring; the dot is the load-bearing signal).
		v.CoreColor = p.Accent
		v.CoreAlpha = 0.95
		v.GlowColor = p.AccentLine
		v.QueueDot = true
	case statemachine.StateStuck:
		v.CoreColor = p.Warn
		v.CoreAlpha = 0.95
		v.GlowColor = p.WarnLine
		v.Icon = IconRefresh
	}
	if prototypeVisuals {
		applyGlassForm(&v, p, s, fadeProgress)
	}
	return v
}

// borderAtRest is the border level a state's frame carries once its transition
// has landed: the "not speaking" border rides every session frame, and is away
// in the two states where it would be a lie - Sleeping (the resting body has
// no border) and Speaking (the assistant IS the one talking, so the liquid
// owns the orb). This single function is both the frozen prototype mapping in
// applyGlassForm and the target the transition ramps toward, so the two can
// never disagree about what the border should end up being.
func borderAtRest(s statemachine.State) float32 {
	switch s {
	case statemachine.StateSleeping, statemachine.StateSpeaking:
		return 0
	}
	return 1
}

// applyGlassForm is the ticket 62 remap: every state becomes the same glass
// material, distinguished by its liquid colour field, its border and its
// motion - not by the visibility of a small dot. It only ever runs in
// prototype mode (see EnablePrototypeVisuals).
func applyGlassForm(v *Visual, p Palette, s statemachine.State, fadeProgress float32) {
	lk := activeLook
	v.Glass = true
	v.BorderAlpha = borderAtRest(s) // resting level; the transition owns the travel
	// v.SizePx is already the right body: stateSize() returned the rest orb in
	// prototype mode and the configured size everywhere else.

	// The state's semantic colour rides on the third liquid blob + the halo,
	// so the 20 states stay distinguishable inside one material.
	v.GlowColor = lk.Glow
	switch s {
	case statemachine.StateSleeping:
		// Resting: fully opaque glass orb (the fix - the old 12px/0.35 dot
		// measured invisible), liquid still, border away. Static frame.
		v.Opacity = 1
		v.GlowColor = mulA(lk.Glow, 0.45)
	case statemachine.StateArmed:
		v.Opacity = 0.92
	case statemachine.StateMuted:
		v.Opacity = 0.85
	case statemachine.StateSettling:
		// Fades to the RESTING body, not to an invisible micro dot. The border
		// is not doubled onto this fade: borderAtRest already lands it.
		v.Opacity = 1 - fadeProgress*(1-RestSettledOpacity)
	case statemachine.StateListening, statemachine.StateSpeaking:
		v.GlowColor = mixLook(lk, v.CoreColor)
	}
	if v.CoreAlpha > 0 {
		// A state with a solid core keeps that core as its identity colour.
		v.GlowColor = mixLook(lk, v.CoreColor)
	}
}

// mixLook keeps the chosen treatment dominant while pulling the halo toward
// the state colour (state identity must survive the shared material).
func mixLook(lk LiquidLook, state Color) Color {
	const w = float32(0.45)
	return Color{
		R: lk.Glow.R*(1-w) + state.R*w,
		G: lk.Glow.G*(1-w) + state.G*w,
		B: lk.Glow.B*(1-w) + state.B*w,
		A: lk.Glow.A,
	}
}

// Animated reports whether the state runs an animation timer (SPEC-08 §2:
// only Listening/Thinking/Acting/Speaking/Confirming start animation timers;
// Warm breathing and the Settling fade are the two SPEC-sanctioned
// exceptions - see anim.go for the full policy).
func Animated(s statemachine.State) bool {
	switch s {
	case statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateActing, statemachine.StateSpeaking,
		statemachine.StateConfirming:
		return true
	}
	return false
}
