package config

// Full config.toml schema (D36 / SPEC-03 sec 3 + sec 3.1 LLM catalog v2, 2026-09-19
// user-approved supplement). The section/key/type/enum inventory below is
// contract-level: it must stay exactly in sync with the SPEC-03 sec 3 table and
// is pinned by tests (per-section valid/unknown-key/wrong-type sets).
//
// D36 rule 3: defaults live ONLY in the `default:"..."` struct tags. The TOML
// file is a serialization of this schema, never a second source of defaults
// (NewDefaults applies the tags; go-toml overlays the file on top).
//
// D36 rule 5: sensitive values are stored as references (`api_key_ref =
// "dpapi:<id>"` or `"env:NAME"`); no struct in this file has a plaintext
// `api_key` field, and none ever may (pinned by boundary_test.go, which also
// enforces the SPEC-03 sec 3.1 storage split: probe results / health / latency
// sampling live in SQLite provider_health (SPEC-02 schema v2, ticket 09),
// never here).

// SchemaVersionCurrent is the config.toml schema version this build writes
// and understands. History:
//
//	v1 - D36 draft model: [llm] with default_provider/fallback_provider/
//	     timeout/temperature and flat providers.{name}.{base_url,model,
//	     api_key_ref}; no text_chain, roles, model catalog or quotas.
//	v2 - 2026-09-19 LLM catalog supplement (SPEC-03 sec 3.1): text_chain,
//	     roles, providers.<name>.models.<id> with capabilities/billing/
//	     quota, [voice] realtime + cloud chains.
const SchemaVersionCurrent = 2

// Tier is the effect level of a config change (D36 three-tier semantics).
type Tier string

const (
	// TierHot applies immediately: the new value is visible through
	// Store.Config() as soon as the reload merge completes.
	TierHot Tier = "hot"
	// TierReload applies immediately AND emits a reload event so the owning
	// subsystem can reload (e.g. voice model swap triggers model load/unload).
	TierReload Tier = "reload"
	// TierRestart is NOT applied at runtime: the old value stays active until
	// process restart. The reload report carries the pending section so the
	// app can surface "restart required".
	TierRestart Tier = "restart"
)

// Direction is the security posture of a change inside a locked section
// (risk/fs/net/plugins). Loosening changes require L2 re-confirmation
// (D36 rule 1); tightening and neutral changes hot-apply with a log line.
type Direction string

const (
	DirNeutral Direction = ""
	DirLoosen  Direction = "loosen"
	DirTighten Direction = "tighten"
)

// Enum vocabularies. Values are validated strictly (SPEC-03 sec 2 rule 2:
// nothing invalid passes silently).
const (
	ThemeDark  = "dark"
	ThemeLight = "light"
	ThemeAuto  = "auto"

	ObserveDebug = "debug"
	ObserveInfo  = "info"
	ObserveWarn  = "warn"
	ObserveError = "error"

	ProxySystem = "system"
	ProxyNone   = "none"
	ProxyManual = "manual"

	OverBudgetPause = "pause"
	OverBudgetWarn  = "warn"

	ProtocolOpenAIChat      = "openai-chat"
	ProtocolOpenAIResponses = "openai-responses"
	ProtocolAnthropic       = "anthropic"

	BillingPayPerToken = "pay-per-token"
	BillingPlan        = "plan"

	// Thinking vocabulary shared by roles.<role>.thinking_intensity and
	// models.<id>.thinking_levels (SPEC-03 sec 3.1).
	ThinkingOff    = "off"
	ThinkingLow    = "low"
	ThinkingMedium = "medium"
	ThinkingHigh   = "high"
)

var (
	themes            = []string{ThemeDark, ThemeLight, ThemeAuto}
	observeLevels     = []string{ObserveDebug, ObserveInfo, ObserveWarn, ObserveError}
	proxyModes        = []string{ProxySystem, ProxyNone, ProxyManual}
	overBudgetModes   = []string{OverBudgetPause, OverBudgetWarn}
	protocols         = []string{ProtocolOpenAIChat, ProtocolOpenAIResponses, ProtocolAnthropic}
	billingModes      = []string{BillingPayPerToken, BillingPlan}
	thinkingIntensity = []string{ThinkingOff, ThinkingLow, ThinkingMedium, ThinkingHigh}
)

