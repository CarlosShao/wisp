package models

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/blake2b"
)

// DeriveKeyID derives the 8-byte minisign key id from a public key (first 8
// bytes of SHA-256): deterministic, so regenerated dev keys with the same
// seed carry the same id, and the id is verified to match across pubkey file
// and signatures. Shared with the signing tool.
func DeriveKeyID(pub ed25519.PublicKey) [sizeKeyID]byte {
	sum := sha256.Sum256(pub)
	var id [sizeKeyID]byte
	copy(id[:], sum[:sizeKeyID])
	return id
}

func keyIDFromPub(pub ed25519.PublicKey) [sizeKeyID]byte {
	return DeriveKeyID(pub)
}

// Hand-written minisign signature verification (C29, D17/D41e: minisign is the
// model/update signing mechanism, distinct from CA code signing). No minisign
// library is pulled in: the file format is a fixed 4-line layout and the
// crypto is Ed25519, so a strict parser plus stdlib crypto/ed25519 (and
// x/crypto/blake2b for the pre-hashed variant) reproduces real minisign
// behavior byte for byte.
//
// Signature file format (jedisct1/minisign spec):
//
//	untrusted comment: <arbitrary>
//	<base64: alg[2] || keyid[8] || sig[64]>          sig = Ed25519 over payload
//	trusted comment: <comment>                        tamper-evident (layer 2)
//	<base64: sig2[64]>                                sig2 = Ed25519 over sig || comment
//
// alg is "Ed" (legacy: payload = file bytes) or "ED" (pre-hashed: payload =
// Blake2b-512(file bytes), the modern minisign default). Both are accepted;
// signatures produced by tools/signmodels use "Ed".
//
// Public key format:
//
//	untrusted comment: <arbitrary>
//	<base64: "Ed" || keyid[8] || pub[32]>
//
// The keyid in the signature must equal the keyid of the verification key, or
// verification fails (cross-key replay guard).

const (
	sigAlgLegacy    = "Ed" // direct Ed25519 over the payload
	sigAlgPrehashed = "ED" // Ed25519 over Blake2b-512 of the payload
)

// minisign sizes (bytes).
const (
	sizeAlg   = 2
	sizeKeyID = 8
	sizeSig   = 64
)

// MinisignPublicKey is a parsed, format-checked minisign public key.
type MinisignPublicKey struct {
	KeyID [sizeKeyID]byte
	Key   ed25519.PublicKey // 32 bytes
}

