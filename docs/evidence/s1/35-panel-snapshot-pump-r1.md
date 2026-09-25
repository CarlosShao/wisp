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

## 2. 出口到底通到哪里（续程，回答派单点名的那三问）

本节由**续程**写。上一程中断前最后一句是"日志管道把任何单条字符串截到 512 字符，所以那个包
乘不了这条管道；正在把出口重构为'有界记录 ＋ 保留的包'"。下面三问逐条给 `file:line`，
每条都带现量命令；**"能交字节"与"有人接"在本节是两栏，不合并**。

### 2.1 今天这份快照字节实际上流到哪里去了

链条（生产路径，每一跳当轮现量）：

| 跳 | 位置 | 是什么 |
|---|---|---|
| 1a 装配 | `cmd/wisp/run.go:428` | `rt.ui.publish = rt.publishPanelSnapshot` |
| 1b 装配 | `cmd/wisp/run.go:550` | `consoleSink{out: rt.stdout, stream: rt.stream, publish: rt.publishPanelSnapshot}` |
| 2 触发点 | `cmd/wisp/run.go:809` / `:823` / `:761` | `consoleApprovalUI.Prompt` 末、`consoleApprovalUI.Update`（仅 Dismissed/Started）、`consoleSink.Publish`（仅 `changed` 的那几类事件） |
| 3 入口 | `cmd/wisp/panel_pump.go:237` | `snap, _, err := rt.pump.Publish()` |
| 4 构造 | `internal/panel/pump.go:205-206` | `Snapshot()` ＋ `json.Marshal` |
| 5 交付 | `internal/panel/pump.go:215` | `p.src.Out(snap, data)` |
| 6 `src.Out` 是谁 | `cmd/wisp/run.go:426` | `Out: rt.bookPanelSnapshot` |
| 7 落点 A | `cmd/wisp/panel_pump.go:152-154` | 全量字节留在**进程内存**（`rt.lastSnap` / `rt.lastSnapBytes` / `snapSeen`） |
| 8 落点 B | `cmd/wisp/panel_pump.go:155` | 一条**有界摘要**进本 run 的持久 ledger |

**"有人接"这一栏（这是本节真正的答复）**：

| 出口 | 非测试调用者 | 现量 |
|---|---|---|
| `(*SnapshotPump).Marshal()`（`pump.go:190`，被叫作"那枚出口函数"） | **0** | 全仓 `Marshal(` 的非测试命中里无一枚指向它；唯一调用者是 `internal/panel/pump_test.go:107` |
| `lastPanelSnapshot()`（`panel_pump.go:215`，全量字节唯一的读法） | **0** | 三处命中全在 `cmd/wisp/panel_pump_test.go:173 / :307 / :360` |
| `snapshotCount()`（`panel_pump.go:247`） | **0** | 命中在 `panel_pump_test.go:417 / :429 / :430` |
| ledger 里那条 `panel: SNAPSHOT` 记录 | **0 枚读者** | `grep -rn "panel: SNAPSHOT" --include=*.go .` → 3 命中：`panel_pump.go:182`、`:189` 是**生成方**，`panel_pump_test.go:56` 是**同一条尺的读方** |

⇒ 结论逐字：**构造与交付这两跳今天真的通了**（`Publish` 由生产自己的三处触发点驱动，见 §5 的
`TestThePumpIsDrivenNotJustAssembled` 读数），**但没有任何第三方读者**：全量字节今天只活在
`rt` 这个 struct 里、只有测试向它伸手；持久件只有一枚有界摘要，而那枚摘要今天的读者也只有
测试自己。**"面板看得到这份快照"这件事今天仍然不成立**，与 §1.2 的四条现量一致（下面 2.1.1 重测）。

#### 2.1.1 §1.2 那四条在中断后的树上仍然成立（续程重测，不照抄）

| §1.2 的断言 | 现在的读数 | 命令 |
|---|---|---|
| WebView2 宿主侧 API 在 Go 里 0 处 | **0 非测试命中**（仍成立） | `grep -rn --include=*.go -iE "PostWebMessage\|WebMessageReceived\|AddHostObjectToScript\|SetVirtualHostNameToFolderMapping\|CallDevToolsProtocol\|NavigateToString" . \| grep -vc _test.go` |
| 无 HTTP/SSE/ws 服务端 | **0 枚真代码**（仍成立） | 上表那条 `grep -rln ... internal/panel cmd/wisp` 返回 1 枚**文件**、命中的是 `internal/panel/pump.go:17` 那行**注释**（新泵自己写的"没有 HTTP/SSE"）。⚠ 这枚 1 是仪器数到自己头上，不是破口；按文件筛会误判，按行看才是 0 |
| `runResident` 不承载面板 | **仍 0**（本批没碰 `cmd/wisp/resident_windows.go`） | `grep -c "internal/panel" cmd/wisp/resident_windows.go` → 0 |
| `PanelBridge` 类型 0 枚声明 | **前后都是 0** | `grep -rn --include=*.go -E "type PanelBridge\|PanelBridge interface" . \| grep -v _test \| wc -l` → 0（`aeba6ff` 副本同为 0） |

