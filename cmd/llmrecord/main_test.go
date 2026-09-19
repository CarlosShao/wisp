package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm/golden"
)

// The recorder is a stub until an API key exists (H2, docs/reports). These
// tests exercise the implemented record path (request -> capture -> format
// write -> parse) against a local server so only the REAL provider endpoint
// is pending.

func TestRecordWithoutKeyFailsExplicitly(t *testing.T) {
	t.Setenv("WISP_LLM_KEY", "")
	err := run("http://127.0.0.1:1/v1", "m", "p", "n", t.TempDir(), "WISP_LLM_KEY", "", 0)
	if err == nil {
		t.Fatal("expected an explicit no-key failure")
	}
	for _, want := range []string{"WISP_LLM_KEY", "deferred", "H2"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q missing %q", err, want)
		}
	}
}

func TestRecordCapturesSSEAndWritesParsableGolden(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			http.NotFound(w, r)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer sk-record" {
			t.Errorf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(
			"data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"recorded!\"},\"finish_reason\":null}]}\n\n" +
				"data: {\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n" +
				"data: [DONE]\n\n"))
	}))
	defer srv.Close()

	outDir := t.TempDir()
	t.Setenv("WISP_LLM_KEY", "sk-record")
	if err := run(srv.URL+"/v1", "mock-small", "say hi", "record-smoke", outDir, "WISP_LLM_KEY",
		"recorder smoke test", 0); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(outDir, "record-smoke.sse")
	rs, err := golden.LoadFile(path)
	if err != nil {
		t.Fatalf("recorded file does not parse: %v", err)
	}
	if len(rs) != 1 || rs[0].Status != 200 {
		t.Fatalf("parsed = %+v, want one 200 response", rs)
	}
	body := string(rs[0].Body)
	for _, want := range []string{`"content":"recorded!"`, "finish_reason\":\"stop\"", "data: [DONE]"} {
		if !strings.Contains(body, want) {
			t.Errorf("recorded body missing %q:\n%s", want, body)
		}
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
