# 245-c1 普查 — 「借走窗口」收尾仪器该长什么样（只读，禁 Go）

> 票 245 已裁乙＝稳态不绑裸 Esc、只在 Confirming 那 2–3 秒借走；诚实残余＝那两三秒里
> 桌面其他程序的 Esc 仍被借走（SPEC-12 §5 在册推迟项）。本腿任务＝量清**借走窗口的收尾
> 仪器**：① 借/还调用点现量；② Confirming 进/出信号链；③ 既有钉有没有一枚盖住
> 「Confirming 期间被借、离开后归还」这一对；④ 若没有，最小新钉形状（只给形状不写码）。
> 行号全部基于起手 HEAD `19b63425`／骨架 `f17b1165` 时点，并发腿落件后会漂（台账有先例）。

## §0 起手锚

- 分支 `dev`；起手 HEAD：`19b634255788f1e04fd5e9d2aad0ba08971583ee`；骨架 commit `f17b1165`。
- 关键在册裁定：
  - 票 245 乙形（`.scratch/wisp/issues/245-*.md:21-26`）＋末尾残余登记五字段（同票 §残余登记）。
  - **10-01 裁定（任务书所指）＝台账 A539 §3**（`docs/reports/pending-and-issues.md`，
    2026-10-02 16:4x 记）：票 258 选形＝形 A，**边界**逐字「rebind 丢 in-flight Esc borrow
    （`ball_windows.go:807-812`）那笔残留**不归本票修**（票 245 在册）」；票面侧同句在
    `.scratch/wisp/issues/258-*.md:46`。⚠ 该条台账日期戳是 10-02、记的是 258 选形裁定；
    本腿未找到更早的独立「10-01 当天」同名裁定条目——若编排者另有所指，以那条为准、本条作线索。
  - 票 246 装配判例 A481：装配根 `cmd/wisp` 建门、以函数值注入球宿主。
- 本腿边界：⛔ 禁跑任何 Go 命令（所有实跑读数格具名弃权，见 §4）；⛔ 不改产码；⛔ 票面勾选框
  一枚不碰（票 245 现读勾 4／未勾 5，本腿零触碰）；`frontend/**`／`design/**` 未进入；冻结件
  （D43 转移表、SLO/thresholds、golden）未触碰。grep 显式根＝`internal/ball`、
  `internal/agent/approval`、`cmd/wisp`、`cmd/balldebug`。

## §1 借还两点的现量（file:line）

### 1.1 借（注册 cancel id ＝ 裸 Esc）

产码唯一链，自上而下：

- 调用点（常驻腿）：`cmd/wisp/resident_approval_windows.go:441` — `b.TakeEscForCancel()`，
  在 `ballCardUI.Prompt`（:428 起）内，**仅当 `p.Level == "L1"`**（:440）；同一函数 :454 再读
  `b.EscTakenOver()`，借被拒时打「裸 Esc 未借到」警句。
- 球 API：`internal/ball/ball_windows.go:873-889` — `TakeEscForCancel`，经 `b.uiRun` 在 ui-sta
  线程同步执行；幂等（`b.escTakenOver` 早退，:875-877）；借被拒时把报告行换成
  `cancelFailedLine`（:881）并同步 `registeredHotkeys`。
- Win32 原语：`internal/ball/hotkey_windows.go:543-545` — `takeEsc` → `takeEscWith`（:531-539：
  先 `unreg(hkCancel)` 清残、再 `reg(hkCancel, escBorrowAcc())`）。`hkCancel = 3`（:32）；
  `escBorrowAcc()`（:524）＝ `{Mods: modNoRepeat, VK: vkEscape(0x1B)}`＝**无修饰键裸 Esc**；
  :541 注释「It is the ONLY path that registers the cancel id (ticket 245)」；idle 注册遍
  （`registerAllWith` :453，`hkCancel` 分支 :457-464）对该 id 只分类不注册（`cancelIdleLine`
  :428，好绑定＝`HotkeyStandby`）。

旁支调用者（非产链）：`cmd/balldebug/main.go:635`（`syncMachineToBall`）；测试：
`internal/ball/live_windows_test.go:126`、`internal/ball/interaction_live_test.go:141`、
`internal/ball/hotkey_live_test.go:310`。

### 1.2 还（注销 cancel id）

