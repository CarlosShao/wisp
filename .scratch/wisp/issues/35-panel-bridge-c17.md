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
- [ ] **Inbound judge must ride the page's own envelope (added 10-07 by the orchestrator from `35-a2` §4; mandatory — no existing nail may substitute for it)**: one Go-side test that (a) takes the raw envelope string verbatim
      from what the page actually posts (`frontend/src/lib/panel.ts:180-184`), (b) enters through the transport the page uses **after** the chosen shape lands — ⛔ not by calling `dispatchRaw` directly and ⛔ not via
      `w.Eval("window.wispDispatch(…)")` (that shape lives at `cmd/wisp/panel_resident_windows_test.go:301/:805`; it is *the test playing the page*, and `35-a2` §4 N6 measured it stays green under shape 甲 ⇒ zero
      discriminating power for this edge), (c) asserts `dispatchRaw` (`panel_host_windows.go:630`) or its downstream was reached **carrying that correlationId**, and (d) asserts an unregistered method name is still refused.
      ★Falsifiability: deleting / no-op-ing the forwarding statement **must** redden this named test — if it stays green the judge is not reading this edge and is rejected. ⛔ N1 (Go roster count) and N2/N3 (page-wording
      counters) may not be claimed as evidence for this edge.
      ★**编排者 10-07 13:4x 落地约束（据 `35-a3` 的具名顶回＋我自己复跑守卫顺序后裁；不改上面四件要求，只钉死用哪枚信封）**：上面 (a) 那句"verbatim from what the page actually posts"若照 `HEAD:frontend/src/lib/panel.ts:181` 那枚 approval 信封直落，**今天永远红**——`panel.approval.request` 不在 dev 的 6 枚名册里，会被**名册守卫**挡下 ⇒ (c) 那格"到了 `dispatchRaw` 且带着关联号"**结构上不可能成立**。⚠ **我把守卫顺序自己复跑了一遍，腿 §3-b 那句"还有第三层：`:135`（来源）与 `:139`（缺 requestId）"对这枚信封不成立**——它"这两道先于路由"那半句是对的，但它把这两层算成那枚 approval 信封**会撞到的**后续拒绝＝错：`ParseComposerRequest`（`internal/panel/bridge.go:126`）的次第是 **解析 `:128` → 名册 判 `:132`／拒 `:133` → 来源 判 `:135`／拒 `:136-:137` → requestId 判 `:139`／拒 `:140-:141`**（三枚 `if` 与三枚 `return fmt.Errorf` 的行号我都用 `grep -n` 对过），名册守卫在最前 ⇒ 那枚信封**只可能在 `:133` 被拒，`:135`/`:139` 两道根本走不到**。⇒ 裁成三支各钉各的文案（这才是"有辨别力"而不是"看它红"）：
        - **(c) 到达性正控改用名册内的那枚信封原文**：`HEAD:frontend/src/lib/panel.ts:246` `sendRequest("panel.mode.request", { to })`（同文件 `:211` 的 `sendRequest` 给每请求带上 `requestId: pc-<seq>-<uuid>` ＋ `source: "panel-composer"`，正是 `:135`/`:139` 那两道要的两枚形状字段）。⛔ 不许改用别的层拼的假信封。
        - **(d)-1 未注册名**＝就用页面今天真发的 `panel.approval.request`（比造一枚 `panel.bogus.request` 强），断红句逐字含 **`方法 … 不是面板 composer 通路的能力入口`**（`:133`）。⚠ **本格的新料不是"它被拒"——那一格已有既有钉**（我 13:5x 现量，腿整族漏了，见下面"收件"一节）：`HEAD:internal/panel/l2_grant_boundary_test.go:1964`/`:1965` 已经拿页面那两枚原始串要求拒、`:1251-1267` 要求 `knownComposerMethod` 拒 11 枚候选名。⇒ (d)-1 只许按**本框的新要求**记功：**经页面用的那枚传输到达**＋**点名是哪一道守卫拒的文案**，⛔ 不许拿"既有那枚已覆盖"来抵。
        - **(d)-2 名册内但来源不符**＝`panel.mode.request` ＋ `source:"panel-composer-x"`，断红句含 **`按伪造/串台拒绝`**（`:136-:137`）。**(d)-3 名册内但缺 requestId**＝同枚去掉 `requestId`，断红句含 **`缺少 requestId`**（`:140-:141`）。⛔ 三支**不许共用一句文案**（先例＝对抗验收那条『门拦下的文案不许复用权限文案』；也＝第 120 条『一个门有多种红法，不许把几桩病压成一格』）。
        ⚠ **顺带一枚新事实（我 13:4x–13:5x 现量，⛔ 不在本框射程，别顺手"修"它）**：页面在 `panel.ts:181` **今天就在发**审批请求，而 Go 侧 `DecideFromPanel`（`internal/agent/approval/gate.go:736`）对 **`allow` 一律拒**（`:737-746`，`PANEL-ALLOW-REJECTED`＋烧 nonce＋`ErrPanelAllow`），对 **`reject` 却走 `:748 g.q.reject(…)`＝能力在、路由不在**。**但我 13:5x 把这一格写得更准：它不是"少接一根线"，是三样东西叠着的**——① 名字形状被既有仪器钉死：`HEAD:internal/panel/l2_grant_boundary_test.go:1251-1267` 要求 `knownComposerMethod` 拒 11 枚候选名（第一枚就是页面原文那枚），`:2274` 那发突变要求"bridge.go 答了它"必须红，而它判的是**词根**不是方向（`grantRouteWords` 10 枚含 `approval`/`decide`/`verdict`，`:192-195`）；② 载荷键也被钉：`grantFieldWords`（`:183-188`）里有 **`outcome`**，而页面的信封正是 `{method, correlationId, outcome}`；③ `:1275-1289` 那道 pool 审计会把"经另一枚函数答"的路由同样点红。⇒ "面板上点拒绝到底能不能真拒绝"是一枚**功能级缺口**，但补它的价钱比"加一枚名册"贵（要避开词根表＋换掉 `outcome` 键形＋过 pool 审计），已并进 `A671` 那枚〔待人拍板〕（与 `Q-76` 同族：入向要不要第五／第六枚方法名，⛔ 动它＝改 C17 白名单＝契约变更）。

