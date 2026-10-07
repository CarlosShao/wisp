# 票 35 · AC#6 验收腿 `35-v2`（非实现者，攻尺＋真窗那一发）

起手 HEAD `14d027df`（A683）。本腿不翻任何 AC 框、不改产码、不动判据。
证据件＝本目录 `logs/*.txt` 20 枚，每件自带 `rc=` 行（⛔ 无 0 字节件）。
交件 commit 号见本程末段。

---

## W1 ★攻那把 1824 行的行为尺 —— 判语：**不成立（这把尺有两面恒真，已用定向突变证明"改错了也绿"）**

打法：⛔ 不改仓里一字。`go test -c -overlay <json>`，overlay 只把 `cmd/wisp/panel_host_windows.go`
指到 `/d/tmp/wisp35v2/mut/m*.go` 的突变副本（每件 `diff` 恰 1 行，先证落地再看颜色，第 127 条）。

| 突变 | 那枚错法 | 五枚点名用例的结果（我这把尺的逐字读数） |
|---|---|---|
| 基线 | 未突变 | 5/5 PASS，rc=0；`nativeExitCalls=1 doorRounds=1 maxPostDepth=3 capTripped=false resolved=1 rejected=0`，`params[0]` 逐字＝`{"method":"panel.mode.request","requestId":"pc-1-3f2b1c0d-9e7a-4c1b-8f14-e45fceea469a","source":"panel-composer","to":"ask_every_step"}` |
| **M-A** | `native.call(cw, message)` → **`native(message)`**（丢 `this`） | **5/5 PASS rc=0**，读数与基线一字不差 |
| **M-B** | `if (inside \|\| typeof window.%[1]s !== "function")` → **`if (inside)`** | **5/5 PASS rc=0** |
| M-C1 | `inside` 在同步 hop **之前**归还 | **3 枚红** rc=1：`capTripped=true maxPostDepth=9 nativeExitCalls=0` |
| M-C2 | `finally { inside = inside; }`（块在、归还剩空） | **1 枚红** rc=1：`TestShapeA3SecondPagePostStillDelivers` |

**恒真面（具名）**

1. **`this` 那一面是恒真的。** 根因不在解释器而在夹具：模型里那枚原生出口桩
   `panel_transport_35r2_test.go:1238` 写的是 `host: func(_ jsValue, args []jsValue)`
   ——recv 参数被直接丢弃；`jsCallValue`（:346）虽然把 recv 传到底，**没有任何一处读它**。
   ⇒ M-A 全绿。真浏览器里 `chrome.webview.postMessage` 是宿主方法，脱离接收者调用按
   Chromium 惯例是 `TypeError: Illegal invocation`，这把尺表达不出这个形状。
2. **属性可写性这一面根本不可表达。** `jsObject.set`（:137）无条件 `o.props[name] = v`；
   全文 `grep -cE 'writable|configurable|defineProperty|Illegal'` = **0**。
   ⇒ 甲③ 的全部机制建立在"覆写 `cw.postMessage` 会生效"上，而"覆写被静默忽略"（sloppy 模式
   对非 writable 属性就是静默失败）在这把尺下**必绿**。这一条不是理论风险：它就是 W2 那发真窗。
3. **`typeof window.wispDispatch !== "function"` 那一支是死支。** 本文件没有任何"门未绑"用例
   （`TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot` 红的是 `msgcb` 那一侧的未绑方法名，不是门缺失），
   删掉它 5/5 照绿 ⇒ 那半支的行为在真窗里"钩子先于 Bind 跑"的形状，尺看不见。

**这一面确有牙（A682 要求的那一面，我复认）**：M-C1 / M-C2 各红在指名用例上；
`TestLegacySubShapeOneHookDiesInAReentryLoop` 报 `re-entered 9 levels (cap 8), native exit called 0 times`。
`try/finally` 与 `Function.prototype.call` 都**真实现了**（:633-651 的 try 支连"finally 跑完后继续抛"都对），
成员调用也确实带接收者（:565-588）——恒真面是**夹具**的，不是解释器的。

**三处对拉（库原文本腿自取，`go env GOMODCACHE` = `D:\work\base\gopath\pkg\mod`，
`github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808`）**
- `pkg/edge/chromium.go:112` 逐字＝`e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")`
  ⇒ 垫片在 **call 时**读属性，与测试里 `libExternalShim35r2` 一致。✔
- `webview.go` Bind 注入的桩（实测 `:449-478`，票与写腿写 `462-478`——**行号漂移，具名报回**）：
  同步 `RPC[seq]={resolve,reject}` → `window.external.invoke(JSON.stringify({id:seq,method:name,params:Array.prototype.slice.call(arguments)}))` → `return promise`。✔ 与 `libBindStubBody35r2` 一致。
