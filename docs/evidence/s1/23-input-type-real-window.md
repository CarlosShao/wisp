# 23 — AC#3 未勾那格：`input.type` 真机 300 字 / 提权那一形（写腿读数）

**日期**：2026-11-14 10:24–10:36 +08
**锚**：HEAD `7bdf09042482b582a98071143f5cf6ec75512407`（`dev`）
**票**：`.scratch/wisp/issues/23-system-window-input-tools.md` AC#3（唯一未勾那格，逐字见 §1）
**写面**：`.scratch/wisp/probes/23/**`（新建）＋ 本文件 ＋ 票 23 Progress log 追加。⛔ 未动 `cmd/**`、`internal/**`、`scripts/**`、`frontend/**`、`models/**`、阈值/golden/契约；⛔ 未翻任何 AC 框；⛔ 未 push。

## §0 一句话结论（三个分句，逐句有读数）

| 票面分句 | 本次真机读数 | 判 |
|---|---|---|
| `input.type` 往 fixture app 聚焦窗口打 300 字 | runA（ASCII）：注入 300 字 → 窗内 **300 字**、`wm_char=300`、读取 sha256 与期望**逐字相同**；runB（CJK）同形同数 | **成立（真机实测）** |
| elevated-process injection → explicit UIPI error | 对**真 High 完整性窗口**（Task Manager，IL 12288）调用 `input.type`：**没有** `ErrElevatedTarget`，回的是 `IsError=false` ＋「已向 taskmgr.exe 注入 15 个字符」 | **不成立（产品缺陷，已定位到 `file:line` 并可复现）** |
| 确认卡三选项 render（「所有应用」danger 样式） | 生产前端源码三档文案 **0 命中**；只有设计稿 `design/old/screens/approval.html:848/:856/:871` | 本腿不可判（`frontend/**` 零写面），仍是裁决表 R3 |

⇒ 本格**不能算满足**：第一分句已量到，第二分句在真机上**失败**（不是"没测"，是"测了，产品没按契约拒"），第三分句归 frontend。**AC 框保持未勾**（翻框须非实现者）。

## §1 票面那格逐字

```
- [ ] input.type: types 300-char text into fixture app focused window; elevated-process
      injection → explicit UIPI error; confirm-card three options render with danger styling
      on "all apps".
```

契约出处：`docs/PLAN.md:3234` 第 2 行逐字「**Windows UIPI 平台硬限制**：非提权进程不能向更高完整性级别窗口发送输入。**无解**」＋「必须在文档与错误文案明示（`permission_denied` +「目标程序以管理员身份运行，Wisp 无法向其输入」）。**不得尝试提权绕过**」。

## §2 仪器（三个 probe，全部只读产品导出面；不重造轮子）

产物在 `.scratch/wisp/probes/23/`，源码三枚、日志在 `logs/`：

1. **`fixture23/`** ＝ 票面允许的 "scripted window" 目标程序（票面原文 "a test target app (notepad/scripted window)"）。独立进程 `wisp23fixture.exe`：Win32 顶层窗 `WISP23-FIXTURE input.type target` ＋ 多行 EDIT 子控件；**从自己进程内**用 `GetWindowTextLengthW/GetWindowTextW` 报回缓冲长度、sha256、`WM_CHAR`/`WM_KEYDOWN` 计数（子类化 EDIT 计数），并自报 `is_foreground`/`focus_hwnd`。50ms 定时落 `dump.txt`，`cmd.txt` 收到 `clear` 时清空缓冲与计数。
2. **`probe23/`** ＝ 读数腿。**不重写注入**：取 `tools.BuiltinSystemEntries(tools.SystemDeps{})`（`SystemDeps{}` → `bridge()` 造零值 `SystemBridge`，`Backend==nil` → `platformSystemBackend()`，与生产 `cmd/wisp/run.go:897` 的 `&tools.SystemBridge{}` 同形），调真正的 `window.focus` / `system.get` / `input.type` 三个 tool 的 `Execute`。注入前**必须**同时满足 `foreground==fixture 顶层窗` 且 `focus==EDIT 控件`（`GetGUIThreadInfo` 读键盘焦点），否则**一个字都不打**（本次 run2–run8 就是这么退出的）。
3. **`diag23/` · `uipi23/`** ＝ 只读谓词/平台读数（§4、§5）。

## §3 300 字那一半：逐字读数（run1，`logs/run1/summary.json`）

环境：注入进程 `probe23.exe` pid 58056，**IL `0x2000`（8192 = Medium）**；`roster_entries=11`（D34 十一枚）。

