# 44-r1 — CostMeter (C23) 落地腿实现证据（逐格 AC ＋ 尺 ＋ 读数）

- **腿**：票 44 落地腿 44-r1（实现腿／写腿）。**日期**：2026-11-02。
- **票面**：`.scratch/wisp/issues/44-costmeter-c23.md`（58 行，`Status: ready-for-agent`，前置 04／10 已 `-done`）。
- **起手 HEAD**：`dcb2a28a3`（`票 391 验收腿 391-v1`）。**索引**：`git diff --cached` 空 ⇒ 干净。
- **写面冲突核（起手那发）**：`git status --porcelain` ⇒ 跟踪文件只有 ` M internal/risk/pathresolver_budget_norace_test.go`（＝在飞腿 **86-r1** 的写面，本腿**未碰未误判**）；未跟踪的他腿面＝`internal/tools/web*.go`＋`launcher_*.go` 一族。**本腿零写这些面**（归因见 §6）。
- **★写面所有权（本腿干活期间盘上自己长出来的旁证，非本腿自述）**：本腿收尾时 `git rev-parse --short HEAD` ＝ **`6e9d12c3b`**（起手 `dcb2a28a3` ⇒ 期间进了 18 枚他腿 commit），而 `git --no-pager diff --name-only dcb2a28a3..HEAD -- internal/observe internal/memory internal/agent` **输出为空**＝**没有一枚 commit 碰本腿三包**。同批台账逐字：`5873aa3e8 编排者A1419：45-r1停手(写面internal/observe与44-r1同包,串行重派)`、`d457fcbac 票 31 落地腿 31-r1（写腿）＝停手上报（零产码）：写面两枚包均被他腿占`（那两枚＝memory／tools）⇒ **internal/observe＋internal/memory 这一波归本腿独占**，45-r1 与 31-r1 串行在后。
- **收尾复跑**：`go vet ./...` rc=0、`go build ./...` rc=0、`go test -count=1` 六包全 ok（`logs/16-vet-all-final.txt`／`logs/17-tests-final.txt`）＝本腿交付面对新 HEAD 仍绿。
- **只 commit、⛔ push**；显式 pathspec 逐枚点名（§6）。临时件只建不删：`D:\tmp\wisp\20261102\44r1\`（`logs/01…14`、`logs/mutation.txt`、`mutate-44r1.ps1`）。

## 0. 起手三核的第 19／20 处（票面点名物与引文出处）

| 核 | 票面点名的串／行／文件 | 现跑判 |
|---|---|---|
| 19 | `:17` 逐字 `cost_daily upsert` | 在 ⇒ `internal/memory/dao_misc.go:130-151` `BumpCostDay`（upsert + add） |
| 19 | `:19-20` 逐字 `"价格未知" marker (never fabricate amounts)` | 在 ⇒ 本腿落 `observe.PriceUnknownMarker`（常量值逐字 `价格未知`，`TestAC4UnknownCardPricesNothing` 钉字符串） |
| 19 | `:21` 逐字 ``Budgets `[cost]`: daily/monthly caps; ≥80% → warn toast + badge; 100% → PAUSE new tasks`` | 在 ⇒ `internal/config/schema.go:556-564` `CostSection{DailyBudget,MonthlyBudget,AlertThreshold default:"0.8",OverBudget default:"pause"}` 四枚叶子在场 |
| 19 | `:27` 逐字 `All math integer micro-units; reconciliation test vs task_log sums exact` | 在 ⇒ 本腿 AC#2 格对账三枚尺（§2 行 2） |
| 20 | `:8` `Spec refs: SPEC-05 §8, §14.1, C23, D35 cost_daily, S6` | SPEC-05 §8 在（`docs/specs/SPEC-05-agent-core.md:164-173`，`:170-172` 逐字含 `日/月累计（cost_daily）`／`预算达 80% 提示、100% 默认暂停新任务（[cost] over_budget 可配仅告警）`／`单任务 200k 预算（防死循环）与费用预算（钱）是两个层级，都要`）；PLAN.md:2711 `cost_daily` 行在；C23＝PLAN.md:1379；S6＝PLAN.md:1435/1449 在 |
| 20 | `:24` `Distinct layers documented…` | 行 24 逐字就是那句（不是 24-25）⇒ 本腿代码注释按 24 引（先前草稿写 24-25 已改） |

**派单存档不在盘上（上报，不影响本腿开工）**：派单笔 `D 单` 指 `.scratch/wisp/probes/orch/20261102/txt/44r1-dispatch.txt`，现跑 `Get-ChildItem .scratch/wisp/probes/orch/20261102/txt` ⇒ 该目录 15 枚里**没有 44r1-dispatch.txt**（最新是 `395v1-dispatch.txt`）。派单正文在腿的收件里是完整的（射程／禁面／尺都点名），所以本腿按收件正文干活，**但"先存档再发"这一条在编排者侧没做到**，具名报在这里。

---

## 1. 射程与落点（为什么是 `internal/observe`）

票面没写包名，本腿按盘上三处指向定落点，⛔ 自造新模块：

1. `docs/specs/SPEC-01-architecture.md:40` 逐字 `（【SPEC】由冻结模块表 1:1 映射，不得增删模块）`，`:67` 逐字 `│   └── observe/  # slog JSONL、脱敏、诊断包、SLO 采样、CostMeter(C23) 落账`；
2. `internal/observe/doc.go:3` 现读逐字 `CostMeter (C23) accounting`；
3. 票前 HEAD 读数（尺＝`git --no-pager show HEAD:internal/agent/cost.go`，本腿已把这段改掉，故标〔票前〕）：`:12-14` 逐字 `The daily/monthly aggregation (cost_daily) and the budget-pause policy belong to the CostMeter ticket (44)`；旁证另有一枚现读在位的：`docs/audit/2026-10-25-frontend-ia-design.md:97` 逐字 `cost_daily DAO 已备、聚合未做=票 44`。

