//go:build windows

package main

// The resident-side halves of the conversation-mode first-enable privacy gate
// (see conversation_privacy_180p.go for the shared gate and the rationale).
//
// FACE 1 - the machine-side gate. startResidentVoice calls
// applyConversationPrivacyGate180P after buildResidentVoice and BEFORE
// rv.start, so the mode the D43 speak-end row selector walks
// (speakDoneFacts180K, resident_voice_engines_windows.go:1279) is always the
// POST-gate value: an unconfirmed true never reaches the dispatch. The card is
// raised at assembly time, where the ball is already live (startResidentBall
// ran long before this line in runResident) and the voice chain has not yet
// opened a feed - an unanswered card costs the boot its wait, never an open mic.
//
// FACE 2 - the scope seam. The gate's confirm surface and the data dir ride
// the residency as fields (set here, before rv.start), and
// voiceResidency.scopeOptions() assigns Hooks.PrivacyConfirm from them, so
// session.Scope.EnableConversation's own contract (first enable through the
// confirm, later enables free, nil fail-closed) is honoured by the same card
// and the same marker instead of a lookalike.

import (
	"fmt"
	"log/slog"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// conversationPrivacyCardTaskID / conversationPrivacyCardTool are the card
// identity, in the exact shape config_reload.go's host-raised confirmation
// uses: a synthetic task id the gate has never seen because no model call
// raised it, admitted for the length of the card and revoked on the way out
// (AskOnTaskRoot's own D47 discipline), plus a dotted name that names the door
// being opened - not a tool the model can call, not a C17 method.
const (
	conversationPrivacyCardTaskID = "host:conversation-privacy"
	conversationPrivacyCardTool   = "voice.conversation_mode"
)

// conversationPrivacyConfirm180P is the L2 card surface the gate asks through:
// this process's own approval gate, on the task root D38(e) step 3 owns, so an
// exit while the card is up refuses it through the ordinary queue funnel
// instead of killing the process with a card pending. The reason says what the
// user is actually agreeing to - the mic stays open between turns - because a
// privacy card that cannot say what it protects is not a confirmation, it is a
// receipt for a click.
func (ra *residentApproval) conversationPrivacyConfirm180P() bool {
	ans, why := ra.AskOnTaskRoot(residentCard{
		TaskID: conversationPrivacyCardTaskID,
		Tool:   conversationPrivacyCardTool,
		Level:  risk.L2,
		Reason: "开启对话模式（voice.conversation_mode=true）：播报结束后麦克风不回 Warm，" +
			"而是停在 Listening 继续收音，等待你紧接着说下一句（D43#28）。「首次开启需要一次 L2 级隐私确认」" +
			"是产品契约（PLAN.md:2211）：确认后本机会记一条「已确认」标记，之后不再弹；拒绝或不做声，" +
			"本次启动对话模式保持关闭，下一次启动会再问。（这张卡由语音链装配签发，与任何一次工具调用无关）",
	})
	if ans != tools.AnswerAllow {
		slog.Warn("voice: conversation_mode 首次开启的 L2 隐私确认未通过",
			"answer", string(ans), "why", why)
		return false
	}
	return true
}

// applyConversationPrivacyGate180P is the machine-side gate application. A mode
// that is already off costs nothing (nothing to confirm); the four gate
// outcomes each print their own line, and only the two ON outcomes leave
// plan.ConversationMode true.
func (rv *residentVoice) applyConversationPrivacyGate180P(dataDir string, confirm func() bool) {
	if rv == nil || rv.residency == nil || !rv.plan.ConversationMode {
		return
	}
	// Both faces share these two values from here on: the scope's
	// Hooks.PrivacyConfirm (scopeOptions) asks the same card and reads the same
	// marker, so the scope seam and the machine-side selector cannot disagree.
	rv.residency.conversationPrivacyDataDir = dataDir
	rv.residency.conversationPrivacyConfirm = confirm

	granted, outcome := conversationPrivacyGate180P(dataDir, confirm)
	if !granted {
		// The selector carries the post-gate value: every refused outcome lands
		// the mode off before the switch below writes its own sentence.
		rv.plan.ConversationMode = false
	}
	switch outcome {
	case convPrivacyAlready:
		slog.Info("voice: conversation_mode 首次开启的隐私确认已在既往启动完成（标记件在）",
			"marker", conversationPrivacyMarkerPath(dataDir), "outcome", string(outcome))
	case convPrivacyGranted:
		slog.Info("voice: conversation_mode 首次开启的 L2 隐私确认已通过（PLAN.md:2211），已写确认标记",
			"marker", conversationPrivacyMarkerPath(dataDir), "outcome", string(outcome))
		fmt.Printf("wisp: voice 对话模式隐私确认已通过（L2，首次开启）：确认标记已写入 %s；删除该文件即撤回「已确认」\n",
			conversationPrivacyMarkerPath(dataDir))
	case convPrivacyDenied:
		slog.Warn("voice: conversation_mode 保持关闭（首次开启的 L2 隐私确认被拒绝/无人应答/超时）",
			"outcome", string(outcome),
			"detail", "D43#28 不走：本进程的播报结束回到 Warm（#26），确认标记未写入，下一次启动会再问")
		fmt.Printf("wisp: voice 对话模式保持关闭：首次开启的 L2 隐私确认没有拿到「允许」（拒绝/没答/超时都会落到这里，fail-closed）。" +
			"本次播报结束按会话模式关闭处理（D43#26），下一次启动会再弹一次。\n")
	case convPrivacyNoSurface:
		slog.Warn("voice: conversation_mode 保持关闭（本进程没有可弹隐私确认的面，fail-closed）",
			"outcome", string(outcome),
			"detail", "与 session.EnableConversation 对 nil PrivacyConfirm 的处理同向：没有确认面就不能静默开一个保持收音的模式")
		fmt.Printf("wisp: voice 对话模式保持关闭：本进程没有可弹 L2 隐私确认的面（fail-closed，不是静默生效）。下一次启动若面恢复会再问。\n")
	}
}
