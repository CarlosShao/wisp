package config

// Ticket 40 AC#3's quota-editing setters (the 后三半 dispatch's third half:
// "quota edits enforced ranges"). The schema has carried the three leaves since
// day one - llm.providers.<name>.models.<id>.quota_daily_micro /
// quota_monthly_micro and llm.providers.<name>.plan_credit_total_micro
// (schema.go:389-390/:414) - and the C23 pre-dispatch gate consumes them
// (internal/observe/quota_44.go), but no settings write owned them, so the
// model page could not edit a cap it could only read about in the file.
//
// THE RANGE, AND WHERE IT COMES FROM (⛔ no invented number). The dispatch
// names the ruler: grep the existing quantities first, then enforce THEM. What
// the codebase already answers about these fields:
//
//   - the unit and the lower bound: int64 micro-USD, 0 = uncapped (schema.go's
//     own field comment) / 0 = not entered and gates nothing
//     (internal/observe/quota_44.go's QuotaProvider reading). A negative spend
//     cap is not a cap: it would sit below every consumption forever, which is
//     not a number the gate has a meaning for - so the negative half of the
//     band is refused BY NAME, exactly the sentence shape the sibling money
//     leaves already use ("a negative price is not a price",
//     settings.go:184; "is not a token count (0 means unknown)", :146).
//   - the upper bound: the int64 type's own maximum. No narrower ceiling
//     exists anywhere in this tree (the C23 gate judges ratios of consumption
//     against the cap, nothing clamps it), and writing one here would be the
//     invented number the dispatch forbids. A value past int64 cannot even
//     parse, so the caller's ParseInt is the boundary's enforcer.
//
// SAME SHAPE AS EVERY SETTER BEFORE THEM. One owned key path each,
// writeOneKey, validate() on the base the write would produce before a byte is
// written, and every refusal leaves config.toml byte-identical. The model
// leaves require the catalog entry (requireCatalogEntry, reason 2 names the
// row that must exist first); the provider leaf requires the row
// (unknownProviderErr).

import "fmt"

// SetModelQuotaDaily writes
// llm.providers.<provider>.models.<model>.quota_daily_micro. Negative is
// refused by name before the file is touched; 0 is a legal write (0 =
// uncapped, the schema's own meaning, and clearing a cap is a real edit).
func (m *Manager) SetModelQuotaDaily(provider, model string, micros int64) ([]string, error) {
	if err := refuseNegativeQuota(provider, model, "quota_daily_micro", micros); err != nil {
		return nil, err
	}
	return m.setModelQuotaLeaf(provider, model, "quota_daily_micro", micros)
}

// SetModelQuotaMonthly writes quota_monthly_micro; the rules are
// SetModelQuotaDaily's (identical band, identical refusals).
func (m *Manager) SetModelQuotaMonthly(provider, model string, micros int64) ([]string, error) {
	if err := refuseNegativeQuota(provider, model, "quota_monthly_micro", micros); err != nil {
		return nil, err
	}
	return m.setModelQuotaLeaf(provider, model, "quota_monthly_micro", micros)
}

func (m *Manager) setModelQuotaLeaf(provider, model, leaf string, micros int64) ([]string, error) {
	return m.writeOneKey(keyPathProviderModelField(provider, model, leaf),
		func(base *Config) bool {
			p, ok := base.LLM.Providers[provider]
			if !ok {
				return false
			}
			spec, ok := p.Models[model]
			if !ok {
				return false
			}
			switch leaf {
			case "quota_daily_micro":
				spec.QuotaDailyMicro = micros
			case "quota_monthly_micro":
				spec.QuotaMonthlyMicro = micros
			default:
				return false
			}
			p.Models[model] = spec
			base.LLM.Providers[provider] = p
			return true
		}, func(base *Config) error {
			return requireCatalogEntry(base, provider, model)
		})
}

// refuseNegativeQuota is the one out-of-range refusal the quota band has: a
// negative spend cap is not a cap (0 means uncapped). The sentence mirrors the
// sibling money leaves' shape ("a negative price is not a price") and names the
// field and the value, so the page knows which leaf got the bad number.
func refuseNegativeQuota(provider, model, leaf string, micros int64) error {
	if micros >= 0 {
		return nil
	}
	return fmt.Errorf(
		"config: llm.providers.%s.models.%s.%s: %d 超出范围（额度不能为负，0 = 不封顶）；未落盘",
		provider, model, leaf, micros)
}

// SetProviderPlanCreditTotal writes
// llm.providers.<provider>.plan_credit_total_micro (the plan-billing credit
// total, manually entered per the schema; the C23 balance is bookkeeping
// against it). Same band as the per-model caps: negative refused by name,
// 0 = not entered is a legal answered value that gates nothing (quota_44.go's
// own reading).
func (m *Manager) SetProviderPlanCreditTotal(provider string, micros int64) ([]string, error) {
	if micros < 0 {
		return nil, fmt.Errorf(
			"config: llm.providers.%s.plan_credit_total_micro: %d 超出范围（credit 总额不能为负，0 = 未录入）；未落盘",
			provider, micros)
	}
	return m.writeOneKey(keyPathProviderField(provider, "plan_credit_total_micro"),
		func(base *Config) bool {
			p, ok := base.LLM.Providers[provider]
			if !ok {
				return false
			}
			p.PlanCreditTotalMicro = micros
			base.LLM.Providers[provider] = p
			return true
		}, func(base *Config) error {
			if _, ok := base.LLM.Providers[provider]; !ok {
				return unknownProviderErr(provider)
			}
			return nil
		})
}
