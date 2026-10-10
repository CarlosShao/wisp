# 票 299 — 面板那张页面**按 URL 去取自己的 js/css**这条路今天没有注册方：窗建得出来、屏上是全白（D29 写死的 `AddWebResourceRequestedFilter` 在产码里 0 枚）

**立票**：2026-10-10 08:5x 编排者（来路＝机主 08:3x 批 `Q-83` 甲「把带页面的版本跑出来」⇒ 我亲跑出厂链＋亲跑真机那一发；件 `.scratch/wisp/probes/orch/2026-10-10-q83-build-panel-verify.md`，账 `A794`）
**性质**：⚠ **这不是新需求，是把一张早就登记过、却没有格子的缺口立成能派单的票。** 两处旧账逐字在盘上：
- `docs/PLAN.md:1009`（D29 那张表的"资源加载"那一行）逐字 `| 资源加载 | \`AddWebResourceRequestedFilter\` 从内存喂 \`embed.FS\` | **不起 localhost HTTP 服务** —— 避免 dsh-tauri 的 iframe-to-localhost 模式（该模式已被 D21 批评） |`
- `docs/evidence/s1/33-35-preflight.md:81` 的 S4 行：「页面资源服务：`AddWebResourceRequestedFilter` → `Assets.Resolve(path)` 喂嵌入字节 | 被调物已存在 `internal/panel/assets.go:75`；**注册方零处** | **完全不存在（票 33）**（Go 侧 API 齐、宿主侧空）」
⇒ 那次普查把它归给票 33，而**票 33 从头到尾没有一格管这件事**（见"为什么票 33 没背"）。

**严重性＝高（用户看得见的那一层）**：今天有三条链（票 33 宿主／票 35 桥与泵／票 145 快照扩字段）都在往那扇窗里送数据，而**窗里那张纸自己取不进来**。⛔ 不许因为"面板还没有真内容"就当低优先级——**面板有内容也画不出来**是同一枚原因。

## 现量（每条带尺；⛔ 引用前先重跑，行号是快照）

1. ★**真机那一发（本票的第一手凭据，我 08:53 亲跑）**：`build/wisp.exe`＝**今天 08:51 出厂链那一发**（`doctor` 报 `commit=29081a13`；`grep -c -a -o 'index-B8yINMF1.js' build/wisp.exe` ⇒ **2**、旧名 `index-BVKlegVD.js` ⇒ **0**、`id="root"` ⇒ **1**；PE 子系统＝**2 GUI**；`sha256sum -c build/SHA256SUMS` ⇒ 四枚全 `OK`）⇒ **页面字节确实在库里、确实被嵌进去了**。合成 `Ctrl+Alt+P` ⇒ 日志 `wisp: panel window is up (panel-hotkey, cold -1.0 ms, hot path 0.0 ms)`，`list_windows` ⇒ `title="Wisp panel"` `bounds 640x260` ⇒ **窗真建出来了**。⚠ **截图＝纯白客户区**；无障碍树 22 枚节点里只有 `region Wisp - Web content` 套了七层**空 `region`**，**零文本、零控件**。
2. **入口 HTML 要两枚外部文件**（尺＝`cat frontend/dist/index.html` 现读）：`<script type="module" crossorigin src="./assets/index-B8yINMF1.js">`（551,989 字节）＋ `<link rel="stylesheet" crossorigin href="./assets/index-yy8KMgdf.css">`（49,540 字节）；`frontend/vite.config.ts:23` 逐字 `  base: "./",` ⇒ **拼法是相对路径**（不是绝对、不是 CDN）。
3. ★**注册方 0 枚**（尺＝`git grep -in 'WebResourceRequested' HEAD -- ':!.scratch'`）⇒ 命中只有：`cmd/wisp/panel_host_windows.go:42`（**注释**）、`docs/PLAN.md:1009`、`docs/evidence/s1/` 三份文档。**产码里 `AddWebResourceRequestedFilter`／回调注册＝0 枚。**
4. **被调物已经齐**（这是本票最省事的地方）：`internal/panel/assets.go:75` `func (a *Assets) Resolve(requestPath string) ([]byte, string, error)`（第二返回值就是 content type）、`:126` `entryRefRe` 就是为匹配 Vite 那两枚引用写的、`:161+` `contentTypeOf` 的注释逐字写着 "Anything unlisted is **refused by the host** rather than served as text/html"、`:130` `Check()` 的 doc **逐字预警过今天这个形状**：「the entry file is there, the page loads, and then **every hashed asset it names 404s, so the panel renders blank while the binary still reports built=true**」。⇒ 缺的只有"宿主把回调挂上"那一跳。⚠ 但 `Check()` 只证"字节在同一棵树里"（它今天**通过**——`cmd/wisp/panel_assets.go:118` 是它唯一的生产调用者），**它证不了"运行时取得到"**，所以三绿同放不是假设，是今天的实况。
5. **今天供给只有一整包**：`cmd/wisp/panel_host_windows.go:486` `w.SetHtml(string(data))`（`serveEntry` 体内）；同文件 `:95-96` 注释逐字 `// SetHtml produces (about:blank): offline, no host name, no port.`；`:49` 逐字 `// Today serving goes through SetHtml over embedded bytes (still offline, still no listener).` ⇒ **文档基是 `about:blank`，相对 URL 没有可解析的基**。
6. **为什么冷启动那发探测"看起来通"**：`panel_host_windows.go:863-866` 那枚自造页是**内联**单文件（`<script>if(window.wispProbeRT)window.wispProbeRT();</script>`，零外部引用）⇒ 它不需要这条路。⚠ **本票⛔ 不许把"探测回来了"当"页面渲染了"的凭据**（票 33 `AC#14` 那句定式同一个意思：「"窗口开出来了"与"回执到页面了"是两维」——这里是第三维：**"回执到了"与"屏上有内容"又是两维**）。
7. ⚠ **一条会挡死"偷懒修法"的现读**：入口 HTML 自带 CSP `<meta http-equiv="Content-Security-Policy" content="default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; ...">` ⇒ **`script-src 'self'` 不放内联脚本**，所以"把 js 内联进一整包再 `SetHtml`"这一支在今天的 CSP 下**天然不成立**；而那枚 `<meta>` 在 `frontend/index.html`＝**冻结面**（`docs/evidence/s1/33-35-preflight.md:123` 逐字："**属冻结**；Go 侧只做响应头注入"）。⇒ **只有 D29 写死的那条甲形走得通**，⛔ 不许为了绕 CSP 去动页面那枚 meta 一字。

