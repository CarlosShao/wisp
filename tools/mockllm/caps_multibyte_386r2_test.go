package main

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

// Ticket 386: the per-cap verdicts for the caps 385-r1 registered unfixed
// (evidence docs/evidence/s1/386r2-*.md):
//
//   - reasoning/reason 200-byte cap, three dialects (chat.go / responses.go /
//     messages.go): unintentional defect, removed. No real provider caps
//     reasoning output at 200 bytes; the cap simulated nothing and cut CJK
//     text mid-UTF-8-sequence exactly like the echo clamp 385 removed.
//   - max_tokens byte cut (chat.go): the cap is intentional emulation (a real
//     chat API that hits max_tokens stops mid-output and reports
//     finish_reason "length") and is kept with its heuristic; only the cut
//     point was unintentionally byte-exact. Real providers truncate at token
//     boundaries and never emit invalid UTF-8, so the cut now lands on a
//     rune boundary. ASCII behavior is byte-identical (pinned below).

func think386Want(sample string) string {
	return "thinking about: [think] " + sample
}

func Test386ReasoningDeliversLongMultibyteInFull(t *testing.T) {
	sample := ticket30Sample()
	want := think386Want(sample)
	if len(want) <= 200 {
		t.Fatalf("fixture is only %d bytes; this test guards the >200-byte shape", len(want))
	}
	ts := newTestServer(t, t.TempDir())

	t.Run("chat_stream_sse", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-small", "stream": true,
			"messages": []map[string]any{{"role": "user", "content": "[think] " + sample}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/chat/completions", string(body))
		reasoning := ""
		for _, c := range sseChunks(t, readAll(t, resp.Body)) {
			data := chunkData(t, c)
			if data == "[DONE]" {
				continue
			}
			var m struct {
				Choices []struct {
					Delta struct {
						ReasoningContent string `json:"reasoning_content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &m); err != nil {
				t.Fatalf("chunk not JSON: %v", err)
			}
			if len(m.Choices) > 0 {
				reasoning += m.Choices[0].Delta.ReasoningContent
			}
		}
		if !utf8.ValidString(reasoning) {
			t.Fatalf("assembled reasoning is not valid UTF-8: a clamp cut mid-rune (tail %q)", tailRunes(reasoning, 20))
		}
		if reasoning != want {
			t.Errorf("assembled reasoning = %d bytes / %d runes, want %d bytes / %d runes (tail lost: %q)",
				len(reasoning), len([]rune(reasoning)), len(want), len([]rune(want)), tailRunes(reasoning, 30))
		}
	})

	t.Run("chat_nonstream_json", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model":    "mock-small",
			"messages": []map[string]any{{"role": "user", "content": "[think] " + sample}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/chat/completions", string(body))
		var out struct {
			Choices []struct {
				Message struct {
					ReasoningContent string `json:"reasoning_content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatal(err)
		}
		if len(out.Choices) == 0 || out.Choices[0].Message.ReasoningContent != want {
			got := ""
			if len(out.Choices) > 0 {
				got = out.Choices[0].Message.ReasoningContent
			}
			t.Errorf("reasoning_content = %d bytes / %d runes, want the sample intact behind the signature (tail %q)",
				len(got), len([]rune(got)), tailRunes(got, 30))
		}
	})

	t.Run("messages_thinking_block", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_tokens": 4096,
			"messages": []map[string]any{{"role": "user", "content": "[think] " + sample}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/messages", string(body))
		var msg struct {
			Content []struct {
				Type     string `json:"type"`
				Thinking string `json:"thinking"`
			} `json:"content"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			t.Fatal(err)
		}
		got := ""
		for _, b := range msg.Content {
			if b.Type == "thinking" {
				got = b.Thinking
			}
		}
		if got != want {
			t.Errorf("messages thinking = %d bytes / %d runes, want the sample intact behind the signature (tail %q)",
				len(got), len([]rune(got)), tailRunes(got, 30))
		}
	})

	t.Run("responses_reasoning_summary", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{"model": "mock-large", "input": "[think] " + sample})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/responses", string(body))
		var rs struct {
			Output []struct {
				Type    string `json:"type"`
				Summary []struct {
					Text string `json:"text"`
				} `json:"summary"`
			} `json:"output"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
			t.Fatal(err)
		}
		got := ""
		for _, item := range rs.Output {
			if item.Type == "reasoning" && len(item.Summary) > 0 {
				got = item.Summary[0].Text
			}
		}
		if got != want {
			t.Errorf("responses reasoning summary = %d bytes / %d runes, want the sample intact behind the signature (tail %q)",
				len(got), len([]rune(got)), tailRunes(got, 30))
		}
	})
}

func Test386MaxTokensTruncationIsRuneSafe(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	// 300 hanzi = 900 bytes; "echo: " + sample = 906 bytes; the ~4-bytes-per-
	// token heuristic reads 226 tokens > max_tokens 100, so the emulation
	// must truncate. Within the 400-byte bound the longest whole-rune prefix
	// is "echo: " + 131 hanzi = 399 bytes; the old byte cut kept 400 bytes and
	// split the 132nd hanzi mid-sequence (the JSON encoder carried the broken
	// tail as U+FFFD).
	const maxTokens = 100
	want := "echo: " + strings.Repeat("测", 131)

	t.Run("hanzi_cut_lands_on_rune_boundary", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-small", "max_tokens": maxTokens,
			"messages": []map[string]any{{"role": "user", "content": strings.Repeat("测", 300)}},
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
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatal(err)
		}
		if len(out.Choices) != 1 {
			t.Fatalf("choices = %d, want 1", len(out.Choices))
		}
		got := out.Choices[0].Message.Content
		if out.Choices[0].FinishReason != "length" {
			t.Errorf("finish_reason = %q, want %q: the max_tokens emulation must keep its length semantics", out.Choices[0].FinishReason, "length")
		}
		if !utf8.ValidString(got) {
			t.Errorf("truncated content is not valid UTF-8: the cut landed mid-rune (tail %q)", tailRunes(got, 20))
		}
		if got != want {
			t.Errorf("content = %d bytes / %d runes, want the rune-boundary prefix %d bytes / %d runes (tail %q)",
				len(got), len([]rune(got)), len(want), len([]rune(want)), tailRunes(got, 20))
		}
		if len(got) > maxTokens*4 {
			t.Errorf("content = %d bytes, exceeds the max_tokens*4 heuristic bound %d", len(got), maxTokens*4)
		}
	})

	t.Run("stream_finish_reason_length", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-small", "stream": true, "max_tokens": maxTokens,
			"messages": []map[string]any{{"role": "user", "content": strings.Repeat("测", 300)}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/chat/completions", string(body))
		text, finish := "", ""
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
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &m); err != nil {
				t.Fatalf("chunk not JSON: %v", err)
			}
			if len(m.Choices) > 0 {
				text += m.Choices[0].Delta.Content
				if m.Choices[0].FinishReason != nil {
					finish = *m.Choices[0].FinishReason
				}
			}
		}
		if finish != "length" {
			t.Errorf("stream finish_reason = %q, want %q", finish, "length")
		}
		if text != want {
			t.Errorf("assembled text = %d bytes / %d runes, want the rune-boundary prefix %d bytes / %d runes (tail %q)",
				len(text), len([]rune(text)), len(want), len([]rune(want)), tailRunes(text, 20))
		}
	})

	t.Run("ascii_cut_byte_identical_to_cap", func(t *testing.T) {
		// Control for the kept cap: for ASCII the rune-boundary cut is the
		// byte cut, so the emulation still truncates exactly at max_tokens*4
		// bytes ("echo: " + 194 a's = 200 bytes = 50*4).
		body, err := json.Marshal(map[string]any{
			"model": "mock-small", "max_tokens": 50,
			"messages": []map[string]any{{"role": "user", "content": strings.Repeat("a", 1000)}},
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
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatal(err)
		}
		want := "echo: " + strings.Repeat("a", 194)
		if len(out.Choices) != 1 || out.Choices[0].Message.Content != want || out.Choices[0].FinishReason != "length" {
			t.Errorf("ascii control: content = %d bytes (want %d), finish_reason = %q (want length) - the kept emulation must still truncate",
				len(out.Choices[0].Message.Content), len(want), out.Choices[0].FinishReason)
		}
	})
}
