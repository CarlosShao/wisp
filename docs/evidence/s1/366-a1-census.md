# 票 366 / 腿 `366-a1` 普查（只读）— 第 10 步从未跑：孤儿可达、迁移原样放行、今天仍能开库

- 腿＝`366-a1`（**只读腿**，只做 `AC#0`）· HEAD＝`dd16b005cc6b5593d18d147a6dfbc36f613e19d6` · 2026-10-23 17:0x–17:4x +08
- 写面＝本件 ＋ 仓外 `D:\tmp\wisp\20261023\366a1\**`。仓内 `.go` **0 字节**；探针只在**仓外导出树**（`git archive HEAD`）里，⛔ 工作树突变窗口。

1. **`PRAGMA foreign_key_check`**：SQLite 的引用完整性检查，**不要求 `foreign_keys` 开关是开的** ⇒ 本腿所有读数都在**业务形状连接**（`dsn()`，FK 关＝与迁移同条件）上取，⛔ 用"开了 FK 才准"的假前提。
2. **`git grep -n '<pat>' HEAD -- <paths>`**：量的是 **HEAD blob**，⛔ 工作树（本腿运行期间另有腿在改 `cmd/wisp/run.go`／`internal/agent/loop.go`）。
3. **`splitSQLStatements(ddlV4)` 的返回长度**：量**代码**（真被执行的语句），⛔ 数量注释里的步骤名。
4. 夹具＝**真 v3 库**：生产链 `Open(dir, withSchemaTarget(3))` 跑 `0→1→2→3`（`ddlV1/ddlV2/ddlV3`）建出，⛔ 手抄 DDL；迁移一律走生产 `Open(dir)`（`migrate()`→`applyMigrationStep`→`applyV4`）。
5. **导出树指纹**：`schema.go`/`open.go`/`dao_tasklog.go`/`dao_toolcall.go`/`retention.go`/`models.go`/`journal.go` 逐枚 `git hash-object` ＝ HEAD blob；改前、gate 变体后、inject 变体后**各验一次**，三次全 `SAME`（§9）。

---

## §1 票面三条逐条复核（行号漂没漂一并报）

### 前提 A：官方 12 步的第 1 步（关 FK）——`196-r3` 真补上了？**是，且早于事务**

| 项 | HEAD 实测（`git grep` on HEAD） | 判 |
|---|---|---|
| 12 步**自陈**（Go 注释） | `schema.go:283-284` 逐字 "The rebuild follows SQLite's documented 12-step / procedure" | ✅ **未漂**（`196-v2`/`196-r2` 引的就是 `:283-284`） |
| 12 步**自陈**（SQL 内注释） | `schema.go:302` "重建走 SQLite 官方 12 步 ALTER 过程：建新表 -> 搬行 -> 删旧表 -> 改名 -> 重建索引。" | ✅ 只列 5 句 |
| 重建**实际语句数** | `len(splitSQLStatements(ddlV4))` ＝ **5**，逐句首行＝`CREATE TABLE task_log_v4`／`INSERT INTO task_log_v4 … SELECT`／`DROP TABLE task_log`／`ALTER TABLE task_log_v4 RENAME TO task_log`／`CREATE INDEX idx_task_log_started` | ✅ 5 句在 |
| 第 1 步落点 | `open.go:46 pragmaForeignKeysOff="foreign_keys(OFF)"` · `:183 func migrationDSN` · `:199 func migrationPragmas`（**删光 shared 里每条 `foreign_keys` 再追加恰好一条 OFF**）· `:365/:402` 两条迁移建连 · `:465 tx, err := db.BeginTx` | ✅ **建连层、`BeginTx` 之前**（r3 引的 `:465` 未漂） |
| 实读 | 迁移连接 `PRAGMA foreign_keys`＝**0**；业务 DSN 串**零 `foreign_keys`**、业务池 `0/0` | ✅ |

> ★**5 句不是缺陷**：`task_log` 只有 1 枚索引、0 触发器、0 视图 ⇒ 12 步里"重建索引/触发器/视图"本来就只有一句要做。缺的是**结尾那道校验**，不是句数。

