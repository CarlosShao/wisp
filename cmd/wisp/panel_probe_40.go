package main

// Ticket 40 AC#3 (first two halves) - the model page's probe leg.
//
// WHAT THIS IS. The page's 测试 button used to be an honest stub sentence
// ("暂不支持连接测试"); this leg replaces the stub with the REAL measurement:
// the same machinery `wisp providers probe` runs (llm.NewResolver over the
// live config, BuildEndpointProvider, llm.RunProbeSuite) with the verdicts
// recorded into provider_health through ticket 09's sink (storeHealthSink).
// The receipt carries the measured facts verbatim - p50 latency, one verdict
// word per attempted capability, every 「声明可用，实测不可用」mismatch's frozen
// label and its error原文 - and nothing that claims a save: a probe writes no
// config key, so unlike the CRUD writes beside it there is no tier sentence
// here at all.
//
// ONE MACHINERY, NO SECOND COPY. providers.go's providersProbe already owns
// the CLI shape of this run; this handler is the panel shape over the same
// pieces (same resolver, same provider build, same sink, same suite). A probe
// that "measured" without those pieces would be the 空转返回OK shape the
// ticket forbids, and the wiring nail in panel_modelmgmt_40_test.go reads this
// file's bytes to keep them on the wire.
//
// THE STORE IS OPENED PER RUN (the cost/privacy handler shape, sessions.go's
// outlet): a chain that never sees a probe click opens no database, and the
// sink write has the same lifetime as the measurement that produced it.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/secret"
)

// modelProbeHandler40 is the probe leg: the three objects a real run needs,
// plus the audit sink its refusals go to.
type modelProbeHandler40 struct {
	mgr     *config.Manager
	secrets *secret.Store
	dir     string
	auditf  panel.AuditFunc
	actor   string
}

// newModelProbeHandler40 wires the leg over the SAME Manager and secret store
// the settings leg and the models-list leg use, so the probe measures exactly
// the row the page shows and the key the write path stored - one truth per
// process, no second loader.
func newModelProbeHandler40(dir string, mgr *config.Manager, secrets *secret.Store, auditf panel.AuditFunc, actor string) *modelProbeHandler40 {
	return &modelProbeHandler40{mgr: mgr, secrets: secrets, dir: dir, auditf: auditf, actor: actor}
}

// probeSuiteBudget is the whole-suite budget, the CLI's own shape (providers
// probe --timeout 60s): a record-producing budget, never a wall-clock deadline
// trick (D42#9 - context carries it, the suite's probes honour it).
const probeSuiteBudget40 = 60 * time.Second

