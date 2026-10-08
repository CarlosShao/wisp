# 35-a5-census — 票 35 `:75`(c)「winlive 读数进 CI」的形状／成本／可选路普查

- 腿：**只读普查腿**，非任何产码作者；**零 go 命令／零开窗／零 push**。
- 锚：`dffb8062b765a32230d983b82b3c02d425a0c2b3`（branch `dev`；起手时 `git diff --cached --name-only | wc -l` = 0）。
  下面所有行号＝这笔 HEAD 的工作树现量（`file:line`），尺日志全在 `logs/`（每把带 `rc=N`）。
- 被普查的框逐字在 `.scratch/wisp/issues/35-panel-bridge-c17.md:75` 的 **(c) 支**：
  "winlive has **no CI carrier**, so either give that edge a CI-reachable stand-in or state 〔仅本机可量〕 on its face
  and re-run it in every wave that touches the transport"。⛔ 本腿不动票面一个字。

---

## 0. 编排者三条转述的复量结果（两条同、两条被推翻）

| # | 转述 | 我的现量 | 判定 |
|---|---|---|---|
| 1a | `ci.yml` 里 `grep -n winlive` 零命中 | `grep -n winlive .github/workflows/ci.yml` → **0 命中／rc=1**（`logs/a1-ci-winlive.txt`） | 同 |
| 1b | `grep -n -- '-tags'` 也零命中 | `ci.yml` → **0 命中／rc=1**（`logs/a1-ci-tags.txt`）；⚠ 补一句：**`scripts/` 与 `tools/d22scan/runtests.sh` 里有 4 处 `-tags`**，但 3 处是 `scripts/spike/bin/*.exe` 二进制、1 处是 `scripts/spike/run.ps1:47` 的通用透传，`grep -rn 'winlive' scripts/` = **0**（`logs/q2-tags.txt`）⇒ 没有任何 CI 通路把 winlive 编进来 | 同（并补强） |
| 2 | `git log --oneline -S winlive -- .github/workflows/` 为空 | 输出 **0 行／rc=0**（`logs/a2-ci-history.txt`） | 同 |
| 3a | "go test 步骤走 `scripts/runtests.sh` 那形" | **`scripts/runtests.sh` 不存在**：`ls scripts/` 14 项名册（`logs/q2-scripts-ls.txt`）无此件；`find . -name 'runtests*'` 只命中 `.scratch/**` 里的历史日志。真件＝**`tools/d22scan/runtests.sh`**（`runtests.sh:75` = `go test -v -count=1 "$@"`，`:86-88` 数 `--- PASS/FAIL/SKIP`，`ci.yml:658` 注释自陈"runtests.sh:98-102 treats ANY `--- SKIP` as fatal"）。CI 的包级测试实际走 **`scripts/portable-tests.sh --scope=core|census|windows`**（`ci.yml:459/:682/:718`）＋ `scripts/wisp-cli-tests.sh`（`:593`）＋ `scripts/winsec-tests.sh`（`:557`） | **冲突**（前一条普查腿 `35/census-resident-panel/summary.md:60` 也抄了同一句错话） |
| 3b | "如果 CI 不做任何处理，那 CI 里的 cmd/wisp 压根没跑过（DLL 没人备）" | **被推翻**：CI 的 windows 腿**专门为此建了三段承载**——`ci.yml:561-564` `actions/cache@v4 path: third_party key: wisp-third_party-${{ hashFiles('deps.toml') }}` → `ci.yml:566-571` `build.ps1 -Env dev`（fetch-deps＋cgo 链接）→ `ci.yml:573-593` `shell: bash / if: ${{ !cancelled() }} / run: bash scripts/wisp-cli-tests.sh`，脚本内 `wisp-cli-tests.sh:73 dll_dir="$root/third_party/sherpa-onnx"`、`:82` 从 `deps.toml` 抽 `[sherpa-onnx.dll.*]` 名册、`:94` GUARD"pinned DLL(s) absent → 红"、`:109 export PATH="$dll_dir:$PATH"`，再委托 `portable-tests.sh`（`:260 scope=(./cmd/wisp/)`，`:257-259` 写明"Only scripts/wisp-cli-tests.sh may call this"）。⇒ **"cmd/wisp 在 CI 里因 0xc0000135 静默没跑"这格不成立**；成立的是**另一半**：见 §Q2 末的"我没验的那格" | **冲突** |

