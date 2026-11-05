package approval_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Ticket 21's L1 route is a 2-3s pre-execution BLOCK with a countdown on the
// ball strip. The number the strip draws is the Remaining field of the events
// this layer emits - and every L1 Remaining this layer emitted was the
// CONFIGURED window length, never the live one: a warning raised 2.5s into a 3s
// block still said "3s left", so a ring driven by it would freeze at full and
// then jump to zero.
//
// The warning below is the one L1 event a user can provoke at a known instant
// (the KWS veto word on a host where ticket 41 has not loaded that channel), so
// it is where the honesty of the number is pinned.
func TestTicket21UnavailableChannelWarningCarriesLiveRemaining(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)

	res := runWindow(t, g, context.Background(), l1Decision("C:/data/a.txt"))
	ui.wait(t) // the block is on screen and counting
	clk.Advance(2500 * time.Millisecond)

	// The user tries the voice veto word. KWS is not loaded here (ticket 41), so
	// the honest answer is the 「语音取消不可用」 warning - and the countdown keeps
	// running, which is why the strip needs the number at THIS instant.
	err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelKWS})
	if !errors.Is(err, approval.ErrChannelUnavailable) {
		t.Fatalf("未加载通道的否决 = %v, want ErrChannelUnavailable", err)
	}

	warn := ui.ofKind(approval.EventWarning)
	if len(warn) != 1 {
		t.Fatalf("醒目提示事件 = %d 条, want exactly 1: %+v", len(warn), warn)
	}
	if !strings.Contains(warn[0].Text, "语音取消不可用") {
		t.Errorf("提示文案必须原样包含「语音取消不可用」，得 %q", warn[0].Text)
	}
	// 2.5s of the 3s block have elapsed, so the live remaining is 0.5s.
	if got, want := warn[0].Remaining, 500*time.Millisecond; got != want {
		t.Errorf("提示事件的 Remaining = %v, want %v (窗口已过 2.5s): "+
			"报整窗 3s 的倒计时环会骗人", got, want)
	}

	// The window is still counting (an unloaded channel changes nothing): close it
	// through a loaded one, and the answer is still the veto, not an execute.
	if err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelEsc}); err != nil {
		t.Fatalf("loaded-channel veto: %v", err)
	}
	if a := mustAnswer(t, res); a.a != tools.AnswerVeto {
		t.Fatalf("answer=%q, want veto: %s", a.a, a.why)
	}
	clk.Advance(10 * time.Second)
}

// The other end of the same ruler: at the instant the block opens, "left" is the
// whole block. Without this half a reading that mixed up elapsed with remaining
// (or clamped everything to zero) would pass the case above.
func TestTicket21WarningAtWindowStartReportsTheWholeBlock(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)

	res := runWindow(t, g, context.Background(), l1Decision("C:/data/b.txt"))
	ui.wait(t)

	if err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelKWS}); !errors.Is(err, approval.ErrChannelUnavailable) {
		t.Fatalf("未加载通道的否决 = %v, want ErrChannelUnavailable", err)
	}
	warn := ui.ofKind(approval.EventWarning)
	if len(warn) != 1 {
		t.Fatalf("醒目提示事件 = %d 条, want exactly 1: %+v", len(warn), warn)
	}
	if got, want := warn[0].Remaining, 3*time.Second; got != want {
		t.Errorf("窗口刚开时提示的 Remaining = %v, want %v（整窗，不是 0 也不是已过时间）", got, want)
	}

	if err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelBall}); err != nil {
		t.Fatalf("loaded-channel veto: %v", err)
	}
	if a := mustAnswer(t, res); a.a != tools.AnswerVeto {
		t.Fatalf("answer=%q, want veto: %s", a.a, a.why)
	}
	clk.Advance(10 * time.Second)
}
