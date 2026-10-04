package approval_test

// Ticket 260-r4 (ledger A598 §2): the cancel channel's LABEL follows the key this
// host actually borrows, and it is resolved when the line is rendered.
//
// WHAT WAS BROKEN, in the user's words. 260-r3 put the configured key into nine
// of the ten "按 Esc" sentences and stopped on the tenth - the card's own channel
// line - because naming it needed a value this package cannot see. Leaving it
// produced a NEW contradiction the orchestrator refused to book overnight:
// approval.go's unloaded branch renders 「标签：不可用句」, so the line read
// 「按 Esc 键：快捷键取消不可用」 - the first clause pointing at a key this machine
// never holds, the second saying no such key exists. Exactly the moment a user
// reads it (the card did not come up / the key was never borrowed).
//
// THE FOUR CELLS THIS FILE OWES (dispatch ①②③④, ⛔ not mergeable into one):
//	① default 档 renders 「按 Esc 键」 byte-for-byte as before;
//	② a host that installed a reader naming another key renders THAT key, and
//	  "Esc" disappears from every one of the four faces;
//	③ the line is not cached: what the host's reader answers at render time is
//	  what prints (the ticket-258 reload chain, seen from this side);
//	④ the UNLOADED line names no key at all, in either 档 - the chosen honest
//	  form, see TestTicket260R4UnloadLineNamesNoKey's header.
// Plus one boundary cell: this seam carries NO authority (a fake reader cannot
// make a channel cancellable), which is what makes a package-level injection
// acceptable in a security layer at all.
//
// WHY THE READER IS INJECTED AND NOT IMPORTED: internal/agent/approval must not
// gain a dependency edge on internal/ball (human-approval face, AGENTS.md §1.2 /
// A598 §2 边界), so the composition root hands a function - the shape ticket 246
// used for the gate and ticket 253 for panel-as-truth-source. cmd/wisp installs
// its ONE existing reader (the ball's hotkey receipt), so no second spelling
// path opens here either.

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

const (
	// oldR4EscLabel is the string this package printed for the loaded cancel
	// channel since before ticket 260 existed. AC#4 ①'s zero-drift half: a user
	// who never touched [hotkey] cancel must not see one character move.
	oldR4EscLabel = "按 Esc 键"
	// r4SeededKey / r4ReloadedKey are two spellings no code in this package can
	// produce by itself - if any face still hard-codes a key, cell ② goes red.
	r4SeededKey   = "Ctrl+Alt+Q"
	r4ReloadedKey = "Ctrl+Alt+W"
	// r4SlotLabel is 档一's front half: the SLOT, never a key. Pinned verbatim so
	// a later leg that "improves" it back into a key name is a red test.
	r4SlotLabel = "取消键通道"
)

// useCancelKeySeam installs a reader for one test and puts the incumbent back
// afterwards. The restoration is not decoration: this seam is process-wide, and
// a leaked reader would wash the readings of every later case in this package
// (including 260-r3's named-residual ruler).
func useCancelKeySeam(t *testing.T, src approval.CancelKeySpelling) {
	t.Helper()
	prev := approval.SetCancelKeySpelling(src)
	t.Cleanup(func() { approval.SetCancelKeySpelling(prev) })
}

func escStatus(t *testing.T, reg *approval.ChannelRegistry) approval.ChannelStatus {
	t.Helper()
	for _, st := range reg.Statuses() {
		if st.Channel == approval.ChannelEsc {
			return st
		}
	}
	t.Fatal("名册里没有取消通道那一行")
	return approval.ChannelStatus{}
}

// TestTicket260R4DefaultLabelIsTheOldSentence is cell ① on the SHIPPED fallback
// path: no reader installed (a host with no hotkey chain), so the label is the
// product default - and it is byte-for-byte the old string, including the two
// characters of padding around the key name.
func TestTicket260R4DefaultLabelIsTheOldSentence(t *testing.T) {
	useCancelKeySeam(t, nil)
	reg := approval.NewChannels(approval.ChannelBall, approval.ChannelEsc)
	st := escStatus(t, reg)
	if !st.Loaded {
		t.Fatal("显式加载了取消通道却报成未加载")
	}
	if st.Text != oldR4EscLabel {
		t.Errorf("默认档已加载的取消通道标签漂移：got %q, want %q", st.Text, oldR4EscLabel)
	}
	// The three labels this leg had no licence to touch.
	for _, want := range []string{"单击悬浮球", "面板拒绝", "说取消词"} {
		found := false
		for _, x := range approval.NewChannels(approval.ChannelBall, approval.ChannelPanel, approval.ChannelKWS).Statuses() {
			if x.Text == want {
				found = true
			}
		}
		if !found {
			t.Errorf("其它通道的标签被改动了，名册里再没有 %q", want)
		}
	}
}

