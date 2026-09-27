# 176-a1b —— 后台任务「起跑口」普查（只读设计核·重派枚·零产码）

- 派单：`.scratch/wisp/dispatches/2026-09-27-235x-readonly-176-a1b-where-does-the-background-start-port-live-retry.md`
- 票面：`.scratch/wisp/issues/176-there-is-no-spawner-runasync-has-zero-production-call-sites-and-taskroster-record-has-zero-writers-so-164-ac4-has-no-object-and-175s-unstamped-path-waits-for-one-writer.md`
- 性质：**只读**。未产码、未改测试、未改判据、未勾任何框。

## 0. 起手锚点（本枚 §0 现跑，逐字贴；不抄派单里的号）

```
$ date "+%Y-%m-%d %H:%M %z" && git rev-parse --short HEAD && git log --oneline -3 && git status --short && git branch --show-current
2026-09-27 23:54 +0800
f6a4eebd
f6a4eebd 派单 176-a1b（重派同一格）：明写"上一枚的通知带着一枚不存在的提交号来过、别继承它的任何叙述"、写面改名防撞、锚点一律现取（我手敲哈希敲错过一次）
1ed0bc14 ledger(A356)：本会话第二枚假"完成通知"落在我自己派的 176-a1 身上——报告通篇像真、提交号 150f798c 查不到、表与台件目录都不存在；我只采纳可复算的两枚地基读数（RunAsync 不存在／Record 零写者），不采纳任何交件叙述、不翻勾，重派一枚换文件名
fbcdf441 派单 176-a1（只读设计核·第一枚"一格一单、上限 45 次"）
dev
```

| 件 | 要求 | 实测 |
|---|---|---|
| `git rev-parse --abbrev-ref HEAD` | 必须 `dev` | `dev` 相符 |
| `git rev-parse HEAD` | 我这边现算应为 `57d6491b` 一带 | **`f6a4eebd38e8ce9d53878e8443ad063645e33b2b`** ⇒ 与派单正文那两个号都不同，**以本枚为准**（`57d6491b` 在 git 里我没查，按 §0.5#3 只当线索） |
| `git status --porcelain -- internal/ cmd/` | 必须空 | **空**（rc=0，零行）⇒ 干净树，可以取证 |
| `go test -count=1 ./internal/risk/ ./internal/tools/` | 本程唯一一次跑测试包 | **`internal/risk` FAIL / `internal/tools` ok** ⇒ 见 §0.1 |
| 票面 176 起手 Progress log | 决定是否只登记自己 | 起手时**只有第 40 行标题、下面零字节** ⇒ 当时无人追加（但见 §0.2，情况在程中变了） |

### 0.1 起手测试读数（**这一枚不是我造成的，我也没去放宽它**）

```
$ go test -count=1 ./internal/risk/ ./internal/tools/
--- FAIL: TestResolvePerCallBudget (2.31s)
    pathresolver_budget_norace_test.go:34: C26 Resolve: 1198746 ns/op = 1.199 ms/op (budget 1.000 ms, 981 samples)
    pathresolver_budget_norace_test.go:37: C26 budget breach: Resolve averages 1.198746ms per call, budget 1ms
FAIL	github.com/CarlosShao/wisp/internal/risk	9.102s
ok  	github.com/CarlosShao/wisp/internal/tools	22.437s
```

- 红因：`internal/risk` 的**每调用耗时预算**（1 ms）在本机被超（实测 1.199 ms，981 样本）。这是时序敏感断言，**不是**本枚引入的（本枚零产码、且 `internal/ cmd/` 起手即净）。
- 处置：按 `AGENTS §1.1`「SLO 阈值 / `thresholds.go` 一字节不许动」，**未动任何阈值、未重跑以求翻绿**。台账该记：**在锚点 `f6a4eebd` 上，`./internal/risk/` 逐包跑是红的**。
- ⚠ 这条与票 177 的「豁免今天在生产里是惰性的」是两回事，别并成一枚读数。

### 0.2 本枚特有的并发事实（**「别人在飞」这一支被实测撞上了**）

- 派单 §0.5#1 说上一枚 `176-a1`「`probes/176/` 目录不存在」。派单 §0.5#2 说：若发现它其实还在跑，**只登记自己那一段、不重抄它的结论**。
- 本枚 23:54 取锚点时，`git status --short` 的未跟踪清单里**没有** `probes/176`（它当时确实不存在，`ls` 报 absent 类）。
- 本枚 00:01 现查（`find .scratch/wisp/probes/176 -maxdepth 2`）：

