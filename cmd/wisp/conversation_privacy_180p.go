package main

// The deferred-disposition 20261110 contract item C2 wiring: voice
// .conversation_mode=true 的首次 L2 隐私确认（PLAN.md:2211 逐字 mandate：
// 「首次开启弹一次隐私说明并按 L2 确认，记录日志」）。
//
// THE SEAM THE CONTRACT NAMES. session.Scope.EnableConversation
// (internal/session/scope.go:520-569) already owns the whole shape: the first
// enable MUST pass Hooks.PrivacyConfirm (an L2-grade card), every later enable
// passes without a prompt (ticket 28 AC4), and a nil hook fails closed. What
// was missing was production: the resident machine walks D43#26/#28 straight
// off the plan (resident_voice_engines_windows.go:1263), bypassing the scope,
// and scopeOptions() assigned no PrivacyConfirm - the mandate had an engine and
// no方向盘.
//
// THE WIRING. One gate function (conversationPrivacyGate180P below), one
// persisted marker, TWO faces that share it:
//
//  1. the machine-side face: startResidentVoice applies the gate to the plan
//     before the voice chain starts, so an unconfirmed mode never reaches the
//     D43#28 row selector at all (the mode the announce end dispatches is the
//     POST-gate value);
//  2. the scope face: voiceResidency.scopeOptions() now assigns
//     Hooks.PrivacyConfirm to the same gate, so session's own contract - first
//     enable through the confirm, nil fail-closed - is honoured with the same
//     marker and the same card instead of a lookalike.
//
// "首次" IS PER DATA DIR, NOT PER PROCESS. EnableConversation's `confirmed`
// flag dies with the scope, so without persistence every boot would ask again -
// which is not "once". The marker is one small file in the data dir, written
// only AFTER an explicit grant; its absence is what makes the next boot ask.
// Deleting the marker re-arms the ask, which is the honest reverse of "once".
//
// THE FAIL DIRECTIONS, each named:
//   - no card surface (confirm == nil): the gate refuses, exactly the direction
//     session.EnableConversation takes for a nil hook (fail-closed, loud);
//   - card denied / abandoned / timed out: the mode stays OFF for this process
//     and the marker is NOT written - the next boot asks again;
//   - the marker write fails after a grant: the grant stands for THIS process
//     (the user did say yes) and the warning names that the next boot will ask
//     again - a lost consent record is a re-ask, never a silent permanent yes.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/observe"
)

// conversationPrivacyMarkerName is the once-marker file. A name, not a hidden
// dotfile: the owner can find and delete it, and deleting it is exactly how a
// "confirmed once" becomes "ask me again".
const conversationPrivacyMarkerName = "conversation-mode.confirmed.txt"

// conversationPrivacyOutcome is the four ways out of the gate, each one the
// boot report and the audit line can name without collapsing two different
// facts into one word (the discipline 票 216 set).
type conversationPrivacyOutcome string

const (
	// convPrivacyAlready: a previous boot's grant is on record; no card, mode on.
	convPrivacyAlready conversationPrivacyOutcome = "already-confirmed"
	// convPrivacyGranted: the L2 card was answered allow; marker written (or a
	// failed write warned about), mode on.
	convPrivacyGranted conversationPrivacyOutcome = "granted"
	// convPrivacyDenied: the card was refused / nobody answered / it timed out;
	// mode OFF this process, marker untouched, next boot asks again.
	convPrivacyDenied conversationPrivacyOutcome = "denied"
	// convPrivacyNoSurface: no confirm function was handed in at all - the same
	// fail-closed direction session.EnableConversation takes for a nil hook.
	convPrivacyNoSurface conversationPrivacyOutcome = "no-surface-fail-closed"
)

// conversationPrivacyMarkerPath is the marker file inside the data dir.
func conversationPrivacyMarkerPath(dataDir string) string {
	return filepath.Join(dataDir, conversationPrivacyMarkerName)
}

// conversationPrivacyConfirmed reports whether a previous grant is on record.
// A directory at that name is NOT a grant ( somebody's odd fixture is not
// consent); only a regular file counts.
func conversationPrivacyConfirmed(dataDir string) bool {
	st, err := os.Stat(conversationPrivacyMarkerPath(dataDir))
	return err == nil && st.Mode().IsRegular()
}

// writeConversationPrivacyMarker records the grant. The body names the key,
// the level and the wall-clock moment - a receipt, not a flag byte, so a human
// who finds the file can tell what it is and when it was given.
func writeConversationPrivacyMarker(dataDir string) error {
	body := fmt.Sprintf("key=voice.conversation_mode\nlevel=L2\ngranted_at=%s\nnote=删除此文件即撤回「已确认」，下次开启会重新弹隐私确认\n",
		observe.WallTimestampUTC(observe.NowWallUTC()))
	return os.WriteFile(conversationPrivacyMarkerPath(dataDir), []byte(body), 0o644)
}

// conversationPrivacyGate180P is the whole gate, shared by both faces. confirm
// is the card surface (the resident approval's AskOnTaskRoot adapter in
// production; a test stub elsewhere). nil is fail-closed, by the same sentence
// session.Scope uses for a nil PrivacyConfirm hook.
func conversationPrivacyGate180P(dataDir string, confirm func() bool) (bool, conversationPrivacyOutcome) {
	if conversationPrivacyConfirmed(dataDir) {
		return true, convPrivacyAlready
	}
	if confirm == nil {
		return false, convPrivacyNoSurface
	}
	if !confirm() {
		return false, convPrivacyDenied
	}
	if err := writeConversationPrivacyMarker(dataDir); err != nil {
		// The grant stands; the record does not. Say which, so the next boot's
		// re-ask is a named consequence and not a mystery.
		fmt.Printf("wisp: voice conversation_mode 隐私确认标记写入失败（%v）：本次进程按已确认运行，下一次启动会再弹一次确认\n", err)
	}
	return true, convPrivacyGranted
}
