# 35-r3 交件｜票 `:63` "Inbound judge must ride the page's own envelope" 的 (d) 三支守卫红

日期 `2026-10-07 21:2x +08`。本腿只写测试：新增 `cmd/wisp/panel_inbound_guards_35r3_test.go`（`//go:build windows`），⛔ 未改产码、⛔ 未碰任何 AC 框、⛔ 未改 `panel_transport_35r2_test.go` 既有断言。锚 commit `9b2551f2`。

---

## ① 三支各落在哪一行 + 三枚用例名 + 绿那发逐字行

守卫现量行号（`internal/panel/bridge.go`，CR=0，与 `git show HEAD:` 逐字节一致；见 `../r3/logs` 前置 start.md）：

| 支 | 守卫 | 判/拒行 | 红句关键词（逐字，本程三支各用一句、互不共用） | 用例名 |
|---|---|---|---|---|
| (d)-1 未注册名 | 名册 | 判 `:132` ⇒ 拒 `:133` | `不是面板 composer 通路的能力入口` | `TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge` |
| (d)-2 名册内来源不符 | 来源 | 判 `:135` ⇒ 拒 `:136-137` | `按伪造/串台拒绝` | `TestInboundSourceGuardRefusesForeignSourceOnPageEdge` |
| (d)-3 名册内缺 requestId | requestId | 判 `:139` ⇒ 拒 `:140-141` | `缺少 requestId` | `TestInboundRequestIDGuardRefusesMissingIDOnPageEdge` |

（解析判 `:128-129` `不是可解析的封套` 不在本程三支之列，未测未改。）

入口＝页面那条边（甲③ 传输 + JS 行为尺）：复用 r2 同包 `main` 的 `fakeDoc35r2`——`installPanelTransport` 先 `Bind(wispDispatch, 门)` 后 `Init(panelPostMessageForwardInit)`，`openDocument()` 在内置 JS 解释器里逐字跑 chromium.go:112 垫片 + webview.go:462-478 Bind 桩 + **产码那枚钩子字符串本身**，`pagePost(envelope)` 即页面 `window.chrome.webview.postMessage(<信封>)`（`frontend/src/lib/panel.ts:179/:218` 真调的那句）。⛔ 不直调 `dispatchRaw`、⛔ 不 `w.Eval("window.wispDispatch(…)")`（N6）。每支断言：`thrown==""` ＋ `doorRounds==1`（信封确由转发钩子送进绑定的门一次，删转发即 `doorRounds=0` ⇒ 三支到达性一起红＝`:63` 的可伪证条款）＋ `lastReply` 逐字含**本支关键词**、且**不含另两支关键词**＋ `len(spy.reqs)==0`（守卫在路由前拒，下游 handler 从未跑＝"是守卫拒的不是 handler 拒的"）。

绿那发逐字行（`logs/mut-baseline-green.txt`，未突变 exe，`rc=0`）：
```
--- PASS: TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge (0.00s)
--- PASS: TestInboundSourceGuardRefusesForeignSourceOnPageEdge (0.00s)
--- PASS: TestInboundRequestIDGuardRefusesMissingIDOnPageEdge (0.00s)
```
三支各自的 `door reply`（同一发 `t.Logf`，逐字）：
- (d)-1 `面板请求被拒绝 [panel.approval.request (无 requestId)]：panel: composer request refused: 方法 "panel.approval.request" 不是面板 composer 通路的能力入口`
- (d)-2 `面板请求被拒绝 [panel.mode.request pc-35r3-2-9e7a-4c1b]：panel: composer request refused: 来源 "panel-composer-x" 不是 "panel-composer"，按伪造/串台拒绝（requestId="pc-35r3-2-9e7a-4c1b"）`
- (d)-3 `面板请求被拒绝 [panel.mode.request (无 requestId)]：panel: composer request refused: 缺少 requestId，无法与审计/卡片对齐（method="panel.mode.request"）`

---

## ② 三支"各点各的"证明（中和任一支 ⇒ 哪枚红、哪两枚仍绿）

手段＝`go test -c -overlay` 在 `/d/tmp/wisp35r3` 造 `bridge.go` 副本、**每份恰一行改动**（把那一支 `if` 的守卫条件前置 `false &&` ⇒ 中和它那枚 `return`，⛔ 仓里一字未改）。三枚 overlay JSON＋三份副本 diff 逐字见 `logs/overlay-diffs.txt`。各覆盖一个 exe，`-test.run 'TestInbound'` 分别跑：

