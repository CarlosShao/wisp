# 派单 2026-09-28 18:4x — 票 197 第一波（两枚写腿：197-r1 实体层／197-r2 流层）

编排者裁定（全部落 `docs/reports/pending-and-issues.md` 的 `A399`／`A401`）：
owner 09-28 18:1x 原话＝「**子代理这么重要的东西，说不要就不要了，还尼玛问问问，别人怎么做的，你就怎么做，别鸡吧天天问我……必须特么做完整功能，明白吗？**」
⇒ `PLAN.md:1544` 的 REJECTED「多 Agent 协作」与 `:1447` 的「S7 之前不得开放多任务并发」**已当场改判**，**不许再拿这两行当不做的理由**；
`PLAN.md` 与 `docs/specs/**` 的文字仍然一字不动（改判只落台账）。

两枚腿**并行**，写面互不相交（r1 只碰 `internal/tools/**`＋`cmd/wisp/**`＋新建文件；r2 只碰 `internal/panel/**`）。
两枚**共用同一份约定**（下面 §0），谁都不许自己改它。

---

## §0 两枚共用的约定（编排者已裁，不是留给你们选的）

- **子代理身份**：一枚子代理＝名册里一条**有自己 `taskID`** 的记录，带 `parentTaskID`（可空＝根任务）、`label`（标题，只放任务名）、`kind`（`root`／`subagent`）。
- **状态维只有一枚、且只能取 D43 那 20 枚名字**：`internal/statemachine/states.go:11-31`，写入必须过 `statemachine.Valid`（`:39`）。
  ⚠ **不许新造第二个状态字段、不许接 `memory.TaskLog.State` 那套旧词表**（`done/cancelled/running/succeeded` 四枚都不在 D43 表上，归票 196）。
- **流键形状写死**：`subagent:<taskID>`（全小写、冒号分隔、`<taskID>` 逐字透传不加工）。r2 按这个形状建流，r1 按这个形状喂流。
- **并发与深度**：同时在跑的子代理上限 **8**（根与其所有后代**共享同一个池**），嵌套**深度默认 1**（子代理不许再派子代理）。超限＝**硬拒并给可读理由**，不许静默排队到看不见的地方。
  依据＝DSH 现量（`packages/subagent/subagent/src/index.ts:201-203`、`continuation-activation.ts:45-50,604-611`）。
- **审批**：**子代理永不自批**。它遇到需要批准的动作，走**父任务同一个审批队列**上抛；**面板来源的"允许"出口一条都不许为它开**（`Q-49` 丙那批判据在场，`AGENTS.md` §1.2 硬禁）。
- **取消**：一枚 `context.CancelFunc`／`AbortSignal` 贯穿到底。**父取消不自动级联到子**——这条要**显式写在给模型读的那段话里**（DSH `tool-subagent-control/src/index.ts:75-79` 的理由：否则模型以为全停了）。
- **taint**：子代理返回的结论文本进父上下文前，**必须用现成的 C25 源名 `risk.SrcTaskOutput`（"task.output"，`internal/risk/provenance.go:92-94`）盖戳**。
  ⚠ **不许新造 `subagent.output` 这种名册外的源名**——那等于让子级输出走"没人认识它"的通道（D30 间接提示注入那条）。

## §1 三条共同的硬规矩

1. **锚点自取**：起手先 `git rev-parse --short HEAD`＋`git status --porcelain` 记下名册，**别信我这里的任何行号**——共享树里别人的提交会动它，行号漂了只算正常，**具名报出来即可**。
2. **只 commit、绝不 push**；commit 必须带**显式 pathspec**（一枚枚路径写全）；禁 `git add -A`/`git add .`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`restore`/`clean`/`worktree`。
   ⚠ 新文件要**先 `git add -- <那一条显式路径>`、再核 `git diff --cached --name-only` 只有它**，然后 commit（`-- <未跟踪文件>` 直接 commit 会报 `pathspec did not match`）。
3. **写面只有 `frontend/**` 与 `design/**` 是零写面**（读可以）。树里那 16 枚 `design/` 未提交删除、`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*` 都是别人的，**别碰别提交**。
   ⚠ **绝不在仓内跑任何删除命令**；临时件只建不删。

---

## 腿 A：`197-r1` 实体层（子代理作为任务真的存在，并且能派出一个）

**目标（这就是"完整功能"的第一层，不许砍）**：模型调一枚新工具 **`task.spawn`**，宿主就**派生出一枚子代理**（一份自己的 `agent.Loop`），
它跑完后**结论回到父任务**，同时它在 `TaskRoster` 里**是一行有状态、可查、可停的记录**。

