//go:build windows

package main

// R382-1's readings (382-w-r2), taken through the production announce path the
// wire-r1 readings already use. The swap under test: speak()'s HOP 1 asks the
// TTSGate (CheckTextWithInheritance) instead of the bare Provenance check, so a
// sentence that quotes a CLOSED task scope's taint is refused where the old
// check could no longer see it. The staging contract is 15e's (sherpa DLLs next
// to the test binary, model files reachable): every case runs with
// -test.run=Test382WR2, never as part of the whole package (which dies at
// MARK299).
//
// The RecordTaintInheritance call below is the readings' stand-in for the
// production accumulation hook R382-1 names: the hook itself has no reachable
// seam in cmd/wisp (mark-time content never crosses internal/tools' bridge, and
// the engine's ScopeTaints is metadata-only by design), which the leg's evidence
// declares as a named stop - this file proves the mechanism, not the missing
// production feed.

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// closedTaskTaint is lowercase ASCII with no whitespace, which makes it its own
// normal form under the contract normalization (internal/risk/taintmatch.go's
// normalizeTaint: full-width folded, whitespace and zero-width dropped, case
// folded). RecordTaintInheritance's contract requires the CALLER to accumulate
// normalized content at mark time, and this package cannot call the unexported
// normalizer - so the link below is recorded with content that is verifiably
// already in the normalized coordinate space, without a second copy of the
// algorithm living here.
const closedTaskTaint = "ax5s3yl4apikey382wr2"

// ttsInheritLeg builds the leg the wire-r1 refusal case builds: the real
// announce path, a recording speaker for the device, and the process's C25
// engine attached - which since R382-1 also builds the leg's TTSGate.
func ttsInheritLeg(t *testing.T) (*residentVoice, *voiceWireSpeaker, *risk.Provenance) {
	t.Helper()
	dataDir := voice15eDataDir(t)
	voice15eWriteConfig(t, dataDir, nil)
	t.Setenv("WISP_MODELS_MANIFEST", voice15eRepoFile(t, filepath.Join("build", "models", "manifest.json")))
	store := voice15eInstallAsr(t, dataDir)

	speaker := &voiceWireSpeaker{rate: 22050}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	sink := &voice15eSink{}
	ball := &residentBall{}
	plan := voicePlan{
		VadModelPath: voiceVadNone(),
		AsrDir:       store,
		Source:       "382wr2 test",
		Tts:          voiceWireTtsFiles(t),
	}
	// No frames at all: these cases are about the gate, not about the
	// microphone, and a leg with no frame channel still owns an announce path.
	frames := make(chan []byte)
	rv := newResidentVoice(frames, plan, ball, sink)
	rv.residency = newVoiceResidency(plan, rv.queue, nil)
	rv.newSpeaker = func(rate int) (voiceSpeaker, error) { return speaker, nil }
	rv.attachProvenance(prov)
	t.Cleanup(rv.stop)
	return rv, speaker, prov
}

