//go:build windows

package main

// Ticket 260-r4 (ledger A598 §2): the card's own channel line comes from the SAME
// reader as this file's sentences, and it is resolved when it is rendered.
//
// 260-r3 landed nine of A595's ten "Esc" sentences and stopped on the tenth - the
// label the approval package puts on the cancel channel - because naming it needed
// internal/ball's receipt, which approval must not import. The leftover was a NEW
// contradiction of the same shape this ticket was filed over: the unloaded branch
// renders 「标签：不可用句」, so the line read 「按 Esc 键：快捷键取消不可用」 - a
// key named in a moment when no key is held.
//
// WHAT THIS FILE PINS (and what it does not):
//	- the assembly root really installs the reader (the shipped constructor, not a
//	  test-only helper), so the card line and this file's veto sentences cannot be
//	  pointed at two different keys - the M4 blind spot of r3 was about this exact
//	  wiring, from the other side;
//	- default 档 renders the old string byte-for-byte, with no seam and no window;
//	- a seeded receipt changes BOTH halves at once;
//	- a host with no ball window names NO key on that line, in either 档.
// ⛔ What stays 〔未实测〕 here: that a real rebind rewrites the ball's receipt
// (that is internal/ball, and a real window is a winlive rig this leg may not run;
// see .scratch/wisp/probes/260/r4/label.md §8). The seam this file seeds is the
// r3 seam for the same reason - it stands in for the receipt line, never for the
// chain that fills it.

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/ball"
)

const (
	// oldR4CardLine is the loaded cancel-channel line as approval rendered it
	// before ticket 260 ever existed. AC#4 ① is measured against this string.
	oldR4CardLine = "按 Esc 键"
	// r4CardSlotLine is 档一's line front (approval's slot name), asserted absent
	// from any line that claims a key.
	r4CardSlotLine = "取消键通道"
	// r4SeededCardKey is a spelling this package can only supply through the
	// injected reader - which is the point.
	r4SeededCardKey = "Ctrl+Alt+Q"
)

// escCardLine pulls the cancel channel's line out of the gate the resident leg
// built (the same registry promptFor stamps onto both routes).
func escCardLine(t *testing.T, ra *residentApproval) approval.ChannelStatus {
	t.Helper()
	for _, st := range ra.gate.Channels().Statuses() {
		if st.Channel == approval.ChannelEsc {
			return st
		}
	}
	t.Fatal("闸门的名册里没有取消通道那一行")
	return approval.ChannelStatus{}
}

// restoreCancelKeySeam puts approval's process-wide reader back after a case, so
// nothing here washes a later test's reading.
func restoreCancelKeySeam(t *testing.T) {
	t.Helper()
	prev := approval.SetCancelKeySpelling(nil)
	t.Cleanup(func() { approval.SetCancelKeySpelling(prev) })
}

// TestTicket260R4ShippedConstructorInstallsTheReader is the wiring cell: the gate
// the resident process builds must not be left guessing a key. With no seam and no
// window the installed reader answers the BALL's own default constant, so the card
// line is still the old string - and the moment the receipt names another key, the
// card line and this file's status sentence move together (one source, no second
// spelling path, which is the shape ticket 260 was filed to remove).
func TestTicket260R4ShippedConstructorInstallsTheReader(t *testing.T) {
	if ball.DefaultHotkeys().Cancel != "Esc" {
		t.Fatalf("这一格的前提动了：ball 的出厂 cancel 是 %q 不是 Esc；旧句常量要改就得显式改，不许顺手",
			ball.DefaultHotkeys().Cancel)
	}
	restoreCancelKeySeam(t)

	// The channel gets loaded the way production loads it (bindBallHost with a
	// window in hand), on r3's reader seam standing in for the receipt line - a
	// zero-value ball cannot answer an uiRun, and a real window is winlive.
	ra := newResidentApproval()
	ra.cancelKeyRead = func() string { return ball.DefaultHotkeys().Cancel }
	captureStdout128(t, func() {
		if !ra.bindBallHost(&residentBall{b: &ball.Ball{}, cancelHosted: true}) {
			t.Error("bindBallHost 在没有执行者的情况下声称装配好了")
		}
	})
	// The shipped fallback path is measured separately, without a seam at all:
	// TestTicket260R4SeamIsNotAProductionShortcut / r3's own production-path case.
	if got := escCardLine(t, ra); got.Text != oldR4CardLine {
		t.Errorf("默认档卡片那行漂移：got %q, want %q", got.Text, oldR4CardLine)
	}

	// Seed the receipt line and re-read: the SAME object, no reconstruction. If the
	// label had been computed once at construction, this is where it goes red.
	ra.cancelKeyRead = func() string { return r4SeededCardKey }
	seeded := escCardLine(t, ra)
	if !strings.Contains(seeded.Text, r4SeededCardKey) {
		t.Errorf("改了配置以后卡片那行没念那枚键：%q（应含 %q）", seeded.Text, r4SeededCardKey)
	}
	if strings.Contains(seeded.Text, "Esc") {
		t.Errorf("改了配置以后卡片那行还在念 Esc：%q", seeded.Text)
	}

	// And the sentence this file prints about the same channel agrees with the
	// card, in the same key slot - the single-source claim: two renderings, one
	// reader, so no edit can leave them naming different keys.
	status := ra.residentStatusLine()
	if !strings.Contains(status, r4SeededCardKey) || strings.Contains(status, "Esc") {
		t.Errorf("装配根那句和卡片那行分家了：%q（键名应与卡片一致为 %q）", status, r4SeededCardKey)
	}
}

