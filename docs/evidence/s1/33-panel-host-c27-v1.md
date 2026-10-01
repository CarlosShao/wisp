# `33-v1` — 面板宿主真起来：非实现者对抗验收表（射程＝票 33 的 AC#1／AC#2／AC#3／AC#4／AC#11／AC#12）

- 程：`33-v1`（**对抗验收腿，非实现者**）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`（票 33）
- 被验腿：`33-r1`（写腿），五枚提交 `e4ed8304`→`697b4fae`→`b1e3d94a`→`fafe0445`→`f0ad3b51`
- ⛔ 它的证据件 `docs/evidence/s1/33-panel-host-c27-r1.md` 属**〔实现方自述＋台件〕**，**不构成任何一格 AC 的翻勾依据**；本表每一格只用我自己现跑的尺。
- ⛔ 本表**不勾票面任何一枚 `- [ ]`／`- [x]`**（勾框归编排者）；⛔ 不写"附条件通过"式补判语。
- 写作顺序（硬要求）：**§A／§B／§C 三节自对抗内容先写满并 commit**，之后才回填 §D 逐格判语。本表判据交付＝`wc -l -c` 这个文件。

---

## 0. 起手同发读数（`date`＋HEAD＋脏项，一次取）

| 尺 | 现量 |
|---|---|
| `date` | `Thu Oct  1 10:35:27 CST 2026` |
| `git rev-parse --short HEAD` | `94bd133e`（**起手锚**；HEAD 会动——只读普查腿 `244-a3` 在往 `.scratch/wisp/probes/244/a3/census.md` 写文档件，所以收尾时 `git diff 94bd133e..HEAD -- cmd internal` 必须为空） |
| `git status --porcelain -- cmd internal` | **0 行**（起手名册＝空；终态必须等于这一枚名册） |

---

## A. 我跑了哪些尺、每条真实读数（逐条现跑，无一条抄编排者或抄实现腿）

1. **起手同发三读数**：见 §0，时刻 `10:35:27`、HEAD `94bd133e`、`cmd internal` 脏项 **0 行**。
2. **被验的五枚提交逐枚复认**（`git show --stat --format="%h %ad %s" --date=format:"%H:%M"`）：
   - `e4ed8304`（docs，`1 file changed 53 insertions`）＝证据件 ⑥⑦⑧ 起手件；
   - `697b4fae`（`10:10`）＝`cmd/wisp/panel_host_windows.go` **+456**／`cmd/wisp/panel_host_gate_test.go` **+126**／`cmd/wisp/panel_host_windows_test.go` **+233**／`go.mod` **+2**／`go.sum` **+6**，共 823 增 0 删；
   - `b1e3d94a`（`10:26`）＝同三枚测试/产码文件的 **112 增／63 删**（专用线程 bring-up 从产码移进 `_test.go`）；
   - `fafe0445`（`10:26`）＝**只碰** `internal/panel/composer_dispatch_test.go`（**97 增／32 删**）＝那枚反转钉单独一枚提交；
   - `f0ad3b51`（docs）＝证据件 69 增／2 删（表里删除列非 0 的是它自己那份文档的 2 行，不是台账追加件 ⇒ 与"追加式文件删除列须为 0"那条纪律不同族，记下来备查）。
3. **实现腿的写面范围**（`git log --oneline 7a41db9b..HEAD -- cmd internal`，`7a41db9b`＝实现腿自报的真实起手 HEAD）＝**恰三枚**（`697b4fae`／`b1e3d94a`／`fafe0445`），与它自述一致；`git rev-parse --short 7a41db9b` 复认存在。
4. **票面 AC 框归零自证**（`git diff --numstat 7a41db9b..HEAD -- .scratch/wisp/issues/33-panel-host-c27.md`）＝**该路径零行输出**＝实现腿对票面**一字未动**；现量票面 `grep -c '^- \[ \]'`＝**12**、`grep -c '^- \[x\]'`＝**0**。⇒ "AC 框一枚没勾"这一条我复认成立（不是它自述）。
5. **裁定与预检的出处我已读**：票面「编排者裁定 J1–J8」全表（`:20-31`）＋新判据 AC#11（`:35`）／AC#12（`:37`）原文；台账 `A489`（`:10120` 起，反转裁与三处会被打假的指认）／`A490`（`:10153` 前，HEAD `7297e521` 上"设计内的红"）／`A491`（`:10153`，票 248 预检，与本程只共享"33 独占 `cmd/wisp`"这一句）。**两处编排者自认的错（`A490` 的文件名写错／`A491` 的措辞缺口）我不重复立案。**
6. **写面形状起手现量**：`cmd/wisp/panel_host_windows.go`＝**411 行**（`697b4fae` 落 456 行、`b1e3d94a` 净 −45）；`cmd/wisp/panel_host_gate_test.go`＝**180 行**；`cmd/wisp/panel_host_windows_test.go`＝**273 行**；`internal/panel/composer_dispatch_test.go`＝**718 行**。⇒ "宿主产码 456 行"那枚自述已过期，现读 411。
7. **工具链前提**（起手 `go env GOFLAGS` 空、`GOMODCACHE`、`GOPROXY` 见 §A#10 同发读数）＋ 基线命令必须带 `PATH=$PWD/third_party/sherpa-onnx:$PWD/build:$PATH`，否则 `0xc0000135`＝用例根本没跑（这条不是我的假设：票面与台账多处记录同形坑，我起手第一发就带上，并在红名册里核对有没有 `--- FAIL` 行）。
8. **三枚冻结件起手未动自证**（`git log --oneline 7a41db9b..HEAD -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go`）＝**零行**＝实现腿没碰那三枚（它的 `composer_dispatch_test.go` 不属冻结名册，反转属编排者在 `A489` 已裁）。

### 基线／归因（9–10）

9. **起手基线一发**（命令＝§0 那三读数之后立刻跑的 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp ./internal/panel -count=1`；起钟 `10:38:44`、收钟 `10:43:19`，HEAD `94bd133e`）：
   - `ok github.com/CarlosShao/wisp/cmd/wisp 201.507s`＝**零枚红**；
   - `FAIL github.com/CarlosShao/wisp/internal/panel 1.474s`＝**恰 4 枚红**，逐名（`--- FAIL` 计数＝4，有 `--- FAIL` 行＝用例真跑了，不是 `0xc0000135` 那形）：
     `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`。
   - 红句逐字（只用红句定归因，⛔ 未读 `frontend/**`／`design/**` 内容）：`approval_test.go:129 Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`；`composer_test.go:74 Go ComposerState emits [git currentModel modelKnown] that interface ComposerState does not declare`；`frontend_hygiene_test.go:216 a second style source appeared in code the bundle actually ships`；`tokens_fourway_test.go:441 read design/assets/tokens.css: open …: The system cannot find the path specified.`
   - ⚠ 第四枚的载体是**冻结件** `internal/panel/tokens_fourway_test.go`（三枚冻结件之一）⇒ 我只判归因、⛔ 一字不动、⛔ 不为它变绿动任何东西。
10. **归因（逐名判"是不是这程造成的"）**：四枚红的**测试文件与 Go 侧输入全不在本腿写面**——`git diff --numstat 7a41db9b..HEAD -- internal/panel`＝只 `composer_dispatch_test.go 97/32`；四枚文件最后修改逐枚：`internal/panel/approval.go 63ef895a`(09-21 13:26)／`approval_test.go f4bf0fae`(09-21 19:02)／`composer_test.go 0e5c0d4a`(09-28 18:18)／`tokens_fourway_test.go 63ef895a`。`tokens_fourway` 那枚的红因直接是工作树删除：`git status --porcelain -- design` 里 **16 枚 ` D`**（含 `design/assets/tokens.css`）；`frontend` 侧另有 4 枚 `M`＋11 枚 `??`。⇒ **4 枚红全非这程造成**，我不算新增红、不放宽、不修。实现腿 ④ 自报的正是这 4 枚同名 ⇒ 它没谎报红名册（这条我复认）。

### 攻击面①：反转钉（11–13）

