package models

// Ticket 121 AC#6 - the cost of the hand-off re-verification, measured rather
// than quoted. acceptor-ticket109b measured 256 MiB through
// Manager.VerifyInstalled at 1.07-1.22 s (220-252 MB/s) while at least three
// other agents were compiling on this host and the page cache could not be
// dropped, so every one of those samples is a warm read. This benchmark exists so
// that number has a reproducible instrument next to it, and so the "startup pays
// for the model bytes twice" statement can be checked by anyone: Ensure's
// cache-hit branch already hashes the whole installed tree, and the hand-off
// point hashes it again.
//
// It is a Benchmark, not a Test, on purpose: `go test` never runs it (neither
// does scripts/portable-tests.sh, which passes no -bench flag), so the tail it
// measures cannot quietly add five seconds to the CI leg, and it can never become
// a SKIP - which that runner treats as fatal. Run it with:
//
//	go test ./internal/models -run '^$' -bench BenchmarkAC6 -benchtime=5x -benchmem
//
// Every sample is reported as its own line on purpose. Averaging is exactly how
// a 1.22 s tail disappears from a 1.15 s "mean", and the number that matters for
// a startup budget is the tail.

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ac6SampleMiB is the size the acceptor's reading used; it is close to the
// largest single model the real manifest pins (asr-streaming is ~237 MB), so the
// per-call numbers below map onto a real install without a unit conversion.
const ac6SampleMiB = 256

// BenchmarkAC6VerifyInstalledHandOffCost times VerifyInstalled over one
// ac6SampleMiB file whose bytes are exactly what the signed manifest pins.
func BenchmarkAC6VerifyInstalledHandOffCost(b *testing.B) {
	body := ac6DeterministicBody(ac6SampleMiB)
	store, mgr := ac6InstalledStore(b, body)

	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		started := time.Now()
		if err := mgr.VerifyInstalled("ac6-fixture"); err != nil {
			b.Fatalf("clean install refused at sample %d: %v", i, err)
		}
		elapsed := time.Since(started)
		b.ReportMetric(elapsed.Seconds(), "s/call")
		// The benchmark framework prints one mean over b.N; the reading this
		// ticket cares about is per sample, so each one is logged verbatim.
		fmt.Fprintf(os.Stderr, "AC6 sample %d: %d MiB in %s (%.0f MB/s), store=%s\n",
			i, len(body)>>20, elapsed.Round(time.Millisecond),
			float64(len(body))/1e6/elapsed.Seconds(), store)
	}
}

// BenchmarkAC6EnsureCacheHitIsTheSameCost is the other half of the "doubles"
// claim: it times the hash Ensure itself already pays on a cache hit over the
// same bytes, so "the hand-off is the same order as what startup already pays"
// is a comparison of two measurements and not a rhetorical move.
func BenchmarkAC6EnsureCacheHitIsTheSameCost(b *testing.B) {
	body := ac6DeterministicBody(ac6SampleMiB)
	_, mgr := ac6InstalledStore(b, body)

	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		started := time.Now()
		if _, err := mgr.Ensure(context.Background(), "ac6-fixture"); err != nil {
			b.Fatalf("cache-hit Ensure failed at sample %d: %v", i, err)
		}
		elapsed := time.Since(started)
		b.ReportMetric(elapsed.Seconds(), "s/call")
		fmt.Fprintf(os.Stderr, "AC6 ensure-cache-hit sample %d: %d MiB in %s (%.0f MB/s)\n",
			i, len(body)>>20, elapsed.Round(time.Millisecond),
			float64(len(body))/1e6/elapsed.Seconds())
	}
}

// ac6DeterministicBody builds MiB mebibytes of reproducible bytes. Deterministic
// so the manifest can pin a sha256 computed here, and patterned at a 1 KiB stride
// rather than all zeros so the hash loop is not measured against the one input
// shape a real model file would never present.
func ac6DeterministicBody(mib int) []byte {
	block := make([]byte, 1024)
	for i := range block {
		block[i] = byte(i*31 + 7)
	}
	out := make([]byte, 0, mib<<20)
	for k := 0; k < (mib<<20)/len(block); k++ {
		out = append(out, block...)
		out[len(out)-1] = byte(k)
	}
	return out
}

// ac6InstalledStore writes body into a fresh store as the single named file of
// the ac6-fixture entry, signs a manifest that pins it, and returns a Manager
// whose cache hit and hand-off both land on those bytes.
func ac6InstalledStore(b *testing.B, body []byte) (store string, mgr *Manager) {
	b.Helper()
	sum := sha256.Sum256(body)
	const name = "ac6.bin"
	mf := ac6SignedManifest(b, fmt.Sprintf(`{
		"manifest_version": 1,
		"models": [
			{
				"id": "ac6-fixture", "purpose": "asr-streaming",
				"urls": ["http://127.0.0.1:1/ac6.bin"],
				"sha256": "%s", "size_bytes": %d,
				"license": "MIT", "quant": "int8"
			}
		]
	}`, hex.EncodeToString(sum[:]), len(body)))

	store = filepath.Join(b.TempDir(), "models")
	dir := filepath.Join(store, "ac6-fixture")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), body, 0o644); err != nil {
		b.Fatal(err)
	}
	// The entry has to name the file it installed, and the loader has to agree
	// the tree is exactly that: if this premise is wrong the timings below would
	// be measuring a refusal, not a verification.
	if !strings.HasSuffix(mf.Models[0].URLs[0], "/"+name) {
		b.Fatalf("fixture url does not end in the installed file name: %s", mf.Models[0].URLs[0])
	}
	var err error
	mgr, err = NewManager(Options{
		DataDir: store, Manifest: mf, VerifySignature: true, Attempts: 1,
		// Any request at all is a bug in this fixture: both timings below are
		// cache-hit measurements, and a silent download would turn the reading
		// into a network benchmark.
		HTTPClient: &http.Client{Transport: panicTransport{}},
	})
	if err != nil {
		b.Fatal(err)
	}
	if err := mgr.VerifyInstalled("ac6-fixture"); err != nil {
		b.Fatalf("fixture does not verify before it is timed: %v", err)
	}
	return store, mgr
}

// ac6SignedManifest signs the given manifest JSON with the same deterministic
// test key the package's other fixtures use, and loads it back through the real
// verification path. It is a local copy rather than a call into
// loadManifestFromString because that helper takes a *testing.T while a benchmark
// has a *testing.B; the key material and the code path are identical, so the
// signature leg measured here is the same one the tests exercise.
func ac6SignedManifest(tb testing.TB, manifest string) *Manifest {
	tb.Helper()
	seed := []byte("wisp-test-minisign-seed-01234567")
	priv := ed25519.NewKeyFromSeed(seed)
	pub := MinisignPublicKey{KeyID: keyIDFromPub(priv.Public().(ed25519.PublicKey)), Key: priv.Public().(ed25519.PublicKey)}

	dir := tb.TempDir()
	mpath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(mpath, []byte(manifest), 0o644); err != nil {
		tb.Fatal(err)
	}
	raw, err := os.ReadFile(mpath)
	if err != nil {
		tb.Fatal(err)
	}
	sig, err := SignPayload(priv, pub, raw, "ac6", "trusted comment: ac6")
	if err != nil {
		tb.Fatal(err)
	}
	if err := os.WriteFile(mpath+".minisig", sig, 0o644); err != nil {
		tb.Fatal(err)
	}
	m, err := LoadSignedManifest(mpath, mpath+".minisig", pub.String("ac6 key"))
	if err != nil {
		tb.Fatalf("LoadSignedManifest: %v", err)
	}
	return m
}
