package config

// Ticket 187's config-side half: the thinking-intensity write leg and the
// read helper the composer-row model/tier picker needs.
//
// WHAT THIS FILE IS. Two setters-shaped facts about [llm.roles.chat]:
//
//	SetRoleChatThinking writes llm.roles.chat.thinking_intensity through the
//	same guarded key-by-key path every settings setter uses (writeOneKey:
//	pre-read the file, validate the base the write would produce, one key,
//	receipt names what the file changed by). The enum is enforced HERE, not
//	on the panel side - the same ruling ticket 320 made for provider_billing
//	("the config setter owns the vocabulary check so the page is told which
//	field it got wrong").
//
//	ChatSwitchOptions is the read face: the pair the chat role names right
//	now, the thinking tier it carries, the catalog entries a picker may
//	offer (enabled rows only), and the levels the CURRENT model itself
//	declares. It reads the in-memory config (Manager.Config()), which is
//	what a snapshot reader should read - the file and memory are kept equal
//	by commitSectionWrite, and the reader must answer with what this
//	process would act on, not with what a re-read of the file would guess.
//
// THE TIER, measured not assumed (ticket 187 AC#1). The thinking tier's
// reader is cmd/wisp/run.go's per-task re-resolution
// (`llm.NewResolver(cfg, nil).ResolveRole(llm.RoleChat)` over the LIVE
// `rt.mgr.Config()`), and commitSectionWrite applies a write to the
// in-memory config first - so a thinking write IS picked up by the next
// task: next_task tier. The chat role's MODEL, by contrast, is resolved
// once at assembly (run.go builds rt.provs from it and never rebuilds):
// a model write is restart tier. That is why the two fields ride the same
// route but do not share one receipt sentence, and why the panel's 187
// roster carries a per-field tier instead of trusting the leg's blanket
// restart stamp.
//
// WHAT THIS FILE DOES NOT DO: no permission-mode surface of any kind. The
// direction verdict in sectionwrite.go is neutral for [llm] keys and this
// file adds no second one; there is no import of the permission or risk
// packages here, and the panel-side nail (config_write_187_test.go and
// settings_187_test.go) keeps it that way by reading this file's bytes.

import (
	"fmt"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// keyPathChatThinking is the one dotted key path SetRoleChatThinking owns.
const keyPathChatThinking = "llm.roles.chat.thinking_intensity"

// SetRoleChatThinking writes llm.roles.chat.thinking_intensity. The value is
// checked against the schema's own thinking vocabulary BEFORE the write, and -
// when the chat role resolves to a catalog model that declares its levels -
// against THAT list too: a panel that offered "high" for a model declaring
// off|low would be offering a request the model cannot answer (ticket 187
// AC#4's 宁缺毋造, applied to the write half).
//
// A model that declares NO thinking_levels constrains nothing (the schema
// defaults the key empty), so the four-value enum alone decides in that shape;
// the picker on the page lists nothing for such a model, which is the same
// measured asymmetry stated on the read half.
//
// THE CAPABILITY BIT DOES NOT VETO HERE - ticket 416's ruling on R-145-3,
// decided on this leg, not inherited. The two faces answer different questions
// and ticket 145 AC#2b already split them: the READ face's question is "does
// this dimension exist for this model, so the page may offer it"
// (thinkingDimensionDeclared answers it with the bit AND the declared list,
// because SPEC-03 section 3.1's probe-first all-false default makes an
// undeclared bit UNKNOWN, and unknown offers nothing - 宁缺毋造); the WRITE
// face's question is "is this stored value legal", and the vocabulary alone
// decides that. Three measured facts close it this way:
//
//   - thinking_intensity is a REQUEST parameter the role carries, not a
//     property the model owns, and the rest of the chain never lets the bit
//     near the value: the resolver carries it verbatim
//     (internal/llm/resolver.go:204) and each adapter validates it against
//     its own vocabulary at request time
//     (internal/llm/openairesponses/request.go:139,
//     internal/llm/anthropic/request.go:124). A setter veto would be the
//     chain's only bit gate on a VALUE, and it would gate nothing:
//     config.toml is hand-editable, so the file door accepts what the setter
//     door refuses - the gate would only divide the two doors.
//   - What a veto would cost buys nothing: the model switcher never touches
//     this key, so a tier rides across models; on a no-bit model every
//     re-assert (even writing "off" back) would be refused while the stored
//     value sat in the file regardless. Storage is where the user's explicit
//     choice lives; the page simply does not offer the option again
//     (MAY-store, never-offer - the same split every mainstream harness
//     draws between a persisted preference and a filtered picker).
//   - So the shape already here IS the mainstream shape: capability gates
//     the OFFER (read face), the vocabulary gates the VALUE (this setter),
//     the adapter gates the REQUEST. Each face answers one question, and the
//     bit decides only the first.
//
// The declared-list check above stays, for a VALUE reason not a bit reason: a
// non-empty ThinkingLevels is the entry's own named vocabulary, and a value
// its own words contradict is refused wherever the write came from - the
// refusal tracks the entry's vocabulary, never its capability bits.
// Round-trip teeth for this whole ruling:
// internal/config/settings_187_r1453_test.go. Registered as R-145-3 in
// docs/evidence/s1/145-ac2b-r1.md; closed by this file, not by a doc.
func (m *Manager) SetRoleChatThinking(level string) ([]string, error) {
	v := strings.TrimSpace(level)
	if v == "" {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: %s: an empty tier is not a tier; the vocabulary is %s. "+
				"第 3 种拒因：校验不过", keyPathChatThinking, vocabList(thinkingIntensity)))
	}
	if !slicesContains(thinkingIntensity, v) {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: %s %q must be one of %s. 第 3 种拒因：校验不过",
			keyPathChatThinking, v, vocabList(thinkingIntensity)))
	}
	return m.writeOneKey(keyPathChatThinking, func(base *Config) bool {
		base.LLM.Roles.Chat.ThinkingIntensity = v
		return true
	}, func(base *Config) error { return requireDeclaredThinkingLevel(base, v) })
}

