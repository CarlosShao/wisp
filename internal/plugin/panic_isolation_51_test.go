package plugin

// Ticket 51 AC#4 - "Panic isolation: JS triggering goja panic -> recovered,
// other plugin unaffected, host alive".
//
// Block A (8a2a40b0f) built the recover boundary and proved it by throwing a
// panic into guardedCall DIRECTLY. That proves the recover function recovers; it
// does not prove the sentence this AC actually writes, which starts from JS: a
// plugin calls a host entry and the Go implementation behind that entry crashes
// (a bad assumption, a nil map write). goja does not turn such a panic into a JS
// exception - exceptionFromValue() answers nil for a non-Exception value and
// handleThrow() re-panics it - so the value leaves RunProgram as a raw Go panic
// and must be caught on the calling goroutine, exactly once, without taking the
// host or a neighbour plugin with it.
//
// What this file pins, all through a real VM and real JS text:
//
//  1. the crash is recovered (the process lives) and the call fails as D37
//     `tool`, which is what the model gets to see;
//  2. the failure is ATTRIBUTED: the error the model reads and the audit line an
//     operator greps both name the entry, its capability and the tool, and the
//     attempt is on surface.Calls() - "something panicked" is not actionable,
//     "wisp.fs.read of tool js.read panicked in its host implementation" is;
//  3. the plugin's own try/catch cannot swallow a host crash (a host bug is not
//     a plugin-controllable exception the plugin can retry against in a loop);
//  4. other plugins are unaffected, sequentially and concurrently, on the same
//     runtime and behind the same enable gate;
//  5. nothing is left running (watchdog joined, registry empty) and the runtime
//     keeps serving - "host alive" is a reading, not an adjective.

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// crashHost is a C24Host whose Go implementation crashes, which is the shape of
// a host bug rather than of a refusal (a refusal is C24Reply.Refusal and is
// covered by c24_test.go). panicVal == nil selects the nil-map-write form: the
// classic nil-dereference host bug, which arrives as a runtime.Error.
type crashHost struct {
	mu       sync.Mutex
	calls    []C24Request
	crashing bool
	panicVal any
	reply    C24Reply
}

// CallC24 serves the request and, when the fixture is set to crash, raises the
// host-bug panic shape selected by panicVal: an explicit panic when panicVal is
// set, otherwise the runtime.Error of a write to a nil map, which staticcheck
// flags as SA5000 (suppressed where it fires; the write is the subject, see
// there).
func (h *crashHost) CallC24(_ context.Context, req C24Request) (C24Reply, error) {
	h.mu.Lock()
	h.calls = append(h.calls, req)
	crashing, val, reply := h.crashing, h.panicVal, h.reply
	h.mu.Unlock()
	if crashing {
		if val != nil {
			panic(val)
		}
		var m map[string]string
		//lint:ignore SA5000 the nil-map write on this line IS the subject: ticket 51 AC#4 requires the host to survive the runtime.Error shape of a Go host bug, not only the explicit panic shape crashHost raises when panicVal is set, and a write to a nil map is the cheapest deterministic way to raise one - removing it removes that half of the assertion. Intentional, registered in pending-and-issues A1773, re-verified independently by wave 2 (A1833).
		m["boom"] = "nil map write" // runtime.Error, no explicit panic
	}
	return reply, nil
}

func (h *crashHost) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.calls)
}

