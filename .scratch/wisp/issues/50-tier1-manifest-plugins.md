# 50 — Tier1 manifest plugins: C14 schema, C16 host_api, install/uninstall lifecycle

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 20-host-bridge-fs-tools
**Parallel slots:** ≤2 sub-agents (A: manifest load/validate/registry; B: install/uninstall/
update/tamper lifecycle + net allowlist)
**Spec refs:** SPEC-07 §4, §6, C14, C16, D19, §14.7, S7

## What to build
Declarative Tier1 plugins: TOML manifests defining tools that map parameters → host capabilities;
the C14 schema validator, C16 host_api semver gate, install/uninstall/update lifecycle with
tamper detection, and plugin net allowlists.

## Key constraints
- Manifest (SPEC-07 §4): schema_version, id, host_api ("^1.0"), name, version, tools[] with
  name/description/risk(input-only)/capabilities/parameters(JSON Schema)/net_allowlist;
  Tier2 adds entry (51); command-kind adds exe/exe_hash/argv_template/extract/timeout_ms/
  env_allowlist (52).
- C16: host_api major mismatch → refuse load with explicit error (never silent skip); minor
  within range → load.
- Install sources: local path / git URL (registry post-S8); flow: download/copy → schema
  validate → host_api check → show FULL capability list + source → explicit user authorize →
  write plugins\<id>\ + plugin_state (hash, exe_hash slot) → register.
- Tamper: installed manifest hash mismatch at load → refuse + alert (prevents in-place swaps).
- Capabilities: declared set enforced at bridge (undeclared → reject); net_allowlist: out-of-
  list domain → auto-L2 (R5); plugin tool taint/results flow like builtins.
- Uninstall: deregister + delete + ask-about-data (default keep) + DisposalScope teardown of
  running instances; update = re-run install authorization; major-incompatible update refused.
- Plugin tools appear in `list_tools` + discovery (38) with risk pills; no plugin can register
  host_api-undeclared bridge methods.

## Out of scope
- goja runtime (51); command tools (52); SDK docs/registry (57); MCP (RESERVED).

## Acceptance criteria
- [ ] Schema validator: per-field positive/negative matrix; unknown keys rejected.
- [ ] host_api: ^1.0 plugin vs host 1.x loads; vs 2.x host refuses with clear error.
- [ ] Install flow e2e incl. authorization UI showing capabilities + source; uninstall teardown
      leaves zero goroutines/files (DisposalScope test).
- [ ] Tamper: bit-flip installed manifest → load refused + alert.
- [ ] net_allowlist: in-list L0 pass; out-of-list → L2 (R5) — e2e with mock plugin.
- [ ] Plugin tool visible in list_tools/discovery; execution via bridge records plugin id.

## Progress log (append-only, newest last)
