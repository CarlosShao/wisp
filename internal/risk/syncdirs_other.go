//go:build !windows

package risk

// registryProbe is a Windows-only P12 probe. macOS (S8) will need its own
// client-config probing (~/Library/CloudStorage etc.); until then the
// portable probes (defaults/config/fixture) run and an empty result keeps the
// fail-closed suspect fallback active.
//
// DEFERRED(P12-macos): sync-client probing for the macOS port (ticket 55).
func registryProbe(e probeEnv) []SyncRoot { return nil }
