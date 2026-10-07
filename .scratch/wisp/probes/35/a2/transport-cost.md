# 35-a2 — Transport-cost survey (read-only leg, ticket 35 / C17 panel bridge)

Leg: `35-a2` (read-only cost census). Shape NOT chosen here (选形归编排者).
Scope = ticket 35's new AC frame "Transport agreement page↔host" (uncheckable by this leg, ⛔ not flipped).
All command transcripts live in `logs/` beside this file; every line number below is re-measured at
HEAD `28872c9a` on branch `dev`, and every page-side number states which tree it came from.

---

## §0 起手锚 — see the appended section at the end of this file (anchor gate landed first, commit `ccb9f0df`)

(§0 content is kept verbatim in the "## §0 起手锚" heading below — the commit ordering is:
`ccb9f0df` = §0 + logs/anchor-reads.txt, landed before any long-running command.)

---

## §1 规格射程 — is the transport written down anywhere? (question ①)

**Ruler used.** Per this repository's rule #111 ("judge whether a sentence exists by what the sentence
SAYS, not by one symbol name"), I scanned two ways and read every hit's FULL line:
- Ruler A — the four candidate transport spellings, verbatim: `postMessage`, `external.invoke`,
  `AddHostObjectToScript` / `HostObject`, `JSON-?RPC` / `jsonrpc`, plus the library-side names
  `CoreWebView2`, `WebMessageReceived`, `EvaluateScript`, `ExecuteScript`, `CallDevToolsProtocol`,
  `wispDispatch`, `window.wisp`.
- Ruler B — meaning-level Chinese/English phrasings a spec would use to *prescribe* a transport:
  `传输`／`通道`／`桥接`／`信封`／`序列化`／`协议`／`绑定`／`注入宿主`／`前端…调用`／`宿主…注入`.

**Transcripts:** `logs/q1-spec-scan-hits.txt`, `logs/q1-transport-spelling-scan.txt`.

### 1.1 Answer, split three ways as asked

**盘上查无此句 — for the authoritative texts.** ⛔ Not "I did not find it": this is a negative result
from a scoped re-run whose command and output are in `logs/q1-transport-spelling-scan.txt`:

```
$ grep -rnE 'postMessage|external\.invoke|AddHostObjectToScript|HostObject|JSON-?RPC|jsonrpc|CoreWebView2|wispDispatch|window\.wisp|WebMessageReceived|WebView2' docs/PLAN.md docs/specs/
```
Every surviving hit is carried by the single token **`WebView2`** (platform/runtime/library), and those
lines are reproduced in full in §1.3 below. Zero hits for `postMessage`, `external.invoke`,
`AddHostObjectToScript`, `HostObject`, `JSON-RPC`, `jsonrpc`, `CoreWebView2`, `wispDispatch`,
`window.wisp`, `WebMessageReceived` **anywhere in `docs/PLAN.md` or `docs/specs/**`**.
Where those spellings DO live: `docs/evidence/s1/**` (adjudication tables) and
`docs/reports/pending-and-issues.md` (the ledger) — i.e. **leg/orchestrator prose, not spec**.
`external.invoke` and `JSON-RPC`/`jsonrpc` = **0 hits in the whole `docs/` tree**, evidence included.

**规格只描述能力，不描述传输 — this is the actual situation.** Everything the frozen texts say about
this bridge is capability-and-shape, never wire-and-API:

- `docs/specs/SPEC-08-ui-ball-panel.md:156-159` (the C17 section, whole line read):
  `156: ### 5.2 C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】`
  `158: `invoke(method, args) → result` + Go→前端事件推送（流式结果/审批请求/任务状态），回复按`
  `159: correlationId 路由；每方法标注 capability 与是否需原生侧授权；未列出方法名 → 拒绝并记日志。`
  → names a *call signature* (`invoke(method,args)→result`) and a *routing key* (correlationId).
  It never says which WebView2 API carries it.
- `docs/PLAN.md:1367` (C17, the frozen contract row, whole line read):
  `| **C17** | **`PanelBridge`** | 前端↔Go 双向通道：`invoke(method, args) → result` + Go→前端事件推送
  （流式结果、审批请求、任务状态）。**回复必须按 correlationId 路由**；前端**必须无状态**（WebView 销毁后
  一切从 Go 侧重读） | D29, D31 |`
  → "双向通道" is a capability word here, not a transport choice.
