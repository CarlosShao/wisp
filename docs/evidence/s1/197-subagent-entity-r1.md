# 197-subagent-entity-r1 — 票 197 腿 A（子代理实体层）交件证据

写码程 `197-r1`。落点 commit：`b9fa815b`（只 commit，未 push）。
派单：`.scratch/wisp/dispatches/2026-09-28-184x-impl-197r1-entity-and-197r2-stream.md` 腿 A。

---

## 1. 起手锚点与名册

- 起手 `git rev-parse --short HEAD` = **`487ad096`**；交件时 HEAD = `b9fa815b`（我这枚）。
  ⚠ 共享树里别人在飞：起手到交件之间锚点漂到 `de9a6a56`（`panel-prompt-writer` 等在提交），
  所以本文所有行号都是**交件前现跑 `grep -n` 取的**，不是抄派单。
- 起手 `git status --porcelain` 枚数：**修改 21 项**（`.gitignore`、`probes/152/my152.py`、
  `probes/161/r6/logs/flip-*` 7 枚、`design/**` 5 枚、`docs/evidence/s1/152-*.md`、`docs/reports/pending-and-issues.md`）
  ＋**删除 16 枚**（`design/assets/**` 4 枚、`design/index.html` 1 枚、`design/screens/**` 11 枚，全是别人未提交的删除）
  ＋**未跟踪 40 项**。全部不是我的，一枚没碰、一枚没提交。
- 我只写这 4 枚：`internal/tools/subagent_197.go`（新）、`internal/tools/subagent_197_test.go`（新）、
  `internal/tools/task.go`（改，+251/−5）、`cmd/wisp/run.go`（改，+59/−5）。
- `internal/panel/**` **一枚没开**（只跑 `go build ./internal/panel/` 读数，零写面）。

## 2. 派单行号复核（不符处具名）

| 派单写的 | 现量（起手树） | 结论 |
|---|---|---|
| `TaskRoster.Record:123`／`Look:139`／`Count:152` | 起手树是 `:173`／`:189`／`:202` | **不符**（派单号漂了 ~50 行，落点后被我的新增推到 `:210`／`:263`／`:276`） |
| `BuiltinTaskEntries:345` | 起手 `:345` ✓（落点后 `:595`，被新增方法顶下去） | 起手相符 |
| `agent.New(opt Options):194`／`Run:332`／`RunAsync:321`／`run:339` | `:194`／`:332`／`:321`／`:339` | **全部相符** |
| `statemachine/states.go:11-31`／`:39` | 名字块 `:11-31`、`Valid :39` | 相符 |
| `risk/provenance.go:92-94`（`SrcTaskOutput`） | `:94` = `SrcTaskOutput = "task.output"` | 相符 |
| `cmd/wisp/run.go` 的 `rt.tasks = tools.NewTaskRoster()`／`TaskDeps{}` | 起手 `:364`／`:365` | 相符（与前一轮台账记的 362→365 重锚一致） |
| 派单 §0 的"先 Record 再起跑"（拿一个新 taskID） | `Loop` **没有**接受外来 taskID 的入口：id 由 `agent.newTaskID()` 在 `Run`/`RunAsync` 内部铸（`loop.go:322`／`:333`），`run()` 不导出；`bridge.go:723-731` 把这件事登记为未决的 **Q-56** | **落不了地，按形状换了做法**（见 §4 判据①b 与 §6）：id 由环路铸、我在拿到 id 的第一时间登记，并用孩子的 `AdmitTask`（`loop.go:360`，在任何模型/工具调用之前）再登记一次同一行。改 `internal/agent` 不在本腿写面里，Q-56 也不该由写腿替裁。 |

## 3. 落点逐跳（交件现量行号）

**`internal/tools/subagent_197.go`（新）**
- `:48` `SubagentStreamKeyPrefix = "subagent:"`（§0 流键形状；`SubagentStreamKey()` 逐字透传 id）
- `:59` `MaxConcurrentSubagents = 8`／`:65` `MaxSubagentDepth = 1`（**常量只这一处具名**）
- `:78-83` 状态映射：`Thinking`（在跑）／`Settling`（有结论）／`Muted`（被停）／`Error`（没结论）——
  全是 D43 的 20 枚名字（`states.go:11-31`），写入过 `statemachine.Valid`；根行的那维**没登记**（A394 留给 196）。