- 调用点（常驻腿）：`cmd/wisp/resident_approval_windows.go:491` — `b.ReleaseEscAfterSession()`，
  在 `ballCardUI.settleOrb`（:483-493）内；settleOrb 先 `cards.AwaitingHuman()` 复查（两张卡
  时第一张撤下不动第二张的借，:484-486），然后还键＋`b.SetState(StateSleeping)`（:492）。
- 球 API：`internal/ball/ball_windows.go:895-907` — 幂等；`internal/ball/hotkey_windows.go:553/558`
  — `releaseEscWith` 只 `unreg(hkCancel)`，**不重绑配置键**（:547-553 注释：重绑＝把产线默认
  裸 Esc 直接放回桌面＝本票缺陷回潮）；报告行回 `cancelIdleLine(b.cancelBinding)`（:904）。
- settleOrb 的触发面＝常驻腿内三条（§2.3）。

### 1.3 残余现场：rebind 丢 in-flight borrow

- `internal/ball/ball_windows.go:807-812` 注释逐字自述：「A rebind drops an in-flight Esc
  borrow (escTakenOver below): the new binding set is the idle set... A host that rebinds while
  a card is waiting therefore has to call TakeEscForCancel again on the next state sync -
  **no caller does that today**... the residual is named in ticket 245」。
- 机制：`RebindHotkeys`（:813-825）＝ `unregisterAll` ＋ `registerAll`（idle 集）＋
  `b.escTakenOver = false`（:822）。即卡片还挂着、bookkeeping 清零、cancel id 真被注销——
  Confirming 剩余时间里既没有 cancel 键、报告又回 standby 行。
- 产码驱动者：`internal/ball/hotkey_reload.go:96`（`HotkeyReloader.Check` →
  `binder.RebindHotkeys`；差异检测 :93、节流 :97）。**装配位现量**：全仓非测试
  `NewHotkeyReloader` 只有 `cmd/balldebug/main.go:237`（`-config` 桥）；常驻腿建球是
  `cmd/wisp/resident_ball_windows.go:171` `Hotkeys: ball.DefaultHotkeys()`，**没有桥**。
  所以「rebind × in-flight borrow」这对今天在产线进程不可达（balldebug 是旁支调试器）；
  A539 形 A 落地后常驻腿才长出这个面，残余才从「注释里的话」变成「可发生的事」。

## §2 Confirming 进出信号链（常驻腿真实触发点）

### 2.1 进（借键与 SetState 的信号源）

1. `residentApproval.askConfirmation`（resident_approval_windows.go:205-250）：live++（:214）→
   `AdmitTextTask`（:230）→ 按 level 走 `gate.PendingWindow`（L1，:242）或 `PendingApproval`
   （L2，:244）。
2. `Gate.PendingWindow`（internal/agent/approval/gate.go:246-344）：admitCheck → openWindow →
   `deadline := g.clock.After(g.window)`（:285，**窗口默认 3s**＝`approval.DefaultL1Window`，
   `internal/agent/approval/queue.go:116`，界 `MinL1Window/MaxL1Window` 2s/3s，:120-122）→
   **`g.ui.Prompt(ctx, p)`（gate.go:288）＝「Confirming 进入」信号本身**。
3. 注入 UI `ballCardUI.Prompt`（resident_approval_windows.go:428-460）：L1 ⇒ `TakeEscForCancel()`
   （:441）→ `b.SetState(StateConfirming)`（:443，经 `stateForCardLevel` :501-506；L2 ⇒
   AwaitingApproval、**不借**，:440 的 `if` 与 ：420-427 注释）。
4. `SetState` 在 `cmd/wisp` 产码里**只有两枚**调用点：`resident_approval_windows.go:443`（进）
   与 `:492`（出）。没有独立的 statemachine 监听者、没有 re-sync 循环——orb 状态与 Esc 借还
   完全由 gate 对注入 UI 的 Prompt/Update 调用驱动。

### 2.2 出的三条路（gate 侧）

- **a. 否决**：`Gate.Veto`（gate.go:415-472）把否决投进 `w.vetoes` → PendingWindow
  `case v := <-w.vetoes`（:296）→ 写 ANSWER-VETO 行（:314）→ `ui.Update(EventDismissed)`
  （:316）→ return AnswerVeto。
