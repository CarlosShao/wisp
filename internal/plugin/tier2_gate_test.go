package plugin

// Ticket 51, the "Enable gate" acceptance criterion. The unit under test is
// internal/plugin/tier2_gate.go; these tests are its only door.
//
// Three claims, one reading each - and each is a reading that a refusal alone
// cannot produce:
//  1. default off means NO runtime object exists (the factory counter is 0),
//     not "one exists and refuses";
//  2. enabling is not a boolean flip: a config.toml that says
//     tier2_enabled = true at boot reaches no runtime until the D36 rule 1
//     verdict has been asked and answered yes IN THIS PROCESS;
//  3. [plugins.<id>].enabled = false removes that plugin's runtime and refuses
//     it, while a sibling plugin keeps running.
//
// Test 4 is the end-to-end half: the real config.Manager (no doubles) is the
// one that decides whether the loosening applied at all.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// gateFixture is a gate plus the count of runtime objects its factory built.
// The counter is the difference between "no VM exists" and "a VM exists that
// politely says no": only the first one reads 0.
type gateFixture struct {
	gate  *Tier2Gate
	built *int
}

func newTestGate(t *testing.T, policy Tier2Policy) gateFixture {
	t.Helper()
	reg := observe.NewRegistry()
	built := 0
	g := NewTier2Gate(policy, WithTier2RuntimeFactory(func() *GojaRuntime {
		built++
		return NewGojaRuntime(WithGojaRegistry(reg))
	}))
	return gateFixture{gate: g, built: &built}
}

// policyFrom is the one-line adapter a loader writes: the applied [plugins]
// section -> the gate's face. It lives in the test because making it
// production code would couple internal/plugin to internal/config; the loader
// that owns both (block B) is where that hop belongs.
func policyFrom(c *config.Config) Tier2Policy {
	p := Tier2Policy{
		Tier2Enabled: c.Plugins.Tier2Enabled,
		Plugins:      make(map[string]Tier2PluginPolicy, len(c.Plugins.Entries)),
	}
	for id, e := range c.Plugins.Entries {
		p.Plugins[id] = Tier2PluginPolicy{Enabled: e.Enabled}
	}
	return p
}

// TestTier2GateDefaultOffBuildsNoRuntime is the "default off (no VM exists)"
// half of the AC. The shipped zero policy is the default: tier2 off, no plugin
// enabled, so a call is refused AND nothing is constructed to refuse it.
func TestTier2GateDefaultOffBuildsNoRuntime(t *testing.T) {
	t.Parallel()
	fx := newTestGate(t, Tier2Policy{})
	g := fx.gate

	snap := g.Snapshot()
	if snap.PolicyEnabled || snap.Confirmed || snap.Enabled {
		t.Fatalf("the shipped default must read off, got %+v", snap)
	}
	if len(snap.Runtimes) != 0 {
		t.Fatalf("default gate already owns runtimes: %v", snap.Runtimes)
	}

	_, err := g.Runtime("demo")
	if got := observe.ClassOfOrInternal(err); got != observe.ClassPermissionDenied {
		t.Fatalf("Runtime while tier2 is off: class = %v (err %v), want permission_denied", got, err)
	} else if !strings.Contains(err.Error(), "tier2_enabled") {
		t.Fatalf("the refusal must name the key that is off, got %q", err)
	}
	if _, err := g.Execute(context.Background(), GojaCall{PluginID: "demo", Source: "1+1"}); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("Execute while tier2 is off: class = %v (err %v), want permission_denied",
			observe.ClassOfOrInternal(err), err)
	}

	if *fx.built != 0 {
		t.Fatalf("tier2 off still built %d runtime object(s): no VM may exist while the gate is off", *fx.built)
	}
	if n := len(g.Snapshot().Runtimes); n != 0 {
		t.Fatalf("refused calls left %d runtime(s) behind", n)
	}
}

