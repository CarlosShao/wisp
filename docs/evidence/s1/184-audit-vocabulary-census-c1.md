# 184-c1 — `tool_call.decision` 取证词汇普查（只读·零产码）

- 本程代号 `184-c1`｜派单 `.scratch/wisp/dispatches/2026-09-28-123x-readonly-184-c1-…-what-would-it-cost.md`｜票面 `.scratch/wisp/issues/184-a-decision-column-of-five-frozen-values-…-zero-production-writers.md`
- 起手锚（派单写定）`58124185`；**实际 HEAD `0662a35a`**（= 后继，非派单锚点）——按实际 HEAD 量，**不改判**。日期现量 `2026-09-28 13:0x`。
- 读面：`docs/evidence/s1/184-audit-vocabulary-census-c1.md`（本件）· 读数 `.scratch/wisp/probes/184/c1/{census.sh,readings.log,d22scan.log,gate-clauses.log}`
- **所有 `internal/**` 行号一律取 `git show HEAD:<path>` 版本**（跑时 `git status --porcelain -- internal/ cmd/` 全程为空）；同机 183-v1 在 `internal/risk/**` 变异，本程不碰、不引其工作树版本。
- ⚠ 本件只给**读数与代价**，不拍板。

---

## §1 step-0 四件 + 票面 8 把尺复量

step-0：date=`2026-09-28 13:03:59`｜branch=`dev`｜HEAD=`0662a35a`｜`git status --porcelain -- internal/ cmd/`=**空**（起手与每格前均复验，全程空）。

8 把尺逐把复量（尺命令见 `readings.log`），**8/8 与票面一致，无一需推翻**：

| # | 尺（票面） | 本程 HEAD 复量读数 | 判定 |
|---|---|---|---|
| 1 | `sed -n '26,38p' journal.go` 五枚值 | `journal.go:30-36` const 块 `allow/allow_session_grant/reject/timeout/batch_aggregated`，注释自书 "the vocabulary is already frozen here" | 一致 |
| 2 | `toolCallDecisions` 闭合 map | `models.go:136-142` 五键 map；`:173-174` 写入校验（非空且不在 map→`invalid tool_call.decision`）；`dao_toolcall.go:45-48` `DecideToolCall` 入口无条件再校验（`""` 亦拒） | 一致 |
| 3 | `SPEC-02:79` / `PLAN:2703` | `SPEC-02:79 decision TEXT -- 'allow'|…`（DDL 注释，非 CHECK）；`PLAN:2703` 取证记录行逐字「**它到底对我的机器做了什么**必须可查」 | 一致 |
| 4 | `GrantID:`=0；`DecideToolCall(` 生产仅 `journal.go:101` | `git grep 'GrantID:' HEAD -- internal/ cmd/` **rc=1（0 命中）**；`DecideToolCall(` 生产调用方唯一 `journal.go:101 _ = j.DecideToolCall(ctx,rowID,decision,nil)`（接口声明 `:22`；测试 `dao_test.go:338` 给过 `&grantID`） | 一致 |
| 5 | `privacy.go:96-104` 拼串 | `internal/memory/privacy.go:99-101` `detail := fmt.Sprintf("%s · %s", r.RiskLevel, r.Tool)` → `if r.Decision != "" { detail += " · " + r.Decision }`；`git ls-files` 坐实在 `internal/memory/`（非 `internal/agent/`） | 一致 |
| 6 | bridge 3 处 / loop 2 处 allow | `bridge.go:374/:386/:400`（`dec.DecisionColumn = agent.DecisionAllow`）；`loop.go:794/:803`（`j.decide(ctx, rowID, DecisionAllow)`） | 一致（**另见 §2 AC#1 补记 bridge.go:945 一枚 `decision = agent.DecisionAllow` 局部兜底**，非 ruler 目标形，不新增可达 allow 源） |
| 7 | `riskColumn` 压平 | `journal.go:150-157` `case RiskL1,RiskL2: return risk; default: return memory.RiskL0`；生产调用方 `journal.go:80`（loop.go:611 亦 `riskColumn(info.RiskLevel)`） | 一致 |
| 8 | `PassThroughUnclassifiedRisk` 生产零赋值 | 非测试命中：`loop.go:125`（字段注释）/`:128`（字段声明）/`:774`（注释）/`:799`（**读点·branch 条件**）；`tools.go:23`、`gate.go:133` 注释；**赋值 `true` 仅 `harness_test.go:116`**，测试里 `false`（`declared_l0_risk_179_test.go:87/157`、`loop_approval_test.go:128`）。生产**无赋值点**→恒零值 false | 一致 |

