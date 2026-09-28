# 调研 2026-09-28 腿 2：DeepSeek harness（DSH）`packages/client/ui-*` 52 枚分包逐块清点

> 只读调研程 `survey-dsh-ui`。参照仓：`D:\work\AI\open source\deepseek-harness`（**该克隆没有 `.git`** ⇒ 全文不引提交历史，每条"他们有"都是 `文件:行`）。
> 产出件只这一枚。未改那边一字、未建那边一文件。
> 本仓现量全部来自 `D:\work\workspace\projects plans\Wisp` 的 `internal/**`、`cmd/**`、`docs/**` grep/ls 真读数；
> **`frontend/**` 与 `design/**` 属界面那支地界，本程不读不转述** —— 纯界面项一律标〔未量·frontend 地界〕。

---

## 0. 口径与分母（先给推导式）

| 量 | 读数 | 推导式（现跑） |
|---|---|---|
| `ui-*` 分包 | **52** | `ls packages/client \| grep -c '^ui-'` = 52 |
| `packages/client` 条目总数 | 66 | `ls packages/client \| wc -l` = 66 |
| 非 `ui-*` 条目 | 14 | `ls packages/client \| grep -v '^ui-'` = `AGENTS.md README.i18n.yaml README.md README.zh.md connection file-upload hmr locale modules resources shortcuts store tsdown.client.ts web` |
| 交叉验算 | 66 − 14 = **52** | 两条独立口径同值 |
| 名册逐枚一行摘要的出处 | 每枚 `ui-<slug>/README.zh.md:2`（frontmatter description）与 `ui-<slug>/package.json:3`（英文 description） | 52 枚**两行号完全一致**（全部 `.zh.md:2`、全部 `package.json:3`）⇒ 名册不是一枚枚手抄，是同一口径抽的 |
| DSH 官方包地图的行数 | **50** 枚 ui-* | `grep -oE '\[`ui-[a-z-]+/`\]' packages/client/README.md \| sort -u \| wc -l` = 50，表体在 `packages/client/README.md:29-88` ⇒ **官方那张包地图漏了 2 枚**（见 §9-④） |
| DSH 会话事件词汇表 | **59** 枚 | `packages/core/session/src/known-event-types.ts:22-82`，条目体在 `:23-81`（`sed -n '23,81p' \| grep -c "^  '"` = 59） |
| DSH Host→浏览器转发事件 | **27** 枚（其中 **2** 枚是 waterfall） | `packages/api/remotes/src/remote-events.ts:20-48`（条目 `:21-47` = 27）；waterfall = `approval/request`(`:22`) 与 `user-questions/request`(`:47`) |
| 本程深读枚数 | 8（派单点名的 8 枚）＋顺手读了 `ui-jobs`、`ui-conversation` 的 README 概述段 | 其余 42 枚只登记（名册在 §1，逐枚标"未深读"） |

派单原话"我 19:5x 现列 52 枚"——**复核结果：52 枚没错，名册逐名也对得上**（我没有在派单名册里找到抄错的枚名）。

---

## 1. 52 枚名册：逐枚"它是干什么的"

格式：包名 ｜ 它是干什么的（摘自该包 `README.zh.md:2`，中文摘要）｜ 深读与否。
"它是干什么的"一列**是照抄该包 README frontmatter 的 description**，不是我概括的（口径统一，便于他抽查）。

| # | 包 | 它是干什么的 | 本程 |
|---|---|---|---|
| 1 | `ui-agent-preset` | 在 Web 中选择 Agent preset 和新任务默认值，查看各模式的说明与声明的配置；创建与修改引导至创造模式 | 未深读 |
| 2 | `ui-approval` | 通过作用域交互路径响应 Host 权限请求的浏览器批准界面 | **深读 §3.1** |
| 3 | `ui-attachment` | 对话 UI 的附件呈现：混合草稿附件栏、文档拖放目标、历史图片画廊与原图灯箱 | 未深读 |
| 4 | `ui-brand-official` | 面向侧栏的官方品牌填充，仅在官方构建中生效 | 未深读 |
| 5 | `ui-chat` | 渲染会话对话节点、历史图片、操作、本地化和滚动状态的浏览器 Chat target | 未深读 |
| 6 | `ui-commands` | Web GUI 的客户端命令 API：`/` 命令 source、三类派发、会话级命令目录、popupSelect 与 action 注册 | 未深读 |
| 7 | `ui-conversation` | Target-neutral 对话装配与浏览器 shell：事件和视图注册表、逐会话 binding、输入状态、slot 与临时 composer takeover | 顺手读概述段（§3.2 需要） |
| 8 | `ui-deliverables` | 改动文件卡片＋逐文件对比 review tab＋交付文件卡片＋收尾正文里可点击的行内文件引用 | 未深读 |
| 9 | `ui-directory-picker-browse` | 应用内目录浏览表面：Miller 分栏「选择工作区目录」对话框 | 未深读 |
| 10 | `ui-directory-picker-native` | 原生目录选择表面：驱动本地 Desktop 或 Host OS 选择器的浏览器半部 | 未深读 |
| 11 | `ui-dockkit` | 停靠布局套件：带**可逆操作**的标签格分裂树、planner、线性历史，以及渲染并驱动它的组件 | **深读 §3.5** |
| 12 | `ui-goal` | 目标界面：composer 上方常驻的 GoalBar，可编辑/暂停/恢复/清除 | 未深读 |
| 13 | `ui-input-trigger` | 输入触发流水线：光标处 `/` 与 `@` 检测、分组候选菜单、把 pick 路由到已注册 source | 未深读 |
| 14 | `ui-jobs` | 会话头部后台任务列表：可展开的流式输出面板、进行中/已结束分组、无保留输出的静态行 | 顺手读概述段（§3.9） |
| 15 | `ui-layout` | 外壳布局：三栏 AppFrame（右栏是贴边面板的轨道）、面板几何服务与主题呈现 | **深读 §3.6** |
| 16 | `ui-message-feedback` | 已定稿助手消息动作行里的 Like/Dislike 对＋两种评分与 `/feedback` 共用的反馈弹窗 | 未深读 |
| 17 | `ui-model-selection` | 模型选择：`/model` 弹窗与 composer 模型位共用一份按提供方分组的会话级目录 | 未深读 |
| 18 | `ui-open-in-app` | "Open In…" 控件：在记住的应用里打开 workspace 目录；文档预览里用默认应用打开/显示单个文件 | 未深读 |
| 19 | `ui-permission-presets` | 权限预设界面：通用设置里的默认行＋切换当前会话的 `/permission` 选择器 | 未深读 |
| 20 | `ui-plan` | plan 模式状态徽章：显示 plan 模式已开启并可将其关闭 | 未深读 |
| 21 | `ui-plugin-manager` | 从 Web 侧栏管理 profile 的插件组合包、它们的行与插件配置 | 未深读 |
| 22 | `ui-primitives` | 共享 React UI 原子件：控件、图标、Markdown 与数学公式渲染、终端/读取/差异/搜索/网页输出卡片（零 Cordis） | 未深读 |
| 23 | `ui-reference` | Web `@file` 与 `@session` 引用 source：候选、排序、原子行内引用 | 未深读 |
| 24 | `ui-renderer` | 浏览器 UI 渲染器：普通 Slot 与可复用 Component Factory 的 React 绑定、组装后的应用根 | 未深读 |
| 25 | `ui-schedule` | 跨会话的任务管理页面＋当前会话的提醒列表 | 未深读 |
| 26 | `ui-session` | 面向 Session Controller 列表、交互状态与逐会话上下文的 React 与 Slot 适配器 | 未深读 |
| 27 | `ui-settings` | 设置领域底座插件：共享配置表单、schema 服务、规范设置 slot 类型约定 | 未深读 |
| 28 | `ui-settings-account` | Desktop 设置里的账号页：DeepSeek 登录状态、浏览器授权与取消；侧栏账号菜单提供退出登录 | 未深读 |
| 29 | `ui-settings-agent-loop` | 插件页上的 Agent 循环设置页：并行工具调用上限 | 未深读 |
| 30 | `ui-settings-general` | 设置外壳、无归属文案与持久化产品引导命名空间：「通用」分区、触发控件框架、引导账本投影 | 未深读 |
| 31 | `ui-settings-models` | 模型设置与产品引导：提供商行、API 密钥管理、模型列表、首次运行弹窗 | 未深读 |
| 32 | `ui-settings-plugin-inventory` | 按作用域分组的**只读**插件清单标签页：preset 组合在前、全局平面收在折叠分组、搜索跨两组 | 未深读 |
| 33 | `ui-settings-plugins` | 「内置插件」设置分区：导航项＋供功能插件注册标签页的标签行 | 未深读 |
| 34 | `ui-settings-shell` | 终端执行器设置页：shell 命名空间的命令超时与单流输出上限 | 未深读 |
| 35 | `ui-settings-subagent` | 子智能体设置页：委派深度与并行容量＋Agent 可为子代理选的模型，同页一次保存 | **深读 §3.4** |
| 36 | `ui-settings-web-search` | DeepSeek 网页搜索提供方设置页：API Key、接口地址、单次请求搜索次数上限 | 未深读 |
| 37 | `ui-shortcuts` | 查看当前窗口可用命令，按操作名/英文别名/键位搜索；录制、清除、恢复快捷键 | **深读 §3.7** |
| 38 | `ui-sidebar` | 侧栏外壳插件：品牌行、New Session、折叠控件、可感知滚动的区域席位、底部固定 Settings 席位 | 未深读 |
| 39 | `ui-sidebar-browser` | 右侧 Sidebar 浏览器 tab：在 sandbox 中访问 HTTP(S) 页面，含 loopback 服务 | 未深读 |
| 40 | `ui-sidebar-documentpreview` | 右侧 Sidebar 文档预览：共享加载与控件，可选 Markdown/代码/图片/PDF/Office/HTML 渲染器，纯文本兜底 | 未深读 |
| 41 | `ui-sidebar-files` | 右侧 Sidebar 文件树 tab：逐层列出会话工作区根，按资源地址把文件打开到 Sidebar | 未深读 |
| 42 | `ui-sidebar-right` | 右侧 Sidebar：每会话一个停靠面、两种呈现形态、导航控制器、tab 类型注册表与 Tab 域 | 未深读（停靠面本体在 §3.6） |
| 43 | `ui-sidebar-terminal` | 在 Web 右侧栏打开、恢复和控制交互式 shell 标签页 | 未深读 |
| 44 | `ui-skill` | Web skill 引用与专属 skill 工具行：`/` 触发的 skill source 与 skill 调用卡片 | 未深读 |
| 45 | `ui-slots` | slot 注册表纯核心：普通扩展 slots、可复用 Component Factory、推导 props 类型、store 席位、渲染器安装约定 | 未深读 |
| 46 | `ui-subagent` | 子代理对话目录、续接路由 UI 与 `@` 引用 source | **深读 §3.2／§3.3（owner 那一问的答案）** |
| 47 | `ui-theme` | 主题与正文字号设置：`--dsw-*` token 样式表、ThemeRuntime 状态、「通用」设置行、插件前引导 | 未深读 |
| 48 | `ui-tool` | Client 工具展示插件：完整调用树的组合、按工具名键控的视图 slot、内置原子工具卡片 | 未深读 |
| 49 | `ui-trajectory` | Trajectory 视图：按轮次组织的事件记录表＋交互式时间概览，注册进对话视图环 | **深读 §3.3** |
| 50 | `ui-user-questions` | `ask_user_question` 功能：接管编辑器的提问 UI 与 plan-review 审批卡片 | **深读 §3.1b** |
| 51 | `ui-workflow-run` | 持久化工作流运行节点：把顶层工作流运行重建为带嵌套成员折叠的独立聊天节点 | 未深读 |
| 52 | `ui-workspace` | 共享 Workspace 浏览器与选择器：分组或扁平的会话行、管理操作、slot 组合的行 action、目录选择 | 未深读 |

