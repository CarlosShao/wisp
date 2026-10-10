# 302-a1 `40-census.md` — 票 302 `AC#0` 只读代价普查（三张表＋一枚必答）

腿＝`302-a1`（只读普查）。射程＝**只有 `AC#0`**（⛔ 碰 `AC#1`..`AC#4`，那几格是落地腿与验收腿的）。
⛔ 本轮产码＝**0 字节**。⛔ 编译面命令（`go build`/`go test`/`go list`/`go vet`/`go run`）＝**0 次**（`go env` 本腿也**没跑**，所以无从自报例外）。
⛔ 开窗、⛔ 装东西、⛔ push。写面＝只新建 `.scratch/wisp/probes/302/a1/**` 里的 `.md`/`.txt`（⛔ `.sh`/`.ps1`/`.go`/`.out`）。

对象层规矩：凡引 `scripts/**`、`.github/**`、`cmd/**`、`internal/**` 一律 `git show HEAD:<path>`／`git grep … HEAD`（**blob 层**）。
CI 字节那几件是 `.scratch/**` 里的归档件，读的是**工作树里的归档文件本体**（它们不是 HEAD 的产码面，逐枚列在 §0）。

---

## 0. 读数台（每把自己尺的射程与 rc 落在这张表里，⛔ 裸数）

| # | 尺 | 射程 | 载体 | rc |
|---|---|---|---|---|
| S1 | 红句词面（CI 侧） | `grep -E "no report .* from the page"` | 工作树归档件 `probes/301/orch/logs/{ci-after-cli-block,ci-baseline-cli-block,ci-windows-block,ci-baseline-windows-block}.txt` ＋两枚整 job 原始日志 | rc=0（命中 3/0/0/0；整 job 3/0） |
| S1′ | 红句词面（代码侧，**整族**） | `git grep -n "awaitReport(" HEAD -- '*.go'` 剥 `.scratch/probes` | HEAD blob，全仓 `*.go` | rc=0（定义 1＋调用 3） |
| S2 | tag 尺·甲 | `git grep -nE "^//go:build.*winlive" HEAD` 全树 | HEAD blob | rc=0（14 命中，其中 2 枚在 `.scratch`，**Go 文件 12**） |
| S2′ | tag 尺·乙（独立第二把） | 逐文件 `git show HEAD:<p> \| head -3` 扫 `go:build.*winlive`，射程＝`git ls-tree -r HEAD -- cmd/wisp internal` 的 **640 枚 .go 全量** | HEAD blob | rc=0（**12**，与 S2 同数） |
| S3 | ledger 尺 | `ledger=(` 段逐行 `grep -c '^    "Test'` | HEAD blob `scripts/portable-tests.sh`（blob 行 590–602） | rc=0（**11 行**） |
| S4 | `-skip` 静默尺 | 归档 CI 字节里 `^--- SKIP: <ledger 名>` 计数 | 工作树 `ci-windows-block.txt`／`ci-after-cli-block.txt` | rc=0（7 枚 ledger 名 ⇒ `--- SKIP` **0 命中**） |
| S5 | SKIP-判红门 | `awk` 取 blob 行 88–110 | HEAD blob `tools/d22scan/runtests.sh` | rc=0（`:98` 判 `skipped≠0` ⇒ `:102 exit 1`） |
| S6 | winlive 编译门覆盖面 | `grep -n "go vet -tags winlive"` | HEAD blob `.github/workflows/ci.yml` | rc=0（`:672` 逐字含 `./cmd/wisp/`） |
| S7 | pin 尺 | `awk` 取 blob 行 206–214 | HEAD blob `scripts/portable-tests.sh` | rc=0（`cli_pin`＝**1 行 import path，零枚测试名**） |
| S8 | 逐枚色差 | `grep -o -- "--- [A-Z]*: <名>"` 两发切片 | 工作树归档件 | rc=0（§表①·差集那节） |
| S9 | 嫌疑面尺 | `git log --oneline cc315261..bcd0a543 -- <4 files>` | commit 图 | rc=0（**5 枚**） |

---

## 表① 全名册——三把尺分开各报一次，⛔ 混用

### (a) 红句词面尺

**CI 那一半（S1）**。词面＝`no report … from the page`。四个切片＋两枚整 job 原始日志各扫一次：

| 切片／整 job | 命中枚数 |
|---|---|
| `ci-after-cli-block.txt`（改后 `--scope=cli`） | **3** |
| `ci-baseline-cli-block.txt`（基线同档） | 0 |
| `ci-windows-block.txt`（改后 `--scope=windows`） | 0 |
| `ci-baseline-windows-block.txt` | 0 |
| `job-114162576251.log`（改后整 test-windows job） | **3** |
| `job-114069831344-baseline.log`（基线整 job） | 0 |

⇒ 名册 **3 枚**，逐枚（＝派单预告的那 3 枚，逐枚确认，⛔ 多⛔ 少）：

1. `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`
2. `TestAC14AwaitedBindingReplyReachesThePage`
3. `TestAC14GoSideEvalPushReachesThePage`

红句逐字（三枚同形，出处尺＝blob `cmd/wisp/panel_resident_windows_test.go:226` 那句 `t.Fatalf`）：
```
no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all). AC#13 cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
```
另两枚＝同句把 `%q`/`%s` 换成 `"ac14r-0"`／`AC#14's reply hop` 与 `"ac14-push"`／`AC#14's push hop`。

**代码那一半（S1′）＝整族，⛔ 抽样**：`awaitReport` 在 HEAD 全树 `*.go` 里**定义 1 处、调用 3 处**，三处全在同一枚文件（blob 行 :325 / :821 / :866），且逐枚落进上面那 3 枚函数。
⇒ **尺 (a) 的整族＝3 枚，且与 CI 那 3 枚同名同枚**（这不是巧合：那 3 枚是这 3 处调用的唯一消费者）。`.scratch/wisp/probes/**` 里的命中（`33/p1/receipt/main.go:290` 等）**⛔ 计入名册**（⛔ 产码面、⛔ 任何档扫的射程）。

**宽一把（同一射程、词面换成"要浏览器真回话"）**——派单让我"看⛔ 还有别的"，答案是**有，2＋1 枚**：

