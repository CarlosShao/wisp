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

## 4. 承重自证：逐字段答"哪一枚用例断言了它来自真来源"（票 145 AC#6 那一问）

全部在**仓外副本** `D:\tmp\wisp-35-pump-r2\mut`（`git -c core.autocrlf=false -c core.eol=lf archive HEAD`
抽的，含 `design/assets/`）里打；仓库工作树**一字未动**（§5.5 给它单独现量）。
驱动件 `D:\tmp\wisp-35-pump-r2\mutate.py`，每发一枚 `.log` ＋ 一枚 `.raw.txt`。

### 4.1 四枚字段各一发（派单点名的 pending／results／composer／generatedAt）

| # | 字段 | 被换掉的真来源读法 | 结果 |
|---|---|---|---|
| **M1** | `pending` | `cmd/wisp/panel_pump.go:70` `Reason: it.Decision.Reason` → 常量句 | **红** `TestRunBooksWithASnapshotOfItsLiveQueue`，红句 `panel_pump_test.go:240: reason: console="R2: 目标路径在授权目录之外: C:\\Windows\\win.ini" packet="hardcoded by M1, not read from the queue", want the same native reason on both surfaces` |
| **M1b** | `pending`（不带 §4.5 那枚夹具补全，再打一遍） | 同上 | **仍红、同一句** ⇒ `pending` 的承重**不依赖** `6e348f1` 那半件 |
| **M2** | `results` | `cmd/wisp/run.go:425` `Results: rt.stream.Chunks` → 常量 chunk | **红**同一枚用例，红句三处：`:321 results correlationId = "M2-fixed", want the task id the run printed "b4520d56-…"`、`:328 chunk text "M2 fixed text", want the provider's reply the console also showed`、`:331 the booked text "M2 fixed text" is not the text this run streamed to the console`。（首发因替换串非法 Go 报 `[setup failed]`，**那不算读数**，修正后重跑才是这一发。） |
| **M3** | `composer.mode` | `cmd/wisp/run.go:423` `Mode: rt.modes.PermissionMode` → 常量 `risk.ModeAskEveryStep` | **红**同一枚用例，红句两枚：`:293 composer.mode.current = "ask_every_step", want the boot read "ask_high_risk"` ＋ `:296 …, want the non-default档 this config set` |
| **M4** | `generatedAt`（注入钟那一支） | `internal/panel/pump.go:179-181` `if p.src.Now != nil { now = p.src.Now }` → 恒常量钟（读数器被摘掉） | **红两枚**：`TestThePumpBuildsThePacketFromWhatTheHostHolds` 响 `pump_test.go:102: generatedAt = "2026-01-01T00:00:00Z", want the injected clock, RFC3339 UTC`；`TestPublishHandsTheBytesToTheAttachedExit` 响 `pump_test.go:292: …, want the injected clock's stamp` |

> 抽取红句的窗口是被测名之前 14 行，因此会把相邻的 `t.Logf` 一起带进来。本表**只收真正的断言行**；
> 已核对的一条噪声：`l2_grant_boundary_test.go:2332` 那句 `ZERO INSTRUMENT COVERAGE…` 在**干净基线**
> （`post_panel.txt`）里同样出现 1 次 ⇒ 它是常驻自述日志、不是变异红。

### 4.2 另外两枚派单点名的形状

| # | 摘掉的是什么 | 结果 |
|---|---|---|
| **M3b** | `composer.workspace` 的真来源：`run.go:424 Workspace: rt.workspaceView`（→ `panel_pump.go:86 panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())`）→ 常量 `panel.UnsetWorkspaceView()` | **红** `TestSnapshotWorkspaceSectionReportsTheNarrowing`，红句 `panel_pump_test.go:369: packet workspace = {Set:false … Reason:未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定}, want the narrowing this run applied` |
| **M5** | `NativeVerdict.CardView` 那条"原样搬运已决断 verdict"的路径（`internal/panel/pump.go:72-85`，走 `CardViewFromDecision(subject, risk.Decision{Level: v.Level, RulesHit: v.RulesHit, …})`）→ **改成在面板侧重算**：`NewApprovalCardView(risk.NewRiskAssessor(), subject)` | **红六处，本票段最重的一发**：`internal/panel` 的 `TestThePumpBuildsThePacketFromWhatTheHostHolds` 响 `pump_test.go:71 level = "L0", want the assessed level rendered by risk`、`:74 rulesHit = [], want the rules verbatim and in the assessor's order`、`:77 reason = "无规则命中（L0 直接执行）" known=true, want the native reason carried as-is`、`:80 sessionOverrideBlocked was lost between the verdict and the card (R4's flag)`；`cmd/wisp` 的 `TestRunBooksWithASnapshotOfItsLiveQueue` 响 `panel_pump_test.go:234 level: packet="L0" console="L2"`、`:240`（reason 对不上）、`:262 the card showed rule R2 on the console but the packet carries []`。 ⇒ **数值后果写出来**：泵一旦在面板侧自判，一枚真 **L2**（R2 命中）的卡在快照里会报成 **L0 · 无规则 · "直接执行"**。这就是"搬运不重算"那条设计的承重证明 |

### 4.3 "Snapshot 四键与 `panel.ts` 未动、双向尺仍绿"——本程自己复算了

| 复算项 | 读数 |
|---|---|
| 四键未动 | `git diff --stat aeba6ff..HEAD -- internal/panel/composer.go` → **空**；`frontend/src/lib/panel.ts` → **空**（且工作树该文件 `git status --porcelain` 空 ⇒ 尺读的是干净树，那发绿说明得了事） |
| 双向尺仍绿 | `TestComposerContractTypesMatchFrontend` 在改后的 `internal/panel` 里 `--- PASS`，`t.Logf` 自报 `Snapshot <-> PanelSnapshot: 4 JSON keys reconciled`（4 枚键，逐字对得上 `composer.go:44-49`） |
| **这把尺不是恒绿的**（M6 反向一发） | 把 `composer.go:48` 的 `json:"generatedAt"` 改名 `generatedAtRenamed` ⇒ **派单点名的 `composer_test.go:74` 响**：`Go Snapshot emits [generatedAtRenamed] that interface PanelSnapshot does not declare`；另一向同时响（`approval_test.go:129` 同一句 ＋ `:132 interface PanelSnapshot reads [generatedAt] that Go Snapshot never sends - those fields render as undefined`）；泵自己的两枚键集钉 `pump_test.go:119` / `:271` 一起响（`TestThePumpBuildsThePacketFromWhatTheHostHolds`、`TestPublishHandsTheBytesToTheAttachedExit`）。⇒ 派单引的 `:73-78` **两向都是活的** |
| ⚠ 口径：`cmd/wisp` 那 3 枚**不响** | M6 下 `cmd/wisp` rc=0 全绿 —— 它们断言的是解析后的 Go struct 字段、不是 JSON 键名。⇒ "键名被改"只有 `internal/panel` 那把尺管得住，**别拿 `cmd/wisp` 的绿当第二枚凭证** |

