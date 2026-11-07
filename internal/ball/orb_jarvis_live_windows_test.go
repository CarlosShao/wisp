//go:build windows && winlive

package ball

// Ticket 65 second visual style (the jarvis particle page): the LIVE proof
// that the orb host serving the jarvis style really renders the particle orb.
// The tooth is the same as TestBallLiveOrbPageRendersPixels: composited
// non-dark pixels over a known-black backdrop on the real desktop - a page
// that merely LOADED proves nothing, and a dark boot here is a real defect
// (delivery, CSP, WebGL or seam), not a cosmetic one. Run explicitly:
//
//	go test -tags winlive ./internal/ball/ -run TestBallLiveOrbJarvisStyleRendersParticles -count=1 -v

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TestBallLiveOrbJarvisStyleRendersParticles: select the jarvis style, create
// the real ball + orb window, and require the jarvis page to light up the
// desktop. The style family is a process-wide global, so the selection is
// restored to glass before this test returns - a later test in the same
// process runs its glass contract unchanged.
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
	var lastBest float64
	var lastSamples int
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
		lit, best, samples := orbPixelAttempt(t, b, attempt)
		if lit {
			testJarvisRuntimeStyleSwitch(t, b)
			return
		}
		lastBest, lastSamples = best, samples
		b.Close()
		orbTestReapLeakedWebView2(t)
		if attempt < attempts {
			t.Logf("attempt %d/%d stayed dark (best %.3f%% over %d sampled px); rebuilding ball and WebView2 instance",
				attempt, attempts, best*100, samples)
		}
	}
	t.Fatalf("jarvis orb page stayed black after %d attempts: best non-dark ratio %.3f%% over %d sampled px (particle page did not render; delivery, CSP, WebGL or seam broken)",
		attempts, lastBest*100, lastSamples)
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
