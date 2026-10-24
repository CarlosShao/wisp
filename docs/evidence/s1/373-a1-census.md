# 票 373 / 腿 `373-a1` 普查（只读，`AC#0`）— "删父不删子"三处**当场造孤儿**，但票面的**总根**判词要修两处

- 腿＝`373-a1`（**只读腿**，只做 `AC#0`；⛔ 落任何 `.go`、⛔ 改票面／台账／`Status:`、⛔ push、⛔ 翻框）。
- 开工 anchor `0fa9392634a2f4dbce44bbd6b6ccd9f37fea335a`；落笔时 HEAD＝`e485af5aa`（期间 `368-v1` 进 3 笔）。`git diff --stat 0fa93926 e485af5a` ＝**3 文件、183 insertions、零 `.go`** ⇒ 下面每条行号在**两个 HEAD 上都成立**。
- 写面＝本件 ＋ 仓外 `D:\tmp\wisp\20261024\373a1\**`。仓内 `.go` **0 字节**：3 枚探针全在仓外，经 `go test -overlay` 注入 `internal/memory`（`overlay-memory.json` 三个 `Replace` 项）。
- **工作树声明**：开工我现跑 `git status --porcelain -- internal/ cmd/ scripts/` ＝**空**；收尾仍**空**。别人的脏（`.scratch/**` 共 1265 行，其中 77 行 `.go`）**⛔ 碰、⛔ 清**。⇒ **无写面冲突，`368-v1` 的 `internal/memory` 读数为干净树**。

**尺的射程（落笔前五问）**：①**对象**＝生产链 `Open(dir)` 建出的**真 v4 库**（⛔ 手抄 DDL）；`PRAGMA foreign_key_check` 数的是**DDL 里声明了的 FK 的悬空行**，**不要求** `foreign_keys` 开着。②**形状**＝`fkc` 取**行数**＋`table|rowid|parent|fkid` 逐枚；行数一律 SQL `COUNT(*)`。③**代码 vs 注释**：调用者判定用 `git grep` 且剔注释行；`PurgeTaskLogs` 的"tool_call 归自己的域"是**注释**，单列。④**总数核对**：每次读数都核 `live_children + orphans == tool_call`（探针内不成立即 `t.Errorf`）。⑤**谁算的**：整数由 SQLite `COUNT(*)` 给，探针只搬运；**无 hash**。

---

## §1 三处逐条复核（含行号漂没漂）

| 票面 | 现跑 `git grep -n` @ `e485af5a` | 漂？ |
|---|---|---|
| `privacy.go:177` | `:177` 逐字 `return s.PurgeTaskLogs(ctx)`（`case PrivacyTaskLog:` 在 `:176`） | ✅ **未漂** |
| `dao_tasklog.go:134 DELETE FROM task_log` | `:134` 逐字 ``res, err := tx.ExecContext(ctx, `DELETE FROM task_log`)``（`PurgeTaskLogs` 定义在 `:131`） | ✅ **未漂** |
| `DeleteTaskLog`（票面**未给行号**） | 定义 `:116`；DELETE 在 **`:118`** `` `DELETE FROM task_log WHERE id=?` `` | 补 `:116`／`:118` |
| `retention.go:197` | `:197` 逐字 ``q := fmt.Sprintf(`DELETE FROM %s WHERE %s IS NOT NULL AND %s < ?`, table, tsExpr, tsExpr)`` | ✅ **未漂** |
| `foreign_keys` 关（`open.go` 的 `dsn()`） | `open.go:166-173` `dsnPragmas()` **四枚**：`busy_timeout`／`journal_mode(WAL)`／`synchronous(NORMAL)`／`wal_autocheckpoint`，**无它**；**现量** `dsn()` 连接 `PRAGMA foreign_keys`＝**0** | ✅ 成立 |

★**票面哪里说得不准（要修的三条）**：
1. **`retention.go:197` 不是"只删子"的那一句，是父子共用的泛型**。同一函数在 `runRetention` 里被调**两次**：父用 `started_at`（`:146`）、子用 `COALESCE(ended_at,started_at,decided_at)`（`:154-155`），**同一 tx**（`:143`）。且 `task_log.started_at` 是 `INTEGER NOT NULL`（现量 `pragma_table_info` `notnull=1`）⇒ 那个 `IS NOT NULL` 在**父**上是恒真，在**子**上等价于"三列不全 NULL"。
2. **`ON DELETE CASCADE` 是死条款 —— 已实证，不只是读代码**：`PRAGMA foreign_key_list(tool_call)` ⇒ `on_delete=NO ACTION`；`ddlV1` 不含该子句（`strings.Contains(ddlV1,"ON DELETE CASCADE")=false`）；**C1 实测**：FK 关时**带 CASCADE 的子照样留 1 行**。
3. **`366-a1` 的 `787` 逐字复现**：`FK=ON` 下 ``DELETE FROM task_log WHERE id='task-A'`` ⇒ `constraint failed: FOREIGN KEY constraint failed (787)`。

