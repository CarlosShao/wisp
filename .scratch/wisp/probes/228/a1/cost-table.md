# 228-a1 — 只读普查：球／托盘宿主缺失那一格的**甲／乙代价表**

- 腿：`228-a1`（只读普查，非写腿）。工作树 `D:\work\workspace\projects plans\Wisp`，分支 `dev`。
- 起手锚点：`git log -1` ＝ 见 §0 现核。
- ⛔ 本腿**一律不改代码**，唯一产物＝这份文件。
- ⛔ 本腿**不跑** `go test`／`go build`／`go vet`／`gofumpt`／`d22scan`、不执行任何二进制；尺只有 `grep`／`sed -n`（只读打印）／`Read`／`wc`／`git log|show|cat-file`。
- ⛔ **不选边、不写"推荐"**：定案是 owner 的事。

---

## 0. 状态与口径（现核锚点／每把尺的口径）

〔待填〕

---

## 1. 常驻那条腿今天能挂什么（`cmd/wisp/resident_windows.go` ＋ `internal/proc` 全部文件）

### 1.1 `runResident()` 逐行读数（现核：`cmd/wisp/resident_windows.go` 共 85 行，票面记 `:24/:33/:83` 全部复认未漂）

| 现核行号 | 读数 |
|---|---|
| `resident_windows.go:24` | `func runResident()` —— 无参数子命令那条腿；由 `cmd/wisp/main.go:60 runResident()` 调 |
| `:27` | `buildinfo.ResolveEnv()` |
| `:33` | `rt, err := proc.Boot(env)` —— **常驻腿是 `proc.Boot` 的两枚生产调用者之一**（另一枚是 `cmd/wisp/slo_windows.go:261`，SLO 量具子命令） |
| `:34-41` | `proc.ErrAlreadyRunning` → `SignalExistingInstance` → 打印一句 → `return`（D42#7） |
| `:57-65` | `installLogSink(rt.Layout.DataDir)` ＋ `defer sink.close()` |
| `:66-76` | `defer rt.Shutdown(false)` ＋ 逐步打印失败数（"10 steps, %d failed"） |
| `:81` | 逐字：`"wisp: empty event loop running; the floating ball arrives in ticket 07 (Ctrl+C exits cleanly)"` —— **码内自陈球还没接** |
| `:83` | `reason := rt.RunEventLoop()` 阻塞；返回后 `:84` 打印一句、函数返回、`defer` 里的 `Shutdown` 跑完 ⇒ **进程退出** |

`RunEventLoop` 本体（`internal/proc/boot_windows.go:120-138`，现核）＝**一条主 goroutine 上的裸 `for {}`**：

- `:121` `signal.NotifyContext(..., os.Interrupt, syscall.SIGTERM)`；
- `:126` 每 tick 用 `windows.WaitForSingleObject(rt.Instance.ActivateEvent(), 50)` 泵"第二次拉起"的激活事件，命中只 `slog.Info`（`:127` 逐字「ball bring-to-front lands with ticket 07」）＋ `ResetActivation()`；
- `:131-136` `select` 信号或 `time.After(50ms)` 继续空转。
- ⛔ **它没有任何扩展点**：不是回调、不是 channel、没有接口、不接 `ShutdownHooks`；`boot_windows.go:115-119` 的注释把它写成 "the empty event loop of ticket 03"，并明说"pumps the activation event **without a dedicated goroutine**"。⇒ 想让它"接单干活"就是**改这个 `for` 循环本身**（改 `internal/proc`），或在 `cmd/wisp` 里**另写一条循环取代它**。

### 1.2 `internal/proc` 全包清点（票 201 普查腿自陈"这包它自己没打开"＝本腿主攻面）

非测试文件 **9 枚 / 现核行数**：`boot_windows.go` 154 · `doc.go` 25 · `envfork.go` 308 · `externalsampler_windows.go` 115 · `jobscope_windows.go` 243 · `shutdown.go` 188 · `singleinstance_windows.go` 141 · `systemprocs_windows.go` 96 · `treemetrics_windows.go` 298（`wc -l` 口径＝这 9 枚被扫文件的行数总和 1568，不含 8 枚 `*_test.go`）。

逐枚导出面（尺＝`grep -n "^func \|^type "` 全包非测试文件）：

- `boot_windows.go`：`Runtime{Env,Layout,Job,Instance,Registry,StartedAt}`（`:28-35`）／`BootOption`／`WithLayout`／`WithRegistry`／`Boot`（`:57`）／`(*Runtime).RunEventLoop`（`:120`）／`(*Runtime).Shutdown`（`:144`）。
- `envfork.go`：`Layout`／`LayoutFor`／`TestDataDir`／`SealableRoot`／`PortableDataDirName`／`ApplyPortableOverride`／`DefaultLayout`／`Summary`／`Summarize`／`EnvBadge`。
- `jobscope_windows.go`：`JobScope`／`OpenJobScope`／`Assign`／`StartInJob`／`TreePrivateBytes`／`TreeProcessCount`／`TreePIDs`／`PrivateBytesByPID`／`Close`／`Closed`。
- `shutdown.go`：`ShutdownStep` 1..10／`TaskWaitTimeout = 3s`（`:72`）／`PanelWaitTimeout = 2s`（`:73`）／**`ShutdownHooks`（`:80-89`，8 个槽位）**／`ShutdownOptions`／`StepRecord`／`RunShutdownSequence`（`:112`）。
- `singleinstance_windows.go`：`ErrAlreadyRunning`／`SingleInstance`／`AcquireSingleInstance`／`ActivateEvent`／`ResetActivation`／`Release`／`SignalExistingInstance`。
- `externalsampler_windows.go`／`systemprocs_windows.go`／`treemetrics_windows.go`：都是**采样/度量**（`observe.TreeReader` 的实现），与任务/审批/面板无关。
- `doc.go:19-20` 逐字划界：`no business process management (scheduler), no service/watchdog logic (watchdog)`。⇒ **`internal/proc` 自己声明它不该有任务通路**。

⭐ **本包唯一现成的扩展槽＝`ShutdownHooks`（`shutdown.go:80-89`）**，而且**D38(e) 的球位与面板位已经留好了**：`:82` `StopHotkeyAndKWS func(...) // step 2 (ball)`、`:86` `DestroyPanel // step 6 (panel)`、`:81` `SchedulerClose // step 1 (agent/scheduler)`、`:83` `CancelTasks // step 3`、`:87` `FlushAndCloseDB // step 7`。但今天 `(*Runtime).Shutdown`（`boot_windows.go:144-154`）**只填了 `CloseJob` 一枚**（`:146-148`），其余 7 枚为 nil ⇒ `RunShutdownSequence` 记成 `Skipped`（`shutdown.go:127-132`）并打一行 "shutdown step skipped (module not present)"。⇒ 这一格是**一根线**级别（把钩子填上），前提是钩子背后的东西存在。

### 1.3 三枚待挂物的具名结论：①审批门 ②任务通路 ③面板桥

尺：`grep -rn "proc.Boot("` 全仓（剥 `.scratch`）＝ **2 枚**：`cmd/wisp/resident_windows.go:33`、`cmd/wisp/slo_windows.go:261`。
`grep -rln "wisp/internal/ball" --include=*.go .`（剥 `.scratch`）＝ **1 枚**：`cmd/balldebug/main.go`（复认票面"现量"表第一行）。

