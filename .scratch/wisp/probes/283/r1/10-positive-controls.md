# 票 283 / 腿 283-r1 —— 两处盲区均选「补尺」：两条新尺 + 三发正控四件套读数

- 腿名：`283-r1`（窄写腿；写面＝`internal/tools/**` 新文件 + 本目录 `*.md`）
- 起手锚 commit `a99f9a39`（`.scratch/wisp/probes/283/r1/00-anchor.md`）；尺 commit `039ec93c`
- 时刻：起手 `2026-10-09 08:39:49 +08`；正控收束 `2026-10-09 08:53:17 +08`
- 零翻框、零 push；票面／台账／`issues/README.md` 一字未改（归编排者）；`frontend/**`／`design/**` 零读零引

## AC#1 查重（三处，逐处具名）

**① 票池**（`.scratch/wisp/issues/`）。尺＝`grep -rn "TaskID" .scratch/wisp/issues/` 逐枚读过。
结论：**没有一张票在册"corr≠taskID 下 TaskID 访问器／身份读点"的尺，也没有整链尺**。相邻而不重复的三处，逐枚具名：
- `259-…-ruler-holes-done.md:76`：`loop.go:654 两枚 id 塌在一起`（票 242 落地**前**的旧形）已判"不新立票、归 242 一族"——那是"载具拆开后产码默认仍同值"的登记，不是访问器尺。
- `236-…-six-cells-….md` AC#3：钉的是"**任务 id 由谁铸造**"（铸造者维），问的不是"读到的是哪一枚 id"；且它自己写明"需要一发真跑组合根的 CLI 接缝用例"仍缺。
- `242-…-zero-rulers.md:33`／`60`：面板出向读面那族，与 `internal/tools` 的身份读点不同宇宙。
- 近邻"发现者而非尺"：`282-…`（其 AC#1 正是种坏后发现盲区的那一发）与票 283 自身。
**② `issues/README.md` 归口句面**。尺＝`grep -n "归口" issues/README.md`＝**零命中**（全文没有"归口"二字）；`grep -n "重复|新立|立票|新票"` 的命中全部属票 262 长名门那一族（`:67-83`），与本两处零交集。⇒ 本两处**没有任何归口句**把它们推给别的票或别的仪器。
**③ 台账**（`docs/reports/pending-and-issues.md`）。尺＝`grep -n "A74[5-9]|A75|283|整链|回退形" …`。
读数：`A745`（242 普查＋四条边界）、`A748`（242-corrland 收＋越格自报）、`A749`（282-v1b 收＋两处新盲区⇒补票 283）；另 `:14400`（259 残留，上文①）。⇒ 台账里这两处**只有"缺口登记"**（A749 原文：`缺"真 loop per-call corr→bridge"整链尺，该支今天没有仪器`），**没有既有尺**可重复。
另核：`grep -rn "withCancel(" --include=*.go internal cmd` 的 `internal` 侧命中只有产码（`cancel.go` 定义、`bridge.go:527` 调用点）——**没有任何既有测试构造过 corr 与 taskID 不同的 handle**。

## AC#2 两处：都选「补尺」（落点＝新文件，零改既有件）

落点：**`internal/tools/ticket283_corr_identity_rulers_test.go`**（新增 227 行；`package tools`，复用包内现成 harness `build221`／`windowGate221`／`fake197Provider`，C5 注入面仍只用假 provider）。
`git show --stat 039ec93c` ＝ `1 file changed, 227 insertions(+)`，**零删除、零改既有断言**（AC#3）。

- **尺一（对应盲区 1：`cancel.go:69` 回退形无仪器）**＝`Test283TaskIDAccessorIsTheTaskIDNotTheCorr`。
  最小形状：一个 `cancelHandle{corr:"corr-283", taskID:"task-283"}` 断言 `TaskID==taskID`、`CorrelationID==corr`；再加"没铸 task id 的 handle（`corr` 在、`taskID` 空）必须答空"——两处身份读点的 fail-closed 拒绝就建立在"空就是空"上。
- **尺二（对应盲区 2：身份改读无整链尺）**＝`Test283IdentityChainThroughTheRealLoop`。
  整链：真 `agent.Loop`（`callCorr` 每调用一枚 corr，形如 `taskID#callID`）→ 真 `bridge.Execute`（corr 与 taskID 分别进 `cancelHandle`）→ 真 `task.spawn`（孩子落名册）→ 真 `task.cancel`（父停自己的孩子）。第二跳的子 id 由 provider **从工具回执正文里取**（脚本型 C5 provider 读 `llm.Request` 的历史，与生产里模型读到的是同一份字面），不是钉死的常量。
  断言三签：`res.TaskID` 全链成立（loop 随机铸 id）、子行 `ParentTaskID==res.TaskID`、cancel 回执**必须走到权威判定通过支**（"没能停掉"，孩子此时已跑完；"拒绝停止"只该在调用者读错身份时出现）；另把两枚调用在 bridge 记账层的行 corr 逐行核为"非空、以 task id 为前缀、不等于它"。
  与 `loop_approval_test.go:219-223` 的关系：那枚（fs.write 路径）钉的是**行 corr 形状**；本尺里同形断言不重复它的劳动，而是本场景的前提守卫——若哪天 loop 把 corr 退回 task id，须由它把"本尺的 corr≠taskID 前提不成立"直接报出来（否则身份断言会静默变空转）。

