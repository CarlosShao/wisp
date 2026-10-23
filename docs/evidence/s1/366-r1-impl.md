# 票 366 / 腿 `366-r1`（写腿）— 第 10 步落地：迁移里跑 `foreign_key_check`，**只记录、⛔ 拒开库**

- 腿＝`366-r1`（写腿＝`AC#2`）· 起始 HEAD＝`4ee57d318d73ea43890400b7aa4dab59b5127eeb` · 2026-10-23 17:5x–18:4x +08
- 形＝〔甲〕**基线化**（编排者已拍）。写面＝`internal/memory/**`（2 枚产码 ＋ 1 枚旧测试签名适配 ＋ 1 枚新用例）＋本件；写面外**0 字节**。
- 仓外树（⛔ 工作树突变窗口）：`D:\tmp\wisp\20261023\366r1\{headtree,posttree,pristine,mut,logs}\**`；变异两发全在 `posttree`，本仓工作树**全程未当过突变靶**。

---

## §1 先复跑 `366-a1` 的核心读数（⛔ 抄它的读数 —— 下面是我自己的发数）

**尺的射程（逐条写清，⛔ 只说"跑过了"）**：夹具＝**真 v3 库**（生产链 `Open(dir, withSchemaTarget(3))` 跑 `0→1→2→3` 建出，⛔ 手抄 DDL）；孤儿走**生产调用**（`PurgeTaskLogs`／`DeleteTaskLog`／raw INSERT 绕 DAO／`RunRetentionOnce`）；迁移＝生产 `Open(dir)`；`fkc`＝在**业务形状连接**（`dsn()`，`foreign_keys` 关）上数 `PRAGMA foreign_key_check` 的**行数**（`366-a1` §1 已证：该检查**不要求**开关是开的）。

| 夹具 | 迁移**前** fkc | 迁移**后** fkc | 迁移后 v | task_log / tool_call | 迁移后 `Open()` | 再开一次 |
|---|---|---|---|---|---|---|
| **`clean`（无孤儿正控）** | **0** | **0** | 4 | 2 / 3 | `err=<nil>` | `err=<nil>` |
| `manual-insert`（raw，绕 DAO） | **1** | **1** | 4 | 2 / 4 | `err=<nil>` | `err=<nil>` |
| `purge-tasklogs`（`PurgeTaskLogs`） | **3** | **3** | 4 | 0 / 3 | `err=<nil>` | `err=<nil>` |
| `delete-tasklog`（`DeleteTaskLog`） | **2** | **2** | 4 | 1 / 3 | `err=<nil>` | `err=<nil>` |
| `retention`（`RunRetentionOnce`） | **1** | **1** | 4 | 2 / 4 | `err=<nil>` | `err=<nil>` |

- 与 `366-a1` 的 **3/2/1/1** 对得上（`purge 3`／`delete 2`／`manual 1`／`retention 1`），**`clean` 读 0** ⇒ 尺不是恒 `>0`。
- 迁移**既不清孤儿、也不造孤儿**：`tool_call.task_id` 逐字保留（`task-A,task-A,task-B` 等）。
- 我这枚 `retention` 夹具的**行数**与 `366-a1` 不同（它 1/2、我 2/4，因我沿用了 `task-A/task-B` 的正控种子），**孤儿数相同（1）** ⇒ 报的是"孤儿行数"，⛔ 不是"表行数"。
- ★`366-a1` §10-1 的**未定案**（官方 12 步第 10 步的逐字条件句，它 `web_fetch` 取不回 sqlite.org）**本腿同样取不回**，故我按票面措辞落，**没有**把它当成已核事实。

## §2 判决：**记在哪儿**（"下一个读它的人怎么看到"）

**记两处，缺一不可**（编排者倾向的那形）：

1. **日志（运维看）** —— 每发被审计的迁移步骤，**提交成功后**打一行：
   - 干净：`level=INFO msg="migration baseline: PRAGMA foreign_key_check found no dangling references" orphan_rows=0 schema_version=4`
   - 有孤儿：`level=WARN msg="migration baseline: PRAGMA foreign_key_check found dangling references; recorded and opening anyway (rows are never deleted here - run \`PRAGMA foreign_key_check\` on the database to list them)" orphan_rows=3 child_tables=tool_call schema_version=4`
   ⇒ **含**"怎么自己看行"的那句 SQL；`child_tables` 指出孤儿在**哪张子表**（本仓＝`tool_call`）。
