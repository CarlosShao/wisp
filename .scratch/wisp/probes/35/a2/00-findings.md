# 35-a2 只读普查：回执那一跳到哪儿为止

派单＝票 35 `:44` 格 / 票 248 AC#8 第二半。只量不裁。

- 锚 sha `98482adc4de13d4a29e20a57bc9e1ca895687cfa`（短 `98482adc`），分支 `dev`，时刻 `2026-10-09 11:44 +0800`。
- 起手在飞字节快照：`git status --porcelain -- internal cmd frontend` ＝ **空输出**（无别人未提交字节）。全程未改任何已跟踪文件。
- ⚠ **自报越界一发**：我跑过 `go env GOMODCACHE`（取模块缓存路径，好读库源码）。**没跑** `go build`／`go vet`／`go test`／`gofmt`。
- 所有读数默认尺＝`git … HEAD -- <下述 pathspec>`；带 pathspec 的一律写明，不写成"全仓"。
- 库源码在仓外（`D:\work\base\gopath\pkg\mod\github.com\jchv\go-webview2@v0.0.0-20260205173254-56598839c808`，`go.mod:19` 声明它 `// indirect`，无 `vendor/`）。下面标 `[库]` 的行号是那枚目录里的 `webview.go`，尺＝`grep -n`。

---

## 1　回程两条路各走到哪儿

**先报一条与派单前提不同的读数：这两条不是两条路，是同一条。** 同步返回的那枚串正是 `msgcb` 拿去 `Eval`  resolve 的那枚值。

ⓐ 同步返回，逐跳（尺 pathspec：`HEAD -- 'cmd/wisp/*.go' 'internal/panel/*.go'`）：

1. `cmd/wisp/panel_host_windows.go:801-804`（`installPanelTransport`）逐字：
   `func (m *PanelManager) installPanelTransport(w pageTransport, ctx context.Context) error {`
   `	if err := w.Bind(panelDispatchBinding, func(raw string) string {`
   `		reply, _ := m.dispatchRaw(ctx, raw)`
   `		return reply`
   ⚠ 票 35:389 引的是 `:713-716`，锚上已是 `:802-804`；`dispatchRaw` 定义在 `:814`（`:821 return m.disp.Handle(ctx, raw)`）。
2. `internal/panel/composer_dispatch.go:155-167` `Handle`：拒时 `return RefusedEnvelopeForUser(req, err), err`，成功时 `return receipt, nil`。
3. `:177-198` `dispatch` 的四枚非设置门逐字 `return "", d.Mode.HandleModeRequest(ctx, req)` 形 ⇒ **只有两枚设置门会填 `receipt`**（`:174-176` 注释自陈）。
4. 那句话本体：`internal/panel/config_handlers.go:431/435`（`tierSentence`）逐字 `return "这一项立即生效。"` / `return "这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。"`——复认，与 `248-w1` 一致。
5. **宿主层对这枚返回串的读者＝非 test 恰好 2 处**：
   - `cmd/wisp/panel_host_windows.go:803-804`（交回绑定，即下面 ⓑ）；
   - `cmd/wisp/panel_inbound.go:164` `reply, err := disp.Handle(ctx, raw)` → `:170`／`:181-182` `fmt.Fprintf(s.stdout, "wisp panel-inbound: 第 %d 行回执：%s\n", lineNo+1, reply)`。这枚是 CLI 文本腿（`wisp panel-inbound`），**不是页面**。
   `dispatchRaw` 非 test 调用者＝**1**（就是 `:803`）。尺＝`git grep -n "dispatchRaw" HEAD -- 'cmd/wisp/*.go' ':!*test*'`。

ⓑ `msgcb → Eval` 那条（本仓非 test 侧）：

