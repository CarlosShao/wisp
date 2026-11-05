package openaichat

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// Ticket 23 AC#2, openai-chat leg: a tool result that carries an image must
// reach the wire.
//
// The chat-completions protocol has no image part inside a role:"tool" message
// (content there is text only), so the adapter's honest mapping - the one every
// mainstream harness ships - is to keep the tool message text and attach the
// image to the follow-up user turn. What must NOT happen is what happened
// before this test existed: the adapter returning "only TextPart" and failing
// the whole request, i.e. one screenshot killing the conversation on that
// provider.
func TestToolResultImageReachesChatWire23(t *testing.T) {
	data := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x23, 0x23}
	ref := "mem:23"
	const meta = `{"mime":"image/png","bytes_len":10,"drm":false}`

	var asked []string
	req := &llm.Request{
		Model: "mock-small",
		Messages: []llm.Message{
			{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "截个图看看"}}},
			{Role: llm.RoleTool, Content: []llm.Content{llm.ToolResultPart{
				ID: "call_e1",
				Content: []llm.Content{
					llm.TextPart{Text: meta},
					llm.ImagePart{MimeType: "image/png", BytesRef: ref, Alt: "桌面截图"},
				},
			}}},
		},
		ResolveBytes: func(got string) (string, []byte, error) {
			asked = append(asked, got)
			return "image/png", data, nil
		},
	}

	a := New(endpointOptionsFor("http://127.0.0.1:1", false, true))
	body, err := a.buildWire(req)
	if err != nil {
		t.Fatalf("buildWire refused a tool result carrying an image: %v", err)
	}
	if len(asked) != 1 || asked[0] != ref {
		t.Fatalf("ResolveBytes calls = %v, want exactly the one ref %q", asked, ref)
	}

	var wire struct {
		Messages []struct {
			Role       string          `json:"role"`
			Content    json.RawMessage `json:"content"`
			ToolCallID string          `json:"tool_call_id"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		t.Fatalf("decode built body: %v", err)
	}
	if len(wire.Messages) != 3 {
		t.Fatalf("wire messages = %d, want 3 (user, tool, user carrying the image): %s",
			len(wire.Messages), body)
	}
	if wire.Messages[1].Role != "tool" || wire.Messages[1].ToolCallID != "call_e1" {
		t.Fatalf("message 1 = %+v, want the role:tool answer to call_e1", wire.Messages[1])
	}
	var toolText string
	if err := json.Unmarshal(wire.Messages[1].Content, &toolText); err != nil {
		t.Fatalf("tool content is not the plain string form: %v (%s)", err, wire.Messages[1].Content)
	}
	if !strings.Contains(toolText, "bytes_len") {
		t.Errorf("tool message lost the metadata text: %q", toolText)
	}
	if strings.Contains(toolText, base64.StdEncoding.EncodeToString(data)) {
		t.Errorf("tool message carries base64 pixels; the protocol only takes text there")
	}

	last := wire.Messages[2]
	if last.Role != "user" {
		t.Fatalf("message 2 role = %q, want user (the image turn)", last.Role)
	}
	var parts []struct {
		Type     string `json:"type"`
		Text     string `json:"text"`
		ImageURL struct {
			URL string `json:"url"`
		} `json:"image_url"`
	}
	if err := json.Unmarshal(last.Content, &parts); err != nil {
		t.Fatalf("user content is not a parts array: %v (%s)", err, last.Content)
	}
	wantURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
	var sawText, sawImage bool
	for _, p := range parts {
		switch p.Type {
		case "text":
			sawText = true
		case "image_url":
			sawImage = true
			if p.ImageURL.URL != wantURL {
				t.Errorf("image part url = %q, want the data URL of the resolved bytes", p.ImageURL.URL)
			}
		}
	}
	if !sawText || !sawImage {
		t.Fatalf("image turn parts = %s, want a text label AND the image_url part", last.Content)
	}
	if strings.Contains(string(body), ref) {
		t.Errorf("the bytes ref %q leaked onto the wire: the adapter must resolve it, not forward it", ref)
	}

	// A tool result that carries only text keeps the old one-message shape:
	// this change must not add an empty user turn to every tool call.
	plain := &llm.Request{
		Model: "mock-small",
		Messages: []llm.Message{
			{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "查一下"}}},
			{Role: llm.RoleTool, Content: []llm.Content{llm.ToolResultPart{
				ID: "call_e2", Content: []llm.Content{llm.TextPart{Text: "22 摄氏度"}},
			}}},
		},
	}
	plainBody, err := a.buildWire(plain)
	if err != nil {
		t.Fatalf("buildWire text-only tool result: %v", err)
	}
	var plainWire struct {
		Messages []map[string]any `json:"messages"`
	}
	if err := json.Unmarshal(plainBody, &plainWire); err != nil {
		t.Fatalf("decode text-only body: %v", err)
	}
	if len(plainWire.Messages) != 2 {
		t.Fatalf("text-only wire messages = %d, want 2 (no synthetic image turn)", len(plainWire.Messages))
	}
}
