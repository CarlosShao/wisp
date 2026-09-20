package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Capability emulation (/__control/capability, ticket 11 AC#6) tests.
//
// The point of the mode is to let the probe suite prove that a probe MEASURES:
// a provider that claims fc/vision but is set to fc=broken / vision=broken must
// fail the real request, not merely report a different status field.

// capsFromState reads the "capabilities" object of /__control/state.
func capsFromState(t *testing.T, ts *httptest.Server) capabilityMode {
	t.Helper()
	resp, err := http.Get(ts.URL + "/__control/state")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var st struct {
		Capabilities capabilityMode `json:"capabilities"`
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &st); err != nil {
		t.Fatalf("decode state %s: %v", b, err)
	}
	return st.Capabilities
}

// postCapabilityRaw POSTs a capability body and returns the status plus body,
// without asserting 200 (Control does that, so the 400 cases need this).
func postCapabilityRaw(t *testing.T, ts *httptest.Server, body string) (int, string) {
	t.Helper()
	resp := postJSON(t, ts.URL+"/__control/capability", body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, string(b)
}

// chatForcedTool asks mockllm's chat route for a FORCED tool call and reports
// whether a tool call came back, plus whether the same turn also carried the
// text answer a degraded (non-tool-calling) model would produce.
func chatForcedTool(t *testing.T, ts *httptest.Server) (toolCall, hasText bool) {
	t.Helper()
	resp := postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"tool_choice":"required",
		  "tools":[{"type":"function","function":{"name":"echo","parameters":{"type":"object"}}}],
		  "messages":[{"role":"user","content":"Call the echo tool with text=ping."}]}`)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	if resp.StatusCode != 200 {
		t.Fatalf("chat status = %d: %s", resp.StatusCode, body)
	}
	return strings.Contains(body, `"tool_calls"`), strings.Contains(body, "echo: Call the echo tool")
}

// imageRequest posts one chat request carrying an image part and returns status+body.
func imageRequest(t *testing.T, ts *httptest.Server) (int, string) {
	t.Helper()
	resp := postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"messages":[{"role":"user","content":[
		  {"type":"image_url","image_url":{"url":"data:image/png;base64,xxx"}},
		  {"type":"text","text":"what color"}]}]}`)
	return resp.StatusCode, readAll(t, resp.Body)
}

// TestCapabilityControlIsAtomic is the guard against a half-applied mode:
// a request with a valid fc and an INVALID vision must be rejected as a whole.
// Before the fix, s.caps.FC was written before vision was validated, so the
// 400 left the server in a capability state nobody asked for.
func TestCapabilityControlIsAtomic(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	postJSON(t, ts.URL+"/__control/reset", `{}`)

	// Pre-request state, established through the same endpoint.
	if code, body := postCapabilityRaw(t, ts, `{"fc":"capable","vision":"capable"}`); code != 200 {
		t.Fatalf("setup capability: status %d: %s", code, body)
	}
	if toolCall, text := chatForcedTool(t, ts); !toolCall || !text {
		t.Fatalf("precondition: capable mode must produce a tool call (toolCall=%v text=%v)", toolCall, text)
	}

	// One field good, one field garbage -> the whole request must be refused.
	code, body := postCapabilityRaw(t, ts, `{"fc":"broken","vision":"nonsense"}`)
	if code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body %s)", code, body)
	}
	if !strings.Contains(body, "vision") {
		t.Errorf("400 detail %q does not name the offending field", body)
	}

	// THE assertion: the rejected request changed nothing.
	got := capsFromState(t, ts)
	if got.FC != "capable" || got.Vision != "capable" {
		t.Errorf("capabilities after the rejected request = %+v, want {capable capable} "+
			"(control endpoint must validate every field before mutating either)", got)
	}
	// Behavioral proof, not just the status field: the server still serves fc.
	if toolCall, _ := chatForcedTool(t, ts); !toolCall {
		t.Error("rejected request silently flipped the served behavior too: forced tool call produced no tool_calls")
	}
}

// TestCapabilityControlAcceptsBothFields checks the happy path stays a no-op
// for omitted fields and applies both when both are valid.
func TestCapabilityControlAppliesBothAndKeepsOmitted(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	postJSON(t, ts.URL+"/__control/reset", `{}`)

	if code, body := postCapabilityRaw(t, ts, `{"fc":"broken","vision":"broken"}`); code != 200 {
		t.Fatalf("status %d: %s", code, body)
	}
	if got := capsFromState(t, ts); got.FC != "broken" || got.Vision != "broken" {
		t.Errorf("capabilities = %+v, want both broken", got)
	}
	// Omitted fields are unchanged.
	if code, body := postCapabilityRaw(t, ts, `{"fc":"capable"}`); code != 200 {
		t.Fatalf("status %d: %s", code, body)
	}
	if got := capsFromState(t, ts); got.FC != "capable" || got.Vision != "broken" {
		t.Errorf("capabilities = %+v, want fc=capable (set) vision=broken (untouched)", got)
	}
	// reset clears the mode.
	postJSON(t, ts.URL+"/__control/reset", `{}`)
	if got := capsFromState(t, ts); got.FC != "capable" || got.Vision != "capable" {
		t.Errorf("capabilities after reset = %+v, want the capable defaults", got)
	}
}

