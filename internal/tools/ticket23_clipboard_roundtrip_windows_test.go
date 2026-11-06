//go:build windows

package tools

// Ticket 23 AC#4's one residual row: the REAL clipboard round trip.
//
// The ticket log (2026-11-03) names it as still open, and
// system_23_clipboard_memory_windows_test.go deliberately stops short of it:
// that file proves the Win32 memory copy moves bytes both directions, and its
// header says a round trip through the SYSTEM clipboard "is a separate,
// hand-run row". This file is that row, turned into a test - with one hard
// condition, because writing CF_UNICODETEXT calls EmptyClipboard, which
// destroys EVERY format on the clipboard, not just text:
//
//	WISP_T23_LIVE_CLIPBOARD=1
//
// Without it the file skips on every machine (the default gate run stays
// hands-off, exactly like the *winlive* rule). An operator who sets the
// variable has opted in; the leg still restores what it can and says out loud
// what it cannot: non-text formats (an image, a file drop) are destroyed by
// EmptyClipboard and can only be restored if the previous content was text -
// CF_UNICODETEXT is the one format this family reads and writes, so text
// before -> same text after; anything else before -> an empty clipboard after,
// never a fabricated replacement.
//
// The round trip runs through the SAME two production entry points the tools
// use (realClipboardWrite -> realClipboardRead), not through a private copy of
// their bodies, so a regression in either one is what turns this red. The
// second leg ties the round-tripped content into AC#4's judge half through the
// family's own C25 source name (SrcClipboardRead), which is the propagation
// the ticket's "taint propagation into R4" wording asks to survive a real
// clipboard.read.

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// liveClipboard23 reports whether the operator opted in.
func liveClipboard23() bool { return os.Getenv("WISP_T23_LIVE_CLIPBOARD") == "1" }

// saveClipboard23 reads the current text content ("" when there is none). A
// clipboard held by another process, or a session with no clipboard at all, is
// a "cannot run here", not a product failure: skip, never fail.
func saveClipboard23(t *testing.T) (prev string) {
	t.Helper()
	prev, err := realClipboardRead()
	if err != nil {
		t.Skipf("读不到当前剪贴板（%v）：本腿只在可独占剪贴板的会话里跑", err)
	}
	if !callBool(procOpenClipboard, 0) {
		t.Skipf("打不开剪贴板：本腿需要能独占剪贴板的会话")
	}
	callUintptr(procCloseClipboard)
	return prev
}

// restoreClipboard23 puts the text back when there was text, and empties the
// clipboard when there was not - the honest restore, since writing an empty
// string would leave a CF_UNICODETEXT where none was.
func restoreClipboard23(t *testing.T, prev string) {
	t.Helper()
	if prev != "" {
		if err := realClipboardWrite(prev); err != nil {
			t.Errorf("恢复原剪贴板内容失败（%q 还在测试内容里）：%v", prev, err)
		}
		return
	}
	if !callBool(procOpenClipboard, 0) {
		t.Errorf("恢复空剪贴板失败：OpenClipboard 打不开")
		return
	}
	defer callUintptr(procCloseClipboard)
	if !callBool(procEmptyClipboard) {
		t.Errorf("恢复空剪贴板失败：EmptyClipboard 被拒绝")
	}
}

// Test23LiveClipboardRoundTripThroughRealOS is the AC#4 row: write the marker
// through the production writer, read it back through the production reader,
// byte-for-byte, then restore what was there.
func Test23LiveClipboardRoundTripThroughRealOS(t *testing.T) {
	if !liveClipboard23() {
		t.Skip("WISP_T23_LIVE_CLIPBOARD!=1：真剪贴板往返是操作者点名才跑的腿（EmptyClipboard 会清掉全部格式）")
	}
	prev := saveClipboard23(t)
	marker := fmt.Sprintf("WISP-T23-ROUNDTRIP-%d·中文剪贴板·42%%", time.Now().UnixNano())
	if err := realClipboardWrite(marker); err != nil {
		t.Fatalf("写入剪贴板失败：%v", err)
	}
	back, err := realClipboardRead()
	if err != nil {
		restoreClipboard23(t, prev)
		t.Fatalf("读回剪贴板失败：%v", err)
	}
	if back != marker {
		restoreClipboard23(t, prev)
		t.Fatalf("往返不一致：读回 %q，写入 %q（真剪贴板没有把字节带到另一头）", back, marker)
	}
	if prev != "" && back == prev {
		// Only possible if the write silently no-op'd; the marker carries a
		// nanosecond stamp, so equality with the previous content is a lie.
		restoreClipboard23(t, prev)
		t.Fatalf("读回的是写入前的旧内容 %q：写入没有生效", prev)
	}
	restoreClipboard23(t, prev)
	if after, err := realClipboardRead(); err != nil || after != prev {
		t.Fatalf("恢复失败：读回 %q err=%v，期望 %q", after, err, prev)
	}
}

// Test23ClipboardReadSourceFeedsR4ThroughWrite closes AC#4's judge half on the
// family's OWN source name: Test23ClipboardWriteIsAFrozenC25Channel marks
// through SrcFSRead; this leg marks through SrcClipboardRead - the source a
// real clipboard.read books - and requires the same R4 exfil verdict on
// clipboard.write. No OS involvement: the taint ledger is pure memory, so this
// leg runs wherever the package builds.
func Test23ClipboardReadSourceFeedsR4ThroughWrite(t *testing.T) {
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	const scope, secret = "t23-clip-scope", "round-tripped-secret-23"
	if !prov.Mark(scope, risk.SrcClipboardRead, "", "剪贴板正文含 "+secret+" 结尾") {
		t.Fatalf("前置失败：SrcClipboardRead 的 Mark 没落污点")
	}
	hit, ok := prov.Inspect(scope, "clipboard.write", map[string]any{
		"text": "外发 " + secret,
	})
	if !ok {
		t.Fatalf("clipboard.read 打的污点经 clipboard.write 外发没有判 R4：读源与外泄通道没接上")
	}
	if hit.Channel != risk.ChClipboard {
		t.Errorf("命中通道 %q，要 %q（F4）", hit.Channel, risk.ChClipboard)
	}
	if hit.SrcTool != string(risk.SrcClipboardRead) {
		t.Errorf("命中来源 %q 没点名 clipboard.read（AC#4 要的就是这个源）", hit.SrcTool)
	}
}
