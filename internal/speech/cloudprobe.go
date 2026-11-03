package speech

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Capability probes for the two voice bits (ticket 61: "audio_in/audio_out
// capability probes defined here feed provider_health (ticket 11's probe
// framework)"). The persistence half is provider_health.probe_json (SPEC-02
// schema v2; memory.ProbeFlags names audio_in/audio_out verbatim); ticket
// 11's writer owns the write. What this file owns is the REQUEST side that
// llm/probe.go deliberately left unimplementable -- its ProbeAudio case says
// the word: "the audio case is registered with its response check and an
// explicit NotImplementable note; the request definition lands together with
// the voice-cascade tickets (15/60/61)".
//
// What the probes do and do NOT prove. A probe drives one REAL canned
// exchange through the same engines the product uses -- so a pass means: the
// endpoint accepts this dialect, the key authorizes it, the audio path
// works. A pass does NOT mean the model is accurate: recognition quality is
// the CER harness's claim (ticket 15's domain), voice quality is the owner's
// P7 score sheet, and neither is a probe's job to fake. The audio_in probe
// therefore does not require a non-empty transcript (a pure tone carries no
// words -- demanding text there would probe the audio CONTENT, not the
// capability); it requires the exchange to complete without a classified
// failure. The audio_out probe requires at least one audio chunk back.

// CapAudioIn / CapAudioOut name the two probeable voice capabilities,
// spelled exactly like the provider_health probe_json keys (and like
// memory.ProbeFlags' Set/Get switch) so ticket 11's writer has no second
// vocabulary to translate through.
const (
	CapAudioIn  = "audio_in"
	CapAudioOut = "audio_out"
)

// AllVoiceProbeCapabilities lists the two capabilities in catalog order.
func AllVoiceProbeCapabilities() []string {
	return []string{CapAudioIn, CapAudioOut}
}

// VoiceProbeResult is one probe outcome. The shape mirrors llm.ProbeResult
// (capability/ok/latency_ms/detail/at) because provider_health rows are
// written once from both legs and one vocabulary beats two (ticket 196).
type VoiceProbeResult struct {
	Capability string `json:"capability"`
	OK         bool   `json:"ok"`
	LatencyMS  int64  `json:"latency_ms"`
	Detail     string `json:"detail,omitempty"`
	At         time.Time
}

// probeLatencyMS reads an elapsed monotonic duration (D42#9: never wall-clock
// differences) as the whole milliseconds provider_health stores. A sub-
// millisecond exchange still reports 1: that column is a MEASURED reading,
// and truncating a live endpoint's round trip to 0 would render it as "no
// reading at all" next to failed probes that legitimately carry 0.
func probeLatencyMS(d time.Duration) int64 {
	ms := d.Milliseconds()
	if ms < 1 {
		return 1
	}
	return ms
}

// ProbeAudioIn drives one canned recognition session against ep: a 200ms
// 440Hz tone at the seam rate, one utterance, and reads the outcome. OK is
// "completed without a classified failure" -- including the honest empty
// final (the tone has no words; an endpoint that 200s and closes is alive).
func ProbeAudioIn(ctx context.Context, ep CloudEndpoint, opts ChainOptions) VoiceProbeResult {
	res := VoiceProbeResult{Capability: CapAudioIn, At: observe.NowWallUTC()}
	engine, err := NewCloudASR(ep, opts)
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	defer engine.Close()

	start := time.Now()
	engine.Begin()
	engine.Feed(probeTone())
	text := engine.Finish()
	res.LatencyMS = probeLatencyMS(time.Since(start))

	if cause := engine.LastFailure(); cause != nil {
		res.Detail = cause.Error()
		return res
	}
	res.OK = true
	res.Detail = "transcript runes=" + strconv.Itoa(len([]rune(text)))
	return res
}

// ProbeAudioOut synthesizes one canned phrase and requires at least one
// audio chunk back. The phrase is spoken aloud into the probe only -- never
// logged -- and its length is seconds-scale.
func ProbeAudioOut(ctx context.Context, ep CloudEndpoint, opts ChainOptions) VoiceProbeResult {
	res := VoiceProbeResult{Capability: CapAudioOut, At: observe.NowWallUTC()}
	engine, err := NewCloudTTS(ep, opts)
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	defer engine.Close()

	start := time.Now()
	chunks, samples := 0, 0
	err = engine.Speak("喂", func(pcm []float32) error {
		chunks++
		samples += len(pcm)
		return nil
	})
	res.LatencyMS = probeLatencyMS(time.Since(start))

	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if chunks == 0 || samples == 0 {
		res.Detail = "endpoint answered with zero audio chunks"
		return res
	}
	res.OK = true
	res.Detail = "chunks=" + strconv.Itoa(chunks) + " samples=" + strconv.Itoa(samples) +
		" first_frame_ms=" + strconv.FormatInt(engine.FirstFrame().Milliseconds(), 10)
	return res
}

// RunVoiceProbeSuite executes every declared voice capability against one
// endpoint and returns the results, skipping (never faking) the ones the
// endpoint's catalog declaration does not carry. The caller maps OK=false +
// Detail into provider_health; the declaration/observed mismatch check is
// ticket 11's, not this file's.
func RunVoiceProbeSuite(ctx context.Context, ep CloudEndpoint, opts ChainOptions) []VoiceProbeResult {
	var out []VoiceProbeResult
	if ep.DeclaresAudioIn {
		out = append(out, ProbeAudioIn(ctx, ep, opts))
	}
	if ep.DeclaresAudioOut {
		out = append(out, ProbeAudioOut(ctx, ep, opts))
	}
	return out
}

// probeTone is the canned utterance: 200ms of 440Hz at 0.4 amplitude,
// float32 at the seam rate. A tone (not silence) is used deliberately: it
// keeps any provider-side energy gate from ending the session before the
// request ever carries audio, and it costs 6.4KB.
func probeTone() []float32 {
	n := audio.TargetRate / 5 // 200ms
	out := make([]float32, n)
	for i := range out {
		out[i] = 0.4 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(audio.TargetRate)))
	}
	return out
}