// requireDeclaredThinkingLevel refuses a tier the chat role's own model does
// not declare, when it declares any. It runs on the base the write would
// produce, so the check is against the file's truth, not this process's
// memory (the same base rule every pre-write check here follows).
func requireDeclaredThinkingLevel(base *Config, level string) error {
	if base == nil {
		return observe.New(observe.ClassConfig, "config: no config to check a thinking tier against")
	}
	role := base.LLM.Roles.Chat
	if role.Provider == "" || role.Model == "" {
		// No pair to resolve: the enum above is the whole vocabulary that
		// applies, and validateCatalog's pair rules stay the pair rules.
		return nil
	}
	p, ok := base.LLM.Providers[role.Provider]
	if !ok {
		// The pair is broken on the base already; validate() refuses the write
		// for that reason, so this check adds no second sentence for it.
		return nil
	}
	spec, ok := p.Models[role.Model]
	if !ok || len(spec.ThinkingLevels) == 0 {
		return nil
	}
	if slicesContains(spec.ThinkingLevels, level) {
		return nil
	}
	return observe.New(observe.ClassConfig, fmt.Sprintf(
		"config: %s: model %s/%s declares thinking_levels %s, and %q is not one of them; "+
			"第 3 种拒因：校验不过（值要改成这一枚模型自己声明的档）",
		keyPathChatThinking, role.Provider, role.Model, vocabList(spec.ThinkingLevels), level))
}

// slicesContains is slices.Contains spelled locally so this file keeps the
// package's existing helper shape (validate.go uses the slices import; the
// settings files have avoided the import so far - one tiny local keeps that).
func slicesContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// ChatModelOption is one catalog entry a picker may offer: the pair a
// role/pair write addresses, plus the entry's own human-facing Display name
// (ticket 145 AC#2b①: ModelSpec.Display verbatim - the source the census
// named; empty stays empty, the page falls back to the model id and this
// reader invents no name). No credential-shaped string can appear here by
// construction.
type ChatModelOption struct {
	Provider string
	Model    string
	// Display is the catalog entry's display name, verbatim from
	// ModelSpec.Display; "" means the entry declares none.
	Display string
}

// thinkingDimensionDeclared answers ticket 145 AC#2b's FIRST question for one
// catalog entry: does this model have a thinking dimension at all?
//
// The ticket splits the tier question in two and forbids answering them with
// one judgement: Capabilities.Thinking decides WHETHER the dimension exists,
// the adapter's vocabulary decides WHICH tiers it holds. This function is the
// first half, and it is a conjunction of two declarations the entry itself
// makes:
//
//   - Capabilities.Thinking is the authority on existence. SPEC-03 section 3.1
//     makes every capability bit default to false and probe-first ("新增模型
//     默认 capabilities 全 false + 引导探测"), so an entry that never declared
//     the bit is not "probably able to think" - it is UNKNOWN, and unknown
//     offers nothing (宁缺毋造). The same shape already governs the speech
//     chain, where a missing AudioIn/AudioOut bit removes the chain element
//     outright (internal/speech/cloudcatalog.go:95,99).
//   - A non-empty ThinkingLevels is the entry's own list of tiers. A model
//     that declares the bit but no tiers says "thinking exists, and I name no
//     tier for it" - the adapter then maps nothing, and an empty offer is the
//     truthful answer rather than a guess at four enum values.
//
// The bit VETOES: levels declared under a false bit are dropped, not trusted.
// That is deliberate - it is the only order in which a contradiction inside
// one catalog entry (levels declared, capability never declared) resolves
// toward "do not offer a request we cannot justify", and it is what makes the
// all-false default safe rather than accidentally permissive.
func thinkingDimensionDeclared(spec ModelSpec) bool {
	if !spec.Capabilities.Thinking {
		return false
	}
	return len(spec.ThinkingLevels) > 0
}

