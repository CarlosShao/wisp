package models

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// ---- test fixtures ---------------------------------------------------------

// tinyModelArchive constants pin internal/models/testdata/tiny-model.tar.bz2
// (deterministic; regenerate with the recipe in testdata/README.md).
const (
	tinyArchiveSHA = "fb043a49e248bba9881648d9362abeeb78e129a2b1fcb4f5c73340317536c54d"
	tinyArchiveLen = 795
)

var tinyArchiveMembers = []ExtractedFile{
	{"model.onnx", "785b0751fc2c53dc14a4ce3d800e69ef9ce1009eb327ccf458afe09c242c26c9", 1024},
	{"tokens.txt", "fa89d000e03baeda19370a433d65fca9665e80f8e76fda996df865662c8aea50", 20},
	{"dict/inner.txt", "489810b7007e1bb4fe29cadd1c0eba2b7b23d8b1ffee00ea7090d93e8668b5be", 19},
	{"README.md", "e97d8cc1824a23548505a48fb1f5a7e2b7d4f931a72f145f6250c58017bb18c8", 21},
}

var vadBody = []byte(strings.Repeat("wisp-vad-fixture-bytes-", 200)) // 4600 bytes

func shaOf(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// manifestJSON is the single test-manifest template. vadURL/kwsURL point at
// the local test servers.
func manifestJSON(vadURL, kwsURL string) string {
	return fmt.Sprintf(`{
		"manifest_version": 1,
		"models": [
			{
				"id": "vad-fixture", "purpose": "vad",
				"urls": ["%s"],
				"sha256": "%s", "size_bytes": %d,
				"license": "MIT", "quant": "fp32"
			},
			{
				"id": "kws-fixture", "purpose": "kws",
				"urls": ["%s"],
				"sha256": "%s", "size_bytes": %d,
				"license": "Apache-2.0", "quant": "int8",
				"files": [{
					"path": "tiny-model.tar.bz2", "urls": ["%s"],
					"sha256": "%s", "size_bytes": %d,
					"archive": {"format": "tar.bz2", "files": [
						{"path": "model.onnx", "sha256": "%s", "size_bytes": 1024},
						{"path": "tokens.txt", "sha256": "%s", "size_bytes": 20},
						{"path": "dict/inner.txt", "sha256": "%s", "size_bytes": 19},
						{"path": "README.md", "sha256": "%s", "size_bytes": 21}
					]}
				}]
			}
		]
	}`,
		vadURL, shaOf(vadBody), len(vadBody),
		kwsURL, tinyArchiveSHA, tinyArchiveLen, kwsURL, tinyArchiveSHA, tinyArchiveLen,
		tinyArchiveMembers[0].SHA256, tinyArchiveMembers[1].SHA256,
		tinyArchiveMembers[2].SHA256, tinyArchiveMembers[3].SHA256)
}

// loadManifestFromString signs json with a deterministic test key and loads it
// back through the real verification path.
func loadManifestFromString(t *testing.T, json string) (*Manifest, MinisignPublicKey) {
	t.Helper()
	seed := []byte("wisp-test-minisign-seed-01234567")
	if len(seed) != ed25519.SeedSize {
		t.Fatalf("test seed length %d", len(seed))
	}
	priv := ed25519.NewKeyFromSeed(seed)
	pub := MinisignPublicKey{KeyID: keyIDFromPub(priv.Public().(ed25519.PublicKey)), Key: priv.Public().(ed25519.PublicKey)}

	dir := t.TempDir()
	mpath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(mpath, []byte(json), 0o644); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(mpath)
	sig, err := SignPayload(priv, pub, raw, "test", "trusted comment: test")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mpath+".minisig", sig, 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadSignedManifest(mpath, mpath+".minisig", pub.String("test key"))
	if err != nil {
		t.Fatalf("LoadSignedManifest: %v", err)
	}
	return m, pub
}

func manifestWithURLs(t *testing.T, base string) *Manifest {
	t.Helper()
	m, _ := loadManifestFromString(t, manifestJSON(base+"/vad-fixture/vad.onnx", base+"/kws-fixture/tiny-model.tar.bz2"))
	return m
}

func signedManifest(t *testing.T) (*Manifest, MinisignPublicKey) {
	t.Helper()
	return loadManifestFromString(t, manifestJSON("http://127.0.0.1:1/vad.onnx", "http://127.0.0.1:1/kws.tar.bz2"))
}

// eventCollector is a thread-safe progress sink (progress fires from the
// Ensure goroutine; Cancel tests read from the test goroutine).
type eventCollector struct {
	mu     sync.Mutex
	events []ProgressEvent
}

func (c *eventCollector) add(ev ProgressEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, ev)
}

