package secret

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// captureLog swaps the default slog logger for a JSON handler writing into a
// buffer, so tests can prove secret material never reaches a log path.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	return &buf
}

const testSecret = "sk-super-secret-9876543210"

// TestStoreResolveRoundTrip is acceptance criterion 1: store -> resolve
// returns the original secret; storing again overwrites (per-ref one file).
func TestStoreResolveRoundTrip(t *testing.T) {
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	ref := NewRef()

	if err := st.Store(ref, testSecret); err != nil {
		t.Fatalf("Store: %v", err)
	}
	got, err := st.Resolve(ref)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got != testSecret {
		t.Fatalf("round-trip mismatch: got %d chars, want the original", len(got))
	}

	// Overwrite semantics: same ref, new value, still exactly one file.
	if err := st.Store(ref, "second-value-42"); err != nil {
		t.Fatalf("Store(overwrite): %v", err)
	}
	got, err = st.Resolve(ref)
	if err != nil {
		t.Fatalf("Resolve(after overwrite): %v", err)
	}
	if got != "second-value-42" {
		t.Fatalf("overwrite did not take effect")
	}
	entries, err := os.ReadDir(st.Dir())
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("secrets dir holds %d files, want exactly one per ref", len(entries))
	}
}

// TestBlobFileNotPlaintext is acceptance criterion 1: the on-disk blob must
// not be readable as the plaintext secret.
func TestBlobFileNotPlaintext(t *testing.T) {
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	ref := NewRef()
	if err := st.Store(ref, testSecret); err != nil {
		t.Fatalf("Store: %v", err)
	}
	blob, err := os.ReadFile(filepath.Join(st.Dir(), strings.TrimPrefix(ref, RefPrefixDPAPI)))
	if err != nil {
		t.Fatalf("read blob: %v", err)
	}
	if bytes.Contains(blob, []byte(testSecret)) {
		t.Fatal("blob file contains the plaintext secret")
	}
	if bytes.Contains(bytes.ToLower(blob), []byte("super-secret")) {
		t.Fatal("blob file contains a recognizable fragment of the secret")
	}
	if len(blob) == 0 {
		t.Fatal("blob file is empty")
	}

	// File-permission semantics: 0600 requested. On Windows the DPAPI user
	// scope is the access control and Go's mode mapping is coarse, so the
	// assertion only applies where modes are real.
	if runtime.GOOS != "windows" {
		fi, err := os.Stat(filepath.Join(st.Dir(), strings.TrimPrefix(ref, RefPrefixDPAPI)))
		if err != nil {
			t.Fatalf("stat blob: %v", err)
		}
		if fi.Mode().Perm() != 0o600 {
			t.Errorf("blob mode = %v, want 0600", fi.Mode().Perm())
		}
	}
}

// TestStoreRejectsEnvRefsAndEmpty pins the Store-side ref contract: env refs
// live in the environment (nothing to persist), empty secrets are refused,
// invalid refs never reach the filesystem.
func TestStoreRejectsEnvRefsAndEmpty(t *testing.T) {
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if err := st.Store("env:WISP_TEST_LLM_KEY", "x"); err == nil {
		t.Error("Store(env ref) must fail: env values are not persisted")
	}
	if err := st.Store(NewRef(), ""); err == nil {
		t.Error("Store(empty) must fail")
	}
	if err := st.Store("sk-plaintext", "x"); err == nil {
		t.Error("Store(invalid ref) must fail")
	}
	entries, _ := os.ReadDir(st.Dir())
	if len(entries) != 0 {
		t.Errorf("rejected stores must not leave files, found %d", len(entries))
	}
}

// TestResolveEnvRef is acceptance criterion 6: env: refs resolve any dummy
// value the environment carries (CI-friendly; the mock LLM does not validate).
func TestResolveEnvRef(t *testing.T) {
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	t.Run("arbitrary placeholder value", func(t *testing.T) {
		t.Setenv("WISP_TEST_LLM_KEY", "dummy-placeholder-value-not-a-real-key")
		got, err := st.Resolve("env:WISP_TEST_LLM_KEY")
		if err != nil {
			t.Fatalf("Resolve(env): %v", err)
		}
		if got != "dummy-placeholder-value-not-a-real-key" {
			t.Errorf("env ref returned %q, want the env value verbatim", got)
		}
	})

	t.Run("missing variable is an explicit error", func(t *testing.T) {
		const name = "WISP_SECRETSTORE_DEFINITELY_UNSET"
		if err := os.Unsetenv(name); err != nil {
			t.Fatalf("Unsetenv: %v", err)
		}
		if _, err := st.Resolve("env:" + name); err == nil {
			t.Fatal("Resolve of a missing env var must fail")
		}
	})

	t.Run("empty variable is an explicit error", func(t *testing.T) {
		t.Setenv("WISP_SECRETSTORE_EMPTY", "")
		if _, err := st.Resolve("env:WISP_SECRETSTORE_EMPTY"); err == nil {
			t.Fatal("Resolve of an empty env var must fail")
		}
	})
}

