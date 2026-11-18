package main

// Ticket 40 AC#3 后三半 (2026-11-18) - the teeth the dispatch names, on the
// production dispatch chain:
//
//	1. chain drag-order PERSISTED: the chain's order is whatever the page sent
//	   last - a reorder write lands on config.toml in the NEW order, and the
//	   reader is the DISK (a fresh LoadFile, diskConfig40), so a reorder that
//	   only moved in-memory state reads back the OLD order and reddens. The
//	   read face is asserted off the same disk truth, so a disk order the
//	   payload never picks up reddens too.
//	2. quota edits ENFORCED RANGES: the band comes from the codebase's own
//	   quantities (schema.go's "0 = uncapped" / quota_44.go's "zero = not
//	   entered", int64 micro-USD end to end) - negative is refused BY NAME and
//	   NOTHING is written (the fresh LoadFile still reads the old leaf), the
//	   type's own maximum is the only ceiling (ParseInt enforces it; no
//	   invented upper number anywhere), and an in-range write round-trips to
//	   disk and to the read face.
//	3. the read face carries the realtime bit and the quota values the page's
//	   capability filter and quota editor consume - the bits come from the
//	   model rows' real fields (spec.Capabilities.*), never a page-side list.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// secondModel40b is a made-up catalog id (the 313 rule: no real model names).
const secondModel40b = "model-under-test-40b"

// wantRestart40b is the 39 AC#1 restart sentence VERBATIM - the wording is a
// machine-checked contract (config_handlers.go's tierSentence), so a copy here
// is a drift alarm, not a second vocabulary.
const wantRestart40b = "这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。需要重启：本次运行不会读取新值。"

// ---------------------------------------------------------------------------
// Half 2: chain drag-order persisted (fresh LoadFile reads back the sent order).
// ---------------------------------------------------------------------------

