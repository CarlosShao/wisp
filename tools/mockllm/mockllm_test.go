package main

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// In-module tests of the mockllm handlers (httptest against the mux). The
// cross-module byte-identical guarantee vs internal/llm/golden is pinned by
// the integration test in the main module.

func newTestServer(t *testing.T, goldenDir string) *httptest.Server {
	t.Helper()
	srv := NewServer(goldenDir, []string{"mock-small", "mock-large"})
	mux := http.NewServeMux()
	srv.Register(mux)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func postJSON(t *testing.T, url, body string, hdr ...map[string]string) *http.Response {
	t.Helper()
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	for _, m := range hdr {
		for k, v := range m {
			req.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func readAll(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// SSE chunks of a response body (split on the blank dispatcher line).
func sseChunks(t *testing.T, body string) []string {
	t.Helper()
	var out []string
	var cur strings.Builder
	sc := bufio.NewScanner(strings.NewReader(body))
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteString(line + "\n")
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

func chunkData(t *testing.T, chunk string) string {
	t.Helper()
	for _, line := range strings.Split(chunk, "\n") {
		if v, ok := strings.CutPrefix(line, "data: "); ok {
			return v
		}
	}
	t.Fatalf("no data line in chunk %q", chunk)
	return ""
}

func TestChatSynthStreamEcho(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	resp := postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"messages":[{"role":"user","content":"hello mock"}]}`)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	chunks := sseChunks(t, readAll(t, resp.Body))

	if n := len(chunks); n < 5 {
		t.Fatalf("chunk count = %d, want >= 5 (role+text+finish+usage+[DONE])\n%v", n, chunks)
	}
	if got := chunkData(t, chunks[len(chunks)-1]); got != "[DONE]" {
		t.Errorf("last chunk = %q, want [DONE]", got)
	}

	// Assemble text deltas; verify finish_reason and usage chunks exist.
	text := ""
	sawFinish, sawUsage := "", false
	for _, c := range chunks[:len(chunks)-1] {
		var m struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
			Usage *struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
			} `json:"usage"`
		}
		if err := json.Unmarshal([]byte(chunkData(t, c)), &m); err != nil {
			t.Fatalf("chunk not JSON: %v", err)
		}
		if len(m.Choices) > 0 {
			text += m.Choices[0].Delta.Content
			if m.Choices[0].FinishReason != nil && *m.Choices[0].FinishReason != "" &&
				*m.Choices[0].FinishReason != "null" {
				sawFinish = *m.Choices[0].FinishReason
			}
		}
		if m.Usage != nil {
			sawUsage = true
			if m.Usage.PromptTokens == 0 || m.Usage.CompletionTokens == 0 {
				t.Error("usage chunk with zero tokens")
			}
		}
	}
	if text != "echo: hello mock" {
		t.Errorf("assembled text = %q, want %q", text, "echo: hello mock")
	}
	if sawFinish != "stop" {
		t.Errorf("finish_reason = %q, want stop", sawFinish)
	}
	if !sawUsage {
		t.Error("no usage chunk in stream")
	}
}

func TestChatSynthToolCallVisionAudioAndMaxTokens(t *testing.T) {
	ts := newTestServer(t, t.TempDir())

	// Forced tool call.
	resp := postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"tool_choice":{"type":"function","function":{"name":"get_weather"}},
		  "tools":[{"type":"function","function":{"name":"get_weather","parameters":{}}}],
		  "messages":[{"role":"user","content":"weather?"}]}`)
	var toolSeen bool
	for _, c := range sseChunks(t, readAll(t, resp.Body)) {
		if strings.Contains(c, `"tool_calls"`) && strings.Contains(c, "get_weather") {
			toolSeen = true
		}
	}
	if !toolSeen {
		t.Error("forced tool_choice produced no tool_calls chunk")
	}

	// Vision marker.
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"messages":[{"role":"user","content":[
		  {"type":"image_url","image_url":{"url":"data:image/png;base64,xxx"}},{"type":"text","text":"what"}]}]}`)
	if !strings.Contains(readAll(t, resp.Body), "vision-ok") {
		t.Error("image part did not produce vision-ok")
	}

	// Audio marker.
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"messages":[{"role":"user","content":[
		  {"type":"input_audio","input_audio":{"data":"xxx","format":"wav"}}]}]}`)
	if !strings.Contains(readAll(t, resp.Body), "audio-ok") {
		t.Error("input_audio part did not produce audio-ok")
	}

	// max_tokens truncation -> finish_reason length.
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"max_tokens":2,"messages":[{"role":"user","content":"a very long answer is coming"}]}`)
	body := readAll(t, resp.Body)
	if !strings.Contains(body, `"finish_reason":"length"`) {
		t.Errorf("max_tokens=2 did not produce finish_reason=length:\n%s", body)
	}

	// [think] marker -> reasoning_content chunks.
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"messages":[{"role":"user","content":"[think] what is 2+2"}]}`)
	if !strings.Contains(readAll(t, resp.Body), "reasoning_content") {
		t.Error("[think] marker did not produce reasoning_content")
	}
}

