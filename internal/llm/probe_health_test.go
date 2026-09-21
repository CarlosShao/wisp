package llm_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
)

// AC#6: the probe suite runs REAL requests against a live mockllm whose
// capabilities are set by /__control/capability, and the verdicts land in
// provider_health.
//
// The trap this file is written against: a "probe" that answers from the
// catalog declaration would be green forever and tell the user nothing. Every
// case below therefore pins three things at once:
//
//  1. the SERVER is provably in the mode (read back from /__control/state,
//     not assumed),
//  2. requests really went out (the server's own per-route counter),
//  3. the verdict is the OPPOSITE of the declaration in the broken cases, and
//     the persisted row says so.

// declaredEverything is the lying catalog: it claims every capability this
// suite probes.
var declaredEverything = config.Capabilities{
	Text: true, Vision: true, Thinking: true, AudioIn: true, FC: true, Stream: true,
}

// probeFixture is one live mockllm plus the SQLite store the verdicts land in.
type probeFixture struct {
	srv   *adaptertest.Mockllm
	store *memory.Store
	sink  llm.HealthSink
}

func newProbeFixture(t *testing.T) *probeFixture {
	t.Helper()
	pf := &probeFixture{srv: adaptertest.StartMockllm(t)}
	store, err := memory.Open(filepath.Join(t.TempDir(), "data"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	pf.store = store
	pf.sink = memoryHealthSink{store: store}
	return pf
}

// memoryHealthSink adapts the C5 write seam to the provider_health DAO. This is
// the exact mapping the composition root performs; it lives here because
// internal/llm must not import the storage layer.
type memoryHealthSink struct{ store *memory.Store }

func (m memoryHealthSink) RecordProbe(ctx context.Context, provider, model string,
	flags llm.MeasuredFlags, ok bool, latencyMS int64, at time.Time,
) error {
	return m.store.UpsertProviderProbe(ctx, provider, model, memory.ProbeFlags{
		Text: flags.Text, Vision: flags.Vision, AudioIn: flags.AudioIn,
		AudioOut: flags.AudioOut, Thinking: flags.Thinking, FC: flags.FC,
	}, ok, latencyMS, at)
}

// setCapability puts the live server into a capability mode and PROVES it
// (state read-back), after a reset so the request counters start clean.
func (pf *probeFixture) setCapability(t *testing.T, body string) {
	t.Helper()
	pf.srv.Reset(t)
	pf.srv.Control(t, "/__control/capability", body)
	resp, err := http.Get(pf.srv.Base + "/__control/state")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var st struct {
		Capabilities map[string]string `json:"capabilities"`
	}
	if err := json.Unmarshal(raw, &st); err != nil {
		t.Fatalf("decode state %s: %v", raw, err)
	}
	want := map[string]string{}
	if err := json.Unmarshal([]byte(body), &want); err != nil {
		t.Fatal(err)
	}
	for k, v := range want {
		if st.Capabilities[k] != v {
			t.Fatalf("server capability %s = %q, want %q (the mock is not in the mode the probe is being tested against: %+v)",
				k, st.Capabilities[k], v, st.Capabilities)
		}
	}
}

// requests counts what the server actually served on one route.
func (pf *probeFixture) requests(t *testing.T, route string) int {
	t.Helper()
	return pf.srv.RouteCount(t, route)
}

// runProbe is the suite under test; it returns the report plus the mismatches
// the notice saw. The sink and notice are injected here so no case can forget
// them, and the provider is built fresh per protocol.
func runProbe(t *testing.T, pf *probeFixture, protocol string, o llm.ProbeSuiteOptions) (llm.ProbeReport, []llm.ProbeMismatch) {
	t.Helper()
	var seen []llm.ProbeMismatch
	if o.Provider == "" {
		o.Provider = "mock"
	}
	if o.Model == "" {
		o.Model = "mock-small"
	}
	o.Sink = pf.sink
	o.Mismatch = func(m llm.ProbeMismatch) { seen = append(seen, m) }
	rep, err := llm.RunProbeSuite(context.Background(),
		mustProvider(t, protocol, pf.srv.Base+"/v1"), o)
	if err != nil {
		t.Fatalf("RunProbeSuite: %v", err)
	}
	return rep, seen
}

// mustFlag fails when a capability was left unmeasured: nil (never probed) and
// measured-false are different facts and the row must say which one it holds.
func mustFlag(t *testing.T, f memory.ProbeFlags, cap string) bool {
	t.Helper()
	v := f.Get(cap)
	if v == nil {
		t.Fatalf("provider_health %s = nil (never probed), want a verdict", cap)
	}
	return *v
}

// TestProbeSuiteMeasuresBrokenFC: the catalog claims fc AND vision; the wire
// can only do vision. The asymmetry is the point - a probe answering from
// config would mark both capable, and a probe answering "always broken" would
// mark both broken.
func TestProbeSuiteMeasuresBrokenFC(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{"fc":"broken"}`)

	rep, mismatches := runProbe(t, pf, "openai-chat", llm.ProbeSuiteOptions{Declared: declaredEverything})

	// 2. the requests went out: fc + vision + thinking on the chat route.
	if n := pf.requests(t, "chat"); n != 3 {
		t.Errorf("mockllm served %d chat requests, want 3 (one probe per capability)", n)
	}

	fc, ok := rep.Verdict(llm.ProbeFC)
	if !ok {
		t.Fatal("no fc verdict")
	}
	if fc.Result.OK {
		t.Errorf("fc probe reported capable against fc=broken; detail=%s", fc.Result.Detail)
	}
	vision, ok := rep.Verdict(llm.ProbeVision)
	if !ok {
		t.Fatal("no vision verdict")
	}
	if !vision.Result.OK {
		t.Errorf("vision probe reported broken although the server is vision-capable: %s", vision.Result.Detail)
	}

	// 1. declared-vs-measured mismatch, fc only, and it carries the copy.
	if len(mismatches) != 1 || mismatches[0].Capability != llm.ProbeFC {
		t.Fatalf("mismatches = %+v, want exactly one for fc", mismatches)
	}
	if mismatches[0].Kind() != "capability_mismatch" || !strings.Contains(mismatches[0].Label(), "fc") {
		t.Errorf("mismatch = %+v label=%q", mismatches[0], mismatches[0].Label())
	}
	if len(rep.Mismatches) != 1 {
		t.Errorf("report mismatches = %+v", rep.Mismatches)
	}

	// 3. the persisted row, read back from SQLite.
	row, err := pf.store.ProviderHealthRow(context.Background(), "mock", "mock-small")
	if err != nil {
		t.Fatalf("provider_health row: %v", err)
	}
	if mustFlag(t, row.Probe, "fc") {
		t.Error("provider_health says fc works; the server said it does not")
	}
	if !mustFlag(t, row.Probe, "vision") {
		t.Error("provider_health says vision is broken; the probe never saw that")
	}
	if row.QuotaState != memory.QuotaUnknown || row.LastProbeAt.IsZero() {
		t.Errorf("row metadata = %+v", row)
	}
	if mustFlag(t, row.Probe, "thinking") != true {
		t.Error("thinking verdict missing: the text turn proves the probe is not blanket-reporting broken")
	}
	if row.Probe.Get("audio_in") != nil {
		t.Error("audio_in must stay nil: no audio request exists to probe yet (C7), so it was never measured")
	}
}

// TestProbeSuiteMeasuresBrokenVision is the mirror image: the same code path,
// the opposite capability broken. A suite that hardcodes either answer cannot
// pass both this and TestProbeSuiteMeasuresBrokenFC.
func TestProbeSuiteMeasuresBrokenVision(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{"vision":"broken"}`)

	rep, mismatches := runProbe(t, pf, "openai-chat", llm.ProbeSuiteOptions{Declared: declaredEverything})

	if n := pf.requests(t, "chat"); n != 3 {
		t.Errorf("mockllm served %d chat requests, want 3", n)
	}
	fc, _ := rep.Verdict(llm.ProbeFC)
	if !fc.Result.OK {
		t.Errorf("fc probe reported broken although the server is fc-capable: %s", fc.Result.Detail)
	}
	vision, _ := rep.Verdict(llm.ProbeVision)
	if vision.Result.OK {
		t.Error("vision probe accepted the 400 a non-vision model returns as success")
	}
	if len(mismatches) != 1 || mismatches[0].Capability != llm.ProbeVision {
		t.Fatalf("mismatches = %+v, want exactly one for vision", mismatches)
	}
	if !strings.Contains(mismatches[0].Detail, "vision=broken") {
		t.Errorf("mismatch detail = %q, want the provider's own refusal", mismatches[0].Detail)
	}
	row, err := pf.store.ProviderHealthRow(context.Background(), "mock", "mock-small")
	if err != nil {
		t.Fatal(err)
	}
	if mustFlag(t, row.Probe, "vision") {
		t.Error("provider_health says vision works; the server 400'd the image")
	}
	if !mustFlag(t, row.Probe, "fc") {
		t.Error("provider_health says fc is broken; only vision was")
	}
}

// TestProbeSuiteHonestProviderKeepsRowClean: with nothing broken the suite
// must record all-usable and emit no mismatch (the same fixtures, so a
// probe that always reports broken cannot pass).
func TestProbeSuiteHonestProviderRecordsNoMismatch(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{"fc":"capable","vision":"capable"}`)

	rep, mismatches := runProbe(t, pf, "openai-chat", llm.ProbeSuiteOptions{Declared: declaredEverything})

	if len(mismatches) != 0 || len(rep.Mismatches) != 0 {
		t.Errorf("mismatches = %+v, want none for an honest provider", mismatches)
	}
	if !rep.AllOK || !rep.Written {
		t.Errorf("report = %+v", rep)
	}
	row, err := pf.store.ProviderHealthRow(context.Background(), "mock", "mock-small")
	if err != nil {
		t.Fatal(err)
	}
	for _, cap := range []string{"fc", "vision", "thinking"} {
		if !mustFlag(t, row.Probe, cap) {
			t.Errorf("provider_health %s = false against a fully capable provider", cap)
		}
	}
	if row.QuotaState != memory.QuotaOK {
		t.Errorf("quota_state = %q, want ok", row.QuotaState)
	}
}

// TestProbeSuiteHonestNegativeIsNotAMismatch: a catalog that declares fc
// broken AND is measured broken is telling the truth - the mismatch event is
// about contradiction, not about failure. This is what keeps the event
// signal-bearing on the panel.
func TestProbeSuiteHonestNegativeIsNotAMismatch(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{"fc":"broken"}`)

	declared := declaredEverything
	declared.FC = false
	_, mismatches := runProbe(t, pf, "openai-chat", llm.ProbeSuiteOptions{
		Declared:     declared,
		Capabilities: []llm.ProbeCapability{llm.ProbeFC},
	})

	if len(mismatches) != 0 {
		t.Errorf("mismatches = %+v, want none: declared false + measured false is honest", mismatches)
	}
	row, err := pf.store.ProviderHealthRow(context.Background(), "mock", "mock-small")
	if err != nil {
		t.Fatal(err)
	}
	if mustFlag(t, row.Probe, "fc") {
		t.Error("the measured verdict must still be recorded as broken")
	}
}