| 待挂物 | 常驻腿今天有吗 | 造它的那套东西在哪 | 缺的是"一根线"还是"一整块" |
|---|---|---|---|
| ①审批门（`internal/agent/approval` 的 `Gate`／`Replies`） | **零**：`resident_windows.go` 全文 85 行不含 `approval` import；`RunEventLoop` 内不含 | 装配点只有一枚：`cmd/wisp/run.go:452 rt.gate = approval.New(approval.Options{...})`，其中 `:454 Channels: approval.NewChannels()`（空参），答复侧挂点在 `run.go:609 rt.attachReplyListener(s.reply, s.replyVeto)`；`attachReplyListener` 本体 `cmd/wisp/approval_reply.go:389`，`:393 observe.NewRoot("approval-reply")`，`:421 observe.Default.Spawn("approval-waiter", "approval", root, ...)` | **一整块**：门本身在（`gate.go` 现成），但"谁在常驻腿里造 Gate、把 Prompt 接到球/托盘、把答复路由回 Gate"这条装配链今天**只存在于 CLI 腿的 `assembleRuntime`**（`run.go:329`）。接进常驻＝要么把 `assembleRuntime` 变成两腿共用（动 `run.go`），要么在常驻腿里重造一份（＝两份装配根）。 |
| ②任务通路（`rt.execute` 那一形） | **零**：`RunEventLoop` 不接单、无 tick 队列；`main.go:60` 之后 `runResident()` 返回即退出 | `cmd/wisp/run.go:775 func (rt *agentRuntime) execute(task string) int`，尾巴 `:872 bg := loop.RunAsync(ctx, task)` → `:879 res := bg.Wait()` → `:915 code := runExitCode(...)` → `:938 return code`；一次性 ctx 超时在 `:858-860`；入口 `runTextTask` `:220 return rt.execute(task)` | **一整块**：`execute` 的入参是**一枚已经活着的 `agentRuntime`**（`:224` 那个大结构体，字段含 gate/bridge/pump/store/modes…），常驻腿没有 `agentRuntime`；且 `execute` 一发即 `return`，没有"下一单从哪来"的任何结构。**任务队列/触发源（D12 定时、D9 非语音通道）在 `internal/proc` 的自陈里被点名不属于它**（`doc.go:19`）。 |
| ③面板桥（`internal/panel` 的 bridge） | **零**：常驻腿不含 panel；`grep -rln "wisp/internal/panel" cmd/`（剥测试）＝ `panel_assets.go`／`panel_inbound.go`／`panel_pump.go`／`run.go` 四枚，**全是 CLI 腿或旁支子命令** | 组装点 `cmd/wisp/run.go:526 rt.pump = panel.NewSnapshotPump(panel.PumpSources{...})`；`Replies` 侧另有 `cmd/wisp/panel_pump.go:62` 调 `Queue.LiveApprovals()`（票面"现量"表第 8 行，复认在场） | **一整块**（数据侧）＋**一根线**（渲染/生命周期侧）：packet 组装（`SnapshotPump`）现成、且 `ShutdownHooks.DestroyPanel` 槽位在（`shutdown.go:86`）；缺的是真正持有 WebView2 的那位宿主（`internal/panel/doc.go:7` 逐字「WebView2 lifecycle on the shared ui-sta STA thread; 3-5 msedgewebview2 …」）与 `msedgewebview2` 子进程进 Job 的通路——`rt.Job.StartInJob` 的**生产调用者今天只有 1 枚**（`cmd/wisp/slo_windows.go:500`，SLO 量具），见 §6③。 |

⚠ 顺带量到的一枚事实（不是本票地界，但直接影响 AC#0 的读表）：**L1 那一发的枚举口今天不存在**。尺＝`grep -rn "replyVeto"` 全仓非 `.scratch`：声明 `run.go:146`、读点 `run.go:609`、注释 `run.go:448`，**其余命中全在 `cmd/wisp/approval_reply_201_test.go`（`:97/:175/:574`）＝测试赋值 ⇒ 生产赋值点 0 枚**（复认票面"现量"表第 9 行）。`approval.DefaultChannels()`（`internal/agent/approval/approval.go:162-163`）逐字 `return NewChannels(ChannelBall, ChannelEsc)`，全仓生产调用者＝**0 枚**（尺＝`grep -rn "DefaultChannels("`：命中 `approval.go:162` 定义、`gate.go:108` 包内兜底、`fakes_test.go:253`、`run.go:454` 是**空参** `NewChannels()`）⇒ 票面"**不是没建，是装配处主动关掉了**"这句本腿复认。

---

## 2. 常驻名额与看门狗（`internal/observe/goroutine.go` 全文 ＋ `RosterReport`／`Unknown` 的判据与消费者）

尺与口径：`wc -l internal/observe/goroutine.go` ＝ **425 行**（票面记 `:44` 有 `ui-sta` ⇒ **复认未漂**）。以下行号均为现核。

**在册形状（全文清点）**

- `:31-36` 五枚类别：`resident`／`on_demand`／`per_task`／`temporary`／`unknown`（注释逐字 "NOT in the roster: leak symptom"）。
- `:40` `const ResidentBaseline = 6`；`:43-45` `ResidentNames = {"ui-sta", "audio-capture", "hotkey-listener", "db-writer", "watchdog", "log-flusher"}`（`ui-sta` 在 `:44`）。
- `:48` `OnDemandNames = {"kws-infer"}`；`:52` `PerTaskPrefixes = {"agent-task-", "tool-exec-", "approval-waiter"}`；`:57-60` `TemporaryNames` 7 枚（含 `panel-host`）。
- `:64-82` `ClassifyGoroutine`：先查常驻表→按需表→**前缀**匹配 per-task→临时表→否则 `CategoryUnknown`。⚠ **它是"名字→类别"的纯字符串函数，没有任何"这一枚属于哪个进程角色"的输入**。
- `:88-149` `Root`：`Ctx`＋`Cancel`＋`pending` 计数；`:132-149` `Root.Wait(timeout)`＝单调 `NewTimeout` 轮询 `Pending()`（2ms 步进）⇒ **AC#4 需要的 join 原语现成**。注释 `:130-131` 逐字："The 3s cap used by the D38(e) shutdown sequence lives in proc, not here."
- `:188-196` `RosterReport{Total, Resident, OnDemand, PerTask, Temporary, Unknown []string, ResidentOverBaseline bool}`。
- `:200-210` `Handle`：`Done() <-chan struct{}`／`Err() error` ⇒ **每枚 Spawn 都返句柄，join 面现成**。
- `:230` `var Default = NewRegistry()` —— **进程级单例**，注释逐字 "so the watchdog sees one picture"。`boot_windows.go:72` 里 `rt.Registry = observe.Default` ⇒ 常驻腿的 `rt.Registry` 与 CLI 腿 `approval_reply.go:421`／`config_reload.go:119` 用的 `observe.Default` **是同一枚对象**。
- `:262-283` `Spawn`：`:269` 分类，`:270-273` **`CategoryUnknown` 只做一件事——`slog.Warn("goroutine outside the D38 roster (leak symptom)")`**；`:276` `live[{name,owner}]++`、`:279` `root.addPending(1)`；`:281` 全仓唯一被豁免的 `go` 语句。⛔ **没有拒绝、没有阈值、没有报警对象**。
- `:405-425` `RosterReport()`：按类别把 `Snapshot()` 的计数相加，`:422` `rep.ResidentOverBaseline = rep.Resident > ResidentBaseline`，`:423` 排序 `Unknown`。⚠ **`Unknown` 与 `PerTask` 在这里只是被填进结构体，没有任何比较**：全仓没有 `PerTask > 3` 之类的判据（尺＝`grep -rn "ResidentBaseline|ResidentOverBaseline|PerTask"`，非 `.scratch`，命中仅 `goroutine.go` 自身＋`boot_windows.go:108`＋两枚测试文件）。

**`RosterReport` 到底谁在读（现跑的尺）**

尺＝`grep -rn --include=*.go "RosterReport|ResidentBaseline|ResidentOverBaseline"` 剥 `.scratch`：

- **生产调用者＝1 枚**：`internal/proc/boot_windows.go:108 if rep := rt.Registry.RosterReport(); rep.ResidentOverBaseline {` → `:110` `return nil, fmt.Errorf("proc: resident goroutine baseline exceeded at boot")` → 调用侧 `resident_windows.go:42-45` 打印 `boot failed` 并 `os.Exit(1)`。
  ⇒ **它是"开机时一次性自检"，不是循环看门狗**；`Boot` 之后没有任何东西再查这枚数。
- 测试调用者：`internal/agent/loop_golden_test.go:286/293`、`internal/observe/goroutine_test.go:74/102/106/196/202/268`。
- **`RosterReport.Unknown` 的生产消费者＝0 枚**（`Unknown` 只在 `goroutine.go:419/423` 被写、被 `goroutine_test.go:204` 断言；`cmd/balldebug/main.go:40`、`cmd/wisp/approval_always.go:33`、`internal/memory/writer.go:24` 三处只是**注释里提到**）。⇒ **没有"把 Unknown 判成什么颜色／是否报警"这回事**，能观察到的只有 `Spawn` 那行 WARN。

**"看门狗在跑"这一句：现跑尺的结果＝不成立**