## §2 六格读数（对应 AC#1..AC#6）

### AC#1 五枚记账点逐枚（问没问过人 / 真机可达 / 别的列能不能分开）

调用链（静态坐实，**未跑 go test**，见 §3）：`cmd/wisp → loop.decideRisk`（`loop.go:622`，agent 内建工具）与 `cmd/wisp → bridge.route`（`bridge.go:307`，插件工具，经 `gate.PendingApproval` 见 `approval/gate.go:236`＋`cmd/wisp/run.go:500`）两条都活着。

| 记账点 | ①当时问过人吗 | ②真机可达 | ③共用 `allow` 后别的列能否分开 |
|---|---|---|---|
| `bridge.go:374` route `case risk.L0` | **否**：直接 `return true`，不进 gate、不弹卡 | 是（L0 插件调用） | `risk_level=L0`；`grant_id` NULL；`correlation_id` 与任务同、不区分人；`outcome/error_class/decided_at` 对 allow 行无信号 ⇒ 无列能说"没问过人" |
| `loop.go:794` decideRisk `case RiskL0` | **否**：注释逐字 "run it, interrupt nobody (D4)" | 是（票 179 新支：fs.read/list/task.output 声明 L0 落此） | 同上（`risk_level=L0`） |
| `bridge.go:386` route `case risk.L1` → `AnswerAllow/AnswerTimeout` | **半**：`PendingWindow` 弹了窗口，但**超时未否决＝放行**（`allow` 可来自"沉默"非"点了允许"） | 是（L1 写类工具） | `risk_level=L1` 能把它和 L0 自动、L2 人批分开；但 L1 内部 `allow` 分不清"人点了允许" vs "超时沉默" |
| `bridge.go:400` route `case risk.L2` → `AnswerAllow` | **是**：`PendingApproval` 审批卡，只有真人 `AnswerAllow` 才走到（超时→`:405 DecisionTimeout`，非 allow） | 是（L2 高危·真人批准） | `risk_level=L2` 可标"这是人批的那枚 allow" ⇒ **risk_level 是唯一能把 L2 人批与 L0 自动分开的列** |
| `loop.go:803` decideRisk `default`（未分级）+ 直通开 | 否 | **否（生产到不了）**：`loop.go:799 if !PassThroughUnclassifiedRisk` 恒真⇒先 `:800 DecisionReject` | 不适用（不可达）；**关键**：这一支本可把"未分级"经 `riskColumn("")`→`L0` 伪装成 `L0+allow`，它不可达＝AC#3 成立的地基 |

补记：`bridge.go:943-948` 的 book 兜底 `decision := dec.DecisionColumn; if decision==""{ decision=agent.DecisionAllow; if kind!=Success{=Reject} }`——`route()` **全分支含 `default:414` 都把 `DecisionColumn` 置成非空**（L0→allow/Deny→reject/L1→…/L2→…/default→reject），故正常路由后 `==""` 永不成立，该兜底在路由过的调用上不触发；即便触发也不经 L0 route。⇒ **不构成额外的 `L0+allow` 可达源**。

### AC#2 `decision` **列**的读者名册（不是 `tools.Decision` 结构体）

尺：`git grep -nE '\.Decision' HEAD -- internal/ cmd/` 再逐枚判列/类型（`readings.log`）。同名不同物坐实：`tools.Decision`（`bridge`/`approval`/`cmd/wisp` 路由用的**结构体**，`panel_pump.go:66-71`、`run.go:465/500`、`approval/batch.go:35`）、`risk.Decision`（`panel/pump.go:93`）**都不是这枚列**。

**读"列"并解释 `allow` 语义的，生产只有一处：**
- `internal/memory/privacy.go:100-101` — 导出面把 `r.Decision` **原样**拼进给用户看的 `detail`：用户读到 `L0 · fs.read · allow`，无从知那是机器自批。**"allow"→解释成"允许"，无任何一层说明来源。**

