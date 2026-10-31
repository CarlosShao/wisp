//go:build windows

package main

// The Path-C barge-in integration surface (ticket 26 re-open, the integration
// leg): the 389 combination gate plus the AEC3 engine, wired into the
// resident speech path per D47 sec 16.12 - "麦克风常开，AEC 消掉播报回声，
// 说话即插话".
//
// WHAT RUNS HERE. During one announcement the monitor owns the announcement
// window:
//
//   - the FAR reference is the announcement itself: speak()'s chunk callback
//     hands every synthesis chunk here (speakFar) before it reaches the
//     device, and the monitor resamples it to the 16k seam rate - D47's own
//     argument ("回声参考是我们自己渲染的 PCM") is what makes this reference
//     exact;
//   - the CAPTURE frames are the real C8 seam: while the monitor is armed the
//     capture gate's D47 conversation window (audio.HalfDuplexGate.
//     SetConversation, the internal/audio half of this integration) keeps
//     capture open across SetSpeaking(true), and the monitor's drain goroutine
//     is the one reader of rv.frames for the window (the consumer goroutine is
//     inside speak() for exactly that span - the single-reader discipline
//     holds by the same bracket that makes it true between announcements);
//   - the judge is the product chain the 59 legs measured: AEC3 cancels the
//     echo from the capture (render-then-capture 10ms lockstep, the stream
//     delay hint measured online by the envelope aligner - the product form
//     of the rig's delaytrack instrument), audio.BargeinGate.Feed folds the
//     (far, near, residual) streams online, the monitor's own silero session
//     segments the RESIDUAL, and an episode start is judged through the gate's
//     real Evaluate. D32:2315's shape is what makes the 400ms budget possible:
//     播放期挂起 ASR、400ms 门禁由 VAD 而非 ASR 满足 - the episode is judged at
//     its start (EnergyPass under the applied cut; the gate stays the single
//     threshold authority through Evaluate), and the text tier is the ASR's
//     business after the barge-in, not the trigger's;
//   - a pass fires the trigger once (TriggerInterrupt - the same channel the
//     Esc cancel key and the ball click use), which aborts the speaker in
//     flight through the 381-r2 abort hop.
//
// WHAT DEGRADES, HONESTLY. The AEC3 binding compiles only under
// `-tags webrtcap` (internal/audio/aec); everywhere else NewApm answers
// ErrUnavailable and the monitor stays inert - the announcement window then
// keeps today's exact D16 shape (the capture gate closes, no conversation
// window, no drain), which is the only honest posture without the AEC layer.
// The same nil-safety covers every other missing piece (no VAD model, no
// scope, no capture gate). Every degradation lands in the readings with its
// reason.
//
// WHAT THE REVERSE CRITERIA FLIP. Two switches exist for the integration
// tests (and only they set them):
//
//   - dynamicOff: the gate's own reverse switch (389 AC#3) - the static
//     loose chain judges, and the collapse-state artifact family the
//     tightened cut exists to eat passes again (the echo self-trigger);
//   - aecBypass: the judge runs on the raw capture - the echo reaches the
//     speech judge uncancelled, its episode spans are what the product mic
//     surface heard, and transcribing them is the verbatim-transcription
//     reproduction.
//
// Threads. arm/release run on the owner goroutine (inside speak, the same
// bracket as SetSpeaking); speakFar runs on the owner goroutine inside the
// synthesis callback; the drain runs on the C11 scope under the frozen
// temporary roster name "tts-infer" (TTS-time inference is exactly what it
// is). The far ring is guarded by its own mutex; the drain's processing
// state is drain-owned (one goroutine, no locks); the readings are written
// before the drain starts or after it joins, except the two atomics the
// drain writes live.
//
// The drain never waits on the far: when the ring does not yet cover the
// next block the capture keeps draining into the pending buffer (bounded by
// the announcement itself - the anchor opens with the first synthesis chunk
// already in the ring) and the next frame retries. Waiting in the frame loop
// would stall the D38d bounded channel and manufacture drops the monitor
// itself caused.