- **b. 到点执行**：`case <-deadline`（:319）→ markStarted → 写 ANSWER-EXPIRED 行（:332）→
  `ui.Update(EventStarted)`（:334）→ return AnswerTimeout（＝执行）。
- **c. 任务取消**：`case <-ctx.Done()`（:340-341）→ 直接 return AnswerReject，**gate 不发任何
  UI 事件**（resident_approval_windows.go:220-226 注释具名这一点）。

### 2.3 收尾汇聚点（还键发生地）

- 路 a/b：`ballCardUI.Update`（resident_approval_windows.go:465-474）对
  `EventDismissed/EventStarted` 做 `cards.Forget` ＋ `settleOrb()`（:467-469）。
- 路 c：无 UI 事件，兜底＝`askConfirmation` 的 defer `ra.ui.settleOrb()`（:216-228，注释逐字
  「A card that left without settling here would leave the orb in Confirming and the desktop
  without its Esc key, which is precisely the residue ticket 245 was filed to remove」）。
- `settleOrb`（:483-493）：`AwaitingHuman()` 复查（还有别的卡在等就早退）→
  `ReleaseEscAfterSession()` → `SetState(StateSleeping)`。
- 另一条非常规出口：`detachBall`（:344-353，球窗拆除前）只卸载 channel 记账，**不还键**——
  还键靠 `Ball.Close` 的 `unregisterAll`（ball_windows.go:944）兜底。
- 退出序列：`cancelTaskRoots`（:283-331）拒绝/作废挂起卡 → `ra.cancel()` → 等 live 归零；
  归还由上述 settleOrb（经路 c 的 defer）或 Close 兜。

## §3 既有钉射程逐枚读断言

结论先行：**任务书问的那一对（Confirming 期间被借、离开后归还）——有，且是最强形**。
`cmd/wisp/resident_approval_live_246_windows_test.go` 的
`TestLive246ConfirmingCardBorrowsEscVetoesAndReturns`（:72-188）盖住它，含第二进程观察。
逐枚：

### 3.1 常驻腿真装配层（winlive，cmd/wisp）

- `TestLive246ConfirmingCardBorrowsEscVetoesAndReturns`（:72-188）：
  - 正控：`-steal` 读数 keydown==0 且 wm_hotkey==1（:76-81，尺自身防盲）；基线：无 Wisp 时
    keydown==1（:83-88，桌面已有他人占裸 Esc＝响亮红非 skip）。
  - idle：`requireIdleCancelSlot246`（:108；helper :350-360：`Live()==3`＋EscTakenOver false）。
  - 真 L1 卡：`AskOnTaskRoot`（:115-121）→ AwaitingHuman（:124）→
    `WaitingState()==Confirming`（:132-135）。
  - **借**：`rb.b.EscTakenOver()`（:136-138）＋ `HotkeyReport().Live()==4`（:139-141）。
  - **(a) 桌面被借走**：第二真进程注入 Esc `keydown_esc==0`（:145-149）。
  - **(b) 否决生效**：AnswerVeto（:157-160）＋理由含「Esc」（:161-163）＋审计 ANSWER-VETO
    `channel=esc`（:164-169）。
  - **(c) 归还**：EscTakenOver false（:173-175）＋ idle 名册三枚（:176）＋ **第二真进程重新
    收到 keydown_esc==1**（:177-182）＋无残留等待（:183-185）。
  - 判定：**「借＝Win32 层真注册」由 Live()==4＋第二进程 keydown==0 两枚合钉；
    「还＝真注销」由 idle Live()==3＋第二进程 keydown==1 合钉**。第二进程读数是
    `escBorrowProbe`（注册一个外来 id 试裸 Esc）替代不了的最强证据。
- `TestLive246ExitRefusesAHangingL2Card`（:195-283）：L2 卡**不借**（:240-242「an L2 card
  borrowed the cancel key: a 300s deadline is not a 2-3s window」）；退出序列后 idle 名册
  （:274）＋ `Problems()==0`（:278-280）。
- `TestLive246ExitAbandonsAHangingL1Window`（:289-341）：L1 借（:317 waitFor EscTakenOver）→
  退出弃窗（RESIDENT-WINDOW-ABANDONED 审计行 ：333-336）→ EscTakenOver false（:337-339）＋
  idle 名册（:340）。**盖住路 c（ctx.Done）那半的收尾**。

### 3.2 球层（winlive，internal/ball）

