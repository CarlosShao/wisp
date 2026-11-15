//go:build windows

package main

// Ticket 46's settle-gate cases (the S6 acceptance's AC#2 "settle gate has no
// production caller" finding, and 46-s6-acceptance.md §1 AC#2's third row).
//
// WHAT THESE CASES ARE FOR, in one sentence: the D32 settle row must be
// verifiable by DRIVING THE PRODUCT'S OWN DISPOSE CASCADE, not by reading a
// source string. Each case below therefore starts at a real
// session.Scope dispose through voiceResidency (the exact hook
// resident_voice_engines_windows.go's scopeOptions installs), and asserts on
// readings the loop and the instrumented release counter produce - so deleting
// the wiring turns these red for the reason the wiring exists, not because a
// comment moved.
//
// Staging: no sherpa DLL, no model, no audio device. The dispose cascade's
// engine step is skipped because this residency never acquired an engine (the
// cascade's "no slot" arm), which is the same arm a process whose engines were
// already released walks.

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/watchdog"
)

// settleTree46 is the tree reader the fixture owns: one footprint number, and a
// read counter so a case can tell "the peak was read" from "a zero was passed".
type settleTree46 struct {
	mu    sync.Mutex
	bytes int64
	err   error
	reads int
}

func (s *settleTree46) ReadTree() (observe.TreeMetrics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reads++
	if s.err != nil {
		return observe.TreeMetrics{}, s.err
	}
	return observe.TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: s.bytes}, nil
}

func (s *settleTree46) readCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.reads
}

// settleLoop46 builds a real watchdog loop over the fixture tree with a short
// settle window, which is the only thing a case shortens: the window is the
// loop's own option (D32's 10s row in production), never a number this test
// re-states as the gate's.
func settleLoop46(t *testing.T, tree *settleTree46, reg *observe.Registry, within time.Duration) *watchdog.Loop {
	t.Helper()
	return watchdog.New(watchdog.Options{
		TreeReader:     tree,
		Registry:       reg,
		SettleWithin:   within,
		SettleInterval: 20 * time.Millisecond,
		StateFn:        func() observe.SLOState { return observe.SLOSleeping },
	})
}

func waitFor46(t *testing.T, limit time.Duration, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(limit)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("46: %s never happened within %v", what, limit)
}

// Test46SettleGateVerifiesTheDisposeCascade is the load-bearing case: a real
// session dispose (the cascade's four frozen steps) must produce one settle
// verification, judged by the loop, with the state the machine is in, the peak
// the tree really held, and a release the instrumented wrapper counted.
//
// Red under each of the three ways this could regress:
//   - no gate armed on the residency (the `_ = wd` shape): Runs stays 0;
//   - the release going back to debug.FreeOSMemory directly: the counter delta
//     is 0, so `released` would be false and the row could never pass;
//   - the peak not being read: Peak is 0 instead of the fixture's footprint.
func Test46SettleGateVerifiesTheDisposeCascade(t *testing.T) {
	const footprint = int64(400 << 20)
	before := observe.FreeOSMemoryCount()

	tree := &settleTree46{bytes: footprint}
	reg := observe.NewRegistry()
	loop := settleLoop46(t, tree, reg, 120*time.Millisecond)
	root := observe.NewRoot("46-settle-gate")
	t.Cleanup(root.Cancel)

	gate := newResidentSettleGate(loop, tree, reg, root,
		func() observe.SLOState { return observe.SLOSleeping })
	if gate == nil {
		t.Fatal("46: newResidentSettleGate refused a loop plus a tree reader; the production assembly passes both")
	}
	rv := newVoiceResidency(voicePlan{}, nil, nil)
	rv.attachSettleGate(gate)

	// The production wake path opens the scope, the production dispose path
	// ends it: nothing here calls the gate directly.
	scope, err := rv.openScope("46 case")
	if err != nil {
		t.Fatalf("46: openScope: %v", err)
	}
	if err := scope.Dispose("46 case dispose"); err != nil {
		t.Fatalf("46: dispose cascade: %v", err)
	}

	waitFor46(t, 10*time.Second, "the settle verification the cascade should have started",
		func() bool { return gate.snapshot().Runs == 1 })

	snap := gate.snapshot()
	if snap.Requests != 1 {
		t.Fatalf("46: the cascade's release step raised %d settle request(s), want exactly 1: the dispose path is not triggering the gate", snap.Requests)
	}
	if snap.Peak != footprint {
		t.Fatalf("46: the report's peak = %d, want the %d bytes the tree really held before the release (a peak the gate did not read is a guess)", snap.Peak, footprint)
	}
	if snap.Target != observe.SLOSleeping {
		t.Fatalf("46: the verification judged state %q, want Sleeping (what the D43 machine answered)", string(snap.Target))
	}
	if snap.Pass {
		t.Fatalf("46: a 400MB tree passed a 25MB Sleeping settle row: the gate is not judging the frozen cap")
	}
	if got := loop.SettleFailures(); got != 1 {
		t.Fatalf("46: the watchdog's consecutive settle-failure count = %d, want 1: the host is not calling the loop's own VerifySettle", got)
	}
	// The release half: the cascade's third step must move the INSTRUMENTED
	// counter (observe.FreeOSMemoryCount), which is the D32 row's own evidence.
	if got := observe.FreeOSMemoryCount(); got <= before {
		t.Fatalf("46: the dispose cascade released through something the instrument cannot see (count %d -> %d): the settle gate would fail by contract on every settle of a healthy process", before, got)
	}
	if tree.readCount() < 1 {
		t.Fatal("46: the peak read never touched the tree reader")
	}
	t.Logf("46-READ cascade -> %s", gate)
}

