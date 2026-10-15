# ticket 319 · 319-v2 非实现者对抗验收（v2：AC#2 两半＋AC#3 复核）

- 角色：验证腿 **319-v2**，非实现者（实现者＝`319-r3` `60f3b1ec`＝"不藏窗"半、`319-r4` `259cce15`＝"抬前台"半）。
  本腿零产码：不改任何跟踪 Go 件、不翻 `Status:`/AC 框、不碰台账/HANDOVER；只追加
  `probes/319/v2/`（6 件）＋本表＋票尾 1 bullet。
- 锚：起步与全程 HEAD＝`259cce15`（腿中未动）。319 面 8 枚与
  `git show 259cce15:<f>` sha256 逐一相等；三枚生产件 sha＝
  `resident_windows.go 816017b6…`、`panel_host_windows.go 9bb2be45…`、
  `panel_resident_windows.go 4154d1d0…`，每发变异跑前后恒等。
- 方法：只读仓＋**仓外**快照（`git archive 259cce15 | tar -x -C /tmp/wisp319v2snap`）
  ＋`-overlay` 注入。本腿自造的两枚活体探针（源码＝`probes/319/v2/01、02`）与
  四枚一行变异（副本在 `C:\Users\swq\AppData\Local\Temp\wisp319v2-mut\`，清单见
  `probes/319/v2/00-index.txt`）都不落进跟踪树；变异跑后 `git diff -- <319 件> | wc -c`＝0。
- 环境：`GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\wisp319v2gotmp'`；
  PATH 前置 `third_party/sherpa-onnx`；机主实例 `wisp.exe` PID 40296 每发活体捕获都在场、未扰
  （tasklist 自证行在每份 run 记录头部）；winlive 一次一枚、一枚一进程；⛔ 未拉第二个 `--env dev`。
- 读数件：`.scratch/wisp/probes/319/v2/`（00–05）。以下引文全是**本腿自己跑出来**的逐字行，
  ⛔ 没有把 r3/r4/编排者的日志当读数引用（它们只被当"待验主张"）。

## 1. 逐格判语（先行）

| 格 | 判语 | 一句话理由（读数在 §2–§7） |
|---|---|---|
| **AC#2-a 不藏窗** | **成立** | 自造活体探针：面板已可见＋prior 在前 ⇒ 生产二次启动路后 `IsShown=true`、HWND 不变、窗口活着；把 hook 拧回 `h.toggle`（变异 A）⇒ 红名点到"被 HIDDEN"；还原绿。 |
| **AC#2-b 抬前台** | **成立** | 同探针 + 自造 prior（同进程真窗、先当前台，`prior=0xACA10AE`）⇒ `GetForegroundWindow()=0x7551138=面板 HWND`、`SetForegroundWindow answered 1`（`raise request landed…set_foreground_ret=1`）；摘 `RequestRaise` 一跳（变异 B）⇒ **只有**前台组红（红名"left the panel in the background"，`prior=0x31712A0, fg=0x31712A0`），三枚可见性断言同跑全绿 ⇒ 两半可分离、这半的断言承重。⛔ 未拿 `IsShown` 当替身。 |
| **AC#2-c 隐藏态** | **成立** | 隐藏面板 ⇒ 二次启动后 `IsShown=true`、同一 HWND（C27 不销毁）、活着，且 `foreground=0x42505B0=面板`（SFW answered 1）。 |
| **AC#2-d 死线重建** | **成立** | 真窗 `WM_SYSCOMMAND/SC_CLOSE` ⇒ 死线、`IsCreated=false hwnd=0`；二次启动 ⇒ **新 HWND 0xA7708A0→0xA210958**（旧窗 `alive=false`）、`IsCreated/IsShown=true`、`foreground=新 HWND`（SFW answered 1）＋生产日志 `starting a fresh panel thread … why="the previous panel thread had exited"`。 |
| **AC#3 三枚新钉** | **成立** | N4-1/N4-2（headless）与 N4-3（winlive，含 r4 第三条断言）逐枚**自造变异**：A/B/C 各按设计红名逐字（§5），不带 overlay 还原全绿；`-run '319'` 两包 rc=0。 |
| **AC#3"toggle 语义未动"** | **附条件**（语义有牙、**接线无钉**） | "可见态仍隐藏"由 (a) N4-1 的 control 半（变异 D ⇒ `:277/:280` 红）＋(b) 活体 `TestTicket255r8…`（变异 D ⇒ `:156` 超时红；现树绿）守着；但**"托盘/热键仍走 `holder.toggle`"那条接线（`resident_windows.go:268-270`）找不到会红的用例 ⇒ 无钉**（见 §6）。 |
| D29 焦点族回归 | **成立** | 八个分段逐名 rc=0（§7），无 314 死区；`TestAC4FocusReturnToPriorWindowGap33r5` PASS(0.89s)、STA PASS、RealWindowHop PASS。 |
| AC#4 门禁＋面核 | **成立** | build/gofumpt/vet/d22scan/长度门全 rc=0；两笔登陆只落 319 面＋probes＋票面，冻结点名零命中（§7）。 |

