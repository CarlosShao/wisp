package config

// Ticket 313's settings-write additions: the two voice cloud chains, the two
// realtime voice head leaves, ONE explicit provider/model row-creation step, and
// the chat role pair. Every one of them is a NEW method; no existing setter
// changed shape (ticket 313's leg was sent with "只许追加", 313r1a-send §2).
//
// WHY THEY EXIST. Ticket 313's owner decided that the model surface has to be
// configurable from the settings page, and the 313-a1 census measured that on a
// config the page has just created there is nothing to write: a provider row no
// one declared cannot be addressed (unknownProviderErr), a cloud-chain element
// whose provider/model pair is not in the catalog is refused by validateCatalog,
// the voice section had zero writable fields at all, and llm.roles.chat cannot
// be reached one half at a time because validateRoles rejects BOTH half states
// ("set both or neither", catalog.go's role switch).
//
// THE IMPLICIT/EXPLICIT LINE (A877 ruling 2, and the reason the comment above the
// three refusal tags is NOT amended). That comment says "the page does not mint a
// provider or a catalog row", and it stays literally true: no existing setter
// grew a create branch, so a write that addresses a missing row still fails with
// reason 2 and still names the file as the place a row comes from
// (unknownProviderErr / requireCatalogEntry, unchanged). What ticket 313 adds is
// a SEPARATE step a caller has to ask for BY NAME - CreateProviderRow - with its
// own key path on the receipt; "I noticed the row was missing so I made one"
// remains impossible. Explicit and implicit are two different things and this
// file keeps them apart.
//
// ONE EXPLICIT STEP, CRASH-ATOMIC, NO HALF ROW (A877 ruling 2, boundary 2).
// CreateProviderRow writes the provider section and, when a model id is given,
// that model's section, through the same guarded key-by-key primitive every
// other setter rides (writeOneKey): validate() runs on the base the write would
// produce before a byte is written, and a failure leaves config.toml
// byte-identical. There is no path that can leave a provider row without the
// model row that was asked for in the same call.
//
// PRESETS ONLY. A provider whose name is not a built-in preset must state its
// protocol explicitly (validateAPIKeyRefs), and this step carries no protocol
// field, so it refuses a non-preset name by name rather than writing a section
// that would not load back on the next start. That is the same "one click bricks
// the next start" failure AC#7 exists to stop.

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// SetVoiceCloudASRChain writes voice.cloud_asr_chain (A877 ruling 1). The value
// is the ordered cloud ASR fallback, one "provider/model" element each; an empty
// chain means "no cloud ASR head", which is a legitimate write and not a refusal.
// Every element is checked before the file is touched, so an element naming a
// provider or a model this config does not list comes back as reason 2 (行不存在)
// with the same guidance a single-field catalog write carries.
func (m *Manager) SetVoiceCloudASRChain(chain []string) ([]string, error) {
	return m.setVoiceChain("voice.cloud_asr_chain", chain,
		func(c *Config, els []string) { c.Voice.CloudASRChain = els })
}

// SetVoiceCloudTTSChain writes voice.cloud_tts_chain; see SetVoiceCloudASRChain
// for the element rules (they are identical, and validateCatalog checks both).
func (m *Manager) SetVoiceCloudTTSChain(chain []string) ([]string, error) {
	return m.setVoiceChain("voice.cloud_tts_chain", chain,
		func(c *Config, els []string) { c.Voice.CloudTTSChain = els })
}

func (m *Manager) setVoiceChain(keyPath string, chain []string,
	set func(*Config, []string),
) ([]string, error) {
	els, err := normalizedChainElements(chain)
	if err != nil {
		return nil, err
	}
	return m.writeOneKey(keyPath, func(base *Config) bool {
		set(base, els)
		return true
	}, func(base *Config) error {
		return requireChainElements(keyPath, els, base)
	})
}

// normalizedChainElements trims each element and refuses an empty one instead of
// dropping it: silently skipping an element the caller sent would be a guess
// about intent, and the caller that wants to clear the chain sends an empty
// chain, which is a different value.
func normalizedChainElements(chain []string) ([]string, error) {
	out := make([]string, 0, len(chain))
	for _, el := range chain {
		t := strings.TrimSpace(el)
		if t == "" {
			return nil, observe.New(observe.ClassConfig,
				"config: cloud chain 里有一枚空元素；空元素不会被跳过（跳过＝替你猜），要清空就传空链")
		}
		out = append(out, t)
	}
	return out, nil
}

// requireChainElements is reason 2 for a chain write: every element must name a
// provider the config declares AND a model that provider's catalog lists
// (catalog.go's checkChainElement is the same ruler validateCatalog runs at
// load). The local sherpa cascade is refused by that ruler too - it is the
// terminal fallback and never a chain member.
func requireChainElements(keyPath string, els []string, base *Config) error {
	if base == nil {
		return observe.New(observe.ClassConfig, "config: no config to check cloud chain elements against")
	}
	for _, el := range els {
		if err := checkChainElement(keyPath, el, base.LLM.Providers); err != nil {
			return observe.Wrap(observe.ClassConfig, err,
				"config: refused to write "+keyPath+"。 "+refusalRowMissing+
					"：链元素点名的 provider/model 必须已经在目录里（先用建行那一步，或手改文件）。")
		}
	}
	return nil
}

