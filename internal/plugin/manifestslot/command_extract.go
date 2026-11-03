package manifestslot

import (
	"encoding/json"

	"github.com/CarlosShao/wisp/internal/plugin"
)

// extractDeclJSON marshals an extract declaration for audit rows and error
// text (never for the child).
func extractDeclJSON(decl map[string]any) string {
	b, err := json.Marshal(decl)
	if err != nil {
		return ""
	}
	return string(b)
}

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