// TestProbeSuiteBrokenOnAllThreeDialects proves the emulation AND the probe are
// protocol-wide: each adapter must see the same two honest failures in its own
// wire shape.
func TestProbeSuiteBrokenOnAllThreeDialects(t *testing.T) {
	for _, protocol := range []string{"openai-chat", "anthropic", "openai-responses"} {
		t.Run(protocol, func(t *testing.T) {
			pf := newProbeFixture(t)
			pf.setCapability(t, `{"fc":"broken","vision":"broken"}`)

			rep, mismatches := runProbe(t, pf, protocol, llm.ProbeSuiteOptions{
				Provider:     "mock-" + protocol,
				Declared:     declaredEverything,
				Capabilities: []llm.ProbeCapability{llm.ProbeFC, llm.ProbeVision},
			})

			route := map[string]string{
				"openai-chat": "chat", "anthropic": "messages",
				"openai-responses": "responses",
			}[protocol]
			if n := pf.requests(t, route); n != 2 {
				t.Errorf("%s: %s served %d requests, want 2 (both probes really ran)",
					protocol, route, n)
			}
			fc, _ := rep.Verdict(llm.ProbeFC)
			vision, _ := rep.Verdict(llm.ProbeVision)
			if fc.Result.OK {
				t.Errorf("%s: fc probe missed the broken tool choice", protocol)
			}
			if vision.Result.OK {
				t.Errorf("%s: vision probe accepted the image rejection as success", protocol)
			}
			if len(mismatches) != 2 {
				t.Errorf("%s: mismatches = %+v, want one per capability", protocol, mismatches)
			}
			row, err := pf.store.ProviderHealthRow(context.Background(), "mock-"+protocol, "mock-small")
			if err != nil {
				t.Fatalf("%s: row: %v", protocol, err)
			}
			if mustFlag(t, row.Probe, "fc") || mustFlag(t, row.Probe, "vision") {
				t.Errorf("%s: provider_health claims a capability the wire refused: %+v", protocol, row.Probe)
			}
		})
	}
}

