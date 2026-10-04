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
// The loaded-channel label (the 第十枚 of A595's ten) is the ONE string 260-r3 did
// NOT change, and it stayed a named residual for exactly one day: making it follow
// the config needed a value this package cannot see without importing internal/ball
// (a new package-level dependency edge, which AGENTS.md §1.2 books as a
// human-approval face). The orchestrator ruled it in ledger A598 §2 with a named
// unfreeze of five things, and 260-r4 landed it: the label is now resolved at READ
// time from a function the composition root installs (approval.SetCancelKeySpelling),
// so 「按 <key> 键」 names the key this host borrows.
//
// The assertion below therefore keeps its byte-exact expectation for the DEFAULT
// 档 (a host that installed no reader still renders the old string, AC#4 ①) and
// gained the counterpart that makes it a moving ruler instead of a frozen one:
// with a reader installed, the same line MUST NOT be that string. That is not a
// relaxation of anything - it is one assertion plus the half that used to be
// missing. See .scratch/wisp/probes/260/r4/label.md §5 for the cell-by-cell map.

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

// TestTicket260R3LoadedCancelLabelIsTheNamedResidual is the ruler 260-r3 left
// behind for the one string it could not fix. The NAME stays because r3's own
// evidence logs cite it (.scratch/wisp/probes/260/r3/logs/approval-*.txt) - the
// residual it names is closed, and what this case pins now is the DEFAULT 档 of
// the new read-time mechanism plus the half that was missing: install a reader,
// and the same line must stop being that string.
func TestTicket260R3LoadedCancelLabelIsTheNamedResidual(t *testing.T) {
	// 260-r4's counterpart first: with a reader installed by the host, this line
	// is NOT the old string. Without this half the case below would be satisfied
	// by a label frozen back into a constant.
	useCancelKeySeam(t, func() string { return "Ctrl+Alt+Q" })
	seeded := ""
	for _, st := range approval.NewChannels(approval.ChannelBall, approval.ChannelEsc).Statuses() {
		if st.Channel == approval.ChannelEsc {
			seeded = st.Text
		}
	}
	if seeded == "按 Esc 键" {
		t.Errorf("装配根装了取消键读口，卡片那行还在念出厂默认：%q（读取时求值被算成了一次）", seeded)
	}

	// Then the default 档, on the shipped fallback path (no reader).
	useCancelKeySeam(t, nil)
	reg := approval.NewChannels(approval.ChannelBall, approval.ChannelEsc)
	for _, st := range reg.Statuses() {
		if st.Channel != approval.ChannelEsc {
			continue
		}
		if !st.Loaded {
			t.Fatal("显式加载了取消通道却报成未加载")
		}
		if st.Text != "按 Esc 键" {
			t.Errorf("默认档的取消通道标签 = %q。这一格钉的是「用户没配 [hotkey] cancel 时一字不动」"+
				"（票 260 AC#4 ①）；若这是有意的换形，改 wording.md §5 与 r4 件 §5 的账再动它。", st.Text)
		}
	}
}
