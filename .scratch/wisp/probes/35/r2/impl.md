# 35-r2 交件（票 35 AC#6 `:52`：甲③ 落地 + 判据从"读字符串"换成"行为尺"）

腿＝`35-r2`（产码腿）。起手锚 `a7f9781c`；我的 commit：`a885b532`（起手锚 1 枚件）／`286a7f30`（产码＋判据 2 枚件）／本笔（证据件）。
⛔ 未 push；⛔ 未开任何真实 WebView2 窗（真窗那一发归 `35-v2`）；⛔ 未翻任何 AC 框；⛔ 未动 `frontend/**`／`internal/panel/bridge.go`／`docs/**`／SLO／golden／`thresholds.go`／票面原句。
逐字读数全在同目录 `logs/`。

---

## ① 新钩子 JS 原文逐字 + 落点 + 形ⓐ 还是形ⓑ

落点＝`cmd/wisp/panel_host_windows.go`：注释块起 `:638`，声明 `:673`，注册行 `w.Init(panelPostMessageForwardInit)` `:700`（`installPanelTransport` 体 `:693`；`dispatchRaw` 现量 `:706`）。
**类型由 `const` 改为 `var`**，因为绑定名必须由 Go 侧常量注入 JS（简报对形ⓑ 的那条 ⛔ 我按同一理由也用到了形ⓐ 上，见下）。

运行时字符串逐字（`%[1]s` 由 `panelDispatchBinding`＝`"wispDispatch"`（`:80`）填入，填后就是第三行 `window.wispDispatch`）：

```js
(function () {
  var cw = window.chrome && window.chrome.webview;
  if (!cw || cw.__wispForwardInstalled) { return; }
  cw.__wispForwardInstalled = true;
  var native = cw.postMessage;
  var inside = false;
  cw.postMessage = function (message) {
    if (inside || typeof window.%[1]s !== "function") { return native.call(cw, message); }
    inside = true;
    try { return window.%[1]s(message); } finally { inside = false; }
  };
})();
```

Go 侧原文即 `fmt.Sprintf` 那一段，逐字见 `cmd/wisp/panel_host_windows.go:673-682`（本腿不重复贴，`git show 286a7f30` 可取）。

**选形＝ⓐ（再入守卫）。理由（四条，全部锚在库原文）：**

1. **帧不被二改。** 形ⓐ 把库自己产生的 RPC 帧**原样**交回原生出口——一次 `JSON.parse` 都不做。形ⓑ 要先 `JSON.parse` 再手造 `{id, method, params:[raw]}`，那就引入了"我的判别规则错了怎么办"这一层新病（页面今天发的 `panel.approval.request` 那枚 `outcome` 键形、以及任何将来带数字 `id` 的页面信封都可能被误判成帧）。甲③ 的原文要求是"库的帧原样交回原生出口"，ⓐ 是它的直读。
2. **守卫的正确性有源码级依据，不是猜测。** 桩体是同步的：`webview.go:462-478` 里 `var seq = RPC.nextSeq++` → `new Promise(...)` 登记 `RPC[seq]` → `window.external.invoke(JSON.stringify({...}))` → `return promise`；`invoke` 就是 `chromium.go:112` 那枚箭头，也是同步。所以"我此刻是不是正站在一次 postMessage 调用里"这一个布尔值足以区分"页面在发"与"库在发帧"，不需要看内容。
3. **seq 仍由库分配。** 形ⓑ 要自增一枚页面侧的 id，与库的 `RPC.nextSeq` 是两个计数器，撞号就会把 `_rpc` 槽错位（回执送到别人的 Promise）。形ⓐ 完全不碰 id。
4. **改名安全照样满足。** 绑定名以 `%[1]s` 从 `panelDispatchBinding` 注入，JS 文本里没有硬编码的 `wispDispatch` 字面量（`grep -c '"wispDispatch"' cmd/wisp/panel_host_windows.go` ＝ 只在 `:80` 那枚常量里）。⚠ 顺带保住的一格：填完后的运行时串仍逐字含 `chrome.webview` / `postMessage` / `window.wispDispatch` 三串 ⇒ `panel_transport_35r1_test.go:127` 那把旧词面正控不会被我改红（简报 ⛔ 不许删它既有断言）。