### 前提 B：第 10 步 `PRAGMA foreign_key_check` 从来没跑过？**是（含注释都没有）**

- `git grep -n 'foreign_key_check' HEAD -- '*.go'` ⇒ **rc=1、0 命中**（连一句注释都没提）；`grep -c` 同 rc=1。
- 对比：`12-step` 在 `.go` 里 6 处命中、`foreign_keys` 41 处 ⇒ 尺不是"扫不到"——**这个词真的不在代码里**。

### 前提 C：`196-v2` 造孤儿库实测"迁移前后都＝1"？**方向真，数字不可当常数**

- 我造 **5 枚**真 v3 库、逐枚迁移，`fkc` 读数为 **3/2/1/1/0**（顺序见 §3）⇒ **迁移前后条数全部相等（Δ=0）**，且**迁移既不清孤儿、也不造孤儿**（`tool_call` 行与 `task_id` 逐字保留）。
- ⇒ "迁移前后都＝1"在**它那一枚库**上成立（我那枚 `manual-insert` 也正好 **1→1**）；但 **`1` 是夹具行数，不是常量**。AC#1 若拿"=1"当判据会误判。

### 前提 D：三条孤儿路径的行号

| 票/件引 | HEAD 实测 | 漂否 |
|---|---|---|
| `PurgeTaskLogs` 的 DELETE `dao_tasklog.go:134` | `:134` 逐字 `res, err := tx.ExecContext(ctx, "DELETE FROM task_log")`（原文为反引号串） | ✅ 未漂 |
| `DeleteTaskLog` 的 DELETE `:118` | `:118` 逐字 `… "DELETE FROM task_log WHERE id=?", id` | ✅ 未漂 |
| 注释"`tool_call` 归自己的域"（票面文字） | `:129-130` "…Forensic / tool_call rows are governed by their own domain and retention window." | ✅ 在 |
| 台账（`pending-and-issues.md:23701`）引"`open.go:146 func dsn`" | **`:154`** | ⚠ **+8 漂**（事实未变：`dsnPragmas()` `:161-168` 四枚 pragma、**无** `foreign_keys`） |
| `196-v3` 引护栏报错串 `schema.go:416-421` | **`:417-422`** | ⚠ **+1 漂** |

---

## §2 核心读数：迁移后 `foreign_key_check`（含**无孤儿正控**）

| 夹具 | 迁移**前** fkc | 迁移**后** fkc | 迁移后 `schema_version` | 迁移后 `task_log`/`tool_call` 行数 |
|---|---|---|---|---|
| **`clean`（无孤儿正控）** | **0** | **0** | 4 | 2 / 3（子行指向 `task-A,task-A,task-B`，全在） |
| `manual-insert`（票面 `196-v2` 那一形） | **1** | **1** | 4 | 2 / 4（`ghost-missing-task`） |
| `purge-tasklogs` | **3** | **3** | 4 | 0 / 3 |
| `delete-tasklog` | **2** | **2** | 4 | 1 / 3 |
| `retention`（§4） | **1** | **1** | 4 | 1 / 2 |

★**正控成立**：无孤儿库读 **0**、有孤儿库读 **>0**，且两向都随"删父/加悬空子行"变化（`purge` 是 `0→3` 同一枚库）⇒ **尺不是恒 >0**。
★尺的**射程**：`foreign_key_check` 只查**当前 DDL 里声明了的** FK；`tool_call.task_id REFERENCES task_log(id)` 声明在 `schema.go:82`（`ddlV1`）⇒ 在射程内。`session_turn.session_id REFERENCES session(id)`（`schema.go:261`）今天**没有任何删除路径**（`DELETE FROM session` 零命中）⇒ 不可达，故 ⛔ 列为孤儿来源。

---

## §3 三条孤儿路径逐条（每发都给：迁移前／迁移后／孤儿行数／指向谁）