func (c *eventCollector) all() []ProgressEvent {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]ProgressEvent(nil), c.events...)
}

func (c *eventCollector) phases() map[Phase]bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := map[Phase]bool{}
	for _, ev := range c.events {
		out[ev.Phase] = true
	}
	return out
}

func newTestManager(t *testing.T, m *Manifest, mirrors []string, override map[string]string) (*Manager, *eventCollector) {
	t.Helper()
	col := &eventCollector{}
	mgr, err := NewManager(Options{
		DataDir:         t.TempDir(),
		Manifest:        m,
		Mirrors:         mirrors,
		VerifySignature: true,
		LocalOverride:   override,
		BackoffBase:     time.Millisecond,
		Attempts:        2,
		Progress:        col.add,
	})
	if err != nil {
		t.Fatal(err)
	}
	return mgr, col
}

// countingServer serves static files, counts hits, records Range headers, and
// can 404 everything or truncate the body mid-transfer (simulated kill).
type countingServer struct {
	srv    *httptest.Server
	mu     sync.Mutex
	hits   int
	ranges []string
	bodies map[string][]byte
	partial int // >0: serve only this many bytes when no Range header, then drop
	notFound bool
}

func newCountingServer(t *testing.T, bodies map[string][]byte) *countingServer {
	t.Helper()
	cs := &countingServer{bodies: bodies}
	cs.srv = httptest.NewServer(http.HandlerFunc(cs.handler))
	t.Cleanup(cs.srv.Close)
	return cs
}

func (cs *countingServer) handler(w http.ResponseWriter, r *http.Request) {
	cs.mu.Lock()
	cs.hits++
	cs.ranges = append(cs.ranges, r.Header.Get("Range"))
	notFound, partial := cs.notFound, cs.partial
	cs.mu.Unlock()

	body, ok := cs.bodies[r.URL.Path]
	if !ok || notFound {
		http.NotFound(w, r)
		return
	}
	if rng := r.Header.Get("Range"); rng != "" {
		var start int64
		if _, err := fmt.Sscanf(rng, "bytes=%d-", &start); err != nil || start < 0 || start > int64(len(body)) {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, int64(len(body))-1, len(body)))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", int64(len(body))-start))
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(body[start:])
		return
	}
	if partial > 0 && partial < len(body) {
		// Declared full length, deliver a prefix, then return: the client sees
		// an unexpected EOF, exactly like a killed downloader.
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body[:partial])
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	_, _ = w.Write(body)
}

func (cs *countingServer) URL() string { return cs.srv.URL }
func (cs *countingServer) Hits() int   { cs.mu.Lock(); defer cs.mu.Unlock(); return cs.hits }
func (cs *countingServer) LastRange() string {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if n := len(cs.ranges); n > 0 {
		return cs.ranges[n-1]
	}
	return ""
}
func (cs *countingServer) setNotFound(v bool) { cs.mu.Lock(); cs.notFound = v; cs.mu.Unlock() }
func (cs *countingServer) setPartial(n int)   { cs.mu.Lock(); cs.partial = n; cs.mu.Unlock() }

// ---- acceptance tests -------------------------------------------------------

