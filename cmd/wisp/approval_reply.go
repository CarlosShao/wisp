package main

// Ticket 201 - the reply listener: the answer side of an approval, wired into
// the assembly that actually runs.
//
// WHY THIS FILE EXISTS. The 09-29 census of the approval reply surface
// (docs/evidence/s1/219-approval-reply-surface-c1.md §②A) measured one
// structural fact and it is the whole subject of this file: the ASKING side is
// live in production - internal/tools/bridge.go calls Gate.PendingWindow for an
// L1 verdict and Gate.PendingApproval for an L2 one, and cmd/wisp's own mode
// switch raises a host-side L2 card - while the ANSWERING side,
// Gate.DecideFromNative / Gate.DecideFromPanel / Gate.Veto, had zero callers
// outside _test.go at census time. Every card this process showed was therefore
// answered by nobody: an L2 card waited out the full C18 deadline and
// auto-rejected, and an L1 window could not be opposed at all. 票 201 :10 names
// the same hole from the other end, row verbatim: | 三种"人能答复"的入口 | **生产零调用者** |.
// This file (with the Replies routes it drives) is what closed it - .Veto,
// .DecideFromNative and .DecideFromPanel all have production call sites today.
//
// WHAT WAS BUILT, and what was deliberately NOT. The census's second
// conclusion is the design constraint: the correct shape of an answer already
// exists, as the two faces approval.Gate publishes (NativeAPI at ui.go:143 and
// PanelAPI at ui.go:155, handed out by Gate.Native() and Gate.Panel()), plus the
// two routers (DecideFromNative / DecideFromPanel). So this file opens NO third
// interface and adds NO authority. It assembles the faces that exist, feeds them
// from the one human surface this tree really has in a running process (the
// operator's console), and books what came out.
//
//   - the allow direction lives only on the native face, and only with the
//     single-use grant that UI.Prompt was handed. PanelItem has no grant field
//     and PanelAPI has no Allow method (ui.go:45-48 spells that out), so the
//     untrusted route cannot carry an allow even if this file wanted to send
//     one; DecideFromPanel then refuses one on the ROUTE alone. That is
//     AGENTS.md §1.2's ban #6 ("allow decisions are native-side only") and it is
//     not loosened here, in either direction, for convenience.
//   - the panel route gets 拒绝 and 查看卡片 and nothing else, which is exactly
//     SPEC-06 §9 layer 3's sentence.
//   - an unanswered L2 still resolves to a REJECT (queue.expire) and this file
//     adds a line that books it; an expired L1 window still EXECUTES
//     (SPEC-06 §2 row L1, pinned green by TestComposedGateBlocksAWriteForTwo
//     Seconds). Changing either polarity is a D4 / SPEC-06 §2 contract act,
//     which ticket 201's AC#2 asks for and which no implementation leg may do
//     quietly - it is registered as an open question in this leg's evidence file
//     instead (docs/evidence/s1/201-reply-listener-r1.md §⑥).
//
// WHY A CONSOLE LISTENER AT ALL, when 票 201 names the ball and the tray first:
// `cmd/wisp run` is the only production assembly in this tree that constructs an
// approval.Gate at all (measured: approval.New has one non-test call site,
// run.go:406). The floating ball is built only by cmd/balldebug, and
// cmd/wisp/resident_windows.go runs an event loop with no gate, no bridge and no
// task - so a ball click handed to this gate today would be a listener with no
// questioner, in a process that never asks. The veto CHANNELS are the honest
// limit of what a console can claim: SPEC-06 §2's four are 单击球 / Esc / KWS
// 否决词 / 面板拒绝, and none of them is a ball, a global hook, a panel or a
// speech model. Ticket 170 (form 甲, owner-approved) changed the EMPTY half of
// that sentence and nothing else: when interactiveStdin proved the process
// really owns a console input buffer, cmdRun declares the cancel-KEY slot
// (esc) as this host's transport and names the key it really holds through the
// read-time spelling seam (「按 回车 键」) - a terminal's cancel key IS the Enter
// the operator hits to feed this stream. A non-TTY stdin keeps the empty
// declaration, and run.go books CLI-VETO-NONE so the posture is said, not
// implied. Either way this file calls Gate.Veto with the channel the HOST
// declared and lets the gate's own ChannelRegistry answer - it invents no fifth
// channel and claims no produced-but-unloaded one (B1).

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// The route labels this assembly stamps onto an approval.Request.Source.
//
// They are LABELS, not authority: gate.go:594-600 says in as many words that no
// branch in that package reads Source, because a security decision keyed on a
// caller-supplied selector is the M-7/C-3 failure shape. They exist so the audit
// line can hold "what the transport claimed" next to "what the route proved",
// which is the only use a claimed source has in this design.
const (
	nativeReplySource = "cmd-wisp-console-native"
	panelReplySource  = "cmd-wisp-console-panel-route"
)

// liveCard is one confirmation this process displayed, as the answering side
// needs it: where it is, what it is, and - native side only - the proof an allow
// has to present. It is the seam's own card type (ticket 201 moved the ledger and
// the route-choosing into internal/agent/approval/replies.go so a non-console
// host has something to be handed); the alias is kept because this file's
// sentences name its fields.
type liveCard = approval.ReplyCard