管道/校验（读列但不解释成人话）：
- `dao_toolcall.go:171` `tc.Decision = decision.String`（扫描建行，所有行读者经此）；`models.go:173`（写前校验）。

**面板／`wisp` 命令面读列吗？—— 不读。** 它们读实时路由结构体 `tools.Decision`（`cmd/wisp/panel_pump.go`、`cmd/wisp/run.go`、`internal/panel/pump.go`），非 `tool_call.decision` 列。取数命令即上 `git grep '\.Decision' HEAD -- cmd/ internal/panel/`（命中全为结构体字段/类型，非列）；⚠ 此判据**未把 frontend/\*\* 计入**（别家地界，本程不读），故只覆盖 Go 侧。

钉**枚数**的判据：**无**（见 AC#4）。比较**具体值**的测试（读列·断言"值本身"）：`cmd/wisp/run_test.go:304/317/408/411`、`internal/agent/{declared_l0_risk_179_test.go:103/173/219/238, forensics_test.go:49/160, loop_golden_test.go:344, truncation_test.go:85}`、`internal/memory/dao_test.go:352`、`internal/tools/{bridge_test.go:709, bridge_junction_windows_test.go:621, loop_approval_test.go:207/247}`。

### AC#3 「`L0` ∧ `allow` 的唯一来源」——**坐实（不推翻）**

生产上 `risk_level=L0 且 decision=allow` 只可能来自 **`bridge.go:374`（route-L0 档位放行）** 或 **`loop.go:794`（decideRisk 声明 L0 档位放行）**，两者皆"未询问"。第二条潜在生产路——`loop.go:803`（未分级经 `riskColumn("")`→`L0` 伪装成 `L0+allow`）——**不可达**，因 `PassThroughUnclassifiedRisk` 生产零赋值、恒 false（尺见 §1 第 8 把、`readings.log`）。真人批准永远是 L1/L2（`risk_level` L1/L2），不落 L0。⇒ 编排者那句**成立**，AC#5(a) 的显示层推断可用。

### AC#4 加一枚 `decision` 值的代价清单（只量不改）

**冻结件只引号、只报"要动"、一字节未改**（尺见 `readings.log`）：

1. `internal/memory/models.go:136-142` 闭合 map — +1 键（写前校验源）。
2. `internal/memory/models.go:81` `ToolCall.Decision` 结构体注释（`'' | allow|…`）同步。
3. `internal/agent/journal.go:30-36` — +1 常量。
4. 实际**记账点**（要 emit 新值的那些 `bridge.go`/`loop.go` 行）。
5. `docs/specs/SPEC-02-data-storage.md:79` — DDL 注释枚举串【**冻结件**，`SPEC-02 §3`】。
6. `docs/PLAN.md:2703` — 取证记录行【**冻结件**】。
7. `internal/memory/schema.go:76` — DDL `decision TEXT` 后的注释（**纯注释**）。

**迁移那一格（含幂等雷）——现量结论：加一枚值今天不需要任何迁移。**
- 尺：`git grep -n 'Enum-shaped columns' HEAD -- internal/memory/schema.go` ⇒ `schema.go:14-16` 逐字写明"枚举形列（risk_level/decision/outcome）在 DAO 层校验，**不作 SQL CHECK**"；`schema.go:76 decision TEXT` 无 CHECK。⇒ 列是自由 TEXT，**加一枚枚举值 = 纯 Go map 变更，DDL 不动、无迁移步**。
- 迁移机制现量：**没有"迁移目录"**（尺：`git ls-files | grep -i migration` rc=1）；迁移是**代码内链** `internal/memory/schema.go:168 var migrationChain = { {0→1 applyV1}, {1→2 applyV2} }`，`SchemaVersionTarget = 2`（`schema.go:129`）。
- "会不会让全部历史迁移判定未执行而重跑 / 新老库缺什么"：若**保持 DAO 校验、不加 SQL CHECK**（推荐，配合 AC#5(a)/(c)）⇒ 新老库都本就存任意 TEXT，**零分叉、零重跑**。风险只在**收紧**时出现：若为对齐枚数给 `decision` 加 SQL CHECK，SQLite 不能 `ALTER ADD CHECK` ⇒ 需新增重建表步 `2→3`，且**必须同时改 `applyV1` 基础 DDL**，否则老库经 `2→3` 迁移、新库从 `applyV1` 建，两条路径的枚数不一致＝新老分叉；且 `schema_version` 水位（`open.go:388-407 applyMigrationStep` 用 `INSERT…ON CONFLICT` 幂等）若步号撞号即触发本项目踩过的"版本号当幂等键撞号"雷。
- **枚数判据：无任何测试钉"共五枚"。** 尺：`git grep -nE 'len(.*[Dd]ecision)|== 5|five decision' HEAD -- '*_test.go'` rc=1（`declared_l0_risk_179_test.go:219/238` 的 `len(rows)` 是行数不是枚数）。加一枚值**不打破任何计数断言**；上面 AC#2 那批"断言具体值"的测试只有在你**复用/改写 `allow` 语义**时才需逐一复看。

