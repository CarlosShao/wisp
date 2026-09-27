# 176-a1 —— "后台任务的起跑口"今天缺哪几件、合法落点在哪（只读设计核，零产码）

- 程：`176-a1`｜派单：`.scratch/wisp/dispatches/2026-09-27-235x-readonly-176-a1-where-does-the-background-start-port-live.md`
- 起手锚点：`fbcdf441764a77ccb917348b057fc5fb5907bebf`（09-27 23:51 现量，branch=`dev`）
- ⚠ **本程锚点中途漂移**：读数件跑成时刻到 `f6a4eebd38e8ce9d53878e8443ad063645e33b2b`（09-27 23:59）。漂移的两枚提交都是**票面与台账**（`git diff --name-only fbcdf441..HEAD -- internal/ cmd/` 输出为空 ⇒ **本表所有 `internal/`／`cmd/` 读数在两枚锚点下同值**，行号亦未移）。
- 台件：`.scratch/wisp/probes/176/a1/census.sh`｜有效读数 `readings2.txt`｜`readings.txt`＝本程第一次跑坏的记录（仓根级数算错），零删除规矩下**原样留着**
- 只读声明：本程 `internal/**`、`cmd/**`、`tools/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**` 零字节改动；勾框 0 枚

---

## 0. 起手五件（逐件，含不漂亮的那一件）

| 件 | 读数 | 命令 |
|---|---|---|
| 时刻 | `Sun Sep 27 23:51:22 CST 2026` | `date` |
| 分支 | `dev` 相符 | `git rev-parse --abbrev-ref HEAD` |
| HEAD | `fbcdf441`（派单正文写 `57d6491b` 一带、委托语写 `fbc441`，**两者都与盘上不符**，以本表为准） | `git rev-parse HEAD` |
| 工作树 `internal/ cmd/` | **空＝干净** | `git status --porcelain -- internal/ cmd/` |
| 包测试 | `internal/tools` **ok** (24.743s)；`internal/risk` **FAIL** | `go test -count=1 ./internal/risk/ ./internal/tools/` |

⚠ **起手那一发红不是我造的、也不归我修**（本程唯一一次跑测试包，除 §4 三枚门禁外没跑过别的）：

```
pathresolver_budget_norace_test.go:34: C26 Resolve: 1098445 ns/op = 1.098 ms/op (budget 1.000 ms, 1020 samples)
pathresolver_budget_norace_test.go:37: C26 budget breach: Resolve averages 1.098445ms per call, budget 1ms
FAIL	github.com/CarlosShao/wisp/internal/risk	8.130s
```

形状＝**机器负载敏感的每调用预算**（budget 1.000 ms，实测 1.098 ms，越线 9.8%）。`thresholds.go`／预算常量在 AC#5 的零字节区，**我一字节没碰、也不许任何人为了变绿碰**。它同时是 §2-Q3 那条"起跑口动 `internal/tools`/`cmd` 会不该惊动 risk 预算"的反面警告：起跑口若真起协程跑活，C26 的Resolve 负载形状没有现量证据 ⇒ 属落地腿需要重跑 `./internal/risk/` 的一点，本程没测。

---

## 1. 第一问：把"没有起跑口"量成两枚可复制的数（＝票 176 AC#1）

**两枚都是 0 ⇒ 票 176 不作废、票 175 排程前提不变。** 停手条件（任一枚 ≥1）**没有触发**。

| # | 判据 | 读数 | 命令（可逐字复制） |
|---|---|---|---|
| R1 | `RunAsync` **存不存在** | **存在**，1 枚定义：`internal/agent/loop.go:321` | `grep -rn 'func (l \*Loop) RunAsync' --include='*.go' internal/` |
| R2 | `RunAsync` **生产调用点** | **0 枚** | `grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ \| grep -v _test.go` |
| R2b | 同尺含测试（对照） | 6 枚，**全部在 `_test.go`**：`internal/agent/control_test.go:81`、`:131`、`internal/agent/forensics_test.go:99`、`internal/agent/loop_golden_test.go:172`、`:251`、`:275` | `grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ \| wc -l` |
| R3 | `TaskRoster.Record` **生产写者** | **0 枚** | `grep -rn '\.Record(' --include='*.go' internal/ cmd/ \| grep -v _test.go` |
| R3b | 同尺含测试（对照） | 16 枚 | `grep -rn '\.Record(' --include='*.go' internal/ cmd/ \| wc -l` |
| R3c | `Record` 定义在不在 | **在**，`internal/tools/task.go:123` | `grep -rn 'func (r \*TaskRoster) Record' --include='*.go' internal/` |