| # | 动作（**生产调用**，非手搓 SQL） | 迁移**前** fkc | 迁移**后** fkc | **孤儿行数** | 子行指向谁（父还在吗） | 今天 `Open()` |
|---|---|---|---|---|---|---|
| ① | `PurgeTaskLogs()` ⇒ `purge_rows_affected=2`、`purge_err=<nil>` | 0 | **3** | **3** | `task-A,task-A,task-B`（**父全没**，`task_log=0`） | **`nil`／v4** |
| ② | `DeleteTaskLog("task-A")` ⇒ `delete_err=<nil>` | 0 | **2** | **2** | `task-A,task-A`（父已删；`task-B` 仍在 ⇒ `task_log=1`） | **`nil`／v4** |
| ③ | **手工 INSERT**（raw SQL，绕开 DAO）`task_id='ghost-missing-task'` ⇒ `raw_insert_err=<nil>` | 0 | **1** | **1** | `ghost-missing-task`（从未存在） | **`nil`／v4** |

- 三条**全部造得出孤儿**（⛔ 无"造不出"），且**迁移原样放行**：`fkc` 行逐字 `tool_call|rowid=N|parent=task_log|fkid=0`。
- ③ 的**同类**（同一根因，非手工）：**生产 DAO `InsertToolCall` 也不校验父存在**——`models.go:160-183 validateToolCall` 只查非空/枚举，`task-missing` 直插成功（`dao_dangling_insert_err=<nil>`、`id=4`）。⇒ ③ 不只在"手工"才可达（见 §4-⑤）。

---

## §4 ★**第四条路径（具名上报）**：retention 清理

- **代码**：`retention.go:146` 删 `task_log WHERE started_at < cutoff`；`:154-155` 删 `tool_call WHERE COALESCE(ended_at, started_at, decided_at) < cutoff`；而 `:197` 的谓词是 `… IS NOT NULL AND …`。
- **可达形状**：`internal/agent/journal.go:75-82` 的 `startCall` 插入的行**三列全 NULL**（`DecidedAt/StartedAt/EndedAt` 一个都不设）⇒ `COALESCE(...)` ＝ NULL ⇒ **永不被选中**，父却按 `started_at` 老化 ⇒ 孤儿。
- **实测**（真 v3 库：`task-old` 起于 31 天前 ＋ 1 枚 NULL 时间戳子行 ＋ 1 枚老时间戳子行 ＋ `task-fresh` 对控）：`RunRetentionOnce` ⇒ `tasklog_deleted=1`、`toolcall_deleted=1`，**留下 1 枚孤儿**（`NULL/NULL/NULL` 那枚，指向已删的 `task-old`），`fkc` **0→1**，迁移后仍 **1**，今天 `Open()`＝`nil`／v4。
- **射程（⛔ 夸大）**：`StartRetentionJob` 今天**零生产调用者**（`git grep` 只命中定义 `retention.go:98` 与 `retention_test.go:198`）⇒ 这条路径**代码可达、尚未接线**；一旦接线即生效。
- **第五条（代码可达、只量到最后一跳）**：`agent/loop.go:390-392` 里 `startTask` 失败**只 Warn 不中止**，随后 `:670/:909/:941` 的 `startCall` 照写 `tool_call`，`task_id` 指向**从未建成的父**（DAO 不校验、FK 关）⇒ 也能造孤儿。本腿只量了最后一跳（`InsertToolCall` 接受悬空 `task_id`），⛔ 冒充端到端。

---

## §5 升级会不会炸：今天能开？做**硬门槛**后还能开？

**今天（HEAD，FK 关）**：5 枚夹具**全部** `Open()` 成功——`open_err=<nil>`、`schema_version=4`、行数与子行 `task_id` **一字未动**（§2/§3 表）⇒ ★**"今天的优先级＝库必须能开"成立**。

**若把第 10 步做成硬门槛**（探针在 `applyMigrationStep` 里 `step.apply` 之后、`UPDATE schema_version`/`Commit` **之前**插 `PRAGMA foreign_key_check`，`>0` 即 `Rollback`＋拒）——**同一批夹具、同一份数据**：

