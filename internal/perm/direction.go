package perm

import "github.com/CarlosShao/wisp/internal/risk"

// RequiresCard reports whether the switch from -> to costs one answered L2 strong
// confirmation. Since R9 (owner 2026-11-05, ledger A1541 item 4: "every loosening
// shows a confirmation card") the answer is exactly "this switch removes
// questions", i.e. IsWidening(from, to). A tightening and a no-change cost nothing
// - getting stricter is never guarded (PLAN.md:132) - and the middle tier no
// longer rides free.
//
// HISTORY, because two earlier rulings shaped this line and neither is an accident:
//
//   - Ticket 367 made the direction the gate; ticket 372 reverted it to the TARGET
//     ("only auto_approve") because no surveyed mainstream harness REFUSES a switch
//     for being wider (docs/evidence/s1/m4-survey.md:96).
//   - R9 is the owner's own product decision (ledger A1541), and it is not a
//     re-litigation of that survey: what a loosening COSTS is a card, and only an
//     assembly with no confirmation channel still turns that into a refusal. The
//     config-file channel (D36's rule 1, internal/config/manager.go) has carded
//     every loosening since ticket 90; R9 makes the runtime channel agree with it
//     instead of letting the same request read differently at the two doors.
//
// It is the ONE definition of "does this switch need a card". perm.Store.Set,
// panel.ModeWriteHandler and command.Roster call JudgeCard below rather than each
// spelling the comparison, so a change to the rule cannot land at one door only.
func RequiresCard(from, to risk.Mode) bool { return IsWidening(from, to) }

// CardGate is what one switch costs the layer judging it. It is the ONE gate
// judgement the transition points apply, and it is a value rather than a pair of
// booleans so "which of the three cases is this" is asked once (ticket 371).
type CardGate uint8

// The three cases JudgeCard can report.
const (
	// GateNoCard: the switch removes no questions (a tightening, or the no-change
	// its caller already handled). Proceed.
	GateNoCard CardGate = iota
	// GateNoChannel: the switch costs a card and this layer holds no confirmation
	// channel. Refuse HERE, before a write leg is reached - the refusal belongs to
	// the layer that could not even try to ask. Since R9 this covers every
	// loosening, so a host with no card surface (the CLI pipe) cannot loosen at
	// all; that is the ruling's accepted cost, named in
	// docs/evidence/s1/R9-impl.md.
	GateNoChannel
	// GateAsk: the switch costs a card and a channel is attached. The asking layer
	// (perm.Store.Set, the only one that raises the card) must obtain an ANSWERED
	// allow before the switch proceeds; this layer itself may proceed to that
	// layer.
	GateAsk
)

// JudgeCard is the ONE predicate every transition point calls:
// perm.Store.Set (store.go), panel.ModeWriteHandler (composer_handlers.go) and
// command.Roster.Judge (roster.go) judge the same switch with the same channel
// term by calling this function instead of each spelling
// `RequiresCard(from, to) && confirm == nil`. Ticket 367 established that
// single-definition discipline for the direction rule; ticket 372 kept it for the
// target rule; ticket 371 extended it to the pair (switch, channel); R9 puts the
// direction back in the first term and adds the roster as a third caller, so a
// change to the RULE cannot land in one gate only.
//
// It does not judge the ANSWER. Whether an allow actually came back is a fact only
// the layer that raises the card can observe - the panel deliberately raises no
// second card (composer_handlers.go's header) - and it is judged there, against
// the Card credential ConfirmFunc hands back (store.go's Set).
func JudgeCard(from, to risk.Mode, channel bool) CardGate {
	if !RequiresCard(from, to) {
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
// It is the gate's first term since R9 (RequiresCard above), and it is also the
// AUDIT vocabulary: "widened/narrowed" is how a switch record is described. The
// ladder order it reads is a premise three packages pin from their own side
// (TestIsWideningOnTheDocumentedLadder here, the ladder test in internal/panel and
// the roster's matrix in internal/command), so a reorder breaks them rather than
// silently flipping one.
func IsWidening(from, to risk.Mode) bool { return to > from }
