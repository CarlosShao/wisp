# 299-a1 · 40-open-questions — 四枚从"已知"退回"待量"的问题

锚＝HEAD `6ef14788`。⛔ 本件不采纳派单 §0 里那四条历史读数；每一枚答案只引本腿现跑的尺（件 `10`/`20`/`30` 里的 file:line 与逐字片段）＋具名标出"这枚还是预测"。

## Q1 甲形自不自足？—— **不自足（判语），但"到底请不请求得到"这枚事实本腿量不到**

**为什么不能只回答"注册了就通了"**：名册里两枚 URL 全是相对形状（件 `10` §1：`src="./assets/index-B8yINMF1.js"`、`href="./assets/index-yy8KMgdf.css"`；`frontend/vite.config.ts:23` 逐字 `  base: "./",`——这一枚本腿自己复跑过）。而今天文档基那件事，仓里与依赖里各有一句逐字：

- `cmd/wisp/panel_host_windows.go:95-96`：`// … panelOrigin is the opaque page origin that` / `// SetHtml produces (about:blank): offline, no host name, no port.`
- `webview2@…/common.go:61-63`（顶层库对 `SetHtml` 的 doc）：`// SetHtml sets the webview content directly.` / `// The origin of the page is `about:blank`.`
- `webview.go:393-395` 逐字 `func (w *webview) SetHtml(html string) {` ⇒ `w.browser.NavigateToString(html)` ⇒ **宿主那一发 `SetHtml` 就是 `ICoreWebView2.NavigateToString`**（`pkg/edge/chromium.go:123`）。

`SetHtml` 那发在产码里的位置＝`panel_host_windows.go:486` `w.SetHtml(string(data))`（`serveEntry` 体内，与票面现量第 5 条同一枚行号，本腿复尺复核为**成立**）。

⇒ 结论链条：**若**文档基真是 `about:blank`（两枚原文都这么写），相对 URL 无基可解析 ⇒ 页面**根本不会发出 `./assets/*` 请求** ⇒ 只注册过滤器＝注册了一枚没人来敲的门。**这一支不成立的可能性极高，但它不是本腿能钉的**：URL 相对解析是浏览器语义，⛔ 不在仓里、⛔ 不在依赖源码里。**具名标为待量**，尺见下面 §Q1.3。

### Q1.1 "只注册不改供给"能不能让 `./assets/*` 真被请求到

按上面三枚原文：**不能**（前题＝`about:blank`）。件 `20` §1 也补了一枚独立的"即使请求到了也不够"的洞：库在 `chromium.go:281-287` 的蹦床里做的是 `if e.WebResourceRequestedCallback != nil {` —— **只调 `AddWebResourceRequestedFilter` 而不设 `Chromium.WebResourceRequestedCallback`（`chromium.go:42`）＝请求来了也没人答**，两半必须配对着动。

### Q1.2 候选形状逐枚（每枚带 file:line＋要动的文件枚数；⛔ 不凭 API 名字推，可达性一律按件 `20` §3 的"宿主握 `webview2.WebView`"重算）

