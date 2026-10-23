# 票 365 落地腿 `365-r1`：两处 `taskLogState` 收成一处权威映射（〔乙〕形）

**腿号** `365-r1`（写腿，只 commit 不 push）｜形＝〔乙〕两处都走票 196 的 `memory.TaskState`
**锚** 落笔时 HEAD `dd16b005c`；腿中途编排者并发落了 `a8155492f`、`771bf451e` 两笔（`git diff --stat a8155492f..771bf451e -- internal/agent internal/memory cmd/wisp` 为空 ⇒ 与本笔无交集）
**环境** Go 1.27.1／`GOFLAGS=`／只 `unset GOTMPDIR`／gofumpt `/d/work/base/gopath/bin/gofumpt.exe`（v0.12.0）
**尺的射程**：改前四数取**仓外导出树** `D:\tmp\wisp\20261023\365r1\treeBefore`（`git archive HEAD`）；`cmd/wisp` 的改前取**工作树＋`-overlay`**（见 §6 为何不用导出树）；rc 全取壳级 `$?`、不走管道。

## §0 一句话裁决

两处裸 `switch` **都收掉**：一根梯子搬进 `agent.TaskLogState`（返回 `memory.TaskState` 常量），`cmd/wisp` 只留一行委托 ⇒ 同一枚 `Status` 两条路径写同一个词；四读数全给、正控三枚未动、反向判据两发转红并逐名报红行、还原后全绿。

## §1 改前三张静态表＋一张实跑表（我现跑，⛔ 抄 `365-a1`）

### 表 A：两处梯子逐字（改前，HEAD 内容）

| 位置 | `case` 集合（顺序照源码） | 默认支 |
|---|---|---|
| `internal/agent/loop.go:1048` `taskLogState(s Status) string` | `Completed→"done"`、`Stuck→"stuck"`、`Cancelled→"cancelled"`、`Truncated→"interrupted"`、`Failed→"error"` | `"done"`（无 `Control` 支） |
| `cmd/wisp/run.go:1361` `taskLogState(s agent.Status) string` | `Completed→"done"`、`Cancelled→"cancelled"`、`Control→"done"` | `"error"`（无 `Stuck`/`Truncated`/`Failed` 支） |

### 表 B：六个 `Status` ＋ 未知值逐枚比对（差集 3 枚）

| `Status` | 底层串 | `internal/agent` | `cmd/wisp` | 同词？ |
|---|---|---|---|---|
| `Completed` | `completed` | `done` | `done` | ✓ |
| `Stuck` | `stuck` | `stuck` | **`error`**（默认支） | **✗** |
| `Cancelled` | `cancelled` | `cancelled` | `cancelled` | ✓ |
| `Failed` | `failed` | `error` | `error`（默认支） | ✓ |
| `Truncated` | `truncated` | `interrupted` | **`error`**（默认支） | **✗** |
| `Control` | `control` | `done`（默认支） | `done` | ✓ |
| 未知值 `Status("")` | — | **`done`**（默认支） | `error`（默认支） | **✗** |

⇒ 票面点的一枚（`Stuck`）＋ 另两枚（`Truncated`、未知值）。与 `365-a1` 一致。

### 表 C：票 196 名册 12 枚 × `Status` 底层串在册否

名册（`internal/memory/taskstate.go:72` `taskStateRoster`）＝`running`·`done`·`succeeded`·`cancelled`·`interrupted`·`stuck`·`error`·`Thinking`·`Settling`·`Muted`·`Stuck`·`Error`（12 枚）。
`completed`✗ `stuck`✓ `cancelled`✓ `failed`✗ `truncated`✗ `control`✗ ⇒ **六个底层串里 3 个不在册** ⇒ 屏幕词与库词只可能在 `stuck`／`cancelled` 上逐字相同。

### 表 D：改前**一发真读数**（缺陷本体，用 `-overlay` 把 `cmd/wisp/run.go` 换回 HEAD 版实跑）

```text
屏幕逐字: wisp run: 任务 b2145bdd-e5fd-4554-bedf-10cab80fcc8f 结束（stuck，1 轮，0 次工具调用，成本 0 CNY（未计价））
库逐字:   task_log.state = "error"
```

同一发任务、同一行 `task_log`：屏幕说 `stuck`、库说 `error`。**每发 1 行**（`TaskLogByID` 单条）。

## §2 改了哪两处（逐行）＋ 权威入口在哪

