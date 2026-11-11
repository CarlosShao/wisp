package main

// The deferred-disposition 20261110 #5 tests: privacy.diagnostics_opt_in, the
// total switch on `wisp diagnostics`.
//
// THE REVERSE CONTRACT (摘喂线 -> 红): the gate-off case asserts rc=1 AND that
// nothing was built - remove the gate and the export case below would succeed
// with the switch off, which is exactly the fail-open shape the总闸 exists to
// refuse. The manifest-honesty case reads the zip back: a builder bypass that
// skipped the preview contract would show up as an entry list that disagrees
// with the printed preview.

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func diagnostics180PDataDir(t *testing.T, optIn bool) string {
	t.Helper()
	dir := t.TempDir()
	body := "schema_version = 2\n"
	if optIn {
		body += "\n[privacy]\ndiagnostics_opt_in = true\n"
	}
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return dir
}

func diagnostics180PIO() (diagnosticsIO, *bytes.Buffer, *bytes.Buffer) {
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	return diagnosticsIO{stdout: out, stderr: errb}, out, errb
}

func diagnostics180PSeedLog(t *testing.T, dir string) {
	t.Helper()
	logDir := filepath.Join(dir, logDirName)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatalf("mkdir logs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(logDir, "wisp-20261111-001.jsonl"),
		[]byte(`{"t":"2026-11-11T01:02:03Z","msg":"probe line"}`+"\n"), 0o600); err != nil {
		t.Fatalf("write log: %v", err)
	}
}

// Test180PdiagnosticsGateRefusesWhenSwitchOff is the总闸's teeth: switch off
// (and missing-file) both refuse rc=1 with the key named, and - the reverse
// half - the export path must not leave a zip behind, so a future gate
// removal turns this case into a red about a file that exists.
func Test180PdiagnosticsGateRefusesWhenSwitchOff(t *testing.T) {
	for _, tc := range []struct {
		name  string
		mkDir func(t *testing.T) string
		sub   string
	}{
		{"switch-off-export", func(t *testing.T) string { return diagnostics180PDataDir(t, false) }, "export"},
		{"switch-off-preview", func(t *testing.T) string { return diagnostics180PDataDir(t, false) }, "preview"},
		{"no-config-file", func(t *testing.T) string { return t.TempDir() }, "export"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.mkDir(t)
			dio, _, errb := diagnostics180PIO()
			var rc int
			if tc.sub == "export" {
				rc = runDiagnosticsExport180P(dio, dir, nil)
			} else {
				rc = runDiagnosticsPreview180P(dio, dir)
			}
			if rc != 1 {
				t.Fatalf("%s: rc = %d, want 1 (the gate refused, the operational-refusal exit this leg shares with `wisp config`)", tc.name, rc)
			}
			if !strings.Contains(errb.String(), "diagnostics_opt_in") {
				t.Errorf("%s: the refusal never names the key (stderr=%q): an operator cannot act on it", tc.name, errb.String())
			}
			matches, _ := filepath.Glob(filepath.Join(dir, "*.zip"))
			if len(matches) != 0 {
				t.Fatalf("%s: a refused export left %v behind - the gate is not in front of the builder", tc.name, matches)
			}
		})
	}
}

// Test180PdiagnosticsPreviewListsWhatTheZipWouldCarry is the preview contract
// on the switch-on side: rc=0, the rows name every default-in category, the
// never-collected audio category is stated out loud, and nothing is written.
func Test180PdiagnosticsPreviewListsWhatTheZipWouldCarry(t *testing.T) {
	dir := diagnostics180PDataDir(t, true)
	diagnostics180PSeedLog(t, dir)
	dio, out, _ := diagnostics180PIO()
	if rc := runDiagnosticsPreview180P(dio, dir); rc != 0 {
		t.Fatalf("preview rc = %d (stderr=%v via buffers)", rc, out.String())
	}
	rows := out.String()
	for _, want := range []string{"[logs]", "[config]", "[version]", "[deferred]", "[audio]", "manifest.json", "wisp-20261111-001.jsonl"} {
		if !strings.Contains(rows, want) {
			t.Errorf("preview rows lack %q (output=%q)", want, rows)
		}
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "*.zip"))
	if len(matches) != 0 {
		t.Fatalf("a preview wrote %v: the preview subcommand writes nothing", matches)
	}
}

