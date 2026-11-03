package speech

import (
	"net/http"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// This file is the shared vocabulary of the ticket 61 cloud voice cascade.
// The C9 seam in engines.go says local sherpa and cloud providers are
// interchangeable; making that true on a network fixes four shapes once so
// cloudasr.go, cloudtts.go and cloudchain.go do not each invent their own:
//
//   - CloudEndpoint: one catalog pair ("provider/model" from a voice chain)
//     resolved to everything an HTTP audio call needs (base URL, key, voice,
//     price card). The catalog is the truth source (ticket 05); the builder
//     in cloudcatalog.go produces these, the engines only consume.
//   - Failover events with the same shape as llm.FailoverEvent: SPEC-03 sec
//     3.1 demands "quota-exhaustion switches must produce a user-visible
//     event" for the text chain and SPEC-03's voice row carries the same
//     never-silent rule for the chains, so both report through one
//     vocabulary -- ticket 196's lesson about two state vocabularies in one
//     product applies to event vocabularies too.
//   - QuotaGate: the pre-dispatch check (ticket 44: an over-quota voice
//     model advances the chain BEFORE dispatch, never into a negative).
//   - VoiceUsage: the C23 metering record, carrying the catalog's
//     audio_in/audio_out price fields so audio lands priced.
//
// What this file deliberately does not do: touch a wire (HTTP shapes live
// with their engines), change a seam (engines.go owns those), or read
// config.toml (cloudcatalog.go adapts a loaded snapshot to these types).

// CloudChainKind names which direction of the voice cascade an event, a
// state or a usage record belongs to. Two values, mirroring the two config
// keys voice.cloud_asr_chain / voice.cloud_tts_chain.
type CloudChainKind string

const (
	// ChainASR is the voice.cloud_asr_chain direction.
	ChainASR CloudChainKind = "asr"
	// ChainTTS is the voice.cloud_tts_chain direction.
	ChainTTS CloudChainKind = "tts"
)

// localSherpaName is the terminal fallback's event identity, spelled exactly
// like config/catalog.go's reserved provider name so one reader tells chain
// members and the fallback apart. It never appears inside a chain: the
// validator rejects it there (the local cascade is always tried last,
// outside the chains).
const localSherpaName = "local-sherpa"

// ChainKey renders provider+model in the "provider/model" element spelling
// the catalog chains use. config's splitChainElement is the parse side of
// this exact format; rendering here keeps chain event names and config
// elements in one vocabulary.
func ChainKey(provider, model string) string {
	return provider + "/" + model
}

// KeyResolver resolves config api_key_ref values ("dpapi:<id>" | "env:NAME").
// *secret.Store satisfies it directly, and so does llm.KeyResolver -- the
// interface is re-declared here so the voice leg does not depend on the
// agent kernel just to name a signature. Contract, same as llm's: plaintext
// exists only inside CloudEndpoint.APIKey, never in a log line or an error
// detail.
type KeyResolver interface {
	Resolve(ref string) (string, error)
}

// HTTPDoer is the one-method shape of *http.Client the cascade uses. Tests
// substitute a recording client; production passes the default proxy-aware
// client built by cascadeClient below -- the same transport rules the C5
// adapters use (llm.NewDefaultHTTPClient), because "国内可达" is an
// environment fact of this product, not a per-protocol preference.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// CloudEndpoint is one fully resolved cloud voice target: a catalog pair
// plus everything its audio endpoint needs. cloudcatalog.go produces these
// from a loaded config snapshot; engines and chains consume them read-only.
type CloudEndpoint struct {
	// Provider and Model are the catalog pair; ChainKey renders them and
	// every event/usage record carries that spelling.
	Provider string
	Model    string
	// BaseURL is the provider's API root, e.g. "https://api.stepfun.com/v1".
	// The audio paths hang off it the way the OpenAI-compatible audio
	// endpoints do (cloudasr.go / cloudtts.go name the exact ones).
	BaseURL string
	// APIKey is the resolved plaintext key. Never logged, never in an error
	// detail (the String method below omits it, matching llm.Endpoint).
	APIKey string
	// ExtraHeaders are the provider's compat.extra_headers, appended
	// verbatim to every request of this endpoint.
	ExtraHeaders map[string]string

	// Voice selects the TTS voice id; empty = provider default. The panel's
	// voice dropdown reaches it through voice.tts.voice.
	Voice string
	// Speed is the TTS rate multiplier (voice.tts.speed, hot-tier); 0 reads
	// as 1.0. ASR endpoints ignore it.
	Speed float64
	// SampleRate is the rate of the PCM a TTS response carries. 0 = the
	// OpenAI audio/speech convention (24 kHz, mono, s16le). The render side
	// resamples from whatever TtsEngine.SampleRate reports.
	SampleRate int
	// Language optionally pins the ASR language hint; "" = model default.
	Language string

	// Billing and the price fields ride along for the usage record: audio
	// lands in the C23 ledger priced by the catalog's own audio_in /
	// audio_out card, without a second catalog read. A zero card reads as
	// free/unpriced exactly like the catalog defines it.
	Billing       string
	PriceAudioIn  float64
	PriceAudioOut float64
	// Capabilities echoes the two catalog bits the builder validated
	// (audio_in for an ASR endpoint, audio_out for a TTS endpoint). It
	// exists so a probe can re-check the DECLARATION against the PROBE
	// without reaching back into config; nothing else reads it.
	DeclaresAudioIn  bool
	DeclaresAudioOut bool
}

// String renders the endpoint WITHOUT the key (log-safe).
func (e CloudEndpoint) String() string {
	return ChainKey(e.Provider, e.Model) + " " + e.BaseURL
}

// ttsPCMDefaultRate is 24000: the rate response_format "pcm" carries on the
// OpenAI-compatible audio/speech wire the first candidates (StepAudio 2.5
// ASR / StepAudio 3 TTS) speak. It is a wire fact, not a knob, which is why
// it is reached through an endpoint default rather than a config key --
// same argument as asrFeatureDim for sherpa.
const ttsPCMDefaultRate = 24000

// ttsSampleRate is the effective PCM rate: the endpoint's override or the
// wire's own convention.
func (e CloudEndpoint) ttsSampleRate() int {
	if e.SampleRate > 0 {
		return e.SampleRate
	}
	return ttsPCMDefaultRate
}

// ttsSpeed is the effective rate multiplier; 0 (unset) reads as 1.0, the
// config default for voice.tts.speed.
func (e CloudEndpoint) ttsSpeed() float64 {
	if e.Speed > 0 {
		return e.Speed
	}
	return 1.0
}

// applyAuth sets the Authorization header when the endpoint carries a key.
// Keyless endpoints (local self-hosted gateways on the LAN) send no header
// rather than an empty Bearer, same as the openai-chat adapter.
func (e CloudEndpoint) applyAuth(req *http.Request) {
	if e.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.APIKey)
	}
	for k, v := range e.ExtraHeaders {
		req.Header.Set(k, v)
	}
}

