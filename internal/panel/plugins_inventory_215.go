package panel

// Ticket 215's carrier: the snapshot's plugins key - the Tier-2 plugin
// inventory the live config holds, each row carrying what the config layer
// can actually prove about it today. Ticket 215's own 09-28 correction splits
// the ticket in two halves and this carrier takes exactly one of them:
//
//   - 215a, the PLUGIN half: the config model HAS the dimension
//     (internal/config's PluginsSection - tier2_enabled plus dynamic
//     [plugins.<id>] sub-tables, 8 schema hits), so the inventory is
//     enumerable from a real source and lands here.
//   - 215b, the SKILL half: the config model has NO key at all
//     (grep -c Skill internal/config/schema.go = 0, measured 2026-11-07) -
//     adding one means touching D36's section tree, which is a contract move
//     needing human approval (AGENTS.md sec 0.2). This carrier registers that
//     absence as a sentence on the wire (SkillState) instead of rendering an
//     empty list that would read as "no skills configured" - 宁缺毋造.
//
// WHAT TODAY'S READER CAN HONESTLY CARRY, AND WHAT IT MUST NOT INVENT.
// The plugin engine (tickets 50/51) is not built: nothing in the repository
// reads PluginsSection.Entries, which is exactly why internal/config's
// unwired.go lockedKeyDisposition names the [plugins] rows "not built" and
// hands them to the orchestrator (2026-11-07 reading: unwired.go:140-145, five
// guard rows in unwiredKeys unchanged). So the inventory lists what the
// config SAYS - the master switch and each entry's own four fields, verbatim -
// and says once per section what the machine can DO about it: nothing. A row
// showing "version 1.0.0, loaded" would be the lying key ticket 215 was filed
// to stop: the schema has no version field (measured: no Version field on
// PluginEntry), and no engine exists to load one.
//
// AC#2 LIVES IN THE SHAPE AND NOW IN A PER-ROW BIT. "配置开着但没生效" is the
// state of every row until the engine lands, and since the 215a write leg it
// is carried TWICE: once by EngineState (the section-level sentence) and once
// by each row's State (the three-state label the menu renders). The State bit
// is DERIVED BY THE READER from the assembly probe the composition root hands
// in - never from Enabled - so "enabled=true" alone can never render as
// effective; the derivation lives in cmd/wisp (the only place that can know
// what this process actually assembled) and its teeth are there.
//
// THE FLIP HALF RIDES THE EXISTING SETTINGS DOOR. Ticket 交付 2 (从菜单里开关)
// landed in the 215a write leg: the write is the settings route ticket 248
// opened (config.set), carrying ONE named field (plugin_enabled, roster
// config_write_215.go) whose selector names the plugin id; the leg it reaches
// is internal/config's own protected section-write entry (Manager.
// SetPluginEnabled over WriteSection), so the D36 direction verdict runs
// before any byte moves: switching OFF is a tighten and writes immediately,
// switching ON from the panel is a loosening and is REFUSED there (the panel
// is a permission input, not an authority - sectionwrite.go's structural
// rule), with the audit line written either way. No new method arrives with
// this ticket and no C17 method name changes.
//
// READER-GATED, THE 213/187/65 PATTERN, NO SECOND PARADIGM. A pump assembled
// without the reader sends no "plugins" key at all - the shape the four-key
// byte nails pin - and a pump that HAS one always sends the key. Nothing here
// opens a door and no C17 method arrives with it.
//
// ORDERING IS THE PUMP'S, CONTENT IS THE READER'S. PluginsSection.Entries is
// a Go map, so the reader cannot promise an order; the pump sorts rows by ID
// so two publishes of the same config marshal byte-identically (a menu whose
// row order reshuffles every push is an audit defect, not a style choice).
// Row CONTENT travels verbatim - the pump normalizes order and nil-ness, and
// invents nothing else.

// The two honest-state sentences, in exactly the spelling a production reader
// sends them. They live here (not at the assembly root) so the reader that
// finally wires this key has ONE canonical sentence per dimension instead of
// forked prose, and so the carrier tests can pin the wire against these very
// constants. Both sentences share unwired.go's lifecycle: when the engine
// assembly (50/51) or the skills config model (215b, D36) lands, the landing
// leg rewrites the sentence in the same commit - an honest sentence must not
// outlive the absence it describes. PluginEngineNotWired215's spelling was
// rewritten in the 215a write leg: the old text said the [plugins] keys have
// NO readers, which stopped being true when the inventory itself and the
// settings route started reading them - what is still true, and what the
// sentence now says, is that the ENGINE is not wired at the assembly root.
const (
	// PluginEngineNotWired215 is EngineState's canonical sentence: the Tier-2
	// engine components exist in the tree but no assembly root wires them, so
	// enabled=true means configured, never running, and each row's State is
	// the assembly probe's answer, not the boolean's.
	PluginEngineNotWired215 = "插件引擎未在装配根接线（票 50/51 组件已在库、装配根未接）：enabled=true 只表示已配置；每行的生效位由装配探针判，探针上没有这枚＝未生效"
	// SkillsNotInConfigModel215 is SkillState's canonical sentence: the config
	// model has no skills dimension to enumerate, and adding one is the D36
	// contract move ticket 215b waits on - registered, not rendered as an
	// empty list that would read as "no skills configured".
	SkillsNotInConfigModel215 = "技能维今天没有配置键，无清单可枚举；补技能清单要先动 D36 配置模型（契约面，待批），本清单不列空表"
)