func TestTamperedModelByteRejectedAndDeleted(t *testing.T) {
	bad := append([]byte{}, vadBody...)
	bad[100] ^= 0x01
	srv := newCountingServer(t, map[string][]byte{
		"/vad-fixture/vad.onnx": bad,
	})
	m := manifestWithURLs(t, srv.URL())
	mgr, _ := newTestManager(t, m, nil, nil)

	_, err := mgr.Ensure(context.Background(), "vad-fixture")
	if err == nil {
		t.Fatal("tampered download accepted (must fail)")
	}
	if !strings.Contains(err.Error(), "sha256 MISMATCH") {
		t.Fatalf("wrong error: %v", err)
	}
	// reject+delete: nothing left behind.
	if _, statErr := os.Stat(filepath.Join(mgr.opts.DataDir, "vad-fixture")); !os.IsNotExist(statErr) {
		t.Fatal("install dir must not exist after integrity failure")
	}
	if _, statErr := os.Stat(filepath.Join(mgr.opts.DataDir, "staging", "vad-fixture")); !os.IsNotExist(statErr) {
		t.Fatal("staging must be cleaned after integrity failure")
	}
}

func TestTamperedManifestFailsBeforeAnyDownload(t *testing.T) {
	seed := []byte("wisp-test-minisign-seed-01234567")
	priv := ed25519.NewKeyFromSeed(seed)
	pub := MinisignPublicKey{KeyID: keyIDFromPub(priv.Public().(ed25519.PublicKey)), Key: priv.Public().(ed25519.PublicKey)}
	dir := t.TempDir()
	mpath := filepath.Join(dir, "manifest.json")
	good := manifestJSON("http://127.0.0.1:1/vad.onnx", "http://127.0.0.1:1/kws.tar.bz2")
	if err := os.WriteFile(mpath, []byte(good), 0o644); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(mpath)
	sig, err := SignPayload(priv, pub, raw, "t", "trusted comment: t")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mpath+".minisig", sig, 0o644); err != nil {
		t.Fatal(err)
	}

	// Tamper: flip one byte of the signed manifest.
	tampered := []byte(good)
	tampered[8] ^= 0x01
	if err := os.WriteFile(mpath, tampered, 0o644); err != nil {
		t.Fatal(err)
	}

	srv := newCountingServer(t, map[string][]byte{})
	_, err = LoadSignedManifest(mpath, mpath+".minisig", pub.String("test key"))
	if err == nil {
		t.Fatal("tampered manifest verified (must fail)")
	}
	if !strings.Contains(err.Error(), "FAILED") {
		t.Fatalf("expected signature failure, got: %v", err)
	}
	if srv.Hits() != 0 {
		t.Fatalf("tampered manifest reached the network (%d hits)", srv.Hits())
	}
}

func TestVerifySignatureFalseIsHardError(t *testing.T) {
	m, _ := signedManifest(t)
	_, err := NewManager(Options{DataDir: t.TempDir(), Manifest: m, VerifySignature: false})
	if err == nil {
		t.Fatal("VerifySignature=false accepted (C29 violation)")
	}
	if !strings.Contains(err.Error(), "NOT disableable") {
		t.Fatalf("wrong error: %v", err)
	}
}