| 中和的守卫（overlay 改的行） | (d)-1 名册 | (d)-2 来源 | (d)-3 requestId | exe rc |
|---|---|---|---|---|
| 基线（无 overlay） | **绿** | **绿** | **绿** | 0 |
| `bridge.go:132` 名册判 → `if false && !knownComposerMethod(r.Method)` | **红** | 绿 | 绿 | 1 |
| `bridge.go:135` 来源判 → `if false && TrimSpace(r.Source)!=ComposerRequestSource` | 绿 | **红** | 绿 | 1 |
| `bridge.go:139` requestId 判 → `if false && TrimSpace(r.RequestID)==""` | 绿 | 绿 | **红** | 1 |

三支今天**分得开**：中和任一支，只有它那枚用例红，另两枚照绿。逐字的红句见 `logs/mut-roster.txt`／`mut-source.txt`／`mut-reqid.txt`：
- 名册被中和时 (d)-1 红句＝`reply does not name guard "不是面板 composer 通路的能力入口" ... 且 leaked a DIFFERENT guard "按伪造/串台拒绝"`——因为 `panel.approval.request` 无 `source`/`requestId`，跳过名册判后**下一道来源判接住它**（这正是票 `:69`/`:71` 说的次第）；它落到了**别支**的文案，而本支断言只认名册那句 ⇒ 红。这恰好证明 (d)-1 断言的是"名册那一道"而非"任何一道拒了就算"。
- 来源被中和时 (d)-2 红句＝`reply="" ... a refusal should never reach the router, but the mode handler ran 1 time(s)`（信封一路过名册、过被中和的来源判、requestId 又合法 ⇒ 抵达路由，路由回了个成功 `""`）。
- requestId 被中和时 (d)-3 红句＝同形（抵达路由）。两处中和下 (d)-1/(d)-3(或 d-2) 都不跟着红，因为它们的失败点在被中和判**之前**（名册）或**之后**（requestId），互不重叠。

---

## ③ 每支断言用的红句关键词（逐字，三支不重句）

1. `不是面板 composer 通路的能力入口`（`:133`）
2. `按伪造/串台拒绝`（`:136-137`）
3. `缺少 requestId`（`:140-141`）

三串互不为子串、无一相同 ⇒ 满足票 `:72`「⛔ 三支不许共用一句文案」与本腿 `assertGuardRedness35r3` 的"本支在、另两支不在"双断言。

---

## ④ (d)-1 与既有钉的差别（本格新料，⛔ 不拿既有钉抵账）

既有钉＝`internal/panel/l2_grant_boundary_test.go`（`HEAD:internal/panel/l2_grant_boundary_test.go:1964`/`:1965` 拿页面原始串要求 `knownComposerMethod` 拒、`:1251-1267` 要求拒 11 枚候选名，第一枚即页面 `panel.approval.request`）。那族判的是 **"名册函数认不认识这个名字"**——**在 `internal/panel` 包内直调 `knownComposerMethod`/`ParseComposerRequest`**，不经任何传输、不跑 JS、不产生页面可见的拒句。

本程 (d)-1 的新面：同枚 `panel.approval.request`（页面今天真发的、`panel.ts:179-184` 原文 `{method, correlationId, outcome}`）**经由页面那条边进来**（甲③ 转发钩子 → 绑定门 → `dispatchRaw` → `Handle` → `ParseComposerRequest`），且 `door reply` **点名是名册那一道**（`不是面板 composer 通路的能力入口`），并由 `doorRounds==1`＋`spy.reqs==0` 钉住"它确实走完了这段传输、却在守卫处、在路由之前被拒"。既有钉走不到这条边，抵不了这格；这格也不覆盖既有钉。

---

## ⑤ 是否有一支落进 `A684` 那两面恒真

`A684` §3 记的这把尺两面恒真＝(M-A) 夹具 `fakeDoc35r2` 的原生 `postMessage` 桩 `:1238` 丢弃 receiver（`native.call(cw,msg)` vs 裸 `native(msg)` 两发都绿）、(M-B) `:137 jsObject.set` 无可写性概念（删 `typeof…!=="function"` 半支静默忽略覆写）。那两面归另一程（票 `:75` 新框），不归本腿修。

