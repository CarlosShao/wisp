package ball

// C21 DesignTokens native side (SPEC-08 §2 "所有颜色/圆角/阴影/缓动必须取自
// C21 DesignTokens；design/assets/tokens.css 为唯一样式真相源").
//
// This file is the machine-checked twin of design/assets/tokens.css: every
// value below is copied verbatim from the CSS custom properties listed in
// docs/evidence/s1/c21-native-tokens.md (the cross-check document). Drawing
// code in this package references ONLY these values - a hardcoded hex/rgb
// literal in any other ball file is a review-rejecting violation (ticket 07
// constraint; grep-audited by TestNoHardcodedColorsInBallPackage).
//
// Both themes ship (D29 wallpaper adaptivity): Dark is the default, Light is
// the "白玻璃" counterpart from the [data-theme="light"] block. Tokens the
// light block does not redefine (onSolid, ballHalo) carry the root values.

// Theme selects the active palette (system theme follow lands with the
// settings wiring; the ball API takes it as a value so the renderer stays
// testable).
type Theme uint8

const (
	ThemeDark Theme = iota
	ThemeLight
)

// Color is a straight (non-premultiplied) RGBA color, 0..1 channels - the
// D2D1_COLOR_F shape. Premultiplied forms are derived where a bitmap needs
// them (ULW path), never stored.
type Color struct {
	R, G, B, A float32
}

// rgba builds a Color from 0..255 channels plus straight alpha 0..1.
func rgba(r, g, b uint8, a float32) Color {
	return Color{float32(r) / 255, float32(g) / 255, float32(b) / 255, a}
}

// hex builds a Color from a packed 0xRRGGBB value plus straight alpha 0..1.
func hex(rgb uint32, a float32) Color {
	return rgba(uint8(rgb>>16)&0xFF, uint8(rgb>>8)&0xFF, uint8(rgb)&0xFF, a)
}

// WithAlpha returns the same RGB with alpha replaced (opacity layering).
func (c Color) WithAlpha(a float32) Color { return Color{c.R, c.G, c.B, a} }

// Premultiplied returns the premultiplied form used by 32bpp ARGB bitmaps
// (UpdateLayeredWindow expects premultiplied alpha).
func (c Color) Premultiplied() Color {
	return Color{c.R * c.A, c.G * c.A, c.B * c.A, c.A}
}

