# `33-r1` — 面板宿主真起来：C27 PanelManager 单例／WebView2 窗口／embed.FS 离线供给

- 程：`33-r1`（写腿，后端主链）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`（票 33）
- 派单：编排者收件正文（本腿射程＝票面 **AC#1..AC#4 + AC#11 + AC#12**；AC#5/6/7/8/9/10 不在本腿）
- 裁定出处：台账 `docs/reports/pending-and-issues.md` **`A486`**（编排者裁定 J1-J8）＋**`A489`**（撞钉预检裁甲形反转）
- 起手时刻／锚：`2026-10-01 09:3x–09:4x +0800`／**实际起手 HEAD `7a41db9b`**（派单写的 `f1e7a3e2` 是它的父提交，一枚 `A489` 撞钉预检台账提交夹在中间＝起手锚漂了一枚，`git merge-base --is-ancestor f1e7a3e2 HEAD` 现量 YES）｜分支 `dev`
- 写作顺序（硬要求）：**本件先落 ⑥⑦⑧ 三节自对抗内容并 commit**，再回填 ①–⑤ 结论表。⑥⑦⑧ 不留任何占位词。

---

## ⑥ 我跑了哪些尺、每条真实读数

> 时刻一律 `+0800`；凡"只有本机有"的读数都带 HEAD。以下每条都是我自己现跑的，没有抄编排者预检。

1. **环境**（`go env` 现量）：GOPROXY=`https://goproxy.cn,direct`；GOSUMDB=`sum.golang.org`；GOMODCACHE=`D:\work\base\gopath\pkg\mod`；GOFLAGS 空。仓内无 `vendor/`。
2. **起手 `go.mod` 的 webview 命中＝0**（`grep -c webview go.mod`，`09:3x` 现跑，退出码 1＝零命中）。⇒ 证实 `l2_grant_boundary_test.go:71`「`grep -c webview go.mod` 是 0」这句在起手时是真的、会被本腿打假。
3. **`go list -m github.com/jchv/go-webview2@latest`＝`v0.0.0-20260205173254-56598839c808`**（现跑，与 `scripts/spike/go.mod:9` 逐字同枚）。该模块无 semver tag ⇒ `go list -m -versions` 为空是正常的，不代表拉不到。
4. **加依赖走的是主模块 `go get`**（`go get github.com/jchv/go-webview2@latest`，带 sherpa PATH，`go get-exit=0`）：`go: added github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808` ＋ `go: added github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1`。`go get` 后 `grep -c webview go.mod`＝**1**（require 行），`grep -ci webview go.sum`＝**2**（h1 + /go.mod 两行）。⛔ **没有手抄任何哈希**、没有从 `scripts/spike/go.sum` 或镜像目录取任何一行；两行 `go.sum` 是 Go 自己经 sumdb 写入的 ⇒ ban #5 未碰。`git diff --numstat -- go.mod go.sum`＝`go.mod 2/0`、`go.sum 6/0`（含 winloader 四条）。
5. **干净能建基线**（`GOFLAGS= go build ./...`，加完依赖、宿主产码之前，`09:4x`）：`build-exit=0`。
6. **AC#12 两口径现量**（`git ls-files` vs `find`）：
   - 已跟踪＝**只有 `frontend/dist/.gitkeep`**（`git ls-files frontend/dist` 恰一行）。
   - 工作树磁盘＝`.gitkeep` + `assets/index-BRKj5OIJ.css` + `assets/index-BVKlegVD.js` + `index.html`（`find frontend/dist -type f` 四行）⇒ 那三枚产物是**别人未跟踪**的前端构建物，不在版本库里。
   - ⇒ `frontend/embed.go:19` 的 `//go:embed all:dist`：在**干净检出**里只匹配到一枚 `.gitkeep`（`panel.Assets.Built()`＝false）；在本工作树里 `go build` 出来的二进制会捎上那三枚未跟踪产物。**"能构建"（AC#11）≠"真有一包入库内容"（AC#12）**，两回事。
