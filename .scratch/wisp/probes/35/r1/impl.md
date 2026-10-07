# 35-r1 交件 — 票 35 AC#6（`:52` 那格，页面↔宿主传输）

只读锚＝起手 `git log -1` = `10899b77`；本件只碰票 35 AC#6（`:52`）。AC#7（`:63` 判据框）与三支
(d)-1/-2/-3 红文案归下一枚腿，本格不加其断言。全程只 commit、⛔ 未 push。

---

## ① 转发语句落在哪一行、原码逐字（`cmd/wisp/panel_host_windows.go`）

选形＝票面 `:150` 编排者已裁的**甲子形①**（转发原始串进**已绑定的** `window.wispDispatch(raw)`，
走库的 binding 回执路）。三条新符号＋一处调用，行号现量：

- `:408` `bringUp` 的入向门改为 `if bindErr := m.installPanelTransport(w, ctx); bindErr != nil {`
  （原 `:405` 的内联 `w.Bind(...)` 挪进 installPanelTransport）。
- `:633` `type pageTransport interface { Bind(name string, f interface{}) error; Init(js string) }`
  ——窄接口，只为让这段 shipping 接线在无浏览器下可被驱动；`webview2.WebView` 结构上满足它。
- `:648` `const panelPostMessageForwardInit`（就是那条转发语句，逐字）：

```
`(function () {
  var cw = window.chrome && window.chrome.webview;
  if (!cw || cw.__wispForwardInstalled) { return; }
  cw.__wispForwardInstalled = true;
  cw.postMessage = function (message) {
    if (typeof window.wispDispatch === "function") { window.wispDispatch(message); }
  };
})();`
```

- `:664` `func (m *PanelManager) installPanelTransport(w pageTransport, ctx context.Context) error`，
  体内 `:671` 即**那条转发语句** `w.Init(panelPostMessageForwardInit)`（先 `w.Bind(panelDispatchBinding,
  func(raw string) string { reply, _ := m.dispatchRaw(ctx, raw); return reply })` 再 Init）。

时序：`installPanelTransport` 在 `bringUp` 里跑于 `firstRoundTripLocked`/`serveEntry` 之前 ⇒ `Init` 早于
首个 `SetHtml`（票 `:144` 的 `:449/:468/:681` 约束）；Init 脚本随每次文档创建重跑，活得过后续 `SetHtml`
重建；`__wispForwardInstalled` 守卫防同一文档内重复挂钩。⛔ 不新造 `C##`、⛔ 不给 C17 名册加名
（转发只把页面既有信封喂进既有门）、⛔ 未答 RPC 解析器绕开 `dispatchRaw`。

## ② 新用例的名字 + 绿那发逐字行

文件：`cmd/wisp/panel_transport_35r1_test.go`（`//go:build windows`）。
用例：`TestPagePostMessageEnvelopeReachesDispatchRawViaTransport`。
入口＝页面真用的那条边：`fakePageControl.deliverPostMessage(<panel.ts:246 在册信封原文>)`
（`{"method":"panel.mode.request","requestId":"pc-1-...","source":"panel-composer","to":"ask_every_step"}`，
正控取自 `git show HEAD:frontend/src/lib/panel.ts:246`，⛔ 未进工作树别处、未改 `frontend/**`）。
fake 的 `msgcb` 逐字复刻库的 `{id,method,params}` 解析＋只路由到已绑方法，故：⛔ 未直调 `dispatchRaw`、
⛔ 未用 `w.Eval("window.wispDispatch(...)")`（N6 那一形，票 `:65`/`:66` 明令禁）。
断 (c)：`dispatchRaw` 下游（`ComposerDispatch.Handle` → `Mode.HandleModeRequest`）被叫**恰好一次**、
带 `req.RequestID == "pc-1-3f2b1c0d-9e7a-4c1b-8f14-e45fceea469a"`、method/source 未被途中篡改、`rpcDeadSlot==0`。

绿那发逐字行：
```
=== RUN   TestPagePostMessageEnvelopeReachesDispatchRawViaTransport
--- PASS: TestPagePostMessageEnvelopeReachesDispatchRawViaTransport (0.00s)
```
（`logs/mutation-restored-green.txt`；整包 `-test.v` 内同用例亦 `--- PASS`，见 `cmd-suite.txt` 抓取。）

## ③ ★删掉转发语句那一发的逐字红输出 + rc（本格的牙）