⇒ 代价总账：**4 处 Go 编辑 + 2 处冻结文档（需人工批准）+ 0 迁移（若保持 DAO 校验）/ +1 重建表迁移（若收紧成 CHECK，且撞号雷风险）**。

### AC#5 不动契约的半条路（逐形可达性＋最坏后果）

- **(a) 导出面文案（`privacy.go:99-101`）——走得通，且不碰契约/枚数。** 因 AC#3 成立，`risk_level==L0 ∧ decision==allow` 可人话化为"自动放行（按档位，未询问）"，**列里的字一字不动**；改的是 `internal/memory/privacy.go`（显示层，非冻结件、非枚举域）。**最坏后果**：只澄清 L0 那一档；`L1 ∧ allow`（超时沉默放行）仍显示 `allow`，`L2 ∧ allow`（真人批）与 `L1 ∧ allow` 在这一形下仍不可分——**半条路只治最常见（只读 L0）那一类误读，非全解**。
- **(b) 真把 `grant_id` 填上——今天 DEAD，卡在 D45 未落。** `approval_grant` 表存在（`schema.go:89`），DAO 有 `InsertGrant`（`dao_misc.go:17`），但**生产无任何调用方**（尺：`git grep 'InsertGrant(' HEAD -- internal/ cmd/ | grep -v dao_misc.go|_test.go` ⇒ **NONE**，`git grep … | grep -v _test` rc=1）。⇒ 没有授权记录生产者就没有可引用的 `grant_id`，此形**须先等 D45 落地**才有意义；单填 `nil→某 id` 是空转。
- **(c) 接受歧义并在 `docs/reports` 明写"取证面不区分自动放行与人工批准"——零改代码零改契约，恒可用，但不解决问题只登记问题。**
- ⚠ 三形皆**未动 `riskColumn`**（票面第 7 把尺已裁其为 `SPEC-02:76 NOT NULL + 域 L0|L1|L2` schema 逼出，非缺陷）。

### AC#6 `_ =` 那一支要不要另立票——**报回：今天不触发，但与"加一枚值"同命**

`journal.go:101 _ = j.DecideToolCall(ctx,rowID,decision,nil)` 的 `decision` 来自 `j.decide` 六枚生产调用方：`loop.go:788/794/800/803/846/882`，**逐枚传五枚已登记常量之一，从不传 `""` 或 map 外值**。⇒ **今天该丢弃的错误不可能触发**（无静默丢行）。但 `dao_toolcall.go:46` 无条件校验（含拒 `""`），故危险窗口恰在"**引入一枚 map 未登记的新值**"时——校验拒→取证行静默不写（又一族"不表现为报错"）。⇒ **建议**：与 AC#4/AC#5(a) 那枚改动**绑定另立票**（改枚举必同时补 map＋处理 `_ =`），**本票不动**；今日无生产触发点。

## §3 本程没测什么（具名）

- **未跑任何 `go test`**（同机 183-v1 在 `internal/risk/**` 变异，遵派单硬约束全免）⇒ 所有可达性（AC#1②/AC#3）为**静态调用链论证**，非运行时证实。
- **未读 `frontend/**`、`design/**`**（别家地界）⇒ AC#2"面板不读列"仅覆盖 **Go 侧**（`cmd/wisp`/`internal/panel`），前端若另读 `tool_call.decision` 不在本程射程。
- **未量 `dec.LevelString()` 对未分级返回 `""` 时 `InsertToolCall` 是否被 `validRiskLevel` 拒**（AC#1 相邻风险，超出本票六格）。
- **未证实 `bridge.go:945` 兜底在所有 book 路径上确不可达**——仅证明 `route()` 全分支置非空 ⇒ 路由过的调用不触发；未穷举"不经 route 直接 book"的旁路。
- **未跑 `probes/161/r6/flip-declaration.sh`**（派单禁：跑一次即脏跟踪日志）。