```
.scratch/wisp/probes/176/a1/census.sh       4869 B  mtime 2026-09-27 23:59:17
.scratch/wisp/probes/176/a1/readings.txt    1876 B  mtime 2026-09-27 23:58:59
.scratch/wisp/probes/176/a1/readings2.txt   3932 B  mtime 2026-09-27 23:59:35
```

  —— 三枚文件的 mtime **全部落在本枚工作时间段内**（本枚 23:56:03 建了自己的 `a1b/logs`），且 `git check-ignore` 说它们**没有被忽略** ⇒ 它们在本枚 23:54 那次 `git status` 之后才出现。**结论：另有一枚 `176-a1` 正在同一格上跑。**
- 本枚**没有读** `a1/` 那三枚文件的任何一个字节（只 `ls`/`find` 取了名字、大小、mtime）。它的任何读数、任何结论，本表一律不继承、不引用、不比对。
- 本枚写面已按 §0.5#2 改名（`a1b`），与它在盘上零重叠。
- 00:01 现查票面 176 的 Progress log：`git status --porcelain -- <票面>` 仍**为空** ⇒ 截至那一刻**它还没往票面上追加过任何东西**，所以本表没有「逐字引它的原文」可引；若它随后追加，以它自己的原文为准，本枚不代述。
- ⚠ **本枚写表期间（约 00:0x）它又前进了一步**：`docs/evidence/s1/176-background-start-port-census-a1.md` **已在盘上出现**（`git status` 报 `??` 未跟踪）。⇒ 「它的表不存在」这句到本枚写表时**已经过期**；「它的提交不存在」仍然成立（`git cat-file -t 150f798c` → `fatal: Not a valid object name`）。**本枚依旧没读它的内容**，也不与它比对 —— 两枚表若要取舍，是编排者按盘上可复算的那把尺裁，不是本枚裁。

## 1. 第一问：两枚数（票 176 AC#1）—— 先测，命令原文旁边全给

| # | 量 | 尺形（可复制） | 实测 |
|---|---|---|---|
| 1 | `RunAsync` **生产调用点** | `grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ \| grep -v _test.go \| wc -l` | **0** |
| 2 | `TaskRoster.Record` **生产写者** | `grep -rn '\.Record(' --include='*.go' internal/ cmd/ \| grep -v _test.go \| wc -l` | **0** |
| 2b | 同一枚、把接收者收窄到在册名 | `grep -rnE '(tasks\|roster\|Roster)\.Record(' --include='*.go' internal/ cmd/ \| grep -v _test.go \| wc -l` | **0** |

⇒ **两枚都是 0，没有任何一枚 ≥1 ⇒ 票 176 的前提活着，不作废**，175 的排程前提不变。本枚据此继续做 §2/§3。

### 1.1 ⚠ 必须更正的一条：`RunAsync` **不是不存在，是存在且零生产调用点**

派单 §1 与本仓库台账（`1ed0bc14` 那句「RunAsync 不存在」）把它写成「**这个符号压根没有**」。**现量不支持这个说法**：

```
$ grep -rn 'RunAsync' --include='*.go' internal/ cmd/
internal/agent/loop.go:320:// RunAsync spawns the task on a roster-named goroutine and returns immediately.
internal/agent/loop.go:321:func (l *Loop) RunAsync(ctx context.Context, input string) *RunningTask {
internal/tools/bridge.go:719:// mints itself (Run/RunAsync -> newTaskID, internal/agent/loop.go:332/:321).
cmd/wisp/run.go:357:	// RunAsync still has zero production call sites (measured, evidence file
internal/agent/control_test.go:81,131 / forensics_test.go:99 / loop_golden_test.go:172,251,275   ← 6 枚调用点，全在测试里
```

- 声明在位：`internal/agent/loop.go:321`。测试侧调用点 **6 枚**（`grep -rn '\.RunAsync(' --include='*_test.go' internal/ cmd/ | wc -l` → `6`）。生产侧 **0 枚**。
- 票面正文（第 11 行「在**生产码里零调用点**」）用的口径是**对的**；要更正的是**台账/派单**那句「不存在」。这是两种不同事实：「不存在」⇒ 落地腿要先造函数；「存在、零调用点」⇒ 落地腿**只需要有人调它**。本枚按后者给建议。
- 对照：`internal/tools/task.go` 的 `Record` 才是「有方法、有语义、零写者」那一枚；`cmd/wisp/run.go:216` 持有字段、`:361` 造表、`:362` 注册读者 —— **全仓 `tools.TaskRoster` 的生产引用只有这三行**（`grep -rn 'tools.TaskRoster\|TaskDeps{' --include='*.go' internal/ cmd/ \| grep -v _test.go` → 只有 `cmd/wisp/run.go:216` 与 `:362`）。名册与 `internal/` 之间**零耦合**。

## 2. 第二问：起跑口的合法落点（四条形状禁区之下的落点裁定）

### 2.1 先例：**有，而且是同一枚形状**

「宿主自己起 goroutine 跑一件活并把结果记账」在本仓**不是没有先例**，而是有 8 处，且**没有一处**用裸 `go func(`：

