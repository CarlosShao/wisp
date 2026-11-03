# 45-r1 实现腿｜internal/observe 诊断守卫 + 诊断包隐私预览（票 45，S6／D42#5#8／D35 artifacts quota）

写面：`internal/observe/**`（+ 本证据件 + 派单存档）。⛔cmd/wisp／⛔internal/agent·memory／⛔internal/tools／⛔internal/plugin／⛔thresholds·golden —— **全部零字节**（`git status --porcelain -- internal/observe` 起手＝空，收尾＝只本腿三枚件）。
HEAD 起手：`d0cf6a7c7`；索引：无 staged 内容（`git diff --cached --stat` 空），工作树只有 `.scratch/**` 未跟踪噪声。

## 起手三核（读数）

| 核 | 尺 | 读数 |
|---|---|---|
| HEAD 最新 | `git log --oneline -3` | `d0cf6a7c7 编排者A1442：44-r5交付…；派45-r1(observe解押)` |
| observe 无他腿脏 | `git status --porcelain -- internal/observe` | **空**（44 已结，24-v2 只读 internal/tools ⇒ 同包冲突已解） |
| 票面所述未落地 | `grep -rniE "disk.?full\|FreeSpace\|FocusAssist\|PreviewDiagnostics" *.go` + 读 `internal/observe/diagnostics.go` | 诊断包 `BuildDiagnosticsBundle` **已在**（08 的 API 桩，注释自己写着"ticket 45 completes the UX wiring"）；**守卫一件不在**：无 free-space 闸、无 WAL 再核、无 artifacts LRU、无预览／逐条 opt-in ⇒ 非"已落地"，不停手 |
| 引文出处 | 票面 `:17-29` 逐行 | `:20` "preview of what's included BEFORE export (§14.4)"／`:21-22` "pre-write free-space check (200MB threshold)…stop L3/log/artifacts writes + user warning, main response path UNAFFECTED"／`:22-23` "WAL kept bounded (startup TRUNCATE + autocheckpoint 1000 from 04, runtime re-check here)"／`:27` "Artifacts quota enforcement point here (LRU at 500MB…)"／`:29` "Log rollover under disk pressure never loses the last error line" —— 全部对得上；watchdog 快照的 attach 形状出处 `internal/watchdog/watchdog.go:551-552` 亦在 |

## 逐格 AC

