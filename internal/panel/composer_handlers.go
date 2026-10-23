package panel

// Ticket 114 AC#2 - the native side of the composer's mode write (R-92-2).
//
// WHY A GATE HERE AT ALL, NOW THAT perm.Store.Has ONE TOO. Both layers judge the
// same thing and they judge it with the same function (perm.IsWidening), so this
// is not a second opinion that can drift - it is the same opinion applied at the
// two points a widening has to cross. The store's gate is what makes the
// invariant true for every caller, including the ones that never come through
// this file; this gate is what refuses a panel request BEFORE a write leg is
// even asked, so the refusal the page sees is a panel refusal with a panel
// error code rather than whatever the store happened to say.
//
// HISTORY THAT MATTERS: until ticket 367 the two layers judged DIFFERENTLY - the
// store asked "is the target auto_approve?" while this file asked "is the
// direction widening?" - so a panel request that widened to ask_high_risk was
// refused here and would have been applied unchanged by Set had anything else
// called it. docs/evidence/s1/114-ac1-status-table.md ②.1 measured exactly that
// hole. The single predicate is the fix for the class, not for that one pair.
//
// THE RULE, in owner's R20 direction-asymmetry as settled for this ticket: a
// request that makes the档 WIDER needs a leg; a request that makes it NARROWER
// never does, and must not be blocked. Refusing to become stricter guards
// nothing, and the档 that removes questions is the only one worth guarding.
//
// THE HANDLER NEVER CALLS Confirm. Asking belongs to the write leg: perm.Store.Set
// raises the single C18 card for auto_approve (R20/M4), and an ask here on top
// of it would be two cards for one click - two 300s windows, two timeouts, two
// rejects. What Confirm's nil-ness buys is exactly AC#2's property, and nothing
// else: an assembly that hands this handler no leg cannot widen anything.
// TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached pins that
// Set is not reached; TestTheModeHandlerRaisesNoSecondCard pins that Confirm is
// not called either.
//
// AUDIT IS NOT THE GATE'S TRADE. Every refusal written here emits its own audit
// line before returning, because the write leg never ran and would otherwise
// have recorded nothing; every request that passes goes to Set, which records
// all three档 (R20). No branch of this file is silent, and adding the gate
// removes no line the wiring would have written.
//
// WHAT THIS IS NOT: a path. Nothing calls this handler yet, because the
// WebView2 "event -> ParseComposerRequest" hop does not exist in this tree
// (status table ①.1; tickets 33/35 are still ready-for-agent). AC#2 is the
// precondition for that hop, not a substitute for it - the panel still cannot
// change the档, and this file must not be read as if it could.