4. `TestPanelHostRealWindowHopAndLifecycle`（`cmd/wisp/panel_host_windows_test.go:614`）——红句逐字
   `cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window`（`:662`）。
   **⛔ 被尺 (a) 抓到**（词面⛔ 含 `no report`），**⛔ 被票面名册抓到**：它基线与改后**两发都红** ⇒ 名集合差尺 `comm -13` 天生看不见它。判据要的仍是**页面那一跳**（一次真 round trip）。
5. `TestPanelHostLatencyPercentilesAC2`（同文件 `:976`）——**派生**一枚：红句逐字
   `no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate`（`:985`）。
   它⛔ 自己建窗，它吃第 4 枚的样本 ⇒ **第 4 枚被 tag 走，第 5 枚必变 `--- SKIP`**（见 §表③·甲 的副作用与 §丙 的形状 3）。

**相邻但⛔ 同族的一枚（只登记、本格⛔ 收）**：
6. `TestAC247LiveMicrophoneLevelsReachTheBallSeam`（`cmd/wisp/resident_audio_247_live_windows_test.go:130`）——要**真麦克风**（`AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0`），⛔ 要窗⛔ 要页面回执。它产 `--- SKIP`（⛔ 静默）⇒ 撞 `runtests.sh:98→:102`。

**`core`／`winsec` 档那一格——证不到，缺的是哪枚读数（明写）**：
尺 (a) 需要 CI 日志字节。盘上只有 `test-windows` 那两发的切片／整 job（`probes/301/orch/logs/` 的 24 枚件里**没有** `test-core`（ubuntu）那一发的任何字节，也**没有** `winsec` 档切片）。
⇒ **尺 (a) 在 `core`／`winsec` 档⛔ 有读数**（⛔ 说"零枚"）。本腿只能给**算出来的**那一半：S1′ 在 HEAD 全树 `internal/**` 里 `awaitReport` **零命中** ⇒ `core` 档里**没有**"等页面回执"这类判据的调用点。⚠ 这是**代码侧词面尺**的结论，⛔ 是 CI 实测，引它时必须带这句。

### (b) tag 尺（`//go:build winlive`）

两把独立尺（S2 全树 HEAD blob／S2′ 逐文件 `head -3`、射程＝`cmd/wisp`＋`internal` 的 640 枚 `.go` **全量，⛔ 抽样**）**同数＝12**。逐枚（首行逐字都是 `//go:build windows && winlive`）：

`cmd/wisp/`（**7**）：
`panel_geometry_255_winlive_test.go`(1 枚 func) · `panel_host_windows_live_test.go`(1) · `panel_transport_live_35v2_windows_test.go`(1) · `resident_approval_live_246_windows_test.go`(3) · `resident_ball_live_228_windows_test.go`(1) · `resident_hotkey_live_258_windows_test.go`(4) · `resident_task_source_live_246_windows_test.go`(2)
`internal/ball/`（**5**）：
`hotkey_cancel_borrow_live_260_test.go`(2) · `hotkey_live_test.go`(4) · `interaction_live_test.go`(4) · `live_guard_windows_test.go`(0，守卫件) · `live_windows_test.go`(6)
⇒ **12 枚文件／29 枚 `func Test`**。

★**对本票这 3 枚，尺 (b) 的答案是"零枚"**：三枚的出处文件首行逐字都是单 tag（⛔ winlive）——
`cmd/wisp/panel_resident_windows_test.go`／`cmd/wisp/panel_pageover_33r10_windows_test.go`／`cmd/wisp/panel_host_windows_test.go` 三枚 blob 第 1 行＝`//go:build windows`。
⇒ 尺 (b) 给本票的唯一产物＝**"甲形该搬去的那个档长什么样"**（12 枚文件、7 枚在 `cmd/wisp`、编译门在 `ci.yml:672`）。

★**派单点名要的那枚核对：`ci.yml` 的 `go vet -tags winlive` 覆盖⛔ 覆盖得到 `cmd/wisp` 的这三枚？——覆盖得到。** 逐字（S6，blob 行 672）：
```
run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/
```
⇒ 三枚搬进 winlive 后**仍进编译面**（⛔ 变成"既⛔ 跑也⛔ 编"的隐形）。⚠ 但那一半只有 `go vet`：AC#2 还要 `go vet … rc=0` 与 `-list` 名册含那三枚**两枚读数**，本腿⛔ 能采（⛔ go 编译面）⇒ **欠读数＝那两枚 rc／名册**，归落地腿。
⚠ 顺带一枚过期读数：33-n1 引的 `ci.yml:388 runs-on: windows-latest` 在 HEAD `6414a4bb` 已腐烂 ⇒ 本腿现量＝`:534`（`test-windows` 那 job 跨 `:533`–`:817`；`--scope=windows` 那步在 `:797`，cli 那步 `scripts/wisp-cli-tests.sh` 在 `:621`，census 在 `:674`）。

★**与〔⛔ 我复跑，出处 `33-n1`〕打一次架**：其 §1 补充读数写 `grep -rln "go:build windows && winlive" --include=*.go .` ⇒ **8 枚文件**。本腿在 HEAD `6414a4bb` 量到 **12**，且 `ci.yml:625-626` 本仓自己的注释逐字写着 `7 under cmd/wisp, 5 under internal/ball` 与 `all twelve` ⇒ **两把独立尺＋一处仓内自述都指 12**。⇒ 8 是 10-01 锚 `a5357fa1` 的**工作树**过期读数。本票一律按 **12** 计。（⚠ 工具口径具名：我另跑了一把窄尺＝行首锚定＋含 `&` 字样的模式，它返回 **0**——**那是本会话 shell 对 `&` 的处理造成的假零**（同一形状换成 `-F` 字面量就有命中），⛔ 是仓里的事实，本件⛔ 用它得数。第三把可信尺＝逐枚 `git show HEAD:<p> | head -1 | grep -q windows`，对 S2′ 那 12 枚跑 ⇒ **12/12** 首行同时含 `windows` 与 `winlive`。）

