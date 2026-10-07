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
