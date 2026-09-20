# T07 对抗验收报告（编排者执行）

> 执行者：orchestrator。实现者 T07-impl 第三棒完成实现后死于配额（exceed quota limit），
> 实现已全部提交（8464da6/04c628c/1390db4/1e80700/78ed454）。时间：2026-09-19T13:56:25Z

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | 转移表对账 | PASS | table.go 人工逐行审计：D43 #1–#40 全在 + #41/#42（D47）；双目标行（#10/12/15/17/18/19/23/25/32/37/40）guard 分支编码；#34/35/36 用 AnyState 通配且不遮蔽具体行；#30 的 To 复合语义（Conversation→Listening）有注释裁决 |
| 2 | 测试 | PASS | ball ok / statemachine ok（go test 实跑）；table_test.go 314 行 + machine_test.go 195 行 |
| 3 | 视觉证据 | PASS（人工签收挂起） | docs/evidence/s1/ball-states/ 截图 + c21-native-tokens.md 对照表 |
| 4 | D22 | PASS | 依赖 ncruces/go-strftime 为 modernc/sqlite 传递依赖（非新增直接依赖） |
| 5 | 诚实偏差 | PASS | handoff 留 10 条偏差（NoNetwork/WatchdogAlert 无出口行=spec-faithful、托盘通用图标、热键默认值待 39 接线等），全部登记于票据 |

**待人工项**：20 态视觉的主观签收（用户本人）→ 已登记 docs/reports/pending-human-review.md。
**VERDICT: PASS**

## Addendum 裁决（2026-09-20，AC 补裁）

> 背景：上表 5 行为"审计项"而非"逐 AC"表，AC#3/AC#4/AC#5/AC#6 从未被裁（票据复核已记为四个空框）。
> 本节由补裁代理实跑实读后补裁。**四框无一可勾**，理由与缺口逐条列明。
> 环境事实（影响判据，先声明）：本机**仅一台显示器**（`\\.\DISPLAY4` 3440×1440 primary，
> `[System.Windows.Forms.Screen]::AllScreens` 实测单条）；用户 H1 实况签收用的
> `balldebug.exe`（PID 29844，10:19:34 启动，MainWindowTitle="Wisp"）**在本节所有实跑期间常驻桌面**。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#3 交互：点球 Sleeping→Listening / Esc·点击从 Confirming 取消 / 焦点不被抢（需用焦点编辑器验证）/ 透明区点击穿透 | **FAIL** | 唯一被点名的候选 `internal/ball/tokens_test.go:239 TestHitTestAndDPIInjection` 只覆盖**第 4 个分句的一半**：它对纯函数 `HitTest()`（`internal/ball/hit.go:35`）在 4 档 DPI 下断言 HTClient/HTTransparent 分区——这是 SPEC-08 §2"hit-test 返回 HTTRANSPARENT"认可的正确机制，但断言止于几何函数本身，**不经窗口、不产生真实鼠标命中**。其余三项零证据：①"点球 → Sleeping→Listening" 只有转移表行 `internal/statemachine/table.go:48 {D43:4, Sleeping, EvSummon → Listening}`（属 AC#1 已勾的穷举表，非交互路径；且**球侧无任何测试把真实点击转成 EvSummon**）；②"Esc/点击从 Confirming 取消"仅存在于 winlive 门控的 `internal/ball/live_windows_test.go:99-112`，且该段对 Esc 注册失败是 `t.Log` 降级（`:104`），**本机上确实走了降级支**（见下方输出）；③"焦点不被抢（verified with a focused editor）"——**全仓没有一处打开焦点编辑器并断言前台窗口未变的测试**，实现侧 `WS_EX_NOACTIVATE`（`internal/ball/ball_windows.go:168`）与 `MA_NOACTIVATE`（`:505`）仅靠代码审读。<br>`go test ./internal/ball/ -run 'TestHitTestAndDPIInjection' -count=2 -v` |

