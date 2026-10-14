# 319-a1 · 面板线程死亡 / 二次启动唤窗 —— 根因调查与最小形（只读腿）

- 锚：`git rev-parse --short HEAD` = **`ccdceff5`**（2026-10-14 21:5x 现读）。
- 现场：`%APPDATA%\wisp-dev\logs\wisp-20261014-001.jsonl`（4417 行；本腿逐字现读，未改动）。
- 面：`cmd/wisp` 面板宿主 + `internal/proc`（唤窗那半）。
- 只读声明：零产码；`%APPDATA%` 与仓内生产面一字未动；未 kill / 未重启 / 未跑第二实例。

---

## 0. 结论（先行）

**AC#0 根因＝定案（不是"未定"）**：15:41:45 面板窗被**关闭**（`WM_SYSCOMMAND`→`DefWindowProc`→`WM_CLOSE`→`DestroyWindow`→`WM_DESTROY`→`PostQuitMessage(0)`），
库泵 `Run()` 见 `WM_QUIT` 即返回；线程 `loop` 因此 return，`finished` 关闭 ⇒ 之后每一次手势都被 `post()` 按"线程已死"拒。
**这不是崩溃、不是电源/会话事件**：Win32 消息序是"正常关窗"的完整档案，且 Application 事件日志与 WebView2 Crashpad **两侧都空**。

**AC#1 判定＝今天无任何复活路径**（`startResidentPanel` 全仓唯一构造点在 `cmd/wisp/resident_windows.go:198`，无第二处、无重建守卫）。
最小形＝把装配根手里的 `*residentPanel` 换成"持有者"，手势处 `isFinished()` ⇒ 重建（同一 `mgr`、同一 `panel-sta` 线程名）。**2 枚产码文件**。

**AC#2 判定＝唤窗路径今天确实不存在**，且 21:40:38 最新一发是在**线程活着、窗口活着**的情况下也不唤 ⇒ 两条缺陷可分离、都现症。
最小形＝把 `internal/proc` 那个激活分支接一个注入型 hook（同 `RegisterShutdownHook` 形状），指向**球/托盘已用的同一个面板手势执行器**；
外加第二实例侧 `AllowSetForegroundWindow(ASFW_ANY)` 一跳（否则 `SetForegroundWindow` 会被前台锁拒）。

---

## 1. AC#0 · `panel thread ending / why: the library pump returned` 为什么返回

### 1.1 双向锚（逐字）

**日志侧（同一 hwnd=0x9e9047e，窗建于 15:33:30.048，`probes` 之外的原始行）**

```
15:41:45.5657899  MARK299 wndproc hwnd=0x9e9047e msg=0x20  known=true   WM_SETCURSOR
15:41:45.5657899  MARK299 wndproc hwnd=0x9e9047e msg=0xa1  known=true   WM_NCLBUTTONDOWN
15:41:45.6853491  MARK299 wndproc hwnd=0x9e9047e msg=0x215 known=true   WM_CAPTURECHANGED
15:41:45.6859938  MARK299 wndproc hwnd=0x9e9047e msg=0x112 known=true   WM_SYSCOMMAND
15:41:45.6859938  MARK299 wndproc hwnd=0x9e9047e msg=0x10  known=true   WM_CLOSE      <-- 同毫秒
15:41:45.6859938  MARK299 wndproc hwnd=0x9e9047e msg=0x90  known=true   WM_UAHDESTROYWINDOW
15:41:45.6865..6929  msg=0x46 0x47 0x86 0x6 0x1c 0x8 0x281 0x282 0x272   拆窗序列
15:41:45.6940082  MARK299 wndproc hwnd=0x9e9047e msg=0x2   known=true   WM_DESTROY
15:41:45.736702   {"msg":"panel thread ending","why":"the library pump returned","window_opened":true}
```

**代码侧（哪一行）**

