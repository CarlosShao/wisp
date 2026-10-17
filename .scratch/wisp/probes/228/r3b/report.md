# 228-r3b 读数交件：五枚 AC#11 停机用例跑到终态颜色 + AC#11 ⓑ 那把尺

任务书（编排者派单逐字存档）：`D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\orch\20261017\txt\228r3b-dispatch.txt`（46 行／11088 字节，本腿已 Read 全文）。
起手锚（现跑 `git log -1`）：`ae732436  Sat Oct 17 12:09:40 2026 +0800`（分支 dev）。
票面未勾尺（现跑 `grep -c '^[[:space:]]*- \[ \]'`）：**11** 枚未勾（与派单点名一致，非抄录）。
台面：所有颜色在**仓外导出台树** `/d/wisp228r3b`（`D:\wisp228r3b`→8.3 短名 `D:\WISP22~1`）里跑，⛔ 未在共享工作树跑 `./cmd/wisp`；取树尺=`git archive ae732436 <八件套> docs`＋手工三件（build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist）；取出后 `.go` 数=**758**（证非空）。两把临时根台面成对：c=默认 C 卷、d=`TMP/TEMP/TMPDIR=/d/wisp228r3b/tdst`；`GOTMPDIR=/d/wisp228r3b/tmp` 先建。`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 全程带上（防 `0xc0000135`）。

## 1. 那五枚逐枚颜色（起手名册=交件名册；estate 逐枚=定向 -run）

名册复认（起手侧 `git grep -h '^func Test' HEAD -- 那两枚件` ⇄ 交件侧定向 `-run`）五枚逐字一致：

| 顶层用例 | 判据格 | surface c | surface d |
|---|---|---|---|
| TestAC11StopBitIsBuiltAtBoot | ⓐ 第一句（位在建） | PASS | PASS |
| TestAC11TrayExitExecutorIsMountedAtTheBallHost | 派单额外钉（挂到球宿主） | PASS | PASS |
| TestAC11TrayExitTriggerWalksTheTenSteps | **ⓑ-1**（进程内读满十步） | PASS | PASS |
| TestAC11StopRequestDuringTheBootWindowIsNotSwallowed | ⓐ 空窗那半 | PASS | PASS |
| TestLive228R3TrayExitStopsTheShippedProcessWithCodeZero | **ⓑ-2**（真机） | **FAIL** | **FAIL** |

定向四枚（不含 winlive）台面成对四数：c `RUN=4 PASS_TOP=4 FAIL_TOP=0 SKIP=0`；d 同。ⓑ-1 逐字 READING（estate-targeted-c.log:20）：`trail=scheduler-close > stop-hotkey-kws > cancel-task-roots > stop-audio > release-speech-sessions > destroy-panel-webview > flush-logs-close-db > free-os-memory > close-job-object > exit`，`reason="stop request (tray Exit item selected in this process)"`＝十步真走满、顺序 1..10、点名谁要求退出。⛔ 它不声称"托盘被点过"（无窗口腿）。

## 2. 整包四数（`go test ./cmd/wisp/ -count=1 -timeout 30m -v`，导出台树）

- surface c：`RUN=208 PASS_TOP=149 FAIL_TOP=4 SKIP=0`，`go-test-rc=1`，包尾 `FAIL ... 197.781s`。
- surface d：`RUN=203 PASS_TOP=146 FAIL_TOP=2 SKIP=0`，`go-test-rc=1`，包尾 `FAIL ... 151.508s`。
- **真窗区整进程死亡（派单点名的那一格，两发不同处＝非确定性）**：c 最后一个 `=== RUN` 是 `TestPanelThreadIsSTAAndExitsCleanly`（无 `--- 结果`行，包即回 FAIL）；d 最后一个是 `TestAC13BringUpSurvivesAReusedThreadQuit`。两枚都是 panel/WebView2 真窗用例 ⇒ `resident_*` estate 用例（文件名排在 `panel_*` 之后）在整包里**从不执行**；所以 estate 的真实颜色取自第 1 节的定向 `-run`，整包四数只是台面终态，⛔ 当"全绿"也⛔ 当"我的字节把包打红"。

### 红名册逐枚归属（必答 ii）
| 红/死 | 出处 | 归谁 |
|---|---|---|
| TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (40.96s，仅 c) | approval L2 二卡时序 | 别格（票 224/226 地界），非 estate |
| TestAC2SealNoticeLandsInTheRunLegLogFile (c+d) | run 腿控制台/日志 | 本机 run-leg，非 estate |
| TestAC14GoSideEvalPushReachesThePage (仅 c) | panel WebView2 Eval push | 别腿真窗（票 33），非 estate |
| **TestCleanCheckoutBuilds_AC11 (c+d)** | `git rev-parse --show-toplevel: exit status 128` | **台面故障**：导出台无 `.git`，非 AC#11 estate，⛔ 计入判据 |
| 整进程死亡（见上） | panel/WebView2 真窗 | 本机真窗区／别腿在飞，非 estate |

⇒ estate AC#11（ⓐ/ⓑ-1 四枚）**零枚本案红**；ⓑ-2 红＝缺真窗交互世界（见 §3）。

## 3. AC#11 ⓑ 那把尺（判据原文在票面 :91-102，本腿逐字读，⛔ 抄转述）

ⓑ 已由编排者拆成两枚分开的钉（票面 :97-100，A495）：ⓑ-1 进程内读满 10 条 StepRecord 断顺序；ⓑ-2 真机只钉"退出码 0＋审计 install/托盘/拆球/D38(e) 尾巴"、⛔ 不读 1..10 全序。
- **ⓑ-1：绿（两台面），且真读到 1..10**（§1 READING 行）。
- **ⓑ-2：红，两台面同拦于一行**＝`resident_tray_stop_live_228_windows_test.go:184`。逐字红句（surface c，live-targeted-c.log:4）："AC#11 ⓑ-2 RED: no popup menu window is owned by pid 42864 after the tray callback message; the break is between the message and internal/ball's showMenu on the ui-sta thread."；surface d 同句、pid 为 42040。
  **是哪一行拦的＋最小说明**：出货常驻子进程（bootResidentLeg，pid 42864/42040）确实起来、`ball window is up`（hwnd 被该 pid 持有、isWindow 通过第 2 步）、第 3 步 `PostMessageW(hwnd, WM_APP+0x202, trayUID, WM_RBUTTONUP)` 返回真值（非 0，未在第 174 行拦），但第 4 步轮询 200 次内**没有 #32768 弹出菜单窗口**归该 pid ⇒ 断点在"消息→`internal/ball` 的 `showMenu`/`TrackPopupMenu`（ui-sta 线程）"这一段：这台机器当下开不出那个可被注入鼠标点选最后一行的交互托盘菜单世界。子进程热键全 `occupied by another program`（本机活实例 PID 21892 握全局热键＋单实例锁，本腿⛔ 未杀它、⛔ 未为它动断言；本案不依赖热键）。⇒ ⓑ-2 这一格＝**跑过、红、具名到行**，非"没跑写成跑过"，归**编排者约真窗**，⛔ 用"缺世界"当通过。

## 4. 突变（牙，必答 i）——只打在导出台树字节上，逐发已还原并报还原哈希

两发落点/用例不同（都在 `internal/proc/boot_windows.go`，不同断言行、驱动不同用例）：
- **M1** 断环路 stop 检查：`stopSignal := rt.StopSignal()`(行 185)→置 nil 通道 ⇒ **ⓑ-1 用例 TestAC11TrayExitTriggerWalksTheTenSteps 点名红**于 `resident_tray_stop_228_windows_test.go:310`（"RunEventLoop never returned...the clause that reads that check"，5s 超时）。⇒ ⓑ-1 **有牙，本腿实测**。
- **M2** 令 Boot 不建停机位：`rt.stopReq = newStopRequest()`(行 111)→注释掉 ⇒ **ⓐ-first 用例 TestAC11StopBitIsBuiltAtBoot 点名红**于 `:137`（"no stop-request bit (StopSignal() answers nil)"）。⇒ ⓐ-first **有牙，本腿实测**。
- 还原哈希：两发 pristine＝`3e2e83fab184ce91640ad7dbc06f8dd116baffe325bd47aa46af2527aad73fe7`（＝HEAD `git show ae732436:internal/proc/boot_windows.go` 的 blob）。⚠ 一处如实记账：两发的 EXIT-trap 曾各自当场 echo `RESTORED sha=3e2e83fa`，但事后核查导出台**一度同时残留两处突变**（跨调用 trap 竞态）⇒ 本腿已**强制还原并逐字 diff=IDENTICAL**，随后在确认 pristine（跑前/跑后同哈希 3e2e83fa）的树上**复跑四枚 estate=4/4/0/0 全绿**（estate-reverify-c.log）。仓内跟踪码 `cmd/**`、`internal/**` 全程 `git status` 干净（见 §5 尺）。
- **本腿未验**：TestAC11TrayExitExecutorIsMountedAtTheBallHost、TestAC11StopRequestDuringTheBootWindowIsNotSwallowed 两枚**未做反形突变** ⇒ 不写它们有牙、也不写它们没牙（⛔ 盖章）。

## 5. 门禁 rc（逐枚）
- `go vet ./cmd/wisp/`（**导出台树**）＝ **rc=0**。
- `sh scripts/d22scan.sh`（**共享工作树**）＝ **rc=0**（clean；examined internal/=241+cmd/=42 production，ban #8 emoji 全段扫；本腿零 .go 跟踪件）。
- `sh scripts/check-path-length-budget.sh`（**共享工作树**）＝ **rc=0**（VERDICT GREEN；denominator tracked=14986、over-budget=60、全部 roster 覆盖）。
- `gofumpt -l`：本腿**零 .go 跟踪件** ⇒ 同尺空列表，⛔ 当"过了"（无凭据）。
- 产码零字节尺（下 §6 随附）。
