package main

// Ticket 40 AC#3 (first two halves) - the teeth the dispatch names:
//
//	1. CRUD round-trip: 模型清单 增/改/删 真落 config.toml（走 config 写门），
//	   读回逐字段相等 - and the reader is the DISK (a fresh LoadFile), so a
//	   write that only moved memory reddens the read-back, and the read-back
//	   face (config.get) is asserted off the SAME bytes, so a disk write the
//	   read face never picks up reddens too.
//	2. Probe button: config.model.probe runs the REAL llm probe machinery
//	   (RunProbeSuite -> provider_health) and echoes the measured facts
//	   verbatim - latency, per-capability verdicts, mismatch labels and the
//	   raw error text. A stub "OK" with nothing behind it has no wire shape it
//	   can wear (an empty Outcomes list with a non-empty error is the only
//	   failure shape; Written=false says nothing was recorded).
//	3. The [llm] endpoint-class CRUD writes answer the restart tier with the
//	   39 AC#1 verbatim sentence; the probe receipt carries NO tier sentence
//	   (it measures, it never saves - tier wording would be the wrong-tier lie).
//
// The probe cases run against the live mockllm (adaptertest), the same server
// the CLI's own probe tests use, so a verdict that moves with the server's
// capability mode is a measurement, not an echo.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// model40 is a made-up catalog id (the 313 rule: the repo never invents real
// model names).
const model40 = "model-under-test-40"

// writeModel40 drives one config.set through the production dispatch chain.
func writeModel40(t *testing.T, disp *panel.ComposerDispatch, id, field, provider, model, value string) string {
	t.Helper()
	raw, err := json.Marshal(map[string]string{
		"method": "config.set", "requestId": id, "source": "panel-composer", "configField": field,
		"configProvider": provider, "configModel": model, "configValue": value,
	})
	if err != nil {
		t.Fatal(err)
	}
	reply, err := disp.Handle(context.Background(), string(raw))
	if err != nil {
		t.Fatalf("config.set %s was refused: %v", field, err)
	}
	return reply
}

// diskConfig40 reads config.toml FRESH from disk. This is the round-trip's
// reader: nothing the leg holds in memory can influence it, so "只改内存不落盘"
// fails here before any other assertion could.
func diskConfig40(t *testing.T, dir string) *config.Config {
	t.Helper()
	cfg, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("fresh config load: %v", err)
	}
	return cfg
}

// wantFreshRow40 is the exact row CreateModelRow must land on disk: the
// schema's own defaults (settings_320.go's applyDefaults), every other leaf
// zero. 逐字段, not a DeepEqual shrug.
func wantFreshRow40() config.ModelSpec {
	return config.ModelSpec{Billing: "pay-per-token", Enabled: true}
}

// assertRowFields40 compares one ModelSpec field by field against want.
func assertRowFields40(t *testing.T, label string, got, want config.ModelSpec) {
	t.Helper()
	if got.Display != want.Display {
		t.Errorf("%s display = %q, want %q", label, got.Display, want.Display)
	}
	if got.ContextWindow != want.ContextWindow {
		t.Errorf("%s context_window = %d, want %d", label, got.ContextWindow, want.ContextWindow)
	}
	if got.MaxOutputTokens != want.MaxOutputTokens {
		t.Errorf("%s max_output_tokens = %d, want %d", label, got.MaxOutputTokens, want.MaxOutputTokens)
	}
	if !reflect.DeepEqual(got.ThinkingLevels, want.ThinkingLevels) {
		t.Errorf("%s thinking_levels = %v, want %v", label, got.ThinkingLevels, want.ThinkingLevels)
	}
	if got.Billing != want.Billing {
		t.Errorf("%s billing = %q, want %q", label, got.Billing, want.Billing)
	}
	if got.Price != want.Price {
		t.Errorf("%s price = %+v, want %+v", label, got.Price, want.Price)
	}
	if got.QuotaDailyMicro != want.QuotaDailyMicro || got.QuotaMonthlyMicro != want.QuotaMonthlyMicro {
		t.Errorf("%s quota = %d/%d, want %d/%d", label, got.QuotaDailyMicro, got.QuotaMonthlyMicro, want.QuotaDailyMicro, want.QuotaMonthlyMicro)
	}
	if got.Enabled != want.Enabled {
		t.Errorf("%s enabled = %v, want %v", label, got.Enabled, want.Enabled)
	}
	if got.Capabilities != want.Capabilities {
		t.Errorf("%s capabilities = %+v, want %+v", label, got.Capabilities, want.Capabilities)
	}
}