- `msgcb`（实测 `:140-168`）：`callbinding` 未绑名 ⇒ `return nil, nil` ⇒ `err != nil` 不成立 ⇒ **成功支**
  `w.Eval("window._rpc[0].resolve(null); ...")`。⇒ **票 `:58-59` 那句 `reject` 是过期句，`A682` 的更正成立**（我按库源码复核，不是按转述）。
- ⚠ 一处模型比浏览器**严**：`onGet`（:1210-1235）把"读一次属性"也计一层深度，所以真实形状的
  `maxPostDepth=3`（页面 1 ＋ 钩子自己那一次 `var native` 1 ＋ 垫片 invoke 里那次读 1）而不是浏览器的 1。
  帽 8 因此是量级代理、不是栈深移植；不构成假绿，但"cap 8"这个数不可移植。

---

## W2 ★真 WebView2 那一发 —— 判语：**成立（真机到达已取到，一发即中）**

载体＝产品路径，逐字：`PanelManager.bringUp`（`cmd/wisp/panel_host_windows.go:318`）→
`installPanelTransport`（调用点 `:408`，函数 `:693`，`w.Init(panelPostMessageForwardInit)` 在 `:700`）→
稳态泵＝**库自己的 `WebView.Run()`**（与 `cmd/wisp/panel_resident_windows.go:299` 同形），跑在
`runtime.LockOSThread()` 且永不解锁的那枚线程上（对齐本仓定案：面板归专用 STA 线程＋库 `Run()` 泵）。
⛔ 本腿不起自泵、⛔ 不调任何产品路径没有的 helper、一次只开一枚窗。

新台件一枚：`cmd/wisp/panel_transport_live_35v2_windows_test.go`（`//go:build windows && winlive`，
文件头写明是验收台件；⛔ 未碰任何既有判据、⛔ 未碰产码）。

**逐字读数（`logs/live-run1.txt`，rc=0，`--- PASS ... (8.25s)`）**
```
35v2 cold bring-up measured on this box: -1.000 ms, hwnd=0x8a0a0e
35v2 canary arrival=true after 1 post(s); door has recorded [pc-35v2-canary]
35v2 DELIVERED TO GO: method="panel.mode.request" requestId="pc-35v2-pagepost-1" source="panel-composer" to="ask_every_step" (attempt 1)
35v2 PAGE-SIDE FACTS (diagnostic lane, asserted for nothing): pc-35v2-DESC-OWN-WR-CF-HOOK-NOERR
35v2 all door recordings: [pc-35v2-canary pc-35v2-pagepost-1 pc-35v2-desc-ping pc-35v2-DESC-OWN-WR-CF-HOOK-NOERR]
35v2 webview processes at start: tree=0 machineNamed=29 pids=0
35v2 webview processes after teardown: tree=0 machineNamed=29 pids=0 (start tree=0 pids=0)
```
判读：**页侧 `chrome.webview.postMessage(<信封>)` 到达 Go 门，第一发即到**。这条到达本身就把
"覆写生效"反推出来了——若钩子没挂上，那枚串会原样进 `msgcb`、`method="panel.mode.request"` 无人绑、
走 `resolve(null)` 支，我的记录仪**一枚都收不到**。诊断道回带 `OWN-WR-CF-HOOK-NOERR`
＝覆写后的属性是自有·可写·可配置、钩子的 `__wispForwardInstalled` 在场、经属性的间接调用不抛。
残留：进程树 0→0（machineNamed 29→29）。

**⛔ 这不是 N6 那一形**（票 `:64-65` 禁的是 `w.Eval("window.wispDispatch(…)")`＝测试在**门**上扮演页面，
对任何选形都恒绿）。我叫的是**属性**，判据全在 Go 侧记录仪。
**控制组（证明这发有辨别力）**：`git show fb2fb802:cmd/wisp/panel_host_windows.go` → overlay 换成甲① 的钩子，
同尺同载体：`logs/live-run-legacy-control.txt`（与并发跑时）与 `logs/live-run-legacy-control-2.txt`（空载复跑）
都是 **58.27s、`35v2 all door recordings: []`、`--- FAIL` rc=1**，外加 3 条"渲染串缺甲③ 形状"的自证红。
⇒ 甲① 在真窗里 0 到达，甲③ 一发即中——这发判据两侧分得开。

**取不到的两格（具名，⛔ 不许写成"应该没问题"）**
1. 覆写**之前**那份原生 `postMessage` 的 descriptor：文档级脚本没有比 `installPanelTransport` 更早的插槽，
   测试拿不到"钩子跑之前"的那次读。所以"可写"是由到达**反推**的，不是直读。
