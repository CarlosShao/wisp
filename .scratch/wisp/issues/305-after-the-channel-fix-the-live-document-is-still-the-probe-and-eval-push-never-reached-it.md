# 票 305 — 面板**冷启动之后"活文档"到底是谁**：通道已经修好（页面能答了），但内嵌入口没交接上来，而 `Eval` 推不进文档那一枚**早在回归之前就红、从来没绿过**

**立票**：2026-10-10 18:0x 编排者（来路＝票 303 修复腿 `303-r1` 交件时**顶回我票面那句"三枚回绿"**；我随后自己跑了三发把它证成）
**性质**：★**缺陷票**，两枚症状、**⛔ 一条因果**（`AC#0` 必须先分开，混一票必然顺手改错那半）。
**为什么要紧**：这就是用户真正会碰到的那一屏——**面板打开了，但里面不是给他看的那页**。通道修好之后红句从"页面没答话"换成了"页面答了话、可它答的是探针那张文档"，⇒ **症状从"链路断"升级成"链路通而内容错"**，这对后面每一张要往面板里放东西的票（186／187／167／145）都是前提。

## 现量（每条带尺；⛔ 引用前先重跑，行号与枚数一律当快照）

- ★**母仓 @ `807497c1`（修复已在树里）**，尺＝`go test ./cmd/wisp/ -count=1 -timeout 420s -v -run '<三枚名>'`，件 `probes/303/orch/r2-orch-targeted-after.txt`（我 18:0x 现跑，`rc=1`）：
  `TestAC14AwaitedBindingReplyReachesThePage` ⇒ **PASS(0.55s)**，页面自己的话逐字 `REPLIED,REPLIED,REPLIED`（`Go's handler was reached by 3 of the 3 real requests`）。
  `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` ⇒ **FAIL(1.30s)**，红句逐字 `after a real cold start the live document contains NONE of the 1 element ids the embedded entry declares (the page itself answered "0")`（出处 `cmd/wisp/panel_resident_windows_test.go:330`）。
  `TestAC14GoSideEvalPushReachesThePage` ⇒ **FAIL(0.54s)**，红句逐字 `Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK"`（`:869`）。
- ★**同一枚仓外 clone 上的前后对照**（这是本票最关键的一把尺：台面相同、只有 commit 不同）：
  `f718e9b6`（＝回归笔 `fb2fb802` 的父）⇒ 件 `probes/303/orch/r3-f718e9b6-three-nails.txt`：nail1 **PASS**、**nail2 FAIL 且红句与今天逐字同形**（`title as "", want "PUSHED-33R5-OK"`）、AC13 **SKIP(0.00s)**。
  `807497c1`（修复后）⇒ 件 `probes/303/orch/r4-clone-at-HEAD-nails.txt`：nail1 **PASS**、nail2 **同一句 FAIL**、AC13 **SKIP**。
  ⇒ **判语①**：nail2 的那枚红**⛔ 是 `fb2fb802` 造成的**——它在回归之前的同一台面上就是同一句话。**"Eval 推不进文档"是一枚独立的、更早的缺陷**（本票症状②）。
  ⇒ **判语②**：AC13 在这枚 clone 里两发**都 SKIP** ⇒ 它的颜色**只有母仓台面能量**（先例＝`frontend/dist` 只跟踪 `.gitkeep`，本机有旧产物、干净 clone 没有）⇒ 本票⛔ 拿 clone 的 SKIP 当"它本来就好"，也⛔ 拿它当"它坏了"。
- ★**母仓改前那一发**（腿件 `probes/303/r1/20-targeted-before-red.txt`，`rc=1`）：三枚全 FAIL 且都是 `no report … from the page within 15s (what DID arrive at the door: nothing at all)` ⇒ 通道那一层确实被 `fb2fb802` 断过、也确实被 `303-r1` 接回来了（我这发的 PASS 是凭据）。⚠ **AC13 在母仓⛔ 已知"改前绿"过**（改前那一发它红的是**另一句**，因为通道断在前面挡住了它）⇒ 归因欠着，见 `AC#0` ②。
- **腿具名的嫌疑笔（〔仅腿报〕，我⛔ 复跑过归因）**：`70b00885`（`bringUp` 把探针文档当最后活文档那一族的引入处）。⛔ 把它当结论引。

## 要建什么（`AC#0` 之前⛔ 任何产码）

- [x] **`AC#0` 只读归因＋分层普查**：交回四张表——
  ① **"活文档"名册**：`bringUp` → `firstRoundTrip` → `serveEntry` 这条装配序里，每一次 `SetHtml`/导航各自把哪份文档变成活的、最后一次是谁（尺＝逐枚调用点文件:行＋调用序；⛔ 按注释读）。
  ② **两枚症状分开归因**（各用一枚**同台面**的前后对照，⛔ 拿 clone 的 SKIP 充数、⛔ 拿不同台面相减）：症状①＝AC13（入口没交接）；症状②＝nail2（`Eval` 推不进）。⚠ 若某一枚在**能判死的台面上从来没有绿过**，就写"它可能⛔ 是一枚回归、而是一枚从没做通过的功能"，并具名说凭据是哪一发。
  ③ **判据是不是同一条**：nail1 绿／nail2 红 是同一枚用例文件里的**两个维度**（awaited reply vs Go-side push，票 33 `A475` 那批裁过"两维⛔ 混一枚"）⇒ 本票必答"push 那一维今天缺的是哪一跳"。
  ④ **代价表**：修①要动 `bringUp` 的文档交接次序（**第二枚未归因的产码改动**）；修②要动 `Eval` 那一路（今天是不是根本没人调用它？先跑调用点尺）。＋**必答**：动②会不会撞上票 33 已裁的那条"面板线程＝专用 STA ＋ 库 `Run()` 泵"（`Eval` 只能在泵线程上跑）——撞就具名报回、由我裁，⛔ 自己绕。
