package main

// 票 174 AC#2b 的端到端判据（写腿 174-r8，脑子里的那条接线是
// cmd/wisp/run.go 的 tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}）。
//
// 这一枚不重建接线、不用假桥：它走 assembleRuntime 交出来的那一座真桥与那一份
// 真 TaskDeps，拿回**模型会读到的那一段回执文本**，问两件事：
//
//	T1 ...174r8SpeaksTheRealReason
//	    出厂授权形状（[fs] allowed_dirs = []，SPEC-03:35 的默认值）下，宿主
//	    登记一枚落在授权根外的真副本文件 ⇒ 回执必须给出精确文案（"不在你被
//	    授权的目录范围内……加进 [fs] allowed_dirs"），且**不许**再说"判定者
//	    未接线"。把 run.go 那枚装配行里的 Paths 字段摘掉（-overlay 突变）这一
//	    枚必须变红，红句里就是"未接线"那句原文——判据的另一半两向。
//	T2 ...174r8StaysQuiet
//	    根内＋真实普通文件 ⇒ 回执里"注意："零次（安静正例；把说实话写成噪音
//	    同样是不合格的写法，这条钉的不是 T1 的重复而是它的反向）。
//
// 两枚都只经生产实例做判定：C26 是 rt.paths 本尊（与 fs.read 同一枚），存在性
// 是只读 os.Stat 的既有实现（internal/tools/task.go），本文件零 filepath.Clean
// /Abs 决策、零配置读取——AC#2 判据③的禁区在这里同样成立。

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/tools"
)

// new174r8EmptyRootsFixture 是 newRunFixture 的"出厂授权形状"兄弟：唯一差别是
// [fs] allowed_dirs 写成空表，也就是文档冻结的出厂默认值。其余字段逐字相同，
// 于是两枚用例的差别只剩"授权根里有没有这枚副本文件"这一件事。
func new174r8EmptyRootsFixture(t *testing.T) *runFixture {
	t.Helper()
	f := &runFixture{t: t, srv: adaptertest.StartMockllm(t), out: &bytes.Buffer{}, err: &bytes.Buffer{}}
	f.dir = t.TempDir()
	cfgPath := filepath.Join(f.dir, configFileName)
	body := fmt.Sprintf(`schema_version = 2

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
enabled = true
context_window = 128000

[fs]
allowed_dirs = []
`, f.srv.Base+"/v1")
	if err := os.WriteFile(cfgPath, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(f.dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}
	return f
}

// spill174r8 造一段必然超过 D15(3) 参考预算（4000 token / 16000 字节）的长输出，
// 并在 artifact 处落一枚真的普通文件，模拟"宿主真把全文落了盘"的那一发。
func spill174r8(t *testing.T, artifact string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(artifact), 0o700); err != nil {
		t.Fatal(err)
	}
	long := strings.Repeat("段落内容", 4000)
	if err := os.WriteFile(artifact, []byte(long), 0o600); err != nil {
		t.Fatal(err)
	}
	return long
}

// dispatchTaskOutput174r8 登记一枚名册行，然后经组合根自己那座桥发一发
// task.output（任务先走组合根自己的 per-task 边界 admitTask，与环路同款），
// 返回模型会读到的那段回执。f.run 之后 f.rtHook 关闭，回执在测试进程里被断言。
//
// receipt174r8 取出回执原文，并且先在测试日志里逐字打一遍——`-v` 的绿箭头和
// 红腿里都能看到同一段文本，这正是"两向读数"要贴的那一物。
func dispatchTaskOutput174r8(t *testing.T, f *runFixture, artifact string) agent.ToolOutcome {
	t.Helper()
	long := spill174r8(t, artifact)
	var out agent.ToolOutcome
	f.rtHook = func(rt *agentRuntime) {
		rt.tasks.Record("t-174r8", tools.TaskOutput{Text: long, ArtifactPath: artifact})
		revoke := rt.admitTask("host-174r8")
		defer revoke()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		args, err := json.Marshal(map[string]string{"task_id": "t-174r8"})
		if err != nil {
			t.Fatal(err)
		}
		var execErr error
		out, execErr = rt.bridge.Execute(ctx, agent.ToolRequest{
			TaskID: "host-174r8", CorrelationID: "host-174r8-corr", CallID: "call-174r8",
			Name: "task.output", Args: args,
		})
		if execErr != nil {
			t.Fatalf("task.output 经真桥执行返回宿主错误（不是工具结果）：%v", execErr)
		}
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	t.Logf("task.output 回执逐字（模型可见的那一段）：\n%s", out.Text)
	return out
}

// T1：授权根外那枚真副本 ⇒ 回执给出精确文案、且不再说"未接线"。
func TestTaskOutputPointerNotice174r8SpeaksTheRealReason(t *testing.T) {
	f := new174r8EmptyRootsFixture(t)
	artifact := filepath.Join(f.dir, "artifacts", "tool-output-call-174r8.txt")
	out := dispatchTaskOutput174r8(t, f, artifact)

	if out.IsError {
		t.Errorf("拿回执这条路是正常结果，不该是错误结果：IsError=true text=%q", out.Text)
	}
	if !out.Truncated {
		t.Errorf("超过 4000 token 的输出必须走截断那一支：Truncated=false text=%q", out.Text)
	}
	if !strings.Contains(out.Text, "全文见 "+artifact) {
		t.Errorf("指针本身照旧要给（只截不指＝不合格）：回执里找不到 %q\n回执：%s",
			"全文见 "+artifact, out.Text)
	}
	if !strings.Contains(out.Text, "不在你被授权的目录范围内") ||
		!strings.Contains(out.Text, "加进 [fs] allowed_dirs") {
		t.Errorf("授权根外的副本路径必须给出精确文案（不在授权范围内＋加进 [fs] allowed_dirs 两条都在），"+
			"实际回执：\n%s", out.Text)
	}
	if strings.Contains(out.Text, "未接线") {
		t.Errorf("判定者已经接进生产（run.go 的 TaskDeps.Paths），回执里不许再出现「未接线」："+
			"\n%s", out.Text)
	}
	if n := strings.Count(out.Text, "注意："); n != 1 {
		t.Errorf("这一发的实情只有一条（路径在授权根外；文件真实存在于盘上），"+
			"注意：应恰好出现 1 次，实际 %d 次：\n%s", n, out.Text)
	}
}

// T2：安静正例——根内＋真实普通文件 ⇒ 零说明、指针仍在。
func TestTaskOutputPointerNotice174r8StaysQuiet(t *testing.T) {
	f := newRunFixture(t, "openai-chat") // allowed_dirs = [数据根]，副本就在根内
	artifact := filepath.Join(f.dir, "artifacts", "tool-output-call-174r8.txt")
	out := dispatchTaskOutput174r8(t, f, artifact)

	if out.IsError {
		t.Errorf("安静正例不该是错误结果：IsError=true text=%q", out.Text)
	}
	if !out.Truncated {
		t.Errorf("超过 4000 token 的输出必须走截断那一支：Truncated=false text=%q", out.Text)
	}
	if !strings.Contains(out.Text, "全文见 "+artifact) {
		t.Errorf("指针必须照旧在：回执里找不到 %q\n回执：%s", "全文见 "+artifact, out.Text)
	}
	if n := strings.Count(out.Text, "注意："); n != 0 {
		t.Errorf("健康的指针必须安静：注意：应出现 0 次，实际 %d 次：\n%s", n, out.Text)
	}
	if strings.Contains(out.Text, "读不到") {
		t.Errorf("根内真文件不该触发任何一条「读不到」的说明：\n%s", out.Text)
	}
}
