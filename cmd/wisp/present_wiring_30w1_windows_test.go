//go:build windows

package main

// 票 30 装配缝腿 30-w-r1 的端到端读数件：真装配根（assembleRuntime）×真环路
// （真 mockllm 子进程×3＋golden C5 缝×1，AGENTS §1.3 的两枚许可缝）×真 present
// 通道适配器，对照 30-r1 的 19 PASS 形状把「四档 → 各通道真触发」钉在 cmd/wisp 侧。
//
// 为什么不 mockllm 回声模式：cmd/wisp 的 prompt 恒带 scene 后缀（ContextMessage 的
// 当前时间块），回声回的是那个后缀块而不是任务文本——30-r1 的「回声＝样本本身」
// 成立在它直驱 agent.New 的零 scene 形状上。这里用同一个真 mockllm 子进程的
// golden 选择器（model "golden/<名>"，server 侧 RouteCount 照读），应答文本由
// 本腿自备的 golden 目录给出（共享 testdata 目录零写入）。
//
// 每案钉什么：
//
//	Test30Wr1TierFourShortCopiesAndSkipsTheRest   tier-4：剪贴板恰一次（内容逐字）、
//	                                              通知 0（旧呈现路已被装配门让位）、无工件
//	Test30Wr1TierThreeMediumNotifiesOnce          tier-3：通知恰一次（双发＝红）、
//	                                              剪贴板 0、无工件
//	Test30Wr1TierOneStructuredLandsTheArtifact    tier-1：工件真落盘＋面板指针行、
//	                                              通知 0（长结果不发通知＝D10）、剪贴板 0
//	Test30Wr1TierTwoGoldenLongLandsAndStripsTheSummary
//	                                              tier-2（golden C5 缝）：>400 汉字＋口播
//	                                              摘要行——单请求、工件剥掉摘要行、通知 0、
//	                                              剪贴板 0
//	Test30Wr1UnwiredHostKeepsTheLegacyPresentation
//	                                              反向判据的另一半：没装配 presenter 的
//	                                              宿主（CLI 形状）completed 结果仍走旧
//	                                              通知路——装配门没有把旧路打死
//	Test30Wr1ClipboardGateRefusesTaintedSentence  闸核负控：污句 ⇒ ChClipboard 拒写、
//	                                              真写半边零调用、ok=false 落日志；
//	                                              同形干净句 ⇒ 恰写一次（闸活着）

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/secret"
)

// presentWr1Clip is the recorder standing in the real write half's seat. The
// gate in front of it is the production one; only the Win32 body is swapped,
// which is the same seam shape runSpec.notify is for the notification side.
type presentWr1Clip struct {
	mu    sync.Mutex
	texts []string
}

func (c *presentWr1Clip) record(text string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.texts = append(c.texts, text)
	return nil
}

func (c *presentWr1Clip) got() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.texts...)
}

// presentWr1Mockllm is the real mockllm subprocess against THIS leg's own
// golden directory (the shared testdata directory is never written). It is
// StartMockllm's launch with one flag different, and the control plane
// (/__control/state) is read the same way for the server-side route counts.
type presentWr1Mockllm struct {
	base string
	cmd  *exec.Cmd
}

