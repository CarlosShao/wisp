//go:build windows && winlive

package ball

import (
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Ticket 62 wave-3: the glass teeth the 62-v1 adversarial acceptance called
// missing. The page-side seam (orbStateScript / orbAudioScript) is pinned as
// STRING SHAPE in orb_test.go; this file pins the thing the owner actually
// sees - the mapped state and the mic envelope must reach VISIBLE PIXELS on
// the real WebView2 orb page, measured on the composited screen over a known
// black backdrop, never "the page loaded".
//
// Thresholds are pinned from the wave-3 probe run (2026-11-07, this box, 961
// sampled px per reading, 3px grid, 4px inset):
//
//	T1 idle(Sleeping)=0.1873 vs thinking(Thinking)=0.4204  -> absdiff 0.2331
//	T2 audio driven at level 0.9 vs silence (probe: +0.0156; green run's
//	   post-feed decay tail read +0.003, so the tooth samples DURING the feed)
//	T3 ordering silence < level 0.5 < level 0.9, each round
//	   (probe: 0.4214 < 0.4329 < 0.4339; the 0.5 -> 0.9 step measured +0.0010,
//	    inside single-read noise, so only the ordering is pinned; the top step
//	    is reported here and in docs/evidence/s1/62-liquid-glass-closure.md,
//	    not asserted)
//
// Floors of 0.02 on the LOW side of every comparison keep a black page from
// passing by coincidence (0 == 0 is not a state difference); the dark-boot
// flake (A1602 family) is handled by the same bounded rebuild-retry shape
// orbPixelAttempt uses - a retry happens only when the page never painted,
// never when a threshold is missed.
func TestBallLiveOrbStateAndAudioReachPixels(t *testing.T) {
	store, err := OpenPositionStore(filepath.Join(t.TempDir(), "ballpos.json"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	requireQuietBallDesktop(t)
	orbTestReapLeakedWebView2(t)
	restored := PrototypeVisualsEnabled()
	t.Cleanup(func() { EnablePrototypeVisuals(restored) })
	EnablePrototypeVisuals(true)

	const attempts = 3
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
		lit, readings := orbGlass62Attempt(t, b, attempt)
		b.Close()
		orbTestReapLeakedWebView2(t)
		if lit {
			idle, think, silence, r05, r09 := readings[0], readings[1], readings[2], readings[3], readings[4]
			t.Logf("attempt %d: idle=%.4f thinking=%.4f silence=%.4f level0.5=%.4f level0.9=%.4f",
				attempt, idle, think, silence, r05, r09)
			// T1: the mapped state changes the picture (frozen 0.10 floor,
			// measured 0.2331).
			if diff := absf(think - idle); diff < 0.10 {
				t.Fatalf("attempt %d: idle %.4f vs thinking %.4f differ by %.4f < 0.10: the pushed state does not reach pixels",
					attempt, idle, think, diff)
			}
			if minf(idle, think) < 0.02 {
				t.Fatalf("attempt %d: idle %.4f / thinking %.4f below the 0.02 alive floor: the page never painted", attempt, idle, think)
			}
			// T2: the mic envelope changes the picture (frozen 0.005 floor;
			// measured +0.0156 during the feed).
			if silence < 0.02 {
				t.Fatalf("attempt %d: silence %.4f below the 0.02 alive floor", attempt, silence)
			}
			if lift := r09 - silence; lift < 0.005 {
				t.Fatalf("attempt %d: silence %.4f vs level-0.9 %.4f lifts only %.5f < 0.005: the pushed audio bands do not reach pixels",
					attempt, silence, r09, lift)
			}
			// T3: ordering across levels (both rounds; the 0.5->0.9 step is
			// NOT pinned - measured +0.0010, single-read noise).
			if r05 <= silence {
				t.Fatalf("attempt %d: level-0.5 %.4f <= silence %.4f: audio does not lift the page monotonically (round 1)",
					attempt, r05, silence)
			}
			if r09 <= r05 {
				t.Fatalf("attempt %d: level-0.9 %.4f <= level-0.5 %.4f: audio does not lift the page monotonically (round 2)",
					attempt, r09, r05)
			}
			return
		}
		if attempt < attempts {
			t.Logf("attempt %d/%d: the orb page never painted; rebuilding the ball and WebView2 instance", attempt, attempts)
		}
	}
	t.Fatalf("orb page never painted after %d attempts: the state/audio pixel teeth could not run", attempts)
}

// orbGlass62Attempt runs one create-measure round: it parks the orb over a
// known-black backdrop, waits for the page to paint, then takes the five
// readings (idle, thinking, silence, level 0.5, level 0.9). It returns
// (painted, [idle, think, silence, r05, r09]); readings are valid only when
// painted. Fatal paths inside are hard regressions (attach loss, window API
// failure), not the dark-boot flake - only the never-painted return retries.
func orbGlass62Attempt(t *testing.T, b *Ball, attempt int) (bool, [5]float64) {
	t.Cleanup(b.Close)

	deadline := time.Now().Add(15 * time.Second)
	var orbHWND windows.HWND
	for {
		orbHWND = b.DebugOrbHWND()
		if orbHWND != 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt %d: orb window did not attach within 15s", attempt)
		}
		time.Sleep(100 * time.Millisecond)
	}

	mon, ok := monitorRectsForWindow(windows.HWND(b.DebugHWND()))
	if !ok {
		t.Fatal("monitorRectsForWindow(ball) returned no monitor")
	}
	work := mon.Work

	ballHWND := uintptr(b.DebugHWND())
	pSetWindowPos.Call(ballHWND, 0,
		uintptr(int32(work.L)+8), uintptr(int32(work.T)+8), 0, 0,
		uintptr(swpNoSize|swpNoZOrder|swpNoActivate))
	var wr rect
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) failed")
	}
	edge := wr.width()
	ox := (int32(work.L) + int32(work.R)) / 2
	oy := (int32(work.T) + int32(work.B)) / 2
	pSetWindowPos.Call(uintptr(orbHWND), orbHwndTopmost,
		uintptr(ox-edge/2), uintptr(oy-edge/2), 0, 0,
		uintptr(swpNoSize|swpNoActivate))
	if r, _, _ := pGetWindowRect.Call(uintptr(orbHWND), unsafePtr(&wr)); r == 0 {
		t.Fatal("GetWindowRect(orb) after the move failed")
	}

	// Known-black floor behind the orb (orbTestCreateBackdrop), created and
	// destroyed on the ball's STA thread, orb re-asserted above it.
	backdrop := make(chan uintptr, 1)
	b.sta.PostTask(func() { backdrop <- orbTestCreateBackdrop(wr) })
	bh := <-backdrop
	if bh == 0 {
		t.Fatal("black backdrop window failed to create")
	}
	pSetWindowPos.Call(uintptr(orbHWND), orbHwndTopmost, 0, 0, 0, 0,
		uintptr(swpNoSize|swpNoMove|swpNoActivate))
	defer func() {
		done := make(chan struct{})
		b.sta.PostTask(func() {
			pDestroyWindow.Call(bh)
			close(done)
		})
		<-done
	}()

	// Poll until the page paints (WebGPU boot is not instant).
	lit := false
	for i := 0; i < 60; i++ {
		time.Sleep(500 * time.Millisecond)
		r, _ := orbTestNonDarkRatio(wr.l+4, wr.t+4, wr.width()-8, wr.height()-8)
		if r > 0.01 {
			lit = true
			break
		}
	}
	if !lit {
		return false, [5]float64{}
	}
	time.Sleep(2 * time.Second)

	sample := func() float64 {
		r, _ := orbTestNonDarkRatio(wr.l+4, wr.t+4, wr.width()-8, wr.height()-8)
		return r
	}
	// meanOf3 denoises one condition: three screen samples 400ms apart.
	meanOf3 := func() float64 {
		sum := 0.0
		for i := 0; i < 3; i++ {
			time.Sleep(400 * time.Millisecond)
			sum += sample()
		}
		return sum / 3
	}
	// sampleWhile feeds the level continuously at the capture rate (30Hz)
	// WHILE sampling, so the reading is the DRIVEN page, not the "zero them
	// on stop" decay tail the page ramps down after feeding ends (the wave-3
	// green run measured exactly that: a post-feed sample read a 0.003 lift
	// where a during-feed probe read 0.0115). SetAudioLevel is documented
	// safe from any goroutine; the feeder stops before the return.
	sampleWhile := func(level float32) float64 {
		stop := make(chan struct{})
		go func() {
			tk := time.NewTicker(33 * time.Millisecond)
			defer tk.Stop()
			for {
				select {
				case <-stop:
					return
				case <-tk.C:
					b.SetAudioLevel(level)
				}
			}
		}()
		sum := 0.0
		for i := 0; i < 3; i++ {
			time.Sleep(400 * time.Millisecond)
			sum += sample()
		}
		close(stop)
		time.Sleep(1500 * time.Millisecond) // let the bands settle back down
		return sum / 3
	}

	// The readings, in a fixed order with a settle after every push. State
	// differences are 0.23 of the ratio - single samples carry them; audio
	// differences are ~0.01 - three-sample means and during-feed sampling
	// carry them.
	b.SetState(statemachine.StateSleeping)
	time.Sleep(2 * time.Second)
	idle := sample()

	b.SetState(statemachine.StateThinking)
	time.Sleep(2 * time.Second)
	think := sample()

	b.SetState(statemachine.StateListening)
	time.Sleep(2 * time.Second)
	silence := meanOf3()

	r05 := sampleWhile(0.5)
	r09 := sampleWhile(0.9)

	return true, [5]float64{idle, think, silence, r05, r09}
}

func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