| # | 位置（改后行号） | 改了什么 |
|---|---|---|
| 1 | `internal/agent/loop.go:1069` `func TaskLogState(s Status) memory.TaskState` | 原 `taskLogState` 的裸 `switch` **换成这一根**：`Completed,Control→TaskStateDone`、`Stuck→TaskStateStuck`、`Cancelled→TaskStateCancelled`、`Truncated→TaskStateInterrupted`、`default→TaskStateError`；**返回值全是 `memory` 常量、零字符串字面量** |
| 2 | `internal/agent/loop.go:1014` | 调用点 `taskLogState(status)` → `string(TaskLogState(status))` |
| 3 | `internal/agent/loop.go:17` | 新增 import `internal/memory`（原文件不 import，`journal.go`/`tools.go` 已 import） |
| 4 | `cmd/wisp/run.go:1366-1368` | 原裸 `switch` **整个删掉**，函数体只剩 `return string(agent.TaskLogState(s))`（`run.go:1196` 调用点未动） |

**权威入口**＝`internal/memory/taskstate.go`：`TaskState`（`:54`）＋ `taskStateRoster`（`:72`，12 枚）＋ `ValidTaskState`（`:136`，DAO 与 SQL CHECK 共用同一集合）。

**为什么梯子必须在 `internal/agent` 而不在 `internal/memory`（派单允许动那层、我核完不能动）**：`agent` 已 import `memory`（`journal.go:8`、`tools.go:11`），`memory` 对 `agent` 的 import 数＝**0** ⇒ 让 `memory` 收 `agent.Status` 会造 **import cycle**。故 `internal/memory` **一个字节未动**，`memory` 仍拥有词表、`agent` 只拥有「把这一枚枚举投影到词表」的那一跳——这也正是铁律要的「一个值域只有一个权威转换点」。

## §3 四条必须给的读数（★每条都是库里的真读数，不是函数返回值）

| 项 | `internal/agent` 路径写进库的 `state` | `cmd/wisp` 路径写进库的 `state` | 改前→改后 |
|---|---|---|---|
| **① `StatusStuck`** | **`"stuck"`**（真 Journal＋真 SQLite，fixture `repeat-echo`） | **`"stuck"`**（真 `wisp run` 一发，见 §4） | 库里 `error`→`stuck`（`cmd` 侧）；`agent` 侧不变 |
| **② `StatusTruncated`** | **`"interrupted"`**（真 Journal＋真 SQLite，fixture `max-tokens-toolcall`） | `"interrupted"`（函数读数；**真 CLI 一发到不了**，见 §9-2） | `error`→`interrupted`（`cmd` 侧）；`agent` 侧不变 |
| **③ 默认支／未知值** | **`"error"`** | **`"error"`** | `agent` 侧 `done`→**`error`**（唯一行为变更，fail-closed）；`cmd` 侧不变 |
| **④ 正控 `Completed`／`Cancelled`／`Failed`** | `"done"`／`"cancelled"`／`"error"` | `"done"`／`"cancelled"`／`"error"` | **三枚全不变**（真 CLI 三发，见 §4） |

**②的「变了什么」如实报**：`StatusTruncated` 的**屏幕词**是 `truncated`（`res.Status` 的原样），**权威词**是 `interrupted`（`taskstate.go:63` 注释逐字：「truncated by the round/token ceiling」）⇒ **两者不同词，且本票不可能让它们同词**：`truncated` 不在 12 枚名册里，写它会被 DAO＋CHECK 双拒，而加词属「改名册＋升 schema」另一件事（派单硬边界 2）。我按权威落 `interrupted`，**不静默**，并把这一条列进 §9-1 供裁。

## §4 「屏幕与库不再两套」逐字（真 `wisp run` 四发，全部走同一组装点）

组装点 `cmd/wisp/run.go:1072-1100`，`Journal: nil`（`:1080`）⇒ 发布二进制里 `internal/agent` 那根梯子**无写者**（`journal.go:139` `if j == nil { return nil }`），`cmd` 是唯一写者。

| 状态 | 屏幕逐字（截任务号） | 库逐字 | 同词？ |
|---|---|---|---|
| **`StatusStuck`** | `wisp run: 任务 <id> 结束（stuck，1 轮，0 次工具调用，成本 0 CNY（未计价））` | `task_log.state = "stuck"` | **✓ 不再两套** |
| `StatusCompleted`（正控） | `…结束（completed，1 轮，…）` | `"done"` | 各说各词表、主张一致 |
| `StatusCancelled`（正控） | `…结束（cancelled，0 轮，…）` | `"cancelled"` | ✓ |
| `StatusFailed`（正控） | `…结束（failed，1 轮，…）` | `"error"` | 各说各词表、主张一致 |

