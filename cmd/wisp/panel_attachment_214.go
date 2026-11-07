package main

// Ticket 214: the attachment ingest door, wired.
//
// WHAT EXISTED AND WHAT WAS MISSING. attachments.go (ticket 92) shipped the
// whole machine - AttachmentBroker.Ingest with its name guard, size ceiling,
// content sniffing and dedup - and, measured by ticket 145's census, had ZERO
// production callers: the pump handed NewComposerState nil, and the
// panel.attachment.add door was refused by name at the inbound router. The two
// composer keys (attachments / attachmentError) were therefore always empty by
// design, pinned by 145's test. This file is the missing caller: the handler
// behind ComposerDispatch.Attachment, the draft ledger the handler records
// into, and the readers the snapshot pump consumes.
//
// THE SHAPE IS TICKET 186'S HOLDER PATTERN. The ingest handler is built when
// the panel chain assembles (cmd/wisp/panel_inbound.go's
// newComposerDispatchChain), which happens BEFORE the task pipeline - and the
// snapshot pump that must display the drafts is built by that later pipeline
// (run.go's assembleRuntime). A ledger created by the caller, filled by the
// handler and bound onto the runtime afterwards (resident_windows.go, next to
// setPanelSnapshotHost) crosses that ordering the same way the workspace
// scope holder does, with the same nil tolerances: a chain assembled without a
// ledger still ingests and audits (the CLI seam passes nil - it has no
// snapshot to show rows in), and a runtime without a bound ledger reports the
// honest empty (the pump's reader answers nil / "" and the packet keeps the
// 145-pinned empty shape).
//
// THE SINK IS THE REAL ONE. The chain opens the process's memory.Store on the
// data root and hands it to the broker, so an accepted attachment lands as a
// content-hashed file inside the SAME artifacts tree every tool artifact
// lives in - the name guard, the O_EXCL write, the private DACL and the
// recursive quota of tickets 76/79 all apply to it, and no second home for
// user bytes is created anywhere. A store that will not open is an assembly
// failure, the same posture the settings leg takes for its secret store.
//
// THE CEILING IS THE COMPILED DEFAULT, ON PURPOSE. Ticket 214's 09-29 ruling:
// the wiring uses the existing default (max = 0 -> the broker applies
// MaxAttachmentBytes) and adds NO config key - a user-tunable limit belongs to
// the D36 config tree and is explicitly left to a later ticket.
//
// WHAT THIS FILE DOES NOT DO: no page-side picker, no @reference parser and no
// message-send door exist yet (frontend is another leg's land; message is
// ticket 35's socket and stays refused by name). Ticket 214's AC#6 names all
// three as still absent, and nothing here pretends otherwise.

import (
	"context"
	"fmt"
	"sync"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// attachmentDraftMaxRows bounds the ledger's memory the way every other
// producer in this repository bounds its own: rows are references (a few
// hundred bytes each), never bytes, but a bound that names itself beats an
// unbounded slice behind a long-lived panel. The oldest row gives way; the
// current refusal sentence never does.
const attachmentDraftMaxRows = 64

// attachmentDraftLedger is the one draft-attachment state a process holds:
// every ingest attempt the panel made this session, in arrival order, plus the
// current refusal sentence. The handler writes it; the snapshot pump reads it
// through the two readers below; nothing else touches it.
type attachmentDraftLedger struct {
	mu    sync.Mutex
	items []panel.AttachmentRef
	// last is the MOST RECENT attempt's refusal sentence, cleared by a later
	// accepted ingest. The clearing is the ruling written on composer.go's
	// AttachmentReason field: an old refusal hanging over a fresh success is
	// the lying-key family, and history is not lost by it because every
	// refusal also lives on its own row.
	last string
	// store is the artifacts store the handler opened on first ingest, held
	// here so the holder of the ledger can close it (closeStore). nil until
	// the first ingest actually opened one.
	store *memory.Store
}

func newAttachmentDraftLedger() *attachmentDraftLedger {
	return &attachmentDraftLedger{items: []panel.AttachmentRef{}}
}

// record books one ingest attempt. nil-tolerant by contract: the CLI seam
// builds its chain without a ledger (no snapshot reads it there), and a nil
// ledger makes recording a no-op, never a panic and never a lost ingest.
func (l *attachmentDraftLedger) record(ref panel.AttachmentRef) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items, ref)
	if len(l.items) > attachmentDraftMaxRows {
		l.items = l.items[len(l.items)-attachmentDraftMaxRows:]
	}
	switch {
	case ref.Stored:
		l.last = ""
	case ref.Reason != "":
		l.last = ref.Reason
	}
}