### 4.4 答不出的字段＝装饰品（本程实测出来的，逐枚给依据）

| 字段 | 裁定 | 依据 |
|---|---|---|
| `generatedAt` 的**生产那一支** | **装饰品（本程新量出来的洞）** | `assembleRuntime`（`cmd/wisp/run.go:421-427`）**从不设 `PumpSources.Now`** ⇒ 生产走的永远是 `internal/panel/pump.go:178` 的 `now := time.Now`。把**这一支**换成常量（M4b：`now := func() time.Time { return time.Date(2020,1,1,…) }`）⇒ **两包 59／63 枚全绿、0 红 0 跳**。M4 之所以响，响的是"注入的钟被读没有"，不是"生产那枚钟是活的"。⇒ `generatedAt` 在**库**一层承重、在 **`wisp run` 的实际装配上不被任何用例管**：`panel_pump_test.go:283` 只断言非空，`assertPacketMatchesLedger145`（`:84-100`）比的是 bytes／sha256／depth／mode，**没比 `at`** |
| `pending[].sessionOverrideBlocked` 的**生产那一支** | **装饰品（M7 量出来），且今天造不出承重** | 见下一行 |
| ↑ **M7** 单发：`cmd/wisp/panel_pump.go:71` `SessionOverrideBlocked: it.Decision.SessionOverrideBlocked` → 常量 `false` | **两包 59／63 枚全绿、0 红 0 跳**（副本内，夹具已按 §4.6 补全，故红因不可能是别处） | **因不只是"少一枚断言"**：`cmd/wisp/run.go:439-440` 逐字写着 "**R4 stays dormant**: probing every result for sensitive sources needs the C25 detector's own wiring (ticket 25)" ⇒ 今天 `wisp run` 这条路上**没有任何生产者能把这一位置真**，真来源与常量 `false` 在本路径上**可观测地等价** ⇒ 任何断言都恒真、写了也是假绿。唯一的现量仪器 `-taint-source` 只挂在 `wisp panel-assets` 那枚 CLI 诊断上（`cmd/wisp/panel_assets.go:44`，自建 Facts），走不到队列。**⇒ 这一格归口票 25（C25 探测器接线），本程不造 mock 顶它**（"用 mock 代替真的"是禁形）。库那一层另有 `pump_test.go:79-81` 钉"搬运不丢位"，M5 里它以 `:80 sessionOverrideBlocked was lost between the verdict and the card (R4's flag)` 响过 ⇒ 承重的是**搬运**、不是**生产来源** |
| `pending[].decidedBy` | **本来就是常量，且不是本批造的** | 值来自 `internal/panel/approval.go:88 DecidedBy: "native"`——纯函数里写死。两枚用例（`pump_test.go:82`、`panel_pump_test.go:221`）断言的是"它等于 native"，不是"它来自某处"。这是票 77 冻结的形状、本批只搬运它 ⇒ 记为**既有事实**，不计入本票段的装饰品 |
| `pending[].callChain` | **诚实的空** | 生产恒 nil（`panel_pump.go` 不填、`pump.go:53-59` 写明理由），两枚用例断言 `len(…) == 0`（`pump_test.go:85`、`panel_pump_test.go:224`）。⇒ 它**没有来源可换**、也就没有"换成常量"这一发可打；缺的是生产者本身（`tools.Decision` 无 chain 字段），那是另一格 |
| `composer.maxAttachmentBytes` | 承重、但来源是**同一枚常量的两次引用** | 生产**不设** `PumpSources.AttachmentMax`（`run.go:421-427`）⇒ `pump.go:167-170` 回落 `MaxAttachmentBytes`；`panel_pump_test.go:278-279` 断言的也是 `panel.MaxAttachmentBytes`。值对（与 `attachments.go` 的 broker 默认一致），但**它没有经过任何"读"的动作** ⇒ 若哪天 broker 的默认与常量分家，这枚包与这条用例都不会响。列 §7 没测 |

### 4.5 `6e348f1` 那半件到底承不承重（派单写死"这条理由不许顺手改回去"，本程独立复核）

**它的理由成立，且现在有了实测支撑**：

| 实验 | 配置夹具状态 | 变异 | 哪枚断言响 |
|---|---|---|---|
| **M3** | 副本里补全（`[risk]` 真写上 `permission_mode = "ask_high_risk"`） | mode 读数器 → 常量 `ask_every_step` | `:293` **与** `:296` 两枚都响 |
| **M3_no_wfix** | **仓库现状**（`[risk]` 里没有 `permission_mode` 这个键） | 同一发变异 | **只有 `:296` 响**，`:293` **不响** |

⇒ `:293` 那枚"对 boot 的 `perm: MODE-READ` 审计行"的交叉核对，在**默认档**下对这一发变异是
**结构性打不响的**（boot 真读出来的就是 `ask_every_step`，硬编码的泵与它对得上眼）。
**今天唯一抓得住"把档硬编码成默认值"的一枚，正是 `6e348f1` 新增的 `:295-298`。**
⇒ 派单那句理由（"默认档正是一枚硬编码 `ask_every_step` 的泵会打出来的值；档位在夹具里是非默认的，
包里的 mode 才必须真有人读它才对"）**复算为真**，本程不改它、不删它。
⚠ 同一次复算也量出它**只落了一半**（夹具写入没做）——后果与那枚新红见 §5.2。

### 4.6 隔离实验 E0（把 §5.2 那枚新红的因与果钉死）

副本里**只**补那一行夹具写入、不改任何生产码：
`internal/panel` rc=0 / 59 枚 / **0 红** ／ `cmd/wisp` rc=0 / 63 枚 / **0 红**。
⇒ 那枚红**不是泵的缺陷**（泵的 mode 读数器一配上非默认档就读对了：M1/M2/M3b 的 booked record 里
`mode:ask_high_risk`），**是夹具写入缺了一行**；而补上它之后，M1／M2／M3／M3b／M5 每一发仍然各自
把同一枚用例打红 ⇒ **补那一行不会把任何一根钉子磨钝**（这一句是给编者的，本程不动手）。

---

## 5. 三门 ＋ 改前改后四数 ＋ 名册差集

### 5.1 取数形状（两支正路里选了哪一支、为什么）

**改前基线在仓外副本取**：

