# 35 — PanelBridge (C17): method whitelist, correlationId routing, resync statelessness

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 33-panel-host-c27, 34-frontend-scaffold
**Parallel slots:** ≤2 sub-agents (A: bridge transport + dispatch + whitelist; B: resync model +
event stream + frontend client SDK)
**Spec refs:** SPEC-08 §5.2, C17, D31 routing, F2 server-side allow rule, D38(d) backpressure

## What to build
The frontend↔Go bidirectional bridge: `invoke(method,args)→result` with the frozen method
whitelist, correlationId request routing, Go→frontend event push (task deltas, tool chips,
approval requests, ball state, cost ticks) with bounded-queue merge, and `panel.resync` full-
state push making the frontend provably stateless.

## Key constraints
- Method whitelist per SPEC-08 §5.2 table (panel.resync, tasks.list/detail, history.query,
  transcript.get, approval.current/queue, approval.decide, config.get/set, grants.list/revoke,
  privacy.purge/export, cost.summary, models.list/delete, diagnostics.export). Unlisted method →
  reject + log (test with fuzzed names).
- **`approval.decide`: `allow` REJECTED for any panel-sourced call — enforced server-side
  (Go) regardless of UI** (F2 layer 3); `reject` allowed from panel. Method metadata table
  carries required capability + native-authorization flag (consumed by dispatcher).
- correlationId routing: concurrent invokes + out-of-order replies land on the right caller;
  10-way concurrency test.
- Event push: bounded queue per panel; overflow MERGES increments (never drops content);
  backpressure counters visible.
- `panel.resync` on EVERY show: Go pushes full state (tasks, approvals, drafts of current
  results, config view, costs); frontend asserts "no cached state across show" — resync-only
  hydration pattern encoded in the client SDK; stale-state test (hide → mutate → show shows
  fresh).