**结构读到的一个事实（不是功能差集，但会影响抄法）**：这 52 枚里只有 `ui-primitives`、`ui-slots`、`ui-dockkit` 三枚在 description 里明写"零 Cordis／不注册服务"（`ui-primitives/README.zh.md:2`、`ui-dockkit/README.zh.md:14`），其余都靠"注册进某个 slot"活着。**界面在这家是一等公民的扩展点，不是一个渲染函数。**

---

## 2. 专答 owner 那一问（先给结论，再给数据形状）

> owner 09-28 18:0x 原话：「**子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面**，能明白吗？这是主流 harness 必做的，不要偷懒」

**答：是 `ui-subagent`，不是 `ui-trajectory`。但 `ui-trajectory` 是那一页的"第二种看法"，两者是同一份数据的两个视图，缺一都不算做完。**

逐条对上"看到状态／点进去／各自的流式工作页"这三段要求：

| owner 要求的那一段 | 在哪一枚 | 具体是什么（带行） |
|---|---|---|
| "必须看到状态" | `ui-subagent` | 会话**页头**一枚后代数量触发器打开的**目录树**，每行一个子代理：状态点（running=共享 ongoing loading 转圈、正常完成=success 绿点、其余=idle 灰点，`SubagentHeaderLineage.tsx:334`）＋ mode 标签＋活动标签＋标题＋**持久化 token 总量**＋**活跃轮次耗时**（`:250-270`）。行文案的组装式在 `:250-252`（`[title, mode, activity].join(' · ')`）与 `:268-270`（`[tokenMetric, durationExact].join(' · ')`） |
| "点击某个子代理" | `ui-subagent` | 两个出口：**点行 = 把主区导航进那枚子会话**（`:272-279` → `openChild({parentSessionId, childSessionId, mode})` → `ctx.uiWorkspace.openSession(address)`，`index.ts:61-63`）；**点行尾箭头 = 在右侧 Sidebar 另开一个 tab 装它**（`:280-285` → `openChildAside` → `ctx.sidebarRight.openResource(..., {preferNewPane: true})`，`index.ts:64-69`） |
| "能看到它们各自的流式工作页面" | `ui-subagent` | 那个 tab 的正文渲染的是**共享的 `conversation.content` Component Factory，view 钉成 Chat**（`sidebar-chat/index.tsx:145` 与 `:121-123`）——**不是一套"子代理专用小视图"，就是把主对话那套组件原样装进子会话**。所以子代理的 delta、工具调用、思考、图片，全部走和父会话一样的流式渲染 |
| 那一页还能换成"事件台账＋时间轴"的看法 | `ui-trajectory` | Trajectory 是注册进 `conversation.view` 环的一枚视图 tab（`ui-trajectory/src/client/index.ts:79-112`，`id: 'trajectory'`、`order: 10`），**它的 inject 是按 sessionId 取绑定的**（`:88-93`：`ctx.sessions.binding(sessionId)` / `ctx.uiConversation.binding(sessionId).target('trajectory')`）⇒ **导航进子会话之后，那枚子会话自己的 Trajectory 也在**（同一套视图环，逐会话一份数据）。但注意：**在右侧 Sidebar 那个 tab 里视图被钉死成 Chat，看不到 Trajectory**（`sidebar-chat/index.tsx:121-123` `FixedChatConversationView` 只渲染 `{ view: 'chat' }`） |

**一句话给 owner**：DSH 的"点进去那页"= 子代理被当成**一份真会话**（有它自己的 `childSessionId`），点进去就是把整台对话外壳搬到它身上；`ui-trajectory` 不是那一页，它是**"任何一份会话"的第二视图（逐事件台账＋可拖选的时间轴）**，父会话和子会话都能开。

---

## 3. 八枚深读：三段式（他们怎么做／我们现在是什么／落到我们身上是什么后果）

### 3.1 `ui-approval` —— 审批那一跳

**他们怎么做**

- 机制：Host 的 `approval/request` 事件以 **waterfall** 形态跨过进程边界，浏览器侧 `ctx.remote.$on('approval/request', function (request, next) { ... })`（`ui-approval/src/client/index.ts:104`）。转发名册逐字写着 `{ event: 'approval/request', mode: 'waterfall' }`（`api/remotes/src/remote-events.ts:22`）。**waterfall = 监听者可以 `next()` 把请求让给下一个应答方**；让出去了就不算答（`contract/slots.ts:144-157` 的 `delegate()` / `isDelegation()`）。
- 数据形状（客户端可见字段，`contract/slots.ts:52-63`）：`toolName` / `callId?` / `reason?` / `displayReason?: { en, [locale] }`（请求方自带的**本地化展示文案**）/ `signal?: AbortSignal`。
- 决定只有两枚：`type ApprovalDecision = 'allowed-once' | 'rejected'`（`contract/slots.ts:66`）。**README 明写这是面板的边界**："面板只提供临时决定——它支持仅本次允许和拒绝；**持久权限策略仍由 Host 侧审批包拥有**"（`README.zh.md:36`）。
- 一次性结算：`PendingApproval` 是一个持有 `Promise<ApprovalDecision>` 的 class（`:88`、`:109-113`），`#settled` 之后**再答一次直接 throw**（`:169`：`pending approval ... is already settled`）；已撤销/被替换的请求不能二次作答，较早请求失败也**不会解锁替代请求**（`README.zh.md:21`）。
- 键盘：详情区聚焦后 **Enter 批准、Esc 拒绝**；插件挂载期间这两个键**不许分配给可编辑快捷键**（`README.zh.md:21`）——界面抢键、快捷键让路、且被钉在快捷键目录里（`ui-shortcuts/README.zh.md:77`：审批和双 Esc 停止行只由已挂载的功能 owner 贡献，不可改）。
- 关联详情：可选注入一枚 `conversation.approval.detail` slot，按 `callId` 渲染"这次批准背后的那次工具调用"（`contract/slots.ts:37-49`）。

**我们现在是什么**

- 审批卡视图字段：**11 枚**（`internal/panel/approval.go:39-59`：`correlationId/tool/args/level/rulesHit/reason/reasonKnown/sessionOverrideBlocked/callChain/decidedBy`）——比他们多（我们带 `rulesHit`、`sessionOverrideBlocked`、`decidedBy`，这三枚他们没有）。
- 应答通道：**4 枚**，闭集（`internal/agent/approval/approval.go:53-61`：`ball`/`esc`/`panel`/`kws`，注释逐字写"the set is closed: a fifth value means a caller invented a selector"）。他们是"waterfall 应答方链"，我们是"多通道同一票"——形状不同，见 §9-⑧。
- `ask_user` / `ask_user_question` 关键词：**0 命中**（`grep -rni "askuser|ask_user|userquestion|user_question" --include=*.go internal cmd` = 0）。
- 面板能否"允许"：**只有拒绝**（`docs/specs/SPEC-08*.md:156-176` C17 名册逐字：`approval.decide` —— 「"allow"拒绝一切面板来源（F2）；仅 reject 可面板发起」）。DSH 的面板侧**能给 `allowed-once`**，但它把这件事放在"请求方是 Host、应答方是浏览器 Remote listener"这一形里，而我们禁的是"面板主动写 allow 到 Go"。**这两条不等价**，见 §9-⑧。

**后果**：我们缺的不是审批卡（我们卡更硬），缺的是**"允许一次"这个动词在界面上有没有出口**。今天界面上没有任何一枚出口能给出"就这一次"，`L2` 只有拒绝；这会导致 owner 第一天自用就撞上"要么永久放宽、要么每次都拒"的两难。

### 3.1b `ui-user-questions` —— agent 反过来问人那一跳

**他们怎么做**