```
$ grep -rn 'go func(' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l
0
$ grep -rn '\.Spawn(\|SpawnCapture(' --include='*.go' internal/ cmd/ | grep -v _test.go
```

| 先例 | 位置 | 形状 | 与本票的距离 |
|---|---|---|---|
| **`Loop.RunAsync` 本身** | `internal/agent/loop.go:324` `l.reg.Spawn("agent-task-"+id, "agent", t.root, func(c context.Context){...})` | 宿主起工人 + 名字进 D38 名册 + owner=`agent` + 结果落在 `RunningTask.result` | **就是这一枚**。它已经把「起」和「收」都做完了，缺的**只有下一跳：没人把 `t.result` 写进 `tools.TaskRoster`** |
| tool-exec 工人 | `internal/agent/loop.go:650` `l.reg.Spawn("tool-exec-"+taskID, "agent", root, ...)` | 按 taskID 命名的一轮工人 | 命名先例（`<用途>-<taskID>`）|
| 内存写者 | `internal/memory/writer.go:79` `q.reg.Spawn("db-writer", "memory", nil, ...)` | 队列→落库的常驻回填 | 「异步产出→回填」的**同族**形状 |
| 音频采集 | `internal/audio/audio.go:180-181` → `wasapimic_windows.go:83` / `wavinjector.go:84` | `SpawnCapture` 收 registry | D47 陪聊路径那条 |
| UI STA / 日志刷盘 / 插件回收 | `internal/ball/ball_windows.go:170`、`internal/observe/logging.go:96`、`internal/plugin/disposal.go:233,296` | 同上 | 一般性 |
| `cmd/balldebug` | `cmd/balldebug/main.go:251,432` | 调试二进制里的 owner=`balldebug` | 组合根侧起工人的先例 |

owner＋recover 的强制点只有一处，**新落点走 `observe.Registry.Spawn` 就自动合规**（ban 1 与 AC#4(iii) 同源）：`internal/observe/goroutine.go:262` 是唯一的 `Spawn` 实现，`:281` 才 `go r.run(...)`，`:296` 起 `recover()` 并把 panic 记成事件（`:299` 带 `Owner`）。⇒ **落地腿不需要自己处理 recover，需要的是别绕过 `Registry`。**

### 2.2 D15(3) `Spiller` 能不能被起跑口复用 —— **能，但不要裸用 taskID 当 call-id**

`Spiller` 是**导出的**、与 `Loop` 无关的一层，宿主可以自持一枚，不需要碰 `Loop` 的任何方法：

- `internal/agent/spill.go:43` `func NewSpiller(dir string, b Budgets) *Spiller`（导出）；`:84` `Prepare(callID, text string) (Spill, error)`；`s.dir` 由调用方给。
- 组合根手里已经有同一枚目录：`cmd/wisp/run.go:609` `ArtifactsDir: filepath.Join(rt.spec.dataDir, "artifacts")`。
- `TaskOutput` 的两栏（`internal/tools/task.go:93-102`：`Text` + `ArtifactPath`）**正好是 `Spill` 的两栏**（`out.Text` 已是 PLAN:431 那形头尾桩、`out.Path` 是副本文件）。⇒ 「复用 `Spiller`」不是勉强搭桥，是同一枚契约的两端。
- `TaskRoster.Record` 的语义与它一致：`task.go:120-121` 明写 last-writer-wins，并点名「D15(3) 对 spill artifacts 已记过同一语义」。

**撞不撞名？本枚现量的结论是：裸用会撞，加前缀不撞。**