11. **旧名只剩注释与归档**：`internal/panel/composer_dispatch_test.go:17`／`:396` 两处注释（＋`:432` 那句"old name now lives only in archived gate logs"）；活仪器名册已换＝`TestPanelHostIsAttachedAndNamesTheWindowHops`（`:433` 函数定义）。全仓 grep 旧名 `TestSliceAAttachesNoHost…`＝**只命中票面／`.scratch/wisp/probes/**` 归档日志／别家证据件／台账**，无任何 `scripts`／`tools`／`*.go` 活取数点 ⇒ 编排者 `A489` 那句"全仓无任何活仪器按旧名取数"我独立复认成立。⛔ 我没把旧名加回来。
12. **命中名册现量**（`go test ./internal/panel -run TestPanelHostIsAttached… -v`，`10:44:36`，HEAD `94bd133e`）：真树 **4 命中**＝`cmd/wisp/panel_host_windows.go:55:imports github.com/jchv/go-webview2` ＋ `go.mod:19:dependency` ＋ `go.sum:9/10:dependency`（`(production=true dependency=true)`）。⇒ **标识符族命中＝0 枚**：这棵树的 production 半边**完全由那一行 import 路径**供给（`hostSignals` 把 `:imports` 记进 production 那支，`:546-556`）。三枚原正控仍按期望响：A 假宿主源件＝**5 命中**、B 假 webview 依赖＝**2 命中**、C CLI-seam 那形＝**0 命中**；反向正控 D（中和宿主标识符）＝`0 hit(s) (production=false dependency=false)` PASS。载体全在 `t.TempDir()`，`go test -overlay` 零枚 ✓。
13. **定向突变 M1（我自己跑的"反转后还有没有牙"）**：`10:46:26` 起手 `md5 e75b9b2e36772ca4cae6043a6dbe7384`，把 `cmd/wisp/panel_host_windows.go:55` 的 import 路径中和（`go-webview2`→`go-xnone2`，一枚 token，一次只改一枚）⇒ 同钉**必红**：`--- FAIL: TestPanelHostIsAttachedAndNamesTheWindowHops (0.17s)`，红句在 `:439`（"no native host / WebView2 message channel is attached…"），读数 `3 -> [go.mod:19 …, go.sum:10 …, go.sum:9 …] (production=false dependency=true)`（dep 半边红不掉＝`go.mod` 是禁区我没动，这恰好也证明**两半是分开判的**）。`10:46:49` `git cat-file blob HEAD:<path> > <path>` 还原 ⇒ md5 与起手拷贝逐字节相同、`git status --porcelain -- cmd internal`＝**0 行**、同钉复跑 `ok internal/panel 0.178s`。⛔ 突变体零枚进提交。
   - **残余缺口（不用突变就能量）**：钉绿不等于有人开窗。#14 现量＝`PanelManager` 生产调用者**0 枚**，而这枚钉在那一态下是绿的 ⇒ 它的牙口只咬"这棵树链接了宿主"，不咬"用户点得开"；它自己的 doc `:413-418` 也是这么写的（没有越界宣称）。

### 攻击面②：宿主 ≠ 接进常驻进程（14–18）

14. **(a) 今天用户有没有任何一条路真能开出面板窗口＝没有**。尺：`grep -rn "PanelManager" --include=*.go .`（去 `scripts/`）＝**只两枚文件**：定义处 `cmd/wisp/panel_host_windows.go` ＋ `cmd/wisp/panel_host_windows_test.go`；`NewPanelManager` 的调用点全仓**唯一一枚**＝`panel_host_windows_test.go:158`（测试）。`grep -rn "PanelManager|panel_host" cmd/wisp/main.go cmd/wisp/run.go cmd/wisp/resident*.go`＝**0 命中**。`cmd/wisp/main.go:89-124` 的分派名册逐枚＝`run/providers/doctor/secret/models/slo/panel-assets/panel-inbound/version/help`，**没有一条起窗口**；`go test` 里起得来的那扇窗（HWND `0x6f0a60` 等）不是产品里点得开的东西。⇒ 与编排者 10:33 那尺（`PanelManager` 三文件 0 命中）同向，我自己重跑取到同一读数。
15. **常驻那一跳没接**：`internal/ball/ball_windows.go:54-55` 有 `OnPanelHotkey`/`OnTrayPanel` 两枚回调字段、`:665/:672/:677` 真 fire（在 ui-sta 的 wndproc 里）；`cmd/wisp/resident_ball_windows.go:127-128` 的回调体逐字＝`recordBallGesture("panel-hotkey")`／`recordBallGesture("tray-open-panel")`＝**只打点，不建窗**。另 `cmd/balldebug/main.go:199` 至今是 `tray: open panel (stub, ticket 33)` 桩。⇒ "宿主真装起来了"与"接进常驻进程"是两件事，本程只做了前一件（这一点实现腿具名说了，我没抓到老早交回，抓的是它的**边界与凭证**，见 #16/#18）。
16. **(b) 专用线程为什么从产码移走（`b1e3d94a`）**：`git show 697b4fae:cmd/wisp/panel_host_windows.go` 第 **433** 行确有 `go func()`，而且**带着** `// owner: this panel host thread; recover below`＋`defer recover()`＋`runtime.LockOSThread`；我把那一版放回工作树（一枚文件，`cp` 备份后 `git show`，读完立刻还原）跑 `sh scripts/d22scan.sh` ⇒ `--- FAIL: TestScannerSelfScanOfRealRepoIsGreen`，红句逐字 `cmd/wisp/panel_host_windows.go:433: [bare-goroutine] bare 'go func(' is banned (D22/D38b): use observe.Registry.Spawn (named, owner, recover boundary)`；仪器谓词在 `tools/d22scan/main.go:700-711`＝**任何** `go <anything>` 都算（owner 注释与 recover **不豁免**，唯一豁免是按文件路径的 allowlist），而 `_test.go` 在 ban #1–5 的走树里被跳过（`main.go:665`）。⇒ 它自述的"撞 ban #1"**成立**，且合规出路只有 `observe.Registry.Spawn`＋**名册内**协程名，而 J1 明令不许新增名册外名 ⇒ 移进测试脚手架是当时唯一不违 J1 的形状。还原后 md5 三处相同（`e75b9b2e…`）、`git status -- cmd internal`＝0 行（`11:01:07`）。
17. **线程形状与注释过期**：`grep -n "LockOSThread" cmd/wisp/panel_host_windows.go`＝**只有 `:101` 那行注释**（"see bringUp's runtime.LockOSThread"），函数体里没有；真的 `runtime.LockOSThread` 在 `panel_host_windows_test.go:48`。⇒ `b1e3d94a` 之后那句注释成了过期指认（`panelTitle` 头上 `:64-66` 那段也还留着 `panelWindowClass`/`panelOrigin` 两个**不存在的常量名**）。现在产码里 `bringUp` 的契约是"跑在调用方所在线程"（`:24`／`:164-167`），而 #14 现量＝**没有那样的调用方**。
18. **(c) 那句 panic 是跑出来的还是推出来的**：它表里 `②#2`／`⑧` 给的是**栈形状串**＋"本机现量"字样，并自陈**那枚用例文件未提交**（`⑧` 逐字"我一度写了这枚再入用例……未提交该文件"）。我扫盘找凭证：`.scratch/wisp/probes/33/r1/`（mtime 09-28）只有 `composer_dispatch.go.pristine`／`mutate.sh`／`logs/`；全仓 `grep -rl "panic: runtime error"` 只命中它自己那份证据件（**没有原始输出、没有复现命令**）。⇒ 判**〔仅自述＋栈形状〕**，登记为下一程前置（要么复跑留栈文，要么按别的架构裁）。它依赖的**库侧前提**我独立核到了：`pkg/edge/chromium.go:96-111` 的 `Embed` 确是 `for { GetMessageW → Translate → Dispatch }` 直到 `inited` 的**阻塞嵌套泵**，`chromium.go:130-131` 的 `Init` 直接解引用 `e.webview.vtbl`（nil ⇒ panic），`webview.go:340-344` 是 `CreateWithOptions → Embed → browser.Resize()`。⇒ "非可重入"这一半有源码依据；"从球泵再入必炸"这一半我没跑（要动 `internal/ball` 写面，超出我只读＋突变射程，见 §C）。

### 攻击面③：两枚 SLO 数自己复跑（19–21）

19. **11 枚干净读数**（同一条用例 `TestPanelHostRealWindowHopAndLifecycle`，四发跑次：`-count=1`／`-count=2`／`-count=3`／`-count=5`，全部带 sherpa PATH、`GOFLAGS=`，HEAD `94bd133e`；阈值一字未动）：
   - 冷（ms）：`937.093`(10:45:13)／`850.154`(10:45:40)／`832.173`(10:45:41)／`778.987`(11:02:21)／`667.225`(11:02:23)／`683.352`(11:02:24)／`739.195`(11:02:56)／`760.377`(11:02:58)／`805.933`(11:02:59)／`725.547`(11:03:00)／`649.102`(11:03:01) ⇒ **中位 760.377／最大 937.093**；
   - 热（ms）：`35.594`／`39.665`／`45.426`／`37.825`／`40.506`／`38.033`／`40.718`／`34.405`／`34.523`／`26.838`／`36.699` ⇒ **中位 37.825／最大 45.426**；
   - 断言位置＝`panel_host_windows_test.go:185-187`（冷 >1500 ⇒ `t.Errorf`）与 `:223-227`（热 >200 ⇒ `t.Errorf`），预算数字与 `docs/PLAN.md:508` 那句「冷 ≤1500ms / 热 ≤200ms」逐字一致；`git diff --numstat 7a41db9b..HEAD -- internal/observe/thresholds.go docs/SLO.md`＝**空**；P11（冷拉起 >2000ms 才重评 L2 卡）**未触发**（11 发最大 937.093）。
   - ⛔ 这族在 CI 无 `-tags winlive` 档（`A489`／`Q-75` 第四次具名），且用例靠 `//go:build windows` 平台分叉不进 ubuntu 编译 ⇒ **我的读数只有本机有效**，每条带时刻＋HEAD。
