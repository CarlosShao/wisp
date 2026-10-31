package watchdog

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The resource watchdog loop (ticket 42; SPEC-01 §3, D32 16.3.2 revision 3,
// D42#10, D43 rows 9/35, SPEC-03 §4.3).
//
// One loop, two duties. The same tick that enforces the resource SLO also
// carries the config-file mtime+size poll (SPEC-03 §4.3): no second
// goroutine, no fsnotify dependency, and config changes keep flowing through
// config.Manager.CheckAndReload - the contract ticket 05 froze (idempotent
// poll, D36 tier semantics, OnReload / OnRestartPending events).
//
// Enforcement shape (D32 16.3.2 revision 3):
//   - thresholds come from the CURRENT state's frozen D32 row, read through
//     the observe accessors (internal/observe/caps.go) - never a second
//     table here;
//   - memory over-cap sustained for SustainTicks consecutive ticks (default
//     60 = 60s at the 1s tick) runs ONE unload pass over the safe whitelist
//     intersected with what is registered; a pass that leaves the tree over
//     cap raises an overrun alert - within cap or alert, never silent;
//   - in Armed the pass may still drop whitelisted components, but KWS is
//     structurally unreachable (it is not on the whitelist) and the
//     keep-alive warning with the disable-KWS suggestion goes out with the
//     sustained breach (D43 row 9);
//   - leak indicators (handles, GDI objects, and goroutines where the frozen
//     table carries a state cap) raise a leak alert with a diagnostic
//     snapshot and never auto-unload anything (D42#10);
//   - USER objects have no frozen absolute limit (thresholds.go records
//     them), so they are logged with the tick metrics, never gated;
//   - settle verification (VerifySettle) re-measures a disposed tree for up
//     to 10s and counts failures; three consecutive failures raise the fatal
//     alert the host maps to D43 row 35 (WatchdogAlert + one-click restart).
//
// The loop never touches the state machine (package doc: no state
// transitions): it raises alerts; the host owns the D43 mapping.

// AlertKind classifies a watchdog alert for the host-side D43 mapping.
type AlertKind string

const (
	// AlertOverrun: memory over cap sustained and the unload pass ran (or,
	// in Armed, the breach itself). In Armed the host maps this to D43 row 9
	// (EvWatchdogOverrun: stay Armed, keep-alive alert); other states have
	// no D43 row, so there the alert is log/UI material only.
	AlertOverrun AlertKind = "overrun"
	// AlertLeak: a leak indicator crossed its frozen cap (D42#10). No D43
	// row; the alert carries the diagnostic snapshot.
	AlertLeak AlertKind = "leak"
	// AlertFatal: three consecutive settle verifications failed. The host
	// maps this to D43 row 35 (EvWatchdogFatal -> WatchdogAlert, whose side
	// effect is the one-click restart affordance).
	AlertFatal AlertKind = "fatal"
)

// Alert is one watchdog verdict handed to Options.OnAlert. The host owns
// the D43 event mapping; the sink must be safe for concurrent calls (the
// settle verification emits from the caller's goroutine) and must not block
// or re-enter the loop.
type Alert struct {
	Kind AlertKind `json:"kind"`
	// Class is the D37 class every watchdog alert belongs to: resource
	// (mapped state WatchdogAlert, retry never - the restart affordance is
	// the remedy). The stable copy key is Class.MessageKey(); the visible
	// text is the deferred error-copy layer's decision.
	Class    observe.ErrorClass `json:"class"`
	State    observe.SLOState   `json:"state"`
	Metric   string             `json:"metric"`
	Measured string             `json:"measured"`
	Limit    string             `json:"limit"`
	// Detail carries the actionable sentence (e.g. the Armed disable-KWS
	// suggestion). User-visible copy per D37 is the deferred error-copy
	// layer's decision (observe.ClassResource.MessageKey), so only stable
	// wording travels here.
	Detail string `json:"detail,omitempty"`
	// Snapshot is the diagnostic snapshot JSON attached to leak alerts
	// (D42#10: alert + snapshot, never auto-unload). nil otherwise.
	Snapshot []byte `json:"snapshot,omitempty"`
	// At is the wall-clock RFC3339 stamp for the persisted record only
	// (observe clock discipline: never used for interval math).
	At string `json:"at"`
}

// Unloadable is one component the watchdog may drop on an unload pass.
type Unloadable interface {
	// Unload attempts to drop the component. dropped=false means nothing
	// was dropped this time (component busy, user viewing the panel, ...);
	// detail says why. A panic is recovered per call and logged - one bad
	// component never kills the loop.
	Unload(ctx context.Context) (dropped bool, detail string)
}

// The frozen safe-unloadable whitelist (D32 16.3.2 revision 3, verbatim):
// the panel WebView if the user is not viewing it, long-output temp buffers,
// the goja VM if enabled, the BM25 index cache. The unload pass walks THIS
// list and looks each name up in the registered set, so a component
// registered under any other name - KWS, audio capture, the SQLite
// db-writer - is structurally unreachable by every pass (ticket 42:
// "no unload of non-whitelisted components ever"; the AC#4 log audit pins it).
const (
	ComponentPanelWebView      = "panel-webview"
	ComponentLongOutputBuffers = "long-output-buffers"
	ComponentGojaVM            = "goja-vm"
	ComponentBM25IndexCache    = "bm25-index-cache"
)

// whitelist is the frozen D32 revision 3 list, in table order.
var whitelist = []string{
	ComponentPanelWebView,
	ComponentLongOutputBuffers,
	ComponentGojaVM,
	ComponentBM25IndexCache,
}

// Whitelist returns a copy of the frozen safe-unloadable whitelist.
func Whitelist() []string { return slices.Clone(whitelist) }

// Options configures the watchdog Loop. Zero values fall back to the frozen
// cadences; the test seam shrinks them, production never does.
type Options struct {
	// TreeReader is the process-tree metric source (C30 JobScope on Windows;
	// injected fakes in tests). Required - New panics without it, because a
	// watchdog with no eyes must not silently pretend to watch.
	TreeReader observe.TreeReader
	// StateFn returns the CURRENT state; thresholds are looked up per tick
	// (D32: the current state's row). nil = always Sleeping.
	StateFn func() observe.SLOState
	// Registry counts goroutines (the D38b registry is the product source of
	// truth) and receives the loop's own named goroutine. nil = observe.Default.
	Registry *observe.Registry
	// Config, when non-nil, is polled every tick (SPEC-03 §4.3: one loop,
	// both duties). nil = resource duty only.
	Config *config.Manager
	// OnAlert receives alerts; the D43 mapping is the host's. nil = log only.
	OnAlert func(Alert)
	// Log; nil = slog.Default().
	Log *slog.Logger
	// Owner recorded on the loop's named goroutine. Default "watchdog".
	Owner string
	// TickEvery is the loop cadence. Default 1s (SPEC-03 §4.3).
	TickEvery time.Duration
	// SustainTicks is how many consecutive over-cap ticks trigger one unload
	// pass. Default 60 (D32 revision 3: sustained 60s at the 1s tick).
	SustainTicks int
	// SettleWithin bounds one settle verification. Default 10s (D32 settle
	// row: back under the target cap within 10s).
	SettleWithin time.Duration
	// SettleInterval is the re-measure cadence inside one verification.
	// Default 250ms (the observe.CheckSettle default).
	SettleInterval time.Duration
}

// Loop is the resource watchdog. Start it once via Start (which spawns the
// named "watchdog" goroutine); give it the current state through
// Options.StateFn; components register with Register; alerts leave through
// Options.OnAlert. overTicks / overAlerted / armedWarned / leakLatched are
// touched only on the loop goroutine; comps and settleFailures are shared
// and mu-guarded.
type Loop struct {
	tree    observe.TreeReader
	stateFn func() observe.SLOState
	reg     *observe.Registry
	cfg     *config.Manager
	onAlert func(Alert)
	log     *slog.Logger
	owner   string

	tickEvery      time.Duration
	sustainTicks   int
	settleWithin   time.Duration
	settleInterval time.Duration

	mu             sync.Mutex
	comps          map[string]Unloadable
	settleFailures int

	// loop-goroutine-only enforcement state (see the Loop doc).
	overTicks   int
	overAlerted bool
	armedWarned bool
	leakLatched map[string]bool
}

// New builds a Loop from opts (see Options for the defaults).
func New(opts Options) *Loop {
	if opts.TreeReader == nil {
		panic("watchdog: Options.TreeReader is required")
	}
	reg := opts.Registry
	if reg == nil {
		reg = observe.Default
	}
	log := opts.Log
	if log == nil {
		log = slog.Default()
	}
	owner := opts.Owner
	if owner == "" {
		owner = "watchdog"
	}
	stateFn := opts.StateFn
	if stateFn == nil {
		stateFn = func() observe.SLOState { return observe.SLOSleeping }
	}
	tickEvery := opts.TickEvery
	if tickEvery <= 0 {
		tickEvery = time.Second
	}
	sustain := opts.SustainTicks
	if sustain <= 0 {
		sustain = 60
	}
	within := opts.SettleWithin
	if within <= 0 {
		within = 10 * time.Second
	}
	interval := opts.SettleInterval
	if interval <= 0 {
		interval = 250 * time.Millisecond
	}
	return &Loop{
		tree: opts.TreeReader, stateFn: stateFn, reg: reg, cfg: opts.Config,
		onAlert: opts.OnAlert, log: log, owner: owner,
		tickEvery: tickEvery, sustainTicks: sustain,
		settleWithin: within, settleInterval: interval,
		comps:       map[string]Unloadable{},
		leakLatched: map[string]bool{},
	}
}

// Register adds a component for unload-pass consideration. Only names on
// the frozen whitelist are ever unloaded: the pass walks the whitelist and
// looks each name up here, so a registration named "kws" (or anything else
// outside the list) is carried for transparency and structurally
// unreachable (ticket 42 AC#4; D32 16.3.2 revision 3).
func (l *Loop) Register(name string, u Unloadable) {
	if u == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.comps[name] = u
}

// Registered lists the registered component names, sorted.
func (l *Loop) Registered() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, 0, len(l.comps))
	for name := range l.comps {
		out = append(out, name)
	}
	slices.Sort(out)
	return out
}