// Config is the whole config.toml (D36 section tree). Field order is the
// canonical marshal order; every field is emitted by Write (no omitempty) so
// the file always round-trips the full struct (GUI double-source guard).
type Config struct {
	// SchemaVersion is the file's schema version (SPEC-03 sec 4.4). It is
	// written by migration/load, never defaulted; hand-editing it down
	// triggers the migration chain on next load.
	SchemaVersion int `toml:"schema_version"`

	App     AppSection     `toml:"app"`
	Ball    BallSection    `toml:"ball"`
	Hotkey  HotkeySection  `toml:"hotkey"`
	Session SessionSection `toml:"session"`
	Voice   VoiceSection   `toml:"voice"`
	Audio   AudioSection   `toml:"audio"`
	LLM     LLMSection     `toml:"llm"`
	Agent   AgentSection   `toml:"agent"`
	Risk    RiskSection    `toml:"risk"` // locked section: loosening needs L2 re-confirm
	FS      FSSection      `toml:"fs"`   // locked section: same
	Net     NetSection     `toml:"net"`  // locked section: same
	Privacy PrivacySection `toml:"privacy"`
	Memory  MemorySection  `toml:"memory"`
	Panel   PanelSection   `toml:"panel"`
	Cost    CostSection    `toml:"cost"`
	// Plugins is excluded from plain toml.Marshal (`toml:"-"`): the dynamic
	// per-plugin sub-tables are appended by MarshalCanonical/marshalPlugins
	// and decoded through the raw-map absorber in parse.go (single code path,
	// round-trip pinned by tests). A `toml:"plugins"` tag here would make the
	// head marshal emit [plugins] with a literal "Entries" key AND the tail
	// emit it again - the double-emission bug the round-trip test pins.
	Plugins PluginsSection `toml:"-"`
	Models  ModelsSection  `toml:"models"`
	Observe ObserveSection `toml:"observe"`
}

// AppSection is [app]. Everything except theme is restart-tier (D36).
type AppSection struct {
	Language string `toml:"language" default:"zh-CN"`
	// Theme is the one hot-tier key of [app] (D36).
	Theme string `toml:"theme" default:"dark"`
	// Autostart is restart-tier (OS registration at boot path).
	Autostart bool `toml:"autostart" default:"false"`
	// SingleInstance is restart-tier (mutex name/registration at boot).
	SingleInstance bool `toml:"single_instance" default:"true"`
	// Portable is decided by portable.txt next to the exe (SPEC-03 sec 5.2),
	// never by config.toml: it is not serialized and a literal `portable`
	// key in the file is rejected as an unknown key. Populated by ticket 06
	// (env layout); consumers must treat it as read-only.
	Portable bool `toml:"-"`
}

// BallPos is ball.position (D36: position{x,y,monitor}).
type BallPos struct {
	// X/Y are the last user-dragged position; 0 means "let the app pick".
	X int `toml:"x"`
	// Y see X.
	Y int `toml:"y"`
	// Monitor names the monitor ("", the primary, by default).
	Monitor string `toml:"monitor"`
}

// BallSection is [ball]; entirely hot-tier.
type BallSection struct {
	// Size is the floating ball diameter in px, range 44-72 (SPEC-03 sec 3).
	Size int `toml:"size" default:"56"`
	// Position is the persisted window position.
	Position BallPos `toml:"position"`
	// OpacityIdle is the idle-state opacity (0..1).
	OpacityIdle float64 `toml:"opacity_idle" default:"0.35"`
	// ClickThrough lets mouse events pass through the idle ball.
	ClickThrough bool `toml:"click_through" default:"true"`
	// HideOnFullscreen hides the ball while a fullscreen window has focus.
	HideOnFullscreen bool `toml:"hide_on_fullscreen" default:"true"`
}

// HotkeySection is [hotkey]; hot-tier (hotkeys re-register on change).
// Empty string = binding unset (feature key disabled). `cancel` temporarily
// takes over during Confirming and is returned at session end (D36).
type HotkeySection struct {
	Summon string `toml:"summon"`
	Mute   string `toml:"mute"`
	Cancel string `toml:"cancel" default:"Esc"`
	Panel  string `toml:"panel"`
}