- Frontend client SDK: typed invoke/events, auto-reconnect on hide/show, error toasts via
  Sonner with badge-fallback rule comment (D42#8).
- Payload hygiene: no secrets/keys in any bridge payload (scan test); long strings truncated
  per log rules.

## Out of scope
- Page implementations (36–40); native card UI (21 already built).

## Acceptance criteria
- [ ] Whitelist fuzz: 50 random/forbidden method names → all rejected + logged; each listed
      method dispatches.
- [ ] Concurrency routing: 10 parallel invokes with shuffled replies → correct pairing.
- [ ] Forged panel allow: direct bridge call `approval.decide{allow:true}` from the webview
      context → server rejects (decision test at Go boundary + browser-context e2e).
- [ ] Resync: hide → tool state changes → show → UI reflects fresh state; SDK forbids cross-show
      caches (architecture test/lint).
- [ ] Backpressure: flood events under blocked consumer → merges, no unbounded memory (heap cap
      asserted), content integrity kept.
- [ ] No secret leakage scan across bridge payloads.
- [ ] **Transport agreement page↔host (added 10-07 by `35-a1`, shape NOT chosen here)**: the page posts its envelope on
      `window.chrome.webview.postMessage` (`frontend/src/lib/panel.ts:141` interface, `:179`/`:218` the two senders), while Go's
      handler is reachable only through the library's own RPC pipe — `w.Bind("wispDispatch", …)` (`cmd/wisp/panel_host_windows.go:80`/`:405`)
      injects `window.wispDispatch(…)` which calls `window.external.invoke(JSON.stringify({id, method, params}))`, and every inbound
      string goes through the library's `msgcb` (go-webview2 `webview.go:140-158`, wired at `:103`). The page never calls
      `wispDispatch` (`git grep wispDispatch HEAD -- frontend` = **0**), and its own envelope, when it reaches `msgcb`, unmarshals
      as an RPC message with `id=0` and an unknown method ⇒ the library replies by evaluating `window._rpc[0].reject(…)` into a page
      that has no such slot. **The request dies with no Go-side log and no page-side receipt.**
      Any fix must be Go-side transport (page files are out of bounds for this fleet), must keep C17's roster unchanged,
      ⛔ and no leg may "make this green" by weakening the whitelist or by answering the RPC parser without the envelope ever
      reaching `dispatchRaw` (`panel_host_windows.go:630`).

## Progress log (append-only, newest last)

### 10-07 编排者现状对拉（票面 09-19 写的东西从没按今天的仓复算过；落账 `A657`）

★**为什么这节写在本票而不是新票**：机主原话「**不要重复劳动哈**」。本票就是 C17 那座桥（`Go→frontend event push` ＋ `panel.resync` 全量推），而字段对接普查（`145-f1`）量出来的头号阻断**正好落在本票的许诺上**——所以我把现量补到本票面上，⛔ 不再另开一枚同义票。

**1. 方向要分两半，今天只有半座桥**（四把尺都我自己现跑，⛔ 不引腿的读数）：
- **页→Go（入向）＝真在**：`cmd/wisp/panel_host_windows.go:405` `w.Bind(panelDispatchBinding, func(raw string) string{…})` → `:630 dispatchRaw` → `internal/panel/composer_dispatch.go:181 … HandleModeRequest(…)`。方法名册现量 **6 枚**（`internal/panel/bridge.go:42-49`：`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send` ＋ 票 248 的 `config.get`／`config.set`）。
- **Go→页（出向）＝整条不存在**：ⓐ`grep -rn --include=*.go -E 'PostWebMessage|EvaluateScript|CreateWebMessageAsJson' cmd/wisp internal/panel`（剔 test）＝**0**，正控＝同尺全仓只命中 `.scratch/wisp/probes/33/p1/q2/main.go` **一枚探针**；ⓑ泵 marshal 出的字节唯一去处＝`cmd/wisp/run.go:725 Out: rt.bookPanelSnapshot` → `panel_pump.go:321` 只写 `rt.lastSnap/lastSnapBytes` ＋ 一行日志，而读者 `(rt *agentRuntime).lastPanelSnapshot()` 带括号数＝**10 处调用全在 `_test.go`、非 test 0 处**；ⓒ`func (m *PanelManager)` 的导出方法 **8 枚**（`IsCreated/IsShown/LastColdMs/LastHotMs/Show/HotShow/Hide/Destroy`）**没有一枚能把数据送进页面**。
⇒ 票面 §"What to build" 里那句 "Go→frontend event push (task deltas, tool chips, approval requests, ball state, cost ticks) with bounded-queue merge" **今天兑现枚数＝0**，而它正是 `panel.resync` 那一格的前置。

**2. `panel.resync` 两侧都是空的（而且我差点被一把坏尺骗了，具名记我）**：尺＝`grep -rni --include=*.go resync cmd internal` 回 **4 枚命中**——**逐行读全文后发现全是同串假阳性**（`internal/risk/provenance_test.go:45/51/62` 的 `fixtureSy`**`reSync`**`Root`、`syncdirs_test.go:155` 的 `NutstoreSy`**`reSync`**）。⇒ **真命中＝0**。页面侧同族尺（`git grep -iE 'resync' dsh/feat/frontend-p0-v2 -- 'frontend/src'`）只有 **1 处注释**：`lib/panel-views.ts:90` 逐字 "…ticket 35's pump owns it, and `panel.resync` (the push…" ⇒ **它在等本票，不是等页面**。⚠ 教训：大小写不敏感的**子串**尺会把驼峰里的同串当命中；判"某个词在不在仓里"要用词边界形（`\bresync\b`）并**逐行读命中全文**。

**3. 页面没有耳朵**：`git grep -E 'webview\.addEventListener|postMessage|onmessage' … -- 'frontend/src'` ＝ **0**；正控＝同树 `addEventListener` 在 8 枚文件命中（`App.tsx` 2 枚）。页面唯一的桥引用是 `lib/panel.ts:154 return window.wispBridge ?? window.chrome?.webview ?? null`，**只用于发起调用**。⇒ **本票的落地集不是单侧活**：出向那一跳即使今天在 `cmd/wisp` 里接上了，页面没有接收器仍然画不出来；而"谁来写页面文件"⛔ **不在我今天的许可里**（机主 10-07 放开的是 `frontend/**` **只读**，写页面他没给、我也没替他扩大解释）。⇒ 这一格是**真待拍板项**，我按第 12 款把"不做"那一支也写出来：ⓐ我把出向跳＋接收器形状**只在证据件里给逐键清单**，页面侧不落；ⓑ等票 274（构建带页面）与分支合并之后再一起落；ⓒ机主点名许可我写页面。**默认＝ⓐ**，⛔ 不许任何腿按ⓒ自行开工。

**4. 六格 AC 的射程要先补一枚地基（不然全是恒真）**：票面六格全 `[ ]`，但其中 `Whitelist fuzz`／`Concurrency routing`／`Forged panel allow` 三格量的是**入向**，今天有产码可攻；`Resync`／`Backpressure` 两格量的是**出向**，而**出向今天零产码** ⇒ ⚠ 在这两格上"写了判据并跑绿"是可能的（因为没有被测物）——**这正是本票最容易自证通过的地方**。⇒ 我给后续腿加一条前置尺：**任何声称覆盖出向的判据，必须能指名它断言了哪个 Go→页调用点**（今天可指名的只有 `.scratch/wisp/probes/33/p1/q2/main.go`，⛔ 探针不算产码）。

**5. 顺带一枚账面过期（不在本票射程，只登记）**：`cmd/wisp/panel_host_windows.go:403-404` 注释逐字写着 "whose whitelist is fixed at the **four** existing methods (bridge.go:42-45…)"，而 `bridge.go:42-49` 今天**有 6 枚**（票 248 那两枚已进册）；同树另一处 `cmd/wisp/run.go:695-697` 还写着 "this tree carries no WebView2 host (tickets 33/35)" 而 `:405` 真在 `Bind`。⇒ 两句都与代码不符，⛔ 我在本票里不改它们（属票 33/248 的注释面），只具名登记，交给那两票的验收腿处置。

**6. 排程**：本票**仍按住**，且它的解锁条件比从前更清楚了——① 票 274（出货 exe 带页面字节）② 票 33 的常驻接线（`NewPanelManager` 已进常驻腿，但"面板能不能真开出来"owner 10-07 已把页面那半挪到另一棵分支）③ 第 3 节那枚 ⓐ/ⓑ/ⓒ 归属。**当前 `cmd/wisp` 写面被 `272-r2` 占着**（包级互斥），本票的写腿不进本轮队列。

### 10-07 11:5x 收只读普查腿 `35-a1`（件 `.scratch/wisp/probes/35/a1/outbound-bridge.md`＝**436 行／44,457 字节**＋15 枚 logs＝2,097 行／120,512 字节；五枚 commit `6f132ca8`／`cdffb610`／`95b1bb0c`／`8976ffab`／`89114888`，我逐枚 `git show --stat` 复认只碰自己的路径）⇒ ★**它顶回我三句，三句我都现量坐实、都算它对；同时它自己的行号引错一处，也算它不对**

**1. ★它推翻我的第一句：上一节 §1 那句「页→Go（入向）＝真在」只对 Go 那半截**。
我复跑的协议链（⛔ 不是转述它的）：库根 `webview.go:103` 把 `chromium.MessageCallback` 系成 `msgcb`；`pkg/edge/chromium.go:233-245` 的 `MessageReceived` 用 `TryGetWebMessageAsString` 取串后**先喂 `MessageCallback`、再 `PostWebMessageAsString` 把同一条原样回丢给页面**；
`msgcb`（`webview.go:140-158`）把串按 `{id, method, params}` 解析——页面那枚信封 `{"method":"panel.approval.request","correlationId":…,"outcome":…}`（`frontend/src/lib/panel.ts:179-186`）**缺 `id`／缺 `params` 但 JSON 合法**，
⇒ Go 的 `json.Unmarshal` **不报错**、解出 `id=0`＋未知 method ⇒ 走 reject 支 ⇒ `w.Eval("window._rpc[0].reject(…)")` 打进**没有这个槽位**的页面 ⇒ JS TypeError 静默；
而那条被 `PostWebMessageAsString` 回丢的信封页面也**没人听**（下条④）。⇒ **Go 侧 `dispatchRaw`（`:630`）与那 6 枚名册今天收不到任何来自页面的请求**；上一节那条链（`:405 Bind` → `:630` → `composer_dispatch.go:181`）在代码上是真的，**但它的第一跳只能由 Go 测试直接调**（正控＝`panel_host_*_test.go` 里真调过），页面从不叫。
⇒ 已落成上面那枚新框 **AC#「Transport agreement page↔host」**（未勾、⛔ 本票不选形），⛔ 后续腿不许用"删白名单"或"让 Go 直接回 RPC 解析器但信封从未到达 `dispatchRaw`"来把它做绿。

**2. ★它推翻我的第二句：上一节那句「出向整条不存在」说过头了**。现量：库的出向原语三枚都在（`Eval`／`Init`／`SetHtml`），**而且每次 binding 回执都真的走 `Dispatch(Eval(…))`**（`webview.go:148/152/156`）；
我们缺的是**无人发起的推送**那一支——我自己数的四把枚数（`cmd`＋`internal`，非 test）：`\.Eval(`＝**0**、`\.Init(`＝**0**、`\.SetHtml(`＝**3**、`\.Dispatch(`＝**7**。
⇒ 正确口径＝**"泵到面板线程的路真在（`Dispatch` 7 枚调用点）、把 JS 推给页面的那一句今天零调用者"**；⛔ 别再写"出向整条不存在"。

**3. ★它推翻我的第三句，也是我自己那把尺的病**：我上一节用 `PostWebMessage|EvaluateScript|CreateWebMessageAsJson` 这族拼写去扫，报"0 命中"⇒ **那是拼错名字尺的假阴性**：
这三枚是 `pkg/edge` 的 **COM vtable 名**（`corewebview2.go:106-109`），不是这套 Go API 的拼写。我复跑同一把尺在**全仓**的命中＝**3 枚**（⛔ 不是 0）：`.scratch/wisp/probes/33/p1/q2/main.go:16` 一句注释、
`internal/panel/composer_dispatch_test.go:461` 一枚**空桩**（`func (w *CoreWebView2) PostWebMessageAsJson(json string) error { return nil }`）、`:561` 一句注释。
⇒ 教训按第 111 条入账：**判"某符号零调用者"要先确认这符号是不是这套 API 的真名**；而"3 枚命中全是注释／空桩"这句才是本件事的正确读数——**它同时说明测试侧有一枚会把出向吞成 `nil` 的桩**，⛔ 谁以后拿"这包绿"抵"页面收到过东西"都要先过这枚桩。
⚠ 另两处顺带坐实：`AddHostObjectToScript` 全仓＝**0**（可页面 `panel.ts:146` 的注释写着桥"Installed by WebView2's AddHostObjectToScript / postMessage pipe"＝**注释描述的机制不存在**，与 §1 那条协议断裂同源）；`go.mod:19` 把 `jchv/go-webview2` 标着 `// indirect`，而 `cmd/wisp/panel_host_windows.go:64`／`panel_resident_windows.go:73` **直接 import 它** ⇒ 账面过期（`go mod tidy` 会把它挪成直接依赖），⚠ 只登记、本轮不动 `go.mod`。

**4. ★它自己的读数有一枚硬伤（页面那半的行号引在**另一棵树**上）**：它写 `frontend/src/main.tsx:97` 以 `<App />` 无 prop 挂载、`App.tsx:91` 的 `snapshot = EMPTY`、常量在 `:70-73`——
dev 上 `main.tsx` 只有 **61 行**、`App.tsx` 只有 **117 行**，**这两个行号在 dev 上根本不存在**；那组号是**分支 `dsh/feat/frontend-p0-v2`** 的版本（`main.tsx` 115 行、`App.tsx` 531 行）。
dev 真身我现量：挂载＝`frontend/src/main.tsx:58`（`{harness === "1" ? <AppHarness /> : harness === "2" ? <Showcase /> : <App />}`，**确实无 prop**）、
默认值＝`App.tsx:70`（`snapshot = EMPTY,`）、`EMPTY` 定义＝`App.tsx:49-53`（`pending: []`／`results: []`／`composer: EMPTY_COMPOSER`／`generatedAt: ""`）。
⇒ **它的结论在 dev 上同样成立，但锚必须按 dev 重钉**；这条要记牢是因为**票 274 正在把 dev 这 64 枚 `frontend/src` 构进 exe**，接收器／桥的落地判据只能钉在 dev 的树上，⛔ 不许拿分支的行号当凭据。
★页面"没有耳朵"我在 dev 上单独复跑：`frontend/src` 里 `addEventListener` 只有 4 处命中且**没有一处是 `message`**（keydown×2／mousemove／`main.tsx:50` 那句 close 按钮），`window._rpc`／`onmessage`／`WebMessageReceived`＝**0** ⇒ "Go 推过去也没人接"成立。

**5. 其余两格它答了、我也复认**：`PanelManager` 非 test 构造点唯一＝`cmd/wisp/panel_resident_windows.go:253`（我同一把尺复跑：其余 11 处全在 `_test.go`）；
`SnapshotPump.Publish` 非 test 调用者只有 `cmd/wisp/panel_pump.go:405`（复跑到；另 `internal/agent/loop.go:1041` 那枚是 `opt.Sink.Publish`、**不是同一座**，别混）。
可抄先例＝同文件 `:349-372` 那个 `post()`（两态路由 `Dispatch`／`tasks`）。§5 那 6 枚代价行＋四条禁区我读进本票面射程，⛔ 一枚都不替它翻。

**6. 它自报未做完的两格我不代填**：真机再入测量（它零编译权限）、`addEventListener('message')` 只证到"本仓无先例"（不等于"WebView2 里不可行"）。⇒ 归后续腿逐枚处置。

**7. 排程（本轮变更）**：本票**仍按住**，但解锁条件少了一枚、多了一枚——
① 票 274（`274-r1` 10-07 11:5x 已开工，出货 exe 带页面字节）② 票 33 的常驻接线（`NewPanelManager` 已进常驻腿）③ ★**新增前置**＝上面那枚 **AC#「Transport agreement」**：桥的信封格式与泵的接收器不接上，本票那六枚端到端判据**在任何一棵树上都无法执行**；
⛔ 本轮不派本票写腿（`cmd/wisp` 写面此刻由 `274-r1` 的构建与门禁占着，且桥的落地形状还没裁）。归属问题（ⓐ 我出表／页面由别人落 ⛔ 默认、ⓑ 我写页面文件〔未放开〕、ⓒ 交给别的会话〔⛔ 禁〕）**维持 ⓐ 不变**，
但按 §1 的新读数修正一句：**协议断裂那一半是纯 Go 侧的活**（改宿主如何收信封），不需要动页面文件、也不需要动 C17 名册 ⇒ 本票不必等机主再给权限。
