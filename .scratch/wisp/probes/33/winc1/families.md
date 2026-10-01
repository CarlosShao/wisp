# winlive 族普查（代号 `winlive-census-1`，只读腿）

- 起手锚：`date` = `Thu Oct  1 18:25:37 CST 2026`（落笔时 18:3x）；`git log -1` = `07cccefa`
  `docs(probes): 33-n1 verdict - does windows-latest carry the WebView2 Runtime`；分支 `dev`。
  ⚠ 交件时 `git log -1` 已变成 `b3eed120`（共享工作树里别的腿在动）。现量：`07cccefa..HEAD` 只 1 发
  `ledger(A508)`、`git diff --name-only` 只命中 `docs/reports/pending-and-issues.md`，
  对那 8 个 winlive 文件的 `git diff --name-only 07cccefa..HEAD -- <8 路> ` ＝ **零命中**
  ⇒ 本文件的 `file:line` 与计数在 `b3eed120` 上同样成立。
- 本程零写入除本文件；零 CPU：没跑 `go test`/`go build`/`go vet`/`d22scan`/任何脚本。
  全部读数＝`grep`/`sed`/`wc`/`head`/`ls`/`git log`/`git ls-files` 与**盘上已记录的他人读数**（逐条注明出处）。
- 未读 `frontend/**`、`design/**` 任何内容（E 族只按代码怎么引 dist 判；`design` 只数了跟踪文件条数＝30，未读一字）。

## 0. 与派单给的尺/数的差异（先报这个，因为它决定名册可不可信）

| 项 | 派单值 | 本腿实测 | 差在哪 |
|---|---|---|---|
| 尺的字面串 | `grep -rl '//go:build winlive' --include=*.go .` | **原样跑＝0 命中** | 仓里**没有**任何文件写 `//go:build winlive`；真形是 `//go:build windows && winlive`（全仓 tag 清单现量：`windows`=130、`!windows`=28、`windows && winlive`=8、`windows && wisp174r1after`=1、`cgo_sherpa`=1、`!race`=1、`!cgo_sherpa`=1）⇒ 用派单那串尺会得到"零枚"，与本族存在的 Fact 相反 |
| 文件数 | 8（前一腿） | **8**（同） | 一致 |
| 用例数 | 22（前一腿） | **21** | 见 §5 第 3 条；`func Test` 逐枚数过、`t.Run` 子测试 **0 枚**（8 文件全 0），所以差的不是子测试口径 |
| 编排者 18:3x 自记（`b3eed120`＝台账 **A508**） | "松散 `grep -rl winlive` 得 **8 枚文件**、`^func` 行 **64 枚**" | 松散尺（`internal cmd tools`）＝**16 枚文件**； tagged 8 枚文件的 `^func`＝**80**、`^func Test`＝**21** | 下表逐枚给尺。**64 我复现不出来**（最接近的组合：`internal/ball` 四枚文件的 `^func`＝50 ＋ `cmd/wisp` 的 13＋3＝66，仍差 2）；台账那句"22 应是只数 `Test*`"与本腿 21 冲突 ⇒ 名册以逐枚表为准 |
| 更早的普查 | `.scratch/wisp/probes/33/n1/verdict.md` 引 S64＝"6 枚文件" | 8 | S64 那次只数到 6 枚 ⇒ 该值过期：本腿现量新增的两枚文件首发＝`cmd/wisp/panel_host_windows_live_test.go`（`git log --reverse` ＝ **`13acad46`** 33-r5）与 `cmd/wisp/resident_task_source_live_246_windows_test.go`（**`2071f59e`** 246-r2） |
| SKIP 判红的位置 | `scripts/runtests.sh:98` | **`scripts/runtests.sh` 不存在**（`find . -name 'runtests*'` 只命中 `.scratch` 日志） | 真身＝**`tools/d22scan/runtests.sh:98`**（行号对、路径错）：`if [ "$skipped" -ne 0 ]` ⇒ 任何顶层 `--- SKIP` 判红；`scripts/portable-tests.sh:433` 同规矩＋一张记账台账（`:320-330`） |
| `runs-on: windows-latest` | `:388` 附近 | **`:388` 逐字命中**，属 job `test-windows`（`:387`）；另有一枚 `:533`（`slo-smoke`）与一枚 `:591` `runs-on: [self-hosted, wisp-slo]` | 一致 |