```
=== RUN   TestHitTestAndDPIInjection
--- PASS: TestHitTestAndDPIInjection (0.00s)
=== RUN   TestHitTestAndDPIInjection
--- PASS: TestHitTestAndDPIInjection (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/ball	0.041s
```
> 该命令 PASS **不构成 AC#3 的证据**：PASS 的是四个分句里唯一被自动化的那半个。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#4 Sleeping 态 CPU≈0（无定时器）——08 采样器可测；**本票处：断言 Sleeping 下无动画定时器句柄存活** | **PARTIAL** | 被点名候选 `internal/ball/tokens_test.go:202 TestAnimationPolicyZeroTimerInSleeping` 断言的是 `AnimationPolicy()` **策略表的返回值**（`:203` `p.Kind != AnimNone \|\| p.PeriodMs != 0` → Fatalf），**不是"存活句柄"这一可观测量**——正是"只对函数返回值断言、未对 observable 断言"的形态。句柄级断言唯一存在于 `internal/ball/live_windows_test.go:95 if b.TimersAlive() { t.Fatal(...) }`，而 ①该文件首行 `//go:build windows && winlive`（`:1`）把它**排除在 `go test ./...` 与 CI 之外**，②本机该测试当前为红（见下方"发现 A"），③`Ball.TimersAlive()`（`internal/ball/ball_windows.go:387`）返回的是内部布尔 `animTimerActive`，该布尔由 `applyStateLocked`（`:303-306`）在 `pSetTimer.Call` 前后置位——即它**镜像**策略表而非独立观测 Win32 定时器。构造链本身经得起审读：`SetTimer` 全包仅 1 处调用（`internal/ball/ball_windows.go:305`，`grep -rn "SetTimer" internal/ball/` 除 `win32_windows.go:34` 的 proc 声明外唯一），且受 `policy.PeriodMs > 0` 守卫，配合 Sleeping 策略 = `{AnimNone,0}` ⇒ Sleeping 拿不到定时器。**但**：AC 字面首句"Sleeping-state CPU≈0"从未在真实 runner 上量得——票 08 报告自记"首次 CI 运行是最终验证点"，且其行 2 只裁了内存（8.7MB）。辅证（不同一类定时器）：状态机侧超时定时器有默认套件内的真断言，见下。<br>`go test ./internal/ball/ -run 'TestAnimationPolicyZeroTimerInSleeping' -count=2 -v` + `go test ./internal/statemachine/ -run 'TestSleeping\|TestTimers\|TestIllegal' -count=2 -v` |

```
=== RUN   TestAnimationPolicyZeroTimerInSleeping
--- PASS: TestAnimationPolicyZeroTimerInSleeping (0.00s)
=== RUN   TestAnimationPolicyZeroTimerInSleeping
--- PASS: TestAnimationPolicyZeroTimerInSleeping (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/ball	0.041s
--- PASS: TestSleepingHasZeroTimers (0.00s)
    --- PASS: TestSleepingHasZeroTimers/boot_into_sleeping (0.00s)
    --- PASS: TestSleepingHasZeroTimers/muted_without_kws (0.00s)
--- PASS: TestTimersArmAndFire (0.13s)
--- PASS: TestIllegalTransitionsRejected (0.02s)
ok  	github.com/CarlosShao/wisp/internal/statemachine	0.320s
```
> `internal/statemachine/machine_test.go:83,93` 断言的是 `Machine.TimersAlive()==0`（真·定时器计数，
> `machine.go:94`），但那是**状态机超时定时器**，AC#4 字面要的是**动画定时器句柄**；不能移花接木。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#5 热键注册 / 配置变更后重注册；静音切换 Muted 态 | **FAIL** | 被点名候选 `internal/ball/hotkey_test.go:51 TestDefaultHotkeys` 与 AC#5 **语义无关**：它只断言 `DefaultHotkeys()` 的四个**默认字符串**（`Cancel=="Esc"`、其余非空），既不注册也不切态。逐项缺口：①"**重注册 on config change**"——生产实现 `func (b *Ball) RebindHotkeys(cfg HotkeyConfig)`（`internal/ball/ball_windows.go:693`）**全仓零调用者**（`grep -rn "RebindHotkeys" --include=*.go .` 仅命中其定义与注释两行；`cmd/`、`internal/config` reload 回调均无），即热配置变更→重注册这条链**尚未接线**，与本表行 5 已登记的偏差"热键默认值待 39 接线"一致；②"**注册**"本身不可断言：唯一注册点 `internal/ball/hotkey_windows.go:165 registerAll` 对失败仅 `slog.Warn`（`:189`）并把失败项从返回映射里丢弃，**没有任何测试断言注册集合非空/四项齐**；本机实测四项**全部注册失败**（见下方 winlive 输出 8 行 WARN），测试照样 PASS；③"**mute 切换 Muted 态**"只有转移表行 `internal/statemachine/table.go:56`（D43 #8 Armed→Muted）/`:62,:64`（D43 #10 Muted→Armed/Sleeping），属 AC#1 已勾范围；球侧 `hkMute` 命中仅 `b.fire(b.opts.Events.OnMuteHotkey)`（`ball_windows.go:576-577`）转成回调，**该回调到 EvMuteKey 的分派无测试**。<br>`go test ./internal/ball/ -run 'TestDefaultHotkeys' -count=2 -v` |

