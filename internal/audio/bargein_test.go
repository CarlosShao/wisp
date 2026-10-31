package audio

import (
	"math"
	"testing"
)

// The barge-in combination gate's unit surface: every constant and decision
// rule here is pinned to the 59-w-r1 / 59-v3 readings (docs/evidence/s1/
// 59wr1-longrun-gates.md sections 4-5, docs/evidence/s1/59-v3-verdict.md
// section 3). The synthesized streams below are deterministic sines - no
// fixture files, no wall clock, no goroutines.

// bargeinSine is a deterministic mono sine at the seam rate.
func bargeinSine(freq float64, n int) []float32 {
	out := make([]float32, n)
	for i := range out {
		out[i] = float32(0.5 * math.Sin(2*math.Pi*freq*float64(i)/float64(TargetRate)))
	}
	return out
}

// bargeinEcho derives the capture of a pure echo: far delayed by lag.
func bargeinEcho(far []float32, lag int) []float32 {
	out := make([]float32, len(far))
	for i := range out {
		j := i - lag
		if j >= 0 {
			out[i] = far[j]
		}
	}
	return out
}

// bargeinGain scales a span by a linear gain.
func bargeinGain(s []float32, g float64) []float32 {
	out := make([]float32, len(s))
	for i, v := range s {
		out[i] = float32(float64(v) * g)
	}
	return out
}

func TestBargeinWindowStatERLEdB(t *testing.T) {
	if !math.IsNaN((BargeinWindowStat{}).ERLEdB()) {
		t.Fatalf("no far-active frames: want NaN, got %v", (BargeinWindowStat{}).ERLEdB())
	}
	if !math.IsInf((BargeinWindowStat{ActiveFrames: 3, NearPow: 1}).ERLEdB(), 1) {
		t.Fatalf("perfect cancellation: want +Inf")
	}
	if !math.IsInf((BargeinWindowStat{ActiveFrames: 3, ResPow: 1}).ERLEdB(), -1) {
		t.Fatalf("zero near: want -Inf")
	}
	got := (BargeinWindowStat{ActiveFrames: 1, NearPow: 100, ResPow: 1}).ERLEdB()
	if want := 20.0; math.Abs(got-want) > 1e-9 {
		t.Fatalf("10*log10(100) = %v, want %v", got, want)
	}
}

// TestSlidingERLEFoldsWindows pins the window folding and the chunk
// invariance: the same stream fed in any chunking must produce the identical
// window list.
func TestSlidingERLEFoldsWindows(t *testing.T) {
	far := bargeinSine(220, BargeinWindowSamples*3)
	farRMS := frameRMSBargein(far)
	// near = the echo at 0dB ERLE; residual = the echo reduced by 28dB.
	near := bargeinEcho(far, 493)
	res := bargeinGain(far, math.Pow(10, -28.0/20.0))

	build := func(chunks ...int) *SlidingERLE {
		acc := NewSlidingERLE(farRMS)
		pos := 0
		for _, c := range chunks {
			hi := pos + c
			if hi > len(far) {
				hi = len(far)
			}
			acc.Observe(far[pos:hi], near[pos:hi], res[pos:hi])
			pos = hi
		}
		return acc
	}
	whole := build(len(far))
	halves := build(777, BargeinWindowSamples+1, 999, len(far))
	if len(whole.Windows) != len(halves.Windows) {
		t.Fatalf("chunk invariance: %d windows vs %d", len(whole.Windows), len(halves.Windows))
	}
	for i := range whole.Windows {
		w, h := whole.Windows[i], halves.Windows[i]
		if w.ActiveFrames != h.ActiveFrames || math.Abs(w.NearPow-h.NearPow) > 1e-6 || math.Abs(w.ResPow-h.ResPow) > 1e-6 {
			t.Fatalf("window %d differs under chunking: %+v vs %+v", i, w, h)
		}
	}
	// Three fully far-active 1s windows at 28dB, plus the partial tail that
	// never closes. A window closes when its 32nd frame lands (32x512=16384
	// crosses the 16000 boundary), so every closed window aggregates 32
	// full frames - the same boundary behavior as the rig's accumulator.
	if got := len(whole.Windows); got != 2 {
		t.Fatalf("3s stream (93 frames) closed %d windows, want 2", got)
	}
	for i, w := range whole.Windows {
		if w.ActiveFrames != 32 {
			t.Fatalf("window %d active frames = %d, want 32", i, w.ActiveFrames)
		}
		if db := w.ERLEdB(); math.Abs(db-28) > 0.5 {
			t.Fatalf("window %d ERLE = %.2f, want ~28", i, db)
		}
	}
}

