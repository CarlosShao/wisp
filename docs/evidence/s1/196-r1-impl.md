# 票 196 落地腿 `196-r1` —— 〔乙〕＝甲 ＋ CHECK 约束

- 形：**〔乙〕**（甲 ＋ 数据层 CHECK）。⛔ 甲、⛔ 丙。
- 锚 `8d3ad434`（编排者拍形那次提交）。写面：`internal/memory/**` ＋ 两包测试 ＋ 本件；`internal/statemachine/**` 只读引用。
- 判据框：本票无（资料型）；下面的读数都是可逐字重跑的命令回声。

## 0. 机主批准（逐字）

> 「好的，这三件事我都允了，同意，你可以继续了」

## 1. 票面四个事实逐条现跑复核 —— **3 枚行号已腐坏**

| # | 票面尺 | 现跑读数 | 漂没漂 |
|---|---|---|---|
| F1 | `sed -n '8p;11,31p' internal/statemachine/states.go` | `:8` ＝ `// The 20 states, names exactly as in D43 / SPEC-08...`；`:11-31` ＝ 20 枚 `StateXxx State = "..."` 恰 21 行（含收尾 `)`） | **没漂** |
| F2 | `sed -n '3055p' docs/PLAN.md` | `:3055` 是 **D42 可靠性表的一行**（WebView2 缺失）；**D43 标题真在 `:3061`**（`### D43 — 状态机权威转移表`），转移表 `:3065-3100` | **漂 +6** |
| F3 | `loop.go:983-998`；`run.go:651-670` | 两处都是**别的东西**：`loop.go:983-998` 是 `func (l *Loop) finish(...)` 的签名与赋值；`run.go:651-670` 是 `tools.Options` 的注释块。真位置见 §2 | **漂** |
| F4 | `models.go:57` 邻域 | `:57` ＝ `// complete.`。**CHECK 的问题属于 `internal/memory/schema.go:53-66`**：`state TEXT NOT NULL` 在 `:57`，DDL 内**无 CHECK** | **漂**（文件也换了） |

★结论：票面「今天的四词＝`done`/`cancelled`/`running`/`succeeded`」**只对一半**——见 §2。

## 2. ★本票核心证据：今天谁能把 `task_log.state` 写成非法值（一发真读数）

改任何一行**之前**，在 `internal/memory` 里放一次性探针跑 `-run TestZZ196ProbeHoleIsOpen`（跑完即删，输出留在 `D:\tmp\wisp\20261023\196r1\logs\05-hole-open-before.txt`）：

```
READING door1 StartTaskLog(state="definitely-not-a-task-state") err = <nil>
READING door1 readback state = "definitely-not-a-task-state" (read err = <nil>)
READING door2 FinishTaskLog(state="definitely-not-a-task-state", errorClass="") err = <nil>
READING door2 readback state = "definitely-not-a-task-state"
READING door3 raw INSERT state="definitely-not-a-task-state" err = <nil>
READING rows now holding "definitely-not-a-task-state" = 3
READING schema_version = 3
```

★**三扇门全开**，而且门 2 比票面想的更宽：`FinishTaskLog` 原来只在 `errorClass != ""` 时才校验，**干净退出那条路一次都没校验过**（`internal/memory/dao_tasklog.go`）。
★列自己怎么写的（同一次读数的 `sqlite_master` 回声）：`state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）` —— **只有 NOT NULL，没有 CHECK**。

★**写面普查（非测试代码）**：真正能写进这一列的只有六个词，⛔ 四个：
`running`（`internal/agent/journal.go:131`）、`done`/`stuck`/`cancelled`/`interrupted`/`error`（`internal/agent/loop.go:1046-1063` 的 `taskLogState`）、`done`/`cancelled`/`error`（`cmd/wisp/run.go:1361-1372` 的同名函数）。
★`succeeded` **今天已没有生产者**，只剩在 `internal/memory/dao_test.go`、`privacy_test.go` 和**旧版本写下的库行**里 ⇒ 名册必须留它，否则老库打不开。

## 3. 〔甲〕Go 侧唯一权威映射：`internal/memory/taskstate.go`（新文件，179 行）

- `type TaskState string` ＋ 7 枚常量：`running`/`done`/`succeeded`/`cancelled`/`interrupted`/`stuck`/`error`。
- `var taskStateRoster = []TaskState{...}` —— **12 枚、定序**：上面 7 枚 ＋ **D43 里任务行真会取的 5 枚** `Thinking`/`Settling`/`Muted`/`Stuck`/`Error`。
  ★D43 那 5 枚不是我发明的：`internal/tools/subagent_197.go:107-110` 已经为"任务记录"选过 `Thinking`/`Settling`/`Muted`/`Error`，我照抄这份读数再加 `Stuck`（因为旧词表本来就有 `stuck`）。
  ★**⛔ 不是 D43 全 20 枚**：`Sleeping`/`Listening` 这类球态进不了任务行（有专门一条腿钉住）。