// Palette is the ball-relevant slice of C21, per theme. Field comments name
// the CSS custom property; the values are the dark/light rows of
// design/assets/tokens.css (checked against the CSS by the evidence table).
type Palette struct {
	// Surfaces and glass edges.
	BGInset     Color // --bg-inset     rgba(6,8,11,.38) | rgba(20,24,28,.06)
	BGOverlay   Color // --bg-overlay   rgba(30,36,44,.50) | rgba(255,255,255,.72)
	GlassRing   Color // --glass-ring   rgba(255,255,255,.14) | rgba(255,255,255,.66)
	GlassHi     Color // --glass-hi     rgba(255,255,255,.18) | rgba(255,255,255,.90)
	GlassHiSoft Color // --glass-hi-soft rgba(255,255,255,.08) | rgba(255,255,255,.60)
	BorderSoft  Color // --border-soft  rgba(255,255,255,.10) | rgba(16,20,24,.12)

	// Text.
	FgPrimary   Color // --fg-primary   #F2F3F5 | #17191C
	FgSecondary Color // --fg-secondary #A6A9B0 | #55585F
	FgTertiary  Color // --fg-tertiary  #71747C | #82858C

	// Semantic colors (all low-saturation per SPEC-08 §6 rule 1).
	Accent      Color // --accent       #86C2B9 | #3E837A
	AccentLine  Color // --accent-line  rgba(134,194,185,.36) | rgba(62,131,122,.38)
	Danger      Color // --danger       #E07A70 | #BC544C
	DangerLine  Color // --danger-line  rgba(224,122,112,.38) | rgba(188,84,76,.36)
	Warn        Color // --warn         #D9B26A | #96762C
	WarnLine    Color // --warn-line    rgba(217,178,106,.38) | rgba(150,118,44,.36)
	Success     Color // --success      #7DB896 | #47805F
	SuccessLine Color // --success-line rgba(125,184,150,.38) | rgba(71,128,95,.36)
	Info        Color // --info         #8AAAD2 | #4A6C96
	InfoLine    Color // --info-line    rgba(138,170,210,.38) | rgba(74,108,150,.36)
	Warm        Color // --warm         #D6B184 | #8F6E3C
	WarmLine    Color // --warm-line    rgba(214,177,132,.38) | rgba(143,110,60,.36)

	// Orb glass (SPEC-08 §2.1 base form: "自发光玻璃滴").
	OrbHi     Color // --orb-hi     rgba(255,255,255,.90) | rgba(255,255,255,.96)
	OrbBody   Color // --orb-body   rgba(255,255,255,.30) | rgba(255,255,255,.62)
	OrbBody2  Color // --orb-body-2 rgba(255,255,255,.07) | rgba(255,255,255,.24)
	OrbRim    Color // --orb-rim    rgba(255,255,255,.30) | rgba(20,24,28,.16)
	OrbShadow Color // --orb-shadow rgba(0,0,0,.45) | rgba(20,24,28,.20)
	BallHalo  Color // --ball-halo  rgba(255,255,255,.10) (not redefined in light)

	// Per-state glow tints (v2 "核心发光色" pairs).
	TintAccentHi   Color // --tint-accent-hi   #C4E8E2 | #7FB5AC
	TintSuccessHi  Color // --tint-success-hi  #BFE3CF | #7FAE93
	TintSuccessLo  Color // --tint-success-lo  #4E8A69 | #35624A
	TintWarmHi     Color // --tint-warm-hi     #F0DCB8 | #C2A071
	TintWarmLo     Color // --tint-warm-lo     #96794C | #6E5732
	TintWarnHi     Color // --tint-warn-hi     #F0DCB8 | #B99A4E
	TintWarnLo     Color // --tint-warn-lo     #8F7434 | #6B5522
	TintDangerHi   Color // --tint-danger-hi   #F3C0BA | #C57E76
	TintDangerLo   Color // --tint-danger-lo   #9C4A43 | #83372F
	TintNeutralHi  Color // --tint-neutral-hi  #B9BCC4 | #9C9FA6
	TintNeutralMid Color // --tint-neutral-mid #6A6D75 | #71747B
	TintNeutralLo  Color // --tint-neutral-lo  #3E4148 | #4A4D54

	// Foreground on solid fills (badge digits, icons on tinted orb).
	OnSolid Color // --on-solid #FFFFFF (not redefined in light)
}

// DarkPalette returns the :root palette of tokens.css.
func DarkPalette() Palette {
	return Palette{
		BGInset:     rgba(6, 8, 11, 0.38),
		BGOverlay:   rgba(30, 36, 44, 0.50),
		GlassRing:   rgba(255, 255, 255, 0.14),
		GlassHi:     rgba(255, 255, 255, 0.18),
		GlassHiSoft: rgba(255, 255, 255, 0.08),
		BorderSoft:  rgba(255, 255, 255, 0.10),

		FgPrimary:   hex(0xF2F3F5, 1),
		FgSecondary: hex(0xA6A9B0, 1),
		FgTertiary:  hex(0x71747C, 1),

		Accent:      hex(0x86C2B9, 1),
		AccentLine:  rgba(134, 194, 185, 0.36),
		Danger:      hex(0xE07A70, 1),
		DangerLine:  rgba(224, 122, 112, 0.38),
		Warn:        hex(0xD9B26A, 1),
		WarnLine:    rgba(217, 178, 106, 0.38),
		Success:     hex(0x7DB896, 1),
		SuccessLine: rgba(125, 184, 150, 0.38),
		Info:        hex(0x8AAAD2, 1),
		InfoLine:    rgba(138, 170, 210, 0.38),
		Warm:        hex(0xD6B184, 1),
		WarmLine:    rgba(214, 177, 132, 0.38),

		OrbHi:     rgba(255, 255, 255, 0.90),
		OrbBody:   rgba(255, 255, 255, 0.30),
		OrbBody2:  rgba(255, 255, 255, 0.07),
		OrbRim:    rgba(255, 255, 255, 0.30),
		OrbShadow: rgba(0, 0, 0, 0.45),
		BallHalo:  rgba(255, 255, 255, 0.10),

		TintAccentHi:   hex(0xC4E8E2, 1),
		TintSuccessHi:  hex(0xBFE3CF, 1),
		TintSuccessLo:  hex(0x4E8A69, 1),
		TintWarmHi:     hex(0xF0DCB8, 1),
		TintWarmLo:     hex(0x96794C, 1),
		TintWarnHi:     hex(0xF0DCB8, 1),
		TintWarnLo:     hex(0x8F7434, 1),
		TintDangerHi:   hex(0xF3C0BA, 1),
		TintDangerLo:   hex(0x9C4A43, 1),
		TintNeutralHi:  hex(0xB9BCC4, 1),
		TintNeutralMid: hex(0x6A6D75, 1),
		TintNeutralLo:  hex(0x3E4148, 1),

		OnSolid: hex(0xFFFFFF, 1),
	}
}

