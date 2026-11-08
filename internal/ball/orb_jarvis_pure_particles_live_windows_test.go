//go:build windows && winlive

package ball

// Ticket 65 pure-particle leg (2026-11-08): the LIVE composited-pixel proof
// that the jarvis orb is PURE PARTICLES - no in-page backdrop of any kind -
// and that the cyan particles stay visible on BOTH a white and a black
// desktop. Run explicitly:
//
//	go test -tags winlive ./internal/ball/ -run TestBallLiveOrbJarvisPureParticles -count=1 -timeout 25m -v
//
// The owner's ruling verbatim (highest priority): no sphere body, no visible
// frame, no backdrop - "直接就是粒子态", particles floating by themselves like
// Jarvis. This REVERSES the previous dark-disc proof
// (TestBallLiveOrbJarvisDarkDiscOverWhiteFloor, deleted in this leg): the
// disc existed to hide the AdditiveBlending wash-out over white; the owner
// rejected the disc itself, so the fix moved into the blending mode (the
// scaffold's wisp-jarvis-normal-blending vite plugin rewrites the upstream's
// ten additive material sites to NormalBlending at build time) and the teeth
// below now prove the DISC IS GONE while the particles SURVIVE WHITE.
//
// The floor is WHITE (0xFFFFFF) - deliberately the owner's defect shape at
// its extreme - and then BLACK, the opposite extreme: "works on any
// background" is the actual Jarvis reading, so the cyan-dominant tooth runs
// twice, once per floor. All samples come from the composited screen (what a
// screenshot shows). Four teeth:
//
//  1. NO DARK DISC: an 8-point sample ring at 42% of the window edge (the
//     disc leg's own ring geometry, where the #050913-family disc used to
//     read luminance ~10-20) must yield at least 4 non-particle samples whose
//     AVERAGE luminance stays ABOVE 200 - i.e. the ring reads the white
//     floor, not a backdrop. Individual particle crossings (a data packet
//     crossing a ring point) are tolerated by the >=4-of-8 gate, exactly as
//     the disc leg tolerated them.
//  2. CORNERS STILL TRANSPARENT (over BOTH floors): the four corners 3 px
//     inside the rect must read the floor within +/-12 per channel. The
//     transparency architecture is untouched: the window's corners keep
//     showing the desktop, which also proves no FULL-window backdrop crept
//     back in.
//  3. CYAN PARTICLES ON WHITE: the cyan-dominant ratio across the orb rect
//     (g-r AND b-r >= jarvisCyanGate, b >= jarvisCyanBright - the upstream
//     cyan palette) must beat jarvisCyanFloor on the WHITE floor. This is
//     the leg's core tooth: under AdditiveBlending the white-floor cyan
//     ratio collapses toward zero (the owner's wash-out); under
//     NormalBlending the particles composite by their own alpha and keep
//     their hue.
//  4. CYAN PARTICLES ON BLACK: the same cyan-dominant tooth re-run with the
//     floor swapped to black. Particle visibility must not lean on either
//     extreme - that is the "particles on any background, like Jarvis"
//     reading, and it pins the palette coverage (cores bright enough) from
//     BOTH sides.
//
// The 2026-11-08 deep-color leg (owner's ruling verbatim: "太淡了，颜色深
// 一点，粒子颜色") added the DEPTH teeth below teeth 3-4: the pale build
// passed all four existence teeth while washing out, so the same white-floor
// cyan population must also read DARK and MORE CHROMATIC, cover > 20% of the
// white-floor rect, and cover > 6% of the black-floor rect.
//
// The 2026-11-08 sharp-fine leg (owner's ruling verbatim: "太糊了，根本不像
// 人家那个demo里那样高清" + "粒子有点粗") RE-CALIBRATED the depth teeth and
// added the FINE-GRAIN teeth, because the shape the teeth sample changed
// physics, not quality: the shipped field went from giant blobs (point floor
// 7.00, ~8 connected components averaging 420.9 px on a static white-floor
// capture) to fine grains (upstream point 1.08, render dpr 3.5 supersample,
// 1652 particles). Measured on the static white-floor captures (116x116,
// finegrain.py, same cyan gate): coarse build share 25.02% / mean row-run
// 10.52px / 8 components at 420.9px; sharp-fine build share 17.35% / mean
// run 2.59px / 175 components at 13.34px. A fine-grain field is
// intrinsically EDGE-DOMINATED: a ~2px dot is mostly partial-alpha blend
// pixels over white, so the population's mean luminance sits above the blob
// build's and its mean chroma below - the old blob-era numbers (lum < 120,
// chroma > 44.2, share > 20%) are not properties of the new shape. The
// re-calibrated teeth still bite their own defects with measured margins:
//   - share > 14% (measured 17.35; the pale deep-less build's 12.28-13.63%
//     sits under the floor; particle loss collapses toward 0);
//   - mean luminance < 135 (measured 122.9; the pale build read 174.0, and
//     an alpha-floor mutation 0.92 -> 0.24 pushes fine-dot blend pixels
//     toward ~230);
//   - mean chroma > 34 (measured 39.7; the saturation mutation 1.9 -> 1.0
//     drops the framebuffer's pure color g-r by ~19%, landing near 32 - red).
// The FINE-GRAIN teeth are new and pin exactly the two words of the owner's
// ruling: mean cyan row-run length must stay SHORT (measured 2.59px static;
// the coarse blob build measured 10.52px - the mutation back to the 7.00
// point floor re-lands there) and the row-run COUNT must stay HIGH (fine
// grains fragment the field; blobs merge it).
//
// Mutation evidence (dispatch tooth 4, RUN 2026-11-08, honest result): the
// blending-mode mutation does NOT turn these teeth red. Setting
// jarvisBlendingPatch back to the additive string, rebuilding and re-running
// measured white-floor cyan 13.63% / black-floor cyan 4.89% - the same range
// as the NormalBlending build (12.28% / 3.54%) - so ALL FOUR teeth stayed
// green under the mutation. Mechanism: for every NON-overlapping particle
// pixel the two blending modes produce the IDENTICAL transparent-canvas
// buffer (rgb = cyan*a, alpha = a against the cleared transparent black), and
// overlap regions are a small share of the sampled rect, so a flat-floor
// hue-ratio metric cannot discriminate the modes. The teeth above therefore
// bite DISC-regressions and particle-loss, not the blending constant:
// tooth 1 separates the deleted disc (measured lum ~17.9 under the disc leg)
// from no-disc (measured 255.0); teeth 3-4 separate "particles gone/washed"
// (near-zero ratio, the old non-dark-tooth blindness fixed by hue) from
// "particles visible on BOTH extremes". NormalBlending is shipped because the
// owner rejected the disc outright and normal (alpha) blending is three.js's
// own default for transparent materials - the mainstream shape for a
// particle system that must composite correctly over an UNKNOWN background -
// not because a flat-floor pixel tooth can pin the constant.
//
// Flake treatment (the file family's established shape): up to three fresh
// ball+WebView2 instances, reaping this process's leaked WebView2 children
// before each rebuild (A1602); only the not-cleared expiry is retryable, hard
// failures are Fatal.

