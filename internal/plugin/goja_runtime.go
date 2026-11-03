package plugin

// Tier-2 goja JS runtime - ticket 51, BLOCK A.
//
// Block A is the runtime envelope, nothing else:
//   - VM lifecycle with per-call isolation (a fresh goja.Runtime per call),
//   - the wall-clock interrupt budget (default 5s) enforced with
//     Runtime.Interrupt() from a named watchdog goroutine,
//   - the result cap (default 256KB) with an explicit truncated flag,
//   - the recover boundary on the calling goroutine that turns a goja panic
//     into a D37 "tool"-class error instead of killing the process.
//
// Block A is NOT (deliberately, so the two legs of ticket 51 do not collide):
//   - the frozen C24 GojaHostAPI surface (wisp.http/fs/clipboard/notify/
//     state/sysinfo entries, capability checks at the bridge, taint flow,
//     audit rows) - that is block B;
//   - the tier2_enabled default-off gate and its L2 re-confirmation - loader
//     side of the ticket;
//   - any memory ceiling. goja has no memory limit and this file does not
//     pretend it has one: the wall-clock budget and the result cap are the
//     ONLY resource ceilings enforced here, and stack depth is left to goja's
//     own limit. A plugin that allocates until the Go heap grows is stopped by
//     the clock, not by a byte ceiling. Claiming a memory bound here would be
//     the D33/F5 fake promise the ticket bans.
//
// Why a new VM per call instead of one cached VM per plugin: the isolation AC
// ("no state crosses calls") cannot be met by a reused VM without wiping the
// global object, and a partial wipe leaves reachable objects in the graph.
// D3's "use-to-wake" lazy VM still holds at the plugin level (no VM exists
// until a call happens); the per-call VM is the cheaper, provable form of it.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"

	"github.com/CarlosShao/wisp/internal/observe"
)

// GojaDefaultWallClockBudget is the default per-call wall-clock interrupt
// budget (ticket 51: 5s). goja checks for interrupts at statement/loop
// boundaries, so the budget is honoured with interrupt-check granularity
// (spike T02: P50 +0.4-0.6ms, max ~12ms), never overshot by more than one
// check point. This is an interrupt of execution, not an SLO measurement.
const GojaDefaultWallClockBudget = 5 * time.Second

// GojaDefaultMaxResultBytes is the default result cap: 256KB. A longer result
// is truncated to this many bytes and the Truncated flag is set; truncation is
// never silent (the call is logged at Warn and the flag is part of the result).
const GojaDefaultMaxResultBytes = 256 * 1024

// gojaWatchdogNamePrefix / Suffix build the D38 roster name of the per-call
// watchdog. PerTaskPrefixes in internal/observe lists "tool-exec-", so the
// generated name classifies as per_task instead of CategoryUnknown (an
// unknown name is a leak symptom by D38b).
const (
	gojaWatchdogNamePrefix = "tool-exec-"
	gojaWatchdogNameSuffix = "-watchdog"
)

// gojaInterruptReason is the value handed to Runtime.Interrupt so the
// classification step can tell a budget hit from a host cancellation without
// parsing message text.
type gojaInterruptReason string

const (
	gojaInterruptWallClock  gojaInterruptReason = "wall_clock_budget_exceeded"
	gojaInterruptHostCancel gojaInterruptReason = "host_cancelled"
)

// gojaRendererSource is the JS used to turn a call result into text. Strings
// pass through raw (a plugin returning a string means that string); anything
// else goes through JSON.stringify and falls back to String(v), so a cyclic or
// exotic value still yields something the model can read instead of failing
// the whole call. It is compiled into a throw-away URL and run in the call's
// own VM, so it defines no globals and cannot leak between calls.
const gojaRendererSource = `(function(v){` +
	`if(typeof v==="string"){return v}` +
	`try{var s=JSON.stringify(v);if(typeof s==="string"){return s}}catch(e){}` +
	`try{return String(v)}catch(e){return ""}}` +
	`)`

const gojaRendererURL = "plugin://wisp/tier2-render.js"

// GojaRuntime is the Block-A runtime envelope. It is stateless apart from its
// configuration, so one instance may serve many plugins and many goroutines.
// Zero value is not usable; use NewGojaRuntime.
type GojaRuntime struct {
	timeout        time.Duration
	maxResultBytes int
	reg            *observe.Registry
}

// GojaOption configures a GojaRuntime.
type GojaOption func(*GojaRuntime)

// WithGojaWallClockBudget overrides the per-call interrupt budget. Use 0 to
// keep GojaDefaultWallClockBudget; a negative value is rejected at Execute
// time (fail closed, never silently "unlimited").
func WithGojaWallClockBudget(d time.Duration) GojaOption {
	return func(r *GojaRuntime) { r.timeout = d }
}