// TestCapabilityBrokenChat proves the emulation is on the RESPONSE, not in a
// status field: fc=broken answers a forced tool_choice with text only, and
// vision=broken rejects the image with a 400 exactly like a non-vision model.
func TestCapabilityBrokenChat(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	postJSON(t, ts.URL+"/__control/reset", `{}`)

	// fc=broken: same request that produced tool_calls above now must not.
	if code, body := postCapabilityRaw(t, ts, `{"fc":"broken"}`); code != 200 {
		t.Fatalf("status %d: %s", code, body)
	}
	toolCall, text := chatForcedTool(t, ts)
	if toolCall {
		t.Error("fc=broken still returned a tool_calls block (the probe could never detect this)")
	}
	if !text {
		t.Error("fc=broken dropped the text answer too: a degraded model answers in text")
	}

	// vision=broken.
	if code, body := postCapabilityRaw(t, ts, `{"vision":"broken"}`); code != 200 {
		t.Fatalf("status %d: %s", code, body)
	}
	status, got := imageRequest(t, ts)
	if status != http.StatusBadRequest {
		t.Errorf("image request status = %d, want 400 (body %s)", status, got)
	}
	if !strings.Contains(got, "vision=broken") {
		t.Errorf("image rejection detail = %q, want it to name the capability mode", got)
	}

	// fc=broken must not affect a plain text turn.
	postJSON(t, ts.URL+"/__control/reset", `{}`)
	if code, body := postCapabilityRaw(t, ts, `{"fc":"broken"}`); code != 200 {
		t.Fatalf("status %d: %s", code, body)
	}
	if status, got := imageRequest(t, ts); status != 200 || !strings.Contains(got, "vision-ok") {
		t.Errorf("vision capable + fc broken: image status = %d body %q, want 200/vision-ok", status, got)
	}
}

// TestCapabilityBrokenOnAllThreeDialects proves the emulation is protocol-wide:
// the probe suite runs on every adapter, so an endpoint that ignores the mode
// would make a lying provider look honest for that dialect only.
func TestCapabilityBrokenOnAllThreeDialects(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	postJSON(t, ts.URL+"/__control/reset", `{}`)
	if code, body := postCapabilityRaw(t, ts, `{"fc":"broken","vision":"broken"}`); code != 200 {
		t.Fatalf("status %d: %s", code, body)
	}

	// Anthropic Messages: forced tool_choice -> no tool_use block, text only.
	resp := postJSON(t, ts.URL+"/v1/messages",
		`{"model":"mock-large","max_tokens":100,
		  "tools":[{"name":"echo","input_schema":{"type":"object"}}],
		  "tool_choice":{"type":"any"},
		  "messages":[{"role":"user","content":"Call the echo tool with text=ping."}]}`)
	msgBody := readAll(t, resp.Body)
	if resp.StatusCode != 200 {
		t.Fatalf("messages status = %d: %s", resp.StatusCode, msgBody)
	}
	if strings.Contains(msgBody, "tool_use") {
		t.Errorf("messages fc=broken returned a tool_use block:\n%s", msgBody)
	}
	if !strings.Contains(msgBody, "echo: Call the echo tool") {
		t.Errorf("messages fc=broken lost the text answer:\n%s", msgBody)
	}

	// Anthropic Messages: image block -> 400.
	resp = postJSON(t, ts.URL+"/v1/messages",
		`{"model":"mock-large","max_tokens":100,"messages":[{"role":"user","content":[
		  {"type":"image","source":{"type":"base64","media_type":"image/png","data":"xxx"}},
		  {"type":"text","text":"what color"}]}]}`)
	msgBody = readAll(t, resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("messages image status = %d, want 400 (body %s)", resp.StatusCode, msgBody)
	}

	// Responses: forced function call -> no function_call item, text only.
	resp = postJSON(t, ts.URL+"/v1/responses",
		`{"model":"mock-large","stream":true,
		  "tools":[{"type":"function","name":"echo","parameters":{"type":"object"}}],
		  "tool_choice":"required",
		  "input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"Call the echo tool with text=ping."}]}]}`)
	rspBody := readAll(t, resp.Body)
	if resp.StatusCode != 200 {
		t.Fatalf("responses status = %d: %s", resp.StatusCode, rspBody)
	}
	if strings.Contains(rspBody, "function_call") {
		t.Errorf("responses fc=broken returned a function_call item:\n%s", rspBody)
	}
	if !strings.Contains(rspBody, "echo: Call the echo tool") {
		t.Errorf("responses fc=broken lost the text answer:\n%s", rspBody)
	}

	// Responses: input_image -> 400.
	resp = postJSON(t, ts.URL+"/v1/responses",
		`{"model":"mock-large","stream":true,"input":[{"type":"message","role":"user","content":[
		  {"type":"input_image","image_url":"data:image/png;base64,xxx"},
		  {"type":"input_text","text":"what color"}]}]}`)
	rspBody = readAll(t, resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("responses image status = %d, want 400 (body %s)", resp.StatusCode, rspBody)
	}
}