2. **迁移结果里的一枚计数（测试/调用方看）** —— `Store.ForeignKeyAudit() ForeignKeyAudit`：`{Checked, SchemaVersion, OrphanRows, ChildTables}`。`Open()` 返回的 `*Store` 上就能读，⛔ 不用刮日志。

**为什么不是"只写日志"**：日志是**人**的通道，断言不了；计数是**尺**的通道，`Checked=false` 还能区分"没跑"与"跑了得 0"——这两件事在"基线化"里必须分得开（第二次启动不跑步骤，若把 0 当"跑过且干净"就是假绿）。
**为什么不落 `task_log`**：`task_log` 是**用户任务面**（会随 retention 老化、会被 `PurgeTaskLogs` 一键清空）⇒ 把运维读数写进去，既会被清掉，又污染审计面。**为什么不落台账/新表**：要动 DDL（硬边界 #2）且本形明说"只记录"，`ForeignKeyAudit` 已够。

## §3 改了什么（逐行；`git diff --numstat`）

| 文件 | numstat | 内容 |
|---|---|---|
| `internal/memory/schema.go` | **+122 / −1** | ①`schema.go:190-201` `migrationStep` 加 `auditForeignKeys bool`（**由步骤自报**，不按版本号硬编码）；②`schema.go:219-222` 3→4 条目置 `true`（**唯一**一枚）；③`schema.go:304-309` `ddlV4` 的 **Go 注释**补"第 10 步不在这条字符串里、它在 `applyMigrationStep`"（**SQL 一字未动**，见 §3 末）；④`schema.go:509-516` `pragmaForeignKeyCheckQuery`；⑤`schema.go:520-567` `type ForeignKeyAudit`（含"为什么⛔ 当门槛"的实测出处）；⑥`schema.go:569-592` `foreignKeyCheck`（在**调用方事务内**读、**必须排空**再写水位；`rowid` 用 `sql.NullInt64` 兜 `WITHOUT ROWID`）；⑦`schema.go:594-605` `foreignKeyChildTables`（去重＋排序）；⑧`sort` 进 import |
| `internal/memory/open.go` | **+113 / −20** | ①`open.go:103-105` `Store.fkAudit`；②`open.go:265-272` `Open` 收 `migrate()` 的审计并落到 `s.fkAudit`；③`open.go:329-345` `ForeignKeyAudit()`（**对外唯一读数口**，注释写清"读的是迁移那一刻，不是现在"）；④`open.go:390-455` `migrate()` 改签名 `(ForeignKeyAudit, error)`（**目标版本已是 v4 时 `Checked=false`**）；⑤`open.go:494-...` `applyMigrationStep` 改签名并在 `step.apply` 之后、**写水位与 `Commit` 之前**跑审计（＝12 步第 10 步的位置），**提交成功后才打日志**（⛔ 报一个被回滚的状态）；⑥`open.go:551-568` `auditForeignKeys`（**永不失败**：读不出就 Warn＋`Checked=false`）；⑦`open.go:570-...` `logForeignKeyAudit` |
| `internal/memory/taskstate_196r2_test.go` | **+2 / −2** | 只适配 `applyMigrationStep` 的新签名（`if _, err :=`／`_, err =`）。**196-r2 的护栏逻辑与断言一字未动**。 |
| `internal/memory/migration_fkcheck_366_test.go` | **新 426 行** | 5 枚用例（§5）。 |

**硬边界自证（⛔ 靠"我没动"这句）**：
- **`ddlV1`–`ddlV4` 的 SQL**：两棵树各打一次指纹 ⇒ `ddlV1 fda736e1…`／`ddlV2 3e9ece70…`／`ddlV3 d46ee18e…`／`ddlV4 c17241cd…`（`len` 4463/903/1286/1357）**逐字相同**（`logs/26-ddl-head.txt` vs `27-ddl-post.txt`，`diff` 空）。
- **12 枚名册**：`TaskStateRoster()` 两树同为 `[running done succeeded cancelled interrupted stuck error Thinking Settling Muted Stuck Error]`。
- **`196-r3` 的 `migrationDSN`／`migrationPragmas`／`dsnPragmas`**：diff 里**零 hunk**。**`196-r2` 的 `requireForeignKeysOff`**：diff 里**零 hunk**。
- **⛔ 清孤儿**：本腿**没有**任何 `DELETE`；四条孤儿夹具迁移后 `tool_call` 行数逐条不变（§4）。

## §4 四条必给读数（改后；仓外 `posttree`，`logs/11-after-probe.txt`）

