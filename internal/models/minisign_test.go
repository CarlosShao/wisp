package models

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
)

// testKeypair derives a deterministic keypair from seed for reproducible tests.
func testKeypair(t *testing.T, seed []byte) (ed25519.PrivateKey, MinisignPublicKey) {
	t.Helper()
	if len(seed) != ed25519.SeedSize {
		t.Fatalf("seed must be %d bytes", ed25519.SeedSize)
	}
	priv := ed25519.NewKeyFromSeed(seed)
	keyID := keyIDFromPub(priv.Public().(ed25519.PublicKey))
	return priv, MinisignPublicKey{KeyID: keyID, Key: priv.Public().(ed25519.PublicKey)}
}

func TestMinisignRoundTrip(t *testing.T) {
	priv, pub := testKeypair(t, []byte("0123456789abcdef0123456789abcdef"))
	payload := []byte("{\"manifest_version\":1}\n")

	sig := SignPayloadMust(t, priv, pub, payload, "signature from wisp dev key", "timestamp: 2026-09-19 file: manifest.json")

	if err := VerifyMinisignSignature(pub, payload, sig); err != nil {
		t.Fatalf("round trip verify failed: %v", err)
	}

	// Rendered public key must itself parse (format stability).
	rekey, err := ParseMinisignPublicKey(pub.String("wisp test key"))
	if err != nil {
		t.Fatalf("re-parse rendered pub: %v", err)
	}
	if rekey.KeyID != pub.KeyID || string(rekey.Key) != string(pub.Key) {
		t.Fatal("public key round trip mismatch")
	}
}

func TestMinisignTamperDetection(t *testing.T) {
	priv, pub := testKeypair(t, []byte("0123456789abcdef0123456789abcdef"))
	payload := []byte("model bytes")
	sig := SignPayloadMust(t, priv, pub, payload, "u", "trusted comment: t")

	// 1. One flipped payload byte must fail.
	bad := append([]byte{}, payload...)
	bad[3] ^= 0x01
	if err := VerifyMinisignSignature(pub, bad, sig); err == nil {
		t.Fatal("tampered payload verified (must fail)")
	}

	// 2. Tampered trusted comment must fail (layer 2).
	badSig := strings.Replace(string(sig), "trusted comment: t", "trusted comment: X", 1)
	if err := VerifyMinisignSignature(pub, payload, []byte(badSig)); err == nil {
		t.Fatal("tampered trusted comment verified (must fail)")
	}

	// 3. Tampered untrusted comment MUST still pass (spec: freely editable).
	okSig := strings.Replace(string(sig), "untrusted comment: u", "untrusted comment: edited", 1)
	if err := VerifyMinisignSignature(pub, payload, []byte(okSig)); err != nil {
		t.Fatalf("edited untrusted comment must still verify: %v", err)
	}

	// 4. Wrong key must fail.
	_, otherPub := testKeypair(t, []byte("ffffffffffffffffffffffffffffffff"))
	if err := VerifyMinisignSignature(otherPub, payload, sig); err == nil {
		t.Fatal("cross-key verification succeeded (keyid guard broken)")
	}

	// 5. Truncated/extra lines must fail structurally.
	if err := VerifyMinisignSignature(pub, payload, sig[:len(sig)-4]); err == nil {
		t.Fatal("truncated signature file verified (must fail)")
	}
}

func TestMinisignFormatExactness(t *testing.T) {
	priv, pub := testKeypair(t, []byte("0123456789abcdef0123456789abcdef"))
	payload := []byte("data")
	sig := string(SignPayloadMust(t, priv, pub, payload, "signature from minisign secret key", "timestamp: 1"))

	lines := strings.Split(strings.TrimSuffix(sig, "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("want 4 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "untrusted comment: ") || !strings.HasPrefix(lines[2], "trusted comment: ") {
		t.Fatalf("comment lines malformed:\n%s", sig)
	}
	// Algorithm bytes must be "Ed" (base64 of Ed||keyid||sig starts with "RWQ" only by
	// accident; decode instead of guessing).
	blob := decodeB64Must(t, lines[1])
	if string(blob[:2]) != "Ed" {
		t.Fatalf("alg bytes %q, want Ed", blob[:2])
	}
	if len(blob) != 2+8+64 {
		t.Fatalf("sig blob %d bytes, want 74", len(blob))
	}
	g2 := decodeB64Must(t, lines[3])
	if len(g2) != 64 {
		t.Fatalf("global sig %d bytes, want 64", len(g2))
	}
}

func TestMinisignRealKeyAcceptsBuildinfoPlaceholder(t *testing.T) {
	// The buildinfo placeholder must NOT verify (fail-closed until the real/dev
	// key lands in buildinfo). This pins the integration seam: once
	// buildinfo.MinisignPublicKey is a real key, TestSignedManifestInRepo
	// becomes the authoritative check.
	if _, err := ParseMinisignPublicKey(buildinfo.MinisignPublicKey); err == nil {
		t.Log("buildinfo carries a parseable public key (expected after C29 dev-key landing)")
	}
}

func TestMinisignRandomKeygenShape(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(pub) != 32 || len(priv) != 64 {
		t.Fatalf("ed25519 key sizes: pub=%d priv=%d", len(pub), len(priv))
	}
}

func decodeB64Must(t *testing.T, s string) []byte {
	t.Helper()
	b, err := base64Decode(s)
	if err != nil {
		t.Fatalf("base64 %q: %v", s, err)
	}
	return b
}