**AC#2 能不能翻？——建议：能翻（成立）。** 两半各自有独立探针、"绿⇒变异红（红名点到那一半）⇒还原绿"三态齐；c/d 两形四读齐；与 `319-v` 的 D1 判定对齐（可见态**不再**被藏，且**新**增"真到前台"读数）。残余三项见 §8（皆已具名、不推翻成立）。
**AC#3 能不能翻？——建议：能翻（成立）。** `319-v` 判 AC#3 成立时钉的 `show` 半；本轮新增的 r3 两枚 headless + r3/r4 的 winlive 第三条断言都逐枚自造变异证明了"SHA 改动 ⇒ 红名 ⇒ 还原绿"。缺口只有 §6 的"接线无钉"一条，属**新增钉子候选**、不属 AC#1/AC#2 钉子的欠账。

## 2. AC#2-a（不藏窗）三态——读数逐字

探针 `TestProbe319v2VisibleStaysShownAndRaised`（`probes/319/v2/01`）。

- **绿（现树，14:57:57）**：`BEFORE: panel hwnd=0x7551138 mgr.IsShown=true alive=true; prior=0xACA10AE foreground=0xACA10AE; resident raises=0`；
  `AFTER: mgr.IsShown=true hwnd=0x7551138 (same=true) alive=true; foreground=0x7551138 (panel=true, prior=0xACA10AE); resident raises 0->1; SetForegroundWindow answered 1` ⇒ **PASS**（0.69s）。
- **变异 A（`resident_windows.go:690` `h.show`→`h.toggle`，14:58:24）**：红名逐字 `AC#2-a RED (the visible panel was HIDDEN by the second launch): the window 0xA7208A0 was shown before the activation and mgr.IsShown() is false after it; … hiding it is the D1 defect` ⇒ **FAIL**（读数 `IsShown=false, raises 0->0`）。这同时复现了 `319-v` 的 D1 现场形。
- **还原绿（14:59:10）**：同绿读数（`hwnd=0xAA5112E`，`raises 0->1`，SFW answered 1）⇒ **PASS**；`git diff` 0 字节、`git status --porcelain -- cmd/wisp/ internal/proc/` 0 行。

## 3. AC#2-b（抬前台）——prior/foreground 句柄与 SFW 返回逐字

- **绿（14:57:57）**：`prior=0xACA10AE`（同进程真窗，先夺前台并断言成功，否则记为 failed measurement）；二次启动后 `foreground=0x7551138`（＝面板 HWND）、`SetForegroundWindow answered 1`；同跑生产日志 `panel host: raise request landed on the panel thread via=second-launch set_foreground_ret=1`。
  另两发独立绿（还原跑 14:59:10）：`prior=0x16F0CBA → foreground=0xAA5112E（panel），answered 1`。
