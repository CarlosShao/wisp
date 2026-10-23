# 365 — 两处 `taskLogState` 各写各的：同一枚 `StatusStuck` 一边写 `stuck`、一边写 `error`

- Status: **DONE（2026-10-23 17:5x 编排者收口＝五格全闭、文件名加 `-done`）** ~~**ready-for-agent**~~
- 立案：2026-10-23 10:5x +08 编排者立
- ★**来路（实现者签名）**：`196-r1` §9-1（件 `docs/evidence/s1/196-r1-impl.md` **139 行**、提交 `6e19f4a43`／`e89d1fe98`／`2c3628679`）
- ★**编排者已现跑核实**（见下）
- 关联：票 **196**（已落地甲＋乙、`cmd/**` 不在它写面）· `cmd/wisp/**` · `internal/agent/**`

## 这是什么（人话）

**同一个任务状态，两条生产路径写进库里的词不一样。**

★**编排者现跑核实的原话形状**：
- `internal/agent/loop.go:1048` `func taskLogState(s Status) string` ⇒ ★**有 `case StatusStuck: return "stuck"`**。
- `cmd/wisp/run.go:1361` `func taskLogState(s agent.Status) string` ⇒ ★**没有 `StatusStuck` 那一支**，默认落 `:1370 return "error"`。
- ⇒ ★★**同一枚 `StatusStuck`：库里一边是 `stuck`、一边是 `error`。**

★**这就是票 196 那个"可观察分裂"的真身**：票 196 的甲（`internal/memory/taskstate.go` 的权威映射 ＋ DAO 校验）**已经盖住"写进库的那一跳"**，
★**但两个 `taskLogState` 函数本身还在各自映射** ⇒ **分裂的源头没除根**。

## 结案判据

- [x] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：
      ①**逐条复核上面两处函数**（★带 `file:line` 现跑复认、⛔ 抄本票）；★**行号会不会漂要现跑**。
      ②★**把"分裂"量成一发真读数**：造一发 `StatusStuck` 走 `cmd/wisp` 那条、一发走 `internal/agent` 那条 ⇒ **逐字贴两行库里读到的 `state`**。
      ③★**穷举两处 `case` 集合的差集**：★**除 `StatusStuck` 外还有没有别的分歧**（★`StatusTruncated`／`StatusFailed`／`StatusControl` 各比一遍）。
      ④**正控**：造一发两边都该写**同一个词**的状态（如 `StatusCompleted`）⇒ ★**两处必须同词** ⇒ **证明尺不是全红**。
- [x] **AC#1 裁形（编排者拍）**：三形各写代价＝
      〔甲〕**`cmd/wisp` 那处改成调用 `internal/agent` 的**（单一实现、`cmd` 侧只留一层薄壳）；
      〔乙〕**两处都改成走票 196 的权威映射 `internal/memory.TaskState`**（`taskstate.go` 已存在）；
      〔丙〕**只补 `cmd/wisp` 缺的那几支**（最小动作、但两处仍在各写各的）。
      ★**未裁形任何实现都⛔ 开工**。
      ★**注意**：★**〔甲〕〔乙〕都要动 `cmd/wisp/**`**（★共享编译面、每波只许一枚写腿）；★**〔乙〕还跨到 `internal/memory`**。
- [x] **AC#2 落地＋牙（写腿）**：按裁定那一形做；★**每处改动配一枚今天响、不修不响的用例**；
      ★**反向判据**＝把新加的那一跳摘掉 ⇒ 用例必须**转红并逐名报红行**；三读数（变异红／还原绿／原样绿）逐字贴 ＋ **还原核**。
      ★★**验收判据是"两处对同一 `Status` 读出同一个词"**（★⛔ 只测"某处改了"）。
- [x] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet` 动过的包／`go test -count=1`（**四数分开报、⛔ 相减**）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。
      ★**格式门禁甲乙两尺都报**（甲＝入库 blob、**判据**；乙＝工作树、**只作诊断**）。
      ★★**`cmd/wisp` 整包四数 ⛔ 当判据**（未判决恒 1、红名册逐发漂移、截断点后整族 `RUN=0`）⇒ 只按名作差＋逐发登记死点。
      ★★**`go test ./cmd/wisp/` 有预存红 `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND、测试进程一次没起）** ⇒ ★**⛔ 当本票的红**（票 196 用原样 HEAD 树复跑同样失败 ⇒ 环境）。
