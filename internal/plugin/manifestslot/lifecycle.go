package manifestslot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Manager is the SPEC-07 §6 lifecycle owner for Tier-1 manifest plugins:
// install -> restart-load -> update -> uninstall, with one C11 DisposalScope
// per live plugin.
//
// Contract mapping (ticket 50):
//   - Install reuses Installer (C14 -> C16 -> full-sheet authorize -> write
//     plugins\<id>\ + plugin_state(hash, exe_hash slot) -> register), and
//     keeps the live declaration set so later calls never run against a
//     stale manifest snapshot.
//   - LoadAll is the restart path: every plugin_state record re-reads its
//     manifest through the AC4 tamper gate (pinned hash) and the C16 gate;
//     refusals are reported per plugin, never silently skipped.
//   - Update re-runs the FULL install authorization flow against the live
//     id (SPEC-07 §6 "更新 = 重走安装授权流程"); a host_api the running host
//     cannot satisfy (major mismatch) is refused with the C16 error and the
//     old install keeps serving.
//   - Uninstall revokes calls, disposes the plugin scope (deterministic
//     teardown of running instances: cancel -> join tracked goroutines ->
//     deferred fns), deletes files + state record, and asks about data via
//     options with the default being KEEP (SPEC-07 §6 / PLAN.md §14.7).
//
// Known boundary, stated once and honestly: tools.Registry (ticket 20, not
// this leg's write surface) has no entry-removal API. After uninstall or a
// shrinking update, stale entry names stay visible in the directory until
// the composed host rebuilds the Registry from plugin_state (the restart
// path above is exactly that rebuild). What this leg guarantees is that a
// stale name can never EXECUTE: Manager.Run refuses it here and now.
type Manager struct {
	// PluginsRoot is the plugins\ directory (SPEC-02 §6 layout).
	PluginsRoot string
	// DataRoot is the optional root holding each plugin's mutable data as
	// <DataRoot>\<id>\ (the "ask-about-data" subject on uninstall). Empty
	// means no separate data tree exists for this host.
	DataRoot string
	// Registry is the composed tools.Registry list_tools reads.
	Registry *tools.Registry
	// Store is the plugin_state seam (hash pins live here).
	Store StateStore
	// Authorizer is the explicit-user-authorization hook (install, update).
	Authorizer Authorizer
	// HostVersion is what the C16 gate checks against.
	HostVersion string
	// Exec is the behaviour seam every declarative call is handed to (the
	// goja host of ticket 51 and the D46 runner of ticket 52 fill it; tests
	// fill a mock). Required: a slot that lists but cannot execute is the
	// failure mode manifestslot.New already refuses.
	Exec Run
	// Alerter receives tamper alerts on the restart path; nil = LogAlerter.
	Alerter plugin.Alerter
	// ScopeParent is the parent ctx of per-plugin DisposalScopes (nil =
	// Background). The C11 hierarchy puts plugin scopes under the task/
	// session roots the composed host owns.
	ScopeParent context.Context

	mu      sync.Mutex
	plugins map[string]*managedPlugin
}

// managedPlugin is one live install: the current declarations, the pinned
// hash, and the scope every running instance of it is tracked under.
type managedPlugin struct {
	manifest plugin.Manifest
	hash     string
	// exeHash is the D46 command-plugin slot (ticket 52 fills it); Tier-1
	// carries it forward untouched so an update never silently clears a pin.
	exeHash string
	defs    map[string]plugin.ToolDef
	scope   *plugin.DisposalScope
	removed bool
}

func (m *managedPlugin) live() bool { return m != nil && !m.removed }

// Lifecycle refusal sentinels.
var (
	ErrPluginNotInstalled = errors.New("manifestslot: update refused: plugin is not installed (install first)")
	ErrPluginNotLive      = errors.New("manifestslot: call refused: plugin is not installed or has been uninstalled")
	ErrToolNotDeclared    = errors.New("manifestslot: call refused: tool is no longer declared by the plugin's current manifest")
	ErrUpdateIDMismatch   = errors.New("manifestslot: update refused: manifest id does not match the plugin being updated")
	ErrUnsafePluginID     = errors.New("manifestslot: unsafe plugin id (must be a single safe path segment)")
)

// LoadVerdict is the per-plugin outcome of the restart path. Every record is
// reported; nothing is skipped silently.
type LoadVerdict struct {
	PluginID string
	Loaded   bool
	Reason   string
}

// UninstallOptions carries the "ask-about-data" answer. Zero value keeps
// plugin data (SPEC-07 §6 default).
type UninstallOptions struct {
	// DeletePluginData also removes <DataRoot>\<id>\ when set. The UI asks;
	// the default answer is keep.
	DeletePluginData bool
}

