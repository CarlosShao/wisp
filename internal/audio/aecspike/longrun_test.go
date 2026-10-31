package aecspike

import (
	"math"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
)

// CI-safe pins for the long-window instruments. The windows rig and the
// gated drivers are separate files; everything here runs on every GOOS.

// fracEcho builds near[i] ~ far[i-(lag+frac)] by linear interpolation, the
// shape a fractional render/loopback phase offset produces.
func fracEcho(far []float32, lag int, frac float64) []float32 {
	near := make([]float32, len(far))
	for i := range near {
		j := float64(i-lag) - frac
		lo := int(math.Floor(j))
		d := j - float64(lo)
		var v float64
		if lo >= 0 && lo < len(far) {
			v += (1 - d) * float64(far[lo])
		}
		if lo+1 >= 0 && lo+1 < len(far) {
			v += d * float64(far[lo+1])
		}
		near[i] = float32(v)
	}
	return near
}

func TestFractionalLagRecoversSubSampleDelay(t *testing.T) {
	far, _ := echoFixture(0, 16000*2, 1, 0)
	for _, tc := range []struct {
		lag  int
		frac float64
	}{{493, -0.3}, {493, 0}, {493, 0.25}, {493, 0.45}, {100, 0.4}} {
		near := fracEcho(far, tc.lag, tc.frac)
		got := FractionalLag(far, near, tc.lag)
		if math.Abs(got-tc.frac) > 0.05 {
			t.Fatalf("lag %d frac %.2f: FractionalLag = %.3f, want within 0.05", tc.lag, tc.frac, got)
		}
	}
	// Degenerate inputs answer 0 rather than dividing by zero.
	if got := FractionalLag(far, far[:10], 5000); got != 0 {
		t.Fatalf("FractionalLag on a truncated near = %.3f, want 0", got)
	}
}

func TestERLEAccumMatchesBaselineERLEdB(t *testing.T) {
	far, near := echoFixture(1600, 16000*4, 0.5, 0.001)
	lag, _ := AlignByEnvelopeCorrelation(far, near, DefaultMaxLag)
	residual, _, _ := SubtractEcho(far, near, lag)
	// The baseline's own convention, over the whole recording.
	want := ERLEdB(far, near, residual, lag)
	// The accumulator with the lag-shifted far stream, fed in ragged chunks
	// (the shape the streaming AEC pass produces).
	rms := frameRMS(far)
	acc := NewERLEAccum(rms)
	// activity stream: far delayed by lag (frames before the lag are silent)
	activity := make([]float32, lag)
	activity = append(activity, far...)
	lo := 0
	for lo < len(near) {
		hi := lo + 1000 + lo%7
		if hi > len(near) {
			hi = len(near)
		}
		aEnd := lo + (hi - lo)
		if aEnd > len(activity) {
			aEnd = len(activity)
		}
		acc.Observe(activity[lo:aEnd], near[lo:hi], residual[lo:hi])
		lo = hi
	}
	got := acc.SpanERLE(0, len(acc.Windows))
	if math.Abs(got-want) > 0.5 {
		t.Fatalf("SpanERLE = %.3f dB, baseline ERLEdB = %.3f dB", got, want)
	}
	if len(acc.Windows) < 3 {
		t.Fatalf("window count = %d, want several 1s windows", len(acc.Windows))
	}
	for _, w := range acc.Windows {
		if w.ActiveFrames == 0 {
			t.Fatalf("window %d: no far-active frames on an active fixture", w.Index)
		}
	}
}

func TestERLEAccumSkipsSilentSpans(t *testing.T) {
	// A far with a genuine silent stretch: windows over the silence carry no
	// active frames and report NaN, not a fake ERLE. Windows close on 512
	// -sample frames, so the fixture is two exact windows of 16384 samples.
	const win = 32 * erleFrameSamples
	far := make([]float32, win*2)
	for i := win; i < win*2; i++ { // active only in the second window
		far[i] = float32(math.Sin(2 * math.Pi * 220 * float64(i) / 16000))
	}
	near := append([]float32(nil), far...)
	residual := make([]float32, len(far))
	acc := NewERLEAccum(frameRMS(far))
	acc.Observe(far, near, residual)
	if len(acc.Windows) != 2 {
		t.Fatalf("windows = %d, want 2", len(acc.Windows))
	}
	if !math.IsNaN(acc.Windows[0].ERLEdB()) {
		t.Fatalf("silent window ERLE = %.2f, want NaN", acc.Windows[0].ERLEdB())
	}
	if math.IsNaN(acc.Windows[1].ERLEdB()) {
		t.Fatalf("active window ERLE is NaN, want a number (perfect cancellation)")
	}
	if got := acc.SpanERLE(0, 2); !math.IsInf(got, 1) {
		// Span over both: the first window has zero active frames, the second
		// cancels perfectly, so the span is +Inf, not NaN.
		t.Fatalf("span ERLE = %v, want +Inf (only the active window counts)", got)
	}
}

func wstat(index int, db float64) ERLEWindowStat {
	return ERLEWindowStat{
		Index:        index,
		NearPow:      math.Pow(10, db/10),
		ResPow:       1,
		ActiveFrames: 312, // a fully active 1s window (32ms frames)
	}
}

