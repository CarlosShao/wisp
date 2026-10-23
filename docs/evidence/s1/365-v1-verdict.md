# 票 365 终裁腿 `365-v1`：`AC#2` **能翻**（两处一根梯子、库里的词真的变了）

- 裁的对象：`365-r1`（`6b9a83c04`，5 文件＝2 改 3 增、**0 件删除**）；件 `docs/evidence/s1/365-r1-impl.md` **150 行**（逐行读过）
- 票面：`.scratch/wisp/issues/365-two-taskLogState-functions-disagree.md` **55 行**（通读）；落笔时 HEAD `0b690c40a`
- 环境：Go 1.27.1 windows/amd64、`GOFLAGS=`、只 `unset GOTMPDIR`、`cmd/wisp` 带 `third_party/sherpa-onnx` PATH
- **尺的射程**：夹具全部在仓外 `D:\tmp\wisp\20261023\365v1\**`、经 `-overlay` 入场；仓内 `.go` 改动 **0 字节**（`git status --porcelain -- internal/agent internal/memory cmd/wisp`＝空，量于**最后一发变异之后**）
- 本腿只写本件＋仓外；未翻框、未改票名／`Status:`／`-done`、未 push、未碰 196/188 原句

## §0 一句话裁决

**`AC#2` 能翻。** 真跑改前 `stuck`/`error`、改后 `stuck`/`stuck`；两根梯子收成一根；`internal/memory` 零字节；我自造三发变异各自红在具名行；四数同尺逐数相同；五门 rc=0。**不挡翻框的四条保留在 §10**（没有一条要求 365 返工）。

## §1 九项判语（读数见 §2–§9）

| # | 项 | 判语 |
|---|---|---|
| 1 | 票面前提新鲜度 ＋ `365-a1` 归因 | 前提**新鲜**；归因**结论对、措辞错**（"无调用者"→"有调用、无写者"），见 §2 |
| 2 | 核心一：库里真的变了 | **成立**（两条路径都给了库读数，见 §3） |
| 3 | 核心二：屏幕与库不再两套 | **成立**（真跑＋改前对照，见 §4） |
| 4 | 核心三：`done`→`error` | **该成立**（fail-closed，见 §5） |
| 5 | "一根梯子"＋`memory` 零改动 | **成立**（见 §6） |
| 6 | 反向判据＋`M2` 射程声明 | 我自造 M1/M2/M-old **全红且具名**；它的射程声明**对**（见 §7） |
| 7 | 四数 | **逐数相同**（见 §8）；整包 `cmd/wisp` 未复跑、只登记死点 |
| 8 | 门禁＋格式甲乙 | 五门 rc=0、两尺 0 行（见 §9） |
| 9 | 与并发腿无交集 | **成立**（见 §9 末） |

## §2 项 1：前提新鲜度 ＋ `365-a1` 归因复核

- **前提（对改前 blob 复量）**：`6b9a83c04^:internal/agent/loop.go:1048 func taskLogState(s Status) string`（`:1052 case StatusStuck:`／`:1053 return "stuck"`）；`6b9a83c04^:cmd/wisp/run.go:1361 func taskLogState(s agent.Status) string`（无 `Stuck`/`Truncated`/`Failed` 支、`:1370 return "error"`）⇒ 票面 F1–F4 **逐字成立**。
- **没腐烂**：`loop.go`＝`31dec50a8…`、`run.go`＝`9fc8fbe97…` 在 `ef3ec49bd`（`365-a1` 的 HEAD）、`dd16b005c`（`365-r1` 锚）、`6b9a83c04^` **三处逐枚相同**；`ef3ec49bd..6b9a83c04^` 内碰这两文件的提交 **0 枚** ⇒ 前提从普查到落地没漂。
- **归因（结论）成立**：唯一 Loop 装配点 `run.go:1072`（`agent.New` `:1100`）带 `Journal: nil`（`:1080`）；`task.spawn` 子环 `internal/tools/subagent_197.go:320 opt := base` → `:347 agent.New(opt)` 继承 nil ⇒ **发布二进制里库中的 `stuck` 行只可能由 cmd 写**。
- ★**归因（措辞）要改**：`TaskLogState` 在发布二进制里**有调用**——`loop.go:1014` 是**实参求值**，每发结束都调；**没有写者**——`journal.go:141 if j == nil { return nil }` 提前返回、**值被丢**。`365-r1` 件 §4 引的 `journal.go:139` 偏 2 行（`:141` 才是那句）。⇒ 台账该记"有调用、无写者"，⛔ "没有调用者"。