| 步 | 逐字位置 |
|---|---|
| `DefWindowProc(WM_SYSCOMMAND)` 发出 `WM_CLOSE` | 本仓无任何代码 post `WM_SYSCOMMAND`/`SC_CLOSE`（`grep -rn "WM_SYSCOMMAND\|SC_CLOSE\|0x112" --include=*.go` 射程含 go-webview2 依赖 ⇒ **0 命中**）⇒ 只能来自系统 |
| `WM_CLOSE` → `DestroyWindow` | `cmd/wisp/panel_control_windows.go:266-269`（`case ctlWmClose:` → `ctlDestroyWindow.Call(hwnd)`） |
| `WM_DESTROY` → `PostQuitMessage(0)` | 同文件 `:270-275`（`case ctlWmDestroy:` → `ctlDropWindow` + `ctlPostQuitMessage.Call(0)`） |
| 泵见 `WM_QUIT` 即返回 | 同文件 `:559-573`（`Run()`：`else if msg.message == ctlWmQuit { return }`） |
| 线程随之收尾 | `cmd/wisp/panel_resident_windows.go:319-320`（`w.Run()` → `rp.teardown("the library pump returned")`） |
| 此后永久拒单 | 同文件 `:280`（`defer close(rp.finished)`）、`:373`（`post()` 首查 `startUpErr() != nil || isFinished()`） |

移植忠实性已逐枚对过上游：`webview.go:242-243`（`WMDestroy` → `w.Terminate()`）、`:347-349`（`Destroy` = `PostMessage(WMClose)`）、`:351-379`（`Run` 见 `WMQuit` 返回）、`:381-383`（`Terminate` = `PostQuitMessage`）。
⇒ **这条链是库合约定死的，不是 299 移植引入的偏差。**

### 1.2 三假说判读

- **甲（关窗 ⇒ 泵返回）＝成立，且是唯一被证据支持的一条。**
  判据①：`0x112` 之后**同毫秒**出现 `0x10`，而本仓/依赖无一处 post `WM_SYSCOMMAND` ⇒ `DefWindowProc` 把该 `SC_` 码翻成了 `WM_CLOSE`（能翻成 `WM_CLOSE` 的只有 `SC_CLOSE` 一族：X 按钮 / Alt+F4 / 系统菜单"关闭" / 标题图标双击 / 任务栏关闭）。
  判据②：同一份日志里 `0x112` 出现过 **12 次**（15:33:32/36/39/44、15:37:04、21:03:31、21:04:15 等），只有 15:41:45 这一次**后面接 `WM_CLOSE`** ⇒ 其余是 `SC_MOVE`/`SC_RESTORE`（15:33:32 那发后面是 `0x215→0x24→0x231(ENTERSIZEMOVE)`＝拖标题栏；21:03:31 那发后面是 `0x13(WM_QUERYOPEN)`+`WM_SIZE` 簇＝还原/最大化）。
  判据③：紧邻的 `0xa1`(按下) → 0.119 s → `0x215`(capture 释放) → `0x112`，与"点非客户区按钮"的按下/抬起形状一致。
  判据④（自洽性）：此后 9 发手势全被拒、进程仍活 ⇒ 正是"窗没了 + 线程没了"的读法。