20. **票面 AC#2 那半句没落**：票面 `:67` 要 "10-run P50/P95 into SLO appendix"；交的仪器**每发只出 1 冷 1 热**，仓里既无分位汇总也无 appendix——尺：`grep -n "P50\|P95\|percentile" cmd/wisp/panel_host_windows_test.go cmd/wisp/panel_host_gate_test.go`＝**0 命中**（我上面那 11 枚是我自己攒的，不是它交的形态）。
21. **读数日志自带的假锚**：`panel_host_windows_test.go:181` 那句 `"(HEAD 7a41db9b, …)"` 里的 HEAD 是**硬编码字面量**，我 11 发都在 HEAD `94bd133e` 上跑，日志仍写 `7a41db9b` ⇒ 任何读者按日志认锚都会认错。

### 攻击面④：无监听端口那两层尺（22–25）

22. **L1 射程**：`panel_host_gate_test.go:26-42` 用 `go/parser` 只读**一枚文件**（`parser.ParseFile(fset, "panel_host_windows.go")`，禁 `net`/`net/http`/`net/netip`/`net/mail`/`net/tcp`）。我按 J7 的字面（"扫**宿主包** import"）另取一发更宽的：`grep -rn '^\s*"net' cmd/wisp/*.go`（去 `_test`）＝**0 命中** ⇒ 包级今天也干净，这个差别不改变结论、但射程与裁定文字不同形。链接闭包另说：`GOFLAGS= go list -deps github.com/jchv/go-webview2` 里 `net`＋`net/netip`（count＝2），`go list -deps ./internal/panel` 同样含两枚 ⇒ "宿主文件没 import net"⛔ 不等于"宿主可达闭包里没人能 listen"。
23. **M2 突变＝真开一枚端口，L2 不响**：`10:47:25` 在 `bringUp` 里加 `net.Listen("tcp","127.0.0.1:0")`（外加 `net` import ＋一枚全局持有）。结果：`--- FAIL: TestPanelHostOpensNoListeningSocketL1`（红句逐字 `panel_host_gate_test.go:38: panel host imports "net", which can open a listening socket…`）——**L1 响**；而 `TestPanelHostRealWindowHopAndLifecycle` **整条 PASS**，`-count=8` 那发 8 枚全 PASS。同发我用仓外 PowerShell 取 **ground truth**：`Get-NetTCPConnection -State Listen` 里 `Pid 28572 Name wisp.test Addr 127.0.0.1 Port 62971`＝**确有一枚 LISTEN**；不跑测试时该筛选＝0 枚（`10:48:19` S1／`10:48:33` S3）。⇒ **"本 pid 零条"在那一态下是假的**，L2 那半是恒真尺。
24. **为什么恒真（仓外 stdlib 探针，`C:/Users/swq/AppData/Local/Temp/33v1-l2probe/`）**：探针自己开 `127.0.0.1:0` 与 `[::1]:0` 两枚监听，然后按同一套 `GetExtendedTcpTable` 行解析读表（`10:50:21` 读数逐字）：
   - `AF_INET class=3 numEntries=42 state2(LISTEN)=42 state2ForOwnPid=1 state10(ESTAB)=0`
   - `AF_INET class=4 numEntries=174 state2(LISTEN)=0 state2ForOwnPid=0 state10(ESTAB)=0`
   - `AF_INET class=5 numEntries=216 state2(LISTEN)=42 state2ForOwnPid=1 state10(ESTAB)=0`
   - `AF_INET6 class=3 numEntries=27 state2(LISTEN)=3`（另 class=5＝4 枚）
   对照仪器（`panel_host_windows_test.go:114-141`）：它写 `tcpTableOwnerPIDAll = 4 // TCP_TABLE_OWNER_PID_ALL`（真值 **5**；仓里既有采样器 `internal/proc/treemetrics_windows.go:57` 逐字就是 `tcpTableOwnerPIDAll = 5`），并写 `tcpStateListen = 10`（10 是 `MIB_TCP_STATE_ESTAB`，**LISTEN＝2**）⇒ **两枚常量都错位**，它实际过滤的是"本 PID 的已建立连接"，任何监听端口都永远数不到。行偏移本身没错（rowSize 24／pidOffset 20 与 `internal/proc:58-59` 同形），所以**修正常量即有牙**——我的 class=5＋state=2 读数就是证据。另有两枚射程缺口：只传 `AF_INET` ⇒ IPv6 那一族整族不在场；重试耗尽时只 `t.Logf` 后 `return 0`（`:145-146`）＝工具失败也算过。
25. **子进程那一维**：这台机上不跑测试时 LISTEN 表里的 `msedgewebview2.exe`＝**0 枚**，跑窗口测试期间也＝0（S2 那发只有我们自己的 pid 一枚）⇒ "webview2 子进程自己开端口所以看不见"这一维今天**没有实际命中**，但仪器结构上只看 `os.Getpid()`（`:230`），修常量时得把"树内 PID"一起放进去——现成口 `internal/proc/jobscope_windows.go:157 TreePIDs()`。
   - 顺路核 `AGENTS.md` §1.2（⛔ 不扩射程，只核那两条）：`grep -n "filepath\." cmd/wisp/panel_host_windows.go`＝**0 命中**（没有 `risk.PathResolver` 之外的 `Clean`/`Abs` 决策，资源路径全走 `panel.Assets` 的字符串口）；明文凭据＝0（文件里唯二像配置的是常量标题 `"Wisp panel"` 与绑定名 `wispDispatch`）。另核两形顺手记下：裸 `go func(`＝0（#16/#17）；"墙钟时间差实现超时"＝0（`firstRoundTripLocked` 的 `t0.Add(5*time.Second)` 与 `time.Now().After(deadline)` 全程带单调读数，注释在 `:168-173`）。

### 攻击面⑤：AC#11／AC#12（26–29）

