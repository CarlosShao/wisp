//go:build windows && winlive

package ball

// Ticket 65 second visual style (the jarvis particle page): the LIVE proof
// that the orb host serving the jarvis style really renders the CYAN particle
// orb over a TRANSPARENT window. Run explicitly:
//
//	go test -tags winlive ./internal/ball/ -run TestBallLiveOrbJarvis -count=1 -timeout 20m -v
//
// The teeth (ticket 65 style/bg rework, 2026-11-07): the previous tooth was
// "non-dark ratio > 1%" (any channel > 30 over a black backdrop). That tooth
// is BLIND to the defect this leg exists for: a page that paints the wrong
// BACKGROUND is nearly 100% non-dark and passes. Replacement teeth, all
// sampled from the composited screen (what a screenshot shows):
//
//  1. TRANSPARENT BACKGROUND: the orb is parked over a KNOWN (51,51,50) floor
//     (orbTestCreateColoredFloor; a black floor cannot serve - transparent-
//     over-black and the black-square defect both read (0,0,0)) and the four
//     corners of the window, 3 px inside (the orchestrator's probe shape),
//     must read the floor within +/-12 per channel. The product orb window is
//     a transparent WebView2 overlay (orb_windows.go: DefaultBackgroundColor
//     a=0, no WS_EX_LAYERED; orbasset.go injects background:transparent
//     !important), so the background tooth is the desktop-match form of
//     TestBallLiveOrbCornersReadTheDesktop - NOT an absolute darkness
//     threshold, which would break the transparent-overlay architecture.
//     A page that paints its own background (white OR black) fails here.
//  2. CYAN PARTICLES: across the orb rect, pixels must be CYAN-DOMINANT -
//     g-r >= jarvisCyanGate AND b-r >= jarvisCyanGate AND b >= jarvisCyanBright
//     (the upstream cyan palette is primary 0x38f4ff / secondary 0x19a8ff /
//     tertiary 0x7fffee: every particle color has green and blue far above
//     red; the additive-blended contribution over the gray floor keeps that
//     dominance). The cyan-dominant share of sampled pixels must exceed
//     jarvisCyanFloor. This tooth catches the washed-out white-ring defect
//     (over a white page the ring reads gray-white: g-r ~ 0) and a recolored
//     orb (aurora/ember hues fail the g/b>r test).
//  3. The old "non-dark > 1%" assertion is GONE from this test: on a page
//     with a wrong (light) background it is trivially always true, and it
//     says nothing about hue. Cyan-dominance at ratio > floor already proves
//     the page is alive AND the right color; corners +/-12 prove the
//     background is the desktop's, not the page's.
//
// Flake treatment (the file family's established shape): up to three fresh
// ball+WebView2 instances, reaping this process's leaked WebView2 children
// before each rebuild (the measured accumulation mechanism, A1602); only the
// dark/no-cyan expiry is retryable, hard failures are Fatal.

import (
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// The jarvis style's probe constants (self-contained: the corners regression's
// file-scope constants belong to orb_live_windows_test.go's evolution, and
// this test pins its OWN teeth so the two can move independently).
const (
	jarvisProbeInset = int32(3)  // corners sampled this far inside the rect
	jarvisCornerTol  = int32(12) // per-channel distance from the floor color
	jarvisCyanGate   = 12        // g-r AND b-r must both exceed this (cyan hue)
	jarvisCyanBright = 80        // blue floor: dark sampling noise never counts
	jarvisCyanFloor  = 0.02      // cyan-dominant share of sampled px must beat this
	// (measured 2026-11-07 on this box: 7.91% = 76/961 sampled px on a healthy
	// boot; the 2% floor keeps ~4x headroom under that reading while a washed-
	// out gray/white ring or a recolored page lands near zero.)
)

// TestBallLiveOrbJarvisStyleRendersParticles: select the jarvis style, create
// the real ball + orb window, and require the teeth above to clear. The style
// family is a process-wide global, so the selection is restored to glass
// before this test returns - a later test in the same process runs its glass
// contract unchanged.
func TestBallLiveOrbJarvisStyleRendersParticles(t *testing.T) {
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
	lastCyan := 0.0
	lastSamples := 0
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
		ok, corners, cyan, samples := orbJarvisPixelAttempt(t, b, attempt)
		if ok {
			testJarvisRuntimeStyleSwitch(t, b)
			return
		}
		lastCorners, lastCyan, lastSamples = corners, cyan, samples
		b.Close()
		orbTestReapLeakedWebView2(t)
		if attempt < attempts {
			t.Logf("attempt %d/%d did not clear both teeth (best cyan ratio %.3f%% over %d sampled px); rebuilding ball and WebView2 instance",
				attempt, attempts, cyan*100, samples)
		}
	}
	t.Fatalf("jarvis orb style regression after %d attempts: corners %v (want each channel within +/-%d of the (51,51,50) floor: transparent background); cyan-dominant ratio %.3f%% over %d sampled px (want > %.2f%%: g-r>=%d AND b-r>=%d AND b>=%d, the upstream cyan palette's hue)",
		attempts, lastCorners, jarvisCornerTol, lastCyan*100, lastSamples,
		jarvisCyanFloor*100, jarvisCyanGate, jarvisCyanGate, jarvisCyanBright)
}

