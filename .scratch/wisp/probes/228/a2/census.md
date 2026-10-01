# 票 228 · 普查腿 `228-a2`：为"GUI 常驻腿的停机触发"装一枚会响的仪器——射程普查

只读普查。零产码、零脚本、零测试、零构建（另一枚验收腿 `33-v1` 正在 `cmd/wisp`／`internal/panel` 跑测试与突变，
本腿**未跑任何** `go build`／`go vet`／`go test`／`scripts/build.ps1`；全部读数来自 Read／Grep／Glob／`wc`／`git`）。

**本件要答的那一个问题**：在今天的仓里，"托盘那一次点击导致了十步有序退出"这一维，能被哪几种仪器形状看见？
各要动什么、今天有什么会因此红、有什么天生看不见。

⛔ 本件不提"新建框架"，所有形状都落在仓里现成可照的件上；⛔ 本件不推荐修法（§⑦ 只摆甲／乙／不做）。
⛔ 未读未引 `frontend/**` 与 `design/**`；未动 `docs/PLAN.md`／`docs/specs/**`／`docs/BUILD.md`／`docs/SLO.md`／
`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件／任何 `*.go`。

---

## 起手锚点

同发取的（一条命令里连续跑；逐字读数见 §⑤ R1）：

- 现读时刻：`2026-09-30 11:06:22 +0800`（`date "+%Y-%m-%d %H:%M:%S %z"` 的原样打印）
- `git log -1 --format=%H` → `0c8b9fddd2c3da647a76f0757fd466ef6767e6c1`
- `git rev-parse --abbrev-ref HEAD` → `dev`
- `.scratch/wisp/probes/228/` 当时只有 `a1/ r1/ v1/`；`a2/` 由本腿新建

> 编排者给的锚点是"AC#11 今天 11:00 由我落盘，提交 `3398e6f1`"。我看到的 HEAD 是 `0c8b9fddd2`，与之不同（见 §⑥ 条 E1）。
> **本件所有 `file:line` 都是相对 `0c8b9fddd2` 的现读行号**，尺命令全部附在 §⑤，编排者可逐条复跑比对。

---

## ① A 托盘点击的可注入面

### A-0 链路全图（每条都是现读行号）

| 段 | 词面锚点 | `file:line` |
|---|---|---|
| 托盘注册（回调消息号与 uID） | `uCallbackMessage: wmAppTray`、`uID: trayUID`、`uFlags: nifMessage \| nifIcon \| nifTip` | `internal/ball/tray_windows.go:37`、`:35`、`:36`（`hWnd: hwnd` 在 `:34`）；`trayUID = 0x5701` 在 `:18`；`nifMessage = 1` 在 `internal/ball/win32_windows.go:265` |
| 消息号本体 | `wmAppTray = wmApp + 0x202 // tray icon callback`；同族 `wmAppTask = wmApp + 0x201` | `internal/ball/win32_windows.go:103`、`:102` |
| 右键分流 | `case wmAppTray:` → `switch lParam & 0xFFFF` → `case wmRButtonUp:` → `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)` | `internal/ball/ball_windows.go:669`、`:670`、`:673`、`:674` |
| 选择值→回调 | `case menuExit:` → `b.fire(b.opts.Events.OnTrayExit)` | `internal/ball/ball_windows.go:682`、`:683` |
| 四枚菜单 id | `menuOpenPanel = 1 / menuMute = 2 / menuPauseWake = 3 / menuExit = 4`（**均非导出**） | `internal/ball/tray_windows.go:20-23` |
| 命令**不回 `WM_COMMAND`** | 注释原文：`wParam of WM_COMMAND after TPM_RETURNCMD is NOT used; we use TrackPopupMenu's return value inline in the wndproc instead` | `internal/ball/tray_windows.go:15-16` |
| 选择值的真来源 | `sel, _, _ := pTrackPopupMenu.Call(menu, tpmRightButton\|tpmReturnCMD\|tpmNoActivate, pt.x, pt.y, 0, hwnd, 0)` → `return uint32(sel)` | `internal/ball/tray_windows.go:101-102`、`:106`（`showMenu` 签名 `:72`；`tpmReturnCMD = 0x0100` 在 `win32_windows.go:129`） |
| 回调执行体 | `func (b *Ball) fire(fn func())`（**非导出**） | `internal/ball/ball_windows.go:728` |
| STA 线程投递 | `func (s *staThread) PostTask(fn func())`（**方法名导出、接收者类型非导出**）→ `pPostMessageW.Call(hwnd, wmAppTask, id, 0)` | `internal/ball/sta_windows.go:127`、`:142`；`type staThread struct` 在 `:24`；`pPostMessageW` 在 `win32_windows.go:33` |
| 宿主注入点 | `Events: ball.Events{... OnTrayExit: recordTrayExit, ...}`（十枚键） | `cmd/wisp/resident_ball_windows.go:122-133`，`OnTrayExit: recordTrayExit` 在 `:131` |
| 今天那枚"零执行者" | `func recordTrayExit()`：`slog.Warn("tray exit requested", "outcome", "ignored", ...)` | `cmd/wisp/resident_ball_windows.go:253`、`:255`（console 那句在 `:256`） |

### A-1 关键判断：**不依赖真鼠标能不能把一次托盘命令送达？**

分四支答，每支给凭据：

1. **`PostMessage`/`SendMessage` 一枚 `WM_COMMAND` ＋ 那枚 id 到那枚窗口 → 不能，这条路在码里不存在。**
   凭据：`tray_windows.go:15-16` 那句注释**明写不用 `WM_COMMAND`**；`menuExit` 的值只从 `TrackPopupMenu` 的**同步返回值**拿到
   （`tray_windows.go:101-106`），`wndProc` 里没有任何 `case wmCommand`（分流只在 `ball_windows.go:669-686` 的 `wmAppTray` 支）。
   ⇒ 这条是任务描述里的一个**假设，本普查把它证伪了**。

2. **同包（`package ball`）＋ tag `winlive`：能把一条真消息送进真 wndproc——形状今天就有。**
   凭据（现读）：`internal/ball/live_guard_windows_test.go:251-258`
   ```
   func sendMouse(b *Ball, msg, wParam, lParam uintptr) uintptr {
       done := make(chan uintptr, 1)
       b.sta.PostTask(func() {
           r, _, _ := pSendMessageW.Call(uintptr(b.hwnd), msg, wParam, lParam)
           done <- r
       })
       return <-done
   }
   ```
   以及它的别名 `sendMessage`（`:261-263`）、`pSendMessageW` 的载入（`live_guard_windows_test.go:36`）。
   PostMessage 变体同样现成：`live_windows_test.go:553-562` 的 `postMsg`（`pPostMessageW.Call(uintptr(b.hwnd), uintptr(msg), 0, lParam)`）。
   ⇒ 于是 `sendMessage(b, wmAppTray, 0, wmRButtonUp)` **这一行今天就能在同包 winlive 台件里写**，它会真进 `ball_windows.go:669-686`。
   **但**：它会停在 `showMenu` 里那句阻塞的 `TrackPopupMenu`（`tray_windows.go:101`），要拿到 `sel == 4` 就得让那枚弹出菜单真被选中。
   仓里**已有**注入真输入的先例：`live_guard_windows_test.go:274-290 injectHotkey`（`pKeybdEvent`）与
   `interaction_live_test.go:35 procMouseEvent = modUser32.NewProc("mouse_event")`。
   ⇒ **可达性与确定性之间的取舍不是我这一腿能定的**，落成 §⑦ F1 甲/乙。

3. **同包、跳过菜单那一步：完全确定性，但看不见"分流"。**
   `b.sta.PostTask(func() { b.fire(b.opts.Events.OnTrayExit) })`——同形用法已在 `internal/ball/hotkey_live_test.go:327-334`
   （`h := b.sta; h.PostTask(func() { bump(&mu, hits, "mute-cb") })`）。
   ⇒ 它**绕过 `ball_windows.go:670-675` 整段**（不读 `lParam`、不调 `showMenu`、不看 `sel`），
   所以它证明的是"回调链通"，不是"托盘那一次点击通"。这是**假绿风险最高的形状**，见 §② 结论。

4. **跨进程（另一个真 `wisp.exe` 的球窗）：句柄今天就能拿到，且已在用。**
   - 导出面凭据：`internal/ball/ball_windows.go:1003 func (b *Ball) DebugHWND() windows.HWND { return b.hwnd }`（进程内可达，
     `interaction_live_test.go:81 hwnd := b.DebugHWND()` 是现成用法）。
   - **跨进程凭据（本普查挖出来的一枚，不在任务描述里）**：`cmd/wisp/resident_ball_live_228_windows_test.go:60-79`
     `findBallWindows228` 用 `FindWindowExW`＋`GetWindowThreadProcessId` 枚举类名 `WispBallWindow`（`:55`）的**顶层窗口**并按 pid 认领（`:129`）。
     ⇒ 顺带把"那枚窗口是不是 message-only"从推断升为现量：`FindWindowExW(0, prev, cls, 0)` 只枚举顶层窗口，
     且 `HWND_MESSAGE` 在 `internal/ball` **零命中**（R8）；`live_windows_test.go:512-515` 还对它调 `SetWindowPos`/`GetWindowRect`。
     ⇒ **它是真顶层窗口，不是 message-only**。
   - 缺的那一点：该文件的 proc 表只有 `FindWindowExW / IsWindow / GetWindowThreadProcessId`（`:45-48`），**没有 `PostMessageW`**；
     要跨进程投 `wmAppTray` 需在 `package main` 的台件里再载一枚 user32 proc（**同库、非新依赖**，形如 `:45` 那三行）。
     且投过去之后仍撞第 2 支那个 `TrackPopupMenu` 阻塞问题，只是这次阻塞在**被测进程**的 STA 线程里。

