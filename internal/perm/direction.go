package perm

import "github.com/CarlosShao/wisp/internal/risk"

// RequiresCard reports whether a switch INTO to costs one answered L2 strong
// confirmation. Exactly one档 does: auto_approve (R20/M4, owner 2026-09-21).
//
// This is the ONE definition of the gate's judgement. It lives in this package,
// next to the switch it guards, because "which target needs a confirmation" is
// permission policy rather than panel vocabulary: perm.Store.Set (the layer every
// caller crosses) and panel.ModeWriteHandler (the layer a panel request crosses
// before a write leg sees it) both call this function instead of each spelling
// the comparison. Ticket 367 established that single-definition discipline for a
// direction predicate; ticket 372 reverted the JUDGEMENT (the survey at
// docs/evidence/s1/m4-survey.md:96 found no mainstream harness that gates on
// "wider") and kept the discipline.
//
// Every non-auto_approve target passes this predicate - widening, narrowing and
// no-change alike. Refusing a narrowing would guard nothing (getting stricter is
// never a risk), and refusing ask_high_risk would block a switch R20/M4 says
// costs no card; with no confirmation channel the READ side keeps that switch
// harmless (PermissionMode collapses to ask_every_step).
func RequiresCard(to risk.Mode) bool { return to == risk.ModeAutoApprove }

// CardGate is what one switch costs the layer judging it. It is the ONE gate
// judgement both transition points apply, and it is a value rather than a pair of
// booleans so "which of the three cases is this" is asked once (ticket 371).
type CardGate uint8

// The three cases JudgeCard can report.
const (
	// GateNoCard: the target is not the one档 R20/M4 cards. Proceed.
	GateNoCard CardGate = iota
	// GateNoChannel: the target costs a card and this layer holds no confirmation
	// channel. Refuse HERE, before a write leg is reached - the refusal belongs to
	// the layer that could not even try to ask.
	GateNoChannel
	// GateAsk: the target costs a card and a channel is attached. The asking layer
	// (perm.Store.Set, the only one that raises the card) must obtain an ANSWERED
	// allow before the switch proceeds; this layer itself may proceed to that
	// layer.
	GateAsk
)

// JudgeCard is the ONE predicate both transition points call:
// perm.Store.Set (store.go) and panel.ModeWriteHandler (composer_handlers.go)
// judge the same target with the same channel term by calling this function
// instead of each spelling `RequiresCard(to) && confirm == nil`. Ticket 367
// established that single-definition discipline for the direction rule; ticket
// 372 kept it for the target rule; ticket 371 extends it to the pair
// (target, channel) so a change to the RULE cannot land in one gate only.
//
// It does not judge the ANSWER. Whether an allow actually came back is a fact only
// the layer that raises the card can observe - the panel deliberately raises no
// second card (composer_handlers.go's header) - and it is judged there, against
// the Card credential ConfirmFunc hands back (store.go's Set).
func JudgeCard(to risk.Mode, channel bool) CardGate {
	if !RequiresCard(to) {
		return GateNoCard
	}
	if !channel {
		return GateNoChannel
	}
	return GateAsk
}

// IsWidening reports whether the switch from -> to removes questions instead of
// adding them: a raw `to > from` over risk.Mode's constants, which is only sound
// because internal/risk/mode.go orders them strict -> loose and says "Do not
// reorder these constants".
//
// It is NOT a gate any more (ticket 372). It is kept because the ladder order it
// reads is a premise two packages pin from both sides (TestIsWideningOnThe
// DocumentedLadder here and the ladder test in internal/panel), and because the
// AUDIT vocabulary still says "widened/narrowed" when it describes a switch. What
// no code may do is refuse a switch because this function returned true: that was
// ticket 367's rule, and it is the one the mainstream survey refuted.
func IsWidening(from, to risk.Mode) bool { return to > from }
