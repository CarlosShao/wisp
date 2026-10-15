package config

// Ticket 320's settings-write additions: the provider billing leaf (AC#1), and
// (below, same file) the row-creation steps its AC#2/AC#3 ask for.
//
// WHY IT EXISTS. Ticket 320's owner asked for three writes the model page v3
// cannot make today (ledger A968); this first one is a leaf the schema already
// had but no setter owned: the provider's billing mode.
//
// SAME SHAPE AS THE SETTERS BEFORE IT. The method here owns exactly one key path
// and rides writeOneKey, so validate() runs on the base the write would produce
// before a byte is written and every refusal leaves config.toml byte-identical.
// The explicit/implicit line settings_313.go's header states is unchanged: no
// older setter grew a create branch.
//
// THE ENUM IS CHECKED HERE, NOT ONLY AT LOAD. billing is a membership check over
// the same vocabulary validateAPIKeyRefs runs, so a value the schema would reject
// comes back as a sentence naming the field instead of as a whole-file validation
// failure the operator has to localize.

import (
	"fmt"
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
