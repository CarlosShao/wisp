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

⇒ `ParseComposerRequest` / `HandleModeRequest` / `RequestWorkspaceSwitch` / `NewAttachmentBroker` / `NewSnapshot` 的**生产调用者仍为 0**（本代理复跑 `grep -rn "HandleModeRequest\|modeWrites"`：除定义与 `cmd/wisp/run.go` 的字段/装配外全是测试；`HandleModeRequest` 的 12 枚调用者全在 `composer_handlers_test.go`）。
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