---

## §2 三处**各一发**的读数（真 v4 库，生产入口）

夹具：`task-A`（2 子）＋`task-B`（1 子），生产 DAO `StartTaskLog`／`InsertToolCall` 落库；基线 `fkc=0`。

| # | 发起（生产入口） | 父删了 | `task_log` | `tool_call` | `fkc` | 孤儿指向哪个 `task_id` |
|---|---|---|---|---|---|---|
| P0 | 基线（对控） | — | 2 | 3 | **0** | — |
| **P1** | `PurgePrivacy(PrivacyTaskLog)` ⇒ `privacy.go:177` ⇒ `:134`，`rows_affected=2` | 全删光 | **0** | **3** | **3** | `task-A,task-A,task-B` |
| **P2** | `DeletePrivacyItem(PrivacyTaskLog,"task-A")` ⇒ `:118` | 删 1 | 1 | 3 | **2** | `task-A,task-A` |
| **P3** | `RunRetentionOnce`（`task-old` 起于 31 天前；子①40 天前、子②三列全 NULL、`task-fresh` 对控） | `task_log_deleted=1` | 1 | 2 | **1** | `task-old`（＝NULL 那枚） |
| **P3b** | 同库再过一发（`Now=+1d`） | `0`／`0` | 1 | 2 | **1** | 同上 ⇒ **永生** |
| **P4** | `RunRetentionOnce`（父 31 天／子 **5 天**） | `1`／`0` | 0 | **1** | **1** | `task-old` |

★**P4 是票面没写的第二个机制**：retention 的**父窗与子窗是两条独立语句** ⇒ 父到期、子未到期 ⇒ 照样留孤儿。P3 是"三列全 NULL"那枚（`368` 之后新行造不出、老行有），P4 是"窗口不齐"那枚。**两条都量到了**。

---

## §3 "删父"路径穷举（现跑，⛔ 抄票面）

**全仓 `DELETE FROM` 共 16 命中／10 处语句**（`internal/` 内，`git grep -n "DELETE FROM"`）。按父／子分类：

- **父 `task_log`：2 处** —— `dao_tasklog.go:118`（单个）、`:134`（全清）＋`retention.go:197` 经 `:146`（按龄）。
- **子 `tool_call`：2 处** —— `dao_toolcall.go:195`（单个）、`:210`（全清）。删子**不产生**孤儿。
- 其余 6 处（`profile`×3、`memory`×2、`approval_grant`×2、`plugin_state`×1、`cost_daily`×1）**无 FK 子表**。

★★★**第四处语句：具名报上** —— `schema.go:353` `DROP TABLE task_log`（在 `ddlV4` 重建里，与 `CREATE`／`INSERT..SELECT`／`RENAME`／`REINDEX` 同一 tx）。**我的判定：它不是第四个孤儿源**，三条理由：①同 tx 且 id 逐行搬过去；②`taskLogStatesOutsideRoster`（`schema.go:362`）在重建**前**先查，超名册的行会让整步失败而不是静默丢行；③`migrationDSN` 把 FK 强制 OFF（`open.go:188-213`），`taskstate_196r2_test.go` 已把这条钉住。**但它是本仓唯一一处"能让 `task_log` 整体消失"的语句 ⇒ 请编排者复核我这判定。**

★**第二对 FK 今天没有删父路径**：`session_turn.session_id REFERENCES session(id)`（`schema.go:286`，`on_delete=NO ACTION`）而**全仓无 `DELETE FROM session`**；新库 `session` 行数＝0。⇒ 孤儿类今天**只有 `task_log`→`tool_call` 一对**。

★★**但票面那句"孤儿只有一个总根"要修 —— 有第二个根，且它不是删父**：
`cmd/wisp/run.go:777` 给 Bridge 传 `Journal: mem` ⇒ `tool_call` 行在**跑任务时**由 `internal/tools/bridge.go:1225` 写入，而父行 `task_log` 在 `run.go:1243`（`bg.Wait()` 于 `:1221` 返回**之后**）才写 ⇒ **先子后父**。实测：**E1** 子落了、父从未落（`StartTaskLog` 失败时 `run.go:1244` 只 `auditf`）⇒ `fkc=1`、孤儿 `[task-never-booked]`；**E2** 同一顺序但父随后落 ⇒ `fkc=0`。⇒ **孤儿＝"子有父无"，父可以是"被删"，也可以是"从来没落"；〔甲〕只治前者。**

