package agent

// Ticket 370: the chat role's own output ceiling has to reach the wire.
//
// Before this ticket the value existed twice and travelled zero times:
// llm.Resolver copied [llm.roles.chat].max_output_tokens into RoleSettings, the
// composition root discarded it (`ep, _, err := res.ResolveRole(...)`), and
// agent.Config had no field to carry it, so Assembler.MaxOutputTokens stayed 0
// and llm.Request.MaxOutputTokens stayed 0 with it. Every request went out with
// no max_tokens at all - which is exactly what cmd/wisp/tasklogstate_365_test.go
// recorded as "wisp run never sends max_tokens".
//
// The assertions here are ordered from the seam outwards, and the last one is
// the only one that cannot pass by bookkeeping: it reads the JSON body the
// openai-chat adapter actually put on the socket. A Config field that is set
// and then dropped, or an Assembler field nobody copies, fails there.
//
// The second half pins the other direction - the value must NOT appear when the
// role set none and the endpoint is the same one - so "always sends something"
// cannot masquerade as a pass.

import (
	"bytes"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

const ticket370MaxOutput = 370

// Test370RoleMaxOutputReachesTheAssemblerAndTheWire is the delivered claim:
// Config.MaxOutputTokens -> Assembler.MaxOutputTokens -> Request.MaxOutputTokens
// -> the "max_tokens" field of the JSON body the provider receives.
func Test370RoleMaxOutputReachesTheAssemblerAndTheWire(t *testing.T) {
	h := newHarness(t, "text-reply", withConfig(func(cfg *Config) {
		cfg.MaxOutputTokens = ticket370MaxOutput
	}))

	// Link 1: New copies the config value onto the assembler.
	if got := h.loop.asm.MaxOutputTokens; got != ticket370MaxOutput {
		t.Fatalf("Assembler.MaxOutputTokens = %d, want %d: agent.New dropped Config.MaxOutputTokens",
			got, ticket370MaxOutput)
	}

	res := h.run("这个上限到底有没有发出去")
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", string(res.Status), res.Message)
	}

	// Link 2: Build puts it on the C5 request.
	req := h.loop.asm.Build(PromptInput{}, nil)
	if req.MaxOutputTokens != ticket370MaxOutput {
		t.Fatalf("Request.MaxOutputTokens = %d, want %d: Assembler.Build did not carry the ceiling",
			req.MaxOutputTokens, ticket370MaxOutput)
	}

	// Link 3 (the one that cannot be faked by struct bookkeeping): the
	// openai-chat adapter renders it as max_tokens on the wire.
	bodies := h.requestBodies()
	if len(bodies) == 0 {
		t.Fatal("no request went out at all")
	}
	if !bytes.Contains(bodies[0], []byte(`"max_tokens":370,`)) {
		t.Fatalf("the first request body carries no \"max_tokens\":370 (ticket 370 field):\n%s",
			trunc(string(bodies[0]), 800))
	}
}

// Test370UnsetRoleMaxOutputSendsNoCeiling is the control. Same harness, same
// fixture, same model - only the config field is left at 0, which is what
// llm.Resolver produces for a role that never set max_output_tokens. The body
// must then carry no max_tokens at all (openaichat's `json:"max_tokens,omitempty"`),
// i.e. the pre-370 behaviour, so the assertion above is measuring the field
// rather than some adapter default that is always present.
func Test370UnsetRoleMaxOutputSendsNoCeiling(t *testing.T) {
	h := newHarness(t, "text-reply")
	if h.loop.asm.MaxOutputTokens != 0 {
		t.Fatalf("Assembler.MaxOutputTokens = %d for an unset config, want 0", h.loop.asm.MaxOutputTokens)
	}
	if res := h.run("没有配上限的时候"); res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", string(res.Status), res.Message)
	}
	bodies := h.requestBodies()
	if len(bodies) == 0 {
		t.Fatal("no request went out at all")
	}
	if bytes.Contains(bodies[0], []byte(`"max_tokens"`)) {
		t.Fatalf("an unset role ceiling still produced a max_tokens field:\n%s",
			trunc(string(bodies[0]), 800))
	}
}

// Test370MaxOutputIsPerLoopNotGlobal guards the sequence of the fix: the
// assembler gets the value inside New, before the Loop is handed out, so two
// loops built from two different configs keep their own ceilings. A fix that
// scribbled on a package-level default (or mutated a shared Assembler) would
// make the second loop report the first one's number.
func Test370MaxOutputIsPerLoopNotGlobal(t *testing.T) {
	first := newHarness(t, "text-reply", withConfig(func(cfg *Config) { cfg.MaxOutputTokens = 111 }))
	second := newHarness(t, "text-reply", withConfig(func(cfg *Config) { cfg.MaxOutputTokens = 222 }))
	if got := first.loop.asm.MaxOutputTokens; got != 111 {
		t.Errorf("first loop ceiling = %d, want 111", got)
	}
	if got := second.loop.asm.MaxOutputTokens; got != 222 {
		t.Errorf("second loop ceiling = %d, want 222", got)
	}
	if got := first.loop.asm.Build(PromptInput{}, nil).MaxOutputTokens; got != 111 {
		t.Errorf("first loop request ceiling = %d after the second loop was built, want 111", got)
	}
	// And the request model is untouched by this ticket: the ceiling rides with
	// the role, not with a new endpoint identity.
	if got := first.loop.asm.Build(PromptInput{}, nil).Model; got != "mock-small" {
		t.Errorf("request model = %q, want mock-small: ticket 370 must not move the model", got)
	}
	_ = llm.RoleChat // the role this value is resolved for; kept explicit for readers
}