**"查无此项"这一支单独交代**：委托语与派单 §1 都留了"若 `RunAsync` 压根不存在，就要把'零调用点'改写成'这枚函数不存在'"的岔口。**盘上不是那一形**：函数存在、签名是 `RunAsync(ctx context.Context, input string) *RunningTask`、`internal/agent/loop.go:321`，注释（`:320`）逐字写着 "spawns the task on a roster-named goroutine and returns immediately"。**是"存在、有实现、有测试、生产零调用者"**，不是"查无此项"。⇒ 票 176 票面第 11 行那句"RunAsync（环路里那个'起一个后台东西'的入口）在生产码里零调用点"**逐字成立，不必改写**。

顺带把票面上两处**不准确的措辞**量出来（不是我的判断，是行号对行号）：

- `cmd/wisp/run.go:357` 的注释与 `internal/tools/bridge.go:719` 的注释是 R2 尺形下**仅有的两处 `RunAsync` 生产文件命中**，两枚都是**注释行**、不是调用 ⇒ "非测试命中只有两行注释"那句（票 164 Progress log 09-27 段）今天仍逐字成立。
- 台账 `A356`（09-27 提交 `1ed0bc14`）记下上一枚 `176-a1` 的假交件里含**"`RunAsync` 不存在"**这一条，编排者当时"只采纳可复算的两枚地基读数（RunAsync 不存在／Record 零写者）"。⚠ **本程现量：`RunAsync` 存在（R1＝1）**。那半句地基读数是错的，**建议把 `A356` 的地基更正成"存在·生产零调用点"**——台账我不动（`docs/reports/**` 是零字节区），在此具名上报。

**再加一枚本程现量（票面没问、但它决定 175 的成色）**：

| R7 | 谁给 `TaskOutput.ArtifactPath` 赋值（生产，排测试） | **0 枚** | `grep -rn 'ArtifactPath:' --include='*.go' internal/ cmd/ \| grep -v _test.go` |

⇒ 票 175 的豁免载具 `task.go:253` `box.set(rec.ArtifactPath)` 只在 `rec.ArtifactPath != ""` 那一支里被走到（`task.go:244`），而生产里**没有任何一处能给 `ArtifactPath` 填非空** ⇒ 那枚精确豁免在生产里**今天恒不可达**。这条我**不引 `177-c1` §2 的表**，是 R3＋R7 两把尺自己推出来的；引文只有 `177-shape-a-impl-r1.md:68` 那句"在生产里今天是惰性接线"作为**同向线索**（线索，非证据）。

---

## 2. 第二问：起跑口的合法落点在哪（四条，逐条指回行号）

### Q1 仓里有没有"宿主自己起协程跑一件活并把结果记账"的先例？

**有，而且不止一枚——但它不在 `go func(` 那一层。**

| 尺 | 读数 | 命令 |
|---|---|---|
| 裸 `go func(` 生产命中（ban 1 口径） | **0 枚**（含测试 26 枚，全在测试与探针里） | `grep -rn 'go func(' --include='*.go' internal/ cmd/ \| grep -v _test.go` |
| `observe.Registry.Spawn` 生产调用点 | **13 枚** | `grep -rn '\.Spawn(' --include='*.go' internal/ cmd/ \| grep -v _test.go` |