```
mkdir -p /d/tmp/wisp-35-pump-r2/base
git -c core.autocrlf=false -c core.eol=lf archive aeba6ff | tar -x -C /d/tmp/wisp-35-pump-r2/base
```

副本字节＝锚点字节，用 `git hash-object --no-filters` 逐枚比 `git rev-parse aeba6ff:<path>`，
**全 1174 枚都核**（不是抽样）：

| 读数 | 值 |
|---|---|
| 树内文件总数（`git ls-tree -r aeba6ff`） | 1174 |
| `hash-object --no-filters` 与锚点 blob 不符的枚数 | **6** |
| 那 6 枚是什么 | `scripts/build.ps1`、`scripts/dev/ball-cycle.ps1`、`scripts/fetch-deps.ps1`、`scripts/sign-models.ps1`、`scripts/slo-check.ps1`、`scripts/spike/run.ps1` —— **全部 `*.ps1`** |
| 差的确实只有行尾吗 | 抽两枚核：`git show aeba6ff:<f> \| tr -d '\r' \| md5sum` 与 `tr -d '\r' <副本/<f> \| md5sum` **相同**（`build.ps1 → 136c48e2…`、`slo-check.ps1 → 2635bcc3…`） |
| 为什么正好 6 枚 | `.gitattributes:2` 写着 `*.ps1 text eol=crlf` —— 那是 **checkout 属性**，单给 `core.autocrlf=false` 压不住，必须连 `-c core.eol=lf` 一起给；而**不带 `--no-filters` 就根本看不出这件事**。派单预告的"那 6 枚 `eol=crlf` 差是预期"复算为真 |

副本里没有 `third_party/sherpa-onnx/`（DLL 不入库）⇒ 从工作树拷三枚进去（`onnxruntime.dll`、
`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`）。

**用的正路＝②（把仓里那几枚 DLL 放进 `PATH` 再 `go test`）**，并且**逐包分开跑**：

```
cd <tree> && export PATH="<tree>/third_party/sherpa-onnx:$PATH" \
  && go test -count=1 -v ./internal/panel/      # 再单独跑 ./cmd/wisp/
```

- **为什么不用①（`scripts/wisp-cli-tests.sh`）**：它末行是 `bash "$portable" --scope=cli`，
  而 `scripts/portable-tests.sh:195` 里 `cli` 那档的范围是 `scope=(./cmd/wisp/)` **一枚包**
  ⇒ 拿不到 `internal/panel/` 那一半，派单要的是**两包四数**。
- ⚠ 代价本程写明：②不经 `portable-tests.sh` 的 GUARD A/B，也不经 `tools/d22scan/runtests.sh`
  那三道（顶层 SKIP 即红、PASS 与 FAIL 同时为 0 即红）。本程用**手工等价尺**补：判红只认
  `^--- FAIL:`、SKIP 与 RUN 各自单列（见 5.2／5.4），四数之外点名册差集。
- ⚠ **一枚取数坑（本程踩过、写下来给别程）**：第一次把两包合成一发
  `go test -v ./internal/panel/ ./cmd/wisp/`，两包并行输出**交错**、`ok\t<pkg>` 尾行落在测试行
  **之后**，按"上一个 ok 行"归包会量出 `internal/panel = 113 枚 / cmd/wisp = 0 枚`这种荒谬读数。
  ⇒ **逐包跑才是可归因的形状**；聚合跑的总枚数（113＝53＋60）虽然对，分包数全错。

### 5.2 四数（改前基线取自锚点副本，改后取自本机工作树）

| 树／包 | `=== RUN`（行首精确） | 顶层 PASS | **顶层 FAIL** | 顶层 SKIP | 子测试 P/F/S | 包 rc |
|---|---|---|---|---|---|---|
| 改前 `internal/panel`（副本 `aeba6ff`） | 99 | **53** | **0** | **0** | 46/0/0 | 0 |
| 改后 `internal/panel`（工作树，代码 = `6e348f1`） | 105 | **58** | **1** | **0** | 46/0/0 | 1 |
| 改前 `cmd/wisp`（副本 `aeba6ff`） | 113 | **60** | **0** | **0** | 53/0/0 | 0 |
| 改后 `cmd/wisp`（工作树，代码 = `6e348f1`） | 116 | **62** | **1** | **0** | 53/0/0 | 1 |

⚠ "代码 = `6e348f1`"这一句本程现量过：`git diff --name-only 6e348f1 HEAD -- cmd/wisp internal/panel internal/agent/approval`
→ **0 行**。HEAD 之后又落了别家的 `docs/**` 提交（写本件这一程自己那几枚也在里面），
但那三枚路径下的**码**一字未动 ⇒ 这一行读数既可标 `6e348f1` 也可标 HEAD，**标错才是要防的**。

两枚红，逐枚点名：

1. **`TestC21DesignTokensFourWayAgree`**（`internal/panel`）——改前就有、**保持红、未修未跳**。
   ⚠ **引用枚数必须带口径**：**锚点／HEAD 副本 0 枚红／本机工作树 1 枚红**，两个数都对、
   拼成单值才错。因由：`tokens_fourway_test.go:441` 读 `design/assets/tokens.css`，
   而那批文件被 owner 以**未提交删除**挪走（`git status --porcelain` 里 16 枚 `D design/**`，
   本程一格没碰）。副本里 `design/assets/` 在归档内 ⇒ 同一枚用例在副本里绿。
   本程另在 `D:\tmp\wisp-35-pump-r2\mut`（**HEAD** 副本）复量：`internal/panel` 59 枚、
   `cmd/wisp` 63 枚、**0 红**（§4.6 的 E0 那一发）⇒ 与"锚点副本 0 枚"同口径，钉住了这条红
   **只由工作树形状造成、与代码版本无关**。
2. **`TestRunBooksWithASnapshotOfItsLiveQueue`**（`cmd/wisp`）——**本批新造的红**，逐字红句：
   `panel_pump_test.go:296: composer.mode.current = "ask_every_step", want the non-default档 this config set`。
   归因与隔离见 §4.6／§5.4 末段：**不是泵的缺陷**，是 `6e348f1` 只落了断言、没落夹具写入
   （`panel_pump_test.go:133` 只追加 `confirm_timeout_sec = 2`，`permission_mode` 一个字没写）。
   **本程不修**（改配置夹具让断言变绿在禁改清单上；改回去派单又明令禁止）⇒ **报回编队定夺**。
   最小闭合集合只有一行，且方向是收紧：把那枚键真写进 `[risk]`。