```
=== RUN   TestDefaultHotkeys
--- PASS: TestDefaultHotkeys (0.00s)
=== RUN   TestDefaultHotkeys
--- PASS: TestDefaultHotkeys (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/ball	0.043s
```
> 该 PASS 只证明"默认值字符串没写错"，与 AC#5 要求的注册/重注册/切态无关。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#6 多显示器：拖到第二屏、持久化、恢复；模拟拔出 → 主屏 | **PARTIAL** | 已真正验证的两段：①"**模拟拔出 → 主屏**"是纯函数级真断言且**在默认套件内**——`internal/ball/position_test.go:82 TestResolvePosition/saved_monitor_detached:_primary_default` 用两显示器夹具（`:9 monitorsFixture` DISPLAY1+DISPLAY2）断言 `ResolvePosition(..., "\\.\DISPLAY9", {9999,9999}, true, 70)` 回落到 `DefaultPosition(mons[0])`，另有钳边/未保存/主屏第二位列出等 6 子例；②"**持久化+恢复**"——`internal/ball/position_test.go:17 TestPositionStoreRoundTrip` 断言按 device 键存取、跨重开持久、**per-monitor 隔离**（`:40` 取 DISPLAY1 必须 miss）与损坏文件恢复；窗口级恢复 `internal/ball/live_windows_test.go:125 TestBallLivePositionPersistence` 移窗→persist→Close→重开→比对真实 `GetWindowRect` 左上角，**实测 PASS（见下）**。未验证的一段：AC 字面"**拖到第二屏**"从未被执行——(a) 本机物理上无第二屏（单 `\\.\DISPLAY4`），故真实跨屏拖拽/`WM_DPICHANGED`/按屏恢复均不可测；(b) 更关键的是**端到端两屏路径根本没有接缝**：启动恢复走 `resolveInitial`→`enumMonitors()`（`internal/ball/monitors_windows.go:40`，包级私有、直调 `EnumDisplayMonitors`，无注入点），`live_windows_test.go:125` 只在主屏上按偏移量移窗，从不进入 DISPLAY2；(c) `position_test.go` 的显示器拓扑只喂给纯 `ResolvePosition`，从未喂进 `Ball`。<br>`go test ./internal/ball/ -run 'TestResolvePosition\|TestPositionStoreRoundTrip' -count=2 -v` + `go test -tags winlive ./internal/ball/ -run 'TestBallLivePositionPersistence' -count=2 -v` |

