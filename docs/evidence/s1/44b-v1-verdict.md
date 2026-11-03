# 票 44b 验证腿 44b-v1 —— 非实现者对抗验收裁决表

- **裁决者**: 44b-v1（≠ 实现者 44b-r1）
- **被裁对象**: commit `3cc7b1d6c`（票 44b per-model 消费永久落库，D35 加列，机主批 A1445）
- **票面**: `.scratch/wisp/issues/44b-cost-per-model-persist.md`
- **实现件**: `docs/evidence/s1/44br1-impl.md`
- **验收时刻**: 2026-11-03（HEAD `05ec59ef9`；三包工作树与 44b-r1 内容一致，其后 `05ec59ef9` 只动依赖）
- **纪律**: 本腿 ⛔改码/票面/台账/push；只 commit 本证据件；派单先存档 `…/20261103/txt/44bv1-dispatch.txt`

---

## 一、真跑测试（FAIL=0）

命令：
```bash
go test ./internal/memory ./internal/observe ./internal/agent -count=1 -timeout 14m
```

| 包 | 结果 | 用时 |
|---|---|---|
| `internal/memory` | `ok` | 30.319s |
| `internal/observe` | `ok` | 5.072s |
| `internal/agent` | `ok` | 7.745s |

**FAIL = 0**。变异还原后复跑：memory 29.499s / observe 5.080s / agent 7.376s，仍全 `ok`；`go build ./...` 静默通过。

---

## 二、逐 AC 验收（票面 5 格）

| # | 票面 AC | 判据落点 | 读数 | 裁决 |
|---|---|---|---|---|
| AC#1 | `cost_daily_model` 随 04 迁移建好（幂等、版本 +1），旧 `cost_daily` 结构与数据不变 | `schema.go` `SchemaVersionTarget=5`；链 `{from:4,to:5,apply:applyV5}`；`ddlV5` 为**新表**，`cost_daily` 主键零改动 | `TestMigration44bIdempotentReopens` PASS；`TestMigration44bAddsModelTableKeepsOldData` PASS | **PASS** |
| AC#2 | Usage 落库 upsert `cost_daily_model`（provider/model 维度），与 `cost_daily` **同批**；整数 micro；未知价 tokens-only+marker | `costmodeldao_44b.go` `BumpCostDayModel`：两条 INSERT 同在**一个 `s.write` 事务**；负数/空 provider/坏 day 先拒后写；micros==0=未知价 tokens-only，marker 归读侧（`observe.PriceUnknownMarker`） | `TestBumpCostDayModelSameBatchBothRows` PASS；`TestCostModelDaysBetweenRejectsBadRanges` PASS；`TestPostModelTaskRefusesBadShapes` PASS | **PASS** |
| AC#3 | 重启后从 `cost_daily_model` 重建 per-model 基线，配额闸读数与重启前一致（fixture） | `agent/permodel_44b.go` `SeedBaseline`→`observe.QuotaLedger.SeedFrom`（月行折进 day/month/plan 桶）；`SetModelStore` 注入口 | `TestSeedBaselineSurvivesRestart` PASS | **PASS** |
| AC#4 | 日/月 rollover 在 model 维度正确；reconciliation vs `task_log` 精确 | `observe.ReconcileModels`（Σ model 行 == cost_daily == task_log）；`MonthFirst`/`DayKey`；ledger `roll()` | `TestMonthFirst` PASS；`TestQuotaLedgerSeedRolloverIsModelCorrect` PASS；`TestQuotaLedgerSeedFromBucketsDayMonthPlan` PASS；`TestPostModelTaskBooksOneBatchUnderTheMeterDay` PASS | **PASS** |
| AC#5 | 迁移测试：旧 `cost_daily` 数据不丢；model 维度从今起累积/具名回填 | `TestMigration44bAddsModelTableKeepsOldData` 逐行核 v4→v5 后 byte-identical；回填策略具名=「从今起累积」（`costmodel_44b.go:101` 注明 v5 前的日=空切片诚实形状） | PASS | **PASS** |

**5 格全 PASS。**

---

## 三、反向判据（承重证明，非表演绿）

摘掉关键件 ⇒ 点名测试必须变红；还原 ⇒ blob 与 HEAD 逐字节相同。

| 变异 | 操作 | 点名结果 | 还原核 |
|---|---|---|---|
| ① 摘 `SeedFrom` | 临时注释 `permodel_44b.go` 内 `g.ledger.SeedFrom(...)` 调用 | `--- FAIL: TestSeedBaselineSurvivesRestart` / `FAIL github.com/.../internal/agent`（exit 1） | `git checkout` 后 `hash-object=f125cfb2…` == HEAD blob ✅ |
| ② 摘 `cost_daily_model` upsert | 删 `BumpCostDayModel` 事务内第二条 INSERT（留全局写=半批） | `--- FAIL: TestBumpCostDayModelSameBatchBothRows`（exit 1）；`TestCostModelDaysBetweenRejectsBadRanges` 仍 PASS（不受本变异影响，符合预期） | `git checkout` 后 `hash-object=7b6e4228…` == HEAD blob ✅ |

两把反判据均**点名红**，证明 AC#2/AC#3 是真承重断言、非空跑。还原后三包复跑全绿。

---

## 四、纪律核查

- `git show 3cc7b1d6c --stat` = **15 文件 +1280/-56**，与题面逐字相符 ✅
- 写面（`--name-only`）= 全部落在 `internal/{memory,observe,agent}` ＋ `docs/evidence/s1/44br1-impl.md` ✅
- ⛔ 未碰 `internal/tools` / `internal/plugin` / `cmd/wisp` / `thresholds` ✅
  （`3cc7b1d6c→HEAD` 里 `internal/tools/document_text.go` 的 90/23 改动属他腿 `05ec59ef9` 依赖票，非本 commit 所为）
- 迁移走 04 schema 版本机制、`ddlV1` 文本与 `cost_daily` 主键零改动（分表不回改主表 PK，主流 harness 做法，与票面 Key constraints 一致）✅
- 本腿只 commit 证据件，⛔ push ✅

---

## 五、非缺陷留档（编排者后续，不阻断结票）

`44br1-impl.md` 已具名两处 **composition-root 咬合**：宿主在任务结束调 `BookModel`（替代 `PostTask`+全局对的成对写，否则 `cost_daily` 双计）、启动调一次 `SeedBaseline`。该落点在 `cmd/wisp`——**44b 票写面之外**（票 Key constraints 与 out-of-scope 均划在 `internal/*`，面板/宿主咬合属编排者一步）。库侧接缝与测试已完备；此为交接注记，非实现缺陷。

---

## 六、结论

- 真跑：三包 `ok`，**FAIL=0**。
- AC#1–AC#5：**逐格 PASS**（读数见第二节）。
- 反向判据①②：均**点名红**，还原 **blob same** ✅。
- 纪律：stat 与题面符、写面仅三包＋证据件、未越界 ✅。
- 放水自查：未放宽任何断言、未改 SLO/golden/thresholds、临时变异已还原并哈希核对。

# **同意结 44b** —— 附注：`cmd/wisp` 宿主咬合（`BookModel` 于任务末、`SeedBaseline` 于启动一次）留编排者落，属票面写面之外、已在实现件具名。

*裁决出自非实现者 44b-v1。*
