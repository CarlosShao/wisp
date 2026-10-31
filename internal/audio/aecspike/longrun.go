package aecspike

import (
	"math"

	"github.com/CarlosShao/wisp/internal/audio"
)

// The long-window measurement math of the 59-w-r1 leg: sub-sample lag
// estimation, the streaming per-window ERLE accumulator, the in-run collapse
// detector (sliding-ERLE sudden-drop alarm), the barge-in gate tier
// combinator, and the phase-continuous far-reference clock that keeps a
// looped far fixture aligned to the captured stream across a whole long run.
// Everything here is deterministic pure DSP over mono float32 samples, so
// the unit tests pin every constant on every GOOS. The windows-only rig and
// the gated drivers live in longrun_windows*.go (aecd3-r1 baseline numbers:
// docs/evidence/s1/aec3r1-remeasure.md; the double-peak verdict this leg
// measures against: docs/evidence/s1/59-v2-verdict.md).

// FractionalLag refines an integer echo lag to sub-sample resolution by
// parabolic interpolation of the sampleScore correlation peak around lag.
// The returned offset lies in [-0.5, 0.5] under the near[i] ~ far[i-lag]
// convention: a positive offset means the echo arrives offset samples LATER
// than the integer lag. The aec3-r1 double peak (30.7dB vs 8.9dB on
// indistinguishable captures) made the render/loopback sub-sample phase the
// leading collapse-suspect; this is the instrument that reads it.
func FractionalLag(far, near []float32, lag int) float64 {
	if lag <= 0 || lag >= len(far)-1 || len(near) == 0 {
		return 0
	}
	alpha := sampleScore(far, near, lag-1)
	beta := sampleScore(far, near, lag)
	gamma := sampleScore(far, near, lag+1)
	denom := alpha - 2*beta + gamma
	if denom == 0 {
		return 0
	}
	delta := 0.5 * (alpha - gamma) / denom
	if delta > 0.5 {
		delta = 0.5
	}
	if delta < -0.5 {
		delta = -0.5
	}
	return delta
}

// erleFrameSamples is the ERLE frame, 32ms at the 16k seam rate: the same
// frame the baseline ERLEdB uses, so accumulator aggregates stay on the
// baseline ruler.
const erleFrameSamples = 512

// erleWindowSamples is the reporting window: 1s at the seam rate. The long
// -run ERLE curve is one point per window.
const erleWindowSamples = 16000

// ERLEWindowStat is one reporting window's far-active power aggregation:
// near power vs AEC-residual power over the window's far-active frames.
type ERLEWindowStat struct {
	Index        int     // window index; StartSample = Index * erleWindowSamples
	NearPow      float64 // sum of near^2 over far-active frames
	ResPow       float64 // sum of residual^2 over far-active frames
	ActiveFrames int     // 32ms frames whose (lag-shifted) far RMS passed the cut
}

// ERLEdB reads the window's ERLE. No far-active frames (a silent span) is
// not a number: the caller skips such windows rather than treating them as
// collapse. Perfect cancellation over the frames is +Inf.
func (w ERLEWindowStat) ERLEdB() float64 {
	if w.ActiveFrames == 0 {
		return math.NaN()
	}
	if w.ResPow <= 0 {
		return math.Inf(1)
	}
	if w.NearPow <= 0 {
		return math.Inf(-1)
	}
	return 10 * math.Log10(w.NearPow/w.ResPow)
}

// ERLEAccum is the streaming per-window ERLE aggregator. It receives the
// captured stream and the AEC output in order (any chunking) plus the
// LAG-SHIFTED far reference (the echo-source-time far, which is what decides
// whether a frame carries echo to cancel - the baseline ERLEdB convention),
// and folds them into 1s windows of 32ms frames. Memory is bounded by the
// current frame and window regardless of run length; the caller reads the
// closed windows out of Windows as they appear.
type ERLEAccum struct {
	farThreshold float64 // frame RMS cut: 0.1 x the far fixture's own RMS
	Windows      []ERLEWindowStat

	farBuf, nearBuf, resBuf []float32
	cur                     ERLEWindowStat
	curSamples              int
}

// NewERLEAccum takes the far fixture's RMS (over one loop - a gapless loop
// makes that identical to the whole-run RMS the baseline ERLEdB thresholds
// on) and arms the accumulator with the baseline's 0.1x cut.
func NewERLEAccum(farLoopRMS float64) *ERLEAccum {
	return &ERLEAccum{farThreshold: 0.1 * farLoopRMS}
}