7. **本机能否真起窗口**（撞钉预检我复跑）：`powershell` 现量 `SESSIONNAME=Console`、`explorer` 进程数 1；WebView2 运行库存在（`C:/Program Files (x86)/Microsoft/EdgeWebView/Application/` 下有 `153.0.4234.48`、`154.0.4258.37`）；跑现成 spike：`scripts/spike/bin/webview2-latency.exe -mode coldonly -datapath <仓外TMP>` ⇒ **`{"mode":"coldonly","coldMs":862.431}`、`exit=0`**（尾部一行 `Failed to unregister class Chrome_WidgetWin_0 Error=1412` 是非致命拆卸噪声）。⇒ **本机能真建 WebView2 窗口**，本腿的窗口类读数属本机真机量。
8. **那枚会打红我的钉的现状**（`internal/panel/composer_dispatch_test.go`）：`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops`（`:424`）当前语义＝"这棵树还没接宿主"，`hostChannelCapabilityHits` 扫产码标识符族 `CoreWebView2/WebView2/WebMessage`（`:504`）＋`go.mod`/`go.sum` 里 `webview`/`msedge`（`:512`），`len(hits)!=0` 即红（`:428-431`）；自带正控：假宿主源件（`:438` 必红）、假 webview 依赖（`:458`，写的是 `github.com/jchv/go-webview2 v0.0.0-20220111092826-5ea30c1aaa37`，必红）、CLI 接缝那形（`:473` 必 0 命中）。carrier 全用 `t.TempDir()`，`go test -overlay` 一枚未用。
9. **不许起第二通道**（`internal/panel/bridge_test.go:157-159`）：前端源件里出现 `window.postMessage(`/`new WebSocket`/`fetch(` ⇒ 红（"ticket 92 forbids it"）；`bridge.postMessage` 调用点必须＝2。⇒ 本腿入向只能走 WebView2 消息通道，⛔ localhost／websocket／任何监听端口。
10. **入向白名单四枚未动**（`internal/panel/bridge.go:42-45`）：`MethodModeRequest/MethodWorkspaceRequest/MethodAttachmentAdd/MethodMessageSend`。本腿零枚新入向方法名（配置读／写属票 248）。
11. **pkg/edge 里那两枚 `log.Fatalf`**（J2 我要避的）：`chromium.go:189`（controller 创建失败）与 `:284`（`WebResourceRequested` 里 `args.GetRequest()` 出错即 `log.Fatal`）。⇒ `:284` 只在有 filter 命中请求时才被回调；`:189` 只在 controller 真失败时；本腿的降级路径不主动引它们可达（见 ⑧）。
12. **会被本腿打假的两处过期指认**（⛔ 一字不改，只登记）：`internal/panel/l2_grant_boundary_test.go:70-71`（冻结件）「`grep -c webview go.mod` 是 0」；`cmd/wisp/approval_always.go:165`「this binary links no WebView2 host」。前者冻结、解冻范围是另一枚具名 `A##`（只覆盖 `:2051`/`:2139` 两枚锚、属票 248）；后者归票 33 自己改口。

## ⑦ 我可能写错的条目（对抗我自己）

1. **独立线程 vs ui-sta**：J1 裁的是"宿主投现成的 `ui-sta`、⛔ 不新增名册外协程名"。我的产码里给宿主一个"在调用方所在线程上 `Embed`"的形状——若**产码**里出现自建的 `runtime.LockOSThread` 常驻面板线程＝我违反了裁定＝跑歪模式。测量用的**测试**临时线程不能读成出货架构（详见 ⑧-J1）。
2. **embed 口径混用**：`find frontend/dist`（工作树磁盘 4 枚）会把"别人未提交的构建物"误算成"入库内容"。AC#12 判的是**已跟踪／干净检出**口径＝1 枚 `.gitkeep`。我若拿"本机能 build 出带页面的 exe"当 AC#12 的证据＝把 ⑥#6 的两种口径混为一谈＝缺陷。本票⛔ 不许拿"构建 rc=0"当"有页面"的证据，也不许我自己造一包产物填 `frontend/dist`（两层禁令）。
3. **回环测试用的是我自己写的极简 HTML**（`NavigateToString` 里一句 `window.chrome.webview.postMessage(...)`），不是真前端页面。它证的是"WebView2 消息通道这一跳 H3/H10 在真宿主上跑通"，**不证**"用户能在界面上点设置"。我若把它写成"设置可用／用户看得见卡片"＝越界宣称。票 33 收件第 3 条明令：窗口开出来≠设置可用（那条路由属票 248）。
4. **焦点回还的"编辑器"**：`GetForegroundWindow` 在我这台上取到的是测试进程的终端／控制台前台窗，不是真实文本编辑器。读数只对"记录前台窗 → 显示面板 → 隐藏 → 前台窗回到原窗"这一**机制**负责，不代表在真编辑器里的手感。带本机时刻，不当通用结论。
5. **反转钉的时序**：我承诺在**单独最后一枚提交**里反转 `TestSliceA…` 并改名。在那枚提交落地之前，整棵树的 `internal/panel` 是红的（宿主产码命中能力尺）。若我中途交件而没打出这枚反转提交，就是把仓库留在红态。⇒ 反转＋改名＋新增反向正控必须成对落在一枚提交，且改名后要具名登记"旧名只存在于归档日志"。
6. **go-webview2 的 GC 隐患**：`NewChromium` 注释自陈 handler 经 `uintptr(unsafe.Pointer(...))` 交给 native 且依赖"Go 不搬堆对象"。这是上游既成事实、非本腿能修，但它是我这族真机读数可靠性的一个已知前提，写出来免得被当成"我验证过内存安全"。
7. **常量尺噪声**：`bridge.go:35` 注释含 `panel.*`，`grep -c 'panel\.'` 起手＝5 不是 4（前一程已具名"尺与描述不符"）。我不用那把字面尺判白名单枚数，用 `bridge.go:42-45` 直接点名。