⚠ 我⛔ 不替编排者裁定"这条缺口要不要报"，但事实形状变了：**缺口不是"cmd/wisp 没跑"，而是"winlive 那 12 枚文件连编译都没人查"**（§Q2 结论）。

---

## Q1 名册：带 `winlive` 构建标签的测试件

尺＝`grep -rln '^//go:build.*winlive' --include=*.go .`（剔 `.scratch` 拷贝）＝**12 枚文件／29 枚 `Test*`**；逐枚表 `logs/q1-tagsummary.tsv`，名册 `logs/q1q3-misc.txt`。
**12 枚的标签行逐字完全相同**（各文件 `:1`）：`//go:build windows && winlive`。

| 文件 | 行数 | `Test*`（名册逐字见 logs） | 需要什么外部条件才不红 |
|---|---|---|---|
| `cmd/wisp/panel_transport_live_35v2_windows_test.go` | 261 | `TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor` ＝ **票 35 `:75`(c) 那一格本尊** | 真桌面＋真窗：文件头 `:13-18` 写明走**产品通路** `PanelManager.bringUp`（`cmd/wisp/panel_host_windows.go:318`）→`installPanelTransport`（`:408`→`:693`→`w.Init(panelPostMessageForwardInit)` 逐字 `:700`），泵＝库自己的 `WebView.Run()`（同常驻进程 `panel_resident_windows.go:299`）**在 lock 死的线程上**（`:188 runtime.LockOSThread()`）。⛔ 不 skip：`:208/:211/:214` 三处 `t.Fatalf`（"35v2 NO LIVE READING…"）。`skip=0／lock=1／deadline 型 2 处`；对照件用 `go test -overlay` 从 `/d/tmp` 跑旧钩子（`:26-27`）。WebView2 Runtime 检查：**0 处**（`grep 'WebView2 Runtime'`＝0）⇒ 装没装都只会被 `Fatalf` 打成红 |
| `cmd/wisp/panel_host_windows_live_test.go` | 74 | `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive` | 真窗＋**进程树计数**：`:49 Fatalf("bringUp for the winlive exit clause")`、`:55 Fatalf("the window is up but our tree holds %d msedgewebview2 process(es)")`；上界**逐字仍是 2 秒**（票 33 r5 判决 §⑤ `docs/evidence/s1/33-panel-host-c27-r5.md:76`） |
| `cmd/wisp/panel_geometry_255_winlive_test.go` | 241 | `TestTicket255RealWindowWidthFollowsTheConfig` | 真窗＋**量像素**：`:108/:121/:125/:154/:162` 五处 `Fatalf`（首扇真建、窗归属线程、窗 rect、dispose 后重建、重建后 rect）。`skip=0` ⇒ 无桌面＝红，不是跳过 |
| `cmd/wisp/resident_approval_live_246_windows_test.go` | 543 | `TestLive246ResidentPipelineRaisesACardAndEscVetoesIt`／`TestLive246ConfirmingCardBorrowsEscVetoesAndReturns`／`TestLive246ExitRefusesAHangingL2Card`／`TestLive246ExitAbandonsAHangingL1Window`／`TestLive246ExitCancelsARunningTask` | 真窗＋**前台焦点必须归它**＋**第二个进程观测台**：`:85-86 Fatalf("ticket 246 RED (desktop state, not our code): an injected … foreground window before Wisp started")`、`:94 Fatalf("this leg could not create the ball window")`、`:143` 注释"(a)+(b): the observer's window **takes the foreground** and injects"⇒ **这一族会抢鼠标/键盘**。运行前提逐字在 `:32`：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -tags winlive ./cmd/wisp …`（**winlive 也要那三枚 DLL**，因为整包在加载期就死）。`deadline 型 4 处` |
| `cmd/wisp/resident_task_source_live_246_windows_test.go` | 589 | （task-source 那一族，名册见 logs） | 同上（同包同载具：无 DLL 就 `0xc0000135`）；`skip=0／deadline 2` |
| `cmd/wisp/resident_ball_live_228_windows_test.go` | 194 | `TestLive228ResidentLegOwnsABallWindowOnTheDesktop` | 跑法逐字 `:20`：`go test -tags winlive ./cmd/wisp -run TestLive228 -v`；`:22` 具名"why 'alone' is a rule and not a preference: **internal/ball's winlive guard**" ⇒ 依赖下面 `live_guard` 那一族 |
| `cmd/wisp/resident_hotkey_live_258_windows_test.go` | 328 | `TestLive258ResidentBootNamesProvenanceAndDefaults`／`…BootTakesConfigHotkeys`／`…OccupiedCombinationNamesTheNewValue`／`…RebindsWithoutRestart` | 全局热键注册 ⇒ **独占热键表**；`skip=4`（本名册里唯一在 cmd/wisp 侧会跳过的）、`lock=1` |
| `internal/ball/live_guard_windows_test.go` | 333 | （守卫台件，无对外 `Test*` 判据，名册见 logs） | ★**"一次只能一扇窗"就是这里定的**：`:125 enumerateBallWindows()` 走 `EnumWindows`＋`GetWindowThreadProcessId`（`:143`），`:158-172 requireQuietBallDesktop`：`:167 Fatalf("foreign Wisp ball window on the desktop: hwnd=%v pid=%d title=%q. Live measurements must own the desktop… Kill pid %d (a leftover balldebug/wisp) and re-run. **t.Skip is not an option here.**")` |
| `internal/ball/hotkey_live_test.go` | 669 | `TestLiveHotkeyRebindEndToEnd`／`TestLiveMuteHotkeyEndToEnd`／`TestLiveHotkeyOccupiedVsNotAttempted`／`TestLiveNeverStealsFocus`／`TestLiveClickSummonsAndDragDoesNot`／`TestLiveConfirmingCancelAndEscReturned` | 真窗＋真热键＋真注入；`skip=4`（部分条件会跳过）；`:TestLiveNeverStealsFocus` 直接断焦点行为 ⇒ 与前台争用互斥 |
| `internal/ball/interaction_live_test.go` | 420 | `TestBallLiveEdgeDock`／`…EdgeDockHover`／`…IdleBorderTransition`／`…PositionPersistence` | 真窗＋鼠标注入（`live_guard_windows_test.go:251 sendMouse`/`:261 sendMessage`）；`skip=2` |
| `internal/ball/live_windows_test.go` | 689 | `TestBallLiveLifecycle`／`…AudioLiquidGate`／`TestLiveSleepingZeroTimerHandles`／`TestLiveTransparentCornerFallsThrough` | 真窗＋GDI 资源计数（`:243 guiResources`、`:333 uiUserObjects`）；`skip=2／deadline 2` |
| `internal/ball/hotkey_cancel_borrow_live_260_test.go` | 110 | `TestLiveCancelBorrowDefaultUnchanged260`／`TestLiveCancelBorrowFollowsConfig260` | 同 hotkey 族（借 Esc 的窗口必须真在） |

**"一次只能开一扇窗"的实测答案＝是，而且不只在文件层**：`internal/ball/live_guard_windows_test.go:162-171` 是**跨进程**判据（按窗口类名 `WispBallWindow` 枚举，见到非己 pid 就 `Fatalf` 并点名要杀的 pid），`grep -rn 'tasklist' <12 枚> = 0` ⇒ 它不查进程表、查的是**桌面窗口**，所以同机并发＝必红，且**红了不会告诉你"是环境问题"以外的任何余地**（cmd/wisp 侧 9 枚里 7 枚 `skip=0`）。

**没答的两格（具名）**：① 逐文件的**超时预算数值**（我只量到 `deadline/timeout` 型语句的**枚数**：transport 2／approval_live 4／task_source_live 2／live_windows 2／live_guard 1，⛔ 没逐枚读那几行去抄秒数）；② 票 35 `:75` 里那两枚毫秒读数（PASS 8.25s／FAIL 58.39s）**我未复跑**，只能引票面。

---

## Q2 CI 今天到底跑不跑得到 cmd/wisp／winlive

**跑得到 cmd/wisp（windows 腿），且专门为那三枚 DLL 备了料**——证据链见 §0-3b，四段皆现量：`ci.yml:561-564`／`:566-571`／`:573-593`／`scripts/wisp-cli-tests.sh:73/:82/:94/:109`。
`ci.yml:577-580` 那段注释自陈的读数是"PASS=33 FAIL=0 SKIP=0 rc=0 once the pinned DLLs are on PATH，without them = ticket 98's `0xc0000135` before a single test runs"＝**这正是编排者本机撞到的那一形，CI 已把它当已知病处理过**。
ubuntu 腿故意不接：`ci.yml:580-583`"On ubuntu the binary builds and runs but 19 of its cases are red … written up on the ticket rather than being wired into this job"；`wisp-cli-tests.sh:62-65` GUARD"this is the windows leg"。

**跑不到 winlive（这才是要紧缺口，而且比"没跑"更硬）**：
- `ci.yml` 里 `-tags` **0 命中**、`scripts/`＋`tools/d22scan/runtests.sh` 里 `-tags` 有 4 处但**无一含 winlive**（`logs/q2-tags.txt`，含 3 枚 `scripts/spike/bin/*.exe` 二进制与 `scripts/spike/run.ps1:47` 的通用透传）。
- ⇒ **winlive 那 12 枚文件／29 枚用例在 CI 里连"能不能编译"都没被检查过**：`runtests.sh:75` 永远不带 tag，`go vet ./...`（`ci.yml:238/:241/:332`）在 ubuntu 上跑，而 `//go:build windows && winlive` 双重排除 ⇒ **改坏这些文件不会让 CI 变红**。这一格比"读数没进 CI"更便宜就能补（ⓓ，见 Q3）。
- 失败会不会让整作业红：**会，而且这台机器上没有任何"可跳过"的口子**——`ci.yml:6`"skippable: no `if: false`, no continue-on-error, no skip flags"、`:34-35`、`:782`、`:869`（D22 mode-6："A skippable job is a job that will one day be skipped"）⇒ 任何被塞进现有 job 的 winlive 步骤，**一旦桌面不安静就是把整个 job 判红**，不是变绿也不是变灰。
- `scripts/slo-check.ps1`（625 行）内 `grep -c 'winlive|tags'` = **0**；`winsec-tests.sh`/`wisp-cli-tests.sh`/`portable-tests*.sh` 亦 0（`logs/q2-dll-staging.txt`）。

**我没验的一格（具名，⛔ 不当结论用）**：`wisp-cli-tests.sh:20`/`ci.yml:577` 那枚 `PASS=33 FAIL=0 SKIP=0` 是**脚本头的自陈**，本腿**没查任何一次真实 run 的日志**（那要 `gh api`，且本腿只测盘上形状）。按规矩（"yaml 里有 ≠ 它跑过"／第 109 条）**"CI 真跑过 cmd/wisp 并绿过"这一条仍待一次 success 读数**。同一条规矩对本普查的 (c) 支同样适用：我量到的只有"yaml 里无＋历史里无＋标签进不来"。

---

## Q3 载体形状 ×4（⛔ 不替编排者选；代价写平列）

GitHub 端**盘上看不见**的东西我单列。今天盘上的 labels 名册（逐字，`gh api repos/CarlosShao/wisp/actions/runners` 本腿现跑＝`logs/q3-runners.txt`）：
**只有一台**：`name=wisp-selfhosted-01|os=Windows|status=online|labels=self-hosted,Windows,X64,wisp-slo`。
`ci.yml:798` 逐字 `runs-on: [self-hosted, wisp-slo]`（⚠ 只 2 枚 label；编排者说的 4 枚是 runner 侧全名册，不是 job 侧选择器）。`ci.yml` 里 `secrets.` **0 命中**、**没有 `permissions:` 块**、**没有 `timeout-minutes`**（`logs/q1q3-misc.txt`）⇒ 今天没有任何作业用密钥；org/repo 默认权限盘上不可见。

| 形状 | 要动的文件 | blast radius（实测成本） | 还需要什么盘上看不到的东西 | 今天有没有证据支持 |
|---|---|---|---|---|
| **ⓐ 塞进现有 job**（`test-windows` 或 `slo-full`） | `.github/workflows/ci.yml`（一步）＋多半要一枚 `scripts/winelive-tests.sh` 包 `PATH=third_party/sherpa-onnx`（同 `wisp-cli-tests.sh:109` 形） | ①**`on:` 是 `pull_request` ＋ `push: branches:[main,dev]`（`ci.yml:10-13`）⇒ 每次提交都开真窗**；②放 `test-windows`（`runs-on: windows-latest`，`ci.yml:506`）**结构不可行**：托管机无交互桌面，而 cmd/wisp 9 枚里 7 枚 `skip=0` 只 `t.Fatalf`（`panel_transport_live_35v2:208/:211/:214`、`panel_host_windows_live:49/:55`、`panel_geometry_255:108-162`）⇒ **永久红，且红因不是产品**；仓里连"无人值守环境能不能真建可见窗"都还挂着〔未定〕（`docs/reports/pending-and-issues.md:10395`）；③放 `slo-full`＝**机主的开发机**（`gh api`＝唯一在线自托管 runner），而 `slo-full` 无 `if:`/无 `needs:`/无 `timeout-minutes`（`logs/q1q3-misc.txt`）⇒ 每次 push 在他屏幕开窗；④争用：`requireQuietBallDesktop`（`live_guard:162-171`）见别家窗就 `Fatalf` 并要人 kill pid，`resident_approval_live_246:143` 还要**抢前台**；⑤`concurrency:`（`ci.yml:40`）只按 sha 分组，**挡不住同机两批发窗**；⑥失败即整作业红（D22 mode-6 禁 `if:`/`continue-on-error`） | runner 必须**开着交互会话**且那天那小时桌面安静；workflow 权限今天没写过（无 `permissions:` ⇒ 走仓库默认，盘上不可见） | **有**：`ci.yml:797-823` 已经存在一台自托管 job 的先例（`env.MINGW64_ROOT: E:\work\base\msys64\mingw64\bin` ＋ `:801-809` 注释逐字记着"这台的 PATH 里没有 msys64，runner 作业环境取不到 gcc"＝自托管机的环境坑已被踩过一次并写死） |
| **ⓑ 独立 job ＋ 独立 label**（只有某台带那 label 才跑） | `ci.yml` 新增 job（`runs-on: [self-hosted, wisp-winelive]`）；⚠ **label 不在 yaml 里生效**，要人在机器上改 runner 配置 | ①**今天只有一台自托管 runner**（现量 `logs/q3-runners.txt`）⇒ 新 label 只能贴到**同一台＝机主屏幕**，开窗成本与 ⓐ 一模一样，只是多了"要有人去 runner 目录改 `.env` 并重启服务"；②**label 并不能把它和 push 解耦**：只要作业还在 `ci.yml` 的 `on: push` 下，触发就照旧 ⇒ 想只在特定时机跑就**必须写 `if:`**，而 D22 mode-6 逐字禁（`ci.yml:34-35`"No `if:`, no continue-on-error, no skip flag anywhere in this block"、`:782`、`:869`）⇒ **ⓑ 走不通除非动契约级禁令**；③⛔ 编排者转述里的 label 名册（4 枚）是 runner 侧，`slo-full` 只用了 2 枚 ⇒ 加 label 属**机器侧变更，不在盘上 diff 里**，后续程看不见 | GitHub 端：runner 注册配置（`.Env`/服务账号）＋**谁能在那台机器上贴 label**；仓库是否允许 job 独占该 runner（无 `permissions:`/无 environment 保护，盘上看不出） | **有（反向）**：`slo-fresh.yml:7-12` 逐字写明"slo-full … can only run on one machine: the self-hosted wisp-slo runner **that is this laptop** (ci.yml, A103)"＝仓里已经承认单点机器 |
| **ⓒ schedule 触发的独立 workflow**（不挂 push） | 新文件 `.github/workflows/wisp-winelive.yml` ＋（几乎必然）一枚 `scripts/*-freshness` 反形件＋票 35 `:75` 的凭据句 | ①**push 成本＝0**（这正是机主屏幕的价钱归零那一支）；②⚠ **cron 只对默认分支生效**：`ci.yml:24-26` 逐字"GitHub evaluates cron ONLY against the config on the DEFAULT BRANCH (main). Until this file reaches main this schedule fires zero times — on dev it is inert by design" ⇒ **在 dev 上它今天就是零发**，这不是可以靠 push 修的 bug；③⚠ **定时运行是 best-effort**：`ci.yml:27-29`（峰值可延迟几分钟），且"一份没人触发的 job 就是 ticket 85/71 那种病"＝`slo-fresh.yml:9-18` 的原话 ⇒ **新 workflow 必须自带"它还活着"的钉子**，否则后续程会把"有这份文件"读成"跑过"（第 109 条）；④时刻选择权在机主：现成先例把窗排在本机 **03:37 本地**（`ci.yml:19-21`"at 19:37 UTC … when the review fleet is asleep … one interference-free shot"）＋`slo-fresh.yml:47-51` 自己的钟是 `cron: '23 */6 * * *'`、`runs-on: ubuntu-latest`（`:62`）＝**先例走托管机，因为它是"问 GitHub 要记录"而不是开窗**；开窗那一支只能落自托管机 ⇒ 要求**那个时刻屏幕开着且没锁、且没有别家的球窗**（`live_guard:167`）；⑤失败即红＝每天/每六小时可能给机主发一封红邮件；⑥`schedule` 会跑**整个 workflow 的全部 job**（`ci.yml:30-33`"A scheduled run executes ALL SIX jobs … narrowing it would need an `if:`, and D22 mode-6 bans that"）⇒ 要么新文件里**只放一个 job**，要么付托管分钟 | GitHub 端：**默认分支必须是 main**（否则零发）、`on.schedule` 只在默认分支配置里被读、runner 在触发那一刻在线（`status=online` 是此刻读数，不是保证）、红邮件/通知策略 | **有**：`ci.yml:36-38`（`schedule: - cron: '37 19 * * *'` ＋ `workflow_dispatch:`）与 `slo-fresh.yml:46-52`（独立钟＋手动入口）＝**两种先例都已经在仓里** |
| **ⓓ 只补"编译面"的 CI 载体**（⛔ 不开窗；本腿新提，形状由 Q2 缺口直接推出） | 一步：在 `test-windows` 里加 `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`（或 `-C` 形的包名册件）＋可选一枚 `scripts/winelive-compile.sh` 记包名册 | ①**零开窗、零争用、零鼠标**——vet 不生成也不加载测试 exe ⇒ 那三枚 DLL 用不上（对照 `wisp-cli-tests.sh:7`：病在**加载期**）；②**买到的只有一格**：`//go:build windows && winlive` 这 12 枚文件今天**没有任何 CI 通路编译过**，ⓓ 让它们改坏就红；③⛔ **它不满足票 35 `:75`(c) 的判据**（(c) 要的是"页面自己的 `chrome.webview.postMessage` 真到 Go 门"的凭据或就地〔仅本机可量〕），⛔ 不许被后续程读成"载体已建"；④代价＝托管 job 多一步（`windows-latest` 上要有 cgo 链，`ci.yml:566-571` 已证明该 job 有 mingw；自托管那台则需 `MINGW64_ROOT`，`ci.yml:810`） | GitHub 端：不需要新东西（不新增 label／不新增 runner／不动触发） | **有**：`go vet ./...` 已是既有步骤形（`ci.yml:238/:241/:332`），加 tag 属**同一形状**的收窄 |

