package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// SSE response writer with pacing (latency) and truncation (fault
// injection). A truncated stream simply stops mid-body: the client sees EOF
// before [DONE], which is exactly the mid-stream disconnect shape.

type sseWriter struct {
	w        http.ResponseWriter
	flusher  http.Flusher
	latency  time.Duration // per-chunk delay
	truncate int           // remaining chunks before cut (0 = unlimited)
	cut      bool          // set when truncation fired
}

func newSSEWriter(w http.ResponseWriter, latencyMS, truncate int) *sseWriter {
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	var f http.Flusher
	if fl, ok := w.(http.Flusher); ok {
		f = fl
		f.Flush()
	}
	return &sseWriter{w: w, flusher: f, latency: time.Duration(latencyMS) * time.Millisecond, truncate: truncate}
}

// chunk writes one "data: <v>\n\n" event. Returns false when the stream was
// truncated or the client went away.
func (s *sseWriter) chunk(v any) bool {
	if s.truncate > 0 {
		if s.truncate == 1 {
			s.cut = true // emit this last chunk, then stop
		}
		s.truncate--
	}
	b, err := json.Marshal(v)
	if err != nil {
		return false
	}
	if _, err := fmt.Fprintf(s.w, "data: %s\n\n", b); err != nil {
		return false
	}
	if s.flusher != nil {
		s.flusher.Flush()
	}
	if s.latency > 0 {
		time.Sleep(s.latency)
	}
	return !s.cut
}

// raw writes one pre-formatted block (for event: prefixed lines).
func (s *sseWriter) raw(block string) bool {
	if _, err := fmt.Fprint(s.w, block); err != nil {
		return false
	}
	if s.flusher != nil {
		s.flusher.Flush()
	}
	if s.latency > 0 {
		time.Sleep(s.latency)
	}
	return true
}

// done writes the [DONE] sentinel (unless truncated).
func (s *sseWriter) done() {
	if s.cut {
		return
	}
	fmt.Fprint(s.w, "data: [DONE]\n\n")
	if s.flusher != nil {
		s.flusher.Flush()
	}
}

// serveFault applies an injected failure (status, retry-after, body).
func serveFault(w http.ResponseWriter, f fault, route string) {
	if f.retryAfter != "" {
		w.Header().Set("Retry-After", f.retryAfter)
	}
	if f.body == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(f.status)
		fmt.Fprintf(w, `{"error":{"message":"mockllm injected failure for %s","type":"mockllm_fault","code":%s}}`,
			route, strconv.Quote("injected_"+strconv.Itoa(f.status)))
		return
	}
	ct := "application/json"
	if f.status == 200 {
		ct = "text/event-stream"
	}
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(f.status)
	fmt.Fprint(w, f.body)
}

// writeJSONError answers a client-side (4xx) problem.
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":{"message":%s,"type":"invalid_request_error","code":"mockllm_bad_request"}}`,
		mustJSONString(msg))
}

func mustJSONString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return `"mockllm error"`
	}
	return string(b)
}