// TestCollapseDetectorShapes pins both alarm shapes, the once-per-streak
// rule, the inactive-window skip, and the recovery.
func TestCollapseDetectorShapes(t *testing.T) {
	d := NewCollapseDetector()
	mk := func(idx int, db float64) BargeinWindowStat {
		// ERLE dB = 10*log10(near/res): 20dB needs near/res = 100.
		return BargeinWindowStat{Index: idx, ActiveFrames: 10, NearPow: math.Pow(10, db/10), ResPow: 1}
	}
	var alarms []CollapseAlarm
	for i := 0; i < 5; i++ {
		alarms = append(alarms, d.Observe(mk(i, 28))...) // steady platform: silent
	}
	if len(alarms) != 0 {
		t.Fatalf("steady platform fired: %+v", alarms)
	}
	// Collapse to 8dB: the drop shape fires first (the baseline ring still
	// holds five 28dB windows, median 28, and 8 <= 28-6), once, when the
	// streak reaches 3 sustained windows.
	for i := 5; i < 8; i++ {
		alarms = append(alarms, d.Observe(mk(i, 8))...)
	}
	if len(alarms) != 1 || alarms[0].Reason != "drop" || alarms[0].WindowIndex != 7 ||
		math.Abs(alarms[0].Baseline-28) > 1e-9 || math.Abs(alarms[0].WindowERLE-8) > 1e-9 {
		t.Fatalf("collapse streak: %+v, want one drop alarm at window 7 (baseline 28)", alarms)
	}
	// The streak continuing does not re-fire.
	if got := d.Observe(mk(8, 8)); len(got) != 0 {
		t.Fatalf("continuing streak re-fired: %+v", got)
	}
	// Recovery resets the streak.
	if got := d.Observe(mk(9, 28)); len(got) != 0 {
		t.Fatalf("recovery fired: %+v", got)
	}
	if d.Sustained() {
		t.Fatalf("sustained after recovery")
	}
	// A second collapse re-arms and fires again.
	for i := 10; i < 13; i++ {
		alarms = append(alarms, d.Observe(mk(i, 8))...)
	}
	if len(alarms) != 2 {
		t.Fatalf("second collapse: %+v, want one more alarm", alarms)
	}
}

// TestCollapseDetectorSustained pins the tighten window: untrusted from the
// streak completing its third window until a window fails both shapes.
func TestCollapseDetectorSustained(t *testing.T) {
	d := NewCollapseDetector()
	mk := func(idx int, db float64) BargeinWindowStat {
		return BargeinWindowStat{Index: idx, ActiveFrames: 10, NearPow: math.Pow(10, db/10), ResPow: 1}
	}
	if d.Sustained() {
		t.Fatalf("sustained on an empty detector")
	}
	d.Observe(mk(0, 28))
	d.Observe(mk(1, 8))
	d.Observe(mk(2, 8))
	if d.Sustained() {
		t.Fatalf("sustained before the streak completes")
	}
	d.Observe(mk(3, 8))
	if !d.Sustained() {
		t.Fatalf("not sustained after three qualifying windows")
	}
	// Inactive windows (no far-active frames) ride through without breaking
	// the streak.
	if got := d.Observe(BargeinWindowStat{Index: 3}); got != nil {
		t.Fatalf("inactive window returned %+v", got)
	}
	if !d.Sustained() {
		t.Fatalf("inactive window broke the streak")
	}
	d.Observe(mk(4, 28))
	if d.Sustained() {
		t.Fatalf("recovery window did not clear the streak")
	}
}

