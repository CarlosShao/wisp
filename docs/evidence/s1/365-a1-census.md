# 票 365 普查腿 365-a1 —— 两处 `taskLogState` 的只读普查（AC#0）

- 腿：**365-a1（只读）**｜HEAD `ef3ec49bd`（`dev`）｜读数 2026-10-23 14:3x–15:0x +08
- 写面：本件 ＋ 仓外 `D:\tmp\wisp\20261023\365a1\**`｜**仓内 `.go` 改动 0 枚**（§7）
- 票面：`.scratch/wisp/issues/365-two-taskLogState-functions-disagree.md`（55 行，已通读）；必读清单 `.scratch/wisp/probes/orch/20261022/txt/ruler-scope-checklist.md`（73 行，已通读）
- 每条读数都写了"哪把尺、量的是哪个对象"（清单五问）

## 1. 两处逐条复核（对 HEAD blob 现读）

尺：`git show HEAD:<path>` 导出到仓外再读；哈希对账 `git rev-parse HEAD:<path>`（SHA-1、40 位）vs `git hash-object <worktree>`，两枚相等 ⇒ 我读的就是 HEAD 的字节。

| # | 票面断言（逐字） | 现读（HEAD blob） | 漂？ |
|---|---|---|---|
| F1 | `internal/agent/loop.go:1048` `func taskLogState(s Status) string` | `:1048` 逐字相同 | **没漂** |
| F2 | 它有 `case StatusStuck: return "stuck"` | `:1052 case StatusStuck:`／`:1053 return "stuck"` | **没漂**（函数在 :1048、case 在 :1052-1053） |
| F3 | `cmd/wisp/run.go:1361` `func taskLogState(s agent.Status) string` | `:1361` 逐字相同 | **没漂** |
| F4 | 它**没有** `StatusStuck` 支、`:1370` 默认 `return "error"` | 只有三支：`:1363` Completed／`:1365` Cancelled／`:1367` Control；`:1369 default:`／`:1370 return "error"` | **没漂** |

哈希：`loop.go`＝`31dec50a83971862b15fde5d1a15ee2f8f889c16`、`run.go`＝`9fc8fbe97ed730d458cf834d584cea6d624d5f5e`，HEAD 与工作树同值 ⇒ **票面四条行号全部仍有效（票面前提新鲜）**。

## 2. `Status` 对账表（逐枚；六枚常量一颗不漏）

尺：常量集合＝`internal/agent/loop.go:40-55` 的 `const` 块，并用 `:58-65 Valid()` 反向核（`Valid` 列的六枚＝常量块六枚；`grep "Status[A-Za-z]* Status ="` 亦只六行）⇒ **没有第七枚**。返回词取自 HEAD blob 的 `case` 体。

| `Status` 常量 | 底层串 | `internal/agent`（loop.go:1048） | `cmd/wisp`（run.go:1361） | 同？ |
|---|---|---|---|---|
| `StatusCompleted` | `completed` | `done`（:1050-1051） | `done`（:1363-1364） | ✅ |
| `StatusStuck` | `stuck` | **`stuck`**（:1052-1053） | **`error`**（落 :1369-1370 默认） | ❌ **分裂** |
| `StatusCancelled` | `cancelled` | `cancelled`（:1054-1055） | `cancelled`（:1365-1366） | ✅ |
| `StatusFailed` | `failed` | `error`（:1058-1059） | `error`（落默认） | ✅ 同词（路径不同） |
| `StatusTruncated` | `truncated` | **`interrupted`**（:1056-1057） | **`error`**（落默认） | ❌ **第二处分裂** |
| `StatusControl` | `control` | `done`（落 :1060-1061 默认；★`loop.go:52-53` 说该支不开 task_log 行） | `done`（:1367-1368 显式支） | ✅ 同词（agent 侧不可达） |
| （非常量）未知值 | — | `done`（:1061） | `error`（:1370） | ❌ 两处默认支相反 |

⇒ **票面只点了 `StatusStuck`；实测差集＝`{StatusStuck, StatusTruncated}`，外加一对相反的默认支。**

## 3. 真读数（本格核心：两条路径各写什么，⛔ 读码推断）

台件（全在仓外、经 `-overlay` 入场，仓内 0 字节 `.go`）：`probe-agent_test.go`／`probe-cli_test.go`／`overlay.json`；日志 `agent-probe.log`／`cli-probe.log`／`cli-census.log`。
- **路径甲＝`internal/agent`**：真 `Loop.Run`，`Options.Journal = memory.Store`（`loop.go:159-160` 是写库开关、`:369` 建 journal）⇒ `finish`→`taskJournal.finishTask`（`journal.go:137`）→`memory.Store.FinishTaskLog`→sqlite；golden `repeat-echo` 驱动真 C22 刹车（8 次同参）。
- **路径乙＝`cmd/wisp`**：真 `runTextTask`（`wisp run` 与常驻腿共用的组合根）⇒ `agentRuntime.execute`（`run.go:1067`）`:1196 state := taskLogState(res.Status)` → `:1205 StartTaskLog`／`:1212 FinishTaskLog` → sqlite；同一枚 golden、同一个刹车。
- 读回尺：`memory.Store.TaskLogByID`（不是 raw SQL、也不是打印推断）。

