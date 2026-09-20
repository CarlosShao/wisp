# 64 — 票 07 遗留缺陷修复：热键接线 / 交互四项 / Sleeping 零定时器实测 / 多显示器实拖

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-20
**Blocked by:** 07-ball-state-machine-core（已 done，但下列四条 AC 从未通过）
**Parallel slots:** ≤2 sub-agents（A：热键接线与错误语义；B：交互/多显示器/零定时器断言）
**Spec refs:** SPEC-08 §2, D43 状态机, SPEC-05 §6（热重载）, 票 39 前置, registry A1–A7
**来源:** 2026-09-20 AC 补裁（`docs/reports/pending-and-issues.md` §"AC 补裁遗留"）+ 实况签收发现

## What to build
票 07 被标 DONE，但补裁查明它的 6 条 AC 里有 4 条**从未被验证**、其中一条标的东西是死代码。
owner 2026-09-20 裁定：**现在插队修，不等 S3**。本票消化 registry 的 A1–A7。

1. **A1 热键接线（最高优先）**：`Ball.RebindHotkeys`（`internal/ball/ball_windows.go:693`）
   **全仓零调用者** → 改 `[hotkey]` 配置后热键根本不会重注册。要接
   `config.Manager` 的热重载回调 → `RebindHotkeys`，并把 `registerAll` 的失败语义拆开：
   **"被他人占用"（ERROR_HOTKEY_ALREADY_REGISTERED=1409）** 与 **"未尝试/其它错误"** 必须给出
   不同提示（现状是一律 `slog.Warn` 后丢弃，用户完全不知道键没生效）。
2. **默认唤起键换掉**：`Ctrl+Alt+W` 在 owner 本机被第三方程序占用（实测：全 84 个候选组合里
   **只有 W 这一个被占**，其余 83 个空闲）。裁定改为 **`Ctrl+Alt+Q`**；`Ctrl+Alt+Space`
   本机也空闲，作为可配置备选写进文档（注意部分中文输入法会抢 Space 组合，故不作默认）。
   改默认值属 SPEC-08 之外的默认参数，但**仍要在票里留痕**，别让它悄悄变。
3. **A2 交互四项**：`HitTest()` 只覆盖了 4 个条款里的 1 个。补齐：单击球 Sleeping→Listening、
   从 Confirming 用 Esc/单击取消、**不抢焦点**（真起一个有焦点的记事本类窗口验证）、
   透明区点击穿透。
4. **A3 Sleeping 零定时器实测**：现有测试断的是**策略表**，不是活句柄。改成实测
   `EnumThreadWindows`/timer 句柄集合为空（`winlive` tag 下），且**不能被 `t.Skip` 掉**。
5. **A4 多显示器实拖**：detach→primary 与持久化/恢复有测，但"拖到第二块屏"从未真执行。
   本机只有单屏（`\\.\DISPLAY4` 3440x1440）——**若无法真机验证就如实留着不勾**，
   并写明需要什么硬件，不许用双屏 mock 冒充。
6. **A5–A7**：见 registry 对应条目，逐条消化或明确登记为"待硬件/待票 39"。

## Key constraints
- **测量必须独占**：本票含 winlive 与句柄/时序断言。跑之前**必须确认没有其它 balldebug/wisp 进程在跑**
  （实况签收期间出过两次污染：`TestBallLiveLifecycle` 按窗口类名 `FindWindowW` 会匹配到演示球；
  "四个热键全失败"也是这么来的）。测试开头加一道前置断言：检测到同名类窗口已存在就 **fail 并提示**，
  不要静默 skip——静默 skip 正是这轮漏检的根因。
- 不改 20 态语义与 D43 迁移表（视觉重做在票 62，本票只管功能与接线）。
- 禁止裸 `go func(`（热键消息回路要有 owner+recover）；禁止墙钟差判超时；禁止 emoji。
- 与票 62 的边界：本票**不碰材质/动画/吸附**，只让"球能被唤起、配置改了能生效"成立。

