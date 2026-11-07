//go:build windows

package ball

// The ball renderer: an ID2D1DCRenderTarget bound to a top-down 32bpp DIB
// selected into a memory DC. Each frame draws the state Visual into the DIB
// and the window presents it with UpdateLayeredWindow (premultiplied alpha,
// SPEC-08 §2). The render target is created with DPI 96/96 so 1 DIP == 1
// physical pixel - all layout math below is in physical pixels and DPI only
// enters through the sizes computed in hit.go / tokens.go.
//
// COM object pointers are stored as unsafe.Pointer (vet-clean provenance);
// syscall args are uintptr words.

import (
	"fmt"
	"log/slog"
	"math"
	"unsafe"

	"golang.org/x/sys/windows"
)

type renderer struct {
	dcRT       unsafe.Pointer // ID2D1DCRenderTarget
	hdcMem     uintptr
	hbm        uintptr // DIB section handle
	bits       uintptr // DIB memory (Windows heap - NOT a Go pointer; GC must not track it)
	capW, capH uint32  // current DIB allocation
	w, h       int32   // physical pixels
	dpi        uint32

	solid    unsafe.Pointer // single solid brush, SetColor per use
	badgeFmt unsafe.Pointer // DWrite text format: badge digits (weight 600)
	microFmt unsafe.Pointer // DWrite text format: countdown/percent (micro)

	// Gradient brush cache for the CURRENT state frame set. Rebuilt by
	// setAssets (state/size/dpi/theme change); lazily extended for animated
	// alphas (breathing/pulse pick the nearest cached variant).
	body unsafe.Pointer
	hi   unsafe.Pointer
	core unsafe.Pointer
	glow map[int]unsafe.Pointer // alpha-permille -> radial glow brush

	// Liquid-glass material (ticket 62). The colour blobs are pre-baked at a
	// fixed angular ladder so an animated frame only picks an index - no COM
	// object is created during an animation tick (D32 CPU discipline).
	liq     [3]liquidLayer
	caustic unsafe.Pointer // contact shadow under the orb

	// Wensity band material (ticket 65): pre-baked segments along the band
	// span, plus the white core and the bottom inner glow. Baked only when
	// the active look is the wensity treatment; released with the rest.
	band        [wensityBandSegs]unsafe.Pointer
	bandCore    unsafe.Pointer
	wensityGlow unsafe.Pointer

	// Jarvis particle material (ticket 65): the base tile glow, the core
	// halo and the fixed-seed particle layout. Baked only when the active
	// style is the jarvis family; released with the rest.
	jarvis jarvisAssets

	// overlayOnly (ticket 65 rework): the WebView2 orb window is the ball's
	// visible face and overlays this window, so frames draw ONLY what the
	// page cannot - the approval badge, progress ring, queue dot, text
	// override - plus a near-invisible hit-test silhouette (the layered
	// window is OS-hit-tested per pixel; all-zero alpha would make the ball
	// unclickable, the measured ticket-65 regression). Set by the window
	// before every draw (ball_windows.go); direct drawFrame callers (tests)
	// default to the full D2D frame.
	overlayOnly bool
}

// liquidLadder is the angular resolution of the pre-baked blob positions.
// 24 steps over a full turn is visually smooth at 30fps (a frame advances at
// most a few steps) and costs 3*24 = 72 tiny brushes per state.
const liquidLadder = 24

// liquidLayer is one colour blob: 24 pre-baked radial gradients placed around
// the orb centre, plus their centres (used when the envelope shrinks the
// visible patch).
type liquidLayer struct {
	brushes [liquidLadder]unsafe.Pointer
	centers [liquidLadder]d2d1Point2F
}

// per-blob motion: relative rotation rate and angular phase offset. Blob B
// counter-rotates, which is what reads as "flowing" rather than "spinning".
var (
	liqRate   = [3]float64{1.0, -1.6, 0.55}
	liqPhase  = [3]float64{0, 2.1, 4.2}
	liqRadius = [3]float32{LiquidRadiusA, LiquidRadiusB, LiquidRadiusC}
	liqOffset = [3]float32{LiquidOffsetA, LiquidOffsetB, LiquidOffsetC}
)

// pal is the ACTIVE C21 palette (dark by default; switched by SetTheme).
// All renderer color reads go through it - the machine rule
// TestNoHardcodedColorsInBallPackage keeps literal colors out of this file.
var pal = DarkPalette()

// SetTheme switches the active C21 palette (D29 system-theme follow). The
// caller must trigger a repaint afterwards.
func SetTheme(t Theme) {
	pal = PaletteFor(t)
}

