# 314 — 「真窗用例进程内残留」死法的定位与修复 r1（实修腿，2026-11-21）

> **本件是实修件。** 全部读数可复跑；命令与驱动逐条在文内，驱动落 `.scratch/wisp/probes/314/`，
> 逐发日志落 `.scratch/wisp/probes/314/logs/`。
> **代码面锚点**：基线 `fb9da6f73`（＝诊断腿 `314-time-correlation-r1.md` 的提交），修复 commit 见本腿回话首行。
> **机主进程**：本腿一次未碰 `wisp.exe`／浏览器；一次未 kill 任何非本腿自己起的进程；每发载具前查 `wisp.test` 有则等 60 秒。

---

## §0 一句话结论（先看这一行）

> **死法＝真窗用例每轮真建一套 WebView2 环境＋浏览器进程，而没有任何 teardown 关闭 controller（pkg/edge 的
> `Chromium` 无 `Close`、vtbl 字段未导出），环境逐轮净增，同进程下一次环境创建即原生崩溃（RC=139，零 `--- FAIL`）。**
> `count=96` 下**单枚自重复就死**的有 3 枚——恰是族里仅有的 3 枚真窗用例；族序死亡（5 秒复现）与整包死亡
> （327/400/403/422 s）都是它的两种表现。修法（只动 `_test.go` 夹具面）：**controller `Close`（走 wrapper 自己的
> vtable 槽位）＋ Destroy ＋ 等 `WM_DESTROY` 真派发＋ 真窗夹具的 profile 叶子改 swept-never-delete**。
> 修复后自证见 §5。

---

## §1 步骤 1：复现确认（3 发全稳）

命令（`repro3.sh`，逐字可复跑；模式经文件传参，避开 `$` 转义链——诊断腿已踩过的空跑坑）：

```bash
D:/work/soft/Git/bin/bash.exe '.scratch/wisp/probes/314/repro3.sh'
# 内容：scripts/wisp-cli-staged-tests.sh -test.run "$(cat pattern-family.txt)" -test.count=96
# pattern-family.txt = ^(TestPanelHost.*|TestAC4.*|TestResync.*|TestTicket255r8RealWindow.*|TestTicket255RealWindow.*)$
```

| 发 | RC | 墙钟（含载具构建） | 死点（最后一个 `=== RUN`） |
|---|---|---|---|
| r1a | **139** | 13 s | `TestAC4PriorFocusSurvivesARefusedPanelSample` |
| r1b | **139** | 13 s | 同一枚 |
| r1c | **139** | 9 s | 同一枚 |

**3/3 稳**，测试段 2–5 s，与诊断腿 B/B2 两臂逐字同形（零 `--- FAIL`、零 panic、零 timeout）。

---

## §2 步骤 2：单枚扫描（13 枚 × `-test.run '^<name>$' -test.count=96`，一次一发）

命令（`singles.sh`；名字运行时从 staged 二进制 `-test.list` 取，不手写）：

| # | 枚名 | 耗时（墙钟） | rc | 死点/判 |
|---|---|---|---|---|
| 1 | TestAC4CloseWaitsForBothGoroutinesItSpawned | 16 s | 0 | 96/96 全活 |
| 2 | TestAC4CloseNamesTheParkedListenerInsteadOfWaitingOnIt | 3 s | 0 | 96/96 全活 |
| 3 | TestAC4Distinctness133NoTwoLegsWearTheSameRow | 46 s | 0 | 96/96 全活 |
| 4 | TestAC4EveryLegIsNailedOrRuled | 38 s | 0 | 96/96 全活 |
| 5 | TestPanelHostOpensNoListeningSocketL1 | 5 s | 0 | 96/96 全活 |
| **6** | **TestPanelHostRealWindowHopAndLifecycle** | 12 s | **139** | **第 5 轮开局死**（4 轮存活，RUN=5/PASS=4） |
| **7** | **TestAC4FocusReturnToPriorWindowGap33r5** | 6 s | **139** | **第 2 轮开局死**（RUN=2/PASS=1） |
| 8 | TestPanelHostLatencyPercentilesAC2 | 3 s | 0 | 96/96 全活 |
| **9** | **TestAC4PriorFocusSurvivesARefusedPanelSample** | 5 s | **139** | **第 1 轮 t.Cleanup 段死**（`--- PASS` 未印出） |
| 10 | TestResyncOnShowDeliversFreshSnapshotFieldByField35g6 | 112 s | 0 | 96/96 全活（无头 stand-in 窗，不建真 WebView2） |
| 11 | TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo | 6 s | 0 | 96/96 全活 |
| 12 | TestAC452ShippedBuildSaysItHasNoEchoCancellation | 3 s | 0 | 96/96 全活 |
| 13 | TestAC452LostMicWalksD43FromListeningIntoError | 54 s | 0 | 96/96 全活 |