// payload40 parses the machine line after a marker out of one reply.
func payload40(t *testing.T, reply, marker string) map[string]any {
	t.Helper()
	idx := strings.Index(reply, marker)
	if idx < 0 {
		t.Fatalf("reply carries no %s line: %q", marker, reply)
	}
	line := strings.SplitN(reply[idx+len(marker):], "\n", 2)[0]
	var payload map[string]any
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("the payload after %s is not JSON: %v", marker, err)
	}
	return payload
}

// ---------------------------------------------------------------------------
// Half 1: CRUD round-trip.
// ---------------------------------------------------------------------------

func Test40ModelMgmtCRUDRoundTripFromDiskAndBack(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	// ADD. The receipt names the key path and the read face must answer it.
	writeModel40(t, disp, "r-40-add", panel.FieldModelRowCreate, "deepseek", model40, "")
	disk := diskConfig40(t, dir)
	row, ok := disk.LLM.Providers["deepseek"].Models[model40]
	if !ok {
		t.Fatalf("the created entry is not on disk:\n%s", body40(t, dir))
	}
	assertRowFields40(t, "created", row, wantFreshRow40())

	// MODIFY, twice, through two different leaf setters. After each write the
	// disk must carry exactly the new leaf and every other leaf must be where
	// it was (field isolation, read back from the DISK).
	writeModel40(t, disp, "r-40-ctx", panel.FieldModelContextWindow, "deepseek", model40, "64000")
	disk = diskConfig40(t, dir)
	want := wantFreshRow40()
	want.ContextWindow = 64000
	assertRowFields40(t, "after context_window", disk.LLM.Providers["deepseek"].Models[model40], want)

	writeModel40(t, disp, "r-40-audio", panel.FieldModelAudioIn, "deepseek", model40, "true")
	disk = diskConfig40(t, dir)
	want.Capabilities.AudioIn = true
	assertRowFields40(t, "after audio_in", disk.LLM.Providers["deepseek"].Models[model40], want)

	// READ-BACK FACE: config.get must answer what the disk now holds. A disk
	// write the read face never picks up reddens here.
	reply, err := disp.Handle(context.Background(), `{"method":"config.get","requestId":"r-40-read","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get: %v", err)
	}
	payload := payload40(t, reply, panel.SettingsPayloadMarker)
	providers, _ := payload["providers"].([]any)
	if len(providers) == 0 {
		t.Fatal("the read face answered no providers at all")
	}
	sawModel := false
	for _, p := range providers {
		rowMap, _ := p.(map[string]any)
		if rowMap["name"] != "deepseek" {
			continue
		}
		models, _ := rowMap["models"].([]any)
		for _, m := range models {
			if m == model40 {
				sawModel = true
			}
		}
		caps, _ := rowMap["modelCapabilities"].(map[string]any)
		bits, _ := caps[model40].(map[string]any)
		if bits == nil || bits["audioIn"] != true {
			t.Errorf("the read face does not answer modelCapabilities.audioIn=true for %s: %v", model40, caps)
		}
	}
	if !sawModel {
		t.Errorf("the read face never lists the disk row %s: %v", model40, payload["providers"])
	}

	// DELETE. The disk must lose exactly that entry, and the read face must
	// stop listing it.
	writeModel40(t, disp, "r-40-del", panel.FieldModelRowDelete, "deepseek", model40, "")
	disk = diskConfig40(t, dir)
	if _, still := disk.LLM.Providers["deepseek"].Models[model40]; still {
		t.Fatalf("the deleted entry is still on disk:\n%s", body40(t, dir))
	}
	if _, keep := disk.LLM.Providers["deepseek"].Models["deepseek-chat"]; !keep {
		t.Fatal("the delete took the seed's own entry with it")
	}
	reply, err = disp.Handle(context.Background(), `{"method":"config.get","requestId":"r-40-read2","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get after delete: %v", err)
	}
	payload = payload40(t, reply, panel.SettingsPayloadMarker)
	for _, p := range payload["providers"].([]any) {
		rowMap, _ := p.(map[string]any)
		if rowMap["name"] != "deepseek" {
			continue
		}
		models, _ := rowMap["models"].([]any)
		for _, m := range models {
			if m == model40 {
				t.Errorf("the read face still lists the deleted entry %s", model40)
			}
		}
	}
}

