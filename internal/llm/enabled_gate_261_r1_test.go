package llm_test

// Ticket 261 leg r1 - rulers for the fix, written by the leg that landed it.
//
// The fix (internal/llm/resolver.go, this leg): the enumeration point
// (DiscoveredModels) skips Enabled==false entries and the selection point
// (resolveEndpoint) refuses them by name. These rulers measure the gate's
// own three hops - enumerate / select / role-fallback - plus the exact error
// shape. The end-to-end reach/cost rulers live in enabled_reach_261_test.go
// (p1's fixture, rewritten by this leg to the post-fix semantics); the
// load/write-side rulers live in internal/config/enabled_261_test.go
// (untouched: leg r1 does not modify internal/config).
//
// Mutation self-proof lives in .scratch/wisp/probes/261/r1/mut/ (an overlay
// copy of resolver.go with both filters removed) - the tracked tree is never
// mutated; see probes/261/r1/impl.md section 3.

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/observe"
)

// gate261Resolver builds the minimal shape: one provider, two entries - the
// disabled one is only reachable by explicit name (never in text_chain), the
// enabled one is the positive control. Written by hand (map literal, not a
// TOML file) so the ruler isolates the resolver from the loader.
func gate261Resolver() (*llm.Resolver, *config.Config) {
	cfg := &config.Config{}
	cfg.LLM.Providers = map[string]config.Provider{
		"p261": {
			Protocol: "openai-chat",
			BaseURL:  "http://127.0.0.1:9/v1",
			Models: map[string]config.ModelSpec{
				"m-on":  {Enabled: true},
				"m-off": {Enabled: false},
			},
		},
	}
	cfg.LLM.TextChain = []string{"p261/m-on"}
	return llm.NewResolver(cfg, nil), cfg
}

// TestTicket261R1EnumerationSkipsDisabled: the enumeration point must not
// list a disabled entry (the ModelSpec.Enabled promise, discovery half).
func TestTicket261R1EnumerationSkipsDisabled(t *testing.T) {
	res, _ := gate261Resolver()
	listed := res.DiscoveredModels("p261")
	joined := strings.Join(listed, ",")
	if !strings.Contains(joined, "m-on") {
		t.Errorf("DiscoveredModels = %v, want the enabled entry still listed", listed)
	}
	if strings.Contains(joined, "m-off") {
		t.Errorf("DiscoveredModels = %v, want the disabled entry removed from discovery", listed)
	}
	if len(listed) != 1 {
		t.Errorf("DiscoveredModels = %v, want exactly the one enabled entry", listed)
	}
	// A provider absent from the catalog still returns nil (no behavior change
	// beyond the filter).
	if got := res.DiscoveredModels("ghost261"); got != nil {
		t.Errorf("DiscoveredModels(unknown provider) = %v, want nil", got)
	}
}

// TestTicket261R1SelectionRefusesDisabledByExactWording: naming a disabled
// entry through text_chain must be refused with the disabled wording - NOT a
// reuse of, nor a substring of, the unknown-model error - and the wire must
// never be built for it (BuildEndpointProvider is downstream of resolveEndpoint,
// so resolveEndpoint's refusal is what keeps dispatch away).
func TestTicket261R1SelectionRefusesDisabledByExactWording(t *testing.T) {
	res, _ := gate261Resolver()
	res.TextChain = []string{"p261/m-off"}
	_, err := res.ResolveChain()
	if err == nil {
		t.Fatal("ResolveChain(disabled) succeeded; the gate is missing")
	}
	if class := adaptertest.ClassOf(err); class != observe.ClassConfig {
		t.Errorf("class = %v, want config", class)
	}
	msg := err.Error()
	if !strings.Contains(msg, "disabled") {
		t.Errorf("error = %q, want the disabled wording", msg)
	}
	if !strings.Contains(msg, "p261") || !strings.Contains(msg, "m-off") {
		t.Errorf("error = %q, want it to name both the provider and the model", msg)
	}
	if strings.Contains(msg, "unknown model") {
		t.Errorf("error = %q, must not reuse the unknown-model wording", msg)
	}
	// Ghost control at the same line: absent entry keeps the old wording, so
	// the two refusals are distinguishable in logs.
	res.TextChain = []string{"p261/ghost"}
	_, ghostErr := res.ResolveChain()
	if ghostErr == nil {
		t.Fatal("ResolveChain(ghost) succeeded")
	}
	if !strings.Contains(ghostErr.Error(), "unknown model") {
		t.Errorf("ghost error = %q, want the unknown-model wording", ghostErr)
	}
	if strings.Contains(ghostErr.Error(), "disabled") {
		t.Errorf("ghost error = %q, must not mention disabled", ghostErr)
	}
}

// TestTicket261R1RoleFallbackAlsoGated: the unset-role fallback walks the
// same resolveEndpoint line (run.go:436 shape). A disabled FIRST element must
// be refused, not silently skipped over - skipping would re-order the chain
// behind the user's back.
func TestTicket261R1RoleFallbackAlsoGated(t *testing.T) {
	res, _ := gate261Resolver()
	res.TextChain = []string{"p261/m-off", "p261/m-on"}
	if _, _, err := res.ResolveRole(llm.RoleChat); err == nil || !strings.Contains(err.Error(), "disabled") {
		t.Errorf("ResolveRole fallback onto a disabled first element = %v, want the disabled refusal (not a silent skip)", err)
	}
	// The enabled-first control: fallback resolves the chain head normally.
	res2, _ := gate261Resolver()
	res2.TextChain = []string{"p261/m-on"}
	if _, _, err := res2.ResolveRole(llm.RoleChat); err != nil {
		t.Errorf("ResolveRole fallback onto the enabled head = %v, want nil", err)
	}
}