// VoiceFailoverEvent is one hop of a cloud chain: cloud to cloud, the
// terminal hop to the local fallback, or a quota skip that happens BEFORE
// anything is dispatched. Emitted synchronously, never silent -- the same
// contract llm.ChainRunner.OnFailover carries for text_chain, with the same
// field meanings.
type VoiceFailoverEvent struct {
	Chain   CloudChainKind
	Index   int    // element being LEFT (or SKIPPED, for a quota skip)
	From    string // "provider/model" being left
	To      string // next element tried; localSherpaName at the terminal hop
	Cause   string // D37 class spelling (observe.ErrorClass) or quotaCauseSkip
	Detail  string // log-safe reason; never the key, never the audio
	Attempt int    // 1-based element position
	Total   int    // cloud chain length (the local fallback is not counted)
}

// VoiceProviderState is the privacy indicator D16 extends to the cascade:
// while a cloud engine is ACTIVE, audio (or the text to be spoken) leaves
// the machine and the host must surface that state; when the local fallback
// holds the line, the privacy posture is restored. The chain reports the
// fact; the ball/panel layer owns the rendering (Conversation red-ring rules
// unchanged -- SPEC-03 sec 3.1's voice row says surface the state, and the
// state machine is not this package's business).
type VoiceProviderState struct {
	Chain  CloudChainKind
	Active string // "provider/model", or localSherpaName on the fallback
	// Cloud is the boolean the UI reads: true exactly when this chain's
	// active engine takes the exchange off the machine.
	Cloud bool
}

// quotaCauseSkip is the Cause spelling for a pre-dispatch quota skip. It is
// deliberately not a D37 class: nothing failed, the dispatch never happened
// -- and "budget" (the class the 100% line maps to on the TEXT side) would
// read as a failure that already cost something. One vocabulary, exact.
const quotaCauseSkip = "quota"

// QuotaGate is the pre-dispatch quota check. Ticket 44's anti-negative rule
// applies to voice the same way as to text: an over-quota model is skipped
// before its turn in the chain. observe/quota_44's ledger is the real
// judge; the composition root wires the adapter. nil gate = allow
// everything (no quota configured).
type QuotaGate interface {
	// Allow reports whether provider/model may be dispatched right now.
	// ok=false says skip it; reason is the user-visible sentence (the same
	// QuotaVerdict.Notice text the ledger produces for text models).
	Allow(provider, model string) (ok bool, reason string)
}

// quotaGateFunc adapts a bare function to QuotaGate.
type quotaGateFunc struct {
	fn func(provider, model string) (bool, string)
}

func (q quotaGateFunc) Allow(provider, model string) (bool, string) {
	if q.fn == nil {
		return true, ""
	}
	return q.fn(provider, model)
}