- **乙（WebView2 崩 / 浏览器进程退出）＝反证。** ①消息序是**程序化销毁**（`WM_CLOSE` 起头），崩溃不会发 `WM_CLOSE`；
  ②`wevtutil`/`Get-WinEvent` Application 日志 15:20–16:00 **零条**（`probes/319/a1/evt-app-1541.ps1` 可重跑；相邻最近两条是 15:21:42 与 15:47:01 的 Security-SPP），**无 Application Error / Windows Error Reporting**；
  ③`%APPDATA%\wisp-dev\panel-webview2\EBWebView\Crashpad\reports\` **空**（该目录自 10-09 起从未有过 dump）。
- **丙（电源/会话事件）＝反证。** 序里无 `WM_QUERYENDSESSION(0x11)` / `WM_ENDSESSION(0x16)`；进程在 15:41 之后仍持续记日志到 21:40。

### 1.3 AC#0 的取证缺口（建议登记，本腿不改）

`MARK299` 那行只打 `hwnd/msg/known`，**不打 wParam** ⇒ `SC_` 码（`SC_CLOSE` vs `SC_MOVE`）只能靠邻居消息反推，不能逐字复读。
建议：给该行加 `wparam`（**这是产码改动，本腿只提形**）。同一条建议也覆盖"当日 12 发 `0x112` 无一处留 wParam"这一事实。

---

## 2. AC#1 · 线程退出后今天能不能重建

- **不能，且是"三处都缺"**：
  1. 构造点唯一：`startResidentPanel` 定义 `panel_resident_windows.go:144`，**唯一生产调用** `resident_windows.go:198`（`panel = startResidentPanel(rt.Registry, rp)` + `:199 defer panel.stop()`）。全仓无第二处、无 watchdog 重拉（`internal/observe` 只有名册统计，**无重启**）。
  2. 判死是终态：`post()`（`:369-392`）首句 `if rp.startUpErr() != nil || rp.isFinished() { return false }` ⇒ 线程 goroutine 一 return，**永远 false**，这正是日志那句 `why` 的出处（`resident_ball_windows.go:122`）。
  3. `Manager` 侧没有可复用的"窗口"概念：`teardown` 会 `mgr.Destroy()` ⇒ `created=false, w=nil, hwnd=0`；`bringUp` 可以再建（`newPanelControl` 重注册同类名，Win32 答"已存在"、首次注册继续服务，见 `panel_control_windows.go:362-363`），**缺的只是"谁来再叫一次线程"**。
- **顺带一枚同族缺陷（比 X 按钮更硬）**：`RequestDispose` 的注释说"keeps the pump running, so a later RequestShow builds a fresh window"（`panel_resident_windows.go:459-461`），**实测不成立** —— `Destroy()` 走 `PostMessageW(hwnd, WM_CLOSE)`（`panel_control_windows.go:605-608`），与用户按 X **同一条路**：`WM_CLOSE→DestroyWindow→WM_DESTROY→PostQuitMessage` ⇒ 泵返回、线程死。
  ⇒ 即"自带的重建路"（dispose+show）在今天的移植下**也是死的**（生产零调用者，故从未暴露；测试用的是自己的泵 `pnlPumpOnce`，见 `panel_host_windows_test.go:67-101`，盖不到这条）。
- **机主今天自己能走的恢复路径（AC#1 要求具名）**：只有"杀进程再拉起"（编排者 20:30/20:35 就是这么救的）。**托盘 Exit 不可用**（15:41:51 / 15:42:00 逐字 `tray exit requested / outcome: ignored`），故**托盘上没有任何停机面**，只能任务管理器结束 `wisp.exe` 或 Ctrl+C（有控制台时）。

### 2.1 最小形（逐枚点名；⛔ 本腿不改码）

**形甲（建议主形）＝手势处惰性重建**，产码面 **2 枚文件**：

- `cmd/wisp/resident_windows.go`
  - 把 `panel` 换成持有者（例：`type panelHolder struct { mu sync.Mutex; reg *observe.Registry; mgr *PanelManager; rp *residentPanel }`，方法 `toggle(via string) bool` / `statusLine() string` / `stop()`）；
  - `toggle`：`rp == nil || rp.isFinished()` ⇒ `rp = startResidentPanel(reg, mgr)`（**同一 `panel-sta` 名、同一 owner、经 `Registry.Spawn`，不新增名册项**——D38b 名册冻结不动，`panel_resident_windows.go:144-156` 原样复用）；
  - `:199 defer panel.stop()` ⇒ `defer holder.stop()`（否则退出时停的是**旧**那枚）；`:257` 的闭包 ⇒ `holder.toggle`；`:388` 的 `panel.statusLine()` ⇒ `holder.statusLine()`。
- `cmd/wisp/panel_resident_windows.go`
  - 只加 Read 级谓词即可（`isFinished()` 已存在，可直接用）；若要钉死"线程是不是被 `errPanelRefusedThread` 退休"可加 `dead() (bool, error)` 供日志说真话。**不必改 `loop`/`post` 的路由**。

**形乙（可选渐小改，只治 X 按钮那一半）＝关窗改"藏"**：`panel_control_windows.go` 的 `ctlWmClose`（`:266-269`）是库逐字移植，要动就按该文件已有的"具名偏差"体例登记（文件头已列 2 条偏差），并且必须区分"用户关"与"宿主 dispose"（否则 `teardown`/`RequestDispose` 也被藏住）⇒ 需要一枚宿主侧意图标志。**收益**：把"线程会死"的频率压到零；**局限**：CoInitializeEx 拒绝、`bringUp` 具名拒绝（33-r9 的退休路）仍会让线程永久死 ⇒ **甲形无论如何都要**。

**不变量检查**：形甲不新增线程名/不新增依赖边/不改 D38(e) 十步（新线程只在**旧线程已 return** 后才 Spawn，名册 live 计数 0→1，`ResidentOverBaseline` 不动）；新 goroutine 经 `LockOSThread` 拿到的是**全新 OS 线程**（Go 对 lock 后 return 的线程是销毁，`panel_host_windows_test.go:90-101` 已把这条写进注释），故 `bringUp` 的前置条件"本线程无可派发消息"成立。

---

## 3. AC#2 · 二次启动为什么唤不醒

- **现读路径**：`cmd/wisp/resident_windows.go:77-84`（`proc.ErrAlreadyRunning` ⇒ `proc.SignalExistingInstance` ⇒ 打印"activated it"⇒ 退出）
  → `internal/proc/singleinstance_windows.go:127-141`（`SetEvent`）
  → `internal/proc/boot_windows.go:127-146`（`RunEventLoop` 每 50 ms `WaitForSingleObject`；命中即**只记账+清除**，并逐字打印 `bring_to_front: no path from this loop to a window; the ball, when a process has one, is its host's`）。
  ⇒ **这条边两头都没接**：回执没有可调用的窗口句柄，`internal/proc` 是叶子包（不知道球/面板存在）。
