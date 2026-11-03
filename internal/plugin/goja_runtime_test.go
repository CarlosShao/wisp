package plugin

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// newTestRuntime returns a runtime on its own registry so the tests can assert
// "no goroutine left behind" without fighting other package tests for the
// default registry.
func newTestRuntime(opts ...GojaOption) (*GojaRuntime, *observe.Registry) {
	reg := observe.NewRegistry()
	r := NewGojaRuntime(append([]GojaOption{WithGojaRegistry(reg)}, opts...)...)
	return r, reg
}

// TestGojaBlockADefaults pins the two frozen Block-A numbers: 5s wall clock and
// 256KB result cap.
func TestGojaBlockADefaults(t *testing.T) {
	t.Parallel()
	r, _ := newTestRuntime()
	if GojaDefaultWallClockBudget != 5*time.Second {
		t.Errorf("GojaDefaultWallClockBudget = %v, want 5s", GojaDefaultWallClockBudget)
	}
	if GojaDefaultMaxResultBytes != 256*1024 {
		t.Errorf("GojaDefaultMaxResultBytes = %d, want %d", GojaDefaultMaxResultBytes, 256*1024)
	}
	if r.WallClockBudget() != GojaDefaultWallClockBudget {
		t.Errorf("budget = %v, want %v", r.WallClockBudget(), GojaDefaultWallClockBudget)
	}
	if r.MaxResultBytes() != GojaDefaultMaxResultBytes {
		t.Errorf("cap = %d, want %d", r.MaxResultBytes(), GojaDefaultMaxResultBytes)
	}
}

// TestGojaBlockAMisconfigRejected: a non-positive budget or cap refuses the
// call with a config error instead of running unbounded.
func TestGojaBlockAMisconfigRejected(t *testing.T) {
	t.Parallel()
	for _, opts := range [][]GojaOption{
		{WithGojaWallClockBudget(0)},
		{WithGojaWallClockBudget(-time.Second)},
		{WithGojaMaxResultBytes(0)},
		{WithGojaMaxResultBytes(-1)},
	} {
		r, _ := newTestRuntime(opts...)
		out, err := r.Execute(context.Background(), GojaCall{PluginID: "cfg", Source: "1"})
		if err == nil || observe.ClassOfOrInternal(err) != observe.ClassConfig {
			t.Fatalf("misconfigured envelope must fail closed with config class, got %v / %v", err, out)
		}
	}
	// No plugin id: the watchdog name and the audit row are keyed on it.
	r, _ := newTestRuntime()
	if _, err := r.Execute(context.Background(), GojaCall{Source: "1"}); err == nil ||
		observe.ClassOfOrInternal(err) != observe.ClassConfig {
		t.Fatalf("call without plugin id must be rejected, got %v", err)
	}
}

// TestGojaPerCallIsolation is AC-1: the SAME runtime, two consecutive calls,
// and a global written by the first must not be visible to the second.
func TestGojaPerCallIsolation(t *testing.T) {
	t.Parallel()
	r, reg := newTestRuntime()
	ctx := context.Background()

	first, err := r.Execute(ctx, GojaCall{
		PluginID: "iso",
		Source:   `(function(){globalThis.__wispProbe = 42; return "written"})()`,
	})
	if err != nil {
		t.Fatalf("first call: %v", err)
	}
	if strings.Trim(first.Value, `"`) != "written" {
		t.Fatalf("first value = %q", first.Value)
	}

	second, err := r.Execute(ctx, GojaCall{
		PluginID: "iso",
		Source:   `(function(){ return typeof globalThis.__wispProbe })()`,
	})
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if got := strings.Trim(second.Value, `"`); got != "undefined" {
		t.Fatalf("state leaked across calls: typeof globalThis.__wispProbe = %q, want \"undefined\"", got)
	}
	if n := reg.Count(); n != 0 {
		t.Fatalf("watchdog goroutine leaked: registry count = %d", n)
	}
}