// body40 reads the raw config bytes for failure messages.
func body40(t *testing.T, dir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		return fmt.Sprintf("(unreadable: %v)", err)
	}
	return string(b)
}

// Test40ModelMgmtEndpointWriteSaysRestartVerbatim pins dispatch item 3 on the
// field the dispatch names: an [llm] endpoint write's receipt carries the 39
// AC#1 restart sentence VERBATIM (tierSentence(EffectiveRestart)) and never a
// live-save wording. The sentence is pinned as the literal the page is shown -
// the wording is the 39 AC#1 machine-checked contract, so a copy here is a
// drift alarm, not a second vocabulary.
func Test40ModelMgmtEndpointWriteSaysRestartVerbatim(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	reply := writeModel40(t, disp, "r-40-endpoint", panel.FieldProviderBaseURL, "deepseek", "", "https://api.example.invalid/v1")
	const wantRestart40 = "这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。需要重启：本次运行不会读取新值。"
	if !strings.Contains(reply, wantRestart40) {
		t.Errorf("the endpoint write's receipt lost the verbatim restart sentence:\n%q\nwant it to contain %q", reply, wantRestart40)
	}
	if strings.Contains(reply, "立即生效") || strings.Contains(reply, "保存即生效") {
		t.Errorf("the endpoint write's receipt claims a live save:\n%q", reply)
	}
}

// ---------------------------------------------------------------------------
// Half 2: the probe button.
// ---------------------------------------------------------------------------

// probeFixture40 is the probe fixture over the PRODUCTION dispatch chain: the
// config names the live mockllm, and newComposerDispatchChain wires the probe
// leg the same way the resident boot does.
type probeFixture40 struct {
	t    *testing.T
	srv  *adaptertest.Mockllm
	disp *panel.ComposerDispatch
	dir  string
}

func newProbeFixture40(t *testing.T) *probeFixture40 {
	t.Helper()
	srv := adaptertest.StartMockllm(t)
	srv.Reset(t)
	f := &probeFixture40{t: t, srv: srv, dir: sealableTempDir124(t)}
	f.writeConfig40(srv.Base + "/v1")
	disp, _, err := newComposerDispatchChain(f.dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("newComposerDispatchChain: %v", err)
	}
	f.disp = disp
	return f
}

