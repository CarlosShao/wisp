# 票 373 / 腿 `373-r1` 落地（写腿，`AC#2`）— 〔甲〕DAO 层显式同删 ＋ **先父后子**

- 腿＝`373-r1`（**写腿/实现者**）。写面＝`internal/memory/**` ＋ `cmd/wisp/run.go`（**只那处写序**）＋ 本件；仓外件全在 `D:\tmp\wisp\20261024\373r1\**`（basetree／base-clean／probe／mut／out）。
- 开工 anchor `7085bb3d4`。开工／收尾 `git diff --cached --name-status` 均**空**；`git status --porcelain -- internal/memory cmd/wisp/run.go` 开工**空**。别人的脏（`.scratch/**` 1265 行、`?? docs/evidence/s1/369-r1-impl.md`）⛔ 碰、⛔ 清。
- ⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。⛔ 动 `ddlV1`–`ddlV4`／schema／`dsnPragmas()`／`migrationDSN()`；⛔ 加 `foreign_key_check` 硬门槛；⛔ 给存量孤儿编时间戳；⛔ 动 `winsec/panel/perm` 等禁区。
- **尺的射程**：①对象＝`Open(dir)` 建出的**真 v4 库**（⛔ 手抄 DDL）；`fkc`＝`PRAGMA foreign_key_check` 的**行数**（**不要求** `foreign_keys` 开）。②形状＝整数＋逐行 `child|rowid|parent|fkid`；行数一律 SQL `COUNT(*)`。③代码 vs 注释：调用者判据用现跑 `git grep`，注释单列。④总数核对：每个 fixture 先读 `fkc=0` 基线，再读改动后的值。⑤谁算的：整数由 SQLite 给，本件无 hash 参与判定（文件 hash 只用于还原核）。

## §1 为什么选〔甲〕，不选〔乙〕（引普查 `787`／吞错两条）

- ★★**普查实测**：`FK=ON` 下插入无父的子 ⇒ **`787`**（`foreign_key_check` 一行字都不假）；而 `internal/tools/bridge.go:1225-1227`（`b.log(...insert failed...)` 后照跑）与 `internal/agent/journal.go:83-85`（`return 0`）**把这条错吞掉**。
- ⇒ ★★**〔乙〕（开 FK ＋ `ON DELETE CASCADE`）会让 `wisp run` 的取证行在"先子后父"的旧写序下静默全丢** ⇒ **新开一个洞来补另一个洞**。且〔乙〕要升 schema `4→5`、`dsnPragmas()` 加一行、`366-a1` 的 `787` 逐字复现、CASCADE **不追溯**存量。
- ★★★**我拍〔甲〕**：**3 个函数、3 条 SQL、0 DDL、0 schema、0 新管道**，三处**已经**在同一 `s.write` tx 内；**用最小的机制面关掉整个孤儿类**，⛔ 为了让"数据库自己保证"而引入一个会吞掉取证行的新失败模式。

## §2 改了什么（逐行）