- `taskStateToD43`（A→B，全函数）：`running→Thinking`、`done→Settling`、`succeeded→Settling`、`cancelled→Muted`、`interrupted→Muted`、`stuck→Stuck`、`error→Error`。
- `taskStateFromD43`（B→A，由名册序推导，多对一时取**名册里靠前**的那个 ⇒ 往返确定，⛔ 依赖 map 序）。
- 对外四枚：`TaskStateRoster()`（还一份拷贝，调用者改不动词表）、`ValidTaskState()`、`TaskStateFor()`、`TaskStateFromD43()`；读侧门 `TaskLog.D43State()`。
- **写走它**：`validateTaskLog`（`internal/memory/models.go`）现在拒非名册词；`FinishTaskLog` 改成**每次都校验**（补上 §2 门 2）。
- **读走它**：任何读面要 D43 名就问 `TaskStateFor`/`D43State`，⛔ 自己写 switch。

## 4. 〔乙〕CHECK 逐字 ＋ 它的显式名册

`internal/memory/schema.go` 新增 **schema v4**：`applyV4` 用 SQLite 官方 12 步重建 `task_log`（建新表→搬行→删旧表→改名→重建索引，5 条语句，全在链条原有的那一枚事务里）。

`buildDDLV4()` 生成的那一列（逐字）：

```
  state          TEXT NOT NULL CHECK (state IN ('running', 'done', 'succeeded', 'cancelled', 'interrupted', 'stuck', 'error', 'Thinking', 'Settling', 'Muted', 'Stuck', 'Error')),
```

★★**名册是具名的、⛔ 藏在 SQL 里**：值表由 `strings.Join(quoted, ", ")` 从 Go 的 `taskStateRoster` **生成**，SQL 里没有第二份手抄。改词表只有一处可改。
★为什么重建而⛔ 改 `ddlV1`：SQLite 没有 `ALTER TABLE ADD CONSTRAINT`，而 `ddlV1` 是对 `docs/specs/SPEC-02-data-storage.md` **逐字冻结**的（`internal/memory/schema_test.go` 的内省腿钉住 15 条语句）。走 v4 步 ⇒ **新库与老库走同一条 3→4**，不会出现"新库 CREATE、老库 ALTER"的分叉（这正是 v3 立下的规矩）。
`applyV4` 先**点名**再动手：库里若有名册外的值，直接报"holds N value(s) outside the v4 roster (…)"并回滚，⛔ 让操作者面对一句不带行身份的 `CHECK constraint failed`。

## 5. 迁移路径 ＋ 老库能打开的实测读数（⛔ 只新库）

老库夹具 `seedLegacyV3` **用生产的 `ddlV1`/`ddlV2`/`ddlV3` 常量**造一个真 v3 库（`schema_version=3`），塞进四词行，再 `Open()`：

```
--- PASS: TestTaskState196LegacyV3DatabaseOpensAndReads
  schema_version = 4                      （3→4 落了）
  rows = 4                                （四词行一枚不少、没重复）
  四词各一读回后在名册里都有 D43 读音
  迁移后的库上 raw INSERT 非法值 ⇒ 被拒  （CHECK 真的进了老库，不只是新库）
```

反向那条腿（`TestTaskState196RefusesToMigrateAnOutOfRosterDatabase`）：塞一枚名册外词 ⇒ `Open()` **失败并点名该值**，且事务回滚到 `schema_version=3`、`rows=2`（没搬、没删）——即"老库打不开"只在**它本来就坏了**的时候发生，且说得出坏在哪。

## 6. CHECK 两向读数

①**写非法值 ⇒ 被拒**（`TestTaskState196CheckRejectsIllegalRawSQL`，原样树，`t.Logf` 逐字）：

```
constraint failed: CHECK constraint failed: state IN ('running', 'done', 'succeeded', 'cancelled', 'interrupted', 'stuck', 'error', 'Thinking', 'Settling', 'Muted', 'Stuck', 'Error') (275)
```

UPDATE 方向同样被拒（腿里另有一发）。

②**写合法值（四词各一）⇒ 通过**（`TestTaskState196CheckAcceptsTheFourLegacyWords`）：

```
legacy word "done"      accepted by the CHECK
legacy word "cancelled" accepted by the CHECK
legacy word "running"   accepted by the CHECK
legacy word "succeeded" accepted by the CHECK
```

★另有 `TestTaskState196CheckAcceptsTheWholeRoster` 把 12 枚全量过一遍 raw SQL ⇒ 甲、乙两张嘴⛔ 会不一致。

## 7. 反向判据：三读数 ＋ 还原核

hash 口径先问清楚：**`sha256sum` ＝ SHA-256 ＝ 64 位十六进制**（⛔ `git hash-object` 的 SHA-1/40 位，两者⛔ 混比）。

