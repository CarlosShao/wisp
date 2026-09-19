package secret

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Reference kinds (SPEC-03 §5.3): config.toml stores references, never
// plaintext values (D36#5). "dpapi:<blob-id>" names a DPAPI blob file under
// <data>\secrets\; "env:NAME" names a process environment variable (CI/test:
// the value may be any placeholder - the mock LLM does not validate keys).
const (
	RefKindDPAPI = "dpapi"
	RefKindEnv   = "env"

	RefPrefixDPAPI = RefKindDPAPI + ":"
	RefPrefixEnv   = RefKindEnv + ":"
)

// ErrInvalidRef marks refs that are not exactly one of the two frozen forms,
// or that name something outside the safe charset (no traversal, no
// separators, no whitespace).
var ErrInvalidRef = errors.New("secret: invalid ref")

// ParseRef splits a reference into (kind, value).
func ParseRef(ref string) (kind, value string, err error) {
	switch {
	case strings.HasPrefix(ref, RefPrefixDPAPI):
		id := strings.TrimPrefix(ref, RefPrefixDPAPI)
		if !ValidBlobID(id) {
			return "", "", fmt.Errorf("%w: %q: blob id must be 1-128 chars of [A-Za-z0-9._-] and not \".\"/\"..\"", ErrInvalidRef, ref)
		}
		return RefKindDPAPI, id, nil
	case strings.HasPrefix(ref, RefPrefixEnv):
		name := strings.TrimPrefix(ref, RefPrefixEnv)
		if name == "" || len(name) > 128 || strings.ContainsAny(name, " \t\r\n=:") {
			return "", "", fmt.Errorf("%w: %q: env name must be 1-128 chars without whitespace, \"=\" or \":\"", ErrInvalidRef, ref)
		}
		return RefKindEnv, name, nil
	default:
		return "", "", fmt.Errorf("%w: %q must start with %q or %q", ErrInvalidRef, ref, RefPrefixDPAPI, RefPrefixEnv)
	}
}

// ValidBlobID reports whether s is a safe single blob file name: 1-128 chars
// of [A-Za-z0-9._-], not "." or "..", no path separators (traversal guard -
// the id is joined directly under <data>\secrets\).
func ValidBlobID(s string) bool {
	if s == "" || len(s) > 128 {
		return false
	}
	if s == "." || s == ".." {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '.' || r == '_' || r == '-':
		default:
			return false
		}
	}
	return true
}

// NewRef mints a fresh random dpapi reference ("dpapi:<hex>"). The id is
// random so that refs do not leak which provider or config table they belong
// to; the blob file is still exactly one file per ref (SPEC-02 §6).
func NewRef() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand failing is an unrecoverable runtime fault; the store
		// must never fall back to predictable ids.
		panic("secret: crypto/rand unavailable: " + err.Error())
	}
	return RefPrefixDPAPI + hex.EncodeToString(b[:])
}
