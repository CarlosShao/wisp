//go:build windows

package main

// Ticket 384, the watchdog host wiring, driven end to end through the REAL
// components the resident boot assembles:
//
//   - a REAL proc.TreeSampler over the REAL C30 Job Object proc.Boot opens,
//     so the loop's ticks read this process's own tree metrics;
//   - a REAL D43 machine (statemachine.New) carrying a TEST sink - the shape
//     ticket 273's own census names as the legal test form, so the mapped
//     effects (kws.keep-alive-alert, ui.one-click-restart) have an observable
//     landing. Production builds the machine with Sink=nil (273 open); the
//     test sink is test-side only and invents nothing;
//   - the REAL host mapping functions sloStateFromD43Machine and
//     dispatchWatchdogAlert - the same functions startResidentWatchdog wires,
//     at test cadence (the Options seam the watchdog doc sanctions:
//     "the test seam shrinks them, production never does").
//
// The production assembly startResidentWatchdog itself is driven in
// Test384StartResidentWatchdogRegistersRealLoop: default cadence, cancelled
// before its first unload pass could ever fall due, read off the D38b roster.

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/watchdog"
)

// lineBuf384 collects the loop's own log lines (the tick reads, the unload
// passes, the mapping receipts). The handler is installed as slog's default
// BEFORE the loop under test is built, because watchdog.New captures
// slog.Default() at construction.
type lineBuf384 struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lineBuf384) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lineBuf384) count(needle string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Count(l.b.String(), needle)
}

// effectLog384 is the test sink: every Effect the machine fires lands here in
// dispatch order. Production carries Sink=nil (273 open); this is the
// test-side observer that makes "the mapping dispatched" a reading.
type effectLog384 struct {
	mu  sync.Mutex
	all []statemachine.Effect
}

func (e *effectLog384) sink(ev statemachine.Effect) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.all = append(e.all, ev)
}

func (e *effectLog384) countByD43(d43 int) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	n := 0
	for _, ev := range e.all {
		if ev.D43 == d43 {
			n++
		}
	}
	return n
}

func (e *effectLog384) has(d43 int, name string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, ev := range e.all {
		if ev.D43 == d43 && ev.Name == name {
			return true
		}
	}
	return false
}

// holdTouched384 allocates size bytes and touches every page, so the private
// working set of this process tree really carries them (the holdLeakFixture
// shape slo_windows.go uses - a real allocation, never a fake tree read).
func holdTouched384(size int) []byte {
	buf := make([]byte, size)
	for i := 0; i < len(buf); i += 4096 {
		buf[i] = 1
	}
	return buf
}

// waitUntil384 polls pred on a monotonic observe budget (D42#9), never a
// wall-clock difference.
func waitUntil384(t *testing.T, budget time.Duration, what string, pred func() bool) {
	t.Helper()
	tm := observe.NewTimeout(budget)
	for {
		if pred() {
			return
		}
		if tm.Expired() {
			t.Fatalf("384: never observed %s within %s", what, budget)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// currentTree384 is one real read of this process's tree.
func currentTree384(t *testing.T, rt *proc.Runtime) observe.TreeMetrics {
	t.Helper()
	m, err := proc.NewTreeSampler(rt.Job).ReadTree()
	if err != nil {
		t.Fatalf("384: real tree read failed: %v", err)
	}
	if m.PrivateWorkingSetBytes <= 0 {
		t.Fatalf("384: real tree read returned a zero private working set")
	}
	return m
}

// runWatchdogE2E384 boots the real runtime, builds a real machine with the
// test sink, and runs the loop at test cadence through the REAL host mapping
// functions. Every phase reads the machine, the sink and the log buffer.
func runWatchdogE2E384(t *testing.T) (*proc.Runtime, *residentVoice, *watchdog.Loop, *effectLog384, *lineBuf384) {
	t.Helper()

	buf := &lineBuf384{}
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(prev) })

	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("proc.Boot(test): %v", err)
	}
	t.Cleanup(func() { rt.Job.Close() })

	sinkLog := &effectLog384{}
	rv := &residentVoice{
		machine: statemachine.New(statemachine.Options{
			Initial: statemachine.StateSleeping,
			Sink:    sinkLog.sink,
		}),
	}

	loop := watchdog.New(watchdog.Options{
		TreeReader:     proc.NewTreeSampler(rt.Job),
		StateFn:        func() observe.SLOState { return sloStateFromD43Machine(rv) },
		Registry:       rt.Registry,
		OnAlert:        func(a watchdog.Alert) { dispatchWatchdogAlert(rv, a) },
		TickEvery:      20 * time.Millisecond,
		SustainTicks:   2,
		SettleWithin:   200 * time.Millisecond,
		SettleInterval: 25 * time.Millisecond,
	})
	root := observe.NewRootFrom(context.Background(), "watchdog-host")
	handle := loop.Start(root.Ctx)
	t.Cleanup(func() { root.Cancel(); <-handle.Done() })

	return rt, rv, loop, sinkLog, buf
}

