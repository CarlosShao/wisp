# 221-r1 落点草案（产码腿，起手即落，骨架先于产码）

锚点：本腿现跑 `git rev-parse --short HEAD` = `a7993a9b`（dev）。
批准依据：台账 `A434`（甲形＝注册 `task.cancel`；乙形三处删句与甲形同一发）。撤销口令「撤 221 甲」。
票面权威＝`.scratch/wisp/issues/221-task-spawn-description-promises-task-cancel-that-is-not-registered.md`
含文末「收 `221-c1`」那节的判据改写（AC#1 分母＝**这次装配真正注册出去的所有工具名的并集**，
不是单枚 `BuiltinTaskEntries`；AC#4＝本票自己那枚 DEFERRED 标记的两条，不承诺与 `SPEC-12 §5` 双向 1:1）。

---

## 0. 撞钉预检（动手前现跑，逐枚读了断言原文）

| 尺 | 读数 | 影响 |
|---|---|---|
| `grep -rn 'task\.cancel' --include='*.go' internal/ cmd/` | 改前 **4 处**：`subagent_197.go:196`（说明书许诺）、`task.go:23`（DEFERRED 标记）、`task.go:576`/`task.go:593`（注释指 D34 那行的 `—` 与"为何不在名册"） | 四处都是本票射程内要改/要核对的文本 |
| `grep -rn '\.Cancel(' --include='*.go' internal/ cmd/ \| grep -v _test` | 改前 **8 处命中，零枚是 `TaskRoster.Cancel`**（全是 `RunningTask.Cancel`／`root.Cancel`／`feedRoot.Cancel`／`replyRoot`／`reloadRoot`）＋`subagent_197.go:39` 的一句注释指它 | AC#2 的分母基线＝**0 枚生产调用者** |
| `internal/tools/ticket175r2_stamp_live_test.go:343` `TestEveryRegisteredToolIsClassifiedForMarking175r2` | **枚举 `BuiltinFSEntries`＋`BuiltinTaskEntries` 的并集**，要求每一枚名字都在 `classified175r2` 表里；新名字的分支逐字写着「新工具请就地回答」 | 注册 `task.cancel` **必然**要求我在该表加一行分类答语——这是答问，不是放宽断言 |
| `internal/tools/subagent_197_test.go:419` `Test197SpawnDescriptionNamesTheRealPoolCap` | 只断 `"上限 %d 枚"` 与 `"深度 %d"` | 删句不碰它 |
| `internal/tools/subagent_197_test.go:797` | 断 `"停掉父任务不会级联"`（另一枚字面量） | **必须原样保住** |
| `internal/tools/subagent_197_test.go:562`／`:747` | 断 `ctx.Done()` 分支文本含 `"没有被级联取消"` | 改 `:389` 那半句时保住这一枚 |
| `internal/tools/subagent_197_test.go:871` | `subagentDepsFieldNames()` 断 `len == 5` | **不碰 `SubagentDeps`**（A434 ① 同句） |
| `internal/tools/subagent_197_test.go:722` `Test197CancelIsPerRowAndNeverCascades` | 已经在 `roster.Cancel` 上测过"停一枚只动那一行/重复取消报没有在跑/不存在的报查不到" | 我的工具复用同一条 `Cancel`，不重造语义 |
| 三枚冻结件 `internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go` | 本腿新增/改名符号（`taskCancel`、`TaskCancelDecl`、`task.cancel` 字样）在 `--include='*.go'` 全量 grep 里**零命中**它们的断言 | 一字不动，也不被顶 |
| `NoGate.PendingWindow`（`internal/tools/gate.go:138`） | **返回 `AnswerReject`**：「L1 确认窗口尚未接入（票 21），已拒绝执行」 | 若 `task.cancel` 按 D34 声明 L1，则**所有用 NoGate 的既有 harness 里它根本进不到 Execute** ⇒ 我的正控用例必须自带一枚会放行的窗口答复器，否则红绿都不是判据 |
| `PLAN.md:2564`（只读） | `task.list / task.cancel | 查看/取消任务 | **L0 / L1** | — | S7` | 名别＝`task.cancel`、级别＝**L1**（`task.list` 那半支仍是 DEFERRED，本票不注册它）、能力列＝`—` ⇒ `Capabilities/Needs = nil`、`PathParams = nil` |
| `docs/specs/SPEC-07-tools-and-plugins.md:70`（只读） | 同一行镜像 | 同上；两份冻结文字**一字不动**（A434 ⑤） |

## 1. 打算动的文件（逐枚写清为什么）

1. `internal/tools/task.go`
   - 加 `taskCancel` 工具（`Name()` 回 `task.cancel`）＋ `taskCancelDecl()`（**不导出**，见 §4 导出名纪律）＋ schema（只有 `task_id` 一枚参数）。
   - `BuiltinTaskEntries` 从 1 枚变 2 枚（`task.output` ＋ `task.cancel`）。
   - 头部 `:18-27` 那段把 `task.cancel` 标为 DEFERRED 的文字**与接线同发摘掉**（AC#4a/4b），`task.list` 那行**一字不动**。
   - 权力边界写在 Execute 里：只有**父任务**能停自己的孩子；停自己、停兄弟、停别人家的孩子一律拒；
     root 行没有父 ⇒ 也拒（"用户能停"那一支按 A434 ⑤ 具名归票 181／票 220，不在本票）。
   - 回执文本必须"写给模型看"：停了谁、没停谁、父取消不级联这一维反向也一样、那一行不会消失。