### 2.2 那枚 512 字符截断是哪一枚既有设施的行为（现量，不信注释）

**是 `internal/observe` 的冻结脱敏规则 4**，不是日志层新加的东西、也不是本票造的：

| 环节 | 位置 |
|---|---|
| 常量 | `internal/observe/redact.go:36-37` `MaxLoggedString = 512`（头注释 `redact.go:26` 就是那句"Long argument strings are truncated (bounded log lines)"） |
| 施加点 | `internal/observe/redact.go:141` `return truncate(s, MaxLoggedString)`，在 `Redactor.String` 里；`redact.go:144-145` `Message()` 直接委派给 `String()` |
| 截断实现 | `internal/observe/redact.go:193-199` `truncate` —— 按 **rune** 计数，并**追加** `...(truncated, %d chars total)` 后缀 |
| 谁逼每一条记录过这一关 | `internal/observe/logging.go:132-138 redactHandler.Handle`，其中 `:133` `h.red.Message(rec.Message)`、`:135` `h.red.Attr(a.Key, a.Value)`；handler 由 `logging.go:104 Handler()` 交出 |
| 生产怎么装上它 | `cmd/wisp/logsink.go:149 observe.InitLog(observe.LogConfig{Dir: dir, Level: logSinkLevel})` → `cmd/wisp/logsink.go:122-127 (*logSink).logger()` 返回 `slog.New(s.pipeline.Handler())` |
| 泵走的就是这条路 | `cmd/wisp/panel_pump.go:155 rt.spec.sink.logger().Info(...)` —— 与 `cmd/wisp/run.go:538`（`auditf`）同一个方法 |

**现量**（仓外副本 `D:\tmp\wisp-35-pump-r2\probe\truncprobe\main.go`，跑的是 `observe.InitLog`
**真管道**、再把它写出的 JSONL 文件读回来，不是调一次函数自说自话）：

| 输入 | 送入 rune 数 | **落到文件里的 rune 数** | 落盘串尾部 |
|---|---|---|---|
| 900×`A` | 900 | **543** | `AAA...(truncated, 900 chars total)` |
| 640-rune 的 JSON 包（含一枚 L2 卡的形状） | 640 | **543** | `...\":1...(truncated, 640 chars total)` |
| 摘要行的形状 | 180 | **180**（没被截） | ` bytes=997 sha256=2d711642b726b044` |

`observe.MaxLoggedString = 512` 由探针直接打印。**两处口径要在引用时带上**：

1. 512 是**保留部分**的上限，落盘的那枚串实为 `512 + 31 = 543` rune（后缀是 `truncate` 自己加的）。
   `cmd/wisp/panel_pump.go:136` 与 `panel_pump_test.go:26` 写的"caps/bounds … at 512 characters"
   对"保留多少"是准的、对"落盘多长"少算一枚后缀。**本程不改这两行注释**（它不是缺陷，是简写），
   只在这里补一句口径。
2. `truncate` 数的是 **rune**，`panelSnapshotSummary` 的钳位 `len(line)`（`panel_pump.go:186`）
   数的是 **byte**。今天这一路 `at/depth/mode/ws/bytes/sha256` 全是 ASCII、`GeneratedAt` 是 RFC3339
   ASCII，所以两者不会分出不同结果；但**它不是同一把尺**。这是本程**没测**的形状（§7 第 5 条）。

**这一批真打出来的包有多大**：`go test -v ./cmd/wisp/ -run TestRunBooksWithASnapshotOfItsLiveQueue`
的 `panel_pump_test.go:216` 那行打印 `packet bytes (831)`，同一条 ledger 记录打印 `bytes:831`。
⇒ **831 > 512 ⇒ 那个包今天乘不进这条管道**，上一程这一句**成立**。
⚠ 但它的**数值**不成立：`cmd/wisp/panel_pump.go:138-139` 与 `cmd/wisp/panel_pump_test.go:27`
两处注释写着"最小包 534 字节 / 带卡的那枚 997"，**实测这一枚带卡的包是 831**（两处注释与
`37a4705` 的 commit message 里那句"最小包 831 字节"互相矛盾）。详见 §6。

### 2.3 它的"有界记录 ＋ 保留的包"落成了没有、落成什么样

**落了，形状是两件都在同一枚函数里**：

