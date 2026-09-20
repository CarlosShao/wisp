package llm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Probe primitive (2026-09-19 supplement): per-capability test requests on
// C5. Ticket 09 defines the CONTRACT (the four capability cases and the
// generic runner); ticket 11 implements the provider-specific probe
// transport where the normalized request needs protocol detail, and writes
// the results into provider_health (SPEC-02 schema v2; the DAO primitive
// lands with this ticket in internal/memory).
//
// Capabilities: fc (function calling), vision, thinking, audio.
//
// DEFERRED(audio, ticket 09): C7 (SPEC-05 sec 3.3) has no audio content
// part, so a normalized audio probe request is not expressible yet. The
// audio case is registered with its response check and an explicit
// NotImplementable note; the request definition lands together with the
// voice-cascade tickets (15/60/61) which own the C7 audio part decision.
// RunProbe(audio) reports OK=false with that reason - it never fakes a pass.

// ProbeCapability names one probeable capability.
type ProbeCapability string

const (
	ProbeFC       ProbeCapability = "fc"
	ProbeVision   ProbeCapability = "vision"
	ProbeThinking ProbeCapability = "thinking"
	ProbeAudio    ProbeCapability = "audio"
)

// AllProbeCapabilities lists the four capabilities in contract order.
func AllProbeCapabilities() []ProbeCapability {
	return []ProbeCapability{ProbeFC, ProbeVision, ProbeThinking, ProbeAudio}
}

// ProbeCase is one capability's minimal request/response check.
type ProbeCase struct {
	Capability  ProbeCapability
	Description string
	// BuildRequest constructs the minimal request. May return nil ONLY when
	// NotImplementable is set (see the audio note above).
	BuildRequest func() (*Request, error)
	// Validate checks the assembled turn after the probe stream.
	Validate func(res TurnResult) error
	// NotImplementable explains why the request cannot be built yet.
	NotImplementable string
}

// tinyPNG is a 1x1 transparent PNG used by the vision probe.
var tinyPNG = mustBase64Decode(
	"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==")

func mustBase64Decode(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("llm: probe fixture broken: %v", err))
	}
	return b
}

func schema(s string) json.RawMessage { return json.RawMessage(s) }

// ProbeCases returns the four capability definitions (contract order).
func ProbeCases() []ProbeCase {
	return []ProbeCase{
		{
			Capability:  ProbeFC,
			Description: "minimal function calling: one echo tool, forced selection, expect a ToolCall",
			BuildRequest: func() (*Request, error) {
				temp := 0.0
				return &Request{
					Model:       "probe",
					Temperature: &temp,
					Messages: []Message{{
						Role:    RoleUser,
						Content: []Content{TextPart{Text: "Call the echo tool with text=ping."}},
					}},
					Tools: []ToolDef{{
						Name:        "echo",
						Description: "echo the given text back",
						Parameters:  schema(`{"type":"object","properties":{"text":{"type":"string"}},"required":["text"]}`),
					}},
					ToolChoice: &ToolChoice{Mode: ToolChoiceRequired},
				}, nil
			},
			Validate: func(res TurnResult) error {
				for _, tc := range res.ToolCalls {
					if tc.Name == "echo" {
						return nil
					}
				}
				return fmt.Errorf("fc probe: no echo tool call in response (stop=%s err=%v)", res.Stop, res.Err)
			},
		},
		{
			Capability:  ProbeVision,
			Description: "minimal vision: 1x1 PNG part, expect non-empty text answer",
			BuildRequest: func() (*Request, error) {
				return &Request{
					Model: "probe",
					Messages: []Message{{
						Role: RoleUser,
						Content: []Content{
							ImagePart{MimeType: "image/png", BytesRef: "probe:vision-1x1", Alt: "1x1 probe image"},
							TextPart{Text: "What color is this image? Answer in one short sentence."},
						},
					}},
					// The probe supplies ResolveBytes directly; production
					// requests get it from the caller that owns the bytes.
					ResolveBytes: func(ref string) (string, []byte, error) {
						return "image/png", tinyPNG, nil
					},
				}, nil
			},
			Validate: func(res TurnResult) error {
				if res.Err != nil {
					return fmt.Errorf("vision probe failed: %v", res.Err)
				}
				if res.Text == "" {
					return fmt.Errorf("vision probe: empty text answer (stop=%s)", res.Stop)
				}
				return nil
			},
		},
		{
			Capability:  ProbeThinking,
			Description: "minimal thinking: intensity=low, expect reasoning delta or a text answer",
			BuildRequest: func() (*Request, error) {
				return &Request{
					Model:             "probe",
					ThinkingIntensity: "low",
					Messages: []Message{{
						Role:    RoleUser,
						Content: []Content{TextPart{Text: "2+2=? Answer with just the number."}},
					}},
				}, nil
			},
			Validate: func(res TurnResult) error {
				if res.Err != nil {
					return fmt.Errorf("thinking probe failed: %v", res.Err)
				}
				if res.Reasoning == "" && res.Text == "" {
					return fmt.Errorf("thinking probe: no reasoning and no text (stop=%s)", res.Stop)
				}
				return nil
			},
		},
		{
			Capability:  ProbeAudio,
			Description: "audio_in round-trip: minimal input_audio part (request deferred, check defined)",
			BuildRequest: func() (*Request, error) {
				return nil, nil // see NotImplementable
			},
			Validate: func(res TurnResult) error {
				if res.Err != nil {
					return fmt.Errorf("audio probe failed: %v", res.Err)
				}
				if res.Text == "" {
					return fmt.Errorf("audio probe: empty transcription answer (stop=%s)", res.Stop)
				}
				return nil
			},
			NotImplementable: "C7 (SPEC-05 sec 3.3) has no audio content part yet; the audio probe request lands with the voice-cascade tickets (15/60/61) that own that decision. The response check is defined here so ticket 11 only adds the request.",
		},
	}
}

