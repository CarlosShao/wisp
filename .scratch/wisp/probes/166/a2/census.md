# 票 166 · a2 只读普查（动一版 schema 的代价表：A 面积 / B 第二条路 / C 归档耦合 / D 取证链 / E 人话后果）

- 子代理：`166-a2`（只读普查，未跑任何 Go 命令，未产码一字）
- 锚点（`date -Iseconds` ＋ `git log -1 --format=%H` **自取**）：`2026-10-03T14:48:39+08:00` 同刻 = `efc130614745510200b07ea53b9ab7832985cb6d`
  - ⚠ 与上一枚腿 `166-a1` 的锚（`70a935ce` / `1a14f203`）之间已越过若干枚 commit ⇒ **本件行号一律取自当刻工作树现读**，不复用 a1 的行号；两处与 a1 不同处以 F 节为准。
- 树状态：**不干净**（`git status --short` 有 4 枚写腿在飞的条目，含 `cmd/wisp/config_readers_255.go`、`cmd/wisp/panel_host_windows.go`、`cmd/wisp/panel_resident_windows.go`）。
- 射程：显式搜索根 `cmd internal tools scripts`（＋ `docs` 只读、`.scratch/wisp/{issues,probes}` 对账）。`frontend/**`、`design/**` **未读、未引**。
- 本程**未跑任何 Go 命令**（无 `go test`/`build`/`vet`/`run`）。必须跑了才知道的格子全部进 G 节并写复量法，一处未猜。

---

## 0. 一句话结论（只量代价，不拍板）

**"加一枚 `session` 表 ＋ `task_log` 加一枚可空列"这一最小形，在现仓规下要动 10 枚生产文件（含 2 枚新建）、17 枚会红的断言点（分布在 4 枚测试文件里，其中 3 枚是 a1 没列出的连环钉）、6 节 spec 文字；而且"完全不动 schema"的那条路也躲不开契约——`docs/PLAN.md:1363` 把 C13 `MemoryStore` 的射程逐字写成"L1 画像／L2 显式／L3 任务日志"，会话存储不在里面。**
第二条路里**只有 ⓐ 的一半可行**（用 `task_log.id` 前缀分组 ⇒ 能列、不红钉），但它给不出"标题／归档位／分叉父指针"任何一个可写的地方；ⓑ（塞 JSON 列）**不可行**（`task_log` 一枚 JSON 列都没有，唯一两枚 JSON 宿主一枚被整体替换、一枚 owner 不对）；ⓒ（只加不加表的可空列）**不可行**——那枚钉**逐字断的是"每枚表的列集合"，不只是表名集合**（`schema_test.go:259-260`）。
归档那一格：**今天连"该启动它的那枚容器"都不存在**——`StartRetentionJob` 需要 `*plugin.DisposalScope`，而 `plugin.NewDisposalScope` 在四枚根里的生产调用者＝**0 枚**。

---

## A. 动一版 schema 到底要动哪几处（最小形：`session` 表 1 枚 ＋ `task_log` 加 1 枚可空列）

### A-1 生产面（逐枚 `file:line` ＋ 它今天写着/断着什么）

| # | 位置 | 今天是什么（引号内为本程 grep 到的逐字） | 为什么必须动 | 不动的后果 |
|---|---|---|---|---|
| A1 | `internal/memory/schema.go:9-17` | 顶部注释逐字："The DDL below is CONTRACT-LEVEL: it must stay byte-identical to the ```sql blocks of docs/specs/SPEC-02-data-storage.md §3 (including comments). No added or removed columns, tables, indexes or constraints" | 这段自陈与新表/新列**正面矛盾**，不重写就是留一条说谎的注释 | 注释继续声明一条已被本 commit 违反的规矩（见 F3：它今天已被自己违反 2 处） |
| A2 | `internal/memory/schema.go:122-129` | `const SchemaVersionTarget = 2`，上面 v1/v2 历史注释逐字列了两步 | 目标版本要变 3（迁移与 `open.go:121 target: SchemaVersionTarget` 都由它驱动） | 链走不到新步 |
| A3 | `internal/memory/schema.go:134-146`（`ddlV2` 之后） | `ddlV2` 是一段裸 SQL 常量，注释逐字"mirrors the ```sql block in docs/specs/SPEC-02-data-storage.md §3 byte-for-byte" | 新表要有一枚同形的 `ddlV3` | 无处放 DDL |
| A4 | `internal/memory/schema.go:168-171` | `var migrationChain = []migrationStep{{from: 0, to: 1, apply: applyV1}, {from: 1, to: 2, apply: applyV2}}` | 追加 `{from: 2, to: 3, apply: applyV3}` | 老库永不变形 |
| A5 | `internal/memory/schema.go:184-191`（`applyV2` 之后） | `applyV2` 的形：`for _, stmt := range splitSQLStatements(ddlV2)` | 要一枚 `applyV3` | 步没有 apply |
| A6 | `internal/memory/schema.go:193-216` | `splitSQLStatements` 注释逐字："no string literal or comment in it contains a semicolon (asserted by TestDDLContractIntrospection...)" | ⚠ **不是要改它，是要受它约束**：新 `ddlV3` 的注释里出现 `;` 会被**静默错切成两句** | 半条 DDL 落库，且没有测试会告诉你（F5：`TestDDLContractIntrospection` 这枚名字在本仓不存在） |
| A7 | `internal/memory/models.go:10-13` | 头注释逐字："Row models for the eight D35 tables (SPEC-02 §3)" | 枚数变了，注释要跟着 | 注释说谎 |
| A8 | `internal/memory/models.go:58-70` | `TaskLog` 结构体**恰 11 个字段**（与 DDL/`contractSpec` 三处双向对得上） | 加 `SessionID` 一枚字段 | 读不到新列 |
| A9 | `internal/memory/models.go:185-197` | `validateTaskLog` 今天只断三件事：`ID` 非空、`State` 非空、`error_class` 走 `observe.ValidateErrorClass` | 新列的校验要放这里（DDL 无 CHECK，见 C1） | 空/非法 session 键静默入库 |
| A10 | `internal/memory/dao_tasklog.go:33-37` | `INSERT INTO task_log(...)` 显式 11 列 ＋ `VALUES(?,?,NULL,?,?,?,?,?,?,?,?)` **恰 11 个占位** | 要写 `session_id` 必须同批改列名表与占位表 | 新列永远 NULL（**不报错**，四枚根里 `SELECT *`＝0 命中，见 D3） |
| A11 | `internal/memory/dao_tasklog.go:80-82`、`:95-97` | 两枚显式 SELECT 列表（各 11 列，`COALESCE(error_class,'')` 收尾） | 要读就得加 | 读路径看不见新列 |
| A12 | `internal/memory/dao_tasklog.go:144-160` | `scanTaskLog` 按**位置** `scan(&tl.ID, &tl.StartedAt, ...)` 11 个目标 | 与 A10/A11 同批 | 错一位＝**静默串列**（不 red，只是数据错着放） |
| A13 | 新文件 `internal/memory/dao_session.go` | 不存在（`ls internal/memory` ⇒ 无 `dao_session*.go`） | 会话行模型 ＋ list/search/rename/archive 原语；可复用现成两枚助手：`models.go:215-219 escapeLike`、`dao_memory.go:52 SearchMemory`（LIKE＋`ESCAPE` 形状）与 `dao_memory.go:69` 的 `LIMIT` 形 | 票面四件事没有一件有存储跳 |
| A14 | `internal/memory/privacy.go:17-33` | `PrivacyDomain` **恰五枚**常量 + `PrivacyDomains()` 逐字"returns the domains in SPEC-02 §5 order" | 会话含用户文本 ⇒ 要么进名册（改四枚 switch：`:47-131`、`:135-165`、`:170-186`、`:198-239`，另 `:262-269 purgeAllDomains`），要么在 spec 里明写它为什么不进 | 隐私页**漏域**（`SPEC-02:167` 逐字要"对 … 全部提供：一键清空 / 导出（JSON）/ 逐条删除"） |
| A15 | `internal/memory/retention.go:30-42`、`:85-93`、`:139-190` | TTL 常量组（`TaskLogTTL = 30 * 24 * time.Hour`）、`RetentionResult` 五枚计数、`runRetention` 四步删行 | 会话表要不要窗口必须在这里有归属（`SPEC-02:152` 节标题逐字"保留期与清理（有归属才有效）"） | 新表成为一张**永不被清**的表；反之若纳入，见 C 节"归档挡不住删除"那一格 |

