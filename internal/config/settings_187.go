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

// ChatSwitchOptions is ticket 187's read face over the chat role: what the
// composer row's model and tier pickers need, read off the in-memory config.
//
// Models lists ONLY enabled catalog entries (ModelSpec.Enabled's promise is
// enforced at resolve time in internal/llm; the read face applies it here so
// the picker cannot offer a row the resolver would refuse). ThinkingLevels
// lists what the CURRENT model declares, verbatim and possibly empty - an
// empty list is "this model declares no tiers", never "the reader failed".
// Reason is the one sentence the empty-catalog shapes carry instead of rows
// (AC#4's explicit "没有可用模型", 宁缺毋造): empty Reason means the catalog
// was read and holds offerable rows.
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
	// ThinkingLevels is the current model's own declared list.
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
				out.ThinkingLevels = append([]string(nil), spec.ThinkingLevels...)
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