⇒ 本仓**结构性没有裸协程**：唯一的 `go` 语句在受管入口 `Registry.Spawn` 内部（`internal/observe/goroutine.go:262` 签名、`:283` `go r.run(...)`、`:286` `run` 的注释逐字 **"run is the single sanctioned goroutine body of the codebase"**，recover 边界装在 `:293-320`）。派单 §2 把"说不出真先例就当无先例"当默认，**盘上是反过来的：先例是制度，裸形才是新闻**。

**四枚离起跑口最近的先例**（按可抄程度排）：

1. `internal/agent/loop.go:324` `t.h = l.reg.Spawn("agent-task-"+id, "agent", t.root, func(c context.Context) { defer close(t.done); t.result = l.run(c, id, input) })` —— **这就是那枚"起一件活、把结果记进调用方持有的东西"的模板**，缺的只是"没人从生产调它"（R2＝0）。
2. `internal/agent/loop.go:650` `h := l.reg.Spawn("tool-exec-"+taskID, "agent", root, ...)` ＋ join 段 `:663-670`（`<-h.Done()` 全join，再 `h.Err()` 取被 recover 的 panic 并降成"工具失败不是任务死亡"）—— **并行 fan-out 后逐枚 join 再回填**的形状先例。
3. `internal/memory/writer.go:79` `q.reg.Spawn("db-writer", "memory", nil, ...)` —— 常驻、`root=nil`（detached）那形。
4. `internal/plugin/disposal.go:233` 与 `:296` —— 带 `scope:` owner 的一次性 worker。

**命名硬约束（写手最容易撞的一枚，本程现量）**：`observe/goroutine.go:43-60` 就是 D38 名册的盘上形状——

```go
ResidentNames   = ["ui-sta","audio-capture","hotkey-listener","db-writer","watchdog","log-flusher"]
OnDemandNames   = ["kws-infer"]
PerTaskPrefixes = ["agent-task-","tool-exec-","approval-waiter"]
TemporaryNames  = ["asr-infer","tts-infer","panel-host","retention-job","model-downloader","memory-extract","disposal-worker"]
```

`ClassifyGoroutine`（`:64-82`）落不进任何一档 ⇒ `CategoryUnknown` ⇒ `Spawn` 里 `:269-272` 直接 `slog.Warn("goroutine outside the D38 roster (leak symptom)")`。**名册里没有"后台任务记录员"这一档**。⇒ 合法出路只有两枚：**(a)** 复用 `agent-task-` 前缀（前缀匹配，`"agent-task-"+id+"-record"` 这类名字天然归 PerTask）；**(b)** 新登记一枚 D38 名——那要动 `docs/PLAN.md` 的 D38，**AC#5 零字节区，owner 单独批准**。本程倾向 (a)，但**这是形状裁定、不是实现**，留给落地腿。

### Q2 `internal/agent` 的 D15(3) 落盘层能不能被起跑口复用？命名会不会撞？

`Spiller` 现量面：`internal/agent/spill.go:34` 结构体、`:43` `NewSpiller(dir, b)`、`:68` `CapRaw`、`:84` `Prepare(callID, text)`、`:105` `name := artifactName(callID, s.nextSequence())`、`:150` `nextSequence`、`:184` `artifactName`。环路里的唯一调用点是 `internal/agent/loop.go:701-707`（D15(3) 对**工具结果**的截断落盘）。

**结论：`Spiller` 作为"字节落盘＋stub"的载体可以复用；作为"按 task id 寻址"的命名层不能直接复用。** 两条现量理由：