⛔ 本表不选路。ⓐⓑⓒ 三类的**共同硬事实**：票 35 `:75`(c) 那一枚 `TestLive35v2…` 是**真窗＋真产品通路＋抢前台**的读数，而它同族的 `requireQuietBallDesktop` 把"桌面归我们独占"写成了 `Fatalf` 判据（`t.Skip is not an option`）⇒ **任何"进 CI"的形状都同时在买"CI 每次跑都要机主让出鼠标"这一条**；不想要这条就只有 ⓓ（换不到 (c) 的凭据）或 Q5 的〔仅本机可量〕。

---

## Q4 此刻桌面争用尺（只读；⛔ 本腿没开窗、没跑 go、没动鼠标）

`logs/q4-tasklist.txt`（`tasklist //NH`，现量）：**总进程 343**；`wisp`=**0**；`msedgewebview2.exe`=**18**；`node`=9；`*test*.exe`=**0**；`conhost`=21。
- 好的一半：**没有 wisp 进程、没有任何测试 exe 在飞** ⇒ `enumerateBallWindows` 那条"别家 `WispBallWindow`"判据（`internal/ball/live_guard_windows_test.go:125-171`）**此刻不会因为残留 wisp 而红**（注：该判据按窗口类名枚举，进程表看不见类名 ⇒ 严格说我只排除到"没有 wisp.exe"这一层）。
- 坏的一半：**此刻盘上有 18 个 `msedgewebview2.exe`**（属别的应用，不是 wisp），而 `resident_approval_live_246_windows_test.go:85-86` 的红条件是"Wisp 起窗**之前**就有别家的前台窗"、`:143` 那段是"观测台自己抢前台再注入按键" ⇒ **真跑一轮这族＝今天这台机器上一定会挪走某个正在前台的窗口**。⚠ 这 18 枚是谁的、会不会被抢焦点，**只有机主知道**。
- **屏幕锁不锁、机主此刻在不在机器前、那个点他能不能让出鼠标＝只有机主知道**（本腿无法从盘上判断，⛔ 不猜）。
- `gh api` 现量 `wisp-selfhosted-01 … status=online` ⇒ runner 在线；**但"runner 在线"≠"有交互桌面会话"**，这格同样只有机主知道。

