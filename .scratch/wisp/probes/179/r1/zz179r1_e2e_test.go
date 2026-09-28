package main

// Ticket 179-r1 - the end-to-end reading ticket 179 exists for, on the CLI
// seam (`wisp run` = runTextTask / assembleRuntime + execute), i.e. the ONLY
// shape where the defect was visible (cmd/wisp never sets
// agent.Config.PassThroughUnclassifiedRisk, so Loop.decideRisk's conflated
// branch refused every declared-L0 builtin).
//
// Real loop, real tools.Bridge, real risk assessor, real approval gate, real
// files, real OpenAI-compatible SSE endpoint (the C5 golden-SSE shape over a
// live HTTP server, per AGENTS §1.3's CLI seam). Zero mocks standing in for
// real components.
//
// Same two-run shape as 176-r1's probe, deliberately with one correction: the
// tool-result branches are ordered so the one that looks for the host-written
// pointer wins BEFORE the catch-all (176's rig put the catch-all first, so the
// pointer leg could never be reached). Run A fills the process-local roster
// with an oversized answer that spills to a real artifact; run B asks the model
// to read it back: task.output -> head+tail+total+pointer -> fs.read on that
// very pointer -> the bytes come back without our own gate refusing it.
//
// Probe, not a tracked test: compiled into cmd/wisp with
// `go test -overlay=.scratch/wisp/probes/179/r1/overlay-e2e.json`, readings
// land next to this file in logs/e2e-readings.txt.

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

var long179r1 = strings.Repeat("WISP179R1-background-output-line.", 700)[:20000]

var pointer179r1Re = regexp.MustCompile(`全文见 ([0-9A-Za-z:\\_./%-]+)`)

// unclassifiedSentence is the fail-closed wording Loop.decideRisk produced for
// a declared L0 call before ticket 179. Its absence is the point of this leg.
const unclassifiedSentence = "风险未分级"

type rig179r1 struct {
	mu      sync.Mutex
	srv     *httptest.Server
	dir     string
	taskID  string
	pointer string
	turns   []string
	calls   []string
}

func (r *rig179r1) handler(w http.ResponseWriter, req *http.Request) {
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
		// Run A: the background task's work is an oversized answer, which the
		// loop spills into a real artifact and points at.
		r.record("run A: 打印超长输出")
		sseContent179r1(w, long179r1)

	case last.Role == "tool" && strings.Contains(text, "全文见"):
		// The model follows the pointer the HOST wrote. This is the call that
		// ticket 179 refused as "unclassified".
		m := pointer179r1Re.FindStringSubmatch(text)
		if m == nil {
			r.record("run B: 桩里没有可续读的指针")
			sseContent179r1(w, "失败：桩里没有可续读的指针")
			return
		}
		r.mu.Lock()
		r.pointer = m[1]
		r.mu.Unlock()
		r.record("run B: tool call fs.read " + m[1])
		sseToolCall179r1(w, "call-179r1-reread", "fs.read", map[string]any{"path": m[1]})

	case last.Role == "tool" && strings.Contains(text, "查不到这个任务"):
		r.record("run B: 名册答「查不到」- 端到端没落地")
		sseContent179r1(w, "失败：名册里没有这条记录")

	case last.Role == "tool":
		// fs.read came back through the door: report the byte count.
		head := text
		if len(head) > 160 {
			head = head[:160]
		}
		r.record(fmt.Sprintf("run B: fs.read 回来 %d 字节，开头 %q", len(text), head))
		sseContent179r1(w, fmt.Sprintf("续读完成：拿到 %d 字节", len(text)))

	case last.Role == "user":
		// The loop appends a scene/user message AFTER every tool result, so
		// "last" is a user row even when the payload that matters is the tool
		// row further up. Scan the whole conversation instead.
		var convo strings.Builder
		longestToolRow := 0
		for _, mm := range body.Messages {
			ml := contentText(mm.Content)
			convo.WriteString(mm.Role + "|" + ml + "\n")
			if mm.Role == "tool" && len(ml) > longestToolRow {
				longestToolRow = len(ml)
			}
		}
		c := convo.String()

		switch {
		case !strings.Contains(c, "call-179r1-task-output"):
			r.record("run B: tool call task.output " + r.taskID)
			sseToolCall179r1(w, "call-179r1-task-output", "task.output",
				map[string]any{"task_id": r.taskID})

		case r.pointer == "" && strings.Contains(c, "全文见"):
			// task.output answered with head+tail+total+pointer: follow it.
			m := pointer179r1Re.FindStringSubmatch(c)
			if m == nil {
				r.record("run B: 有「全文见」但正则没抓到指针")
				sseContent179r1(w, "失败：指针抓不到")
				return
			}
			r.mu.Lock()
			r.pointer = m[1]
			r.mu.Unlock()
			r.record("run B: tool call fs.read " + m[1])
			sseToolCall179r1(w, "call-179r1-reread", "fs.read", map[string]any{"path": m[1]})

		case r.pointer != "":
			// fs.read已经回来：报字节数并收尾。
			r.record(fmt.Sprintf("run B: fs.read 回来 %d 字节（会话里最长的那条 tool 行）", longestToolRow))
			sseContent179r1(w, fmt.Sprintf("续读完成：最长 tool 行 %d 字节", longestToolRow))

		default:
			r.record("run B: task.output 回来了但会话里没有「全文见」指针 - 续读断在这里")
			sseContent179r1(w, "task.output 没有给出可续读的指针")
		}

	default:
		r.record("run: 意外形状，回文本")
		sseContent179r1(w, text)
	}
}

