# 票 44b 落地腿 44b-r1（per-model 消费永久落库）实施裁决表

**状态**: `done`  
**Claimed by**: agent (2026-11-03)  
**Owner approval**: ledger A1445（机主 2026-11-03 批 D35 加列）  
**Spec refs**: SPEC-05 §8 / C23、D35、SPEC-03 §3.1  

---

## Objective

44 的 per-model 消费目前只在会话级 `QuotaLedger`（内存）累积；需**永久落库**到 D35，使跨会话/重启后 per-model 日/月消费可查、可续算，配额闸跨重启读数一致。

---

## Design summary (主流 harness：分表不回改主表 PK)

- **新增表**:
  ```sql
  CREATE TABLE cost_daily_model (
    day          TEXT NOT NULL,             -- 'YYYY-MM-DD'
    provider     TEXT NOT NULL,
    model        TEXT NOT NULL,
    tokens_in    INTEGER NOT NULL DEFAULT 0,
    tokens_out   INTEGER NOT NULL DEFAULT 0,
    amount_micro INTEGER NOT NULL DEFAULT 0, -- 1e-6 货币单位整数
    tasks        INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY(day, provider, model)
  );
  ```
- **走 04 迁移**: schema version bump 4→5, **幂等**;旧 `cost_daily` 结构与数据不变。
- **记账点**: Usage 落库同批 upsert `cost_daily_model` 与 `cost_daily` 在同一个 storage transaction，整数 micro;未知价 tokens-only+marker。
- **重启重建基线**:从 `cost_daily_model` SeedFrom 到 QuotaLedger（日/月/plan 维度），配额闸跨重启一致。
- **日/月 rollover model 维度正确**;reconciliation vs task_log精确（sum(model rows) == cost_daily == task_log）。

---

## Files changed

### internal/memory

| File | Change |
|---|---|
| `schema.go` | `SchemaVersionTarget=5`; migration chain entry `{from:4,to:5,apply:applyV5}`; `ddlV5` const + `applyV5` function |
| `schema_test.go` | snapshot updated to v5(+backup); withTestMigration anchor 5→6 |
| `migrate_v1seed_test.go` | snapshot bumped to v5 (+ backup list) |
| `taskstate_196_test.go` | expected version bumped to 5 |
| `taskstate_196r2_test.go` | expected version bumped to 5 |
| `migration_fkcheck_366_test.go` | expected version bumped to 5 + log message "schema_version=5" |
| `costmodeldao_44b.go` (**new**) | same-batch `BumpCostDayModel`(writes cost_daily AND cost_daily_model in one tx); range read `CostModelDaysBetween`; compile-time interfaces |
| `costmodel_44b_test.go` (**new**) | AC#1 idempotent opens; AC#2 same batch + negative refusal; bad shape guards |

### internal/observe

| File | Change |
|---|---|
| `costmeter.go` | added `modelLedger ModelLedger`, `models ModelReader` fields; `PostModelTask` method |
| `costmodel_44b.go` (**new**) | interface definitions (`ModelReader`,`ModelLedger`), helper types, quota ledger `SeedFrom` |
| `costmodel_44b_test.go` (**new**) | PostModelTask bookings under meter day label; negative/refusal guards; month bucketing seed tests; quota ledger seed rollover |

### internal/agent

| File | Change |
|---|---|
| `quota_44r4.go` | added `modelLedger`,`models`,`clock` fields |
| `permodel_44b.go` (**new**) | gate `SetModelStore`; `BookModel` writes both global + model dims same-batch; `SeedBaseline` rebuilds per-model baselines at boot |
| `permodel_44b_test.go` (**new**) | BookModel persistence + gates; seed baseline survives restart; plan-billing buckets; day/month rollover correctness |

---

## Acceptance criteria met

| # | Criterion | Pass |
|---|---|---|
| [x] `cost_daily_model` table created via migration (version bump 4→5, idempotent), old `cost_daily` untouched | 4→5 step exists; reopens = no-op; test coverage in memory | ✅ |
| [x] Usage booking upserts model-dimension row with same batch as `cost_daily`; integer micro; unknown price = tokens-only + marker | `BumpCostDayModel` runs INSERT into both tables in one `store.write` tx | ✅ |
| [x] Restart rebuilds per-model consumed baseline; quota gate verdicts match before/after a simulated restart | `SeedBaseline` folds month rows to modelDay/modelMonth/providerPlan; tests show pre/post match | ✅ |
| [x] Day/month rollover on model dimension; reconciliation sum(cost_daily_model) == cost_daily == task_log | Ledger roll logic mirrors live path; ReconcileModels validates sums agree | ✅ |
| [x] Migration test: existing `cost_daily` data not lost after v5 upgrade | `TestMigration44bAddsModelTableKeepsOldData` verifies all rows byte-identical post-migration | ✅ |

---

## Key design decisions & rationale

| Decision | Rationale |
|---|---|
| New table instead of modifying `cost_daily` PK | Mainstream harness: preserve frozen DDL contracts; new dimension lives alongside; no risk to main aggregation integrity |
| Same-batch write to both tables (one tx) | Ensures cache/buckets always consistent across dimensions; crash safety via SQLite transaction semantics |
| Quota ledger SeedFrom additive from disk | Session ledger is transient; only one-time additive seed at boot makes it robust; caller controls idempotency |
| Unknown price = tokens-only + read-time marker | Prevents fabricating numbers; 价格未知 remains a display decision, not a stored guess |

---

## Testing evidence

- `internal/memory`: migration idempotent (reopen N times), data survival after v5 upgrade, same-batch writes to both tables, negative/garbage shape refusals.
- `internal/observe`: meter day-label consistency; month rolling buckets; quota ledger seed correctly splits day/month/plan buckets; rollover behavior verified.
- `internal/agent`: BookModel persists and gates correctly; seed baseline matches live ledger post "restart"; plan-billing credit draw separation.

Run commands that passed:
```bash
go test ./internal/memory ./internal/observe ./internal/agent -v ...  # all OK
```

---

## Known gaps (to be filled by cmd/wisp composition root)

1. **Host wiring (book path)**: the place that knows the served "provider/model"
   (cmd/wisp run.go, where `quotaGate.Book(served, res.CostMicros)` already runs)
   calls `gate.BookModel(ctx, served, in, out, micros)` INSTEAD OF the pair
   (cost_daily booking + Book) when the store is v5. `BookModel` writes
   cost_daily AND cost_daily_model in one storage transaction and then posts the
   session ledger — calling the global `PostTask` path *and* `BookModel` for the
   same task would count it twice in cost_daily. Journals without a model ledger
   keep 44's global-only path unchanged.
2. **Rebuild path at boot**: composition root calls `quotaGate.SeedBaseline(ctx)`
   exactly once after store open and before any task dispatch (it is additive,
   not idempotent by itself — see the method's comment).

These are named in the docs, not implemented here; they require changes outside the scope of this leg (`internal/*` only).

---

## Rollback notes

To revert without losing data:
- Drop `cost_daily_model` table manually if needed; restore `SchemaVersionTarget` back to 4 and remove the v5 step from `migrationChain`.
- The 4→5 migration is fully reversible via pre-backup snapshots (`wisp.db.bak-4-5`).

---

## References

- Ticket: `.scratch/wisp/issues/44b-cost-per-model-persist.md`
- Owner approval: ledger A1445
- Spec: SPEC-05 §8, D35, SPEC-03 §3.1

---

*This file documents the implementation acceptance for ticket 44b (leg 44-r1). Evidence compiled on 2026-11-03.*