---

## Q5 〔仅本机可量〕的既有记账形状（可直接复用的形）

★**先例出处（票 33 那一次，四要素都在）**：`docs/evidence/s1/33-panel-host-c27-r5.md:74-79`，判决正文在 `.scratch/wisp/issues/33-panel-host-c27.md:308`（"裁定＝乙：那一支移出默认档、进 `winlive`…上界仍写 2 秒"）与台账 `docs/reports/pending-and-issues.md:10282`。
`r5.md:78` 逐字：

> ⚠ 代价逐字写进文件与本判决：**`winlive` 在 CI 零岗位 ⇒ 这一支从此〔仅本机可量、CI 永看不见〕**；引用它的任何表必须带这句（撤销口令＝「33 退净断言回默认档」）。

拆成**可复用的四要素**（后续程要照抄的是这四件，不是一句话）：
1. **写两处，不写一处**：同一句既进**判决件**，也进**测试文件自己的头注释**（`r5.md:78` 的"逐字写进文件与本判决"；对照现物＝`cmd/wisp/panel_host_windows_test.go:6`、`:605/:609` 与 `cmd/wisp/panel_resident_windows_test.go` 族里都带着"winlive has no CI job"这句）。
2. **带撤销口令**（`r5.md:78` 的"撤销口令＝「…」"）——让"以后想搬回默认档"有一个具名动作，而不是重新论证一遍。
3. **欠账具名**：`r5.md:79` 逐字写"本腿**没有**跑 `winlive` 档…读数只有 `33-r4` 的（我未复跑，出处＝`33-panel-host-c27-r4.md` §②）…这条欠账也记在 §⑦ 第 10 条" ⇒ **"谁的读数"必须点名到哪一腿、哪一节，并声明本腿未复跑**（正是第 108 条"死腿的 logs 只是它取过"的处置形）。
4. **凡引用必带这句**（"引用它的任何表必须带这句"）⇒ 票面／矩阵／台账里那一格被引用时必须重复标注，防止被读成"CI 里有"。