func (r *rig179r1) record(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.turns = append(r.turns, line)
	if i := strings.Index(line, "tool call "); i >= 0 {
		r.calls = append(r.calls, line[i+len("tool call "):])
	}
}

func (r *rig179r1) dump() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return strings.Join(r.turns, " ; ")
}

// contentText renders a message's content whether the adapter sent a string or
// a content-part array. cmd/wisp's tracked sources do not define it (176-r1's
// probe carried its own copy), so this probe owns its own.
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

func sseStart179r1(w http.ResponseWriter) {
	fmt.Fprint(w, "data: "+`{"choices":[{"index":0,"delta":{"role":"assistant","content":""}}]}`+"\n\n")
}

func sseContent179r1(w http.ResponseWriter, text string) {
	sseStart179r1(w)
	for _, chunk := range split179r1(text, 2000) {
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

func sseToolCall179r1(w http.ResponseWriter, id, tool string, args map[string]any) {
	sseStart179r1(w)
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

func split179r1(s string, n int) []string {
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

// newRig179r1 writes the config the CLI reads and starts the SSE endpoint.
func newRig179r1(t *testing.T) *rig179r1 {
	t.Helper()
	r := &rig179r1{}
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

func TestDeclaredL0ReadBackOnTheCLISeam179r1(t *testing.T) {
	r := newRig179r1(t)
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

	rec, ok := rt.tasks.Look(r.taskID)
	if !ok {
		t.Fatalf("名册里查不到 %s（Count()=%d）：task.output 在真机上仍然只能拒绝",
			r.taskID, rt.tasks.Count())
	}
	if rec.ArtifactPath == "" {
		t.Fatal("超长输出没有产物路径（只截不指）")
	}
	body, err := os.ReadFile(rec.ArtifactPath)
	if err != nil {
		t.Fatalf("登记的路径读不回来：%v", err)
	}
	if string(body) != long179r1 {
		t.Errorf("产物里的字节换了：%d vs %d", len(body), len(long179r1))
	}

	cardsBefore := rt.windowCount()
	codeB := rt.execute("把刚才那个后台任务的输出读回来")
	console := out.String() + errb.String()
	r.mu.Lock()
	turns := strings.Join(r.turns, "\n")
	pointer := r.pointer
	calls := append([]string{}, r.calls...)
	r.mu.Unlock()
	if codeB != 0 {
		t.Fatalf("run B（模型续读）退出 %d\nconsole:\n%s\nrig-turns:\n%s", codeB, console, turns)
	}
	if got := rt.windowCount(); got != cardsBefore {
		t.Errorf("续读这一路上下了 %d 枚确认卡（want 0）", got-cardsBefore)
	}

	// The whole point: no declared-L0 refusal anywhere on the real CLI.
	if strings.Contains(console, unclassifiedSentence) ||
		strings.Contains(turns, unclassifiedSentence) {
		t.Errorf("真机 CLI 上仍然出现「%s」\nconsole:\n%s\nrig-turns:\n%s",
			unclassifiedSentence, console, turns)
	}
	if len(calls) != 2 || !strings.HasPrefix(calls[0], "task.output") ||
		!strings.HasPrefix(calls[1], "fs.read") {
		t.Errorf("模型这两发应为 task.output 然后 fs.read，got %v", calls)
	}
	if pointer != filepath.ToSlash(rec.ArtifactPath) && pointer != rec.ArtifactPath {
		t.Errorf("模型续读用的路径 %q 不是名册登记的那条 %q", pointer, rec.ArtifactPath)
	}
	var readBack string
	for _, line := range strings.Split(turns, "\n") {
		if strings.Contains(line, "fs.read 回来 ") {
			readBack = line
		}
	}
	if readBack == "" {
		t.Errorf("没有「fs.read 回来」那一行：续读没走通\nrig-turns:\n%s\nconsole:\n%s", turns, console)
	}
	for _, marker := range []string{"已拒绝执行", "查不到这个任务", "路径授权判定者未接线"} {
		if strings.Contains(console, marker) || strings.Contains(turns, marker) {
			t.Logf("已登记的形状：尾态里出现「%s」那一类句子（见 logs/e2e-readings.txt 逐字）", marker)
		}
	}
	writeReadings179r1(t, r, rec.ArtifactPath, len(body), readBack, console, code, codeB)
}

func writeReadings179r1(t *testing.T, r *rig179r1, pointer string, artifactBytes int,
	readBack, console string, codeA, codeB int,
) {
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
	fmt.Fprintf(&b, "HEAD 修法后：run A 退出 %d / run B 退出 %d\n", codeA, codeB)
	fmt.Fprintf(&b, "task id（环路铸）: %s\nartifact（真落盘）: %s（%d 字节）\n",
		r.taskID, pointer, artifactBytes)
	fmt.Fprintf(&b, "「%s」出现次数: console=%d rig-turns=%d\n", unclassifiedSentence,
		strings.Count(console, unclassifiedSentence), strings.Count(r.dump(), unclassifiedSentence))
	r.mu.Lock()
	turns := strings.Join(r.turns, "\n")
	calls := strings.Join(r.calls, " | ")
	r.mu.Unlock()
	fmt.Fprintf(&b, "模型调用序列: %s\n续读那一行: %s\n", calls, readBack)
	b.WriteString("--- rig turns ---\n" + turns + "\n")
	b.WriteString("--- console (run A + run B) ---\n" + console + "\n")
	if err := os.WriteFile(filepath.Join(dir, "e2e-readings.txt"), []byte(b.String()), 0o600); err != nil {
		t.Errorf("读数写不出去：%v", err)
	}
}
