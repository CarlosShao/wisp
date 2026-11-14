package main

// Ticket 192, the wiring leg: the three steps that turn a landed preview stack
// into a reachable product feature. The stack itself lives in
// internal/panel/preview_192.go (leg 192-ac12-r1, evidence
// docs/evidence/s1/192-preview-ac12-r1.md) and is complete: four classes, bytes
// instead of paths, the truncation pointer, and the one C25 door. What was
// missing was every caller - measured before this leg:
//
//	PreviewResolver production call sites          0 (only probes and tests)
//	QuoteForModel production call sites            0
//	panelAssetLookup preview branch                none (:563 answered the embed only)
//
// so in the product a preview answered 404 and the quote door could not be
// reached at all. This file is the three callers, and it deliberately holds no
// preview LOGIC of its own: every rule (which classes, the caps, the sandbox
// policy, the pointer, the proof that a quote is really in the artifact) stays
// in internal/panel where it is judged. The subject here is assembly.
//
// STEP 1 of the three lives in cmd/wisp/panel_control_windows.go: the resource
// filter's lookup gained the preview branch, and the response builder takes the
// resolver's own header block (panelAssetResponse.Headers) instead of the panel
// CSP. That is the step that decides whether the HTML sandbox is real: the
// class-specific policy is built where the class is known, and a host that
// "prepends nothing of its own for a preview" (PreviewDocument.Headers' own
// doc) is exactly what makes a previewed document unable to speak. Publishing
// the panel's fallback policy over a preview instead would look like a served
// preview and be a page that can still postMessage: the AC#2 mutation the
// previous leg's M3 aims at.
//
// STEP 2 is panelPreviewSource below plus the withArtifactPreview option: the
// composition root builds ONE source per data root (cmd/wisp/panel_resident_windows.go's
// newResidentPanelManager) and hands it to the host, so the window's filter and
// the message door read one artifacts directory, one name guard and one set of
// budgets. The store opens lazily on the first preview request - the same
// posture ticket 214's ingest door takes - because a chain that never previews
// anything should not open a database.
//
// STEP 3 is panelPreviewQuoteDoor below: the message handler that names
// QuoteForModel. The page never had a way to say "this text came from
// artifact X", and minting one was not available to this leg (C17 gains no
// method, and the composer envelope is shared by every answered route, so a new
// field would be readable from all of them - the same ruling ticket 248 wrote
// down for credentials). So the door reads the association off the ONE key the
// message route already carries and nobody reads today: `path`, holding the
// SAME same-origin preview URL the page used to fetch the bytes it is now
// quoting from. That is not a second address space: it is PreviewURL(name), the
// very string the preview response was served from, so "the text the user
// selected" and "the document those bytes came out of" are one address.
//
// WHAT THE DOOR DOES NOT DO. It does not relay a message to the model: the
// message dispatch socket is ticket 35's, still unbuilt (see
// cmd/wisp/panel_inbound.go's chain, where this door replaces a nil socket and
// everything that is not a preview quote still refuses BY NAME). So the door's
// success path stamps the quote and then refuses the relay with a sentence that
// says both halves. That refusal is the honest state: what this leg buys is
// that no preview text can reach a model except through QuoteForModel - the
// stamp is in the process's own C25 engine under panelPreviewQuoteScope, and
// the moment ticket 35 lands its relay, the text it forwards is already marked
// (the half that is still missing - carrying that mark into a TASK scope - is
// named in the evidence file, not papered over here).

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// panelPreviewQuoteScope is the C25 scope id a preview quote is stamped under.
//
// Why a named process scope and not a task id: the panel's message route runs
// in a process whose task pipeline may not exist yet (the resident boot
// assembles the panel first), and a scope id invented per request would be a
// roster nobody owns. This constant is the panel surface's own scope, the shape
// the present leg's presentClipboardScope (resident_present_windows.go) and
// resident_voice's voiceAnnounceScope already use for surfaces that are not a
// task. It is NOT opened: risk.Provenance.Mark records into a named scope
// without one (internal/risk/provenance.go logs exactly that shape), so this
// file opens no scope and leaves the G5 "opened and never closed" family
// untouched.
const panelPreviewQuoteScope = "panel-preview-quote"

// panelPreviewLookup is the one question the panel window asks the preview
// stack: "what is at this same-origin path". *panel.PreviewResolver satisfies
// it as it stands, and the interface is declared here rather than in
// internal/panel so the resolver keeps knowing nothing about a WebView.
type panelPreviewLookup interface {
	Lookup(ctx context.Context, requestPath string) panel.PreviewDocument
}

// panelAssetResponse is one answer the resource filter serves:
// (bytes, content type, status, reason phrase, header block).
//
// Headers is the field this leg added, and it is what makes the preview's CSP
// real. An empty Headers means "the host's own header set applies"
// (panelResponseHeaders), which is the whole pre-192 behaviour and stays the
// answer for the embedded bundle, for a miss, and for a preview refusal - none
// of those has a class-specific policy to publish. A NON-empty Headers is a
// complete header block, assembled by the preview stack from PreviewDocument.
// Headers, and the host then adds nothing: one document, one policy, chosen
// where the class is known.
type panelAssetResponse struct {
	Data        []byte
	ContentType string
	Status      int
	Reason      string
	Headers     string
}