func newRenderer(hwnd windows.HWND, w, h int32, dpi uint32) (*renderer, error) {
	if err := ensureFactories(); err != nil {
		return nil, err
	}
	r := &renderer{w: w, h: h, dpi: dpi, glow: map[int]unsafe.Pointer{}}

	// Memory DC + top-down 32bpp DIB (UpdateLayeredWindow needs premultiplied
	// ARGB; D2D renders premultiplied into the DIB).
	hdcScreen, _, _ := pGetDC.Call(0)
	if hdcScreen == 0 {
		return nil, fmt.Errorf("ball: GetDC(NULL) failed")
	}
	defer pReleaseDC.Call(0, hdcScreen)

	r.hdcMem, _, _ = pCreateCompatibleDC.Call(hdcScreen)
	if r.hdcMem == 0 {
		return nil, fmt.Errorf("ball: CreateCompatibleDC failed")
	}
	bmi := bitmapInfo{}
	bmi.header = bitmapInfoHeader{
		size:     uint32(unsafe.Sizeof(bmi.header)),
		width:    w,
		height:   -h, // negative = top-down
		planes:   1,
		bitCount: 32,
	}
	var bits uintptr
	hbm, _, err := pCreateDIBSection.Call(r.hdcMem, unsafePtr(&bmi), dibRGBColors, unsafePtr(&bits), 0, 0)
	if hbm == 0 {
		pDeleteDC.Call(r.hdcMem)
		return nil, fmt.Errorf("ball: CreateDIBSection(%dx%d): %v", w, h, err)
	}
	r.hbm = hbm
	r.bits = bits
	pSelectObject.Call(r.hdcMem, hbm)

	// DC render target with premultiplied alpha (the ULW combination).
	props := d2d1RenderTargetProps{
		format:    dxgiFormatB8G8R8A8,
		alphaMode: d2d1AlphaModePremultiplied,
		dpiX:      96,
		dpiY:      96,
	}
	var dcRT unsafe.Pointer
	hr, _, _ := comCall(d2dFactory, slotFactoryCreateDCRT, unsafePtr(&props), unsafePtr(&dcRT))
	if hr != 0 {
		r.release()
		return nil, fmt.Errorf("ball: CreateDCRenderTarget hr=0x%x", hr)
	}
	r.dcRT = dcRT

	bindRect := rect{l: 0, t: 0, r: w, b: h}
	hr, _, _ = comCall(r.dcRT, slotDCRTBindDC, uintptr(r.hdcMem), unsafePtr(&bindRect))
	if hr != 0 {
		r.release()
		return nil, fmt.Errorf("ball: BindDC hr=0x%x", hr)
	}

	// Solid brush.
	white := colorF(pal.OnSolid)
	brushProps := identityBrushProps()
	var solid unsafe.Pointer
	hr, _, _ = comCall(r.dcRT, slotRTCreateSolidColorBrush,
		unsafePtr(&white), unsafePtr(&brushProps), unsafePtr(&solid))
	if hr != 0 {
		r.release()
		return nil, fmt.Errorf("ball: CreateSolidColorBrush hr=0x%x", hr)
	}
	r.solid = solid

	// DWrite text formats (family/sizes from C21; DirectWrite only - no GDI
	// text, SPEC-08 §2).
	if err := r.createTextFormats(); err != nil {
		r.release()
		return nil, err
	}
	return r, nil
}

func (r *renderer) createTextFormats() error {
	family := utf16(FontFamily)
	locale := utf16("zh-CN")
	mk := func(weight uint32, size float32) (unsafe.Pointer, error) {
		var tf unsafe.Pointer
		hr, _, _ := comCall(dwriteFactory, slotDWriteCreateTextFormat,
			unsafePtr(family), 0, uintptr(weight), 0 /*normal*/, 5, /*normal stretch*/
			f32bits(size), unsafePtr(locale), unsafePtr(&tf))
		if hr != 0 {
			return nil, fmt.Errorf("ball: CreateTextFormat(%v) hr=0x%x", size, hr)
		}
		// Center align horizontally + vertically.
		_, _, _ = comCall(tf, slotTFSetTextAlignment, 2)      // DWRITE_TEXT_ALIGNMENT_CENTER
		_, _, _ = comCall(tf, slotTFSetParagraphAlignment, 2) // DWRITE_PARAGRAPH_ALIGNMENT_CENTER
		return tf, nil
	}
	var err error
	if r.badgeFmt, err = mk(600, float32(BadgeFontPx)); err != nil {
		return err
	}
	if r.microFmt, err = mk(500, float32(FontSizeMicroPx)); err != nil {
		return err
	}
	return nil
}

// resize rebinds the DC render target to a new physical size (DPI change /
// state size change).
func (r *renderer) resize(w, h int32) error {
	if w == r.w && h == r.h {
		return nil
	}
	bindRect := rect{l: 0, t: 0, r: w, b: h}
	hr, _, _ := comCall(r.dcRT, slotDCRTBindDC, uintptr(r.hdcMem), unsafePtr(&bindRect))
	if hr != 0 {
		return fmt.Errorf("ball: BindDC(%dx%d) hr=0x%x", w, h, hr)
	}
	r.w, r.h = w, h
	return nil
}

func (r *renderer) release() {
	r.dropAssets()
	comRelease(&r.solid)
	comRelease(&r.badgeFmt)
	comRelease(&r.microFmt)
	comRelease(&r.dcRT)
	if r.hbm != 0 {
		pDeleteObject.Call(r.hbm)
		r.hbm = 0
	}
	if r.hdcMem != 0 {
		pDeleteDC.Call(r.hdcMem)
		r.hdcMem = 0
	}
}

// ---------------------------------------------------------------- gradients