26. **AC#1 那枚"进程树"尺**：票面 `:65-66` 写的是 **process-tree** child-count；仪器 `countWebviewChildren`（`panel_host_windows_test.go:89-105`）是 `CreateToolhelp32Snapshot` 上按**进程名**全机数 `msedgewebview2.exe`，不按父 PID 走树。现量（`11:02:17`→`11:02:29`，tasklist 采样）：不跑测试＝**14 枚**（别人会话的 WebView2 进程）、跑我们一枚窗口期间＝**20 枚**（我们贡献约 6）、销毁后回 **14**。⇒ `before != after`（`:213-215`）与"≤2s 回落基线"（`:263-272`）都建立在 14 枚外部噪声上：别人起一扇窗能假红、别人退一扇窗能假绿。另外**全用例从未比较 hide→re-show 前后的 HWND 是否同一枚句柄**（`:175`/`:237`/`:260` 只判 0／非 0）⇒ 票面"reuses window"这一句没有按身份钉住；这一格里真有牙的是状态断言（`:206-208` `Hide` 后 `IsCreated` 必须仍 true＝hide 不 destroy；`:195-197` mode 处理器 `count()!=1` ⇒ Fatalf）。
27. **AC#4 焦点尺**：Show 后唯一断言是"前台窗非 0"（`:242-243`），**回还那一跳整条只有 `t.Logf`**（`:249-251`），且 `prior := windows.GetForegroundWindow()` 取在 `HotShow` 之后 ⇒ 我 11 发里 `foreground after Show == panel hwnd == foreground prior to test` 逐发同一枚句柄（`0x6f0a60`／`0x11fd0638`／`0x1008065e`／`0x1260a8e`／`0x509079a`／`0x124e08a2`／`0x1a4f057a`／`0x515092a`…）＝"上一次的前台窗"就是面板自己，`Hide` 把焦点还给它刚藏起来的那扇窗。尺：`grep -rn "prevFocus|SetForegroundWindow" --include=*.go cmd internal`（去产码文件）＝测试侧**只命中那句 `t.Logf`**（`:250`）⇒ 全仓零枚焦点回还断言。
28. **冷启路径会把页面盖掉（产码形状，非仪器）**：`bringUp` 先 `serveEntry()`（`:221` ⇒ `SetHtml(真入口字节)`），紧接着 `firstRoundTripLocked`（`:227` ⇒ `:372-375` 又一枚 `SetHtml`，内容是它自己写的探针文档）⇒ **冷启动结束后窗里显示的是探针文档而不是面板页**；`firstRoundTripLocked` 的调用点全仓唯一＝`:227`（grep 现量）。这一条同时削 AC#3 的"离线供给"与 AC#12 的"有页面可发"的现实意义。
29. **AC#11 我自己复尺＋我自己跑的反向证**：`git archive --format=tar HEAD | tar -x` 到仓外 `C:/Users/swq/AppData/Local/Temp/33v1-clean`（副本里 `find frontend/dist -type f` 恰一行 `.gitkeep`）⇒ `GOFLAGS= go build ./...` **exit 0**（`10:52:16`→`10:52:27`；`10:52:43`→`10:52:46` 复跑再一次）；**反向证**：只在**副本**里把 embed 目标改名（`mv frontend/dist frontend/dist-neutralised-33v1`，零删除、⛔ 未改任何前端文件一字）⇒ `go build ./...` **exit 1**，红句逐字 `frontend\embed.go:19:12: pattern all:dist: no matching files found`（`11:01:21`→`11:01:25`；跑完改回）。⇒ **AC#11 这枚判据有牙齿**，票面 `:35` 那句"判据自带反向证"由我在仓外满足，不依赖实现腿（它表里把这枚反向证推给"frontend 的主人"）。附带两枚读数：① `GOFLAGS= go mod tidy -diff` 在该副本 **exit 1**，差异＝把 `github.com/jchv/go-webview2` 从 `// indirect` 块移进直接 require 块（它只跑了 `go get`、没跑 `tidy`）；tidy 另外会抹掉 `golang.org/x/sys v0.47.0` 两行，而那两行**起手就在**（`git show 7a41db9b:go.sum | grep -c "x/sys v0.47.0"`＝**2**）＝不是这程带来的；② 仓里/CI **没有** tidy 门（`grep -rn "go mod tidy" .github/workflows scripts/*.sh`＝0 命中）⇒ 不红，但谁跑一次 tidy 就造出一枚与本腿写面无关的 diff。
30. **AC#12 复尺（能力形，不问退出码／不问文案）**：同一枚 HEAD 两口径各跑一发产码调用口 `wisp panel-assets`（仓外副本＋本机工作树）：
   - 清检出（只有入库件）⇒ `wisp panel-assets: assets NOT BUILT`（exit 1，`10:55:01`）；`-check` ⇒ `wisp panel-assets: panel: embedded assets are not built (run npm run build in frontend/)`（exit 1，`10:55:03`）⇒ `Assets.Built()==false`、入口 `index.html` **不存在**；
   - 工作树（HEAD＋别人那 3 枚未跟踪产物）⇒ `panel assets embedded: 4 files, entry=index.html built=true`（exit 0，`10:55:19`）；`-check` ⇒ `entry=index.html built=true 2 asset refs resolve [./assets/index-BVKlegVD.js ./assets/index-BRKj5OIJ.css]`（`10:55:20`）；
   - 入库口径：`git ls-files frontend/dist`＝**1 枚 `frontend/dist/.gitkeep`**（与编排者现量同）。
   ⇒ **"只匹配到占位"与"真有一包页面产物"这两态确实可分**，而且分它的仪器是既有那枚：`internal/panel/assets.go:56`（`fs.Stat(tree, EntryFile)` 定 `built`）＋`internal/panel/assets_test.go:22 TestAnchorOnlyBundleIsNotBuiltAndFailsClosed`（**逐名复认存在**——按"注释里的测试名一律当待验断言"那条纪律查的）。⛔ 我没有造任何页面产物、⛔ 未读前端内容（只数条目名与存在性）。
31. **腿交的那枚 AC#12 用例是空尺**：`panel_host_gate_test.go:52-76` 全文只有 `t.Logf`——非锚条目 `t.Logf`（`:73`）、计数 `t.Logf`（`:75`），除 `git ls-files` 报错（`:62` `t.Fatalf`）之外**没有一条 `t.Errorf`** ⇒ 它注释里那句"if a real bundle is ever committed, fails closed on a half-bundle"在代码里**没有对应物**。⇒ AC#12 的牙在 `internal/panel` 那侧（#30）与产码 fail-closed（`serveEntry` 在 `!Built()` 时返回错误 ⇒ 窗里根本没有页面），不在这枚用例里。
32. **CI 分母具名（这一格是我读出来的错位）**：`cmd/wisp` 只在 **windows cli 腿**跑——`scripts/portable-tests.sh:195` `scope=(./cmd/wisp/)`；`ci.yml:475` `bash scripts/wisp-cli-tests.sh` 落在 `test-windows` 作业（`ci.yml:388 runs-on: windows-latest`）；`scripts/wisp-cli-tests.sh:65` 自陈"this is the windows leg of the cmd/wisp gate"。`ci.yml:341` 的 ubuntu 腿（`ci.yml:278 ubuntu-latest`）跑 `--scope=core`，那份清单（`portable-tests.sh:169-180`）里是 `./internal/panel/...`（`:179`），**没有 `./cmd/wisp/`**。⇒ **L1／AC#11／AC#12 三枚用例都不在 ubuntu core scope**，与 `panel_host_gate_test.go:3-6` 文件头那句 "These run in the ubuntu core scope too (J7: the L1 half must have a real CI denominator)" 矛盾＝一句过期指认（⛔ 我不动文件，只登记）。反转钉本身在 `internal/panel` ⇒ 它在 ubuntu core **有**分母 ✓，且它的扫描不看 `//go:build` 约束，所以在 ubuntu 上同样从 `panel_host_windows.go:55` 取到命中。
33. **依赖与哈希（J5／ban #5）**：`git diff --numstat 7a41db9b..HEAD`＝`go.mod 2/0`、`go.sum 6/0`（我另跑一遍逐字读 diff）。新增 6 行＝webview2 的 `h1:`＋`/go.mod`、winloader 的 `h1:`＋`/go.mod`，**外加 2 枚 `golang.org/x/sys` 老版本的 `/go.mod` 行**（`v0.0.0-20200810151505`／`v0.0.0-20210218145245`）。那两枚是**模块图解析的产物**（`go mod tidy` 保留它们，见 #29 的 comm 读数）——手抄 4 行不会产生它们 ⇒ 佐证它真走了 `go get`＋sumdb，没从镜像目录取件。⛔ 我全程没跑 `go get`/`go mod`（`go.mod`／`go.sum` 我一个字没动）。⚠ 两处读数要更正它表里的话：普查那句"4 行 go.sum"少计 2 行；`go.mod` 里 `github.com/jchv/go-webview2` 被记成 **`// indirect`** 而它是 `cmd/wisp` 的直接 import（#29）。
34. **我自己的门禁终态**（HEAD `94bd133e`，全部现跑）：`bash scripts/d22scan.sh` **exit 0**、逐字 `clean - no D22 ban violations; live scope work: bans #1-5 internal/=224, bans #1-5 cmd/=33, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=476, ban #8 cmd/=78`（与它 ④ 自述的分母逐字相同＝复认）；`GOFLAGS= go vet ./cmd/wisp ./internal/panel` **exit 0**；`/d/work/base/gopath/bin/gofumpt.exe -l cmd/ internal/panel/`＝**空**（`11:05:24`）。⛔ staticcheck 我没复跑：本机二进制是 `2025.1.1 (0.6.1)`，而 CI 钉的是 `2026.2.1`，且 `ci.yml:181-200` 注释自陈旧版在 go1.27 的 V4 export data 上**导入期就崩**（"cannot decode …, export data version 4 is greater than maximum supported version 2"）⇒ 跑它只会得到一枚假结论。
35. **冻结件与禁区自证（我的，不是引它的）**：`git diff --numstat 7a41db9b..HEAD -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go internal/observe/thresholds.go docs/SLO.md docs/BUILD.md docs/PLAN.md docs/specs`＝**全空**；`bridge.go` 入向常量枚数我重跑 `grep -cE '=\s*"panel\.' internal/panel/bridge.go`＝**4**；`internal/panel` 零平台分叉我重跑（无 `*_windows.go` 之类文件名、非 test 文件里无 `//go:build`）＝**成立**；本腿写面⛔ 未碰构建链（`git diff --name-only 7a41db9b..HEAD -- scripts/build.ps1 .github/workflows docs/BUILD.md`＝**0 行**，只有 `go.mod`）。

### 两枚"停手上报"的独立复认（36–37）