// panelPreviewSource is the process's ONE preview stack: the artifacts store
// opened on the data root the first time a preview is asked for, the resolver
// built over it, and the C25 marker this process's engine was bound as.
//
// It is a holder in the same sense as attachmentDraftLedger and ticket 186's
// workspace holder: it exists before the task pipeline does, it is handed to
// both surfaces that need it (the window's filter and the message door), and
// the marker lands in it afterwards from the assembly root, because the
// process's provenance engine (agentRuntime.c25) is assembled by the task
// pipeline, not by the panel boot. Until that bind, QuoteForModel refuses by
// name - "no C25 engine, so the quote cannot be stamped, so it must not enter a
// model" - which is the fail-closed state, not a stub.
type panelPreviewSource struct {
	dir   string
	audit panel.AuditFunc

	once     sync.Once
	store    *memory.Store
	resolver *panel.PreviewResolver
	openErr  error

	mu     sync.Mutex
	marker panel.TaintMarker
}

// newPanelPreviewSource192 builds the holder. No store is opened here: the
// first preview request does that, so a panel that never previews anything
// costs no SQLite open (the rule newAttachmentIngestHandler214 states for its
// own door).
func newPanelPreviewSource192(dir string, audit panel.AuditFunc) *panelPreviewSource {
	return &panelPreviewSource{dir: dir, audit: audit}
}

// resolver192 opens the store and builds the resolver exactly once. The guard
// is memory.ValidArtifactName - the store's own bare-name rule, the same
// function the ingest door injects - so a preview name that spells a path is
// refused by the artifacts store's own judgement and never by a rule spelled
// here.
func (p *panelPreviewSource) resolver192() (*panel.PreviewResolver, error) {
	if p == nil {
		return nil, fmt.Errorf("预览源未装配：没有 artifacts 目录，预览与引用都无从谈起")
	}
	p.once.Do(func() {
		store, err := memory.Open(p.dir)
		if err != nil {
			p.openErr = fmt.Errorf("预览存储不可用（%s）：%w", p.dir, err)
			return
		}
		p.store = store
		resolver, rerr := panel.NewPreviewResolver(store, memory.ValidArtifactName)
		if rerr != nil {
			p.openErr = fmt.Errorf("预览解析器不可用（%s）：%w", p.dir, rerr)
			return
		}
		p.resolver = resolver
	})
	return p.resolver, p.openErr
}

// Lookup answers one intercepted request path through the resolver, turning a
// failure to open the store into a document the user can read. It is what makes
// this type a panelPreviewLookup, i.e. what the window's resource filter calls.
func (p *panelPreviewSource) Lookup(ctx context.Context, requestPath string) panel.PreviewDocument {
	resolver, err := p.resolver192()
	if err != nil {
		return panel.PreviewDocument{Status: http.StatusInternalServerError, Reason: err.Error()}
	}
	return resolver.Lookup(ctx, requestPath)
}

// bindMarker192 hands the door the process's C25 engine. Called by the assembly
// root once the task pipeline exists (cmd/wisp/resident_windows.go), never by
// the panel boot itself, and nil is a no-op: an unbound marker keeps the door
// refusing, which is the honest answer for a process with no engine.
func (p *panelPreviewSource) bindMarker192(marker panel.TaintMarker) {
	if p == nil || marker == nil {
		return
	}
	p.mu.Lock()
	p.marker = marker
	p.mu.Unlock()
}

// marker192 reads the bound engine. nil means "this process has none yet".
func (p *panelPreviewSource) marker192() panel.TaintMarker {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.marker
}

// panelPreviewAssetResponse turns one resolver answer into the filter's own
// answer. The header block is the resolver's, joined with CRLF because that is
// the one-header-set-per-response form CreateWebResourceResponse takes
// (panelResponseHeaders builds its set the same way) - this function adds no
// header of its own and knows no policy.
//
// A refusal's body is the resolver's own sentence: a 415 or a 413 the user can
// read beats an empty pane that looks like a rendering bug, and the sentence
// never contains a filesystem path (PreviewDocument.Reason is built from the
// bare artifact name the caller sent).
func panelPreviewAssetResponse(doc panel.PreviewDocument) panelAssetResponse {
	resp := panelAssetResponse{
		Data:        doc.Body,
		ContentType: doc.ContentType,
		Status:      doc.Status,
		Reason:      http.StatusText(doc.Status),
		Headers:     strings.Join(doc.Headers, "\r\n"),
	}
	if len(resp.Data) == 0 && doc.Status != http.StatusOK && doc.Reason != "" {
		resp.Data = []byte(doc.Reason)
	}
	return resp
}