// setAssets (re)builds the cached gradient brushes for one state visual.
// Called on state/size/dpi/theme changes - never per animation frame.
func (r *renderer) setAssets(v Visual) {
	r.dropAssets()
	r.body = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(mulA(pal.OrbBody, v.Opacity))},
		{0.58, colorF(mulA(pal.OrbBody2, v.Opacity))},
		{0.84, colorF(mulA(pal.OrbBody2, 0))},
	}, r.center(), d2d1Point2F{0, 0}, v.SizePx/2, v.SizePx/2)
	// Wensity dark glass body (ticket 65): the reference orb is a near-black
	// translucent sphere, not the white prototype shell (pal.OrbBody is white
	// .30 - that is what made it read as a pale media button). Built from
	// existing C21 dark tokens (BGInset / TintNeutralLo) so no palette/CSS row
	// change is needed; only the glass body is affected.
	if wensityActive() {
		r.body = r.radialBrush([]d2d1GradientStop{
			{0.00, colorF(mulA(pal.BGInset, v.Opacity))},
			{0.50, colorF(mulA(pal.TintNeutralLo, 0.94*v.Opacity))},
			{0.90, colorF(mulA(pal.BGInset, 0.96*v.Opacity))},
			{1.00, colorF(mulA(pal.BGInset, 0))},
		}, r.center(), d2d1Point2F{0, 0}, v.SizePx/2, v.SizePx/2)
	}
	r.hi = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(mulA(lipOf(v), v.Opacity))},
		{0.46, colorF(mulA(lipOf(v), 0))},
	}, d2d1Point2F{float32(r.w) * 0.32, float32(r.h) * 0.24}, d2d1Point2F{0, 0},
		v.SizePx*0.46, v.SizePx*0.46)
	r.core = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(mulA(mulA(v.CoreColor, v.CoreAlpha), v.Opacity))},
		{0.74, colorF(mulA(v.CoreColor, 0))},
	}, d2d1Point2F{r.center().x, r.center().y - v.SizePx*0.06}, d2d1Point2F{0, 0},
		v.SizePx*0.70, v.SizePx*0.70)
	if v.Glass {
		r.setLiquid(v)
		if wensityActive() {
			r.setWensity(v)
		}
	}
	if activeStyle == StyleJarvis {
		r.setJarvis(v)
	}
}

// lipOf is the specular colour: the look's highlight, which reads on top of
// the liquid (pure white would vanish on a light desktop without it).
func lipOf(v Visual) Color {
	if v.Glass {
		return activeLook.Hi
	}
	return pal.OrbHi
}

// setLiquid bakes the shared liquid material: the contact shadow always,
// and the three colour blobs only for the blob treatments - the wensity
// band look (ticket 65) replaces the blobs with its pre-baked band
// (setWensity), so its bake is skipped here on purpose.
func (r *renderer) setLiquid(v Visual) {
	lk := activeLook
	c := r.center()
	R := v.SizePx / 2
	// Contact shadow: what makes a light-coloured orb read on a white
	// desktop. Sits under everything else, offset downward.
	deep := mulA(lk.Caustic, v.Opacity)
	r.caustic = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(deep)},
		{0.55, colorF(mulA(lk.Caustic, GlassCaustic*v.Opacity))},
		{1.00, colorF(mulA(lk.Caustic, 0))},
	}, d2d1Point2F{c.x, c.y + R*0.42}, d2d1Point2F{0, 0}, R*1.25, R*0.85)
	if wensityActive() {
		return
	}
	colors := [3]Color{
		mulA(lk.BlobA, v.Opacity),
		mulA(lk.BlobB, v.Opacity),
		mulA(v.CoreColor, v.Opacity),
	}
	if v.CoreAlpha <= 0 {
		colors[2] = mulA(lk.BlobC, v.Opacity)
	}
	for bi := 0; bi < 3; bi++ {
		col := colors[bi]
		br := liqRadius[bi] * R
		or := liqOffset[bi] * R
		for i := 0; i < liquidLadder; i++ {
			a := float64(i) / liquidLadder * 2 * math.Pi
			p := d2d1Point2F{c.x + or*float32(math.Cos(a)), c.y + or*float32(math.Sin(a))}
			r.liq[bi].centers[i] = p
			r.liq[bi].brushes[i] = r.radialBrush([]d2d1GradientStop{
				{0.00, colorF(col)},
				{0.42, colorF(mulA(col, 0.72))},
				{0.78, colorF(mulA(col, 0.18))},
				{1.00, colorF(mulA(col, 0))},
			}, p, d2d1Point2F{0, 0}, br, br)
		}
	}
}

// liquidIndex picks the ladder step of one blob at the current liquid angle.
// The angle is accumulated by the window (audio-driven), never by a clock in
// the renderer.
func liquidIndex(v Visual, bi int) int {
	a := v.LiquidAngle*liqRate[bi] + liqPhase[bi]
	a = math.Mod(a, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}
	return int(a/(2*math.Pi)*liquidLadder) % liquidLadder
}

func (r *renderer) dropAssets() {
	comRelease(&r.body)
	comRelease(&r.hi)
	comRelease(&r.core)
	comRelease(&r.caustic)
	comRelease(&r.bandCore)
	comRelease(&r.wensityGlow)
	for i := 0; i < wensityBandSegs; i++ {
		comRelease(&r.band[i])
	}
	r.dropJarvis()
	for bi := 0; bi < 3; bi++ {
		for i := 0; i < liquidLadder; i++ {
			comRelease(&r.liq[bi].brushes[i])
		}
	}
	for _, b := range r.glow {
		comRelease(&b)
	}
	r.glow = map[int]unsafe.Pointer{}
}

// glowBrush returns (lazily creating) the outer-glow radial gradient at the
// requested peak alpha (permille 0..1000). Animated states step through a
// small alpha ladder instead of recreating a brush per frame.
// haloExtents returns the (fill, falloff) radii of the outer glow. The window
// is a square (edge = orb + 2*RingMarginPx) while the halo fills 1.5*R out to
// a 1.9*R falloff, so a full-size body gets its halo sliced by the rect and
// reads as a framed square instead of a ball (owner 09-30: 「四周边边角角的，
// 有边框」). Both radii clamp to the inscribed circle together, which leaves
// the recorded Sleeping numbers untouched (its 1.9*R already fits) and makes
// the clipped states end at zero alpha.
func (r *renderer) haloExtents(R float32) (float32, float32) {
	fill, grad := R*1.5, R*1.9
	if fit := float32(r.w)/2 - 1; grad > fit {
		grad = fit
	}
	if fit := float32(r.h)/2 - 1; grad > fit {
		grad = fit
	}
	if fill > grad {
		fill = grad
	}
	return fill, grad
}

