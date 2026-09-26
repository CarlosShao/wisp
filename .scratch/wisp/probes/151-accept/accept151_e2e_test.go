package main

// 151-accept E2E shot (attack point 1 of the acceptance dispatch).
//
// The implementation record (§8#5) says the CLI leg can never reach
// mark->OpenTask end to end, because mockllm only answers a request that
// carries tool_choice, and the loop sends none. That is true of mockllm's
// DETERMINISTIC SYNTHESIS path. It is NOT true of mockllm's golden replay
// path: a request whose model is "golden/<name>" replays
// internal/llm/testdata/golden/<name>.sse verbatim, tool_calls included,
// and the multi-`# @response` grammar serves one section per request.
//
// So this case authors its own golden (turn 1 = model-initiated fs.read on a
// sensitive file inside the run's allowlist, turn 2 = a plain text answer),
// points config.toml at it, and drives the real composition root. If the
// C25 scope the loop's own task id earns is really closed at the task
// boundary, CloseTask's audit line must report was_open=true dropped=1
// for that very id — i.e. a TAINTED scope, not the dropped=0 shape the
// implementation's own T1 pins.
//
// Readings are written to stderr of the test log; the probe copy of the
// raw output lives in .scratch/wisp/probes/151-accept/.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
)

// accept151Rows reads back the forensics rows this run booked, so the reading
// names what the model actually asked for instead of trusting the golden.
func accept151Rows(t *testing.T, dir, stdout string) string {
	t.Helper()
	m := accept151TaskRe.FindStringSubmatch(stdout)
	if m == nil {
		return "no task id"
	}
	st, err := memory.Open(dir)
	if err != nil {
		return "open store: " + err.Error()
	}
	defer st.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	list, err := st.ListToolCalls(ctx)
	_ = m[1]
	if err != nil {
		return "list: " + err.Error()
	}
	var sb strings.Builder
	for _, r := range list {
		fmt.Fprintf(&sb, "[task=%s tool=%s args=%s level=%s decision=%s outcome=%s class=%s] ",
			r.TaskID, r.Tool, r.ArgsJSON, r.RiskLevel, r.Decision, r.Outcome, r.ErrorClass)
	}
	return sb.String()
}

const accept151Secret = "身份证号码 110101199003071234 的档案摘录，长度足够形成可匹配的片段"

var accept151TaskRe = regexp.MustCompile(`任务 (\S+) 结束`)

// accept151Golden writes the two-turn golden into THIS tree's golden dir
// (the one mockllm is started with) and returns the golden name.
func accept151Golden(t *testing.T, secretPath string) string {
	t.Helper()
	name := "accept151-e2e-fsread" // file stem; the wire selector is "golden/"+name
	dir := filepath.Join("..", "..", "internal", "llm", "testdata", "golden")
	body := fmt.Sprintf(`# wisp golden sse v1
# @scenario 151-accept: turn 1 issues fs.read on %s, turn 2 answers in text
# @response 200
data: {"choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_151a","type":"function","function":{"name":"fs.read","arguments":"{\"path\":\"%s\"}"}}]},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}

data: {"choices":[],"usage":{"prompt_tokens":21,"completion_tokens":9,"total_tokens":30}}

data: [DONE]

# @response 200
data: {"choices":[{"index":0,"delta":{"content":"已经读到那份档案摘录了。"},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: {"choices":[],"usage":{"prompt_tokens":21,"completion_tokens":9,"total_tokens":30}}

data: [DONE]
`, filepath.ToSlash(secretPath), filepath.ToSlash(secretPath))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".sse"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return "golden/" + name
}

func TestAccept151E2ELoopTaskClosesATaintedScope(t *testing.T) {
	srv := adaptertest.StartMockllm(t)
	dir := t.TempDir()
	secretPath := filepath.Join(dir, "idcard.txt")
	if err := os.WriteFile(secretPath, []byte(accept151Secret), 0o600); err != nil {
		t.Fatal(err)
	}
	golden := accept151Golden(t, secretPath)

	cfgPath := filepath.Join(dir, configFileName)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/%[1]s"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = "%[2]s/v1"
api_key_ref = "dpapi:acme"

[llm.providers.acme.models."%[1]s"]
context_window = 128000

[fs]
allowed_dirs = ["%[3]s"]
`, golden, srv.Base, filepath.ToSlash(dir))
	if err := os.WriteFile(cfgPath, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}

	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	code := runTextTask(runSpec{
		argv:    []string{"读一下 idcard.txt 那份档案"},
		stdout:  out,
		stderr:  errb,
		dataDir: dir,
		notify:  func(title, msg string) error { return nil },
	})
	stdout, stderr := out.String(), errb.String()
	t.Logf("exit=%d\n--- stdout ---\n%s\n--- stderr ---\n%s", code, stdout, stderr)

	if code != 0 {
		t.Fatalf("run exited %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}
	rows := accept151Rows(t, dir, stdout)
	t.Logf("tool_call rows: %s", rows)

	m := accept151TaskRe.FindStringSubmatch(stdout)
	if m == nil {
		t.Fatalf("no task id on the status line, so this run never completed a loop task\nstdout:\n%s", stdout)
	}
	id := m[1]

	// Did the model-initiated read really execute? The run books a tool_call
	// audit line per dispatch; name it instead of trusting the golden.
	dispatched := strings.Contains(stderr, "fs.read") || strings.Contains(stdout, "fs.read")
	if !dispatched {
		t.Errorf("这一轮里没有一发 fs.read 被派发（金样本没生效），后面的 dropped 断言不成立\nstdout:\n%s\nstderr:\n%s", stdout, stderr)
	}

	// The nail: the loop's own task id must close a scope that carried 1 mark.
	want := fmt.Sprintf("[audit] tools: C25 scope closed task=%s was_open=true dropped=1 open_scopes=0", id)
	if !strings.Contains(stderr, want) {
		lines := []string{}
		for _, l := range strings.Split(stderr, "\n") {
			if strings.Contains(l, "C25 scope closed") {
				lines = append(lines, l)
			}
		}
		t.Errorf("端到端这一轮没有在任务边界上关掉一枚带污点的 scope：stderr 里找不到 %q；实际看到的关闭审计行=%q", want, lines)
	}
}