// LightPalette returns the [data-theme="light"] palette of tokens.css.
func LightPalette() Palette {
	p := DarkPalette() // values the light block does not redefine carry over
	p.BGInset = rgba(20, 24, 28, 0.06)
	p.BGOverlay = rgba(255, 255, 255, 0.72)
	p.GlassRing = rgba(255, 255, 255, 0.66)
	p.GlassHi = rgba(255, 255, 255, 0.90)
	p.GlassHiSoft = rgba(255, 255, 255, 0.60)
	p.BorderSoft = rgba(16, 20, 24, 0.12)

	p.FgPrimary = hex(0x17191C, 1)
	p.FgSecondary = hex(0x55585F, 1)
	p.FgTertiary = hex(0x82858C, 1)

	p.Accent = hex(0x3E837A, 1)
	p.AccentLine = rgba(62, 131, 122, 0.38)
	p.Danger = hex(0xBC544C, 1)
	p.DangerLine = rgba(188, 84, 76, 0.36)
	p.Warn = hex(0x96762C, 1)
	p.WarnLine = rgba(150, 118, 44, 0.36)
	p.Success = hex(0x47805F, 1)
	p.SuccessLine = rgba(71, 128, 95, 0.36)
	p.Info = hex(0x4A6C96, 1)
	p.InfoLine = rgba(74, 108, 150, 0.36)
	p.Warm = hex(0x8F6E3C, 1)
	p.WarmLine = rgba(143, 110, 60, 0.36)

	p.OrbHi = rgba(255, 255, 255, 0.96)
	p.OrbBody = rgba(255, 255, 255, 0.62)
	p.OrbBody2 = rgba(255, 255, 255, 0.24)
	p.OrbRim = rgba(20, 24, 28, 0.16)
	p.OrbShadow = rgba(20, 24, 28, 0.20)

	p.TintAccentHi = hex(0x7FB5AC, 1)
	p.TintSuccessHi = hex(0x7FAE93, 1)
	p.TintSuccessLo = hex(0x35624A, 1)
	p.TintWarmHi = hex(0xC2A071, 1)
	p.TintWarmLo = hex(0x6E5732, 1)
	p.TintWarnHi = hex(0xB99A4E, 1)
	p.TintWarnLo = hex(0x6B5522, 1)
	p.TintDangerHi = hex(0xC57E76, 1)
	p.TintDangerLo = hex(0x83372F, 1)
	p.TintNeutralHi = hex(0x9C9FA6, 1)
	p.TintNeutralMid = hex(0x71747B, 1)
	p.TintNeutralLo = hex(0x4A4D54, 1)
	return p
}

// PaletteFor returns the palette of a theme.
func PaletteFor(t Theme) Palette {
	if t == ThemeLight {
		return LightPalette()
	}
	return DarkPalette()
}

// ------------------------------------------------------- liquid glass looks

// LiquidLook is one candidate colour treatment of the glass orb (ticket 62).
// The owner's reference images (图1/图2) are not on disk, so several
// treatments ship in the prototype and the owner picks on the real desktop
// with `balldebug -stay -look <name>` instead of the code betting on one.
//
// Every colour is a token here: the renderer draws with these values only
// (C21 / TestNoHardcodedColorsInBallPackage).
type LiquidLook struct {
	Name    string // -look flag value
	BlobA   Color  // dominant liquid body (largest, slowest)
	BlobB   Color  // counter-rotating body
	BlobC   Color  // accent wisp (state tint rides on this one)
	Deep    Color  // far-side refraction tint
	Rim     Color  // outer edge light - the contrast anchor on a white desktop
	Lip     Color  // inner bright lip just inside the rim
	Hi      Color  // specular highlight
	Caustic Color  // contact shadow / caustic under the orb
	Glow    Color  // outer halo
}

