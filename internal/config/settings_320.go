package config

// Ticket 320's settings-write additions: the provider billing leaf (AC#1), the
// non-preset provider row (AC#2), and (below, same file) the catalog-entry step
// its AC#3 asks for.
//
// WHY THEY EXIST. Ticket 320's owner asked for three writes the model page v3
// cannot make today (ledger A968): the first is a leaf the schema already had but
// no setter owned; the second is the non-preset twin of ticket 313's
// CreateProviderRow. The explicit/implicit line settings_313.go's header states
// is unchanged: no older setter grew a create branch, and a row-creating step is
// asked for BY NAME, with its own key path on the receipt.
//
// SAME SHAPE AS THE SETTERS BEFORE THEM. Every method here owns exactly one key
// path (or, for the row creators, one section whose name it owns) and rides
// writeOneKey, so validate() runs on the base the write would produce before a
// byte is written and every refusal leaves config.toml byte-identical.
//
// THE ENUMS ARE CHECKED HERE, NOT ONLY AT LOAD. Each vocabulary gate below is a
// membership check over the same list validateAPIKeyRefs runs, so a value the
// schema would reject comes back as a sentence naming the field instead of as a
// whole-file validation failure the operator has to localize.

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// SetProviderBilling writes llm.providers.<name>.billing (ticket 320 AC#1). The
// value is one of the schema's two billing modes; anything else is refused before
// the file is touched, because a third spelling would fail validate() on the next
// load and turn a settings click into an Unconfigured start.
func (m *Manager) SetProviderBilling(provider, billing string) ([]string, error) {
	v := strings.TrimSpace(billing)
	if !isBillingMode(v) {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: llm.providers.%s.billing: %q 不在枚举内（%s）；这一枚只接受这几种写法",
			provider, v, vocabList(billingModes)))
	}
	return m.writeOneKey(keyPathProviderField(provider, "billing"), func(base *Config) bool {
		p, ok := base.LLM.Providers[provider]
		if !ok {
			return false
		}
		p.Billing = v
		base.LLM.Providers[provider] = p
		return true
	}, func(base *Config) error {
		if _, ok := base.LLM.Providers[provider]; !ok {
			return unknownProviderErr(provider)
		}
		return nil
	})
}

// isBillingMode asks the schema's own list, never a copy of the two strings.
func isBillingMode(v string) bool {
	for _, m := range billingModes {
		if m == v {
			return true
		}
	}
	return false
}

// CreateCustomProviderRow is ticket 320 AC#2: the EXPLICIT creation step for a
// provider whose name is not one of the built-in presets. It is the non-preset
// twin of ticket 313's CreateProviderRow, and it exists because the two halves
// that step is missing are exactly the two a non-preset row must carry:
// validateAPIKeyRefs refuses a non-preset section with no protocol, and a row
// with no base_url addresses nothing.
//
// The row is seeded from the schema's own `default:"..."` tags first (the same
// rule CreateProviderRow follows), then protocol and base_url are written over
// the seed, and the optional model id - when given - is one catalog entry seeded
// the same way. All of that is ONE guarded write: the receipt names the section
// the file gained, and a failure moves nothing.
//
// FOUR REFUSALS, EACH BY NAME, because four different things are wrong four
// different ways: a preset name (provider_row_create owns those), an empty
// base_url, a protocol outside the enum, and a row that already exists (this step
// creates; it is not a rewrite tool - silently replacing a declared provider
// would drop the billing, quotas and catalog the operator wrote). An empty name
// is the fifth, and it is refusal 2's shape: there is no row to create.
func (m *Manager) CreateCustomProviderRow(provider, protocol, baseURL, model string) ([]string, error) {
	name := strings.TrimSpace(provider)
	proto := strings.TrimSpace(protocol)
	url := strings.TrimSpace(baseURL)
	id := strings.TrimSpace(model)
	if name == "" {
		return nil, observe.New(observe.ClassConfig,
			"config: "+refusalRowMissing+"：建行必须点名服务商名")
	}
	if _, ok := LookupPreset(name); ok {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: %q 是内置预设名，预设行请走 provider_row_create（那一步按预设带出 protocol 与 base_url）；"+
				"本步只建非预设行，并且要求显式写 protocol", name))
	}
	if url == "" {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: 建 [llm.providers.%s] 的值必须是「协议,base_url」两半，base_url 空着不会被写成默认：%s",
			name, refusalRowMissing))
	}
	if !isProtocol(proto) {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: protocol %q 不在枚举内（%s）；非预设行必须显式写 protocol，否则这份文件下一次加载会不过",
			proto, vocabList(protocols)))
	}
	ownedKey := "llm.providers." + name
	return m.writeOneKey(ownedKey, func(base *Config) bool {
		row := Provider{}
		applyDefaults(reflect.ValueOf(&row).Elem())
		row.Protocol = proto
		row.BaseURL = url
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
				"config: llm.providers.%s already exists；这一步是「建行」，不是改写已有行——"+
					"已有行请用逐字段写口（protocol / base_url / 目录模型那几枚）", name))
		}
		return nil
	})
}

// isProtocol asks the schema's own vocabulary, never a copy of the three names.
func isProtocol(v string) bool {
	for _, p := range protocols {
		if p == v {
			return true
		}
	}
	return false
}