- `docs/PLAN.md:2968-2973` (D39's C17 deepening, "补四项", whole block read): ① method whitelist
  ② capability + native-authorization flag per method ③ correlationId routing + backpressure/merge
  ④ forced statelessness via `panel.resync` on every `show`. ⛔ No transport clause in any of the four.
- `docs/PLAN.md:2434`: `**配套**：C17 `PanelBridge` 必须给出**方法白名单**（只有列出的方法可被前端调用）`
  → "被前端调用" again names the caller, not the mechanism.
- `docs/specs/SPEC-08-ui-ball-panel.md:235-238` (test decisions for C17) pin *behaviours to test*
  (correlationId routing, whitelist refusal, resync-equality-after-reopen, `approval.decide({allow:true})`
  server-side rejection) — ⛔ not a wire format.

**规格里写死了的东西（there IS one, and it is narrower than the ticket implies）— the library and the
resource-loading path are fixed, the message pipe is not:**
- `docs/specs/SPEC-08-ui-ball-panel.md:145`: `- `jchv/go-webview2`（MIT，纯 Go 无 cgo）；**单例
  `PanelManager`：一会话至多一个 WebView2 窗口，隐藏而非销毁**——…` → **the dependency is named in the
  spec** (this matters for shape 乙 in §2.2).
- `docs/specs/SPEC-08-ui-ball-panel.md:148-149`: `AddWebResourceRequestedFilter` from `embed.FS`,
  **不起 localhost HTTP 服务** → the *resource* channel IS prescribed (and `docs/PLAN.md:1008` names the
  library again). The *message* channel is not.
- `docs/specs/SPEC-08-ui-ball-panel.md:154` / `PLAN.md:1377` (C27): "多任务共用单窗口，按 correlationId
  分区渲染" — again a rendering/routing rule, not a transport.

### 1.2 What the spec does NOT let me say

⛔ I am not claiming the spec *permits* any transport either. Two constraints exist and bite all three
shapes, and they are contract/禁止-list constraints rather than spec sentences:
- `docs/specs/SPEC-08-ui-ball-panel.md:148-149` forbids a localhost HTTP server as the channel
  (also enforced by a test — §3 nail N5).
- `docs/PLAN.md:1374` (C24) has the "未列出的宿主入口即判契约违规" idea for the *plugin* JS surface;
  the panel's binding-name surface is guarded by `inbound_roster_253_test.go` (§3 nail N1).

### 1.3 The remaining `WebView2` lines in PLAN.md, read in full so nobody has to re-guess them
(all reproduced in `logs/q1-spec-scan-hits.txt` + the §1 batch read)
`:46` Tauri/WebView2 resident-RSS rebuttal · `:1008` dependency table row (`jchv/go-webview2`, MIT,
pure Go, "完整 WebView2 COM 绑定") · `:1015`/:`:1021` cold-start budget and the "停更 7 个月可接受，
它是对微软版本化稳定 COM ABI 的薄绑定" judgement · `:1339`,`:1499`,`:1563`,`:3049` Runtime-missing
degradation · `:1568`,`:1936`,`:1955`,`:2043`,`:3227` P11 cold-start / Environment-sharing risk ·
`:1593`,`:1793`,`:1813`,`:2352` SLO/Job-object/`C27` single-window · `:2820`,`:2824` (D38) STA-thread
ownership of the WebView2 window · `:2227`,`:2244`,`:2259` ball/panel visuals.
⇒ **None of these is a transport clause.** `docs/specs/SPEC-01:70` only repeats
`AddWebResourceRequestedFilter`; `SPEC-09:45/46/53` only repeats Runtime/leak/degradation.

### 1.4 One-line verdict for ①
**规格只描述能力不描述传输；传输方式在 `docs/PLAN.md` 与 `docs/specs/**` 里查无此句（三个候选拼写全部 0 命中，
命中都在 evidence/ledger 里）。唯一被规格写死的是依赖名 `jchv/go-webview2`（SPEC-08:145）与资源通道
`AddWebResourceRequestedFilter` + 禁 localhost HTTP（SPEC-08:148-149）。** ⛔ 我没有、也不会改 `docs/**` 一字。

---

## §2 三形代价（⛔ 只报代价，不选形 — question ②）

### 2.0 The break being costed, re-confirmed by this leg (not copied)
Page sender `frontend/src/lib/panel.ts:179`/`:218` → `bridge.postMessage(JSON.stringify({method,
correlationId|requestId, …}))` with `bridge = window.wispBridge ?? window.chrome?.webview` (`:152-155`).
WebView2's `postMessage` lands in `pkg/edge/chromium.go:233` `MessageReceived` → fed to
`MessageCallback` = `webview.go:139` `msgcb` → `json.Unmarshal` into `rpcMessage{ID, Method, Params}`
→ the page envelope has **no `id`/no `params` but is valid JSON**, so unmarshal succeeds with `id=0`
and `method="panel.approval.request"`/`"panel.composer.*"`, none of which is a *binding* name
(bindings present: `wispDispatch` `panel_host_windows.go:80/:405`, `wispProbeRT` `:668`) →
`callbinding` errors → `webview.go:148-150` `Dispatch(Eval("window._rpc[0].reject(…)"))` into a page
with no `window._rpc` slot ⇒ silent TypeError. ⛔ `dispatchRaw` (`:630`) is never entered, so no Go log.
Independently, `chromium.go:242-245` echoes the same string back with `PostWebMessageAsString` and the
page has no listener (dev `frontend/src`: `addEventListener` 4 hits, none `message`;
`window._rpc`/`onmessage`/`WebMessageReceived` = 0).

### 2.1 甲 — `w.Init(<shim>)` from Go, forwarding the page's existing `postMessage` into `wispDispatch`

**Files and lines it would touch (Go side only, all in `cmd/wisp`):**
- `cmd/wisp/panel_host_windows.go` — one new statement beside the existing `Bind` block (`:405-412`);
  `Init` is a method on the same `webview2.WebView` value `w` created at `:386` and stored at `:415`
  (interface method set confirmed: `Init(js string)` is exported, `common.go:68`).
  ⇒ cost = **1 file, ~1 statement + a JS-literal constant**.
- No page file. No `bridge.go`. No `go.mod`. ⇒ **the C17 roster is untouched** (see the whitelist
  question below).
- The receipt half: `window._rpc` must exist *before* `msgcb` ever evaluates `window._rpc[0].reject(…)`.
  `Bind`'s own injected script (`webview.go:462-479`) already creates
  `var RPC = window._rpc = (window._rpc || {nextSeq: 1})`, so a shim injected via `Init` that references
  `window._rpc` needs either to run after Bind's script (same mechanism, order not guaranteed by the
  library) or to create the slot itself. ⇒ **cost = the shim must be defensive about `window._rpc`**
  (a guard `window._rpc = window._rpc || {nextSeq:1}` inside the shim), ⛔ 不能假定 Bind 注入的脚本先跑。

**`Init` timing constraint, measured rather than assumed:**
- `edge.Chromium.Init` (`pkg/edge/chromium.go:130-136`) is `AddScriptToExecuteOnDocumentCreated` —
  a **document-creation hook, not an eval**. It therefore runs on every document/navigation created
  *after* the call and is naturally ordered before any `SetHtml`-served page's script executes,
  provided the `Init` call happens before the page is loaded.
- Live call order in this host today: `NewWithOptions` `:386` → `Bind(panelDispatchBinding,…)` `:405`
  (which itself calls `Init` `webview.go:462`) → `m.w = w` `:415` → `SetHtml` at **three** sites:
  `:449` (`serveNotBuiltNoticeLocked`), `:468` (`serveEntry`), `:681` (round-trip probe page).
  ⇒ **甲's `Init` must sit between `:386` and the first `SetHtml`**, i.e. beside `:405`. ⚠ Because
  `AddScriptToExecuteOnDocumentCreated` is *document-creation*-scoped, an `Init` issued after
  `SetHtml` would not retroactively reach the already-created document — so "Init 之前还是之后" has a
  real answer only for the first `SetHtml`; the later `SetHtml` (`serveEntry` on a re-show) each create
  a new document and *will* see it. ⇒ cost note: **the shim must be re-`Init`-able and idempotent**.

**Does this JS injection bypass C17's whitelist? (the load-bearing question)**
- The shim forwards the **raw envelope string** into the *already-bound* `wispDispatch(raw)`, whose Go
  body is `dispatchRaw(ctx, raw)` (`:405-407`, `:630`) → `panel.ParseComposerRequest` →
  `knownComposerMethod` (`internal/panel/bridge.go:146`, its case at `:148` lists all six) →
  the roster of 6 (the four `panel.*` at `bridge.go:42-45` + `MethodConfigGet`/`MethodConfigSet`
  at `bridge.go:66-67`). ⛔ It **forwards raw**; it does not add a *method* name. The whitelist decision is made by code that already exists and already goes red when
  a name appears without registration (nail N1).