- **变异 B（`resident_windows.go:614` `tookRaise := true`，14:58:50）**：
  `AFTER: mgr.IsShown=true hwnd=0x578046E (same=true) alive=true; foreground=0x31712A0 (panel=false, prior=0x31712A0); resident raises 0->0; SetForegroundWindow answered 0` ⇒ 红名逐字 `AC#2-b RED (the second launch left the panel in the background): …` ⇒ **FAIL**。
  ⚠ 这是**唯一**红的一组：同跑的 AC#2-a 三枚断言（shown/same-HWND/alive）全绿 ⇒ 断言承重且可分离。⚠ 红文里的 `answered 0` 是 `mgr.LastRaiseSetForeground()` 的零值（变异下 raise 从未跑），**不是**一次被拒绝的 Win32 调用——如实登记。
- **仓内 N4-3（`TestSecondLaunchActivationOnAVisiblePanelDoesNotHideIt319`）本腿独立复跑**：绿 15:01:03（`fg=0x1E2309C2=panel, prior=0x3EE1292, answered 1`）；变异 A ⇒ `:433 the second launch HID the panel` 红；变异 B ⇒ `:442 … left the panel in the background` 红（`IsShown` 保持 true）；还原绿 15:01:52。⇒ 仓内钉与我的探针互证。

## 4. AC#2-c/d——两形各两向读数

- **c 隐藏态（14:59:23，绿）**：`BEFORE: hwnd=0x42505B0 mgr.IsShown=false alive-while-hidden=true; prior=0x3EC1292 foreground=0x3EC1292`；
  `AFTER: mgr.IsShown=true hwnd=0x42505B0 (same=true) alive=true; foreground=0x42505B0 (panel=true); SetForegroundWindow answered 1` ⇒ **PASS**（0.80s）。向 1＝重新显示（同柄，C27）；向 2＝到前台。
- **d 死线重建（14:59:35，绿）**：`BEFORE: old hwnd=0xA7708A0 closed via WM_SYSCOMMAND/SC_CLOSE (alive after close=false); mgr.IsCreated=false mgr hwnd=0x0; prior=0xA590FB2 foreground-before-signal=0xA590FB2 (prior-in-front=true)`；
  `AFTER: rebuilt hwnd=0xA210958 (new=true) mgr.IsCreated=true mgr.IsShown=true alive=true; old window alive=false; foreground=0xA210958 (rebuilt-in-front=true); SetForegroundWindow answered 1` ⇒ **PASS**（7.57s）。向 1＝**新窗**（≠旧、旧已死）；向 2＝到前台。生产日志：`starting a fresh panel thread for this request via=second-launch why="the previous panel thread had exited"`＋`raise request landed … set_foreground_ret=1`。
- 与 `319-v` D1 对齐：可见态在本轮**不再**被藏（§2），且三形（可见/隐藏/死线）都拿到"真到前台"读数（`319-v` 当时只对冷形／跨进程探针拿到）。

## 5. AC#3 三钉＋r4 第三条断言——逐枚"弄坏⇒红→还原⇒绿"

变异说明：A＝`h.show`→`h.toggle`；B＝摘 `RequestRaise`（`tookRaise := true`）；C＝摘 `panelHolder.show` 重建条件里的 `isFinished()`；D＝`RequestToggle` 的 `IsShown` 支置假（toggle 语义牙，§6）。全部仓外副本＋`-overlay`。

- **N4-1 `TestSecondLaunchActivationDoesNotTakeTheHideRoad319`（headless）**
  绿（14:59:58）：`control toggle hides=1 shows=0; … after activation: shows=1 toggles=1 hides=1 raises=1; queued=3; spawns=1` ⇒ PASS。
  变异 A：红名 `:289 did not take the show road` / `:292 took the toggle road … HIDES it` / `:295 asked the panel thread to HIDE` / `:298 never asked for the raise hop` / `:301 queue holds 2, want 3` ⇒ FAIL。
  变异 B：只红 `:298`＋`:301`（shows=1/toggles=1/hides=1/raises=0）⇒ FAIL。
  变异 C：不受影响 ⇒ PASS（该层线程活着）。
