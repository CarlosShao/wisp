# 299-a1 · 30-two-shape-cost-table — 甲（注册资源过滤器）／乙（内联整包）代价表

锚＝HEAD `6ef14788`；对象可达性判语的依据全在件 `20`，资源枚数依据在件 `10`，本件不重复取证、只算代价。
⛔ 本腿不起 localhost HTTP（D29 逐字禁），所以表里没有"丙＝本地端口"那一支——那是**禁区**，不是候选。

## 1. 终表（四问逐枚答，⛔ 不给"推荐"）

| 代价项 | 甲＝注册资源过滤器（票面 D29 写死那条） | 乙＝内联整包 |
|---|---|---|
| **要动哪几枚文件（产码）** | **≥2 枚，且第 1 枚是重写级别**：① `cmd/wisp/panel_host_windows.go` —— 必须换掉"宿主手里握着什么"：`:148` `w webview2.WebView` ＋ `:386` `webview2.NewWithOptions(...)` 这一对拿不到注册对象（件 `20` §3），要拿到 `*edge.Chromium` 就得自己 `edge.NewChromium()`（`pkg/edge/chromium.go:47`）并自己承担建窗/STA 线程/消息泵——库里那些动作在 `webview.go:97-130`（`chromium := edge.NewChromium()` 是该函数体内第 6 行）与 `CreateWithOptions`（`webview.go:269`）里，是**非导出类型的方法**，包外一行都调不到；② 回调本体＋`Assets.Resolve` 接线（同一枚文件或新增一枚 `cmd/wisp/panel_asset_filter_windows.go`，枚数由落地腿定）。**`internal/panel/assets.go` 可以 0 字节**：`Resolve` 第二返回值就是 content type（`assets.go:93`），⛔ 不需要为"导出 `contentTypeOf`"动它。测试面另计（票面 AC#2 ⓐⓑ 至少要 1 枚新用例＋1 枚反控）。 | **2 枚**：① `internal/panel/assets.go`（新增一枚导出函数把 entry＋`./assets/*` 拼成一整包文档）；② `cmd/wisp/panel_host_windows.go` 的 `serveEntry`（`:472-488`，今天它只做 `m.assets.Resolve(panel.EntryFile)` ＋ `w.SetHtml(string(data))`）。 |
| **会不会撞那枚冻结的 CSP `<meta>`** | **不撞，但甲必须让文档换个 origin 才有效**（这是本件最重要的一格，展开见 §2）。meta 本体一字不动（`frontend/index.html`＝冻结面；尺＝`git grep -n -i 'Content-Security-Policy\|CSP' HEAD -- ':!.scratch'` ⇒ 37 枚命中里 **`cmd/`／`internal/`／`tools/`／`scripts/` 产码侧 0 枚**，只有 `frontend/index.html` 2 枚＋docs/evidence 若干）。 | **撞，而且是结构性撞、不是代价问题**：入口 meta 逐字 `script-src 'self'`（⛔ 无 `'unsafe-inline'`，`style-src` 才有），乙形要写的正是 `<script type="module">` 内联体 ⇒ 天然被拒。票面现量第 7 条给的唯一出口是"Go 侧只做响应头注入"（`docs/evidence/s1/33-35-preflight.md:123`），但**响应头注入只活在甲形里**——`CreateWebResourceResponse` 的 headers 参数在 `pkg/edge/corewebview2.go:169`，而 `SetHtml` 走的是 `NavigateToString`（件 `20` §3 末）⇒ **根本没有一枚可注入的响应对象**。⇒ 乙的两条出路（动冻结 meta／改用甲的响应通道）一条撞禁区、一条不再是乙。 |
| **要不要新增包级依赖边** | **包级 +1，模块级 +0（若走甲-1）**：`cmd/wisp` 现在只导顶层包（尺＝`git grep -n 'jchv' HEAD` ⇒ 产码/测试 5 枚文件全是 `webview2 "github.com/jchv/go-webview2"`，`panel_host_windows.go:64`），甲-1 要新增 `github.com/jchv/go-webview2/pkg/edge` 这一枚**包级 import**；`go.mod`/`go.sum` 同模块不改一字节 ⇒ AC#4 的 `git show --stat` 名册（只含 `cmd/wisp/**`＋`internal/panel/**`＋`probes/299/**`）**容得下**。⚠ 但如果落地时选的是甲-2（换/patch 依赖版本、`replace`、或另引一枚封装库）⇒ 动 `go.mod`/`go.sum`/`deps.toml`/`allowlist.txt`，**同时撞 AC#4 名册与"⛔ 不引新依赖边"两格**。 | **0 枚**（纯 Go 字符串拼装，`strings`/`bytes` 已在 stdlib import 面里）。 |
| **要不要动 `coldStartPageHandover` 的次序（票 33 `AC#13` 已落，⛔ 本票不许动）** | **次序不动，但注册点必须落在它之前、且落在那枚线程上**。`coldStartPageHandover` 是 `panel_host_windows.go:447-454`，体内逐字 `rtMs := m.firstRoundTrip(ctx, t0)` → `if err := m.serveEntry(); err != nil {` → `m.serveNotBuiltNotice()`；它的注释（`:429`）逐字 "AC#13's order, and it is the whole fix"。甲形要的是**第三件事**（注册过滤器），插在 `bringUp` 里、`Embed` 之后、`coldStartPageHandover` 之前即可 ⇒ 两发 `SetHtml` 的相对次序**一枚都不必换**。⚠ 两个真约束：① 事件处理器是库在 controller-completion 回调里挂的（`chromium.go:186` 体内 `:211`），CoreWebView2 只有在消息泵跑过之后才存在 ⇒ 注册**必须派发到这枚面板线程**（`webview.go:108` 记下 `w.mainthread`），⛔ 不能从别的 goroutine 抢跑；② 若甲形把供给形状从 `SetHtml` 改成"导航到虚拟主机/自定义 scheme URL"，`serveEntry` 那一发就**变成多余**——这正是票面 `AC#3` 写死的"⚠ 停手上报由我裁，腿不许自己换供给形状"。本件据 §2 判断：**这一支不是"可能会变成多余"，是甲形成立的必要条件**，所以 AC#3 那格裁度**必然会被触发**，请编排者提前裁。 | **不动**（乙只换 `serveEntry` 里喂进去的字节，`firstRoundTrip` 那发探测页与它的相对次序原样保留）。⚠ 另有一笔没被上面四问覆盖的代价：内联把 551,989＋49,540 字节（件 `10` 实测）塞进单发 `NavigateToString`，与 D32 的"面板拉起 ≤500ms"预算同框——本腿**没有**量过这一发，⛔ 把这句当读数。 |