**面积小结**：生产文件 **10 枚**（`schema.go`、`models.go`、`dao_tasklog.go`、`privacy.go`、`retention.go`、`open.go`（A-3 的注释锚）、`doc.go` ＋ 新建 `dao_session.go`、`dao_session_test.go`，另加隐私/保留两族的同批改动）；其中"改的是注释而非代码"的有 3 处（A1/A7/A-钉里的 `:29`、`:44`）。

### A-2 会红的断言（17 枚，逐枚给坐标 ＋ 它今天断什么）

| # | 位置 | 今天逐字断什么 | 为什么红 | 改法（只描述，不代拍） |
|---|---|---|---|---|
| P1 | `schema_test.go:52`（`contractSpec` map，条目 `:53-165`） | 注释逐字"contractSpec mirrors SPEC-02 §3 exactly: nine tables … Any drift between schema.go and the spec fails here" | 新表要加条目；`task_log` 条目（`:81-96`，恰 11 枚 `colSpec`）要加一枚 | 加条目（＝**改断言**，须非实现者复判） |
| P2 | `schema_test.go:168-171` | `contractTableOrder` 恰九枚名 | 变十枚 | 同上 |
| P3 | `schema_test.go:195-196` | `if len(gotTables) != len(contractTableOrder) { t.Errorf("table count = %d, want %d: %v", …) }` | 10 != 9 | 与 P2 同批 |
| P4 | `schema_test.go:198-208` | 逐名核 `missing table %q` ⇒ 表**名册**双向都断 | 只加表不改名册则红 | 同上 |
| P5 | `schema_test.go:259-260` | `checkColumns`：`if len(got) != len(spec.cols) { t.Fatalf("column count = %d, want %d\n…") }` | **只加一枚可空列、不加任何表，也在这里红**（Bⓒ 的关键） | 加列规格；⛔ 不许把它改窄 |
| P6 | `schema_test.go:262-276` | 逐位核 `col[%d] name` / `type` / `notnull` / `pk` | 列名、**顺序**、可空性都是契约（新列必须追加在 `PRAGMA table_info` 的末尾一位） | 同上 |
| P7 | `schema_test.go:222-224` | `if v != "2" { t.Errorf("schema_version = %q, want %q", v, "2") }` | 版本值钉（幂等键那枚行的值） | "2"→"3" |
| P8 | `schema_test.go:29`、`:44` | 注释逐字"(D35, eight tables)" / "nine tables" | 枚数过期 | 改注释 |
| P9 | `schema_test.go:229-231` | `if n := len(splitSQLStatements(ddlV1)); n != 15 { … want 15 (8 tables + 7 indexes) }` | **只有当你把新东西写进 `ddlV1` 才红**；走独立 `ddlV3` 这一条不红（分野见 Bⓒ） | 优先走独立 v3 段 |
| P10 | `schema_test.go:233-235` | `len(splitSQLStatements(ddlV2)) != 1` | 同上（ddlV2 不该被改） | 不动 `ddlV2` |
| P11 | `schema_test.go:341-343` | `TestMigrateFreshDatabaseCreatesCurrentVersion` 里 `if v != "2"` | 版本值 | "2"→"3" |
| P12 | `schema_test.go:350-358` | `wantBackups := []string{"wisp.db.bak-0-1", "wisp.db.bak-1-2"}` ＋ `:351` `len(entries) != len(wantBackups)` | 新步多产生一枚 `wisp.db.bak-2-3`（`open.go:493-495 backupPath` 逐字 `%s.bak-%d-%d`） ⇒ 按名册数就红 | 名册加一项 |
| P13 | `schema_test.go:369-380` | `withTestMigration` 逐字 `migrationChain = append(saved, migrationStep{from: 2, to: 3, apply: next})` | **连环钉（a1 未列）**：生产链一旦有 2→3，`nextStep`（`open.go:379-386`）返回**第一个** `st.from == from` ⇒ 测试那步永不执行 | 测试脚手架锚位改 `from:3, to:4` |
| P14 | `schema_test.go:406-411`、`:419`、`:433`、`:446-452` | `TestMigrationChainWithBackup`：`withSchemaTarget(3)`、断 `backup/wisp.db.bak-2-3`、断 `v != 3`、断 `schema_probe_v2` 存在（失败消息逐字"test-only v2 table missing after migration"） | 受 P13 牵连：probe 表建不出来 ⇒ 这一枚**由绿转红** | 同批改锚（4 个点） |
| P15 | `schema_test.go:479-482`、`:494-496` | `TestMigrationFailedStepIsAtomic`：`t.Fatal("Open succeeded despite failing migration step")`、`want 2 (atomic rollback at v2, target v3)` | 受 P13 牵连：生产 2→3 不会失败 ⇒ Open 成功 ⇒ 红。⚠ **这一枚正是票面 AC#3 要的"崩溃原子性"现成形状**——它一红，AC#3 的证据必须先修它 | 同批改锚 |
| P16 | `migrate_v1seed_test.go:68-70` ＋ `:106-118` | `if v != "2" { … want 2 }`；`wantBackups := []string{"wisp.db.bak-0-1", "wisp.db.bak-1-2"}` ＋ 枚数比对 | 版本值 ＋ 备份名册（与 P11/P12 同物，**另一枚文件**） | 两处同批 |
| P17 | `privacy_test.go:53-60` | `for _, d := range PrivacyDomains()` ＋ `:58` `if len(items) == 0 { t.Fatalf("list %s: no items after seeding", d) }` | **一旦会话进隐私名册（A14），这枚循环要求 `:20-49` 那段 seed 真的产出一枚会话行**，否则**整枚测试 Fatal 中止**（不是 Error，后面的 delete/export/purge 三段都不跑） | seed 段加一发 ＋ 断言加一域 |

**会不会红的邻居（本程读到、判为不红，理由具名）**：
- `schema_test.go:510-562 TestMigrationNewerSchemaIsUnmigratable`：`:525` 伪造 `value='99'`，99 > 任何 target ⇒ 目标升到 3 后仍成立。
- `schema_test.go:564-586`、`:591-608`：与版本枚数无关。
- `internal/memory/dao_test.go`、`retention_test.go`、`concurrent_test.go`：INSERT 全用显式列名（`privacy_test.go:31,34` 那种 `StartTaskLog(TaskLog{…})` 形＋`dao_tasklog.go:33` 的显式列表），加一枚**可空**列不改变既有语句的合法性。
- `internal/config/boundary_test.go:56-70 TestBoundaryNoRuntimeObservationFields`：禁的是 config 结构体字段名子串（`health|probe|latency|…`），"session" 不在禁列 ⇒ 不红；但若落地把会话标题塞进 `config.toml`，要重读这枚尺的 owner 判断（它断的是"运行观测态不进 config"）。

### A-3 SPEC 侧要动的节（票面 AC#5 冻的正是这些文件）

