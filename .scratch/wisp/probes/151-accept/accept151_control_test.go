package main

// 151-accept CONTROL for the E2E shot.
//
// The E2E case (accept151_e2e_test.go) got a REAL model-initiated fs.read all
// the way into the loop's turn; the console printed "[工具 fs.read -> error]"
// and CloseTask reported was_open=false dropped=0 - i.e. the call never reached
// the bridge, so mark/OpenTask never ran. Static reading at 4e16976:
//
//	internal/tools/fs.go:298      FSReadDecl -> Declared: risk.L0
//	internal/tools/bridge.go:224  ToolInfo.RiskLevel = levelString(e.Decl.Declared)
//	internal/agent/loop.go:776-789  decideRisk: only L1/L2 pass the switch; the
//	                              default branch rejects unless
//	                              Config.PassThroughUnclassifiedRisk is set, and
//	                              cmd/wisp/run.go:588-595 never sets it.
//
// That inference becomes a measurement with one control: the SAME leg, the SAME
// golden machinery, only the tool swapped for a call whose declared level is
// L1/L2 (fs.write, new file). If THAT reaches the bridge, it was the L0
// declaration that stopped fs.read - not the fixture, not my golden.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/secret"
)

var accept151ToolLineRe = regexp.MustCompile(`\[工具 (\S+) -> (\S+)\]`)

// accept151Run writes <goldenName>.sse into THIS tree's golden dir (body built
// once the data dir is known), points a real config.toml at it and drives the
// real composition root.
func accept151Run(t *testing.T, goldenName string, body func(dir string) string) (int, string, string, string) {
	t.Helper()
	srv := adaptertest.StartMockllm(t)
	dir := t.TempDir()
	gdir := filepath.Join("..", "..", "internal", "llm", "testdata", "golden")
	if err := os.WriteFile(filepath.Join(gdir, goldenName+".sse"),
		[]byte(body(filepath.Join(dir, "out.txt"))), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/golden/%[1]s"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = "%[2]s/v1"
api_key_ref = "dpapi:acme"

[llm.providers.acme.models."golden/%[1]s"]
context_window = 128000

[fs]
allowed_dirs = ["%[3]s"]
`, goldenName, srv.Base, filepath.ToSlash(dir))
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(cfg), 0o600); err != nil {
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
		argv:    []string{"把这行字写进 out.txt"},
		stdout:  out,
		stderr:  errb,
		dataDir: dir,
		notify:  func(title, msg string) error { return nil },
	})
	return code, out.String(), errb.String(), dir
}

func accept151ToolWriteGolden(target string) string {
	return fmt.Sprintf(`# wisp golden sse v1
# @scenario 151-accept control: turn 1 issues fs.write (declared L1), turn 2 answers
# @response 200
data: {"choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_151w","type":"function","function":{"name":"fs.write","arguments":"{\"path\":\"%[1]s\",\"content\":\"written-by-the-loop\"}"}}]},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}

data: {"choices":[],"usage":{"prompt_tokens":21,"completion_tokens":9,"total_tokens":30}}

data: [DONE]

# @response 200
data: {"choices":[{"index":0,"delta":{"content":"写好了。"},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: {"choices":[],"usage":{"prompt_tokens":30,"completion_tokens":3,"total_tokens":33}}

data: [DONE]
`, filepath.ToSlash(target))
}

func TestAccept151ControlL1DeclaredCallReachesTheBridge(t *testing.T) {
	code, stdout, stderr, dir := accept151Run(t, "accept151-e2e-fswrite",
		func(d string) string { return accept151ToolWriteGolden(d) })
	t.Logf("exit=%d\n--- stdout ---\n%s\n--- stderr ---\n%s", code, stdout, stderr)

	written := filepath.Join(dir, "out.txt")
	data, err := os.ReadFile(written)
	toolLines := accept151ToolLineRe.FindAllString(stdout, -1)
	closed := []string{}
	for _, l := range strings.Split(stderr, "\n") {
		if strings.Contains(l, "C25 scope closed") {
			closed = append(closed, strings.TrimSpace(l))
		}
	}
	t.Logf("control reading: tool_lines=%v file_written=%v err=%v close_lines=%q",
		toolLines, len(data), err, closed)

	if err != nil {
		t.Errorf("控制读数没有成立：这条腿上的 fs.write 没有落盘（%v），所以它和 fs.read 一样没能到达桥；两形的差别不能只归给声明档位", err)
	}
	if !strings.Contains(stdout, "fs.write") {
		t.Errorf("环路没有派发 fs.write，stdout=%q", stdout)
	}
}