- [ ] **`AC#1` 症状①（入口交接）落地**：只在 `AC#0` ② 归因之后开工。**硬约束**＝⛔ 放宽 `TestAC13…` 那句判据（"live document 必须含内嵌入口声明的那些 element id"是票 33 立的行为判据，⛔ 改成"含任一 id"或改成 `t.Skip`）；⛔ 为了让 AC13 绿而动 `firstRoundTrip` 探针文档的存在性（探针那发⛔ 是"最后活文档"这件事本身就是要修的对象）；⛔ 动 `frontend/src/**`（页面源码不在本仓射程，读一律 `git show <ref>:<path>`）。
- [ ] **`AC#2` 症状②（`Eval` 推不进）落地**：同上，⛔ 放宽"页面自己报 title"那句判据（判据必须**由页面报回**，先例＝票 33 `A502` 形）。若 `AC#0` 判出"这一维从来没有绿过 ⇒ 是缺功能⛔ 是回归"，本格**自动降级为⛔ 开工**，由我另立落地票。
- [ ] **`AC#3` 反恒真成对两发**：每枚症状各一对＝(a) 撤掉修复 ⇒ 目标用例当场按预期红（**具名红句逐字引**，⛔ 只报枚数）；(b) `cmp` 逐字节还原 ⇒ 复绿。⚠ 外加一发**跨维度负控**：修①的那笔⛔ 可能让 nail2 变绿——若变绿，必须具名说明为什么（否则就是"两枚症状其实同源"，本票的分层假设作废）。
- [ ] **`AC#4` 门禁与越界**：`go vet ./cmd/wisp/` 与 `go vet -tags winlive ./cmd/wisp/ ./internal/ball/` 各 rc=0；`sh scripts/d22scan.sh` rc=0；格式名册**两把并排**（工作树／HEAD blob，各自写射程）；**整包 `go test ./cmd/wisp/ -v` 改前改后各一发取交集＝新增红 0 枚**（⚠ 两发必须**同一台面**，见下面"规矩"里那条）；`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺）；每把门禁件自落一行 `rc=N`；⛔ 零 push；commit 显式 pathspec 写在 `$( … )` **之外**、只落自家 `probes/305/<腿名>/**`（票 301 `AC#4b`）。
- [ ] **`AC#5` 与票 33 的对账**：本票收掉的两枚症状对应票 33 的 `AC#13`／`AC#14` 哪几格，逐枚写"⛔ 翻票 33 那格／翻本票本格"，⛔ 两本账同时记同一件事（先例＝`A817`/票 302 那处"同一物理缺陷链上只记一次"）。

## 边界（⛔ 塞进本票）

- ⛔ **`fb2fb802` 那一跳的回执通道**＝票 303 已闭（`AC#3` 的修复），本票只在它之上判"内容对不对"。
- ⛔ **`--scope=cli` 里这三枚该不该搬档**＝票 302（现在**仍未裁**，且它的甲／乙裁语本来就按在本票之后）。⚠ 顺序⛔ 倒：先修内容，再谈搬档——搬档会把"链路通而内容错"这枚进步藏起来。
- ⛔ **`msedgewebview2.exe` 那 18/24 枚**＝机主自己的应用（`A821`＋`A823`），⛔ 残留、⛔ 本票射程。
- ⛔ **`internal/risk` 那 12 枚 CI 红**＝`A817` §3〔待归因〕。

## 规矩（本票全程）

- 子代理**只 commit、不 push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 在仓内建 worktree 或 checkout；临时件**只建不删**；证据件⛔ 叫 `.out`。
- ★**台面规则（本票立的一条新规矩，来自我这次自己抓到的坑）**：**前后对照必须在同一枚台面上跑**（母仓↔母仓、clone↔clone）。母仓带 `frontend/dist` 旧产物、干净 clone 没有 ⇒ 同一枚用例会给出**不同颜色**（这次实测：AC13 在 clone 两发都 SKIP、在母仓是 FAIL）。⇒ 凡跨台面比较，读数⛔ 可比，必须具名写成"两把台面"。
- ★★**台面还有第二轴（`305-a1` 现量、编排者复跑确认，2026-10-11）**：上面那条只覆盖了一根轴（母仓 ↔ clone）。盘上另有＝**同一枚母仓里那枚未跟踪 dist 自己会漂**——现量尺＝`wc -c frontend/dist/index.html`＋mtime＋`git cat-file -s HEAD:frontend/dist/index.html`＝`exists on disk, but not in 'HEAD'`（**1044 字节／mtime 2026-10-10 08:51**，HEAD 里只有 `.gitkeep`）。⇒ **"母仓↔母仓"的跨时间对照⛔ 天然同台面**：凡引母仓的绿↔红对照（含 10-01 那批 240-0-0 绿读数，它用的是**旧字节**），必须**同时**交出那枚文件的字节数与 mtime，否则那两发⛔ 可比。
- ★★★**台面第三根轴（`305-a2` 现量、编排者复跑，2026-10-11 08:0x）**：ⓐ母仓↔干净 clone、ⓑ母仓那枚未跟踪 dist 会漂之外，盘上还有＝**干净 clone／`git archive` 导出树⛔ 带 `third_party/sherpa-onnx` 与 `build` 这两枚 DLL 目录**（尺＝`git ls-files third_party/sherpa-onnx`＝**0**、`git ls-files build`＝**0**）。⇒ 把上面那条 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 照抄进 clone ⇒ 解析成 clone 里⛔ 存在的路径 ⇒ `0xc0000135` **且零 `--- FAIL`**＝"用例根本没跑"那形假绿。⇒ 两条落地规矩：① 在 clone／导出树里跑 `cmd/wisp`，PATH 一律写**母仓绝对** `/d/…` 形；② **第一发必须先断言那枚具名用例真发出了 `--- PASS/--- FAIL/--- SKIP`**，⛔ 拿退码或日志行数当"跑过了"。正文与判死配方＝本票"编排者收 `305-a2`"那一节第 2／3 节。
- ⚠ **嫌疑笔那一格已闭（2026-10-11，⛔ 改"⚠ 腿给的嫌疑笔 `70b00885` 是〔仅腿报〕"那句原句）**：`305-a1` 现量＋我自己读 `git show 70b00885 -- cmd/wisp/panel_host_windows.go` ＝`bringUp` 尾段四行**逐字搬进** `coldStartPageHandover`（同语句、同顺序、＋41/−11 除注释外⛔ 语义变化）⇒ **否证**，⛔ 再当归因结论引；判死那一截改走"窗口内 bisect"（`416d9d56`↔`f718e9b6`，见本票"编排者收 `305-a1`"那一节第 2 节）。
- ⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件；⛔ 为变绿放宽任何断言；⛔ 造 `--- SKIP`。
- 跑 `cmd/wisp` 测试必带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（★shell 形 `/d/…`，`D:/…` 那形 `0xc0000135` 且零 `--- FAIL`）；跑真机/整包前 `tasklist` 现量 `wisp.exe`／`balldebug.exe`＝0。
- 裁决者≠实现者（`SPEC-12 §4.3`）：`AC#1`/`AC#2`/`AC#3` 的判语⛔ 由落地腿自勾。
- ⚠ 腿给的嫌疑笔 `70b00885` 是〔仅腿报〕，本票 `AC#0` ② 未跑之前，任何文案⛔ 把它当归因结论写。

next=`303-v1`（非实现者裁票 303 `AC#2`/`AC#3`，含"这修复⛔ 放宽判据"那道必答题）→ `305-a1`（`AC#0` 只读，⛔ Go 编译面按补位腿规矩自报；若它要"改前必红"的夹具色，需先与编排者抢台面）→ 我裁 → 视 `AC#0` ② 的结论决定 `305-r1` 开⛔ 开
---

## 编排者收 `305-a1`（2026-10-11 07:3x，三笔 `5f891575`→`53c854e3`→`52a5eb23`，件 `probes/305/a1/**`；⛔ 我自己的复跑件＝这一节正文，逐枚尺写死）⇒ **`AC#0` 由我翻**＋★**裁 `AC#2` ⛔ 降级（那两枚是回归、有绿点）**＋★**更正我票面标题里那句"从来没绿过"**＋★**新规矩一枚：母仓↔母仓 对照必须带第二轴**

### 0. 台面与越界（⛔ 抄腿的自陈；腿⛔ 跑任何 Go 编译面，这一条我也替它核了）

- 起手锚＝`a81c2980`（＝我上一程推上去那笔，与派单一致）；`git status --porcelain -- internal cmd docs scripts .github`＝**0 行**；腿 `logs/l7` 自陈 28 枚行号锚在 HEAD blob 上逐枚 `git show` 核过全 MATCH。
- 逐笔越界尺＝`git show --name-only --format=<hash>`（⛔ 区间尺）我自己在三笔上重跑：**名册每一行都在 `.scratch/wisp/probes/305/a1/**` 下**；新建 `.go`＝**0** 枚（⇒ `gofumpt`/`d22scan` 分母零移动）；`.out`＝**0** 枚。
- 编译面自查：我复核它的 `logs/l3` 里只有 `go list -m -f '{{.Dir}}'` 那一发（`rc_golist=0`），⛔ `go test`／⛔ `go build`／⛔ `go vet` ⇒ 票 305 的"改前必红"那一族颜色**仍然只有我手里那五发**，⛔ 被第二人碰过。

### 1. 承重读数我自己在盘上现量（⛔ 抄腿；每条带尺）

| 表 | 腿的读数 | 我这发的独立复跑 |
|---|---|---|
| ① 活文档三份 | 最后一次写文档＝`:486` 入口（**1044 B**）或 `:467` 公告（**181 B**），⛔ `:892` 探针（**136 B**） | `sed -n '467,468p;484,488p;890,894p' cmd/wisp/panel_host_windows.go` 逐字复现三处 `SetHtml`；两枚字面量我把两段拼接串喂 `python len(encode())` 现算＝**181／136 字节逐字相等**；入口那枚 **1044 B** 量的是 `frontend/dist/index.html`（`wc -c`＝1044，mtime **2026-10-10 08:51**），而 `git cat-file -s HEAD:frontend/dist/index.html` ＝ **`exists on disk, but not in 'HEAD'`** ⇒ 那 1044 是**未跟踪的本机产物** |
| ① 只由 `SetHtml` 决定 | 生产面⛔ 任何导航调用点 | `grep -rn --include=*.go '\.Navigate(' cmd internal \| grep -v _test.go` ＝ **0 行**（库内唯一的 `browser.Navigate` 在依赖 `webview.go:390`，本仓⛔ 调用者）⇒ 成立 |
| ② ★曾绿 | `AC13` 与 `nail2` 在 **`416d9d56`**（2026-10-01 15:41，整包 **rc=0／240-0-0／红册空**）里**两枚逐字 PASS** ⇒ "从来没绿过"⛔ 成立 | `grep -n "PASS: TestAC13ColdStartEndsOnTheEmbeddedEntry…\|PASS: TestAC14GoSideEvalPush…"` 我自己跑＝**`:201`＋`:205`**，另两发同形 **`:241`/`:245`**、**`:277`/`:281`**（件＝`docs/evidence/s1/33-panel-host-c27-r7.md`，身份行 `:10` 逐字 `HEAD：416d9d56977a89a070e29e23391b1a2f02bdc9b5`、`:23` 逐字 `发 1 已填（rc=0／240-0-0／红册空）`；`git log -1 416d9d56` ＝ `Thu Oct 1 15:39:23 2026 +0800`）⇒ **成立，这是本程最重的一枚** |
| ② 嫌疑笔 | `70b00885`＝纯搬函数、⛔ 改任何交接 | `git show 70b00885 -- cmd/wisp/panel_host_windows.go` 我自己读 diff＝`bringUp` 尾段那四行（`firstRoundTripLocked` → `serveEntry`/`serveNotBuiltNoticeLocked`）**逐字原样搬进新函数 `coldStartPageHandover`**，同语句同顺序，＋41/−11 里除注释⛔ 其它语义 ⇒ **否证那枚嫌疑**（腿自己把它当初〔仅腿报〕处理，做对了） |
| ③ push 缺两截 | 库侧 `ExecuteScript(script, 0)`⛔ 回程＋与 AC13 共享"最后一份文档⛔ 落地" | 那枚文件**⛔ 在仓里**：路径＝module cache `D:/work/base/gopath/pkg/mod/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808/pkg/edge/chromium.go:144`，我现读到 `e.webview.vtbl.ExecuteScript.Call(…, _script, 0,)` ⇒ 第三枚参数（回调）传 **0**＝无回程 ✓；`go.mod:8` 逐字 `github.com/jchv/go-webview2 v0.0.0-…`，仓里**⛔ `vendor/` 目录** ⇒ 要回程＝**改依赖**（fork／`replace`），⛔ 本票射程（见第 3 节） |
| ④ 生产调用者 | `EVAL_PROD_HITS`＝0 | `grep -rn --include=*.go '\.Eval(' cmd internal \| grep -v _test.go` ＝ **0 行** ✓ ⇒ 修②今天⛔ 生产调用者受益（受益面＝测试判据本身） |
| 收益边界 | 入口页的 `./assets/*.js|css`⛔ 有注册方 ⇒ AC13 绿⛔ 等于面板有内容 | `frontend/dist/index.html:19-20` 逐字 `<script type="module" crossorigin src="./assets/index-B8yINMF1.js">` 与 `<link … href="./assets/index-yy8KMgdf.css">`，`:23` 是 `<div id="root"></div>`；仓里 `AddWebResourceRequestedFilter`/`SetVirtualHostName…` 的**生产调用点＝0**（只有 `panel_host_windows.go:42` 与 `internal/panel/doc.go:3` 两行注释提到它）；`NavigateToString` 形同 `about:` 文档 ⇒ 相对资源取不到 ✓ **同一枚缺口早就有票＝票 299（`299-panel-page-cannot-fetch-its-own-assets-no-webresourcerequested-registration.md`，5 格未勾，停在 `Q-84`）⇒ ⛔ 新票** |

⚠ **腿自己两处引用对不上盘**（记档，⛔ 改它的件）：① 它 `90-final.md` §D.1 把 `nail2` 那枚 PASS 写成 `33-panel-host-c27-r7.md:207`，**盘上是 `:205`**（它自己 `logs/l4` 里写的就是 `:205`＝最终件里手滑）；② 它报的其余 28 枚锚我自己抽查（`:467`/`:486`/`:892`/`:19-20`/`chromium.go:144`/`go.mod:8`）逐枚命中。★定式＝**同一枚事实在同一腿的两个件里给了两个行号**，下游一律以"我自己 `grep -n` 那一把"为准。

### 2. ★裁：`AC#2` ⛔ 降级（票面那句自动降级条件⛔ 触发）

票面 `AC#2` 写的是"若 `AC#0` 判出'这一维从来没有绿过 ⇒ 是缺功能⛔ 是回归'，本格自动降级为⛔ 开工"。现量把它**否掉**了：`416d9d56`（10-01）那一发里 `nail2` 与 `AC13` **两枚都 PASS**（`:205`/`:201`，整包 240-0-0），而 `f718e9b6`（10-10，`fb2fb802` 的父发）同一枚干净 clone 上 `nail2` 已是那句 `title=""` ⇒ **症状② 是一枚可回归、有绿点的缺陷，窗口＝`416d9d56`（绿）↔ `f718e9b6`（红）**。⇒ `AC#2` **保持开工射程**、⛔ 降级。
⚠ 三条边界同时落下：
1. **落地腿第一发＝在这个窗口里定位**（`git log 416d9d56..f718e9b6 -- cmd/wisp internal/panel` 逐枚判），⛔ 上来就新造机制；嫌疑笔 `70b00885` 已由我与腿各自否证（纯搬函数），⛔ 再当结论引。
2. **⛔ 动依赖面**：`Eval` 无回程那截在**第三方模块**里（`go-webview2`，无 `vendor/`）⇒ 要它带回程＝fork 或 `replace`＝**换依赖＝人工批准**，⛔ 本票自决。
3. **⛔ 动泵形状**：腿的必答我认＝"让页面先自报就绪再推"那一形⛔ 撞 `A502` 的 P1（`Eval` 本来就在泵线程上跑，甲形⛔ 是答案）；"要回程"那一形**同时**撞依赖面（上面第 2 条）与泵面（`A798` 甲形逐字"泵形状要重来一遍"），且与〔待人拍板〕`Q-84` 同格 ⇒ 若定位结论是"必须补回程才能修"，**回来找我裁、由我摆给 owner**，⛔ 落地腿自己绕。

### 3. ★裁：两枚症状＝**部分同源**，⛔ 合并、⛔ 各修

腿给的判死料我复核成立：`nail2` 的红句在四发里**一句没变**，而那四发的"最后一份文档"在两轴上是**不同的**（干净 clone＝`:467` 公告 181 B；母仓＝`:486` 入口 1044 B）⇒ `nail2` ⛔ 依赖"入口交接成功"。同时两枚**共享**上游那一截（最后一次 `SetHtml` 的那份文档⛔ 是 `Eval` 落到的那份）。⇒ 本格成论＝**修①⛔ 自动闭②**；票面 `AC#3` 那枚**跨维度负控**（"修①的那笔若让 `nail2` 变绿必须具名解释为什么"）从此是**判死那一发的凭据**，落地腿⛔ 用推理替代它。
⇒ `AC#1`／`AC#2` **两格都保持**，⛔ 把两格并成一格、⛔ 因"共享上游"只修一截。

### 4. ★更正我自己两处（⛔ 删原句，逐处具名）

1. **本票标题**那句"`Eval` 推不进文档那一枚**早在回归之前就红、从来没绿过**"——后半句**⛔ 成立**（见第 2 节 `:205`）。前半句（⛔ 由 `fb2fb802` 造成）照旧成立。标题原句⛔ 改，以本节为准；下次任何文案引"从来没绿过"＝引这一条更正。
2. **票 303 的登记**（`AC#3` 下面那条 22:2x blockquote）里我写的"…且属还没做出来的功能"——同样**⛔ 成立**，已在票 303 就地追加更正。⚠ 那枚登记的**结论⛔ 动摇**：`AC#3` 那半句之所以⛔ 能在票 303 里满足，靠的是**归口**（那两枚症状属票 305），⛔ 靠"做不到"。
3. **我派单里那句"票 33 的 `A475` 那批裁过'两维⛔ 混一枚'"**——**⛔ 对不上盘**（记我）：台账 `:9890` 那枚 `A475` 是**票 228-r1／球第一次住进会跑任务的那条腿**，⛔ 面板两维；"两维⛔ 混一枚"的真身＝**10-01 那道 P2／`A502`**（台账 `:10287`）。⇒ 以后派单引这格一律写 `A502 P2`。⚠ 腿自己把它列成冲突第 3 条并⛔ 改我票面一字，处置正确。

### 5. ★新规矩一枚（这条是腿现量出来的，我认并落档）

票面"台面规则"原来只写了**一根轴**（母仓带旧 `frontend/dist` ↔ 干净 clone ⛔）。盘上还有**第二轴**＝**同一枚母仓里那枚未跟踪 dist 自己会漂**（现量：`wc -c`＝1044、mtime **2026-10-10 08:51**，⛔ 在 HEAD）。⇒ **"母仓↔母仓"的跨时间对照⛔ 天然同台面**；凡引母仓的绿↔红对照，必须**同时**交出 `wc -c frontend/dist/index.html` 与它的 mtime，否则那两发⛔ 可比。已写进票面"规矩"段的新一条（⛔ 改原句）。

### 6. 欠着的读数（⛔ 算腿的欠账，具名归我）

腿响亮报了四件⛔ 做到的，全部合规、且**三件是我禁的射程**（⛔ 编译面、⛔ 真窗、⛔ winlive 档）：
- **R2**＝同一枚台面的绿↔红对照那一发；**R1/R4**＝`Eval` 那一跳的机制三形（脚本被丢／落到另一份文档实例／文档在两发间被重建）静态读码分⛔ 开；**R3**＝10-01 那份 dist 字节⛔ 可复现（未跟踪、已漂）＝腿具名报"做不到"，我认；
- **winlive 档那一发**（凭据逐字 `panel_host_windows_live_test.go:18` 的 `winlive has NO CI job`）。
⇒ 这四件**都只能由我在 `cmd/wisp` 车道上跑**，而该车道今天⛔ 空：写腿 `306-r1` 正在 `internal/audio` 里做三对突变（`go list -deps ./cmd/wisp/` 已证 audio 在 `cmd/wisp` 的导入图里）⇒ **本程⛔ 跑、⛔ 跑 `d22scan`**（它扫 `internal/`，那枚临时突变会让一把 docs-only 的程吃假红），全部等 `306-r1` 交完再取。⚠ 这不是腿的欠账，是先例"派单里就禁掉的资源⛔ 算腿的欠账，记编排者"。

### 7. 现态与 next

- **现态**＝票 305 **1 勾（`AC#0`）／5 未勾**（`AC#1`／`AC#2`／`AC#3`／`AC#4`／`AC#5`），⛔ 改 `-done`。票 306 **1 勾／5 未勾**（写腿在飞）；票 303 **6 勾／1 未勾**＋登记＋本程第 4 节那条更正；票 300 **6 勾／1 未勾**（`AC#4` 那半按住）。
- **在飞**＝只剩 `306-r1`（它已落三笔：起手锚、`AC#1` 夹具、`AC#3` 注释；还欠 `AC#2`/`AC#2b` 的三对突变凭据与 `AC#4` 门禁）。交件后我自己复跑再派 `306-v1`。
- **next**＝收 `306-r1` → 派 `306-v1` → **补票 300 `AC#4` 那半**（此时 `cmd/wisp` 面才空）→ **我在 `cmd/wisp` 车道取票 305 第 6 节那四发**（R1/R2/R4＋winlive）→ 派 `305-r1`（按第 2、3 节的三条边界）→ 票 303 改名 `-done`＋索引 → `296-r2` → `295-r1` → `294-r1` → `294-v1` → `298-v1`。⛔ 动依赖面、⛔ 动泵形状、⛔ 动 `Q-84`（未答＝默认不做）、⛔ 动 D32 那枚 `1500`。
- Progress log：
  - [2026-10-11 07:3x +0800] agent=编排者 did=收 `305-a1`（三笔 `5f891575`→`52a5eb23`，逐笔名册我自己重跑＝**越出 `probes/305/a1/**` 0 格**、新 `.go` 0、`.out` 0、⛔ push、⛔ 编译面只 `go list -m`）＋**我自己现量承重读数**：三份 `SetHtml` 逐字复现（`:467`/`:486`/`:892`），两枚字面量我自己拼接现算＝公告 **181 B**／探针 **136 B** 逐字相等，入口 **1044 B** 量的是**未跟踪本机产物**（`git cat-file -s HEAD:frontend/dist/index.html`＝`exists on disk, but not in 'HEAD'`，mtime 2026-10-10 08:51）；`grep '\.Navigate(' cmd internal` 非测试面 **0 行** ⇒ 活文档只由那三枚 `SetHtml` 决定；`grep '\.Eval(' cmd internal` 非测试面 **0 行**；★**最关键那枚我 grep 出来了**＝`33-panel-host-c27-r7.md:201`（AC13 **PASS 0.97s**）与 **`:205`**（nail2 **PASS 0.50s**），身份 `:10`＝HEAD `416d9d56`、`:23`＝整包 `rc=0／240-0-0／红册空`（`git log -1 416d9d56`＝10-01 15:39），另两发同形 `:241/:245`、`:277/:281`；`git show 70b00885 -- cmd/wisp/panel_host_windows.go` 我自己读 diff＝尾段四行逐字搬进 `coldStartPageHandover`、同语句同顺序 ⇒ **嫌疑笔否证**；库侧 `ExecuteScript(…, _script, 0,)` 现读于 module cache `…go-webview2@v0.0.0-20260205173254-56598839c808/pkg/edge/chromium.go:144`（`go.mod:8` 逐字命中、仓⛔ `vendor/`）⇒ 要回程＝改依赖；入口页 `:19/:20` 那两枚相对 `./assets/*.js|css` 的注册方＝**生产调用点 0**（只剩两行注释）⇒ 缺口⛔ 新票、**归口＝票 299（5 格未勾，停 `Q-84`）**。⇒ **`AC#0` 翻**；★**裁 `AC#2` ⛔ 降级**（症状②＝可回归、绿点 `416d9d56`、窗口 ↔ `f718e9b6`），三条边界＝先 bisect／⛔ 动依赖面／⛔ 动泵形状（甲形⛔ 撞 `A502` P1，乙形撞依赖＋泵＋`Q-84` 同格 ⇒ 回来找我）；★**裁两枚症状＝部分同源、⛔ 合并**（nail2 四发一句没变而"最后文档"在两轴上不同 ⇒ 修①⛔ 自动闭②，判死凭据＝`AC#3` 那枚跨维度负控）；★**更正我自己三处**＝本票标题"从来没绿过"⛔ 成立／票 303 登记里"属还没做出来的功能"⛔ 成立（就地追加、结论⛔ 动摇，靠的是归口⛔ "做不到"）／派单把 `A475` 当"两维⛔ 混一枚"（真身＝**`A502` 的 P2**，台账 `:9890`⇄`:10287`）；★**新规矩**＝母仓↔母仓 跨时间对照⛔ 天然同台面，必须同时交 `wc -c frontend/dist/index.html`＋mtime（腿现量到的第二轴）；⚠ 腿的最终件把 nail2 的 PASS 写成 `:207`、盘上 **`:205`**（它自己 `logs/l4` 写对过）＝同腿两件两个行号，下游以我自己 `grep -n` 那把为准；四件欠读数（R1／R2／R4＋winlive 档）⛔ 算腿的欠账、具名归我，且⛔ 本程跑（`cmd/wisp` 导入图含 `internal/audio`、写腿 `306-r1` 正在做三对突变 ⇒ 那一发会取到被临时改过的产码；同因本程⛔ 跑 `d22scan`，docs-only 的程会吃假红）。现态＝票 305 **1 勾／5 未勾**，⛔ `-done`；在飞只剩 `306-r1`；next＝收 `306-r1`→`306-v1`→补票 300 `AC#4`→我在 `cmd/wisp` 面取那四发→派 `305-r1`→票 303 改名 `-done`＋索引→`296-r2`→`295-r1`→`294-r1`→`294-v1`→`298-v1`；本程零产码、⛔ push（`306-r1` 在飞，同一把理由）
  - [2026-10-11 08:0x +0800] agent=305-a2（只读普查腿，第二程） did=交齐三张表＋配角件，件 `probes/305/a2/**`（三笔 `a3560d6b`→`05b37976`→本笔）。★本腿零产码、零新 `.go`、零 `.out`、⛔ push、⛔ 任何 Go 编译面（`go build`／`go vet`／`go test`／`d22scan`／`gofumpt` 全 0 条，**`go env`／`go list` 亦 0 条**）、⛔ 真窗；越界尺逐笔 `git show --name-only --format=`＝越出自家目录 **0 格**（`logs/41-overreach-per-commit.txt`）。⇒ **表①**：窗口 `416d9d56..f718e9b6`＝963 枚且**⛔ merge**（`git rev-list --count --merges`＝0、`--first-parent --count`＝963 ⇒ 线性，逐枚 detach 安全）；名册尺 `git log --oneline 416d9d56..f718e9b6 -- cmd/wisp internal/panel`＝**37 枚**，真碰文档链的只有 **4 枚**（★定罪尺＝`git grep -l -E 'SetHtml|wispDispatch|currentWindow|bringUp|firstRoundTrip|serveEntry' 416d9d56 -- cmd internal`＝5 枚文件**全在 `cmd/wisp`**）＝`f7d28ef0`／`7a0236b3`／`4658dbb6`／`cf95c799`；另三枚列"可能不能"＝`0d87a681`（入向门白名单扩容）／`480b970d`（`ip/pump.go`＝快照泵 `SnapshotPump`，读码判它无关、脱罪链已写进表里供跑打）／`32e74479`（`panel_pump_test.go` 种子与 ctx 时长）；排除 **30 枚**（＝那张排除名册的行数，⛔ 等于"已排除"，它只排除"产码面直接改这两条路径"，⛔ 排除整包序/线程污染那族）。头号嫌疑（各带读到的 diff 内容锚）＝**症状① `f7d28ef0`**（`bringUp` 在建窗前新增第二道闸：`+ if dirty, closedHwnd := staleCloseQueued(); dirty {` → `+ return fmt.Errorf("panel host: refusing to create the panel window: …", closedHwnd)` ⇒ 命中即第 5/6 步那两次 `SetHtml` 永远跑不到）／**症状② `7a0236b3`**（`var errPanelRefusedThread`＋`rp.post(func() { rp.showOnThread(via) })`，`showOnThread` 认哨兵即**判死常驻面板线程**；同一笔在 AC13（`:313`）与 nail2（`:855`）之间插进 179 行会开真窗的新用例）＋形状不同的第三枚 **`4658dbb6`**（`bringUp` 的 `WindowOptions{Title: panelTitle, Width: 420, Height: 260}` 整块换成 `m.windowOptions()`，体内每次建窗现调 `panelGeometrySource`→`config.LoadFile(cfgPath, nil)`＝STA 线程建窗那一跳里多了盘读；笔自标"半成品·未验证"）。⇒ **表②**：生产面 `SetHtml` 只有三枚调用点（`:467` 公告／`:486` 入口／`:892` 探针），装配序＝先探针后入口、入口报错才公告 ⇒ **按语句最后一次⛔ 是探针**；生产面 `.Navigate(`/`NavigateToString`＝**0 行**、`.Eval(`＝**0 行**（两把我各自独立复跑）；链内唯一手工泵点＝`panel_host_windows.go:897` 的 `pnlPumpOnce()`，而 `post()` 建窗后一律走 `w.Dispatch`（该文件头注释逐字"Run() is that queue's reader"）、首扇窗 `bringUp` 跑在 `drainTasks()` 里 ⇒〔**读码推到、⛔ 判死**〕最后一次 `SetHtml` 之后这条链⛔ 泵。★**新量到一格（能省一发）**：三份文档指纹互异——探针 136 B ⛔ `<title>` ⇒ `document.title==""`；公告 181 B ⇒ `"panel assets unavailable"`；入口 1044 B ⇒ `"Wisp"` 且唯一一枚 `id="root"`（`entryIDProbes` 因此只交回 1 枚，与红句"NONE of the 1 element ids"对上）——而 nail2 红句逐字 `the page reports its title as ""` 我在 `probes/303/orch/r3-f718e9b6-three-nails.txt:17` 现读到 ⇒ **症状②自己已答出"那一刻在场的是探针"**；AC13 那位图判据分⛔ 开探针与公告 ⇒ **修①若只让 `serveEntry` 被执行而⛔ 让它落地，nail2 照旧红**（＝跨维度负控那一发的读法写死）。⇒ **表③**：9 点二分名册（7 枚候选＋绿红两端点，生成尺逐字给在件里）＋发 0 台面自检／发 1 绿端点复现（＝R2）／发 2–4 二分带**单调性闸门**（出现非单调即改逐枚全跑 7 发、⛔ 报"二分定位到 X"）／发 5 AC13 与 nail2 同枚对照／发 6 负控前置读数；每发必交 `sha`＋clone 路径＋`-run` 串逐字＋`rc`＋`--- PASS/FAIL/SKIP` 逐字红句＋`wc -c`/`sha256sum`/mtime 三枚。★**台面第三轴一枚（票面原有两根，盘上还有第三根；建议进票面"规矩"段，我⛔ 改票面一字）**：`git ls-files third_party/sherpa-onnx`＝**0 行**、`git ls-files build`＝**0 行** ⇒ 干净 clone **⛔ 带这两枚未跟踪 DLL 目录**，照抄 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 进 clone＝两枚路径⛔ 存在＝`0xc0000135` 且**零 `--- FAIL`**（＝用例根本没跑那形假绿）⇒ 配方要求 PATH 两枚条目写**母仓绝对 `/d/…` 路径**（或把两枚目录复制进 clone），并由发 0 先断言具名用例真发出了 `--- PASS/FAIL/SKIP` 行。⛔ 做到的**逐条具名**（⛔ 写成"没做"、⛔ 用静态读码冒充色）：① **R2** 同台面绿↔红对照那一发；② **R1／R4** `Eval` 那一跳三形分离（我只用内容排除了"落到另一枚控件"那一支，另三形＋我新添的第四形都要跑）；③ **R3** 10-01 那枚 dist 字节复现（我只交出今日三枚读数 1044 B／`sha256 9b7856b9…074f`／mtime 2026-10-10 08:51:10.5，并具名说"发 1 用今日字节跑出来的 AC13 绿⛔ 等于 R3"）；④ 整包序那一族（`cf95c799` 新增 `settleTreeReading` 的 3 s 采样等待、`32e74479` 的 ctx 30s→45s、`7a0236b3` 那 179 行新用例）的颜色——⛔ cmd/wisp 车道⛔ 能量。与任务书不符处（以盘上为准）：起手 HEAD＝**`e4740e35`**（⛔ `6284a489`）；`70b00885` **根本⛔ 在这 37 枚名册里**（S3 name-only 尺，⛔ 需要再判一次，"已否证"那格照旧）；绿点六枚锚我自己重跑过且**与任务书逐字对上**（`33-panel-host-c27-r7.md:201` AC13 PASS 0.97s／`:205` nail2 PASS 0.50s，另 `:241`/`:245`、`:277`/`:281`，身份行 `:10`＝HEAD `416d9d56…`）；票面"红句在四发里一句没变"我**只逐字重读了 1 发**（`r3-f718e9b6`），另三发引自票面、⛔ 在我手里。现态＝票 305 仍 **1 勾／5 未勾**（`AC#0` 那格由编排者翻），本腿⛔ 翻任何框、⛔ 改原句、⛔ `-done`；依赖面（`ExecuteScript(…, _script, 0,)` 在第三方模块、仓⛔ `vendor/`）与泵形状（专用 STA＋库 `Run()` 泵）与 `Q-84` 三格一律⛔ 动、⛔ 提议。

## 编排者收 `305-a2`（2026-10-11 08:0x，三笔 `a3560d6b`→`05b37976`→`17dd4a6f`，件 `probes/305/a2/**`＝五份正文＋`logs/` 38 枚；⛔ 本节目＝我自己复跑的读数）⇒ `AC#0` 那四张表⛔ 重开；★**台面⛔ 只有两根轴，第⏂根本轮量出来了**；判死那一发的配方我照它跑

### 1. 承重读数我自己现跑（⛔ 抄腿；每条带尺）

| 腿的读数 | 我这发 |
|---|---|
| 窗口名册＝**37 枚**（总 963 里碰 `cmd/wisp`／`internal/panel` 的） | `git log --oneline 416d9d56..f718e9b6 -- cmd/wisp internal/panel \| wc -l`＝**37** ✓ |
| 四枚真碰文档链的嫌疑 | `git log -1` 逐枚存在：`f7d28ef0`（10-01 16:20）／`7a0236b3`（10-01 18:04）／`4658dbb6`（10-03 13:58）／`cf95c799`（10-03 17:06）✓ |
| ★三份文档指纹互异⇒`nail2` 红句**自己点名在场的是探针**（136 B 那枚⛔ 带 `<title>` ⇒ `title==""`） | 我现跑＝探针那段（`cmd/wisp/panel_host_windows.go:888-900`）`grep -ci '<title'`＝**0**；`frontend/dist/index.html`＝**1** ⇒ **成**——这一枚值钱：它把"在场的是哪份文档"从**跑出来的**变成**红句里就写着的** |
| `70b00885` **根本不在**这 37 枚名册里 | 尺＝`git log --format=%H 416d9d56..f718e9b6 --name-only -- cmd/wisp internal/panel \| grep -c '^70b00885'`＝**0** ⇒ ⛔ 任何人再把它当本票嫌疑引（`AC#0` 那格已闭，见"规矩"段那条已闭注） |
| 它⛔ 跑任何 Go 命令、⛔ 动工作树 | 逐笔 `git show --name-only --format=` 我自己重跑＝名册全在 `probes/305/a2/**` 下；新建 `.go`＝**0**；`.out`＝**0**；0 字节件＝**0**（`-type f` 那把尺）；它引用名册时把别家的两笔（我的 `b55a3ff6`、`risk-attrib-1` 的 `f4fd71bc`）如实标成"别家交错"、⛔ 碰 ⇒ 越界 0 格 |

### 2. ★新规矩一枚（腿现量、我复跑）＝**台面有第三根轴**

ⓐ母仓↔干净 clone（已写）、ⓑ母仓那枚未跟踪 `frontend/dist/index.html` 会漂（已写）、**⏂干净 clone／`git archive` 导出树⛔ 带 `third_party/sherpa-onnx` 与 `build` 这两枚 DLL 目录**——尺＝`git ls-files third_party/sherpa-onnx`＝**0**、`git ls-files build`＝**0**（我 08:0x 现跑复认）。⇒ 把 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 照抄进 clone ⇒ 解析成 clone 里⛔ 存在的路径 ⇒ `0xc0000135` **且零 `--- FAIL`**＝"用例根本没跑"那形假绿。⇒ 两条落地规矩（已追加进上面"规矩"段一条）：跑 `cmd/wisp` 时 PATH 写**母仓绝对** `/d/…` 形；**第一发必须先断言具名用例真发出了 `--- PASS/FAIL/SKIP`**（⛔ 拿 `rc` 或行数当"跑过了"）。

### 3. 判死那一发＝**归我跑**（配方射程⛔ 在此重复，正文＝`probes/305/a2/30-deciding-run-recipe.md`）

- 序列 9 点＝`416d9d56`(G) → `f7d28ef0` → `7a0236b3` → `0d87a681` → `480b970d` → `4658dbb6` → `cf95c799` → `32e74479` → `f718e9b6`(R)；窗口线性（`--first-parent --count`＝963＝总枚数 ⇒ 逐枚 detach 安全）。
- ★**非单调即作废**：腿自带一道闸门＝二分结果若⛔ 单调，改成逐枚全跑，⛔ 报"二分定位到 X"。这条我照守。
- 我这边的顺序＝发 0 台面自检（先断言用例真发出颜色）→ **发 1＝绿端点复现（＝我欠的那枚 R2，同一台面）** → 发 2–4 `nail2` 二分 → 发 5 `AC13` 同枚对照（母仓台面，⛔ 交 `wc -c frontend/dist/index.html`＋mtime＝第ⓑ轴）→ 发 6 跨维度负控前置读数。
- 头号嫌疑（腿给的排序，⛔ 定罪、只算"先验"）＝症状①`f7d28ef0`（建窗前多出一道会整枚 `return` 的闸 ⇒ 两次 `SetHtml` 可以永远跑不到）；症状②`7a0236b3`（哨兵 `errPanelRefusedThread`＋`showOnThread` 认哨兵即判死常驻面板线程；同笔在 `AC13` 与 `nail2` 之间插进 179 行开真窗的新用例）；第三形`4658dbb6`（`WindowOptions` 整块换成每次建窗现调 `panelGeometrySource`→`config.LoadFile`＝盘读进 STA 建窗那一跳）。⚠ 三枚**全是 10-01～10-03 的面板宿主／尺寸那一族**，与票 33／票 255 同族 ⇒ 定位到之后**归哪张票修**要另裁（⛔ 默认塞进 305）。

### 4. 它响亮⛔ 做到的五件 ⇒ ⛔ 算腿的欠账，**具名归我**（先例＝"派单里就禁掉的资源⛔ 算腿的欠账"）

①任何颜色读数（本腿零测试，全部标〔读码推到〕，⛔ 有一格写成"已判死"——这一条它做得对）；②R2 同台面绿↔红那一发；③R1／R4 `Eval` 那一跳的机制三形＋它新添的第四形（它只内容排除了"落到另一枚控件"那一支）；④R3＝10-01 那份 dist 字节⛔ 可复现（未跟踪、已漂）；⑤整包序那一族的颜色（`cf95c799`／`32e74479`／`7a0236b3` 那 179 行）。⇒ ①–⑤ 全在**我的 `cmd/wisp` 车道**上，⛔ 第二人碰过。

### 5. 现态与 next

票 305 现态⛔ 变＝**1 勾（`AC#0`）／5 未勾**（`AC#1`／`AC#2`／`AC#3`／`AC#4`／`AC#5`），⛔ `-done`；腿⛔ 翻框（改前改后各量一次＝`logs/43-append-verify.txt`）。next＝**我跑那 6 发**（上面第 3 节）→ 按结果裁落点归属 → 派 `305-r1`（三条边界照 `AC#0` 收腿那一节：先 bisect／⛔ 动依赖面／⛔ 动泵形状）。同批在飞补位＝`306-v1`（票 306 五格终裁）、`307-a1`（票 307 `AC#0`）、`303-v2`（补那张与 AC 编号 1:1 的七行裁决表）。
- Progress log：
  - [2026-10-11 08:0x +0800] agent=编排者 did=收 `305-a2`（三笔，逐笔越界我自己重跑＝越出 `probes/305/a2/**` **0 格**、新 `.go` 0、`.out` 0、0 字节件 0、⛔ 编译面只读）＋我自己复跑它的承重尺＝窗口 **37** 枚✓／四枚嫌疑存在✓／★探针面 `<title>` 命中 **0**、入口 **1** ⇒ "`nail2` 红句自己点名在场的是探针"**成**／`70b00885` 在 37 枚名册里 **0** 命中 ⇒ ⛔ 再引；★**台面第⏂根轴**（`git ls-files third_party/sherpa-onnx`＝0、`build`＝0 ⇒ clone／导出树⛔ 带 DLL ⇒ 照抄 `$PWD` PATH＝`0xc0000135` 零 `--- FAIL`＝假绿形）已复跑并追加进"规矩"段⇒ 落地两条：PATH 写母仓绝对 `/d/…`、第一发必须先断言具名用例真发出 `--- PASS/FAIL/SKIP`；判死那一发按它的 9 点序列**归我跑**（发 0 自检→发 1 绿端点＝我欠的 R2→发 2–4 `nail2` 二分→发 5 `AC13` 母仓对照带 dist 字节＋mtime→发 6 跨维度负控），★非单调即改逐枚全跑、⛔ 报"二分定位到 X"；它响亮⛔ 做到的五件⛔ 算腿欠账、具名归我。现态＝票 305 **1 勾／5 未勾**⛔ `-done`；在飞＝`306-v1`；next＝我跑那六发→裁落点归属→`305-r1`；⛔ push（照 `A830` §6 那两条理由）
