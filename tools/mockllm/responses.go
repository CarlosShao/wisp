package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// POST /v1/responses - OpenAI Responses API, framework level (the real
// adapter is ticket 11; mockllm only needs a stable shape for tests).

func (s *Server) handleResponses(w http.ResponseWriter, r *http.Request) {
	if f, ok := s.takeFault(); ok {
		serveFault(w, f, "responses")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "read body: "+err.Error())
		return
	}
	var req struct {
		Model  string          `json:"model"`
		Input  json.RawMessage `json:"input"`
		Stream bool            `json:"stream"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "responses body is not valid JSON: "+err.Error())
		return
	}

	if name := goldenName(req.Model, r.Header.Get("X-Wisp-Golden")); name != "" {
		s.serveGoldenGeneric(w, name, "responses")
		return
	}

	text := "echo: " + responsesInputText(req.Input)
	lat, trunc, _ := s.snapshot()
	respID := s.reqID("resp")
	if !req.Stream {
		writeJSON(w, map[string]any{
			"id": respID, "object": "response", "model": req.Model,
			"output": []map[string]any{{
				"type": "message", "role": "assistant",
				"content": []map[string]string{{"type": "output_text", "text": text}},
			}},
			"usage": map[string]int{"input_tokens": 12, "output_tokens": len(text) / 4},
		})
		return
	}
	sse := newSSEWriter(w, lat, trunc)
	if !sse.chunk(map[string]any{"type": "response.created", "response": map[string]any{"id": respID}}) {
		return
	}
	for _, piece := range splitChunks(text, 16) {
		if !sse.chunk(map[string]any{"type": "response.output_text.delta", "delta": piece}) {
			return
		}
	}
	if !sse.chunk(map[string]any{
		"type": "response.completed",
		"response": map[string]any{
			"id":    respID,
			"usage": map[string]int{"input_tokens": 12, "output_tokens": len(text) / 4},
		},
	}) {
		return
	}
	sse.done()
}

// responsesInputText extracts a text echo source from the flexible Responses
// input (string, or array of {role,content} items).
func responsesInputText(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var items []struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(raw, &items); err == nil {
		out := ""
		for _, it := range items {
			c := decodeChatContent(it.Content)
			if c.text != "" {
				out = c.text // last wins, mirrors chat echo
			}
		}
		return out
	}
	return ""
}

// serveGoldenGeneric replays a golden file for the non-chat endpoints.
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
	lat, trunc, _ := s.snapshot()
	_ = lat // golden latency directive wins below
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