// SettleFailures returns the current consecutive settle-failure count
// (diagnostics; reset by a passing verification or by raising the fatal
// alert).
func (l *Loop) SettleFailures() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.settleFailures
}

// Start spawns the loop on the D38b resident roster name "watchdog" through
// observe.Registry.Spawn - named, owner-carried, with the recover boundary
// that is the ticket's self-protection. The loop exits when ctx is
// cancelled. This is the only goroutine the package ever starts.
func (l *Loop) Start(ctx context.Context) *observe.Handle {
	root := observe.NewRootFrom(ctx, "watchdog")
	return l.reg.Spawn("watchdog", l.owner, root, l.Run)
}

// Run is the loop body: one tick per TickEvery carrying both duties. A
// panic escaping the body is recovered by the Spawn boundary (the process
// survives; the panic record names this goroutine); a panicking component
// is recovered per call inside the pass.
func (l *Loop) Run(ctx context.Context) {
	t := time.NewTicker(l.tickEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			l.tick(ctx)
		}
	}
}

// tick runs both duties of the one loop: the config mtime+size poll first
// (cheap, and its events should not wait behind enforcement), then the
// resource sample and enforcement.
func (l *Loop) tick(ctx context.Context) {
	l.pollConfig()
	l.enforceResources(ctx)
}