// nativeCards is the ledger the native surface fills and the reply listener
// reads. It exists because the grant cannot travel any other way: the answer
// arrives on a different goroutine from the one that displayed the card, so the
// surface that received the proof has to hold it until someone spends it.
//
// It is now a thin binding onto approval.Replies rather than a second ledger, and
// the reason is this file's own rule: a second copy of the grant-spending and
// route-choosing logic is how one of them grows a leniency the other does not
// have. The methods below keep their lowercase names because they are the
// console's vocabulary; the authority lives one package over.
type nativeCards struct{ h *approval.Replies }

func newNativeCards() *nativeCards { return &nativeCards{h: approval.NewReplies()} }

// bind attaches the ledger to the Gate whose UI it was injected into and, on the
// reply leg, to the veto channel this host really wired. Empty channel and empty
// source labels are legal: the gate's own registry then refuses an undeliverable
// veto, and the seam stamps its own transport names in the audit.
func (n *nativeCards) bind(g *approval.Gate, vetoChannel approval.Channel, nativeSource, panelSource string) {
	if n == nil || n.h == nil {
		return
	}
	n.h.Attach(approval.HostBinding{
		Gate: g, VetoChannel: vetoChannel,
		NativeSource: nativeSource, PanelSource: panelSource,
	})
}

// record books one displayed card straight off the Prompt the gate handed this
// surface, so a field a card carries cannot be forgotten by whoever is typing the
// literal (ticket 201 AC#4 needs Paths for the rule text, and the console leg had
// no use for them before).
func (n *nativeCards) record(p approval.Prompt) {
	if n == nil || n.h == nil {
		return
	}
	n.h.Record(p)
}

// look reads one card back. A missing entry means this surface never displayed
// it, or it already left the screen - in both cases there is no grant to spend.
func (n *nativeCards) look(corr string) (liveCard, bool) {
	if n == nil || n.h == nil {
		return liveCard{}, false
	}
	return n.h.Look(corr)
}

// forget drops one entry (answered, dismissed, or handed off to execution).
func (n *nativeCards) forget(corr string) {
	if n == nil || n.h == nil {
		return
	}
	n.h.Forget(corr)
}

// pending lists what this host is holding, oldest displayed first.
func (n *nativeCards) pending() []liveCard {
	if n == nil || n.h == nil {
		return nil
	}
	return n.h.Pending()
}

// waitingState is ticket 201 AC#6's producer: the D43 state name this instant
// deserves, read off the live queue and this host's own ledger rather than
// scripted. ("", false) means nobody is being waited on.
func (n *nativeCards) waitingState() (statemachine.State, bool) {
	if n == nil || n.h == nil {
		return "", false
	}
	return n.h.WaitingState()
}

// replySurface is the host's answer side, assembled from the faces the gate
// already publishes. It is the object a native surface (today the console,
// tomorrow the ball's buttons and the tray menu) and a panel host are handed,
// and the method set is the whole boundary:
//
//	allow / session / reject / veto   the native route, approval.Gate.Native()
//	                           faces and the native routers; allow and session
//	                           both need the grant. `session` is ticket 224's
//	                           D45-2 third answer: the same release, plus one
//	                           approval_grant row per path the card named.
//	panelReject / panelAllow   the panel route, DecideFromPanel only.
//	head / view                the panel's read side (PanelAPI), no grant, no allow.
//	revoke                     the 长期 rule's undo (ticket 226 AC#5, approval_revoke.go):
//	                           no card, no grant - it only ever SHRINKS the
//	                           allowlist, through Manager.SetAllowedDirs.
//
// There is no method here that lets a panel-sourced decision execute anything,
// and no path that mints or forwards a grant to the panel route. That is not a
// check this type remembers to run; it is the absence of a door. The same holds
// for the REMEMBERED variant: `session` reaches NativeAPI.AllowSession, and
// PanelAPI has no Allow method for a session-flavoured sibling to attach to.
type replySurface struct {
	gate  *approval.Gate
	live  *nativeCards
	audit func(format string, args ...any)
	// rt is the run this surface belongs to: the 长期 branch needs its config
	// write door, its stdout and its audit sink (approval_always.go), and a
	// surface that could not name the run it answers for could not persist
	// anything either.
	rt *agentRuntime
	// seat is WHO OWNS the 长期 branch's three objects (gate, config write door,
	// stdout, and the档 the widening card is booked under) when that owner is
	// not a run. Ticket 219 AC#1 needs exactly this: the resident process's
	// double-click shape assembles no task pipeline at all, so there is no
	// *agentRuntime to hand a surface - yet the two answer items that DO work
	// there (允许一次 / 拒绝) and the third one (长期允许) all belong to the same
	// verb table. nil means "derive it from rt", which is the console shape.
	// See widenSeat and answerSurfaceWithoutPipeline.
	seat *widenSeat
	// review is the operator's own copy of what was answered, newest last
	// (ticket 219 AC#5's third destination: 理由要能回看). A surface built by a
	// case that does not care leaves it nil and record() then books the audit
	// line only - no invented history.
	review *replyReview
	// vetoChannel is the channel this host's cancel transport really is. Empty
	// means this assembly wired none of SPEC-06 §2's four, and Gate.Veto's own
	// ChannelRegistry says so back.
	vetoChannel approval.Channel
	// ctx is the listener's own context; it carries no authority.
	ctx context.Context
	// parked is the listener's own statement of WHERE it is (ticket 228 AC#4):
	// true while it sits inside its blocking read on the operator's stream, where
	// a cancelled context cannot reach it, and false while it is handling a line.
	// It exists so the run's teardown (agentRuntime.close) can tell "nothing is in
	// flight, no cancel can arrive" from "a line is being handled" - the pre-ticket
	// comment asserted the first case from outside and nothing measured it.
	parked atomic.Bool
}

