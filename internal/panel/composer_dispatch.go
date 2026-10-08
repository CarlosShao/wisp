package panel

// Ticket 33 slice A (H4 + H5) - the inbound hop's router.
//
// WHAT THIS IS. bridge.go already reads one postMessage string and says whether
// Go answers the method inside it (ParseComposerRequest, and the closed set in
// knownComposerMethod). What never existed is the next sentence: taking a parsed
// request to the handler that owns it. The 33-a1 census named that gap as the
// single common hole across H4-H9 (docs/evidence/s1/33-inbound-hop-design-a1.md
// §4), and it is the reason the mode write leg at composer_handlers.go:141 has
// been written and assembled with nobody able to ring it. This file is that one
// hop: raw in, method + args out, one door per whitelisted method. Four when
// slice A shipped; the two settings doors ticket 248 added make six, and the
// roster they are read from is bridge.go's, never a copy written here.
//
// IT IS NOT A SECOND GATE, and the reason is checkable rather than asserted.
// The whitelist decision is made exactly once, inside ParseComposerRequest, and
// this file reaches it and nothing else - it never calls knownComposerMethod
// itself, it never re-derives the answered set, and it spells no method name of
// its own: every case label below is one of bridge.go's exported constants. That
// last part is not tidiness. This package's own AST instrument
// (l2_grant_boundary_test.go, poolJudgedByRealGuard) collects every route-shaped
// name written anywhere in the package and asks the RUNNING guard about each one;
// a name the guard answers but the guard's own case list never named is reported
// as "a second switch/if chain ... exactly the shape a handler registry takes".
// Writing a route string here instead of using the constant would trip that nail
// on purpose, so the file cannot grow a private vocabulary by accident.
//
// THE default BRANCH REFUSES. A request whose method reached this function but
// matches no case is refused and audited, never accepted quietly. Today the
// branch is unreachable from the outside, because parse refuses unlisted names
// first - it exists for the day bridge.go gains a fifth method and this file has
// not been taught it, which is the failure that would otherwise be silent: a
// whitelisted name that nobody handles is a request the user made and no one
// answered. Same for a listed name with no handler attached in this assembly
// (workspace / attachment / message have no handler in this tree yet - see
// unattached). Fail-closed is the only state that costs nothing to add.
//
// ATTRIBUTION IS NOT THIS FILE'S TO TOUCH. The source and requestId inside an
// envelope are how the audit trail says who asked; a router that "helpfully"
// normalised them would manufacture the attribution it is supposed to carry.
// Nothing here assigns to any field of ComposerRequest. A tampered source never
// reaches the routing below at all - bridge.go's source check refuses it, and
// that refusal is the same door, not a copy of it (AGENTS.md §1.2: no
// panel-side L2 "allow", and Q-49's family is exactly what this defends).
//
// WHAT THIS DOES NOT DO, stated so nobody infers it from the file's existence:
// it opens no window and imports no WebView2 symbol (slice A is the half of the
// hop that needs neither); it writes no reply back to the page, which is H10 and
// belongs to the host. Slice A shipped with zero production callers; both callers
// exist today - cmd/wisp/panel_host_windows.go's dispatchRaw and
// cmd/wisp/panel_inbound.go (`wisp panel-inbound`, one envelope per stdin line) -
// and the slice-A reading, with its grep lines, is in
// docs/evidence/s1/33-minimal-inbound-hop-r1.md §④.

import (
	"context"
	"errors"
	"fmt"
)

// ErrNoHandlerAttached is the refusal for a method the whitelist answers while
// this assembly wires no handler for it. It is not ErrComposerRequest: the
// envelope was fine, the machine is incomplete, and collapsing those two would
// tell the user their own click was malformed.
var ErrNoHandlerAttached = errors.New("panel: 该方法在名册内，但本机未接入处理器")

// ModeRequestHandler answers one mode request. *ModeWriteHandler satisfies this
// shape as it stands, which is the point: the router reaches the handler ticket
// 114 built, not a parallel one.
type ModeRequestHandler interface {
	HandleModeRequest(ctx context.Context, req ComposerRequest) error
}

