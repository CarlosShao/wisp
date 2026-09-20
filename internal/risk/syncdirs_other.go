//go:build !windows

package risk

// pathHandleVerify is false here: pathresolver_other.go has no handle-based
// resolution yet, so C26 can only offer lexical forms. Sync-root membership
// therefore stays best-effort on these platforms (documented in
// docs/PRECHECK.md §P12 alongside the DEFERRED marker below); failing every
// candidate closed would flag every write, not just the risky ones.
const pathHandleVerify = false

// registryProbe is a Windows-only P12 probe. macOS (S8) will need its own
// client-config probing (~/Library/CloudStorage etc.); until then the
// portable probes (env/defaults/config/fixture) run and an empty result keeps
// the fail-closed suspect fallback active.
//
// DEFERRED(P12-macos): sync-client probing for the macOS port (ticket 55) —
// acceptance needs a CloudStorage/~/Library config probe whose AC is written
// into ticket 55 (see its Constraints addendum).
func registryProbe(e probeEnv) []SyncRoot { return nil }