2. 原生出口在丢 `this` 时到底抛不抛：属性已被覆写，测试侧再也取不到原生引用 ⇒ M-A 那枚真机风险
   **本腿无法在真窗里判定**，只能标"尺看不见＋真机未证"。
3. 附带一格非本腿射程的观察：`cold bring-up = -1.000 ms`＝我这枚线程在 `bringUp` 期间不起泵，
   产品那枚冷探针的往返要等 `Run()` 才有戳（票 33 AC#2 的形状问题，⛔ 不是 AC#6 的判据，交回你决定要不要开格）。

---

## W3 复量那枚 flake —— 判语：**"批次 flake" 定性成立；35-r2 给的凭据形状不成立**

我这把尺（`logs/pkg-panel-1.fails.txt`／`solo-panel-*.txt`，exe＝HEAD 全量含 r2）：
- 单跑 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` ×3：**3/3 PASS rc=0**（2.02／2.03／2.21s）。
- **同一枚 exe、同一棵树**整包两发：run1 红名册 6 枚（含 `TestTicket223…`），run2 红名册 **5 枚（不含它）**。
  ⇒ 决定性的不是"基线 vs r2"，而是"**同树两发不一致**"＝间歇，坐实"批次 flake"。
- 机制旁证（读码，⛔ 不替代实测）：`config_reload_223_test.go` 是排序第 6 枚文件、
  `panel_transport_35r2_test.go` 第 34 枚 ⇒ r2 的用例在 223 **之后**跑；r2 文件无 `init()`、无全局写入
  （只有自己的 `var jsPuncts`）⇒ 顺序运行里它没有碰 223 的通道。223 的红句是
  `config_reload_223_test.go:545: the card does not name risk.permission_mode`＝内容断言不是时延阈值，
  符合"被批次的共享态（mockllm／日志 sink／配置）碰到"那一族。
- ⛔ 与间歇门 `TestResolvePerCallBudget` 不同名不同物：本程 `gate-internal.txt` 里它 **0 命中**（未红）。
- **顶回一处凭据**：35-r2 报的凭据是"单跑 r2 3/3 绿、基线 2/2 绿"。"基线绿＋r2 侧红"这个形状本身
  更像**新红**而不是 flake——它证不了 flake。真凭据是"同树两发不一致"，本腿现量到了。
- 没做到的那一轴（具名）：我只摘了 r2 的**判据文件**（overlay 换成 `package main` 空件，`nor2.test.exe`），
  没单独隔离 r2 的**产码改动**那一轴；因"同树两发不一致"已足以排除"确定性新红"，我没再造第三枚 exe。
- 整包另有 5 枚红两发都在（`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、
  `TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`、
  `TestAC4FocusReturnToPriorWindowGap33r5`、`TestPanelHostRealWindowHopAndLifecycle`）＝
  `cmd/wisp` 整包在本机的既有红，⛔ 不是本程的账（它们写进红名册会让人以为票 35 弄坏了它们）。

---

## W4 禁区与渲染串 —— 判语：**成立**

- `git diff --numstat a7f9781c..HEAD -- frontend/` ＝ **空**；`-- internal/panel/` ＝ **空**（`logs/00-anchor.txt`）。
- `internal/panel/bridge.go` 未被这两笔碰：blob `bebe8e702a85c640641551e93dd09936530950f4`，且
  `git diff a7f9781c..HEAD -- internal/panel/bridge.go` 输出 0 行。
- **名册六枚逐字**（`internal/panel/bridge.go:42-45` ＋ `:66-67`，`case` 那一枚 `:148` 六名齐）：
  `panel.mode.request` · `panel.workspace.request` · `panel.attachment.add` · `panel.message.send` ·
  `config.get` · `config.set`。⚠ 具名报回：票 `:60`／写腿注释说"名册六枚"在数值上对，但那两枚 config 的名字
  **不是** `panel.config.*`——按 `panel.` 前缀去 grep 只会搜到 4 枚（我第一发就这么错，现量改 `MethodConfig`）。
- `m.dispatchRaw(` 在 `cmd/wisp/panel_host_windows.go` 计数 **1**。
- **运行时真实渲染串**（`logs/live-run1.txt` 第 3-16 行，台件里 `t.Logf(panelPostMessageForwardInit)` 打出来的，不是读源码）：
```
(function () {
  var cw = window.chrome && window.chrome.webview;
  if (!cw || cw.__wispForwardInstalled) { return; }
  cw.__wispForwardInstalled = true;
  var native = cw.postMessage;
  var inside = false;
  cw.postMessage = function (message) {
    if (inside || typeof window.wispDispatch !== "function") { return native.call(cw, message); }
    inside = true;
    try { return window.wispDispatch(message); } finally { inside = false; }
  };
})();
```
  `%[1]s` 三处全换成 `wispDispatch`、整串 `strings.Contains(rendered,"%")` = false（台件里是 Errorf，未触发）、
  首 `(function () {` 尾 `})();` 与 `:638` 注释许诺逐字一致，无游离引号、无被换行打断的字符串字面量。