**判读**：死的就是族里仅有的 3 枚真窗用例 ⇒ **「进程内重复建面板」成立**，两枚次序假设不成立为唯一死因
（450 的「单跑 0 死」是 `count=14` 下的读数；`count=96` 下残留累积到第 1–5 轮即崩）。
族序复现（§1）与整包历史死亡（327/400/403/422 s，死点 `TestAC4PriorFocus…`×2 与 `TestResyncOnShow…35g6`×1）
是同一残留机制的两种表现：**死于"前一枚残留环境存在时的下一次环境创建"**。

---

## §3 步骤 3：MARK299 是谁打的

`grep -rn "MARK299" --include=*.go cmd internal` ⇒ **全部命中产码 `cmd/wisp/panel_control_windows.go`**（票 299-r1 的调试标记，共 8 处）：

- `:277` `ctlWndProc` 顶部：**每条到达面板窗窗口过程的消息打一行**（`known=true` 表示该 hwnd 在 `ctlWindows` 注册表里有属主）；
- `:424/:444/:450/:452/:454/:456/:462` `newPanelControl` 建窗各步（before CreateWindowExW / after CreateWindowExW / after ShowWindow / after UpdateWindow / after SetFocus / before Embed / after Embed）。

**宿主与消息循环**：`ctlWndProc` 是 `windows.NewCallback` 的类级回调（类名 `wisp-panel-host`），消息由**拥有该窗口的线程的消息循环**派发——产码里是 `residentPanel.loop`（`panel_resident_windows.go:397`）锁定的 STA 线程跑 `panelControl.Run()`（`GetMessageW` 泵），夹具里是 `hostThreadHarness.runOnThread` 的锁定线程（`pnlPumpOnce` 每 5ms 一拍）。
⇒ 日志里成片 `MARK299 wndproc` 说明**当时有个活的 WebView2 面板窗正被某条 STA 泵着**；死前最后一批 MARK299 全部属于**前一枚用例留下的窗**，而死点那枚**没走到自己的建窗**（日志里它之后零 `MARK299 before CreateWindowExW`）。

---

## §4 机制（读数定案，非推测）

**证据 A——前一枚的窗从未被销毁。** r1a 日志逐窗对照：

| 窗 | 建 | WM_CLOSE(0x10) | WM_DESTROY(0x2) | 判 |
|---|---|---|---|---|
| lifecycle 的 0x9f0c40 | `MARK299 before/after CreateWindowExW` | **有**（113 行） | **有**（124 行） | 拆干净 |
| AC4FocusReturn 的 0x15507cc | 有（128–163 行） | **无** | **无** | **残留**（deferred `mgr.Destroy()` 只 POST WM_CLOSE，`hh.stopHostThread()` 紧随其后，泵的下一次 5ms 拍永远没来） |

**证据 B——环境逐轮净增、浏览器进程永不退出。** single-06 每轮 `AC#1 after Destroy` 的树内 webview 计数：
`0→6`、`6→11`、`11→16`、`16→22`——**每轮净增 5–6 个 msedgewebview2，基线永不回落**；第 5 轮 `=== RUN` 后即死。

**证据 C——死点都落在"下一次环境创建"。** 三份单枚死亡日志的死亡处均在 `=== RUN` 之后、本轮任何
`MARK299 before CreateWindowExW` 之前（bringUp 的 `CreateWebView2Environment` 阶段），或（single-09）
落在 t.Cleanup 段——`rp.stop()` 的 WM_QUIT 路径让 teardown 迟发的 WM_CLOSE 永无泵可派发，窗与浏览器**活着**
跨过 cleanup，随后 `os.RemoveAll(profile)` 在活浏览器脚下删目录，进程死、`--- PASS` 未印。

**证据 D——single-07 全日志零 0x10/0x2**（两轮、两窗，grep 全文无命中）⇒ deferred 拆除竞态 100% 复现。

---

## §5 步骤 4＋5：修法与自证

### §5.1 修了什么（全部 `_test.go`，产码零改动）

