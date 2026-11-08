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
// cyan population must now also read DARK (mean luminance < 120 vs the pale
// baseline's 174.0) and MORE CHROMATIC (mean min(g-r,b-r) > 44.2, the pale
// baseline's measured value), cover > 20% of the white-floor rect, and cover
// > 6% of the black-floor rect; see the constants block for the provenance
// of every number.
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
	// white floor: cyan-dominant share must beat the pure-particle leg's 2%
	// existence floor by an order of magnitude (wash-out collapses coverage;
	// the deep build measures ~2x this floor).
	jarvisDeepWhiteRatioFloor = 0.20
	// black floor: the deep build's dots cover far more than the
	// bright-palette build's measured 3.3-3.5% share.
	jarvisDeepBlackRatioFloor = 0.06
	// mean luminance of the white-floor cyan population: the pale baseline
	// measured 174.0; the deep build must read DARKER than this.
	jarvisDeepLumCeiling = 120.0
	// mean min(g-r,b-r) of the white-floor cyan population must beat the
	// pre-leg asset's measured 44.2.
	jarvisDeepChromaFloor = 44.2
)

// pureReport carries one attempt's readings for the failure line.
type pureReport struct {
	corners     [4][3]int32
	ringLum     float64
	ringOK      bool
	whiteCyan   float64
	whiteLum    float64
	whiteChroma float64
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
	t.Fatalf("jarvis pure-particle deep-color regression after %d attempts: corners %v (want each channel within +/-%d of the floor: still transparent on BOTH floors); "+
		"white-floor ring non-particle average luminance %.1f (want > %d with >=%d/%d non-particle samples: NO dark backdrop may return); "+
		"white-floor cyan-dominant ratio %.2f%% (want > %.2f%% existence and > %.2f%% deep), particle mean luminance %.1f (want < %.0f: DEEPER than the pale 174.0 baseline) "+
		"and mean chroma %.1f (want > %.1f: MORE saturated than the pale 44.2 baseline); "+
		"black-floor cyan-dominant ratio %.2f%% (want > %.2f%% existence and > %.2f%% deep: particles visible on any background)",
		attempts, last.corners, jarvisCornerTol,
		last.ringLum, jarvisPureRingLumFloor, jarvisPureRingMinSamples, jarvisPureRingProbeCount,
		last.whiteCyan*100, jarvisCyanFloor*100, jarvisDeepWhiteRatioFloor*100,
		last.whiteLum, jarvisDeepLumCeiling,
		last.whiteChroma, jarvisDeepChromaFloor,
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
		// Tooth 3 + deep-color teeth: cyan particles keep their hue over
		// WHITE, and the field is DEEP - large, dark, more chromatic than
		// the pale pre-leg baseline.
		rep.whiteCyan, rep.whiteLum, rep.whiteChroma, rep.samples = orbPureParticleStats(hdc, wr)
		pReleaseDC.Call(0, hdc)
		if cornersOK && rep.ringOK && rep.ringLum > float64(jarvisPureRingLumFloor) &&
			rep.whiteCyan > jarvisCyanFloor &&
			rep.whiteCyan > jarvisDeepWhiteRatioFloor &&
			rep.whiteLum < jarvisDeepLumCeiling &&
			rep.whiteChroma > jarvisDeepChromaFloor {
			t.Logf("attempt %d: white-floor phase: corners %v within +/-%d of WHITE; ring %d/%d non-particle samples averaging lum %.1f (> %d: no backdrop); cyan-dominant %.2f%% (existence > %.2f%%, deep > %.2f%%); particle lum %.1f (< %.0f) chroma %.1f (> %.1f)",
				attempt, rep.corners, jarvisCornerTol, nonParticle, jarvisPureRingProbeCount,
				rep.ringLum, jarvisPureRingLumFloor, rep.whiteCyan*100, jarvisCyanFloor*100, jarvisDeepWhiteRatioFloor*100,
				rep.whiteLum, jarvisDeepLumCeiling, rep.whiteChroma, jarvisDeepChromaFloor)
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