5. **`cmd/wisp` 同包直接调执行体：可达，且不需要任何窗口。**
   `recordTrayExit` 是 `package main` 的包级零参函数（`resident_ball_windows.go:253`）⇒ 同包台件可直接调。
   但它同样**不经过托盘分流**，与第 3 支同性质（更弱：连 `b.fire` 的 recover/线程语义都不走）。

### A-2 三枚导出面事实（决定"外部包能做什么"）
- `Events` 的九枚回调字段名**全导出**（`ball_windows.go:50-59`，含 `OnTrayExit` 在 `:58`）⇒ 外部包能注入，不能读回。
- `wndProc`（`:549`）、`fire`（`:728`）、`uirun`（`:741`）、`showMenu`（`tray_windows.go:72`）**全非导出**
  ⇒ 跨包无法喂消息；`internal/ball` 的同包台件是唯一走径。
- `Ball.Close()` 导出（`:935`）、`SetTrayChecks`（`:919`）、`SetTrayTip`（`:927`）导出——它们动的是托盘**显示**，不动命令分流。

---

## ② B 三族（实为四族）台件逐枚可否照抄

### B-1 族①：真机进程族（起真 `wisp.exe`，读 console＋落盘日志＋退出码）

| 成员 | `file:line` | setup | 触发 | 断言面 | 可否照抄 |
|---|---|---|---|---|---|
| `TestAC228ResidentLegReportsAndBooksItsBall` | `cmd/wisp/resident_ball_228_windows_test.go:53` | `buildWispForTest`（`:54`）＋`bootResidentLeg`（`:56`）＋`t.TempDir()` | `leg.breakToLoop()`（`:89`）＝`GenerateConsoleCtrlEvent(CTRL_BREAK_EVENT, pid)`（定义在 `resident_sink_nail_127_windows_test.go:252-254`） | console 里 `ballUpClaim`/`ballAbsentClaim` **恰一次**（`:66-81`）＋退出码 `leg.waitErr`（`:99-102`）＋`readResidentSink`（`:104`）里 created/refused/stopped/install 的**索引关系**（`:117-167`） | **形状可照、触发器不可照** |
| `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` | 同文件 `:184` | 同上 | 同上，但**投递时刻**由 `observe.CountLogFiles(sinkDir) > 0` 锚定（`:190-196`） | 退出码＝0 且点名 `0xc000013a`（`:205-209`）＋ install/shutdown 记录存在（`:211-220`） | 同上（这一枚的**时刻锚**是空窗那一坑唯一现成的可照件，见 §③） |
| `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` | `cmd/wisp/resident_sink_nail_127_windows_test.go:565` | 同一套 helper（`:566-569`） | `breakToLoop`（`:581`） | record 0 必须是 `residentEarlyResolverMsg`（`:608`）、`installIdx == 1`（`:619`）、trail 非空（`:629`）、**最后一条记录必须是 shutdown 记录**（`:633-635`）、console 有 `residentShutdownLine`（`:637`） | **形状可照**；它已经是"审计顺序"最接近真机的一枚，但**步号序列只进 `t.Logf`**（`:640-641`，见 §④ D-2） |
| `TestAC246ShippedResidentProcessOwnsItsCancelStep` | `cmd/wisp/resident_approval_246_windows_test.go:184` | 同一套 helper（`:185-188`） | （该用例只读 boot 报告，不断开退出） | console 同时含 ball 姿态句 **和** `cancelStepRosterClaim = "3:cancel-task-roots"`（`:274`，断言在 `:190-207`） | 可照；它钉的是"进程自己承认注册了哪一步"，不是退出顺序 |

**共用骨架的逐字锚点**（照抄时不用另找）：
`bootResidentLeg` `resident_sink_nail_127_windows_test.go:217-240`（`exec.Command(exe)` 无参、`WISP_ENV=test`＋`WISP_TEST_DATA_DIR`、`CREATE_NEW_PROCESS_GROUP` `:222`、`t.Cleanup(leg.stop)` `:238`）；
`breakToLoop` `:252-254`；`exitedWithin(attempts)`＝50ms 定拍、注释明写"no wall-clock deadline arithmetic, D42#9" `:256-265`；
`pollUntil127` `:317`；`readResidentSink` `:331`；`installRecordIndex127` `:144`。

**这一族的致命上限（现量）**：`buildWispForTest`（`cmd/wisp/secret_argv_windows_test.go:161-199`）的构建命令逐字是
```
174:	cmd := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp")
```
**零 `-ldflags`**（我把 `ldflags`／`H=windowsgui` 两个 pattern 在整个 `cmd/wisp` 里搜过，只命中无关文本）。
⇒ 台件二进制永远是 CUI 子系统，**这一族天生看不见 GUI 子系统下的那一支**：不是"没测"，是**产物形态决定了 `GenerateConsoleCtrlEvent` 无处可投**。
（这一条从 §⑥ 的"待核"升为现量。）

### B-1b 族①的第四枚成员（本普查新发现，任务描述里没有）

`cmd/wisp/resident_ball_live_228_windows_test.go:89` `TestLive228ResidentLegOwnsABallWindowOnTheDesktop`，tag `//go:build windows && winlive`（`:1`）。
它同时具备**真机进程**与**真窗句柄**两面：`bootResidentLeg`（`:101`）→ `pollUntil127` 等姿态句（`:109-116`）→ `findBallWindows228` 按 `leg.pid()` 认领窗口（`:126-134`）→ `len(mine) != 1` 即红（`:139-142`）→ `breakToLoop()`（`:172`）→ **要求窗口灭**（`:182-191` `isWindow228`）。
⇒ **这一枚是 AC#11 ⓑ 最接近的现成形状**：它已经把"另一个进程的球窗"当被测量；唯一缺口是把 `:172` 那一次控制台信号**换成一次托盘退出**。
它同时把两枚限制写死在 header 里（`:22-29`）：winlive 两包不得共享桌面（`go test -tags winlive ./...` 会并发跑），且"alone 是规矩不是偏好"。

### B-2 族②：AST／静态走形族（读源码断形状，跨平台）

成员：`cmd/wisp/resident_ball_228_test.go`（`:198 TestAC228ResidentLegIsTheBallHost` 四支、`:275 TestAC228BallHostAnswersEveryGesture`），
名册 `ballEventCallbacks228` 十枚逐字在 `:50-53`。
同族其它件（`grep -n 'go/parser|go/ast' --include=*_test.go` 共 14 枚文件，逐字命中）：
`cmd/wisp/leg_sink_gate_131_test.go`、`cmd/wisp/leg_dispatch_gate_133_test.go`、`cmd/wisp/panel_host_gate_test.go`、
`cmd/wisp/panel_inbound_33_test.go`、`cmd/wisp/panel_assets_143_test.go`、`cmd/wisp/subagent_carrier_197_test.go`、
`internal/panel/l2_grant_boundary_test.go`（⛔ 冻结件，本腿只列名不动）、`internal/panel/composer_dispatch_test.go`、
`internal/ball/tokens_table_test.go`、`internal/tools/task_state_188_test.go`、`internal/projctx/projctx_test.go`、
`internal/agent/spill_path_invariant_test.go`、`internal/memory/artifacts_path_invariant_test.go`。

可照判断：**能照，且它自己就是"会响"的一族**——`resident_ball_228_test.go:319-322` 在宿主设的回调数**超过**名册数时红：
```
319:	if len(set) > len(ballEventCallbacks228) {
320:		t.Fatalf("AC#1 RED (the instrument, not the code): the host sets %d callbacks but this file lists %d; "+
```
⇒ 若修法**新加一枚 `Events` 键**（例如 `OnTrayExitRequested`），这一枚会立刻红并要求人改名册；
若修法**复用现有 `OnTrayExit`**，它一声不吭。这是本普查对"这一族能盯到什么程度"的准确边界。
**上限**：它读的是 `Events` 复合字面量的**键名**（`:288-304`），**完全看不见执行体**——
`recordTrayExit` 今天只打一句 `outcome=ignored`，而 `:275` 这一枚**今天是绿的**。⇒ 这就是任务里说的那形"回调非 nil 的假绿"。

### B-3 族③：`winlive` tag 族（`internal/ball` 同包，真窗口）

成员（逐字 tag `//go:build windows && winlive`）：`interaction_live_test.go:1`（用例 `:53`、`:133`、`:282`、`:351`）、
`hotkey_live_test.go:1`（`:121`、`:225`、`:318`、`:386`）、`live_windows_test.go:1`、`live_guard_windows_test.go:1`（工具件：`newLiveBall:187`、`sendMouse:251`、`sendMessage:261`、`injectHotkey:274`、`pSendMessageW:36`）。
header 的自我定位逐字（`interaction_live_test.go:11-15`）：
```
11: // Each test drives the real window procedure through posted/sent Win32
12: // messages or injected input, and asserts on what the machine and the window
13: // ended up as. Where a clause needs a human hand that no harness can supply
14: // (a second monitor, the felt experience of focus) the test says so in this
15: // file rather than faking a proof: see the notes at each clause.
```
现读的**真**上限（不是我上一次的错读，见 §⑥ E4）：`interaction_live_test.go:14` 把"没有 harness 能供的那只手"写作
**a second monitor** 与 **the felt experience of focus** 两例，**没有写托盘菜单**。
⇒ 但"今天这一族有没有覆盖托盘"仍是否定答案，凭据换成尺：**全仓 `*_test.go` 里 `wmAppTray`／`showMenu`／`menuExit` 零命中**（R13）——
即：托盘命令分流**没有任何台件**（三族皆无）碰过。`tray` 一词在测试里只以"被拆掉的图标"出现（`resident_ball_228_windows_test.go:45` 那条日志字面量、`:176` 注释）。

