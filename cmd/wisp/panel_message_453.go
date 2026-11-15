package main

// Ticket 453: the attachment's EXIT.
//
// WHAT WAS MISSING (ticket 214 AC#5, measured at the 453 dispatch anchor): the
// ingest leg really stored the bytes - the artifacts tree, the name guard,
// O_EXCL, the private DACL, the recursive quota, content-hash naming - but the
// half that carries a stored file into a conversation did not exist:
//
//	OutgoingMessage{...}   production constructions: 0
//	ForAgent(...)          production callers:      0
//	panel.message.send     handler:                 the 192 preview door only,
//	                                                which refuses by name
//
// So the user could attach a file, see it in the composer, press send, and the
// model never received it. This file is that half, in the form the 214 evidence
// (docs/evidence/s1/214-ac5-stamp-feasibility.md §3-§4) pre-decided: the handler
// body lands here, the rendering is the existing OutgoingMessage.ForAgent, and
// the C25 stamp happens under the task id the text loop itself mints.
//
// THE PATH, one sentence: panel.message.send -> this door -> the ledger's own
// stored rows (never the page's claims) -> OutgoingMessage.ForAgent -> the
// resident task source (the one task slot) -> agentRuntime.admitTask, where the
// loop-minted task id first exists -> risk.Mark(taskID, SrcPanelAttachment, ...)
// -> agent.Loop.Run(rendered), whose first C5 request now carries
// "[附件] name=... path=...".
//
// WHAT THIS DOOR IS NOT: a second relay for preview quotes. A request whose path
// is a preview address is forwarded to the 192 door UNCHANGED (same stamp, same
// refusals), so there is still exactly one QuoteForModel caller in the tree. The
// plain-message branch is the new one, and it is only reachable when the boot
// bound a relay; without one it refuses by name, exactly as before.

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

// panelMessageRelay453 is the one seam this door needs from the task pipeline:
// hand the composed text (already rendered with ForAgent) and the rows it was
// rendered from to whatever owns the task slot. The resident source implements
// it; a boot with no pipeline binds nothing and the door refuses by name.
type panelMessageRelay453 func(ctx context.Context, rendered string, marks []panel.AttachmentRef) error

// panelMessageDoor453 is the message socket of the resident chain. It wraps the
// 192 preview door rather than replacing it: preview requests still go through
// the ONLY QuoteForModel caller, and the composite only adds the branch the
// ticket exists for.
type panelMessageDoor453 struct {
	preview *panelPreviewQuoteDoor
	drafts  *attachmentDraftLedger
	audit   panel.AuditFunc
	actor   string

	mu    sync.RWMutex
	relay panelMessageRelay453
}

func newPanelMessageDoor453(preview *panelPreviewQuoteDoor, drafts *attachmentDraftLedger, audit panel.AuditFunc, actor string) *panelMessageDoor453 {
	return &panelMessageDoor453{preview: preview, drafts: drafts, audit: audit, actor: actor}
}

// bindRelay453 is the late bind. The panel chain is assembled before the task
// pipeline exists (resident_windows.go:359 vs :548), so the relay cannot be a
// constructor argument: the boot binds it on the first line where the source
// object exists. A nil relay is a legal state (CLI seam, refused pipeline) and
// means "refuse by name", never "drop".
func (d *panelMessageDoor453) bindRelay453(relay panelMessageRelay453) {
	if d == nil {
		return
	}
	d.mu.Lock()
	d.relay = relay
	d.mu.Unlock()
}