import (
	"context"
	"log/slog"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/speech"
)

// voiceBargeinConvGate is the conversation-window seam: the one method the
// monitor needs beyond the capture gate's SetSpeaking. Production hands the
// *audio.HalfDuplexGate (which carries SetConversation since this leg); test
// recorders that do not implement it leave the seam nil and the monitor runs
// without the window (the readings say so) - the frames they want judged are
// pushed into the seam directly.
type voiceBargeinConvGate interface {
	SetConversation(on bool)
}

// voiceBargeinRoster is the frozen temporary-roster name the drain is booked
// under (internal/observe/goroutine.go's TemporaryNames): TTS-time inference.
const voiceBargeinRoster = "tts-infer"

// The monitor's shape constants, all audio-clock quantities.
const (
	// voiceBargeinAlignFar is the far span the delay measurement runs on
	// (300ms of announcement), and voiceBargeinAlignMaxLag the search span
	// around the anchor (the waveOut buffer latency plus the acoustic/hardware
	// return lag - the rig measured 493 samples of the same quantity).
	voiceBargeinAlignFar    = 4800
	voiceBargeinAlignMaxLag = 4800
	// voiceBargeinAlignStrength is the trust floor for the measured alignment:
	// below it the far/capture relationship is not an echo and the lag answers
	// zero (the readings carry the strength either way).
	voiceBargeinAlignStrength = 0.2
	// voiceBargeinFarCap bounds the far ring (60s of announcement at the seam
	// rate - the announce path's own text cap cannot reach it).
	voiceBargeinFarCap = 60 * audio.TargetRate
	// voiceBargeinStreamCap bounds the recorded streams (the equivalence-grid
	// and margin readings' material) at 30s per stream per announcement.
	voiceBargeinStreamCap = 30 * audio.TargetRate
	// voiceBargeinEpisodeCap bounds the recorded episode spans (the M2
	// verbatim-transcription material), each at voiceBargeinSpanCap.
	voiceBargeinEpisodeCap = 4
	voiceBargeinSpanCap    = 16 * audio.TargetRate
	// voiceBargeinOnsetRing is the rolling residual the onset RMS reads from
	// (a start event's Total is the MinSpeech run plus at most one window).
	voiceBargeinOnsetRing = 2 * 16000
	// The residual noise gate's shape (the announcement-window segmenter).
	// The finding that pins it: the AEC'd residual during an announcement is
	// CONTINUOUS - the echo residue never lets the endpointer see 700ms of
	// silence, so the first residual episode stays open and every later
	// onset (the talker's) never produces a start event (measured: round 3
	// of the first (a) run judged one window-0 episode and went silent for
	// the rest of the window). The gate closes the steady floor and keeps
	// the dynamics: the floor is the MINIMUM window RMS over the last two
	// seconds (the syllable dips of anything speech-shaped keep it honest),
	// a window passes when it clears the floor by voiceBargeinGateStep, and
	// a passing window holds the gate open for voiceBargeinGateHold windows
	// (512ms - longer than MinSpeech, so a real utterance's inter-syllable
	// dips do not break its speech run). What the gate discriminates, on
	// this surface, measured: the AEC'd residue is flat (its floor ≈ its
	// level ⇒ gated out), the raw echo (the M2 bypass) and the talker carry
	// speech-shaped dynamics above their own dips ⇒ they pass. silero and
	// the Endpointer still run on the gated residual - the D32 shape keeps
	// its VAD; the gate only stops the echo floor from owning it.
	voiceBargeinFloorWindows = 64 // 2s of 32ms windows
	voiceBargeinGateStep     = 2.0
	voiceBargeinGateHold     = 16 // 512ms
)

// voiceBargeinEpisode is one judged episode's reading. Span is the judged
// residual audio, transient and in-memory only (D16: buffers are never
// persisted and never logged) - it is the M2 reading's material.
type voiceBargeinEpisode struct {
	StartSample   int     // on the gate's fed-sample clock
	Samples       int     // the opening run (MinSpeech-dominated)
	RmsDbfs       float64 // onset RMS over the opening run
	ThresholdDbfs float64 // the cut Evaluate applied
	TightenActive bool    // a collapse channel was untrusted at the start window
	EnergyPass    bool    // the onset cleared the applied cut
	Suppressed    bool    // the tightened cut is what kept it out
	Span          []float32
}