- **N4-2 `TestPanelHolderShowRebuildsADeadThread319`（headless）**
  绿：`spawns=2; rebuilt thread shows=1 toggles=0 hides=0 raises=1 queued=2 same manager=true` ⇒ PASS。
  变异 C：红名 `:355 the spawn was asked 1 time(s) … want exactly 2` / `:358 the holder hands out the dead thread again` / `:361 the rebuild's spawn returned no thread at all` ⇒ FAIL。
  变异 A：`:370 rebuilt thread got a toggle request` / `:376 did not get the raise request` / `:379 queue holds 1` ⇒ FAIL。变异 B：`:376`＋`:379` ⇒ FAIL。
- **N4-3 `TestSecondLaunchActivationOnAVisiblePanelDoesNotHideIt319`（winlive）**：见 §3 末（绿⇒A 红 `:433`⇒B 红 `:442`⇒还原绿）。其中 **r4 的第三条断言**（`GetForegroundWindow()==面板 HWND`）在变异 B 下**正是红的那个**，且在变异 A 下与"被藏"同时红——即它不是装饰断言。
- **还原**：15:04:02 `go test -run '319' ./cmd/wisp/ ./internal/proc/ -count=1` rc=0（6 枚 319 用例全 PASS）；生产件 sha 恒等。

## 6. "toggle 语义未动"钉名盘点（派单第 4 条末问）

- **有牙的语义钉（点名）**：
  1. `TestSecondLaunchActivationDoesNotTakeTheHideRoad319` 的 **control 半**（headless）：`h.toggle` 对已可见态 ⇒ `hides=1 shows=0`。变异 D 打红：`:277 the toggle control did not take the hide road (hides=0, want 1)`＋`:280 … took the show road (shows=1) while the manager reported the panel visible`（15:02:10）。
  2. `TestTicket255r8RealWindowFollowsTheConfigAcrossAReshow`（winlive，`panel_reshow_255r8_winlive_test.go`）：`rp.RequestToggle("255r8-hide")` 后 `waitPanelTrue(!IsShown)` —— 可见面板真被隐藏、再 reshow。现树绿（15:02:24，`same hwnd: true`）；变异 D ⇒ `:156 timed out after 15s waiting for the panel to hide - this is a failed measurement, not a pass`（15:02:35）——**是超时红、不是红名**，如实登记。
- **无钉**：production 接线 `resident_windows.go:268-270`（`withPanelHost(func(via string) bool { return holder.toggle(via) })`，托盘/热键/球的手势入口）若改成 `holder.show`，我找不到任何会红的用例：`TestBallPanelGesturesReachThePanelThread` 是手工搭 hook（用 `rp.RequestToggle`），不走运行期装配。⇒ **"仍走 toggle"这一半＝无钉**；若要补，建议加一枚读 `withPanelHost` 装配的 headless 钉（归编排者裁）。

## 7. 门禁、D29 焦点族、面核

- 门禁（15:04–15:05，全 rc=0）：`go build ./...`；gofumpt v0.12.0 对 319 面 8 件 `-l` **空**（另：本腿两枚探针源 gofumpt 空）；`go vet -tags winlive ./cmd/wisp/ ./internal/proc/`；`sh scripts/d22scan.sh` clean；`scripts/check-path-length-budget.sh --with-self-test` **VERDICT GREEN**（13187 跟踪路径，over=60＝roster=60，not in roster=0）。
- D29 焦点族（分段逐名，无一撞 314 死区）：`33r11` rc=0（5 PASS）｜`TestTicket255r1` rc=0（6 PASS）｜`TestTicket255[A-Z]`（255 parse/roster/receipt 扫族）rc=0｜`TestAC4PriorFocusSurvivesARefusedPanelSample` rc=0｜`TestAC4FocusReturnToPriorWindowGap33r5` rc=0（PASS 0.89s）｜resident 族（`TestBallPanelGesturesReachThePanelThread`＋`TestBallGestureWithoutPanelHostStillRecords`＋`TestPanelThreadNameIsNotInResidentRoster`）rc=0｜`TestPanelThreadIsSTAAndExitsCleanly` rc=0｜`TestPanelHostRealWindowHopAndLifecycle` rc=0。另 `-run '319'` 两包 rc=0。（裸全包 `go test ./cmd/wisp/` 未跑＝314 死区，见 §9。）
- 面核：`60f3b1ec`＝票面＋`probes/319/r3/**`＋`cmd/wisp/{panel_death_replay_319_winlive_test.go, panel_holder_319_windows_test.go, resident_windows.go}`；`259cce15`＝票面＋`probes/319/r4/**`＋`cmd/wisp/{两枚测试件, panel_host_windows.go, panel_resident_windows.go, resident_windows.go}`。两笔并集的**冻结点名零命中**（`docs/PLAN.md`／`docs/specs/**`／阈值／golden／`frontend/**`／`design/**`）。工作树与仓外快照 8 件逐枚 sha256 **相等**；319 面 `git status` 0 行。

