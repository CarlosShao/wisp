# 票 255 · AC#4 格 · 非实现者验收（255-v2，2026-10-08）

被审对象：实现腿 `255-r1` 自报三笔 `d37e75bd` → `8b32060b` → `18655988`。
本腿现量 HEAD 起手值 `6320a16a`（⛔ 派单给的 `e05b8a2b` 已漂，在 HEAD~2，见 `logs/00-anchor-gitlog.txt`，rc=0）。
内容锚＝票面 `:20` 那句「让 `[panel] width` 真生效（这才是 owner 那句"改了要有反应"）」（本腿 `grep -n 'AC#4'` 现量，`logs/11-ticket255-cites.txt` rc=0）。

---

## Q1 落地性与产码射程

**结论一句：甲形落在该落的那一格、四枚禁区一枚没碰，所以"该翻 AC#4"这件事成立；但实现腿的产码名册**漏报一枚**（`cmd/wisp/panel_resident_windows.go` 3/3），翻勾时应按五枚文件而不是三枚记。**

凭哪把尺：
- `git show --stat` 逐枚（`logs/` 起手那发，rc=0）＋`git show --numstat --format="" 8b32060b`（`logs/25`，rc=0）：
  `config_readers_255.go 35/15`·`panel_geometry_255_test.go 8/8`·`panel_host_windows.go 93/4`·`panel_reshow_255r1_windows_test.go 563/0`·**`panel_resident_windows.go 3/3`**（自报没列这枚）。
- `git diff d37e75bd..18655988 -- cmd/wisp/panel_host_windows.go` 全文（`logs/02`，rc=0）。
- 依赖边：`sed -n '/^import (/,/^)/p' cmd/wisp/panel_host_windows.go` ⇒ import 块只有 `internal/panel` + `go-webview2` + `x/sys/windows`，**没有 `internal/config`**；文件里三处 `config` 全在注释里（`grep -n "internal/config|..."` rc=0，见批量输出）。装配根仍是唯一接缝：`cmd/wisp/panel_resident_windows.go:253` 用 `withGeometrySource(panelGeometrySource(dataDir))` 递闭包（`logs/07`，:253 现量命中）。
- 线程：`requestGeometryOnReshow` 体内**零 `go func`**，只走 `w.Dispatch(func(){...})`（`logs/02` 第 122-124 行）；`sh scripts/d22scan.sh` rc=0 clean（`logs/15`）⇒ D38／禁裸 goroutine 那条自动扫也背书。
- C27：diff 里没有任何 Hide/Destroy/重建；`m.w` 只在 bringUp(:414) 写、Destroy(:725) 置 nil，新代码 `m.mu.Lock()` 取下 `w` 副本再判 nil（`logs/26`，rc=0）⇒ 生命周期没动，也没有跨锁调用。
- `HintFixed`：调用点实参是 `webview2.HintNone`（`cmd/wisp/panel_host_windows.go:614`，本腿 `grep -n "SetSize" ` 现量 rc=0）。
- 锁形状补一句：`m.geometry` 只在构造期 `withGeometrySource`(:200) 写，之后再无赋值 ⇒ 新代码不持锁读它与 `windowOptions` 既有形状一致，**不是新增竞态**。

**读错了最坏会放行什么**：把"漏计一枚产码文件"当无害 ⇒ 后续程按三枚名册核射程，`residentPanelGeometryNote` 那句文案改动（含 `state=per-create-and-per-reshow` 这种对外可见的口径变更）就没人审过。

⛔ 本腿没翻票面（翻勾是编排者的活）；本腿的意见＝**该翻**，前提是把名册改成五枚。

---

## Q2 这枚仪器有没有牙（主活）

**结论一句：四问能力形整体有牙——本腿自造两形，M1 直接咬红（两枚用例、七条断言）；但 M2 咬不红，`TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne` 的"only one"是**单文件射程**，同包另一枚文件里新加一条 `HintFixed` 的 resize 路，它照绿 ⇒ 这一格具名恒真面。**

本腿自己造的错法（⛔ 不复用实现腿那三发；两形都是他们没列的）：