- 跨进程形状：`user-questions/request` 同样是 **waterfall**（`api/remotes/src/remote-events.ts:47`），事件签名逐字写着"Return an answer to claim the request or call `next()` to delegate. Scope-filtered dispatch: agent-scoped listeners receive only that agent."（`packages/interaction/user-questions/src/types.ts:79-94`）⇒ **提问是按 Agent 作用域投递的**，多个界面（TUI/Web/移动）挂上来就是天然的应答方链，第一个认领者赢。
- 问题形状（`src/types.ts:36-51`）：`id` / `question` / `detail?` / `header?` / `options?: {label, description?}[]` / `multiSelect?` / `intent?`。
- **`intent` 这一枚很关键**（`src/types.ts:22-33`）：注释逐字写"an intent changes **presentation only**, never the protocol"。目前只有一种 intent：`plan-review`，且批准项是**按 label 点名**而不是按位置——"Named rather than positional so no UI infers the verdict from option order"（`:27-29`）。
- 答案形状（`src/types.ts:54-67`）：`answers: [{id, selected: string[], custom?}]`。**多选允许 `selected` 与 `custom` 同时携带**（`ui-user-questions/README.zh.md:32`），单选自定义互斥。
- 跳过／取消是两种不同的东西：跳过=给这一项发**既有的空 `{ selected: [] }`**；关闭=以 `ASK_CANCELLED` **拒绝整个等待**（`README.zh.md:32`，代码 `contract/slots.ts:106-111`、`:188-194`）。取消码两枚：`ASK_ABORTED`（宿主撤走）／`ASK_CANCELLED`（人关的）。
- 草稿持久性刻意的窄（`README.zh.md:40`、`:94`）：题号/已选/custom/显式跳过状态只放在**非持久化 slot 存储**里，键是"待处理请求的本地渲染标识"；**从不写 Host、`localStorage`、磁盘**。
- **每次只有一个请求拥有编辑器**；后续待处理请求留在会话快照里，等前一枚落定再显示（`README.zh.md:95`）。

**我们现在是什么**

- 0 命中（见 3.1）。`internal/agent/approval/ui.go:71-87` 有 `Prompt`/`Update` 两方法，但那是**审批**的展示 seam，不是"agent 问一个带选项的问题并等人回答"。
- `PLAN.md:1123`（账里已引过）自己写着子代理没有提问通道 ⇒ 票 197 判据"不许自带允许出口"。

**后果**：这一跳我们**整条都不存在**。它和 3.1 是同一枚缺口的两面：DSH 用**同一套 waterfall 机制**同时服务"要批准"和"要回答"，只是 payload 不同。我们要做的话，形状上应该复用 `internal/agent/approval` 已有的作用域路由＋correlationId，而不是新造一条通路。

### 3.2 `ui-subagent` —— 名册与那一页（数据形状，owner 硬要求）

**名册条目的数据形状**（服务端拥有、客户端只读，`packages/subagent/subagent/src/projection-types.ts:10-19`）：

```
SubagentCatalogEntry =
  { id: SessionId, createdAt: number }
  & ( { mode:'one-shot';   label?: string }   // 一次性，label 可缺
    | { mode:'continuable'; label: string }   // 可续，label 必在
    | { mode:'unknown';     label?: string } )// 模式未知仍可见可点
```

**计时条目的形状**（`projection-types.ts:22-37`）：

```
SubagentTimingProjection =
  { settledMs: number,                          // 已结束轮次累计毫秒
    active?: { since, through },                // 当前未结束轮次的同一切面起止
    lastTurnCompleted?: boolean }               // 最近一个已结束轮次是否 completed
```

投影注册三枚（`projection-types.ts:68-85`）：`subagentCatalog`（父会话的**直接**子项，按 catalog 事件顺序，**排除 fork 继承来的事实**）、`subagentTiming`、`subagent`（身份，`null` 表示"没有合法 descriptor"且刻意可序列化——注释逐字写：JSON 传输里 `undefined` 会被丢掉、过期身份会活下来）。

**客户端把这三枚拼成一行的过程**（全部有行）：

1. 快照类型 `SubagentCatalogSnapshot`＝投影快照去掉 `values/state`、加 `state:'loading'|'ready'|'error'`、每条 entry 加 `activity:'running'|'inactive'`（`SubagentHeaderLineage.tsx:18-21`）。
2. `activity` 怎么变：**先取统一 UI status，取不到才退到 Session 摘要**——`statuses.get(entry.id)?.running ?? summaries[entry.id]?.running === true ? 'running' : 'inactive'`（`:459-461`）。这就是"状态不建第二套、从事件推导"的落点。
3. token 总量＝**四个互不重叠的桶之和**：`uncachedInputTokens + outputTokens + cacheReadTokens + cacheWriteTokens`（`:67-74`；桶定义在 `llm/token-meter/src/projection.ts:14`、`:71`）。
4. 耗时：`settledMs + max(0, end - active.since)`，`end` 在 running 时取当前秒、inactive 时取 **`active.through`**（`activityDuration`，`:77-91`）。README 逐字钉住这条：**"被中断的未结束轮次以其同一切面的 `active.through` 为上界，绝不使用更新的会话元数据"**（`README.zh.md:66`）。
5. 秒针：只有该层**含 running child** 才起 1 秒 `setInterval`，折叠或关菜单就释放（`:205-211`、`README.zh.md:66`）。⇒ **不是墙钟差实现超时**，是"事件起点＋本地 tick 显示"。
6. 目录不预取根：打开下拉不请求根目录；**展开某个分支或重试失败读取才 `refreshProjection(parentSessionId)`**（`README.zh.md:60`、`:28-40`）；实时投影帧自动更新所有已加载层级，"无需菜单订阅或重复成员查询"。
7. "已知叶子"的判定很硬：**只有这一行自己的目录加载完且为空，它才是叶子**（`:195-198`）；加载中/失败的行仍保持可展开（`README.zh.md:58`）。
8. hover 时序（派单里被点名"抄不出来"的那格，**其实有原文**）：**悬停 150ms 打开、指针离开触发器与目录后 120ms 关闭；点击则钉住，直到点外部或按 Escape**（`README.zh.md:34`）。
9. 键盘：`ArrowRight/Left` 展开折叠，`ArrowUp/Down`/`Home`/`End`/`Escape` 导航或关闭（`README.zh.md:36`、代码 `:286-299`）；行的 `aria-label` 是 `label · secondary · metrics` 三段拼接（`:313`），树用 `role="treeitem"`＋`aria-level`＋`aria-current`（`:309-314`）。

**它自己承认的边界**（`README.zh.md:115-116`）：目录**不区分**失败/取消/拒绝/token 耗尽/尚无已结束轮次（这几样全合并显示成"非完成的 inactive"）；`@label` 只是显示文本，**故意不获得继续执行语义**（label 会重复/改名）。

**我们现在是什么**

- **名册侧已经不是 0 了**：`internal/tools/task.go` 有 `TaskRoster`，方法 **15 枚**（`:210 Record`、`:243 WatchRow`、`:263 Look`、`:276 Count`、`:300 Descendants`、`:337 PublishSubagent`、`:350 MarkRoot`、`:358 TryAcquireSubagentSlot`、`:385 InFlightSubagents`、`:397 RunningSubagentIDs`、`:415 AttachCancel`、`:429 DetachCancel`、`:442 Cancel`）。行结构 `TaskOutput` **6 枚字段**：`Text / ArtifactPath / State / ParentTaskID / Label / Kind`（`internal/tools/task.go:95-140`），`TaskRecord{TaskID, Out}`（`:289-292`）。
- **`task.spawn` 已经落地**：`internal/tools/subagent_197.go` **452 行**，含 `MaxConcurrentSubagents = 4`（`:73`）、`MaxSubagentDepth = 1`（`:79`）、四态映射到 D43 名字（`:94-99`：Thinking/Settling/Muted/Error）、`subagentToolProvider` 把子级的工具目录里**删掉 `task.spawn`**（`:423 subagentHiddenTool`）。
- **缺的三样，逐条带现量**：
  - `TaskOutput` **没有 token 用量、没有任何计时字段**（6 枚字段名逐名见上）⇒ 名册行**给不出 DSH 那一行的"右侧两列"**。
  - 会话事件持久层：**`internal/session/` 只有 1 枚文件 `doc.go`**，其 `:16-17` 逐字写着 `DEFERRED(SessionScope): implemented by ticket 28. This ticket only freezes the package boundary.` ⇒ **没有 seq、没有事件日志、没有投影**，所以"计时"和"状态"这两维在**源头**就还没有承载体。
  - 环回事件：`internal/agent/sink.go:20-38` 的 `EventKind` 只有 **9 枚**（`text_delta/reasoning_delta/tool_start/tool_end/reminder/stuck/control/error/done`），`Event` 结构（`:44-63`）**没有 `Seq`、没有 `Timestamp`、没有 `TurnID`**；`RecordSink` 是纯内存数组（`:73-95`）。对照 DSH 的 **59 枚持久事件类型**＋`SessionSeq`。

**后果**：票 197 的**工具层已经落地**，界面上"点进去那一页"仍然不是画不画的问题——**它没有可导航的载体**：我们的 `SubagentStreamKey(taskID)`（`internal/panel/pump.go:386`、`internal/tools/subagent_197.go:105`）只到"每条子流一行文本"，而 `Snapshot` 的 `results` 元素只有 3 枚 JSON 键 `correlationId/text/done`（`internal/panel/composer.go:72-76`），**"这一行属于哪枚子代理"这句话今天不在线上**（`pump.go:300-306` 的注释自己承认："carrying the marker into the packet is ticket 197's carrier leg (197-r3)"）。⇒ 抄 DSH 那一形，**最该抄的恰恰是"子代理是一份可导航的会话，不是一条流"**——我们的 key 是 `subagent:<taskID>`（流为中心），他们的是 `{parentSessionId, childSessionId, mode}`（会话为中心）。

### 3.3 `ui-trajectory` —— 逐事件台账（"那一页"的第二视图）

**他们怎么做**（数据形状逐条）

