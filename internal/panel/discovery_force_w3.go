package panel

// Readface wave 3's wire shapes (2026-11-17, ticket 38 AC#3 / AC#5). They live
// beside the envelope for the same reason the approval card's view does: the
// page parses these keys, so the vocabulary is declared ONCE, next to the
// roster that answers them - never re-invented in the handler.

// ToolRosterRow is one tools.list row: the registry's own truth per entry.
// RiskPill is rendered by internal/risk's Level vocabulary (L0/L1/L2 - the
// same words the confirmation card shows); no handler may spell a level word.
type ToolRosterRow struct {
	// Name is the C1 'namespace.action' name the page echoes back on
	// tool.launch.
	Name string `json:"name"`
	// Title is the tool description's first line - what the model itself
	// reads to decide whether to call the tool.
	Title string `json:"title"`
	// RiskPill is the entry's declared risk level (the R1 floor the D34
	// table pins), rendered by internal/risk's truth. It is a FLOOR, never a
	// verdict: the per-call C19 judge still owns the real gate.
	RiskPill string `json:"riskPill"`
	// ArgSummary is the JSON Schema's top-level property names, comma
	// separated ("无参数" when the tool takes none).
	ArgSummary string `json:"argSummary"`
}

// ForceRuleRow is one force rule the rules.force face reports.
type ForceRuleRow struct {
	// Source is "toml" (a [[force_rule]] row of force-rules.toml) or
	// "override" (the palette chip; Pattern carries the task id it rode).
	Source string `json:"source"`
	// Pattern is the toml rule's match pattern, or the override's task id.
	Pattern string `json:"pattern"`
	// Mode is force_tool | force_chat (D11③'s own two words).
	Mode string `json:"mode"`
}

// ForceRulesView is the rules.force payload: the effective rules plus the
// mode the NEXT prompt will actually adopt for the requested task text - the
// same ForceDecision the prompt injection consumed (one resolver, no drift).
type ForceRulesView struct {
	// Mode is "force_tool" | "force_chat", or "" for the default (no force).
	Mode string `json:"mode"`
	// Source is "override" | "toml" | "default".
	Source string `json:"source"`
	// Pattern is the winning toml rule's pattern ("" for override/default).
	Pattern string `json:"pattern"`
	// Instruction is the exact block the next prompt will carry ("" = the
	// prompt is unchanged).
	Instruction string `json:"instruction"`
	// Rules is every effective rule with its source (override rows first
	// after the stable sort; the face does not rank them).
	Rules []ForceRuleRow `json:"rules"`
}
