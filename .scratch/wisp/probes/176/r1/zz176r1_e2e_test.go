package main

// Ticket 176-r1 - the FIRST end-to-end reading through the CLI seam
// (`wisp run` = runTextTask / assembleRuntime + execute), for the pipeline
// ticket 176 exists to close:
//
//	run a background task -> its oversized output spills to a real artifact
//	-> the model then calls task.output -> it gets head + tail + total + the
//	   pointer -> and the model takes THAT pointer to fs.read without our own
//	   C25/R4 gate refusing it (ticket 177's precise exemption, first time on a
//	   production path rather than a hand-built bridge).
//
// Injection surface is the one AGENTS §1.3 lists for this leg: the CLI seam,
// driven by a real OpenAI-compatible SSE endpoint (the C5 golden-SSE shape over
// a live HTTP server). Real loop, real bridge, real risk.Provenance, real
// files, real approval gate. Zero mocks standing in for real components.
//
// Both runs share one assembled runtime on purpose: ticket 164 定案② says the
// roster is process-local, so "the model calls task.output about the task that
// just finished" can only be true inside one process. Run A goes through
// runTextTask (the whole `wisp run` entry, exit code and all); run B is the
// same composition root's own execute() called again with the runtime this
// process still holds.
//
// This file is a probe, not a tracked test: it is compiled into cmd/wisp with
// `go test -overlay` from .scratch/wisp/probes/176/r1/overlay-e2e.json, and it
// writes its readings next to itself in logs/e2e-readings.txt.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/secret"
)

// long176r1 is the background task's answer: 20000 ASCII bytes, over the
// D15(3) spill threshold at the reference window (4000 tokens / 16000 bytes),
// and under the raw hard cap.
var long176r1 = strings.Repeat("WISP176R1-background-output-line.", 700)[:20000]

// pointer176r1Re reads the pointer out of a stub the HOST wrote.
var pointer176r1Re = regexp.MustCompile(`全文见 ([0-9A-Za-z:\\_./%-]+)`)

type rig176r1 struct {
	mu      sync.Mutex
	srv     *httptest.Server
	dir     string
	taskID  string
	pointer string
	turns   []string // one entry per provider request: what this rig answered
	calls   []string // tool calls the model made, in order
}