func TestChatFaultInjectionAndControlState(t *testing.T) {
	ts := newTestServer(t, t.TempDir())

	// Queue two 429s with retry-after, then a normal request succeeds.
	resp := postJSON(t, ts.URL+"/__control/fail_next",
		`{"status":429,"times":2,"retry_after":"7"}`)
	if resp.StatusCode != 200 {
		t.Fatalf("fail_next status = %d", resp.StatusCode)
	}

	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","messages":[{"role":"user","content":"x"}]}`)
	if resp.StatusCode != 429 {
		t.Fatalf("faulted request status = %d, want 429", resp.StatusCode)
	}
	if ra := resp.Header.Get("Retry-After"); ra != "7" {
		t.Errorf("Retry-After = %q, want 7", ra)
	}

	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","messages":[{"role":"user","content":"x"}]}`)
	if resp.StatusCode != 429 {
		t.Fatalf("second faulted request status = %d, want 429", resp.StatusCode)
	}

	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","messages":[{"role":"user","content":"x"}]}`)
	if resp.StatusCode != 200 {
		t.Fatalf("post-fault request status = %d, want 200", resp.StatusCode)
	}

	// State reflects the counters.
	stateResp, err := http.Get(ts.URL + "/__control/state")
	if err != nil {
		t.Fatal(err)
	}
	defer stateResp.Body.Close()
	var state struct {
		RequestsTotal uint64            `json:"requests_total"`
		Routes        map[string]uint64 `json:"routes"`
		QueuedFaults  int               `json:"queued_faults"`
	}
	if err := json.NewDecoder(stateResp.Body).Decode(&state); err != nil {
		t.Fatal(err)
	}
	if state.Routes["chat"] != 3 {
		t.Errorf("routes.chat = %d, want 3", state.Routes["chat"])
	}
	if state.QueuedFaults != 0 {
		t.Errorf("queued_faults = %d, want 0 (all consumed)", state.QueuedFaults)
	}

	// Reset clears counters.
	postJSON(t, ts.URL+"/__control/reset", `{}`)
	stateResp2, err := http.Get(ts.URL + "/__control/state")
	if err != nil {
		t.Fatal(err)
	}
	defer stateResp2.Body.Close()
	var state2 struct {
		RequestsTotal uint64 `json:"requests_total"`
	}
	if err := json.NewDecoder(stateResp2.Body).Decode(&state2); err != nil {
		t.Fatal(err)
	}
	if state2.RequestsTotal != 0 {
		t.Errorf("requests_total after reset = %d, want 0", state2.RequestsTotal)
	}
}

func TestChatTruncateCutsStreamBeforeDone(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	postJSON(t, ts.URL+"/__control/truncate", `{"chunks":2}`)

	resp := postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"mock-small","stream":true,"messages":[{"role":"user","content":"one two three four five"}]}`)
	body := readAll(t, resp.Body)
	if strings.Contains(body, "[DONE]") {
		t.Error("truncated stream still delivered [DONE]")
	}
	if n := len(sseChunks(t, body)); n != 2 {
		t.Errorf("chunk count after truncate{2} = %d, want 2", n)
	}
}

func TestChatGoldenReplayByteIdentity(t *testing.T) {
	ts := newTestServer(t, "testdata/golden")

	// Raw golden mode: served body must equal the recorded section bytes.
	resp := postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"golden/tool-call","stream":true,"messages":[]}`)
	raw1 := readAll(t, resp.Body)

	recorded, err := os.ReadFile(filepath.Join("testdata", "golden", "tool-call.sse"))
	if err != nil {
		t.Fatal(err)
	}
	// The recorded file is header + directives + body; the served body is the
	// section after the @response directive.
	if i := strings.Index(string(recorded), "data: "); i < 0 {
		t.Fatal("fixture has no data lines")
	} else {
		want := string(recorded[i:])
		if raw1 != want {
			t.Errorf("served golden bytes differ from fixture:\n--- served ---\n%q\n--- fixture ---\n%q", raw1, want)
		}
	}

	// Sections 2..n of the ladder: 429, 429, 200, then exhausted.
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"golden/backoff-429","stream":true,"messages":[]}`)
	if got := resp.StatusCode; got != 429 {
		t.Errorf("ladder response 2 status = %d, want 429", got)
	}
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"golden/backoff-429","stream":true,"messages":[]}`)
	if got := resp.StatusCode; got != 429 {
		t.Errorf("ladder response 3 status = %d, want 429", got)
	}
	if ra := resp.Header.Get("Retry-After"); ra != "1" {
		t.Errorf("ladder Retry-After = %q, want 1", ra)
	}
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"golden/backoff-429","stream":true,"messages":[]}`)
	if resp.StatusCode != 200 {
		t.Errorf("ladder response 4 status = %d, want 200", resp.StatusCode)
	}
	// Fifth -> exhausted.
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"golden/backoff-429","stream":true,"messages":[]}`)
	if resp.StatusCode != 500 {
		t.Errorf("exhausted golden status = %d, want 500", resp.StatusCode)
	}

	// Reset clears the cursor: replay from the top works again.
	postJSON(t, ts.URL+"/__control/reset", `{}`)
	resp = postJSON(t, ts.URL+"/v1/chat/completions",
		`{"model":"golden/backoff-429","stream":true,"messages":[]}`)
	if resp.StatusCode != 429 {
		t.Errorf("after reset golden status = %d, want 429 (fresh cursor)", resp.StatusCode)
	}
}

