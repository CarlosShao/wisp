package approval_test

// Ticket 260-r3 (ledger A595 §1, 档一): the veto-channel availability lines.
//
// Two of A595's ten strings live in this package, and they are the two whose
// honest form is NOT another key name: both are reached only when the cancel
// channel is not loaded, so there is no borrowed key to print. The unavailable
// line used to read 「Esc 取消不可用」; it now reads 「快捷键取消不可用」, naming
// the slot the way the other three lines name theirs (语音取消不可用 / 面板取消不
// 可用 / 悬浮球取消不可用), and it holds for every value [hotkey] cancel may take.
//
// The loaded-channel label (channelNames[ChannelEsc], 「按 Esc 键」) is the ONE
// string of the ten this leg did NOT change, and it is not changed here either:
// making it follow the config needs a value this package cannot see without
// either importing internal/ball (a new package-level dependency edge, which
// AGENTS.md §1.2 books as a human-approval face) or a late write to the package
// map that gate.go:314 / gate.go:477 and report.go:142 read unlocked. The
// orchestrator's ruling is owed; wording.md §8 carries the three candidate forms
// and their costs. The assertion below pins today's text so the follow-up leg
// sees it move - it is a NAMED RESIDUAL RULER, not an approval of the wording.

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent/approval"
)

// TestTicket260R3EscUnavailableLineNamesNoKey is 档一 for this package: the
// unavailable half of the cancel line says the slot cannot fire, and does not
// point at a key this host is not holding.
func TestTicket260R3EscUnavailableLineNamesNoKey(t *testing.T) {
	reg := approval.NewChannels(approval.ChannelBall)
	var esc, ballLine approval.ChannelStatus
	for _, st := range reg.Statuses() {
		switch st.Channel {
		case approval.ChannelEsc:
			esc = st
		case approval.ChannelBall:
			ballLine = st
		}
	}
	if esc.Loaded {
		t.Fatal("NewChannels(ball) 却把取消通道报成已加载")
	}
	if ballLine.Text != "单击悬浮球" {
		t.Errorf("对照那行（本腿没碰的悬浮球通道标签）漂移了：got %q, want %q", ballLine.Text, "单击悬浮球")
	}
	// SplitN, not Index+"：" slicing: a byte index into a full-width colon cuts
	// the rune in half and the comparison then fails on mojibake, not on wording.
	parts := strings.SplitN(esc.Text, "：", 2)
	if len(parts) != 2 {
		t.Fatalf("未加载那行的形状不是「标签：不可用句」：%q", esc.Text)
	}
	why := parts[1]
	if why != "快捷键取消不可用" {
		t.Errorf("取消通道不可用句 = %q, want %q（档一：只说通道不可用，不指枚键名）", why, "快捷键取消不可用")
	}
	if strings.Contains(why, "Esc") {
		t.Errorf("取消通道不可用句还在指枚 Esc：%q", why)
	}
	// The three lines outside this leg's scope must still read exactly as before
	// (A595 §1 names only the two Esc lines in this file). All four channels are
	// unloaded here so each one prints its own availability sentence.
	for _, want := range []string{"语音取消不可用", "面板取消不可用（票 37 未接入）", "悬浮球取消不可用"} {
		found := false
		for _, st := range approval.NewChannels().Statuses() {
			if strings.Contains(st.Text, want) {
				found = true
			}
		}
		if !found {
			t.Errorf("其它通道的不可用句被改动了，名册里再没有 %q：%+v", want, approval.NewChannels().Statuses())
		}
	}
}

// TestTicket260R3LoadedCancelLabelIsTheNamedResidual is the ruler that keeps the
// one string this leg could not fix from drifting silently. See the header: this
// is a booked debt, and the leg that lands the fix flips this expectation on
// purpose (it must also flip the two gate.go / one report.go read faces, which is
// why it needs the orchestrator's word first).
func TestTicket260R3LoadedCancelLabelIsTheNamedResidual(t *testing.T) {
	reg := approval.NewChannels(approval.ChannelBall, approval.ChannelEsc)
	for _, st := range reg.Statuses() {
		if st.Channel != approval.ChannelEsc {
			continue
		}
		if !st.Loaded {
			t.Fatal("显式加载了取消通道却报成未加载")
		}
		if st.Text != "按 Esc 键" {
			t.Errorf("已加载的取消通道标签 = %q。若这一行是新形（票 260 条件②的收尾腿改的），"+
				"记得同时改 wording.md §8 那笔账与本件；若是无意的漂移，改回来。", st.Text)
		}
	}
}