// Test384ResidentWatchdogEndToEnd is AC#1/AC#2's machine: real ticks, a real
// unload pass over a real over-cap tree, the D43 #9 Armed self-loop (stay
// Armed + keep-alive effect), and the D43 #35 fatal row (WatchdogAlert +
// one-click restart) after three consecutive settle failures.
func Test384ResidentWatchdogEndToEnd(t *testing.T) {
	rt, rv, loop, sinkLog, buf := runWatchdogE2E384(t)

	// Phase 0 - the loop really ticks. The tick line is Debug-level, so it
	// only lands when this test's own handler is the default; counting real
	// tick lines is the "the goroutine is alive" reading.
	waitUntil384(t, 3*time.Second, "real watchdog tick lines", func() bool {
		return buf.count(`"msg":"watchdog: tick"`) >= 3
	})
	t.Logf("384-READ real tick lines=%d (the loop goroutine is alive and ticking)",
		buf.count(`"msg":"watchdog: tick"`))
	t.Logf("384-READ D38b roster CountByName(watchdog)=%d", rt.Registry.CountByName("watchdog"))
	if got := rt.Registry.CountByName("watchdog"); got != 1 {
		t.Fatalf("384: D38b roster CountByName(watchdog) = %d, want 1", got)
	}

	// Phase 1 - over cap in Sleeping (cap 25MB): hold 64MB of touched pages,
	// sustain 2 ticks, and the REAL unload pass runs over the whitelist
	// (nothing registered on this loop, so the pass walks the whitelist and
	// finds nothing to drop) and logs its completion. The tree stays over cap,
	// so the loop raises AlertOverrun for a state with no D43 row - the host
	// mapping records it as log/UI material and dispatches NOTHING (that is
	// the negative half of the mapping: only the two named rows dispatch).
	// The hold is dropped explicitly afterwards: Go keeps the pages until the
	// references die AND the instrumented release runs, and phase 2 needs the
	// tree genuinely back under the Armed cap before its own breach starts.
	hold64 := holdTouched384(64 << 20)
	waitUntil384(t, 5*time.Second, "the real unload pass over the sustained Sleeping breach",
		func() bool { return buf.count(`"msg":"watchdog: unload pass complete"`) >= 1 })
	t.Logf("384-READ Sleeping unload passes=%d, unload-action lines=%d, no-row alert receipts=%d",
		buf.count(`"msg":"watchdog: unload pass complete"`),
		buf.count(`"msg":"watchdog: unload action"`),
		buf.count(`"msg":"watchdog host: alert recorded (no D43 row; log/UI material)"`))
	if n := sinkLog.countByD43(9) + sinkLog.countByD43(35); n != 0 {
		t.Fatalf("384: Sleeping-phase over-cap must not dispatch D43 watchdog rows, got %d", n)
	}
	if rv.machine.State() != statemachine.StateSleeping {
		t.Fatalf("384: Sleeping-phase over-cap moved the machine to %s, want Sleeping", rv.machine.State())
	}
	hold64 = nil
	_ = hold64 // the reference is dead from here; the release below returns the pages
	settleBeforeSample(200)
	waitUntil384(t, 3*time.Second, "the tree back under the Armed cap (a real within-cap reading)",
		func() bool {
			return currentTree384(t, rt).PrivateWorkingSetBytes < 80<<20
		})

	// Phase 2 - the D43 #9 Armed self-loop. The loop has just observed
	// within-cap ticks (the sustain counter reset), so the Armed breach starts
	// a fresh episode: move the REAL machine to Armed the way the voice leg
	// does (EvKwsEnabled, row #5), hold 128MB over the Armed 90MB cap and wait
	// for the mapped keep-alive effect and the state that STAYS Armed.
	if _, err := rv.machine.Dispatch(statemachine.EvKwsEnabled, nil); err != nil {
		t.Fatalf("384: EvKwsEnabled (#5) refused from Sleeping: %v", err)
	}
	if rv.machine.State() != statemachine.StateArmed {
		t.Fatalf("384: machine state after #5 = %s, want Armed", rv.machine.State())
	}
	hold128 := holdTouched384(128 << 20)
	waitUntil384(t, 5*time.Second, "D43 #9 kws.keep-alive-alert effect (Armed over-cap, mapped)",
		func() bool { return sinkLog.has(9, "kws.keep-alive-alert") })
	t.Logf("384-READ D43#9 effects(kws.keep-alive-alert)=%d, machine stays %s",
		sinkLog.countByD43(9), rv.machine.State())
	if st := rv.machine.State(); st != statemachine.StateArmed {
		t.Fatalf("384: after D43 #9 the machine must STAY Armed, got %s", st)
	}
	if buf.count(`"msg":"watchdog host: D43 #9 Armed overrun mapped"`) < 1 {
		t.Fatalf("384: the #9 mapping receipt never landed in the log")
	}
	hold128 = nil
	_ = hold128
	settleBeforeSample(100)

	// Phase 3 - the D43 #35 fatal row. Three CONSECUTIVE settle verifications
	// fail by contract (released=false fails the release gate), the third
	// raises AlertFatal, the host mapping dispatches EvWatchdogFatal from
	// AnyState, the machine lands in WatchdogAlert and the one-click-restart
	// effect fires. VerifySettle runs on the caller's goroutine by design
	// (the dispose path's shape).
	peak := currentTree384(t, rt).PrivateWorkingSetBytes
	ctx := context.Background()
	var lastBackWithinCapMS int64
	for i := 1; i <= 3; i++ {
		pass, rep, err := loop.VerifySettle(ctx, observe.SLOSleeping, false, peak)
		if err != nil {
			t.Fatalf("384: VerifySettle %d errored: %v", i, err)
		}
		if pass || rep.Pass {
			t.Fatalf("384: VerifySettle %d with released=false must fail by contract", i)
		}
		lastBackWithinCapMS = rep.BackWithinCapMS
	}
	t.Logf("384-READ settle failures=3 (BackWithinCapMS=%d, released=false fails by contract); D43#35 effects(ui.one-click-restart)=%d, machine lands in %s",
		lastBackWithinCapMS, sinkLog.countByD43(35), rv.machine.State())
	waitUntil384(t, 3*time.Second, "D43 #35 ui.one-click-restart effect",
		func() bool { return sinkLog.has(35, "ui.one-click-restart") })
	if st := rv.machine.State(); st != statemachine.StateWatchdogAlert {
		t.Fatalf("384: after D43 #35 the machine state = %s, want WatchdogAlert", st)
	}
	if buf.count(`"msg":"watchdog host: D43 #35 fatal mapped (WatchdogAlert + one-click restart)"`) < 1 {
		t.Fatalf("384: the #35 mapping receipt never landed in the log")
	}
}