| 候选 | 依赖侧出处 file:line | 要动哪几枚文件 | 可达性判语（现量） |
|---|---|---|---|
| **a. 自定义 scheme ＋ 过滤器**（`wisp://panel/…`，文档与资源同 origin） | 注册：`pkg/edge/corewebview2.go:383`、`pkg/edge/chromium.go:292`；导航：`pkg/edge/chromium.go:116 func (e *Chromium) Navigate(url string)`；造响应：`pkg/edge/corewebview2.go:169`；应答：`...RequestedEventArgs.go:27 PutResponse` | **2 枚产码起**：`cmd/wisp/panel_host_windows.go`（换掉 `:386 NewWithOptions` 的对象取得方式＋把供给从 `:486 SetHtml` 换成带 host 的 `Navigate`）＋新增 1 枚回调文件；测试面另加 ≥2 枚 | ⚠ **越界**：`Navigate` 只在 `*edge.Chromium` 上，而宿主握的 `webview2.WebView` 只有 `Navigate`/`SetHtml`（`common.go:56-63`）——顶层 `Navigate` **是**可用的（`webview.go:389-391` ⇒ `w.browser.Navigate`），所以"导航到自定义 scheme URL"这一半**现有对象就能做**；**不能做的只有"注册"那一半**。⇒ 甲形卡的是注册对象的**所有权**，⛔ 卡在导航 |
| **b. `SetVirtualHostNameToFolderMapping`** | `pkg/edge/ICoreWebView2_3.go:22`；导出的取得口：同文件 `:58 func (e *Chromium) GetICoreWebView2_3()`、`:46 func (i *ICoreWebView2) GetICoreWebView2_3()` | 1 枚（`panel_host_windows.go`）＋**一枚落盘面**（要么改构建链把 dist 复制出去，要么运行期解包） | 同样需要 `*edge.Chromium`（`:58` 挂在 `Chromium` 上）⇒ 与 a 同一枚门槛。**且它要 `folderPath`＝磁盘目录** ⇒ 造出票面 AC#1 禁的"第二枚资源真相源"，`//go:embed all:dist`（`frontend/embed.go:19`）就不再是唯一的一枚。**〔越界候选〕，本腿⛔ 采纳**（理由已写在件 `30` §4） |
| **c. `NavigateToString`** | `pkg/edge/chromium.go:123`（＝今天 `SetHtml` 走的同一条通道，`webview.go:394`） | 0 枚（已经是它） | ⛔ **不是候选**：它与现状同一枚 origin 语义，换了名不换形状 |
| **d. 给入口页添 `<base href>`** | — | ⛔ 撞 `frontend/index.html` 冻结面；且入口 meta 逐字带 `base-uri 'none'`（件 `10` §1 那枚文档里读到的）⇒ **注入 `<base>` 本身会被 CSP 拒** | ⛔ 不成立，两重 |

### Q1.3 这枚待量的尺（一条命令就能钉死，⛔ 本腿没跑，具名）

现读 `document.baseURI`／`location.href` 的通道**已经在仓里**：`WebView.Eval`（`common.go:70-73`）＋宿主的 `installPanelTransport`（`panel_host_windows.go:801`）＋测试侧的 `evalOnPanelThread`（`panel_resident_windows_test.go:322` 那一发）与 `probeJS` 自报形状（同文件 `:295-303`，逐字 `window.wispDispatch(`+`reportJSEnv("ac13-probe", "b"))`）。⇒ **把 `probeJS` 的 `b` 换成 `document.baseURI + '|' + location.href`，一发就知道基是不是 `about:blank`**；这一发要么归落地腿（它独占 go 编译面），要么归编排者的真窗那一格（AC#5）。⛔ 本票 AC#0 不许改码，所以本腿只报尺名与改法，不落任何代码。

## Q2 入口那枚 CSP `<meta>` 今天到底生不生效？—— **今天没有任何仪器能看见它；两形各一枚验法如下**

⛔ 不按"文档基是 about:blank"推理（派单明令）。分两问：

**(甲) 它今天有没有被解析？** 一枚不碰 `frontend/**` 的验法：**`Init` 注入一枚"探针哨兵"**——`WebView.Init`（`common.go:65-68`，doc 逐字 "It is guaranteed that code is executed before window.onload"）走的是 `AddScriptToExecuteOnDocumentCreated` 那一族宿主侧注入，**不经过文档内联许可**；先 `Init("window.__wispCspProbe=1")`，再 `Eval("String(window.__wispCspProbe)")`：拿得到 ⇒ 页面脚本上下文活着、能读；然后**在文档字符串里临时加一枚内联 `<script>window.__inline=1</script>`（改的是 Go 侧 `serveEntry` 喂进去的字节、⛔ 不改 `frontend/**`），`Eval("String(window.__inline)")`：拿不到 ⇒ `<meta>` 的 `script-src 'self'` 正在生效；拿得到 ⇒ 它没生效（乙形的那道 CSP 关卡就不存在）。

**(乙) 两形各怎么验（生效／不生效两种预期）**