- 一行的种类是**闭集 7 枚**：`system / user / context / compacted / message / tool / subtool`（`trajectory-record.ts:9-16`）。**`subtool` 单独一枚** ⇒ 嵌套子工具在同一张台账里有自己的行。
- 一行携带的字段（`trajectory-record.ts:40-105`，共 **28 枚可选字段**），与"看它干了什么"直接相关的：`index`（1 起的 `#N`）、`recordId`、`kind`、`text`、`opensTurn`、**`sourceSeq`（源事件 seq，用于跨行导航）**、`inputDetail`、`outputDetail`、**`thinkingDetail`**、`sourceBlocks`/`outputBlocks`（**保留模型侧原始块序**）、`schemaDetail`（**调用时刻的模型可见工具 schema**）、`assistantMetrics`（TTFT/解码吞吐的原料，`:19-26`：`stepStartTime/firstTokenTime/completedTime/outputTokens`）、`callId`、`toolName`、`isError`、`timeSeconds`、**`startedAt`（真实开始时刻）**、`input/cacheRead/cacheWrite/output/think`（五枚 token 桶）。
- 行身份解析优先级（`:112-117`）：`recordId` → `kind+callId` → `kind+seq` → `kind+index` ⇒ **往前补页时键不变**（`README.zh.md:64`："语义行键与 ARIA 索引在向前补页后保持不变"）。
- 失败字段做了脱敏：`displayFailure()` 对 `code === 'AUTH'` **把 message 清空**（`trajectory-event-projection.ts:127-138`，注释逐字："Provider AUTH messages may echo a masked or partially preserved credential... never retain it in UI state"）。⇒ **错误文本进界面之前过一道"可能带凭据"的判断**，这一条对我们是现成的判据形状。
- 时间轴：Overview 从左到右投影真实开始时间与耗时，助手条区分 TTFT 与解码时间，悬停 500ms 出精确值；**拖选区间→表格聚焦该闭区间内任何时刻活跃的记录**；滚轮缩放时间域；右键清除选区、放大态下右键拖动平移（`README.zh.md:44`）。
- 虚拟化：初始只从"尾部结束的 50 个 target Node"派生；只挂可见行窗口＋少量缓冲；**仅含内容更新的流式帧保持键与高度、复用测量、不重复写末尾滚动位置**（`README.zh.md:64`）。
- 它自己承认：**进行中不虚构耗时**（`partial`/`runningCalls` 只显示开始标记，Time 留空，`README.zh.md:104`），**不提供锚点深链接**。
- 层位：`inject` 逐名读服务＝`['slots','sessions','uiSession','uiConversation','locale']`（`index.ts:41`）；台账是**纯投影**，"既不读取也不改变 Chat 会话快照"（`README.zh.md:62`）。数据来自 Host 侧会话窗口，经 WebSocket `$events`/Remote 流到浏览器（`client/connection/README.zh.md:32`：一元调用走 HTTP POST，`/api/remote.mux` WebSocket 持有逻辑流）。**没有子进程、没有 OS API。**

**我们现在是什么**：`trajectory` 关键词在 Go 侧 **0 命中**（`grep -rni trajectory --include=*.go internal cmd` = 0）。事件词汇 9 枚 vs 59 枚（§3.2）。**〔未量·frontend 地界〕**：面板有没有"逐事件台账"那一屏，本程不读。

### 3.4 `ui-settings-subagent` —— 把上限做成可设项

**他们怎么做**

- 一页两半、一次保存：`subagent` 与 `subagent-model-selection` 两个 Host 命名空间收在**同一次保存**下；**两次写入相互独立**——Host 拒的那半留草稿并报失败，另一半照常落地；被更新修订号超越的模型草稿**报冲突并要求放弃，而不是覆盖较新的路由**（`README.zh.md:32`）。
- 字段就两枚：`maxDepth`（最小 0）、`maxActiveSubagents`（最小 1）（`src/client/subagent-limits-card-controller.ts:10-13`、`:48` 的 `limitField('maxDepth', 0)` / `limitField('maxActiveSubagents', 1)`），且 `Number.isSafeInteger && !Object.is(v,-0)` 才收（`:36`）。
- 计口径写在界面上：**容量统计同一主 Agent 下所有层级存活的子代理，主 Agent 不计入；深度让位于工具自己的上限**（`README.zh.md:28`）。每个标签旁有信息按钮展开规则（深度是一张两行示例表）。
- 页面存活条件：`ctx.configForms.whileServed` 监视两个命名空间，**Host 服务其中任一枚，页面就存在，并只显示被服务的部分**（`README.zh.md:12`、`:42`）。
- 模型路由的边界：开启时至少一条路由；关闭保留已选；**已存储但没有适配器公布的路由留在"已保存但当前不可用"分组且仍可移除**；加载失败的提供方被报告而不隐藏其他提供方，并提供重试（`README.zh.md:30`）。

**我们现在是什么**：`MaxConcurrentSubagents = 4`、`MaxSubagentDepth = 1` 两枚 `const`（`internal/tools/subagent_197.go:73`、`:79`），注释逐字说明 4 这个数字**来自 `internal/tools/bridge.go` 的 `MaxToolConcurrency`（冻结契约、不许从名册侧上调）**，并且故意写成字面量 4 而不是别名，让判据比较"两次独立读数"。配置面：`internal/config/schema.go` 里 grep `subagent` = **0 命中**。

**后果**：这两枚数字**今天是编译期的**，owner 要改一次就得重新编译。DSH 那一形告诉我们：**上限＋计口径＋"谁让位于谁"这三句话要一起摆到设置页**，只放一个数字框会让人误读（他们连"主 Agent 不计入"都要写成字段旁的说明）。⚠ 这条会**顶到 D36 配置模型的 section 树**（`PLAN.md:2715`），属契约面，不是纯界面单。

### 3.5 `ui-dockkit` —— 停靠／分栏／可撤销布局

**他们怎么做**

- 两层：**引擎是纯逻辑**（无框架、无 DOM、无宿主概念），**组件只渲染快照＋上报意图**（`README.zh.md:31`、`:40`）。
- 状态：归一化递归分裂树 `LayoutState{nodes, tabs, rootId, floats(自底向上), activePaneId, expanded, mode}`（`src/contract/types.ts:105-118`）。`PaneId/SplitId/TabId` 是 **brand 字符串，只有 `Mint` 或库自己的 DOM 往返能产出**（`README.zh.md:33`）⇒ 三种 id 互相不可替、不能拿裸串充数。浮窗不是第二个概念：**它就是 `host:'float'` 的格，容量一个 tab，绘制时不带 tab 条**（同段）。
- **可逆操作**：`applyOp(state, op)` 返回**下一状态＋撤销它的操作**；逆操作在**执行那一刻**捕获，因为到撤销时"操作前的状态已经不存在了"（`README.zh.md:34`、`types.ts:197-201`）。op 携带它创建的 id ⇒ `replay(initial, ops)` 复现同一棵树，**引擎不读时钟不读随机源**（`README.zh.md:35`）。
- 操作词表 **16 枚**（`types.ts:137-192`）：`split / merge / openTab / closeTab / moveTab / reorderTab / focusTab / focusPane / resize / float / unfloat / moveFloat / resizeFloat / setExpanded / setMode / insertPane / insertTab / restoreFocus`（数一下是 18 个 type，**逐名列表**：`:140 split`、`:149 merge`、`:151 openTab`、`:153 closeTab`、`:155 moveTab`、`:157 reorderTab`、`:159 focusTab`、`:161 focusPane`、`:163 resize`、`:165 float`、`:167 unfloat`、`:169 moveFloat`、`:171 resizeFloat`、`:173 setExpanded`、`:175 setMode`、`:178 insertPane`、`:184 insertTab`、`:187 restoreFocus` ＝ **18 枚**）。
- 历史：`Sequencer` 维护线性历史，**每个意图一条记录**（一次手势产出的多个操作一起后退／前进）；连续纯焦点记录合成一步（`FocusOpType = 'focusTab'|'focusPane'|'restoreFocus'`，`:194-195`）；后退后的新记录丢弃前进分支（`README.zh.md:36`）。⇒ **关键设计：组件绝不上报拖动帧，一次手势只在松手时上报一条净结果**（`README.zh.md:40`）。
- 可选规则 `planSettle`：保证意图之后每个停靠格都有内容；被清空的格并掉，被清空的根格用嵌入方工厂重新播种（不传工厂就只并格）（`README.zh.md:37`）。
- 策略以 props 进入：`canSplit`（整面格预算）/`canAddTab(paneId)`/`canCloseTab(tabId)`/`minPaneFraction`（Sidebar 用 0.2）/`hideSplitWhenBlocked`（`README.zh.md:55-57`）。
- 内容身份是二元组 `(kind, contentId)`：`findContentTab`、`planOpenContent` 会**聚焦已存在的 tab 而不是再开一个**（`README.zh.md:59`）⇒ 同一枚子代理不会被开成两页。
- 几何分工：**布局状态只携带比例、从不携带像素**；`geometry.ts` 的 `halvesFit` 是算术，`measure.ts` 在每次提交后与尺寸变化时读矩形（`README.zh.md:75`）。分栏要求"两个可用半格都容得下不可收缩部分"，规则逐字写在 `:75`（这一整段是踩坑记录，不是设计偏好）。
- 无障碍现状（他们自己登记）：**分隔条没有 `separator` 角色，没有键盘路径去分栏/移动/浮出**；触控未调优；尺寸语义刻意精简（无吸附、无优先级、无首选尺寸）（`README.zh.md:98-101`）。
- 层位：**浏览器主线程内纯计算＋DOM**，无子进程、无服务端、无 OS API。**静态链接**发布（`README.zh.md:83`），且**是内部引擎、导出无任何版本承诺**（`README.zh.md:14`）。

**我们现在是什么**

