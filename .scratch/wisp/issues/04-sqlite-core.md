# 04 — SQLite storage core: D35 schema v1, WAL, single writer, retention

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: schema+DAO+writer goroutine; B: RetentionJob + crash/privacy tests)
**Spec refs:** SPEC-02 §2–§8, D35, §14.11, D20, C13, C23

## What to build
The `memory` module's storage core: all eight D35 tables with exact DDL, WAL configuration,
serialized single-writer goroutine (`db-writer`), schema versioning/migration, RetentionJob
(30/400-day cleanup + artifacts 500MB LRU), and the privacy page data operations (purge/export/
delete-one) exposed as module APIs (GUI pages come in ticket 40).

## Key constraints
- Driver: `modernc.org/sqlite` (pure Go — cgo budget stays with sherpa; SPEC-02 §2).
- DDL exactly per SPEC-02 §3 (schema_meta, profile, memory, task_log, tool_call, approval_grant,
  cost_daily, plugin_state incl. `exe_hash`). No added/removed columns.
- Pragmas: `journal_mode=WAL`, `synchronous=NORMAL`, `busy_timeout=3000`,
  `wal_autocheckpoint=1000`, `wal_checkpoint(TRUNCATE)` at startup. All writes serialized through
  the `db-writer` goroutine (no ad-hoc SQLITE_BUSY retry anywhere).
- Connections opened on demand, not resident (D20); verify WAL `-wal/-shm` interplay across
  clean exit / kill / power-loss reopen.
- `schema_meta.schema_version` migration chain; each step single transaction + pre-backup
  `wisp.db.bak-<from>-<to>`; unmigratable → fail startup into Unconfigured (never overwrite).
- RetentionJob: run once 5min after boot + every 24h (`retention-job` goroutine, DisposalScope
  governed): task_log/tool_call >30d, cost_daily >400d, expired grants kept 30d for audit,
  `artifacts\` LRU at 500MB.
- Time: persisted timestamps = wall clock; retention math = monotonic (D42#9).
- L1 profile ≤20 rows enforced app-layer with LRU-by-updated_at eviction that WRITES A LOG line
  (eviction must be visible, D35 rule).

## Out of scope
- Panel privacy UI (40); grant lifecycle logic (49); cost aggregation logic (44) — tables only.

## Acceptance criteria
- [ ] Contract test per table: columns/indexes via sqlite_master introspection match SPEC-02 §3.
- [ ] Concurrency: 2 writers + 4 readers × 10s → zero busy errors, WAL bounded by autocheckpoint.
- [ ] Crash recovery: kill mid-write → reopen passes `integrity_check`, checkpoint truncates WAL.
- [ ] Retention boundary tests (29/30/31 days; 399/400/401 days).
- [ ] Privacy ops API tests: purge-all / export-JSON / delete-one for profile, memory, task_log,
      tool_call, artifacts.
- [ ] Profile LRU eviction logs (test asserts log line).

## Progress log (append-only, newest last)
