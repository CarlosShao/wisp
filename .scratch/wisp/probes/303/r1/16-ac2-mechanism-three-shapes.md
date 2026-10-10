# 303-r1 `AC#2` — 机制三形分开答（甲/乙/丙）

承重读数均本腿现量，件＝`14-ac2-shape3-reentry-read.txt`（headless 行为尺）＋`20-targeted-before-red.txt`（母仓 HEAD 三枚真窗红）；
归因承重的成对＝编排者件 `.scratch/wisp/probes/303/orch/r1-orch-pair-summary.txt`（(a) `fb2fb802` rc=1 `nothing at all`／(b) 父 `f718e9b6` rc=0 `REPLIED,REPLIED,REPLIED`），本腿在母仓复跑 HEAD＝同红，见 `20-`。
第一手现场＝`git show fb2fb802 -- cmd/wisp/panel_host_windows.go`（件 `10-culprit-prod-diff.txt`）：那笔**只**把已有的 `w.Bind(panelDispatchBinding, …)` 原样搬进 `installPanelTransport`，并**新增** `w.Init(panelPostMessageForwardInit)`。⇒ green→red 的唯一 delta 是一枚**页面 JS 钩子**（见下）。

关键库事实（读码＋35r2 夹具逐字）：绑定出口链＝`webview.go:472` 桩体 `window.external.invoke(JSON.stringify({id,method:"wispDispatch",params:[…]}))` → `chromium.go:112` `window.external={invoke:s=>window.chrome.webview.postMessage(s)}` ⇒ **绑定门 `window.wispDispatch(...)` 的唯一致命副作用是它内部会调 `chrome.webview.postMessage(RPC帧)`**（`panel_transport_35r2_test.go:70/85` 逐字复刻）。本票受害三枚的页面 JS（`ac14AwaitJS`／`reportJSEnv`／AC#13 探针）**直接调 `window.wispDispatch(env)`**，不经 `postMessage` 门面。

---

## 甲 · 页面 JS 压根没跑
① 能把甲与乙/丙分开的读数：一枚**不经那道门**的页面可观测物——例如把 `document.title` 由 plain `Eval` 设值后用 Go 侧 `Eval` 取回（`nail2` 现在只经门回读，取不回门就分辨不了）。若门回读为零但门外 Eval 能读回 title ⇒ 排除甲；两者皆零 ⇒ 不能排除甲。
② 现量到了没有：**已用因果对间接否证，未跑门外 Eval-return 探针。**
   - 否证依据：green→red 的唯一差异是一枚纯 JS 钩子（`w.Init(forward)`）；撤掉这一行（父发）⇒ 复绿。若页面 JS 从不执行，加/不加一枚 JS 钩子**不可能**翻转结果 ⇒ 甲与"这一笔造成的红"矛盾。
   - 辅证（现量）：`20-targeted-before-red.txt` 三枚均 `window_opened=true`、`panel thread exited cleanly shows=1`、`cold -1.0 ms`；`TestAC13…` 打印 `AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]` ⇒ 入口文档已建、探针脚本已被投。
   - **缺的读数（具名）**：一枚门外 `Eval` 回读（`typeof window.chrome.webview` / 直接取 `document.title` 的同步 Eval 返回值）——**由非实现者 `303-v1` 补，或后续真窗夹具补**；本腿不为它动产码。甲已被因果对否证，这枚只是把"读码推"升级成"门外读数"。

