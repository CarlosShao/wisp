package main

// 313-v1 acceptance probe (non-implementer). Overlay-added into cmd/wisp so it
// rides the production legs: real config.Manager, real secret.Store, the real
// panel.ComposerDispatch exactly as run.go:431 / panel_inbound.go:279 assemble
// it (newConfigStore). No fake stands in for the code under judgement; the only
// stub is nowhere here.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/secret"
)

func v1Leg(t *testing.T, dir string) (*panel.ComposerDispatch, *config.Manager, *[]string) {
	t.Helper()
	mgr, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	secrets, err := secret.NewStore(dir)
	if err != nil {
		t.Fatalf("secret.NewStore: %v", err)
	}
	audit := []string{}
	auditf := panel.AuditFunc(func(format string, args ...any) {
		audit = append(audit, fmt.Sprintf(format, args...))
	})
	leg := newConfigStore(mgr, secrets, auditf, "probe-313v1")
	disp := &panel.ComposerDispatch{
		Config: &panel.ConfigWriteHandler{Store: leg, Audit: auditf, Actor: "probe-313v1"},
		Audit:  auditf,
	}
	return disp, mgr, &audit
}

func v1PayloadLine(t *testing.T, reply string) string {
	t.Helper()
	i := strings.Index(reply, panel.SettingsPayloadMarker)
	if i < 0 {
		t.Fatalf("no payload marker in reply: %q", reply)
	}
	return strings.SplitN(reply[i+len(panel.SettingsPayloadMarker):], "\n", 2)[0]
}

func v1Set(t *testing.T, disp *panel.ComposerDispatch, id, field, provider, model, value string) string {
	t.Helper()
	raw, err := json.Marshal(map[string]string{
		"method": "config.set", "requestId": id, "source": "panel-composer",
		"configField": field, "configProvider": provider, "configModel": model, "configValue": value,
	})
	if err != nil {
		t.Fatalf("marshal env: %v", err)
	}
	reply, err := disp.Handle(context.Background(), string(raw))
	if err != nil {
		t.Fatalf("config.set %s (%s): %v", id, field, err)
	}
	return reply
}

func v1Get(t *testing.T, disp *panel.ComposerDispatch, id string) string {
	t.Helper()
	raw := fmt.Sprintf(`{"method":"config.get","requestId":%q,"source":"panel-composer"}`, id)
	reply, err := disp.Handle(context.Background(), raw)
	if err != nil {
		t.Fatalf("config.get %s: %v", id, err)
	}
	return reply
}

// TestV1ProbeReadPayloadAndPageFlow drives the exact three-step flow the page
// spec (313-panel-ask §3) recommends, on an empty config, then reads the payload
// back. Payload lines are dumped verbatim for the three-way key comparison.
func TestV1ProbeReadPayloadAndPageFlow(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte("schema_version = 2\n"), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	disp, _, audit := v1Leg(t, dir)

	reply := v1Get(t, disp, "v1-get-empty")
	fmt.Printf("V1-EMPTY-REPLY-BEGIN\n%s\nV1-EMPTY-REPLY-END\n", reply)
	fmt.Printf("V1-EMPTY-PAYLOAD %s\n", v1PayloadLine(t, reply))

	// Step 1: explicit row creation, with the model half so the picker has an id.
	r := v1Set(t, disp, "v1-s1", panel.FieldProviderRowCreate, "stepfun", "model-under-test-313", "")
	fmt.Printf("V1-RECEIPT-ROW %s\n", r)
	// Step 2: credential through the write-only envelope key.
	canary := "sk-v1probe-canary-313"
	rawCred, _ := json.Marshal(map[string]string{
		"method": "config.set", "requestId": "v1-s2", "source": "panel-composer",
		"configField": panel.FieldProviderCredential, "configProvider": "stepfun", "credentialValue": canary,
	})
	rc, err := disp.Handle(context.Background(), string(rawCred))
	if err != nil {
		t.Fatalf("credential write: %v", err)
	}
	fmt.Printf("V1-RECEIPT-CRED %s\n", rc)
	// Step 3: the pair.
	r = v1Set(t, disp, "v1-s3", panel.FieldRoleChatPair, "stepfun", "model-under-test-313", "")
	fmt.Printf("V1-RECEIPT-PAIR %s\n", r)

	body, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	for name, surface := range map[string]string{"receipt-cred": rc, "disk": string(body)} {
		if strings.Contains(surface, canary) {
			t.Errorf("canary appears in %s:\n%s", name, surface)
		}
	}
	for _, a := range *audit {
		if strings.Contains(a, canary) {
			t.Errorf("canary appears in an audit line: %s", a)
		}
	}
	fmt.Printf("V1-DISK-AFTER-THREE-STEPS-BEGIN\n%s\nV1-DISK-AFTER-THREE-STEPS-END\n", string(body))

	reply = v1Get(t, disp, "v1-get-full")
	fmt.Printf("V1-FULL-PAYLOAD %s\n", v1PayloadLine(t, reply))
}