// ParseMinisignPublicKey parses a public key in minisign two-line format
// (exactly what `minisign -p` writes and what buildinfo.MinisignPublicKey
// embeds as its second line). comment is ignored.
func ParseMinisignPublicKey(pub string) (MinisignPublicKey, error) {
	var out MinisignPublicKey
	lines := strings.Split(strings.TrimSuffix(pub, "\n"), "\n")
	if len(lines) != 2 {
		return out, fmt.Errorf("minisign: public key must be 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "untrusted comment: ") {
		return out, fmt.Errorf("minisign: public key line 1 must start %q", "untrusted comment: ")
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(lines[1]))
	if err != nil {
		return out, fmt.Errorf("minisign: public key base64: %w", err)
	}
	if len(raw) != sizeAlg+sizeKeyID+ed25519.PublicKeySize {
		return out, fmt.Errorf("minisign: public key blob is %d bytes, want %d", len(raw), sizeAlg+sizeKeyID+ed25519.PublicKeySize)
	}
	if string(raw[:sizeAlg]) != sigAlgLegacy {
		return out, fmt.Errorf("minisign: public key alg %q, want %q", raw[:sizeAlg], sigAlgLegacy)
	}
	copy(out.KeyID[:], raw[sizeAlg:sizeAlg+sizeKeyID])
	out.Key = ed25519.PublicKey(raw[sizeAlg+sizeKeyID:])
	return out, nil
}

// String renders the key back to minisign two-line format (stable round trip).
func (k MinisignPublicKey) String(comment string) string {
	blob := make([]byte, 0, sizeAlg+sizeKeyID+len(k.Key))
	blob = append(blob, sigAlgLegacy...)
	blob = append(blob, k.KeyID[:]...)
	blob = append(blob, k.Key...)
	return "untrusted comment: " + comment + "\n" + base64.StdEncoding.EncodeToString(blob) + "\n"
}

// KeyIDHex returns the 8-byte key id as 16 hex characters (log/diagnostic use;
// never the key material itself).
func (k MinisignPublicKey) KeyIDHex() string {
	return hex.EncodeToString(k.KeyID[:])
}

// VerifyMinisignSignature verifies sigBytes (the whole .minisig file) against
// payload (the exact bytes of the signed file) using pub. Every layer is
// checked: key id match, Ed25519 over the payload, then Ed25519 over
// (sig || trusted comment). Any mismatch fails the whole verification: there
// is no "trusted comment is only advisory" soft path.
func VerifyMinisignSignature(pub MinisignPublicKey, payload, sigBytes []byte) error {
	sig, err := parseMinisignSigFile(sigBytes)
	if err != nil {
		return err
	}
	if sig.keyID != pub.KeyID {
		return fmt.Errorf("minisign: signature keyid %x does not match public key keyid %x", sig.keyID, pub.KeyID)
	}

	var payloadHash hash.Hash
	switch sig.alg {
	case sigAlgLegacy:
		payloadHash = nil
	case sigAlgPrehashed:
		h, err := blake2b.New512(nil)
		if err != nil {
			return fmt.Errorf("minisign: blake2b: %w", err)
		}
		payloadHash = h
	default:
		return fmt.Errorf("minisign: unsupported signature algorithm %q (want Ed or ED)", sig.alg)
	}
	digest := payload
	if payloadHash != nil {
		_, _ = payloadHash.Write(payload)
		digest = payloadHash.Sum(nil)
	}
	if !ed25519.Verify(pub.Key, digest, sig.sig) {
		return fmt.Errorf("minisign: Ed25519 verification FAILED (payload signature invalid)")
	}
	if !ed25519.Verify(pub.Key, append(append([]byte{}, sig.sig...), sig.trustedComment...), sig.sig2) {
		return fmt.Errorf("minisign: trusted-comment signature verification FAILED")
	}
	return nil
}

type minisignSig struct {
	alg            string
	keyID          [sizeKeyID]byte
	sig            []byte
	trustedComment []byte
	sig2           []byte
}

// parseMinisignSigFile strictly parses the 4-line minisign signature file.
func parseMinisignSigFile(b []byte) (minisignSig, error) {
	var out minisignSig
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	if len(lines) != 4 {
		return out, fmt.Errorf("minisign: signature file must be 4 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "untrusted comment: ") {
		return out, fmt.Errorf("minisign: signature line 1 must start %q", "untrusted comment: ")
	}
	if !strings.HasPrefix(lines[2], "trusted comment: ") {
		return out, fmt.Errorf("minisign: signature line 3 must start %q", "trusted comment: ")
	}
	blob, err := base64.StdEncoding.DecodeString(strings.TrimSpace(lines[1]))
	if err != nil {
		return out, fmt.Errorf("minisign: signature base64: %w", err)
	}
	if len(blob) != sizeAlg+sizeKeyID+sizeSig {
		return out, fmt.Errorf("minisign: signature blob is %d bytes, want %d", len(blob), sizeAlg+sizeKeyID+sizeSig)
	}
	out.alg = string(blob[:sizeAlg])
	copy(out.keyID[:], blob[sizeAlg:sizeAlg+sizeKeyID])
	out.sig = blob[sizeAlg+sizeKeyID:]
	out.trustedComment = []byte(lines[2][len("trusted comment: "):])
	sig2, err := base64.StdEncoding.DecodeString(strings.TrimSpace(lines[3]))
	if err != nil {
		return out, fmt.Errorf("minisign: global signature base64: %w", err)
	}
	if len(sig2) != sizeSig {
		return out, fmt.Errorf("minisign: global signature is %d bytes, want %d", len(sig2), sizeSig)
	}
	out.sig2 = sig2
	return out, nil
}

// SignPayload is the single helper shared with the signer tool: it produces
// the full 4-line minisign signature file for payload under priv, using the
// legacy "Ed" algorithm and key id derived from pub. Exported so
// tools/signmodels (separate module) and in-repo tests cannot drift apart.
//
// comment fields are recorded verbatim; trustedComment becomes tamper-evident.
func SignPayload(priv ed25519.PrivateKey, pub MinisignPublicKey, payload []byte, untrustedComment, trustedComment string) ([]byte, error) {
	if len(priv) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("minisign: secret key is %d bytes, want %d", len(priv), ed25519.PrivateKeySize)
	}
	sig := ed25519.Sign(priv, payload)
	sig2 := ed25519.Sign(priv, append(append([]byte{}, sig...), []byte(trustedComment)...))
	blob := make([]byte, 0, sizeAlg+sizeKeyID+sizeSig)
	blob = append(blob, sigAlgLegacy...)
	blob = append(blob, pub.KeyID[:]...)
	blob = append(blob, sig...)
	var b bytes.Buffer
	b.WriteString("untrusted comment: " + untrustedComment + "\n")
	b.WriteString(base64.StdEncoding.EncodeToString(blob) + "\n")
	b.WriteString("trusted comment: " + trustedComment + "\n")
	b.WriteString(base64.StdEncoding.EncodeToString(sig2) + "\n")
	return b.Bytes(), nil
}
