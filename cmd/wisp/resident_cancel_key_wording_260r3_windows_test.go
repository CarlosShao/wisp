//go:build windows

package main

// Ticket 260-r3 (ledger A595 §1): the two wording judges for the cancel key.
//
// WHAT THIS FILE PINS, in the user's words. Ticket 260 形ⓐ made the confirmation
// window borrow the key that [hotkey] cancel CONFIGURES (default 档 still the bare
// Esc). Ten human-facing sentences still said "Esc" by fiat. This file is the pair
// of rulers the ticket demands, kept as two separate cells on purpose:
//
//	① default 档: every sentence that names a key renders BYTE-FOR-BYTE the
//	  sentence this file printed before the change - a user who never touched
//	  [hotkey] cancel must not see a single character move;
//	② seeded 档: with the cancel binding changed to another combination, the same
//	  sentence names THAT key and no longer contains "Esc" anywhere.
//
// WHY THE SEEDED CASE USES THE cancelKeyRead SEAM, and what it does NOT prove.
// The production value comes from the ball's own hotkey report
// (residentCancelKeySpelling below reads that report's cancel line, which is the
// line TakeEscForCancel fills from the SAME cancelBorrow it hands to
// RegisterHotKey). A live ball window is a winlive rig - this leg is forbidden to
// run one (no owner word for "register a hot key on the real desktop today"), and
// CI never runs that tag either. So what executes here is the plumbing from the
// label to the sentence; the plumbing from the ball's report to the label is
// pinned by reading plus the winlive claim constant, and is booked as 〔未实测〕 in
// .scratch/wisp/probes/260/r3/wording.md §5/§6 rather than claimed as measured.
// TestTicket260R3ProductionPathHasNoSeam is the guard that the seam is not a way
// to fake the shipped process: the constructor must leave it nil.
//
// The three 档一 sentences (no channel loaded at all) are a THIRD shape, per
// A595 §1 boundary ④: their honest form names no key at all, so they are asserted
// as "says the slot is empty, never mentions a key" in both 档, and they are NOT
// held to the zero-drift rule - the old wording itself was the lie.

import (
	"errors"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/ball"
)

// The four sentences as cmd/wisp printed them BEFORE this leg, copied verbatim out
// of git show 57a33804:cmd/wisp/resident_approval_windows.go (lines 287, 309, 311,
// 495 of that anchor). They are the zero-drift half of cell ①: if this leg's
// rewording moves even one character in the default 档, these assertions fail.
const (
	old260r3LoadedAudit = "resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；" +
		"单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
	old260r3RejectedSaid = "按 Esc 否决卡片 corr-260r3 被拒：桌面拒绝原因"
	old260r3DoneSaid     = "按 Esc 否决了卡片 corr-260r3（L1 / fs.write）：该调用未执行，答案已入审计"
	old260r3StatusLine   = "审批门已装配进本进程（取消通道：Esc 已加载；等待中的确认项：0）"

	// The two 档一 sentences as they read BEFORE this leg, kept here so the test
	// can say out loud which words it deliberately dropped.
	old260r3NoBallAudit  = "四条否决通道全部保持未加载（Esc 无处可借，卡片无处可呈）"
	old260r3NoExecutor   = "Esc 通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
	seeded260r3CancelKey = "Ctrl+Alt+Q"
)

// card260r3 is the ledger card the two veto sentences are built from. Only the
// three fields the sentences print are filled, which is what vetoByEsc reads back
// from its own ledger.
func card260r3() approval.ReplyCard {
	return approval.ReplyCard{
		CorrelationID: "corr-260r3",
		Level:         "L1",
		Tool:          "fs.write",
	}
}

