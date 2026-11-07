//go:build windows && winlive

package ball

// Ticket 65 jarvis dark-disc leg (2026-11-07): the LIVE composited-pixel proof
// of the in-page dark disc. Run explicitly:
//
//	go test -tags winlive ./internal/ball/ -run TestBallLiveOrbJarvisDarkDisc -count=1 -timeout 25m -v
//
// The owner's ruling: the jarvis particles must composite over a DARK DISC
// inside the page, not over whatever surface sits behind the transparent
// window - over the owner's white windows AdditiveBlending washed the cyan
// ring to grey-white (JARVIS-PRODUCT-REAL.png). The disc therefore lives in
// the page (#jarvis-disc, border-radius 50%, demo dark-blue family) and the
// window corners must KEEP showing the desktop.
//
// The floor is WHITE (0xFFFFFF), deliberately the opposite extreme of the
// existing corners regression's (51,51,50) gray: if the disc is gone, the
// particles composite over white and every tooth below collapses in the
// direction the owner measured on his machine. Three teeth, all sampled from
// the composited screen (what a screenshot shows):
//
//  1. CORNERS STILL TRANSPARENT: the four corners 3 px inside the rect must
//     read the WHITE floor within +/-12 per channel. The disc's circle (96%
//     of the edge) cannot reach them: the probes sit on the diagonal at 66.5%
//     of the edge (70.7% asymptotically), far outside the 48%-of-edge disc
//     radius. This tooth is what keeps the disc from degenerating into a
//     full-window backdrop (which would break the transparency architecture).
//  2. THE DISC IS REALLY THERE AND REALLY DARK: an 8-point sample ring at
//     42% of the window edge (inside the disc's 48% radius with a ~6 px rim
//     margin, outside the particles' 37.1%-of-edge projection) must yield at
//     least 4 non-particle samples whose average luminance stays under 60 - over the
//     #050913-family disc that reads ~10-20; with no disc the same samples
//     read the white floor (>=249) or washed-white particles. Particle
//     contamination of individual ring points (a data packet crossing) is
//     tolerated by the >=4-of-8 gate, not by loosening the luminance floor.
//  3. CYAN PARTICLES KEEP THEIR HUE: the cyan-dominant ratio across the orb
//     rect (g-r AND b-r >= 12, b >= 80 - the upstream cyan palette) must beat
//     jarvisCyanFloor, the same tooth the existing jarvis proof carries - on
//     the dark disc the additive ring stays cyan instead of washing out.
//
// Mutation plan (tooth ④ of the dispatch): removing the #jarvis-disc element
// from the template leaves teeth 1 and 3 untouched (corners stay transparent,
// particles stay cyan-ish over white only near the ring - over white the hue
// washes and tooth 3 usually fails too) but tooth 2 goes red immediately:
// every ring sample reads the white floor. Restore the bytes and it goes green.
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

// discTestFloorColor is the WHITE known floor, as a COLORREF (0x00BBGGRR):
// 0x00FFFFFF = r=255, g=255, b=255. White is the extreme the owner's real
// desktop sits at (his JARVIS-PRODUCT-REAL.png defect shape); if the disc is
// missing, the window reads white through its transparent pixels and the
// teeth collapse exactly as the product defect showed.
const discTestFloorColor = 0x00FFFFFF

// discTestWhiteFloorWndProc: the white floor never handles anything itself.
func discTestWhiteFloorWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

var (
	discTestWhiteFloorProcPtr = windows.NewCallback(discTestWhiteFloorWndProc)
	discTestWhiteFloorOnce    sync.Once
	discTestWhiteFloorErr     error
)