// SessionSection is [session]; hot-tier.
type SessionSection struct {
	WarmTimeoutSec      int `toml:"warm_timeout_sec" default:"90"`
	SettlingSec         int `toml:"settling_sec" default:"3"`
	ConversationIdleSec int `toml:"conversation_idle_sec" default:"30"`
}

// WakeWord is voice.wake_word. Thresholds are per-keyword confidences
// (0..1), keywords[i] pairs with thresholds[i]. Veto words immediately stop
// a recording session; they are hot-tier (D36: thresholds and veto words).
type WakeWord struct {
	Enabled bool `toml:"enabled" default:"false"`
	// Keywords are the KWS phrases, in priority order.
	Keywords []string `toml:"keywords"`
	// Thresholds holds one confidence threshold per keyword.
	Thresholds []float64 `toml:"thresholds"`
	// VetoWords abort listening the moment they are heard.
	VetoWords []string `toml:"veto_words" default:"取消,停下,别"`
}

// ASRConfig is voice.asr. Changing provider/model is reload-tier (model
// load/unload); the local sherpa cascade is the final voice fallback and is
// never a member of the cloud chains (SPEC-03 sec 3).
type ASRConfig struct {
	Provider string `toml:"provider" default:"local-sherpa"`
	Model    string `toml:"model"`
}

// TTSConfig is voice.tts. Provider/voice swaps are reload-tier; speed is hot.
type TTSConfig struct {
	Provider string  `toml:"provider" default:"local-sherpa"`
	Voice    string  `toml:"voice"`
	Speed    float64 `toml:"speed" default:"1.0"`
}

// AECConfig is voice.aec (echo cancellation for conversation mode).
type AECConfig struct {
	Enabled bool `toml:"enabled" default:"true"`
	// EchoRef selects the echo reference source; "self-render" loops the
	// app's own audio output back as the reference.
	EchoRef string `toml:"echo_ref" default:"self-render"`
}

// RealtimeConfig is voice.realtime (C32, ticket 60). The provider/model pair
// must exist in the llm catalog when realtime is enabled; api_key_ref follows
// the reference-only rule (D36 rule 5).
type RealtimeConfig struct {
	Enabled   bool   `toml:"enabled" default:"false"`
	Provider  string `toml:"provider"`
	Model     string `toml:"model"`
	APIKeyRef string `toml:"api_key_ref"`
	BaseURL   string `toml:"base_url"`
}

// VoiceSection is [voice]. Reload-tier: enabled, wake_word.enabled/keywords,
// asr.*, tts.provider/voice, conversation_mode, aec.*, realtime.* and both
// cloud chains (anything that re-shapes the audio/model pipeline). Hot-tier:
// wake_word.thresholds, wake_word.veto_words, tts.speed, punctuation.
type VoiceSection struct {
	Enabled bool `toml:"enabled" default:"true"`
	// WakeWord tuning; see WakeWord for the reload/hot split.
	WakeWord WakeWord `toml:"wake_word"`
	// ASR / TTS local engines (sherpa cascade, the terminal voice fallback).
	ASR ASRConfig `toml:"asr"`
	TTS TTSConfig `toml:"tts"`
	// Punctuation restores punctuation on ASR output.
	Punctuation bool `toml:"punctuation" default:"true"`
	// ConversationMode keeps the mic open between turns.
	ConversationMode bool `toml:"conversation_mode" default:"false"`
	// AEC echo cancellation.
	AEC AECConfig `toml:"aec"`
	// Realtime is the C32 realtime voice service head of the voice chain.
	Realtime RealtimeConfig `toml:"realtime"`
	// CloudASRChain is the ordered cloud ASR fallback ("provider/model"
	// elements), tried before the local sherpa cascade. Never contains
	// "local-sherpa" (validation rejects it: the local cascade is the
	// terminal fallback and lives outside the chains, SPEC-03 sec 3.1).
	CloudASRChain []string `toml:"cloud_asr_chain"`
	// CloudTTSChain is the cloud TTS counterpart of CloudASRChain.
	CloudTTSChain []string `toml:"cloud_tts_chain"`
}

