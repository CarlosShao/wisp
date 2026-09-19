package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// POST /v1/chat/completions (OpenAI Chat Completions shape).
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if f, ok := s.takeFault(); ok {
		serveFault(w, f, "chat")
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "read body: "+err.Error())
		return
	}
	var req chatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "chat body is not valid JSON: "+err.Error())
		return
	}

	// Golden mode?
	if name := goldenName(req.Model, r.Header.Get("X-Wisp-Golden")); name != "" {
		s.serveGoldenChat(w, r, name)
		return
	}
	lat, trunc, _ := s.snapshot()
	s.synthChat(w, &req, lat, trunc)
}

type chatRequest struct {
	Model       string          `json:"model"`
	Messages    []chatMessage   `json:"messages"`
	Tools       []chatTool      `json:"tools"`
	ToolChoice  json.RawMessage `json:"tool_choice"`
	Stream      bool            `json:"stream"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature any             `json:"temperature"`
	Stop        any             `json:"stop"`
}

type chatMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
	// assistant tool_calls
	ToolCalls []chatToolCall `json:"tool_calls,omitempty"`
}

type chatToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type chatTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Parameters  json.RawMessage `json:"parameters"`
	} `json:"function"`
}

// chatContent decodes the two content spellings (string or parts array).
type chatContent struct {
	text    string
	image   bool
	audio   bool
	thought bool // "[think]" marker requests reasoning_content output
}

func decodeChatContent(raw json.RawMessage) chatContent {
	var c chatContent
	if len(raw) == 0 {
		return c
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		c.text = s
		c.thought = strings.Contains(s, "[think]")
		return c
	}
	var parts []struct {
		Type     string `json:"type"`
		Text     string `json:"text"`
		ImageURL struct {
			URL string `json:"url"`
		} `json:"image_url"`
		InputAudio struct {
			Data   string `json:"data"`
			Format string `json:"format"`
		} `json:"input_audio"`
	}
	if err := json.Unmarshal(raw, &parts); err != nil {
		return c
	}
	for _, p := range parts {
		switch p.Type {
		case "text":
			c.text += p.Text
			if strings.Contains(p.Text, "[think]") {
				c.thought = true
			}
		case "image_url":
			c.image = true
		case "input_audio":
			c.audio = true
		}
	}
	return c
}

// planChat derives the deterministic response from the request.
func planChat(req *chatRequest) (answer string, toolName string, toolArgs string, reasoning string, finish string) {
	lastUser := chatContent{}
	for _, m := range req.Messages {
		if m.Role == "user" {
			lastUser = decodeChatContent(m.Content)
		}
	}
	switch {
	case lastUser.image:
		answer = "vision-ok"
	case lastUser.audio:
		answer = "audio-ok"
	default:
		answer = "echo: " + lastUser.text
		if len(answer) > 400 {
			answer = answer[:400]
		}
	}
	if lastUser.thought {
		reasoning = "thinking about: " + lastUser.text
		if len(reasoning) > 200 {
			reasoning = reasoning[:200]
		}
	}

	// Tool call decision: forced tool / required / auto-with-[tool] marker.
	if len(req.Tools) > 0 && req.ToolChoice != nil {
		var tc string
		if err := json.Unmarshal(req.ToolChoice, &tc); err == nil {
			switch tc {
			case "required":
				toolName = req.Tools[0].Function.Name
			case "auto":
				if strings.Contains(lastUser.text, "[tool]") {
					toolName = req.Tools[0].Function.Name
				}
			}
			return // bare string handled
		}
		var obj struct {
			Type     string `json:"type"`
			Function struct {
				Name string `json:"name"`
			} `json:"function"`
		}
		if err := json.Unmarshal(req.ToolChoice, &obj); err == nil && obj.Type == "function" {
			toolName = obj.Function.Name
		}
	}
	if toolName != "" {
		toolArgs = fmt.Sprintf(`{"text":%s}`, mustJSONString(lastUser.text))
	}

	// max_tokens truncation: ~4 chars per token heuristic.
	if req.MaxTokens > 0 && toolName == "" {
		if len(answer)/4 > req.MaxTokens {
			answer = answer[:req.MaxTokens*4]
			finish = "length"
		}
	}
	if finish == "" {
		if toolName != "" {
			finish = "tool_calls"
		} else {
			finish = "stop"
		}
	}
	return answer, toolName, toolArgs, reasoning, finish
}

func estimateTokens(b []byte) int { return len(b) / 4 }

func (s *Server) synthChat(w http.ResponseWriter, req *chatRequest, lat int, trunc int) {
	answer, toolName, toolArgs, reasoning, finish := planChat(req)
	reqID := s.reqID("chatcmpl")
	created := time.Now().UTC().Unix()

	if !req.Stream {
		msg := map[string]any{"role": "assistant", "content": answer}
		if toolName != "" {
			msg["content"] = nil
			msg["tool_calls"] = []map[string]any{{
				"id": "call_mock_1", "type": "function",
				"function": map[string]string{"name": toolName, "arguments": toolArgs},
			}}
		}
		if reasoning != "" {
			msg["reasoning_content"] = reasoning
		}
		writeJSON(w, map[string]any{
			"id": reqID, "object": "chat.completion", "created": created, "model": req.Model,
			"choices": []map[string]any{{
				"index": 0, "message": msg, "finish_reason": finish,
			}},
			"usage": chatUsage(answer, toolArgs, req),
		})
		return
	}

	sse := newSSEWriter(w, lat, trunc)
	// 1. role chunk
	if !sse.chunk(chunkBase(reqID, created, req.Model, map[string]any{"role": "assistant", "content": ""})) {
		return
	}
	// 2. reasoning chunks
	if reasoning != "" {
		for _, piece := range splitChunks(reasoning, 24) {
			if !sse.chunk(chunkBase(reqID, created, req.Model, map[string]any{"reasoning_content": piece})) {
				return
			}
		}
	}
	// 3. tool call chunk (single full fragment, index 0)
	if toolName != "" {
		tc := map[string]any{
			"index": 0, "id": "call_mock_1", "type": "function",
			"function": map[string]string{"name": toolName, "arguments": toolArgs},
		}
		if !sse.chunk(chunkBase(reqID, created, req.Model, map[string]any{"tool_calls": []any{tc}})) {
			return
		}
	}
	// 4. text chunks
	for _, piece := range splitChunks(answer, 16) {
		if !sse.chunk(chunkBase(reqID, created, req.Model, map[string]any{"content": piece})) {
			return
		}
	}
	// 5. finish chunk
	if !sse.chunk(chunkWithFinish(reqID, created, req.Model, finish)) {
		return
	}
	// 6. usage chunk (include_usage semantics)
	u := chatUsage(answer, toolArgs, req)
	if !sse.chunk(map[string]any{
		"id": reqID, "object": "chat.completion.chunk", "created": created, "model": req.Model,
		"choices": []any{}, "usage": u,
	}) {
		return
	}
	sse.done()
}

func chunkBase(id string, created int64, model string, delta map[string]any) map[string]any {
	return map[string]any{
		"id": id, "object": "chat.completion.chunk", "created": created, "model": model,
		"choices": []map[string]any{{"index": 0, "delta": delta, "finish_reason": nil}},
		"usage":   nil,
	}
}

func chunkWithFinish(id string, created int64, model, finish string) map[string]any {
	return map[string]any{
		"id": id, "object": "chat.completion.chunk", "created": created, "model": model,
		"choices": []map[string]any{{"index": 0, "delta": map[string]any{}, "finish_reason": finish}},
		"usage":   nil,
	}
}

type usageBlock struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	TotalTokens         int `json:"total_tokens"`
	PromptTokensDetails struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
}

func chatUsage(answer, toolArgs string, req *chatRequest) usageBlock {
	var u usageBlock
	u.PromptTokens = 12
	u.CompletionTokens = len(answer)/4 + len(toolArgs)/4
	if u.CompletionTokens == 0 {
		u.CompletionTokens = 1
	}
	u.TotalTokens = u.PromptTokens + u.CompletionTokens
	return u
}

// splitChunks cuts s into rune-safe pieces of about n runes.
func splitChunks(s string, n int) []string {
	if s == "" {
		return nil
	}
	r := []rune(s)
	var out []string
	for i := 0; i < len(r); i += n {
		end := i + n
		if end > len(r) {
			end = len(r)
		}
		out = append(out, string(r[i:end]))
	}
	return out
}

// serveGoldenChat replays testdata/golden/<name>.sse section-by-section.
func (s *Server) serveGoldenChat(w http.ResponseWriter, r *http.Request, name string) {
	if !validGoldenName(name) {
		writeJSONError(w, http.StatusBadRequest, "invalid golden name")
		return
	}
	path := filepath.Join(s.goldenDir, name+".sse")
	src, err := os.ReadFile(path)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("golden %q not found in %s", name, s.goldenDir))
		return
	}
	sections, err := parseGolden(src)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "golden parse: "+err.Error())
		return
	}
	sec, ok := s.nextGolden(name, sections)
	if !ok {
		writeJSONError(w, http.StatusInternalServerError,
			fmt.Sprintf("golden %q exhausted (%d responses served; use /__control/reset)", name, len(sections)))
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
	// Stream the recorded bytes chunk-wise (pacing + truncation apply).
	_, trunc, _ := s.snapshot()
	sse := newSSEWriter(w, sec.latencyMS, trunc)
	for _, chunk := range sseChunksOf(sec.body) {
		if !sse.raw(chunk) {
			return
		}
	}
}

// validGoldenName rejects path traversal in golden names.
func validGoldenName(name string) bool {
	if name == "" || strings.ContainsAny(name, "/\\") || strings.Contains(name, "..") {
		return false
	}
	return true
}

// sseChunksOf splits a recorded SSE body into raw chunks (one per SSE event
// block) so pacing/truncation behave like a live stream.
func sseChunksOf(body string) []string {
	var out []string
	var cur strings.Builder
	for _, line := range strings.SplitAfter(body, "\n") {
		if line == "" {
			continue
		}
		cur.WriteString(line)
		if line == "\n" {
			out = append(out, cur.String())
			cur.Reset()
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}