- **寻址键不同型**：`spill.go:81-83` 逐字写着 callID **"it is model-supplied"**，`:26-29` 逐字写着文件名是"model-supplied tool-call id 的单射编码"、理由是"两枚 id 共享一枚名 ⇒ 一次工具调用覆盖另一次存的产物（票 79, C25）"。后台任务的产物**没有 call id**（它是"一次任务"、不是"一次工具调用"）。委托语引的 `177-c1` 那句"当年无答案"我**不当证据**——盘上自认在 `internal/tools/task.go:86-92` 的注释里，而它给的解答正是 `TaskOutput` 这枚映射（"This record is that mapping, and the host - not the model - fills it"）。
- **会撞**：`spill.go:119-133` 的覆盖语义**刻意允许"同一枚 id 覆盖自己早先的产物"**（retried tool call 那一支）。若起跑口把 **task id 当 callID 传进同一枚 `Spiller` 实例**，task id 与模型 call id 落在**同一命名空间**里 ⇒ 一旦相等就走"覆盖"那一支。加上 `sequence`（`:38` `mu`+`:150 nextSequence`）挂在**实例**上，两个生产者共用一枚 Spiller 会共享序号，`:195-196` 的溢出改道（`cutEncodedID`＋`digestArtifactID`）也随之易位。⇒ **复用载体时至少要 (i) 用宿主铸的、不可能等于模型 call id 的 id 形状，或 (ii) 换一枚 `NewSpiller(dir2, b)` 实例写进不同目录**；本程不替写手挑，标**未定，属落地腿自己裁**。

⚠ **禁区 (ii) 的红线就在这儿最近**：`spill.go:29-31` 逐字已经把立场写进码里——"This is a HOST-INTERNAL artifact write (D34 note(2)): **it is deliberately not a gated tool**"。起跑口若要"让模型能读自己的后台产物"，正确方向是走 `task.output`（L0，已在册）＋宿主内部落盘；**把它做成一枚受门控的 artifacts 写入 Tool 就是 `AGENTS §1.2` 的逐字禁止项**。见 §5(ii)。

### Q3 `cmd/wisp/run.go` 里 `rt.tasks` 与工具注册点的相对位置：`Record` 那一步接在哪儿最不外溢？

现量位置（同一枚文件，行号差＝外溢距离的度量）：

| 行 | 内容 |
|---|---|
| `:213-216` | `tasks *tools.TaskRoster` 字段＋注释（"the only thing task.output reads"、定案②不跨重启） |
| `:361` | `rt.tasks = tools.NewTaskRoster()` |
| `:362-367` | `for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks}) { reg.Register(e) ... }` |
| `:354-360` | 注释逐字：**"⚠ WHO FILLS IT IS NOT THIS TICKET"**、"Until a spawner records into it, every call here answers 查不到这个任务" |
| `:564-582` | `admitTask(taskID string) func()` —— 本组合根**唯一的 per-task 边界** |
| `:586` | `func (rt *agentRuntime) execute(task string) int` |
| `:588-617` | `agent.New(agent.Options{... AdmitTask: rt.admitTask, Registry: observe.NewRegistry(), ...})` |
| `:622` | `res := loop.Run(ctx, task)` —— **今天生产里唯一真正跑环路的那一发，同步** |

**最不外溢的落点＝`execute()` 里那一侧，形状是"调现成的 `RunAsync`、拿它返回的 `RunningTask`、join 完之后 Record"**：

- 不需要新面：`RunningTask` 已经有 `ID`（`loop.go:288`，**导出**）、`Root()`（`:297`）、`Cancel()`（`:300`）、`Wait()`（`:308-315`，注释逐字交代它 join `t.done` **＋** `t.h.Done()`，即 D38e"完成＝join counter 抽干"）、`Pending()`（`:320`）。**起跑口要的"取消口、id、完成信号"三件全在现成的返回值上**，一字节 `internal/agent` 都不用动 ⇒ 天然避开禁区 (i)。
- **id 必须用环路铸的那枚，不能宿主自造**：`:564-570` 的注释逐字说明 `admitTask` 之所以拥有收口，是因为它"是本根里唯一被交到环路 task id 的地方——**环路在 Run 内部生成 id，此点之上没人知道它**"。宿主自造 id ⇒ 掉进同一段 `:573-575` **已经登记的开放端**（"a task id a host-internal caller invents when it dispatches on the bridge directly (no admission, no boundary)... recorded as an open end, not silently folded into this line"）⇒ D47 门没登记、C25 scope 没人关。那**既不是"起跑口"，还额外欠一发绕过**——本程把它标成落地腿必须回避的形状。
- Record 的时机：`Wait()` 返回之后（D38e 已 join）是唯一能同时满足 §3 取消判据与"last writer wins"（`task.go:121-122`）的时点；放在 join 之前就是 §3 红色那一发。

