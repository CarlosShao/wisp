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
- [x] **Transport agreement page↔host (added 10-07 by `35-a1`, shape NOT chosen here)**: the page posts its envelope on
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
- [x] **Inbound judge must ride the page's own envelope (added 10-07 by the orchestrator from `35-a2` §4; mandatory — no existing nail may substitute for it)**: one Go-side test that (a) takes the raw envelope string verbatim
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

- [ ] **Behavioural yard must model `this` and property writability, and `:52`'s live proof must stop being machine-local (added 10-07 by the orchestrator from `35-v2` W1/W2; ledger `A684`; mandatory — it is NOT satisfied by the current 1824-line interpreter)**: `35-v2` attacked the behavioural judge `cmd/wisp/panel_transport_35r2_test.go` with two one-line mutants of the shipped hook (via `go test -overlay`, repo untouched) and **both stayed green** — M-A replaced `native.call(cw, message)` with a bare `native(message)` (a real browser may reject that as an "Illegal invocation"; the fixture's native stub at `:1238` discards the receiver, and `jsObject.set` at `:137` has no concept of `writable`/`configurable` — zero occurrences of `writable|configurable|defineProperty` in the whole file), M-B deleted the `typeof … !== "function"` half of the guard. ⇒ (a) the fixture must **fail loudly** on both shapes (drop-the-receiver and silently-ignored-override), ⛔ 不许靠"多写几枚断言"糊过去 — each face needs its own named red; (b) keep the re-entry/sequencing face it already has (M-C1/M-C2 went red as designed, that half has teeth); (c) `:52`'s decisive reading today lives **only** in `-tags winlive` on this desktop (`cmd/wisp/panel_transport_live_35v2_windows_test.go`, PASS 8.25s with the page's own `chrome.webview.postMessage` reaching the Go door vs the `fb2fb802` hook recording `[]` and FAILing at 58.39s) — winlive has **no CI carrier**, so either give that edge a CI-reachable stand-in or state 〔仅本机可量〕 on its face and re-run it in every wave that touches the transport. ⛔ 本框不许用"解释器绿"来满足 (c)。

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

**2026-10-07 18:3x 编排者（据 `35-r1` 三笔 `f718e9b6`／`fb2fb802`／`1e42fc61` ＋ 我自己现量，台账 `A681`）**：AC#6 的传输形状＝票 `:150` 已裁的**甲子形①** 已落 Go 侧——`cmd/wisp/panel_host_windows.go:633` 窄接口 `pageTransport`、`:648` `panelPostMessageForwardInit`（把页面 `window.chrome.webview.postMessage` 转发进**已绑定的** `window.wispDispatch`）、`:664` `installPanelTransport`（先 `w.Bind` 门后 `w.Init` 钩子）、`:408` `bringUp` 改调它；`frontend/**` 一字节未动（`git diff --stat 10899b77 -- frontend/` ＝空）、C17 名册未变。新用例 `cmd/wisp/panel_transport_35r1_test.go:195` `TestPagePostMessageEnvelopeReachesDispatchRawViaTransport` 覆盖的是**本票 `:63` 那格的 (a)(b)(c)**（正控逐字＝`git show HEAD:frontend/src/lib/panel.ts:246`；入口是页面的那条 postMessage 边，⛔ 未直调 `dispatchRaw`、⛔ 未用 `w.Eval("window.wispDispatch(…)")`；断下游 `modeSpy` 恰一次且带 `requestId`），**(d) 三支守卫红文案（`bridge.go:133`／`:136-137`／`:140-141`）一支都没覆盖**，腿在 impl.md §④ 具名写了"未覆盖"、归 `35-r2`。⇒ **AC#6 与 AC#7 两框都不翻**：AC#6 的非实现者验收腿 `35-v1` 在飞（要取的是"真 WebView2 里那条转发钩子到底跑不跑"——本用例的 `forwardingInstalled()` 是**读串**不是跑 JS，`-tags winlive` 零读数；以及 `Init` 早于首个 `SetHtml` 的顺序证），AC#7 缺 (d)。⚠ 本机口径（我 18:36 现跑）：`go vet ./cmd/wisp/ ./internal/panel/`＝rc=0，而 `go test ./cmd/wisp/`＝**`exit status 0xc0000135`**（sherpa/onnxruntime 原生 DLL 不在 PATH）⇒ 派 `cmd/wisp` 写腿的简报必须自带"把 `%GOMODCACHE%\…\sherpa-onnx-go-windows@v1.13.8\lib\x86_64-pc-windows-gnu` 那三枚 DLL 拷到测试 exe 旁、CWD＝`cmd/wisp`"这句，否则整包红名册那一格任何腿都只能交 `0xc0000135`。

**2026-10-07 19:0x 编排者更正裁断（据 `35-v1` 的 G2＋我自己读库原文；台账 `A682`；⛔ 本条不改任何判据文字，只改 `:150` 那枚子形的选择与一处细节描述）**：验收腿 `35-v1`（`4b5c756a`／`2f9fe600`）判 **AC#6 今天不该翻**，我按它判，并把我自己的定性写死在这里——