// noteParked is the listener telling its owner where it is. Called on both sides
// of every read by runReplyLoop and nowhere else: the loop is the only thing that
// knows, so a second writer would be a second truth source.
func (s *replySurface) noteParked(p bool) {
	if s == nil {
		return
	}
	s.parked.Store(p)
}

// parkedInRead reports whether the listener is inside its blocking read right
// now. false also covers "this surface has no listener" (a host that only built
// the verb table, e.g. the resident leg's tray answers), which is honest: there
// is then nothing parked to name.
func (s *replySurface) parkedInRead() bool {
	return s != nil && s.parked.Load()
}

// allow spends this card's native grant through the native router, via the seam
// (internal/agent/approval/replies.go). This file formats the sentence; the
// grant lookup, the route choice and the ledger write are the seam's.
func (s *replySurface) allow(corr string) (string, error) {
	card, _ := s.live.look(corr)
	err := s.live.h.Allow(s.ctx, corr)
	switch {
	case errors.Is(err, approval.ErrNoTrackedCard):
		s.record("REFUSED", corr, "", "native/allow", "本机没有这张卡的记录")
		return "", fmt.Errorf("没有找到待答复的卡片 %q：它可能已经结束，从未显示的卡片这里也没有令牌可花", corr)
	case errors.Is(err, approval.ErrRouteHasNoAllow):
		s.record("REFUSED", corr, card.Tool, "native/allow", "该路线没有允许动作")
		return "", fmt.Errorf("%s 是执行前阻止窗口，它只有否决、没有允许（SPEC-06 §2 B1）；"+
			"要停下它就答 veto %s", card.Level, corr)
	case err != nil:
		s.record("REFUSED", corr, card.Tool, "native/allow", err.Error())
		if errors.Is(err, approval.ErrBadGrant) {
			return "", fmt.Errorf("原生令牌无效（缺失/已用/与本次请求不绑定）：%w；"+
				"这张卡需要重新显示才能被允许", err)
		}
		return "", err
	}
	s.record("ANSWERED", corr, card.Tool, "native/allow", "原生侧允许")
	return "已允许 " + corr + "（" + card.Tool + "）：调用已放行，这一发会执行", nil
}

// session answers with D45-2's third option, 「本会话内允许」 (ticket 224).
//
// It is a separate verb rather than a flag on `yes` because the two answers
// store different things and the operator has to be able to tell them apart from
// the transcript: `yes` releases one call, `session` releases one call AND
// writes one approval_grant row per path the card named, keyed to the identity
// this process minted at boot.
//
// What it can cover is narrower than the button, and the contract is why:
// SPEC-06 §8.3's first bullet ("L2 永不进入任何持久授权（含会话级）") means a
// stored rule is consulted for an L1 question and never for an L2 one. This
// console hands out L2 cards, so the honest reading of `session` on an L2 card
// is "allow this call, and remember the scope for the L1 questions this session
// would otherwise ask about the same tool and path". It is NOT "never ask about
// this path again": the next L2 card for that path still arrives, and the line
// below says so rather than letting the verb over-promise.
//
// The grant is written through the same native route as `yes` and needs the same
// single-use nonce, so AGENTS.md §1.2's ban #6 holds for the remembered scope
// exactly as it holds for the immediate one: no panel-shaped caller reaches this
// method, and there is no panel-shaped method for it to reach.
func (s *replySurface) session(corr string) (string, error) {
	card, _ := s.live.look(corr)
	err := s.live.h.AllowSession(s.ctx, corr)
	switch {
	case errors.Is(err, approval.ErrNoTrackedCard):
		s.record("REFUSED", corr, "", "native/allow-session", "本机没有这张卡的记录")
		return "", fmt.Errorf("没有找到待答复的卡片 %q：会话授权必须由一张显示过的卡片推出，"+
			"本机账上没有它，也就没有任何东西可记", corr)
	case errors.Is(err, approval.ErrRouteHasNoAllow):
		s.record("REFUSED", corr, card.Tool, "native/allow-session", "该路线没有允许动作")
		return "", fmt.Errorf("%s 是执行前阻止窗口，它只有否决、没有允许（SPEC-06 §2 B1）；"+
			"要停下它就答 veto %s", card.Level, corr)
	case err != nil:
		s.record("REFUSED", corr, card.Tool, "native/allow-session", err.Error())
		if errors.Is(err, approval.ErrBadGrant) {
			return "", fmt.Errorf("原生令牌无效（缺失/已用/与本次请求不绑定）：%w；"+
				"能记下授权的只有显示过这张令牌的那一发，重新显示后再答", err)
		}
		return "", err
	}
	s.record("ANSWERED", corr, card.Tool, "native/allow-session",
		fmt.Sprintf("原生侧允许并记入本会话，paths=%d", len(card.Paths)))
	return "已按「本会话内允许」答复 " + corr + "（" + card.Tool + "）：" +
		"这一发已放行，卡片上那些路径对本会话后续的 L1 询问不再重复提问" +
		"（L2 永不被会话授权覆盖，SPEC-06 §8.3 第 1 条）", nil
}