## 为什么票 33 没背这一格（具名对拉，⛔ 不是抢票）

票 33 现有未勾的三格各管一件事：`AC#12`＝"能构建不等于有页面产物"（管 **dist 里有没有字节**）、`AC#13`＝"冷启动探测页不许盖住真页面"（管 **两发 `SetHtml` 的次序**）、`AC#14`＝"Go→页面这一跳由谁投递"（管 **`Dispatch`/`Run()` 队列**）。⇒ 三格**没有一格**管"页面自己按 URL 取它的 `.js`/`.css`"。`AC#13` 判据①里那句"或改用资源请求过滤器"是**并列备选路**、不是判据本体（次序那一支已经由 `coldStartPageHandover` 落地，见 `panel_host_windows.go:429-456` 的注释逐字 "AC#13's order, and it is the whole fix"）。

## 要建什么

- [ ] **AC#0 只读普查（⛔ 不许只答"没找到"）**：① 把入口 HTML 在**运行时**会去取的**每一枚** URL 现抽成名册（尺＝对 `frontend/dist/index.html` 跑 `entryRefRe` 同形正则＋对 js 产物再扫一层 `import(`／`new URL(...import.meta.url)` 形状，⚠ 枚数必须写清"哪把尺＋射程文件＋含不含注释"），并逐枚判 `Assets.Resolve` 取不取得到；② 把依赖侧**可用的注册 API** 逐枚列名（`AddWebResourceRequestedFilter`／`add_WebResourceRequested`／`WebResourceRequestedEventData`／`put_...`／`GetResponse` 之类，**带依赖模块 file:line**，⛔ 只写名字不算交付），并具名回答"注册需要哪一枚控制器对象、本仓宿主现在手里有没有那枚对象"；③ 出**两形代价表**（甲＝注册过滤器；乙＝内联整包）各自"要动哪几枚文件／会不会撞冻结的 CSP meta／要不要新增依赖边"。⛔ 不改码。
- [ ] **AC#1 落地（甲形＝D29 写死那条）**：宿主在**建窗之后、供页之前**注册资源过滤器，回调里走 `Assets.Resolve(path)` ＋ `contentTypeOf`（⛔ 不许自己拼 content type、⛔ 不许新造第二枚资源真相源——`embed.FS` 那棵树就是唯一的一枚），把 `./assets/*` 从内存喂回去。⛔ **不起 localhost HTTP**（D29 逐字禁）、⛔ 不加协程（D38(b) 名册零膨胀）、⛔ 不动 `frontend/**` 一字（含那枚 CSP `<meta>`）。完成判据＝**真机一次"屏上有内容"**：同一枚 `Ctrl+Alt+P`，截图里要出现页面自己的文本节点（无障碍树里 `region Wisp - Web content` 之下**不再是七层空 region**）。
- [ ] **AC#2 判据要有牙（⛔ 不许拿"注册被调用过"当凭据）**：ⓐ 由**页面自己报回**"我的 js 与 css 都到了"（形如票 33 `AC#14` 要求的"由页面自报"定式，先例＝`A502` 那一枚），断言的是**页面报的内容**，不是 Go 侧调用计数；ⓑ **反控**＝把注册那一行摘掉 ⇒ 指名用例**必须红**（⛔ 恒绿＝它对这件事不敏感，本仓有先例 `A710`）；ⓒ 现有那几枚"窗建出来了／回执到了"的用例**一枚都不许多变绿**——它们管的是别的维度（本票现量第 6 条）。
- [ ] **AC#3 界**：⛔ 不动 `coldStartPageHandover` 的次序（票 33 `AC#13` 已落的那一发）、⛔ 不动泵/`Run()` 那一格（票 33 `AC#14`）、⛔ 不动 `Assets.Check()` 的语义。若注册之后 `serveEntry` 的 `SetHtml` 那一发**变成多余**（改成导航到虚拟主机 URL），⚠ **停手上报由我裁**，腿不许自己换供给形状。
- [ ] **AC#4 门禁＋越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；格式两把尺并排（工作树＋**仓外** blob 各一把，⛔ blob 不许落在仓内）新增 0 枚；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/panel/ -count=1 -v` 改前改后各 ≥2 发取交集＝新增红 0 枚（⚠ 三数必带尺名：`^--- PASS` 顶层 vs `--- PASS` 含子测试）；每枚门禁件自落一行 `rc=N`；`git show --stat` 名册只含 `cmd/wisp/**`＋`internal/panel/**`＋`probes/299/**`；⛔ `frontend/**`（含 `dist/**` 生成物）／`design/**`／三枚冻结件／`thresholds.go`／`allowlist.txt`／D43 表／`PLAN.md` 零字节；⛔ 零 push、commit 必带显式 pathspec。
- [ ] **AC#5 〔仅本机可量，⛔ 归编排者跑，派单就要写〕**：改前／改后各一次真窗截图对比（改前＝本票现量第 1 条那发纯白；改后＝屏上出现页面文本），＋跑前 `tasklist` `wisp.exe`／`balldebug.exe`＝0。⚠ 这一格⛔ 不算腿的欠账（真窗与截图设备在机主这台机器上）。

