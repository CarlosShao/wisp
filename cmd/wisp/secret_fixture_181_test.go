package main

// 票 98 C1（orphan-red repair, ticket 181 leg）— the SAME seam as
// secret_fixture_180_test.go, for the fixtures that leg left behind. The roster
// and the reasoning are in docs/evidence/s1/orphan-reds-cmdwisp-180.md §5; this
// file only adds the two shapes those call sites need, and does not restate the
// attribution.
//
// What is deliberately NOT done (identical to 180):
//   - no //go:build windows and no _windows_test.go rename: that would remove
//     the cases from the CLI denominator `scripts/portable-tests.sh --scope=cli`
//     runs, which ticket 98 AC#2 names as the forbidden shortcut;
//   - no t.Skip, no assertion, expectation or deadline touched here or at any
//     call site: only the ORIGIN of the ref moves;
//   - Windows keeps today's bytes exactly: same ref string, same blob id, same
//     DPAPI Store call, same on-disk shape.
//
// On every other platform the fixture uses the OTHER frozen ref kind,
// secret.RefPrefixEnv, which internal/secret/refs.go documents for exactly this
// ("env:NAME names a process environment variable (CI/test: the value may be any
// placeholder ...)"). The value is always the suite's single fake key
// (run_test.go:53), so nothing real ever reaches an environment variable.
//
// 撤销：at each call site replace `keyRef := fixtureKeyRef181(...)` with the
// host's former inline secret.NewStore + st.Store pair, put the literal ref back
// into the config string, replace `fixtureKeyDrop181(...)` with the former
// NewStore + st.Delete pair, and delete this file.

import (
	"os"
	"runtime"
	"testing"

	"github.com/CarlosShao/wisp/internal/secret"
)

// fixtureKeyRef181 seeds key under the dpapiRef this fixture's config.toml used
// to hard-code, and returns the api_key_ref that config must now carry.
//
// This is the "the fixture already spells the whole ref" variant of
// fixtureKeyRef180 (which takes a blob id). Several hosts write the ref twice -
// once into the config body, once into the Store call - so the ref is passed
// through unchanged instead of being rebuilt: on Windows the bytes written to
// config.toml and to the blob are the ones those hosts wrote before, and the
// returned string is the argument itself.
func fixtureKeyRef181(t *testing.T, dataDir, dpapiRef, key string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		st, err := secret.NewStore(dataDir)
		if err != nil {
			t.Fatal(err)
		}
		if err := st.Store(dpapiRef, key); err != nil {
			t.Fatal(err)
		}
		return dpapiRef
	}
	t.Setenv(fixtureKeyEnv180, key)
	return secret.RefPrefixEnv + fixtureKeyEnv180
}

// fixtureKeyDrop181 retires the key fixtureKeyRef181 planted, so a case can
// measure "this ref now names nothing" on either platform.
//
// Windows deletes the blob the ref names, which is what those cases did inline.
// Elsewhere the ref names an environment variable, so the same premise is "the
// variable is gone": os.Unsetenv plus a restore, rather than an empty value,
// because Resolve words "unset" and "empty" as two different errors and the
// case under test is about the Unconfigured path, not about either sentence.
func fixtureKeyDrop181(t *testing.T, dataDir, ref string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		st, err := secret.NewStore(dataDir)
		if err != nil {
			t.Fatal(err)
		}
		if err := st.Delete(ref); err != nil {
			t.Fatal(err)
		}
		return
	}
	old, had := os.LookupEnv(fixtureKeyEnv180)
	if err := os.Unsetenv(fixtureKeyEnv180); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if had {
			_ = os.Setenv(fixtureKeyEnv180, old)
			return
		}
		_ = os.Unsetenv(fixtureKeyEnv180)
	})
}
