# 票 46 · S6 两处空洞 + 一句假话：接缝落点与红证据

- 票：`.scratch/wisp/issues/46-*.md` S6 验收；空洞出处 `docs/evidence/s1/46-s6-acceptance.md` §1 AC#2（三条命令）
- 写腿：`agent-61`，2026-11-15；写面 `cmd/wisp/**`、`internal/observe/**`（⛔ `thresholds.go`）、本文件
- 本文件只记「接缝现在落在哪一行、哪枚测试会红、红句逐字」，不重开票 46 的其它验收项

## 0. 一句话账

| 空洞 | 本发落点 | 会红的测试（删掉那一行就红） |
|---|---|---|
| ① 收口闸零生产执行（`cmd/wisp/resident_windows.go` 的 `_ = wd`） | 装配 `cmd/wisp/resident_watchdog_windows.go:80-81`；触发 `cmd/wisp/resident_voice_engines_windows.go:525-552`；启动行 `cmd/wisp/resident_windows.go:889` | `Test46SettleGateVerifiesTheDisposeCascade`／`Test384StartResidentWatchdogRegistersRealLoop` |
| ② `NewDiskGatedWriter` 无生产装配（旧读数：全仓 2 处命中，同一函数自检） | `internal/observe/logging.go:150`（装）+ `cmd/wisp/logsink.go:210`（填 `LogConfig.Disk`） | `Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink` |
| ③ `internal/observe/guards_45.go` 那句「`LogConfig.Disk is set`」是假话 | 字段真造出来了（`LogConfig.Disk *DiskGuard`），句子改为四行可核的组装说明 | ②的红测试同时是③的读法：字段没接上，②立刻红 |

## 1. 空洞①：收口闸现在跑在哪，为什么是那里

**赋值点（三行，逐字）**

- `cmd/wisp/resident_watchdog_windows.go:80` `gate := newResidentSettleGate(loop, tree, rt.Registry, root, stateFn)`
- `cmd/wisp/resident_watchdog_windows.go:81` `rv.attachSettleGate(gate)`
- `cmd/wisp/resident_voice_engines_windows.go:549-551`
  `if gate != nil { gate.request(settleRequest{peak: peak, released: released}) }`

**触发点为什么是「dispose 级联第三步 FreeOSMemory」而不是别处**
D32 那一行问的是「dispose 之后树回到目标态 cap 以下、且释放路径真跑过」。全进程里唯一同时满足「引擎已卸载、堆已交回、面板已/即将卸下」的时刻就是 C11 级联第 3 步（`resident_voice_engines_windows.go:525`）。在它之前采样会把「还没释放的大树」当成 dispose 后的读数；在它之后（比如等某个定时器）则要在没有事件的地方猜一个时刻。级联本身是四步固定顺序，第三步是唯一带「内存已交回」语义的那一步。

**为什么测量不在这条线程上**
`watchdog.Loop.VerifySettle` 会阻塞到整个窗口结束（`internal/watchdog/settle.go`，默认 `SettleWithin` 10s），压在级联线程上等于让 dispose 卡 10 秒。所以 `request()` 只登记（`resident_settle_gate_windows.go`），测量走 `observe.Registry.Spawn("watchdog", "settle-gate", root, …)`——D38b 的 6 枚冻结名册不变，用 owner 记账，先例是 `cmd/wisp/config_reload.go:142`。

**为什么不硬编码 Sleeping**
目标态读 D43 机现态（`stateFn` → `sloStateFromD43Machine`），只有读到非采样态才回落 Sleeping。Armed 的 cap 是 90MB、Sleeping 是 25MB，硬编码 Sleeping 会把一台 KWS 常驻的机器判成每发都不过。

**三件配套决定（都是为了让这枚闸门在健康进程上不误报）**
1. **`released` 用增量**：`resident_voice_engines_windows.go:542` 取 `observe.FreeOSMemoryCount()`，`:548` 算增量，`:544` 的释放走 `freeOSMemoryNow()` → `observe.ReleaseMemory()`（`:257`，插桩包装）。谁哪天把它改回 `debug.FreeOSMemory()`，`released` 就是 false，闸门按契约判失败——这是**故意的红**，不是要圆过去的谎。
2. **合并**：同一时刻只跑一发校验；窗口内的新请求记 `coalesced`。
3. **唤醒取消**：`openScope` 里 `resident_voice_engines_windows.go:596` `r.abortSettleCheck("session reopened: " + reason)`——被唤醒的树已不是 dispose 后那棵，取消**不计失败**（`finish` 里 `cancelled ≠ failures`），否则健康进程会被普通唤醒一路走到 D43#35。