可照判断：**这一族是唯一能把消息喂进真 wndproc 的**；但它是 `package ball` 的台件，
而"十步有序退出"住在 `cmd/wisp`＋`internal/proc` ⇒ 单用这一族只能证"点击到了回调"，证不了"回调走到了十步"。

### B-4 具名结论（任务要的那句）

- **今天没有任何一族能"不加新东西"完整覆盖"托盘 → 有序退出"**：
  族①（含 ①b）在真进程里但只有控制台触发器（`secret_argv:174` 无 `-ldflags` ⇒ CUI；`breakToLoop` 是 `GenerateConsoleCtrlEvent`）；
  族③能把消息喂进真 wndproc（`live_guard:251-263`）但**不在那个会走十步的进程里**，且要拿 `menuExit` 得让 `TrackPopupMenu` 真返回 4。
- **只能覆盖"回调非 nil"这种假绿的＝族②**，凭据是 `resident_ball_228_test.go:275-322` 它今天对 `recordTrayExit` 的零执行体**一声不吭**。
- **最省的一格**（不是建议，是现量）：族① 的 ①b（`resident_ball_live_228:172`）那一次触发器替换，
  与族① 的 127 那枚 `stepsOf127` 从 `t.Logf`（`:640-641`）升成断言——这两处的**被测量已经在码里**，缺的只是触发源与一句断言。

---

## ③ C 空窗那一坑的可测性

### C-1 空窗的两端（现读）
- 入口注册：`cmd/wisp/resident_windows.go:106-108`
  ```
  106:	bootExit := make(chan os.Signal, 1)
  107:	signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)
  108:	defer signal.Stop(bootExit)
  ```
- 空窗消费点：`resident_windows.go:194-203`（`var reason string` / `select` / `case sig := <-bootExit:` `:196` / `default:` `:199` / `reason = rt.RunEventLoop()` `:202`）。
- 环路本体：`internal/proc/boot_windows.go:127-152`——**每拍 50ms**（`WaitForSingleObject(rt.Instance.ActivateEvent(), 50)` `:133` ＋ `select { case <-sigCtx.Done(): return "signal"; case <-time.After(50*time.Millisecond): }` `:145-150`）；
  唯一停机输入是 `:128` 的 `signal.NotifyContext(..., os.Interrupt, syscall.SIGTERM)`。
- 十步全挂在正常返回上（现量）：`resident_windows.go:70 defer sink.close()`、`:72-82 defer rt.Shutdown(false)`、`:137 defer rb.stop()`、`:144 defer ra.detachBall()`、`:202 rt.RunEventLoop()`。

### C-2 现成有没有**任何**一台件能制造"球已起、环路未进"的窗口？
**有，且只有一枚**：`cmd/wisp/resident_ball_228_windows_test.go:184`。机制复述（逐字行号，非只给名字）：
1. `:185-188` 起真 `wisp.exe`（无参＝常驻腿），`dataDir := t.TempDir()`，`sinkDir := logSinkDir(dataDir)`（定义 `cmd/wisp/logsink.go:87`）。
2. **时刻锚＝磁盘事实，不是时间**：`:190-196`
   ```
   190:	if !pollUntil127(200, func() bool {
   191:		n, err := observe.CountLogFiles(sinkDir)
   192:		return err == nil && n > 0
   193:	}) {
   ```
   header `:179-182` 逐字给出这样挑的理由：*"The trigger is the sink FILE appearing, which is the earliest observable moment inside that window
   (the rolling writer opens its file at install time, ~270ms before the loop is entered on the bench machine)"*。
   ⇒ 也就是说它把"球已起、环路未进"这段窗口的**存在性绑在 sink 文件出现这一刻**，而不是绑在 sleep 上。
3. `:197-200` 立刻 `leg.breakToLoop()`＝`GenerateConsoleCtrlEvent(CTRL_BREAK_EVENT, leg.pid())`（`127:252-254`）。
4. `:201-209` 退出码断言（`exitedWithin(400)` 定拍；`leg.waitErr != nil` 即红，并把 `0xc000013a` 点名成"什么都没注册"的症状）。
5. `:211-220` 读落盘记录，要求 install 记录与 `residentShutdownRecord`（`127:131 = "shutdown step"`）都在。
⇒ **删除 `resident_windows.go:106-108` 那块 `signal.Notify`，这一枚是红的那一个名字**（`:181-183` 逐字这么承诺）。

另一枚用**同一个锚点**的是 `resident_sink_nail_127_windows_test.go:574-577`，但它随后不宣称自己在空窗内（它要的是顺序）。
⇒ **空窗这一维今天只由一枚台件钉住，且只用控制台那一支。**
还有一格**完全无人钉**：`resident_windows.go:103-105` 注释自承留下的最后一段缝——"the few instructions between the check below and
`NotifyContext` inside `RunEventLoop` - which costs one more Ctrl+C"。现量：没有任何台件打第二次触发。

### C-3 若停机请求位也接进这段，最小可测形状是什么
**不选边**（§⑦ F3）。可测性的**共同骨架**是清楚的，逐条给现量：
- 时刻锚：沿用 `observe.CountLogFiles`（`:191`）——它不需要新接口，`internal/observe` 已被 `cmd/wisp` 台件导入。
- 投递物：**这是唯一缺口**。今天这一族的投递物只有 `GenerateConsoleCtrlEvent` 一枚；GUI 子系统下它不可用（B-1 的 `-ldflags` 现量）。
- 判据面：`exitedWithin`＋`leg.waitErr`＋`readResidentSink`＋`residentShutdownRecord` 四件套全部现成可照（`:201-220`）。
⇒ 换句话说：**"空窗里按下的退出被吞掉"这一维，可测性不缺判据、不缺时刻锚，只缺一枚能在 GUI 子系统投递的触发物。**

---

## ④ D 十步与审计的可测面

### D-1 新加一枚触发时，这些会不会自动变红？——**编排者的判断成立：不会。**
凭据（全部现读）：
1. `RunShutdownSequence(hooks ShutdownHooks, opts ShutdownOptions) []StepRecord`（`internal/proc/shutdown.go:112`）
   与 `Runtime.Shutdown(fast bool)`（`boot_windows.go:164`）的签名里**没有任何触发源概念**。
2. `TestShutdownOrderAudit`（`internal/proc/shutdown_test.go:11`）**手搓** `ShutdownHooks{}` 八槽（`:13-27`）后直接调 `RunShutdownSequence`（`:30`），
   全程不引用 `cmd/wisp`、不引用 `Runtime`、不引用信号。⇒ 它与"谁要求退出"**结构性无关**。
3. `hookableRoster` 是**闭集八枚**（`shutdown_hooks.go:47-56`，注释 `:44-46` 逐字"built from the sequence's own constants, so it cannot drift into offering a ninth hookable step"）；
   `Register` 的三道拒绝面在 `:89/:92/:95`（sealed / slot unknown / already registered）。新触发**不注册新步**⇒ 这三面一声不吭。
4. `shutdown_hooks_test.go` 的三枚真名（我先前引错了一枚，见 §⑥ E9）：
   `TestShutdownHookSetRunsEveryRegisteredStepInOrder:31`、`TestShutdownHookSetRefusesEveryShapeThatWouldLie:68`、`TestShutdownHookSetZeroValueIsUsable:120`。
   它们同样只测 hook 集，不测触发源。
5. `recordTrayExit` 今天那句 `outcome=ignored`（`resident_ball_windows.go:255`）**没有任何台件引用**：
   `grep -n 'tray exit requested|outcome.*ignored' --include=*.go` 只命中 `:255`/`:256` 两行（生产文件自身）。
   ⇒ **把"ignored"改成"executed"今天撞不到任何断言。**

### D-2 但有两枚**会**自动变红（这才是有用的一半）
1. **起第七枚常驻协程**：`ResidentNames` 是六枚（`internal/observe/goroutine.go:42-43`，注释逐字 "the 6 resident roster names (D38b)"），
   `goroutine.go:65` 判名、`:422` 算 `rep.ResidentOverBaseline = rep.Resident > ResidentBaseline`；
   `internal/proc/boot_windows.go:115-118` 在 Boot 里**直接拒绝启动**；`internal/observe/goroutine_test.go:268-269` 钉 `ResidentBaseline == len(ResidentNames)`。
   ⇒ 这正是 `resident_ball_windows.go:250-252` 注释里那句"a seventh live goroutine that is not in the frozen D38b resident roster (observe.ResidentNames)"的**仪器落点**：它不需要新台件就会红。
2. **给 `internal/ball.Events` 新加一枚键**：`resident_ball_228_test.go:319-322`（见 B-2 的逐字读数）。
3. 附带一枚**非红**但会被打的东西：`tools/d22scan` ban #1（裸 `go`）的射程逐字是"production scope internal/ + cmd/ (non-test, non-testdata)"（`tools/d22scan/main.go:5-6`），
   而 `main.go:44-46` 逐字写明**只有 ban #8** 的射程"`_test.go` INCLUDED"。⇒ **测试文件里为并发而写的 `go func()` 不会被 ban #1 打**；生产文件里的会被。