// panelPreviewQuoteDoor answers panel.message.send. It is the ONLY caller of
// PreviewResolver.QuoteForModel in this tree, and it is the door ticket 192
// AC#2 asks for: a preview's text may become model context through one place,
// and that place proves the text is in the artifact before it stamps.
//
// The envelope it reads, and why: `path` carries the same-origin preview URL
// the page used to fetch the bytes (PreviewURL(name), optionally with the
// offset/limit the user scrolled to - the query is ignored here, the whole
// artifact is the subject). `text` carries the selection. Everything else is
// ignored. A message that carries no preview address is NOT this door's: the
// message relay is ticket 35's socket, so it keeps the refusal by name.
type panelPreviewQuoteDoor struct {
	src   *panelPreviewSource
	audit panel.AuditFunc
	actor string
}

func newPanelPreviewQuoteDoor(src *panelPreviewSource, audit panel.AuditFunc, actor string) *panelPreviewQuoteDoor {
	return &panelPreviewQuoteDoor{src: src, audit: audit, actor: actor}
}

// HandleMessageRequest is the whole door. Its outcomes, in the order they are
// decided:
//
//	not a preview address      refused by name (ticket 35's relay, unbuilt)
//	artifact not servable      refused with the resolver's own status + sentence
//	no C25 engine / no proof   refused by QuoteForModel's own rules
//	proof + stamp              the mark is live, and the RELAY is refused by
//	                           name, because forwarding it is ticket 35's
//
// The last line is the honest half of this leg and the reason it is written as
// a refusal rather than a receipt: a handler that answered "ok" here would be
// telling the user their text went to the model while no code in this process
// can carry it there. What the refusal buys is the stamp - so the day the relay
// lands, the text it forwards is already marked by the door, and a relay that
// skips this door is a mutation the criterion in
// cmd/wisp/panel_preview_192_windows_test.go can aim at.
func (d *panelPreviewQuoteDoor) HandleMessageRequest(ctx context.Context, req panel.ComposerRequest) error {
	quoted := req.Text
	requestPath := previewPathPart192(req.Path)
	if !strings.HasPrefix(requestPath, panel.PreviewPathPrefix) {
		err := fmt.Errorf("%w: %q 的派发通路属票 35，本机未接入；这条请求里也没有预览引用（path 不是以 %q 开头的预览地址），因此本机既没有送出消息也没有盖戳",
			panel.ErrNoHandlerAttached, req.Method, panel.PreviewPathPrefix)
		d.record(req, "", err)
		return err
	}
	doc := d.src.Lookup(ctx, requestPath)
	if doc.Status != http.StatusOK {
		err := fmt.Errorf("%w: 预览引用指向的产物取不回来（status %d）：%s",
			panel.ErrPreviewRefused, doc.Status, doc.Reason)
		d.record(req, "", err)
		return err
	}
	resolver, oerr := d.src.resolver192()
	if oerr != nil {
		d.record(req, doc.Name, oerr)
		return oerr
	}
	// The one door. Everything QuoteForModel proves - the text is non-empty, it
	// is really inside THIS artifact's bytes (read back here, never taken from
	// the page), the artifact is small enough to have been checked - happens
	// inside this call, and its refusals are its own sentences.
	if _, qerr := resolver.QuoteForModel(ctx, d.src.marker192(), panelPreviewQuoteScope, doc.Name, quoted); qerr != nil {
		d.record(req, doc.Name, qerr)
		return qerr
	}
	// The stamp is live (QuoteForModel returns an error rather than a silent
	// "marked nothing" when Mark records nothing), so the audit line can say so
	// without echoing the text: the length is a fact, the content is the leak
	// the 35r7 ruler hunts, and no audit line in this package carries it.
	d.record(req, doc.Name, nil)
	return fmt.Errorf("%w: 这段文字已过 C25 盖戳（artifact:%s，scope=%s，%d 字），但消息派发通路属票 35，本机未接入：它没有进模型，也不会被静默丢掉",
		panel.ErrNoHandlerAttached, doc.Name, panelPreviewQuoteScope, len([]rune(quoted)))
}

// record writes the one audit line per attempt. It never carries the quoted
// text, only the artifact it was proven against and the outcome - the same
// discipline the 35r7 raw-leak ruler holds every panel line to.
func (d *panelPreviewQuoteDoor) record(req panel.ComposerRequest, artifact string, err error) {
	if d == nil || d.audit == nil {
		return
	}
	outcome := "stamped"
	if err != nil {
		outcome = "refused"
	}
	d.audit("panel: PREVIEW-QUOTE request=%q actor=%q artifact=%q textRunes=%d outcome=%s err=%v",
		req.RequestID, d.actor, artifact, len([]rune(req.Text)), outcome, err)
}

// previewPathPart192 cuts the envelope's path at its fragment. A fragment is
// never sent to a server, so it is not part of the address; the query IS
// (offset/limit), and the resolver is the one that parses it.
func previewPathPart192(raw string) string {
	raw = strings.TrimSpace(raw)
	if i := strings.IndexByte(raw, '#'); i >= 0 {
		return raw[:i]
	}
	return raw
}