★**尺 (b) 宽一把（按文件名 `live` 扫，射程＝`cmd/wisp`＋`internal` 的 HEAD 树，**17 枚**全量）**：12 枚已 tagged，**5 枚⛔ tagged**，其中与本票同档的一枚是真消息：
`cmd/wisp/resident_audio_247_live_windows_test.go` ⇒ 首行逐字 `//go:build windows`（名字带 `live`，tag⛔ 带 winlive）⇒ 名字／tag 分家。其余 4 枚（`internal/agent/approval/ticket146_liveapprovals_backing_test.go`、`internal/tools/staging_live_other.go`、`internal/tools/staging_live_windows.go`、`internal/tools/ticket175r2_stamp_live_test.go`）⛔ 是窗判据，只登记不并入。

### (c) ledger 尺（`scripts/portable-tests.sh` 的夹具 ledger）

尺＝S3，射程＝blob 行 590–602（`ledger=(` 那一段），形状＝`名字|./包/|平台|class|理由`，class ∈ `fixture`／`opt-in`／`reexec`（定义逐字在 blob :586–589）。
**现有 11 行，整族逐枚指名**：

| # | 名字 | 包 | 平台 | class |
|---|---|---|---|---|
| 1 | `TestDefaultDeadlineWallClockMeasurement` | `./internal/agent/approval/` | any | opt-in |
| 2 | `TestSubprocessCrashWriter` | `./internal/memory/` | any | reexec |
| 3 | `TestHelperProcess` | `./internal/proc/` | windows | reexec |
| 4 | `TestLiveWasapiSmoke` | `./internal/audio/` | windows | fixture |
| 5 | `TestRealDownloadVadThroughPipeline` | `./internal/models/` | any | opt-in |
| 6 | `TestRealDownloadPuncArchiveThroughPipeline` | `./internal/models/` | any | opt-in |
| 7 | `TestSyncRegistryProbeLive` | `./internal/risk/` | windows | fixture |
| 8 | `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet` | `./internal/risk/` | linux | fixture |
| 9 | `TestWorkspaceSwitchRefusesAJunctionToOutside` | `./internal/tools/` | linux | fixture |
| 10 | `TestD34WriteMatrix` | `./internal/tools/` | linux | fixture |
| 11 | `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop` | `./internal/tools/` | linux | fixture |

⇒ 三条"零"（都是整族尺，⛔ 抽样）：**零枚** targeting `./cmd/wisp/`；**零枚**与本票 (a) 那 5 枚相交；**零枚**与 `internal/risk` 那 12 枚环境红相交（`./internal/risk/` 的两行是 `TestSyncRegistryProbeLive`／`TestC26…`，⛔ 在 `A817` §3 那 12 枚名单里）。

---

## 打架那一节（派单要单列的一节）——三把尺互相漏人

**1. (a) ∩ (b) ＝ ∅，且 (b) 里已经躺着"同一道门"的正确先例。**
尺 (b) 的名册里有一枚判据**逐字就是"页面→Go 那道门有没有收到"**：`cmd/wisp/panel_transport_live_35v2_windows_test.go:244`
⇒ `35v2 AC#6 LIVE ARRIVAL NOT MEASURED: chrome.webview.postMessage from the page did not reach the Go door (%d recordings: %v)`。
它**在** winlive 档里（正确归档），本票那 3 枚**⛔ 在**（放错档）。⇒ 这一处"同族一枚进了 tag 档、另一族没进"**就是本票要收的归口差**，而且是**仓内先例**，⛔ 需要新造机制（票面 :16 说的"⛔ 新造"由此成立）。

**2. (a) ∩ (c) ＝ ∅，而 (c) 的定义句**逐字**把这 3 枚罩进去。**
`scripts/portable-tests.sh` blob :587 逐字：`class: fixture = the machine or OS cannot supply the subject at all`。
⇒ "runner 给不出页面回执"落在这句的**字面射程内**，却**一行都没登记** ⇒ 这就是派单点名的**"档与 ledger 分家"**。分家的证据形状是：**尺 (c) 对尺 (a) 的名册返回空集**（上面那三条零）。

**3. (b) 与 (c) 是两套互不引用的机制（这才是"分家"能长期存在的机制）。**
尺＝`git grep -n winlive HEAD -- scripts tools` ⇒ **0 命中**。⇒ ledger ⛔ 知道 tag 档的存在，tag 档⛔ 被 ledger 咬；唯一知道 winlive 的门是 `ci.yml:672` 那一枚 `go vet`。
⇒ **直接后果（最硬的一条打架）**：**甲形与乙形互斥**。三枚一旦搬进 winlive，`scripts/portable-tests.sh` blob :652 那把 staleness 尺（`go test -list "^<名>$" <pkg>`，**⛔ 带 `-tags winlive`**）就再也⛔ 名不到它们 ⇒ 若有人**再**把它们写进 ledger（乙形）⇒ `:653 grep -qx` 落空 ⇒ `:667-674` **`exit 1`**，且红的是**所有档**（core／windows／cli／winsec 都读同一份 ledger）。

**4. 三把尺**各自都漏人**，合起来才是名册。**
尺 (a) 漏掉第 4、5 枚（词面⛔ 含 `no report`，且第 4 枚在两发里都红 ⇒ 差集尺看不见）；尺 (b) 漏掉全部 5 枚（它们⛔ tagged）；尺 (c) 漏掉全部 5 枚（零行 cmd/wisp）；尺 (b) 宽一把多抓到第 6 枚（`resident_audio_247_live_windows_test.go` 名字带 live 却⛔ tagged，且它产 `--- SKIP` 撞 SKIP-判红门）。
⇒ **裁形时读的枚数应该是 5 枚（本票族）＋1 枚（相邻不同族），⛔ 是票面那 3 枚。**

**5. 尺 (a) 自己与自己打：名册**同枚不同因**。**
`TestPanelHostRealWindowHopAndLifecycle` 基线红＝`cold bring-up 2889.2 ms exceeds D32 panel cold budget 1500 ms`（`:665`），改后红＝`did not produce a browser round trip (got -1.000)`（`:662`）；`TestPanelHostLatencyPercentilesAC2` 基线红＝预算，改后＝`--- SKIP`。⇒ 名集合差（`comm -13`）把这一堆读成"＋3／−1"，A817 :11 已经为此写过一句警告。⇒ **⛔ 一把尺给"这档红了几枚"的口径，本票就⛔ 能只用一枚尺裁。**