// TestTicket260R3DefaultWordingIsTheOldSentence is cell ①: nothing a default-档
// user can see may move. The three sentences that do NOT go through a ball window
// run on the real production reader (no seam installed), so what is compared here
// is the shipped fallback path, not a stand-in.
func TestTicket260R3DefaultWordingIsTheOldSentence(t *testing.T) {
	if ball.DefaultHotkeys().Cancel != "Esc" {
		t.Fatalf("the premise of this whole cell moved: ball's product default cancel binding is %q, not Esc; "+
			"update the old-sentence constants deliberately, not by accident", ball.DefaultHotkeys().Cancel)
	}
	ra := newResidentApproval()

	if got := vetoRejectedLine(ra.cancelKeySpelling(), "corr-260r3", errors.New("桌面拒绝原因")); got != old260r3RejectedSaid {
		t.Errorf("默认档否决被拒句漂移：got %q, want %q", got, old260r3RejectedSaid)
	}
	if got := vetoDoneLine(ra.cancelKeySpelling(), card260r3()); got != old260r3DoneSaid {
		t.Errorf("默认档否决成功句漂移：got %q, want %q", got, old260r3DoneSaid)
	}

	// The boot-time pair, reached through bindBallHost with a window that is only
	// a placeholder: the leg's own reader is replaced (a zero-value ball cannot
	// answer an uiRun), and the value handed back is the product default the
	// ball's chain resolves to, so the rendered sentence is still the default 档.
	ra2 := newResidentApproval()
	ra2.cancelKeyRead = func() string { return ball.DefaultHotkeys().Cancel }
	said := captureStdout128(t, func() {
		if !ra2.bindBallHost(&residentBall{b: &ball.Ball{}, cancelHosted: true}) {
			t.Error("bindBallHost 在没有执行者的情况下声称装配好了")
		}
	})
	want := "wisp: [audit] " + old260r3LoadedAudit + "\n"
	if !strings.Contains(said, want) {
		t.Errorf("默认档装配回执漂移：got %q, want it to contain %q", said, want)
	}
	if line := ra2.residentStatusLine(); line != old260r3StatusLine {
		t.Errorf("默认档状态句漂移：got %q, want %q", line, old260r3StatusLine)
	}
}

// TestTicket260R3SeededKeyReplacesEsc is cell ②: with the cancel binding moved to
// another combination, every key-naming sentence names THAT key and never says
// Esc. This is the cell with teeth against a re-hardcoded literal - if any of the
// four sentences goes back to printing a fixed "Esc", the corresponding assertion
// below goes red.
func TestTicket260R3SeededKeyReplacesEsc(t *testing.T) {
	ra := newResidentApproval()
	ra.cancelKeyRead = func() string { return seeded260r3CancelKey }

	cases := []struct {
		name string
		got  string
	}{
		{"否决被拒句", vetoRejectedLine(ra.cancelKeySpelling(), "corr-260r3", errors.New("桌面拒绝原因"))},
		{"否决成功句", vetoDoneLine(ra.cancelKeySpelling(), card260r3())},
	}
	for _, tc := range cases {
		if !strings.Contains(tc.got, seeded260r3CancelKey) {
			t.Errorf("%s 没有念出配置里那枚键：%q（应含 %q）", tc.name, tc.got, seeded260r3CancelKey)
		}
		if strings.Contains(tc.got, "Esc") {
			t.Errorf("%s 仍在念 Esc：%q", tc.name, tc.got)
		}
	}

	// The boot-time pair, same placeholder window as above, only the key changed.
	ra2 := newResidentApproval()
	ra2.cancelKeyRead = func() string { return seeded260r3CancelKey }
	said := captureStdout128(t, func() {
		ra2.bindBallHost(&residentBall{b: &ball.Ball{}, cancelHosted: true})
	})
	for _, line := range []string{said, ra2.residentStatusLine()} {
		if !strings.Contains(line, seeded260r3CancelKey) {
			t.Errorf("改过键以后的装配句没有念那枚键：%q", line)
		}
		if strings.Contains(line, "Esc") {
			t.Errorf("改过键以后的装配句仍在念 Esc：%q", line)
		}
	}
}

