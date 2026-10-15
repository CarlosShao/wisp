package main

// 票 165-r2（装配半格）· 尺 b 与 c：装配可达尺＋批准只批这一批（装配版）。
//
// 尺 a（调用点枚数尺，plan_wiring_165r2_test.go）量"那一行在不在"；这一枚量"那一行
// 接上的东西真能被打到"，走 newReplyHost -> assembleRuntime 的真装配：真 registry、真
// bridge、真 C18 卡、真 tool_call 行、真会话授权表。接线没接上时两枚用例红在同一个地方
// ——这一发够不到装配，不是红在某句文案上。
//	b) AssembledBridgeReachesPlanPresent：registry 有没有 plan.present、声明面是不是
//	   冻结行那一枚（L2＋常驻）、模型面 Tools() 有没有报它，再把这一发真打一遍：升卡、
//	   原生侧答 yes、清单文本回得来、tool_call 落库、授权表零行。
//	c) UnlistedWriteRingsItsOwnCardAfterAnApprovedPlan（票面 AC#3① 装配版）：先批方案
//	   （只列一步），再打一发清单外的写（夹具里没被 [fs] 罩住那一半 ⇒ R2 走 L2，响的
//	   形式是一张可数的卡），必须自己再起一张卡；答 no 后不许执行、不许落盘、零授权行。
// 答复语序 verb <corr> [reason]（splitReplyLine 取第二段当 corr）：本腿第一版打反了，红成
// REPLY REFUSED ＋ ANSWER-EXPIRED；那条台面与摘掉那一行的正控两发逐字在证据件 §9 与 probes/165/r2/05、06。

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

func plan165r2Args(t *testing.T, steps ...string) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(map[string]any{"summary": "165-r2 装配探针", "steps": steps})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func writeArgs165r2(t *testing.T, path, content string) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(map[string]any{"path": filepath.ToSlash(path), "content": content})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// ask165r2 drives one tool call through the bridge THIS assembly built, waits for
// the card it raises, and answers that card from the native side of the reply
// stream. It returns the outcome, the tool name the card carried, and whether a
// card arrived at all - the honest "响了没有" reading, never an inference from a
// log line.
func ask165r2(t *testing.T, rt *agentRuntime, pw io.Writer, task, corr, tool string,
	args json.RawMessage, verb, reason string,
) (agent.ToolOutcome, string, bool) {
	t.Helper()
	var (
		out  agent.ToolOutcome
		done = make(chan struct{})
	)
	go func() {
		defer close(done)
		out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: task, CorrelationID: corr, CallID: corr, Name: tool, Args: args,
		})
	}()
	saw := waitForCard187(rt.liveCards, corr, 10*time.Second)
	asked := tool
	if saw {
		card, _ := rt.liveCards.look(corr)
		asked = card.Tool
	}
	line := verb + " " + corr
	if reason != "" {
		line += " " + reason
	}
	if _, err := fmt.Fprintln(pw, line); err != nil {
		t.Errorf("写答复 %q 失败：%v", line, err)
	}
	<-done
	return out, asked, saw
}

// grantsByMinted165r2 reads this boot's session grant rows back off disk.
func grantsByMinted165r2(t *testing.T, h *replyHost, minted string) []memory.ApprovalGrant {
	t.Helper()
	if minted == "" {
		t.Fatal("这一次装配没读到会话身份，授权表无从量起")
	}
	rows, err := h.openStore().ListGrantsBySession(context.Background(), minted)
	if err != nil {
		t.Fatalf("ListGrantsBySession(%s): %v", minted, err)
	}
	return rows
}