| 夹具 | 硬门槛后 `Open()` | 逐字读数 | 拒后库状态 |
|---|---|---|---|
| `clean` | **能开** | `open_err=<nil>`、`schema_version=4` | v4 |
| `purge-tasklogs` | **⛔ 打不开** | `…366-a1 probe gate: PRAGMA foreign_key_check reports 3 orphan row(s) after the 3->4 rebuild; refusing to open the database (SQLite 12-step ALTER, step 10)` | **仍 v3**、`task_log=0`/`tool_call=3` 原样 |
| `delete-tasklog` | **⛔ 打不开** | 同上，`reports 2 orphan row(s)` | **仍 v3**、1/3 原样 |
| `manual-insert` | **⛔ 打不开** | 同上，`reports 1 orphan row(s)` | **仍 v3**、2/4 原样 |
| `retention` | **⛔ 打不开** | 同上，`reports 1 orphan row(s)` | **仍 v3**、1/2 原样 |

- ★**门槛位置决定后果**：本探针把检查放在**事务内、提交前**（＝12 步的第 10 步），拒 ⇒ `schema_version` **停在 3**；`migrate()` 的早退是 `open.go:374 if version == s.target` ⇒ **每次启动都会再拒**（永久打不开，不是一次性告警）。若放在**提交后**，下次启动 `version==target` 直接早退 ⇒ 只拒一次，第二次就开了（"告警"而非"门槛"；这半句是读 `open.go:374` 的机理推断，未另建树实测）——AC#1 若选硬门槛，**这两处不能混**。

---

## §6 正控：`196-r3` 的修补**真在**（无孤儿 ＋ FK 开 ⇒ 必须能开）

探针变体树＝HEAD ＋ `dsnPragmas()` **只加一行** `"foreign_keys(1)"`（＝"将来那一改"真的发生；`diff -u` 见 `diff-open-inject.txt`）：

| 读数 | 实值 |
|---|---|
| 业务 DSN 串 | `…wal_autocheckpoint(1000)&_pragma=foreign_keys(1)&_txlock=immediate` |
| 迁移 DSN 串 | `…wal_autocheckpoint(1000)&_pragma=foreign_keys(OFF)&_txlock=immediate`（注入的 `(1)` **没进**） |
| 迁移连接实读 `PRAGMA foreign_keys` | **0** |
| `clean`（无孤儿＋FK 开）`Open()` | **`err=<nil>`**、`schema_version=4`、`task_log=2`/`tool_call=3`、子行仍 `task-A,task-A,task-B`、`fkc=0` |
| 同一树上业务池 `s.reader`/`s.writer` | **1 / 1**（调用方的 FK 开**被尊重**，迁移连接没被连累） |
| 同树上再插悬空子行（DAO） | `memory: insert tool_call: constraint failed: FOREIGN KEY constraint failed (787)` ⇒ 业务侧**真在拦** |

⇒ ★**"注释写得好"被排除**：迁移专用 DSN 在**建连层**把 FK 关掉、`BeginTx` 之前生效，**同一进程里业务连接仍是 1**。
★**顺带的代价读数（AC#1 会用到）**：FK 真开时**三条孤儿路径当场 787、一行都写不进去**——`purge_rows_affected=0`＋`purge_err=…(787)`、`delete_err=…(787)`、`raw_insert_err=…(787)`、retention `deleted=0`。⇒ "把 FK 打开"不是免费修法：**它会同时让「清空任务日志」「删单条任务」变成永远失败**（子行还在时）。

---

## §7 `AC#0③` 口径：今天有多少**实际库**可能有孤儿？——★**本机量不出，如实说**

