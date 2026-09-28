# `33-a1` — 面板入向那一跳（网页点一下 → Go 收到）：地基普查与代价表

- 程：`33-a1`（只读设计核·零产码）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`
- 派单：`.scratch/wisp/dispatches/2026-09-28-150x-readonly-33-a1-inbound-hop-design-core.md`
- 起手时刻：`2026-09-28 14:55 +0800`｜起手 HEAD：`ee0ef8ec`（派单锚 `1f22f1aa`，**已漂**，见 §1）
- 性质：**只读** ⇒ AC 框一枚未勾、`internal/**`／`cmd/**`／`go.mod` 零字节不动。

---

## 1. 起手锚＋写面闸门

| 项 | 现量（本程自己跑，非引用） |
|---|---|
| 起手 `date` | `2026-09-28 14:55 +0800` |
| 起手 HEAD | `ee0ef8ec`（`docs(evidence): 174-c2 correct the tool-call tally in section 9 (33 of 35, exploration stopped at 23)`） |
| 派单给定锚 | `1f22f1aa` —— **已漂**（按实际 HEAD 做，不改判） |
| 漂移核法 | `git merge-base --is-ancestor 1f22f1aa HEAD` = **YES**（同一条线，非分叉）；`git log --oneline 1f22f1aa..HEAD \| wc -l` = **4**（含本程自己那枚 `40aee084`） |
| 起手时在本锚之后的他人 commit | `ad759e94`（182-c1 面板字段普查）· `6479d540`（174-c2 接线代价普查）· `ee0ef8ec`（174-c2 更正计数）—— **全是别家的证据件/探针**，本程未读作依据、未 add、未改动 |
| 写面闸门（起手） | `git status --porcelain -- internal/ cmd/ go.mod` → **空**（14:55 实测） |
| 分支 | `dev` |
| 本程性质 | 只读设计核：**AC 框一枚未勾**（票 33 的 8 格全部保持未勾），`internal/**`／`cmd/**`／`go.mod`／`go.sum`／`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden **零字节改动** |

判读法沿用票 33/35 预检那张表的三档：**已存在**（有实现且有生产听众）／**半存在**（有实现但没人取用，或缺一半）／**完全不存在**（给不出 `file:line`）。
⚠ **行号一律本程现跑**：票 33/35 预检（`docs/evidence/s1/33-35-preflight.md`，锚 `dbbc822`）里的行号已确认漂移两处——它写 `ParseComposerRequest` 在 `bridge.go:77`（**现为 `:84`**）、`rt.modeWrites` 装配在 `run.go:378`（**现为 `:417`**）。本报告不复用它那两个号，凡引用都重跑。

---

## 2. Q1 规格原文：这一跳被定成谁的责任

以下全部**逐字抄**（含加粗与标点），行号现跑。

### 2.1 `docs/specs/SPEC-08-ui-ball-panel.md`

`:141` 节标题：

> `## 5. 面板（WebView2，S5）`

`:143` 节标题 —— **这就是"把它定成宿主层责任"的那一节**：

> `### 5.1 宿主（D29/C27）`

`:145-147`：

> `- \`jchv/go-webview2\`（MIT，纯 Go 无 cgo）；**单例 \`PanelManager\`：一会话至多一个 WebView2 窗口，`
> `  隐藏而非销毁**——既绕开 Environment 共享问题，又把绝大多数交互压到热路径（冷 ≤1500ms /`
> `  热 ≤200ms，D32 修正）。`

`:148-149`：

> `- 资源加载：\`AddWebResourceRequestedFilter\` 从 \`embed.FS\` 喂前端产物——**不起 localhost HTTP`
> `  服务**（避免 dsh-tauri 的 iframe-to-localhost 模式，D21/D29）。`

`:150-151`：

> `- **前端必须无状态**：WebView 销毁/隐藏后状态全丢；每次 \`show\` Go 侧推 \`panel.resync\` 全量状态；`
> `  前端不得缓存任何跨 show 的业务状态（历史读 SQLite、配置读 TOML）。`

`:152-153`：

> `- WebView2 Runtime 缺失 → **无面板模式**：D10 短结果 + 落文件仍可用 + L2 走原生降级卡；`
> `  引导安装 Evergreen Runtime，**不得自动下载安装器**（D42#11）。`

`:154`：

> `- 多任务共用单窗口，按 correlationId 分区渲染。`

`:156` 节标题（**白名单今天定到哪一步，答案就写在这行的方括号里**）：

> `### 5.2 C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】`

`:158-159`：

> `\`invoke(method, args) → result\` + Go→前端事件推送（流式结果/审批请求/任务状态），回复按`
> `correlationId 路由；每方法标注 capability 与是否需原生侧授权；未列出方法名 → 拒绝并记日志。`

`:167`（那张表里唯一涉及"面板侧能不能批准"的一行）：

> `| \`approval.decide\` | invoke | **「allow」拒绝一切面板来源（F2）；仅 \`reject\` 可面板发起** |`

`:138`（承载分层表 L2 那一行，说明宿主失败不许吃掉能力）：

> `| L2 强确认卡 | WebView（面板）/ 原生降级卡（面板不可用） | 完整命令/参数/风险说明是 CSS 强项；此时已在工作态；**L2 能力不得因面板缺失消失（C27）** |`

`:235-236`（测试侧对 C17 的要求）：

> `- C17：correlationId 路由测试（并发 10 个请求乱序回复）、方法白名单拒绝、resync 后前端状态与`
> `  Go 侧一致（面板「无状态」用重开面板断言）。`

### 2.2 `docs/PLAN.md`

`:1367` —— C17 契约行（**"网页事件进 Go"这一跳在冻结契约里的正身就是这一句**）：

> `| **C17** | **\`PanelBridge\`** | 前端↔Go 双向通道：\`invoke(method, args) → result\` + Go→前端事件推送（流式结果、审批请求、任务状态）。**回复必须按 correlationId 路由**；前端**必须无状态**（WebView 销毁后一切从 Go 侧重读） | D29, D31 |`

`:1377` —— C27 契约行（宿主本体）：

> `| **C27** | **\`PanelManager\`** | 单例，**唯一** WebView2 窗口持有者；一会话至多一个面板窗口，**隐藏而非销毁**；面板不可用 → **L2 降级为原生最简确认卡，L2 能力不得消失**；多任务共用，按 correlationId 分区渲染 | D32/D39（补 WebView 无主问题） |`

`:2434-2435` —— 白名单的强制出处：

> `**配套**：C17 \`PanelBridge\` 必须给出**方法白名单**（只有列出的方法可被前端调用），`
> `每个方法标注所需 capability 与是否需要原生侧二次授权。未列出的方法名 → 直接拒绝并记日志。`

`:2968-2970` —— D39 契约深化里 C17 补的项：

> `**C17 \`PanelBridge\` — 补四项**：`
> `① **方法白名单**（只有列出的方法可被前端调用，未列出 → 拒绝并记日志）`
> `② 每个方法标注所需 capability 与**是否需要原生侧二次授权**（F2 的落地）`

`:981` —— D29 标题行（选型已定）：

> `**Raycast/cmdk 视觉语言 + React + TypeScript + Tailwind + shadcn/ui；WebView 宿主 \`jchv/go-webview2\`。**`

`:1273` —— 模块表里 `panel` 那一行（宿主只能按需起）：

> `| \`panel\` | 按需 WebView 宿主、命令面板、结果面板、配置编辑器 | 常驻（**禁止**） |`

`:1568` —— P11（冷拉起那条待定项原文，见 §5）：

> `**若实测冷拉起 >2s，必须重新评估 L2 确认卡是否改回原生**（会牺牲 D29 的视觉质量）。这个决定不能在 S5 才发现要做`

### 2.3 票面怎么分的（这决定"缺的那一行归谁"）

- 票 33 `:11-14` What to build 点名的是：`singleton PanelManager owning at most ONE WebView2 window per session (hide-don't-destroy), embed.FS resource serving via \`AddWebResourceRequestedFilter\` (no localhost server), cold/hot show paths, focus return, and the no-runtime fallback declaration consumed by 37's native card.`
- 票 33 `:33` Out of scope：`- Frontend bundle contents (34); bridge protocol (35); panel pages (36–40).`
- 票 35（`.scratch/wisp/issues/35-panel-bridge-c17.md`）`:7` Parallel slots：`≤2 sub-agents (A: bridge transport + dispatch + whitelist; B: resync model +`；`:6` `**Blocked by:** 33-panel-host-c27, 34-frontend-scaffold`。

### 2.4 ⚠ 本程查到的一枚**规格级缺口**（只报，不填）

**"网页事件 → Go 里那段 raw JSON 由谁取出"这一句，在 SPEC-08 §5.1、票 33、票 35 三处的票面/规格里都没有逐字点名。** 现量：

- SPEC-08 §5.1（`:143-154`）那五条宿主责任点名了 生命周期 / 资源服务（`AddWebResourceRequestedFilter`，**出向**）/ 无状态 / Runtime 缺失 / correlationId 分区——**没有一条写"消息接收"**；
- 票 33 全文 `grep -nE "WebMessage|postMessage|inbound|receive"` = **0 命中**（本程跑：票 33 只提到 `AddWebResourceRequestedFilter`）；
- 票 35 把 `bridge transport + dispatch + whitelist` 划进 slot A，但 **transport 要的那只手（WebView2 COM 事件回调）长在票 33 的窗口上**，而票 33 的 Out of scope 又明写 bridge protocol 归 35。

⇒ 断口的**规格根因**不是"有人偷懒"，是**责任矩阵里少一行**：入向接收器介于两票之间、两票都没逐字认领。这是要给 owner/编排者定的一句话（本程不替它定，也不改任何票面契约）。

### 2.5 在册待定项：`C24 GojaHostAPI` 初始集与 `C17` 方法白名单**今天定到哪一步**

- **C17 白名单：两份名单都存在，且交集为 0。**
  - 规格侧：SPEC-08 §5.2 那张表列了 12 枚 invoke 方法名 + 5 枚事件推送名（`:163-174`），节标题自带 `【SPEC 提案，S5 定稿走契约批准】` ⇒ **未批准**。
  - 代码侧：`internal/panel/bridge.go:42-45` 是一枚**闭集四常量**（现跑）：`MethodModeRequest = "panel.mode.request"`、`MethodWorkspaceRequest = "panel.workspace.request"`、`MethodAttachmentAdd = "panel.attachment.add"`、`MethodMessageSend = "panel.message.send"`，另有 `:30` `ComposerRequestSource = "panel-composer"`。
  - 交叉实测（本程逐枚跑 `grep -c` 于规格、`git grep -l` 于非测试产码）：**§5.2 那 12 枚名字在非测试 Go 代码里 0 命中；Go 这 4 枚名字在 `docs/specs/**` 里 0 命中**（只出现在工单与 `docs/evidence/s1/**` 里）。⇒ 白名单**不是"快定稿了"，是两条平行线**；定稿=选一条或合并，属契约面，须人工批准。
- **C24 `GojaHostAPI`：初始集一个名字都还没定。** 现量：`git grep -ln "GojaHostAPI\|host_api"` 在产码里只命中 `internal/config/manager.go`、`internal/config/schema.go`（配置字段，不是宿主 API 初始集），其余命中全在 `docs/PLAN.md` 与 `docs/specs/**`。⇒ 与 `AGENTS §2` 那条待定项一致（S7 切片卡），**本单只报状态，不定稿**。

---

## 3. Q2 依赖面与可复用地基

### 3.1 要引什么依赖（现量）

| 尺 | 读数 | 说明 |
|---|---|---|
| `grep -c webview go.mod` | **0** | 根模块（`go.mod:1` `module github.com/CarlosShao/wisp`）直依赖只有 5 枚：`sherpa-onnx-go v1.13.8`（`:8`）· `go-toml/v2 v2.2.4`（`:9`）· `x/crypto v0.57.0`（`:10`）· `x/sys v0.48.0`（`:11`）· `x-term v0.46.0`（`:12`），`:16` 起是 indirect 块 |
| 非测试产码里 `WebMessageReceived\|AddScriptMessageHandler\|PostWebMessage\|CreateCoreWebView2` | **0 命中**（排除 `scripts/spike`） | 复跑 181-c1 那件事：入向接收器整层不存在 |
| `github.com/jchv/**` 在 `go.mod`／`go.sum`／`deps.toml` | **三处全 0 命中** | 引依赖＝新动作，不是补登记 |

**要新增的确切坐标**（本程从仓里已有的独立模块现抄，不靠记忆）：`scripts/spike/go.mod:9` = `github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808`，`:18` 带一枚 indirect `github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1`。

⚠ **本程最重要的一枚 Q2 现量**：**这枚依赖今天已经在仓里，只是不在根模块。**
`scripts/spike/go.mod:1` 是独立模块 `github.com/CarlosShao/wisp/scripts/spike`（`go 1.27`／`toolchain go1.27.1`），
它 require 了 go-webview2，而且**真编译过**——`scripts/spike/bin/webview2-latency.exe` 就在目录里（与 `goja-caps.exe`、`xy-verdict.exe` 同批）。
⇒ "引这枚依赖能不能在本机解析并链接"这个问题**已经有人在 S0 回答过**，写腿不必再探。
版本斜差：spike 钉 `x/sys v0.42.0`、根模块 `v0.48.0`（根更新）⇒ 并进根模块**不需要降级**任何东西。

### 3.2 可复用地基（逐枚带尺，全部 `grep -n "^func "` 现量）

| 地基 | 今天提供什么（尺） | 对这一跳的可复用度 |
|---|---|---|
| `internal/panel/assets.go` | 10 枚 func：`BuiltinAssets :43`、`Resolve :75`、`Manifest :100`、`Check :137`、`Built :65`、`contentTypeOf :165`；`EntryFile :29 = "index.html"`（注释 `:28` 逐字：`// EntryFile is what the host loads for the panel's root URL.`） | **直接可用，零改动**：出向资源那一半的**被调物齐备**，缺的只是"注册 `AddWebResourceRequestedFilter` 的那只手" |
| `internal/panel/bridge.go` | 5 枚 func：`NewRequestID :58`、`ParseComposerRequest :84`、`knownComposerMethod :104`（闭集 `:42-45` 四常量 + 来源 `:30`）、`RefusedEnvelopeForUser :115` | **直接可用**：入向封套/白名单/拒答全在，**没有派发**（`:104` 只答"认不认得"） |
| `internal/panel/composer_handlers.go` | 7 枚符号：`ModeWriter :73`、`ModeConfirm :83`、`ModeWriteHandler :86`、`modeIsWidening :107`、`HandleModeRequest :111`、`actor :148`、`record :158` | **已存在、生产听众 0**（票 114 AC#2 的原生侧门） |
| `internal/panel/workspace.go` | 5 枚：`PathScope :37`、`AuditFunc :50`、`UnsetWorkspaceView :55`、`WorkspaceViewFromRoot :60`、`RequestWorkspaceSwitch :76`、`currentAfter :117` | 工作区切换的原生腿**已写好**（含 reparse/审计），等的只是入向调用方＝票 186 的活 |
| `internal/panel/pump.go` | 13 枚 func（`Snapshot :150`、`Marshal :204`、`Publish :215`、`StreamLog` 一族） | **出向管子今天真在跑**（`cmd/wisp/panel_pump.go` 10 枚 func 是其生产消费者）⇒ 入向只补另一半，**不许重画** |
| `internal/ball/sta_windows.go` | `//go:build windows :1`；10 枚符号：`staThread :24`、`newSTAThread :38`、`start :47`、`releaseCOM :103`、`waitStarted :111`、`threadID :120`、`PostTask :127`、`runTask :146`、`quit :156` | ⚠ **形状对、门不通**：类型与构造函数**全未导出** ⇒ 包外无法把 func 投到 `ui-sta`。D38a 要 WebView2 起在共享 `ui-sta` 上 ⇒ 写腿要么在 `internal/ball` 开一枚新导出（**碰球地界，须先报备**），要么 panel 自己再起一枚 STA 线程（**与 D38a 的"共享"相悖，属要人拍的一枚**） |
| `internal/proc/**` | 9 枚非测试文件、**43 枚 func**：`jobscope_windows.go` 11（`OpenJobScope :71`、`Assign :89`）、`treemetrics_windows.go` 7、`envfork.go` 9、`boot_windows.go` 5、`singleinstance_windows.go` 5、`externalsampler_windows.go` 3、`shutdown.go` 2、`systemprocs_windows.go` 1、`doc.go` 0 | Job Object 已可用（WebView2 的 3–5 枚子进程必须入 Job），⚠ 但 `Assign` 收 `*os.Process`，**COM 内部拉起的 `msedgewebview2` 子进程怎么入 Job 今天没有任何既有做法**（沿用 33-35 预检 U2，本程仍未验证——无宿主代码可跑） |
| `internal/winsec/**` | 5 枚非测试文件、**69 枚 func**：`winsec_windows.go` 25、`resolve.go` 22、`winsec.go` 13、`winsec_other.go` 8、`placement_windows.go` 1 | C26 `PathResolver` 地界；面板侧已经通过 `workspace.go:37 PathScope` 接口接上 ⇒ **入向接线不需要新写路径安全**（也不许在 `PathResolver` 之外用 `filepath.Clean\|Abs`，那是 `AGENTS §1.2` 的禁止项） |
| `scripts/spike/webview2-latency/main.go` | 327 行 / 13 枚 func；**已经跑通 JS→Go 回调**：`bindPing :118` 里 `w.Bind("spikePing", func(){…})` `:119`，浏览器侧由 `w.Eval("typeof spikePing === 'function' && spikePing()")` `:134` 触发；`ShowWindow` 掩/显 `:208-214`；`makeWebview :148` 用 `webview2.NewWithOptions` `:149` | ⚠ **一枚必须抄进新宿主的陷阱**，注释 `:113-115` 逐字：`// NOTE: w.Dispatch is NOT usable with a manual message pump - it posts a` / `// thread message that only go-webview2's own Run() loop interprets. Bind +` / `// Eval ride on window messages, which PeekMessage pumping does process.` ⇒ 常驻腿自己泵消息时**不能用 `Dispatch`** |

### 3.3 落到"新写几枚／复用哪几枚"

- **零改动即可复用（9 枚）**：`bridge.go` · `composer.go` · `composer_handlers.go` · `workspace.go` · `attachments.go` · `assets.go` · `approval.go` · `pump.go` · `cmd/wisp/panel_pump.go`。
- **新写（4 处，按地界分属两票）**：
  1. `internal/panel/host_windows.go`（票 33）——环境/控制器创建、`AddWebResourceRequestedFilter`→`assets.go:75 Resolve`、**消息接收回调→把 raw 文本交出去**、显示/隐藏/销毁、焦点归还、创建失败→`panel.unavailable`＋原生 L2 flag；
  2. `internal/panel/composer_dispatch.go`（票 35 Go 半边）——`raw → ParseComposerRequest → 按 `bridge.go:42-45` 四方法路由 → 拒答走 `RefusedEnvelopeForUser`；
  3. `cmd/wisp` 侧宿主生命周期挂载（`run.go`／`resident_windows.go`／`main.go`／`cmd/balldebug/main.go:199` 那枚现成桩句 `"tray: open panel (stub, ticket 33)"`）；
  4. 依赖与许可登记：根 `go.mod`＋`go.sum`＋`deps.toml` 三处（今天三处都还没有）。
- **可能连带（需先报备，本程不推荐顺手做）**：`internal/ball` 开 STA 投递口，或 `scripts/portable-tests.sh` 的 win scope 加 `./internal/panel/`（后者是**改门禁形状**，要人拍）。

---

## 4. Q3 最小闭环逐环表（网页 `postMessage` → `HandleModeRequest` 真被调用）

环号用 H（本程自己的编号，与 `docs/evidence/s1/33-35-preflight.md` 的 S 编号**不是一套**，映射见末行）。
"⚠ 真宿主？"= 这一环能不能在没有 WebView2 的前提下被建出**并可验**。

| # | 这一环必须存在的东西 | 今天有什么（现跑行号） | 缺什么 | 归属 | ⚠ 真宿主？ |
|---|---|---|---|---|---|
| H1 | 页面把请求封成一段 JSON 文本发出来 | Go 侧看得见的契约面就是封套结构体：`bridge.go:69-79`（`Method json:"method"` `:70`、`RequestID json:"requestId"` `:71`、`Source json:"source"` `:72`、`To/Path/Text :74-76`、`Attachments :78`）；来源常量 `:30 = "panel-composer"` | （前端侧的发送腿**属禁区，本程不读不引**，见 §11） | 票 34/36–40 | — |
| H2 | WebView2 控件被创建（Environment + Controller + 导航到虚拟 host），且起在 `ui-sta` 上 | **零**：`grep -c webview go.mod` = 0；非测试代码里 `CreateCoreWebView2` = 0 命中 | 全新宿主文件（建议 `internal/panel/host_windows.go`）＋ STA 投递口（§3.2 那格冲突） | **票 33** | **是** |
| H3 | **消息接收**：把网页发来的那段 raw 文本取进 Go | **零**：`WebMessageReceived\|AddScriptMessageHandler\|PostWebMessage` 非测试（排除 spike）= 0 命中 | 接收回调本体。**手段已在 spike 里跑通**：`scripts/spike/webview2-latency/main.go:119` 的 `w.Bind("spikePing", func(){…})` | **票 33 手段 / 票 35 语义，两票票面都没逐字点名（§2.4）** | 建要真宿主，**验不要**（见 §7 Q6） |
| H4 | `raw string → ParseComposerRequest` | 被调物 `bridge.go:84`；**生产调用者 0**（本程复跑：非测试命中只有注释与定义本身） | "接 raw 的那只手"——一枚 `func (d Dispatcher) Handle(raw string) string` 形的入口 | 票 35 | **否** |
| H5 | 方法 → 处理器**派发** | 只有"认不认得"：`knownComposerMethod :104`（`switch` 命中 `:42-45` 四常量则 true，否则拒）；`RefusedEnvelopeForUser :115` 给拒答文案 | 路由表本身（四方法分流）。本程现量：`git grep -nE "func Dispatch\|type Dispatcher\|func Route"` = **0 命中**；`func Serve*` 命中 5 枚但**全在 `internal/llm/adaptertest/harness.go:344,363,385` ＋ `internal/llm/golden/replay.go:129,139`**（LLM 测试夹具，与面板无关）⇒ 面板侧确实没有派发器 | 票 35 | **否** |
| H6 | `panel.mode.request` 的原生处理器（票 114 AC#2 的门） | **已存在**：`composer_handlers.go:111 HandleModeRequest`（变宽判定 `:107 modeIsWidening`，**写腿调用点 `:141 h.Modes.Set(ctx, to, PanelModeOrigin, h.actor())`——全仓这一族唯一的一处**，审计 `:158 record`）；装配齐 `cmd/wisp/run.go:417-423`（`Modes: modeStore` / `Confirm:` 与 Store 同一条腿 / `Audit: rt.auditf`），字段 `:233` | **只缺 H4/H5**——处理器与真写腿都就位，没人按下门铃；`HandleModeRequest` 非测试调用者 0（本程复跑 `git grep -n "HandleModeRequest" -- '*.go' ':!*_test.go'`：命中只有定义与注释） | 票 114（已交） | **否** |
| H7 | `panel.workspace.request` 的处理器 | 原生腿**已写好**：`workspace.go:76 RequestWorkspaceSwitch`（配 `PathScope :37`、`AuditFunc :50`、`currentAfter :117`） | handler＋路由（票 186 的 AC#2 正在选方法名） | 票 186 | **否** |
| H8 | `panel.attachment.add` 的处理器 | 守卫齐（本程现量）：`attachments.go` 的 `AttachmentBroker :169`、`NewAttachmentBroker :179`、`Ingest :199`、`AcceptedMIMETypes :153`、`ArtifactSink :59`、`NameGuard :165` | handler＋路由 | 票 92 后续 | **否** |
| H9 | `panel.message.send` 的处理器 | **白名单里有名字（`bridge.go:45`、`:106`），处理器零**：本程 `git grep -n "MethodMessageSend" -- '*.go'` ⇒ 非测试命中只有那两行白名单，其余 7 枚**全在测试文件**里 | handler 本体 | 票 35/92 | **否** |
| H10 | 回执回灌 Go→页面 | 出向**数据侧**在跑：`pump.go:150 Snapshot`／`:204 Marshal`／`:215 Publish`，消费者 `cmd/wisp/panel_pump.go:241 publishPanelSnapshot`（听众 `run.go:449`、`:597`），**落点是日志账**（`panel_pump.go:159-164 bookPanelSnapshot` → `logger().Info(panelSnapshotSummary(...))`） | 把同一包送去页面的**手段**（`Eval`/消息投递）＋ `RefusedEnvelopeForUser` 的第一枚生产调用者 | 票 33 手段 + 票 35 泵 | 手段要、**泵不要** |

**一句话形状**：H4–H9 **六环里唯一的公共缺口是"有人把 raw 交进来并按方法分流"**，它们一环都不需要 WebView2；
H2/H3/H10 的手段那三枚才需要。⇒ 这就是 §7 那张"能不能先动一半"的答案来源。
**与 preflight 的编号映射**：H1≈S1/S2 · H2≈S3 · H3≈S5 · H4≈S6 · H5≈S7 · H6≈S8/S9/S11 · H7/H8/H9≈S14 · H10≈S12/S13。

---

## 5. Q4 冷拉起代价：`AGENTS §2` 那条待定项今天能给什么现量

待定项原文（`docs/PLAN.md:1568` P11 行内，逐字）：

> `**若实测冷拉起 >2s，必须重新评估 L2 确认卡是否改回原生**（会牺牲 D29 的视觉质量）。这个决定不能在 S5 才发现要做`

### 5.1 能给：这枚现量**已经在仓里，而且是真宿主量出来的**

出处 `docs/evidence/s0/02-spike-report.md`（`:3` 逐字：`> 完成日期：2026-09-19 · 执行代理：T02-impl`），§3.4「WebView2（JSON 06，两轮）」`:144-152`，逐字抄：

> `| 量 | run1 | run2 |`
> `| cold（12 子进程真冷）P50 / P95 / max | 879.7 / 1041.6 / 1116.8 ms | 1125.7 / 1256.4 /（n=12） |`
> `| 全进程内首次 create（run 内） | 1808.9 ms（紧跟 12 子进程后，竞争态） | 1163.7 ms |`
> `| hot show P50 / P95 | 25.7 / 49.5 ms | 71.4 / 79.5 ms |`
> `| hot 浏览器往返 p50 | 3.5 ms | — |`
> `| recreate P50 / P95 | 955.7 / 1221.3 ms | 859.0 / 955.5 ms |`

时序定义同一份报告 `:82-84`：cold = `NewWithOptions`（建窗+show+`Embed` 阻塞至…）、hot = `ShowWindow(SW_HIDE→SW_SHOW)` + 一次泵。

**判读（只按上面这些读数说话）**：
- **2s 触发线：两 run 都没越过**（最大单枚 = 1808.9 ms，且它标的是"紧跟 12 子进程后的竞争态"，不是 cold 主体）。⇒ 本单**不能**给出"必须重评 L2 卡回原生"的证据，也**不能**给出"永远不必重评"——两 run 都在同一台机、同一枚 spike，距今 9 天且未复跑。
- **D32 的冷 ≤1500ms**：cold 主体 P95 为 1041.6 / 1256.4 ms ⇒ **两 run 均在预算内**；但那枚"进程内首次 create（竞争态）"1808.9 ms 已越过 1500ms ⇒ 预算达标**依赖机器不忙**，这一格到票 33 落地时仍要按 `A103` 的安静性三连重测（hot show 两 run P50 差 2.8 倍＝25.7 vs 71.4，本机噪声水平就在这）。

### 5.2 本单**唯一一枚可算"入向这一跳自己的"代理指标**

`hot 浏览器往返 p50 = 3.5 ms`（同一份表 `:151`）——它的量法就是 **JS 调 Go 绑定再回调**那一趟（`scripts/spike/webview2-latency/main.go:118-134`：`w.Bind("spikePing", …)` 由 `w.Eval("typeof spikePing === 'function' && spikePing()")` 触发，`browserRT` 计时间差）。
⇒ **H3 这一跳的传输成本量级 = 毫秒个位，相对 hot 预算 200ms 可忽略**；入向建起来不会因为"传输慢"而威胁任何预算。

### 5.3 没有的部分要说清

除上面两枚，**冷拉起的其它未知项都要"落地腿"才能测**：Runtime 缺失/过旧的失败形状（本程无宿主代码可跑，沿用 preflight U4/U2 的"未验证"）、`msedgewebview2` 子进程入 Job 的方式、CSP 响应头与 `<meta>` 同时存在时取哪条。**本程一行未测。**
⚠ 台账现量：`grep -n "P11" docs/reports/pending-and-issues.md` = **0 命中** ⇒ P11 只挂在 `AGENTS §2`／`SPEC-12 §4.2` 那张待定项表上，**真相源台账里没有对应 `A##` 条目**；这枚现量上面的 §5.1 若要用来自偿这条待定项，得由编排者入册（本程不代拍）。

---

## 6. Q5 owner 要的那四枚"能点的按钮"各自等这一跳的哪一部分

| 那一枚 | 工单 | 它等 §4 里的哪一环 | 建完最小入向（H4/H5）之后**能不能立刻动** | 还缺的、不在这一跳里的东西 |
|---|---|---|---|---|
| **切工作树／分支** | 票 186（`.scratch/wisp/issues/186-...md`，AC#1–#9 全未勾） | **H4+H5**（把 raw 交进来并按 `panel.workspace.request` 分流）＋ **H7**（handler 本体；原生腿 `workspace.go:76` 已写好） | **能动 Go 侧那一半**：请求→门→`RequestWorkspaceSwitch`→审计，全链在假宿主下可建可验；`perm.Store.Set` 型的生产调用者由 0 抬成 ≥1 | ① 快照**读面**（票 181：Go 侧今天零 git 读面）② 要不要跑外部 `git` 二进制（票 186 AC#1 自陈"这是一枚新的宿主能力，今天全仓零次"，归不归 `allowlist.txt`/`d22scan` 要说清）③ **`AC#2` 选方法名**若结论是"新增白名单条目"＝契约面，要 owner 落 `A##` ④ 真机那一格（H2/H3/H10）⑤ 界面控件＝`frontend/**`，另一会话 |
| **改模型／改思考档位** | 票 187（`:21`–`:29` AC#1–#9 全未勾） | **H4+H5** ＋ **一枚全新的 handler**（白名单 `bridge.go:42-45` 里**没有** `panel.model.request`／`panel.thinking.request` 这两个名字）——同 H9 那一形："有名无实现"再往前一步，是**连名都还没有** | **半能**：路由骨架、写腿、生效判据（AC#1 的读者普查）都能在假宿主下做；**但方法名要先入档** | ① **C17 白名单新增条目＝契约变更**（票 187 `:22` AC#2 自己写了"先具名报我、我落 `A##` 再写"）② D36 生效级别定档（热改/重载/重启，与这一跳无关）③ AC#3 那枚"改模型绝不可能顺带放宽权限档位"的反向判据 ④ 清单要真（AC#4，来自 `llm.providers.*.models.*` 与 `ThinkingLevels`）⑤ 真机 + 界面控件 |
| **面板里的档位与工作区＝权限输入口** | 票 92（`accepted-done`）＋ R20 口径 | **不等**——它给的是这一跳**必须遵守的边界**，不是要这一跳才能做的事 | **今天就是终态**（口径已在票 114 `:50` AC#5 落成格："面板侧**只许显示 + 发起请求**（R20 的明写）。工作区/档位两个输入口不得变成授权口"） | 唯一还欠它的是**可见性证据**：票 114 `:52` AC#6 要"真机**差分截屏**（改档位前后各一张）"⇒ 那一格要 **H2/H3/H10 真宿主**，假宿主给不出截屏 |
| **原生侧门已立、入向没接** | 票 114（`in-progress`；AC#1 `:34`、AC#2 `:38` 已勾 `[x]`，AC#3–AC#7 未勾） | **正是 H3→H4→H5 这一段**：门（`composer_handlers.go:111`）与写腿（`:141`）在等有人把 raw 递进来 | **能动到剩真机格**：AC#4（门二扩展名覆盖面）与入向无关、今天就能做；AC#5 在 H4/H5 落地当天可钉；AC#3 的"真机上变宽必须走 C18 原生卡、超时 300s 判拒"与 AC#6 的差分截屏**必须等真宿主** | `panel.unavailable` 事件的 sink 归口（预检 ④.3 第 1 条，写手不该自选）；票 37（原生降级卡的消费方）还没做 ⇒ flag 有置位者、零读取者 |

**四枚合起来的一句话**：**H4+H5 这一枚"最小入向"同时解开票 186 的 Go 侧、票 187 的骨架、票 114 的 AC#5**；
唯一被它**解不开**的是"肉眼在真机上看得到按钮生效"那一族格（票 92 的可见性、票 114 AC#3/AC#6、票 33 的 6 格真机 AC）——那些必须等 H2/H3/H10 的真宿主。

---

## 7. Q6 拆不拆得开：**能**——先建最小入向、不引 WebView2 依赖（结论 + 落点 + 判据形状）

### 7.1 结论

**存在这条路，而且它比"先做真宿主"在四枚现量上都更硬。** 一句话：**入向这一跳的"语义"和"手段"是可分的**——
`ParseComposerRequest`（`bridge.go:84`）的输入就是**一枚 `string`**（`:84` 签名现量：`func ParseComposerRequest(raw string) (ComposerRequest, error)`），
它不 import 任何 WebView2 符号；H4–H9 全链同理。⇒ **用一枚"假宿主/测试宿主"把 raw 文本喂进去，router 与写腿就能真跑、真红真绿。**

### 7.2 四条支撑现量

1. **依赖与语义无耦合**：`grep -c webview go.mod` = 0，而 `internal/panel` 的非测试产码里 `WebMessageReceived|CreateCoreWebView2` 命中 0（§3.1 复跑）⇒ 今天的 H4–H9 六环**一环都不引用 WebView2 符号**。
2. **写腿调用点已经写好在等**：`composer_handlers.go:141` 那一枚 `h.Modes.Set(...)` 是全仓唯一一处（本程 `git grep -n "\.Set(" -- cmd/** internal/** ':!*_test.go'` 里 mode 相关只命中它）⇒ 只要 H4/H5 有生产听众，`perm.Store.Set` 的生产调用者当天由 0 变 ≥1。
3. **CI 分母正好相反**（这条最实用）：`internal/panel` 在 **ubuntu/core scope**（预检 ③.3 现量：`scripts/portable-tests.sh:179` 的 glob、`core_pin :141`、CI 步 `ci.yml:266`/run `:288`），而带 `//go:build windows` 的宿主文件**在 windows scope 零用例分母**（`win_pin` `:151-160` 不含 panel）⇒ **最小入向那半有分母、真宿主那半今天没分母**（要它就得先改门禁形状，那是要人拍的）。
4. **依赖字节已在仓内**（§3.1）：就算要真宿主，`scripts/spike/go.mod:9` 已经钉过同一枚坐标并真编译过 `scripts/spike/bin/webview2-latency.exe` ⇒ 拆成两步不会让引依赖这件事变难。

### 7.3 落点（建议形状，不是派单）

- **新建** `internal/panel/composer_dispatch.go`：`raw string → ParseComposerRequest → 按 `bridge.go:42-45` 四方法分流 → 处理器 → 返回给用户的那句话`（拒答走 `RefusedEnvelopeForUser :115`）。全平台中性。
- **给它一枚生产听众，且必须具名是"诊断入口"不是"面板"**（这条防的是禁止清单 1.3"用 mock 代替真的来假报完成"）：形状沿用预检 §⑤ 的推荐——在 `cmd/wisp` 的 CLI 面上开一枚诊断子命令（同族先例 `cmd/wisp/panel_assets.go` 就是 `wisp panel-assets` 那条诊断腿），把 stdin/参数里那段同形 JSON 交给 dispatcher。票面与交付说明里要**逐字写明"这一格结的是 router+写腿，结不了'面板真的能点'"**。
- **不动 `bridge.go`**（白名单 `:104` 已是闭集；`bridge.go:5-9` 是明写的规则，别再开第二条 pipe）；**不动 `run.go`** 的 `rt.modeWrites` 装配（`:417-423` 原样留给宿主）。

### 7.4 判据形状（三发"拿掉就红"，缺一不算结）

1. 摘掉路由到 mode 处理器那一行 ⇒ "ModeWriter 被调用一次"那枚断言**必须红**。⚠ **看被调次数，不看日志字符串**（预检 §⑤ 同一条）。
2. 把 `knownComposerMethod`（`bridge.go:104`）的默认分支从"拒"改成"放行" ⇒ 白名单外方法名那枚用例**必须红**。
3. 把假宿主喂入的封套里 `source` 从 `"panel-composer"`（`:30`）改成别的 ⇒ 必须走 `RefusedEnvelopeForUser` 的拒答路径**且不是静默丢弃**（`bridge.go:93-96` 那支就是为这写的）。
4. ⚠ 反面对手（本程加的一条，防"假宿主被当成真面板"）：交付里必须有一枚具名断言说**真窗口那六格（票 33 AC `:36-42`）一格都没结**；否则按"能力装上就以为通路通了"那族病判退回。

### 7.5 这条路**结不了**的（要说清，别让它被当成地基完工）

真窗口生命周期/时序/焦点/netstat 无监听/runtime-missing fixture（票 33 的 6 格 AC）、`panel.unavailable` 事件、票 114 AC#3/AC#6 的真机与截屏、票 92 的可见性证据、以及票 33 **AC#7/AC#8** 那两格（要真有人消费落盘摘要/读全量字节才有 hook）。⇒ **票 186/187 不是"整体等地基"，是"能先动 Go 侧那一半"。**

---

## 8. 推荐拆法与代价（**不拍板**，只摊清）

先摆这一条约束（票 92/R20 口径的逐字原文，`docs/reports/pending-and-issues.md:1263`）：

> `- **A84⑥ 对 owner 的口径**：R20（三档 / 默认最严 / 档位持久化 / 只有全自动要 L2 / 显示在主面板输入框 / 不换 git 分支）`

三片，顺序即建议顺序：

| 片 | 内容 | 新写 | 引依赖 | CI 分母 | 要人先批的枚数 | 结不了什么 |
|---|---|---|---|---|---|---|
| **A 最小入向**（建议先做） | `internal/panel/composer_dispatch.go`（H4+H5）＋ 一枚**具名诊断入口**当生产听众（§7.3） | 1 产码 + 1 同包测试 | **0**（`go.mod` 一字不动） | **有**：`internal/panel` 在 ubuntu core scope（预检 ③.3） | **0** | 票 33 的真机 6 格、票 114 AC#3/#6、票 92 的可见性、AC#7/AC#8 |
| **B 真宿主**（票 33 本体） | `internal/panel/host_windows.go`：环境/控制器/资源回调（接 `assets.go:75`）/**H3 消息接收**/显示隐藏销毁/焦点归还/创建失败→`panel.unavailable` | 1–2 产码（带 `//go:build windows`）+ 用例 | **1 枚新依赖 + 1 indirect**（`jchv/go-webview2` + `go-winloader`）+ `go.sum` + `deps.toml` 许可登记 | **今天没有**（`win_pin` 不含 panel） | **2**：① 要不要把 `./internal/panel/` 加进 win scope（改门禁形状）② STA 归属（复用 `ui-sta` 要 `internal/ball` 开导出 vs panel 自起第二枚 STA，与 D38a 的"共享"相悖） | `msedgewebview2` 子进程入 Job 的方法（U2 未验证，可能要再一发 spike）；CSP 头 vs `<meta>` 冲突（U4）；票 37 消费方未做 ⇒ flag 零读取者 |
| **C 回灌**（H10） | `Eval`/消息投递 + `RefusedEnvelopeForUser` 的第一枚生产调用者 | 1 处改宿主 | B 已引 | 同 B | 同 B | 前端接收端（`frontend/**`，本程禁区，另一会话） |