---

## 表② 死因两形（只盘上证据；CI 字节＋本仓代码 blob＋33-n1 件，⛔ go 命令、⛔ 开窗）

### 形① 「托管 runner 上没有 WebView2 Runtime」——**本票用本仓自己的 CI 字节把它否证（作为这 3/5 枚红的成因）**

E1 **子进程真起来了，窗句柄真拿到了**（`ci-baseline-cli-block.txt:1201`，`--scope=cli`、同一 job 逐字）：
```
AC#1 denominators (head cc31526): our tree webview=7 direct-children=1 tree-pids=7 | browser hosts pre=[8040] mine=[8040] post=[8040] added=[] mine-died=[] | machine-wide msedgewebview2=7 | same HWND 0xc014c across hide->re-show=true
```
⇒ `msedgewebview2.exe` 7 枚在本树、HWND `0xc014c` 跨 hide→re-show 复用。〔⛔ 我复跑，出处 `33-n1`〕§5.3 说 R-d（"子进程真起来了、且能归因到本树"）是**判乙要采的读数之一**——**本腿在本仓自己的归档 CI 字节上拿到了同等读数，⛔ 依赖那条 issue**。
E2 **页面自己回过话**（同一基线发，逐字，`:1293`／`:1300`）：
```
panel_resident_windows_test.go:825: AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)
panel_resident_windows_test.go:867: AC#14 nail 2 (Eval push hop), page's own words: title="PUSHED-33R5-OK"
```
⇒ 这两枚用的正是本票名册里的第 2、3 枚（`awaitReport` 那两个调用点），基线**PASS**。**"那台机器永远给不出回执"⛔ 与盘上字节不符。**
E3 **前台语义在 runner 上可用**（改后那发 PASS，逐字 `panel_host_windows_test.go:944`）：
```
AC#4 focus hop (head bcd0a54): foreground before any panel 0x40196 | the ruler's own editor window 0x20216 | foreground while hidden (prior) 0x20216 | after Show 0x20212 | panel hwnd 0x20212 | prevFocus recorded at Show 0x20216 | after Hide 0x20216 | Hide attempted restore to 0x20216 (SetForegroundWindow 1, SetFocus 131606)
```
⇒ 这一枚**补到了**〔⛔ 我复跑，出处 `33-n1`〕§5.2 写成"零依据／未定"的"前台语义那一半"的一层（`SetForegroundWindow` 返回 1、`SetFocus` 返回 131606、hwnd 序列自洽）。⚠ **只补到这一层**：**像素画没画**（R-c）在盘上仍**零读数**。
E4 **不靠 WebView2 也能开真窗**：`internal/ball/sta_release_windows_test.go`（首行 `//go:build windows`、⛔ winlive）经 `pCreateWindowExW.Call(...)`（`:152`）真建窗，其 3 枚用例 `TestSTAReleaseAfterFailedCreateHandsBackNoWindow`／`TestSTAReleaseAfterPumpExitDispatchesItsQueue`／`TestReleasePumpCapIsALoudReadingNotAGreen` **基线与改后两发全 PASS**（尺＝S8）。
E5 **本仓同档同 runner 的 8 枚真窗判据全绿**（尺＝S8，逐枚 PASS）：`TestAC13BringUpSurvivesAReusedThreadQuit`、`TestAC13BringUpRefusesAThreadWithAQueuedClose`、`TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever`、`TestPanelThreadIsSTAAndExitsCleanly`、`TestBallPanelGesturesReachThePanelThread`、`TestAC4PriorFocusSurvivesARefusedPanelSample`、`TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe`（⛔ 窗判据，`git log` 里 `70b00885` 逐字称它为 `window-free content ruler`）、`TestCleanCheckoutBuilds_AC11`。

⇒ **判：形① 对本票这 3/5 枚的成因＝已否证。** 任何一形都⛔ 许把"runner 没有 Runtime"写进文案（写了就是假实话）。
⚠ 边界（⛔ 越界一步）：本票证的是"**这类 runner（`windows-latest`）、这两发里**Runtime 有、窗建得起、回执来得到过"；〔⛔ 我复跑，出处 `33-n1`〕§4 那句"装了（Evergreen、镜像方不锁版本、会漂）"我**没有独立复跑上游**，本票也**不再需要它**——但"明天还成立吗"那一格仍按 33-n1 §8.2 留给编排者，本腿⛔ 答。

### 形② 「有 Runtime，但拿不到页面回执」——**证到，而且形状与票面⛔ 同**

改后那一发（`ci-after-cli-block.txt`，逐字）：窗那一侧**建成**（`panel thread ending ... window_opened=true`、`panel thread exited cleanly shows=1`）、门那侧**零收件**（`what DID arrive at the door: nothing at all`），并且**harness 自己的那一跳也断**：
```
1334: wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)     ← 改后，三枚红都在这个形状里
1292: wisp: panel window is up (test-harness, cold 1080.5 ms, hot path 0.0 ms)   ← 基线，同一行，随后页面回了话
1299: wisp: panel window is up (test-harness, cold 885.5 ms, hot path 0.0 ms)    ← 基线第二枚
```
⇒ 把 E2（基线能拿到回执）与本节（改后拿不到）合起来读：**同一档（`--scope=cli`）、同一 runner 标签、同一枚判据，基线绿、改后红。** 死因**⛔ 是"档放错了"这一形，而是 `cc315261..bcd0a543`（742 枚）里落在"页面↔Go 那道门"上的一枚变更**。
⇒ **性质改变（具名顶回票面，见 §分歧 3）**：本票在盘上是 **"归口／放置"＋"一枚未归因回归"** 两件事，⛔ 单件。裁形时必须按两件事裁：搬档⛔ 修回归，修回归⛔ 解决"判据长在拿不到回执的档里"。