// TestV1ProbeMultiElementChainAndSeedPayload: a two-element comma chain written
// through the door (the comma split has no repository pin - this is the probe
// that reads it), on a seeded catalog with two models, plus the payload dump.
func TestV1ProbeMultiElementChainAndSeedPayload(t *testing.T) {
	dir := t.TempDir()
	seed := `schema_version = 2

[llm.providers.stepfun]
protocol = "openai-chat"
base_url = "https://api.stepfun.com/v1"
api_key_ref = "env:STEPFUN_KEY"

[llm.providers.stepfun.models.step-2]
context_window = 32000

[llm.providers.stepfun.models.step-2-mini]
context_window = 16000
`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(seed), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	disp, _, _ := v1Leg(t, dir)
	r := v1Set(t, disp, "v1-s-asr", panel.FieldVoiceCloudASRChain, "", "", "stepfun/step-2,stepfun/step-2-mini")
	fmt.Printf("V1-RECEIPT-2ELT %s\n", r)
	r = v1Set(t, disp, "v1-s-asr1", panel.FieldVoiceCloudTTSChain, "", "", "stepfun/step-2")
	fmt.Printf("V1-RECEIPT-TTS %s\n", r)

	body, _ := os.ReadFile(filepath.Join(dir, configFileName))
	fmt.Printf("V1-DISK-2ELT-BEGIN\n%s\nV1-DISK-2ELT-END\n", string(body))
	line := v1PayloadLine(t, v1Get(t, disp, "v1-get-2elt"))
	fmt.Printf("V1-SEED-PAYLOAD %s\n", line)
	var payload struct {
		Voice struct {
			CloudASRChain []string `json:"cloudASRChain"`
			CloudTTSChain []string `json:"cloudTTSChain"`
		} `json:"voice"`
	}
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("payload: %v", err)
	}
	if len(payload.Voice.CloudASRChain) != 2 ||
		payload.Voice.CloudASRChain[0] != "stepfun/step-2" || payload.Voice.CloudASRChain[1] != "stepfun/step-2-mini" {
		t.Errorf("the two-element chain did not round-trip: %v", payload.Voice.CloudASRChain)
	}
}

// TestV1ProbeUnionRefusalText reads the door's refusal sentence for an unlisted
// name at baseline: it must name the full union of both rosters.
func TestV1ProbeUnionRefusalText(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte("schema_version = 2\n"), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	disp, _, _ := v1Leg(t, dir)
	raw := `{"method":"config.set","requestId":"v1-u","source":"panel-composer","configField":"v1_no_such_field"}`
	_, err := disp.Handle(context.Background(), raw)
	if err == nil {
		t.Fatal("an unlisted field was accepted")
	}
	fmt.Printf("V1-UNION-REFUSAL %v\n", err)
}

// TestV1ProbeRefusalsLeaveTheFileByteIdentical: every 313 refusal surface, with
// the sha256 of config.toml before and after each one; a legal control write
// proves the ruler is sensitive.
func TestV1ProbeRefusalsLeaveTheFileByteIdentical(t *testing.T) {
	dir := t.TempDir()
	seed := `schema_version = 2

[llm.providers.stepfun]
protocol = "openai-chat"
base_url = "https://api.stepfun.com/v1"

[llm.providers.stepfun.models.step-2]
context_window = 32000
`
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(seed), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	disp, mgr, _ := v1Leg(t, dir)
	hash := func() string {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read for hash: %v", err)
		}
		s := sha256.Sum256(b)
		return hex.EncodeToString(s[:])
	}
	doorErr := func(id, field, provider, model, value string) error {
		raw, _ := json.Marshal(map[string]string{
			"method": "config.set", "requestId": id, "source": "panel-composer",
			"configField": field, "configProvider": provider, "configModel": model, "configValue": value,
		})
		_, err := disp.Handle(context.Background(), string(raw))
		return err
	}
	attempts := []struct {
		name string
		run  func() error
	}{
		{"door: role_chat_pair without model", func() error { return doorErr("v1-r1", panel.FieldRoleChatPair, "stepfun", "", "") }},
		{"door: role_chat_pair without provider", func() error { return doorErr("v1-r2", panel.FieldRoleChatPair, "", "step-2", "") }},
		{"leg: SetRoleChatPair half", func() error { _, e := mgr.SetRoleChatPair("stepfun", ""); return e }},
		{"leg: SetVoiceCloudASRChain ghost model", func() error { _, e := mgr.SetVoiceCloudASRChain([]string{"stepfun/ghost"}); return e }},
		{"door: voice_cloud_asr_chain ghost element", func() error { return doorErr("v1-r5", panel.FieldVoiceCloudASRChain, "", "", "stepfun/ghost") }},
		{"door: voice_cloud_asr_chain local-sherpa", func() error { return doorErr("v1-r6", panel.FieldVoiceCloudASRChain, "", "", "local-sherpa/anything") }},
		{"door: provider_row_create non-preset", func() error { return doorErr("v1-r7", panel.FieldProviderRowCreate, "inventco", "m", "") }},
		{"door: provider_row_create existing row", func() error { return doorErr("v1-r8", panel.FieldProviderRowCreate, "stepfun", "step-2", "") }},
		{"leg: CreateProviderRow non-preset", func() error { _, e := mgr.CreateProviderRow("inventco", "m"); return e }},
		{"leg: SetVoiceCloudTTSChain empty element", func() error { _, e := mgr.SetVoiceCloudTTSChain([]string{"stepfun/step-2", " "}); return e }},
	}
	before := hash()
	fmt.Printf("V1-HASH-BEFORE %s\n", before)
	for _, a := range attempts {
		err := a.run()
		after := hash()
		fmt.Printf("V1-REFUSAL %s | refused=%v hash-unchanged=%v | err=%v\n", a.name, err != nil, after == before, err)
		if err == nil {
			t.Errorf("%s was accepted", a.name)
		}
		if after != before {
			t.Errorf("%s changed config.toml (%s -> %s)", a.name, before, after)
		}
	}
	if _, err := mgr.SetVoiceCloudASRChain([]string{"stepfun/step-2"}); err != nil {
		t.Fatalf("control legal write: %v", err)
	}
	moved := hash() != before
	fmt.Printf("V1-HASH-AFTER-CONTROL %s moved=%v\n", hash(), moved)
	if !moved {
		t.Error("the control write did not move the hash: the ruler is dead")
	}
}
