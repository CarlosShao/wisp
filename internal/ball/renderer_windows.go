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
}

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
	comRelease(&r.solid)
	for _, b := range r.glow {
		comRelease(&b)
	}
	r.glow = map[int]unsafe.Pointer{}
	comRelease(&r.core)
	comRelease(&r.hi)
	comRelease(&r.body)
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
	r.hi = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(mulA(pal.OrbHi, v.Opacity))},
		{0.46, colorF(mulA(pal.OrbHi, 0))},
	}, d2d1Point2F{float32(r.w) * 0.32, float32(r.h) * 0.24}, d2d1Point2F{0, 0},
		v.SizePx*0.46, v.SizePx*0.46)
	r.core = r.radialBrush([]d2d1GradientStop{
		{0.00, colorF(mulA(mulA(v.CoreColor, v.CoreAlpha), v.Opacity))},
		{0.74, colorF(mulA(v.CoreColor, 0))},
	}, d2d1Point2F{r.center().x, r.center().y - v.SizePx*0.06}, d2d1Point2F{0, 0},
		v.SizePx*0.70, v.SizePx*0.70)
}

func (r *renderer) dropAssets() {
	comRelease(&r.body)
	comRelease(&r.hi)
	comRelease(&r.core)
	for _, b := range r.glow {
		comRelease(&b)
	}
	r.glow = map[int]unsafe.Pointer{}
}

// glowBrush returns (lazily creating) the outer-glow radial gradient at the
// requested peak alpha (permille 0..1000). Animated states step through a
// small alpha ladder instead of recreating a brush per frame.
func (r *renderer) glowBrush(v Visual, permille int) unsafe.Pointer {
	if b, ok := r.glow[permille]; ok {
		return b
	}
	g := mulA(v.GlowColor, float32(permille)/1000*v.Opacity)
	b := r.radialBrush([]d2d1GradientStop{
		{0.45, colorF(g)},
		{1.00, colorF(mulA(v.GlowColor, 0))},
	}, r.center(), d2d1Point2F{0, 0}, v.SizePx*0.95, v.SizePx*0.95)
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

// mulA scales a color's alpha (straight alpha, D2D does the premultiply).
func mulA(c Color, a float32) Color {
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	return Color{c.R, c.G, c.B, c.A * a}
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

	// 1. Outer glow (box-shadow 0 0 28px in ball.html). Speaking breathes the
	// glow (success 1.6s); every other state holds a static glow.
	permille := 850
	if v.Icon == IconVolume {
		permille = 775 + int(225*0.5*(1+math.Cos(phase*2*math.Pi))) // 550..1000
	}
	if g := r.glowBrush(v, permille); g != nil {
		r.fillEllipse(c, R*1.5, R*1.5, g)
	}

	// 2. Glass body + 3. highlight + 4. state core tint.
	if r.body != nil {
		r.fillEllipse(c, R, R, r.body)
	}
	if r.hi != nil {
		r.fillEllipse(d2d1Point2F{float32(r.w) * 0.32, float32(r.h) * 0.24},
			v.SizePx*0.46, v.SizePx*0.46, r.hi)
	}
	if r.core != nil && v.CoreAlpha > 0 {
		r.fillEllipse(d2d1Point2F{c.x, c.y - v.SizePx*0.06},
			v.SizePx*0.70, v.SizePx*0.70, r.core)
	}

	// 5. Rim (edge light line).
	r.setSolid(mulA(pal.OrbRim, v.Opacity))
	r.circle(c, R-0.75*s, 1*s)

	// 6. State ring.
	if v.RingColor.A > 0 {
		ringAlpha := v.Opacity
		if v.RingPulse {
			// Danger pulse: 2s cycle, alpha swells and decays (ball.html).
			ringAlpha *= 0.55 + 0.45*float32(0.5+0.5*math.Cos(phase*2*math.Pi))
		}
		r.setSolid(mulA(v.RingColor, ringAlpha))
		r.circle(c, R+5*s, v.RingWidth*s)
	}

	// 7. Listening waves: three staggered expanding rings.
	if v.Icon == IconAudioLines {
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
	if v.Sweep {
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

	// 10. Icon line-art.
	r.drawIcon(v, c, R, s)

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
	r.arc(c, radius, 0, 360, width)
}

// arc draws a stroked circular arc (degrees) as a polyline.
func (r *renderer) arc(c d2d1Point2F, radius float32, startDeg, endDeg, width float32) {
	const segs = 48
	if endDeg < startDeg {
		startDeg, endDeg = endDeg, startDeg
	}
	span := endDeg - startDeg
	n := int(span/360*segs) + 2
	prev := r.pointOn(c, radius, startDeg)
	for i := 1; i <= n; i++ {
		a := startDeg + span*float32(i)/float32(n)
		next := r.pointOn(c, radius, a)
		r.line(prev, next, width)
		prev = next
	}
}

func (r *renderer) pointOn(c d2d1Point2F, radius, deg float32) d2d1Point2F {
	rad := deg * float32(math.Pi) / 180
	return d2d1Point2F{c.x + radius*float32(math.Cos(float64(rad))),
		c.y + radius*float32(math.Sin(float64(rad)))}
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
