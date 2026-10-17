//go:build windows

package proc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Boot assembles the resident runtime (ticket 03): env resolution -> env
// layout -> Job Object (C30) -> per-session single instance (D42#7) ->
// goroutine registry (D38b). Each stage doubles as the init-time self-check
// for the parts of the runtime that already exist.
//
// Boot returns ErrAlreadyRunning when another Wisp instance of this session
// owns the mutex; the caller must then SignalExistingInstance and exit.
type Runtime struct {
	Env       buildinfo.Env
	Layout    Layout
	Job       *JobScope
	Instance  *SingleInstance
	Registry  *observe.Registry
	StartedAt time.Time // wall clock, for boot records only

	// shutdownHooks is the registration seam of the D38(e) sequence (ticket 246
	// AC#4). Boot always fills it, so a runtime built through Boot can hand a
	// step to the module that owns it instead of leaving seven steps with no
	// producer. It carries no order: RunShutdownSequence still walks 1..10.
	shutdownHooks *ShutdownHookSet

	// activationHook is what one round of the activation pump runs when the
	// second-instance event fires (ticket 319 AC#2). Nil is the honest pre-319
	// shape - the signal is booked and cleared only - and that branch is kept
	// verbatim below, because a runtime nobody injected into must keep behaving
	// exactly as it did. See SetActivationHook.
	//
	// The mutex is not decoration: RunEventLoop reads the field on its own
	// goroutine while an assembly root (today: cmd/wisp's resident boot, before
	// the loop starts) writes it, and a later ticket that reads a second-launch
	// condition may write it from anywhere. The field is one function value;
	// the lock makes the handoff as race-free as the shape allows.
	activationMu   sync.Mutex
	activationHook func() string

	// stopReq is ticket 228 AC#11's stop-request bit (form 乙-1, ruling A493).
	// Boot always fills it, for the same reason shutdownHooks is always filled:
	// the thing that RECORDS a request is a callback the assembly root installs
	// long before the loop runs - the tray's Exit item is handed to ball.New by
	// startResidentBall, hundreds of milliseconds before RunEventLoop is called -
	// so a bit created inside the loop would drop every request that arrives in
	// that window. It is a latch and not a waiter: see stop_windows.go for why a
	// callback on the ui-sta thread may not wait on anything.
	//
	// It adds a TRIGGER and no STEP. RunShutdownSequence still walks 1..10
	// (shutdown.go), the ten-step chain still hangs off the resident leg
	// returning normally, and hookableRoster keeps its closed eight slots
	// (shutdown_hooks.go).
	stopReq *stopRequest
}

// BootOption adjusts Boot (tests, embedders).
type BootOption func(*bootConfig)

type bootConfig struct {
	layout   *Layout
	registry *observe.Registry
}

// WithLayout overrides the env layout (unique mutex names for tests; the
// portable-mode override is ticket 06 and will use the same seam).
func WithLayout(l Layout) BootOption {
	return func(c *bootConfig) { c.layout = &l }
}

// WithRegistry overrides the process registry (tests).
func WithRegistry(r *observe.Registry) BootOption {
	return func(c *bootConfig) { c.registry = r }
}