func (d *panelMessageDoor453) relay453() panelMessageRelay453 {
	if d == nil {
		return nil
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.relay
}

// HandleMessageRequest is the door. Outcomes, in the order they are decided:
//
//	preview address     -> the 192 door, verbatim (one QuoteForModel caller stays one)
//	stale/forged rows   -> refused before anything is rendered
//	no relay bound      -> refused by name (the pre-453 behaviour, still honest)
//	nothing to send     -> refused before the task slot is taken
//	otherwise           -> OutgoingMessage.ForAgent -> relay -> nil (accepted)
//
// The rows are selected from the LEDGER by ID, never taken from the envelope:
// the page can echo a row, it cannot mint one. A row it names that the ledger
// does not hold as Stored is a refusal, not a best effort.
func (d *panelMessageDoor453) HandleMessageRequest(ctx context.Context, req panel.ComposerRequest) error {
	if d == nil {
		return fmt.Errorf("%w: 消息门不存在", panel.ErrNoHandlerAttached)
	}
	if d.preview != nil && strings.HasPrefix(previewPathPart192(req.Path), panel.PreviewPathPrefix) {
		return d.preview.HandleMessageRequest(ctx, req)
	}
	rows, selErr := d.drafts.storedSelection453(req.Attachments)
	if selErr != nil {
		err := fmt.Errorf("面板消息没有送出：%w", selErr)
		d.record453(req, 0, err)
		return err
	}
	relay := d.relay453()
	if relay == nil {
		err := fmt.Errorf("%w: 消息中继未接入（本进程没有任务管线，或任务入口没有起来）：已收下的附件 %d 枚留在本机，但它们没有进模型，也不会被静默丢掉",
			panel.ErrNoHandlerAttached, len(rows))
		d.record453(req, len(rows), err)
		return err
	}
	if strings.TrimSpace(req.Text) == "" && len(rows) == 0 {
		err := fmt.Errorf("%w: 这条消息既没有文字，也没有已收下的附件，本机没有内容可送出", panel.ErrComposerRequest)
		d.record453(req, 0, err)
		return err
	}
	// THE production construction point of panel.OutgoingMessage. It carries
	// the ledger truth, and the rendering below is ForAgent's contract - the
	// same sentences the agent's own shape test pins
	// (panel_attachment_214_test.go:253).
	msg := panel.OutgoingMessage{Text: req.Text, Attachments: rows}
	rendered := msg.ForAgent(d.drafts.artifactsDir453())
	if strings.TrimSpace(rendered) == "" {
		err := fmt.Errorf("%w: 这条消息渲染后是空的（文字与附件都没有可送出的内容）", panel.ErrComposerRequest)
		d.record453(req, len(rows), err)
		return err
	}
	if err := relay(ctx, rendered, rows); err != nil {
		err = fmt.Errorf("面板消息没有送进模型：%w", err)
		d.record453(req, len(rows), err)
		return err
	}
	d.record453(req, len(rows), nil)
	return nil
}

// record453 writes one audit line per attempt. It never carries the text, the
// file name or the artifact coordinate: only counts, the request id and the
// outcome, the same discipline panel_preview_192.go:307 holds.
func (d *panelMessageDoor453) record453(req panel.ComposerRequest, attachments int, err error) {
	if d == nil || d.audit == nil {
		return
	}
	outcome := "accepted"
	if err != nil {
		outcome = "refused"
	}
	d.audit("panel: MESSAGE-SEND request=%q actor=%q attachments=%d textRunes=%d outcome=%s err=%v",
		req.RequestID, d.actor, attachments, len([]rune(req.Text)), outcome, err)
}

// storedSelection453 turns the page's claimed rows into the ledger's own rows.
//
// The page read them out of the snapshot, so for an honest page this is an
// identity map. For a page that invents one, it is a refusal: the ledger's
// Stored rows are the only proof that bytes exist on disk, and the message must
// not be able to name a file the ingest leg never wrote.
func (l *attachmentDraftLedger) storedSelection453(claimed []panel.AttachmentRef) ([]panel.AttachmentRef, error) {
	if len(claimed) == 0 {
		return nil, nil
	}
	if l == nil {
		return nil, fmt.Errorf("%w: 这条消息声称带了 %d 枚附件，但本进程没有附件草稿账本可以核对", panel.ErrComposerRequest, len(claimed))
	}
	known := make(map[string]panel.AttachmentRef, len(claimed))
	for _, row := range l.rows() {
		if row.Stored && row.ID != "" {
			known[row.ID] = row
		}
	}
	out := make([]panel.AttachmentRef, 0, len(claimed))
	seen := make(map[string]bool, len(claimed))
	for _, claim := range claimed {
		row, ok := known[claim.ID]
		if !ok {
			return nil, fmt.Errorf("%w: 附件 %q（id=%q）不在本进程已收下的草稿里；页面声称的引用不能当收件凭据",
				panel.ErrComposerRequest, claim.Name, claim.ID)
		}
		if seen[row.ID] {
			return nil, fmt.Errorf("%w: 附件 %q 在一条消息里被引用了两次", panel.ErrComposerRequest, row.Name)
		}
		seen[row.ID] = true
		out = append(out, row)
	}
	return out, nil
}

// artifactsDir453 is the directory ForAgent prefixes onto each artifact name.
// A ledger with no store (every boot before the first ingest, the CLI seam) has
// no directory and no stored rows; empty is the honest answer there.
func (l *attachmentDraftLedger) artifactsDir453() string {
	if l == nil {
		return ""
	}
	l.mu.Lock()
	store := l.store
	l.mu.Unlock()
	if store == nil {
		return ""
	}
	return store.ArtifactsDir()
}

// stagingPanelMessage453 is the runtime's one-slot handoff from the door to the
// task boundary. The loop mints the task id inside RunAsync, so the mark cannot
// be made when the message is accepted; it is made in admitTask, the first hook
// that receives that id (run.go). These fields hold the rows across that gap.
func (rt *agentRuntime) stagePanelMessage453(rows []panel.AttachmentRef) {
	if rt == nil {
		return
	}
	cp := make([]panel.AttachmentRef, len(rows))
	copy(cp, rows)
	rt.msgMu.Lock()
	rt.pendingMarks453 = cp
	rt.msgMu.Unlock()
}

func (rt *agentRuntime) clearPanelMessage453() {
	if rt == nil {
		return
	}
	rt.msgMu.Lock()
	rt.pendingMarks453 = nil
	rt.msgMu.Unlock()
}

func (rt *agentRuntime) takePanelMessage453() []panel.AttachmentRef {
	if rt == nil {
		return nil
	}
	rt.msgMu.Lock()
	defer rt.msgMu.Unlock()
	rows := rt.pendingMarks453
	rt.pendingMarks453 = nil
	return rows
}

// markPanelMessage453 is the C25 stamp. It runs under the task id the loop
// minted, so the scope it writes to is the one the bridge opened for that task
// and the one CloseTask will close. Content is the attachment's context
// representation (display name + artifact coordinate), not the raw bytes: a
// binary PNG has no text to match, and the coordinate is what a message would
// have to reproduce to carry the file outward.
func (rt *agentRuntime) markPanelMessage453(taskID string) {
	rows := rt.takePanelMessage453()
	if len(rows) == 0 || taskID == "" || rt.c25 == nil {
		return
	}
	for _, row := range rows {
		if row.Artifact == "" {
			continue
		}
		content := row.Name
		if content != "" {
			content += "\n"
		}
		content += row.Artifact
		if !rt.c25.Mark(taskID, risk.SrcPanelAttachment, "artifact:"+row.Artifact, content) {
			if rt.spec.sink != nil {
				rt.auditf("wisp run: 面板附件的 C25 来源戳没有记上（task=%s artifact=%s）：脱敏引擎拒绝了这次 Mark", taskID, row.Artifact)
			}
		}
	}
}