- `ls internal/watchdog/` ＝ **只有 `doc.go` 一枚文件，零实现**。`internal/watchdog/doc.go:18` 逐字 `DEFERRED(watchdog loop/thresholds): implemented by ticket 42.` ⚠ 票面把这行记在 `internal/observe/goroutine.go:18` 一带——**位置错**：现核 `goroutine.go` 全文**零枚 `DEFERRED(` 标记**（尺＝`grep -rn "DEFERRED" internal/observe/` ⇒ 命中只有 `doc.go:22`（JSONL sink/redaction/diagnostics/CostMeter）与 `diagnostics.go:23/28/52/179`（登记表载体本身））。`goroutine.go:23-25` 那句是 "AST scan lands with ticket 08; a grep-level regression test runs meanwhile"。
- `grep -rn --include=*.go "internal/watchdog"` 非 `.scratch` ⇒ 命中全是**注释与一枚 reachability 测试**（`cmd/wisp/config_reload.go:27/31`、`internal/config/manager.go:19`、`internal/models/assembly_reachability_121_test.go:41`），**零 import 者、零调用者**。
- 但**名额名 `watchdog` 今天确实在被用**：`cmd/wisp/config_reload.go:119 rt.reloadHandle = observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)`（票 223 交付的**配置热加载 1s tick**，`config_reload.go:87` 注释逐字 "SPEC-03 §4.3's 1s watchdog tick"；`internal/config/manager.go:18-19` 逐字承认"the resident 'watchdog' goroutine cmd/wisp spawns for `wisp run` … internal/watchdog's own loop is still ticket 42's"）。⇒ 这一枚住在**CLI 腿**（`startConfigReload` 由 `assembleRuntime` 调用），**常驻腿没有它**（`resident_windows.go` 不 import config）。

**唯一真在数 goroutine 的门（不是看门狗，是量具）**

- `internal/observe/thresholds.go`（**冻结件，本腿只读**）：`:36` `goroutineLimitSleeping = 6 // D38b resident baseline`、`:37` `goroutineLimitArmed = 7 // + kws-infer`。
- 判据在 `thresholds.go:156-184` `goroutineVerdict`：`Sleeping` 行 `Pass: rep.GoroutinesMax <= goroutineLimitSleeping, Gate: true`（`:161-162`）；`Armed` 行 `<=7`（`:174-175`）；**其余态 `Limit: "record", Gate: false`**（`:179-182`，Note 逐字 "per-task counts vary; D38b roster report governs"）。
- ⚠ **`GoroutinesMax` 的口径＝总数，不是常驻类别**：`sampler.go:325` `rep.GoroutinesMax = maxInt(..., sm.Goroutines)`，而 `sm.Goroutines` 来自 `sampler.go:365`（及 `:515` settle 同形）`Goroutines: s.reg.Count()`；`goroutine.go:344-352` `Count()` 把 `live` 里**所有** key 的计数相加。⇒ **同一枚"6"在两处含义不同**：`ResidentOverBaseline` 数的是"名字落在常驻表里的总数"，SLO 门数的是"这枚 registry 上活着的 goroutine 总数"。这是 §6① 的直接证据。
- **谁评估这枚门**：`grep -rn "observe.NewSampler"` 剥 `.scratch` ⇒ **2 枚，都在 `cmd/wisp/slo_windows.go`（`:283` in-tree、`:385` out-of-tree）**；`SampleState(` 生产调用者 2 枚（`:387`、`:436`）、`CheckSettle(` 1 枚（`:985`）。⇒ **产品运行时里没有任何东西在评估 SLO 门**；门只活在 `wisp slo` 这把量具与测试里。
- ⚠ 顺带量到：`thresholds.go:166-171` 逐字承认 goroutine "are only readable in-process, so an out-of-tree row reports the measuring skeleton's own roster" ⇒ 外测那一行的 `<=6` **数的不是被侧的那条腿**。

---

## 3. 上限按什么计名（`docs/PLAN.md` D38 那一节现读；⛔ PLAN.md 只读一字不改）

尺＝`Read docs/PLAN.md` 现核。**票面记的段落起于 `:2818`、止于 `:2848`——复认**；但票面记的两枚具体行号**有漂**：

- `PLAN.md:2818` 逐字标题：`### D38 — 并发与线程模型（**原文只有一行，而 D32 的可达性完全依赖它**）`。
- `:2824`（票面复认）逐字：`一个 STA/UI 线程｜拥有：分层窗口 + 托盘 + **WebView2 窗口** + 所有 Win32 消息循环`，末句 `→ 悬浮球绘制与面板**必须在同一 STA 线程**。决定：**共用一个 STA 线程，不建第二个 D2D factory**`。
- `:2831-2832` 常驻 6 枚逐字：`ui-sta · audio-capture(pinned) · hotkey-listener · db-writer · watchdog · log-flusher`；`:2833` 按需 +1 `kws-infer`（`Armed` 态）；`:2834` 每任务 3 枚 `agent-task-<id>`／`tool-exec-<id>`／`approval-waiter`；`:2835-2836` 临时 6 枚。
- `:2838`（票面复认未漂）逐字：`**总数上限：常驻 6 + 每任务 3。** 超出即为泄漏征兆 → 看门狗记录并进 SLO 采样与诊断包。`
- `:2842-2843`（⚠ 票面记 `:2845`，**现核漂 3 行**）逐字：`**任务完成必须等所有派生 goroutine 退出（WaitGroup）才算完成** —— 否则「任务结束了但 goroutine 还在跑」正是 **RSS 不回落的头号原因**（直接关联 D32 的 10s）。`
- `:2852-2866` 关停顺序 1..10 全表（与 `internal/proc/shutdown.go:11-31` 的注释同源）。
- `:2870-2871` 逐字：`先退主进程再等子进程 → **孤儿 msedgewebview2.exe 常驻吃几百 MB**（Job Object 兜住，但仍须显式等待…）`。

**"按什么计名"在文字里读到什么**：

- `:2838` 这枚上限是**分成两堆写的**（常驻 6 ＋ 每任务 3），**它没有一句话规定这两堆按"进程角色"分别计**——它只说了"总数上限"与"超出即为泄漏征兆"。⇒ 票面 §"冻结文字怎么说的"最后一条的判断在本腿现读下成立：**"在一发即退的 CLI 腿里起 `ui-sta` 会不会被看门狗判成泄漏"这件事，`PLAN.md` 没有覆盖**（未定义即停点，见 §6①）。
- ⚠ 文字里那句"→ **看门狗**记录并进 SLO 采样与诊断包"预设了一枚看门狗。本腿 §2 的现跑尺结果：这枚看门狗**零实现**（`internal/watchdog/` 只有 `doc.go`，`DEFERRED(watchdog loop/thresholds): implemented by ticket 42` 在 `internal/watchdog/doc.go:18`）。⇒ **`:2838` 那句"看门狗记录"今天是空指**；`PLAN.md:2274-2281`（D32 第 3 点）里那套"阈值按态查表／连续 3 次回落失败 → WatchdogAlert"同样零实现（`internal/statemachine/events.go:22` 有 `EvWatchdogOverrun`、`:56` 有 `EvWatchdogFatal` 两枚事件名在册，但那只是状态机的边，不是运行时的看门狗）。

---

## 4. D32 的资源口径（25MB／40MB 与 handle 数，**从 PLAN.md 现读**，不引票面转述）

现核锚点：`PLAN.md:2221 ## 16.3 D32 — 资源 SLO 与延迟预算（**整体取代 D18 的两张表**）`。

- **口径本体**（`:2223-2235`，`### 16.3.1 度量口径`）：`:2232-2233` 逐字 `唯一合法口径：Wisp 进程树私有内存 = Σ(常驻主进程 + 所有子进程) 的 Private Bytes`；`:2235` 逐字 `实现方式必须是 Job Object（C30），不是遍历进程快照`（并给了 `JOB_OBJECT_MEMORY_INFO.ProcessMemoryUsed` 这个 API 名）。`:2227-2228` 逐字：WebView2 派生 3–5 个 `msedgewebview2.exe`、各 ~60–120MB，面板一开树内瞬间多 200–500MB。
- **附带指标**（`:2242-2245`）：逐字 `GDI 对象数 · User 对象数 · 句柄数 · goroutine 数 · 线程数` 全部纳入 SLO 采样与诊断包；`:2244-2245` 点名"分层窗口 + Direct2D + WebView2 的典型泄漏是 GDI/User 对象，RSS 上完全看不出来，直到 `CreateWindowEx` 开始失败"。⇒ **这一句是本票"球挂进哪条腿"的直接代价维度**。
- **两个数本体**（`:2247` `### 16.3.2 内存与 CPU 上限（按态查表）`、表体 `:2253-2261`）：
  - `:2255` `Sleeping`（空闲，KWS 关，**无子进程**）＝ **≤ 25MB / ≤ 40MB**（Y／X 两档并存）、CPU ≤ 0.5%、附加硬约束逐字含 `句柄 <300 · GDI <200 · goroutine ≤6`。
  - `:2256` `Armed` ＝ 90/110MB、goroutine ≤7；`:2259` 面板开启 ＝ ≤600MB（含 WebView2 3–5 子进程）；`:2260` 工作峰值 ≤700MB；`:2261` 回落＝`Settling` 触发后 10s 内回到该态上限，且 `必须显式 debug.FreeOSMemory()`。
  - `:2249-2251`：25MB 是"路径 Y（语音子进程）"那一档、40MB 是"路径 X（同进程 cgo）"那一档，**X 放宽是因为 onnxruntime DLL 卸不掉**。⚠ 这两档的取舍由 `PLAN.md:2025`（§15 第 4 项定案）写明**"不在方案里预选路径"、"X/Y 由 S0 spike 判定"** ⇒ **25/40 这枚数今天仍是两档并存，不是已定的一档**。
