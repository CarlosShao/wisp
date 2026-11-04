package panel

// Ticket 114 AC#2 - the native side of the composer's mode write (R-92-2).
//
// WHY A GATE HERE AT ALL, NOW THAT perm.Store HAS ONE TOO. Both layers judge the
// same thing and they judge it with the same function (perm.RequiresCard), so this
// is not a second opinion that can drift - it is the same opinion applied at the
// two points a mode request has to cross. The store's gate is what makes the
// invariant true for every caller, including the ones that never come through
// this file; this gate is what refuses a panel request BEFORE a write leg is
// even asked, so the refusal the page sees is a panel refusal with a panel
// error code rather than whatever the store happened to say.
//
// HISTORY THAT MATTERS, TWICE. Until ticket 367 the two layers judged
// DIFFERENTLY - the store asked "is the target auto_approve?" while this file
// asked "is the direction widening?" - so a panel request that widened to
// ask_high_risk was refused here and would have been applied unchanged by Set had
// anything else called it. 367 made both layers judge the DIRECTION and closed
// that gap, but "any widening needs a leg" is not what the owner settled
// (2026-09-21: only the switch INTO full auto costs a card), and no surveyed
// mainstream harness gates on "wider" (docs/evidence/s1/m4-survey.md:96). Ticket
// 372 reverted the JUDGEMENT to the target and kept the single predicate:
// perm.RequiresCard is defined once and called here and in Set.
//
// THE RULE, in owner's R20/M4: a request whose TARGET is auto_approve needs a
// leg; every other target passes, wider or narrower. The missing leg is made safe
// for the middle档 on the READ side instead: perm.Store.PermissionMode answers
// ask_every_step when the store has no confirmation channel, so the档 a legless
// process actually runs is the strictest one even after ask_high_risk has been
// persisted (ticket 372-r2).
//
// THE HANDLER NEVER CALLS Confirm. Asking belongs to the write leg: perm.Store.Set
// raises the single C18 card for auto_approve (R20/M4), and an ask here on top
// of it would be two cards for one click - two 300s windows, two timeouts, two
// rejects. What Confirm's PRESENCE buys is exactly that property, and nothing
// else: an assembly that hands this handler no leg cannot reach auto_approve.
// TestOnlyTheAutoApproveTargetIsRefusedWhenNoL2ConfirmLegIsAttached pins that Set
// is not reached for that target; TestTheModeHandlerRaisesNoSecondCard pins that
// Confirm is not called either.
//
// TICKET 371: THE SAME RULE, AND WHY THE ANSWER IS NOT JUDGED HERE. Since 371 the
// card must come back with an ANSWERED allow (perm.Card), not merely a nil error.
// This layer judges the half it can observe - the target and whether a channel is
// attached - by calling the SAME predicate the store calls (perm.JudgeCard), so
// the rule has one definition and two call sites. It cannot judge the answer
// itself: it raises no card, so "was an allow answered" is a fact only the layer
// that asked can see. What it owes the operator instead is that the store's
// refusal is NAMED here with the same classification the store booked
// (perm.RefusalResult) and recorded with this request's id, so the two records
// cannot describe one refusal as two different things.
//
// AUDIT IS NOT THE GATE'S TRADE. Every refusal written here emits its own audit
// line before returning, because the write leg never ran and would otherwise
// have recorded nothing; every request that passes goes to Set, which records
// all three档 (R20). A refusal the write leg made on this request's behalf is
// recorded here too - it is the layer that knows the requestId, and one refusal
// deserves one line at each layer that judged it, both carrying the same
// perm.Result. No branch of this file is silent, and adding the gate removes no
// line the wiring would have written.
//
// WHO CALLS THIS NOW: the earlier "nothing calls this handler yet" was true
// while the inbound hop did not exist, and it is not true now. Both production
// doors reach this handler through one router - the page's door at
// panel_host_windows.go:929/948, `wisp panel-inbound` at panel_inbound.go:173 -
// and ComposerDispatch.dispatch routes the mode case here (composer_dispatch.go:204).

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

// ErrNoL2Confirm is the refusal AC#2 exists for: a request for the one target
// that costs a card arrived with no confirmation leg attached.
var ErrNoL2Confirm = errors.New("panel: 确认腿未接入，档位不得变宽")

// ErrNoModeWriter means there is no write leg to ask for at all.
var ErrNoModeWriter = errors.New("panel: 档位写入腿未接入")

// ErrNoCurrentMode means the档 in force could not be read, so a refusal would
// have no attributable starting point and the audit line would name nothing.
// Failing closed here is not paranoia about a nil: a writer that reports a档 it
// is not actually running would put the wrong "from" on every record this handler
// writes.
var ErrNoCurrentMode = errors.New("panel: 当前档位不可读")