// looks are the candidate treatments. Order is the -look cycling order.
var looks = []LiquidLook{
	{ // vivo 蓝心小V reference family: electric indigo / violet / cyan.
		Name: "aurora", BlobA: hex(0x4B49FF, 0.95), BlobB: hex(0x8B5CF6, 0.90),
		BlobC: hex(0x22D3EE, 0.85), Deep: hex(0x1E1B7A, 0.60),
		Rim: hex(0x141A2E, 0.72), Lip: hex(0xE8ECFF, 0.75), Hi: rgba(255, 255, 255, 0.92),
		Caustic: hex(0x1B1F3A, 0.35), Glow: hex(0x6D7CFF, 0.55),
	},
	{ // glacier: cyan / teal / deep blue, cooler and calmer.
		Name: "glacier", BlobA: hex(0x22D3EE, 0.95), BlobB: hex(0x2DD4BF, 0.85),
		BlobC: hex(0x1D4ED8, 0.90), Deep: hex(0x0B3B66, 0.60),
		Rim: hex(0x0E2233, 0.72), Lip: hex(0xE6FBFF, 0.70), Hi: rgba(255, 255, 255, 0.92),
		Caustic: hex(0x0B2233, 0.34), Glow: hex(0x38BDF8, 0.52),
	},
	{ // nebula: magenta / indigo / pink, the loud end of the range.
		Name: "nebula", BlobA: hex(0xD946EF, 0.92), BlobB: hex(0x6366F1, 0.90),
		BlobC: hex(0xF472B6, 0.80), Deep: hex(0x3B0764, 0.55),
		Rim: hex(0x1A1030, 0.74), Lip: hex(0xFCE7FF, 0.70), Hi: rgba(255, 255, 255, 0.92),
		Caustic: hex(0x241033, 0.34), Glow: hex(0xC026D3, 0.50),
	},
	{ // solar: amber / coral / violet, for a warm reading of "cool colour".
		Name: "solar", BlobA: hex(0xF59E0B, 0.92), BlobB: hex(0xFB7185, 0.88),
		BlobC: hex(0x7C3AED, 0.80), Deep: hex(0x4A1D3A, 0.55),
		Rim: hex(0x241019, 0.72), Lip: hex(0xFFEFD6, 0.70), Hi: rgba(255, 255, 255, 0.92),
		Caustic: hex(0x2A1206, 0.32), Glow: hex(0xFDBA74, 0.50),
	},
}

// LookByName resolves a -look value (case-insensitive); ok=false on unknown.
func LookByName(name string) (LiquidLook, bool) {
	for _, l := range looks {
		if sameName(l.Name, name) {
			return l, true
		}
	}
	return LiquidLook{}, false
}

