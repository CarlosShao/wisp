package main

// 票 47 装配腿（cmd/wisp 写面）的端到端判据，写法与 174r8 的两枚同款：
// 不重建接线、不用假桥，走 assembleRuntime 交出来的那一座真桥、那份真
// TaskDeps 与那枚真 C20 scheduler（run.go 装配块里 scheduler.New 的产物），
// 拿回**模型会读到的那一段回执文本**。
//
// 修的是什么缝：票 47 的实现腿（1cc67664b）把 TaskDeps.Scheduler 与
// NewBridgeDeps 都落了，但 run.go:653 那行 TaskDeps 没带 Scheduler 字段，
// 生产 task.list 一直答 fail-closed 的「调度器未接线」。本腿把真实例接进
// 装配行，两枚用例各钉一半：
//
//	T1 Test47AssemblyTaskListAnswersTheRealQueue
//	    空队列 ⇒ 回执是「队列为空」的正常成功，且**不许**再出现「未接线」。
//	    这是非 fail-closed 的判据本体：把 run.go 装配行里的 Scheduler 字段
//	    摘掉，这枚必须变红，红句里就是「未接线」那句原文。
//	T2 Test47AssemblyTaskListSeesARowAndCancelStopsIt
//	    经生产 scheduler 实例真 Submit 一发（体挂在自己的 root ctx 上，取消
//	    才结算）+ SetParent 到宿主任务名下 ⇒ task.list 回执里必须有那一行
//	    （id＋标题都在），task.cancel 以父任务身份发出停止必须成功——
//	    list/cancel 两条腿都经真桥走到真调度器，才算端到端。
//
// 两枚都只经生产实例判定：桥是装配根自己那座，调度器是装配根自己那枚，
// caller 身份走 admitTask 这道宿主自己的 per-task 边界（与环路同款）。

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/scheduler"
)

// dispatchTask47 经组合根自己那座桥发一发 task.*，返回模型会读到的回执。
// 调用方先 admitTask 登记宿主任务身份（D47 门拒掉没被收编的调用者）。
func dispatchTask47(t *testing.T, rt *agentRuntime, taskID, corr, callID, name string, args json.RawMessage) agent.ToolOutcome {
	t.Helper()
	revoke := rt.admitTask(taskID)
	defer revoke()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out, execErr := rt.bridge.Execute(ctx, agent.ToolRequest{
		TaskID: taskID, CorrelationID: corr, CallID: callID,
		Name: name, Args: args,
	})
	if execErr != nil {
		t.Fatalf("%s 经真桥执行返回宿主错误（不是工具结果）：%v", name, execErr)
	}
	t.Logf("%s 回执逐字（模型可见的那一段）：\n%s", name, out.Text)
	return out
}

// T1：空队列 ⇒ 正常成功＋「队列为空」，且零「未接线」。
func Test47AssemblyTaskListAnswersTheRealQueue(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	var out agent.ToolOutcome
	f.rtHook = func(rt *agentRuntime) {
		if rt.scheduler == nil {
			t.Fatal("装配根没有构造票 47 的 C20 scheduler（run.go 装配块的 scheduler.New 没落）：task.list 只能答 fail-closed")
		}
		out = dispatchTask47(t, rt, "host-47-list", "host-47-list-corr", "call-47-list",
			"task.list", json.RawMessage(`{}`))
	}
	if code := f.run("看看现在有什么任务在排队"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if out.IsError {
		t.Errorf("空队列是正常答案（调度器真的在，只是没任务），不该是错误结果：IsError=true text=%q", out.Text)
	}
	if strings.Contains(out.Text, "未接线") {
		t.Errorf("调度器已接进生产（run.go 装配块的 TaskDeps.Scheduler），回执里不许再出现「未接线」：\n%s", out.Text)
	}
	if !strings.Contains(out.Text, "队列为空") {
		t.Errorf("空队列必须如实答「队列为空」（这是真调度器的答案，不是占位文案）：\n%s", out.Text)
	}
}

// T2：真 Submit 一发 ⇒ task.list 看得到那一行，task.cancel 以父任务身份停得掉。
func Test47AssemblyTaskListSeesARowAndCancelStopsIt(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// task.cancel 声明 L1：票 375 的守卫要「有人可能答得上」才开窗口，所以这
	// 一发接一枚附着但沉默的控制台（与 panel_l1_blocked_220_test.go 同款），
	// 窗口无人否决到点即执行＝L1 路自己的语义，不是给 cancel 开后门。
	f.reply = strings.NewReader("")
	pinL1WindowCeiling220(t, f.dir)
	var listed, cancelled agent.ToolOutcome
	var rowID string
	f.rtHook = func(rt *agentRuntime) {
		// 经生产 scheduler 实例真收一发：体挂在任务自己的 root ctx 上，
		// 所以 task.cancel 发出的停止就是让它结算的那一下。
		id, err := rt.scheduler.Submit(scheduler.TaskSpec{
			Label: "票47装配校验任务",
			Exec: func(ctx context.Context, _ *scheduler.TaskRun) scheduler.TaskResult {
				<-ctx.Done()
				return scheduler.TaskResult{Status: "cancelled", Text: "校验任务已停"}
			},
		})
		if err != nil {
			t.Fatalf("生产 scheduler 拒收任务：%v", err)
		}
		rowID = id
		// 宿主把这一行的父记成 CALLER 的任务 id：task.cancel 的「只有父任务
		// 能停」检查读的就是这份数据（bridge.go 的 authority join）。
		if err := rt.scheduler.SetParent(id, "host-47-row"); err != nil {
			t.Fatalf("SetParent 失败：%v", err)
		}
		listed = dispatchTask47(t, rt, "host-47-row", "host-47-row-corr", "call-47-row-list",
			"task.list", json.RawMessage(`{}`))
		args, err := json.Marshal(map[string]string{"task_id": id})
		if err != nil {
			t.Fatal(err)
		}
		cancelled = dispatchTask47(t, rt, "host-47-row", "host-47-row-corr", "call-47-row-cancel",
			"task.cancel", args)
	}
	if code := f.run("看看队列里有什么，然后把我派的那发停掉"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if listed.IsError {
		t.Errorf("队列里有一行真任务，task.list 不该答错误：\n%s", listed.Text)
	}
	if strings.Contains(listed.Text, "未接线") {
		t.Errorf("task.list 必须读到真调度器的行，不许答「未接线」：\n%s", listed.Text)
	}
	if !strings.Contains(listed.Text, rowID) || !strings.Contains(listed.Text, "票47装配校验任务") {
		t.Errorf("task.list 回执里必须出现真收下的那一行（id %s 与标题都在）：\n%s", rowID, listed.Text)
	}
	if !strings.Contains(listed.Text, "调度器任务 1 行") {
		t.Errorf("本进程只收过这一发，行数必须是 1（多记＝把别的行算进来了）：\n%s", listed.Text)
	}
	if cancelled.IsError {
		t.Errorf("父任务身份的 task.cancel 必须成功停掉自己派的那行：\n%s", cancelled.Text)
	}
	if !strings.Contains(cancelled.Text, "发出停止") || !strings.Contains(cancelled.Text, rowID) {
		t.Errorf("task.cancel 回执必须点名这一行并说明停止已发出：\n%s", cancelled.Text)
	}
}