2. `internal/tools/subagent_197.go` — 乙形三处（`:196`、`:197` 尾巴那个括号、`:389`）**同一发**改：
   删掉无条件的"可以单独停它"许诺，换成当下为真的话（只有派生它的父任务能停它），
   并保住 `"停掉父任务不会级联"` 与 `"没有被级联取消"` 两枚字面量。
3. `internal/tools/ticket175r2_stamp_live_test.go` — 在 `classified175r2` 里为 `task.cancel` 就地答一句（那枚普查尺要求的形状，不是放宽）。
4. **新** `internal/tools/task_cancel_221_test.go` — AC#1 常驻能力尺＋AC#2 生产调用者与终态／审计＋AC#3 两枚正控。
5. `cmd/wisp/run.go` — 只改 `:420-431` 那三行注释（它现在写着"只有 task.output 在这一族"，接完线就成了谎）；
   注册本身由 `BuiltinTaskEntries` 自动带上生产装配路径，**不新增调用行**。

不动：`frontend/**`、`design/**`（零读零写）、`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、
`allowlist.txt`、`SubagentDeps` 字段集、C17 白名单。

## 2. AC 逐格落点（UNJUDGED 起点，完成后逐格改状态）

- AC#1 **UNJUDGED**：常驻尺＝新用例 `Test221EveryPromisedTaskNameIsRegistered`。
  分母＝`BuiltinFSEntries`＋`BuiltinFSWriteEntries`＋`BuiltinTaskEntries`＋`BuiltinSubagentEntries` 的并集；
  分子＝每枚 `Description()`＋`Parameters()` 文本里 `task\.[a-z_]+` 词根。
  ⚠ **未修码读数先跑**：本腿先落测试、后落产码，改前这把尺必须恰好红一枚＝`task.cancel`（红句要点名是哪一句、缺哪一枚）。
- AC#2 **UNJUDGED**：尺＝`grep -rn '\.Cancel(' --include='*.go' internal/ cmd/ | grep -v _test` 里属于 `TaskRoster` 的枚数，改前 0 ⇒ 改后 ≥1（`taskCancel.Execute`），
  且 `cmd/wisp/run.go:432` 那条装配路径真注册它（另加一枚装配侧用例证明从 bridge 能 dispatched 到）。
  终态：停掉的孩子那一行 state 落 `Muted`（借名，票面写明）、行不消失、`ParentTaskID` 还在、它自己的流 `Close` 掉。
- AC#3 **UNJUDGED**：两枚正控都走真桥真名册真 task id——
  ① 以"孩子"的身份停兄弟 ⇒ 拒；② 以"孩子"的身份停自己 ⇒ 拒；两枚都要证明**没动到目标行**（句柄还在、状态还在）。
  ＋一条钉住父取消不级联（保 `:334` 那形＋`:197` 那句原样）。
- AC#4 **UNJUDGED**：(a) `task.go` 头部那段 `task.cancel` 的 DEFERRED 文字随接线同发摘掉；
  (b) 同一 commit 内完成；尺＝摘标记后 `grep -n 'DEFERRED' internal/tools/task.go` 只剩 `task.list` 那一行。
  ⛔ 本票不承诺与 `SPEC-12 §5` 双向 1:1（票面改写⑤：§5 里没有 `task.*` 行，那枚仪器归票 225）。
- AC#5 **UNJUDGED**：整包 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` 到终态，
  名册尺＝`grep -cE '^[[:space:]]*--- (PASS|FAIL)'` 两层分开报，分母对拉 `=== RUN`；
  失败包尺 `grep -P '^FAIL\t'`。已知在册红（`internal/panel` 4＋`internal/ball` 1）不算本腿新增。
  ⚠ 今晚 CI 同机跑（`slo-full` 上一回 1h40m）⇒ 计时类红约 20:30 前一律〔待复量〕，安静复量三发再归因。

## 3. 已知的"我不做"那一格（不许变隐形缺格）

- **"用户能停子代理"**：A434 ⑤ 具名归票 181／票 220，本票射程外。
- **`task.list`**：仍 DEFERRED，本票一字不动它的标记。
- **"被谁停"的上屏字段**：A434 ④ ⛔ 不新增；第一版只走**现成审计 sink**
  ＝桥每调用一行 `tools: call ... task=<调用者> tool=task.cancel`（`Options.Logf`＝`rt.auditf`）
  ＋ `tool_call` 行 `ArgsJSON`（记目标 id）。本腿要用用例证明这条审计链真落地，而不是口头说"审计可查"。
- **`task.spawn` 等孩子的时限那一半**（A434 末段）＝票 222 AC#3 剩半格，两支都要另批，不在本票。

## 4. 导出名纪律（要不要另落 `A##`）

**本腿刻意不新增任何导出标识符**：工具类型 `taskCancel`、schema 变量 `taskCancelSchema`、
Decl 构造函数 `taskCancelDecl`（小写）全部不导出，`BuiltinTaskEntries` 的签名不变
（只多返回一枚 Entry，向后兼容、不改字段、不改函数名）。
若实现中撞到必须新增的导出名 ⇒ **停手上报**，写进"没做完"，不自造。