- ★**`:150` 我选的甲子形①（把页面 `postMessage` 折进已绑定的 `window.wispDispatch`）落码后在真浏览器里是死循环，不是"生效"**。凭据（`github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808`，我逐字读过）：`pkg/edge/chromium.go:112` ＝ `e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")` ⇒ 库的桩**唯一出口就是被这枚钩子覆写掉的那枚属性**；`webview.go:462-478` 的桩同步建 `window._rpc[seq]` 再 `window.external.invoke(…)` 并**返回 Promise**（⛔ 同步拿回执是我此前的误读）。⇒ 已落的钩子每跑一次就再生成一枚嵌套 RPC 帧、再进自己，**原生出口一次都没被调** ⇒ Go 侧零到达。
- ⇒ **改裁甲③**（同一条 Go 侧一文件的射程内）：转发钩子**先把原生 `cw.postMessage` 存成局部引用**，再用再入守卫把**库自己产生的 RPC 帧**原样交给原生出口，只有页面的裸信封才折进 `wispDispatch`。⛔ 三条禁区一字不放宽（不动页面文件、C17 名册不变、不许靠弱化白名单或空答 RPC 换绿）。
- ★同批硬要求（这格的牙换形）：判据里那把 `forwardingInstalled()`（`cmd/wisp/panel_transport_35r1_test.go:127`）**只认字符串**，`v1` 具名报了"两发突变红都落在这道词面正控、到达性断言根本没轮到"⇒ 从此**必须建模再入**：fake 的 `external.invoke` 在**调用时**取当前 `chrome.webview.postMessage`、桩同步注册 `_rpc` 槽，并带**再入深度帽**；**今天这条旧钩子在新尺下必红，甲③ 下才绿**。词面检测只留作正控#1，⛔ 不许再当"到达了"。
- ⚠ 票面 `:58-59` 那句"replies by evaluating `window._rpc[0].reject(…)`"**细节过期**：`msgcb`（`webview.go:139-160`）对未绑名走 `callbinding → nil, nil`（**不是 error**）⇒ 成功支 ⇒ Eval `window._rpc[0].resolve(null); window._rpc[0] = undefined`，而页面没有 `_rpc[0]` ⇒ 实际是一枚被吞掉的 TypeError。**原句保留不改**（它是判据的来由），本条只做具名更正；"请求死了、Go 侧无日志、页面无回执"那句结论不变。
- ⚠ 行号漂移具名（`35-v1` 现量，⛔ 以后派单不许再引旧号）：`dispatchRaw` `:630`→**`:677`**；`:144` 里的 `:449/:468/:681`→**`:448/:467/:671`**；C17 名册在 `bridge.go` 守卫里是**六枚**（四枚 composer＋票 248 的两枚 config），不是四枚。
- ⇒ 排程：`35-r2`（甲③＋行为尺＋旧形必红反证）→ `35-v2`（非实现者，含真窗载体那一发：走 `bringUp` 真路径、由 Go 门报逐字到达，⛔ 不许 `Eval` 扮页面）→ `35-r3`＝AC#7 的 (d) 三支（今天 0/3）。**AC#6／AC#7 两框都保持未勾。**

**10-07 20:1x｜写腿 `35-r2` 进展登记（只追加，⛔ 不改本票上面任何原句；AC#6/AC#7 两框本腿一个字没动）**：甲③ 已按编排者 `a7f9781c` 那笔的改裁落地，选**形ⓐ（先存原生出口＋再入守卫）**，落点现量＝`cmd/wisp/panel_host_windows.go:638`（注释）／`:673`（声明，`const`→`var`，绑定名由 `panelDispatchBinding` 以 `fmt.Sprintf` 注入 JS，JS 文本零硬编码 `wispDispatch`）／`:700`（`w.Init(...)`）／`dispatchRaw` `:706`。判据已换成**行为尺**：新件 `cmd/wisp/panel_transport_35r2_test.go` 内置一枚只为此而生的 JS 子集解释器，**逐字执行** `chromium.go:112` 的 `external.invoke` 垫片原文、`webview.go:462-478` 的 `Bind` 桩原文（按库的拼接式生成）、以及**产码运行时那枚钩子字符串本身**（⛔ 不抄文本），`msgcb`/`callbinding` 按 `webview.go:131-168` 复刻（未绑名＝`resolve(null)` 成功支，⛔ 非票 `:58-59` 旧句的 `reject`），再入帽 8。**★硬判据两发都取到**：把 `fb2fb802` 那段旧钩子喂进新尺 ⇒ 到达性断言逐字红 `AC#6 RED: ... re-entered 9 levels (cap 8) with the native exit called 0 times`（`nativeExitCalls=0 doorRounds=0 capTripped=true capDepth=9`）；甲③ ⇒ 绿，`nativeExitCalls=1 doorRounds=1 maxPostDepth=3 resolved=1`，帧里 `params[0]` 逐字＝`panel.ts:246` 信封、`dispatchRaw` 下游带 `requestId`/`method`/`source`。**四发定向突变全按预期红**（M1 中和 `w.Init` 那句／M2 旧形进产码格／M3 摘掉 `var native` 的捕获⇒再入 9 层、原生 0 次／M4 摘掉 `finally` 的归还⇒第二封变裸信封死于未绑槽），源码三行尺还原（md5 `3995d6bbb20ec251651544c872f5d359` 前后同、`git diff --exit-code`=0）。三条禁区现量：`git diff a885b532..HEAD -- frontend/`＝0、`-- internal/panel/`＝0、名册逐枚仍是那六枚、`m.dispatchRaw(` 调用点恰 1。门禁：`d22scan` rc=0 clean（`cmd/` 105→106＝本腿那枚新判据）、`go vet ./cmd/wisp/ ./internal/panel/` rc=0、`./internal/...` 起终逐名红集合 **DIFF-EMPTY**（已知 5 ＋ 间歇 `TestResolvePerCallBudget`）。⚠ 具名两笔：①整包 `cmd/wisp` 名册作差多出的那一枚 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 经单独复跑（r2 件 3/3 绿、基线件 2/2 绿）定性为**批次时序 flake，与本程无关**（该用例体内 `installPanelTransport|dispatchRaw|Eval|wispDispatch|bringUp` 命中 0）；②`35-r1` 那枚判据一字未动且照旧绿。⛔ **AC#6 仍不该翻**：本尺证的是"属性可覆写时甲③ 送达、旧形死循环"，而**真 WebView2 里 `chrome.webview.postMessage` 让不让覆写**这一支仍未取数——那一发连同 `bringUp` 真实载体的页面式 `postMessage` 到达，归 `35-v2`。逐字读数与推理全在 `.scratch/wisp/probes/35/r2/impl.md`＋`logs/`。⛔ 未 push。