| 文件 | 处 | 改法 |
|---|---|---|
| `internal/memory/dao_tasklog.go:124` `DeleteTaskLog` | `:126` 新增 | 父删**前**加 ``DELETE FROM tool_call WHERE task_id=?``；与 `:128` 的父删**同一** `s.write` 闭包/tx |
| `internal/memory/dao_tasklog.go:157` `PurgeTaskLogs` | `:160` 新增 | 父删**前**加 ``DELETE FROM tool_call``（全清＝这一发的闭包：父全删光后每一子都是子） |
| `internal/memory/retention.go:159` | 新增调用 | `deleteChildrenOfOlderParents`（`:232`）：``DELETE FROM tool_call WHERE task_id IN (SELECT id FROM task_log WHERE started_at < ?)``，在 `:163` 父删**之前**；`ToolCallDeleted` 折入（`:88` 注释说明），`:174/:179` 的原"按子窗老化"保留并 `+=` |
| `cmd/wisp/run.go:1108-1152` | 新增类型 | `taskLifecycleJournal`：只转发 `StartTaskLog`（`:1123`，`WithoutCancel`＋5s 界），其余四法惰性（`tool_call` 行仍归 bridge；终态关闭仍归本文件） |
| `cmd/wisp/run.go:1213` | 新增一行 | **只给 root loop** 装 journal（`loopOpts := opts` 的副本上装；`rt.loopOpt` 仍是 `Journal:nil` ⇒ `task.spawn` 的子环路写入面**一字不变**） |
| `cmd/wisp/run.go:1334-1338` | 改写序 | Wait 之后**不再无条件 StartTaskLog**：`rootRows.opened()` 为真时只 `Finish`；仅当循环没开行（控制词 `loop.go:349` 早退／开行写失败）才补开＋关（E1 形**回补**为 `fkc=0`） |
| `internal/memory/dao_tasklog.go:140-155`（件三） | 注释改成事实 | 原逐字 `// PurgeTaskLogs removes every task row (privacy: one-click clear). Forensic / // tool_call rows are governed by their own domain and retention window.` ⇒ 新注释写明：清空**含** `tool_call`；旧说法**没有契约出处**（现跑：`git grep -n PrivacyTaskLog -- docs/specs docs/PLAN.md` **0 命中**，rc=1；命中只在 `docs/evidence/**` 5 文件＋台账 9 行，都是本票讨论） |
| `internal/memory/migration_fkcheck_366_test.go:107-183` | 夹具改 raw 播种 | `purge-tasklogs`／`delete-tasklog` 改用 `t366RawParentDelete`（只删父，复刻 pre-373）；`retention` 直接手写"父已没、三列全 NULL 的 legacy 子"。**仍测原来那件事**：`3/2/1/1` 的 **legacy 孤儿库**照样被 3→4 迁移**记录**且不删行；变的只是"谁来写这份 legacy 状态"（同 `368-r1` 的 `raw 播种` 办法） |

- 新增用例：`internal/memory/cascade_373_test.go`（367 行，5 枚顶层／2 子例）。**⛔ 新增 `cmd/wisp` 测试文件**（写面明写"`cmd/wisp/**` 的其他文件"禁改）⇒ 件一的用例走**仓外 `-overlay`**（同 `373-a1` 普查的办法），见 §5。

## §3 三处各一发（改前 `fkc` → 改后 `0`）

夹具：`task-A`（2 子）＋`task-B`（1 子），生产 DAO 落库；两棵树 same probe（`probe/p373r1_probe_test.go`，before＝`base-clean` 导出树，after＝现树 `-overlay`）。

| 发起（生产入口） | 改前 `task_log`／`tool_call`／`fkc` | 改后 | 判据 |
|---|---|---|---|
| `PurgePrivacy(PrivacyTaskLog)`（`:177`→`:157`） | 0／**3**／**3**（`rows=2`） | 0／**0**／**0** | `fkc==0`＋子表 0 |
| `DeletePrivacyItem(PrivacyTaskLog,"task-A")`（→`:124`） | 1／**3**／**2** | 1／**1**／**0** | `fkc==0`；`task-B` 的子**原样** |
| `RunRetentionOnce`（父 31d／子 5d） | 0／**1**／**1**（`task_log_deleted=1`／`tool_call_deleted=0`） | 0／**0**／**0**（`1/1`） | `fkc==0` |

（附：普查 `P3` 的 pre-368 三列全 NULL 子亦量了：改前 `1/0/fkc=1` ⇒ 改后 `1/1/fkc=0`。）

## §4 父子窗口不同时那一发（★票面点名）

- ★**父 31d／子 5d**（同一库，`RunRetentionOnce` 单发）：**改前** `task_log_deleted=1`／`tool_call_deleted=0`／留孤儿 `1`／`fkc=1`（`child_ids=[task-old]`）⇒ **改后** `1/1`／`tool_call=0`／**`fkc=0`**。**子跟着父一起没**。
- 生产片段逐字：`retention.go:159` 的 cascade 在前、`:163` 父删在后，同一 `s.write`（`runRetention` 的 `:147` 闭包）。

## §5 先父后子那一发（E1：改前 `fkc=1` → 改后 `0`）

- 尺：`runTextTask` 真入口（`zz315` harness 的 golden SSE 脚本化**一次 `fs.read`**），另开一条 SQLite 连接在**任务运行期间**每 2ms 读 `count(task_log)`／`count(tool_call)`／`PRAGMA foreign_key_check`；`onRuntime` 钩子起表（store 已开）。
- **改前**（`base-clean` 导出树）：`poll_samples=12`、**`max_fkc_during_run=1`**、`child_without_parent_seen=true`（首次 `12:46:06.190`）⇒ 父行整段运行期**不存在**，子行已提交（**E1 形**）。
  - 红行逐字：`p373r1_e1_probe_test.go:109: E1: max foreign_key_check during the run = 1, want 0: at least one child was committed before the parent that names it`／`:113: E1: observed tool_call>0 with task_log=0 during the run (first at 12:46:06.190): the parent row landed after the child, so a crash in that window leaves an orphan`。
