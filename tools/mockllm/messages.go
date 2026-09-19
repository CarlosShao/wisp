package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// POST /v1/messages - Anthropic Messages protocol, framework level (the real
// adapter is ticket 11). Synthesis mirrors the chat endpoint: echo / vision
// / audio / tool use, with message_start + content blocks + message_delta
// usage semantics (cumulative output tokens).

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	if f, ok := s.takeFault(); ok {
		serveFault(w, f, "messages")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "read body: "+err.Error())
		return
	}
	var req struct {
		Model      string          `json:"model"`
		System     json.RawMessage `json:"system"`
		Messages   []chatMessage   `json:"messages"`
		Tools      []chatTool      `json:"tools"`
		ToolChoice json.RawMessage `json:"tool_choice"`
		MaxTokens  int             `json:"max_tokens"`
		Stream     bool            `json:"stream"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "messages body is not valid JSON: "+err.Error())
		return
	}

	if name := goldenName(req.Model, r.Header.Get("X-Wisp-Golden")); name != "" {
		s.serveGoldenGeneric(w, name, "messages")
		return
	}

	answer := "echo: " + lastUserText(req.Messages)
	if len(answer) > 400 {
		answer = answer[:400]
	}
	toolName, toolArgs := "", ""
	if len(req.Tools) > 0 && len(req.ToolChoice) > 0 {
		var tc string
		if json.Unmarshal(req.ToolChoice, &tc) == nil && (tc == "any" || tc == "required") {
			toolName = req.Tools[0].Function.Name
		}
		var obj struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}
		if json.Unmarshal(req.ToolChoice, &obj) == nil && obj.Type == "tool" {
			toolName = obj.Name
		}
		if toolName != "" {
			toolArgs = fmt.Sprintf(`{"text":%s}`, mustJSONString(lastUserText(req.Messages)))
		}
	}
	promptTokens, completionTokens := 12, len(answer)/4+len(toolArgs)/4
	if completionTokens == 0 {
		completionTokens = 1
	}
	msgID := s.reqID("msg")
	lat, trunc, _ := s.snapshot()

	if !req.Stream {
		content := []map[string]any{}
		if toolName != "" {
			content = append(content, map[string]any{
				"type": "tool_use", "id": "toolu_mock_1", "name": toolName,
				"input": json.RawMessage(toolArgs),
			})
		} else {
			content = append(content, map[string]any{"type": "text", "text": answer})
		}
		stop := "end_turn"
		if toolName != "" {
			stop = "tool_use"
		}
		writeJSON(w, map[string]any{
			"id": msgID, "type": "message", "role": "assistant", "content": content,
			"model": req.Model, "stop_reason": stop,
			"usage": map[string]int{"input_tokens": promptTokens, "output_tokens": completionTokens},
		})
		return
	}

	// Stream (anthropic SSE dialect); newSSEWriter writes the status header.
	sse := newSSEWriter(w, lat, trunc)
	if !sse.raw("event: message_start\n") {
		return
	}
	if !sse.chunk(map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id": msgID, "type": "message", "role": "assistant",
			"usage": map[string]int{"input_tokens": promptTokens, "output_tokens": 0},
		},
	}) {
		return
	}
	if toolName != "" {
		if !sse.chunk(map[string]any{
			"type": "content_block_start", "index": 0,
			"content_block": map[string]any{
				"type": "tool_use", "id": "toolu_mock_1", "name": toolName, "input": map[string]any{},
			},
		}) {
			return
		}
		if !sse.chunk(map[string]any{
			"type": "content_block_delta", "index": 0,
			"delta": map[string]any{"type": "input_json_delta", "partial_json": toolArgs},
		}) {
			return
		}
		if !sse.chunk(map[string]any{"type": "content_block_stop", "index": 0}) {
			return
		}
	} else {
		if !sse.chunk(map[string]any{
			"type": "content_block_start", "index": 0,
			"content_block": map[string]any{"type": "text", "text": ""},
		}) {
			return
		}
		for _, piece := range splitChunks(answer, 16) {
			if !sse.chunk(map[string]any{
				"type": "content_block_delta", "index": 0,
				"delta": map[string]any{"type": "text_delta", "text": piece},
			}) {
				return
			}
		}
		if !sse.chunk(map[string]any{"type": "content_block_stop", "index": 0}) {
			return
		}
	}
	stop := "end_turn"
	if toolName != "" {
		stop = "tool_use"
	}
	if !sse.chunk(map[string]any{
		"type":  "message_delta",
		"delta": map[string]any{"stop_reason": stop},
		"usage": map[string]int{"output_tokens": completionTokens}, // cumulative
	}) {
		return
	}
	if !sse.chunk(map[string]any{"type": "message_stop"}) {
		return
	}
}

// lastUserText returns the last user message's text (string or parts).
func lastUserText(messages []chatMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return decodeChatContent(messages[i].Content).text
		}
	}
	return ""
}