**启动行**（可核的一句）：`cmd/wisp/resident_windows.go:889`
`wisp: 收口闸：装配（触发＝会话 dispose 级联第 3 步 FreeOSMemory；目标态＝D43 机现态，非采样态回落 Sleeping；窗口＝watchdog 的 SettleWithin；证据＝watchdog settle verification 日志行＋三次连续失败升 D43#35）；安全卸载白名单组件：[panel-webview]`
没有 watchdog 循环的进程打的是 `收口闸：未装配（…D32 settle 行在这次启动里没有验证者）`。

## 2. 空洞②：磁盘门控装配在哪，为什么不写「具名不装配」

**装配（两行，逐字）**

- `cmd/wisp/logsink.go:209-210` `guard, guardNote := diskGuardForDataDir(dataDir)` / `obs.Disk = guard`
- `internal/observe/logging.go:149-151` `if cfg.Disk != nil { sink = NewDiskGatedWriter(w, cfg.Disk) }`

探针是 `cmd/wisp/diskfree_windows.go:24` `dataVolumeFreeBytes`（`windows.GetDiskFreeSpaceEx`，非 windows 走 `diskfree_other.go` 返回错误 ⇒ 装上但不门控，安装记录里写 `disk_gated:false` 并带拒因）。阈值一字节没动：`observe.DiskGuardThresholdDefault`（200MB），安装记录里 `disk_guard="volume … 75587309568 bytes free at install, threshold 209715200"`。

**为什么不写「具名不装配，理由是 X」**：那个理由找不到一句真话。票 46 S6 的 AC 说「低于阈值时日志/L3/artifacts 停止写入」；只装配会得到「这条承诺在任何进程里都没有执行者」。性价比也不支持：装配它只花一个探针文件 + 一个字段（`LogConfig.Disk`），而③那句假话**只有把字段造出来才能变成真话**（见 §3）。

## 3. 那句假话改成什么

**旧句**（`internal/observe/guards_45.go`，装配前逐字）：
`LogConfig.Disk is set` —— 而 `internal/observe/logging.go` 的 `type LogConfig` 当时只有 `Dir/Level/RollSizeMB/RollDays/RedactPaths`，`Disk` 零命中；所以那句话描述的是一个不存在的字段，读者只有它自己下面两行自检。

**新句**（同一处，现在是四行可核的组装说明，逐字见 `guards_45.go:338-353`）：
- `LogConfig.Disk (*DiskGuard, logging.go:LogConfig)` 承载闸门；
- `InitLogWithRegistry` 在该字段非 nil 时**恰好**把这枚包装装在 JSON handler 与滚动写入器之间；
- `cmd/wisp/logsink.go` 的 `installLogSink` 是**唯一**的生产填充点；
- 探针在 `cmd/wisp/diskfree_windows.go`。
并明写「旧的这句描述的是一个不存在的字段，所以这个构造器当时的读者只有本包两处自检」。

字段本体（`internal/observe/logging.go:75-88`，逐字节选）：
```go
// Disk is the D42#5 write-admission guard for the log class (guards_45.go).
// Nil is the pre-ticket-45 shape: the file takes every record with no
// pre-write free-space verdict, which is what every caller before the cmd
// wiring leg passed. When it is set, InitLogWithRegistry puts
// NewDiskGatedWriter between the JSON handler and the rolling writer, so a
// refused record is HELD (bounded) and replayed instead of being dropped -
// the sentence guards_45.go's NewDiskGatedWriter doc used to state about a
// field that did not exist until this one.
// …
Disk *DiskGuard
```

## 4. 装配本身暴露出的新缺陷：0% CPU 挂死（重入死锁）

**这是本发最值钱的一件东西**，也是「空洞②为什么必须真装配」的答案：这个 bug 只有在闸门真挂在日志管线上时才存在。

**现象（读数）**：`Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink` 在 1MB 探针下不动了。
`pid=64600 name=wisp.test start=17:49:14 age≈18.6min cpu=0.27s threads=6 memMB=22`，
`-test.v` 最后一条 `=== RUN` 就是它；`-test.timeout=90m` 意味着它还会再挂 71 分钟；`Stop-Process` 后收 `FAIL … 1188.00s`。
**0.27 秒 CPU / 18.6 分钟 ⇒ 它在等，不是在算。**

