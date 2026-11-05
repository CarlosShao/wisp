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
// THE INSTRUMENT TEETH (leg 320-r3, ticket 320 AC#4's remaining gaps; widened by
// ticket 323 AC#1/AC#2): the fetch port answers separate questions and each used
// to have no world that could prove it - "the key never lands on disk"
// (Test320FetchLegWritesTheKeyNowhereOnDisk, now with the file-set + per-file-sha
// ruler and the "no new file carries the credential" ruler beside the substring
// search), "a slow provider is cut off by the budget"
// (Test320FetchLegTimesOutWithinItsBudget), "the shipped budget is the ten-second
// budget the page can survive" (Test320FetchLegBudgetIsTheOneThePageWaitsOn),
// "a non-200 is classified without relaying the upstream body"
// (Test320FetchLegClassifiesNon200WithoutRelayingTheBody), and "the receipt
// carries the ids and neither the upstream text nor the key"
// (Test320FetchLegCarriesIDsAndNoUpstreamTextOrKey). Every one of them is
// red on the head this leg started from or under the mutation named in
// docs/evidence/s1/320-r3-instrument.md and docs/evidence/s1/323-a0-ac1-ac2.md; a
// tooth whose red cannot be produced is a claim, not an instrument.

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
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

// dirSnapshot320 is one reading of a directory: every regular file's path
// relative to the root, and the sha256 of its bytes. Ticket 323 AC#1 names this
// ruler for a reason the substring detector above cannot cover: a fetch leg can
// change the on-disk state without ever writing the credential's bytes (it can
// rewrite config.toml, append to a secret blob, or drop a scratch file holding
// the key re-encoded), and "the file set and every file's sha are the same
// before and after" is the shape that sees all of those. A fetch is a read-only
// hop: it reads config.toml, it reads one secret, it writes nothing.
type dirSnapshot320 map[string][sha256.Size]byte