## 禁区

- ⛔ **不许用"塞一张静态占位页"来交差**（那会把这一格从"没做"变成"看起来做了"——与本仓"宁缺毋造"那条定式同族）。
- ⛔ 不起本地 HTTP 服务／不开监听端口（D29 逐字）；⛔ 不引新依赖边（票 244 那族尺）；⛔ 不动 `frontend/index.html` 的 CSP meta（冻结）；⛔ 不改 `//go:embed all:dist` 的活模式。
- ⛔ 不许顺手翻票 33 `AC#12`/`AC#13`/`AC#14` 任何一格（本票现量第 6、7 条与票 33 `AC#12` 的今日进展**只写 Progress log，框不动**）。
- git：只 commit 不 push；显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件、不建 worktree；证据件只新建 `.md`（⛔ `.out` 被根 `.gitignore` 第 8 行全仓忽略）。

## 排程与串行

- 写面＝`cmd/wisp/panel_host_windows.go`（＋同包测试）⇒ ⛔ 不与 `296-v1`（独占 Go 编译面）并发；⛔ 不与 `293-r1`（同包 `resident_windows.go`）同批。
- 队列位置：**按在 `296-v1` 交回并裁完之后**，排在 `293-r1` 之前还是之后由我下一轮裁——裁的根据是"这一格挡的是用户看得见的那一层，`293` 挡的是隐私看得见的那一层"，两枚都高，⛔ 不并发。
-  本票 `AC#0` 是只读腿，⛔ 禁跑 go 编译面；`AC#1`/`AC#2` 是落地腿，独占 Go 面。

**Status:** **未开工**（本票 08:5x 立，六格 `AC#0`..`AC#5` 全未勾；凭据已在盘上＝本票"现量"第 1 条那一发真窗，⛔ 那不是"待验假设"）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）