// TestTier2GateEnableNeedsL2Confirm is the "enabling without confirm is
// impossible" half. Three refusals (boot-time true, nil hook, declined hook)
// then the one path that works.
func TestTier2GateEnableNeedsL2Confirm(t *testing.T) {
	t.Parallel()
	policy := Tier2Policy{
		Tier2Enabled: true,
		Plugins:      map[string]Tier2PluginPolicy{"demo": {Enabled: true}},
	}
	fx := newTestGate(t, policy)
	g := fx.gate

	// (a) The file says true; that is data, not a verdict.
	if _, err := g.Runtime("demo"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("a config.toml that says tier2_enabled = true reached a runtime without the L2 verdict: %v", err)
	} else if !strings.Contains(err.Error(), "L2") {
		t.Fatalf("the refusal must name the missing L2 re-confirmation, got %q", err)
	}

	// (b) No hook at all = deny (fail closed), never "assume approved".
	if err := g.Enable(nil); err == nil {
		t.Fatal("Enable(nil) was accepted: a missing confirmation hook must deny")
	} else if observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("Enable(nil): class = %v, want permission_denied", observe.ClassOfOrInternal(err))
	}
	if g.Snapshot().Enabled {
		t.Fatal("Enable(nil) turned the gate on")
	}

	// (c) The human said no.
	asked := 0
	if err := g.Enable(func() bool { asked++; return false }); err == nil {
		t.Fatal("a declined L2 card was accepted")
	}
	if asked != 1 {
		t.Fatalf("the confirmation hook was called %d time(s), want 1", asked)
	}
	if g.Snapshot().Enabled || g.Snapshot().Confirmed {
		t.Fatal("a declined L2 card still enabled the gate")
	}
	if _, err := g.Runtime("demo"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("Runtime after a declined card: err = %v, want permission_denied", err)
	}
	if *fx.built != 0 {
		t.Fatalf("refusals built %d runtime(s)", *fx.built)
	}

	// (d) The one path that works.
	if err := g.Enable(func() bool { asked++; return true }); err != nil {
		t.Fatalf("approved L2 card was rejected: %v", err)
	}
	snap := g.Snapshot()
	if !snap.Enabled || !snap.Confirmed || !snap.PolicyEnabled {
		t.Fatalf("after an approved card the gate must read on: %+v", snap)
	}
	rt, err := g.Runtime("demo")
	if err != nil {
		t.Fatalf("Runtime after approval: %v", err)
	}
	if rt == nil {
		t.Fatal("Runtime after approval returned a nil runtime and a nil error")
	}
	if *fx.built != 1 {
		t.Fatalf("runtime factory ran %d time(s) for one plugin, want 1 (lazy, then cached)", *fx.built)
	}
	rt2, err := g.Runtime("demo")
	if err != nil || rt2 != rt {
		t.Fatalf("second Runtime call did not reuse the envelope: %v / %v", rt2 == rt, err)
	}
	if *fx.built != 1 {
		t.Fatalf("the runtime is re-created per call: factory ran %d time(s)", *fx.built)
	}

	// The runtime the gate hands out is a working block-A envelope.
	out, err := g.Execute(context.Background(), GojaCall{PluginID: "demo", Source: `(function(){return "hi"})()`})
	if err != nil {
		t.Fatalf("gated Execute: %v", err)
	}
	if strings.Trim(out.Value, `"`) != "hi" {
		t.Fatalf("gated Execute value = %q, want \"hi\"", out.Value)
	}
	if out.Truncated {
		t.Fatal("a two-byte result was reported as truncated")
	}
}