**10-07 21:0x｜验收腿 `35-v2` 进展登记（只追加，⛔ 原句一字未动、⛔ 本腿不翻任何 AC 框）**：`-tags winlive` 真窗那一发**取到了**——载体＝`bringUp`（`panel_host_windows.go:318`）→ `installPanelTransport`（`:408`→`:693`→`w.Init(panelPostMessageForwardInit)` `:700`），泵＝库自己的 `WebView.Run()`（与 `panel_resident_windows.go:299` 同形，⛔ 本腿不自泵），页面上下文里调的是**属性** `window.chrome.webview.postMessage(<信封>)` 而不是门（⛔ 非 N6 那一形，判据全部由 Go 侧门记录给出）。读数：`pc-35v2-pagepost-1` 第 1 发即到达，`method="panel.mode.request" source="panel-composer" to="ask_every_step"` 逐字段正确，8.25s `--- PASS` rc=0，webview 进程树 0→0（machineNamed 29→29）；诊断道回带 `pc-35v2-DESC-OWN-WR-CF-HOOK-NOERR`（覆写后的属性＝自有·可写·可配置·钩子在场）。★控制组（`fb2fb802` 的甲①钩子，`go test -overlay` 换文件，⛔ 仓里一字未改）同尺同载体空载复跑 58.27s、门记录 `[]`、`--- FAIL` rc=1 ⇒ 这发判据对"甲③ vs 甲①"有辨别力。⚠两格本腿取不到并具名：覆写**之前**那份原生 descriptor（文档级脚本没有更早的插槽 ⇒ "可写"是被"到达"反推的）、以及原生出口丢 `this` 时到底抛不抛（覆写后拿不到原生引用）。★攻尺（`panel_transport_35r2_test.go` 1824 行）四发定向突变：`native(message)`（丢 this）与删 `typeof` 半支 **两发全 5/5 PASS rc=0＝这把尺的两面恒真**（模型的 `native` 桩 `:1238` 把 recv 参数丢弃、`:137 jsObject.set` 无 writability 概念，全文 0 处 `writable|configurable|defineProperty`）；`inside` 在同步 hop 之前归还 ⇒ 3 枚红（`capTripped=true maxPostDepth=9 nativeExitCalls=0`）、`finally { inside = inside; }` ⇒ `TestShapeA3SecondPagePostStillDelivers` 红 ⇒ A682 要求的那一面（再入/时序）**确有牙**。三处对拉自取库原文坐实：`chromium.go:112` 逐字、Bind 的桩同步 `window._rpc[seq]`→`external.invoke`→返回 Promise、`msgcb` 未绑名走 `callbinding return nil, nil` ⇒ **resolve** 支（票 `:58-59` 那句 `reject` 是过期句，`A682` 的更正成立）。禁区与渲染串：`frontend/`、`internal/panel/` 对 `a7f9781c..HEAD` 均空，`m.dispatchRaw(` 计数 1，名册六枚逐字（`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`/`config.get`/`config.set`），渲染出的 JS 里 `%[1]s` 三处全换成 `wispDispatch`、整串无游离 `%`。门禁：d22scan rc=0、`go vet ./cmd/wisp/ ./internal/panel/` rc=0、`go test ./internal/...` rc=1 且红名册恰＝已知 5 枚（`TestResolvePerCallBudget` 本程未红）。证据＝`.scratch/wisp/probes/35/v2/`（`verdict.md`＋logs 20 枚，每件自带 rc）。⛔ 未 push、⛔ 未动 AC 框、新增文件只有 `cmd/wisp/panel_transport_live_35v2_windows_test.go` 一枚（`//go:build windows && winlive`，标注验收台件）。