- 行号对拉（本腿现量，供后续腿钉）：`:638` 注释 ✔／`:673` `var panelPostMessageForwardInit = fmt.Sprintf(` ✔／
  `:693` `func (m *PanelManager) installPanelTransport` ✔／`:700` `w.Init(panelPostMessageForwardInit)` ✔／
  模板体在 `:673-684`（`:680` 那一行是 `inside || typeof` 支）。票 `:54` 的 `:80`/`:405`、`:62` 与 `:66` 的 `:630`
  已漂到 `:80`（绑定名常量，未漂）／`:408`／`:696`——⛔ 本腿不改票面原句，只具名登记。

---

## 门禁（终态，件内自落 rc）

| 门 | 读数 |
|---|---|
| `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`（scope ban #8：internal/=514、cmd/=107）→ `logs/gate-d22scan.txt` |
| `go vet ./cmd/wisp/ ./internal/panel/` | **rc=0**（输出零行，件里那一行 `rc=0` 是本腿自落的）→ `logs/gate-vet.txt` |
| `go test -count=1 ./internal/...` | **rc=1**，红名册恰 **5 枚**＝已知 5 枚（`logs/gate-internal.fails.txt`）：`TestApprovalCardViewJSONKeysMatchFrontendTypes` · `TestC21DesignTokensFourWayAgree` · `TestC21TableColourRowsMatchTokensCSS` · `TestComposerContractTypesMatchFrontend` · `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`；与 `logs/known-red.txt` 按名字集合作差＝**空**。另两行 `github.com/CarlosShao/wisp/internal/{ball,panel}` 是包级 FAIL 行，不是用例名。间歇门 `TestResolvePerCallBudget` 本程 0 命中。 |
| 五枚点名行为尺（HEAD 全量，`panel.test.exe`，CWD＝`cmd/wisp`，DLL 三枚在 exe 旁） | **5/5 PASS rc=0** → `logs/ruler-named5.txt` |

本机可测性照单执行：`go test ./cmd/wisp/` 直接跑仍 `exit status 0xc0000135`，全部读数走
`go test -c` ＋ `/d/tmp/wisp35v2/`（`panel.test.exe` 非 winlive／`live.test.exe` ＝`-tags winlive`／
四枚突变 exe／`legacy.test.exe` 控制组），临时件只建不删、全在 `/d/tmp`。
`git status` 起手 754 枚（本机 autocrlf 幻影），产码 CR 计数 0＝工作树与 HEAD 一致。

---

## 总结论（你问的两问）

**① 票 35 AC#6 现在能不能翻？——能翻。** 凭据不是"解释器绿"，是 W2 那发真机到达：
产品 `bringUp` 装的钩子、库的 `Run()` 泵、页面侧 `chrome.webview.postMessage` 一发即中、
Go 侧门逐字段收到 `panel.mode.request`／`pc-35v2-pagepost-1`／`panel-composer`／`ask_every_step`，
且控制组（甲①）同尺同载体 0 到达＝这发判据两侧分得开；残留进程树 0→0。
翻框时建议带上三条限定：(a) 判据是"测试在页面里叫属性"，不是真人点鼠标（页面开机不发东西）；
(b) 到达＝逐字段相等，字节级 verbatim 仍只在 35r2 那把模型尺里；
(c) 甲③ 依赖"覆写生效"，这一条今天由到达反推、不是 descriptor 直读。

**② 那把 1824 行的行为尺有没有被攻出恒真面？——有，三面，其中两面已用突变证明。**
1. `this`／接收者：M-A（`native(message)`）5/5 全绿 rc=0 —— 夹具 `:1238` 把 recv 丢掉，尺表达不出 `Illegal invocation`。
2. 属性可写性：`:137 jsObject.set` 无条件成功、全文 0 处 writable/configurable/defineProperty —— "覆写被静默忽略"这一族必绿。
3. `typeof … !== "function"` 半支：M-B 删掉它 5/5 全绿 rc=0 —— 没有"门未绑"用例，那一支在本尺下是死码。
A682 真正要的那一面（再入／时序）**有牙**：M-C1 三枚红（capTripped=true、nativeExitCalls=0）、
M-C2 红在指名用例、旧形在本尺下报"再入 9 层、原生出口 0 次"并在真窗控制组里 0 到达。
⇒ 这把尺可以留在仓里当代替"读词面"的那把旧尺，但⛔ 不许再把它当"浏览器会不会同意"的凭据——
"浏览器同不同意"那一面今天只有 W2 那发真窗能判，而 winlive 在本仓不进 CI。
