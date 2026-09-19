package observe

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestDiagnosticsBundleCollectsAndRedacts(t *testing.T) {
	dir := t.TempDir()
	// Seed a redacted-at-rest log plus a config file that would leak if the
	// bundle skipped redaction.
	logDir := filepath.Join(dir, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatal(err)
	}
	logLine := "{\"level\":\"INFO\",\"msg\":\"provider call\",\"api_key\":\"sk-SEEDLEAK1234567890\"}\n"
	if err := os.WriteFile(filepath.Join(logDir, logFilePrefix+"20260919-001"+logFileExt), []byte(logLine), 0o644); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(dir, "config.toml")
	cfg := "[llm.mock]\napi_key_ref = \"env:MOCK_KEY\"\n# comment mentions sk-SEEDCFG9999999999 inline\n"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	slo := []byte(`{"state":"Sleeping","pass":true}`)

	out := filepath.Join(dir, "diag.zip")
	path, err := BuildDiagnosticsBundle(BundleOptions{
		OutPath:     out,
		LogDir:      logDir,
		ConfigPath:  cfgPath,
		RedactPaths: true,
		SLOSnapshot: slo,
		Version: map[string]string{
			"version": "0.0.0-dev", "commit": "test", "env": "test",
		},
		Deferred: []DeferredItem{{ID: "D-test", Ticket: "45", Summary: "bundle UX"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if path != out {
		t.Fatalf("bundle path mismatch: %s", path)
	}

	zr, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	found := map[string]string{}
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(rc)
		_ = rc.Close()
		found[f.Name] = string(data)
	}
	for _, want := range []string{"manifest.json", "version.json", "deferred.json", "slo-snapshot.json", "config.redacted.toml"} {
		if _, ok := found[want]; !ok {
			t.Errorf("bundle missing entry %q (have %v)", want, keys(found))
		}
	}
	var oneLog bool
	for name, content := range found {
		if strings.HasPrefix(name, "logs/") {
			oneLog = true
			if strings.Contains(content, "sk-SEEDLEAK1234567890") {
				t.Errorf("seeded log key survived into the bundle: %s", name)
			}
		}
	}
	if !oneLog {
		t.Errorf("no log entry in bundle: %v", keys(found))
	}
	if strings.Contains(found["config.redacted.toml"], "sk-SEEDCFG9999999999") {
		t.Error("seeded config inline key survived into the bundle")
	}
	var man struct {
		Contents []struct{ Entry string } `json:"contents"`
	}
	if err := json.Unmarshal([]byte(found["manifest.json"]), &man); err != nil {
		t.Fatalf("manifest is not JSON: %v", err)
	}
	if len(man.Contents) < 5 {
		t.Fatalf("manifest contents incomplete: %s", found["manifest.json"])
	}
	var def []DeferredItem
	if err := json.Unmarshal([]byte(found["deferred.json"]), &def); err != nil {
		t.Fatalf("deferred.json is not a register: %v", err)
	}
	if len(def) != 1 || def[0].ID != "D-test" {
		t.Fatalf("deferred register wrong: %s", found["deferred.json"])
	}
}

func TestDiagnosticsBundleFailClosedOnSurvivingKey(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	// The abort path is defensive: it fires when the post-mask detector
	// still matches. Pin the contract by injecting a detector pattern the
	// masker does not consume (same-package variable swap) - this keeps the
	// fail-closed branch exercised without weakening the real rules.
	cfg := "note = ZZZZTRIGGERMARKER\n"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	orig := snapshotDetectRe
	snapshotDetectRe = regexp.MustCompile(`ZZZZTRIGGERMARKER`)
	t.Cleanup(func() { snapshotDetectRe = orig })

	_, clean := Redactor{}.ConfigSnapshot(cfg)
	if clean {
		t.Fatal("precondition: injected detector must flag the fixture")
	}
	_, err := BuildDiagnosticsBundle(BundleOptions{
		OutPath:    filepath.Join(dir, "no.zip"),
		ConfigPath: cfgPath,
	})
	if err == nil {
		t.Fatal("bundle must abort (fail-closed) on unclean config snapshot")
	}
	if _, err := os.Stat(filepath.Join(dir, "no.zip")); !os.IsNotExist(err) {
		t.Fatal("aborted bundle must not leave a zip behind")
	}
}

func TestDiagnosticsBundleRecordsMissingPieces(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "minimal.zip")
	if _, err := BuildDiagnosticsBundle(BundleOptions{OutPath: out}); err != nil {
		t.Fatal(err)
	}
	zr, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	var manifest []byte
	for _, f := range zr.File {
		if f.Name == "manifest.json" {
			rc, _ := f.Open()
			manifest, _ = io.ReadAll(rc)
			_ = rc.Close()
		}
	}
	if !bytes.Contains(manifest, []byte("not collected")) {
		t.Fatalf("missing pieces must be recorded with reasons: %s", manifest)
	}
	if !bytes.Contains(manifest, []byte("ticket 45")) {
		t.Fatalf("absent deferred register must be recorded (ticket 45): %s", manifest)
	}
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
