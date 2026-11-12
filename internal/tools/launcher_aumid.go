package tools

import "regexp"

// windowsAUMIDPattern is the conservative shape gate for !Name!App ids
// (package family name + app id, no path characters). It lives in a
// platform-neutral file so EVERY GOOS classifies an AUMID-shaped app.launch
// target the same way (resolveApp order: direct path, AUMID shape, name
// resolver); what a platform DOES with an aumid target stays
// platform-specific — launcher_windows.go activates it, launcher_other.go
// answers the honest "Windows-only capability (D7)" refusal.
var windowsAUMIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{1,255}![A-Za-z0-9._-]{1,255}$`)