// TestFarCorrelationDetector pins the sliding correlation: a pure echo keeps
// correlation at 1 (including under the lag shift), foreign energy pulls it
// below the floor, all-silent spans are NaN and never mark untrusted.
func TestFarCorrelationDetector(t *testing.T) {
	far := bargeinSine(200, BargeinWindowSamples*3)
	const lag = 493
	echo := bargeinEcho(far, lag)

	d := NewFarCorrelationDetector(lag)
	d.Observe(far, echo)
	if len(d.Values) < 2 {
		t.Fatalf("correlation windows closed: %d", len(d.Values))
	}
	for i, v := range d.Values {
		if math.IsNaN(v) || math.Abs(v-1) > 1e-6 {
			t.Fatalf("pure echo window %d corr = %v (all=%v), want ~1", i, v, d.Values)
		}
	}
	if d.Untrusted() {
		t.Fatalf("pure echo marked untrusted")
	}

	// Foreign energy: the whole second window's capture replaced by an
	// uncorrelated signal (a different frequency). Windows close on the
	// 32-frame boundary, so window 0 [0,16384) catches the first 384
	// polluted samples (2.3 percent - a light touch that stays above the
	// floor), while window 1 is fully polluted.
	polluted := append([]float32(nil), echo...)
	for i := BargeinWindowSamples; i < 2*BargeinWindowSamples; i++ {
		polluted[i] = float32(0.5 * math.Sin(2*math.Pi*1800*float64(i)/float64(TargetRate)))
	}
	d2 := NewFarCorrelationDetector(lag)
	d2.Observe(far, polluted)
	if v := d2.Values[0]; math.IsNaN(v) || v <= BargeinFarCorrFloor {
		t.Fatalf("lightly-touched window 0 corr = %v, want above %.1f (still trusted)", v, BargeinFarCorrFloor)
	}
	if v := d2.Values[1]; !d2.Untrusted() || v >= BargeinFarCorrFloor {
		t.Fatalf("polluted window corr = %v, want below %.1f and untrusted", v, BargeinFarCorrFloor)
	}

	// All-silent span: NaN, never untrusted.
	d3 := NewFarCorrelationDetector(0)
	d3.Observe(make([]float32, BargeinWindowSamples*2), make([]float32, BargeinWindowSamples*2))
	if !math.IsNaN(d3.Values[0]) || d3.Untrusted() {
		t.Fatalf("silent span: values %v untrusted %v", d3.Values, d3.Untrusted())
	}
}

// TestBargeinGateStaticTiers pins the static chain: energy cut, ASR text
// gate, and the tier combination on episodes.
func TestBargeinGateStaticTiers(t *testing.T) {
	g := NewBargeinGate(BargeinOptions{DynamicGating: false})
	if g.opts.EnergyThresholdDbfs != BargeinNominalEnergyDbfs {
		t.Fatalf("nominal cut = %v, want %v", g.opts.EnergyThresholdDbfs, BargeinNominalEnergyDbfs)
	}
	cases := []struct {
		rms  float64
		text string
		pass bool
	}{
		{-26.70, "天气怎么样", true}, // the control-26 reading: passes
		{-36.58, "天气怎么样", true}, // the control-36 reading: reachable on the loose nominal cut
		{-32.20, "听我", true},    // the collapse artifact's static shape: passes at -45 (what the dynamic gating exists to eat)
		{-32.20, "", false},     // artifact without text: the text gate holds it
		{-50.00, "真话", false},   // too quiet even for the loose cut
	}
	// Recompute the expectation table honestly against the nominal cut:
	for i, c := range cases {
		got := g.Evaluate(BargeinEpisode{StartSample: 0, RmsDbfs: c.rms, AsrText: c.text})
		want := c.rms > BargeinNominalEnergyDbfs && c.text != ""
		if got.Pass != want || got.Pass != c.pass {
			t.Fatalf("case %d (rms %.2f text %q): pass=%v want %v", i, c.rms, c.text, got.Pass, want)
		}
		if got.ThresholdDbfs != BargeinNominalEnergyDbfs {
			t.Fatalf("case %d: threshold %v, want the nominal", i, got.ThresholdDbfs)
		}
		if got.TightenActive {
			t.Fatalf("case %d: tighten active with dynamic gating disabled", i)
		}
	}
}

// buildCollapseStream renders a 6s AEC path: two good-platform windows
// (28dB) then three collapse windows (8dB) then a good window - the rig's
// "converged then lost" shape. Returns the three streams.
func buildCollapseStream() (far, near, res []float32) {
	far = bargeinSine(220, BargeinWindowSamples*6)
	near = bargeinEcho(far, 493)
	res = make([]float32, len(near))
	copy(res, near)
	// Good windows: residual = echo reduced 28dB. Collapse windows [2,5):
	// residual = the echo at 8dB (below the 15dB absolute floor).
	good := math.Pow(10, -28.0/20.0)
	bad := math.Pow(10, -8.0/20.0)
	for i := range res {
		w := i / BargeinWindowSamples
		g := good
		if w >= 2 && w < 5 {
			g = bad
		}
		res[i] = float32(float64(near[i]) * g)
	}
	return far, near, res
}