本腿三支**都不落进那两面**：每支只依赖"信封作为不透明串经转发钩子 → 门 → `Handle` → `ParseComposerRequest` → 回复串"这条**已被 `A683`/`A684` 复量、确有牙的**投递面（读数 `doorRounds`/`lastReply`），不碰 receiver 绑定、不碰属性可写性。定向突变（§②）三支各点各的红，正是这条投递面有辨别力的正面证据。⚠ 一个**依赖前提**具名：本三支的到达（`doorRounds==1`）建立在"转发钩子确实把裸信封折进了绑定门"这一支上，而"真 `chrome.webview.postMessage` 允许被钩子覆写"那一面今天只有 `-tags winlive` 真窗取过数（票 `:75`/`A684`），在本尺里是**被 M-B 那面盲区覆盖的**——即若 M-B 成立（覆写被静默忽略），本三支在本尺内会一致地绿而不自知。故本三支证的是"守卫次第＋各支各的文案"这一维，⛔ 不得被读成"传输可覆写性"的证据（那一维恒真面已知、另程处置）。此点如实标注，不当"有牙"。

---

## ⑥ 门禁数 + 逐名作差

- `sh scripts/d22scan.sh`：**rc=0** clean（`cmd/` 覆盖 108 枚 Go 文件含本腿新判据；`ban #8 emoji` 对 comments 豁免，我的 `⇒`/`⛔`/`★` 全在注释里、字符串字面量零符号）。件 `logs/gate-d22scan.txt`。
- `go vet ./cmd/wisp/ ./internal/panel/`：**rc=0**（成功＝零输出，rc 行自带）。件 `logs/gate-vet.txt`。
- `go test -count=1 ./internal/...`：**rc=1**，逐名红名册（`logs/internal-reds.txt`，5 枚）＝起手已知 5 枚**集合作差空**：
  `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestC21TableColourRowsMatchTokensCSS`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`。间歇门 `TestResolvePerCallBudget` 本程 0 命中；间歇 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 不记本程账。（本腿新件在 `cmd/wisp`，不影响 `internal/...` 结果，故 internal 起终同一集合。）
- 本腿三支：未突变 harness 运行 3/3 **绿 rc=0**（`logs/mut-baseline-green.txt`）；`cmd/wisp` 整包 harness 运行在后台（`logs/cmd-wisp-full.txt`）——该包含大量依赖真窗/前台环境的既有红（`cold bring-up=-1.000`、`TestAC4FocusReturnToPriorWindowGap33r5` 等），⛔ 非本腿账，本腿只以 `-test.run TestInbound` 那发 3/3 绿为准。
- harness：`go test -c -o /d/tmp/wisp35r3/panel.test.exe ./cmd/wisp/`（build rc=0），三枚 sherpa DLL 拷到 exe 旁，CWD=`cmd/wisp`。⛔ 未开任何真实 WebView2 窗口（`-tags winlive` 归验收腿）。

---

## ⑦ git diff --numstat + commit 号 + 每笔文件枚数

（本行以下为 commit 后回填；锚笔 `9b2551f2` 已交 1 枚 `start.md`。交件笔见回报。）
- 交件笔触及：`cmd/wisp/panel_inbound_guards_35r3_test.go`（新增）、`.scratch/wisp/probes/35/r3/{impl.md,logs/*}`（新增）、`.scratch/wisp/issues/35-panel-bridge-c17.md`（**只追加**一段进展登记，⛔ 不改原句、不翻任何 AC 框）。
- 每笔 `git show --stat <号>` 数件数见回报正文。证据件一律 `.txt`／`.md`（无 `.out`）。

---

## ⑧ 与简报的冲突

- 简报 §1 三支行号与盘上 `bridge.go` **完全一致**（`:132/:133`、`:135/:136-137`、`:139/:140-141`；解析判 `:128-129` 已在射程外）。⛔ 无冲突。
- 简报 §2 (d)-2/(d)-3 说"名册内"用 `panel.mode.request`——与票 `:70`/`:72` 一致（同枚信封，去掉 requestId），无冲突。
- 简报对既有钉的表述（(d)-1 的"新料不是它被拒"）＝票 `:71` 原句，一致。
- 唯一需具名的口径：简报把入口要求写成"经 JS 行为尺那条路径"。本腿复用 r2 的 `fakeDoc35r2` 行为尺（同包、未改其断言）而非 r1 的 `forwardingInstalled()` 词面正控——后者是 r1 的正控#1（`panel_transport_35r1_test.go:127`），本腿用**真跑 JS** 的尺到达，判别力更强、且不依赖词面串，符合 `:64`「⛔ 不直调 dispatchRaw」。无冲突，只是把"哪把尺"钉清楚。
- ⛔ 本腿不改产码（`git diff HEAD -- cmd/wisp/panel_host_windows.go internal/panel/**` 应为空）、不翻 AC 框、`frontend/**` 零改。