⚠ **本程会不会动 `cmd/**`：不会，一字节未动**（派单 §5 把 `cmd/**` 列为禁改、且没具名放开）⇒ **未动、留给落地腿**。但要看清：这条落点**要求动 `cmd/wisp/run.go`**，而 `internal/tools` 对 `TaskRoster` 是 consume-only（`:340-341` 注释自陈）、`internal/agent` 撞禁区 (i)——**三条路里只有 `cmd/**` 是活的**。所以落地腿的派单**必须具名放开 `cmd/wisp/run.go`**，否则它无处可接、只能去撞 (i)。这一条排在 §next 第一位。

顺带一枚现量供写手安心：G2 腿（`OpenTask|CloseTask` 生产出现点）今天 `want_n 2`／实测 2，两枚是 `cmd/wisp/run.go:564`（注释）与 `:581`（`rt.bridge.CloseTask(taskID)`）。起跑口若在 `admitTask` 的收口里顺带处理"任务终结"，**这两枚计数会变** ⇒ 落地腿要么把 G2 基线连同读数一起更新（那是**声明变更**，需按票面上"改判据"的规矩走），要么设计成不新增 `OpenTask/CloseTask` 字面。本程没替它裁。

### Q4 面板侧（`internal/panel/**` 与泵）有没有一条在跑的"异步产出→回填"管子可以搭车？

**现量：有管子，但不同型，搭不上。**

- 管子在场：`cmd/wisp/run.go:438` `rt.stream = panel.NewStreamLog(panel.DefaultStreamKeys)`；`:439-446` `rt.pump = panel.NewSnapshotPump(panel.PumpSources{Verdicts: rt.liveVerdicts, Mode: ..., Workspace: ..., Results: rt.stream.Chunks, Out: rt.bookPanelSnapshot})`；`:446` `rt.ui.publish = rt.publishPanelSnapshot`。泵的面：`internal/panel/pump.go:145` `NewSnapshotPump`、`:215` `Publish()`、`:282` `NewStreamLog`、`:290` `Append(key, text)`、`:328` `Chunks()`。`internal/panel/**` 本程**只读到 `pump.go` 的函数签名行**（没读实现体、没跑面板测试、`tokens_fourway_test.go` 一行没开——别家常红地界）。
- 为什么搭不上：这条管子搬的是**固定键的快照块**（`DefaultStreamKeys`），起跑口要搬的是**按 task id 寻址的产物字节＋一条 `TaskOutput`**。名册的公开面只有 `Look(taskID)`（`task.go:139`）与 `Count()`（`:152`），而 `Count` 的注释 `:149-152` 逐字写明它只为宿主自查存在、**"giving task.output a 'list everything' mode would smuggle that row back in"**（那张不许捞回的表＝`task.list`，DEFERRED 于 `PLAN.md:1531`）。⇒ **把名册塞进面板泵＝从面板侧把 `task.list` smuggle 回来的风险形状**（这句是**我从注释推的风险，不是现量读数**）⇒ **未定，属落地腿自己裁**。
- 硬约束：票面 AC#5 把 `internal/panel/**` 列为**零字节区**。⇒ 落点**不许**落在面板侧，这条不需要论证。

---

## 3. 第三问：取消语义挂在哪、红／绿两向长什么样（＝票 176 AC#3）

**D31 取消总线在场的面（现读行号）**：