前置（打之前，两条独立读数都指向同一窗）：

| 读数 | 值 |
|---|---|
| fixture pid / 顶层 hwnd / EDIT hwnd | 29516 / `0x82e1476` = 137237622 / `0x6f10fba` = 116461498 |
| fixture 进程 / IL | `wisp23fixture.exe` / `0x2000` |
| 顶层窗标题 / 类 | `WISP23-FIXTURE input.type target` / `Wisp23FixtureWindow` |
| fixture 自报前台 | `foreground_method=already-foreground`，`is_foreground=1` |
| probe 侧 `GetForegroundWindow` / `GetGUIThreadInfo.hwndFocus` | `0x82e1476` / `0x6f10fba` → `foreground_is_fixture=true`、`focus_is_edit=true` |
| **产品侧** `system.get` | `focused_hwnd=137237622`，`focused_process="wisp23fixture.exe"`，`focused_title="WISP23-FIXTURE input.type target"` |
| **产品侧** `window.focus` | `IsError=false`，text「已切到窗口：WISP23-FIXTURE input.type target」 |

对照（证明原始 SendInput 仪器本身活着）：focus 在 fixture 时直接 `SendInput` 一个 `Q` → 落窗 `chars=1`、`wm_char=1`。

**runA — 300 个 ASCII 字符（`"0123456789abcdef"` ×18 ＋ `"0123456789AB"`）**

| 量 | 期望 | 实测 |
|---|---|---|
| 字符数 / UTF-16 单元 | 300 / 300 | 300 / 300 |
| sha256 | `1f47881ec9b983bda34c61a1a419b2483bd0091dc3c8f5a994d65b1c792fedbf` | `1f47881e…fbf`（**逐字相同**） |
| 产品回话 | — | `IsError=false`，「已向 wisp23fixture.exe 注入 300 个字符」 |
| 注入耗时 | — | 1632 ms（600 个 SendInput 事件） |
| 窗内 `WM_CHAR` / `WM_KEYDOWN` | — | 300 / 300 |
| 清空后起点 | 0 | 0（先 `clear` 再打） |

**runB — 300 个 CJK 字符（U+4E00 起连续 300 枚，均 BMP＝1 UTF-16 单元）**

| 量 | 期望 | 实测 |
|---|---|---|
| 字符数 / UTF-16 单元 | 300 / 300 | 300 / 300 |
| sha256 | `84056f1d0140b27916fd6811b594712249973a76881169a15473942946910e0c` | 相同 |
| 产品回话 | — | `IsError=false`，「已向 wisp23fixture.exe 注入 300 个字符」 |
| 注入耗时 | — | 1598 ms |
| 窗内 `WM_CHAR` / `WM_KEYDOWN` | — | 300 / 300 |

⇒ 票面第一分句**量到并逐字对上**：300 字（ASCII 与 CJK 各一轮）真的落进聚焦窗，不是"回了成功但看不到"。

## §4 提权那一形：真机上**没有** UIPI 错误（本格的真结论）

### §4.1 目标不是假的：真 High 完整性窗口

`EnumWindows` 扫到 14 枚可见有标题顶层窗，其中 **2 枚高于本进程 IL(8192)**：

```
hwnd=0x1706be pid=6904  name=taskmgr.exe   il=0x3000(12288) title="Task Manager"
hwnd=0x10346 pid=12780  name=cua-driver.exe il=0x3000(12288) title="Cua.AgentCursorOverlay.default"
```

用 `window.focus` 把它切到前台（`foreground=0x1706be`）后调：

```
input.type  {"text":"wisp-uipi-probe","hwnd":1509054}
→ IsError=false
→ text = 已向 taskmgr.exe 注入 15 个字符
```

**契约要的是**（`internal/tools/system_bridge.go:163` 的冻结句）：

> 目标程序以管理员身份运行，Wisp 无法向其输入（Windows UIPI 平台硬限制：普通权限进程不能向更高完整性级别的窗口发送输入，本工具不会尝试提权绕过）。请把目标应用以普通权限运行后重试

实测**一个字都没出现**；该次调用期间 fixture（另一进程）字符数不变（`elevated_fixture_chars_unchanged=true`）。也就是说：用户会读到"注入成功"，而更高完整性目标**静默收不到**——正是 `PLAN.md:3234` 明令不许发生的形态。

### §4.2 根因（只读定位到行，附可复现读数）

`internal/tools/system_backend_windows.go:728-742`：

