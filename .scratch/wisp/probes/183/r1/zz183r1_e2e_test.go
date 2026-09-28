package main

// Ticket 179-r2 - the one reading ticket 179's AC#8 still lacks: on the REAL
// CLI seam, the model follows the pointer the host itself wrote (task.output's
// "全文见 <artifact>") with fs.read, and the copy comes back byte for byte
// without our own gate refusing it.
//
// Same shape as 179-r1's probe, deliberately with ONE correction and zero
// production changes: 179-r1's rig looked for its call ids inside message
// *content*, but a call id travels in the tool_call_id FIELD of the tool row
// (internal/agent/loop.go:715 -> toolResultMessage(c.ID, ...)). That is why its
// run B re-issued task.output forever. This rig keys its state machine on that
// field, so the fs.read leg is actually reached.
//
// Real loop, real tools.Bridge, real risk assessor, real approval gate, real
// files, real OpenAI-compatible SSE endpoint (the C5 golden-SSE shape over a
// live HTTP server, AGENTS §1.3's CLI seam). Zero mocks standing in for real
// components.
//
// Probe, not a tracked test: compiled into cmd/wisp with
// `go test -overlay=.scratch/wisp/probes/179/r2/overlay-e2e.json`; readings land
// next to this file in logs/ (logdir from runtime.Caller, never from CWD).

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

// long183a1 is the run-A answer: big enough that the host spills it into a real
// artifact (SpillTokens = 4000 at a 128k window) and points at that artifact.
// long183a1D1 is 183-a1's M-D1 payload: 20000 bytes of pure CJK, so NO 8-rune
// window of the artifact path occurs anywhere in the marked body. Base payload
// repeats "...-output-line." whose "-output-" is also an 8-rune window of the
// path "tool-output-agent-task-...". If the exemption dies only with the base
// payload and survives here, the blocker is the duplicated fragment, not the
// carrier / not the span coordinates / not the Chinese notice.
var long183a1D1 = strings.Repeat("甲乙丙丁戊己庚辛壬癸子丑寅卯辰巳午未申酉戌亥", 500)[:20000]

// 183-r1 AC#2: the payload is pinned back to 179-r2's BASE body verbatim -
// strings.Repeat("WISP179R2-background-output-line.", 700)[:20000]. Its
// "-output-" is one of the artifact path's own 8-rune windows, which is the
// whole defect of ticket 183. 183-a1's M-D1 CJK payload (long183a1D1) proves
// the opposite half only by REMOVING the twin fragment, and dispatch 183-r1
// section 3 forbids crediting a green to a changed body, so this copy uses the
// base payload and nothing else in the rig changed.
var long183a1 = strings.Repeat("WISP179R2-background-output-line.", 700)[:20000]

// headReadBytes183a1 is the second leg's fs.read max_bytes: small enough that
// the read-back stays WHOLE in the context (no re-spill), so the model literally
// holds those bytes and they can be compared byte for byte with the artifact.
const headReadBytes183a1 = 10000

const (
	modeRunA183a1 = "runA"
	modeRunB183a1 = "runB"

	idOut183a1  = "call-183a1-out"
	idFull183a1 = "call-183a1-full"
	idHead183a1 = "call-183a1-head"
)

var (
	pointer183a1Re    = regexp.MustCompile(`全文见 ([0-9A-Za-z:\\_./%-]+)`)
	totalBytes183a1Re = regexp.MustCompile(`总长 (\d+) 字节`)
)

// forks183a1 is the four-way divergence ticket 179's AC#8 says must be reported
// separately: the unclassified refusal this ticket fixed, an L2 refusal (a
// different defect), an R4 taint hit (tickets 177/162), and a roster miss.
var forks183a1 = []string{"风险未分级", "L2 级", "R4", "查不到这个任务"}

type rig183a1 struct {
	mu            sync.Mutex
	srv           *httptest.Server
	dir           string
	mode          string
	taskID        string
	pointer       string
	rows          map[string]string
	turns         []string
	calls         []string
	console       string
	codeA         int
	codeB         int
	cardsA        int
	cardsB        int
	artifactPath  string
	artifactBytes int
	dumpInfo      string
}