36. **上报①（`pkg/edge` 控制器定位）＝理由不成立**。它说的类型事实我逐字核到了：模块缓存 `…/pkg/edge/ICoreWebView2Controller.go:59` 确为 `func (i *ICoreWebView2Controller) PutBounds(bounds w32.Rect) error`，而 `w32` 在 `internal/w32/`（`package w32`，`w32.go:140 type Rect struct`）⇒ 主模块**确实 import 不到那枚类型**。但它由此推的"所以纯 edge 宿主给不了尺寸、只能渲染空白、只能退回高层 API"**不成立**，因为**根本不需要构造那枚类型**：`pkg/edge/chromium_amd64.go:12-22` 是导出的 `func (e *Chromium) Resize()`，它自己 `var bounds w32.Rect` ＋ `User32GetClientRect(e.hwnd)` ＋ 裸 vtbl 调 `PutBounds`；`pkg/edge/chromium.go:47` 有导出构造器 `NewChromium()`，`:292` 有导出的 `AddWebResourceRequestedFilter`；而高层 API 给控制器定位用的正是同一枚 `Resize()`（`webview.go:429` `w.browser.Resize()`，`webview.go:340-344` 创建链路也是 `Embed` 后直接 `Resize()`）。⇒ **这一枚不是依赖边界，是产品形状／依赖用法的决策**（要不要在产码里自建窗口层、吃 edge 层那枚独占泵）。它把"多文件页面那半做不了"报成硬边界，报重了。
37. **上报②（泵再入 panic）＝前提有源码依据、结论〔仅自述〕**。库侧可核的那半在 #18（`chromium.go:96-111` 阻塞嵌套泵／`chromium.go:130-131` nil 解引用／`resident_ball_windows.go:127-128` 回调确实只打点）。它没留原始栈、没留复现命令、且自陈那枚用例**未提交**（#18）⇒ 我判**〔仅自述〕**并登记为下一程前置。而且它按票面处置对了：J1 那一格逐字要求"证不出来就停手上报、⛔ 不许改成独立线程绕过"，它既没造第二线程（`grep` 现量：产码里没有 `LockOSThread`／名册里没有新协程名）也没把 `internal/ball` 的写面自己动了。⇒ 这一枚属**依赖／架构形状问题**（泵的归属权），不是我们产品的形状问题；要不要授权动 `internal/ball` 交泵权＝契约面（可能触 D38），归编排者摆 owner，我不选。


---

## B. 我可能写错的条目（对抗我自己）

1. **"go test 里起得来的窗口" ≠ "产品里点得开"**。本表最可能的错就是把宿主原语的真机读数（测试在专用锁线程上建的窗）直接读成"用户点设置看得见"。⇒ 判 AC#1/AC#2/AC#4 时每一格都要具名写清**调用者是谁**（`cmd/wisp` 产码里枚数几枚、非用例调用者逐枚点名），⛔ 两件事不许合并成一句。
2. **我会不会把 `internal/panel` 起手那 4 枚红算成这程造成的**。归因尺＝`git diff 7a41db9b..HEAD -- internal/panel`（只有一枚测试文件 97/32）＋红句源码点名；⛔ 我不读 `frontend/**`／`design/**` 内容（两层禁令），只用红句、`git status` 的 `D` 行与断言源码定归因。判语只写"是不是这程造成"，不为它们放宽任何断言，也不动三枚冻结件让它们变绿。
3. **定向突变的还原风险**。我只改产码一枚、读完红名**立刻** `git cat-file blob HEAD:<path> > <path>` 还原（⛔ 禁 `checkout`/`restore`/`stash`），还原后取 `md5sum` 与起手拷贝逐字节比对，再复跑同发改绿；突变体一枚都不进提交。若还原后 hash 不同 ⇒ 这一格我判"判不动"并具名上报，不留"应该没问题"。
4. **SLO 读数的自相关**：`-count=2` 的第二发可能吃到第一发的 WebView2 缓存/常驻子进程，把"冷"读成"不冷"或把"热"读得更热。⇒ 两发各带时刻＋HEAD，并具名写这族在 CI **无 `-tags winlive` 档**（`A489`／`Q-75`）⇒ 只有本机有效；⛔ 阈值一字节不动（`internal/observe/thresholds.go` 起手只读不改）。
5. **AC#12 的两种口径混用**：`find frontend/dist`（磁盘，含他人**未跟踪**构建物）与 `git ls-files frontend/dist`（入库）差得很远。判 AC#12 只用**入库/干净检出**口径；⛔ 我不造页面产物填 `frontend/dist`（那是界面侧那枚 agent 的活），也⛔ 不把 AC#11 的 `rc=0` 读成"有内容"。
6. **"无监听端口"那把尺可能恒真**：如果 L2 那半按**测试进程 PID** 过滤，而真开端口的是 webview2 的**子进程**，或它取 PID 的方式取错，那"本 pid 零条 LISTEN"是**必然成立**的空尺。⇒ 我必须读它的射程（它数的是哪个 PID、有没有把行数当命中），并做一发**真开端口的仓外载体/产码突变**证明它会响；光看"零条"绝不算凭据。
7. **我会不会把它的自述当成它跑出来的**：那句泵再入 panic（`edge.Chromium.Init` nil `webview`）我只有复认"证据件里有没有真栈原文＋可复现命令"的义务；若只有栈的形状没有复现尺 ⇒ 标〔仅自述〕并登记为下一程前置，⛔ 不替它追认也不替它推翻。
8. **载体必须在仓外**：正控/反向正控/清检出全部建在 `C:/Users/swq/AppData/Local/Temp/`（或 `t.TempDir()`），⛔ 零枚 `go test -overlay`，⛔ 绝不在仓内建 worktree／checkout。
9. **HEAD 会动**：`244-a3` 正在往 `.scratch/wisp/probes/244/a3/census.md` 写并 commit 文档件。⇒ 每条读数都带取数时刻＋当时的 `git rev-parse --short HEAD`；收尾必用**起手锚** `94bd133e` 跑 `git diff 94bd133e..HEAD -- cmd internal`（应空），而不是用"最新 HEAD"自证。
10. **反转钉的射程我可能读歪**：新语义是"产码标识符族 ≥1 命中 **且** 主模块 `go.mod` 有 webview 模块"。我要分清**谓词**（`hostChannelCapabilityHits` 那类）与**判定形**（`==0` 改成 `!=0`）分别在哪一行，并逐枚数正控的命中数；只看函数名改名就算"已反转"是空判。
11. **`AGENTS.md` §1.2 的两条顺手核不是我该扩大射程的口子**：读 `panel_host_windows.go` 时只核"明文 API 密钥"与"`risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策"两形；若发现别的越界形状，**具名登记**交编排者，⛔ 不自己判它的票、⛔ 不动票 248／票 228 的活。
12. **我会不会把"这格勾不了"写成"这格红"**：AC#12 与 AC#3 的过滤器半、AC#1 的 dispose 半是**归口未落定／依赖边界**，判语要写"证据形状成立否＋归口"，勾不勾归编排者。
13. **M1 我只主张得对它主张的那一半**：我用 import 路径中和来抓反转钉的牙，而 `go.mod` 属禁区我没动 ⇒ dep 半边在突变态下仍是 `true`，所以我只说"production 半边按期望红了（读数 `(production=false dependency=true)`）"，⛔ 不说"两半都响过"。真要两半都响，得允许动 `go.mod`（禁区）或由下一程在仓外副本里做。
14. **IPv6 那一半我不越结论**：探针里 `AF_INET6` 表确有 `state=LISTEN` 行（class=3 共 3 枚），但我按 48 字节行layout取的 own-pid 归属**没对上**（`state2ForOwnPid=0` 而我明知自己开了 `[::1]` 监听）⇒ 可能是我的 offset 错、也可能是 Windows 的归属口径。所以我只主张源码事实那一句"仪器只传 `AF_INET`、IPv6 整族不在射程"，⛔ 不主张"这台机上已经漏了一枚 IPv6 端口"。
15. **我在工作树里复现了历史版本**（d22scan ban #1 那一发）：放回的是 `697b4fae` 那枚**已提交过**的文件版，不是我的突变体，但我仍按突变纪律处理——一枚文件、先 `cp` 备份、读完立刻 `git cat-file blob HEAD:` 还原、三处 md5 相同、`git status -- cmd internal`＝0 行，暴露窗口 `11:00:48`→`11:01:07`（19 秒）。⚠ 残余风险：若那一瞬有别的腿跑 `cmd/` 门禁会读到它；台账 `A491` 声明 33-r1 已交、写面只有我在动，我把时刻写死在这里供回查。
16. **`go mod tidy -diff` 走的是网络元数据**：那一发在**仓外副本**里跑，会 `go: downloading` 若干模块的 `.mod`（GOPROXY 是本机既成设置的镜像）。⛔ 我没有把任何哈希写进仓内文件、⛔ `go.mod`／`go.sum` 一字未动（终态 `git status` 可核）；ban #5 指的是"把值抄进文件"，我这一发只读、无落笔。若编排者认为连"镜像取元数据"都算污染，请把这一条具名退回，我不自辩。
17. **我这 11 枚 SLO 读数不是它的交付物**：票面要的是腿交的 10-run 汇总（#20 现量＝仓里没有）。我把自己的样本记在 §A#19 并标"我自己攒的"，⛔ 不拿它替实现腿补一个"P50/P95 已交"。