- **改后**（现树 `-overlay`）：`poll_samples=16`、**`max_fkc_during_run=0`**、`child_without_parent_seen=false`、运行后 `fkc=0`、`task_log_rows=1`／`child_rows=1`。**⇒ 改后 `fkc=0`。**
- 机制逐行：`run.go:1213` 把 journal 装到 root loop ⇒ 环路的 `j.startTask`（`internal/agent/loop.go:391`，在 `AdmitTask` 之后、**第一次模型调用之前**）写父行；任何 `tool_call` 行只能在其后。**用环路既有的写入点**，⛔ 造 id（id 仍由环路铸）。
- 语义边界：任务终态/成本/摘要仍由 `run.go` 在 Wait 后**一次**写（`:1334`），journal 的 `FinishTaskLog` 惰性；`started_at` 由"写完时刻"变成"任务开始时刻"（父行 state=`running`），这是 `journal.startTask` 的同形，也是 retention 按父龄删的正确起点。
- ⚠ **时间戳不足以做判据**：`started_at`／`decided_at` 都是**秒粒度**，改前改后两值都同秒（`1792816889`／`1792816918`）⇒ **判据只认运行期 `fkc` 观测**，⛔ 拿时间戳大小当证明。

## §6 正控（尺能读脏，也能读干净）

- **改前**（`base-clean`）：先 raw 删父弄脏到 `task_log=0/tool_call=3/fkc=3`，再**单 tx** 发 ``DELETE FROM tool_call WHERE task_id NOT IN (SELECT id FROM task_log)`` ⇒ `child_delete=<nil> parent_delete=<nil> commit=<nil>`、`fkc=0`、`rows=[]`、`tool_call=0`。
- **改后**（现树）：同一弄脏路径改用生产 `PurgeTaskLogs` ⇒ `n=0`（父已无）、`tool_call=0`、`task_log=0`、`fkc=0`。⇒ 3／2／1 这些读数不是尺坏。

## §7 "删的是谁的"（逐次，多集差集，已排序）

| 那一发 | 删掉的 `tool_call.task_id` | 留下的 |
|---|---|---|
| P1 `PurgeTaskLogs` | `[task-A, task-A, task-B]`（3 行，全部） | `[]` |
| P2 `DeleteTaskLog(task-A)` | `[task-A, task-A]`（2 行） | `[task-B]` |
| P4 retention（父 31d／子 5d） | `[task-old]`（1 行） | `[]`（父也删） |
| 用例里的 retention（父 31d／子 5d＋2d） | `[task-old, task-old]`（2 行） | `[task-live]`（活父的活子**不动**） |
| 正控（继承孤儿） | `[task-A, task-A, task-B]`（父已不在，仍被清） | `[]` |
| P2 的"父不存在"分支 `DeleteTaskLog(task-never-booked)` | **0 行**（`ErrNotFound` ⇒ 整笔回滚） | `[task-B, task-never-booked]` |

## §8 事务边界（逐处证明；`commit` 读数）

- ★**代码面**：三处的子删与父删都在**同一个** `s.write(ctx, func(ctx, tx *sql.Tx) error {...})` 闭包里（`dao_tasklog.go:125-134`／`:158-168`、`retention.go:147-190`）——该闭包就是一次 `BeginTx`…`Commit`，返回错误即回滚（`open.go:385` 的 queue 是唯一写路径）。
- ★**行为面（SQLite 触发器的真回滚）**：种 `BEFORE DELETE ON task_log … RAISE(ABORT)` ⇒ 子删**已跑**、父删被拒 ⇒ `DeleteTaskLog` 返错后 `tool_call` 仍 **3** 行／`task_log` 2／`fkc=0`；摘掉触发器再发同调用 ⇒ `tool_call=1`（**证"3"是回滚，不是级联没跑**）。retention 同法：被拒后 `tool_call=1`／`task_log=1`；摘触发器再跑 ⇒ `0/0`。
- ★**`commit` 读数**：B2 原始复现（同 tx、子先父后）⇒ `child_delete=<nil> parent_delete=<nil> commit=<nil>`、`fkc=0`；`B1`（父先）在普查里是 `787`（本腿未改 FK，未复跑）。
- ⚠ **诚实披露**：DAO 侧不向调用者返回 commit（queue 内部提交），所以上面的 commit 读数来自 B2 复现；DAO 侧的等价证据是"触发器拒绝父删 ⇒ 子删被回滚"这一发。另：`runRetention` 的 `res.ToolCallDeleted` 在闭包内**先赋值后提交**，事务失败时该计数不代表落库（**既有形状**，本腿未改；判定以库为准）。