| 面 | 位置 |
|---|---|
| `CancelBus` 是 Bridge 的注入字段 | `internal/tools/bridge.go:78-83`（注释逐字"Cancel is the D31 seam (see cancel.go): the veto poll a running tool"）、`:131`、`:176` |
| 每调用的 ctx 挂载／收口 | `bridge.go:452` `withCancel(ectx, cancelHandle{corr: orDefault(req.CorrelationID, req.TaskID), bus: b.cancel})`、`:459` `defer b.cancel.Complete(...)` |
| 部分完成的账 | `Result.AppliedSteps` 由**工具自己填**（`internal/tools/tool.go:54-59`）、bridge 搬运 `:508`、`cancelText` 取报告 `:622-643`（`:631` 短路：`len(res.AppliedSteps)==0 && !b.cancel.Vetoed(corr)` 就无文案）、`approval/report.go:14-15`＋`:117` `UnreportedSteps = rep.Vetoed && len(rep.AppliedSteps)==0` |
| 取消的退码形状 | `bridge.go:443-444`（`ClassCancelled`＋`OutcomeKindCancelled`）、`:516-525`、`:815`、`:830-831` |
| 任务侧取消口 | `loop.go:300` `RunningTask.Cancel() → t.root.Cancel()`；root 于 `:323` `observe.NewRootFrom(ctx, id)` 铸 |
| "完成"的仓内定义 | `goroutine.go:84-87` 注释逐字"task completion means `Root.Wait()` drains to zero, which is what keeps RSS fall-back (D32) reachable"；`Root` `:88-94`（`Ctx`/`Cancel`/`pending`）、`Pending()` `:119`、`Wait(timeout)` `:132`；`Handle.Done()` `:206`、`Err()` `:210` |

**起跑口落地那天，判据能挂的接缝（按外溢从小到大排，本程只裁形状）**：

1. **`Wait()` 之后才 Record**（`loop.go:308-315` 已把 D38e 的双 join 做完 ⇒ 物理上"归零后才写"）。零新面、零 `internal/` 改动。
2. Record 之前问一次 `t.Root().Pending()`（`goroutine.go:119`）或 `ctx.Err()` —— 同一条判据的断言面，不改语义。
3. 让名册自己变成"取消即拒写"（`Record` 里加终态闸）—— ⚠ 这会**改掉 `task.go:121-122` 逐字声明的 "Last writer wins"**（票 164 已勾语义邻近），属契约轴邻近面。**未定，上报，不由落地腿单方面选**。

⚠⚠ **本程最该带给写手的一枚陷阱（现量）**：仓里**今天已经有一处合法的"取消之后仍然写"**——`internal/agent/loop.go:956-958`

```go
func terminalWriteCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), terminalWriteTimeout)
}
```

调用点 `:666-672`，注释逐字："**a task cancelled mid-execution must still book its in-flight calls instead of leaving them pending forever**"。⇒ **"取消后不许还有人在写"这句判据的射程必须显式排除 `terminalWriteCtx` 那一支**，否则绿色那一发会先把自家唯一合法的收尾读成违规、逼写手去松掉一条真要求。这是 §3 的形状硬约束。

**红／绿两向（成对，写手照这两行做）**：

- **绿**（正向）：起跑口起一枚慢任务 → `Cancel()` → 等到 `Wait()` 返回／`Pending()==0` → 断言 `Roster.Look(id)` 自取消时刻起**不再变化**（要么根本没有记录、要么是一条终态记录），且再 `Roster.Count()` 不增。
- **红**（变异，当场必响）：把 `Record` 从 `Wait()` 之后挪到之前，或删掉 `<-h.Done()` 那一发 join（`loop.go:663-665` 的形状），使取消后仍在写。**不许用"跳过断言"造红，也不许用 mock 代替真名册**（`AGENTS §1.3` 只在接缝注入）。
- **今天响不响：不响，而且不可能响。** R2＝0、R3＝0 ⇒ 生产里没有起跑口、没有可被取消的后台写者 ⇒ 这一格**今天连被测对象都没有**。盘上最接近的资产是 `internal/agent/control_test.go:77 TestDefaultControlCancelsRunningTask`（`:94-95` 断言环路 `Status == StatusCancelled`），但它**一字节不碰 `TaskRoster`** ⇒ **不是这一格的替身**，不许拿它当"取消语义已有覆盖"。⇒ 本表把 AC#3 写成**"待起跑口同批诞生"，未通过、不许翻勾**（票 164 定案⑤登记的正是这枚形状，票 176 AC#3 逐字不许再推回去）。

