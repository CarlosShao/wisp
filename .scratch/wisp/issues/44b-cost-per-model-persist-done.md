# 44b — CostMeter 补：per-model 消费永久落库（D35 加列）

**Status:** **已结案 done**（2026-11-03 11:5x 编排者：44b-v1 非实现者验收 5 格全 PASS〔memory 30.3s/observe 5.1s/agent 7.7s FAIL=0，commit `3cc7b1d6c`〕）
**Claimed by:** —
**Last update:** 2026-11-03
**Blocked by:** —（44 已结，quota 组件层已交）
**Spec refs:** SPEC-05 §8 / C23、D35、SPEC-03 §3.1
**Owner approval:** 机主 2026-11-03 批「D35 加列」（台账 A1445）。

## What to build
44 的 per-model 消费目前只在会话级 `QuotaLedger`（内存）累积；需**永久落库**到 D35，使跨会话 / 重启后 per-model 日 / 月消费可查、可续算，配额闸跨重启读数一致。

## Key constraints
- **数据模型（编排者定；主流 harness：分表优于改主表 PK）**：新增表
  `cost_daily_model(day TEXT, provider TEXT, model TEXT, tokens_in INTEGER, tokens_out INTEGER, amount_micro INTEGER, tasks INTEGER, PRIMARY KEY(day, provider, model))`；
  **不回改 `cost_daily` 主表 PK**。迁移走 04 的 schema 版本机制（版本 +1，**幂等**）。
- 记账点：C6 Usage 落库时同步 upsert `cost_daily_model`（provider/model 维度），与 `cost_daily` **同批**；整数 micro 单位；未知价 tokens-only + "价格未知" marker（沿用 44）；reconciliation 与 `task_log` 精确。
- 启动 / 重启后从 `cost_daily_model` **重建 per-model 已消费基线**进 `QuotaLedger`，使配额闸跨会话正确。
- 日 / 月 rollover 在 model 维度正确。

## Out of scope
- 面板成本页展示（40）；价格研究（starter 表够）。

## Acceptance criteria
- [x] `cost_daily_model` 表随 04 迁移建好（幂等、schema 版本 +1），旧 `cost_daily` 结构与数据不变。
- [x] Usage 落库 upsert `cost_daily_model`（provider/model 维度），与 `cost_daily` 同批；整数 micro；未知价 tokens-only+marker。
- [x] 重启后从 `cost_daily_model` 重建 per-model 已消费基线，配额闸读数与重启前一致（fixture）。
- [x] 日 / 月 rollover 在 model 维度正确；reconciliation vs `task_log` 精确。
- [x] 迁移测试：已有 `cost_daily` 数据不丢；model 维度从今起累积或具名回填策略。

## Progress log (append-only, newest last)

- [2026-11-03 11:5x +08] agent=编排者（结案）did=44b-v1 5格全PASS(memory/observe/agent FAIL=0);44b-r1 3cc7b1d6c D35 cost_daily_model分表+同批upsert+SeedBaseline重建基线.