### 10-07 20:5x 编排者收 `35-v2` ⇒ ★**翻 `:52` 那枚框（Transport agreement）**＋★**我自己一处编号错误具名作废**（我这几轮写的"AC#6/AC#7"是**序号口径**，盘上这节现量 8 枚框：`:42/:44/:45/:47/:49/:51/:52/:63` ⇒ `:52` 其实是第 **7** 枚、`:63` 是第 **8** 枚；⛔ 从今起本票只用"行号＋名字"引框，A681/A682/A683 与票面那三段注里的"AC#6/AC#7"字样**按此改读**，原文不改）＋⛔ `:63` 仍不翻（(d) 三支 0/3）

**凭据（非实现者腿 `35-v2`，两笔 `8b75816a`／`3eeff4a7`，24 枚件）**：真窗那一发**取到了**，载体＝产品路径 `bringUp`→`installPanelTransport`→`w.Init`，泵＝库 `Run()`，`LockOSThread` 不解锁、一次一窗。绿那发（`logs/live-run1.txt`，`rc=0`，8.25s）：页面自己调 `window.chrome.webview.postMessage(<composer 信封>)`（台件 `:57`／`:64`，⛔ 不是 `w.Eval("window.wispDispatch(…)")` 那枚 N6），**第 1 发就到 Go 门**，`method="panel.mode.request" requestId="pc-35v2-pagepost-1" source="panel-composer" to="ask_every_step"` 逐字段正确；诊断道回带 `OWN-WR-CF-HOOK-NOERR` ⇒ 真 WebView2 里 `chrome.webview.postMessage` 是**自有、可写、可配置**的属性且钩子挂载无错。红那发＝控制组（`logs/live-run-legacy-control.txt`，`rc=1`，58.39s）＝同一台件同一载体、只把钩子换成 `fb2fb802` 那枚甲①（`go test -overlay`，⛔ 仓里一字未改）⇒ 门记录 `[]`、四发都不到 ⇒ **这发判据两侧分得开**。webview 进程树 0→0 零残留。禁区三面我自己也复量过：`frontend/` 与 `internal/panel/` 对 `a7f9781c..HEAD` 均空、`bridge.go` blob `bebe8e70` 未动、`m.dispatchRaw(` 恰 1 枚、渲染串里 `%[1]s` 三处全换成 `wispDispatch`、无游离 `%`。⇒ **`:52` 成立并翻勾**；⚠ 但这枚的牙今天**只活在 `-tags winlive` 的本机台件里**（winlive 无 CI 载体）⇒ 我把"补 CI 可达替身或明写〔仅本机可量〕并逐波复跑"连同 W1 那两面恒真一起立成**新框**（本节上面那枚 `- [ ]`），⛔ 不许把 `:52` 的翻勾读成"这格在 CI 有牙"。

**`35-v2` 攻那把 1824 行行为尺的结果**：**M-A**（`native.call(cw, message)` → 裸 `native(message)`，丢 `this`）与 **M-B**（删 `typeof … !== "function"` 半支）两发**都 5/5 PASS＝恒真面**，根因在**夹具不在解释器**——原生桩 `:1238` 把 receiver 丢了、`jsObject.set`（`:137`）没有可写性概念、全文 0 处 `writable|configurable|defineProperty`。**再入／时序那一面确有牙**：M-C1（在 hop 之前归还 `inside`）三枚红、M-C2（`finally { inside = inside }`）红在 `SecondPagePost`。另具名一枚偏严：`onGet` 把"读属性"也计一层 ⇒ 再入帽 8 不可移植。⇒ 第 127 条预言的形状兑现了：**"读字符串"换成"手写行为模拟"之后，新尺自己立刻成了新的假绿面**；写腿自报的"四发突变都红"证明不了覆盖面，只有非实现者自己造反形才测得出。