```go
h, err := windows.OpenProcess(processQueryLimitedInfo, false, pid)   // :733 process handle
...
defer windows.CloseHandle(h)                                        // :739
return integrityOf(windows.Token(h)) > self                          // :741 ← 把"进程句柄"当"令牌句柄"
```

`integrityOf`（`:757-776`）第一条就是 `windows.GetTokenInformation(tok, …)`，而该 API 只接受 **token** 句柄；进程句柄进去必然失败 → 返回 `0` → `0 > self` 恒假。于是 `realProcessElevated` 对**任何能 OpenProcess 的 pid 都答 false**，只有 OpenProcess 自己失败那条分支才答 true。

`diag23` 同 run 并排量（`logs/diag23-run1.txt`）：

```
self pid=48572 integrity=0x2000
pid=6904  Taskmgr.exe   shipped_integrity=0 shipped_err="The handle is invalid." corrected_integrity=0x3000 corrected_sid=S-1-16-12288 → shipped_says_elevated=false corrected_says_elevated=true
pid=12780 cua-driver.exe shipped_integrity=0 …                                  corrected_integrity=0x3000 …                        → false / true
pid=69160 pwsh.exe      shipped_integrity=0 …                                  corrected_integrity=0x2000 …                        → false / false
```

`uipi23` 同形复读（`logs/uipi23-taskmgr-run1.txt`）：`target_integrity_shipped_shape=0 (error=sizing call returned no size)`、`target_integrity_corrected=12288`。

### §4.3 为什么门禁是绿的（诚实交代）

`internal/tools/ticket23_input_uipi_windows_test.go`：

- `:66-72` 的拒绝腿把 `processElevated` **打桩成 `func(uint32) bool { return true }`** ⇒ 断的是"拒绝序"，不是谓词；
- `:129` 用真谓词的那条腿只测"本进程自己的窗口放行"——那格**正确答案就是 false**，谓词恒假时它照样绿。

⇒ 裁决表 `docs/evidence/s1/23-verdict.md` §7「AC#3 已闭：UIPI 拒绝序」这句在**真机**上不成立；本条是本腿新量到的、需要写腿修 `internal/tools/system_backend_windows.go:741`（把 `integrityOf(windows.Token(h))` 改成先 `OpenProcessToken(h, TOKEN_QUERY, &tok)` 再 `integrityOf(tok)`）。**本腿写面禁令不含 `internal/**`，未改一个字节。**

## §5 平台那一半：UIPI 在这台机器上确实是硬的（可复读）

`uipi23 1509054`（`logs/uipi23-taskmgr-run1.txt`）：

| 读数 | 值 | 含义 |
|---|---|---|
| 注入方 / 目标 IL | 8192 (Medium) / 12288 (High, `S-1-16-12288`) | 真的跨边界 |
| `AttachThreadInput(myThread, taskmgrThread=26548, TRUE)` | **0，`Access is denied.`** | UIPI 对"跨进程输入耦合"的正面拒绝 |
| `SetForegroundWindow(taskmgr)` | ok（本来就是前台） | — |
| `SendInput` 一个 `Z` | **inserted=2/2，err=nil** | **插入 ≠ 送达**：计数器只数入队 |
| 目标线程键盘焦点（前后） | `0x9fc058e class="DirectUIHWND"`（列表区，非文本框） | 见下方诚实边界 |

诚实边界（不许夸大）：目标窗内焦点落在 **DirectUIHWND 列表**上（不是搜索框），所以"屏幕上看不到 `Z`/`wisp-uipi-probe`"这条**单独**不能证明"被 UIPI 丢弃"——非文本控件本来也不显示字符。本腿能确证的只有：①`AttachThreadInput` 被拒（`Access is denied.`）；②`SendInput` 的返回值只表示入队。若要正面量"插入成功但送达被丢"，需要可控的 High 目标（＝起提权进程＝UAC/机主在场）或可控的 Low 注入方（起低完整性进程），**本腿不做**。

留档截图：`logs/screen-after-elevated-attempt.png`（Task Manager 搜索框仍为空；同上，仅作旁证）。

## §6 第三分句（确认卡三选项 + danger 样式）：本腿不可判

只读普查（`grep` 三档文案）：

- 生产 `frontend/**`：**0 命中**（`仅本次` / `本会话内允许向此应用` / `所有应用`）。
- 只有设计稿 `design/old/screens/approval.html:848`「仅本次」、`:856`「本会话内允许向 WINWORD.EXE 注入」、`:871`「本会话内允许向所有应用注入」＋ `color:var(--danger)`；规范 `docs/specs/SPEC-06-security-gatekeeping.md:117/:120`（第三条须"视觉降级＋二次点击展开"）。
- `cmd/wisp/approval_reply.go:239/:282/:655` 是**答复动词**（yes/session/no/veto），不是卡片渲染。