// snapDir320 walks root and returns the snapshot. Paths are relative to root so
// a temp directory's random name never becomes part of a comparison.
func snapDir320(t *testing.T, root string) dirSnapshot320 {
	t.Helper()
	snap := dirSnapshot320{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		rel, rerr := filepath.Rel(root, path)
		if rerr != nil {
			return rerr
		}
		snap[rel] = sha256.Sum256(b)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	return snap
}

// diffSnapshot320 names every difference in BOTH directions - a file that
// disappeared and a file that appeared are the same class of event for this
// ruler - and returns nil when the two readings are identical.
func diffSnapshot320(before, after dirSnapshot320) []string {
	var out []string
	for name, sum := range before {
		got, ok := after[name]
		if !ok {
			out = append(out, name+" (disappeared)")
			continue
		}
		if got != sum {
			out = append(out, name+" (bytes changed)")
		}
	}
	for name := range after {
		if _, ok := before[name]; !ok {
			out = append(out, name+" (appeared)")
		}
	}
	sort.Strings(out)
	return out
}

// credentialsInNewFiles320 answers AC#1's second half literally: after the
// fetch, which files EXIST that did not exist before and carry the credential in
// their name or bytes? config.toml is excluded because the AC excludes it; it is
// not a loophole, because the snapshot ruler above already refuses any change to
// config.toml at all.
func credentialsInNewFiles320(t *testing.T, root string, before dirSnapshot320, needle string) []string {
	t.Helper()
	var hits []string
	for name := range snapDir320(t, root) {
		if _, existed := before[name]; existed {
			continue
		}
		if name == configFileName {
			continue
		}
		path := filepath.Join(root, name)
		if strings.Contains(name, needle) {
			hits = append(hits, path+" (file name)")
			continue
		}
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if strings.Contains(string(b), needle) {
			hits = append(hits, path)
		}
	}
	sort.Strings(hits)
	return hits
}

// Test320FetchLegWritesTheKeyNowhereOnDisk is ticket 320 AC#4's a4m5 tooth.
// Ticket 323 AC#1 widened it from a substring search to two rulers: the file set
// with every file's sha (identical before and after the fetch) and "no file that
// appeared carries the credential". The verdict of 320-v1 was that the handler could write the resolved key to a file
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
	// Snapshot AFTER the fixture's own writes (the seed, the secret store, the
	// base_url) so the ruler measures the fetch leg and nothing else.
	before := snapDir320(t, dir)

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

	// AC#1 ruler one: the file set and every file's sha are the same before and
	// after. This is the half the substring search cannot see - a leg that
	// rewrites config.toml or drops a file with the key re-encoded leaves no
	// plaintext behind and would pass a bytes-only search.
	if diff := diffSnapshot320(before, snapDir320(t, dir)); len(diff) > 0 {
		t.Errorf("the fetch leg changed the config root; a fetch is a read-only hop (reads config.toml, reads one secret, writes nothing): %v", diff)
	}
	// AC#1 ruler two, named literally: no NEW file other than config.toml carries
	// the credential. config.toml is not a loophole - ruler one already refuses
	// any change to it.
	if hits := credentialsInNewFiles320(t, dir, before, key); len(hits) > 0 {
		t.Errorf("a file appeared under the config root carrying the resolved credential: %v", hits)
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

	// The honest reading is taken here, BEFORE any control file exists, so the
	// two controls below can be read as their own diffs and a leak above cannot
	// contaminate them (a control that folds the leak into its own count reports
	// "the ruler is blind" when the ruler in fact saw the leak).
	afterFetch := snapDir320(t, dir)

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
	// Same positive control for the snapshot ruler, in both directions: a new
	// file must appear in the diff, and a rewritten file must show as changed.
	// A ruler that cannot see a planted difference is a claim, not an instrument.
	afterPlant := snapDir320(t, dir)
	if diff := diffSnapshot320(afterFetch, afterPlant); len(diff) != 1 ||
		!strings.Contains(diff[0], "planted-key-320.txt") || !strings.Contains(diff[0], "appeared") {
		t.Fatalf("the snapshot ruler cannot see a file planted at %s (diff = %v); it proves nothing about a rewrite",
			planted, diff)
	}
	appended := filepath.Join(dir, configFileName)
	if werr := os.WriteFile(appended, []byte(settingsSeed+"\n# planted\n"), 0o600); werr != nil {
		t.Fatalf("rewrite the control file: %v", werr)
	}
	if diff := diffSnapshot320(afterPlant, snapDir320(t, dir)); len(diff) != 1 ||
		!strings.Contains(diff[0], configFileName) || !strings.Contains(diff[0], "bytes changed") {
		t.Fatalf("the snapshot ruler cannot see a rewritten %s (diff = %v); it proves nothing about a rewrite",
			configFileName, diff)
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
	var attempted atomic.Int32
	var answered atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempted.Add(1)
		select {
		case <-r.Context().Done():
		case <-release:
			// Only a server that got as far as producing an answer counts here:
			// ticket 323 AC#2 asks for "upstream hits == 0", which is the world
			// where no upstream ANSWER landed, not a world where the request was
			// never sent (a fixture nobody calls cannot prove a timeout).
			answered.Add(1)
		case <-time.After(5 * time.Second):
			// A production leg that ignored its budget would otherwise hang this
			// test into the package timeout; a bounded hold keeps the red clean.
			answered.Add(1)
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
	if n := attempted.Load(); n != 1 {
		t.Errorf("the slow endpoint was called %d time(s), want exactly 1 (no retry storm behind one fetch)", n)
	}
	if n := answered.Load(); n != 0 {
		t.Errorf("the upstream produced %d answer(s) for a request the budget cut off, want 0", n)
	}
}

// Test320FetchLegBudgetIsTheOneThePageWaitsOn pins the PRODUCTION budget value,
// which the timeout tooth above cannot see: that tooth overrides the variable to
// 120ms and restores it, so turning the shipped 8s into 8h leaves it green (320-v1's
// a4m3, re-measured on this tree in docs/evidence/s1/323-a0-ac1-ac2.md: rc=0,
// 12/12 PASS). The pair together is the instrument - one tooth proves the leg
// honours the variable, this one proves the variable is a budget the page can
// survive, and A972 ruling 3 is the number it is pinned to.
func Test320FetchLegBudgetIsTheOneThePageWaitsOn(t *testing.T) {
	// A972 ruling 3: the one synchronous call this page-blocking hop makes has to
	// stay at or under ten seconds.
	const a972Ceiling = 10 * time.Second
	if modelsFetchBudget320 <= 0 {
		t.Fatalf("the fetch budget is %v; a non-positive budget would refuse every fetch", modelsFetchBudget320)
	}
	if modelsFetchBudget320 > a972Ceiling {
		t.Errorf("the fetch budget is %v; A972 ruling 3 caps this page-blocking call at %v, and a budget nobody can reach in a test is a budget nobody watches",
			modelsFetchBudget320, a972Ceiling)
	}
}

// Test320FetchLegClassifiesNon200WithoutRelayingTheBody is the classification
// tooth: 320-v1 found that no test ever made the provider answer non-200, so the
// class could fall out of the sentence unnoticed. The marker is the upstream's
// own words, and the point is the pair - the class IS in the page-visible
// sentence, the upstream body is NOT.
func Test320FetchLegClassifiesNon200WithoutRelayingTheBody(t *testing.T) {
	const marker = "UPSTREAM_LEAK_320XYZ"
	// Ticket 323 AC#3 names two shapes, not one: a credential-shaped refusal
	// (401) and a provider-shaped one (500). They have to land in DIFFERENT D37
	// classes, because the page's next step differs (fix the key vs wait/retry),
	// and 320-v1 found no test at this hop ever made the upstream answer non-200
	// at all - so neither class was reachable from a page-visible sentence.
	cases := []struct {
		name   string
		status int
		class  observe.ErrorClass
	}{
		{"401", http.StatusUnauthorized, observe.ClassAuth},
		{"500", http.StatusInternalServerError, observe.ClassProvider},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(`{"error":{"message":"` + marker + `: the provider's own words"}}`))
			}))
			defer srv.Close()

			dir := settingsDir248(t)
			disp, mgr := newModelsLeg320(t, dir)
			t.Setenv("DEEPSEEK_KEY", uniqueKey320(t))
			if _, err := mgr.SetProviderBaseURL("deepseek", srv.URL+"/v1"); err != nil {
				t.Fatalf("SetProviderBaseURL: %v", err)
			}

			refused, err := disp.Handle(context.Background(), modelsEnvelope320("r-320-"+tc.name, "deepseek"))
			if err == nil {
				t.Fatalf("a %d answer was accepted: %q", tc.status, refused)
			}
			if !strings.Contains(refused, "拉取失败") {
				t.Errorf("the refusal does not say the fetch failed: %q", refused)
			}
			if !strings.Contains(refused, "deepseek") {
				t.Errorf("the refusal does not name the provider: %q", refused)
			}
			if !strings.Contains(refused, string(tc.class)) {
				t.Errorf("the refusal does not carry the D37 class %q for a %d (ruling 4 wants the class, not the body): %q",
					tc.class, tc.status, refused)
			}
			if strings.Contains(refused, marker) {
				t.Errorf("the page-visible refusal relays the upstream body's own words: %q", refused)
			}
		})
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
