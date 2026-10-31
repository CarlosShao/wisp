package panel

// Ticket 102's rewrite account at the two display-boundary consumers the A380
// static criterion caught reading Result.Canonical without it (102r2: the hole
// re-opened by composer.go:127 and display_scrub_216.go). The static criterion
// only pins that the account's SYMBOL is read - ticket 102's own mutation
// reading (M1: "只挡符号消失，不挡条件被改") is why the behavioural halves get
// their own nails here:
//
//   - a rewritten row that arrives without narration leaves the boundary WITH
//     the account sentence, and never with 「没有读数」 stamped onto it (the
//     account proves the dimension was read);
//   - a rewritten row whose narration quotes the C26 coordinate crosses the
//     boundary washed, exactly like the coordinate itself.

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

func TestRewrittenWorkspaceRowLeavesTheBoundaryWithItsAccountSentence(t *testing.T) {
	rewritten := WorkspaceView{
		Set:       true,
		Spelling:  `%WISP102R2%\notes`,
		Canonical: `D:\elsewhere\notes`,
		Rewritten: true,
	}
	composer := NewComposerState(risk.ModeAskEveryStep, rewritten, nil, MaxAttachmentBytes)
	snap := NewSnapshot(nil, nil, composer, time.Now())

	ws := snap.Composer.Workspace
	if !ws.Rewritten {
		t.Fatalf("fixture premise broken: the account did not survive the boundary: %+v", ws)
	}
	if strings.TrimSpace(ws.Reason) == "" {
		t.Errorf("a rewritten row crossed the display boundary with no narration at all: %+v", ws)
	}
	if ws.Reason == workspaceNoReadingReason {
		t.Errorf("the no-reading sentence was stamped onto a row whose account proves the dimension was read: %q", ws.Reason)
	}
	if !strings.Contains(ws.Reason, "改写") {
		t.Errorf("the fallback narration does not say the rewrite out loud: %q", ws.Reason)
	}

	// The unchanged half: a row with no account and no reading still gets the
	// no-reading sentence, not the account one.
	unread := NewSnapshot(nil, nil, ComposerState{}, time.Now())
	if unread.Composer.Workspace.Reason != workspaceNoReadingReason {
		t.Errorf("the no-reading net moved: %q", unread.Composer.Workspace.Reason)
	}
}

func TestRewrittenNarrationQuotingTheCoordinateIsWashedAtTheBoundary(t *testing.T) {
	const hostile = `C:\data` + "\u202e" + `reports`
	rewritten := WorkspaceView{
		Set:       true,
		Spelling:  `%WISP102R2%\reports`,
		Canonical: hostile,
		Rewritten: true,
		Reason: "已收窄到该工作区，但这条路径被 C26 的展开步改写过（env）：下面的坐标 " +
			hostile + " 是改写后的树，不是使用者命名的那串 " + `%WISP102R2%\reports`,
	}
	composer := NewComposerState(risk.ModeAskEveryStep, rewritten, nil, MaxAttachmentBytes)
	snap := NewSnapshot(nil, nil, composer, time.Now())

	ws := snap.Composer.Workspace
	if strings.ContainsRune(ws.Reason, 0x202e) {
		t.Errorf("the rewritten narration carried the raw coordinate past the boundary: %q", ws.Reason)
	}
	if !strings.Contains(ws.Reason, `\u202e`) {
		t.Errorf("the narration lost the visible escape of the coordinate it quotes: %q", ws.Reason)
	}
	if got := ws.Canonical; strings.ContainsRune(got, 0x202e) || !strings.Contains(got, `\u202e`) {
		t.Errorf("the coordinate itself is no longer washed: %q", got)
	}

	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	if i := bytes.Index(data, []byte{0xe2, 0x80, 0xae}); i >= 0 {
		t.Errorf("the packet still holds a raw U+202E at offset %d: the account-gated wash is missing", i)
	}
}