⚠ 判红只认 `--- FAIL:` 这条规矩在本批**当场救了一次**：两棵树的 `cmd/wisp` 输出里都有
一行 `        [FAIL] onnxruntime.dll version            file version 1.28.2.0 does not match build pin unknown`
（`TestAC2RealProcessRefusesOnEveryLegWithoutAppData` 内部 `t.Log` 打的自检文本），
**改前改后同一枚、字节相同** ⇒ 它是文本噪声不是红；按 `[FAIL]` 或按 `file:line` 判会凭空多一枚红。
另：`=== RUN` 若按"行里含 `=== RUN`"数会各多 1 枚（99→99、113→114），多的那枚是落在日志正文里的
同名字符串 ⇒ 本表用的是**行首精确**那支。

### 5.3 三门（本机，工作树 = `6e348f1`）

| 门 | 命令 | 读数 |
|---|---|---|
| `gofmt` | `gofmt -l internal/panel cmd/wisp internal/agent/approval` | **空**（rc=0） |
| `go vet` | `go vet ./internal/panel/ ./cmd/wisp/` | **空**（rc=0） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**，`d22scan: clean - no D22 ban violations` |

d22scan 那一枚要补三行，因为它是 `set -eu`（`scripts/d22scan.sh:38`）且**第一步就是正控**
（`:14` "the seeded-violation positive"）——正控红则真扫描根本不跑：

- **正控先响过才算数**：`runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70, '[no tests to run]'=0`
  ⇒ 真扫描确实跑了。
- **live scope 分母（带 HEAD）**：`bans #1-5 internal/=205 cmd/=23`、`ban #6 frontend/=46`、
  `ban #7 internal/tools/=18`、`ban #8 design/=32 frontend/=46 internal/=410 cmd/=42`。
- ⚠ `design/=32` 与门钉 r5 那件记的 **30** 差 2 枚：**因由已核不是本批**——本批对 `design/**`
  零改动（§3.2），差的是 owner 未跟踪的 `design/old/`、`design/doubao/demo/{lib,screenshots}/`
  等目录（`git status --porcelain` 里 4 枚 `??`）。分母随锚点变 ⇒ 引用时必须带"哪一版、工作树还是树内"。

### 5.4 名册差集（两本，逐枚列新增／消失）

**本 1：跑出来的名册（`-v` 里的顶层结果行）**

| 包 | 顶层枚数 改前 → 改后 | **新增**（逐枚） | **消失** | 子测试 改前 → 改后 | 新增／消失 |
|---|---|---|---|---|---|
| `internal/panel` | 53 → **59** | `TestThePumpBuildsThePacketFromWhatTheHostHolds`、`TestAPumpWithNoReadersSaysSoInsteadOfInventingState`、`TestAnUnreadableModeFromALiveReaderStillRendersUnknown`、`TestTheStreamLogMergesInsteadOfDropping`、`TestPublishWithoutAnExitReportsInsteadOfPretending`、`TestPublishHandsTheBytesToTheAttachedExit` | **无** | 46 → 46 | **+0 / −0** |
| `cmd/wisp` | 60 → **63** | `TestRunBooksWithASnapshotOfItsLiveQueue`、`TestSnapshotWorkspaceSectionReportsTheNarrowing`、`TestThePumpIsDrivenNotJustAssembled` | **无** | 53 → 53 | **+0 / −0** |

⇒ **上一程自称的"6 枚新用例入名册 53→59"核过：成立**（六枚逐枚点得出名字、无一枚消失、
子测试本 0 增 0 减）；`cmd/wisp` 自称的"60 → 63"同样成立。两包都**没有 panic 吞读数**的形状
（顶层结果行枚数与 `=== RUN` 顶层枚数逐包相等：99=53+46、105=59+46、113=60+53、116=63+53）。

**本 2：源码里的静态声明（`^func Test`，含被 build tag 挡掉的）**

| 包 | 改前 → 改后 | 新增 | 与"本 1"的差 |
|---|---|---|---|
| `internal/panel` | 53 → 59 | 同上六枚 | **0**（两本一致） |
| `cmd/wisp` | 63 → 66 | 同上三枚 | **−3**，改前改后**同一枚数** |

⇒ 那 3 枚的口径已查明、**不是本批造成**：`cmd/wisp/secret_dataroot_119b_test.go` 带
`//go:build !windows`、内含 3 枚 `func Test`，本机（windows）根本不编译它们；
其余 6 枚带 tag 的测试件是 `windows` tag，本机照跑。⇒ **引用"某包几枚用例"必须说得出是哪一本。**

### 5.5 本程有没有污染被验物

| 检查 | 读数 |
|---|---|
| 仓库工作树里本程改过的代码件 | `git status --porcelain -- cmd/wisp internal/panel internal/agent/approval` → **空**（11 发变异全在 `D:\tmp\…\mut` 里，每发后还原） |
| 变异副本有没有停在改过的状态 | 五枚被改文件逐枚 `diff <(git show HEAD:<f>) mut/<f>` → **五枚全等** |
| 本程的 commit 各带了什么 | 每枚 `git commit -q -F - -- docs/evidence/s1/35-panel-snapshot-pump-r1.md`；`git show --numstat` 逐枚只列这一枚路径，删除列 **0** |
| 别人的东西有没有被卷走 | `design/**` 16 枚 owner 未提交删除 ＋ 4 枚未跟踪目录 ＋ `docs/reports/frontend-session-log.md`（别家在写）**全程原样**，未出现在任何一次暂存清单里 |
| 加载期失败有没有被偷偷避开 | 四份读数里 `0xc0000135`／`error while loading shared libraries` **各 0 命中**（两棵树都装了 DLL）⇒ 派单预告的那枚"改前也红、一条尺都不执行"的形状**本程没遇到**，也因此**没有复现它**（列 §7） |

---

## 6. 上游简报与上一程自述里，哪一句不成立（逐句复算）

判据：**每句都当未验证断言重走一遍**；不成立的写清"不成立的是数值还是结论"，因为这两件事的处置不一样。