- **"代码自注 bring-to-front 落它名下"＝外证**：`internal/proc/singleinstance_windows.go:17`「bringing the ball to the front is the ball's job, **ticket 07**」＋`:84`「ActivateEvent ... (ball, ticket 07)」。
  而 `internal/ball/doc.go:24` 的 ticket 07＝「static states + basic animations」，`anim.go:5`/`tokens.go:10` 同族 —— **07 名下没有 bring-to-front 的实现**。
  运行时那句带 `(ball bring-to-front lands with ticket 07)` 的旧文案已在 `15d2b2cd`（228-r1 AC#7）改成现句；**注释里的旧归属没跟着改**（本腿只报）。
- **今天最新的活体读数（关键）**：日志最后一行 `21:40:38.9085` 又是一发 `activation requested by second launch`，而此时**面板线程活着、窗口活着**（hwnd 0x8fd0784 建于 21:00:11，无 `panel thread ending`）⇒ **AC#2 与 AC#1 是两条独立缺陷，双击没反应不能全赖线程已死**。当日累计：`activation requested by second launch` **12 发**、`panel gesture could not be handed over` **9 发**、`panel thread ending` **1 发**。

### 3.1 最小形（逐枚点名；⛔ 本腿不改码）

1. **接一条桥（形甲）**：`internal/proc/boot_windows.go`
   - 加注入位（形状抄本文件已有的 `RegisterShutdownHook`，**func 值、不新增依赖边**）：`func (rt *Runtime) SetActivationHook(f func() string)`；
   - 把 `:133-142` 那一段抽成可测方法 `func (rt *Runtime) handleActivation() bool`：命中挂钩就跑一次并把它的**回执逐字入日志**（谁被要求举到前台、Win32 答了什么），随后 `ResetActivation()`；没挂钩时保留今天那句实话，**但不再说"no path"**。
   - ⚠ 不碰 D38(e) 十步，不新增 goroutine（轮询仍在事件循环里）。
2. **接线**：`cmd/wisp/resident_windows.go`
   - `rt.SetActivationHook(func() string { return holder.toggleReport("second-launch") })` —— **复用球/托盘同一条执行器**（`:257` 那枚 `panel.RequestToggle` 就是；`Manager.Show` 里已有 `pnlShowWindow(SW_SHOW)` + `pnlSetForeground`，`panel_host_windows.go:518-556`）⇒ "唤窗"不是新能力，只是**没接**。
3. **前台权限那一跳（否则 SetForegroundWindow 会被前台锁拒）**：第二实例在 `SetEvent` **之前**调 `AllowSetForegroundWindow(ASFW_ANY)`。
   落点二选一：`internal/proc/singleinstance_windows.go`（与 `SignalExistingInstance` 同函数体，最少面）或 `cmd/wisp/resident_windows.go:77-84` 那个 `ErrAlreadyRunning` 分支。
   **判据**：把 `SetForegroundWindow` 的返回值继续"只记不重试"（`panel_host_windows.go:686-690` 的既有姿态），还有前台锁就可能失败——**不许用重试把它拧绿**。