- `artifactName`（`spill.go:184-199`）渲染 `tool-output-<encode(id)>.txt`；编码是 `[^a-z0-9_-]` 全大写字节百分号转义（`spill.go:167-175` 的理由是 NTFS 大小写折叠），**注入性：`name(a)==name(b) ⇒ a==b`**（超长 id 那支退化为 sha256 抗撞，`:181-183` 自己承认）。
- `nextSequence()`（`spill.go:150-155`）**只在 id 编码为空串时**才进入名字（`:186-194`）。taskID 永不为空 ⇒ **序列计数器不构成共享冲突**（两枚 `Spiller` 实例各自计数也无妨）。
- 但**命名空间是共享的**：`newTaskID()`（`loop.go:1107-1119`）造的是小写十六进制 UUID 形（`xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx`），全落在字面集 `[a-z0-9_-]` 里 ⇒ **taskID 编码后就是它自己**，`tool-output-<taskid>.txt` 与 `tool-output-<tool-call-id>.txt` 躺在**同一目录、同一前缀**。工具 call id 是模型侧供给的（`spill.go:81` 自己写着 "it is model-supplied"），今天**没有任何一道闸门断言过它不会长成 UUID 形**。真撞上的后果不是报错，是 `spill.go:117-131` 那条 `ErrExist` ⇒ **temp 写＋rename 的 last-writer-wins**：任务副本被一次同 id 工具调用的字节悄悄换掉，正是票 79 当年修的 C25 来源断裂。
- ⇒ **本枚的裁定（形状，不是码）**：起跑口**不要**把 taskID 裸送进 `Prepare`。给它加逻辑前缀 —— 这一手在本仓**已有权威先例且是同一枚 `RunAsync` 干的**：`loop.go:324` 用的是 `"agent-task-"+id`，`loop.go:650` 用 `"tool-exec-"+taskID`。照抄这个命名法（`Prepare("agent-task-"+taskID, ...)`）就把共享命名空间问题关掉了，且**不新增任何导出面**（`artifactName` 是包私有的，别去导出它）。
- ⚠ 另一条**不落在本票但要点名**：`internal/tools/task.go:331` 的 `pointerNotice` 用 `os.Stat(raw)` 判存在性，而 raw 是从名册里取的字符串。它不规范化路径（`task.go:309` 明写禁手写 `filepath.Clean/Abs`），所以 AC#4(iv) 不撞；但**起跑口是第一次往那栏写字符串的人** —— 那栏一旦被填，`os.Stat` 就开始对宿主自己写的字节生效。落地腿要清楚这一点，别把它当成「已经在验证来源」。

### 2.3 `cmd/wisp/run.go`：填 `Record` 的那一步接在哪儿最不外溢

现状（现量行号）：

```
cmd/wisp/run.go:216   tasks    *tools.TaskRoster        ← 字段
cmd/wisp/run.go:361   rt.tasks = tools.NewTaskRoster()   ← 造表（在工具注册那一族里）
cmd/wisp/run.go:362   tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks})  ← 唯一的读者注册
cmd/wisp/run.go:357   // ... RunAsync still has zero production call sites ...    ← 现场留话：谁填不归票 164
cmd/wisp/run.go:575   func (rt *agentRuntime) admitTask(taskID string) func()     ← 唯一被递到 taskID 的边界
cmd/wisp/run.go:591   loop, err := agent.New(agent.Options{ ... })                ← 环路每次 execute 现装
cmd/wisp/run.go:604   AdmitTask: rt.admitTask
cmd/wisp/run.go:605   Registry:  observe.NewRegistry()
cmd/wisp/run.go:622   res := loop.Run(ctx, task)                                   ← 今天唯一的生产起跑（同步）
cmd/wisp/run.go:769   func (c consoleSink) Publish(e agent.Event)                  ← 出事件的落点
```

三个候选，**按外溢从小到大排**：

| 候选 | 位置 | 判定 |
|---|---|---|
| **甲（推荐）** | `execute()` 里 `loop.Run` 那一发换成 `RunAsync` 拿 `*RunningTask`，在其收尾处 `rt.tasks.Record(res.TaskID, tools.TaskOutput{Text: res.Text, ArtifactPath: ...})` | **改动面＝`cmd/wisp/run.go` 内一处、约三行**。`RunningTask` 已把需要的东西全公开了：`ID`（`loop.go:287`）、`Wait() Result`（`:308`）、`Cancel()`（`:300`）、`Pending()`（`:318`）；`Result` 有 `TaskID`/`Text`/`Status`（`:83-86`）。**不新增 `Loop` 的任何方法 ⇒ 禁区 (i) 结构上碰不到。** |
| 乙 | 挂在 `consoleSink.Publish` 的 `agent.EvDone` 分支（`run.go:798`） | 看似更「搭车」，**但它是**有损的：`loop.go:922` 是 `summary := res.Text`，`:923-924` **当 `res.Text==""` 时把 `summary` 换成状态文案 `msg`**，`:939` 发的就是这个 `summary`。⇒ 被取消/被刹车的任务会以「状态文案」冒充「任务输出」进名册。**本枚判它不合格**，除非落地腿愿意在 Event 上再加一栏（那是 C 级契约面，不归本票）。 |
| 丙 | 扩 `admitTask`（`:575`）在返回的撤销闭包里顺手 `Record` | **不可行**：闭包拿不到 `Result`（它在 `run()` 内被 `defer`，`loop.go:366`），只能看见 taskID。它能拥有的恰是 `CloseTask` 那种「不含输出的收尾」，见 `:581` 与 `:564-570` 那段「WHY THIS HOOK OWNS THE CLOSE」的现量理由。 |

