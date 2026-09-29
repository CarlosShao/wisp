package main

// Ticket 201 AC#4/AC#5 - the 长期 branch, and the two things it owes.
//
// WHY A SECOND CARD IS NOT REDUNDANCY. PLAN.md:1645-1646 freezes one sentence
// about a loosening of [risk]/[fs]/[net]/[plugins]: 热加载放宽「必须触发 L2 级重
// 新确认，不得静默生效」. So an answer of 「一直」 to an fs card cannot be one
// confirmation doing two jobs: the card that lets THIS call run is about the call,
// and the stored rule widens the allowlist for every future call. This file makes
// them two events, in this order:
//
//	always <编号>      (1) prints the exact line that would be stored, derived
//	                       from the card's own Paths - AC#4's 「卡上出现规则文
//	                       本」, and the reason 票 219's 理由框 exists
//	                  (2) raises a second, L2-level card about that line, which
//	                       the operator then answers on its own correlation id.
//	                       That is PLAN.md:1645-1646's re-confirmation, and it is
//	                       AC#4's second confirmation.
//
// The write only happens if that second card is ALLOWED on the native route. An
// unanswered second card dies at the C18 deadline exactly like any other L2 card
// (queue.expire -> reject), so 不答 = nothing stored, which is the direction this
// repo requires.
//
// WHY THE CARD IS RAISED OFF THE LISTENER GOROUTINE. Gate.PendingApproval blocks
// until an answer or the deadline. The only thing that can answer it, in a console
// run, is the reply listener - so raising it from inside the listener parks the
// listener and the card dies unanswered every time. The widening therefore runs as
// the listener's own continuation on a derived root, and the loop keeps reading.
//
// IT CARRIES THE LISTENER'S ROSTER NAME ON PURPOSE. D38b's roster (internal/
// observe/goroutine.go:52) has exactly one approval-side per-task entry,
// "approval-waiter", and an unrostered name is what RosterReport reports as
// Unknown - i.e. a leak symptom. Minting a second roster name is a contract move
// this leg may not make quietly, so the continuation shares the existing entry and
// the gap is written up in this leg's evidence file instead: if the orchestrator
// wants the widening worker counted separately, that is a D38b edit, not a line of
// Go.
//
// WHAT IS STILL TRUE AFTER THE WRITE. The line lands in config.toml and is read by
// the NEXT start. The running process keeps the C26 canonicalizer it built at
// assembly (cmd/wisp/run.go:386-389), so this call and every later call in this run
// is still judged by the old allowlist - the promise 一直 makes is "next time", and
// the card says so. Making it live would be a second feature (re-feeding
// tools.NewPathCanonicalizer mid-run), and it is reported as such.

