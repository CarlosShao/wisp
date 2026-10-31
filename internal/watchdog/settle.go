package watchdog

import (
	"context"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// settleFatalConsecutive is the D43 row 35 trigger: three CONSECUTIVE failed
// settle verifications raise the fatal alert, which the host maps to
// EvWatchdogFatal (WatchdogAlert + one-click restart).
const settleFatalConsecutive = 3

// VerifySettle is the D32 settle-row re-measure the watchdog owns (ticket
// 42): after a Settling dispose it re-measures the tree for up to
// SettleWithin (default 10s, the D32 "back under the target cap within 10s"
// row) through observe.CheckSettle, and counts a failed verification.
//
// released states whether the dispose path invoked observe.ReleaseMemory;
// false fails the check by contract - the gate exists to prove the release
// path ran - and the FreeOSMemory counter (ticket 08) must be > 0. The
// instrumented wrapper is the evidence: a dispose calling
// debug.FreeOSMemory directly is invisible to it, which is the point.
//
// Three consecutive failures raise exactly one AlertFatal and reset the
// count; three MORE failures raise it again. The D43 row 35 machine row is
// from AnyState, so re-dispatching the mapped event is legal and the host
// can treat the alert as idempotent.
//
// This call blocks for up to SettleWithin - it runs on the caller's
// goroutine (the dispose path), never on the loop tick, so the config poll
// keeps running while a verification measures. The trade is deliberate: a
// dispose happens at most once per settle, and the mtime+size fingerprint
// loses nothing while the poll waits.
func (l *Loop) VerifySettle(ctx context.Context, target observe.SLOState, released bool, peak int64) (bool, *observe.SettleReport, error) {
	if !target.Valid() {
		return false, nil, observe.New(observe.ClassConfig,
			"watchdog: unknown settle target state "+string(target))
	}
	s := observe.NewSampler(l.tree, l.reg)
	rep, err := s.CheckSettle(ctx, target, l.settleWithin, l.settleInterval, released, observe.TreeMetrics{}, peak)
	if err != nil {
		return false, nil, err
	}

	l.mu.Lock()
	l.settleFailures++
	n := l.settleFailures
	fatal := !rep.Pass && n >= settleFatalConsecutive
	if rep.Pass || fatal {
		l.settleFailures = 0
	}
	l.mu.Unlock()

	l.log.Info("watchdog: settle verification",
		"target_state", string(target),
		"pass", rep.Pass,
		"back_within_cap_ms", rep.BackWithinCapMS,
		"final_bytes", rep.FinalBytes,
		"cap_bytes", rep.CapBytes,
		"free_os_memory_count", rep.FreeOSMemoryCount,
		"sample_errors", rep.SampleErrors,
		"consecutive_failures", n,
		"error_class", string(observe.ClassResource))

	if fatal {
		l.emit(Alert{
			Kind: AlertFatal, Class: observe.ClassResource, State: target,
			Metric:   "settle_fall_back",
			Measured: fmt.Sprintf("%d consecutive failures", settleFatalConsecutive),
			Limit:    fmt.Sprintf("back under cap within %s with the FreeOSMemory counter > 0", l.settleWithin),
			Detail:   "three consecutive settle verifications failed; the one-click restart affordance is required (D43 row 35)",
			At:       observe.WallTimestampUTC(observe.NowWallUTC()),
		})
	}
	return rep.Pass, rep, nil
}