三条禁区（票 `:52`）现量见 §⑤。

---

## ② 行为模型的语义对拉表（四件事各对哪一行库源码）

模型＝`cmd/wisp/panel_transport_35r2_test.go`：`fakeDoc35r2`（文档＋库的管子＋msgcb）＋一枚**只为跑这三段脚本而存在**的 JS 子集解释器（子集清单写在该文件头注释里）。

| # | 简报要求建模的那件事 | 库里逐字的出处 | 模型里怎么做（⛔ 不是读字符串） |
|---|---|---|---|
| ① | `chrome.webview.postMessage` 是一枚**可被覆写的属性**，存的是"当前值" | `pkg/edge/chromium.go:112`（该属性就是库的唯一出口）；`webview.go:435→chromium.Init→AddScriptToExecuteOnDocumentCreated`（`chromium.go:130-136`） | `webview.set("postMessage", native)` 是一枚普通属性槽；钩子对它的赋值走 `jsSetProp`；`openDocument()` 每次新建一个 window（＝新文档重跑全部脚本）。槽的读取经 `cw.onGet` 拦截，**读取时快照"当时的值"**（`captured := cur`）⇒ `var native = cw.postMessage` 抓到的是原生出口，之后的读取抓到的才是包装器 |
| ② | `window.external.invoke` 在**调用时**去取那枚属性的当前值 | `chromium.go:112` 逐字：`window.external={invoke:s=>window.chrome.webview.postMessage(s)}` | **这行原文逐字由解释器执行**（常量 `libExternalShim35r2`，`openDocument` 先跑它再跑 host 排队的脚本）。箭头体里的 `window.chrome.webview.postMessage` 每一次调用都重新走一遍成员读取 ⇒ 再入环不是我写的，是脚本自己撞出来的 |
| ③ | `wispDispatch` 桩**同步**登记 `window._rpc[seq]` 再 `external.invoke`，返回 Promise | `webview.go:462-478`（`Bind` 注入的桩），拼接式 `webview.go:462`，`jsString`＝`webview.go:137` | `Bind` 用库的拼接式生成脚本原文（`libBindStubScript35r2`＋`libBindStubBody35r2`，token 序列＝库的，只把 raw string 的前导 Tab 归一成空格＝JS 不读），**由解释器执行**；`new Promise(executor)`／`RPC[seq]={resolve,reject}`／`RPC.nextSeq++`／`Array.prototype.slice.call(arguments)`／`JSON.stringify` 都是模型里的真对象 |
| ④ | `msgcb` 的分支（未绑名 ⇒ `resolve`，⛔ 不是 `reject`） | `webview.go:139-160`（`msgcb`）＋`:131-135`（`rpcMessage`）＋`:163-168`（`callbinding` 未绑名 `return nil, nil`）＋`:157`（成功支 Eval 文本） | `fakeDoc35r2.msgcb` 用同一个 `rpcMessage` 形状解析；未绑名 ⇒ `callbinding` 回 `(nil,nil)` ⇒ 拼出**库那一句原文** `window._rpc[0].resolve(null); window._rpc[0] = undefined` **交给解释器执行** ⇒ 页面没有 `_rpc[0]` ⇒ `evalThrew++`（＝A682 §2 说的"被吞掉的 TypeError"）。票 `:58-59` 那句旧 `reject` 描述我按盘面库原文作废（`reject` 只在参数形状坏时走 `:144-151`，模型里也有那一条支） |
| ⑤ | 再入深度帽 | 真浏览器里是调用栈上限（`RangeError: Maximum call stack size exceeded`） | `postMessageHopCap = 8`：`onGet` 返回的 hop 在**调用时**自增 `postDepth`、返回时自减；超过帽就记 `capTripped/capDepth/capMessage` 并抛 `jsPanic`。⛔ 是计数帽，不是墙钟（AGENTS.md §1.2 那一条） |