**反形（防被后续程读成"做过了"）——仓里已有三种，缺一都算没记牢**：
- 票面自带的复跑义务：`issues/35-panel-bridge-c17.md:75`(c) 逐字 "**re-run it in every wave that touches the transport**" ⇒ 〔仅本机可量〕**不是一次性注记，是每波的欠账**，需要台账里逐波留名（票 33 的先例是一行欠账，波次靠 `HANDOVER`/台账推）。
- 派单级禁跑句的既有写法（两枚，形可直接抄）：`issues/273-…md:25`"⛔ **`winlive` 未批**：任何'真开窗口…'"、`issues/274-…md:34`"⛔ **不许动 `winlive`**：任何'真开面…'" ⇒ 任何触碰传输层的写腿派单**都要带上这两条形之一**，否则"仅本机可量"会被下一波当成"它已被 CI 覆盖"。
- 仪器级反形：`tools/d22scan/runtests.sh:22-27`（零 `--- PASS`/`--- FAIL` ＝ fatal、**任何 `--- SKIP` ＝ fatal**、`go test` 退码原样传）＋`ci.yml:658` 的注释 ⇒ 一旦有人把 (c) 改成"没桌面就 `t.Skip`"来骗 CI，**这道严格跑器自己就红**——这是"仅本机可量"不会悄悄变成假绿的既有理由，⛔ 不许绕过它去加 skip。

