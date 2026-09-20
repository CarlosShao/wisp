package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// POST /v1/messages - the Anthropic Messages protocol shape (ticket 11 owns
// this endpoint's dialect fidelity: flat tool definitions, block content,
// cache_control markers, thinking blocks). Synthesis mirrors the chat
// endpoint with message_start + content blocks + message_delta usage
// semantics (cumulative output tokens).

// msgTool is an Anthropic tool declaration: flat name/input_schema (NOT the
// chat {function:{name}} nesting).
type msgTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

// msgBlock is one Anthropic content block.
type msgBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text"`
	CacheCtl  json.RawMessage `json:"cache_control"`
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
	ToolUseID string          `json:"tool_use_id"`
	Content   json.RawMessage `json:"content"`
	IsError   bool            `json:"is_error"`
	Source    struct {
		Type      string `json:"type"`
		MediaType string `json:"media_type"`
		Data      string `json:"data"`
	} `json:"source"`
}

type msgMessage struct {
	Role    string     `json:"role"`
	Content []msgBlock `json:"content"`
	Raw     json.RawMessage
}

// Unmarshal accepts both the string and the block-array content forms.
func (m *msgMessage) UnmarshalJSON(b []byte) error {
	type alias struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	m.Role, m.Raw = a.Role, a.Content
	var s string
	if err := json.Unmarshal(a.Content, &s); err == nil {
		m.Content = []msgBlock{{Type: "text", Text: s}}
		return nil
	}
	return json.Unmarshal(a.Content, &m.Content)
}

// msgPlan is the deterministic answer derived from a Messages request.
type msgPlan struct {
	answer   string
	reason   string
	toolName string
	toolArgs string
	// image records that the last user turn carried an image block: the
	// capability emulation reads it to answer like a non-vision model.
	image bool
}

