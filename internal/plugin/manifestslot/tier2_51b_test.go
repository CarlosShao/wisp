package manifestslot

// Ticket 51 block B - the loader's Tier-2 arm.
//
// The Manifest slot already registers a Tier-1 plugin's declarative tools
// (ticket 50). A Tier-2 plugin has the same manifest face, but its tools run
// inside goja through the frozen C24 surface, so the slot's Run seam has to be
// filled by the C24 executor instead of the declarative runner. Tier2Run is
// that one-line adapter; this test proves the composition actually runs, and
// that the capability face the slot publishes is the manifest's own (the
// bridge re-checks it per call - C3 is not ours to enforce here).

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/tools"
)

// tier2Recorder is the host side; it answers one canned reply per entry.
type tier2Recorder struct {
	calls []plugin.C24Request
}

func (h *tier2Recorder) CallC24(_ context.Context, req plugin.C24Request) (plugin.C24Reply, error) {
	h.calls = append(h.calls, req)
	return plugin.C24Reply{Value: map[string]any{"data": "hello"}}, nil
}

// TestTier2RunFeedsTheManifestSlot is the wiring reading: a Tier-2 manifest
// becomes a registry entry whose Execute reaches the JS bundle through the C24
// surface, and the registered Decl carries the manifest's declarations so the
// bridge's C3 check sees the same face.
func TestTier2RunFeedsTheManifestSlot(t *testing.T) {
	t.Parallel()
	reg := observe.NewRegistry()
	gate := plugin.NewTier2Gate(plugin.Tier2Policy{
		Tier2Enabled: true,
		Plugins:      map[string]plugin.Tier2PluginPolicy{"com.example.js": {Enabled: true}},
	}, plugin.WithTier2GateOptions(plugin.WithGojaRegistry(reg)))
	if err := gate.Enable(func() bool { return true }); err != nil {
		t.Fatalf("enable: %v", err)
	}

	host := &tier2Recorder{}
	ex, err := plugin.NewTier2Executor(gate, host, func(context.Context, string) (string, error) {
		return `function js_echo(params){ return wisp.fs.read({path: params.path}) }`, nil
	})
	if err != nil {
		t.Fatalf("executor: %v", err)
	}

	m := plugin.Manifest{
		SchemaVersion: plugin.CurrentSchemaVersion,
		ID:            "com.example.js",
		HostAPI:       "^1.0.0",
		Name:          "JS Tools",
		Version:       "0.1.0",
		Tools: []plugin.ToolDef{{
			Name:         "js.echo",
			Description:  "Read a file through the C24 surface.",
			Risk:         plugin.RiskL0,
			Capabilities: []string{"fs.read"},
			Parameters:   map[string]any{"type": "object"},
		}},
	}
	provider, err := New(m, Tier2Run(ex))
	if err != nil {
		t.Fatalf("provider: %v", err)
	}
	if err := tools.NewRegistry().RegisterProvider(provider); err != nil {
		t.Fatalf("register: %v", err)
	}

	entry, ok := provider.Lookup("js.echo")
	if !ok {
		t.Fatalf("the Tier-2 tool was not registered")
	}
	if len(entry.Decl.Needs) != 1 || entry.Decl.Needs[0] != tools.CapFSRead {
		t.Fatalf("Decl.Needs = %v, want the manifest's [fs.read]", entry.Decl.Needs)
	}
	if entry.Decl.PluginID != "com.example.js" {
		t.Fatalf("entry lost its plugin id: %q", entry.Decl.PluginID)
	}

	res, err := entry.Tool.Execute(context.Background(), json.RawMessage(`{"path":"D:/work/notes.txt"}`), nil)
	if err != nil {
		t.Fatalf("tool execute: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool reported an error: %q", res.Text)
	}
	if !strings.Contains(res.Text, `"data":"hello"`) {
		t.Fatalf("tool text = %q, want the host reply", res.Text)
	}
	if len(host.calls) != 1 || host.calls[0].Entry != "wisp.fs.read" {
		t.Fatalf("host calls = %+v, want one wisp.fs.read", host.calls)
	}
	if n := reg.Count(); n != 0 {
		t.Fatalf("watchdog goroutine leaked: %d", n)
	}
}
