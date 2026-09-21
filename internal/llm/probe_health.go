package llm

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 11 AC#6, consumer side: turn probe.go's per-capability PRIMITIVE
// (ticket 09's contract: one minimal real request per capability, response
// checked) into the provider_health WRITE PATH (SPEC-02 schema v2) plus the
// declared-vs-measured mismatch event (SPEC-05 sec 3.1: 「声明 PASS 实测 FAIL」
// 必须可见). SPEC-05 spells that pair as a check mark and a cross; D22 ban #8
// keeps those glyphs out of internal/, so this quotation is ASCII by design.
//
// The load-bearing property is that a verdict is MEASURED: RunProbeSuite forms a
// verdict only from RunProbe's outcome, i.e. from a request that really went out
// and an answer that really came back. The declared bits are read for exactly
// one purpose - naming the mismatch - and never feed a verdict. A probe whose
// answer could be derived from config would pass its own tests while telling the
// user nothing true.
//
// Silence is structural rather than documented: a sink-less or notice-less run
// is a configuration error, because a probe nobody can read is not AC#6.

// MeasuredFlags is the MEASURED capability bit-set of one provider/model pair,
// i.e. the payload of provider_health.probe_json. A nil flag means "never
// probed", which is a DIFFERENT fact from measured-false; that distinction is
// why the audio bit stays nil until a C7 audio part exists (ticket 61).
//
// It mirrors memory.ProbeFlags field for field. internal/llm deliberately does
// not import the storage layer (C5 must not know about SQLite), so the
// composition root adapts the two - memoryHealthSink in probe_health_test.go is
// the exact mapping the wiring code performs.
type MeasuredFlags struct {
	Text     *bool
	Vision   *bool
	AudioIn  *bool
	AudioOut *bool
	Thinking *bool
	FC       *bool
}

// Set returns a copy of f with one capability's measured verdict recorded.
func (f MeasuredFlags) Set(cap string, v bool) MeasuredFlags {
	switch cap {
	case "text":
		f.Text = &v
	case "vision":
		f.Vision = &v
	case "audio_in":
		f.AudioIn = &v
	case "audio_out":
		f.AudioOut = &v
	case "thinking":
		f.Thinking = &v
	case "fc":
		f.FC = &v
	}
	return f
}

// Get returns the measured verdict for a capability (nil = never probed).
func (f MeasuredFlags) Get(cap string) *bool {
	switch cap {
	case "text":
		return f.Text
	case "vision":
		return f.Vision
	case "audio_in":
		return f.AudioIn
	case "audio_out":
		return f.AudioOut
	case "thinking":
		return f.Thinking
	case "fc":
		return f.FC
	}
	return nil
}

// capabilityFlagKey maps a probe capability to its provider_health key. The
// audio probe tests input audio, so it maps to audio_in.
func capabilityFlagKey(cap ProbeCapability) string {
	if cap == ProbeAudio {
		return "audio_in"
	}
	return string(cap)
}

// declaredCapability reads one bit out of the catalog declaration. It exists
// ONLY to name mismatches; it never feeds a verdict.
func declaredCapability(c config.Capabilities, cap ProbeCapability) bool {
	switch cap {
	case ProbeFC:
		return c.FC
	case ProbeVision:
		return c.Vision
	case ProbeThinking:
		return c.Thinking
	case ProbeAudio:
		return c.AudioIn
	}
	return false
}

// HealthSink is the provider_health write seam. Its signature mirrors
// memory.Store.UpsertProviderProbe so the composition adapter is a field copy.
// probe_json is REPLACED wholesale by the DAO, so a run writes exactly once per
// pair with every verdict it holds - never once per capability, which would
// erase the earlier ones.
type HealthSink interface {
	RecordProbe(ctx context.Context, provider, model string, flags MeasuredFlags,
		ok bool, latencyMS int64, lastProbeAt time.Time) error
}

// ProbeMismatch is the structured 「声明 PASS / 实测 FAIL」fact: the catalog claims a
// capability the wire refused to demonstrate. The seam emits DATA; the ball and
// the panel (ticket 40) own rendering it, but Label() is the single frozen copy
// so no surface invents its own wording (the RetryNotice rule).
type ProbeMismatch struct {
	Provider   string
	Model      string
	Capability ProbeCapability
	Declared   bool
	Measured   bool
	Detail     string    // the probe's own failure detail (adapter-classified already)
	At         time.Time // wall clock: a record, never a deadline (D42#9)
}

// Kind is the stable discriminator a state event carries, so a surface can route
// on it without string-matching the copy.
func (ProbeMismatch) Kind() string { return "capability_mismatch" }

// Label is the user-visible copy for one mismatch.
func (m ProbeMismatch) Label() string {
	return fmt.Sprintf("能力实测不符：%s/%s 声明支持 %s，实测不可用", m.Provider, m.Model, m.Capability)
}

// ProbeOutcome is one capability's verdict next to the declaration it confirms or
// contradicts.
type ProbeOutcome struct {
	Capability ProbeCapability
	Declared   bool
	Result     ProbeResult
	// Attempted is false when no request could be sent at all (audio today);
	// then no verdict is recorded rather than a fabricated one.
	Attempted bool
}