校准那一格（证明这把尺不是恒绿）＝`TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot`：**不装任何转发**，页面裸信封自己走到原生出口（`nativeCalls=1`）、msgcb 解析成功但名未绑（`deadSlots=1`）、Eval 文本逐字＝`window._rpc[0].resolve(null); window._rpc[0] = undefined`、且它自己 `evalThrew=1`、door 0 次、下游 0 次。

---

## ③ ★硬判据两发的逐字输出

### 3.1 旧钩子（`fb2fb802` 那段 JS 原文）喂进新尺 ⇒ 必红

两发都在，一把比一把硬：

**(a) 常驻的那枚反形钉**（`TestLegacySubShapeOneHookDiesInAReentryLoop`，逐字常量 `legacySubShapeOneHook`＝`fb2fb802:648` 那段）——它自己绿，但它打印的就是旧形的红读数：

```
panel_transport_35r2_test.go:1772: sub-shape ① under the behavioural yard: re-entered 9 levels (cap 8), native exit called 0 times, door fired 0 times, page error="RangeError: maximum call stack size exceeded - chrome.webview.postMessage re-entered 9 levels (cap 8) with the native exit called 0 times; a page in a real browser dies here"
--- PASS: TestLegacySubShapeOneHookDiesInAReentryLoop (0.00s)
```

**(b) 把旧形放进产码那一格真跑到达性断言 ⇒ 到达性那把尺自己红**（突变 M2：`w.Init(panelPostMessageForwardInit)` → `w.Init(legacySubShapeOneHook)`；逐字 `logs/mut-m2.txt`）：

```
m2 run-rc=1   fails=3   passes=0
    panel_transport_35r2_test.go:1677: AC#6 RED: the page's own postMessage threw in the document: RangeError: maximum call stack size exceeded - chrome.webview.postMessage re-entered 9 levels (cap 8) with the native exit called 0 times; a page in a real browser dies here. nativeExitCalls=0 doorRounds=0 maxPostDepth=9 capTripped=true capDepth=9 unparsable=0 unboundSlots=0 evalThrew=0 resolved=0 rejected=0 scriptThrows=[]
    panel_transport_35r2_test.go:1709: AC#6 RED: page post #1 threw: RangeError: ... re-entered 9 levels (cap 8) with the native exit called 0 times ...
```

⇒ **这条就是简报点名的那一问（恒真性第 108 条那一问）的答案**：判据不是"删了转发就红"，而是"形状错在新尺下也红"，且红句里带**再入层数**与**原生出口被调次数**。

### 3.2 甲③ 新钩子 ⇒ 绿（逐字 `logs/shape-a3-green.txt`，源自同一次运行）

```
--- PASS: TestPagePostMessageEnvelopeReachesDispatchRawViaTransport (0.00s)     ← 35r1 那枚旧判据，一字未动，照旧绿
    panel_transport_35r2_test.go:1683: shape ③ delivery: nativeExitCalls=1 doorRounds=1 maxPostDepth=3 capTripped=false capDepth=0 unparsable=0 unboundSlots=0 evalThrew=0 resolved=1 rejected=0 scriptThrows=[] | frame="{\"id\":1,\"method\":\"wispDispatch\",\"params\":[\"{\\\"method\\\":\\\"panel.mode.request\\\",...
--- PASS: TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3 (0.00s)
--- PASS: TestShapeA3SecondPagePostStillDelivers (0.00s)
--- PASS: TestShapeA3ForwardingHookIsIdempotentInOneDocument (0.00s)
--- PASS: TestLegacySubShapeOneHookDiesInAReentryLoop (0.00s)
--- PASS: TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot (0.00s)
run rc=0
```

