package plugin

// Tier-2 enable gate (ticket 51, acceptance criterion "Enable gate"):
// the one door a Tier-2 goja runtime can be reached through, and the three
// fail-closed rules behind it.
//
//  1. DEFAULT OFF, AND OFF MEANS NOTHING EXISTS. [plugins].tier2_enabled
//     ships false (SPEC-07 §4, D19). While the gate is off no GojaRuntime
//     object is built at all - not "one is built that refuses" - which is
//     D3's use-to-wake taken literally. Snapshot().Runtimes is the reading;
//     the factory counter in the tests is what tells the two states apart,
//     because a refusal alone cannot.
//
//  2. ENABLING TAKES THE L2 VERDICT, NOT A BOOLEAN. D36 rule 1 makes
//     [plugins] a locked section: turning tier 2 on is a loosening and the
//     config manager asks its L2 hook before writing the value
//     (internal/config/manager.go: planLocked -> ConfirmLocked). That covers
//     RELOADS only. The initial load raises no card, so a hand-edited
//     config.toml that says tier2_enabled = true loads as true - "the applied
//     config says on" is data, not an approval. Enable therefore asks the
//     confirmation hook itself: nil and false both deny, and a denial leaves
//     the gate (and the absence of any runtime) exactly where it was. The
//     hook is func() bool on purpose - how the card is raised is the caller's
//     business (panel today, SPEC-07 §6), and this package must not grow a
//     dependency on internal/panel or internal/agent to ask a human.
//     A tightening (tier2 off) CLEARS the record: a stale yes is not a yes.
//
//  3. PER-PLUGIN KEYS. [plugins.<id>].enabled is the second switch. Tier 2 on
//     plus a missing key, or enabled = false, is a refusal that names the key
//     path - never "on by default". ApplyPolicy drops the runtime of any
//     plugin that was just disabled or removed and leaves its neighbours'
//     envelopes alone.
//
// What this file is NOT:
//   - not the C24 GojaHostAPI surface (block B, whose initial set is frozen
//     only after owner approval - AGENTS.md §2 "undefined means stop");
//   - not a thread-affinity guard: Execute runs on the caller's goroutine,
//     same rule and same recover boundary as block A (goja_runtime.go);
//   - not a memory bound. goja has no memory limit and nothing here pretends
//     otherwise; the only ceilings are block A's wall clock and result cap.
//
// Wiring status, said once so no reader has to guess: the loader that would
// call this gate is block B's, and block B waits on the C24 initial set. Until
// then the gate has no production caller BY DESIGN - it is the component that
// hop plugs into, with the refusal behaviour and the per-plugin teardown nailed
// today. The alternative (a gate that guesses at a loader shape) is exactly
// the shape other tickets in this repo were filed to remove.

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Tier2PluginPolicy is the per-plugin face the gate needs from
// [plugins.<id>]: the enable switch and nothing else. Capabilities,
// net_allowlist and host_api are the C24 surface's and the bridge's business
// (block B), so they are absent here rather than present and unchecked.
type Tier2PluginPolicy struct {
	Enabled bool
}

// Tier2Policy is the applied [plugins] section, copied by value: a caller's
// later edit of its config snapshot must not reach into a live gate.
type Tier2Policy struct {
	Tier2Enabled bool
	Plugins      map[string]Tier2PluginPolicy
}

// Tier2GateSnapshot is the gate's observable state: what the applied config
// says, whether this process holds a verdict, the effective switch, and the
// plugin ids that currently own a runtime envelope (sorted). JSON-tagged
// because evidence and a panel reading are the obvious consumers.
type Tier2GateSnapshot struct {
	PolicyEnabled bool     `json:"policy_enabled"`
	Confirmed     bool     `json:"confirmed"`
	Enabled       bool     `json:"enabled"`
	Runtimes      []string `json:"runtimes"`
}

// Tier2GateOption configures a gate.
type Tier2GateOption func(*Tier2Gate)

// WithTier2RuntimeFactory overrides how one plugin's runtime envelope is
// built. Tests use it to count constructions ("no VM exists" is otherwise
// unfalsifiable); a host uses it to stamp non-default block-A options on
// every Tier-2 plugin at once.
func WithTier2RuntimeFactory(f func() *GojaRuntime) Tier2GateOption {
	return func(g *Tier2Gate) {
		if f != nil {
			g.factory = f
		}
	}
}

// WithTier2GateOptions forwards block-A runtime options to every envelope the
// gate builds (default: NewGojaRuntime()).
func WithTier2GateOptions(opts ...GojaOption) Tier2GateOption {
	return func(g *Tier2Gate) {
		g.runtimeOpts = append(g.runtimeOpts, opts...)
	}
}

// Tier2Gate is the enable gate. The zero value is not usable (the maps are
// nil and the default denies everything anyway); use NewTier2Gate. One gate
// serves every plugin and may be read from any goroutine.
type Tier2Gate struct {
	mu          sync.Mutex
	policy      Tier2Policy
	confirmed   bool
	runtimes    map[string]*GojaRuntime
	factory     func() *GojaRuntime
	runtimeOpts []GojaOption
}

// NewTier2Gate builds a gate over an applied policy. It starts CLOSED even
// when policy.Tier2Enabled is true: the value came from a file, and only a
// verdict can open the gate (Enable).
func NewTier2Gate(policy Tier2Policy, opts ...Tier2GateOption) *Tier2Gate {
	g := &Tier2Gate{
		policy:   copyTier2Policy(policy),
		runtimes: map[string]*GojaRuntime{},
	}
	g.factory = func() *GojaRuntime { return NewGojaRuntime(g.runtimeOpts...) }
	for _, opt := range opts {
		if opt != nil {
			opt(g)
		}
	}
	return g
}

