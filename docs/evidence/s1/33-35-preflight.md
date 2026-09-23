# 票 33 / 票 35（Go 半边）— 开工前只读预检

**性质**：只读预检子代理。**本轮零写码、零测试运行、零 git 变更（除本文件）**；五问里每一处"存在/不存在"都是本代理自己 grep/读原文得出，不是引用他人转述。
**锚定 sha**：`dbbc822`（任务书给定）。交件时 HEAD 已到 `2c7fb54`；`git diff --name-only dbbc822 HEAD` = 4 枚文件（票 130/134 票面、`cmd/wisp/early_log_nail_130_windows_test.go`、本台账），**没有一枚是本报告引用的 `internal/panel/**`、`cmd/wisp/run.go`、`cmd/wisp/panel_assets.go`、`frontend/**`、`docs/specs/**`、`scripts/**`** ⇒ 下列行号在两处通用。
**工作树**：`git status --porcelain` 仅 1 枚未跟踪 `docs/evidence/s1/131-adversarial-acceptance.md`（他人产物，本代理未 add、未读作依据）。
**边界遵守**：`frontend/**` 全程只读；未写 `cmd/wisp/**`、`internal/panel/**`、`internal/winsec/**`、`internal/risk/**`；未碰冻结件（`docs/PLAN.md`、`docs/specs/*.md`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`）。
**本段无需时序/RSS 读数**（`A103`）⇒ 全文不报任何耗时。
**行号口径**：`file:line` 一律为 HEAD/锚定 sha 工作树的真实行号（写本文件前逐条 grep 复核，未沿用 114 表的行号未核）。

---

## ① 今天到底有什么

### ①.1 `internal/panel/` 现有文件清单与职责（10 枚非测试 + 7 枚测试）

| 文件 | 职责（一句话） | 关键锚点 |
|---|---|---|
| `internal/panel/doc.go` | 包边界声明：C27 单例宿主 + C17 桥 + embed 资源；**`:16` 自陈 "DEFERRED(host/bridge): implemented by ticket 33 (host), ticket 35 (bridge)"** | `doc.go:1-18` |
| `internal/panel/assets.go` | 只读视图 over 嵌入的 `dist/`：`BuiltinAssets` `:43`、`Resolve` `:75`、`Manifest` `:100`、`Check` `:137`、`Built` `:65`、`EntryFile` `:29` | `:3-9` 注释明写"The WebView2 host (ticket 33) calls Resolve" |
| `internal/panel/bridge.go` | 唯一的网页→原生封套读法：`ParseComposerRequest` `:77`、闭集白名单 `knownComposerMethod` `:97`（四常量 `:35-38`）、`RefusedEnvelopeForUser` `:108`、`NewRequestID` `:51` | 全文件**只有解析与拒绝，没有派发** |
| `internal/panel/composer.go` | 面板输入行的视图模型 + 请求形状：`Snapshot` `:44`、`NewSnapshot` `:63`、`ModeView` `:93`、`NewModeView` `:110`、`ModeUnknownView` `:120`、`ModeRequest` `:128`、`ModeRequest.Parse` `:139`、`WorkspaceView` `:161`、`WorkspaceRequest` `:180`、`NewComposerState` `:197`、`OutgoingMessage` `:213` | `:3-27` 规则：面板只"显示 + 发起请求"，**无 setter** |
| `internal/panel/composer_handlers.go` | 票 114 AC#2 的原生侧门：`ModeWriteHandler` `:86`、`HandleModeRequest` `:111`、`modeIsWidening` `:107`、门在 `:133`、写入在 `:141`；本地小接口 `ModeWriter` `:73` / `ModeConfirm` `:83`（**不 import `internal/perm`**，依赖面未动） | `:36-40` 自陈"WHAT THIS IS NOT: a path. Nothing calls this handler yet" |
| `internal/panel/attachments.go` | 附件入站的全部守卫：`AttachmentBroker` `:169`、`NewAttachmentBroker` `:179`、`Ingest` `:199`、`DecodeAttachmentPayload` `:381`、`AcceptedMIMETypes` `:153`、`MaxAttachmentBytes` `:47` | 无 handler、无宿主听众 |
| `internal/panel/workspace.go` | 工作区请求的原生腿：接口 `PathScope` `:37`、`AuditFunc` `:50`、`RequestWorkspaceSwitch` `:76`、`currentAfter` `:117`、`WorkspaceViewFromRoot` `:60` | 有实现，**没人调**（除测试） |
| `internal/panel/approval.go` | L2 卡的渲染模型：`ApprovalSubject` `:27`、`ApprovalCardView` `:39`、`NewApprovalCardView` `:64`、`CardViewFromDecision` `:72`（判定全部外包 `internal/risk`） | 票 35 的地界（推送载荷）尚未接 |
| 测试 7 枚 | `approval_test.go` / `assets_test.go` / `attachments_test.go` / `bridge_test.go` / `composer_test.go` / `composer_handlers_test.go` / `workspace_test.go` / `frontend_hygiene_test.go` / `tokens_fourway_test.go` | 其中 `composer_test.go`、`frontend_hygiene_test.go`、`bridge_test.go` 会**读 `frontend/` 源码做对账**（只读） |

### ①.2 谁 import 它（现核）

`grep -rln "wisp/internal/panel" --include=*.go .` 的**非测试**命中 = **2 枚**（比 114 表当时的 1 枚多一枚，因为票 114 AC#2 已落地）：

1. `cmd/wisp/panel_assets.go:23` import；用到的符号：`panel.NewApprovalCardView`（`:47`）、`panel.BuiltinAssets`（`:67`）、`panel.EntryFile`（`:103`、`:114`）。
2. `cmd/wisp/run.go:229` 字段 `modeWrites *panel.ModeWriteHandler`，装配点 `:378`。

⇒ `ParseComposerRequest` / `HandleModeRequest` / `RequestWorkspaceSwitch` / `NewAttachmentBroker` / `NewSnapshot` 的**生产调用者仍为 0**（本代理复跑 `grep -rn "HandleModeRequest\|modeWrites"`：除定义与 `cmd/wisp/run.go` 的字段/装配外全是测试；`HandleModeRequest` 的 10 枚调用者全在 `composer_handlers_test.go`）。
⇒ `cmd/wisp` 非测试代码里 `.Set(` 命中 = **0**（全 `cmd/` 只有 `cmd/llmrecord/main.go:94-96` 的 HTTP header，与档位无关）⇒ `perm.Store.Set`（`internal/perm/store.go:175`）生产仍 0。

### ①.3 `go:embed` 产物今天是怎么被"宿主"取的

链路（每一跳都有行号，**最后一跳断在"宿主不存在"**）：

```
npm run build -> frontend/dist/            （dist 被 gitignore，仅 dist/.gitkeep 跟踪）
  -> frontend/embed.go:19  //go:embed all:dist   -> frontend.Dist()  :23
  -> internal/panel/assets.go:44 fs.Sub(..., "dist")  -> panel.BuiltinAssets() :43
  -> 今天的唯一消费者 = cmd/wisp/panel_assets.go:67（`wisp panel-assets` 诊断子命令）
       -render  :74  assets.Resolve(path)
       -manifest :85  assets.Manifest()
       -check   :97  assets.Check()
  -> 契约里预期的消费者 = WebView2 的 AddWebResourceRequested 回调（票 33）：**今天零处**
```
`frontend/embed.go:5-7` 与 `internal/panel/assets.go:6-9` 都把自己写成"等票 33 来调"，即**产物侧齐备、取用侧为空**。
本机现存陈旧构建：`frontend/dist/assets/index-UL9kYvYl.js` + `index-CEH-Pz8P.css`（未跟踪）。

### ①.4 有没有任何一处已经在创建 WebView2 控件 —— **答：零**

三重否证（全部本代理实跑）：

1. **依赖根本没声明**：`go.mod` 的 require 块只有 sherpa-onnx / go-toml / x-crypto / x-sys / x-term / modernc-sqlite（`go.mod:6-14`），**没有 `github.com/jchv/go-webview2`**；`go.sum` 与 `deps.toml` 同样零命中（`deps.toml` 只登记 sherpa-onnx / onnxruntime 两组）。⇒ 票 33 的第一件事是**引一枚新的第三方依赖**，不是写代码。
2. **代码零命中**：`grep -rni "webview" --include=*.go .` 的**全部**命中都是注释/字符串/文档（`cmd/wisp/panel_assets.go:13,42`、`internal/panel/doc.go:1,7`、`internal/panel/assets.go:6`、`frontend/embed.go:5`、`internal/config/schema.go:531`、`internal/ball/doc.go:21` 等）；`grep -rn "CreateCoreWebView2\|WebMessageReceived\|AddWebResourceRequested\|AddHostObjectToScript\|PostWebMessage" --include=*.go .` 非测试命中 = **0**。
3. **连承载它的进程侧都还没有**：常驻腿 `cmd/wisp/resident_windows.go:24-85` 只做 Boot + 日志 sink + 空事件循环，`RunEventLoop()` 在 `:83`；全仓 `ball.New(` 的**生产**调用者只有调试壳 `cmd/balldebug/main.go:189` ⇒ 球都还没进常驻，面板宿主更没有。

**已经为票 33 准备好的插槽**（写手不必自己找）：
- 托盘/热键两个回调已经存在且今天没人填：`internal/ball/ball_windows.go:55` `Options.Events.OnTrayPanel`、左键与右键菜单项在 `:667`、`:672` 触发；面板热键 `OnPanelHotkey` 字段在 `:54`、触发在 `:660`（默认 `Ctrl+Alt+P`，`internal/ball/hotkey_windows.go:62`）。
- 现成的桩文案在 `cmd/balldebug/main.go:199`：`"tray: open panel (stub, ticket 33)"`。
- STA 线程：**`ui-sta` 归 `internal/ball` 私有**（`internal/ball/ball_windows.go:170` `Spawn("ui-sta", "ball", …)`；线程类型 `staThread` `internal/ball/sta_windows.go:24` 与投递方法 `(*staThread).PostTask` `:127` **均未导出**，只能经未导出字段 `Ball.sta` 可达 ⇒ 包外无法投递）。对 panel 唯一导出的钩子是 `internal/ball/d2d_windows.go:206 FactoriesForPanel()`（只给 D2D/DWrite 工厂指针）⇒ **今天没有任何"把 func 投递到 ui-sta"的口子**，票 33 要么经 `observe.Registry` 自己再起一枚 STA 线程，要么在 `internal/ball` 开一枚新导出（后者会碰球的地界，需先报备）。
- Job Object（WebView2 的 3-5 枚子进程必须入 Job）已可用：`internal/proc/jobscope_windows.go:71 OpenJobScope`、`:89 Assign(p)`；但 `Assign` 收的是 `*os.Process` ⇒ **对"由 COM 内部拉起的 msedgewebview2 子进程"如何入 Job，今天没有任何既有做法**（未验证）。

---

## ② 那一跳要经过哪几枚文件（施工图：网页点档位 → `rt.modeWrites`）

判读法：**已存在** = 有实现且有生产听众；**半存在** = 有实现但没人取用/缺一半；**完全不存在** = 缺文件或缺函数，给不出 `file:line`。
S1–S2 属 TS（冻结，见第③问），列出来只为让写手知道"上游今天已经把什么形状发出来了"。

| # | 必须存在的那一段 | 落点（file:line） | 状态 |
|---|---|---|---|
| S1 | 点击档位控件 → `requestModeSwitch(to)` | `frontend/src/components/composer.tsx:32`（import）、`:79`（调用） | **已存在（TS，冻结）** |
| S2 | 封套成形 + 送出：`sendRequest` 写 `{method, requestId, source:"panel-composer", to}` 后 `bridge.postMessage(JSON.stringify(…))` | `frontend/src/lib/panel.ts:236-238`（`requestModeSwitch`）、`:211-224`（`sendRequest`）、取桥 `:152-154` | **已存在（TS，冻结）** |
| S3 | WebView2 控件被创建（environment + controller + 导航到虚拟 host，`ui-sta` 上） | 无任何文件；依赖也未声明（`go.mod:6-14`） | **完全不存在（票 33）** |
| S4 | 页面资源服务：`AddWebResourceRequestedFilter` → `Assets.Resolve(path)` 喂嵌入字节 | 被调物已存在 `internal/panel/assets.go:75`；**注册方零处** | **完全不存在（票 33）**（Go 侧 API 齐、宿主侧空） |
| S5 | **事件接收**：`WebMessageReceived` → 从事件参数里取出那段 raw JSON 文本 | 全仓 `grep "WebMessageReceived\|PostWebMessage"` 非测试 = 0 | **完全不存在 ← 114 表 ①.1 指的断点正身** |
| S6 | Go 侧入口：`raw string → panel.ParseComposerRequest(raw)` | 被调物 `internal/panel/bridge.go:77`；**生产调用者 = 0**（6 枚调用点全在 `bridge_test.go`） | **半存在**：函数在、"接 raw 的那只手"不在 |
| S7 | 方法 → 处理器**派发**（按 `bridge.go:35-38` 四常量分流的 switch/表） | `knownComposerMethod` `bridge.go:97` 只答"认不认得"、**不路由**；全仓 `type Dispatcher\|func Dispatch\|func Route\|func Serve` = 0 命中 | **完全不存在（票 35 的 Go 半边正体）** |
| S8 | `panel.mode.request` 的原生处理器本体（含票 114 的 AC#2 门） | `internal/panel/composer_handlers.go:111`（门 `:133`、写 `:141`） | **已存在但生产听众 0**（`HandleModeRequest` 除测试外无人调；文件自陈 `:36-40`） |
| S9 | 载荷守卫 `ModeRequest.Parse`（档位名闭集） | `internal/panel/composer.go:139`，被 S8 在 `:122` 调 | **已存在** |
| S10 | 装配根把 handler 交出去 | 字段 `cmd/wisp/run.go:229`、装配 `:378`（`Modes: modeStore` / `Confirm:` 与 Store 同一条腿 / `Audit: rt.auditf`） | **半存在**：装配齐，但**没有任何取用点**——`rt.modeWrites` 在 `cmd/wisp` 非测试代码里只被写、不被读 |
| S11 | 真正的写腿 + L2 确认 + 审计 | `internal/perm/store.go:175`；腿 `cmd/wisp/run.go:354-356` → `rt.confirmModeSwitch` `:437` | **已存在**（生产 `.Set(` 仍 0，只能由 S8 抬成 ≥1） |
| S12 | **回执回灌 Go→页面**：拒答文案或新快照推给 webview | `RefusedEnvelopeForUser` `bridge.go:108` 生产 0；快照构造者 `NewSnapshot` `composer.go:63` / `NewComposerState` `:197` 生产 0；发送手段（`PostWebMessageAsString`）0 处 | **完全不存在（票 33 手段 + 票 35 泵）** |
| S13 | 页面消费快照的 listener（`PanelSnapshot` 的接收端） | TS 只有类型与 prop：`frontend/src/lib/panel.ts:129`、`frontend/src/App.tsx:45`；`grep onmessage\|addEventListener frontend/src` = **0** | **完全不存在，且属冻结面**（不是我们该补的） |
| S14 | 另两条输入的 handler：`panel.workspace.request` / `panel.attachment.add` | workspace 有原生腿无 handler：`internal/panel/workspace.go:76`；attachment 有 broker/解码无 handler：`attachments.go:199`、`:381` | **半存在**（守卫齐、处理器与派发缺） |

**结论形状**：S1–S2、S8–S9、S11 已经在；**缺的是中间那截 S3→S7 与回头的 S12**。
按票面地界：S3/S4/S5/S12 的手段 = **票 33**；S6/S7/S14 的语义 = **票 35 的 Go 半边**；S13 = 外部前端 agent。
最小通路（只做 Go、不碰 TS、不开真窗口）成立的**充要条件**是：S5 可以被打桩——即用一段与 `panel.ts:211-224` 产物**逐字同形**的字符串直接喂 S6，门与派发就都能绿；真窗口只影响 AC 的"真机那一格"，不影响 S6/S7 的可验性（分母证据见第③问）。

---

## ③ Go 半边与 TS 半边在哪里切开

判据来自 `docs/reports/pending-and-issues.md:3309-3312`（`A102③` "正解落在 `frontend/**` …一律不派"）与 `:3313-3318`（`A102③b` 明写"票 33 与票 35 的 Go 半边不冻，但一旦要往 `frontend/` 里写东西就必须停下上报"），复述于 `docs/reports/frontend-handoff.md:80-81`。**判据按目录，不按票号。**

### ③.1 属 Go（编队可以做）

| 落点 | 干什么 | 备注 |
|---|---|---|
| **新建** `internal/panel/host*.go`（票 33） | WebView2 环境/控制器、`AddWebResourceRequested` 回调调 `assets.go:75 Resolve`、显示/隐藏/销毁、焦点归还、创建失败→`panel.unavailable` | 必然带 `//go:build windows`（见 ③.3 的分母账） |
| **新建** `internal/panel/*_dispatch.go`（票 35 Go 半边） | S6（raw → `ParseComposerRequest`）+ S7（方法→处理器）+ S14（workspace/attachment 处理器） | 全平台中性， ubuntu 有分母 |
| `internal/panel/bridge.go` | 预计 0–10 行：白名单 `:97` 已是闭集，路由属于新文件 | 别再开第二条 pipe（`bridge.go:5-9` 是明写的规则） |
| `cmd/wisp/run.go` | 把 `rt.modeWrites`（`:229`/`:378`）交给宿主的派发；宿主生命周期挂进常驻腿 | 与票 130 写手同文件——**开工前必查 `git status`** |
| `cmd/wisp/resident_windows.go` / `main.go` / `cmd/balldebug/main.go:199` | 真起宿主、填 `OnTrayPanel`/`OnPanelHotkey` | 桩句已存在 |
| `internal/ball/ball_windows.go:54-55`、`sta_windows.go:24/127` | 若要复用 `ui-sta`：需新导出投递口 | **碰球的地界，先报备再动** |
| `go.mod` / `go.sum` / `deps.toml` | 引 `github.com/jchv/go-webview2`（MIT）并登记许可 | 今天**三处都还没有**（`go.mod:6-14`、`deps.toml` 只登 sherpa/onnxruntime）；⚠ 禁改清单里的 `tools/d22scan/**`、`allowlist.txt` 不在本行范围内，别顺手改它们 |
| 同包测试文件 | S6/S7/门 的用例 | `internal/panel` 在 ubuntu core scope ⇒ 有分母 |

### ③.2 属 TS（我们必须停手，命中即整段上报）

| 落点 | 缺什么 | 出处 |
|---|---|---|
| `frontend/src/lib/panel.ts` | **接收端**：没有任何 `onmessage`/`addEventListener`（实测 `grep -rn "onmessage\|addEventListener(" frontend/src` = 0），只有 `PanelSnapshot` 类型 `:129-137` | S13 |
| `frontend/src/App.tsx` | `snapshot?: PanelSnapshot` 只是 prop（`:45`，EMPTY `:34`），没人把推送接进来 | S13 |
| `frontend/src/components/composer.tsx` | 请求腿已齐（`:32`、`:79`）⇒ **这一枚今天不需要动** | S1 |
| `frontend/index.html` | 票 33 约束 `:29-30` 要 "CSP header/**meta** injection point"：`<meta>` 那一半在 `frontend/index.html`，**属冻结**；Go 侧只做响应头注入（`AddWebResourceRequested` 回调里加 header） | ⚠ 拆半做，`skipped=frontend(owner-delegated)` 记账 |
| `frontend/embed.go` | ⚠ 它是 **Go 文件但在 `frontend/**` 目录下** ⇒ 按 `A102③` 的"按目录"判据它在冻结面内。今天看**不需要改**（`//go:embed all:dist` `:19` 已够）；若写手认为要改（换 embed 根/加 CSP 注入点），**停手上报**，别按"这是 .go 所以能做"解释 | 边界易踩点 |
| `frontend/dist/**`、`frontend/package.json`、`frontend/scripts/**`、`frontend/fixtures/**` | 产物与构建 | 票 34/36–40 全部在冻结批（`A102③` `:3311`） |

### ③.3 "只做 Go 半边，Go 侧能不能独立成立并有分母"——**分两支，答案不同**

**能独立成立且有分母的那支（S6/S7/S8/S9/S11 + 回灌的构造侧）**：
- 请求封套的形状不必等 TS：TS 已经在发（`panel.ts:211-224` 的 `sendRequest`），Go 侧只需读同形状字符串；`internal/panel/bridge_test.go`、`composer_test.go` 已有把两侧字面量对账的门（`bridge_test.go:131-146` 遍历 4 枚 Go 常量要求 `panel.ts` 出现对应字面量，`composer_test.go` 的结构钉扫 `frontend/src`）⇒ **"只做 Go 侧"不破坏任何既有对账**。
- 分母：`internal/panel` 在 **ubuntu core scope**（`scripts/portable-tests.sh:179` 的 glob、`core_pin` 于 `:141`），CI 步名 `Portable package tests (core scope; the list and its guards live in scripts/portable-tests.sh)`（`ci.yml:266`，run 于 `:288`）⇒ 真执行，不是 `GOOS=linux go vet` 那种只编译。
- 〔引台账〕票 114 交件时该包读数 `RUN 152/PASS 92/FAIL 0/SKIP 0`（`pending-and-issues.md:3562`，`A110④`）——**本代理未复跑**，写手要用的话自己复现。

**不能独立成立的那支（S3/S4/S5 的 WebView2 宿主 + 真窗口 AC）**：
- ⚠ **关键分母缺陷**：`internal/panel` **不在 windows scope**——`win_pin`（`scripts/portable-tests.sh:151-160`）只有 ball/config/perm/plugin/proc/risk/secret/llmrecord，`--scope=windows` 的 glob（`:185-189`）也没有 `./internal/panel/`。
  ⇒ 带 `//go:build windows` 的宿主文件：**linux 腿（`go vet ./...` `ci.yml:122`）结构性排除它，windows 腿不测这个包** ⇒ **零用例分母**；唯一碰到它的是编译（`cgo build smoke`，`ci.yml:395`，run `:400`，会 link 整个 `wisp.exe`）。
  要给它分母，只有两条路：① 把它挂到 `cmd/wisp` 的 CLI 步（`ci.yml:402`/`:422`，scope `cli` = `./cmd/wisp/` 于 `portable-tests.sh:191-196`，windows-only）；② 把 `./internal/panel/` 加进 `win_pin` + windows scope（**这是改门禁形状，本代理不动、写手也不该顺手动，需编排者/owner 点头**）。
- 真机类 AC（票 33 AC 的 `:36-42`：show/hide 生命周期、冷≤1500ms/热≤200ms、焦点往返、netstat 无监听、runtime-missing fixture）**没有一枚能在 CI 上独立结**，且 `:38` 那格是时序读数（受 `A103` 本机 self-hosted runner 污染，须走安静性三连）。
- ⚠ **票 35 的 Blocked by 有一格是名义上的**：票面 `:7` 写 `Blocked by: 33-panel-host-c27, 34-frontend-scaffold`，票 34 状态仍 `ready-for-agent`，但 `frontend/` 产物实际已入库（票 77/92 交出 `src/App.tsx`、`src/lib/panel.ts`、`src/components/composer.tsx`、`dist/` 构建链）⇒ 票 35 的 Go 半边**不必等票 34**；反过来票 35 的 `parallel slots B`（`:7-8` 的"resync 模型 + 事件流 + **frontend client SDK**"）落在 ③.2，**这一整块要等 TS 半边，别开工**。

**一句话**：**派发应当只派"Go 语义半边 + 有 ubuntu 分母"那一段**；WebView2 宿主本体可以做、也能在本机编译与人工验，但**今天没有任何 CI 步会给它用例分母**，这一条要么先改 win scope（要授权），要么在票面上明写"该腿分母=本机人工，CI 无"。

---

## ④ 失败面：宿主起不来会怎样

### ④.1 今天的代码：无可失败之物（先说实话）

`wisp` 的三条腿里没有任何一处会去起宿主：常驻腿 `cmd/wisp/resident_windows.go:24-85` 只 Boot + 装日志 sink（`:57`）+ 跑空事件循环（`:83`）；`cmd/wisp/main.go:39` 里 WebView2 只出现在帮助文本；`cmd/balldebug/main.go:199` 的 tray-panel 回调是 `fmt.Println` 的桩。⇒ **"宿主起不来"这个失败面今天不存在**，所以票 33 落地时它是**新引入的面**，不是一条已存在的性质——这正是第十次那族病要防的形状（能力装上就以为通路通了）。

### ④.2 契约**原文**（不引转述）

- `docs/specs/SPEC-08-ui-ball-panel.md:152-153`：「**WebView2 Runtime 缺失 → 无面板模式**：D10 短结果 + 落文件仍可用 + L2 走原生降级卡；引导安装 Evergreen Runtime，**不得自动下载安装器**（D42#11）。」
- `docs/specs/SPEC-08-ui-ball-panel.md:138`（承载分层表 L2 行）：「…**L2 能力不得因面板缺失消失（C27）**」
- `docs/specs/SPEC-09-platform-windows.md:46`（第 11 行）：「WebView2 Runtime 缺失/过旧 | go-webview2 创建失败 | 无面板模式（C27 原生 L2 降级卡）+ 引导 Evergreen 安装；不自动下载安装器 | S5」
- 票 33 `:26-28`：「Runtime-missing detection: creation failure → `panel.unavailable` event + no-crash; L2 native fallback flag set (37 consumes); guide user to Evergreen install, never auto-download.」

⇒ **答任务问的那句：有，而且不止一条 AC。** 票 33 AC `:41`：「Runtime-missing fixture (rename/mask loader) → fallback event, app alive, native L2 flag on.」是**明写要造降级路径的验收格**；C27 的"L2 能力不得因面板缺失消失"是契约级要求（SPEC-08 `:138`）。

### ④.3 落地时的真实缺口（写手会撞上的三件事）

1. **`panel.unavailable` 事件今天没有承载物**：全仓 `grep -rn "panel.unavailable\|unavailable" --include=*.go internal/panel internal/agent/approval cmd/wisp` 无任何宿主事件符号（`internal/agent/approval/approval.go:95-103` 的 `unavailableText` 是**语音通道**的词，不是一回事）。要新建，且要选 sink——`internal/observe` 的日志 sink 只在 `wisp slo`/常驻腿装（`cmd/wisp/resident_windows.go:57`），事件若只走 `Logf` 则控制台腿看得见、常驻腿看不见（memory 里那条"安全告警默认只到 stderr"同族）。〔引台账〕
2. **"native L2 flag" 的消费方（票 37）还没做**：票 37 状态 `ready-for-agent`（`ls .scratch/wisp/issues/37-approval-ui-l2.md` 的 Status 行）⇒ 票 33 能**置位**这个 flag，但**没有任何人读它**，AC#5 的"native L2 flag on"只能测到"flag 被设"，测不到"降级卡真出现"。**这条要如实写成交件说明，不许因为框能勾就当结清**（本仓第九次就是这么抓的）。
3. **fail-closed 一侧今天已经偏严**：无宿主 ⇒ L2 卡片只能超时被拒（`cmd/wisp/run.go:118-123` 自陈"console run has no native channel … always resolves to a reject"、`:427-431` 注释、实现 `:437` `confirmModeSwitch` 走 `rt.gate`）。所以"降级路径"在产品感受上今天不是"退化"而是"什么都不许"：`perm.Store.Set`（`internal/perm/store.go:191-197`）对 `auto_approve` 在 `Confirm == nil` 时直接拒；票 114 的门 `composer_handlers.go:133` 对**变宽**在 `Confirm == nil` 时拒。⇒ 宿主失败面**不会**放宽任何权限，只会少可用性（这点对本票是好消息，写手不必担心降级路径变成后门）。
4. ⚠ 未验证：本机是否装了 WebView2 Runtime、`CreateCoreWebView2EnvironmentWithOptions` 在哪种失败下返回什么错误码——**本代理没有跑任何宿主代码**（不存在）。

---

## ⑤ 最小可交付切法：**一件**事

> 建议 = **S6 + S7 一起，并且当场给它一枚生产调用者**：新建 `internal/panel/composer_dispatch.go`（把一段 host-message 文本 → 白名单派发 → 票 114 的门 → 返回给用户的那句话），加同包用例，并在 `cmd/wisp/panel_assets.go` 开**一枚诊断入口**当它的生产听众。**除此之外一件都不做。**

**为什么是这一件**：
- 它是②表里唯一"上不下下"的那格：S8/S9/S11 全在等 S6/S7，而 S3/S4/S5 只在等这一件落地后才有意义；做了它，`ParseComposerRequest`（`bridge.go:77`）、`knownComposerMethod`（`:97`）、`RefusedEnvelopeForUser`（`:108`）、`HandleModeRequest`（`composer_handlers.go:111`）四条的**生产调用者同时从 0 抬成 ≥1**（114 表 ①.4 数的正是这四条）。
- 它**零 `frontend/`**：封套形状已由 `panel.ts:211-224` 定死，且 `bridge_test.go:131-146` 已经在拿 Go 常量去 `panel.ts` 里逐字找字面量 ⇒ 测试输入可以直接用与产物同形的字符串，不需要动 TS 一行。
- 它**双平台都有分母**：逻辑落在 `internal/panel`（`scripts/portable-tests.sh:179` 的 core glob、`core_pin:141`、CI 步 `ci.yml:266`/run `:288`），ubuntu 腿真执行；诊断入口落在 `cmd/wisp`，其用例在 windows CLI 步（`ci.yml:402`，run `:422`，scope `portable-tests.sh:191-196`）有分母。⇒ 不需要先动门禁（③.3 那两条路一条都不用走）。
- 它**不碰热文件**：`cmd/wisp/run.go` 一行不动（`rt.modeWrites` 那枚装配 `:378` 原样留着等宿主），避开与票 130 写手的文件级冲突；本代理读数时 `cmd/wisp/logsink.go` + `internal/observe/logging.go` 已被他人改动（`git status` 实测）。
- 它**天然带一个 fail-closed 的可证明形状**：诊断入口不接 C18 原生通道 ⇒ `Confirm` 传 nil，票 114 的门（`composer_handlers.go:133`）对"变宽"必拒、对"变严"必放，两端都能写成断言（`perm.New` 的 `Confirm` 允许 nil：`internal/perm/store.go:191-197` 那一支就是为此写的）。构造 `perm.Store` 只需 `config.NewManager(path, res)`（`internal/config/manager.go:78`）+ `perm.New`，**不必经 `assembleRuntime`** ⇒ 能做到不碰 `run.go`。〔静态核对签名，**未编译验证**〕

**它能不能独立红/绿：能。** 判据固定三发（写手自证，缺一不算结）：
1. 把派发里"route 到 mode handler"那一行摘掉 ⇒ "ModeWriter 被调用了一次"那枚断言**必须红**（判"没接上"看被调次数，不看日志字符串——114 表 ②.3 的同一条）；
2. 把 `knownComposerMethod` 的默认分支从"拒"改成"放行" ⇒ 白名单外方法名的那枚用例**必须红**（票 35 AC#1 `:42-43` 的 Go 侧那一半）；
3. 还原 ⇒ 复绿。⚠ **不许**用"日志里出现了拒绝字样"当断言替身。

**它故意不结的东西（要如实写进交件说明）**：面板**仍然**改不了档位——真机上没有那只手（S3/S5 未做），AC 里"真窗口/时序/焦点/无监听"六格一格都不动。这恰是 **`票 133`（`.scratch/wisp/issues/133-...md`，`R-121-1`）要造的尺能看见的形状**："只摘掉分发那一跳、两包用例一条不红"（票 133 `:15-27`）。⇒ 建议写手在交件里显式指回票 133，而不是自己顺手把那把尺做掉（那是另一张票的一个行为）。

**反例（为什么不是"先做 WebView2 宿主本体"）**：S3/S4/S5 的 6 枚 AC 全是真机/时序类，本包在 windows 腿零用例分母（③.3），且 `A102③` 已把"用眼睛看界面"那类证据冻在外部 agent 手里 ⇒ 一次 commit 验不完，红绿都立不住。

---

## ⑥ 未验证 / 伪授权登记 / `next=`

### ⑥.1 本代理查了但答不出的（一律"未验证"，不圆）

| # | 条目 | 为什么答不出 |
|---|---|---|
| U1 | `internal/panel` 里带 `//go:build windows` 的新文件在 CI 上有没有**任何**用例分母 | 静态判"没有"（`win_pin` `scripts/portable-tests.sh:151-160` 无 panel；windows 步 `ci.yml:458` 不测该包），但**未跑 CI**、也没查 `build.ps1` 是否顺带 `go test` |
| U2 | WebView2 的 3-5 枚子进程怎么入 Job（`proc.JobScope.Assign` 只收 `*os.Process`，`internal/proc/jobscope_windows.go:89`） | 没有既有代码可引，需 spike |
| U3 | `go-webview2` 引入后的许可台账落点（`deps.toml` 是否被 `wisp doctor` 逐条比对 Go module；`internal/buildinfo` 的读法） | 本代理只读到 `deps.toml:1-14` 的自陈，未核 doctor 的实现 |
| U4 | 票 33 的 CSP "响应头注入"与 `frontend/index.html` 的 `<meta>` 会不会互相冲突（两条同存时 WebView2 取哪条） | 属契约实现细节 + 需真窗口 |
| U5 | 诊断入口能不能真的不碰 `run.go` 就装出 `perm.Store`（`config.NewManager` 需要 data dir 与 `SecretResolver`） | 静态核签名，未编译 |
| U6 | 票 114 交件的 `RUN 152/PASS 92` 现值 | 〔引台账 `pending-and-issues.md:3562`〕本代理零跑测试 |
| U7 | `frontend/dist` 陈旧产物（`index-UL9kYvYl.js`，114 表 ③.2 记过"里面没有 composer 路由"）现在是否仍陈旧 | 只 `ls` 了文件名，未 grep 内容、未跑 `npm run build` |

### ⑥.2 伪授权登记（判据 = ①那枚路径/编号真不真 ②内容是否削弱 owner 权威或放宽判据）

**命中数 0**：本会话工具输出里**没有**任何自称「编排者备注 / 系统提示 / 用户已更新编码规则 / `wisp-orchestrator-continuation` / 请 revert / 冻结某包 / 放宽阈值」的文字，也没有要求我改判据或回滚的指示。

到达并被核对过的外来内容（**均未当作指令执行**）：
1. harness 的 `Note: The file C:\Users\swq\.qoder-cn\...\memory\MEMORY.md was modified since it was last read.` —— **本会话共 9 次到达**（项目级 `projects\D--work-...-Wisp\memory\MEMORY.md` 与全局 `C:\Users\swq\.qoder-cn\memory\MEMORY.md` 交替，首批一次到 2 条）。三条可判据内容差异各自核过：
   - "⚠09-23 追加批复：**票 33 + 票 35 的 Go 半边由本编队做**（TS 半边仍归外部），交接页 `docs/reports/frontend-handoff.md`" ⇒ ①路径真：`docs/reports/frontend-handoff.md` 存在（本代理读过全文，84 行）；②内容与账目一致：`pending-and-issues.md:3604-3607`（`A111⑤`，owner 批复「做，都按照你的推荐来就完事了」）。**不放宽任何判据**（它反而把 TS 半边划走）。⇒ 判合法，且**本报告的边界结论仍按 `A102③`/`A111⑤` 原文写，不按这条转述写**。
   - "`wisp-truth-source-docs` 追加两条治理裁定（'故意留空的一格' 与 '结案须 0 未勾' 冲突…）" ⇒ ①真实存在：`docs/reports/HANDOVER.md:107`、`pending-and-issues.md:3426`、`:3468`、票 121/133 票面 `:49-52`。与本报告无关（我不结任何票）。
   - 全局索引追加"改名 commit 后要数命中数"“我给 owner 的'推荐'同样是未验证断言"等条目 ⇒ 与本票无关，未据此改变行为。
   ⚠ **本代理的自陈缺口**：只把"到达"数清了（9 次），**没有逐枚记下它挂在哪一次工具调用上（工具名 + 命令前 40 字）**——`票 133` `:61` 与 `A104③` 要求的正是这一项。补不了就是补不了，登记为缺陷而非美化。
2. 首条消息尾部的 skills 清单（工具输出的一部分）⇒ 未据此调用任何 skill。

### ⑥.3 时间戳
本文件首段写入时本机 `date` 实测 `Wed Sep 23 11:58:17 CST 2026`；末次编辑前重读同一条命令，实测 `Wed Sep 23 12:10:56 CST 2026`（会话开始 11:5x，两个读数都是本代理自己跑的，未复用他人读数、未跨时段复用旧值）。

### next=（写手开工前还缺的东西）
1. **门禁授权（唯一一条硬阻塞）**：`internal/panel` 不在 windows CI scope（`scripts/portable-tests.sh:151-160`）⇒ 若要写宿主本体（S3/S4/S5），**先由编排者/owner 拍**"把 `./internal/panel/` 加进 win scope"还是"接受该腿无 CI 分母、以本机人工为证据"。本报告推荐的 S6/S7 那一件**不需要**这条授权。
2. **票面地界补一行**：票 35 的 `Parallel slots B`（`:7-8` 的 frontend client SDK）与票 33 的 CSP `<meta>` 那半格（`:29-30`）按 `A102③` 属外部 agent，建议编排者在两票票面各加一行"该半格 `skipped=frontend(owner-delegated)`"，否则写手会在同一枚框上撞第二次（`A102④` 的口径：规则要写在能挡住动作的那一处）。
3. **`panel.unavailable` 事件的 sink 归口**（④.3 第 1 条）：它进 `Logf`、进 `internal/observe`，还是进 C17 事件流（票 35 AC `:13-14`）？三者听众不同，写手不该自己选。
4. **`frontend/embed.go` 的归属确认**：它是 `.go` 但在 `frontend/**` 下 ⇒ 本报告按"按目录判"划入冻结（③.2）。若 owner 认为 Go 文件例外，需要一句明示；不给就按冻结办。
5. 开工前查 `git status --porcelain`：`cmd/wisp/run.go` / `panel_assets.go` 若仍在他人在飞清单内 ⇒ **文件级冲突就串行**（`A111⑤` 的原始理由）。
