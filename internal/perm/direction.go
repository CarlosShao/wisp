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