### D-3 "审计记录里 1..10 都在、且顺序对"这句话今天由谁在钉
| 维 | 谁钉 | 逐字凭据 | 射程 |
|---|---|---|---|
| 执行顺序 1..9 ＋ records 恰 10 | `internal/proc/shutdown_test.go:11` | `:36-40 want` 九枚；`:41-48` 逐位比；`:52-53 len(records) != 10` | **进程内**，hooks 手搓，与触发源无关 |
| 无钩子时仍 10 条、8/10 永不 skipped | 同文件 `:158 TestShutdownNilHooksAreRecordedSkipped` | `:160-174` | 进程内 |
| fast path 只跳第 7 步 | 同文件 `:73 TestShutdownFastPathSkipsOnlyStep7` | `:104-117` | 进程内 |
| 注册面→真序列 | `internal/proc/shutdown_hooks_test.go:31` | （名册见 D-1 第 4 条） | 进程内 |
| **真机进程里 step 3 真的被执行且 records 恰 10** | `cmd/wisp/resident_approval_246_windows_test.go:132` | `:145-148 records := rt.Shutdown(false); len(records) != 10`；`:149-158 step3 := records[proc.StepCancelTasks-1]` 且要求 `!Skipped`/`Err==nil`/`Name=="cancel-task-roots"` | **进程内**（但它走的是真 `proc.Boot`＋真 `RegisterShutdownHook`＋真 `Shutdown`，是**离真机最近的一枚十步读数**） |
| **真机进程（子进程）里的 1..10 顺序** | **没有人钉** | `resident_sink_nail_127_windows_test.go:623-641`：trail 非空 `:629`、**最后一条是 shutdown 记录** `:633-635`、console 有 `residentShutdownLine` `:637`；**步号只进 `t.Logf`（`:640-641`，用 `stepsOf127` 读 `r.Step`）** | 只钉"顺序在 sink 关闭之前"，不钉 1..10 |

⇒ **结构性障碍（本普查最要紧的一格读数）**：`internal/proc` 的日志面**只为异常步出声**——
skipped 才写（`shutdown.go:130 "shutdown step skipped (module not present)" step name`）、
超时/失败才写（`:145`、`:147`）、fast 才写（`:172`）；
**执行成功的步零日志**（`shutdown.go:140-150` 里 err==nil 分支不 log），
第 8 步（`:178-180`）与第 10 步（`:186`）**只 append record，零日志**。
⇒ 所以在**盘面上**看到"1..10 全在"这件事**今天结构上不可能**——不是缺台件，是缺**被测量**。
⇒ `sinkInstallRecord` **已经有 `Step` 字段**（`stepsOf127` 读 `r.Step`，`resident_sink_nail_127_windows_test.go:652-658`），
  即解析器已能把步号从 jsonl 里拿出来；缺的只是"每步都落一行"。
⇒ 结论交给编排者：**AC#11 ⓑ 那句"审计顺序仍是 1..10"在真机侧目前无法读**，
  要么读面换到进程内（照 `resident_approval_246_windows_test.go:145-158`），要么改 `internal/proc` 的日志面（那是动 D38(e) 附近的东西）。两枚都不是普查腿能定的，见 §⑦ F5。

### D-4 常驻腿今天真注册的步（供"该看到什么"对表）
`resident_windows.go:157-161` 注册 `proc.StepCancelTasks`（step 3）；`boot_windows.go:169-171` 若 `CloseJob` 无人注册则**由 Shutdown 自己补上**（⇒ step 9 恒有主）；
`resident_windows.go:177 startResidentTaskSource(rt, ra)` 据其注释拥有 step 1 与 7（`resident_windows.go:166-167` 逐字"books the two D38(e) steps that entry owns (1 and 7)"）；
boot 报告把注册表打印出来：`resident_windows.go:181-187`（`registered := rt.RegisteredShutdownSteps()`，逐枚 `fmt.Sprintf("%d:%s", int(s), s.Name())`）。
⇒ step 2 **不在 hook 集里**，它由 `defer rb.stop()`（`:137`）在序列之前做（注释 `:127-131` 与 `resident_ball_windows.go:177-182` 双向对上）。

---

## ⑤ 我跑了哪些尺、每条真实读数

> 每条＝完整尺＋逐字读数。所有 `file:line` 都由带 `-n` 的词面尺反取（Read 工具的输出是 `cat -n` 绝对行号，同样逐字可核）。
> 本腿**未跑**任何 `go`／`build`／`test`／`vet`／`ps1`。

### R1 锚点尺（同发取）
尺：`date "+%Y-%m-%d %H:%M:%S %z"; echo ---LOG---; git log -1 --format=%H; echo ---BRANCH---; git rev-parse --abbrev-ref HEAD; echo ---DIRS---; ls -la ".scratch/wisp/probes/228/"; echo ---PROBES---; ls ".scratch/wisp/probes/"`
逐字读数（节选）：
```
2026-09-30 11:06:22 +0800
---LOG---
0c8b9fddd2c3da647a76f0757fd466ef6767e6c1
---BRANCH---
dev
---DIRS---
total 24
drwxr-xr-x 1 swq 197609 0 Sep 29 16:05 a1
drwxr-xr-x 1 swq 197609 0 Sep 30 16:19 r1
drwxr-xr-x 1 swq 197609 0 Sep 30 16:56 v1
```
⇒ `a2/` 当时不存在。`wc` 同时给出七枚主体文件行数（见 R2）。

### R2 尺：主体文件行数/字节数
尺：`wc -l -c internal/ball/tray_windows.go cmd/wisp/resident_ball_windows.go cmd/wisp/resident_windows.go internal/proc/shutdown.go internal/proc/shutdown_hooks.go internal/proc/shutdown_test.go internal/ball/ball_windows.go`
逐字读数：
```
 109  3090 internal/ball/tray_windows.go
  271 13792 cmd/wisp/resident_ball_windows.go
  205 10283 cmd/wisp/resident_windows.go
  188  6815 internal/proc/shutdown.go
  158  6183 internal/proc/shutdown_hooks.go
  175  5745 internal/proc/shutdown_test.go
 1003 31709 internal/ball/ball_windows.go
 2109 77617 total
```

### R3 尺：停机入口／现有 stop-request 位
尺（Grep，`-n`，glob `*.go`）：`RunEventLoop|NotifyContext|shutdownRequest|stopRequest`
逐字读数（全部 7 行命中，无分页损失）：
```
cmd\wisp\resident_windows.go:89:	// The handler that stops this process is installed inside RunEventLoop, and
cmd\wisp\resident_windows.go:104:	// check below and NotifyContext inside RunEventLoop - which costs one more
cmd\wisp\resident_windows.go:202:		reason = rt.RunEventLoop()
cmd\wisp\resident_sink_nail_127_windows_test.go:250:// turns that into os.Interrupt, proc.RunEventLoop returns "signal", and the
internal\proc\boot_windows.go:122:// RunEventLoop is the empty event loop of ticket 03: it blocks until an exit
internal\proc\boot_windows.go:127:func (rt *Runtime) RunEventLoop() string {
internal\proc\boot_windows.go:128:	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
```
⇒ **`shutdownRequest`／`stopRequest` 全仓零命中**⇒ AC#11 要加的位是全新面；`RunEventLoop` 唯一定义 `boot_windows.go:127`，唯一输入 `:128`。

### R4 尺：`internal/ball` 导出面
尺（Grep，`-n`，path `internal/ball`）：`^func \(|^func [A-Z]|^\tOn[A-Z]`（120 条分页首 120）
关键逐字读数（其余为 dock/hotkey/anim，与本件无关）：
```
internal\ball\ball_windows.go:50:	OnClickBall     func()
internal\ball\ball_windows.go:55:	OnTrayPanel     func()
internal\ball\ball_windows.go:56:	OnTrayMute      func()
internal\ball\ball_windows.go:57:	OnTrayPauseWake func()
internal\ball\ball_windows.go:58:	OnTrayExit      func()
internal\ball\ball_windows.go:59:	OnDragEnd       func()
internal\ball\ball_windows.go:143:func New(opts Options) (*Ball, error) {
internal\ball\ball_windows.go:549:func (b *Ball) wndProc(hwnd, m, wParam, lParam uintptr) uintptr {
internal\ball\ball_windows.go:728:func (b *Ball) fire(fn func()) {
internal\ball\ball_windows.go:741:func (b *Ball) uiRun(fn func()) {
internal\ball\ball_windows.go:919:func (b *Ball) SetTrayChecks(muted, pausedWake bool) {
internal\ball\ball_windows.go:927:func (b *Ball) SetTrayTip(tip string) {
internal\ball\ball_windows.go:935:func (b *Ball) Close() {
internal\ball\ball_windows.go:1003:func (b *Ball) DebugHWND() windows.HWND { return b.hwnd }
```

### R5 尺：（本腿第一次提交前未单列，内容并入 R4/R6/R7/R8；编号留空是为了不与首发 §⑤ 的号串错位）