**(1) `clean` 库迁移后：`fkc=0`、计数 `0`**（两形都给：全新库 0→4、干净 v3 库 3→4）
```
[clean] OPEN#1 open_err=<nil>  audit checked=true  schema_version=4  orphan_rows=0  child_tables=[]
[clean] AFTER#1 fkc=0  task_log=2 tool_call=3 schema_version=4
```

**(2) 四条孤儿路径各一发：库照常打开、计数如实**

| 路径 | `open_err` | `audit.Checked` | `audit.SchemaVersion` | **`audit.OrphanRows`** | `child_tables` | 迁移后独立 `fkc` | task_log/tool_call |
|---|---|---|---|---|---|---|---|
| `purge-tasklogs` | **`<nil>`** | true | 4 | **3** | `[tool_call]` | 3 | 0 / 3（**一行未删**） |
| `delete-tasklog` | **`<nil>`** | true | 4 | **2** | `[tool_call]` | 2 | 1 / 3 |
| `manual-insert` | **`<nil>`** | true | 4 | **1** | `[tool_call]` | 1 | 2 / 4 |
| `retention` | **`<nil>`** | true | 4 | **1** | `[tool_call]` | 1 | 2 / 4 |

★"独立 `fkc`"＝迁移**之后**另开一条业务连接再数一遍；**与审计计数逐条相等**（两把尺同值，⛔ 一把尺自证）。

**(3) "永久锁死"不复发：同一发孤儿库连开两次，两次都 `err=<nil>`**

```
[purge-tasklogs] OPEN#1 open_err=<nil>  audit checked=true   orphan_rows=3  schema_version=4
[purge-tasklogs] OPEN#2 open_err=<nil>  audit checked=false  orphan_rows=0  schema_version=4
[purge-tasklogs] AFTER#2 fkc=3  task_log=0 tool_call=3 schema_version=4   ← 孤儿仍在，数据一字未动
```
其余三形（`manual-insert`／`delete-tasklog`／`retention`）与 `clean` 的 `OPEN#1/OPEN#2` **全部 `err=<nil>`**。
★第二次 `Checked=false` 是**如实的**：`open.go` 的 `version == s.target` 早退 ⇒ 这一发**没跑任何步骤**；把 0 说成"跑过且干净"就是假绿。

## §5 新用例逐名 ＋ 为什么**今天响**（不修不响）

| 用例 | 钉的是什么 | 摘掉新跳后为什么红 |
|---|---|---|
| `TestMigration366AuditRecordsCleanDatabase`（`fresh database`／`clean v3 database` 两子例） | 正控：干净库必须 `Checked=true`、`OrphanRows=0`、独立 `fkc=0`；行数与子行 `task_id` 逐字保留 | `Checked=false`、`orphan_rows` 无从谈起 |
| `TestMigration366AuditRecordsOrphansAndStillOpens`（`purge-tasklogs`／`delete-tasklog`／`manual-insert`／`retention` 四子例） | **四条孤儿路径**：`err=<nil>`、v=4、计数 **3/2/1/1**、`child_tables=[tool_call]`、独立 `fkc` 同值、**行数不变**（⛔ 清孤儿） | 同上；且若把审计改成"`>0` 就拒"，这里**当场红**（见 §6 M2） |
| `TestMigration366AuditNeverLocksTheDatabase` | **本形与硬门槛的分水岭**：同一库连开两次都 `err=<nil>`；第二次 `Checked=false`；孤儿仍在、`tool_call.task_id` 不变 | 第一次 `Checked=false`/`OrphanRows=0` ⇒ 红 |
| `TestMigration366AuditIsLogged`（`clean`／`orphans` 两子例） | **日志半**：`orphan_rows=0`（Info）与 `orphan_rows=2`＋`child_tables=tool_call`＋`opening anyway`（Warn），且**无 `level=ERROR`** | 日志一行都不打 ⇒ 4 条 `strings.Contains` 全红 |
| `TestMigration366AuditFlagIsOnTheRebuildStep` | 审计只挂在**重建那一步**（`3->4` 恰好一枚），`Store` 只存一枚基线 | 不红（它钉的是声明面，不是那一跳；**如实记**） |

## §6 反向判据：三读数 ＋ **还原核**

变异**全部在仓外 `posttree`**（`cp` 变体 → 跑 → 还原），本仓工作树**未当靶**。两发变异都落在**最终字节**（gofumpt 之后）上。