// TestProbeSuiteCapableOnAllThreeDialects is the other half of the measurement:
// a probe that reported broken for a provider it merely failed to understand
// would be just as useless as one that echoes the config. Each dialect's own
// forced-tool-choice spelling must be recognised as a real tool call.
func TestProbeSuiteCapableOnAllThreeDialects(t *testing.T) {
	for _, protocol := range []string{"openai-chat", "anthropic", "openai-responses"} {
		t.Run(protocol, func(t *testing.T) {
			pf := newProbeFixture(t)
			pf.setCapability(t, `{"fc":"capable","vision":"capable"}`)

			rep, mismatches := runProbe(t, pf, protocol, llm.ProbeSuiteOptions{
				Provider:     "cap-" + protocol,
				Declared:     declaredEverything,
				Capabilities: []llm.ProbeCapability{llm.ProbeFC, llm.ProbeVision},
			})
			fc, _ := rep.Verdict(llm.ProbeFC)
			vision, _ := rep.Verdict(llm.ProbeVision)
			if !fc.Result.OK {
				t.Errorf("%s: fc probe reported broken against a capable provider: %s",
					protocol, fc.Result.Detail)
			}
			if !vision.Result.OK {
				t.Errorf("%s: vision probe reported broken against a capable provider: %s",
					protocol, vision.Result.Detail)
			}
			if len(mismatches) != 0 {
				t.Errorf("%s: mismatches = %+v, want none", protocol, mismatches)
			}
			row, err := pf.store.ProviderHealthRow(context.Background(), "cap-"+protocol, "mock-small")
			if err != nil {
				t.Fatal(err)
			}
			if !mustFlag(t, row.Probe, "fc") || !mustFlag(t, row.Probe, "vision") {
				t.Errorf("%s: provider_health = %+v, want both measured usable", protocol, row.Probe)
			}
		})
	}
}