⇒ **成本计所在包＝`internal/observe`**；`cmd/wisp/**` 是派单点名禁面（本腿零写）；`internal/agent`／`internal/memory` 只加本票需要的入账点与读口（§5）。

**依赖方向锁死了接口形状**：`internal/memory` 已 import `observe`（`dao_misc.go:196` 用 `observe.NowWallUTC`、`dao_tasklog.go:9`），`internal/config` 也 import `observe`（`settings.go:183` 用 `observe.New`），`internal/llm` 同样（`events.go:130` 用 `*observe.Error`）⇒ **observe ⛔ 反向 import memory/config/llm**。所以 `CostMeter` 对存储只认三枚窄接口（`observe.DayLedger/DayReader/TaskReader`＋修复口 `DayWriter`），由 `*memory.Store` 结构性满足（`internal/memory/costdao_44.go:122-127` 编译期钉死），与 `internal/agent/journal.go:18-26` 的 `Journal`＋`var _ Journal = (*memory.Store)(nil)` 同形。

**日界历法**：`day` 桶＝注入历（默认 `time.Local`，测试注入 `FixedZone`）。理由写进代码注释：日预算是人说的"今天"，⛔ 让"今天"在 08:00 翻篇。落盘时间戳仍一律 wall UTC 秒（D42#9）。**残余（不摆机主、只具名）**：若要把日界做成配置项（如 `[cost] day_timezone`）＝动 D36 配置模型文字＝契约面，本腿不动。

---

## 2. 逐格 AC（尺 ＋ 读数 ＋ 判）

尺都可逐字重跑；读数取自 `D:\tmp\wisp\20261102\44r1\logs\01-observe-44.txt`／`02-memory-44.txt`／`03-agent-44.txt`／`06-full-pkgs.txt`。