- **实现侧口径与冻结件读数**（本腿只读，一字未动）：
  - `internal/proc/jobscope_windows.go:125-142` `TreePrivateBytes()` 的实现是**逐 pid 调 psapi `GetProcessMemoryInfo` 取 `PrivateUsage` 相加**（`:223-242`），**不是** `PLAN.md:2235` 点名的 `JOB_OBJECT_MEMORY_INFO.ProcessMemoryUsed`。
  - `internal/proc/treemetrics_windows.go:65-70` 逐字说明主进程**不**被 assign 进 Job（"assigning self to a KILL_ON_JOB_CLOSE job would kill the process at graceful shutdown"），而是作为树根手工并入统计（`:80-88`）。⇒ **口径＝"Job 里的子进程 ＋ 自己"，与 `PLAN.md:2233` 的 Σ(主＋子) 一致，与 `:2235` 点名的 API 名不一致**（事实级，供写腿知道数从哪来）。
  - `internal/observe/thresholds.go`（**冻结件，只读**）：`:19` `memCapSleeping = 25 << 20` ⇒ **冻结件里的 Sleeping 只有 25MB 一档，没有 40MB 那一档**；`:33-34` `gdiLimitAll = 200`、`handleLimitAll = 600` ⚠ **与 `PLAN.md:2255` 的"句柄 <300"不同**——差别写在 `thresholds.go:10-14` 的 Ruling 1 里，逐字：`D32 wrote "handles <300" pre-measurement; the measured layered-window + D2D + DWrite stack alone holds ~420 handles, so every state that carries the ball window stack gates at <600 (measured 420 + margin)`。
  ⇒ **本票最重要的一枚"量过的数"就在这里**：**"带球窗口栈"的句柄数＝约 420（实测，出自 `thresholds.go:11-12` 的 Ruling 1 记录）**，且这枚数是"整个分层窗口＋D2D＋DWrite 栈"的，不是"CLI 腿里多起一枚 STA"的增量。⚠ **口径要点明**：这是**被写进冻结件注释的实测结论**（2026-09-19 编排者裁定），本腿**没有**复跑任何测量（本腿禁跑二进制），也没有在仓里找到那台 spike 的原始读数文件——见 §8。
  - `:36-37` `goroutineLimitSleeping = 6`／`goroutineLimitArmed = 7`；`:156-184` `goroutineVerdict` 只在 Sleeping／Armed 两态 `Gate: true`，其余态 `Limit: "record"`。
  - `:24` `memCapWorkPeak = 700 << 20 // TARGET (gate=false) until S3/S5`——冻结件自己承认 700MB 未验收（与 `PLAN.md:2025` 第⑤支一致）。
- ⚠ **D32 的"句柄 <300／GDI <200／goroutine ≤6"这一行是按"Sleeping 态的常驻进程"写的**，`PLAN.md` 里**没有任何一行**给出"CLI 腿在一发任务期间宿主一枚球"的内存或句柄读数。⇒ 该格在 §7 表里一律写"**没量过**"，只有"带球窗口栈 ≈420 句柄"这一枚是量过的、且它来自冻结件注释而非本腿。

---

## 5. 球侧的收口形状（`internal/ball/*` ＋ `cmd/balldebug/main.go` 量具）

尺＝逐枚打开读全文／关键段；票面记的 `:664-681`、`:86-91`、`:20-23`、`:47`、`:78-85`、`:904` **全部复认未漂**。

**`internal/ball/ball_windows.go`（现核 1000 行级；关键形状）**

- `:49-60` `Events` 九枚回调全是**裸 `func()`**（宿主传进来的闭包），`:35-46` `EventKind` 十枚枚举含 `EvTrayPanel/EvTrayMute/EvTrayPauseWake/EvTrayExit`。⇒ **`internal/ball` 不 import 审批包**（复认票面"现量"表第 5 行：它只把"用户选了哪一项"发出去，语义由宿主定）。
- `:63-78` `Options`：`Registry *observe.Registry`，`:159-161` nil ⇒ `observe.Default`。
- `:81` 注释逐字 `One Ball per process`；`:137-140` `activeBall atomic.Pointer[Ball]` ＋ 单一 `ballWndProc` 回调 ⇒ **同进程第二枚球不可用**；`:949-955` `registerBallClass` 走 `sync.Once`。
- `:170-172` **`ui-sta` 的 Spawn 形状**（票面 AC#4 指的"那一形"）：`b.sta.handle = opts.Registry.Spawn("ui-sta", "ball", nil, func(_ context.Context) { b.sta.start(...) })`。⚠ **root 传的是 `nil`** ⇒ 这一枚**不挂任何 `observe.Root` 的 pending 计数**，只能靠 `Handle.Done()` join，`Root.Wait()` 对它无效。
- `:196` `exStyle = wsExLayered|wsExTopmost|wsExToolWindow|wsExNoActivate`、`:212` `wsPopup`（`:197-207` 注释解释了为什么必须是 POPUP）。
- `:306` `SetState(s statemachine.State)` ＝ D43 态的驱动入口（`balldebug` 全程用它：`cmd/balldebug/main.go:287/317/352/489/498/505/511/517/523/632`）。⚠ 这就是票 201 AC#6b"球真进那一态"要动的**唯一现成入口**。
- `:664-681` `wmAppTray`：左键 → `OnTrayPanel`；右键 → `showMenu(...)` 返回值 `menuOpenPanel/menuMute/menuPauseWake/menuExit` 四分支各 `b.fire(...)`。⚠ `:723 func (b *Ball) fire(fn func())` 是回调派发点。
- `:904-937` **`Close()` 的正确收口形状**（票面记 `:904`，复认）：`closed.Swap(true)` 幂等 → `PostTask` 里 `stopAnimTimer`+`stopLiquidTimer`+`unregisterAll(hwnd)`+`tray.remove()`+`rend.release()`+`DestroyWindow`（`:924` 现在**检查返回值并 slog.Error**）→ `<-done`（`:933`）→ `:934-936` **`<-b.sta.handle.Done()` join ui-sta 线程**。⇒ 注释 `:903` 逐字"After Close returns, the process holds no ball-side USER/GDI/D2D objects"——**球侧自己是把收口做完的，缺的是宿主有没有调它**。

**`internal/ball/sta_windows.go`（162 行）**

- `:47` `start` **阻塞直到消息泵退出**：`:52-53` `runtime.LockOSThread()`、`:62` `CoInitializeEx(STA)`、`:77-93` `GetMessage` 循环、`:80-82` `c == 0`（WM_QUIT）才 `break`、`:95` `releaseCOM()`。票面记 `:47` 复认。
- `:127-143` `PostTask`（跨线程投递，`hwnd == 0` 时**就地执行**并注释说明是为了不让调用方死等于死线程）；`:156-162` `quit()`＝`PostMessage(wmNull)`＋`PostQuitMessage(0)`。

**`internal/ball/tray_windows.go`（109 行）**