**代价的诚实形状**：片 A 是"2 枚新文件、零依赖、有分母、零批准"；片 B 是"1–2 枚新文件、1 枚新依赖、**今天零分母、2 枚要先批**"。
把 A 与 B 分成两片的代价 = **H3 那枚接收器要在 B 落地后改一行改成调 A 的 router**（不是重写）；
不拆的代价 = 为了拿到 router 的绿，先替 owner 批下"改门禁 + 开球的地界导出"这两枚他还没点头的东西。

---

## 9. 本程没测什么（逐名，不圆）

| # | 没测/没核的条目 | 为什么 |
|---|---|---|
| N1 | **`frontend/**` 与 `design/**` 一行未读** | 派单禁区（本程工具调用里对它零次读取）⇒ §4 的 H1 只从 Go 侧的 `ComposerRequest` json tag 反推形状，**不代表前端真发的是那个形状**；那一半要由前端会话对账（既有对账门：`bridge_test.go`/`composer_test.go` 会读 `frontend/` 源码，本程没跑它们） |
| N2 | **`publishPanelSnapshot` 是否经 `pump.Publish()`→`Marshal()`** | 本程只现量到 `panel_pump.go:241` 的定义与 `:232` 的注释逐字 `// publishPanelSnapshot puts one snapshot on the wire. Called where a state the`，**没读函数体** ⇒ **票 33 AC#7 的开窗条件（"`Marshal()` 的第一枚生产调用者"）今天算不算已满足，本程不判**，交回编排者 |
| N3 | 本机 WebView2 Runtime 现在装没装、`msedgewebview2` 子进程怎么入 Job、CSP 响应头与 `<meta>` 谁赢 | 无宿主代码可跑（§4 H2 不存在）；沿用预检 U2/U4 的"未验证" |
| N4 | `deps.toml` 是否被 `wisp doctor` 逐条比对 Go module | 沿用预检 U3，本程只现量到 `grep -rn "jchv\|webview" deps.toml` = 0 命中 |
| N5 | 票 33 AC#7/AC#8 那两格提到的分支覆盖（`panel_pump.go:186-192` 摘要超 440 字符那支、`run.go:738` `EvToolStart` 那支） | 本程**没跑任何覆盖率/变异**；`EvToolStart` 零发射者那条是票面转述，我未复量 |
| N6 | 全仓编译/`go vet`/任何 `go test` 之外的产码验证 | 只读程；§10 那三条门禁是本程唯一跑过的执行 |
| N7 | preflight 里我**没引用**的其它行号（如 `attachments.go:381`、`pump.go` 各行） | 只复核了本报告实际写进行号的那几枚；未复核的一律不引 |
| N8 | 门禁 baseline 之外的包、SLO 台件、`scripts/portable-tests.sh` 本体 | 派单只点名三条门禁；⚠ **未跑 `probes/161/r6/flip-declaration.sh`（禁令）**，未改 `probes/154/gate-clauses.sh`（禁令） |