---

## C. 判不动的地方（逐条甲／乙／不做＋现量）

1. **AC#1 的"session dispose  destroys ＋ WebView children exit ≤2s"那一半**：票面 `:65-66` 要求 dispose 触发。现量（起手 `grep -rn "PanelManager" cmd/wisp/main.go cmd/wisp/resident_windows.go cmd/wisp/run.go`）**0 命中**（编排者 10:33 同尺，我自己会再跑一发）⇒ 若真没有会话级 dispose 的生产触发者，那"会话结束拆窗"这半**在产品里不可达**，我只能判"显式 `Destroy` 机制＋子进程回落"这一半。**请裁**：甲＝按机制判这半、那半登记未落地；乙＝要求这程先接 dispose 跳（超出它自述射程）；不做＝不放宽、不宣称。
2. **J1 的第一枚用例（票面 `:24`：落地腿"第一枚用例就证泵期间球仍能出帧、消息仍被派发"）**：实现腿自述这一枚**未提交**（本机 panic，改落进 ②/⑧）。⇒ 本程能判的是"证没证／停手上报合不合规"（它自述证不出就停手、未造第二线程绕＝符合票面），但**"证不出来"这一句是它跑出来的还是推出来的**要靠复认（见 §C#3）。真裁"要不要授权动 `internal/ball` 交泵权"＝契约/别模块写面，**归编排者摆 owner**，我不选。
3. **那句 panic 的可复现性**：我要在它的表里找**真栈原文＋复现命令**；现读 `33-panel-host-c27-r1.md` ②#2 给的是栈的**形状串**（`ballWndProc→…→webview2.NewWithOptions→CreateWithOptions(webview.go:340)→Embed→Init(chromium.go:131)`）与时刻口径，**没有贴原始 panic 输出、没有复现命令行**。⇒ 判〔仅自述〕；下一程前置：要么复跑一发留下栈文，要么改成别的架构裁。⛔ 我不自己复现那发（要动 `internal/ball` 写面，超出我只读＋突变射程）。
4. **AC#3 的"多文件页面经 `AddWebResourceRequestedFilter` 服务"那一半**：它交回的是 `pkg/edge` 控制器 `PutBounds` 吃模块内私有 `w32.Rect`（`pkg/edge/ICoreWebView2Controller.go:59`）⇒ 主模块给不了尺寸 ⇒ 只能退 `SetHtml`。这枚我要独立复认（读依赖模块真实签名），并明确写**它是依赖边界问题还是产品形状问题**——这句决定编排者要不要摆 owner，我不替他选，但**必须给"如果成立，属哪一类"**。
5. **AC#4 的 manual 那一半**（票面 `:69`「automated ＋ manual」）：我没有真编辑器在场上，`GetForegroundWindow` 取到的是控制台前台。⇒ 机制可判（记录前窗＋回还调用），"手感"判不动，标〔仅本机机制层〕；⛔ 不许写成"焦点回还在真编辑器里成立"。
6. **AC#12 的勾**：产物归界面侧那枚 agent，本编队⛔ 不写 `frontend/**` ⇒ 判据能否区分"只匹配到占位"与"真有一包页面产物"我可以复尺，但**这一格在"谁把 dist 填上"落定前勾不了**（票面 `:37` 逐字），我把这件事写成"判据形状是否成立＋现量归口未落"，不当红、不当绿。
7. **`frontend/dist` 磁盘态那三枚未跟踪产物**（`git status` 里属他人地界）：我不动、不提交、不删、不读内容；它们存在这一事实只用来解释"本机能 build 出带页面的 exe"为什么**不是** AC#12 的凭据。⛔ 我不打开那些文件、⛔ 结论不引到 `frontend/**` 身上。
8. **winlive 档在 CI 零岗位**（`A489`／`Q-75`）：本表所有真机读数只在开机这台机器成立，每条带时刻＋HEAD；要不要单开 workflow＝契约级，已由编排者登记为默认"不做"，我不重开这一裁。

9. **那句 panic 我自己不复现**（甲／乙／不做）：现量＝`internal/ball` 里"把闭包投到 ui-sta"那枚原语**存在但不可用**——`ball_windows.go:312/317/325/333` 都在调 `b.sta.PostTask(...)`，而 `staThread` 是**未导出类型**、`b.sta` 是未导出字段 ⇒ 包外零枚公开口（实现腿那句"零枚"我复认成立，但它没说出"原语已在、只差一枚导出的包装"这一层，这会让编排者以为要大动 `internal/ball`）。⛔ 不做＝我不自己复现那发再入（要动别模块写面）；⛔ 不做＝我不替编排者选"授权交泵权／换异步创建绑定／接受独立线程改名册"这三条里的任何一条。甲＝本表按〔仅自述〕登记并作为下一程前置；乙＝要我现在复现则需追加只读腿授权。
10. **AC#2 该按哪种形状结**：票面 `:67` 的 "10-run P50/P95 into SLO appendix" 与交的"单发断言"（#20，`grep P50|P95`＝0 命中）不同形。我⛔ 不自己判它算不算满足——甲＝按"数值成立＋口径混（J3 三口径没分开报）＋本机有效"结，把汇总器推给下一程；乙＝要求本票补一枚 10-run 汇总器再判；不做＝不拿我这 11 枚替它充当交付物。**交裁。**
11. **J2 我只判到"入口自带 fatal"这一层**，AC#5 的夹具（rename/mask loader）不在我射程：`webview2.NewWithOptions` 在模块里带三枚 `log.Fatal`（`webview.go:115/120/125`），edge 层另带 `chromium.go:173`（环境创建失败＝运行库缺失那一支）与 `:295`；它表里只具名了两枚、且行号与我逐字读的差一枚（它写 `:189`，真身 `:188`；`:284` 对得上）。AC#5（no-crash 降级）⛔ 不在我射程 ⇒ 我判 J2 那句"不走那两枚会 fatal 的 API"**不成立**，但"这格要不要现在开"交裁——因为 AC#5 的夹具（rename/mask loader）是另一程的活。
12. **staticcheck 这道门我复跑不了**：本机 `/d/work/base/gopath/bin/staticcheck.exe 2025.1.1 (0.6.1)`，CI 钉 `2026.2.1`（`ci.yml:181-200` 自陈旧版在 go1.27 export data V4 上导入期即崩）⇒ 判"这一门未复认"，交 CI；甲＝就此收；乙＝要我装对版（属工具链动作，不在我只读＋突变射程）。
13. **AC#12 的勾**：产物归口（界面侧那枚 agent）未落定 ⇒ 票面 `:37` 逐字"这一格在'谁把 dist 填上'落定前勾不了"，我照办：判语只写"判据形是否可分＋现量＋腿那枚用例有没有牙"，⛔ 不勾、⛔ 不造产物、⛔ 不读前端内容。
14. **"面板能点"这句话我给不出任何凭据**：#14（0 生产调用者）＋#28（冷启末页被探针页覆盖）＋#30（入库无页面）⇒ 三件事叠加时，本票射程内**不存在**"用户看得见主面板"这条路径。owner 那句「起码我要能看到主面板，我要点击设置」要的是这条路径，它归 J1 那一格与票 248；我不把本票这六格读成"owner 的要求已满足"，也不反过来宣称"永远满足不了"（那是编排者的裁）。


---

## D. 逐格判语（六格；⛔ 本表不勾票面任何一枚框）

| 格 | 判语 | 凭据（全指本表现跑） |
|---|---|---|
| **AC#1** | **部分成立／判据不齐全：`⛔ 不给勾`** | #14/#26/#28，J1 行 |
| **AC#2** | **数值成立、形状与票面不同形：`⛔ 不给勾`（口径与交付形状交裁）** | #19/#20/#21，§C#10 |
| **AC#3** | **不成立（"无监听端口"那一维今天没有任何会响的检；离线供给那半未落且被自覆盖）：`⛔ 不给勾`** | #22–#25/#28/#36 |
| **AC#4** | **不成立（ automated 半是一枚任何焦点行为下都不会红的尺）：`⛔ 不给勾`** | #27/#14，§C#5 |
| **AC#11** | **成立（我自己清检出复建＋我自己仓外反向证）** | #29，#32/#35 |
| **AC#12** | **判据可分／现量＝入库没有页面／腿那枚用例无牙：`⛔ 今天勾不了`（归口未定）** | #30/#31/#28，§C#13 |

