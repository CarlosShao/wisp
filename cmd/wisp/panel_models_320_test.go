package main

// Ticket 320 AC#4 at the assembly root: the page's "拉取模型列表" reaches the
// fetch leg through the production dispatch route, and the receipt is the shape
// the page parses (one sentence plus ONE line beginning with the marker
// panel.ModelsListPayloadMarker, whose JSON carries models[].id).
//
// WHAT IS REAL, and why each piece is not a fake: one config.Manager over a
// seeded config.toml (the same fixture the 248/313/320 settings tests use), one
// secret.Store (the env-backed ref the seed already declares), the real
// ComposerDispatch route (params.parse -> dispatch -> the fetch leg), and a real
// HTTP round trip to a test server that speaks the OpenAI-shaped /models body.
// The only substituted thing is the remote provider itself, which is what a test
// server is for.
//
// WHAT THE REFUSAL CASES PIN: an unknown provider, a row with no base_url, and a
// ref that does not resolve must each come back as a named refusal WITHOUT
// touching the network - so the test server counts requests and the count has to
// stay at zero for those cases.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/secret"
)

// modelsPath320 turns one fetch into the receipt's parsed halves.
func modelsPath320(t *testing.T, disp *panel.ComposerDispatch, provider string) (string, map[string]any) {
	t.Helper()
	env := `{"method":"config.models.list","requestId":"r-320-models","source":"panel-composer","configProvider":` +
		jsonString320(provider) + `}`
	receipt, err := disp.Handle(context.Background(), env)
	if err != nil {
		t.Fatalf("config.models.list for %q was refused: %v", provider, err)
	}
	idx := strings.Index(receipt, panel.ModelsListPayloadMarker)
	if idx < 0 {
		t.Fatalf("the receipt carries no %q line, so the page cannot parse it: %q", panel.ModelsListPayloadMarker, receipt)
	}
	line := receipt[idx+len(panel.ModelsListPayloadMarker):]
	if i := strings.IndexAny(line, "\r\n"); i >= 0 {
		line = line[:i]
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("the marker line is not JSON (%v): %q", err, line)
	}
	return receipt, payload
}

func jsonString320(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// newModelsLeg320 builds the production pair the fetch leg needs, over the
// ticket-248 settings seed.
func newModelsLeg320(t *testing.T, dir string) (*panel.ComposerDispatch, *config.Manager) {
	t.Helper()
	mgr, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	secrets, err := secret.NewStore(dir)
	if err != nil {
		t.Fatalf("secret.NewStore: %v", err)
	}
	return &panel.ComposerDispatch{Models: &modelsListHandler320{mgr: mgr, secrets: secrets}}, mgr
}

func Test320ModelsListFetchesThroughTheRealLeg(t *testing.T) {
	var hits atomic.Int32
	var sawAuth atomic.Value
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		sawAuth.Store(r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"m-b","owned_by":"vendor-b"},{"id":"m-a","owned_by":"vendor-a"}]}`))
	}))
	defer srv.Close()

	dir := settingsDir248(t)
	disp, mgr := newModelsLeg320(t, dir)
	t.Setenv("DEEPSEEK_KEY", "test-key-320")
	if _, err := mgr.SetProviderBaseURL("deepseek", srv.URL+"/v1"); err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}

	receipt, payload := modelsPath320(t, disp, "deepseek")
	if n := hits.Load(); n != 1 {
		t.Fatalf("the endpoint was called %d time(s) for one fetch, want exactly 1", n)
	}
	if got, _ := sawAuth.Load().(string); got != "Bearer test-key-320" {
		t.Errorf("the request carried Authorization %q, want the RESOLVED key (the ref is never sent)", got)
	}
	if strings.Contains(receipt, "test-key-320") {
		t.Errorf("the receipt echoes the credential; the reply is page-visible and must never carry it: %q", receipt)
	}
	if got, want := payload["provider"], "deepseek"; got != want {
		t.Errorf("payload provider = %v, want %q", got, want)
	}
	models, _ := payload["models"].([]any)
	if len(models) != 2 {
		t.Fatalf("payload carries %d model(s), want 2: %v", len(models), payload)
	}
	first, _ := models[0].(map[string]any)
	if first["id"] != "m-a" {
		t.Errorf("first id = %v, want m-a (DiscoverModels sorts by id)", first["id"])
	}
	if first["ownedBy"] != "vendor-a" {
		t.Errorf("first ownedBy = %v, want vendor-a (the page's field spelling, not internal/llm's)", first["ownedBy"])
	}
}

func Test320ModelsListRefusalsDoNotTouchTheNetwork(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	dir := settingsDir248(t)
	disp, mgr := newModelsLeg320(t, dir)

	// (a) an unknown provider: refused by name, no request.
	_, err := disp.Handle(context.Background(),
		`{"method":"config.models.list","requestId":"r-320-a","source":"panel-composer","configProvider":"ghost"}`)
	if err == nil || !strings.Contains(err.Error(), "ghost") {
		t.Errorf("an unknown provider was not refused by name: %v", err)
	}
	// (b) a known row with no base_url: refused by name, no request.
	if _, err := mgr.SetProviderBaseURL("deepseek", ""); err != nil {
		t.Fatalf("SetProviderBaseURL(\"\"): %v", err)
	}
	_, err = disp.Handle(context.Background(),
		`{"method":"config.models.list","requestId":"r-320-b","source":"panel-composer","configProvider":"deepseek"}`)
	if err == nil || !strings.Contains(err.Error(), "base_url") {
		t.Errorf("a row with no base_url was not refused by name: %v", err)
	}
	// (c) a row pointed at the test server but whose ref does not resolve: refused
	// by name, no request - the Unconfigured posture.
	dir2 := settingsDir248(t)
	disp2, mgr2 := newModelsLeg320(t, dir2)
	if _, err := mgr2.SetProviderBaseURL("deepseek", srv.URL+"/v1"); err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}
	t.Setenv("DEEPSEEK_KEY", "")
	_, err = disp2.Handle(context.Background(),
		`{"method":"config.models.list","requestId":"r-320-c","source":"panel-composer","configProvider":"deepseek"}`)
	if err == nil || !strings.Contains(err.Error(), "密钥") {
		t.Errorf("an unresolved key ref was not refused by name: %v", err)
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("the endpoint was called %d time(s) across three refusals, want 0", n)
	}
}

// Test320ModelsListNeedsAProvider pins the empty-selector refusal, which is the
// shape the page cannot produce but a hand-written envelope can.
func Test320ModelsListNeedsAProvider(t *testing.T) {
	dir := settingsDir248(t)
	disp, _ := newModelsLeg320(t, dir)
	_, err := disp.Handle(context.Background(),
		`{"method":"config.models.list","requestId":"r-320-d","source":"panel-composer"}`)
	if err == nil || !strings.Contains(err.Error(), "configProvider") {
		t.Errorf("a fetch with no provider was not refused by name: %v", err)
	}
}
