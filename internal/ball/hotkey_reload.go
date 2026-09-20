//go:build windows

package ball

// Hotkey reload bridge (ticket 64 A1): the thing that makes "the user edited
// [hotkey] in config.toml" become "the new key is live and the old one is
// dead" on a running instance.
//
// Before this file the chain did not exist: Ball.RebindHotkeys had zero
// callers in the whole repo, so a hotkey edit changed nothing until restart.
// The bridge is deliberately free of any config import - internal/ball owns
// no business/config knowledge (SPEC-01 §3) - so the host supplies two
// closures and one optional refresh hook:
//
//	m, _ := config.NewManager(path, res)
//	r := ball.NewHotkeyReloader(b, b.HotkeyConfig(), func() ball.HotkeyConfig {
//		h := m.Config().Hotkey
//		return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{
//			Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel})
//	})
//	r.Refresh = func() error { _, err := m.CheckAndReload(); return err }
//	go r.Run(ctx, time.Second) // panic-isolated inside; the host owns the goroutine
//
// or, on a host that already polls config every tick (the wisp watchdog),
// call r.Check() from that tick.
// Why the bridge polls the section instead of riding config.Manager.OnReload:
// [hotkey] is HOT-tier (D36), and OnReload only fires for RELOAD-tier
// sections - a hotkey edit would never reach a callback installed there. The
// host's poll (CheckAndReload, driven by the watchdog tick) is the trigger on
// a running instance; OnReload is still honored via Sections() for hosts that
// prefer to be told, and both paths converge on the same diff, so a double
// notification costs nothing.

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"time"
)

// HotkeyBinder applies a new binding set to the live registrations. *Ball
// implements it; the seam exists so the diff logic is testable with no window.
type HotkeyBinder interface {
	RebindHotkeys(cfg HotkeyConfig) HotkeyReport
}

// HotkeySource returns the currently effective bindings (the host's mapping
// of [hotkey], already defaulted via ApplyHotkeyDefaults).
type HotkeySource func() HotkeyConfig

// HotkeyReloader diffs the config-sourced bindings against what the ball is
// holding right now and re-registers on a real change.
type HotkeyReloader struct {
	binder HotkeyBinder
	src    HotkeySource

	// Refresh runs before every Check (hosts hook config.Manager.CheckAndReload
	// here so the poll drives BOTH the config reload and the rebind). nil =
	// the host refreshes the config itself.
	Refresh func() error

	mu      sync.Mutex
	applied HotkeyConfig
	exists  bool
	last    HotkeyReport
	rebound int
}

// NewHotkeyReloader wires a ball to a hotkey source. The applied set starts
// as the ball's CURRENT bindings, so constructing the bridge never re-registers
// anything by itself.
func NewHotkeyReloader(binder HotkeyBinder, current HotkeyConfig, src HotkeySource) *HotkeyReloader {
	return &HotkeyReloader{binder: binder, src: src, applied: current, exists: true}
}

// Check pulls the effective bindings and re-registers when they differ from
// what is live. It returns true when a rebind happened. Safe from any
// goroutine EXCEPT the ball's UI thread (the rebind is synchronous on it, so
// calling from there would deadlock) - hosts call it from their own tick.
func (r *HotkeyReloader) Check() (bool, HotkeyReport) {
	if r.Refresh != nil {
		if err := r.Refresh(); err != nil {
			slog.Warn("ball: hotkey bridge config refresh failed; keeping current hotkeys", "err", err)
		}
	}
	if r.src == nil {
		return false, r.Report()
	}
	cfg := r.src()
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.exists && reflect.DeepEqual(r.applied, cfg) {
		return false, r.last // nothing changed: do not touch Win32
	}
	rep := r.binder.RebindHotkeys(cfg)
	r.applied, r.exists, r.last, r.rebound = cfg, true, rep, r.rebound+1
	for _, line := range rep.Problems() {
		slog.Error("ball: hotkey reload", "binding", line)
	}
	slog.Info("ball: hotkeys rebound after config change",
		"summon", cfg.Summon, "mute", cfg.Mute, "cancel", cfg.Cancel, "panel", cfg.Panel,
		"live", len(rep.Live()))
	return true, rep
}

// Sections is the config.Manager.OnReload-shaped hook: it reacts when the
// reload-tier notification names [hotkey] (it does not today, hence the poll)
// or when the host wants a push instead of a pull.
func (r *HotkeyReloader) Sections(sections []string) {
	for _, s := range sections {
		if s == "hotkey" {
			r.Check()
			return
		}
	}
}

// OnReload returns the hook installed as config.Manager.OnReload.
func (r *HotkeyReloader) OnReload() func(sections []string) { return r.Sections }

// Run is the poll loop for hosts without a watchdog tick. Every tick is
// panic-isolated (one bad config read must not kill the app's hotkeys), and
// the loop stops when ctx is cancelled. Callers own the goroutine:
// `go r.Run(ctx, time.Second)`.
func (r *HotkeyReloader) Run(ctx context.Context, every time.Duration) {
	if every <= 0 {
		every = time.Second
	}
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			r.checkSafe()
		}
	}
}

// checkSafe keeps a panic in a rebind from taking the process down: the
// bridge runs on the host's goroutine, whose owner is this reloader.
func (r *HotkeyReloader) checkSafe() {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("ball: hotkey reload panicked; hotkeys unchanged", "panic", rec)
		}
	}()
	r.Check()
}

// Applied returns the binding set currently live on the ball.
func (r *HotkeyReloader) Applied() HotkeyConfig {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.applied
}

// Report returns the last registration outcome.
func (r *HotkeyReloader) Report() HotkeyReport {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.last
}

// Rebinds counts the actual Win32 re-registrations (a test seam for "the
// bridge must not churn the hotkeys on every tick").
func (r *HotkeyReloader) Rebinds() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rebound
}
