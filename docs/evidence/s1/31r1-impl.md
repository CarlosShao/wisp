# 票 31 落地腿 31-r1（写腿）＝停手上报件：写面撞车＋三处契约面未批

腿＝`31-r1`（写腿）。票面＝`.scratch/wisp/issues/31-reminders.md`（42 行，`Status: ready-for-agent`，五格全未勾）。
射程＝提醒 `set/cancel/list` ＋ 持久化 ＋ 跳钟鲁棒。⛔ `cmd/wisp/**`。
**结论先给**：本票**今天不能落地**，且不是"难"，是三处硬闸（§2 写面撞车〔现扩到 `internal/memory` ＋ `internal/tools` 两枚包〕、§3 存储契约、§4 名册契约），
逐格 AC 因此全部停在"✗ 未落地（具名阻塞）"，**没有一格被我用部分实现糊成绿**。
原始读数全部落 `.scratch/wisp/probes/31/r1/logs/00-start-readings.txt`（下称"件 00"，[R#] 指其小节）。

---

## §0 起手三核（尺＋读数）

| 核 | 尺（逐字可重跑） | 读数 |
|---|---|---|
| HEAD 现盘最新 | `git log --oneline -5`／`git rev-parse HEAD`／`git branch --show-current`／`git rev-list --left-right --count origin/dev...HEAD` | `23e63fe0a…`、分支 `dev`、首行即本腿派单笔「编排者 A1417…派 31-r1+45-r1 补写腿」；`0  0` ⇒ HEAD＝`origin/dev`、未推 0 [R1][R3] |
| 索引干净 | `git diff --cached --name-only` | 空 [R2] |
| 票面串/行/文件在不在（第 19 处） | 逐枚 `Get-Content`/`Select-String` | 票面在；前置 `04-sqlite-core-done.md`、`30-result-routing-d10-done.md` 均在＝"04／30 已 done"成立；SPEC-07:68、PLAN.md:2568（D34 行）、PLAN.md:2594/2606/2612（16.5.4）、D42 第 9 行、D38（b）名册——逐字命中 [R7] |
| 引文出处对不对（第 20 处） | 档位三向对账＋措辞对账＋行号现跑 | ① `L1 / L0 / L1` 三处逐字一致（票面 L16 ⇄ PLAN.md:2568 ⇄ SPEC-07:68）；② 票面 L21 的局限句「进程未运行则不会触发」**与 PLAN.md:2612 逐字「进程没在跑就不会触发」不同**＝票面是转述，实现按 PLAN 权威串走；③ `AGENTS.md` D 表称 D42＝`PLAN.md:3035` **已腐烂**：现读 L3035 是别的正文，真表头 L3045、真第 9 行 **L3057**（AGENTS.md 自己要求"引用前先现跑那把尺"，本腿照做）[R8] |

另记一件**编排者侧缺件**（不构成本腿判据、不据此改任何格）：派单存档
`.scratch/wisp/probes/orch/20261102/txt/31r1-dispatch.txt` **不在盘**（同目录 14 枚他腿 dispatch 件俱在）。
AGENTS.md ③「派单先存档再发」归编排者动作，本腿只登记、不代写。[R16]

---

## §1 本票射程（从票面读出的包），⛔ cmd/wisp

票面要求三样东西各自有落点：**工具动词**、**一条 SQLite 行**、**进程内定时器（named/owned、SessionScope 无关）**。

| 落点 | 今天状态 | 本腿可否写 |
|---|---|---|
| `internal/tools/reminder*.go`（三枚 Tool ＋ `Decl` L1/L0/L1 ＋ 参数 Schema ＋ `BuiltinReminderEntries()` builder） | 起手尺（19:2x）＝`git status --porcelain -- internal/tools` 空；**19:38:51 复尺＝被占**：多出一枚他腿未跟踪件 `internal/tools/webnet.go`（2696B，ssrfGuard，`//go:build !windows`，与本票零关系）⇒ **也停手** [R18][R19] | ✗ 停手（19:4x 起） |
| `internal/memory/`（`reminder` 表＝schema v5 迁移 ＋ DAO） | **他腿在飞**（见 §2） | ✗ 停手 |
| `cmd/wisp/run.go:571`（`reg := tools.NewRegistry()` 组装根＝注册与拉起定时器） | 面干净但⛔ 禁区 [R18] | ✗ 不动（交完面归下一枚占 `cmd/wisp` 色的腿；先例＝`plan_present.go:138 BuiltinPlanEntries()` 只交货不接线）[R13] |

机器侧今天**零实现**（分母尺）：`git grep -nE "reminder\.(set|list|cancel)" -- "*.go"` ⇒ 命中 **1 行且唯一命中在测试注释**
（`tier_consistency_172_ruler_test.go:129`），滤掉 `*_test.go` ⇒ 非测试 **0** ⇒ 与票 172 普查"未落地抽 3 枚…`reminder.list`"同向 [R17]。

---

## §2 停手闸①：写面与在飞腿同包（＝dispatch 明令的"停手上报"形）

**两枚包都被占，不是一枚**：`internal/memory`（44-r1，起手即撞）＋`internal/tools`（一枚**归属待认**的在飞件，19:38:51 落盘，＝我起手尺之后才长出来的）。

`git status --porcelain` 现读，44-r1（成本计）在飞的字节**正落在本票持久化必须落的包里**：

```
 M internal/agent/cost.go
 M internal/agent/loop.go
?? internal/agent/costbook_44.go
?? internal/memory/costdao_44.go        ← 与本票持久化同包
?? internal/observe/costbudget.go
?? internal/observe/costmeter.go
?? internal/observe/pricing.go
?? internal/observe/pricing.json
```

- mtime 尺（判"正在写"，非"已交未收"）：`costdao_44.go 18:46:46`、`costbudget.go 18:55:42`、`cost.go 19:04:08`、
  `costbook_44.go 19:05:51`；本腿起手＝`19:07:35` ⇒ **起手前 2 分 44 秒她还在落笔** [R5]。
- `go build ./internal/...` rc=0 [R6] ——但**这不等于"色可取"**：本票 AC#1/#2 的判据要在 `internal/memory` 里跑
  `go test`，那一发会把她的**未跟踪文件一起编译进同一包**，红句归属拿不到干净证明。
  台账同族亏已记多次（`A862`／`A866`／`A1042`／`A1012`：同一包两腿＝名册与颜色互相洗）。
- 编排者自己在台账把她那面记窄了（L30315「44-r1…正在写 internal/agent/cost.go+loop.go」），
  **盘上现读多出 `internal/memory/costdao_44.go`** ⇒ 这条更正归编排者，本腿按盘上尺停手。

⇒ 派单原文「若与在飞腿同包或需越界，停手上报」命中，**本腿对 `internal/memory` 零写入**。

★**同一把尺 19:4x 复跑，撞面扩到本票另一枚包**（[R19][R20]）：起手那次（19:2x）`internal/tools` 是干净的，
`19:38:51` 现量多出 `?? internal/tools/webnet.go`（2696B，`//go:build !windows` 的 ssrfGuard，正文逐字
「See webnet_windows.go for platform-specific implementations」而 **`webnet_windows.go` 盘上不存在**）。
该件与票 31 零交集（本票不碰 net/ssrf），但它落在本腿**唯一还可用的那枚包**里 ⇒
**`internal/tools` 也停手**。另 44-r1 一侧自起手后新出四枚 `*_44_test.go`（`cost`／`costbook`／`costdao`／`costmeter`／`pricing`）＝她仍在写。
⇒ 本票射程内三枚包里，**两枚（memory／tools）现被他腿占、一枚（cmd/wisp）是禁区**＝今天没有任何可落地写面；
`git diff --cached --name-only` 复尺仍为空 ⇒ 本腿⛔ add 任何他件，只提自己两枚文档。

---

## §3 停手闸②：持久化唯一合规落点＝存储契约面（D35／SPEC-02 §3）

"只能落 `internal/memory`"不是偏好，是三把尺的量法（[R9]）：

1. `git grep -l "CREATE TABLE" -- "*.go"`（滤 `_test`）⇒ 全仓**唯一**命中＝`internal/memory/schema.go`。
2. `open.go:100-101` 的 `reader`／`writer` 两枚 `*sql.DB` **皆小写字段**；包外无 `Exec/Raw/DB/AddMigration` 之类导出缝
   （导出面只有 `Store` 的各 DAO 方法）。
3. 迁移链组件 `migrationChain:225`／`ddlV1..V4`／`applyV1..V4`／`splitSQLStatements:497` **全 package-private**。

⇒ 要那"一条 SQLite 行"，必须编辑 `schema.go` 并**同时移动四枚既有契约钉**（[R10]）：

| 钉 | 位置 | 今天钉着 |
|---|---|---|
| 版本水位 | `schema.go:159` | `const SchemaVersionTarget = 4` |
| 表名册（排他） | `schema_test.go:201-205` | `contractTableOrder` 十一枚，注释逐字 "Exactly the eleven tables exist (**plus none other**)" |
| 版本读数 | `schema_test.go:257` | `if v != "4"` |
| 逐版语句数 | `schema_test.go:263-279` | v1=15／v2=1／v3=5／v4=5 |

而 `SPEC-02 §3` 标题逐字＝「Schema（**契约级，不得自行增删字段**；【SPEC 细化】给出类型）」，**该节没有 `reminder` 表**。

- 先例（同族怎么走的，[R11]）：`427c62d0a`（票 166-r1 加 v3 `session`/`session_turn`）与 `6e19f4a43`（票 196-r1 加 v4 CHECK）
  都**同一笔里改了这两枚钉**，并把契约文字欠账**具名留在注释**——`schema.go:152-158` 逐字承认
  「⚠ D35 is a frozen contract… SPEC-02 §3 still lists nine tables; this v3 roster extension carries a **spec-note debt** … names for the owner」。
  台账那一步＝L22848「★摆机主 1 件（票 196 契约文字）」。
- ⇒ 本票同形：**Go 侧可落，契约文字那一步＝`SPEC-02 §3` 那张"契约级"表要加 `reminder` 表**（现枚该节列 **9 张表**，逐字见 `schema.go:152-153`「SPEC-02 §3 still lists nine tables」；机器侧名册 `schema_test.go:201-205` 今天 11 张 ⇒ 落地后 12 张），**必须摆机主批**（铁律第 4 条的三类之一）。
- ★**更要紧的"未定义即停"**：票面与 SPEC **都没定义这张表的形状**（表名、列名、有没有正文列、状态字有哪些值、
  missed 怎么记）——票只说「one SQLite row」。按 AGENTS.md §0「方案没覆盖的情况 → 停下来问，**不得自行假设**」，
  本腿**不发明 schema**；§5 只给"待人批的草案"，不作为已定方案落码。

---

## §4 停手闸③：两处名册/分母钉（落地必然牵动，非本腿能自行改）

1. **goroutine 名册（D38（b）＝冻结契约）**。票面 L23-24 要求定时器是「SessionScope-independent task timer goroutine（**named, owned**）」。
   现读名册（PLAN.md:2838-2846 ＋ `observe/goroutine.go:43-60`）＝常驻 6／按需 1／每任务 3／临时 6，**没有任何定时器名**；
   `ResidentBaseline = 6` 亦为钉（D32 各态 goroutine 预算逐字在 PLAN.md:2261-2262，**加第 7 枚常驻＝动 D32/D38**）。
   名册外新名 `ClassifyGoroutine` 答 `CategoryUnknown`（`goroutine.go:81`，注释：名册外＝泄漏症状）。
   - 先例：票 197 生产面就是名册外新名——`observe.Default.Spawn("subagent-finish-"+bg.ID, "tools", watchRoot, …)`（`internal/tools/subagent_197.go`）。
   - ⇒ 两条路都在盘上可指：**甲**＝照 197 形起 `reminder-timer`（⛔ 改契约文字，代价＝名册读数多一枚 Unknown，需编排者裁是否算债）；
     **乙**＝往 D38（b）加一行（＝**改契约文字**，摆机主）。本腿**不代裁**（§6 摆出来）。
2. **票 172 分母钉**：`tier_consistency_172_ruler_test.go:34-44` 逐字「这几个数是分母（AC#3）：**任何一处变了都是契约变更**，
   钉住它们是为了不让分母悄悄缩水」，钉 `tier172TierEntries=41 / tier172DocNames=40 / tier172LandedNumb=9`；
   机器侧装配件 `tier_consistency_172_machine_test.go:15-27` 只注册 FS/Task/Subagent 三族。
   - ⇒ `reminder.*` 落地后，"已落地枚数"这一维**要么 9→12（动分母＝契约变更，须裁）**，要么保持 9（则该尺对 reminder 三枚永不对账＝假绿形状）。
     本腿不自作主张改分母。

---

## §5 逐格 AC（状态／阻塞因／计划用例名／反向判据）

> 统一状态：**✗ 未落地**。下面的"用例名/反向判据"是给下一枚（面 cleared 后）的施工图，**不是本腿的完成声明**。

### AC#1 — set → fires at T+delta with notify（fixture clock）；cancel removes；list shows pending
- 阻塞因：§2（`internal/memory` 撞车）＋§3（表形状未定义）。
- 计划落点：`internal/tools/reminder.go` 三枚 Tool；时钟缝照现成形状注入——`internal/agent/approval/approval.go:20-26` 的 `Clock`、
  `internal/session/scope.go:107` 的 `Clock/Timer(AfterFunc)`、`internal/llm/ratelimit.go:150 NewBucketLimiterWithClock`，**⛔ 新造第四套**。
- 计划用例：`Test31SetFiresAtDeltaOnFixtureClock`／`Test31CancelRemovesPending`／`Test31ListShowsPendingOnly`（真 `*memory.Store` ＋ fixture clock，照 `helpers_test.go:49 openStore` 的形）。
- 反向判据（摘实现⇒点名用例红）：M1 摘掉 `AfterFunc` 的 arm ⇒ 三枚全红；M2 把 `cancel` 写成只改内存不落行 ⇒ `list` 仍列出红。

### AC#2 — Persistence：行活过 Warm dispose；进程重启⇒只有「你错过了 N 条提醒」计数，⛔ 逐条补发
- 阻塞因：§3 是本格主闸（表＋列名＋状态值**票面未定义**）＋§2 同包。
- 计划用例：`Test31RowSurvivesWarmDispose`（`session.Scope.Dispose` 后行仍在、定时器仍在）／`Test31RestartReportsMissedCountOnly`（重启喂两枚过期行 ⇒ 只出一条计数句，通知端口调用次数＝1）。
- 反向判据：M3 摘掉 boot 扫描 ⇒ 计数句消失红；M4 把 boot 改成逐条 fire ⇒ "调用次数＝1"那条断言红（**⛔ 删用例，只加断言**）。

### AC#3 — 排除项：周期性 ⇒ Schema 拒；触发执行 ⇒ Schema／prompt 契约拒（"只做前半句"写在文档行为里）
- 阻塞因：低（本格可在 `internal/tools` 单面包内做完）——但与 §2/§3 同属一票，**不拆半交**（拆＝AC#1/#2 空转，属假绿）。
- 计划形状：`Parameters()` 里 `additionalProperties:false` ＋**根本没有** `recurrence/cron/action/tool` 属性；四道 NOT 各配一条显式拒绝注释，
  指向 `SPEC-12 §5` 的 **REJECTED 行**（逐字「REJECTED | 周期性调度器（cron） | D12 有意取舍；一次性 `reminder` 已移出为窄例外（16.5.4）」）。
  ⚠ **⛔ 用 `DEFERRED(D-xx)` 标记**——AGENTS.md §1.1 要求该标记与 `SPEC-12 §5` 登记表 1:1 双向对得上，这四条是 REJECTED 不是 DEFERRED。
- 计划用例：`Test31SchemaRejectsRecurrence`／`Test31SchemaRejectsActionOnFire`＋拒绝句逐字含"只做前半句"。
- 反向判据：M5 把 `additionalProperties` 摘掉 ⇒ 两枚 Schema 用例红。

### AC#4 — 跳钟：pending 期间系统时间 +5h ⇒ 仍按单调正确触发（D42#9，现读 PLAN.md:3057）
- 阻塞因：§2/§3（同一条行＋同一个缝）；纪律尺本身已具备（`internal/observe/clock.go`＋`tools/d22scan` 两条禁令）。
- 计划用例：`Test31ClockJumpFiveHoursStillFiresOnMonotonic`——缝内 wall `Now()` 前跳 +5h 而单调通道未走 ⇒ fire 时刻仍＝`T+delta`；
  另加 wake/重挂时**重算剩余**一发。
- 反向判据：M6 把剩余时间改成 `wall target − wall now` ⇒ 本用例红（＝把 `time.Since` 单调性摘掉的那一刀，d22scan 也会在同一处点名）。

### AC#5 — set 确认输出里必须有局限句（16.5.4 诚实条）
- 阻塞因：低；串以 **PLAN.md:2612 逐字**「进程没在跑就不会触发」为准（票面 L21 是转述，措辞不同，见 §0 核④②）。
- 计划用例：`Test31SetConfirmationCarriesCaveatVerbatim`（`strings.Contains` 用 PLAN 权威串，⛔ 用票面转述串）。
- 反向判据：M7 摘掉该句 ⇒ 用例红；M8 把句子换成票面转述措辞 ⇒ 同一发红（钉"出处＝PLAN 权威串"这件事本身）。

---

## §6 摆上来的三件（本腿不代裁、不自填）

1. **让面（两枚包，不是一枚）**：① 44-r1 交件入库（`internal/memory/costdao_44*.go` 等收进提交）⇒ `internal/memory` 空出来；
   ② **`internal/tools/webnet.go` 的归属请编排者现认**（19:38:51 落盘、2696B、`//go:build !windows` 的 ssrfGuard，
   正文指向的 `webnet_windows.go` 盘上不存在 ⇒ 看着像一枚未收尾的在飞件；今天 `go build ./internal/tools/` rc=0＝还没人引用它）。
   两枚面都空出来，本票才可落 AC#1/#2。
   撤销口令形：「31 继续，memory＋tools 面已让」。
2. **契约文字（D35／SPEC-02 §3）**：新增 `reminder` 表**及其列定义**属改契约文字 ⇒ **摆机主批**（铁律第 4 条三类之一）。
   待人批草案（**未落码、仅作裁定用**，`text` 列是"消息正文该不该进库"的实义问题）：
   ```sql
   CREATE TABLE reminder (
     id           TEXT PRIMARY KEY,      -- 铸 id（与 session.id 同形，⛔ 派生）
     text         TEXT NOT NULL,         -- 提醒正文（list/重启计数句要用）
     fire_at_wall INTEGER NOT NULL,      -- 落盘只用 wall clock（D42#9）
     created_at   INTEGER NOT NULL,      -- unix 秒
     state        TEXT NOT NULL DEFAULT 'pending'  -- 'pending'|'missed'|'consumed'（值域待裁，是否照 task_log 加 CHECK 亦待裁）
   );
   CREATE INDEX idx_reminder_state ON reminder(state);
   ```
   同批必动：`SchemaVersionTarget 4→5`、`migrationChain` 加 `{from:4,to:5}`、`schema_test.go` 三枚钉（十一枚表名→十二、`"4"`→`"5"`、逐版语句数加 v5 一行）。
   欠账句照 `schema.go:152-158` 先例具名留在注释。撤销口令形：「reminder 表不批，票 31 退回无持久化形」。
3. **定时器名册**：甲（照 197 形，`reminder-timer` 落在名册外＝`CategoryUnknown`）／乙（D38（b）加一行＝契约文字）。
   **本腿推荐甲**：主流 harness 同形（宿主自有 goroutine 登记面按功能增量、⛔ 为省事先把定时器塞进任一 Scope），
   代价**如实写明**＝名册里多一枚 Unknown，需要编排者把它登记成"已知非泄漏"，否则 `RosterReport` 会把症状读数常态化。
   撤销口令形：「31 改裁乙（动 D38 文字）」。

另有两件**归编排者的更正**（本腿只具名，不动他件）：① L30315 对 44-r1 写面的记载少了 `internal/memory/costdao_44.go`；
② `AGENTS.md` §5 D 表 D42 行号已腐烂（真＝`PLAN.md:3057`，非 3035）。

---

## §7 本腿零做清单（对照派单禁区，逐条自证）

- ⛔ 零产码：本程未新建/未修改任何 `.go`；改动只有本件＋件 00 两枚文档。
- ⛔ 未碰 `thresholds.go`／golden／SLO；⛔ 放宽任何断言；⛔ 删用例；⛔ `t.Skip`；⛔ `frontend/**`、`design/**`。
- ⛔ 未碰 `cmd/wisp/**`、`internal/panel/**`、他腿在飞八枚文件（只读）。
- ⛔ `git add -A`／`.`；两处均显式 pathspec 逐枚点名；⛔ amend/reset/rebase/stash/checkout ./clean；⛔ push。
- 临时件只建不删：`D:\tmp\wisp\20261102\31r1\`（本件未用中间产物，目录留空建档）。

**五格终量**：`- [ ]` ×5 → 仍 ×5（本腿未翻任何框；翻框权归非实现者裁腿）。