- ⚠ But it **does add a second caller of an existing door, and it adds one *binding-visible* JS symbol**
  (`window.wispBridge.postMessage`). Two facts the orchestrator must weigh, both measured:
  ① `wispDispatch` is invoked by the library through `window.external.invoke(JSON.stringify({id,method,params}))`
  (`webview.go:472-476`). A shim that instead calls `window.external.invoke` directly **with a forged
  `{id,method:"wispDispatch",params:[raw]}`** would reach `dispatchRaw` while *impersonating the library
  RPC pipe* — that shape is a new inbound path in everything but name, and it is the shape §4's mutation
  ⓑ is designed to catch. A shim that calls `window.wispDispatch(raw)` uses the real binding promise path.
  ⇒ **记一行代价：甲 有两条落地子形（`window.wispDispatch(raw)` vs 手造 `{id,method,params}`），
  后者等价于往白名单外加了一条"以 RPC 形状进来的路"，⛔ 本腿不裁，但它是必须被 §4 的突变照出来的一支。**
  ② ⛔ 我没有、也不会往 `bridge.go` 加名字；三形中任何一支若加名，nail N1 立刻红。

**"这条边以后谁能看见" (who can see this edge later) — measured scanner coverage:**
- `internal/panel/composer_dispatch_test.go:588-660` `hostChannelCapabilityHits` reads **production Go
  sources only**; it explicitly skips directories `.git, .scratch, node_modules, dist, third_party,
  scripts, docs, frontend, design, build` (`:605`) and skips comment lines (`:649`). ⇒ A JS shim that
  lives **inside a Go string literal** is visible to it only through the *dependency* half
  (`hostModuleTokens`, lowercased substring over non-comment lines `:652-657`), and its *identifier*
  half cannot see JS at all. Its own doc comment says the blind spot out loud (`:584-587`): "A host
  reached through a hand-rolled COM vtable that spells none of these names in an identifier is outside
  this ruler's sight".
- `inbound_roster_253_test.go` scans Go AST for declared names, guard `case` labels, router `case`
  labels and dotted string literals in production Go — so a dotted route name invented **inside the shim
  JS** would be caught by the "running guard answers a dotted literal not on any roster" denominator
  (`:507-517`) only if it also appears in Go; a shim-side-only invented name would be **invisible to it**.
  ⇒ 记代价：**甲把一条语义边放进了所有现有仪器都不解析的语言（JS 字符串）里；今天没有任何一把尺能看见它。**
- Page side: `bridge_test.go:162` and `composer_test.go:510` count `postMessage`/host call sites in
  `frontend/src`; they keep counting 2 under 甲 (page untouched) — ⚠ which is exactly why **those two
  nails cannot be used as evidence that the transport works** (they'd read the same green before and
  after the fix). Registered for §4.

### 2.2 乙 — don't use the library's RPC pipe (own `AddWebMessageReceived` / another COM route)