| AC | 判 | 尺（重跑命令） | 读数 |
|---|---|---|---|
| **#1** 多轮 golden 任务：per-task tokens/cost 与手算夹具**逐单位**相等 | **达成** | `go test ./internal/agent -count=1 -v -run TestAC1GoldenPerTaskCostMatchesHandComputed` | `--- PASS`（真 SSE 回放 `testdata/golden/tool-then-text.sse`＋真 C5 adapter＋真 SQLite store）。手算写在用例里：turn1 `{150,14,0}`→440，turn2 `{180,10,150}`→310（缓存 150 走 cache 价、非缓存 30 走 in 价），合计 **750 micros**；断言 `res.CostMicros==750` **且** `task_log` 行 `{330,24,750}` 同数；夹具自身先自检（`wantMicros!=750 ⇒ Fatalf`）＝手算动了也会当场喊 |
| **#1** 同上（数学面） | 达成 | `go test ./internal/observe -count=1 -v -run 'AC1PriceMicros'` | 2 枚 PASS：六个手算子形（含 `1 token×2.5e6 → 2.5 → half-up = 3`、`cached>in` 夹到 `in`、`Cached=0 走 In 价`）＋整数溢出**饱和不回转** |
| **#2** 日翻篇＋月翻篇（夹具时钟）；与 task_log 对账精确 | **达成** | `go test ./internal/observe -count=1 -v -run 'AC2Today\|AC2DayRollover\|AC2MonthRollover\|AC2Reconcile'`；`go test ./internal/memory -count=1 -v -run AC2`；`go test ./internal/agent -count=1 -v -run 'AC2Reconciles'` | observe 4 枚 PASS：23:59→00:00 同一 store 只换时钟 ⇒ 两枚 `tasks=1` 分家、新日 `micros=11` 不吸收昨日 500；11-30→12-01 ⇒ 月范围 `2026-11-01..2026-11-30`→`2026-12-01..2026-12-01` 且**旧日行逐字未变**；对账先 AGREES 再植一枚"只进 base 不进 aggregate"的游离行 ⇒ `Agree=false` 且数字对得上（590 vs 632），`RefreshDay` 修复后再 AGREES。memory 3 枚 PASS（真 SQLite：范围过滤/排序/畸形端点拒绝；`ended_at` 首秒含末秒含、次日首秒不含、`ended_at IS NULL` 不含）。agent 1 枚 PASS（环路跑两发真任务后对账 AGREES，`tasks=2`） |
| **#3** 80% 告警事件；100% → 新任务被拒＋**显式文案**＋运行中的任务跑完；warn-only 备选模式核过 | **达成**（界面侧 toast/badge 归票 40，见 §4 残余） | `go test ./internal/observe -count=1 -v -run 'AC3Budget\|AC3OverBudget\|AC3Binding\|AC3Unset'`；`go test ./internal/agent -count=1 -v -run 'AC3'` | 阈值是**整数基点**（0.8 只在 `BudgetFrom` 边界转一次成 8000）：789→ok／799→ok／**800→warn**／999→warn／**1000→over＋拦**／1500→over。真环路三枚 PASS：① `TestAC3NewTaskRefusedBeforeTheProvider` ⇒ 被拒任务 `Status=failed`、`Err.Class=budget`（D35 枚举行 31）、文案含 `暂停新任务`＋`运行中的任务不受影响`＋恢复时点，**`second.requests()==0`＝一字节都没出机器**，且 `task_log` 有行（`error_class='budget'`）、对账仍 AGREES（`tasks=2`、钱没涨）；② `TestAC3RunningTaskFinishesPastTheLimit` ⇒ 用 sink 在**首个 text delta 时**把日预算压到 1 micro（飞行中翻越），该任务仍 `completed`（闸只在入口读一次，源码 `loop.go` 注释逐字写明"nothing below this line consults the budget again"）；③ `TestAC3WarnModeStillRunsTasks` ⇒ 超 100% 仍放行、provider 真被够到、`Notice` 自报模式 |
| **#4** 未知价 provider → tokens-only ＋ marker（⛔ 造数） | **达成** | `go test ./internal/observe -count=1 -v -run 'AC4Unknown\|AC4Table\|PricingRejects'`；`go test ./internal/agent -count=1 -v -run 'AC4Unknown'` | 全零卡 ⇒ `PriceMicros==0`（**不是**"免费"：`Price.Known()` 为 false；真免费必须由行 `billing:"free"` **明说**，`TestPricingRejectsFabricatedAndMalformedRows` 把 9 种伪造/畸形形状逐枚拒掉，含"没 source 的价"）；marker 常量逐字 `价格未知`；环路侧 `config.Price{}` ⇒ `res.CostMicros==0` 而 `res.Usage` 三枚仍在，日聚合仍入账（`tasks=1 / micros=0`） |
| **#5** 配额：per-model 日/月闸；plan-credit 记账（手录总额−消耗，"估算"marker）；**派发前防负**（remaining 0 → 跳过、链推进、可见事件）；plan vs pay-per-token 分账 | **停手上报（不动手）** | 尺：`grep -n "CREATE TABLE" internal/memory/schema.go` ⇒ `task_log:65 / cost_daily:114`；再 `sed -n '65,77p;114,120p'` | **结构上量不到 per-model**：`task_log` 12 列里没有 provider/model 列，`cost_daily` 只有 `day/tokens_in/tokens_out/cost_micros/tasks`，`provider_health.quota_state` 只是文本态（`:174`）⇒ per-model 日/月花费**无法从盘上任何持久面还原**。补这一维必须给 D35 冻结表**加列或加表**＝`SPEC-02 §3` 逐字同步面＋`internal/memory/schema.go:11-15` 的"⛔ added/removed columns，introspection 测试逐字节对契约"⇒ 属"改契约文字"（AGENTS.md §0 第 2 条／`SPEC-12 §4.1`），**本腿⛔自改**。另两半：`over_budget` 全局那一层本腿已落（AC#3），per-model 层与 `plan_credit_total_micro` 分账等这一格裁定；链推进归票 11（票面 `:43` Out of scope 逐字 `chain stepping logic (11)`）。**已在场的接口位**：`observe.BudgetParams` 是 per-cap 值对象（不是单例），`CostMeter.Judge` 收两个 `DayTotal` ⇒ 拿到持久面后可直接按 model 复用，不需重写闸 |
| **#6** 价格表版本 bump → 记录在案，旧聚合不变 | **达成（记录面见 §4）** | `go test ./internal/observe -count=1 -v -run AC6Embedded`；`go test ./internal/agent -count=1 -v -run AC6PriceTable` | 版本＝`pricing.json` 的 `version`（现值 1），`ParsePriceTable` 拒 `version<1`；`CostMeter.RecordPriceTableVersion()` 落一条 slog（`version/updated/models/recomputes_aggregates=false`）；环路真 store 用例：入账后 `CostDay` 读数 `before`，把表 bump＋1 并 `Rebuild` ⇒ 版本随表走、`after == before`（struct 相等），再 `RefreshDay` 也只从 `task_log` 重算 ⇒ `healed.CostMicros` 仍是入账那一笔，**没有任何路径按新价重算历史** |

