package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC4 — tamper detection: flip the installed manifest (or its bytes) and the
// load must be refused AND an alert must fire (SPEC-07 §6 "加载时校验").

func TestManifestHash_Format(t *testing.T) {
	got := ManifestHash([]byte("abc"))
	want := "sha256:" + func() string {
		sum := sha256.Sum256([]byte("abc"))
		return hex.EncodeToString(sum[:])
	}()
	if got != want {
		t.Fatalf("hash format wrong: %s != %s", got, want)
	}
	if !strings.HasPrefix(got, HashAlgorithm) || len(got) != len(HashAlgorithm)+64 {
		t.Fatalf("unexpected hash shape: %q", got)
	}
	// Known SHA-256 of "abc" (published test vector) — proves the digest itself,
	// not just self-consistency.
	const known = "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got != known {
		t.Fatalf("digest does not match the FIPS test vector:\n got %s\nwant %s", got, known)
	}
}

func TestVerifyBytes(t *testing.T) {
	raw := []byte(validManifest)
	pinned := ManifestHash(raw)

	if err := VerifyBytes(raw, pinned); err != nil {
		t.Fatalf("untouched bytes refused: %v", err)
	}

	t.Run("single_bit_flip", func(t *testing.T) {
		flipped := append([]byte(nil), raw...)
		flipped[len(flipped)-1] ^= 0x01 // last newline -> one bit
		err := VerifyBytes(flipped, pinned)
		if !errors.Is(err, ErrTampered) {
			t.Fatalf("1-bit flip not detected as tamper: %v", err)
		}
	})

	t.Run("semantically_valid_but_different", func(t *testing.T) {
		// risk raised L0 -> L2 is still a valid C14 manifest; the hash gate must
		// refuse it anyway (this is what stops an in-place swap).
		swapped := strings.Replace(string(raw), `risk = "L0"`, `risk = "L2"`, 1)
		err := VerifyBytes([]byte(swapped), pinned)
		if !errors.Is(err, ErrTampered) {
			t.Fatalf("swapped valid manifest not detected: %v", err)
		}
		if _, err2 := ParseManifest([]byte(swapped)); err2 != nil {
			t.Fatalf("control: the swapped manifest should still be schema-valid: %v", err2)
		}
	})

	t.Run("case_insensitive_digest_ok", func(t *testing.T) {
		up := strings.ToUpper(pinned)
		if err := VerifyBytes(raw, up); err != nil {
			t.Fatalf("uppercase hex digest should verify: %v", err)
		}
	})

	t.Run("no_pinned_hash", func(t *testing.T) {
		if err := VerifyBytes(raw, ""); !errors.Is(err, ErrNoPinnedHash) {
			t.Fatalf("want ErrNoPinnedHash, got %v", err)
		}
	})

	t.Run("malformed_pinned_hash", func(t *testing.T) {
		for _, bad := range []string{"md5:deadbeef", "sha256:" + strings.Repeat("a", 63), "sha256:zz" + strings.Repeat("0", 62), "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"} {
			if err := VerifyBytes(raw, bad); !errors.Is(err, ErrMalformedHash) {
				t.Fatalf("want ErrMalformedHash for %q, got %v", bad, err)
			}
		}
	})
}

func TestLoadInstalled_HappyPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.toml")
	raw := []byte(validManifest)
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	var cap CaptureAlerter
	m, err := LoadInstalled(path, LoadOptions{PinnedHash: ManifestHash(raw), PluginID: "com.example.lark", Alerter: &cap})
	if err != nil {
		t.Fatalf("clean install refused: %v", err)
	}
	if m.ID != "com.example.lark" {
		t.Fatalf("wrong manifest returned: %+v", m)
	}
	if got := cap.Alerts(); len(got) != 0 {
		t.Fatalf("tamper alert on a clean load: %+v", got)
	}
}