### R6 尺：`recordTrayExit` 与 `Events` 注入点
尺（Grep，`-n`，path `cmd/wisp`）：`recordTrayExit|OnTrayExit|outcome.*ignored`
```
resident_ball_windows.go:131:			OnTrayExit:      recordTrayExit,
resident_ball_windows.go:246:// recordTrayExit is the tray's 退出 item. The request is recorded; the stop is
resident_ball_windows.go:253:func recordTrayExit() {
resident_ball_windows.go:254:	const why = "this process leaves when its event loop returns, which today only a console signal (Ctrl+C) asks for; the tray Exit item has no stop path attached to it"
resident_ball_windows.go:255:	slog.Warn("tray exit requested", "outcome", "ignored", "why", why)
resident_ball_windows.go:256:	fmt.Printf("wisp: tray exit requested: %s\n", why)
```
⇒ 包级零参非导出函数 ⇒ `package main` 同包台件可直接调；且 `:255/:256` 是全仓唯一命中（见 R18 ⇒ 改它不撞断言）。

### R7 尺：托盘菜单四枚 id 与 `TrackPopupMenu` 阻塞返回
尺：Read `internal/ball/tray_windows.go` 全文（109 行，与 R2 相符）
关键逐字（行号为该文件绝对行号）：
```
15-16: // Tray menu command ids (wParam of WM_COMMAND after TPM_RETURNCMD is NOT
16:    // used; we use TrackPopupMenu's return value inline in the wndproc instead).
18:	trayUID = 0x5701
20-23:	menuOpenPanel = 1 / menuMute = 2 / menuPauseWake = 3 / menuExit = 4
36:		uFlags:           nifMessage | nifIcon | nifTip,
37:		uCallbackMessage: wmAppTray,
72:func showMenu(hwnd windows.HWND, muted, pausedWake bool) uint32 {
97:	pSetForegroundWindow.Call(uintptr(hwnd))
100:	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
101:	sel, _, _ := pTrackPopupMenu.Call(menu, tpmRightButton|tpmReturnCMD|tpmNoActivate,
102:		uintptr(pt.x), uintptr(pt.y), 0, uintptr(hwnd), 0)
105:	pPostMessageW.Call(uintptr(hwnd), wmNull, 0, 0)
106:	return uint32(sel)
```

### R8 尺：常量表与消息分流（含 message-only 判定）
尺（Grep，`-n`，glob `*.go`）：`wmAppTray|wmAppTask|nifMessage|tpmReturnCMD|pPostMessageW|pSendMessageW`
```
internal\ball\win32_windows.go:33:	pPostMessageW                  = modUser32.NewProc("PostMessageW")
internal\ball\win32_windows.go:102:	wmAppTask  = wmApp + 0x201 // run a closure posted to the STA thread
internal\ball\win32_windows.go:103:	wmAppTray  = wmApp + 0x202 // tray icon callback
internal\ball\win32_windows.go:129:	tpmReturnCMD   = 0x0100
internal\ball\win32_windows.go:265:	nifMessage = 1
internal\ball\sta_windows.go:142:	pPostMessageW.Call(uintptr(hwnd), wmAppTask, uintptr(id), 0)
internal\ball\live_windows_test.go:557:			pPostMessageW.Call(uintptr(b.hwnd), uintptr(msg), 0, lParam)
internal\ball\live_guard_windows_test.go:36:	pSendMessageW             = modUser32.NewProc("SendMessageW")
internal\ball\live_guard_windows_test.go:254:		r, _, _ := pSendMessageW.Call(uintptr(b.hwnd), msg, wParam, lParam)
internal\ball\ball_windows.go:669:	case wmAppTray:
internal\ball\ball_windows.go:688:	case wmAppTask:
```
⇒ `HWND_MESSAGE` 未出现在任何命中里（尺里带了这条 pattern 的一枚变体见 R9b）⇒ **不是 message-only**。

### R9 尺：`ball_windows.go` 托盘分流的绝对行号
尺：Read `internal/ball/ball_windows.go` offset=600 limit=120（绝对行号由 `cat -n` 型尺给）
```
669:	case wmAppTray:
670:		switch lParam & 0xFFFF {
671:		case wmLButtonUp: // left click = open panel (SPEC-08 §7)
672:			b.fire(b.opts.Events.OnTrayPanel)
673:		case wmRButtonUp:
674:			sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)
675:			switch sel {
676:			case menuOpenPanel:
677:				b.fire(b.opts.Events.OnTrayPanel)
682:			case menuExit:
683:				b.fire(b.opts.Events.OnTrayExit)
686:		return 0
```
**R9b 复核尺**（专门补 message-only）：Grep（`-n`，path `internal/ball`）：`notification area|tray menu|go:build`
⇒ `go:build` 全量命中（见 R12）；**`notification area` 与 `tray menu` 在 `internal/ball` 零命中**；`HWND_MESSAGE` 亦零命中。
⇒ 这次尺**推翻了我第一次提交里 R11 的一段引用**，见 §⑥ E4。

### R10 尺：228 真机两枚用例的机制
尺（Grep，`-n`，path `cmd/wisp`）：`bootResidentLeg|breakToLoop|GenerateConsoleCtrlEvent|exitedWithin|pollUntil127|CountLogFiles|residentShutdownRecord|buildWispForTest`
```
resident_ball_228_windows_test.go:54:	exe := buildWispForTest(t)
resident_ball_228_windows_test.go:56:	leg := bootResidentLeg(t, exe, dataDir)
resident_ball_228_windows_test.go:61:	sawVerdict := pollUntil127(200, func() bool {
resident_ball_228_windows_test.go:89:	if err := leg.breakToLoop(); err != nil {
resident_ball_228_windows_test.go:94:	if !leg.exitedWithin(400) {
resident_ball_228_windows_test.go:104:	recs := readResidentSink(t, sinkDir)
resident_ball_228_windows_test.go:162:		firstShutdown := firstIndexOfContains228(recs[stoppedIdx+1:], residentShutdownRecord)
resident_ball_228_windows_test.go:184:func TestAC228ExitRequestDuringBootStillLeavesThroughD38E(t *testing.T) {
resident_ball_228_windows_test.go:190:	if !pollUntil127(200, func() bool {
resident_ball_228_windows_test.go:191:		n, err := observe.CountLogFiles(sinkDir)
resident_ball_228_windows_test.go:197:	if err := leg.breakToLoop(); err != nil {
resident_ball_228_windows_test.go:216:	shutdownIdx := firstIndexOfContains228(recs, residentShutdownRecord)
```
（`resident_sink_nail_127_windows_test.go` 的 helper 定义行号见 B-1 表末"共用骨架"。）

### R11（首发编号保留）尺：空窗用例的完整 Read
尺：Read `cmd/wisp/resident_ball_228_windows_test.go` 全文 245 行 ＋ Read `cmd/wisp/resident_ball_228_test.go` 全文 344 行。
关键逐字（前者 `:179-183` header、`:190-196` 锚、`:205-209` 退出码、`:216-220` 记录；后者 `:50-53` 名册、`:275-322` 反沉默支）。
**同一次我误把 winlive 注释当命中抄进 R11——已在 R9b/R13 反取真词面并改正，见 §⑥ E4。**

### R12 尺：三族成员与 tag
尺（Grep，`-n`，path `internal/ball`）：`go:build`（含 -A 3）＋（Grep，`-n`，glob `*.go`）：`func newLiveBall|func liveHotkeys|func requireFreeOfDock|go:build windows && winlive`
```
internal\ball\interaction_live_test.go:1://go:build windows && winlive
internal\ball\hotkey_live_test.go:1://go:build windows && winlive
internal\ball\live_windows_test.go:1://go:build windows && winlive
internal\ball\live_guard_windows_test.go:1://go:build windows && winlive
internal\ball\live_guard_windows_test.go:187:func newLiveBall(t *testing.T, opts Options) *Ball {
internal\ball\hotkey_live_test.go:41:func liveHotkeys() HotkeyConfig {
internal\ball\interaction_live_test.go:394:func requireFreeOfDock(t *testing.T, b *Ball) {
cmd\wisp\resident_ball_live_228_windows_test.go:1://go:build windows && winlive
cmd\wisp\resident_task_source_live_246_windows_test.go:1://go:build windows && winlive
cmd\wisp\resident_approval_live_246_windows_test.go:1://go:build windows && winlive
```
⇒ 逐字纠正首发 R11：tag 是 `windows && winlive`，**不是** `winlive` 单独。⇒ `cmd/wisp` 里有**三枚** winlive 件（首发我只数到 `internal/ball` 那两族）。

### R13 尺：全仓测试对托盘命令分流的覆盖（本件的核心否定读数）
尺（Grep，`-n`，glob `*_test.go`）：`wmAppTray|showMenu|menuExit|recordTrayExit|OnTrayExit|tray`
⇒ **`wmAppTray`／`showMenu`／`menuExit`／`recordTrayExit` 在全部 `*_test.go` 零命中。**
`OnTrayExit` 唯一命中是名册字符串：`cmd\wisp\resident_ball_228_test.go:52:	"OnTrayPanel", "OnTrayMute", "OnTrayPauseWake", "OnTrayExit", "OnDragEnd",`
`tray` 的实质命中：
```
cmd\wisp\resident_ball_228_windows_test.go:45:	ballStoppedMsg = "ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
cmd\wisp\resident_ball_228_windows_test.go:176:// used to get 0xc000013a - no D38(e) trail, no tray removal, no hot key
cmd\wisp\approval_seam_201_test.go:8:// a ball click handler or a tray menu item is not a terminal. So both cases below
```
⇒ 结论：**今天没有任何台件碰过托盘命令分流**（三族皆无）；`approval_seam_201_test.go:8` 那句注释是本件找到的唯一一处把"tray menu item"当**语义排除项**写的地方。