## 8. 分歧与残余（逐处登记，皆不推翻判语）

- **D2 升级读数（新）**：15:05–15:08 期间**无关进程**（Tabbit Browser PID 28348，hwnd 0x191204B8）持前台。此时：
  - 固定树上跑 d 形：生产 ASFW 自己**被拒**——`level=WARN msg="proc: AllowSetForegroundWindow(ASFW_ANY) was refused; the resident's SetForegroundWindow may be denied by the foreground lock" err="Access is denied."`；随后 raise 落地 `set_foreground_ret=0`；重建出的新窗 shown/alive 但 `foreground` 仍是 Tabbit ⇒ 我的探针在固定树上也红。
  - ⇒ ① 这一形里 ASFW 的"缺/拒"与"抢不到前台"**同刻可读**，给 `319-v` 的 D2（必要性未证）补上一个隔离形状（但**仍不是**现场链：现场第二实例由前台 shell 拉起，本测试进程不是 ⇒ 现场链条不可在进程内复现，残余保留）。② 同一环境下变异 B 的 d 跑与固定树同红 ⇒ **该两发不具判别力**，不能拿它们说"d 依赖 raise 跳"；c 在变异 B 下停在 prior 安装的 failed measurement（`:98 the prior window … did not take the foreground`），同样不构成结论。d/c 对 raise 跳的敏感性结论以 §3/§4 的"我们进程持前台"环境读数为准（b 半有明确结论；c/d 两形由 `mgr.Show` 的 `!shown` 支与重建冷路径承载，且 c 的隐藏-显示语义另有 255r8 钉）。
- **残余 1（承自 `319-v`/r4）**：真**两进程**形仍未升格为仓内常驻钉（`319-v` 判暂不升格；`r3/r4` 亦未做）。本腿复核：AC#2 的三形证据都在同进程节拍内，跨进程那一层只有 `319-v` 的仓外探针读数。
- **残余 2**：`RunEventLoop` 的 goroutine 在本腿两枚探针里也留到进程尾（无导出 stop）——winlive 一枚一进程，不污染；已在探针注释与 r2/v 处具名。
- **残余 3**：探针 02 的**运行期修订**与随件版本的差异（注释缩短＋4 处 Logf 折行；零语句/断言/格式串改动）已在 `00-index.txt` 明记；引用行号属运行期修订。
- **残余 4**：§6 的接线无钉。

## 9. 具名没做

- 未跑裸全包 `go test ./cmd/wisp/`（314 死区；派单要求按名单分段，已做）。
- 未重跑 winlive 冷形 `TestSecondLaunchActivationRaisesThePanel319`（不在本派单名单；r4 读数未当本腿读数）。
- 未把真两进程形升成仓内钉子（残余 1；需先定 helper 子进程形状）。
- 未在"无关进程持前台"环境下取 c/d 的绿读数（Windows 前台锁拒绝安装 prior，见 §8；改取固定树/变异树同红对照）。
- 未碰 `frontend/**`、`design/**`、台账、HANDOVER、AC 框、`Status:`；未扰 `wisp.exe 40296`；未拉第二个 `--env dev`。