**归口（一句＋理由）**：这条判据落 **176**、不落 164——**判据必须与它的被测对象同票**：164 已按定案⑤把 AC#4 记成"今天连被测对象都没有"，而对象（会写的后台尾巴）只在 176 落地那天诞生；挂在 164 名下就会再长出票面 §20 那句最坏形状——"三张票都以为自己不归口"。

---

## 4. 门禁三枚（只跑这三枚；没跑第 4 枚）

| 门禁 | rc | 读数 |
|---|---|---|
| `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=207, bans #1-5 cmd/=23, ban #6 frontend/=85, ban #7 internal/tools/=20, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=427, ban #8 cmd/=45`（正向对照 `runtests.sh` 内嵌：top-level PASS=34 FAIL=0 SKIP=0） |
| `sh .scratch/wisp/probes/154/gate-clauses.sh` | **0** | 逐字：**`# 腿数＝14 声明与实测不符＝0`**｜`# 腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`｜`# 腿数断言＝相符（名册上每一腿都记了账）`｜`# 基线过期枚数＝0`｜`# 聚合退码＝0` |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **0** | `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` |

与本程相关的两腿逐字（同一次跑出的读数，锚点＝**起手 `fbcdf441`**，口径＝`gate-clauses.sh` 自己的名册）：

```
# ok   腿=G3 声明=quiet 基线=0枚 实测=0枚
# ok   腿=G2 声明=ring 基线=2枚 实测=2枚
```

**没跑的第 4 枚**：`.scratch/wisp/probes/161/r6/flip-declaration.sh`（跑一次就脏跟踪日志）。**没跑全仓 `./...`**；除起手那一发外只跑过 §4 三枚门禁。G3 腿 pattern 与行号现量＝`gate-clauses.sh:363-366`（`want G3 quiet` / `want_n 0` / `run "G3 Q-56 那一支落地：Loop 上出现收 taskID 的导出方法"` / pattern `'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' '' internal/agent`）——派单 §2 写的"`:363-366` 一带"**逐字对上**。

---

## 5. 两条形状禁区的自证（含**两枚标红**）

**(i) 不给 `internal/agent.Loop` 加收"收 `taskID` 的导出方法"——本程结论不违反，但尺子比禁令窄，必须标红。**

- 自证：R6 现量（我**没用腿的读数、自己跑了同一把 pattern**）：`grep -rnE '^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' internal/agent` ⇒ **0 命中**，与腿报的 quiet／基线 0 一致。
- 我的落点（§2-Q3）拿的是**现成导出的 `RunningTask.ID`（`loop.go:288`）**＋`Wait()`／`Cancel()`／`Root()`，**不新增任何 Loop 方法**。
- 🔴 **标红 1（仪器比规矩窄）**：G3 是按**字面 `taskID`** 抓的正则。把参数写成 `id string`、`correlationID`、或包进一枚 struct，**G3 读数纹丝不动（仍 0）**，而票面 AC#4(i) 禁的是**"收 taskID 的导出方法"那一形**、不是那个拼写。⇒ **"腿没响"绝不等于"合规"**；派给写手时必须把这句一起给，否则最可能的跑歪路径就是"加个方法、改个参数名、G3 安静 ⇒ 当成尺子误报"。这正是票面 §24「不按推荐来的坏处」点名的那一发。

**(ii) 不把宿主内部 artifacts 写入做成受门控的 Tool——本程结论不违反，但离那条线只有一次"顺手加个工具"的距离，标红。**

- 自证：我的落点里 `task.output`（L0，已注册 `run.go:362`）仍是**唯一的读回面**，产物落盘继续走 `Spiller`（`spill.go:29-31` 逐字"HOST-INTERNAL artifact write ... deliberately not a gated tool"），d22scan 的 ban 7（`internal/tools/` 那族）今天 rc=0。
- 🔴 **标红 2（给写手的预防针）**：§2-Q2 说"载体可复用 `Spiller`"，若落地腿为了**让 C25 的 hostPathBox 认下后台产物**而把"写 artifacts"暴露成一枚**受门控 Tool**，那就正面撞 `AGENTS §1.2` 的逐字禁止项。合法方向是**保持宿主内部写入**，豁免沿 177 甲形（`task.go:249-255` 那枚惰性载具）等 175 的常驻正控用例——**这条必须写进落地腿派单，不能留给写手临场判断**。