## §3 项 2：两条路径各写什么（我的夹具；读回尺＝`memory.Store.TaskLogByID`，⛔ 函数返回值）

- **甲（`internal/agent`）**：真 `Loop.Run`＋真 `Journal`（`Options.Journal`＝真 `memory.Store`）＋真 SQLite，golden `repeat-echo` 驱动真 C22 刹车。
- **乙（`cmd/wisp`）**：真组合根 `runTextTask`（`wisp run` 与常驻腿共用）＋`adaptertest` 真 HTTP 供应商；`context_window = 1` ⇒ 冻结的 200k C22 预算缩到 2 token ⇒ 第一发超支 ⇒ `brakeStuck` ⇒ 真 `StatusStuck`、退出码 3。

```text
Z365V1-AGENT | label=Stuck | fixture=repeat-echo | Status="stuck" | task=575f1ae9-51e7-48e2-9903-20ef84ab8262 | task_log.state="stuck"
Z365V1-CMD   | exit=3 | screen="wisp run: 任务 cee2ca3d-36ca-4245-aaa4-a75316f4a30d 结束（stuck，1 轮，0 次工具调用，成本 0 CNY（未计价））"
Z365V1-CMD   | task=cee2ca3d-36ca-4245-aaa4-a75316f4a30d | task_log.state="stuck"
```

同一夹具另三枚（甲路径库读数，改前/改后都量了）：`Truncated`→`"interrupted"`、`Completed`→`"done"`、`Failed`→`"error"`。
⇒ **同一枚 `StatusStuck`：两条路径都写 `"stuck"`**（票面缺陷在库里已不存在）。

**`Truncated` 在 `cmd` 侧只有函数读数，够不够？——我独立判：够**，理由**不是"不可达"**而是**委托身份**：`run.go:1367` 是一行委托，`run.go:1196 state := taskLogState(res.Status)` 的返回值原样进 `:1201 State:` 与 `:1212 FinishTaskLog(…, state, …)`，**函数读数与库读数之间没有第二跳**；`interrupted` 在 12 枚名册内，我的甲路径夹具已把 `interrupted` 真写进同一版 v4 库 ⇒ 该写必落。

## §4 项 3：屏幕与库——真跑 ＋ 改前对照（`-overlay` 换回 `6b9a83c04^` 的 `run.go`，`git hash-object`＝`9fc8fbe97…` 逐字节）

```text
改前：screen = wisp run: 任务 276918e8-fe98-4180-81d6-32c0c38905ab 结束（stuck，1 轮，0 次工具调用，成本 0 CNY（未计价））
      库     = task_log.state="error"      ← 票面缺陷原样复现（屏幕 stuck／库 error）
改后：screen = wisp run: 任务 cee2ca3d-36ca-4245-aaa4-a75316f4a30d 结束（stuck，1 轮，0 次工具调用，成本 0 CNY（未计价））
      库     = task_log.state="stuck"      ← 同一发，两侧同词
```

## §5 项 4：`agent` 侧 `done`→`error`（本笔唯一行为变更）——**判：该成立**

① `runExitCode`（`run.go:77-95`）对未知 `Status` 落 `default: return 1`（失败）⇒ 库若写 `done` 就与退出码、与屏幕自相矛盾；② 名册语义 `done`＝"loop/CLI completed"（`taskstate.go:60`），未知结局没有"出了答案"的证据，写 `done` 是**把无名结局报成成功**；③ 状态列是审计/统计用，**假成功比假失败更贵**（抬高成功率、藏 bug），fail-closed 与本仓 bridge/gate 的既有形状一致；④ 今天生产路径**不可达**未知值（终态只由六枚常量触发：`loop.go:397/459/500/510/541/971/980`；全仓无 `Status(...)` 构造）⇒ 这是**结构钉**不是改存量；`StatusControl` 旧默认的答案 `done` 被**显式保留**（`:1071`）、且 Control 不开行（`:52-53`）⇒ **零回归**。结论：**不该改回 `done`**。

## §6 项 5："一根梯子"＋`memory` 零改动