- 本仓**没有**任何 `cmd/` 或 `internal/` 的非 test 代码调 `Eval`。尺＝`git grep -nE "\.Eval\(|msgcb|WebMessageReceived" HEAD -- '*.go' ':!.scratch/**' ':!*_test.go'` → 命中 2 行：`cmd/wisp/panel_host_windows.go:750`（注释，非调用）与 `[库]`无关的 `scripts/spike/webview2-latency/main.go:134` `w.Eval("typeof spikePing === 'function' && spikePing()")`。**`Eval` 的非 test 真调用者：产码 0 枚、spike 1 枚。**
- 调用者全在库：`[库] webview.go:139 func (w *webview) msgcb(msg string) {` → `:147 if res, err := w.callbinding(d); err != nil {` → `:151 } else if b, err := json.Marshal(res); err != nil {` → `:157 w.Eval("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined")`。`id` ＝ `:146 id := strconv.Itoa(d.ID)`，即**库自己的整数 seq**，来自绑定桩 `[库] webview.go:465 var seq = RPC.nextSeq++;` 与 `:472-476 window.external.invoke(JSON.stringify({id: seq, method: name, params: …}))`。
- ⇒ 回程帧上**不带**页面那枚 `pc-<seq>-<uuid>`；配对键是库的 seq，与本仓 `requestId` 无映射。
- 交付前提（投递泵）：`cmd/wisp/panel_resident_windows.go:299 w.Run()`；该文件 `:21-23` 逐字 "webview.Dispatch only appends a closure to the library's private queue and Run() is that queue's ONLY reader, so a host that pumps by hand can never deliver a Go -> page reply"。库侧 `[库] webview.go:351 func (w *webview) Run()`、`:443 func (w *webview) Dispatch(f func())`。
- 页面侧终点（对象层读，尺＝`git show HEAD:frontend/src/lib/panel.ts`）：`:141 postMessage(message: string): void;`、`:211 function sendRequest(method: string, payload: Record<string, unknown>): void`、`:218-225` 只 `bridge.postMessage(JSON.stringify({method, requestId: nextRequestId(), source: "panel-composer", ...payload}))`，返回值无处接。宿主转发钩逐字 `cmd/wisp/panel_host_windows.go:790 try { return window.%[1]s(message); } finally { inside = false; }`——**钩子是把 promise 原样 return 给页面的**，接不接是页面侧的事。

> 本节尺与读数：`dispatchRaw` 非 test 调用者 **1**；宿主层 reply 串非 test 读者 **2**（其一为 stdout 文本腿）；本仓产码 `Eval` 非 test 调用者 **0**；页面 `postMessage` 签名 **void**（`panel.ts:141`）；Promise 到页面调用点 **有**（`panel_host_windows.go:790`），页面取用 **0 处**。

---

## 2　`NewRequestID()` 现状