变异＝把 `:671` 的 `w.Init(panelPostMessageForwardInit)` 换成 `_ = panelPostMessageForwardInit`
（转发不安装），恢复＝逐字节回到 commit `fb2fb802`（`git diff` 空）。红那发逐字：
```
=== RUN   TestPagePostMessageEnvelopeReachesDispatchRawViaTransport
    panel_transport_35r1_test.go:209: AC#6 RED: installPanelTransport registered no page->wispDispatch forwarding hook (no Init string reroutes window.chrome.webview.postMessage into window.wispDispatch). Without it the page's envelope is mis-parsed by msgcb and never reaches dispatchRaw (ticket 35:52-59).
--- FAIL: TestPagePostMessageEnvelopeReachesDispatchRawViaTransport (0.00s)
FAIL
```
**run rc = 1**（红）。⇒ 恒真性成立：判据读的就是这条边，删了转发它就红。
（`logs/mutation-deleted-red.txt`／`logs/mutation-restored-green.txt`。N1/N2/N3/N6 一律未用作证据。）

## ④ 三支守卫文案（`bridge.go:133`/`:136-137`/`:140-141`）覆盖情况

**本格未覆盖任何一支——具名写"未覆盖"：**
- `bridge.go:133`（不在册，"方法 … 不是面板 composer 通路的能力入口"）＝**未覆盖**（AC#7 的 (d)-1）。
- `bridge.go:136-137`（来源，"按伪造/串台拒绝"）＝**未覆盖**（AC#7 的 (d)-2）。
- `bridge.go:140-141`（缺 requestId，"缺少 requestId"）＝**未覆盖**（AC#7 的 (d)-3）。

原因（票 `:69-72`）：本格的正控信封 `panel.mode.request` 是**在册名**，在名册守卫 `:132` 处**判真通过**，
根本走不到那三道 `return`；名册守卫在最前，所以它结构上不可能触发任何一支拒。三支红各点各的断言属
AC#7，归下一枚腿。`internal/panel/bridge.go` 一字未动 ⇒ 三支文案仍是各点各的、没被我合并（见票
"一个门不许把几桩病压成一格"）。`bridge.go:133` 那格"未注册名被拒"已有既有钉（`l2_grant_boundary_test.go`
`1964/1965`/`1251-1267`），本枚不重复、也不抵账。

## ⑤ 门禁数（⛔ 按名字集合比，不按枚数）

