//go:build windows

package risk

import (
	"golang.org/x/sys/windows/registry"
)

// WOW64 access view masks. golang.org/x/sys/windows/registry does not export
// the KEY_WOW64_* view flags (winreg.h: KEY_WOW64_64KEY=0x0100,
// KEY_WOW64_32KEY=0x0200), so we define them. We enumerate both views because
// 32/64-bit installers spell the same client records differently.
const (
	keyWow6464Key = 0x0100
	keyWow6432Key = 0x0200
)

// pathHandleVerify reports that C26 can OS-verify paths here
// (GetFinalPathNameByHandle), so an existing-but-unverifiable write target
// spelling is a real anomaly and fails closed as sync-suspect.
const pathHandleVerify = true

// registryProbe runs the portable P12 registry logic (syncdirs.go
// registryProbeFor) against HKCU. Reading with explicit access views matters:
// 32/64-bit installers spell the same client records differently, and
// golang.org/x/sys/windows/registry does not export the KEY_WOW64_* masks.
// Any probe error yields no roots, which keeps the caller on the fail-closed
// suspect-fallback path.
func registryProbe(e probeEnv) []SyncRoot { return registryProbeFor(liveRegistry{}) }

// liveRegistry is the registryValueSource over the real HKCU hive.
type liveRegistry struct{}

// subKeys enumerates both WOW64 views and returns the first non-empty list
// (a 32-bit-only installer record is still a configured root).
func (liveRegistry) subKeys(key string) []string {
	for _, view := range []uint32{keyWow6464Key, keyWow6432Key} {
		k, err := registry.OpenKey(registry.CURRENT_USER, key, registry.READ|view)
		if err != nil {
			continue
		}
		names, err := k.ReadSubKeyNames(-1)
		k.Close()
		if err == nil && len(names) > 0 {
			return names
		}
	}
	return nil
}

func (liveRegistry) string(key, name string) string {
	for _, view := range []uint32{keyWow6464Key, keyWow6432Key} {
		if v := regStringView(registry.CURRENT_USER, key, name, view); v != "" {
			return v
		}
	}
	return ""
}

func regStringView(root registry.Key, path, name string, view uint32) string {
	k, err := registry.OpenKey(root, path, registry.READ|view)
	if err != nil {
		return ""
	}
	defer k.Close()
	v, _, err := k.GetStringValue(name)
	if err != nil {
		return ""
	}
	return v
}