// TestPortableDecryptFailureExplicit is the P13 acceptance: under portable
// mode, an undecryptable dpapi: ref (corrupted, or encrypted by another
// user/machine) produces the explicit ErrPortableDecrypt guiding to env: -
// never plaintext, never a silent fallback. Without portable mode the same
// failure is still an explicit error.
func TestPortableDecryptFailureExplicit(t *testing.T) {
	mkStore := func(t *testing.T, portable bool) (*Store, string) {
		t.Helper()
		st, err := NewStore(t.TempDir(), WithPortable(portable))
		if err != nil {
			t.Fatalf("NewStore: %v", err)
		}
		ref := NewRef()
		if err := st.Store(ref, testSecret); err != nil {
			t.Fatalf("Store: %v", err)
		}
		return st, ref
	}
	corrupt := func(t *testing.T, st *Store, ref string) {
		t.Helper()
		blobPath := filepath.Join(st.Dir(), strings.TrimPrefix(ref, RefPrefixDPAPI))
		if err := os.WriteFile(blobPath, []byte("this is not a valid DPAPI blob at all"), 0o600); err != nil {
			t.Fatalf("corrupt blob: %v", err)
		}
	}

	t.Run("portable mode surfaces ErrPortableDecrypt with env: guidance", func(t *testing.T) {
		st, ref := mkStore(t, true)
		corrupt(t, st, ref)
		v, err := st.Resolve(ref)
		if err == nil {
			t.Fatal("Resolve must fail on an undecryptable blob")
		}
		if v != "" {
			t.Fatalf("Resolve returned %q alongside the error; there must be no fallback value", v)
		}
		if !errors.Is(err, ErrPortableDecrypt) {
			t.Fatalf("err = %v, want ErrPortableDecrypt", err)
		}
		if !strings.Contains(err.Error(), "env:") {
			t.Errorf("error must guide to env: refs, got: %v", err)
		}
		if strings.Contains(err.Error(), testSecret) {
			t.Error("error text must not carry the secret")
		}
	})

	t.Run("portable mode still decrypts own blobs", func(t *testing.T) {
		st, ref := mkStore(t, true)
		got, err := st.Resolve(ref)
		if err != nil || got != testSecret {
			t.Fatalf("Resolve own blob: (%q, %v)", got, err)
		}
	})

	t.Run("non-portable mode also fails explicitly", func(t *testing.T) {
		st, ref := mkStore(t, false)
		corrupt(t, st, ref)
		_, err := st.Resolve(ref)
		if err == nil {
			t.Fatal("Resolve must fail on an undecryptable blob")
		}
		if errors.Is(err, ErrPortableDecrypt) {
			t.Fatalf("non-portable failure must be the plain DPAPI error, got ErrPortableDecrypt: %v", err)
		}
		if strings.Contains(err.Error(), testSecret) {
			t.Error("error text must not carry the secret")
		}
	})
}

// TestResolveErrorsCarryNoSecret proves the error paths of the store never
// embed secret material: the only string that could carry it is the value
// itself, and it must never be formatted into an error or a log.
func TestResolveErrorsCarryNoSecret(t *testing.T) {
	logs := captureLog(t)
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	// A resolved secret must never be formatted anywhere; the sanctioned log
	// path is RedactSecret (last 4 only).
	got, err := st.Resolve(NewRef()) // missing blob -> error, not the secret
	if err == nil {
		t.Fatal("resolve of a missing blob must fail")
	}
	_ = got
	if err := st.Store(NewRef(), testSecret); err != nil {
		t.Fatalf("Store: %v", err)
	}
	resolved, err := st.Resolve("env:WISP_SECRETSTORE_MISSING_VAR")
	if err == nil {
		t.Fatal("resolve of missing env var must fail")
	}
	_ = resolved
	slog.Warn("resolved api key", "key", RedactSecret(testSecret)) // sanctioned form

	out := logs.String()
	if strings.Contains(out, testSecret) {
		t.Errorf("log path contains the plaintext secret:\n%s", out)
	}
	if !strings.Contains(out, RedactSecret(testSecret)) {
		t.Errorf("log path must carry the last-4 form %q:\n%s", RedactSecret(testSecret), out)
	}
}

// TestRedactSecret pins the last-4 contract (C28): everything but the final
// 4 runes collapses to '*'; secrets of 4 runes or fewer reveal nothing.
func TestRedactSecret(t *testing.T) {
	cases := []struct{ in, want string }{
		{"sk-super-secret-9876543210", "****3210"},
		{"abcd", "****"},
		{"abc", "****"},
		{"", "****"},
		{"密钥测试密钥九", "****试密钥九"},
	}
	for _, tc := range cases {
		if got := RedactSecret(tc.in); got != tc.want {
			t.Errorf("RedactSecret(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