| AC（票面 `:35-41`） | 本腿落点 | 尺（可逐字重跑） | 读数 |
|---|---|---|---|
| **AC#1** 包内容＝预览所述；密钥／转写／音频默认不在；opt-in 逐条＋逐条预览 | `diagnostics.go`：`PreviewDiagnosticsBundle`＋`BundlePreview`／`PreviewCategory`／`PreviewItem`＋`BundleOptIn`＋`checkBundleItemName`＋builder 第 6/7 段（guard 读数、逐条转写／快照、音频声明） | `go test ./internal/observe -run TestAC1BundleContents -count=1 -v`；比的是**条目名清单**，不是文案 | **1 PASS／0 FAIL（0.04s）**：默认包 7 枚条目与 `prev.Entries` 逐字同集；无 `transcripts/`·`watchdog/`·`audio*` 条目；超 `MaxLogBytes` 的日志两路径**同跳**（预览带 Reason，zip 无此件）；opt-in 一枚 ⇒ 只多那一枚，正文逐字＝`Redactor.String(原文)`（同一个不可关脱敏器，无第二套格式）；`"" .. . ../secrets a/b a\b .hidden` 七枚畸形名 ⇒ 预览与构造**都**红；预览失败面＝构造失败面（不可读日志目录两红） |
| **AC#2** disk-full fixture：闸在阈值跳、L3+日志停、活任务照跑、腾出空间即恢复 | `guards_45.go`：`DiskGuard`／`SpaceProbe`（注入＝票面"small virtual disk"）／`DiskGuardThresholdDefault=200<<20`／`WriteClass{Essential,Log,Artifacts,L3}`／`ErrDiskSpaceLow`／`DiskStatus` | `go test ./internal/observe -run 'TestAC2DiskGuard' -count=1 -v` | **3 PASS／0 FAIL**：199MiB ⇒ log/artifacts/l3 三枚全红且 `errors.Is(ErrDiskSpaceLow)`；**essential 恒 nil（`main response path UNAFFECTED` 唯一落点，单列一枚断言）**；警告文案含"live task is unaffected"＋`TripAt` 已落盘形状；`p.set(3GiB)` ⇒ 下一次判定即恢复（不重启）；探针失败＝gated 拒写**且**不冒充 ErrDiskSpaceLow、Status 记 `ProbeError` 而非 0 字节；无探针／无 volume＝构造期红（fail-closed）；缓存只在健康期生效（窗口内 20 次判定 0 次重探，过窗必重探，tripped 期必重探） |
| **AC#2 附**":29 日志轮转不丢最后一条错误行"（flusher 协同） | `guards_45.go`：`diskGatedWriter`（`NewDiskGatedWriter`，`heldLogTailMaxBytes=bufferedLogBytes`）——拒收的字节**上尾不丢**，恢复后先重放再写新行；尾满丢**最旧**并计数 | `go test ./internal/observe -run 'TestGatedLogs' -count=1 -v` | **2 PASS／0 FAIL**：压力期 ERROR 行不入文件、`HeldStats` held>0／lost=0；恢复后文件同时含四行且**重放顺序在触发写之前**（时间线不断）；4× 上限风暴 ⇒ held ≤ 64KiB、lost>0（成本是数字不是沉默）；"读了有空、写下 ENOSPC" ⇒ 字节回尾、`Flush` 再放、尾清零 |
| **AC#3** WAL 持续写下有界；按需 checkpoint | `guards_45.go`：`CheckWalFile`＋`WalAutoCheckpointPages=1000`（04 的启动值）＋`WalMaxBytesDefault`=2×autocheckpoint＋动作名 `none`／`wal_checkpoint(AUTO)`／`wal_checkpoint(TRUNCATE)` | `go test ./internal/observe -run TestAC3Wal -count=1 -v` | **1 PASS／0 FAIL**：16 次持续写逐个读数 ⇒ 首次越界那枚 `OverBound=true` 且动作点名 TRUNCATE，半界处先报 AUTO 折叠，末读仍越界（不衰减）；无 -WAL＝健康；**父目录不在＝`ReadingFail` 记账**（Windows 两者都回"not exist"，故本腿问父目录，不把不可读算成 0 字节） |
| **AC#5** artifacts LRU 超额驱逐；驱逐有日志 | `guards_45.go`：`EnforceArtifactsQuota`＋`ArtifactsQuotaDefault=500<<20`（D35）＋`ArtifactsEvictionReport` | `go test ./internal/observe -run TestAC5Artifacts -count=1 -v` | **1 PASS／0 FAIL**：300B/120B  ⇒ 逐枚 `slog.Warn("observe: artifacts quota eviction"…)`，驱逐序 `a-old-1.bin,b-old-2.bin`（mtime 最旧先，理由写在件首：NTFS atime 不可靠 ⇒ 只能用"最后写入"）；件真没了／最新件活／`BytesFreed=200`／`TotalAfter ≤ quota`；不超额时零驱逐；`quota<=0`＝取 500MB 而非关闸；目录**绕过**不记 Refused；空 dir＝配置红 |
| **AC#4** Focus-Assist 抑制时 badge+任务list 仍更新 | **不在本腿**：这条的写面是通知／球／面板侧（`⛔cmd/wisp`、`⛔internal/plugin`），`internal/observe` 里没有可落的守卫件；`grep -rniE "FocusAssist" internal cmd`＝**0 命中**（best-effort 状态读与"无条件双通道"降级都还没人落） | 同上 grep | **未落地 ⇒ 上报编排者另派一枚写腿**（票面 `:7-8` 的 B 槽）。本腿**不硬做**：跨面会撞"每波只许一枚写腿"的共享编译面规矩 |