```
--- PASS: TestResolvePosition (0.00s)
    --- PASS: TestResolvePosition/saved_monitor_detached:_primary_default (0.00s)
    --- PASS: TestResolvePosition/saved_visible_on_its_monitor:_kept (0.00s)
    --- PASS: TestResolvePosition/saved_partially_off_right_edge:_clamped (0.00s)
    --- PASS: TestResolvePosition/nothing_saved:_primary_default (0.00s)
    --- PASS: TestResolvePosition/saved_point_on_no_monitor:_that-monitor_default_impossible_->_primary (0.00s)
    --- PASS: TestResolvePosition/primary_detection_when_listed_second (0.00s)
--- PASS: TestPositionStoreRoundTrip (0.00s)
ok  	github.com/CarlosShao/wisp/internal/ball	0.041s
--- PASS: TestBallLivePositionPersistence (0.81s)
--- PASS: TestBallLivePositionPersistence (0.48s)
PASS
ok  	github.com/CarlosShao/wisp/internal/ball	1.365s
```

### 补裁期间发现的两个新问题（原报告未记，均在本票范围内）

**发现 A（阻断 AC#4/AC#5 的可裁定性）：`TestBallLiveLifecycle` 在本机稳定复红，且该测试非 hermetic。**

```
=== RUN   TestBallLiveLifecycle
    live_windows_test.go:104: Esc registration busy (another app holds it); takeover not assertable this run
    live_windows_test.go:119: handles: base=104 peak-suite=376 (gate <600)
    live_windows_test.go:65: ball window still alive after Close
--- FAIL: TestBallLiveLifecycle (2.21s)
=== RUN   TestBallLiveLifecycle
    live_windows_test.go:104: Esc registration busy (another app holds it); takeover not assertable this run
    live_windows_test.go:119: handles: base=376 peak-suite=378 (gate <600)
    live_windows_test.go:65: ball window still alive after Close
--- FAIL: TestBallLiveLifecycle (1.93s)
FAIL	github.com/CarlosShao/wisp/internal/ball	4.211s
```
（命令：`go test -tags winlive ./internal/ball/ -run 'TestBallLiveLifecycle' -count=2 -v`，另 `-count=1` 单跑亦 FAIL）
- **红因须打折**：`live_windows_test.go:65` 用的是 `findBallWindow()`＝`FindWindowW("WispBallWindow", 0)`
  （`:30-34`），该调用**跨进程按类名匹配全桌面**。用户 H1 签收的 `balldebug.exe`（PID 29844，10:19:34 起）
  正是同类名同类窗口，且早于本代理所有实跑存在 ⇒ **"still alive after Close" 不能证明是测试自身窗口的泄漏**。
  这是被测件的缺陷：**live 测试用类名做全局查找，故与任何并发常驻球互扰、不可并行、不可裁定**。
- **但同一次实跑里有一条不受该混淆影响的独立信号**：`handlesOfProcess()` 走
  `GetProcessHandleCount(windows.CurrentProcess())`（`:36-40`），**是进程内计数，别进程无法抬高**；
  `-count=2` 两次观测到 run2 起点 base 由 104 升至 376（两次独立调用分别复现 374/376），
  ⇒ `Ball.Close()`（`ball_windows.go:762-786`，内部 `pDestroyWindow.Call` 返回值未检查）之后
  本进程仍留有约 270 个句柄未回收，与票面约束"release path must bound non-Sleeping handle growth"相悖。
  此项**未定案**：可能是真泄漏，也可能是 D2D/USER 对象待 GC finalizer；需一个"桌面上无其他 Wisp 球"的环境复跑才能归类。

**发现 B（AC#5 的运行时佐证）**：上述 winlive 实跑期间四项默认热键
（`Ctrl+Alt+W`/`Ctrl+Alt+M`/`Esc`/`Ctrl+Alt+P`）**全部注册失败**，8 行 WARN 被测试吞掉、结论仍 PASS。
本机占用属环境事实（`pending-and-issues.md` 已记 `Ctrl+Alt+W` 被占），**不据此判产品缺陷**；
但它具体演示了 AC#5 为何不可勾：注册失败在本仓是**静默降级**，无任何断言会因此变红。

**补裁小结**：AC#3 **FAIL**、AC#4 **PARTIAL**、AC#5 **FAIL**、AC#6 **PARTIAL** —— 四框**全部保持未勾**。
AC#2（20 态视觉人工签收）按指示不属本次补裁范围，原样保留。