| 位置 | 今天逐字写着 | 要动到什么 |
|---|---|---|
| `docs/specs/SPEC-02-data-storage.md:22` | "## 3. Schema（**契约级，不得自行增删字段**；【SPEC 细化】给出类型）" | "不得**自行**"留了人工批准的门，但这一节本身是 spec ⇒ 动它＝人工批准 |
| 同文件 `:24-124` | 一枚 ```sql 块（ddlV1 的镜像，含 8 CREATE TABLE + 7 CREATE INDEX） | 若不动 ddlV1 文字，这一节要**追加**一段 v3 说明而非改这段 |
| 同文件 `:126-131` | 逐字"schema v2 迁移 DDL（与 `internal/memory/schema.go` 的 `ddlV2` **逐字节一致**；票 09 已落地）" ＋ `:133-147` 第二枚 sql 块 | 新表照这个**先例形**加第三段（这是仓里唯一一枚"spec 追加而不动老块"的成例，v2 就是这么干的） |
| 同文件 `:152-163`（§4） | 保留期表逐枚列五行 ＋ 节标题"有归属才有效" | 会话表要么进这张表，要么明写为什么不进（否则 C 节那一格又落空一次） |
| 同文件 `:165-168`（§5） | 逐字点名 `profile`/`memory`/`task_log`/`tool_call`/`artifacts` 五域 | 隐私页域名单（与 A14/P17 同一枚事的两面） |
| 同文件 `:170-185`（§6） | 目录布局图只有 `wisp.db (+ -wal/-shm)` | **只有走"第二个 db 文件"那条才要动**（见 Bⓔ） |
| 同文件 `:203-209`（§8） | 逐字"**每张表一组契约测试**：字段/索引存在性（用 `sqlite_master` 自省）＋ 保留期清理边界" | 新表**必须**带来一枚新契约测试（这是**加**断言，不是放宽；不遵守＝违 §8） |
| `docs/PLAN.md:2692-2719`（D35） | `:2695` 逐字"**八张表**，字段/索引/保留期如下（契约级，不得自行增删字段）"；`:2699-2706` 八行表；`:2710-2714` 规则 1（保留期归属）；`:2715-2716` 规则 2（隐私页域） | 枚数一句＋表一行＋两条规则各一读 |
| `docs/PLAN.md:1363`（C13 那一行） | 逐字：`\| C13 \| MemoryStore \| L1 画像（≤20 条）/ L2 显式 / L3 任务日志 \| D20 \|` | **会话不在 C13 射程** ⇒ 只要把会话 DAO 放进 `internal/memory`，就是**扩 C13＝C1–C32 契约面＝人工批准**（`AGENTS.md` §1.1）。⚠ 这一条**连"不动 schema"那一形都躲不开**，是 A 节里最早撞上的一枚 |

### A-4 计数（要给 owner 的那一个数）

- 生产文件：**10 枚**（含 2 枚新建）；
- 会红的断言点：**17 枚**，分布在 `internal/memory/schema_test.go`（13 枚）、`internal/memory/migrate_v1seed_test.go`（2 枚）、`internal/memory/privacy_test.go`（1 枚，且是 `t.Fatalf` 级）、以及 `schema_test.go:29/:44` 注释（1 枚）；
- 其中 **3 枚（P13/P14/P15）是"测试脚手架自身的锚位"**，a1 未列出——它们不是"加一列所以数不对"，而是**"借用 2→3 当测试专用步"的形状与生产 2→3 正面撞车**；
- spec 侧：**6 节 SPEC-02 ＋ 2 处 PLAN.md**（其中 `PLAN.md:1363` 是 C 级契约面，不是 D 级决策文字）。

---

## B. 有没有"不动契约那一版"的第二条路（三形各给结论 ＋ 两形登记）

### ⓐ 不加表，用 `task_log` 上的一枚已有列／既有 id 前缀分组 —— **一半可行**

先指认列名（三处双向核过：DDL `schema.go:53-65`、`contractSpec` `schema_test.go:82-94`、结构体 `models.go:58-70`，**恰 11 枚**）：
`id` `started_at` `ended_at` `state` `query_text` `summary_text` `cost_tokens_in` `cost_tokens_out` `cost_amount_micro` `currency` `error_class`。
⇒ **没有一枚是会话键、分组键或标题**。剩下的可能只有两条腿：

1. **腿一：把会话编码进 `task_log.id` 的字符串形状。** 不红 P1–P12（列数/表数/版本都不动）。代价逐枚：
   - 铸点 `internal/agent/loop.go:1123 newTaskID()`（注释 `:1121` 逐字"the task_log.id column is"）；
   - 等值匹配的两处读路径：`dao_toolcall.go:111-112`（`WHERE task_id=?`）、`dao_tasklog.go:82`（`WHERE id=?`）；
   - `cmd/wisp/approval_reply_201_test.go:681` 逐字在用 `strings.TrimSuffix(corr, "-corr")` 反推 task id ⇒ 那枚用例的 id 推导链跟着变（本程只登记形状，会不会真红要跑，见 G3）。
   - ⛔ **致命局限（盘上有据）**：`internal/session/session.go:19-20` 那枚 `session.ID` 是**故意**不跨重启的（票面 a1 已引，本程复量：`session.Mint()` 生产铸点只有 `cmd/wisp/run.go:473` 一枚）⇒ id 前缀只能把"同一发里派生出的任务"归堆，**用户明天新开一发时没有任何东西告诉他"这属于昨天那一串"**。所以腿一给得出"列"，给不出"跨重启回到昨天的对话"。
2. **腿二：拿 `state` 当归档位。** `validateTaskLog`（`models.go:185-197`）只断非空；DDL 里 **CHECK 零枚**（尺：`grep -rn "CHECK(" internal/memory --include=*.go | wc -l` ⇒ **0**）⇒ 加第五个词**门不会响**。⚠ 这正是票 196（状态词表，未派、待 owner 拍）警告的形状：`state` 今天已有多套词不相认，再塞一词是把功能建在未定案的地基上。

⇒ **ⓐ 结论**：**只能给出"分组／列"这半件事。** 改名（标题）、归档位、分叉父指针在 11 枚列里**没有合法宿主**：
- 写 `query_text` ＝覆盖已脱敏的用户原话（DDL 注释 `schema.go:58` 逐字"已脱敏；模式匹配掩码是尽力而为，须标注（§14.4）"），违"真相源只追加"（票面硬约束 3）；
- 写 `summary_text` ＝`dao_tasklog.go:58-59` 逐字 `summary_text=COALESCE(?, summary_text)` 且生产写点在 `cmd/wisp/run.go:1135` 每次任务收尾都带摘要 ⇒ 用户改的名**下一次收尾就被覆盖**（a1 的 (a) 支，本程按 SQL 原文复量成立）。

### ⓑ 塞进已有 JSON/state 列 —— **不可行**

全库 JSON 形状列逐枚指认（谁在写）：
| 列 | 表 | 写者 | 能不能当宿主 | 逐字依据 |
|---|---|---|---|---|
| `args_json` | `tool_call` | `dao_toolcall.go:24`（NOT NULL，逐次调用一行） | **不能**：owner 是"一次工具调用"，且它被取证读路径整列拉走（`:159 toolCallSelect`）——往里塞会话＝污染"它到底对我机器做了什么"那一族 | `schema.go:74` 注释逐字"长字符串截断（§5.1 脱敏规则）" |
| `probe_json` | `provider_health` | `dao_providerhealth.go:116-122`，逐字 `ON CONFLICT(provider, model) DO UPDATE SET probe_json = excluded.probe_json` | **不能**：**每次探测整体替换**，塞进去的字段活不过下一次 `UpsertProviderProbe`；owner 是 provider/model | `internal/llm/probe_health.go:114` 注释逐字"probe_json is REPLACED wholesale by the DAO"；解码侧 `dao_providerhealth.go:207` 是普通 `json.Unmarshal`（**不拒未知键**，所以塞进去时也不会报错——静默丢失比报错更坏） |
| `capabilities_json` / `net_allowlist_json` | `plugin_state` | `dao_misc.go:204-209` | **不能**：键在插件 id 上，且 `:193` 逐字断 `plugin_state.capabilities_json is required` | `schema.go:115-116` |
| （会话目标表 `task_log`） | — | — | **一枚 JSON 列都没有**（11 枚名册见 ⓐ） | 三处双向核过 |

尺寸那一问（票面点名要答的）：
- **今天没有任何尺寸检查可撑爆**：四枚根里对 `query_text`/`summary_text` 的长度上限＝0 枚；唯一 `truncate` 助手在 `dao_providerhealth.go:143,223`（只管 `last_error`），`privacy.go:251 firstLine` 只管列表标签截一行。
- 真实代价在**读路径形状**：`ListTaskLogs`（`dao_tasklog.go:94-97`）**无 `LIMIT`、无分页**，整表 `append` 进 `[]TaskLog`（尺：`grep -n "LIMIT" internal/memory/dao_tasklog.go` ⇒ **0 命中**；对照 `dao_memory.go:69` 是带 `LIMIT ?` 的）。⇒ 把会话正文塞进 `task_log` 任何既有列，等于让"列历史"这一发随用量线性变胖变慢，而它今天已经在被 `privacy.go:76`、`:214` 两处全量走。

### ⓒ 只加不加表的可空列 —— **不可行**；那枚钉断的是"每枚表的列集合"，不只是表名集合

逐字读那枚钉（`TestSchemaContractIntrospection`，`schema_test.go:173`）——**两件事都断**：
- 表名册：`:179` 那条 `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'` ＋ `:195` 枚数 ＋ `:198-208` 逐名"missing table %q"；
- **每枚表的列集合**：`:211-214` 对 `contractSpec` 里**每一枚表**跑 `checkColumns` ⇒ `:259-260` `if len(got) != len(spec.cols) { t.Fatalf("column count = %d, want %d…") }` ＋ `:262-276` 逐位核 `name/dataType/notNull/pk`；
- 另断索引枚数（`:298-300` `%s index count = %d … want %d`）与版本值（`:222`）。
⇒ **"只加一列"照样在 P5 红**，且列的**位置**也被 `:262` 断（`ALTER TABLE … ADD COLUMN` 落在末尾，`PRAGMA table_info` 也是末尾，这一形是匹配的；把列插中间则红 name 位序）。
⇒ **一条真实的"少红一枚"分野（不是绕过，量清楚给实现者）**：新列/新表若写进**独立的 `ddlV3`**（`ALTER TABLE task_log ADD COLUMN` ＋ `CREATE TABLE session`）而**不改 `ddlV1` 的文字**，则 P9/P10（`ddlV1`=15 句、`ddlV2`=1 句）**不红**；若图省事把新列直接加进 `ddlV1` 的 `CREATE TABLE task_log`，则 P9 由 15 变 16 红，**且新库靠 CREATE、老库靠 ALTER 两条形状分叉**——票面 AC#3 后半句"分别推演新库/老库各缺哪列"要防的正是这一形。
⛔ **禁区登记**：`contractSpec` 的枚数、`:222` 的版本值、`:195`/`:259` 的比较，任何"改窄"都是**放宽判据**（`AGENTS.md` §1.2 同族规矩）。这三枚只许**加**、不许**减**。本程一律不动它们，只登记"这 17 枚的改动必须人工批准 ＋ 非实现者复判"。

