// Command llmrecord records golden SSE fixtures from a REAL provider
// (ticket 09 recorder CLI). It POSTs one streaming chat completion and
// writes the raw response bytes as testdata/golden/<name>.sse in the
// "wisp golden sse v1" format (internal/llm/golden), so the unit-test
// replayer and tools/mockllm can replay it verbatim.
//
// STATUS: STUB - live recording is DEFERRED until an LLM API key is
// available (docs/reports/pending-and-issues.md item H2; see also ticket 09
// progress log). The full record path (request -> capture -> format write)
// is implemented and exercised against mockllm today; only the real
// provider endpoint is pending. Without a key the command fails with an
// explicit error and never writes a file.
//
// Usage:
//
//	llmrecord -base-url https://api.openai.com/v1 -model gpt-4o \
//	          -name tool-call -prompt "call the echo tool" \
//	          -out testdata/golden [-key-env WISP_LLM_KEY]
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

func main() {
	baseURL := flag.String("base-url", "https://api.openai.com/v1", "provider base URL (OpenAI-compatible)")
	model := flag.String("model", "", "model id to record")
	prompt := flag.String("prompt", "", "user prompt for the recorded turn")
	name := flag.String("name", "", "golden name (file becomes <out>/<name>.sse)")
	out := flag.String("out", "testdata/golden", "output directory")
	keyEnv := flag.String("key-env", "WISP_LLM_KEY", "environment variable holding the API key")
	scenario := flag.String("scenario", "", "free-text scenario note written into the file")
	timeout := flag.Duration("timeout", 2*time.Minute, "record timeout")
	flag.Parse()

	if err := run(*baseURL, *model, *prompt, *name, *out, *keyEnv, *scenario, *timeout); err != nil {
		fmt.Fprintf(os.Stderr, "llmrecord: %v\n", err)
		os.Exit(1)
	}
}

func run(baseURL, model, prompt, name, outDir, keyEnv, scenario string, timeout time.Duration) error {
	if model == "" || name == "" || prompt == "" {
		return fmt.Errorf("-model, -name and -prompt are all required")
	}

	// H2 (docs/reports/pending-and-issues.md): no key -> explicit failure,
	// never a partial or synthetic recording.
	key := os.Getenv(keyEnv)
	if key == "" {
		return fmt.Errorf("environment variable %s is empty: live golden recording is deferred until an LLM API key is provided (H2); "+
			"the mockllm + replayer harness needs no key and is already in place", keyEnv)
	}

	// Build the wire request via the normalized seam path so the recorded
	// bytes match what the adapter really sends.
	temp := 0.0
	req := &llm.Request{
		Model:       model,
		Temperature: &temp,
		Messages: []llm.Message{{
			Role:    llm.RoleUser,
			Content: []llm.Content{llm.TextPart{Text: prompt}},
		}},
	}
	wire, err := buildRecordPayload(req)
	if err != nil {
		return err
	}

	if timeout <= 0 {
		timeout = 2 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(baseURL, "/")+"/chat/completions", bytes.NewReader(wire))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w (class %s)", err, observe.ClassOfOrInternal(llm.ClassifyTransportError(err)))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return fmt.Errorf("provider returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// Capture the SSE body verbatim (that IS the golden payload).
	var body bytes.Buffer
	if _, err := io.Copy(&body, resp.Body); err != nil {
		return fmt.Errorf("read stream: %w", err)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(outDir, name+".sse")
	var sb strings.Builder
	sb.WriteString("# wisp golden sse v1\n")
	sb.WriteString("# @scenario recorded from " + hostOf(baseURL) + " model=" + model)
	if scenario != "" {
		sb.WriteString(": " + scenario)
	}
	sb.WriteString("\n")
	sb.WriteString("# @response 200\n")
	bodyText := strings.ReplaceAll(body.String(), "\r\n", "\n")
	if !strings.HasSuffix(bodyText, "\n") {
		bodyText += "\n"
	}
	for _, line := range strings.Split(strings.TrimSuffix(bodyText, "\n"), "\n") {
		if strings.HasPrefix(line, "# @") {
			line = "# @" + line // escape directive-shaped body lines
		}
		sb.WriteString(line + "\n")
	}
	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "recorded %s (%d bytes of SSE)\n", path, body.Len())
	return nil
}

// buildRecordPayload mirrors openaichat's streaming request shape (the
// recorder records against OpenAI-compatible endpoints; the stream_options
// flag is what makes providers emit the usage chunk the fixtures rely on).
func buildRecordPayload(req *llm.Request) ([]byte, error) {
	type wireMsg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	messages := make([]wireMsg, 0, len(req.Messages))
	for _, m := range req.Messages {
		text := ""
		for _, p := range m.Content {
			if t, ok := p.(llm.TextPart); ok {
				text += t.Text
			}
		}
		messages = append(messages, wireMsg{Role: string(m.Role), Content: text})
	}
	payload := map[string]any{
		"model":    req.Model,
		"messages": messages,
		"stream":   true,
		"stream_options": map[string]bool{
			"include_usage": true,
		},
	}
	return json.Marshal(payload)
}

func hostOf(u string) string {
	s := strings.TrimPrefix(strings.TrimPrefix(u, "https://"), "http://")
	if i := strings.IndexByte(s, '/'); i >= 0 {
		s = s[:i]
	}
	return s
}