---

## §4 〔甲〕vs〔乙〕的代价

### 〔甲〕DAO 层显式同删 —— **3 个函数、3 条 SQL、0 处 DDL、0 处 schema 变更**

| 要改处 | 怎么改 | 事务 |
|---|---|---|
| `dao_tasklog.go:116-127` `DeleteTaskLog` | 在 `:118` **前**加 `DELETE FROM tool_call WHERE task_id=?` | 已在 `s.write(...)` 单 tx 内，**无需新管道** |
| `dao_tasklog.go:131-142` `PurgeTaskLogs` | 在 `:134` **前**加 `DELETE FROM tool_call`（整表） | 同上 |
| `retention.go:139-165` `runRetention` | 父删（`:146`）前先删**这批父**的子；或父删后补一次孤儿扫 | 父子**已经**在同一 tx（`:143`），只需多一句 |

**"父子同删必须在同一事务里吗"——必须在，且三处都**已经**满足**：否则中途崩就留孤儿。**B2 实测**（`FK=ON`、同一 tx、**子先父后**）：两句皆 `nil`、`commit=nil`、`fkc=0`；**B1 实测**父先子后（或不同 tx）⇒ `787`。⇒ **顺序是硬约束：先子后父。**

**风险（两条，都是语义不是实现）**：①`PurgeTaskLogs` 的注释**明写** `tool_call` 归自己的域（`dao_tasklog.go:129-130`）⇒〔甲〕**改的是这条产品语义**（清空任务日志会连带删取证行）⇒ `AC#1` 必须把它写成显式裁决，⛔ 当实现细节。②retention 若写成"全库 `NOT IN` 孤儿扫"，会顺带删掉 E1 形状（父从未落）的行 —— 也许正是想要的，但要明说。③〔甲〕**不治** E1（先子后父）那个根。

### 〔乙〕开 `foreign_keys` ＋ DDL 加 `ON DELETE CASCADE`

| 问题 | 现跑答案 |
|---|---|
| 要动 `ddlV1`（**冻结常量**）吗？ | ★**不必**。`ddlV4` 已立"重建而不改 `ddlV1`"的先例（`schema.go:304-324`）⇒ 用 `ddlV5` 重建 `tool_call` 加 CASCADE，`ddlV1` 文本与**语句数 15**（现量）一字节不动。⇒ **不触发"停手上报③"** |
| 要不要升 schema？ | ★**要**。`SchemaVersionTarget` 4→5（`schema.go:159`）＋新 `ddlV5`（v4 同形 5 条语句）＋迁移链加一步 |
| `dsn` 那行改多少？ | `dsnPragmas()` 加 `"foreign_keys(1)"`（`open.go:166-173`，**一行**）。`migrationDSN`／`migrationPragmas` 会**剔重后强制 OFF**（`:192-213`）⇒ v4 rebuild **不会**因此炸 |
| `366-a1` 的 `787` 复现？ | ★**能，逐字**（B1）。且 **C1** 证 CASCADE 是**必需**不是可选：`FK_ON`＋plain FK 父删 ⇒ `787`；`FK_ON`＋CASCADE ⇒ 子 0 |
| ★**它没量到的两个新代价** | ①**`cmd/wisp/run.go` 必须先改写序**：`FK=ON` 下 **B3 实测** `INSERT tool_call`（父未落）⇒ **`787`**，而 `bridge.go:1225-1227` 与 `journal.go:83-85` **把错吞了**（只 log／返回 0）⇒ **每发 `wisp run` 的取证行静默全丢**。②CASCADE **不追溯**清存量孤儿（只作用于将来的 DELETE） |

| 维度 | 〔甲〕DAO 显式同删 | 〔乙〕FK＋CASCADE |
|---|---|---|
| 改动面 | **3 个函数、≤3 条 SQL**；0 DDL、0 schema、0 迁移 | 1 行 DSN ＋ 新 `ddlV5`（5 语句）＋ `SchemaVersionTarget` ＋ 迁移链一步 ＋ **`cmd/wisp/run.go` 改写序** |
| 风险 | 语义：**改了"清空任务日志"的边界**（要 `AC#1` 裁）；漏一处即旧病复发 | 全局开关：任何"先子后父"的写者当场 `787`，且**两处吞错**会把它变成静默丢失；升级不可逆 |
| 要不要动契约 | **不动 DDL**；但要**明写** purge 的语义 | **动 schema 契约（v5）**＋ `SPEC-02 §3` 的 `tool_call` 块要加 CASCADE ⇒ **改契约文字，须摆机主** |
| 治 E1（父从未落） | ⛔ 治不了 | ✅ 治得了（`787` 是硬错，可惜被吞） |