- `sh scripts/d22scan.sh`：**rc=0**；自带正控 `PASS=35 FAIL=0 SKIP=0`（`logs/d22scan-final.txt`，
  含 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77`）；扫描 clean，
  本格两枚文件都在 `cmd/` 射程内被读到（ban #8 cmd/=105 Go files）。
- `go vet ./cmd/wisp/`：起手 **rc=0**（elapsed ≈1.6s，`logs/vet-start.txt`）。
- `go vet ./cmd/wisp/ ./internal/panel/`：终态 **rc=0**（`logs/vet-final.txt`）。
- 整包 `go test -count=1 ./cmd/wisp ./internal/...` **两发**（起手 / 终态，各 -v 落文件）：
  - 起手 rc=1、终态 rc=1（内部既有红，见下）。
  - **逐名红名册起手集合 == 终态集合**（作差 = ∅）。两发红名（按名）：
    `TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestC21DesignTokensFourWayAgree`、
    `TestC21TableColourRowsMatchTokensCSS`、`TestComposerContractTypesMatchFrontend`（＝票面点名的
    `composer_test.go:73-77` 既有红，**本格没碰它、它也不归本格**）、
    `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestResolvePerCallBudget`。
  - `TestResolvePerCallBudget`（`internal/risk/pathresolver_budget_norace_test.go`）＝贴预算线的间歇门，
    这两发恰都出现；按简报它不计账（出现/消失都不算）。本格**未新增任何具名红**。
- ⚠ **`cmd/wisp` 整包在这台机器跑不起来＝`exit status 0xc0000135`（STATUS_DLL_NOT_FOUND）**，起手与
  终态同发同现：`go test -c ./cmd/wisp` 能编（rc=0），但测试 exe 载入时找不到 sherpa/onnxruntime 原生
  DLL（在 `%GOMODCACHE%\...\sherpa-onnx-go-windows@v1.13.8\lib\x86_64-pc-windows-gnu`）。这是本机 PATH
  条件、非某条断言、非本格产物。为把 `cmd/wisp` 测起来，我把该三枚 DLL 放到测试 exe 旁（Windows 先在
  exe 目录搜 DLL），并以 `cmd/wisp` 为 CWD 跑（源码扫描类用例需要）。这是临时件（只建不删，均在 `/d/tmp`，
  ⛔ 未落进仓）。
- DLL-harness 里 `cmd/wisp` 全族跑：**本用例确定性 PASS**；但若干**真窗口生命周期用例是批次不确定的**
  （叠在一进程里跑多个真实 WebView2 窗口 → 焦点/泵互扰，20s 超时）：with-change 那发红了
  {`TestPanelHostRealWindowHopAndLifecycle`,`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`,
  `TestAC14AwaitedBindingReplyReachesThePage`,`TestAC14GoSideEvalPushReachesThePage`}；而**原始码**（不含我
  改动）单独跑这四枚时，`TestPanelHostRealWindowHopAndLifecycle` 红、三枚 AC13/AC14 **绿**；原始码整包跑又
  先红了另一枚 `TestAC4FocusReturnToPriorWindowGap33r5`。⇒ 红哪一枚在变、与我的改动无关（我的改动只在
  `bringUp` 里多一句 `w.Init`，不建窗不抢焦点）。我据此起的那份"原始码整包"后台跑到 memory/session 段
  被我 TaskStop（同机 self-hosted runner 会抢 CPU），未取完——**具名：原始码整包 `cmd/wisp` 全族红名册
  未跑到尾，故 with-change vs baseline 的整包逐名差集未坐实**；坐实的是"我的用例确定性绿＋删转发确定性红＋
  这些窗口用例本就在原始码上也随机红/绿"。

## ⑥ `git diff --numstat` + commit 号 + 每笔文件枚数（越界写面 = 0）

起手锚 `10899b77` → 本腿三笔：
- `f718e9b6` 起手骨架：**3 枚**（均 `.scratch/wisp/probes/35/r1/`：start.md + logs/vet-start.txt +
  logs/porcelain-start.txt），754+6+23 行；⛔ 未碰 Go。
- `fb2fb802` 实现：**2 枚**（`cmd/wisp/panel_host_windows.go` +52/-5、`cmd/wisp/panel_transport_35r1_test.go`
  +243/-0）。
- 本笔交件：`.scratch/wisp/probes/35/r1/logs/*`（d22scan-final、vet-final、test-start/final、red-start/final、
  mutation-deleted-red / mutation-restored-green、cmdwisp-suite-* 两份红名册）＋本 impl.md。
- `git diff --numstat f718e9b6~1 HEAD` 落在代码面只有 `cmd/wisp` 那两枚；`git diff --stat 10899b77 -- frontend/`
  = **空**（`frontend/**` 一字节未动）。票面/台账/`docs/**`/`design/**`/`scripts/**`/`go.mod` 均未动 ⇒ **越界写面 = 0**。

## ⑦ 与本简报任何一处冲突／过期（以盘上原文为准）

1. **简报"AC#7 归下一枚腿"与 §2/§5-④ 要求点名三支守卫的张力**：我按票 `:63`+`:69-72` 取"本格只做
   (c) 到达性正控，三支 (d) 红文案属 AC#7"，未加其断言；§5-④ 要求的"三支各点各的有没有被我覆盖"已如实
   写**三支皆未覆盖**。若编排者其实要 35-r1 就把三支也测掉，那与"AC#7 归下一枚腿"矛盾，我按后者（不动 AC 框）。
2. **简报 `:1` "在已绑定的 `window.wispDispatch` 之前加一条转发语句"的措辞**与票 `:150` 甲子形①（"转发
   原始串进**已绑定的** `window.wispDispatch(raw)`"）读起来方向不同：我按票 `:150` 的功能裁——转发把页面
   `postMessage` 喂进**已存在**的 `window.wispDispatch`；`installPanelTransport` 里先 `w.Bind`（库注入
   wispDispatch 桩）后 `w.Init`（转发钩子），二者都是 AddScriptToExecuteOnDocumentCreated、按加入顺序在
   document-created 执行，转发钩子里对 `window.wispDispatch` 的调用发生在页面运行时（页面真发 postMessage 时），
   与注册先后无关 ⇒ 功能正确。简报那句"之前"我按"在页面用到之前装好"理解，未把它读成 JS 语句字面顺序。
3. **简报/只读腿 §1 的行号我复跑坐实**：`:52`=119、`:63`=225、`:69`=731、`:73`=848（`awk` 量）——
   裁定三支的正文确在 `:69-73`；`git grep -c wispDispatch HEAD -- frontend`＝**0**（页面从不调 wispDispatch，
   与 `:52` 一致）；`bridge.go` 守卫次第 解析`:128`→名册判`:132`拒`:133`→来源判`:135`拒`:136-137`→requestId
   判`:139`拒`:140-141`，我照 `:69` 逐字对过（未合并文案）。
4. **简报"终态整包 go test ./cmd/wisp ./internal/..." 在这台机上 `cmd/wisp` 发根本载入不了**
   （0xc0000135）——若编排者预期这一格 `cmd/wisp` 会"跑出红名册"，那是**过期/环境不符**：它在起手与终态都
   只能以 `0xc0000135` 出现，`cmd/wisp` 对逐名红名册的**贡献恒为空**。我以"红名册按名集合比＝∅ 差"交付，
   另附 DLL-harness 的窗口用例批次不确定性说明。
5. 未跑完项具名：原始码 `cmd/wisp` 整包全族红名册（后台跑到 session 段被我停）；`cmd/wisp` 的 winlive 真机
   用例（本腿未加 `-tags winlive`，按票 `:154` 只登记〔仅本机可量〕）。这两项都不是 AC#6 的牙，牙（③）已跑。