func (r *renderer) glowBrush(v Visual, permille int, radius float32) unsafe.Pointer {
	if b, ok := r.glow[permille]; ok {
		return b
	}
	g := mulA(v.GlowColor, float32(permille)/1000*v.Opacity)
	b := r.radialBrush([]d2d1GradientStop{
		{0.45, colorF(g)},
		{1.00, colorF(mulA(v.GlowColor, 0))},
	}, r.center(), d2d1Point2F{0, 0}, radius, radius)
	r.glow[permille] = b
	return b
}

func (r *renderer) radialBrush(stops []d2d1GradientStop, center, offset d2d1Point2F, rx, ry float32) unsafe.Pointer {
	if r.dcRT == nil || len(stops) == 0 {
		return nil
	}
	var coll unsafe.Pointer
	hr, _, _ := comCall(r.dcRT, slotRTCreateGradientStops,
		unsafePtr(&stops[0]), uintptr(len(stops)), d2d1Gamma22, d2d1ExtendModeClamp, unsafePtr(&coll))
	if hr != 0 {
		return nil
	}
	defer comRelease(&coll)
	props := d2d1RadialGradientProps{center: center, originOffset: offset, radiusX: rx, radiusY: ry}
	brushProps := identityBrushProps()
	var brush unsafe.Pointer
	hr, _, _ = comCall(r.dcRT, slotRTCreateRadialGradient,
		unsafePtr(&props), unsafePtr(&brushProps), uintptr(coll), unsafePtr(&brush))
	if hr != 0 {
		return nil
	}
	return brush
}

func (r *renderer) center() d2d1Point2F {
	return d2d1Point2F{float32(r.w) / 2, float32(r.h) / 2}
}

// setSolid colors the shared solid brush (SetColor - no allocation).
func (r *renderer) setSolid(c Color) {
	col := colorF(c)
	_, _, _ = comCall(r.solid, slotSolidSetColor, unsafePtr(&col))
}

// ------------------------------------------------------------ frame drawing