**机制**：`diskGatedWriter.Write` 是被 slog 的 JSON 处理器**持锁**调用的（`slog` 的 `commonHandler.Handle` 全程持 `h.mu` 并调用 writer）。旧 `DiskGuard.Check` 在拒绝时 `slog.Warn(warning)`；那枚记录又流回同一个处理器 ⇒ 同一 goroutine 对不可重入的 `sync.Mutex` 自锁。**不是 panic，是静默挂死。**

**后果不止测试**：`installLogSink` 自己的安装记录（`slog.Info("wisp: persistent log sink installed"…)`）就足以触发第一发拒绝。也就是说**真机磁盘一旦低于 200MB，下一次启动会挂在「装日志」这一行**，而且不报错、不留日志（日志正是被挂住的那条链）。

**修法（两处，都很小）**
1. `internal/observe/guards_45.go:201` 新增 `admit(c)` = `check(c, false)`：同一套准入判定与同一套状态迁移，**不发记录**；`diskGatedWriter.Write`（`guards_45.go:362`）改用它。`Check` 的语义一字不改（它的读者——包括本包 11 处测试调用——照旧）。
2. 警告改由**管线自己的 flusher** 发：`internal/observe/logging.go:169` `w.eachTick = p.publishDiskState`，`logging.go:189` `publishDiskState` 每拍读 `guard.Status()`，**只在跨越时**发一枚（tripped→一枚 WARN；恢复→一枚 WARN）。flusher 是唯一一个「永远不在 Write 内部」的协程，所以它发的记录能安全地回到管线；而且这枚记录在盘满时会被**同一枚 gated writer 扣住**，恢复后随尾巴一起落到文件里——「解释这段空白的那一行」正好在能写的时候写进去。
   - **口径修正（我自己先读错、后按 diff 改正）**：`Check` 原本就是「只在跨越时发一枚」（`warning` 只在 `low && !g.tripped` 那支被赋值，still-tripped 支只刷新 `g.warning`）。所以这里**没有**「逐调刷屏」要修——要修的**只有发的地点**，而地点正是死锁的成因。
3. 为了让测试能看见「扣住 ≠ 还没写」，`internal/observe/logging.go:257` 新增 `LogPipeline.DiskHeldStats()`（扣住字节 / 丢掉字节 / 已回放字节 / 有没有闸门）。盘空时一个被缓冲的记录和被扣住的记录在文件系统上一模一样，这个读数是唯一的分辨手段。

**红向仍然在**：`Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink` 的第一条断言就是 `gated == false` 就红（见 §6 的红句逐字）。

## 5. 新增/收紧的断言（没有一条放松）

| 断言 | 位置 | 删掉什么会红 |
|---|---|---|
| 装配函数必须返回闸门 | `cmd/wisp/resident_watchdog_windows_test.go:376` | `startResidentWatchdog` 不再 `return gate` |
| 闸门必须**挂到** residency 上（新） | 同上 `:386` | `resident_watchdog_windows.go:81` 那行 |
| 级联必须发出一发请求、peak 是读到的真值、目标态是机现态、loop 的失败计数 +1、插桩计数有增量 | `cmd/wisp/resident_settle_gate_46_windows_test.go`（`Test46SettleGateVerifiesTheDisposeCascade`） | `:550` 的 `gate.request(...)`；把释放改回 `debug.FreeOSMemory()`；`readPeak()` |
| 目标态不是硬编码（Armed≠Sleeping，非采样态回落 Sleeping） | 同上（`Test46SettleGateJudgesTheStateTheMachineIsIn`） | 把 `settleTarget()` 写成常量 |
| 唤醒取消 ≠ 失败 | 同上（`Test46WakeCancelsTheSettleVerification`） | `:596` 的 `abortSettleCheck` |
| 启动行必须点名触发/目标态/窗口/证据；未装配要明说 | 同上（`Test46SettleGateLineNamesItsTrigger`） | `line()` 改回空串或去掉关键字 |
| 盘满时安装记录被扣住（`gated`、`held>0`、`lost=0`、文件里没有）；恢复后回放且时序正确；健康对照立刻落盘；真探针给出真实空闲字节且阈值一字节未动 | `cmd/wisp/logsink_diskgate_46_windows_test.go`（2 枚） | `logsink.go:210` 的 `obs.Disk = guard`；`logging.go:150` 的包装；`FlushNow` 或回放路径 |

`internal/observe` 侧没有修改任何断言（`go test ./internal/observe -count=1` → `ok … 5.057s`，含 `guards_45_test.go` 全部 11 处 `Check` 调用）。