**另外三笔小账**：①名册六枚的真名＝四枚 composer **＋ `config.get`／`config.set`**（⛔ 不带 `panel.` 前缀——按 `panel.` grep 只得 4 枚，与 10-02 那条"config 两名掉在名册尺外"同源）；②W3：它顶回 `35-r2` 那对"基线绿＋r2 侧红"**证不了 flake**（凭据形状不成立），但用**同一 exe 同一棵树整包两发**（一发含 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`、一发不含）＋单跑 3/3 绿坐实"间歇"，我按"间歇成立、产码轴未隔离具名"收；③它台件打印的 `cold bring-up=-1.000ms` 判为**它自己 bringUp 期间不起泵的伪影**，⛔ 不开新格、不记产品账（本条即处置）；④`3eeff4a7` 提交信息里把 `a7f9781c` 打成 `a7f981c`（腿自报），按纪律不 amend。

**10-07 21:2x｜写腿 `35-r3` 交件登记（票 `:63` 的 (d) 三支守卫红；⛔ 只追加，本腿一个字没动 AC 框、未翻任何勾；引框只用行号＋名字，遵 `A684` §2）**：新增 `cmd/wisp/panel_inbound_guards_35r3_test.go`（`//go:build windows`）一枚，⛔ 不改产码（`git diff HEAD -- cmd/wisp/panel_host_windows.go internal/panel/` 空）、⛔ 不碰 `frontend/**`、⛔ 不改 `panel_transport_35r2_test.go` 既有断言。三支入口一律走页面那条边：复用 r2 同包 `main` 的 `fakeDoc35r2` 行为尺，`installPanelTransport`→`openDocument()`（逐字跑 chromium.go:112 垫片＋webview.go:462-478 Bind 桩＋产码 `panelPostMessageForwardInit` 钩子串本身）→`pagePost(envelope)`＝页面 `window.chrome.webview.postMessage(<信封>)`（`panel.ts:179/:218` 真调的那句）；⛔ 未直调 `dispatchRaw`、⛔ 未用 `w.Eval("window.wispDispatch(…)")`（N6）。三枚用例（绿那发 rc=0）：
- **(d)-1** `TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge`＝页面今天真发的 `panel.approval.request` 原文（`panel.ts:179-184` `{method,correlationId,outcome}`），断 `doorRounds==1`＋回复逐字含 **`不是面板 composer 通路的能力入口`**（`:133`）、不含另两支关键词、且 `spy.reqs==0`（守卫在路由前拒）。⚠ 本格新料≠"它被拒"（既有钉 `internal/panel/l2_grant_boundary_test.go:1964`/`:1251-1267` 在 internal/panel 包内直调名册函数就已拒它、走不到这条传输边）＝**经页面那条边到达＋点名是名册那一道**，⛔ 未拿既有钉抵账。
- **(d)-2** `TestInboundSourceGuardRefusesForeignSourceOnPageEdge`＝名册内 `panel.mode.request`＋`source:"panel-composer-x"`，断回复含 **`按伪造/串台拒绝`**（`:136-137`）。
- **(d)-3** `TestInboundRequestIDGuardRefusesMissingIDOnPageEdge`＝同枚去掉 `requestId`、source 合法，断回复含 **`缺少 requestId`**（`:140-141`）。
**恒真性自查（三支各点各的，`go test -overlay` 于 `/d/tmp/wisp35r3` 各改一行、⛔ 仓里一字未改，diff 见 `logs/overlay-diffs.txt`）**：中和名册判（`:132`→`if false && !knownComposerMethod(...)`）⇒ **只有 (d)-1 红**、(d)-2/(d)-3 绿；中和来源判（`:135`→`if false && TrimSpace(r.Source)!=...`）⇒ **只有 (d)-2 红**、(d)-1/(d)-3 绿；中和 requestId 判（`:139`→`if false && TrimSpace(r.RequestID)==""`）⇒ **只有 (d)-3 红**、(d)-1/(d)-2 绿。三支今天**分得开**（逐字红句见 `logs/mut-roster.txt`／`mut-source.txt`／`mut-reqid.txt`）。⚠ 具名恒真前提：三支只依赖"不透明串经转发钩子→门→Handle→ParseComposerRequest→回复串"这条已复量有牙的投递面（`doorRounds`/`lastReply`），**不落进** `A684` 那两面恒真（丢 receiver／无可写性概念），但其"到达"以"真 `postMessage` 允许被钩子覆写"为前提——那一面只有 `-tags winlive` 真窗取过数、正落在 M-B 盲区里，故本三支证的是守卫次第＋各支文案这一维，⛔ 不得读成"传输可覆写性"证据（那一维另程处置，票 `:75`）。门禁：`d22scan` rc=0 clean（`cmd/` 108 枚含本判据，符号全在注释内）、`go vet ./cmd/wisp/ ./internal/panel/` rc=0、`go test -count=1 ./internal/...` rc=1 且逐名红名册＝起手已知 5 枚集合作差空（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestC21TableColourRowsMatchTokensCSS`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`，间歇门不计账）；cmd/wisp 走 harness（`go test -c`＋三枚 sherpa DLL 拷 exe 旁＋CWD=`cmd/wisp`），本腿 `-test.run TestInbound` 那发 3/3 绿。⛔ 未 push、未开真窗。逐字读数与推理全在 `.scratch/wisp/probes/35/r3/impl.md`＋`logs/`。