**V2-M1「取到配置那一环假通」**＝把 `requestGeometryOnReshow` 里的 `options := m.windowOptions()` 摘掉，改发编译期常量 `panelWidthPx/panelHeightPx`（overlay 件在仓外 `D:/tmp/wisp255v2/m1-panel_host_windows.go`）。
命令：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -overlay D:/tmp/wisp255v2/m1.json -run 'TestTicket255' -count=1 ./cmd/wisp/`
关键输出（`logs/20`）：`--- FAIL: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch`（五个子例红，两条断言各咬一遍：`asked for 420x260, want 517x331` ＋ `geometry source was consulted 0 times, want exactly 1`），`--- FAIL: TestTicket255r1ShowOnAnExistingWindowSendsTheResize`。**M1_rc=1** ⇒ 有牙。

**V2-M2「AST 那枚只许 1 枚调用点的尺是不是恒真面」＝是，且能具名说它怎么坏都照绿**：在同包一枚**兄弟文件**（overlay 虚拟路径 `cmd/wisp/zz_probe_255v2_windows.go`，盘上从没写过这个文件）给 `PanelManager` 加第二条 `w.Dispatch(func(){ w.SetSize(..., webview2.HintFixed) })`。
命令：`go test -overlay D:/tmp/wisp255v2/m2.json -run 'TestTicket255' -count=1 -v ./cmd/wisp/`（`logs/21`）
关键输出：`--- PASS: TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne (0.00s)`，整发 **M2_rc=0**，12 枚 `TestTicket255*` 全绿。
正控（⛔ 不许把"没响"当"没生效"，所以先证 overlay 真进了编译）：同一路径换成含未定义符号的件，`go vet -overlay D:/tmp/wisp255v2/m2bad.json ./cmd/wisp/` ⇒ `vet.exe: ... m2bad.go:5:33: undefined: thisSymbolDoesNotExist255v2`，**positive_control_rc=1**（`logs/22`）；同一命令换回 M2 那份是 rc=0 clean ⇒ 那枚兄弟文件确实在 package build 里，绿是真绿。
坏在哪：该尺自己的失败文案写着 "every extra one is a second place a hint could be picked wrong"，`parseHostFile255r1` 只读 `panel_host_windows.go` 一枚文件（`hostFile255r1` 常量＋`logs/16` 现量：255 一族所有 AST 尺都只读这一枚文件，⛔ 无任何尺全包扫 resize 调用点）。⇒ **名字里的"the host"是宿主级承诺，量的是单文件级事实。**
同族另一向（没跑，只具名）：那枚 `isPackageSelector255r1(args[2],"webview2","HintNone")` 认的是"包选择器字面名"，若把 import 改个别名而行为不变，它会**假红**；这属"词面钉"的脆，不属"恒真"，但同一枚尺两头都不精确。

其余三问本腿复认有牙（不重复实现腿的表）：`geoSink255r1.Dispatch` 同步跑闭包并用 `depth` 分流 `size`/`bare`（测试件 :104-116/:146-156），且有一枚自我证伪用例 `TestTicket255r1SinkTellsDispatchFromABareCall`（:319）钉住"两形分不开 ⇒ 上面那把变恒真"⇒ 摘掉 Dispatch 会被**运行时＋AST 双响**，这格不是恒真面。
"任何分支都跑"那一形：本腿推演未实跑——`countCall255r1(createBranch,...)!=0` 只封 `!created` 分支体内，**在 if/else 之后再无条件补一发**结构尺看不见；但 `TestTicket255r1ShowOnAnExistingWindowSendsTheResize` 数的是"恰好 1 次发送＋源恰好读 1 次"，会红 ⇒ 判定＝**有牙，但牙在 (d) 不在 (c)**。

**读错了最坏会放行什么**：真 bug 形＝以后任何人在 `cmd/wisp` 别处（比如常驻腿或重建路径）加第二条尺寸路并选 `HintFixed`，用户拖窗边和最大化按钮就没了，而 255 的尺**一声不响**，交接件还会写着"AC#4 的 hint 钉住了"。

---

## Q3 宣称有没有超出做到的

**结论一句：没有超出——"等于配置宽度"这句话它明写不宣称；`客户区/外框` 两形是本腿独立读库代码复认成立的；`〔仅本机可量〕` 用在该用的位置上（边框差没被写成"量过"）。真正没被任何尺覆盖的是"真窗上这一发到底动了没有"，而实现件 :4 自己就具名放弃了这个宣称。**

凭哪把尺：
- 库代码现量（`logs/23`，rc=0）：`webview.go:320-332` 把 `opts.Width/Height` **原样**喂给 `CreateWindowExW` 的 `nWidth/nHeight`（style `0xCF0000`），建窗那份＝**外框**；`AdjustWindowRect` 在整个 `webview.go` 只出现 **1 次**，就在 `SetSize` 的 `else`（HintNone/HintMin/HintMax 之外那一支，`logs/03`，:421-431）：先 `AdjustWindowRect(0,0,w,h)` 再 `SetWindowPos(...SWPNoMove|SWPFrameChanged)` 再 `browser.Resize()` ⇒ `SetSize` 那份＝**客户区**，两形确实不是同一个宽度，与 10-02 裁定（外框≠客户区）同向。
- `SWPNoMove` 在实参里 ⇒ 这一发**不改位置**，只改尺寸（`logs/03` :428-430）。
- "不许宣称等于配置宽度"：`requestGeometryOnReshow` 头注第 1 条逐字写着 `"the on-screen width now equals [panel] width" is not a claim this file makes - only "the pair the host resolves right now was sent" is`（`logs/02` :73-75）；`config_readers_255.go` 的 `[panel]` 行也是 "the same width is not the same on-screen pixels"（`logs/05` :76）。本腿对全文件 `grep -n '真机|实测|立即生效|等于配置'`（`logs/25`，rc=0）：`impl.md:4` 逐字 **⛔ 本件不宣称量过"真浏览器里的宽度"**。⇒ 名册里**零枚**把边框差写成"真机量过"。
- 窗口射程现量（`logs/24`）：`grep -rln requestGeometryOnReshow cmd/wisp/` ⇒ 只 3 枚文件命中，**`cmd/wisp/panel_geometry_255_winlive_test.go` 不在里面**（那枚 winlive 用 `GetWindowRect`（外框）量的只有建窗路 `TestTicket255RealWindowWidthFollowsTheConfig`，且默认 build tag 不编译）。⇒ "发给真窗之后真窗变了"这一格今天**无人量过**，实现件也没宣称它被量过。

**读错了最坏会放行什么**：把这格当"已验收真窗"⇒ 下一程按"AC#4 完成"去写 `DECISIONS`／面板文案，撞上第一次 re-show 就把窗口外框宽度挪了（边框差）、用户以为配置没生效，而那正是 owner 立票要治的病。

---

## Q4 名册与锚自洽

**结论一句：一票顶回实现腿失败、顶回派单成功——`304` 改动前只有 `:17` 一枚引用、改动后 0 枚；但它新写的两枚内容锚**有一枚现在真跑不命中**（`Width: 420` 单空格形 vs 原文 `Width:  420,` 双空格），且票面那族 `file:LINE [token]` 14 枚本腿逐枚命中。**

- 派单那条我独立复认（⛔ 未引编排者那把尺）：`grep -n '304' cmd/wisp/config_readers_255.go` ⇒ **rc=1 零命中**；`git show 18655988:...| grep -c 304` ⇒ `0`；`git show d37e75bd:...| grep -n 304` ⇒ **只有 `17:// (cmd/wisp/panel_host_windows.go:304-305 at HEAD 67ab595d ...`** 一枚（`logs/06`，rc=0）。⇒ **实现腿对，派单"两处（:141 与 :17）"错**；另 `git diff 18655988..HEAD -- cmd/wisp/config_readers_255.go` 空（rc=0）⇒ 那枚文件自交件后没被后续腿挪过，号不会漂。
- 两枚新内容锚命中性（`logs/08`，逐条真跑）：
  - `git show 67ab595d:cmd/wisp/panel_host_windows.go | grep -n 'Height: 260'` ⇒ `305: Height: 260,` **rc=0 命中** ✓
  - 同件 `grep -n 'Width: 420'`（照 `config_readers_255.go:18` 与 `:146` 写的字面）⇒ **rc=1 零命中** ✗；`grep -n 'Width:  420'`（双空格）⇒ `304:` rc=0。⇒ **锚的一半是坏的**：文件教人去 grep 的字面在原锚点找不到东西。
  - 连带后果（本腿具名，属"认字符串不认行为"同族）：`panel_geometry_255_test.go` 那枚防漂针仍扫 `"Width:  420,"`（双空格＋逗号），而文案里现在是单空格形 ⇒ **针绿是因为它和目标串差了个空格**，不是因为它证明了那件东西不存在。
  - `NewPanelManager took no config parameter` 这半句：`git show 67ab595d:...| grep -n 'func NewPanelManager'` ⇒ `:175 (disp, assets, dataPath)` **rc=0 成立**（`logs/06`）。