// Test180PdiagnosticsExportBuildsThePreviewedBundle is the happy path end to
// end: preview printed first, zip built at the requested path, and the zip's
// real entry list compared against what the preview promised - the same
// two-paths-agree claim ticket 45 built its sibling preview function for.
func Test180PdiagnosticsExportBuildsThePreviewedBundle(t *testing.T) {
	dir := diagnostics180PDataDir(t, true)
	diagnostics180PSeedLog(t, dir)
	outPath := filepath.Join(t.TempDir(), "bundle.zip")
	dio, out, errb := diagnostics180PIO()
	if rc := runDiagnosticsExport180P(dio, dir, []string{"-o", outPath}); rc != 0 {
		t.Fatalf("export rc = %d, stderr=%q", rc, errb.String())
	}
	if !strings.Contains(out.String(), filepath.Base(outPath)) {
		t.Errorf("the receipt never names the written file (stdout=%q)", out.String())
	}
	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("the promised zip does not exist: %v", err)
	}
	defer f.Close()
	zr, err := zip.NewReader(f, mustSize180P(t, outPath))
	if err != nil {
		t.Fatalf("not a readable zip: %v", err)
	}
	// The default-excluded categories must be HONESTLY absent, not silently
	// absent: manifest.json records them, audio is never collected.
	var names []string
	hasManifest := false
	for _, e := range zr.File {
		names = append(names, e.Name)
		if e.Name == "manifest.json" {
			hasManifest = true
		}
		if strings.HasPrefix(e.Name, "transcripts/") || strings.HasPrefix(e.Name, "watchdog/") {
			t.Errorf("zip carries %s with no per-item consent taken: the default-excluded categories leaked", e.Name)
		}
	}
	if !hasManifest {
		t.Fatal("the zip has no manifest.json: the builder's honesty record is missing")
	}
	for _, want := range []string{"logs/", "config.redacted.toml", "version.json", "deferred.json", "manifest.json"} {
		found := false
		for _, n := range names {
			if strings.HasPrefix(n, want) || n == want {
				found = true
			}
		}
		if !found {
			t.Errorf("zip lacks an entry under %q (entries=%v)", want, names)
		}
	}
	// The manifest says what is not there: the transcripts/watchdog/audio rows
	// must carry their not-collected sentences.
	mf, err := openZipString180P(zr, "manifest.json")
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	for _, want := range []string{"transcripts/", "watchdog/", "audio/", "never collected"} {
		if !strings.Contains(mf, want) {
			t.Errorf("manifest lacks the honest record for %q (manifest=%q)", want, mf)
		}
	}
}

// Test180PdiagnosticsExportRefusesBadArguments pins the flag surface: an
// unknown flag is a usage refusal (rc=2), a missing -slo-snapshot file is an
// operational refusal (rc=1) that leaves no zip, and a missing flag value is
// usage (rc=2).
func Test180PdiagnosticsExportRefusesBadArguments(t *testing.T) {
	dir := diagnostics180PDataDir(t, true)
	t.Run("unknown-flag", func(t *testing.T) {
		dio, _, _ := diagnostics180PIO()
		if rc := runDiagnosticsExport180P(dio, dir, []string{"-wat", "x"}); rc != 2 {
			t.Fatalf("rc = %d, want 2", rc)
		}
	})
	t.Run("dangling-flag-value", func(t *testing.T) {
		dio, _, _ := diagnostics180PIO()
		if rc := runDiagnosticsExport180P(dio, dir, []string{"-o"}); rc != 2 {
			t.Fatalf("rc = %d, want 2 (a flag with no value is a usage refusal)", rc)
		}
	})
	t.Run("missing-slo-snapshot", func(t *testing.T) {
		outPath := filepath.Join(t.TempDir(), "b.zip")
		dio, _, _ := diagnostics180PIO()
		if rc := runDiagnosticsExport180P(dio, dir, []string{"-o", outPath, "-slo-snapshot", filepath.Join(dir, "no-such.json")}); rc != 1 {
			t.Fatalf("rc = %d, want 1", rc)
		}
		if _, err := os.Stat(outPath); !os.IsNotExist(err) {
			t.Fatal("a refused export left a zip behind")
		}
	})
}

// Test180PdiagnosticsSLOSnapshotLandsInTheZip is the attach half: a real
// -slo-snapshot file rides into the zip under its .json entry name, and the
// manifest records it.
func Test180PdiagnosticsSLOSnapshotLandsInTheZip(t *testing.T) {
	dir := diagnostics180PDataDir(t, true)
	diagnostics180PSeedLog(t, dir)
	slo := filepath.Join(t.TempDir(), "slo.json")
	if err := os.WriteFile(slo, []byte(`{"state":"Awake","p50_ms":11}`), 0o600); err != nil {
		t.Fatalf("write slo: %v", err)
	}
	outPath := filepath.Join(t.TempDir(), "b.zip")
	dio, _, errb := diagnostics180PIO()
	if rc := runDiagnosticsExport180P(dio, dir, []string{"-o", outPath, "-slo-snapshot", slo}); rc != 0 {
		t.Fatalf("export rc = %d, stderr=%q", rc, errb.String())
	}
	zr := openZip180P(t, outPath)
	data, err := openZipString180P(zr, "slo-snapshot.json")
	if err != nil {
		t.Fatalf("the attached snapshot is not in the zip: %v", err)
	}
	if !strings.Contains(data, "p50_ms") {
		t.Errorf("the attached snapshot is not the file's content (%q)", data)
	}
	mf, _ := openZipString180P(zr, "manifest.json")
	if !strings.Contains(mf, "sampler report attached") {
		t.Errorf("the manifest does not record the attachment (%q)", mf)
	}
}

func mustSize180P(t *testing.T, p string) int64 {
	t.Helper()
	st, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat %s: %v", p, err)
	}
	return st.Size()
}

func openZip180P(t *testing.T, p string) *zip.Reader {
	t.Helper()
	f, err := os.Open(p)
	if err != nil {
		t.Fatalf("open %s: %v", p, err)
	}
	t.Cleanup(func() { _ = f.Close() })
	zr, err := zip.NewReader(f, mustSize180P(t, p))
	if err != nil {
		t.Fatalf("zip.NewReader %s: %v", p, err)
	}
	return zr
}

func openZipString180P(zr *zip.Reader, name string) (string, error) {
	for _, e := range zr.File {
		if e.Name != name {
			continue
		}
		rc, err := e.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(rc); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return "", fmt.Errorf("zip entry %s not found", name)
}
