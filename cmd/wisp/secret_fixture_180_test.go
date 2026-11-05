package main

// 票 98 C1（orphan-red repair, ticket 180 leg）— ONE portable seam for the fake
// key the composition-root fixtures below put into config.toml.
//
// The red this closes: docs/evidence/s1/98-attribution.md §2 C1 (76 of that
// census' 81 ubuntu reds). Those hosts wrote `api_key_ref = "dpapi:<id>"` and
// then really called secret.Store(); internal/secret/protect_other.go is a
// DOCUMENTED fail-closed stub off Windows (C28 / SPEC-02 §6 — the platform rule
// belongs to tickets 06/08, not to this file), so on the POSIX leg the fixture's
// own t.Fatal fired before any case could measure anything. The key reference is
// the HOST's seam (it exists so `wisp run` can reach the mock LLM at all), not
// any affected case's subject.
//
// What is deliberately NOT done:
//   - no //go:build windows and no _windows_test.go rename: that would remove
//     the cases from the test-core denominator `scripts/portable-tests.sh
//     --scope=cli` runs, which ticket 98 AC#2 names as the forbidden shortcut;
//   - no assertion, expectation or deadline is touched here or at any call site;
//   - Windows keeps today's bytes exactly: same blob id, same DPAPI Store call,
//     same on-disk shape, same value.
//
// On every other platform the fixture uses the OTHER frozen ref kind —
// secret.RefPrefixEnv, which internal/secret/refs.go documents for precisely
// this ("env:NAME names a process environment variable (CI/test: the value may
// be any placeholder ...)"). The value is always the suite's single fake key
// (run_test.go:53), so nothing real ever reaches an environment variable.
//
// 撤销：at each call site replace `keyRef := fixtureKeyRef180(...)` with the
// host's former inline secret.NewStore + st.Store pair, put the literal ref back
// into the config string, and delete this file.

import (
	"runtime"
	"testing"

	"github.com/CarlosShao/wisp/internal/secret"
)

// fixtureKeyEnv180 names the process variable carrying the suite's fake key on
// platforms without DPAPI. One fixed name, not one per case: cmd/wisp cases run
// sequentially and t.Setenv restores whatever was there before.
const fixtureKeyEnv180 = "WISP_TEST_FIXTURE_KEY"

// fixtureKeyRef180 seeds key for blob id blobID and returns the api_key_ref the
// fixture's config.toml must carry.
//
// On Windows the ref is the same dpapi: blob the fixtures wrote before, so the
// DPAPI path keeps its production coverage there (secret.Store holds no handle
// to release: its reads and writes open the blob per call, unlike memory.Store,
// which owns the SQLite file). Elsewhere it is `env:` plus the variable set for
// this case - an env ref has no blob to name, and secret.Store rejects env refs
// outright (internal/secret/store.go), which is why the seeding branches here
// instead of being one call.
func fixtureKeyRef180(t *testing.T, dataDir, blobID, key string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		st, err := secret.NewStore(dataDir)
		if err != nil {
			t.Fatal(err)
		}
		ref := secret.RefPrefixDPAPI + blobID
		if err := st.Store(ref, key); err != nil {
			t.Fatal(err)
		}
		return ref
	}
	t.Setenv(fixtureKeyEnv180, key)
	return secret.RefPrefixEnv + fixtureKeyEnv180
}
