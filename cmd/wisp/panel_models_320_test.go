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
//
// THE FOUR INSTRUMENT TEETH (leg 320-r3, ticket 320 AC#4's remaining gaps): the fetch port
// answers four separate questions and each used to have no world that could
// prove it - "the key never lands on disk" (Test320FetchLegWritesTheKeyNowhereOnDisk),
// "a slow provider is cut off by the budget" (Test320FetchLegTimesOutWithinItsBudget),
// "a non-200 is classified without relaying the upstream body"
// (Test320FetchLegClassifiesNon200WithoutRelayingTheBody), and "the receipt
// carries the ids and neither the upstream text nor the key"
// (Test320FetchLegCarriesIDsAndNoUpstreamTextOrKey). Every one of them is
// red on the head this leg started from or under the mutation named in
// docs/evidence/s1/320-r3-instrument.md; a tooth whose red cannot be produced is
// a claim, not an instrument.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
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

// modelsEnvelope320 is the page's envelope for one fetch, built the same way for
// the accepted and the refused cases so no test depends on string surgery.
func modelsEnvelope320(requestID, provider string) string {
	return `{"method":"config.models.list","requestId":` + jsonString320(requestID) +
		`,"source":"panel-composer","configProvider":` + jsonString320(provider) + `}`
}

// uniqueKey320 returns a credential that exists ONLY in this process. The
// no-plaintext tooth below walks the config root and the working tree for these
// bytes, so the value cannot be a literal any file already holds - a literal in
// this very file would make the scan fail on itself and the instrument would have
// to be weakened to a name check, which is exactly the tooth it replaces.
func uniqueKey320(t *testing.T) string {
	t.Helper()
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("crypto/rand: %v", err)
	}
	return "sk-320-" + hex.EncodeToString(b[:])
}

// filesCarrying320 returns every regular file under root whose name or bytes
// contain needle. Nothing here knows HOW a credential would have been written, so
// a dump next to the secret store and a relative-path dump in the working tree
// are both caught: the detector is the filesystem, not a list of expected paths.
func filesCarrying320(t *testing.T, root, needle string) []string {
	t.Helper()
	var hits []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.Contains(filepath.Base(path), needle) {
			hits = append(hits, path+" (file name)")
		}
		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		if strings.Contains(string(b), needle) {
			hits = append(hits, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return hits
}

// Test320FetchLegWritesTheKeyNowhereOnDisk is ticket 320 AC#4's a4m5 tooth. The
// verdict of 320-v1 was that the handler could write the resolved key to a file
// and the whole suite stayed green, so this test makes the filesystem the judge:
// one successful fetch and one failed fetch run, then every file under the
// config data root and under the working tree is read and the key's own bytes are
// looked for. It also pins the receipt side (page-visible text carries no key) and
// closes with a positive control, because a detector that cannot see a key it was
// handed is not an instrument.
func Test320FetchLegWritesTheKeyNowhereOnDisk(t *testing.T) {
	key := uniqueKey320(t)
	var fail atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"message":"upstream refused"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"m-a"}]}`))
	}))
	defer srv.Close()

	dir := settingsDir248(t)
	disp, mgr := newModelsLeg320(t, dir)
	t.Setenv("DEEPSEEK_KEY", key)
	if _, err := mgr.SetProviderBaseURL("deepseek", srv.URL+"/v1"); err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}

	receipt, _ := modelsPath320(t, disp, "deepseek")
	if strings.Contains(receipt, key) {
		t.Errorf("the successful fetch's receipt carries the resolved key: %q", receipt)
	}

	fail.Store(true)
	refused, err := disp.Handle(context.Background(), modelsEnvelope320("r-320-disk", "deepseek"))
	if err == nil {
		t.Fatalf("the 500 was accepted; a refused fetch must never look like a success: %q", refused)
	}
	if strings.Contains(refused, key) || strings.Contains(err.Error(), key) {
		t.Errorf("the refused fetch's page-visible sentence carries the resolved key: %q", refused)
	}

	wd, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatalf("os.Getwd: %v", wdErr)
	}
	for _, root := range []string{dir, wd} {
		if hits := filesCarrying320(t, root, key); len(hits) > 0 {
			t.Errorf("the fetch leg left the resolved credential in %d file(s) under %s: %v",
				len(hits), root, hits)
		}
	}

	// Positive control: the detector must see a key planted on disk.
	planted := filepath.Join(dir, "planted-key-320.txt")
	if werr := os.WriteFile(planted, []byte(key), 0o600); werr != nil {
		t.Fatalf("plant the control key: %v", werr)
	}
	hits := filesCarrying320(t, dir, key)
	if len(hits) != 1 || hits[0] != planted {
		t.Fatalf("the detector cannot see a key planted at %s (hits = %v); it proves nothing about a leak",
			planted, hits)
	}
}

// Test320FetchLegTimesOutWithinItsBudget is ticket 320 AC#4's a4m3 tooth. 320-v1
// measured that the budget could be turned into eight hours and nothing failed,
// because no fixture ever answered slowly. Here the test server holds its answer
// until the request is cancelled, the budget is shrunk for this run only, and the
// refusal has to come back classified as a network deadline instead of the test
// hanging or the call returning a success.
func Test320FetchLegTimesOutWithinItsBudget(t *testing.T) {
	old := modelsFetchBudget320
	modelsFetchBudget320 = 120 * time.Millisecond
	t.Cleanup(func() { modelsFetchBudget320 = old })

	release := make(chan struct{})
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	defer srv.Close()
	defer close(release)

	dir := settingsDir248(t)
	disp, mgr := newModelsLeg320(t, dir)
	t.Setenv("DEEPSEEK_KEY", uniqueKey320(t))
	if _, err := mgr.SetProviderBaseURL("deepseek", srv.URL+"/v1"); err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}

	start := time.Now()
	refused, err := disp.Handle(context.Background(), modelsEnvelope320("r-320-timeout", "deepseek"))
	elapsed := time.Since(start)
	if err == nil {
		t.Fatalf("a provider slower than the 120ms budget was not refused: %q", refused)
	}
	if !strings.Contains(refused, "拉取失败") {
		t.Errorf("the timeout refusal does not say the fetch failed: %q", refused)
	}
	if !strings.Contains(refused, string(observe.ClassNetwork)) {
		t.Errorf("the timeout refusal does not carry the D37 network class: %q", refused)
	}
	if elapsed > 2*time.Second {
		t.Errorf("the budget did not cut the call short: it took %v for a 120ms budget", elapsed)
	}
	if n := hits.Load(); n != 1 {
		t.Errorf("the slow endpoint was called %d time(s), want exactly 1", n)
	}
}

// Test320FetchLegClassifiesNon200WithoutRelayingTheBody is the classification
// tooth: 320-v1 found that no test ever made the provider answer non-200, so the
// class could fall out of the sentence unnoticed. The marker is the upstream's
// own words, and the point is the pair - the class IS in the page-visible
// sentence, the upstream body is NOT.
func Test320FetchLegClassifiesNon200WithoutRelayingTheBody(t *testing.T) {
	const marker = "UPSTREAM_LEAK_320XYZ"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"message":"` + marker + `: the provider's own words"}}`))
	}))
	defer srv.Close()

	dir := settingsDir248(t)
	disp, mgr := newModelsLeg320(t, dir)
	t.Setenv("DEEPSEEK_KEY", uniqueKey320(t))
	if _, err := mgr.SetProviderBaseURL("deepseek", srv.URL+"/v1"); err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}

	refused, err := disp.Handle(context.Background(), modelsEnvelope320("r-320-500", "deepseek"))
	if err == nil {
		t.Fatalf("a 500 answer was accepted: %q", refused)
	}
	if !strings.Contains(refused, "拉取失败") {
		t.Errorf("the refusal does not say the fetch failed: %q", refused)
	}
	if !strings.Contains(refused, "deepseek") {
		t.Errorf("the refusal does not name the provider: %q", refused)
	}
	if !strings.Contains(refused, string(observe.ClassProvider)) {
		t.Errorf("the refusal does not carry the D37 class for a 5xx (ruling 4 wants the class, not the body): %q", refused)
	}
	if strings.Contains(refused, marker) {
		t.Errorf("the page-visible refusal relays the upstream body's own words: %q", refused)
	}
}

