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
	"io"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
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

// widenSeat is everything the 长期 branch needs from whoever owns it, and the
// reason it is a VALUE instead of a read through replySurface.rt.
//
// Ticket 219 AC#1: the resident process's double-click shape assembles no task
// pipeline, so its answer surface has no *agentRuntime behind it - and until this
// leg that also meant the tray's 长期允许 item could not reach this verb at all
// (the tray's items were not even wired in that shape). The seat lets the SAME
// verb body run over the objects that DO exist there: the resident gate, the
// boot's config Manager, the process's own stdout and the档 its panel chain
// screens under. Nothing in this struct is a route or a decision - it is four
// handles, so a second copy of the widening logic (the thing this file's header
// refuses) never has to be written for the second host.
type widenSeat struct {
	gate *approval.Gate
	// mgr is a READER, not a value: on the run seam the Manager is assigned
	// during assembly, after the surface exists (approval_reply.go's
	// newReplySurface is called from the same function that builds it), so
	// capturing it by value here would capture the nil it had at that moment.
	mgr   func() *config.Manager
	out   io.Writer
	mode  func() risk.Mode
}

// widenSeatOf answers "who owns the 长期 branch for this surface": the seat a
// host handed in, else the run this surface belongs to, else nil (a surface with
// neither refuses the verb - see always()).
func (s *replySurface) widenSeatOf() *widenSeat {
	if s == nil {
		return nil
	}
	if s.seat != nil {
		return s.seat
	}
	if s.rt == nil {
		return nil
	}
	rt := s.rt
	return &widenSeat{
		gate: rt.gate,
		mgr:  func() *config.Manager { return rt.mgr },
		out:  rt.stdout,
		mode: rt.permissionMode,
	}
}

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
	seat := s.widenSeatOf()
	if seat == nil || seat.gate == nil || seat.mgr == nil || seat.mgr() == nil {
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
		s.runWidening(ctx, seat, corr, dir, rule)
	})
	return "要存的规则：" + rule + "\n  这会再起一张 L2 卡（第二重确认，PLAN.md:1645-1646），" +
		"卡上会印自己的编号；同意它才真写进 config.toml，" +
		"不答按 C18 超时落拒绝、什么都不写。" +
		"注意：写入后由下一次启动读取，本次运行仍按现有 [fs] 名单判定。", nil
}

// runWidening raises the second card, waits for it off the listener's goroutine,
// and writes only on a native allow. Every branch books one audit line, because
// 「存了没有」 must be a reading and not an inference.
//
// seat is a parameter (ticket 219 AC#1): the caller already resolved which gate,
// write door, stdout and档 own this widening, so this body never has to know
// whether it is running inside a task pipeline or in the resident process's
// no-pipeline shape.
func (s *replySurface) runWidening(ctx context.Context, seat *widenSeat, fromCorr, dir, rule string) {
	revoke := seat.gate.AdmitTextTask(allowWidenTaskID)
	defer revoke()
	// A seat with no stdout is legal (a host that only has the audit sink): the
	// sentences below are duplicates of the record lines, so discarding them
	// loses no fact a reader cannot get from the ledger.
	out := seat.out
	if out == nil {
		out = io.Discard
	}
	mode := risk.DefaultMode()
	if seat.mode != nil {
		mode = seat.mode()
	}
	ans, why := seat.gate.PendingApproval(ctx, tools.Decision{
		TaskID: allowWidenTaskID,
		Tool:   allowWidenTool,
		Level:  risk.L2,
		Reason: rule + "（这是「一直」要求的 L2 级重新确认；原来那张卡 " + fromCorr +
			" 的答复与这张无关，各答各的）",
		Paths: []string{dir},
		Mode:  mode,
	})
	if ans != tools.AnswerAllow {
		if why == "" {
			why = string(ans)
		}
		s.record("WIDEN-REFUSED", fromCorr, allowWidenTool, "native/always",
			fmt.Sprintf("answer=%s dir=%q detail=%q", ans, dir, why))
		fmt.Fprintf(out,
			"wisp run: 长期规则未存储（第二张卡的答案是 %s，最坏后果就是什么都不写）：%s\n", ans, why)
		return
	}
	if err := seat.mgr().AddAllowedDir(dir); err != nil {
		s.record("WIDEN-WRITE-FAILED", fromCorr, allowWidenTool, "native/always",
			fmt.Sprintf("dir=%q err=%v", dir, err))
		fmt.Fprintf(out, "wisp run: 第二张卡批准了，但写盘失败，配置一个字节都没改：%v\n", err)
		return
	}
	s.record("WIDEN-APPLIED", fromCorr, allowWidenTool, "native/always",
		fmt.Sprintf("dir=%q section=fs key=allowed_dirs effect=next-start", dir))
	fmt.Fprintf(out,
		"wisp run: 长期规则已写入 config.toml 的 [fs] allowed_dirs：%s（下一次启动生效，本次运行仍按旧名单）\n", dir)
}