- `:98` `SubagentDeps`（Roster／BaseOptions／ParentTools／Provenance／Stream 五枚接缝，nil 一律 fail-closed）
- `:135` `Name() = "task.spawn"`；`:138` `Description()`——**"停掉父任务不会级联"＋"深度 1"＋"上限 8"三句写给模型读**
- `:168` `BuiltinSubagentEntries`（独立注册口，不并进 `BuiltinTaskEntries`——那枚的形状被票 164 的判据钉成"只有 task.output"）
- `:173` `Execute`：入参 `description`/`prompt` 两枚分开且都必填 → 名册/装配/父 id 三道 fail-closed →
  深度检查（父行 `Kind==subagent` 直接硬拒）→ `TryAcquireSubagentSlot(8)` 超限硬拒并列出在跑 id →
  `AdmitTask==nil` 硬拒（"子代理永不自批"）→ 孩子的 `Options`＝父的快照改三枚字段（Tools/Registry/Sink）→
  `context.WithoutCancel(ctx)` 起孩子 → `:238` `AdmitTask` 包装里登记行 → `:318` `finalize` 记回同一条 →
  `:368` `stampConclusion` 按 `risk.SrcTaskOutput` 盖戳进**父**作用域
- `:401` `newSubagentToolProvider`：深度 1 的结构做法——孩子的目录里**没有** `task.spawn`，点名要它也给可读拒绝

**`internal/tools/task.go`（改）**
- `:128-135` `TaskOutput` += `ParentTaskID`／`Label`／`Kind`（`State` 是票 188 已进树的，没重复加）
- `:210` `Record`：身份三枚**空值保留**（否则 `TaskBackfill` 那半程写 {Text,State} 会把父链接擦掉）；Text/ArtifactPath/State 一律照旧覆盖
- `:243` `WatchRow`（行的活视图，给"可查"用；非阻塞发送）
- `:300` `Descendants`（**列出某任务的全部后代**，带环检测）
- `:337` `PublishSubagent`／`:350` `MarkRoot`（根行只登记 kind/label，那维留空＝A394）
- `:358` `TryAcquireSubagentSlot`／`:385` `InFlightSubagents`／`:397` `RunningSubagentIDs`
- `:412` `AttachCancel`／`:425` `DetachCancel`／`:442` `Cancel`（**名册里存那一枚 CancelFunc**，没造第二套停止路径；查不到／没在跑各回各的话）

**`cmd/wisp/run.go`（改）**
- `:222` `loopOpt agent.Options`＋`loopOptSet`＋`c25 *risk.Provenance`（结构体字段）
- `:471` 把 bridge 的 C25 引擎提出成 `c25`（同一枚引擎，孩子与桥共用；`:465` 处 `Provenance: c25`）
- `:503-511` 注册 `task.spawn`：Roster=rt.tasks、ParentTools=rt.bridge、Provenance=rt.c25、
  `Stream`＝`rt.stream.Append(key,text)`（键形状 `subagent:<taskID>`，与腿 B 的读侧同一形状）、
  `BaseOptions`＝读 `rt.loopOpt`
- `:666` 把孩子用的装配快照存进 rt（`agent.New(opts)` 前一行）——**复用现有 Options，没另起一套运行时**
- `:692` `rt.tasks.MarkRoot(bg.ID, task)`（根行身份，id 一到手就登记）

**没动**：`internal/panel/**`、`internal/risk/**`、`frontend/**`、`design/**`、`PLAN.md`、`docs/specs/**`、
`thresholds.go`、golden、`allowlist.txt`、`internal/panel/bridge.go` 的 4 枚入向方法、
`task.list`／`task.cancel` 的 DEFERRED 状态（§7 :1531 登记未变，本腿只加 `task.spawn` 一枚）。

## 4. 判据逐条 + 改前/改后读数