// pollConfig is the watchdog's config duty (SPEC-03 §4.3): hand the tick to
// config.Manager.CheckAndReload - the ticket 05 contract - and log what it
// applied. A failed reload keeps the current config (the manager adopts the
// stat, so there is no per-tick re-erroring to suppress).
func (l *Loop) pollConfig() {
	if l.cfg == nil {
		return
	}
	rep, err := l.cfg.CheckAndReload()
	if err != nil {
		l.log.Warn("watchdog: config poll failed; keeping current config",
			"err", err.Error(), "error_class", string(observe.ClassConfig))
		return
	}
	if rep == nil {
		return
	}
	l.log.Info("watchdog: config poll applied changes",
		"hot", rep.Hot, "reload", rep.Reload, "restart", rep.Restart,
		"locked", len(rep.Locked))
}

// enforceResources samples the tree once and enforces the frozen D32 row of
// the current state: leak indicators first (they alert and never unload),
// then the memory sustain rule. A failed or zero-footprint read is not a
// measurement: it is logged and skipped, never counted as within-cap
// (fail-closed) and never counted toward the sustain window.
func (l *Loop) enforceResources(ctx context.Context) {
	m, err := l.tree.ReadTree()
	if err != nil {
		l.log.Warn("watchdog: tree read failed; tick skipped",
			"err", err.Error(), "error_class", string(observe.ClassResource))
		return
	}
	if m.PrivateWorkingSetBytes <= 0 {
		l.log.Warn("watchdog: tree read returned a zero private working set; tick skipped",
			"error_class", string(observe.ClassResource))
		return
	}
	st := l.stateFn()
	if !st.Valid() {
		l.log.Warn("watchdog: current state is not a D32 sampling state; tick skipped",
			"state", string(st))
		return
	}
	goroutines := l.reg.Count()
	l.log.Debug("watchdog: tick",
		"state", string(st),
		"tree_private_bytes", m.PrivateWorkingSetBytes,
		"mem_cap_bytes", observe.MemCapBytes(st),
		"cpu_limit_percent", observe.CPULimitPercent(st),
		"gdi_objects", m.GDIObjects,
		"user_objects", m.USERObjects,
		"handles", m.Handles,
		"goroutines", goroutines,
		"threads", m.Threads,
		"tree_processes", m.PIDs,
	)

	l.checkLeak(st, m, "handles", m.Handles, observe.HandleLimit())
	l.checkLeak(st, m, "gdi_objects", m.GDIObjects, observe.GDILimit())
	if limit, gated := observe.GoroutineLimit(st); gated {
		l.checkLeak(st, m, "goroutines", goroutines, limit)
	}

	l.sustainMemory(ctx, st, m)
}