| 输入 | 甲写进库的 `state` | 乙写进库的 `state` | 判 |
|---|---|---|---|
| `StatusStuck` | **`stuck`** | **`error`** | ❌ 分裂复现 |
| `StatusTruncated` | **`interrupted`** | **`error`** | ❌ 第二处也复现 |
| `StatusCompleted` | **`done`** | **`done`** | ✅ 正控（§4） |

逐字（探针 stdout，原样）：

```
PROBE365A1-AGENT | label=stuck | fixture=repeat-echo | Status="stuck" | task=d5315ee3-8616-46ca-92c4-a5ef424983c1 | task_log.state="stuck" | ended=true | error_class=""
PROBE365A1-CLI | label=stuck | fixture=repeat-echo | exit=3 | task=013a5dfb-acca-4767-9fce-f7dc693942db | task_log.state="error" | ended=true | error_class=""
```

★**同一发任务，屏幕与库各说一套**（`wisp run` 自己的状态行，逐字）：

```
wisp run: 任务 013a5dfb-acca-4767-9fce-f7dc693942db 结束（stuck，8 轮，7 次工具调用，成本 0 CNY（未计价））
```

⇒ 屏幕上 `stuck`、库里 `error`。★这不是纸上推演，是**同一发任务一次跑出来**的两行字。

每发任务只 1 行 `task_log`（`PROBE365A1-CENSUS | … | rows=1`，五库五发逐发核过）⇒ 乙路径不是"先被甲写成 `stuck`、再被乙覆盖"：`cmd/wisp` 的 Loop 是 `Journal: nil`（§6），乙是**唯一**写者（stderr 也没有 `task_log open failed` 那行主键冲突迹象）。

## 4. 正控（硬闸：证明尺不是全红）

同一套台件、同一把读回尺、同一个 `StatusCompleted`（golden `text-reply`，纯文本答复）：

```
PROBE365A1-AGENT | label=completed | fixture=text-reply | Status="completed" | task=c9a39f5c-0edc-4144-a6a5-bbda3888a396 | task_log.state="done" | ended=true | error_class=""
PROBE365A1-CLI | label=completed | fixture=text-reply | exit=0 | task=3c9843ab-1a0d-49c0-a08b-760202a3eb54 | task_log.state="done" | ended=true | error_class=""
```

⇒ **两处同词 `done`**、退出码 0 ⇒ 尺能读出"同"，分裂只落在 §2 那两枚上。

## 5. 票 196 的名册会不会挡住这件事

尺三把、对象不同（⛔ 混报）：Go 侧名册＝`memory.TaskStateRoster()`；库侧约束＝`sqlite_master` 里 `task_log` 的 DDL；库的实际拒绝＝raw SQL 绕过 DAO 的 INSERT。

- 名册逐字 12 枚（`len=12`）：`["running" "done" "succeeded" "cancelled" "interrupted" "stuck" "error" "Thinking" "Settling" "Muted" "Stuck" "Error"]`
- `stuck`：`ValidTaskState("stuck")=true`；`error`：`ValidTaskState("error")=true` ⇒ **两个词都在册**。
- 两口都收（实测）：DAO `StartTaskLog(state="stuck") err=<nil>`／`(state="error") err=<nil>`；raw SQL `INSERT state="stuck" err=<nil>`／`state="error" err=<nil>`。
- ⇒ **分裂不会被库拒绝**：两条路径写出的两种词都合法（§3 那两行就是真进了库的）。票 196 的乙只保证"词在册"，⛔ 保证"两处同词"。
- **若将来统一成名册外的词**（最容易落的形状＝两个函数的 `Status` 底层串 `failed`／`truncated`／`completed`，或自造词）⇒ 两口都拒，逐字：
  - 甲（DAO，`models.go:194-197`）：`memory: invalid task_log.state "failed" (want one of: running|done|succeeded|cancelled|interrupted|stuck|error|Thinking|Settling|Muted|Stuck|Error)`
  - 乙（CHECK，raw SQL）：`constraint failed: CHECK constraint failed: state IN ('running', 'done', 'succeeded', 'cancelled', 'interrupted', 'stuck', 'error', 'Thinking', 'Settling', 'Muted', 'Stuck', 'Error') (275)`
  - DDL 行逐字：`state          TEXT NOT NULL CHECK (state IN ('running', 'done', 'succeeded', 'cancelled', 'interrupted', 'stuck', 'error', 'Thinking', 'Settling', 'Muted', 'Stuck', 'Error')),`
