package manifestslot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/tools"
)

// AC3 — the install leg's backend flow (ticket 50, SPEC-07 §6, PLAN.md §14.7):
//
//	copy/download -> C14 schema validate -> C16 host_api gate -> explicit
//	authorization (full capability list + source) -> write plugins\<id>\ +
//	plugin_state (pinned manifest hash) -> register -> list_tools visibility.
//
// The authorization UI itself belongs to the panel/frontend leg (machine
// owner's external agent). This file owns only the backend contract around it:
// an Authorizer must exist and its decision is load-bearing — no authorize, no
// write, no state, no registration.
//
// Scope notes:
//   - This file is the fresh-install half. The update flow re-runs this same
//     authorization flow against a live install (SPEC-07 §6), and the full
//     uninstall teardown (DisposalScope of running instances, ask-about-data)
//     lives in lifecycle.go — same package, same Authorizer contract.
//   - The on-disk layout is exactly SPEC-07 §6: <root>\<id>\manifest.toml.
//   - exe_hash is pinned for a manifest that wraps CLIs (D46, ticket 52):
//     before anything lands, PinManifestExe resolves each `exe` through C26,
//     refuses an unstable path, hashes the bytes and compares them with any
//     exe_hash the manifest declared. The resulting pins go into plugin_state
//     and into the live view, where every execution re-verifies them.

// InstalledRecord is what the install leg pins in plugin_state: the §14.7
// columns needed for the AC4 tamper gate (id/version/hash), plus the D46 pin.
type InstalledRecord struct {
	ID      string
	Version string
	Hash    string // "sha256:<hex>" over the exact manifest bytes
	// ExeHash is the D46 command-plugin slot (ticket 52): the pinned sha256 of
	// every wrapped CLI, keyed by tool name. Empty for a manifest without
	// command tools, and then nothing about the record changed (the value is
	// the same shape ticket 50 reserved).
	ExeHash string
}

// StateStore is the plugin_state seam the installer must record through.
// internal/memory owns the table; internal/plugin MUST NOT import it (import
// cycle via memory -> plugin, see the package doc in manifestslot.go), so the
// composed leg supplies the adapter.
type StateStore interface {
	SaveInstalled(ctx context.Context, rec InstalledRecord) error
	DeleteInstalled(ctx context.Context, id string) error
	// ListInstalled feeds Manager.LoadAll's restart path (SPEC-07 §6
	// load-time verification runs over exactly this list).
	ListInstalled(ctx context.Context) ([]InstalledRecord, error)
}

// InstallInfo is the full review sheet the Authorizer must see before
// anything lands: identity, source, and the union capability face.
type InstallInfo struct {
	PluginID     string
	Name         string
	Version      string
	Source       string // what the user is being asked to trust (path/URL)
	Capabilities []string
	Manifest     plugin.Manifest
	// PreviousVersion is non-empty only on an update's sheet: "update =
	// re-run install authorization" must tell the user what is being
	// replaced (SPEC-07 §6).
	PreviousVersion string
}

// Authorizer decides whether the install proceeds. Returning an error refuses
// the install with nothing written (explicit user authorize, SPEC-07 §6).
type Authorizer interface {
	Authorize(ctx context.Context, info InstallInfo) error
}

// Installer wires the AC3 backend flow. All fields are required.
type Installer struct {
	// PluginsRoot is the plugins\ directory; the id from the manifest becomes
	// its single path segment underneath (SPEC-07 §6).
	PluginsRoot string
	// Registry is the composed tools.Registry list_tools reads.
	Registry *tools.Registry
	// Store records the pinned install state (plugin_state.hash).
	Store StateStore
	// Authorizer is the backend hook the panel/frontend leg implements.
	Authorizer Authorizer
	// HostVersion is what the C16 gate checks against (buildinfo.Version).
	HostVersion string
	// Runner is the execution seam handed to the registered slot (tickets
	// 51/52 fill it; the install leg passes it through untouched).
	Runner Run
}

// InstallResult reports one completed install.
type InstallResult struct {
	Manifest plugin.Manifest
	Path     string // <PluginsRoot>\<id>\manifest.toml
	Hash     string // the pinned plugin_state.hash value
	// ExePins are the D46 pins by tool name (empty when the plugin wraps no
	// CLI). They are what every call re-verifies.
	ExePins map[string]plugin.ExePin
}

// Refusal sentinels for the install leg.
var (
	ErrNotAuthorized  = errors.New("manifestslot: install refused by the authorizer")
	ErrAlreadyInstall = errors.New("manifestslot: plugin already installed (update re-runs authorization in its own leg)")
	ErrToolNameTaken  = errors.New("manifestslot: a declared tool name is already registered (C1 global uniqueness)")
)

