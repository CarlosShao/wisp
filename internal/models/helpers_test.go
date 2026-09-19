package models

import (
	"crypto/ed25519"
	"encoding/base64"
	"testing"
)

// SignPayloadMust wraps SignPayload for tests (fatal on format errors).
func SignPayloadMust(t *testing.T, priv ed25519.PrivateKey, pub MinisignPublicKey, payload []byte, untrusted, trusted string) []byte {
	t.Helper()
	sig, err := SignPayload(priv, pub, payload, untrusted, trusted)
	if err != nil {
		t.Fatalf("SignPayload: %v", err)
	}
	return sig
}

// base64Decode is the test-facing alias of std base64 (kept tiny on purpose).
func base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
