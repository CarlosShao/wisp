package main

// Ticket 167 AC#5's teeth: a failed send keeps its own words, and those words
// never overwrite anything the user wrote afterwards.
//
// The two directions the dispatch names, as measurable Go facts:
//
//	失败时整块回填覆盖新字 ⇒ 红   - the failure path writes ONLY its own slot;
//	                               every other composer field is byte-identical
//	                               across a failed send, so a mutation that
//	                               backfills the failed text into any other field
//	                               (the "whole-block refill" shape) reddens the
//	                               field-diff pin below.
//	正控：不失败时草稿不被抹      - a SUCCEEDING send leaves the stored original
//	                               untouched; only a later failure replaces it,
//	                               with that failure's OWN text, never a merge.
//
// The wire half is pinned too: both keys are omitempty, so "no send face" is
// absence on the wire, and the reader that exists always sends the flag.

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// errRelayDown167 is the teeth's own refusal: the relay the failure branches
// are driven with.
var errRelayDown167 = errors.New("167: relay down")

func sendReq167(requestID, text string) panel.ComposerRequest {
	return panel.ComposerRequest{
		Method: panel.MethodMessageSend, RequestID: requestID,
		Source: panel.ComposerRequestSource, Text: text,
	}
}

// Test167AC5FailedSendKeepsItsOwnWordsSeparately: the original is booked
// verbatim on the failure branches, survives a later SUCCESS (the positive
// control), and is replaced - not merged - by a later failure's own text.
func Test167AC5FailedSendKeepsItsOwnWordsSeparately(t *testing.T) {
	door := newPanelMessageDoor453(nil, nil, nil, "test")
	failRelay := func() panelMessageRelay453 {
		return func(ctx context.Context, rendered string, rows []panel.AttachmentRef) error {
			return errRelayDown167
		}
	}
	okRelay := func() panelMessageRelay453 {
		return func(ctx context.Context, rendered string, rows []panel.AttachmentRef) error {
			return nil
		}
	}

	door.bindRelay453(failRelay())
	const original = "原文甲：这一段是用户先写的字\n第二行也在此"
	if err := door.HandleMessageRequest(t.Context(), sendReq167("167-draft-1", original)); err == nil {
		t.Fatal("the failing relay accepted the send: the failure branch was never exercised")
	}
	got, reqID, failures := door.failed.snapshot()
	if got != original {
		t.Errorf("失败原文没有逐字存住：slot=%q want=%q（截断、改写或合并都不算存住）", got, original)
	}
	if reqID != "167-draft-1" {
		t.Errorf("failed under request %q, want 167-draft-1", reqID)
	}
	if failures != 1 {
		t.Errorf("failures = %d, want 1", failures)
	}
	if r := door.failedSendReading(); !r.Known || r.Text != original {
		t.Errorf("reading = %+v, want Known=true with the original verbatim", r)
	}

	// The positive control: a send that does NOT fail leaves the stored
	// original alone. Wiping it here would be the "draft eaten on success" lie.
	door.bindRelay453(okRelay())
	if err := door.HandleMessageRequest(t.Context(), sendReq167("167-draft-2", "后来新写的乙")); err != nil {
		t.Fatalf("the succeeding relay refused: %v", err)
	}
	if got, _, _ := door.failed.snapshot(); got != original {
		t.Errorf("成功发送抹掉了草稿：slot=%q want=%q（不失败时草稿不被抹）", got, original)
	}

	// The newest failure replaces the slot with ITS OWN text: a merge of the
	// old draft into the new send's words is the overwrite family this cell
	// forbids, in the only direction the Go side can produce it.
	door.bindRelay453(failRelay())
	const second = "失败丙"
	if err := door.HandleMessageRequest(t.Context(), sendReq167("167-draft-3", second)); err == nil {
		t.Fatal("the failing relay accepted the second send")
	}
	if got, reqID, failures := door.failed.snapshot(); got != second || reqID != "167-draft-3" || failures != 2 {
		t.Errorf("slot after the second failure = %q (req %q, failures %d), want exactly %q from 167-draft-3 (no merge with the older draft)", got, reqID, failures, second)
	}

	// A blank text books nothing: an empty slot must not dress itself up as a
	// recovery offer, and a blank send must not wipe the slot either.
	door.bindRelay453(failRelay())
	if err := door.HandleMessageRequest(t.Context(), sendReq167("167-draft-4", "   ")); err == nil {
		t.Fatal("the blank send was accepted")
	}
	if got, _, failures := door.failed.snapshot(); got != second || failures != 2 {
		t.Errorf("blank-text failure moved the slot: text=%q failures=%d, want %q / 2", got, failures, second)
	}
}