// AudioSection is [audio]; hot-tier (audio device reopen on change).
type AudioSection struct {
	// InputDevice is the capture device name; "default" = OS default.
	InputDevice string `toml:"input_device" default:"default"`
	SampleRate  int    `toml:"sample_rate" default:"16000"`
	// HalfDuplex is HARD-CODED true (D36/D16): the value is read-only and
	// writing `false` in config.toml is a load error, never a setting.
	HalfDuplex bool `toml:"half_duplex" default:"true"`
	// MicMutedDefault starts every session muted.
	MicMutedDefault bool `toml:"mic_muted_default" default:"true"`
}

// Role is one llm.roles.<name> entry (SPEC-03 sec 3.1). provider+model together
// reference llm.providers.<provider>.models.<model>; a role with both empty
// is unset (feature falls back/disabled); exactly one set is an error.
// Temperature default (0.7) is shared by all roles; per-role call-time
// refinement belongs to the llm/agent modules, not to a second default store.
type Role struct {
	Provider string `toml:"provider"`
	Model    string `toml:"model"`
	// Temperature is the sampling temperature.
	Temperature float64 `toml:"temperature" default:"0.7"`
	// ThinkingIntensity is one of off|low|medium|high (enum enforced).
	ThinkingIntensity string `toml:"thinking_intensity" default:"off"`
	// MaxOutputTokens caps the response; 0 = provider default.
	MaxOutputTokens int `toml:"max_output_tokens"`
}

// Roles assigns models to the four fixed agent roles (SPEC-03 sec 3.1).
type Roles struct {
	Chat          Role `toml:"chat"`
	MemoryExtract Role `toml:"memory_extract"`
	Handoff       Role `toml:"handoff"`
	Summarize     Role `toml:"summarize"`
}

// RetryConfig is llm.retry.
type RetryConfig struct {
	Max       int `toml:"max" default:"3"`
	BackoffMS int `toml:"backoff_ms" default:"1000"`
}

// Compat carries per-provider interoperability switches (SPEC-03 sec 3.1).
// loose tolerates missing fields (stream_options/usage/tool_choice);
// allow_missing_usage accepts responses without usage blocks; extra_headers
// is appended verbatim to every request.
type Compat struct {
	Loose             bool              `toml:"loose" default:"false"`
	AllowMissingUsage bool              `toml:"allow_missing_usage" default:"false"`
	ExtraHeaders      map[string]string `toml:"extra_headers"`
}

// Capabilities is the declared capability bit-set of a catalog model
// (SPEC-03 sec 3.1). Declarations may be hand-filled but must be verified by
// probe (ticket 11); probe results themselves live in SQLite provider_health
// and deliberately have no field here (storage split).
type Capabilities struct {
	Text     bool `toml:"text"`
	Vision   bool `toml:"vision"`
	AudioIn  bool `toml:"audio_in"`
	AudioOut bool `toml:"audio_out"`
	Realtime bool `toml:"realtime"`
	Thinking bool `toml:"thinking"`
	FC       bool `toml:"fc"` // function calling
	Stream   bool `toml:"stream"`
}

// Price is the per-model price card. Units: micro-USD per 1M tokens
// (C23/ticket 44 consumes these to compute cost_micros; 0 = free/unknown).
type Price struct {
	In       float64 `toml:"in"`
	Out      float64 `toml:"out"`
	Cached   float64 `toml:"cached"`
	AudioIn  float64 `toml:"audio_in"`
	AudioOut float64 `toml:"audio_out"`
}