**开工前先现读这几处**（都在树里，逐枚读到行）：
`internal/tools/task.go`（`TaskOutput` 结构、`TaskRoster.Record:123`/`Look:139`/`Count:152`、`BuiltinTaskEntries:345`、`TaskDeps`）
／`internal/agent/loop.go`（`New(opt Options):194`、`Run:332`、`RunAsync:321`、`run(ctx,taskID,input):339`）
／`cmd/wisp/run.go`（`rt.tasks = tools.NewTaskRoster()` 与 `TaskDeps{Roster:…, Paths:…}` 那两行、`agent.New` 的装配点）
／`internal/statemachine/states.go:11-39`／`internal/risk/provenance.go:84-94`。

**要做的事**（按序，每步都带判据）：
1. **名册扩身份**：`TaskOutput` 加 `ParentTaskID`／`Label`／`Kind` 三枚（`State` 那枚票 188 已进树，别重复加）。
   给名册加一枚**"列出某任务的全部后代"**的读口（界面那一层要点进去，靠它拿树）。
2. **派生入口 `task.spawn`**：新文件 `internal/tools/subagent_197.go`＋声明/注册；
   入参＝**`description`（短标题）与 `prompt`（自包含任务书）两枚分开**（minimax 现量 `agent-tools/src/desktop/builtin-defs.ts:541-547` 把它们写成两个字段，不许混用）；
   实现＝用**同一套宿主装配**构造一份子 `agent.Loop`（复用现有 `Options`，别另起一套运行时）、拿一个新 `taskID`、
   **先 `Record` 一行状态进名册、再起跑**（⚠ DSH/Step-Code 都自述过这个坑：不预先发布名册，"实时列表会在调用方最需要它的几十秒里保持不可见"，见 `subagent/execute.ts:376-382`）；
   跑完把结论 `Record` 回同一条、状态过 `statemachine.Valid`、结论文本按 §0 的 taint 规矩盖 `task.output` 源名后回给父任务。
3. **上限与深度**：按 §0 的 8／1 写死（**常量放一处具名**，别散在各文件），超限硬拒并给可读理由；子代理里拿不到 `task.spawn` 这一枚工具（深度＝1 的结构做法，比运行期计数更硬）。
4. **取消**：`task.spawn` 返回的句柄要能被 `task.cancel` 认（名册里存 cancel 函数或 ctx，别新造第二套停止路径）；**父取消不级联**这句要出现在给模型的说明文本里。
5. **审批**：子的审批请求走**父任务同一个队列**；证明它**没有**自己的允许出口（见判据 AC#5）。
6. **判据（`internal/tools/subagent_197_test.go`，每条都要正控）**：
   - 派生后名册里**确有**一行、且 `ParentTaskID`／`Kind`／`State` 三枚都是真值（⚠ 票 181 `AC#7` 那枚教训＝"字段恒 false/恒空"要能被判红，不许只测"结构体有这个字段"）；
   - 造 **9 枚**子代理 ⇒ 第 9 枚被硬拒且理由可读（正控）；8 枚以内全跑通（反控）；
   - 子级拿不到 `task.spawn`（正控：让子尝试派一枚 ⇒ 要失败并说明原因）；
   - 子结论带 taint 源，**去掉盖戳那一步的变异要能被判红**；
   - 取消语义：停父 ⇒ 子继续（这条要真测，别只写注释）；停子 ⇒ 只死那一行；
   - 子不许自带允许出口：造一枚"子自己批自己"的假腿 ⇒ **要红**。
7. **门禁**（交件前自己跑，读数逐名报）：`go build ./...`／`go vet ./internal/tools/ ./cmd/wisp/`／
   `go test -count=1 ./internal/tools/`／`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./cmd/wisp/`
   ⚠ **`cmd/wisp` 必须带那枚 PATH**，否则会拿到 `exit status 0xc0000135`（缺 DLL，**根本没跑任何用例**）——那不是红，那是仪器瞎了（`A400` 刚记过）。
   ／`gofumpt -l` 空／`sh scripts/d22scan.sh` 判"clean 或逐名"。
   ⚠ **`internal/observe/nobarego_test.go` 在场**：裸 `go func(` 无 owner/recover 会红，起 goroutine 走它认的那一形。
8. **证据件**：`docs/evidence/s1/197-subagent-entity-r1.md`（起手名册／落点逐跳到行／判据逐条与正控变异读数／门禁改前改后／**没测到什么**清单）。

