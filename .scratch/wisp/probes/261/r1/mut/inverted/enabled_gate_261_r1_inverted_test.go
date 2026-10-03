package llm_test

// SCRATCH BENCH ONLY - inverted-form rulers for the constant-truth check
// (impl.md section 3). This file never compiles inside the tracked tree.
// Each assertion below is the OPPOSITE of enabled_gate_261_r1_test.go: if
// these pass against the GATED tracked resolver, the real rulers are
// insensitive and must be rewritten.

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
)

func gate261ResolverInv() *llm.Resolver {
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
	return llm.NewResolver(cfg, nil)
}

func TestTicket261R1EnumerationSkipsDisabled(t *testing.T) {
	res := gate261ResolverInv()
	listed := strings.Join(res.DiscoveredModels("p261"), ",")
	if strings.Contains(listed, "m-on") {
		t.Errorf("INVERTED: enabled entry %q listed, want NOT listed", listed)
	}
	if !strings.Contains(listed, "m-off") {
		t.Errorf("INVERTED: disabled entry missing from %q, want it listed", listed)
	}
}

func TestTicket261R1SelectionRefusesDisabledByExactWording(t *testing.T) {
	res := gate261ResolverInv()
	res.TextChain = []string{"p261/m-off"}
	if _, err := res.ResolveChain(); err != nil {
		t.Errorf("INVERTED: ResolveChain(disabled) = %v, want nil (admitted)", err)
	}
	res.TextChain = []string{"p261/ghost"}
	if _, err := res.ResolveChain(); err != nil {
		t.Errorf("INVERTED: ResolveChain(ghost) = %v, want nil (admitted)", err)
	}
}

func TestTicket261R1RoleFallbackAlsoGated(t *testing.T) {
	res := gate261ResolverInv()
	res.TextChain = []string{"p261/m-off", "p261/m-on"}
	if _, _, err := res.ResolveRole(llm.RoleChat); err != nil {
		t.Errorf("INVERTED: ResolveRole(disabled head) = %v, want nil", err)
	}
}