- `TestLiveConfirmingCancelAndEscReturned`（interaction_live_test.go:133-274）：D43 走到
  Confirming（#4→#11→#15→#17，:181-190）→ 借（:200-205）→ `wmHotkey(hkCancel)` 否决到 Acting
  （:211-214）→ 还（:219-222）→ 再进 Confirming 用 `injectBinding("Esc")` 物理注入否决
  （:225-240）→ 再还（:242-250）→ 单击否决路径收尾同样还（:253-273）。**球层 B1 全对**。
- `TestBallLiveLifecycle`（live_windows_test.go:44-143）：idle 三枚（:74-78）＋借/还对
  （:118-134）。
- `TestLiveHotkeyRebindEndToEnd`（hotkey_live_test.go:121-220）：rebind 后 idle 集 3 枚
  （:171-175）；**rebind 前**名册 idle（:149）——**未在 borrow 挂起中 rebind**。
- `TestLiveHotkeyOccupiedVsNotAttempted`（:225-314）：1409/不尝试分族（:252-267）；借/还
  （:309-313）。
- 共用量尺：`requireEscBorrowed`（hotkey_live_test.go:90-116：Live()==4＋VK 0x1B＋
  `escBorrowProbe`＝外来 id 试注册裸 Esc 必须回 1409）；`requireEscReturned`（:456-462）；
  `requireIdleRoster`（:429-452：idle 三枚＋cancel standby 行＋裸 Esc 桌面自由）；
  `escBorrowProbe`（:57-81）；正控/基线协议同 246。
- 无桌面层（hotkey_status_test.go）：`TestCancelBorrowRoundTrip`（:202-266：借 1 次、还
  **不重绑** attemptsOf(hkCancel)==1、再借无还必败）；`TestCancelBorrowFailureIsAProblemLine`
  （:268-）；`TestStandbyIsNotDisabledAndNotAProblem`（:291-）；假注册器分类族（:320-）。
  rebind 桥确定性层：`TestHotkeyReloaderRebindsOnConfigChange`（:454-503，改一次→一次 rebind、
  不变→零 churn）、`TestHotkeyReloaderSectionsHookIgnoresOtherSections`（:507-520）。

### 3.3 差距判定（逐格对答任务书③）

- 「Confirming 期间被借、离开后归还」**这一对：被盖**（3.1 第一枚，最强形）。
- **未被任何现役钉盖住的一对＝rebind × in-flight borrow**：`TestLiveHotkeyRebindEndToEnd`
  的 rebind 全部发生在 idle；246 族三枚的 L1 卡都在无 rebind 的窗口里借/还。全仓 grep
  `RebindHotkeys`（测试调用点 :247/:272/:297 与 recorder :444）没有一处把它与
  `EscTakenOver()==true` 并置。这正是 ball_windows.go:807-812 注释自己承认、A539 边界点名
  不修的格子。
- 次级空档（非任务书直问，但仪器设计要知）：双卡重叠借还（settleOrb 的 AwaitingHuman 早退
  分支）没有专钉——interaction_live 的三段借/还是串行的，不是重叠的。

## §4 新钉形状建议＋我可能判错的条目＋量不到的格子

### 4.1 任务书问的最小新钉：不需要

「Confirming 期间被借、离开后归还」已经有
`TestLive246ConfirmingCardBorrowsEscVetoesAndReturns`（§3.1 第一枚），两枚断言（借＝Win32
真注册、还＝真注销）都在里面且带第二进程桌面读数。若编排者要的「收尾仪器」指的就是这一
对，答案是**复用，不新增**；新腿要做的只是把它当撞钉预检的正控。

### 4.2 若仪器目标是把 A539 形 A 落地后的 rebind 残余也网住（建议形状，不写码）

落在 `cmd/wisp` 的 winlive 族（复用 246 live 文件全部基建：`buildEscListener246` /
`observeEsc246` / `captureAudit246` / `waitFor246` / `requireIdleCancelSlot246`），一例覆盖
两头：

- **编排**：起球＋门（同 246 live 第一例）→ 起 L1 卡（AskOnTaskRoot）→ 等
  AwaitingHuman＋EscTakenOver（复用 waitFor246）→ **在卡挂起中**调
  `rb.b.RebindHotkeys(ball.DefaultHotkeys())`（或等价 HotkeyConfig）——这一步模拟形 A 桥
  （`HotkeyReloader.Check` 在卡挂着时 diff 出一次 rebind）。