## 2. 甲形的真门槛：不是"能不能注册"，是"文档有没有 origin"

件 `10` 的名册里两枚 URL 都是**相对**引用（`./assets/…`，`frontend/vite.config.ts:23` 逐字 `  base: "./",`）。宿主注释（`panel_host_windows.go:95-96`）逐字写着 `// SetHtml produces (about:blank): offline, no host name, no port.`，顶层库的 doc 也逐字写着（`common.go:62`）`// The origin of the page is `about:blank`.`。⇒ **只注册过滤器、不改供给形状，`./assets/*` 能不能被请求到**这件事是甲形的生死题，本腿把它作为 **Q1 单列在件 `40`**，⛔ 不在本件凭 API 名字推。

对代价表的直接影响：**甲-1 的"要动哪几枚文件"里还欠一枚"文档 origin 形状"的改动**（自定义 scheme 或 `Navigate` 到某枚有 host 的 URL），而那一枚改动的性质就是 AC#3 要裁的那一件事。把两格合起来读，甲形的真实射程＝**注册（可达性）＋ origin（可请求性）两枚门槛同时过**，缺一枚都不会让屏上有内容。

## 3. 已存在的仪器会怎么反应（避免"改完才发现撞钉"）

- `cmd/wisp/panel_host_gate_test.go:30` `TestPanelHostOpensNoListeningSocketL1` 用 `go/parser` 只读**那一枚宿主文件**，体内逐字断言 `panel host imports %q, which can open a listening socket; D29/AC#3 forbid a localhost server`（`:42`）⇒ 甲形/乙形都⛔ 不能给宿主文件添 `net`／`net/http`。
- 同文件 `:77` `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` 已经在**钉件 `10` §0 那件事**：`:143` 逐字 `AC#12's two shapes have just been conflated: the binary carries a page (built=true) while this tree tra…`。⇒ "HEAD 只跟踪锚／exe 里带着真页"这一区分**已有仪器看护**，本腿 §0 的读数与它同向，不是新造判据。
- `internal/panel/assets.go:126` 的 `entryRefRe` 若因页面形状变化（乙形内联后入口不再写 src/href）而枚数漂到 0，`Check()`（`:137-160`）会返回**空 served 而不报错** ⇒ 乙形有一枚"仪器静默变宽"的次生代价；甲形不动入口文档，`Check()` 语义原样（票面 AC#3 那格"不动 `Assets.Check()` 的语义"对甲自动满足）。

## 4. 〔越界候选〕两枚（写出来，⛔ 本腿不采纳）

1. **`SetVirtualHostNameToFolderMapping`**（`pkg/edge/ICoreWebView2_3.go:22`；可达性比甲形好——`Chromium` 上有导出的 `GetICoreWebView2_3()`，同文件 `:58`，ICoreWebView2 上也有 `:46`）：它给文档一枚真 host ⇒ 相对 URL 有基 ⇒ `'self'` 有归属。⛔ 但它的入参是 `folderPath string`——**一枚磁盘目录**，而 `//go:embed all:dist` 的字节只在 exe 里，要用它就先得把 dist 落盘 ⇒ 直接造出票面 AC#1 明令禁止的"第二枚资源真相源"，还要添一枚写临时目录的面（`risk.PathResolver` 规矩、`AGENTS.md` §1.2 那条 `filepath.Clean|Abs` 禁区）。⇒ 记为**越界候选**，交编排者裁，本腿⛔ 当方案。
2. **`Navigate("data:text/html,…")`**：库顶层接口确实支持（`common.go:56-59` 的 doc 逐字提到 data URI）。⛔ 入口 meta 是 `default-src 'none'` 且 `script-src` 只有 `'self'`、`img-src` 才带 `data:` ⇒ data 文档里的脚本没有许可；再加 551,989 字节进 URL，形状本身就越界。

## 5. 本节终态

**乙形不是"贵"，是"不成立"**：它唯一合规的 CSP 出口（响应头注入）本身只在甲形的通道里存在。**甲形成立，但票面 AC#1 那句只覆盖了甲形两枚门槛里的一枚**（注册可达性），另一枚（文档 origin／可请求性）落在 AC#3 明令"停手上报由我裁"的射程内 ⇒ 本腿预测 AC#3 那格**必然会被触发**，建议编排者在派落地腿之前先把"供给形状能不能从 `SetHtml` 换成带 host 的导航"裁掉。
