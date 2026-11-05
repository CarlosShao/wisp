//go:build !windows

package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
)

// Non-Windows protector: the other half of protect.go's seam.
//
// WHAT THIS FILE IS FOR. Until this seam, every protect/unprotect call on a
// non-Windows build returned errDPAPIUnavailable ("DPAPI is only available on
// Windows; use an env: ref on this platform"). That was a deliberate CI shim
// (ticket 08): the package only had to COMPILE on ubuntu so its importers could
// run there. The cost was that the C28 store itself was unreachable on the
// POSIX leg - one file per ref, the no-plaintext invariant, the portable-mode
// P13 error, and every command that writes a credential - so ten cmd/wisp
// cases could only ever fail there.
//
// WHAT IT IS. An AEAD seal (AES-256-GCM, random 96-bit nonce) whose key is
// derived from the machine name and the store directory. The store therefore
// keeps its observable contract on any platform: one opaque blob per ref, no
// plaintext on disk, an explicit error when a blob cannot be opened (never a
// fallback value), and a blob bound to the machine and store that wrote it.
//
// WHAT IT IS NOT, stated here so no later reader quotes it as more:
//   - NOT a confidentiality boundary. The key comes from inputs a local reader
//     can also read (host name, directory), because a non-Windows platform has
//     no DPAPI-equivalent user secret to derive from without a keyring service
//     or a key file - and a key file inside the data root is refused by the
//     store's own no-extra-files invariant (ticket 63 AC#2, ticket 248 AC#3).
//     Windows keeps DPAPI (C28) and is the only platform Wisp ships; D7 keeps
//     Linux out of the support matrix. This seal exists so that the Linux leg
//     is a COMPILE target that can be exercised, not so that a Linux install
//     can be called secure.
//   - NOT a second C28 mechanism. C28's text names Windows DPAPI and the macOS
//     Keychain (S8) for the supported platforms; it says nothing about Linux,
//     and this file changes no contract sentence, no ref spelling (C27) and no
//     error class (D37).

const (
	// posixSealVersion is the first byte of every blob this file writes. It is
	// there so a blob from somewhere else (another mechanism, a truncated file,
	// a plaintext fragment) is refused by name instead of reaching the AEAD.
	posixSealVersion = 1
	posixNonceLen    = 12
	posixKeyLen      = 32
	posixSealDomain  = "wisp/secret/posix-seal/v1"
)

// posixProtector seals under a key bound to (machine, store directory).
type posixProtector struct{ dir string }

// platformProtector is the non-Windows half of protect.go's seam.
func platformProtector(dir string) protector { return posixProtector{dir: dir} }

// protect seals plaintext into a self-describing blob: version, nonce,
// ciphertext with its GCM tag. An empty plaintext is refused for the same
// reason the Windows protector refuses it (the store also refuses it one layer
// up; this keeps the two mechanisms answering alike).
func (p posixProtector) protect(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, errors.New("secret: posix seal: empty plaintext")
	}
	aead, err := p.aead()
	if err != nil {
		return nil, err
	}
	out := make([]byte, 1+posixNonceLen+len(plaintext)+aead.Overhead())
	out[0] = posixSealVersion
	nonce := out[1 : 1+posixNonceLen]
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("secret: posix seal: nonce: %w", err)
	}
	aead.Seal(out[:1+posixNonceLen], nonce, plaintext, nil)
	return out, nil
}

// unprotect opens a blob written by protect. Every failure is an explicit
// error naming the shape, never a fallback: a corrupted file, a blob from
// another machine or a blob from another store all land here, and the store
// turns them into the P13 guidance under portable mode (unchanged).
func (p posixProtector) unprotect(blob []byte) ([]byte, error) {
	if len(blob) == 0 {
		return nil, errors.New("secret: posix seal: empty blob")
	}
	const headerLen = 1 + posixNonceLen
	if len(blob) < headerLen || blob[0] != posixSealVersion {
		return nil, fmt.Errorf("secret: posix seal: not a wisp seal (want version %d, %d-byte header, got %d bytes)",
			posixSealVersion, headerLen, len(blob))
	}
	aead, err := p.aead()
	if err != nil {
		return nil, err
	}
	plain, err := aead.Open(nil, blob[1:headerLen], blob[headerLen:], nil)
	if err != nil {
		return nil, errors.New("secret: posix seal: blob cannot be opened here " +
			"(made by another machine or another store, or the file is corrupted)")
	}
	return plain, nil
}

// aead builds the AEAD for this store's derived key. The key is computed per
// call rather than cached: it is one HKDF over two short strings, the store is
// not on a hot path, and a cached key would be one more piece of process state
// to reason about (D38).
func (p posixProtector) aead() (cipher.AEAD, error) {
	key, err := posixSealKey(p.dir)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("secret: posix seal: cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("secret: posix seal: gcm: %w", err)
	}
	return gcm, nil
}

// posixSealKey derives the seal key from the machine name and the store
// directory (HKDF-SHA256, domain-separated). Both inputs are stable for the
// life of a store and differ across machines and stores, which is the property
// the store's contract needs - not secrecy of the inputs (see the file header).
// The derivation deliberately reads no environment variable and writes no file:
// a `wisp secret set` must leave exactly one file behind (ticket 63 AC#2).
func posixSealKey(dir string) ([]byte, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("secret: posix seal: host name: %w", err)
	}
	key, err := hkdf.Key(sha256.New, []byte(host), []byte(posixSealDomain), dir, posixKeyLen)
	if err != nil {
		return nil, fmt.Errorf("secret: posix seal: derive key: %w", err)
	}
	return key, nil
}