`Stuck` 一发的复现配方（不是造假 `Result`）：`[llm.providers.acme.models.m1] context_window = 1` ⇒ `BudgetsFor` 把冻结的 200k C22 预算缩到 **2 token** ⇒ 第一发响应即超支 ⇒ `brakeStuck(BrakeToken)` ⇒ 真 `StatusStuck`、退出码 3。

## §5 新用例逐名 ＋「为什么今天响」

`internal/agent/tasklogstate_365_test.go`
1. `Test365TaskLogStateSpeaksThe196Roster` —— **今天响**：改前默认支把 `Status("")` 与未知值写成 `done`，本用例要 `error`。
2. `Test365TaskLogStateCoversEveryStatusTheLoopCanProduce` —— 守护：枚举长一枚就落默认支，先把「沉默」标出来（今天不响，属结构钉）。
3. `Test365JournalRowsTheRosterWordOnTheAgentPath` —— **路径证据**（真循环＋真 Journal＋真库，四枚状态逐枚读回 `state`）；今天**不响**（改前该路径本就写 `stuck`/`interrupted`/`done`/`error`，这是要保住的性质，不是要改的）。

`cmd/wisp/tasklogstate_365_test.go`
4. `Test365TaskLogStateSpeaksOneWordOnBothPaths` —— **今天响**（`stuck`、`truncated` 两行：`cmd` 写 `error`）；即 AC#2「两处对同一 `Status` 读出同一个词」。
5. `Test365RunRowsAStuckTaskAsStuck` —— **今天响**（库 `error` vs 屏幕 `stuck`）。
6. `Test365RunKeepsItsWordForTheControls` —— **正控**：三发真跑（completed／failed／cancelled），今天与改后都必须绿。
7. `Test365TruncatedCannotBeReachedThroughThisCommand` —— 把 §9-2 的断链写成会响的钉子：预算是 0、真 CLI 到不了天花板；哪天有人把字段接上，本用例转红、逼出真读数。

## §6 反向判据：三读数 ＋ 还原核（全部走 `-overlay`，工作树零突变）

- **M1「摘掉 cmd 那一跳」**（`overlay` 把 `cmd/wisp/run.go` 换回 `HEAD` 版＝裸 `switch`）：`rc=1`，逐名红行：

```text
tasklogstate_365_test.go:67: cmd/wisp taskLogState("stuck") = "error", want "stuck"
tasklogstate_365_test.go:75: the two paths disagree on "stuck": cmd/wisp writes "error", the loop writes "stuck"
tasklogstate_365_test.go:67: cmd/wisp taskLogState("truncated") = "error", want "interrupted"
tasklogstate_365_test.go:75: the two paths disagree on "truncated": cmd/wisp writes "error", the loop writes "interrupted"
--- FAIL: Test365TaskLogStateSpeaksOneWordOnBothPaths
tasklogstate_365_test.go:198: task_log.state = "error" for a stuck task, want "stuck": the screen said "stuck" and the row says "error"
tasklogstate_365_test.go:205: screen and row disagree about the same task: screen "stuck", row "error"
--- FAIL: Test365RunRowsAStuckTaskAsStuck
--- PASS: Test365RunKeepsItsWordForTheControls   ← 尺不是全红
```

- **M2「把 agent 默认支改回 `done`」**（`overlay`，`diff -u` 仅 `-return memory.TaskStateError`／`+return memory.TaskStateDone` 一行）：`rc=1`，红行 `TaskLogState("failed") = "done"`、`TaskLogState("") = "done"`、`TaskLogState("a-status-this-build-does-not-know") = "done"`，＋ `task_log.state = "done" for status "failed", want "error"`（子用例 `disconnect-mid-toolcall`）。
  ★**射程声明**：`failed` 那两行红是 M2 把**显式 `Failed` 支折进默认支**造成的，比改前更狠；**改前真实的红只有未知值两行**（改前 `Failed` 有显式 `case → "error"`）。我不把 M2 冒充改前读数。
- **还原核**：不带 `-overlay` 重跑同一 `-run 'Test365'` ⇒ `internal/agent` `rc=0`、`cmd/wisp` `rc=0`；`git status --porcelain -- internal/agent cmd/wisp` ＝ 本笔 4 个文件（2 改 2 增），**无第三处改动**。

