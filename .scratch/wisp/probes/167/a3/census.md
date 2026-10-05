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

## 2. Q2 — 分子与分母各自的真源

> 派单原话："分子＝这次对话目前占了多少（谁算的、什么单位），分母＝总共可用（谁给的）"。本节**逐枚列候选并具名写它今天能不能在 `execute()` 之外被读到**；拿不到的写"今天不存在"，⛔ 不写"应该有"。

### 2.1 分子候选（8 枚，全在盘上，逐枚现量）

| # | 候选 | `file:line` ＋ 被指字符逐字 | 谁算的／单位 | execute 之外读得到吗（判语＋凭据） |
|---|---|---|---|---|
| N1 | **整段历史的 token 估算**（＝票面语义上唯一真"占用分子"） | `internal/agent/compress.go:110` 逐字 `func (c *Compressor) TotalTokens(hist []llm.Message) int {`（体内 `:113` 逐字 `n += ApproxTokensOf(m.Content)`） | `Compressor` 算的；单位＝**估算 token**（utf-8 字节／4，`internal/agent/budgets.go:129` 逐字 `func ApproxTokens(s string) int { return len(s) / 4 }`） | **读不到（对象私有）**。`Compressor` 只挂在 `internal/agent/loop.go:189` 逐字 `comp     *Compressor`（小写字段），`Loop` 的导出面上没有一枚报数方法（尺：`grep -rn "^func (l \*Loop) [A-Z]" internal/agent/*.go` ⇒ 八枚：`Budgets/History/Reset/Steer/RunAsync/Run` ＋ `prompt.go:280 AttachProjectInstructions` ＋ `prompt.go:288 ProjectInstructionManifest`）。**但同一枚量可以在包外复算**：入参 `hist` 有导出读口 `internal/agent/loop.go:263` 逐字 `func (l *Loop) History() []llm.Message {`，算子有导出面 `internal/agent/budgets.go:150` 逐字 `func ApproxTokensOfMessage(m llm.Message) int { return ApproxTokensOf(m.Content) }` ⇒ 复算不需要动 `internal/agent` 一字；⚠ 前提仍然是"有人握着那枚 `*Loop`"，即 §1.2 的第 0 跳 |
| N2 | 压缩触发布尔 | `internal/agent/compress.go:119-120` 逐字 `func (c *Compressor) Need(hist []llm.Message) bool {`／`	return c.TotalTokens(hist) > c.b.HistoryCompressTokens` | 同一个算子的布尔半 | 读不到。产码唯一调用点 `internal/agent/loop.go:400` 逐字 `if l.comp.Need(hist) {` ⇒ **整数在这一步被算出来又被 `if` 吃掉，没有变量接住**（§1.1 跳②，类别 ②） |
| N3 | 压缩前后整数对 | `internal/agent/compress.go:69` 逐字 `TokensBefore        int`、`:70` 逐字 `TokensAfter         int`；产点 `compress.go:175` 逐字 `rep := CompressionReport{TokensBefore: c.TotalTokens(hist)}`、`compress.go:212` 逐字 `rep.TokensAfter = c.TotalTokens(out)` | `Compressor.Compress` 算的；单位＝估算 token | **`execute()` 内读得到、`execute()` 外读不到**。载体是导出的：`internal/agent/loop.go:100` 逐字 `Compression CompressionReport`（`Result` 的字段），`res` 在 `cmd/wisp/run.go:1106` 逐字 `res := bg.Wait()` 之后才存在。⚠ 两枚硬伤：① **只在压缩真跑过时才填**（`internal/agent/loop.go:409-411` 逐字 `} else if rep.Ran {`／`				l.replaceHistory(nh)`／`					res.Compression = rep`），没压缩过＝零值；② 时序在任务**结束之后** ⇒ 画不出"进行中"。零值当读数＝票面 AC#2 点名的退化形状 |
| N4 | 跨轮累加的已用 token | `internal/agent/guard.go:132` 逐字 `func (g *Guard) TokensUsed() int { return g.tokensIn + g.tokensOut }`；累加在 `internal/agent/guard.go:149-150` 逐字 `	g.tokensIn += u.InputTokens`／`	g.tokensOut += u.OutputTokens` | `Guard`；单位＝**计费 token（provider 报的真值），但是逐轮相加** | 读不到，而且**语义就不是占用**：`Guard` 是 `run()` 的局部（`internal/agent/loop.go:374` 逐字 `guard := NewGuard(l.guard)`）。⚠ 每一轮的 `u.InputTokens` 都含整段前文 ⇒ 第 k 轮这个数里同一段文字被算了 k 次 ⇒ **多轮任务上它可以合法地大于窗口**；拿它除以分母会画出一根 300% 的条 |
| N5 | 成本口径的 usage | `internal/agent/cost.go:41` 逐字 `		InputTokens:  c.Usage.InputTokens + u.InputTokens,`；出口 `internal/agent/loop.go:963` 逐字 `		TokensIn: cost.Usage.InputTokens, TokensOut: cost.Usage.OutputTokens,` | `Cost`；单位＝计费 token（同样累加） | `execute()` 内读得到：`cmd/wisp/run.go:1125` 逐字 `				CostTokensIn: int64(res.Usage.InputTokens), CostTokensOut: int64(res.Usage.OutputTokens),` —— 它今天被写进 task_log，**那一格是成本页的料**。票面 line 37 逐字「不做成本页（`D32` 那两个数已有读法，界面归界面）」⇒ ⛔ 不许把 N5 当占用分子复用 |
| N6 | 事件载荷里的两枚整数 | `internal/agent/sink.go:53` 逐字 `	TokensIn  int`、`:54` 逐字 `	TokensOut int`（`Event` 的字段） | 环路在终态填：`internal/agent/loop.go:911`（`brakeStuck`）与 `:963`（`finish`）两处；单位＝计费 token | **这是唯一一枚天然跨过 `execute()` 的载体**——sink 由装配根递进去（`cmd/wisp/run.go:998` 逐字 `		Sink:     consoleSink{out: rt.stdout, stream: rt.stream, publish: rt.publishPanelSnapshot},`），`cmd/wisp/run.go:1260` 逐字 `func (c consoleSink) Publish(e agent.Event) {` 跑在任务 goroutine 上，往 runtime 自有对象里写是现成形状（同函数 `:1266` 逐字 `			c.stream.Append(e.TaskID, e.Text)` 就是这么把 results 送出来的）。⚠ 但今天**两枚整数没有任何消费者**：尺 `grep -rn "\.TokensIn" --include=*.go cmd \| grep -v _test` ⇒ **0 命中**（`cmd/wisp/run.go:1125` 那枚是 `res.Usage`，不是 `e.TokensIn`），且只在 `EvStuck`／`EvDone` 两枚终态事件上非零 ⇒ 类别 **②（有载体、传的是零值／终态值）** |
| N7 | **项目说明子预算的分子** | `internal/projctx/projctx.go:121` 逐字 `	UsedTokens   int`；写在 `internal/projctx/projctx.go:431` 逐字 `	b.UsedTokens = used`（累加处 `:402` 逐字 `			used += o.Tokenize(renderSegment(f).text)`） | `projctx.Loader`；单位＝token，且算子是**调用方递进来的那一枚**（`internal/projctx/projctx.go:171-173` 注释逐字 `// Tokenize is the caller's token heuristic (agent.ApproxTokens), so the`／`// budget math cannot drift away from the loop's own math.`；装配 `cmd/wisp/run.go:1063` 逐字 `		Tokenize:     agent.ApproxTokens,`）⇒ 与 N1 同单位 | **今天完整读得到，一路到 `internal/panel`**：`cmd/wisp/panel_pump.go:108` 逐字 `func (rt *agentRuntime) instructionBundle() *projctx.Bundle {`（`:117` 逐字 `	return rt.instrLoader.Last()`）→ 装配 `cmd/wisp/run.go:718` 逐字 `		Instructions: rt.instructionBundle,` → 泵 `internal/panel/pump.go:177` 逐字 `	Instructions func() *projctx.Bundle` → `internal/panel/pump.go:308` 逐字 `		snap.Instructions = InstructionsSectionFromBundle(p.src.Instructions())`。**两个数就在泵手里**，但出口把它们的字段名册收窄了：`internal/panel/instructions_200.go:170` 逐字 `	sec := &InstructionsSection{Files: files, Reason: b.Skipped}`（整个 `:158-194` 只读 `b.Skipped`／`b.SkipReason`／`b.Files`）⇒ **类别 ②（载体有、数在、映射丢了）**。⚠ 语义边界：它只统计说明文件，不含对话历史 ⇒ ⛔ 不许直接把它当整段对话的分子画（那会低报） |
| N8 | 请求级预留估算 | `internal/llm/provider.go:111` 逐字 `func (r *Request) EstimateTokens() int {`，末尾 `:138-140` 逐字 `	tokens := n / 4`／`	if r.MaxOutputTokens > 0 {`／`		tokens += r.MaxOutputTokens` | `llm.Request` 自己；单位＝估算 token **再加输出预留** | 读不到。产码唯一调用点 `internal/llm/ratelimit.go:222` 逐字 `	reserved := req.EstimateTokens()`（本地限流自用）；`req` 是 `internal/agent/loop.go:415` 逐字 `		req, err := l.buildRequest(ctx)` 的局部。⚠ **含 `MaxOutputTokens` 这一项 ⇒ 它压根不是"已占多少"**，是"准备占多少"；`internal/agent/budgets.go:125-128` 注释明确它只是与 `ApproxTokens` 同规则的那个同源量 |