- **我们有一个叫 `dock` 的文件，但它不是这个 dock**：`internal/ball/dock.go` 讲的是球的**贴边/钳制**（`:166` 逐字 "Dock axis last, so the clamp cannot undo it"）——`grep -rniE "dock" --include=*.go internal cmd` 那 **442 行**绝大多数落在这族与 `docs/` 用词上。⇒ **报"我们有 dock"会是假阳性，这条得按文件名报。**
- 面板结构：`Snapshot` 是 **5 个固定 JSON 键**（`internal/panel/composer.go:57-69`：`pending/results/composer/generatedAt/instructions?`），`NewSnapshot` 的注释把 composer 段称作 MANDATORY（`:78-82`）。⇒ **形状上是固定四／五区，不是一棵可分裂的树**。
- 撤销：`grep -rniE "undo" --include=*.go internal cmd | grep -v _test.go` = **20 行**，**逐名全不是布局撤销**——审批否决那一族（`internal/agent/approval/gate.go:214`、`report.go:13`/`:59`/`:133`，B1 的"绝不把取消渲染成干净撤销"）与球贴边钳制（`dock.go:166`）。另：`grep -rniE "layoutop|applyop|dockingsurface" --include=*.go internal cmd` = **0**；`inverse` 只有 **1 行**，是 `internal/config/parse.go:26` 讲 TOML 序列化反方向的注释，与布局无关。⇒ **今天没有任何一枚界面操作带逆操作。**

**后果**：这条**不建议直接开抄的票**——DSH 自己说这是"内部引擎、导出随时会变"（`README.zh.md:14`），抄它的 API 等于抄一份无承诺的接口。**但它回答了一个我们必须自己拍板的问题**：如果面板四块固定区要保留，那么"点开一枚子代理"就只有两种落点——顶替主区，或者新造一格。**DSH 两种都给了**（主区导航＋右侧 tab，`index.ts:61-69`），且用 `(kind,contentId)` 去重保证不重复开。⇒ 真正该开的票是**"子代理会话的两种落点＋去重"**，而不是"停靠系统"。

### 3.6 `ui-layout` —— 三栏外壳与挤压次序

**他们怎么做**

- 三栏 AppFrame＋数值全在明处：`CENTER_MIN=400`、`SIDEBAR_MIN=264`、`SIDEBAR_MAX=420`、`SIDEBAR_DEFAULT=280`、`SIDEBAR_COLLAPSED=56`、`SIDEBAR_AUTO_COLLAPSE=1024`、`RIGHTBAR_MIN=300`、`RIGHTBAR_DEFAULT_RATIO=0.45`（`src/client/columns.ts:11-29`，8 枚具名常量）。右栏上限 70%（`README.zh.md:30`）。
- **确定的挤压次序**（这段对我们最有用）："为给中栏保留 400px，框架**先**把右侧面板缩减至 300px，**再**报告空间不足使占用方将其关闭，**最后**才进一步压缩中栏"（`README.zh.md:30`）；变宽**不自行重新展开**，全屏隐藏宽度手柄但不自行释放轨道（`:54`）。
- 窗口 chrome 座：macOS 下收起的侧栏**整列隐藏**，框架在左上角挂单一 `shell.leading` 座（红绿灯旁），并发布 `--dsh-frame-leading-clearance`（chrome 占的带宽）与 `--dsh-frame-top-clearance`（48px，**发布在根元素上，好让 portal 到 body 的浮层也读得到**）（`README.zh.md:36`）。Windows Electron 用 `data-windows-titlebar` 预留顶栏高度、移除收起轨道，内容区**只有左上角 16px 圆角、其余直角**，并发布 `--dsh-windows-content-radius`/`--dsh-windows-sidebar-width`（`README.zh.md:38`）。**每个 chrome 行自己打 `data-window-drag`，于是行自己的盒子就是可拖几何——行的空白段可拖，控件保持可点**（`:36`）。
- `Mod+B` 切左栏；模态弹窗打开时不可用；终端输入优先（`README.zh.md:28`）。
- 边界登记：面板几何是**瞬时状态**（刷新恢复默认、宽度偏好是整个框架共用的一份而非每会话）；极窄窗口下中栏仍可能小于 400px；挤压重排期间**无滚动锚定**（`README.zh.md:89-92`）。

**我们现在是什么**：`--dsh-frame-top-clearance` 同形物 0（Go 侧不产 CSS 变量）；`grep -rni "hide_on_fullscreen"` 在 `docs/` 有规格（`SPEC-08*.md:229` `[ball] hide_on_fullscreen`）但那是**球**不是面板。**窗口拖拽几何由谁声明**这一条：〔未量·frontend/WebView2 地界〕。

**后果**：**"谁让位、按什么顺序让位"这句话我们必须自己写出来**。我们面板＋球两块面并存，比 DSH 多一层；DSH 只写了右栏/中栏两级的次序，**没有写"贴边面板 vs 悬浮球"的次序**。这是我们自己要拍的空格，不是抄来的。

### 3.7 `ui-shortcuts` —— 快捷键全可改

**他们怎么做**

- 可改到什么程度：**除固定输入操作外全部可录制、可清除、可恢复默认**；`Mod+/` 打开速查（无模态在最前时；速查在最前时它反过来关闭自己）（`README.zh.md:12`、`:30`）。
- 持久化分平台两形（`README.zh.md:36`）：**Web = 当前来源的 `localStorage` 项 `dsh.keybindings.v1`**（键名常量真在 `packages/client/shortcuts/src/client/storage.ts:6`）；**Desktop = `userData/keybindings.json`**。
- 录键规则细到可直接当判据抄（`:36`）：Windows/macOS Desktop 可录**一至两枚按住时间重叠的非修饰键**＋修饰键，**第三枚直接拒**；**释放第一个键即保存当前组合**，所以先后按 A、B 只录进 A；保存拒被可编辑或固定操作占用的组合（含相互重叠的单键与双键）；写失败**保留录键框焦点和草稿，可重试**；松开非修饰键后可直接录新组合，不必再点。
- 读失败不静默：**禁用编辑，并通过系统提示说明当前使用的键位及受影响的配置位置**；内容损坏时先备份再修；配置由更新版本写入时应升级 Harness；**"恢复全部默认"不覆盖无法读取的数据**（`:36`）。
- 冲突与安全面：macOS Web 允许录被报告为 dead key 的 `Option+Command+N`，**但输入法组合和普通重音输入仍受保护**（`:36`）；IME 组合期间 Enter 只确认候选（`:28`）。
- 固定操作只有 **3 枚**：`fixed.move`(↑↓)、`fixed.select`(Enter)、`fixed.dismiss`(Esc)，`src/client/fixed.ts:10-17`——即"菜单移动/确认/关闭"这三条键位是**只读保留**，其余都可改。审批与双 Esc 停止行**由已挂载的功能 owner 贡献**，不在可编辑集里（`README.zh.md:77`）。
- 排序规则：核心操作用**固定产品顺序，不受插件注册/卸载/重挂影响**；未列入该顺序的扩展命令按命令 ID 稳定排在核心之后（`:32`）。计数与恢复：底部显示"恢复全部默认"＋**当前运行环境与平台的自定义项数（包括已移除的绑定和当前未挂载的命令）**；无自定义项时隐藏计数并禁用该操作（`:36`）。
- 层位：**浏览器/桌面 renderer 进程内**，写 `localStorage` 或 `userData` 文件；按键的实际生效在另一枚服务包 `client/shortcuts`（`packages/client/README.md:33`：`ctx.shortcuts` 注册应用键盘命令）。

**我们现在是什么**（⚠ **派单里那句"我们完全没有"不准，我这腿量出来是"有 4 枚可改、但没有编辑器"**）

- **可改键位已经有了，而且是热改**：`internal/config/schema.go:176-183` 的 `[hotkey]` section 有 **4 枚字符串字段** `Summon / Mute / Cancel(默认 "Esc") / Panel`，注释逐字写 "hot-tier (hotkeys re-register on change)"、"Empty string = binding unset (feature key disabled)"、"`cancel` temporarily takes over during Confirming and is returned at session end (D36)"。
- **热改链条也在**：`internal/ball/hotkey_reload.go`（**新建的桥**）定义 `HotkeyBinder`(`:44-46`)、`HotkeySource`(`:49-50`)、`HotkeyReloader.Check()`(`:80+`)；`Ball.RebindHotkeys(cfg) HotkeyReport` 在 `internal/ball/ball_windows.go:801`。`:42-44` 注释还专门解释了**为什么是轮询而不是 `OnReload`**：`[hotkey]` 是 HOT-tier，`OnReload` 只对 RELOAD-tier 触发，"a hotkey edit would never reach a callback installed there"。
- **键位串的形状也在**：`HotkeyConfig{Summon,Mute,Cancel,Panel string}`（`internal/ball/hotkey_windows.go:42-47`，示例 `"Ctrl+Alt+Q"`），`HotkeyBinding{Name,ID,Binding}`＋`HotkeyReport.Bindings()/Binding(id)`（`:264-301`）——**已经有一份"当前哪些键注册成功了"的可读报告**。
- **`AltSummonSpace = "Ctrl+Alt+Space"`（`:53`）注释逐字记录了 IME 抢占**："several Chinese IMEs (Sogou/WeChat-style shift-space and ctrl-space commits) grab Ctrl+Alt+Space before the app can"，且 `DefaultHotkeys()`（`:61-62`）注释记着 owner R10 拍板：**在 owner 本机扫过 84 枚候选组合，只有 `Ctrl+Alt+W` 被第三方占用**。⇒ **这一条我们比 DSH 早踩：他们 README 里的 IME 保护是 Web 层的，我们有 Win32 全局层的实测台账。**
- **真正缺的三样，逐条带现量**：
  1. **没有"命令目录＋速查窗"**：`grep -rn "shortcut" --include=*.go internal cmd` 里没有一枚"列出所有可绑命令"的目录类型（只有 `ChannelEsc` 这类通道常量）。**〔未量·frontend 地界〕**：面板里有没有那一扇窗。
  2. **没有录键 UI／冲突拒／恢复默认**：`keybind` 大小写不敏感命中 **14 行**（`internal/ball/hotkey_windows.go` 与 `internal/ball/hotkey_reload.go` 的 `HotkeyBinding` 家族＋1 行测试注释），**逐名都是"注册"侧，没有一枚是"用户改键"侧的校验**。
  3. **热改桥在真实宿主里还没接上**：`NewHotkeyReloader` 的生产调用者**只有 1 处**——`cmd/balldebug/main.go:237`；`cmd/wisp/` 下 **0 处**（`grep -rn "NewHotkeyReloader" --include=*.go .` 的完整结果＝balldebug 1 行＋`internal/ball/hotkey_live_test.go` 2 行＋`hotkey_status_test.go` 2 行＋定义文件自身）。⇒ 按本仓那套三档证据强度，这一格是 **〔建了但没接〕**。