---

## 3. 反向判据（摘实现／挪阈值 ⇒ 点名用例红）＋还原核

脚本 `D:\tmp\wisp\20261102\44r1\mutate-44r1.ps1`（M1-M6）＋同目录 `logs/mutation.txt`（M2p/M7/M8 追加）。每发＝备份→改一处→跑点名用例→还原→再跑。

| 发 | 动了哪一处（逐字） | 红在谁身上（具名） | 还原 |
|---|---|---|---|
| M1 | `costbudget.go` `case spent >= cap:` → `case spent > cap:` | `TestAC3BudgetBoundaries79And80And100`：`spent 1000: level "warn", want "over"` ＋ `AllowNewTask true, want false` | `restored exit=0 ok` |
| M2 | `costbook_44.go` 入账提前 return（PowerShell 反引号转义把源码写坏，构建失败＝无效探针，作废重做） | 无效（build failed），不采信 | 还原后 `ok` |
| M2p | `costbook_44.go` `if l.costMeter == nil {` → `if true {` | `TestLoopBooksFinishedTaskIntoCostDay`：`cost_daily[2026-11-02] = {…Tasks:0}, want exactly one booked task`；`TestAC2ReconcilesAgainstTaskLog`：`agg={…Tasks:0} base={TokensIn:660 TokensOut:48 Micros:1500 Tasks:2}`；`TestAC4UnknownPriceKeepsTokens` 同族红 | `restored exit=0 ok` |
| M3 | `pricing.go` `q := (num + million/2) / million` → `q := num / million` | `TestAC1PriceMicrosMatchesHandComputedFixture`：`half-micro total rounds half-up: PriceMicros = 2, want 3` | ok |
| M4 | `costdao_44.go` `WHERE day BETWEEN ? AND ?` 删掉 | `TestAC2CostDaysBetweenFiltersAndOrders`：`range returned 4 rows (…2026-10-31…, …2026-12-01…), want exactly the two days inside it` | ok |
| M5 | `SetCostDay` 的 `tokens_in=excluded.tokens_in` → `tokens_in=tokens_in+excluded.tokens_in` | `TestSetCostDayOverwritesAndRefusesNegatives`：`SetCostDay did not replace the row: {TokensIn:130 …}` | ok |
| M6 | `PriceMicros` 的 `if cached > t.In { cached = t.In }` 删除 | `TestAC1PriceMicrosMatchesHandComputedFixture`：`cached cannot exceed input: PriceMicros = 9223372036854775807, want 50`（＝饱和路径同时被验到） | ok |
| M7 | `preDispatchBudget` 判空反向（闸永不拦） | `TestAC3NewTaskRefusedBeforeTheProvider`：`costgate_44_test.go:111: status = completed (), want failed: 100% of a pause-mode cap must reject a new task` | `restored exit=0 ok` |
| M8 | `PostTask` 的 `if micros < 0 {` → `if micros == 0 || micros < 0 {`（零成本不入账） | `TestAC3NewTaskRefusedBeforeTheProvider`：`the refusal broke the reconciliation invariant: agg={…Tasks:1} base={…Tasks:2}` | ok |

