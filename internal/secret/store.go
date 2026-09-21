package secret

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// ErrPortableDecrypt is returned when a dpapi: ref cannot be decrypted under
// portable mode (P13, SPEC-02 §6): portable installs move between machines
// and users, and a DPAPI blob created elsewhere can never be decrypted here.
// The remedy is an env: reference - plaintext is never used as a fallback.
var ErrPortableDecrypt = errors.New(
	"secret: portable mode: DPAPI blob cannot be decrypted (created by another user or on another machine); " +
		"store this key as an env: reference instead (P13: there is no plaintext fallback)")

// Store is the DPAPI-backed secret store (C28): one blob file per ref under
// <dataDir>\secrets\ (SPEC-02 §6 layout). Create with NewStore.
type Store struct {
	dir      string
	portable bool
}

// Option adjusts a Store.
type Option func(*Store)

// WithPortable marks the store as running under portable mode (P13): dpapi
// decrypt failures then surface as ErrPortableDecrypt with explicit guidance
// to env: refs.
func WithPortable(p bool) Option {
	return func(s *Store) { s.portable = p }
}

// NewStore creates (0700 semantics) the secrets directory under dataDir.
// dataDir is the per-env data root resolved by proc (the caller owns the
// env-fork policy).
func NewStore(dataDir string, opts ...Option) (*Store, error) {
	s := &Store{dir: filepath.Join(dataDir, "secrets")}
	for _, o := range opts {
		o(s)
	}
	// DPAPI user scope is the crypto boundary, but the blob file is not the
	// key: whoever can read it can still try to decrypt it, and on Windows
	// MkdirAll's 0o700 never restricted anybody (ticket 89).
	if err := winsec.PrivateDirAll(s.dir, 0o700); err != nil {
		return nil, fmt.Errorf("secret: create %s: %w", s.dir, err)
	}
	return s, nil
}

// Dir returns the blob directory. File names are blob ids (refs), not secret
// material; diagnostics may print them.
func (s *Store) Dir() string { return s.dir }

// Store encrypts secret with DPAPI (CurrentUser scope) and writes it as the
// blob file for ref ("dpapi:<blob-id>"), creating or overwriting it. env:
// refs are rejected: their values live in the process environment, there is
// nothing to persist. On Windows the DPAPI user scope is the access control;
// the 0600 mode is best-effort semantics for non-Windows tooling.
func (s *Store) Store(ref, secret string) error {
	kind, id, err := ParseRef(ref)
	if err != nil {
		return err
	}
	if kind != RefKindDPAPI {
		return fmt.Errorf("secret: store: %q: env refs live in the process environment; only dpapi: refs are persisted", ref)
	}
	if secret == "" {
		return errors.New("secret: store: refusing to store an empty secret")
	}
	blob, err := protect([]byte(secret))
	if err != nil {
		return err
	}
	if err := winsec.PrivateFile(s.blobPath(id), blob, 0o600); err != nil {
		return fmt.Errorf("secret: store: write blob: %w", err)
	}
	return nil
}

// Resolve returns the secret a ref names (SPEC-03 §5.3). "env:NAME" reads the
// process environment (missing or empty -> explicit error, never a default);
// "dpapi:<blob-id>" decrypts the blob file. Under portable mode a decrypt
// failure returns ErrPortableDecrypt (P13) - never a fallback value.
func (s *Store) Resolve(ref string) (string, error) {
	kind, value, err := ParseRef(ref)
	if err != nil {
		return "", err
	}
	switch kind {
	case RefKindEnv:
		v, ok := os.LookupEnv(value)
		if !ok {
			return "", fmt.Errorf("secret: resolve %q: environment variable %s is not set", ref, value)
		}
		if v == "" {
			return "", fmt.Errorf("secret: resolve %q: environment variable %s is empty", ref, value)
		}
		return v, nil
	case RefKindDPAPI:
		blob, err := os.ReadFile(s.blobPath(value))
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("secret: resolve %q: no blob file under %s", ref, s.dir)
		}
		if err != nil {
			return "", fmt.Errorf("secret: resolve %q: read blob: %w", ref, err)
		}
		plain, err := unprotect(blob)
		if err != nil {
			if s.portable {
				return "", fmt.Errorf("secret: resolve %q: %w", ref, ErrPortableDecrypt)
			}
			return "", fmt.Errorf("secret: resolve %q: %w", ref, err)
		}
		return string(plain), nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidRef, ref)
	}
}

// Exists reports whether a dpapi: ref has a blob file in this store. env: refs
// are rejected (they name process state, not a file). It exists so a caller
// can tell "not stored yet" from "unreadable" without building the blob path
// itself - blobPath stays the single traversal choke point.
func (s *Store) Exists(ref string) (bool, error) {
	kind, id, err := ParseRef(ref)
	if err != nil {
		return false, err
	}
	if kind != RefKindDPAPI {
		return false, fmt.Errorf("secret: exists: %q: env refs live in the process environment; only dpapi: refs have blob files", ref)
	}
	switch _, err := os.Stat(s.blobPath(id)); {
	case err == nil:
		return true, nil
	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	default:
		return false, fmt.Errorf("secret: exists %q: %w", ref, err)
	}
}

// Delete removes the blob file a dpapi: ref names (one file per ref, SPEC-02
// §6 - deleting the ref's file is the whole of "forget this secret"). An
// absent blob reports fs.ErrNotExist wrapped, so callers can distinguish
// "already gone" from an I/O failure; env: refs are rejected because there is
// nothing on disk to remove. The error text carries the ref and the directory,
// never secret material (C28).
func (s *Store) Delete(ref string) error {
	kind, id, err := ParseRef(ref)
	if err != nil {
		return err
	}
	if kind != RefKindDPAPI {
		return fmt.Errorf("secret: delete: %q: env refs live in the process environment; there is no blob to delete", ref)
	}
	if err := os.Remove(s.blobPath(id)); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("secret: delete %q: no blob file under %s: %w", ref, s.dir, err)
		}
		return fmt.Errorf("secret: delete %q: %w", ref, err)
	}
	return nil
}

// blobPath joins the id under the secrets dir. ParseRef already validated the
// charset; this is the choke point that keeps traversal out of file paths.
func (s *Store) blobPath(id string) string {
	if !ValidBlobID(id) {
		// Unreachable through the public API; a hard stop because this string
		// becomes a file path.
		panic(fmt.Sprintf("secret: unsafe blob id %q reached blobPath", id))
	}
	return filepath.Join(s.dir, id)
}