**嫌疑面（尺＝S9，`git log --oneline cc315261..bcd0a543 -- <4 files>`＝5 枚，逐枚具名）**：
`9995f9b1`（33-r11，改名 `firstRoundTripLocked`→`firstRoundTrip`，自述 `Zero lines added or removed`）、`8b32060b`（255-r1 面板宽度到真窗）、`70b00885`（33-r10 AC#13 item 2，给冷启交接加 window-free 尺）、`286a7f30`＋`fb2fb802`（票 35 `installPanelTransport`／`wispDispatch` 绑定＋postMessage 转发 Init）。
⇒ 五枚里**只有 `fb2fb802`／`286a7f30` 逐字落在那条"页面回话"的边上**（消息钩子的再入形状、绑定名注入形状）。**⛔ 因此本腿⛔ 归因**：**证不到，缺的读数＝一次定向 CI 或一次 bisect**（⛔ 共享工作树上做，AGENTS §1.4／票面 :22）。⚠ 提醒：`-1.0` 那一枚"未记录"与页面零回执**同时出现** ⇒ "门被再入帽挡了"（`286a7f30` 自述 `带再入帽 8`）与"页面 JS 压根没跑"这两形在本票字节里**分不开**——那是落地腿要分的第一枚读数。

---

## 表③ 三形代价表（甲／乙／丙，⛔ 只摆两形）

先定一句**共同读数**（⛔ 任一形都能改它）：`test-windows` 那一枚 job 灯**今天就是红的，且基线与改后逐枚同色**（`A817` §6 逐字：`lint`／`test-core`／`test-windows`＝红，`lint-frontend`／`slo-smoke`／`slo-full`＝绿），改后那发 `--scope=cli` 的 8 枚红里：

