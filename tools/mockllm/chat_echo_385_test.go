package main

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

// Ticket 385: the synthesized echo used to hard-clamp at 400 BYTES
// (len(answer) > 400 -> answer[:400]). 400 bytes is not a rune boundary for
// CJK text, so the clamp cut mid-UTF-8-sequence - the JSON encoder then
// carried the broken tail as U+FFFD - and made long samples physically
// undeliverable, which is why ticket 30's 402-hanzi case had to detour to an
// in-process C5 provider (internal/present/e2e_present_30_test.go). The
// clamp is gone in all three dialects; these tests pin that the exact sample
// the C5 seam case scripts comes back in full, so that after stripping the
// mock's "echo: " signature the two paths' outputs are byte-identical.

// ticket30Sample is the tier-2 sample from internal/present's
// e2e_present_30_test.go: 402 hanzi (1206 bytes) plus the same-response
// summary line. Behind the "echo: " signature that is 1243 bytes; the old
// clamp kept 400 bytes and cut the 395th byte inside the 132nd hanzi.
func ticket30Sample() string {
	return strings.Repeat("测", 402) + "\n口播摘要：长文已归档"
}

func tailRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[len(r)-n:])
}

func Test385EchoDeliversLongMultibyteInFull(t *testing.T) {
	sample := ticket30Sample()
	want := "echo: " + sample
	if len(want) <= 400 {
		t.Fatalf("fixture is only %d bytes; this test guards the >400-byte shape", len(want))
	}
	ts := newTestServer(t, t.TempDir())

	t.Run("chat_stream_sse", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-small", "stream": true,
			"messages": []map[string]any{{"role": "user", "content": sample}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/chat/completions", string(body))
		text := ""
		for _, c := range sseChunks(t, readAll(t, resp.Body)) {
			data := chunkData(t, c)
			if data == "[DONE]" {
				continue
			}
			var m struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &m); err != nil {
				t.Fatalf("chunk not JSON: %v", err)
			}
			if len(m.Choices) > 0 {
				text += m.Choices[0].Delta.Content
			}
		}
		if !utf8.ValidString(text) {
			t.Fatalf("assembled text is not valid UTF-8: a clamp cut mid-rune (tail %q)", tailRunes(text, 20))
		}
		if text != want {
			t.Errorf("assembled text = %d bytes / %d runes, want %d bytes / %d runes (tail lost: %q)",
				len(text), len([]rune(text)), len(want), len([]rune(want)), tailRunes(text, 30))
		}
	})

	t.Run("chat_nonstream_json", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model":    "mock-small",
			"messages": []map[string]any{{"role": "user", "content": sample}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/chat/completions", string(body))
		var out struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatal(err)
		}
		if len(out.Choices) != 1 || out.Choices[0].Message.Content != want {
			t.Errorf("content = %d bytes, want %d bytes with the sample intact",
				len(out.Choices[0].Message.Content), len(want))
		}
	})

	t.Run("messages_and_responses_dialects", func(t *testing.T) {
		// Anthropic Messages, non-stream (shape pinned by
		// TestMessagesAndResponsesSmoke): content[0].text is the echo.
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_tokens": 4096,
			"messages": []map[string]any{{"role": "user", "content": sample}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/messages", string(body))
		var msg struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			t.Fatal(err)
		}
		if len(msg.Content) == 0 || msg.Content[0].Text != want {
			t.Errorf("messages echo = %d bytes, want the sample intact behind the signature",
				len(msg.Content[0].Text))
		}

		// Responses, non-stream: the message output item's output_text is
		// the echo.
		body, err = json.Marshal(map[string]any{"model": "mock-large", "input": sample})
		if err != nil {
			t.Fatal(err)
		}
		resp = postJSON(t, ts.URL+"/v1/responses", string(body))
		var rs struct {
			Output []struct {
				Type    string `json:"type"`
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			} `json:"output"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
			t.Fatal(err)
		}
		got := ""
		for _, item := range rs.Output {
			if item.Type == "message" && len(item.Content) > 0 {
				got = item.Content[0].Text
			}
		}
		if got != want {
			t.Errorf("responses echo = %d bytes, want the sample intact behind the signature", len(got))
		}
	})
}
