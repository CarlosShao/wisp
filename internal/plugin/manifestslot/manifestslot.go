// Package manifestslot is the C4 Manifest slot (internal/tools.Provider) over
// the declarative Tier-1 tools of ONE installed plugin.
//
// Why this is a leaf package and not a file inside internal/plugin: the
// dependency direction is fixed by what already exists. internal/tools pulls in
// internal/agent, which pulls internal/memory, which imports internal/plugin
// (memory/retention.go). So internal/plugin MUST NOT import internal/tools -
// that is an import cycle, not a style choice. The adapter therefore sits above
// both and imports both, and neither package learns about the other:
//
//	plugin (manifest grammar, C14/C2/C3 declarations)
//	   ^                                          ^
//	   |                                          |
//	manifestslot  ----> tools (Registry, Decl, Result, Provider)
//
// Nothing here invents a seam. tools.Provider (Kind/List/Lookup), tools.Entry
// (Tool + Decl), tools.KindManifest and the Decl sidecar are the shape
// internal/tools/registry.go already declares, and its own comment says the
// host-side facts of a plugin tool live in Decl: "read from the manifest by the
// Tier-1/Tier-2 providers (tickets 50/51)".
//
// Scope (ticket 50 AC6 only): registration and identity. It is not an installer
// (AC3 stays open), it does not re-validate a manifest (plugin/validator.go owns
// C14, plugin/tamper.go owns the hash), and it cannot execute anything by
// itself, because a Tier-1 tool is declarative: the behaviour arrives as a Run.
package manifestslot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Run executes one declarative Tier-1 call for the plugin that owns it. This is
// the seam the later tickets fill (ticket 51's goja host, ticket 52's D46
// command wrapper). The plugin id travels with the call, so one Run can never
// serve two plugins without knowing which one it is answering.
//
// C3 is NOT checked here. The bridge re-resolves the declared capability face
// per call ahead of Execute (tools.Bridge.checkCaps), because the host's view can
// shrink after registration. That is why Needs below is the manifest's own face
// and nothing more.
type Run func(ctx context.Context, pluginID string, def plugin.ToolDef, params json.RawMessage) (tools.Result, error)

// Provider is the Manifest slot over one plugin's declarative tools.
type Provider struct {
	pluginID string
	entries  []tools.Entry
	byName   map[string]tools.Entry
}

var _ tools.Provider = (*Provider)(nil)

// New turns a parsed manifest into registry entries.
//
// The manifest must carry its id: a plugin tool with no owner is exactly the
// provenance gap AC6 closes, so this refuses instead of registering anonymous
// entries. run must not be nil either - a slot that can list but never execute
// would make the tool visible to the model and then fail at the first call.
func New(m plugin.Manifest, run Run) (*Provider, error) {
	if m.ID == "" {
		return nil, fmt.Errorf("manifestslot: empty plugin id (AC6 requires every plugin tool to name its owner)")
	}
	if run == nil {
		return nil, fmt.Errorf("manifestslot: plugin %q: nil Run", m.ID)
	}
	p := &Provider{pluginID: m.ID, byName: map[string]tools.Entry{}}
	for _, def := range m.Tools {
		e, err := entryFor(m.ID, def, run)
		if err != nil {
			return nil, err
		}
		p.entries = append(p.entries, e)
		p.byName[def.Name] = e
	}
	return p, nil
}

// Kind names the C4 slot this provider serves.
func (p *Provider) Kind() tools.Kind { return tools.KindManifest }

// PluginID is the plugin this slot is the view of.
func (p *Provider) PluginID() string { return p.pluginID }

// List implements tools.Provider in the manifest's declaration order, which is
// also the order a capability review shows the user.
func (p *Provider) List() []tools.Entry {
	out := make([]tools.Entry, len(p.entries))
	copy(out, p.entries)
	return out
}

// Lookup implements tools.Provider.
func (p *Provider) Lookup(name string) (tools.Entry, bool) {
	e, ok := p.byName[name]
	return e, ok
}

// entryFor reads one ToolDef into one registry entry. The values come ONLY from
// the manifest: the loader never raises a risk level (C2: a declaration is the
// R1 lower bound, C19 owns the verdict) and never widens a capability face.
func entryFor(pluginID string, def plugin.ToolDef, run Run) (tools.Entry, error) {
	level, err := declaredLevel(def.Risk)
	if err != nil {
		return tools.Entry{}, fmt.Errorf("manifestslot: plugin %q tool %q: %w", pluginID, def.Name, err)
	}
	caps := make([]tools.Capability, 0, len(def.Capabilities))
	for _, c := range def.Capabilities {
		cap := tools.Capability(c)
		if !cap.Valid() {
			return tools.Entry{}, fmt.Errorf(
				"manifestslot: plugin %q tool %q: capability %q is outside the frozen C3 set", pluginID, def.Name, c)
		}
		caps = append(caps, cap)
	}
	if def.Parameters == nil {
		return tools.Entry{}, fmt.Errorf(
			"manifestslot: plugin %q tool %q: no parameter schema (C14)", pluginID, def.Name)
	}
	schema, err := json.Marshal(def.Parameters)
	if err != nil {
		return tools.Entry{}, fmt.Errorf("manifestslot: plugin %q tool %q: parameter schema: %w", pluginID, def.Name, err)
	}
	return tools.Entry{
		Tool: &declTool{pluginID: pluginID, def: def, run: run, schema: schema},
		Decl: tools.Decl{
			Capabilities: caps,
			Needs:        caps,
			Declared:     level,
			Provider:     tools.KindManifest,
			PluginID:     pluginID,
		},
	}, nil
}

// declaredLevel maps the C2 spelling to the host's risk.Level. An unknown
// spelling is an error rather than the safest-looking default, so a typo cannot
// silently register a dangerous tool as L0.
func declaredLevel(s string) (risk.Level, error) {
	switch s {
	case plugin.RiskL0:
		return risk.L0, nil
	case plugin.RiskL1:
		return risk.L1, nil
	case plugin.RiskL2:
		return risk.L2, nil
	}
	return risk.L0, fmt.Errorf("risk %q is not one of %s/%s/%s (C2)",
		s, plugin.RiskL0, plugin.RiskL1, plugin.RiskL2)
}

// declTool is one declarative Tier-1 tool as the registry sees it: the four C1
// methods, with Execute handed to the plugin's Run.
type declTool struct {
	pluginID string
	def      plugin.ToolDef
	run      Run
	schema   json.RawMessage
}

var _ tools.Tool = (*declTool)(nil)

// Name implements C1. Global uniqueness is the Registry's job, which is what
// stops a plugin from shadowing a builtin or another plugin's tool.
func (t *declTool) Name() string { return t.def.Name }

// Description implements C1: the model-facing text, verbatim from the manifest.
func (t *declTool) Description() string { return t.def.Description }

// Parameters implements C1: the manifest's JSON Schema, marshaled once here.
func (t *declTool) Parameters() tools.JSONSchema { return t.schema }

// Execute implements C1 by handing the call to the plugin runner with its own
// plugin id, so the seam can never attribute work to the wrong owner.
func (t *declTool) Execute(ctx context.Context, params json.RawMessage, _ func(delta string)) (tools.Result, error) {
	return t.run(ctx, t.pluginID, t.def, params)
}