**会不会动 `cmd/**`：会 —— 甲、乙、丙三候选的落点全在 `cmd/wisp/run.go` 里。**
⇒ **本枚未动 `cmd/**`（也未动 `internal/**`），按派单 §2.3 的口径写「未动、留给落地腿」**。派单明写：动 `cmd/**` 需要编排者**在派单里具名放开**，而本枚（只读腿）**没有**这道放开。
⇒ **排程后果（要进 `next=`）**：176 的落地腿**必须**带一句对 `cmd/wisp/run.go`（至少）的具名解冻；否则写手只有两条坏路 —— 要么越权动 `cmd/**`，要么退回去给 `Loop` 加方法（撞禁区 (i)）。**这一条是本枚最该被读到的一条。**

### 2.4 面板侧有没有「异步产出→回填」的管子可以搭车 —— **没有，别再找**

- 唯一像样的泵是 `internal/panel/pump.go`：`NewSnapshotPump`（`:145`）、`Snapshot()`（`:150`）、`Marshal()`（`:204`）、`Publish()`（`:215`）、`Publishes()`（`:240`）。它的出口是 `PumpSources.Out func(snapshot Snapshot, data []byte) error`（`:126-130`）—— **同步函数指针，不是队列，也不是 goroutine**。整个 `internal/panel/` 里 `grep -n 'go func\|Ticker'` 命中的只有注释（`pump.go` 的 `:41` 那行是 `EventTick` 这个名字里的 Tick）。
- 泵自己的文件头就把它不是什么都写了：`:15-23`「WHAT THIS FILE IS NOT: the transport. There is no Go -> page channel in this tree today … 这个泵停在字节边界」；`:6-13` 说 `Snapshot` 只有四栏、`NewSnapshot` 曾经零非测试调用者。⇒ **搭车点不存在**；`publishPanelSnapshot`（`run.go:594` 递进 `consoleSink` 那一枚）是**同步**推一次快照，与「后台任务的输出晚点回来」这个形状无关。
- 因此：**名册的写入方不能设计成「由面板泵驱动」**。真要一条异步回填管，本仓今天唯一在跑的同类是 `internal/memory/writer.go:79` 那枚 `db-writer`（队列→SQLite），它是 §2.1 的**先例**而不是**载具**（它的队列语义属于 memory，把 task 输出塞进去＝动数据模型，那是 D35/票外）。
- ⚠ 本枚**未读** `internal/panel/tokens_fourway_test.go`（派单 §2.4 指名别碰，别家常红地界），也未读 `frontend/**`、`design/**`（禁读区）。

## 3. 第三问：取消语义挂在哪，红／绿两向长什么样（票 176 AC#3）

### 3.1 在场的机器（现量，逐枚指回行号）

| 面 | 位置 | 它今天保证什么 |
|---|---|---|
| 取消请求入口 | `internal/agent/loop.go:300` `func (t *RunningTask) Cancel() { t.root.Cancel() }` | 工人自己的 root ctx —— **前提是你手里有那枚 `*RunningTask`**（正是起跑口会给出的东西） |
| ctx 传播 | `loop.go:323` `observe.NewRootFrom(ctx, id)`；`loop.go:345-352` 在 root 自己的 ctx 下跑，并把 `ctx` 换成 `root.Ctx`；`loop.go:377` 每轮开头 `if ctx.Err() != nil` → `finish(... StatusCancelled, "任务已取消")` | 环路**在轮边界**看得见取消 |
| 收尾仍要落盘 | `loop.go:930-931` `terminalWriteCtx(ctx)` + `:932` 的注释（终态行必须活过取消它的那一发；2s 是 ctx deadline 不是墙钟差 ⇒ D42#9） | **这一枚是本判据的正面教材**：它已经解决过「取消之后还要写一次」的问题，解决法是**换一枚 detach 过的 ctx 并显式注明哪一写是被赦免的** |
| D31 否决总线 | `internal/tools/cancel.go:21-31` `CancelBus{Vetoed, Started, Report, Complete}`；`cancel.go:63-69 Vetoed(ctx)`；`cancel.go:74-82 Stopped(ctx)`（「ctx 没了 **或** 用户投了否决票」两形合一） | 运行中的**工具**在步边界自查的接缝 |
| 总线接线 | `internal/tools/bridge.go:83` 字段、`:459 defer b.cancel.Complete(orDefault(req.CorrelationID, req.TaskID))`、`:523` 否决后改判 outcome、`:631/634` 出报告 | 桥只在**调用**粒度上用它，id 是 correlation id（= taskID，`loop.go:356` `j := newTaskJournal(..., taskID, taskID)`） |
| AppliedSteps | `internal/agent/approval/report.go:14/33-34/98`；`internal/agent/tools.go:74-79`；`internal/tools/fs_write.go` 每一发返回都带 `AppliedSteps: s.snapshot()`；`report.go:117 UnreportedSteps = Vetoed && len(AppliedSteps)==0` | 「部分完成后必须报已落地的步」——**这套语义是「已经写了什么」，不是「别再写了」**。别指望它替你把关后续写入 |

### 3.2 判据该挂的接缝（裁一定二）

