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

## 4. Q3 最小闭环逐环表

（待填：从 `postMessage` 到 `HandleModeRequest` 被真调用，逐环 = 环名 / 现在有什么 / 缺什么 / 最近邻指到现有行）

---

## 5. Q4 冷拉起代价

（待填：`AGENTS §2` 那条待定项今天能给什么现量）

---

## 6. Q5 与 owner 那四枚按钮的依赖关系表

（待填：票 186／187／92+R20／114 逐枚）

---

## 7. Q6 能不能不引 WebView2 先建最小入向

（待填：结论 + 落点 + 判据形状，或具名"必须真宿主"）

---

## 8. 推荐拆法与代价（不拍板）

（待填）

---

## 9. 本程没测什么（逐名）

（待填）

---

## 10. 门禁终态

（待填：`d22scan.sh` rc 与 `ban #8 internal/` examined、`gate-clauses.sh` 红腿名册、`go test ./internal/panel/ ./internal/config/`；2 枚已知红照实记）

---

## 11. 被拒调用＋零删除自证＋工具调用终值

（待填）

---

## 12. next：派写腿之前还缺什么、哪几枚要人先批准

（待填）