// Boot performs the init-time self-checks and assembles the runtime.
func Boot(env buildinfo.Env, opts ...BootOption) (*Runtime, error) {
	cfg := &bootConfig{}
	for _, o := range opts {
		o(cfg)
	}

	// Self-check: env value is one of the three frozen values.
	if env != buildinfo.EnvProd && env != buildinfo.EnvDev && env != buildinfo.EnvTest {
		return nil, fmt.Errorf("proc: unknown env %q", env)
	}

	rt := &Runtime{Env: env, StartedAt: observe.NowWallUTC()}
	rt.shutdownHooks = &ShutdownHookSet{}
	// Ticket 228 AC#11 ⓐ: the stop-request bit is built HERE, not inside the event
	// loop, because the callback that records into it is installed by the assembly
	// root before the loop exists. A runtime nobody booted keeps a nil bit and the
	// nil-safe answers in stop_windows.go, which is the honest "nothing can ask
	// this runtime to stop" shape rather than a panic.
	rt.stopReq = newStopRequest()
	if cfg.registry != nil {
		rt.Registry = cfg.registry
	} else {
		rt.Registry = observe.Default
	}

	// Self-check: layout resolution (per-env fork incl. the portable-mode
	// override; test env resolves a real data dir and registers no mutex).
	layout, err := DefaultLayout(env)
	if err != nil {
		return nil, err
	}
	rt.Layout = layout
	if cfg.layout != nil {
		rt.Layout = *cfg.layout
	}

	// Self-check: the Job Object opens and carries KILL_ON_JOB_CLOSE.
	job, err := OpenJobScope()
	if err != nil {
		return nil, err
	}
	rt.Job = job

	// Self-check: single instance acquisition (skipped when disabled).
	if rt.Layout.MutexEnabled {
		si, err := AcquireSingleInstance(rt.Layout.MutexName, rt.Layout.ActivateEventName)
		if err != nil {
			_ = job.Close()
			if errors.Is(err, ErrAlreadyRunning) {
				return nil, ErrAlreadyRunning
			}
			return nil, err
		}
		rt.Instance = si
	}

	// Registry sanity: nothing but the boot-time thread is running; the
	// roster check must not flag resident overflow.
	if rep := rt.Registry.RosterReport(); rep.ResidentOverBaseline {
		_ = job.Close()
		return nil, fmt.Errorf("proc: resident goroutine baseline exceeded at boot: %+v", rep)
	}
	return rt, nil
}

// RunEventLoop blocks until this runtime is asked to leave, and pumps the
// activation event without a dedicated goroutine (a short WaitForSingleObject per
// iteration keeps the loop responsive while the D38b resident roster stays at its
// frozen baseline). Returns the shutdown reason.
//
// THE COMPLETE ROSTER OF WHAT CAN END THIS LOOP, since ticket 228 AC#11 (the
// doc used to enumerate "Ctrl+C / console close / WM_ENDSESSION landing in ticket
// 43", which stopped being the whole list the moment the tray got an executor):
//
//   - os.Interrupt, which is what a console Ctrl+C becomes;
//   - SIGTERM, registered on the same signal.NotifyContext below;
//   - a stop request recorded on this runtime (Runtime.RequestStop) - ticket
//     228 AC#11's tray Exit item is its only production caller today, and the
//     reason the request is a bit rather than a call into Shutdown is in
//     stop_windows.go: the requester runs on the ui-sta thread and may wait on
//     nothing;
//   - console close and WM_ENDSESSION, which still land in ticket 43: this loop
//     has no handler for either, and saying so here is why AC#11's executor ends
//     the loop instead of calling os.Exit on the user's window.
//
// Every one of them leaves the same way: this function returns and the caller's
// deferred chain runs the frozen D38(e) order. Since ticket 319 the activation
// round is handleActivation - still in this loop, still no goroutine, still one
// 50 ms wait per tick.
func (rt *Runtime) RunEventLoop() string {
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	stopSignal := rt.StopSignal()

	for {
		rt.handleActivation()
		select {
		case <-sigCtx.Done():
			return "signal"
		case <-stopSignal:
			// Ticket 228 AC#11: the request is a latch, so a stop asked for before
			// this select existed is still honoured here. The return value names
			// who asked, because the line cmd/wisp prints off it is the only
			// sentence the trail has for "why did this process leave".
			return "stop request (" + rt.StopReason() + ")"
		case <-time.After(50 * time.Millisecond):
			// loop tick: activation pump + signal check + stop-request check
		}
	}
}

// SetActivationHook attaches this runtime's second-launch executor. It is the
// same injection shape RegisterShutdownHook set for the D38(e) sequence: the
// module that owns the capability hands the process a FUNCTION VALUE at assembly
// time, so internal/proc stays a leaf - it learns nothing about balls, panels,
// WebView2 or internal/ball, and no package-level dependency edge appears
// (`GOOS=windows go list -deps ./cmd/wisp` is unchanged by this field).
//
// What the hook must do, and what it may not: it runs inline on the event-loop
// goroutine for one round of the pump, so it has to return - a hook that blocks
// stalls the exit-signal check with it. It answers a receipt string, which
// handleActivation logs verbatim; it must not claim work it did not do, because
// that receipt is the only sentence the field has for "who was asked to raise
// what". A nil hook restores the honest pre-319 behaviour (book and clear only).
func (rt *Runtime) SetActivationHook(f func() string) {
	rt.activationMu.Lock()
	defer rt.activationMu.Unlock()
	rt.activationHook = f
}