| 半件 | 落成什么 | 位置 |
|---|---|---|
| 有界记录 | `panelSnapshotSummary(snap, data)`，字段 `at / depth / pending(前 3 枚 id，超了就 +N) / mode / ws / results=总数/done 数 / bytes / sha256` | `cmd/wisp/panel_pump.go:163-194` |
| 自钳位 | `summaryClamp = 440`，超钳时**先丢 ids**、保留计数与指纹 | `cmd/wisp/panel_pump.go:186-192`、`:196-199` |
| 指纹 | `sha256Short` 取 16 枚 hex | `cmd/wisp/panel_pump.go:203-206` |
| 保留的包 | `rt.lastSnap / rt.lastSnapBytes / rt.snapSeen`，读法 `lastPanelSnapshot()`（互斥锁保护） | `cmd/wisp/panel_pump.go:152-154`、`:215-222`，字段在 `cmd/wisp/run.go:237-245` |
| 尺 | `assertPacketMatchesLedger145` 把"落盘那条记录说的字节数/指纹/深度/档"对上"保留的那枚包" | `cmd/wisp/panel_pump_test.go:84-100`（被 `:211`、`:311` 两处调用） |

**没有为它新造通道**（这一条按派单的"不许自选发明替代管子"逐形现量）：本批 5 枚非测试件
（`internal/panel/pump.go`、`internal/panel/bridge.go`、`internal/agent/approval/pending_read.go`、
`cmd/wisp/panel_pump.go`、`cmd/wisp/run.go`）**新增行里 `http` / `websocket` / `ListenAndServe` /
`event-stream` / `os.WriteFile` / `ioutil.WriteFile` / `postMessage` / `webview` 命中 0 枚**。
正控：同一条尺打在 `94071ff`（真引入服务端形状的提交）上是 **7 枚命中** ⇒ 这把尺是活的。
⇒ **不需要停手上报的那一支没有发生**：这半件落在"构造 ＋ 可测出口函数"的授权面内，
把字节交给页面仍然留在票 33／票 35 的出站那一跳，本程与本批都没造替代管。

一处**要报名排者的取舍**（不是缺陷，但它改变"出口"这个词的重量）：

- `bookPanelSnapshot` 无论如何都 `return nil`（`cmd/wisp/panel_pump.go:156`）。
  ⇒ `Publish` 的"出口拒收"分支（`internal/panel/pump.go:215-218`）在生产里**打不到**，
  `publishPanelSnapshot` 那行 stderr 兜底（`cmd/wisp/panel_pump.go:238-241`）也就**不会因 ledger 而响**。
  今天唯一会响的分支是 `Out == nil`（`pump.go:211-214`），钉它的是
  `internal/panel/pump_test.go:219 TestPublishWithoutAnExitReportsInsteadOfPretending`。
  读法：**"快照没能送达"这件事今天只有一种失败形状会被说出口**，ledger 写失败、写错、
  被别家截断这三种形状都还沉默。这一格是票 35 出站那一跳的事，本程不修（修它就是替
  那根管子做决定）。

---

## 3. 落地集与没落的那半（续程核过的版本）

### 3.1 落地集（三枚 commit，逐枚现量）

| commit | 文件 | 行 | 这一枚做了什么（续程读过、不是转述 commit message） |
|---|---|---|---|
| `5821e24` | `internal/panel/pump.go` | 341（新） | `NativeVerdict` / `PumpSources` / `SnapshotPump`（`Snapshot` / `Marshal` / `Publish` / `Publishes`）/ `StreamLog`（有界合并） |
| | `internal/panel/pump_test.go` | 309（新） | 6 枚用例（名册见 §5） |
| | `internal/agent/approval/pending_read.go` | 60（新） | `(*Queue).LiveApprovals()` 只读枚举 `statePending` 项，返回 `tools.Decision` 原样拷贝 |
| | `internal/panel/bridge.go` | 8/1 | **只改注释**：把点名不存在的 `TestComposerMethodNamesMatchFrontend`… 的旧指向改指真存在的两枚门（自证见 §3.3）；本程复算：这枚文件的**行为**未变（`ParseComposerRequest` 与其 4 枚方法常量逐字未动） |
| `37a4705` | `cmd/wisp/panel_pump.go` | 252（新） | 四个读数器接活对象 ＋ 有界摘要出口 `bookPanelSnapshot` |
| | `cmd/wisp/panel_pump_test.go` | 430（新） | 3 枚用例 ＋ 两把 helper 尺 |
| | `cmd/wisp/run.go` | 94/6 | `agentRuntime` 加 `pump/stream/lastSnap*/snapSeen`；装配根 `:420-428`；`consoleSink` 加 `stream/publish`；`consoleApprovalUI` 加 `publish` |
| `6e348f1` | `cmd/wisp/panel_pump_test.go` | 16/5 | 把 `shortenL2Window145` 改名 `nonDefaultConfig145`，**新增**一枚断言 `:295-298`（档必须是夹具设的非默认档） |