| 件 | 改法 |
|---|---|
| `cmd/wisp/panel_host_windows_test.go` | 新增 `panelTestControllerVtbl`/`panelTestController`（镜像 pkg/edge 未导出 vtbl，26 槽）＋ `panelCloseWebView2Controller`（`ICoreWebView2Controller::Close`，槽位 24：被产码行为钉死——槽 23 `NotifyParentWindowPositionChanged` 每条 WM_MOVE 都走、槽 25 `GetCoreWebView2` 每次 Embed 都走、槽 12 `MoveFocus` 是 Focus() 的路，官方接口在 23/25 之间有且仅有 Close）＋ `panelHostTeardownClean`（on-thread Close → Destroy → 同线程拆编辑窗 → 等双窗 `IsWindow==0`）＋ `hostThreadHarness.pumpDead`（泵死了就不再等 task port，t.Fatal 路径不挂 30s） |
| 同上（lifecycle 用例） | body 内 `mgr.Destroy()` → `panelHostTeardownClean`；defer 同步替换 |
| 同上（TestAC4FocusReturnToPriorWindowGap33r5） | 三个 defer → `stopHostThread` + `panelHostTeardownClean(t, hh, editor, mgr)`（编辑窗也改为**同线程**拆——原 cross-thread DestroyWindow 本来就默默失败） |
| `cmd/wisp/panel_assets_299_windows_test.go` | defer 对替换为 `panelHostTeardownClean` |
| `cmd/wisp/panel_resident_windows_test.go` | `startPanelForTest`：`uniqueProfileLeaf+RemoveAll` → `panelProfileLeaf`（swept-never-delete，与族内既有纪律同形）；cleanup `rp.stop` → `stopResidentPanelClean`（on-panel-thread Close → Destroy → 等 WM_DESTROY → 再 stop）；两枚 AC13 reused-thread 用例：leaf 同改、RemoveAll 删、`mgr.Destroy()` 前补 `panelCloseWebView2Controller`；`editorOnPanelThread` 的 cleanup：跨线程 DestroyWindow（静默失败，编辑窗靠线程死亡代拆）→ 经 `rp.post` 在面板线程上真拆＋等 `IsWindow==0` |

**迭代一读（诚实记录）**：只修窗与 profile 时，`v-single-priorfocus`（驻留路径）仍死——第 8 轮开局、死点移进 `newPanelControl` 的 `ShowWindow` 中段（`after CreateWindowExW` 已打、`after ShowWindow` 未打，最后一条消息 0x7 WM_SETFOCUS；日志 `logs/v-single-priorfocus.log`，RUN=8/PASS=7、每轮 WM_CLOSE/WM_DESTROY 各 7 次、线程 7 次 `exited cleanly`）。比修复前（第 1 轮即死）已从"轮轮残留"降到"7 轮后才爆"，说明残留源还剩一个：**该用例在面板线程上建的编辑窗**仍由线程死亡代拆。补上编辑窗同线程真拆（上表最后一行）后，`v2-single-priorfocus` RC=0、96 轮全活。

### §5.2 自证读数（verify.sh 阶梯，逐发）

（本节由实修腿在阶梯跑完后回填——见回话首行与 `.scratch/wisp/probes/314/logs/verify-summary.txt`。）

### §5.3 产码侧建议（⛔ 本腿未改产码，按派单停手写建议）

1. **`residentPanel.stop()/teardown()`（`panel_resident_windows.go:682/:724`）的拆除链有同形缺陷**：`Terminate` 的 WM_QUIT 让 `Run()` 返回后，teardown 才 POST 的 WM_CLOSE 永无泵可派发 ⇒ 生产驻留进程每次 dispose 都留一个活窗＋活 WebView2，下一次 show 的 recreate 路径在注释里已被判"死路"（:662-669），累积环境正是 314 整包死亡在测试里的形状。建议：teardown 在 `Run()` 返回**之前**（WM_QUIT 入队前）先走 Close+Destroy 并等 `WM_DESTROY`，与 `stopResidentPanelClean` 同序。
2. **pkg/edge `Chromium` 不暴露 `Close`**：建议在产码装配层给 dispose 路径补一条 controller Close（本腿的 vtable 槽位论证可直接搬），或在 fork 里给 `ICoreWebView2Controller` 包一层 `Close`。

---

## §6 硬约束核对（本腿自证）

| 约束 | 核对 |
|---|---|
| ⛔ 改非 `_test.go` 产码 | `git show --stat` 仅 `cmd/wisp/**_test.go`＋本件＋`.scratch/wisp/probes/314/**` |
| ⛔ 动 thresholds/golden/allowlist/PLAN/specs | 未动 |
| ⛔ 放宽/删断言、⛔ `-skip` | 零断言改动；零 `-skip`（跳过仅来自载具既有夹具回退，本腿未新增） |
| ⛔ push | 只 commit 不 push |
| ⛔ emoji | 本件零 emoji（d22scan GREEN） |
| ⛔ kill 非本腿进程 | 未 kill；载具一次一发；单发最长（整包）<25 分钟上限 |
| 载具跑前查 `wisp.test` | 每发前查（驱动内置），有则等 60 秒 |

## §7 撤销口令

「314 实修腿撤」＝ revert 本腿 commit（测试文件回到 `fb9da6f73` 形状）；`.scratch/wisp/probes/314/**` 只建不删，留档。