★M8 值得留一句：它证明"被拒的那一笔也必须入账"不是随手写的断言，而是**对账不变式的承重件**（否则 `cost_daily.tasks` 与 `task_log` 永远差被拒的枚数）。

---

## 4. 本腿交付的口径边界（写明，不含糊）

- **告警事件（toast/badge）**：Go 侧出口＝`observe.BudgetVerdict{Level,Notice,Caps}`＋被拒时 `publish(Event{Kind:EvError,Text:Notice})`。渲染与角标是 `frontend/**`＋票 40 的面（派单⛔面＋AGENTS.md §1.5）⇒ **本腿不勾"界面可见"那一半**，只保证它有真数据与真事件。
- **版本"记录在案"的载体**：slog JSONL（诊断面）＋内存里的 `PriceTable.Version`。`cost_daily`/`task_log` **没有**可放版本号的列 ⇒ 逐日持久钉版本同样要动 D35（与 AC#5 同一枚墙，具名不自改）。本腿因此把"记录"落在不改契约的那一面，并把 `recomputes_aggregates=false` 显式写进那条记录。
- **`pricing.json` 现收 0 枚价**：票面 `:20` 逐字 `never fabricate amounts`、`:42` 把"价格研究"划到 out of scope，而本机取不到任何 provider 价卡（`web_fetch https://api-docs.deepseek.com/quick_start/pricing` ⇒ `resolves to a non-public IP address`）。⇒ 本腿交**机制**（版本／校验／覆盖合并／marker），行由机主按公开价卡填，**没有**任何编造的数。覆盖路径已核：`WithOverrides` 与表内行同一套校验（要 currency、要 source、⛔ 负价），且**不改接收者**（`TestAC4TableLookupMissesSayUnknownAndCostNoMoney` 钉住）。
- **单一算钱口**：`internal/agent/cost.go` 的 `PriceMicros` 改为**逐字委托** `observe.PriceMicros`（整数），环路/`guard`/`cmd/wisp costLine` 三个嘴读同一笔数；`CardOf` 把 D36 冻结的 `float64` 价卡**在边界一次**换成整数微单位（注释写明：那是"率"不是"钱"，钱全程 int64）。
- **`cmd/wisp` 侧残余（⛔ 本腿写面）**：`cmd/wisp/run.go:1342-1349` 装配仍不填 `Config.Price/Currency`（台账 `docs/reports/pending-and-issues.md:11218` 的 261-p1 段已具名同一件事，逐字 `run.go:1010-1017 装配从不填 Price/Currency`；行号随文件长大已移到 `:1342-1349`，本腿现读到的是那一处）⇒ `wisp run` 真跑时价卡恒零＝走 marker 路，`[cost]` 四枚叶子也还没进 `agent.Config`（本腿提供 `observe.BudgetFrom`＋`Loop.SetCostBudget` 两个口，等 cmd/wisp 那批）。这是**接线**那一寸，不是算法那一寸。

---

## 5. 改动文件（逐枚点名）

**新增（12 枚代码／测试＋1 枚本证据件）**

