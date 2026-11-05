package openairesponses

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// Ticket 23 AC#2, openai-responses leg: function_call_output takes an ARRAY of
// input items, and input_image is one of them - so the pixels stay attached to
// the call they came from. Before this test the adapter returned "only
// TextPart" and failed the request.
func TestToolResultImageReachesResponsesWire23(t *testing.T) {
	data := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x23}
	ref := "mem:23"

	req := &llm.Request{
		Model: "mock-small",
		Messages: []llm.Message{
			{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "截个图看看"}}},
			{Role: llm.RoleTool, Content: []llm.Content{llm.ToolResultPart{
				ID: "call_e1",
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
		Input []struct {
			Type   string `json:"type"`
			CallID string `json:"call_id"`
			Output []struct {
				Type     string `json:"type"`
				Text     string `json:"text"`
				ImageURL string `json:"image_url"`
			} `json:"output"`
		} `json:"input"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		t.Fatalf("decode built body: %v", err)
	}
	var out *struct {
		Type   string `json:"type"`
		CallID string `json:"call_id"`
		Output []struct {
			Type     string `json:"type"`
			Text     string `json:"text"`
			ImageURL string `json:"image_url"`
		} `json:"output"`
	}
	for i := range wire.Input {
		if wire.Input[i].Type == "function_call_output" {
			out = &wire.Input[i]
		}
	}
	if out == nil {
		t.Fatalf("no function_call_output item in %s", body)
	}
	if out.CallID != "call_e1" {
		t.Errorf("call_id = %q, want call_e1", out.CallID)
	}
	wantURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
	var sawText, sawImage bool
	for _, item := range out.Output {
		switch item.Type {
		case "input_text":
			sawText = true
		case "input_image":
			sawImage = true
			if item.ImageURL != wantURL {
				t.Errorf("input_image url = %q, want the data URL of the resolved bytes", item.ImageURL)
			}
		}
	}
	if !sawText || !sawImage {
		t.Fatalf("function_call_output.output = %+v, want input_text AND input_image", out.Output)
	}
	if strings.Contains(string(body), ref) {
		t.Errorf("the bytes ref %q leaked onto the wire", ref)
	}
}