// drawFrame paints one full frame of the visual. phase in [0,1) is the
// animation phase (waves/sweep/pulse); unused by static states. All colors
// come from the C21 palette via Visual - zero hardcoded values here.
func (r *renderer) drawFrame(v Visual, phase float64) {
	if r.dcRT == nil {
		return
	}
	_, _, _ = comCall(r.dcRT, slotRTBeginDraw)
	defer func() {
		hr2, _, _ := comCall(r.dcRT, slotRTEndDraw, 0, 0)
		if hr2 != 0 {
			slog.Error("ball: EndDraw failed", "hr", fmt.Sprintf("0x%x", hr2))
		}
	}()

	transparent := d2d1ColorF{}
	_, _, _ = comCall(r.dcRT, slotRTClear, unsafePtr(&transparent))

	c := r.center()
	R := v.SizePx / 2
	s := float32(r.dpi) / 96 // token px -> physical px scale

	if r.overlayOnly {
		// The orb window owns the picture; this window draws the hit-test
		// silhouette plus the host-only overlays further down. The
		// silhouette's radius is HitTest's clickable circle verbatim (the
		// state body + click tolerance, extended to the 0.75-edge halo law):
		// every pixel the user can see and click must carry nonzero alpha or
		// the OS skips this window before WM_NCHITTEST is consulted. Alpha
		// 1/255 is invisible (0.4% blend) but nonzero, which is what the
		// per-pixel hit test reads. The silhouette's alpha byte must be
		// EXACTLY 1 - WithAlpha replaces the token alpha outright, because
		// mulA(OrbBody, 1/255) would scale 0.62 down to a byte that rounds
		// to zero, and an all-zero-alpha layered window is skipped by the
		// OS hit test entirely (the measured ticket-65 regression).
		r.setSolid(pal.OrbBody.WithAlpha(1.0 / 255.0))
		orbPx := float64(v.SizePx) * float64(s)
		hr := orbPx/2 + float64(ClickTolerancePx)*float64(s)
		if e := math.Ceil(orbPx * 0.75); e > hr {
			hr = e
		}
		r.fillEllipse(c, float32(hr), float32(hr), r.solid)
	} else if v.Glass {
		// The style family picks the material: glass draws the ticket-62/68
		// glass body (with the wensity band as its default treatment), jarvis
		// draws the ticket-65 particle orb (jarvis_windows.go).
		if activeStyle == StyleJarvis {
			r.drawJarvis(v, c, R, s, phase)
		} else {
			r.drawGlass(v, c, R, s)
		}
	} else {
		// Legacy flat body (kept for non-glass visuals).
		hr, hg := r.haloExtents(R)
		if g := r.glowBrush(v, 850, hg); g != nil {
			r.fillEllipse(c, hr, hr, g)
		}
		r.fillEllipse(c, R, R, r.body)
		r.fillEllipse(d2d1Point2F{float32(r.w) * 0.32, float32(r.h) * 0.24},
			v.SizePx*0.46, v.SizePx*0.46, r.hi)
		if v.CoreAlpha > 0 && r.core != nil {
			r.fillEllipse(d2d1Point2F{c.x, c.y - v.SizePx*0.06},
				v.SizePx*0.70, v.SizePx*0.70, r.core)
		}
		r.setSolid(mulA(pal.OrbRim, v.Opacity))
		r.circle(c, R-0.75*s, 1*s)
	}

	// 6. State ring. (Skipped in orb mode: the page renders the state; a D2D
	// ring around the web orb would draw the state twice.)
	if !r.overlayOnly && v.RingColor.A > 0 {
		ringAlpha := v.Opacity
		if v.RingPulse {
			// Danger pulse: 2s cycle, alpha swells and decays (ball.html).
			ringAlpha *= 0.55 + 0.45*float32(0.5+0.5*math.Cos(phase*2*math.Pi))
		}
		r.setSolid(mulA(v.RingColor, ringAlpha))
		r.circle(c, R+5*s, v.RingWidth*s)
	}

	// 7. Listening waves: three staggered expanding rings. (Orb mode: the
	// page's audio bands are the audio face.)
	if !r.overlayOnly && v.Icon == IconAudioLines {
		for i := 0; i < 3; i++ {
			t := math.Mod(phase+float64(i)/3, 1)
			alpha := v.Opacity * float32(1-t) * 0.8
			if alpha > 0 {
				r.setSolid(mulA(pal.AccentLine, alpha))
				r.circle(c, R+5*s+float32(t)*16*s, RingStrokePx*s)
			}
		}
	}

	// 8. Thinking light band: 2px accent band flowing at the orb bottom.
	// (Orb mode: the page's thinking state is the face.)
	if !r.overlayOnly && v.Sweep {
		bw := v.SizePx * 0.4
		x := -bw + float32(math.Mod(phase, 1))*(v.SizePx+bw)
		r.setSolid(mulA(pal.Accent, 0.9*v.Opacity))
		r.fillRect(d2d1RectF{
			l: c.x - R + x, t: c.y + R*0.55,
			r: c.x - R + x + bw, b: c.y + R*0.55 + ThinkingBandPx*s,
		})
	}

	// 9. Progress ring (Downloading).
	if v.Progress >= 0 {
		r.setSolid(mulA(pal.BorderSoft, v.Opacity))
		r.circle(c, R+5*s, 2*s)
		if v.Progress > 0 {
			r.setSolid(mulA(pal.Accent, v.Opacity))
			r.arc(c, R+5*s, -90, -90+360*v.Progress, 2*s)
		}
	}

	// 10. Icon line-art. The wensity glass look has NO glyphs (the reference is
	// a minimal dark orb); state reads through the ring and band behaviour, not
	// the prototype's line icons (triangle/bars/mic - what made it look like a
	// media button). Blob and generic styles keep their icons. (Orb mode: the
	// page carries the face; no icon.)
	if !r.overlayOnly && !(v.Glass && wensityActive()) {
		r.drawIcon(v, c, R, s)
	}

	// 11. Approval badge (top-right, danger disc + white digit).
	if v.BadgeCount > 0 {
		bx, by := c.x+R*0.82, c.y-R*0.82
		br := BadgeDiameterPx * s / 2
		r.setSolid(mulA(pal.BGInset, v.Opacity))
		r.fillEllipse(d2d1Point2F{bx, by}, br+BadgeBorderPx*s, br+BadgeBorderPx*s, r.solid)
		r.setSolid(mulA(pal.Danger, v.Opacity))
		r.fillEllipse(d2d1Point2F{bx, by}, br, br, r.solid)
		r.setSolid(mulA(pal.OnSolid, v.Opacity))
		r.drawText(fmt.Sprintf("%d", v.BadgeCount), r.badgeFmt,
			d2d1RectF{bx - br, by - br, bx + br, by + br})
	}

	// 12. Queued info dot (bottom-left).
	if v.QueueDot {
		dx, dy := c.x-R*0.8, c.y+R*0.8
		dr := QueueDotPx * s / 2
		r.setSolid(mulA(pal.BGInset, v.Opacity))
		r.fillEllipse(d2d1Point2F{dx, dy}, dr+BadgeBorderPx*s, dr+BadgeBorderPx*s, r.solid)
		r.setSolid(mulA(pal.Info, v.Opacity))
		r.fillEllipse(d2d1Point2F{dx, dy}, dr, dr, r.solid)
	}

	// 13. Center text override (L1 countdown seconds / download percent).
	if v.BadgeText != "" {
		r.setSolid(mulA(pal.OnSolid, v.Opacity))
		r.drawText(v.BadgeText, r.microFmt, d2d1RectF{c.x - R, c.y - R, c.x + R, c.y + R})
	}
}