// TestTier2GatePerPluginDisable is "per-plugin disable works": one plugin's
// switch falling drops only that plugin, and a plugin that never had a key is
// refused with the key path named.
func TestTier2GatePerPluginDisable(t *testing.T) {
	t.Parallel()
	policy := Tier2Policy{
		Tier2Enabled: true,
		Plugins: map[string]Tier2PluginPolicy{
			"demo":  {Enabled: true},
			"other": {Enabled: true},
		},
	}
	fx := newTestGate(t, policy)
	g := fx.gate
	if err := g.Enable(func() bool { return true }); err != nil {
		t.Fatalf("Enable: %v", err)
	}
	demo, err := g.Runtime("demo")
	if err != nil {
		t.Fatalf("demo runtime: %v", err)
	}
	other, err := g.Runtime("other")
	if err != nil {
		t.Fatalf("other runtime: %v", err)
	}
	if demo == other {
		t.Fatal("two plugins share one runtime envelope")
	}

	// demo's enable key goes false; other's stays true.
	g.ApplyPolicy(Tier2Policy{
		Tier2Enabled: true,
		Plugins: map[string]Tier2PluginPolicy{
			"demo":  {Enabled: false},
			"other": {Enabled: true},
		},
	})
	if got := g.Snapshot().Runtimes; len(got) != 1 || got[0] != "other" {
		t.Fatalf("after disabling demo the live runtimes = %v, want [other]", got)
	}
	if _, err := g.Runtime("demo"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("disabled plugin still got a runtime: %v", err)
	} else if !strings.Contains(err.Error(), "[plugins.demo].enabled") {
		t.Fatalf("the refusal must name the per-plugin key, got %q", err)
	}
	kept, err := g.Runtime("other")
	if err != nil {
		t.Fatalf("the sibling plugin was collateral damage: %v", err)
	}
	if kept != other {
		t.Fatal("the sibling plugin's envelope was rebuilt: a neighbour's disable must not touch it")
	}
	if out, err := g.Execute(context.Background(), GojaCall{PluginID: "other", Source: "6*7"}); err != nil {
		t.Fatalf("sibling Execute after the neighbour was disabled: %v", err)
	} else if out.Value != "42" {
		t.Fatalf("sibling Execute value = %q, want \"42\"", out.Value)
	}

	// A plugin with no key at all is not silently "on by default".
	g.ApplyPolicy(Tier2Policy{
		Tier2Enabled: true,
		Plugins:      map[string]Tier2PluginPolicy{"other": {Enabled: true}},
	})
	if _, err := g.Runtime("demo"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("a plugin with no enable key got a runtime: %v", err)
	} else if !strings.Contains(err.Error(), "no [plugins.demo].enabled key") {
		t.Fatalf("the refusal must say the key is missing, got %q", err)
	}

	// The gate-wide switch going off drops the last runtime and clears the
	// verdict, so re-enabling needs a fresh card - a stale yes is not a yes.
	g.ApplyPolicy(Tier2Policy{})
	snap := g.Snapshot()
	if snap.PolicyEnabled || snap.Confirmed || snap.Enabled || len(snap.Runtimes) != 0 {
		t.Fatalf("after tier2 went off the gate must be fully closed: %+v", snap)
	}
	if err := g.Enable(func() bool { return true }); err == nil {
		t.Fatal("a confirmation turned tier2 on while the applied config says false")
	}

	// Re-opening tier 2 must not resurrect the old yes: the tightening cleared
	// the record, so the gate is closed again until a fresh card is answered.
	g.ApplyPolicy(Tier2Policy{
		Tier2Enabled: true,
		Plugins:      map[string]Tier2PluginPolicy{"other": {Enabled: true}},
	})
	if g.Snapshot().Confirmed || g.Snapshot().Enabled {
		t.Fatalf("a stale L2 verdict survived the tightening: %+v", g.Snapshot())
	}
	if _, err := g.Runtime("other"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("a re-opened tier 2 handed out a runtime on the old verdict: %v", err)
	}
}

// TestTier2GateBootTrueIsNotAnApproval is the e2e seam against the real
// config.Manager: a hand-edited config.toml that says tier2_enabled = true
// loads as true (the initial load raises no card), and the gate still hands
// out nothing until the L2 hook has answered in this process.
func TestTier2GateBootTrueIsNotAnApproval(t *testing.T) {
	t.Parallel()
	path := writeTier2Config(t, `schema_version = 1

[plugins]
tier2_enabled = true

[plugins.demo]
enabled = true
`)
	m, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("fixture must load (tier2_enabled is not an unwired key): %v", err)
	}
	policy := policyFrom(m.Config())
	if !policy.Tier2Enabled {
		t.Fatal("fixture: the initial load was supposed to carry tier2_enabled = true")
	}
	if !policy.Plugins["demo"].Enabled {
		t.Fatal("fixture: [plugins.demo] enabled = true did not reach the policy face")
	}

	fx := newTestGate(t, policy)
	if _, err := fx.gate.Runtime("demo"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("boot-time tier2_enabled = true reached a runtime with no L2 card: %v", err)
	}
	if *fx.built != 0 {
		t.Fatalf("boot-time true built %d runtime(s) with no card on record", *fx.built)
	}
	if err := fx.gate.Enable(func() bool { return true }); err != nil {
		t.Fatalf("Enable after the L2 card: %v", err)
	}
	if _, err := fx.gate.Runtime("demo"); err != nil {
		t.Fatalf("Runtime after the L2 card: %v", err)
	}
}

