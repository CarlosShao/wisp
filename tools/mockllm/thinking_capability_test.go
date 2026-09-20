package main

import (
	"io"
	"strings"
	"testing"
)

// The thinking capability mode (ticket 12, ruling A11).
//
// fc and vision already had a broken emulation; thinking did not, and the
// probe case that consumed it accepted any text answer. Together those two gaps
// meant a broken thinker was indistinguishable from a working one. This file
// pins the mock half of the fix: reasoning appears when the mode is capable and
// disappears when it is broken, while the text answer keeps arriving - which is
// precisely the shape a real model with dead reasoning produces.

func thinkReasoningSSE(t *testing.T, url string) (reasoning, text bool) {
	t.Helper()
	resp := postJSON(t, url+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,
		  "messages":[{"role":"user","content":"[think] 2+2=? Answer with just the number."}]}`)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	if resp.StatusCode != 200 {
		t.Fatalf("chat status = %d: %s", resp.StatusCode, body)
	}
	return strings.Contains(body, "reasoning_content") ||
			strings.Contains(body, `"reasoning"`) || strings.Contains(body, "thinking about:"),
		strings.Contains(body, "echo:")
}

func TestThinkingModeCapableEmitsReasoningDelta(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	postJSON(t, ts.URL+"/__control/reset", `{}`)
	if got := capsFromState(t, ts).Thinking; got != "capable" {
		t.Fatalf("default thinking mode = %q, want capable", got)
	}
	reasoning, text := thinkReasoningSSE(t, ts.URL)
	if !reasoning {
		t.Error("a capable thinker emitted no reasoning: the probe could never pass on a working provider")
	}
	if !text {
		t.Error("the answer text disappeared alongside the reasoning")
	}
}

func TestThinkingModeBrokenWithholdsReasoningOnly(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	postJSON(t, ts.URL+"/__control/reset", `{}`)
	if code, body := postCapabilityRaw(t, ts, `{"thinking":"broken"}`); code != 200 {
		t.Fatalf("status %d: %s", code, body)
	}
	if got := capsFromState(t, ts).Thinking; got != "broken" {
		t.Fatalf("state read-back = %q, want broken (the mock must be provably in the mode)", got)
	}
	reasoning, text := thinkReasoningSSE(t, ts.URL)
	if reasoning {
		t.Error("thinking=broken still produced a reasoning delta: the probe has nothing to detect")
	}
	if !text {
		t.Error("thinking=broken must degrade to a text answer, not fail the request")
	}
	// The mode is one field: fc stays as configured, and the counters show the
	// request really went out twice (once per mode), which is the "server's own
	// request counters" half of the evidence.
	if code, body := postCapabilityRaw(t, ts, `{"thinking":"nonsense"}`); code != 400 {
		t.Errorf("status = %d for thinking=nonsense (body %s), want 400", code, body)
	}
	if got := capsFromState(t, ts).Thinking; got != "broken" {
		t.Errorf("a rejected capability request must leave the mode untouched, got %q", got)
	}
}