// WithGojaMaxResultBytes overrides the result cap, in bytes.
func WithGojaMaxResultBytes(n int) GojaOption {
	return func(r *GojaRuntime) { r.maxResultBytes = n }
}

// WithGojaRegistry spawns the watchdog through a specific registry instead of
// observe.Default (tests, isolated subsystems).
func WithGojaRegistry(reg *observe.Registry) GojaOption {
	return func(r *GojaRuntime) {
		if reg != nil {
			r.reg = reg
		}
	}
}

// NewGojaRuntime builds the runtime with the ticket's defaults (5s wall clock,
// 256KB result cap, default goroutine registry).
func NewGojaRuntime(opts ...GojaOption) *GojaRuntime {
	r := &GojaRuntime{
		timeout:        GojaDefaultWallClockBudget,
		maxResultBytes: GojaDefaultMaxResultBytes,
		reg:            observe.Default,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

// WallClockBudget returns the configured per-call interrupt budget.
func (r *GojaRuntime) WallClockBudget() time.Duration { return r.timeout }

// MaxResultBytes returns the configured result cap.
func (r *GojaRuntime) MaxResultBytes() int { return r.maxResultBytes }

// GojaCall is one Tier-2 tool invocation. Source is the plugin bundle text;
// Function names the entry the bundle exposes (empty = take the value of the
// last statement, which keeps plain-script plugins working). Args are Go
// values converted by the VM into JS values. Scope is optional: pass the
// caller's tool/task DisposalScope to have the watchdog joined by that scope,
// otherwise Execute owns a per-call scope and disposes it.
type GojaCall struct {
	PluginID string
	Source   string
	Function string
	Args     []any
	Scope    *DisposalScope
}

// GojaResult is the capped, rendered outcome of one call. FullBytes is the
// size before the cap, ReturnedBytes after it, so a truncated result is
// measurable rather than silently shorter.
type GojaResult struct {
	Value         string `json:"value"`
	Truncated     bool   `json:"truncated"`
	FullBytes     int    `json:"full_bytes"`
	ReturnedBytes int    `json:"returned_bytes"`
	Interrupted   bool   `json:"interrupted"`
}

// GojaWatchdogName is the D38 roster name the per-call watchdog is spawned
// under, derived from the plugin id (non-name characters replaced by '-').
func GojaWatchdogName(pluginID string) string {
	return gojaWatchdogNamePrefix + sanitizeGoroutineName(pluginID) + gojaWatchdogNameSuffix
}

func sanitizeGoroutineName(s string) string {
	if s == "" {
		return "unnamed"
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9',
			r == '-', r == '_', r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	return b.String()
}

// Execute runs one isolated call. It never returns a raw goja error or a
// panic: every failure is a classified D37 error (tool-class for a JS throw,
// a syntax error, an interrupt or an internal panic; cancelled for a host
// cancellation; config for a misconfigured envelope).
func (r *GojaRuntime) Execute(ctx context.Context, call GojaCall) (out GojaResult, err error) {
	if r == nil {
		return out, observe.New(observe.ClassInternal, "plugin: Execute called on a nil *GojaRuntime")
	}
	if r.timeout <= 0 || r.maxResultBytes <= 0 {
		return out, observe.New(observe.ClassConfig, fmt.Sprintf(
			"plugin %q: tier2 runtime is misconfigured (wall-clock budget %v, result cap %d bytes); both must be positive, the runtime refuses to run unlimited",
			call.PluginID, r.timeout, r.maxResultBytes))
	}
	if strings.TrimSpace(call.PluginID) == "" {
		return out, observe.New(observe.ClassConfig,
			"tier2 goja call carries no plugin id: the watchdog name, the logs and the audit row are keyed on it")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return out, observe.Wrap(observe.ClassCancelled, ctx.Err(), fmt.Sprintf(
			"plugin %q: tier2 call was cancelled by the host before it started", call.PluginID))
	}

	// Last-resort recover boundary, ON THE CALLING GOROUTINE: a panic that
	// escapes the guarded block (including one from the disposal tail) becomes
	// a tool-class error instead of taking the process down. Registered first
	// so it runs last, after the teardown defers below.
	defer func() {
		if p := recover(); p != nil {
			out = GojaResult{}
			err = r.panicError(call.PluginID, p)
		}
	}()

	scope := call.Scope
	owned := scope == nil
	if owned {
		scope = NewDisposalScope("tool:goja:"+call.PluginID, ctx, WithRegistry(r.reg))
		defer func() {
			if res := scope.Dispose(); res.Incomplete {
				slog.Warn("tier2 call scope did not dispose cleanly",
					"plugin", call.PluginID, "timed_out", res.TimedOut, "failed_steps", len(res.Failed))
			}
		}()
	}

	// Per-call VM: the isolation mechanism. Nothing is cached, nothing is
	// reused, so no global, closure or object graph survives into the next
	// call of the same plugin.
	vm := goja.New()
	stop := r.armWatchdog(scope, vm, call.PluginID)
	defer stop()

	rendered, gerr := r.guardedCall(call.PluginID, func() (string, error) {
		v, e := r.runProgram(vm, call)
		if e != nil {
			return "", e
		}
		return renderGojaValue(vm, v)
	})
	if gerr != nil {
		return r.classify(call.PluginID, gerr)
	}

	out.FullBytes = len(rendered)
	if len(rendered) > r.maxResultBytes {
		// Truncate on a rune boundary so the model never receives a broken
		// UTF-8 sequence, and say so - in the result and in the log.
		out.Value = truncateAtRuneBoundary(rendered, r.maxResultBytes)
		out.Truncated = true
		out.ReturnedBytes = len(out.Value)
		slog.Warn("tier2 result truncated",
			"plugin", call.PluginID,
			"full_bytes", out.FullBytes,
			"returned_bytes", out.ReturnedBytes,
			"cap_bytes", r.maxResultBytes,
		)
		return out, nil
	}
	out.Value = rendered
	out.ReturnedBytes = len(rendered)
	return out, nil
}

// armWatchdog starts the named goroutine that enforces the wall-clock budget.
// Only Runtime.Interrupt() is called from it - the VM itself is touched by the
// calling goroutine alone, which is goja's documented concurrency rule. The
// returned stop function disarms the timer and joins the goroutine, so a call
// never leaves a watchdog behind.
func (r *GojaRuntime) armWatchdog(scope *DisposalScope, vm *goja.Runtime, pluginID string) func() {
	done := make(chan struct{})
	var once sync.Once
	h := scope.Go(GojaWatchdogName(pluginID), func(ctx context.Context) {
		timer := time.NewTimer(r.timeout)
		defer timer.Stop()
		select {
		case <-done:
			return
		case <-ctx.Done():
			// The host cancelled. Prefer a clean stop if the call already
			// finished on this same wake-up.
			select {
			case <-done:
				return
			default:
			}
			vm.Interrupt(gojaInterruptHostCancel)
		case <-timer.C:
			vm.Interrupt(gojaInterruptWallClock)
		}
	})
	return func() {
		once.Do(func() { close(done) })
		<-h.Done()
	}
}

// guardedCall is the recover boundary for the VM work. A panic raised inside
// goja (its own invariant break, or one from a Go value it walks) is converted
// to an internal marker that classify turns into a tool-class error, with the
// stack logged so the failure is diagnosable and never silent. Running the
// program here, on the caller's goroutine, is what keeps the boundary
// meaningful: a panic that crossed a spawned goroutine would land on the
// registry boundary instead (class internal), not on this tool result.
func (r *GojaRuntime) guardedCall(pluginID string, fn func() (string, error)) (out string, err error) {
	defer func() {
		if p := recover(); p != nil {
			out = ""
			err = &gojaPanic{pluginID: pluginID, recovered: p, stack: string(debug.Stack())}
		}
	}()
	return fn()
}

// gojaPanic marks a recovered goja panic between the guarded call and the
// classifier.
type gojaPanic struct {
	pluginID  string
	recovered any
	stack     string
}

func (g *gojaPanic) Error() string {
	return fmt.Sprintf("goja panic in plugin %q: %v", g.pluginID, g.recovered)
}

func (r *GojaRuntime) panicError(pluginID string, p any) error {
	slog.Error("tier2 goja call panicked at the outer boundary",
		"plugin", pluginID, "recovered", fmt.Sprint(p), "stack", string(debug.Stack()))
	return observe.New(observe.ClassTool, fmt.Sprintf(
		"plugin %q: tier2 runtime panicked (%v); recovered at the calling-goroutine boundary, host and other plugins unaffected",
		pluginID, p))
}

// runProgram compiles and runs the bundle, then calls the declared entry
// function when one is named.
func (r *GojaRuntime) runProgram(vm *goja.Runtime, call GojaCall) (goja.Value, error) {
	if call.Source == "" {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"plugin %q: tier2 call has an empty source", call.PluginID))
	}
	prg, err := goja.Compile("plugin://"+call.PluginID+"/tier2.js", call.Source, false)
	if err != nil {
		return nil, err
	}
	v, err := vm.RunProgram(prg)
	if err != nil || call.Function == "" {
		return v, err
	}
	fnVal := vm.Get(call.Function)
	if fnVal == nil || goja.IsUndefined(fnVal) || goja.IsNull(fnVal) {
		return nil, observe.New(observe.ClassTool, fmt.Sprintf(
			"plugin %q: entry function %q is not defined by the bundle", call.PluginID, call.Function))
	}
	fn, ok := goja.AssertFunction(fnVal)
	if !ok {
		return nil, observe.New(observe.ClassTool, fmt.Sprintf(
			"plugin %q: entry %q exists but is not callable (got %s)", call.PluginID, call.Function, fnVal.String()))
	}
	argv := make([]goja.Value, 0, len(call.Args))
	for _, a := range call.Args {
		argv = append(argv, vm.ToValue(a))
	}
	return fn(goja.Undefined(), argv...)
}

// renderGojaValue turns a JS result into text inside the call's own VM.
func renderGojaValue(vm *goja.Runtime, v goja.Value) (string, error) {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return "", nil
	}
	prg, err := goja.Compile(gojaRendererURL, gojaRendererSource, false)
	if err != nil {
		return "", err
	}
	rv, err := vm.RunProgram(prg)
	if err != nil {
		return "", err
	}
	fn, ok := goja.AssertFunction(rv)
	if !ok {
		return "", fmt.Errorf("tier2 renderer did not evaluate to a function")
	}
	out, err := fn(goja.Undefined(), v)
	if err != nil {
		return "", err
	}
	if out == nil || goja.IsUndefined(out) || goja.IsNull(out) {
		return "", nil
	}
	return out.String(), nil
}

