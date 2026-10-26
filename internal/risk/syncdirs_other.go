//go:build !windows

package risk

// pathHandleVerify is false here: it is the alternative to the Windows
// "unverifiable existing ancestor fail-closes" rule (syncdirs.go), and ticket
// 347 did NOT change it. C26 now resolves paths on these platforms (realpath,
// pathresolver_other.go), so an existing directory normally comes back
// Resolved; what this flag still covers is the residue where the OS refuses to
// re-open a component a plain Lstat accepted (EACCES on a parent, a dangling
// link), and there the lexical form has always been the answer on POSIX.
// Turning it on is a separate decision with its own denominator - it tightens
// sync-root membership for that residue - so it stays as it was.
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