import (
	"context"
	"errors"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The host-side identity of a widening card. Same shape as the mode switch's
// (run.go:618-640, which is the precedent this copies): a task id the gate has
// never seen because no model call raised it, admitted for exactly the length of
// the confirmation, and a 'namespace.action' tool name that names the door being
// widened rather than a tool the model can call.
const (
	allowWidenTaskID = "host:fs-allow-widen"
	allowWidenTool   = "config.allow_dir"
)

// always is the 长期 verb (ticket 201 AC#4). It prints the rule text and hands the
// re-confirmation to the widening worker; it stores nothing by itself.
func (s *replySurface) always(corr string) (string, error) {
	card, ok := s.live.look(corr)
	if !ok {
		s.record("REFUSED", corr, "", "native/always", "本机没有这张卡的记录")
		return "", fmt.Errorf("没有找到这张卡 %q：长期落点必须由一张显示过的卡片推出，本机账上没有它", corr)
	}
	if card.Level != "L2" {
		s.record("REFUSED", corr, card.Tool, "native/always", "该路线没有允许动作，也就没有可存的长期规则")
		return "", fmt.Errorf("%s 是执行前阻止窗口，只有否决；「一直」是长期允许，得先有允许这一动作（SPEC-06 §2 B1）", card.Level)
	}
	if s.rt == nil || s.rt.mgr == nil {
		s.record("REFUSED", corr, card.Tool, "native/always", "配置写面未装配")
		return "", errors.New("配置写面没有装配：这一发能允许，但无处可存长期规则，所以什么都不写")
	}
	dir, rule, ok := card.WideningRule()
	if !ok {
		s.record("REFUSED", corr, card.Tool, "native/always",
			fmt.Sprintf("可存规则不唯一或不可用 paths=%d", len(card.Paths)))
		return "", fmt.Errorf("这张卡涉及 %d 个目录（影响路径：%v）；"+
			"「一直」只能存一条明确的规则，目录不唯一就不存，请改用 yes 只放行这一次",
			len(card.Paths), card.Paths)
	}
	rule = "把这一发存成长期规则会放宽 [fs]：" + rule

	// The card the operator is answering right now is NOT retired: it still has
	// its own answer to give, and the widening card is a separate question with a
	// separate id. The ledger keeps both.
	root := observe.NewRootFrom(s.ctx, "approval-always")
	s.record("STAGED", corr, card.Tool, "native/always", rule)
	observe.Default.Spawn("approval-waiter", "approval", root, func(ctx context.Context) {
		defer root.Cancel()
		s.runWidening(ctx, corr, dir, rule)
	})
	return "要存的规则：" + rule + "\n  这会再起一张 L2 卡（第二重确认，PLAN.md:1645-1646），" +
		"卡上会印自己的编号；同意它才真写进 config.toml，" +
		"不答按 C18 超时落拒绝、什么都不写。" +
		"注意：写入后由下一次启动读取，本次运行仍按现有 [fs] 名单判定。", nil
}

// runWidening raises the second card, waits for it off the listener's goroutine,
// and writes only on a native allow. Every branch books one audit line, because
// 「存了没有」 must be a reading and not an inference.
func (s *replySurface) runWidening(ctx context.Context, fromCorr, dir, rule string) {
	revoke := s.rt.gate.AdmitTextTask(allowWidenTaskID)
	defer revoke()
	ans, why := s.rt.gate.PendingApproval(ctx, tools.Decision{
		TaskID: allowWidenTaskID,
		Tool:   allowWidenTool,
		Level:  risk.L2,
		Reason: rule + "（这是「一直」要求的 L2 级重新确认；原来那张卡 " + fromCorr +
			" 的答复与这张无关，各答各的）",
		Paths: []string{dir},
		Mode:  s.rt.permissionMode(),
	})
	if ans != tools.AnswerAllow {
		if why == "" {
			why = string(ans)
		}
		s.record("WIDEN-REFUSED", fromCorr, allowWidenTool, "native/always",
			fmt.Sprintf("answer=%s dir=%q detail=%q", ans, dir, why))
		fmt.Fprintf(s.rt.stdout,
			"wisp run: 长期规则未存储（第二张卡的答案是 %s，最坏后果就是什么都不写）：%s\n", ans, why)
		return
	}
	if err := s.rt.mgr.AddAllowedDir(dir); err != nil {
		s.record("WIDEN-WRITE-FAILED", fromCorr, allowWidenTool, "native/always",
			fmt.Sprintf("dir=%q err=%v", dir, err))
		fmt.Fprintf(s.rt.stdout, "wisp run: 第二张卡批准了，但写盘失败，配置一个字节都没改：%v\n", err)
		return
	}
	s.record("WIDEN-APPLIED", fromCorr, allowWidenTool, "native/always",
		fmt.Sprintf("dir=%q section=fs key=allowed_dirs effect=next-start", dir))
	fmt.Fprintf(s.rt.stdout,
		"wisp run: 长期规则已写入 config.toml 的 [fs] allowed_dirs：%s（下一次启动生效，本次运行仍按旧名单）\n", dir)
}

// permissionMode reads the档 this run is screening under, so the widening card
// carries the mode the audit trail expects (tools.Decision.Mode's own comment:
// 「why was I not asked?」 has to be answerable from the record). An unreadable
// mode is booked as the strictest档 rather than the safest-sounding one, which is
// config.PermissionMode's own fail-closed reading.
func (rt *agentRuntime) permissionMode() risk.Mode {
	if rt == nil || rt.mgr == nil {
		return risk.DefaultMode()
	}
	return rt.mgr.Config().PermissionMode()
}

// bookWaitingState is ticket 201 AC#6's production caller: every time this run
// shows or retracts a card, it asks the seam what the D43 state of this instant
// is and books the answer.
//
// WHY A LINE IN THE LEDGER IS THE HONEST FORM. The 20-state visual for both of
// these rows exists (internal/ball/statevisual.go:177 renders Confirming, :181
// renders AwaitingApproval) but the ball is built only by cmd/balldebug, and this
// binary links no WebView2 host (internal/panel/pump.go:16, ticket 33 unclaimed) -
// so the producer this ticket owes is the one that NAMES the state, not one that
// draws it. The name comes out of the frozen D43 table via statemachine, never a
// new coinage, and a display surface or a log reader can act on it. The panel's
// own readable answer is the snapshot's existing "pending" key, which run.go
// republishes on exactly these two events.
func (rt *agentRuntime) bookWaitingState(via string) {
	if rt == nil || rt.liveCards == nil {
		return
	}
	st, awaiting := rt.liveCards.waitingState()
	if !awaiting {
		rt.auditf("approval: WAITING-STATE state=none awaiting=false tracked=%d producer=%q",
			len(rt.liveCards.pending()), via)
		return
	}
	card, _ := rt.liveCards.h.AwaitingHuman()
	rt.auditf("approval: WAITING-STATE state=%q awaiting=true corr=%q tool=%s level=%s paths=%d producer=%q",
		st, card.CorrelationID, card.Tool, card.Level, len(card.Paths), via)
}

// waitingStateName renders the same reading for any caller that wants the word
// rather than the audit line. It exists so the day a ball leg is wired into this
// assembly, the call it needs is already exported from this package's own ledger
// and not re-derived from a queue it cannot see.
func (rt *agentRuntime) waitingStateName() (statemachine.State, bool) {
	if rt == nil || rt.liveCards == nil {
		return "", false
	}
	return rt.liveCards.waitingState()
}