// ProbeCaseFor returns one capability's case.
func ProbeCaseFor(cap ProbeCapability) (ProbeCase, bool) {
	for _, c := range ProbeCases() {
		if c.Capability == cap {
			return c, true
		}
	}
	return ProbeCase{}, false
}

// ProbeResult is one probe outcome. Persistence into provider_health is the
// caller's job (ticket 11; DAO primitive in internal/memory).
type ProbeResult struct {
	Capability ProbeCapability `json:"capability"`
	OK         bool            `json:"ok"`
	LatencyMS  int64           `json:"latency_ms"`
	Detail     string          `json:"detail,omitempty"`
	At         time.Time       `json:"at"` // wall clock, for the persisted record only
}

// RunProbe executes one capability probe against p and classifies the
// outcome. A probe that cannot even be attempted (audio today) returns
// OK=false with the reason in Detail - never a fake pass. Callers substitute
// the real model id into the built request before probing (RunProbeSuite does
// that for them).
func RunProbe(ctx context.Context, p LlmProvider, cap ProbeCapability) ProbeResult {
	cse, ok := ProbeCaseFor(cap)
	if !ok {
		return ProbeResult{Capability: cap, At: observe.NowWallUTC(), Detail: "unknown capability"}
	}
	return runProbeOnModel(ctx, p, cse, "")
}

// runProbeOnModel executes one probe case; when model is non-empty it replaces
// the case's placeholder model id, so the request names the model that is
// actually being measured.
func runProbeOnModel(ctx context.Context, p LlmProvider, cse ProbeCase, model string) ProbeResult {
	res := ProbeResult{Capability: cse.Capability, At: observe.NowWallUTC()}
	req, err := cse.BuildRequest()
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if req == nil {
		res.Detail = "not implementable: " + cse.NotImplementable
		return res
	}
	if model != "" {
		req.Model = model
	}

	start := time.Now()
	collector := NewTurnCollector()
	streamErr := p.Stream(ctx, req, func(ev StreamEvent) error {
		collector.Observe(ev)
		return nil
	})
	res.LatencyMS = time.Since(start).Milliseconds()
	turn := collector.Result()

	if streamErr != nil {
		res.Detail = streamErr.Error()
		return res
	}
	if err := cse.Validate(turn); err != nil {
		res.Detail = err.Error()
		return res
	}
	res.OK = true
	return res
}