## Progress log (append-only, newest last)

### 10-07 编排者现状对拉（票面 09-19 写的东西从没按今天的仓复算过；落账 `A657`）

★**为什么这节写在本票而不是新票**：机主原话「**不要重复劳动哈**」。本票就是 C17 那座桥（`Go→frontend event push` ＋ `panel.resync` 全量推），而字段对接普查（`145-f1`）量出来的头号阻断**正好落在本票的许诺上**——所以我把现量补到本票面上，⛔ 不再另开一枚同义票。

**1. 方向要分两半，今天只有半座桥**（四把尺都我自己现跑，⛔ 不引腿的读数）：
- **页→Go（入向）＝真在**：`cmd/wisp/panel_host_windows.go:405` `w.Bind(panelDispatchBinding, func(raw string) string{…})` → `:630 dispatchRaw` → `internal/panel/composer_dispatch.go:181 … HandleModeRequest(…)`。方法名册现量 **6 枚**（⚠ **10-07 13:4x 由 `35-a3` 顶回后我复跑更正**：枚数对、**行号区间错**——`internal/panel/bridge.go:42-45` 只有 **4 枚**（`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`），`config.get`/`config.set` 两枚在 **`:66`/`:67`**（`:46-65` 是那段"为什么不加 `panel.` 前缀"的注释），`knownComposerMethod` 的 case 表在 `:148`；我原来写的"`:42-49`"是**没数过就顺手写的区间**，属第 119 条那一族）。
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
> ⚠ **10-07 12:3x 就地更正（原句不抹；凭据见文末「收 `35-a2`」§1）**：上面那 **7 枚**是"全仓所有同名 `Dispatch(`"的总数，⛔ 不是这座桥的泵调用点。属 WebView2 面板线程的只有 **1 枚**＝`cmd/wisp/panel_resident_windows.go:360`（`post()` 内，`:349-372`），其余 6 枚是状态机族的同名方法。⇒ 那句应读作**"泵到面板线程的路真在（`Dispatch` 调用点 1 枚）、把 JS 推给页面的那一句今天零调用者"**。

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

### 10-07 12:3x 收只读普查腿 `35-a2`（件 `.scratch/wisp/probes/35/a2/transport-cost.md`，**744 行／49,788 字节**＋`logs/`；commit `f88144b3`；★它顶回我票面一句枚数，那一句**它对**）＋编排者裁形

