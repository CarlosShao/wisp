# 44-v1 — CostMeter (C23) 非实现者验收表（裁决腿）

- **腿**：票 44 验证腿 44-v1 ＝ 非实现者验收（裁决者 ≠ 实现者，`AGENTS.md` §0 第 3 条／`SPEC-12 §4.3` #1/#3）。
- **被验对象**：44-r1 两枚 commit —— `8827e54ca`（写腿，15 枚文件）＋ `dc88642e0`（追加更正，1 枚文件＝同一份证据件）。
- **票面**：`.scratch/wisp/issues/44-costmeter-c23.md`（58 行，6 枚 `- [ ]` 勾框）。
- **日期**：2026-11-02。**本腿零改码**：⛔ 改实现／票面／台账／push；反向判据那几发的"改→红→还原"见 §4，还原后 `git status` 对 15 枚交付件＝空。
- **临时件**：`D:\tmp\wisp\20261102\44v1\`（`logs/01…23`、`bak/` 六枚原件备份、`head/` ＝ `git archive HEAD` 的隔离副本）。只建不删（`issues/README` 规则 8）。

---

## 1. 起手三核 ＋ 两枚 commit 与题面

| 核 | 尺（可逐字重跑） | 读数 |
|---|---|---|
| HEAD 现盘最新 | `git rev-parse --short HEAD` | 起手 `96e8cb278` ≥ 要求的 `885811c79` ✅。干活期间进到 `547d9da00`（22-r1 落地）＋编排者 A1427/A1428/A1429 —— **他腿在飞，本腿零碰** |
| 索引干净 | `git diff --cached --name-only`；`git diff --name-only`；`git status --porcelain \| ? {-notmatch '^\?\?'}` | 三把全空 ✅（只有未跟踪的scratch/他腿面） |
| 定位两枚 | `git log --oneline --grep="44-r1"` | `8827e54ca`＋`dc88642e0` ✅（其余命中是编排者记账与 45-r1/31-r1 的引用） |
| 盘＝HEAD | `git diff --stat HEAD -- <15 枚>` | 空 ✅ 交付件在盘上与提交一致，期间无人改过 |

**`--stat` 对题面**：`8827e54ca` ＝ **15 files changed, 3165 insertions(+), 14 deletions(-)** ＝ 14 枚代码／测试跨 `internal/observe`(6)／`internal/memory`(2)／`internal/agent`(6) ＋ 1 枚证据件 `docs/evidence/s1/44r1-impl.md`。与派单题面"15 文件跨 3 包"**一致** ✅。`dc88642e0` ＝ 同一份证据件 `9 insertions(+), 1 deletion(-)`＝**只追加、⛔ amend** ✅。

> ⚠ 一处**别误判**：`git diff --name-status 8827e54ca^..dc88642e0` 会多出一枚 `docs/evidence/s1/279r1-stop.md`，那是夹在两枚之间的 **279-r1**（`d27b3ff43`）写的，不是 44-r1 的。尺＝逐枚 `git show --name-only`：两枚 44-r1 commit 各自的文件名单里都没有它。

---

## 2. 真跑（点名）

### 2.1 派单指定的那把尺

```
go test ./internal/observe ./internal/memory ./internal/agent -run 'Cost|cost|Pricing|pricing|Budget|budget' -count=1 -timeout 20m
```
读数：`ok observe / ok memory / ok agent`，**FAIL=0** ✅（`logs/01-named-run.txt`）。

**但这把尺盖不住本票，必须具名**（`-v` 逐条点名，`logs/02-named-verbose.txt`）：18 命中里**只有 11 枚属于 44**，另 **7 枚是他票的**（`TestGrantCostPluginStateDAO`、`TestLoopGuardTokenBudgetScalesWithWindow`、`TestGoldenBudgetExhaustionStuck`、`TestGoroutineBudgetDrains`、`TestD39SectionBudgetsAndTotalCeiling`、`TestEnforceTotalTerminatesOnABelowBudgetSection`、`TestSpillPointerNoticeKeepsHeadTailAndBudget174`），而 44 自己的 36 枚里**有 25 枚这尺根本扫不到**，其中包括最承重的 `TestAC3NewTaskRefusedBeforeTheProvider`／`TestAC3RunningTaskFinishesPastTheLimit`／`TestAC3WarnModeStillRunsTasks`／`TestAC2DayRolloverMovesTheBucket`／`TestAC2MonthRolloverResetsTheMonthRange`／`TestAC1PriceMicrosMatchesHandComputedFixture`。**⇒ 光看这条命令的 FAIL=0，预算闸整个摘掉也不会红。** 所以本腿另跑全名册（§2.2），并建议编排者以后把 44 的尺改成 `-run '^(?:AC1|AC2|AC3|AC4|AC6|…)$'` 或直接点全名。

### 2.2 全 36 枚名册（逐枚点名，`logs/03…05-*-roster.txt`）

名册来源＝6 枚交付测试文件里 `^func Test` 实抽（7+15+4+2+5+3＝**36**，与 r1 自述"observe 22／memory 4／agent 10"一致）。

| 包 | RUN | PASS | FAIL | SKIP |
|---|---|---|---|---|
| `internal/observe` | 22 | **22** | **0** | 0 |
| `internal/memory` | 4 | **4** | **0** | 0 |
| `internal/agent` | 10 | **10** | **0** | 0 |
| 合计 | 36 | **36** | **0** | 0 |

36 枚名字（全部 `--- PASS`）：`TestAC1PriceMicrosMatchesHandComputedFixture`、`TestAC1PriceMicrosIsIntegerAndSaturates`、`TestAC4UnknownCardPricesNothing`、`TestAC4TableLookupMissesSayUnknownAndCostNoMoney`、`TestPricingRejectsFabricatedAndMalformedRows`、`TestAC6EmbeddedTableIsValidAndVersioned`、`TestOverrideValidationMatchesTableValidation`、`TestAC2TodayReadsOnlyItsOwnDay`、`TestAC2DayRolloverMovesTheBucket`、`TestAC2MonthRolloverResetsTheMonthRange`、`TestAC2ReconcileAgreesAndDisagrees`、`TestAC2ReconcileRefusesWithoutBothMouths`、`TestPostTaskRefusesNegativesAndNamesTheDay`、`TestAC3BudgetBoundaries79And80And100`、`TestAC3OverBudgetWarnModeNeverBlocksAndNamesItself`、`TestAC3BindingCapIsTheTighterOneAndBothAreListed`、`TestAC3UnsetBudgetsAreUncappedNotBlocking`、`TestBudgetFromRejectsBadConfig`、`TestBasisPointsSaturateInsteadOfWrapping`、`TestDayBoundsUseTheInjectedCalendar`、`TestMeterRefusesWithoutReader`、`TestDayKeyIsTheShapeMemoryValidates`、`TestAC2CostDaysBetweenFiltersAndOrders`、`TestAC2SumTasksEndedBetweenCountsOnlyFinishedRows`、`TestAC2ReconcileAgainstRealStoreAndRepair`、`TestSetCostDayOverwritesAndRefusesNegatives`、`TestAC1GoldenPerTaskCostMatchesHandComputed`、`TestAC4UnknownPriceShowsTokensOnlyAndPricedFalse`、`TestLoopBooksFinishedTaskIntoCostDay`、`TestAC2ReconcilesAgainstTaskLog`、`TestAC3BudgetGateReadsTheBookedSpend`、`TestAC6PriceTableBumpLeavesAggregatesAlone`、`TestAC4UnknownPriceKeepsTokens`、`TestAC3RunningTaskFinishesPastTheLimit`、`TestAC3NewTaskRefusedBeforeTheProvider`、`TestAC3WarnModeStillRunsTasks`。

**不是空跑**：`internal/agent` 那 10 枚走真 `memory.Open` SQLite ＋真 C5 golden SSE 回放（`internal/agent/testdata/golden/tool-then-text.sse`，最后由 `ae23da673` 票 10 提交，本票零改）⇒ 按 `issues/README` 1.3"不许用 mock 代替真的"这条**过**。

---

## 3. 逐格 AC（票面 6 枚勾框 ＋ 派单 5 格）

### 3.1 票面六格

| # | 票面逐字 | 判 | 本腿读数（自己复跑，非转抄 r1） |
|---|---|---|---|
| **#1** | `Multi-turn golden task: per-task tokens/cost match hand-computed fixture to the unit` | **PASS** | 数学：`PriceMicros`＝`bits.Mul64`＋**一次**除法＋half-up，溢出**饱和不回转**；六子形手算逐条在测试里带 `byHand` 串（`1000*2.5e6+200*10e6=4500`、缓存走 cache 价 `3300`、`Cached=0` 走 In 价 `2500`、`cached>in` 夹到 `in`＝`50`、`1 token`＝`10`、半微 `2.5→3`）。环路：`TestAC1GoldenPerTaskCostMatchesHandComputed` 真字节回放 ⇒ `{330,24,150}` tokens 与 **750 micros** 逐单位相等，且同数落进 `task_log` 行；夹具自检（`wantMicros!=750 ⇒ Fatalf`）在场。per-turn 求和＝4500+3300=7800 也在 observe 侧钉住。**"全仓一枚算钱口"复核**：`git grep math.Round -- internal/...` 只剩两处**边界换算**（`cost.go:95` 把 D36 冻结的 float64 **率**换成整数微单位、`costbudget.go:66` 把 alert_ratio float 换成基点），钱的**金额**全程 int64；`agent.PriceMicros` 逐字委托 `observe.PriceMicros`。✅ |
| **#2** | `Aggregates: day rollover + month rollover correct (fixture clocks); reconciliation with task_log exact` | **PASS** | 日：夹具时钟 23:59→00:00，同一 store 只换钟 ⇒ `11-02.tasks=1 / 11-03.tasks=1` 分家、新日 `micros=11` **不吸收**昨日 500。月：11-30→12-01 ⇒ 范围 `2026-11-01..2026-11-30`→`2026-12-01..2026-12-01`，旧日行逐字不动。对账：`Reconcile` 先 `Agree=true`，植一枚 base-only 游离行 ⇒ `Agree=false` 且两侧数字对得上（读数里 `micros=632 tasks=3` 那次 refresh 可见），`RefreshDay` 从 `task_log` 重算后再 `Agree=true`。真 SQLite 侧边界也钉住了：`ended_at` 首秒**含**、末秒**含**、次日首秒**不含**、`ended_at IS NULL`（在跑）**不含**、空区间读取＝零行不报错、`from>=to` 空区间＝报错拒绝。入账口＝`finish()` 在 `finishTask` 成功分支里 `PostTask`（失败不入账＝两嘴不分裂），`cost_daily` 从此有生产写手（GAP-19 那格）。✅ |
| **#3** | `Budget: 80% warn event; 100% → new task rejected with explicit message + running task finishes; warn-only mode alternative verified` | **部分 PASS（⚠ 见右）** | 阈值：整数基点，`789/799→ok`、`800→warn`、`999→warn`、`1000/1500→over`，且 warn/over 都必须带非空 `Notice`、ok 必须空 `Notice`（反向也钉）。环路三枚真跑：① `TestAC3NewTaskRefusedBeforeTheProvider`＝`startTask` 之后、provider 之前拒，`Status=failed`、`Err.Class=budget`、文案含 `暂停新任务`、**`second.requests()==0`**（一字节没出机器）、`task_log.error_class='budget'`＋`summary_text` 带文案、对账仍 `Agree` 且 `tasks=2`、钱没涨；② `TestAC3RunningTaskFinishesPastTheLimit`＝`flipSink` 在**首个 text delta**把日预算压到 1 micro（飞行中翻越），该任务仍 `completed`，且先 `if fs.flipped()==0 { t.Fatal }` 钉住"这一发真测到了"；③ warn 模超 100% 仍放行、provider 真被够到、`Notice` 自报模式。**⚠ 未满足的那半＝"80% warn event"没有事件出口**：见 §3.2 第 2 项。 |
| **#4** | `Unknown-price provider → tokens-only + marker (no invented currency amounts)` | **PASS** | 全零卡 `Known()=false` ⇒ `PriceMicros` 恒 0（**不是**"免费"）；真免费必须行上明写 `billing:"free"`；marker 常量逐字 `价格未知`（票面 `:19-20`），`TestAC4UnknownCardPricesNothing`／环路侧 `TestAC4UnknownPriceKeepsTokens` 钉住字符串与"tokens 还在、钱＝0、日行仍入账 `tasks=1/micros=0`"。表校验拒九种伪造／畸形形（无 source 的价、负价、无 currency、key 不是 `provider/model`、版本<1、缺 models 对象…）。`WithOverrides` 与表内行同一套校验且**不改接收者**。✅ **但注意**：`pricing.json` 现收 **0 枚价**（`"models": {}`），这是票面 `:20` `never fabricate amounts`＋`:42` 把价研究划 out-of-scope 的**正确**结果，机制齐、行由机主按公开价卡填。 |
| **#5** | `Quota: per-model daily/monthly gates; plan-credit bookkeeping…; pre-dispatch anti-negative test; plan vs pay-per-token tracked separately per model` | **FAIL（未交付；停手上报本身成立）** | 本腿**独立复跑**了 r1 的那堵墙，实测为真：`internal/memory/schema.go` 里 `task_log` ＝ **11 列**（`id/started_at/ended_at/state/query_text/summary_text/cost_tokens_in/cost_tokens_out/cost_amount_micro/currency/error_class`），**无 provider／无 model**；`cost_daily` ＝ `day/tokens_in/tokens_out/cost_micros/tasks`；`provider_health.quota_state` 只是文本态 ⇒ per-model 日/月花费**盘上无处还原**，这一点具名**准确**。补齐要给 D35 冻结表加列/加表＝改契约文字（`SPEC-12 §4.1`／`AGENTS.md` §0 第 2 条）＝**人工批准面**，⛔ 写腿自改 ⇒ **停手上报是正确处置，不是偷懒**。链推进那半票面 `:43` 逐字划给票 11。**但 AC#5 这一格作为一个勾框并未满足**：per-model 闸、`plan_credit_total_micro` 分账、"remaining 0 → 跳过派发＋链推进＋可见事件"这三样今天**都不在盘上**。✅ 的是"没违令"，❌ 的是"没交付"。 |
| **#6** | `Pricing-table version bump → recorded, old aggregates unchanged` | **PASS（记录载体具名降级）** | `pricing.json` 带 `version`（现 1）、`ParsePriceTable` 拒 `version<1`；`RecordPriceTableVersion()` 落一条 slog（`version/updated/models/recomputes_aggregates=false`）；环路真 store 用例：入账后取 `before`，表 bump＋1 并 `Rebuild` ⇒ `PriceTableVersion()` 跟上新值而 `after == before`（struct 相等），再 `RefreshDay` 也只从 `task_log` 重算 ⇒ `healed.CostMicros` 仍是入账那一笔，**没有任何路径按新价重算历史**。⚠ "recorded"的载体只有 slog JSONL＋内存里的 `PriceTable.Version`，**没有逐日持久钉版本**（`cost_daily`/`task_log` 没有可放版本号的列，与 AC#5 同一堵墙，r1 §4 已具名）。按票面字面（版本记录在案＋旧聚合不变）判 PASS，欠账见 §3.2 第 3 项。 |

### 3.2 派单五格

| 格 | 判 | 一句话读数 |
|---|---|---|
| ① per-task 成本累加正确（定价×tokens） | **PASS** | §3.1 #1 |
| ② 日聚合（按日 bucket 汇总） | **PASS** | §3.1 #2；`DayKey` 由注入历法定桶、落盘时间戳仍 wall UTC（D42#9） |
| ③ 预算 80-100% 预警阈值行为 | **PASS** | 阈值本身（800 基线、整数基点、800=warn/1000=over、pause/warn 两模、拦在 provider 之前）全部量到且突变能红；见 §3.1 #3 |
| ④ 零回归 | **PASS** | §5 |
| ⑤ 纪律 | **PASS（有 1 项具名不实已自纠＋1 项编排者侧缺口）** | §6 |

### 3.3 三处具名欠账（⛔ 不算"通过"，只是写清楚落点）

1. **AC#5 整格未交付**（per-model 日/月闸＋plan-credit 分账＋防负夹具）。**前提卡在 D35 冻结表要不要加列/加表＝机主批准面**，不由写腿裁。⇒ 结不结 44 的钥匙在编排者／机主手里：要么裁"AC#5 拆出去另开一票"，要么批加列。
2. **`[cost]` 80% 那一发今天没有事件出口**（票面 `:21` 写的是 `≥80% → warn toast + badge`，AC#3 写的是 `80% warn event`）。尺＝`git grep -n "publish(Event" -- internal/agent/loop.go internal/agent/costbook_44.go` ⇒ 成本路径**只有 `loop.go:415` 那一枚 `EvError`，且只在 100% 被拒时发**；warn 态 `AllowNewTask=true` ⇒ `preDispatchBudget` 直接返回 nil ⇒ **总线上什么都没有**。`CostBudget()` 返回的 `BudgetVerdict{Level:"warn", Notice:"费用已用到预算的 80%…"}` 是真的、也测到了（`TestAC3BudgetBoundaries…` 钉非空），但**没有任何生产者把它发出去**，渲染归票 40 属另一层。⇒ **80% 预警这一寸＝机制在、出口未接**，本格不给"整格 PASS"。
3. **装配未接，产品今天不咬**：`git grep -rn "CostBudget|SetCostBudget|BudgetFrom" -- cmd internal`（排除测试与定义自身）⇒ **零命中**。`cmd/wisp/run.go:1342-1349` 仍不填 `Config.Price`／`[cost]`（台账 11218 行 261-p1 段早具名同一件事）⇒ 真跑 `wisp run` 时价卡恒零（走 marker）、预算恒 uncapped。**机制齐全、口子已交（`BudgetFrom`＋`Loop.SetCostBudget`），差 cmd/wisp 那一寸装配**——那是派单点名禁面，本腿与写腿都⛔写，不算违令，但必须记账。
4. **测试件一处过期注释（小）**：`internal/agent/cost_44_test.go:14-19` 那段头注释写 `Usage.Add max-merges → totals: in=180, out=10` 与 `PriceMicros = 180*2 + 10*10 = 460 micros`，而同函数体逐字断言 `{330,24,150}` 与 **750**（累加、不是 max-merge）。**可执行断言是对的**（本腿量过），错的是注释——一句讲错钱语义的注释该改，归 44-r1 顺手一发。

---

## 4. 反向判据（本腿自己做的五发，非转抄 r1 的 M1-M8）＋还原核

每发＝备份原件→改一处→跑点名用例→**逐字节还原**→再跑。备份在 `44v1\bak\`，红句在 `44v1\logs\11…17`。

| 发 | 动哪一处（逐字） | 红在谁身上（具名） | 还原核 |
|---|---|---|---|
| **V1 摘成本累加** | `costbook_44.go:62` `if l.costMeter == nil {` → `if true \|\| l.costMeter == nil {` | **4 枚红**：`TestLoopBooksFinishedTaskIntoCostDay`（`cost_daily[2026-11-02]={…Tasks:0}, want exactly one booked task`）、`TestAC2ReconcilesAgainstTaskLog`（`agg={…Tasks:0} base={660/48/1500/Tasks:2}`）、`TestAC3BudgetGateReadsTheBookedSpend`（100% 仍放行、`level "ok" want over`、`gate read day spend 0, want 750`）、`TestAC4UnknownPriceKeepsTokens`（`want 1 task / 0 micros`） | `sha256=E71CB0CE0C8699CD` ＝基线；`git status` 空；4 枚复跑 `ok` ✅ |
| **V2a 改坏 100% 阈值** | `costbudget.go:158` `case spent >= cap:` → `case spent > cap:` | `TestAC3BudgetBoundaries79And80And100`（`spent 1000: level "warn", want "over"` ＋ `AllowNewTask true, want false`）；agent 侧连带 `TestAC3NewTaskRefusedBeforeTheProvider`（`status = completed (), want failed`）＋ `TestAC3BudgetGateReadsTheBookedSpend` | `A847C48FC2DD8360` ＝基线；`git status` 空 ✅ |
| **V2b 改坏 80% 预警线** | `costbudget.go:160` `case c.Basis >= p.AlertBasisPoints:` → `case c.Basis > …` | `TestAC3BudgetBoundaries79And80And100`（`spent 800: level "ok", want "warn"` ＋ `warn verdict with no user-visible notice`）——**80% 那一格单独有牙的尺** | `A847C48FC2DD8360` ＝基线 ✅ |
| **V3 改坏定价舍入** | `pricing.go:120` `q := (num + million/2) / million` → `q := num / million` | `TestAC1PriceMicrosMatchesHandComputedFixture`（`half-micro total rounds half-up: PriceMicros = 2, want 3`）⇒ **AC#1 那格"定价×tokens"确实是手算夹具钉住的** | `27BB153ECE1E8F22` ＝基线 ✅ |
| **V4 摘掉派发前闸** | `costbook_44.go:141` `if !rep.Metered \|\| rep.Verdict.AllowNewTask {` → 前面加 `true \|\|` | `TestAC3NewTaskRefusedBeforeTheProvider`（`costgate_44_test.go:111: status = completed (), want failed: 100% of a pause-mode cap must reject a new task`）；同时另两枚（飞行中翻越／warn 模）仍绿 ⇒ 用例分工正确，没有互相掩盖 | `E71CB0CE0C8699CD` ＝基线 ✅ |

**还原总核**：`git status --porcelain -- <15 枚交付件>` ＝**空**；六枚被改文件 `Get-FileHash` 与 `bak\` 逐枚 `OK`（`costbook_44.go`/`costbudget.go`/`pricing.go`/`loop.go`/`costdao_44.go`/`costmeter.go`）；还原后 36 枚名册再跑一遍仍 **36/36 PASS**（`logs/` 尾段）。**⇒ 五发全是有效探针（能红）、且盘上零残留。**

---

## 5. 零回归

| 尺 | 读数 |
|---|---|
| `go test ./internal/observe ./internal/memory ./internal/agent -count=1 -timeout 20m`（**共享工作树**） | ⚠ **第一发不是 44 的红**：`internal/plugin` **未跟踪**的在飞面（50-r1c 正在写 `manifest.go`/`tamper.go`/`validator.go`，先后报 `undefined: ValidationErrors`、`undefined: strings`）把依赖它的 memory/agent 拖成 `[build failed]`；`go list -deps` 证＝memory／agent 依赖 `internal/plugin`，**observe 不依赖**（observe 始终 `ok 4.187s`）。该面本腿零碰、44-r1 也零碰（`git log -- internal/plugin/manifest.go` ＝ 从未提交）。 |
| 同尺，**隔离 HEAD 副本**（`git archive HEAD \| tar -x` 到 `44v1\head`，`HEAD=547d9da00`，未跟踪他腿面天然不在里面） | `ok observe 4.932s / ok memory 23.127s / ok agent 4.911s` ＝ **三包全绿、FAIL=0**（`logs/23-headcopy-fullpkgs.txt`）✅ |
| `go vet ./...` ＋ `go build ./...`（同一隔离副本） | **均 rc=0** ✅ ⇒ `cmd/wisp` 编译面也没被 `loop.go`/`cost.go` 的改动打破 |
| `gofumpt.exe -l`（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）对 §5 十三枚 .go | **输出空** ✅（含被本腿突变又还原的三枚＝还原后仍格式干净） |
| D22 门 `tools/d22scan -root <隔离副本>` | `clean - no D22 ban violations`；`ban#8 internal/=769 Go files`、`cmd/=186`；rc=0 ✅ |
| `DEFERRED(` 计数 base vs HEAD | `32` vs `32` ＝ **零增删** ✅（与 `SPEC-12 §5` 1:1 那条没被破坏） |

---

## 6. 纪律（逐字核）

| 项 | 尺 | 读数 |
|---|---|---|
| 两枚 `--stat` 与题面一致 | `git show --stat 8827e54ca`／`dc88642e0` | 15 枚（14 代码／测试跨 3 包＋1 证据件）／1 枚证据件追加 ✅ |
| 写面仅 observe/memory/agent | `git diff --name-only` 两枚各自的文件名，剥掉 `^internal/(observe\|memory\|agent)/` 与 `^docs/evidence/s1/` | 余集＝**空** ✅（`docs/evidence/s1/44r1-impl.md` 是裁决表落点，`AGENTS.md` §1.5 允许） |
| ⛔ `thresholds.go` | `git rev-parse 8827e54ca^:…` vs `HEAD:…` | blob **同为 `e2677b11b5a5adff8b4eb87c36d31286a274471d`** ＝ **一字未动**；其最后一次提交是 `00bbb76fe`（票 66），早于 44-r1 ✅ |
| ⛔ golden | `git diff --name-only base..HEAD -- "*testdata*"` | **空** ✅；AC#1 用的 `internal/agent/testdata/golden/tool-then-text.sse` 最后提交 `ae23da673`（票 10）✅ |
| ⛔ SLO | `git diff --stat base..HEAD -- docs/SLO.md internal/observe/thresholds.go internal/slo` | **空** ✅（`docs/SLO.md` 此仓不存在，属 `AGENTS.md` §4 已注明的索引过期，不是本票丢件） |
| ⛔ `cmd/wisp`／`internal/tools`／`internal/panel`／`frontend`／`design` | 同写面剥集 | 全零 ✅（22-r1 在写 `internal/tools`、296-r1 在写 `cmd/wisp`，本腿与 44-r1 都没碰，也没误判成对方的红） |
| ⛔ 改票面／台账／push | `git log -- .scratch/wisp/issues/44-costmeter-c23.md` 最后提交＝建票那两枚，早于 44-r1；两枚 44-r1 commit 对 `docs/reports/pending-and-issues.md` 的点名 diff＝**空**（range 里那几处台账改动全来自编排者 `96e8cb278`/`f59327fdd`/`a18ff4546`/`61046b65a`，已逐枚归属核过）；票面 6 枚勾框**一枚未翻**、`Status` 仍 `ready-for-agent` ✅；两枚 commit 题面均写"只 commit、⛔ push" ✅ |
| EOL | `git ls-files --eol` 交付件 | 全 `i/lf w/lf` ✅ |
| **`dc88642e0` 那枚自纠属不属实** | `git status --porcelain -- internal/agent/tools.go internal/agent/approval/pending_read.go`＝空；`git ls-files --eol` 读 `i/lf w/crlf attr/text eol=lf`；`gofumpt -l` 仍点名这两枚；两枚 44-r1 commit 都不含它们 | **属实** ✅ —— 首枚证据件把这两枚写成"他腿的未格式化"确实不实，真因是工作树 CRLF 违 `.gitattributes *.go eol=lf` 的既有债（与"哪条腿在飞"无关），追加更正把 §6 改对了。这枚追加合规：⛔ amend、只追加、只动自己的证据件。 |
| 编排者侧缺口（不是 44-r1 的账，转告） | `Test-Path .scratch/wisp/probes/orch/20261102/txt/44r1-dispatch.txt` | **仍 False** ⇒ r1 上报的"先存档再发没做到"**今天也没补**，`20261102/txt/` 里只有 `244c1/244r1` 等。 |

---

## 7. 是否同意结 44

**不同意现在整票结案。**（按派单要求具名，⛔ 放水）

**同意的部分**：票面 **AC#1、AC#2、AC#4 三格判 PASS**，AC#6 判 PASS（记录载体降级已具名）；派单三格 ①②③ 判 PASS，④ 零回归 PASS，⑤ 纪律 PASS。36/36 点名全绿、五发独立突变五发能红且逐字节还原、写面与禁区逐字干净、`thresholds.go` blob 逐字相同。**交付是真的**：真 SQLite＋真 golden 回放＋"一字节没出机器"的 `requests()==0`，不是 mock 顶包；`dc88642e0` 那枚自纠我也独立复核为属实。

**不同意结的三处具名欠账**：

1. **AC#5 整格未交付**（per-model 日/月配额闸、`plan_credit_total_micro` plan-credit 分账、"remaining 0 → 跳过派发＋链推进＋可见事件"三样盘上都没有）。**写腿的停手上报我判为正确**——我独立复跑确认 `task_log` 11 列无 provider/model、`cost_daily` 5 列、`schema.go` 顶部逐字禁增删列，补这维＝改 D35 契约文字＝机主批准面。但"上报得体"≠"格子满足"：**这一格要么由机主先裁 D35 加列/加表再补一票，要么由编排者把 AC#5 正式拆出本票**，二选一之后 44 才谈得上结。⛔ 不许直接翻 AC#5 的框。
2. **AC#3 只满足四分之三**：80% 预警今天**没有事件出口**（成本路径只发 100% 被拒那一枚 `EvError`，warn 态 `preDispatchBudget` 返回 nil ⇒ 总线无消息；§3.3 第 2 项那把尺可逐字重跑）。票面 `:21` 与 AC#3 都点名要 warn **event**。这一寸要么补一个 warn 事件的发布点（仍在 `internal/agent` 写面内、⛔ 不碰装配），要么具名记给票 40／装配票。
3. **产品今天不咬**：`cmd/wisp/run.go` 不填 `Config.Price`／`[cost]`，非测试代码里 `SetCostBudget`/`BudgetFrom` **零调用者** ⇒ 真跑价卡恒零、预算恒 uncapped。这是 ⛔ 派单禁面，不算 44-r1 违令，但**结 44 之前该由编排者明确把"装配"那一发排出去**（r1 §7 已经把口子名字都起好了：`BudgetFrom`＋`Loop.SetCostBudget`＋`CostMeter.SetBudget/Rebuild`）。

另附两处**不挡结票但要点名**的小账：`internal/agent/cost_44_test.go:14-19` 那段头注释把 tokens 说成 max-merge、把钱算成 460，与同函数体断言的 `{330,24,150}/750` 矛盾（断言对、注释错，归 44-r1 顺手一发）；派单存档 `44r1-dispatch.txt` 编排者侧至今仍未补。

**给编排者的一句话**：44 的**算法与聚合那半可以收**（AC#1/#2/#4/#6 ＋ ③ 阈值行为，证据硬、突变能红、还原干净），**卡口在 AC#5 的 D35 那一寸**——那是机主批准面，不是任何写腿能自己拍的东西。

---

*本表由非实现者腿 44-v1 产，只 commit 本件、⛔ push。*