逐文件尺（`grep -c`，`windows && winlive` 命中的 8 枚）：

| 文件 | `^func` | `^func Test` |
|---|---|---|
| `cmd/wisp/panel_host_windows_live_test.go` | 1 | 1 |
| `cmd/wisp/resident_approval_live_246_windows_test.go` | 13 | 3 |
| `cmd/wisp/resident_ball_live_228_windows_test.go` | 3 | 1 |
| `cmd/wisp/resident_task_source_live_246_windows_test.go` | 13 | 2 |
| `internal/ball/hotkey_live_test.go` | 18 | 4 |
| `internal/ball/interaction_live_test.go` | 8 | 4 |
| `internal/ball/live_guard_windows_test.go` | 17 | **0**（纯共享守卫） |
| `internal/ball/live_windows_test.go` | 7 | 6 |
| 合计 | **80** | **21** |

## 1. 名册（8 文件／21 枚，`file:line` 为 `func Test` 声明处）

| # | 用例 | 出处 |
|---|---|---|
| — | （无测试函数的共享守卫文件） | `internal/ball/live_guard_windows_test.go:1`（0 枚 `func Test`，供下面 14 枚用） |
| 1 | `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive` | `cmd/wisp/panel_host_windows_live_test.go:36` |
| 2 | `TestLive228ResidentLegOwnsABallWindowOnTheDesktop` | `cmd/wisp/resident_ball_live_228_windows_test.go:89` |
| 3 | `TestLive246ConfirmingCardBorrowsEscVetoesAndReturns` | `cmd/wisp/resident_approval_live_246_windows_test.go:72` |
| 4 | `TestLive246ExitRefusesAHangingL2Card` | 同文件 `:195` |
| 5 | `TestLive246ExitAbandonsAHangingL1Window` | 同文件 `:289` |
| 6 | `TestLive246ResidentPipelineRaisesACardAndEscVetoesIt` | `cmd/wisp/resident_task_source_live_246_windows_test.go:77` |
| 7 | `TestLive246ExitCancelsARunningTask` | 同文件 `:238` |
| 8 | `TestLiveHotkeyRebindEndToEnd` | `internal/ball/hotkey_live_test.go:121` |
| 9 | `TestLiveHotkeyOccupiedVsNotAttempted` | 同文件 `:225` |
| 10 | `TestLiveMuteHotkeyEndToEnd` | 同文件 `:318` |
| 11 | `TestLiveSleepingZeroTimerHandles` | 同文件 `:386` |
| 12 | `TestLiveClickSummonsAndDragDoesNot` | `internal/ball/interaction_live_test.go:53` |
| 13 | `TestLiveConfirmingCancelAndEscReturned` | 同文件 `:133` |
| 14 | `TestLiveNeverStealsFocus` | 同文件 `:282` |
| 15 | `TestLiveTransparentCornerFallsThrough` | 同文件 `:351` |
| 16 | `TestBallLiveLifecycle` | `internal/ball/live_windows_test.go:44` |
| 17 | `TestBallLivePositionPersistence` | 同文件 `:147` |
| 18 | `TestBallLiveIdleBorderTransition` | 同文件 `:211` |
| 19 | `TestBallLiveEdgeDock` | 同文件 `:316` |
| 20 | `TestBallLiveEdgeDockHover` | 同文件 `:474` |
| 21 | `TestBallLiveAudioLiquidGate` | 同文件 `:625` |

## 2. 六族归属（每枚给凭据：import 什么／调什么／断什么）

**汇总（六个数）**：A＝**1** Ｂ＝**0** Ｃ＝**19** Ｄ＝**0** Ｅ＝**0** Ｆ＝**1**（合计 21）