## ⑧ 判不动的地方（逐条甲／乙／不做＋现量）

- **J1（宿主投 ui-sta + 第一枚用例证"泵期间球仍出帧/消息仍派发"）——本腿撞到一枚写面外依赖**：
  - 现量：`internal/ball` 对 `cmd/wisp` 公开的、能"把任意闭包投到 ui-sta 上执行"的钩子＝**零枚**（`grep -rn "^func (b \*Ball) [A-Z]" internal/ball/ball_windows.go` 全列出的是 SetState/SetBadge/TimersAlive/Close/DebugHWND 等，没有 `PostTask`/`RunOnSTA` 一类通用投递口；唯一的在-STA-执行入口是 `Options.Events` 那几个回调，它们只在真实 `WM_HOTKEY` 到来时被 ball 的 wndproc inline 调用）。
  - ⇒ **出货路径**（原判甲·可落，现量后**改判受阻**）：把宿主接在现成的 `OnPanelHotkey`／`OnTrayPanel` 回调里（`resident_ball_windows.go:127-128` 今天只 `recordBallGesture`）。回调确在 ui-sta 上，但**从它再入跑 go-webview2 的阻塞嵌套泵会 panic**（见 ② 硬伤 2）⇒ 这条出货接线本腿**未做**、不是"可落"，改列停手上报。⛔ 未改成独立线程绕过（那是 J1 明令禁止的形状）。
  - ⇒ **测试路径**（同受阻）：要在自动化用例里"证泵期间球仍派发"，得让 `Embed` 真跑在 ball 的 ui-sta 上；上面现量＝写面内没有把闭包投到 ui-sta 的公开口，且再入会 panic。我一度写了这枚再入用例（本机现量 nil-webview panic 栈已入 ②），因它会把 `cmd/wisp` 留红、且证的是"落不了地"这一件事，**未提交**该文件（是我本程新写的，非在册钉），改把结论落进 ②/⑧。⇒ **判乙·请裁**：落地 ui-sta 复用要么授权动 `internal/ball` 交泵权（别的模块写面／可能触 D38），要么换/加能异步创建的绑定（连 ② 硬伤 1 一起）。本腿交的是"宿主原语 + 独立线程上的真窗口回环/生命周期/延迟/焦点/无端口读数"，ui-sta 共驻那一半明确停在"库不支持、未造第二线程"。
  - ⛔ 不做：独立常驻面板线程（改名册＝契约面）。
- **AC#5（运行时缺失降级 no-crash）／AC#6（CSP）／AC#7（快照跨瞬间）／AC#8（泵死分支）／AC#9（入向已接线的界面意义）／AC#10（[panel] 宽高 scale 生产者）**：都不在本腿射程（派单限定 AC#1–4/11/12）。本腿只**不制造回归**、**不宣称**。J2 那两枚 `log.Fatalf`：本腿宿主在"环境/控制器创建失败"时走**具名事件 + 不 fatal**（`:189` 只在 edge 内部 controller 创建失败才 fatal，属库路径；本腿不新增可达 `log.Fatalf` 的产品分支，见交付）。
- **winlive / CI 分母**（J7）：CI `winlive` **零岗位**（第四次具名，`ci.yml` 零命中）。⇒ L1（`go/ast` 扫宿主包 import 禁 `net`/`net/http`）落 ubuntu core scope 有真分母；L2（`GetExtendedTcpTable` 且过滤 `dwState==LISTEN`）＋窗口类读数（冷/热/焦点/回环）只在有桌面的本机跑、每条带时刻＋HEAD。⛔ 不加 `testing.Short()`／环境变量豁免假装绿：窗口类用例靠 `//go:build windows` 平台分叉天然不进 ubuntu 编译（那是合法平台门，不是假装跳）。