// discTestCreateWhiteFloor registers (once) and shows a white WS_POPUP window
// covering the orb's rect - the dark-disc proof's KNOWN floor. Same creation
// law as orbTestCreateColoredFloor: WS_EX_TOPMOST (the orb is re-asserted
// above it), created and painted on the ui-sta thread via Ball.sta.PostTask -
// a window on a thread nobody pumps never paints.
func discTestCreateWhiteFloor(wr rect) uintptr {
	discTestWhiteFloorOnce.Do(func() {
		brush, _, _ := orbTestCreateSolidBrush.Call(discTestFloorColor)
		wc := wndClassEx{
			size:      uint32(unsafe.Sizeof(wndClassEx{})),
			wndProc:   discTestWhiteFloorProcPtr,
			instance:  windows.Handle(moduleHandle()),
			cursor:    windows.Handle(loadArrowCursor()),
			bkgnd:     windows.Handle(brush),
			className: utf16("WispOrbPxWhiteFloor"),
		}
		if r, _, _ := pRegisterClassExW.Call(unsafePtr(&wc)); r == 0 {
			discTestWhiteFloorErr = fmt.Errorf("RegisterClassExW(white-floor) failed")
		}
	})
	if discTestWhiteFloorErr != nil {
		return 0
	}
	hwnd, _, _ := pCreateWindowExW.Call(0x8, // WS_EX_TOPMOST, below the re-asserted orb
		unsafePtr(utf16("WispOrbPxWhiteFloor")), unsafePtr(utf16("Wisp Orb White Floor")),
		uintptr(wsPopup),
		uintptr(wr.l), uintptr(wr.t), uintptr(wr.width()), uintptr(wr.height()),
		0, 0, uintptr(moduleHandle()), 0)
	if hwnd == 0 {
		return 0
	}
	pShowWindow.Call(hwnd, swShownoactivate)
	return hwnd
}

// TestBallLiveOrbJarvisDarkDiscOverWhiteFloor is the live proof: the jarvis
// orb parked over a white floor keeps its corners transparent AND composites
// its particles over a dark in-page disc. The style family is a process-wide
// global, so the selection is restored to glass before the test returns.
func TestBallLiveOrbJarvisDarkDiscOverWhiteFloor(t *testing.T) {
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
	var lastCorners [4][3]int32
	lastDiscLum, lastCyan := 0.0, 0.0
	lastNonParticle := 0
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
		ok, corners, discLum, cyan, nonParticle := orbDiscPixelAttempt(t, b, attempt)
		if ok {
			return
		}
		lastCorners, lastDiscLum, lastCyan, lastNonParticle = corners, discLum, cyan, nonParticle
		b.Close()
		orbTestReapLeakedWebView2(t)
		if attempt < attempts {
			t.Logf("attempt %d/%d did not clear all three teeth (best disc-ring lum %.1f, cyan %.2f%%, non-particle ring samples %d); rebuilding ball and WebView2 instance",
				attempt, attempts, discLum, cyan*100, lastNonParticle)
		}
	}
	t.Fatalf("jarvis dark-disc regression after %d attempts: corners %v (want each channel within +/-%d of the WHITE floor: still transparent); "+
		"disc-ring non-particle samples %d/%d with average luminance %.1f (want >=%d samples averaging < %d: the dark disc really exists); "+
		"cyan-dominant ratio %.2f%% (want > %.2f%%: particles keep their hue over the disc)",
		attempts, lastCorners, jarvisCornerTol, jarvisDiscMinSamples, jarvisDiscProbeCount,
		lastDiscLum, jarvisDiscMinSamples, jarvisDiscLumFloor, lastCyan*100, jarvisCyanFloor*100)
}

