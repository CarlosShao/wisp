package anthropic

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
)

// AC#2: cache-breakpoint behavior proven from CAPTURED REQUEST BODIES served
// by the live mockllm (/__control/last_request), not from adapter bookkeeping.
// A breakpoint that was computed but never reached the wire, or one stamped on
// the volatile suffix, fails these tests.

// stableSystem is the cache-stable prefix (SPEC-05 sec 4.1 sections 1/5/6/2).
func stableSystem() []llm.Content {
	return []llm.Content{
		llm.TextPart{Text: "①identity: Wisp voice agent. Never write code."},
		llm.TextPart{Text: "⑤safety: confirm before destructive actions."},
		llm.TextPart{Text: "⑥style: short spoken answers."},
	}
}

func turn(msgs ...llm.Message) *llm.Request {
	return &llm.Request{
		Model:            "mock-small",
		System:           stableSystem(),
		Messages:         msgs,
		CacheBreakpoints: []int{-1}, // after System, before the conversation
		Tools: []llm.ToolDef{{Name: "get_weather", Description: "d",
			Parameters: json.RawMessage(`{"type":"object"}`)}},
	}
}

func userMsg(text string) llm.Message {
	return llm.Message{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: text}}}
}

// TestCacheBreakpointReachesTheWire asserts the marker's exact placement in
// the request body the provider actually received.
func TestCacheBreakpointReachesTheWire(t *testing.T) {
	proc := adaptertest.StartMockllm(t)
	proc.Reset(t)
	proc.StreamGolden(t, unit(), "anthropic-tool-call", turn(userMsg("hi")))

	body, hdr := proc.LastRequest(t, "messages")
	var sent map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &sent); err != nil {
		t.Fatalf("captured body is not JSON: %v\n%s", err, body)
	}

	// 1. The system block array carries exactly one cache_control, on its
	//    LAST block (the prefix boundary).
	var sys []struct {
		Type         string          `json:"type"`
		Text         string          `json:"text"`
		CacheControl json.RawMessage `json:"cache_control"`
	}
	if err := json.Unmarshal(sent["system"], &sys); err != nil {
		t.Fatalf("system field: %v\n%s", err, sent["system"])
	}
	if len(sys) != 3 {
		t.Fatalf("system blocks = %d, want the 3 stable prefix sections", len(sys))
	}
	for i, b := range sys {
		last := i == len(sys)-1
		if last != (len(b.CacheControl) > 0) {
			t.Errorf("system block %d cache_control present=%v, want %v",
				i, len(b.CacheControl) > 0, last)
		}
		if last && string(b.CacheControl) != `{"type":"ephemeral"}` {
			t.Errorf("cache_control = %s, want {\"type\":\"ephemeral\"}", b.CacheControl)
		}
	}

	// 2. Nothing in the conversation is cached (suffix never cached).
	if bytes.Contains([]byte(string(sent["messages"])), []byte("cache_control")) {
		t.Errorf("cache_control leaked into the conversation suffix:\n%s", sent["messages"])
	}
	// 3. Tools stay unmarked too (a tool block breakpoint is not declared).
	if bytes.Contains([]byte(string(sent["tools"])), []byte("cache_control")) {
		t.Errorf("cache_control leaked into the tool definitions:\n%s", sent["tools"])
	}

	// 4. The protocol's auth scheme is x-api-key + anthropic-version (not a
	//    Bearer token): a wrong header means the request would 401 live.
	if hdr.Get("Anthropic-Version") != AnthropicVersion {
		t.Errorf("anthropic-version = %q, want %q", hdr.Get("Anthropic-Version"), AnthropicVersion)
	}
	if hdr.Get("x-api-key") == "" {
		t.Error("x-api-key header missing (the seam sends the key here)")
	}
	if auth := hdr.Get("Authorization"); auth != "" {
		t.Errorf("Authorization header = %q, want none on the anthropic protocol", auth)
	}
	proc.Reset(t)
}