框数尺现量：本票 **8 枚框／0 勾**（本轮新增第 8 枚＝上面那枚 **Inbound judge must ride the page's own envelope**）。⛔ 未 `-done`。

**1. ★它改我票面 `:96` 那句「泵到面板线程的路真在（`Dispatch` 7 枚调用点）」——枚数错了，按 1 枚算（原句不抹，就地认）**：
我那句把全仓所有 `\.Dispatch(` 当一个通道数了。它的尺（`logs/q2-init-semantics.txt`）＋**我自己的复跑**：非 test 全仓 `Dispatch(`＝**7 枚**这个总数**是对的**，但其中属于 **WebView2 面板线程泵**的只有 **1 枚**＝`cmd/wisp/panel_resident_windows.go:360 w.Dispatch(fn)`（在 `post()` 里，`:349-372`，两态路由 `Dispatch`／`tasks`）；**其余 6 枚是状态机那族的同名方法**，不是这座桥。
⇒ 正确口径＝**"面板线程的 `Dispatch` 调用点今天＝1 枚；把 JS 推给页面的那一句零调用者"**。⛔ 谁以后要论证"出向有路"，只能引那一枚，⛔ 不许再引"7 枚"给我这句话充数——**这类"同名符号凑分母"正是第 107 条的形状，我这次是自己踩的**。

**2. 规格射程（它的问题①，我复跑对回）**：`docs/PLAN.md` 与 `docs/specs/**` 里**没有任何一句写死传输方式**（三个候选拼写全 0 命中，命中只在 evidence/ledger 层）；被规格写死的只有**依赖名**与**资源通道**两枚邻居＝`SPEC-08:145`（`jchv/go-webview2`）与 `SPEC-08:148-149`（`AddWebResourceRequestedFilter` ＋ 禁 localhost HTTP）。⇒ **选形不被规格绑**，但乙形若把库从 `go.mod` 摘掉会同时撞 `SPEC-08:145` 与既有钉 N4。⛔ 它与我都没改 `docs/**` 一字。

**3. 三形代价（它只报代价，我据此裁；表在它 §2.4）**：甲＝`cmd/wisp/panel_host_windows.go` **1 枚文件**、一条语句＋一枚 JS 常量、页面 0 枚、`go.mod` 不动、名册不变；⚠ 时序陷阱＝`Init` 必须早于首次 `SetHtml`（`:449/:468/:681`）且要活得过后续 `SetHtml` 重建，且需 `window._rpc` 守卫。乙＝**替换宿主创建＋泵整条路**（`:386`→`:415` 及 pump/`Embed`/`Resize`，数枚文件），要重实现 STA 泵（该文件头 `:29-44` 逐字记着冻结／嵌套泵的坑），若保留包装还要动 `go.mod`/`go.sum`＝**人工批准**，并在 `TestPanelHostIsAttachedAndNamesTheWindowHops` 上**真红**。丙＝只改页面 `frontend/src/lib/panel.ts` 两枚发送点（`:179/:218`）＋`:139-155` 类型，⛔ **页面写权未放开**，且它自己算过"丙仍然拿不到回执"。

**4. ★★它最值钱的一格（恒真性，我复认）＝今天没有任何一枚既有钉能对"页面的请求到不到 Go"分辨**：N2/N3（`bridge_test.go:162`、`composer_test.go:510/:518-522`）数的是**页面词面**（甲不动页面 ⇒ 恒绿）；N1（`inbound_roster_253_test.go:419`）数的是**Go 名册**（与传输无关 ⇒ 恒绿）；★**N6 是同一族的更强一击——`cmd/wisp/panel_resident_windows_test.go:301/:801/:805/:808/:864` 这些真机用例今天就是"从 `Eval` 里直接叫 `window.wispDispatch(…)`"**，也就是**测试自己扮演了页面**，⛔ 因此它对"真页面的请求能不能到 Go"**零分辨力**（甲落地后照旧绿）。⇒ 我据此**新增上面那枚框**，并把它写成"判据必须自己经这条边、且删掉转发语句必须红"。⛔ 后续任何腿不许拿 N1/N2/N3/N6 抵这一格。