---

## 10. 门禁终态

（待填：三条现量 + 红腿名册 + `internal/panel` 那 2 枚已知红照实记；读数取在本程最后一枚内容 commit 之后）

---

## 11. 被拒调用＋零删除自证＋工具调用终值

（待填）

---

## 12. next：派写腿之前还缺什么、哪几枚要人先批准

**缺的东西（按能不能自己走分档）**：

1. **可立即派**（不需任何批准）＝ §7 那片"最小入向"：`internal/panel/composer_dispatch.go` + 具名诊断入口 + §7.4 三发判据；零依赖、有 CI 分母、不碰 `bridge.go` 也不碰 `run.go:417-423`。
2. **要人先批准（4 枚，逐枚具名）**：
   - **① C17 方法白名单定稿路线**——§2.5 那枚现量：规格的 12 枚名字与代码的 4 枚名字**交集为 0**。三支候选（合流 / 各管各层并写明 / 以代码那 4 枚为准补进 SPEC-08 §5.2）都要改契约面 ⇒ **owner 拍**。（票 186 AC#2 与票 187 AC#2 都挂它。）
   - **② STA 归属**——复用 `ui-sta` 要在 `internal/ball` 开新导出（碰球地界，`staThread :24`/`PostTask :127` 全未导出）；panel 自起一枚与 D38a"共享"相悖 ⇒ **owner 或球票写手拍**。
   - **③ CI 分母**——把 `./internal/panel/` 加进 win scope 是**改门禁形状**；不加就要在票面明写"该腿分母＝本机人工，CI 无" ⇒ **编排者/owner 拍**。
   - **④ 是否允许跑外部 `git` 二进制**（票 186 AC#1 那枚"新的宿主能力，今天全仓零次"）⇒ 关系到 `allowlist.txt`/`d22scan` 边界，**owner 拍**（本程不碰那两件）。
3. **编排者可自行补、不必 owner**：
   - **责任矩阵补一行**：§2.4 那枚"入向接收器 H3 两票都没逐字认领"——在票 33 或票 35 票面各加一行归属即可（票面不是契约，但**改动要留痕**）。
   - **`panel.unavailable` 的 sink 归口**（进 `Logf`／`internal/observe`／C17 事件流，三者听众不同，写手不该自选）。
   - **`frontend/embed.go` 归属确认**（`.go` 但在 `frontend/**` 下，预检按"按目录判"划入冻结）。
   - **P11 的台账条目**：`grep -n "P11" docs/reports/pending-and-issues.md` = 0 命中（§5.3）⇒ 若要拿 §5.1 的现量自偿那条待定项，需编排者先落 `A##`。
4. **本程对票 33 的唯一写动作**：追加一段 Progress log（§11 记其落点），**AC 框一格未勾**、原句一字未改。
