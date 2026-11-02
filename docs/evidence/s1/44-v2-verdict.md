# 44-v2 — CostMeter (C23) 非实现者综合复验（裁决腿 ≠ 实现腿）

- **腿**：票 44 验证腿 44-v2 ＝ 非实现者综合复验（`AGENTS.md` §0 第 3 条／`SPEC-12 §4.3` #1/#3）。本腿零改码：⛔ 实现／票面／台账／push；§5 的突变探针跑完**逐字节还原**（hash 核见表格）。
- **被验对象**：补账三腿 ＋ 证据 —— `19b48978b`（44-r2 observe/quota_44.go，337 行新文件）· `5ae48abc9`（44-r3 AC#3 EvWarn 出口＋run.go 咬合，7 files）· `a11e63bee`（44-r4 agent/quota_44r4.go AC#5 接线，2 files 493 行）· `b14ee2a95`（44-r4 证据件）；底账仍含 44-r1 `8827e54ca`＋`dc88642e0`。
- **票面**：`.scratch/wisp/issues/44-costmeter-c23.md`（58 行、6 枚 `- [ ]` 勾框，**一枚未翻**）。
- **日期**：2026-11-02。**临时件**：`D:\tmp\wisp\20261102\44v2\`（logs/01…12、bak/ 四枚原件、cmdwisp.test.exe 等构建物）。只建不删。

---

## 1. 起手三核 ＋ commit 与题面 ＋ 纪律（逐字核）

| 核 | 尺 | 读数 |
|---|---|---|
| HEAD 现盘最新 | `git rev-parse HEAD` | 起手 `c00045680` ≥ 要求 `65768beca` ✅；全程未再前进 |
| 索引干净 | `git status --porcelain` 剥 `^??` | 空 ✅（只有未跟踪的 scratch/他腿面） |
| 盘＝HEAD（交付件） | `git diff --stat HEAD -- observe/memory/agent 新面＋cmd/wisp/run.go` | 空 ✅ |
| 四枚 `--stat` 对题面 | `git show --stat` 逐枚 | r2＝observe/quota_44.go 单文件 337 行 ＝题面「observe/quota_44.go」✅；r3＝7 files＝agent 事件出口 3 枚＋run.go＋cmd 测试＋证据件，与「EvWarn＋run.go 咬合」✅；r4＝quota_44r4.go＋测试 ＝「AC#5 接线」✅；b14ee2a95＝仅证据件 ✅。**⚠ 两处具名（不判违令）**：① r3 新增 `cmd/wisp/costwiring_44r3_test.go`——派单写面文字是「observe/memory/agent＋cmd/wisp/run.go」，这枚测试超出字面一寸，但它钉的正是 run.go 的咬合、不碰他腿文件，且 r3 题面自报；② r2 是 checkpoint，**无独立证据件**（读数并入 r3 证据 §配额段与 r4 证据），纪律上算薄。**⚠ 编排者侧仍欠**：`probes/orch/*/txt/` 里 44r2/r3/r4 派单存档＝零命中（44-v1 §6 已记 44r1 缺档，同一规矩没补）。 |
| 写面 | 四枚 `--name-only` 剥 `^internal/(observe|memory|agent)/`、`^cmd/wisp/run\.go$`、`^cmd/wisp/costwiring_44r3_test\.go$`（上格已具名）、`^docs/evidence/s1/` | 余集＝**空** ✅；`cmd/wisp` 本波只有 44-r3 写过（`git log -6 -- cmd/wisp` 首条即 `5ae48abc9`）⇒「每波一枚写腿」没破 |
| ⛔ thresholds | `git rev-parse` 五点对比（r1 base／r2／r3／r4／HEAD） | `internal/observe/thresholds.go` blob 全程 `e2677b11b5a5adff8b4eb87c36d31286a274471d`＝**逐字一字未动** ✅ |
| ⛔ golden | `git diff --name-only 8827e54ca^..HEAD -- "*testdata*"` | 空 ✅ |
| ⛔ 票面/台账/push | 票面最后提交＝建票两枚（早于全部 44 腿）；四枚 44 腿文件名单不含 `docs/reports/*`；题面各写「⛔ push」 | ✅（range 里台账改动归编排者 A1435–A1438） |
| DEFERRED | `git grep -o "DEFERRED(D…)" base..HEAD` 对比 | base 4 → HEAD 9，**新增 5 枚全部＝`DEFERRED(D34-1)`，归 24-r2a**（`-S` 单提交命中）；44 三腿**自加 0** ✅ |
| EOL | `git ls-files --eol` 四枚承重 .go | 全 `i/lf w/lf` ✅ |
| gofumpt | 对 9 枚交付/被改文件 | 输出空 ✅（含突变还原后的四枚） |
| d22scan（tools/d22scan 独立模块，`go run . -root ..\..`） | **1 处命中，不在 44**：`internal/tools/search_files.go:259` phantom-citation（引不存在的 `24r3-impl.md`）；`git ls-files` 证该文件**未跟踪＝不在 HEAD**，属票 24 在飞面（24-v1 在读，本腿⛔碰）。⇒ **44 交付面零命中 ✅；树级那枚红归 24**。 |

---

## 2. 真跑（派点名命令逐字）

### 2.1 定向尺

```
go test ./internal/agent ./internal/observe ./internal/memory -run 'Cost|AC5|Budget|Price|Quota|Wiring' -count=1 -timeout 25m
```

**FAIL=0** ✅（`logs/01`）。`-v` 名册（`logs/02`）＝**35 枚 PASS、0 FAIL、0 SKIP**，含全部新面：AC5×7、AC3-warn/under/refusal×4（r3）、AC1×3、AC4×5、AC6×2、`TestLoopBooksFinishedTaskIntoCostDay` 等。**照旧具名**（与 44-v1 §2.1 同一条尺病）：`AC2DayRollover/MonthRollover/Reconcile*` 不匹配这五词、`Budget|Quota` 又误捞七枚他票（`TestGrantCostPluginStateDAO`、`TestGoroutineBudgetDrains`、`TestD39SectionBudgets…` 等）⇒ **本票真尺＝下面 §2.2 全包回归＋名册**，编排者以后点票 44 建议直接 `TestAC[1-6]`。

### 2.2 三包全回归

```
go test ./internal/agent ./internal/observe ./internal/memory -count=1
```

`ok agent 4.718s / ok observe 4.378s / ok memory 23.007s` ＝ **三包全绿、FAIL=0** ✅（`logs/07`）。共享工作树今天没有他腿未跟踪面挡编译（v1 那次 `internal/plugin` 挡包的情形今天不存在）。

### 2.3 产品咬合真跑（cmd/wisp）

**环境事实先具名**：`go test ./cmd/wisp` 在本腿 shell 里**进程起不来**（`exit status 0xc0000135`＝STATUS_DLL_NOT_FOUND）——test 二进制经 cgo 面依赖仓根的 `onnxruntime.dll`／`sherpa-onnx-*.dll`（三枚**未跟踪** staged 件），`go test` 把二进制拷进临时目录、CWD 在包目录，Loader 找不到。⇒ 复现办法＝`go test -c -o` 后从 `cmd/wisp` 目录带 `PATH=<仓根>` 直跑（`logs/06`）：

```
--- PASS: Test44R3CostWiringCarriesPriceAndCaps
--- PASS: Test44R3CostWiringSaysUnknownPriceAndBadBudget
--- PASS: Test44R3RunGoAppliesTheWiring   (3/3 PASS, rc=0)
```

**边界一句**：CLI 端到端整跑（真 `wisp run` 出账）走 DPAPI＋真 HTTP 端点、无离线注入口，不属本腿可为——「咬合」由三枚钉死：run.go 调用点静态真引用（下格）＋装配层测试＋agent 包内真 golden 回放的预警/拒绝环路（r1 已验、今天随全包重跑绿）。

---

## 3. 逐格复验 AC#1–AC#6（今天自己的读数，非转抄）

| # | 票面逐字（缩） | 判 | 读数 |
|---|---|---|---|
| **#1** | per-task tokens/cost ＝ 手算夹具到单位 | **PASS** | 名册 `TestAC1×3` 今天绿；`PriceMicros` 单一算钱口、整型微单位、半微 up 进位由 M-尺（44-v1 V3）钉过，本次零改动文件（`pricing.go` 不在三腿名单）。 |
| **#2** | 日/月滚桶＋对账 task_log exact | **PASS** | 三包全回归覆盖 `TestAC2×7` 全绿（含 Day/Month/Rollover/Reconcile 真 SQLite）；三腿未动这些文件，盘＝HEAD 已核。 |
| **#3** | 80% warn event；100% 拒＋在跑任务收工；warn-only 备选 | **PASS（欠账已补）** | v1 的缺半＝warn 无事件出口。今天：`sink.go` 新 `EvWarn` 种类；`costbook_44.go:146-148` 发布点（`Metered && AllowNewTask && Level≠OK` ⇒ `publish(EvWarn, TaskID, Notice 逐字)`）；`run.go` consoleSink 有 reader（`[预警]`）；`TestAC3WarnBudgetPublishesWarnEvent`／`…FiresPastTheLimitInWarnMode`／`UnderTheAlertLinePublishesNothing`（反向防噪）／`RefusalStaysAnError`（100% 仍只 EvError 不双发）四枚今天全绿，且 **§5-M3 摘掉发布点⇒恰这两枚红、另两枚绿**＝事件出口有独立牙齿。100% 拒绝＋飞行中收工（r1 两枚）随全包重跑绿。**整格补齐。** |
| **#4** | 未知价→tokens-only＋「价格未知」marker | **PASS** | `TestAC4×5` 全绿；新增佐证：咬合层「无卡＝零卡＝Known false＝0 micro＋marker」由 `Test44R3CostWiringSaysUnknownPriceAndBadBudget` 今天真跑钉住（`go test` 语义下 3/3 绿）。 |
| **#5** | per-model 日/月闸；plan-credit 分账（手录总额−消费＋「估算」marker）；**防负夹具＝remaining 0→跳过派发＋链推进＋可见事件**；80%/100% 事件；plan/pay 分模型分账 | **PASS（格子按票面文字成立；三处归属具名 §4.1）** | 三层闸：`observe/quota_44.go` `JudgeModel`＝daily/monthly/plan-credit 取更紧线、remaining 夹 0、`AllowDispatch` 只在已设 cap 的 100% 关；shape 零发明核＝`config/schema.go:372-373/396-397` **现盘对字存在**、`memory.QuotaExhausted`/`RecordProviderError` 系既有 DAO。接线：`agent/quota_44r4.go` `GateProvider` 在 inner.Stream **之前** `Decide()`——skip＝inner 零调用＋返回 `ClassBudget`（＝票 11 `llm.ChainRunner` 步进形状，其自带 `llm_test.go:458-484` 喂 ClassBudget 元素验证跳下一家，本腿对字核过）＋`EvWarn` 可见事件＋`provider_health.quota_state=exhausted` 落盘（D35 冻结表内既有列，⛔ 加列）。7 枚 AC5 测试今天全绿＋本腿 §5 四发突变各自能红。plan vs pay 物理分桶（三 map）＋「估算」marker 进 stop 文案，M2a/M2b 红证。 |
| **#6** | 价表版本 bump 记录在案、旧聚合不变 | **PASS（v1 记录载体降级原样，不翻案）** | `TestAC6×2` 全绿；三腿未动价表与聚合面；「无逐日持久钉版本列」与 AC#5 同一堵 D35 墙，v1 已判按票面字面 PASS，本腿维持。 |

---

## 4. 44-v1 三处欠账逐条对账（★派单点名）

**欠账① AC#5 整格缺 → 已补（到票面文字射程内）**。
per-model 日/月闸＋plan-credit 分账＋「remaining 0→跳过派发＋链推进＋可见事件」＋80%/100% 事件＋plan/pay 分账——五子项全有实现与测试（§3#5），本腿独立复跑 7/7 绿＋自做 4 发突变验证牙齿。三处**剩余归属**不属本腿可判「未补」也⛔ 藏进绿里：
　(a) **`QuotaGate`/`GateProvider`/`QuotaWiring` 在 cmd/wisp 零调用者**（`git grep` 非测试面只有 CostWiring/SetCostBudget 两枚命中，配额三枚零命中，本腿实量）——派发链在宿主组装，cmd/wisp 写面归他腿、r4 已具名。**这与 v1 欠账③同形状、不同寸**：v1 那寸（价卡＋预算）已接；这寸是**配额闸的装配**，结 44 应由编排者排一发装配腿（接口名 r4 已备齐）。
　(b) **per-model 消费永久落库**仍卡 D35 冻结表（`task_log` 11 列无 model 维，v1 复跑证过、三腿一字未动）——会话级 ledger＋盘上 `quota_state` 信号是冻结 schema 内的诚实解，永久明细行＝**机主批准面**，必须按 SPEC-12 §5 登记推迟或正式拆票，⛔ 静默当已完成。
　(c) realtime 轮同 gate＝票面 `:33-34` 括号点名**归票 60**；球/面板渲染归票 40（事件已上总线、console 先有 reader）。

**欠账② AC#3 warn 无事件出口 → 已补**。发布点在（`costbook_44.go:146-148`）、事件有 reader（run.go `[预警]`）、四枚测试今天绿、摘发布点恰红在事件两枚（M3）。v1 那把尺（`git grep publish(Event)` 成本路径只见 100% 一枚 EvError）**今天重跑不成立**：warn 路径有生产者了。✅

**欠账③ 产品不咬 → 按本派单口径已补**。`cmd/wisp/run.go:1329` `card, budget := agent.CostWiring(cfg, ep.Provider, ep.Model)`＋`:1357` `Price: card,`＋`:1390` `loop.SetCostBudget(budget)`＝**静态核引用成立**（非测试代码零调用者的判据今天只剩余配额闸那一寸，见 ①(a)）；顺序与「读活配置」（reload tick 够得到）由 `Test44R3RunGoAppliesTheWiring` 在真跑 3/3 里钉住；摘掉装配（M4）⇒ 该测试红在三处点名。价卡零卡不猜、坏 `[cost]` 不 clamp 均反向核过。✅（配额闸宿主寸另账。）

---

## 5. 反向判据（本腿自己四发，非转抄 r3/r4）＋还原总核

每发＝`bak\` 备份→改一处→跑点名测试→逐字节还原→`git hash-object` 对基线。

| 发 | 动哪 | 红在谁（具名） | 还原核 |
|---|---|---|---|
| **M1 摘派发前闸**（`quota_44r4.go` `if v.AllowDispatch` → `if true \|\| …`） | **remaining 0 仍派发** ⇒ `TestAC5RemainingZeroSkipsDispatchChainAdvances`＋`TestAC5PlanCreditBalance…`＋`TestAC5ExhaustedSignal…` **3 枚红**（派单元账：「remaining0 仍派发⇒红」量到）；warn/负数/空config 三枚仍绿＝分工不互掩 | hash `b93e862e…`＝基线 ✅ |
| **M2a plan 消费串进 pay 桶**（`quota_44.go` `if m.Billing==BillingPlan && prov!=nil` → `if prov!=nil`） | 分账隔离破 ⇒ `TestAC5PlanCredit…` 红在两处（sibling verdict=plan-credit 90%＋sibling 被闸 calls=0） | hash `27c1753d…` ✅ |
| **M2b 余额算术错**（`remaining = rem` → `remaining = cap`） | **余额错⇒红**：`TestAC5PlanCredit…`（1000≠100）＋`TestAC5NegativeSpend…`（1000≠600）两枚红 | hash `27c1753d…` ✅ |
| **M3 摘 80% 发布点**（publish 行 → `_ =`） | `TestAC3WarnBudgetPublishesWarnEvent`＋`…WarnMode…` 红，`RefusalStaysAnError`／`UnderTheAlertLine…` 绿 | hash `39134bba…` ✅ |
| **M4 run.go 不填价**（`card:=config.Price{}` 断开 CostWiring 调用） | `Test44R3RunGoAppliesTheWiring` 红在缺串＋order(-1)＋live-config(-1) 三点名；另两枚适配测试绿＝红得准。**r3 证据 X2 当时被 24 在飞面挡没拿到干净红、点名要验收腿补——本腿今天 `internal/tools` 已收口（`8523211a0` 在 HEAD），真红补到了** | hash `904683e0…` ✅（＝r3 自录基线同值） |

**还原总核**：四枚 hash 逐一对回基线；`git diff --stat HEAD`（全部 44 面＋cmd/wisp）＝**空**；还原后定向复跑 `ok agent 0.551s`＋cmd 44R3 3/3 绿。`bak\`＋logs 留存 `D:\tmp\wisp\20261102\44v2\`。

---

## 6. 零回归汇总

三包全回归（§2.2）全绿；cmd/wisp 三枚（§2.3）全绿；`go build ./cmd/wisp` rc=0（突变窗口内也编译得动＝形状合法）；gofumpt 净；DEFERRED 零增删（44 面）；thresholds/golden 逐字不动。**新面小账两条（不挡结、具名）**：① `observe/quota_44.go` 的 337 行在 **observe 包内零直测**（`internal/observe` 无 quota 测试文件），其数学经 agent 的 7 枚间接覆盖——M1/M2a/M2b 证明牙齿在，但日后挪调用方时 observe 侧缺本包守门；② `QuotaLedger.roll()` 的换日/换月清桶行为**没有直接测试钉住**（名册里无配额 rollover 用例；日/月滚测试属全局 `cost_daily`，不是配额桶）。建议装配腿或补票顺手两枚。

---

## 7. 是否同意结 44

**同意结 44**——六格按票面文字全部 PASS，且不是把 v1 的三处欠账「判成没欠」：它们是真补了（§4 逐条），补的部分本腿都今天独立复跑＋自己下突变验过牙。

**结票必须带走的三条登记（缺一即本裁决的同意不成立，落点全在编排者/机主手，⛔ 任何写腿能自行拍）**：

1. **配额闸的 cmd/wisp 装配一寸**（§4①(a)）：排一发装配腿把 `QuotaWiring`＋`NewQuotaGate`＋逐链 `GateProvider` 包上宿主派发链。v1 用「机制在、产品不咬」拦下过 `[cost]` 那寸，同一把尺今天量的是配额这寸——不拦格子、拦结票后的寂静：此寸须有票号才许结。
2. **per-model 消费永久落库**＝D35 冻结表，机主批准面（§4①(b)）：按 `SPEC-12 §5` 五字段登记推迟、或正式拆成加列票，二选一，⛔ 就地遗忘。
3. **realtime 轮 gate＝票 60、warn/降级渲染＝票 40**：两处票面自带的归属转登记，别记在 44 的完成里。

**不挡结票但记账的小项**：`cmd/wisp/costwiring_44r3_test.go` 超题面写面字面一寸（内容正当、自报在案）；44-r2 checkpoint 无独立证据件；`probes/orch/20261102/txt/` 的 44r2/r3/r4 派单存档未补（v1 起就缺 r1 的，同一规矩连续两轮没走）；d22scan 今天整树那 1 枚红在未跟踪的 `internal/tools/search_files.go`（票 24 面，归 24 腿收口前改引或建件，⛔ 44 的账）；observe 侧配额直测与 ledger rollover 测试两条小账（§6）；`go test ./cmd/wisp` 进程级 DLL 装载问题＝未跟踪根置 DLL＋临时目录 loader 的环境事实，CI/ legs 若点名跑 cmd/wisp 建议统一用 `go test -c`＋根 PATH 姿势或把三枚 DLL 收进工具链文档。

**一句话**：v1 卡的是「三处没交付」，今天三处都交了、交得能红能还原；44 的算法层、聚合层、预算层、配额层、事件层、装配层（价＋预算寸）至此闭环——按上面三条登记走完，本腿同意编排者翻六枚勾框、结 44。

---

*本表由非实现者裁决腿 44-v2 产；只 commit 本件、⛔ push、⛔ 改码/票面/台账。*