| # | 那句原话（出处） | 复算结果 |
|---|---|---|
| E1 | "日志管道把任何单条字符串截到 512 字符"（上一程中断前最后一句） | **成立**。现量见 §2.2（900 rune → 落盘 543，含 `truncate` 自加的 31 rune 后缀）。唯一要补的是口径：512 是**保留部分**的上限、不是落盘总长 |
| E2 | "`NewSnapshot`/`NewComposerState`/`NewModeView`/`ModeUnknownView`/`UnsetWorkspaceView` 从此有**包外**调用者"（`5821e24` 的 commit message） | **一半不成立**，逐枚见 §6.1。真正成立的表述是"这些建造者从此有**非测试**调用者、且那条链由生产驱动"；"包外"只对 `UnsetWorkspaceView` 一枚字面为真，`NewModeView` 那枚**前后一枚数、调用点没动过** |
| E3 | "ledger 单串上限 512 而最小包 **831** 字节"（`37a4705` message，成立）／代码注释写"最小包 **534** 字节、带卡那枚 **997**"（`cmd/wisp/panel_pump.go:138-139`、`cmd/wisp/panel_pump_test.go:27`） | **结论成立、注释里两枚数值都不成立**。实测：库层最小包 **585 字节 / 551 rune**；带卡的包实测七档 **803／811／812／830／831／837／893**（§6.2 有逐枚来路），**534 与 997 一次都没出现**。⇒ 这是**注释里的读数腐坏**、不是判断错：即使按 rune 口径 551 > 512，那个包仍装不进 ledger，所以"落有界摘要"的决定仍然对。**本程不改那两行注释**（改码不在本程三件交付物里，且上一程已把 831 写进 commit message、注释与它自相矛盾，谁对要由实现者那一族定）⇒ **报回** |
| E4 | "6 枚新用例入 `internal/panel` 名册（53 → 59）" | **成立**，且本程补了派单要的两本差集（§5.4）：六枚逐枚点得出名字、**消失 0 枚**、子测试 46 → 46（+0/−0）；`cmd/wisp` 自称的 60 → 63 同样成立（新增三枚具名、子测试 53 → 53） |
| E5 | "`Snapshot` 四键与 `panel.ts` 未动、双向尺 `composer_test.go:73-78` 仍绿" | **成立**，且本程加了一发反向证明它**不是恒绿**（§4.3 的 M6：改一枚 json 键 → `composer_test.go:74` 与另一向 `approval_test.go:132` 同时响）。⚠ 附一条口径：`cmd/wisp` 那三枚在 M6 下**不响**，别把它们的绿当第二枚凭证 |
| E6 | "逐字段拿另一路生产读数对账，**无一是本文件填的值**"（`37a4705` message） | **过宽**。真的是另一路生产读数的有 5 处（level／tool／reason／rulesHit 对 `consoleCard145` 那枚控制台卡、mode 对 `perm: MODE-READ` 审计行、results 的 correlationId 对 run 打印的 task id、results 的 text 对 stdout 流出的正文）；**但另有 4 处比的是本文件自己写的字面量或常量**：`panel_pump_test.go:218` `card.CorrelationID != "pump-corr"`（那枚 id 是本文件交给真 `bridge.Execute` 的）、`:295` `!= risk.ModeAskHighRiskName`、`:327` `strings.Contains(chunk.Text, "echo:")`、`:279` `!= panel.MaxAttachmentBytes`。⇒ **这四处不是缺陷**（值确实穿过生产对象，只是终点是一枚本文件写的期望值），但"无一"那两个字说过满了 ⇒ 按现量报名排者 |
| E7 | "顺手把 `workspace canonical` 其实是 `foldPath` 折过的比较键这一处复算不符**钉进用例与证据件**"（`37a4705` message） | **用例那半成立，证据件那半在落盘时不存在**。`cmd/wisp/panel_pump_test.go:372-378` 逐字写着"written up in `docs/evidence/s1/35-panel-snapshot-pump-r1.md` **§3**"，而本件当时只有 §0／§1（`ddf8b6f`，派单自己也是这样记的）。⇒ **这正是派单点名要防的那一族（注释预先引用尚未产出的读数），上一程自己犯了一次。** 本程：①把这条 finding 真核为假（§6.1′）②写进本节 ③**不改那枚已提交的注释**，只点名"编号与指向仍不符（它指 §3、内容在本节）" |
| E8 | §1.1"那份 C17 冻结名册 Go 侧 0 枚"，配的是 6 枚名字的抽样 | **仍然成立，且本程从 6 枚量到 18 枚全量**（§6.3）。两枚看似非零的逐枚处置掉了。⚠ 另记一条仪器坑：原命令 `grep -rn --include=*.go "panel.resync"` 里的 `.` 是**正则通配**，会误命中 `internal/agent/approval/queue.go:288` 那种路径 ⇒ 复测改用 `-F`（固定串） |
| E9 | §1 的总结论"通道不存在"——**这条现在是否已被它自己改变**？ | **一半变了、一半没变，逐格给**（§6.4）。⇒ 准确说法：**"那份包有没有人在生产里造"这一格被本批改掉了；"Go → 面板这条通道不存在"没有被改掉**（`PanelBridge` 类型声明前后各 0 枚、C17 名册 18 枚全 0、`ParseComposerRequest` 非测试调用者前后各 0、`Marshal()`／`lastPanelSnapshot()`／`snapshotCount()` 非测试调用者各 0） |
| E10 | 派单"`internal/panel/` **改前就有 1 枚红**" | **只在其中一棵树成立**：本机工作树 1 枚、锚点副本 0 枚（§5.2 已按两口径写）。这一条派单自己也预告了"两个数都对、拼成单值才错"⇒ 本程照写，并**另在 HEAD 副本复量第三格**（59／63 枚、0 红）钉住"它只由工作树形状造成" |
| E11 | 派单"上一程说 `Go→面板` 这一跳不存在（WebView2 API 0 命中、无 HTTP/SSE/ws、`runResident` 原本不 import `panel`）" | **三条全部复算为真**（§2.1.1）。其中"原本"两字有下文了：`runResident` 现在**仍然**不 import `internal/panel`（0 命中），被接上的是 **`wisp run` 那条腿**、不是常驻腿 ⇒ 若有人以为"常驻进程已经在推快照"，那是**不成立**的 |

### 6.1 E2 的逐枚现量（"包外调用者"这一句为什么只算一半成立）

改前的数取自 `D:\tmp\wisp-35-pump-r2\base`（锚点副本），改后取自工作树；
尺：`grep -rn --include=*.go "<sym>(" .` 去 `_test.go`、去注释行、去 `func ` 定义行。

| 建造者 | 改前非测试调用点 | 改后非测试调用点 | 新增那几枚在**包外**吗 |
|---|---|---|---|
| `NewSnapshot` | **0** | 2（`internal/panel/pump.go:138`、`:182`） | ❌ 同包 |
| `NewComposerState` | **0** | 1（`pump.go:171`） | ❌ 同包 |
| `NewModeView` | **1**（`composer.go:202`） | **1**（还是 `composer.go:202`） | ❌ **一枚都没新增**——那句对它不成立 |
| `ModeUnknownView` | **0** | 1（`pump.go:175`） | ❌ 同包 |
| `UnsetWorkspaceView` | 2 | 4（新增 `cmd/wisp/panel_pump.go:84` 与 `pump.go:163`） | ✅ **`cmd/wisp` 那一枚是包外** |
| `CardViewFromDecision` | 1 | 2（新增 `pump.go:72`） | ❌ 同包 |
| `WorkspaceViewFromRoot` | 1 | 2（新增 `cmd/wisp/panel_pump.go:86`） | ✅ 包外 |