**分子判语（只摆料）**：票面语义的占用分子（这段对话现在吃多少）在盘上只有一种真算法＝N1，它今天**既不在任何 runtime 字段上、也不在包外任何导出读口上**；N3／N6／N7 三枚是"已有载体但传的不是那一格"，N4／N5 是"同名不同义"（计费累加），N8 是"预留量"。**任何把 N4／N5／N8 当分子的实现都是在造数**，正是票面 AC#2「不许退化成 0 或空串」的反面变体（退化成"一个真存在但说的是另一件事的数"）。

### 2.2 分母候选（8 枚）

| # | 候选 | `file:line` ＋ 被指字符逐字 | 谁给的 | execute 之外读得到吗 |
|---|---|---|---|---|
| D1 | **配置声明的模型窗口（唯一带具名"未知"的一枚）** | `internal/config/schema.go:353` 逐字 `	// ContextWindow in tokens; 0 = unknown.` ＋ `:354` 逐字 `	ContextWindow int \`toml:"context_window"\`` | 用户在 `config.toml` 的 `[llm.providers.*].models.*.context_window` 里写的 | **读得到，零新连线**：`internal/llm/resolver.go:141` 逐字 `		ContextWindow: spec.ContextWindow,` 把它抄进 `Endpoint`（声明 `internal/llm/resolver.go:56` 逐字 `	ContextWindow int`），装配根存在 runtime 字段 `cmd/wisp/run.go:302` 逐字 `	endpoint llm.Endpoint`（赋于 `cmd/wisp/run.go:441` 逐字 `	rt.endpoint = ep`，发生在 `assembleRuntime` 里、早于任何 `execute()`）⇒ **`rt.endpoint.ContextWindow` 今天就能被快照 reader 直接读，且 `0` 是它自带的具名"未知"** |
| D2 | provider 自报的窗口（D1 的回退源） | `internal/llm/provider.go:50` 逐字 `	// MaxContextWindow of the model (0 = unknown); D15 thresholds scale by it.` ＋ `:51` 逐字 `	MaxContextWindow int` | C5 接缝 `LlmProvider.Info()` | **读得到**：`rt.provs` 是 runtime 字段（`cmd/wisp/run.go:300` 逐字 `	provs    []llm.LlmProvider // built chain, kept for the probe path`）。⚠ 今天 `execute()` 之外的消费者有没有它是另一回事（本腿未量到产码 reader，写"今天存在但无人读"） |
| D3 | **缩放后窗口**（环路实际拿来缩放阈值的值） | `internal/agent/budgets.go:54` 逐字 `	ContextWindow int`；产点 `internal/agent/budgets.go:85-88` 逐字 `func BudgetsFor(ctxWindow int) Budgets {`／`	if ctxWindow <= 0 || ctxWindow > ReferenceContextWindow {`／`		ctxWindow = ReferenceContextWindow` | `BudgetsFor()`，输入是 D1→D2 的回退结果（`internal/agent/loop.go:218`／`:220`） | **读不到**（只有那枚 `execute()` 局部的 `loop` 有：`internal/agent/loop.go:260` 逐字 `func (l *Loop) Budgets() Budgets { return l.b }`；三枚产码调用者全在 `execute()` 内＝`cmd/wisp/run.go:1056`／`:1062`／`:1109`，⛔ 没有一枚取 `.ContextWindow`）。★ **而且就算读到也不能当"分母已知"**：D1 写 0（未知）会被折成 128000、任何 >128000 的声明也折成 128000 ⇒ 它**恒非零、恒不报未知**。⇒ 用 D3 当分子分母的除数是票面 AC#2 禁止的那一发（把未知当已知） |
| D4 | 压缩触发线 | `internal/agent/budgets.go:74` 逐字 `	HistoryCompressTokens int`；值 `:107` 逐字 `	b.HistoryCompressTokens = s(refHistoryTokens)`，源 `:36` 逐字 `	refHistoryTokens  = 12000 // D15(4) compression trigger` | D15(4) 冻结参考数（按窗口缩放） | 同 D3：只在 `Budgets` 里。⚠ **它不是"总共可用"**，是"到这儿就该压"：拿它当分母会把一根还剩一大截的条画成满格 ⇒ 具名排除 |
| D5 | 每任务 token 刹车 | `internal/agent/budgets.go:79` 逐字 `	TokenBudget int`；`:108` 逐字 `	b.TokenBudget = s(refTokenBudget)`，源 `:37` 逐字 `	refTokenBudget    = 200000`；消费 `internal/agent/guard.go:104-105`（`NewGuard` 里 `cfg.TokenBudget <= 0` 才回落到它） | C22 | 同 D3。⛔ 与 N4 同族（刹车 vs 用量＝成本轴），不是窗口轴 |
| D6 | 128k 标定基准 | `internal/agent/budgets.go:18` 逐字 `const ReferenceContextWindow = 128000` | 本包常量 | 全局可读（导出常量）。⚠ **它是校准尺上的刻度、不是任何一台模型的窗口**；拿它当分母＝凭空造一个"总共多少"，正是票面硬约束 1「分母不知道就不许画」要拦的形状 ⇒ 具名排除 |
| D7 | D39 提示词段总预算 | `internal/agent/budgets.go:67` 逐字 `	PromptTotal int`；消费 `cmd/wisp/run.go:1062` 逐字 `		BudgetTokens: loop.Budgets().PromptTotal,` → `internal/projctx/projctx.go:170` 逐字 `	BudgetTokens int`（注释 `:168-169` 逐字 `// BudgetTokens is the caller's existing context-budget number: the loop`／`// hands over its scaled D39 prompt total. No budget constant lives here.`） | 缩放后的 D39 表 | **execute 外读得到——但只读得到它的 N7 副本**（`projctx.Bundle.BudgetTokens`，与 N7 同一条链，泵里现成）。⚠ 语义＝说明文件那一格的天花板（128k 下 2300），是对话窗口的子集 ⇒ 只能当"项目说明占用"这一枚**独立小格**的分母，⛔ 不能当整根上下文条的分母 |
| D8 | 组装器的同名别名 | `internal/agent/prompt.go:181` 逐字 `func (a *Assembler) TotalTokens() int {`，体 `:182` 逐字 `	return a.budgets.PromptTotal` | `Assembler` | 读不到（`l.asm` 私有，`internal/agent/loop.go:187` 逐字 `	asm      *Assembler`）。⚠ 注释 `:179-180` 自述 `// TotalTokens reports the section budget sum actually granted by the scaled`／`// table (used by tests to assert the D39 ceiling).` ⇒ 它是测试用别名，不是出口 |