| 变异 | 内容 | 读数 |
|---|---|---|
| **M1＝摘掉那一跳** | `open.go` 里 `audit = auditForeignKeys(ctx, tx, step.to, logger)` 整行去掉（审计永不取） | **红**：4 枚顶层用例／10 个子例；红行 `migration_fkcheck_366_test.go:207,232,278(×4),333,380,383,403(×4)` |
| **M2＝把记录改成门槛**（硬边界 #1） | 取完审计后 `if audit.OrphanRows > 0 { return audit, fmt.Errorf(...) }` | **红**：3 枚顶层用例／6 个子例（`RecordsOrphans…` 四子例 `Open refused a database with N orphan row(s)`、`NeverLocks…`、`IsLogged/orphans`）；红行 `:271(×4),329,392` |

**三读数**：① 变异红（上表两行，rc=1）② 还原绿（`logs/23-final-restore-green.txt`：`rc=0`、`PASS 5 / FAIL 0`）③ 原样绿（`logs/20-final-pristine-green.txt`：`rc=0`、`PASS 5 / FAIL 0`）。
**还原核（时机＝最后一发变异之后量的）**：`sha256sum` 四枚逐字相等 ——
`open.go 3d9f2655…`／`schema.go 3f2e62ea…`／`taskstate_196r2_test.go a34a2591…`／`migration_fkcheck_366_test.go 8730ae29…`
（`posttree` 与工作树各算一次，**四次比较全同**；⛔ 在变异之间量，那量到的是变异产物）。

## §7 四数：`go test -count=1 ./internal/memory/`（改前／改后**分开**，⛔ 相减）

| 读数 | 改前（HEAD 工作树） | 改后（交付字节） |
|---|---|---|
| rc | **0** | **0** |
| `ok` 行 | `ok … 21.252s` | `ok … 22.341s` |
| `=== RUN`（含子例） | **101** | **114** |
| `--- PASS`（**顶层**） | **59** | **64** |
| `--- FAIL` | **0** | **0** |
| `--- SKIP` | **2** | **2** |

★**尺的射程**：`--- PASS` 只数**顶层**（子例是缩进的 `    --- PASS`，本表**不计**）；`=== RUN` 含子例 ⇒ 两列**不同尺**，故分开列、⛔ 用差数当结论。

## §8 门禁 rc（逐门分开、⛔ 相减、⛔ 走管道取 rc）

| 门 | 命令 | rc |
|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **0**（0 字节输出） |
| G2 | `GOFLAGS= go vet ./internal/memory/` | **0**（0 字节输出） |
| G3 | `GOFLAGS= go test -count=1 ./internal/memory/` | **0**（见 §7） |
| G4 | `sh scripts/d22scan.sh` | **0**（`clean - no D22 ban violations`；`internal/` 637 枚 Go 文件） |
| G5 | `sh scripts/check-path-length-budget.sh` | **0**（`VERDICT GREEN`；新文件 46 字符，远在 hat 内） |

**格式两尺**：**甲＝入库 blob（判据）**：把 `git show <commit>:<path>` 取出到仓外再 `gofumpt -l` ⇒ **0 命中**；**乙＝工作树（只作诊断）**：`/d/work/base/gopath/bin/gofumpt.exe -l internal/memory` ⇒ **0 字节输出**（`[\\/]` 过滤后 0 行）。
★**CRLF 尺⛔ 用 `gofumpt -l`**（它读工作树）：`git show HEAD:<p> | tr -cd "\r" | wc -c` ⇒ `open.go=0`／`schema.go=0`／`taskstate_196r2_test.go=0`（**LF**，与入库一致）。

**提交后复跑（同一提交 `d4ce46f97` 的字节，⛔ 未改一行）**：G1–G5 **全 0**；G5 这一发才**看得见**新文件（它量 **tracked** 路径，提交前那枚 `??` 不在射程内）⇒ `VERDICT GREEN`；G3 本次 `ok … 26.134s`（机器上另有腿在飞，**rc 不变**，时长⛔ 当判据）。

## §9 纪律格

- **只 commit、⛔ push**；`add`/`commit` 两处**显式 pathspec**（5 枚文件名字面点名，⛔ glob）。
- ⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ 碰票 196 原句。
- **未清任何孤儿**（硬边界 #3）：本腿**没有**"必须清孤儿"的情形，故**无需停手上报**；"`>0` 才做得到"的情形**也未出现**（记录形不依赖它）。
- **留给下一位的未定案**：①`366-a1` §10-1 的官方 12 步逐字条件句仍未核（本机取不回 sqlite.org）；②本腿的审计**只在"跑过重建的那一发 Open"里有值**，之后 `Checked=false`——若将来要"每次启动都报"，那是另一张票（会在每次启动多一次全表检查的开销），⛔ 本腿自行扩。