| 变异 | 摘掉哪一跳 | 读数 | 逐名报红 |
|---|---|---|---|
| M1 | 让 CHECK 不再约束（`CHECK (state IS NOT NULL OR state IN (…)`，恒真） | **红** rc=1 | `TestTaskState196CheckRejectsIllegalRawSQL`（`taskstate_196_test.go:196`）、`TestTaskState196LegacyV3DatabaseOpensAndReads`（`:369`） |
| M2 | 关掉 DAO 名册校验（`if false && !ValidTaskState(...)`） | **红** rc=1 | `TestTaskState196DAORejectsIllegalOnStart`（`:147`）、`TestTaskState196DAORejectsIllegalOnFinish`（`:169`） |
| M3 | 把 `succeeded` 从名册摘掉 | **红** rc=1 | 7 枚：`TestTaskLogLifecycleAndErrorClassValidation`、`TestPrivacyOpsAllDomains`、`TestRetentionBoundaries`、`TestTaskState196RosterIsExplicit`、`…MappingBothDirections`、`…CheckAcceptsTheFourLegacyWords`、`…LegacyV3DatabaseOpensAndReads` |

★M2 那发还顺带证明了**纵深**：关掉甲之后写入**仍被拒**，因为乙先开口；红的只是"拒绝语没点名"。
★M3 那发证明了名册的旧词半边**吃重**：摘掉 `succeeded`，老库迁移直接按名拒绝（`task_log.state holds 1 value(s) outside the v4 roster (succeeded)`）。

- **还原绿**：`sha256sum -c` 三枚（`schema.go`/`models.go`/`taskstate.go`）全 `OK`；随后 `go test ./internal/memory/ ./internal/statemachine/` ⇒ `ok 33.301s` / `ok 0.317s`，rc=0。
- **原样绿**：同上一发即原样树读数（还原后跑的就是原样树）。

## 8. 门禁 rc（`unset GOTMPDIR` 后逐条跑）

| 门 | rc | 说明 |
|---|---|---|
| `go build ./...` | **0** | |
| `go vet ./...` | **0** | |
| `sh scripts/d22scan.sh` | **0** | 本件写出后复跑；先前那次 rc=1 的唯一 finding 就是本件还没落盘（`[phantom-citation] … 196-r1-impl.md`） |
| `bash scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`，17612 条 tracked，over-budget 61 全在名册 |
| `gofumpt -l internal cmd tools scripts frontend` | **1** | ★**⛔ 我造成的**：3 枚**写面外**文件在工作树里是 CRLF（`internal/agent/approval/pending_read.go` CR=131、`internal/agent/tools.go` CR=225、`cmd/wisp/models.go` CR=389，各自 CR 数＝各自行数）。`git status` 对这三枚**无输出**（LF blob 与工作树经 eol 过滤后一致）⇒ 是别人/环境写出的 CRLF，⛔ 我的写面，我没碰。我的 7 枚文件 CR **全 0**、`gofumpt -l` 单点全清。 |
| `go test ./cmd/wisp/` | **1（预存）** | `exit status 0xc0000135` ＝ STATUS_DLL_NOT_FOUND，**测试进程一次没起**（0.12–0.19s）。★已用 `git archive HEAD` 导到仓外的**原样 HEAD 树**复跑：**同样 0xc0000135** ⇒ 与本次改动无关。 |
| `go test ./internal/{memory,statemachine,agent,session,tools}/` | **0** | memory `ok`、statemachine `ok`、agent `ok 7.225s`、session `ok 6.403s`、tools `ok 30.095s` |

## 9. 未做／欠账（点名，⛔ 偷偷做掉）

1. ★**`cmd/wisp/run.go:1361-1372` 有第二个 `taskLogState`，与 `internal/agent/loop.go:1046-1063` 的映射⛔ 一致**（同一枚 `StatusStuck`，一边写 `stuck`、一边写 `error`）。这就是票面说的"可观察分裂"的**真身**，但 `cmd/**` **在我的写面之外** ⇒ 我没动它。甲的插口（DAO）已经把两条生产路径都盖住，**收敛这两个函数本身仍是一件待开票的事**。
2. ★`docs/specs/SPEC-02-data-storage.md` §3 的 `task_log` 块**仍然没有 CHECK**；`docs/PLAN.md` D35 同理。机主批准的是"触 D35"，⛔ 等于"改契约文字"，而契约文字须另走批准 ⇒ 这是一条**spec-note 欠账**，已在 `internal/memory/schema.go` 头部注释里写明来源与本票形状。
3. ★名册里 D43 的 5 枚是**为将来迁移预留**，今天没有生产者写它们；`error` 与 `Error`、`stuck` 与 `Stuck` 在同一列并存是**刻意的**（两套词各自保留），收敛要等第 1 条那件事。
4. ★`internal/tools`（票 188 的地盘）**零改动**：`task_state_188_test.go` 的 `legacyTaskStateWords` 名册仍按 A394 归 196，但 188 那枚新字段不归我。
