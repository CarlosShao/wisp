package panel

// Ticket 187's carrier: the composer row's model and thinking-tier dimension
// as a snapshot section.
//
// WHAT THIS IS. The composer row already carries composer.currentModel - the
// model id the RUNNING endpoint resolved at assembly. What ticket 187's picker
// needs besides that is the switchable half: the catalog rows a picker may
// offer, the tier the chat role carries for the NEXT task, and the levels the
// current model itself declares. This file is the wire shape for that half;
// the values are read off the live config by the assembly root through the
// reader this package declares on PumpSources (the same division ticket 248's
// credential dimension rides: internal/panel does not import internal/config,
// the reader arrives already shaped).
//
// THE TIER IS THE SECTION'S OWN FACT, AND IT IS PER-FIELD. Measured, not
// assumed (ticket 187 AC#1): the chat role's model is resolved ONCE at
// assembly (cmd/wisp run.go builds the provider chain from it and never
// rebuilds), so a model switch is restart tier; the thinking tier is
// re-resolved PER TASK from the live config (run.go's per-task
// ResolveRole over rt.mgr.Config()), so a tier write is next_task. The
// composer.currentModel / modelSwitch.current pair is exactly how the
// "displayed value vs the value in force" drift of ticket 146's family is
// kept HONEST instead of hidden: until a restart the two may differ, and
// the wire carries both so the page can say so rather than lie.
//
// READER-GATED, THE 213 PATTERN, NO SECOND PARADIGM. A pump assembled
// without the reader sends no "modelSwitch" key at all - the shape the
// four-key byte nails pin for the no-reader assembly - and a pump that HAS
// one always sends the key, empty lists rendered as []. Nothing here opens
// a door: the section names what may be asked, the ask itself rides
// config.set (the settings route ticket 248 opened) with the new
// role_chat_thinking field, and NO C17 method is added by this carrier.
//
// 宁缺毋造, AC#4's shape on the read half: Models lists only what the
// catalog actually holds enabled; ThinkingLevels lists only what the
// CURRENT model declares (possibly fewer than four, possibly none); and
// the no-offerable-rows shapes - no providers at all, every row disabled,
// config unreadable - travel as Reason, never as an empty list dressed up
// as "nothing to choose".

// ModelOption is one offerable catalog row: the pair a role_chat_pair write
// addresses. Provider plus model id, nothing else - no price card, no
// capability bits, and by construction no credential-shaped string.
type ModelOption struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

// ModelSwitchSection is the snapshot's modelSwitch key.
type ModelSwitchSection struct {
	// Current is the chat role's pair as the live config holds it RIGHT NOW
	// - what a restart (model writes) or the next task (tier writes) will
	// act on. Empty both = the role is unset; that is a state, never a
	// default silently filled in.
	Current ModelOption `json:"current"`
	// CurrentThinking is llm.roles.chat.thinking_intensity verbatim; ""
	// means the file carries no tier key (the schema default spells "off"
	// on disk) - the wire says what is there, not what a default would be.
	CurrentThinking string `json:"currentThinking"`
	// Models is the offerable catalog, sorted by the config side. Empty
	// WITH a Reason is the measured "no rows can be offered"; empty without
	// one would be a lie, so the config read face always sets Reason when
	// Models comes back empty.
	Models []ModelOption `json:"models"`
	// ThinkingLevels is what the CURRENT model declares, verbatim: the
	// picker may offer only these (possibly fewer than four; empty means
	// the model declares no tiers, and the page then offers no tier
	// choices - it does not invent the four-value vocabulary).
	ThinkingLevels []string `json:"thinkingLevels"`
	// Reason is the one sentence the no-offerable-rows shapes carry
	// (no [llm.providers] section / every row disabled / config unreadable).
	// Empty = the catalog was read and holds offerable rows.
	Reason string `json:"reason,omitempty"`
}
