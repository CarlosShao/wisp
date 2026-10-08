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
// 否决词 / 面板拒绝, none of which exists in a terminal. This file therefore
// calls Gate.Veto with the channel the HOST declares (hostVetoChannel, empty on
// the console leg) and lets the gate's own ChannelRegistry answer - it invents
// no fifth channel and claims no loaded one, so nothing here can render an
// unavailable channel as available (B1).

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
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

// maxTrackedCards is the ledger's own ceiling, read off the seam that enforces
// it (approval.MaxTrackedCards). The queue bounds pending items at its own
// DefaultMaxPending and an L1 window leaves the ledger on its dismissal or
// handoff event, so the ceiling exists only for the case where an event never
// arrives (a UI that returned an error mid-flight). Eviction is the safe
// direction: an evicted card has no grant held for it, so it can be refused but
// never allowed by this surface.
const maxTrackedCards = approval.MaxTrackedCards

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
	// vetoChannel is the channel this host's cancel transport really is. Empty
	// means this assembly wired none of SPEC-06 §2's four, and Gate.Veto's own
	// ChannelRegistry says so back.
	vetoChannel approval.Channel
	// ctx is the listener's own context; it carries no authority.
	ctx context.Context
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
func (s *replySurface) reject(corr, reason string) (string, error) {
	return s.refuse(corr, reason, false)
}

// refuse is the shared body of the two reject routes, kept as one function
// because the ONLY difference between them is which router the decision travels
// through - and a second copy is how one of them grows a leniency the other does
// not have.
func (s *replySurface) refuse(corr, reason string, viaPanel bool) (string, error) {
	card, _ := s.live.look(corr)
	route, routeText := "native", "原生侧"
	var err error
	if viaPanel {
		route, routeText = "panel", "面板路线"
		err = s.live.h.PanelReject(s.ctx, corr, reason)
	} else {
		err = s.live.h.Reject(s.ctx, corr, reason)
	}
	if err != nil {
		s.record("REFUSED", corr, card.Tool, route+"/reject", err.Error())
		return "", err
	}
	s.record("ANSWERED", corr, card.Tool, route+"/reject", reason)
	if reason == "" {
		return "已拒绝 " + corr + "（" + routeText + "，未填理由）：未执行", nil
	}
	return "已拒绝 " + corr + "（" + routeText + "）：未执行，理由已随结果回给模型", nil
}

// panelReject is the panel route's refusal - SPEC-06 §9 layer 3's 「拒绝」
// button, and the only decision the page may send.
func (s *replySurface) panelReject(corr, reason string) (string, error) {
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
func (s *replySurface) record(kind, corr, tool, route, detail string) {
	if s == nil || s.audit == nil {
		return
	}
	s.audit("approval: REPLY %s corr=%q tool=%s route=%q detail=%q", kind, corr, tool, route, detail)
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
			"面板路线：panel-no <编号> [理由]、panel-yes <编号>（这一条只会被服务端 API 拒绝）、"+
			"head / view <编号> 看卡片；quit 退出监听。\n")
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
func runReplyLoop(ctx context.Context, s *replySurface, in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 0, 4096), 1<<16)
	for sc.Scan() {
		if ctx.Err() != nil {
			return
		}
		line := strings.TrimSpace(sc.Text())
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

// handle maps one verb onto the surface. The verbs are the host's own grammar -
// not a wire protocol, and nothing here is a C17 method name.
func (s *replySurface) handle(verb, corr, arg string) (string, error) {
	switch verb {
	case "yes":
		return s.allow(corr)
	case "session":
		return s.session(corr)
	case "no":
		return s.reject(corr, sanitizeReplyText(arg))
	case "veto":
		return s.veto(corr)
	case "panel-no":
		return s.panelReject(corr, sanitizeReplyText(arg))
	case "panel-yes":
		return s.panelAllow(corr)
	case "always":
		return s.always(corr)
	case "head":
		return s.head()
	case "view":
		return s.view(corr)
	case "help":
		return "yes/session/no/veto（原生侧；session＝本会话内允许，D45-2 第三枚答复）· " +
			"always <编号>（长期：先把要存的规则印出来，再走一张 L2 重新确认卡）· " +
			"panel-no/panel-yes/head/view（面板路线，无允许）· quit", nil
	default:
		return "", fmt.Errorf("未知答复指令 %q（yes/session/no/veto/always/panel-no/panel-yes/head/view/help/quit）", verb)
	}
}

// replyReasonMax bounds one typed reason. The reason travels into the tool
// result text the model reads, so the ceiling is not a display choice; 160 runes
// is deliberately below observe's own 512-character log bound so the ledger and
// the model-visible text cannot disagree about how much of a sentence arrived.
const replyReasonMax = 160

// sanitizeReplyText washes one operator-typed reason.
//
// Why wash a human's own words at all: the reason is relayed into the model's
// tool result verbatim (bridge.go's orDefault(why, …)), and a control character
// in that string is how a reply starts spelling things to the model that the
// operator never typed. This is the same family as 票 216's control-character
// wash, applied at the surface that produces the text rather than downstream.
// Only the shape is washed: the words, their order and their length survive, so
// 「为什么拒绝」 stays the operator's sentence rather than something this file
// paraphrased.
func sanitizeReplyText(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	runes := 0
	pendingSpace := false
	for _, r := range s {
		if unicode.IsControl(r) {
			r = ' '
		}
		if r == ' ' {
			// A run of whitespace collapses to one separator, so a reason typed
			// with a trailing newline cannot become a two-line string in the
			// model's tool result.
			pendingSpace = b.Len() > 0
			continue
		}
		if pendingSpace {
			b.WriteRune(' ')
			runes++
			pendingSpace = false
		}
		b.WriteRune(r)
		runes++
		if runes >= replyReasonMax {
			break
		}
	}
	return b.String()
}
