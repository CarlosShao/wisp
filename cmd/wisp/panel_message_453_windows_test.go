//go:build windows

package main

// Ticket 453 AC#1 and AC#2's end-to-end teeth.
//
// The BEFORE half of this file was written and run before the message door's
// relay existed, so it is the measurement that the exit did not exist. At that
// HEAD the production chain (newResidentPanelManager) really stored the PNG -
// the audit line says stored=true artifact="attachment-....png" - and
// panel.message.send then refused with the 192 preview door's sentence "...
// 本机未接入". The exact red is quoted in docs/evidence/s1/453-ac1-no-exit.md.
//
// The file now carries four cases, each aimed at a place AC#2/AC#3 needs a
// behaviour (not a source string) to be load bearing:
//
//	Test453PanelMessageCarries...   ingest -> send -> OutgoingMessage.ForAgent
//	                                -> real agent loop -> the C5 request text
//	Test453WithoutARelay...         the nil-relay branch refuses by name and
//	                                runs no model (delete the nil check -> panic)
//	Test453AForgedRow...            the ledger, not the page, is the source of truth
//	Test453PanelMessageMarks...     admitTask stamps the loop-minted task id in C25
//	Test453BindingTheRelay...       bindPanelMessageRelay453 walks the chain and
//	                                attaches the source's submit method

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// png453 is the smallest bytes the ingest door's sniffer accepts as a PNG: the
// eight-byte signature is the whole decision (internal/panel/attachments.go:120),
// and a real image is not needed to prove which file arrived.
var png453 = []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 0}

// attachmentEnvelope453 is the exact wire shape frontend/src/lib/panel.ts:709
// sends (panel.attachment.add with name/declaredMime/sizeBytes/dataBase64).
func attachmentEnvelope453(t *testing.T, rid, name string) string {
	t.Helper()
	raw, err := json.Marshal(map[string]any{
		"method":       panel.MethodAttachmentAdd,
		"requestId":    rid,
		"source":       panel.ComposerRequestSource,
		"name":         name,
		"declaredMime": "image/png",
		"sizeBytes":    len(png453),
		"dataBase64":   base64.StdEncoding.EncodeToString(png453),
	})
	if err != nil {
		t.Fatalf("marshal attachment envelope: %v", err)
	}
	return string(raw)
}

// messageEnvelope453 is the exact wire shape frontend/src/lib/panel.ts:736 sends:
// text plus the stored rows the page read out of the snapshot.
func messageEnvelope453(t *testing.T, rid, text string, rows []panel.AttachmentRef) string {
	t.Helper()
	raw, err := json.Marshal(map[string]any{
		"method":      panel.MethodMessageSend,
		"requestId":   rid,
		"source":      panel.ComposerRequestSource,
		"text":        text,
		"attachments": rows,
	})
	if err != nil {
		t.Fatalf("marshal message envelope: %v", err)
	}
	return string(raw)
}

// residentChain453 assembles the same chain the resident boot assembles and
// ingests one PNG through the production attachment door.
func residentChain453(t *testing.T) (*PanelManager, *attachmentDraftLedger, panel.AttachmentRef, string) {
	t.Helper()
	dir := newInboundDataDir33(t, "ask_every_step")
	host, _, _, drafts, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("the resident panel chain would not assemble: %v", err)
	}
	if host == nil || host.disp == nil {
		t.Fatal("newResidentPanelManager returned no dispatch chain")
	}
	if drafts == nil {
		t.Fatal("newResidentPanelManager returned no attachment ledger")
	}
	t.Cleanup(func() { _ = drafts.closeStore() })

	ctx := context.Background()
	if _, err := host.dispatchRaw(ctx, attachmentEnvelope453(t, "r-453-ingest", "shot-453.png")); err != nil {
		t.Fatalf("the ingest door refused a valid PNG: %v", err)
	}
	var row panel.AttachmentRef
	for _, got := range drafts.rows() {
		if got.Stored {
			row = got
			break
		}
	}
	if row.Artifact == "" {
		t.Fatalf("the ingest door reported no stored row: rows=%+v", drafts.rows())
	}
	return host, drafts, row, dir
}