// TestBargeinGateDynamicTighten is the product-shape core reading: with the
// dynamic gating on, an episode in an untrusted window is judged at the
// tightened cut (the -32.2dBFS artifact dies, the -26.7dBFS control passes),
// while the same artifact outside any untrusted window keeps the nominal
// cut.
func TestBargeinGateDynamicTighten(t *testing.T) {
	far, near, res := buildCollapseStream()
	g := NewBargeinGate(BargeinOptions{
		DynamicGating: true,
		LagSamples:    493,
		FarRMS:        frameRMSBargein(far),
	})
	for lo := 0; lo < len(far); lo += 4000 {
		hi := lo + 4000
		if hi > len(far) {
			hi = len(far)
		}
		g.Feed(far[lo:hi], near[lo:hi], res[lo:hi])
	}
	wins := g.UntrustedWindows()
	if len(wins) != 5 { // five full windows closed
		t.Fatalf("windows closed: %d, want 5", len(wins))
	}
	for i, u := range wins {
		if i >= 4 && !u {
			t.Fatalf("window %d (alarm sustained from the third collapse window) not untrusted: %v", i, wins)
		}
		if i < 4 && u {
			t.Fatalf("window %d untrusted before the streak completed: %v", i, wins)
		}
	}
	// The collapse-state artifact (the rig reading -32.2dBFS, "listen to me")
	// starting in the untrusted window: eaten by the tightened cut.
	art := g.Evaluate(BargeinEpisode{StartSample: 4 * BargeinWindowSamples, RmsDbfs: -32.20, AsrText: "listen to me"})
	if art.Pass || !art.TightenActive || !art.SuppressedByDynamic {
		t.Fatalf("collapse artifact: %+v, want suppressed by the dynamic gate", art)
	}
	if art.ThresholdDbfs != BargeinTightenedEnergyDbfs {
		t.Fatalf("applied threshold %v, want %v", art.ThresholdDbfs, BargeinTightenedEnergyDbfs)
	}
	// The positive control (the rig reading -26.70dBFS) in the same
	// untrusted window: real speech survives the tightened cut.
	ctl := g.Evaluate(BargeinEpisode{StartSample: 4 * BargeinWindowSamples, RmsDbfs: -26.70, AsrText: "How is the weather today"})
	if !ctl.Pass || !ctl.TightenActive || ctl.SuppressedByDynamic {
		t.Fatalf("positive control in an untrusted window: %+v, want pass", ctl)
	}
	// The same artifact level in a trusted window: judged at the nominal
	// cut, where it passes on energy and text (the static -45 shape).
	tr := g.Evaluate(BargeinEpisode{StartSample: 0, RmsDbfs: -32.20, AsrText: "listen to me"})
	if !tr.Pass || tr.TightenActive {
		t.Fatalf("trusted-window artifact: %+v, want the nominal shape", tr)
	}
	// Quiet speech in a trusted window: the loose nominal keeps -36.58
	// reachable (the control-36 shape the static -30 gate ate).
	q := g.Evaluate(BargeinEpisode{StartSample: 0, RmsDbfs: -36.58, AsrText: "How is the weather today"})
	if !q.Pass {
		t.Fatalf("quiet speech at the nominal cut: %+v, want pass", q)
	}
	// ...and the tightened cut eats it while untrusted (the accepted cost
	// the 59-v3 verdict recorded).
	q2 := g.Evaluate(BargeinEpisode{StartSample: 4 * BargeinWindowSamples, RmsDbfs: -36.58, AsrText: "How is the weather today"})
	if q2.Pass || !q2.SuppressedByDynamic {
		t.Fatalf("quiet speech in an untrusted window: %+v", q2)
	}
}