- [2026-10-10 10:00:25 +08] agent=299-a1 did=AC#0 只读普查交齐四件（`.scratch/wisp/probes/299/a1/{00,10,20,30,40}-*.md`，锚＝HEAD `6ef14788`，五笔 commit `e7368e10`/`677eea04`/`7354b9a8`/`38ed114e`/`68c7dd24`（＋本条 append 为第六笔）；⛔ 改码、⛔ 翻框、⛔ 编译面 0 发，只跑过 `go env`）。读数六枚：①★**名册尺在 HEAD 上不成立**——`git show HEAD:frontend/dist/index.html` ⇒ `fatal: … exists on disk, but not in 'HEAD'`，`git ls-tree -r --name-only HEAD -- frontend/dist` ⇒ **1 枚**（`.gitkeep`，`frontend/.gitignore:12 dist/*`），故名册取自工作树 08:51 那一发＝2 枚（`./assets/index-B8yINMF1.js`／`.css`，`entryRefRe` 同形尺含注释行、第二层 `import(`/`new URL(`/`import.meta.url`/`.wasm`/`Worker(` 全 0 新增；`fetch(`＝modulepreload 垫片今天 0 员）②★**注册对象＝`ICoreWebView2` 非 Controller**（`pkg/edge/corewebview2.go:383`），而宿主握的是 `webview2.WebView`（`panel_host_windows.go:148`/`:386`），`common.go:26-84` 十二枚成员零枚过滤器、底层 `*edge.Chromium` 被非导出字段 `webview.go:52 browser` 吃掉 ⇒ **API 齐（`AddWebResourceRequestedFilter`/`WebResourceRequestedCallback chromium.go:42`/`GetUri`/`CreateWebResourceResponse corewebview2.go:169`/`PutResponse`——派单点名的 `add_WebResourceRequested`/`WebResourceRequestedEventData`/`get_Response` 在本绑定 0 命中），缺的是所有权不是行数** ③**乙形结构性不成立**：唯一合规 CSP 出口（响应头注入）只活在甲形通道里，`SetHtml`＝`NavigateToString`（`webview.go:394`）无可注入响应 ④**Q1 甲形不自足**：两枚 URL 相对（`vite.config.ts:23 base: "./"`）而 `about:blank` 有两处原文（`host.go:95-96`＋`common.go:62`）⇒ 只注册不换供给大概率无人敲门；待钉的那一发＝`probeJS` 里把 `b` 换成 `document.baseURI`（尺名已交，⛔ 本腿未跑）⑤**Q3**：`//go:build windows` 那几枚**确实进 CI 执行面**（`ci.yml:605`→`portable-tests.sh cli` 逐字 `scope=(./cmd/wisp/)` 无 `-run`），winlive 13 枚只进 `ci.yml:655 go vet -tags winlive`；两枚过期指认＝`panel_host_windows_test.go:5-8` 自称 winlive 层但首行 tag 只有 `windows`、`wisp-cli-tests.sh` 的"33 枚"对不上 HEAD 现数（133+151+13）；一枚预测（⛔ 非读数）＝`TestAC13ColdStart…` 在 CI 上会 `t.Skipf`（`:317`）而 `pt.sh` 台账 0 条 AC13 条目、`pt.sh:18` 规定 SKIP 即 fatal ⑥**Q4 AC#1 时序成立、三支不成立**＝缺注册对象取得方式／漏写回调配对（只调过滤器则 `chromium.go:286` 判空什么都不做）／`contentTypeOf` 非导出（`assets.go:165`）宿主在 `package main` 调不到，可用只有 `Resolve` 第二返回值。另：CSP 今日零枚仪器可见（尺＝`git grep -n -i 'Content-Security-Policy\|CSP' HEAD -- ':!.scratch'` ⇒ 37 命中、产码侧 0），且 `TestAC13…` 的探针 id 取自入口静态字节、今天只有 `id="root"` 一枚 ⇒ **它对"js 跑没跑"无牙，AC#2 ⓐ 不可复用该形状**。next=⛔ 编排者裁三件事：①AC#3 那格（供给形状能否从 `SetHtml` 换成带 host 的导航）按本腿判是**必然触发**而非偶然，请在派落地腿前先裁；②`cmd/wisp` 若要注册就必得自持 `*edge.Chromium`（`pkg/edge/chromium.go:47`）＝包级 import +1（模块边 +0，AC#4 名册容得下）还是走依赖改动（撞"不引新依赖边"）；③矛盾①②与 Q4 三支的票面措辞归你落格，本腿⛔ 动票面一字。


