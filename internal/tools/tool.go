package tools

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// C1 Tool contract (SPEC-07 §2, Pi AgentTool semantics per D21#1). The four
// methods below are the frozen contract: Name / Description / Parameters
// (JSON Schema) / Execute(ctx, params, onUpdate).
//
// Everything the HOST needs but the contract does not give a tool (its
// capability declarations, its R1 declared risk floor, its C22 per-tool
// timeout) lives in the sidecar Decl, never as extra methods on Tool - adding
// methods to Tool would force every ticket-50/51 plugin shim to grow
// host-only plumbing.
type Tool interface {
	// Name is globally unique and 'namespace.action' shaped (C1).
	Name() string
	// Description is what the model reads to decide whether to call this.
	Description() string
	// Parameters is the JSON Schema for the argument object.
	Parameters() JSONSchema
	// Execute runs the call. onUpdate streams progress deltas to the host
	// (may be nil). A returned error is a HOST-internal failure (D37 class
	// internal): it means the tool could not answer at all. A failure the
	// model can act on goes back in Result.Text with IsError set instead.
	Execute(ctx context.Context, params json.RawMessage, onUpdate func(delta string)) (Result, error)
}

// JSONSchema is a JSON Schema document for one tool's parameter object.
type JSONSchema = json.RawMessage

// Result is one tool's answer, before the bridge adds risk/gate/provenance
// framing to it.
type Result struct {
	// Text is the payload the host puts back into the context (the loop caps
	// and spills it, D15(3) - a tool must not pre-truncate to a guessed
	// budget).
	Text string
	// IsError marks a tool-level failure: the bridge books it as
	// outcome=error with class "tool" so the model can self-correct (D37)
	// instead of the task failing.
	IsError bool
	// Truncated is set by a tool that shortened its own output.
	Truncated bool
	// Origin names where the result came from for C25 provenance marking
	// (the canonical path / URL / window title). Empty falls back to the
	// bridge's own view of the call's path parameter.
	Origin string
	// AppliedSteps lists what the tool already did before it stopped.
	// D31: a cancel after partial work MUST report the applied steps and
	// never claim a clean undo. No builtin tool in this segment fills it
	// (fs.read/fs.list have no side effects); the fs.write/trash/move family
	// (ticket 20 segment 2) is the first consumer.
	AppliedSteps []string
}

// Decl is the host-side declaration for one Tool: the metadata the bridge
// consults on every call. It is frozen at registration by the Builtin
// provider and read from the manifest by the Tier-1/Tier-2 providers
// (tickets 50/51).
type Decl struct {
	// Capabilities are the C3 tokens this tool DECLARED at registration.
	// The bridge compares them against Needs on every call.
	Capabilities []Capability
	// Needs are the C3 tokens this tool's Execute actually uses. A token in
	// Needs but not in Capabilities is a hard reject (SPEC-07 §2).
	Needs []Capability
	// Declared is the R1 lower bound only; C19 owns the verdict (R1's whole
	// point is that a declaration is never a conclusion).
	Declared risk.Level
	// PathParams name the parameter keys whose string values are affected
	// filesystem paths. The bridge, not the tool, feeds them to C19 as
	// Facts.Paths so R2/R3 can judge them.
	PathParams []string
	// Timeout is the C22 per-tool budget (manifest timeoutMs); zero falls
	// back to Options.DefaultTimeout. It is honored by context deadline,
	// never by a wall-clock comparison.
	Timeout time.Duration
	// Resident marks a D15(1) builtin: always injected into the request.
	// Non-resident entries are selected by local retrieval (D15(2)).
	Resident bool
	// Provider is which C4 slot supplied this entry.
	Provider Kind
}