- 已有的"必须归还"判据（我们的 D22 式形状）在 `SPEC-08*.md:226-227` 与 `internal/agent/approval/approval.go:55`（`ChannelEsc`）：**Esc 同时是"取消键"（`[hotkey].cancel`）和"Confirming 期间的临时接管键"**，而它还是审批四通道之一。

**后果**：这不是"从零做快捷键"，是**"把已有的 4 枚可改键补上三件配套"**：命令目录、录键校验、真实宿主里的热改接线。⚠ 但有一条**必须他拍板**（我不越权）：一旦 `Cancel`（=今天的 `Esc`）可改，`AGENTS.md` §1.2/`SPEC-08*.md:226-227` 钉的"Esc 在 Confirming 期间临时接管、会话结束必须归还原绑定"这条规则的前提就变了——**它假设 `Esc` 是那个可被临时接管的键**。DSH 那一形里这件事的答案是现成的：**审批挂载期间 Enter/Esc 不许分配给可编辑快捷键**（`ui-approval/README.zh.md:21`），且审批与双 Esc 停止行**在快捷键目录里是只读贡献项、不可改**（`ui-shortcuts/README.zh.md:77`）——即"**可改集与不可改集分开，安全相关的进不可改集**"。我们已经有这个形状（四通道是闭集），只是没人把它说成快捷键策略。**这一格写进差集清单时请写成"缺策略"，不要写成"缺功能"。**

### 3.8 顺手多读的两枚（不在派单 8 枚里，如实标注）

**`ui-jobs`**（后台任务那一枚，和票 197 直接同题，`ui-jobs/README.zh.md:11`、`:28-38`）

- 名册**只有一份**：`ctx.jobs` 镜像的 `job.list` **流**是唯一名册，"不存在需要 join 的第二份名册"（`:28`）——这一句正是 DSH 用**投影**替代"两份名册拼合"的写法。
- **两击式停止控件**：首击武装，**三秒内的确认击**才调 `ctx.jobs.kill`；行状态经名册流收敛（先 `stopping`，再进已结束分组，detail 携带 `cancelled by the user`）（`:30`）。
- **人类停止任务要显式告诉模型**："该 kill 不在模型的播报台账里认领任何东西，owner agent 照常收到标准完成通知——**模型被明确告知用户停止了它的任务，而不是留给它去猜**"（`:30`）。这一条和我们的 C22"触发必须显式告知用户"（`internal/agent/sink.go:10-15`）是**同一句话的反方向**：我们只管告诉人，他们还要告诉模型。
- 收起即停流：输出观测流跟随展开/收起/卸载/弹层关闭开合（`:34`、`:38`）。展开面板复制的是**命令**不是输出（`:34`）。
- 未使用不占位："在会话看得到至少一个 job 之前它不渲染任何东西，因此普通对话不会为未使用的能力长出控件"（`:26`）。

**`ui-conversation`**（那一页的宿主，`ui-conversation/README.zh.md:40`、`:67`）

- 视图环：target 包通过 declaration merge 扩 snapshot 与 Location data map，再 `ctx.uiConversation.events.register(...)` ＋ `views.register(...)`；读自己的数据走 `binding(sessionId).target(targetId)`（`:40`）。
- **View 偏好是持久化的**：Session 首次绑定或缓存 Session 变 current 时，shell 在渲染前**读取持久化 View 偏好**，激活已注册的偏好 View 或 Chat fallback（`:67`）；`ConversationStoreState.view: string | null`（`src/client/contract/views.ts:18-22`），`ViewTab{id,label}`（`:7`）。
⇒ 这条决定"点进去那页"的体验：**你上次看子代理用的是台账视图，下次点进来还是台账视图**。

---

## 4. 层位表（判断能不能落到 Win32＋Go 上用）

| 机制 | 跑在哪一层 | 出处 | 落到我们身上的对应层 |
|---|---|---|---|
| 名册/计时/token（子代理状态） | **服务端进程内投影**（Host 折叠事件成投影值）→ 逐会话快照推给浏览器 | `subagent/subagent/src/projection-types.ts:68-85`；`api/session-controller` 的 `projectionsBySession` | 我们的等价物**不存在**：`internal/session/` 只有 `doc.go` |
| 审批/提问跨进程 | **Host→浏览器 Remote Event waterfall**（WebSocket `$events` 流；一元走 HTTP POST） | `api/remotes/src/remote-events.ts:22,47`；`client/connection/README.zh.md:32`；`ui-approval/src/client/index.ts:104` | 我们是 Go 侧 channel＋correlationId＋WebView2 桥（`SPEC-08*.md:156-176`）——**形状不同，但可承接 waterfall 的"多个应答方按作用域让位"语义** |
| 停靠/分栏/撤销 | **浏览器主线程内纯逻辑＋DOM**，无服务端、无子进程、无 OS API | `ui-dockkit/README.zh.md:31-40` | 可原样落 WebView2，也可只抄"逆操作在执行时捕获"这一条规则 |
| 三栏几何/窗口 chrome | **CSS＋DOM＋桌面 preload 打的标记**（`data-platform`/`data-windows-titlebar`/`data-window-drag`） | `ui-layout/README.zh.md:36-38` | Windows 标题栏与拖拽几何在我们身上是 WebView2 宿主＋`internal/ball` 的地界 |
| 快捷键录制与持久化 | **renderer 进程内**；Web 写 `localStorage['dsh.keybindings.v1']`，Desktop 写 `userData/keybindings.json` | `client/shortcuts/src/client/storage.ts:6`；`ui-shortcuts/README.zh.md:36` | 我们是 `config.toml`（D36）；DPAPI 保护范围要另判 |
| Trajectory 台账 | **浏览器侧纯投影**（读 Host 会话窗口的事件，不改 Chat 快照）＋虚拟化列表 | `ui-trajectory/README.zh.md:62-64` | 需要"带 seq 的持久事件"当源头，**我们先决条件是 `internal/session` 落地** |
| 子代理页里的工具执行 | **Host 侧**（本包逐字声明："本包绝不接收宿主上下文，也不调用面向模型的工具"） | `ui-subagent/README.zh.md:70` | 与我们的 D33/票 197 "子级不自批、批准上抛父作用域"同形 |

---

## 5. 我这一腿没读到什么（具名到包）

**深读源码为 0 的 44 枚**（只读了各自 `README.zh.md:2` 的 frontmatter description 与 `package.json:3`，没读 `src/**`、没读 `tests/**`）：

`ui-agent-preset` `ui-attachment` `ui-brand-official` `ui-chat` `ui-commands` `ui-deliverables` `ui-directory-picker-browse` `ui-directory-picker-native` `ui-goal` `ui-input-trigger` `ui-message-feedback` `ui-model-selection` `ui-open-in-app` `ui-permission-presets` `ui-plan` `ui-plugin-manager` `ui-primitives` `ui-reference` `ui-renderer` `ui-schedule` `ui-session` `ui-settings` `ui-settings-account` `ui-settings-agent-loop` `ui-settings-general` `ui-settings-models` `ui-settings-plugin-inventory` `ui-settings-plugins` `ui-settings-shell` `ui-settings-web-search` `ui-sidebar` `ui-sidebar-browser` `ui-sidebar-documentpreview` `ui-sidebar-files` `ui-sidebar-right` `ui-sidebar-terminal` `ui-skill` `ui-slots` `ui-theme` `ui-tool` `ui-workflow-run` `ui-workspace`
（＋ `ui-conversation`、`ui-jobs` 只读了 README 概述段与一两个 contract 文件，未读全部源码）

**具体没读到的东西（会影响结论强度的，逐条列）**

1. **`ui-*` 之外的一切**：`client/connection`、`client/store`、`client/modules`、`client/resources`、`client/web`、`client/shortcuts`、`client/file-upload`、`client/hmr`、`client/locale`、`client/test-runtime` —— 我只引了 `connection/README.zh.md:32` 与 `shortcuts/src/client/storage.ts:6` 两处，其余整树未碰。
2. **Host 半边的实现**：`packages/subagent/subagent/src/index.ts`、`tool-subagent/src/index.ts`、`child-agent.ts`、`continuation-activation.ts`、`interaction/user-approval/src/index.ts` —— 我只读到 `projection-types.ts`、`types.ts`、`remote-events.ts` 和 `subagent/src/index.ts:482` 那一行 `@Remote` 声明。**编排者 A401 引的那几行（`tool-subagent/src/index.ts:262,270`、`subagent/src/index.ts:201-203`、`child-agent.ts:172-176,231-235`、`continuation-activation.ts:45-50,604-611`）本程一枚都没核。**
3. **`ui-trajectory/TrajectoryTable.tsx`（3513 行）与 `layout.ts`（1179 行）没通读**，只读了 `trajectory-record.ts`、`trajectory-event-projection.ts`、`index.ts` 全文与 README。⇒ §3.3 里所有"检查器/代码视图/搜索"的细节都出自 README 的自述，**没有第二处代码佐证**。
4. **测试没跑、没读**：52 枚每枚都有 `tests/`（例：`ui-subagent/tests/{browser-plugin,conversation-ui,sidebar-chat}.client.spec.tsx`），我一个没打开。⇒ "他们的行为被钉住了"这件事我只看到 README 里写了"由本包 spec 直接断言"，**没验**。
5. **`ui-dockkit` 的 `planner.ts`/`operations.ts`/`sequence.ts` 没读实现体**（只读了 `contract/types.ts` 全文与 `operations.ts` 的两个导出签名 `:423 applyOp`、`:462 replay`）。⇒ "18 枚 op 都能给出逆操作"是**契约声明**，不是**我逐 op 读到的实现**。
6. **`ui-sidebar-*` 五枚全未深读**（`ui-sidebar-right` 是停靠面的真正嵌入方，我只从 `ui-dockkit` 侧看它）。
7. **DSH 的 TUI／Desktop 半侧**：`ui-approval` 的"持久权限策略由 Host 侧审批包拥有"（`README.zh.md:36`）指的包我没读；`ui-settings-account` 提到的 Desktop 授权流我也没读。