// reject refuses through the native router. The reason is not decoration: the
// queue hands it to the waiting call and internal/tools/bridge.go passes it back
// to the model verbatim (orDefault(why, …) on the reject branches), which is the
// channel 票 219's reason box is going to need.
func (s *replySurface) reject(corr string, reason replyReason) (string, error) {
	return s.refuse(corr, reason, false)
}

// refuse is the shared body of the two reject routes, kept as one function
// because the ONLY difference between them is which router the decision travels
// through - and a second copy is how one of them grows a leniency the other does
// not have.
func (s *replySurface) refuse(corr string, reason replyReason, viaPanel bool) (string, error) {
	card, _ := s.live.look(corr)
	route, routeText := "native", "原生侧"
	var err error
	if viaPanel {
		route, routeText = "panel", "面板路线"
		err = s.live.h.PanelReject(s.ctx, corr, reason.text)
	} else {
		err = s.live.h.Reject(s.ctx, corr, reason.text)
	}
	if err != nil {
		s.record("REFUSED", corr, card.Tool, route+"/reject", err.Error())
		return "", err
	}
	s.record("ANSWERED", corr, card.Tool, route+"/reject", reason.text)
	if reason.text == "" {
		return "已拒绝 " + corr + "（" + routeText + "，未填理由）：未执行", nil
	}
	if note := reason.truncationNote(); note != "" {
		return "已拒绝 " + corr + "（" + routeText + "）：未执行，理由已随结果回给模型（" + note + "）", nil
	}
	return "已拒绝 " + corr + "（" + routeText + "）：未执行，理由已随结果回给模型", nil
}

// panelReject is the panel route's refusal - SPEC-06 §9 layer 3's 「拒绝」
// button, and the only decision the page may send.
func (s *replySurface) panelReject(corr string, reason replyReason) (string, error) {
	return s.refuse(corr, reason, true)
}

// panelAllow exists so the panel route's refusal is reachable, not to soften it.
// A WebView host that is handed {allow:true} for approval.decide forwards it here
// and gets back ErrPanelAllow plus a burned nonce: the answer is refused on the
// ROUTE, before any grant is looked at, and a grant that surfaced on this path is
// treated as leaked (gate.go's DecideFromPanel branch, queue.revokeGrants).
// Nothing on this method's happy path exists, because it has none.
func (s *replySurface) panelAllow(corr string) (string, error) {
	card, _ := s.live.look(corr)
	offered, err := s.live.h.PanelAllow(s.ctx, corr)
	s.record("PANEL-ALLOW-REFUSED", corr, card.Tool, "panel/allow",
		fmt.Sprintf("claimed_source=%q grant_offered=%v err=%v", panelReplySource, offered, err))
	if err == nil {
		// Unreachable while gate.DecideFromPanel refuses on the route. Kept as a
		// loud fault rather than a silent success: if that branch ever stops
		// refusing, the sentence this function prints would be the bug report.
		return "", errors.New("安全故障：面板路线的「允许」没有被拒绝，请立即停住这条通路")
	}
	if offered {
		// The cost of ringing this door is stated, not hidden: a grant that
		// appeared on the untrusted route is treated as leaked, so the nonce this
		// card was holding is gone. The ledger entry STAYS, because the card is
		// still on screen and still pending - what a native attempt now meets is
		// ErrBadGrant, which is the reading this path exists to produce: a token
		// that passed through the panel route is not spendable anywhere else.
		return "", fmt.Errorf("面板路线不得允许（F2 第三层）：%w；"+
			"这张卡的一次性令牌已按泄露处理烧掉，要允许它得重新显示", err)
	}
	return "", fmt.Errorf("面板路线不得允许（F2 第三层）：%w", err)
}

// veto cancels an L1 window through the gate's veto funnel, on the channel THIS
// host declared. It never claims a channel the assembly did not wire: with an
// empty host channel the gate's own registry answers 「未知取消通道」, which is a
// refusal the operator hears rather than a name this file invented.
func (s *replySurface) veto(corr string) (string, error) {
	card, ok := s.live.look(corr)
	if !ok {
		s.record("REFUSED", corr, "", "veto", "本机没有这张卡的记录")
		return "", fmt.Errorf("没有找到待否决的卡片 %q：窗口可能已经结束", corr)
	}
	err := s.live.h.Veto(corr)
	if err != nil {
		s.record("REFUSED", corr, card.Tool, "native/veto", err.Error())
		return "", fmt.Errorf("否决未能送达：%w", err)
	}
	s.record("ANSWERED", corr, card.Tool, "native/veto", string(s.vetoChannel))
	return "已否决 " + corr + "（" + card.Tool + "）：窗口已取消，未执行", nil
}

// head reads the queue's displayed item through the panel-safe projection: what
// the page may show, which by construction carries no grant and no allow.
func (s *replySurface) head() (string, error) {
	item, ok := s.live.h.Head()
	if !ok {
		return "队头没有待审批的卡片（当前深度 " + fmt.Sprint(s.gate.Queue().Depth()) + "）", nil
	}
	return panelItemLine(item), nil
}

// view reads one card the same way, by the name the operator was given.
func (s *replySurface) view(corr string) (string, error) {
	item, ok := s.live.h.View(corr)
	if !ok {
		return "", fmt.Errorf("卡片 %q 不在待审批队列里（已结束、已作废，或从来不是 L2 卡）", corr)
	}
	return panelItemLine(item), nil
}