### R14 尺：AST／静态族成员
尺（Grep，`-n`，glob `*_test.go`，files_with_matches）：`go/parser|go/ast` ⇒ 14 枚文件（名单见 B-2）。

### R15 尺：真机 helper 定义面
尺（Grep，`-n`，path `cmd/wisp`）：`func bootResidentLeg|func .*breakToLoop|func buildWispForTest|func readResidentSink|func installRecordIndex127|ldflags|H=windowsgui`
```
cmd\wisp\resident_sink_nail_127_windows_test.go:144:func installRecordIndex127(recs []sinkInstallRecord) int {
cmd\wisp\resident_sink_nail_127_windows_test.go:217:func bootResidentLeg(t *testing.T, exe, dataDir string) *residentLeg {
cmd\wisp\resident_sink_nail_127_windows_test.go:252:func (l *residentLeg) breakToLoop() error {
cmd\wisp\resident_sink_nail_127_windows_test.go:331:func readResidentSink(t *testing.T, dir string) []sinkInstallRecord {
cmd\wisp\secret_argv_windows_test.go:161:func buildWispForTest(t *testing.T) string {
cmd\wisp\resident_task_source_246_windows_test.go:455:func bootResidentLegWithEnv(t *testing.T, exe string, env []string) *residentLeg {
```
⇒ **`ldflags` 与 `H=windowsgui` 在 `cmd/wisp` 零命中**（尺里带了这两条 pattern）⇒ B-1 末的 CUI 结论成立。

### R16 尺：`buildWispForTest` 与 `bootResidentLeg`/`breakToLoop` 逐字
尺：Read `cmd/wisp/secret_argv_windows_test.go:150-199`；Read `cmd/wisp/resident_sink_nail_127_windows_test.go:210-290`
```
secret_argv_windows_test.go:174:	cmd := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp")
resident_sink_nail_127_windows_test.go:219:	cmd := exec.Command(exe)
resident_sink_nail_127_windows_test.go:221:	cmd.Env = append(os.Environ(), "WISP_ENV=test", procTestDataDirEnv+"="+dataDir)
resident_sink_nail_127_windows_test.go:222:	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
resident_sink_nail_127_windows_test.go:238:	t.Cleanup(leg.stop)
resident_sink_nail_127_windows_test.go:253:	return windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, l.pid())
resident_sink_nail_127_windows_test.go:257:// through its own shutdown path (no wall-clock deadline arithmetic, D42#9).
```
（`secret_argv:158-160` 的 header 逐字："compiles the real binary (the ticket's claim is about the shipped command line, not a test harness's)"。）

### R17 尺：winlive 消息注入面的逐字
尺：Read `internal/ball/live_guard_windows_test.go:215-290`；Read `internal/ball/interaction_live_test.go:1-135`；Read `internal/ball/live_windows_test.go:500-590`
```
live_guard_windows_test.go:251-258	sendMouse（见 A-1 第 2 支的逐字块）
live_guard_windows_test.go:261-263	sendMessage = sendMouse 的别名
live_guard_windows_test.go:274-275	injectHotkey(... ) { press := func(v uintptr, up bool) { ... pKeybdEvent.Call(v, 0, f, 0) } }
interaction_live_test.go:35:var procMouseEvent = modUser32.NewProc("mouse_event")
interaction_live_test.go:58-78:	newLiveBall(t, Options{... Events: Events{OnClickBall: func(){ ... m.Dispatch(statemachine.EvSummon, nil) ...}}})
interaction_live_test.go:81:	hwnd := b.DebugHWND()
interaction_live_test.go:89-91:	sendMouse(b, wmLButtonDown/wmMouseMove/wmLButtonUp, ...)
interaction_live_test.go:117-121:	读回 SetState 与 m.State() 两维断言
live_windows_test.go:553-562:	postMsg := func(msg uint32) { b.sta.PostTask(func() { pPostMessageW.Call(...) }) }
hotkey_live_test.go:327-328:	h := b.sta ; h.PostTask(func() { bump(&mu, hits, "mute-cb") })
```
⇒ 这一族的**断言风格**现量：回调里 `mu.Lock(); clicks++`（`interaction_live_test.go:63-65`）＋ 状态机与球面双读（`:117-122`）。

### R18 尺：`recordTrayExit` 那句 ignored 有没有被钉
尺（Grep，`-n`，glob `*.go`）：`tray exit requested|outcome.*ignored`
```
cmd\wisp\resident_ball_windows.go:255:	slog.Warn("tray exit requested", "outcome", "ignored", "why", why)
cmd\wisp\resident_ball_windows.go:256:	fmt.Printf("wisp: tray exit requested: %s\n", why)
```
⇒ 仅生产文件自身两处 ⇒ **没有任何台件引用这两句**（D-1 第 5 条）。

### R19 尺：`internal/proc` 的日志面（决定盘面上能看见什么）
尺（Grep，`-n`，path `internal/proc`，-B 2）：`slog\.(Info|Warn|Error)`
```
internal\proc\shutdown.go:130:			slog.Info("shutdown step skipped (module not present)", "step", int(step), "name", rec.Name)
internal\proc\shutdown.go:145:				slog.Error("shutdown step abandoned wait (deadline exceeded)", "step", int(step), "name", rec.Name)
internal\proc\shutdown.go:147:				slog.Error("shutdown step failed", "step", int(step), "name", rec.Name, "err", err)
internal\proc\shutdown.go:172:		slog.Warn("fast shutdown: skipping non-critical log/db flushes (D38e fast path)")
internal\proc\boot_windows.go:139:				slog.Info("activation requested by second launch",
```
⇒ **成功执行的步零日志**（`shutdown.go:140-150` 只有 err!=nil 分支出声）；第 8/10 步（`:178-180`、`:186`）零日志。⇒ D-3 的结构性障碍。

### R20 尺：真机侧"审计顺序"断言面的逐字
尺（Grep，`-C 6`，`-n`，path `cmd/wisp`，glob `*_test.go`）：`residentShutdownRecord`
```
resident_sink_nail_127_windows_test.go:131:	residentShutdownRecord = "shutdown step"
resident_sink_nail_127_windows_test.go:625:		if strings.Contains(r.Msg, residentShutdownRecord) {
resident_sink_nail_127_windows_test.go:629:	if len(trail) == 0 {
resident_sink_nail_127_windows_test.go:634:	if !strings.Contains(last.Msg, residentShutdownRecord) {
resident_sink_nail_127_windows_test.go:637:	if !leg.stdout.has(residentShutdownLine) {
resident_sink_nail_127_windows_test.go:640:	t.Logf("RESIDENT LEG RECORDED ON DISK: %d shutdown record(s), steps %v, first=%q last=%q",
resident_sink_nail_127_windows_test.go:641:		len(trail), stepsOf127(trail), trail[0].Msg, trail[len(trail)-1].Msg)
resident_sink_nail_127_windows_test.go:652:func stepsOf127(recs []sinkInstallRecord) []int {
resident_sink_nail_127_windows_test.go:655:		out = append(out, r.Step)
resident_ball_228_windows_test.go:162 / :216（见 R10）
```
⇒ `sinkInstallRecord` **有 `Step` 字段**（`:655`）⇒ 步号可读；`:640` 只 `t.Logf` ⇒ **步序今天不是断言**。

### R21 尺：常驻协程名册（第 2 枚"会自动红"）
尺（Grep，`-n`，glob `*.go`）：`ResidentNames|ResidentOverBaseline|residentBaseline`
```
internal\observe\goroutine.go:42:// ResidentNames are the 6 resident roster names (D38b), in table order.
internal\observe\goroutine.go:43:var ResidentNames = []string{
internal\observe\goroutine.go:65:	if slices.Contains(ResidentNames, name) {
internal\observe\goroutine.go:195:	ResidentOverBaseline bool     // resident count > 6: leak symptom (D38b)
internal\observe\goroutine.go:422:	rep.ResidentOverBaseline = rep.Resident > ResidentBaseline
internal\observe\goroutine_test.go:106:	if rep.ResidentOverBaseline {
internal\observe\goroutine_test.go:268:	if ResidentBaseline != len(ResidentNames) {
internal\observe\goroutine_test.go:269:		t.Fatalf("ResidentBaseline = %d but ResidentNames has %d entries", ResidentBaseline, len(ResidentNames))
internal\proc\boot_windows.go:115:	if rep := rt.Registry.RosterReport(); rep.ResidentOverBaseline {
internal\agent\loop_golden_test.go:293:	if rep.ResidentOverBaseline {
cmd\wisp\resident_ball_windows.go:251:// resident roster (observe.ResidentNames) - both are new contracts, and ticket
```

