package approval

// Ticket 375: 「能问」is NOT 「非 nil」.
//
// 372 closed one hole and its test was PRESENCE: a nil Confirm (no second leg at
// all) collapses the effective mode to ask_every_step, so no mode can silence a
// question that has no leg behind it (perm/store.go's downgrade). Presence is
// not the same fact as answerability. Four production assemblies hand the store
// a NON-nil leg, and at the moment of a call that leg may have nobody behind it
// at all - a console run with no interactive input stream, a resident process
// whose floating ball window never came up. Those runs kept the looser档, and on
// the L1 route the window's own polarity finished the job: 「窗口到点未被否决」
// MEANS EXECUTE, so "nobody could see this card" was resolved as consent
// (375-a1's READ-S: strictest档, no answer source, 2.04s, file written).
//
// This file is the action-layer half of the fix (AC#1 形 乙). It adds ONE fact
// and no authority:
//
//	AnswerSource - the host's LIVE statement, asked at the moment of the call,
//	               that a human could answer a card raised right now.
//
// Two consumers read the SAME value, which is why it is one type and not two
// independent injections:
//
//   - Gate.PendingWindow (the L1 route) refuses to open a window nobody could
//     see. This is what turns READ-S ("strictest档 + no answer source" ->
//     2.04s -> written) into a refusal.
//   - UnsilencedMode is the mode read the bridge consults once per call: a mode
//     that would only SILENCE a question answers with the strictest档 when
//     nobody can answer. This is what makes READ-F ("auto_approve + no answer
//     source" -> 0 cards -> written) reach the gate at all.
//
// Both halves are needed and neither is sufficient: without the guard, the
// downgrade only buys READ-S's 2s-then-execute; without the downgrade, the L1
// route is never reached for a silenced call. The ticket's boundary 2 is
// untouched by both: the L2 route keeps its own fail-closed surface check and
// its 到点即拒 polarity, because an unanswered L2 card is already a reject.
//
// nil and false are deliberately NOT the same answer:
//
//	nil    the host declared nothing about answerability. Every host built
//	       before this ticket - and every test fixture that only exercises the
//	       routing - is in this state, and its behaviour is kept byte for byte.
//	       Refusing all of them would not be a fix, it would be a purge
//	       (positive control: with an answer source, L1 and L2 do not move).
//	false  the host DECLARED that no one can answer right now. Only the host can
//	       report this about itself, and it is the state this ticket is about.
//
// The production composition roots therefore all declare one (cmd/wisp/run.go's
// console leg answers from its own reply surface; the resident leg answers from
// the floating ball it bound at boot). Both are closures, not snapshots: the
// resident ball is bound 85 lines AFTER the panel chain is handed its leg
// (resident_windows.go:224 vs :309), which is exactly why a boot-time reading
// could not cover assembly D (census 375-a1 §②).

import (
	"sync"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// AnswerSource reports, at the moment it is called, whether a human could answer
// a confirmation card this process would raise right now. It is a read face and
// conveys no authority: it decides nothing by itself, it only lets the two
// consumers above tell 「有一条腿」 from 「这条腿此刻答得出」.
//
// It must not block (the bridge reads it on the tool-call path) and it must be
// safe for concurrent use.
type AnswerSource func() bool

// answerable is the L1 route's reading of that fact. An undeclared host answers
// true, i.e. keeps the behaviour it had before ticket 375.
func (g *Gate) answerable() bool {
	if g.answer == nil {
		return true
	}
	return g.answer()
}

// AnswerSource exposes the predicate this gate was built with (nil when the host
// declared nothing). It is a read face for assembly rulers: "this process can
// tell whether a card is answerable at all" has to be checkable from outside
// this package, or a host that forgot to install one stays invisible until the
// hole is used - the same reason Gate.Window and Gate.Queue are read faces.
func (g *Gate) AnswerSource() AnswerSource { return g.answer }

// UnsilencedMode wraps a tools.ModeSource into the mode read the bridge should
// be wired with (tools.Options.Modes).
//
// It answers the wrapped source's mode unchanged, EXCEPT when
//
//	the host declared that nobody can answer (answer() == false), and
//	the configured mode is looser than ask_every_step,
//
// in which case it answers ask_every_step - 372's own downgrade target, for
// 372's own reason: a mode that removes questions must not be honoured by a
// process that cannot ask them. The L1 route then reaches the gate instead of
// being silenced, and the gate's own guard refuses it, so the call is REFUSED
// rather than run unasked.
//
// It books ONE audit line when it starts downgrading (and again only after the
// configured mode changed or the answer source came back and was lost): the read
// face is called once per tool call, so a line per call would bury the ledger.
func UnsilencedMode(src tools.ModeSource, answer AnswerSource, logf func(string, ...any)) tools.ModeSource {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &unsilencedMode{src: src, answer: answer, logf: logf}
}

// unsilencedMode is UnsilencedMode's state: the wrapped source plus the one
// memory the audit line needs (which configured mode was last announced as
// downgraded).
type unsilencedMode struct {
	src    tools.ModeSource
	answer AnswerSource
	logf   func(string, ...any)

	mu           sync.Mutex
	announced    risk.Mode
	announcedSet bool
}

// PermissionMode implements tools.ModeSource.
func (u *unsilencedMode) PermissionMode() risk.Mode {
	if u.src == nil {
		return risk.DefaultMode()
	}
	configured := u.src.PermissionMode()

	if u.answer != nil && !u.answer() && configured != risk.DefaultMode() {
		u.announce(configured)
		return risk.DefaultMode()
	}
	u.announceClear()
	return configured
}

// announce books the downgrade once per (configured mode) episode, so the line
// is a reading of 「这一档被降下来了」 rather than a per-call echo.
func (u *unsilencedMode) announce(configured risk.Mode) {
	u.mu.Lock()
	already := u.announcedSet && u.announced == configured
	u.announced, u.announcedSet = configured, true
	u.mu.Unlock()
	if already {
		return
	}
	u.logf("approval: MODE-UNANSWERABLE configured=%s effective=%s because=no-answer-source "+
		"(本进程此刻没有能答复的界面：档位不许把问题静默掉，L1 那一发会按拒绝处理)",
		configured, risk.DefaultMode())
}

// announceClear re-arms the line: once the host can answer again (or the
// configured mode is already the strictest one), a later downgrade is a new
// event and is booked again.
func (u *unsilencedMode) announceClear() {
	u.mu.Lock()
	u.announcedSet = false
	u.mu.Unlock()
}