### A｜WebView2 运行库／真窗 —— 1 枚
- `#1`：import `internal/panel`；`NewPanelManager(..., dataPath)` ＋ `hostThreadHarness.bringUp`（`cmd/wisp/panel_resident_windows.go` 那条专用 STA 线程）；断的是 `readWebviewTree(t, self).TreeWebview >= 1`（本进程树的 `msedgewebview2` 子进程数）与 Destroy 后 2s 内回到基线（`:54`、`:70`）。⛔ 不读像素、不看焦点。
  盘上已知脆弱性：`docs/evidence/s1/33-panel-host-c27-r5.md` §⑤ 记录它 33-r4 整包 6 发红 1 发／单独 5 发全绿 ⇒ **负载敏感**，正是它被搬出默认档的理由（上界仍 2s，未放宽）。
  另：它只数 `msedgewebview2` 进程名，⛔ 不钉运行库版本 ⇒ 台账 A508 的"只用 Evergreen、不用 Fixed Version"那一支前置对本枚不构成额外改动（现量：`grep -n 'version\|FixedVersion\|ReleaseNotes' cmd/wisp/panel_host_windows_live_test.go`＝0 命中）。

### B｜音频设备／麦克风 —— 0 枚
- 现量：8 个 winlive 文件里 `grep -niE 'wavein|asio|microphone'`＝**零命中**；`internal/audio/` 整包 `grep -rn winlive`＝**零命中**。
- 真依赖麦克风的用例**不在本族**：`internal/audio/hotplug_test.go:1` 是 `//go:build windows` ＋ `WISP_LIVE_MIC=1` 闸门，且已写进 `scripts/portable-tests.sh:323` 的 skip 台账。
- 唯一带"audio"名字的是 `#21`，它只调 `b.SetAudioLevel(0.9)`（`live_windows_test.go:637/648/656/685`）＝**喂合成电平**，不开设备 ⇒ 归 C，别被名字骗进 B。

### C｜悬浮球／GDI／屏幕像素比对 —— 19 枚
凭据按"读什么"分三档（同族内差别就是这三档，CI 可行性也不同）：
- **C1 只要建得出窗**（8 枚＝`#4` `#5` `#11` `#16` `#17` `#18` `#19` `#21`）：`#4` `#5`（在进程内 `startResidentBall` ＋ `proc.Boot/Shutdown`，断 `AwaitingHuman`/`WaitingState`/10 条退出记账/`HotkeyReport().Problems()`；`resident_approval_live_246_windows_test.go:195-341`）、`#16`～`#19`＋`#11`（窗存在性 `IsWindow`、`GetWindowRect`、`monitorRectsForWindow` work-area 几何、`SetWindowLongPtrW` 子类化数 `WM_TIMER`；`live_windows_test.go:44/147/211/316`、`hotkey_live_test.go:386`）、`#21`（`live_windows_test.go:625`：`StartHidden` 建窗 ＋ `SetAudioLevel` 喂合成电平，断 `liquidTimerActive`/`DebugTimersAlive`＝**消息泵与计时器**，不开任何设备）。
- **C2 要全局热键在册**（5 枚）：`#8` `#9` `#10` `#13`（`RegisterHotKey`/`UnregisterHotKey` 真 API、`escBorrowProbe` 在 spare id 上试注裸 Esc＝**桌面级独占资源**，`hotkey_live_test.go:57-81`）与 `#2`（真起 `wisp.exe`，按 pid 认领 `WispBallWindow`，`resident_ball_live_228_windows_test.go:60-79`）。`live_guard_windows_test.go:162-176` 的守卫**对任何非自建球窗判红且明写"t.Skip is not an option here"**。
- **C3 要真实输入流／前台／像素命中**（6 枚）：`#12` `#15`（`SendMessageW` 走真 wndproc ＋ `WM_NCHITTEST` ＋ `WindowFromPoint` 读分层窗的 alpha 命中＝**屏幕像素比对**，`interaction_live_test.go:342/370-376`）、`#14`（`GetForegroundWindow`＋`SetCursorPos`＋`mouse_event`，`:307-344`）、`#20`（`SetCursorPos` 回执校验，失败即 `t.Skipf`，`live_windows_test.go:495-502/577`）、`#3` `#6`（**第二个进程**注入裸 Esc，`#6` 还要 `armEscObserver246r2` 抢前台：`resident_task_source_live_246_windows_test.go:99-103`）。