⇒ **成立的那半**（也是本票段真正要的那格）：`panel.Snapshot` 今天有了一条**从运行进程出发可达**的建造路径——
`cmd/wisp/run.go:421` `panel.NewSnapshotPump(...)` → `Publish()` → `Snapshot()` → `NewSnapshot(...)`。
`NewSnapshotPump` 本身非测试调用点 1 枚、且在包外（`run.go:421`）。
⇒ **不成立的那半**：五枚里只有 2 枚真 gained 包外调用者，1 枚（`NewModeView`）什么也没多。
**这不是行为缺陷、是自述过宽**，写下来是为了让下一程别拿这句话当"名册已外扩"的凭据。

### 6.1′ E7 里那条 finding 本身：复算为**真**

`internal/tools/paths_workspace.go`：`:82 SetWorkspaceRoot(root)` 里 `:86 f := foldPath(root)`、`:93 p.workspace = f`；
`:40 WorkspaceRoot()` 直接 `return p.workspace`。⇒ 泵经
`cmd/wisp/panel_pump.go:86 panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())` 拿到的是
**折过的比较键（小写形）**、不是 C26 授权的那个拼写。用例 `panel_pump_test.go:379-387` 因此
按 `strings.EqualFold` 比、并在两枚拼写**相同**时才 `t.Log` 记"没复现"。
本程现量：`post_cmd.txt` 里那句 `fold-key finding did not reproduce on this run` **0 次出现**
⇒ 本机这一发两枚拼写确实不同 ⇒ finding **可复现**，不是纸面推断。

### 6.2 E3 那两枚数值的逐枚来路（免得不成立被读成"没量过"）

| 观测到的包大小 | 来路 |
|---|---|
| 585 B / 551 rune | 库层最小包（`pending:[]`、`results:[]`、默认 composer），从 M4 那发的 `pump_test.go:292` 红句里取串量字节 |
| 803 | M5（面板侧重算那发，卡里 rules 变空、reason 变短） |
| 811 / 812 | M1 / M1b（reason 换成常量句） |
| 830 | E0 / M3b / M4 / M4b / M7 |
| 831 | 工作树现状（`post_cmd.txt`）、M3、M3_no_wfix —— 也就是 `37a4705` message 里那枚 831 |
| 837 | M6（json 键改名，多出来的 6 字节就是 `Renamed`） |
| 893 | M2（results 被换成常量 chunk） |
⇒ 七档全是**同一枚用例的同一枚采样**在不同变异下的值，且每一枚的 `packet bytes (N)` 与打印串的 UTF-8 长度**逐枚相等**（脚本比过，13/13 match）⇒ `len(data)` 这个字段自己是诚实的；不诚实的是那两行注释里的 534／997。

### 6.3 E8 的 18 枚全量（改前只抽了 6 枚）

尺：`grep -rnF --include=*.go "<名字>" internal cmd`，去 `_test.go`、去 `tools/d22scan`（那是 ban #6 的扫描模式串，不是路由）。

| 名字 | 非测试命中 | 处置 |
|---|---|---|
| `panel.resync` `tasks.list` `tasks.detail` `history.query` `transcript.get` `approval.current` `approval.decide` `config.get` `config.set` `grants.list` `grants.revoke`(见下行) `privacy.purge` `privacy.export` `cost.summary` `models.list` `models.delete` `diagnostics.export` | **各 0** | ⇒ 名册确实一枚都没落地 |
| `approval.queue` | 1 | **前缀巧合**：命中的是 `internal/statemachine/events.go:43 EvQueueDrained Event = "approval.queue-drained"`（D43 转移表的事件名），不是 `approval.queue` 那枚方法 |
| `grants.revoke` | 2 | **不是字符串**：命中的是 `internal/agent/approval/queue.go:288`／`:374` 的 Go 方法调用 `it.grants.revoke()` |

### 6.4 "0 枚调用者"这一条被它自己改掉了什么（派单点名要答的那格）

| 那一格 | 改前 | 改后 | 判定 |
|---|---|---|---|
| `Snapshot` 有没有生产可达的建造者 | **无**（`NewSnapshot` 非测试调用点 0） | **有**（`run.go:421` → `Publish` → `Snapshot()`） | **已被本批改掉** |
| 那份包有没有被真运行驱动 | 无 | 有（三处生产触发点，`snapshotCount()` 由 0 走到 ≥1，`TestThePumpIsDrivenNotJustAssembled` 钉住） | **已被本批改掉** |
| 全量字节有没有人接 | — | **仍无**（`lastPanelSnapshot()`／`Marshal()`／`snapshotCount()` 非测试调用者各 0） | 未改掉 |
| `Go → 面板`这条通道 | 不存在（§1.2 四条） | **仍不存在**（§2.1.1 四条重测，全 0） | 未改掉 |
| `PanelBridge` 类型声明 | 0 | **0** | 未改掉 |
| C17 冻结名册（18 枚） | 0 | **0** | 未改掉 |
| `ParseComposerRequest` 非测试调用者 | 0 | **0**（前后同一枚数，且 `run.go:225` 那句"Nothing calls it yet"**未被推翻**） | 未改掉 |

⇒ 给编排者的一句话：**§1 那份探针的结论没有被推翻，被推翻的是它的第 4 格（"那份包只有测试能造"）。**
如果下一程或票面把这两件事混成"通道已存在"，那是错的；反过来，如果还有人拿 §1.1 说
"`Snapshot` 没人造"，那也过期了——过期时刻就是 `5821e24`／`37a4705`。

---

## 7. 本程没测什么（按"如果我漏了它，谁会先被骗"排序）

1. **没跑 CI 那支形状（正路①）**。`scripts/wisp-cli-tests.sh` 与它背后的
   `portable-tests.sh`（GUARD A"这一档必须有已编译的测试文件"、GUARD B"cmd/wisp 必须打出自己的
   锚定顶层结果行"）与 `tools/d22scan/runtests.sh`（顶层 SKIP 即红、PASS 与 FAIL 同时为 0 即红）
   **都以手工等价尺替代了、没有以 CI 同形跑过**。⇒ 本件的 5.2 那张表**不等于 CI 会打出的那张表**；
   尤其"某包一枚测试都没有"这一类失败，本程那把尺抓不到。