**5. 它给的三发定向突变与那枚已知盲区（我采纳进验收口径）**：ⓐ 删掉/中和转发 ⇒ **只有新判据能红**（旧钉全绿）；ⓑ "什么都不解析直接 ack" ⇒ 既有三枚会红（`inbound_roster_253_test.go:465`、`composer_dispatch_test.go:201/:231/:126`），⚠ 但若新判据仍从 `Eval`／直调 `dispatchRaw` 进去，ⓑ 会**在被绕过的地方红、在真边上绿** ⇒ ⓑ 的验收必须同时指认"新判据也红"；ⓒ（它自己加的第三发，成本 0）＝把垫片改成手造 `window.external.invoke({id:0,method:'wispDispatch',params:[raw]}` ⇒ **今天一把尺都不红**（`external.invoke` 全仓 0 命中、`inbound_roster_253` 不看 JS、`hostChannelCapabilityHits` 跳字符串与 `frontend/`）⇒ **这条子形是已知盲区**，我因此把 ⛔ 它禁进选形（见 §6）。

**6. ★编排者裁形＝甲，且只准**甲的子形①**（转发原始串进**已绑定的** `window.wispDispatch(raw)`，走库的 binding 回执路）**：
- **为什么不选乙**：乙要重写宿主创建＋泵、动 `go.mod`（人工批准）、摘库即撞 `SPEC-08:145`＋N4，收益只是"不经库的 RPC 管道"——而甲子形①本来就不需要页面的 `{id,method}` 信封进那根管。**代价不成比例。**
- **为什么不选丙**：页面写权未放开（机主 10-07 放开的只有 `frontend/**` **只读**）；且它自己算出丙没有回执。⇒ 丙只作为对照行留在它 §2.3。
- **⛔ 禁甲的子形②**（手造 `{id,method:'wispDispatch',params:[raw]}` 直接喂 `external.invoke`）：§5 ⓒ 已证这条**今天没有任何尺看得见**，用它＝把一条无牙的路选成产品边。若以后要用，**必须先给它单独造一枚判据**再谈选形。
- **⛔ 落地禁区（写进派单）**：不许新增点分法名（N1 会在加名的那一刻红，而名册属 C17 面＝人工批准）；不许把新判据写成"来者都 ack"；不许动 `internal/panel/assets.go` 的 fail-closed；⛔ 不许用 `w.Eval("window.wispDispatch(…)")` 那种"测试扮演页面"的形状当凭据（N6 教训）；⛔ 不动 `frontend/**` 一个字节；⛔ 不给 CI 加 `-tags winlive`（真机那两枚 `winlive` 用例只登记〔仅本机可量〕）。

**7. 它复认并扩了我票面的两处，另交回五枚"只是注释过期"（⛔ 都不是门，写腿不许为对上它们而改代码，也不许删注释换绿）**：M1 `panel_host_windows.go:402-404` 写着白名单"fixed at the **four** existing methods"，`bridge.go:42-47` 现量 **6 枚**（与我票面 §1 同数）；M2 页面 `panel.ts:146` 那句"Installed by WebView2's **AddHostObjectToScript** / postMessage pipe"——`AddHostObjectToScript` 全仓 **0 命中**＝**注释描述的机制不存在**；M3 `composer_dispatch_test.go:461` 那枚空桩拼的是 `PostWebMessageAsJson`，而库里真名是 **`PostWebMessageAsJSON`**（大写 JSON，`corewebview2.go:106`），且它挂在 `writeHostCarrier` 现造的 `type CoreWebView2 struct{}` 上 ⇒ **不是库类型的实现**（★这条接着修我上一节：我按第 116 条把 vtable 名与 Go API 名分开了，但没验"这枚桩到底实没实现那个接口"）；M4 `PLAN.md:1747` 指向不存在的 `docs/contracts/`（非本票射程）；M5 `go.mod:19` 把直接 import 的库标成 `// indirect`（只登记，本轮不动 `go.mod`）。

**8. 它顶回我与 `35-a1` 的页面行号（我复量它对）**：`hostBridge()` 真身＝`frontend/src/lib/panel.ts:152-155`（`:146` 是 `declare global` 里那枚过期注释行），"不回东西、只等 PanelSnapshot push"那句＝`:162-166`；票框里引的 `:141`/`:179`/`:218` 三枚**照旧正确**。⇒ 我那枚新框 §5.1 的引用锚按它改。