### D｜托盘图标／Shell 通知 —— 0 枚
- 现量：8 文件 `grep -niE 'tray|NotifyIcon|Toast'`＝**零命中**（只命中 "stray" 一词的两处注释）；全仓 `grep -rnE '^func Test.*[Tt]ray'`＝**零命中**。
- 产码**有**托盘与通知：`internal/ball/tray_windows.go`、`cmd/wisp/notify_windows.go`。⇒ 这一格是**覆盖缺口**（今天谁都没测，本机也没测），不是"搬不动"。

### E｜依赖 `frontend/dist` 真页面产物 —— winlive 内 0 枚，但雷在默认档
- winlive 内 `grep -nE 'dist|frontend|gitkeep|go:embed'` 8 文件＝**零命中**。`#1` 只数 webview 子进程，页面建不出来也照样判。
- **真正的 E 族今天已在 CI 的编译射程里**（`//go:build windows`，无 winlive）：`cmd/wisp/panel_resident_windows_test.go:317` 逐字 `t.Skipf("AC#13 has no subject in this tree: the embed resolves no entry...")`；判据＝`internal/panel/assets.go:56` 以 `fs.Stat(tree, "index.html")` 定 `built`，清检出只有 `.gitkeep` ⇒ `Built()=false`（`assets.go:65`、自证用例 `internal/panel/assets_test.go:22`）。同一 `cli` 岗位（`scripts/wisp-cli-tests.sh:113` → `--scope=cli` → `./cmd/wisp/` → `runtests.sh:98`）**任何 SKIP 判红**。
- 默认档同类 skip 现量（`grep -c 't\.Skip'`，均为 `//go:build windows` 或无标签）：`panel_resident_windows_test.go`=3、`panel_host_gate_test.go`=1、`panel_host_windows_test.go`=1、`resident_sink_nail_127_windows_test.go`=1、`leg_sink_nail_131_windows_test.go`=1 ⇒ 合计 **7 处**；台账 `scripts/portable-tests.sh:320-330` 里**没有一行**给 cmd/wisp 记过账。
- ⇒ 结论：**"一推送就红"这一类不是 winlive 带进来的，是默认档带进来的**；给 CI 加 winlive 之前先把这一格判清，否则新步的红会被算到 winlive 头上。

### F｜其实什么都不依赖 —— 1 枚
- `#7` `TestLive246ExitCancelsARunningTask`（`resident_task_source_live_246_windows_test.go:238`）：建 `httptest` golden 服务（`:291`）＋写 `config.toml`＋DPAPI blob（`:302-305`）＋ `buildWispForTest` 起真 `wisp.exe`，断言**只读控制台文本与日志 sink**（`:252/267/272-283`）；不枚举窗、不注键、不比像素。窗起不起得来只影响被启动进程的自我报告，不影响任何一条断言。
- 搬进默认档的代价：`cmd/wisp` 默认档（`//go:build windows`＋无标签，实测 **160 枚**顶层用例）**新增 1 枚**；形状与既有默认档孪生用例同形（`resident_task_source_246_windows_test.go:358` `TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline` 已经用同一个 `bootResidentLegWithEnv`），前置条件（mingw＋`third_party/sherpa-onnx/*.dll`）已由 `ci.yml:441/:448` 备好。⇒ **这一枚今天就该在 CI 跑**。
- 其余 20 枚**都不满足 F**：每一枚都至少读一次 Win32（`IsWindow`/`EnumWindows`/`GetWindowRect`/`RegisterHotKey`/`GetForegroundWindow`/`WindowFromPoint`/进程树窗枚举）或依赖注入键。

## 3. A 族再分一刀 ＋ 必须人眼的