// TestCachedPrefixIsByteStableAcrossTurns is the cost statement behind the
// breakpoint: across two consecutive turns the cached prefix must be
// byte-identical, or the provider cache never hits.
func TestCachedPrefixIsByteStableAcrossTurns(t *testing.T) {
	proc := adaptertest.StartMockllm(t)
	a := New(endpointOptionsFor(proc.Base+"/v1", false, false))

	// Turn 1: one user message.
	proc.Reset(t)
	if _, _, err := adaptertest.Drain(t, a, turn(userMsg("今天天气怎么样"))); err != nil {
		t.Fatal(err)
	}
	first, _ := proc.LastRequest(t, "messages")

	// Turn 2: the same prefix, a longer conversation and a volatile suffix
	// message (time/focus section) appended AFTER the breakpoint.
	r2 := turn(userMsg("今天天气怎么样"),
		llm.Message{Role: llm.RoleAssistant, Content: []llm.Content{
			llm.ToolUsePart{ID: "call_a1", Name: "get_weather",
				Input: json.RawMessage(`{"city":"Zhuhai"}`)}}},
		llm.Message{Role: llm.RoleTool, Content: []llm.Content{
			llm.ToolResultPart{ID: "call_a1",
				Content: []llm.Content{llm.TextPart{Text: "sunny 27C"}}}}},
		userMsg("现在几点"))
	r2.CacheBreakpoints = []int{-1}
	if _, _, err := adaptertest.Drain(t, a, r2); err != nil {
		t.Fatal(err)
	}
	second, _ := proc.LastRequest(t, "messages")

	fieldOf := func(body, key string) string {
		var m map[string]json.RawMessage
		if err := json.Unmarshal([]byte(body), &m); err != nil {
			t.Fatalf("captured body: %v", err)
		}
		return string(m[key])
	}
	if p1, p2 := fieldOf(first, "system"), fieldOf(second, "system"); p1 != p2 {
		t.Errorf("cached prefix drifted between turns:\n turn1: %s\n turn2: %s", p1, p2)
	}
	if fieldOf(first, "system") == "" {
		t.Fatal("no system field captured")
	}
	// The first conversation turn is inside the cached prefix too (the
	// breakpoint is after System, so the prefix is System only); the earlier
	// user message must nevertheless be byte-identical in both bodies.
	var msgs1, msgs2 []json.RawMessage
	if err := json.Unmarshal([]byte(fieldOf(first, "messages")), &msgs1); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(fieldOf(second, "messages")), &msgs2); err != nil {
		t.Fatal(err)
	}
	if string(msgs1[0]) != string(msgs2[0]) {
		t.Errorf("first conversation turn changed between turns:\n %s\n %s", msgs1[0], msgs2[0])
	}
	if len(msgs2) <= len(msgs1) {
		t.Errorf("turn 2 should carry more history: %d vs %d", len(msgs2), len(msgs1))
	}
	// Only ONE breakpoint total in each body.
	if n := strings.Count(first, "cache_control"); n != 1 {
		t.Errorf("turn 1 body has %d cache_control markers, want exactly 1", n)
	}
	if n := strings.Count(second, "cache_control"); n != 1 {
		t.Errorf("turn 2 body has %d cache_control markers, want exactly 1", n)
	}
	proc.Reset(t)
}

// TestCacheReadUsagePopulatesCachedTokens proves Usage.CachedTokens is mapped
// from the provider's cache-read counter (the fixture reports
// cache_read_input_tokens: 8), via the live second runner.
func TestCacheReadUsagePopulatesCachedTokens(t *testing.T) {
	proc := adaptertest.StartMockllm(t)
	events := proc.StreamGolden(t, unit(), "anthropic-tool-call", turn(userMsg("hi")))
	var found llm.Usage
	for _, ev := range events {
		if ev.Type == llm.EvUsage {
			found = ev.Usage.Add(found)
		}
	}
	if found.CachedTokens != 8 {
		t.Errorf("cached_tokens = %d, want 8 (mapped from cache_read_input_tokens); events %v",
			found.CachedTokens, events)
	}
	if found.InputTokens != 21 {
		t.Errorf("input_tokens = %d, want 21 verbatim (anthropic reports non-cached input separately; see the package usage-mapping note)",
			found.InputTokens)
	}
}

