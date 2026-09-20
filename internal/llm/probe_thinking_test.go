package llm_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
)

// thinkingProbeModel is the model id these two cases measure; the provider name
// is runProbe's default ("mock").
const thinkingProbeModel = "mock-small"

// thinkingRow reads back the persisted provider_health record.
func thinkingRow(t *testing.T, pf *probeFixture) memory.ProviderHealth {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row, err := pf.store.ProviderHealthRow(ctx, "mock", thinkingProbeModel)
	if err != nil {
		t.Fatalf("provider_health for mock/%s: %v", thinkingProbeModel, err)
	}
	return row
}

// Ruling A11's second half: ticket 09's thinking check accepted a plain text
// answer, so it could not detect a broken thinker. The check now requires a
// reasoning delta, and both directions are pinned - here without a network at
// all (the deterministic layer), and against the live mockllm below.

// TestThinkingCheckRejectsPlainTextAnswer is the never-skip layer: the response
// check itself is the assertion, so no skipped integration case can hide a
// regression back to "text is good enough".
func TestThinkingCheckRejectsPlainTextAnswer(t *testing.T) {
	cse, ok := llm.ProbeCaseFor(llm.ProbeThinking)
	if !ok {
		t.Fatal("no thinking probe case")
	}
	if err := cse.Validate(llm.TurnResult{Text: "4", Stop: "end_turn"}); err == nil {
		t.Error("a text-only answer passed the thinking check: a broken thinker would still be reported as working")
	} else if !strings.Contains(err.Error(), "no reasoning delta") {
		t.Errorf("unexpected refusal: %v", err)
	}
	if err := cse.Validate(llm.TurnResult{Reasoning: "2+2 is 4", Text: "4"}); err != nil {
		t.Errorf("a real reasoning delta must pass: %v", err)
	}
	if err := cse.Validate(llm.TurnResult{}); err == nil {
		t.Error("an empty turn must not pass")
	}
}

// TestThinkingProbeRequestCarriesTheHarnessMarker pins why the prompt text has
// the [think] marker: the openai-chat adapter deliberately does not map
// ThinkingIntensity onto the wire, so without the marker a chat-route thinking
// probe is indistinguishable from a plain request and the mock would answer
// with no reasoning.
func TestThinkingProbeRequestCarriesTheHarnessMarker(t *testing.T) {
	cse, _ := llm.ProbeCaseFor(llm.ProbeThinking)
	req, err := cse.BuildRequest()
	if err != nil || req == nil {
		t.Fatalf("build: %v", err)
	}
	if req.ThinkingIntensity != "low" {
		t.Errorf("thinking_intensity = %q, want low", req.ThinkingIntensity)
	}
	var text string
	for _, part := range req.Messages[0].Content {
		if tp, ok := part.(llm.TextPart); ok {
			text += tp.Text
		}
	}
	if !strings.Contains(text, "[think]") {
		t.Errorf("the thinking probe prompt lost its harness marker: %q", text)
	}
}

// TestProbeSuiteMeasuresBrokenThinking is the live half: the SERVER is proven to
// be in thinking=broken mode, the requests are proven to have gone out, and the
// verdict the suite records is the opposite of the declaration.
func TestProbeSuiteMeasuresBrokenThinking(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{"thinking":"broken"}`)

	rep, mismatches := runProbe(t, pf, "openai-chat", llm.ProbeSuiteOptions{
		Model:        thinkingProbeModel,
		Declared:     declaredEverything,
		Capabilities: []llm.ProbeCapability{llm.ProbeThinking},
	})
	if n := pf.requests(t, "chat"); n == 0 {
		t.Fatal("the thinking probe sent no request: it cannot have measured anything")
	}
	if len(rep.Outcomes) != 1 {
		t.Fatalf("outcomes = %+v, want the one thinking verdict", rep.Outcomes)
	}
	o := rep.Outcomes[0]
	if o.Result.OK {
		t.Errorf("thinking measured usable while the server withholds every reasoning delta: %+v", o.Result)
	}
	if !strings.Contains(o.Result.Detail, "no reasoning delta") {
		t.Errorf("detail = %q, want the missing-reasoning refusal", o.Result.Detail)
	}
	if len(mismatches) != 1 || mismatches[0].Capability != llm.ProbeThinking {
		t.Fatalf("mismatches = %+v, want the one declared-vs-measured event", mismatches)
	}
	if !mismatches[0].Declared || mismatches[0].Measured {
		t.Errorf("mismatch = %+v, want declared true / measured false", mismatches[0])
	}
	if got := mustFlag(t, thinkingRow(t, pf).Probe, "thinking"); got != false {
		t.Errorf("persisted thinking flag = %v, want false", got)
	}
}

// TestProbeSuiteThinkingCapableIsNotTheSameAsBroken closes the other direction
// against the same fixture: with the mode off the record, the identical request
// is measured usable. The pair is what makes the check a measurement.
func TestProbeSuiteThinkingCapableIsNotTheSameAsBroken(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{"thinking":"capable"}`)
	rep, mismatches := runProbe(t, pf, "openai-chat", llm.ProbeSuiteOptions{
		Model:        thinkingProbeModel,
		Declared:     declaredEverything,
		Capabilities: []llm.ProbeCapability{llm.ProbeThinking},
	})
	if len(rep.Outcomes) != 1 || !rep.Outcomes[0].Result.OK {
		t.Fatalf("a thinking server must measure usable: %+v", rep.Outcomes)
	}
	if len(mismatches) != 0 {
		t.Errorf("no mismatch may be reported when the measurement agrees: %+v", mismatches)
	}
	if got := mustFlag(t, thinkingRow(t, pf).Probe, "thinking"); got != true {
		t.Errorf("persisted thinking flag = %v, want true", got)
	}
}
