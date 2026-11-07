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
// AC#2 LIVES IN THE SHAPE, NOT IN A LABEL. "配置开着但没生效" is the default
// state of every row on this wire, stated once by EngineState rather than
// hidden behind an enabled=true that would read as "running". The sentence is
// the reader's to produce (internal/config's own knowledge of its consumers)
// and the pump's to carry verbatim; neither side derives an "effective" bit,
// because there is no assembly to be effective IN. When the engine lands
// (50/51), its leg updates unwired.go's disposition rows and this sentence in
// the same commit - a load-time honest sentence must not outlive the
// capability that justified it, same rule the guard table itself follows.
//
// THE FLIP HALF IS REGISTERED, NOT BUILT. Ticket 交付 2 (从菜单里开关) needs a
// config write route (internal/config's Manager + the locked-section audit +
// 三档审计), and no write field arrives with this carrier: panel carries no
// config writer beyond the known fields (config_write_*.go's rosters), adding
// one is a C17-adjacent contract move, and the ticket's own 判据 #3 (开关跨档
// 一致) cannot be measured against an engine that does not exist. Registered
// as the third not-wired dimension in docs/evidence/s1/215-skills-plugins-inventory.md;
// the carrier invents no toggle to dress that gap up.
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
// (50/51) or the skills config model (215b, D36) lands, the landing leg
// rewrites the sentence in the same commit - an honest sentence must not
// outlive the absence it describes.
const (
	// PluginEngineNotWired215 is EngineState's canonical sentence: no consumer
	// of the [plugins] keys exists today, so enabled=true means configured,
	// never running.
	PluginEngineNotWired215 = "插件引擎未装配（票 50/51）：[plugins] 各键今天没有任何读者，enabled=true 只表示已配置、不表示在跑"
	// SkillsNotInConfigModel215 is SkillState's canonical sentence: the config
	// model has no skills dimension to enumerate, and adding one is the D36
	// contract move ticket 215b waits on - registered, not rendered as an
	// empty list that would read as "no skills configured".
	SkillsNotInConfigModel215 = "技能维今天没有配置键，无清单可枚举；补技能清单要先动 D36 配置模型（契约面，待批），本清单不列空表"
)

// PluginRow is one [plugins.<id>] entry on the wire.
type PluginRow struct {
	// ID is the plugin id exactly as the section spells it - the <id> of
	// [plugins.<id>] - because the flip half (when a later ticket builds it)
	// must be able to name the same entry the config layer would write.
	ID string `json:"id"`
	// Enabled is plugins.<id>.enabled as the live config holds it - the
	// switch's own words. true is a configuration fact, NOT "running": no
	// engine reads this key today (unwired.go:142 names the row), so an
	// enabled row is precisely ticket 215 AC#2's "配置开着但没生效" and
	// EngineState says so on the same packet instead of letting the boolean
	// dress itself up.
	Enabled bool `json:"enabled"`
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