- `:17-24` `trayUID = 0x5701` ＋ 四枚命令 id `menuOpenPanel=1 / menuMute=2 / menuPauseWake=3 / menuExit=4`（票面 `:20-23` 复认）。
- `:86-91` 四枚标签逐字 `"打开面板" / "静音" / "暂停唤醒" / "退出"` ⇒ **这四处是字符串字面量，按票面禁区那条属 d22scan ban #8 射程**（本腿不跑仪器，只点名形状）。
- `:5-6` 文件头注释逐字 `left click = open panel (no-op stub until ticket 33)` ⇒ **票面"托盘有个空 stub 那句是错的"这一条在本腿得到双重复认**：ball 侧只有注释，真正那句 `fmt.Println("tray: open panel (stub, ticket 33)")` 在**消费者** `cmd/balldebug/main.go:199`。
- `:62-67` `remove()`＝`Shell_NotifyIconW(NIM_DELETE)`（由 `Close()` 在 STA 线程上调）。

**量具 `cmd/balldebug/main.go`（663 行，⛔ 本腿只读未执行）**

- `:78-85` `handleCount()`＝`kernel32!GetProcessHandleCount`（`:76`）；`:174` 启动即打一行 `balldebug: start handles=%d`。`:8-9` 文件头逐字：它打印 handle 数是"so the window-stack SLO gate (<600, docs/SLO.md) is measurable" ⇒ **量具现成但住在 debug 件里**（复认票面"现量"表第 3 支）。
- `:33-47` 自陈：两枚自有 goroutine 名**故意不借产品名册**（`balldebug-hotkey-bridge`／`balldebug-level-feeder`），后果逐字 "one 'goroutine outside the D38 roster' WARN per spawn plus an entry in RosterReport.Unknown"，并点名"要消音得改 `internal/observe/goroutine.go` 的 `TemporaryNames`，那不是本票的决定"。⇒ **这是本腿在 §2 之外独立复认"Unknown 没有报警消费者、只有 WARN 日志"的一枚现场证据。**
- `:188-211` 宿主接法：`statemachine.New` ＋ `ball.New(Options{Events: ..., Registry: observe.Default})`，九枚回调全是 `fmt.Println` 或 `gesture(b, m, kind)`（`:591`）；`:234/266/358` 三条退出路径都调 `b.Close()`。⇒ **今天全仓唯一一份"怎么当球的宿主"就是这枚 debug 件**（写腿若接球，参照系只有它）。

---

## 6. 必须具名答的三问（AC#0 前置，逐问给"读数／无读数"）

### ① 常驻名额上限怎么算——按进程角色分别计，还是全树一个数？CLI 腿里起 `ui-sta` 会不会被判泄漏？

**读数（代码侧，具名到行）**：

- **没有任何一处按"进程角色"计名。** 尺＝`grep -rn "ResidentBaseline|ResidentOverBaseline|ClassifyGoroutine"` 剥 `.scratch` ⇒ 判据只有 `internal/observe/goroutine.go:422`（`rep.ResidentOverBaseline = rep.Resident > ResidentBaseline`）这一枚比较；`ResidentBaseline` 是包级 `const = 6`（`:40`），`ClassifyGoroutine` 的**唯一入参是名字字符串**（`:64`），签名里没有角色／进程类型／env。⇒ 所谓"上限"是**一枚全局数字**，且它的计数口径是"名字落在 `ResidentNames` 那六枚里的存活数"（`:410-411`）。
- 计数载体是**进程级单例 registry**：`:230 var Default = NewRegistry()`；`internal/proc/boot_windows.go:72 rt.Registry = observe.Default`。⇒ 同一进程内两条腿（若并存）共用一张表；**跨进程各数各的**，没有任何"全树一个数"的实现（`RosterReport` 只读本进程 `live` map，`goroutine.go:405-407`）。
- ⚠ **同一枚"6"有两套互不相同的口径**（本腿认为这是 AC#0 最该被 owner 知道的一枚事实）：
  1. `ResidentOverBaseline`＝**只数常驻表命名列**（`goroutine.go:422`）；
  2. SLO 的 `goroutines <=6`＝**数该 registry 的存活总数**（`internal/observe/thresholds.go:161-162` 用 `rep.GoroutinesMax`，而它来自 `internal/observe/sampler.go:365 Goroutines: s.reg.Count()`，`goroutine.go:344-352` `Count()` 把所有 key 相加）。
  ⇒ 口径 2 会把 `approval-waiter`／`tool-exec-`／`subagent-finish-*` 一并算进那 6 枚里；口径 1 不会。
- **"CLI 腿里起 `ui-sta` 会不会被判成泄漏"——按代码今天的形状：不会有任何东西来判它。** 三条现跑尺：
  - `ui-sta` **在册**（`goroutine.go:44`）⇒ 不触发 `CategoryUnknown` 的 WARN（`goroutine.go:269-273`）；
  - 唯一会返回"泄漏"结论的生产代码是 `boot_windows.go:108` 那枚**开机一次性**自检，而 `grep -rn "proc.Boot("` ⇒ **`cmd/wisp/run.go` 不在调用者名单里**（调用者只有 `resident_windows.go:33`、`slo_windows.go:261`）⇒ **CLI 腿今天没有任何判据路径**；
  - `wisp run` 里常驻名列今天**已经在用 2～3 枚**：`log-flusher`（`observe/logging.go:96`，由 `installLogSink`→`observe.InitLog` 触发，`cmd/wisp/logsink.go:149`；`run.go:195` 与 `resident_windows.go:57` 两腿都装）、`watchdog`（`config_reload.go:119`，**只在 CLI 腿**）、`db-writer`（`memory/writer.go:79`，**首次写才懒启动**）。⇒ 加一枚 `ui-sta` 后口径 1 约 3～4／6，**不会**触发 `ResidentOverBaseline`；但**口径 2 在 `Sleeping`／`Armed` 两态是硬门**（`Gate: true`），常驻腿若同时把 `approval-waiter` 常设成"随时可答复"，总数就会是 6＋1＝7 > 6 ——这一枚是**真数**不是误报，但它只在有人跑 `wisp slo` 时才成立（见 ②）。
- ⛔ **判不了的部分（明写）**：**"这 6 枚该按角色分别计还是全树一个数"这个问题，`PLAN.md` 与代码都没有覆盖**——`PLAN.md:2838` 只写"总数上限：常驻 6＋每任务 3"（见 §3），代码里没有"角色"这个输入。要定案缺的是**"票 42 那枚看门狗按什么口径数"这一行读数**（今天不存在，`internal/watchdog/` 只有 `doc.go`）。**本腿不按口味替它填。**

### ② 看门狗会不会误报？谁读 `RosterReport`、阈值从哪来、超了以后发生什么、生产调用者点数

- **生产调用者点数（现跑的尺）**：
  - `grep -rn --include=*.go "RosterReport" | 剥 .scratch` ⇒ **生产 1 枚**：`internal/proc/boot_windows.go:108`；其余 8 处命中全在测试（`internal/agent/loop_golden_test.go:286/293`、`internal/observe/goroutine_test.go:74/102/106/196/202`）与 3 处注释（`cmd/balldebug/main.go:40`、`cmd/wisp/approval_always.go:33`、`internal/memory/writer.go:24`）。
  - `grep -rn "observe.NewSampler"` ⇒ **生产 2 枚**，同在 `cmd/wisp/slo_windows.go`（`:283` 内测、`:385` 外测）；`SampleState(` 生产 2 枚（`:387`、`:436`）、`CheckSettle(` 1 枚（`:985`）。
  - `ls internal/watchdog/` ⇒ **只有 `doc.go`，零实现**；`grep -rn "internal/watchdog"` ⇒ **零 import 者**。`doc.go:18` 逐字 `DEFERRED(watchdog loop/thresholds): implemented by ticket 42.`
  ⇒ ⛔ **"看门狗在跑"这一句今天不成立**；本腿按现跑尺具名写：**周期性看门狗＝零调用者、零实现**。
- **阈值从哪来**：两枚来源、都是编译期常量——`goroutine.go:40 ResidentBaseline = 6`（给 boot 自检用）；`thresholds.go:36-37 goroutineLimitSleeping = 6 / goroutineLimitArmed = 7`（**冻结件，本腿一字未动、只读**）。⚠ **没有任何一处从 config 读阈值**（`internal/config/schema.go:605` 有 `SLOSampleIntervalSec` 这一枚采样间隔字段，但它的消费者只有量具侧；阈值本体不在配置里）。
- **超了以后发生什么**（逐条具名）：
  - 开机自检超 ⇒ `boot_windows.go:110` 返回 error（并先 `job.Close()`，`:109`）⇒ `resident_windows.go:42-45` 打 `"boot failed"` 并 **`os.Exit(1)`** ⇒ **拒绝启动**，这是今天唯一"超了会怎样"的真实后果。
  - 名字不在册 ⇒ `goroutine.go:271-273` **只打一行 `slog.Warn`**；无拒绝、无报警对象、无计数（`Unknown` 生产消费者＝0，见 §2）。
  - SLO 门不过 ⇒ `thresholds.go:162` `Pass:false, Gate:true` ⇒ 由 `wisp slo` 折成退出码（`cmd/wisp/slo_windows.go:290 exitCode := 0`、`:337 exitCode = 1`、`:339 return exitCode`）⇒ **只有量具变红，产品不会有任何动作**。