// Enable asks the D36 rule 1 L2 re-confirmation and opens the gate only when
// the hook says yes AND the applied config still turns tier 2 on. The hook is
// called with the gate unlocked: it may block on a human, and nothing has
// changed until it answers (the same ordering internal/config uses for
// ConfirmLocked).
func (g *Tier2Gate) Enable(confirm func() bool) error {
	if confirm == nil {
		return gateRefusal("tier 2: enabling needs the D36 rule 1 L2 re-confirmation and no hook was supplied " +
			"(nil = deny; the gate never assumes an approval)")
	}
	if !confirm() {
		return gateRefusal("tier 2: the L2 re-confirmation for [plugins].tier2_enabled was declined; " +
			"the gate stays off and no goja runtime is created")
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.policy.Tier2Enabled {
		return gateRefusal("tier 2: [plugins].tier2_enabled is false in the applied config - a confirmation " +
			"cannot turn on what the file keeps off; the loosening has to be applied by the config layer first")
	}
	g.confirmed = true
	return nil
}

// Runtime returns the plugin's runtime envelope, building it on first use
// (D3 use-to-wake). Every refusal happens BEFORE the factory runs, so a
// closed gate leaves no object behind to tear down.
func (g *Tier2Gate) Runtime(pluginID string) (*GojaRuntime, error) {
	if strings.TrimSpace(pluginID) == "" {
		return nil, observe.New(observe.ClassConfig,
			"tier2 gate: a runtime is keyed on the plugin id and none was supplied")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	switch {
	case !g.policy.Tier2Enabled:
		return nil, gateRefusal("plugin %q: tier 2 is off - [plugins].tier2_enabled is false in the applied "+
			"config, so no goja runtime exists for it (default-off gate, ticket 51)", pluginID)
	case !g.confirmed:
		return nil, gateRefusal("plugin %q: tier 2 holds no L2 re-confirmation - [plugins].tier2_enabled = true "+
			"in the file is data, not a verdict; D36 rule 1 applies the loosening only after the card answers "+
			"yes in this process", pluginID)
	}

	entry, ok := g.policy.Plugins[pluginID]
	if !ok {
		return nil, gateRefusal("plugin %q: no [plugins.%s].enabled key - tier 2 being on does not turn a "+
			"plugin on; every Tier-2 plugin carries its own switch (SPEC-07 §4)", pluginID, pluginID)
	}
	if !entry.Enabled {
		return nil, gateRefusal("plugin %q: [plugins.%s].enabled = false - the per-plugin switch is off, so it "+
			"gets no runtime even while tier 2 is on", pluginID, pluginID)
	}

	if rt, live := g.runtimes[pluginID]; live {
		return rt, nil
	}
	rt := g.factory()
	g.runtimes[pluginID] = rt
	return rt, nil
}

// Execute is the call door: gate first, then the block-A envelope. A refused
// call never reaches goja, so "disabled" means "no code of this plugin ran",
// not "its code ran and was told no".
func (g *Tier2Gate) Execute(ctx context.Context, call GojaCall) (GojaResult, error) {
	rt, err := g.Runtime(call.PluginID)
	if err != nil {
		return GojaResult{}, err
	}
	return rt.Execute(ctx, call)
}

// ApplyPolicy adopts a new applied [plugins] section, dropping the runtime of
// every plugin that is now disabled, removed, or behind a closed tier-2
// switch. Neighbours keep their envelopes. Going to tier2 off also clears the
// L2 record, so re-enabling after a tightening needs a fresh card.
//
// The only teardown the gate owns is that envelope reference: block A creates
// a fresh VM per call and disposes its own per-call DisposalScope, so there is
// no idle VM, goroutine or scope to join here - which is why ApplyPolicy is
// synchronous and cannot block.
func (g *Tier2Gate) ApplyPolicy(policy Tier2Policy) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.policy = copyTier2Policy(policy)
	if !g.policy.Tier2Enabled {
		g.confirmed = false
	}
	for id := range g.runtimes {
		entry, ok := g.policy.Plugins[id]
		if !g.policy.Tier2Enabled || !ok || !entry.Enabled {
			delete(g.runtimes, id)
		}
	}
}

// Snapshot returns the gate's state. Runtimes is always non-nil so a JSON
// reading says [] rather than null.
func (g *Tier2Gate) Snapshot() Tier2GateSnapshot {
	g.mu.Lock()
	defer g.mu.Unlock()
	snap := Tier2GateSnapshot{
		PolicyEnabled: g.policy.Tier2Enabled,
		Confirmed:     g.confirmed,
		Enabled:       g.policy.Tier2Enabled && g.confirmed,
		Runtimes:      make([]string, 0, len(g.runtimes)),
	}
	for id := range g.runtimes {
		snap.Runtimes = append(snap.Runtimes, id)
	}
	sort.Strings(snap.Runtimes)
	return snap
}

// gateRefusal builds the D37 permission_denied error a gate refusal carries:
// the class says "policy, not bug", and D37's retry column for it is
// RetryNever, so no caller may retry its way through a closed gate.
func gateRefusal(format string, args ...any) error {
	return observe.New(observe.ClassPermissionDenied, fmt.Sprintf(format, args...))
}

// copyTier2Policy detaches a policy from the caller's map.
func copyTier2Policy(policy Tier2Policy) Tier2Policy {
	out := Tier2Policy{Tier2Enabled: policy.Tier2Enabled}
	if len(policy.Plugins) > 0 {
		out.Plugins = make(map[string]Tier2PluginPolicy, len(policy.Plugins))
		for id, entry := range policy.Plugins {
			out.Plugins[id] = entry
		}
	}
	return out
}