// UninstallReport is what one completed teardown did (and what it kept).
type UninstallReport struct {
	Disposal    plugin.DisposeResult
	DataDir     string
	DataDeleted bool
}

func (m *Manager) required() error {
	switch {
	case m.PluginsRoot == "":
		return errors.New("manifestslot: manager: PluginsRoot is required")
	case m.Registry == nil:
		return errors.New("manifestslot: manager: Registry is required")
	case m.Store == nil:
		return errors.New("manifestslot: manager: Store is required")
	case m.Authorizer == nil:
		return errors.New("manifestslot: manager: Authorizer is required (SPEC-07 §6: explicit user authorize, never auto-allow)")
	case m.HostVersion == "":
		return errors.New("manifestslot: manager: HostVersion is required")
	case m.Exec == nil:
		return errors.New("manifestslot: manager: Exec is required")
	}
	return nil
}

func (m *Manager) installer() Installer {
	return Installer{
		PluginsRoot: m.PluginsRoot,
		Registry:    m.Registry,
		Store:       m.Store,
		Authorizer:  m.Authorizer,
		HostVersion: m.HostVersion,
		Runner:      m.Run,
	}
}

func (m *Manager) alerter() plugin.Alerter {
	if m.Alerter == nil {
		return plugin.LogAlerter{}
	}
	return m.Alerter
}

// safePluginID re-checks the C14 charset gate for any id that touches the
// filesystem, so a caller-supplied id can never escape PluginsRoot/DataRoot.
// (risk.PathResolver owns real path decisions; this is the shape guard, same
// one Installer.Uninstall applies.)
func safePluginID(id string) error {
	if id == "" || id != filepath.Base(id) || id == "." || id == ".." {
		return fmt.Errorf("%w: %q", ErrUnsafePluginID, id)
	}
	return nil
}

// Run is the Manager's Run: the seam every registered slot executes through.
// It resolves the CURRENT declaration of the tool (so an updated manifest
// retires the old one's allowlist immediately), refuses anything that is not
// live, and hands the call a ctx bound to the plugin's DisposalScope, so
// uninstalling a plugin deterministically cancels its running instances.
func (m *Manager) Run(ctx context.Context, pluginID string, def plugin.ToolDef, params json.RawMessage) (tools.Result, error) {
	m.mu.Lock()
	mp := m.plugins[pluginID]
	m.mu.Unlock()
	if !mp.live() {
		return tools.Result{}, fmt.Errorf("%w: %s", ErrPluginNotLive, pluginID)
	}
	current, declared := mp.defs[def.Name]
	if !declared {
		return tools.Result{}, fmt.Errorf("%w: %q (plugin %q)", ErrToolNotDeclared, def.Name, pluginID)
	}

	callCtx, cancelCall := context.WithCancel(mp.scope.Ctx())
	defer cancelCall()
	done := make(chan struct{})
	defer close(done)
	// Tracked watcher: whichever side dies first (the caller's ctx or the
	// plugin scope being disposed) cancels the derived call ctx. Registered
	// through the scope, so Dispose joins it (no bare goroutine).
	mp.scope.Go("plugin-call-watch", func(scopeCtx context.Context) {
		select {
		case <-ctx.Done():
		case <-scopeCtx.Done():
		case <-done:
		}
		cancelCall()
	})
	return m.Exec(callCtx, pluginID, current, params)
}

// Scope exposes a live plugin's DisposalScope so an Exec implementation can
// register long-lived instances for deterministic teardown (SPEC-07 §6
// uninstall rule). It reports false once the plugin is uninstalled.
func (m *Manager) Scope(pluginID string) (*plugin.DisposalScope, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mp := m.plugins[pluginID]
	if !mp.live() {
		return nil, false
	}
	return mp.scope, true
}