**10-07 21:5x｜验收腿 `35-v3` 交件登记（⛔ 只追加，本腿一个字没动 AC 框、⛔ 未翻任何勾；引框按 `A684` §2 只用行号＋名字）**：本腿＝非实现者，验 (d) 那一支。八发突变**全自造**（`D:/tmp/wisp35v3` 副本＋`go test -c -overlay`，⛔ 未复用 r3 那三份副本、⛔ 仓里一字未改；四枚守卫/产码件在八发跑完后 `git hash-object` 与 HEAD blob 逐枚相同＝见 `logs/repo-untouched-after-overlays.txt`）。**V1 成立**：中和 `:132`／`:135`／`:139` 三道判的条件本身（每份恰 1 行）⇒ 各**只有它那枚红**、另两枚绿（rc=1×3）。**V2 成立（本程最值钱的一发，⛔ 无恒真面可报）**：文案合并两形都咬得住——(甲) 三句都换成 `面板请求被拒绝` ⇒ **3/3 红**（走"本支关键词必须在"那条 prong）；(乙) 三句都换成**同一句、但把三枚关键词全写进那一句** ⇒ **3/3 红**（走"另两支关键词必须不在"那条 prong）。⇒ 简报担心的"这组判据只是 in-string 找关键词、看不见文案合并"**没有兑现**：它有两枚各自独立的 prong。**V3**：①摘掉产码转发钩子（`w.Init(panelPostMessageForwardInit)`→`_ = panelPostMessageForwardInit`）⇒ **3/3 红、`doorRounds=0`**＝`:63` 的 ★Falsifiability 由非实现者复现；②`spy.reqs==0` 不是恒真（正控＝守卫不动、只把 (d)-3 信封补成合法 ⇒ 红句 `the mode handler ran 1 time(s)`），但它在所有可造形状里与 mine-absent **同发**＝佐证级、不是独立尺；③★**具名两面看不见**：`native.call(cw, message)`→裸 `native(message)`（`A684` 的 M-A 形）本腿实测 **3/3 绿 rc=0**；夹具 `writable|configurable|defineProperty` 现量 **0 处**（M-B 结构性盲区，A684 已造反证）。⇒ ⚠ **这三枚绿不许被读成"传输在真浏览器里成立"**，那一格仍只由 `:52` 的 winlive 凭据撑着（票 `:75`）。**V4**：禁区现量＝`git diff --stat 1d70fb2b..HEAD -- cmd/wisp/panel_host_windows.go internal/panel/ frontend/` **空**、`bridge.go` blob `bebe8e70`／r2 尺 blob `7185ab56` 逐字节未变、`m.dispatchRaw(` 产码 1 枚、r3 件里 `dispatchRaw`/`.Eval(` 命中**全在注释**（`:24`/`:43`，调用形状 0）；三枚关键词在 `bridge.go` 各出现**恰 1 次**且⛔ 互不为子串（Python 对拉 total=3）。★**一处具名顶回（腿对、原文需改）**：`A685` §2 末句说中和名册判那发 (d)-1 "撞上 requestId 判"——盘上逐字 reply 是 **`来源 "" 不是 "panel-composer"，按伪造/串台拒绝（requestId=""）`**＝**来源判 `:135`/`:136-137`** 接住的（impl.md §② 的表述才对）。★**另一处自我更正并具名**：本腿原猜"既有钉看不见名册被中和"＝**错**——同一份 M-R overlay 跑 `TestGrantWireShapesAreRefusedAtTheDoor` **rc=1**（它五枚线形里 4 枚带 `requestId`＋`source`，中和后一路 accept），未突变对照 rc=0。⇒ 不抵账仍成立，但**理由只能写这两条**：既有钉 ①不经传输、②只判 `err != nil` 说不出**哪一道**答的（对页面真发的那枚 `{method,correlationId,outcome}`，名册判被中和它照样绿，而 (d)-1 在这一发是红的）。**V5**：`d22scan` rc=0 clean、`go vet ./cmd/wisp/ ./internal/panel/` rc=0、`go test -count=1 ./internal/...` rc=1 且红名册与既有 5 枚**双向作差皆空**（known=5 got=5），间歇门 `TestResolvePerCallBudget`/`TestTicket223…` 在 internal 那发 0 命中不计账；`cmd/wisp` 整包（本腿独占）截数点 **250 PASS／5 FAIL**＝`A684` §3-W3 那族窗口依赖既有红（`TestPanelHostRealWindowHopAndLifecycle`/`TestAC4FocusReturnToPriorWindowGap33r5`/`TestAC13…`/`TestAC14…`×2），⛔ 不记本程账；**最终读数 rc=1／338 PASS／5 FAIL／1 SKIP**（三支在整包同发 3 PASS；间歇门 `TestTicket223…` 这一发 PASS；逐名存 `logs/cmd-wisp-full-summary.txt`，原始 1500+ 行按输出预算留 `D:/tmp/wisp35v3/cmd-wisp-full.txt`、⛔ 不入库亦未删）；三支基线本腿自跑 **3 PASS／0 FAIL／rc=0**。证据 `.scratch/wisp/probes/35/v3/{verdict.md,logs/*}`（21 枚件，⛔ 无 `.out`、0 字节件＝0、每件自落 rc 行）。**本腿结论＝`:63` 的 (d) 我这格无异议、可以翻，⛔ 翻勾由编排者做**；射程一句话＝这三支咬得住**守卫中和／文案合并／转发整段摘除**三种形状，咬不住**夹具看不见的 JS 语义层**（`this` 绑定、属性可写性）。⛔ 未 push、未开真窗。