---

## 成本表（最关键三行，其余在上）

| 成本项 | 现量 |
|---|---|
| **每次 push 在机主屏幕开窗＋抢前台** | `ci.yml:10-13` push 触发 [main,dev]；唯一自托管 runner＝机主这台（`gh api` 现量 `wisp-selfhosted-01 online`，labels `self-hosted,Windows,X64,wisp-slo`）；`slo-full` 无 `if:`（`ci.yml:797-798`）⇒ ⓐ/ⓑ 都躲不掉这一条 |
| **红了不是产品红，是环境红** | cmd/wisp 的 9 枚 winlive 用例里 **7 枚 `t.Skip`＝0**、只有 `Fatalf`（如 `panel_transport_live_35v2:208/:211/:214`）；跨进程独占判据 `internal/ball/live_guard_windows_test.go:162-171`"t.Skip is not an option here"；D22 mode-6 禁 `continue-on-error`（`ci.yml:34-35/:869`） |
| **今天的真实缺口不是"cmd/wisp 没跑"，而是"winlive 连编译都没人查"** | CI 已为 cmd/wisp 备齐 DLL 三段（`ci.yml:561-564/:566-571/:573-593`＋`wisp-cli-tests.sh:73/:82/:94/:109`）；但 `-tags` 在 `ci.yml` 0 命中、`scripts/`＋`runtests.sh` 的 4 处 `-tags` 无一含 winlive ⇒ 12 枚文件／29 枚用例**改坏不红**（ⓓ 一步可补，零开窗，⛔ 不满足 (c) 凭据） |