### ⓓⓔ 两种"不红任何上述钉"的形 —— **只登记，不推荐，且都仍要人工批准**

| 形 | 为什么上面 17 枚全不响 | 为什么本程判它**比加表更坏**（逐枚有据） |
|---|---|---|
| ⓓ 把会话当额外 key 的**行**塞进 `schema_meta`（`key TEXT PRIMARY KEY, value TEXT NOT NULL`，列一枚不加） | `readSchemaVersion`（`open.go:422-470`）只读 `key='schema_version'`；`applyMigrationStep`（`open.go:402-408`）只 UPSERT 那一枚键；**没有任何测试数 `schema_meta` 的行**（尺：`grep -rn "schema_meta" cmd internal tools scripts --include=*.go \| wc -l` ⇒ **20 行**，逐行要么锁 `key='schema_version'`、要么只核表存在性） | 表的**用途**是契约文字：`schema.go:25` DDL 内注释逐字"-- schema 版本与迁移水位（§14.5 的地基）"、`PLAN.md:2699` 该行逐字"schema 版本与迁移水位"。往里放用户原文＝**把钉当空子钻**（改的是语义不是枚数），而且它**天然绕开隐私名册**（`privacy.go:20-33` 五域里没有它，P17 那族尺根本不数它）⇒ 用户"一键清空"清不掉它，正面撞 `SPEC-02:167` 那句"全部提供"。**判：违 §5/§3 语义，须人工批准；本程判它不可取。** |
| ⓔ 第二个 SQLite 文件（如 `data\sessions.db`） | `TestSchemaContractIntrospection` 只自省 `s.reader`（`schema_test.go:174-175`），而 `s.reader` 就是 `wisp.db`（`open.go:199` ＋ `:39 dbFileName = "wisp.db"`）⇒ 第二张库**不被那枚钉看见** | 逐枚代价：(1) `backupDatabase`（`open.go:501-525`）只拷 `s.dbPath` 与 `-wal/-shm` ⇒ 新库**没有迁移前备份**，票面硬约束 1"崩了不能留半个库"在这一形里**没有任何防线**；(2) `startupCheckpoint`（`open.go:228-243`）只 checkpoint `wisp.db`；(3) `privacy.go` 五域与 `ExportPrivacy:198` 都不认它 ⇒ 隐私页漏域（同 ⓓ）；(4) `PLAN.md:2692` 标题逐字"D35 — 数据模型（SQLite，`%APPDATA%\wisp\wisp.db`）"＋`SPEC-02:172-185` 布局图只有 `wisp.db` ⇒ 要动的是 **D35 本身**＋§6，动的节比加表还多。⚠ 负向句自量：`internal/winsec/acl_windows_test.go:222` 与 `internal/winsec/production_windows_test.go:42` 都只按 `{"wisp.db","wisp.db-wal","wisp.db-shm"}` **三枚名**巡 ACL，**不是**"目录里只许这三枚"的穷举断言 ⇒ 这两枚尺**看不见**第二枚 db（本程自己跑过 `grep -rn "wisp.db" cmd internal tools scripts docker` 才敢这么写）。 |

### B 的一句话结论

- **可行**：ⓐ 的腿一（`task_log.id` 前缀分组）——给"列"，不给"回到昨天"；
- **不可行**：ⓐ 的改名/归档/分叉（11 枚列里没有合法宿主，两条可写路径都有具名的覆盖/污染后果）、ⓑ（`task_log` 零 JSON 列；两枚宿主一枚被整体替换一枚 owner 不对）、ⓒ（钉逐字断"每枚表的列集合"）；
- **技术上不响钉但必须人工批准、本程判不可取**：ⓓ、ⓔ；
- **且没有任何一形躲得过 `PLAN.md:1363` 的 C13 射程**（会话 ∉ L1/L2/L3）。

---

## C. 归档/保留的真实耦合（同时说清"没接线"和"接上之后会怎样"）

### C-1 "今天有没有任何机制真去清它"——答：**没有，且连该启动它的那枚容器都不存在**

本程自己复跑的尺与读数（不照抄上一枚腿）：
- `grep -rn "StartRetentionJob" cmd internal tools scripts --include='*.go' | grep -v _test` ⇒ **4 行**，逐行都是 `internal/memory/retention.go` 自家（`:95` 注释、`:98` 定义、`:100` panic 串、`:234` 注释）⇒ **生产调用者 0 枚**；
- 含测试的全集 ⇒ 5 行，唯一调用者是 `internal/memory/retention_test.go:198`；
- `grep -rn "RunRetentionOnce" cmd internal tools scripts --include='*.go' | grep -v _test` ⇒ **2 行**（`:233` 注释、`:235` 定义）⇒ **生产调用者 0 枚**。⚠ 更正一处（F9）：`RunRetentionOnce` 的**测试**调用者不是 0——`retention_test.go:150` 有一枚真调用（a1 那句"2 行都在注释里"只在 `grep -v _test` 的射程里成立）；
- **更硬的一层（a1 未量）**：`StartRetentionJob` 的第一枚实参是 `*plugin.DisposalScope`（`retention.go:98`，且 `:99-101` 逐字在没有它时 `panic("memory: StartRetentionJob requires a DisposalScope")`）。尺：`grep -rn "NewDisposalScope" cmd internal tools scripts --include='*.go' | grep -v _test` ⇒ **2 行，都在 `internal/plugin/disposal.go:152/:155`（注释＋定义）** ⇒ **生产里没有一枚 DisposalScope 被创建** ⇒ 今天"把归档接上"不是补一行调用，而是要先把 C11 的 scope 装配在 cmd 侧建出来。
- `retention.go:233-234` 那句注释逐字"tests and the doctor command" ⇒ 与盘上不符：**`cmd/wisp/doctor.go` 里 retention 零命中**（尺：`grep -rni "retention" cmd/wisp/doctor.go` ⇒ **0**）。**这是一枚注释在声明一条不存在的接线**（`internal/observe/goroutine.go:58` 的 D38 名册也确实留着 `"retention-job"` 这枚名字，`goroutine_test.go:256` 还把它分类为 `CategoryTemporary`——名册认得、生产没人启动）。

⇒ 结论：**"归档"这件事今天就算做出来，也没有任何东西会去清它**。30 天窗口（`retention.go:32 TaskLogTTL = 30 * 24 * time.Hour`）是一张**空头支票**；连带后果是 `task_log`/`tool_call` 行**当前只增不减**（唯一会删它们的生产路径是 `PurgeTaskLogs`/`DeleteTaskLog`，而调用它们的隐私 API 生产调用者＝0，尺见 F6）。

### C-2 规格在哪一节要求它

`docs/specs/SPEC-02-data-storage.md:152` 节标题逐字"**## 4. 保留期与清理（有归属才有效）**"；`:154-158` 那张表逐枚给 `task_log`/`tool_call` 30 天、`cost_daily` 400 天、`approval_grant` 失效行 30 天，清理者一栏逐字写 `memory.RetentionJob`；`:162-163` 逐字"**进程启动后延迟 5 分钟跑一次 + 每 24h 一次**……挂在 `retention-job` 临时 goroutine"。
`docs/PLAN.md:2710-2713`（D35 规则 1）逐字"**保留期的清理必须有归属，否则等于没有**"。
⇒ 规格**要求了清理者是谁**，代码里**清理者存在但没人雇它上班**。属"安全/数据结论只因功能没接线才成立"的形状（见 C-3）。

### C-3 两种说法都要讲，不许只讲一半