// activationHookRead returns the registered hook, or nil.
func (rt *Runtime) activationHookRead() func() string {
	rt.activationMu.Lock()
	defer rt.activationMu.Unlock()
	return rt.activationHook
}

// handleActivation runs ONE round of the activation pump: wait a bounded 50 ms
// for the second-instance event, and when it fired, hand it to whatever this
// runtime was given and clear the event so the next launch can signal again. It
// is the testable half of what used to be an inline branch of RunEventLoop
// (ticket 319 AC#2's harness drives one round at a time), and it reports whether
// this round consumed a signal: false means the event was quiet.
//
// The pre-319 sentence is the no-hook branch below, and it is kept word for word
// except for the part that was describing a lack ("no path from this loop to a
// window"). That was true of this package and false of the process that hosts
// it, which had the path all along (the ball's own panel executor). What is
// still true is stated instead: this runtime was handed no hook.
func (rt *Runtime) handleActivation() bool {
	if rt.Instance == nil {
		return false
	}
	if code, err := windows.WaitForSingleObject(rt.Instance.ActivateEvent(), 50); err != nil || code != windows.WAIT_OBJECT_0 {
		return false
	}
	if hook := rt.activationHookRead(); hook != nil {
		receipt := hook()
		slog.Info("activation requested by second launch",
			"handled_here", "the registered activation hook ran",
			"receipt", receipt)
		_ = rt.Instance.ResetActivation()
		return true
	}
	slog.Info("activation requested by second launch",
		"handled_here", "booked and cleared only",
		"bring_to_front", "no activation hook is registered with this runtime (Runtime.SetActivationHook); the signal was booked and cleared only")
	_ = rt.Instance.ResetActivation()
	return true
}

// Shutdown runs the D38(e) sequence with the hooks this runtime owns today
// (Job close; single-instance release is an OS no-op at process death but is
// released explicitly after the sequence). Later tickets register their hooks
// on the same sequence - the order itself is frozen and audited.
//
// Since ticket 246 AC#4 that registration entry exists: RegisterShutdownHook
// attaches a step's owner, and the set is read here as a snapshot. Two facts
// did not move: RunShutdownSequence still walks steps 1..10 in the frozen
// order, and a step nobody registered is still recorded as skipped - registering
// is what changes the record, never the promise that the record is honest.
func (rt *Runtime) Shutdown(fast bool) []StepRecord {
	hooks := ShutdownHooks{}
	if rt.shutdownHooks != nil {
		hooks = rt.shutdownHooks.Hooks()
	}
	if hooks.CloseJob == nil && rt.Job != nil {
		hooks.CloseJob = func(context.Context) error { return rt.Job.Close() }
	}
	records := RunShutdownSequence(hooks, ShutdownOptions{Fast: fast})
	if rt.Instance != nil {
		_ = rt.Instance.Release()
	}
	return records
}

// RegisterShutdownHook attaches one hook to this runtime's D38(e) sequence.
// Ticket 246 AC#4's shape: the module that owns a step hands it to the process
// at assembly time, so the audit trail says what ran instead of saying
// "skipped" seven times.
//
// The errors it returns are the guarantees: a step that has no hook slot
// (8 and 10), a slot that already has an owner, a nil hook, and a sequence that
// has already started, are each refused. See shutdown_hooks.go for why each one
// is a refusal rather than a quiet overwrite.
func (rt *Runtime) RegisterShutdownHook(step ShutdownStep, hook func(ctx context.Context) error) error {
	if rt.shutdownHooks == nil {
		rt.shutdownHooks = &ShutdownHookSet{}
	}
	return rt.shutdownHooks.Register(step, hook)
}

// RegisteredShutdownSteps lists the steps this process has an owner for, in
// sequence order. The boot report reads it so "which steps really run" is a
// fact taken off the registration, not off a comment.
func (rt *Runtime) RegisteredShutdownSteps() []ShutdownStep {
	if rt.shutdownHooks == nil {
		return nil
	}
	return rt.shutdownHooks.Registered()
}