func messageDoor453(t *testing.T, host *PanelManager) *panelMessageDoor453 {
	t.Helper()
	door, ok := host.disp.Message.(*panelMessageDoor453)
	if !ok {
		t.Fatalf("the resident chain's message socket is %T, want *panelMessageDoor453", host.disp.Message)
	}
	return door
}

// recordingProvider453 is the smallest C5 seam that can answer, and it records
// the text of every request it is handed BEFORE emitting anything: the model
// input is the assertion surface, not the loop's result.
type recordingProvider453 struct {
	mu   sync.Mutex
	seen []string
}

func (p *recordingProvider453) Info() llm.ProviderInfo {
	return llm.ProviderInfo{
		Provider:         "probe-453",
		Model:            "probe-453",
		Protocol:         "probe",
		MaxContextWindow: 128000,
	}
}

func (p *recordingProvider453) Stream(_ context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	var b strings.Builder
	for _, msg := range req.Messages {
		for _, c := range msg.Content {
			if txt, ok := c.(llm.TextPart); ok {
				b.WriteString(txt.Text)
				b.WriteString("\n")
			}
		}
	}
	p.mu.Lock()
	p.seen = append(p.seen, b.String())
	p.mu.Unlock()
	for _, ev := range []llm.StreamEvent{
		{Type: llm.EvTextDelta, Text: "看到了这张图"},
		{Type: llm.EvStop, Stop: llm.StopEndTurn},
		{Type: llm.EvDone},
	} {
		if err := emit(ev); err != nil {
			return err
		}
	}
	return nil
}

func (p *recordingProvider453) requests() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]string, len(p.seen))
	copy(out, p.seen)
	return out
}

// Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn is AC#1's
// measurement and AC#2's end-to-end tooth.
func Test453PanelMessageCarriesTheAcceptedAttachmentIntoTheModelsNextTurn(t *testing.T) {
	host, drafts, row, _ := residentChain453(t)
	door := messageDoor453(t, host)

	prov := &recordingProvider453{}
	loop, err := agent.New(agent.Options{
		Provider: prov,
		Tools:    agent.NewEchoProvider(),
		Sink:     &agent.RecordSink{},
		Journal:  nil,
		Registry: observe.NewRegistry(),
		Logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
		Config: agent.Config{
			Model:                       "probe-453",
			ContextWindow:               128000,
			MaxRounds:                   2,
			ArtifactsDir:                drafts.artifactsDir453(),
			PerToolTimeout:              2 * time.Second,
			PassThroughUnclassifiedRisk: true,
			Currency:                    "CNY",
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	door.bindRelay453(func(ctx context.Context, rendered string, marks []panel.AttachmentRef) error {
		res := loop.Run(ctx, rendered)
		if res.Err != nil {
			return fmt.Errorf("agent loop failed: %w", res.Err)
		}
		return nil
	})

	if _, err := host.dispatchRaw(context.Background(), messageEnvelope453(t, "r-453-send", "看看这张图", []panel.AttachmentRef{row})); err != nil {
		t.Fatalf("panel.message.send has no exit: the accepted attachment %q (artifact %q) cannot reach the model, so the next turn's input never sees it: %v",
			row.Name, row.Artifact, err)
	}

	seen := prov.requests()
	if len(seen) == 0 {
		t.Fatal("the message was accepted but the C5 seam was never called: no model turn exists")
	}
	path := filepath.Join(drafts.artifactsDir453(), row.Artifact)
	all := strings.Join(seen, "\n---\n")
	for _, want := range []string{"[附件]", "name=" + row.Name, "path=" + path, "看看这张图"} {
		if !strings.Contains(all, want) {
			t.Fatalf("the model's next turn input does not contain %q.\nThe input was:\n%s", want, all)
		}
	}
}

// Test453WithoutARelayTheMessageIsRefusedByNameAndRunsNoModel pins the pre-453
// shape: a chain whose relay was never bound must say so, not drop the message.
// Deleting the nil-relay branch does not quietly deliver - it panics on the
// call - which is why this case is the behaviour tooth for that branch.
func Test453WithoutARelayTheMessageIsRefusedByNameAndRunsNoModel(t *testing.T) {
	host, _, row, _ := residentChain453(t)
	_, err := host.dispatchRaw(context.Background(), messageEnvelope453(t, "r-453-unbound", "看看这张图", []panel.AttachmentRef{row}))
	if !errors.Is(err, panel.ErrNoHandlerAttached) {
		t.Fatalf("unbound relay error = %v, want an ErrNoHandlerAttached refusal", err)
	}
	if !strings.Contains(err.Error(), "消息中继未接入") {
		t.Fatalf("unbound relay refusal = %q, want the named '消息中继未接入' sentence", err.Error())
	}
}

// Test453AForgedRowIsRefusedBeforeAnyRendering pins the selection's source of
// truth: the page can echo a ledger row, it cannot mint one.
func Test453AForgedRowIsRefusedBeforeAnyRendering(t *testing.T) {
	host, _, _, _ := residentChain453(t)
	forged := panel.AttachmentRef{
		ID: "att-forged", Name: "not-there.png", MIME: "image/png",
		Kind: "image", SizeBytes: 12, Artifact: "attachment-deadbeef.png", Stored: true,
	}
	_, err := host.dispatchRaw(context.Background(), messageEnvelope453(t, "r-453-forged", "看看", []panel.AttachmentRef{forged}))
	if !errors.Is(err, panel.ErrComposerRequest) {
		t.Fatalf("forged row error = %v, want an ErrComposerRequest refusal", err)
	}
	if !strings.Contains(err.Error(), "不在本进程已收下的草稿里") {
		t.Fatalf("forged row refusal = %q, want the ledger's own sentence", err.Error())
	}
}

// Test453PanelMessageMarksTheAttachmentsUnderTheTaskScope is AC#3's behaviour:
// the mark is made under the id admitTask receives, and the egress check reports
// it as panel.attachment. Remove the markPanelMessage453 call from admitTask and
// this reddens; mark under another source and the SrcTool assertion reddens.
func Test453PanelMessageMarksTheAttachmentsUnderTheTaskScope(t *testing.T) {
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	rt := &agentRuntime{
		c25:    prov,
		bridge: tools.New(tools.Options{Provenance: prov}),
		gate:   approval.New(approval.Options{}),
	}
	const taskID = "task-453-c25"
	rt.stagePanelMessage453([]panel.AttachmentRef{{
		ID: "att-1", Name: "shot-453.png", MIME: "image/png", Kind: "image",
		SizeBytes: 12, Artifact: "attachment-1b56b50ac4e976f4.png", Stored: true,
	}})
	revoke := rt.admitTask(taskID)
	defer revoke()

	hit, ok := prov.Inspect(taskID, "notify", map[string]any{"text": "转到别处：shot-453.png"})
	if !ok {
		t.Fatalf("no taint after a panel attachment reached task %q: the attachment's name can leave unmarked", taskID)
	}
	if hit.SrcTool != risk.SrcPanelAttachment {
		t.Fatalf("taint source = %q, want %q", hit.SrcTool, risk.SrcPanelAttachment)
	}
	if rt.takePanelMessage453() != nil {
		t.Fatal("the staged rows survived admission: a second task would inherit this message's marks")
	}
}

// Test453BindingTheRelayWalksTheResidentChain pins the boot's bind function: it
// finds the composite socket on the real chain and attaches the source's submit
// method. A source with no runtime is the smallest object that proves the attach
// happened - the refusal it returns is SubmitPanelMessage453's own sentence, not
// the door's unbound one.
func Test453BindingTheRelayWalksTheResidentChain(t *testing.T) {
	host, _, row, _ := residentChain453(t)
	if messageDoor453(t, host).relay453() != nil {
		t.Fatal("the chain arrived with a relay already bound: the boot has nothing left to bind")
	}
	bindPanelMessageRelay453(host, &residentTaskSource{})
	if messageDoor453(t, host).relay453() == nil {
		t.Fatal("bindPanelMessageRelay453 left the door unbound")
	}
	_, err := host.dispatchRaw(context.Background(), messageEnvelope453(t, "r-453-bind", "看看", []panel.AttachmentRef{row}))
	if err == nil || !strings.Contains(err.Error(), "没有任务管线") {
		t.Fatalf("a bound relay over a source with no runtime = %v, want SubmitPanelMessage453's '没有任务管线' refusal", err)
	}
}