// TestProbeSuiteRefusesSilentRuns pins the "never silent" rule structurally:
// without a sink or without a notice the suite refuses to run at all, and it
// sends no traffic while refusing.
func TestProbeSuiteRefusesSilentRuns(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{}`)
	probe := mustProvider(t, "openai-chat", pf.srv.Base+"/v1")

	if _, err := llm.RunProbeSuite(context.Background(), probe, llm.ProbeSuiteOptions{
		Provider: "mock", Model: "mock-small", Declared: declaredEverything,
		Mismatch: func(llm.ProbeMismatch) {},
	}); err == nil {
		t.Error("a run with no HealthSink must fail: an unread probe is not a measurement")
	} else if !strings.Contains(err.Error(), "HealthSink") {
		t.Errorf("err = %v", err)
	}
	if _, err := llm.RunProbeSuite(context.Background(), probe, llm.ProbeSuiteOptions{
		Provider: "mock", Model: "mock-small", Declared: declaredEverything,
		Sink: pf.sink,
	}); err == nil {
		t.Error("a run with no mismatch notice must fail (SPEC-05 sec 3.1)")
	} else if !strings.Contains(err.Error(), "mismatch") {
		t.Errorf("err = %v", err)
	}
	if n := pf.requests(t, "chat"); n != 0 {
		t.Errorf("rejected runs sent %d requests, want 0", n)
	}
}

// TestProbeSuiteSinkFailurePropagatesAndIsNotWritten: a provider_health write
// that fails must never be swallowed - the caller has to learn the verdict is
// unpersisted.
func TestProbeSuiteSinkFailurePropagates(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{}`)
	boom := errors.New("disk on fire")

	rep, err := llm.RunProbeSuite(context.Background(),
		mustProvider(t, "openai-chat", pf.srv.Base+"/v1"), llm.ProbeSuiteOptions{
			Provider: "mock", Model: "mock-small", Declared: declaredEverything,
			Capabilities: []llm.ProbeCapability{llm.ProbeFC},
			Sink:         failingSink{err: boom},
			Mismatch:     func(llm.ProbeMismatch) {},
		})
	if err == nil || !errors.Is(err, boom) {
		t.Fatalf("err = %v, want the sink's error to propagate", err)
	}
	if rep.Written {
		t.Error("report claims the verdict was written")
	}
	if !rep.Outcomes[0].Result.OK {
		t.Errorf("the measurement itself was lost with the write: %+v", rep.Outcomes[0])
	}
	if _, err := pf.store.ProviderHealthRow(context.Background(), "mock", "mock-small"); err == nil {
		t.Error("nothing should have reached provider_health")
	}
}

