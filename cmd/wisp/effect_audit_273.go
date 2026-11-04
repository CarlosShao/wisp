package main

// Ticket 273, wave 1: the D43 side-effect audit face of `wisp models ensure`.
//
// WHAT WAS MISSING, stated as a count rather than as an intention. The D43
// table has carried SideEffects on 39 of its rows since it was frozen, and
// statemachine.Machine.Dispatch has offered every one of them to Options.Sink
// since ticket 07 (internal/statemachine/machine.go:162-164). Across the whole
// non-test tree there was not ONE statemachine.Options with a Sink in it: four
// construction sites, four with Initial alone. The table was a promise with
// zero recipients, so "the machine fired model.verify-sha256-signature" was a
// statement no process could make after the fact.
//
// Wave 1 collects exactly ONE name - model.verify-sha256-signature, row #37 -
// for the one machine that fires it in production today, and it collects it
// into a channel that already exists: the rolling JSONL audit listener
// modelsEnsure installs at models.go:339 (<data root>\logs\wisp-*.jsonl, one
// rolling writer under observe.InitLog with redaction on). No new channel is
// created here, which is the whole point: a fourth listener is a fourth place
// nobody reads.
//
// WHAT THIS FILE DOES NOT DO, because each omission is a way the wiring could
// be a fake instead of a fact:
//
//   - It does not add, rename or reorder a side-effect name. The table, C17's
//     allowlist and the 50 names are frozen contract (ticket 273's forbidden
//     list); this file only formats what Dispatch already hands it.
//   - It does not EXECUTE the effect. The success branch of row #37 already
//     re-verified the installed tree before the dispatch - internal/models
//     bridge.go:58 runs Manager.VerifyInstalled on the way out, and Dispatch
//     only fires after that returned nil. This is a recorder, so the verdict
//     survives the terminal; it is not a second implementation of the check.
//   - It does not invent an effect for the failed branch. Row #37's
//     EvDownloadFailed row (internal/statemachine/table.go:267) carries no
//     SideEffects, so a failed walk fires nothing and this file books nothing
//     (ticket 273 R3). The failure verdict stays the stderr sentence at
//     models.go:372-375, unchanged by this ticket.
//
// THE CHANNEL, named precisely because "log it" has two spellings in this
// package and only one of them is right. The record goes through the logf
// closure handOffModel receives - sinkLogger.Info over the pipeline
// installLogSink started (models.go:344-346). It is deliberately NOT
// slog.Default(): installLogSink makes the default a tee over that same file
// AND stderr (logsink.go:153-161), so a record written that way would print on
// the owner's terminal a second time, and the same call shape would put a
// resident-process record into whichever listener installed last (ticket 273
// R2; the "audit trail already has its own console form" note at
// logsink.go:116-127 is the same reasoning one ticket earlier).
//
// WHY `at` IS A COUNTER AND NOT A CLOCK (D42#9): a wall-clock reading is not a
// measurement this repository accepts as a bound, and a timestamp would invite
// comparing records across processes when the only question a reader of one
// file has is "which of these fired first in THIS walk". The counter is per
// sink (per machine), starts at 1, and never decreases. Nothing else here reads
// a clock.

import (
	"sync/atomic"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// effectAuditRecordPrefix273 is the tag that makes an audit line identifiable
// without re-parsing every message in the file. modelsEnsure's logf prepends
// "models: " on its own, so a settled line reads
//
//	models: effect d43=37 name=model.verify-sha256-signature at=1
//
// and the fields after the tag are the three ticket 273 asks for, in the
// key=value spelling this package's other records already use
// (`models: hand-off OK id=... state=... root=...`). Structured slog attributes
// were not available on this seam: logf's shape is fixed by handOffModel's two
// production callers (models.go:348 and resident_voice_windows.go:2035), and
// widening it to a *slog.Logger is a signature change outside this ticket's
// write face.
const effectAuditRecordPrefix273 = "effect "

// effectAuditSink273 returns the statemachine.Sink that books one record per
// fired side effect through logf.
//
// The sink contract is Dispatch's, not this file's: it must not block and must
// not re-enter Dispatch (internal/statemachine/machine.go:31-34). One Info call
// on a writer the leg already writes to is the whole body, so it satisfies
// that by construction - there is nothing to post to and no lock to take.
//
// nil logf is the degraded branch: installLogSink refused, modelsEnsure has
// already said 持久日志未启用 on stderr, and the sink records nothing. Writing
// to stderr anyway would be a second channel invented for the failure case,
// which is exactly the shape ticket 273 forbids. The counter still advances, so
// `at` keeps meaning "how many effects this machine fired", not "how many
// happened to be written".
func effectAuditSink273(logf func(string, ...any)) statemachine.Sink {
	var at atomic.Uint64
	return func(e statemachine.Effect) {
		n := at.Add(1)
		if logf == nil {
			return
		}
		logf("%sd43=%d name=%s at=%d", effectAuditRecordPrefix273, e.D43, e.Name, n)
	}
}