## §7 四数改前／改后（分开报、⛔ 相减）

| 包 | 尺 | `RUN_top` | `PASS_top` | `FAIL_top` | `SKIP_top` | rc |
|---|---|---|---|---|---|---|
| `internal/agent` | 改前（treeBefore `a8155492f`） | 131 | 102 | 0 | 0 | 0 |
| `internal/agent` | 改后 全量 | 138 | 105 | 0 | 0 | 0 |
| `internal/agent` | 改后 `-skip '^Test365'`（同尺对同尺） | **131** | **102** | **0** | **0** | 0 |
| `internal/memory` | 改前 / 改后 | 101 / **101** | 59 / **59** | 0 / **0** | 2 / **2** | 0 / 0 |
| `cmd/wisp`（`-run '^TestRunTextTask'`，带 DLL PATH） | 改前（`-overlay` HEAD `run.go`） | — | 3 | 0 | 0 | 0 |
| `cmd/wisp`（同上） | 改后 | — | **3** | **0** | **0** | 0 |

`agent` 全量 +7 `RUN`/+3 `PASS` ＝ 本笔新增 3 个顶级用例＋4 个子用例；`-skip '^Test365'` 后**逐数相同** ⇒ 无回归。`memory` 未动一字节、逐数相同。

**指令那一发整包**：`go test -count=1 ./internal/agent/ ./internal/memory/ ./cmd/wisp/`（带 PATH）⇒ `internal/agent ok 8.007s`、`internal/memory ok 25.111s`、`cmd/wisp FAIL 327.637s`（死点逐字 `MARK299 wndproc ... / MARK299 before CreateWindowExW`，WebView2 硬终止）⇒ **`cmd/wisp` 整包四数⛔当判据**，只作死点登记。

## §8 门禁 rc（逐门分开）

| 门 | 命令 | rc |
|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** |
| vet | `go vet ./internal/agent/ ./cmd/wisp/` | **0** |
| 格式 乙（工作树，诊断尺） | `gofumpt -l` 本笔 4 个 `.go` | **0**（无输出） |
| 格式 甲（入库 blob，判据尺） | 见下（`git cat-file blob` 取 index 版再 `gofumpt -l`） | **0** |
| d22scan | `sh scripts/d22scan.sh` | **0**（`clean - no D22 ban violations`；`internal/` 635 个 Go 文件、`cmd/` 149 个） |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **0**（`VERDICT GREEN`，最长 185 字符、rostered） |

## §9 停手上报／未决（★不自行假设，列清供裁）

1. **★`StatusTruncated` 屏幕词（`truncated`）与权威词（`interrupted`）不同词** —— 派单停手条件①。我**没有**停手不落：权威词在 `taskstate.go:63` 是**定义好的**（不是我的假设）、`runExitCode`（`run.go:80`）本来就把 `Stuck`/`Truncated` 归同一档「环路刹车＝3」，且「同词」在本票内**不可能达成**（名册无 `truncated`、加词属另一件事）。落的是权威词，**分歧逐字上报**：要不要让屏幕也说库的词，属「改屏幕文案」＝另一张票。
2. **`StatusTruncated` 在 `cmd/wisp` 路径上今天不可达**：`[llm.roles.chat].max_output_tokens` 已由 `llm/resolver.go:205` 落进 `llm.Endpoint`，但 `agent.Config` 没有承接字段、`agent.Assembler.MaxOutputTokens`（`prompt.go:132`）**全仓无赋值** ⇒ `llm.Request.MaxOutputTokens` 恒 0 ⇒ 供应商无从「撞天花板」。修它要动 `cmd/wisp` 组装（超出本笔写面「只 `run.go` 的 `taskLogState`」）⇒ 只报不修。
3. **`365-a1` 的归因我复核成立**：发布组合里 `Journal: nil`（`run.go:1080`）⇒ `internal/agent` 那根梯子在已发布二进制里无调用者；票 196 读侧名册（`TaskStateFor`／`TaskStateFromD43`／`D43State`）在 `cmd/**`＋`internal/**` 非测试代码里**仍 0 调用者**（现跑 `git grep`）。
4. **未动**：12 枚名册、`internal/tools/subagent_197.go:451-457／480-486`（第三处同形映射、另一把尺）、`internal/memory` 任何字节、schema。