## 6. 我自跑的读数

### 6.1 新增测试（绿）

```
=== RUN   Test46SettleGateVerifiesTheDisposeCascade
    resident_settle_gate_46_windows_test.go:149: 46-READ cascade -> settle-gate requests=1 coalesced=0 runs=1 cancelled=0 failures=1 last_target=Sleeping last_pass=false back_ms=-1
--- PASS (0.13s)
=== RUN   Test46SettleGateJudgesTheStateTheMachineIsIn  --- PASS (0.13s)  [armed_machine PASS 0.06s / no_sampling_row PASS 0.06s]
=== RUN   Test46WakeCancelsTheSettleVerification  --- PASS (0.00s)
=== RUN   Test46SettleGateLineNamesItsTrigger  --- PASS (0.00s)
=== RUN   Test384StartResidentWatchdogRegistersRealLoop
    resident_watchdog_windows_test.go:383: 46-READ settle gate armed: 收口闸：装配（…）
--- PASS (0.00s)
=== RUN   Test46DiskGuardIsAssembledWithTheRealVolumeProbe
    logsink_diskgate_46_windows_test.go:93: 46-READ disk guard: volume …\001, 75587309568 bytes free at install, threshold 209715200
--- PASS (0.00s)
=== RUN   Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink  --- PASS (0.03s)
```

真校验那一发的 watchdog 日志（Test46SettleGateVerifiesTheDisposeCascade 里，逐字）：
`watchdog: settle verification target_state=Sleeping pass=false back_within_cap_ms=-1 final_bytes=419430400 cap_bytes=26214400 free_os_memory_count=1 sample_errors=0 consecutive_failures=1 error_class=resource`
——`final_bytes=419430400` 就是夹具给的 400MB，`cap_bytes=26214400` 是 Sleeping 的 25MB，`back_within_cap_ms=-1` 是「窗口内没回到 cap」。

### 6.2 红齿（把那一行临时拆掉，跑一次，红了，再装回）

| 拆掉什么 | 红句（逐字） | 耗时 |
|---|---|---|
| `cmd/wisp/logsink.go:210`（`obs.Disk = guard` 改成 nil） | `46: this pipeline has no gated writer at all (LogConfig.Disk never reached InitLog): nothing is holding bytes, so the log class is ungated`（`logsink_diskgate_46_windows_test.go:130`） | 0.03s |
| `cmd/wisp/resident_watchdog_windows.go:81`（`rv.attachSettleGate(gate)`） | `46: the residency holds settle gate 0x0, want the one this assembly built (0x362153830210): the dispose cascade would verify nothing`（`resident_watchdog_windows_test.go:391`） | 0.00s |
| `cmd/wisp/resident_voice_engines_windows.go:550`（`gate.request(...)`） | `46: the settle verification the cascade should have started never happened within 10s`（`resident_settle_gate_46_windows_test.go:122`） | 10.00s |

### 6.3 挂死那一发（修前，作为发现记在这里）

```
=== RUN   Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink
FAIL	github.com/CarlosShao/wisp/cmd/wisp	1188.00s
```
pid 64600、18.6 分钟、0.27s CPU、`Stop-Process` 收尾。修后同一条 0.03s 绿（§6.1）。

### 6.4 整包（A1882 三件）

**结论先说：三发整包都没跑到终点锚，三发都是 native AV（退出码原值已抓到），本发⛔不声称整包绿。**

| 发 | 形态 | 时长 | 退出码（原值） | 最后一条 `=== RUN` | 终点锚 `TestZZ315A1` | `--- FAIL` |
|---|---|---|---|---|---|---|
| A | 搬运工 `wisp-cli-staged-tests.sh` | 393 s | 未抓（载具合成 `FAIL`） | `TestAC4PriorFocusSurvivesARefusedPanelSample` | 0 命中 | 0（352 枚 `--- PASS`） |
| B | 搬运工 | 425 s | 未抓 | `Test145ArgsJSONComesFromTheBookedForensicsRow` | 0 命中 | 0（357 枚 `--- PASS`） |
| C | **照载具的形直跑**（cwd=`cmd/wisp`＋三个 pin DLL 在二进制旁＋5 个 `WISP_*`，`Start-Process -Wait -PassThru` 重定向到文件、⛔ 管道） | ≈420 s | **−1073741819 ＝ `0xC0000005`（native ACCESS_VIOLATION）** | `Test145SinkFeedLinesCarryTheEventsIntoTheRunKey` | **0 命中** | **0**（353 枚 `--- PASS`） |