// sustainMemory is the D32 sustain rule: consecutive over-cap ticks accrue;
// at SustainTicks one unload pass runs over whitelist-registered components
// and the tree is re-read. Within cap closes the episode ("within cap");
// still over raises the overrun alert, deduped per episode ("or alert").
// In Armed the keep-alive warning with the disable-KWS suggestion goes out
// with the sustained breach itself (D43 row 9), regardless of the pass
// outcome.
func (l *Loop) sustainMemory(ctx context.Context, st observe.SLOState, m observe.TreeMetrics) {
	capBytes := observe.MemCapBytes(st)
	if m.PrivateWorkingSetBytes <= capBytes {
		l.overTicks = 0
		l.overAlerted = false
		l.armedWarned = false
		return
	}
	l.overTicks++
	if l.overTicks < l.sustainTicks {
		return
	}
	before := m.PrivateWorkingSetBytes
	l.log.Info("watchdog: memory over cap sustained; running unload pass",
		"state", string(st),
		"metric", "tree_private_bytes",
		"measured_bytes", before,
		"cap_bytes", capBytes,
		"sustained_ticks", l.overTicks,
	)
	dropped := l.unloadPass(ctx, st)
	after := before
	if m2, err := l.tree.ReadTree(); err == nil && m2.PrivateWorkingSetBytes > 0 {
		after = m2.PrivateWorkingSetBytes
	} else if err != nil {
		// Say the after-number is stale instead of dressing it up as fresh.
		l.log.Warn("watchdog: post-pass tree read failed; after-metrics not refreshed",
			"err", err.Error(), "error_class", string(observe.ClassResource))
	}
	l.log.Info("watchdog: unload pass complete",
		"state", string(st),
		"dropped", dropped,
		"before_bytes", before,
		"after_bytes", after,
		"cap_bytes", capBytes,
	)
	if st == observe.SLOArmed {
		// In Armed the keep-alive warning IS the state's overrun response
		// (D43 row 9): the generic still-over alert below stays off, or the
		// user would hear the same breach twice.
		l.warnArmedKWS(st, before, capBytes)
	}
	if after > capBytes && !l.overAlerted && st != observe.SLOArmed {
		l.overAlerted = true
		l.emit(Alert{
			Kind: AlertOverrun, Class: observe.ClassResource, State: st,
			Metric:   "tree_private_bytes",
			Measured: fmt.Sprintf("%.1fMB", mb(after)),
			Limit:    fmt.Sprintf("<=%.0fMB", mb(capBytes)),
			Detail:   "unload pass ran over the safe whitelist and the tree is still over cap",
			At:       observe.WallTimestampUTC(observe.NowWallUTC()),
		})
	}
	l.overTicks = 0
}

// warnArmedKWS is the D43 row 9 / D32 revision 3 Armed rule: the overrun
// warning goes out with the disable-KWS suggestion, and KWS itself is
// structurally un-unloadable (not on the whitelist). Deduped per over-cap
// episode; the decision stays with the user.
func (l *Loop) warnArmedKWS(st observe.SLOState, measured, capBytes int64) {
	if l.armedWarned {
		return
	}
	l.armedWarned = true
	l.emit(Alert{
		Kind: AlertOverrun, Class: observe.ClassResource, State: st,
		Metric:   "tree_private_bytes",
		Measured: fmt.Sprintf("%.1fMB", mb(measured)),
		Limit:    fmt.Sprintf("<=%.0fMB", mb(capBytes)),
		Detail: "KWS retained (never unloaded in Armed, D32 revision 3); " +
			"suggest disabling KWS instead - the decision stays with the user",
		At: observe.WallTimestampUTC(observe.NowWallUTC()),
	})
}

