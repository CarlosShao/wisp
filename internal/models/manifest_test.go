package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func goodEntryJSON(id string) string {
	return `{
		"id": "` + id + `", "purpose": "vad", "urls": ["https://example.com/` + id + `/silero_vad.onnx"],
		"sha256": "` + strings.Repeat("ab", 32) + `", "size_bytes": 671089,
		"license": "MIT", "quant": "fp32"
	}`
}

func parseOK(t *testing.T, json string) *Manifest {
	t.Helper()
	m, err := ParseManifest([]byte(json))
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	return m
}

func TestManifestParseSingleArtifact(t *testing.T) {
	m := parseOK(t, `{"manifest_version":1,"models":[`+goodEntryJSON("vad-silero")+`]}`)
	e, err := m.FindModel("vad-silero")
	if err != nil {
		t.Fatal(err)
	}
	if e.ArtifactCount() != 1 {
		t.Fatalf("artifact count %d", e.ArtifactCount())
	}
	f := e.Artifact(0)
	if f.Path != "silero_vad.onnx" || f.SizeBytes != 671089 {
		t.Fatalf("implicit file wrong: %+v", f)
	}
	inst := e.InstalledFiles()
	if len(inst) != 1 || inst[0].Path != "silero_vad.onnx" {
		t.Fatalf("installed files wrong: %+v", inst)
	}
}

func TestManifestParseMultiArtifactWithArchive(t *testing.T) {
	m := parseOK(t, `{"manifest_version":1,"models":[{
		"id": "kws", "purpose": "kws",
		"urls": ["https://github.com/x/a.tar.bz2"],
		"sha256": "`+strings.Repeat("11", 32)+`", "size_bytes": 100,
		"license": "Apache-2.0", "quant": "int8",
		"files": [{
			"path": "a.tar.bz2", "urls": ["https://github.com/x/a.tar.bz2"],
			"sha256": "`+strings.Repeat("11", 32)+`", "size_bytes": 100,
			"archive": {"format": "tar.bz2", "files": [
				{"path": "encoder.onnx", "sha256": "`+strings.Repeat("22", 32)+`", "size_bytes": 90},
				{"path": "tokens.txt", "sha256": "`+strings.Repeat("33", 32)+`", "size_bytes": 10}
			]}
		}]
	}]}`)
	e, _ := m.FindModel("kws")
	inst := e.InstalledFiles()
	if len(inst) != 2 || inst[0].Path != "encoder.onnx" {
		t.Fatalf("installed files wrong: %+v", inst)
	}
}

func TestManifestRejectsStructuralDrift(t *testing.T) {
	cases := map[string]string{
		"version":       `{"manifest_version":2,"models":[` + goodEntryJSON("a") + `]}`,
		"empty":         `{"manifest_version":1,"models":[]}`,
		"bad-purpose":   strings.Replace(goodEntryJSON("a"), `"purpose": "vad"`, `"purpose": "voip"`, 1),
		"bad-quant":     strings.Replace(goodEntryJSON("a"), `"quant": "fp32"`, `"quant": "int4"`, 1),
		"no-license":    strings.Replace(goodEntryJSON("a"), `"license": "MIT"`, `"license": ""`, 1),
		"short-hash":    strings.Replace(goodEntryJSON("a"), strings.Repeat("ab", 32), strings.Repeat("ab", 16), 1),
		"dup-id":        `{"manifest_version":1,"models":[` + goodEntryJSON("a") + `,` + goodEntryJSON("a") + `]}`,
		"traverse-path": strings.Replace(goodEntryJSON("a"), "silero_vad.onnx", "../../etc/passwd", 1),
	}
	for name, json := range cases {
		if _, err := ParseManifest([]byte(json)); err == nil {
			t.Errorf("%s: accepted (must reject)", name)
		}
	}
}

func TestManifestCombinedDigestContract(t *testing.T) {
	// Entry-level sha256 = sha256 over concatenated artifact bytes in order.
	b1 := []byte("artifact-one")
	b2 := []byte("artifact-two")
	h := sha256.New()
	h.Write(b1)
	h.Write(b2)
	want := hex.EncodeToString(h.Sum(nil))

	e := ModelEntry{ID: "m", Purpose: "tts", License: "L", Quant: "fp32"}
	got := e.CombinedDigest([][]byte{b1, b2})
	if got != want {
		t.Fatalf("combined digest %s, want %s", got, want)
	}
}