| 形状 | 生效时的预期读数 | 不生效时的预期读数 |
|---|---|---|
| **甲（注册过滤器＋带 host 的文档）** | `./assets/*` 请求到达、`Resolve` 返回字节；但若 origin 与 `'self'` 不匹配 ⇒ 控制台出现 `Refined to load the script … Content Security Policy directive: "script-src 'self'"` 那一族拒绝，**页面仍空白**；`document.getElementById('root').children.length`＝**0** | 无拒绝消息、`root.children.length` **>0**（React 真挂载） |
| **乙（内联整包）** | 内联 `<script>` 被拒 ⇒ `root.children.length`＝0，且 §(甲) 那枚 `window.__inline` 读不到 | `window.__inline`＝`1`、`root.children.length`>0 |

两行的断言都必须**由页面自己报回**（票面 AC#2 ⓐ 的定式），⛔ 不能由 Go 侧数调用。

**今天仓里有没有仪器能看见它？——没有。** 尺＝`git grep -n -i 'Content-Security-Policy\|CSP' HEAD -- ':!.scratch'` ⇒ **37 枚命中**，分布：`docs/PLAN.md` 8、`docs/reports/pending-and-issues.md` 5、`docs/evidence/s1/**` 若干、`frontend/index.html` 2、`docs/specs/SPEC-08` 2 …；**`cmd/`／`internal/`／`tools/`／`scripts/` 产码侧命中＝0 枚**（尺＝同一发结果里 `grep -E 'HEAD:(cmd|internal|tools|scripts)/'` ⇒ 空）。⇒ "CSP 生效与否"这件事在今天的 CI/本机两边都是**不可见量**，本腿把它归为"要先造仪器、再谈判据"。

⚠ 顺手上的一枚仪器牙齿读数（与本问同族，具名报回）：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 的探针 id 来自**入口 HTML 静态字节**（`entryIDProbes`＝`panel_resident_windows_test.go:265-276`，`assets.Resolve(panel.EntryFile)` 之后扫 id），而今天入口文档里**只有 1 枚 id**＝`<div id="root"></div>`（件 `10` §1 全文读数）。⇒ 它断言的是"`#root` 在不在文档里"，而 `#root` 是**静态存在**的：⛔ 它对"js 到底跑没跑"**不敏感**，也就是它对今天这枚纯白缺陷**没有牙**。这不是它该背的格子（票 33 AC#13 管的是两发 `SetHtml` 的次序），但票 299 **AC#2 ⓐ 不能复用这一形状**——要断言的是 `root.children.length`／页面自报的资源到位，⛔ 再数一次 id 在场。

## Q3 那几枚"窗建出来了／回执到了"的用例到底进不进 CI？—— **逐名册；⛔ 本机红≠CI 红，本腿两枚都没跑**

尺①＝`git show HEAD:<file>` 首行 `//go:build`；尺②＝`git show HEAD:.github/workflows/ci.yml` 里对应那一步。

| 用例 | 出处 file:line | tag | 进不进 CI |
|---|---|---|---|
| `TestPanelHostRealWindowHopAndLifecycle` | `cmd/wisp/panel_host_windows_test.go:614` | **`//go:build windows`**（首行逐字，本腿复核） | **进**（见下面"矛盾①"）——它开的是**真 WebView2 窗** |
| `TestAC3ListeningSocketRulerSeesItsOwnListener` | 同文件 `:1036` | 同上 | 进 |
| `TestAC14AwaitedBindingReplyReachesThePage` | `cmd/wisp/panel_resident_windows_test.go:812` | `//go:build windows` | 进 |
| `TestAC14GoSideEvalPushReachesThePage` | 同文件 `:855` | 同上 | 进 |
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | 同文件 `:314` | 同上 | **进，但会具名 skip**（矛盾②） |
| `TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe` | `cmd/wisp/panel_pageover_33r10_windows_test.go:241` | `//go:build windows` | 进（headless 那枚，`coldStartPageHandover` 拆出来的产物） |
| `TestPanelHostAC1*` 那一族"2s 子进程退出"条款 | `cmd/wisp/panel_host_windows_live_test.go` | **`//go:build windows && winlive`** | ⛔ **不进执行面**，只进编译面（见下） |
| winlive 全族（13 枚顶层用例） | 尺＝`git grep -lE '//go:build.*winlive' HEAD -- '*.go'` ⇒ **cmd/wisp 7 枚文件＋internal/ball 5 枚文件** | `windows && winlive` | ⛔ 只 `go vet` |