// writeConfig40 lays the fixture config with one provider row at baseURL.
func (f *probeFixture40) writeConfig40(baseURL string) {
	f.t.Helper()
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/mock-small"]

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q

[llm.providers.acme.models.mock-small]
enabled = true
context_window = 128000

[llm.providers.acme.models.mock-small.capabilities]
text = true
fc = true
vision = true
thinking = true
`, baseURL)
	if err := os.WriteFile(filepath.Join(f.dir, configFileName), []byte(body), 0o600); err != nil {
		f.t.Fatal(err)
	}
}

func (f *probeFixture40) probe(provider, model string) (string, error) {
	f.t.Helper()
	raw, err := json.Marshal(map[string]string{
		"method": "config.model.probe", "requestId": "r-40-probe", "source": "panel-composer",
		"configProvider": provider, "configModel": model,
	})
	if err != nil {
		f.t.Fatal(err)
	}
	return f.disp.Handle(context.Background(), string(raw))
}

func Test40ModelProbeMeasuresTheRealProviderAndWritesProviderHealth(t *testing.T) {
	f := newProbeFixture40(t)

	reply, err := f.probe("acme", "mock-small")
	if err != nil {
		t.Fatalf("probe refused: %v", err)
	}
	payload := payload40(t, reply, panel.ModelProbePayloadMarker)
	if payload["provider"] != "acme" || payload["model"] != "mock-small" {
		t.Errorf("the payload does not name the measured pair verbatim: %v", payload)
	}
	if payload["written"] != true {
		t.Errorf("the run's verdicts did not reach provider_health: %v", payload)
	}
	latency, _ := payload["latencyMS"].(float64)
	if latency < 0 {
		t.Errorf("latencyMS = %v, want a measured number >= 0", payload["latencyMS"])
	}
	outcomes, _ := payload["outcomes"].([]any)
	if len(outcomes) == 0 {
		t.Fatal("the payload carries no per-capability verdicts (an empty run cannot read as a success)")
	}
	for _, o := range outcomes {
		om, _ := o.(map[string]any)
		if om["attempted"] != true || om["ok"] != true {
			t.Errorf("a capable server measured %+v, want attempted+ok for every capability", om)
		}
	}
	if !strings.Contains(reply, "实测 p50") {
		t.Errorf("the receipt does not carry the measured latency sentence: %q", reply)
	}
	// The receipt measures; it never saves. No tier wording may appear on a
	// probe receipt (the wrong-tier lie, this direction).
	if strings.Contains(reply, "需要重启") || strings.Contains(reply, "立即生效") || strings.Contains(reply, "设置已写入") {
		t.Errorf("the probe receipt carries save/tier wording:\n%q", reply)
	}
	// The measurement really recorded: provider_health holds the row.
	store, err := memory.Open(f.dir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = store.Close() }()
	h, err := store.ProviderHealthRow(context.Background(), "acme", "mock-small")
	if err != nil {
		t.Fatalf("provider_health: %v", err)
	}
	if h.Probe.FC == nil || !*h.Probe.FC {
		t.Errorf("provider_health fc = %+v, want measured true", h.Probe)
	}
}

// repoint40 rewrites the fixture config's base_url to baseURL on disk and
// rebuilds the chain, so the probe resolves the endpoint the way production
// does (fresh Manager over the file).
func (f *probeFixture40) repoint40(baseURL string) {
	f.t.Helper()
	f.writeConfig40(baseURL)
	disp, _, err := newComposerDispatchChain(f.dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		f.t.Fatalf("newComposerDispatchChain (repoint): %v", err)
	}
	f.disp = disp
}

// append40 inserts one TOML line after anchor and rebuilds the chain.
func (f *probeFixture40) append40(line, anchor string) {
	f.t.Helper()
	path := filepath.Join(f.dir, configFileName)
	body, err := os.ReadFile(path)
	if err != nil {
		f.t.Fatal(err)
	}
	rewritten := strings.Replace(string(body), anchor, anchor+"\n"+line, 1)
	if rewritten == string(body) {
		f.t.Fatalf("anchor %q not found", anchor)
	}
	if err := os.WriteFile(path, []byte(rewritten), 0o600); err != nil {
		f.t.Fatal(err)
	}
	disp, _, err := newComposerDispatchChain(f.dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		f.t.Fatalf("newComposerDispatchChain (append): %v", err)
	}
	f.disp = disp
}

func Test40ModelProbeReportsTheMismatchAndTheRawDetail(t *testing.T) {
	f := newProbeFixture40(t)
	f.srv.Control(f.t, "/__control/capability", `{"thinking":"broken"}`)

	reply, err := f.probe("acme", "mock-small")
	if err != nil {
		t.Fatalf("probe refused: %v", err)
	}
	payload := payload40(t, reply, panel.ModelProbePayloadMarker)
	mismatches, _ := payload["mismatches"].([]any)
	found := false
	for _, m := range mismatches {
		mm, _ := m.(map[string]any)
		if mm["capability"] == "thinking" {
			found = true
			if mm["label"] == "" || mm["detail"] == "" {
				t.Errorf("the thinking mismatch carries no frozen label or detail: %+v", mm)
			}
		}
	}
	if !found {
		t.Errorf("the payload carries no thinking mismatch: %v", mismatches)
	}
	if !strings.Contains(reply, "能力实测不符") || !strings.Contains(reply, "thinking") {
		t.Errorf("the receipt does not announce the declared-vs-measured fact:\n%q", reply)
	}
	// Declared-true + measured-false: the receipt pairs the words, the page's
	// "declared ✓ / measured ✗" fact in words the emoji ban allows.
	if !strings.Contains(reply, "声明可用") || !strings.Contains(reply, "实测不可用") {
		t.Errorf("the receipt does not pair declared vs measured wording:\n%q", reply)
	}
}

func Test40ModelProbeFailureCarriesTheRealErrorNeverOK(t *testing.T) {
	// A base URL nothing listens on: the suite's own probes fail, and the
	// receipt must carry that failure原文 per capability - never a fabricated
	// pass. (The suite measures; a dead endpoint is a measured "实测不可用"
	// with the real error in the detail, exactly the shape the CLI's probe
	// report prints. The detail's exact wording is the stack's own answer on
	// this machine - the tooth pins that it is ECHOED, not what it says.)
	f := newProbeFixture40(t)
	f.repoint40("http://127.0.0.1:1/v1")

	reply, err := f.probe("acme", "mock-small")
	if err != nil {
		t.Fatalf("a dead endpoint is a measured failure, not a route refusal: %v", err)
	}
	payload := payload40(t, reply, panel.ModelProbePayloadMarker)
	if payload["allOK"] != false {
		t.Errorf("allOK = %v, want false against a dead endpoint", payload["allOK"])
	}
	outcomes, _ := payload["outcomes"].([]any)
	if len(outcomes) == 0 {
		t.Fatal("the dead-endpoint run carries no verdicts at all")
	}
	detail := ""
	for _, o := range outcomes {
		om, _ := o.(map[string]any)
		if om["ok"] != false || om["attempted"] != true {
			t.Errorf("a dead endpoint measured %+v, want attempted+fail", om)
		}
		d, _ := om["detail"].(string)
		if d == "" {
			t.Errorf("a failed outcome carries no error原文: %+v", om)
		}
		if detail == "" {
			detail = d
		} else if d != detail {
			t.Errorf("outcomes disagree about the failure (%q vs %q) - one of them is not the real answer", detail, d)
		}
	}
	if !strings.Contains(reply, detail) {
		t.Errorf("the receipt does not echo the failure原文 %q:\n%q", detail, reply)
	}
	if !strings.Contains(reply, "实测不可用（细节：") {
		t.Errorf("the receipt does not carry the measured failure with its raw detail:\n%q", reply)
	}
	if strings.Contains(reply, "实测可用") {
		t.Errorf("the receipt claims a capability worked against a dead endpoint:\n%q", reply)
	}
}

// Test40ModelProbeUnconfiguredRefIsNotSilentlyKeyless pins the STRUCTURAL
// failure path: a ref that does not resolve refuses before any request is
// sent, and the refusal sentence carries the error原文 (Handle's err path
// renders it through RefusedEnvelopeForUser - the same shape every settings
// refusal takes; no payload line, no verdicts, nothing recorded).
func Test40ModelProbeUnconfiguredRefIsNotSilentlyKeyless(t *testing.T) {
	f := newProbeFixture40(t)
	f.append40(`api_key_ref = "dpapi:absent-40"`, "[llm.providers.acme]")

	reply, err := f.probe("acme", "mock-small")
	if err == nil {
		t.Fatal("a probe with an unresolvable key must not answer success")
	}
	if !strings.Contains(reply, "Unconfigured") {
		t.Errorf("the refusal does not carry the Unconfigured原文:\n%q", reply)
	}
	if strings.Contains(reply, panel.ModelProbePayloadMarker) {
		t.Errorf("a refused run carries a payload line:\n%q", reply)
	}
	if strings.Contains(reply, "实测可用") || strings.Contains(reply, "探测完成") {
		t.Errorf("a refused run reads as a measurement:\n%q", reply)
	}
}

func Test40ModelProbeRefusesByNameAndWiringNail(t *testing.T) {
	f := newProbeFixture40(t)

	// Unknown provider / unknown model: refused BY NAME, no measurement.
	if _, err := f.probe("ghostco", "mock-small"); err == nil || !strings.Contains(err.Error(), "没有服务商") {
		t.Errorf("an unknown provider was not refused by name: %v", err)
	}
	if _, err := f.probe("acme", "no-such-model-40"); err == nil || !strings.Contains(err.Error(), "没有模型") {
		t.Errorf("an unknown model was not refused by name: %v", err)
	}

	// The wiring nail: the production chain answers the door through the real
	// machinery. Source pins (the same shape Test40PrivacyWiringNail does):
	// unplugging any of these lines must redden this pin.
	nail40(t, filepath.Join("cmd", "wisp", "panel_inbound.go"),
		"newModelProbeHandler40(dir, mgr, secrets, auditf, actor)",
		"ModelProbes: probes")
	nail40(t, filepath.Join("cmd", "wisp", "panel_probe_40.go"),
		"llm.RunProbeSuite(", "storeHealthSink{store: store}", "llm.NewResolver(cfg, h.secrets)")
}

// nail40 reads one production file from the repo root and requires each pin
// byte-for-byte. The failure names the file and the pin, so an unplug is a
// one-line fix away from red-to-green.
func nail40(t *testing.T, rel string, pins ...string) {
	t.Helper()
	root := panelRepoRoot40(t)
	src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	for _, pin := range pins {
		if !strings.Contains(string(src), pin) {
			t.Errorf("%s lost the production pin %q - the probe door is no longer on the real machinery", rel, pin)
		}
	}
}

func panelRepoRoot40(t *testing.T) string {
	t.Helper()
	// The package's own directory is cmd/wisp; the repo root is two up.
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(filepath.Dir(dir))
}