// The three states a row can honestly carry, in exactly the spelling the
// menu renders them (the page shows the wire value verbatim - one spelling,
// no per-page copy that can drift). The derivation is the composition
// root's, from its own assembly probe; these constants only fix the words.
const (
	// PluginRowLive215: this round's assembly actually carries the plugin -
	// the probe answered with this id AND the current config still matches
	// what was assembled.
	PluginRowLive215 = "本轮生效"
	// PluginRowRestart215: the probe carries this id but the config has moved
	// away from what was assembled (e.g. switched off after assembly), and
	// the consumer only re-reads config at restart - the row says so instead
	// of pretending the old state already ended.
	PluginRowRestart215 = "需重启"
	// PluginRowNotLive215: enabled=true on the config and the probe has
	// nothing for this id - the exact "配置开着但没生效" shape A408 names
	// (and the honest state of every row until the engine lands), plus the
	// plain "switched off" rows.
	PluginRowNotLive215 = "未生效"
)

// PluginRow is one [plugins.<id>] entry on the wire.
type PluginRow struct {
	// ID is the plugin id exactly as the section spells it - the <id> of
	// [plugins.<id>] - because the flip half (when a later ticket builds it)
	// must be able to name the same entry the config layer would write.
	ID string `json:"id"`
	// Name is the installed plugin's manifest name (plugins\<id>\manifest.toml,
	// read by the composition root's manifest reader). Empty means "no
	// install on disk for this id" - a measured fact, never invented: the
	// installer (manifestslot) has no production caller yet, so a config row
	// without an install is the NORMAL state and must not be dressed up with
	// a made-up name (宁缺毋造, the rule the carrier's header states).
	Name string `json:"name"`
	// Version is the manifest's version string, read with Name from the same
	// manifest bytes and empty for the same reason. Both come from the
	// manifest's PARSED fields only - the reader extracts two typed strings
	// and nothing else, so manifest prose can never ride the wire as config.
	Version string `json:"version"`
	// Note is the manifest reader's refusal line for a plugin whose manifest
	// exists on disk but refused to load (unreadable, C14-refused, or an id
	// that does not match its directory). Empty means either "no install on
	// disk" (Name/Version empty too) or "manifest parsed clean" (Name/
	// Version filled): the row never dresses a refusal up as an identity.
	Note string `json:"note,omitempty"`
	// Enabled is plugins.<id>.enabled as the live config holds it - the
	// switch's own words. true is a configuration fact, NOT "running": the
	// engine is not wired at the assembly root, so an enabled row is
	// precisely ticket 215 AC#2's "配置开着但没生效" and BOTH the State bit
	// below and EngineState say so on the same packet instead of letting the
	// boolean dress itself up.
	Enabled bool `json:"enabled"`
	// State is the row's honest 生效与否 label, one of the three constants
	// above, derived by the reader from the composition root's assembly
	// probe - never from Enabled alone. A packet whose State contradicts
	// Enabled without a probe behind it is the A408 family; the cmd/wisp
	// tests pin the derivation in both directions (probe empty => not live
	// even when enabled; probe present => live).
	State string `json:"state"`
	// Capabilities is plugins.<id>.capabilities, verbatim, in file order -
	// the grant names the user wrote, which nothing consumes today.
	Capabilities []string `json:"capabilities"`
	// NetAllowlist is plugins.<id>.net_allowlist, verbatim, in file order -
	// the private egress allowlist the config grants and no egress leg reads.
	NetAllowlist []string `json:"netAllowlist"`
	// HostAPI is plugins.<id>.host_api, verbatim, in file order - the C24
	// surfaces the config grants and no host exposes yet.
	HostAPI []string `json:"hostApi"`
}

// PluginInventorySection is the snapshot's plugins key (ticket 215a).
type PluginInventorySection struct {
	// Tier2Enabled is plugins.tier2_enabled as the live config holds it. A
	// true here is still not "plugins are running" - the same unwired.go rows
	// own that distinction for the master switch as for the entries.
	Tier2Enabled bool `json:"tier2Enabled"`
	// Rows is one row per [plugins.<id>] entry. Empty means the config
	// declares no plugin entries - a measured fact, never "unknown": the
	// reader either read the section or was not handed one at all, and the
	// nil-reader case is the pump's key-absence rule, not an empty list.
	Rows []PluginRow `json:"rows"`
	// EngineState is the one sentence that keeps enabled=true from lying: who
	// consumes these keys today (nobody - the engine is ticket 50/51) and
	// therefore what "enabled" means on this wire. Produced by the reader,
	// carried verbatim; a packet with rows and a blank EngineState is the
	// 文案在、控件不在 family A408 names, and the carrier tests pin the
	// sentence's presence on the wire.
	EngineState string `json:"engineState"`
	// SkillState is the skill half's honest registration (215b): the config
	// model has no skills dimension at all (grep -c Skill
	// internal/config/schema.go = 0), so there is no skill list to enumerate
	// and adding one is a D36 contract move awaiting its own approval. The
	// sentence says that instead of an empty list that would read as "no
	// skills configured".
	SkillState string `json:"skillState"`
}
