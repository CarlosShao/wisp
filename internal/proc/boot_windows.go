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
	Env      buildinfo.Env
	Layout   Layout
	Job      *JobScope
	Instance *SingleInstance
	Registry *observe.Registry
	StartedAt time.Time // wall clock, for boot records only
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
				slog.Info("activation requested by second launch (ball bring-to-front lands with ticket 07)")
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
func (rt *Runtime) Shutdown(fast bool) []StepRecord {
	hooks := ShutdownHooks{}
	if rt.Job != nil {
		hooks.CloseJob = func(context.Context) error { return rt.Job.Close() }
	}
	records := RunShutdownSequence(hooks, ShutdownOptions{Fast: fast})
	if rt.Instance != nil {
		_ = rt.Instance.Release()
	}
	return records
}