---

（①–⑤ 结论表：宿主产码形状、四枚 AC 的逐格判据与读数、门禁终态、票面 AC 框自证——待宿主与用例落地后回填。⑥⑦⑧ 已按硬要求在表之前写满并 commit。）

---

## ① 宿主产码形状（落了什么、没落什么）

- 产码：`cmd/wisp/panel_host_windows.go`（`//go:build windows`，`PanelManager` 单例）。
  - C27 单窗：每会话至多一枚 WebView2 窗口，`Show`/`Hide`（hide-不-destroy）/`Destroy`（dispose 才真拆）；`HotShow` 走复用路径并另记 hot 延迟。
  - 入向那一跳（本票主链）：页面 `postMessage` 的一枚原生封套经**绑定函数** `wispDispatch`（= WebView2 消息通道，非第二通道）进 Go → 唯一那道 `panel.ComposerDispatch.Handle` → 绑定返回值即回灌页面的回执 ⇒ H2（控件真建）/H3（页面 raw 进 Go）/H10（回执回页面）落在一枚真宿主上。**零枚新入向方法名**（`bridge.go:42-45` 四枚常量一字未动，配置读/写归票 248）。
  - 离线供给：资源从 `panel.Assets`（`//go:embed all:dist`）经 `SetHtml`/`Navigate` 直喂控件，⛔ 无 localhost/无 websocket/无监听端口（`net`/`net/http` 零 import，见 ⑥#11 L1 尺）。
  - 焦点：`Show` 记录 `GetForegroundWindow` 前窗，`Hide` 回还。
- ⛔ 未落（且本腿明确不做，非遗漏）：
  - **AC#3 的"多文件页面经 `AddWebResourceRequestedFilter` 服务"这一半**——见下方"两枚库硬伤"。当前离线供给是 SetHtml 单文档，不是资源过滤器。
  - **J1=甲 的"宿主投现成 ui-sta"出货接线**——见下方"两枚库硬伤"。出货走 OnPanelHotkey 回调（在 ui-sta 上再入 bringUp）在本机 panic，停手上报，未造第二线程绕过。
  - AC#5（运行时缺失 no-crash）/AC#6（CSP）/AC#7/AC#8/AC#9/AC#10：不在本腿射程，未做未宣称。
- 提交序（三枚，各带显式 pathspec）：`697b4fae` 宿主+依赖+测试（本枚**故意**让 `internal/panel` 的宿主钉由绿转红）→ `b1e3d94a` 消 ban#1（专用线程 bring-up 从产码移进 `_test.go` 脚手架）+ 修 AC#11 抽取 → `fafe0445` **单独最后一枚**把 `TestSliceA…` 期望反转成"宿主必须已接上"、改名、加反向正控。

## ② 两枚 go-webview2 库硬伤（AC#3 过滤器半 + J1 ui-sta 复用半，都是停手上报）

1. **`pkg/edge` 的 `AddWebResourceRequestedFilter` 能调，但同层 `ICoreWebView2Controller.PutBounds` 的参数类型是 go-webview2 的 internal `w32.Rect`**（`pkg/edge/ICoreWebView2Controller.go:59`），主模块 import 不到该 internal 包、`.vtbl` 也不导出 ⇒ 外部模块给不了控制器尺寸 ⇒ 纯 edge 宿主渲染空白。要一枚**能渲染**的窗口只能走高层 `webview2.WebView` API，而它**不暴露**资源过滤器。⇒ 本腿离线供给退成 SetHtml；真正的"多文件 embed 经过滤器服务"需要一枚依赖/库决策（replace/fork 或换绑定），属人工批准，⛔ 不自造形状。
2. **`embed` 的窗口 bring-up 是阻塞式嵌套 `GetMessageW`**（高层 `webview.go:352 Run`／edge `chromium.go:97` 的嵌套泵），它必须是该线程**唯一**的泵。J1=甲 要"从球的 ui-sta 回调里建窗"＝在球已有的 `DispatchMessageW` 里**再入**跑这枚嵌套泵 ⇒ 本机现量 `edge.Chromium.Init` nil-webview panic（栈：`ballWndProc→staThread.start DispatchMessageW→…→PanelManager.bringUp→webview2.NewWithOptions→CreateWithOptions(webview.go:340)→Embed→Init(chromium.go:131)`）。⇒ 按票面"证不出来就停手上报、⛔ 不许改成独立线程绕过"处置：出货 ui-sta 接线**未做**、未造第二线程；测量用专用锁线程只活在 `_test.go` 脚手架（不进出货、不登名册）。
   - 若要落地 ui-sta 复用，需其一：(a) edge 的异步创建路径（受硬伤 1 的 `w32.Rect` 尺寸问题连带卡住）；(b) 改 `internal/ball` 让 ui-sta 交出泵权（别的模块写面 + 可能触 D38 契约）。均判**请裁**。