---

## §5 正控（硬闸）：父子同删之后 `fkc` **必须** 0

先把库弄脏到 `fkc=3`，再在**同一 tx** 里发一句 ``DELETE FROM tool_call WHERE task_id NOT IN (SELECT id FROM task_log)`` ⇒ **`fkc=0`、`rows=[]`、`tool_call=0`**。⇒ **尺能读出"干净"**，上面的 3／2／1 不是尺坏。

---

## §6 "存量孤儿有多少"：★**今天能量，读数是 0，但是空库意义上的 0**

- ★**`368-a1` 的读数要更正（不是错，是不全）**：它量 `%APPDATA%\wisp` 只有 `logs\`、无 `wisp.db` —— **字面为真**。但真实数据根是 **`%APPDATA%\wisp-dev\`**，那里有 `config.toml`／`secrets\`／`artifacts\`／`wisp.db`。
- **口径（可复跑）**：`wisp-dev\wisp.db` 139,264 B、`wisp.db-wal` **0 字节**、`-shm` 32,768 ⇒ 干净收尾，**单文件拷贝即完整**。拷贝 SHA-256 与源**逐位相同**（`DA475F5D…4831`）；以 `mode=ro`＋`PRAGMA query_only=1` 读，**读后 hash 不变**（证明只读）。
- **读数**：`schema_version=4`；`task_log=0`／`tool_call=0`／`session=0`／`session_turn=0`；`PRAGMA foreign_key_check` **0 行**；孤儿 **0**；三列全 NULL **0**；`ended_at IS NULL` **0**。
- ★**不是"被人清过库"**：同源备份 `bak-1-2`（v1）与 `bak-3-4`（v3）也是 `task_log=0`／`tool_call=0` ⇒ 这台机器**从未落过一枚任务行**。
- ⇒ **存量孤儿＝0，⛔ 当作"级联已生效"的证据**；〔甲〕〔乙〕都**不需要**本机回填。别处机器的存量仍未知（本腿只有这一台）。

---

## §7 门禁 rc（全部 `cmd > out 2>&1; rc=$?`，⛔ 走管道）

| 门 | rc | 证据 |
|---|---|---|
| 探针 5 次（`-overlay`，仓外） | **0／0／0／0／0** | `out/probe.log`／`probe2.log`／`probe4.log`／`probe5.log`／`probe7-real.log` |
| `go build ./...` | **0** | `out/G1-build.log`＝**0 字节** |
| `go vet ./internal/memory` | **0** | `out/G2-vet.log`＝**0 字节** |

（`go test ./...` 全量**本腿不跑**：`373-a1` 只读、且别的腿在飞；探针的 `-run` 只选 `TestP373A1*`。）

---

## §8 提交与写面

- 提交＝本件**一枚**（`docs/evidence/s1/373-a1-census.md`）；逐枚 `.go` ＝ **0**（`git diff --cached --name-status` 交件时现跑核过）。
- 仓内改动的文件**总数＝1**；其余产出全在 `D:\tmp\wisp\20261024\373a1\**`（探针 3 枚、overlay 1 枚、日志 8 枚、真实库拷贝 3 枚）。

---

## §9 未定义即停（★请编排者在 `AC#1` 里逐条裁）

1. ★**〔甲〕会改"清空任务日志"的产品边界**：`dao_tasklog.go:129-130` 的注释明写 `tool_call` 归自己的域，而〔甲〕会让它连带被删。这是**语义裁决**，⛔ 我拍。
2. ★**第二个根（先子后父，E1）不在票面三处里**：〔甲〕治不了它，〔乙〕能治但会被两处吞错变成静默丢行。⇒ 要不要在本票一并治，或另开一小格。
3. ★**〔乙〕若被选**：`cmd/wisp/run.go` 的写序**必须**同批改（把 `StartTaskLog` 提到 `RunAsync` 之前，父行 `state='running'`，与 `journal.startTask` 同形），否则 `wisp run` 取证行全丢。这属"改契约文字"之外的**行为改动**，请一并裁。
4. ★**第四处语句 `schema.go:353 DROP TABLE task_log` 已具名**；我判它**不是**孤儿源（§3），但这个判定请复核。
