package main

// Ticket 273, wave 2: the fan-out seam the resident voice leg hands to the D43
// machine as Options.Sink.
//
// WHY A FAN-OUT AND NOT A SINK THAT JUST ACTS. The frozen table fires 50 effect
// names (62 delivery points, 33 rows), and they are not all the same kind of
// thing: some are OBSERVATIONS the owner wants recorded, and a few are ACTIONS
// that change the outside world. Wave 1 (effect_audit_273.go) booked one name
// into the rolling JSONL listener. That shape is right for a record and wrong
// for an action: a log line cannot push a panel frame. Conversely a sink that
// executed everything would execute the actions the resident voice leg already
// performs inline on the very paths that fire these effects. This file is the
// one place where a reader can see both halves at once, and where "audit only"
// is the DEFAULT rather than a promise in a comment.
//
// THE DRIVE SET IS THE FOUR NAMES WITH ZERO DRIVE IMPLEMENTATION TODAY. The
// resident voice leg already owns panel pushes, the panel keep-alive, the
// conversation ring and the mic gate; what it does not have is a recipient that
// turns the fired effect into that work. Those four are exactly:
//
//	panel.stream-push         D43 #15 (table.go:96-110, both guard branches)
//	panel.keep-alive          D43 #26 (table.go:185-189)
//	mic.enable                D43 #28 (table.go:209-213)
//	conversation.ring-solid   D43 #28 (table.go:209-213)
//
// Everything else - mic.keep-off and tts.stop first among them - is auditOnly:
// booked through audit and NEVER handed to drive, even when the caller injected
// one. R1 is the reason and it is a double-execution reason, not a style
// preference. mic.keep-off is row #26's inline branch in the resident leg: the
// same walk that fires the effect has already closed the mic. Driving it again
// would close it twice, and for tts.stop that second execution is audible. The
// default-deny spelling (map lookup, missing = false) is what makes a future
// frozen-table name audit-only until somebody edits this file on purpose.
//
// INJECTION, NOT GLOBALS. Both halves are func fields on the value the leg
// builds, so a test can count calls and production can pass the real panel/mic
// closures. A nil audit means "no recorder installed" and a nil drive means "no
// executor installed" - each branch is skipped on its own, which is what lets
// the leg ship the audit half before the drive half (or the reverse) without a
// second type. The drive set itself is the only package-level state here, and it
// is a literal read-only map: no mutex, no clock, no goroutine, no I/O.
//
// CONTRACT INHERITED FROM machine.go:31-34: a Sink must not block and must not
// re-enter Dispatch. This value adds neither a lock nor a queue, so it holds
// that by construction - each call is at most two function calls into code the
// caller supplied. audit runs FIRST, before drive, so a drive that fails still
// leaves the effect booked; ordering the other way would lose the record of the
// effect that the failure was about.

import "github.com/CarlosShao/wisp/internal/statemachine"

// effectFanout273 is one machine's recipient: audit for every fired effect,
// drive only for the allowlisted four. Both funcs are injected (nil = that
// half is not installed); the type itself carries no state.
type effectFanout273 struct {
	drive func(statemachine.Effect)
	audit func(statemachine.Effect)
}

// sink is the statemachine.Sink shape. It books the effect, then - only when
// the name is in drive273 and an executor was injected - lets the action half
// run. A name outside drive273 is never driven.
func (f effectFanout273) sink(e statemachine.Effect) {
	if f.audit != nil {
		f.audit(e)
	}
	if drive273[e.Name] && f.drive != nil {
		f.drive(e)
	}
}

// drive273 is the allowlist of effect names whose action half the resident
// voice leg may execute. Membership is deliberate and small: adding a name here
// is the same edit as promising that the inline path does NOT already do it,
// because that promise is what keeps R1 (no double execution) true.
//
// The four names are spelled exactly as the frozen D43 rows spell them
// (internal/statemachine/table.go); a typo here would silently demote a drive
// effect to auditOnly, so the test file pins both the set size and the fact
// that every key occurs in statemachine.Table.
var drive273 = map[string]bool{
	"panel.stream-push":       true,
	"panel.keep-alive":        true,
	"conversation.ring-solid": true,
	"mic.enable":              true,
}