断到的东西（不是词面）：原生出口**恰被调 1 次**；那一次携带的是**库自己生成的**帧 `{"id":1,"method":"wispDispatch","params:[...]}`，其 `params[0]` 逐字＝`pageModeRequestEnvelope`（＝`HEAD:frontend/src/lib/panel.ts:246` 那枚 `panel.mode.request` 原文，35r1 已对拉过）；门只开 1 次；`dispatchRaw` 下游 1 次且带 `requestId`/`method`/`source`（AC#7 (c) 那三枚字段）；`_rpc[1]` 被 `resolve` 恰 1 次；`maxPostDepth=3`（页面→invoke→native），远在帽下。

---

## ④ 定向突变（四发，逐字 `logs/mut-m1..m4.txt`）

还原三行尺：`md5-good=3995d6bbb20ec251651544c872f5d359` ＝ `md5-after=3995d6bbb20ec251651544c872f5d359`（逐字节同），`git diff --exit-code -- cmd/wisp/panel_host_windows.go` ＝ `git-diff-exit-code=0`。四发突变都是同一份 good 源码上 `sed` 出来再建一枚 mutant exe（⛔ 不动判据、⛔ 不动断言换绿）。

| 发 | 突变 | 该红的断言 | 现量 |
|---|---|---|---|
| M1 | `w.Init(panelPostMessageForwardInit)` → `w.Init("/* no forwarding */")`（＝简报点名的"删掉/中和 `w.Init` 那句"） | 到达性 | **fails=3 passes=0**；红句：`AC#6 RED: the bound door fired 0 time(s), want 1. nativeExitCalls=1 ... unboundSlots=1 evalThrew=1`——即票 `:52-59` 那枚死法自己走出来了（信封到了原生出口、msgcb 解析得过、名未绑、`_rpc[0]` 不存在） |
| M2 | 把**旧形**放进产码那一格：`w.Init(legacySubShapeOneHook)` | 到达性 | **fails=3**；`capTripped=true capDepth=9 nativeExitCalls=0 doorRounds=0`（§3.1(b)） |
| M3 | 中和"先存原生出口"：`return native.call(cw, message)` → `return cw.postMessage(message)` | 到达性 | **fails=3**；同 M2 形：`capDepth=9 nativeExitCalls=0` ⇒ **甲③ 的那枚"先存"是真牙**，不是装饰 |
| M4 | 中和再入帽的归还：`finally { inside = false; }` → `finally { /* reset dropped */ }` | 第二发投递 | **fails=1 passes=1**：第一发仍绿（`TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3` PASS），第二发红 `door fired 1 time(s), want 2 ... unboundSlots=1 evalThrew=1` ⇒ "守卫不归还＝第二封页面信变成裸信封"这条病被单独钉住 |

---

## ⑤ 三条禁区的现量凭据（`git diff a885b532..HEAD`）

1. `git diff --stat a885b532..HEAD -- frontend/` ＝ 0 行；`--numstat` 亦 0 行 ⇒ **`frontend/**` 一字节未动**。
2. `git diff --stat a885b532..HEAD -- internal/panel/` ＝ 0 行 ⇒ `internal/panel/bridge.go` 未动。名册现量**六枚**（逐字，行号是本腿现跑的）：`:42` `MethodModeRequest = "panel.mode.request"`、`:43` `MethodWorkspaceRequest = "panel.workspace.request"`、`:44` `MethodAttachmentAdd = "panel.attachment.add"`、`:45` `MethodMessageSend = "panel.message.send"`、`:66` `MethodConfigGet = "config.get"`、`:67` `MethodConfigSet = "config.set"` ⇒ **⛔ 未加名**。
3. ⛔ 未把 RPC 解析器空答掉却不信封到达 `dispatchRaw`：`git grep -n "m\.dispatchRaw(" -- cmd/wisp` ＝ **恰 1 枚**（`panel_host_windows.go:695`，仍在 `installPanelTransport` 的绑定闭包里）。模型里未绑名的那一支**照旧不叫门**（§② ④／校准用例 `doorRounds=0`），新支路一条也没加。

---

## ⑥ 门禁数 + 逐名红名册作差