- `cmd/wisp/run.go` 内 `switch` 共 5 处：`:77`（Status→**退出码**，不同对象）、`:83`（error_class）、`:1340`（压缩态）、`:1413`/`:1550`（事件 kind）⇒ **无一处 Status→state 词**；全仓 `case agent.Status` 在 `cmd/**` 非测试代码 **0 命中**。
- `TaskLogState`（`loop.go:1069`）的 `return` **全是 `memory.TaskState*` 常量**；函数体内引号 **0**（`sed -n 1047,1092p | grep '"'` 只命中注释 6 行）⇒ "零字符串字面量"成立。
- **`internal/memory` 零改动**：`git diff --numstat 6b9a83c04^ HEAD -- internal/memory` 空；`6b9a83c04 HEAD` 亦空。
- 唯一 Status→`task_log.state` 转换点＝`agent.TaskLogState`；两个写者只传它的值：`journal.go:144`←`loop.go:1014`、`run.go:1205/:1212`←`run.go:1196`。
- **import cycle 断言成立**：`internal/memory` 对 `internal/agent` 的 import **＝0**（grep 只命中注释）；`agent` 已 import memory（`journal.go:8`、`tools.go:11`、`loop.go:17`）⇒ 让 `memory` 收 `agent.Status` 确会造环 ⇒ **〔乙〕的层级只能这样落**（词表仍属 `memory`，投影只此一跳）。
- 第三处同形投影 `internal/tools/subagent_197.go:451-457/480-486` 写的是 `TaskOutput.State`（D43 名、`statemachine.Valid`）⇒ **⛔ 同一列、不计第二根梯子**；但它是 `Status` 的第二个投影（票外，`365-r1` §9-4 已报）。

## §7 项 6：我自造的反向判据（⛔ 只看 `M1`/`M2`）＋ `M2` 射程声明

三发变异**全部在仓外、经 `-overlay`、无工作树突变**；rc 全取壳级 `$?`。

- **M1＝换回改前 `run.go`**（真字节）⇒ `rc=1`，逐名红行：`tasklogstate_365_test.go:67: cmd/wisp taskLogState("stuck") = "error", want "stuck"`／`:75: the two paths disagree on "stuck": cmd/wisp writes "error", the loop writes "stuck"`／`:67: …("truncated") = "error", want "interrupted"`／`:75: … "truncated" …`／`:198: task_log.state = "error" for a stuck task, want "stuck": the screen said "stuck" and the row says "error"`／`:205: screen and row disagree about the same task: screen "stuck", row "error"`；`--- FAIL: Test365TaskLogStateSpeaksOneWordOnBothPaths`、`--- FAIL: Test365RunRowsAStuckTaskAsStuck`；**正控三子用例 PASS**（尺不是全红）。
- **M2＝默认支 `TaskStateError`→`TaskStateDone`（一行）** ⇒ `rc=1`：`…:51: TaskLogState("failed") = "done", want "error"`／`…:51: TaskLogState("") = "done", want "error"`／`…:51: TaskLogState("a-status-this-build-does-not-know") = "done", want "error"`／`--- FAIL: Test365TaskLogStateSpeaksThe196Roster`／`…:124: task_log.state = "done" for status "failed", want "error"`＋`--- FAIL: …/disconnect-mid-toolcall`。
- **M-old＝把新函数改回改前那套 arms（显式 `Failed→Error`、default→`Done`）** ⇒ `rc=1`，**只有两行红**：`TaskLogState("") = "done", want "error"`、`TaskLogState("a-status-this-build-does-not-know") = "done", want "error"`；`Test365JournalRows…` **四子用例全 PASS**（含 `disconnect-mid-toolcall` failed→`error`）。
- ⇒ **它的射程声明对**：`failed` 那两行红是 `M2` 造成的，**改前真实的红只有未知值两行**，它没拿 `M2` 冒充改前。措辞小瑕：当前函数**没有**显式 `Failed` 支（默认支本来就替 `Failed` 答题），说"把显式 `Failed` 支折进默认支"宜改成"默认支的答案被 `M2` 改成 `done`、`Failed` 跟着变"——**结论不变**。
- 还原核：不带 `-overlay` 重跑 `-run '^Test365'` ⇒ `internal/agent` rc=0、`cmd/wisp` rc=0（8 枚用例全 PASS），且 `git status` 三包仍空。

## §8 项 7：四数（分开报、⛔ 相减；全量走 `-v`，`RUN_top` 只数顶格 `=== RUN   Test`）