## Acceptance criteria
- [ ] 改 `[hotkey]` 后热键真的重注册（端到端：改配置→按新键→球响应；旧键不再响应）。
      **机制已证、生产触发者缺**：见 Progress log 2026-09-20 段 2 的 R12 裁决——`cmd/wisp` 既不建球也不持
      `config.Manager`，装配根接桥是票 12 的活，故本框不勾。
- [x] 注册失败两类语义分开：被占用 vs 未尝试，各自给出用户可见提示（不再只有一行 Warn）。
- [ ] 默认唤起键为 `Ctrl+Alt+Q`，且 `Ctrl+Alt+Space` 作为可配置备选写进文档。
      默认值已改并被 `TestDefaultSummonHotkeyIsQNotW` 钉住；**文档半边未做**：列有 `[hotkey]` 的三份文档
      （SPEC-03、SPEC-08、PLAN）对本票全部冻结，无处可写。见 Progress log。
- [x] 交互四项各有真机测试且**不被 skip**：单击唤起、Confirming 取消、不抢焦点、透明区穿透。
- [x] Sleeping 零活动定时器句柄为**实测断言**（winlive），策略表断言保留但不作为唯一证据。
      （残留：winlive 仍不进 CI —— A3 的 CI 半边转 nightly，见 Progress log。）
- [ ] 多显示器：真拖到第二屏验证并留证；本机无第二屏则**保持未勾**并写明所需硬件。
      本机物理单屏（`\\.\DISPLAY4` 3440x1440），**硬件缺席**，未做任何双屏 mock。
- [x] 测试前置：检测到外部同窗口/演示进程时 fail-fast 并提示，禁止静默 skip。
- [x] registry A1–A7 逐条标注"已修/待硬件/移交票 xx"，不许无声消失。
- [ ] 对抗验收由非实现者执行，报告含与本表 **1:1 的裁决表**（README 规则 6）。

## Progress log (append-only, newest last)

### 2026-09-20 段 1（前一位代理，150 轮上限前落地：`901f334`/`0b6ae92`/`2d063f0`/`307e24e`）
- 接线 `HotkeyReloader`（config.Manager ↔ `RebindHotkeys`），拆开注册失败语义（`HotkeyTaken`=真 1409 /
  `HotkeyUnparsable`+`HotkeyDisabled`=未尝试），默认唤起键 `Ctrl+Alt+W`→`Ctrl+Alt+Q`（R10）。
- winlive 实况测试：热键链四项 + 交互四项 + 零定时器句柄实测 + 污染 fail-fast 前置。
- **真缺陷**：球窗以 `dwStyle=0` 建，OS 白送 `WS_CAPTION|WS_BORDER` 隐形边框，layered 窗不画非客户区，
  于是 client 只有 56x33 而所有客户区规则都按 WINDOW 矩形算 → 点击落点比球心偏 (8, 31)，
  点球正中返回 `HTTRANSPARENT` 穿到桌面。改 `WS_POPUP`（`2d063f0`）。
- live 断言改按测试自身 `DebugHWND()` 而非窗口类名（A5 的假红根因）+ 每轮句柄日志（A6 取证用）。

### 2026-09-20 段 2（续做代理，预算 40-70 次工具调用）
- **`cmd/balldebug/main.go` 复检落地**（`620f565`）：段 1 被回退的 42 行编辑重新应用并补完——
  `-config <path>` 开真 `config.Manager`、经 `ApplyHotkeyDefaults` 喂 `NewHotkeyReloader`、
  `Refresh=CheckAndReload`、`OnReload` 挂桥、1s 轮询；启动打印 `[hotkey]` 全部结果
  （`Problems()` + `hotkeySummary` 每键一行）。缺失的 `hotkeySummary` 已写；轮询协程不裸奔
  （`runHotkeyBridge`：main 持有、recover、cancel→join）。
  `b.HotkeyReport()`/`b.ConfiguredHotkeys()`/`ball.NewHotkeyReloader` 三者**均已存在**，reloader 未残缺落地。