## 正控四件套（逐发；日志 `/tmp/283-seedA|seedB|seedB2|seedC|final.txt`）

基线：`go test ./internal/tools/ -run 'Test283' -count=1 -v` ⇒ `rc=0`，`=== RUN` 2 枚、2 PASS（跑法含 PATH 只注 `third_party/sherpa-onnx`）。

**种子 A** —— `cancel.go:69` `return h.taskID` ⇒ `return h.corr`（盲区 1 的原始形状）
- 种前 hash＝`87ac0624`（与 A749 复核值一致）；sed 复量＝`return h.corr`
- 红：`rc=1`，两尺都 FAIL；红句：
  - `ticket283_corr_identity_rulers_test.go:40: TaskID = "corr-283", want the task id: a handle whose corr differs from its task id must not answer with the corr`
  - `…_test.go:194: 子代理行的 ParentTaskID = "93309383-…#call-283-spawn", want 宿主 task id "9330938…"`
  - `…_test.go:206: task.cancel 的回执没有走到权威判定通过那一支: "本用例不带 L2 队列"`
- 还原：`return h.taskID` ⇒ hash `87ac0624` 逐字回；`git status --porcelain -- internal/tools/cancel.go` 空

**种子 B** —— `subagent_197.go:262` `parentID := TaskID(ctx)` ⇒ `CorrelationID(ctx)`（盲区 2 的一处身份读点）
- 种前 hash＝`8c9266a7`；sed 复量＝`parentID := CorrelationID(ctx)`
- 红（两跑同形，确定性）：`rc=1`；单元尺 PASS（靶向），整链尺 FAIL：
  - `…_test.go:194: 子代理行的 ParentTaskID = "035538fa-…#call-283-spawn", want 宿主 task id "035538fa-…"`（第二跑 `c64aa21d-…` 同形）
  - `…_test.go:206: …: "本用例不带 L2 队列"`（**观测登记**：该种子下取消跳被上游 L2 判定拦下——该回执文本逐字出自 `windowGate221` 的 Approval 支——未及 Execute；此 L2 的上游通道本程未追到底，主红句是 `:194`，归因不受影响）
- 还原：hash `8c9266a7` 逐字回；`git status --porcelain -- internal/tools` 空

**种子 C** —— `task.go:708` `caller := TaskID(ctx)` ⇒ `CorrelationID(ctx)`（盲区 2 的另一处身份读点）
- 种前 hash＝`3ec8d498`；sed 复量＝`caller := CorrelationID(ctx)`
- 红：`rc=1`；单元尺 PASS；整链尺 FAIL，走的是干净通道（Execute 内部身份拒绝）：
  - `…_test.go:206: …: "拒绝停止：4d3e9318-… 是任务 7c04d1c2-… 派生的孩子，不是调用者 7c04d1c2-…#call-283-cancel 的孩子——只有父任务能停自己的孩子…"`
  - `…_test.go:209: 父任务停自己的孩子被身份判定拒了（调用者读到的不是 task id）: "拒绝停止…"`
- 还原：hash `3ec8d498` 逐字回；`git status --porcelain -- internal/tools` 空

收束复跑：三枚 hash 复原为 `87ac0624 8c9266a7 3ec8d498`（与 A749 复核值逐枚一致），porcelain 空，`rc=0`、2 PASS。

## AC#3 不许放宽（读数）

本腿**只新增一枚 `_test.go`**（`039ec93c`：+227/−0），既有任何断言零删改；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）零碰；`cmd/wisp` 只读、一枚 `go test` 未对它跑；突变全部逐字还原（上面三发）。

## 没做的（具名）

1. **seed B `:206` 的上游 L2 判定通道**未追到底（登记为观测；`rc`/红句已留档，主红 `:194` 独立成立）。
2. `loop.go` 的 `CorrelationID: callCorr(...)` 突变**未跑**（`internal/agent` 在写面外；由本尺 `:218-226` 行断言按形状拦截——该断言在 corr==taskID 时不成立——但本程未取这发读数）。
3. 不跑整包（在册红会混淆；按名靶向）；不碰票面／台账／README；零 push。