func TestResumeContinuesAtByteOffset(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/vad-fixture/vad.onnx": vadBody,
	})
	srv.setPartial(2_000)
	m := manifestWithURLs(t, srv.URL())
	col := &eventCollector{}
	// Attempts:1 -> the first Ensure dies mid-transfer exactly like a killed
	// process; only the SECOND Ensure (restart) can resume.
	mgr, err := NewManager(Options{
		DataDir: t.TempDir(), Manifest: m, VerifySignature: true,
		BackoffBase: time.Millisecond, Attempts: 1, Progress: col.add,
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = col

	_, err = mgr.Ensure(context.Background(), "vad-fixture")
	if err == nil {
		t.Fatal("mid-transfer kill must fail the first run")
	}
	staged := filepath.Join(mgr.opts.DataDir, "staging", "vad-fixture", "vad.onnx")
	st, statErr := os.Stat(staged)
	if statErr != nil {
		t.Fatalf("partial must survive a mid-transfer kill for resume: %v", err)
	}
	if st.Size() != 2_000 {
		t.Fatalf("staged %d bytes, want 2000", st.Size())
	}

	// Restart: the downloader must resume at the exact byte offset.
	srv.setPartial(0)
	dir, err := mgr.Ensure(context.Background(), "vad-fixture")
	if err != nil {
		t.Fatalf("resume run failed: %v", err)
	}
	if got := srv.LastRange(); got != "bytes=2000-" {
		t.Fatalf("resume Range header %q, want bytes=2000-", got)
	}
	out, readErr := os.ReadFile(filepath.Join(dir, "vad.onnx"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if shaOf(out) != shaOf(vadBody) {
		t.Fatal("resumed file hash mismatch")
	}
	if _, statErr := os.Stat(staged); !os.IsNotExist(statErr) {
		t.Fatal("staging must be cleaned after success")
	}
}

func TestMirrorFailoverAndProgressEvents(t *testing.T) {
	primary := newCountingServer(t, map[string][]byte{})
	primary.setNotFound(true) // primary mirror: 404
	fallback := newCountingServer(t, map[string][]byte{
		"/vad-fixture/vad.onnx": vadBody,
	})
	m := manifestWithURLs(t, fallback.URL())
	mgr, col := newTestManager(t, m, []string{primary.URL()}, nil)

	dir, err := mgr.Ensure(context.Background(), "vad-fixture")
	if err != nil {
		t.Fatalf("failover failed: %v", err)
	}
	if primary.Hits() == 0 {
		t.Fatal("primary mirror was never tried")
	}
	if fallback.Hits() == 0 {
		t.Fatal("fallback was never used")
	}
	if out, readErr := os.ReadFile(filepath.Join(dir, "vad.onnx")); readErr != nil || shaOf(out) != shaOf(vadBody) {
		t.Fatal("installed content wrong after failover")
	}
	phases := col.phases()
	for _, want := range []Phase{PhaseConnecting, PhaseDownloading, PhaseVerifying, PhaseInstalling, PhaseDone} {
		if !phases[want] {
			t.Fatalf("progress phase %q missing; got %v", want, phases)
		}
	}
}

func TestCancelCleansStagingNoPartialsOutside(t *testing.T) {
	release := make(chan struct{})
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(vadBody)))
		_, _ = w.Write(vadBody[:4096])
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	defer slow.Close()
	defer close(release)

	m := manifestWithURLs(t, slow.URL)
	mgr, col := newTestManager(t, m, nil, nil)

	done := make(chan error, 1)
	go func() {
		_, err := mgr.Ensure(context.Background(), "vad-fixture")
		done <- err
	}()
	waitForPhase(t, col, PhaseDownloading)
	mgr.Cancel("vad-fixture")

	err := <-done
	if err == nil {
		t.Fatal("cancelled Ensure returned success")
	}
	if !strings.Contains(err.Error(), "cancel") {
		t.Fatalf("wrong error: %v", err)
	}
	// No partials outside staging (staging itself must be gone too).
	walkErr := filepath.Walk(mgr.opts.DataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(filepath.ToSlash(path), "/staging/") {
			t.Errorf("leftover staging file: %s", path)
		}
		return nil
	})
	if walkErr != nil {
		t.Fatal(walkErr)
	}
}