// sameName is a case-insensitive ASCII compare (tokens.go carries no imports
// so that it stays the plain data twin of tokens.css).
func sameName(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// LookNames lists the selectable treatments (for the debug flag help).
func LookNames() []string {
	out := make([]string, 0, len(looks))
	for _, l := range looks {
		out = append(out, l.Name)
	}
	return out
}

// DefaultLook is the treatment used when the consumer does not choose.
func DefaultLook() LiquidLook { return looks[0] }

// activeLook is the chosen treatment (process-wide, like pal).
var activeLook = DefaultLook()

// SetLook selects the liquid treatment; a repaint must follow. Returns false
// for an unknown name so the caller can report it.
func SetLook(name string) bool {
	l, ok := LookByName(name)
	if !ok {
		return false
	}
	activeLook = l
	return true
}

// Look returns the active liquid treatment.
func Look() LiquidLook { return activeLook }

// Liquid blob placement (fractions of the orb radius; the renderer rotates
// and swells these with the audio envelope - never re-allocates brushes).
const (
	// Rest orb: Sleeping keeps a real glass body instead of the 12px micro
	// dot, but smaller than the active orb so "asleep" still reads.
	SleepRestRatio = 0.62

	// Blob radii / offsets as fractions of the orb radius. They overlap on
	// purpose: the blend of three soft fields is what reads as liquid, and a
	// mostly-filled orb survives a white desktop far better than a shell.
	LiquidRadiusA = 0.72
	LiquidRadiusB = 0.62
	LiquidRadiusC = 0.50
	LiquidOffsetA = 0.22
	LiquidOffsetB = 0.30
	LiquidOffsetC = 0.40

	// Glass edge weights (physical px at 96 DPI, DPI-scaled at draw time).
	GlassRimPx   = 1.2 // dark outer edge light (contrast on white)
	GlassLipPx   = 1.0 // bright inner lip
	GlassCaustic = 0.30
	BorderRingPx = 1.8 // the "not speaking" border that fades in

	// Audio envelope -> liquid motion.
	SwimLevelGain   = 0.55 // blob offset swell per unit level
	SpinLevelGain   = 1.0  // rotation speed gain per unit level
	BorderOpenMs    = 220  // border fade-in (SPEC-08 dur-slow family)
	BorderCloseMs   = 180  // liquid converge when the voice stops
	SummonFlowMs    = 900  // one-shot liquid flow burst on summon
	DockAnimMs      = 160  // edge-dock squash / pop-back
	DockOverlapFrac = 0.42 // fraction of the orb left visible when docked
)

// Geometry + motion tokens (CSS lengths/ms verbatim; px at 96 DPI, scaled by
// the per-monitor DPI before use - never rescaled at call sites).
const (
	// --ball-size / --ball-sm / --ball-lg (config range 44-72, default 56).
	BallSizeDefaultPx = 56
	BallSizeSmallPx   = 44
	BallSizeLargePx   = 72
	BallSizeMinPx     = 44 // [ball] size_min
	BallSizeMaxPx     = 72 // [ball] size_max

	// Sleeping micro-dot: SPEC-08 §2 - "Sleeping 态直径缩到 12px 微点
	// (opacity 0.35)"; ball.html pins it at 12px regardless of user size.
	// Ticket 62 keeps the constants (the contract is only amended after owner
	// sign-off) but the prototype no longer renders them: the rest body is the
	// glass orb at SleepRestRatio, and this is its floor.
	SleepingDotPx = 12
	SleepOpacity  = 0.35

	// SleepingRestMinPx floors the resting orb so a small user size cannot
	// make it invisible again; RestSettledOpacity is where the Settling fade
	// lands (the resting orb must stay readable on a white desktop).
	SleepingRestMinPx  = 30
	RestSettledOpacity = 0.95

	// --dur-fast/base/slow. Motion budget: no duration above 300ms in UI
	// motion; per-state animation PERIODS (2.4s breathing etc.) are loops
	// sanctioned by SPEC-08 §2, not UI transitions.
	DurFastMs = 120
	DurBaseMs = 180
	DurSlowMs = 260

	// Animation periods (SPEC-08 §2.1 / ball.html).
	WarmBreathPeriodMs   = 2400 // Warm 2.4s opacity 0.55<->0.7
	SpeakingBreathMs     = 1600 // Speaking success breathe 1.6s
	ListeningWaveMs      = 2400 // three staggered rings, 0.8s offsets
	ThinkingSweepMs      = 1200 // bottom 2px light band, 1.2s
	ConfirmingPulseMs    = 2000 // danger pulse ring 2s
	FirstRunGuidePulseMs = 2600 // guide pulse 2.6s
	SettlingFadeMs       = 260  // opacity 1 -> 0.35 fade (SPEC-08 §2.1)

	// Warm breathing opacity band (SPEC-08 §2.1).
	WarmOpacityLow  = 0.55
	WarmOpacityHigh = 0.70

	// Frame cap: <=30fps (ticket 07 animation discipline).
	MaxAnimFPS    = 30
	MinFrameMs    = 1000 / MaxAnimFPS // 33ms
	WarmBreathFPS = 10                // compositor-friendly slow tick for Warm

	// Stroke weights.
	SlashStrokePx    = 1.5 // Muted slash (SPEC-08 §2.1)
	IconStrokePx     = 1.5 // Lucide-style icon stroke (SPEC-08 §6 rule 10)
	RingStrokePx     = 1.5 // Listening/Confirming ring width (ball.html)
	ConvRingStrokePx = 2.0 // Conversation ring "不得弱于 Confirming"
	ThinkingBandPx   = 2.0 // Thinking bottom light band height

	// Badge/dot geometry (ball.html: 17px min-width badge, 9px queue dot).
	BadgeDiameterPx = 17
	BadgeFontPx     = 10
	BadgeBorderPx   = 2
	QueueDotPx      = 9

	// Fonts: --font-sans family stack head + sizes used by the ball.
	FontFamily      = "Microsoft YaHei UI" // --font-sans CJK head on Windows
	FontSizeMonoPx  = 12.5                 // --t-mono
	FontSizeMicroPx = 11                   // --t-micro (digits/latin only)
	CountdownFontPx = 10                   // micro countdown digits under the ring
)
