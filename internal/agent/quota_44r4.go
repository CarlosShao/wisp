package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44 AC#5 wiring (leg 44-r4): the pre-dispatch quota gate.
//
// The math lives in observe/quota_44.go (44-r2): per-model day/month caps and
// the per-provider plan-credit line, both judged PRE-DISPATCH (票 44:33,
// SPEC-03 §3.1 防负＝派发前检查). This file is the one missing inch: the
// gate plugged in front of an actual dispatch, so a model whose remaining hit
// 0 is SKIPPED (zero bytes leave the machine for it) while the text_chain
// keeps advancing instead of breaking the chain — 票 44:34-35 "over-quota
// model is SKIPPED via the text_chain (never dispatched into debt)" and
// :39-40 "quota exhaustion ... funnels through the same 'advance text_chain'
// mechanism (ticket 11 owns stepping; C23 owns the quota signal)".
//
// The signal C23 owns is a classified observe.ClassBudget error returned
// before the element's Stream is ever called. That is exactly the shape
// ticket 11's llm.ChainRunner steps over — its own tests feed a ClassBudget
// element into the chain and prove the walk moves to the next entry
// (internal/llm/llm_test.go, chainProvider "quota" case) — so wrapping each
// chain element needs no change to the stepping code and no loop change.
//
// The disk-side bookkeeping stays inside the D35-frozen schema: remaining
// detail is session-scoped in the ledger (44-r2's persistence note), while
// the permanent per-model quota signal is provider_health.quota_state
// (SPEC-02 §3), written through the existing memory DAO method — never a new
// column. plan vs pay-per-token consumption remain separate buckets by
// construction (observe.QuotaSpend), which is AC#5's "tracked separately".

// QuotaSignalWriter is the assertion-discovered disk mouth
// (memory.Store.RecordProviderError, already implemented and tested there).
// Nil means "no store wired": the gate still gates, it just says nothing to
// disk.
type QuotaSignalWriter interface {
	RecordProviderError(ctx context.Context, provider, model, class, detail, quotaState string) error
}

// QuotaGate holds the quota layer of the C23 pre-dispatch check.
type QuotaGate struct {
	params    observe.QuotaParams
	cards     map[string]observe.QuotaModel
	providers map[string]observe.QuotaProvider
	ledger    *observe.QuotaLedger
	// notify is the visible-event mouth (票 44:35 "visible fallback event",
	// SPEC-03 §3.1 不得静默): every warn and every skip is announced before
	// the gate acts on it.
	notify func(Event)
	// signals is the disk-side quota_state writer (may be nil).
	signals QuotaSignalWriter
	// logW receives signal failures; they never block a dispatch.
	logW func(string, ...any)
}

// NewQuotaGate wires cards/providers from config-shape inputs plus the
// session ledger and the event mouth. A nil ledger starts a fresh one on the
// ledger's own clock (fixture clocks pass their own, AC#2 style).
func NewQuotaGate(params observe.QuotaParams, cards map[string]observe.QuotaModel,
	providers map[string]observe.QuotaProvider, ledger *observe.QuotaLedger,
	notify func(Event),
) *QuotaGate {
	if ledger == nil {
		ledger = observe.NewQuotaLedger(nil)
	}
	if cards == nil {
		cards = map[string]observe.QuotaModel{}
	}
	if providers == nil {
		providers = map[string]observe.QuotaProvider{}
	}
	return &QuotaGate{
		params: params, cards: cards, providers: providers,
		ledger: ledger, notify: notify,
	}
}

// SetSignalWriter points the exhausted-state disk write at the store
// (internal/memory's provider_health DAO).
func (g *QuotaGate) SetSignalWriter(w QuotaSignalWriter, logf func(string, ...any)) {
	g.signals = w
	g.logW = logf
}

// Ledger exposes the session bookkeeping so the host (or the loop's finish
// path) can post each dispatch's spend into the same buckets the gate reads.
func (g *QuotaGate) Ledger() *observe.QuotaLedger { return g.ledger }

// QuotaWiring adapts the live config into gate inputs (the same reading
// discipline as CostWiring: nothing is guessed). A model card without its own
// billing inherits the provider's; a missing provider billing defaults to
// pay-per-token, which is the schema default. The alert line comes from the
// single [cost] alert_threshold through observe, so the global budget and the
// quota layers cannot drift two different 80%s apart (44-r2's rule).
func QuotaWiring(cfg *config.Config) (map[string]observe.QuotaModel, map[string]observe.QuotaProvider, observe.QuotaParams) {
	if cfg == nil {
		return nil, nil, observe.QuotaParams{}
	}
	cards := map[string]observe.QuotaModel{}
	providers := map[string]observe.QuotaProvider{}
	for pname, p := range cfg.LLM.Providers {
		pBilling := p.Billing
		if pBilling == "" {
			pBilling = observe.BillingPayPerToken
		}
		providers[pname] = observe.QuotaProvider{Provider: pname, PlanCreditTotalMicro: p.PlanCreditTotalMicro}
		for mname, m := range p.Models {
			billing := m.Billing
			if billing == "" {
				billing = pBilling
			}
			cards[pname+"/"+mname] = observe.QuotaModel{
				Provider:      pname,
				Model:         mname,
				Billing:       billing,
				DailyMicros:   m.QuotaDailyMicro,
				MonthlyMicros: m.QuotaMonthlyMicro,
			}
		}
	}
	var params observe.QuotaParams
	if b, err := observe.BudgetFrom(cfg.Cost.DailyBudget, cfg.Cost.MonthlyBudget,
		cfg.Cost.AlertThreshold, cfg.Cost.OverBudget); err == nil {
		params = observe.NewQuotaParamsFrom(b)
	}
	return cards, providers, params
}