// QuotaFunc wraps fn as a QuotaGate (host convenience, UsageFunc's twin).
func QuotaFunc(fn func(provider, model string) (bool, string)) QuotaGate {
	return quotaGateFunc{fn: fn}
}

// VoiceUsage is one billable cloud exchange, metered through C23. The two
// audio counters come from the audio clock (samples over the rate the audio
// ran at), never wall time (D42#9). Tokens are the provider's own usage
// numbers when the response carried any: 0 means "not reported", and a
// ledger must not read that as "free" unless the card says so.
type VoiceUsage struct {
	Chain CloudChainKind
	// Key is the "provider/model" that served the exchange.
	Key string
	// Billing is the catalog billing mode (pay-per-token | plan).
	Billing string
	// AudioSecondsIn is the utterance an ASR call received; AudioSecondsOut
	// the speech a TTS call produced.
	AudioSecondsIn  float64
	AudioSecondsOut float64
	// InputTokens / OutputTokens are the provider's usage block when
	// present; 0 = the response carried none.
	InputTokens  int
	OutputTokens int
	// PriceAudioIn / PriceAudioOut are the catalog price card fields of
	// the model that served the call, copied onto the record so the ledger
	// can price audio without a second catalog read.
	PriceAudioIn  float64
	PriceAudioOut float64
	// FirstFrame is how long the caller waited for the first TTS chunk
	// (SPEC-04 sec 8's measured quantity). It is provider-dependent by
	// contract -- the ticket forbids inventing an SLO gate here -- so it
	// travels as a reading, never as a verdict. ASR usage carries 0.
	FirstFrame string // duration in ms as text? no -- see field type below.
}

// UsageSink receives every successful cloud exchange. nil = the host does
// not meter (tests, or a cloud-less deployment); chains and engines check
// for nil and never block on a sink.
type UsageSink interface {
	RecordVoiceUsage(VoiceUsage)
}

// funcUsageSink adapts a bare function to UsageSink.
type funcUsageSink struct{ fn func(VoiceUsage) }

func (f funcUsageSink) RecordVoiceUsage(u VoiceUsage) {
	if f.fn != nil {
		f.fn(u)
	}
}

// UsageFunc wraps fn as a UsageSink, mirroring the observe.PanicSink style.
func UsageFunc(fn func(VoiceUsage)) UsageSink { return funcUsageSink{fn: fn} }

// ChainCallbacks groups the host-facing observers a chain reports through.
// Every field may be nil. All callbacks run synchronously on the owner
// goroutine (D38b: the chain owns no goroutine of its own that could absorb
// a slow host -- cloud ASR's transport goroutine reports state to the chain
// under a lock, it does not call these), and must not block long.
type ChainCallbacks struct {
	// OnFailover receives every hop (cloud to cloud, quota skip, cloud to
	// local fallback). Called before the next element acts.
	OnFailover func(VoiceFailoverEvent)
	// OnProviderState receives the ACTIVE state whenever the chain's head
	// changes: the initial pick, every step, the local fallback.
	OnProviderState func(VoiceProviderState)
	// OnUsage is the C23 meter, fed after each successful exchange.
	OnUsage UsageSink
	// Quota is the pre-dispatch gate.
	Quota QuotaGate
}

// ChainOptions are the constructor-time knobs the cloud engines and chains
// share. The zero value is valid: default proxy-aware HTTP client, no
// callbacks, transport goroutines on observe.Default detached from any task
// root (the engine cancels its own on Close).
type ChainOptions struct {
	// Client overrides the HTTP transport (tests); nil = the default
	// proxy-aware client. When set, ProxyMode/ProxyURL are ignored: exactly
	// one of the two paths is ever live.
	Client HTTPDoer
	// ProxyMode / ProxyURL come from the [net] section and feed the default
	// client (system/manual/none -- the same rules the C5 adapters use).
	ProxyMode string
	ProxyURL  string
	// Registry receives the cloud ASR transport goroutine; nil =
	// observe.Default so the watchdog sees one picture.
	Registry *observe.Registry
	// Root optionally joins the spawned transport's pending counter and
	// supplies its cancellation (the voice loop's task root). nil = the
	// engine runs its transport detached and stops it on Close.
	Root *observe.Root
	// Callbacks is the chain-level observer group (engines built through
	// the catalog builders receive the chain's own copy).
	Callbacks ChainCallbacks
}

// registry returns the registry to spawn on (never nil).
func (o ChainOptions) registry() *observe.Registry {
	if o.Registry != nil {
		return o.Registry
	}
	return observe.Default
}

// resolveClient returns the HTTP doer for a cascade object: the test/host
// override, or the default client built with the same proxy rules as the
// C5 adapters (llm.NewDefaultHTTPClient -- one transport policy for every
// outbound HTTP leg of the product).
func (o ChainOptions) resolveClient() HTTPDoer {
	if o.Client != nil {
		return o.Client
	}
	return llm.NewDefaultHTTPClient(o.ProxyMode, o.ProxyURL)
}