---

## 6. 最值得开的 3 枚票（是什么／后果／在哪一层）

### 票 A：**子代理的"身份载体"从"一条流"换成"一份可导航会话"**（改形状，不是加控件）

- **是什么**：DSH 的寻址是 `{parentSessionId, childSessionId, mode}`（`client/ui-subagent/src/client/index.ts:61-69`），资源地址 `dsh-resource://subagentchat/session/<child>?parent=<parent>&mode=<mode>`（`sidebar-chat/index.tsx:21,51-57`），**去重靠二元组 `(kind, contentId)`**（`ui-dockkit/README.zh.md:59`）。我们今天的寻址是 `subagent:<taskID>` 一枚**流键**（`internal/panel/pump.go:379,386`、`internal/tools/subagent_197.go:105`），`ResultChunk` 只有 `correlationId/text/done` 三枚 JSON 键（`internal/panel/composer.go:72-76`）——**"这行是谁的"这句话不在包上**（`pump.go:300-306` 自己承认载体系在 197-r3）。
- **后果**：不改形状，票 197 的界面腿只能画出一列文本；画完也**做不到 owner 说的"点进去"**，因为点进去需要一个可导航的会话身份。更要命的是：**流键合不了、也切不走**——owner 想看的是"那一枚子代理的工作"，我们给的是"混在一个 results 数组里的若干行"。
- **在哪一层**：Go 侧 `internal/tools`（名册已有 `Descendants`/`PublishSubagent`）＋ `internal/panel`（Snapshot 载体）＋ **C17 名册要不要新增"打开某子会话"这一个方法**（`SPEC-08*.md:156-176` 那张表今天有 `tasks.list`/`task.detail`，**没有"把这枚会话装进某一格"**）—— ⇒ 这一枚**碰契约**，不是纯界面。

### 票 B：**"上限＋计口径＋谁让位于谁"三件一起摆进设置面**（`maxDepth` / `maxActiveSubagents` 现在是编译期常量）

- **是什么**：DSH 一页两半一次保存，两个命名空间独立写、冲突报不覆盖（`ui-settings-subagent/README.zh.md:12,32`），字段带最小值与安全整数校验（`subagent-limits-card-controller.ts:36,48`），**并把"容量统计所有层级、主 Agent 不计入；深度让位于工具自己的上限"写成字段旁的说明**（`README.zh.md:28`）。我们是 `MaxConcurrentSubagents = 4`／`MaxSubagentDepth = 1` 两枚 `const`（`internal/tools/subagent_197.go:73,79`），`internal/config/schema.go` 里 `subagent` **0 命中**。
- **后果**：owner 想放宽一次深度就得重新编译；而他真会想做这件事——因为**深度 1 意味着 §3.2 那棵"任意深度的展开树"（`SubagentHeaderLineage.tsx:368-394` 是递归渲染）我们根本用不上**。更要紧的是 4 这个数字**是从冻结契约 `MaxToolConcurrency` 借来的**（`:60-72` 注释逐字说明），把它做成设置项会**打开一条"名册侧上调到桥之上"的路**——那正是票 211 记的那枚谎。
- **在哪一层**：`internal/config`（D36 的 section 树）＋ `internal/tools` 的常量读取点＋ `internal/panel` 的 `config.get/config.set`（C17 已有方法）。**碰 D36 与 D38/D32 那两枚上限的归属** ⇒ 契约级，人工批准。

### 票 C：**人类否决要同时告诉模型**（`ui-jobs` 那枚两击停止的**反向告知**）

- **是什么**：DSH 的 job 停止控件是**两击**（首击武装、3 秒内确认击才 kill），并且**明确告诉模型"用户停了你的任务"**，而不是留给它猜（`ui-jobs/README.zh.md:30`，附决策笔记链接）。我们：`grep -rnE "两击|confirm.{0,20}[Ss]top|arm.{0,20}[Kk]ill" --include=*.go internal cmd` = **3 行**，**逐名看全不相关**（`internal/agent/approval/batch_test.go:141`、`internal/risk/pathshape_portable_test.go:195` 等是审批确认的**测试注释**，不是"两击停止"控件）⇒ **可数的硬读数：两击式武装停止在 Go 侧 0 个生产落点**。告知模型这一半：`grep -rn "cancelled by the user" --include=*.go internal cmd` = **0**。我们的 `EvStuck`（`internal/agent/sink.go:30-32`）钉的是"C22 刹车必须**告诉人**"（`:10-15` 注释逐字："the C22 rule 触发必须显式告知用户 (never stop silently)"）——**方向只有一条**。
- **后果**：人按了停止，模型下一轮看到的是"任务自己结束了"。这不是体验问题，**是正确性问题**：模型的后续判断会建立在"它以为自己说完了"上面。这一枚也**顺带解决 D45/D31 里"谁取消了这一发"的记账**——取消原因要有一个人可看、一个模型可看的两份出口。
- **在哪一层**：纯 Go 侧（`internal/agent` 的事件＋`internal/panel` 的 composer 状态），**界面只在右侧多一个"已确认/已中止"的读回**；不碰 C17、不碰 D34。**这一枚是我们今天就能自己做完的一枚。**

**（备选，第 4 名给 `ui-user-questions`）**：`ask_user_question` 整条 0 命中。它的票要排在**"提问走不走审批那条作用域路由"**拍板之后，否则会把两条通路做成两套。

---

## 7. 一张对照表（给差集清单直接并）

| 界面这块 | DSH 有（行） | 我们现在（现量） | 差 |
|---|---|---|---|
| 子代理名册一行带 token＋耗时 | `SubagentHeaderLineage.tsx:253-268`；`tokenTotal:67-74`；`activityDuration:77-91` | `TaskOutput` 6 字段无 token 无计时（`task.go:95-140`） | **全缺** |
| 子代理状态从事件推导 | `:459-461`（UI status 优先、退到摘要） | `TaskOutput.State` 是 `statemachine.State`，写入方是 host（`:104-115`） | 形似而**源头不同**：他们没有"第二套状态"，我们也不该有（票 196/A394 已按住） |
| 点进去＝那份会话自己的页面 | `index.ts:61-69`＋`sidebar-chat/index.tsx:145` | 无（`Snapshot` 5 固定键，`composer.go:57-69`） | **全缺** |
| 事件台账＋时间轴（TTFT/解码/吞吐） | `ui-trajectory` 全栈；59 枚持久事件 | 9 枚瞬态 EventKind、无 seq 无时间戳（`sink.go:20-63`） | **缺源头**（`internal/session` 只有 `doc.go`） |
| agent 反问人（带选项、可跳过、可自定义） | `user-questions/types.ts:36-67`＋`remote-events.ts:47` | 0 命中 | **全缺** |
| 审批只有"仅本次允许／拒绝"，持久策略在别处 | `ApprovalDecision:66`＋`README.zh.md:36` | 面板只能 reject，allow 被结构禁（`SPEC-08*.md:156-176`） | **方向相反**，见 §9-⑧ |
| 停靠／分栏／可撤销布局 | 18 枚 op＋逆操作执行时捕获（`types.ts:137-192`；`README.zh.md:34`） | Go 侧无同形；面板固定 5 键 | **无此概念**（但见 §6 票 A 的建议：不要抄整套） |
| 快捷键全可改＋冲突拒＋恢复默认 | `ui-shortcuts/README.zh.md:36`；`fixed.ts:10-17` | **已有 4 枚可改键＋热改桥**：`config.toml [hotkey] summon/mute/cancel/panel`（`internal/config/schema.go:179-183`）、`HotkeyReloader`（`internal/ball/hotkey_reload.go`）、`RebindHotkeys`（`ball_windows.go:801`）；**缺**命令目录/录键校验/真实宿主接线（`NewHotkeyReloader` 生产调用只在 `cmd/balldebug/main.go:237`，`cmd/wisp/` 0 处） | **缺三件配套，不缺底座** |
| 后台任务两击停止＋告知模型 | `ui-jobs/README.zh.md:30` | 停止告知人（`sink.go:10-15,30-32`），告知模型 0 | **半缺**（票 C） |
| 上限做成设置项 | `ui-settings-subagent` 整包 | 2 枚编译期常量 | **全缺**（票 B） |

---

## 8. 关于"我们的规矩"的两条提前报备（不是缺口，是撞车预告）

1. **DSH 每秒 tick 的写法会踩我们的 D22 禁令**：它的计时是 `settledMs + (now - active.since)`，`now` 由 1 秒 `setInterval` 推（`SubagentHeaderLineage.tsx:205-211`、`:87-90`）。**这是墙钟差**。`AGENTS.md` §1.2 禁的是"**用墙钟时间差实现超时**"——显示不算超时，但**我们的写手如果照这一形抄，门禁那一关会把它读成违规形状**。A401 里编排者已经预判到并自裁"不许照抄"（`pending-and-issues.md` A401 第③条），**我这腿确认他那句是对的，并且给出可用的替代**：DSH 自己的存储侧本来就是"事件起止＋累计"形（`projection-types.ts:22-37` 的 `settledMs/since/through`），**tick 只在 renderer 一侧**——照这个分层抄就不撞。
2. **零 emoji 规范在这块地界有具体形状**：DSH 全程用**具名 SVG 图标组件**（`SubagentHeaderLineage.tsx:167-188` 是一枚手画 `SubagentSwitcherIcon`，`:9-11` 从 primitives 引 `IconChevronDown/Right/Refresh`），没有一枚字符图标。这与我们 `lucide-react` 的选型一致（`SPEC-08*.md:176-186`），**但 DSH 的 `StateDot` 是共享组件、三态一名**（`ongoing/done/idle`，`:334`），且**每行为状态图标预留同一 14px 列宽**（`README.zh.md:36`）。⇒ 我们界面那支若要摆 14 种状态的视觉语言（`SPEC-08*.md:186`），这是**同一形问题的同一解法**，值得点名给他们，不是我们该开的票。