// TestTicket260R4SeededKeyReplacesEscOnEveryFace is cell ② across ALL FOUR faces:
// the card's availability line (Statuses), the L1 window's veto sentence
// (gate.go), the L2 card's veto sentence (gate.go) and the applied-steps report's
// channel attribution (report.go). Re-hardcoding a literal in any one of them -
// or leaving one reading the old map - reddens that sub-case by name.
func TestTicket260R4SeededKeyReplacesEscOnEveryFace(t *testing.T) {
	useCancelKeySeam(t, func() string { return r4SeededKey })

	// Face 1: the card / ball strip.
	line := escStatus(t, approval.NewChannels(approval.ChannelBall, approval.ChannelEsc))
	assertNamesSeededKey(t, "卡片通道行", line.Text)

	// Face 2: the L1 window's veto sentence.
	ui := newFakeUI()
	g, _, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	res := runWindow(t, g, context.Background(), l1Decision("C:/data/r4-l1.txt"))
	ui.wait(t)
	if err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelEsc}); err != nil {
		t.Fatalf("L1 否决被拒：%v", err)
	}
	assertNamesSeededKey(t, "L1 否决审计句", mustAnswer(t, res).why)

	// Face 3: the L2 card's veto sentence.
	ui2 := newFakeUI()
	g2 := ticket87Gate(t, ui2)
	d := l2Decision("C:/elsewhere/r4-l2.txt")
	d.CorrelationID = "corr-r4-l2"
	res2 := runApproval(t, g2, context.Background(), d)
	p := ui2.wait(t)
	if err := g2.Veto(approval.Veto{CorrelationID: p.CorrelationID, Channel: approval.ChannelEsc}); err != nil {
		t.Fatalf("L2 卡片否决被拒：%v", err)
	}
	assertNamesSeededKey(t, "L2 否决审计句", mustAnswerWithin(t, res2, vetoBudget, "r4 L2 veto").why)

	// Face 4: D31's applied-steps report (the exported re-render path).
	rep := approval.CancellationReport{
		Tool: "fs.write", CorrelationID: "corr-r4-report", Level: "L1",
		Vetoed: true, Channel: approval.ChannelEsc, StartedBeforeVeto: true,
		AppliedSteps: []string{"已写入 C:/data/r4.txt"},
	}
	assertNamesSeededKey(t, "回执报告的否决通道行", rep.TextFor())
}

// TestTicket260R4LabelFollowsTheReaderAcrossAReload is cell ③ - as far as this
// package can be measured. What ticket 258's bridge changes on a hot reload is
// the ball's receipt, i.e. what the injected reader ANSWERS; from this side the
// defect it must not leave behind is a label computed once. So the reader here is
// backed by a holder a test can move, and the assertion is that two renders on
// either side of the move print two different keys, plus that the reader is
// actually CALLED per render (a value frozen into a var answers the first render
// and every later one, which is what this counts against).
//
// ⛔ What this does NOT prove, stated so nobody reads it as proven: that a real
// rebind rewrites the receipt's cancel line. That half is internal/ball's and
// needs a real window (winlive, not run by this leg) - booked as 〔未实测〕 in
// .scratch/wisp/probes/260/r4/label.md §8.
func TestTicket260R4LabelFollowsTheReaderAcrossAReload(t *testing.T) {
	var mu sync.Mutex
	held := "Esc"
	calls := 0
	useCancelKeySeam(t, func() string {
		mu.Lock()
		defer mu.Unlock()
		calls++
		return held
	})
	reg := approval.NewChannels(approval.ChannelBall, approval.ChannelEsc)

	before := escStatus(t, reg).Text
	if before != oldR4EscLabel {
		t.Errorf("热重载前那行 = %q, want %q", before, oldR4EscLabel)
	}
	mu.Lock()
	held = r4ReloadedKey
	mu.Unlock()

	after := escStatus(t, reg).Text
	if after != "按 "+r4ReloadedKey+" 键" {
		t.Errorf("热重载后那行陈旧了：got %q, want %q（同一台机器换了配置，标签必须跟着走）",
			after, "按 "+r4ReloadedKey+" 键")
	}
	if after == before {
		t.Errorf("重载前后一行字没变，说明标签被算了一次：%q", after)
	}
	if calls < 2 {
		t.Errorf("取值函数被调用 %d 次，少于两枚渲染所需的次数：标签不是在读取时求值的", calls)
	}
}