// WorkspaceRequestHandler answers one workspace request. The native leg it will
// call exists (workspace.go's RequestWorkspaceSwitch); the handler in front of it
// is ticket 186's work, so this file declares the socket and nothing more.
type WorkspaceRequestHandler interface {
	HandleWorkspaceRequest(ctx context.Context, req ComposerRequest) error
}

// AttachmentRequestHandler answers one attachment request. The guards exist
// (attachments.go's DecodeAttachmentPayload and AttachmentBroker.Ingest); the
// handler is a later ticket's, and a nil here refuses rather than pretending.
type AttachmentRequestHandler interface {
	HandleAttachmentRequest(ctx context.Context, req ComposerRequest) error
}

// MessageRequestHandler answers one message request.
type MessageRequestHandler interface {
	HandleMessageRequest(ctx context.Context, req ComposerRequest) error
}

// ConfigRequestHandler answers one settings request (ticket 248 AC#1/AC#8).
//
// Its signature is not the same as the other four, and the reason is a ruling
// rather than a slip: the credential write's VALUE must never live on the shared
// ComposerRequest (bridge.go's envelope serves every answered method, so a value
// field added there is readable from the message route too, and the frozen
// instrument's verdict vocabulary has no word that stops a value-shaped key).
// The leg therefore receives the raw envelope as well and decodes the credential
// into a write-only shape it owns alone, uses it once, and keeps nothing. Every
// other field this route reads comes off the parsed req.
//
// It returns the one sentence the user is told, where the other four handlers
// return none: for settings there is no second voice to defer to - H10's page
// text is the *approval* card's story, and the awaited binding reply is a plain
// return value from this hop (cmd/wisp/panel_host_windows.go's bind closure
// hands whatever Handle returns back to the page). Ticket 248 AC#8 needs a
// visible "when does this take effect", and a return value is the only
// page-reachable surface this tree has for it today.
type ConfigRequestHandler interface {
	HandleConfigRequest(ctx context.Context, req ComposerRequest, raw string) (string, error)
}

// ComposerDispatch is the inbound entry for one postMessage string: the host's
// receive callback hands it the text it was given and reads back the one
// sentence the user should see.
//
// It is assembled as a struct literal like ModeWriteHandler is, and every
// handler field is optional in the same way: nil is a refusal with a name, not a
// nil dereference and not a silence.
type ComposerDispatch struct {
	// Mode answers the mode request. In production this is the *ModeWriteHandler
	// that cmd/wisp assembles with the real write leg attached.
	Mode ModeRequestHandler
	// Workspace answers the workspace request. Untouched by this slice.
	Workspace WorkspaceRequestHandler
	// Attachment answers the attachment request. Untouched by this slice.
	Attachment AttachmentRequestHandler
	// Message answers the message request. Untouched by this slice.
	Message MessageRequestHandler
	// Config answers the two settings requests (ticket 248). nil is the same
	// refusal as the three above, not a silent drop: the two names are on the
	// whitelist, so a machine with no settings leg attached must say so by name.
	Config ConfigRequestHandler
	// Audit is the project's existing "[audit]" sink (workspace.go's AuditFunc).
	// A refusal is returned whether or not it is attached; a nil sink silences
	// the line, never the refusal.
	Audit AuditFunc
}

// Handle is the whole hop: raw postMessage text in, the sentence for the user
// out. The error is the same refusal that text describes, so a caller can count
// refusals without parsing prose.
//
// An accepted request returns no text except for the two settings routes: this
// slice deliberately does not invent a success line for mode / workspace /
// attachment / message, because deciding what the page is told about THOSE is
// H10's business and a router that writes its own confirmation would be the
// second voice in the room. A settings request has no other voice at all yet, and
// ticket 248 AC#8 requires "when this takes effect" to reach a page-visible
// surface, so the config handler's receipt is what comes back (dispatch's own
// "" for the four older doors is unchanged).
func (d *ComposerDispatch) Handle(ctx context.Context, raw string) (string, error) {
	// One gate, the existing one. Everything parse refuses never gets routed.
	req, err := ParseComposerRequest(raw)
	if err != nil {
		d.record(req, err)
		return RefusedEnvelopeForUser(req, err), err
	}
	receipt, err := d.dispatch(ctx, req, raw)
	if err != nil {
		return RefusedEnvelopeForUser(req, err), err
	}
	return receipt, nil
}