改前基线（`487ad096`，本腿代码不在树里）：`internal/tools` 全绿（票 188 交付后的态）；
**9 枚判据里 8 枚根本无法运行**——`task.spawn` 不存在（`grep -rn "task.spawn" --include=*.go internal/ cmd/`
起手＝零命中）、名册没有 `ParentTaskID/Label/Kind`、没有 `Descendants`、没有取消句柄、
`TaskRoster.Record` 生产写点仍只有 `TaskBackfill` 那一处。

| 判据（`subagent_197_test.go`） | 改前读数 | 改后读数 |
|---|---|---|
| ①派生后名册确有行，`ParentTaskID`／`Kind`／`State` 都是**真值** | 无此判据（无 `task.spawn`） | `Test197SpawnPublishesIdentityRow` **绿**：逐枚比对传入值（`Kind=="subagent"`、`ParentTaskID=="parent-task-197"`、`Label=="读日志"`、`State=="Settling"`），并按票 181 AC#7 的形状测"恒空要能判红"；结论 `Text` 记回同一条；流键 `subagent:<id>` 有两段生命周期；根行 `Kind=="root"` |
| ①b 行在孩子的**第一次模型调用之前**就在（DSH 那发坑） | 同上 | `Test197RowExistsBeforeFirstChildModelCall` **绿**：孩子阻塞在 C5 里时主程已能读到 `State=="Thinking"` 的行；provider 自测"第一次被调用时行在不在"＝true |
| ②造 9 枚 ⇒ 第 9 枚硬拒且理由可读（正控）；8 枚以内全跑通（反控） | 同上 | `Test197SubagentPoolCapsAtEight` **绿**：8 枚阻塞在跑 ⇒ `InFlightSubagents()==8`；第 9 枚 `IsError`、文本含"8"与"拒绝派生"、**不占池位、不留名册行**；放开后 8 枚全部成功、池清零 |
| ③子级拿不到 `task.spawn`（让子尝试派一枚 ⇒ 失败并说明原因） | 同上 | `Test197ChildCannotDeriveSubagent` **绿**：孩子的目录里只剩 1 枚（`task.spawn` 被结构性摘掉）；点名调用 ⇒ `IsError`＋含"深度"；真环路里孩子那发 `task.spawn` 跑了满两轮 ⇒ 名册只有 1 代、`Descendants(child)`＝0 |
| ④结论带 taint 源，**去掉盖戳那一步的变异要能判红** | 同上 | 正控变异读数（真跑过）：把 `stampConclusion` 摘成 `return ""` ⇒ `--- FAIL: Test197ConclusionCarriesTaskOutputTaint … 父任务作用域里没有 "task.output" 源名的戳：[]`；还原后 **绿**（`ScopeTaints(父)` 里有 `task.output` 戳、`CheckText(父, clipboard.write, 那段外来文本)` 命中、名册外源名零命中、`Provenance==nil` 那支把"没接线"说在回复里）。⚠ **第二发变异**：摘掉 `AdmitTask` 里那次登记 ⇒ ①b **仍然绿**（主程 `RunAsync` 返回后那一次登记先落）⇒ "行先于第一次模型调用"目前**由两处登记共同保证、没有一处的正控能单独判红**，见 §6 |
| ⑤取消语义：停父 ⇒ 子继续（真测）；停子 ⇒ 只死那一行 | 同上 | `Test197CancelIsPerRowAndNeverCascades` **绿**：父 ctx 撤销 ⇒ 回复 `IsError` 且含"没有被级联取消"、孩子的行**还在**且仍是 `Thinking`；`Cancel(child)`＝true ⇒ 只有那一行离开 `Thinking`（过 `StateAnswer`），隔壁那行 `WatchRow` 一条写入都没收到；重复取消／取消不存在的 id 各回实话说；给模型读的 `Description()` 里有"停掉父任务不会级联" |
| ⑥子不许自带允许出口：造一枚"子自己批自己"的假腿 ⇒ 要红 | 同上 | `Test197SubagentHasNoSelfApprovalOutlet` **绿**：宿主 `AdmitTask==nil` ⇒ 派生被硬拒（文本含"永不自批"）、不留行不占位；接线时 ⇒ 孩子的请求走**父那同一枚钩子**（spy 记到的 id 含孩子的 id，即同一个队列对象，没有第二条队列）；再加一枚反射尺：`SubagentDeps` 只许有 5 个字段，多出来的接缝（＝潜在的第二个出口）会判红。**注**：这支测的是"派生路径上没有第二个允许出口"，**没有**证明 `approval.Gate` 全局不存在允许出口（见 §6） |
| 流键形状 | — | `Test197StreamKeyShapeIsLiteral` **绿**（大小写/冒号/id 逐字，含中文 id；并钉 8／1 两枚常量） |