// panelItemLine renders approval.PanelItem for the console. The last sentence is
// not filler: 「显示卡片」 is half of what the web side is allowed to do, and the
// projection's whole security claim is what it does NOT contain (ui.go:45-48), so
// the reader is told which fields this surface can never see.
func panelItemLine(item approval.PanelItem) string {
	reason := item.Reason
	if strings.TrimSpace(reason) == "" {
		reason = "（判定未给出理由，卡片必须明说它缺信息）"
	}
	return fmt.Sprintf("卡片 %s｜%s｜%s｜深度 %d｜理由：%s｜影响路径 %d 条"+
		"｜此投影不含令牌、不含允许（PanelAPI 没有 Allow 方法）",
		item.CorrelationID, item.Level, item.Tool, item.Depth, reason, len(item.Paths))
}

// record books one answered-or-refused reply to the run's audit trail.
//
// The family is "approval: REPLY ..." and it is deliberately separate from the
// approval layer's own ANSWER-* lines: those say what a settled card resolved to,
// this says what a host surface was ASKED to do and whether it got that far - so
// an unanswered question, a refused answer and an applied one are three different
// readings instead of one inference. Reason text travels verbatim (observe's
// handler bounds it), because 「谁以什么理由答的」 is the sentence the ledger owes
// an auditor and 票 219's reason box is built on top of it.
//
// when=<RFC3339> is ticket 219 AC#2's 什么时候: the widening's WIDEN-APPLIED line
// must let an auditor read WHO (corr=), WHEN (when=), WHICH LINE (dir= inside
// detail=) and WHICH SOURCE (route=) off one line. The durable JSONL copy was
// already stamped by the slog handler; the console copy carries its own stamp
// now, so both copies answer the same four questions.
//
// TICKET 219 AC#5's third destination hangs off THIS function, on purpose: a
// review list filled anywhere else would be a second truth source about what was
// answered, and the two would drift the first time a route was added. The audit
// line stays the durable record; review is the operator's readable copy of it.
func (s *replySurface) record(kind, corr, tool, route, detail string) {
	if s == nil {
		return
	}
	when := time.Now()
	if s.audit != nil {
		s.audit("approval: REPLY %s corr=%q tool=%s route=%q detail=%q when=%s",
			kind, corr, tool, route, detail, when.Format(time.RFC3339))
	}
	s.review.add(replyDecision{
		when: when, kind: kind, corr: corr, tool: tool, route: route, detail: detail,
	})
}

// decisions is the operator's own copy: what this process answered, in what
// order, and with what reason (ticket 219 AC#5's third destination - 「今天它替
// 我做了哪些决定」). Before this verb the reason had exactly two readers, the
// audit sink and the model, and the human who typed it had no way back to it:
// tool_call's thirteen frozen columns carry no reason (internal/memory/schema.go)
// and a reply is not a tool call, so the only honest read face available in this
// tree is the answer surface's own list.
//
// It says what it is and is not: the list is IN-PROCESS and bounded
// (replyReviewKeep), the durable copy is the audit family, and it never
// re-derives a decision from the queue - a card still waiting is not in here,
// which is the difference between "what it decided" and "what it is asking".
func (s *replySurface) decisions() (string, error) {
	if s == nil || s.review == nil {
		return "本进程没有答复回看账：这个装配没有接答复面，也就没有可回看的理由", nil
	}
	all := s.review.snapshot()
	if len(all) == 0 {
		return "本进程还没有答复过任何卡片（空账不是错误：一张都没答过）", nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "本进程最近 %d 条答复（新的在前；落盘的证据是日志里 approval: REPLY 那一族行）：", len(all))
	for i := len(all) - 1; i >= 0; i-- {
		e := all[i]
		detail := strings.TrimSpace(e.detail)
		if detail == "" {
			detail = "（未填理由）"
		}
		fmt.Fprintf(&b, "\n  %s %s corr=%q tool=%s route=%q 理由/明细=%q",
			e.when.Format(time.RFC3339), e.kind, e.corr, e.tool, e.route, detail)
	}
	return b.String(), nil
}

// attachReplyListener puts the answer side on this run. Called from
// assembleRuntime when a reply source was handed in (production: an interactive
// console; the CLI tests: a scripted reader, which is the injection seam
// AGENTS.md §1.3 allows for `wisp run`).
//
// vetoChannel is the host's OWN statement of which of SPEC-06 §2's four channels
// its cancel transport really is, and a console run has none of them: no ball
// window, no global Esc hook, no panel host, no KWS model - which is also why
// run.go hands the gate an empty ChannelRegistry. Empty is the honest value:
// Gate.Veto then answers with its own ChannelError instead of this file naming a
// channel it cannot produce, and naming one would put a fabricated attribution
// into both the audit line and the reason the model reads back. The GUI leg that
// really owns the ball or the Esc hook (tickets 07/77/92) sets the channel AND
// marks it loaded through Gate.Channels() in the same step, because the two
// belong together: a loaded-but-unproduced channel over-promises on the card
// (B1's rule), a produced-but-unloaded one silently drops the user's click.
//
// It is a listener, not a poller: one goroutine blocked on the operator's stream
// costs nothing until a card is on screen, and the reply lands on the queue from
// there. It never answers on the user's behalf - there is no path through this
// file that produces a decision without a line arriving on in.
func (rt *agentRuntime) attachReplyListener(in io.Reader, vetoChannel approval.Channel) {
	if rt == nil || rt.gate == nil || in == nil {
		return
	}
	// The console leg owns its ledger, so it binds here. A host that handed its
	// gate AND its ledger in (ticket 246 AC#7) does not: re-attaching would stamp
	// this file's console transport names onto a ledger another surface filled,
	// and the audit line would then name the wrong transport for a real card.
	surface := rt.newReplySurface(nil, vetoChannel, rt.spec.gate == nil)
	// The goroutine goes through the registry like every other one in this
	// process (D22 ban #1): name from the D38 per-task roster, owner named, and
	// the recover boundary installed by Registry.run rather than a local recover.
	rt.replyHandle = observe.Default.Spawn("approval-waiter", "approval", rt.replyRoot,
		func(ctx context.Context) {
			runReplyLoop(ctx, surface, in, rt.stdout)
		})
	fmt.Fprintf(rt.stdout,
		"wisp run: 答复监听已接入（卡片上给了编号）。原生侧：yes <编号> / no <编号> [理由] / veto <编号>；"+
			"L1 确认窗口存活期间，敲任意内容或直接回车＝停（票 170）；"+
			"面板路线：panel-no <编号> [理由]、panel-yes <编号>（这一条只会被服务端 API 拒绝）、"+
			"head / view <编号> 看卡片；"+
			"revoke <规则或目录> 撤销一条长期规则（票 226，只减不加，help 看写法）；"+
			"quit 退出监听。\n")
}