**禁区**：`internal/panel/**`（腿 B 的地盘，**这枚文件都别开**）；`frontend/**`／`design/**`；`PLAN.md`／`docs/specs/**`；`internal/risk/**` 不许改（只用现成源名）；
`thresholds.go`／golden／`allowlist.txt`；C17 白名单既有名字（`internal/panel/bridge.go` 那 4 枚入向方法**不许动**，本腿不接界面）。
**不许新增 D34 之外的第二枚工具**（只加 `task.spawn` 这一枚；`task.list`／`task.cancel`／`task.output` 是 D34 现表里的三行，别新造同类名）。

## 腿 B：`197-r2` 流层（每枚子代理一条自己的流，溢出别再"合并着吞"）

**现量（起手自己复算一遍，别信我写的行号）**：`internal/panel/pump.go` 里 `StreamLog.Append(key,text)`（约 `:301`）、`Chunks()`（约 `:339`）、
`DefaultStreamKeys = 32`（约 `:275`），**溢出语义＝"合并而不是丢"（约 `:287-290`）**。
⚠ 对子代理这一形**是错的**：两枚子代理的输出被并成一条＝**界面在说谎**（owner 要的就是"点某枚看到它自己的流"）。
⚠ 我跑过预检：**`DefaultStreamKeys`/`maxKeys` 在全仓测试里零命中 ⇒ 没有既有钉，改它不会有人替你响**，所以正控必须你自己造。

**要做的事**：
1. **重裁溢出语义**（二选一，选完具名说清为什么）：甲＝上限随子代理数走（键集合可增长，按活跃任务数封顶）；
   乙＝到上限时**显式标"已截断"并保留每个键各自的头尾**，绝不把两条流并成一条。**不许选"把上限调大就完事"**（那只是把同一个谎言推后一轮）。
2. **每子一条流**：按 §0 写死的键形状 `subagent:<taskID>`；**这条腿不需要 r1 的代码**，只要按形状喂与读。
3. **判据（`internal/panel/subagent_stream_197_test.go`）**：
   - 造 **N＋1 枚**子代理（N＝当前上限）⇒ **没有任何一条被并进别人的**（这是票 197 第 2 层那发结构性判据，必须真测）；
   - 正控：故意让两枚子代理同时写 ⇒ 两条流各自完整、互不污染（若你的修法选了"合并"分支，这条要红）；
   - 截断／增长那一支要有一枚能判红的尺（按你选的甲或乙写），**并具名说明"不修会怎样、修好什么不会坏"**；
   - ⚠ **不许把 `pump_test.go:111-124` 与 `:270-276` 那两枚"快照四枚字节键"钉改松**——它们与本腿无关（本腿动的是 `StreamLog`，不是 `Snapshot` 顶层键集），**动了就是违规**（解冻口径见 `A388`/`A389`）。
4. **门禁**：`go build ./...`／`go vet ./internal/panel/`／`go test -count=1 ./internal/panel/`
   ⚠ 这包**今天有 3 枚在册常红**（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）——**逐名报，不当绿、不修、不 Skip**，它们归界面那侧与色表那格。
   ／`gofumpt -l` 空／`sh scripts/d22scan.sh`。
5. **证据件**：`docs/evidence/s1/197-subagent-stream-r2.md`。

**禁区**：`internal/tools/**`／`cmd/wisp/**`（腿 A 的地盘，尤其 `cmd/wisp/run.go` 那两行它正在动）；`Snapshot` 顶层键集（这枚腿**不加键**，载体那一层是下一程 197-r3）；
`internal/panel/tokens_fourway_test.go` **一字不动**；`frontend/**`／`design/**`；`PLAN.md`／`docs/specs/**`。

---

## 两枚共同的交件要求

- **工具调用预算 30 枚**（超了具名报"超了几枚、花在哪"，**不许为此放宽任何判据或断言**）；骨架（首条产码前的现读）争取 ≤8 枚。
- **交件正文必须包含**：① 起手锚点号与名册枚数 ② 落点逐跳到 `文件:行号` ③ 判据逐条的**改前/改后读数**（含正控那几发）
  ④ 门禁逐包读数（红就逐名）⑤ **被拒调用枚数** ⑥ **"我没测到什么"**（至少 3 条，写具体的）。
- 落不了地的地方**当场报回来换形状**，不要自己砍需求、也不要停下来等 owner——**owner 这一程不参与**（他那句话就是"别问我"）。