- **"没接线"这一半**：今天（a）没有会话表；（b）有会话表也没人去清；（c）连启动它所需的 scope 都没在 cmd 侧被建过。所以"用户的归档会话会不会 30 天后消失"这一问，**当前答案是"什么都不会发生"**——既不会因归档而豁免删除，也不会因到期而消失。
- **"接上之后会怎样"这一半**：一旦有人把 `StartRetentionJob` 接进 cmd 的启动序列（并且先建出 scope），三件事同时变真：
  1. `deleteOlderThan`（`retention.go:195-203`）谓词逐字 `DELETE FROM %s WHERE %s IS NOT NULL AND %s < ?`——**只看时间列，不看 `state`、不看任何"归档位"** ⇒ **归档不会多活一天**，用户"归档"的会话会在窗口到期时被**连行删掉**（`DELETE`，不可撤销）。
  2. 若会话表被纳入窗口 ⇒ 上面那一发落到会话行；若**不**纳入 ⇒ 新表成为一张**永不被清**的表（`SPEC-02 §4` 的"有归属才有效"再落空一次，且 `privacy.go` 的删除族也得同批改）。
  3. `tool_call` 的年龄用的是 `COALESCE(ended_at, started_at, decided_at)`（`retention.go:154-155`）⇒ 会话删除与取证删除**不同步**：30 天那条线在两张表上算出的时刻不一样。
- ⚠ 一处"已经在跑"的近邻，别混：**artifacts 文件配额是真的有生产触发点**——`enforceArtifactsQuota` 在写路径上被 `artifacts.go:78-84` 逐次 `PutArtifact` 调用（注释逐字"Quota is enforced on the write path, not only by the retention job"）。⇒ **文件那一族的 LRU 有效、DB 行那一族的清理无效**，两族今天状态相反，说"保留期没跑"时必须限定为 DB 行。

---

## D. 改名会不会断取证链

### D-1 先把"改名"这一动落到哪枚列说清

改名**不碰主键**：目标是"标题"（票面 A3 那一格），今天盘上**没有标题列**（`schema.go:53-65` 的 11 枚名册）。⇒ **"改名断不断链"的第一答案是"当前无处可改"**；会断链的是**另一种实现**——把"改名"实现成**改 `task_log.id`**，或把 id 形状改掉（Bⓐ 腿一）。下面把这一动的耦合面逐枚钉出。

### D-2 指向 `task_log.id` 的每一枚（含软引用）

| 类型 | 位置 | 逐字 | 动 id 会怎样 |
|---|---|---|---|
| **FK 声明**（唯一一枚） | `schema.go:71` | `task_id TEXT NOT NULL REFERENCES task_log(id)` | **不拦、不级联**（见 D-3） |
| 索引 | `schema.go:85` `idx_tool_call_task ON tool_call(task_id, seq)`；`:86 idx_tool_call_corr ON tool_call(correlation_id)` | — | 孤儿行照旧被索引，只是指向空 |
| 读路径（取证唯一入口） | `dao_toolcall.go:111-112` | `toolCallSelect + " WHERE task_id=? ORDER BY seq"` | 改 id 后**这一发查不到任何行**＝"它到底做了什么"那屏变空白，**不报错** |
| 写路径 | `dao_toolcall.go:24`（INSERT 带 `task_id`）；`internal/agent/journal.go:76`（`TaskID: t.taskID`）；`:130`（`ID: t.taskID`，即 task_log 主键） | — | 同一串要同步改，改不齐就分叉成两个 id 空间 |
| **软引用（别名，不是分工）** | `internal/agent/loop.go:356` 逐字 `j := newTaskJournal(l.opt.Journal, taskID, taskID) // C18: correlation == task id` | `tool_call.correlation_id`（`schema.go:82`）今天**等于** task id | 任何按 correlation 归组的想法＝按 task 归组（a1 F8，本程复量成立） |
| 与 task 无关的两族（改名碰不到） | `approval_grant.session_id`（`schema.go:94`，铸点 `cmd/wisp/run.go:473 session.Mint()`）；`cost_daily(day)`（`schema.go:103`，写者 `dao_misc.go:132` 一带） | — | 授权键在"一次开机"上、成本按天聚合 ⇒ **改名/改 task id 都不动这两族** |
| 面板侧键（内存，不落库） | `cmd/wisp/run.go:1266,1293` 一带用 task id 当 stream key（a1 B1 坐标） | — | 改名不影响（不是持久引用） |

### D-3 `foreign_keys` 那把尺：**本程自跑复量，不照抄**

- 派单给的原尺 `grep -rn "foreign_keys" --include=*.go internal cmd` ⇒ **0 命中**（本程实跑）。
- 宽到四枚根＋build：`grep -rn "foreign_keys" cmd internal tools scripts docker build` ⇒ **4 行**，逐字全是 `Binary file build/wisp.exe matches`、`build/wisp228.exe`、`build/wisp77.exe`、`build/wisp77_partial.exe` ⇒ **源码零命中**（编译产物里的串是驱动自带符号，不等于我们开了它——这一句本程**不**当"已开"用）。
- 连接串那侧自读：`open.go:146-156 dsn()` 逐字只挂四枚 pragma——`busy_timeout`、`journal_mode(WAL)`、`synchronous(NORMAL)`、`wal_autocheckpoint`，另 `:153 "&_txlock=immediate"`；**`_pragma=foreign_keys(...)` 一枚没有**。驱动是 `modernc.org/sqlite`（`open.go:19` 空导入、`:37 DriverName = "sqlite"` 注释逐字"the modernc.org/sqlite registration name"）。
- ⇒ 可断言的部分：**本仓源码从未请求开启外键约束**。SQLite 的默认值（关）属"跑了才知道"的那一格 ⇒ 登记进 G1，并给复量法（`PRAGMA foreign_keys` 一发读数）。
- ⇒ 对 D 的结论（按当前证据）：**`REFERENCES` 是一枚声明了但本仓没有一行代码去激活的约束**；因此"改 id"的失败形状是**静默孤儿行**（不报错、不级联、`ListToolCallsByTask` 返回空），而不是"被数据库挡住"。⛔ 落地腿禁区一句话：**任何"改名"都不许实现成改 `task_log.id`。**

### D-4 另一枚会静默的形状（`SELECT *` 尺）

`grep -rn "SELECT \*" cmd internal tools scripts --include=*.go` ⇒ **0 命中**：全库读路径都写显式列名。⇒ 好处是加列不会自动改变既有行的形状；坏处是**新列不会自己出现在读路径**——忘了同改 `dao_tasklog.go:33-37`（INSERT 列名＋11 占位）、`:80-82`、`:95-97`（两枚 SELECT）、`:144-160`（按位置 Scan 的 `scanTaskLog`）这四处，**没有任何一枚测试会响**，只会静默串列或永远 NULL。这一格属"改了不报"，比红更危险。

---

## E. 料齐之后怎么问（零术语的人话后果，只写盘上有据的）

### E-1 如果什么都不做（会话继续只活在内存里）

`internal/agent/loop.go:187` 那枚 `history []llm.Message` 是对话的唯一载体（`History()` 只给内存副本，包内在 `:399`、`:567`、`:573` 三处读它拼 prompt），**没有任何一行代码把它写进库**（`grep -rn "\.History()" … \| grep -v _test` ⇒ 只有那三枚自家调用）。⇒ 他会看见的：

1. **关掉程序再打开，昨天那次对话就没有了。** 不是"找不到"，是**盘上从来没有过**。他能拿回的只有每一发的收尾一行（`task_log.query_text` 是他被脱敏过的原话、`summary_text` 是模型写的收尾摘要），而且——
2. **今天连那一行他也看不见**：列它的 SQL 有（`dao_tasklog.go:93`），但**没有任何出口在用它**（包外生产调用者 0；CLI 十枚命令名册 `main.go:89-124`：run/providers/doctor/secret/models/slo/panel-assets/panel-inbound/version/help，没有一枚列历史；面板快照字段只有 `composer.go:74,91` 那两枚可选 `instructions`/`tasks` 加 pending/results/composer）。⇒ 他要**看自己的历史，只能拿 SQLite 浏览器手动开 `%APPDATA%\wisp\wisp.db`**。
3. **改名／归档／分叉这三件事，他现在连"没有"都感觉不到**，因为没有可改的东西；而**"分叉"在今天的盘上没有可分叉的对象**（历史本体不落库）。
4. **附带一枚他此刻没抱怨但会撞上的**：既然清理没人启动（C-1），`task_log`/`tool_call` 行**只增不减**。今天看不出问题（用得少），**用得越久库越大**，而规格里那张"30 天"的表（`SPEC-02:154-158`）是给这台机器**承诺过的**。
   ⛔ 不许把它说成"他会不满意／他会流失"——那是推测，本件只给上面四条盘上有据的后果。

### E-2 如果做（动 schema），最坏代价是什么（三条，逐条给凭据）