// TestTier2GateReloadL2ControlsTheRuntime is the reload half of the same seam:
// the real config.Manager declines the loosening (no card answered), then
// approves it, and the gate follows the APPLIED config - never the bytes on
// disk.
func TestTier2GateReloadL2ControlsTheRuntime(t *testing.T) {
	t.Parallel()
	path := writeTier2Config(t, `schema_version = 1

[plugins]
tier2_enabled = false

[plugins.demo]
enabled = true
`)
	m, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("off fixture: %v", err)
	}
	fx := newTestGate(t, policyFrom(m.Config()))
	g := fx.gate
	if g.Snapshot().PolicyEnabled {
		t.Fatal("fixture: the off config must start the gate off")
	}

	// Flip the file on. The L2 hook declines.
	var asked []string
	m.ConfirmLocked = func(section string, keys []string) bool {
		asked = append(asked, section+":"+strings.Join(keys, ","))
		return false
	}
	overwriteTier2Config(t, path, `schema_version = 1

[plugins]
tier2_enabled = true
# first reload, declined by the human

[plugins.demo]
enabled = true
`)
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload (declined): %v", err)
	}
	if len(rep.Locked) != 1 || rep.Locked[0].Section != "plugins" {
		t.Fatalf("the [plugins] loosening must raise the L2 card, got %+v", rep.Locked)
	}
	if rep.Locked[0].Approved {
		t.Fatal("a declined card was recorded as approved")
	}
	if !containsKey(rep.Locked[0].Keys, "plugins.tier2_enabled") {
		t.Fatalf("the card did not name plugins.tier2_enabled: %v", rep.Locked[0].Keys)
	}
	if len(asked) != 1 {
		t.Fatalf("the L2 hook was asked %d time(s), want 1 (%v)", len(asked), asked)
	}
	if m.Config().Plugins.Tier2Enabled {
		t.Fatal("a declined L2 card still applied tier2_enabled = true")
	}
	g.ApplyPolicy(policyFrom(m.Config()))
	if g.Snapshot().PolicyEnabled || g.Snapshot().Enabled {
		t.Fatalf("the gate followed the bytes instead of the applied config: %+v", g.Snapshot())
	}
	if _, err := g.Runtime("demo"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("Runtime after a declined loosening: %v", err)
	}

	// Same write, this time approved.
	m.ConfirmLocked = func(section string, keys []string) bool { return true }
	overwriteTier2Config(t, path, `schema_version = 1

[plugins]
tier2_enabled = true
# second reload, approved by the human (SP2)

[plugins.demo]
enabled = true
`)
	rep2, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload (approved): %v", err)
	}
	if len(rep2.Locked) != 1 || !rep2.Locked[0].Approved {
		t.Fatalf("the approved loosening was not recorded: %+v", rep2.Locked)
	}
	if !m.Config().Plugins.Tier2Enabled {
		t.Fatal("an approved L2 card did not apply tier2_enabled = true")
	}
	g.ApplyPolicy(policyFrom(m.Config()))
	if !g.Snapshot().PolicyEnabled {
		t.Fatal("gate did not pick up the applied loosening")
	}
	// The config layer said yes; the runtime still needs the process-local
	// verdict, which is the gate's own record.
	if _, err := g.Runtime("demo"); observe.ClassOfOrInternal(err) != observe.ClassPermissionDenied {
		t.Fatalf("applied config alone reached a runtime with no process verdict: %v", err)
	}
	if err := g.Enable(func() bool { return true }); err != nil {
		t.Fatalf("Enable after the approved reload: %v", err)
	}
	if _, err := g.Runtime("demo"); err != nil {
		t.Fatalf("Runtime after the approved reload: %v", err)
	}
	if *fx.built != 1 {
		t.Fatalf("runtime factory ran %d time(s), want 1", *fx.built)
	}
}

func writeTier2Config(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	overwriteTier2Config(t, path, body)
	return path
}

func overwriteTier2Config(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func containsKey(keys []string, want string) bool {
	for _, k := range keys {
		if k == want {
			return true
		}
	}
	return false
}