// voiceBargeinReadings is the post-release snapshot (copies).
type voiceBargeinReadings struct {
	Armed         bool
	Reason        string // why the monitor never armed (the honest posture)
	ConvWindow    bool   // the conversation window was actually flipped
	AecAvailable  bool   // the AEC3 binding was really in the build
	AecBypass     bool   // the M2 mutant ran the judge on raw capture
	DynamicOff    bool   // the M1 mutant judged with the static loose chain
	Lag           int    // measured render-to-capture lag (samples)
	LagStrength   float64
	FarSamples    int // far ring size at release
	FarCuts       int // far samples the ring cap dropped
	Processed     int // capture samples pushed through the judge
	FramesDrained int // C8 frames the monitor read for the window
	Episodes      []voiceBargeinEpisode
	Triggers      int
	StreamCapped  bool
	// Untrusted is the online gate's per-window verdict sequence (the
	// equivalence grid's online half: the same windows re-cut offline through
	// ObserveWindow must answer the same flags).
	Untrusted []bool
	// TriggerAt is when the monitor fired the interrupt trigger (the (a)
	// reading's t0; zero when no trigger fired).
	TriggerAt time.Time
}

// voiceBargeinMonitor is one announcement's barge-in judge. A zero value is
// not usable: build with newVoiceBargeinMonitor, arm per announcement, and
// always release.
type voiceBargeinMonitor struct {
	frames   <-chan []byte
	conv     voiceBargeinConvGate
	scope    *plugin.DisposalScope
	vadModel string
	trigger  func()

	// The reverse-criteria switches (test-only; production leaves them off).
	aecBypass  bool
	dynamicOff bool

	// armed state (owner goroutine, before the drain starts)
	cancel    context.CancelFunc
	handle    *observe.Handle
	det       *speech.Detector
	apm       *aec.Apm
	resampler *audio.Resampler

	// far ring (farMu guards it; the drain polls coverage, never waits)
	farMu    sync.Mutex
	far      []float32
	farDone  bool
	farCuts  int
	farTotal int

	// readings (written before the drain starts or after it joins, except
	// the two atomics the drain writes live)
	mu           sync.Mutex
	armed        bool
	reason       string
	convWindow   bool
	aecAvailable bool
	lag          int
	lagStrength  float64
	episodes     []voiceBargeinEpisode
	triggers     int
	processed    int
	untrusted    []bool
	streamFar    []float32
	streamNear   []float32
	streamRes    []float32
	streamCapped bool

	framesDrained atomic.Int64
	triggerAtNano atomic.Int64
}

// newVoiceBargeinMonitor wires the monitor to the leg's seams. Everything is
// nil-tolerant downstream: a monitor without a frame channel, a VAD model, a
// scope or a trigger answers every arm with its reason instead of arming.
func newVoiceBargeinMonitor(frames <-chan []byte, vadModel string, scope *plugin.DisposalScope, trigger func()) *voiceBargeinMonitor {
	m := &voiceBargeinMonitor{
		frames:   frames,
		vadModel: vadModel,
		scope:    scope,
		trigger:  trigger,
	}
	return m
}

// setReverseCriteria installs the integration tests' mutant switches. It
// must run before arm.
func (m *voiceBargeinMonitor) setReverseCriteria(aecBypass, dynamicOff bool) {
	if m == nil {
		return
	}
	m.aecBypass, m.dynamicOff = aecBypass, dynamicOff
}

// attachConvGate hands the monitor the capture gate's conversation window.
// The assembly calls it from attachCaptureGate with the same object; nil (a
// gate without SetConversation) degrades to no window.
func (m *voiceBargeinMonitor) attachConvGate(g voiceBargeinConvGate) {
	if m == nil {
		return
	}
	m.conv = g
}