1. **要人工批准，且批的是一串契约面**：`AGENTS.md` §1.1 逐字"C1–C32 / D1–D47 契约变更须人工批准"——这里至少四枚具名：`D35`（`PLAN.md:2692-2719`）、`C13`（`PLAN.md:1363` 的 L1/L2/L3 射程）、`SPEC-02 §3`（`:22` 标题"不得自行增删字段"）、`SPEC-02 §4/§5`（`:154-158`、`:167`）。**最坏情况是 owner 不批 ⇒ 票 166 的 5 件事里 4 件（搜/改名/归档/分叉）在现仓规下无解，只能闭"列"那一件（Bⓐ 腿一）。**
2. **17 枚会红的断言里有 3 枚是"改测试自己的锚"**（P13/P14/P15），而且 P15 那枚正是票面 AC#3 要的**崩溃原子性证据**：不修它就没法交 AC#3；修它＝改断言＝必须非实现者复判。⇒ 最坏形状是实现者为了变绿**放宽那枚原子性断言**——那是 `AGENTS.md` §1.2 的禁区，必须在派单里写明。
3. **迁移失败会不会丢历史**：按代码读，**不会丢，但会拒启动**。`applyMigrationStep`（`open.go:390-418`）单事务＋失败回滚；`migrate`（`:324-341`）在每步之前**先关连接、先写整库备份**（`backupDatabase:501-525`，命名 `wisp.db.bak-<from>-<to>`）；版本比二进制新 ⇒ `ErrSchemaUnmigratable`（`open.go:64-66` 逐字"refusing to touch wisp.db (restore the file from backup\ or remove it…)"，且 `TestMigrationNewerSchemaIsUnmigratable` 逐字断"the database file must be untouched"）。⇒ **真正的风险不是"丢历史"，而是"打不开"**：球进 `Unconfigured` 态、用户看到"要么从备份恢复要么删库"的指引，而这一步**今天没有任何界面替他做**（票 40 的隐私/诊断页那族 API 生产调用者＝0）。⚠ "备份真的完整"这一格只有跑了才知道（备份是文件级拷贝，`-wal/-shm` 存在时才带上，`open.go:511-517`）⇒ 登记进 G4，不替它背书。

### E-3 一句问法（摆给 owner 的原话）

"要不要让用户明天打开程序时能回到昨天那次对话？要——就批一次 schema 契约变更（动 D35 ＋ C13 ＋ SPEC-02 §3/§4/§5，并批准 17 枚断言同批改）；不要——那么票 166 里改名/归档/分叉三件应当**从票面上删掉**，而不是留成永远勾不掉的格子。"

---

## F. 推翻清单（派单与本会话里每一句都复量；不成立的具名推翻并给真身）