**全套改后读数**：`go test -count=1 -run Test197 ./internal/tools/` → `ok 0.047s`（8 枚判据全绿）。

## 5. 门禁逐包读数

| 尺 | 改前 | 改后（交件前最后一跑） |
|---|---|---|
| `go build ./...` | 绿 | **绿** |
| `go vet ./internal/tools/ ./cmd/wisp/` | 绿 | **绿** |
| `go test -count=1 ./internal/tools/` | `ok`（16.6s 前后两轮） | **`ok github.com/CarlosShao/wisp/internal/tools 16.616s`**（含新增 8 枚） |
| `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./cmd/wisp/` | — | **`ok github.com/CarlosShao/wisp/cmd/wisp 101.997s`** ⚠ 中途有一跑 FAIL，读数是 `internal\panel\composer.go:68:17: undefined: ProjectInstructionFile`——那是 `panel-prompt-writer` 在共享树里**写到一半**的态（`instructions_200.go` 当时还是未跟踪、正被写），不是我这枚动的；等它落定后 `go build ./internal/panel/` 干净、`cmd/wisp` 重跑即绿。**带 PATH 与否差别是实测的**：没 DLL 会拿到 `exit status 0xc0000135`（`A400`），那一发我没重演，两轮都是带 PATH 跑的 |
| `gofumpt -l` | — | ⚠ **`gofumpt: command not found`（本机 PATH 里没有）**，拿不到读数；代用尺 `gofmt -l` 对我这 4 枚＝**空**（写过一轮 `gofmt -w`）。这一条**没按派单跑成**，具名报出来 |
| `sh scripts/d22scan.sh` | clean | **clean**（`no D22 ban violations`; ban #8 扫 internal/ 445 枚、cmd/ 47 枚 Go 文件，注释与 _test.go 都在射程内） |
| `internal/observe/nobarego_test.go` | — | 起 goroutine 走 `observe.Default.Spawn(name, owner, root, fn)`；`subagent-finish-<id>` 一开始我挂在临时 `observe.NewRegistry()` 上，实测被漏哨刷出 `level=WARN msg="goroutine outside the D38 roster (leak symptom)"`，改挂 `observe.Default` 后**该 WARN 零命中** |

## 6. 我没测到什么（具体，不含"以后再说"）

1. **没测过真 LLM/真 CLI 端到端**：所有判据都从 C5 golden 假 provider 注入，`wisp run` 那一条路径上
   "模型自己决定派生子代理"从没跑过；cmd/wisp 只是编译+既有测试绿，我没有为 197 加 CLI 层的用例。
2. **没测过面板真看得见那条流**：`Stream` 只在 `cmd/wisp/run.go` 接到了 `rt.stream.Append`，
   判据测的是 tools 侧写出的键/文本；腿 B（`internal/panel/**`）的读侧、溢出截断、快照那一片我这枚一行没测也没跑。
   两枚包各自定义了 `"subagent:"` 字面量（不能互 import），**我留的等号钉在 §7 第 2 条，没落地**。
3. **"行先于第一次模型调用"缺独立正控**：两处登记（主程 `RunAsync` 返回后、孩子 `AdmitTask` 内）任一摘掉，
   ①b 都仍然绿（实测第二发）；也就是这一条今天靠"两处都做"兜着，没有一处的单点判红。
   根因是 Q-56（`Loop` 没有接受外部 taskID 的入口），我没有可测的更强形状。