// newReplySurface builds this run's answer side and returns it, WITHOUT starting
// a reader. attachReplyListener is the console's caller and does the spawn;
// ticket 246 AC#7's resident leg needs the same surface driven from a stream that
// also carries task texts, and a second verb table would be the second copy
// approval_reply.go's own rule refuses ("a second copy is how one of them grows
// a leniency the other does not have").
//
// root names the owner of the reply side's lifetime: nil keeps today's
// behaviour (a detached process-level root named "approval-reply"), and a host
// that must be cancellable from its own task root passes one derived from it.
// rebindLedger is the same question run.go asks: only the assembly that built
// the ledger gets to say which transport and which veto channel it answers on.
func (rt *agentRuntime) newReplySurface(root *observe.Root, vetoChannel approval.Channel, rebindLedger bool) *replySurface {
	if rt == nil || rt.gate == nil {
		return nil
	}
	if root == nil {
		root = observe.NewRoot("approval-reply")
	}
	rt.replyRoot = root
	surface := &replySurface{
		gate:        rt.gate,
		live:        rt.liveCards,
		audit:       rt.auditf,
		rt:          rt,
		vetoChannel: vetoChannel,
		ctx:         root.Ctx,
		review:      newReplyReview(),
	}
	rt.reply = surface
	// Bind the seam to this gate before anything can be answered: the ledger was
	// built at assembly (it has to be, because the UI that fills it is one of the
	// gate's own options) and it only learns the routers now. The host label
	// travels with the binding so the audit line names this transport, and the
	// veto channel travels with it too - see the paragraph above about the two
	// statements being made together.
	if rebindLedger {
		rt.liveCards.bind(rt.gate, vetoChannel, nativeReplySource, panelReplySource)
	}
	// Declaring a transport IS asserting it is up, so the two statements are made
	// here together rather than in two places a refactor can pull apart: the card
	// may now render 「按 Esc 键（已加载）」 because this host really will deliver
	// that cancel, and Gate.Veto's registry check will let it land.
	//
	// The host that injected its gate already made that statement - and made it
	// conditionally, on a ball window that Win32 really gave it
	// (resident_approval_windows.go's bindBallHost). Such a host passes the empty
	// channel here, so this line cannot load a channel that process does not own.
	if vetoChannel != "" {
		rt.gate.Channels().SetLoaded(vetoChannel, true)
		// Ticket 170: the armed posture is said, not implied - the mirror of the
		// CLI-VETO-NONE line run.go books when a host declares nothing.
		rt.auditf("approval: CLI-VETO-ARMED channel=%s（票 170：本装配声明了 L1 否决通道，"+
			"窗口存活期间输入即停）", vetoChannel)
	}
	return surface
}

// splitReplyLine is the console's line grammar, extracted from runReplyLoop by
// ticket 246 AC#7 so the resident leg's single console reader can use the SAME
// parse. It is two Cut calls and not an authority: which verbs mean what is
// replySurface.handle's business, and a host that parsed lines its own way would
// be a second copy of the one thing that decides whether an answer counts.
func splitReplyLine(line string) (verb, corr, arg string) {
	verb, rest, _ := strings.Cut(line, " ")
	corr, arg, _ = strings.Cut(strings.TrimSpace(rest), " ")
	return verb, corr, strings.TrimSpace(arg)
}