## 乙 · JS 跑了但绑定名 `wispDispatch` 没接上
① 能把乙与甲/丙分开的读数：活页上 `typeof window.wispDispatch`。== "function" ⇒ 排除乙；== "undefined" ⇒ 乙。行为侧：乙时 ③ 钩子守卫 `typeof window.wispDispatch !== "function"` 为真 ⇒ 原始页信封**逐字交给 native 出口** ⇒ msgcb 见 `method:"panel.mode.request"`、id=0、未绑 ⇒ 死在 unbound 槽。
② 现量到了没有：**已否证（三条）。**
   - 代码：`fb2fb802` 的产码半边**没动** `Bind`，也没动绑定名常量 `panelDispatchBinding = "wispDispatch"`（`panel_host_windows.go:80`，逐字见 `10-culprit-prod-diff.txt`）。
   - 因果：父发同一枚 `Bind` 产出 `REPLIED,REPLIED,REPLIED` ⇒ 绑定名被正确注入且可调 ⇒ 门接得上。
   - HEAD 自洽：若 HEAD 上 `wispDispatch` 真未定义，③ 钩子的 `typeof !== "function"` 那一支会把一切路由到 native ⇒ **与父发同行为⇒应绿**；而 HEAD 现量**红**（`20-`）⇒ 反证 `wispDispatch` 在 HEAD 确为函数 ⇒ 乙不成立。
   - 辅证（现量，headless）：`TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent` ⇒ PASS，逐字 `door absent: the guard's typeof half routed the page's envelope to the native exit byte-for-byte (frames=1, no throw)`——这正是"若真是乙会看到的形"，与真窗 HEAD 的 `nothing at all`（门被进入但信封被双包）不同路。

## 丙 · 钩子接上了但门被再入帽挡住（本票活动形状）
① 能把丙与甲/乙分开的读数：门**被进入**（stub 的 `invoke→postMessage(RPC帧)` 触发 ③ 覆盖层）但门收到的**是被二次包裹的信封**（`method:"wispDispatch"` 而非 `panel.mode.request`），且**原生出口在直接 `wispDispatch` 入口下 0 次 / 递归爆栈**。
② 现量到了没有：**是，量到了（①旧形的再入是硬读数）。**
   - 现量（headless 行为尺，件 `14-`）：`TestLegacySubShapeOneHookDiesInAReentryLoop` ⇒ PASS 并逐字打印 `sub-shape ① under the behavioural yard: re-entered 9 levels (cap 8), native exit called 0 times, door fired 0 times, page error="RangeError: maximum call stack size exceeded - chrome.webview.postMessage re-entered 9 levels (cap 8) with the native exit called 0 times"`。＝丙的先例形（与 `286a7f30` 提交说明"再入 9 层、原生出口 0 次"一字对得上）。
   - HEAD 现状的机制（读码＋真窗读数合）：35-r2 换上的 ③ 钩子用布尔 `inside` 区分"页 postMessage"与"库发帧"，但 `inside` **只在覆盖层自己把页信封折进 `wispDispatch` 那一支被置真**。受害三枚**直接调 `window.wispDispatch(env)`**，从不经过覆盖层入口 ⇒ `inside` 恒假 ⇒ 桩内 `invoke→postMessage(RPC帧)` 撞上覆盖层时被当成"页又发了一次"折回 `wispDispatch` ⇒ 门收到双包 RPC 帧（`method:"wispDispatch"`，非在册四法）⇒ `ComposerDispatch.Handle` 于路由前丢弃 ⇒ `describeRequests` 空 ⇒ `no report … (what DID arrive at the door: nothing at all)` ＋ `firstRoundTrip` 报哨兵 `-1.0 ms`（`20-` 逐字 `cold -1.0 ms`）。
   - 为何 35-r2 的 headless 夹具当时"绿"：夹具的投递入口是 `window.chrome.webview.postMessage(envelope)`（`panel_transport_35r2_test.go:1384`，postMessage 源），`maxPostDepth=3 capTripped=false nativeExitCalls=1 doorRounds=1`（件 `14-` 现量）——**它从不模拟直接 `wispDispatch` 入口**，那正是真窗断口所在。⇒ 判语本体三句一字未动。

## 三形小结
甲：因果对已否证，门外 Eval-return 读数缺（归 `303-v1`／真窗夹具）。
乙：代码+因果+HEAD 自洽三重否证。
丙：**活动形状，有硬读数**（`re-entered 9 levels / native exit 0 / door fired 0`；HEAD 的直接 `wispDispatch` 入口 `inside` 未置真 ⇒ 门收到双包 RPC ⇒ `nothing at all`＋`-1.0` 哨兵）。
本票产品路径也断的判据不靠任一形假设，靠 `AC#1` 因果对＋此处真窗 `nothing at all`。