| # | 原断言（出处） | 本程复跑读数（尺＋命中数） | 判定 |
|---|---|---|---|
| F1 | "盘上根本没有'会话'这个实体"（a1 §0／本派单转述） | `ls internal/memory` ⇒ 九张表的 DAO 名册（`dao_tasklog/dao_memory/dao_profile/dao_toolcall/dao_misc/dao_providerhealth`）里**没有会话 DAO**；DDL `schema.go:22-146` 逐枚读过，无对话/消息/标题列；`internal/agent/loop.go:187` 是对话唯一载体 | **成立**。⚠ 但需补一枚**词面冲突**（派单没提、a1 只点到"无实体"）：**规格里"会话"另有所指**——`docs/specs/SPEC-05-agent-core.md:127-137`（"## 5. 会话与保活（B4/C31）……唤起 → 用户显式结束，或空闲 90s"）与 `docs/specs/SPEC-00-product-overview.md:143`（逐字"会话（Session） \| 唤起 → 显式结束或 90s 空闲"）说的是**一次开机内的语音会话**，不是票 166 要的跨重启对话。⇒ 落地腿在 spec 里加"会话"二字会**和既有术语撞车**，必须先定名（进 H1） |
| F2 | "`internal/memory/dao_tasklog.go:93` 有 `ListTaskLogs`，包外生产调用者 0" | `grep -rn "ListTaskLogs" cmd internal tools scripts --include='*.go' \| grep -v _test` ⇒ **4 行**：`:92` 注释、`:93` 定义、`privacy.go:76`、`privacy.go:214`（两枚调用都在 `internal/memory` 家里）；含测试 5 行 | **成立**（行号 `:93` 复量正确，本程锚点上仍对得上） |
| F3 | "`schema.go` 顶部注释逐字声明 DDL 与 SPEC-02 §3 逐字节相同" | 逐字读到 `schema.go:10-13`（引文见 A-1 表 A1）；`doc.go:9-10`、`models.go:13` 各重述一次 | **成立（注释确实这么写）**，但两条本程新量：① **没有任何测试执行这条**——四枚根里读 `docs/specs/**` 的 Go 文件＝0 枚（`grep -rn "docs/specs" cmd internal tools scripts --include='*.go'` ⇒ **4 行，全是注释**，无一枚打开文件比对）；② **它今天已被自己违反 2 处**（本程 `diff` spec 的 ```sql 块 `:25-123` 与 `schema.go:23-120`）：spec `:32` 是 `≤20 行` 而 Go `:29` 是 `<=20 行`；spec `:121` 是 `不符 → 拒绝加载` 而 Go `:118` 是 `不符 -> 拒绝加载`。⇒ **判读更新**："逐字节相同"是**意图声明**（且已不成立），真正的硬约束是 `schema_test.go` 那族自省钉 ＋ spec 标题那句"不得自行增删字段"。这**不放宽任何一格**（还是人工批准），只是把"要动什么"从"改字节"更正为"改语义＋改名册" |
| F4 | "那张钉断的是**表名集合**还是**每枚表的列集合**？"（派单 ⓑⓒ 的关键一问） | 逐字读 `schema_test.go:173-236`：`:195` 表枚数 ＋ `:198-208` 逐名 ＋ `:211-214`→`:259-260` 逐表**列枚数**（`t.Fatalf("column count = %d, want %d")`）＋ `:262-276` 逐位列名/类型/可空/主键 ＋ `:298-300` 索引枚数 ＋ `:222` 版本值 ＋ `:229/:233` 语句枚数 | **推翻"也许只断表名"这一支**：**两样都断**。⇒ "只加一枚可空列"不能绕过（Bⓒ）。顺带更正 a1 E5 的坐标：a1 写 `:173`（对，函数起点）与 `:259`（对，列枚数那行）、`:195`（对，表枚数那行）、`:222/:229/:233`（全对）；**a1 漏了 `:262-276` 逐位属性核**与 `:298-300` 索引枚数，也漏了 P11–P17 那 7 枚 |
| F5 | "`StartRetentionJob` 生产调用者＝0"（派单 C，a1 F11） | `grep -rn "StartRetentionJob" cmd internal tools scripts --include='*.go' \| grep -v _test` ⇒ 4 行全在 `retention.go` 自家（`:95/:98/:100/:234`）；全集 5 行，唯一调用 `retention_test.go:198`。**另加两把 a1 没跑的尺**：`grep -rn "NewDisposalScope" … \| grep -v _test` ⇒ **2 行，全在 `internal/plugin/disposal.go:152/:155`**（注释＋定义）；`grep -rni "retention" cmd/wisp/doctor.go` ⇒ **0** | **成立，且比 a1 更硬一层**：不光"没人启动它"，**它要求的那枚 `DisposalScope` 在生产里根本没被创建过** ⇒ 接线的代价不是补一行调用。另推翻 `retention.go:233-234` 那句注释（"tests and the doctor command"）：**doctor 不叫它**，注释在声明一条不存在的接线 |
| F6 | "归档/清理/隐私 API 都只在自家与测试里"（a1 A4） | `grep -rn "ListPrivacy\|ExportPrivacy\|PurgePrivacy\|DeletePrivacyItem" … \| grep -v _test` ⇒ **5 行，全是 `privacy.go` 自家注释/定义**（`:45/:47/:133/:135/:167/:170/:197/:198` 里非调用的那些 ＋ `:264` 自家 `purgeAllDomains`）；`grep -rn "purgeAllDomains" …` ⇒ **2 行，注释＋定义**（连测试都没叫它） | **成立**。⛔ 负向句自量：`ListPrivacy` 的命中只在测试（`privacy_test.go:53/74/92/138` 一族与 `artifacts_path_invariant_test.go` 巡过），本程自己跑过 `grep` 才写 |
| F7 | "a1 F2：`internal/store` 这枚包不存在，真身是 `internal/memory`" | `ls internal` ⇒ **23 枚目录**，名册逐枚看过（agent/audio/ball/buildinfo/config/llm/memory/models/observe/panel/perm/plugin/proc/projctx/risk/secret/session/speech/statemachine/streamkey/tools/watchdog/winsec），无 `store`；`find . -maxdepth 3 -type d -name "store"` ⇒ **0**（含 `-not -path "./.git/*"` 兜底） | **成立**（名册级推翻有效：本件全程按 `internal/memory` 引） |
| F8 | "a1 F4：生产可达工具＝恒在 8 枚（fs 六 + task 二）＋门控 1（`fs.delete`）＋子代理 1（`task.spawn`）" | 名册自量：`internal/tools/fs.go:124 fs.read`、`:190 fs.list`；`fs_edit.go:70 fs.edit`；`fs_write.go:234 fs.write`、`:369 fs.trash`、`:440 fs.move`、`:583 fs.delete`（`:598` 逐字"未获授权（[fs] delete_enabled=false），已拒绝执行"）；`task.go:484 task.output`、`:625 task.cancel`；`subagent_197.go:187 task.spawn`。注册三处循环：`cmd/wisp/run.go:519`、`:545`、`:791`（`reg.Register`） | **成立**（10 枚名册逐枚对上；票面引的"六枚"确为过期读数，台账 `A308` 那句需重读） |
| F9 | "a1 F11 补句：`RunRetentionOnce` ⇒ 2 行都在注释里，零调用" | `grep -rn "RunRetentionOnce" cmd internal tools scripts --include='*.go' \| grep -v _test` ⇒ 2 行（`:233/:235`）**对**；**但**含测试：`retention_test.go:150` 有一枚**真调用** `s.RunRetentionOnce(ctx, …)` | **半推翻**：**生产**调用者 0 成立；"**唯一调用者是测试里一枚/零调用**"这句措辞不实——测试调用者是有的，只是被 `grep -v _test` 排除。引用时请写"生产 0 枚 ＋ 测试 1 枚（`retention_test.go:150`）" |
| F10 | "a1 F10：`docs/PLAN.md` 十四态表真身＝`:3481-3494`，尺＝`awk 'NR>=3474 && NR<=3500 && /^| \*\*/ {c++}'` ⇒ 14" | 本程**复跑那一把原尺** ⇒ **27**（不是 14）。原因：awk 正则里的 `/^| \*\*/` 被读成"行首**或** ` **`"⇒ 匹配区间内**每一行**（3474–3500 恰 27 行）。**换成正确尺** `awk 'NR>=3474 && NR<=3500 && /^\| \*\*/ {c++}'` ⇒ **14**，逐枚首列名册＝思考中／SSE 流式输出／推理过程／工具调用／工具调用（展开）／审批等待（L2）／L2 确认卡（面板内）／L1 阻止窗口／L2 原生降级卡／错误／Stuck／成本／已取消 被打断／注入检出，行号首尾正是 `:3481` 与 `:3494` | **结论成立、尺是坏的**：坐标 `:3481-3494` 复量对；但 a1 写下的那把尺**不可复现**（跑出来 27）。⚠ 派单说"两枚坏尺今天栽在 `find -name` 上"——**这是第三枚**：`|` 未转义的 awk 正则。后续引用请写正确尺 |
| F11 | "a1 F7：`grep -rn '两枚改名' docs .scratch/wisp/issues` ⇒ 0" | 本程复跑同一把尺 ⇒ **1 行**，命中正是**票 166 面自己**（`.scratch/wisp/issues/166-…:18`）。另 `grep -rln "崩溃原子" docs .scratch` ⇒ **2 枚文件**（票面 ＋ a1 本件，与 a1 一致）；`grep -rn "minimax" docs .scratch \| wc -l` ⇒ **262 行**（分布在 `docs/reports/missing-features-2026-09-29-v4.md` 83 行、`survey-2026-09-28-composer-plus-menu.md` 27 行、`pending-and-issues.md` 27 行等），其中与"迁移"同现 **4 行** | **推翻两处**：① "两枚改名 ⇒ 0" 不成立（1 枚，且是票面自身——引"仓里没有对应物"这句时须写成"除票面外 0 枚"）；② "minimax 只出现在审批面"**偏窄**（262 行命中，含 4 行与"迁移"同现，需要逐枚读一遍才能判它是不是同一件事——本件**未读**，登记进 G7）。⚠ a1 的主判断"这条教训在本仓没有对应防线"**仍成立**（`journal_mode(WAL)` 是连接 pragma `open.go:147-151`，不是迁移手段），只是它的两条支撑读数得按上面条更正 |
| F12 | "a1 A3：`task_log` 十二枚列" | 三处双向数过：DDL `schema.go:53-65`、`contractSpec` `schema_test.go:82-94`、`models.go:58-70` ⇒ **恰 11 枚**（`id/started_at/ended_at/state/query_text/summary_text/cost_tokens_in/cost_tokens_out/cost_amount_micro/currency/error_class`）；`dao_tasklog.go:33-37` 的 INSERT 也是 11 列名＋11 占位 | **推翻**：**11 枚，不是 12 枚**（a1 自己列出的名册也只有 11 个名字，枚数写错）。这一枚对 Bⓒ/P5 有实际后果：**"加一枚列"是从 11 到 12，`contractSpec` 里那条 `t.Fatalf` 断的正是这个数** |
| F13 | "a1 A6：`Loop.Reset()` 全仓零生产调用者" | `grep -rn "\.Reset()" cmd internal tools --include='*.go' \| grep -v _test` ⇒ 4 行：`internal/llm/golden/replay.go:112`、`tools/d22scan/main.go:175,183`、`tools/mockllm/chat.go:431`（逐枚看过，全是 `cur.Reset()` 一类的 regexp/buffer 复位） | **成立**：`(*Loop).Reset`（定义在 `internal/agent/loop.go:265`）**没有任何调用者**，生产或测试都没有 ⇒ 连"会话结束"这个动作都没人按 |
| F14 | "a1 D2：`internal/panel/git_test.go:385` 那把前缀尺看不见已存在的 `config.*`" | 本程未复跑那三枚行的**内容**（属面板面，不是本票 A–E 射程），只复量了一条相关事实：入向六枚名册的补尺 `inbound_roster_253_test.go` 在派单里被 a1 具名，本件**不动那一格** | **未复核（不在本件射程）**，登记 G6。引用请引 a1，本件不背书也不推翻 |
| F15 | "`docs/PLAN.md:1874` 逐字：SQLite 历史'必须存'"（a1 A6 末行，本件 E-1 依赖它） | 本程读到 `docs/PLAN.md:1874` 一带为 §14.4 隐私族；该行逐字内容**未在本件锚点上重新打印核对**（`PLAN.md` 行号在 a1 与本程两个锚之间有漂移） | **量不到逐字复现**（见 G5）：给复量尺 `grep -n "必须存" docs/PLAN.md`。**本件的 E-1 不建立在这句上**——它建立在 `loop.go:187`＋`open.go` 的落盘面，已自量 |
| F16 | 本派单自己的一句："`internal/memory/retention.go:98 StartRetentionJob`" | `retention.go:98` 逐字 `func (s *Store) StartRetentionJob(scope *plugin.DisposalScope, cfg RetentionConfig) {` | **成立**（派单这一句坐标对） |

---

## G. 量不到的格子（禁跑 Go 的硬闸门下如实登记 ＋ 复量法；不含任何猜测）

| # | 量不到的那一格 | 为什么量不到（＋此刻在飞的写腿） | 复量法（无并发时跑） |
|---|---|---|---|
| G1 | **SQLite 外键到底开没开**（D-3 的那半句"默认关"） | 只能读源码；`PRAGMA foreign_keys` 的返回值是运行时读数 | 对真库发一发 `PRAGMA foreign_keys`（现代码里那四枚 pragma 之外的第五枚），读回 0/1；或在一次性小程序里 `INSERT INTO tool_call(task_id,…) VALUES('不存在',…)` 看是否被挡。**在得到读数之前，D-3 只可写"本仓源码从未请求开启"**，不许写"确认不生效" |
| G2 | P1–P17 那 17 枚钉**当前是绿的还是红的** | 需 `go test ./internal/memory/`，此刻 4 枚写腿在飞（`255-r3` `cmd/wisp`、`261-p1` `internal/llm`＋`internal/config`、`174-r2` `internal/agent`、`171-r3` 只写 probes） | 写腿清空后 `go test ./internal/memory/ -run 'TestSchemaContractIntrospection|TestMigrate|TestMigration|TestPrivacy' -v`，逐枚记 verdict；**基线四数之外要点名册差集**（票 145 AC#5 同规矩） |
| G3 | P13/P14/P15 那三枚连环钉**会不会以我预测的方式红**（本件只读到形状，没读到 verdict） | 同上；且这三枚要**先有生产 2→3 步**才谈得上撞车 | 复量法＝在**分支上真加一枚 no-op 的 `{from:2,to:3}`**（不改任何列）跑 `./internal/memory/` 看哪三枚红——这是最便宜的一发，且**不碰 schema 契约**（no-op 步不动表），仍属"要人工批准才能这么合入"的临时实验件，只在本地跑 |
| G4 | 票面 **AC#3 崩溃原子性现量**：中途被杀后"要么全成要么全不成"、**新库/老库各缺哪列**的实测 | 需要真起真 kill 的 subprocess | 照现成两枚形状各跑一发：`schema_test.go:458 TestMigrationFailedStepIsAtomic`（改锚后）与 `concurrent_test.go:179`（注释逐字"Crash recovery with a real subprocess"）；水位分别取"新库 0→1→2→3"与"老库已 2、只跳 3"，同批验 `schema_meta.schema_version` 与 `PRAGMA table_info(task_log)` 枚数**要么都动要么都不动**；并核 `backup/wisp.db.bak-2-3` 是否含 `-wal/-shm`（`open.go:511-525`） |
| G5 | `docs/PLAN.md:1874` 与 `docs/PLAN.md:2710-2716` 的**逐字**（F15 欠的那一格） | `PLAN.md` 行号在两个锚点间漂移，本件只读了 2686-2730 一节，**没读到 1874 那一行** | `grep -n "必须存" docs/PLAN.md` ＋ `grep -n "保留期的清理必须有归属" docs/PLAN.md`（词面尺，不跑 Go 就能补，本程**未跑**是因为不在本件 A–E 射程内必用；见 I 节） |
| G6 | 面板入向名册那三枚尺（`git_test.go:385/:394/:517`、`inbound_roster_253_test.go`、`composer_test.go:393`）在本锚上的行号与 verdict | 属面板面（票 145/194/253 的射程），本件**未复量**，以免与 174-r2/255-r3 在飞的写腿互洗 | 见 a1 D2/E1–E3 的坐标，改动前由**面板面的只读腿**在当刻锚上重跑一次；本件不背书 |
| G7 | "`minimax` 那 262 行命中里，有没有**一枚真的**在讲"迁移不是崩溃原子的"这条教训"（F11 ②的收尾） | 要逐枚读 262 行才能判，本件只做计数未做语义判读（不做"大概没有"这种话） | `grep -rn "minimax" docs .scratch \| grep "迁移"` 的 4 行逐枚打开上下文读；若确无对应物，请把票 166 硬约束 1 那句改成"外部经验、本仓无对应尺" |
| G8 | `scripts/portable-tests.sh` 的 census 读数（`internal/memory` 是否被覆盖率名册完整认领、新文件 `dao_session.go` 会不会落在尺外） | census 分支要跑 Go | `bash scripts/portable-tests.sh --scope=census`；重点核 `internal/memory` 在 scope 列表与 `core_pin` 两处的解析集（a1 E8 报的是 `internal/session/` 缺席那一格） |
| G9 | `tools/d22scan` 在本锚的被扫文件计数与 ban#8 现量（本件 A–E 若落地会新增 `internal/` 文件） | 禁跑 | 独立 module 里 `./tools/d22scan/d22scan.exe -root .`（别在根目录 `go vet ./tools/d22scan/`） |
| G10 | 会话界面那一腿今天有没有桩（票面"连列出来都没有出口"若前端已画了，归口会变） | 两层禁令：`frontend/**`、`design/**` **不读不引** | 归 owner 委托的前端会话自取；本件不代答，也不想象界面形状 |

---

## H. 未定义即停上报（碰到就停手，不按自己的判断填）

1. **"会话"这个词在规格里已另有所指**（F1）：`SPEC-05:127-137`／`SPEC-00:143` 的"会话"＝一次开机内的语音会话（90s 空闲即结束），票 166 的"会话"＝跨重启的对话。⇒ **落地前必须由人定名**（`session` 这枚表名会与 `internal/session`、`approval_grant.session_id` 正面撞名；`conversation` 则是**新造词**，spec 里没有）。本件 A 节全部按"`session` 表"这一**待批**假名测量，改名不影响面积结论。
2. **schema 契约 vs AC#5 互斥**（a1 H1，本件 F3 已把"逐字节"那一半更正为"自省钉 ＋ spec 标题那句"，**但互斥不解除**）：`SPEC-02:22` 的"不得自行增删字段"＋`schema_test.go` 17 枚钉 vs 票 166 AC#5 的 `docs/specs/**` 零字节。⇒ 需 owner 拍：**批一段 spec 追加（照 v2 的先例形 `SPEC-02:126-147`）**，还是"会话不进 SQLite"，还是"票 166 的 AC 缩到只做'列'"。
3. **C13 的射程不含会话**（A-3 最后一行，`PLAN.md:1363`）：即便一列都不加、只读现有 `task_log`，把会话语义放进 `internal/memory` 也是**扩 C13＝契约面**。⇒ 需人拍：C13 那行要不要改写，或会话 DAO 归哪一枚包（新包 ⇒ 又撞 `AGENTS.md` §2 里"Go module 组织名/命名核查（S1 建仓前，**阻塞**）"那一族未定案项）。
4. **票 196 未拍，归档形无地基**（Bⓐ 腿二）：`archived` 若进 `state` 列，是给一枚**无 CHECK（`grep -rn "CHECK(" internal/memory --include=*.go` ⇒ 0）**、词表已有两套的列加第五个词。⇒ 196 未定案前不写 166 的归档形。
5. **`SPEC-05:144` 那句"禁对话摘要堆积"与"存会话摘要"会不会正面打架**：该行逐字（在 L1 记忆里）"禁同步提取（+300–1500ms）、禁 embedding、**禁对话摘要堆积**"。若会话层的实现是"每轮存一条摘要"，读法上会被这枚禁令挡。⇒ **需人拍这一句的射程**（它禁的是 L1 提取路径，还是所有会话级摘要堆积？本件不判）。
6. **分叉 AC#4 今天无可验对象**（票面对象不存在，a1 H5 同格）：历史本体不落库 ⇒ 分叉无原件可吞。⇒ 建议编排者把 AC#4 读成**条件 AC**，别让人造桩勾它（`issues/README`：不许拿 mock 代替真的假报完成）。
7. **"归档 30 天后清"这一形两头都不成立**（C-1／C-3）：谓词不看归档位（`retention.go:195-203`）＋ 清理没人启动（`StartRetentionJob`/`NewDisposalScope` 生产 0 枚）。⇒ 落地若写"归档的会话到期后消失"，那是**在承诺一条没有主人的规矩**（`SPEC-02:152` 节标题就是"有归属才有效"）。

---

## I. 自查

- **三节写满**：F（推翻清单）**16 条**，逐条给了尺＋命中数＋判定，其中 **5 条推翻/更正**（F9 半推翻、F10 尺坏、F11 两处推翻、F12 数字推翻、F3 判读更新），**2 条登记为不背书**（F14 未复核、F15 未复现）；G（量不到）**10 格**，每格都给了复量法或归属；E（人话后果）三段齐全（E-1 四条后果、E-2 三条最坏代价、E-3 一句问法）。
- **负向句全部自量**，未照抄上一枚腿：`foreign_keys`（0 源码／4 二进制）、`CHECK(`（0）、`SELECT *`（0）、`docs/specs` 被 Go 读取（4 行全注释、0 枚比对）、`NewDisposalScope` 生产调用（0）、`retention` 在 doctor（0）、`internal/store`（0，含 `find -name store` 兜底）、`CREATE VIRTUAL TABLE`（0）、`LIMIT` 在 `dao_tasklog.go`（0）、`RunRetentionOnce` 测试调用（**非 0**，F9 据此更正）。
- **引号里的字全部本程 grep/Read 到过**（含 spec 的 `:22/:131/:152/:162/:167`、`PLAN.md:1363/:2695/:2699`、`retention.go:233`、`open.go:64-66/:137`、`schema.go:10-13/:25/:193-195`、`schema_test.go:44/:177/:259/:299/:451/:481`、`dao_providerhealth.go` 与 `probe_health.go:114` 的"REPLACED wholesale"）。
- 计数尺一律 `;` 串接 ＋ `| wc -l` 收尾；零命中处用文字明写"0"，未用 `grep -c` 接 `&&`。
- **未跑任何 Go 命令**（`go test`/`build`/`vet`/`list`/`run` 一律未跑）；未产码一字；未改 `docs/**`、`internal/**`、`cmd/**` 任何文件；未碰任何 `^- [ ]` 勾选框（本件表格全部用 `| … |` 形，无勾选框）；未读未引 `frontend/**`、`design/**`。
- 全文无密钥值；出现的只有配置键名（`api_key_ref`、`[fs] delete_enabled`、`[privacy] keep_transcript` 一类）。
- 本件落 `.scratch/wisp/probes/166/a2/census.md`，单文件 commit，显式 pathspec；未 push。
- ⛔ 本件**没有替 owner 选形**：A 是面积、B 是可行性、C/D 是耦合，判"批不批"不在这个子代理的权限里。