### R22 尺：`tools/d22scan` 两枚 ban 的射程差别
尺（Grep，`-C 3`，`-n`，path `tools/d22scan`）：`HasSuffix\(.?"_test\.go"|non-test|prodScope|scopes =|bannedScope` ＋（Grep，`-C 4`）：`go func|ban 1`
```
tools\d22scan\main.go:5:// Bans (PLAN.md D22 fourth-round additions), production scope internal/ +
tools\d22scan\main.go:6:// cmd/ (non-test, non-testdata) unless stated:
tools\d22scan\main.go:44://	                      _test.go INCLUDED (D23). This is the one ban whose
tools\d22scan\main.go:45://	                      scope is NOT the "production, non-test" default
tools\d22scan\main.go:707:					"bare `go func(` is banned (D22/D38b): use observe.Registry.Spawn (named, owner, recover boundary)")
tools\d22scan\scan_test.go:295:	seedFile(t, root, "internal/ok/emoji_test.go", ...)
tools\d22scan\scan_test.go:309:		"internal/ok/emoji_test.go", // _test.go
```
⇒ ban #1（裸 `go`）走 `main.go:5-6` 的**production、非测试**默认射程；ban #8（emoji）是**唯一**明确含 `_test.go` 的一枚（`:44-46`）。

### R23 尺：246 那枚"进程内十步读数"的逐字
尺：Read `cmd/wisp/resident_approval_246_windows_test.go:118-207`
```
132:func TestAC246CancelStepHookRunsOnTheRealShutdownSequence(t *testing.T) {
133:	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(observe.NewRegistry()))
138:	if got := rt.RegisteredShutdownSteps(); len(got) != 0 {
141:	if err := rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots); err != nil {
145:	records := rt.Shutdown(false)
146:	if len(records) != 10 {
149:	step3 := records[proc.StepCancelTasks-1]
150:	if step3.Skipped {
156:	if step3.Name != "cancel-task-roots" {
184:func TestAC246ShippedResidentProcessOwnsItsCancelStep(t *testing.T) {
185:	exe := buildWispForTest(t)
190:	said := pollUntil127(400, func() bool {
274:	cancelStepRosterClaim = "3:cancel-task-roots"
```
⇒ **这是全仓唯一一枚"用真 `proc.Boot`＋真 `Shutdown` 读满 10 条 record"的 `cmd/wisp` 台件**（进程内）⇒ §⑦ F5 甲支的现成可照形状就是它。

### R24 尺：票面进度与提交
尺：`git commit` 第一发 → `34c28c5c8f8a756f9f7f9721ce849e7e05b80bdf`（2 files changed, 326 insertions；census 首发 312 行／26552 字节）。
本文件的第二发改的是 §①②③④ 正文与 §⑤⑥⑦ 的更正；两发都带显式 pathspec，未 push。

---

## ⑥ 我可能写错的条目（对抗我自己）

- **E1｜HEAD 与编排者给的 `3398e6f1` 不一致。** 现量 R1：`0c8b9fddd2...`。我**没有**跑 `git log --oneline -5`，所以不能断定是"另有提交落在其上"还是"记错了那枚"。
  ⇒ 影响：本件所有行号相对 `0c8b9fddd2`。尺全在 §⑤，可逐条复跑。
- **E2｜现读时刻的日期翻篇。** R1 的 `date` 原样打印 `2026-09-30 11:06:22 +0800`，而会话侧给出的当前日期在后续翻了到 10-01。
  我只按尺面写，不解释成"今天"。若编排者需要"10-01 11:00 那一次 AC#11 落盘"的严格时刻，**我这把尺取不到**。
- **E3｜第一次提交里我把"窗口是否 message-only"写成推断。** 现已升为现量：`HWND_MESSAGE` 在 `internal/ball` 零命中（R9b）＋
  `resident_ball_live_228_windows_test.go:69` 的 `FindWindowExW(0, prev, cls, 0)` 只能枚举**顶层窗口**这一事实＋`live_windows_test.go:512-515` 对它调 `SetWindowPos/GetWindowRect`。
  ⇒ 但"它的样式位里有什么（WS_EX_LAYERED/TOPMOST）"我**仍未逐字读 `createOnSTA`**（`ball_windows.go:182-283`）；本件结论不依赖它。
- **E4｜⛔ 第一次提交的 R11 有一段引用不成立——我已自查并改正。**
  首发里我写了 `internal/ball/interaction_live_test.go:7 "The two shapes this family CANNOT test:"` 与 `:12 "- the tray menu (needs a real mouse in the notification area; the unit"`，
  并把它当成"winlive 族自带承认测不了托盘菜单"的承重墙。**现量推翻**：
  该文件 `:11-15` 的真词面是"Each test drives the real window procedure through posted/sent Win32 messages or injected input..."，
  被点名为"没有 harness 能供的那只手"的是 **a second monitor** 与 **the felt experience of focus**（`:13-14`），**不是托盘菜单**；
  且 `notification area`、`tray menu` 在 `internal/ball` 目录**零命中**（R9b）。
  ⇒ 首发把 `//go:build winlive` 也写错了，真词面是 `//go:build windows && winlive`（R12）。
  ⇒ **这段错读的教训具名**：我把"我以为应当存在的自我承认"抄成了读数。本件的结论现在**只**建立在 R13 那种"pattern 在 `*_test.go` 零命中"的否定尺上，
  不再建立在任何注释文字上。**若编排者已经在上一发读过 R11，请以本条为准覆盖之。**
- **E5｜首发 R9 里我引了 `shutdown_hooks_test.go:64 TestShutdownHookSetCoversEveryNamedStep`。**
  **该函数不存在**（R/D-1 第 4 条给的真名册：`:31`、`:68`、`:120`）；`TestShutdownHookSetCoversEveryNamedStep` 这个名字**只出现在 `shutdown_hooks.go:108` 的注释里**。
  ⇒ 首发是从注释推的函数名＋从记忆推的行号，两枚都不该写。已改正。
- **E6｜"winlive 族今天零托盘覆盖"这句现在是强读数（否定尺），但仍有一条我没排除的可能**：
  某枚台件用**数字字面量**（`0x400+0x202`／`1026`）而不是符号名投递托盘消息，那样我的 pattern 抓不到。
  我没做"数值常量反查"这一形（成本／收益我不裁）。⇒ 若这条对裁定承重，值得让另一枚腿复核。
- **E7｜§② B-1 表里 `resident_approval_246_windows_test.go:132` 被我归入"真机进程族"是**不准确的**。**
  现量（R23）：`:132` 那枚是**进程内**（`proc.Boot` 直接调），`:184` 才是真机子进程。B-1 表与 B-4 结论都按这两枚分开算。
  （首发 E6 说"我没逐字读 246 与 127:565"，现已逐字读完，那条待核撤销。）
- **E8｜首发 E8（`recordTrayExit` 同包可调）仍成立**，但它成立与否**取决于修法**：若新位被做成 `Runtime` 的方法或需要参数，
  "直接调"这一支的形状就变。⇒ 留在 §⑦ F2，不作为已定案。
- **E9｜§② B-4 那句"没有任何一族能不加新东西完整覆盖"是本件的承重结论，它建立在两枚否定尺上**：
  R13（测试里零托盘符号）＋ R15（`cmd/wisp` 里零 `ldflags`/`H=windowsgui`）。
  两枚都是 pattern-零命中形，**天然可能被别的拼写绕过**（例如别的构建 helper 文件、别的 tag 名）。
  我搜的范围是 `cmd/wisp` 与全仓 `*_test.go`，**没做全仓 `-r` 的构建脚本普查**。
- **E10｜我对 `TrackPopupMenu` 阻塞行为的描述来自码形（同步取返回值）与注释，不是来自 MSDN 或实跑。**
  ⇒ 我在 A-1 第 2 支只说"它停在 `showMenu`、要拿 `sel==4` 得让菜单真被选中"，
  没有断言"必然挂死"或"能用 `WM_CANCELMODE` 之类干净退出"。⇒ 这一格留 §⑦ F1 甲支判不了的原因里。
- **E11｜本发自查出的两处引用级错误**（正文已改，记录在此以免被当成没发生过）：
  1. `tray_windows.go` 的 `uID`／`uFlags`／`uCallbackMessage` 三枚行号，我在 A-0 表与 R7 里各写过一种错配。
     真值最终由**两把独立的尺同向**才定：Grep（`-n`，pattern `uFlags|uCallbackMessage|uID:`）→ `35 / 36 / 37`；
     窄幅 Read（`offset=26 limit=16`）→ 同值。⇒ 这正是任务里"行号必须反取词面"那条纪律防的形。
  2. §⑦ F7 我最初把 `requireQuietBallDesktop` 引成 `live_guard_windows_test.go:11`。真凭据是 `resident_ball_live_228_windows_test.go:22-23`
     （那枚 header 在**指认别的文件**；`live_guard_windows_test.go:11` 的真词面是 "the stray balldebug.exe made TestBallLiveLifecycle go red for the wrong"）。
  ⇒ 其余未做双尺复核的行号仍可能有这类 off-by-one；**每处都附了尺命令，请以尺为准**。

---

## ⑦ 判不动的地方（逐条甲／乙／不做＋现量）

> ⛔ 只摆形状与现量，一律不推荐。

- **F1｜"一次托盘命令能否不靠真鼠标确定性地送达 `menuExit` 那一支"——判不了到底。**
  现量：`tray_windows.go:15-16`（不用 `WM_COMMAND`）、`:101-106`（`TrackPopupMenu` 同步返回）、
  `ball_windows.go:669-686`（分流）、`live_guard_windows_test.go:251-263`（`sendMouse`/`sendMessage` 形状）、
  `live_windows_test.go:553-562`（`postMsg` 形状）、`interaction_live_test.go:35` ＋ `live_guard_windows_test.go:274-290`（真输入注入先例）。
  - **甲**：同包 winlive 台件 `sendMessage(b, wmAppTray, 0, wmRButtonUp)`，再用真输入（键盘或鼠标注入）在弹出的菜单上选"退出"。
    判不了的原因：`TrackPopupMenu` 在返回前泵消息且阻塞调用线程（码形如此，见 E10）；
    弹出菜单能否被稳定聚焦/导航、`SetForegroundWindow`（`:97`）在无人值守会话下的行为、以及返回 0（被 dismiss）与返回 4 的区分——**我不跑就不敢定**。
  - **乙**：同包 winlive 台件跳过 `showMenu`，直接 `b.sta.PostTask(func(){ b.fire(b.opts.Events.OnTrayExit) })`。
    判不了的原因：它**不经过 `:670-675` 的分流**，"算不算 AC#11 ⓑ 要的仪器"是**裁定**问题，不是测量问题。
  - **丙**：为测试开一枚导出面（把"托盘命令"变成可投递的值，例如 `Ball` 上一枚 `ForTest` 方法）。
    判不了的原因：**是否为测试开导出面**是契约级决定；且它会与 `resident_ball_228_test.go:319-322` 的名册纪律相互牵连。
  - **不做**：只在 `cmd/wisp` 同包调 `recordTrayExit()`。现量：可达（R6）；但它连 `b.fire` 的线程语义都不走，是**最弱的一支**。
- **F2｜"修法长什么样决定哪一族能变红"——本普查拒绝猜修法。**
  现量：R3（全仓零 stop 位）＋ B-4。⇒ 完整覆盖必须**先有**可投递的停机请求；它的形状（位在 `Runtime` 上还是包级、谁读、`select` 几支）
  直接决定 §② 每一枚的可照性。**先有修法才谈得上仪器**。
- **F3｜空窗那一坑（`resident_windows.go:106-108` 与 `:194-203`）的最小可测形状——判不了。**
  现量：`boot_windows.go:127-152`（环路已每拍 50ms ⇒ 加一支 `case` 不需新协程）、`resident_windows.go:194-203`（已有 `select`＋`default`）、
  `resident_ball_228_windows_test.go:190-196`（时刻锚 `CountLogFiles`）。
  - **甲**：位＝`chan struct{}`，同时接进 `:196` 那支 `select` 与 `RunEventLoop` 的 `:145` `select`；测试沿用时刻锚后**关那枚 channel**。
    判不了的原因：要能在测试里关它，测试必须持有它 ⇒ 要么导出、要么在 `cmd/wisp` 同包构造 `runResident`；
    而后者与仓里既有判断冲突——`resident_ball_228_test.go:17-20` 逐字写了为什么不做进程内测（"runResident() never returns ... would mean a test reimplementing its surroundings"）。
  - **乙**：位＝原子 bool，由环路的 50ms 拍轮询；测试侧仍走真进程，但触发物换成别的可投递物（命名事件／管道／文件哨兵）。
    判不了的原因：三枚候选**每枚都是新契约**，且命名事件会撞 D38b 名册（R21：`boot_windows.go:115-118` 会因超基线**直接拒绝 Boot**）与 ban #1（R22）。
  - **不做**：只在环路内那一支装仪器，空窗那一支不装。判不了的原因：AC#11 把空窗写成"隐藏的坑"，**不装算不算缺陷由编排者判**。
- **F4｜要不要给台件产第二枚 `-H=windowsgui` 产物——判不了。**
  现量：`secret_argv_windows_test.go:174` 逐字无 `-ldflags`（R16）；`buildWispForTest` 被至少四枚用例共用
  （`resident_ball_228_windows_test.go:54`/`:185`、`resident_sink_nail_127_windows_test.go:566`、`resident_approval_246_windows_test.go:185`、`resident_ball_live_228_windows_test.go:90`）。
  判不了的原因：**共用者名册我只数到上面这五处，没做全仓反查**；改一枚共用 build helper 的产物形态会把这五处的读数一起洗掉
  （尤其 CUI→GUI 之后 stdout/stderr 的去向，而这一族**全靠 `leg.stdout` 的字面句子**断言，`resident_ball_228_windows_test.go:25-28` 逐字声明"三枚字面量是故意抄的"）。
- **F5｜"审计顺序仍是 1..10"这一维在真机侧要不要新做读面——判不了。**
  现量：D-3 全表。结构性障碍＝成功步不落日志（R19：只有 `shutdown.go:130/145/147/172` 四支出声，`:178-180` 与 `:186` 零日志）；
  `sinkInstallRecord.Step` 已在解析器里（R20 `:655`），而 `stepsOf127` 今天只进 `t.Logf`（`:640`）。
  - **甲**：改在**进程内**读满 10 条——照抄 `resident_approval_246_windows_test.go:145-158` 那枚形状（真 `proc.Boot`＋真 `RegisterShutdownHook`＋`rt.Shutdown(false)` 取 `[]StepRecord`）。
    判不了的原因：进程内测不到"托盘那一次点击"（它压根没有窗口腿），只能测"触发之后走满十步"——**这两维是否都要在同一枚仪器里，是裁定问题**。
  - **乙**：让 `internal/proc` 对**每一步**都落一行（把 `records` 变成盘面事实），再把 `:640` 那枚 `t.Logf` 升成断言。
    判不了的原因：这是在 D38(e) 的日志面上加东西，会同时改变 127 一族看到的记录数与顺序
    （`resident_sink_nail_127_windows_test.go:619-641` 那组"installIdx==1 / 最后一条必须是 shutdown 记录"的断言对**记录集合**敏感）。
  - **不做**：ⓑ 只钉"回调触发的那一支也走到了 D38(e) 尾巴"（沿用 `residentShutdownRecord` 存在＋退出码 0）。
    判不了的原因：够不够，由编排者按 AC#11 的措辞判。
- **F6｜若修法复用 `OnTrayExit` 而不新加 `Events` 键，则 §② 族② 一声不吭——这是可接受的还是缺陷，判不了。**
  现量：`resident_ball_228_test.go:319-322`（只有**多于**名册才红）、`:275-318`（只查键名与非空）。
  判不了的原因：那一枚的 AC#1 措辞是"每个 gesture callback 都要有执行者"，**不是**"每个执行者都要有行为"；把它扩到行为维是改判据。
- **F7｜winlive 两包不得共享桌面这一限制，会不会让新仪器变成"只能单跑"的台件——判不了。**
  现量：`resident_ball_live_228_windows_test.go:22-29` 逐字（"Run it alone... the two live suites may not share a desktop"，并在 `:23` 点名
  `live_guard_windows_test.go` 的 `requireQuietBallDesktop` 与 ticket 64 rule 1）。
  判不了的原因：CI 怎么排这两族是**构建/调度**问题，本腿没读构建脚本（禁令范围内我只读了 `tools/d22scan` 与 `cmd/wisp` 的台件）。

---

## 交件判语

**射程量到哪**：
1. A（§①）答满：链路 12 段全部带现读行号；"不靠真鼠标能否送达"给了**五支**分判，其中"投 `WM_COMMAND`＋id"这一支被**证伪**（`tray_windows.go:15-16`），
   "同包 sendMouse 进真 wndproc"这一支被**证实形状现成**（`live_guard_windows_test.go:251-263`）；`hWnd` 可达性两维（进程内 `DebugHWND:1003`、跨进程 `findBallWindows228:60-79`）都有凭据；
   窗口**不是** message-only 已从推断升为现量。
2. B（§②）答满并**多交一族**：任务给的三族之外，普查挖出第四枚成员 `resident_ball_live_228_windows_test.go:89`（真机进程 × 真窗句柄），
   它是现有码里离 AC#11 ⓑ 最近的一枚。具名结论：族②（AST）＝假绿那一族；**三族（实为四族）今天都不能"不加新东西"完整覆盖这一维**。
3. C（§③）答满：能制造"球已起、环路未进"的台件**全仓只有一枚**（`:184`），机制已逐行复述（时刻锚是 `observe.CountLogFiles` 这个磁盘事实，不是 sleep）；
   并指出**另有一格无人钉**（`resident_windows.go:103-105` 自承留下的"NotifyContext 之前那几条指令"缝）。
4. D（§④）答满：编排者"不会自动变红"的判断**证实**，凭据五条；同时交出**两枚会红的**（第七枚常驻协程 → `boot_windows.go:115-118` 拒绝 Boot；新增 `Events` 键 → `resident_ball_228_test.go:319-322`）
   与一枚射程差别（ban #1 不含 `_test.go`，ban #8 含）。"1..10 都在且顺序对"这句话**分进程内／真机两半答**，
   并交出本件最硬的一格：**真机侧这句话目前结构性读不到**——因为成功步不落日志（`shutdown.go:140-150`），不是缺台件。

**哪几格我没答（不写成答了）**：
- `TrackPopupMenu` 在无真指针下能否稳定返回 `4`：**没答**（E10、F1 甲）。
- 跨进程投递托盘消息需要新载 proc 之后**实际能否生效**（UIPI／会话桌面限制）：**没答**，我没有跑任何东西。
- `createOnSTA`（`ball_windows.go:182-283`）里的具体窗口样式位：**没逐字读**（E3；本件结论不依赖它）。
- `buildWispForTest` 的**完整共用者名册**：**没做全仓反查**（我只列了 5 处已见调用点，F4）。
- 用数值字面量（而非符号名）投托盘消息的既有台件：**没排除**（E6）。
- 票面 AC#11 的**其它判据支**（ⓐ／ⓒ…）：**不在本腿射程**，本腿只量"ⓑ 那枚会响的仪器"这一格。
- 我自己的**两处假读数**（R11 的 winlive 注释、`shutdown_hooks_test.go:64`）已在 §⑥ E4/E5 具名改正；
  若编排者已消费过第一发，**以本发的 §⑤⑥ 为准**。