func TestCollapseDetectorAlarmsOnMidRunDrop(t *testing.T) {
	d := NewCollapseDetector()
	var alarms []CollapseAlarm
	for i := 0; i < 10; i++ {
		alarms = append(alarms, d.Observe(wstat(i, 30))...)
	}
	if len(alarms) != 0 {
		t.Fatalf("steady 30dB: %d alarms, want 0", len(alarms))
	}
	fired := 0
	for i := 10; i < 20; i++ {
		got := d.Observe(wstat(i, 9))
		alarms = append(alarms, got...)
		fired += len(got)
	}
	if fired != 1 {
		t.Fatalf("drop to 9dB sustained: %d alarms, want exactly 1", fired)
	}
	a := alarms[len(alarms)-1]
	if a.Reason != "drop" || a.WindowIndex != 12 {
		t.Fatalf("alarm = %+v, want drop at window 12 (3rd sustained)", a)
	}
	// Back to good, then drop again: the streak re-arms and fires again.
	for i := 20; i < 23; i++ {
		if got := d.Observe(wstat(i, 30)); len(got) != 0 {
			t.Fatalf("recovery window %d alarmed: %+v", i, got)
		}
	}
	if got := d.Observe(wstat(23, 9)); len(got) != 0 {
		t.Fatalf("1st re-drop window alarmed early: %+v", got)
	}
	_ = d.Observe(wstat(24, 9))
	if got := d.Observe(wstat(25, 9)); len(got) != 1 {
		t.Fatalf("re-drop streak did not re-fire: %+v", got)
	}
}

func TestCollapseDetectorAlarmsOnColdLow(t *testing.T) {
	// The cold collapse shape: never converged, no good baseline to drop
	// from - the absolute floor carries the alarm.
	d := NewCollapseDetector()
	var fired int
	var last *CollapseAlarm
	for i := 0; i < 10; i++ {
		got := d.Observe(wstat(i, 9))
		fired += len(got)
		if len(got) > 0 {
			a := got[0]
			last = &a
		}
	}
	if fired != 1 || last == nil || last.Reason != "low-absolute" || last.WindowIndex != 2 {
		t.Fatalf("cold-low: fired=%d last=%+v, want one low-absolute at window 2", fired, last)
	}
	// A 20dB engine sits above the absolute floor and never alarms cold.
	d2 := NewCollapseDetector()
	for i := 0; i < 10; i++ {
		if got := d2.Observe(wstat(i, 20)); len(got) != 0 {
			t.Fatalf("20dB cold alarmed: %+v", got)
		}
	}
}

func TestCollapseDetectorSkipsInactiveWindows(t *testing.T) {
	d := NewCollapseDetector()
	silent := ERLEWindowStat{Index: 0, ActiveFrames: 0}
	for i := 0; i < 5; i++ {
		if got := d.Observe(silent); len(got) != 0 {
			t.Fatalf("inactive window alarmed: %+v", got)
		}
	}
	// Two low windows, an inactive one, then the third low window: the
	// streak survives the gap and fires on the third low window.
	_ = d.Observe(wstat(1, 9))
	_ = d.Observe(wstat(2, 9))
	_ = d.Observe(silent)
	got := d.Observe(wstat(3, 9))
	if len(got) != 1 || got[0].Reason != "low-absolute" {
		t.Fatalf("streak across inactive windows: %+v, want low-absolute", got)
	}
}

func TestGateTiersCounts(t *testing.T) {
	eps := []GateEpisode{
		{StartSample: 0, Samples: 16000, RmsDbfs: -50, AsrText: ""},       // nothing passes
		{StartSample: 32000, Samples: 8000, RmsDbfs: -20, AsrText: ""},    // energy only
		{StartSample: 64000, Samples: 8000, RmsDbfs: -50, AsrText: "为什么"}, // text only
		{StartSample: 96000, Samples: 8000, RmsDbfs: -20, AsrText: "我是"},  // both
	}
	g := GateTiers(eps, -30)
	if g.Starts != 4 || g.EnergyPass != 2 || g.TextPass != 2 || g.AllPass != 1 {
		t.Fatalf("GateTiers = %+v, want 4/2/2/1", g)
	}
	if Rate10Min(3, 18) != 100 {
		t.Fatalf("Rate10Min(3, 18s) = %f, want 100", Rate10Min(3, 18))
	}
	if Rate10Min(5, 0) != 0 {
		t.Fatalf("Rate10Min on zero duration = %f, want 0", Rate10Min(5, 0))
	}
}

func TestFarClockLoopsWithContinuousPhase(t *testing.T) {
	// 3 samples at 3kHz "mix rate" against the 16k seam rate: the fixed
	// -point step is fractional, so the wraparound interpolation is exercised
	// within a few hundred samples.
	const mixRate = 3000
	mix := []float32{1, 2, 4}
	f := NewFarClock(mix, mixRate)
	// At pos p (samples): x = p*mixRate/16000, idx = floor(x) mod 3.
	for i := 0; i < 1000; i++ {
		x := float64(i) * mixRate / audio.TargetRate
		idx := int(math.Floor(x)) % len(mix)
		frac := x - math.Floor(x)
		want := float32(float64(mix[idx])*(1-frac) + float64(mix[(idx+1)%len(mix)])*frac)
		if got := f.Next(); math.Abs(float64(got-want)) > 1e-6 {
			t.Fatalf("sample %d: got %f, want %f", i, got, want)
		}
	}
	// A single-sample loop is constant, not a division by zero.
	f1 := NewFarClock([]float32{0.5}, mixRate)
	for i := 0; i < 10; i++ {
		if v := f1.Next(); v != 0.5 {
			t.Fatalf("single-sample loop sample %d = %f, want 0.5", i, v)
		}
	}
	// NextFrame fills exactly.
	f2 := NewFarClock(mix, mixRate)
	buf := make([]float32, 7)
	f2.NextFrame(buf)
	if buf[0] != 1 {
		t.Fatalf("NextFrame first sample = %f, want 1", buf[0])
	}
}