// TestC24HostCrashIsRecoveredAndAttributed is AC#4's first half: JS reaches the
// host, the host crashes, the call is recovered and the failure says which entry
// and which tool crashed. The plugin wraps the call in try/catch on purpose - a
// host crash must not become a value the plugin can swallow.
func TestC24HostCrashIsRecoveredAndAttributed(t *testing.T) {
	t.Parallel()
	const (
		plugin = "boom"
		tool   = "js.read"
	)
	log := &c24Log{}
	host := &crashHost{crashing: true, panicVal: "host exploded while reading"}
	surface := NewC24Surface(plugin, tool, []string{"fs.read"}, host, WithC24Logf(log.Logf))
	reg := observe.NewRegistry()
	r := NewGojaRuntime(WithGojaRegistry(reg))

	out, err := r.Execute(context.Background(), GojaCall{
		PluginID: plugin,
		Source: `(function(){` +
			` try { return "swallowed:" + wisp.fs.read({path:"C:/a"}) }` +
			` catch (e) { return "swallowed:" + e }` +
			` })()`,
		Surface: surface,
	})

	if err == nil {
		t.Fatalf("a crashing host must fail the call, got value %q", out.Value)
	}
	if out.Value != "" {
		t.Fatalf("the plugin's own try/catch swallowed a host crash: %q", out.Value)
	}
	if class := observe.ClassOfOrInternal(err); class != observe.ClassTool {
		t.Errorf("host crash class = %q, want tool (D37: the model decides whether to retry)", class)
	}

	// Attribution: the text the model reads has to name the entry and the tool,
	// otherwise the only recourse is "tier 2 is broken somewhere".
	msg := err.Error()
	for _, want := range []string{C24FSRead, "fs.read", plugin, tool} {
		if !strings.Contains(msg, want) {
			t.Errorf("the failure does not name %q, so nobody can tell which host entry crashed: %q", want, msg)
		}
	}

	// Attribution in the audit trail: one greppable line, entry and capability
	// on it, the same shape as C24-CALL / C24-REFUSE / C24-HOSTERR.
	if n := log.CountContaining("C24-HOSTPANIC"); n != 1 {
		t.Errorf("C24-HOSTPANIC audit lines = %d, want 1; log:\n%s", n, log.Joined())
	}
	for _, want := range []string{"entry=" + C24FSRead, "capability=fs.read", "plugin=" + plugin, "tool=" + tool} {
		if !strings.Contains(log.Joined(), want) {
			t.Errorf("the audit line does not carry %q; log:\n%s", want, log.Joined())
		}
	}

	// The attempt is a call the VM really made: it belongs on Calls(), otherwise
	// the one call that took the host down is the one call missing from the log.
	if got := surface.Calls(); len(got) != 1 || got[0] != C24FSRead {
		t.Errorf("surface.Calls() = %v, want exactly [%s]", got, C24FSRead)
	}

	// Host alive: the per-call watchdog is joined and no goroutine survives.
	if n := reg.Count(); n != 0 {
		t.Errorf("goroutines left behind after a recovered host crash: %d", n)
	}
	if out2, err2 := r.Execute(context.Background(), GojaCall{
		PluginID: "after", Source: `(function(){ return "alive" })()`,
	}); err2 != nil || out2.Value != "alive" {
		t.Errorf("the runtime stopped serving after a recovered crash: value=%q err=%v", out2.Value, err2)
	}
}

// TestC24HostRuntimeErrorIsRecoveredAndAttributed is the same claim for the other
// panic shape a host bug produces: a runtime.Error (nil map write / nil
// dereference) rather than an explicit panic("..."). No try/catch on the JS side
// this time, so this also pins the plain "the call just fails" path.
func TestC24HostRuntimeErrorIsRecoveredAndAttributed(t *testing.T) {
	t.Parallel()
	log := &c24Log{}
	host := &crashHost{crashing: true} // nil panicVal => nil map write
	surface := NewC24Surface("crash", "js.write", []string{"fs.write"}, host, WithC24Logf(log.Logf))
	reg := observe.NewRegistry()
	r := NewGojaRuntime(WithGojaRegistry(reg))

	out, err := r.Execute(context.Background(), GojaCall{
		PluginID: "crash",
		Source:   `(function(){ return wisp.fs.write({path:"C:/a", data:"x"}) })()`,
		Surface:  surface,
	})
	if err == nil {
		t.Fatalf("a runtime.Error from the host must fail the call, got %q", out.Value)
	}
	if class := observe.ClassOfOrInternal(err); class != observe.ClassTool {
		t.Errorf("class = %q, want tool", class)
	}
	msg := err.Error()
	if !strings.Contains(msg, C24FSWrite) || !strings.Contains(msg, "js.write") {
		t.Errorf("the failure does not attribute the crash to %s of tool js.write: %q", C24FSWrite, msg)
	}
	if !strings.Contains(log.Joined(), "entry="+C24FSWrite) {
		t.Errorf("no audit line names %s; log:\n%s", C24FSWrite, log.Joined())
	}
	if got := host.Count(); got != 1 {
		t.Errorf("host calls = %d, want exactly 1 (one crashing call, no retry by the runtime)", got)
	}
	if n := reg.Count(); n != 0 {
		t.Errorf("goroutines left behind: %d", n)
	}
}