type failingSink struct{ err error }

func (f failingSink) RecordProbe(context.Context, string, string, llm.MeasuredFlags,
	bool, int64, time.Time,
) error {
	return f.err
}

// TestProbeSuiteAudioStaysUnprobed: attempting the deferred audio case must
// record NO verdict rather than a fabricated one, and must not drag the run
// verdict down.
func TestProbeSuiteAudioStaysUnprobed(t *testing.T) {
	pf := newProbeFixture(t)
	pf.setCapability(t, `{}`)

	rep, mismatches := runProbe(t, pf, "openai-chat", llm.ProbeSuiteOptions{
		Declared:     declaredEverything,
		Capabilities: []llm.ProbeCapability{llm.ProbeAudio, llm.ProbeFC},
	})
	audio, _ := rep.Verdict(llm.ProbeAudio)
	if audio.Attempted {
		t.Error("audio claims to have been probed")
	}
	if audio.Result.Detail == "" {
		t.Error("audio must carry the deferral reason")
	}
	if rep.Flags.Get("audio_in") != nil {
		t.Error("audio_in verdict recorded without a request")
	}
	for _, m := range mismatches {
		if m.Capability == llm.ProbeAudio {
			t.Error("an unattempted capability was reported as a mismatch")
		}
	}
	if !rep.AllOK {
		t.Error("the unprobed audio case must not decide the run verdict")
	}
	if n := pf.requests(t, "chat"); n != 1 {
		t.Errorf("served %d chat requests, want 1 (only fc went out)", n)
	}
}
