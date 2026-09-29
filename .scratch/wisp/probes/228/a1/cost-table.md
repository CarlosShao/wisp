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

〔待填〕

---

## 3. 上限按什么计名（`docs/PLAN.md` D38 那一节现读；⛔ PLAN.md 只读一字不改）

〔待填〕

---

## 4. D32 的资源口径（25MB／40MB 与 handle 数，**从 PLAN.md 现读**，不引票面转述）

〔待填〕

---

## 5. 球侧的收口形状（`internal/ball/*` ＋ `cmd/balldebug/main.go` 量具）

〔待填〕

---

## 6. 必须具名答的三问（AC#0 前置，逐问给"读数／无读数"）

### ① 常驻名额上限怎么算（按进程角色分别计／全树一个数？CLI 腿里起 `ui-sta` 会不会被判泄漏）

〔待填〕

### ② 看门狗会不会误报（谁读 `RosterReport`、阈值从哪来、超了以后发生什么、生产调用者点数）

〔待填〕

### ③ 谁的进程持有 Job Object

〔待填〕

---

## 7. 甲／乙代价表（本腿唯一产物）

- **甲＝把审批门与任务通路接进"常驻"那条腿**（贴 D2 原意）
- **乙＝让 `wisp run` 在一发任务期间宿主球并正确收口**

〔待填〕

---

## 8. 我没读到／判不了（诚实列，不编）

〔待填〕