// classify maps whatever the guarded call produced onto the D37 taxonomy. A JS
// throw, a syntax error and a goja panic are all tool-class: the outcome goes
// back to the model, which decides whether to try again (D37 tool row,
// RetryAgentDecided). A wall-clock hit is tool-class too, and sets the
// Interrupted flag so the caller can tell "the plugin is wrong" from "the
// plugin ran too long".
func (r *GojaRuntime) classify(pluginID string, cause error) (GojaResult, error) {
	var ie *goja.InterruptedError
	if errors.As(cause, &ie) {
		if reason, ok := ie.Value().(gojaInterruptReason); ok && reason == gojaInterruptHostCancel {
			return GojaResult{}, observe.Wrap(observe.ClassCancelled, cause, fmt.Sprintf(
				"plugin %q: tier2 call was cancelled by the host while the VM was running", pluginID))
		}
		return GojaResult{Interrupted: true}, observe.New(observe.ClassTool, fmt.Sprintf(
			"plugin %q: tier2 call exceeded its %v wall-clock budget and was interrupted at the next interrupt check point "+
				"(goja offers no memory ceiling; the clock and the %d-byte result cap are the only limits enforced)",
			pluginID, r.timeout, r.maxResultBytes))
	}
	var gp *gojaPanic
	if errors.As(cause, &gp) {
		slog.Error("tier2 goja panic recovered at the calling-goroutine boundary",
			"plugin", gp.pluginID, "recovered", fmt.Sprint(gp.recovered), "stack", gp.stack)
		return GojaResult{}, observe.New(observe.ClassTool, fmt.Sprintf(
			"plugin %q: tier2 runtime panicked (%v); recovered, host and other plugins unaffected",
			pluginID, gp.recovered))
	}
	var oe *observe.Error
	if errors.As(cause, &oe) {
		// Already classified by the envelope itself (config, tool: missing
		// entry function, ...). Pass it through unchanged.
		return GojaResult{}, cause
	}
	return GojaResult{}, observe.New(observe.ClassTool, fmt.Sprintf(
		"plugin %q: tier2 call failed: %s", pluginID, jsExceptionText(cause)))
}

// jsExceptionText renders a JS-side failure the way a plugin author reads it:
// the exception value's own text ("TypeError: x is not a function"), or the
// compiler's file:line message for a syntax error.
func jsExceptionText(cause error) string {
	// A JS throw surfaces as *goja.Exception; Value() is the thrown value
	// itself (an Error object renders as "TypeError: ..."), which is the text
	// the plugin author recognises. Block B's audit rows can pull Value() for
	// the structured form.
	var ex *goja.Exception
	if errors.As(cause, &ex) {
		if v := ex.Value(); v != nil {
			if s := v.String(); s != "" {
				return s
			}
		}
		return ex.Error()
	}
	return cause.Error()
}

// truncateAtRuneBoundary cuts s to at most max bytes without splitting a UTF-8
// rune; the result is always valid UTF-8 when the input is.
func truncateAtRuneBoundary(s string, max int) string {
	if len(s) <= max {
		return s
	}
	cut := max
	for cut > 0 && !isRuneStart(s[cut]) {
		cut--
	}
	return s[:cut]
}

func isRuneStart(b byte) bool { return b&0xC0 != 0x80 }