**另两枚需要 owner 具名放开的写面**（不是禁区，是授权缺口）：`cmd/wisp/run.go`（§2-Q3 唯一活路）、以及**可能**需要 `internal/tools/task.go`（§3 接缝选项 3，会改"last writer wins"）。D38 名若走"新登记一枚"那一支（§2-Q1）还要 `docs/PLAN.md`＝契约变更。

---

## 6. 我没读／没测什么（照实列）

- **没读**：`internal/panel/**` 实现体（只 `grep` 到 `pump.go` 的函数签名行）；`internal/observe/goroutine.go` 的 `run`/`RosterReport` 全文（只读 `:244-300`、`:43-82`）；`frontend/**`、`design/**`（规矩禁读，一行没开）；`internal/panel/tokens_fourway_test.go`（一行没开）。
- **没测**：全仓 `./...`；任何面板/risk 之外的包测试；真机 `wisp run` 端到端（起跑口不存在，端到端本无从跑）；D31 取消在**后台任务**这一路的实际行为（无对象）；`internal/risk` 那发预算红的成因（起手一次性读数，未复跑、未归因）；`ArtifactPath` 的"路径是假的／是枚目录"那一形（票 164 已登记为射程缺口，本程未测）。
- **没算**：票 175 判据本体的形状（那是 175 的腿）；`task.list`／`task.cancel` 的实现（票面 §本票不解决）；跨进程重启名册（票 164 定案②已裁 v1 不做）。
- ⚠ **`risk-level green` ≠ `end-to-end works`**：本程三枚门禁全绿只说明"没引入新形状违规、成对普查言行一致"。**177 那枚精确豁免在生产里仍然惰性**（R3＝0 且 R7＝0），`task.output` 真机上仍然只会回答"查不到这个任务"。
- **同一格有另一枚在飞**：`.scratch/wisp/dispatches/2026-09-27-235x-readonly-176-a1b-...-retry.md` ＋ `.scratch/wisp/probes/176/a1b/logs/`（gate1.log／q1-counts.txt…q7-evdone.txt，本程**只看目录清单、未读内容、未采纳任何结论**）。它的写面按派单是 `-a1b.md`＋`probes/176/a1b/**`，与本表不撞名；**票面 176 的 Progress log 是唯一共享面**，我追加时标题写 `176-a1`。

---

## 7. next=（175＋176 落地腿派之前还缺哪几件，按轻重排）

1. **具名放开 `cmd/wisp/run.go` 给落地腿**——§2-Q3 已量清：`internal/tools` consume-only、`internal/agent` 撞禁区 (i)，**三条路只剩 `cmd/**` 活的**。不放开＝逼写手去撞 (i)。
2. **把 §5 两枚标红逐字抄进落地腿派单**（G3 安静≠合规；artifacts 写入不许做成受门控 Tool）。
3. **更正 `A356` 的地基读数**：`RunAsync` **存在**（`loop.go:321`），不是"不存在"；票面 AC#1 措辞无需改写。
4. **裁 §3 接缝选项 3 要不要**（名册"取消即拒写"会改 `task.go:121-22` 的 last-writer-wins）——它决定 AC#3 判据能不能只靠 (1)(2) 成立。
5. **裁 D38 名册那一枚**：复用 `agent-task-` 前缀，还是新登记（后者＝契约变更、owner 批准）。
6. **交代与 `176-a1b` 的关系**（同一格两枚在飞；先落的那枚把表名占掉之后，另一枚的进度与裁决表怎么算）。
7. 起跑口动码那天**重跑 `./internal/risk/`**：起手那发预算红（1.098 对 1.000 ms）与本程无关，但后台协程会改 C26 的负载形状，没有现量证据。