**挂点＝`tools.TaskRoster` 自己那一层，配合 `Stopped(ctx)`；不是挂桥、不是挂 Loop。**

理由（全部指回现量）：
- `Record` 是**唯一**能把输出写进名册的门（`task.go:123`），且它今天对 ctx 一无所知（签名里没有 ctx）。判据要成立，`Record` 那条路上必须有一处「问一句这个任务还活着吗」。
- 桥那一侧不行：`CancelBus` 的键是 **correlation id**（`cancel.go:40 corr`），语义是**单次工具调用**的否决，不是一枚后台任务的取消；把 taskID 塞进去当 veto 键＝在 D31 上私造第二枚生命周期，那是契约轴（AC#5）。
- `Loop` 那一侧也不行：任何「`func (l *Loop) X(ctx, taskID)`」形状的守卫都直接撞禁区 (i)。

⇒ **落地形状（写手照这段做）**：`Record` 改成接受一次「还在不在」的判据（或名册持一枚按 taskID 登记的活体集合），判定语用现成的 `tools.Stopped(ctx)` / `RunningTask.Root()` 的 ctx —— 二者都已导出，都不需要新面。

### 3.3 红／绿两向（成对，判据本体）

| 向 | 造法（注入面只用 C5/C8/C17/CLI，见 `AGENTS §1.3`） | 期望 | 在未修码上响不响？ |
|---|---|---|---|
| **绿** | 起一枚**会一直产出行**的后台任务（golden SSE 里放一次 `fs.write` 式的多步工具），等它写下第一步后 `Cancel()`，再 `Wait()`（`Wait` 已保证 join 到 root 排空，`loop.go:308-314` / `Pending()` `:318`），然后断言名册里那一 taskID 的记录**在取消之后再没有变化**（或压根没有新记录） | 取消之后**零次** `Record` 生效 | **不响** —— 今天没有起跑口，测试**无法**在纯生产码路径上造出「一枚被取消的后台任务」，此判据**不可能通过**，也不可能红 |
| **红** | 同一台件上故意留一处「取消后仍在写」：把 `Record` 调用放到 `Cancel()` 之后、且不问任何判据的那一发（等价于今天 `Wait()` 回来无脑写名册） | 当场响 | **同样不响** —— 它需要的是同一枚起跑口 |

⇒ **明写（派单 §3 要求的那一句）**：**这一发今天「不可能响」，不是「已经绿」。** 它属**「待起跑口同批诞生」**，与票 176 AC#2 的「同批」是同一条硬约束的两面：**写者落了、取消判据没有 ⇒ 那次交付判不通过**。本枚**没有**跑过任何红/绿台件（本枚只读，且写面不含测试文件），因此这里给的是**判据形状**，不是**读数**；不许被引用成「已验收」。

### 3.4 归口一句

**归票 176。** 理由：票 176 AC#3 逐字要求「起跑口一落地……必须**同批**做成红／绿两向判据——不许再写成『等 164 那格以后再说』」，而票 164 AC#4 的登记形状正是「等的对象今天不存在」；判据的**被测对象**（往名册里写的那一发）由 176 造出来，164 只是第一个**指出**它缺位的人 —— 让 164 再收一次，就是把「三张票都以为自己不归口」那个最坏形状（票面 §为什么值得做）留在原地。

## 4. 门禁三枚（未跑第 4 枚）

```
$ sh scripts/d22scan.sh                                                    rc=0
  d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=207, bans #1-5 cmd/=23,
  ban #6 frontend/=85, ban #7 internal/tools/=20, ban #8 design/=39, ban #8 frontend/=85,
  ban #8 internal/=427, ban #8 cmd/=45
$ sh .scratch/wisp/probes/154/gate-clauses.sh                              rc=0
  # 腿数＝14 声明与实测不符＝0
  # 腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0
  # 腿数断言＝相符（名册上每一腿都记了账）
$ bash tools/d22scan/runtests.sh -C tools/d22scan ./...                    rc=0
  runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0
```

（**未跑** `.scratch/wisp/probes/161/r6/flip-declaration.sh` —— 跑一次就脏跟踪日志。全仓 `./...` **未跑**。）

### 4.1 ⚠ 给落地腿的一条实测警告：**G2 那枚数会把你的注释打进红**

同一枚 `gate-clauses.sh` 里现量（锚 `f6a4eebd38e8ce9d53878e8443ad063645e33b2b`）：

```
## G2 OpenTask/CloseTask 的生产出现点（排 bridge.go 自身；今日应只剩 cmd/wisp/run.go 那一发＋它上面的注释）
$ git grep -nE -w 'OpenTask|CloseTask' f6a4eebd... -- internal/**/*.go cmd/**/*.go :!*_test.go :!internal/tools/bridge.go
f6a4eebd...:cmd/wisp/run.go:564:// WHY THIS HOOK OWNS THE CLOSE. CloseTask's own comment says the composition
f6a4eebd...:cmd/wisp/run.go:581:		rt.bridge.CloseTask(taskID)
```

