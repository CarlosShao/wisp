package manifestslot

// The Manifest slot's Tier-2 arm (ticket 51 block B).
//
// A Tier-2 plugin's manifest has the same C14 face as a Tier-1 one; what
// differs is what its tools DO: instead of a declared host action they run a JS
// entry inside goja through the frozen C24 surface (internal/plugin/c24.go).
// The slot itself does not care - it publishes tools.Entry values and calls
// whatever Run it was built with - so the whole Tier-2 wiring is this adapter
// from "one manifest tool call" to "one C24Executor call".
//
// Two boundaries are deliberate:
//
//   - The adapter converts the executor's ERROR into an is_error Result instead
//     of returning it. tools.Bridge treats a non-nil error from Tool.Execute as
//     a provider fault ("工具内部故障", class tool) and would lose the D37 class
//     the executor already set (permission_denied for a closed gate, config for
//     a missing bundle). The model needs the refusal's own text, and the row
//     needs the right class, so the refusal travels as data.
//
//   - It does not touch capabilities or risk. The manifest's declaration is
//     already on Decl.Needs (entryFor), and tools.Bridge.checkCaps re-reads it
//     per call; Tier2Run passing def.Capabilities to the surface is how the same
//     declaration decides which C24 entries exist in the VM.

import (
	"context"
	"encoding/json"

	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Tier2Run is the Run a Tier-2 provider is built with. A nil executor is
// refused at call time rather than at construction because New's signature
// (ticket 50) takes a Run value: the slot stays compatible, and a host that
// forgot the executor sees a readable tool error instead of a nil panic.
func Tier2Run(ex *plugin.Tier2Executor) Run {
	return func(ctx context.Context, pluginID string, def plugin.ToolDef, params json.RawMessage) (tools.Result, error) {
		if ex == nil {
			return tools.Result{
				Text:    "tier-2 executor is not wired: the plugin's JS runtime was never composed",
				IsError: true,
			}, nil
		}
		out, err := ex.Execute(ctx, plugin.Tier2ToolRequest{
			PluginID: pluginID,
			Tool:     def.Name,
			Declared: def.Capabilities,
			Params:   params,
		})
		if err != nil {
			// Already D37-classed by the gate/runtime. As data it reaches the
			// model with its own wording; as a Go error it would be rewritten.
			return tools.Result{Text: err.Error(), IsError: true}, nil
		}
		return tools.Result{
			Text:      out.Text,
			IsError:   out.IsError,
			Truncated: out.Truncated,
		}, nil
	}
}