func presentWr1StartMockllm(t *testing.T, goldenDir string) *presentWr1Mockllm {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("go toolchain: %v", err)
	}
	exe := filepath.Join(t.TempDir(), "mockllm-30wr1.exe")
	build := exec.Command(goBin, "build", "-o", exe, ".")
	build.Dir = filepath.Join(adaptertest.RepoRoot(), "tools", "mockllm")
	if out, berr := build.CombinedOutput(); berr != nil {
		t.Fatalf("build tools/mockllm: %v\n%s", berr, out)
	}
	m := &presentWr1Mockllm{}
	m.cmd = exec.Command(exe, "-addr", "127.0.0.1:0", "-print-addr", "-golden-dir", goldenDir)
	// The address handoff runs through a StdoutPipe read on its own booked
	// goroutine - the same shape adaptertest.StartMockllm uses. A polled
	// bytes.Buffer would race the exec package's own copier goroutine (the
	// reading this file's first -race shot caught).
	stdout, perr := m.cmd.StdoutPipe()
	if perr != nil {
		t.Fatal(perr)
	}
	m.cmd.Stderr = os.Stderr
	if err := adaptertest.StartChildGuarded(t, m.cmd); err != nil {
		t.Fatalf("start mockllm: %v", err)
	}
	// Cleanup is registered before the readiness wait so every failure path
	// below still reaps the child and joins the reader spawn.
	spawnRoot := observe.NewRoot("adaptertest:mockllm-30wr1")
	ready := make(chan error, 1)
	reg := observe.NewRegistry()
	h := reg.Spawn("mockllm-addr-reader-30wr1", "test", spawnRoot,
		func(_ context.Context) {
			sc := bufio.NewScanner(stdout)
			for sc.Scan() {
				if line, ok := strings.CutPrefix(sc.Text(), "MOCKLLM_ADDR="); ok {
					m.base = "http://" + strings.TrimSpace(line)
					ready <- nil
					return
				}
			}
			ready <- errors.New("mockllm exited before printing its address")
		})
	t.Cleanup(func() {
		_ = m.cmd.Process.Kill()
		_, _ = m.cmd.Process.Wait()
		if left := spawnRoot.Wait(2 * time.Second); left != 0 {
			t.Errorf("mockllm stdout reader still pending after cleanup: %d", left)
		}
	})
	budget := observe.NewTimeout(20 * time.Second)
	timer := time.NewTimer(budget.Remaining())
	defer timer.Stop()
	var startErr error
	select {
	case startErr = <-ready:
	case <-h.Done():
		if perr := h.Err(); perr != nil {
			startErr = fmt.Errorf("mockllm stdout reader failed: %w", perr)
		} else {
			startErr = <-ready
		}
	case <-timer.C:
		startErr = errors.New("mockllm did not become ready in 20s")
	}
	if startErr != nil {
		_ = m.cmd.Process.Kill()
		t.Fatal(startErr)
	}
	return m
}

// routeCount reads the server-side counter for one route (chat).
func (m *presentWr1Mockllm) routeCount(t *testing.T) int {
	t.Helper()
	resp, err := http.Get(m.base + "/__control/state")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var state struct {
		Routes map[string]int `json:"routes"`
	}
	if err := json.Unmarshal(b, &state); err != nil {
		t.Fatalf("decode /__control/state: %v (%s)", err, b)
	}
	return state.Routes["chat"]
}

// presentWr1Fixture is one assembled composition root plus the recorders the
// wiring can be read through: the notification poster (runSpec's own seam),
// the clipboard write half, and the runtime handle the panel packet is read
// back from after the run closed its store.
type presentWr1Fixture struct {
	t    *testing.T
	dir  string
	out  *bytes.Buffer
	errb *bytes.Buffer

	ml  *presentWr1Mockllm // the mockllm shape
	rep *golden.Replayer   // the golden-server shape

	mu   sync.Mutex
	rt   *agentRuntime
	clip presentWr1Clip

	// notices records every poster call as title NUL body.
	notices []string
	// postWire runs inside onRuntime after the presenter is wired - the seam the
	// taint-marking case uses to poison the gate's scope once it exists.
	postWire func(*agentRuntime)
}