## 反向判据＋还原核（尺与读数）

尺：把 `guards_45.go` 里 `low := free < g.threshold` 替成 `low := false`（摘掉守卫的跳闸判定），只跑点名用例，再逐字节还原。

- 摘后：`TestAC2DiskGuardTripsAtThresholdAndRecovering`／`TestAC2DiskGuardCachesOnlyWhileHealthy`／`TestGatedLogsHoldTheLastLineAndReplayItAfterRecovery`（+ 其余断言）**4 枚用例红、15 条断言句**，头三句逐字：`guard did not stop write class "log" at 199MiB free (threshold is 200MiB)`、`the guarded write reached the file while the guard was tripped`、`after the recheck window the guard must look again and trip`。
- 还原：`RESTORED-OK`（与 `D:\tmp\wisp\20261103\45r1\guards_45.go.orig` 逐字节同）⇒ 同批用例 `ok 0.058s`。
- ⇒ 这四处断言**有牙齿**：绿不是断言自己造的洞。

## 全量读数与别的邻居面

- `go test ./internal/observe -count=1` ⇒ **ok 4.258s**（本包既有的 130／132／136／44 等件**无一转红**；`diagnostics_test.go` 三枚旧用例仍绿——本腿只往 manifest 里加行，未改既有断言一字）。
- `go build ./cmd/wisp ./internal/watchdog` ⇒ exit 0（`BundleOptions` 加字段是**加法式**改动，旧调用零改动）。
- `go vet ./internal/observe` ⇒ 静默。
- `gofumpt -w internal/observe/guards_45.go internal/observe/diagnostics.go internal/observe/guards_45_test.go` ⇒ 已跑（`/d/work/base/gopath/bin/gofumpt.exe`，第一趟只重排了测试件一枚，第二趟重排新增件）。
- SLO 阈值／golden／`thresholds.go`／`go.mod`／`go.sum`：**零字节**。本腿的 200MB／500MB／8MiB 是**写入准入数**（票面自己给的），不是 D32 的门，也不在任何门文件里。
- 禁形：无裸 `go func(`、无 `filepath.Clean|Abs` 决策（只用 `Join`＋基名规则）、无墙钟差实现超时（`now` 读单调钟，落盘串走 `WallTimestampUTC`）、零 emoji。

## 本腿留下的两处（诚实记账，不当完成报）

1. **装配缝未接**：闸与日志管的接缝是导出的 `NewDiskGatedWriter(w io.WriteCloser, g *DiskGuard) *diskGatedWriter`，但 `LogConfig` 还没有 `Disk *DiskGuard` 字段、`InitLog` 还没把 sink 换上——**cmd 侧装配属另一枚写腿**（与 44-r5"配额闸装进 cmd/wisp 派发链"同一分法）。装配那枚要做的：`winsec`/`x/sys` 的 `GetDiskFreeSpaceEx` 作探针注进 `NewDiskGuard`、`InitLog` 里 `w` 经 `NewDiskGatedWriter` 包一层再交给 `slog.NewJSONHandler`，并把 `FlushNow`／`Close` 改走包外层的 `Flush`。
2. **AC#4 未落地**（见上表），票面 `:24-26` 的"预期在 Win11 读不到 ⇒ 无条件双通道＋PRECHECK 诚实记一句"整条待另一写面。

改动文件：`internal/observe/guards_45.go`（新）、`internal/observe/guards_45_test.go`（新）、`internal/observe/diagnostics.go`（改）、`docs/evidence/s1/45r1-impl.md`（本件）、`.scratch/wisp/probes/orch/20261103/txt/45r1-dispatch.txt`（派单存档）。