| 文件 | 行 | 是什么 |
|---|---|---|
| `internal/observe/pricing.go` | 326 | C23 价卡／整数微单位算法（`PriceMicros`，`bits.Mul64`＋饱和）／版本化价格表（embed）／config 覆盖合并／`价格未知` marker 与其校验 |
| `internal/observe/pricing.json` | 6 | 起始价格表（version 1、models `{}`、notes 写明填法与 ⛔ 造数） |
| `internal/observe/costmeter.go` | 426 | `CostMeter`：日／月聚合读、`DayKey`/历法、入账（`PostTask`）、对账（`Reconcile`）、修复（`RefreshDay`）、三枚接口＋`DayWriter` |
| `internal/observe/costbudget.go` | 230 | `[cost]` 钱预算：`BudgetParams`/`BudgetFrom`、整数基点阈值、`Judge` 的 ok/warn/over 与 pause/warn 两模、显式文案、溢出饱和 |
| `internal/memory/costdao_44.go` | 128 | `CostDaysBetween`／`SumTasksEndedBetween`／`SetCostDay` ＋编译期"Store 是 meter 的嘴"四枚断言 |
| `internal/agent/costbook_44.go` | 154 | 环路入账装配（`costMeterFor` 断言发现，⛔ 不宽化 `Journal`）、`bookCostDay`、`CostReport/CostBudget`、`preDispatchBudget` 派发前闸、读失败为何放行（写进注释） |
| `internal/observe/pricing_44_test.go` | 259 | AC#1/#4/#6 数学与表（7 枚） |
| `internal/observe/costmeter_44_test.go` | 533 | AC#2/#3 翻篇／对账／阈值（15 枚） |
| `internal/memory/costdao_44_test.go` | 265 | AC#2 真 SQLite 三枚＋`SetCostDay` 一枚 |
| `internal/agent/cost_44_test.go` | 120 | AC#1 golden 逐单位一枚＋AC#4 一枚 |
| `internal/agent/costbook_44_test.go` | 285 | 入账／对账／闸读数／marker／版本（6 枚） |
| `internal/agent/costgate_44_test.go` | 202 | AC#3 三枚（飞行中翻越／入口拒绝且零出网＋可审计＋对账不破／warn 模式） |

（数＝`git status --porcelain` 现跑：`??` 13 枚＝本表 12 枚代码／测试件＋本证据件；`M` 2 枚＝`cost.go`／`loop.go`。）

**修改（2 枚，`git diff --numstat`）**

| 文件 | numstat | 改了什么 |
|---|---|---|
| `internal/agent/cost.go` | `54 14` | `PriceMicros` 改为委托 observe 整数算；新增 `CardOf`/`TokensOf`/`wholeMicros`；`Cost` 加 `Priced bool`（marker 的来源位）；顶部注释改成"单一算钱口"那段 |
| `internal/agent/loop.go` | `30 0` | `Loop` 加 `costMeter` 字段并在 `New()` 里装；`run()` 在 `startTask` 之后插派发前预算闸；`finish()` 在 `finishTask` 成功分支里入账当天聚合（失败不入，注释写明为何） |

**未碰的禁区（现跑自证）**：`internal/observe/thresholds.go`（`git status` ⇒ 无本腿改动；它是 SLO 阈值，`ObserverCost` 那几行是"采样器自身开销"，**不是**预算阈值——所以本票没撞禁改面）；`testdata/golden/**` 零改（只读回放）；`internal/panel/**`、`frontend/**`、`design/**`、`cmd/wisp/**` 零写；`C1–C32`／`D1–D47`／D43 转移表零改；`DEFERRED(` 标记零增删（`internal/observe/doc.go:22` 那条 `DEFERRED(JSONL sink/redaction/diagnostics/CostMeter)` 本腿**没动**——它自称"implemented by ticket 08"而 44/45 各管其中两件，台账 `225-deferred-registry-audit.md` 第 13 行已把这条判成〔标记本身该摘／改写〕，摘标记属人批面）。

---

## 6. 门禁读数（现跑）