func planMessagesFrom(messages []msgMessage, tools []msgTool, choice json.RawMessage) msgPlan {
	var p msgPlan
	text, image := "", false
	for _, m := range messages {
		if m.Role != "user" {
			continue
		}
		text, image = "", false
		for _, b := range m.Content {
			switch b.Type {
			case "text":
				text += b.Text
			case "image":
				image = true
			case "tool_result":
				var inner []msgBlock
				if err := json.Unmarshal(b.Content, &inner); err == nil {
					for _, ib := range inner {
						text += ib.Text
					}
				}
			}
		}
	}
	switch {
	case image:
		p.answer = "vision-ok"
	default:
		p.answer = "echo: " + text
		if len(p.answer) > 400 {
			p.answer = p.answer[:400]
		}
	}
	p.image = image
	if strings.Contains(text, "[think]") {
		p.reason = "thinking about: " + text
		if len(p.reason) > 200 {
			p.reason = p.reason[:200]
		}
	}

	// Forced tool selection (tool_choice any | tool) with tools declared.
	if len(tools) > 0 && len(choice) > 0 {
		var bare string
		if json.Unmarshal(choice, &bare) == nil && (bare == "any" || bare == "required") {
			p.toolName = tools[0].Name
		}
		var obj struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}
		if json.Unmarshal(choice, &obj) == nil {
			switch obj.Type {
			case "tool":
				p.toolName = obj.Name
			case "any":
				// Anthropic spells the forced choice as the OBJECT form
				// {"type":"any"}, which is exactly what the C5 adapter sends
				// for tool_choice=required. A real model has to answer that
				// with a tool_use block, so the mock has to as well - or an
				// fc probe on this dialect could never measure "capable".
				if p.toolName == "" {
					p.toolName = tools[0].Name
				}
			}
		}
		if p.toolName == "" &&
			(len(choice) == 0 || strings.Contains(string(choice), `"auto"`)) {
			// auto: answer in text (the real protocol's default too).
		}
		if p.toolName != "" {
			argText := text
			if argText == "" {
				argText = "ping"
			}
			p.toolArgs = fmt.Sprintf(`{"text":%s}`, mustJSONString(argText))
		}
	}
	return p
}

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	// The body is read and recorded BEFORE any response byte is produced
	// (ticket 11): a faulted or golden-served request must not leave the
	// previous /__control/last_request capture in place.
	s.recordRequest("messages", body, r.Header)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "read body: "+err.Error())
		return
	}
	if f, ok := s.takeFault(); ok {
		serveFault(w, f, "messages")
		return
	}
	var req struct {
		Model     string          `json:"model"`
		System    json.RawMessage `json:"system"`
		Messages  []msgMessage    `json:"messages"`
		Tools     []msgTool       `json:"tools"`
		Choice    json.RawMessage `json:"tool_choice"`
		MaxTokens int             `json:"max_tokens"`
		Stream    bool            `json:"stream"`
		Thinking  json.RawMessage `json:"thinking"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "messages body is not valid JSON: "+err.Error())
		return
	}
	if name := goldenName(req.Model, r.Header.Get("X-Wisp-Golden")); name != "" {
		s.serveGoldenGeneric(w, name, "messages")
		return
	}

	p := planMessagesFrom(req.Messages, req.Tools, req.Choice)
	// Capability emulation (ticket 11 AC#6), Messages dialect: the Anthropic
	// shape of the same two honest failures a real provider produces. It runs
	// AFTER the golden branch, so byte-pinned replay can never be disturbed.
	caps := s.capability()
	if caps.thinkingBroken() {
		p.reason = ""
	}
	if caps.visionBroken() && p.image {
		writeJSONError(w, http.StatusBadRequest,
			"mockllm capability mode: vision=broken rejects image input")
		return
	}
	if caps.fcBroken() && p.toolName != "" {
		// A model without tool use silently ignores a forced tool_choice:
		// same text answer, no tool_use block, stop_reason end_turn.
		p.toolName, p.toolArgs = "", ""
	}
	promptTokens := 12 + len(string(body))/200
	completionTokens := len(p.answer)/4 + len(p.toolArgs)/4 + 1
	msgID := s.reqID("msg")
	lat, trunc, _ := s.snapshot()

	if !req.Stream {
		content := []map[string]any{}
		if p.reason != "" {
			content = append(content, map[string]any{"type": "thinking", "thinking": p.reason})
		}
		if p.toolName != "" {
			content = append(content, map[string]any{
				"type": "tool_use", "id": "toolu_mock_1", "name": p.toolName,
				"input": json.RawMessage(p.toolArgs),
			})
		} else {
			content = append(content, map[string]any{"type": "text", "text": p.answer})
		}
		stop := "end_turn"
		if p.toolName != "" {
			stop = "tool_use"
		}
		writeJSON(w, map[string]any{
			"id": msgID, "type": "message", "role": "assistant", "content": content,
			"model": req.Model, "stop_reason": stop,
			"usage": map[string]int{"input_tokens": promptTokens, "output_tokens": completionTokens},
		})
		return
	}

	// Stream (Anthropic SSE dialect: event: + data: per frame).
	sse := newSSEWriter(w, lat, trunc)
	raw := func(ev string, payload any) bool {
		// sse.event applies /__control/truncate accounting (sse.raw does not,
		// so a raw-only writer would make the named dialects uncuttable).
		return sse.event(ev, payload)
	}
	idx := 0
	if !raw("message_start", map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id": msgID, "type": "message", "role": "assistant", "model": req.Model,
			"content": []any{},
			"usage": map[string]int{"input_tokens": promptTokens,
				"cache_read_input_tokens": 0, "output_tokens": 0},
		},
	}) {
		return
	}
	if p.reason != "" {
		if !raw("content_block_start", map[string]any{"type": "content_block_start",
			"index": idx, "content_block": map[string]any{"type": "thinking", "thinking": ""}}) {
			return
		}
		if !raw("content_block_delta", map[string]any{"type": "content_block_delta",
			"index": idx, "delta": map[string]any{"type": "thinking_delta", "thinking": p.reason}}) {
			return
		}
		if !raw("content_block_stop", map[string]any{"type": "content_block_stop", "index": idx}) {
			return
		}
		idx++
	}
	stop := "end_turn"
	switch {
	case p.toolName != "":
		stop = "tool_use"
		if !raw("content_block_start", map[string]any{"type": "content_block_start",
			"index": idx, "content_block": map[string]any{
				"type": "tool_use", "id": "toolu_mock_1", "name": p.toolName,
				"input": map[string]any{}}}) {
			return
		}
		if !raw("content_block_delta", map[string]any{"type": "content_block_delta",
			"index": idx, "delta": map[string]any{"type": "input_json_delta",
				"partial_json": p.toolArgs}}) {
			return
		}
		if !raw("content_block_stop", map[string]any{"type": "content_block_stop", "index": idx}) {
			return
		}
	default:
		if !raw("content_block_start", map[string]any{"type": "content_block_start",
			"index": idx, "content_block": map[string]any{"type": "text", "text": ""}}) {
			return
		}
		for _, piece := range splitChunks(p.answer, 16) {
			if !raw("content_block_delta", map[string]any{"type": "content_block_delta",
				"index": idx, "delta": map[string]any{"type": "text_delta", "text": piece}}) {
				return
			}
		}
		if !raw("content_block_stop", map[string]any{"type": "content_block_stop", "index": idx}) {
			return
		}
	}
	if !raw("message_delta", map[string]any{
		"type":  "message_delta",
		"delta": map[string]any{"stop_reason": stop, "stop_sequence": nil},
		"usage": map[string]int{"output_tokens": completionTokens},
	}) {
		return
	}
	raw("message_stop", map[string]any{"type": "message_stop"})
}