---

## 9. 我认为编排者会误读的地方（他的哪条结论可能被推翻）

**① 他 A397/A400 那句"Go 侧子代理实体 0 命中"——今天已经不成立了**
- 他那句的现量口径是 `grep subagent|Subagent|SubAgent` 在非测试码 0 命中（`pending-and-issues.md:8564`、`HANDOVER.md:1253`）。
- **我 20:2x 现读**：`grep -rln subagent --include=*.go internal cmd` = **6 枚文件**，其中 **3 枚不是测试**：`internal/tools/subagent_197.go`（**452 行**，`task.spawn` 工具、深度 1、池 4、D43 四态、taint 盖章、子级工具目录删 `task.spawn`）、`internal/tools/task.go`（`TaskRoster` **15 个方法**，含 `PublishSubagent:337`/`Descendants:300`/`TryAcquireSubagentSlot:358`/`Cancel:442`）、`internal/panel/pump.go`。非测试命中 **71 行**。
- ⇒ **A397 那句要加时间戳**（"截至 09-28 18:1x"），否则下个读台账的人会以为实体层还没动工、把已完成的格子再开一遍。**票 197 的 A 段与 B 段（名册＋流键形状）看起来已经落码**，C 段（界面）与 197-r3（载体）没有。

**② 他那句"只有 DSH 真做到点进去的子代理页面"——DSH 这一侧我**只能证实"DSH 确实做到了"，**不能证实"只有"**
- 证实的部分（写得很硬）：`ui-subagent` 确实给了**两种落点**（主区导航 **＋** 右侧 Sidebar tab），而且那一页**就是主对话那套组件**（`sidebar-chat/index.tsx:145` 直接 `renderFactorySlot('conversation.content', ...)`）。**这一点比 A401 写的还要多**：A401 只写"页面直接复用主对话那套组件"，**没写"同一份数据还能在右侧再开一格，且按 `(kind,contentId)` 去重不会开重"**。⇒ 抄的时候如果只做主区导航，就少做了他们那一形的一半。
- 不能证实的部分：**"只有"是个跨家断言，另三家不在我这腿射程**。而且我这腿还读到**一处会让人把"只有"读歪的东西**：`ui-workflow-run`（名册第 51 枚，"把顶层工作流运行重建为带嵌套成员折叠的独立聊天节点"）与 `ui-jobs`（"没有保留输出的任务保持静态，**包括回答已交给模型的 subagent**"，`ui-jobs/README.zh.md:11`）——**DSH 内部对 subagent 就有三条不同的展示通路**（名册树／Sidebar 页／jobs 静态行）。"DSH 做了那一页"是真的，但**"那一页在 DSH 里也不是唯一形状"**。

**③ A401 的两处 `文件:行` 我核到偏差（他说过"抽查到我写歪一条就整批判未验证"，所以我主动报）**
- `client/ui-subagent/src/client/index.ts:46-47`（他用来钉"独立于输入框的 Stop"）：**第 46-47 行是一段注释**（`:45-47` 三行注释，`:48` 才是那句 `return owner.session?.running === true ? null : { reason: 'parent-unavailable' }`）。"独立 Stop"这件事的**真出处是 README 叙述**（`ui-subagent/README.zh.md:40`：输入区与 Send 被禁用，**但独立的 Stop 保持可用**）与 Host 侧 `@Remote('interruptByParent')`（`subagent/subagent/src/index.ts:482`，客户端面在 `api/session-controller/src/client/sessions/session.ts:343-350`）。⇒ **结论方向没错，行号要换。**
- `subagent/src/types.ts:250-262`（他用来钉"running/waiting/settled 推导表"）：那一段实际是 **`SubagentStopReasonMap`**（现读 `:252` 是它的 `export interface` 行，成员 `:254-262` 五枚结束原因 `completed/aborted/error/max-tokens/refusal`），**不是 running/waiting/settled 的推导表**。真正的"activity 怎么定"在客户端 `SubagentHeaderLineage.tsx:459-461`（UI status 优先、退摘要），**而"目录不区分失败/取消/拒绝/token 耗尽"恰恰是他们登记的限制**（`README.zh.md:115`）。⇒ **A401 那句"状态从事件推导、不建第二套状态机"仍然成立，但他引的那段原文推不出它。**

**④ 他可能被 DSH 的官方包地图带偏：那张表只有 50 枚，不是 52**
- `packages/client/README.md:29-88` 那张"Package map"里 **ui-* 只列了 50 行**（`grep -oE '\[`ui-[a-z-]+/`\]' | sort -u | wc -l` = 50），**磁盘上有 52 枚**。**缺席的两枚：`ui-settings-account`、`ui-sidebar-terminal`**（`grep -c 'ui-settings-account/'` 与 `'ui-sidebar-terminal/'` 在 README.md 里都 = 0）。
- ⇒ 他若照着那份"官方名册"数覆盖率，会**少算两枚**，而且少掉的正好是"账号页"和"侧栏终端"这两枚——后者对我们有直接意义（**我们面板今天没有终端格**）。

**⑤ "嵌套树的 hover 时序抄不出对象"这句该撤**
- A401 写"抄不出对象的那几格（如嵌套树的 hover 时序）不硬造"。**这一格有原文**：悬停 **150ms** 打开、离开触发器与树后 **120ms** 关闭、点击则钉住直到外部点击或 Esc（`ui-subagent/README.zh.md:34`），点击固定的树同时保留键盘退出（`:34`）。⇒ 这一格**不用自己发明**。

**⑥ 派单自己那句"快捷键……我们完全没有"——我这腿量出来不对**
- 派单原文（`docs/specs/../.scratch/wisp/dispatches/…three-named-uncovered-surfaces.md:30` 之后那段"要问四件事"的第三条）写的是「`ui-shortcuts` 回答"快捷键是否全可改"（**我们完全没有**）」。
- 现读：我们**有 4 枚可改键**（`[hotkey] summon/mute/cancel/panel`，`internal/config/schema.go:179-183`，hot-tier），**有键位串的解析与注册报告**（`internal/ball/hotkey_windows.go:42-47` 的 `HotkeyConfig`、`:264-301` 的 `HotkeyBinding`/`HotkeyReport`），**有热改桥**（`internal/ball/hotkey_reload.go`，`:42-44` 还写清了"为什么必须轮询而不是 `OnReload`"），**甚至有一份 84 枚候选组合在 owner 本机的占用实测**（`hotkey_windows.go:56-60`）。
- 真正没有的是**命令目录＋录键 UI＋冲突校验＋"恢复默认"**，以及**真实宿主里的接线**（`NewHotkeyReloader` 的生产调用只有 `cmd/balldebug/main.go:237` 一处，`cmd/wisp/` 下 0 处）⇒ **〔建了但没接〕，不是〔没有〕。**
- ⇒ 并进差集清单时如果写成"我们完全没有快捷键功能"，owner 会以为要新建一整支；实际是**一条已存在的链缺最后三件配套**。这一格我建议按 `A401` 他自己那套三档证据强度标成 **〔建了但没接〕**。

**⑦ 一处"我们比他们强"的地方，别被差集清单倒过来写**
- 审批卡的字段面我们 **11 枚** vs 他们 **5 枚**（`internal/panel/approval.go:39-59` vs `ui-approval/src/client/contract/slots.ts:52-63`）：我们多 `rulesHit`（**命中了哪几条冻结规则、按评估者顺序**）、`reasonKnown`（**没有理由时不许渲染成"没风险"**）、`sessionOverrideBlocked`（R4 taint）、`callChain`、`decidedBy`（**这一票来自哪一侧**）。
- 层位上我们也多一维：他们的应答面只有"浏览器"（＋TUI 组装），我们有 **4 枚闭集通道** ball/esc/panel/kws（`internal/agent/approval/approval.go:53-61`），其中 `kws` 是**语音否决**——**这一条 DSH 的 52 枚里没有任何一枚对应物**（语音在他们是 Host 侧链路，不进 `ui-*`）。
- ⇒ 报差集时请把这一格报成"我们没有的那一维"，不要报成"他们的审批更完整"。

**⑧ 一处方向性风险，我要说清楚**
- DSH 的面板侧**可以给 `allowed-once`**（`ApprovalDecision = 'allowed-once' | 'rejected'`，`ui-approval/src/client/contract/slots.ts:66`），而**我们 C17 明写 `approval.decide` 的 allow 拒绝一切面板来源（F2）**（`SPEC-08*.md:156-176`）。**这两条很容易被读成"DSH 允许面板批准、我们禁了，所以我们更严"**——**不是同一回事**：DSH 那条通路是"浏览器作为 Host waterfall 的一个应答方"，它的**安全边界在 Host 侧的持久策略**（`README.zh.md:36` 逐字："持久权限策略仍由 Host 侧审批包拥有"）。⇒ 真正该问的是**"我们的 allow 出口该落在哪一层"**（原生卡？球的长按？语音取消词的反义？），而不是"要不要学他们放开面板"。**这一问碰 Q-49（面板侧来源的 L2「允许」零仪器覆盖），我不越权答，只标出来。**

---

*本程只写这一枚文件；未 commit；未碰 `pending-and-issues.md`、未碰别人的证据件、未碰源码；`frontend/**` 与 `design/**` 未读未转述。*
