package watchdog

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TestThreeConsecutiveSettleFailuresReachWatchdogAlertWithRestart is AC#2:
// the fixture dispose never falls back under the Warm cap, three
// verifications fail in a row, the third raises exactly one fatal alert, and
// the host mapping (alert -> D43 row 35, EvWatchdogFatal) lands the machine
// in WatchdogAlert with the one-click restart effect fired. That mapping is
// the seam the resident wiring owns; this test pins both halves of it.
func TestThreeConsecutiveSettleFailuresReachWatchdogAlertWithRestart(t *testing.T) {
	// The D37 resource row the alert belongs to: one mapped state
	// (WatchdogAlert), retry never - the restart affordance is the remedy.
	if ms := observe.ClassResource.MappedStates(); len(ms) != 1 || ms[0] != "WatchdogAlert" {
		t.Fatalf("ClassResource.MappedStates() = %v, want [WatchdogAlert] (D37 resource row)", ms)
	}
	if observe.ClassResource.RetryPolicy() != observe.RetryNever {
		t.Errorf("ClassResource.RetryPolicy() = %v, want RetryNever", observe.ClassResource.RetryPolicy())
	}

	tree := &fakeTree{}
	tree.set(overCapTree(400 << 20)) // over the Warm cap for the whole window

	observe.ReleaseMemory() // the dispose path's instrumented release: the counter is > 0

	log, buf := captureLogger()
	als := &alertLog{}
	l := New(Options{
		TreeReader:     tree,
		Registry:       observe.NewRegistry(),
		Log:            log,
		TickEvery:      time.Millisecond,
		SettleWithin:   30 * time.Millisecond,
		SettleInterval: 5 * time.Millisecond,
	})

	// The host seam: AlertFatal maps to D43 row 35. The sink records every
	// fired effect so the restart affordance is asserted, not assumed.
	effects := make(chan statemachine.Effect, 16)
	sm := statemachine.New(statemachine.Options{
		Initial: statemachine.StateWarm,
		Sink:    func(e statemachine.Effect) { effects <- e },
	})
	defer sm.Close()

	ctx := context.Background()
	onAlert := func(a Alert) {
		als.add(a)
		if a.Kind == AlertFatal {
			if _, err := sm.Dispatch(statemachine.EvWatchdogFatal, nil); err != nil {
				t.Errorf("dispatching EvWatchdogFatal failed: %v", err)
			}
		}
	}
	l.onAlert = onAlert

	verify := func() *observe.SettleReport {
		t.Helper()
		pass, rep, err := l.VerifySettle(ctx, observe.SLOWarm, true, 400<<20)
		if err != nil {
			t.Fatalf("VerifySettle: %v", err)
		}
		if pass {
			t.Fatalf("VerifySettle passed while the fixture tree never fell under cap")
		}
		if rep == nil || rep.BackWithinCapMS != -1 {
			t.Fatalf("BackWithinCapMS = %+v, want -1 (never within the window)", rep)
		}
		return rep
	}

	verify() // failure 1
	if got := l.SettleFailures(); got != 1 {
		t.Fatalf("consecutive failures = %d, want 1", got)
	}
	if got := len(als.of(AlertFatal)); got != 0 {
		t.Fatalf("fatal alerts after 1 failure = %d, want 0", got)
	}
	verify() // failure 2
	if got := l.SettleFailures(); got != 2 {
		t.Fatalf("consecutive failures = %d, want 2", got)
	}
	verify() // failure 3 -> fatal
	if got := l.SettleFailures(); got != 0 {
		t.Fatalf("consecutive failures after the fatal = %d, want 0 (reset)", got)
	}
	fatals := als.of(AlertFatal)
	if len(fatals) != 1 {
		t.Fatalf("fatal alerts = %d, want exactly 1", len(fatals))
	}
	a := fatals[0]
	if a.Class != observe.ClassResource {
		t.Errorf("fatal alert class = %q, want resource (D37 row: mapped state WatchdogAlert, retry never)", a.Class)
	}
	if a.State != observe.SLOWarm || a.Metric != "settle_fall_back" {
		t.Errorf("fatal alert = %s/%s, want Warm/settle_fall_back", a.State, a.Metric)
	}
	if !strings.Contains(a.Detail, "restart") {
		t.Errorf("fatal alert detail must name the restart affordance, got %q", a.Detail)
	}

	if got := sm.State(); got != statemachine.StateWatchdogAlert {
		t.Fatalf("machine state = %s, want WatchdogAlert (D43 row 35)", got)
	}
	select {
	case e := <-effects:
		if e.D43 != 35 || e.Name != "ui.one-click-restart" {
			t.Errorf("fired effect = %+v, want D43 35 / ui.one-click-restart", e)
		}
	default:
		t.Fatal("no effect fired: the restart affordance did not happen")
	}

	// The consecutive rule: a pass resets the count, so two more failures
	// after a success must NOT raise a second fatal yet.
	tree.set(overCapTree(100 << 20)) // within the Warm cap now
	pass, _, err := l.VerifySettle(ctx, observe.SLOWarm, true, 400<<20)
	if err != nil || !pass {
		t.Fatalf("the success fixture failed: pass=%v err=%v", pass, err)
	}
	tree.set(overCapTree(400 << 20))
	verify()
	verify()
	if got := len(als.of(AlertFatal)); got != 1 {
		t.Fatalf("fatal alerts after fail/2 with a pass in between = %d, want still 1", got)
	}
	verify() // third consecutive -> the second fatal
	if got := len(als.of(AlertFatal)); got != 2 {
		t.Fatalf("fatal alerts = %d, want 2 (the count re-arms after a pass)", got)
	}
	if !strings.Contains(buf.String(), "watchdog: settle verification") {
		t.Errorf("every verification must be logged; log:\n%s", buf.String())
	}
}

// TestVerifySettleReleasedFalseFailsByContract pins the release half: a
// dispose that did not invoke observe.ReleaseMemory fails the verification
// even when the memory came back - the gate exists to prove the release
// path ran (D32 settle row / C11 step 4).
func TestVerifySettleReleasedFalseFailsByContract(t *testing.T) {
	tree := &fakeTree{}
	tree.set(overCapTree(100 << 20)) // within the Warm cap already
	l := New(Options{
		TreeReader:     tree,
		Registry:       observe.NewRegistry(),
		SettleWithin:   20 * time.Millisecond,
		SettleInterval: 5 * time.Millisecond,
	})
	pass, rep, err := l.VerifySettle(context.Background(), observe.SLOWarm, false, 400<<20)
	if err != nil {
		t.Fatalf("VerifySettle: %v", err)
	}
	if pass || rep == nil || rep.Pass {
		t.Fatalf("released=false must fail the check by contract, got pass=%v rep.Pass=%v", pass, rep != nil && rep.Pass)
	}
}

// TestVerifySettleRejectsUnknownTarget keeps the D32 vocabulary closed: an
// unknown settle target is a config-class error, never a guessed state row.
func TestVerifySettleRejectsUnknownTarget(t *testing.T) {
	l := New(Options{TreeReader: &fakeTree{}, Registry: observe.NewRegistry()})
	_, _, err := l.VerifySettle(context.Background(), observe.SLOState("Nonsense"), true, 0)
	if err == nil {
		t.Fatal("VerifySettle accepted an unknown target state")
	}
	if got := observe.ClassOfOrInternal(err); got != observe.ClassConfig {
		t.Errorf("error class = %s, want config", got)
	}
}