**9. 它自报未做完的格我不代填**（归 `35-v1`／写腿）：`Init` 语义的两条源码级承重（它 §0 `:391` 具名"单参数是否原样交给 Go 闭包"只读到 `:140-158`／`:450-482`）、真机再入测量（零编译权限）、以及 §5.2 那枚空桩今天被谁依赖（只登记未动）。

**10. 排程（本轮变更）**：`cmd/wisp` 写面**现在空**（`274-r1` 已交完并入库 `35633445`；`272-v1` 已收）⇒ **本票的落地腿解除按住**，可派 `35-r1`（唯一写面＝`cmd/wisp/panel_host_windows.go` ＋ 一枚新 `cmd/wisp/*_test.go` 判据；⛔ 包级互斥：`cmd/wisp` 同一时刻只一枚写手，且它与 `272` 的 AC#5 整包复跑**不得同批**）。⛔ 只 commit 不 push。

### 收 `35-a3`（编排者 10-07 13:5x 收件；只读普查腿；⛔ 本票 AC 框一枚没翻——AC#0 那句"落地腿自己现跑复认"仍归落地腿）

**1. 凭据与纪律**＝`.scratch/wisp/probes/35/a3/field-reconciliation.md`（**246 行**）＋同目录 31 枚现量件；三笔提交 `6a976294`→`0f0fdaf2`→`06408f3b`。全程零 `go` 命令（派单禁的，它 §6-b 具名报了因此判不动的那两枚双向尺）；分支读数只走对象层（`git show <ref>:<path>`／`git grep <pat> <ref> --`／`ls-tree`／`cat-file -e`），⛔ 零 `checkout`／`switch`／worktree／进那棵工作树；`frontend/**`＋`design/**` 零写入；⛔ 未 push。我复跑的归位尺：`git status --porcelain -- .scratch/wisp/probes/35/a3`＝**0 行**（它的写面确实全入库）。

**2. 三把我自己复跑的尺，复认为真**：① `git diff HEAD 16c2f038 -- frontend/src/lib/panel.ts`＝**0 行**（两棵树各 311 行）⇒ 它 §0 那句"页面对账两棵树共用同一枚契约件"成立，这条比我的派单转述更极端、我派单那句"dev 上是较旧的页面源码"只对到装配层（`App.tsx` 117↔531、`frontend/src` 64↔85 枚）；② `git grep -n -E 'panel\.approval\.request' HEAD -- '*.go' ':!*_test.go'`＝**0**，且 `HEAD:internal/panel/bridge.go:146-152` 的 `case` 在 **`:148`** 一次列全 6 枚不含它 ⇒ 入向名册确实不答这枚；③ 三道守卫行号逐枚 `grep -n` 对过＝解析 `:128`／名册 判 `:132` 拒 `:133`／来源 判 `:135` 拒 `:136-137`／requestId 判 `:139` 拒 `:140-141`，`DecideFromPanel`＝`HEAD:internal/agent/approval/gate.go:736`（allow 拒 `:737-746`＋`ErrPanelAllow`，reject 走 `:748 g.q.reject(…)`）。⚠ 它引的 `gate.go:640-642` 是**注释**、我引的 `:736` 是**代码**——两层都对，⛔ 以后别混着引。

**3. ★它顶回我票面 `:77` 那枚行号区间＝对，我已就地更正**（`:42-45` 只有 4 枚、`config.get`/`config.set` 在 `:66`/`:67`、`:46-65` 是那段"为什么不加 `panel.` 前缀"的注释、`case` 表在 `:148`；枚数我一贯是对的、区间是没数过顺手写的）⇒ 记进第 119 条那一族。

