package audio

import "math"

// The Path-C barge-in combination gate (ticket 26 re-open, first work item):
// the product form of the gate chain the 59-w-r1 leg measured on the rig.
//
// Static tiers (the 59-w-r1 shape): an episode passes when it carries VAD
// (the caller's speech Detector start - this package never segments audio
// itself), enough energy, and a non-empty ASR transcription supplied by the
// caller. The static energy threshold alone is NOT the product gate: the
// 59-v3 verdict pinned both constants - the static -30 dBFS cut zeroes false
// triggers on its sample set but eats -36 dBFS quiet speech and leaves only
// 2.2-5.6 dB of margin on the collapse-state artifacts, while the loose -45
// dBFS cut holds the quiet speech but lets collapse-state artifacts through
// at 11.5 events/10min against 5.4 on the good state.
//
// The dynamic gating closes that gap with the three collapse criteria the
// 59-w-r1 leg productized (its longrun.go section 5, all pure audio-clock
// quantities, causal, no future samples):
//
//  1. sliding-ERLE sudden-drop alarm (median -6dB / absolute 15dB /
//     sustained 3 windows) - the CollapseDetector below, constants pinned to
//     the same values the rig leg validated on 75 events;
//  2. sliding correlation against the far reference, floor 0.8 - the
//     reference-free capture-autocorrelation discriminator's product form
//     (the rig leg's own conversion note): the product always has the true
//     far (it is playing the TTS), so the capture-vs-far correlation replaces
//     the loop-period self-correlation, and one sliding quantity covers both
//     "foreign energy present" and "alignment strength" (the 60b head at
//     strength 0.65 made a whole-run collapse immediately visible);
//  3. either channel sustained marks the window as UNTRUSTED: while untrusted,
//     the energy cut tightens from the nominal -45 to the pinned -30.
//
// The tighten/untrusted semantics survive the double-talk case the rig leg
// called out ("the product must survive"): a real barge-in depresses ERLE and
// fires the alarm, but real speech measured -26.70 dBFS RMS (control-26) and
// passes the tightened -30 cut; the collapse-state artifacts it eats measured
// -32.1 dBFS and below. On the loose -45 axis this is exactly the reading the
// dispatch asks for: collapse-state false triggers drop to the good-state
// level or below, and removing the dynamic gating restores the 11.5/10min
// baseline (the reverse criterion).
//
// Everything here is deterministic pure DSP over mono float32 samples at the
// 16k seam rate, sample-count driven (never wall-clock), goroutine-free.
// Constants carry the evidence citations; nothing is tunable away from the
// measured values without editing them in the open. Numbers:
// docs/evidence/s1/59wr1-longrun-gates.md sections 4-5,
// docs/evidence/s1/59-v3-verdict.md section 3.

// BargeinFrameSamples is the gate's frame, 32ms at the 16k seam rate - the
// frame the baseline ERLE ruler aggregates on.
const BargeinFrameSamples = 512

// BargeinWindowSamples is the gate's reporting window: 1s at the seam rate.
const BargeinWindowSamples = 16000

// The collapse-detector tuning, pinned to the values the 59-w-r1 leg derived
// from the double peak (good 29.0-30.7dB vs collapsed 8.88-8.91dB) and
// validated on 75 events: 15dB absolute separates the peaks with margin,
// 6dB below a rolling median catches mid-run losses above the floor, three
// consecutive 1s windows ride out word-gap flicker.
const (
	BargeinCollapseDropDb     = 6
	BargeinCollapseLowAbsDb   = 15
	BargeinCollapseSustain    = 3
	bargeinCollapseBaseWindow = 8
	bargeinCollapseMinArm     = 3
)

// BargeinFarCorrFloor is the sliding far-correlation floor: below it the
// capture carries foreign energy (the rig leg's ground-truth calibration -
// strength 0.92 was enough to drive ERLE to 15; 0.8 is loose in that
// direction, as that leg recorded honestly).
const BargeinFarCorrFloor = 0.8

// The two energy constants of the 59-v3 verdict, section 3: the loose
// nominal cut (quiet speech reachable) and the tightened cut the dynamic
// gating applies while a collapse channel is untrusted.
const (
	BargeinNominalEnergyDbfs   = -45
	BargeinTightenedEnergyDbfs = -30
)