## ③ 四枚 AC 逐格判据 + 读数（本机真机，带时刻／HEAD `fafe0445`）

| AC | 判据（会响的尺） | 本机现量 |
|---|---|---|
| **AC#1 生命周期** | `TestPanelHostRealWindowHopAndLifecycle`：二次 show 复用同窗（`hide→re-show` 后 `msedgewebview2.exe` 子进程数＝hide 前）；dispose 后 `IsCreated=false`、HWND=0、子进程 ≤2s 回落基线 | 子进程数稳定通过；Destroy 后 ≤2s 回落通过（无新子进程） |
| **AC#2 冷/热延迟** | 同测：`LastColdMs ≤1500`、`LastHotMs ≤200`（D32 面板行；`thresholds.go` 零延迟毫秒字段、一字未动，读数只进本表） | **冷 717.9ms**（`10:05`，另一发 `10:04` 837.2ms）／**热 36.9ms**（另一发 26.6ms）。均 <1500/<200；P11（冷拉起 >2s 才重评 L2）未触发 |
| **AC#3 离线+无监听端口** | L1＝`TestPanelHostOpensNoListeningSocketL1`（AST 扫宿主产码不 import `net`/`net/http`，ubuntu 有分母）；L2＝本机 `GetExtendedTcpTable` 过滤 `dwState==LISTEN` 且 pid∈本进程 ⇒ 须 0 | L1 PASS；L2 PASS（本 pid 零枚 LISTEN）。⚠ 过滤器那半见 ②硬伤 1，SetHtml 单文档已离线无端口，多文件页服务未落 |
| **AC#4 焦点回还** | 同测：`Show` 后取前台＝面板 HWND；`Hide` 后尝试回还前窗（记录句柄并 `SetForegroundWindow`+`SetFocus`） | 现量：Show 后前台＝面板 hwnd（`0x3a07d6`）。回还机制在（前窗句柄已记录并恢复）；Windows 前台锁可能不让非前台进程强夺焦点，此处只钉机制不硬断真编辑器手感 |
| H2/H3/H10 | 同测：真封套 `{method:panel.mode.request,...,source:panel-composer}` → `dispatchRaw` → `recordingModeHandler.count()==1` | PASS：mode 处理器被走到 1 次；回执经绑定返回 |

⚠ 上表读数只在**有桌面的本机**成立（winlive 在 CI 里零岗位，第四次具名）；靠 `//go:build windows` 平台分叉天然不进 ubuntu 编译，⛔ 未加 `testing.Short()`/环境豁免假装绿。回环用的是我自己写的极简 HTML（`NavigateToString`/`SetHtml`），证的是 WebView2 通道这一跳在真宿主跑通，**不证**用户能在真界面点设置（票 33 收件第 3 条：窗口开出来≠设置可用，那条路由属票 248）。

## ④ 门禁终态（收尾必跑，`10:2x`，HEAD `fafe0445`）