- **A1 实跑为证**（未重做，只跑）：`go test -tags winlive ./internal/ball/ -run 'TestLiveHotkey|TestLiveMute'`
  → `TestLiveHotkeyRebindEndToEnd` PASS **且未走 skip 分支**：真 `config.toml` 手改 summon→`Ctrl+Alt+R`，
  驱动 `CheckAndReload`+`Check`，`RegisteredHotkeys()` 集合真变（`hkSummon` VK 由 Q 变 R）、一次编辑恰好
  一次 rebind、注入新键真唤起、**注入旧键 Q 不再唤起**（计数停在 2）。同轮 `TestLiveHotkeyOccupiedVsNotAttempted`
  以真 `RegisterHotKey` 蹲键拿到真 1409→`HotkeyTaken`+"occupied by another program"，unset/垃圾串判为未尝试。
- **A1 的运行实例证明**：`balldebug -config <tmp>/config.toml -state Sleeping -cycle-ms 6000` 起来后
  1.5s 内手改文件，日志出 `ball: hotkeys rebound after config change summon=Ctrl+Alt+R ... live=4`，
  进程正常收尾（`closed handles=369` / `OK`，桥已 join）。启动行 `hotkeys live=4/4 summon=... mute=... cancel=... panel=...`。
- **R12 裁决（不勾框，说清谁触发）**：改动后"真配置编辑→跑着实例的热键重注册"这条路**在 balldebug 上可达**
  （上条即证据），winlive 亦证明集合真变。但**生产触发者不存在**：全仓 `ball.New(` 只有
  `cmd/balldebug/main.go:171` 一处，`cmd/wisp` 用一次性 `config.LoadFile`（`run.go:189`、`providers.go:90`），
  既不持 `config.Manager` 也不建球。⇒ 装配根装桥属票 12：需在 `ball.New` 旁建 `config.NewManager` +
  `NewHotkeyReloader`，并把 `Check()` 挂在 watchdog tick 上——**不能只挂 `OnReload`**：`[hotkey]` 是 HOT 档
  （`internal/config/manager.go:203` 的 rest 表整体覆盖 `cur.Hotkey`，故轮询读得到新值，但 `OnReload`
  只对 RELOAD 档触发）。已登记，不外溢到票 12 的文件。
- **A1c**：`DefaultHotkeys().Summon == "Ctrl+Alt+Q"` 与"不得为 W/Space"已由 `hotkey_test.go:55`
  （`TestDefaultHotkeys`，断言体在 :60-63）钉住；
  `AltSummonSpace` 常量注释写明"可配置、非默认（中文输入法抢 Space）"。**用户可见文档没动**——
  含 `[hotkey]` 的 SPEC-03/SPEC-08/PLAN 全在本票冻结清单内（SPEC-08 还带着今天的 INTERIM 标记），
  `docs/tickets/` 为空、其余非冻结文档无一字提及 summon。→ 该框保持未勾，待 owner 解冻或票 39（配置 GUI）落地。