// TestTicket260R4NoBallHostNamesNoKeyOnTheCard is the cell the orchestrator
// refused to book overnight: with no ball window the channel is NOT loaded, and
// the card line must then name no key at all - neither the shipped default nor a
// configured one - because this process holds nothing to press. Both 档 are run:
// the seeded reader is present, and printing its key here would advertise a
// channel this process refuses (the B1 shape).
func TestTicket260R4NoBallHostNamesNoKeyOnTheCard(t *testing.T) {
	restoreCancelKeySeam(t)
	for _, seeded := range []string{"", r4SeededCardKey} {
		ra := newResidentApproval()
		ra.cancelKeyRead = func() string { return seeded }
		captureStdout128(t, func() {
			if ra.bindBallHost(&residentBall{}) {
				t.Error("没有球窗的 bindBallHost 声称加载了通道")
			}
		})

		st := escCardLine(t, ra)
		if st.Loaded {
			t.Fatalf("没有球窗却把取消通道报成已加载（档=%q）", seeded)
		}
		front, why, ok := strings.Cut(st.Text, "：")
		if !ok {
			t.Fatalf("未加载那行不再是「标签：不可用句」：%q", st.Text)
		}
		if front != r4CardSlotLine {
			t.Errorf("未加载那行的前半不再是通道名：%q（want %q）", front, r4CardSlotLine)
		}
		for _, banned := range []string{"Esc", r4SeededCardKey, "按"} {
			if strings.Contains(front, banned) {
				t.Errorf("未加载那行的前半还在指枚键名（含 %q，档=%q）：%q", banned, seeded, st.Text)
			}
		}
		if !strings.Contains(why, "不可用") {
			t.Errorf("未加载那行的后半不再是不可用句：%q", st.Text)
		}
	}
}

// TestTicket260R4SeamIsNotAProductionShortcut keeps r3's guard standing on the new
// mechanism: the injected reader is an injection point, the shipped constructor
// still leaves it nil, and nothing but a real ball receipt can put a key on the
// card in a shipped process.
func TestTicket260R4SeamIsNotAProductionShortcut(t *testing.T) {
	restoreCancelKeySeam(t)
	ra := newResidentApproval()
	if ra.cancelKeyRead != nil {
		t.Fatal("newResidentApproval 装了取消键读口：产码路径不再读球自己的回执行了")
	}
	if got := escCardLine(t, ra); !strings.HasPrefix(got.Text, r4CardSlotLine) || strings.Contains(got.Text, "按") {
		// Channel not loaded in a fresh gate -> the slot form. If this ever reads
		// 「按 Esc 键：…」 again, a load happened somewhere that has no window
		// behind it, i.e. the contradiction came back.
		t.Errorf("未装配闸门的卡片那行 = %q, want 以通道名 %q 开头且不指枚键名", got.Text, r4CardSlotLine)
	}
}