- **Dependency verdict, measured:** ⛔ no `go.mod`/`go.sum` edit and **no module hand-edit is needed to
  get at the COM layer** — `github.com/jchv/go-webview2/pkg/edge` is an exported package in the same
  module already required (`go.mod:19`), and `edge.NewChromium()`, `Chromium.MessageReceived`,
  `ICoreWebView2.AddWebMessageReceived` (`corewebview2.go:108`) are all exported.
  ⇒ but ⚠ **`edge.Chromium` is a whole host implementation, not a hook**: today `cmd/wisp` uses the
  high-level wrapper (`webview2 "github.com/jchv/go-webview2"` at `panel_host_windows.go:64` and
  `panel_resident_windows.go:73`, the only two non-test importers). Driving `pkg/edge` directly means
  replacing `NewWithOptions` (`:386`), `Embed`, `Resize` (`chromium_amd64.go:12`, referenced in this
  file's own header comment `:47-52`), the message pump and the STA wiring.
- **There is no narrow version of 乙 inside the wrapper:** the `WebView` interface method set is exactly
  `Run, Terminate, Dispatch, Destroy, Window, SetTitle, SetSize, Navigate, SetHtml, Init, Eval, Bind`
  (`common.go:30-83`) — it exposes **no accessor for the `*edge.Chromium` or the `ICoreWebView2`**, and
  `webview`'s fields are unexported (`webview.go:99-110`: `w.browser = chromium` is private).
  `chromium.MessageCallback` is a **single slot** set once at `:103`, so 乙 cannot even *co-register* a
  handler on the existing window. ⇒ **fork-or-no-fork verdict: 乙 does NOT need a fork to exist, but it
  DOES need to bypass the wrapper entirely; a "fork to add a second callback / to stop the echo" would
  be needed to keep today's host shape.**
- **Cost if it takes the fork branch (as asked):** a `replace` in `go.mod` + `go.sum` churn +
  ⛔ *手改模块依赖本仓禁* ⇒ the fork itself is a permission-level blocker, not just a line count.
- **Whose nail turns red on 乙 (this is the sharp part):**
  `internal/panel/composer_dispatch_test.go:405-412` states `TestPanelHostIsAttachedAndNamesTheWindowHops`
  requires **BOTH** "a production identifier from the family … **AND** a webview module in the main
  module's `go.mod`. **Either signal missing ⇒ red**". So removing `jchv/go-webview2` from `go.mod`
  turns that nail red by design; and `SPEC-08:145` names that library in the spec, so it is also a
  contract-level change (人工批准), not a code choice.
- `panel_host_gate_test.go:37` forbids imports `net, net/http, net/netip, net/mail, net/tcp` in
  `panel_host_windows.go` — ⛔ 乙 as "another COM route" doesn't need them, but 乙-as-alternative-shape
  "served over localhost" would die there (N5).

### 2.3 丙 — page calls `window.wispDispatch(…)` instead (⛔ page files are NOT in this fleet's remit; recorded as the comparison row)

- **Page surface, measured on dev:** the senders are **1 file, 2 call sites**:
  `frontend/src/lib/panel.ts:179` and `:218` (plus the interface/global declaration `:139-155` which
  would need a shape change if the return type is to become a promise).
- **dev vs branch — stronger than "same lines", and it is now measured by blob hash:**
  `frontend/src/lib/panel.ts` is **byte-identical on both trees**:
  `git hash-object frontend/src/lib/panel.ts` = `git rev-parse HEAD:…` = `git rev-parse dsh/feat/frontend-p0-v2:…`
  = **`5ea999d5444aa13cadbd3d94a895b44c2bfdf697`** (`logs/q2-shape-c-page-cost.txt`).
  ⇒ **一份改动，两棵树同一处，且不需要 tree-specific patch**（sender 行在两树同为 `:179`/`:218`，
  `panel.ts` 同为 311 行）。分支对 `frontend/src` 的总差异是 37 files / 7229 insertions（`git diff --stat`
  尾部见同一 log），⛔ 但不含这枚文件。
  File counts do differ (`git ls-files frontend/src`: dev **64** — the batch ticket 274 builds into the exe —
  vs branch **85**), so the *roster of page files* is tree-dependent while **丙's touched file is not**.
- **Cost 丙 does *not* pay that 甲 pays:** 丙 needs no injected JS, no `window._rpc` guard, no ordering
  constraint against `SetHtml`. **BUT 丙 still does not get receipts**: `window.wispDispatch(raw)` returns
  a *promise* (the library's Bind script returns a Promise, `webview.go:464-478`), so 丙 must either
  ignore it (the page's own contract at `panel.ts:162-166` says "return nothing") or create
  `window._rpc` slots — i.e. **丙 still leaves the reply half unsolved**, same as 甲.
- **⛔ Permission:** `frontend/**` is read-only for this fleet (owner 10-07; orchestrator's §3 of the
  ticket records ⓐ/ⓑ/ⓒ with ⓐ = default). 丙 is unaffordable without a permission change.

### 2.4 Cross-shape cost summary (no selection)

| | files to touch (Go) | page files | roster | dep/`go.mod` | ordering hazard | who can see the edge later |
|---|---|---|---|---|---|---|
| 甲 | `cmd/wisp/panel_host_windows.go` (1 file, ~1 stmt + JS const) | 0 | unchanged (⚠ two sub-shapes, see §2.1②) | none | `Init` must precede first `SetHtml` (`:449`/`:468`/`:681`); must survive later `SetHtml` re-creation; `window._rpc` guard required | ⛔ no existing scanner parses the JS; `hostChannelCapabilityHits` skips comments/strings & `frontend/`; roster nail only sees Go |
| 乙 | replaces the host's creation+pump path (`:386`→`:415`, plus pump/`Embed`/`Resize`) — several files | 0 | unchanged | `replace`/fork + `go.mod`+`go.sum` if the wrapper is kept; ⛔ 手改模块依赖 = 需人工批准 | must re-implement the STA pump (header comment `:29-44` documents the freeze/nested-pump trap) | red on `TestPanelHostIsAttachedAndNamesTheWindowHops` if the library leaves `go.mod`; spec-named library (`SPEC-08:145`) |
| 丙 | 0 | `frontend/src/lib/panel.ts` — 2 sender lines (`:179`,`:218`) + `:139-155` typing; **same lines on dev and on branch** | unchanged | none | none | page-side only; ⛔ write permission not granted |

---

## §3 撞哪几枚既有钉（⛔ "会让门变红" 与 "只是注释过期" 分两格 — question ③）

### 3.1 会让门变红（assertion has teeth; transcript `logs/q3-nails-scan.txt`, `logs/q3-nails-2.txt`)

| # | 钉（文件＋行） | 断言原文（节选，逐字） | 甲 | 乙 | 丙 |
|---|---|---|---|---|---|
| N1 | `internal/panel/inbound_roster_253_test.go:419` `TestFullInboundMethodRosterIsClosed` | `:437` "bridge.go declares the inbound method name %s … the closed roster does not carry it"; `:465` "the roster carries %s but the running guard `knownComposerMethod` refuses it"; `:500/:503` router↔roster both directions; `:515` "the running guard answers the dotted literal %q … on no roster" | ⛔ 不红 **only if no new dotted name is added to Go** — 红 the moment a shim-side or Go-side new route name is registered | 不红 | 不红（除非页面开始发未注册名 ⇒ 由 N2/N3 抓） |
| N2 | `internal/panel/bridge_test.go:162` | `"postMessage call sites = %d, want 2 (the approval request and sendRequest)"` | **不红**（页面未动）——⚠ 正因如此它**不能当甲的凭据** | 不红 | **红**（改成 `wispDispatch` 后 `postMessage` 计数＝0） |
| N3 | `internal/panel/composer_test.go:510` 与 `:518-522` | `:510` "host call sites in frontend/src = %d, want 2 (the approval request and sendRequest)"; `:518` "sendRequest is called with a computed route…"; `:522` **"the renderer names a route the Go side does not answer"** | 不红 | 不红 | **红**（同一把词面尺） |
| N4 | `internal/panel/composer_dispatch_test.go:405-412` `TestPanelHostIsAttachedAndNamesTheWindowHops` (impl `:436` `hostChannelCapabilityHits`, controls `:455-465`, `:474-476`, `:490-499`, `:518-531`) | `:409-412` "a native host / WebView2 message channel IS attached … **AND** a webview module in the main module's go.mod. **Either signal missing => red**" | 不红 | **红**（若库从 `go.mod` 消失／宿主标识被中和；`:518-531` 那枚反向对照正是把标识中和成 0 命中来证明此尺有牙） | 不红 |
| N5 | `cmd/wisp/panel_host_gate_test.go:30-45` `TestPanelHostOpensNoListeningSocketL1` | `:42` "panel host imports %q, which can open a listening socket; D29/AC#3 forbid a localhost server"（forbidden = `net, net/http, net/netip, net/mail, net/tcp`，`:37`） | 不红（`Init`/`Eval` 不引新 import） | ⚠ 仅当乙改走 localhost 才红；纯 COM 路不红 | 不红 |
| N6 | `cmd/wisp/panel_resident_windows_test.go:166-185` (`evalOnPanelThread`, 内含 `:175 w.Eval(js)`) 与调用点 `:301`, `:801`, `:805`, `:808`, `:864` | `:301` `parts = append(parts, "window.wispDispatch("+reportJSEnv("ac13-probe", "b")+");")`; `:805` `"Promise.resolve(window.wispDispatch(ENVS[i])).then(function(v){ return 'REPLIED'; }),"`; `:199` 注释 "So the reports ride window.wispDispatch, the product's own…" | ⚠ **这些用例今天就是"从 Eval 里直接叫 `window.wispDispatch`"**，即**测试自己扮演了丙**；甲落地后它们**照旧绿**，⛔ 因此**它们对"页面的请求能不能到 Go"零分辨力**（见 §4 的恒真指控） | 红：`w.Eval` 走的库通道若被换掉，这几枚真机用例的形状就不再对应产品边 | ⚠ 甲/丙下不红；若有人**顺手把这些 Eval 注入改成依赖垫片**，则这几枚会从"产品边的证明"退化成"垫片的证明" ⇒ 必须显式记账 |
| N7 | `cmd/wisp/panel_host_gate_test.go:177-185` `TestAC1SessionDisposeHasAProductionTrigger_AC1` | `:185` "…no non-test file calls Destroy ON THE MANAGER. The %d .Destroy() call(s)…" | 不红 | ⚠ 乙重写宿主创建路径易把 `Destroy` 触发点搬走 ⇒ 需现量确认 | 不红 |

### 3.2 只是注释过期（⛔ 不是门，写腿不许为了让它"对上"去动代码；也不许反过来删注释换绿）

| # | 位置 | 过期内容（逐字摘录） | 与什么不符 |
|---|---|---|---|
| M1 | `cmd/wisp/panel_host_windows.go:402-404` | "the raw envelope is parsed by panel.ParseComposerRequest inside Handle, whose whitelist is fixed at the **four** existing methods (bridge.go:42-45, ticket 248 owns config methods)" | `bridge.go:42-47` 现量 **6 枚**（编排者票面 §5 已登记，本腿复认） |
| M2 | `frontend/src/lib/panel.ts:146` | "Installed by WebView2's **AddHostObjectToScript** / postMessage pipe." | `AddHostObjectToScript` 全仓（含 lib）＝ 0 命中；本腿在 `logs/q1-transport-spelling-scan.txt` 复跑确认只在 evidence/ledger 出现 |
| M3 | `internal/panel/composer_dispatch_test.go:459-461` 桩 | `func (w *CoreWebView2) PostWebMessageAsJson(json string) error { return nil }` | 库的真名是 `PostWebMessageAsJSON`（`corewebview2.go:106`，大写 JSON），且该桩挂在 `writeHostCarrier` 现造的 `type CoreWebView2 struct{}` 上，**不是库类型的实现** |
| M4 | `docs/PLAN.md:1747` | "docs/contracts/C1..C31.md …" | 交付物不存在（AGENTS.md §4 已登；非本票射程） |
| M5 | `go.mod:19` | `github.com/jchv/go-webview2 … // indirect` | `cmd/wisp/panel_host_windows.go:64`、`panel_resident_windows.go:73` 直接 import ⇒ 账面过期（`go mod tidy` 会挪）。⛔ 本腿不动 `go.mod` |

---

## §4 恒真性进攻计划（question ④ — 本票最硬的规矩）

**前提读数：** 甲 之后，"页面请求到达 `dispatchRaw`" 这件事**今天没有任何一枚具名用例会红**——
N2/N3 数页面词面（甲 不动页面 ⇒ 恒绿），N6 从 `Eval` 里自己叫 `window.wispDispatch`（绕过页面 ⇒ 恒绿），
N1 数 Go 名册（与传输无关 ⇒ 恒绿）。⇒ **判据必须新增，不许借旧钉报功。** 下面是两发定向突变与它们
各自必须红在哪一枚具名用例。

**⚠ 新增判据的前置形状（我只报代价与射程，不写实现）：** 一枚 **Go 侧**用例，它 ① 取到**页面实际会发的
那串原始信封**（`{"method":"panel.approval.request","correlationId":…,"outcome":…}`，逐字取自
`frontend/src/lib/panel.ts:180-184` 现量），② 走**垫片之后的入口**而不是 `window.wispDispatch`，
③ 断言 `dispatchRaw`（`panel_host_windows.go:630`）或其下游**被调用过且带那枚 correlationId**，
④ 断言**未注册名仍被拒**（不许把 ③ 写成"来者都 ack"）。

| 突变 | 改什么 | 必须红在哪 | 为什么这一枚抓得住它 |
|---|---|---|---|
| ⓐ **垫片消失／不再转发**：删掉那条 `w.Init(<shim>)` 语句，或把垫片里的转发调用注释掉 | 新判据（上面那枚） | ⛔ 若新判据仍绿 ⇒ **该判据读的不是这条边**，判为不合格并退回（这正是票面 §4 那句"必须能指名它断言了哪个 Go→页调用点"的入向版） | 信封重新只能到 `msgcb`（`webview.go:139`）而 `msgcb` 不认识它 ⇒ `dispatchRaw` 零调用 ⇒ "被调用过且带 correlationId" 断言必然拿不到记录 |
| ⓑ **"什么都不解析直接 ack"**：把 `panel_host_windows.go:405-407` 的 Bind 闭包（或 `:630 dispatchRaw` 的第一跳）改成 `return "ok"` 而不进 `ParseComposerRequest`／`knownComposerMethod` | **既有钉即已覆盖**：`internal/panel/inbound_roster_253_test.go:465`（"the roster carries %s but the running guard `knownComposerMethod` refuses it"）与 `internal/panel/composer_dispatch_test.go:201` `TestUnlistedMethodNameIsRefusedAndAudited`、`:231` `TestRosterMismatchBackstopRefusesInsteadOfAccepting`、`:126` `TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce` | ⚠ **但这三枚只守 Go 内部**，⛔ 不能抵"页面真的进来了"：若判据仍从 `Eval`/测试直调 `dispatchRaw`（N6 那个形状），ⓑ 会**在被绕过的地方红、在真正的边上是绿的** ⇒ 所以 ⓑ 的验收必须同时指认"新判据也红"。这一条就是今天那枚假绿的形状 | 白名单/路由的三条既有钉会红；新判据若也红才说明它真的经过这条边 |
| ⓒ 顺手第三发（本腿自己加的攻击，成本 0）：把垫片改成"用 `window.external.invoke` 手造 `{id:0,method:'wispDispatch',params:[raw]}`"（§2.1 的子形②） | ⛔ 今天**没有一把尺会红**：`external.invoke` 全仓 0 命中（`logs/q1-transport-spelling-scan.txt` 同尺的 Go/前端 scope），`inbound_roster_253` 不看 JS，`hostChannelCapabilityHits` 跳字符串 | ⇒ 记为**已知盲区**并具名交给编排者：要么新判据必须断言"进入的是 binding 的返回路径（有 Promise 回执）"，要么显式接受这条子形并在票上写清它没尺 | — |

**⛔ 明写（按票面要求）：** 若上面任一判据**换成反形还全绿**，那把尺对这件事不敏感，**不许当凭据**，
也不许用它去勾票 35 的任何一格。特别地：**N6（`panel_resident_windows_test.go:301/801/805/808/864`）
与 N2/N3 今天都不能用来证明这条边**——前者从 `Eval` 直调 binding，后者只数页面词面。

---

## §5 顺手补的两格（question ③ in the brief's §3）

### 5.1 那句 "responses arrive as a fresh PanelSnapshot push"（`frontend/src/lib/panel.ts:162-166`，引 `PLAN.md:1044`）要兑现，Go 侧最少要哪几枚调用点

现量：**"push" 的今天生产者枚数＝0；听众枚数＝0。** 逐项给枚数，⛔ 不写"接上就行"：

- `Dispatch(` 非 test 全仓 **7 枚**（`logs/q2-init-semantics.txt`），但**其中只有 1 枚是 WebView2 的**：
  `cmd/wisp/panel_resident_windows.go:360 w.Dispatch(fn)`（在 `post()` 内，`:349-372`，两态路由
  `wRef` → `Dispatch` / `tasks` 队列）。其余 6 枚是**状态机**的 `Dispatch`，与面板线程无关：
  `cmd/balldebug/main.go:618`、`internal/models/bridge.go:40,46,59,64`、`internal/statemachine/machine.go:202`。
  ⇒ ⚠ **这是对票面 §2 "泵到面板线程的路真在（`Dispatch` 7 枚调用点）" 的一句具名修正：
  面板线程的 `Dispatch` 调用点今天＝ 1 枚，不是 7 枚。** 编排者裁形时请按 1 枚算。
- `Eval(` 非 test ＝ **0 枚**（本腿复跑同尺，命中只有 `panel_resident_windows_test.go:175` 一枚 test 侧）。
  ⇒ **推送的那一句 `Eval` 由谁调，今天无人**：可用的最小形状是 `rp.post(func(){ w.Eval(<snapshot JS>) })`
  ——即"经唯一的 `:360` 那条 `post()` 落到面板线程，再由该线程调 `Eval`"。理由写在这枚文件的头部注释
  `cmd/wisp/panel_host_windows.go:36-42`：`Dispatch` 只 append 到库私有队列，**唯一读者是 `Run()`**，
  手工泵的主机永远不排空 ⇒ 不经 `post()` 的 `Eval` 不会被送达。
- **最少调用点清单（3 枚，缺一不可）**：
  ① `cmd/wisp/panel_pump.go:321` 一带——今天 `Out: rt.bookPanelSnapshot`（`cmd/wisp/run.go:725`）之后
  只写 `rt.lastSnap/lastSnapBytes`＋一行日志，**必须在此新增一次对外调用**（读者
  `(rt *agentRuntime).lastPanelSnapshot()` 非 test＝0 处，票面 §1⑦ 已复认）。
  ② `cmd/wisp/panel_resident_windows.go:349-372 post()`——唯一能把闭包送上面板线程的路（`:360`）。
  ③ 一个 `w.Eval(...)` 调用点（**目前不存在，需新建**；`w` 只能从 `rp.wRef`／
  `mgr.currentWindow()`（`panel_host_windows.go:302`）取）。
  ⇒ ⛔ 三枚之外还要第 4 件事才可能"被看见"：**页面得有耳朵**（dev `frontend/src` 里
  `addEventListener` 4 枚命中、无一为 `message`；`window._rpc`/`onmessage`/`WebMessageReceived` = 0），
  而页面写权限不在本轮 ⇒ **量到的结论是：出向今天最少差 3 枚 Go 调用点 + 1 枚页面接收器；
  只补 Go 那 3 枚，页面仍然画不出来。** 这一格我不裁形，只报枚数。

### 5.2 `internal/panel/composer_dispatch_test.go:461` 那枚空桩——今天被谁依赖（⛔ 只登记，未动）

- 它出现在 `writeHostCarrier(...)` 注入的**fixture 源码字符串**里（`:455-462`），载体是临时造的
  `type CoreWebView2 struct{}`（`:457`）＋ `WebMessageReceived`（`:459`）＋ `PostWebMessageAsJson`（`:461`）。
- 依赖它的用例：`TestPanelHostIsAttachedAndNamesTheWindowHops`（`:434`）的**第一枚正向对照**
  （`:436` 与 `:455-471`：`hostChannelCapabilityHits(t, a)` 必须 `len(got) != 0`，否则
  `:464` 报 "POSITIVE CONTROL RED (the ruler, not the product)"）。同函数还有另外三枚载体
  `:474`（假 `go.mod`/`go.sum`）、`:490`（CLI seam 反向）、`:518`（标识中和反向）。
- ⇒ **正确读数（⛔ 别写错）：** 这枚桩**不在生产路径上、也不是产品对端**；它只是让"宿主族标识存在"
  这条尺的正向对照能构造出来。**它的"吞成 nil"效应只影响 fixture 树，不影响任何真实出向送达。**
  真正的危险不是它，而是 §4/§3 那两枚：`test` 里 `PostWebMessageAsJson` 拼写与库的
  `PostWebMessageAsJSON` 不同（M3），**任何人若拿"这包绿"抵"页面收到过东西"，那把尺连这个符号都不指向库的类型。**
- ⛔ 未改动，只登记。

---

## §6 门禁与还原自证

**Static gate only, as ordered.** Full transcript: `logs/gate-d22scan.txt`.
```
$ sh scripts/d22scan.sh
rc=0
```
Verdict line read verbatim from the log: `d22scan: clean - no D22 ban violations`, and it
examined — bans #1-5 `internal/` **228** production Go files, `cmd/` **38**; ban #6 `frontend/`
**85** text files; ban #7 `internal/tools/` **23**; ban #8 `design/` **39**, `frontend/` **85**,
`internal/` **514**, `cmd/` **104** Go files (comments and `_test.go` included).
It also reports `skipped as git-ignored: 1 file(s) … frontend/dist/assets/` (by `frontend/.gitignore`).

⚠ **One thing I checked because it could have stepped on `274-r1`:** step 1 of `scripts/d22scan.sh:50-51`
runs `tools/d22scan/runtests.sh -C tools/d22scan ./...` — the `-C` pins the package dir, so the only
`ok` line in my log is `github.com/CarlosShao/wisp/tools/d22scan` (1 package, 19.5s).
**It is NOT a repo-wide `go test`; I ran no other test binary and asked for no single-case run.**
The `PASS=35` in that line is this instrument's own self-test count, not 35 repository packages.

**⛔ §3 表的性质（必须在裁决时写在前面）：** 那张"红/不红"表是**读断言原文得到的静态判定**，
本腿没有跑过 `internal/panel` 或 `cmd/wisp` 任何一枚用例。

**还原自证 / 写入面：**
```
$ git status --porcelain -- cmd internal scripts tools .github docs frontend
 M docs/reports/pending-and-issues.md
 M scripts/build.ps1
```
— identical to the opening anchor (same two entries, both **not mine**: ledger = orchestrator,
`scripts/build.ps1` = `274-r1`). ⛔ 本腿未碰产码、未碰 `docs/**`、未碰 `go.mod`/`go.sum`、未碰 `frontend/**`
（只读）。写入面只有 `.scratch/wisp/probes/35/a2/transport-cost.md` 与 `.scratch/wisp/probes/35/a2/logs/`。
⛔ 未 push。

---

## §7 未做完的格（逐枚具名，⛔ 不写成结论）

1. **甲 的实机送达验证**：`Init` 注入的垫片在真机 WebView2 上是否真的在第一次 `SetHtml` 前执行、
   且对 `serveEntry` 重建的新文档仍生效——⛔ 本腿零编译零真机，未量。**取不到。**
2. **丙 的完整页面代价面**：只量到 `panel.ts` 的 2 枚 sender + `:139-155` 类型块；
   ⛔ 未量页面里是否还有依赖 `WispHostBridge` 返回形状的用例（`frontend/**` 只读，未做全量 TS 依赖分析）。
3. **`callbinding` 对 `wispDispatch` 的参数解包细节**：`webview.go:162-180` 我**没有整段读**（只读到
   `:140-158` 与 `:450-482`）⇒ "甲 的 `window.wispDispatch(raw)` 单参数是否被原样交给 Go 闭包"
   这一格**未现量，取不到**，引用前需补读。
4. **真机再入测量／`addEventListener('message')` 在 WebView2 里可行性**：⛔ 本腿未做（`35-a1` 也声明未做）。
5. **单枚用例的实跑**：⛔ 未跑任何 `go test`（本轮整包门禁归 `274-r1`）；§3 表里所有"红/不红"
   都是**读断言原文得出的静态判定**，非运行读数——这一条要在裁决时写在前面。
6. **`docs/PLAN.md:1044`（页面注释引的那句）**：本腿只按票面引用它，**未整行读原文** ⇒ 它与
   "push" 承诺的对应关系**未复证**，属 `35-a1`/票 36-40 射程。
## §0 起手锚 (landing gate — first commit, before any long-running command)

```
$ date
Wed Oct  7 11:55:16 CST 2026

$ git rev-parse --short HEAD
28872c9a

$ git branch --show-current
dev

$ git status --porcelain -- cmd internal scripts tools .github docs frontend
 M docs/reports/pending-and-issues.md
 M scripts/build.ps1
```

Pre-existing dirty entries NOT touched by this leg (registered as-is, per §Git 纪律):
`design/**`, `.gitignore`, several `.scratch/**`, `docs/reports/pending-and-issues.md`,
`scripts/build.ps1` (another leg `274-r1` is writing it).

### 8 anchors, re-read live at HEAD `28872c9a` (verbatim)

**Anchor 1 — `go.mod:15-21`** (the dependency marking):
```
15:
16: require (
17:	github.com/dustin/go-humanize v1.0.1 // indirect
18:	github.com/google/uuid v1.6.0 // indirect
19:	github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect
20:	github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1 // indirect
21:	github.com/k2-fsa/sherpa-onnx-go-linux v1.13.8 // indirect
```

**Anchor 2 — library root `webview.go:99-110`** (`MessageCallback` wiring):
```
99: 	w.bindings = map[string]interface{}{}
100: 	w.autofocus = options.AutoFocus
101:
102: 	chromium := edge.NewChromium()
103: 	chromium.MessageCallback = w.msgcb
104: 	chromium.DataPath = options.DataPath
105: 	chromium.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)
106:
107: 	w.browser = chromium
108: 	w.mainthread, _, _ = w32.Kernel32GetCurrentThreadID.Call()
109: 	if !w.CreateWithOptions(options.WindowOptions) {
110: 		return nil
```

**Anchor 3 — library `webview.go:139-160`** (`msgcb`, the RPC parser):
```
139: func (w *webview) msgcb(msg string) {
140: 	d := rpcMessage{}
141: 	if err := json.Unmarshal([]byte(msg), &d); err != nil {
142: 		log.Printf("invalid RPC message: %v", err)
143: 		return
144: 	}
145:
146: 	id := strconv.Itoa(d.ID)
147: 	if res, err := w.callbinding(d); err != nil {
148: 		w.Dispatch(func() {
149: 			w.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
150: 		})
151: 	} else if b, err := json.Marshal(res); err != nil {
152: 		w.Dispatch(func() {
153: 			w.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
154: 		})
155: 	} else {
156: 		w.Dispatch(func() {
157: 			w.Eval("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined")
158: 		})
159: 	}
160: }
```

**Anchor 4 — library `webview.go:450-482`** (`Bind`, the injected script):
```
450: func (w *webview) Bind(name string, f interface{}) error {
451: 	v := reflect.ValueOf(f)
452: 	if v.Kind() != reflect.Func {
453: 		return errors.New("only functions can be bound")
454: 	}
455: 	if n := v.Type().NumOut(); n > 2 {
456: 		return errors.New("function may only return a value or a value+error")
457: 	}
458: 	w.m.Lock()
459: 	w.bindings[name] = f
460: 	w.m.Unlock()
461:
462: 	w.Init("(function() { var name = " + jsString(name) + ";" + `
463: 		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
464: 		window[name] = function() {
465: 		  var seq = RPC.nextSeq++;
466: 		  var promise = new Promise(function(resolve, reject) {
467: 			RPC[seq] = {
468: 			  resolve: resolve,
469: 			  reject: reject,
470: 			};
471: 		  });
472: 		  window.external.invoke(JSON.stringify({
473: 			id: seq,
474: 			method: name,
475: 			params: Array.prototype.slice.call(arguments),
476: 		  }));
477: 		  return promise;
478: 		}
479: 	})()`)
480:
481: 	return nil
482: }
```

**Anchor 5 — library `webview.go:435-448`** (the four outbound primitives):
```
435: func (w *webview) Init(js string) {
436: 	w.browser.Init(js)
437: }
438:
439: func (w *webview) Eval(js string) {
440: 	w.browser.Eval(js)
441: }
442:
443: func (w *webview) Dispatch(f func()) {
444: 	w.m.Lock()
445: 	w.dispatchq = append(w.dispatchq, f)
446: 	w.m.Unlock()
447: 	_, _, _ = w32.User32PostThreadMessageW.Call(w.mainthread, w32.WMApp, 0, 0)
448: }
```

**Anchor 6 — `pkg/edge/chromium.go:233-248`** (`MessageReceived`, echo back to page):
```
233: func (e *Chromium) MessageReceived(sender *ICoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr {
234: 	var message *uint16
235: 	_, _, _ = args.vtbl.TryGetWebMessageAsString.Call(
236: 		uintptr(unsafe.Pointer(args)),
237: 		uintptr(unsafe.Pointer(&message)),
238: 	)
239: 	if e.MessageCallback != nil {
240: 		e.MessageCallback(w32.Utf16PtrToString(message))
241: 	}
242: 	_, _, _ = sender.vtbl.PostWebMessageAsString.Call(
243: 		uintptr(unsafe.Pointer(sender)),
244: 		uintptr(unsafe.Pointer(message)),
245: 	)
246: 	windows.CoTaskMemFree(unsafe.Pointer(message))
247: 	return 0
248: }
```

**Anchor 7 — `pkg/edge/corewebview2.go:104-112`** (COM vtable names; note `AsJSON` casing):
```
104: 	CapturePreview                         ComProc
105: 	Reload                                 ComProc
106: 	PostWebMessageAsJSON                   ComProc
107: 	PostWebMessageAsString                 ComProc
108: 	AddWebMessageReceived                  ComProc
109: 	RemoveWebMessageReceived               ComProc
110: 	CallDevToolsProtocolMethod             ComProc
111: 	GetBrowserProcessID                    ComProc
112: 	GetCanGoBack                           ComProc
```

**Anchor 8 — `cmd/wisp/panel_host_windows.go`**, four live reads:

`:78-82` (binding name):
```
78:
79: const (
80: 	panelDispatchBinding = "wispDispatch"
81: 	// panelWidthPx / panelHeightPx are what the host asks WebView2 for when NOBODY
82: 	// handed it a geometry source (票 255 AC#4: these two are now "the value when
```
`:400-412` (the one inbound door, and the comment the orchestrator already flagged as stale):
```
400: 	// The one inbound door. JS -> Go via a bound function over the WebView2
401: 	// message channel; the returned string is the receipt that reaches the page.
402: 	// No new inbound method name is invented here: the raw envelope is parsed by
403: 	// panel.ParseComposerRequest inside Handle, whose whitelist is fixed at the
404: 	// four existing methods (bridge.go:42-45, ticket 248 owns config methods).
405: 	bindErr := w.Bind(panelDispatchBinding, func(raw string) string {
406: 		reply, _ := m.dispatchRaw(ctx, raw)
407: 		return reply
408: 	})
409: 	if bindErr != nil {
410: 		w.Destroy()
411: 		return fmt.Errorf("panel host: bind dispatch door: %w", bindErr)
412: 	}
```
`:440-470` (`SetHtml` call sites — the Init-timing question):
```
440: // the user is never left looking at the round-trip probe page, and it still opens
441: // no socket and loads nothing over the network.
442: func (m *PanelManager) serveNotBuiltNoticeLocked() {
443: 	m.mu.Lock()
444: 	m.mu.Unlock()  [tail: w := m.w / if w == nil { return }]
445: 	w.SetHtml(`<!doctype html>... panel assets unavailable ...`)
...   serveEntry(): w.SetHtml(string(data))
```
(full literal text of `:440-470` is in `logs/anchor-reads.txt`; the two `SetHtml` statements are at
`panel_host_windows.go:449` and `:468`, third one recorded in §2.)
`:625-635` (`dispatchRaw`):
```
625: 		}
626: 	}
627:
628: // dispatchRaw runs one page envelope through the router. Kept separate from the
629: // bind closure so a test can call the same code path the page reaches.
630: func (m *PanelManager) dispatchRaw(ctx context.Context, raw string) (string, error) {
631: 	if m.disp == nil {
632: 		return "panel host: no inbound router attached", fmt.Errorf("panel host: nil router")
633: 	}
634: 	if ctx == nil {
635: 		ctx = context.Background()
```
`:665-672` (`wispProbeRT` second binding):
```
665: 	// A second, one-shot binding the probe page calls; its Go side just signals
666: 	// the round trip completed. Using a binding keeps the probe on the same
667: 	// WebView2 channel the real door uses.
668: 	if err := w.Bind("wispProbeRT", func() string {
669: 		select {
670: 		case <-done:
671: 		default:
672: 			close(done)
```

**Anchor 9 — `internal/panel/bridge.go:38-52`** (the roster, 6 entries):
```
38: // test that does not exist in this repository; ticket 35's snapshot pump was the
39: // step that checked it, and the behaviour was already covered - only the name
40: // was wrong.)
41: const (
42: 	MethodModeRequest      = "panel.mode.request"
43: 	MethodWorkspaceRequest = "panel.workspace.request"
44: 	MethodAttachmentAdd    = "panel.attachment.add"
45: 	MethodMessageSend      = "panel.message.send"
46: 	// MethodConfigGet and MethodConfigSet are ticket 248 AC#1's two settings
47: 	// routes. Three facts about the naming, because each one is a decision
48: 	// somebody could later "fix" by accident:
49: 	//
50: 	//   - The names are NOT invented here. They are the two spellings the
51: 	//     product spec already carries (docs/specs/SPEC-08-ui-ball-panel.md:168,
52: 	//     read-only for this leg), so this file follows a document rather than
```

**Anchor 10 — frontend on the `dev` tree** (line counts re-measured:
`panel.ts` 311, `main.tsx` 61, `App.tsx` 117 — the brief's page line numbers below are dev-true):
```
frontend/src/lib/panel.ts:139-155
139: /** Shape of the object WebView2 installs on the host page (C17). */
140: interface WispHostBridge {
141:   postMessage(message: string): void;
142: }
143:
144: declare global {
145:   interface Window {
146:     /** Installed by WebView2's AddHostObjectToScript / postMessage pipe. */
147:     wispBridge?: WispHostBridge;
148:     chrome?: { webview?: WispHostBridge };
149:   }
150: }
151:
152: function hostBridge(): WispHostBridge | null {
153:   if (typeof window === "undefined") return null;
154:   return window.wispBridge ?? window.chrome?.webview ?? null;
155: }
```
```
frontend/src/lib/panel.ts:162-166 (the "push" promise)
162: /**
163:  * Send one bridge request and return nothing: responses arrive as a fresh
164:  * PanelSnapshot push, never as a return value, because a panel that keeps a
165:  * copy of the answer would be a second state holder (PLAN.md:1044).
166:  */
```
```
frontend/src/lib/panel.ts:179-186 (sender 1)   /   :218-225 (sender 2)
179:   bridge.postMessage(
180:     JSON.stringify({
181:       method: "panel.approval.request",
182:       correlationId,
183:       outcome,
184:     }),
185:   );
186: }
...
218:   bridge.postMessage(
219:     JSON.stringify({
220:       method,
221:       requestId: nextRequestId(),
222:       source: "panel-composer",
223:       ...payload,
224:     }),
225:   );
```
```
frontend/src/main.tsx:56-60
56: createRoot(document.getElementById("root")!).render(
57:   <StrictMode>
58:     {harness === "1" ? <AppHarness /> : harness === "2" ? <Showcase /> : <App />}
59:   </StrictMode>,
60: );
```
```
frontend/src/App.tsx:49-54 (EMPTY)  /  :69-73 (default prop)
49: const EMPTY: PanelSnapshot = {
50:   pending: [],
51:   results: [],
52:   composer: EMPTY_COMPOSER,
53:   generatedAt: "",
54: };
...
69: export default function App({
70:   snapshot = EMPTY,
71: }: {
72:   snapshot?: PanelSnapshot;
73: }) {
```

**Anchor 11 — `internal/panel/composer_dispatch_test.go:455-465`** (the empty stub):
```
455: 		writeHostCarrier(t, a, filepath.Join("internal", "panel"), "host_windows.go", `package panel
456:
457: type CoreWebView2 struct{}
458:
459: func (w *CoreWebView2) WebMessageReceived(sender, args any) error { return nil }
460:
461: func (w *CoreWebView2) PostWebMessageAsJson(json string) error { return nil }
462: `)
463: 		if got := hostChannelCapabilityHits(t, a); len(got) == 0 {
464: 			t.Errorf("POSITIVE CONTROL RED (the ruler, not the product): a tree whose production source hosts a "+
465: 				"CoreWebView2 and takes WebMessageReceived answered 0 capability hits, so the check above is a "+
```

### Corrections to the brief's page anchors (dev-true, found while re-reading)
The brief and `35-a1` write `hostBridge()` at `panel.ts:146-153` and the push promise at `:159-163`.
Live dev reads: `hostBridge()` is `:152-155` (`:146` is the stale comment line inside `declare global`),
and the "return nothing / PanelSnapshot push" sentence is `:162-166`. `:141`/`:179`/`:218` in the AC
frame are correct as written.

### Library-vs-stub spelling finding (new this leg, affects §2 shape 乙 and §5 stub cell)
The real COM vtable field is `PostWebMessageAsJSON` (capital `JSON`, `corewebview2.go:106`).
The stub at `composer_dispatch_test.go:461` is spelled `PostWebMessageAsJson` and is declared on a
locally synthesised `type CoreWebView2 struct{}` inside a `writeHostCarrier` fixture string —
⛔ it is NOT an implementation of the library's type. Recorded, not touched.