// TestGojaWallClockInterrupt is AC-2: a spinning loop is stopped by the
// watchdog's vm.Interrupt at the budget, becomes a tool-class error, sets
// Interrupted, leaves the host alive and joins its goroutine.
func TestGojaWallClockInterrupt(t *testing.T) {
	t.Parallel()
	budget := 5 * time.Second // the ticket's default, asserted at full size
	r, reg := newTestRuntime()
	if r.WallClockBudget() != budget {
		t.Fatalf("default budget drifted: %v", r.WallClockBudget())
	}

	start := observe.NowWallUTC()
	out, err := r.Execute(context.Background(), GojaCall{
		PluginID: "spin",
		Source:   `function spin(){ var i = 0; while (true) { i++; } }`,
		Function: "spin",
	})
	elapsed := observe.NowWallUTC().Sub(start)

	if err == nil {
		t.Fatalf("infinite loop must be interrupted, got result %q", out.Value)
	}
	if class := observe.ClassOfOrInternal(err); class != observe.ClassTool {
		t.Errorf("interrupt error class = %q, want tool (D37: the model decides whether to retry)", class)
	}
	if !out.Interrupted {
		t.Errorf("Interrupted flag not set")
	}
	// The budget is honoured with interrupt-check granularity (spike T02: max
	// ~12ms), never early; the ceiling only bounds a runaway implementation.
	if elapsed < budget {
		t.Errorf("interrupted after %v, before the %v budget", elapsed, budget)
	}
	if elapsed > budget+3*time.Second {
		t.Errorf("interrupt took %v, far past the %v budget", elapsed, budget)
	}
	if n := reg.Count(); n != 0 {
		t.Errorf("watchdog goroutine still live after Execute returned: %d", n)
	}
}

// TestGojaInterruptIsPerCallBounded checks the watchdog does not outlive the
// budget it armed: two fast calls must not be interrupted by the first call's
// timer.
func TestGojaInterruptIsPerCallBounded(t *testing.T) {
	t.Parallel()
	r, _ := newTestRuntime(WithGojaWallClockBudget(200 * time.Millisecond))
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		out, err := r.Execute(ctx, GojaCall{PluginID: "fast", Source: `(function(){ return 6*7 })()`})
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		if out.Interrupted {
			t.Fatalf("call %d wrongly reported interrupted", i)
		}
		if got := strings.Trim(out.Value, `"`); got != "42" {
			t.Fatalf("call %d value = %q", i, out.Value)
		}
		time.Sleep(150 * time.Millisecond) // past the armed budget of the previous call
	}
}

// TestGojaResultCapTruncation is AC-3: a result bigger than the cap comes back
// truncated, flagged, with sizes reported, never silently shorter.
func TestGojaResultCapTruncation(t *testing.T) {
	t.Parallel()
	capBytes := 256 * 1024
	r, _ := newTestRuntime()
	if r.MaxResultBytes() != capBytes {
		t.Fatalf("default cap drifted: %d", r.MaxResultBytes())
	}

	// ~1MB of ASCII built inside the VM.
	out, err := r.Execute(context.Background(), GojaCall{
		PluginID: "big",
		Source:   `(function(){ var chunk = ""; for (var i = 0; i < 1024; i++) { chunk += "x" }; var s = ""; for (var j = 0; j < 512; j++) { s += chunk }; return s })()`,
	})
	if err != nil {
		t.Fatalf("big result call: %v", err)
	}
	if !out.Truncated {
		t.Fatalf("Truncated flag must be set for a %d-byte result over a %d-byte cap", out.FullBytes, capBytes)
	}
	if out.FullBytes != 512*1024 {
		t.Fatalf("FullBytes = %d, expected the pre-cap size", out.FullBytes)
	}
	if out.ReturnedBytes != len(out.Value) {
		t.Fatalf("ReturnedBytes = %d but value is %d bytes", out.ReturnedBytes, len(out.Value))
	}
	if out.ReturnedBytes != capBytes {
		t.Fatalf("returned %d bytes, want exactly the cap %d", out.ReturnedBytes, capBytes)
	}

	// Multi-byte safety: a cut that would land inside a rune backs off to the last
	// rune boundary instead.
	rSmall, _ := newTestRuntime(WithGojaMaxResultBytes(5))
	outSmall, err := rSmall.Execute(context.Background(), GojaCall{
		PluginID: "utf8",
		// "aaa" then three 3-byte runes: cut at 5 would split a rune.
		Source: `(function(){ return "aaa\u00e9\u00e9\u00e9" })()`,
	})
	if err != nil {
		t.Fatalf("utf8 call: %v", err)
	}
	if !outSmall.Truncated {
		t.Fatalf("expected truncation")
	}
	// The renderer passes a JS string through raw, so the Go value is the unquoted text.
	want := "aaa\u00e9\u00e9\u00e9" // "aaa" + three 2-byte runes = 9 bytes
	// Cap 5 lands on s[5], which is a rune-start (the second U+00E9), so the
	// cut keeps "aaa" + one complete rune.
	if got := outSmall.Value; got != want[:5] || len(got) != 5 {
		t.Fatalf("value = %q (%d bytes), want %q (%d bytes)", got, len(got), want[:5], 5)
	}
}