// ChatSwitchOptions is ticket 187's read face over the chat role: what the
// composer row's model and tier pickers need, read off the in-memory config.
//
// Models lists ONLY enabled catalog entries (ModelSpec.Enabled's promise is
// enforced at resolve time in internal/llm; the read face applies it here so
// the picker cannot offer a row the resolver would refuse). ThinkingLevels
// lists what the CURRENT model declares, verbatim and possibly empty - an
// empty list is "this model has no thinking dimension" (thinkingDimensionDeclared
// above: the capability bit or the declared tiers said no), never "the reader
// failed". Reason is the one sentence the empty-catalog shapes carry instead of
// rows (AC#4's explicit "没有可用模型", 宁缺毋造): empty Reason means the
// catalog was read and holds offerable rows.
type ChatSwitchOptions struct {
	// CurrentProvider / CurrentModel are llm.roles.chat as this process holds
	// it right now. Empty both = the role is unset (text_chain would carry the
	// model at resolve time); that state is reported, not guessed at.
	CurrentProvider string
	CurrentModel    string
	// CurrentThinking is roles.chat.thinking_intensity verbatim; "" = unset,
	// which the schema's default spells "off" on disk but which this reader
	// does not invent - the wire carries what the file holds.
	CurrentThinking string
	// CurrentDisplay is the CURRENT model's own Display name as the catalog
	// entry declares it (ticket 145 AC#2b①), found on the same walk that
	// fills Models; "" when the pair names no catalog entry - the reader
	// invents no name for a pair the catalog does not hold.
	CurrentDisplay string
	// Models is the offerable list, sorted by provider then model so a page
	// cannot read order as preference.
	Models []ChatModelOption
	// ThinkingLevels is the current model's own declared list, and it is
	// EMPTY unless thinkingDimensionDeclared says the model has the dimension
	// at all (ticket 145 AC#2b: the bit decides existence, the vocabulary
	// decides which tiers). An empty list is an answer, not a failure.
	ThinkingLevels []string
	// Reason explains the no-offerable-rows shapes; empty otherwise.
	Reason string
}

// ChatSwitchOptions reads the picker's data source off the live config.
func (m *Manager) ChatSwitchOptions() ChatSwitchOptions {
	var out ChatSwitchOptions
	cfg := m.Config()
	if cfg == nil {
		out.Reason = "配置不可读：没有可用模型（这不是「目录为空」）"
		return out
	}
	out.CurrentProvider = cfg.LLM.Roles.Chat.Provider
	out.CurrentModel = cfg.LLM.Roles.Chat.Model
	out.CurrentThinking = cfg.LLM.Roles.Chat.ThinkingIntensity

	if len(cfg.LLM.Providers) == 0 {
		out.Reason = "配置里没有任何 [llm.providers] 节：没有可用模型（行在配置文件里建，面板不替你造）"
		return out
	}

	// The declared levels of the CURRENT model, found on the same walk that
	// collects the offerable rows - one pass, so the two answers cannot be
	// read from different moments.
	providers := make([]string, 0, len(cfg.LLM.Providers))
	for name := range cfg.LLM.Providers {
		providers = append(providers, name)
	}
	sort.Strings(providers)
	for _, pname := range providers {
		p := cfg.LLM.Providers[pname]
		ids := make([]string, 0, len(p.Models))
		for id := range p.Models {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			spec := p.Models[id]
			if pname == out.CurrentProvider && id == out.CurrentModel {
				// Ticket 145 AC#2b defect (1): the dimension EXISTS only
				// when the entry declares the bit AND names a tier. Both
				// halves are the entry's own words; neither is inferred.
				if thinkingDimensionDeclared(spec) {
					out.ThinkingLevels = append([]string(nil), spec.ThinkingLevels...)
				}
				// Ticket 145 AC#2b①: the current model's display name rides
				// the same walk - one pass, so the name and the pair cannot
				// be read from different moments.
				out.CurrentDisplay = spec.Display
			}
			if !spec.Enabled {
				continue
			}
			out.Models = append(out.Models, ChatModelOption{Provider: pname, Model: id, Display: spec.Display})
		}
	}
	if len(out.Models) == 0 {
		out.Reason = "目录里的条目都处于 enabled=false：没有可用模型（在条目里写 enabled = true 即可重新启用）"
	}
	return out
}