// Install runs the AC3 flow through Installer and puts the result under a
// live DisposalScope.
func (m *Manager) Install(ctx context.Context, raw []byte, source string) (InstallResult, error) {
	if err := m.required(); err != nil {
		return InstallResult{}, err
	}
	res, err := m.installer().Install(ctx, raw, source)
	if err != nil {
		return InstallResult{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.plugins == nil {
		m.plugins = map[string]*managedPlugin{}
	}
	if _, dup := m.plugins[res.Manifest.ID]; dup {
		// Narrow race with a concurrent install/uninstall of the same id;
		// the file/state side already landed, the live view refuses to
		// double-register.
		return res, fmt.Errorf("manifestslot: plugin %q already has a live install", res.Manifest.ID)
	}
	m.plugins[res.Manifest.ID] = &managedPlugin{
		manifest: res.Manifest,
		hash:     res.Hash,
		defs:     defsOf(res.Manifest),
		scope:    m.newScope(res.Manifest.ID),
	}
	return res, nil
}

func (m *Manager) newScope(pluginID string) *plugin.DisposalScope {
	parent := m.ScopeParent
	if parent == nil {
		parent = context.Background()
	}
	return plugin.NewDisposalScope("plugin:"+pluginID, parent)
}

func defsOf(manifest plugin.Manifest) map[string]plugin.ToolDef {
	defs := make(map[string]plugin.ToolDef, len(manifest.Tools))
	for _, t := range manifest.Tools {
		defs[t.Name] = t
	}
	return defs
}

// LoadAll is the restart path over plugin_state: every installed record
// re-reads its manifest through the AC4 tamper gate, then the C16 gate, then
// registers its slot. Each outcome is reported by name; a refused plugin
// never blocks the others, and none of this is silent.
func (m *Manager) LoadAll(ctx context.Context) ([]LoadVerdict, error) {
	if err := m.required(); err != nil {
		return nil, err
	}
	recs, err := m.Store.ListInstalled(ctx)
	if err != nil {
		return nil, fmt.Errorf("manifestslot: load: read plugin_state: %w", err)
	}
	verdicts := make([]LoadVerdict, 0, len(recs))
	for _, rec := range recs {
		v := LoadVerdict{PluginID: rec.ID}
		if err := safePluginID(rec.ID); err != nil {
			v.Reason = err.Error()
			verdicts = append(verdicts, v)
			continue
		}
		path := filepath.Join(m.PluginsRoot, rec.ID, "manifest.toml")
		manifest, err := plugin.LoadInstalled(path, plugin.LoadOptions{
			PinnedHash: rec.Hash,
			PluginID:   rec.ID,
			Alerter:    m.alerter(),
		})
		if err != nil {
			v.Reason = fmt.Sprintf("tamper/C14 gate: %v", err)
			verdicts = append(verdicts, v)
			continue
		}
		if err := manifest.CompatibleWith(m.HostVersion); err != nil {
			v.Reason = err.Error()
			verdicts = append(verdicts, v)
			continue
		}
		provider, err := New(manifest, m.Run)
		if err != nil {
			v.Reason = fmt.Sprintf("slot build: %v", err)
			verdicts = append(verdicts, v)
			continue
		}
		if err := m.Registry.RegisterProvider(provider); err != nil {
			v.Reason = fmt.Sprintf("registry: %v", err)
			verdicts = append(verdicts, v)
			continue
		}
		m.mu.Lock()
		if m.plugins == nil {
			m.plugins = map[string]*managedPlugin{}
		}
		m.plugins[rec.ID] = &managedPlugin{
			manifest: manifest,
			hash:     rec.Hash,
			exeHash:  rec.ExeHash,
			defs:     defsOf(manifest),
			scope:    m.newScope(rec.ID),
		}
		m.mu.Unlock()
		v.Loaded = true
		v.Reason = "loaded"
		verdicts = append(verdicts, v)
	}
	return verdicts, nil
}

// UpdateResult carries the outcome of an update decision.
type UpdateResult struct {
	InstallResult
	NoChange bool // bytes identical to the pinned install
}

// Update re-runs the full install authorization flow against a live plugin
// (SPEC-07 §6). Order, and what each refusal preserves:
//
//	C14 parse      -> bad new manifest: refused, old install untouched
//	C16 gate       -> host cannot satisfy the NEW host_api (major mismatch):
//	                  refused with the explicit C16 error, old install untouched
//	authorize      -> declined: refused, old bytes/state/registrations intact
//	write + verify -> disk mismatch: old bytes restored, old pin intact
//	pin + swap     -> new hash pinned, live declarations swapped atomically
//
// Tool names the new manifest dropped stay visible in the Registry until the
// host rebuilds it (boundary note above) but are refused at Run; tool names
// it adds register now. Capability-FACE changes on a same-name tool take
// effect at execution (Run uses the current def); the bridge's registration
// snapshot updates on host restart.
func (m *Manager) Update(ctx context.Context, id string, raw []byte, source string) (UpdateResult, error) {
	if err := m.required(); err != nil {
		return UpdateResult{}, err
	}
	if err := safePluginID(id); err != nil {
		return UpdateResult{}, err
	}
	m.mu.Lock()
	mp := m.plugins[id]
	m.mu.Unlock()
	if !mp.live() {
		return UpdateResult{}, fmt.Errorf("%w: %s", ErrPluginNotInstalled, id)
	}

	manifest, err := plugin.ParseManifest(raw)
	if err != nil {
		return UpdateResult{}, fmt.Errorf("manifestslot: update refused by C14: %w", err)
	}
	if manifest.ID != id {
		return UpdateResult{}, fmt.Errorf("%w: manifest declares %q, updating %q", ErrUpdateIDMismatch, manifest.ID, id)
	}
	if err := manifest.CompatibleWith(m.HostVersion); err != nil {
		return UpdateResult{}, err // explicit C16 refusal; the old install keeps serving
	}

	newHash := plugin.ManifestHash(raw)
	path := filepath.Join(m.PluginsRoot, id, "manifest.toml")
	if newHash == mp.hash {
		return UpdateResult{
			InstallResult: InstallResult{Manifest: manifest, Path: path, Hash: mp.hash},
			NoChange:      true,
		}, nil
	}

	info := reviewSheet(manifest, source)
	info.PreviousVersion = mp.manifest.Version
	if err := m.Authorizer.Authorize(ctx, info); err != nil {
		return UpdateResult{}, fmt.Errorf("%w: update of %s declined: %v", ErrNotAuthorized, id, err)
	}

	oldBytes, err := os.ReadFile(path)
	if err != nil {
		return UpdateResult{}, fmt.Errorf("manifestslot: update: read current manifest %s: %w", path, err)
	}
	restore := func() {
		if err := os.WriteFile(path, oldBytes, 0o644); err != nil {
			// The old install may now be half-written; say so loudly rather
			// than let a later load discover it silently.
			slog.Error("manifestslot: update rollback failed; installed manifest may be inconsistent",
				"path", path, "err", err)
		}
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return UpdateResult{}, fmt.Errorf("manifestslot: update: write %s: %w", path, err)
	}
	onDisk, err := os.ReadFile(path)
	if err != nil || plugin.VerifyBytes(onDisk, newHash) != nil {
		restore()
		return UpdateResult{}, fmt.Errorf("manifestslot: update refused: installed manifest does not match its new hash (old bytes restored)")
	}
	if err := m.Store.SaveInstalled(ctx, InstalledRecord{
		ID: id, Version: manifest.Version, Hash: newHash, ExeHash: mp.exeHash,
	}); err != nil {
		restore()
		return UpdateResult{}, fmt.Errorf("manifestslot: update: re-pin plugin_state for %s: %w", id, err)
	}

	m.mu.Lock()
	mp.manifest = manifest
	mp.hash = newHash
	mp.defs = defsOf(manifest)
	m.mu.Unlock()

	// Land any tool names the new manifest ADDS (duplicate names are skipped
	// by the Registry; the stale-name boundary is documented above).
	provider, err := New(manifest, m.Run)
	if err != nil {
		return UpdateResult{}, fmt.Errorf("manifestslot: update: slot build after re-pin (state is new, registry is old; restart will reconcile): %w", err)
	}
	if err := m.Registry.RegisterProvider(provider); err != nil {
		return UpdateResult{}, fmt.Errorf("manifestslot: update: register after re-pin (state is new, registry is old; restart will reconcile): %w", err)
	}
	return UpdateResult{InstallResult: InstallResult{Manifest: manifest, Path: path, Hash: newHash}}, nil
}

// Uninstall is the deterministic teardown (AC3): revoke -> dispose the
// plugin scope (cancel running instances, join tracked goroutines, run
// deferred cleanup) -> delete state record + installed files -> honour the
// ask-about-data answer (default KEEP).
func (m *Manager) Uninstall(ctx context.Context, id string, opt UninstallOptions) (UninstallReport, error) {
	if err := m.required(); err != nil {
		return UninstallReport{}, err
	}
	if err := safePluginID(id); err != nil {
		return UninstallReport{}, err
	}
	m.mu.Lock()
	mp := m.plugins[id]
	if mp != nil {
		mp.removed = true
		delete(m.plugins, id)
	}
	m.mu.Unlock()

	rep := UninstallReport{}
	if mp != nil && mp.scope != nil {
		rep.Disposal = mp.scope.Dispose()
		if rep.Disposal.Incomplete {
			return rep, fmt.Errorf("manifestslot: uninstall %s: disposal incomplete (running instances did not wind down cleanly): %+v", id, rep.Disposal.Failed)
		}
	}

	if err := m.installer().Uninstall(ctx, id); err != nil {
		return rep, err
	}

	if m.DataRoot != "" {
		dataDir := filepath.Join(m.DataRoot, id)
		rep.DataDir = dataDir
		if opt.DeletePluginData {
			if err := os.RemoveAll(dataDir); err != nil {
				return rep, fmt.Errorf("manifestslot: uninstall %s: delete plugin data %s: %w", id, dataDir, err)
			}
			rep.DataDeleted = true
		}
		// Default answer: data kept, DataDir reported so the UI can say
		// where it is and how much it holds (SPEC-02 §6 uninstall note).
	}
	return rep, nil
}