- 调用形状尺（逐字）＝`git grep -nE '[A-Za-z0-9_]\.NewRequestID\(|[^a-zA-Z_.]NewRequestID\(' HEAD -- '*.go' ':!**_test.go'` → **1 命中**：`HEAD:internal/panel/bridge.go:80:func NewRequestID() string {`＝**声明本身**，不是调用。⇒ **非 test 调用者 0 枚。**
- 词面尺（逐字）＝`git grep -n "NewRequestID" HEAD -- '*.go' ':!**_test.go'` → **2 命中**（`bridge.go:78` 文档注释、`bridge.go:80` 声明）。全 `*.go` 计数字尺 `git grep -cn "NewRequestID"` → `bridge.go:2` + `bridge_test.go:1`（测试唯一调用者＝`internal/panel/bridge_test.go:48 return envelopeJSON(t, m, NewRequestID(), ComposerRequestSource, extra)`）。
- 两把尺在"非 test 调用者＝0"上**一致**；差别具名：词面尺会把声明行与文档注释行也算成命中（调用形状尺的第 2 分支 `[^a-zA-Z_.]NewRequestID\(` 恰好也吃下 `func NewRequestID(` 这行，所以两把在这里都不是纯调用尺——**枚数请以"命中行全为声明/注释"这条复核读，别把 1 当"有一个调用者"**）。
- 页面 `pc-<seq>-<uuid>`（`frontend/src/lib/panel.ts:205-209`，逐字 `return \`pc-${requestSeq.toString()}-${uuid ?? "no-randomuuid"}\`;`）在 Go 侧的读者——按语义变体量：
  - **有读者，但只读不配对。**（⚠ 我这一条的原始 grep 输出被 `head -30` 截断过，下面列的是**已复认的子集**不是全集；全集枚数见本节末读数行，子集里每一处我都在截断前看到了原文。） `internal/panel/bridge.go:93 RequestID string \`json:"requestId"\``（被 `ParseComposerRequest` 解析）、`:136`/`:139-140` 缺 id 即按"无法与审计/卡片对齐"拒、`:158-160 id := r.RequestID` 进日志。消费点：`internal/panel/composer_dispatch.go:221/233-234/247`、`composer_handlers.go:113/118/122/135/163`、`config_handlers.go:252/271-272/282-283/295/300/310`。**全是拒因文案与审计行的相关键**，无一处把它登记成"等回执的槽"。
  - `cmd/wisp/*.go` 非 test 里 `requestId|request_id` 命中 **0**。尺＝`git grep -ciE 'requestId|request_id' HEAD -- 'cmd/wisp/*.go' ':!*test*'`（rc=1，无匹配）。
  - 字面 `pc-` 前缀在 Go 侧的读者＝**test 专属**（`cmd/wisp/panel_transport_35r1_test.go:54-55` 把它当信封常量、`panel_inbound_guards_35r3_test.go:80-81`），产码 0 枚。
  - `bridge.go:73-76` 逐字自陈这枚 id 的性质："It is a log correlation key, not a nonce and not authority"；`:85 return "rid-unavailable"` 是取随机失败支，无时钟回退。

> 本节尺＝上面三条逐字串；读数：非 test 调用者 **0**（两把尺一致，注释/声明行差异已具名）、`RequestID` 在 `internal/panel` 非 test 命中 **27 行／4 枚文件**（尺＝`git grep -n "RequestID" HEAD -- 'internal/panel/*.go' ':(exclude)internal/panel/*_test.go'`；`bridge.go` 6／`composer_dispatch.go` 3／`composer_handlers.go` 5／`config_handlers.go` 13，含声明与注释行，**无一处登记成"等回执的槽"**）、在 `cmd/wisp` 非 test **0 处**。

---

## 3　请求态容器的候选落点（只给落点与代价区间，不选形、不写码）

- `ComposerDispatch` 结构体**全部字段＝6 枚**，逐字（`internal/panel/composer_dispatch.go:123-140`）：
  `type ComposerDispatch struct {` / `	Mode ModeRequestHandler` / `	Workspace WorkspaceRequestHandler` / `	Attachment AttachmentRequestHandler` / `	Message MessageRequestHandler` / `	Config ConfigRequestHandler` / `	Audit AuditFunc` / `}`
  （锚上结构体闭括号在 `:140`；票 35:389 引的 `:121-139` 差 2 行，是注释/声明边界漂移，字段名与枚数一致。）
- 生产装配点＝**1 枚**：`cmd/wisp/panel_inbound.go:272 return &panel.ComposerDispatch{`。尺＝`git grep -n "ComposerDispatch{" HEAD -- '*.go' ':!.scratch/**'` → 命中 37 行，其中非 test **1 枚**（其余 36 枚在 `_test.go`）。
- 最少要动的文件（**落点清单，非推荐形**）：
  1. `internal/panel/composer_dispatch.go`——结构体 + `Handle`/`dispatch` 两个签名（`:155`、`:177`）。
  2. `internal/panel/bridge.go`——若配对键要复用 `NewRequestID`（今天 0 调用者）或 `ComposerRequest.RequestID`。
  3. `cmd/wisp/panel_host_windows.go`——`:741-744 type pageTransport interface { Bind…; Init… }`（今天**没有**出向原语）与 `:802-804` 绑定闭包。
  4. `cmd/wisp/panel_inbound.go:272`——唯一装配点，字段一加就要动。
  5. `frontend/src/lib/panel.ts:141`/`:211`——页面侧那一跳（**不在 Go 射程**，`frontend/**` 我只走对象层读）。