// TestBargeinGateDynamicDisabled is the reverse criterion's product shape:
// with the dynamic gating off, the gate is the static tier chain and the
// tightened cut never applies. The judgment channels keep observing (the
// untrusted flags still record), but no decision reads them.
func TestBargeinGateDynamicDisabled(t *testing.T) {
	far, near, res := buildCollapseStream()
	g := NewBargeinGate(BargeinOptions{
		DynamicGating: false,
		LagSamples:    493,
		FarRMS:        frameRMSBargein(far),
	})
	for lo := 0; lo < len(far); lo += 16000 {
		hi := lo + 16000
		if hi > len(far) {
			hi = len(far)
		}
		g.Feed(far[lo:hi], near[lo:hi], res[lo:hi])
	}
	art := g.Evaluate(BargeinEpisode{StartSample: 4 * BargeinWindowSamples, RmsDbfs: -32.20, AsrText: "listen to me"})
	if !art.Pass || art.TightenActive || art.SuppressedByDynamic {
		t.Fatalf("artifact with the gate off: %+v, want the static pass (the reverse reading)", art)
	}
	if art.ThresholdDbfs != BargeinNominalEnergyDbfs {
		t.Fatalf("threshold with the gate off: %v, want the nominal", art.ThresholdDbfs)
	}
}

// TestBargeinGateLagShift pins the ERLE channel's internal ring: a capture
// that IS the far delayed by 493 samples, with a perfectly cancelling
// residual, reads high ERLE only if the shift is right.
func TestBargeinGateLagShift(t *testing.T) {
	far := bargeinSine(200, BargeinWindowSamples*3)
	echo := bargeinEcho(far, 493)
	res := make([]float32, len(echo)) // perfect cancellation
	g := NewBargeinGate(BargeinOptions{
		DynamicGating: true,
		LagSamples:    493,
		FarRMS:        frameRMSBargein(far),
	})
	g.Feed(far, echo, res)
	wins := g.UntrustedWindows()
	// 48000 samples = 93 full frames: two windows close (a third would need
	// 96), the tail stays in the open window.
	if len(wins) != 2 {
		t.Fatalf("windows closed: %d, want 2", len(wins))
	}
	for i, u := range wins {
		if u {
			t.Fatalf("window %d untrusted under perfect cancellation: %v", i, wins)
		}
	}
	if db := g.erle.Windows[0].ERLEdB(); !math.IsInf(db, 1) {
		t.Fatalf("ERLE with the correct lag shift = %v, want +Inf (ResPow 0)", db)
	}
	// The correlation channel sees the pure echo too.
	if v := g.fcorr.Values[0]; math.IsNaN(v) || math.Abs(v-1) > 1e-6 {
		t.Fatalf("correlation under the gate's own shift = %v, want ~1", v)
	}
}

// TestBargeinGateWindowSeam pins the window-level input path: the same
// collapse shape fed as aggregated windows (engine supervision numbers, or
// recorded windows on the offline readings path) drives the same tighten
// decisions as the raw-frame Feed.
func TestBargeinGateWindowSeam(t *testing.T) {
	g := NewBargeinGate(BargeinOptions{DynamicGating: true})
	mk := func(idx int, db float64) BargeinWindowStat {
		return BargeinWindowStat{Index: idx, ActiveFrames: 10, NearPow: math.Pow(10, db/10), ResPow: 1}
	}
	for i := 0; i < 4; i++ {
		db := 28.0
		if i >= 1 {
			db = 8
		}
		g.ObserveWindow(mk(i, db), 0.95)
	}
	wins := g.UntrustedWindows()
	if !wins[3] || wins[0] {
		t.Fatalf("window-seam untrusted flags: %v", wins)
	}
	art := g.Evaluate(BargeinEpisode{StartSample: 3 * BargeinWindowSamples, RmsDbfs: -32.20, AsrText: "listen to me"})
	if art.Pass || !art.TightenActive || !art.SuppressedByDynamic {
		t.Fatalf("window-seam artifact: %+v, want suppressed by the dynamic gate", art)
	}
	// A NaN correlation (no comparable signal) never marks untrusted by
	// itself: five clean windows with an unmeasurable correlation stay
	// trusted.
	g2 := NewBargeinGate(BargeinOptions{DynamicGating: true})
	for i := 0; i < 5; i++ {
		g2.ObserveWindow(mk(i, 28), math.NaN())
	}
	if u := g2.UntrustedWindows(); u[4] {
		t.Fatalf("NaN correlation marked untrusted: %v", u)
	}
}