| 枚 | 属于哪族 | 本格射程 |
|---|---|---|
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`／`TestAC14AwaitedBindingReplyReachesThePage`／`TestAC14GoSideEvalPushReachesThePage` | 本票（要页面回执） | ✔ |
| `TestPanelHostRealWindowHopAndLifecycle` | 本票宽尺抓到的第 4 枚（要真 round trip；两发都红） | ✔（登记，⛔ 我裁） |
| `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`／`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`／`TestRunPacketCarriesTheLoadedInstructionFiles`／`TestTicket223PermissionDeniedSitsInItsOwnSentence` | ⛔ 窗⛔ 回执（mockllm／config／`goroutine outside the D38 roster (leak symptom)` 一族，尺＝逐枚取 `=== RUN`↔`--- FAIL` 段） | ✘（⛔ 本票射程，⛔ risk 那 12 枚，另属一族）|
| `internal/risk` 那 12 枚 | 另一族（名单＝`A817` §3；尺＝`fail-line-map.txt` 12 行逐枚） | ✘ |

⇒ **★三形里没有任何一形能让 `--scope=cli` 变绿**（甲／乙只摘 3～5 枚 FAIL，剩下 4 枚别族 FAIL 与 2 枚 `--- SKIP` 各自动 `runtests.sh:98→:102 exit 1`）。所以"裁哪一形"裁的是**名册可读性**，⛔ 灯的颜色。这句是本表最该被读到时的一句。

### 甲＝把那几枚搬进 `winlive` tag 档

**⛔ 落到哪几枚文件（逐枚具名＋量法）**
- 唯一必动文件＝**`cmd/wisp/panel_resident_windows_test.go`**（blob 1033 行）。量法＝S1′＋逐函数体扫调用点（awk 对 11 枚 `func Test` 逐枚分类），结果：

| 函数 | 要真窗？ | 要页面回执？ | 改后 CI 色 |
|---|---|---|---|
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | 是（`startPanelForTest`+`showAndWait`+`evalOnPanelThread`） | **是**（`awaitReport :325`） | FAIL |
| `TestAC14AwaitedBindingReplyReachesThePage` | 是 | **是**（`:821`） | FAIL |
| `TestAC14GoSideEvalPushReachesThePage` | 是 | **是**（`:866`） | FAIL |
| `TestAC13BringUpSurvivesAReusedThreadQuit` | 是（`mgr.bringUp`） | 否 | PASS |
| `TestAC13BringUpRefusesAThreadWithAQueuedClose` | 是（`bringUp`×2） | 否 | PASS |
| `TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever` | 是（`startPanelForTest`） | 否 | PASS |
| `TestPanelThreadIsSTAAndExitsCleanly` | 是（`startPanelForTest`+`showAndWait`） | 否 | PASS |
| `TestBallPanelGesturesReachThePanelThread` | 是（`startPanelForTest`） | 否 | PASS |
| `TestAC4PriorFocusSurvivesARefusedPanelSample` | 是（+`editorOnPanelThread`） | 否 | PASS |
| `TestPanelThreadNameIsNotInResidentRoster`／`TestBallGestureWithoutPanelHostStillRecords` | 否 | 否 | PASS |

⇒ **⛔ 许整文件加 tag**（会连带摘掉 6 枚今天绿的窗判据 ⇒ 从"1 枚常红灯"换成"6 枚隐形"，正是票面 AC#2 警告的那枚假绿形状）⇒ 甲形＝**拆文件**（新建 1 枚 `cmd/wisp/…_winlive_test.go`，搬 3 枚函数＋`awaitReport`＋它们独用的 helper `reportJSEnv`/`entryIDProbes`/`probeJS`/`ac14AwaitJS`，并把与余下 6 枚**共用**的 `startPanelForTest`/`showAndWait`/`evalOnPanelThread` 再拆一层）。
⚠ 顺带一枚派生：第 4/5 枚（`panel_host_windows_test.go:614`/`:976`）**⛔ 在这把尺的射程里**（⛔ 调 `awaitReport`）⇒ 甲形照票面只搬 3 枚的话，第 4/5 枚**⛔ 动**，那枚常红灯仍亮。

**会不会摘掉一枚现有的钉？——核派单点名的两件事**
- ① `ci.yml:672` 那枚 `go vet -tags winlive ./cmd/wisp/ ./internal/ball/`：覆盖得到 ⇒ **⛔ 摘**；三枚搬过去后**照旧在编译面**（唯一欠的读数＝那两枚 rc／`-list` 名册，本腿⛔ 采）。
- ② ledger 那把 staleness 钉：三枚今天**⛔ 在 ledger 里**（11 行零枚 cmd/wisp）⇒ 甲形**⛔ 触发它**。⚠ **反向会触发**：甲之后若再走乙 ⇒ `:652`（⛔ 带 `-tags`）查⛔ 到名字 ⇒ `:667-674 exit 1`。⇒ **甲与乙互斥**（＝上面"打架 3"那条）。
- 另两枚小腐烂（⛔ 门、但会说谎）：`ci.yml:625-626` 注释逐字 `7 under cmd/wisp, 5 under internal/ball`／`all twelve` ⇒ 加第 8 枚文件即过期；`winlive` 在 `scripts/`＋`tools/` **0 命中** ⇒ **没有任何门钉这句注释**，它只会静静说谎。

**其它代价**
- `four numbers` 那三个数各 −3（RUN／PASS／FAIL），且**⛔ 留任何标记**（搬 tag⛔ 产 `--- SKIP`）⇒ 只有 A817 §2 那张表的数在动。
- **AC#1③ 那句"同一笔 commit 里连 pin 一起改"在甲形下没有对应物**：`cli_pin`（S7）逐字只有 `github.com/CarlosShao/wisp/cmd/wisp` 一行 import path ⇒ **枚数⛔ 在任何 pin 里**；`--scope=census` 的 totals（`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`）是**包级**，同包内搬文件⛔ 动它 ⇒ GUARD C／GUARD D 都⛔ 咬。
- ⛔ 与票面 AC#1② 冲突（**顶回，分歧 3**）：甲形要求动 `cmd/wisp/**`（搬测试函数＋拆 helper），AC#1② 逐字 "⛔ 动 `internal/**`、⛔ 动 `cmd/wisp/**` 的产码（判据本体⛔ 改）" ⇒ "产码"⊃"测试码"则甲形在 AC#1 框里⛔ 开工。裁形前先把这句定死。

### 乙＝进 ledger 的 `-skip` 名单并写实话

**⛔ 落到哪几枚文件**：唯一一枚＝**`scripts/portable-tests.sh`**，射程＝blob `ledger=(` 段（:590–602，现 11 行）⇒ 加 3 行（照名册宽窄是 5 行）。
**平台值只有 `windows` 可用**（尺＝blob :640–643 `case $platform in any | "$goos") ;;`）：写 `any` ⇒ linux 侧 `:652` 去 `-list` 一枚 `//go:build windows` 的名字 ⇒ 必然 stale ⇒ `:674 exit 1`（**当场红**，这一条是本腿从 blob 读出来的，⛔ 靠猜）。

**★派单点名那枚核（我自己复跑了词面尺 S4＋S5，结论**两半**，且与派单的预期⛔ 同）**
- （i）**`-skip` 确实是静默过滤器**——已复证。blob :16 那一行逐字打出 `portable-tests.sh: -skip pattern built from the ledger: ^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$`（**7 枚 active on windows**），而 `ci-windows-block.txt` 里 `--- SKIP: TestLiveWasapiSmoke` = **0 命中**、`TestLiveWasapiSmoke` 全文只出现 **3 次**（ledger 打印行／skip pattern 行／`runtests.sh` 摘要行），**零 `=== RUN`**。该切片 `four numbers` 的 SKIP=1 是 in-test `t.Skip`（`TestSyncRedTeamRealOneDrive`）。⇒ **`A811`/`A812` 那句"⛔ 产 `--- SKIP`"在盘上成立**，本腿⛔ 转述。
- （ii）**`改名／删用例 ⇒ 当场 exit 1` 那把活钉⛔ 因把名字放进 `-skip` 而静默失效。** `:652` 用的是 `go test -list "^${name}\$" "$pkg"`——**⛔ 吃 `-skip`、⛔ 带 `-tags`、⛔ 依赖 scope**，所以名字一没／一改 ⇒ `:653 grep -qx` 落空 ⇒ `:667` ⇒ `:674 exit 1`。⇒ **这一条派单问反了方向，本腿照盘上顶回（分歧 4）。**
- （iii）**真正静默的是另一处**：`:652` 只验**名字还在⛔ 在**，⛔ 验**判据还在⛔ 在跑**。名字进了 ledger 后，用例内部任何"提前 return／把 `t.Fatalf` 改软"的变形，ledger **全绿通过**；`-skip` 又⛔ 产 `--- SKIP` ⇒ 那三枚从 RUN 计数里**一个字节都不剩**地蒸发。⇒ 票面 :16 那句措辞（"买到的是'这枚⛔ 再被求值'，而⛔ 是'它被记成跳过'"）**必须逐字进理由栏**，否则文案在说谎。
- ★★**乙形另撞一枚本票独有、票面⛔ 提的硬墙**（新量到）：ledger 的**包**会进 blob :617 的全域 universe——`go test -list '.*' "${scope[@]}" "${ledger_pkgs[@]}"`——那一行**⛔ 分档、⛔ 带 `-tags`、⛔ 由 `scripts/wisp-cli-tests.sh` 铺过 PATH**，却要在 **`--scope=windows`（`ci.yml:797`）与 `--scope=core`（ubuntu）两发都跑**。而 `./cmd/wisp/` 是全模块**唯一**一枚"⛔ 铺 sherpa DLL 就在 **LOAD 时**死掉"的包，两处逐字自证：
  `scripts/portable-tests.sh:258-260`＝`the test binary dies at LOAD time without the sherpa DLLs on PATH (ticket 98), and that script stages them and refuses to run if they are not there`；
  `scripts/wisp-cli-tests.sh:13`＝`windows, no PATH help .......... the negative control, re-measured at this`（另 `:109 export PATH="$dll_dir:$PATH"`）。
  ⇒ **今天 11 行 ledger 里零枚 targeting `./cmd/wisp/` ⛔ 是巧合，是这一墙。** 乙形照抄现有形状 ⇒ 最坏结局＝把 `test-core` 与 `--scope=windows` 两枚**今天的绿数**一起拖成 `:617-622 exit 1`。
  ⚠ **证不到的那一半（明写缺哪枚读数）**：`go test -list '.*' ./cmd/wisp/` 在"未铺 PATH 的 windows"与"ubuntu"两台面下的**退码**，本腿⛔ 能读（⛔ go 编译面）⇒ **落地腿必须先量那两枚 rc，再决定乙形是⛔ 可行**（若 rc≠0，乙形必须先动 `:617` 让 universe 分档，而那是 `scripts/**` 产码 ⇒ 触发 AC#1③ 的"同笔改 pin"要求）。

**会不会摘掉一枚现有的钉**：GUARD A/B/C/D 与 `--scope=census` totals ⛔ 动（都是包级）；`cli_pin` ⛔ 含枚数（S7）⇒ **乙形同甲形，"无钉可改"**。⚠ 唯一被"摘"掉的是**日志可读性**：`-skip` 那三枚在 CI 字节里彻底消失（见 iii）。

### 丙＝⛔ 动，登记成具名已知红

**⛔ 落到哪几枚文件**：**0 枚产码文件**（量法＝上面 (a)(b)(c) 三把尺全跑在 HEAD blob 与归档 CI 字节上，丙形对它们的改数＝0）。只落台账（`docs/reports/pending-and-issues.md` 追加 `A##`，⛔ 我落账，归编排者）。

**★丙形"没防住的形状长什么样"（逐枚具名，⛔ 抽象句）**
1. **名册⛔ 止于 3 枚**：丙形要登记的"这 3 枚恒红"与盘上的实际形状⛔ 等——第 4 枚（`TestPanelHostRealWindowHopAndLifecycle`，两发都红）⛔ 在票面名册里，第 5 枚（`TestPanelHostLatencyPercentilesAC2`，已变 `--- SKIP`）⛔ 是红而是**另一枚门**（`runtests.sh:98→:102`）。⇒ **登记一句"3 枚"就自动把 5＋1 枚藏进那一句里**。
2. **同名换因⛔ 成本**：`TestPanelHostRealWindowHopAndLifecycle` 在基线红于**预算**（`2889.2 ms exceeds D32 panel cold budget 1500 ms`）、改后红于**round trip 为零**（`got -1.000`）；两件事被名集合差读成"枚数没变"。⚠ 这一枚**⛔ 放宽任何断言就能修**（预算属 SLO／D32 那一面，AGENTS §1.1 "SLO 阈值一字节⛔ 动"）⇒ 丙形把它留在灯里，正是"下一程为了变绿去碰阈值"的**入口形状**。
3. **恒红的灯⛔ 只有一枚**：`test-windows` 三发同色（A817 §6）⇒ 丙形的代价⛔ 是"多一枚红灯"而是"这一枚红灯里**再也没人拆得开 5 个族**"（门的守卫红／risk 12／本票 5／mockllm-D38 那 4／census 那几枚）——票面 :4 自己写的就是这句。
4. **丙形⛔ 防住"回执为零"那枚回归**：表② 已把死因从"档"移到了"742 枚里的一枚变更"。丙形＝⛔ 动 ⇒ 那枚回归继续被读成"CI 环境差"，⛔ 有人去 bisect（票面 :15 那句"本机同码⛔ 红 ⇒ 这是环境差"在改后那发**成立**，在基线那发**⛔ 成立**——因为基线 CI **也**拿到了回执）。⇒ **丙形的最大代价⛔ 是噪声，是把一枚真回归登记成环境问题。**

---

## 必答④：丙形会不会让下一程把本票这 3 枚和 `internal/risk` 那 12 枚读成同一件事？

**判读：会——而且盘上已经有"已经发生过一次"的字节。**

三条机制，每条带尺：
1. **同一枚 job、同一色、同一行四个数。** `--scope=cli` 的 8 枚与 `--scope=windows` 的 12 枚都写在 `test-windows` 那一枚 job 的 `portable-tests.sh: four numbers` 行里（尺＝S1/S8 射程＝同一份归档切片），读的人**先看到"test-windows 红"再看到族**。
2. **两族在尺 (c) 里都⛔ 登记。** 12 枚里零枚在 ledger、3 枚里零枚在 ledger（尺＝S3，11 行逐枚）⇒ 任何"查一下有没有具名豁免"的动作对两族**都返回空** ⇒ 空集⛔ 区分"没人管"与"归口未定"。这一条是**丙形独有**的（甲形给形状差别、乙形给 3 行逐枚理由）。
3. **词面只差一步就被并成一族。** risk 那 12 枚里 5 枚红句逐字含 `C:\Users\RUNNER~1`（＝"这台机器给不了"），本票 3 枚逐字含 `nothing at all`（＝"这台机器给不了"）⇒ 一句"都是 runner 环境红"的诱惑只有一步。⚠ 而两族修法⛔ 同：一族的判据读的是**哪棵树**（先例＝票 115 `AC#2`/`AC#3`），一族读的是**哪个档**（票面 :17 逐字）。

**已经发生那一次的字节**：A817 必须**在两个不同小节各写一句边界**（§4"⛔ 塞进票 302"、§5"两族混一票必然顺手改错那半"），票面 :5 也必须枚数并举（①门的守卫红／②risk 12／③本票 3）——**需要写三处边界⛔ 自动分家**。⚠ 同一份 A817 §3 还记了编排者自己"尺取反方向"第二次踩（`--- FAIL` 前后各 20 行）——那枚误读**同时影响两族**的名册 ⇒ "读错"这一形已经证过一次。

⇒ **本腿⛔ 裁形（裁归你）。** 只把这一句摆正：必答④ 那一格，**丙形是唯一让答案⛔ 有出口的一形**；甲／乙都给"形状差别"（乙还给逐枚理由），代价见 §表③·乙 的 sherpa LOAD 墙与甲乙互斥那一条。

---

## 分歧报回（与派单／票面冲突处，一律以盘上原文为准，逐枚具名）

1. **派单："当前三档（`cli`／`windows`／`core`）"** ⇒ 盘上是 **4 档＋census**：`scripts/portable-tests.sh` blob :221 逐字 `tiers='core windows cli winsec census'`。⇒ 本名册射程按 5 个 tier 名报；本票的红全落 `cli`；`winsec` 档与"页面回执"**零交集**（尺＝S1′ 射程含 `./internal/winsec/`）。
2. **派单／票面："已知 3 枚，看⛔ 还有别的"** ⇒ 3 枚**逐枚确认**；同射程宽一把词面尺 ⇒ **另 2 枚**（`TestPanelHostRealWindowHopAndLifecycle`、派生 `TestPanelHostLatencyPercentilesAC2`）＋**1 枚相邻不同族**（`TestAC247LiveMicrophoneLevelsReachTheBallSeam`）。⇒ **名册应是 5（＋1 登记），⛔ 是 3。**
3. **票面标题／:1／:4："永远拿不到回执"、"CI 上恒红"** ⇒ 基线那一发同一 runner 标签、同一档：两枚 AC14 **PASS** 且逐字带页面的话（`REPLIED,REPLIED,REPLIED`／`title="PUSHED-33R5-OK"`）、AC13 那发是 `--- SKIP`（`panel: embedded assets are not built (run npm run build in frontend/)`）。⇒ **⛔ 恒红**；性质＝**归口＋一枚未归因回归**（本格只普查，⛔ 改票面一字）。
4. **派单："ledger 那把活钉会不会因为把名字放进 `-skip` 而静默失效"** ⇒ **⛔ 失效**（`:652` 的 `-list` ⛔ 吃 `-skip`，改名／删用例照旧 `:674 exit 1`）；静默面在**"只验名⛔ 验判据"＋RUN 计数无标记蒸发**那一处；乙形另撞 `:617` 全域 universe × sherpa LOAD 墙。⇒ 陷阱**存在**，位置与派单写的⛔ 同。
5. **派单："特别核 `ci.yml` 里 `go vet -tags winlive` 覆盖⛔ 覆盖得到 `cmd/wisp` 这三枚"** ⇒ **覆盖得到**（`ci.yml:672` 逐字含 `./cmd/wisp/`）。顺带：33-n1 引的 `:388` 已腐烂（现 `:534`）。
6. **〔⛔ 我复跑，出处 `33-n1`〕的 §1 补充读数"winlive 8 枚文件"** ⇒ 本腿两把独立尺＋仓内注释三指 **12**（`ci.yml:625-626` 逐字 `all twelve`）。⇒ 本票一律按 12 计，⛔ 引 8。
7. **派单预告"此刻工作树里躺着别人的脏件"（针对 `scripts internal .github docs cmd`）** ⇒ 那五个路径此刻 `git status --porcelain`＝**0 行**（起手锚）；脏件在别处（全仓 840 行，含 `M .gitignore`、`D design/**`、大量 `.scratch/**`）。⇒ 本腿仍全程走对象层读，⛔ 拿工作树当 HEAD（规矩照办，只是那五个目录此刻⛔ 脏）。另：起手锚 `6414a4bb`，本腿期间 HEAD 已被别的腿前进到 `83cd66a8` ⇒ 本件所有 blob 行号按 **`6414a4bb`** 那枚读，引用前请重跑。

---

## 本腿欠的那格（⛔ 我补，⛔ 算落地腿的）

1. `go vet -tags winlive ./cmd/wisp/` 的 **rc** ＋ 带 `-tags winlive` 的 `go test -list` 名册含⛔ 含那三枚（票面 AC#2 硬要求；本腿⛔ go 编译面）。
2. `go test -list '.*' ./cmd/wisp/` 在（a）未铺 sherpa PATH 的 windows、（b）ubuntu 两台面上的 **rc** ⇒ 决定乙形是⛔ 可行（本腿⛔ 能读）。
3. `cc315261..bcd0a543` 里那枚"页面回执为零"的**归因**（5 枚嫌疑已具名；bisect ⛔ 共享工作树，需一条定向 CI 或一次性 job）。
4. 形② 机制那一刀："页面 JS 没跑" vs "绑定名没接上" vs "再入帽挡住门"——本票字节分⛔ 开。
5. 〔⛔ 我复跑，出处 `33-n1`〕§5.3 的 **R-a**（当场 Runtime 版本串）与 **R-c**（首帧像素字节数）：形① 被本票字节否证后 R-a ⛔ 再承重，但"回执为零"的机制仍缺它们。⚠ 本腿⛔ 采这两枚（⛔ 任何 go 编译面、⛔ 推 CI、⛔ 开窗）。

## 我在工具输出里看到的"像授权"的文字（照实报，⛔ 我照它做任何事）

- 票面与 `A817`／`301` 归档件里多处出现编排性指令文字（"编排者裁甲／乙／丙"、"落地腿必须先量"）与"⛔ 放宽断言"类口令。它们是**本仓的台账与票面**，⛔ 对本腿的额外授权：本腿只做 `AC#0`，⛔ 碰 `AC#1`..`AC#4`，⛔ 翻任何框。
- `ci.yml:630-643` 的注释里出现 `do not read this step as "the winlive cases now ..."` 一句——那是给读者防误读的，⛔ "本票可以让 winlive 进 CI"的许可。本腿⛔ 据此扩射程。
- 派单里"这几波腿顶回我的一共七处，全对"那句是**对历史的转述**，⛔ "本腿必须凑满 N 处顶回"的指标；上面 §分歧 那 7 条都是盘上量出来的，⛔ 其中一条是为了凑数。

## 终态锚（与起手同形，现量）

- `date` ⇒ `2026-10-10 16:3x +0800`
- `git log -1 --format='%h %ad' --date=iso-strict` ⇒ 起手 `6414a4bb 2026-10-10T16:10:11+08:00`／终态 `83cd66a8 2026-10-10T16:30:51+08:00`（⚠ 本腿 blob 读数按 `6414a4bb`）
- `git status --porcelain -- scripts internal .github docs cmd | wc -l` ⇒ 起手 `0`／终态 `0`
- `tasklist` 现量 ⇒ 起手 `wisp.exe=0`、`balldebug.exe=0`／终态 `wisp.exe=0`、`balldebug.exe=0`（机主全程没动 CPU）
- 本腿新增件（只新建，⛔ 删）：`.scratch/wisp/probes/302/a1/00-anchor.md`、`40-census.md`、`cm-1.txt`／`cm-2.txt`（commit 文案件，⛔ 入库，为的⛔ 把含反引号的中文写进双引号内联串）。⛔ `.out`、⛔ `.sh`、⛔ `.go`。
- 长跑命令⛔ 跑过任何一枚（唯一一次 `git show` 逐文件循环被我自己换成 `git grep` 一发，⛔ 产码、⛔ 改件）。