func TestLocalOverrideIntegrityNoNetwork(t *testing.T) {
	m, _ := signedManifest(t)
	local := t.TempDir()
	if err := os.WriteFile(filepath.Join(local, "vad.onnx"), vadBody, 0o644); err != nil {
		t.Fatal(err)
	}

	// Transport that fails the test on ANY request: override must be offline.
	blocking := http.Client{Transport: panicTransport{}}
	mgr, err := NewManager(Options{
		DataDir: t.TempDir(), Manifest: m, VerifySignature: true,
		LocalOverride: map[string]string{"vad-fixture": local},
		HTTPClient:    &blocking,
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := mgr.Ensure(context.Background(), "vad-fixture")
	if err != nil {
		t.Fatalf("local_override failed: %v", err)
	}
	if got != local {
		t.Fatalf("override dir %s, want %s", got, local)
	}

	// Tampered override: integrity must still bite.
	bad := append([]byte{}, vadBody...)
	bad[7] ^= 0x01
	if err := os.WriteFile(filepath.Join(local, "vad.onnx"), bad, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := mgr.Ensure(context.Background(), "vad-fixture"); err == nil {
		t.Fatal("tampered local_override accepted (integrity check bypassed)")
	}
}

type panicTransport struct{}

func (panicTransport) RoundTrip(*http.Request) (*http.Response, error) {
	panic("network call made during local_override (must never happen)")
}

func TestCacheHitReverifiedEveryTime(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/vad-fixture/vad.onnx": vadBody,
	})
	m := manifestWithURLs(t, srv.URL())
	mgr, _ := newTestManager(t, m, nil, nil)

	if _, err := mgr.Ensure(context.Background(), "vad-fixture"); err != nil {
		t.Fatal(err)
	}
	hits := srv.Hits()
	if hits == 0 {
		t.Fatal("server never hit")
	}
	// Second Ensure: cache hit, zero extra requests.
	if _, err := mgr.Ensure(context.Background(), "vad-fixture"); err != nil {
		t.Fatal(err)
	}
	if srv.Hits() != hits {
		t.Fatalf("cache hit made network calls (%d -> %d)", hits, srv.Hits())
	}
	// Flip one byte of the INSTALLED model: the cache check must reject it,
	// then heal by re-downloading VERIFIED bytes (never load tampered ones).
	p := filepath.Join(mgr.opts.DataDir, "vad-fixture", "vad.onnx")
	raw, _ := os.ReadFile(p)
	raw[10] ^= 0x01
	_ = os.WriteFile(p, raw, 0o644)
	if _, err := mgr.Ensure(context.Background(), "vad-fixture"); err != nil {
		t.Fatalf("tampered cache must be detected and healed: %v", err)
	}
	fixed, _ := os.ReadFile(p)
	if shaOf(fixed) != shaOf(vadBody) {
		t.Fatal("tampered installed model was not replaced with verified bytes")
	}
	if srv.Hits() <= hits {
		t.Fatal("healing did not re-download from the mirror")
	}
}

func TestArchiveModelExtractsFlattensInstalls(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	mgr, _ := newTestManager(t, m, nil, nil)

	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("archive install failed: %v", err)
	}
	for _, want := range tinyArchiveMembers {
		p := filepath.Join(dir, filepath.FromSlash(want.Path))
		if err := verifyFileHash(p, want.SHA256, want.SizeBytes); err != nil {
			t.Errorf("extracted %s: %v", want.Path, err)
		}
	}
	// Flattening: no tiny-model-2024-01-01 top dir left.
	if _, err := os.Stat(filepath.Join(dir, "tiny-model-2024-01-01")); !os.IsNotExist(err) {
		t.Error("archive top dir was not flattened")
	}
	// The archive itself must not sit next to the extracted files.
	if _, err := os.Stat(filepath.Join(dir, "tiny-model.tar.bz2")); !os.IsNotExist(err) {
		t.Error("archive file leaked into install dir")
	}
	// Cache hit without network.
	hits := srv.Hits()
	if _, err := mgr.Ensure(context.Background(), "kws-fixture"); err != nil {
		t.Fatal(err)
	}
	if srv.Hits() != hits {
		t.Fatal("second Ensure hit the network")
	}
}

// ---- helpers ---------------------------------------------------------------

func mustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func waitForPhase(t *testing.T, col *eventCollector, phase Phase) {
	t.Helper()
	for i := 0; i < 500; i++ {
		for _, ev := range col.all() {
			if ev.Phase == phase {
				return
			}
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for phase %q; events: %v", phase, col.all())
}