// drawGlass paints the liquid-glass body: contact shadow, halo, shell, the
// three flowing colour blobs, far-side refraction, state core, specular
// highlight, the dark rim + bright lip that make it read on a light desktop,
// and the "not speaking" border at whatever stage its transition reached.
//
// Everything here is a cached brush fill or a polyline stroke: no COM object
// is created, so an animated frame costs the same as a static one.
func (r *renderer) drawGlass(v Visual, c d2d1Point2F, R, s float32) {
	lk := activeLook

	// The edge dock squeezes the whole body along the dock axis (dock.go's
	// law: the tab is the squashed orb tucked against the boundary, and
	// DockPos places the window with the same DockSquash factor). Before the
	// R15#4 re-dispatch only the liquid/highlight obeyed it and the shell
	// stayed round, so a landed tab leaked ~82% of the rest body onto the
	// screen instead of keeping DockOverlapFrac (R15#4/Q-4).
	dockX, dockY := float32(1), float32(1)
	switch v.Dock {
	case EdgeLeft, EdgeRight:
		dockX = DockSquash(v.DockProgress)
	case EdgeTop, EdgeBottom:
		dockY = DockSquash(v.DockProgress)
	case EdgeNone:
	}

	// 0. Contact shadow / caustic - the separation from a white wallpaper.
	if r.caustic != nil {
		r.fillEllipse(d2d1Point2F{c.x, c.y + R*0.30}, R*1.25, R*0.85, r.caustic)
	}

	// 1. Outer halo. Speaking/Listening breathe it (the phase comes from the
	// window's own accumulator); other states hold a static one.
	permille := 850
	if v.Icon == IconVolume {
		permille = 775 + int(225*0.5*(1+math.Cos(v.LiquidAngle*0.5))) // 550..1000
	}
	permille += int(150 * v.LiquidLevel)
	if permille > 1000 {
		permille = 1000
	}
	// The wensity look is a DARK sphere (reference sits on near-black bg); the
	// bright prototype halo (850) reads as a pale disc washing the orb out.
	// Keep only a faint lift for separation from the desktop.
	if wensityActive() {
		permille = 320
	}
	// Snap onto the cached 25-permille ladder (<=18 variants over the whole
	// audio-driven range): glowBrush creates a missing key on demand, and an
	// animated frame must not COM-create brushes (D32 CPU discipline) - after
	// the first visit to each rung every frame is a cache hit.
	permille -= permille % 25
	fillR, gradR := r.haloExtents(R)
	if g := r.glowBrush(v, permille, gradR); g != nil {
		r.fillEllipse(c, fillR, fillR, g)
	}

	// 2. Glass shell (cool translucent body) - squashed with the body.
	r.fillEllipse(c, R*dockX, R*dockY, r.body)

	// 3. The liquid. The wensity band look (ticket 65) draws its horizontal
	// rainbow band here instead of the three blobs; the blob treatments keep
	// the ticket-62 field exactly as delivered. In the blob branch a louder
	// envelope shrinks the visible patch (the surface "gathers") and the
	// summon burst widens it (flow). Both gains are C21 tokens
	// (LiquidGatherPerLevel / SummonFlowSpread) - ticket 74 lifted them out of
	// here, where they were bare literals no table row named.
	if wensityActive() {
		r.drawWensityBand(v, c, R, dockX, dockY)
	} else {
		shrink := 1 - LiquidGatherPerLevel*v.LiquidLevel
		flow := 1 + SummonFlowSpread*v.SummonFlow
		for bi := 0; bi < 3; bi++ {
			lay := &r.liq[bi]
			i := liquidIndex(v, bi)
			if b := lay.brushes[i]; b != nil {
				rr := R * liqRadius[bi] * shrink * flow
				r.fillEllipse(lay.centers[i], rr*dockX, rr*dockY, b)
			}
		}
	}

	// 4. Far-side refraction: the deep tone gathering at the bottom edge.
	if r.core != nil && v.CoreAlpha > 0 {
		r.fillEllipse(d2d1Point2F{c.x, c.y - v.SizePx*0.06},
			v.SizePx*0.70*dockX, v.SizePx*0.70*dockY, r.core)
	}

	// 5. Specular highlight (upper-left), the glass cue. The wensity look omits
	// this big white blob (the reference carries a thin TOP ARC instead, laid by
	// drawWensityEdges); a large pale specular is what washed the ball out.
	if !wensityActive() {
		r.fillEllipse(d2d1Point2F{c.x - R*0.30, c.y - R*0.38},
			R*0.46*dockX, R*0.36*dockY, r.hi)
	}

	// 6. Rim: dark edge light (contrast on white) + bright inner lip. The
	// rings follow the dock squash with the shell they outline.
	r.setSolid(mulA(lk.Rim, v.Opacity))
	r.arcEx(c, (R-0.8*s)*dockX, (R-0.8*s)*dockY, 0, 360, GlassRimPx*s)
	r.setSolid(mulA(lk.Lip, 0.85*v.Opacity))
	r.arcEx(c, (R-2.6*s)*dockX, (R-2.6*s)*dockY, 0, 360, GlassLipPx*s)

	// 6b. Wensity asymmetric refraction (ticket 65): the bright lip lives
	// only on the bottom-right rim (~0.97R, the reference's signature) - the
	// other three directions stay the dark rim step 6 just laid with lk.Rim -
	// plus the top white arc highlight at 0.95R.
	if wensityActive() {
		r.drawWensityEdges(v, c, R, s, dockX, dockY)
	}

	// 7. The border that fades in when nobody speaks (owner's "丝滑引入边框").
	if v.BorderAlpha > 0 {
		r.setSolid(mulA(lk.Lip, v.BorderAlpha*v.Opacity))
		r.arcEx(c, (R+3.5*s)*dockX, (R+3.5*s)*dockY, 0, 360, BorderRingPx*s)
		r.setSolid(mulA(lk.Rim, 0.55*v.BorderAlpha*v.Opacity))
		r.arcEx(c, (R+5.8*s)*dockX, (R+5.8*s)*dockY, 0, 360, 1*s)
	}
}

func (r *renderer) fillEllipse(c d2d1Point2F, rx, ry float32, brush unsafe.Pointer) {
	if brush == nil {
		return
	}
	e := d2d1Ellipse{center: c, radiusX: rx, radiusY: ry}
	_, _, _ = comCall(r.dcRT, slotRTFillEllipse, unsafePtr(&e), uintptr(brush))
}

func (r *renderer) fillRect(rc d2d1RectF) {
	_, _, _ = comCall(r.dcRT, slotRTFillRectangle, unsafePtr(&rc), uintptr(r.solid))
}