**分母判语（只摆料）**：今天**唯一在 `execute()` 之外可读、且自带具名"未知"**的分母＝**D1 `rt.endpoint.ContextWindow`**（`0 = unknown` 写在 `internal/config/schema.go:353` 那一行注释里），次选 D2。**环路真正缩放所用的那个数（D3）今天读不到，而且它在结构上不可能报出"未知"**——这是本腿对票面 AC#2 最值钱的一枚读数：**未知这一维必须在 D1／D2 那一层就抓住，一旦穿过 `BudgetsFor` 就永久丢失**。

### 2.3 单位一致性（派单点名"什么单位：token 数？字节？"）

- 分子真源 N1 与分母 D1／D2／D3 **同为"估算 token"轴**（bytes／4），可以直接相除；`internal/agent/budgets.go:125-128` 注释逐字写明这是全包统一的口径：`// ApproxTokens is the token heuristic used for every budget decision in this` / `// package. It deliberately matches llm.Request.EstimateTokens (utf-8 bytes /` / `// 4) so the loop's budget math and the seam's pre-flight estimate cannot` / `// disagree; exact accounting comes from Usage events (C23).`
- N4／N5／N6 是**计费 token（provider 报的真值）轴**，且是跨轮累加 ⇒ 与 D 轴相除在算术上就不成立。⛔ 同一格里不许混。
- 字节数在盘上另有两枚（`internal/projctx/projctx.go:94` 逐字 `	Bytes int \`json:"bytes"\`` 与 `internal/panel/instructions_200.go:54` 逐字 `	Bytes int \`json:"bytes"\``）—— 具名排除，它们不是 token 轴。