func Test165r2AssembledBridgeReachesPlanPresent(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	const task, corr = "t165r2-plan", "t165r2-plan-corr"
	var (
		registered, advertised, sawCard bool
		names                           []string
		declared                        risk.Level
		resident                        bool
		cardTool                        string
		out                             agent.ToolOutcome
		minted                          string
	)
	h.rtHook = func(rt *agentRuntime) {
		for _, e := range rt.bridge.Registry().List() {
			names = append(names, e.Tool.Name())
		}
		entry, ok := rt.bridge.Registry().Lookup(planToolName165r2)
		registered, declared, resident = ok, entry.Decl.Declared, entry.Decl.Resident
		list, err := rt.bridge.Tools(context.Background())
		if err != nil {
			t.Errorf("bridge.Tools 读不出来：%v", err)
			return
		}
		for _, i := range list {
			if i.Name == planToolName165r2 {
				advertised = i.Resident && i.RiskLevel == memory.RiskL2
			}
		}
		if rt.session == nil {
			t.Error("装配出来的运行没有会话授权账：「零授权行」这一发读不到")
			return
		}
		minted = rt.session.SessionID().String()
		revoke := rt.gate.AdmitTextTask(task)
		defer revoke()
		out, cardTool, sawCard = ask165r2(t, rt, pw, task, corr, planToolName165r2,
			plan165r2Args(t, "改 A 文件", "改 B 文件"), "yes", "")
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if !registered {
		t.Errorf("接线 RED：生产装配的 registry 里没有 %s；名册 %d 枚：%v"+
			"（这正是开工前 AC#1 的未修码形状：这份清单在 wisp run 里够不到）",
			planToolName165r2, len(names), names)
		return
	}
	if declared != risk.L2 || !resident {
		t.Errorf("装配读到的声明面变了：Declared=%v Resident=%v，冻结行（PLAN.md 的 D34 那一行／SPEC-07 §3）写的是 L2 且常驻",
			declared, resident)
	}
	if !advertised {
		t.Errorf("模型面读不到它：bridge.Tools() 里没有 %s，或它不是 L2/常驻那一项", planToolName165r2)
	}
	if !sawCard {
		t.Errorf("接线 RED：这一发没升起 L2 卡（card tool=%q，回文=%q）⇒「停下等人点头」在装配后仍不成立",
			cardTool, out.Text)
	} else if cardTool != planToolName165r2 {
		t.Errorf("升起的卡报的工具名是 %q，要 %s", cardTool, planToolName165r2)
	}
	if out.IsError {
		t.Errorf("答了 yes 却没拿到清单文本：%s", out.Text)
	} else {
		for _, want := range []string{"1. 改 A 文件", "2. 改 B 文件", "只限上面这 2 步", "清单之外"} {
			if !strings.Contains(out.Text, want) {
				t.Errorf("清单文本读不到 %q：%q", want, out.Text)
			}
		}
	}
	row := toolCallByCorr187(t, h, corr)
	if row.RiskLevel != memory.RiskL2 || row.Decision != agent.DecisionAllow {
		t.Errorf("tool_call 行 = %+v，要 L2/allow：这一发的档位与答复没落到生产记录上", row)
	}
	rows := grantsByMinted165r2(t, h, minted)
	if len(rows) != 0 {
		t.Errorf("批了一份清单就铸出 %d 行会话授权：批准只批这一发", len(rows))
	}
	t.Logf("装配可达读数：registry %d 枚含 %s；卡=%s；行 risk=%s decision=%s；授权行 %d 枚",
		len(names), planToolName165r2, cardTool, row.RiskLevel, row.Decision, len(rows))
}

func Test165r2UnlistedWriteRingsItsOwnCardAfterAnApprovedPlan(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	target := filepath.Join(h.outside, "not-in-plan-165r2.txt")
	const pTask, pCorr, wCorr = "t165r2-batch", "t165r2-batch-plan", "t165r2-batch-write"
	var (
		planOut, writeOut   agent.ToolOutcome
		planCard, writeCard string
		sawPlan, sawWrite   bool
		minted              string
	)
	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Error("装配出来的运行没有会话授权账")
			return
		}
		minted = rt.session.SessionID().String()
		revoke := rt.gate.AdmitTextTask(pTask)
		defer revoke()
		planOut, planCard, sawPlan = ask165r2(t, rt, pw, pTask, pCorr, planToolName165r2,
			plan165r2Args(t, "只改 in-plan.txt"), "yes", "")
		if !sawPlan {
			return
		}
		writeOut, writeCard, sawWrite = ask165r2(t, rt, pw, pTask, wCorr, "fs.write",
			writeArgs165r2(t, target, "清单里没有这一步"), "no", "清单里没有这一步，不批")
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if !sawPlan {
		t.Errorf("接线 RED：批方案那一发没升起卡（回文=%q）⇒ 装配里没有这一发，c) 的后续读不到",
			planOut.Text)
		return
	}
	if planOut.IsError {
		t.Fatalf("方案那一发被拒了：%s", planOut.Text)
	}
	if !sawWrite {
		t.Errorf("AC#3① RED：批完方案之后，清单外的写没有再升起一张卡 ⇒「方案被批准」变成了「后续写都被批准」：%+v",
			writeOut)
	} else if writeCard != "fs.write" {
		t.Errorf("第二张卡的工具名是 %q，要 fs.write", writeCard)
	}
	if !writeOut.IsError {
		t.Errorf("被答 no 的清单外写居然执行了：%+v", writeOut)
	}
	// The two answers have to be the human's, not the clock's: D43 #24's timeout
	// branch also ends in "not executed", so the assertions above would read green
	// while nobody in fact said no. This is the shape this case's first draft hit.
	audit := h.err.String()
	for _, want := range []string{
		"approval: ANSWER-ALLOW corr=" + pCorr + " tool=plan.present route=native decision=allow",
		"approval: ANSWER-REJECT corr=" + wCorr + " tool=fs.write decision=reject",
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("audit 里读不到这一行 %q（两张卡都要是人答的，不是钟到点的）；got:\n%s", want, audit)
		}
	}
	if strings.Contains(audit, "ANSWER-EXPIRED corr="+wCorr) {
		t.Errorf("清单外那一张卡是到点落拒绝的（ANSWER-EXPIRED），不是被人答 no 的 ⇒ 这一发没量到点头：%s", wCorr)
	}
	if _, err := os.Stat(target); err == nil {
		t.Errorf("清单外的写落了盘：%s", target)
	}
	rows := grantsByMinted165r2(t, h, minted)
	if len(rows) != 0 {
		t.Errorf("批方案铸出了 %d 行会话授权：批准范围被撑大了", len(rows))
	}
	t.Logf("AC#3① 装配读数：方案卡=%s（yes）→ 清单外卡=%s（no，IsError=%v）→ 授权行 %d 枚 → 目标文件未落盘",
		planCard, writeCard, writeOut.IsError, len(rows))
}