// TestTicket260R4UnloadLineNamesNoKey is cell ④, and the shape choice is stated
// rather than implied. Two candidate forms existed: (甲) the front half names
// THIS host's key even when the channel is unloaded, (乙) the front half names
// the SLOT and no key. ⛔ 甲 is not buildable honestly: the one source this
// process has for the key is the ball's hotkey receipt, and the receipt exists
// only once a window exists - the no-window host therefore falls back to the
// product default, so 甲 would print 「按 Esc 键：快捷键取消不可用」 in exactly
// the branch this cell is about, i.e. it re-creates the contradiction. 乙 is what
// landed, and it is the same rule 260-r3 applied to the two 档一 sentences in
// cmd/wisp (取消键无处可借 / 取消键通道保持未加载).
//
// The judgement therefore is: THE FRONT HALF NAMES NO KEY IN EITHER 档 - neither
// the default nor the configured one - while the loaded half keeps naming the
// key. Both 档 are checked, and the line keeps the 「标签：不可用句」 shape.
func TestTicket260R4UnloadLineNamesNoKey(t *testing.T) {
	for _, seeded := range []string{"", r4SeededKey} {
		spelling := seeded
		useCancelKeySeam(t, func() string { return spelling })
		st := escStatus(t, approval.NewChannels())
		if st.Loaded {
			t.Fatal("NewChannels() 却把取消通道报成已加载")
		}
		parts := strings.SplitN(st.Text, "：", 2)
		if len(parts) != 2 {
			t.Fatalf("未加载那行的形状不再是「标签：不可用句」：%q", st.Text)
		}
		front, why := parts[0], parts[1]
		if front != r4SlotLabel {
			t.Errorf("未加载那行的前半不再是通道名：got %q, want %q（档一：这一支没有键可指）", front, r4SlotLabel)
		}
		for _, banned := range []string{"Esc", r4SeededKey, "按"} {
			if strings.Contains(front, banned) {
				t.Errorf("未加载那行的前半还在指枚键名（含 %q）：%q", banned, st.Text)
			}
		}
		if why != "快捷键取消不可用" {
			t.Errorf("不可用半句被改动：got %q, want %q（那是 260-r3 的档一句，本腿不许动）", why, "快捷键取消不可用")
		}
		// The line must still be a REFUSAL-shaped statement, i.e. nothing in it
		// may read as "loaded": the loaded form must not appear in the unloaded
		// line, and vice versa.
		if strings.Contains(st.Text, oldR4EscLabel) {
			t.Errorf("未加载那行还带着已加载形：%q", st.Text)
		}
	}
	// And the loaded branch of the same registry says the OPPOSITE, so the two
	// forms cannot be collapsed into one by a later edit.
	useCancelKeySeam(t, func() string { return r4SeededKey })
	loaded := escStatus(t, approval.NewChannels(approval.ChannelEsc))
	if !strings.Contains(loaded.Text, r4SeededKey) || strings.Contains(loaded.Text, "不可用") {
		t.Errorf("已加载那行不再是「念当前那枚键」：%q", loaded.Text)
	}
}

// TestTicket260R4SpellingSeamCarriesNoAuthority is the boundary that makes a
// package-level injection acceptable inside this layer at all: the seam names a
// key for a human, and NOTHING may consume it as an availability answer. An
// unloaded channel stays uncancellable with a reader installed, the warning the
// user gets still names no key, and no answer on the allow side is reachable.
func TestTicket260R4SpellingSeamCarriesNoAuthority(t *testing.T) {
	useCancelKeySeam(t, func() string { return "F13" })
	ui := newFakeUI()
	g, _, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	res := runWindow(t, g, context.Background(), l1Decision("C:/data/r4-auth.txt"))
	p := ui.wait(t)

	st, ok := channelText(p, approval.ChannelKWS)
	if !ok {
		t.Fatal("提示里没有语音通道那一行")
	}
	if st.Loaded {
		t.Fatal("装了一枚读口就把未加载通道报成可用")
	}
	err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelKWS})
	if err == nil {
		t.Fatal("SECURITY: 装了读口之后，未加载的语音通道竟然否决成功")
	}
	if !errors.Is(err, approval.ErrChannelUnavailable) {
		t.Errorf("未加载通道的否决 = %v, want ErrChannelUnavailable", err)
	}
	// The window is still live and still must not be answerable as an allow.
	select {
	case got := <-res:
		t.Fatalf("未加载通道的否决竟然答复了窗口：%+v（窗口必须继续倒计时，这是 B1 的形）", got)
	default:
	}
	if n := g.Queue().Depth(); n != 0 {
		t.Errorf("queue depth = %d, want 0 (this route is an L1 window)", n)
	}
	for _, e := range ui.ofKind(approval.EventWarning) {
		if strings.Contains(e.Text, "F13") {
			t.Errorf("不可用回执念出了一枚根本没加载的键：%q", e.Text)
		}
	}

	// A reader naming a key this host never registered cannot answer the live
	// window either: only the loaded channel can, and it answers with a refusal.
	if err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelEsc}); err != nil {
		t.Fatalf("loaded-channel veto: %v", err)
	}
	if got := mustAnswer(t, res); got.a == tools.AnswerAllow {
		t.Errorf("SECURITY: a veto produced %v", got.a)
	}

	// A reader that answers nothing (empty string) must not print 「按  键」.
	useCancelKeySeam(t, func() string { return "" })
	if got := escStatus(t, approval.NewChannels(approval.ChannelEsc)).Text; got != oldR4EscLabel {
		t.Errorf("读口回空串时标签 = %q, want 落回出厂那枚 %q", got, oldR4EscLabel)
	}
}

// ------------------------------------------------------------------ helpers

func assertNamesSeededKey(t *testing.T, face, got string) {
	t.Helper()
	if !strings.Contains(got, r4SeededKey) {
		t.Errorf("%s 没有念出装配根那枚键：%q（应含 %q）", face, got, r4SeededKey)
	}
	if strings.Contains(got, "Esc") {
		t.Errorf("%s 还在念 Esc：%q", face, got)
	}
}