func (r *rig183a1) handler(w http.ResponseWriter, req *http.Request) {
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

	// THE fix over 179-r1's rig: take the call id from the tool_call_id FIELD.
	r.mu.Lock()
	for _, m := range body.Messages {
		if m.Role == "tool" && m.CallID != "" {
			r.rows[m.CallID] = contentText183a1(m.Content)
		}
	}
	mode := r.mode
	taskID := r.taskID
	pointer := r.pointer
	outRow := r.rows[idOut183a1]
	fullRow := r.rows[idFull183a1]
	headRow := r.rows[idHead183a1]
	r.mu.Unlock()

	w.Header().Set("Content-Type", "text/event-stream")

	if mode == modeRunA183a1 {
		r.record183a1("run A: 交回超长正文，由宿主自己落产物并写指针")
		sseContent183a1(w, long183a1)
		return
	}

	switch {
	case outRow == "":
		r.record183a1("run B: tool call task.output " + taskID)
		sseToolCall183a1(w, idOut183a1, "task.output", map[string]any{"task_id": taskID})

	case pointer == "":
		m := pointer183a1Re.FindStringSubmatch(outRow)
		if m == nil {
			r.record183a1("run B: task.output 的正文里没有「全文见」指针 - 续读断在这里")
			sseContent183a1(w, "失败：task.output 没给出可续读的指针")
			return
		}
		r.mu.Lock()
		r.pointer = m[1]
		r.mu.Unlock()
		r.record183a1("run B: tool call fs.read " + m[1] + "（整份，不带 max_bytes）")
		sseToolCall183a1(w, idFull183a1, "fs.read", map[string]any{"path": m[1]})

	case fullRow == "":
		r.record183a1("run B: 整份 fs.read 那一发没有回来")
		sseContent183a1(w, "失败：整份读回那一发没有回来")

	case headRow == "":
		r.record183a1(fmt.Sprintf("run B: tool call fs.read %s（max_bytes=%d，整段留在上下文里）",
			pointer, headReadBytes183a1))
		sseToolCall183a1(w, idHead183a1, "fs.read",
			map[string]any{"path": pointer, "max_bytes": headReadBytes183a1})

	default:
		r.record183a1(fmt.Sprintf("run B: 续读收尾 - 整份那一发的 tool 行 %d 字节、有界那一发的 tool 行 %d 字节",
			len(fullRow), len(headRow)))
		sseContent183a1(w, "续读完成")
	}
}

func (r *rig183a1) record183a1(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.turns = append(r.turns, line)
	if i := strings.Index(line, "tool call "); i >= 0 {
		r.calls = append(r.calls, line[i+len("tool call "):])
	}
}

func (r *rig183a1) snapshot() (turns, calls, pointer string, rows map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := make(map[string]string, len(r.rows))
	for k, v := range r.rows {
		cp[k] = v
	}
	return strings.Join(r.turns, "\n"), strings.Join(r.calls, " | "), r.pointer, cp
}