- **会不会误报（结论级读数）**：**今天不会**——因为根本没有周期性看门狗；会误报的那枚风险是**未来票 42** 的口径选择问题（①里那两套"6"的差别就是它的输入）。本腿**不预测它会怎么选**。

### ③ 谁的进程持有 Job Object

- 尺＝`grep -rn --include=*.go "CreateJobObject|OpenJobScope|AssignProcessToJobObject|KILL_ON_JOB_CLOSE" | 剥 .scratch 与 scripts/spike` ⇒ 生产命中全在 `internal/proc/jobscope_windows.go`（`:72 :77 :83 :100`）与 `boot_windows.go:87`。
- **创建者＝`proc.Boot`**（`boot_windows.go:86-91`），持有者就是调用 Boot 的那枚进程，而 Boot 的生产调用者只有 **2 枚**：`cmd/wisp/resident_windows.go:33`（常驻腿）、`cmd/wisp/slo_windows.go:261`（`wisp slo` 量具）。⇒ **常驻那条腿持有 Job Object；`wisp run` 那条腿今天既不创建、也拿不到。**
- **另一条腿拿不拿得到：拿不到。** `jobscope_windows.go:72` 逐字 `windows.CreateJobObject(nil, nil)`＝**匿名 Job**（名字为 NULL），全仓没有 `OpenJobObject` 调用（尺同上）⇒ 第二个进程无法按名字打开它。⇒ 谁想在 `wisp run` 里让子进程进 Job，就得**自己再开一枚 Job**，而那要经过 `proc.Boot` ⇒ 顺带撞上单实例互斥：`envfork.go:61 MutexName = Local\wisp-single-instance`（`:63 MutexEnabled: true`，仅 test env 关，`:22`），`singleinstance_windows.go:51-54` 遇 `ERROR_ALREADY_EXISTS` 直接返 `ErrAlreadyRunning`。⇒ **常驻实例在跑时，CLI 腿若走 Boot 会被"另一个实例正在运行"拒掉**（这一枚是 §7 表里"乙"支的一项具体代价）。
- **子进程入 Job 的生产通路今天只有 1 枚**：`cmd/wisp/slo_windows.go:500 rt.Job.StartInJob(cmd)`；`grep -rn "StartInJob|\.Assign("` 剥测试与 spike ⇒ **`wisp run` 与常驻腿都没有任何子进程入 Job 的调用点**。⇒ 今天 `wisp run` 起的子进程（若有）**不受 `KILL_ON_JOB_CLOSE` 保护**（`PLAN.md:2237-2238` 点名的"孤儿 `msedgewebview2.exe` 常驻吃几百 MB"那一型，在 CLI 腿没有兜底；`PLAN.md:2870-2871` 同一句）。
- ⚠ **D32 的内存口径依赖这枚 Job**：`jobscope_windows.go:125-142 TreePrivateBytes()` 与 `treemetrics_windows.go:65-70`（主进程不 assign 进 Job，作为树根手工并入）。⇒ **"哪条腿持有 Job"直接决定 D32 那两枚数（25MB／40MB）量的是哪棵树**（§4）。

---

## 7. 甲／乙代价表（本腿唯一产物）

> ⛔ **本腿不选边、不写"推荐"。** 下面每一栏都是"读数＋`file:line`"；读不到的地方一律写"没量过／判不了"。
> 口径先说清：**"枚数"＝本腿按现核调用点列出的、该支一定要碰的文件数**，其中"新建"另计；它不是工时估算。

### 7.0 两支共同的前置（事实级结论，不是选择）

1. **AC#4 那一格两支都要做**：`rt.close()`（`cmd/wisp/run.go:694-707`）今天只 `replyRoot.Cancel()`（`:695-697`）＋`reloadRoot.Cancel()`（`:701-703`）＋`store.Close()`（`:704-706`），**`replyHandle`（`run.go:265` 声明、`approval_reply.go:421` Spawn）与 `reloadHandle`（`run.go:273`、`config_reload.go:119`）都不 join**；`:688-693` 注释自陈逐字 "The cancel is not a join … the goroutine dies with the process"。⇒ 这就是 **D38(c)（`PLAN.md:2842-2843`）那条违约**，**没有任何一支会"顺带"修掉它**（甲支多一处可挂靠：D38(e) 第 3 步 `shutdown.go:158` 的 `CancelTasks`＋`TaskWaitTimeout = 3s`（`shutdown.go:72`）＋`Root.Wait`（`observe/goroutine.go:132-149`），但**这枚钩子今天填它的人＝0**（`boot_windows.go:144-154` 只填 `CloseJob`））。
2. **答复入口枚数今天＝1（控制台 stdin）**：`main.go:147 reply := interactiveStdin()` → `run.go:609 rt.attachReplyListener(s.reply, s.replyVeto)` → `approval_reply.go:421-424` 那个 goroutine 的内容是 `runReplyLoop(ctx, surface, in, rt.stdout)`，本体 `:433-436` 是 **`bufio.Scanner` 阻塞读 `io.Reader`**。⇒ 两支都要新增"原生点选 → `Replies`"那一条入口（现成方法：`replies.go:313 Allow`／`:336 Reject`／`:419 Veto`；gate 侧 `gate.go:626 DecideFromNative`／`:638 DecideFromPanel`／`:387 Veto`）。
3. **L1 那一发的只读枚举口两支撞同一枚位置**（票 220 AC#2 甲乙形）：本腿不重复裁，只点名"谁先造谁"这一格仍悬着。
4. **`internal/panel` 不许放新增可解码结构体**（票面禁区段；`l2_grant_boundary_test.go` 冻结件），两支同此约束。
5. ⚠ **球侧自己是收口完整的**：`Ball.Close()`（`internal/ball/ball_windows.go:904-937`）在 STA 线程上依次 `unregisterAll`→`tray.remove()`→`rend.release()`→`DestroyWindow`（`:924` 已检查返回值）→`sta.quit()`，然后 `<-done`（`:933`）＋ `<-b.sta.handle.Done()`（`:934-936`）；`staThread.start`（`sta_windows.go:47-96`）只在 `GetMessage` 返回 0（WM_QUIT，`:80-82`）时退出泵。⇒ **缺的不是球不会关门，是宿主有没有去敲它**（全仓今天只有 `cmd/balldebug/main.go:234/266/358` 三处调 `b.Close()`，而 `internal/ball` 的生产引用者＝**1 枚＝那枚 debug 件**）。

### 7.1 甲＝把审批门与任务通路接进"常驻"那条腿（贴 D2 原意，`PLAN.md:74/83-88/103`）