func Test40VoiceChainReorderPersistsToDiskAndReadsBackSameOrder(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	// A second catalog row, so the chain has two elements to reorder.
	writeModel40(t, disp, "r-40b-row", panel.FieldModelRowCreate, "deepseek", secondModel40b, "")

	first40b := "deepseek/deepseek-chat"
	second40b := "deepseek/" + secondModel40b

	// ASR chain, saved in one order...
	writeModel40(t, disp, "r-40b-asr-1", panel.FieldVoiceCloudASRChain, "", "", first40b+","+second40b)
	// ...then the DRAG: the same two elements, sent back in the new order.
	writeModel40(t, disp, "r-40b-asr-2", panel.FieldVoiceCloudASRChain, "", "", second40b+","+first40b)

	// The reader is the DISK. Only-memory reorders read back the old order here.
	disk := diskConfig40(t, dir)
	got := disk.Voice.CloudASRChain
	if len(got) != 2 || got[0] != second40b || got[1] != first40b {
		t.Errorf("the reordered ASR chain did not land on disk in the new order: %v", got)
	}

	// TTS chain: same door, same order rule.
	writeModel40(t, disp, "r-40b-tts-1", panel.FieldVoiceCloudTTSChain, "", "", first40b+","+second40b)
	writeModel40(t, disp, "r-40b-tts-2", panel.FieldVoiceCloudTTSChain, "", "", second40b+","+first40b)
	disk = diskConfig40(t, dir)
	got = disk.Voice.CloudTTSChain
	if len(got) != 2 || got[0] != second40b || got[1] != first40b {
		t.Errorf("the reordered TTS chain did not land on disk in the new order: %v", got)
	}

	// The read face answers the same disk order - a page reopening after the
	// restart must see exactly what was saved.
	reply, err := disp.Handle(context.Background(), `{"method":"config.get","requestId":"r-40b-read","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get: %v", err)
	}
	payload := payload40(t, reply, panel.SettingsPayloadMarker)
	voice, _ := payload["voice"].(map[string]any)
	asr, _ := voice["cloudASRChain"].([]any)
	if len(asr) != 2 || asr[0] != second40b || asr[1] != first40b {
		t.Errorf("the read face does not answer the saved ASR order: %v", voice["cloudASRChain"])
	}
	tts, _ := voice["cloudTTSChain"].([]any)
	if len(tts) != 2 || tts[0] != second40b || tts[1] != first40b {
		t.Errorf("the read face does not answer the saved TTS order: %v", voice["cloudTTSChain"])
	}

	// The gate keeps its existence ruler on a reorder too: an element the
	// catalog does not list is refused by name, nothing is written.
	_, err = disp.Handle(context.Background(), mustJSON40b(t, map[string]string{
		"method": "config.set", "requestId": "r-40b-ghost", "source": "panel-composer",
		"configField": panel.FieldVoiceCloudASRChain, "configValue": "deepseek/no-such-model-40b",
	}))
	if err == nil || !strings.Contains(err.Error(), "references unknown model") {
		t.Errorf("a ghost element was not refused on the reorder path: %v", err)
	}
	disk = diskConfig40(t, dir)
	if len(disk.Voice.CloudASRChain) != 2 || disk.Voice.CloudASRChain[0] != second40b {
		t.Errorf("a refused reorder changed the disk chain: %v", disk.Voice.CloudASRChain)
	}
}

// mustJSON40b marshals one envelope or fails the test.
func mustJSON40b(t *testing.T, env map[string]string) string {
	t.Helper()
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// ---------------------------------------------------------------------------
// Half 3: quota edits, ranges enforced, refusals write nothing.
// ---------------------------------------------------------------------------

func Test40QuotaOutOfRangeRefusedByNameAndWritesNothing(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	cases := []struct {
		name            string
		field           string
		provider, model string
		value           string
		wantErr         []string
	}{
		{
			name: "negative daily", field: panel.FieldModelQuotaDaily,
			provider: "deepseek", model: "deepseek-chat", value: "-1",
			wantErr: []string{"超出范围", "quota_daily_micro", "未落盘"},
		},
		{
			name: "negative monthly", field: panel.FieldModelQuotaMonthly,
			provider: "deepseek", model: "deepseek-chat", value: "-5",
			wantErr: []string{"超出范围", "quota_monthly_micro", "未落盘"},
		},
		{
			name: "negative plan credit", field: panel.FieldProviderPlanCredit,
			provider: "deepseek", model: "", value: "-9",
			wantErr: []string{"超出范围", "plan_credit_total_micro", "未落盘"},
		},
		{
			name: "non-integer", field: panel.FieldModelQuotaDaily,
			provider: "deepseek", model: "deepseek-chat", value: "abc",
			wantErr: []string{"解析不了"},
		},
		{
			// Past int64: the type's own ceiling is the only upper bound the
			// codebase answers, and the parse is its enforcer (no invented cap).
			name: "past int64", field: panel.FieldModelQuotaMonthly,
			provider: "deepseek", model: "deepseek-chat", value: "99999999999999999999",
			wantErr: []string{"解析不了"},
		},
	}
	for _, tc := range cases {
		env := map[string]string{
			"method": "config.set", "requestId": "r-40b-range", "source": "panel-composer",
			"configField": tc.field, "configProvider": tc.provider, "configValue": tc.value,
		}
		if tc.model != "" {
			env["configModel"] = tc.model
		}
		_, err := disp.Handle(context.Background(), mustJSON40b(t, env))
		if err == nil {
			t.Errorf("%s: an out-of-range quota write was accepted", tc.name)
			continue
		}
		for _, want := range tc.wantErr {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("%s: the refusal does not name %q: %v", tc.name, want, err)
			}
		}
	}

	// The tooth: a refused write must leave the disk untouched - every quota
	// leaf still reads its pre-write value (0 = uncapped/not entered) on a
	// FRESH load.
	disk := diskConfig40(t, dir)
	spec := disk.LLM.Providers["deepseek"].Models["deepseek-chat"]
	if spec.QuotaDailyMicro != 0 || spec.QuotaMonthlyMicro != 0 {
		t.Errorf("refused quota writes still landed on the model row: %d/%d",
			spec.QuotaDailyMicro, spec.QuotaMonthlyMicro)
	}
	if disk.LLM.Providers["deepseek"].PlanCreditTotalMicro != 0 {
		t.Errorf("a refused credit write still landed on the row: %d", disk.LLM.Providers["deepseek"].PlanCreditTotalMicro)
	}

	// The empty value is refused by the GATE (the field takes a value), before
	// any leg runs.
	_, err := disp.Handle(context.Background(), mustJSON40b(t, map[string]string{
		"method": "config.set", "requestId": "r-40b-empty", "source": "panel-composer",
		"configField": panel.FieldModelQuotaDaily, "configProvider": "deepseek",
		"configModel": "deepseek-chat", "configValue": "   ",
	}))
	if err == nil {
		t.Errorf("an empty quota value cleared the gate")
	}

	// Wiring nails: the door case and the setter family are production lines;
	// unplugging any of them reddens this pin.
	nail40(t, "cmd/wisp/panel_config_store.go",
		"case panel.FieldModelQuotaDaily, panel.FieldModelQuotaMonthly, panel.FieldProviderPlanCredit:",
		"s.mgr.SetModelQuotaDaily(w.Provider, w.Model, n)",
		"s.mgr.SetModelQuotaMonthly(w.Provider, w.Model, n)",
		"s.mgr.SetProviderPlanCreditTotal(w.Provider, n)")
	nail40(t, "internal/config/settings_quota_40.go",
		"超出范围（额度不能为负，0 = 不封顶）",
		"超出范围（credit 总额不能为负，0 = 未录入）")
}

func Test40QuotaInRangeRoundTripLandsOnDiskAndReadFace(t *testing.T) {
	dir := settingsDir248(t)
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	writeModel40(t, disp, "r-40b-daily", panel.FieldModelQuotaDaily, "deepseek", "deepseek-chat", "1000000")
	writeModel40(t, disp, "r-40b-monthly", panel.FieldModelQuotaMonthly, "deepseek", "deepseek-chat", "20000000")
	writeModel40(t, disp, "r-40b-credit", panel.FieldProviderPlanCredit, "deepseek", "", "5000000")

	disk := diskConfig40(t, dir)
	spec := disk.LLM.Providers["deepseek"].Models["deepseek-chat"]
	if spec.QuotaDailyMicro != 1000000 || spec.QuotaMonthlyMicro != 20000000 {
		t.Errorf("the in-range quotas did not land on disk: %d/%d", spec.QuotaDailyMicro, spec.QuotaMonthlyMicro)
	}
	if disk.LLM.Providers["deepseek"].PlanCreditTotalMicro != 5000000 {
		t.Errorf("the plan credit did not land on disk: %d", disk.LLM.Providers["deepseek"].PlanCreditTotalMicro)
	}

	// [llm] writes answer the restart tier with the verbatim sentence.
	reply, err := disp.Handle(context.Background(), mustJSON40b(t, map[string]string{
		"method": "config.set", "requestId": "r-40b-tier", "source": "panel-composer",
		"configField": panel.FieldModelQuotaDaily, "configProvider": "deepseek",
		"configModel": "deepseek-chat", "configValue": "3000000",
	}))
	if err != nil {
		t.Fatalf("in-range write refused: %v", err)
	}
	if !strings.Contains(reply, wantRestart40b) {
		t.Errorf("the quota receipt lost the verbatim restart sentence:\n%q", reply)
	}

	// The read face answers the disk values - the page's quota editor has a
	// real 当前 to show instead of a guess.
	reply, err = disp.Handle(context.Background(), `{"method":"config.get","requestId":"r-40b-read","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get: %v", err)
	}
	payload := payload40(t, reply, panel.SettingsPayloadMarker)
	var daily, monthly, credit float64
	found := false
	for _, p := range payload["providers"].([]any) {
		row, _ := p.(map[string]any)
		if row["name"] != "deepseek" {
			continue
		}
		found = true
		quotas, _ := row["modelQuotas"].(map[string]any)
		entry, _ := quotas["deepseek-chat"].(map[string]any)
		if entry == nil {
			t.Fatalf("the read face carries no modelQuotas entry for the row: %v", quotas)
		}
		if c, ok := entry["dailyMicro"].(float64); !ok || c != 3000000 {
			t.Errorf("modelQuotas.dailyMicro = %v, want 3000000 (the disk's current value)", entry["dailyMicro"])
		}
		daily, _ = entry["dailyMicro"].(float64)
		monthly, _ = entry["monthlyMicro"].(float64)
		credit, _ = row["planCreditTotalMicro"].(float64)
	}
	if !found {
		t.Fatal("the read face never answered the deepseek row")
	}
	if daily != 3000000 || monthly != 20000000 || credit != 5000000 {
		t.Errorf("read face values = %v/%v/%v, want 3000000/20000000/5000000", daily, monthly, credit)
	}
}