// Test320FetchLegCarriesIDsAndNoUpstreamTextOrKey is the remaining pair of AC#4's
// claims in one round trip: on success the receipt carries the discovered ids
// (that is what the page renders), and on failure neither the upstream's text nor
// the resolved key appears in anything the page can read. The same test server
// serves both halves so the two receipts are produced by one leg, not two.
func Test320FetchLegCarriesIDsAndNoUpstreamTextOrKey(t *testing.T) {
	const marker = "UPSTREAM_LEAK_320XYZ"
	key := uniqueKey320(t)
	var fail atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if fail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"message":"` + marker + `: the provider's own words"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"m-b","owned_by":"vendor-b"},{"id":"m-a","owned_by":"vendor-a"}]}`))
	}))
	defer srv.Close()

	dir := settingsDir248(t)
	disp, mgr := newModelsLeg320(t, dir)
	t.Setenv("DEEPSEEK_KEY", key)
	if _, err := mgr.SetProviderBaseURL("deepseek", srv.URL+"/v1"); err != nil {
		t.Fatalf("SetProviderBaseURL: %v", err)
	}

	receipt, payload := modelsPath320(t, disp, "deepseek")
	for _, id := range []string{"m-a", "m-b"} {
		if !strings.Contains(receipt, `"`+id+`"`) {
			t.Errorf("the receipt does not carry the discovered id %q, so the page has nothing to list: %q", id, receipt)
		}
	}
	if models, _ := payload["models"].([]any); len(models) != 2 {
		t.Errorf("the marker line carries %d model(s), want 2: %v", len(models), payload)
	}
	if strings.Contains(receipt, key) {
		t.Errorf("the success receipt carries the resolved key: %q", receipt)
	}

	fail.Store(true)
	refused, err := disp.Handle(context.Background(), modelsEnvelope320("r-320-echo", "deepseek"))
	if err == nil {
		t.Fatalf("the 500 was accepted: %q", refused)
	}
	for _, needle := range []string{marker, key} {
		if strings.Contains(refused, needle) {
			t.Errorf("the page-visible refusal carries %q: %q", needle, refused)
		}
	}
}