| 栏目 | 读数（每栏带 file:line） |
|---|---|
| **今天缺的整块** | ①**装配根**：`assembleRuntime`（`cmd/wisp/run.go:329`）的唯一生产调用者是 `runTextTask`（`run.go:207`），常驻腿（`cmd/wisp/resident_windows.go` 全文 85 行）不 import config/approval/panel/agent/ball，只有 `buildinfo`＋`proc`（`:10-11`）。②**任务入口**：`RunEventLoop`（`internal/proc/boot_windows.go:120-138`）是主 goroutine 上的裸 `for`，无回调／无 channel／无接口；`internal/proc/doc.go:19` 逐字自陈 `no business process management (scheduler)`。③**把任务文本递给已跑实例的载荷通道**：第二枚拉起只会 `SignalExistingInstance`（`resident_windows.go:37`），激活事件是**一枚无载荷 event 对象**（`singleinstance_windows.go:85 ActivateEvent`／`:93 ResetActivation`，`CreateEvent(nil, 1, 0, name)`＝manual-reset 事件，没有名字管道／socket／共享内存），命中后端只是 `slog.Info`（`boot_windows.go:127`）。④**能答复的控制台**：常驻腿自陈"double click the icon, no terminal attached, stderr going nowhere"（`resident_windows.go:47-49`）⇒ 7.0#2 那条 stdin 通路在常驻腿**天然不存在**。⑤**关停钩子**：`ShutdownHooks` 8 枚槽位（`shutdown.go:80-89`，其中 `:82` step 2 归 ball、`:86` step 6 归 panel）今天只填了 `CloseJob`（`boot_windows.go:146-148`）⇒ 7 枚记为 `Skipped`（`shutdown.go:127-132`）。 |
| **要动的文件枚数** | 现核点名 **7 枚改 ＋ 1 枚新建 ＝ 8**：`cmd/wisp/resident_windows.go`（宿主接线＋事件循环取代或扩展）、`cmd/wisp/run.go`（把 `assembleRuntime`/`execute` 抽成两腿共用，或常驻腿引它）、`cmd/wisp/approval_reply.go`（veto channel 声明与原生答复入口）、`cmd/wisp/config_reload.go`（`startConfigReload` 今天只被 `run.go:620` 调，常驻腿要另调）、`internal/proc/boot_windows.go`（循环本体／或新增 hook 面）、`internal/proc/singleinstance_windows.go`（若"第二次拉起带任务文本"要真做就得加载荷通道）、`internal/proc/shutdown.go` 的使用侧（填钩子，顺序本体不许动）＋ **新建 1 枚 cmd/wisp 球宿主文件**（参照系只有 `cmd/balldebug/main.go:188-211`）。⚠ 其中 `internal/proc/**` 与 `cmd/wisp/run.go` 是否与在飞的票 223／224／214／226／220 撞文件，**排程由编排者判**，本腿只点名。 |
| **会不会改退出语义** | **不改 CLI 腿**（`main.go:85 os.Exit(cmdRun(args[1:]))` 保持一发即退）。**改的是常驻腿的"退出原因"集合**：今天只有 `"signal"`（`boot_windows.go:132-133`），托盘「退出」（`tray_windows.go:23 menuExit=4` → `ball_windows.go:677-678 OnTrayExit`）今天**没有消费者可让它关停**（`resident_windows.go` 无 `Events`）。⇒ 这是"新增退出触发者"，不是"改已有语义"。 |
| **goroutine 名额与看门狗代价** | 常驻腿今天真实存活＝**1 枚**（`log-flusher`，经 `resident_windows.go:57 installLogSink` → `observe/logging.go:96`）。甲做完后名册口径 1（只数常驻表名，`goroutine.go:410-411/422`）＝`ui-sta`＋`audio-capture`＋`db-writer`＋`watchdog`＋`log-flusher` ≈ **5～6，正好不严格大于 6**；⚠ `hotkey-listener` 这一枚名册名**全仓生产零 spawn**（尺＝`grep -rn '"hotkey-listener"'` 只命中名册与测试）⇒ 球的热键注册其实发生在 `ui-sta` 线程里（`ball_windows.go:181-183 createOnSTA` 含 hotkeys；`unregisterAll` 在 `Close()` `:913`）。真正要 owner 知道的一枚是**口径 2**：SLO 的 `Sleeping` 行按**存活总数** `<=6` 判（`thresholds.go:161-162` ← `sampler.go:365 s.reg.Count()`），⇒ 若为"随时可答复"把 `approval-waiter`（名册算 per-task，`goroutine.go:52`）**常设**在常驻进程里，Sleeping 那一行的总数就是 6＋1＝**7 > 6 会判红**。⛔ **这个 7 是本腿按名册推的形状，不是实测读数**（本腿禁跑量具）。看门狗代价＝**今天零**（周期性看门狗零实现，§6②），风险全在票 42 还没定的那枚计数口径。 |
| **Job Object 归属会不会变** | **不变，且是甲的既有优势**：常驻腿已经持有（`resident_windows.go:33` → `boot_windows.go:87 OpenJobScope` → `jobscope_windows.go:71-86`），关停第 9 步已经接好（`boot_windows.go:146-148 hooks.CloseJob`；`shutdown.go:183`）。⇒ 面板 WebView2 那 3–5 个 `msedgewebview2.exe` 子进程**有地方进 Job**（`jobscope_windows.go:110-123 StartInJob`；D32 的动机原文 `PLAN.md:2237-2238/2870-2871`）。⚠ 但今天 `StartInJob` 的生产调用者只有 `cmd/wisp/slo_windows.go:500` 一枚，产品两条腿都还没有子进程入 Job 的调用点。 |
| **D38(c) 违约在哪一支被修** | **不在甲顺带修**：见 7.0#1。甲只是**多一处可挂靠**（D38(e) step 3 的 3s bounded wait，`shutdown.go:72/:158` ＋ `Root.Wait`，`goroutine.go:132-149`）；`ui-sta` 那一枚还特别**不能靠 root 计数**——`ball_windows.go:170` 的 `Spawn("ui-sta", "ball", nil, ...)` **root 传 nil** ⇒ `Root.Wait` 看不见它，只能 `Handle.Done()`（`goroutine.go:206-210`）。 |
| **D32 那两个数上的代价：量过／没量过** | **部分量过**（同一份 S0 spike 报告，`docs/evidence/s0/02-spike-report.md` §3.1，口径＝private WS 中位数 MB／`GetProcessHandleCount` 句柄数）：`:105` 空 Go ＝ 6.89MB／handles 142／threads 10；`:109` **＋分层窗口 D2D/DWrite ＝ 15.95MB，GDI 5、User 8、handles 387、threads 21**；`:110` **＋托盘＋热键＋Job ＝ 15.98MB、User 14、handles 407**；`:111` 再渲 30 帧 16.07MB"无泄漏"；`:117` idle-y（cgo＋DLL＋全 shell 栈）16.59/16.45MB。⇒ **"球＋托盘那堆东西"本身的增量是量过的（约 ＋9MB、约 ＋265 句柄、＋11 线程，相对空 Go），甲支吃的是同一枚增量**；⚠ **没量过的是**"装配根搬进常驻腿之后、任务跑完回到 Sleeping 的那棵树"（`TreePrivateBytes` 口径见 `jobscope_windows.go:125-142`），以及 §6① 那枚"总数 7"的形状。⚠ 另一处**数字不一致要点名**：`thresholds.go:11-12` 写"measured … ~420 handles"，而 spike 那一行是 `handles 407`（`02-spike-report.md:110`）——本腿不裁谁对，只报两处读数不同源。 |

### 7.2 乙＝让 `wisp run` 在一发任务期间宿主球并正确收口

