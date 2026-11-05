package anthropic

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// Ticket 23 AC#2, anthropic leg: Anthropic's tool_result block natively accepts
// image blocks, so the faithful mapping keeps the image INSIDE the tool result
// (no synthetic user turn, unlike openai-chat). Before this test the adapter
// returned "only TextPart" and failed the request.
func TestToolResultImageReachesAnthropicWire23(t *testing.T) {
	data := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x23}
	ref := "mem:23"

	req := &llm.Request{
		Model:           "mock-small",
		MaxOutputTokens: 64,
		Messages: []llm.Message{
			{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "截个图看看"}}},
			{Role: llm.RoleTool, Content: []llm.Content{llm.ToolResultPart{
				ID: "toolu_1",
				Content: []llm.Content{
					llm.TextPart{Text: `{"mime":"image/png","bytes_len":9,"drm":false}`},
					llm.ImagePart{MimeType: "image/png", BytesRef: ref, Alt: "桌面截图"},
				},
			}}},
		},
		ResolveBytes: func(got string) (string, []byte, error) {
			if got != ref {
				t.Fatalf("ResolveBytes ref = %q, want %q", got, ref)
			}
			return "image/png", data, nil
		},
	}

	a := New(endpointOptionsFor("http://127.0.0.1:1", false, true))
	body, err := a.buildWire(req)
	if err != nil {
		t.Fatalf("buildWire refused a tool result carrying an image: %v", err)
	}

	var wire struct {
		Messages []struct {
			Role    string `json:"role"`
			Content []struct {
				Type      string `json:"type"`
				ToolUseID string `json:"tool_use_id"`
				Content   []struct {
					Type   string `json:"type"`
					Text   string `json:"text"`
					Source struct {
						Type      string `json:"type"`
						MediaType string `json:"media_type"`
						Data      string `json:"data"`
					} `json:"source"`
				} `json:"content"`
			} `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		t.Fatalf("decode built body: %v", err)
	}
	if len(wire.Messages) != 2 {
		t.Fatalf("wire messages = %d, want 2 (user, tool-result user): %s", len(wire.Messages), body)
	}
	if wire.Messages[1].Role != "user" || len(wire.Messages[1].Content) != 1 {
		t.Fatalf("second wire message = %+v, want the one tool_result block", wire.Messages[1])
	}
	block := wire.Messages[1].Content[0]
	if block.Type != "tool_result" || block.ToolUseID != "toolu_1" {
		t.Fatalf("block = %+v, want tool_result for toolu_1", block)
	}
	wantData := base64.StdEncoding.EncodeToString(data)
	var sawText, sawImage bool
	for _, inner := range block.Content {
		switch inner.Type {
		case "text":
			sawText = true
		case "image":
			sawImage = true
			if inner.Source.Type != "base64" || inner.Source.MediaType != "image/png" {
				t.Errorf("image source = %+v, want base64/image/png", inner.Source)
			}
			if inner.Source.Data != wantData {
				t.Errorf("image source data = %q, want the base64 of the resolved bytes", inner.Source.Data)
			}
		}
	}
	if !sawText || !sawImage {
		t.Fatalf("tool_result content = %+v, want a text block AND an image block", block.Content)
	}
	if strings.Contains(string(body), ref) {
		t.Errorf("the bytes ref %q leaked onto the wire", ref)
	}
}