func TestModelsEndpoint(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	resp, err := http.Get(ts.URL + "/v1/models")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got struct {
		Data []struct {
			ID      string `json:"id"`
			Object  string `json:"object"`
			OwnedBy string `json:"owned_by"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if len(got.Data) != 2 || got.Data[0].ID != "mock-small" || got.Data[1].ID != "mock-large" {
		t.Errorf("models = %+v", got.Data)
	}
	if got.Data[0].Object != "model" || got.Data[0].OwnedBy != "mockllm" {
		t.Errorf("model entry shape wrong: %+v", got.Data[0])
	}
}

func TestMessagesAndResponsesSmoke(t *testing.T) {
	ts := newTestServer(t, t.TempDir())

	// Anthropic-style non-stream.
	resp := postJSON(t, ts.URL+"/v1/messages",
		`{"model":"mock-large","max_tokens":100,"messages":[{"role":"user","content":"hi anthropic"}]}`)
	var msg struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		StopReason string `json:"stop_reason"`
		Usage      struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		t.Fatal(err)
	}
	if msg.StopReason != "end_turn" || len(msg.Content) == 0 || msg.Content[0].Text != "echo: hi anthropic" {
		t.Errorf("messages non-stream = %+v", msg)
	}

	// Anthropic-style stream: message_start/content_block_delta/message_stop.
	resp = postJSON(t, ts.URL+"/v1/messages",
		`{"model":"mock-large","max_tokens":100,"stream":true,"messages":[{"role":"user","content":"hello"}]}`)
	body := readAll(t, resp.Body)
	for _, want := range []string{"message_start", "content_block_start", "content_block_delta", "message_delta", "message_stop"} {
		if !strings.Contains(body, want) {
			t.Errorf("anthropic stream missing %s", want)
		}
	}

	// Responses API stream smoke.
	resp = postJSON(t, ts.URL+"/v1/responses",
		`{"model":"mock-large","stream":true,"input":"respond to this"}`)
	body = readAll(t, resp.Body)
	for _, want := range []string{"response.created", "response.output_text.delta", "response.completed", "[DONE]"} {
		if !strings.Contains(body, want) {
			t.Errorf("responses stream missing %s", want)
		}
	}
}

func TestGoldenFormatParserMirrorsFixture(t *testing.T) {
	// The mockllm parser must produce the same section structure as the
	// main-module parser on the same fixture (checked structurally here;
	// byte identity through the adapter is the integration test's job).
	src, err := os.ReadFile(filepath.Join("testdata", "golden", "backoff-429.sse"))
	if err != nil {
		t.Fatal(err)
	}
	sections, err := parseGolden(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 3 {
		t.Fatalf("sections = %d, want 3", len(sections))
	}
	if sections[0].status != 429 || sections[0].retryAfter != "1" {
		t.Errorf("section 0 = %+v", sections[0])
	}
	if !strings.HasPrefix(sections[0].body, `{"error"`) {
		t.Errorf("section 0 body = %q", sections[0].body)
	}
	if sections[2].status != 200 || !strings.Contains(sections[2].body, "data: [DONE]") {
		t.Errorf("section 2 = %+v", sections[2])
	}
	_ = io.Discard
}