import (
	"fmt"
	"math"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// pureTestFloorKind is a known-color floor window class (the pure-particle
// proof's KNOWN backdrops), registered once per color on first use. The
// white floor is the owner's defect extreme; the black floor is the opposite
// extreme, so particle visibility is proven from both sides.
type pureTestFloorKind struct {
	class string
	color uint32 // COLORREF 0x00BBGGRR
	once  sync.Once
	err   error
}

var (
	pureTestWhiteFloor = &pureTestFloorKind{class: "WispOrbPxWhiteFloor", color: 0x00FFFFFF}
	pureTestBlackFloor = &pureTestFloorKind{class: "WispOrbPxBlackFloor", color: 0x00000000}
)

// pureTestFloorWndProc: the floors never handle anything themselves. One
// shared trivial wndProc serves both classes.
func pureTestFloorWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

var pureTestFloorProcPtr = windows.NewCallback(pureTestFloorWndProc)

// create registers (once) and shows a WS_POPUP floor window of this kind
// covering wr - created and painted on the ui-sta thread via Ball.sta.PostTask
// (a window on a thread nobody pumps never paints). WS_EX_TOPMOST: the orb is
// re-asserted above it.
func (k *pureTestFloorKind) create(wr rect) uintptr {
	k.once.Do(func() {
		brush, _, _ := orbTestCreateSolidBrush.Call(uintptr(k.color))
		wc := wndClassEx{
			size:      uint32(unsafe.Sizeof(wndClassEx{})),
			wndProc:   pureTestFloorProcPtr,
			instance:  windows.Handle(moduleHandle()),
			cursor:    windows.Handle(loadArrowCursor()),
			bkgnd:     windows.Handle(brush),
			className: utf16(k.class),
		}
		if r, _, _ := pRegisterClassExW.Call(unsafePtr(&wc)); r == 0 {
			k.err = fmt.Errorf("RegisterClassExW(%s) failed", k.class)
		}
	})
	if k.err != nil {
		return 0
	}
	hwnd, _, _ := pCreateWindowExW.Call(0x8, // WS_EX_TOPMOST, below the re-asserted orb
		unsafePtr(utf16(k.class)), unsafePtr(utf16("Wisp Orb Pure-Particle Floor")),
		uintptr(wsPopup),
		uintptr(wr.l), uintptr(wr.t), uintptr(wr.width()), uintptr(wr.height()),
		0, 0, uintptr(moduleHandle()), 0)
	if hwnd == 0 {
		return 0
	}
	pShowWindow.Call(hwnd, swShownoactivate)
	return hwnd
}

// The pure-particle proof's ring geometry: the disc leg's own 8-point ring at
// 42% of the window edge - where the deleted disc used to read luminance
// ~10-20, and where the white floor must now read >200.
const (
	jarvisPureRingProbeFrac  = 0.42 // ring radius, fraction of the window edge
	jarvisPureRingProbeCount = 8    // ring points
	jarvisPureRingMinSamples = 4    // at least this many non-particle ring samples required
	jarvisPureRingLumFloor   = 200  // non-particle ring samples must average ABOVE this (white floor shows through)
	jarvisPureBlackBudgetSec = 20   // black-phase poll budget
)

// Ticket 65 deep-color leg (2026-11-08, owner's ruling verbatim: "太淡了，
// 颜色深一点，粒子颜色"): the four teeth above prove EXISTENCE and hue, not
// DEPTH - the pinned pre-leg asset (sha 754037b9...) passed all four while
// its white-floor particle field read pale. Measured on that committed asset
// over a white floor (the owner's defect shape at its extreme): particle
// mean luminance 174.0, mean min(g-r,b-r) chroma 44.2, ball-center luminance
// 172. These constants pin the DEPTH the owner asked for, from the same
// composited-screen samples: the white-floor cyan population must be LARGE
// (wash-out collapses coverage), DARK (mean luminance well under the pale
// baseline) and MORE CHROMATIC than the baseline (the pale wash desaturates
// toward grey); the black-floor population must cover far more than the
// bright-palette build's measured 3.3-3.5%. Separating mutation (measured,
// evidence 65-jarvis-pure-deep-color.md): reverting the build's saturation
// factor 1.9 -> 1.0 drops the field's pure-color chroma (0,88,158) ->
// (0,71,105) in the framebuffer, dragging the composited chroma mean under
// the floor - the chroma tooth goes red while the darkening comes from the
// palette hexes and the alpha/point floors.
const (
	// white floor: cyan-dominant share. RE-CALIBRATED by the sharp-fine leg
	// (2026-11-08): the blob-era 20% floor was calibrated on a field of
	// giant points whose interiors dominate; the fine-grain field measures
	// 17.35% on the static white-floor capture, and the pale deep-less
	// build's 12.28-13.63% sits under the new floor. Particle loss still
	// collapses the share toward 0, which is what this tooth bites.
	jarvisDeepWhiteRatioFloor = 0.14
	// black floor: the deep build's dots cover far more than the
	// bright-palette build's measured 3.3-3.5% share.
	jarvisDeepBlackRatioFloor = 0.06
	// mean luminance of the white-floor cyan population. RE-CALIBRATED by
	// the sharp-fine leg: fine grains are edge-dominated (partial-alpha
	// blends over white), so the same deep color reads brighter than in the
	// blob build - measured 122.9 on the static capture. Still bites: the
	// pale build read 174.0 and an alpha-floor mutation pushes fine-dot
	// blend pixels toward ~230.
	jarvisDeepLumCeiling = 135.0
	// mean min(g-r,b-r) of the white-floor cyan population. RE-CALIBRATED
	// by the sharp-fine leg: the blob-era 44.2 measured a field of
	// interior-dominated points; the fine-grain field measures 39.7. Still
	// bites: the saturation mutation (1.9 -> 1.0) drops the pure color's
	// g-r ~19%, landing near 32 - red.
	jarvisDeepChromaFloor = 34.0
)

// Sharp-fine fine-grain teeth (2026-11-08, owner's ruling verbatim "粒子
// 有点粗" + "太糊了"): the coarse blob build measured a mean cyan row-run of
// 10.52px with only 8 connected components (static white-floor capture); the
// sharp-fine build measures 2.59px mean run and hundreds of fragments. The
// teeth pin that shape: short runs (fine grains, not blobs) and many runs
// (dense grains, not a sparse film). Sampling is 1px-stride rows so a run is
// a real particle cross-section, unlike the 3px-stride share sampler.
const (
	// mean cyan run length across the sampled rows must stay SHORT.
	jarvisFineRunLenCeiling = 6.0
	// total cyan runs across the sampled rows must stay HIGH (a window-full
	// of separate grains, not a merged mass or an empty field).
	jarvisFineRunsFloor = 100
	// rows sampled at 1px stride for the fine-grain metrics.
	jarvisFineRows = 20
)

// pureReport carries one attempt's readings for the failure line.
type pureReport struct {
	corners     [4][3]int32
	ringLum     float64
	ringOK      bool
	whiteCyan   float64
	whiteLum    float64
	whiteChroma float64
	fineRuns    int
	fineRunLen  float64
	blackCyan   float64
	blackLum    float64
	blackChroma float64
	samples     int
}

// TestBallLiveOrbJarvisPureParticlesOverWhiteAndBlackFloors is the live
// proof: the jarvis orb parked over a white floor shows NO dark disc and
// keeps cyan particles, then over a black floor keeps the same cyan
// particles. The style family is a process-wide global, so the selection is
// restored to glass before the test returns.
func TestBallLiveOrbJarvisPureParticlesOverWhiteAndBlackFloors(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	orbTestReapLeakedWebView2(t)

	if !SetStyle(styleJarvisName) {
		t.Fatal("SetStyle(jarvis) refused")
	}
	t.Cleanup(func() {
		SetStyle(styleGlassName)
	})

	const attempts = 3
	var last pureReport
	for attempt := 1; attempt <= attempts; attempt++ {
		b, err := New(Options{
			SizePx:      BallSizeDefaultPx,
			Initial:     statemachine.StateSleeping,
			Store:       store,
			WindowTitle: liveBallTitle(t),
			Hotkeys:     liveHotkeys(),
		})
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		ok, rep := orbPurePixelAttempt(t, b, attempt)
		if ok {
			return
		}
		last = rep
		b.Close()
		orbTestReapLeakedWebView2(t)
		if attempt < attempts {
			t.Logf("attempt %d/%d did not clear all teeth (ring lum %.1f ringOK=%v white cyan %.2f%% lum %.1f chroma %.1f black cyan %.2f%%); rebuilding ball and WebView2 instance",
				attempt, attempts, rep.ringLum, rep.ringOK, rep.whiteCyan*100, rep.whiteLum, rep.whiteChroma, rep.blackCyan*100)
		}
	}
	t.Fatalf("jarvis sharp-fine pure-particle regression after %d attempts: corners %v (want each channel within +/-%d of the floor: still transparent on BOTH floors); "+
		"white-floor ring non-particle average luminance %.1f (want > %d with >=%d/%d non-particle samples: NO dark backdrop may return); "+
		"white-floor cyan-dominant ratio %.2f%% (want > %.2f%% existence and > %.2f%% deep), particle mean luminance %.1f (want < %.0f) "+
		"and mean chroma %.1f (want > %.1f); "+
		"fine grains: %d cyan runs (want >= %d) with mean run length %.2fpx (want < %.1f: the coarse 7.00-floor build measured 10.52px, "+
		"the sharp-fine build 2.59px on the static capture); "+
		"black-floor cyan-dominant ratio %.2f%% (want > %.2f%% existence and > %.2f%% deep: particles visible on any background)",
		attempts, last.corners, jarvisCornerTol,
		last.ringLum, jarvisPureRingLumFloor, jarvisPureRingMinSamples, jarvisPureRingProbeCount,
		last.whiteCyan*100, jarvisCyanFloor*100, jarvisDeepWhiteRatioFloor*100,
		last.whiteLum, jarvisDeepLumCeiling,
		last.whiteChroma, jarvisDeepChromaFloor,
		last.fineRuns, jarvisFineRunsFloor, last.fineRunLen, jarvisFineRunLenCeiling,
		last.blackCyan*100, jarvisCyanFloor*100, jarvisDeepBlackRatioFloor*100)
}

// orbPurePixelAttempt runs one create-and-sample round of the four teeth: it
// waits for the orb to attach, centers the ball over the orb's future rect
// (the ball is the anchor followOrb aligns the orb to), lays the WHITE floor
// directly below the orb and polls the white teeth, then swaps in the BLACK
// floor and polls the black teeth. Fatal paths inside are hard failures, not
// the cold-boot flake - the expiry return is the retryable one. Both floors
// are destroyed here (their pump is this ball's STA thread), before the
// caller can Close the ball.
func orbPurePixelAttempt(t *testing.T, b *Ball, attempt int) (bool, pureReport) {
	t.Cleanup(b.Close)

	deadline := time.Now().Add(15 * time.Second)
	var orbHWND windows.HWND
	for {
		orbHWND = b.DebugOrbHWND()
		if orbHWND != 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt %d: orb window did not attach within 15s (WebView2 runtime missing or Embed refused)", attempt)
		}
		time.Sleep(100 * time.Millisecond)
	}

	mon, ok := monitorRectsForWindow(windows.HWND(b.DebugHWND()))
	if !ok {
		t.Fatal("monitorRectsForWindow(ball) returned no monitor")
	}
	work := mon.Work

	ballHWND := uintptr(b.DebugHWND())
	var bwr rect
	if r, _, _ := pGetWindowRect.Call(ballHWND, unsafePtr(&bwr)); r == 0 {
		t.Fatal("GetWindowRect(ball) failed")
	}
	if r, _, _ := pSetWindowPos.Call(ballHWND, 0,
		uintptr((int32(work.L)+int32(work.R))/2-bwr.width()/2),
		uintptr((int32(work.T)+int32(work.B))/2-bwr.height()/2), 0, 0,
		uintptr(swpNoSize|swpNoZOrder|swpNoActivate)); r == 0 {
		t.Fatal("SetWindowPos(ball) failed")
	}
	b.followOrb()

	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}

	// The corner probes: 3 px inside, the family's established shape. The
	// particles' widest projection is 37.1% of the edge; the corners sit on
	// the diagonal at 66.5% - no particle and (with the disc gone) no painted
	// pixel reaches them.
	corners := [4][2]int32{
		{wr.l + jarvisProbeInset, wr.t + jarvisProbeInset},
		{wr.r - jarvisProbeInset - 1, wr.t + jarvisProbeInset},
		{wr.l + jarvisProbeInset, wr.b - jarvisProbeInset - 1},
		{wr.r - jarvisProbeInset - 1, wr.b - jarvisProbeInset - 1},
	}
	// The no-disc ring: 8 points at 42% of the window edge around the center
	// - the exact geometry where the deleted disc used to read ~10-20.
	cx := (wr.l + wr.r) / 2
	cy := (wr.t + wr.b) / 2
	edge := wr.width()
	ring := make([][2]int32, 0, jarvisPureRingProbeCount)
	for i := 0; i < jarvisPureRingProbeCount; i++ {
		ang := float64(i) * (math.Pi / 4.0)
		ring = append(ring, [2]int32{
			cx + int32(float64(edge)*jarvisPureRingProbeFrac*math.Cos(ang)),
			cy + int32(float64(edge)*jarvisPureRingProbeFrac*math.Sin(ang)),
		})
	}

	destroyFloor := func(hwnd uintptr) {
		done := make(chan struct{})
		b.sta.PostTask(func() {
			pDestroyWindow.Call(hwnd)
			close(done)
		})
		<-done
	}

	// ---- Phase 1: WHITE floor - no disc + transparent corners + cyan.
	white := pureTestWhiteFloor
	whiteCh := make(chan uintptr, 1)
	b.sta.PostTask(func() { whiteCh <- white.create(wr) })
	fhWhite := <-whiteCh
	if fhWhite == 0 {
		t.Fatal("white floor window failed to create")
	}
	pSetWindowPos.Call(fhWhite, uintptr(orbHWND), 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))

	fr, fg, fb := 255, 255, 255
	deadline = time.Now().Add(25 * time.Second)
	var rep pureReport
	for {
		hdc, _, _ := pGetDC.Call(0)
		if hdc == 0 {
			t.Fatal("GetDC(NULL) failed")
		}
		cornersOK := true
		for i, pt := range corners {
			r, g, bl := orbTestPixelAt(hdc, pt[0], pt[1])
			rep.corners[i] = [3]int32{r, g, bl}
			if absi(int(r)-fr) > int(jarvisCornerTol) ||
				absi(int(g)-fg) > int(jarvisCornerTol) ||
				absi(int(bl)-fb) > int(jarvisCornerTol) {
				cornersOK = false
			}
		}
		// Tooth 1: the ring reads the WHITE floor - no dark backdrop.
		nonParticle, lumSum := 0, 0
		for _, pt := range ring {
			r, g, bl := orbTestPixelAt(hdc, pt[0], pt[1])
			if int(g)-int(r) >= jarvisCyanGate &&
				int(bl)-int(r) >= jarvisCyanGate &&
				bl >= jarvisCyanBright {
				continue // a particle crossed this ring point this frame
			}
			nonParticle++
			lumSum += (int(r) + int(g) + int(bl)) / 3
		}
		rep.ringOK = nonParticle >= jarvisPureRingMinSamples
		if rep.ringOK {
			rep.ringLum = float64(lumSum) / float64(nonParticle)
		}
		// Tooth 3 + deep-color teeth + fine-grain teeth: cyan particles keep
		// their hue over WHITE, the field is DEEP (re-calibrated for the
		// fine-grain shape), and the grains are FINE - short runs, many
		// runs, exactly the owner's "粒子有点粗"/"太糊了" inverted.
		rep.whiteCyan, rep.whiteLum, rep.whiteChroma, rep.samples = orbPureParticleStats(hdc, wr)
		rep.fineRuns, rep.fineRunLen = orbPureFineGrainStats(hdc, wr)
		pReleaseDC.Call(0, hdc)
		if cornersOK && rep.ringOK && rep.ringLum > float64(jarvisPureRingLumFloor) &&
			rep.whiteCyan > jarvisCyanFloor &&
			rep.whiteCyan > jarvisDeepWhiteRatioFloor &&
			rep.whiteLum < jarvisDeepLumCeiling &&
			rep.whiteChroma > jarvisDeepChromaFloor &&
			rep.fineRunLen < jarvisFineRunLenCeiling &&
			rep.fineRuns >= jarvisFineRunsFloor {
			t.Logf("attempt %d: white-floor phase: corners %v within +/-%d of WHITE; ring %d/%d non-particle samples averaging lum %.1f (> %d: no backdrop); cyan-dominant %.2f%% (existence > %.2f%%, deep > %.2f%%); particle lum %.1f (< %.0f) chroma %.1f (> %.1f); fine grains: %d runs (want >= %d) mean run %.2fpx (want < %.1f)",
				attempt, rep.corners, jarvisCornerTol, nonParticle, jarvisPureRingProbeCount,
				rep.ringLum, jarvisPureRingLumFloor, rep.whiteCyan*100, jarvisCyanFloor*100, jarvisDeepWhiteRatioFloor*100,
				rep.whiteLum, jarvisDeepLumCeiling, rep.whiteChroma, jarvisDeepChromaFloor,
				rep.fineRuns, jarvisFineRunsFloor, rep.fineRunLen, jarvisFineRunLenCeiling)
			break
		}
		if time.Now().After(deadline) {
			destroyFloor(fhWhite)
			return false, rep
		}
		time.Sleep(300 * time.Millisecond)
	}

	// ---- Phase 2: BLACK floor - the same cyan tooth from the opposite
	// extreme (particles on ANY background, like Jarvis). The black floor
	// slots directly below the orb - above the white one - and the white
	// floor is destroyed.
	black := pureTestBlackFloor
	blackCh := make(chan uintptr, 1)
	b.sta.PostTask(func() { blackCh <- black.create(wr) })
	fhBlack := <-blackCh
	if fhBlack == 0 {
		destroyFloor(fhWhite)
		t.Fatal("black floor window failed to create")
	}
	pSetWindowPos.Call(fhBlack, uintptr(orbHWND), 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	destroyFloor(fhWhite)

	fr, fg, fb = 0, 0, 0
	deadline = time.Now().Add(jarvisPureBlackBudgetSec * time.Second)
	for {
		hdc, _, _ := pGetDC.Call(0)
		if hdc == 0 {
			t.Fatal("GetDC(NULL) failed")
		}
		cornersOK := true
		for i, pt := range corners {
			r, g, bl := orbTestPixelAt(hdc, pt[0], pt[1])
			rep.corners[i] = [3]int32{r, g, bl}
			if absi(int(r)-fr) > int(jarvisCornerTol) ||
				absi(int(g)-fg) > int(jarvisCornerTol) ||
				absi(int(bl)-fb) > int(jarvisCornerTol) {
				cornersOK = false
			}
		}
		// Tooth 4 + deep-color: the same cyan particles over BLACK, now with
		// the deep build's coverage floor.
		rep.blackCyan, rep.blackLum, rep.blackChroma, _ = orbPureParticleStats(hdc, wr)
		pReleaseDC.Call(0, hdc)
		if cornersOK && rep.blackCyan > jarvisCyanFloor && rep.blackCyan > jarvisDeepBlackRatioFloor {
			t.Logf("attempt %d: black-floor phase: corners within +/-%d of BLACK; cyan-dominant %.2f%% (existence > %.2f%%, deep > %.2f%%); particle lum %.1f chroma %.1f",
				attempt, jarvisCornerTol, rep.blackCyan*100, jarvisCyanFloor*100, jarvisDeepBlackRatioFloor*100,
				rep.blackLum, rep.blackChroma)
			destroyFloor(fhBlack)
			return true, rep
		}
		if time.Now().After(deadline) {
			destroyFloor(fhBlack)
			return false, rep
		}
		time.Sleep(300 * time.Millisecond)
	}
}