// TestMidConversationBreakpointPlacement covers CacheBreakpoints >= 0: the
// marker goes on the LAST block of that message.
func TestMidConversationBreakpointPlacement(t *testing.T) {
	proc := adaptertest.StartMockllm(t)
	proc.Reset(t)
	a := New(endpointOptionsFor(proc.Base+"/v1", false, false))
	req := turn(userMsg("stable tool-index section"), userMsg("volatile tail"))
	req.CacheBreakpoints = []int{0} // after Messages[0]
	if _, _, err := adaptertest.Drain(t, a, req); err != nil {
		t.Fatal(err)
	}
	body, _ := proc.LastRequest(t, "messages")
	var sent map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &sent); err != nil {
		t.Fatal(err)
	}
	var msgs []struct {
		Role    string `json:"role"`
		Content []struct {
			Type         string          `json:"type"`
			CacheControl json.RawMessage `json:"cache_control"`
		} `json:"content"`
	}
	if err := json.Unmarshal(sent["messages"], &msgs); err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("messages = %d, want 2", len(msgs))
	}
	if len(msgs[0].Content[0].CacheControl) == 0 {
		t.Error("Messages[0] carries no cache_control despite the declared breakpoint")
	}
	if len(msgs[1].Content[0].CacheControl) != 0 {
		t.Error("the message AFTER the breakpoint must not be cached")
	}
	if strings.Count(body, "cache_control") != 1 {
		t.Errorf("expected exactly one marker, body:\n%s", body)
	}
	proc.Reset(t)
}

// TestBreakpointLimit: more than the protocol's 4 markers is our bug (strict)
// or a documented trim (compat.loose).
func TestBreakpointLimit(t *testing.T) {
	mk := func(n int) *llm.Request {
		msgs := make([]llm.Message, 0, n)
		breaks := make([]int, 0, n)
		for i := 0; i < n; i++ {
			msgs = append(msgs, userMsg("m"))
			breaks = append(breaks, i)
		}
		r := turn(msgs...)
		r.CacheBreakpoints = breaks
		return r
	}
	a := New(endpointOptionsFor("http://127.0.0.1:1/v1", false, false))
	if _, _, err := adaptertest.Drain(t, a, mk(5)); err == nil {
		t.Error("5 breakpoints must fail in strict mode (protocol max is 4)")
	} else if !strings.Contains(err.Error(), "exceed the protocol maximum") {
		t.Errorf("err = %v, want the actionable max-breakpoints message", err)
	}
	loose := New(endpointOptionsFor("http://127.0.0.1:1/v1", true, false))
	if _, _, err := adaptertest.Drain(t, loose, mk(5)); err == nil ||
		!strings.Contains(err.Error(), "dial tcp") {
		t.Errorf("loose mode must trim and proceed to the transport, got %v", err)
	}
	// Out-of-range index is always an internal error.
	bad := turn(userMsg("m"))
	bad.CacheBreakpoints = []int{7}
	if _, _, err := adaptertest.Drain(t, a, bad); err == nil ||
		!strings.Contains(err.Error(), "out of range") {
		t.Errorf("err = %v, want out-of-range breakpoint rejection", err)
	}
}

// TestInfoReportsExplicitBreakpoints is the C5 cache-capability bit the agent
// core reads to decide whether to declare breakpoints at all.
func TestInfoReportsExplicitBreakpoints(t *testing.T) {
	info := New(endpointOptionsFor("http://127.0.0.1:1/v1", false, false)).Info()
	if info.Cache.Mode != llm.CacheExplicit || !info.Cache.Breakpoints {
		t.Errorf("cache bits = %+v, want explicit/breakpoints", info.Cache)
	}
	if info.Protocol != Protocol {
		t.Errorf("protocol = %s", info.Protocol)
	}
}
