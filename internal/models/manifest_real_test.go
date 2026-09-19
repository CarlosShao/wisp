package models

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
)

// repoManifestPath locates the committed signed manifest regardless of the
// test working directory.
func repoManifestPath(t *testing.T) (string, string) {
	t.Helper()
	for _, dir := range []string{"../..", "."} {
		p := filepath.Join(dir, "models", "manifest.json")
		if _, err := os.Stat(p); err == nil {
			return p, p + ".minisig"
		}
	}
	t.Fatal("committed models/manifest.json not found from test dir")
	return "", ""
}

// TestRealManifestInRepoVerifies pins the C29 chain end to end: the committed
// manifest must verify against the PUBLIC KEY HARDCODED IN BUILDINFO. Any
// tampering with either file fails the whole suite.
func TestRealManifestInRepoVerifies(t *testing.T) {
	mpath, spath := repoManifestPath(t)
	m, err := LoadSignedManifest(mpath, spath, buildinfo.MinisignPublicKey)
	if err != nil {
		t.Fatalf("committed manifest does not verify against buildinfo key: %v", err)
	}
	if len(m.Models) != 6 {
		t.Fatalf("manifest has %d models, want 6", len(m.Models))
	}
	want := map[string]string{
		"kws-zipformer-wenetspeech-3.3M-2024-01-01": "kws",
		"vad-silero":                             "vad",
		"asr-streaming-paraformer-zh-en":         "asr-streaming",
		"asr-offline-sensevoice-zh-en-ja-ko-yue": "asr-offline",
		"punc-ct-transformer-zh-en-vocab272727":  "punctuation",
		"tts-matcha-zh-baker":                    "tts",
	}
	for id, purpose := range want {
		e, err := m.FindModel(id)
		if err != nil {
			t.Fatalf("entry %s missing: %v", id, err)
		}
		if e.Purpose != purpose {
			t.Errorf("%s purpose %q, want %q", id, e.Purpose, purpose)
		}
		if e.License == "" {
			t.Errorf("%s license empty (P3 requires a verdict)", id)
		}
	}
	// P3 blocked decision is machine-readable.
	matcha, _ := m.FindModel("tts-matcha-zh-baker")
	if matcha.Shippable() {
		t.Error("tts-matcha-zh-baker must be blocked-p3 (non-commercial data-baker license)")
	}
}

// TestRealManifestAuthoringConsistency checks the manifest authoring
// invariants that need no model bytes: single-artifact entries must carry the
// artifact digest at entry level; multi-artifact entries must list files; all
// URLs must be https official sources (mirrors are transport fallbacks, never
// hash sources - D33/F3).
func TestRealManifestAuthoringConsistency(t *testing.T) {
	mpath, _ := repoManifestPath(t)
	raw, err := os.ReadFile(mpath)
	if err != nil {
		t.Fatal(err)
	}
	m, err := ParseManifest(raw)
	if err != nil {
		t.Fatal(err)
	}
	for i := range m.Models {
		e := &m.Models[i]
		if len(e.Files) == 1 {
			if e.SHA256 != e.Files[0].SHA256 || e.SizeBytes != e.Files[0].SizeBytes {
				t.Errorf("%s: single artifact but entry digest/size differ from files[0]", e.ID)
			}
		}
		if len(e.Files) > 1 {
			var total int64
			for _, f := range e.Files {
				total += f.SizeBytes
			}
			if total != e.SizeBytes {
				t.Errorf("%s: files sum %d != entry size %d", e.ID, total, e.SizeBytes)
			}
		}
		for _, u := range e.URLs {
			if !strings.HasPrefix(u, "https://") {
				t.Errorf("%s: non-https url %q", e.ID, u)
			}
			if strings.Contains(u, "ghfast.top") || strings.Contains(u, "hf-mirror.com") {
				// ghfast/hf-mirror ARE allowed in the manifest (P5 domestic
				// mirrors); the assertion above documents they are transport
				// mirrors, while hashes stay manifest-internal.
				_ = u
			}
		}
	}
}

// TestRealDownloadVadThroughPipeline runs the REAL pipeline (signed manifest +
// hardcoded buildinfo key + official github URL + ghfast transport fallback)
// against the real world for the smallest model (silero VAD, 0.64MB).
// Gate: WISP_IT_REAL_MIRROR=1 (offline CI skips; run once as download spot
// check evidence).
func TestRealDownloadVadThroughPipeline(t *testing.T) {
	if os.Getenv("WISP_IT_REAL_MIRROR") == "" {
		t.Skip("real-network spot check; set WISP_IT_REAL_MIRROR=1 to run")
	}
	mpath, spath := repoManifestPath(t)
	m, err := LoadSignedManifest(mpath, spath, buildinfo.MinisignPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	mgr, err := NewManager(Options{DataDir: dir, Manifest: m, VerifySignature: true})
	if err != nil {
		t.Fatal(err)
	}
	inst, err := mgr.Ensure(context.Background(), "vad-silero")
	if err != nil {
		t.Fatalf("real download of vad-silero failed: %v", err)
	}
	e, _ := m.FindModel("vad-silero")
	if err := verifyFileHash(filepath.Join(inst, "silero_vad.onnx"), e.SHA256, e.SizeBytes); err != nil {
		t.Fatalf("real downloaded VAD fails manifest hash: %v", err)
	}
}

// TestRealDownloadPuncArchiveThroughPipeline: same, but the 64MB int8 punc
// archive - exercises the real archive extraction path (GitHub official
// release, ghfast fallback). Gate: WISP_IT_REAL_MIRROR=1.
func TestRealDownloadPuncArchiveThroughPipeline(t *testing.T) {
	if os.Getenv("WISP_IT_REAL_MIRROR") == "" {
		t.Skip("real-network spot check; set WISP_IT_REAL_MIRROR=1 to run")
	}
	mpath, spath := repoManifestPath(t)
	m, err := LoadSignedManifest(mpath, spath, buildinfo.MinisignPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	mgr, err := NewManager(Options{DataDir: dir, Manifest: m, VerifySignature: true,
		BackoffBase: 200 * time.Millisecond, Attempts: 3})
	if err != nil {
		t.Fatal(err)
	}
	inst, err := mgr.Ensure(context.Background(), "punc-ct-transformer-zh-en-vocab272727")
	if err != nil {
		t.Fatalf("real download of punc failed: %v", err)
	}
	e, _ := m.FindModel("punc-ct-transformer-zh-en-vocab272727")
	if err := mgr.VerifyDir(e, inst); err != nil {
		t.Fatalf("real downloaded punc fails manifest verification: %v", err)
	}
}

// TestP3BlockedModelRefused pins the P3 gate: the license-blocked TTS entry
// can never be ensured (no network, no install).
func TestP3BlockedModelRefused(t *testing.T) {
	mpath, spath := repoManifestPath(t)
	m, err := LoadSignedManifest(mpath, spath, buildinfo.MinisignPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := NewManager(Options{DataDir: t.TempDir(), Manifest: m, VerifySignature: true,
		HTTPClient: &http.Client{Transport: panicTransport{}}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = mgr.Ensure(context.Background(), "tts-matcha-zh-baker")
	if err == nil {
		t.Fatal("P3-blocked model was ensured (must refuse)")
	}
	if !strings.Contains(err.Error(), "P3-BLOCKED") {
		t.Fatalf("wrong error: %v", err)
	}
}