// runReplyLoop reads the operator's answers until the stream ends or the run's
// reply root is cancelled.
//
// Ticket 228 AC#4: the loop also reports WHERE it is (noteParked on both sides of
// every read), because its owner's teardown must join it and this is the one fact
// that decides whether a join can accomplish anything. The read itself is
// unchanged - same scanner, same buffer, same per-line ctx check - so the grammar
// and the refusals below are byte-for-byte what they were.
func runReplyLoop(ctx context.Context, s *replySurface, in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 0, 4096), 1<<16)
	for {
		s.noteParked(true)
		more := sc.Scan()
		s.noteParked(false)
		if !more {
			break
		}
		if ctx.Err() != nil {
			return
		}
		line := strings.TrimSpace(sc.Text())
		// Ticket 170 form 甲: while an L1 block window is live on this gate, ANY
		// input the operator feeds this stream is the console's 停 - a bare
		// Enter included, which is exactly the key the card names on this leg
		// (main.go's consoleCancelKey). quit keeps its control meaning: it stops
		// the listener, it does not answer a window. When no window is live the
		// line falls through to the verb grammar untouched, so an L2 reply typed
		// a moment too late is still parsed as before.
		if line != "quit" && !strings.HasPrefix(line, "quit ") {
			if s.vetoLiveL1Windows(out) {
				continue
			}
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		verb, corr, arg := splitReplyLine(line)
		if verb == "quit" {
			fmt.Fprintln(out, "wisp run: 答复监听已退出；卡片仍在等待，未答复的那一张按各自的超时处理")
			return
		}
		text, err := s.handle(verb, corr, arg)
		if err != nil {
			// A refused reply is printed and the loop keeps going: a typo must
			// never be a way to lose the ability to answer a live card.
			fmt.Fprintf(out, "wisp run: 答复未被接受（%s）：%v\n", line, err)
			continue
		}
		fmt.Fprintf(out, "wisp run: %s\n", text)
	}
	if err := sc.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(out, "wisp run: 答复输入流已断（%v），本轮之后无人能再答复\n", err)
	}
}

// vetoLiveL1Windows is ticket 170 form 甲's mechanism, and it is deliberately
// the dumbest possible reading of 说「停」: if the gate is holding an L1 block
// window open right now, whatever just arrived on the operator's stream cancels
// it - the operator does not have to remember a verb or a card number in a
// 2-second window, they just have to say stop. It cancels EVERY live window, so
// one Enter stops everything currently asking (the fail direction of a mistake
// here is "the write did not happen", never "it did").
//
// What this is NOT: it never consumes a line when no window is live (the return
// value is the only thing that reroutes a line away from the verb grammar), it
// never claims a channel - the veto still travels through Replies.Veto on the
// channel THIS host declared, so an undeclared host hears the same registry
// refusal as before (which is the AC#5 reverse reading: pull the declaration
// out and the same input executes) - and it never touches an L2 card: the queue
// is not a window, and an L2 answer keeps its verb grammar.
func (s *replySurface) vetoLiveL1Windows(out io.Writer) bool {
	if s == nil || s.gate == nil {
		return false
	}
	wins := s.gate.LiveL1Windows()
	if len(wins) == 0 {
		return false
	}
	for _, w := range wins {
		text, err := s.veto(w.CorrelationID)
		if err != nil {
			// The window may have expired between the listing above and this
			// veto, or the host never declared a channel; s.veto already booked
			// the refusal in the audit, so here it is only said out loud.
			fmt.Fprintf(out, "wisp run: 答复未被接受（L1 窗口否决）：%v\n", err)
			continue
		}
		fmt.Fprintf(out, "wisp run: %s\n", text)
	}
	return true
}

// handle maps one verb onto the surface. The verbs are the host's own grammar -
// not a wire protocol, and nothing here is a C17 method name.
func (s *replySurface) handle(verb, corr, arg string) (string, error) {
	switch verb {
	case "yes":
		return s.allow(corr)
	case "session":
		return s.session(corr)
	case "no":
		return s.reject(corr, washReason(arg))
	case "veto":
		return s.veto(corr)
	case "panel-no":
		return s.panelReject(corr, washReason(arg))
	case "panel-yes":
		return s.panelAllow(corr)
	case "always":
		return s.always(corr)
	case "revoke":
		// Ticket 226 AC#5: the rule text (or the bare directory) arrives spread
		// across the corr and arg slots - splitReplyLine cuts on single spaces,
		// and a rule text is full of them. Rejoin before parsing: a revocation
		// names a line, it does not answer a card.
		return s.revokeRule(strings.TrimSpace(corr + " " + arg))
	case "decisions":
		// Ticket 219 AC#5's third destination: the operator reads back what was
		// answered and why. No correlation id - this is a list, not a decision.
		return s.decisions()
	case "head":
		return s.head()
	case "view":
		return s.view(corr)
	case "help":
		return "yes/session/no [理由]/veto（原生侧；session＝本会话内允许，D45-2 第三枚答复）· " +
			"always <编号>（长期：先把要存的规则印出来，再走一张 L2 重新确认卡）· " +
			"revoke <规则或目录>（撤销一条长期规则：revoke [fs] allowed_dirs += \"D:/某目录\" 或 revoke D:/某目录，" +
			"只减不加、下一次启动生效，票 226）· " +
			"decisions（回看本进程最近答复过的卡片与理由，票 219）· " +
			"L1 窗口存活时任意输入（含回车）＝否决（票 170）· " +
			"panel-no/panel-yes/head/view（面板路线，无允许）· quit", nil
	default:
		return "", fmt.Errorf("未知答复指令 %q（yes/session/no/veto/always/revoke/decisions/panel-no/panel-yes/head/view/help/quit）", verb)
	}
}

