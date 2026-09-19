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
	SizePx     float32 // orb diameter at 96 DPI (scaled by DPI later)
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
}

// stateSize returns the configured orb size clamped to the token range.
// Sleeping overrides to the fixed micro dot (ball.html: "不随用户配置的尺寸放大").
func stateSize(configuredPx int, s statemachine.State) float32 {
	if s == statemachine.StateSleeping {
		return SleepingDotPx
	}
	if configuredPx < BallSizeMinPx {
		configuredPx = BallSizeMinPx
	}
	if configuredPx > BallSizeMaxPx {
		configuredPx = BallSizeMaxPx
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
	return v
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