// Test384StartResidentWatchdogRegistersRealLoop drives the PRODUCTION
// assembly function: real runtime, real machine (Sink=nil, the 273-open
// shape), real panel host object. It reads the registration and the D38b
// roster - the "production starts the watchdog goroutine" fact - and cancels
// the root long before the default 60-tick sustain could ever fire.
func Test384StartResidentWatchdogRegistersRealLoop(t *testing.T) {
	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("proc.Boot(test): %v", err)
	}
	t.Cleanup(func() { rt.Job.Close() })

	rv := &residentVoice{machine: statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})}
	panel := &PanelManager{} // zero value: IsCreated/IsShown answer false under the zero mutex

	wd, wdRoot := startResidentWatchdog(rt, rv, panel)
	t.Cleanup(wdRoot.Cancel)

	t.Logf("384-READ production assembly registered=[%s], D38b CountByName(watchdog)=%d (273 open: production machine Sink=nil)",
		strings.Join(wd.Registered(), ","), rt.Registry.CountByName("watchdog"))
	if got := strings.Join(wd.Registered(), ","); got != watchdog.ComponentPanelWebView {
		t.Fatalf("384: registered components = %q, want exactly %q (the real panel host registers; "+
			"long-output-buffers/goja-vm/bm25-index-cache have no production object yet)", got, watchdog.ComponentPanelWebView)
	}
	if got := rt.Registry.CountByName("watchdog"); got != 1 {
		t.Fatalf("384: D38b roster CountByName(watchdog) = %d, want 1 (the loop goroutine is really running)", got)
	}

	// The panel component's no-window arm through the REAL registered object.
	dropped, detail := panelWebviewUnloadable{mgr: panel}.Unload(context.Background())
	if dropped {
		t.Fatalf("384: Unload on a never-created panel must not claim a drop (detail %q)", detail)
	}

	// The state mapping: Sleeping and Armed by name; a session-active state
	// (Listening, reached through #5 then #7 - the voice leg's own rows)
	// answers the zero SLOState (the fail-closed skip), never a smaller
	// state's cap.
	if got := sloStateFromD43Machine(rv); got != observe.SLOSleeping {
		t.Fatalf("384: Sleeping maps to %s, want Sleeping", got)
	}
	if _, err := rv.machine.Dispatch(statemachine.EvKwsEnabled, nil); err != nil {
		t.Fatalf("384: EvKwsEnabled (#5) refused from Sleeping: %v", err)
	}
	if got := sloStateFromD43Machine(rv); got != observe.SLOArmed {
		t.Fatalf("384: Armed maps to %s, want Armed", got)
	}
	if _, err := rv.machine.Dispatch(statemachine.EvWakeWord, nil); err != nil {
		t.Fatalf("384: EvWakeWord (#7) refused from Armed: %v", err)
	}
	if got := sloStateFromD43Machine(rv); got != observe.SLOState("") {
		// Listening is a session-active state with no D32 sampling row: the
		// mapping must answer the zero value, which the loop treats as
		// "not a D32 sampling state; tick skipped".
		t.Fatalf("384: Listening maps to %q, want the zero SLOState (fail-closed skip)", got)
	}
	if got := sloStateFromD43Machine(nil); got != observe.SLOSleeping {
		t.Fatalf("384: nil machine maps to %s, want Sleeping", got)
	}
}

// Test384PanelWebviewComponentArms is the whitelist component's contract on
// the real host object shape: nil host and windowless host both skip, and the
// shown arm skips because the user may be looking at it. The shown arm is
// driven through the zero manager's flags being false - a windowed Destroy
// reading is a live-window reading (winlive, not granted here).
func Test384PanelWebviewComponentArms(t *testing.T) {
	dropped, detail := panelWebviewUnloadable{mgr: nil}.Unload(context.Background())
	if dropped || detail != "panel host not assembled" {
		t.Fatalf("384: nil-manager arm = (%v, %q), want (false, panel host not assembled)", dropped, detail)
	}
	empty := &PanelManager{}
	dropped, detail = panelWebviewUnloadable{mgr: empty}.Unload(context.Background())
	if dropped || detail != "no panel window exists" {
		t.Fatalf("384: windowless arm = (%v, %q), want (false, no panel window exists)", dropped, detail)
	}
}