- 另读数：大写 D43 拼写 `Stuck`／`Error` **也在册**（`ValidTaskState=true`，raw INSERT `<nil>`）⇒ "统一成 D43 名"库不会拒，但那已是**另一个词**（小写旧词的读者按小写读）。
- ⇒ 交 AC#1：甲/丙不产生新词；**乙若把词改成名册外 ⇒ 必须先改名册（`internal/memory/taskstate.go` 的 `taskStateRoster`）并升 schema／迁移，否则写不进库**。

## 6. 归因（为什么分裂还活着）

1. **发布组合里只有一个 Loop 装配点**：`cmd/wisp/run.go:1072-1100`，且 `Journal: nil`（`:1080`；同义注释在 `internal/tools/bridge.go:1151`——tool_call 行由 bridge 写，Loop 必须 nil 否则双写）⇒ **`internal/agent` 那处 `taskLogState` 在已发布二进制里没有调用者**，它被 `Options.Journal != nil` 这道门关着（`loop.go:159-160`），而唯一装配点把门关死了。我这次拿真读数，正是**把 `memory.Store` 挂成 `Journal`** 才让那条链真跑起来（§3 路径甲）。
2. `task.spawn` 的子环路继承父装配（`internal/tools/subagent_197.go:320 opt := base` → `:347 agent.New(opt)`）⇒ 子也没有 Journal。
3. ⇒ **今天库里所有 `StatusStuck` 行都是 `error`**（cmd 侧写的）；`stuck`／`interrupted` 在当前组合下无生产者——与票 196 的"Census 四旧词"里没有它们吻合。
4. **票 196 的甲（权威映射）也没有生产者**：`TaskStateFor`／`TaskStateFromD43`／`TaskLog.D43State` 在 `cmd/**`＋`internal/**` 的非测试代码里 **0 个调用者**（现跑 grep 只命中 `internal/memory/taskstate.go` 自身）⇒ 名册只当了守门的。**这就是"源头没除根"的准确形状。**
5. ★**第三处同形映射（具名上报）**：`internal/tools/subagent_197.go` 有**两份一模一样**的 `switch res.Status`（`:451-457` `finalize`、`:480-486` `answer`），把 `StatusStuck` 映射成 **`Settling`**（`subagentStateSettled`，`:108`）——同一枚 `StatusStuck` 的**第三种答案**。★但它量的对象不同：名册 `TaskOutput.State`（D43 名、判据 `statemachine.Valid`、票 188；`task.go:162-166` 明说"memory 侧那套旧词表属票 196，本维不接它"）⇒ ⛔ 与 `task_log.state` 同一把尺。列出来只为让 AC#1 裁形时别以为"只有两处"。

## 7. 门禁 rc／提交／0 枚 `.go`／正控尺

命令一律 `cmd > log 2>&1; rc=$?`（⛔ 管道取 rc），日志在 `D:\tmp\wisp\20261023\365a1\gate-*.log`：

| 门 | rc |
|---|---|
| `GOFLAGS= go build ./...` | **0** |
| `GOFLAGS= go vet ./internal/agent ./cmd/wisp` | **0** |
| `GOFLAGS= go test ./internal/agent/ ./internal/memory/ -count=1` | **0** |
| `GOFLAGS= go test ./internal/memory/ -run TestTaskState196 -count=1` | **0** |
| `sh scripts/d22scan.sh` | **0** |
| `sh scripts/check-path-length-budget.sh` | **0** |

- **仓内 `.go` 改动 0 枚**：`git status --porcelain -- '*.go'` 共 77 行、**全部是 `??`**（`?? .scratch/…` 的历史探针残留、非本腿）；tracked 的 `M/D/A` **0 行**。两发探针都在 `D:\tmp\wisp\20261023\365a1\`，靠 `-overlay` 入场。
- ★**已知环境红（⛔ 记本票头上）**：`go test ./cmd/wisp/` **不加** `third_party/sherpa-onnx` 到 PATH ⇒ 整包 `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND、测试进程一次没起、0.133s）。病因：`sherpa-onnx-go-windows` 走 cgo（`CGO_ENABLED=1`）。本件 cmd/wisp 的读数**都是加了 PATH 之后**跑的：`export PATH="…/third_party/sherpa-onnx:$PATH"`，不改仓内一个字节。
- **正控尺**：§4（`StatusCompleted` 两处同词 `done`，两发都是真库读数）⇒ 不是全红。
- **提交**：`b04cbeaba`（本件入库，117 行、只这一枚文件、显式 pathspec）；本行随后由一枚补记提交追加（`git log -1 -- docs/evidence/s1/365-a1-census.md`）。⛔ push。