- A1882 ①**终点锚**：三发都**没跑到** `TestZZ315A1`。
- A1882 ②**最后一条 `=== RUN`**：见上表（三发三个不同名字）。
- A1882 ③**本波新测试在截断点的哪一侧**：`Test46DiskGuard*`／`Test46DiskGatedWriter*` 在**截断点之前**，两发都 `--- PASS`（读数为 §6.1）；`Test46SettleGate*`／`Test46Wake*`／`Test384StartResidentWatchdogRegistersRealLoop` 在**截断点之后**，整包里本轮**没跑到**——它们的绿来自 §6.1 的定向发（`-test.run=…`，一次一枚）。
- 归属：退出码 `0xC0000005` 是本机已知族，台账 `A1877`（今天 13:2x，**早于本发改动**）记的就是同一枚，且明写「AV 没被 WER 捕获、本机 CrashDumps 空」。本发三发的 Application 日志在 18:24／18:32／18:41 三个时刻**零** `Application Error(1000)` 事件（同日本机 WER 在记别的），与 A1877 那句一致 ⇒ **不能拿"没有 WER 记录"当"不是 AV"**，也不能反过来把它算成本发引入。
- 三发的死亡时刻都是 **393／425／≈420 s**、死点名字各不相同 ⇒ 与**时间**相关、与**具体用例**无关（A1877 的原话是「死在收尾段」）。⇒ 整包的"绿"必须以**跑到锚点为前提**，本发没有这个前提，故只报"截断在第 N 秒、最后一条 RUN 是谁"。

存档：`.scratch/wisp/probes/orch/20261115/46-full-log.txt`（A 发）、`46-run3-out.txt` / `46-run3-err.txt` / `46-run3-exit.txt`（C 发）。

## 7. 代价、边界、没做的

1. **族**：一发真校验在窗口里是第 7 枚协程（`Registry.Spawn("watchdog", "settle-gate")`）。D38b 的 6 枚冻结名册一字节没动，用 owner 记账；代价是那 10 秒里 `ResidentOverBaseline` 会读真、真机 SLO 采样若正好落在窗口内会看到 +1 协程。这是「窗口内真测量」的必然账，不是泄漏（窗口结束协程就退）。
2. **没做（有理由）**：不把 flusher 的每拍 `Flush` 接到 `drain()`。`drain()` 不查闸门（它的契约是「flush 不丢字节」），若每拍调它，盘仍低于阈值时缓存会被倒进盘——等于把 D42#5 的「低于阈值日志停止写入」取消掉。现在的回放时机是：**第一发被放行的写**（`guards_45.go` 的 `Write` 路径先 drain 再写）或显式 `FlushNow`/`Close`。
3. **留给票 45 的旧账**：`DiskGuard.Check` 仍在**跨越时**发一枚记录（这是它原本的语义，一字未改，本包 11 处测试调用照旧走它）。装配后的产品里，日志类只走 `admit`（quiet），发记录的是 flusher；两处不会重复发同一枚，因为一条记录只有一个来源：盘满/恢复的那一枚由 `publishDiskState` 从 `Status()` 的跨越里发，`Check` 只有「非日志类调用者」才会走。要不要让 `Check` 也把「发」交出去，是票 45 的判定面，本发不动它的语义。
4. **未验**：真机「盘真低于 200MB」那一发没跑（会把机主的盘撑满，且要用真机窗口）。已验的是：真探针读数（75.5GB 空闲、阈值 209715200）+ 夹具翻转下的扣住/回放/对照。真机那一发的判据留给 §8 的复跑命令 + 真机窗口。
5. 本机三条**不属本发**的已知红照旧：`TestAC4FocusReturnToPriorWindowGap33r5`、`TestAC246DevLegIgnoresTheTestTaskInjection`、`Test320FetchLegWritesTheKeyNowhereOnDisk`。

## 8. 复跑命令（逐条可核）

```sh
# 新增的六枚 + 装配函数那一枚（一次一枚，经搬运工）
"D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh "-test.run=Test46" "-test.timeout=180s"
"D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh "-test.run=Test384StartResidentWatchdogRegistersRealLoop" "-test.timeout=180s"
# observe 侧（含 guards_45 全部旧断言）
go test ./internal/observe -count=1
# 整包（同一时刻只跑一个）
"D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh "-test.timeout=1800s"
```

注：裸 `bash` 在本机命中 WSL（`/bin/bash` 不存在，假红），必须用 `D:\work\soft\Git\bin\bash.exe`。