// replyReasonMax bounds one typed reason. The reason travels into the tool
// result text the model reads, so the ceiling is not a display choice; 160 runes
// is deliberately below observe's own 512-character log bound so the ledger and
// the model-visible text cannot disagree about how much of a sentence arrived.
//
// TICKET 219 AC#5: THE CEILING IS NOT ALLOWED TO BE SILENT. Until this leg the
// wash dropped every rune past 160 with no mark anywhere: 200 typed, 160
// delivered, and neither the model, the audit line nor the operator could tell a
// 160-rune reason from a 200-rune one (ticket 219's 2026-11-15 收件 measured it:
// 200 in / 160 out, zero hits for any truncation marker). The fix keeps the
// ceiling and turns the cut into a READING - washReplyText reports how many runes
// it dropped, replyReason carries that number, and the text itself ends with a
// sentence naming both numbers. The marker is appended ON TOP of the 160 on
// purpose: a mark that could itself be cut off is the same silence one layer
// deeper.
const replyReasonMax = 160

// replyTruncationMark opens that sentence. It is ONE spelling so the audit line,
// the model-visible text and the review list cannot drift into three ways of
// saying the same thing, and it is a sentence rather than an ellipsis because
// "…" is a display convention a reader may miss while a count is not.
const replyTruncationMark = "…(理由过长，已截断"

// replyReason is one operator-typed reason as this surface carries it: the washed
// text that will be relayed, plus how much of what the operator typed did NOT
// fit. The second number is what makes 「空理由」 and 「超长理由」 two states a
// caller can branch on instead of two strings a reader has to guess at (AC#5's
// "空理由与超长理由各有确定行为").
type replyReason struct {
	text string
	// dropped counts the runes of the operator's own sentence the ceiling
	// removed, after whitespace collapsing. Zero means text is everything they
	// typed.
	dropped int
}

// truncationNote is the sentence a caller adds when the reason was cut, or ""
// when it was not.
func (r replyReason) truncationNote() string {
	if r.dropped <= 0 {
		return ""
	}
	return fmt.Sprintf("理由超过 %d 字上限，已截断 %d 字（截断标记在正文末尾，模型与审计都看得到）",
		replyReasonMax, r.dropped)
}

// washReason turns one raw argument into the value the surface answers with.
func washReason(raw string) replyReason {
	text, dropped := washReplyText(raw)
	return replyReason{text: text, dropped: dropped}
}

// sanitizeReplyText is the historical name for the wash, kept because this
// package's own comments and 216's display-scrub notes already point at it. It
// returns the marked text; use washReason where the dropped count matters.
func sanitizeReplyText(s string) string {
	text, _ := washReplyText(s)
	return text
}

// washReplyText washes one operator-typed reason and reports how many of its
// runes the ceiling removed.
//
// Why wash a human's own words at all: the reason is relayed into the model's
// tool result verbatim (bridge.go's orDefault(why, …)), and a control character
// in that string is how a reply starts spelling things to the model that the
// operator never typed. This is the same family as 票 216's control-character
// wash, applied at the surface that produces the text rather than downstream.
// Only the shape is washed - never the meaning: the words, their order and their
// length survive, which is what keeps 「为什么拒绝」 the operator's sentence
// rather than something this file paraphrased.
//
// The cut is now countable and stated (ticket 219 AC#5): a reason over the
// ceiling comes back as the first replyReasonMax runes of the washed text plus
// replyTruncationMark's sentence, and the second return value is the number of
// runes dropped. The dropped runes themselves are NOT kept anywhere, and that
// limit is said out loud in the marker rather than implied - storing a second,
// unbounded copy of operator text was not this ticket's to add.
func washReplyText(s string) (string, int) {
	if s == "" {
		return "", 0
	}
	washed := make([]rune, 0, len(s))
	pendingSpace := false
	for _, r := range s {
		if unicode.IsControl(r) {
			r = ' '
		}
		if r == ' ' {
			// A run of whitespace collapses to one separator, so a reason typed
			// with a trailing newline cannot become a two-line string in the
			// model's tool result.
			pendingSpace = len(washed) > 0
			continue
		}
		if pendingSpace {
			washed = append(washed, ' ')
			pendingSpace = false
		}
		washed = append(washed, r)
	}
	if len(washed) <= replyReasonMax {
		return string(washed), 0
	}
	dropped := len(washed) - replyReasonMax
	text := string(washed[:replyReasonMax]) +
		fmt.Sprintf("%s %d 字，原文共 %d 字)", replyTruncationMark, dropped, len(washed))
	return text, dropped
}

// replyReviewKeep is how many answered replies the in-process review list holds.
// A display window, not an SLO number: the durable record of every reply is the
// audit line, and this exists so the operator can ask "what did it decide for me
// just now" without opening a log file.
const replyReviewKeep = 64

// replyDecision is one line of that list.
type replyDecision struct {
	when   time.Time
	kind   string
	corr   string
	tool   string
	route  string
	detail string
}

// replyReview is the bounded review list behind the `decisions` console verb -
// ticket 219 AC#5's third destination (「今天它替我做了哪些决定」). It is a mutex
// and not a channel because two writers reach it (the reply listener's goroutine
// and the ui-sta tray thread) and the reader is a third.
type replyReview struct {
	mu      sync.Mutex
	entries []replyDecision
}

func newReplyReview() *replyReview { return &replyReview{} }

// add appends one answer and keeps only the newest replyReviewKeep of them. nil
// is legal: a surface a case built by hand has no review list and books nothing.
func (r *replyReview) add(d replyDecision) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, d)
	if over := len(r.entries) - replyReviewKeep; over > 0 {
		r.entries = append(r.entries[:0], r.entries[over:]...)
	}
}

// snapshot copies the list oldest-first, so the caller never holds the lock.
func (r *replyReview) snapshot() []replyDecision {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]replyDecision(nil), r.entries...)
}