import (
	"context"
	"errors"
	"fmt"

	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// PanelModeOrigin is the origin stamp on every switch the panel asked for. It is
// perm's audit vocabulary (internal/perm/store.go: "panel" / "ball" / "cli" /
// "config-file"), and the write leg owns the record - this only names the door.
const PanelModeOrigin = "panel"

// ErrNoL2Confirm is the refusal AC#2 exists for: a widening request arrived with
// no confirmation leg attached.
var ErrNoL2Confirm = errors.New("panel: 确认腿未接入，档位不得变宽")

// ErrNoModeWriter means there is no write leg to ask for at all.
var ErrNoModeWriter = errors.New("panel: 档位写入腿未接入")

// ErrNoCurrentMode means the档 in force could not be read, so "is this request
// wider?" has no answer. Failing closed here is not paranoia about a nil: a
// writer that reports a loose档 it is not actually running would otherwise turn
// every widening request into a narrowing one and switch this whole gate off.
var ErrNoCurrentMode = errors.New("panel: 当前档位不可读")

// ModeWriter is the write leg a panel mode request is handed to, and it is the
// only one: *perm.Store satisfies this shape as it stands, so the handler reaches
// the same Set the CLI reaches instead of a second channel (composer.go's rule).
// Declared locally, with no import of internal/perm, so this package's
// dependency face does not move for one ticket's gate.
type ModeWriter interface {
	// PermissionMode is the档 currently in force.
	PermissionMode() risk.Mode
	// Set switches and persists, confirms if R20/M4 says so, and audits.
	Set(ctx context.Context, to risk.Mode, origin, actor string) error
}

// ModeConfirm is one L2 strong-confirmation leg (R20/M4). The handler reads its
// nil-ness and never its body; the shape mirrors the confirmation the write leg
// takes so an assembly can hand both the same thing.
type ModeConfirm func(ctx context.Context, from, to risk.Mode, origin, actor string) error

// ModeWriteHandler is the native handler for one panel.mode.request.
type ModeWriteHandler struct {
	// Modes is the single write leg. nil = every request is refused.
	Modes ModeWriter
	// Confirm is the L2 leg whose absence closes the gate for widening
	// requests. Never called here; see the header.
	Confirm ModeConfirm
	// Audit is the project's existing "[audit]" sink (workspace.go's AuditFunc).
	// nil is tolerated the way perm tolerates it: the refusal still returns.
	Audit AuditFunc
	// Actor names who triggered it, as perm's Switch record wants it. Empty
	// stays "unknown" - an anonymous request is auditable as anonymous, which
	// beats inventing a name.
	Actor string
}

// HandleModeRequest answers one parsed panel.mode.request (bridge.go's
// MethodModeRequest). It returns the refusal it made, never a swallowed one.
func (h *ModeWriteHandler) HandleModeRequest(ctx context.Context, req ComposerRequest) error {
	if req.Method != "" && req.Method != MethodModeRequest {
		err := fmt.Errorf("panel: %q 不是档位请求，本处理器不回答（requestId=%q）", req.Method, req.RequestID)
		h.record(req, err, "", "")
		return err
	}
	if h == nil || h.Modes == nil {
		err := fmt.Errorf("%w: 面板的档位请求没有写入腿，requestId=%q 被拒绝", ErrNoModeWriter, req.RequestID)
		h.record(req, err, "", "")
		return err
	}
	to, err := ModeRequest{To: req.To, CorrelationID: req.RequestID}.Parse()
	if err != nil {
		h.record(req, err, "", "")
		return err
	}
	from := h.Modes.PermissionMode()
	if !from.Valid() {
		refused := fmt.Errorf("%w: 无法判断 %q 比现在更宽还是更严，按 fail-closed 拒绝", ErrNoCurrentMode, to)
		h.record(req, refused, "(不可读)", to.String())
		return refused
	}
	if perm.IsWidening(from, to) && h.Confirm == nil {
		refused := fmt.Errorf("%w: 请求 %q 要把档位从 %s 放宽到 %s，而本机没有接入 L2 确认腿（R20/M4：只有变宽要确认）",
			ErrNoL2Confirm, req.RequestID, from, to)
		h.record(req, refused, from.String(), to.String())
		return refused
	}
	// The leg is attached, so R20/M4's ask is the write leg's to make - exactly
	// one card, raised by Set, for the one档 that costs it. Set judges the same
	// direction with the same predicate and would refuse the same request on its
	// own if nothing had; this layer exists so the refusal carries a panel error.
	if err := h.Modes.Set(ctx, to, PanelModeOrigin, h.actor()); err != nil {
		return err
	}
	return nil
}

// actor is the audit's who, defaulted rather than invented.
func (h *ModeWriteHandler) actor() string {
	if h == nil || h.Actor == "" {
		return "unknown"
	}
	return h.Actor
}

// record writes the one audit line a refused request deserves. A nil sink is
// tolerated (perm's Options.Logf is optional too) but never turns a refusal into
// a silence at this level - the error still goes back to the caller.
func (h *ModeWriteHandler) record(req ComposerRequest, err error, from, to string) {
	if h == nil || h.Audit == nil {
		return
	}
	h.Audit("panel: MODE-REFUSED request=%q from=%q to=%q origin=%q actor=%q err=%v detail=%q",
		req.RequestID, from, to, PanelModeOrigin, h.actor(), err,
		"档位未变化：写入腿没有被调用")
}