// ModelSpec is llm.providers.<name>.models.<model-id>: one catalog entry.
// New models default to all-false capabilities (probe-first, sec 3.1) and
// enabled=true; billing/quota mirror the provider-level fields.
type ModelSpec struct {
	// Display is the human-facing name; empty = model id.
	Display string `toml:"display"`
	// Capabilities: declared bits, verified by probe (ticket 11).
	Capabilities Capabilities `toml:"capabilities"`
	// ContextWindow in tokens; 0 = unknown.
	ContextWindow int `toml:"context_window"`
	// MaxOutputTokens in tokens; 0 = provider default.
	MaxOutputTokens int `toml:"max_output_tokens"`
	// ThinkingLevels enumerates the thinking levels this model supports
	// (off|low|medium|high; enum enforced on every element).
	ThinkingLevels []string `toml:"thinking_levels"`
	// Billing overrides the billing mode at model granularity.
	Billing string `toml:"billing" default:"pay-per-token"`
	// Price is the price card (micro-USD per 1M tokens).
	Price Price `toml:"price"`
	// QuotaDailyMicro / QuotaMonthlyMicro are per-model spend caps in
	// micro-USD; 0 = uncapped. Enforcement is a C23 pre-dispatch hook
	// (ticket 44); 100% hard-stops and walks the fallback chain.
	QuotaDailyMicro   int64 `toml:"quota_daily_micro"`
	QuotaMonthlyMicro int64 `toml:"quota_monthly_micro"`
	// Enabled removes the model from discovery/selection without deleting
	// its catalog entry.
	Enabled bool `toml:"enabled" default:"true"`
}

// Provider is llm.providers.<name> (SPEC-03 sec 3 + sec 3.1). Names matching a
// built-in preset (openai, anthropic, deepseek, qwen, zhipu, moonshot,
// siliconflow, openrouter, ollama, minimax, mimo, stepfun) inherit
// protocol/base_url defaults for any field left empty; unknown names must
// set protocol explicitly.
type Provider struct {
	// Protocol is openai-chat|openai-responses|anthropic (enum enforced).
	Protocol string `toml:"protocol"`
	// BaseURL is the API endpoint; empty = preset default.
	BaseURL string `toml:"base_url"`
	// APIKeyRef is "dpapi:<blob-id>" or "env:NAME" (D36 rule 5). Plaintext
	// keys are structurally impossible here. Empty = keyless provider
	// (e.g. local ollama).
	APIKeyRef string `toml:"api_key_ref"`
	// Billing is pay-per-token|plan (enum enforced). Plan mode tracks
	// against PlanCreditTotalMicro (manually entered; the C23 balance is
	// always labelled an estimate).
	Billing              string `toml:"billing" default:"pay-per-token"`
	PlanCreditTotalMicro int64  `toml:"plan_credit_total_micro"`
	// Compat switches; see Compat.
	Compat Compat `toml:"compat"`
	// RPM / TPM are local token-bucket limits (self-throttling before the
	// provider's 429 does); 0 = unlimited.
	RPM int `toml:"rpm"`
	TPM int `toml:"tpm"`
	// Models is the model catalog of this provider, keyed by model id.
	Models map[string]ModelSpec `toml:"models"`
}

// LLMSection is [llm] (catalog schema v2, SPEC-03 sec 3.1); entirely hot-tier
// (api_key_ref changes trigger re-resolution of the reference, never a
// restart).
type LLMSection struct {
	// TextChain is the ordered text fallback chain, elements
	// "provider/model". 401/403/exhausted-quota/5xx-after-retries all skip
	// to the next element (sec 3.1). Elements must reference existing
	// catalog pairs (validation names the failing element).
	TextChain []string `toml:"text_chain"`
	// Roles assigns models to the fixed agent roles.
	Roles Roles `toml:"roles"`
	// TimeoutMS is the per-request LLM timeout.
	TimeoutMS int `toml:"timeout_ms" default:"60000"`
	// Retry is the per-request retry policy.
	Retry RetryConfig `toml:"retry"`
	// Providers is the provider registry keyed by provider name.
	Providers map[string]Provider `toml:"providers"`
}

// LoopGuard is agent.loop_guard (C22 gradient braking).
type LoopGuard struct {
	// RepeatThresholds is the ascending escalation ladder (default 3,5,8).
	RepeatThresholds []int `toml:"repeat_thresholds" default:"3,5,8"`
}

// AgentSection is [agent]; hot-tier.
type AgentSection struct {
	// MaxRounds caps agent loop iterations (50 fallback, D36).
	MaxRounds int `toml:"max_rounds" default:"50"`
	// TokenBudget caps tokens per task; 0 = uncapped (not recommended).
	TokenBudget int `toml:"token_budget" default:"200000"`
	// PerToolTimeoutMS overrides per-tool timeouts; 0 = each tool's own
	// default applies.
	PerToolTimeoutMS int `toml:"per_tool_timeout_ms"`
	// LoopGuard is the C22 repeat detector ladder.
	LoopGuard LoopGuard `toml:"loop_guard"`
	// SteeringEnabled allows mid-run user steering.
	SteeringEnabled bool `toml:"steering_enabled" default:"true"`
}

