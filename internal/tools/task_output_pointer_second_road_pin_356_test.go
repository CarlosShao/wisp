package tools

import (
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 356 AC#2: a marker that belongs ONLY to the second road back, so that
// deleting that clause turns this test red. The 174-r1 marker set
// {"不在你被授权的目录范围内", "[fs] allowed_dirs", "L2 卡"} all survive deleting
// "要么批下那一张卡" (174-v1 MUT-5, reproduced): "L2 卡" and "allowed_dirs" live
// in the FIRST half of the sentence, so the second road can vanish silently.
// This pin names the tail clause itself, which exists nowhere else.
func Test356SecondRoadClauseIsPinned(t *testing.T) {
	outside := tempCanonical(t)
	body := asciiRun(20000)
	realFile := mustWriteFile(t, filepath.Join(outside, "tool-output-356-pin.txt"), body)

	roster := NewTaskRoster()
	roster.Record("bg-356-pin", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, realFile)})
	b := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer(nil, nil), NoGate{})

	out := callTaskOutput(t, b, "bg-356-pin")
	notice := announceOf(out.Text)
	// Second road ONLY. Deleting "要么批下那一张卡" (or any part that carries it)
	// MUST turn this red; that is the whole point of the pin.
	if !strings.Contains(out.Text, "要么批下那一张卡") {
		t.Fatalf("356 AC#2: second-road clause '要么批下那一张卡' missing; deleting it MUST turn this red. notice=%q", notice)
	}
	// First road still present: the two are a pair, not a slogan.
	if !strings.Contains(out.Text, "[fs] allowed_dirs") {
		t.Fatalf("356 AC#2: first-road marker '[fs] allowed_dirs' missing. notice=%q", notice)
	}
}