// ModeWriter is the write leg a panel mode request is handed to, and it is the
// only one: *perm.Store satisfies this shape as it stands, so the handler reaches
// the same Set the CLI reaches instead of a second channel (composer.go's rule).
// Declared locally, with no import of internal/perm, so this package's
// dependency face does not move for one ticket's gate.
type ModeWriter interface {
	// PermissionMode is the EFFECTIVE档 in force: a store with no confirmation
	// channel answers ask_every_step (ticket 372-r2), which is what this handler
	// must name in a refusal or an audit line.
	PermissionMode() risk.Mode
	// Set switches and persists, confirms if R20/M4 says so, and audits.
	Set(ctx context.Context, to risk.Mode, origin, actor string) error
}

// ModeConfirm is one L2 strong-confirmation leg (R20/M4), in perm's own
// credential contract (ticket 371): it reports the answered allow as a perm.Card,
// the same value perm.ConfirmFunc returns and the same bit perm.Store.Set judges.
// The handler reads its presence and never its body; the shape mirrors the
// confirmation the write leg takes so an assembly can hand both the same thing,
// and so the credential cannot be dropped between the two layers.
type ModeConfirm func(ctx context.Context, from, to risk.Mode, origin, actor string) (perm.Card, error)

// ModeWriteHandler is the native handler for one panel.mode.request.
type ModeWriteHandler struct {
	// Modes is the single write leg. nil = every request is refused.
	Modes ModeWriter
	// Confirm is the L2 leg whose absence closes the gate for a request whose
	// target costs a card (R20/M4: auto_approve, and only it). Never called here;
	// see the header. Since ticket 371 it is the credential-returning shape, so a
	// limb that reports "no answered allow" can reach the store's gate through the
	// same contract instead of being flattened to an error on the way.
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
	if from == to {
		// A no-change is not a request for a card. Set records it as no-change and
		// returns without asking; refusing it here would make this layer stricter
		// than the store about a request that changes nothing. It is also the row
		// that lets a legless host ask for the档 it is ALREADY effectively running
		// (config auto_approve, read-side downgrade to ask_every_step) and have
		// that stricter value land on disk.
		return h.Modes.Set(ctx, to, PanelModeOrigin, h.actor())
	}
	// Ticket 372: the gate is the TARGET, and it is the same predicate
	// perm.Store.Set applies. Ticket 371 folded the channel term into that
	// predicate too (perm.JudgeCard), so "which targets cost a card, and does
	// this layer hold a channel for them" is answered once, not spelled twice.
	// A request for auto_approve with no leg is refused here, before the write
	// leg sees it, so the page gets a panel error code; every other target
	// reaches Set, which judges the same rule and the ANSWER as well.
	if perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel {
		refused := fmt.Errorf("%w: 请求 %q 要把档位从 %s 切到 %s（R20/M4：只有全自动要确认），而本机没有接入 L2 确认腿",
			ErrNoL2Confirm, req.RequestID, from, to)
		h.record(req, refused, from.String(), to.String())
		return refused
	}
	// The leg is attached, so R20/M4's ask is the write leg's to make - exactly
	// one card, raised by Set, for the one target that costs it, and only an
	// ANSWERED allow settles it (ticket 371). Set judges the same rule with the
	// same predicate and refuses on its own; this layer exists so a refusal made
	// before Set carries a panel error code, and so the refusal Set makes is
	// booked here with this request's id and the same perm.Result.
	if err := h.Modes.Set(ctx, to, PanelModeOrigin, h.actor()); err != nil {
		h.recordStoreRefusal(req, err, from.String(), to.String())
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
//
// result= is perm.RefusalResult's classification, i.e. the SAME value perm.Store
// booked on its own MODE-SWITCH row for the same refusal. That is what makes
// "two gates, one rule" a reading instead of a claim: a card that was never
// answered reads result=card-not-answered here and there, and a missing channel
// reads result=refused-confirm at both.
func (h *ModeWriteHandler) record(req ComposerRequest, err error, from, to string) {
	h.recordAs(req, err, from, to, "档位未变化：写入腿没有被调用")
}

// recordStoreRefusal books the refusal the WRITE LEG made on this request's
// behalf. The store wrote its own row (it always does); this one exists because
// the store's row cannot name the panel's requestId, and because the panel is a
// layer that judged the same rule and must not leave its half of the story
// silent (ticket 371-b). It carries the same perm.Result the store booked.
func (h *ModeWriteHandler) recordStoreRefusal(req ComposerRequest, err error, from, to string) {
	h.recordAs(req, err, from, to, "档位未变化：写入腿按同一条规则拒绝了这一发")
}

// recordAs is the one emitter both refusal paths share, so the two lines cannot
// drift into two vocabularies.
func (h *ModeWriteHandler) recordAs(req ComposerRequest, err error, from, to, detail string) {
	if h == nil || h.Audit == nil {
		return
	}
	h.Audit("panel: MODE-REFUSED request=%q from=%q to=%q origin=%q actor=%q result=%s err=%v detail=%q",
		req.RequestID, from, to, PanelModeOrigin, h.actor(), perm.RefusalResult(err), err, detail)
}
