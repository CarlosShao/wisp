package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// POST /v1/responses - the OpenAI Responses API shape (ticket 11 owns this
// endpoint's dialect fidelity: input items, flat tool definitions, reasoning
// items, function_call items, terminal status).

// respInputItem is one `input` item of a Responses request.
type respInputItem struct {
	Type      string          `json:"type"`
	Role      string          `json:"role"`
	Content   json.RawMessage `json:"content"`
	CallID    string          `json:"call_id"`
	Name      string          `json:"name"`
	Arguments string          `json:"arguments"`
	Output    json.RawMessage `json:"output"`
}

// respTool is a Responses tool declaration: flat name/parameters.
type respTool struct {
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// respPart is one content part of a message item.
type respPart struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	ImageURL string `json:"image_url"`
}

func (p respPart) textish() bool {
	return p.Type == "input_text" || p.Type == "output_text" || p.Type == "text" ||
		p.Type == "refusal"
}

func (s *Server) handleResponses(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	// Read + record before any response byte (ticket 11; see handleMessages).
	s.recordRequest("responses", body, r.Header)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "read body: "+err.Error())
		return
	}
	if f, ok := s.takeFault(); ok {
		serveFault(w, f, "responses")
		return
	}
	var req struct {
		Model        string          `json:"model"`
		Instructions string          `json:"instructions"`
		Input        json.RawMessage `json:"input"`
		Tools        []respTool      `json:"tools"`
		Choice       json.RawMessage `json:"tool_choice"`
		Stream       bool            `json:"stream"`
		Reasoning    json.RawMessage `json:"reasoning"`
		MaxTokens    int             `json:"max_output_tokens"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "responses body is not valid JSON: "+err.Error())
		return
	}
	if name := goldenName(req.Model, r.Header.Get("X-Wisp-Golden")); name != "" {
		s.serveGoldenGeneric(w, name, "responses")
		return
	}

	items := respInputItems(req.Input)
	text, image, thought := lastRespUserText(items)
	// Capability emulation (ticket 11 AC#6), Responses dialect: the same two
	// honest failures a real provider produces, in this protocol's shape.
	// Runs AFTER the golden branch, so byte-pinned replay is untouched.
	caps := s.capability()
	if caps.visionBroken() && image {
		writeJSONError(w, http.StatusBadRequest,
			"mockllm capability mode: vision=broken rejects image input")
		return
	}
	answer := "echo: " + text
	if image {
		answer = "vision-ok"
	}
	if len(answer) > 400 {
		answer = answer[:400]
	}
	reason := ""
	// thinking=broken withholds the reasoning summary item and its delta,
	// exactly as the chat and Messages dialects do above.
	if thought && !s.capability().thinkingBroken() {
		reason = "thinking about: " + text
		if len(reason) > 200 {
			reason = reason[:200]
		}
	}
	toolName, toolArgs := "", ""
	if len(req.Tools) > 0 {
		var bare string
		if json.Unmarshal(req.Choice, &bare) == nil && (bare == "required" || bare == "any") {
			toolName = req.Tools[0].Name
		}
		var obj struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}
		if json.Unmarshal(req.Choice, &obj) == nil && obj.Type == "function" {
			toolName = obj.Name
		}
		if toolName != "" {
			argText := text
			if argText == "" {
				argText = "ping"
			}
			toolArgs = fmt.Sprintf(`{"text":%s}`, mustJSONString(argText))
		}
	}
	if caps.fcBroken() && toolName != "" {
		// No function-calling support: the forced tool_choice is silently
		// ignored and the turn answers as an ordinary message item.
		toolName, toolArgs = "", ""
	}

	inTokens := 12 + len(string(body))/200
	outTokens := len(answer)/4 + len(toolArgs)/4 + 1
	if reason != "" {
		outTokens += len(reason) / 4
	}
	respID := s.reqID("resp")
	lat, trunc, _ := s.snapshot()

	if !req.Stream {
		output := []map[string]any{}
		if reason != "" {
			output = append(output, map[string]any{"type": "reasoning", "id": "rs_mock_1",
				"summary": []map[string]string{{"type": "summary_text", "text": reason}}})
		}
		if toolName != "" {
			output = append(output, map[string]any{"type": "function_call", "id": "fc_mock_1",
				"call_id": "call_mock_1", "name": toolName, "arguments": toolArgs})
		} else {
			output = append(output, map[string]any{"type": "message", "role": "assistant",
				"content": []map[string]string{{"type": "output_text", "text": answer}}})
		}
		writeJSON(w, map[string]any{
			"id": respID, "object": "response", "model": req.Model,
			"status": "completed", "output": output,
			"usage": map[string]any{"input_tokens": inTokens, "output_tokens": outTokens},
		})
		return
	}

	sse := newSSEWriter(w, lat, trunc)
	if !sse.event("response.created", map[string]any{
		"type": "response.created",
		"response": map[string]any{"id": respID, "object": "response",
			"status": "in_progress", "model": req.Model},
	}) {
		return
	}
	outIdx := 0
	if reason != "" {
		if !sse.event("response.output_item.added", map[string]any{
			"type": "response.output_item.added", "output_index": outIdx,
			"item": map[string]any{"id": "rs_mock_1", "type": "reasoning",
				"summary": []any{}},
		}) {
			return
		}
		if !sse.event("response.reasoning_summary_text.delta", map[string]any{
			"type": "response.reasoning_summary_text.delta", "item_id": "rs_mock_1",
			"output_index": outIdx, "summary_index": 0, "delta": reason,
		}) {
			return
		}
		if !sse.event("response.output_item.done", map[string]any{
			"type": "response.output_item.done", "output_index": outIdx,
			"item": map[string]any{"id": "rs_mock_1", "type": "reasoning",
				"summary": []map[string]string{{"type": "summary_text", "text": reason}}},
		}) {
			return
		}
		outIdx++
	}
	if toolName != "" {
		if !sse.event("response.output_item.added", map[string]any{
			"type": "response.output_item.added", "output_index": outIdx,
			"item": map[string]any{"id": "fc_mock_1", "type": "function_call",
				"call_id": "call_mock_1", "name": toolName, "arguments": ""},
		}) {
			return
		}
		if !sse.event("response.function_call_arguments.delta", map[string]any{
			"type": "response.function_call_arguments.delta", "item_id": "fc_mock_1",
			"output_index": outIdx, "delta": toolArgs,
		}) {
			return
		}
		if !sse.event("response.output_item.done", map[string]any{
			"type": "response.output_item.done", "output_index": outIdx,
			"item": map[string]any{"id": "fc_mock_1", "type": "function_call",
				"call_id": "call_mock_1", "name": toolName, "arguments": toolArgs},
		}) {
			return
		}
	} else {
		if !sse.event("response.output_item.added", map[string]any{
			"type": "response.output_item.added", "output_index": outIdx,
			"item": map[string]any{"id": "msg_mock_1", "type": "message",
				"role": "assistant", "content": []any{}},
		}) {
			return
		}
		for _, piece := range splitChunks(answer, 16) {
			if !sse.event("response.output_text.delta", map[string]any{
				"type": "response.output_text.delta", "item_id": "msg_mock_1",
				"output_index": outIdx, "content_index": 0, "delta": piece,
			}) {
				return
			}
		}
		if !sse.event("response.output_item.done", map[string]any{
			"type": "response.output_item.done", "output_index": outIdx,
			"item": map[string]any{"id": "msg_mock_1", "type": "message",
				"role":    "assistant",
				"content": []map[string]string{{"type": "output_text", "text": answer}}},
		}) {
			return
		}
	}
	status := "completed"
	if !sse.event("response.completed", map[string]any{
		"type": "response.completed",
		"response": map[string]any{
			"id": respID, "object": "response", "status": status,
			"usage": map[string]any{"input_tokens": inTokens, "output_tokens": outTokens},
		},
	}) {
		return
	}
	sse.done()
}

// respInputItems decodes the flexible `input` field (string or item array).
func respInputItems(raw json.RawMessage) []respInputItem {
	if len(raw) == 0 {
		return nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return []respInputItem{{Type: "message", Role: "user",
			Content: mustJSONStringAsArray(s)}}
	}
	var one respInputItem
	if err := json.Unmarshal(raw, &one); err == nil && one.Type != "" {
		return []respInputItem{one}
	}
	var items []respInputItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil
	}
	return items
}

// mustJSONStringAsArray wraps plain text as a single input_text part array.
func mustJSONStringAsArray(s string) json.RawMessage {
	b, _ := json.Marshal([]respPart{{Type: "input_text", Text: s}})
	return b
}

// lastRespUserText returns the text (and modality flags) of the last user
// message item, plus whether any item asks for thinking.
func lastRespUserText(items []respInputItem) (text string, image, thought bool) {
	for _, it := range items {
		if it.Type == "function_call_output" || it.Type == "function_call" {
			continue
		}
		if it.Role != "user" {
			continue
		}
		text, image = "", false
		var parts []respPart
		if err := json.Unmarshal(it.Content, &parts); err == nil {
			for _, p := range parts {
				switch {
				case p.textish():
					text += p.Text
				case p.Type == "input_image" || p.ImageURL != "":
					image = true
				}
			}
		} else {
			var s string
			if err := json.Unmarshal(it.Content, &s); err == nil {
				text = s
			}
		}
	}
	for _, it := range items {
		if it.Role == "user" && strings.Contains(text, "[think]") {
			thought = true
		}
	}
	return text, image, thought
}

// serveGoldenGeneric replays a golden file for the non-chat endpoints
// (ticket 09 shape, kept verbatim when this handler was upgraded for
// ticket 11's dialect fidelity).
func (s *Server) serveGoldenGeneric(w http.ResponseWriter, name, route string) {
	sections, err := s.loadGolden(name)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sec, ok := s.nextGolden(name, sections)
	if !ok {
		writeJSONError(w, http.StatusInternalServerError,
			fmt.Sprintf("golden %q exhausted; use /__control/reset", name))
		return
	}
	if sec.retryAfter != "" {
		w.Header().Set("Retry-After", sec.retryAfter)
	}
	if sec.status != 200 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(sec.status)
		fmt.Fprint(w, sec.body)
		return
	}
	_, trunc, _ := s.snapshot()
	sse := newSSEWriter(w, sec.latencyMS, trunc)
	for _, chunk := range sseChunksOf(sec.body) {
		if !sse.raw(chunk) {
			return
		}
	}
}

// loadGolden reads and parses a golden file (shared by all endpoints).
func (s *Server) loadGolden(name string) ([]goldenResponse, error) {
	if !validGoldenName(name) {
		return nil, fmt.Errorf("invalid golden name %q", name)
	}
	path := s.goldenDir + "/" + name + ".sse"
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("golden %q not found in %s", name, s.goldenDir)
	}
	return parseGolden(src)
}
