package main

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

// Ticket 391: the Messages and Responses dialects parsed max_tokens /
// max_output_tokens but never applied them - the chat dialect's max_tokens
// truncation emulation (ticket 386) had no twin on either side, so the same
// request parameter silently did nothing on two of the three dialects. The
// fix mirrors chat's shape: the ~4-bytes-per-token cap is intentional
// emulation and stays, the cut lands on a rune boundary, and the termination
// uses each protocol's OWN spelling - Messages stop_reason "max_tokens",
// Responses status "incomplete" with incomplete_details.reason
// "max_output_tokens". Neither protocol has a "length" word (that is the
// chat dialect's spelling); the C5 adapters map all three onto the same
// normalized max_tokens stop.

func Test391MessagesMaxTokensTruncates(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	// Same arithmetic as the chat twin (ticket 386): 300 hanzi = 900 bytes;
	// "echo: " + sample = 906 bytes, read as 226 heuristic tokens > 100, so
	// the cap must truncate. Within the 400-byte bound the longest whole-rune
	// prefix is "echo: " + 131 hanzi = 399 bytes.
	const maxTokens = 100
	want := "echo: " + strings.Repeat("测", 131)
	input := strings.Repeat("测", 300)

	t.Run("nonstream_stop_reason_is_protocol_word", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_tokens": maxTokens,
			"messages": []map[string]any{{"role": "user", "content": input}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/messages", string(body))
		var msg struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
			StopReason string `json:"stop_reason"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			t.Fatal(err)
		}
		if msg.StopReason != "max_tokens" {
			t.Errorf("stop_reason = %q, want %q: the Messages dialect reports the protocol's own stop_reason spelling, not the chat word %q",
				msg.StopReason, "max_tokens", "length")
		}
		got := ""
		for _, b := range msg.Content {
			if b.Type == "text" {
				got = b.Text
			}
		}
		if !utf8.ValidString(got) {
			t.Errorf("truncated text is not valid UTF-8: the cut landed mid-rune (tail %q)", tailRunes(got, 20))
		}
		if got != want {
			t.Errorf("text = %d bytes / %d runes, want the rune-boundary prefix %d bytes / %d runes (tail %q)",
				len(got), len([]rune(got)), len(want), len([]rune(want)), tailRunes(got, 20))
		}
		if len(got) > maxTokens*4 {
			t.Errorf("text = %d bytes, exceeds the max_tokens*4 heuristic bound %d", len(got), maxTokens*4)
		}
	})

	t.Run("stream_message_delta_stop_reason", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_tokens": maxTokens, "stream": true,
			"messages": []map[string]any{{"role": "user", "content": input}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/messages", string(body))
		text, stop := "", ""
		for _, c := range sseChunks(t, readAll(t, resp.Body)) {
			var m struct {
				Type  string `json:"type"`
				Delta struct {
					Text       string `json:"text"`
					StopReason string `json:"stop_reason"`
				} `json:"delta"`
			}
			if err := json.Unmarshal([]byte(chunkData(t, c)), &m); err != nil {
				t.Fatalf("frame not JSON: %v", err)
			}
			switch m.Type {
			case "content_block_delta":
				text += m.Delta.Text
			case "message_delta":
				stop = m.Delta.StopReason
			}
		}
		if stop != "max_tokens" {
			t.Errorf("stream stop_reason = %q, want %q", stop, "max_tokens")
		}
		if !utf8.ValidString(text) {
			t.Errorf("assembled text is not valid UTF-8: the cut landed mid-rune (tail %q)", tailRunes(text, 20))
		}
		if text != want {
			t.Errorf("assembled text = %d bytes / %d runes, want the rune-boundary prefix %d bytes / %d runes (tail %q)",
				len(text), len([]rune(text)), len(want), len([]rune(want)), tailRunes(text, 20))
		}
	})

	t.Run("ascii_cut_byte_identical_to_cap", func(t *testing.T) {
		// Control for the kept cap: for ASCII the rune-boundary cut is the
		// byte cut, so the emulation still truncates exactly at
		// max_tokens*4 bytes ("echo: " + 194 a's = 200 bytes = 50*4).
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_tokens": 50,
			"messages": []map[string]any{{"role": "user", "content": strings.Repeat("a", 1000)}},
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/messages", string(body))
		var msg struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
			StopReason string `json:"stop_reason"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			t.Fatal(err)
		}
		got := ""
		for _, b := range msg.Content {
			if b.Type == "text" {
				got = b.Text
			}
		}
		want := "echo: " + strings.Repeat("a", 194)
		if got != want || msg.StopReason != "max_tokens" {
			t.Errorf("ascii control: text = %d bytes (want %d), stop_reason = %q (want max_tokens) - the kept emulation must still truncate",
				len(got), len(want), msg.StopReason)
		}
	})
}