// HandleModelProbe runs the real probe suite for one (provider, model) pair
// and answers with the receipt + payload line internal/panel's renderer shapes.
// Every refusal names what was missing; a run that cannot complete reports the
// real error text - never a bare OK.
func (h *modelProbeHandler40) HandleModelProbe(ctx context.Context, req panel.ComposerRequest) (string, error) {
	if h == nil || h.mgr == nil {
		err := fmt.Errorf("panel: 探测腿未接入（requestId=%q），未发生任何测量", req.RequestID)
		h.record(req, err)
		return "", err
	}
	provider := strings.TrimSpace(req.ConfigProvider)
	model := strings.TrimSpace(req.ConfigModel)
	if provider == "" || model == "" {
		err := fmt.Errorf("panel: 探测必须指名服务商与模型（configProvider/configModel），收到 provider=%q model=%q（requestId=%q）",
			provider, model, req.RequestID)
		h.record(req, err)
		return "", err
	}
	cfg := h.mgr.Config()
	if cfg == nil {
		err := fmt.Errorf("panel: 配置不可读，探测没有目录可查（requestId=%q）", req.RequestID)
		h.record(req, err)
		return "", err
	}
	prov, ok := cfg.LLM.Providers[provider]
	if !ok {
		err := fmt.Errorf("panel: 目录中没有服务商 %q，未发生任何测量（requestId=%q）", provider, req.RequestID)
		h.record(req, err)
		return "", err
	}
	spec, ok := prov.Models[model]
	if !ok {
		err := fmt.Errorf("panel: %s 的目录中没有模型 %q，未发生任何测量（requestId=%q）", provider, model, req.RequestID)
		h.record(req, err)
		return "", err
	}
	// The endpoint resolves the way the text path does: the pair under test as
	// a one-element chain, so the key comes out of the store and an
	// unresolvable ref fails loudly instead of probing unauthenticated (the
	// providers probe CLI's own Unconfigured rule).
	res := llm.NewResolver(cfg, h.secrets)
	res.TextChain = []string{provider + "/" + model}
	eps, err := res.ResolveChain()
	if err != nil {
		// The structural refusals ride the refusal sentence (Handle's err
		// path): the error原文 reaches the page verbatim through
		// RefusedEnvelopeForUser, the same shape every settings refusal
		// takes. No payload line, no verdicts, nothing recorded.
		h.record(req, err)
		return "", err
	}
	endpointProvider, err := llm.BuildEndpointProvider(eps[0], chainOptions(cfg))
	if err != nil {
		h.record(req, err)
		return "", err
	}
	store, err := memory.Open(h.dir)
	if err != nil {
		err = fmt.Errorf("存储不可用，判定无法写入 provider_health：%v", err)
		h.record(req, err)
		return "", err
	}
	defer func() { _ = store.Close() }()

	runCtx, cancel := context.WithTimeout(ctx, probeSuiteBudget40)
	defer cancel()
	rep, runErr := llm.RunProbeSuite(runCtx, endpointProvider, llm.ProbeSuiteOptions{
		Provider: provider,
		Model:    model,
		Declared: spec.Capabilities,
		Sink:     storeHealthSink{store: store},
		// SPEC-05 §3.1: the contradiction must be loud. The CLI announces it on
		// stderr; the panel shape announces it on the receipt itself (the
		// renderer prints every mismatch's frozen label + detail原文), so the
		// callback stays silent here and the report's own list is the source.
		Mismatch: func(llm.ProbeMismatch) {},
		Now:      func() time.Time { return time.Now().UTC() },
	})
	if runErr != nil {
		// A run the suite could not complete is an ANSWER, not a route
		// refusal: the receipt + payload carry the error原文 (the payload's
		// error field) so the page shows what failed instead of a bare
		// refusal sentence. Nothing was measured, nothing was recorded.
		h.record(req, runErr)
		out := panel.ModelProbeResult{
			Provider:   provider,
			Model:      model,
			Outcomes:   []panel.ModelProbeOutcome{},
			Mismatches: []panel.ModelProbeMismatch{},
			Error:      runErr.Error(),
		}
		return panel.RenderModelProbeReceipt(out) + "\n" + panel.RenderModelProbePayload(out), nil
	}
	out := panel.ModelProbeResult{
		Provider:   rep.Provider,
		Model:      rep.Model,
		LatencyMS:  rep.LatencyMS,
		AllOK:      rep.AllOK,
		Written:    rep.Written,
		Outcomes:   make([]panel.ModelProbeOutcome, 0, len(rep.Outcomes)),
		Mismatches: make([]panel.ModelProbeMismatch, 0, len(rep.Mismatches)),
	}
	for _, o := range rep.Outcomes {
		out.Outcomes = append(out.Outcomes, panel.ModelProbeOutcome{
			Capability: string(o.Capability),
			Declared:   o.Declared,
			OK:         o.Result.OK,
			Attempted:  o.Attempted,
			Detail:     o.Result.Detail,
		})
	}
	for _, m := range rep.Mismatches {
		out.Mismatches = append(out.Mismatches, panel.ModelProbeMismatch{
			Capability: string(m.Capability),
			Detail:     m.Detail,
			Label:      m.Label(),
		})
	}
	h.record(req, nil)
	return panel.RenderModelProbeReceipt(out) + "\n" + panel.RenderModelProbePayload(out), nil
}

// record writes one audit line per outcome. Names the method, request id and
// outcome - never a provider error body (those live on the receipt the page
// shows, and the audit surface stays value-free like every other route's).
func (h *modelProbeHandler40) record(req panel.ComposerRequest, err error) {
	if h == nil || h.auditf == nil {
		return
	}
	outcome := "ok"
	if err != nil {
		outcome = "refused"
	}
	actor := h.actor
	if actor == "" {
		actor = "unknown"
	}
	h.auditf("panel: MODEL-PROBE request=%q provider=%q model=%q actor=%q outcome=%q",
		req.RequestID, req.ConfigProvider, req.ConfigModel, actor, outcome)
}
