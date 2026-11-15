# 46c · 审计后三件补齐（写腿实现件）

- 腿：agent-76（写腿），2026-11-15 22:1x
- 依据：`docs/evidence/s1/46b-audit-20261115.md` 末三节（必须回退 1 条／可以留／另记 2 条）＋ `.scratch/wisp/issues/46-s6-acceptance.md` 审计收件节（六项成立 4／不成立 2，两条裁定）
- 三件全部来自审计，零自选活。

## 三件的落点（file:line，盘上现读）

| # | 审计条目 | 产码落点 |
|---|---|---|
| 1 | `DiskHeldStats()` 无生产读者（审计 §3「不成立」，裁定①接生产读者不删） | 读者：`cmd/wisp/diskgate_diag_46_windows.go:43` `diskHeldStatsLine46`；打印点：`cmd/wisp/resident_windows.go:895`（resident 启动行，紧跟收口闸行，**无条件打印**） |
| 2 | 启动行 `line()` 两态是假话（审计 §1 ★一条／另记 #1） | 三态：`cmd/wisp/resident_settle_gate_windows.go:283` `line()`；挂载标记：`resident_voice_engines_windows.go:632`（挂上时 `markAttached`）＋ `resident_settle_gate_windows.go` 的 `attachSettleGate` 拒绝分支（带 buildResidentVoice verdict 具名原因 `markUnattached`）；装配点：`resident_watchdog_windows.go:86` |
| 3 | 跳闸裸 `ErrDiskSpaceLow` 不分类（审计 §2 ★／另记 #2） | `internal/observe/guards_45.go:254` `Wrap(ClassResource, fmt.Errorf("%w: …", ErrDiskSpaceLow, …), "disk-full guard refused a gated write")` |

### 件 1 的选点理由（为什么不是 `wisp doctor`/`wisp diagnostics`）
两者都是**独立 CLI 进程**：读不到常驻进程内存里的 `LogPipeline`/guard。诊断包的 `DiskStatus` 缝是给装配时持有 `*DiskGuard` 的调用方的；CLI 腿不装管线（`main.go:173` 自己的注释：diagnostics 腿「reaches no installLogSink」）。持有门控写入器的唯一进程是 resident（`installLogSink`，resident_windows.go:173）⇒ 生产读者落在它的启动行：每次启动读 `DiskHeldStats()` 打一行。

**四句原文（三态＋无管线态）**：
- 正常（未被扣住，always-one-line 契约）：`日志磁盘门控读数：未被扣住（门控健康：扣住 0 条／丢弃 0 条／回放 0 条）`
- 异常（被扣住）：`日志磁盘门控读数：日志被扣住 N 条／丢弃 N 条／回放 N 条`
- 未装门控（`obs.Disk = nil` ⇒ `DiskHeldStats` ok=false）：`日志磁盘门控读数：未装门控（安装时未装配磁盘守卫，日志写入不经任何扣留/回放门）`
- 无管线（installLogSink 拒绝分支，nil 安全）：`日志磁盘门控读数：无（持久日志管线未装配，本次启动没有可读的门控计数）`

### 件 2 的三态原文（`line()` 现读）
1. **未装配**（g == nil）：`收口闸：未装配（本进程没有 watchdog 循环，D32 settle 行在这次启动里没有验证者）`
2. **已装配但未挂载**（新）：`收口闸：已装配但未挂载（<具名原因>）：dispose 级联第 3 步摸不到这枚闸，D32 settle 行在这次启动里不会被触发验证`
   具名原因取自 boot 自身：`语音腿未构造无常驻（<buildResidentVoice 的 verdict>）`／`本进程没有语音腿（residentVoice 为 nil）`——首装未跑 `wisp models ensure` 的常态上打的是**真话**。
3. **装配**（已挂载）：`收口闸：装配（触发＝会话 dispose 级联第 3 步 FreeOSMemory；目标态＝D43 机现态，非采样态回落 Sleeping；窗口＝watchdog 的 SettleWithin；证据＝watchdog settle verification 日志行＋三次连续失败升 D43#35）`

### 件 3 的读法
`Wrap` 的 `Error.Unwrap` 暴露被包的 `fmt.Errorf` ⇒ `errors.Is(err, ErrDiskSpaceLow)` 链路不变（哨兵身份未丢）；`ClassOf(err)` 现在答 `(ClassResource, true)`（此前 `(–, false)`）。

## 三枚牙（行为牙，A1880：断言对象全是运行读数，零源码字符串）

| 测试名 | 位置 | 钉什么 |
|---|---|---|
| `Test46DiskHeldStatsLineChangesWhenTheGateIsAbsent` | `cmd/wisp/diskgate_diag_46_windows_test.go` | 走生产装配 `installLogSink`（真 guard 缝 `diskGuardForDataDir`）四态读行：健康⇒未被扣住；`obs.Disk=nil`（经装配缝实现）⇒行**变**成未装门控；跳闸⇒扣住/丢弃/回放读数；无管线⇒「无」句 |
| `Test46BootRowSaysAssembledButUnmountedWithoutAResidency` | 同上 | 生产挂载半边：无 residency ⇒ attach 返回 false 且 `line()` 须是「已装配但未挂载」＋含 boot verdict；真 residency ⇒ 装配句回归；nil gate ⇒ 未装配句不变 |
| `Test46DiskGuardTripIsClassifiedResource` | `internal/observe/guards_45_diskclass_46_test.go` | 真闸真探针跳闸：`ClassOf(refusal)==ClassResource` **且** `errors.Is(refusal, ErrDiskSpaceLow)` 两读数同发成立 |

