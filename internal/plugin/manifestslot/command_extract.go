package manifestslot

import (
	"encoding/json"

	"github.com/CarlosShao/wisp/internal/plugin"
)

// commandSchemaFor is the model-facing parameter schema of a command tool: the
// manifest's schema narrowed to the placeholders the argv template substitutes
// (plugin.ParameterSchema owns that rule).
func commandSchemaFor(def plugin.ToolDef) json.RawMessage {
	if !def.IsCommand() {
		b, _ := json.Marshal(def.Parameters)
		return b
	}
	b, err := json.Marshal(plugin.ParameterSchema(def))
	if err != nil {
		b, _ = json.Marshal(def.Parameters)
	}
	return b
}