// Test382WR2TTSInheritanceRefusesClosedTaskTaint is the form-B refusal reading:
// a task scope's taint is recorded into the gate, the scope closes the way
// production closes it (admitTask's revocation), and a sentence quoting the
// taint is refused BEFORE synthesis - counted, named by source, and never
// loaded an engine.
func Test382WR2TTSInheritanceRefusesClosedTaskTaint(t *testing.T) {
	rv, speaker, prov := ttsInheritLeg(t)

	scope := prov.OpenScope("task-382wr2-closed")
	if scope == nil {
		t.Fatal("OpenScope refused the task scope")
	}
	if !prov.Mark("task-382wr2-closed", risk.SrcSearchContent, "search", closedTaskTaint) {
		t.Fatal("Mark refused the taint")
	}
	rv.ttsGate.RecordTaintInheritance("task-382wr2-closed", [][]rune{[]rune(closedTaskTaint)}, []string{risk.SrcSearchContent})
	if err := scope.Close(); err != nil {
		t.Fatalf("close the task scope: %v", err)
	}

	if !rv.announce("读到的内容：" + closedTaskTaint) {
		t.Fatalf("the queue refused the sentence, so the gate was never asked: %v", rv.problemsSnapshot())
	}
	rv.queue.run()
	voice15eWaitFor(t, "the gate to refuse the closed-task sentence", func() bool {
		return rv.ttsGated.Load() == 1
	})
	if got := rv.ttsGatedOk.Load(); got != 0 {
		t.Fatalf("the gate PASSED %d closed-task sentences; the inheritance scan is not on the path", got)
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
	t.Logf("INHERIT refused=%d passed=%d spoken=%d samples=%d slot_loads=%d",
		rv.ttsGated.Load(), rv.ttsGatedOk.Load(), rv.ttsSpoken.Load(), speaker.written(), rv.residency.slot.Loads())
}

// Test382WR2TTSCleanSentencePlaysThroughTheGate is the positive control: with
// the gate live AND armed with an inheritance link whose taint shares nothing
// with the sentence, a clean announcement still passes the gate, is synthesised
// and reaches the recording speaker.
func Test382WR2TTSCleanSentencePlaysThroughTheGate(t *testing.T) {
	rv, speaker, _ := ttsInheritLeg(t)

	rv.ttsGate.RecordTaintInheritance("task-382wr2-other", [][]rune{[]rune(closedTaskTaint)}, []string{risk.SrcSearchContent})

	if !rv.announce("语音播报自检：今天没有需要提醒的事项。") {
		t.Fatalf("the queue refused a clean sentence: %v", rv.problemsSnapshot())
	}
	rv.queue.run()
	voice15eWaitFor(t, "the clean sentence to be spoken", func() bool {
		return rv.ttsSpoken.Load() == 1
	})
	if got := rv.ttsGated.Load(); got != 0 {
		t.Fatalf("the gate refused %d clean sentences", got)
	}
	if got := rv.ttsGatedOk.Load(); got != 1 {
		t.Fatalf("the gate passed %d sentences, want 1", got)
	}
	if speaker.written() == 0 {
		t.Fatal("the clean sentence produced no samples on the device")
	}
	t.Logf("CLEAN passed=%d refused=%d spoken=%d samples=%d",
		rv.ttsGatedOk.Load(), rv.ttsGated.Load(), rv.ttsSpoken.Load(), speaker.written())
}

// Test382WR2TTSOpenScopeDirectCheckStillRefuses pins the half of the gate the
// swap must not lose: a taint still live in the announce scope is refused
// through the TTSGate's direct check, exactly as the old bare CheckText did.
func Test382WR2TTSOpenScopeDirectCheckStillRefuses(t *testing.T) {
	rv, speaker, prov := ttsInheritLeg(t)

	const secret = "AKIAIOSFODNN7EXAMPLE"
	if !prov.Mark(voiceAnnounceScope, risk.SrcFSRead, `D:\work\secret.txt`, "凭据是 "+secret) {
		t.Fatal("Mark refused the taint")
	}

	if !rv.announce("读到的内容：凭据是 " + secret) {
		t.Fatalf("the queue refused the sentence, so the gate was never asked: %v", rv.problemsSnapshot())
	}
	rv.queue.run()
	voice15eWaitFor(t, "the gate to refuse the live-scope sentence", func() bool {
		return rv.ttsGated.Load() == 1
	})
	if got := rv.ttsGatedOk.Load(); got != 0 {
		t.Fatalf("the gate PASSED %d live-scope sentences; the direct check is not on the path", got)
	}
	if rv.ttsSpoken.Load() != 0 || speaker.written() != 0 {
		t.Fatalf("a live-scope sentence reached the device: spoken=%d samples=%d",
			rv.ttsSpoken.Load(), speaker.written())
	}
	found := false
	for _, p := range rv.problemsSnapshot() {
		if strings.Contains(p, string(risk.ChTTS)) {
			found = true
		}
	}
	if !found {
		t.Fatalf("the refusal was not recorded with its channel: %v", rv.problemsSnapshot())
	}
	t.Logf("DIRECT refused=%d passed=%d spoken=%d samples=%d",
		rv.ttsGated.Load(), rv.ttsGatedOk.Load(), rv.ttsSpoken.Load(), speaker.written())
}