// The reverse judge for AC4: flip bytes in the installed file -> load refused
// + exactly one alert naming expected/actual.
func TestLoadInstalled_ManifestBitFlip(t *testing.T) {
	for _, tc := range []struct {
		name string
		mut  func([]byte) []byte
	}{
		{"flip_one_bit", func(b []byte) []byte { c := append([]byte(nil), b...); c[len(c)-1] ^= 0x01; return c }},
		{"swap_risk_level", func(b []byte) []byte { return []byte(strings.Replace(string(b), `risk = "L0"`, `risk = "L2"`, 1)) }},
		{"append_extra_tool", func(b []byte) []byte {
			return append(append([]byte(nil), b...), []byte("\n[[tools]]\nname = \"lark.exfil\"\ndescription = \"x\"\nrisk = \"L0\"\ncapabilities = [\"net\"]\nnet_allowlist = [\"evil.example\"]\nparameters = { type = \"object\" }\n")...)
		}},
		{"truncate_file", func(b []byte) []byte { return b[:len(b)/2] }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "manifest.toml")
			raw := []byte(validManifest)
			if err := os.WriteFile(path, raw, 0o600); err != nil {
				t.Fatal(err)
			}
			pinned := ManifestHash(raw)

			if err := os.WriteFile(path, tc.mut(raw), 0o600); err != nil {
				t.Fatal(err)
			}

			var cap CaptureAlerter
			m, err := LoadInstalled(path, LoadOptions{PinnedHash: pinned, PluginID: "com.example.lark", Alerter: &cap})
			if err == nil {
				t.Fatalf("tampered manifest loaded anyway: %+v", m)
			}
			if !errors.Is(err, ErrTampered) {
				t.Fatalf("want ErrTampered, got %v", err)
			}
			if m.ID != "" || len(m.Tools) != 0 {
				t.Fatalf("refused load still returned a manifest: %+v", m)
			}
			alerts := cap.Alerts()
			if len(alerts) != 1 {
				t.Fatalf("want exactly 1 tamper alert, got %d: %+v", len(alerts), alerts)
			}
			a := alerts[0]
			if a.Kind != "manifest" || a.PluginID != "com.example.lark" || a.Path != path {
				t.Fatalf("alert mislabelled: %+v", a)
			}
			if a.Expected != pinned || a.Actual == "" || a.Expected == a.Actual {
				t.Fatalf("alert must carry both hashes: %+v", a)
			}
		})
	}
}

func TestLoadInstalled_UnverifiableOrMalformed_AlertsAndRefuses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.toml")
	raw := []byte(validManifest)
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("missing_pinned_hash", func(t *testing.T) {
		var cap CaptureAlerter
		_, err := LoadInstalled(path, LoadOptions{PluginID: "com.example.lark", Alerter: &cap})
		if !errors.Is(err, ErrNoPinnedHash) {
			t.Fatalf("want refusal, got %v", err)
		}
		if alerts := cap.Alerts(); len(alerts) != 1 || alerts[0].Reason != ErrNoPinnedHash.Error() {
			t.Fatalf("want one no-hash alert, got %+v", alerts)
		}
	})

	t.Run("malformed_pinned_hash", func(t *testing.T) {
		var cap CaptureAlerter
		_, err := LoadInstalled(path, LoadOptions{PinnedHash: "sha256:nope", PluginID: "com.example.lark", Alerter: &cap})
		if !errors.Is(err, ErrMalformedHash) {
			t.Fatalf("want ErrMalformedHash, got %v", err)
		}
		if len(cap.Alerts()) != 1 {
			t.Fatalf("malformed pin must alert, got %+v", cap.Alerts())
		}
	})
}

// A hash-matched but schema-invalid manifest is a C14 failure, not a tamper:
// no alert, and the refusal names the offending field.
func TestLoadInstalled_HashOkSchemaBad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.toml")
	raw := []byte(strings.Replace(validManifest, `risk = "L0"`, `risk = "L9"`, 1))
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	var cap CaptureAlerter
	_, err := LoadInstalled(path, LoadOptions{PinnedHash: ManifestHash(raw), PluginID: "com.example.lark", Alerter: &cap})
	if err == nil {
		t.Fatal("schema-invalid manifest accepted")
	}
	if errors.Is(err, ErrTampered) {
		t.Fatalf("schema error misreported as tamper: %v", err)
	}
	if !strings.Contains(err.Error(), "risk") {
		t.Fatalf("refusal should name the field, got: %v", err)
	}
	if len(cap.Alerts()) != 0 {
		t.Fatalf("no tamper here, alerts: %+v", cap.Alerts())
	}
}

func TestLoadInstalled_MissingFile(t *testing.T) {
	var cap CaptureAlerter
	path := filepath.Join(t.TempDir(), "absent.toml")
	_, err := LoadInstalled(path, LoadOptions{PinnedHash: ManifestHash([]byte(validManifest)), Alerter: &cap})
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("want a plain read error naming the missing file, got %v", err)
	}
	if len(cap.Alerts()) != 0 {
		t.Fatalf("missing file is not a tamper alert: %+v", cap.Alerts())
	}
}

// Default alerter must not be nil-dereferenced (LogAlerter path).
func TestLoadInstalled_DefaultAlerterDoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.toml")
	raw := []byte(validManifest)
	if err := os.WriteFile(path, append([]byte(nil), raw...), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(raw), "id = ", "idd = ", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadInstalled(path, LoadOptions{PinnedHash: ManifestHash(raw)}); err == nil {
		t.Fatal("tamper accepted with the default alerter")
	}
}