// ---------------------------------------------------------------------------
// Half 1's read face: the realtime bit rides the payload, read off the model
// rows' real fields.
// ---------------------------------------------------------------------------

// quotaSeed40b is the fixture: one non-preset row, two models - one declaring
// audio_in + realtime, one text-only. The bits are the schema's own fields.
const quotaSeed40b = `schema_version = 2

[llm]
text_chain = ["acme/text-only-40b"]

[llm.providers.acme]
protocol = "openai-chat"
base_url = "http://127.0.0.1:1/v1"

[llm.providers.acme.models.voice-model-40b]
enabled = true

[llm.providers.acme.models.voice-model-40b.capabilities]
audio_in = true
audio_out = true
realtime = true

[llm.providers.acme.models.text-only-40b]
enabled = true
context_window = 64000

[llm.providers.acme.models.text-only-40b.capabilities]
text = true
`

func Test40ReadFaceCarriesRealtimeBitAndQuotasFromTheModelRows(t *testing.T) {
	dir := sealableTempDir124(t)
	if err := writeFile40b(dir, quotaSeed40b); err != nil {
		t.Fatal(err)
	}
	var audit []string
	disp, _ := newRealSettingsLeg(t, dir, &audit, nil)

	reply, err := disp.Handle(context.Background(), `{"method":"config.get","requestId":"r-40b-caps","source":"panel-composer"}`)
	if err != nil {
		t.Fatalf("config.get: %v", err)
	}
	payload := payload40(t, reply, panel.SettingsPayloadMarker)
	for _, p := range payload["providers"].([]any) {
		row, _ := p.(map[string]any)
		if row["name"] != "acme" {
			continue
		}
		caps, _ := row["modelCapabilities"].(map[string]any)
		voice, _ := caps["voice-model-40b"].(map[string]any)
		if voice == nil || voice["realtime"] != true || voice["audioIn"] != true || voice["audioOut"] != true {
			t.Errorf("the voice model's declared bits did not ride the read face: %v", voice)
		}
		text, _ := caps["text-only-40b"].(map[string]any)
		if text == nil || text["realtime"] != false || text["audioIn"] != false {
			t.Errorf("the text model's bits did not ride the read face as answered falses: %v", text)
		}
		// The quota read face: a fresh row answers every listed model, 0 rides
		// as the real uncapped value.
		quotas, _ := row["modelQuotas"].(map[string]any)
		entry, _ := quotas["voice-model-40b"].(map[string]any)
		if entry == nil || entry["dailyMicro"] != float64(0) || entry["monthlyMicro"] != float64(0) {
			t.Errorf("the quota read face did not answer the fresh row's zeros: %v", entry)
		}
		if row["planCreditTotalMicro"] != float64(0) {
			t.Errorf("planCreditTotalMicro = %v, want the answered 0 (0 = not entered)", row["planCreditTotalMicro"])
		}
	}
}

func writeFile40b(dir string, body string) error {
	return os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600)
}