### 10-07 21:4x 编排者收 `35-v3`（非实现者，`2078ec14`／`8d496dd0`，24 枚件）⇒ ★**翻本票 `:63` 那枚框**（现量 9 框＝**2 勾／7 未勾**，⛔ 不加 `-done`）＋★**更正我自己 `A685` §2 末句的一处错**（我先在盘上复跑那一把才落笔）＋★**勾旁必须钉着的那句限制**

**凭据（逐条我自己复读过原文）**：(a)(b)(c) 半格由 `35-v1`/`35-v2` 复认、`:52` 的真窗那一发在 `A684` 已入账；今天新取的是 **(d) 三支**：
- **V1 成立**——`35-v3` 自造三份**恰一行**的 overlay 副本（`/d/tmp/wisp35v3`，跑完 `git hash-object` 与 HEAD blob **逐枚相同**＝仓里一字未改）中和 `bridge.go:132`／`:135`／`:139` 三道判的**条件本身**（不是删文案）⇒ 每发**恰 1 枚红、另两枚绿**。
- **V2 成立，且我担心的那面没兑现**——我派单专门要求它做一发 `35-r3` 没做过的反形：**把三句拒绝语合并成同一句**。两形都咬住了：三句都换成泛化的 `面板请求被拒绝` ⇒ **3/3 红**（走"本支关键词必须在"）；三句都换成**同一句但含全部三枚关键词** ⇒ **3/3 红**（走"另两支关键词必须不在"）。⇒ 这组判据**不是只会 in-string 找关键词**，"不许把几桩病压成一格"这条今天有机读的牙。
- **V3① 由非实现者复现了这格自己的可变真性**——摘掉产码那句 `w.Init(panelPostMessageForwardInit)` ⇒ **3/3 红、`doorRounds=0`**。⇒ 这三支确实骑在传输那条边上，不是绕进去的。
- ★**勾旁钉着的限制（腿要求我写，我认）**：V3③ 实测 `native.call(cw, msg)` → 裸 `native(msg)`（`A684` 的 **M-A** 形）**3/3 照绿**，夹具里 `writable|configurable|defineProperty` 现量 **0 处**（**M-B**）⇒ **这三支的绿不含"真浏览器里那条转发成立"**；那一格只由 `:52` 的真窗凭据撑，JS 语义层（`this`／可写性）归本票 `:75` 那枚新框。
- **V4 成立**：三方禁区差 `1d70fb2b..HEAD` 全空、`bridge.go` blob `bebe8e70` 与 r2 尺 blob `7185ab56` 未变、`m.dispatchRaw(` 产码 1 枚、r3 件里 `dispatchRaw`／`.Eval` 命中**全在注释**（`:24`／`:43`）；三枚关键词各恰 1 次且互不为子串。**V5**：d22scan rc=0、vet rc=0、`./internal/...` 红名册与既有 5 枚**双向作差皆空**（间歇门本程 0 命中）；`cmd/wisp` 整包 rc=1／338 PASS／5 FAIL／1 SKIP＝窗口依赖的既有红族，⛔ 不记本程账。

**★我自己那处错（先复跑再落笔，第 122 条）**：`A685` §2 末句我写"信封往前走了一步**撞上 requestId 判**"——盘上逐字 reply（`probes/35/r3/logs/mut-roster.txt:5`）是 **`来源 "" 不是 "panel-composer"，按伪造/串台拒绝`** ⇒ 接住它的是**来源判 `:135`/`:136-137`**，不是 requestId 判；我那句是把回显前缀里那个 `(无 requestId)`（只是显示"这封信没带单号"）读成了"报错的那道门"。`impl.md` §② 的原句才是对的。**`A685` 原文不改**（不改已提交的历史），本条即具名更正；这条也正是第 111 条那一族：**看长行里的关键词位置，别按关键词的字面猜因果**。

**⚠ 一处口径被腿收紧（采它）**：它原本要论证"既有钉 `l2_grant_boundary_test.go` 看不见名册被中和"，实测**不成立**（M-R 下 `TestGrantWireShapesAreRefusedAtTheDoor` rc=1，4/5 线形一路 accept）⇒ "不抵账"这个结论仍成立，但**理由只能写两条**：既有钉**不经传输**、且它**说不出是哪一道门答的**。以后引这格区别⛔ 不许再说成"既有钉看不见名册被拆"。

### 10-08 08:2x 编排者收写腿 `35-r4` 的**盘上遗产**（代提笔 `2fc5f5c9`，32 枚件）⇒ ★**这格今天有了自己的牙读数**，但⛔ **本框 `:75` 一个字没动、0 枚勾翻**（实现方自证不算凭据，非实现者 `35-v4` 待派）

**为什么是编排者代提**：`35-r4` 死于**每日额度**（39 次调用／442 万 token／35 分 10 秒，最后一件落盘 10-07 22:36），⛔ 不是任务失败、不是服务端掐大上下文。它把夹具做完了并自己跑绿（`logs/r4-run.txt`＝12 枚 0 FAIL），差的是**交付件与 commit**；那枚 0 字节的 `countershape-driver.txt` 就是它死在"正要跑反形"那一刻的痕迹（我没有把 0 字节件入库，只在这里具名）。