**4. ★我推翻它两处，并推翻我自己一处**：
- 它 §3-b 那句"还有第三层：`:135`（来源）与 `:139`（缺 requestId）"对**这枚** approval 信封**不成立**——名册守卫在最前，那枚信封走不到后两道（它"这两道先于路由"那半句是对的，我上面框里那句"`:135`/`:139` 先拒是错的"**写重了，已按此改准＝记我**：更正别人时不许把原句换成一枚更错的判断，定性词不许重于证据，对腿也一样）。
- 它 §3-b 那句"全仓命中只有 `composer_test.go:417` 与 `frontend_hygiene_test.go:32`"是**假枚举**。尺＝`git grep -n -E 'panel\.approval\.request' HEAD -- '*.go'`（**不加** `':!*_test.go'`）现量＝**17 行／3 枚文件**，其中 `HEAD:internal/panel/l2_grant_boundary_test.go` 一枚占 **15 行**。形＝它自己 §7 刚认过的那族（把加了过滤的清单当成全仓清单）。
- 它 §7 那句"起手登记的 8 枚 ` M probes/161/r6/logs/flip-*.txt` 在终态已从 status 消失＝别的腿在此期间自己入库"（第三笔 `06408f3b` 又重述一遍）**不可复现**：我 13:54 现量那 8 枚**仍是 ` M`**、另多出 1 枚 `?? flip-7.txt`（共 9 行），而 `git log -- .scratch/wisp/probes/161/r6/logs/flip-1.txt` 最新一枚是 **09-27 16:37 `422c1bc7`**，13:40–13:44 之间没有任何提交动过它。⇒ 它自证那三枚 dirty 数（752／781／752）里**"差集归位解释"那一半降档读**（原句不抹），"只有自己的写面被动过"那一半我复跑为真（第 1 条那把尺＝0 行）。

**5. ★这一格里最硬的料是它和我先前都不知道的一族**：`internal/panel/l2_grant_boundary_test.go` 今天已经把"Go 答审批路由"钉成违规，而且**判的是词根不是方向**——`grantRouteWords` **10 枚**（`:192-195`，含 `approval`/`decide`/`verdict`）、`grantFieldWords` **19 枚**里含 **`outcome`**（`:183-188`，正是页面那枚信封的键名）、候选名 **11 枚**第一枚就是页面原文（`:1251-1267` 要求 `knownComposerMethod` 拒它们）、`:2274` 那发突变要求"bridge.go 答了它"必须被点名、`:1275-1289` 那道 pool 审计还会抓"经另一枚函数答"的路由。⇒ **它 §4.2 那张"落地会撞到的既有钉"表（11 组）整族漏了这一组**（尺＝`grep -c l2_grant` 在它 246 行正文＝**0**），而这恰恰是 35-r1 最不能撞的一族。已把这族写进本票上面"落地约束"那一条的 (d)-1 与〔新事实〕段（含"(d)-1 不许拿既有那枚已覆盖来抵账"）。

**6. 五档计数按它的口径收下，并逐档标凭据层级**：同名同义 **41**／同名不同义 **0**／异名同义 **7**（5 对型名＋2 对跨层）／Go 有页面不读 **15 项（展开 45 子键）**／页面读、Go 没有 **4**。⚠ 只有**最后那一档 4 枚**我抽验过（两枚：`panel.approval.request` 用第 2 条那把尺、`snapshot.view` 读 `HEAD:internal/panel/composer.go:57-92` 全结构确认无 view），**41／0／7 三档我没逐枚复跑**，且"同名不同义＝0"它自陈是**读码判**（⛔ 零运行时）⇒ 后续任何腿要引这三档必须复跑，不许当常量。

**7. 它"照今天这张表落地后仍拿不到输入的格子"六枚收下**（L2「允许」结构性无解／审批答复整跳不在名册＝C17 面／九屏里 7 屏 `fed:false`／git chip 手传字面量且与 `GitView` 不同名不同形／模型与密钥格页面零声明零读取／`view` Go 侧零键）。⇒ **对本票的直接后果**：35-r1 只解决"名册内在册那几枚能不能到"，⛔ 不许在它里面顺手扩快照字段——那要撞 `HEAD:internal/panel/pump_test.go:123`、`:291` 与 `HEAD:internal/panel/subagent_roster_197_test.go:214` 那三枚"四键字节钉"（逐字 `!= "composer,generatedAt,pending,results"`，我 13:5x 现量到三处而非它报的两处），属另一枚具名解冻面。

**8. 排程不变**：`cmd/wisp` 写面仍被 `274-v1` 占着（本票 AC#5 那半还欠我一次 `go vet ./cmd/wisp/`，⛔ 不许与它同批）⇒ **35-r1 继续按住**，派单要点已钉死在本票上面那五条落地约束里。它 §8 报的四条原文冲突我逐条裁过：第 1 条成立（已改票面）、第 2 条成立（派单那张候选起点表作废，真起点＝它 §1.0 那 9 枚）、第 3 条成立且比转述更强、第 4 条那句"`panel_host_windows.go:205` 注释把页面 postMessage 说成'第二条腿'"我复跑为真但⛔ 注释我没动，交 35-r1/验收腿按"不许删注释换绿"处置。