## 拆掉时各自红句（三发 -overlay 突变，替换件只落 %TEMP%，⛔ 未落仓；载具照走）

| 拆掉什么 | 红句逐字 | 耗时 |
|---|---|---|
| 突变1：`logsink.go` `obs.Disk = guard` → `obs.Disk = nil; _ = guard`（审计点名的同一突变） | `46: the healthy gated pipeline's row = "日志磁盘门控读数：未装门控（安装时未装配磁盘守卫，日志写入不经任何扣留/回放门）", want the 未被扣住 sentence: a line that only speaks under pressure teaches no baseline and the gate's existence is unprovable from the boot output` | FAIL 0.02s |
| 突变2：删 `resident_voice_engines_windows.go` 的 `g.markAttached()`（挂载标记） | `46: the mounted gate's row = "收口闸：已装配但未挂载（语音腿未构造无常驻（语音链未构造：[voice] enabled=false（配置来源 config.toml）））：dispose 级联第 3 步摸不到这枚闸，D32 settle 行在这次启动里不会被触发验证", want the 装配 sentence: a real residency must not be described as unmounted or the operator stops trusting the row` | FAIL 1.00s |
| 突变3：`guards_45.go` 跳闸那支改回裸 `fmt.Errorf` | `46: the tripped refusal reads as class "" (ok=false), want "resource": the disk-full trip carries no D37 class at the boundary` | FAIL 0.037s |

## 读数（三把尺）

**尺1 载具点名（一次一发，全部经 `scripts/wisp-cli-staged-tests.sh`，跑前 `Get-Process wisp.test`＝无）**：
46 全家＋384 装配 **8 PASS / 0 FAIL**（`Test46DiskHeldStatsLine…` 0.04s／`…Unmounted…` 0.00s／`…HoldsAndReplays…` 0.02s／`…DisposeCascade` 0.13s／`…JudgesTheState…` 0.13s／`…WakeCancels…` 0.00s／`…LineNamesItsTrigger` 0.00s／`Test384StartResidentWatchdogRegistersRealLoop` 0.00s）。

**尺2 observe 全包**：`go test ./internal/observe -count=1` ⇒ `ok 4.890s`（含新牙）。

**尺3 整包（A1882 口径）**：
- `RUN` 计数 **502**，`--- FAIL` 计数 **0**，`--- SKIP` 计数 **0**；
- **终点锚命中**：`TestZZ315A1RootLoopFirstBridgeCall` 在本波之前的独立载具发里 **PASS 0.30s**（同发 `TestZZ315A1SiblingChildFirstBridgeCall` 0.31s／`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` 0.22s）；
- 最后一条 `=== RUN` ＝ **`Test145ArgsJSONComesFromTheBookedForensicsRow`**（整包在该枚**中途**二进制无声断掉：无 `--- FAIL`、无 panic、无 SIGSEGV 行；该枚与本波零交叠——票 145 panel 面的 forensics 行测试）；
- **本波新测试位置＝截断点之前**：`Test46DiskHeldStatsLineChangesWhenTheGateIsAbsent`（22:05:50 PASS 0.02s）／`Test46BootRowSaysAssembledButUnmountedWithoutAResidency`（PASS 0.00s）／`Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink`（22:06:12 PASS 0.02s）全部落在截断点（22:08:4x 起 145 族）之前且 PASS；
- 截断枚单跑复量（载具形）：`Test145ArgsJSONComesFromTheBookedForensicsRow` **PASS 0.16s** ⇒ 中途断是**包级环境事故**（此前 22:06–22:08 连续 7 枚真窗口 panel 用例刚跑完），不是断言红；按 A1886「cpu 不涨＝在等」与诚实记账原则，整包判读记为「502 RUN / 0 断言红 / 1 枚中途截断（单跑复量绿）」，**不宣称整包全绿**。

**七禁门**：`cd tools/d22scan && go run . -root <repo>` ⇒ `clean - no D22 ban violations`（internal 405＋cmd 87 全查）。

## 硬约束对账
- ⛔ 未动 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43／`frontend/**`／`design/**`（`git status` 现核：本腿只写下列 pathspec 内的件）；
- ⛔ 未放宽任何断言（新牙全部是新增读数；唯一既有测试改动是 `Test46SettleGateLineNamesItsTrigger` 为到达「装配」态改走生产挂载半边——断言集合未减，新增挂载读数一行）；
- ⛔ 未 push；新协程零枚（本波三件全部无新 goroutine；`markAttached/markUnattached` 是锁内两次赋值）；
- 串行纪律：`go test ./cmd/wisp` 全部经载具一次一发，跑前现查 `wisp.test` 进程；两个 sh 门用 `D:\work\soft\Git\bin\bash.exe`。

## 撤销口令
- 件 1：`git revert` 本提交即可（boot 行一行＋一个 helper 文件）；
- 件 2：删 `line()` 的中态分支与 `markAttached/markUnattached` 即复原两态；
- 件 3：`guards_45.go:254` 改回裸 `fmt.Errorf` 即复原（审计件的撤销口令原文）。
