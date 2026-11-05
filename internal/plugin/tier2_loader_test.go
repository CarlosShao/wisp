package plugin

// Ticket 51 block B - the loader wiring tests.
//
// The C24 surface is only reachable from JS; what a real caller needs is the
// hop from "one manifest tool was called" to "one JS function ran with the
// plugin's own grants". That hop is Tier2Executor, and these tests are its
// only door. The gate is the SAME Tier2Gate ticket 51's AC#5 shipped: this file
// proves the gate's default-off / L2 / per-plugin rules still stand in front of
// the new surface, not beside it.

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// openTier2Gate builds the state a host has after the owner answered the L2
// card: policy on, this plugin's key on, the verdict recorded.
func openTier2Gate(t *testing.T) (*Tier2Gate, *observe.Registry) {
	t.Helper()
	reg := observe.NewRegistry()
	g := NewTier2Gate(Tier2Policy{
		Tier2Enabled: true,
		Plugins:      map[string]Tier2PluginPolicy{"demo": {Enabled: true}},
	}, WithTier2GateOptions(WithGojaRegistry(reg)))
	if err := g.Enable(func() bool { return true }); err != nil {
		t.Fatalf("enable: %v", err)
	}
	return g, reg
}

// bundleSource is the composition root's read of the plugin's entry file. It
// counts calls so "a refused call never read the plugin's code" is measurable.
func bundleSource(bundle string, calls *int) Tier2BundleSource {
	return func(context.Context, string) (string, error) {
		*calls++
		return bundle, nil
	}
}

// TestTier2ExecutorNeedsEveryWiring: a loader that is missing a gate, a host or
// a source is an assembly bug, refused at construction (fail closed, named).
func TestTier2ExecutorNeedsEveryWiring(t *testing.T) {
	t.Parallel()
	gate, _ := openTier2Gate(t)
	src := bundleSource("function echo(){}", new(int))
	host := &c24Recorder{}

	if _, err := NewTier2Executor(nil, host, src); err == nil ||
		observe.ClassOfOrInternal(err) != observe.ClassConfig {
		t.Fatalf("nil gate: err = %v, want config", err)
	}
	if _, err := NewTier2Executor(gate, nil, src); err == nil ||
		observe.ClassOfOrInternal(err) != observe.ClassConfig {
		t.Fatalf("nil host: err = %v, want config", err)
	}
	if _, err := NewTier2Executor(gate, host, nil); err == nil ||
		observe.ClassOfOrInternal(err) != observe.ClassConfig {
		t.Fatalf("nil source: err = %v, want config", err)
	}
	if _, err := NewTier2Executor(gate, host, src); err != nil {
		t.Fatalf("a fully wired executor was refused: %v", err)
	}
}