4. **回退可见提示**（当面板起不来 / 线程退休 / WebView2 缺）：现有 API 里**够用的最小一枚**是 `internal/ball` 的 `SetTrayTip`（`ball_windows.go:961`，零新 API、加性）；更强的托盘气泡（`Shell_NotifyIcon` + `NIF_INFO`）需要给 `internal/ball` 加新方法 ⇒ **另立腿/另立票**，不要顺手塞进本票。
5. **⛔ 面归外部 agent**：任何"页面内闪烁/高亮"要在 `frontend/**` 落地 ⇒ **面归外部 agent**，本编队零写面。⛔ 本票不改 `frontend/**`。

---

## 4. AC#3 钉子建议（常驻判据 + 弄坏⇒红）

- **N1（AC#1，headless）**：持有者接受**可注入的 spawn 函数**；测试塞一枚 `finished` 已关的假 `residentPanel`，断言 `toggle()` **返回 true 且 spawn 被调第 2 次**。
  变异自证：把 `isFinished()` 那半删掉 ⇒ 必须红（不许靠真窗）。
- **N2（AC#2，headless）**：`handleActivation()` 单测——`SignalExistingInstance` → 跑一轮 → 断言注入 hook 恰被调 1 次（via 逐字 `second-launch`）、`ResetActivation` 生效（第二轮零调用）。
  变异自证：去掉 hook 调用 ⇒ 红。
- **N3（AC#0 现象级钉，winlive 门内）**：真窗建好后 `PostMessageW(hwnd, WM_SYSCOMMAND, SC_CLOSE, 0)` ⇒ 断言 ①线程 `finished` 关闭 ②随后一次手势**重建出窗**。这是 AC#1 票面"真机重放该形"的正面回答；headless 部分只能退到 N1。
- 现有相关钉：`TestPanelThreadIsSTAAndExitsCleanly`（`panel_resident_windows_test.go:958-989`）只覆盖 `stop()` 那条出口，**不覆盖泵自发返回**——N1/N3 是补它。

## 5. 具名没做

- 未 kill / 未重启 / 未跑第二实例 / 未触碰 PID 40296 与任何 `%APPDATA%` 文件（只读现读）。
- 未改任何生产码、未动 `Status:`/票名/台账/AC 框、未碰 `frontend/**`、`design/**`。
- 未逐枚验证 `MARK299` 里非关键消息号的 winuser.h 值（0x20/0xa0/0x215/0x112/0x90/0x231/0x13/0x83/0x85/0x14 等按标准表注；**载重的 0x10/0x2/0xa1/0x216 已在库 `internal/w32/w32.go:87-96` 与本仓 `panel_control_windows.go:119-127` 两侧对过**）。
- 未取证 `SC_` 码本身（日志无 wParam）；未做 15:41 那一刻的进程级抓包/句柄快照（现场已不可复现）。
- 未评估"重建线程 vs 保活窗口"两条最小形的性能/时序成本（需真机发）；未验证 `AllowSetForegroundWindow` 在本机的实际授予效果（需真机发，且不得扰动机主在跑的实例）。

## 6. 可重跑尺

```bash
# 锚
git rev-parse --short HEAD
# 关窗那一串（逐字）
grep -n "msg=0x10 \|msg=0x112\|msg=0x2 \|panel thread ending" "$APPDATA/wisp-dev/logs/wisp-20261014-001.jsonl" | sed -n '1,40p'
# 三计数
for p in "panel thread ending" "panel gesture could not be handed over" "activation requested by second launch"; do
  printf '%s = %s\n' "$p" "$(grep -c "$p" "$APPDATA/wisp-dev/logs/wisp-20261014-001.jsonl")"
done
# 事件日志（零命中＝乙被反证）
powershell -NoProfile -ExecutionPolicy Bypass -File .scratch/wisp/probes/319/a1/evt-app-1541.ps1
# Crashpad 空
find "$APPDATA/wisp-dev/panel-webview2/EBWebView/Crashpad/reports" -type f | wc -l
```