- 票 255 那族 `file:LINE [token]` 名册：本腿自造尺（`grep -oE '[A-Za-z0-9_./-]+\.go:[0-9]+ \[[^]]+\]'` 抽全枚 ⇒ 逐枚 `sed -n "Lp"` 取该行再 `grep -qF token`）⇒ **14/14 全 HIT**（`logs/09`＋`logs/09b`，rc=0）。含 `panel_host_windows.go:262 [Width:  uint(width)]`、`panel_resident_windows.go:207 [cfg.Panel.Width]`、`internal/ball/ball_windows.go:64 [SizePx]`、`internal/observe/logging.go:50 [RedactPaths]`（后两枚第一轮被我按 cmd/wisp 前缀解析成 MISS，改名册后重跑＝HIT，`logs/09b`）。行号没被这次改动挪动：diff 第一枚 hunk `@@ -229,10 +229,10 @@` 净 0 行，:+392 之后的 +6 行不影响它（`:392` 现量＝`WindowOptions: m.windowOptions(),`，"reached from the create at :392" 仍成立，`logs/07` rc=0）。
- "今天没有 resize 路"那句：**是被它自己这次改动变假的，而它在同一枚 commit 里把那句改了**——`config_readers_255.go` 旧 bullet 逐字 "this host contains no MoveWindow / SetWindowPos / SetBounds call, so an already-created window keeps its geometry until it is destroyed and rebuilt"（`logs/05` :53-57），而甲形经库的 `SetSize` 正好就是 `SetWindowPos`（`logs/03` :428 现量）⇒ 该句若留着即刻为假；`panel_geometry_255_test.go:355` 的 "(there is no resize path today)" 同批改成 "the resize 票 255-r1 added runs on a show request"（`logs/05` :88-91）。⇒ **改漂是真改、方向是对的，且没留给下一程当假账。**