// Decide is the pre-dispatch judgment for one "provider/model" chain entry.
// An entry nobody gave a quota card is no-card = uncapped and allowed
// (44-r2's rule: never block real work with an invented cap).
func (g *QuotaGate) Decide(key string) observe.QuotaVerdict {
	m, ok := g.cards[key]
	if !ok {
		provider, model := splitQuotaKey(key)
		m = observe.QuotaModel{Provider: provider, Model: model}
	}
	var prov *observe.QuotaProvider
	if q, ok := g.providers[m.Provider]; ok {
		prov = &q
	}
	return g.params.JudgeModel(m, prov, g.ledger.Spend())
}

// GateProvider wraps one text_chain element with the pre-dispatch check.
// The host builds the chain (llm.BuildChain / ChainRunner, ticket 11's
// surface) and wraps each element with the element's "provider/model" name;
// an over-quota element is refused before its Stream runs, so it is skipped
// and the chain advances. Warn verdicts dispatch anyway, with the event out
// first — same 80/100 semantics as the [cost] gate.
func (g *QuotaGate) GateProvider(p llm.LlmProvider, key string) llm.LlmProvider {
	return &gatedChainElement{inner: p, key: key, gate: g}
}

type gatedChainElement struct {
	inner llm.LlmProvider
	key   string
	gate  *QuotaGate
}

func (e *gatedChainElement) Info() llm.ProviderInfo { return e.inner.Info() }

func (e *gatedChainElement) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	v := e.gate.Decide(e.key)
	if v.Level == observe.LevelWarn && e.gate.notify != nil {
		e.gate.notify(Event{Kind: EvWarn, Text: v.Notice})
	}
	if v.AllowDispatch {
		return e.inner.Stream(ctx, req, emit)
	}
	// Remaining hit 0 on a configured cap: never dispatch into debt.
	if e.gate.notify != nil {
		e.gate.notify(Event{Kind: EvWarn, Text: v.Notice})
	}
	e.gate.writeExhaustedSignal(ctx, e.key, v.Notice)
	// The classified error is what ticket 11's ChainRunner steps on; the
	// inner provider was not called, so zero bytes left the machine for
	// this model while the chain keeps going.
	return observe.New(observe.ClassBudget, fmt.Sprintf("agent: 配额派发前闸跳过 %s：%s", e.key, v.Notice))
}

// writeExhaustedSignal persists the disk-side reading through the existing
// provider_health DAO (quota_state = 'exhausted'). A failure is logged, never
// fatal: the gate already refused the dispatch; the disk note is forensics.
func (g *QuotaGate) writeExhaustedSignal(ctx context.Context, key, notice string) {
	if g.signals == nil {
		return
	}
	provider, model := splitQuotaKey(key)
	if err := g.signals.RecordProviderError(ctx, provider, model, "budget", notice, "exhausted"); err != nil {
		if g.logW != nil {
			g.logW("agent: quota exhausted signal not written", "err", err, "model", key)
		}
	}
}

// Book posts one dispatch's spend (integer micros) into the session ledger:
// the model's day/month buckets and, when the card bills as a plan, the
// provider's credit bucket — the plan-credit 分账记账 leg of AC#5, separate
// from pay-per-token by construction. Negative micros are refused by the
// ledger (anti-negative bookkeeping), zero books nothing.
func (g *QuotaGate) Book(key string, micros int64) error {
	m, ok := g.cards[key]
	if !ok {
		provider, model := splitQuotaKey(key)
		m = observe.QuotaModel{Provider: provider, Model: model}
	}
	return g.ledger.Post(m, micros)
}

// Chain reads the quota layer across an ordered text_chain: the first entry
// still dispatchable plus the skipped ones in order, each with its
// user-visible reason. This is the query half of "chain advanced, not
// broken": a fully-exhausted chain returns chosen == "" with every entry in
// skips, and it is ticket 11's runner that then ends the task — never this
// gate.
func (g *QuotaGate) Chain(chain []string) (chosen string, skips []observe.QuotaSkip) {
	sp := g.ledger.Spend()
	_, skips, _ = g.params.SelectChain(chain, g.cards, g.providers, sp)
	return g.pickFirstAllowed(chain, skips), skips
}

func (g *QuotaGate) pickFirstAllowed(chain []string, skips []observe.QuotaSkip) string {
	skipped := make(map[string]bool, len(skips))
	for _, s := range skips {
		skipped[s.Model] = true
	}
	for _, k := range chain {
		if !skipped[k] {
			return k
		}
	}
	return ""
}

func splitQuotaKey(key string) (provider, model string) {
	if i := strings.IndexByte(key, '/'); i >= 0 {
		return key[:i], key[i+1:]
	}
	return key, ""
}