- 本机 `%APPDATA%\wisp\` **存在但只有 `logs/`**（无 `wisp.db`）；`%LOCALAPPDATA%\wisp`、`%USERPROFILE%\.wisp` **不存在** ⇒ **没有一枚真用户库可跑 fkc** ⇒ ⛔ 编数字。
- 能给的只有**判据（不是计数）**：一枚库有孤儿 ⇔ 它历史上跑过 ①`PurgeTaskLogs` ②`DeleteTaskLog` ③绕开 DAO 的写 ④retention（接线后）⑤`startTask` 失败后继续跑，**且当时 FK 关**。
- ★**可达性不是边缘**：①是**用户一键**（`privacy.go:177` → `PurgePrivacy(PrivacyTaskLog)`），而 `tool_call` 按设计**不清**（`dao_tasklog.go:129-130`）⇒ **凡点过"清空任务日志"的库，必然留孤儿**。
- ★**从文件本身判不出来**：SQLite 不记录"这条子行当年有没有父" ⇒ 唯一的自然时机是**迁移时跑一次 fkc**（这正是第 10 步的位置）；这也解释了为什么"基线化"（AC#1〔甲〕）在工程上比"事后追查"更可行。

---

## §8 门禁 rc（逐门分开、⛔ 相减、⛔ 走管道取 rc）

| 门 | 命令 | rc |
|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **0**（0 字节输出） |
| G2 | `GOFLAGS= go vet ./internal/memory/` | **0**（0 字节输出） |
| G3 | `GOFLAGS= go test -count=1 ./internal/memory/` | **0**（`ok … 20.469s`；HEAD 导出树另跑一发 `25.694s`，亦 0） |
| G4 | `sh scripts/d22scan.sh` | **0** |
| G5 | `sh scripts/check-path-length-budget.sh` | **0** |

★G1–G5 在**共享工作树**上跑，取 rc 用 `cmd > out 2>&1; rc=$?`（⛔ `| tail`）。跑门禁时工作树里有**别的腿**在飞的改动：`--numstat`＝`cmd/wisp/run.go 7/11`、`internal/agent/loop.go 35/12`（⛔ 本腿的；本腿运行开始时 `git status -uno` 里没有这两枚）。⇒ 另在**HEAD 导出树**上复跑 G1–G3：**0/0/0**。

---

## §9 纪律格：提交号 · 逐枚 `.go`（应为 0）· 正控尺

- 提交号：`<填于提交后>`（只 commit、⛔ push；`add`/`commit` 两处显式 pathspec，逐个文件名字面点名）。
- **本腿落进仓的 `.go`＝0 枚**。判据（三把，⛔ 只看一把）：①`git status --porcelain -- '*.go'` 里没有本腿新增/修改（只有别的腿的两枚 `M` 与一批别的腿的 `??`）②本件之外本腿只写了 `docs/evidence/s1/366-a1-census.md` ③探针文件路径在 `D:\tmp\wisp\20261023\366a1\probes\**`（仓外）。
- **正控尺（导出一致性）**：`git hash-object` 逐个比 `git rev-parse HEAD:<path>`——`schema.go 60c9b69c…`／`open.go 43726beb…`／`dao_tasklog.go 8ef59d1c…`／`dao_toolcall.go e53b7485…`／`retention.go 2846718e…`／`models.go 2166333a…`／`journal.go 82d6db71…`；gate、inject 两个变体用完**都还原并复验**（3 次全 `SAME`）。
- 变体 diff 两枚（仓外）：`diff-open-gate.txt`（+18 行：门槛＋计数函数）、`diff-open-inject.txt`（+1 行：`foreign_keys(1)`）。

---

## §10 未定案／请编排者拍（⛔ 本腿自行假设）

1. ★**"官方 12 步第 10 步"的逐字条件句我没能复核**：本环境 `web_fetch` 对本机 DNS 一律解析成非公网 IP，`sqlite.org` 取不回来（3 次全失败）⇒ 票面/件里"第 10 步＝重建后 `foreign_key_check`"这句**我只能核到仓内二手**（`196-v2-verdict.md:72`）。若官方原文是"**若原本开着 FK** 才要查"，那本仓（FK 一直关）加这一步就是**超出官方步骤的额外保险**——结论不变（孤儿真实存在、迁移不校验），但**措辞要改**，请裁。
2. ★**第四条路径（retention）要不要单独立票**：今天没接线，但 `journal.startCall` 的 NULL 时间戳形状是**生产代码**产出的。
3. ★**"把 FK 打开"这条路的代价**已量（§6 末）：`PurgeTaskLogs`/`DeleteTaskLog` 会当场 787 ⇒ 不能当"顺手修"。