// rows returns a copy of the ledger in arrival order. nil answers nil, which
// the pump renders as the empty list - never null on the wire.
func (l *attachmentDraftLedger) rows() []panel.AttachmentRef {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]panel.AttachmentRef, len(l.items))
	copy(out, l.items)
	return out
}

// reason returns the current refusal sentence. nil answers "".
func (l *attachmentDraftLedger) reason() string {
	if l == nil {
		return ""
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.last
}

// holdStore books the artifacts store the handler opened, so the ledger's
// holder can close it (a test reclaiming its temp root). Production never
// calls closeStore: the handle lives as long as the process, the way every
// other store on the data root does.
func (l *attachmentDraftLedger) holdStore(s *memory.Store) {
	if l == nil || s == nil {
		return
	}
	l.mu.Lock()
	l.store = s
	l.mu.Unlock()
}

// closeStore closes the held store, once. nil-safe everywhere.
func (l *attachmentDraftLedger) closeStore() error {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	s := l.store
	l.store = nil
	l.mu.Unlock()
	if s == nil {
		return nil
	}
	return s.Close()
}

// newAttachmentIngestHandler214 builds the production AttachmentRequestHandler
// over the data root named by dir. The store is opened LAZILY, on the first
// ingest, and stays open for the process: a chain that never answers an
// attachment (the mode/settings-only tests, a pipe that carries no files)
// opens no SQLite database, and a store that will not open refuses THE INGEST
// - loudly, audited, with a ledger row - instead of failing the whole chain's
// assembly for a door the request never used. The ceiling is the compiled
// default (max = 0 -> the broker applies MaxAttachmentBytes), ticket 214's
// 09-29 ruling: no config key in this wiring.
func newAttachmentIngestHandler214(dir string, drafts *attachmentDraftLedger, audit panel.AuditFunc, actor string) *attachmentIngestHandler214 {
	return &attachmentIngestHandler214{dir: dir, drafts: drafts, audit: audit, actor: actor}
}

// attachmentIngestHandler214 answers panel.attachment.add: decode the wire
// payload, hand it to the broker, book the outcome in the draft ledger, audit
// the line. The error it returns is the router's refusal - it already wraps
// panel.ErrAttachmentRejected, so the caller can classify without prose
// matching - and an accepted request returns nil exactly like the other
// answered doors (the router adds no second voice for this route).
type attachmentIngestHandler214 struct {
	// dir is the data root the artifacts store opens over, on first ingest.
	dir string
	// drafts may be nil (the CLI seam): the ingest still happens and the audit
	// line still lands; only the snapshot row is absent, which is honest in a
	// process with no snapshot pump.
	drafts *attachmentDraftLedger
	audit  panel.AuditFunc
	actor  string

	once    sync.Once
	broker  *panel.AttachmentBroker
	openErr error
}

// brokerOnce opens the store and builds the broker exactly once, on the first
// ingest that reaches the door. The guard is the store's own (memory.
// ValidArtifactName - the same function the store applies at write time,
// ticket 76's single-owner rule), so an accepted attachment lands as a
// content-hashed file inside the SAME guarded artifacts tree every tool
// artifact lives in (tickets 76/79), never in a second home.
func (h *attachmentIngestHandler214) brokerOnce() (*panel.AttachmentBroker, error) {
	h.once.Do(func() {
		store, err := memory.Open(h.dir)
		if err != nil {
			h.openErr = fmt.Errorf("附件存储不可用（%s）：%w", h.dir, err)
			return
		}
		h.broker, h.openErr = panel.NewAttachmentBroker(store, memory.ValidArtifactName, 0)
		if h.broker != nil && h.drafts != nil {
			// The ledger the caller holds takes the store's lifetime, so a test
			// (and only a test - production keeps the handle for the process's
			// life, the way every other open store on this root does) can close
			// it before its temp root is reclaimed.
			h.drafts.holdStore(store)
		}
	})
	return h.broker, h.openErr
}

func (h *attachmentIngestHandler214) HandleAttachmentRequest(ctx context.Context, req panel.ComposerRequest) error {
	ref, err := h.ingest(ctx, req)
	h.drafts.record(ref)
	if h.audit != nil {
		h.audit("panel: ATTACHMENT-ADD request=%q actor=%q name=%q stored=%v deduplicated=%v bytes=%d artifact=%q err=%v",
			req.RequestID, h.actor, req.AttachmentPayload.Name, ref.Stored, ref.Deduplicated,
			ref.SizeBytes, ref.Artifact, err)
	}
	return err
}

// ingest is the decode-then-store half, split so the ledger always gets a row:
// a payload that cannot even be decoded (a lied-about size, an undecodable
// base64) is still an attempt the user made, and AC#2's rule is that a refusal
// is visible with its reason, not swallowed between the decode and the store.
func (h *attachmentIngestHandler214) ingest(ctx context.Context, req panel.ComposerRequest) (panel.AttachmentRef, error) {
	src, err := panel.DecodeAttachmentPayload(req.AttachmentPayload)
	if err != nil {
		// The decode failed, so no broker verdict exists: the row carries the
		// payload's own display name and the decode refusal, and the request id
		// stands in for the broker's content-hash id the row will never have.
		return panel.AttachmentRef{ID: req.RequestID, Name: req.AttachmentPayload.Name, Reason: err.Error()}, err
	}
	broker, err := h.brokerOnce()
	if err != nil {
		return panel.AttachmentRef{ID: req.RequestID, Name: src.DisplayName, Reason: err.Error()}, err
	}
	return broker.Ingest(ctx, src)
}

// bindAttachmentDrafts is the resident boot's one write of the runtime's draft
// source, the sibling of setPanelSnapshotHost: the pump is built inside
// assembleRuntime, the ledger exists before it, and this bind crosses the
// ordering the same way the workspace holder does. Called after assembly, on
// the same boot line, before any task can run.
func (rt *agentRuntime) bindAttachmentDrafts(l *attachmentDraftLedger) {
	if rt == nil {
		return
	}
	rt.attMu.Lock()
	rt.attachments = l
	rt.attMu.Unlock()
}

// attachmentDraftRows is the pump's reader for PumpSources.Attachments. nil
// runtime or no bound ledger answers nil - the honest no-source state the
// packet renders as the empty list, exactly what ticket 145 pinned for the
// unwired tree and what a console `wisp run` still is today.
func (rt *agentRuntime) attachmentDraftRows() []panel.AttachmentRef {
	if rt == nil {
		return nil
	}
	rt.attMu.Lock()
	l := rt.attachments
	rt.attMu.Unlock()
	return l.rows()
}

// attachmentDraftReason is the pump's reader for PumpSources.AttachmentReason.
// nil answers "" - never an invented sentence.
func (rt *agentRuntime) attachmentDraftReason() string {
	if rt == nil {
		return ""
	}
	rt.attMu.Lock()
	l := rt.attachments
	rt.attMu.Unlock()
	return l.reason()
}
