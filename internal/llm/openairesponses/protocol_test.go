package openairesponses

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
)

// Protocol-shape proofs for the Responses adapter, read back from the
// request body the live mockllm actually received (AC#2's negative half: an
// implicit-cache provider must not sprout cache_control markers).

func baseTurn() *llm.Request {
	return &llm.Request{
		Model:  "mock-small",
		System: []llm.Content{llm.TextPart{Text: "identity"}, llm.TextPart{Text: "style"}},
		Messages: []llm.Message{{
			Role:    llm.RoleUser,
			Content: []llm.Content{llm.TextPart{Text: "hi"}},
		}},
		Tools: []llm.ToolDef{{
			Name: "get_weather", Description: "d",
			Parameters: json.RawMessage(`{"type":"object"}`),
		}},
		// Declared even though this protocol cannot honor it: the seam's
		// capability bit says breakpoints have no effect here.
		CacheBreakpoints: []int{-1, 0},
	}
}

func TestRequestBodyIsResponsesDialect(t *testing.T) {
	proc := adaptertest.StartMockllm(t)
	proc.Reset(t)
	proc.StreamGolden(t, unit(), "responses-tool-call", baseTurn())
	body, hdr := proc.LastRequest(t, "responses")

	var sent map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &sent); err != nil {
		t.Fatalf("captured body: %v\n%s", err, body)
	}
	if _, ok := sent["messages"]; ok {
		t.Error("a Responses request must not carry the chat `messages` field")
	}
	if _, ok := sent["input"]; !ok {
		t.Fatal("missing the `input` item array")
	}
	if inst := string(sent["instructions"]); !strings.Contains(inst, "identity") ||
		!strings.Contains(inst, "style") {
		t.Errorf("System did not flatten into instructions: %s", inst)
	}
	if strings.Contains(body, "cache_control") {
		t.Error("cache_control appeared on an implicit-cache protocol (C5 capability says breakpoints have no effect here)")
	}
	if strings.Contains(body, "max_tokens") {
		t.Errorf("the chat field max_tokens leaked into a Responses body: %s", body)
	}
	if auth := hdr.Get("Authorization"); !strings.HasPrefix(auth, "Bearer ") {
		t.Errorf("Authorization = %q, want the Bearer scheme", auth)
	}

	// Tools use the flat Responses shape (no {function:{...}} nesting).
	var tools []map[string]any
	if err := json.Unmarshal(sent["tools"], &tools); err != nil || len(tools) != 1 {
		t.Fatalf("tools field: %v %s", err, sent["tools"])
	}
	if tools[0]["type"] != "function" || tools[0]["name"] != "get_weather" {
		t.Errorf("tool declaration = %v, want flat {type:function,name:...}", tools[0])
	}
	if _, nested := tools[0]["function"]; nested {
		t.Errorf("chat-style nested function object in a Responses tool: %v", tools[0])
	}
	proc.Reset(t)
}

// TestToolTurnsRoundTripThroughTheWire covers the function_call /
// function_call_output item shapes the agent loop depends on.
func TestToolTurnsRoundTripThroughTheWire(t *testing.T) {
	proc := adaptertest.StartMockllm(t)
	proc.Reset(t)
	req := baseTurn()
	req.Messages = append(req.Messages,
		llm.Message{Role: llm.RoleAssistant, Content: []llm.Content{
			llm.ToolUsePart{
				ID: "call_a1", Name: "get_weather",
				Input: json.RawMessage(`{"city":"Zhuhai"}`),
			},
		}},
		llm.Message{Role: llm.RoleTool, Content: []llm.Content{
			llm.ToolResultPart{
				ID:      "call_a1",
				Content: []llm.Content{llm.TextPart{Text: "sunny"}},
			},
		}},
	)
	proc.StreamGolden(t, unit(), "responses-tool-call", req)
	body, _ := proc.LastRequest(t, "responses")
	var sent struct {
		Input []struct {
			Type      string          `json:"type"`
			CallID    string          `json:"call_id"`
			Name      string          `json:"name"`
			Arguments string          `json:"arguments"`
			Output    json.RawMessage `json:"output"`
		} `json:"input"`
	}
	if err := json.Unmarshal([]byte(body), &sent); err != nil {
		t.Fatalf("body: %v\n%s", err, body)
	}
	var call, out bool
	for _, it := range sent.Input {
		switch it.Type {
		case "function_call":
			call = it.CallID == "call_a1" && it.Name == "get_weather" &&
				it.Arguments == `{"city":"Zhuhai"}`
		case "function_call_output":
			out = it.CallID == "call_a1" && strings.Contains(string(it.Output), "sunny")
		}
	}
	if !call {
		t.Errorf("assistant tool_use did not encode into a function_call item: %s", body)
	}
	if !out {
		t.Errorf("tool result did not encode into a function_call_output item: %s", body)
	}
	proc.Reset(t)
}

