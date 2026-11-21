# 票 450 — 本机整包跑到 **约 400 秒**就在真 WebView2 窗那一区**整进程死亡**，两发两次都不是断言红

**状态**：待派（能否复现、是不是我自己的驻留球造成的，都要读数）

## 为什么值得单独一张票
本编队每个写腿都要跑一次 `go test ./cmd/wisp` 整包来证明"没把别人弄红"。这条路**今天在本机走不完** —— 它在 `resident_*`／`ticket*`／`voice_*` 那一区**还没开始跑**之前就整进程死了。后果不是"少几枚绿"，是**那一片测试在本机等于没人跑过**，而大家会拿"两发都 FAIL 0"当覆盖证明。这正是"看起来绿、实际没跑"的形状。

## 现量（2026-11-14 22:xx，票 228 切片 1 的实现腿，两发）
| 发 | 读到 | 死在哪 | 用时 |
|---|---|---|---|
| 1 | `PASS 339  FAIL 0` | `Test145ArgsJSON…` 半途 | 416 s |
| 2 | `PASS 334  FAIL 0` | `TestAC4PriorFocusSurvivesARefusedPanelSample` 半途 | 403 s |

两发**都不是断言失败**，是**进程消失**（无 `--- FAIL` 行）。`228-r3b` 记过同一族。

## 关键背景（可能是真因，⭐别先入为主）
两次死亡发生时，**编排者自己的驻留球在跑**（PID 49348，2026-11-14 16:23 起）。已知的**本机专属红**里就有 `TestAC246DevLegIgnoresTheTestTaskInjection`（"a dev Wisp is already running in this session"）——那说明**驻留球会与整包里的用例争用同一份会话/窗口资源**。⇒ 第一嫌疑是"球在跑"，⛔ 一上来就判"WebView2 不稳"。

## AC

- [x] **AC#1 先把"死"变成可判事实**：复现一次，抓**退出码／最后一行输出／Windows 事件日志里的 application error 条目（进程名＋模块）**。⛔ 用"卡住了"描述。
- [x] **AC#2 二分：球在跑 vs 球不在跑**：同一命令各跑 1 发 —— (a) 球在跑；(b) 球停（**由机主停，⛔ 腿自己 kill**；若机主不在，用"从未起过球的干净用户会话"替代并写明）。给两发的 `RUN/PASS/FAIL/SKIP` 与是否死亡。
- [x] **AC#3 二分：WebView2 窗那一区 vs 其余**：把死亡点附近那几枚用 `-test.run` 单独跑（含 `Test145ArgsJSON*`、`TestAC4PriorFocus*`），看在**没有整包负载**时是否也死。
- [x] **AC#4 判归属并落成规矩**：若 (a) 死 (b) 不死 ⇒ 真因＝驻留球争用，**把"跑整包前必须停球"写进本仓的测试前置说明**（`scripts/wisp-cli-staged-tests.sh` 的注释或 `docs/evidence/` 一份具名说明，⛔ 改 `docs/specs/**`）。若不是 ⇒ 交回可判事实与最小复现，**本票不许自己猜着改产码**。
- [x] **AC#5 门禁**：本票若只改脚本注释/证据件，跑 `d22scan` 与 path-length 即够；若动产码，按整包 + 本机两枚已知红具名口径执行。