// Observe folds one chunk (any length) of the (farActivity, near, residual)
// streams. Windows close on sample count alone, so gaps in far activity
// never stall the curve.
func (a *ERLEAccum) Observe(farActivity, near, residual []float32) {
	a.farBuf = append(a.farBuf, farActivity...)
	a.nearBuf = append(a.nearBuf, near...)
	a.resBuf = append(a.resBuf, residual...)
	for len(a.farBuf) >= erleFrameSamples {
		frameFar := a.farBuf[:erleFrameSamples]
		if frameRMS(frameFar) > a.farThreshold {
			a.cur.ActiveFrames++
			for k := 0; k < erleFrameSamples; k++ {
				n := float64(a.nearBuf[k])
				r := float64(a.resBuf[k])
				a.cur.NearPow += n * n
				a.cur.ResPow += r * r
			}
		}
		a.curSamples += erleFrameSamples
		if a.curSamples >= erleWindowSamples {
			a.cur.Index = len(a.Windows)
			a.Windows = append(a.Windows, a.cur)
			a.cur = ERLEWindowStat{}
			a.curSamples = 0
		}
		rest := len(a.farBuf) - erleFrameSamples
		copy(a.farBuf, a.farBuf[erleFrameSamples:])
		copy(a.nearBuf, a.nearBuf[erleFrameSamples:])
		copy(a.resBuf, a.resBuf[erleFrameSamples:])
		a.farBuf = a.farBuf[:rest]
		a.nearBuf = a.nearBuf[:rest]
		a.resBuf = a.resBuf[:rest]
	}
}

// SpanERLE aggregates the windows [from, to) into one ERLE: the same
// far-active power-sum ratio the baseline ERLEdB reports for a whole
// recording, restricted to the span. from/to are window indices.
func (a *ERLEAccum) SpanERLE(from, to int) float64 {
	if from < 0 {
		from = 0
	}
	if to > len(a.Windows) {
		to = len(a.Windows)
	}
	var nearPow, resPow float64
	frames := 0
	for _, w := range a.Windows[from:to] {
		nearPow += w.NearPow
		resPow += w.ResPow
		frames += w.ActiveFrames
	}
	if frames == 0 {
		return math.NaN()
	}
	if resPow <= 0 {
		return math.Inf(1)
	}
	return 10 * math.Log10(nearPow/resPow)
}

// CollapseDetector is the in-run collapse alarm of the 59-w-r1 leg: the
// product-usable prototype of "detect a converged-then-collapsed AEC while
// it is running". It consumes one ERLE window at a time (causal: only the
// past) and fires two alarm shapes, each once per sustained streak:
//
//   - "drop": the window fell dropDb below the median of the recent good
//     windows (a converged engine losing convergence mid-run);
//   - "low-absolute": the window sits below lowAbsDb outright (the cold
//     collapse shape - the 8.9dB peak against the 29-30.7dB one, which this
//     cut separates by construction).
//
// Windows without far-active frames are skipped without resetting a streak
// (the echo is absent there; absence is not recovery).
type CollapseDetector struct {
	dropDb   float64
	lowAbsDb float64
	sustain  int

	base []float64 // the last baseWindow active-window ERLEs
	run  int       // consecutive windows satisfying an alarm condition
}

// CollapseAlarm is one fired alarm. It is emitted once, when the streak
// reaches the sustain length; the streak continuing does not re-fire.
type CollapseAlarm struct {
	WindowIndex int
	Reason      string  // "drop" or "low-absolute"
	WindowERLE  float64 // the first window of the streak that completed it
	Baseline    float64 // the detector's baseline at fire time (drop only)
}

// collapseDetectorTuning is the leg's working cut set, derived from the
// double peak the 59-v2 verdict recorded (good 29.0-30.7dB, collapsed
// 8.88-8.91dB): 15dB absolute separates the peaks with a wide margin, 6dB
// below a rolling median catches mid-run losses well above the absolute
// floor, and three consecutive 1s windows ride out word-gap flicker.
const (
	collapseDropDb        = 6
	collapseLowAbsDb      = 15
	collapseSustain       = 3
	collapseBaseWindow    = 8
	collapseMinArmWindows = 3
)

// NewCollapseDetector arms the detector with the leg's tuning.
func NewCollapseDetector() *CollapseDetector {
	return &CollapseDetector{
		dropDb:   collapseDropDb,
		lowAbsDb: collapseLowAbsDb,
		sustain:  collapseSustain,
	}
}