| 尺 | 读数 |
|---|---|
| `gofumpt.exe -l`（绝对 `/d/work/base/gopath/bin/gofumpt.exe`）对 §5 的 13 枚 | 输出**空**＝干净（`logs/10-gofumpt.txt`）。全目录跑只剩他腿的 `internal/agent/approval/pending_read.go`、`internal/agent/tools.go` 两枚未格式化 ⇒ **不是本腿面** |
| `go vet ./...` | rc=0，无输出（`logs/11-vet.txt`）⇒ 含 `cmd/wisp` 测试面编译干净 |
| `go build ./...` | rc=0 |
| D22 门 `sh scripts/d22scan.sh` 的两步（正控测试另计）：`cd tools/d22scan; go run . -root <repo>` | `clean - no D22 ban violations`；`internal/=766 Go files`（含注释与 `_test.go` 的 ban#8 覆盖）`cmd/=186`；rc=0（`logs/12-d22scan.txt`） |
| `sh scripts/check-path-length-budget.sh` | `VERDICT GREEN`（18338 tracked paths、over-budget=62 全在册、roster＝tree） |
| `go test -count=1` 六包 | `ok observe 5.447s / ok memory 25.347s / ok agent 6.280s / ok llm 47.478s / ok config 2.339s / ok agent/approval 0.893s`（`logs/14-final-tests.txt`） |
| `-race`（本腿新代码路径） | `ok observe 1.110s / ok agent 3.158s`（`logs/04…05-race-*.txt`） |

**红在别处的两包（归因，⛔ 当自己的活、⛔ 顺手修；两枚都是"他腿在飞的未跟踪面"，本腿零写）**：
1. `internal/tools`：19:5x 那发 `go test` 报 `[build failed]`（`webmd.go:178 undefined: extractAttr`／`p.links undefined`）；20:2x 复跑 `go build ./internal/tools` **rc=0**（他腿补上了），但 `go vet ./internal/tools` 仍 **rc=1**，红在它的**测试件**：`app_launch_windows_test.go:100:81: undefined: windows.STILL_ACTIVE`。该族文件全是未跟踪态（`?? internal/tools/web*.go`、`launcher_*.go`，票 22／24／385 一族在飞）⇒ 不是本腿造成，也不是本腿能判终态的面。
2. `internal/panel TestC21DesignTokensFourWayAgree` FAIL（20:2x 复跑仍红，`logs/15-panel-c21-foreign-red.txt`）：红句逐字 `docs/evidence/s1/c21-native-tokens.md:58: the table names Palette.GlassRing, but internal/ball/tokens.go holds no such colour literal in DarkPalette()`，同族还有 `:59/:60`（`GlassHi`／`GlassHiSoft`）与 Light 侧 `:103/:104/:105` 三条，收尾两行 `checks 1-4 reconciled 22 of 22 mapped pairs (78 colour rows read, 72 agreed with internal/ball/tokens.go)` ⇒ C21 设计 token 四面（frontend／evidence／ball／css）名册不一致，与本票包面**无交集**（本腿零写 panel/ball/frontend/design）。

---

## 7. 残余（下一批发腿可直接接的清单）

1. **AC#5 per-model／plan-credit 配额**：先解"持久面没有 provider/model 维度"（D35 加列或加表＝人批），本腿的 `BudgetParams`/`Judge`/`DayTotal` 三件套已可按 model 复用。
2. **`cmd/wisp` 装配**：`run.go:1342` 填 `Config.Price`（从 `observe` 价卡取）＋把 `cfg.Cost` 四枚叶子经 `observe.BudgetFrom` 交给 `Loop.SetCostBudget`；`[cost]` 热重载走 `CostMeter.SetBudget/Rebuild`（已为此上锁，`mu sync.RWMutex` 只护 budget/table）。
3. **票 40 成本页**：读 `CostMeter.Today/MonthToDate/Range/Reconcile` ＋ `BudgetVerdict.Caps`（两行都在，⛔ 只给"咬住的那一枚"）。
4. **票 60／61 音频价**：`Price.AudioIn/AudioOut` 已在 config，`observe.Price` 尚无 audio 两枚叶子 ⇒ 到那两票一起加，别在这里先造。
5. **票 145 行 12（成本快照）**：面板字段面。
6. 若机主要"逐日持久钉价格表版本"：与 AC#5 同一枚墙（要列）。

**勾框未翻**：票面 6 枚 `- [ ]` 本腿一枚都没勾（只写实现与证据）；§2 的"判"列是给非实现者验收腿的落点，AC#5 的停手上报需要编排者／机主先裁 D35 那一寸。