`want G2 ring` / `want_n 2`（`:358-362`）。⇒ **基线枚数＝2，且其中一行只是注释。** 起跑口的落地腿几乎一定会在 `cmd/wisp/run.go` 里新写一段「取消后要关谁的范围 / 不调 CloseTask 会怎样」的说明，**那会命中 `OpenTask|CloseTask` 字面、把 2 顶成 ≥3**，于是这枚腿按「声明与实测不符」响。
⇒ **两条出路，都必须由编排者选，不许写手临场定**：(a) 落地腿避开 `OpenTask`/`CloseTask` 字面（改说「那一枚 C25 任务范围的关法」）；(b) 认定这是**契约轴上的基线变更**（`want_n` 是登记值），走人工批准重登。**本枚倾向 (a)**：它不动任何登记数字。

## 5. 两条形状禁区的自证（本枚结论有没有撞上去）

| 禁区 | 本枚给的落点会不会撞 | 自证 |
|---|---|---|
| **(i)** `agent.Loop` 上加「收 `taskID` 的导出方法」 | **不撞** | 推荐形状（§2.3 甲）用的是**已经存在**的 `RunAsync(ctx, input string) *RunningTask`（`loop.go:321`，签名里没有 `taskID`）＋ `RunningTask` 已有的 4 枚公开面 ＋ `tools.TaskRoster.Record`。G3 的尺形现量在 `probes/154/gate-clauses.sh:363-366`：`'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID'`，射程 `internal/agent`，`want G3 quiet` / `want_n 0`，**本枚实测命中行数＝0**。 |
| **(ii)** 宿主内部 artifacts 写入做成受门控的 Tool | **不撞** | 起跑口是**组合根里的一次直调**（`execute()` 内），不注册进 `reg`、不进 D34 权威表、没有 `Decl`、不碰 C1/C19/C3。§2.3 乙/丙 之所以差，是因为它们**同样**不解决「谁调」，而不是因为它们在禁区上。⚠ 顺带钉住一条：本枚**没有**建议新增任何 `task.*` 工具名（`task.list`/`task.cancel` 仍按 `internal/tools/task.go:19-25` 留在 DEFERRED 未注册，`PLAN.md:1531`）。 |
| **(iii)** 裸 `go func(` | **不撞** | §2.1 全表都走 `observe.Registry.Spawn`；本仓生产码裸 `go func(` 今天**实测 0 枚**（命令在 §2.1）。 |
| **(iv)** `risk.PathResolver` 之外做文件系统决策 | **不撞（并提醒）** | §2.2 复用 `Spiller` 是把它当**字节落盘器**用（内部 `filepath.Join(s.dir, name)`，`spill.go:106`，属 `internal/agent` 既有实现，不是新决策点）；名册侧的路径判断只走 `PathCanonicalizer`（`task.go:323-330`）。**新增的 `os.Stat` 一类「存在性」念头不要往 `cmd/` 里带。** |

⚠ **不许把「G3 quiet」当成合规证明**：那枚尺是**文本**尺，只匹配同一行里出现字面 `taskID` 的导出方法。把参数改名成 `id` / `tid` / `task` 就能躲过它而**照样违反 AC#4(i) 的实质**。G3 安静＝尺子没响，≠ 形状合法。这条是本枚给对抗验收用的，不是给写手用的台阶。

## 6. 本枚**没**读 / **没**测（照实列）

- **未读** `internal/panel/tokens_fourway_test.go`（派单指名别碰）；**未读、未引** `frontend/**` 与 `design/**`（`AGENTS §1.2` 禁区，本枚全程没打开过）。
- **未读** `.scratch/wisp/probes/176/a1/**` 三枚文件的**内容**（只取了名字/大小/mtime，见 §0.2）—— 上一枚的读数一律不继承。
- **未跑** 任何全仓 `./...`；**未跑** `flip-declaration.sh`；**除 §0 那一发以外没再跑过任何测试包**（`internal/agent`、`internal/observe`、`internal/memory`、`internal/panel`、`cmd/wisp` 五族**没跑过**，所以对它们的健康度本枚**不下任何结论**）。
- **未测** §3.3 那发红/绿 —— 未修码上它**不可能响**（§3.3 已明写），本枚也没有去造它。
- **未读** `docs/PLAN.md` 的任何正文（含 `:1531`、`:2564`、`:431-432`）与 `docs/specs/**` 任何一份 —— 本表里这些号都是**从代码注释里转引**的，属**未验证断言**，落地腿要用时**自己去 PLAN.md 现量**。
- **未读** 上游三份材料（`177-shape-a-impl-r1.md`、`177-c25-r4-path-exemption-c1.md`、`probes/177/r1/logs/**`）—— 派单只让它们当线索；本表所有行号都出自本枚自己的 grep。
- **没查** 派单里那句 `57d6491b` 在 git 里到底是什么（按 §0.5#3 只当线索处理）。
- **未量** `task.output` 在真机上的端到端行为（那要 C17/CLI 注入台件，超出只读腿）。
- ⚠ **一条必须被说成「不是端到端可用」的**：`d22scan` 与 14 枚腿**全绿**，说的只是「没有形状违规、声明与实测相符」。**票 177 那枚豁免在生产里今天仍然是惰性的**（`task.output` 还没进 C25 名册，且没有任何人往名册里写过东西），**票 176 的起跑口今天仍然不存在**。⇒ 门禁绿 **不等于** 「后台任务读输出这件事能用」；能用那天要等 §2.3 甲那一发落地、且 §3.3 那对判据同批在场。