// Observe folds one window and returns any alarm that just completed its
// streak.
func (d *CollapseDetector) Observe(w ERLEWindowStat) []CollapseAlarm {
	db := w.ERLEdB()
	if math.IsNaN(db) {
		return nil
	}
	d.base = append(d.base, db)
	if len(d.base) > collapseBaseWindow {
		d.base = d.base[len(d.base)-collapseBaseWindow:]
	}
	var alarms []CollapseAlarm
	drop := len(d.base) >= collapseMinArmWindows && db <= median(d.base)-d.dropDb
	low := db < d.lowAbsDb
	if drop || low {
		d.run++
		if d.run == d.sustain {
			reason := "low-absolute"
			baseline := math.NaN()
			if drop {
				reason = "drop"
				baseline = median(d.base)
			}
			alarms = append(alarms, CollapseAlarm{
				WindowIndex: w.Index, Reason: reason,
				WindowERLE: db, Baseline: baseline,
			})
		}
	} else {
		d.run = 0
	}
	return alarms
}

// median of a small slice (the detector's baseline ring). Averaging the two
// middle values keeps it stable on even counts.
func median(vs []float64) float64 {
	if len(vs) == 0 {
		return math.NaN()
	}
	s := append([]float64(nil), vs...)
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
	mid := len(s) / 2
	if len(s)%2 == 1 {
		return s[mid]
	}
	return (s[mid-1] + s[mid]) / 2
}

// GateEpisode is one VAD utterance episode drawn from a gate material: the
// product shape is Detector + Endpointer over the AEC output, and each
// closing event hands over the utterance's own samples.
type GateEpisode struct {
	StartSample int     // seam-rate sample index of the start event
	Samples     int     // episode length in samples (from the closing event)
	RmsDbfs     float64 // episode RMS in dBFS over the closing event's audio
	AsrText     string  // product streaming ASR over the episode audio
}

// GateTierCount is the three-tier barge-in gate tally over one material.
// Tier 1 (VAD alone) fires on every start by construction; tier 2 adds the
// energy gate; tier 3 adds the ASR text gate (non-empty transcription).
// A false trigger on echo-only material is any tier that fires.
type GateTierCount struct {
	Starts     int
	EnergyPass int
	TextPass   int
	AllPass    int
}

// GateTiers tallies the episodes against the energy threshold.
func GateTiers(eps []GateEpisode, energyThresholdDbfs float64) GateTierCount {
	var g GateTierCount
	g.Starts = len(eps)
	for _, e := range eps {
		energy := e.RmsDbfs > energyThresholdDbfs
		text := e.AsrText != ""
		if energy {
			g.EnergyPass++
		}
		if text {
			g.TextPass++
		}
		if energy && text {
			g.AllPass++
		}
	}
	return g
}

// Rate10Min normalizes an event count over a duration to events per 10
// minutes (the false-trigger table's unit since 59-r1).
func Rate10Min(count int, seconds float64) float64 {
	if seconds <= 0 {
		return 0
	}
	return float64(count) * 600 / seconds
}

// FarClock produces the 16k far reference for a long run: the far fixture
// was rendered LOOPED at the mix rate, and the reference the AEC engine
// deserves at seam-frame i is the looped fixture at mix-position i*mixRate
// /16000 - a phase-continuous interpolation with wraparound, not a re-loop
// of the 16k fixture (whose loop length differs from the mix-rate loop by
// rounding and would drift the reference against the echo by one sample per
// loop iteration).
type FarClock struct {
	mix  []float32 // mono far at the mix rate, looping
	step int64     // (mixRate << 32) / 16000, 32.32 fixed point
	pos  int64
}

// NewFarClock arms the clock over one loop of the mix-rate mono far.
func NewFarClock(mixFar []float32, mixRate int) *FarClock {
	return &FarClock{
		mix:  mixFar,
		step: (int64(mixRate) << 32) / int64(audio.TargetRate),
	}
}

// Next returns the next seam-rate far sample, looping with continuous phase.
func (f *FarClock) Next() float32 {
	n := int64(len(f.mix))
	if n == 0 {
		return 0
	}
	if n == 1 {
		return f.mix[0]
	}
	ip := (f.pos >> 32) % n
	frac := uint64(f.pos & 0xffffffff)
	a := f.mix[ip]
	b := f.mix[(ip+1)%n]
	v := a + float32(frac)*(b-a)/float32(uint64(1)<<32)
	f.pos += f.step
	return v
}

// NextFrame fills buf with the next len(buf) far samples.
func (f *FarClock) NextFrame(buf []float32) {
	for i := range buf {
		buf[i] = f.Next()
	}
}
