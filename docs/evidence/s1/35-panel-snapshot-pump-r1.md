# 35 — 面板快照泵 r1：可达性探针（第 0 阶段）与落地

本件按派单第 0 阶段的要求先做只读盘点、落盘、再动生产码。
每一节裁完单独 commit，故本文件在收口前是分节的、带写节时刻。

---

## 0. 锚点与自量 sha

| 项 | 读数 | 怎么量的 |
|---|---|---|
| 进场锚点（本程开工时 HEAD） | `aeba6ff` | `git rev-parse --short HEAD`（本程第一次工具调用） |
| 写本节时 HEAD | `36734e4` | `git rev-parse --short HEAD`（共享树，期间别家落了 `921791a`/`4149971`/`65256e6`/`36734e4` 四枚，全在 `docs/**` 与心跳件里） |
| 锚点后本程射程内有无别家改动 | **无** | `git diff --stat aeba6ff..HEAD -- internal/panel cmd/wisp internal/agent/approval` 输出为空（现量于本节） |
| 工作树进门状态 | `design/**` 16 枚未提交删除 + 4 枚未跟踪（别家/owner 的活） | `git status --porcelain`，本程一格没碰 |
| 临时件目录 | `D:\tmp\wisp-35-pump-r1`（只建不删） | `base_gate.txt`、`base_roster_panel.txt`、`base_roster_cmd.txt`、`base_panel_names.txt`、`base_cmd_names.txt` |
| `frontend/src/lib/panel.ts` 工作树状态 | **干净**（`git status --porcelain -- frontend/src/lib/panel.ts` 空） | 双向尺（`composer_test.go:48`）读工作树，所以这一行必须留证，否则那发绿说明不了任何事（票 145 §4.4 末段） |

### 0.1 改前基线（本机、真实工作树，PATH 加了 `third_party/sherpa-onnx`）

`go test -count=1 ./internal/panel/ ./cmd/wisp/` → **rc=1**，`--- FAIL:` 共 **1** 枚：

- `--- FAIL: TestC21DesignTokensFourWayAgree`（`internal/panel`）——派单已预告：本机
  `design/assets/tokens.css` 被 owner 挪走所致。**保持红、不修不跳**，本程后面每一遍读数里它都应在。

名册基线（判红绿只认 `--- FAIL:`，`-v` 逐枚抽）：`internal/panel` 顶层 **53** 枚、
`cmd/wisp` 顶层 **60** 枚，除上面那一枚外全绿、顶层 0 枚 SKIP。
两枚名册文件落在临时目录（`base_panel_names.txt` / `base_cmd_names.txt`），§5 用差集。

---

## 1. 可达性探针：三条答案

判据是派单写死的那一句：**说不出"上一次真跑过它的 run id／调用点"，就当那条通道不存在。**

### 1.1 `C17 PanelBridge` 的**实际**方法集

**答案：`PanelBridge` 这个类型在本仓不存在，方法集 = 0 枚；今天树里唯一一枚"方法名册"是 composer 入站那 4 枚，它不是 C17。**

现量（`grep -rn --include=*.go "PanelBridge" .`，剔 `_test.go`）：

| 命中 | 位置 | 是什么 |
|---|---|---|
| 声明（`type PanelBridge` / `func .*PanelBridge`） | **无源**：0 处 | — |
| `internal/panel/doc.go:2` | 包注释 | "PanelBridge (C17) with method whitelist" —— 名字，不是实现 |
| `internal/panel/doc.go:9` | 包注释 | "PanelBridge method whitelist + … + push backpressure/merge (D38d)" |
| `internal/agent/approval/gate.go:595` | 注释 | "**a future PanelBridge**" —— 实现者自己承认它将来才有 |
| `internal/panel/doc.go:16` | `DEFERRED(host/bridge): implemented by ticket 33 (host), ticket 35 (bridge)` | 登记在案的推迟，与上面四行互相印证 |

`internal/panel/bridge.go` 这个文件名会误导：它**不是** bridge 的实现，而是 composer **入站封套的解析器**
（`:77 ParseComposerRequest`）。它的名册（`:97-103 knownComposerMethod`，逐枚现量）只有 4 枚：

| 常量 | 行 | 字面值 |
|---|---|---|
| `MethodModeRequest` | `bridge.go:35` | `panel.mode.request` |
| `MethodWorkspaceRequest` | `:36` | `panel.workspace.request` |
| `MethodAttachmentAdd` | `:37` | `panel.attachment.add` |
| `MethodMessageSend` | `:38` | `panel.message.send` |

而票 35 的 C17 冻结名册（`panel.resync` / `tasks.list` / `tasks.detail` / `history.query` /
`transcript.get` / `approval.current` / `approval.queue` / `approval.decide` / `config.get` /
`config.set` / `grants.list` / `grants.revoke` / `privacy.purge` / `privacy.export` /
`cost.summary` / `models.list` / `models.delete` / `diagnostics.export`）在 Go 侧**一枚都不存在**
（现量：`grep -rn --include=*.go` 剔 `_test.go`，`panel.resync`=0、`tasks.list`=0、`cost.summary`=0、
`history.query`=0、`grants.list`=0、`diagnostics.export`=0；`approval.decide` 的 3 处命中全在
`tools/d22scan/main.go`，那是 ban #6 的**扫描模式串**，不是路由）。

