package panel

// The composer's view model (ticket 92): what the input row of the main panel
// shows, and what a user's intent looks like on its way OUT.
//
// The one-sentence design rule, from owner's R20/M5 ruling and PLAN.md:1588:
// the panel displays authority and requests it, it never holds or grants it.
// Two of the three things in the composer are therefore permission inputs -
// the mode (which decides whether the agent asks at all) and the workspace
// (which decides which tree R2/R3 judge against) - and a renderer that could
// write either would be a compromised renderer granting itself a permanent
// approval pass, which is worse than approving one operation.
//
// So this file contains no setter, no "ok" path, and no second channel:
//
//	mode          ModeView is read-only text + the legal names. A change is a
//	              ModeRequest, which reaches perm.Store.Set - the same call the
//	              CLI makes, which still costs the R20/M4 L2 strong confirmation
//	              and still writes the audit line.
//	workspace     WorkspaceView reports what native resolved. A change is a
//	              WorkspaceRequest handled by internal/tools' path scope (C26
//	              Resolve + the ticket 102/107 rewrite accounting), never by the
//	              page picking a string that then takes effect.
//	attachments   The page names a source; the native side reads, sniffs, stores
//	              (see attachments.go) and answers with a reference.
//
// The JSON keys of every type here are reconciled against frontend/src/lib/
// panel.ts by TestComposerContractTypesMatchFrontend in composer_test.go, the
// same two-way check ticket 77 pinned for the approval card.

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Snapshot is the one push the panel receives per state change (ticket 35 owns
// the pump). It is the whole truth: nothing on screen lives outside it. Moved
// here from approval_test.go by ticket 92, because a view model that only a
// test can build is a view model production never sends.
type Snapshot struct {
	Pending     []ApprovalCardView `json:"pending"`
	Results     []ResultChunk      `json:"results"`
	Composer    ComposerState      `json:"composer"`
	GeneratedAt string             `json:"generatedAt"`
}

// ResultChunk is one streamed assistant chunk.
type ResultChunk struct {
	CorrelationID string `json:"correlationId"`
	Text          string `json:"text"`
	Done          bool   `json:"done"`
}

// NewSnapshot builds the one push from native state. The composer section is
// MANDATORY (ticket 92 AC#4): the mode and the workspace a run is actually
// using travel with every snapshot, so a panel that closes and reopens - or a
// WebView process that dies - recovers exactly what Go holds rather than an
// empty input row that invites the user to assume the default.
func NewSnapshot(pending []ApprovalCardView, results []ResultChunk,
	composer ComposerState, now time.Time) Snapshot {
	if pending == nil {
		pending = []ApprovalCardView{}
	}
	if results == nil {
		results = []ResultChunk{}
	}
	if composer.Mode.Current == "" {
		// An unreadable mode is never rendered as the safest-sounding one.
		composer.Mode = ModeView{Current: "unknown", Names: risk.ModeNames()}
	}
	if composer.AttachmentMIMEs == nil {
		composer.AttachmentMIMEs = AcceptedMIMETypes()
	}
	return Snapshot{
		Pending:     pending,
		Results:     results,
		Composer:    composer,
		GeneratedAt: now.UTC().Format(time.RFC3339),
	}
}

// ModeView is the read-only rendering of the R20 permission mode.
//
// It has exactly three fields, and that is the security statement: a view model
// with no setter cannot be written through. The panel may show the current档
// and the legal names; every path from there to perm.Store.Set goes through the
// native confirmation, and the panel is not on it.
type ModeView struct {
	// Current is risk.Mode.String(): "ask_every_step" | "ask_high_risk" |
	// "auto_approve", or "unknown" when the native read failed.
	Current string `json:"current"`
	// Names is risk.ModeNames() verbatim: the display list, not a choice set.
	Names []string `json:"names"`
	// L2ConfirmNames lists the transitions that cost an L2 strong confirmation
	// (R20/M4). Displayed so the request is not a surprise, enforced natively.
	L2ConfirmNames []string `json:"l2ConfirmNames"`
}

// L2ConfirmTarget is the single mode whose arrival needs R20/M4's L2 card.
// It is named once here so the UI text and the native gate cannot disagree
// about which档 removes questions.
const L2ConfirmTarget = risk.ModeAutoApproveName

// NewModeView renders a mode the native side actually holds.
func NewModeView(m risk.Mode) ModeView {
	return ModeView{
		Current:        m.String(),
		Names:          risk.ModeNames(),
		L2ConfirmNames: []string{L2ConfirmTarget},
	}
}

// ModeUnknownView is what AC#4's "never read as safe" case renders as: a mode
// that could not be read is reported as unknown, not as the default档.
func ModeUnknownView() ModeView {
	return ModeView{Current: "unknown", Names: risk.ModeNames(), L2ConfirmNames: []string{L2ConfirmTarget}}
}

// ModeRequest is the panel's INTENT to change the mode. It carries no authority:
// handling one means calling perm.Store.Set with origin "panel", and Set is the
// place R20/M4's confirmation and the audit line live. A handler that wrote the
// mode anywhere else would be a second channel, which this ticket forbids.
type ModeRequest struct {
	To string `json:"to"`
	// CorrelationID ties the request to whatever card the gate raises.
	CorrelationID string `json:"correlationId"`
}

// Parse resolves the requested mode with risk's own parser, so an unknown
// spelling is an error rather than a default. The empty string is refused here
// even though risk.ParseMode accepts it as the default: that tolerance belongs
// to a config key that may be absent, and a REQUEST with no target is not an
// absent key - it is a page that sent nothing.
func (r ModeRequest) Parse() (risk.Mode, error) {
	to := strings.TrimSpace(r.To)
	if to == "" {
		return 0, fmt.Errorf("panel: mode request carries no target (one of %s)",
			strings.Join(risk.ModeNames(), "/"))
	}
	m, err := risk.ParseMode(to)
	if err != nil {
		return 0, err
	}
	if !m.Valid() {
		return 0, fmt.Errorf("panel: mode %q is not one of %s", r.To, strings.Join(risk.ModeNames(), "/"))
	}
	return m, nil
}

// WorkspaceView reports the active workspace as native resolved it.
//
// Spelling is what the user named, Canonical is what the path scope authorised,
// and the two booleans are ticket 102's account: a workspace whose expansion
// moved it onto another tree, or that crossed a reparse point, must be visible
// as such in the panel instead of being presented as "your folder".
type WorkspaceView struct {
	// Set is false while no workspace has been chosen, i.e. while the config's
	// [fs] allowed_dirs roots alone decide scope.
	Set bool `json:"set"`
	// Spelling is the user's own string, verbatim.
	Spelling string `json:"spelling"`
	// Canonical is the C26-resolved tree currently in force ("" when unset).
	Canonical string `json:"canonical"`
	// Reparse reports an authorised reparse traversal (exception-listed).
	Reparse bool `json:"reparse"`
	// Rewritten reports that expansion (%VAR%, ~) substituted part of the
	// spelling, so Canonical may name a different tree than Spelling does.
	Rewritten bool `json:"rewritten"`
	// Reason explains an unset or narrowed workspace; never empty when Set is
	// false, because "no workspace" is a fact the user must be able to read.
	Reason string `json:"reason"`
}

// WorkspaceRequest is the panel's intent to narrow the session to a folder.
type WorkspaceRequest struct {
	Path string `json:"path"`
}

// AttachmentRef is the composer's per-attachment row (see attachments.go).

// ComposerState is the composer section of a snapshot.
type ComposerState struct {
	Mode             ModeView        `json:"mode"`
	Workspace        WorkspaceView   `json:"workspace"`
	Attachments      []AttachmentRef `json:"attachments"`
	AttachmentMIMEs  []string        `json:"acceptedAttachmentMimes"`
	MaxAttachmentB   int64           `json:"maxAttachmentBytes"`
	AttachmentReason string          `json:"attachmentError"`
}

// NewComposerState assembles the composer section from native state.
func NewComposerState(mode risk.Mode, ws WorkspaceView, atts []AttachmentRef, maxBytes int64) ComposerState {
	if atts == nil {
		atts = []AttachmentRef{}
	}
	return ComposerState{
		Mode:            NewModeView(mode),
		Workspace:       ws,
		Attachments:     atts,
		AttachmentMIMEs: AcceptedMIMETypes(),
		MaxAttachmentB:  maxBytes,
	}
}

// OutgoingMessage is one composer send: text plus the attachment references the
// native side already accepted. The panel never ships bytes inside a message
// body and the agent never sees a path the artifacts store did not approve.
type OutgoingMessage struct {
	Text        string          `json:"text"`
	Attachments []AttachmentRef `json:"attachments"`
}

// ForAgent renders what the agent receives for one composed message. Every
// accepted attachment is named with the artifact file that holds its bytes;
// every refused one is repeated verbatim, because AC#2's judgement is "do not
// eat the user's intent": a dropped attachment has to be visible in the very
// text the model reads.
//
// Video is described as stored bytes and nothing more - semantic understanding
// of video is Q-28, which is not this ticket.
func (m OutgoingMessage) ForAgent(dir string) string {
	text := strings.TrimRight(m.Text, " \t\r\n")
	var b strings.Builder
	b.WriteString(text)
	var refused []string
	for _, a := range m.Attachments {
		if !a.Stored {
			refused = append(refused, fmt.Sprintf("%s: %s", a.Name, a.Reason))
			continue
		}
		if text != "" || b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("[附件] name=%s mime=%s kind=%s bytes=%d path=%s",
			a.Name, a.MIME, a.Kind, a.SizeBytes, joinArtifact(dir, a.Artifact)))
		if a.Deduplicated {
			b.WriteString(" (与已有附件同内容，未重复写入)")
		}
	}
	if len(refused) > 0 {
		b.WriteString("\n[附件未送达] ")
		b.WriteString(strings.Join(refused, "; "))
		b.WriteString("\n（这些附件没有进入本轮，请用户确认后重试或改用受支持的类型）")
	}
	return b.String()
}

// joinArtifact builds a display path for a bare artifact name. It does not
// re-resolve or act on anything - the store already refused a name that is not
// bare (ticket 76), so this is string shaping for a message, not a path leg.
// filepath.Join, not a hard-coded separator: ticket 75 caught exactly that in
// this repository, and a comparison string no OS call can open is a lie on
// POSIX even though the product ships on Windows.
func joinArtifact(dir, name string) string {
	if dir == "" {
		return name
	}
	return filepath.Join(dir, name)
}