// TestTicket260R3UnloadBranchesNameNoKey is A595 §1 boundary ④'s 档一: the two
// sentences printed when the cancel channel is NOT loaded must say the slot is
// empty and must not name any key - not Esc, and not the configured one either,
// because there is nothing borrowed to name. Both 档 are checked, and the seeded
// key is checked for ABSENCE: an "improvement" that prints the configured key in
// these two branches would be advertising a channel this process refuses.
func TestTicket260R3UnloadBranchesNameNoKey(t *testing.T) {
	for _, key := range []string{ball.DefaultHotkeys().Cancel, seeded260r3CancelKey} {
		ra := newResidentApproval()
		ra.cancelKeyRead = func() string { return key }

		noBall := captureStdout128(t, func() {
			ra.bindBallHost(&residentBall{})
		})
		if !strings.Contains(noBall, "取消键无处可借") {
			t.Errorf("没有球窗那一支没有说「取消键无处可借」（旧句逐字 %q）：got %q", old260r3NoBallAudit, noBall)
		}
		if strings.Contains(noBall, "Esc") {
			t.Errorf("没有球窗那一支还在指枚一枚键：%q", noBall)
		}

		// The other 档一 branch needs rb.b non-nil and no injected executor; the
		// ball pointer is a placeholder the branch never touches.
		ra2 := newResidentApproval()
		ra2.cancelKeyRead = func() string { return key }
		noExecutor := captureStdout128(t, func() {
			ra2.bindBallHost(&residentBall{b: &ball.Ball{}})
		})
		if !strings.Contains(noExecutor, "取消键通道保持未加载") {
			t.Errorf("没有执行者那一支没有说「取消键通道保持未加载」（旧句逐字 %q）：got %q", old260r3NoExecutor, noExecutor)
		}
		if strings.Contains(noExecutor, "Esc") {
			t.Errorf("没有执行者那一支还在指枚一枚键：%q", noExecutor)
		}
		// 档一 must not name the key either: the honest form is "无处可借", not
		// "your configured key", because nothing was borrowed.
		if strings.Contains(noBall, key) || strings.Contains(noExecutor, key) {
			t.Errorf("未加载的两支指枚了键名 %q：%q / %q", key, noBall, noExecutor)
		}
	}
}

// TestTicket260R3CancelKeyComesFromTheBallChain is the single-source guard for the
// half this file cannot reach with a window: with no ball at all the reader must
// answer with the BALL's own default constant (not a literal this package copied),
// and an empty cancel binding must resolve to that same value - which is exactly
// the input internal/ball's resolveCancelBorrow maps to the default.
func TestTicket260R3CancelKeyComesFromTheBallChain(t *testing.T) {
	if got, want := residentCancelKeySpelling(nil), ball.DefaultHotkeys().Cancel; got != want {
		t.Errorf("没有球窗时的取消键名 = %q, want ball 自己的默认那枚 %q", got, want)
	}
	if got := residentCancelKeySpelling(nil); got == "" {
		t.Error("取消键名是空串：那会印出「按  否决卡片」这种第二枚必谎的尺")
	}
	// The same value the two veto sentences print with no seam installed - so the
	// fallback above is the shipped path, not a test-only constant.
	ra := newResidentApproval()
	if got := ra.cancelKeySpelling(); got != ball.DefaultHotkeys().Cancel {
		t.Errorf("产码路径（无读口、无球窗）念 %q, want ball 的默认 %q", got, ball.DefaultHotkeys().Cancel)
	}
}

// TestTicket260R3ProductionPathHasNoSeam keeps the injected reader an injection
// point and nothing else: the shipped constructor must leave it nil, so no process
// ever prints a key it did not read off the ball's report. Without this assertion
// the seam above could quietly become the production path.
func TestTicket260R3ProductionPathHasNoSeam(t *testing.T) {
	if ra := newResidentApproval(); ra.cancelKeyRead != nil {
		t.Fatal("newResidentApproval 装了取消键读口：产码路径不再读球自己的回执行了")
	}
	if ra := newResidentApprovalWithConfig(""); ra.cancelKeyRead != nil {
		t.Fatal("newResidentApprovalWithConfig 装了取消键读口：同上")
	}
	// And the shipped path, with no seam and no window: a veto on a channel this
	// process never loads is refused, and the refusal sentence is built by
	// production code end to end. It must still念 the ball's default key in the
	// key slot, and its reason half must carry 档一's new wording (the reason is
	// approval.ChannelError.Msg, i.e. this package's unavailableText - so the
	// blast radius of that one string reaches here, through the %v).
	ra := newResidentApproval()
	ra.cards.Record(approval.Prompt{CorrelationID: "corr-260r3-real", Level: "L1", Tool: "fs.write"})
	said := ra.vetoByEsc()
	if !strings.HasPrefix(said, "按 "+ball.DefaultHotkeys().Cancel+" 否决卡片 ") {
		t.Errorf("产码路径那句没有念 ball 自己那枚默认键：%q", said)
	}
	if strings.Contains(said, seeded260r3CancelKey) {
		t.Fatalf("没装读口却念出了种子键名：%q", said)
	}
	if strings.Contains(said, "Esc 取消不可用") {
		t.Errorf("被拒原因的档一句子还是旧形：%q", said)
	}
	if !strings.Contains(said, "快捷键取消不可用") {
		t.Errorf("被拒原因没有说「快捷键取消不可用」（档一那枚在本包的另一半）：%q", said)
	}
}