// unloadPass walks the frozen whitelist in table order and unloads the
// components registered under those names. It never walks the registered
// set: a non-whitelisted registration cannot be reached by construction,
// which is what the AC#4 log audit proves.
func (l *Loop) unloadPass(ctx context.Context, st observe.SLOState) []string {
	l.mu.Lock()
	comps := make(map[string]Unloadable, len(l.comps))
	for k, v := range l.comps {
		comps[k] = v
	}
	l.mu.Unlock()

	var dropped []string
	for _, name := range whitelist {
		u, ok := comps[name]
		if !ok {
			continue
		}
		d, detail, err := l.safeUnload(ctx, name, u)
		switch {
		case err != nil:
			l.log.Error("watchdog: unload action failed",
				"component", name, "state", string(st), "outcome", "error",
				"detail", err.Error(), "error_class", string(observe.ClassInternal))
		case d:
			dropped = append(dropped, name)
			l.log.Info("watchdog: unload action",
				"component", name, "state", string(st), "outcome", "dropped",
				"detail", detail)
		default:
			l.log.Info("watchdog: unload action",
				"component", name, "state", string(st), "outcome", "skipped",
				"detail", detail)
		}
	}
	return dropped
}

// safeUnload recovers a panicking component so one bad Unload cannot kill
// the watchdog loop (ticket 42: panic-isolated). The Spawn boundary remains
// the outer net for everything else.
func (l *Loop) safeUnload(ctx context.Context, name string, u Unloadable) (dropped bool, detail string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("unload of component %q panicked: %v", name, rec)
		}
	}()
	d, det := u.Unload(ctx)
	return d, det, nil
}

// checkLeak is one D42#10 leak indicator: over the frozen cap it raises a
// leak alert with a diagnostic snapshot once per episode (rising edge,
// re-armed when the metric falls back under). It never unloads anything.
func (l *Loop) checkLeak(st observe.SLOState, m observe.TreeMetrics, metric string, measured, limit int) {
	if measured <= limit {
		l.leakLatched[metric] = false
		return
	}
	if l.leakLatched[metric] {
		return
	}
	l.leakLatched[metric] = true
	snap, serr := l.leakSnapshot(st, m, metric, measured, limit)
	if serr != nil {
		// The alert goes out even when the snapshot fails; the failure is
		// recorded instead of swallowing the indicator.
		l.log.Error("watchdog: leak snapshot failed",
			"metric", metric, "err", serr.Error(), "error_class", string(observe.ClassInternal))
	}
	l.log.Info("watchdog: leak indicator over cap; alert + diagnostic snapshot, no auto-unload",
		"state", string(st), "metric", metric,
		"measured", measured, "limit", limit,
		"snapshot_bytes", len(snap), "error_class", string(observe.ClassResource))
	l.emit(Alert{
		Kind: AlertLeak, Class: observe.ClassResource, State: st,
		Metric:   metric,
		Measured: fmt.Sprint(measured),
		Limit:    fmt.Sprintf("<=%d", limit),
		Detail:   "leak indicator crossed its frozen cap; diagnostic snapshot attached, nothing auto-unloaded",
		Snapshot: snap,
		At:       observe.WallTimestampUTC(observe.NowWallUTC()),
	})
}

// leakSnapshot is the D42#10 diagnostic snapshot attached to a leak alert:
// the crossing sample itself (tree metrics), the goroutine roster (the D38b
// registry - a roster of unknowns is the classic leak symptom) and the
// crossed row. JSON, in the shape the diagnostics bundle attaches as its
// slo-snapshot entry (observe.BuildDiagnosticsBundle's SLOSnapshot).
func (l *Loop) leakSnapshot(st observe.SLOState, m observe.TreeMetrics, metric string, measured, limit int) ([]byte, error) {
	snap := struct {
		At         string                  `json:"at"`
		State      observe.SLOState        `json:"state"`
		Metric     string                  `json:"metric"`
		Measured   int                     `json:"measured"`
		Limit      int                     `json:"limit"`
		Goroutines int                     `json:"goroutines"`
		RuntimeGo  int                     `json:"runtime_goroutines"`
		Roster     []observe.GoroutineInfo `json:"roster,omitempty"`
		Tree       observe.TreeMetrics     `json:"tree"`
	}{
		At:         observe.WallTimestampUTC(observe.NowWallUTC()),
		State:      st,
		Metric:     metric,
		Measured:   measured,
		Limit:      limit,
		Goroutines: l.reg.Count(),
		RuntimeGo:  runtime.NumGoroutine(),
		Roster:     l.reg.Snapshot(),
		Tree:       m,
	}
	return json.MarshalIndent(snap, "", "  ")
}

// emit hands an alert to the host sink, in tick order, on the goroutine
// that produced it.
func (l *Loop) emit(a Alert) {
	if l.onAlert != nil {
		l.onAlert(a)
	}
}

func mb(b int64) float64 { return float64(b) / (1 << 20) }