### 逐格展开

- **AC#1 生命周期**。真有牙的两条：`Hide` 后 `IsCreated` 必须仍 true（`panel_host_windows_test.go:206-208`，＝票面 "hide-don't-destroy" 那一半），以及 `Destroy` 后 `IsCreated`/HWND 归零（`:257-262`）；H3/H10 的处理器计数（`:195-197`）也真。三条不齐：ⓐ 票面要 **process-tree** child-count，仪器是全机按名数 `msedgewebview2.exe`（#26：基线 14 枚属别人会话、我们贡献约 6），别人起一扇窗能假红、退一扇窗能假绿；ⓑ **从未比较 hide→re-show 前后的 HWND 身份** ⇒ "second show reuses window" 这句没按身份钉住（#26）；ⓒ "session dispose destroys" 在产品里**不可达**——`Destroy` 的调用者只有测试（#14），常驻的 `OnPanelHotkey`/`OnTrayPanel` 至今只打点（#15），J1 那一跳未接。⇒ 判语＝宿主**原语**在测试线程上表现正确，生命周期**没接进进程**；两件事不许合并（这格的读数全部来自 `go test` 起的那扇窗）。
- **AC#2 冷/热延迟**。我这台上 11 枚读数最大冷 937.093／中位 760.377，最大热 45.426／中位 37.825 ⇒ 两枚预算（1500／200）都**未越界**，P11 未触发；阈值与 `internal/observe/thresholds.go` 我都复认**零改动**，没有为变绿放宽任何断言（#19）。三条不齐：ⓐ 票面要 10-run P50/P95 落 appendix，交的仪器每发只出 1 冷 1 热、仓里无分位汇总（#20）；ⓑ J3 要三口径分开报，交的是**一枚混口径的数**（那枚 cold 含 `SetHtml` 探针页往返＝`cold-embed` 那一形）；ⓒ 读数日志里的 HEAD 是硬编码字面量（#21），且 winlive 在 CI 零岗位 ⇒ 只有本机有效。⇒ 数值我认，交付形状交裁（§C#10）。
- **AC#3 离线供给＋无监听端口**。**这格我判不成立，而且它的"PASS"是仪器造的假 PASS。** L1（AST 扫单文件 imports）真会响——M2 突变下它逐字红（#23）；但 L2 那半**恒真**：同一枚突变下宿主进程确实持有 `127.0.0.1:62971` 的 LISTEN（PowerShell ground truth），用例却整条 PASS，根因是两枚常量都错位（class 该 5 写成 4、state 该 2 写成 10，而仓里 `internal/proc/treemetrics_windows.go:57` 就是 `= 5` 的对照；我 stdlib 探针给出 class=3/5＋state=2 时 own-pid 命中 1 枚）＋只读 `AF_INET`＋重试耗尽静默返 0（#23/#24）。"⛔ 光看我们的 PID 零条"确实⛔ 不等于没开端口，而且这里连"看"都没看对。供给那一半另有两笔：`AddWebResourceRequestedFilter` 在产码里**零使用**（票面 `:49-50` 那句没落，J6 只免了路径名、没免能力），只 `SetHtml` 单文档；而 `bringUp` 在喂完入口字节后又用探针文档 `SetHtml` **把它盖掉**（#28）⇒ 冷启结束显示的既不是真页面也不是"没有页面"，是一枚自造探针；清检出口径下则连入口字节都没有（#30）。它为此交回的停手上报（`w32.Rect` 不可 import ⇒ 多文件页做不了）**理由不成立**：`(*edge.Chromium).Resize()` 是导出的、高层 API 自己就用它给控制器定位（#36）⇒ 那是**依赖用法／产品形状**决策，不是依赖边界。
- **AC#4 焦点回还**。Show 后唯一断言是"前台窗句柄非 0"，回还那一跳**只有 `t.Logf`**；全仓 `prevFocus` 零枚断言（#27）。结构上还有一层：`prior` 取在 `HotShow` 之后，⇒ 我 11 发里"上一次的前台窗"逐发就是面板自己，`Hide` 把焦点还给它刚藏起来的那扇窗——这条用例在任何焦点行为下都不会红。加上 #14（用户点不开这扇窗）⇒ 票面那句 "automated ＋ manual" 两半都没有凭据。机制代码存在（记录句柄＋`SetForegroundWindow`/`SetFocus`，`panel_host_windows.go:269/311-313`），**证据为零**。
- **AC#11 清检出能建**。**成立，且这格的判据有牙齿——是我自己证的**：仓外 `git archive HEAD` 副本 `GOFLAGS= go build ./...` exit 0（两发），把 embed 目标在**副本**里改名（零删除、未动前端一字）⇒ exit 1 红句 `frontend\embed.go:19:12: pattern all:dist: no matching files found`（#29）。用例本体也是真断言（build 失败 ⇒ `t.Fatalf`）。三条附带读数：⚠ 三枚 gate 用例（L1/AC#11/AC#12）只在 **windows cli 腿**跑、不在 ubuntu core scope，与文件头 "These run in the ubuntu core scope too (J7…)" 矛盾（#32）＝一句过期指认，登记不改文件；⚠ `go mod tidy -diff` 在 HEAD 上 exit 1（direct/indirect 归位差一行；两枚起手就有的 `x/sys v0.47.0` 行会被 tidy 抹掉，不是这程带来的），仓里 CI 没有 tidy 门 ⇒ 不红但会造 diff（#29/#33）；⚠ 判语范围＝字面的"能建"，⛔ 不得读成"有内容"——那正是 AC#12 立的理由。
- **AC#12 能建≠有页面**。判据**可分**（能力形，不问退出码）：清检出 `wisp panel-assets` ⇒ `assets NOT BUILT`／`embedded assets are not built`（两发 exit 1，入口不存在），工作树 ⇒ `4 files, entry=index.html built=true`；入库口径 `git ls-files frontend/dist`＝1 枚 `.gitkeep`（#30）。分这两态的既有仪器在 `internal/panel/assets.go:56`＋`assets_test.go:22 TestAnchorOnlyBundleIsNotBuiltAndFailsClosed`（逐名复认存在）。但**腿交的那枚用例是空尺**：`panel_host_gate_test.go:52-76` 全文只有 `t.Logf`、零枚 `t.Errorf`，注释那句 "fails closed on a half-bundle" 在代码里没有对应物（#31）⇒ 它不构成这格的凭据；产码侧真会 fail-closed（`serveEntry` 在 `!Built()` 返回错误）。实现腿**没有**拿 `rc=0` 冒充"有内容"（它 ⑥#6/⑰ 明确分了两口径）＝没有造假，但它自己那枚仪器是恒真。归口：产物由界面侧那枚 agent 带，票面 `:37` 逐字"这一格在'谁把 dist 填上'落定前勾不了"⇒ **今天勾不了**，与票 248 AC#9 同一枚归口。

### J1–J8 裁定落地形（逐条判语）