// TestC24HostCrashLeavesOtherPluginsRunning is AC#4's second half: a crashing
// plugin must not disturb a neighbour, whether the neighbour runs before, after,
// or at the same moment on the same runtime envelope.
func TestC24HostCrashLeavesOtherPluginsRunning(t *testing.T) {
	t.Parallel()
	reg := observe.NewRegistry()
	r := NewGojaRuntime(WithGojaRegistry(reg))
	ctx := context.Background()

	crashing := &crashHost{crashing: true, panicVal: "boom"}
	healthy := &crashHost{reply: C24Reply{Value: map[string]any{"ok": true}}}
	quiet := &c24Log{}

	callBoom := func() (GojaResult, error) {
		return r.Execute(ctx, GojaCall{
			PluginID: "boom",
			Source:   `(function(){ return wisp.fs.read({path:"a"}) })()`,
			Surface:  NewC24Surface("boom", "js.read", []string{"fs.read"}, crashing, WithC24Logf(quiet.Logf)),
		})
	}
	callAlive := func() (GojaResult, error) {
		return r.Execute(ctx, GojaCall{
			PluginID: "alive",
			Source:   `(function(){ return JSON.stringify(wisp.sysinfo.get()) })()`,
			Surface:  NewC24Surface("alive", "js.sysinfo", []string{"sysinfo"}, healthy, WithC24Logf(quiet.Logf)),
		})
	}

	// Sequential: before, across, and after the crash.
	if out, err := callAlive(); err != nil || out.Value != `{"ok":true}` {
		t.Fatalf("healthy neighbour before the crash: value=%q err=%v", out.Value, err)
	}
	if _, err := callBoom(); err == nil {
		t.Fatalf("the crashing plugin must fail")
	}
	if out, err := callAlive(); err != nil || out.Value != `{"ok":true}` {
		t.Fatalf("healthy neighbour after the crash: value=%q err=%v", out.Value, err)
	}

	// Concurrent: crash calls and healthy calls interleaved on one runtime.
	var wg sync.WaitGroup
	bad := make(chan string, 64)
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < 3; k++ {
				out, err := callAlive()
				if err != nil {
					bad <- "healthy call failed during concurrent crashes: " + err.Error()
					return
				}
				if !strings.Contains(out.Value, `"ok":true`) {
					bad <- "healthy value during concurrent crashes: " + out.Value
					return
				}
			}
		}()
	}
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := callBoom(); err == nil {
				bad <- "a crashing host call reported success"
			}
		}()
	}
	wg.Wait()
	close(bad)
	for msg := range bad {
		t.Error(msg)
	}

	if n := reg.Count(); n != 0 {
		t.Errorf("goroutines left behind after concurrent crashes: %d", n)
	}
	if out, err := callAlive(); err != nil || out.Value != `{"ok":true}` {
		t.Errorf("the runtime stopped serving a healthy plugin: value=%q err=%v", out.Value, err)
	}
}

// routedHost answers per plugin so one loader (which takes one host) can serve a
// crashing plugin and a healthy one at the same time.
type routedHost struct {
	mu    sync.Mutex
	calls []C24Request
}