⇒ 与裁决表 `docs/evidence/s1/23-verdict.md` §8 的 R3 一致：卡片渲染归 `frontend/**`（机主外部 agent），本编队零写面。`internal/tools/system.go:63-67` 也逐字写着"AC3 confirm-card UI NOT COVERED here"。

## §7 重跑录像（run2–run8 为什么没数）

run1 之后我试着再跑几轮，**全部按安全闸退出、一个字没打**（每轮 `logs/runN/probe.log` 末行 `ABORTED=focus could not be placed…`）：

- run2：前台是 Task Manager（elevated），`SetForegroundWindow` / `SwitchToThisWindow` / `AttachThreadInput` 全失败 → 拒绝打字；
- run3/run4：前台是另一会话的模态对话框（`#32770`「复习结束 -> 听力练习」，属别的 pwsh，未动它）→ 同样抢不到前台；
- run6–run8：一次 `raise+click` 兜底把鼠标点到 Task Manager 上，反倒把它固定成前台；此后 `ShowWindow(SW_MINIMIZE)`（跨进程）与 Shell 的 `ToggleDesktop` 对 elevated 窗都**无效**（UIPI 挡消息），前台回不到普通窗口。

⇒ **要再产一轮 300 字读数，需要机主把前台让回普通窗口**（点一下任意普通窗口，或关掉/最小化任务管理器）。这一步本腿**不代做**（关别人的 Task Manager 属机主动作）；run1 的读数已足够判这一格，不作弊补跑。

## §8 判语与欠账

1. **AC#3 不满足**，且失败点是产品侧、不是仪器侧：`input.type` 对真 High 窗口没有显式 UIPI 错误。修法见 §4.2（一行句柄语义），**归写腿**（`internal/tools/system_backend_windows.go:741`），修完后本文件 §4 的同一发可用来复核（`probe23` 原样复跑即可，无需 UAC）。
2. 300 字那半已闭合（ASCII＋CJK 各 300/300，sha 逐字相同）。
3. 卡片三选项仍归 frontend（R3）。
4. 复跑命令：
   ```powershell
   cd <repo>
   go build -o .scratch\wisp\probes\23\logs\wisp23fixture.exe ./.scratch/wisp/probes/23/fixture23
   go build -o .scratch\wisp\probes\23\logs\probe23.exe      ./.scratch/wisp/probes/23/probe23
   go build -o .scratch\wisp\probes\23\logs\diag23.exe       ./.scratch/wisp/probes/23/diag23
   go build -o .scratch\wisp\probes\23\logs\uipi23.exe       ./.scratch/wisp/probes/23/uipi23
   .\.scratch\wisp\probes\23\logs\probe23.exe .\.scratch\wisp\probes\23\logs\wisp23fixture.exe .scratch\wisp\probes\23\logs\runX
   .\.scratch\wisp\probes\23\logs\diag23.exe <taskmgr-pid>
   .\.scratch\wisp\probes\23\logs\uipi23.exe <taskmgr-hwnd>
   ```
5. 产物说明：`logs/run1/fixture/dump.txt` 的**最后一行**是 `chars=0 … wm_char=300`——这是 fixture 收到 `WM_CLOSE`、EDIT 控件先被销毁、守护线程在退出前才落最后一发造成的；中段读数（`logs/run1/probe.log:10:24:38` 那几行）才是注入后的实测值，两者不矛盾（同因也解释了 `text.txt` 为 0 字节）。该 artifact 记在此处，免得后人误读。
6. 副作用（如实记）：起过 8 次 fixture 窗（全部关闭/杀干净）；对 Task Manager 发起过两次输入尝试（15 字与 1 个 `Z`），屏上无变化；**未尝试提权、未弹 UAC、未创建任何提权进程**；未改任何系统设置；鼠标位置在 click 兜底里已存/已复原。

## §9 写面与撤销

- 新增：`.scratch/wisp/probes/23/**`（源码 4 枚、日志与读数）＋ 本文件 ＋ 票 23 Progress log 一段。
- 对产品代码、契约、阈值、golden：**零字节**。
- 撤销口令：`git revert <本提交>` 一枚即整体复原；单删 `docs/evidence/s1/23-input-type-real-window.md` 亦不影响仓库其它部分。