// permissionMode reads the档 this run is screening under, so the widening card
// carries the mode the audit trail expects. tools.Decision.Mode's own comment
// spells the requirement out: the field is 「the permission mode this call was
// screened under」, answerable from the record rather than from config.toml's
// current contents, because 「why was I not asked?」 is the question it exists
// for. The D36 re-confirmation card next door (config_reload.go) books this same
// read.
//
// Ticket 377: it asks the STORE, not the config. perm.Store.PermissionMode is
// the EFFECTIVE档 - the read side ticket 372 collapses to ask_every_step when
// this process was handed no confirmation channel - while
// config.Config.PermissionMode answers config.toml's value. On the `wisp run`
// seam the two are equal by construction (confirm is the confirmModeSwitch
// method value, never nil, so 372's collapse cannot fire); they come apart on
// the `wisp panel-inbound` seam, which has no leg at all, and on any future
// assembly that forgets one. Reading the file there would book a档 this process
// never ran - a record that says "全自动" about a run that asked about
// everything.
//
// An unreadable mode is booked as the strictest档: with no store there is no
// evidence a channel exists, and the store itself already fail-closes to
// ask_every_step when it cannot answer (perm.Store.PermissionMode).
func (rt *agentRuntime) permissionMode() risk.Mode {
	if rt == nil {
		return risk.DefaultMode()
	}
	if rt.modes != nil {
		return rt.modes.PermissionMode()
	}
	return risk.DefaultMode()
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
	// Ticket 228 AC#3: the same reading now DRIVES the orb when this process has
	// one. Until this line the function was the producer of a name and nothing
	// else - the ball entered Confirming/AwaitingApproval only because the card
	// surface happened to post the state itself, so "the state was driven" had
	// no witness a mutation could red.
	//
	// Order: the audit line is booked first, on purpose. The record has to exist
	// even if the projection panics or a host's sink is the thing that fails
	// next, which is the same reason consoleApprovalUI.Prompt books before it
	// prints the card.
	//
	// Only the entering direction is projected: with nobody awaiting, this
	// function has no state to name (the line above says so), and the orb's exit
	// edge belongs to the surface that settles the card (ballCardUI.settleOrb
	// reads the same ledger). A second writer for that edge would be two truth
	// sources for "what the orb does now".
	rt.projectWaitingState(st)
}

// waitingStateSink is ticket 228 AC#3's projection signature: the D43 name this
// run just read, handed to a host that holds a surface able to show it. It is a
// named type so run.go's field does not pull internal/statemachine into a file
// that has no other business with the state vocabulary.
type waitingStateSink func(statemachine.State)

// bindWaitingStateSink hands this run the ball projection of its waiting state
// (ticket 228 AC#3). The resident leg calls it once, immediately after the
// runtime it assembled owns the card ledger, and it is the only writer.
func (rt *agentRuntime) bindWaitingStateSink(fn waitingStateSink) {
	if rt == nil {
		return
	}
	rt.waitingMu.Lock()
	defer rt.waitingMu.Unlock()
	rt.waitingSink = fn
}

// projectWaitingState drives one D43 waiting state onto whichever surface this
// host bound. nil sink is legal and means what it means in the rest of this
// assembly: nobody handed this run a ball, so the state is named in the audit
// line above and goes no further.
func (rt *agentRuntime) projectWaitingState(s statemachine.State) {
	if rt == nil || s == "" {
		return
	}
	rt.waitingMu.Lock()
	sink := rt.waitingSink
	rt.waitingMu.Unlock()
	if sink != nil {
		sink(s)
	}
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