## §9 新用例逐名 ＋ 为什么今天响（5 枚顶层／2 子例）

| 用例 | 今天为什么响（`base-clean` 树实跑红行，逐字） |
|---|---|
| `TestTicket373PurgeTaskLogsTakesTheChildren` | `:109 tool_call rows = 3 … want 0`／`:114 children removed = [], want [task-A task-A task-B]`／`:118 foreign_key_check after the purge = 3, want 0` |
| `TestTicket373DeleteTaskLogTakesOnlyThatTasksChildren` | `:141 …= [task-A task-A task-B], want [task-B] only`／`:145`／`:149 fkc = 2`／`:165`（not-found 回滚支） |
| `TestTicket373RetentionTakesChildrenOutsideTheirWindow` | `:198 tool_call deleted = 0, want 2`／`:203`／`:208`／`:211 fkc = 2`／`:222`（第二发） |
| `TestTicket373PositiveControlTheRulerReadsDirtyThenClean` | `:253 tool_call rows = 3 … want 0`／`:260 fkc = 3`（脏态先自证 `fkc=3`） |
| `TestTicket373ChildDeleteRollsBackWithTheParent`（2 子例） | `:328 tool_call rows after the unblocked delete = 3, want 1`／`:361 … unblocked pass = 1, want 0` |
| 件一：`p373r1_e1_probe_test.go`（仓外，`-overlay`） | `:109 E1: max foreign_key_check during the run = 1 …`／`:113 … child_without_parent_seen …` |

## §10 反向判据（摘掉新加那一跳 ⇒ `fkc` 回 `>0`）＋ 还原核

变异全在**仓外**（`mut/**`＋`-overlay`），工作树**零突变窗口**；每发只摘一跳：

| 变异 | 摘掉什么 | 读数（逐字红行） |
|---|---|---|
| M1 | `DeleteTaskLog` 的 ``DELETE FROM tool_call WHERE task_id=?`` | `cascade_373_test.go:141/145/149`（`[task-A task-A task-B]`／`[]`／`fkc=2`） |
| M2 | `PurgeTaskLogs` 的 ``DELETE FROM tool_call`` | `:109/114/118`（`3`／`[]`／`fkc=3`）＋ `:253/260`（正控 `3`／`fkc=3`） |
| M3 | retention 的 `deleteChildrenOfOlderParents` 调用 | `:198/203/208/211/222`（`deleted=0`／留 `task-old`／`fkc=2`） |
| M4 | `run.go:1213` 的 `loopOpts.Journal = rootRows` | `p373r1_e1_probe_test.go:109/113`（`max_fkc=1`、`child_without_parent_seen=true`） |

- ★**还原核（时机＝最后一发变异 M4 之后）**：M4 于 `12:46:35` 结束 ⇒ 之后现树重跑：`restore-e1-after-mutations.log`（E1 **PASS**，`max_fkc=0`）、`after-memory-full.log`（78 PASS／0 FAIL）、`after-cmdwisp-targeted.log`（10 PASS／0 FAIL）、`build-all.log`／`vet.log` 均 0 字节。文件 SHA-256（`sha256sum`）：`dao_tasklog.go 9a22f389…fdc4`／`retention.go a771c0f9…9479`／`cascade_373_test.go 4f28956b…24f0`／`migration_fkcheck_366_test.go 86e1984…3374`／`run.go 11406974…8263`。

## §11 四数（改前／改后分开，⛔ 相减）＋ 门禁 rc ＋ 格式甲乙两尺