func presentWr1WriteConfig(t *testing.T, dir, baseURL, model string) {
	t.Helper()
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/%s"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models."%s"]
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]
`, model, baseURL+"/v1", model, filepath.ToSlash(dir))
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}
}

// presentWr1GoldenUnit writes one mockllm golden unit whose reply text is the
// scripted one, and returns the model id that selects it.
func presentWr1GoldenUnit(t *testing.T, dir, name, text string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name+".sse"),
		[]byte("# wisp golden sse v1\n# @scenario 30wr1 "+name+"\n# @response 200\n"+presentWr1TurnSSE(t, text)+"\n"),
		0o600); err != nil {
		t.Fatal(err)
	}
	return "golden/" + name
}

// presentWr1MockllmFixture assembles a run against the real mockllm subprocess
// with this leg's own golden directory (one unit per case, fresh server).
func presentWr1MockllmFixture(t *testing.T, unitName, replyText string) *presentWr1Fixture {
	t.Helper()
	f := &presentWr1Fixture{t: t, out: &bytes.Buffer{}, errb: &bytes.Buffer{}}
	f.dir = t.TempDir()
	goldenDir := filepath.Join(f.dir, "golden")
	if err := os.MkdirAll(goldenDir, 0o700); err != nil {
		t.Fatal(err)
	}
	model := presentWr1GoldenUnit(t, goldenDir, unitName, replyText)
	f.ml = presentWr1StartMockllm(t, goldenDir)
	presentWr1WriteConfig(t, f.dir, f.ml.base, model)
	return f
}

// presentWr1GoldenFixture is the tier-2 shape: the reply comes from the golden
// C5 seam (internal/llm/golden's replayer over real HTTP - the injection
// surface AGENTS §1.3 names), because >400 hanzi is beyond mockllm's echo and
// the seam gives the exact reply text - the same named bypass 30-r1's E2E
// recorded for this tier.
func presentWr1GoldenFixture(t *testing.T, goldenText string) *presentWr1Fixture {
	t.Helper()
	f := &presentWr1Fixture{t: t, out: &bytes.Buffer{}, errb: &bytes.Buffer{}}
	f.dir = t.TempDir()
	srv, rep, err := golden.ServerFromBytes([]byte(goldenText))
	if err != nil {
		t.Fatalf("golden.ServerFromBytes: %v", err)
	}
	t.Cleanup(srv.Close)
	f.rep = rep
	presentWr1WriteConfig(t, f.dir, srv.URL, "m1")
	return f
}

func (f *presentWr1Fixture) run(task string) int {
	f.t.Helper()
	return runTextTask(runSpec{
		argv:    []string{task},
		stdout:  f.out,
		stderr:  f.errb,
		dataDir: f.dir,
		onRuntime: func(rt *agentRuntime) {
			f.mu.Lock()
			f.rt = rt
			f.mu.Unlock()
			// The REAL assembly function, driven with nil ball/voice the way a
			// degraded boot leaves them: the nil ports are the router's own
			// skipped notes, and the wired ones are the production adapters.
			if dp, _ := assembleResidentPresenter(&residentTaskSource{run: rt}, nil, nil, f.clip.record); dp != nil {
				rt.presenter = dp
			}
			if f.postWire != nil {
				f.postWire(rt)
			}
		},
		notify: func(title, body string) error {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.notices = append(f.notices, title+"\x00"+body)
			return nil
		},
	})
}

func (f *presentWr1Fixture) noticeList() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.notices...)
}

// artifacts lists the names in the run's artifacts directory.
func (f *presentWr1Fixture) artifacts() []string {
	f.t.Helper()
	ents, err := os.ReadDir(filepath.Join(f.dir, "artifacts"))
	if err != nil {
		f.t.Fatalf("read the artifacts dir: %v", err)
	}
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	return names
}

// artifactBody reads one artifact's bytes back from disk.
func (f *presentWr1Fixture) artifactBody(name string) string {
	f.t.Helper()
	b, err := os.ReadFile(filepath.Join(f.dir, "artifacts", name))
	if err != nil {
		f.t.Fatalf("read artifact %s: %v", name, err)
	}
	return string(b)
}

// pointerRows reads the last booked packet's results section and returns the
// rows the presenter's panel adapter appended (keyed by the artifact name).
func (f *presentWr1Fixture) pointerRows() []panel.ResultChunk {
	f.t.Helper()
	rt := f.rt
	if rt == nil {
		f.t.Fatal("no runtime was captured")
	}
	snap, _, ok := rt.lastPanelSnapshot()
	if !ok {
		f.t.Fatal("no panel snapshot was ever booked")
	}
	var out []panel.ResultChunk
	for _, r := range snap.Results {
		if strings.HasPrefix(r.CorrelationID, "result-") {
			out = append(out, r)
		}
	}
	return out
}

func (f *presentWr1Fixture) mustGreen(code int, stage string) {
	f.t.Helper()
	if code != 0 {
		f.t.Fatalf("%s: exit %d\nstdout:\n%s\nstderr:\n%s", stage, code, f.out.String(), f.errb.String())
	}
}

// presentWr1LogHas reads the run's own persistent sink (the rolling JSONL
// under <data>\logs) and reports whether any record carries the needle.
func presentWr1LogHas(t *testing.T, dir, needle string) bool {
	t.Helper()
	ents, err := os.ReadDir(filepath.Join(dir, "logs"))
	if err != nil {
		t.Fatalf("read the log dir: %v", err)
	}
	for _, e := range ents {
		b, err := os.ReadFile(filepath.Join(dir, "logs", e.Name()))
		if err != nil {
			continue
		}
		if strings.Contains(string(b), needle) {
			return true
		}
	}
	return false
}

// presentWr1TurnSSE renders one openai-chat reply (the golden format's own
// byte shape, the one the pinned tool-call fixture and zz315's builder share).
func presentWr1TurnSSE(t *testing.T, text string) string {
	t.Helper()
	line := func(v any) string {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal sse chunk: %v", err)
		}
		return "data: " + string(b) + "\n\n"
	}
	var b strings.Builder
	b.WriteString(line(map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{"role": "assistant", "content": ""}, "finish_reason": nil,
	}}}))
	b.WriteString(line(map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{"content": text}, "finish_reason": nil,
	}}}))
	b.WriteString(line(map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{}, "finish_reason": "stop",
	}}}))
	b.WriteString(`data: {"choices":[],"usage":{"prompt_tokens":11,"completion_tokens":3,"total_tokens":14}}` + "\n\n")
	b.WriteString("data: [DONE]\n")
	return b.String()
}

// ---- the four tiers ---------------------------------------------------------

func Test30Wr1TierFourShortCopiesAndSkipsTheRest(t *testing.T) {
	f := presentWr1MockllmFixture(t, "short", "把这句复制一下")
	code := f.run("随便说一句")
	f.mustGreen(code, "tier-4 short")
	if n := f.ml.routeCount(t); n != 1 {
		t.Errorf("mockllm chat requests = %d, want 1 (single-request rule)", n)
	}
	writes := f.clip.got()
	if len(writes) != 1 {
		t.Fatalf("clipboard writes = %d, want 1 on tier 4: %v", len(writes), writes)
	}
	if writes[0] != "把这句复制一下" {
		t.Errorf("clipboard carried %q, want the reply text verbatim", writes[0])
	}
	if n := f.noticeList(); len(n) != 0 {
		t.Errorf("notifications = %v, want 0: a tier-4 result takes the clipboard and the ball, not a toast", n)
	}
	if names := f.artifacts(); len(names) != 0 {
		t.Errorf("artifacts = %v, want none on tier 4", names)
	}
}

func Test30Wr1TierThreeMediumNotifiesOnce(t *testing.T) {
	const medium = "这是一段中等长度的中文结果文本，用来落进第六十一到四百个汉字的那一档。它应该以系统通知交付，同时更新悬浮球角标与面板任务表，并且绝不碰剪贴板，也不落任何工件文件，播报只念首句。"
	f := presentWr1MockllmFixture(t, "medium", medium)
	code := f.run("给我一段中等长度的说明")
	f.mustGreen(code, "tier-3 medium")
	if n := f.ml.routeCount(t); n != 1 {
		t.Errorf("mockllm chat requests = %d, want 1 (single-request rule)", n)
	}
	notices := f.noticeList()
	if len(notices) != 1 {
		t.Fatalf("notifications = %d, want exactly 1 on tier 3 (a second one is the legacy path double-firing): %v",
			len(notices), notices)
	}
	if !strings.Contains(notices[0], "第六十一到四百个汉字的那一档") {
		t.Errorf("the toast carried %q, want the display text", notices[0])
	}
	if writes := f.clip.got(); len(writes) != 0 {
		t.Errorf("clipboard writes = %d, want 0: tiers 1-3 never touch the clipboard", len(writes))
	}
	if names := f.artifacts(); len(names) != 0 {
		t.Errorf("artifacts = %v, want none on tier 3", names)
	}
}

func Test30Wr1TierOneStructuredLandsTheArtifact(t *testing.T) {
	const table = "| 名 | 值 |\n|---|---|\n| 甲 | 1 |\n| 乙 | 2 |\n| 丙 | 3 |"
	f := presentWr1MockllmFixture(t, "table", table)
	code := f.run("给我一张表")
	f.mustGreen(code, "tier-1 structured")
	if n := f.ml.routeCount(t); n != 1 {
		t.Errorf("mockllm chat requests = %d, want 1 (single-request rule)", n)
	}
	names := f.artifacts()
	if len(names) != 1 {
		t.Fatalf("artifacts = %v, want exactly 1 landed file on tier 1", names)
	}
	body := f.artifactBody(names[0])
	if body != table {
		t.Errorf("artifact %s = %q, want the display text verbatim", names[0], body)
	}
	if n := f.noticeList(); len(n) != 0 {
		t.Errorf("notifications = %v, want 0: a long result lands a file, it does not toast", n)
	}
	if writes := f.clip.got(); len(writes) != 0 {
		t.Errorf("clipboard writes = %d, want 0 on tier 1", len(writes))
	}
	rows := f.pointerRows()
	if len(rows) != 1 {
		t.Fatalf("panel pointer rows = %d, want the one the tier-1 adapter appends", len(rows))
	}
	if !strings.Contains(rows[0].Text, "全文已存：") {
		t.Errorf("pointer row text = %q, want the landed-file sentence", rows[0].Text)
	}
}

func Test30Wr1TierTwoGoldenLongLandsAndStripsTheSummary(t *testing.T) {
	const (
		longBody   = "这一段超长的回复正文用来落在四百个汉字以上的第二档：它必须落成文件、在面板展示、把口播摘要剥出显示文本，而且全过程不得为生成摘要再发第二次请求。"
		spokenLine = "口播摘要：长文已归档，请查收"
	)
	f := presentWr1GoldenFixture(t, "# wisp golden sse v1\n# @scenario 30wr1 tier-2 long\n"+
		"# @response 200\n"+presentWr1TurnSSE(t, strings.Repeat(longBody, 7)+"\n"+spokenLine)+"\n")
	code := f.run("给我一段长文")
	f.mustGreen(code, "tier-2 long")
	if len(f.rep.Requests) != 1 {
		t.Errorf("golden requests = %d, want 1 (the summary must come from the same response)", len(f.rep.Requests))
	}
	names := f.artifacts()
	if len(names) != 1 {
		t.Fatalf("artifacts = %v, want exactly 1 landed file on tier 2", names)
	}
	body := f.artifactBody(names[0])
	if !strings.Contains(body, "超长的回复正文") {
		t.Errorf("artifact %s does not carry the long body:\n%s", names[0], body)
	}
	if strings.Contains(body, spokenLine) {
		t.Errorf("artifact %s still carries the spoken-summary line - the display text must have it stripped:\n%s", names[0], body)
	}
	if n := f.noticeList(); len(n) != 0 {
		t.Errorf("notifications = %v, want 0 on tier 2", n)
	}
	if writes := f.clip.got(); len(writes) != 0 {
		t.Errorf("clipboard writes = %d, want 0 on tier 2", len(writes))
	}
	if rows := f.pointerRows(); len(rows) != 1 {
		t.Fatalf("panel pointer rows = %d, want the one the tier-2 adapter appends", len(rows))
	}
}

// Test30Wr1UnwiredHostKeepsTheLegacyPresentation is the reconciliation's other
// half: a host that wired no presenter (the CLI's shape) still presents a
// completed result through the legacy notification path, unchanged. If this
// reds, the assembly gate above did not keep the else branch alive.
func Test30Wr1UnwiredHostKeepsTheLegacyPresentation(t *testing.T) {
	f := presentWr1MockllmFixture(t, "short", "把这句复制一下")
	code := runTextTask(runSpec{
		argv:    []string{"随便说一句"},
		stdout:  f.out,
		stderr:  f.errb,
		dataDir: f.dir,
		notify: func(title, body string) error {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.notices = append(f.notices, title+"\x00"+body)
			return nil
		},
	})
	f.mustGreen(code, "unwired host")
	notices := f.noticeList()
	if len(notices) != 1 {
		t.Fatalf("notifications = %d, want 1 on the legacy path (no presenter wired)", len(notices))
	}
}

// ---- the clipboard gate (the port's documented obligation) ------------------

func Test30Wr1ClipboardGateRefusesTaintedSentence(t *testing.T) {
	const tainted = "SK-wisp-30wr1-tainted-secret"

	t.Run("poisoned", func(t *testing.T) {
		f := presentWr1MockllmFixture(t, "tainted", "凭据是 "+tainted)
		f.postWire = func(rt *agentRuntime) {
			if !rt.c25.Mark(presentClipboardScope, risk.SrcFSRead, `D:\evil\notes.txt`, "凭据是 "+tainted) {
				t.Errorf("Mark into %s did not take: the negative control would prove nothing", presentClipboardScope)
			}
		}
		code := f.run("把密钥复制给我")
		f.mustGreen(code, "tainted clipboard")
		if writes := f.clip.got(); len(writes) != 0 {
			t.Fatalf("clipboard writes = %v, want 0: the ChClipboard gate must refuse before the real write", writes)
		}
		if !strings.Contains(f.errb.String(), "present: 剪贴板写入被 ChClipboard 闸拦下") {
			t.Errorf("the refusal is not named on the run's own audit surface:\n%s", f.errb.String())
		}
		// The routed receipt's ok=false lands on the persistent sink (slog's
		// mirror writes to the process stderr, which this fixture does not
		// hold), so the reading comes off the run's own log file on disk.
		if !presentWr1LogHas(t, f.dir, `"ok":false`) {
			t.Errorf("the routed receipt never said ok=false on the persistent sink:\n%s", f.errb.String())
		}
	})

	t.Run("clean", func(t *testing.T) {
		f := presentWr1MockllmFixture(t, "clean", "把这句复制一下")
		code := f.run("把密钥复制给我")
		f.mustGreen(code, "clean clipboard")
		writes := f.clip.got()
		if len(writes) != 1 {
			t.Fatalf("clipboard writes = %d, want exactly 1: the gate must pass a clean short result", len(writes))
		}
		if writes[0] != "把这句复制一下" {
			t.Errorf("clipboard carried %q, want the reply text verbatim", writes[0])
		}
	})
}