// TestCapabilityModeCannotLeakIntoGoldenReplay verifies (it does not assume)
// the claim in the control.go doc comment: capability modes are synthesis-only
// and must not disturb the byte-pinned golden path, which ticket 09's recorded
// fixtures depend on.
func TestCapabilityModeCannotLeakIntoGoldenReplay(t *testing.T) {
	recorded, err := os.ReadFile(filepath.Join("testdata", "golden", "tool-call.sse"))
	if err != nil {
		t.Fatal(err)
	}
	want := string(recorded)
	i := strings.Index(want, "data: ")
	if i < 0 {
		t.Fatal("fixture has no data lines")
	}
	want = want[i:]

	for _, mode := range []string{`{}`, `{"fc":"broken","vision":"broken"}`} {
		ts := newTestServer(t, "testdata/golden")
		postJSON(t, ts.URL+"/__control/reset", `{}`)
		if code, body := postCapabilityRaw(t, ts, mode); code != 200 {
			t.Fatalf("capability %s: status %d: %s", mode, code, body)
		}
		// The request body is deliberately one the SYNTHESIS path would answer
		// differently under the mode (forced tool choice + image part): only a
		// leak can change these bytes.
		for _, route := range []string{"/v1/chat/completions", "/v1/messages", "/v1/responses"} {
			postJSON(t, ts.URL+"/__control/reset", `{}`) // fresh golden cursor
			if code, body := postCapabilityRaw(t, ts, mode); code != 200 {
				t.Fatalf("capability %s after reset: status %d: %s", mode, code, body)
			}
			resp := postJSON(t, ts.URL+route, goldenSelector(route))
			served := readAll(t, resp.Body)
			if resp.StatusCode != 200 {
				t.Fatalf("%s golden replay status = %d: %s", route, resp.StatusCode, served)
			}
			if served != want {
				t.Errorf("capability mode %s leaked into golden replay on %s:\n--- served ---\n%q\n--- fixture ---\n%q",
					mode, route, served, want)
			}
			if !strings.Contains(served, "tool_calls") {
				t.Errorf("golden tool-call replay on %s lost its tool_calls block (mode %s)", route, mode)
			}
		}
	}
}

// goldenSelector is a request body that selects golden/tool-call on each
// dialect while carrying exactly the content the capability modes act on
// (forced tool choice + image input).
func goldenSelector(route string) string {
	switch route {
	case "/v1/messages":
		return `{"model":"golden/tool-call","stream":true,"max_tokens":100,
		  "tools":[{"name":"echo","input_schema":{"type":"object"}}],
		  "tool_choice":{"type":"any"},
		  "messages":[{"role":"user","content":[
		    {"type":"image","source":{"type":"base64","media_type":"image/png","data":"xxx"}},
		    {"type":"text","text":"Call the echo tool."}]}]}`
	case "/v1/responses":
		return `{"model":"golden/tool-call","stream":true,
		  "tools":[{"type":"function","name":"echo","parameters":{"type":"object"}}],
		  "tool_choice":"required",
		  "input":[{"type":"message","role":"user","content":[
		    {"type":"input_image","image_url":"data:image/png;base64,xxx"},
		    {"type":"input_text","text":"Call the echo tool."}]}]}`
	default:
		return `{"model":"golden/tool-call","stream":true,"tool_choice":"required",
		  "tools":[{"type":"function","function":{"name":"echo","parameters":{"type":"object"}}}],
		  "messages":[{"role":"user","content":[
		    {"type":"image_url","image_url":{"url":"data:image/png;base64,xxx"}},
		    {"type":"text","text":"Call the echo tool."}]}]}`
	}
}
