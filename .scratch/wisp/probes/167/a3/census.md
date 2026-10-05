# 167-a3 — 占用条那一枚的第 0 跳：只读取证（⛔ 只做占用；序号／停止／草稿／崩溃自救一律未碰）

> 派单：只读腿 `167-a3`，工作目录 `D:\work\workspace\projects plans\Wisp`，分支 `dev`，**共享工作树**。
> 票面：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md`（**AC 框一枚未碰**）。
> 本件只答四问（Q1 第 0 跳名册／Q2 分子分母真源／Q3「未知」该长什么样／Q4 撞钉名册），**范围＝票 167 AC#2 占用那一枚**。
> ⛔ 硬闸照单执行：**零 `go` 命令**（不 build／test／vet／gofumpt／run／scripts）；**未读 `frontend/**` 与 `design/**`**，本件不转述那两棵树的任何内容（凡"页面侧声明了什么键"一律写成注释自述或〔预测〕）；未改任何已跟踪文件；未碰 `docs/reports/pending-and-issues.md`。
> 前一腿读数只当线索：`167-c2` §1 与 §6.1（本件 §1.0 逐条复量并给出"仍成立／已移动"的判语）。所有 `file:line` 均为**本件现量**，⛔ 无一枚抄自前趟（c2 引的 `loop.go:253` 今天已是 `:260`，见 §1.0）。

## 0. 起手锚（逐字读数）

```
date "+%H:%M:%S%z"        -> 11:01:29+0800   （起手第一枚；下文 §5 各尺另记各自时刻）
git rev-parse --short HEAD -> c6cf66e6
git rev-parse --abbrev-ref HEAD -> dev
git log -1 --format='%H %ad' --date=iso
                         -> c6cf66e64849955bf92a4b3356c096ea81c35563 2026-10-05 10:58:27 +0800
git status --porcelain cmd internal tools docs scripts
                         -> （空，11:11:29 复量同值）
```

起手读数解读：产码五根（`cmd`/`internal`/`tools`/`docs`/`scripts`）在锚 `c6cf66e6` 上全部干净 ⇒ 本腿读到的字节即 HEAD 字节。本件唯一写面＝`.scratch/wisp/probes/167/a3/census.md`（目录 `a3` 由本腿 `mkdir -p` 建出，此前不存在）。

⚠ 与本腿写面**不相交**的他程遗留：`.scratch/wisp/probes/167/` 下已有 `a1`／`a2`／`c2`／`c3`／`r1`（`r1` 是空目录）；本腿一律未写、未改、未提交它们。

## 1. Q1 — 第 0 跳名册：从"占用真正被算出来的那一处"到 `internal/panel` 的快照构造

### 1.0 先复核 `167-c2` 那句"loop 是 execute 局部变量"（派单点名要复量，不许信它说的）

**判语：这句在今天仍然成立，且它比 c2 写得更死——但 c2 的三枚行号全部已经移动。**

现量（本腿，锚 `c6cf66e6`，11:11:29）：

- `cmd/wisp/run.go:990` 逐字 `func (rt *agentRuntime) execute(parent context.Context, task string) int {` —— 签名**未变**（c2 引 `run.go:990` 同值）。
- `cmd/wisp/run.go:1023` 逐字 `loop, err := agent.New(opts)` —— `loop` 是 `execute()` 的**函数内局部变量**。c2 §1.4 把它记在 `run.go:1023`（同值），但 §6.1 又记在"一带"，本件钉死为 `:1023`。
- `internal/agent/loop.go:260` 逐字 `func (l *Loop) Budgets() Budgets { return l.b }` —— ⚠ **c2 引的是 `loop.go:253`，今天已移到 `:260`**（差 7 行）。
- `internal/agent/loop.go:263` 逐字 `func (l *Loop) History() []llm.Message {` —— c2 引 `:256`，同样已移动。
- **`agentRuntime` 全字段过目**（`cmd/wisp/run.go:266-377`，本腿逐行读完）：结构体里**没有任何字段是 `*agent.Loop`**（尺：`grep -rn "agent.Loop" --include=*.go cmd | grep -v _test.go` ⇒ 全树只命中 `cmd/wisp/run.go:21` 的**注释**一行，零字段声明）。⇒ c2 那句"`runtime` 里没有任何字段握着它"**成立**。
- `execute()` 的产码调用者＝**2 枚**：`cmd/wisp/run.go:262` 逐字 `return rt.execute(s.taskCtx, task)`（`wisp run` 腿）与 `cmd/wisp/resident_task_source_windows.go:462` 逐字 `code := src.run.execute(root.Ctx, text)`（常驻腿）。每调一次铸一枚新 `loop`，返回即不可达。
- 唯一"握着装配但**不**握着环路"的字段：`cmd/wisp/run.go:292-293` 逐字 `loopOpt    agent.Options` / `loopOptSet bool`，写在 `cmd/wisp/run.go:1022` 逐字 `rt.loopOpt, rt.loopOptSet = opts, true`，读在 `cmd/wisp/run.go:789` 逐字 `BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet },`。**这是值拷贝的 `Options`，不含 `*Loop`、不含 `Budgets`、不含 `history`。**

⚠ 本腿补一枚 c2 没量的对偶事实：**子代理那一侧的环路连 `Options` 拷贝都不留**。`internal/tools/subagent_197.go:319` 逐字 `child, err := agent.New(opt)`、`:344` 逐字 `bg := child.RunAsync(childCtx, prompt)`，两处都是 `taskSpawner` 方法内的局部变量，`internal/tools` 里没有任何字段握着 `child` ⇒ **子代理那一发任务的上下文用量今天从任何地方都读不到**（连"局部变量"这一格都没有）。占用条若只报根环路，必须在出口里具名说"只覆盖根任务"（归 §8 判语三行第 1 行）。

### 1.1 名册（①→⑨，一跳一行；类型标记按派单四分类）

分类图例：**①已有载体**／**②已有载体但传的是零值**／**③根本没有载体（要新加字段／参数）**／**④有载体但装配根没连线**。

| # | 跳（全路径 `file:line` ＋ 被指字符逐字） | 是什么 | 类别 |
|---|---|---|---|
| ① | `internal/agent/compress.go:110` 逐字 `func (c *Compressor) TotalTokens(hist []llm.Message) int {`（体内 `:112-115` 累加 `ApproxTokensOf(m.Content)`） | **分子的真正算处**：整段历史的 token 估算（单位＝token，算法＝`internal/agent/budgets.go:129` 逐字 `func ApproxTokens(s string) int { return len(s) / 4 }`，即 utf-8 字节／4） | ③ 对包外不可见：`Compressor` 挂在 `internal/agent/loop.go:189` 逐字 `comp     *Compressor`，**小写字段、无导出读口**（`grep -rn "^func (l \*Loop) [A-Z]" internal/agent/*.go` ⇒ 导出面只有 `Budgets/History/Reset/Steer/RunAsync/Run/AttachProjectInstructions/ProjectInstructionManifest` 八枚，无 `Occupancy`／`Used`） |
| ② | `internal/agent/loop.go:399` 逐字 `hist := l.History()` ＋ `internal/agent/loop.go:400` 逐字 `if l.comp.Need(hist) {` | 环路每轮**真的算了一次**分子（`compress.go:119-120` 逐字 `func (c *Compressor) Need(hist []llm.Message) bool {`／`return c.TotalTokens(hist) > c.b.HistoryCompressTokens`），但**结果只用作布尔分流，那个整数没有离开 `run()`** | ② 算了、丢了（丢弃点在 `loop.go:400` 的 `if` 条件里，整数没有变量接住） |
| ③ | `internal/agent/compress.go:69-70` 逐字 `TokensBefore        int`／`TokensAfter         int`（`CompressionReport` 的两枚字段，产在 `compress.go:175` 与 `:212`） | **唯一被带出 `run()` 的 token 整数对**，但它只在压缩真跑过一次时才有值（`internal/agent/loop.go:409-411` 逐字 `} else if rep.Ran {`／`l.replaceHistory(nh)`／`res.Compression = rep`） | ① 载体在（`internal/agent/loop.go:100` 逐字 `Compression CompressionReport` 是 `Result` 的一枚字段），**且 `res` 在 `execute()` 里拿得到**（`cmd/wisp/run.go:1106` 逐字 `res := bg.Wait()`）；⚠ **`res` 又是 `execute()` 的局部变量、且是任务跑完之后才有** ⇒ 对"进行中"的读数它等于 ②（值多半是零值，且时序上晚了半拍） |
| ④ | `internal/agent/budgets.go:54` 逐字 `ContextWindow int`（`Budgets` 的一枚字段）；产在 `internal/agent/budgets.go:85-88` 逐字 `func BudgetsFor(ctxWindow int) Budgets {`／`	if ctxWindow <= 0 || ctxWindow > ReferenceContextWindow {`／`		ctxWindow = ReferenceContextWindow` | **分母**（缩放后窗口）。⚠ 这一枚**永远非零**：未知（≤0）与大于 128k 都被折成 `internal/agent/budgets.go:18` 逐字 `const ReferenceContextWindow = 128000` ⇒ **拿它当"分母已知"是假读数**（详见 §2.2） | ③ 对 `cmd/wisp` 而言：读口是导出的（`internal/agent/loop.go:260` `Budgets()`），但**握着它的只有 `execute()` 的局部变量 `loop`** |
| ⑤ | `internal/agent/loop.go:218` 逐字 `window := opt.Config.ContextWindow`／`:220` 逐字 `window = info.MaxContextWindow` | 分母的两级回退（配置值 → provider 报的窗口），发生在 **`New()` 里、环路存在之前** | ① 上游载体已在 runtime：`cmd/wisp/run.go:302` 逐字 `endpoint llm.Endpoint`（赋在 `cmd/wisp/run.go:441` 逐字 `rt.endpoint = ep`），`llm.Endpoint.ContextWindow` 声明在 `internal/llm/resolver.go:56` 逐字 `ContextWindow int` ⇒ **`rt.endpoint.ContextWindow` 在 `execute()` 之外今天就可读，零新连线**；`rt.provs[0].Info().MaxContextWindow` 同样可读（`internal/llm/provider.go:51` 逐字 `MaxContextWindow int`） |
| ⑥ | `cmd/wisp/run.go:1012` 逐字 `ContextWindow: ep.ContextWindow,` | 装配根把**未缩放的原始窗口**递给环路；这就是第 0 跳的**上游端点** | ①（同一枚值的载体在 `rt` 上，见 ⑤） |
| ⑦ | `cmd/wisp/run.go:1023` 逐字 `loop, err := agent.New(opts)` ＋ `cmd/wisp/panel_pump.go:93` 逐字 `func (rt *agentRuntime) setInstructionLoader(l *projctx.Loader) {`（调点在 `cmd/wisp/run.go:1070` 逐字 `rt.setInstructionLoader(instrLoader)`） | **第 0 跳本身今天不存在**：占用没有 `execute()→runtime` 的交接口。但**同型的先例已经存在两枚**：`rt.setInstructionLoader`（runtime 字段 `cmd/wisp/run.go:361-362` 逐字 `instrMu     sync.Mutex`／`instrLoader *projctx.Loader`，读口 `cmd/wisp/panel_pump.go:108` `instructionBundle()`）与 `rt.loopOpt`（`:1022`）。两枚都是"execute() 里造出对象 → 存进 runtime → 泵经由闭包读"的形状 | ③ **要新加**：`agentRuntime` 需要一枚 `*agent.Loop`（或等价的读数载体）字段＋互斥；⛔ 不需要新增任何 `Loop` 方法（分子可由 `History()`＋包外同规则函数 `internal/agent/budgets.go:150` 逐字 `func ApproxTokensOfMessage(m llm.Message) int { return ApproxTokensOf(m.Content) }` 在 `cmd/wisp` 侧复算） |
| ⑧ | `internal/panel/pump.go:138` 逐字 `type PumpSources struct {`（体内**没有**任何占用读口；`grep -rniE "occupanc|contextused|usedtokens|windowused|prompttokens" --include=*.go cmd internal tools` 排除 `_test.go` ⇒ `cmd`／`internal/panel`／`internal/agent` **零命中**，命中的只有 `internal/llm/openaichat/wire.go`、`internal/projctx/projctx.go`、两条注释与 `tools/mockllm`） | 泵的读口名册 | ③ `PumpSources` 要新加一枚 reader 槽（先例：`internal/panel/pump.go:177` 逐字 `Instructions func() *projctx.Bundle`、`:190` 逐字 `Tasks func() TaskRosterState`）。⚠ 装配点**全仓只有一枚**：`cmd/wisp/run.go:699` 逐字 `rt.pump = panel.NewSnapshotPump(panel.PumpSources{`（尺：`grep -rn "NewSnapshotPump" --include=*.go .` 排除 `_test.go` ⇒ 产码只此一处） |
| ⑨ | `internal/panel/pump.go:234` 逐字 `func (p *SnapshotPump) Snapshot() Snapshot {`（填法先例 `:304-309`／`:310-328`）→ `internal/panel/composer.go:57` 逐字 `type Snapshot struct {`（6 枚字段：`:58 Pending`／`:59 Results`／`:60 Composer`／`:61 GeneratedAt`／`:74 Instructions *InstructionsSection \`json:"instructions,omitempty"\``／`:91 Tasks *TaskRosterSection \`json:"tasks,omitempty"\``）与 `internal/panel/composer.go:235` 逐字 `type ComposerState struct {` | 快照构造与**载体字段名册** | ③ 顶层与 composer 段**都没有占用格**（`ComposerState` 14 枚 JSON 键名册＝`mode`/`workspace`/`attachments`/`acceptedAttachmentMimes`/`maxAttachmentBytes`/`attachmentError`/`git`/`currentModel`/`modelKnown`/`credentialState`/`credentialKnown`＋`Git GitView`/`Credential CredentialState` 两枚结构体字段，逐枚过目 `composer.go:236-269`，无 used/total 任何一格） |
| ⑩ | `internal/panel/pump.go:337` 逐字 `func (p *SnapshotPump) Marshal() ([]byte, error) {` → `:348` 逐字 `func (p *SnapshotPump) Publish() (Snapshot, []byte, error) {` → `cmd/wisp/panel_pump.go:319` 逐字 `func (rt *agentRuntime) bookPanelSnapshot(snap panel.Snapshot, data []byte) error {` | 出口。今天唯一生产落点＝账本一行摘要（`cmd/wisp/panel_pump.go:331` `panelSnapshotSummary`）＋内存 `rt.lastSnap`（`cmd/wisp/run.go:369` 逐字 `lastSnap      panel.Snapshot`）；**transport 不存在**（`internal/panel/pump.go:15-23` 自述，本腿未改其一字） | ① 出向通道已有（与占用无关，不动它） |

### 1.2 ④ 类（有载体但装配根没连线）今天在册的三枚，逐条判归属

派单要求"标它是④"。占用这一枚**没有任何一格属于④**——它连①都没有；但同一棵树上另有三枚④，占用条的落地会顺手经过其中一枚，必须具名分清归属，免得下一腿把别人的断口当自己的债背走：

1. `internal/panel/pump.go:200` 逐字 `L1Windows func() []L1WindowWait` —— 槽在、`pump.go:323-325` 有消费，装配根**零连线**（尺：`grep -rn "L1Windows:" --include=*.go .` ⇒ 唯一命中 `internal/panel/subagent_blocked_220_test.go:81`，产码 0 枚）。**归票 220，不属本枚。**
2. `internal/panel/pump.go:190` 逐字 `Tasks func() TaskRosterState` —— 已连（`cmd/wisp/run.go:724` 逐字 `Tasks: rt.taskRosterState,`），但它载的**状态维**今天恒空：`internal/tools/task.go:355` 逐字 `func (r *TaskRoster) MarkRoot(taskID, label string) {` 的注释 `:351` 逐字 `// State is deliberately NOT written here: the root task's producer still files` ⇒ 根行走 `internal/tools/task.go:158` 逐字 `func (o TaskOutput) StateAnswer() (statemachine.State, string) {` 的"没登记"那一支。**这就是 c2 说的"空闲与占用不可区分"那一维，归票 196／`A394`，⛔ 本枚不许顺手改。** 占用条**不得**借这一维表达"有没有一发在跑"。
3. `cmd/wisp/resident_task_source_windows.go:167` 逐字 `running bool`（写 `:444` 逐字 `src.running = true`、`:474` 逐字 `src.running = false`；判据 `:440` 逐字 `if src.running {`）—— 常驻腿**有**一枚真 busy 位，但它是 `residentTaskSource` 的字段、**`wisp run` 腿没有对应物**，且没有任何快照 reader 读它。⇒ 这枚是"占用条能不能说'此刻没有任务在跑，所以占用未知'"的真上游，但它的**接线属序号／停止那一腿的写面**（`src.running` 与 `src.seq` 同族，票面 §9 第 5 条已把 r2 排成"占用／序号／停止一起看得见"）。**本枚只登记，不动手。**

### 1.3 一句话结论（Q1）

**第 0 跳＝⑦＋⑧＋⑨ 三处新建，一跳都省不掉**：`execute()` 里没有把环路（或它的两个读数）交到 runtime 的任何一行 ⇒ 今天直接给 `Snapshot`／`ComposerState` 加快照字段，读到的**必然**是恒空。分子的算处（①）在 `internal/agent` 包内、结果被 `if` 吃掉（②）；分母（④）在包内可读但值本身被 `BudgetsFor` 折成了"永远已知"；**唯一今天已经 crossed 出 `execute()` 的 token 整数对是 `projctx.Bundle{BudgetTokens, UsedTokens}`，而它是项目说明那一格的子预算、不是对话窗口**（详 §2.3，那枚映射在 `internal/panel/instructions_200.go:158-194` 里把两个数都丢了）。