| 栏目 | 读数（每栏带 file:line） |
|---|---|
| **今天缺的整块** | ①**等待结构**：`runTextTask` 结尾 `run.go:220 return rt.execute(task)`；`execute` 内部 `:872 bg := loop.RunAsync(ctx, task)` → `:879 res := bg.Wait()` → `:938 return code`，ctx 是一次性超时（`:858-860 context.WithTimeout(..., cfg.LLM.TimeoutMS)`）。⇒ **没有任何"球活着期间这段等待由谁驱动"的结构**，这正是票面"为什么不许就在 wisp run 里挂个球"第 1 支，本腿复认。②**关停序列**：`RunShutdownSequence`（`shutdown.go:112`）的生产调用者只有 `boot_windows.go:149`，而它只被 `resident_windows.go:67` 与 `slo_windows.go:278/307/316/323/330` 触发 ⇒ **CLI 腿今天不走 D38(e) 十步**，乙要么手搓（`PLAN.md:2868-2871` 逐条点名的后果由谁承担要说清），要么让 CLI 腿也 Boot（见下条 Job 归属）。③**球宿主**：`internal/ball` 生产引用者＝1 枚（`cmd/balldebug/main.go`），`cmd/wisp/**` 零 import。④**答复通路形状**：7.0#2 同一条（stdin scanner vs 原生回调）。⑤**`replyVeto` 生产赋值点＝0**（`run.go:146` 声明、`:609` 读、其余命中全在 `approval_reply_201_test.go:97/175/574`）＋ `run.go:454 Channels: approval.NewChannels()` 空参 ⇒ 全部 `loaded=false`；`approval.go:162-163 DefaultChannels()` 生产零调用者。这一枚**是"一根线"**（把 `ChannelBall`/`ChannelEsc` 设上、`approval_reply.go:415-417 SetLoaded` 同一步做），票面"AC#2 那半格可以早做"在本腿读数下成立。 |
| **要动的文件枚数** | 现核点名 **4 枚改 ＋ 1 枚新建 ＝ 5**：`cmd/wisp/run.go`（execute 尾的等待、`rt.close()` 的 join、宿主接线、`replyVeto` 赋值）、`cmd/wisp/main.go`（`cmdRun` `:145-158` 的答复来源与退出码语义）、`cmd/wisp/approval_reply.go`（原生答复入口取代/并存 stdin loop）、`cmd/wisp/config_reload.go`（若关停要等它）＋ **新建 1 枚 cmd/wisp 球宿主文件**。⚠ 若要 Boot 拿 Job，还要动 `internal/proc/**`（**枚数与甲同侧**，见 §6③ 的单实例冲突）。⚠ 雷点具名：`run.go`／`approval_always.go` 与票 223／224／226／214 同文件互斥（票面禁区段），乙支撞文件的概率比甲支更集中在 `run.go` 一枚上。 |
| **会不会改退出语义** | **会，而且这是乙的核心代价**：球要活着 ⇒ `run.go:220` 之后必须挂等待 ⇒ CLI 腿不再"一发即退"。这与两处已定文字相互拉扯：`PLAN.md:74`（D2 拓扑：常驻原生进程树）与 `PLAN.md:1424`（S1 验收写明 `wisp run "..."` CLI 可用）。⚠ 还要新答一个问题：**"任务完成"和"球消失"谁先？** 今天没有任何文字规定 `wisp run` 的球寿命（`PLAN.md:2831-2838` 的"常驻"名额语义在 CLI 腿上无定义，见 §6①）。退出码路径本身可不动（`main.go:85 os.Exit(cmdRun(...))`），但**何时 `os.Exit` 的语义变了**。 |
| **goroutine 名额与看门狗代价** | 每枚 CLI 进程多 1 枚 `ui-sta`（`ball_windows.go:170` Spawn，名在册 `goroutine.go:44` ⇒ **不触发 Unknown WARN**）。⚠ 关键形状：`ui-sta` 的 Spawn **root＝nil** ⇒ 不进 `Root.pending`，`Root.Wait` 无效，只能 `Handle.Done()` join（`goroutine.go:206/132-149`）⇒ 乙的收口必须自己写这一行 join。名册口径 1 在 CLI 腿今天约 2～3（`log-flusher` via `run.go:195`、`watchdog` via `config_reload.go:119`＋`run.go:620`、`db-writer` 懒起 `memory/writer.go:79`）⇒ 加 `ui-sta` 不触 `ResidentOverBaseline`。口径 2（`Sleeping <=6` 总数）在乙支**不适用**（一发即退的进程不是 Sleeping 态），但⚠ **本腿顺带量到一枚既有名册破口**：`internal/tools/subagent_197.go:372` Spawn 的名字是 `"subagent-finish-"+bg.ID`，**不匹配任何 `PerTaskPrefixes`（`goroutine.go:52`）也不在任何表里 ⇒ `CategoryUnknown` ⇒ 每次 spawn 一行 `slog.Warn`（`:271-273`）**。看门狗代价＝**今天零**（同甲，§6②）。 |
| **Job Object 归属会不会变** | **要么不变但失去兜底，要么变且撞互斥——两难是硬的读数**：`wisp run` 今天不 Boot（`grep "proc.Boot("`＝只有 `resident_windows.go:33`／`slo_windows.go:261`）⇒ **CLI 进程既不持有也拿不到 Job**（`jobscope_windows.go:72 CreateJobObject(nil, nil)`＝匿名，全仓无 `OpenJobObject`，见 §6③）。若为拿 Job 而让 CLI 腿 Boot ⇒ 撞上 `envfork.go:61 Local\wisp-single-instance`＋`singleinstance_windows.go:51-54 ErrAlreadyRunning` ⇒ **常驻实例在跑时 `wisp run` 会被"另一个实例正在运行"拒掉**（`MutexEnabled=false` 只有 test env，`envfork.go:22`）。若不走 Boot ⇒ 面板/子进程无 `KILL_ON_JOB_CLOSE` 兜底（`PLAN.md:2237-2238` 点名的孤儿形态无解），且 **D32 的树口径直接失效**：`treemetrics_windows.go:71-74` job 不在／已关 ⇒ `ReadTree` 返 `ClassResource` 错误（`:90-94` 注释逐字 fail-closed："an unmeasurable window must never turn into a silent 0-byte pass"）⇒ `wisp slo` 在乙形态下量不到那棵树。 |
| **D38(c) 违约在哪一支被修** | **不顺带修，而且乙支更急**：GUI 收口是有硬约束的——`sta_windows.go:78-82` 泵只认 WM_QUIT、`sta.quit()` 在 `:156-162`、`Ball.Close()` 的两步 join 在 `ball_windows.go:933-936`；而 `rt.close()`（`run.go:694-707`）今天连一根答复 goroutine 都不等（自陈 `:672-675` 一带"the goroutine dies with the process"，现核文本在 `:688-693`）⇒ **票面那句"进程退出时球窗口是怎么没的，没人说得清"在乙支是必须先还的账**（AC#4 那一格），在甲支则只是"多一处可挂靠"。 |
| **D32 那两个数上的代价：量过／没量过** | **部分量过、且比甲更贴**：同一份 spike 报告的 ④⑤ 基线是**"纯 Go 进程"里加分层窗口＋D2D＋DWrite＋托盘＋热键＋Job**（`02-spike-report.md:109-110`：15.95MB/handles 387 → 15.98MB/handles 407，GDI 5、User 8→14、threads 21）⇒ 这正是"CLI 腿里多一枚 STA 线程＋一块 D2D 面"那一形的读数，`PLAN.md:2025` 第①支的 25MB（路径 Y）在它之内有约 9MB 余量。⚠ **但口径要说清**：那是 `scripts/spike/shell-baseline` 那枚**旁支件**（`scripts/spike/shell-baseline/main.go:169` 有同类 Job 处理）在 S0 的实测，**不是 `wisp run` 本体的读数**；`wisp run` 挂球后连同 config tick／log pipeline／D2D 重渲 30 帧的**总增量＝没量过**。⛔ 本腿**不估数**、**不跑量具**（`cmd/balldebug/main.go:78-85 handleCount()` ＋ `:174` 的打印点是现成量具，只点名）。 |

### 7.3 并排速览（只有本腿能证的部分）

| 维度 | 甲（接进常驻腿） | 乙（`wisp run` 宿主球） |
|---|---|---|
| 改动的文件枚数（现核点名） | 7 改＋1 新建＝**8** | 4 改＋1 新建＝**5**（若要 Boot 再加 `internal/proc/**`） |
| 改不改退出语义 | 不改 CLI 腿；新增常驻腿的退出触发者（托盘「退出」今天无消费者） | **改**（`run.go:220` 之后必须挂等待；球寿命与"任务完成"先后今天无文字） |
| D2 拓扑贴合度 | 贴 `PLAN.md:74/83-88/103`（球、托盘、Job、Agent 核心循环同一枚常驻进程） | 与 `PLAN.md:74` 的拓扑不同源；`PLAN.md` **没有一句按 CLI 子命令分别许可或禁止**（票面原话，本腿复认＝D2 是拓扑图不是许可表） |
| goroutine 名额 | 名册口径 1 ≈ 5～6／6（不严格超）；口径 2 有 6＋1＝7 的**形状**（未实测） | ＋1 枚 `ui-sta`；口径 2 不适用（不是 Sleeping 态）；另见既有 `subagent-finish-*` 不在册一枚 |
| Job Object | **已有**（`boot_windows.go:87`），子进程可进（`StartInJob` 现成但产品零调用点） | **拿不到**（匿名 Job）；Boot 换 ⇒ 撞单实例互斥 |
| D38(c) 违约（`run.go:694-707`） | 不顺带修；多一处可挂靠（D38(e) step 3 的 3s bounded wait） | 不顺带修；且必须先还才谈得上 GUI 收口 |
| D32 代价 | 球栈增量**量过**（15.95→15.98MB、387→407 句柄，spike §3.1）；搬装配根后的 Sleeping 树**没量过** | 同一份球栈读数**形状最贴**（纯 Go 进程＋窗口栈）；但"不是 `wisp run` 本体"，总增量**没量过** |
| 面板（L0/WebView2）那半 | 有地方进 Job；`ShutdownHooks.DestroyPanel` 槽位在（`shutdown.go:86`） | 子进程无兜底；`wisp slo` 的树口径失效（`treemetrics_windows.go:71-74`） |

---

## 8. 我没读到／判不了（诚实列，不编）

〔待填〕