| 尺 | 读数 | 件 |
|---|---|---|
| `sh scripts/d22scan.sh` | **rc=0**，末行逐字＝`d22scan: clean - no D22 ban violations; live scope work: ... ban #8 cmd/=106`（cmd/ 从 v1 的 105 枚变 106＝本腿新增那枚判据文件；⛔ 零违例、零 SKIP） | `logs/gate-d22scan.txt`（251 行） |
| `go vet ./cmd/wisp/ ./internal/panel/` | **rc=0**，输出 0 字节 | `logs/gate-vet.txt` |
| `go test -count=1 ./internal/...` 第一发 | rc=1，逐名红＝6 | — |
| `go test -count=1 ./internal/...` 终态第二发 | rc=1，逐名红＝6；**起终作差 `DIFF-EMPTY`（roster-diff-rc=0）** | `logs/gate-internal-roster-diff.txt`（0 字节＝空差，逐字名册在下一行） |
| 终态逐名红名册（6 枚＝已知 5 ＋ 间歇 1） | `TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestC21DesignTokensFourWayAgree`、`TestC21TableColourRowsMatchTokensCSS`、`TestComposerContractTypesMatchFrontend`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（已知既有红）＋ `TestResolvePerCallBudget`（简报点名的间歇门，这两发起终都出现＝按规矩不计账，⛔ 我也没为它改任何东西） | 同上 |
| `cmd/wisp` 整包（harness 两发 exe，⛔ 不带 `-tags winlive`） | 基线件（`a885b532` 建的 `panel-base.exe`）：PASS=233，红/跳名册 6 枚；本程件（`panel-r2.exe`）：PASS=237（＝233＋我新交 5 枚里在整包中另计的那几枚，逐名见 `logs/cmd-wisp-suite-r2.txt`），红/跳名册 7 枚；**名册作差只多出 1 枚** `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` | `logs/cmd-wisp-suite-baseline.txt`／`logs/cmd-wisp-suite-r2.txt`／`logs/cmd-wisp-suite-roster-diff.txt` |
| 那枚多出来的红（定性＝批次时序 flake，⛔ 与本程无关，已复跑取证） | `cmd/wisp/config_reload_223_test.go:535` 那枚用例**不碰传输**：`grep -A30` 它体内 `installPanelTransport\|dispatchRaw\|Eval\|wispDispatch\|bringUp` 命中 **0**。单独复跑：**r2 exe 3/3 rc=0**、**base exe 2/2 rc=0**；它在整包里红的那发读数 `--- FAIL: ... (2.08s)`（走的是 L2 再确认＋秒级等待，同包 1700+ 行日志的批处理时序敏感）。⇒ 登记为既有间歇，不计成本程的账；`./internal/...` 那五枚已知红与此无关 | 复跑件 `/d/tmp/wisp35r2/t223-r2-{1,2,3}.txt`、`t223-base-{1,2}.txt`（临时件只建不删，⛔ 不在仓里） |
| 基线就有的 5 枚真窗族红（两发同现，非本程造成） | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`、`TestPanelHostRealWindowHopAndLifecycle`（＋`TestTicket223...` 见上一行）；另有 1 枚 SKIP `TestPanelHostLatencyPercentilesAC2` | 同上两份 suite 件 |
| 本机 `go test ./cmd/wisp/` 直跑 | 照简报：载入即死 `exit status 0xc0000135`（缺原生 DLL）⇒ 本腿全程用 harness（`go test -c` ＋ 三枚 DLL 置 exe 旁 ＋ CWD＝`cmd/wisp`），harness 建件 6 枚全 rc=0（基线、r2、M1-M4） | `/d/tmp/wisp35r2/build-*.txt`、`logs/mut-*.txt` 首行 |

---

## ⑦ `git diff --numstat` + 每笔 commit 与件数

| 笔 | 内容 | 文件枚数（`git show --name-only` 自数） |
|---|---|---|
| `a885b532` | 起手锚 `start.md` | 1 |
| `286a7f30` | `cmd/wisp/panel_host_windows.go`（＋42/−13）＋`cmd/wisp/panel_transport_35r2_test.go`（＋1824/−0，新） | 2 |
| 本笔 | `impl.md` + `logs/`（M1-M4、绿发、旧形反形钉、校准、门禁、两套 cmd/wisp 全量、名册作差） | 见 `git show --stat` |

`git diff --numstat a885b532..286a7f30`（已入库部分）：

```
42	13	cmd/wisp/panel_host_windows.go
1824	0	cmd/wisp/panel_transport_35r2_test.go
```

---

## ⑧ 与简报的冲突或过期（⛔ 我按盘面做，具名报回）

1. **简报的行号全部现量成立**：`panel_host_windows.go:648` 当时确实是那枚 `const panelPostMessageForwardInit`（改后注释块起 `:638`、声明 `:673`、`w.Init` `:700`）；`dispatchRaw` ＝ `:706`（＝v1 报的 `:677` 再 +29，因为我在这段里写了 29 行新注释）。⇒ 我这笔之后的引用一律用现量行号。
2. **简报对 `msgcb` 未绑名的更正（A682 §2）我按盘面执行**：模型走 `resolve(null)` 那一支（`webview.go:157`），⛔ 没按票 `:58-59` 的旧句写成 `reject`。票面原句我一个字没动。
3. **简报"名册六枚"与盘面一致**（`bridge.go:42-45`＋`:66-67`），与 v1 的读数也一致；A682 §4 第①条那处顶回已生效，本轮没有新冲突。
4. **⛔ 我没有放宽任何既有断言，也没有删 `panel_transport_35r1_test.go` 任何东西**：`git diff a885b532..HEAD -- cmd/wisp/panel_transport_35r1_test.go` 空，且它在 M1 下仍按自己的词面正控红（那是它既有形状，本腿不改它）。
5. **本腿自报一处缺陷（不可修，⛔ 不 amend）**：`a885b532` 那笔的 commit **信息末尾**被我把一段 shell 片段写进去了（`-- cmd/wisp/panel_host_windows.go 2>/dev/null; true`）。commit 的**内容**是干净的（`git show --name-only` ＝ 只有 `start.md` 1 枚），只有文字脏。按 git 纪律（已推送/已入库历史不改写、要更正就追加）我不 amend，在此具名登记。
6. **一处判断我扩了简报**：简报只要求形ⓑ 用 Go 侧常量注入绑定名；我在形ⓐ 也这么做了（`fmt.Sprintf` ＋ `const`→`var`）。代价＝该常量不再是编译期常量（本文件内只有 `installPanelTransport` 一处用它，`grep -n panelPostMessageForwardInit` 现量：`:638` 注释、`:673` 声明、`:700` 使用）。如果编排者认为形ⓐ 不必如此，回退只需把 `%[1]s` 换回字面量并把 `var` 换回 `const`，**判据不用动**（旧形反形钉与四条到达性断言都不依赖这个选择）。
7. **`chrome.webview.postMessage` 在真 WebView2 里到底让不让覆写（是不是 configurable）＝本腿仍未证**，这是 A682 §2 末句留的第二支。行为尺只证"若可覆写，则甲③ 送达、旧形死循环"；那一发的真窗凭据归 `35-v2`（`-tags winlive`，用 `bringUp` 的真实载体，⛔ 不许 `Eval` 扮页面）。若真窗测出"不可覆写"，甲③ 需要退回另择形，本文件那枚 `legacy` 反形钉仍成立、不白写。
8. 行为模型自带一枚 JS 子集解释器（新判据文件 1824 行的主体），子集清单写在文件头注释里；它只为让"库的两段脚本"和"产码的一段脚本"在同一次调用栈里相遇而存在，⛔ 不是通用引擎（无闭包词法外的特性、无原型链、无正则、无 try/catch 只有 finally）。评审若判它"过度"，减不出更小的能同时满足 §②四件事与 §③ 硬判据的形状。

（完）