// arm builds the engines and opens the announcement window. It answers
// whether the monitor is really judging; every refusal lands in the readings
// with its reason. Runs on the owner goroutine before SetSpeaking(true) so
// the conversation window precedes the close request.
func (m *voiceBargeinMonitor) arm() bool {
	if m == nil || m.armed {
		return false
	}
	switch {
	case m.frames == nil:
		m.reason = "没有 C8 帧通道（本进程没有采集腿）"
	case m.vadModel == "":
		m.reason = "VAD 模型路径为空（语音链未装配）"
	case m.scope == nil:
		m.reason = "没有 C11 作用域（语音链未启动）"
	case m.trigger == nil:
		m.reason = "没有打断触发器"
	default:
		m.reason = ""
	}
	if m.reason != "" {
		return false
	}
	det, err := speech.NewDetector(speech.DetectorConfig{ModelPath: m.vadModel})
	if err != nil {
		m.reason = "VAD 没起来：" + err.Error()
		return false
	}
	var apm *aec.Apm
	if !m.aecBypass {
		a, aerr := aec.NewApm(audio.TargetRate)
		if aerr != nil {
			_ = det.Close()
			m.reason = "AEC3 不可用（本构建没有 -tags webrtcap 的绑定）：" + aerr.Error()
			return false
		}
		apm = a
	}
	// The D47 conversation window: capture stays open across SetSpeaking(true)
	// and the monitor owns the announcement window's frames. Muted still wins
	// (the gate's own precedence), and without the window the monitor judges
	// nothing - the frames it would need never exist.
	if m.conv != nil {
		m.conv.SetConversation(true)
		m.convWindow = true
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.det = det
	m.apm = apm
	m.armed = true
	m.aecAvailable = apm != nil
	// Per-announcement reset: the monitor object is reused across the leg's
	// announcements, and each window's readings must be its own.
	m.mu.Lock()
	m.reason = ""
	m.lag, m.lagStrength = 0, 0
	m.episodes = nil
	m.triggers = 0
	m.processed = 0
	m.untrusted = nil
	m.streamFar, m.streamNear, m.streamRes = nil, nil, nil
	m.streamCapped = false
	m.mu.Unlock()
	m.framesDrained.Store(0)
	m.triggerAtNano.Store(0)
	m.farMu.Lock()
	m.far = m.far[:0]
	m.farDone = false
	m.farCuts = 0
	m.farTotal = 0
	m.farMu.Unlock()
	m.handle = m.scope.Go(voiceBargeinRoster, func(scopeCtx context.Context) {
		m.drain(ctx, scopeCtx)
	})
	slog.Info("voice: barge-in monitor armed for the announcement window",
		"roster", voiceBargeinRoster, "aec", m.aecAvailable, "aec_bypass", m.aecBypass,
		"dynamic_gating", !m.dynamicOff, "conv_window", m.convWindow)
	return true
}

// speakFar is the far-reference tap: speak()'s chunk callback hands every
// synthesis chunk here before it reaches the device. The chunk is resampled
// to the seam rate (the streaming resampler keeps absolute coordinates, so
// chunk boundaries never drift) and appended to the far ring. Runs on the
// owner goroutine.
func (m *voiceBargeinMonitor) speakFar(samples []float32, engineRate int) {
	if m == nil || !m.armed || len(samples) == 0 {
		return
	}
	if m.resampler == nil {
		if engineRate <= 0 {
			return
		}
		m.resampler = audio.NewResampler(engineRate, audio.TargetRate)
		if m.resampler == nil {
			return
		}
	}
	conv := m.resampler.Process(audio.FloatToPCM16(samples))
	far := make([]float32, len(conv))
	for i, v := range conv {
		far[i] = float32(v) / 32768
	}
	m.farMu.Lock()
	m.farTotal += len(far)
	if room := voiceBargeinFarCap - len(m.far); room > 0 {
		if len(far) > room {
			m.farCuts += len(far) - room
			far = far[:room]
		}
		m.far = append(m.far, far...)
	} else {
		m.farCuts += len(far)
	}
	m.farMu.Unlock()
}

// release closes the announcement window: the conversation window goes first
// (the D16 half-duplex close is back while SetSpeaking(true) still holds, so
// the frames behind the close are pre-release ones), then the drain is told
// the far is complete, cancelled, and joined, and the engines come down.
// Runs on the owner goroutine; speak() declares its defer so this runs
// BEFORE the capture gate's SetSpeaking(false).
func (m *voiceBargeinMonitor) release() {
	if m == nil || !m.armed {
		return
	}
	if m.conv != nil && m.convWindow {
		m.conv.SetConversation(false)
	}
	m.farMu.Lock()
	m.farDone = true
	m.farMu.Unlock()
	if m.cancel != nil {
		m.cancel()
	}
	if m.handle != nil {
		<-m.handle.Done()
	}
	m.closeEngines()
	m.armed = false
	r := m.readings()
	slog.Info("voice: barge-in monitor released",
		"processed_samples", r.Processed, "episodes", len(r.Episodes),
		"triggers", r.Triggers, "lag", r.Lag, "lag_strength", r.LagStrength,
		"frames_drained", r.FramesDrained, "far_samples", r.FarSamples,
		"far_cuts", r.FarCuts)
}

// closeEngines releases the engines the arm built (C11: sessions are
// released, never leaked). Runs after the drain joined.
func (m *voiceBargeinMonitor) closeEngines() {
	if m.apm != nil {
		_ = m.apm.Close()
		m.apm = nil
	}
	if m.det != nil {
		_ = m.det.Close()
		m.det = nil
	}
	m.resampler = nil
}

// readings snapshots the readings (copies; spans deep-copied).
func (m *voiceBargeinMonitor) readings() voiceBargeinReadings {
	if m == nil {
		return voiceBargeinReadings{}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	r := voiceBargeinReadings{
		Armed:         m.armed,
		Reason:        m.reason,
		ConvWindow:    m.convWindow,
		AecAvailable:  m.aecAvailable,
		AecBypass:     m.aecBypass,
		DynamicOff:    m.dynamicOff,
		Lag:           m.lag,
		LagStrength:   m.lagStrength,
		Triggers:      m.triggers,
		Processed:     m.processed,
		FramesDrained: int(m.framesDrained.Load()),
		StreamCapped:  m.streamCapped,
	}
	m.farMu.Lock()
	r.FarSamples = len(m.far)
	r.FarCuts = m.farCuts
	m.farMu.Unlock()
	for _, ep := range m.episodes {
		ep.Span = append([]float32(nil), ep.Span...)
		r.Episodes = append(r.Episodes, ep)
	}
	r.Untrusted = append([]bool(nil), m.untrusted...)
	if n := m.triggerAtNano.Load(); n != 0 {
		r.TriggerAt = time.Unix(0, n)
	}
	return r
}

// FramesDrained is the announcement-window frame count the monitor read; it
// is a field of the snapshot now, taken off the same atomic the drain
// advanced (read after the join, so the value is final).

// streams hands the recorded (far, near, residual) streams out - the
// equivalence-grid and margin readings' material. Nil slices when nothing
// was processed. Copies, so the caller cannot reach the monitor's state.
func (m *voiceBargeinMonitor) streams() (far, near, res []float32) {
	if m == nil {
		return nil, nil, nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]float32(nil), m.streamFar...),
		append([]float32(nil), m.streamNear...),
		append([]float32(nil), m.streamRes...)
}

// bargeinDrain holds the drain goroutine's whole processing state. One
// goroutine owns it from spawn to join - no locks on any of it.
type bargeinDrain struct {
	m           *voiceBargeinMonitor
	anchor      int
	anchorSet   bool
	capCount    int
	pending     []float32
	procSamples int
	lagDone     bool
	lag         int
	gate        *audio.BargeinGate
	vadBuf      []float32
	onsetRing   []float32
	eps         []speech.Endpoint
	triggered   bool
	floorRing   []float64
	floorN      int
	gateHold    int
}

// drain is the announcement window's frame reader and judge: the one reader
// of rv.frames while the consumer goroutine is inside speak(). Returns when
// the private ctx is cancelled (release) or the scope ctx is done (dispose).
func (m *voiceBargeinMonitor) drain(ctx, scopeCtx context.Context) {
	d := &bargeinDrain{m: m}
	defer func() {
		// The online gate's per-window verdicts are the equivalence grid's
		// online half; the drain is the only reader of the gate, and this is
		// its last act before the readings freeze.
		if d.gate != nil {
			m.mu.Lock()
			m.untrusted = d.gate.UntrustedWindows()
			m.mu.Unlock()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-scopeCtx.Done():
			return
		case frame, ok := <-m.frames:
			if !ok {
				return
			}
			m.framesDrained.Add(1)
			pcm, derr := audio.DecodeFrame(frame)
			if derr != nil {
				continue
			}
			d.capCount += len(pcm)
			// The anchor is the first frame read after far appeared (the first
			// synthesis chunk reached the tap): the announcement window starts
			// there, and the aligner's search span absorbs the residual offset
			// (device latency plus the return lag). Frames before it are read
			// and counted, not judged.
			if !d.anchorSet {
				m.farMu.Lock()
				hasFar := len(m.far) > 0
				m.farMu.Unlock()
				if hasFar {
					d.anchorSet = true
					d.anchor = d.capCount - len(pcm)
				} else {
					continue
				}
			}
			d.pending = append(d.pending, voiceSamples(pcm)...)
			if !d.lagDone {
				if !d.measure() {
					continue
				}
			}
			d.process()
		}
	}
}

// measure runs once per announcement: the render-to-capture lag online, then
// the gate armed with it. It answers whether the judge is ready yet.
func (d *bargeinDrain) measure() bool {
	m := d.m
	m.farMu.Lock()
	farLen := len(m.far)
	farDone := m.farDone
	m.farMu.Unlock()
	if farLen < voiceBargeinAlignFar {
		if farDone {
			// The announcement ended before 300ms of far existed: nothing to
			// align, nothing to judge.
			return true
		}
		return false
	}
	if d.capCount-d.anchor < voiceBargeinAlignFar+voiceBargeinAlignMaxLag {
		return false
	}
	m.farMu.Lock()
	farHead := append([]float32(nil), m.far[:voiceBargeinAlignFar]...)
	var farSum float64
	for _, v := range m.far[:voiceBargeinAlignFar] {
		farSum += float64(v) * float64(v)
	}
	farRMS := math.Sqrt(farSum / float64(voiceBargeinAlignFar))
	m.farMu.Unlock()
	nearWin := append([]float32(nil), d.pending[:voiceBargeinAlignFar+voiceBargeinAlignMaxLag]...)
	measured, strength := audio.AlignByEnvelopeCorrelation(farHead, nearWin, voiceBargeinAlignMaxLag)
	if measured < 0 || strength < voiceBargeinAlignStrength {
		measured = 0
	}
	d.lag = measured
	d.lagDone = true
	// The announcement window's operating point. The integration surface's
	// own first staged run measured the shape that pins it: under HEALTHY
	// matched-echo cancellation both dynamic channels read trusted (the
	// far-correlation tracks the echo, the ERLE sits ~23dB - no collapse
	// alarm), so the loose nominal -45 rides, and the engine's residual floor
	// on this surface (-42.11 dBFS measured) sits ABOVE it - the echo
	// self-triggered at the nominal cut. The S4 no-self-trigger validation
	// point 59-v3 §3 recorded is tier-3@-30 (「不自激」字面达标), and the
	// tightened cut is a margin exchange, not a suppression: real speech
	// measured -26.70 dBFS passes -30 (control-26), the residual floor does
	// not. So the window's gate arms AT the tightened level; the dynamic
	// channels still run and still mark windows (their readings stay live).
	// The M1 reverse constructs at the loose axis with the dynamic gating
	// off - exactly the pre-dynamic static shape whose 11.5/10min baseline
	// the 59-w rig measured.
	if m.dynamicOff {
		d.gate = audio.NewBargeinGate(audio.BargeinOptions{
			DynamicGating: false,
			LagSamples:    d.lag,
			FarRMS:        farRMS,
		})
	} else {
		d.gate = audio.NewBargeinGate(audio.BargeinOptions{
			EnergyThresholdDbfs: audio.BargeinTightenedEnergyDbfs,
			TightenedEnergyDbfs: audio.BargeinTightenedEnergyDbfs,
			DynamicGating:       true,
			LagSamples:          d.lag,
			FarRMS:              farRMS,
		})
	}
	m.mu.Lock()
	m.lag, m.lagStrength = d.lag, strength
	m.mu.Unlock()
	slog.Info("voice: barge-in alignment measured",
		"lag_samples", d.lag, "lag_ms", d.lag*1000/audio.TargetRate,
		"strength", math.Round(strength*1000)/1000,
		"far_rms_dbfs", math.Round((20*math.Log10(farRMS))*100)/100)
	return true
}

// process pushes whole 10ms blocks through the judge while the far ring
// covers the next block and pending capture holds one. Never waits: what
// cannot be judged yet stays buffered for the next frame's retry.
func (d *bargeinDrain) process() {
	m := d.m
	if d.gate == nil {
		return
	}
	frame := aec.FrameSamples(audio.TargetRate) // 160: the 10ms frame at the seam rate
	for {
		m.farMu.Lock()
		avail := len(m.far)
		m.farMu.Unlock()
		if avail < d.procSamples+frame {
			return
		}
		if len(d.pending) < frame {
			return
		}
		m.farMu.Lock()
		farBlock := append([]float32(nil), m.far[d.procSamples:d.procSamples+frame]...)
		m.farMu.Unlock()
		nearBlock := append([]float32(nil), d.pending[:frame]...)
		var res []float32
		if m.apm != nil {
			if err := m.apm.ProcessRender(farBlock); err != nil {
				slog.Warn("voice: barge-in AEC render refused", "err", err.Error())
				return
			}
			r, err := m.apm.ProcessCapture(nearBlock, d.lag*1000/audio.TargetRate)
			if err != nil {
				slog.Warn("voice: barge-in AEC capture refused", "err", err.Error())
				return
			}
			res = r
		} else {
			// aecBypass: the judge runs on the raw capture - the M2 shape.
			res = nearBlock
		}
		d.pending = d.pending[frame:]
		d.procSamples += frame
		d.gate.Feed(farBlock, nearBlock, res)
		d.recordStream(farBlock, nearBlock, res)
		d.onsetRing = append(d.onsetRing, res...)
		if len(d.onsetRing) > voiceBargeinOnsetRing {
			cut := len(d.onsetRing) - voiceBargeinOnsetRing
			copy(d.onsetRing, d.onsetRing[cut:])
			d.onsetRing = d.onsetRing[:voiceBargeinOnsetRing]
		}
		d.vadBuf = append(d.vadBuf, res...)
		for len(d.vadBuf) >= speech.WindowSamples {
			win := d.vadBuf[:speech.WindowSamples]
			d.vadBuf = d.vadBuf[speech.WindowSamples:]
			d.eps = m.det.Accept(d.noiseGate(win), d.eps[:0])
			for _, ep := range d.eps {
				if ep.Kind != speech.EndpointStart {
					continue
				}
				d.judge(ep)
			}
		}
		m.mu.Lock()
		m.processed = d.procSamples
		m.mu.Unlock()
	}
}

// noiseGate is the residual segmenter's floor gate: it returns the window
// when it carries dynamics above the tracked floor (or the hold says a run
// is open), and digital silence when it is the steady floor itself. The
// floor is the minimum window RMS over the trailing two seconds; the VAD
// downstream only ever sees what passed.
func (d *bargeinDrain) noiseGate(win []float32) []float32 {
	var sum float64
	for _, v := range win {
		sum += float64(v) * float64(v)
	}
	rms := math.Sqrt(sum / float64(len(win)))
	if d.floorRing == nil {
		d.floorRing = make([]float64, voiceBargeinFloorWindows)
	}
	d.floorRing[d.floorN%voiceBargeinFloorWindows] = rms
	d.floorN++
	floor := math.MaxFloat64
	for i, n := 0, min(d.floorN, voiceBargeinFloorWindows); i < n; i++ {
		if d.floorRing[i] < floor {
			floor = d.floorRing[i]
		}
	}
	if d.gateHold > 0 {
		d.gateHold--
	}
	if rms >= floor*voiceBargeinGateStep {
		d.gateHold = voiceBargeinGateHold
		return win
	}
	if d.gateHold > 0 {
		return win
	}
	return voiceBargeinSilence
}

// voiceBargeinSilence is the gate's closed shape: one window of digital
// silence, shared (read-only).
var voiceBargeinSilence = make([]float32, speech.WindowSamples)

// judge evaluates one episode start through the gate's real decision path
// and fires the trigger once. The episode carries no ASR text by design
// (D32:2315: the 400ms gate is VAD-satisfied; the text tier is the ASR's
// business after the barge-in) - Evaluate stays the single threshold
// authority for the applied cut.
func (d *bargeinDrain) judge(ep speech.Endpoint) {
	m := d.m
	n := ep.Total
	if n <= 0 || n > len(d.onsetRing) {
		n = len(d.onsetRing)
	}
	span := append([]float32(nil), d.onsetRing[len(d.onsetRing)-n:]...)
	var sum float64
	for _, v := range span {
		sum += float64(v) * float64(v)
	}
	rms := math.Sqrt(sum / float64(max(1, len(span))))
	rmsDbfs := math.Inf(-1)
	if rms > 0 {
		rmsDbfs = 20 * math.Log10(rms)
	}
	dd := d.gate.Evaluate(audio.BargeinEpisode{
		StartSample: d.procSamples - ep.Total,
		Samples:     ep.Total,
		RmsDbfs:     rmsDbfs,
		AsrText:     "",
	})
	read := voiceBargeinEpisode{
		StartSample:   d.procSamples - ep.Total,
		Samples:       ep.Total,
		RmsDbfs:       rmsDbfs,
		ThresholdDbfs: dd.ThresholdDbfs,
		TightenActive: dd.TightenActive,
		EnergyPass:    dd.EnergyPass,
		Suppressed:    dd.SuppressedByDynamic,
	}
	if len(span) > voiceBargeinSpanCap {
		span = span[:voiceBargeinSpanCap]
	}
	read.Span = span
	m.mu.Lock()
	if len(m.episodes) < voiceBargeinEpisodeCap {
		m.episodes = append(m.episodes, read)
	}
	m.mu.Unlock()
	slog.Info("voice: barge-in episode judged",
		"rms_dbfs", math.Round(rmsDbfs*100)/100, "threshold", dd.ThresholdDbfs,
		"tighten", dd.TightenActive, "energy_pass", dd.EnergyPass,
		"start_window", (d.procSamples-ep.Total)/audio.BargeinWindowSamples)
	if dd.EnergyPass && !d.triggered {
		d.triggered = true
		m.triggerAtNano.Store(time.Now().UnixNano())
		m.mu.Lock()
		m.triggers++
		m.mu.Unlock()
		slog.Info("voice: BARGE-IN triggered (stop the announcement)",
			"rms_dbfs", math.Round(rmsDbfs*100)/100, "threshold", dd.ThresholdDbfs,
			"tighten", dd.TightenActive, "onset_samples", ep.Total)
		if m.trigger != nil {
			m.trigger()
		}
	}
}

// recordStream appends one processed block to the readings' material (capped
// once, silently - the cap note lands in the evidence).
func (d *bargeinDrain) recordStream(far, near, res []float32) {
	m := d.m
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.streamCapped {
		return
	}
	if len(m.streamFar)+len(far) > voiceBargeinStreamCap {
		m.streamCapped = true
		return
	}
	m.streamFar = append(m.streamFar, far...)
	m.streamNear = append(m.streamNear, near...)
	m.streamRes = append(m.streamRes, res...)
}