// dispatch is the routing table, kept apart from Handle so the backstop branch
// is reachable from inside this package: from outside, parse gets there first,
// and a default branch that no test can ever run is a branch nobody can prove
// refuses anything.
//
// The string it returns is the page-visible receipt, and only the settings doors
// ever fill it; every other branch returns "" with its outcome, so the sentence
// Handle hands back for those stays exactly what it was before ticket 248.
func (d *ComposerDispatch) dispatch(ctx context.Context, req ComposerRequest, raw string) (string, error) {
	switch req.Method {
	case MethodModeRequest:
		if d == nil || d.Mode == nil {
			return "", d.unattached(req)
		}
		return "", d.Mode.HandleModeRequest(ctx, req)
	case MethodWorkspaceRequest:
		if d == nil || d.Workspace == nil {
			return "", d.unattached(req)
		}
		return "", d.Workspace.HandleWorkspaceRequest(ctx, req)
	case MethodAttachmentAdd:
		if d == nil || d.Attachment == nil {
			return "", d.unattached(req)
		}
		return "", d.Attachment.HandleAttachmentRequest(ctx, req)
	case MethodMessageSend:
		if d == nil || d.Message == nil {
			return "", d.unattached(req)
		}
		return "", d.Message.HandleMessageRequest(ctx, req)
	case MethodConfigGet:
		if d == nil || d.Config == nil {
			return "", d.unattached(req)
		}
		return d.Config.HandleConfigRequest(ctx, req, raw)
	case MethodConfigSet:
		if d == nil || d.Config == nil {
			return "", d.unattached(req)
		}
		return d.Config.HandleConfigRequest(ctx, req, raw)
	default:
		return "", d.rosterMismatch(req)
	}
}

// rosterMismatch is the default branch's whole behaviour, kept as its own
// function so the branch is one line: what the mutation runs point at is this
// refusal, and a backstop nobody can aim a mutation at is not a backstop.
func (d *ComposerDispatch) rosterMismatch(req ComposerRequest) error {
	// Reachable only if bridge.go and this file disagree about the roster.
	// Refusing is the whole behaviour; the sentence says which two files to go
	// look at.
	err := fmt.Errorf("%w: 名册与派发表不一致（requestId=%q），处理器未调用", ErrRosterMismatch, req.RequestID)
	d.record(req, err)
	return err
}

// ErrRosterMismatch is the default branch's error: a method arrived here that
// the whitelist answered but this file has no case for.
var ErrRosterMismatch = errors.New("panel: 方法在白名单内而派发表未列出")

// unattached refuses a listed method whose handler socket is nil, and says so in
// the audit. This is the honest state of three of the six doors in this tree.
func (d *ComposerDispatch) unattached(req ComposerRequest) error {
	err := fmt.Errorf("%w: 方法 %q 的处理器未接入（requestId=%q），档位/工作区/附件/消息/设置均未变化",
		ErrNoHandlerAttached, req.Method, req.RequestID)
	d.record(req, err)
	return err
}

// record writes one audit line per refusal. Mirrors composer_handlers.go's
// record: a branch that returns an error without a line is how a dropped request
// comes to look like a handled one.
func (d *ComposerDispatch) record(req ComposerRequest, err error) {
	if d == nil || d.Audit == nil {
		return
	}
	d.Audit("panel: INBOUND-DISPATCH request=%q method=%q source=%q origin=%q err=%v detail=%q",
		req.RequestID, req.Method, req.Source, PanelModeOrigin, err,
		"处理器未被调用：该请求没有任何一侧发生变化")
}
