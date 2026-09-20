package main

import (
	"encoding/json"
	"net/http"
)

// Control endpoints (/__control/*) and /v1/models. Control is intentionally
// tiny JSON: fault queues, per-chunk latency, stream truncation, reset.

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

// POST /__control/fail_next
// body: {"status":429, "times":2, "retry_after":"1", "body":"..."}
// times defaults to 1; body overrides the synthesized error payload.
func (s *Server) handleFailNext(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Status     int    `json:"status"`
		Times      int    `json:"times"`
		RetryAfter string `json:"retry_after"`
		Body       string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "fail_next body: "+err.Error())
		return
	}
	if req.Status < 100 || req.Status > 599 {
		writeJSONError(w, http.StatusBadRequest, "fail_next.status out of range")
		return
	}
	if req.Times <= 0 {
		req.Times = 1
	}
	s.mu.Lock()
	for i := 0; i < req.Times; i++ {
		s.faults = append(s.faults, fault{status: req.Status, retryAfter: req.RetryAfter, body: req.Body})
	}
	queued := len(s.faults)
	s.mu.Unlock()
	writeJSON(w, map[string]any{"queued": queued, "status": req.Status, "times": req.Times})
}

// POST /__control/latency {"ms": 200}
func (s *Server) handleLatency(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MS int `json:"ms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MS < 0 {
		writeJSONError(w, http.StatusBadRequest, "latency body must be {\"ms\":>=0}")
		return
	}
	s.mu.Lock()
	s.latency = req.MS
	s.mu.Unlock()
	writeJSON(w, map[string]any{"latency_ms": req.MS})
}

// POST /__control/truncate {"chunks": 3}
func (s *Server) handleTruncate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Chunks int `json:"chunks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Chunks < 0 {
		writeJSONError(w, http.StatusBadRequest, "truncate body must be {\"chunks\":>=0}")
		return
	}
	s.mu.Lock()
	s.truncate = req.Chunks
	s.mu.Unlock()
	writeJSON(w, map[string]any{"truncate_after_chunks": req.Chunks})
}

// POST /__control/reset clears every injection and counter.
func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.faults = nil
	s.latency = 0
	s.truncate = 0
	s.mu.Unlock()
	s.resetGolden()
	writeJSON(w, map[string]any{"reset": true})
}

// GET /__control/last_request?route=messages|responses|chat
// Answers the captured body of the most recent request on that route
// (ticket 11: cache-breakpoint and probe assertions read the wire the
// adapter actually sent). {present:false} after a reset.
func (s *Server) handleLastRequest(w http.ResponseWriter, r *http.Request) {
	route := r.URL.Query().Get("route")
	if route == "" {
		writeJSONError(w, http.StatusBadRequest, "last_request needs a ?route= parameter")
		return
	}
	rec, ok := s.lastRequestBody(route)
	if !ok {
		writeJSON(w, map[string]any{"route": route, "present": false})
		return
	}
	writeJSON(w, map[string]any{
		"route": route, "present": true,
		"body":   rec.Body,
		"header": rec.Header,
	})
}

// GET /__control/state
func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	lat, trunc, queued := s.snapshot()
	total, routes := s.stats()
	writeJSON(w, map[string]any{
		"requests_total":        total,
		"routes":                routes,
		"queued_faults":         queued,
		"latency_ms":            lat,
		"truncate_after_chunks": trunc,
	})
}

// GET /v1/models - the auto-discovery target.
func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	data := make([]map[string]string, 0, len(s.models))
	for _, m := range s.models {
		data = append(data, map[string]string{
			"id": m, "object": "model", "owned_by": "mockllm",
		})
	}
	writeJSON(w, map[string]any{
		"object": "list",
		"data":   data,
	})
}