// Install runs the AC3 flow end to end and returns the manifest only after the
// slot is registered, i.e. after the tool is visible to list_tools. Every
// failure path leaves no file, no state record, and no registration behind.
func (in Installer) Install(ctx context.Context, raw []byte, source string) (InstallResult, error) {
	switch {
	case in.PluginsRoot == "":
		return InstallResult{}, errors.New("manifestslot: installer: PluginsRoot is required")
	case in.Registry == nil:
		return InstallResult{}, errors.New("manifestslot: installer: Registry is required")
	case in.Store == nil:
		return InstallResult{}, errors.New("manifestslot: installer: Store is required")
	case in.Authorizer == nil:
		return InstallResult{}, errors.New("manifestslot: installer: Authorizer is required (SPEC-07 §6: explicit user authorize, never auto-allow)")
	case in.HostVersion == "":
		return InstallResult{}, errors.New("manifestslot: installer: HostVersion is required")
	case in.Runner == nil:
		return InstallResult{}, errors.New("manifestslot: installer: Runner is required")
	}

	m, err := plugin.ParseManifest(raw)
	if err != nil {
		return InstallResult{}, fmt.Errorf("manifestslot: install refused by C14: %w", err)
	}
	if err := m.CompatibleWith(in.HostVersion); err != nil {
		return InstallResult{}, err
	}

	// C1 global uniqueness, checked before anything lands. The Registry itself
	// skips names it already has when a provider registers (internal/tools/registry.go),
	// so a plugin shadowing a builtin or another plugin's tool would half-
	// install silently — the install leg refuses the whole manifest instead.
	// This runs after the already-installed check so a reinstall of the same
	// plugin reports the real condition (its own names collide with itself)
	// instead of a misleading self-shadowing refusal.
	dir := filepath.Join(in.PluginsRoot, m.ID)
	path := filepath.Join(dir, "manifest.toml")
	if _, statErr := os.Stat(dir); statErr == nil {
		return InstallResult{}, fmt.Errorf("%w: %s", ErrAlreadyInstall, m.ID)
	} else if !os.IsNotExist(statErr) {
		return InstallResult{}, fmt.Errorf("manifestslot: inspect %s: %w", dir, statErr)
	}
	for _, t := range m.Tools {
		if _, taken := in.Registry.Lookup(t.Name); taken {
			return InstallResult{}, fmt.Errorf("%w: %q already registered (plugin %q)",
				ErrToolNameTaken, t.Name, m.ID)
		}
	}

	// D46: pin every wrapped executable BEFORE anything lands, so a command
	// tool whose exe is movable, missing, not an image, or not the bytes the
	// manifest named is refused with no file, no state and no registration
	// (ticket 52, PLAN.md §16.5.6 red-team case 1).
	exePins, err := plugin.PinManifestExe(m, nil)
	if err != nil {
		return InstallResult{}, fmt.Errorf("manifestslot: install refused by D46 pinning: %w", err)
	}
	exeHash := plugin.EncodeExePins(exePins)

	if err := in.Authorizer.Authorize(ctx, reviewSheet(m, source)); err != nil {
		return InstallResult{}, fmt.Errorf("%w: %v", ErrNotAuthorized, err)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return InstallResult{}, fmt.Errorf("manifestslot: write plugins\\%s: %w", m.ID, err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		os.RemoveAll(dir)
		return InstallResult{}, fmt.Errorf("manifestslot: write manifest %s: %w", path, err)
	}

	hash := plugin.ManifestHash(raw)
	// Re-read the bytes we just pinned: what plugin_state.hash vouches for must
	// be exactly what the AC4 load gate will re-hash.
	onDisk, err := os.ReadFile(path)
	if err != nil || plugin.VerifyBytes(onDisk, hash) != nil {
		in.rollback(ctx, dir, m.ID)
		return InstallResult{}, fmt.Errorf("manifestslot: installed manifest does not match its pinned hash (refused before registration)")
	}
	if err := in.Store.SaveInstalled(ctx, InstalledRecord{ID: m.ID, Version: m.Version, Hash: hash, ExeHash: exeHash}); err != nil {
		in.rollback(ctx, dir, m.ID)
		return InstallResult{}, fmt.Errorf("manifestslot: record plugin_state for %s: %w", m.ID, err)
	}

	provider, err := New(m, in.Runner)
	if err != nil {
		in.rollback(ctx, dir, m.ID)
		return InstallResult{}, fmt.Errorf("manifestslot: build registry slot for %s: %w", m.ID, err)
	}
	if err := in.Registry.RegisterProvider(provider); err != nil {
		in.rollback(ctx, dir, m.ID)
		return InstallResult{}, fmt.Errorf("manifestslot: register %s: %w", m.ID, err)
	}
	return InstallResult{Manifest: m, Path: path, Hash: hash, ExePins: exePins}, nil
}

// Uninstall is the state/file half of teardown: the record goes and the
// installed files go. Callers who own live plugin instances must use
// Manager.Uninstall (lifecycle.go), which additionally revokes calls and
// disposes the plugin's C11 DisposalScope before reaching this.
func (in Installer) Uninstall(ctx context.Context, id string) error {
	if in.Store == nil || in.PluginsRoot == "" {
		return errors.New("manifestslot: uninstaller: PluginsRoot and Store are required")
	}
	// The id charset was enforced by C14 at install; a caller-supplied id is
	// still re-checked here before it touches the filesystem.
	if id == "" || id != filepath.Base(id) {
		return fmt.Errorf("manifestslot: uninstall: unsafe plugin id %q", id)
	}
	if err := in.Store.DeleteInstalled(ctx, id); err != nil {
		return fmt.Errorf("manifestslot: delete plugin_state %s: %w", id, err)
	}
	return os.RemoveAll(filepath.Join(in.PluginsRoot, id))
}

func (in Installer) rollback(ctx context.Context, dir, id string) {
	_ = in.Store.DeleteInstalled(ctx, id)
	_ = os.RemoveAll(dir)
}

// reviewSheet renders the full capability list + source the user must approve
// (ticket 50 Key constraints: "show FULL capability list + source"). The face
// is the union over tools, deduped and sorted for a stable review sheet.
func reviewSheet(m plugin.Manifest, source string) InstallInfo {
	seen := map[string]bool{}
	var caps []string
	for _, t := range m.Tools {
		for _, c := range t.Capabilities {
			if !seen[c] {
				seen[c] = true
				caps = append(caps, c)
			}
		}
	}
	sort.Strings(caps)
	return InstallInfo{
		PluginID:     m.ID,
		Name:         m.Name,
		Version:      m.Version,
		Source:       source,
		Capabilities: caps,
		Manifest:     m,
	}
}