- **隐形边框的消费者审计（`2d063f0` 后续，逐个查）**：全仓 `GetWindowRect`/`ScreenToClient` 消费点与
  差分工具截图区域共查 9 处，**无一仍假设旧几何、也无一处做过边框补偿**（`grep AdjustWindowRect|
  DwmGetWindowAttribute|GetClientRect` 在产品码里**零命中**——正因从没人补偿，边框才能把点击悄悄挪走）：
  1. `wndProc wmNCHitTest`：客户坐标 vs `HitTest(..., wr.width(), ...)`，`HitTest` 把圆心放在
     `(sizePx/2, sizePx/2)`——只有 client==window 才成立，`WS_POPUP` 后成立。
  2. `renderFrame`/`present`：`UpdateLayeredWindow` 用 `wr.l/t` + `rend.w/h`；位图尺寸即窗口尺寸，现一致。
  3. `wmMouseMove` 拖拽：`wr.l+dx, wr.t+dy` + `SWP_NOSIZE`——位移量与尺寸无关，干净。
  4. `dock_windows.go dockGeometry`：`edgePx = wr.width()`、`orbR = edgePx/2 - RingMarginPx*scale`，
     `DockPos` 的切边（`work.L + off - half`，half=edgePx/2）**整条链都按 WINDOW 矩形**；有边框时
     edgePx=72 而真实绘制只有 56x33 → 切边断言会偏 (8,31)。现与绘制面重合。
     实况证据：`TestBallLiveEdgeDock`（`live_windows_test.go:338-350` 逐边算 `wr.L±half∓drawn` 与
     `m.Work` 比切边）在边框修复**之后**实跑 PASS（本轮 `-run 'TestBallLive|TestLive'` 15 项全绿）。
  5. `dockMoveTo`/`recenterAt`/`moveWindow`：读写同为窗口矩形，干净。
  6. `persistPosition`：存 `wr.l/wr.t`；`ResolvePosition` 用 `wr.width()`；`reclampPosition` 用窗口矩形
     与 work 区求交——三处一致。**唯一残留事实（非缺陷）**：边框时代存下的 X/Y 是"含边框窗口原点"，
     修复后同值放置的是无边框窗口，老用户第一次启动球会相对旧视觉平移 (8,31)，拖一次即自愈。
     记此留痕，不写迁移（无数据可区分新旧存量，且 D42 迁移面不在本票）。
  7. `cmd/balldebug/diff_windows.go`：截图区 = `windowRect ± margin`，对称扩张，边框只挪裁切框不偏测量；
     `shot_windows.go` 直取窗口矩形，同理。
  8. `interaction_live_test.go` 的球心/透明角：`wr.l+wr.width()/2`、`wr.l+1, wr.t+1`——修复后才是真球心
     （这正是当初抓到处）。
  9. `live_windows_test.go` 断言全部走 `DebugHWND()` 的窗口矩形（`307e24e`），无类名/无客户区混用。
- **A1d**：`TestLiveMuteHotkeyEndToEnd` PASS 且未走 skip 分支——注册表 `hkMute` live →
  `WM_HOTKEY(mute)` → `OnMuteHotkey` → `EvMuteKey` → 机器 `Muted` 且球渲染 `Muted`；再按一次
  （真注入 `Ctrl+Alt+M`）走 D43 #10 回 `Armed`。
- **A6 句柄增长定案（安静桌面，桌面确无他球，前置 fail-fast 未触发）**：
  `go test -tags winlive ./internal/ball/ -run TestBallLiveLifecycle -count=3` 全 PASS，逐轮
  `base / afterNew / afterClose`：**120 / 363 / 369（USER 1→6）**，**369 / 372 / 370（USER 6→6）**，
  **370 / 372 / 370（USER 6→6）**。判读：**不是逐球泄漏**——registry 里"104 升到 376/374"是
  把"进程内尚无球时的底座"和"建过一次 D2D/DirectWrite/COM 之后的常驻底座"当成了两个球；
  第一颗球带来 ~243 句柄 + 5 USER 的一次性进程级底座，第二、三颗球 **delta=1、delta=0**，
  起点 369→370（+1，噪声级），Close 后回落到起点值。⇒ **不立泄漏修复票**；
  D32 的 Sleeping 零定时器未被触碰（`TestLiveSleepingZeroTimerHandles` 绿，本轮无阈值放宽）。
- 门禁：`gofmt -l` 空、`go vet ./internal/ball/ ./cmd/balldebug/` RC0、`go vet -tags winlive ./internal/ball/` RC0、
  `go test -count=2 ./internal/ball/` → `ok ... 0.062s`；winlive 15 项 `-run 'TestBallLive|TestLive'` 全 PASS。
- **未做/未勾**：生产装配根接桥（票 12）、`Ctrl+Alt+Space` 的用户可见文档（文档冻结）、双屏实拖（硬件缺席）、
  winlive 进 CI（A3 残留半）、对抗验收（须非实现者）。
