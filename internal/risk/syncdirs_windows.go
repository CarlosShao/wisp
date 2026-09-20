//go:build windows

package risk

import (
	"strings"

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

// registryProbe reads the per-client account keys under HKCU (P12: the
// configured roots survive client relocation, which path-string matching
// alone cannot). Values are read with explicit access views because 32/64-bit
// hosts spell the same layout differently; any probe error yields no roots,
// which keeps the caller on the fail-closed suspect-fallback path.
func registryProbe(e probeEnv) []SyncRoot {
	var out []SyncRoot
	// OneDrive: HKCU\Software\Microsoft\OneDrive\Accounts\<Personal|Business#>\UserFolder
	out = append(out, oneDriveRoots()...)
	// Dropbox installer record: HKCU\Software\Dropbox\Dropbox "Path"
	if v := regString(registry.CURRENT_USER, `Software\Dropbox\Dropbox`, "Path"); v != "" {
		out = append(out, SyncRoot{Provider: "Dropbox", Path: v, Source: "registry"})
	}
	// Nutstore (坚果云): HKCU\Software\Nutstore "InstallPath" is the client
	// install dir, not the sync root; sync roots live in the (undocumented,
	// version-specific) client config -> P12 residual, covered by
	// default-location probe + suspect fallback. See docs/PRECHECK.md.
	return out
}

func oneDriveRoots() []SyncRoot {
	const accounts = `Software\Microsoft\OneDrive\Accounts`
	var out []SyncRoot
	for _, view := range []uint32{keyWow6464Key, keyWow6432Key} {
		k, err := registry.OpenKey(registry.CURRENT_USER, accounts, registry.READ|view)
		if err != nil {
			continue
		}
		subkeys, err := k.ReadSubKeyNames(-1)
		if err == nil {
			for _, sk := range subkeys {
				if !strings.EqualFold(sk, "Personal") && !strings.HasPrefix(sk, "Business") {
					continue // unknown nodes (e.g. "ConsumingAccounts") carry no UserFolder
				}
				acc, err := registry.OpenKey(k, sk, registry.READ)
				if err != nil {
					continue
				}
				if v, _, err := acc.GetStringValue("UserFolder"); err == nil && v != "" {
					out = append(out, SyncRoot{Provider: "OneDrive", Path: v, Source: "registry"})
				}
				acc.Close()
			}
		}
		k.Close()
		if len(out) > 0 {
			break
		}
	}
	return out
}

func regString(root registry.Key, path, name string) string {
	k, err := registry.OpenKey(root, path, registry.READ)
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