- [x] **AC#4 纪律格**：⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 196／票 188 任何原句。

## ★硬提醒（⛔ 踩）

1. ★★**`cmd/wisp/**` 每波只许一枚写腿**、且 `go test ./cmd/wisp` 重包同一时刻只跑一个。
2. ★**⛔ 顺手把两处都改成 D43 词**：★**票 196 的名册已经定了 12 枚**（旧词 7 ＋ D43 里任务行真会取的 5 枚）⇒ ★**落地要跟那张名册对齐、⛔ 自造第三套**。
3. ★**票 196 的乙（CHECK 约束）已在库上生效** ⇒ ★**任何新词必须先进名册、否则写不进库**（★那是好事、也是约束）。

## 排程

`365-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `365-r1`（写腿＝`AC#2`）→ `365-v1`（非实现者裁）。

## Progress log

- [2026-10-23 10:5x +08] agent=`编排者` did=立案（来路＝`196-r1` §9-1；★两处函数编排者已现跑核实：`internal/agent/loop.go:1048` 有 `case StatusStuck: return "stuck"`、`cmd/wisp/run.go:1361` 无那一支且 `:1370` 默认返 `"error"`）next=派 `365-a1` 只读腿。

## ★★`AC#2` 判语（2026-10-23 17:5x 编排者，凭非实现者 `365-v1`，件 `docs/evidence/s1/365-v1-verdict.md` 119 行、`97bd0cc91`＋`4db5ee7ce`）

★**`AC#2` 能翻**（不缺件、不返工）。四条保留不挡翻框。

**① 九项逐项（判语自造夹具、⛔ 抄实现腿的件）**
- ★**前提新鲜**：`6b9a83c04^` 的 `loop.go:1048`／`run.go:1361` 逐字成立；两 blob `31dec50a8`／`9fc8fbe97` 在 **三处同哈希**、其间 **0 提交**碰它们。
- ★★**库里真的变了**（尺＝`memory.Store.TaskLogByID`）：甲 `internal/agent`（真 Loop ＋ 真 Journal ＋ 真 SQLite）`StatusStuck ⇒ "stuck"`；
  乙 `cmd/wisp`（真组合根 `runTextTask`、`context_window=1` ⇒ 真 `brakeStuck`、`exit=3`）屏幕「结束（stuck，1 轮，…）」、库 **`"stuck"`**；
  另三枚甲路径：`Truncated→interrupted`／`Completed→done`／`Failed→error`。
- ★★**屏幕与库不再两套**：★**`-overlay` 换回 `6b9a83c04^` 的 `run.go`（hash 逐字节 `<9fc8fbe97>`）** ⇒ 改前屏幕 `stuck`／库 **`error`**（缺陷原样复现）；改后屏幕 `stuck`／库 **`stuck`**。
- ★★**`done`→`error` 判成立（四条理由）**：①`runExitCode` 对未知 `Status` **本就返回 1（失败）** ⇒ 库写 `done` 会**与退出码／屏幕矛盾** ②名册 `done`＝"出答案" ③**假成功比假失败更贵** ④未知值在生产**不可达**（终态只由六枚常量触发）⇒ `Control→done` 显式保留、Control 不开行 ⇒ **零回归**。
- ★**一根梯子成立**：`cmd/wisp/run.go` 无 Status→状态词的 `switch`（★全仓 `case agent.Status` 在 `cmd/**` 非测试 **0 命中**）；`TaskLogState` 返回全是 `memory.TaskState*` 常量、**函数体引号 0**；★**`internal/memory` 零改动**（`git diff --numstat` 空）；★**import cycle 断言成立**（`memory→agent` import ＝ 0）。
- ★**反向判据判语自造三发**：`M1`（换回改前 `run.go`）`rc=1`、红在 `tasklogstate_365_test.go:67/75/198/205`、**正控三子用例 PASS**；
  `M2`（默认支 `Done`）`rc=1`、红 `:51×3＋:124＋disconnect-mid-toolcall`；★**`M-old`（把新函数改回改前 arms）`rc=1`、只有未知值两行红、路径证据四子用例全 PASS**
  ⇒ ★★**"它的 `M2` 射程声明经核成立"**（`failed` 红是**变异造成**、**改前真实红只有未知值两行**）。