- A 族只有 1 枚（`#1`），且它**只要能建窗＋能数子进程**就能判，⛔ 不需要任何人眼。⇒ **A 族必须人眼＝0 枚**。
- 派单里举的"退净 2 秒"例子＝`#1` 那一支：它是**进程计数**（`TreeWebview`），不是"看着它退净"。它的真难点是负载下的时序抖动（§2 A 那条 1/6 红），不是眼睛。
- 全族**必须人眼＝1 枚**：`#14 TestLiveNeverStealsFocus`（`interaction_live_test.go:282`）。文件里逐字写了边界（`:280-281`）："The half no harness can assert - that a human does not **FEEL** the focus move - **is the owner sign-off, not this test**."⇒ 可机判的那半（`WS_EX_NOACTIVATE`、`WM_MOUSEACTIVATE=MA_NOACTIVATE`、前台窗柄不变、球不当前台）**今天已经在断**；只有"感觉不到"那一半属〔仅本机可量／人眼〕，⛔ 不许伪装成 CI 能测。
- 另有两处**形状性缺件**（不是用例，别算进枚数）：`interaction_live_test.go:13-15` 具名"needs a human hand that no harness can supply (a **second monitor**, the felt experience of focus)"；`#20` 悬停那一支在拿不到物理指针时**主动 Skipf**（`live_windows_test.go:495/577`）并具名把行为判据退回确定性用例 `TestDockHoverPopBackWalksTheRampHome`（`dock_test.go`，默认档、不 Skip）。

## 4. 代价与落点（若把"能进 CI"的那些加进 `.github/workflows/ci.yml`）

**能进 CI 的集合（本腿判）**：A(1) ＋ C1(8) ＋ F(1) ＝ **10 枚**；C2(5) 要独占桌面、C3(6) 要注入/前台/像素 ⇒ 今天**不建议**进 hosted 岗位。
**落点**：现有 job `test-windows`＝`ci.yml:387-527`（`runs-on: windows-latest` 在 `:388`），步序：`:398` winsec 闸门 → `:441` third_party 缓存 → `:448` cgo build smoke → `:455` cmd/wisp CLI → `:477` portable windows → `:513` junction。新步**必须排在 `:448` 之后**（`buildWispForTest` 要 mingw ＋ `third_party/sherpa-onnx/*.dll` 在场，缺任一即 `t.Fatalf`，`cmd/wisp/secret_argv_windows_test.go:179/:187`）。

要动的五处（本腿一字未改）：
1. `ci.yml` `test-windows` 内插一步（`:475` 与 `:477` 之间），或给 `wisp-slo` 自持机加一个 job —— `ci.yml:563-566` 逐字称那台是 **"interactive desktop-capable machine"**，`:591` 已是 `runs-on: [self-hosted, wisp-slo]` ⇒ **hosted 不行时它才是自然落点**。
2. 新脚本 `scripts/winelive-tests.sh`（具名 `-run` 白名单，把上面 10 枚逐名写死）；不能图省事 `go test -tags winlive ./...`（下条 4）。
3. `scripts/portable-tests.sh:320-330` 那张 skip 台账：winlive 里 8 处 `t.Skipf`（`hotkey_live_test.go`=4、`interaction_live_test.go`=2、`live_windows_test.go`=2）若进 `runtests.sh`，**任何一处触发即整步红**（`tools/d22scan/runtests.sh:98`）。要么记账，要么这些用例不能进 hosted。
4. **串行化**：`go test` 默认多包并跑；`internal/ball` 的守卫对**任何**非自建 `WispBallWindow` 判红（`live_guard_windows_test.go:162-176`），而 `#2` 故意在别的进程造一枚，`cmd/wisp` 默认档也会起真 `wisp.exe` 造球窗（`resident_ball_228_windows_test.go:54/:185`）。该文件头 `:22-29` 逐字规定"alone is a rule and not a preference" ⇒ 需要 `-p 1` ＋逐包逐枚跑。
5. 全仓 `if:`／`continue-on-error` 禁用：`ci.yml:6-8` 逐字 "Per D22 mode-6 NO job is configured skippable"＋"Runner platform limits (**no interactive desktop / audio on hosted runners**) are handled by SUBSET choice, never by skipping a job" ⇒ 新步一红就是 PR 闸门红；而同一份注释已经**先验地断言 hosted 没有交互桌面**，与 33-n1 的旁证相互打脸（§6 第 1 条）。