// Test46SettleGateJudgesTheStateTheMachineIsIn pins the second half of honesty:
// the target is READ, not hardcoded. An Armed process (KWS resident) is judged
// against the Armed cap, and a machine in a state with no D32 sampling row of
// its own falls back to Sleeping - saying so instead of silently substituting.
func Test46SettleGateJudgesTheStateTheMachineIsIn(t *testing.T) {
	reg := observe.NewRegistry()
	root := observe.NewRoot("46-settle-target")
	t.Cleanup(root.Cancel)

	cases := []struct {
		name  string
		state observe.SLOState
		want  observe.SLOState
	}{
		{"armed machine", observe.SLOArmed, observe.SLOArmed},
		{"no sampling row", observe.SLOState(""), observe.SLOSleeping},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tree := &settleTree46{bytes: 40 << 20}
			loop := settleLoop46(t, tree, reg, 60*time.Millisecond)
			gate := newResidentSettleGate(loop, tree, reg, root, func() observe.SLOState { return c.state })
			gate.request(settleRequest{peak: 40 << 20, released: true})
			waitFor46(t, 10*time.Second, "the verification to finish", func() bool { return gate.snapshot().Runs == 1 })
			if got := gate.snapshot().Target; got != c.want {
				t.Fatalf("46: machine state %q was judged as %q, want %q", string(c.state), string(got), string(c.want))
			}
		})
	}
}

// Test46WakeCancelsTheSettleVerification is the "do not walk D43 #35 on a
// healthy process" half: a wake that reopens the session ends the in-flight
// verification (its tree is no longer the disposed one), the loop counts no
// failure, and the cancellation is booked as a cancellation.
func Test46WakeCancelsTheSettleVerification(t *testing.T) {
	tree := &settleTree46{bytes: 400 << 20}
	reg := observe.NewRegistry()
	loop := settleLoop46(t, tree, reg, 2*time.Second) // long enough that the cancel lands inside it
	root := observe.NewRoot("46-settle-cancel")
	t.Cleanup(root.Cancel)

	gate := newResidentSettleGate(loop, tree, reg, root,
		func() observe.SLOState { return observe.SLOSleeping })
	rv := newVoiceResidency(voicePlan{}, nil, nil)
	rv.attachSettleGate(gate)

	if _, err := rv.openScope("46 cancel case"); err != nil {
		t.Fatalf("46: first openScope: %v", err)
	}
	gate.request(settleRequest{peak: 400 << 20, released: true})

	// The wake: a second epoch opens while the verification is inside its
	// window. This is the production line (openScope), not gate.abort.
	if _, err := rv.openScope("46 wake during settle"); err != nil {
		t.Fatalf("46: second openScope: %v", err)
	}
	waitFor46(t, 10*time.Second, "the cancelled verification to return",
		func() bool { return gate.snapshot().Runs == 1 })

	snap := gate.snapshot()
	if snap.Cancelled != 1 {
		t.Fatalf("46: cancelled = %d, want 1: the wake did not end the in-flight verification", snap.Cancelled)
	}
	if snap.Failures != 0 {
		t.Fatalf("46: a cancelled verification was booked as %d settle failure(s); a wake would then walk D43 #35 on a healthy process", snap.Failures)
	}
	if got := loop.SettleFailures(); got != 0 {
		t.Fatalf("46: the watchdog's own counter = %d after a cancelled verification, want 0 (watchdog/settle.go returns before it counts)", got)
	}
}

// Test46SettleGateLineNamesItsTrigger keeps the boot row honest without reading
// source: the row this assembly prints must name the trigger and the evidence,
// and an unassembled process must print the refusal instead of an empty line.
func Test46SettleGateLineNamesItsTrigger(t *testing.T) {
	if got := (*residentSettleGate)(nil).line(); got == "" {
		t.Fatal("46: an unassembled process printed an empty settle-gate row")
	} else if want := "未装配"; !strings.Contains(got, want) {
		t.Fatalf("46: an unassembled process's row = %q, want it to say %q", got, want)
	}
	tree := &settleTree46{bytes: 30 << 20}
	reg := observe.NewRegistry()
	root := observe.NewRoot("46-settle-line")
	t.Cleanup(root.Cancel)
	gate := newResidentSettleGate(settleLoop46(t, tree, reg, 50*time.Millisecond), tree, reg, root, nil)
	// 46 post-audit (three-state boot row): an assembly that never ran the
	// mount half is NOT the assembled posture any more - it is
	// assembled-but-unmounted, which is exactly the state the audit found the
	// old two-state line lying about. Reach the assembled posture the way the
	// boot does: through the production attach half onto a real residency
	// (the mount result is the row's own mount reading).
	rv := &residentVoice{residency: newVoiceResidency(voicePlan{}, newVoiceOwnerQueue(), nil)}
	if !rv.attachSettleGate(gate) {
		t.Fatal("46: the production attach half refused to mount the gate on a real residency")
	}
	line := gate.line()
	for _, want := range []string{"FreeOSMemory", "Sleeping", "SettleWithin", "D43#35"} {
		if !strings.Contains(line, want) {
			t.Fatalf("46: the assembled row = %q, want it to name %q", line, want)
		}
	}
}