| # | 门 | 改前（`base-clean`，`git archive HEAD`） | 改后（现树） |
|---|---|---|---|
| 1 | `go test -count=1 -v ./internal/memory/` | **73 PASS／0 FAIL／2 SKIP**，rc=0 | **78 PASS／0 FAIL／2 SKIP**，rc=0 |
| 2 | `cmd/wisp` 定向（`PATH="$PWD/third_party/sherpa-onnx:$PATH"`；`-run "TestRunTextTask|Test365|TestHostDispatchThroughTheAssembledBridge|TestComposedGateBlocksAWriteForTwoSeconds|TestMissingBlobFailsUnconfiguredNeverSilently"`） | **10 PASS／0 FAIL／0 SKIP**，rc=0 | **10 PASS／0 FAIL／0 SKIP**，rc=0 |
| 3 | `go build ./...` | —（改前树未跑） | **rc=0**（log 0 字节） |
| 4 | `go vet ./internal/memory/ ./cmd/wisp/` | — | **rc=0**（log 0 字节） |

- 门禁：`d22scan`（`cd tools/d22scan && go run . -root ../../`）＝**rc=0**，`clean - no D22 ban violations`（`internal/` 643 Go／`cmd/` 150 Go…）。
- **格式甲**：`gofmt -l <5 文件>` ⇒ **空**，rc=0。**格式乙**：`/d/work/base/gopath/bin/gofumpt.exe -l <5 文件>` ⇒ **空**，rc=0。（⛔ 用 `gofumpt -l .` 判 HEAD：会扫 `.scratch` 探针 ⇒ rc=2。）
- 门禁注意：`cmd/wisp` 整包会死 `MARK299` ⇒ **只跑定向 `-run`**；`go test` 一律 `GOFLAGS=`、⛔ `CGO_ENABLED=0`、只 `unset GOTMPDIR`。

## §12 残余与未做（⛔ 当已解决）

- **存量孤儿**：本机真实数据根（`%APPDATA%\wisp-dev`）`task_log=0/tool_call=0/fkc=0`（`373-a1` 读数）⇒ **回填 0 行、⛔ 编时间戳**；别处机器的存量孤儿**本腿无法量**。⚠ 本腿的"存量 0"是**空库意义上的 0**，⛔ 当"级联已生效"的证据。
- **子环路仍未开 `task_log` 行**（`rt.loopOpt.Journal` 保持 `nil`）：`task.spawn` 的子任务 `tool_call` 行因此仍无父——那是**另一个写面决定**（会给子任务新增 `task_log` 行、改变隐私页列表），⛔ 本票顺手改。已具名上报。
- **开行写失败**（DB 故障）时回落到"父后到"：`run.go:1335` 的补开在子行之后，FK 关着 ⇒ 仍能把 `fkc` 回补到 0（这也是改前对**每个**任务的行为）；真正的"父从未落"只剩"开行写入失败"这一条 store 故障路径。
- `PurgeTaskLogs` 语义：清空任务日志**连带** `tool_call`（件三裁决）；`PrivacyToolCall` 域仍在，`PurgeToolCalls` 仍可单独清。**⛔ 改 spec 文字**（本腿写面禁 `docs/specs/**`）。
- `schema.go:353 DROP TABLE task_log`（`ddlV4` 重建内）**不判为孤儿源**：同 tx、id 逐行搬、重建前 `taskLogStatesOutsideRoster` 先查、`migrationDSN` 强制 FK OFF（采纳 `373-a1` 判定）。
- `366` 夹具两枚改为 raw 播种（§2）：**仍测**"legacy 孤儿库被迁移记录且不删行"，**不再测**"今天的生产路径会造孤儿"（后者已由 §3 的 0 读数从正面钉住）。

## §13 提交

- 本件 ＋ 4 个 `.go`（`cmd/wisp/run.go`、`internal/memory/dao_tasklog.go`、`internal/memory/retention.go`、`internal/memory/migration_fkcheck_366_test.go`）＋ 1 个新测试（`internal/memory/cascade_373_test.go`）**同一笔**；`git add`／`commit` 均**逐个文件名字面点名**（⛔ glob／`-A`／`.`）；主题逐字以 `票 373 落地腿 373-r1` 开头；**⛔ push**。
- `git diff --numstat`：`cmd/wisp/run.go 98/4`、`internal/memory/dao_tasklog.go 32/3`、`internal/memory/retention.go 40/3`、`internal/memory/migration_fkcheck_366_test.go 42/36`；新文件 `cascade_373_test.go` 367 行。
