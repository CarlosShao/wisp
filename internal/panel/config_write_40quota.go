package panel

// Ticket 40 AC#3's quota-editing roster (the 后三半 dispatch's third half).
//
// Same two rules every roster since 313 follows, and the same reason this is a
// THIRD file and not more rows in an older one: the frozen ticket-248 pin
// (config_route_248_test.go) keeps meaning what it meant, each ticket's
// additions live in their own roster, and the write gate asks every roster
// through the shared selector columns in config_write_313.go.
//
// WHAT EACH FIELD OWNS.
//
//	model_quota_daily_micro  llm.providers.<p>.models.<m>.quota_daily_micro,
//	  the C23 per-model daily spend cap. Value = the micro-USD integer as a
//	  decimal string; internal/config's SetModelQuotaDaily owns the band
//	  (negative refused by name, 0 = uncapped) and the catalog-entry ruler.
//	model_quota_monthly_micro  the monthly twin, identical band and ruler.
//	provider_plan_credit_total_micro  llm.providers.<p>.
//	  plan_credit_total_micro, the plan-billing credit total. Provider
//	  selector, no model; the band lives in the same setter family.
//
// No field name here reaches a locked family and none carries a word of the
// approval vocabulary the frozen grant-boundary instrument keys on - the same
// two asserts config_313_test.go runs over the 313 roster run over this one
// through the shared accessor (the names are quota vocabulary, not grant
// vocabulary). The tier is [llm]'s restart stamp, exactly like every other
// catalog leaf (the model path is built once at assembly; panel_config_store.go
// stamps restart before anything field-specific re-adjudicates).

import "sort"

const (
	// FieldModelQuotaDaily is llm.providers.<p>.models.<m>.quota_daily_micro
	// (ticket 40's quota-editing half). Value = micro-USD integer string;
	// 0 = uncapped.
	FieldModelQuotaDaily = "model_quota_daily_micro"
	// FieldModelQuotaMonthly is the monthly twin of FieldModelQuotaDaily.
	FieldModelQuotaMonthly = "model_quota_monthly_micro"
	// FieldProviderPlanCredit is llm.providers.<p>.plan_credit_total_micro.
	// Provider selector only; value = micro-USD integer string; 0 = not
	// entered (an answered value that gates nothing, quota_44.go's reading).
	FieldProviderPlanCredit = "provider_plan_credit_total_micro"
)

// fields40Quota is this file's roster, sorted at the accessor. Kept apart from
// the older rosters on purpose (see the header).
var fields40Quota = []string{
	FieldModelQuotaDaily,
	FieldModelQuotaMonthly,
	FieldProviderPlanCredit,
}

// The selector columns of this roster, same shape as every roster's maps.
var (
	fields40QuotaNeedsProvider = map[string]bool{
		FieldModelQuotaDaily:    true,
		FieldModelQuotaMonthly:  true,
		FieldProviderPlanCredit: true,
	}
	fields40QuotaNeedsModel = map[string]bool{
		FieldModelQuotaDaily:   true,
		FieldModelQuotaMonthly: true,
	}
	fields40QuotaTakesValue = map[string]bool{
		FieldModelQuotaDaily:    true,
		FieldModelQuotaMonthly:  true,
		FieldProviderPlanCredit: true,
	}
)

// WritableFields40Quota returns ticket 40's quota roster, sorted. It is the
// accessor the union and the tests read, the same pattern every roster keeps.
func WritableFields40Quota() []string {
	out := append([]string(nil), fields40Quota...)
	sort.Strings(out)
	return out
}

// writableField40Quota asks this file's roster, not a copy of it - the term
// knownWritableField's union consults (config_handlers.go).
func writableField40Quota(field string) bool {
	for _, f := range fields40Quota {
		if f == field {
			return true
		}
	}
	return false
}
