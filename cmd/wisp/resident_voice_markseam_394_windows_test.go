//go:build windows

package main

// 票 394-r1: the production mark feed R382-1's readings stood in for. The
// 382wr2 file below had to record its link BY HAND with a hand-picked
// normalized constant, because mark-time content never crossed internal/tools'
// bridge and the engine's ScopeTaints is metadata-only by design. The 394 seam
// (internal/risk/markseam_394.go) reclaims that obligation inside risk: the
// gate borrows the SAME engine attachProvenance hands it, and every mark that
// engine commits lands in the inheritance chain already normalized.
//
// This test is the same rig and the same reading shape as
// Test382WR2TTSInheritanceRefusesClosedTaskTaint — with TWO differences that
// are the whole point:
//  1. NO RecordTaintInheritance call anywhere: the chain is fed by a plain
//     prov.Mark, the funnel all four production mark sites go through;
//  2. the marked content is RAW production-shaped text (full-width, spaces,
//     mixed case) while the announcement quotes the HALF-WIDTH spelling — the
//     raw forms share nothing, so only the seam's normalized copy can match.
//
// The staging contract is 15e's (sherpa DLLs next to the test binary, model
// files reachable): run with -test.run Test394R1, never as part of the whole
// package (which dies at MARK299).

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// mark394Raw is the raw content a production mark would carry: full-width
// Latin/digits, ideographic spaces. Its normal form is what the seam stores.
const mark394Raw = "查询结果包含　ＡＰＩＫＥＹ３９４Ｒ１ＸＹＺＱ　请勿外传"

// announce394Quote quotes the taint's normalized (half-width) spelling — raw
// vs raw the two texts share nothing.
const announce394Quote = "已读取内容：apikey394r1xyzq，请注意保管"

// Test394R1SeamRecordsProductionMark is the feed half: the production mark
// alone (no manual record) puts exactly one link in the gate's chain.
func Test394R1SeamRecordsProductionMark(t *testing.T) {
	rv, _, prov := ttsInheritLeg(t)

	scope := prov.OpenScope("task-394r1-closed")
	if scope == nil {
		t.Fatal("OpenScope refused the task scope")
	}
	if !prov.Mark("task-394r1-closed", risk.SrcSearchContent, "web", mark394Raw) {
		t.Fatal("Mark refused the raw taint")
	}
	if got := rv.ttsGate.CountInheritanceLinks(); got != 1 {
		t.Fatalf("the seam recorded %d links, want exactly 1 from the production mark alone (no manual record)", got)
	}
}

// Test394R1SeamFedRefusalRefusesClosedTaskPoison is the AC#2 reading: a closed
// task scope's production-marked taint refuses the announcement BEFORE
// synthesis — the INHERIT shape, fed by the seam instead of a manual record.
func Test394R1SeamFedRefusalRefusesClosedTaskPoison(t *testing.T) {
	rv, speaker, prov := ttsInheritLeg(t)

	scope := prov.OpenScope("task-394r1-closed")
	if scope == nil {
		t.Fatal("OpenScope refused the task scope")
	}
	if !prov.Mark("task-394r1-closed", risk.SrcSearchContent, "web", mark394Raw) {
		t.Fatal("Mark refused the raw taint")
	}
	if err := scope.Close(); err != nil {
		t.Fatalf("close the task scope: %v", err)
	}

	if !rv.announce(announce394Quote) {
		t.Fatalf("the queue refused the sentence, so the gate was never asked: %v", rv.problemsSnapshot())
	}
	rv.queue.run()
	voice15eWaitFor(t, "the seam-fed gate to refuse the closed-task sentence", func() bool {
		return rv.ttsGated.Load() == 1
	})
	if got := rv.ttsGatedOk.Load(); got != 0 {
		t.Fatalf("the gate PASSED %d closed-task sentences; the seam feed is not on the path", got)
	}
	if rv.ttsSpoken.Load() != 0 || speaker.written() != 0 {
		t.Fatalf("a closed-task sentence reached the device: spoken=%d samples=%d",
			rv.ttsSpoken.Load(), speaker.written())
	}
	if got := rv.residency.slot.Loads(); got != 0 {
		t.Fatalf("the TTS engine loaded %d times for a refused sentence; the gate must run BEFORE synthesis", got)
	}
	found := false
	for _, p := range rv.problemsSnapshot() {
		if strings.Contains(p, string(risk.ChTTS)) && strings.Contains(p, "from closed task scope") {
			found = true
		}
	}
	if !found {
		t.Fatalf("the refusal was not recorded with its form-B source: %v", rv.problemsSnapshot())
	}
	t.Logf("SEAM-PRODUCTION refused=%d passed=%d spoken=%d samples=%d slot_loads=%d links=%d",
		rv.ttsGated.Load(), rv.ttsGatedOk.Load(), rv.ttsSpoken.Load(), speaker.written(),
		rv.residency.slot.Loads(), rv.ttsGate.CountInheritanceLinks())
}