// TestTier2ExecutorRunsATier2ToolEndToEnd is the wiring claim: manifest params
// arrive as the JS function's argument, the entry sees only its declared
// grants, the host is called through the C24 surface, and the host's answer is
// what the tool returns.
func TestTier2ExecutorRunsATier2ToolEndToEnd(t *testing.T) {
	t.Parallel()
	gate, reg := openTier2Gate(t)
	host := &c24Recorder{reply: func(C24Request) C24Reply {
		return C24Reply{Value: map[string]any{"data": "hello"}}
	}}
	reads := 0
	ex, err := NewTier2Executor(gate, host, bundleSource(
		`function echo(params){ return wisp.fs.read({path: params.path}) }`, &reads))
	if err != nil {
		t.Fatalf("executor: %v", err)
	}

	out, err := ex.Execute(context.Background(), Tier2ToolRequest{
		PluginID: "demo",
		Tool:     "echo",
		Declared: []string{"fs.read"},
		Params:   json.RawMessage(`{"path":"D:/work/notes.txt"}`),
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if out.IsError {
		t.Fatalf("call reported an error outcome: %q", out.Text)
	}
	if out.Text != `{"data":"hello"}` {
		t.Fatalf("tool text = %q, want the host reply rendered as JSON", out.Text)
	}
	calls := host.Calls()
	if len(calls) != 1 {
		t.Fatalf("host saw %d call(s), want 1", len(calls))
	}
	if calls[0].Tool != "echo" || calls[0].PluginID != "demo" {
		t.Fatalf("host call lost its owner: %+v", calls[0])
	}
	wantArgs := map[string]any{"path": "D:/work/notes.txt"}
	if !jsonEqual(calls[0].Args, wantArgs) {
		t.Fatalf("host args = %#v, want %#v", calls[0].Args, wantArgs)
	}
	if n := reg.Count(); n != 0 {
		t.Fatalf("watchdog goroutine leaked: %d", n)
	}
}

// TestTier2ExecutorRefusesBeforeItReadsTheCode: the gate is in FRONT of the
// loader, so a refused call reads no bundle and calls no host. This is the
// difference between "tier 2 is off" and "the plugin ran and was told no".
func TestTier2ExecutorRefusesBeforeItReadsTheCode(t *testing.T) {
	t.Parallel()
	host := &c24Recorder{}
	reads := 0
	src := bundleSource(`function echo(){ return "ran" }`, &reads)

	cases := []struct {
		name string
		gate *Tier2Gate
		want string
	}{
		{
			name: "default off",
			gate: NewTier2Gate(Tier2Policy{}, WithTier2GateOptions(WithGojaRegistry(observe.NewRegistry()))),
			want: "tier2_enabled",
		},
		{
			name: "policy on but no L2 verdict",
			gate: NewTier2Gate(Tier2Policy{
				Tier2Enabled: true,
				Plugins:      map[string]Tier2PluginPolicy{"demo": {Enabled: true}},
			}, WithTier2GateOptions(WithGojaRegistry(observe.NewRegistry()))),
			want: "L2",
		},
		{
			name: "per-plugin switch off",
			gate: func() *Tier2Gate {
				g := NewTier2Gate(Tier2Policy{
					Tier2Enabled: true,
					Plugins:      map[string]Tier2PluginPolicy{"demo": {Enabled: false}},
				}, WithTier2GateOptions(WithGojaRegistry(observe.NewRegistry())))
				if err := g.Enable(func() bool { return true }); err != nil {
					t.Fatalf("enable: %v", err)
				}
				return g
			}(),
			want: "[plugins.demo].enabled",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ex, err := NewTier2Executor(tc.gate, host, src)
			if err != nil {
				t.Fatalf("executor: %v", err)
			}
			_, err = ex.Execute(context.Background(), Tier2ToolRequest{
				PluginID: "demo", Tool: "echo", Declared: []string{"fs.read"},
				Params: json.RawMessage(`{}`),
			})
			if err == nil {
				t.Fatalf("a refused call succeeded")
			}
			if class := observe.ClassOfOrInternal(err); class != observe.ClassPermissionDenied {
				t.Fatalf("class = %q, want permission_denied", class)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("refusal %q does not name %q", err, tc.want)
			}
		})
	}
	if reads != 0 {
		t.Fatalf("a refused call read the bundle %d time(s)", reads)
	}
	if host.Count() != 0 {
		t.Fatalf("a refused call reached the host")
	}
}

// TestTier2ExecutorMissingEntryIsAToolError: the tool name is the JS entry
// (block A's Function seam); a bundle that does not expose it is the plugin's
// bug, returned as a tool-class error the model can read.
func TestTier2ExecutorMissingEntryIsAToolError(t *testing.T) {
	t.Parallel()
	gate, _ := openTier2Gate(t)
	host := &c24Recorder{}
	ex, err := NewTier2Executor(gate, host, bundleSource(`var notATool = 1`, new(int)))
	if err != nil {
		t.Fatalf("executor: %v", err)
	}
	_, err = ex.Execute(context.Background(), Tier2ToolRequest{
		PluginID: "demo", Tool: "echo", Declared: []string{"fs.read"},
		Params: json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatalf("missing entry was accepted")
	}
	if class := observe.ClassOfOrInternal(err); class != observe.ClassTool {
		t.Fatalf("class = %q, want tool", class)
	}
	if !strings.Contains(err.Error(), "echo") {
		t.Fatalf("the error does not name the missing entry: %v", err)
	}
}

// TestTier2ExecutorUndeclaredCapabilityStaysUndefined: the loader passes the
// manifest's declaration through, and the surface refuses the entry the
// manifest did not declare - the plugin cannot grant itself a capability.
func TestTier2ExecutorUndeclaredCapabilityStaysUndefined(t *testing.T) {
	t.Parallel()
	gate, _ := openTier2Gate(t)
	host := &c24Recorder{}
	ex, err := NewTier2Executor(gate, host, bundleSource(
		`function echo(){ return typeof wisp.fs.write }`, new(int)))
	if err != nil {
		t.Fatalf("executor: %v", err)
	}
	out, err := ex.Execute(context.Background(), Tier2ToolRequest{
		PluginID: "demo", Tool: "echo", Declared: []string{"fs.read"},
		Params: json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if got := strings.Trim(out.Text, `"`); got != "undefined" {
		t.Fatalf("undeclared entry = %q, want undefined", got)
	}
	if host.Count() != 0 {
		t.Fatalf("an undeclared entry reached the host")
	}
}

// TestTier2ExecutorEmptyIdentityIsRefused: the audit row, the gate key and the
// surface are all keyed on plugin id + tool name, so a call missing one is a
// config error rather than a call attributed to nobody.
func TestTier2ExecutorEmptyIdentityIsRefused(t *testing.T) {
	t.Parallel()
	gate, _ := openTier2Gate(t)
	ex, err := NewTier2Executor(gate, &c24Recorder{}, bundleSource("function echo(){}", new(int)))
	if err != nil {
		t.Fatalf("executor: %v", err)
	}
	for _, req := range []Tier2ToolRequest{
		{Tool: "echo"},
		{PluginID: "demo"},
	} {
		if _, err := ex.Execute(context.Background(), req); err == nil ||
			observe.ClassOfOrInternal(err) != observe.ClassConfig {
			t.Fatalf("req %+v: err = %v, want config", req, err)
		}
	}
}

// TestTier2EntryNameIsTheDottedToolNameSanitized pins the JS-entry convention:
// a C1 tool name is namespaced and a JS identifier cannot hold a dot, so the
// loader derives one deterministic name per tool. Without a pinned mapping the
// manifest and the bundle would each guess.
func TestTier2EntryNameIsTheDottedToolNameSanitized(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"echo":    "echo",
		"js.echo": "js_echo",
		"a.b.c":   "a_b_c",
		"9lives":  "_9lives",
		"$ok":     "$ok",
	}
	for in, want := range cases {
		if got := Tier2EntryName(in); got != want {
			t.Errorf("Tier2EntryName(%q) = %q, want %q", in, got, want)
		}
	}
}

// jsonEqual compares two decoded JSON values without depending on map ordering.
func jsonEqual(a, b any) bool {
	ab, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(ab) == string(bb)
}