- 会跟着变色的仪器（登记，不评判）：
  - `pageTransport` 的 fake **4 枚**：`cmd/wisp/panel_pageover_33r10_windows_test.go:87`、`panel_reshow_255r1_windows_test.go:136`、`panel_transport_35r1_test.go:96`、`panel_transport_35r2_test.go:1229`（尺＝`git grep -n "func (.*) Bind(name string, f interface{}) error" HEAD -- 'cmd/wisp/*.go'`）。接口一宽，这 4 枚全要补。
  - 门名册尺：`cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go`（`:34` 断言 `installPanelTransport` 在生产源里**恰好 1 枚**，`:87` 逐字写死 `wispDispatch` 那一格的语义）。
  - 回程形状尺：`panel_transport_35r1/r2_test.go`、`panel_inbound_guards_35r3_test.go`、`panel_transport_live_35v2_windows_test.go`。
  - **`_rpc`／resolve 在 `cmd/wisp` 的 test 里 0 命中**（尺＝`git grep -nE '_rpc|resolve' HEAD -- 'cmd/wisp/*.go'`，命中的 8 行全是英文注释里的 "resolve" 一词，无一枚建模库的 resolve）⇒ 今天**没有任何仪器**在假页面上验过"同步返回真的落到了页面"。
- 代价区间（**粗估，未实测；只按文件枚数给band**）：
  - 支（i）"Go 侧一行不动，只补页面 await"＝Go 产码 **0 枚文件**，页面 1 枚文件（`panel.ts`）。
  - 支（ii）"在 `ComposerDispatch` 加请求态容器 + `Handle` 侧配对"＝**4 枚产码文件**（落点 1/2/3/4）＋ **≥4 枚 fake** ＋ 名册尺，行band **约 200–500 行**。
  - 支（iii）"另立异步帧/推送（走 `Eval`/事件名）"＝上面全部 **+  `pageTransport` 加出向原语**（4 枚 fake 全动）+ C17 名册/`bridge.go` 白名单可能要新事件名，行band **约 300–600 行**（与票 35:389 乙读法给的 3–4 枚/300–600 同量级，但我这枚的前置不是 `:47`，而是 `pageTransport` 今天没有出向方法）。
  - 判据句"什么时候生效"这一句**已在 Go 侧**（`config_handlers.go:431/435`），三条支都不需要新写这句话本身——这句是量出来的，不是估的。

> 本节尺＝上面 pathspec 逐字串；读数：字段 **6**、非 test 装配点 **1**、`pageTransport` fake **4**、`_rpc` 建模 **0**、`pageTransport` 方法数 **2**（`Bind`+`Init`，无出向原语）。

---

## 4　这条链上已经有人立过票吗（只登记，不裁归属）

尺＝`git grep -niE '回执|回程|并发配对|outbound|requestId' HEAD -- '.scratch/wisp/issues/*.md'` ＋ `git grep -n "H10" HEAD -- '.scratch/wisp/issues/*.md'` ＋ `ls` 文件名族词筛。命中：

- **票 33**（`.scratch/wisp/issues/33-panel-host-c27.md`）＝唯一 `H10` 出现处（5 格）：
  - `:239` 逐字："**H10** 回执回灌页面（`Handle` 那句现在只到 stdout）、"
  - `:238-240` 把"界面点一下后端真收到"拆成 `H2`/`H3`/`H10`/`H1` 四跳未落，并注明 `**H1** 前端发送腿（frontend/**，另一会话）`。
  - `:122` 逐字："H3（`WebMessageReceived` 把页面那段 raw 取进 Go）、H10（回执回灌页面）**三枚未落**"；`:111` "别把这一格当成'入向已接线'"；`:150` "真宿主那片（H2/H3/H10 的手段）今天**在 windows scope 零用例分母**"。