// BargeinWindowStat is one closed 1s window's far-active power aggregation.
type BargeinWindowStat struct {
	Index        int     // window index; starts at sample Index * BargeinWindowSamples
	NearPow      float64 // sum of near^2 over far-active frames
	ResPow       float64 // sum of residual^2 over far-active frames
	ActiveFrames int     // 32ms frames whose (lag-shifted) far RMS passed the cut
}

// ERLEdB reads the window's ERLE. No far-active frames (a silent span) is
// not a number: the caller skips such windows - absence of echo is not
// collapse. Perfect cancellation is +Inf.
func (w BargeinWindowStat) ERLEdB() float64 {
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

// SlidingERLE is the streaming per-window ERLE aggregator: it receives the
// captured stream and the AEC output in order (any chunking) plus the
// LAG-SHIFTED far reference (echo-source time, the baseline convention that
// decides whether a frame carries echo), and folds them into 1s windows of
// 32ms frames. Memory is bounded by the current frame and window whatever
// the run length; closed windows append to Windows.
type SlidingERLE struct {
	farThreshold float64

	farBuf, nearBuf, resBuf []float32
	cur                     BargeinWindowStat
	curSamples              int
	Windows                 []BargeinWindowStat
}

// NewSlidingERLE arms the accumulator with the 0.1x far-RMS frame cut the
// baseline ERLEdB thresholds on.
func NewSlidingERLE(farLoopRMS float64) *SlidingERLE {
	return &SlidingERLE{farThreshold: 0.1 * farLoopRMS}
}

// Observe folds one chunk (any length) of the (farActivity, near, residual)
// streams. Windows close on sample count alone, so gaps in far activity
// never stall the curve.
func (a *SlidingERLE) Observe(farActivity, near, residual []float32) {
	a.farBuf = append(a.farBuf, farActivity...)
	a.nearBuf = append(a.nearBuf, near...)
	a.resBuf = append(a.resBuf, residual...)
	for len(a.farBuf) >= BargeinFrameSamples {
		frameFar := a.farBuf[:BargeinFrameSamples]
		if frameRMSBargein(frameFar) > a.farThreshold {
			a.cur.ActiveFrames++
			for k := 0; k < BargeinFrameSamples; k++ {
				n := float64(a.nearBuf[k])
				r := float64(a.resBuf[k])
				a.cur.NearPow += n * n
				a.cur.ResPow += r * r
			}
		}
		a.curSamples += BargeinFrameSamples
		if a.curSamples >= BargeinWindowSamples {
			a.cur.Index = len(a.Windows)
			a.Windows = append(a.Windows, a.cur)
			a.cur = BargeinWindowStat{}
			a.curSamples = 0
		}
		rest := len(a.farBuf) - BargeinFrameSamples
		copy(a.farBuf, a.farBuf[BargeinFrameSamples:])
		copy(a.nearBuf, a.nearBuf[BargeinFrameSamples:])
		copy(a.resBuf, a.resBuf[BargeinFrameSamples:])
		a.farBuf = a.farBuf[:rest]
		a.nearBuf = a.nearBuf[:rest]
		a.resBuf = a.resBuf[:rest]
	}
}

// frameRMSBargein is the frame RMS over a float32 span.
func frameRMSBargein(f []float32) float64 {
	if len(f) == 0 {
		return 0
	}
	var sum float64
	for _, v := range f {
		sum += float64(v) * float64(v)
	}
	return math.Sqrt(sum / float64(len(f)))
}

// CollapseAlarm is one fired collapse alarm, emitted once when the streak
// reaches the sustain length; the streak continuing does not re-fire.
type CollapseAlarm struct {
	WindowIndex int
	Reason      string  // "drop" or "low-absolute"
	WindowERLE  float64 // the first window of the streak that completed it
	Baseline    float64 // the detector's baseline at fire time (drop only; NaN on the absolute floor)
}

// CollapseDetector is the sliding-ERLE sudden-drop alarm (criterion 1). It
// consumes one closed window at a time (causal: only the past) and fires two
// alarm shapes, each once per sustained streak:
//
//   - "drop": the window fell dropDb below the median of the recent good
//     windows (a converged engine losing convergence mid-run);
//   - "low-absolute": the window sits below lowAbsDb outright.
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

// NewCollapseDetector arms the detector with the pinned tuning.
func NewCollapseDetector() *CollapseDetector {
	return &CollapseDetector{
		dropDb:   BargeinCollapseDropDb,
		lowAbsDb: BargeinCollapseLowAbsDb,
		sustain:  BargeinCollapseSustain,
	}
}

// Observe folds one window and returns any alarm that just completed its
// streak.
func (d *CollapseDetector) Observe(w BargeinWindowStat) []CollapseAlarm {
	db := w.ERLEdB()
	if math.IsNaN(db) {
		return nil
	}
	d.base = append(d.base, db)
	if len(d.base) > bargeinCollapseBaseWindow {
		d.base = d.base[len(d.base)-bargeinCollapseBaseWindow:]
	}
	var alarms []CollapseAlarm
	drop := len(d.base) >= bargeinCollapseMinArm && db <= medianBargein(d.base)-d.dropDb
	low := db < d.lowAbsDb
	if drop || low {
		d.run++
		if d.run == d.sustain {
			reason := "low-absolute"
			baseline := math.NaN()
			if drop {
				reason = "drop"
				baseline = medianBargein(d.base)
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

// Sustained reports whether a collapse streak is currently unbroken (the
// gate's tighten condition: from the streak completing its sustain length
// until a window fails both alarm shapes).
func (d *CollapseDetector) Sustained() bool {
	return d.run >= d.sustain
}

// medianBargein of a small slice (the detector's baseline ring). Averaging
// the two middle values keeps it stable on even counts.
func medianBargein(vs []float64) float64 {
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

// FarCorrelationDetector is criteria 2+3 in their product form: the sliding
// correlation of the capture against the lag-shifted far reference, one
// value per closed 1s window. A value below the floor means foreign energy
// is in the capture (the correlation is the alignment strength itself, so
// one quantity discriminates both directions the rig leg split with two
// instruments). It consumes the same frame-aligned streams as SlidingERLE.
type FarCorrelationDetector struct {
	lag int // lag-shift applied to the far stream before correlation

	farBuf, nearBuf []float32
	ring            []float32 // lazily built lag ring
	ringPos         int
	ab, an, bn      float64 // cross/self power accumulators over the open window
	pending         int     // samples folded into the open window
	last, floor     float64
	Values          []float64 // one entry per closed window
}

// NewFarCorrelationDetector arms the detector with the far lag (0 = none).
func NewFarCorrelationDetector(lag int) *FarCorrelationDetector {
	return &FarCorrelationDetector{lag: lag, floor: BargeinFarCorrFloor, last: math.NaN()}
}

// Observe folds one chunk of the (far, capture) streams. Windows close on
// sample count alone. Far samples are lag-shifted internally (a ring of
// `lag` samples delays the far stream into the correlation input), so the
// caller feeds the raw render-time far.
func (d *FarCorrelationDetector) Observe(far, capture []float32) {
	d.farBuf = append(d.farBuf, far...)
	d.nearBuf = append(d.nearBuf, capture...)
	// Ring of `lag` samples delaying the far stream (lag==0 shortcuts).
	delay := d.delayRing()
	pos := d.ringPos
	for len(d.farBuf) >= BargeinFrameSamples {
		for k := 0; k < BargeinFrameSamples; k++ {
			v := d.farBuf[k]
			var shifted float32
			if delay == nil {
				shifted = v
			} else {
				shifted = delay[pos]
				delay[pos] = v
				pos = (pos + 1) % len(delay)
			}
			n := float64(d.nearBuf[k])
			s := float64(shifted)
			d.ab += s * n
			d.an += s * s
			d.bn += n * n
			d.pending++
		}
		if d.pending >= BargeinWindowSamples {
			c := math.NaN()
			if d.an > 0 && d.bn > 0 {
				c = d.ab / math.Sqrt(d.an*d.bn)
			}
			d.Values = append(d.Values, c)
			d.last = c
			d.ab, d.an, d.bn = 0, 0, 0
			d.pending = 0
		}
		rest := len(d.farBuf) - BargeinFrameSamples
		copy(d.farBuf, d.farBuf[BargeinFrameSamples:])
		copy(d.nearBuf, d.nearBuf[BargeinFrameSamples:])
		d.farBuf = d.farBuf[:rest]
		d.nearBuf = d.nearBuf[:rest]
	}
	d.ringPos = pos
}

// Untrusted reports whether the most recent closed window sits below the
// floor (foreign energy present / alignment collapsed). Windows with no
// comparable signal (all-silent spans) are NaN and do not mark untrusted.
func (d *FarCorrelationDetector) Untrusted() bool {
	return !math.IsNaN(d.last) && d.last < d.floor
}

// ObserveValue folds one already-computed window correlation (NaN when the
// span has no comparable signal). This is the window-level input for callers
// that aggregate upstream instead of feeding raw frames.
func (d *FarCorrelationDetector) ObserveValue(v float64) {
	d.Values = append(d.Values, v)
	d.last = v
}

// delayRing lazily materializes the lag ring; nil when lag==0.
func (d *FarCorrelationDetector) delayRing() []float32 {
	if d.lag <= 0 {
		return nil
	}
	if d.ring == nil {
		d.ring = make([]float32, d.lag)
	}
	return d.ring
}

// BargeinEpisode is one VAD utterance episode handed over for a barge-in
// decision. The caller owns segmentation (speech Detector + Endpointer, the
// product shape) and ASR; this package only judges. StartSample is the
// seam-rate sample index of the start event; RmsDbfs is the episode RMS in
// dBFS over the closing event's audio; AsrText is the episode transcription.
type BargeinEpisode struct {
	StartSample int
	Samples     int
	RmsDbfs     float64
	AsrText     string
}

// BargeinOptions configures a BargeinGate. The zero value is unusable; use
// NewBargeinGate, which pins every field to the measured constants.
type BargeinOptions struct {
	EnergyThresholdDbfs float64 // nominal episode-energy cut (loose: quiet speech reachable)
	TightenedEnergyDbfs float64 // energy cut while a collapse channel is untrusted
	DynamicGating       bool    // false = the gate is the static tier chain only (the reverse-criterion switch)
	LagSamples          int     // render-to-capture lag the far reference is shifted by (0 = none)
	FarRMS              float64 // RMS of one far loop (arms the 0.1x frame cut; <=0 disables the ERLE channel)
	FarCorrelationFloor float64 // sliding far-correlation floor
}

// BargeinDecision is one episode's gate reading.
type BargeinDecision struct {
	Pass                bool    // VAD start (given) AND energy AND text under the applied cut
	EnergyPass          bool    // episode RMS over the applied cut
	TextPass            bool    // non-empty transcription
	ThresholdDbfs       float64 // the cut actually applied (nominal, or tightened under dynamic gating)
	TightenActive       bool    // a collapse channel was untrusted at the episode's start window
	SuppressedByDynamic bool    // the tightened cut is what kept the episode out
}

// BargeinGate is the Path-C barge-in combination gate: the static tier chain
// (VAD given, energy, ASR text) plus the three-criteria dynamic gating. The
// caller feeds the AEC-path streams frame-aligned (far reference, capture,
// AEC output) through Feed as they arrive, and hands episodes through
// Evaluate when the Endpointer closes them. State machine inputs are
// sample-counts only.
type BargeinGate struct {
	opts    BargeinOptions
	erle    *SlidingERLE
	det     *CollapseDetector
	fcorr   *FarCorrelationDetector
	ring    []float32 // lag ring delaying the far stream into the ERLE channel
	ringPos int
	wins    []bool // per closed window: a collapse channel was untrusted
}

// NewBargeinGate arms the gate with the measured constants.
func NewBargeinGate(o BargeinOptions) *BargeinGate {
	if o.EnergyThresholdDbfs == 0 {
		o.EnergyThresholdDbfs = BargeinNominalEnergyDbfs
	}
	if o.TightenedEnergyDbfs == 0 {
		o.TightenedEnergyDbfs = BargeinTightenedEnergyDbfs
	}
	if o.FarCorrelationFloor == 0 {
		o.FarCorrelationFloor = BargeinFarCorrFloor
	}
	g := &BargeinGate{opts: o, det: NewCollapseDetector()}
	if o.FarRMS > 0 {
		g.erle = NewSlidingERLE(o.FarRMS)
	}
	g.fcorr = NewFarCorrelationDetector(o.LagSamples)
	g.fcorr.floor = o.FarCorrelationFloor
	return g
}

// Feed advances the dynamic channels with one chunk of the AEC path. far is
// the render-time far reference (the gate applies the lag shift internally,
// the rig pass's ring delay); near is the capture; residual is the AEC
// output. The far-correlation channel runs first so that window i's verdict
// is in place when the ERLE window i closes on the same sample boundary;
// the ERLE channel only folds when all three streams are present and armed.
func (g *BargeinGate) Feed(far, near, residual []float32) {
	n := len(near)
	if n == 0 {
		return
	}
	if len(far) > n {
		far = far[:n]
	}
	if len(residual) > n {
		residual = residual[:n]
	}
	if far != nil {
		g.fcorr.Observe(far, near)
	}
	if g.erle != nil && far != nil && residual != nil && len(far) == n && len(residual) == n {
		shifted := far
		if g.opts.LagSamples > 0 {
			if g.ring == nil {
				g.ring = make([]float32, g.opts.LagSamples)
			}
			shifted = make([]float32, n)
			for k, v := range far {
				shifted[k] = g.ring[g.ringPos]
				g.ring[g.ringPos] = v
				g.ringPos = (g.ringPos + 1) % len(g.ring)
			}
		}
		g.erle.Observe(shifted, near, residual)
	}
	for len(g.wins) < len(g.erle.Windows) {
		w := g.erle.Windows[len(g.wins)]
		g.det.Observe(w)
		i := len(g.wins)
		g.wins = append(g.wins, g.det.Sustained() || g.fcorrUntrustedAt(i))
	}
	// Windows closed by the correlation channel while the ERLE channel is
	// not armed (FarRMS unset) still mark their slot untrusted.
	for len(g.wins) < len(g.fcorr.Values) {
		i := len(g.wins)
		g.wins = append(g.wins, g.fcorrUntrustedAt(i))
	}
}

// fcorrUntrustedAt reports the far-correlation verdict for window i if its
// value has closed yet.
func (g *BargeinGate) fcorrUntrustedAt(i int) bool {
	if i >= len(g.fcorr.Values) {
		return false
	}
	v := g.fcorr.Values[i]
	return !math.IsNaN(v) && v < g.fcorr.floor
}

// Evaluate judges one episode. The tighten state is read at the episode's
// start window (the barge-in decision moment): a window index of
// StartSample / BargeinWindowSamples. An episode that starts before any
// window closed is judged at the nominal cut (nothing observed yet - no
// alarm can be sustained on zero windows).
func (g *BargeinGate) Evaluate(ep BargeinEpisode) BargeinDecision {
	d := BargeinDecision{ThresholdDbfs: g.opts.EnergyThresholdDbfs}
	w := ep.StartSample / BargeinWindowSamples
	if w < 0 {
		w = 0
	}
	if g.opts.DynamicGating && w < len(g.wins) && g.wins[w] {
		d.TightenActive = true
		d.ThresholdDbfs = g.opts.TightenedEnergyDbfs
	}
	d.EnergyPass = ep.RmsDbfs > d.ThresholdDbfs
	d.TextPass = ep.AsrText != ""
	d.Pass = d.EnergyPass && d.TextPass
	if !d.Pass && d.TightenActive && ep.RmsDbfs > g.opts.EnergyThresholdDbfs && ep.AsrText != "" {
		d.SuppressedByDynamic = true
	}
	return d
}

// ObserveWindow is the gate's window-level seam: a caller that aggregates
// upstream (an engine's own supervision numbers, or an offline re-cut over
// recorded windows) hands over one closed window's ERLE stat and the
// far-correlation value measured over the same span (NaN when not
// measurable - silent spans, or a reference-free window). The decision path
// (Evaluate) is the same one raw Feed drives.
func (g *BargeinGate) ObserveWindow(w BargeinWindowStat, farCorr float64) {
	g.det.Observe(w)
	g.fcorr.ObserveValue(farCorr)
	g.wins = append(g.wins, g.det.Sustained() || (!math.IsNaN(farCorr) && farCorr < g.fcorr.floor))
}

// Windows reports the per-window untrusted flags (the readings surface for
// the offline re-cut and the integration panel).
func (g *BargeinGate) UntrustedWindows() []bool {
	out := make([]bool, len(g.wins))
	copy(out, g.wins)
	return out
}