⇒ **要接的那根管子缺的第一层：`C17` 的宿主侧对象本身。它归票 35。**
⇒ 顺带：`bridge.go:32-33` 那枚"空承诺"（注释点名 `TestComposerMethodNamesMatchFrontend`，
"全仓 0 处声明"现量已核）确有其实、点名有虚——它承诺的那道门今天真实存在，只是名字叫
`composer_test.go:502 TestTheRendererHoldsExactlyOneDoorToTheHost` ＋ `:533 TestPlantedRendererDoorShapesGoRed`
（后者把同一个仪器对准一棵种坏的树做正控）。本程在 §2 落地集里**只改这一处注释的指向**，不改名册一字。

### 1.2 "Go -> 面板"这一跳今天到底存不存在

**答案：不存在。面板今天不由任何宿主承载，只由 `internal/panel` 的资产解析器出静态字节；全仓没有任何一条运行时通道能让 Go 主动把一份 JSON 推给已加载的页面。**

四条独立现量（任一条不足以判死，四条齐全）：

1. **WebView2 的宿主侧调用在 Go 里 0 处。**
   `grep -rn --include=*.go -i "PostWebMessage|WebMessageReceived|AddHostObjectToScript|SetVirtualHostNameToFolderMapping|CallDevToolsProtocol|NavigateToString" .`
   → **0 命中**。`grep -rn --include=*.go -i webview` 的非测试命中全部是**注释**与
   `internal/proc`（Job Object 里那 3-5 枚 `msedgewebview2` 子进程的配额），无一处是宿主实现。
2. **资产解析器在生产里唯一的调用者是那枚命令行诊断。**
   `panel.BuiltinAssets()` 非测试调用点：`cmd/wisp/panel_assets.go:88`；
   `(*Assets).Resolve` 非测试调用点：`cmd/wisp/panel_assets.go:95`（`-render`）与包内
   `assets.go:138/154`（`Check` 自查）。`assets.go:6-8` 自己写着
   "**The WebView2 host (ticket 33) calls Resolve**" —— 将来时。
3. **没有任何 HTTP/SSE/websocket 服务端。**
   `grep -rln --include=*.go "websocket|text/event-stream|http.Server|ListenAndServe" internal cmd`
   命中全部在 `internal/llm/*`（`anthropic`/`openaichat`/`openairesponses`/`golden`/`adaptertest`）
   与 `cmd/llmrecord`，**方向是出站到 LLM**，不是到页面；`internal/panel` 与 `cmd/wisp` **0 命中**。
   这与 `assets.go:8-9` 引的 D29 一致："there is no local HTTP server to fall back on"。
4. **常驻 GUI 路径不承载面板。** `cmd/wisp/resident_windows.go:24 runResident` 做的是
   `buildinfo.ResolveEnv -> proc.Boot -> installLogSink -> 空事件循环 -> D38(e) 收口`，
   全程不 import `internal/panel`（`grep -n "internal/panel" cmd/wisp/resident_windows.go` 0 命中）。

并且**这三条是本仓自己说出来的**，不是本件推断（逐字）：

- `internal/panel/composer_handlers.go:36-38`：
  "WHAT THIS IS NOT: a path. **Nothing calls this handler yet, because the WebView2
  "event -> ParseComposerRequest" hop does not exist in this tree** (status table ①.1;
  tickets 33/35 are still ready-for-agent)."
- `cmd/wisp/run.go:225-226`（`modeWrites` 字段注释）：
  "**Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does not exist
  in this tree** (tickets 33/35)".
- `cmd/wisp/panel_assets.go:53-56`：
  "this prints **exactly the JSON the WebView2 host pushes. Ticket 35 replaces the printing
  with a bridge push**; nothing else about the payload changes".

**判据回话**：说不出上一次真跑过快照推送的 run id／调用点 ⇒ 按派单写死的规则，**这条通道不存在**。
唯一能说出的"跑过"是 `wisp panel-assets -l2` 打印一枚卡片 JSON，那是 CLI 诊断、不是推送。

### 1.3 缺的是哪一层、归哪张票（未定义即停，不造替代）

| 缺的层 | 具体缺什么 | 归口 | 本程动不动 |
|---|---|---|---|
| L1 宿主 | WebView2 环境创建、C27 单例窗口、`Resolve` 接进 `AddWebResourceRequestedFilter` | 票 33（`33-panel-host-c27`，仍 ready-for-agent） | **不动** |
| L2 出站通道 | `PostWebMessageAsJson` 一类的宿主回调，即"Go 主动把字节交给页面" | 票 35 名下（`doc.go:16` 已登记 `DEFERRED(host/bridge)`） | **不动**，只把字节准备到它的入口 |
| L3 名册 | C17 冻结白名单的 Go 侧实现（含 `panel.resync`） | 票 35；且 `C17` 白名单**定稿**列在 AGENTS §2／`SPEC-12 §4.2` 的"未定义即停"清单里 | **不动** |
| L4 数据 | 快照的**构造与真来源** | 本程 | **动，这就是交付的那半** |

**本程明确不做的事**（派单禁形）：不写文件让前端轮询、不加本地 HTTP/SSE、不做 `panel.resync`、
不给 `Snapshot` 加第五枚键、不加 `currentView`、不进票 145 那 6 枚禁入字段。

⇒ 第 0 阶段的结论是：**通道不存在，按派单给的"那半"交付**——让 `Snapshot` 在生产里第一次被真的
构造出来、由一枚可测的出口函数交出字节，最后一公里（把字节交给页面）留在票 33＋票 35 名下。
本节先落盘，再动生产码。

---