// orbPureParticleStats samples the orb rect (4 px inset, 3 px stride - the
// file family's established grid) and returns the cyan-dominant share
// (g-r >= jarvisCyanGate AND b-r >= jarvisCyanGate AND b >= jarvisCyanBright),
// that population's mean luminance and mean min(g-r,b-r) chroma (the
// deep-color teeth's dark-and-chromatic ruler), and the total sampled pixel
// count.
func orbPureParticleStats(hdc uintptr, wr rect) (ratio, lumMean, chromaMean float64, total int) {
	cyan, lumSum, chromaSum := 0, 0, 0
	total = 0
	for y := wr.t + 4; y < wr.b-4; y += 3 {
		for x := wr.l + 4; x < wr.r-4; x += 3 {
			r, g, bl := orbTestPixelAt(hdc, x, y)
			total++
			if int(g)-int(r) >= jarvisCyanGate &&
				int(bl)-int(r) >= jarvisCyanGate &&
				bl >= jarvisCyanBright {
				cyan++
				lumSum += (int(r) + int(g) + int(bl))
				cr, cb := int(g)-int(r), int(bl)-int(r)
				if cb < cr {
					cr = cb
				}
				chromaSum += cr
			}
		}
	}
	if total == 0 || cyan == 0 {
		return 0, 0, 0, total
	}
	return float64(cyan) / float64(total),
		float64(lumSum) / float64(cyan) / 3.0,
		float64(chromaSum) / float64(cyan),
		total
}

