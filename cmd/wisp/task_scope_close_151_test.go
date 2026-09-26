package main

// 票 151 AC#2 的两枚承重读数。两枚都走**真组合根**（真 config、真 store、真
// approval.Gate、真 tools.Bridge、真 agent.Loop、真 mockllm），不用假对象交回。
//
//	T1 TestCompositionRootClosesTheLoopTasksTaintScope
//	   一整轮真 `wisp run`，要求任务结束时那发 CloseTask 的审计行带着**本轮
//	   环路自己的 task id** 出现——即"关闭"这味真的长在环路任务的生命周期边界上，
//	   而不是只长在一个测试会去调的函数里。
//	T2 TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit
//	   同一条装配好的桥上：一枚宿主侧任务先读走一份敏感内容（污点挂在它身上），
//	   再走组合根的 per-task 边界收尾，然后让**另一枚新任务**发一次干净的、
//	   allowlist 内的 fs.read。摘掉关闭那一味之后第二发干净调用是 L2/R4
//	   （"包含来自 unbound-scope"）并被 D47 拒掉，带着关闭的这一版必须是 L0。
//
// 承重读数（变异＝只把 admitTask 里那一发 CloseTask 摘掉，其余一字不动）：
// 两枚都红，红句原文见 .scratch/wisp/probes/151/ac2-mutation-noclose.txt。

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
)

// req151 is one host-side tool request carrying its own task id, on the same
// shape cmd/wisp's other host-dispatch tests use.
func req151(task, name, path string) agent.ToolRequest {
	args, err := json.Marshal(map[string]any{"path": filepath.ToSlash(path)})
	if err != nil {
		panic(err)
	}
	return agent.ToolRequest{
		TaskID: task, CorrelationID: task, CallID: "call-" + task,
		Name: name, Args: args,
	}
}

// write151 puts a file inside the run's allowlisted data dir.
func write151(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCompositionRootClosesTheLoopTasksTaintScope(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	id := f.taskID()
	want := "[audit] tools: C25 scope closed task=" + id + " "
	if !strings.Contains(f.err.String(), want) {
		t.Errorf("任务结束时没有关闭它自己的 C25 污点 scope：stderr 里找不到以 %q 开头的审计行\nstderr:\n%s",
			want, f.err.String())
	}
}

func TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	source := write151(t, f.dir, "idcard.txt",
		"身份证号码 110101199003071234 的档案摘录，长度足够形成可匹配的片段")
	plain := write151(t, f.dir, "plain.txt", "与上一轮完全无关的普通内容")

	var (
		firstLevel  string
		firstErr    error
		secondLevel string
		secondText  string
		secondErr   error
	)
	f.rtHook = func(rt *agentRuntime) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// 第 1 轮：一枚宿主侧任务读走敏感内容 -> 污点挂在它自己的 scope 上。
		out1, err := rt.bridge.Execute(ctx, req151("host-151-taint", "fs.read", source))
		firstLevel, firstErr = out1.RiskLevel, err
		if err != nil || out1.IsError {
			t.Fatalf("第一轮敏感读没有执行: out=%+v err=%v", out1, err)
		}
		// 组合根的 per-task 边界收尾（环路在任务结束时 defer 的就是这一发）。
		rt.admitTask("host-151-taint")()
		// 第 2 轮：另一枚新任务、参数干净、路径在 allowlist 内。
		out2, err2 := rt.bridge.Execute(ctx, req151("host-151-clean", "fs.read", plain))
		secondLevel, secondText, secondErr = out2.RiskLevel, out2.Text, err2
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}

	if firstErr != nil {
		t.Fatal(firstErr)
	}
	if firstLevel != "L0" {
		t.Errorf("第一轮敏感读 judged %s, want L0 (allowlisted fs.read)", firstLevel)
	}
	if secondErr != nil {
		t.Fatal(secondErr)
	}
	if secondLevel != "L0" || strings.Contains(secondText, "unbound-scope") {
		t.Errorf("上一轮结束后它的污点 scope 还挂在表上，于是这一枚新任务的干净调用被 C25 的 unbound fail-closed 顶到 L2："+
			"第二发 judged %s, want L0；text=%q", secondLevel, secondText)
	}
	if !strings.Contains(secondText, "普通内容") {
		t.Errorf("第二轮没有读到文件内容，text=%q", secondText)
	}
}
