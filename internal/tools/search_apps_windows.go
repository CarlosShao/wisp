//go:build windows

package tools

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// init registers the Windows enumeration source. On other platforms
// platformAppSources stays nil and the PATH floor alone runs.
func init() {
	platformAppSources = append(platformAppSources, appPathsRegistrySource{})
}

// search_apps_windows.go is the Windows enumeration half of D34 `search.apps`.
//
// It reads exactly one thing: the sub-key NAMES of the two App Paths hive
// locations the shell itself consults (HKLM then HKCU). That source is chosen
// because it is the SAME hive `windowsAppPathsResolver` (launcher_windows.go,
// ticket 22, untouched) resolves app.launch names against - so anything listed
// here is provably consumable by app.launch, which is ticket 24 AC#4's actual
// claim. Nothing in this file launches, writes, or widens a target shape: the
// listing is names, and the resolver downstream decides whether one is real.
//
// Registry value DATA is not trusted here: the values can carry arguments and
// quoting, and this tool must not become a second path-decision site. The
// resolved absolute executable comes back through the shared resolver, on the
// launch side's own rules.

// appPathsRegistrySource enumerates installed Win32 apps registered with the
// shell (the shape "打开记事本" answers to).
type appPathsRegistrySource struct{}

// Label implements AppSource.
func (appPathsRegistrySource) Label() string { return "App Paths 注册表" }

// appsPathsKey is the hive path the shell uses for the "run by name" table.
const appsPathsKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths`

// Apps implements AppSource. Both hives missing is an error (an enumeration
// that could not read anything is not "no apps installed"); one hive missing is
// normal and simply contributes nothing.
func (appPathsRegistrySource) Apps(ctx context.Context) ([]InstalledAppCandidate, error) {
	var (
		out   []InstalledAppCandidate
		seen  = make(map[string]bool)
		holes int
	)
	for _, hive := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		k, err := registry.OpenKey(hive, appsPathsKey, registry.READ)
		if err != nil {
			holes++
			continue
		}
		names, err := k.ReadSubKeyNames(-1)
		_ = k.Close()
		if err != nil {
			holes++
			continue
		}
		for _, raw := range names {
			if len(out) >= defaultMaxAppCandidates {
				return out, nil
			}
			name := strings.TrimSpace(raw)
			// The same shape veto the resolver applies (no path characters, no
			// wildcards): a key that is not a bare name cannot be one either.
			if name == "" || strings.ContainsAny(name, `/\\*?"<>|`) {
				continue
			}
			key := strings.ToLower(name)
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, InstalledAppCandidate{Name: name, Source: "App Paths 注册表"})
		}
	}
	if holes == 2 {
		return out, errNoAppPathsHive
	}
	return out, nil
}

// errNoAppPathsHive says the enumeration had nothing to read at all, so the
// caller reports a failed search instead of an empty app list.
var errNoAppPathsHive = errors.New("App Paths 注册表不可读")