CI 那侧的尺（逐字）：

- 执行面：`ci.yml:605` `run: bash scripts/wisp-cli-tests.sh`（step 名 `cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)`，在 job `test-windows`＝`ci.yml:516`）→ 该脚本末行 `bash "$portable" --scope=cli`（逐字）→ `scripts/portable-tests.sh` 的 `cli` 分支逐字 `scope=(./cmd/wisp/)`（`pt.sh:256-261`）⇒ **整个 `./cmd/wisp/` 包跑进 CI，没有 `-run` 过滤器** ⇒ 上面所有 `//go:build windows` 的用例都在分母里。
- winlive 面：`ci.yml:655` `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`（step 名逐字 `winlive compile gate (go vet -tags winlive, ticket 111 AC#11)`）⇒ **⛔ 从不执行，只保证编译**。

**矛盾①（具名报回，⛔ 我裁）**：`panel_host_windows_test.go` 的头注释逐字写着 "These create an actual WebView2 window, so they belong to the winlive tier, which has NO CI job (recorded in the evidence table): they run only on this desktop session"（`:5-8`），又说 "a window that cannot be created here is a red, not a skip"（`:10-11`）——**但这枚文件的 tag 只有 `windows`，没有 `winlive`**（首行逐字 `//go:build windows`）。⇒ 注释声称的层级与文件实际层级不一致：按实际 tag，它**每天在 CI 的 windows 腿上开真窗**。要么注释过期、要么 tag 该补 `winlive`；两种都动的是别人的账，本腿只报。

**矛盾②（预测，⛔ 读数）**：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 在 `:317` 逐字 `t.Skipf("AC#13 has no subject in this tree: the embed resolves no entry, …")`，而触发前题在 CI 上**天然成立**：`npm run build` 只存在于 node-only 的 `lint-frontend` job（`ci.yml:971-973` `name: build (vite build -> frontend/dist, the bytes go:embed carries)`，同 job 注释 `ci.yml:931` 逐字 "every step here is node-only"）⇒ `test-windows` 腿上 `frontend/dist` 只有 `.gitkeep` ⇒ `Built()==false` ⇒ 这发 skip。而 `pt.sh:18` 的规则逐字 "any top-level `--- SKIP` in what actually ran is fatal"，尺＝`grep -c 'ColdStart\|AC13\|AC14'` 在 `scripts/portable-tests.sh` 里＝**0**（skip 台账 `ledger=(` 在 `pt.sh:589`，名册里没有任何 panel/AC13/AC14 条目）。⇒ **预测：这枚用例在今天的 CI windows 腿上会把整步打红**；⛔ 本腿没跑 `go test`（禁），⛔ 把这句当读数，也⛔ 与"本机红"混写——本机那一发由 `296-v1`/编排者量，本腿只交尺。
⚠ 同框的第二枚过期指认：`scripts/wisp-cli-tests.sh` 头注释逐字 "on windows its **33** top-level cases go PASS=33 FAIL=0 SKIP=0"（`:7-10` 附近），而本腿在 HEAD 上数的顶层用例枚数（尺＝对 `git ls-tree -r --name-only HEAD cmd/wisp | grep '_test\.go$'` 逐枚 `git show HEAD:<f> | grep -c '^func Test'`，按 tag 分箱）＝**非 windows-tag 133 ＋ windows-only 151 ＋ winlive 13**（78 枚 `_test.go`）⇒ "33" 是那一次测量的时代读数，⛔ 今天的分母。

## Q4 票面 AC#1 的措辞成不成立？—— **时序那一支成立；对象／回调配对／`contentTypeOf` 可达性三支不成立或不完整**

票面 `AC#1` 逐字（本腿从 `git show HEAD:<票面>` 现读）：「宿主在**建窗之后、供页之前**注册资源过滤器，回调里走 `Assets.Resolve(path)` ＋ `contentTypeOf`」。