// RiskSection is the locked [risk] section (D36 rule 1: any loosening needs L2 re-confirm).
type RiskSection struct {
	// ConfirmTimeoutSec is how long an L2 confirmation card stays open.
	ConfirmTimeoutSec int `toml:"confirm_timeout_sec" default:"300"`
	// L1WindowSec is the auto-approve window for L1 actions; raising it
	// auto-approves more, so an increase is a loosening change.
	L1WindowSec int `toml:"l1_window_sec" default:"2"`
	// ShellEnabled turns on the shell tool at all.
	ShellEnabled bool `toml:"shell_enabled" default:"false"`
	// AllowShellString permits free-string shell commands (vs. argv arrays).
	AllowShellString bool `toml:"allow_shell_string" default:"false"`
	// ShellAllowlist is the command allowlist for shell tools.
	ShellAllowlist []string `toml:"shell_allowlist"`
	// BlacklistOverrides are B-tier blacklist exemptions (per entry).
	BlacklistOverrides []string `toml:"blacklist_overrides"`
	// PermissionMode is the user-facing permission mode (ticket 90, R20/M1):
	// "ask_every_step" (default, strictest) | "ask_high_risk" | "auto_approve".
	// It is the ONLY persisted permission preference in the system: R20/M3
	// makes a manually chosen mode survive a new session and a restart, and
	// nothing else may borrow that channel - a D45 session grant stays
	// session-scoped and dies with the process (PLAN.md:1640, ticket 49), so
	// the two are covered by separate tests, never one "persistence" test.
	// An unknown value is a load error (validateRisk), not a fallback.
	PermissionMode string `toml:"permission_mode" default:"ask_every_step"`
}

// FSSection is the locked [fs] section.
type FSSection struct {
	// AllowedDirs are the readable/writable roots. Additions are loosening.
	AllowedDirs []string `toml:"allowed_dirs"`
	// ReparsePointExceptions whitelists specific symlink/junction paths
	// that would otherwise be refused (per path, not per pattern).
	ReparsePointExceptions []string `toml:"reparse_point_exceptions"`
	// DeleteEnabled turns delete-capable tools on.
	DeleteEnabled bool `toml:"delete_enabled" default:"false"`
}

// ProxyConfig is net.proxy.
type ProxyConfig struct {
	// Mode is system|none|manual (enum enforced). Moving away from "none"
	// opens a new egress path = loosening; moving to "none" = tightening.
	Mode string `toml:"mode" default:"system"`
	// URL is the manual proxy URL (used when mode == "manual").
	URL string `toml:"url"`
}

// NetSection is the locked [net] section.
type NetSection struct {
	// Allowlist is the egress host allowlist; additions are loosening.
	Allowlist []string `toml:"allowlist"`
	// BlockPrivateRanges blocks RFC1918/link-local targets; disabling it
	// (true->false) exposes the intranet = loosening.
	BlockPrivateRanges bool `toml:"block_private_ranges" default:"true"`
	// Proxy routing; see ProxyConfig for the direction rules.
	Proxy ProxyConfig `toml:"proxy"`
}

// PrivacySection is [privacy]. keep_transcript/keep_audio are HARD-CODED
// false (D16/D36): read-only display, writing true is a load error, and no
// code path may implement keep_audio=true (SPEC-03 sec 7).
type PrivacySection struct {
	// RedactPaths masks filesystem paths in user-visible text.
	RedactPaths bool `toml:"redact_paths" default:"false"`
	// DiagnosticsOptIn consents to diagnostic bundle contents.
	DiagnosticsOptIn bool `toml:"diagnostics_opt_in" default:"false"`
	// RetentionDays is the L3/forensics retention window (D35: 30 days).
	RetentionDays int `toml:"retention_days" default:"30"`
	// KeepTranscript is hard-coded false; true in the file = load error.
	KeepTranscript bool `toml:"keep_transcript" default:"false"`
	// KeepAudio is hard-coded false; true in the file = load error.
	KeepAudio bool `toml:"keep_audio" default:"false"`
}