**时长（⚠ 全是估，且估的基数是 owner 本机盘上读数，不是 runner）**：
- `internal/ball` winlive 整包：本机记录 **11.047 s**（`docs/evidence/s1/239-ball-square-fix-v1.md:32`，同发 rc=1 因默认档 CSS 那枚）；指名四枚子集 **2.638–2.707 s**（`245-...-r1.md:49`、`-v1.md:64`）。逐枚 PASS 最大值：`TestBallLiveIdleBorderTransition` 4.41–4.51 s、`TestBallLiveLifecycle` 1.90–2.14 s、`TestLiveSleepingZeroTimerHandles` 1.59–1.62 s、`TestBallLiveAudioLiquidGate` 1.11–1.14 s，其余 <1 s。
- `cmd/wisp` winlive 7 枚逐枚 PASS 记录：`#2` 3.03/3.10/3.96 s、`#3` 6.59/6.72/7.85 s、`#4` 0.07 s、`#5` 0.06 s、`#6` 10.34 s（另两发 FAIL 9.78/10.14 s）、`#7` 20.79/20.84 s；`#1` 无记录 ⇒ 按同族默认档冷启 715.489 ms＋2 s 上界估 **2–3 s**。逐枚合计 ≈ **45 s**（`246-...-v1.md:121` 那发 `TestLive246` 三枚合计 8.08 s）。
- 隐性大头＝**编译**：`#2` `#6` `#7` 各调一次 `buildWispForTest`（每枚重新 `go build ./cmd/wisp` ＋复制 DLL），`#3` `#6` 另建 esclistener 台件（`resident_approval_live_246_windows_test.go:405`）。默认档 `cmd/wisp` 整包本机 148.168/137.440 s（`228-...-v1.md:13`）里已含同种开销可参照。
- **新增总时长估：1.5–4 分钟**（含编译；不含 runner 首次拉 mingw/DLL 的方差）——量级＝"现有 `test-windows` 整 job 的零头"，⛔ 别宣称它是秒级。

## 5. 我可能判错的条目（附"如果错了后果"）

1. **F 只 1 枚**。判据＝"断言里不读任何 Win32 窗口/键/像素事实"。若按更宽口径（"只要起真进程就算依赖桌面"）⇒ F＝0。后果：若反过来把 C1 那 8 枚也算成 F，默认档凭空多 8 枚、且**这 8 枚全都要真建球窗**（`#4` `#5` 走 `startResidentBall`，其余走 `New(Options…)`）⇒ runner 建不出窗就是成排红。口径再宽一档把 `#4/#5` 也当 F，后果同形。
2. **C1/C2/C3 的三分是本腿造的**，仓里没有这个名字。尤其 `#2`：它同时数窗（C）与读进程自报的 `hotkeys live x/4`（C2 味道）。后果：分档名字被当成契约引用后，别人复算会得到不同归属。
3. **21 vs 前一腿的 22**：我把 `live_guard_windows_test.go` 算作 0 枚（它确实只有 helper：`requireQuietBallDesktop`/`newLiveBall`/`subclassTimers`…）。若 22 才对，漏的那枚要么已被删（`git log` 显示这些文件 10-01 仍在动），要么前一腿把 helper 或已不存在的用例计入了。后果：CI 名册少钉一枚 ⇒ "覆盖了整族"是谎。
4. **`#18` `#19` `#21` 用 `StartHidden: true`**（`live_windows_test.go:216/321/629` 现量）⇒ 严格说它们不要求"可见"，只要求"能 `CreateWindowEx`"；而 `#16` `#17` `#20` 是可见＋topmost（`:54/155/479` 无该选项，文件头 `:5-6` 逐字写 "the REAL ball window on the user's desktop (hidden=false, topmost)"）。若 hosted runner 连隐藏窗都建不出，我把它们放进"能进 CI 的 10 枚"就是错的。后果：新岗位第一发即红。
5. **默认档那 7 处 `t.Skip` 不等于 7 枚会 SKIP**：多数是条件支（例如 `panel_host_windows_test.go:794` 只在聚合为空时跳过）。后果：把"CI 已因 E 类红"当现量读，而它只是我从代码推的预测（§6 第 2 条）。
6. **`ci.yml:459` 那句"33 top-level cases"与我的 160 冲突**。我**没有**跑 `--scope=census` 去复核（CPU 禁令）。后果：若 33 才是 `cli` 岗位的真名册，则 `cli` 岗位今天根本没跑 `panel_resident_windows_test.go` 那 11 枚 ⇒ §2-E 的"红已在默认档"推论降级为"待核"。
7. **时长全部来自盘上本机读数**，`#1` 那一枚根本没有记录。后果：CI 排队时间被低估，`test-windows` 从"几分钟"变成十几分钟，触发 `slo-fresh.yml` 那类时限假设。