// Mismatch reports the AC#6 fact: declared usable, measured unusable.
func (o ProbeOutcome) Mismatch() bool { return o.Attempted && o.Declared && !o.Result.OK }

// ProbeReport is one probe run for a provider/model pair.
type ProbeReport struct {
	Provider   string
	Model      string
	Outcomes   []ProbeOutcome
	Mismatches []ProbeMismatch
	// Flags carries only the capabilities that were actually attempted.
	Flags MeasuredFlags
	// LatencyMS is the median of the attempted probes (provider_health
	// latency_ms_p50).
	LatencyMS int64
	// AllOK is what the sink received as the run-level verdict: every attempted
	// capability measured usable.
	AllOK   bool
	Written bool
	// Probed lists the capabilities attempted, in order.
	Probed []string
}

// Verdict returns one capability's outcome.
func (r ProbeReport) Verdict(cap ProbeCapability) (ProbeOutcome, bool) {
	for _, o := range r.Outcomes {
		if o.Capability == cap {
			return o, true
		}
	}
	return ProbeOutcome{}, false
}

// DefaultProbeCapabilities is what a run attempts: the capabilities with a
// request definition today. Audio is left out because its request is deferred
// (C7 has no audio part); attempting it would only record "not implementable",
// which is not a verdict.
var DefaultProbeCapabilities = []ProbeCapability{ProbeFC, ProbeVision, ProbeThinking}

// ProbeSuiteOptions configures one run.
type ProbeSuiteOptions struct {
	// Provider and Model are the provider_health primary key AND the model id
	// substituted into every probe request (probe.go's contract: the caller
	// supplies the real model).
	Provider string
	Model    string
	// Declared is the catalog's claim, used to name mismatches only.
	Declared config.Capabilities
	// Capabilities defaults to DefaultProbeCapabilities.
	Capabilities []ProbeCapability
	// Sink writes provider_health, Mismatch announces every 「声明 PASS 实测 FAIL」.
	// Both are required: a run that cannot record, or cannot be loud about a
	// contradiction, fails instead of quietly doing nothing.
	Sink     HealthSink
	Mismatch func(ProbeMismatch)
	// Now is the wall-clock source for the persisted record only.
	Now func() time.Time
}

// RunProbeSuite measures what a provider/model pair can actually do and records
// the verdicts. It returns the report even when the write fails, so a caller can
// still surface what was measured.
func RunProbeSuite(ctx context.Context, p LlmProvider, o ProbeSuiteOptions) (ProbeReport, error) {
	if p == nil {
		return ProbeReport{}, errors.New("llm: probe suite needs a provider to measure")
	}
	if o.Provider == "" || o.Model == "" {
		return ProbeReport{}, errors.New("llm: probe suite needs the provider/model key it is measuring")
	}
	if o.Sink == nil {
		return ProbeReport{}, errors.New("llm: probe suite requires a HealthSink (an unread probe is not a measurement)")
	}
	if o.Mismatch == nil {
		return ProbeReport{}, errors.New("llm: probe suite requires a mismatch notice (SPEC-05 sec 3.1: declared-vs-measured must be visible)")
	}
	caps := o.Capabilities
	if len(caps) == 0 {
		caps = DefaultProbeCapabilities
	}
	at := observe.NowWallUTC()
	if o.Now != nil {
		at = o.Now()
	}

	rep := ProbeReport{Provider: o.Provider, Model: o.Model, AllOK: true}
	var latencies []int64
	for _, c := range caps {
		cse, ok := ProbeCaseFor(c)
		if !ok {
			return rep, fmt.Errorf("llm: unknown probe capability %q", c)
		}
		res := runProbeOnModel(ctx, p, cse, o.Model)
		out := ProbeOutcome{
			Capability: c, Declared: declaredCapability(o.Declared, c),
			Result: res, Attempted: cse.NotImplementable == "",
		}
		rep.Outcomes = append(rep.Outcomes, out)
		rep.Probed = append(rep.Probed, string(c))
		if !out.Attempted {
			// No request went out, so nothing is recorded: the flag stays nil
			// (never probed) and the run verdict is not judged on it.
			continue
		}
		latencies = append(latencies, res.LatencyMS)
		rep.Flags = rep.Flags.Set(capabilityFlagKey(c), res.OK)
		if res.OK {
			continue
		}
		rep.AllOK = false
		if out.Mismatch() {
			m := ProbeMismatch{
				Provider: o.Provider, Model: o.Model, Capability: c,
				Declared: true, Measured: false, Detail: res.Detail, At: at,
			}
			rep.Mismatches = append(rep.Mismatches, m)
			o.Mismatch(m)
		}
	}
	rep.LatencyMS = medianMS(latencies)

	if err := o.Sink.RecordProbe(ctx, o.Provider, o.Model, rep.Flags, rep.AllOK,
		rep.LatencyMS, at); err != nil {
		return rep, fmt.Errorf("llm: record probe for %s/%s: %w", o.Provider, o.Model, err)
	}
	rep.Written = true
	return rep, nil
}

// medianMS is provider_health's latency_ms_p50 over one probe run.
func medianMS(v []int64) int64 {
	if len(v) == 0 {
		return 0
	}
	s := append([]int64(nil), v...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	return s[len(s)/2]
}