- **票 35**（`35-panel-bridge-c17.md`）：
  - `:389`（就是被派的 `:44` 格）逐字："**一块没写**（不是'写了零调用者'）：回执今天只能同步返回 `cmd/wisp/panel_host_windows.go:713-716`，`ComposerDispatch` 结构体（`composer_dispatch.go:121-139`）**没有任何请求态容器**；页面侧配对靠库的 `window._rpc[seq]`（⛔ 不是本仓产码）"；并载两条代价支与"今天 0 枚"仪器读数。
  - `:294` 逐字："去掉 `return` 全绿是**等价形**（页面两处把 postMessage 当语句用、答复回程走 `msgcb→Eval`）⇒ 只有页面哪天 `await` 那次 post 的返回值才变成新面"。
  - `:235`/`:221` 是 `requestId` 守卫那两格（`:221` `TestInboundRequestIDGuardRefusesMissingIDOnPageEdge`）。
- **票 248**（`248-…settings-route…md`）：`:56` AC#8 判据句本体；`:59-60` AC#10 乙形要求"在设置页那几项的**回执文案**里逐字写明……"。
- **票 255**（`255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md`）＝票名自带 "the panel host never receives config"。
- **票 92**（`92-panel-composer-mode-attachments-workspace-done.md`）`:275` 逐字："`ParseComposerRequest`（`source:"panel-composer"` + `requestId` 的 fail-closed 校验）"。
- 其余命中（`12:111`、`22:13`、`113:60-61`、`116:23/86/107`、`119:53/77`、`221:189/263`、`261`）经逐行核对＝**同词不同义**（`outbound` 指 HTTP 头、`配对` 指 sha↔数字/路径 `%` 折叠、`回执` 指编排者翻勾记录），不构成这条链的立票。

> 本节尺＝上面两条逐字串；读数：这条链上**已有同族票 5 枚**（33 / 35 / 248 / 255 / 92），`H10` 唯一归属票＝**33**。⛔ 归属判断留给编排者。

---

## 5　我没答的／量不到的

1. **"resolve 真到页面"我没有实跑证据**：需要真 WebView2 宿主 + 消息泵在跑，本腿禁跑 `go build`/`go test`。我只能给出**代码形状**（`[库] webview.go:147-157`）＋本仓两处**别人记录的现量**（`panel_host_windows.go:35-39`、`panel_resident_windows.go:21-25` 引 33-p1 §A R25/R26 "arriving in ~107ms with Run()"）。这两处是**读来的**，不是我复认的。
2. **票 35:389 那两条代价支的行数band 我无法验证**（不写码、不试编译）。第 3 节我的 band 是按"文件枚数 × 常识行数"粗估，标了未实测，别当凭据。
3. **库会不会在别的调用形状上丢掉返回串**：`callbinding`（`[库] webview.go:147`）我只读了调用点，没读它函数体对 `func(string) string` 单返回形状的 marshal 细节；`json.Marshal(res)` 成功支我复认（`:151/:157`），失败支（`:153 reject`）在什么输入下触发＝量不到。
4. **页面两处 `postMessage` 调用点中 `requestApprovalResolution`（`panel.ts:179`）那一处的返回串语义**我没追：它不在 AC#8 这条腿上，超出射程。
5. **`msgcb` 的 `id` 与页面 `requestId` 之间有没有我漏掉的映射表**：尺是词面 + 语义变体两把，若映射写在非 Go 侧（如前端某处自己维护 seq↔requestId 表），本腿的 pathspec 看不到——我在 `frontend/**` 只读了 `panel.ts` 一个文件。
6. **未答**：这活该落票 35 还是另立（派单明令不裁）。