## 6. 判不动的地方（甲／乙／不做，配现量）

1. **hosted runner 到底有没有交互桌面／能不能接受 `keybd_event`/`SetCursorPos`/前台切换**。
   甲＝按 33-n1 的旁证读数（runner-images issue 自述，非独立复核）＋仓内 `ci.yml:7-8` 的相反先验断言；
   乙＝在 CI 上跑**一枚最小探针**（建一枚 hidden 窗 ＋ `RegisterHotKey` 一枚 spare 组合 ＋ `SetCursorPos` 回执校验，三行断言，⛔ 不碰产品判据）拿到实测；
   不做＝把 21 枚整族塞进 `windows-latest`（正是 33-n1 的 ⛔ 判语）。
   现量：8 文件／21 枚；winlive 内 `t.Skipf` 8 处；守卫文件对污染判红且明写不许 Skip。
2. **`test-windows` 岗位今天是不是绿的**（E 类默认档 skip 是否已经判红）。
   甲＝按 `panel_resident_windows_test.go:317` ＋ `runtests.sh:98` 推"一推即红"；
   乙＝取一次 `test-windows` 的**步级**日志（`gh api` 那条路，本腿没走，因为不在"纯读命令"清单里）；
   不做＝拿"票面/注释说过 33 枚 PASS"当今天的事实。
   现量：默认档 windows 层 skip 站点 7 处；`scripts/portable-tests.sh:320-330` 台账 11 行，**无一行属 `./cmd/wisp/`**。
3. **这族该放 hosted 还是 `wisp-slo` 自持机**。
   甲＝按 `ci.yml:563-566` 的"interactive desktop-capable"字样落 `wisp-slo`；
   乙＝先做第 1 条的最小探针，再决定（若 hosted 能建窗，C1＋A 那 9 枚就不必占自持机）；
   不做＝两头都挂（同一枚用例两份读数，红绿口径会分叉）。
   现量：`runs-on: windows-latest` 两枚（`:388`、`:533`）；`[self-hosted, wisp-slo]` 一枚（`:591`）；cron 只认默认分支（`:24` 逐字），dev 上调度为零发。
4. **`#1` 能不能承受整包负载**（它移出默认档的原因就是 1/6 红）。
   甲＝按 r5 §⑤ 记录，判定"单独跑可以、整包不行"；乙＝进 CI 后拿 20 发步级读数再定；不做＝把它的红当 flake 加重试（r5 逐字禁止：没放宽 2s、没加 `t.Logf` 降级）。
5. **D 族（托盘/通知）零覆盖**：本腿只数了"没有用例"，没数"该有几枚"。甲＝按票 228/244 的"recorded, not executed"口径当已知缺口；乙＝查 `internal/ball/tray_windows.go` 与 `cmd/wisp/notify_windows.go` 的判据该落哪个切片（那要读产码，本腿没展开，属下一枚腿）。现量：全仓 `^func Test.*Tray`＝0 命中。