// SetVoiceRealtimeProvider writes voice.realtime.provider (A877 ruling 1's β
// shape: two single-key fields, never one pair field, so the one-key-per-call
// discipline is not bent). No second rule is invented here: whether a half pair
// is legal is validateRealtime's own answer (a disabled head may carry a stale
// pair; an enabled one may not), and that is the ruler this write runs against.
func (m *Manager) SetVoiceRealtimeProvider(provider string) ([]string, error) {
	v := strings.TrimSpace(provider)
	return m.writeOneKey("voice.realtime.provider", func(base *Config) bool {
		base.Voice.Realtime.Provider = v
		return true
	}, nil)
}

// SetVoiceRealtimeModel writes voice.realtime.model; the provider half keeps its
// own field and its own write (see SetVoiceRealtimeProvider).
func (m *Manager) SetVoiceRealtimeModel(model string) ([]string, error) {
	v := strings.TrimSpace(model)
	return m.writeOneKey("voice.realtime.model", func(base *Config) bool {
		base.Voice.Realtime.Model = v
		return true
	}, nil)
}

// CreateProviderRow is ticket 313's EXPLICIT row-creation step (A877 ruling 2).
// It is the one and only way a provider row (and, when model is non-empty, one
// catalog entry under it) can appear through a settings write, and it has to be
// asked for by name: an existing row is refused rather than rewritten, because
// silently replacing a declared provider would drop the models, billing and
// quotas the operator wrote.
//
// Both sections land through one guarded write, so the receipt names both key
// paths and a failure moves nothing. The row is seeded from the schema's own
// `default:"..."` tags (the single source per D36 rule 3) rather than from zero
// values: a section emitted with an explicit empty billing would override the
// default on the next load, which is not what "create a row" means.
func (m *Manager) CreateProviderRow(provider, model string) ([]string, error) {
	name := strings.TrimSpace(provider)
	id := strings.TrimSpace(model)
	if name == "" {
		return nil, observe.New(observe.ClassConfig,
			"config: "+refusalRowMissing+"：建行必须点名服务商名（本步只接受内置预设名之一）")
	}
	if _, ok := LookupPreset(name); !ok {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: %q 不是内置预设名，这一步不能替它建行：非预设的 [llm.providers.%s] 必须自己写 protocol，"+
				"而这一步没有 protocol 字段，写出来的节下一次加载会不过。 %s：%s",
			name, name, refusalRowMissing, guidanceProviderRow))
	}
	ownedKey := "llm.providers." + name
	return m.writeOneKey(ownedKey, func(base *Config) bool {
		row := Provider{}
		applyDefaults(reflect.ValueOf(&row).Elem())
		row.Models = nil
		if id != "" {
			spec := ModelSpec{}
			applyDefaults(reflect.ValueOf(&spec).Elem())
			row.Models = map[string]ModelSpec{id: spec}
		}
		if base.LLM.Providers == nil {
			base.LLM.Providers = map[string]Provider{}
		}
		base.LLM.Providers[name] = row
		return true
	}, func(base *Config) error {
		if base == nil {
			return observe.New(observe.ClassConfig, "config: no config to create a provider row in")
		}
		if _, exists := base.LLM.Providers[name]; exists {
			return observe.New(observe.ClassConfig, fmt.Sprintf(
				"config: llm.providers.%s already exists; 这一步是「建行」，不是改写已有行——"+
					"已有行请用逐字段写口（base_url / api_key_ref / 目录模型那几枚）", name))
		}
		return nil
	})
}

// SetRoleChatPair writes llm.roles.chat's provider AND model in one guarded
// write. It exists because there is no other way to reach the pair: validateRoles
// refuses BOTH half states ("provider set without model", "model set without
// provider"), so from an unset role the two single-key setters cannot reach a
// legal state - the first one to land always fails validation. The pair is one
// logical unit in this schema ("set both or neither"), so one write owns both
// dotted keys and the receipt names both.
//
// ⚠ This is the one place ticket 313 writes two keys from one field, and it is
// reported as such (A877 ruling 1 kept the one-key-per-call discipline for the
// realtime head, where the β two-field shape DOES work; here it structurally
// cannot). If the ruling for this shape goes the other way, the page-side spec in
// docs/evidence/s1/313-panel-ask.md names exactly what stops working.
func (m *Manager) SetRoleChatPair(provider, model string) ([]string, error) {
	p := strings.TrimSpace(provider)
	id := strings.TrimSpace(model)
	if p == "" || id == "" {
		return nil, observe.New(observe.ClassConfig,
			"config: llm.roles.chat 的 provider 与 model 必须一起点名（set both or neither）；半对不会写进文件")
	}
	return m.writeOneKey("llm.roles.chat", func(base *Config) bool {
		base.LLM.Roles.Chat.Provider = p
		base.LLM.Roles.Chat.Model = id
		return true
	}, func(base *Config) error {
		return requireCatalogEntry(base, p, id)
	})
}