- ★**四数**：`internal/agent` 改前 `131/102/0/0` → 改后 `-skip '^Test365'` **`131/102/0/0`**（全量 `138/105/0/0`）；`internal/memory` `101/59/0/2` 两发相同；
  `cmd/wisp -run '^TestRunTextTask'` 改前／改后均 **3 PASS**、`rc=0`。
- ★**门禁复跑**：`build`／`vet`／`d22scan`／`pathlen` **全 `rc=0`**；格式**甲（index blob）0 行、乙（工作树）0 行**、index 与 HEAD 逐枚同哈希。
- ★★**旁证（判语给的一条硬提醒）**：★**`HEAD` blob 是 CRLF**（`grep -c $'\r'` **1223/1223**、**1578/1578**）⇒ ★**"gofumpt 输出 0" ⛔ 不能反推是 LF**。
- ★**与并发腿无交集成立**：`git diff --stat a8155492f 771bf451e -- …` 空。

**② ★★★判语纠正我两条记录（我现跑核实、都对）**
1. ★★**`365-a1` 的归因措辞要改**：★**`TaskLogState` 在发布二进制里**有调用****（`loop.go:1014` 是**实参求值**、每发结束都调）
   ★**但它**无写者****：★`journal.go:141` **`if j == nil { return nil }`** 提前返回、**值被丢**。
   ★★**所以正确记法是「有调用、无写者」、⛔ 「没有调用者」**；★并纠正 `365-r1` 件 §4 引的 `journal.go:139` **偏 2 行**（`:141` 才是那句）。
   ★**结论与〔乙〕形仍成立**（★该梯子在任何挂 Journal 的装配里**都是活的**、且现在正是 `cmd` 委托的那根）。
   ★★**我现跑核**：`loop.go:1014` `j.finishTask(wctx, string(TaskLogState(status)), …)` ✓ ／ `journal.go:141` `if j == nil { return nil }` ✓
2. ★★**票 370 要补两条**：★①**一旦接通天花板，屏幕词 `truncated` 与库词 `interrupted` 会立刻在同一发任务上可见** ⇒ **370 的验收要带"屏幕说哪个词"的判决**；
   ★★②**370 标题"永远到不了"比我量到的强** —— ★**证到的是"没有**我们设的**天花板"**（★`Assembler.MaxOutputTokens` 全仓无赋值 ⇒ `max_tokens` 被 `omitempty` 省略），
   ★**供应商仍可按自身模型上限返回 `length`**（`loop.go:490` `if turn.Stop == llm.StopMaxTokens`）。
   ★★**我现跑核**：`loop.go:490` 确是 `StopMaxTokens` ✓ ⇒ **我的标题写强了、已改。**

**③ ★判语判我两处裁（都认）**
- ★`StatusTruncated` 屏幕词≠权威词：★**我裁得对**（权威词在 `taskstate.go:63`；`truncated` 不在 12 枚名册、加词＝改名册＋升 schema；`AC#2` 判据是"两处函数同词"）⇒ **实现腿没停手不落是对的**。
- ★★`Truncated` 在 `cmd` 侧只有**函数读数**：★**判够**，★**但依据应改** —— ★**不是"不可达"、是**委托身份****（`run.go:1367` 一行委托；`run.go:1196` 的返回值**原样进** `:1201 State:` 与 `:1212 FinishTaskLog(…,state,…)`、**中间没有第二跳**）。
  ⇒ ★★**这正是 `A1188②` 那族的反面**：★**"委托一跳"能替代"端到端读数"当且仅当**中间没有第二跳**、且这一点被逐行核过。**

**④ ★四条保留（不挡翻框）**：见判语 §12（★**都属"可以更严"、⛔ "不成立"**）。