| 尺 | RUN | PASS | FAIL | SKIP | rc |
|---|---|---|---|---|---|
| `internal/agent` 改前（`6b9a83c04^` 仓外导出树） | **131** | **102** | **0** | **0** | 0 |
| `internal/agent` 改后 `-skip '^Test365'`（同尺对同尺） | **131** | **102** | **0** | **0** | 0 |
| `internal/agent` 改后全量 | 138 | 105 | 0 | 0 | 0 |
| `internal/memory` 改前 / 改后 | 101 / **101** | 59 / **59** | 0 / **0** | 2 / **2** | 0 / 0 |
| `cmd/wisp -run '^TestRunTextTask'` 改前（overlay 旧 `run.go`） | 3 | **3** | 0 | 0 | 0 |
| 同上 改后 | 3 | **3** | 0 | 0 | 0 |

- **整包 `cmd/wisp`：本腿未复跑**（死点判据外、327s 级）⇒ 只登记 `365-r1` 件 §7 的逐字死点 `MARK299 … before CreateWindowExW`（WebView2 硬终止），**⛔ 当本票判据**。

## §9 项 8／9：门禁 rc ＋ 与并发腿无交集

| 门 | rc | 读数 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | 空输出 |
| `go vet ./internal/agent/ ./internal/memory/ ./cmd/wisp/` | **0** | 空输出 |
| `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；`internal/` 635 个 Go、`cmd/` 149 个 |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；最长 185 字符、rostered=61、not-in-roster=0 |
| 格式**甲**（入库 blob，判据） | **0** | `gofumpt -l` 0 行；index blob 与 HEAD 逐枚同哈希（loop `a04bab206`、run `d7bf4aabb`、测试 `357fc1e1e`/`b325d8aa5`）|
| 格式**乙**（工作树，诊断） | **0** | `gofumpt -l` 0 行 |

- ★**行尾旁证**（⛔ 用 `gofumpt` 判）：HEAD blob 是 **CRLF**（每行带 CR：`loop.go` 1223/1223、`run.go` 1578/1578，用 `grep -c $'\r'` 量）⇒ **"`gofumpt -l` 输出 0"⛔ 不能反推"是 LF"**；本条不影响上表两尺读数。
- **交集**：`git diff --stat a8155492f 771bf451e -- internal/agent internal/memory cmd/wisp` **空**；两笔分别只碰 `docs/reports/pending-and-issues.md`、`docs/evidence/s1/114-r1-impl.md`；`6b9a83c04^..HEAD` 内碰这三包的提交只有 `6b9a83c04` 一枚 ⇒ **它自报"与并发腿无交集"成立**。

## §10 保留（四条，**都不挡翻框**）

1. **`Truncated` 的 `cmd` 侧只有函数读数**：判**够**，依据是 §3 的委托身份（不是"不可达"）。另：`365-r1` §9-2 的"不可达"准确说是**没有我们设的天花板**（`Assembler.MaxOutputTokens` 全仓无赋值 ⇒ `openaichat` 的 `max_tokens` 被 `omitempty` 省略），供应商仍可按自身模型上限返回 `length`（`loop.go:490` 判 `turn.Stop == llm.StopMaxTokens`）⇒ **票 370 标题的"永远到不了"比量到的强**；370 的 AC#0 ④（手工设 `Request.MaxOutputTokens=1` 撞出 `Truncated`）恰好能把它量准。
2. **屏幕词 `truncated` ≠ 库词 `interrupted`**：编排者的裁**对**（权威词 `taskstate.go:63`；`truncated` 不在 12 枚名册、加词＝改名册＋升 schema；`AC#2` 判据是"两处函数同词"）。**但** 370 一旦接通天花板，这一对会立刻在同一发任务上可见 ⇒ **370 的验收里该带"屏幕说哪个词"的判决**（370 票面今天没写这一条）。
3. **`done`→`error`** 是本笔唯一行为变更，判**成立**（§5）。
4. 票面 5 格仍空、`Status:` 仍 `ready-for-agent`（我按纪律未动）⇒ **翻框／改 `Status:` 是编排者的动作**。

## §11 提交与越界

- 本件 `docs/evidence/s1/365-v1-verdict.md`（本腿**唯一**仓内写面）；首笔提交 **`97bd0cc91`**（119 行、1 文件、0 删除），本行由补记提交追加第二笔。
- 逐笔越界：本腿共 2 笔（首笔 `97bd0cc91` 本件全文；补记笔＝只改本 §11 这一行），**每笔只有这一个文件**、无 `D`/`A` 越界；`internal/**`、`cmd/**`、`scripts/**`、`.github/**`、`.scratch/**`、`docs/reports/**`、`docs/PLAN.md`、`docs/specs/**` **零改动**；未 push。