- **断言 N1（借被 rebind 丢＝红向正控）**：rebind 后立即 `HotkeyReport().Live()==3`
  （cancel 回 standby 行）且 `EscTakenOver()==false`——把 ball_windows.go:807-812 注释的
  现状钉成**具名红灯**而非隐含行为。若 A539 后续裁定此残余要修，这枚钉反过来当修完后的
  绿钉（期望改 `Live()==4`＋EscTakenOver true），红绿两向都由它承担。
- **断言 N2（桌面读数）**：rebind 后立刻用 `observeEsc246(t, rig, "-watch")` 读第二进程——
  残余未修时 keydown_esc==1（桌面拿回了键，Wisp 的 Confirming 却没有 cancel 键＝用户侧
  可见的缺陷形状）；**窗口还开着**（AwaitingHuman 仍 true）是这枚钉区别于 246 第一例的
  关键前提。注意窗口 3s（DefaultL1Window）会先到点执行，observeEsc246 自身 READY→READING
  流程要排在 deadline 内，或把该例的断言窗口压在借期内（246 第一例的 :145 是排在内的先例）。
- **断言 N3（收尾不变量）**：rebind 之后卡片再走完（否决或到点），`settleOrb` →
  `ReleaseEscAfterSession` 幂等早退（escTakenOver 已 false）→ 终态 idle 名册三枚＋
  `Problems()==0`（复用 requireIdleCancelSlot246＋Problems 检查，246 :278 同形）。这枚钉的
  价值＝保证 rebind 残余**不腐蚀归还路径本身**（报告行不烂、close 后无泄漏）。
- 若还要一枚「归还后不再复活」的负向钉：卡撤下后 `RegisteredHotkeys()` 里
  `hkCancel(3)` 不在场（断言「还＝真注销」的 Win32 名册面，与第二进程读数互为正反）。

### 4.3 我可能判错的条目（自查）

- §0 的「10-01 裁定」出处：我只找到 A539（10-02 落笔）里的边界句；若编排者指的是另一条
  当天台账，本报告 §1.3 的裁定引用就挂错了条（内容不受影响，出处要改）。
- 246 live 第一例的 `HotkeyReport().Live()==4` 断言（:139）读的是**报告**不是
  `RegisteredHotkeys()`；「Win32 层真注册」在该例里严格说是靠第二进程 keydown==0 兜底的，
  单看报告行可被「报告烂了」骗过。新钉若要纯 Win32 面，应加 `RegisteredHotkeys()[hkCancel]`
  读数（4.2 N1 的括号项），不要只信报告。
- 「rebind × borrow 今天不可达」依赖两个现量：常驻腿无 `NewHotkeyReloader`（已 grep 复核）
  且 `253-r1` 等并发腿没有正在落形 A——本腿只对起手 HEAD 负责。
- 行号漂移：本仓多腿并发、追加式修改会推移行号（台账多次记录先例）；所有 file:line 以
  `19b63425` 为准。
- `TestLive246ExitAbandonsAHangingL1Window`（:289-341）没有 `bindBallHost` 返回值检查
  （:301 直接调），严格说它对「门绑定成功」的假设比第一例弱——不影响本报告结论，但若有人
  拿它当模板要注意。

### 4.4 量不到的格子（禁 Go ⇒ 实跑读数具名弃权）

- 任何 `go test -tags winlive` 的当日读数（246 live 族三枚＋球层热键族）——本腿零实跑，
  只做了静态断言读法；「这套钉今天是不是绿的」本腿不背书。
- `-steal`/`-watch` 正控与基线在**当前这台机器**的实际数值（是否有人占裸 Esc）——同上弃权。
- N1-N3 三枚建议钉的可行性（能否在 3s 窗口内排完 observe 流程、rebind 与借的线程交错是否
  会撞 `uiRun` 死锁防护）——需真机才能验证，本腿只给形状。
- 常驻腿形 A 落地后的实际 rebind 频率（桥 tick 1s × config mtime）——属票 258 落地腿的量。

## 附：本腿产出与提交

- 本文件（census.md）五节；骨架 commit `f17b1165`（起手 HEAD `19b63425`）。
- 零产码、零票面勾选、零冻结件触碰；未 push。