func (r *rig176r1) handler(w http.ResponseWriter, req *http.Request) {
	if !strings.HasSuffix(req.URL.Path, "/chat/completions") {
		http.NotFound(w, req)
		return
	}
	var body struct {
		Messages []struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
			Name    string          `json:"name"`
			CallID  string          `json:"tool_call_id"`
		} `json:"messages"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "rig cannot read request: "+err.Error(), http.StatusBadRequest)
		return
	}
	last := body.Messages[len(body.Messages)-1]
	text := contentText(last.Content)

	w.Header().Set("Content-Type", "text/event-stream")
	r.mu.Lock()
	started := r.taskID != ""
	r.mu.Unlock()

	switch {
	case !started && last.Role == "user":
		// Run A: this IS the background task's work - an oversized answer.
		r.record("run A: reply with the oversized task output")
		sseContent(w, long176r1)

	case last.Role == "user":
		// Run B: the model is asked to read that task back, and it calls the
		// roster tool with the id the loop minted.
		var rl []string
		for _, mm := range body.Messages {
			ml := contentText(mm.Content)
			if len(ml) > 40 {
				ml = ml[:40]
			}
			rl = append(rl, mm.Role+"|"+mm.Name+"|"+mm.CallID+"|"+ml)
		}
		r.record("run B: tool call task.output " + r.taskID + " || roles: " + strings.Join(rl, " ; "))
		sseToolCall(w, "call-176r1-task-output", "task.output",
			map[string]any{"task_id": r.taskID})

	case last.Role == "tool":
		head := text
		if len(head) > 240 {
			head = head[:240]
		}
		r.record("tool 回来的内容：" + head + fmt.Sprintf(" （共 %d 字节）", len(text)))
	case last.Role == "tool" && strings.Contains(text, "查不到这个任务"):
		r.record("run B: the roster answered 查不到 - end to end did NOT land")
		sseContent(w, "失败：名册里没有这条记录")

	case last.Role == "tool" && strings.Contains(text, "全文见"):
		// The model follows the pointer the host printed. This is the leg
		// ticket 175/177 built the exemption for.
		m := pointer176r1Re.FindStringSubmatch(text)
		if m == nil {
			r.record("run B: no pointer found in the stub - cannot continue")
			sseContent(w, "失败：桩里没有可续读的指针")
			return
		}
		r.mu.Lock()
		r.pointer = m[1]
		r.mu.Unlock()
		r.record("run B: tool call fs.read " + m[1])
		sseToolCall(w, "call-176r1-reread", "fs.read", map[string]any{"path": m[1]})

	case last.Role == "tool":
		// The pointer came back through the door; report what was read.
		head := text
		if len(head) > 120 {
			head = head[:120]
		}
		r.record(fmt.Sprintf("run B: fs.read 回来 %d 字节，开头 %q", len(text), head))
		sseContent(w, fmt.Sprintf("续读完成：拿到 %d 字节", len(text)))

	default:
		r.record("run: unexpected shape, replying with text")
		sseContent(w, text)
	}
}

func (r *rig176r1) record(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.turns = append(r.turns, line)
	if i := strings.Index(line, "tool call "); i >= 0 {
		r.calls = append(r.calls, line[i+len("tool call "):])
	}
}

// contentText renders a message's content whether the adapter sent a string or
// a content-part array.
func contentText(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &parts); err == nil {
		var b strings.Builder
		for _, p := range parts {
			b.WriteString(p.Text)
		}
		return b.String()
	}
	return ""
}

func sseStart(w http.ResponseWriter) {
	fmt.Fprint(w, "data: "+`{"choices":[{"index":0,"delta":{"role":"assistant","content":""}}]}`+"\n\n")
}

func sseContent(w http.ResponseWriter, text string) {
	sseStart(w)
	for _, chunk := range split176r1(text, 2000) {
		b, _ := json.Marshal(map[string]any{
			"choices": []any{map[string]any{
				"index": 0,
				"delta": map[string]any{"content": chunk},
			}},
		})
		fmt.Fprintf(w, "data: %s\n\n", b)
	}
	fmt.Fprint(w, "data: "+`{"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`+"\n\n")
	fmt.Fprint(w, "data: "+`{"choices":[],"usage":{"prompt_tokens":120,"completion_tokens":120,"total_tokens":240}}`+"\n\n")
	fmt.Fprint(w, "data: [DONE]\n\n")
}

func sseToolCall(w http.ResponseWriter, id, tool string, args map[string]any) {
	sseStart(w)
	encoded, _ := json.Marshal(args)
	b, _ := json.Marshal(map[string]any{
		"choices": []any{map[string]any{
			"index": 0,
			"delta": map[string]any{"tool_calls": []any{map[string]any{
				"index": 0, "id": id, "type": "function",
				"function": map[string]any{"name": tool, "arguments": string(encoded)},
			}}},
		}},
	})
	fmt.Fprintf(w, "data: %s\n\n", b)
	fmt.Fprint(w, "data: "+`{"choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`+"\n\n")
	fmt.Fprint(w, "data: "+`{"choices":[],"usage":{"prompt_tokens":120,"completion_tokens":120,"total_tokens":240}}`+"\n\n")
	fmt.Fprint(w, "data: [DONE]\n\n")
}

func split176r1(s string, n int) []string {
	var out []string
	for len(s) > 0 {
		cut := n
		if cut > len(s) {
			cut = len(s)
		}
		out = append(out, s[:cut])
		s = s[cut:]
	}
	return out
}

// newRig176r1 writes the config the CLI reads and starts the SSE endpoint.
func newRig176r1(t *testing.T) *rig176r1 {
	t.Helper()
	r := &rig176r1{}
	r.srv = httptest.NewServer(http.HandlerFunc(r.handler))
	t.Cleanup(r.srv.Close)
	r.dir = t.TempDir()
	cfg := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
context_window = 128000

[fs]
allowed_dirs = [%q]
`, r.srv.URL+"/v1", filepath.ToSlash(r.dir))
	if err := os.WriteFile(filepath.Join(r.dir, configFileName), []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(r.dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestBackgroundTaskToTaskOutputToFsReadOnTheCLISeam176r1(t *testing.T) {
	r := newRig176r1(t)
	out, errb := &strings.Builder{}, &strings.Builder{}
	var rt *agentRuntime

	code := runTextTask(runSpec{
		argv:      []string{"把这份清单整份打印出来"},
		stdout:    out,
		stderr:    errb,
		dataDir:   r.dir,
		onRuntime: func(got *agentRuntime) { rt = got },
		notify:    func(string, string) error { return nil },
	})
	if code != 0 {
		t.Fatalf("run A（后台任务本体）退出 %d\nstdout:\n%s\nstderr:\n%s", code, out, errb)
	}
	m := taskIDRe.FindStringSubmatch(out.String())
	if m == nil {
		t.Fatalf("run A 没打出状态行：%s", out)
	}
	r.mu.Lock()
	r.taskID = m[1]
	r.mu.Unlock()

	// (1) the write side actually wrote: the roster holds this task's answer and
	// the path of the copy that really landed.
	rec, ok := rt.tasks.Look(r.taskID)
	if !ok {
		t.Fatalf("起跑口落地之后名册里还是查不到 %s（Count()=%d）：task.output 在真机上仍然只能拒绝",
			r.taskID, rt.tasks.Count())
	}
	if rec.Text != long176r1 {
		t.Errorf("名册里的正文与任务打印的不一致：%d 字节 vs %d", len(rec.Text), len(long176r1))
	}
	if rec.ArtifactPath == "" {
		t.Fatal("超长输出没有产物路径（只截不指）")
	}
	if base := filepath.Base(rec.ArtifactPath); !strings.Contains(base, "agent-task-"+r.taskID) {
		t.Errorf("产物名 %q 不带 agent-task- 前缀：后台记录与模型 call id 又共用命名空间", base)
	}
	body, err := os.ReadFile(rec.ArtifactPath)
	if err != nil {
		t.Fatalf("登记的路径读不回来：%v", err)
	}
	if string(body) != long176r1 {
		t.Errorf("产物里换了字节：%d vs %d", len(body), len(long176r1))
	}

	// (2) the model's turn: task.output -> pointer -> fs.read on that pointer.
	cardsBefore := rt.windowCount()
	codeB := rt.execute("把刚才那个后台任务的输出读回来")
	if codeB != 0 {
		t.Fatalf("run B（模型续读）退出 %d\nstdout:\n%s\nrig-turns:\n%s", codeB, out, r.dump())
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if got := rt.windowCount(); got != cardsBefore {
		t.Errorf("续读这一路上下了 %d 枚确认卡（want 0）：宿主自己写下的指针被自家门拒了",
			got-cardsBefore)
	}
	if r.pointer != filepath.ToSlash(rec.ArtifactPath) && r.pointer != rec.ArtifactPath {
		t.Errorf("模型续读用的路径 %q 不是名册登记的那条 %q", r.pointer, rec.ArtifactPath)
	}
	if len(r.calls) != 2 || !strings.HasPrefix(r.calls[0], "task.output") ||
		!strings.HasPrefix(r.calls[1], "fs.read") {
		t.Errorf("模型这两发调用应为 task.output 然后 fs.read，got %v", r.calls)
	}
	var readBack string
	for _, line := range r.turns {
		if i := strings.Index(line, "fs.read 回来 "); i >= 0 {
			readBack = line
		}
	}
	if readBack == "" {
		t.Fatal("没有「fs.read 回来」那一行：续读这条路没走通")
	}
	whole := out.String() + errb.String()
	if strings.Contains(whole, "查不到这个任务") {
		t.Errorf("真机这一路还是出现了「查不到这个任务」：\n%s", whole)
	}
	if strings.Contains(whole, "路径授权判定者未接线") {
		t.Log("已登记的形状：task.output 的指针附带「C26 判定者未接线」那句 fail-closed 说明" +
			"（TaskDeps.Paths 在 assembleRuntime 里没接，改它超出本票具名放开的写面）")
	}
	writeReadings176r1(t, r, rec.ArtifactPath, whole)
}

// writeReadings176r1 drops the reading next to this probe (logdir = the
// directory this source file lives in, per the dispatch's §6 note).
func writeReadings176r1(t *testing.T, r *rig176r1, pointer, console string) {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Skip("runtime.Caller failed: 不落读数")
	}
	dir := filepath.Join(filepath.Dir(thisFile), "logs")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Errorf("读数目录建不出来：%v", err)
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, "task id（环路铸）: %s\nartifact（真落盘）: %s\n", r.taskID, pointer)
	for _, line := range r.turns {
		b.WriteString(line + "\n")
	}
	b.WriteString("--- console (run A + run B) ---\n" + console)
	if err := os.WriteFile(filepath.Join(dir, "e2e-readings.txt"), []byte(b.String()), 0o600); err != nil {
		t.Errorf("读数写不出去：%v", err)
	}
}

func (r *rig176r1) dump() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return strings.Join(r.turns, " || ")
}