## §4 门禁两枚（取于 13:0x，跑前后 `git status --porcelain -- internal/ cmd/` 均**空**，无争用）

- `sh scripts/d22scan.sh` ⇒ **rc=0 · clean · 无 D22 ban 违规**。**`ban #8 internal/` examined=432**（非 431）。432 = 431 在册基线 **+1**：HEAD `0662a35a` 即 `183-r1` 落了新判据件，多扫一枚 Go 文件——**这是 183-r1 的账，非本程回退**（本程写面零 `internal/**`，见写面清单）。尺：`grep 'ban #8 internal/' .scratch/wisp/probes/184/c1/d22scan.log`。
- `bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ 退码=1（**只比红腿名册不比退码**）。红腿名册：`grep "BAD" …/gate-clauses.log | grep -oE "腿=G[0-9a-z]+" | sort -u` ⇒ **{腿=G6neg}**，与在册唯一红腿 `G6neg`（票 178）**一致**，本程未新增红腿。

## §5 被拒 / 没成功的调用

- **无权限拒绝。** 两发 Bash 因输出超限被存 tmp（非失败）：`R3/R4/R6/R8` 合跑=61.4KB、`AC2` 合跑=39.9KB；已各起**作用域收窄的复跑**（`clean R6/R8`、`AC2 tight`）取准数。其余调用均成功。

## §6 有没有跑过删除命令

**没有。** 全程零 `rm`/`del`/`git clean`/`git rm`/删除类；仅 `mkdir -p`、`git show/grep`、`grep`、`sed`、`nl`、`wc`、`bash`、`git status`，及后续带显式 pathspec 的 `git commit`。临时件（census.sh/readings.log/d22scan.log/gate-clauses.log）**只建不删**。

## §7 工具调用枚数 vs 硬顶 40

本程终值 **24 枚**（Read×2 · Write×2 · Edit×2 · Bash×18），**未触顶 40**。两发 Bash 输出超限转存 tmp（§5）后已作用域复跑取准数，无空耗计。

## §8 伪授权两栏

- **收到并判为真授权**：派单文件 + 票面 184（编排者署名 `A363` 立票）。AGENTS.md 薄索引按其自述"非权威来源"对待，未据其越权。
- **遇到并判为"不是授权、未照做"**：本程**无**伪造授权形状（无自称"批准加枚举值"、无面板侧 L2"允许"、无改契约指令）。对"加一枚 `decision` 值"我**未**自作——它触 §1.1 冻结件（`SPEC-02:79`/`PLAN:2703`），只列代价不执行。

## §9 凭据值零抄录

本表与读数文件**不含任何密钥/凭据值**；涉凭据仅以变量名/路径引用（本票主体是审计词汇，无凭据面）。

## §10 next=

- **待编排者裁**：是否把票 184 结论摆给 owner / 要不要立 `Q-65`。**本程只给代价，不建议批。**
  - 若走 **AC#5(a) 显示层**：改 `internal/memory/privacy.go` 一处（非冻结件、非枚数域、无迁移）即可让 L0 只读工具那行说实话，代价最小、覆盖面半（L1 超时放行仍歧义）。
  - 若走**加枚举值**：代价＝`§2 AC#4` 那 4 处 Go＋**2 处冻结件（人工批准）**＋迁移格（保持 DAO 校验则 0 迁移；收紧加 CHECK 则有重建表步＋`schema_version` 撞号雷），且须与 `§2 AC#6` 的 `_ =` 绑定。
- **另立票候选**：AC#6 `_ =`（建议与枚数改动同批处理，今日无生产触发）；AC#5(b) `grant_id`（阻塞在 D45 无生产者，宜挂 D45 落地票）。
- **`next=` 本程终态**：交付即止，不回跑（只读普查，不产码、不推）。