// orbJarvisPixelAttempt runs one create-and-sample round of BOTH jarvis teeth
// against the LIVE desktop: it waits for the orb window to attach, centers the
// orb over a known (51,51,50) floor, and polls the composited screen pixels
// until the corners read the floor (transparency half) AND the cyan-dominant
// ratio beats jarvisCyanFloor (particle-color half), or the attempt budget
// expires. It returns (cleared, last corner readings, best cyan ratio, last
// sample count). Fatal paths inside are hard regressions (attach loss,
// window-API failure), not the cold-boot flake - only the expiry return is
// retryable. The floor is destroyed on return (its pump is this ball's STA
// thread), so the caller samples nothing after this function comes back.
func orbJarvisPixelAttempt(t *testing.T, b *Ball, attempt int) (bool, [4][3]int32, float64, int) {
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

	// Center the overlay in the work area: the ball window is the anchor the
	// orb follows (orb_windows.go follow()'s z-law), so move the BALL to the
	// center and let followOrb() align the orb over it - the same shape
	// orbCornerAttempt uses for the glass style's corners proof.
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

	// Known floor below the orb, created and painted on the ui-sta thread
	// (a window on a thread nobody pumps never paints). Z-order for the
	// sampling: ball > orb > floor - the floor goes in topmost, then slots
	// DIRECTLY BELOW the orb (hWndInsertAfter = orb), the same stack the
	// corners regression holds on the live desktop. It is destroyed here -
	// before the caller can Close the ball - because its pump is this ball's
	// STA thread.
	floor := make(chan uintptr, 1)
	b.sta.PostTask(func() { floor <- orbTestCreateColoredFloor(wr) })
	fh := <-floor
	if fh == 0 {
		t.Fatal("colored floor window failed to create")
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

	corners := [4][2]int32{
		{wr.l + jarvisProbeInset, wr.t + jarvisProbeInset},
		{wr.r - jarvisProbeInset - 1, wr.t + jarvisProbeInset},
		{wr.l + jarvisProbeInset, wr.b - jarvisProbeInset - 1},
		{wr.r - jarvisProbeInset - 1, wr.b - jarvisProbeInset - 1},
	}
	fr := int32(orbTestFloorColor & 0xff)
	fg := int32((orbTestFloorColor >> 8) & 0xff)
	fb := int32((orbTestFloorColor >> 16) & 0xff)

	// WebView2 boot is not instant; poll until both teeth clear or the
	// deadline proves this attempt never will.
	deadline = time.Now().Add(25 * time.Second)
	bestCyan := 0.0
	samples := 0
	var readings [4][3]int32
	for {
		hdc, _, _ := pGetDC.Call(0)
		if hdc == 0 {
			t.Fatal("GetDC(NULL) failed")
		}
		// Transparency half: all four corners read the floor color.
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
		// Particle-color half: cyan-dominant pixels across the orb rect.
		// This replaces the old non-dark>1% tooth: hue, not brightness.
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
		samples = total
		ratio := 0.0
		if total > 0 {
			ratio = float64(cyan) / float64(total)
		}
		if ratio > bestCyan {
			bestCyan = ratio
		}
		if cornersOK && ratio > jarvisCyanFloor {
			t.Logf("attempt %d: jarvis orb style proof: corners %v within +/-%d of (51,51,50); cyan-dominant %d/%d sampled px (%.2f%% > %.2f%% floor)",
				attempt, readings, jarvisCornerTol, cyan, total, ratio*100, jarvisCyanFloor*100)
			return true, readings, bestCyan, samples
		}
		if time.Now().After(deadline) {
			return false, readings, bestCyan, samples
		}
		time.Sleep(300 * time.Millisecond)
	}
}

// testJarvisRuntimeStyleSwitch exercises Ball.ApplyStyle on the LIVE window
// after the pixel proof: switch glass (the siri page) and back to jarvis,
// requiring the orb window to survive both swaps without a reattach - the
// in-place navigation seam, not a destroy/create cycle.
func testJarvisRuntimeStyleSwitch(t *testing.T, b *Ball) {
	t.Helper()
	hwnd := b.DebugOrbHWND()
	if hwnd == 0 {
		t.Fatal("orb window lost before the style-switch exercise")
	}
	if !b.ApplyStyle(styleGlassName) {
		t.Fatal("ApplyStyle(glass) refused")
	}
	if !b.ApplyStyle(styleJarvisName) {
		t.Fatal("ApplyStyle(jarvis) refused")
	}
	if ActiveStyle() != StyleJarvis {
		t.Fatalf("ActiveStyle after ApplyStyle(jarvis) = %v", ActiveStyle())
	}
	// The posted STA tasks need a beat to run; then the SAME window handle
	// must still be live - setStyle is a navigation, never a reattach.
	time.Sleep(500 * time.Millisecond)
	if b.DebugOrbHWND() != hwnd {
		t.Fatalf("orb window changed across ApplyStyle round-trip: %v -> %v", hwnd, b.DebugOrbHWND())
	}
	if got := OrbDebugLine(b.DebugOrbAttached(), b.DebugOrbState()); got != "orb-surface=webview2-orb orb-state=idle" {
		t.Errorf("debug face after style round-trip = %q", got)
	}
}
