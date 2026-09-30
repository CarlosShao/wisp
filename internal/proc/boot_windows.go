//go:build windows

package proc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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

// RunEventLoop is the empty event loop of ticket 03: it blocks until an exit
// signal (Ctrl+C / console close / WM_ENDSESSION landing in ticket 43) and
// pumps the activation event without a dedicated goroutine (a short
// WaitForSingleObject per iteration keeps the loop responsive while the D38b
// resident roster stays at its frozen baseline). Returns the shutdown reason.
func (rt *Runtime) RunEventLoop() string {
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		if rt.Instance != nil {
			if code, err := windows.WaitForSingleObject(rt.Instance.ActivateEvent(), 50); err == nil && code == windows.WAIT_OBJECT_0 {
				// What this loop does with an activation is all it does: book it
				// and clear it. It has no window to bring forward - the floating
				// ball belongs to whichever leg hosts it (cmd/wisp's resident
				// path does, through internal/ball), and no path connects this
				// event to that window in either direction today.
				slog.Info("activation requested by second launch",
					"handled_here", "booked and cleared only",
					"bring_to_front", "no path from this loop to a window; the ball, when a process has one, is its host's")
				_ = rt.Instance.ResetActivation()
			}
		}
		select {
		case <-sigCtx.Done():
			return "signal"
		case <-time.After(50 * time.Millisecond):
			// loop tick: activation pump + signal check
		}
	}
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
