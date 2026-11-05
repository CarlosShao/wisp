package plugin

// Tier-2 loader wiring - ticket 51, BLOCK B.
//
// Block A is the envelope (one fresh VM per call, wall clock, result cap,
// recover). Block B's C24 surface is only reachable from inside a VM. This file
// is the hop between them and the only production caller of Tier2Gate today:
//
//	manifest tool call
//	   -> Tier2Executor.Execute
//	      -> Tier2Gate.Execute           (AC#5: default off / L2 / per-plugin)
//	         -> GojaRuntime.Execute      (block A: clock, cap, recover)
//	            -> installed C24 surface (block B: 9 entries, 6 groups)
//	               -> C24Host.CallC24    (host: C3 bridge, C19, taint, audit)
//
// What the loader does NOT do, and must not: it neither judges risk nor checks
// capabilities. Both belong to the bridge, which already re-resolves the
// manifest's declared face per call (tools.Bridge.checkCaps) and owns the
// tool_call row. Doing either here would create the second dispatch path the
// bridge exists to prevent. The executor's job is to bring the four facts the
// runtime needs - plugin id, tool name, the declaration, the bundle - and to
// hand back what the VM produced.
//
// The bundle source is the composition root's read of the plugin's entry file
// (SPEC-07 §4's Tier-2 `entry`). It is a seam rather than a path literal because
// the plugins directory belongs to the installer/lifecycle, not to this package.
//
// The JS entry convention: a Tier-2 bundle exposes one top-level callable per
// declared tool, named by Tier2EntryName (below). The mapping exists because C1
// tool names are namespaced ("js.echo") and a JS identifier cannot contain a
// dot; sanitizing is the same rule block A uses for its watchdog name, so one
// spelling serves the audit line and the entry alike.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Tier2EntryName maps a manifest tool name onto its JS entry function name:
// every rune outside [A-Za-z0-9_$] becomes '_' ("js.echo" -> "js_echo"), and a
// leading digit is prefixed so the result is always a legal identifier.
func Tier2EntryName(tool string) string {
	if tool == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(tool))
	for _, r := range tool {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '$':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	name := b.String()
	// A JS identifier may not start with a digit; a leading underscore keeps
	// the name callable instead of turning a legal tool name into a syntax
	// error at load time.
	if name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	}
	return name
}

// Tier2BundleSource returns one plugin's JS bundle text. An error means the
// bundle could not be read (missing entry file, I/O fault): the call is refused
// as a host-side config fault, because no plugin code can run without it.
type Tier2BundleSource func(ctx context.Context, pluginID string) (string, error)

// Tier2ToolRequest is one manifest tool call, as the loader receives it.
//
// Declared is the manifest's capability list for this tool - the same list that
// travels on tools.Decl.Needs to the bridge. The surface grants entries from it;
// it never widens it.
type Tier2ToolRequest struct {
	PluginID string
	Tool     string
	Declared []string
	Params   json.RawMessage
}

// Tier2ToolOutcome is the capped result of one Tier-2 call plus the two
// readings an operator needs: how big it really was, and whether the source
// touched anything outside the frozen surface.
type Tier2ToolOutcome struct {
	Text               string
	IsError            bool
	Truncated          bool
	FullBytes          int
	ReturnedBytes      int
	ContractViolations []string
}

// Tier2ExecutorOption configures a Tier2Executor.
type Tier2ExecutorOption func(*Tier2Executor)

// WithTier2ExecutorLogf forwards the audit sink to every surface the executor
// builds (default: each surface's own slog sink).
func WithTier2ExecutorLogf(f func(format string, args ...any)) Tier2ExecutorOption {
	return func(e *Tier2Executor) {
		if f != nil {
			e.logf = f
		}
	}
}

// Tier2Executor is the loader arm. It is immutable after construction and safe
// for concurrent use: the gate and the host are shared, every call builds its
// own surface and VM.
type Tier2Executor struct {
	gate   *Tier2Gate
	host   C24Host
	source Tier2BundleSource
	logf   func(format string, args ...any)
}