**落了什么（只动测试夹具这一面，`cmd/wisp/panel_transport_35r2_test.go` ＋245／−3；产码零改动）**：
- `:137` `jsObject` 加 `unwritable map[string]bool`；`:149-160` `markNonWritable`（**Go 侧台件**）；`:163-166` `set` 读表后**静默拒绝覆写**。
- `:1285-1291` 原生出口桩**改成认 receiver**，丢了就 `panic("TypeError: Illegal invocation: chrome.webview.postMessage lost its receiver - called with …")`。
- `:1926` 载具 `runShippedHookInOneWorld35r4` 装的是**生产接线**（`installPanelTransport`：绑门＋把产码钩子字符串当值读入），不是另造一份钩子。
- 三枚具名用例：`:1946` M-A／`:1976` M-B／`:2026` 门不在时回落原生。

**★四发反形是我代跑的（腿那一版四发全部没跑成）**：它的 `logs/leg/mut-summary.txt` 四节都写 `PASS=0 FAIL=0 buildrc=1 runrc=127`＝**编译就没成**，不是"突变没红"。根因一行：台件 `mk_overlay` 里**被替换侧路径走了 `cygpath -w`**，Windows 反斜杠直接进 JSON ⇒ `go: parsing overlay JSON: invalid escape sequence \w in string`。我的 `scripts/orchestrator-countershape.sh` **唯一改动＝那两条路径也走 `-m`**，突变源／`-test.run` 名册／断言沿用腿的原件（`logs/orch/orch-mut-landing.txt` 现量四枚仍各 `changed-line-count=2`＝恰一行被换）。颜色（`logs/orch/orch-mut-summary.txt`，buildrc 全 0）：

| 反形（逐字改了哪一行） | 颜色 |
|---|---|
| `ma-hook`：产码 `native.call(cw, message)` → `native(message)` | **9 FAIL／3 PASS**，红名含 `TestForwardingHookMustNotLoseTheNativeExitReceiver` |
| `ma-detector-gone`：夹具 `if recv != cw {` → `if false {`（**与 ma-hook 同发**） | **12 PASS／0 FAIL**＝摘掉探测器后那 9 枚全转绿 |
| `mb-set-gone`：夹具 `exists && o.unwritable[name] {` → `exists && false {` | **1 FAIL／11 PASS**＝恰 `TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit` |
| `mb-typeof-gone`：产码 `if (inside || typeof window.%[1]s !== "function")` → `if (inside)` | **1 FAIL／11 PASS**＝恰 `TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent` |

⇒ 我把两处形状差别**具名写出不抹平**：M-A 是**级联红**（探测器自己 panic ⇒ 凡走原生出口的一起红），它的归因**只能靠** `ma-detector-gone` 那一对（摘探测器 ⇒ 全绿）来支撑，⛔ 不许读成"红了 9 枚所以更凶"；M-B 与 typeof 两发各**恰红 1 枚且名字唯一**，是更强的形。

**门（我现跑）**：`go vet ./cmd/wisp/` **rc=0**；`cd tools/d22scan && go run . -root ../../` **rc=0 clean**（ban #8 覆盖 `cmd/` 108 枚 Go 文件含注释与 `_test.go`）；格式那格**具名带着既有红交件**——`gofmt -l` 会列这枚文件，但**同样 4 行漂移在 HEAD 版本里就在**（HEAD `:775-776`／工作树 `:807-808`，只是被插入顶位移了；CR 计数两版皆 0）⇒ 本笔**不新增**未格式化面，那 4 行是 `35-r2` 落下的，归票 275 那族格式化欠账，⛔ 不在本程顺手改、⛔ 不许为变绿放宽任何断言（`ci.yml:209` 那条注释自己记着 run `37406757402` step 8 gofmt `[failure]` ⇒ step 9 `[skipped]`）。

**⛔ 这格不许被读成什么（腿自己的话逐字在 `impl.md` §3，这里只摘要）**：它买的是"丢了 `this`、或赋值被静默拒绝的钩子，不再被 certify 成送达"；⛔ 它**不买**"WebView2 真会绑 receiver／真让页面换掉 `postMessage`"——那一格在 `:52` 的真窗凭据，两格⛔ 不许互相借光（`A684` §3／`A686` §4 同一条）。`Object.defineProperty` 没做成页面可见方法，`getOwnPropertyDescriptor`／`configurable`／`delete`／accessor 属性**全部未建模**（`impl.md` §4 列全）⇒ 钩子若写成"先 `delete` 再赋值"这类形状，这枚尺今天**看不见**。

**留给 `35-v4`（非实现者）的四问，⛔ 我不代答**：①级联红能不能算 M-A 的牙、那对摘探测器读数够不够；②两枚"恰红 1 枚"的名字唯一性要**自造混文案反形**去顶（本票 `:63` 的 V2 就是这么做的）；③§4 五样未建模里有没有哪样能让这三枚在产码坏掉时照样绿（＝`A684` 两面之外的第三面）；④整包红名册（我这版与 `35-v3` 记的那五枚逐名作差，见台账 `A689`）里有没有新名字要归本程账。
