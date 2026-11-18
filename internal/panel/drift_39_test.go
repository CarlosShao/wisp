package panel

// Ticket 39 AC#4's drift-handshake teeth, panel half (2026-11-18):
//
//   - the config.set envelope's configBaseSHA reaches the store as
//     SettingWrite.BaseSHA, trimmed, verbatim - the gate never guesses it;
//   - the config.get payload carries the leg's fingerprint as baseSHA, and
//     keeps the key OFF the wire when the leg answered none (absent =
//     "this leg did not answer", never a fake hash);
//   - an accepted write's reload event carries the leg's post-write
//     fingerprint as baseSHA, and omits it when the leg answered none.

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// captureStore39 records the SettingWrite it was handed, so a test can judge
// the envelope->store passthrough instead of the leg's own behaviour.
type captureStore39 struct {
	captured SettingWrite
	view     SettingsView
	apply    SettingWriteResult
}

func (s *captureStore39) ReadSettings(context.Context) (SettingsView, error) { return s.view, nil }
func (s *captureStore39) ApplySetting(_ context.Context, w SettingWrite) (SettingWriteResult, error) {
	s.captured = w
	return s.apply, nil
}
func (s *captureStore39) StoreCredential(context.Context, string) (SettingWriteResult, error) {
	return SettingWriteResult{}, nil
}

func driftRequest39(field, value, baseSHA string) ComposerRequest {
	return ComposerRequest{
		Method: MethodConfigSet, RequestID: "rid-drift-39", Source: ComposerRequestSource,
		ConfigField: field, ConfigValue: value, ConfigBaseSHA: baseSHA,
	}
}

func TestConfigSetEnvelopeCarriesBaseSHAToTheStore(t *testing.T) {
	store := &captureStore39{apply: SettingWriteResult{Field: "ball_size", Written: []string{"ball.size"}, Tier: EffectiveRestart}}
	h := &ConfigWriteHandler{Store: store}
	if _, err := h.HandleConfigRequest(context.Background(), driftRequest39("ball_size", "90", "  abc123def  "), ""); err != nil {
		t.Fatalf("config.set refused: %v", err)
	}
	if store.captured.BaseSHA != "abc123def" {
		t.Errorf("the store received BaseSHA %q, want the envelope's echo trimmed to %q",
			store.captured.BaseSHA, "abc123def")
	}
	// The legacy shape: no fingerprint on the envelope = none on the write -
	// the pre-39 callers keep working byte-identically.
	store2 := &captureStore39{}
	h2 := &ConfigWriteHandler{Store: store2}
	if _, err := h2.HandleConfigRequest(context.Background(), driftRequest39("ball_size", "90", ""), ""); err != nil {
		t.Fatalf("legacy config.set refused: %v", err)
	}
	if store2.captured.BaseSHA != "" {
		t.Errorf("an envelope without configBaseSHA produced BaseSHA %q, want empty", store2.captured.BaseSHA)
	}
}

func TestConfigGetPayloadCarriesBaseSHAOnlyWhenAnswered(t *testing.T) {
	// The answered half: the fingerprint rides verbatim.
	h := &ConfigWriteHandler{Store: &captureStore39{view: SettingsView{
		Readable: true, BaseSHA: "deadbeef01",
	}}}
	reply, err := h.HandleConfigRequest(context.Background(),
		ComposerRequest{Method: MethodConfigGet, RequestID: "rid-drift-39", Source: ComposerRequestSource}, "")
	if err != nil {
		t.Fatalf("config.get refused: %v", err)
	}
	var payloadLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, SettingsPayloadMarker) {
			payloadLine = strings.TrimPrefix(line, SettingsPayloadMarker)
		}
	}
	var payload struct {
		BaseSHA string `json:"baseSHA"`
	}
	if err := json.Unmarshal([]byte(payloadLine), &payload); err != nil {
		t.Fatalf("payload is not JSON: %v\n%s", err, payloadLine)
	}
	if payload.BaseSHA != "deadbeef01" {
		t.Errorf("payload baseSHA = %q, want the leg's own fingerprint verbatim", payload.BaseSHA)
	}

	// The unanswered half: the key stays OFF the wire - an absent fingerprint
	// is "the leg did not answer", never a fabricated hash.
	h2 := &ConfigWriteHandler{Store: &captureStore39{view: SettingsView{Readable: true}}}
	reply2, err := h2.HandleConfigRequest(context.Background(),
		ComposerRequest{Method: MethodConfigGet, RequestID: "rid-drift-39", Source: ComposerRequestSource}, "")
	if err != nil {
		t.Fatalf("config.get refused: %v", err)
	}
	for _, line := range strings.Split(reply2, "\n") {
		if strings.HasPrefix(line, SettingsPayloadMarker) && strings.Contains(line, `"baseSHA"`) {
			t.Errorf("a leg that answered no fingerprint put baseSHA on the wire: %s", line)
		}
	}
}

func TestReloadEventCarriesBaseSHAOnlyWhenAnswered(t *testing.T) {
	h := &ConfigWriteHandler{Store: &captureStore39{
		apply: SettingWriteResult{Field: "ball_size", Written: []string{"ball.size"}, Tier: EffectiveRestart, BaseSHAAfter: "postwrite-sha"},
	}}
	reply, err := h.HandleConfigRequest(context.Background(), driftRequest39("ball_size", "90", ""), "")
	if err != nil {
		t.Fatalf("config.set refused: %v", err)
	}
	var eventLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, ConfigReloadEventMarker) {
			eventLine = strings.TrimPrefix(line, ConfigReloadEventMarker)
		}
	}
	var ev struct {
		Event   string `json:"event"`
		BaseSHA string `json:"baseSHA"`
	}
	if err := json.Unmarshal([]byte(eventLine), &ev); err != nil {
		t.Fatalf("event line is not JSON: %v\n%s", err, eventLine)
	}
	if ev.Event != "config.reload" || ev.BaseSHA != "postwrite-sha" {
		t.Errorf("reload event = %+v, want event config.reload with the post-write fingerprint", ev)
	}

	// The unanswered half: no baseSHA key, the page's cache stays put.
	h2 := &ConfigWriteHandler{Store: &captureStore39{
		apply: SettingWriteResult{Field: "ball_size", Written: []string{"ball.size"}, Tier: EffectiveRestart},
	}}
	reply2, err := h2.HandleConfigRequest(context.Background(), driftRequest39("ball_size", "90", ""), "")
	if err != nil {
		t.Fatalf("config.set refused: %v", err)
	}
	for _, line := range strings.Split(reply2, "\n") {
		if strings.HasPrefix(line, ConfigReloadEventMarker) && strings.Contains(line, `"baseSHA"`) {
			t.Errorf("a write result without a post-write fingerprint put baseSHA on the event: %s", line)
		}
	}
}