1. **"建窗之后"——成立**。CoreWebView2 是 controller-completion 回调的产物：`pkg/edge/chromium.go:186` `func (e *Chromium) CreateCoreWebView2ControllerCompleted(res uintptr, controller *ICoreWebView2Controller) uintptr {`，注册所需的 `e.webview` 与事件挂载（`:211 vtbl.AddWebResourceRequested.Call`）都在这个体内，且尾部 `:227` 才 `atomic.StoreUintptr(&e.inited, 1)`。⇒ 在那之前**没有可注册的对象**，票面这句的时序与依赖语义同向。
2. **"供页之前"——成立，但比它写的更严**：处理器是库无条件挂的（同上一格），所以"先注册再供页"是**唯一**有效次序；再加一层本腿新量到的约束（件 `20` §3）：CoreWebView2 只有在**消息泵跑过之后**才存在，而 `bringUp` 是在一枚专属面板线程上跑的（库在 `webview.go:108` 记下 `w.mainthread`）⇒ 注册那一发**必须派发到那枚线程**，⛔ 从别的 goroutine 抢跑就不是"供页之前"，是"对象还不存在"。
3. **"宿主注册"——按真实 API 语义不成立（缺对象）**。本仓宿主握的是 `panel_host_windows.go:148` `w webview2.WebView`，`:386` 由 `webview2.NewWithOptions(...)` 造。枚接口成员（`common.go:26-84`）只有 `Run`/`Terminate`/`Dispatch`/`Destroy`/`Window`/`SetTitle`/`SetSize`/`Navigate`/`SetHtml`/`Init`/`Eval`/`Bind`；底层那枚 `*edge.Chromium` 被 `webview.go:52` 的**非导出字段** `browser`（类型＝非导出接口 `browser`，`webview.go:38-47`）吃掉了。⇒ 票面这句**省略了"宿主如何取得注册对象"这一环**，而那一环不是补一行调用，是换所有权（要么自己 `edge.NewChromium()`＋自建窗，要么改依赖）。件 `30` §1 已把这枚代价落到文件枚数。
4. **"回调里走 `Assets.Resolve(path)` ＋ `contentTypeOf`"——半句按包边界不成立（具名冲突，⛔ 不改票面一字）**：`contentTypeOf` 是 `internal/panel` 的**非导出**函数（`internal/panel/assets.go:165` 逐字 `func contentTypeOf(name string) string {`），宿主在 `package main` ⇒ ⛔ 调不到。可用的只有 `Resolve` 的第二返回值（`assets.go:75` `func (a *Assets) Resolve(requestPath string) ([]byte, string, error)`、`assets.go:93` `return data, contentTypeOf(name), nil`）。⇒ 正确说法应是"回调里走 `Assets.Resolve(path)` 拿 `(bytes, contentType)`"，`contentTypeOf` 这半句**只能作为"Resolve 内部做什么"的说明**存在，⛔ 作为宿主可点的调用。票面现量第 4 条自己也引了 `assets.go:161+` 那句注释，可见作者本意是"content type 那一层不许自己拼"——**意图成立，写法越界**。
5. **"把 `./assets/*` 从内存喂回去"这一目标——只靠注册达不到**（Q1）：两枚 URL 都是相对形状而文档基按仓内/库内两句原文是 `about:blank`；且只调 `AddWebResourceRequestedFilter` 而不设 `WebResourceRequestedCallback`（`chromium.go:42`）＝请求来了没人答（`chromium.go:286` 的判空）。⇒ 本腿判：**AC#1 需要补三枚从句**（取得注册对象的方式／回调配对／供给形状与 origin 的关系），其中第三枚从句一旦写清就直接命中票面 `AC#3` 的"停手上报由我裁"——那格不是偶然会触发，是**必然触发**。

**本节（本件）终态**：Q1 甲形不自足（高置信、待一发 `baseURI` 读数钉死）／Q2 今天零枚仪器看得见 CSP，两形各一枚验法已交／Q3 `//go:build windows` 那几枚**确实进 CI 执行面**（`scope=(./cmd/wisp/)`，无 `-run`），`winlive` 13 枚只进编译面；顺带两枚过期指认（注释层级与 tag 不符、"33 枚"分母过期）与一枚"AC#13 探针无牙"的读数／Q4 AC#1 的时序成立、其余三支不完整，具名错处＝**注册对象缺失**、**回调配对漏写**、**`contentTypeOf` 跨包不可达**。