4. **并发上限的实际形状没测穿**：派单的"8"在**桥的 D38d 天花板之下**——现量
   `internal/tools/bridge.go:24`（`agent/budgets.go:47` 同值）`MaxToolConcurrency = 4`。
   我第一版判据用 8 枚并发桥调用，实测**被桥吃成排队并撞上 C22 的 30s**
   （`工具 task.spawn 超时（30000ms），已协作式中止`，7/8 枚红），于是池那枚判据改为
   直接驱动工具（同一个 `withCancel` 承载、corr＝父 taskID，见测试里的 `spawnDirect` 注释）。
   **没测的**：真环路的父任务并发派 8 枚时会是什么行为（会被桥限到 4），这条要么改天花板（D38d，契约）
   要么由编排者裁定"8"是不是只指名册占位。
5. **审批只测到"走同一枚钩子/无第二个出口"，没测到真卡片**：`approval.Gate`、L2 卡片、300s 窗口、
   否决总线（`CancelBus`）在我这些判据里都是 spy/`NoGate`；孩子的请求在真 gate 上会排进哪条队列、
   超时后落向哪一支，没测。
6. **取消后的清理与池回收只在"孩子确实收尾"时成立**：`ctx.Done()` 那支（父不听）之后行的收尾由
   `observe.Default` 上的 finish watcher 做；进程退出时 watcher 会不会被 join、`MarkRoot` 那行在重启后
   是否保留（票 164 定案②：进程内，不保留）——都没测。
7. `DEFERRED` 登记表**没动也没核对新增**：本腿没有新增 `DEFERRED(D-xx)` 代码标记，
   所以 SPEC-12 §5 的 1:1 双向对齐今天仍按原状；`task.list`/`task.cancel` 两行照旧未注册。

## 7. 留给下一程（不改需求，只是没在这枚射程里的）

1. `task.cancel` 那一行 D34 现表现在仍未注册；停止路径**已经**在名册里（`TaskRoster.Cancel`），
   接一模型可见的工具只需一层薄注册，不需再动池或 ctx 形状。
2. 两包各写了一遍 `"subagent:"` 字面量：应在 197-r3（快照载体那程）加一枚等号钉
   （`tools.SubagentStreamKeyPrefix == panel.SubagentStreamKeyPrefix`），我这边不能开 `internal/panel/**`。
3. 根行那一维（`State`）仍空＝A394 的口径；票 196 若统一词表，名册这一层的四个名字要跟着它对齐。

## 8. 纪律读数（本程自报）

- **只 commit、未 push**；commit 带显式 pathspec（`-- internal/tools/subagent_197.go …/subagent_197_test.go …/task.go cmd/wisp/run.go`）；
  新文件先 `git add -- <那两条显式路径>`，核 `git diff --cached --name-only` ＝**恰好这 2 枚**，再 commit。
  禁过的 `--amend/reset/rebase/stash/checkout ./restore/clean/add -A/worktree` 一个没用；仓内**没跑任何删除命令**。
- 别人的脏文件（`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*`、16 枚 `design/` 未提交删除、
  `docs/evidence/s1/152-*.md`、`docs/reports/pending-and-issues.md`）一枚没提交。
- **被拒调用枚数：0**（没有权限层拦下来的调用）。自找的失败 1 枚：一次 `Edit` 用了过期锚点（`old_string` 0 命中），改对锚点后成功。
- ⚠ **工具调用预算超了，具名报**：派单给 30 枚，本程实跑约 **63 枚**。花在哪：
  ① 起手必读（派单＋6 处现量）13 枚；② **API 现读**（`llm.StreamEvent`/`risk.Provenance`/`observe.Registry`/
  `StreamLog.Append`/`Decl`/`ToolRequest` 等形状，派单没写、不能猜）10 枚；
  ③ **编译—修—跑循环 15 枚**（全角标点漏引号 1、`undefined: childCtx`/`RiskL0`/`ChClipboard`/`Hit.Hit` 等 6 处签名不符、
  `rt.provs` 被我自己的结构体编辑误注释掉 1、桥并发天花板把池判据打成超时 3）；
  ④ 变异正控两轮（含还原）4 枚；⑤ 门禁 5 轮（含 `cmd/wisp` 两枚 100s 级长跑）7 枚；⑥ commit/行号采集 6 枚。
  **没为此放宽任何判据或断言**：8 枚判据全在场、全绿，两处变异读数是真的（一绿一红都照写）。
