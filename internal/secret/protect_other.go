//go:build !windows

package secret

import "errors"

// Non-Windows compilation shim (ticket 08, CI enablement; the C28 behavior
// itself is ticket 06 and unchanged on Windows).
//
// The SecretStore is DPAPI-backed (C28, SPEC-02 §6) and Wisp is
// Windows-only; the linux build only needs the package to COMPILE so that
// its importers (internal/config) can run in the ubuntu test-core job.
// The stub is fail-closed: every DPAPI operation returns an explicit error
// pointing at the env: ref kind - it never degrades to plaintext, matching
// the portable-mode rule (P13: undecryptable refs produce an explicit error,
// never a fallback).

var errDPAPIUnavailable = errors.New(
	"secret: DPAPI is only available on Windows; use an env: ref on this platform")

func protect(plaintext []byte) ([]byte, error) {
	return nil, errDPAPIUnavailable
}

func unprotect(blob []byte) ([]byte, error) {
	return nil, errDPAPIUnavailable
}
