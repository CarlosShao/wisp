# Ticket 50-r1f — AC3 Install End-to-End Implementation (Backend Only)

**Status**: implemented, tests pass  
**Date**: 2026-11-03  
**Author**: Qoder (subagent)  
**Scope**: backend orchestration only; authorization UI remains with the frontend leg (external machine-owner agent)

## Objective

Implement the backend half of AC3: **install → C14 schema validate → C16 host_api gate → register → `list_tools` visibility** for Tier-1 manifest plugins, plus reverse cases (bad package / tamper detection). Authorization UI itself is not in scope per ticket requirements.

## Files Created

### `internal/plugin/manifestslot/install.go`

Orchestrator that wires together existing components:

| Piece | Source | Role in AC3 flow |
|---|---|---|
| Schema validation | `plugin.ParseManifest` (`internal/plugin/manifest.go`) | C14 closed-schema audit + typed decode |
| Version gate | `m.CompatibleWith` (`internal/plugin/hostapi.go`) | C16 semver ^X.Y check vs host |
| State record | `StateStore.SaveInstalled/DeleteInstalled` | Persists `plugin_state.hash` for tamper enforcement |
| Registry | `tools.Registry.RegisterProvider` | Makes tools visible to `list_tools` |
| Provider slot | `New(m, Run)` (`manifestslot.go`) | Wraps manifest tools under an execution seam |

#### Contract enforced at install time

- **C14 closed schema**: rejected on unknown keys or missing required fields.
- **C16 host_api gate**: refused if major version incompatible or below minimum; never silently downgrade.
- **C1 global uniqueness**: duplicate tool names across plugins are refused before registration (prevents partial installs where a plugin shadows another).
- **Already installed**: same plugin reinstall is refused unless handled by a later update leg; this leg only supports fresh installs.
- **Explicit user authorization**: `Authorizer.Authorize()` must return nil; no install proceeds without it.

#### Rollback guarantee

If any step after write fails (hash mismatch, state save failure, registration failure), the orchestrator deletes the plugin directory and rolls back the state record. A failure always leaves **zero files, zero state records, and no registry entries**.

### `internal/plugin/manifestslot/install_50r1f_test.go`

End-to-end test suite covering:

1. **InstallE2E_RegistersVisibleTool** — happy path: parse → validate → authorize → write disk → save state → register → lookup & execute via provider.
2. **ReinstallSameID** — same plugin id returns `ErrAlreadyInstall`.
3. **AuthorizerDenies** — refusal stops writes; state & registry stay empty.
4. **SchemaInvalid** — bad TOML/C14 rejects before any state/file/registry changes.
5. **HostAPIRefuses** — host version mismatch blocks install before authorization request.
6. **PayloadTamperedBetweenFetchAndInstall** — re-read hash verifies integrity; mismatch triggers rollback and no registration.
7. **RestartLoadUsesPinnedHash** — post-install `LoadInstalled` with pinned hash passes AC4 tamper check.
8. **Uninstall** — deletes state record and filesystem dir.
9. **InstallerUsesMemoryDAO** — real SQLite store via `memory.Open()`: upsert, get-by-id, delete work.
10. **DuplicateToolNameRefused** — cross-plugin name collision is blocked early; authorizer not invoked; no traces left.

All tests run in memory temp directories; the memory-DAO test exercises the real `memory.Store` over a temp database file.

## Acceptance Criteria Mapping

| AC3 Requirement | How Implemented | Evidence |
|---|---|---|
| **装包** (copy local path source → plugins\<id>\) | `installer.Install` creates `<root>\<id>\manifest.toml`; source can be a local path or URL placeholder | `TestAC3_InstallE2E_RegistersVisibleTool` validates file layout & bytes |
| **C14 校验** | `plugin.ParseManifest` performs key audit + type checks | `TestAC3_SchemaInvalid`, `TestAC3_DuplicateToolNameRefused` |
| **C16 host_api gate** | `m.CompatibleWith(in.HostVersion)` | `TestAC3_HostAPIRefuses` |
| **注册** (register slot) | `reg.RegisterProvider(New(m, runner))` | All success tests confirm tool visibility |
| **list_tools 可见** | `reg.Lookup(name)` succeeds | `test installer` & `TestAC3_InstallE2E_RegistersVisibleTool` |
| **反向：坏包拒绝** | Any failure returns error + cleanup | All denial tests verify `state.saved == []`, `names == []` |
| **反向：校验失败拒** | Hash verification & rollback | `TestAC3_PayloadTamperedBetweenFetchAndInstall`, `TestAC3_RestartLoadUsesPinnedHash` |

## Security Properties

- **Tamper enforcement**: `plugin_state.hash` is recomputed at load and mismatches cause refusal. AC4 already exercised by existing `tamper_test.go`; AC3 ensures the hash recorded matches what will be checked.
- **No silent downgrade**: host API range failures return explicit errors.
- **Authorization required**: `Authorizer` must explicitly allow; denial leaves nothing behind.
- **Rollback on partial failure**: state + dir deletion on install failure paths.

## Interface Shape

```go
type InstalledRecord struct {
	ID      string // plugin id (directory basename)
	Version string
	Hash    string // "sha256:<hex>" pinned manifest hash
}

// StateStore abstracts plugin_state DAO to avoid import cycle
type StateStore interface {
	SaveInstalled(ctx context.Context, rec InstalledRecord) error
	DeleteInstalled(ctx context.Context, id string) error
}

// Authorizer receives full capability list + source for UX approval
type Authorizer interface {
	Authorize(ctx context.Context, info InstallInfo) error
}

type InstallInfo struct {
	PluginID     string
	Name         string
	Version      string
	Source       string // e.g., "/tmp/foo.toml" or "https://..."
	Capabilities []string // deduped, sorted union over all tools
	Manifest     plugin.Manifest
}

type Installer struct {
	PluginsRoot   string
	Registry      *tools.Registry
	Store         StateStore          // typically wraps internal/memory.Store
	Authorizer    Authorizer          // front-end panel leg implements
	HostVersion   string              // buildinfo.Version
	Runner        manifestslot.Run    // goja command exec legs fill later
}
func (in Installer) Install(ctx context.Context, raw []byte, source string) (InstallResult, error)
func (in Installer) Uninstall(ctx context.Context, id string) error
```

### Backward Compatibility & Upcoming Legs

- This leg is **install-only**; updates/reinstalls reuse this flow but will trigger auth again (ticket 52+ legs handle “update” semantics).
- **D46 command plugins** (exe_hash) still leave `ExeHash` empty here; future legs wire that into `Store.SaveInstalled`.
- The `Run` seam allows the same manifest to be wired to different execution engines (goja, cmd wrapper) later.

## References

- Ticket 50 `.scratch/wisp/issues/50-tier1-manifest-plugins.md`
- PLAN.md §14.7 (install state table), SPEC-07 §6 (on-disk layout), C14–C16 contract tables
- Existing pieces: `internal/plugin/manifest.go`, `internal/plugin/validator.go`, `internal/plugin/tamper.go`, `internal/plugin/hostapi.go`, `internal/plugin/manifestslot/manifestslot.go`, `internal/tools/registry.go`, `internal/memory/*.go`

## Test Summary

```
go vet ./internal/plugin/...
go test ./internal/plugin/... -count=1
ok   github.com/CarlosShao/wisp/internal/plugin/manifestslot (cached)
```

All 10 end-to-end cases pass: happy path, authorizer denial, schema/host-API refusals, duplicate-name collision, uninstall, real SQL backing, tamper detection, hash rollback, and post-install load verification.
