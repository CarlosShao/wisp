package manifestslot

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ac6Manifest is a schema-valid Tier-1 manifest (C14) whose single declarative
// tool is named in the C1 'namespace.action' shape the Registry enforces.
const ac6Manifest = `
schema_version = 1
id = "com.example.lark"
host_api = "^1.0.0"
name = "Lark Tools"
version = "0.1.0"

[[tools]]
name = "lark.search"
description = "Search the connected workspace."
risk = "L0"
capabilities = ["net"]
net_allowlist = ["api.example.com"]

[tools.parameters]
type = "object"
[tools.parameters.properties]
query.type = "string"
`

func nopRun(context.Context, string, plugin.ToolDef, json.RawMessage) (tools.Result, error) {
	return tools.Result{}, nil
}

// TestAC6PluginToolsEnterTheRegistry is the visibility half of AC6: parse a
// manifest, register its slot, and the plugin tool is on the same directory
// list_tools reads - carrying the manifest's risk floor, its capability face and
// the owning plugin id.
func TestAC6PluginToolsEnterTheRegistry(t *testing.T) {
	m, err := plugin.ParseManifest([]byte(ac6Manifest))
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	var calls int
	provider, err := New(m, func(_ context.Context, pluginID string, def plugin.ToolDef, params json.RawMessage) (tools.Result, error) {
		calls++
		if pluginID != m.ID {
			t.Errorf("runner saw plugin id %q, want %q", pluginID, m.ID)
		}
		if def.Name != "lark.search" {
			t.Errorf("runner saw tool %q", def.Name)
		}
		return tools.Result{Text: "ran " + string(params)}, nil
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if provider.Kind() != tools.KindManifest {
		t.Errorf("Kind() = %q, want %q", provider.Kind(), tools.KindManifest)
	}
	reg := tools.NewRegistry()
	if err := reg.RegisterProvider(provider); err != nil {
		t.Fatalf("RegisterProvider: %v", err)
	}

	var names []string
	for _, e := range reg.List() {
		names = append(names, e.Tool.Name())
	}
	if strings.Join(names, ",") != "lark.search" {
		t.Fatalf("Registry.List() = %v, want the plugin tool visible to list_tools", names)
	}

	e, ok := reg.Lookup("lark.search")
	if !ok {
		t.Fatal("Lookup('lark.search') not found after registration")
	}
	if e.Decl.Provider != tools.KindManifest {
		t.Errorf("Decl.Provider = %q, want %q", e.Decl.Provider, tools.KindManifest)
	}
	if e.Decl.PluginID != "com.example.lark" {
		t.Errorf("Decl.PluginID = %q, want com.example.lark", e.Decl.PluginID)
	}
	if len(e.Decl.Capabilities) != 1 || e.Decl.Capabilities[0] != tools.CapNet {
		t.Errorf("Decl.Capabilities = %v, want [net] read from the manifest", e.Decl.Capabilities)
	}
	if len(e.Decl.Needs) != 1 || e.Decl.Needs[0] != tools.CapNet {
		t.Errorf("Decl.Needs = %v, want the declared face and nothing wider", e.Decl.Needs)
	}
	if got := e.Tool.Description(); got != "Search the connected workspace." {
		t.Errorf("Description() = %q", got)
	}
	var schema map[string]any
	if err := json.Unmarshal(e.Tool.Parameters(), &schema); err != nil {
		t.Fatalf("Parameters() is not JSON: %v", err)
	}
	if schema["type"] != "object" {
		t.Errorf("parameter schema type = %v, want object", schema["type"])
	}
	res, err := e.Tool.Execute(t.Context(), json.RawMessage(`{"query":"hi"}`), nil)
	if err != nil || res.IsError {
		t.Fatalf("Execute: err=%v res=%+v", err, res)
	}
	if calls != 1 {
		t.Fatalf("runner calls = %d, want 1 (Execute delegates to the seam)", calls)
	}
}

// TestAC6UnregisteredPluginToolIsInvisible is the reverse check: registering is
// what makes a plugin tool visible, so an installed-but-unregistered plugin's
// tool is absent from the directory and uncallable.
func TestAC6UnregisteredPluginToolIsInvisible(t *testing.T) {
	m, err := plugin.ParseManifest([]byte(ac6Manifest))
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	if _, err := New(m, nopRun); err != nil {
		t.Fatal(err)
	}
	reg := tools.NewRegistry()
	if _, ok := reg.Lookup("lark.search"); ok {
		t.Fatal("a plugin tool that was never registered is visible in the directory")
	}
	if n := len(reg.List()); n != 0 {
		t.Fatalf("Registry.List() = %d entries, want 0 before any slot lands", n)
	}
	entry, err := entryFor("com.example.lark", m.Tools[0], nopRun)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Lookup(entry.Tool.Name()); ok {
		t.Fatal("building an entry made it visible without registration")
	}
}

// TestAC6RefusesAnOwnerlessSlot closes the identity from the other side: no id,
// no runner, no slot.
func TestAC6RefusesAnOwnerlessSlot(t *testing.T) {
	m := plugin.Manifest{Tools: []plugin.ToolDef{{
		Name: "lark.search", Description: "d", Risk: plugin.RiskL0,
		Capabilities: []string{"net"}, Parameters: map[string]any{"type": "object"},
	}}}
	if _, err := New(m, nopRun); err == nil {
		t.Fatal("New accepted a manifest with no plugin id")
	} else if !strings.Contains(err.Error(), "plugin id") {
		t.Fatalf("err = %v, want it to name the missing plugin id", err)
	}
	if _, err := New(plugin.Manifest{ID: "com.example.lark"}, nil); err == nil {
		t.Fatal("New accepted a nil Run")
	}
}

// TestAC6RefusesADeclarationOutsideTheContract proves the shim cannot launder a
// manifest the validator would reject: an unknown capability token, an unknown
// risk spelling or a missing parameter schema is an error at slot-build time,
// never a silent default to the safest-looking value.
func TestAC6RefusesADeclarationOutsideTheContract(t *testing.T) {
	cases := []struct {
		name string
		def  plugin.ToolDef
		want string
	}{
		{"unknown capability", plugin.ToolDef{
			Name: "lark.search", Description: "d", Risk: plugin.RiskL0,
			Capabilities: []string{"teleport"}, Parameters: map[string]any{"type": "object"},
		}, "teleport"},
		{"unknown risk", plugin.ToolDef{
			Name: "lark.search", Description: "d", Risk: "L9",
			Capabilities: []string{"net"}, Parameters: map[string]any{"type": "object"},
		}, "L9"},
		{"no parameter schema", plugin.ToolDef{
			Name: "lark.search", Description: "d", Risk: plugin.RiskL0,
			Capabilities: []string{"net"},
		}, "C14"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(plugin.Manifest{ID: "com.x", Tools: []plugin.ToolDef{c.def}}, nopRun)
			if err == nil {
				t.Fatalf("slot accepted %s", c.name)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, want it to name %q", err, c.want)
			}
		})
	}
}

// TestAC6ToolNameShapeIsTheRegistrysRuleNotTheShims records the seam that AC3
// still has to close: the C14 manifest grammar allows a multi-dot tool name, the
// C1 registry rule does not, so such a manifest is refused at registration
// rather than silently registered or silently truncated.
func TestAC6ToolNameShapeIsTheRegistrysRuleNotTheShims(t *testing.T) {
	def := plugin.ToolDef{
		Name: "com.example.lark.search", Description: "d", Risk: plugin.RiskL0,
		Capabilities: []string{"net"}, Parameters: map[string]any{"type": "object"},
	}
	p, err := New(plugin.Manifest{ID: "com.example.lark", Tools: []plugin.ToolDef{def}}, nopRun)
	if err != nil {
		t.Fatalf("the shim itself does not judge name grammar: %v", err)
	}
	if err := tools.NewRegistry().RegisterProvider(p); err == nil {
		t.Fatal("the Registry accepted a tool name outside C1's 'namespace.action'")
	} else if !strings.Contains(err.Error(), "more than one") {
		t.Fatalf("err = %v, want the C1 name-shape refusal", err)
	}
}
