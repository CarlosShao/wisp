//go:build !windows

package risk

// DEFERRED(macOS/Linux): SPEC-06 §4 requires realpath + lstat symlink reject
// for non-Windows platforms. Until that ticket lands, handle-based resolution
// is unsupported here and Resolve falls back to the lexical pipeline, which
// still classifies (fail-closed for the blacklist) but does not detect
// reparse traversal on these platforms.
func resolveHandle(p string) (string, bool) { return "", false }

func reparseComponents(p string) []string { return nil }