- `GOFLAGS= go build ./...`：exit 0。
- `gofumpt -l cmd/ internal/panel/`：空（含 `tools/d22scan` 无关）。
- `bash scripts/d22scan.sh`：**exit 0**，"clean - no D22 ban violations"；live 分母 `bans#1-5 internal/=224 · cmd/=33`，`ban#8 internal/=476 · cmd/=78`（文件枚数口径）。（修 ban#1 前它曾 exit 1，红在 `cmd/wisp/panel_host_windows.go:433 bare-goroutine`＝本腿产码裸 `go func(`；已把该线程移进 `_test.go` 脚手架消掉，未放宽任何断言。）
- `PATH=…sherpa… go test ./cmd/wisp`：**ok 184.420s，0 枚红**（含 `TestPanelHostRealWindowHopAndLifecycle`、L1、AC#11 干净检出 `go build ./... ok in 15.58s`、AC#12 记录）。
- `PATH=…sherpa… go test ./internal/panel`：**FAIL，恰 4 枚**，逐名：
  1. `TestComposerContractTypesMatchFrontend`（在册常红，读 `frontend/**`）
  2. `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（在册常红，读 `design/**`）
  3. `TestC21DesignTokensFourWayAgree`（在册常红，读 `design/assets/tokens.css`，起手即在、归因＝工作树他人删除）
  4. `TestApprovalCardViewJSONKeysMatchFrontend`（**本腿新增可见但非本腿造成**：它拿 Go `ApprovalCardView` JSON 键对 `frontend/**` TS 接口比对，红因＝他人会话前端 TS 漂移；本腿未动 `approval*.go` 亦未动 `frontend/**`）
  - 本腿的 `TestPanelHostIsAttachedAndNamesTheWindowHops`（原 `TestSliceA…`）＝**PASS**。4 枚红全属 `frontend/**`/`design/**` 地界，⛔ 我未修、未当绿、未放宽。

## ⑤ 票面 AC 框自证 + 边界自证

- `grep -c '^- \[ \]' .scratch/wisp/issues/33-panel-host-c27.md` 起手现量；本腿**未勾任何 AC 框、未改任何 AC 原句**（勾归编排者）。
- `grep -cE '=\s*"panel\.' internal/panel/bridge.go`＝4（常量四枚，起手终态一致）；`wc -l internal/panel/composer_dispatch.go` 本腿未动。
- 三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）：`git diff HEAD~3 HEAD -- <三件>` 应零改动——本腿只碰 `composer_dispatch_test.go` 一枚（反转那枚），`l2_grant_boundary_test.go:70-71` 一句一字未改（仅在此具名登记为"本腿 go.mod 落地后该『grep -c webview go.mod 是 0』指认过期，编排者已裁、解冻属票 248 的另一枚 A##"）。
- `docs/PLAN.md`／`docs/specs/**`／`docs/BUILD.md`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`：本腿零改动。`frontend/**`／`design/**`：零写面（只被 L1/AC 测试**读到**，且 AC#11 的仓外副本读）。
- 明文凭据：零枚；`cmd/wisp/panel_host_windows.go` 无 `filepath.Clean/Abs` 做安全决策（尺寸全走 WebView2/`panel.Assets` 的字符串路径）。

## ⑥（终态补记）收尾尺逐条真实读数

（在起手 ⑥#1–#12 之外，终态补跑并原样入表：）

13. `git rev-parse --short HEAD` 终态＝`fafe0445`；三枚本腿提交＝`697b4fae`/`b1e3d94a`/`fafe0445`；`git status --porcelain -- cmd internal` 终态 **0 行**＝等于起手名册（起手亦 0 行）。
14. `GOFLAGS= go build ./...` 终态 exit 0。`gofumpt -l cmd/ internal/panel/` 空。
15. `bash scripts/d22scan.sh` 终态 exit 0（修 ban#1 前曾 exit 1，见 ④）。
16. `PATH=… go test ./cmd/wisp` `ok 184.420s`；`./internal/panel` FAIL 恰 4 枚（⑤/④ 逐名，全属他人地界）。
17. 依赖增量终态：`grep -c webview go.mod`＝**1**、`go.sum`＝**2**（`go get` 经 sumdb 写就，非手抄）；`grep -rlniE "webview2|WebMessage" --include=*.go internal cmd tools | grep -v _test` 较起手多本腿一枚产码 `cmd/wisp/panel_host_windows.go`（它 import `github.com/jchv/go-webview2`）。
18. AC#12 终态现量：`git ls-files frontend/dist`＝1 枚 `frontend/dist/.gitkeep`（`TestEmbeddedDistCleanCheckoutHasPlaceholderOnly_AC12` 记录 tracked entries=1；本机能 build 出带页面 exe 是因工作树有他人**未跟踪** `dist/assets/*`+`index.html`，与"入库有内容"是两回事）。
19. 反转钉自身现量：`TestPanelHostIsAttachedAndNamesTheWindowHops` PASS；正控 A(5命中)/B(2命中)/C(0命中) 仍响；反向正控 D（中和宿主标识符）报 `(production=false dependency=false)`＝反转后的"必须已接上"确能红（谓词一字未动）。