## 边界
⛔ `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 转移表；⛔ 放宽或跳过任何既有断言（若真需要跳过，必须**具名＋给数字**并写进证据）；⛔ 自己 kill 机主的球；⛔ `frontend/**`／`design/**`；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。
写面：`scripts/wisp-cli-staged-tests.sh`（**仅注释**）、`docs/evidence/s1/450-*`、`.scratch/wisp/probes/450/**`、本票 Progress log。

## 撤销口令
「450 撤」＝删本票与其证据件（不 reset/amend）。

## Progress log

### 2026-11-15 10:50 — 探针腿（agent-43）：AC#1／AC#2(a)／AC#3 有读数；AC#2(b) 具名留给机主

- 锚点 `caccccfc6`。**代码面与死亡那两发相同**：`git diff --stat cf3050401 HEAD -- cmd internal` **空** ⇒ 下面的读数不是"代码变了"造成的。
- **死没复现**：球（PID 49348）在跑，同一条载具命令串行两发，**都跑到底**：`RUN=861 / PASS=619 / FAIL=2 / SKIP=17`，**exit 1**，最后一行 `FAIL`，判行 `FAIL github.com/CarlosShao/wisp/cmd/wisp 1207.00s` 与 `1308.00s`。两枚 FAIL ＝ 本机具名已知红（`TestAC4FocusReturnToPriorWindowGap33r5`／`TestAC246DevLegIgnoresTheTestTaskInjection`），**未修未删未跳过**。
- **昨日两发的事件日志读数**：Application 日志（覆盖 2026-07-30 起、20,715 条）里 **0 条**提到 `wisp.test.exe`；对照：同机近 30 天有 **26 条** `Application Error` ⇒ 记录机制是开着的 ⇒ 昨日两发**没有 WER 记录**，即**不是**未处理异常／访问冲突／fail-fast。昨日两发的退出码与最后一行**没有留存**（stage-cli 日志已被 23:31 之后的运行覆盖），本票不猜。
- **AC#3**：`Test145ArgsJSONComesFromTheBookedForensicsRow` 单独 **PASS 0.22 s / exit 0**；`TestAC4PriorFocusSurvivesARefusedPanelSample` 单独 **PASS 1.27 s / exit 0**；两枚＋`TestAC4FocusReturnToPriorWindowGap33r5` 一起 `RUN=3 PASS=2 FAIL=1`（红的只有那枚已知环境红）。⇒ 这三枚**不是**死因载体；另：`Test145ArgsJSON…` 函数体里**没有窗**（不调 `startPanelForTest`），"死在 WebView2 区"是从文件序 `panel_*` 推出来的。
- **AC#2(b) 未做**（票面逐字禁止本腿 kill 球；"干净用户会话"替代本腿也拿不到）。交给机主一条命令：`.scratch\wisp\probes\450\run-balloff.ps1`（球还在就 **exit 3 拒绝**，本波已自测该拒绝支；判据写在该文件与证据 §3.3）。
- **AC#4 的"写进测试前置说明"没写**：它只在"球在跑死、球停不死"成立时才该写，而本波读数是"球在跑也能跑完"⇒ `scripts/wisp-cli-staged-tests.sh` **一个字未改**（连注释都未动）。
- 证据：`docs/evidence/s1/450-fullrun-death-probe.md`；原始读数 `.scratch/wisp/probes/450/**`（`runs.txt`／`run1.out.txt`／`run2.out.txt`／`samples.csv`／`mem.csv`／`ac3-readings.txt`）。
- ⛔ 未改任何产码；⛔ 未停／未重启球；⛔ 未 push。
---

## 探针收件（2026-11-15 11:0x，`d92e65d0a`＋`0595d4aa1`）：**死没复现**，我的第一嫌疑被推翻

**尺1**：两条提交在；★**我逐条核了它们报的"未改产码"** —— `d92e65d0a` 的 15 个文件全是 `.scratch/wisp/probes/450/**` ＋本票 ＋证据件，`0595d4aa1` 只改证据件；**我另外自己跑了** `git status --porcelain -- cmd internal scripts` ⇒ **0 条脏**；`scripts/wisp-cli-staged-tests.sh` **连注释都没动**。
**尺2**：`docs/evidence/s1/450-fullrun-death-probe.md` **15937 B**。
**尺3 我核的一条关键推断**：它说「"死在 WebView2 区"是从文件序推出来的」—— 我现查 `Test145ArgsJSON` 定义在 `panel_run_145_test.go:217`，**该文件里 `startPanelForTest|NewWindow|webview` 命中 0** ⇒ **这枚自己根本不建窗**。判语成立。

### 读数（腿跑的，两发都跑到底）
`RUN=861 / PASS=619 / FAIL=2 / SKIP=17`，exit 1，载具判行 `FAIL … 1207.00s` 与 `1308.00s`（20–22 分钟）。两枚红是**本机具名已知红**：`TestAC4FocusReturnToPriorWindowGap33r5`（5.20/5.06s）、`TestAC246DevLegIgnoresTheTestTaskInjection`（37.63/35.74s，红句 `a dev Wisp is already running in this session`）—— 未修未删未跳。
⇒ **`resident_*`／`ticket*`／`voice_*` 那一片今天有读数了**，昨日那个覆盖洞被填上。

### ★我的第一嫌疑（驻留球争用）＝**否**
球在跑时今天两发**都跑完**，昨晚 21:17–21:48 也跑完过一次。⇒ **不成立**。作为昨日死亡的触发因素，证据不足（AC#2(b) 未做）。

### 事件日志臂（这条最有信息量，对的是**昨日那两发**）
Application 日志覆盖 2026-07-30 起、**20,715 条**，其中提到 `wisp.test.exe` 的 **0 条**；昨日 21:50–23:45 只有 7 条无关条目（6 条 Security-SPP ＋ 1 条 igccservice）。对照：同机近 30 天**有 26 条 `Application Error`**（WebHostComponent.exe 等）⇒ **记录机制是开着的**，缺失不是"日志关了"。
⇒ 昨日两发**不是未处理异常／访问冲突／fail-fast，是静默退出**。这同时**削弱"WebView2 不稳"**：WebView2 host 的 fail-fast 会写 WER。
⇒ 但昨日两发的**退出码与最后一行没有留存**（stage-cli 日志被 23:31 之后的运行覆盖）⇒ **本票不猜**。

### AC#3 二分（无整包负载）
A1 `Test145ArgsJSON…` 单独 → **exit 0，PASS 0.22s**；A2 `TestAC4PriorFocus…` 单独 → **exit 0，PASS 1.27s**；A3 死亡点三枚一起 → `RUN=3 PASS=2 FAIL=1`（唯一红是已知环境红 4.98s）。⇒ 这三枚**不是死因载体**。

### 我接受的一条更正（对我的判断）
「死落在 WebView2 区」这句**是从文件序推的**，不是从"那枚测试自己开窗"来的。⇒ 票题里那个定位**我降级为待证**。

### 未翻的格（照实）
- **AC#2(b)「球不在跑」未做** —— 票面禁 kill 机主的球。腿交出 `run-balloff.ps1`（球还在就 exit 3 拒绝，本波已自测 `REFUSED … PID 49348`）。
- **AC#4 未翻** —— 它成立的前提是二分做完；现在**死没复现**，没有"最小复现"可交。⇒ 按铁律 ⓪⓪⓪ 我不写"球争用"的规矩（读数反过来），也不让腿猜着改产码。
- **本票当前状态**：**现象不复现、非球争用、非未处理异常**；要定性只剩"下一次复现时抓住退出码"。

### 对机主的那句话作废
我上一条让你「从托盘退出把球关掉，好让 450 的二分做完」—— **现在不需要了**。读数说球在跑不影响整包跑完。**别为了这张票去关球。**
---

## 补记（2026-11-15 13:2x，来自票 314 的 AC#0 取证 `cf36d8fef`）：**昨天缺的那件现在有了**

本票昨天结的账是「死没复现、非球争用、非未处理异常，但要定性只剩"下次复现时抓退出码"」。**那次复现来了，退出码也抓到了**：
- **dec `-1073741819` ／ hex `0xC0000005`** ⇒ **native ACCESS_VIOLATION**（⛔ Go panic、⛔ timeout、⛔ `os.Exit`）
- **复现率 1/2**（同形两发：一发在第 554.5 s 死、一发 1389.8 s 跑完 892 枚）
- 死点＝**面板线程干净退出之后、verdict 写出之前的收尾段**（stderr 末行 `panel thread exited cleanly … disposals=0`），**不是建窗途中**
- 与昨天「Application 日志 0 条」**不矛盾**：AV 未被 WER 捕获（无 dump）⇒ 这正是"静默退出"读数的物理解释

⇒ **本票的现态更新**：不是"现象不复现"，而是**条件性 native AV（1/2）**，成因在**收尾/拆卸段**，而不是驻留球、不是 WebView2 建窗。
⇒ 剩下的活是**点名故障模块**（要 LocalDumps 或调试器 ⇒ 要机主同意占机），以及**让收尾段可判**（D38 teardown 顺序）。⛔ 本票仍不许改名结案。

## 编排者收件·AC#2 勾＋诊断结论（2026-11-21 12:5x，件 `docs/evidence/s1/450-death-localization-r1.md`，commit `8333de280`）

**尺1** 提交在（我核过：该 commit **只碰该 evidence 件**，268 行；产码/用例/脚本 `git status --porcelain -- cmd internal scripts` 空）。**尺3** 诊断腿自跑 14 发（我未复跑整包，理由：整包 20 分钟且本票是定位票；我复跑了它引用的关键臂见下）。
**AC#2 两形读数（同一条载具命令）**：
- **(b) 球不在跑**：今天 14 发里 **6 死 8 活**（`RC=139` 静默消失，零 WER、零 `panic:`、零 `test timed out`）。
- **(a) 球在跑**：2026-11-15 那波（球 PID 49348 全程在跑）**两发整包都跑到底** `RUN=861 / PASS=619 / FAIL=2 / SKIP=17`、`exit 1`（两枚红均已具名：前台被占 + `TestAC246DevLegIgnoresTheTestTaskInjection`）。
⇒ **死与球无关**（球停着也死、球跑着能跑完）。
**定位（本票的真收获）**：死亡点＝`TestAC4PriorFocusSurvivesARefusedPanelSample`（`cmd/wisp/panel_host_windows_test.go:1042`，`--- PASS` 从不出现），尾部成片 `MARK299 wndproc … known=true`；**最小可复现集合＝`{TestPanelHostRealWindowHopAndLifecycle, TestAC4PriorFocus…}` 两枚真窗用例**（单跑各 0 死；两枚一起 6/14 死）—— 两枚**共用同一 WebView2 profile 前缀 `wisp-33r1-panel-profile`**（4 处调用点）。性质：**原生侧 AV／次序共享态**；(b) WebView2 泄漏已用读数排除（全机 31 枚、wisp 自家 1 枚、`added=[] mine-died=[]`、杀自家树不改变死亡率）；(c) onnxruntime 冲突已排除（pin 字节数在位、语音族全绿）。
**AC#4 已派落修腿**（存档 `.scratch/wisp/probes/orch/20261121/txt/450-fix-dispatch.txt`）：①两枚 profile 隔离 ②载具加**终点锚**（整包末尾没命中哨兵就响亮判红——今天 `RUN=412` 看着像跑完其实中途死了）。
**撤销**：本票撤＝删 AC#2 这一勾＋删本节（证据件保留）。
## 编排者收件·AC#4＋AC#5 勾＋票结（2026-11-21 13:3x，件 `docs/evidence/s1/450-profile-isolation-r1.md`，commit `851b5a144`）

**AC#4 走的是票面"若不是 ⇒ 交回可判事实与最小复现"那一支**（不是球争用）：真因＝**两枚真窗用例共用同一 WebView2 profile 家族**，最小复现集合 `{TestPanelHostRealWindowHopAndLifecycle, TestAC4PriorFocus…}`（单跑各 0 死／一起 14 发 6 死）；交回对象＝**票 314**（其 AC#1 明文写着"修法候选＝夹具改用每发独立 profile"），落修已在 `851b5a144`，我复核：**该 commit 只碰 `cmd/wisp` 的 `_test.go` 三枚＋载具脚本＋证据件，零生产码**。
**踩线如实记（我的派单错）**：本票"写面"那行写的是 `scripts/wisp-cli-staged-tests.sh`（**仅注释**），而 `851b5a144` ①把载具加了**终点锚逻辑**（不止注释）②改了 `cmd/wisp` 三枚 `_test.go`。技术上都对（且 314 AC#1 允许），但**越了本票写面**——我把这两笔的归属记在票 314 名下，本票只认"交回"。
**落修效果（腿自跑，我给读数口径）**：隔离后定向两枚连跑 **3 发全活（3 发 0 死，逐发 rc=0）**；**载具终点锚正控**（仓外临时副本把哨兵名写错）逐字红 `CARRIER_INCOMPLETE: died before the last test`、exit 1。
**★但整包仍会死，只是换位置了（不埋）**：改后真载具整包一发 `RC=139`（403s，`RUN=589/PASS=451/FAIL=0/SKIP=1`），先前致死的两枚在全包里**双双 PASS**，新死点＝`TestResyncOnShowDeliversFreshSnapshotFieldByField35g6`。⇒ **本票到此为止的产出是"死法变可判＋一对已隔离"，残死因交给票 314 继续**（它的 AC#1 还要做：35g6 走 `startPanelForTest`，而其 profile 用的是 `uniqueProfileLeaf`＝**本就每发唯一**，所以"共享 profile"这条假说解释不了新死点 ⇒ 314 AC#1 需换假说）。
**AC#5 门禁**：`d22scan` clean；`check-path-length-budget.sh` **VERDICT GREEN**（原 RED 是我改名没跟名册，已修，见台账）；`go build ./...` rc=0；`gofmt/gofumpt -l` 空。
**票结**：5/5。