// orbPureFineGrainStats samples jarvisFineRows evenly spaced rows of the orb
// rect at 1px stride (a run is a real particle cross-section, unlike the
// share sampler's 3px grid) and returns the total number of cyan runs and
// their mean length in px - the sharp-fine teeth's ruler: the coarse 7.00
// point-floor build merges the field into few long runs (static capture:
// 10.52px mean), the sharp-fine build keeps hundreds of short ones (2.59px
// mean). A row with no cyan at all contributes nothing; the runs floor in
// the caller keeps a empty/merged field from passing on a technicality.
func orbPureFineGrainStats(hdc uintptr, wr rect) (runs int, meanRunLen float64) {
	top, bottom := wr.t+8, wr.b-8
	lenSum := 0
	for i := 0; i < jarvisFineRows; i++ {
		y := top + (bottom-top)*int32(i)/int32(jarvisFineRows-1)
		inRun := false
		for x := wr.l + 8; x < wr.r-8; x++ {
			r, g, bl := orbTestPixelAt(hdc, x, y)
			cy := int(g)-int(r) >= jarvisCyanGate &&
				int(bl)-int(r) >= jarvisCyanGate &&
				bl >= jarvisCyanBright
			if cy {
				lenSum++
				if !inRun {
					runs++
					inRun = true
				}
				continue
			}
			inRun = false
		}
	}
	if runs == 0 {
		return 0, 0
	}
	return runs, float64(lenSum) / float64(runs)
}