// contentText183a1 renders a message's content whether the adapter sent a string
// or a content-part array. cmd/wisp's tracked sources do not define it, so this
// probe owns its own copy (same rule as 179-r1's probe, own names).
func contentText183a1(raw json.RawMessage) string {
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

func sseStart183a1(w http.ResponseWriter) {
	fmt.Fprint(w, "data: "+`{"choices":[{"index":0,"delta":{"role":"assistant","content":""}}]}`+"\n\n")
}

func sseContent183a1(w http.ResponseWriter, text string) {
	sseStart183a1(w)
	for _, chunk := range split183a1(text, 2000) {
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

func sseToolCall183a1(w http.ResponseWriter, id, tool string, args map[string]any) {
	sseStart183a1(w)
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

func split183a1(s string, n int) []string {
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

// newRig183a1 writes the config the CLI reads and starts the SSE endpoint.
func newRig183a1(t *testing.T) *rig183a1 {
	t.Helper()
	r := &rig183a1{rows: map[string]string{}, mode: modeRunA183a1}
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

func TestRereadHostPointerOnCLISeam183a1(t *testing.T) {
	r := newRig183a1(t)
	// 179-r1's readings never landed because the leg died before its writer ran.
	// Here the reading file is written by defer, so every path lands it.
	defer writeReadings183a1(t, r)

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
	r.codeA = code
	if code != 0 {
		t.Fatalf("run A（后台任务本体）退出 %d\nstdout:\n%s\nstderr:\n%s", code, out, errb)
	}
	m := taskIDRe.FindStringSubmatch(out.String())
	if m == nil {
		t.Fatalf("run A 没打出状态行：%s", out)
	}
	r.mu.Lock()
	r.taskID = m[1]
	r.mode = modeRunB183a1
	r.mu.Unlock()

	rec, ok := rt.tasks.Look(r.taskID)
	if !ok {
		t.Fatalf("名册里查不到 %s（Count()=%d）", r.taskID, rt.tasks.Count())
	}
	if rec.ArtifactPath == "" {
		t.Fatal("超长输出没有产物路径（只截不指）")
	}
	body, err := os.ReadFile(rec.ArtifactPath)
	if err != nil {
		t.Fatalf("登记的路径读不回来：%v", err)
	}
	r.mu.Lock()
	r.artifactPath = rec.ArtifactPath
	r.artifactBytes = len(body)
	r.mu.Unlock()
	if string(body) != long183a1 {
		t.Errorf("产物里的字节换了：%d vs %d", len(body), len(long183a1))
	}

	cardsBefore := rt.windowCount()
	codeB := rt.execute("把刚才那个后台任务的输出读回来")
	r.codeB = codeB
	r.console = out.String() + errb.String()
	r.mu.Lock()
	r.cardsA, r.cardsB = cardsBefore, rt.windowCount()
	r.mu.Unlock()

	turns, calls, pointer, rows := r.snapshot()
	if codeB != 0 {
		t.Fatalf("run B（模型续读）退出 %d\nconsole:\n%s\nrig-turns:\n%s", codeB, r.console, turns)
	}
	if r.cardsB != cardsBefore {
		t.Errorf("续读这一路上下了 %d 枚确认卡（want 0）", r.cardsB-cardsBefore)
	}
	fullRow, headRow := rows[idFull183a1], rows[idHead183a1]
	if fullRow == "" {
		t.Errorf("整份 fs.read 那一发的 tool 行没回来：续读没走通\nrig-turns:\n%s\nconsole:\n%s", turns, r.console)
	}
	if headRow == "" {
		t.Errorf("有界 fs.read 那一发的 tool 行没回来\nrig-turns:\n%s", turns)
	}

	// The point of the ticket: no declared-L0 refusal on the real CLI, and none
	// of the three other divergence shapes either.
	for _, fork := range forks183a1 {
		if strings.Contains(r.console, fork) || strings.Contains(turns, fork) ||
			strings.Contains(rows[idOut183a1]+fullRow+headRow, fork) {
			t.Errorf("真机 CLI 的续读那一发上出现「%s」\nconsole:\n%s\nrig-turns:\n%s", fork, r.console, turns)
		}
	}
	if strings.Contains(fullRow+headRow, "已拒绝执行") ||
		strings.Contains(fullRow+headRow, "路径无法解析") {
		t.Errorf("fs.read 那一发被拒了（逐字见读数文件）\n整份行：%q\n有界行：%q", clip183a1(fullRow), clip183a1(headRow))
	}
	callList := []string{}
	if calls != "" {
		callList = strings.Split(calls, " | ")
	}
	if len(callList) != 3 || !strings.HasPrefix(callList[0], "task.output") ||
		!strings.Contains(callList[1], "fs.read") || !strings.Contains(callList[2], "fs.read") {
		t.Errorf("模型这三发应为 task.output 然后两发 fs.read，got %v", callList)
	}
	if pointer != filepath.ToSlash(rec.ArtifactPath) && pointer != rec.ArtifactPath {
		t.Errorf("模型续读用的路径 %q 不是名册登记的那条 %q", pointer, rec.ArtifactPath)
	}

	// Byte-for-byte leg 1: the whole pointer read back. The loop re-spills a
	// 20000-byte tool result through the same D15(3) budget layer (not a gate),
	// so the evidence that fs.read really returned every byte is the host's own
	// stub: the byte count it announces, and the artifact it wrote to hold that
	// result - which must equal the file the pointer pointed at.
	if fullRow != "" {
		tm := totalBytes183a1Re.FindStringSubmatch(fullRow)
		pm := pointer183a1Re.FindStringSubmatch(fullRow)
		if tm == nil || pm == nil {
			t.Errorf("整份那一发的 tool 行不是「总长 N 字节 + 全文见」那一形：%q", clip183a1(fullRow))
		} else {
			var announced int
			if _, err := fmt.Sscanf(tm[1], "%d", &announced); err != nil {
				t.Errorf("总长那一段解析不出来：%q", tm[1])
			}
			second := strings.ReplaceAll(pm[1], "/", string(filepath.Separator))
			bodyB, errB := os.ReadFile(second)
			if errB != nil {
				t.Errorf("整份读回时宿主新落的那枚产物读不回来：%v", errB)
			} else {
				r.mu.Lock()
				r.dumpInfo = fmt.Sprintf("整份那一发：宿主宣告总长 %d 字节；宿主为那一发新落的产物 %s（%d 字节）；与指针所指文件逐字节相等＝%v",
					announced, second, len(bodyB), bytes.Equal(body, bodyB))
				r.mu.Unlock()
				if announced != len(body) || len(bodyB) != len(body) {
					t.Errorf("整份读回的字节数不等：产物 %d / 宿主宣告总长 %d / 新产物 %d",
						len(body), announced, len(bodyB))
				}
				if !bytes.Equal(body, bodyB) {
					t.Errorf("整份读回的内容与产物不逐字节相等（%d vs %d 字节）", len(body), len(bodyB))
				} else {
					t.Logf("整份读回逐字节相等：%d == %d（sha256 %s）", len(body), len(bodyB), sha183a1(string(bodyB)))
				}
			}
		}
	}

	// Byte-for-byte leg 2: the read-back that stays whole inside the model's own
	// context, compared with the artifact's first N bytes.
	if headRow != "" {
		want := body[:headReadBytes183a1]
		if len(headRow) != len(want) {
			t.Errorf("有界读回的字节数不等：产物前缀 %d / 读回 %d（max_bytes=%d）",
				len(want), len(headRow), headReadBytes183a1)
		}
		if !bytes.Equal([]byte(headRow), want) {
			t.Errorf("有界读回的前 %d 字节与产物不逐字节相等", headReadBytes183a1)
		} else {
			t.Logf("有界读回逐字节相等：%d 字节，sha256 %s", len(headRow), sha183a1(headRow))
		}
	}
}

func sha183a1(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func clip183a1(s string) string {
	if len(s) > 400 {
		return s[:400] + "…[剪]"
	}
	return s
}

func headTail183a1(s string, front bool) string {
	if len(s) <= 200 {
		return s
	}
	if front {
		return s[:200]
	}
	return s[len(s)-200:]
}

// logDir183a1 pins the reading dir to this probe's own folder. Under -overlay
// runtime.Caller reports the VIRTUAL path (cmd/wisp/zz183a1_e2e_test.go), which
// is how 179-r2's first run misplaced its readings into cmd/wisp/logs - so the
// probe's real location is only accepted when the path actually contains this
// probe's directory; otherwise the module root (found by walking up to go.mod)
// supplies it. Never the bare CWD of the test process.
func logDir183a1() (string, bool) {
	if _, thisFile, _, ok := runtime.Caller(0); ok {
		if abs, err := filepath.Abs(thisFile); err == nil {
			if strings.Contains(filepath.ToSlash(abs), "/.scratch/wisp/probes/183/r1/") {
				return filepath.Join(filepath.Dir(abs), "logs"), true
			}
		}
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return filepath.Join(wd, ".scratch", "wisp", "probes", "183", "r1", "logs"), true
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", false
		}
		wd = parent
	}
}

func writeReadings183a1(t *testing.T, r *rig183a1) {
	t.Helper()
	dir, ok := logDir183a1()
	if !ok {
		t.Error("定不出读数目录：读数不落盘")
		return
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Errorf("读数目录建不出来：%v", err)
		return
	}
	turns, calls, pointer, rows := r.snapshot()
	var b strings.Builder
	b.WriteString("179-r2 真机 CLI 续读读数（HEAD 修法后，零产码改动）\n")
	fmt.Fprintf(&b, "run A 退出 %d / run B 退出 %d / 确认卡 %d -> %d\n", r.codeA, r.codeB, r.cardsA, r.cardsB)
	fmt.Fprintf(&b, "task id（环路铸）: %s\n名册登记的产物: %s（%d 字节）\npayload 长度: %d\n",
		r.taskID, r.artifactPath, r.artifactBytes, len(long183a1))
	fmt.Fprintf(&b, "模型调用序列: %s\n模型续读用的路径: %s\n%s\n", calls, pointer, r.dumpInfo)
	for _, fork := range forks183a1 {
		fmt.Fprintf(&b, "分岔「%s」出现次数: console=%d rig-turns=%d tool-rows=%d\n", fork,
			strings.Count(r.console, fork), strings.Count(turns, fork),
			strings.Count(rows[idOut183a1]+rows[idFull183a1]+rows[idHead183a1], fork))
	}
	b.WriteString("--- rig turns ---\n" + turns + "\n")
	b.WriteString("--- 回执原文：task.output 那一发（tool_call_id=" + idOut183a1 + "）---\n")
	b.WriteString(rows[idOut183a1] + "\n")
	b.WriteString("--- 回执原文：整份 fs.read 那一发（tool_call_id=" + idFull183a1 + "）---\n")
	b.WriteString(rows[idFull183a1] + "\n")
	b.WriteString("--- 回执原文：有界 fs.read 那一发（tool_call_id=" + idHead183a1 + "）---\n")
	fmt.Fprintf(&b, "（该行 %d 字节，sha256=%s；逐字前 200 字节：%q；逐字后 200 字节：%q）\n",
		len(rows[idHead183a1]), sha183a1(rows[idHead183a1]),
		headTail183a1(rows[idHead183a1], true), headTail183a1(rows[idHead183a1], false))
	b.WriteString("--- console (run A + run B) ---\n" + r.console + "\n")
	if err := os.WriteFile(filepath.Join(dir, "e2e-readings.txt"), []byte(b.String()), 0o600); err != nil {
		t.Errorf("读数写不出去：%v", err)
	}
}