// circle draws a stroked circle as a polyline (DrawLine's strokeWidth is a
// stack-slot float = ABI-safe; DrawEllipse's is NOT, see the ABI note in
// d2d_windows.go).
func (r *renderer) circle(c d2d1Point2F, radius, width float32) {
	r.arcEx(c, radius, radius, 0, 360, width)
}

// arc draws a stroked circular arc (degrees) as a polyline.
func (r *renderer) arc(c d2d1Point2F, radius float32, startDeg, endDeg, width float32) {
	r.arcEx(c, radius, radius, startDeg, endDeg, width)
}

// arcEx draws a stroked ellipse arc (degrees) as a polyline, with per-axis
// radii so the dock-squashed body's rings can follow the shell they outline
// (same polyline law as arc; still ABI-safe, geometry lives in registers).
func (r *renderer) arcEx(c d2d1Point2F, rx, ry float32, startDeg, endDeg, width float32) {
	const segs = 48
	if endDeg < startDeg {
		startDeg, endDeg = endDeg, startDeg
	}
	span := endDeg - startDeg
	n := int(span/360*segs) + 2
	prev := r.pointOnEx(c, rx, ry, startDeg)
	for i := 1; i <= n; i++ {
		a := startDeg + span*float32(i)/float32(n)
		next := r.pointOnEx(c, rx, ry, a)
		r.line(prev, next, width)
		prev = next
	}
}

func (r *renderer) pointOn(c d2d1Point2F, radius, deg float32) d2d1Point2F {
	return r.pointOnEx(c, radius, radius, deg)
}

func (r *renderer) pointOnEx(c d2d1Point2F, rx, ry, deg float32) d2d1Point2F {
	rad := deg * float32(math.Pi) / 180
	return d2d1Point2F{
		c.x + rx*float32(math.Cos(float64(rad))),
		c.y + ry*float32(math.Sin(float64(rad))),
	}
}

// line: DrawLine takes D2D1_POINT_2F BY VALUE - two floats packed into one
// register word (x in the low half, y in the high half), per the Win64 ABI.
func (r *renderer) line(p0, p1 d2d1Point2F, width float32) {
	p0v := packedPoint(p0)
	p1v := packedPoint(p1)
	_, _, _ = comCall(r.dcRT, slotRTDrawLine,
		p0v, p1v, uintptr(r.solid), f32bits(width), 0)
}

// packedPoint packs a D2D1_POINT_2F for a by-value register slot.
func packedPoint(p d2d1Point2F) uintptr {
	return uintptr(math.Float32bits(p.x)) | uintptr(math.Float32bits(p.y))<<32
}

func (r *renderer) drawText(s string, fmtPtr unsafe.Pointer, rc d2d1RectF) {
	if fmtPtr == nil {
		return
	}
	u16, _ := windows.UTF16FromString(s)
	_, _, _ = comCall(r.dcRT, slotRTDrawText,
		unsafePtr(&u16[0]), uintptr(utf16Len(s)), uintptr(fmtPtr), unsafePtr(&rc),
		uintptr(r.solid), 0, 0)
}