// Test167AC5FailureTouchesNoOtherComposerField is the "no wholesale refill"
// pin: across a failed send (and a succeeding one), the composer section's
// ENTIRE key set outside the two failed-send keys is byte-identical. A mutation
// that backfills the failed text into any other field - the shape that lets a
// page clobber the user's newer typing - moves some key here and goes red.
func Test167AC5FailureTouchesNoOtherComposerField(t *testing.T) {
	door := newPanelMessageDoor453(nil, nil, nil, "test")
	pump := panel.NewSnapshotPump(panel.PumpSources{FailedSend: door.failedSendReading})
	composerJSON := func() []byte {
		data, err := json.Marshal(pump.Snapshot().Composer)
		if err != nil {
			t.Fatalf("marshal composer: %v", err)
		}
		return data
	}
	composerKeys := func(data []byte) map[string]json.RawMessage {
		m := map[string]json.RawMessage{}
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("unmarshal composer: %v (%s)", err, data)
		}
		return m
	}
	before := composerKeys(composerJSON())

	door.bindRelay453(func(ctx context.Context, rendered string, rows []panel.AttachmentRef) error {
		return errRelayDown167
	})
	if err := door.HandleMessageRequest(t.Context(), sendReq167("167-iso-1", "失败原文")); err == nil {
		t.Fatal("the failing relay accepted the send")
	}
	after := composerKeys(composerJSON())
	for key := range after {
		if key == "failedSendText" || key == "failedSendKnown" {
			continue
		}
		beforeVal, wasThere := before[key]
		if !wasThere {
			t.Errorf("failed send ADDed a composer key %q - the failure path grew a field (wholesale-refill shape)", key)
			continue
		}
		if string(beforeVal) != string(after[key]) {
			t.Errorf("failed send CHANGED composer key %q: %s -> %s（失败路径只许写自己的槽位）", key, beforeVal, after[key])
		}
	}
	for key := range before {
		if _, ok := after[key]; !ok {
			t.Errorf("failed send DROPped composer key %q", key)
		}
	}
	if _, ok := after["failedSendText"]; !ok {
		t.Error("the failure path wrote no failedSendText key: the slot did not reach the wire")
	}

	// Same pin across a SUCCESS: the success path touches nothing either.
	door.bindRelay453(func(ctx context.Context, rendered string, rows []panel.AttachmentRef) error {
		return nil
	})
	if err := door.HandleMessageRequest(t.Context(), sendReq167("167-iso-2", "新写的")); err != nil {
		t.Fatalf("the succeeding relay refused: %v", err)
	}
	afterOK := composerKeys(composerJSON())
	for key := range afterOK {
		if key == "failedSendText" || key == "failedSendKnown" {
			continue
		}
		if string(before[key]) != string(afterOK[key]) {
			t.Errorf("succeeding send CHANGED composer key %q: %s -> %s", key, before[key], afterOK[key])
		}
	}
}

// Test167AC5WireShapeIsReaderGated: no reader, no keys; a live reader always
// sends the flag; the text key appears exactly when an original is on file.
func Test167AC5WireShapeIsReaderGated(t *testing.T) {
	composerKeysOf := func(src panel.PumpSources) map[string]json.RawMessage {
		data, err := json.Marshal(panel.NewSnapshotPump(src).Snapshot().Composer)
		if err != nil {
			t.Fatalf("marshal composer: %v", err)
		}
		m := map[string]json.RawMessage{}
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("unmarshal composer: %v", err)
		}
		return m
	}

	none := composerKeysOf(panel.PumpSources{})
	if _, ok := none["failedSendText"]; ok {
		t.Error("no reader but the wire carries failedSendText - the reader-gated absence rule is gone")
	}
	if _, ok := none["failedSendKnown"]; ok {
		t.Error("no reader but the wire carries failedSendKnown")
	}

	face := composerKeysOf(panel.PumpSources{FailedSend: func() panel.FailedSendReading {
		return panel.FailedSendReading{Text: "", Known: true}
	}})
	if v, ok := face["failedSendKnown"]; !ok || string(v) != "true" {
		t.Errorf("send face present but failedSendKnown = %s (present=%v), want the literal true", v, ok)
	}
	if _, ok := face["failedSendText"]; ok {
		t.Errorf("no failure on file but the wire carries failedSendText: an empty recovery offer is the lying-key family")
	}

	full := composerKeysOf(panel.PumpSources{FailedSend: func() panel.FailedSendReading {
		return panel.FailedSendReading{Text: "原文", Known: true}
	}})
	if v, ok := full["failedSendText"]; !ok || string(v) != `"原文"` {
		t.Errorf("failedSendText = %s (present=%v), want the verbatim original", v, ok)
	}
}

// Test167AC5FailedSendFacesAreWiredAtBothPumps is the wiring nail: the slot is
// only a real exit when BOTH pump assemblies read it - the pipeline pump through
// the runtime (the resident boot binds the door's ledger onto it) and the
// degraded no-pipeline pump straight off the chain's message door. Unplugging
// either bind reddens the matching line here.
func Test167AC5FailedSendFacesAreWiredAtBothPumps(t *testing.T) {
	runSrc, err := os.ReadFile(filepath.Join("run.go"))
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	if !strings.Contains(string(runSrc), "FailedSend: rt.failedSendReading") {
		t.Error("run.go's pump literal no longer reads the failed-send slot: the pipeline shape lost the exit")
	}
	residentSrc, err := os.ReadFile(filepath.Join("resident_windows.go"))
	if err != nil {
		t.Fatalf("read resident_windows.go: %v", err)
	}
	if !strings.Contains(string(residentSrc), "bindFailedSendLedger453(rp, src)") {
		t.Error("resident_windows.go no longer binds the door's ledger onto the runtime: the resident pipeline shape lost the exit")
	}
	pumpSrc, err := os.ReadFile(filepath.Join("panel_resident_pump.go"))
	if err != nil {
		t.Fatalf("read panel_resident_pump.go: %v", err)
	}
	if !strings.Contains(string(pumpSrc), "disp.Message.(*panelMessageDoor453)") {
		t.Error("residentPumpSources65 no longer reads the message door: the degraded shape lost the exit")
	}
}