// NewTier2Executor requires all three wiring points. A half-wired executor is
// refused here rather than failing on the first call, and every refusal is a
// config-class error naming the missing piece (an assembly bug, not a plugin
// bug).
func NewTier2Executor(gate *Tier2Gate, host C24Host, source Tier2BundleSource, opts ...Tier2ExecutorOption) (*Tier2Executor, error) {
	if gate == nil {
		return nil, observe.New(observe.ClassConfig,
			"tier2 loader: no enable gate - every Tier-2 call must pass the default-off gate (ticket 51 AC#5)")
	}
	if host == nil {
		return nil, observe.New(observe.ClassConfig,
			"tier2 loader: no C24 host - a surface with nowhere to call is not a runtime")
	}
	if source == nil {
		return nil, observe.New(observe.ClassConfig,
			"tier2 loader: no bundle source - the loader cannot invent a plugin's code")
	}
	e := &Tier2Executor{gate: gate, host: host, source: source}
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}
	return e, nil
}

// Execute runs one manifest tool call through the gate, the runtime and the C24
// surface. A gate refusal returns the gate's own permission_denied error before
// the bundle is read; a JS/protocol failure returns a tool-class error with the
// runtime's text.
func (e *Tier2Executor) Execute(ctx context.Context, req Tier2ToolRequest) (Tier2ToolOutcome, error) {
	if e == nil {
		return Tier2ToolOutcome{}, observe.New(observe.ClassInternal, "tier2 loader: Execute on a nil executor")
	}
	if strings.TrimSpace(req.PluginID) == "" {
		return Tier2ToolOutcome{}, observe.New(observe.ClassConfig,
			"tier2 tool call carries no plugin id: the gate key, the audit row and the surface are keyed on it")
	}
	if strings.TrimSpace(req.Tool) == "" {
		return Tier2ToolOutcome{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"plugin %q: tier2 tool call carries no tool name (it is also the JS entry to call)", req.PluginID))
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// Gate FIRST: a refused call must not read the plugin's code at all, so
	// "tier 2 is off" means none of this plugin's JS ran (ticket 51 AC#5).
	if _, err := e.gate.Runtime(req.PluginID); err != nil {
		return Tier2ToolOutcome{}, err
	}

	source, err := e.source(ctx, req.PluginID)
	if err != nil {
		return Tier2ToolOutcome{}, observe.Wrap(observe.ClassConfig, err, fmt.Sprintf(
			"plugin %q: cannot read the tier-2 bundle", req.PluginID))
	}
	if strings.TrimSpace(source) == "" {
		return Tier2ToolOutcome{}, observe.New(observe.ClassConfig, fmt.Sprintf(
			"plugin %q: the tier-2 bundle is empty", req.PluginID))
	}

	args, err := tier2Params(req)
	if err != nil {
		return Tier2ToolOutcome{}, err
	}

	opts := []C24Option{}
	if e.logf != nil {
		opts = append(opts, WithC24Logf(e.logf))
	}
	surface := NewC24Surface(req.PluginID, req.Tool, req.Declared, e.host, opts...)

	out, err := e.gate.Execute(ctx, GojaCall{
		PluginID: req.PluginID,
		Source:   source,
		Function: Tier2EntryName(req.Tool),
		Args:     args,
		Surface:  surface,
	})
	res := Tier2ToolOutcome{
		Text:               out.Value,
		Truncated:          out.Truncated,
		FullBytes:          out.FullBytes,
		ReturnedBytes:      out.ReturnedBytes,
		ContractViolations: out.ContractViolations,
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

// tier2Params turns the manifest call's JSON parameters into the single JS
// argument the entry function receives. Empty params pass no argument (so a
// parameterless tool reads the same as a script).
func tier2Params(req Tier2ToolRequest) ([]any, error) {
	if len(req.Params) == 0 {
		return nil, nil
	}
	var decoded any
	if err := json.Unmarshal(req.Params, &decoded); err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err, fmt.Sprintf(
			"plugin %q tool %q: parameters are not JSON", req.PluginID, req.Tool))
	}
	return []any{decoded}, nil
}