2. **没跑 `-race`**。`grep -rn -- "-race" scripts/portable-tests.sh tools/d22scan/runtests.sh` → **0 命中**，
   本程也没单加。而 `StreamLog` 与 `SnapshotPump` 都是带锁对象、`Snapshot()` 又
   **不持泵锁**地先后调 `Verdicts()` 与 `Results()` ⇒ 一份包可以混合两个瞬间的状态
   （`pending_read.go:38-41` 的注释承认了这件事、并说这是设计）。**承认了 ≠ 量过**：
   本程没造这一发竞态用例、也没造这一发变异。
3. **`panelSnapshotSummary` 那条"超钳就丢 ids"的分支从未被走到**。
   `grep -rn "summaryClamp\|440" --include=*_test.go cmd/wisp` → **0 命中**；
   13 发变异 + 工作树那发里，落盘记录里的 `pending=` 字段**只出现过一枚值** `pump-corr`
   （`grep -hoE "pending:[^ ]*" … | sort -u` → 1 行）。⇒ `cmd/wisp/panel_pump.go:186-192`
   那一整块（深度 >3、或多枚长 correlation id 把行撑过 440）今天是**零执行**的分支。
   ⚠ 按本仓教训"改后全绿可能是那条新路径零次执行"，这一格**不能**被算成"已经守住了"。
4. **rune/byte 分叉那一发没造**（§2.2 口径第 2 条）。今天所有字段都是 ASCII ⇒
   `len(line) > summaryClamp`（byte）与 `truncate`（rune）不会分出不同结果；
   一旦 correlation id 或 mode 名里出现非 ASCII，两把尺会开始说不一样的话。未测。
5. **没测多任务并发**。`DefaultStreamKeys = 32`（`pump.go:265`）对今天的唯一生产者（`wisp run`
   一个 task）是 1∶32 的空档；`TestTheStreamLogMergesInsteadOfDropping` 用 50 枚 key 量了折叠
   本身，但**没有任何一发**量"生产里真的会有几路流"。
6. **没测常驻腿、便携腿、Linux 腿**。`cmd/wisp/resident_windows.go` 今天不 import `internal/panel`
   （§2.1.1 现量 0）⇒ 被接上的只有 `wisp run`；三枚 `*_other_test.go` / `//go:build !windows`
   那 3 枚（§5.4 本 2 的差）**本程未在容器里跑**。
7. **没跑 SLO／D32 那两个数**，也没跑 `wisp slo`（且按台账既有裁定：`wisp slo` 全仓 100% 不被
   `go test` 执行 ⇒ 这一条腿只有推送触发的 `slo-check.ps1` 管，本程未推、未读）。
8. **没跑全树门禁**：`go vet ./...`、全仓 `go test ./...`、`staticcheck`/`gofumpt` 都**未跑**
   （本程范围是派单点名的两包＋三门）。
9. **没复现派单预告的加载期失败**（`0xc0000135`）：两棵树都装了 DLL，四份读数各 0 命中。
   ⇒ 那一句在本件里是**引用**、不是复算。
10. **变异只打在 `internal/panel` ＋ `cmd/wisp` 两包**（外加 `approval/pending_read.go` 被 M7 间接触碰）。
    没有把每一发变异打遍"引用过这些符号的其它包"（`go build ./...` 意义上的全树），
    也没有做"删掉整枚 `pump.go` 会不会有别的包变红"这一发**移除级**实验
    （M5/M6 是替换级、不是移除级）。
11. ~~没量那枚新红在别的机器上的形状／它吃不吃 `design/**`~~ —— **本程把这一刀补了**，记在这儿因为它差点变成误归因：
    在 HEAD 副本里（`design/assets/tokens.css` **存在**，`ls` 现量过）单跑那枚用例
    `go test -count=1 -v -run 'TestRunBooksWithASnapshotOfItsLiveQueue' ./cmd/wisp/`
    ⇒ **仍红、同一句 `panel_pump_test.go:296`，且它是唯一一句红**（`:293` 不响，与 §4.5 的 `M3_no_wfix` 同形）。
    ⇒ 那枚红与 `design/**` 的形状**无关**，§5.2 第 2 条的归因（夹具写入缺一行）成立。
    **仍未测的**只剩：这一枚红在别的机器／别的 runner 上是否同样唯一（本程只有本机一块地）。
12. **没测那枚有界摘要在真实 ledger 文件里的行长**。`ledgerSummaries145`（`panel_pump_test.go:53-79`）
    读了真文件、比了字段，但**没断言过那条 `msg` 的 rune 长度** ⇒ "摘要自己永远不会被 `truncate` 截到"
    这件事今天是**推论**（440 < 512），不是断言。§2.2 的探针量到的是**探针自己造的**那枚摘要（180 rune），
    不是 `wisp run` 真落的那一枚。

### 7.1 一条仪器异常（登记，不据此下任何结论）

本轮若干次 `grep -c` 为 0 时，工具回显里除 `0` 之外还出现过一句
`No matches found`——那**不是** bash `grep` 的输出形状。本程处置：凡关键计数都用
`| wc -l` 或 python 复算过一遍（§3.2、§5.2、§6.3 的每个零都是这么来的），
未把该句当结果读。⚠ 出处：`Bash` 工具、命令形如 `grep '^--- FAIL' pre_test_v.txt`。
只登记形状，不猜成因。

---

## 8. 给 owner 的一段（不夹术语）

**这次改的是什么（一句人话）**：程序内部现在每次"该显示在助手面板上的那张状态表"变了，
就当场真的把那张表算一遍、交出去一次。以前这张表**只有测试代码能算**，真跑起来从来没算过。

**你屏幕上的东西变了吗**：**没有，一点没变。** 那张表现在交到程序自己手里，还没有一条路能通到
你看到的那块网页。把网页装进窗口、以及把这张表推过去，是另外两张工单（票 33、票 35 的另一半）。