func (h *routedHost) CallC24(_ context.Context, req C24Request) (C24Reply, error) {
	h.mu.Lock()
	h.calls = append(h.calls, req)
	h.mu.Unlock()
	switch req.PluginID {
	case "boom":
		panic("the host implementation crashed")
	case "alive":
		return C24Reply{Value: map[string]any{"ok": true}}, nil
	}
	return C24Reply{Refusal: "unknown plugin"}, nil
}

// TestC24HostCrashThroughTheLoaderKeepsTheGateServing is the AC#4 claim at the
// hop a real caller uses: the enable gate hands out one runtime envelope per
// plugin, the loader builds one surface per call, and a host crash in one plugin
// neither closes the gate nor disturbs the other plugin behind it.
func TestC24HostCrashThroughTheLoaderKeepsTheGateServing(t *testing.T) {
	t.Parallel()
	reg := observe.NewRegistry()
	gate := NewTier2Gate(Tier2Policy{
		Tier2Enabled: true,
		Plugins: map[string]Tier2PluginPolicy{
			"boom":  {Enabled: true},
			"alive": {Enabled: true},
		},
	}, WithTier2GateOptions(WithGojaRegistry(reg)))
	if err := gate.Enable(func() bool { return true }); err != nil {
		t.Fatalf("enable: %v", err)
	}

	source := func(_ context.Context, pluginID string) (string, error) {
		if pluginID == "boom" {
			return `function js_read(){ return wisp.fs.read({path:"a"}) }`, nil
		}
		return `function js_ping(){ return JSON.stringify(wisp.sysinfo.get()) }`, nil
	}
	log := &c24Log{}
	ex, err := NewTier2Executor(gate, &routedHost{}, source, WithTier2ExecutorLogf(log.Logf))
	if err != nil {
		t.Fatalf("build executor: %v", err)
	}
	ctx := context.Background()

	out, err := ex.Execute(ctx, Tier2ToolRequest{
		PluginID: "boom", Tool: "js.read", Declared: []string{"fs.read"},
	})
	if err == nil {
		t.Fatalf("a crashing host must fail the tool call, got %+v", out)
	}
	if class := observe.ClassOfOrInternal(err); class != observe.ClassTool {
		t.Errorf("loader error class = %q, want tool", class)
	}
	if !strings.Contains(err.Error(), C24FSRead) {
		t.Errorf("the tool error does not name the crashed entry %s: %q", C24FSRead, err.Error())
	}
	if n := log.CountContaining("C24-HOSTPANIC"); n != 1 {
		t.Errorf("C24-HOSTPANIC lines = %d, want 1; log:\n%s", n, log.Joined())
	}

	// The gate and the runtime envelope are still serving the other plugin.
	good, err := ex.Execute(ctx, Tier2ToolRequest{
		PluginID: "alive", Tool: "js.ping", Declared: []string{"sysinfo"},
	})
	if err != nil {
		t.Fatalf("the neighbour plugin was affected by the crash: %v", err)
	}
	if !strings.Contains(good.Text, `"ok":true`) {
		t.Errorf("neighbour value = %q, want the host's answer", good.Text)
	}

	// A second crash is contained the same way: no state accumulates per crash.
	if _, err := ex.Execute(ctx, Tier2ToolRequest{
		PluginID: "boom", Tool: "js.read", Declared: []string{"fs.read"},
	}); err == nil {
		t.Errorf("the second crashing call reported success")
	}
	if n := log.CountContaining("C24-HOSTPANIC"); n != 2 {
		t.Errorf("C24-HOSTPANIC lines after two crashes = %d, want 2", n)
	}

	if snap := gate.Snapshot(); len(snap.Runtimes) != 2 {
		t.Errorf("gate runtimes = %v, want both plugins still served", snap.Runtimes)
	}
	if n := reg.Count(); n != 0 {
		t.Errorf("goroutines left behind: %d", n)
	}
}