**读错了最坏会放行什么**：单空格那枚锚被人照抄去 grep ⇒ 得"零命中"⇒ 结论翻转成"67ab595d 根本没有硬编码宽度，票 255 整票前提不成立"，这是一条能把已定案翻坏的假凭据。

---

## Q5 门禁

**结论一句：`go vet ./cmd/wisp/` rc=0、`sh scripts/d22scan.sh` rc=0 clean、`gofmt -l` 对五枚被改文件零输出 rc=0；行尾符尺先跑过，派单点名的五枚 CR 幻影逐枚复认；"inbound_raw_leak_35r7_test 解析报错"那一句本腿复现不出来 ⇒ 判它认错件。整包跑到终态那一发的数见下表（⛔ 不是 `-run` 单跑）。**

| 门禁 | 命令原文 | rc | 证据 |
|---|---|---|---|
| vet | `go vet ./cmd/wisp/` | **0**（零输出） | `logs/15` |
| d22scan | `sh scripts/d22scan.sh` | **0** `d22scan: clean - no D22 ban violations`，examined 266 production Go files | `logs/15` |
| gofmt（目标件） | `gofmt -l cmd/wisp/panel_host_windows.go cmd/wisp/config_readers_255.go cmd/wisp/panel_geometry_255_test.go cmd/wisp/panel_reshow_255r1_windows_test.go cmd/wisp/panel_resident_windows.go` | **0**（零输出） | `logs/15` |
| 行尾符尺（前置） | `tr -cd '\r' < F \| wc -c` vs `git show HEAD:F \| tr -cd '\r' \| wc -c` | 见下 | `logs/15` |
| 整包终态（跑到终，⛔ 非 `-run` 单跑） | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` | **1**（939 行；`FAIL cmd/wisp 577.352s`·`FAIL internal/ball 1.200s`·`FAIL internal/panel 9.204s`·`FAIL internal/risk 8.396s`，其余包 ok；`--- FAIL` 共 **11 枚**；⛔ **零 `panic:`／零 `0xc0000135`** ⇒ DLL 那两枚 PATH 前缀确实生效） | `logs/10`＋`logs/30`，rc=0 |

CR 现量（工作树 CR／HEAD blob CR，`logs/15`）：
`pending_read.go 131/0`·`internal/agent/tools.go 225/0`·`internal/risk/provenance.go 1115/0`·`internal/tools/bridge.go 1280/0`·`cmd/wisp/models.go 334/0` ⇒ **派单点名的五枚逐枚对上，确属本机 autocrlf 幻影，不是红**；本次被改的五枚文件工作树与 HEAD **都是 0 CR** ⇒ 交件本身没带幻影。
`internal/panel/inbound_raw_leak_35r7_test.go`：`gofmt -l` **零输出 rc=0**，CR **0/0** ⇒ **实现腿那句"解析报错"本腿复现不出来**。它那发是 `gofmt -l` 全仓（`probes/255/r1/logs/gofmt-repo`，10 行）——**判它把全仓输出里的别枚件认到这枚头上**（同一批它还交过 `logs/gofmt-preexisting-check` 5 行）。⚠ 这一条属"别的腿（35-v7）的件"，本腿⛔ 不修、只具名。

**读错了最坏会放行什么**：拿不带那两枚 DLL 的整包当终态 ⇒ `0xc0000135` 静默、**零 `--- FAIL`**，会把"根本没跑"读成"全绿"。

---

## Q6 不新增红

**结论一句：成立——本腿整包 11 枚 `--- FAIL` 里，`cmd/wisp` 那 5 枚**逐枚 ⊂** `111-c2` 现量的窗口依赖名册 15 枚；对照跑（把甲形那一发调用摘掉）同族仍红、名与耗时基本同形；另 6 枚红全在 `cmd/wisp` 之外，逐枚有具名外因（工作树里 `design/**` 被别的腿删掉／前端契约漂／机器计时预算），**与 255-r1 零重合**，因为它三笔改动一枚 `internal/**` 都没碰。⚠ 具名一处对不上：`TestAC4FocusReturnToPriorWindowGap33r5` 在带改动那发红、在对照跑**绿** ⇒ 属"在册红名册不是稳定集合"那一族（票 236 第 5 格），不是新增红，但本腿不能把它百分百排除在 timing 扰动之外。**

### 我这把尺（不是抄实现腿的表）
- 整包终态：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` ⇒ **rc=1**，`--- FAIL` 11 枚，`logs/10`（939 行）／摘要 `logs/30`（rc=0）。
- 名册现抽（⛔ 不引派单转述）：`grep -oE 'Test[A-Za-z0-9_]+' .scratch/wisp/probes/111/c2/logs/03-roster-rerun-vs-33v4.md` 的表体逐枚＝**15 枚**（12 枚 `windows` ＋ 3 枚 `windows && winlive`；`TestPanelHostLatencyPercentilesAC2` 那行逐字写"**不在名册**"）⇒ `logs/14`＋`logs/24` 区间，rc=0。
  3 枚 `winlive`（`TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor`／`TestTicket255RealWindowWidthFollowsTheConfig`／`TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`）默认 build tag **不编译** ⇒ 整包跑里不可能红。

### `cmd/wisp` 5 枚 vs 名册 15 枚（作差表，`logs/32` rc=0）
| 我这发红的 | 在 15 枚里？ | 红因（原文行号） | 与几何有无关系 |
|---|---|---|---|
| `TestPanelHostRealWindowHopAndLifecycle` 5.08s | **在** | `panel_host_windows_test.go:660/:662` `cold bring-up ... -1.000 ms`／`the message channel did not come up` | 无（`-1` 哨兵形，111-c2 同形在册） |
| `TestAC4FocusReturnToPriorWindowGap33r5` 5.15s | **在** | `:948` `the panel did not take the foreground on Show` | 无（前台竞争；⚠ 对照跑绿，见下） |
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 20.01s | **在** | `panel_resident_windows_test.go:325` `no report "ac13-probe" ... within 15s` | 无 |
| `TestAC14AwaitedBindingReplyReachesThePage` 20.01s | **在** | `:821` `no report "ac14r-0" ... within 15s` | 无 |
| `TestAC14GoSideEvalPushReachesThePage` 20.01s | **在** | `:866` `no report "ac14-push" ... within 15s` | 无 |
⇒ **5/5 ⊂ 15**，**名册那侧多出的 10 枚在我这发全绿或不编译**；**没有一枚红落在名册之外**（就 `cmd/wisp` 而言）。

### 对照跑（我自己复跑的，⛔ 未引用实现腿那件）
`PATH=... go test -overlay D:/tmp/wisp255v2/m3.json -count=1 ./cmd/wisp/`（m3＝当前 `panel_host_windows.go` 只摘掉 `Show` 里那一发 `m.requestGeometryOnReshow()`，其余逐字节相同；盘上文件未动）⇒ **rc=1**，`logs/31`／摘要 `logs/34`（rc=0）：
- 同族仍红且同名同耗时（±0.02s）：lifecycle 5.07s／AC13ColdStart 20.02s／AC14Awaited 20.02s／AC14GoSideEval 20.01s ⇒ **摘掉甲形不治好它们 ⇒ 这 4 枚与 255-r1 无因果**。
- `TestTicket255r1ShowOnAnExistingWindowSendsTheResize` **0.00s 变红** ⇒ 这一发同时是**对照的正控**：证明 m3 overlay 真进了编译，也证明那枚 (d) 尺**确实咬"Show 没走到"那一形**（⇒ Q2 的"Show 真走到吗"不是恒真面）。
- 两处对不上，具名：① `TestAC4FocusReturnToPriorWindowGap33r5` 对照跑**绿**（带改动那发红）；② `TestTicket223HandEditedFsLooseningCostsAnL2Card` 2.21s **只在对照跑红**（票 223 的 L2 卡片尺，与面板几何无关）。⇒ 两枚都属本仓已登记的"红名册不是稳定集合"那一族（票 236 第 5 格／`A451:9626` 第三次推翻），⛔ 我不据此判 255-r1 新增红。

### 另外 6 枚红的归属（逐枚具名，全在 255-r1 射程外）
`git show --stat` 三笔里 `internal/**` **零枚**（`logs/` 那发 numstat：5 枚全在 `cmd/wisp/`）⇒ 6 枚必然另有因：
- `internal/ball/TestC21TableColourRowsMatchTokensCSS` ＋ `internal/panel/TestC21DesignTokensFourWayAgree`：逐字 `read design/assets/tokens.css: open ...: The system cannot find the path specified` ⇒ 起手工作树就有 ` D design/assets/tokens.css`（`git status --porcelain -- design` 现量；`git ls-tree -r HEAD` 显示 HEAD 仍有 30 枚 `design/` 文件）＝**别的腿在飞的删除**，`logs/01` 已登记。⛔ 不提交、不还原。
- `internal/panel/TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`：`Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare` ⇒ 前端声明面（`frontend/**`，机主 10-07 只放开读那一格）在飞。
- `internal/panel/TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`：`a second style source appeared in code the bundle actually ships`（57 files）⇒ 同上前端面。
- `internal/risk/TestResolvePerCallBudget`：`2892185 ns/op = 2.892 ms/op (budget 1.000 ms, 402 samples)` ⇒ 机器耦合的 SLO 计时尺，**阈值一字节没动**。
⚠ 我这发的树 ≠ 实现腿那发的树：我在 HEAD `6320a16a`（其上还有 `7b95f381`/`e05b8a2b` 111-c2、`e922e0d1` 35-r7、`086b24ce`/`6320a16a` 182-c2），他们跑在 `18655988`。⇒ 这 6 枚里若有属于后续腿的，归口在他们那格，⛔ 不记到 255-r1 头上。

### 与实现腿自报对不上的两处（具名）
1. "8 分钟 `panic: timed out`"：**本腿两发都没复现**（`grep -nE 'panic:|timed out|0xc0000135' logs/10` ⇒ **零命中**，`logs/30` rc=1 是 grep 无命中的自然返回码）。我这两发都在 10 分钟默认预算内跑完（577s／505s）。
2. "5 枚同名同耗时全红"：同名 4 枚成立，第 5 枚（AC4 focus）在对照跑绿 ⇒ 见上。

**读错了最坏会放行什么**：把这 11 枚当"255 的账"退回实现腿 ⇒ 他们去动窗口用例或动 SLO 阈值（两边都是禁区）；反过来若只报"cmd/wisp 5 枚 ⊂ 15 ⇒ 全绿"，会漏掉 `design/**` 被删会让 `internal/ball`/`internal/panel` 的两枚 C21 尺在 CI 上也红这一枚外向影响。

---

## 我没量到的／不成立的派单前提

**不成立的派单前提（具名顶回，共 5 条）**
1. **HEAD 号漂**：派单写 `e05b8a2b`，现量 `6320a16a`（其下有 `e922e0d1`、`e05b8a2b`）⇒ 起手已按第 108 条自己复跑（`logs/00-anchor-gitlog.txt`，rc=0）。
2. **"实现腿自报三件改动"名册不全**：`8b32060b` 实际动 **5 枚文件**，含派单与自报三行里都没点名的 `cmd/wisp/panel_resident_windows.go`（3/3，改的是对外可见的 `PANEL-GEOMETRY` note 文案＋`state=per-create-and-per-reshow` 口径）⇒ 它不是"文案"，是一枚**产码文件**。
3. **"派单说 `config_readers_255.go:141` 与 `:17` 两处引用 `:304`"**：本腿独立尺复认 ⇒ 改动前（`d37e75bd`）**只有 `:17` 一枚**命中 `304`，改动后**零枚**（`logs/06`，rc=0）。**实现腿对、派单错**；与编排者自己那把同向，但⛔ 未引用他那把。
4. **"删掉的那句『今天没有 resize 路』是不是被它自己这次改动变假的"**：这句问法暗含两个分支，现量落在"**当时为真、被本次甲形（库内 `SetWindowPos`）变假、并在同一枚 commit 内同步改掉**"（`logs/05`＋`logs/03`）⇒ 属该做的自洽，⛔ 记成缺陷。
5. **"`inbound_raw_leak_35r7_test` 解析报错"**：`gofmt -l` 零输出 rc=0、CR 0/0（`logs/15`）⇒ 判实现腿把全仓 `gofmt -l`（它那发 `probes/255/r1/logs/gofmt-repo`，10 行）输出里的别枚件认到这枚头上。这格归 35-v7，⛔ 本腿不修。

**我没量到的（具名，⛔ 不算已验）**
- **真窗上这一发到底动了没有**：`grep -rln requestGeometryOnReshow cmd/wisp/` 只 3 枚命中，winlive 件不在内（`logs/24`）⇒ "让 width 真生效"在**屏幕像素**这层今天无人量过；`impl.md:4` 自己具名放弃 ⇒ 是**在册缺口**不是谎报。缺口形状＝要不要补一枚 `winlive` 的 re-show 用例，归编排者裁。
- **边框差的具体像素**：`〔仅本机可量〕` 用得对，本腿也没开真窗。
- **`m.geometry` 被构造后重新赋值**这条路：今天只 `:200` 一处写 ⇒ 无新增竞态；`d22scan` 无锁项、全仓零 `-race`（票 33 A704 在册），本腿⛔ 未上机证明。
- **AST hint 尺的"假红方向"**（改 import 别名而行为不变会红）：只推演，未实跑。
- **`-tags winlive` 那一档**：整包跑不含（111-c2 同标"不编译"）⇒ 无人复测。
- **票面翻勾**：⛔ 本腿没动 `.scratch/wisp/issues/**`（只读），只给"该翻"的意见＋"名册须按五枚记"这个前提。
- **"任何分支都跑"那一形只推演未实跑**（结构尺 (c) 封不住、运行尺 (d) 会响；见 Q2）⇒ 若要把它当已验，得补一发 overlay 实跑。

### 收尾自检（`logs/35`，rc=0）
- 三枚被 overlay 的产码文件盘上未变：`git hash-object cmd/wisp/panel_host_windows.go` == `git rev-parse HEAD:...`（起手已量 `098560a3...` 相等，rc=0）；`config_readers_255.go`/`panel_resident_windows.go` 同法复量。
- 收尾 `git status --porcelain` 行数与起手 `logs/01`（771 行）作差，只允许本腿那两笔 probes/255/v2 commit 造成差；`design/**` 删除等**在飞改动逐枚仍在册、⛔ 未提交未还原**。
- ⛔ 零 push。