**这次多出来的、你唯一能看见的东西**：跑一次 `wisp run`，你的日志目录（默认是
`%APPDATA%\wisp\logs\`）里，每有一次面板状态变动就多**一行短记录**——时间、几张待确认卡片、
权限档位、工作区有没有被收窄、回答流了几段、这张表多少字节、一个指纹。
里面**不写对话正文、不写密钥**（密钥、音频、网页正文那类走的是早就有的那套强制遮蔽，这次一个字没动）。
最坏后果形状：**日志文件比你以为的长得快一点**（一次回答里，工具每开始一次、结束一次各多一行）。
不是泄露，是体积。

**一条要你注意、但今天不是事故的事**：这张表里带着"待确认卡片"的内容——工具名、**完整参数**、
判定的等级、命中了哪几条规则、以及为什么。今天这些字只打印在你自己的终端上，那张表也还没送到任何页面。
等票 33 把那根管子接上那天，这些字段就是页面上会显示的那些字的来源。
最坏后果形状：**将来**那块面板会显示这些内容；**今天**没有任何东西被送到页面——
这一句我是量出来的，不是听说的（能送到页面的那条通道，枚数还是零）。

**安全上"没做"的那件事，说清楚**：这次**没有**顺手接上任何一条"从面板进来的批准"。
我没有拿"我搜了没搜到"当证据——我用一把**会亮的尺**量的：同一条尺打在历史上真接过入站路由的两枚
提交上分别亮 **6 次**和 **3 次**，打在这一批上只亮 **1 次**，而那 1 次是一枚**只读**的枚举
（把队列里已经存在的卡片抄一份出来），它不允许任何东西、也不会把任何东西挪出"待确认"。
所以"面板侧来源的允许"这一条硬禁线**没被碰**。最坏后果形状：如果将来有人在这条线上接错一次，
这块面板就能自己批自己——那正是这套设计从头到尾在防的那件事，所以它被单独验了一遍。

**一处诚实的空白**：卡片上有个标记，意思是"这次不能靠本次会话的默认放行"。今天**没有任何一条测试
真的验过它从真实队列抄对了**。原因不是漏写测试，是**这条路径上今天产生不出"真"值**
（要等那个内容来源探测器，票 25）。所以我没有补一枚"永远为真、所以永远绿"的测试来把这一格勾掉——
那种绿是本仓最贵的一种假。

**门禁里现在有一枚新红，需要你点头**。上一程被切断时留下一枚只写了一半的测试配置：
它加了一条**更严**的检查，但忘了把那条检查需要的前提写进测试用的配置文件里，所以那条检查现在**恒红**。
- **甲（我倾向这个）**：把那行配置补上。方向是**收紧**，不放松任何断言。我在一份仓外副本里试过：
  补完之后两包 59／63 枚全绿；而且我故意把六个字段分别从"读真实值"改成"填假值"，
  **每一发都仍然把测试打红** ⇒ 补那一行不会磨钝任何一根钉子。
- **乙**：把那条新检查删掉。方向是**放松**——删掉之后，"把权限档位偷偷写死成默认值"这一种错
  就再也没有任何测试抓得住。**这一句我也是实测出来的**：我照乙的状态打了一发，
  唯一还会响的那根针正好就是被删的那根。

我两个都没做。两个都超出这一程的授权（为了让红变绿去改测试配置在禁改清单上；
删别人写下的断言也在我这程的禁改里）。

**你需要做什么**：回一个"**甲**"或"**乙**"就行。其余不用你动手。

> **编排者处置（17:1x 追加，原段落一字不抹）**：**这一格不该问 owner，我按甲做了**（commit `1485921`）。
> 理由一句话：**乙 那一支是"删掉一条断言换取红色消失"**，而"为了变绿放宽断言"在 owner 自己定的硬禁清单上
> （`AGENTS §1.1`）⇒ 这一支**没有第二个可选项**，摆给他看等于让他替我做一件本不该有选项的事。
> 落地后的现量两发（本机把仓里 `third_party/sherpa-onnx` 三枚 DLL 放进 PATH）：
> `go test ./cmd/wisp/ -run TestRunBooksWithASnapshotOfItsLiveQueue` ⇒ **rc=0，PASS**，
> 且包体里 `"current":"ask_high_risk"`、日志行 `mode:ask_high_risk` ⇒ 那条断言是**被真读出来的值满足的**，
> 不是被常量喂绿的；`go test ./internal/panel/ ./cmd/wisp/ -count=1` ⇒ `cmd/wisp` **ok（71.0 s）**，
> `internal/panel` 唯一红仍是 `TestC21DesignTokensFourWayAgree`（本机样式表被挪走，改前也红，保持红）。
> ⚠ **一句归因要在前面**：那枚恒红**是我造的**——我代提那程的 WIP 时只跑了 `go vet`、没跑该包测试
> （见 `6e348f1` 正文里我自己写的那句"我没跑这个包的用例"）。⇒ 结论没变（**代提方向仍是对的**：
> 共享树里留未提交件的真风险是被下一枚不带 pathspec 的 commit 卷走），但**"vet 过"不等于"测试过"**
> 这一条今天又付了一次学费，落进 §7 那一族。

---

## 9. 临时件路径（只建不删，请编排者一次清）

全部在**仓外**，仓库目录内未建任何 worktree／checkout／临时件：

```
D:\tmp\wisp-35-pump-r2\
  base\                      锚点 aeba6ff 的归档副本（含从工作树拷进去的 third_party/sherpa-onnx/*.dll 三枚）
  mut\                       HEAD 的归档副本 = 11 发变异的靶树（同上三枚 DLL；每发后已还原，末态 == HEAD）
  probe\                     HEAD 的归档副本 ＋ truncprobe\main.go（512 截断的真管道探针）
  anchor_files.txt           aeba6ff 全 1174 枚路径清单（副本逐枚比对的输入）
  pre_test_v.txt             两包合跑那一发（**已被 §5.1 记为不可归因，别引用它的分包数**）
  pre_panel.txt  post_panel.txt
  pre_cmd.txt    post_cmd.txt     四数与名册差集的原始 -v 读数
  d22scan_post.txt                 d22scan 全量输出（正控 + 真扫描）
  gates.sh                           上面四发的驱动
  mutate.py                          11 发变异的驱动（含每发的锚点字符串）
  mutation_summary.txt / _summary2.txt  变异的屏显汇总
  mutlogs\                           每发一枚 .log（红名＋红句）与一枚 .raw.txt（原始 -v），共 13 组。
                                     ⚠ 其中 M2_results.{log,raw.txt} 是**已被作废的那一发**
                                     （替换串非法 Go → `[setup failed]`，§4.1 明写它不算读数），
                                     留着是让"哪一发不算数"本身可核；M2_results_fix.* 才是 M2 的正式读数。
                                     另有 orig\ 目录 = 五枚被改文件的原件备份
```

⚠ 三条给清理者的话：①`base\`／`mut\`／`probe\` 里各有三枚 17 MB 上下的 `onnxruntime.dll` 拷贝，
是**体积大头**，`mutlogs\*.raw.txt` 是**读数本体**、被 §2／§4／§5 逐枚引用，删了这些格子就从
〔可复现〕掉回〔仅自述〕；②本程**一个 `rm` 都没执行**（连续两程因这枚被记名）；
③`mut\` 已核回 HEAD（§5.5 那五行 `diff`），如果编排者要复算任何一发变异，直接改 `mutate.py` 的
参数重跑即可，靶树是干净的。

---