// TestThinkingIntensityMapsToReasoningEffort: the enum reaches the protocol's
// own parameter (the chat adapter deliberately has no such field).
func TestThinkingIntensityMapsToReasoningEffort(t *testing.T) {
	tests := []struct {
		in    string
		want  string
		blank bool
	}{
		{in: "low", want: `"effort":"low"`},
		{in: "medium", want: `"effort":"medium"`},
		{in: "high", want: `"effort":"high"`},
		{in: "off", blank: true},
		{in: "", blank: true},
	}
	for _, tc := range tests {
		proc := adaptertest.StartMockllm(t)
		proc.Reset(t)
		req := baseTurn()
		req.ThinkingIntensity = tc.in
		proc.StreamGolden(t, unit(), "responses-thinking", req)
		body, _ := proc.LastRequest(t, "responses")
		switch {
		case tc.blank && strings.Contains(body, `"reasoning"`) &&
			strings.Contains(body, `"effort"`):
			t.Errorf("thinking_intensity %q must omit the reasoning object: %s", tc.in, body)
		case !tc.blank && !strings.Contains(body, tc.want):
			t.Errorf("thinking_intensity %q did not map to %s: %s", tc.in, tc.want, body)
		}
		proc.Reset(t)
	}
}

// TestStopSequencesAreRejectedNotSilentlyDropped guards the "tolerate per the
// compat switch, or error out" rule (SPEC-03 sec 3.1).
func TestStopSequencesAreRejectedNotSilentlyDropped(t *testing.T) {
	req := baseTurn()
	req.StopSequences = []string{"\nUser:"}
	strict := New(endpointOptionsFor("http://127.0.0.1:1/v1", false, false))
	if _, _, err := adaptertest.Drain(t, strict, req); err == nil ||
		!strings.Contains(err.Error(), "stop-sequence") {
		t.Errorf("strict: err = %v, want an actionable unsupported-field error", err)
	}
	looseSrv, _ := adaptertest.ServeBody(t, "# wisp golden sse v1\n# @response 200\ndata: "+
		`{"type":"response.created","response":{"id":"r","status":"in_progress"}}`+"\n\n"+
		"data: "+`{"type":"response.output_item.added","output_index":0,"item":{"id":"m","type":"message","role":"assistant","content":[]}}`+"\n\n"+
		"data: "+`{"type":"response.output_text.delta","item_id":"m","delta":"ok"}`+"\n\n"+
		"data: "+`{"type":"response.completed","response":{"id":"r","status":"completed","usage":{"input_tokens":3,"output_tokens":1}}}`+"\n\n"+
		"data: [DONE]\n")
	loose := New(endpointOptionsFor(looseSrv, true, false))
	_, turn, err := adaptertest.Drain(t, loose, req)
	if err != nil {
		t.Fatalf("loose: %v", err)
	}
	if turn.Text != "ok" || turn.Stop != llm.StopEndTurn {
		t.Errorf("loose turn = %+v", turn)
	}
}

// TestInfoReportsImplicitNoBreakpoints pins the C5 cache capability of this
// protocol (the agent core must not declare breakpoints for it).
func TestInfoReportsImplicitNoBreakpoints(t *testing.T) {
	info := New(endpointOptionsFor("http://127.0.0.1:1/v1", false, false)).Info()
	if info.Cache.Mode != llm.CacheImplicit || info.Cache.Breakpoints {
		t.Errorf("cache bits = %+v, want implicit / no breakpoints", info.Cache)
	}
	if info.Protocol != Protocol {
		t.Errorf("protocol = %s", info.Protocol)
	}
}