## 7. 交件自证（编排者要求的十四项，逐条）

1. step-0 五件：见 §0 那张表（含红的那一枚已如实贴）。
2. 没读／没测：见 §6。
3. 第一问两枚数：§1，命令原文齐全；**「查无此项」的单列更正见 §1.1（RunAsync 存在、零生产调用点）**。
4. 落点四问：§2.1／§2.2／§2.3／§2.4，每条指回行号。
5. 取消判据红/绿形状：§3.2／§3.3；今天响不响＝**不可能响**。
6. 归口：§3.4 —— 归 176。
7. 门禁三枚：§4，三枚 rc 全 0；另附 §4.1 一枚 G2 陷阱。
8. 两条禁区自证：§5，无标红项。
9. 被拒/没成功的调用：**1 次** —— `Read` 一枚工单时把文件名抄短了（`176-...-zero-writers.md`，真名以 `...-zero-writers-so-164-ac4-...` 才对），报「File does not exist」，取数**之前**发生的，改全长路径后成功。无其它失败。
10. 删除命令：**一枚也没跑过**（零 `rm`/`unlink`/`git clean`/`git restore`）。
11. 工具调用顶：**未超**（预算 ≤45）—— 含本表两处事后更正、提交与提交后校验在内**不超过 38 次**；确切数在最终回禀里给（写文件时没法自证之后的两次）。
12. 伪授权两栏：本程**没有**任何一条来自用户/上游的授权叙述被本枚当作依据；唯一的「授权」是本枚写面清单（派单 §5 + §0.5#2 改名）。**未采纳**派单 §0.5#1 转述的「台账说 RunAsync 不存在」为证据，而是现量后**推翻**了它（§1.1）。出处：本表 §1.1、§6 末条。
13. 凭据值：**零抄录**（全程未打印任何 env 值、未打开配置里的密钥字段；本表只出现变量名与目录名）。
14. `next=`：见 §8。

## 8. next= —— 175＋176 落地腿派出去**还缺**哪几件（按轻重排）

1. **缺一道对 `cmd/wisp/run.go` 的具名解冻**（最重）。§2.3 三个候选全在 `cmd/**` 里，而本枚（与写手的默认口径）都没放开它。**不给这句，落地腿只能二选一：越权，或者给 `Loop` 加方法去撞禁区 (i)。**
2. **缺一次「G2 基线要不要重登」的裁定**（§4.1）。写手在 `cmd/wisp/run.go` 里新写一句含 `OpenTask|CloseTask` 字面的注释就会把 `want_n 2` 顶红。本枚倾向「避开字面、不动登记值」，但那是编排者的选择，不是本枚能替他定的。
3. **缺 `internal/agent` 那批包的健康度读数**（本枚按 §0 的唯一一发跑包限制没跑它们）。起跑口一旦改成 `RunAsync`，被直接改判的就是 `internal/agent`（`control_test.go` / `forensics_test.go` / `loop_golden_test.go` 是今天唯一在调 `RunAsync` 的三族）—— 派写手前要先知道它们在锚点上绿不绿。
4. **缺 AC#2「同批」的可执行定义**：票 175 那枚「外来的读回来要盖戳」判据的**常驻正控用例**落在哪个包、叫什么名，目前只有票 175 侧的叙述，没有本仓现量的在册文件。写手同批交付时要知道往哪加。
5. **缺 §3.2 那枚判据的落点确认**：`Record` 要不要接 ctx / 名册要不要持活体集合 —— 这是形状决定，会动 `internal/tools/task.go` 的导出签名（`Record` 今天不收 ctx）。若算契约轴（AC#5）就需要单独批准；若不算，也要有人明说一句「不收」。
6. **（轻）`task.output` 进 C25 名册** 那件事仍没归口 —— 不在本票面 AC 里，但它决定 177 的豁免什么时候从惰性变成活的。