> ⚠ `6e348f1` 那一枚**只落了断言、没落夹具写入**——它的注释（`panel_pump_test.go:120-125`）
> 声称给了两个值，代码里 `:133` 只往 `[risk]` 追加了 `confirm_timeout_sec = 2` 一个，
> `permission_mode` **一个字都没写**。后果与隔离实验见 §5.2 与 §4 的 `M3_no_wfix`。
> **本程没有、也不会把它"顺手改回去"**；派单写死了这条理由不许推翻，本节只报它落了一半。

### 3.2 冻结面：本批有没有碰到不该碰的（逐条现量）

| 面 | 命令 | 读数 |
|---|---|---|
| `internal/panel/composer.go`（`Snapshot` 四键的家） | `git diff --stat aeba6ff..HEAD -- internal/panel/composer.go` | **空** |
| `frontend/src/lib/panel.ts` | 同上换路径 | **空**（且 `git status --porcelain -- frontend/src/lib/panel.ts` 空 ⇒ 双向尺读的工作树是干净的，那发绿说明得了事） |
| `internal/risk/**`（含 `rules_gateway.go`） | `git diff --stat aeba6ff..HEAD -- internal/risk/` | **空** |
| `thresholds.go` / golden / `allowlist.txt` | `git diff --name-only aeba6ff..HEAD \| grep -iE "thresholds\|golden\|allowlist"` | **0 命中** |
| `internal/panel/l2_grant_boundary_test.go`（别族依据件） | `git diff --stat aeba6ff..HEAD -- internal/panel/l2_grant_boundary_test.go` | **空** |
| `docs/PLAN.md` / `docs/specs/**` | `git diff --name-only aeba6ff..HEAD` 全 13 枚路径里无这两类 | **0 命中** |
| `design/**`（owner 16 枚未提交删除） | `git status --porcelain` 逐枚原样 | **一格没碰、也没算进任何零命中** |
| `DEFERRED(` 标记总数 | 前后各数一遍 | **29 → 29**（本批 0 枚新增；缺的那根管子早登记在 `internal/panel/doc.go:16 DEFERRED(host/bridge): implemented by ticket 33 (host), ticket 35 (bridge)`，该件未被本批改动）⇒ `SPEC-12 §5` 登记表无需新增条目 |

### 3.3 `bridge.go` 那处注释的新指向（本程自证，不留同型 bug）

派单点名的这一条要单独证：**"注释预先引用尚未产出的读数"是本仓假绿前身**，而这次改的正是那类 bug。

```
$ grep -n "func TestTheRendererHoldsExactlyOneDoorToTheHost\|func TestPlantedRendererDoorShapesGoRed" \
      internal/panel/composer_test.go
502:func TestTheRendererHoldsExactlyOneDoorToTheHost(t *testing.T) {
533:func TestPlantedRendererDoorShapesGoRed(t *testing.T) {
$ sed -n '502p;533p' internal/panel/composer_test.go      # 行号与声明同行，两枚都对
$ grep -rn "TestComposerMethodNamesMatchFrontend" --include=*.go . | wc -l
0                                                          # 旧指向仍 0 处（改前也是 0，见 §1.1）
```

⇒ 两枚真存在、真在那两行、与 `bridge.go:31-36` 的**新**指向一字不差。这一处**成立**。

### 3.4 没落的那半（仍然没落，且本程没替它做决定）

| 缺的层 | §1.3 的归口 | 现在的状态 |
|---|---|---|
| L1 宿主（WebView2 环境、C27 单例、`Resolve` 接进请求过滤） | 票 33 | **未落**，本批 0 处宿主代码（§2.1.1 现量 WebView2 API 仍 0 命中） |
| L2 出站通道（Go 主动把字节交给页面） | 票 35 名下 | **未落**。今天 `src.Out` 接的是 ledger，全量字节只有测试在读（§2.1） |
| L3 C17 名册（含 `panel.resync`） | 票 35，且白名单**定稿**在 AGENTS §2「未定义即停」清单里 | **未落**：`PanelBridge` 类型声明前后各 0 枚（§2.1.1）；`bridge.go` 仍只有入站那 4 枚 composer 方法名 |
| L4 快照的构造与真来源 | 本票段 | **落了**，承重自证见 §4 |

⇒ 派单写死的那一条本程**照做**：把字节送到界面是票 33／票 35 出站那一跳，本程没有、也不该
自选发明替代管子（`§2.3` 那条零命中的现量就是这条的凭证）。

---