// MemorySection is [memory]; hot-tier.
type MemorySection struct {
	L1Enabled       bool   `toml:"l1_enabled" default:"true"`
	L1Max           int    `toml:"l1_max" default:"20"`
	L3RetentionDays int    `toml:"l3_retention_days" default:"30"`
	ExtractModel    string `toml:"extract_model"`
}

// PanelSection is [panel]; hot-tier.
type PanelSection struct {
	Enabled bool `toml:"enabled" default:"true"`
	// Width is the panel width in px.
	Width int `toml:"width" default:"640"`
	// Height in px; 0 = auto from content.
	Height int `toml:"height"`
	// KeepAliveInSession keeps the WebView2 warm during a session.
	KeepAliveInSession bool `toml:"keep_alive_in_session" default:"true"`
	// Scale is the UI zoom factor; 0 = follow system DPI.
	Scale float64 `toml:"scale"`
}

// CostSection is [cost]. Budgets are micro-USD; 0 = no budget. Reaching 100%
// pauses new tasks by default (C23 hook, ticket 44).
type CostSection struct {
	DailyBudget    int64   `toml:"daily_budget"`
	MonthlyBudget  int64   `toml:"monthly_budget"`
	AlertThreshold float64 `toml:"alert_threshold" default:"0.8"`
	// OverBudget is pause|warn (enum enforced).
	OverBudget string `toml:"over_budget" default:"pause"`
}

// PluginEntry is plugins.<id> - one Tier-2 plugin's grants (D36: a locked change
// reloads the plugin and re-confirms its capabilities; loosening edits gate
// through the L2 hook like every locked section).
type PluginEntry struct {
	// Enabled loads the plugin at all.
	Enabled bool `toml:"enabled"`
	// Capabilities are the granted capability names.
	Capabilities []string `toml:"capabilities"`
	// NetAllowlist is the plugin's private egress allowlist.
	NetAllowlist []string `toml:"net_allowlist"`
	// HostAPI names the C24 GojaHostAPI surfaces granted to the plugin.
	HostAPI []string `toml:"host_api"`
}

// PluginsSection is the locked [plugins] section: the fixed tier2_enabled key plus dynamic
// per-plugin sub-tables ([plugins.<id>]). Go-toml cannot capture dynamic
// sub-tables into a struct field, so Entries is populated/marshaled manually
// by parse.go (single code path, covered by round-trip tests); the Config
// field carries `toml:"-"` so the plain struct marshal never touches it, and
// Entries itself stays untagged: the diff walk treats it as a transparent
// map hanging directly under "plugins".
type PluginsSection struct {
	Tier2Enabled bool `toml:"tier2_enabled" default:"false"`
	// Entries is keyed by plugin id.
	Entries map[string]PluginEntry
}

// ModelsSection is [models]; hot-tier (threshold/switch class).
type ModelsSection struct {
	// Dir is the local model store directory; empty = default under the
	// data dir.
	Dir string `toml:"dir"`
	// Mirror lists model-mirror base URLs (dev default comes from the env
	// layout, SPEC-03 sec 5.2, not from this struct).
	Mirror []string `toml:"mirror"`
	// VerifySignature is HARD-CODED true (C29): model signature checks are
	// not disableable; `false` in the file is a load error, never a setting.
	VerifySignature bool `toml:"verify_signature" default:"true"`
	// LocalOverride maps model-id -> local file path (sideloading).
	LocalOverride map[string]string `toml:"local_override"`
}

// RollConfig is observe.roll.
type RollConfig struct {
	SizeMB int `toml:"size_mb" default:"10"`
	Days   int `toml:"days" default:"7"`
}

// ObserveSection is [observe]; hot-tier.
type ObserveSection struct {
	// Level is debug|info|warn|error (enum enforced).
	Level string `toml:"level" default:"info"`
	// Roll is log rotation.
	Roll RollConfig `toml:"roll"`
	// SLOSampleIntervalSec is the watchdog SLO sampling cadence (D32).
	SLOSampleIntervalSec int `toml:"slo_sample_interval_sec" default:"5"`
}