| 裁定 | 判语 | 凭据 |
|---|---|---|
| **J1** 宿主投现成 `ui-sta`（甲，第一枚用例就证"泵期间球仍出帧"） | **未落地；停手上报的姿势合规，凭证不合规。** 它⛔ 没造第二线程（产码无 `LockOSThread`、名册无新协程名，#16/#17 我复认）、⛔ 没自己动 `internal/ball`；但票面要求的那枚"证泵期间球仍派发"的用例**不存在**（自陈未提交），panic 判〔仅自述〕（#18），且"库不支持"这句只核到前提（嵌套泵非可重入）没核到结论。⚠ 它没说出"投递原语已在、只差一枚导出包装"这一层（§C#9）——这一句会影响编排者要不要摆 owner。 | #14–#18/#37 |
| **J2** 降级＋响亮，⛔ 任何路径不许让 `log.Fatalf` 可达 | **不成立。** 它只具名库内 2 枚 fatal（行号与我逐字读的差一枚：真身 `chromium.go:188`／`:284`），而产码入口 `webview2.NewWithOptions` **自带三枚 `log.Fatal`**（模块 `webview.go:115/120/125`）＋edge `:173`（＝运行库缺失那一支）／`:295`。⇒ "我们不走那两枚会 fatal 的 API"没有具名依据。AC#5 本体不在我射程，⛔ 我不判它。 | #11(§C)/§A#37 相邻读数 |
| **J3** 三口径分开报／P11 未触发不改契约 | **两半成立、一半不成立。** 阈值文件与 SLO 文档零改动 ✓、P11 未触发 ✓（11 发最大 937.093）；三口径**没分开报**（一枚混口径的 cold）。它把 `cold-embed 1808.9 > 1500` 留在 `A489` 的已知代价那一格，没拿来当本票读数 ✓。 | #19/#20 |
| **J4** 落点甲：宿主进 `cmd/wisp`，`internal/panel` 保持零平台分叉 | **成立。** 新产码/测试全在 `cmd/wisp/`，`internal/panel` 无平台分叉文件名、非 test 文件无 `//go:build`，`internal/panel` 侧唯一改动是那枚钉（97/32）。`panel_assets.go:13` 那句"expected to call the same panel.Assets API"它照做了（`serveEntry` 走 `panel.Assets`）。 | #10/#35/#28 |
| **J5** 依赖只走 `go mod` 解析，⛔ 手抄哈希（ban #5） | **成立（形状对、措辞两处要更正）。** `go.mod 2/0`＋`go.sum 6/0`；多出的 2 枚是 `x/sys` 老版本 `/go.mod` 行、由图解析产生且被 `tidy` 保留 ⇒ 佐证走了 `go get`；`go.mod` 里记成 `// indirect`（直接 import 该在直接块）。⛔ 我不判"它有没有从 spike 的 go.sum 手抄"（值本身与 spike 相同是**同版本必然相同**，不构成抄的证据；能构成的只有上面那 2 枚图解析产物）。 | #33/#29 |
| **J6** embed 路径按活接缝；"清检出能建"升硬判据 | **成立。** 票面那行 `assets/web/dist` 未被当判据（产码走 `panel.Assets` 的 `all:dist`）；AC#11 由我独立证有牙（#29）。⛔ 它没有把 `frontend/**` 提交进来糊这一格（票面写面清单里零枚前端件，#A 全清单）。 | #29/#30/#A#4 |
| **J7** L1 进仓当常驻用例（ubuntu core 有分母）／L2 只留本机证据并具名登记 | **半不成立。** L1 落了但落在 **windows cli 腿**、不在 ubuntu core（#32，且文件头自述与 CI 接线矛盾）；L2 落了但**恒真**（#24）＝比"只留本机证据"更糟；"CI 里没有 winlive 这一半"它具名登记了 ✓；词面 `netstat` 门没做 ✓（`grep -rn netstat --include=*.go cmd internal tools`＝0 命中）。 | #22–#25/#32 |
| **J8** 33 先、与 248／244 不合批 | **成立。** 三枚产码/测试提交各自带显式 pathspec（`697b4fae`／`b1e3d94a` 只碰 `cmd/wisp/*`＋`go.mod`/`go.sum`；`fafe0445` 只碰那枚钉）；构建链零改动（`scripts/build.ps1`／`.github/workflows`／`docs/BUILD.md` 全 0 行）；`bridge.go` 四枚入向常量未动（枚数＝4，我重跑）；票 248/228 的活一枚未做（我没发现越界写面）。 | #A#2/#35/#A#4 |

### 这程造成的 vs 不是这程造成的（一句话账）

- **这程造成**：产码形状三条（L2 常量错位是它写的尺、探针页覆盖入口字节是它的 `bringUp`、焦点尺不回还断言是它的用例）＋依赖记成 `// indirect`（tidy 不净）＋一处注释过期（`:101` 的 `LockOSThread`、`:64-66` 的两个不存在的常量名）。
- **不是这程造成**：`internal/panel` 起手那 4 枚红（#10 逐名归因，含 1 枚落在冻结件上）；`x/sys v0.47.0` 那两枚 tidy 会抹的行（#29）；go-webview2 库内的 fatal 站点与 nil 解引用（库自身形状）。
- **实现腿没有做的事（要记清，免得下一程重复怀疑）**：没勾框、没改 AC 原句（票面零行改动，#A#4）、没动三枚冻结件与阈值/golden/allowlist（#35）、没读前端内容造页、没为变绿放宽断言（M2 那发它没跑过，但我这一发证明：**它的用例即便在这形下也不会红**——这一点它表里写成 "L2 PASS（本 pid 零枚 LISTEN）" 属**未验证的自我确认**，不是它拧判据）。

---


## E. 收尾三把尺（终态现量）

| 尺 | 起手 | 终态（`11:12:07`，HEAD `345b058b`） |
|---|---|---|
| `git status --porcelain -- cmd internal` | **0 行**（`10:35:27`／HEAD `94bd133e`） | **0 行**＝等于起手名册 ✓（我全程只在工作树里做过两枚突变，逐枚 `git cat-file blob HEAD:<path> > <path>` 还原并 md5 比对：`e75b9b2e36772ca4cae6043a6dbe7384` 三处相同） |
| `git diff 94bd133e..HEAD -- cmd internal` | （锚＝起手 HEAD） | **空**（`--name-only` 枚数＝**0**）✓——HEAD 期间从 `94bd133e` 走到 `0d993d66`／`3398e6f1`／`34c28c5c`／`345b058b`（别枚只读腿与编排者的 docs 件），`cmd`／`internal` 一字未变 ⇒ 我的锚与读数仍然同一条线 |
| 三节自对抗非空＋占位符 `grep -c` 为 0 | — | §A **37 条**（逐条真读数＋时刻＋HEAD）／§B **17 条**／§C **14 条**（逐条甲／乙／不做＋现量）。尺：`grep -c "填写中｜待填｜TODO｜未判"`（把词换成全角竖线以免自我指涉）现量＝**1 命中，且那一枚命中就是本行自己写下的这条尺名**（`11:12:56` 现跑）⇒ **除去本行以外为 0** ✓。另两枚 `占位` 字样是 AC#12 的**领域用语**（`.gitkeep` 那枚占位文件，逐字在 §A#30／§C#6），不是空格。全文 `wc -l -c`＝**218 行／59,916 字节**（`11:12:56`，交付判据取这一发） |

### 其它收尾读数（同发取，`11:11`–`11:12`）

- 交付判据（不是回执）：`wc -l -c docs/evidence/s1/33-panel-host-c27-v1.md`。
- 票面卫生：`git diff --numstat 7a41db9b..HEAD -- .scratch/wisp/issues/33-panel-host-c27.md`＝实现腿**零行**；我自己全程没碰票面任何一枚 `- [ ]`／`- [x]`（勾框归编排者）。
- 我这两枚提交的 pathspec 都是显式单文件（`docs/evidence/s1/33-panel-host-c27-v1.md`）；`git add -A`／`.`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean` **零次**；⛔ 未 push。
- 我只读未动的禁区清单终态：`go.mod`／`go.sum`／`internal/observe/thresholds.go`／`docs/SLO.md`／`docs/BUILD.md`／`docs/PLAN.md`／`docs/specs/**`／三枚冻结件／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ 全部由 `git status`＋`git diff --numstat` 双尺现量为零改动（`frontend/dist` 里那 3 枚别人未跟踪产物我只数过条目名、⛔ 未读内容）。
- 我跑过的全部构建/测试都在 `PATH` 带 `third_party/sherpa-onnx:build` 的前提下取数（§A#9/#19/#23/#26），无一条是 `0xc0000135` 那种"用例根本没跑"的空读数。

### 交给编排者的三行结论（我不翻勾、不补判语）

1. **六格里只有 AC#11 我判成立**（且我自己仓外证了它的牙齿）；AC#1／AC#2／AC#12 是"读数真、判据形不齐／归口未落"；**AC#3 与 AC#4 我判不成立**——AC#3 的"无监听端口"那一维今天**没有任何会响的检**（L2 两枚常量错位，我已用真开端口的突变＋仓外探针把根因定位到行），AC#4 的焦点回还**一条断言都没有**。
2. **"宿主真装起来了"与"用户能开出一个窗口"之间差的不是细节，是一整跳**：`PanelManager` 生产调用者＝0 枚，常驻的 `OnPanelHotkey`/`OnTrayPanel` 仍只打点，所以 owner 那句"起码我要能看到主面板"在本票这六格里**没有被满足**，也⛔ 不许被这六格的任何绿读成满足。
3. **两枚停手上报里只有一枚是真的**：ui-sta 再入那枚我判〔仅自述，前提有源码依据〕，需要下一程留栈文或按 §C#9 那枚"只差导出包装"的小口重裁；`w32.Rect` 那枚**理由不成立**（导出的 `(*edge.Chromium).Resize()` 就是高层 API 自己用来定位控制器的那一步）⇒ 它把一枚产品形状决策报成了依赖边界，我按派单要求把这一句定性写死在这里，供你决定要不要摆 owner。