func Test391ResponsesMaxTokensTruncates(t *testing.T) {
	ts := newTestServer(t, t.TempDir())
	const maxOutputTokens = 100
	want := "echo: " + strings.Repeat("测", 131)
	input := strings.Repeat("测", 300)

	t.Run("nonstream_status_incomplete", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_output_tokens": maxOutputTokens,
			"input": input,
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/responses", string(body))
		var rs struct {
			Status            string `json:"status"`
			IncompleteDetails *struct {
				Reason string `json:"reason"`
			} `json:"incomplete_details"`
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
		if rs.Status != "incomplete" {
			t.Errorf("status = %q, want %q: hitting max_output_tokens is an incomplete response, not a completed one",
				rs.Status, "incomplete")
		}
		if rs.IncompleteDetails == nil || rs.IncompleteDetails.Reason != "max_output_tokens" {
			t.Errorf("incomplete_details = %+v, want reason %q", rs.IncompleteDetails, "max_output_tokens")
		}
		got := ""
		for _, item := range rs.Output {
			if item.Type == "message" {
				for _, c := range item.Content {
					if c.Type == "output_text" {
						got = c.Text
					}
				}
			}
		}
		if !utf8.ValidString(got) {
			t.Errorf("truncated text is not valid UTF-8: the cut landed mid-rune (tail %q)", tailRunes(got, 20))
		}
		if got != want {
			t.Errorf("text = %d bytes / %d runes, want the rune-boundary prefix %d bytes / %d runes (tail %q)",
				len(got), len([]rune(got)), len(want), len([]rune(want)), tailRunes(got, 20))
		}
		if len(got) > maxOutputTokens*4 {
			t.Errorf("text = %d bytes, exceeds the max_output_tokens*4 heuristic bound %d", len(got), maxOutputTokens*4)
		}
	})

	t.Run("stream_terminal_event_is_response_incomplete", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_output_tokens": maxOutputTokens, "stream": true,
			"input": input,
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/responses", string(body))
		text, terminal, status, reason := "", "", "", ""
		for _, c := range sseChunks(t, readAll(t, resp.Body)) {
			if data := chunkData(t, c); data == "[DONE]" {
				continue
			}
			var ev struct {
				Type     string `json:"type"`
				Response *struct {
					Status            string `json:"status"`
					IncompleteDetails *struct {
						Reason string `json:"reason"`
					} `json:"incomplete_details"`
				} `json:"response"`
				Delta string `json:"delta"`
			}
			if err := json.Unmarshal([]byte(chunkData(t, c)), &ev); err != nil {
				t.Fatalf("frame not JSON: %v", err)
			}
			switch ev.Type {
			case "response.output_text.delta":
				text += ev.Delta
			case "response.completed", "response.incomplete", "response.failed":
				terminal = ev.Type
				if ev.Response != nil {
					status = ev.Response.Status
					if ev.Response.IncompleteDetails != nil {
						reason = ev.Response.IncompleteDetails.Reason
					}
				}
			}
		}
		if terminal != "response.incomplete" {
			t.Errorf("terminal event = %q, want %q: the truncated stream must end on the incomplete terminal, not the completed one",
				terminal, "response.incomplete")
		}
		if status != "incomplete" || reason != "max_output_tokens" {
			t.Errorf("terminal status = %q reason = %q, want incomplete / max_output_tokens", status, reason)
		}
		if !utf8.ValidString(text) {
			t.Errorf("assembled text is not valid UTF-8: the cut landed mid-rune (tail %q)", tailRunes(text, 20))
		}
		if text != want {
			t.Errorf("assembled text = %d bytes / %d runes, want the rune-boundary prefix %d bytes / %d runes (tail %q)",
				len(text), len([]rune(text)), len(want), len([]rune(want)), tailRunes(text, 20))
		}
	})

	t.Run("ascii_cut_byte_identical_to_cap", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"model": "mock-large", "max_output_tokens": 50,
			"input": strings.Repeat("a", 1000),
		})
		if err != nil {
			t.Fatal(err)
		}
		resp := postJSON(t, ts.URL+"/v1/responses", string(body))
		var rs struct {
			Status string `json:"status"`
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
			if item.Type == "message" {
				for _, c := range item.Content {
					if c.Type == "output_text" {
						got = c.Text
					}
				}
			}
		}
		want := "echo: " + strings.Repeat("a", 194)
		if got != want || rs.Status != "incomplete" {
			t.Errorf("ascii control: text = %d bytes (want %d), status = %q (want incomplete) - the kept emulation must still truncate",
				len(got), len(want), rs.Status)
		}
	})
}