// TestGojaPanicBoundary is AC-4: a panic raised inside the VM work is recovered
// on the calling goroutine and converted to a tool-class error, with the stack
// kept for diagnosis. JS throws are tool-class too.
func TestGojaPanicBoundary(t *testing.T) {
	t.Parallel()
	id := "panic"
	r, reg := newTestRuntime()

	// The guarded block that Execute routes all VM work through.
	rendered, err := r.guardedCall(id, func() (string, error) {
		panic("simulated goja internal panic")
	})
	if rendered != "" {
		t.Fatalf("panicking call must produce no value, got %q", rendered)
	}
	var gp *gojaPanic
	if !errors.As(err, &gp) {
		t.Fatalf("guardedCall must convert a panic into *gojaPanic, got %T: %v", err, err)
	}
	if gp.pluginID != id {
		t.Errorf("panic marker lost the plugin id: %q", gp.pluginID)
	}
	if !strings.Contains(gp.stack, "simulated goja internal panic") && !strings.Contains(gp.stack, "goja_runtime_test.go") {
		t.Errorf("panic marker has no usable stack: %q", gp.stack)
	}
	out, classified := r.classify(id, err)
	if observe.ClassOfOrInternal(classified) != observe.ClassTool {
		t.Errorf("panic must classify as tool, got %q", observe.ClassOfOrInternal(classified))
	}
	if out.Value != "" {
		t.Errorf("a panicked call must not return a value: %q", out.Value)
	}

	// End to end through Execute: a JS throw is a tool error, the host lives.
	_, thrown := r.Execute(context.Background(), GojaCall{
		PluginID: id,
		Source:   `(function(){ throw new Error("plugin said no") })()`,
	})
	if thrown == nil {
		t.Fatalf("JS throw must produce an error")
	}
	if observe.ClassOfOrInternal(thrown) != observe.ClassTool {
		t.Errorf("JS throw class = %q, want tool", observe.ClassOfOrInternal(thrown))
	}
	if !strings.Contains(thrown.Error(), "plugin said no") {
		t.Errorf("JS throw message must reach the caller: %q", thrown.Error())
	}

	// A syntax error is likewise a tool error (the model gets to see it).
	if _, syn := r.Execute(context.Background(), GojaCall{PluginID: id, Source: "var = =;"}); syn == nil ||
		observe.ClassOfOrInternal(syn) != observe.ClassTool {
		t.Errorf("syntax error must classify as tool, got %v", syn)
	}
	if n := reg.Count(); n != 0 {
		t.Errorf("no goroutine may survive a panic boundary: %d", n)
	}
}

// TestGojaHostCancelInterrupts verifies a host cancellation lands on the
// cancelled class rather than the tool class.
func TestGojaHostCancelInterrupts(t *testing.T) {
	t.Parallel()
	r, reg := newTestRuntime(WithGojaWallClockBudget(30 * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	out, err := r.Execute(ctx, GojaCall{
		PluginID: "cancel",
		Source:   `function spin(){ while (true) {} }`,
		Function: "spin",
	})
	if err == nil {
		t.Fatalf("cancelled host ctx must stop the call, got %q", out.Value)
	}
	if class := observe.ClassOfOrInternal(err); class != observe.ClassCancelled {
		t.Errorf("class = %q, want cancelled (%v)", class, err)
	}
	if out.Interrupted {
		t.Errorf("host cancellation is not a wall-clock overrun; Interrupted must stay false")
	}
	if n := reg.Count(); n != 0 {
		t.Errorf("watchdog leaked: %d", n)
	}
}

// TestGojaWatchdogRosterName pins that the watchdog name classifies inside the
// D38 roster (an unknown name is a leak symptom by D38b).
func TestGojaWatchdogRosterName(t *testing.T) {
	t.Parallel()
	name := GojaWatchdogName("com.example/My Plugin")
	if got := observe.ClassifyGoroutine(name); got != observe.CategoryPerTask {
		t.Errorf("watchdog name %q classifies as %q, want per_task", name, got)
	}
	if strings.ContainsAny(name, " /") {
		t.Errorf("watchdog name kept illegal characters: %q", name)
	}
}