// orbDiscPixelAttempt runs one create-and-sample round of the three teeth: it
// waits for the orb to attach, centers the ball over the orb's future rect
// (the ball is the anchor followOrb aligns the orb to), lays the WHITE floor
// directly below the orb, and polls until all three teeth clear or the budget
// expires. Fatal paths inside are hard failures, not the cold-boot flake - the
// expiry return is the retryable one. The floor is destroyed here (its pump is
// this ball's STA thread), before the caller can Close the ball.
func orbDiscPixelAttempt(t *testing.T, b *Ball, attempt int) (bool, [4][3]int32, float64, float64, int) {
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

	floor := make(chan uintptr, 1)
	b.sta.PostTask(func() { floor <- discTestCreateWhiteFloor(wr) })
	fh := <-floor
	if fh == 0 {
		t.Fatal("white floor window failed to create")
	}
	pSetWindowPos.Call(fh, uintptr(orbHWND), 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	defer func() {
		done := make(chan struct{})
		b.sta.PostTask(func() {
			pDestroyWindow.Call(fh)
			close(done)
		})
		<-done
	}()

	// The corner probes: 3 px inside, the family's established shape.
	corners := [4][2]int32{
		{wr.l + jarvisProbeInset, wr.t + jarvisProbeInset},
		{wr.r - jarvisProbeInset - 1, wr.t + jarvisProbeInset},
		{wr.l + jarvisProbeInset, wr.b - jarvisProbeInset - 1},
		{wr.r - jarvisProbeInset - 1, wr.b - jarvisProbeInset - 1},
	}
	fr := int32(discTestFloorColor & 0xff)
	fg := int32((discTestFloorColor >> 8) & 0xff)
	fb := int32((discTestFloorColor >> 16) & 0xff)

	// The disc-sample ring: 8 points at 42% of the window edge around the
	// center - inside the disc (48% radius, ~6 px rim margin), outside the
	// particle projection (37.1% of the edge along the axes the ring touches).
	cx := (wr.l + wr.r) / 2
	cy := (wr.t + wr.b) / 2
	edge := wr.width()
	ring := make([][2]int32, 0, jarvisDiscProbeCount)
	for i := 0; i < jarvisDiscProbeCount; i++ {
		ang := float64(i) * (math.Pi / 4.0)
		ring = append(ring, [2]int32{
			cx + int32(float64(edge)*jarvisDiscProbeFrac*math.Cos(ang)),
			cy + int32(float64(edge)*jarvisDiscProbeFrac*math.Sin(ang)),
		})
	}

	deadline = time.Now().Add(25 * time.Second)
	var readings [4][3]int32
	for {
		hdc, _, _ := pGetDC.Call(0)
		if hdc == 0 {
			t.Fatal("GetDC(NULL) failed")
		}
		// Tooth 1: corners still transparent over the WHITE floor.
		cornersOK := true
		for i, pt := range corners {
			r, g, bl := orbTestPixelAt(hdc, pt[0], pt[1])
			readings[i] = [3]int32{r, g, bl}
			if absi(int(r)-int(fr)) > int(jarvisCornerTol) ||
				absi(int(g)-int(fg)) > int(jarvisCornerTol) ||
				absi(int(bl)-int(fb)) > int(jarvisCornerTol) {
				cornersOK = false
			}
		}
		// Tooth 2: the dark disc really exists under the particles.
		nonParticle, lumSum := 0, 0
		for _, pt := range ring {
			r, g, bl := orbTestPixelAt(hdc, pt[0], pt[1])
			lum := (int(r) + int(g) + int(bl)) / 3
			if int(g)-int(r) >= jarvisCyanGate &&
				int(bl)-int(r) >= jarvisCyanGate &&
				bl >= jarvisCyanBright {
				continue // a particle crossed this ring point this frame
			}
			nonParticle++
			lumSum += lum
		}
		discLum := 0.0
		if nonParticle > 0 {
			discLum = float64(lumSum) / float64(nonParticle)
		}
		// Tooth 3: cyan particles keep their hue across the orb rect.
		cyan, total := 0, 0
		for y := wr.t + 4; y < wr.b-4; y += 3 {
			for x := wr.l + 4; x < wr.r-4; x += 3 {
				r, g, bl := orbTestPixelAt(hdc, x, y)
				total++
				if int(g)-int(r) >= jarvisCyanGate &&
					int(bl)-int(r) >= jarvisCyanGate &&
					bl >= jarvisCyanBright {
					cyan++
				}
			}
		}
		pReleaseDC.Call(0, hdc)
		ratio := 0.0
		if total > 0 {
			ratio = float64(cyan) / float64(total)
		}
		if cornersOK && nonParticle >= jarvisDiscMinSamples && discLum < float64(jarvisDiscLumFloor) && ratio > jarvisCyanFloor {
			t.Logf("attempt %d: jarvis dark-disc proof: corners %v within +/-%d of the WHITE floor; disc ring %d/%d non-particle samples averaging lum %.1f (< %d); cyan-dominant %d/%d sampled px (%.2f%% > %.2f%%)",
				attempt, readings, jarvisCornerTol, nonParticle, jarvisDiscProbeCount,
				discLum, jarvisDiscLumFloor, cyan, total, ratio*100, jarvisCyanFloor*100)
			return true, readings, discLum, ratio, nonParticle
		}
		if time.Now().After(deadline) {
			return false, readings, discLum, ratio, nonParticle
		}
		time.Sleep(300 * time.Millisecond)
	}
}