// drawIcon renders the Lucide-style line glyphs (stroke IconStrokePx, color
// OnSolid - or fg-secondary for the Muted slash per ball.html). Primitive
// shapes only: lines, circles, ellipse fills.
func (r *renderer) drawIcon(v Visual, c d2d1Point2F, R, s float32) {
	g := IconStrokePx * s
	iconColor := mulA(pal.OnSolid, v.Opacity)
	switch v.Icon {
	case IconNone:
		return
	case IconSlash:
		// SPEC-08 §2.1: 1.5px slash. Color fg-secondary (ball.html).
		r.setSolid(mulA(pal.FgSecondary, v.Opacity))
		r.line(d2d1Point2F{c.x - R*0.45, c.y + R*0.45}, d2d1Point2F{c.x + R*0.45, c.y - R*0.45}, SlashStrokePx*s)
		return
	case IconX:
		r.setSolid(iconColor)
		r.line(d2d1Point2F{c.x - R*0.35, c.y - R*0.35}, d2d1Point2F{c.x + R*0.35, c.y + R*0.35}, g)
		r.line(d2d1Point2F{c.x + R*0.35, c.y - R*0.35}, d2d1Point2F{c.x - R*0.35, c.y + R*0.35}, g)
	case IconAudioLines:
		r.setSolid(iconColor)
		for i := -1; i <= 1; i++ {
			h := R * (0.28 + 0.18*float32(2-absi(i))) // center bar tallest
			x := c.x + float32(i)*R*0.28
			r.line(d2d1Point2F{x, c.y - h}, d2d1Point2F{x, c.y + h}, g)
		}
	case IconVolume:
		r.setSolid(iconColor)
		// speaker box + one sound arc
		r.line(d2d1Point2F{c.x - R*0.30, c.y - R*0.14}, d2d1Point2F{c.x - R*0.30, c.y + R*0.14}, g)
		r.line(d2d1Point2F{c.x - R*0.30, c.y - R*0.14}, d2d1Point2F{c.x - R*0.05, c.y - R*0.32}, g)
		r.line(d2d1Point2F{c.x - R*0.30, c.y + R*0.14}, d2d1Point2F{c.x - R*0.05, c.y + R*0.32}, g)
		r.line(d2d1Point2F{c.x - R*0.05, c.y - R*0.32}, d2d1Point2F{c.x - R*0.05, c.y + R*0.32}, g)
		r.arc(c, R*0.42, -50, 50, g)
	case IconMic:
		r.setSolid(iconColor)
		// capsule + stand
		r.circle(d2d1Point2F{c.x, c.y - R*0.12}, R*0.22, g)
		r.line(d2d1Point2F{c.x - R*0.22, c.y - R*0.12}, d2d1Point2F{c.x - R*0.22, c.y + R*0.05}, g)
		r.line(d2d1Point2F{c.x + R*0.22, c.y - R*0.12}, d2d1Point2F{c.x + R*0.22, c.y + R*0.05}, g)
		r.arc(d2d1Point2F{c.x, c.y + R*0.05}, R*0.22, 0, 180, g)
		r.line(d2d1Point2F{c.x, c.y + R*0.27}, d2d1Point2F{c.x, c.y + R*0.42}, g)
	case IconWifiOff:
		r.setSolid(iconColor)
		r.arc(d2d1Point2F{c.x, c.y + R*0.3}, R*0.5, -135, -45, g)
		r.arc(d2d1Point2F{c.x, c.y + R*0.3}, R*0.3, -125, -55, g)
		r.fillEllipse(d2d1Point2F{c.x, c.y + R*0.32}, 1.5*s, 1.5*s, r.solid)
		r.line(d2d1Point2F{c.x - R*0.45, c.y - R*0.35}, d2d1Point2F{c.x + R*0.45, c.y + R*0.15}, g)
	case IconKey:
		r.setSolid(iconColor)
		r.circle(d2d1Point2F{c.x - R*0.12, c.y - R*0.12}, R*0.22, g)
		r.line(d2d1Point2F{c.x + R*0.04, c.y + R*0.04}, d2d1Point2F{c.x + R*0.38, c.y + R*0.38}, g)
		r.line(d2d1Point2F{c.x + R*0.24, c.y + R*0.24}, d2d1Point2F{c.x + R*0.38, c.y + R*0.10}, g)
	case IconArrowUpRight:
		r.setSolid(iconColor)
		r.line(d2d1Point2F{c.x - R*0.28, c.y + R*0.28}, d2d1Point2F{c.x + R*0.28, c.y - R*0.28}, g)
		r.line(d2d1Point2F{c.x + R*0.02, c.y - R*0.28}, d2d1Point2F{c.x + R*0.28, c.y - R*0.28}, g)
		r.line(d2d1Point2F{c.x + R*0.28, c.y - R*0.02}, d2d1Point2F{c.x + R*0.28, c.y - R*0.28}, g)
	case IconAlertTriangle:
		r.setSolid(iconColor)
		r.line(d2d1Point2F{c.x, c.y - R*0.34}, d2d1Point2F{c.x - R*0.40, c.y + R*0.26}, g)
		r.line(d2d1Point2F{c.x, c.y - R*0.34}, d2d1Point2F{c.x + R*0.40, c.y + R*0.26}, g)
		r.line(d2d1Point2F{c.x - R*0.40, c.y + R*0.26}, d2d1Point2F{c.x + R*0.40, c.y + R*0.26}, g)
		r.line(d2d1Point2F{c.x, c.y - R*0.08}, d2d1Point2F{c.x, c.y + R*0.06}, g)
		r.fillEllipse(d2d1Point2F{c.x, c.y + R*0.18}, 1.2*s, 1.2*s, r.solid)
	case IconRefresh:
		r.setSolid(iconColor)
		r.arc(c, R*0.34, 30, 330, g)
		// arrowhead at the open end (30 degrees)
		tip := r.pointOn(c, R*0.34, 30)
		r.line(d2d1Point2F{tip.x - R*0.14, tip.y - R*0.02}, tip, g)
		r.line(d2d1Point2F{tip.x - R*0.02, tip.y - R*0.14}, tip, g)
	}
}

// maxWindowEdge is the largest window edge this renderer must back (max
// configured orb + ring margins, DPI-scaled).
func maxWindowEdge(dpi uint32) int32 {
	return int32(WindowEdgePx(BallSizeMaxPx, dpi))
}

// ensureDCAndDIB creates the memory DC and the backing DIB at (at least)
// w x h. Existing DC/DIB are reused when they are large enough; the DIB is
// recreated only when the requested size outgrows the allocation (DPI rise).
func (r *renderer) ensureDCAndDIB(hdcScreen uintptr, w, h int32) error {
	if r.hdcMem == 0 {
		r.hdcMem, _, _ = pCreateCompatibleDC.Call(hdcScreen)
		if r.hdcMem == 0 {
			return fmt.Errorf("ball: CreateCompatibleDC failed")
		}
	}
	if r.hbm != 0 && int32(r.capW) >= w && int32(r.capH) >= h {
		return nil // reuse
	}
	// (Re)allocate.
	if r.hbm != 0 {
		pDeleteObject.Call(r.hbm)
		r.hbm = 0
	}
	bmi := bitmapInfo{}
	bmi.header = bitmapInfoHeader{
		size:     uint32(unsafe.Sizeof(bmi.header)),
		width:    w,
		height:   -h, // negative = top-down
		planes:   1,
		bitCount: 32,
	}
	var bits uintptr
	hbm, _, err := pCreateDIBSection.Call(r.hdcMem, unsafePtr(&bmi), dibRGBColors, unsafePtr(&bits), 0, 0)
	if hbm == 0 {
		return fmt.Errorf("ball: CreateDIBSection(%dx%d): %v", w, h, err)
	}
	r.hbm = hbm
	r.bits = bits
	r.capW, r.capH = uint32(w), uint32(h)
	pSelectObject.Call(r.hdcMem, hbm)
	return nil
}
