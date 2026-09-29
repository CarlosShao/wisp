# Wisp 外部对标 · 第五版块 A —— 移动端那一屏（openchamber）

> **本腿＝`v5-block-A`，只做块 A（移动端那一屏）**。块 B（VS Code 扩展那一屏）与块 C（DSH 那几十枚 `ui-*` 分包）由别的腿做，本腿不碰。
> ⛔ 不重抄第四版已有条目（`docs/reports/missing-features-2026-09-29-v4.md`，663 行）；重了就算缺陷。
> ⛔ 我方 `frontend/**` 与 `design/**` 零读零引零枚举；我方判断只引 `docs/**` 行号或本腿现读的 Go 文件行号。
> 本腿不跑任何编译／测试／门禁；不 commit、不 push。
> 骨架先落＝防死在轮次上限；每填完一节即回写本文件。

---

## §0 本轮取数的家与出处（版本／快照先摆明）

| 项 | 值 | 尺／出处 |
|---|---|---|
| 家 | **openchamber**（第五家，2026-09-28 进本机） | 台账 `docs/reports/pending-and-issues.md:8737`（A405） |
| 快照落点 | `D:\work\AI\open source\openchamber`（**在本仓之外**） | 现读 `ls` |
| 取法 | ⛔ 不用 `git clone`（本机撞 schannel），用 `curl` 取主干 tarball `codeload...refs/heads/main`（41,543,271B）解包 | 台账 A405 `:8749-8754` |
| 对应提交 | `80c888eb63c9bd92c96488b13df14520f2e1c68a`（`api.github.com/.../commits/main` 现取，非快照内读史） | 台账 A405 `:8754` |
| `.git` 有无 | **无**（`test -e .git`＝NO） ⇒ 快照零历史、零 node_modules；凡"上游怎么改"一律不写 | 本腿现验 |
| 根版本 | `2.0.3`（`package.json:version`，CHANGELOG 顶节 `[2.0.3] - 2026-09-28`） | 本腿现读 |
| 移动端壳版本 | `@openchamber/mobile` **`1.13.2`**（`packages/mobile/package.json`） | 本腿现读 |
| 共享 UI 版本 | `@openchamber/ui` **`2.0.3`**（`packages/ui/package.json`） | 本腿现读 |
| 本腿地界 | `packages/ui/src/apps/*Mobile*`（移动屏逐控件）＋ 移动复用的 `components/views/ChatView` 与 `components/chat/PermissionCard`／`PermissionDock` 正文 ＋ `components/sections/**`（移动可见的设置页）＋ `lib/i18n/messages/en.ts`（原文文案） | —— |
| `apps/` 目录尺 | `find packages/ui/src/apps -type f`＝**52 枚**，其中带 mobile/widget 字样 **36 枚**（与第四版 §5-4 更正一致，本腿沿用不重议分母） | 本腿现量 |
| 出厂装没装 | 移动壳是**独立发布物**（Capacitor，`@openchamber/mobile` 1.13.2，走 App Store／Play），不随 web/桌面分发；web/桌面装的是 `@openchamber/ui` 2.0.3 | `packages/mobile` 存在＋版本号为证 |

**纪律句**：每一处"别人怎么做"＝点名 openchamber ＋ `文件:行`；转述前 `sed`/Read 逐条验原文；结论标〔已证／建了但没接／仅文档〕；快照无 `.git` ⇒ 禁读历史。凭据/API key 绝不进输出。

---

## §1 表一 · 移动端那一屏 · 逐块逐控件表

> 尺：`packages/ui/src/apps/MobileApp.tsx` 现读全 1,394 行 → 逐层落到各件；行号＝本腿 `sed -n`/Read 现验原文。
> 路径前缀（为省字，下表"层"列内省略）：`UI/apps/` ＝ `packages/ui/src/apps/`；`UI/chat/` ＝ `packages/ui/src/components/chat/`；`en.ts` ＝ `packages/ui/src/lib/i18n/messages/en.ts`；`en.settings.ts` ＝ 同目录 `en.settings.ts`。
> **出厂装没装**：`L0` 连接屏与 `L5/L6` 的 instances 层只在 **Capacitor 原生壳**（`@openchamber/mobile` 1.13.2，另装）出现；web/PWA 走同一套 `L1–L5` 但没有连接屏（`MobileApp.tsx:1270` 以 `isNativeMobileApp` 分叉）。
> 层判定：**〔已证〕**＝代码在产路径上挂着；**〔建了但没接〕**＝字符串/组件在、无渲染者；**〔仅文档〕**＝只在 md 里。

### 1.A 进入页面之前（冷启动 → 连接屏 → 进壳）

| # | 控件（界面原文） | 层（组件文件:行） | 点了/发生什么 | 后端对应物 | 判定 |
|---|---|---|---|---|---|
| A1 | （无文字）启动图 + `Connecting to device:` ＋实例名＋跳动的点 | `UI/apps/MobileApp.tsx:1305-1322`；文案 `en.ts:158` `mobile.connect.splash.connectingTo` | 冷启动静默回连上一次那台；连上就进壳，连不上掉到连接屏 | `autoConnectLastInstance()` `MobileApp.tsx:920-956` | 〔已证〕 |
| A2 | 「Use another server」（`mobile.connect.cancelPassword`） | `MobileApp.tsx:1286-1295`（恢复屏）＋`UI/apps/MobileConnectionWelcome.tsx:189-197` | 清掉 runtime endpoint，回到连接屏重选 | `switchRuntimeEndpoint({apiBaseUrl:''})` `MobileApp.tsx:1290` | 〔已证〕 |
| A3 | 「Unable to reach server」（`sessionAuth.error.networkTitle`）＋原生版说明句 | `MobileApp.tsx:1281,1284`；`en.ts:2785`；原生专用文案 `mobile.connect.recovery.description` `en.ts:145` | 8s（浏览器）/4s（原生）等不到就绪才显出这一屏 | 超时门 `MobileApp.tsx:1210-1226` | 〔已证〕 |
| A4 | 顶部黄条：`Couldn't reach {label}…`／`Access to {label} has expired or was revoked. Sign in again.` | `MobileConnectionWelcome.tsx:146-160`；`en.ts:143-144` | 说明"为什么被踢回连接屏"（两种：不可达／授权被撤） | 回连结果 `MobileApp.tsx:933-937`（`unreachable`/`needs-login`） | 〔已证〕 |
| A5 | 「Scan QR code」大按钮（`mobile.connect.scanQr`） | `MobileConnectionWelcome.tsx:202-218`；`en.ts:146`；**主路径**，手动录入默认收起（`:36-38`） | 拉起全屏取景框→识别→直接发起连接 | `scanConnectionQr()`（MLKit/CameraX，见前一份移动件），成功后 `conn.connect` `:89` | 〔已证〕 |
| A6 | 取景框 ＋ "取消" | `MobileConnectionWelcome.tsx:134` → `UI/apps/MobileQrScannerOverlay.tsx` | 中止识别回连接屏 | `AbortController` `:120` | 〔已证〕 |
| A7 | 扫描失败的四种说法（权限没开／码不对／不是本家的码／这台设备不支持） | `MobileConnectionWelcome.tsx:95-106`；`en.ts:150-153` | 落到 `error` 行内提示 | 同文件 status 分支 | 〔已证〕 |
| A8 | 「Saved connections」列表每一行（服务器图标＋名字＋地址或 `via OpenChamber Relay`） | `MobileConnectionWelcome.tsx:222-261`；`en.ts:154,156` | 点＝连这一台，行尾转圈＋文字变 `Connecting...` | `conn.connect({id,candidates,clientToken,label})` `:238` | 〔已证〕 |
| A9 | 「Connect by address」折叠开关（`mobile.connect.manual.toggle`） | `MobileConnectionWelcome.tsx:265-275`；`en.ts:149` | 展开手填三格 | 纯前端 `:38` | 〔已证〕 |
| A10 | 手填三格：`Server URL`／名字／`Client token`（＋提示 `Only needed if…`） | `MobileConnectionWelcome.tsx:282-316`；`en.ts:132-136,297-298` | 支持把 `openchamber://…` 配对链接粘进地址格自动拆回三段 | `parseConnectionPayload` `:54-69`；`conn.redeemPairingConnection` | 〔已证〕 |
| A11 | 「Connect」提交钮（`mobile.connect.connectButton`） | `MobileConnectionWelcome.tsx:318-320`；`en.ts:139` | 发起连接 | `conn.connect` | 〔已证〕 |
| A12 | 密码那一跳：锁图标卡＋`Password` 格＋「Unlock and connect」 | `MobileConnectionWelcome.tsx:162-198`；`en.ts:137-140` | 令牌被拒时先解锁再连 | `conn.submitPassword` `:122-125` | 〔已证〕 |
| A13 | **隐藏诊断**：长按 logo | `MobileConnectionWelcome.tsx:40-45,140` → `MobileConnectionDebugPanel`；文案 `en.ts:178-182`（`Connection log`／`Copy`／`Close`／`No connection events yet.`） | 打开连接事件日志面板（被踢回连接屏也进得去） | 事件本 `logMobileConnectEvent`（`MobileApp.tsx:743,747,762,886,971` 等 8 处实调用） | 〔已证〕 |

### 1.B 顶栏（进壳后常驻，`UI/apps/MobileHeader.tsx`，全 182 行本腿读穿）

| # | 控件 | 层 | 点了发生什么 | 后端对应物 | 判定 |
|---|---|---|---|---|---|
| B1 | 左：三横线钮（`list-unordered`） | `MobileHeader.tsx:94-102`，aria `mobile.sessions.openSheetAria` | 手机＝拉出会话抽屉；平板＝开关常驻侧栏 | 纯前端 `MobileApp.tsx:504` | 〔已证〕 |
| B2 | 中：会话标题＝**下拉触发器**（带会翻转的 v 形箭头） | `MobileHeader.tsx:105-130`（箭头 `:122-128`）；无标题时 `mobile.sessions.untitled`，草稿屏 `sessions.switcher.draftTitle` `:55-57` | 开「最近会话」浮层 | `MobileSessionSwitcher` `:175-179`；行内流式态 `sessions.sidebar.session.status.active/unread` `MobileSessionSwitcher.tsx:64` | 〔已证〕 |
| B3 | 右1：**上下文环形进度**按钮（会话元数据） | `MobileHeader.tsx:136-144` → `MobileSessionMetadata.tsx:44-90`（18px 环、按 0-100 上色）；aria `mobile.header.openMetadataAria` `en.ts:188` | 弹出元数据面板：`Context`／`Usage`／时间格式 | `findLatestContextFill`、`useQuotaStore` `MobileSessionMetadata.tsx:12-14` | 〔已证〕 |
| B4 | （多跑概览开着时）B3 **整枚消失**，标题改说这一跑的名字 | `MobileHeader.tsx:50-51,136`；`MobileApp.tsx:133-141` | 概览盖住聊天时不描述"不在屏上的会话" | `useMultiRunTitle` | 〔已证〕 |
| B5 | 右2：空间准入钮（`SpaceAccessButton`） | `MobileHeader.tsx:146-150` | 隔离空间开才挂（`MobileApp.tsx:1374`） | `SpaceAccessDialog` | 〔已证，默认关：`isolatedSpacesEnabled` 门〕 |
| B6 | 右3：工作区钮（`pencil-ruler-2`），**有未提交改动时右上角一颗琥珀点** | `MobileHeader.tsx:152-172`（点 `:166-171`）；aria 两变体 `en.ts:186-187` | 拉出工作区抽屉 | `useGitStore` 脏判定 `:41-45` | 〔已证〕 |

### 1.C 会话抽屉／常驻侧栏（`UI/apps/MobileSessionsSheet.tsx`，2,642 行，本腿按控件点名）

| # | 控件（原文） | 层 | 行为 | 判定 |
|---|---|---|---|---|
| C1 | 搜索框（复用桌面 `SessionSearchInput`） | `MobileSessionsSheet.tsx:41,403,464,1534`；结果行下面加一行 `Project · branch` `:1534` | 搜会话，命中行两行高 | 〔已证〕 |
| C2 | `View` 切换：`Grouped`／`Timeline` | `MobileSessionsSheet.tsx:150-151,2386`；`en.ts:210-212` | 换分组视图／时间线视图 | 〔已证〕 |
| C3 | 分区标题：`Projects`／`Chats`／`Recent`／`Worktrees`／`Switch project` | `en.ts:205-209` | 只是分区；`Recent` 是移动抽屉新增（`CHANGELOG.md:15`） | 〔已证〕 |
| C4 | 日期分组：`Today`／`Yesterday`／`Earlier this week`／`Older` | `en.ts:201-204` | 时间线视图的分组 | 〔已证〕 |
| C5 | `Show archived ({count})`／`Hide archived` | `en.ts:219-220` | 展开归档 | 〔已证〕 |
| C6 | `New chat in {project}`／`Start new chat`／`New chat` | `en.ts:200,222-223` | 起新会话（进草稿屏） | 〔已证〕 |
| C7 | **行左滑**露出的动作槽：archive→rename→**Copy ID**→delete 顺序 | `UI/apps/MobileSessionSwipe.tsx:12`（"Four slots…archive, delete, rename and Copy ID"）、`:114-117`（**故意的排序**：archive 打头，破坏性的 delete 绝不落在半拉的行程下；Copy ID 压尾）、`:61,88`（过半才露出） | 拖拽露出＋点触 | 〔已证〕 |
| C8 | 钉住／取消钉住（pin 动作） | `MobileSessionSwipe.tsx:17` `MobileSessionPinAction` | 钉住置顶 | 〔已证，`CHANGELOG.md:15` 记为移动抽屉新件〕 |
| C9 | `Reorder projects`／`Done` ＋拖拽把手＋上下键＋移除 | `MobileSessionsSheet.tsx:1718,1880`；`en.ts:224-232`；**引导句逐字** `en.ts:226` | 进入编辑态重排项目/工作树 | 〔已证〕 |
| C10 | 抽屉底部条：当前实例名／`Instances`／`Usage`／`Settings`／`Update` | `MobileApp.tsx:384-393`（`sessionsFooter`） | 分别开 `L6` 的全屏层；`Update` 只在**非**原生壳出现 `MobileApp.tsx:374-376` | 〔已证〕 |
| C11 | 空态三句：`No projects yet`／`No sessions yet`／`No matches` | `en.ts:213-218` | 每枚空态都是句子 | 〔已证〕 |
| C12 | 平板：侧栏左右两条**拖宽把手** | `MobileApp.tsx:491-498,574-581`；`en.ts` `sidebar.resize.leftPanelAria/rightPanelAria` | 改宽 | 〔已证〕 |

### 1.D 聊天正文与"要人处置"那一族浮层（复用桌面件，`ChatView`→`ChatContainer`→`ChatInput`）

| # | 控件 | 层 | 行为 | 判定 |
|---|---|---|---|---|
| D1 | 审批坞（`Permission required` 标题＋可收起＋**每枚请求一颗点**、当前那颗实心） | `UI/chat/PermissionDock.tsx:85-137`（收起钮 `:90-98`，点阵 `:112-132`），挂点 `UI/chat/ChatInput.tsx:4161-4165` | 一次答一枚，答完下一枚顶上 `:44-48`；批起来的请求读作"一坞走一遍"而不是"一摞卡" `:16-25` | 〔已证〕 |
| D2 | 坞上的三枚动作钮 | `UI/chat/PermissionCard.tsx:373-393`（dock 变体）／`:405-450`（内联变体） | 见 §5 | 〔已证〕 |
| D3 | `From subagent` 小徽标 | `PermissionDock.tsx:100-104`；`en.ts:2324` | 这条请求来自子会话 | 〔已证〕（判据 `usePermissionResponse.ts:10-18` 按 `parentID`） |
| D4 | `Request {index}: {tool}` 步序文字＋`{current} of {total}` | `PermissionDock.tsx:105-109`；`en.ts:3369`（stepAria）＋`en.ts:3347`（`'{current} of {total}'`，**复用表单坞那一枚**） | 批量请求走到第几枚 | 〔已证〕 |
| D5 | 安全网拦下时卡顶一行琥珀说明 `The safety net held this action for you to decide · <风险类>` | `PermissionCard.tsx:300-308`；`routing.i18n.ts:5`；七个风险类名 `PermissionCard.tsx:55-63` | 告诉人"为什么这次要问" | 〔已证〕 |
| D6 | 子代理请求的**内联卡**（`PermissionCard`，非坞） | `PermissionCard.tsx:453-495`；注释逐字 "The inline card the BTW sheet shows for its child session's requests" | 在 BTW 分身层里逐条答 | 〔已证〕 |
| D7 | 「Parent」返回钮（左上角悬浮，箭头＋`Parent`） | `ChatContainer.tsx:1008-1022`；`en.ts:2285-2288` | 回到父会话 | 〔已证〕 |
| D8 | 排队消息条（`QueuedMessageChips`）／表单坞（`FormDock`） | `ChatInput.tsx:4166-4175` | **优先级写死**：BTW > 审批 > 表单 > 队列 > 建议（注释逐字 `:4159-4160`） | 〔已证〕 |

### 1.E 输入区（**移动专用分支**，`ChatInput.tsx` 的 `isMobile` 面）

| # | 控件 | 层 | 行为 | 判定 |
|---|---|---|---|---|
| E1 | 模型钮（logo＋名，无 `>` 箭头） | `ChatInput.tsx:3642-3652` → `UI/chat/MobileModelButton.tsx:23-24`；`autoModel`／`selectModel` | 点开移动面板选模型 | 〔已证〕 |
| E2 | 智能体钮：**单击换下一个、按住 500ms 开面板** | `MobileAgentButton.tsx:20`（`LONG_PRESS_MS=500`）、`:39-50`（用 pointer 事件防键盘收起）；挂载 `ChatInput.tsx:3647-3651` | 两用手势 | 〔已证〕 |
| E3 | 附件钮（点开**动作面板**而不是菜单） | `ComposerFooter.tsx:151-165` ＋ `ChatInput.tsx:3658-3666`（先标开再 blur，专门为了键盘收起时不塌） | 选文件/issue/PR/Linear | 〔已证〕 |
| E4 | **盾牌钮**（会话权限档循环：ask→safety→auto） | `ComposerFooter.tsx:166-171`（在 `composer-mobile-actions` 这个移动分支里，分支起点 `:147`）→ `PermissionAutoAcceptButton.tsx:30-34`（三枚不同盾牌图标）；循环逻辑 `UI/chat/permissionAutoAccept.ts:15-42` | 草稿上改草稿档、会话上落会话档；没有会话时先逼你开会话 `permissionAutoAccept.ts:29-37` | 〔已证〕，详见 §5 |
| E5 | 会话目标钮＋目标计数器 | `ComposerFooter.tsx:172-178`（`SessionGoalButton`/`SessionGoalObjectiveCounter`） | 设/看本轮目标 | 〔已证〕 |
| E6 | 话筒钮（听写开关） | `ComposerFooter.tsx:185-203`；`chat.dictation.start` | 全局事件驱动 `ChatInput.tsx:3654-3656` | 〔已证〕 |
| E7 | **发送／排队／停止**三态主钮 | `ComposerActionButtons.tsx:69`（`Send message`）／`:97`（`Queue message`）／`:109`（`Stop generating`）；三句 `en.ts:2406-2408`；挂载 `ComposerFooter.tsx:205-214` | 跑着＝排队＋停止双钮，闲时＝发送 | 〔已证〕 |
| E8 | 输入框**最多几行**在手机上被单独收紧，且随"绑住的那枚元素"重新定位 | `ChatInput.tsx:2535-2537`（`MAX_MOBILE_COMPOSER_LINES`、`[data-composer-bound]`） | 布局约束 | 〔已证〕 |
| E9 | 把手上拉＝**全屏输入模式** | `ChatContainer.tsx:940-942` 注释逐字 "this now covers mobile too: the mobile composer enters the same fullscreen-input mode via its drag handle" | 长需求全屏写 | 〔已证〕（`.tsx` 内 drag handle 具体挂载点本腿未逐行读，见 §8） |
| E10 | **移动行评论输入层**（`MobileCommentComposer`） | `ChatInput.tsx:184-186,1065-1068`（`isMobile && status==='open'` 时接管发送） | 指着某一处说话 | 〔已证〕 |
| E11 | 键盘弹起时把表单钉在可视视口 | `ChatInput.tsx:3668-3672` `useMobileViewportPin`；注释逐字"移动浏览器是平移可视区而不是重排布局" | 防遮挡 | 〔已证〕 |
| E12 | 草稿屏 starter 快捷指令条 | `MobileApp.tsx:269-279`（键盘起来时是否留住的判据：平板竖屏或接了实体键盘）；设置键 `Show Starters on New Session Screen`（`en.settings.ts:2137`） | 首屏可点 | 〔已证〕 |
| E13 | 桌面版才有的那套"草稿呈现层"（`DraftPresetChips`）在手机上**不挂** | `ChatInput.tsx:3504-3506`（`&& !isMobile`）→ 挂载点 `:4153-4157` | 同一功能两条面 | 〔已证〕（形状级风险，见 §7） |

### 1.F 工作区抽屉（`UI/apps/MobileWorkspaceDrawer.tsx`，全 318 行读穿）

| # | 控件 | 层 | 行为 | 判定 |
|---|---|---|---|---|
| F1 | 五枚标签：`Changes`／`Files`／`Terminal`／`Notes`／`MCP`（**可拖序**的 `SortableTabsStrip`；五枚放不下＝只有当前枚带字，其余塌成图标） | `MobileWorkspaceDrawer.tsx:197-224`（`:220-223` 那条取舍写死）；`mobile.menu.*` | 切页 | 〔已证〕 |
| F2 | 抽屉右上 X（`Close`） | `:227-235`；`en.ts:185` | 关 | 〔已证〕 |
| F3 | **看过的页不卸载**（重开落在原处：开着的 diff、编辑一半的文件、挂着的终端） | `:141-151,238-283` | 状态保活 | 〔已证〕 |
| F4 | MCP 页：`add`＋`refresh` 两枚（刷新**最少转 500ms**，免得闪一下看不出来） | `:55-77`（`:45` 那个 `minSpinPromise`） | 加服务器→跳设置页预填草稿（`MobileApp.tsx:395-431`）／刷新状态 | 〔已证〕 |
| F5 | 左边缘回滑＝关抽屉／Esc＝关（终端页占用 Esc 时不抢、Vim 编辑态也不抢） | `:132-139`、`:180-188` | 手势 | 〔已证〕 |
| F6 | 改动页可直接开进"某个文件的 diff" | `MobileApp.tsx:151-152,221-225`；`:245` 用 path＋staged 当 key 重挂 | 从卡片跳进来 | 〔已证〕 |

### 1.G 全屏层与手势（壳层）

| # | 控件 | 层 | 行为 | 判定 |
|---|---|---|---|---|
| G1 | 左边缘滑＝会话抽屉、右边缘滑＝工作区抽屉 | `MobileApp.tsx:325-334`（`useEdgeSwipe`） | 单手换层 | 〔已证〕 |
| G2 | **安卓返回键的分层栈**（计划→全屏层→工作区→会话抽屉→才退出）；设置页自己那套下钻先逐级返回 | `MobileApp.tsx:336-369`（`:351` 先问 `settingsBackRef`） | 系统手势语义 | 〔已证〕 |
| G3 | 实例管理全屏层（`Instances`） | `MobileApp.tsx:605-619`；**仅原生壳**（`showCapacitorOnlyFeatures` `:157`） | 多服务器管理 | 〔已证〕 |
| G4 | 用量全屏层 `Usage`／设置全屏层 `Settings`／更新全屏层 `Update` | `MobileApp.tsx:621-675`；平板上这四层改成**居中对话框**而不是全屏（`surfaceVariant` `:380`） | 同四层 | 〔已证〕 |
| G5 | 冷启动"续上次的会话"期间**用 logo 盖住**中间那一闪（6s 安全阀） | `MobileApp.tsx:1037-1098`（`:1066` 超时阀） | 不闪草稿 | 〔已证〕 |
| G6 | 网络切换重探（安卓切 Wi-Fi **不后台化**，所以靠 `online` 事件补一次）＋**宽限期阶梯 4s→10s** 再判死 | `MobileApp.tsx:831-852`、`:780-814` | 免误判断连 | 〔已证〕 |
| G7 | 实体键盘判定 → 决定首屏快捷指令留不留 | `MobileApp.tsx:192,269-279`（`useHardwareKeyboard`） | 平板/折叠屏 | 〔已证〕 |
| G8 | 折叠屏＝**按尺寸类不按机型**（展开＝平板，合上＝手机，壳不重启） | `MobileApp.tsx:183-203`（注释逐字 "A SIZE class, not a device check"） | 布局分叉 | 〔已证〕 |

---

## §2 表二 · 移动端设置页 · 逐项表（标签原文＋改了什么行为＋我方对应物）

### 2.0 先说清"移动能改哪几页"——这一格本身就是别人替我们量好的数

**手机不是"把桌面设置整页塞进小屏"，是一份写死的白名单**：`UI/apps/MobileApp.tsx:88-111` 的 `MOBILE_SETTINGS_PAGES` 常量逐名列出 **22 页**，装配点在 `MobileApp.tsx:638`（`visiblePageSlugs={[...MOBILE_SETTINGS_PAGES]}`），过滤器在 `components/views/SettingsView.tsx:246-255`（先按白名单裁、再按 `isAvailable(ctx)` 裁、**最后再单独把 `shortcuts` 在手机上摘掉** `:254`）。

| 处置 | 页（slug ＋ 原文标题） | 证据 |
|---|---|---|
| **手机有**（22） | `general` General · `appearance` Appearance · `chat` Chat · `notifications` Notifications · `sessions` Sessions · `routing` Routing · `git` Git · `magic-prompts` Magic Prompts · `snippets` Snippets · `behavior` Behavior · `agents` Agents · `commands` Commands · `mcp` MCP · `plugins` Plugins · `skills.installed` Skills · `skills.catalog` Skills Catalog · `providers` Providers · `web-search` Web search · `usage` Usage · `voice` Voice · `integrations` Integrations · `about` **About** | 白名单 `MobileApp.tsx:88-111`；标题 `lib/settings/metadata.ts:69-226` |
| **手机没有**（5） | `home`（搜索页，永远不给任何端）· `projects` · `remote-instances` · `shortcuts` · `tunnel` · `extensions` | `SettingsView.tsx:249`（`!== 'home'`）／`:254`（shortcuts 摘掉）／白名单不含 projects·remote-instances·tunnel／`metadata.ts:233` `extensions` 的 `isAvailable` 逐字写着 `!ctx.isVSCode && !ctx.isMobile` |
| **只在手机上才有的页** | `about`（`metadata.ts:225` 的 `isAvailable` 逐字 `ctx.isMobile && !ctx.isVSCode`）——桌面/网页端**根本没有这一页**，移动壳里它是"查并装服务器更新"的入口（`MobileApp.tsx:636-637` 注释逐字写明这个分工） | 〔已证〕 |

> **条目总数口径（本腿现量，别当"读穿了"）**：全仓设置项注册表 `lib/settings/search.ts` 里可解析出 **181 枚条目**（尺：`grep -c "^    id: '"`＝181，再按 `page:` 配对本腿现算）；其中落在上面 22 页内的 **158 枚**（页外 23 枚＝`projects` 11／`tunnel` 6／`extensions` 3／`remote-instances` 2／`shortcuts` 1）。**本腿只逐枚核了与移动形态直接相关的 12 枚**（下表），其余按页登记，未逐条读正文（见 §8）。

### 2.1 与移动形态直接相关的那几枚（逐枚，标签＝原文）

| 设置项（标签原文） | 页 | 改了会得到什么 | 移动专属？ | 我方对应物 | 证据 |
|---|---|---|---|---|---|
| **`Mobile Keyboard Behavior`**（提示原文：`Default browser behavior is safest. Resize content asks supported browsers to shrink the app when the on-screen keyboard opens.`） | Appearance | 换键盘弹起时整页是缩还是平移 | **只在"手机＋网页端＋非桌面"出现**（`ctx.isMobile && ctx.isWeb && !ctx.isDesktop && !ctx.isVSCode`） | **没有** | `lib/settings/search.ts:115-120`；标签 `en.settings.ts:2059-2060` |
| **`Input Bar Offset`**（提示：`Raise input bar to avoid OS-level screen obstructions like home bars.`） | Appearance | 把输入条整体上抬，躲开全面屏底部横条 | **`isAvailable: (ctx) => ctx.isMobile`＝手机独有**；注释逐字 "Only the mobile composer applies this offset" | **没有** | `search.ts:163-169`；`en.settings.ts:2084-2085` |
| **`Install App Name`** ＋ **`Install Orientation`** | Appearance | 装成应用时的名字与锁定的屏幕方向（改了提示"要重装 PWA"） | 面向"装到手机上"这一形态 | **没有**（我们不装到别的设备） | `search.ts:99-112`（两枚 id：`appearance.pwa-install-name`／`appearance.pwa-orientation`）；标签 `en.settings.ts:2050,2054` |
| **`Show sessions as tabs in the header`**（`Session tabs`） | General | 关掉的分支逐字写着"回到朴素的会话标题" | **手机上被摘**（`!ctx.isMobile`）＝同一功能两端不同形 | **没有** | `search.ts:185-190`；`en.settings.ts:2089-2092` |
| **`Terminal Quick Keys`**（提示：`Show Esc, Ctrl, Arrows in terminal view`） | General | 终端上方一排快捷键 | **`!ctx.isMobile`＝手机上被摘**（但工作区抽屉里有 Terminal 页，见 §1.F1） | **没有** | `search.ts:193-198`；`en.settings.ts:2088,2093` |
| **`Allow Prompting Subagent Sessions`** | Chat | 点进子代理后能不能**当场往里说话**；默认只读 | 移动与桌面共用，但手机上"钻进去"是唯一入口（§4） | **没有** | `search.ts`（`chat.subagent-read-only-banner`）；标签 `en.settings.ts:2136`；执行侧 `ChatContainer.tsx:1024-1028` `resolveChatPromptReadOnly(...allowPromptingSubagentSessions...)` |
| **`Permissions`**（三档芯片：`Ask every time`／`Safety net`／`Accept everything`） | Sessions | 新会话的**默认**权限档；与聊天框那枚盾牌同一套档位（§5） | 共面，但移动壳里这是唯一能改默认的地方 | **半个**：我们有 `risk.permission_mode`（默认 `ask_every_step`，`internal/config/schema.go:470`，本腿现读），**但没有那一屏**——我们面板里根本没有设置页（尺：`grep -rn "config.html\|settings" internal/panel/*.go`＝0 命中） | `PermissionDefaultModeField.tsx:109-119`；**标签在 `routing.i18n.ts:65,67-69`**（`'settings.sessions.permissions.defaultMode': 'Permissions'` ＋ 三枚芯片 `mode.ask/safety/auto`），**不在 `en.settings.ts`** |
| **`Notification Delivery`／`Notification Events`／`Background Push Notifications`**（事件四枚：`Agent Completion`／`Subagent Completion`／`Agent Errors`／`Agent Questions`） | Notifications | 总开关＋"正在用要不要吵"＋**逐类事件开关**＋**后台推送单独一开关** | 移动壳上推送那一节是**真通道**（APNs），不是浏览器通知 | **没有**：`internal/` 下 23 个包（本腿现数 `ls internal/`）**无 notify/推送类包**；PLAN 里"通知"全部指 Windows 侧（如 `PLAN.md:272`「必须同时更新悬浮球角标」） | `search.ts`（notifications 三枚）；标签 `en.settings.ts:1835-1853,1886-1888` |
| **`Generate Session Recap`／`Enable Session Goals`／`Send Starters on New Session Screen`／`Wide Chat Layout`／`Enable Spellcheck in Text Inputs`／`Send shortcut`** | Chat | 摘要／目标／首屏快捷指令／宽栏／拼写检查／发送键语义 | **最后两枚在手机上被摘**（`chat.spellcheck`、`chat.enter-to-send` 的 `isAvailable` 都是 `!ctx.isMobile`）——手机没有"回车发送"这回事 | **没有** | `search.ts:397-401` 与 `:411-416`；标签 `en.settings.ts:2137,2129,2109,2151,2154` |
| **`Interface Font Size`** | Appearance | 界面字号 | **`!ctx.isMobile`＝手机上不给调**（但终端/编辑器字号没这限制） | **没有**（v4 已记：我们没有任何字号键，本腿不复算） | `search.ts:130-134` |
| **`Sessions in work`** ＋ `Show sessions in work` ＋ `Move sessions into work automatically`（说明逐字：`Jev reads each message you send and moves the session into work when you ask for a change, report a bug, or discuss a concrete change. Questions and research don't count. Jev never marks work done…`） | Sessions | 侧栏里给"真在改东西"的会话单开一区，**旁路小模型判定、且明说"永不替你判完"** | 移动抽屉里就是那个 `In work` 区（`CHANGELOG.md:9`） | **没有**（我们无"会话分区"这一形；`internal/session/` 本腿现量**只有 `doc.go`**） | `en.settings.ts:1015-1019` |
| **`Send anonymous usage reports`** | General | 匿名用量上报 | 共面 | **没有**（`internal/config/schema.go` 无遥测键；本腿按"有没有"问，未逐键扫） | `search.ts`（`appearance.usage-reports`）；标签 `en.settings.ts:2164` |

> ⚠ **两层设置别混成一栏**（沿用上一轮的教训，本腿在移动侧再验一遍）：聊天框里那枚盾牌（§1.E4）是**会话级即时档**，Settings → Sessions 那三枚芯片（本节第 7 行）是**新会话默认值**；两者共用同一套 `PermissionMode` 类型（`PermissionAutoAcceptButton.tsx:18` 与 `PermissionDefaultModeField.tsx:117-119` 引用同一枚举），但**写入面不同**。

---

## §3 表三 · 状态与文案表（外部原句 → 我方对应或"没有"）

> 外部一侧全部是**逐字原句**（引号内即界面文案），行号＝`packages/ui/src/lib/i18n/messages/` 下 `en.ts` 或 `routing.i18n.ts`（后者是 Jev／安全网那一族文案的单独注册表）。
> 我方一侧只引 `docs/**` 行号或本腿现读的 Go 行号；**"我方界面实际印什么字"这一维本腿不读**（前端地界），凡只到"状态名在表上"为止的，都写清到这一层为止。

| 状态类 | 外部原句（逐字）＋行号 | 什么时候出现 | 我方对应 |
|---|---|---|---|
| **连接中** | `"Connecting..."` `en.ts:142`；闪屏上另有一句 `"Connecting to device:"` ＋实例名 ＋跳动点 `en.ts:158` | 点连接／冷启动自动回连 | 无对位状态（我们不连远端）。最接近的是下载：状态名 `Downloading`（`internal/statemachine/states.go:24`）〔已证＝名字在表上〕 |
| **重连中（不清屏）** | `"Reconnecting…"` `en.ts:29`；**实现层的决定更值得抄**：注释逐字 "`isConnected` is a LIVE flag that flips false on every transient SSE/WS drop… we must NOT blank the whole app to a loader" `MobileApp.tsx:1247-1252` | 网络抖一下 | **没有**这一形。我们有 `NoNetwork` 状态名（`states.go:26`），但"掉线时界面**不许**清空"这条判据本腿在 `docs/**` 里没找到对应文字〔未量到〕 |
| **离线／够不着** | `"Could not reach that OpenChamber server."` `en.ts:160`；黄条版 `"Couldn't reach {label}. Check that the server is running."` `en.ts:143`；整屏版标题 `"Unable to reach server"` `en.ts:2785` ＋**原生专用**说明 `"Could not connect to the saved server. Check that it is running, or pick another instance."` `en.ts:145` | 探活失败 | 有状态名 `NoNetwork`（`states.go:26`）；**没有人话句子**（前端地界，不读） |
| **需要凭据** | `"This server needs a password or client token."` `en.ts:161`；按钮 `"Unlock and connect"` `en.ts:140`；占位 `"OpenChamber password"` `en.ts:138` | 令牌没带／要解锁 | **没有**（我们是本机 DPAPI，`internal/secret/` 在仓） |
| **授权过期／被撤** | `"Access to {label} has expired or was revoked. Sign in again."` `en.ts:144` | 探活明确返回 `needs-login` 才敢这么说；注释逐字 "tell the user why they land back on the connect screen **instead of silently bouncing them**" `MobileApp.tsx:773-775` | **没有**对位（我们无跨设备令牌）。同族判语值得抄：**掉回起点必须说为什么** |
| **被拒（密码错）** | `"Could not unlock that server. Check the password."` `en.ts:162`；`"Incorrect password. Try again."` `en.ts:2791` | 解锁失败 | **没有** |
| **被限流** | `"Too many attempts"` ＋ `"Please wait {minutes} minute before trying again."`（复数另一条） `en.ts:2784,2786-2787` | 连错太多次 | **没有**（本腿在 `docs/**` 未量到"尝试次数上限"这一形〔未量〕） |
| **排队（人话回执，含失败）** | `"Queued messages"`／`"edit"`／`"send"`／`"Remove from queue"`／`"Drag to reorder"` `en.ts:2278-2282`；**失败那条**：`"Couldn't queue the message. It's back in the composer."` `en.ts:2283` | 跑着的时候再发一条 | 有状态名 `Queued`（`states.go:28`）；**"没排进去、话退回输入框"这一句的回执形状＝没有** |
| **等批准** | `"Permission required"` `en.ts:3354`（卡标题）；侧栏同名一枚 `en.ts:619`；子代理来源徽标 `"From subagent"` `en.ts:2324` | 审批坞／内联卡 | 有**两枚**状态名 `Confirming`（L1 快确认）／`AwaitingApproval`（L2 队列）（`states.go:21-22`）；球上还有"队列深度角标"的硬要求（`PLAN.md:1327,1334`） |
| **被安全网拦下（人话）** | `"The safety net held this action for you to decide"` ＋ `· <风险类>` `routing.i18n.ts:5`；七个风险类逐名 `PermissionCard.tsx:55-63` | 自动接受开着但被安全网扣住 | **没有**对位句子（我们有 `internal/risk/`，档位判定在，"拦下时要说这一句"没有） |
| **降级（判不了怎么办）** | `"The safety net could not check an action, so it is waiting for you"` `routing.i18n.ts:14` | 分类器跑不起来／超时 | **方向与我方同源且更硬**：v2 已记过"明说判不了就拦、绝不因为查不到就放行"；本腿在移动侧再证一枚实例。我方 `internal/risk/mode.go:44-56` 零值落最严档〔已证在产码里〕，**但没有这句人话** |
| **子代理名册四态** | `"is working"` `en.ts:3260`／`"needs permission"` `:3261`／`"asked a question"` `:3262`／`"Done"` `:3248` | 工作状态面板的子代理区 | **没有**（我们的子代理名册一行给不出花费与计时——v4 §块3 已记，本腿不重抄） |
| **上下文用量"还不知道"** | `"Context usage appears once the session starts."` `en.ts:190`；未知值逐字是一个全角破折号 `UNKNOWN_VALUE = '\u2014'` `MobileSessionMetadata.tsx:41`；**第三种状态**：压缩过之后"填充量未知"→画一条空轨道 `ContextDisplay` 的 `state:'compacted'` `MobileSessionMetadata.tsx:33-39` | 元数据面板 | **没有**（我方有 token 数字面，但"压缩后填充量改成未知态、宁可不画"这一形本腿在 `docs/**` 未量到〔未量〕） |
| **会话读不出来** | `"Session could not be loaded"` ＋ `"The conversation could not be fetched — the server may be offline or unreachable. Nothing is lost; retry once it is back."` ＋ `"Try again"` `en.ts:2303-2306`；授权过期另有专句 `:2305` | 拉历史失败 | **没有**（"Nothing is lost"这种**先安抚再给出口**的三段句形我们没有） |
| **超时** | 审批无超时文案（坞上的三枚钮只有 `isResponding` 转圈 `PermissionCard.tsx:444-448`）；**服务端**把超时算成拒绝在文档里（本腿不据 md 写断言） | 等回答 | 我方有键 `confirm_timeout_sec` 默认 300（`internal/config/schema.go:450`，v4 已引）＋ **D43 转移表**"超时 → 判拒绝"（`PLAN.md:1327`）；**没有**面向人的超时句子 |
| **两端版本不匹配（拒绝启动时也说人话）** | `"OpenCode v2 required"`／`"Your connected server is running OpenCode {version}. Update it to OpenCode {minimum} or newer, then reconnect."`／`"OpenCode is bundled with OpenChamber. Update OpenChamber to get OpenCode v2."` `en.ts:17-23` | 移动壳连上老服务器 | **没有**对位门（`OpenCodeCompatibilityGate` 包在整个 App 外层，`MobileApp.tsx:1,1393`）。我方 `internal/models/` 有版本目录但无"连不上就说这句"的门〔未量到同形〕 |
| **压缩中／完成／失败** | `"Compacting the conversation"`／`"Conversation compacted"`／`"Compaction failed"` `en.ts:3375-3377` | 上下文压缩 | 有状态名 `Settling`（`states.go:23`）〔已证＝名字在表上〕；三句人话**没有** |
| **空态（每一枚都写成句子）** | `"No projects yet"`＋`"Add a project to start chatting with your code."`；`"No sessions yet"`＋`"Start your first chat to see it here."`；`"No matches"`＋`"Try a different search term."` `en.ts:213-218`；`"No saved connections yet."` `en.ts:155`；改动页占位 `"Working-tree review, sync, and commit actions will live here."` `en.ts:292` | 任何没内容的屏 | **没有**这一屏族（前端地界不读；此处只登记"别人每枚空态带下一步"这一形状） |
| **相机权限没开** | `"Camera access is off. Enable it in Settings to scan a QR code."` `en.ts:150` | 扫码前 | **没有**（我们无扫码配对） |
| **诚实的"这条通道今天不可用"（我方独有，反向登记）** | — | — | `internal/agent/approval/approval.go:103-114` 四条逐字：`"语音取消不可用"`／`"面板取消不可用（票 37 未接入）"`／`"悬浮球取消不可用"`／`"Esc 取消不可用"`〔本腿现读，已证在产码里〕。**移动这一侧没有对位物**：它把所有"不可用"都藏在门后面，不写出来——这一格我们不必向它对齐 |

---

## §4 必答一：移动端能不能点进某个子代理，看它各自的流式工作页面

# 答案：**有**〔已证〕——但入口只有一条，且进去默认**只读**。

| 断言 | 证据（全部本腿现读 `sed`/Read 验过原文） |
|---|---|
| 聊天记录里那一枚"派活"的工具调用，渲染成一行**可点**的按钮，原文文案是 `"Open {type} subtask"`（`{type}` ＝子代理名，如 `Open Explore subtask`） | 按钮：`UI/chat/message/parts/ToolPart.tsx:1071-1083`（`<button … onClick={handleOpenSession}>` ＋ `external-link` 图标 ＋ `t('chat.toolPart.openSubtask', …)`）；文案：`en.ts:2522` `'chat.toolPart.openSubtask': 'Open {type} subtask'` |
| **手机上点它就是"就地换成那个子会话"**（不是开侧栏标签）——移动分支被点名写在代码里 | `ToolPart.tsx:1022-1031`：`if (isEmbeddedSessionChat() \|\| isMobile \|\| runtime?.runtime.isVSCode) { setCurrentSession(sessionId, currentDirectory); return; }`，注释逐字 "In contexts with no ContextPanel … **or single-surface layouts (mobile, VS Code), navigate in place.**"；子会话 ID 的来源 `ToolPart.tsx:1902-1914`（`taskSessionId`） |
| 进去之后是**同一个完整 `ChatView`**（所以流式正文、工具行、审批坞一并都在），不是阉割版只读预览页 | 移动壳只有一个 `ChatView`：`UI/apps/MobileApp.tsx:11`（import）＋`:511`（挂载），它按 `currentSessionId` 取会话（`ChatContainer.tsx:936-938`），换 ID 就换页 |
| 进去之后**左上角浮一枚「Parent」返回钮**，能回到父会话 | `ChatContainer.tsx:1008-1022`（`returnToParentButton`，`arrow-left` 图标）→ 挂载点 `:1621`；文案 `en.ts:2285-2288`：`'Return to parent session'`／`'Return to: {title}'`／按钮字面 `'Parent'`；判据 `parentSession = useParentSession(currentSessionId,…)` `ChatContainer.tsx:993` |
| **默认只读**：钻进去之后不能当场对它说话，要在设置里开 | 只读判定 `ChatContainer.tsx:1024-1028` `resolveChatPromptReadOnly(currentSession, embeddedAllowPrompting ?? allowPromptingSubagentSessions, readOnly)`；开关默认值 `stores/useUIStore.ts:1447` `allowPromptingSubagentSessions: false`；开关的原文标签 `"Allow Prompting Subagent Sessions"` `en.settings.ts:2136` |
| 子代理**在跑但还没被点进去**时，父会话这边也看得见它在干什么（活尾巴直接进那一行） | `ToolPart.tsx:1002-1010,1064-1069`（`TaskSummaryEntriesList` 逐条活动尾随渲染）；还没出活动时的两句占位 `:1051` |
| 子代理**在等批准**时，那条请求会**冒到父会话的审批坞**上并打标 | 取数面按"当前会话＋其所有后代子会话"作用域：`ChatContainer.tsx:838-841` 注释逐字 "only subscribe to permissions/forms for the current session **+ descendant subagent sessions**"；坞上的徽标 `PermissionDock.tsx:100-104` 文案 `"From subagent"`（`en.ts:2324`）；判据 `usePermissionResponse.ts:10-18`（`sourceSession?.parentID === currentSessionId`） |

### 4.x 三处**必须限定**的地方（不写清就会被读歪）

1. **工作状态面板里那一整块 `Subagents` 区，手机上根本挂不出来。**
   那块区的行是可点的、也写了移动分支（`UI/chat/work-status/WorkStatusSubagentsSection.tsx:64-76` `openChildSession`，`:102` `onClick`，行状态四字 `"needs permission"`／`"asked a question"`／`"is working"`／`"Done"` `:107-115`）——**但整块面板的前置条件是 `workStatusPanelMountable = !isMobile && !isVSCode && chatSurfaceMode !== 'mini-chat'`**（`ChatContainer.tsx:959`）。
   ⇒ 手机上"看有哪些孩子在跑、谁卡住了"这一屏**没有**，只剩聊天记录里那一行 `Open … subtask`。这一形对我们尤其值钱，因为它自己的注释把理由写死了：*"Running subagents and, more importantly, their **blockers**: a permission request raised by a child session **has no representation in the transcript**, so this panel is the only place it becomes visible"*（`WorkStatusSubagentsSection.tsx:21-25`）——**它承认"子会话的阻塞在正文里没有表示"，而那块补位面板在手机上缺席**。
2. **同一套东西在桌面上不是"翻页"而是"侧栏开一格"**（`ToolPart.tsx:1032-1038` 走 `openContextPanelTab(…, readOnly: true)`）⇒ 移动/桌面是**两套呈现面**，不是同一套的两尺寸（→ §7）。
3. **移动侧那条"点进去"的入口在正文里，不在名册里**：会话抽屉不按父子折叠子会话成可展开族谱（本腿在 `MobileSessionsSheet.tsx` 里没量到父子树；侧栏的父子树件在桌面端：`components/session/sidebar/sessions/SessionTreeItem.tsx` 存在，**是否被移动抽屉引用＝未读到**）。这一条留给编排者（§8）。
4. **新出现的反面素材（一处 i18n 漏口，就在子代理那一屏上）**：占位句 `'Waiting for subagent activity...'` 与 `'No subagent session id on task metadata.'` 是**硬编码英文字面量**、没走 `t()`（`ToolPart.tsx:1051`）。同文件其余文案都走 `t()` ⇒ 这一枚是漏网，不是设计。

---

## §5 必答二：移动端审批卡「本次／本次会话内／长期」＋拒绝／允许的理由框

### 5.1 三档：**只有两档在卡上；「本次会话内」不在卡上，它是另一枚独立的会话级开关**

| owner 问的那一档 | 移动屏上有没有 | 证据 |
|---|---|---|
| **「本次」（只这一次）** | **有**〔已证〕 | 动作类型逐字 `export type PermissionReply = "once" \| "always" \| "reject"`（`UI/lib/opencode/model.ts:366`）；钮 `"Allow once"` `en.ts:3370`，挂载 `UI/chat/PermissionCard.tsx:387-391`（坞变体）／`:407-417`（内联变体） |
| **「长期」（以后这类都别问）** | **有**〔已证〕，**且形状特殊**：钮上**直接把要存成哪条规则印出来**——文案逐字 `"Always: {patterns}"`（`en.ts:3372`），取服务端回的 `permission.save` 前 2 条＋省略号，全量挂到悬停提示（`PermissionCard.tsx:349-360` `useAlwaysLabel`）；无规则可存时退回 `"Always allow"`（`en.ts:3371`） | 同上 |
| **「本次会话内」** | **没有这一枚按钮**〔已证＝卡上无此档〕。**但它并没有消失，而是被搬去另一枚控件**：聊天框里的**盾牌钮**在**会话级**循环三档 `ask`／`safety`／`auto`，原文标签逐字 `"Permissions: ask every time"`／`"Permissions: safety net, ask only before risky actions"`／`"Permissions: accept everything"` | 类型 `stores/utils/permissionAutoAccept.ts`（`PermissionMode`，由 `PermissionAutoAcceptButton.tsx:18` 引用）；三枚图标 `PermissionAutoAcceptButton.tsx:30-34`（`shield-user`／`shield-star`／`shield-check`）；文案 `routing.i18n.ts:58-60`；一次按压的行为 `UI/chat/permissionAutoAccept.ts:15-42`（草稿改草稿档、会话落会话档、没会话先逼你开 `:29-37`）；挂载在**移动分支**里 `composer/ui/ComposerFooter.tsx:147`（`isMobile` 分支起点）→`:166-171` |
| 会话档的**默认值**另有一处可改 | Settings → Sessions 的三枚芯片 `Ask every time`／`Safety net`／`Accept everything` | `components/sections/openchamber/PermissionDefaultModeField.tsx:109-119`；标题 `"Permissions"` 在 **`routing.i18n.ts:65`**，三枚芯片字面在 `:67-69` |
| 卡上还有第四种**状态条**（不是动作） | 安全网扣住时那一行 `"The safety net held this action for you to decide · <风险类>"` | `PermissionCard.tsx:300-308`；`routing.i18n.ts:5`；七个风险类 `PermissionCard.tsx:55-63` |

**移动屏上这一族的实际渲染件是"坞"而不是"卡"**：`UI/chat/PermissionDock.tsx`（挂在 `ChatInput.tsx:4161-4165`，位置＝**输入框正上方的浮层**），一次答一枚、答完下一枚顶上（`:44-48`），批量请求用"一坞走一遍＋每枚一颗点"表达（`:111-132`）。内联的 `PermissionCard` 在移动上只服务于 **BTW 分身层里的子会话请求**——注释逐字 "The BTW sheet keeps the inline `PermissionCard` for its child session's requests"（`PermissionDock.tsx:22-24`）。

### 5.2 理由输入框：**没有**〔已证〕

| 查法 | 结果 |
|---|---|
| 在两张卡／坞里找输入面 | `grep -n "textarea\|<Input\|Textarea" PermissionCard.tsx PermissionDock.tsx` ＝ **NO MATCH**（本腿现跑） |
| 答复那条**传输路**带不带话 | 不带：`respondToPermission(sessionId, requestId, response)`（`UI/sync/session-actions.ts:2085-2099`）三个参数，落到 `replyToPermission(sessionId, requestId, response, {directory})`——**没有任何 reason/message/comment 形参** |
| 那 `permission.message` 是什么？ | 是**它自己的**说明、只读显示：`PermissionCard.tsx:311-316`，注释逐字 *"v2 lets **the agent** explain in its own words why it needs this"*。**方向是单向的（它→人），不是（人→它）** |
| 键盘快捷键能不能代打字表达意见 | 不能。三枚动作只有 `alt+enter`／`alt+shift+enter`／`alt+backspace`（`usePermissionResponse.ts:52-70`，键位映射逐字在 `:56-59`，且**只有最新那张卡吃键盘** `:8,54`）；而且这三枚钮上的键帽提示在窄屏**直接不渲染**（`className="… hidden sm:inline"`，`PermissionCard.tsx:379,385,390,416,429,441`）⇒ **手机上看不见快捷键提示** |
| 结论档位 | 理由框＝**〔已证＝不存在〕**（不是"建了没接"：传输路上连字段都没有；区别于同仓另一处 `permission-toast.ts` 那种"有字段无渲染者"的形状） |

### 5.3 我方对照（只引 `docs/**`／本腿现读 Go）

- 我们的三档是**另一套切法**：`docs/PLAN.md:2028` 逐字三档＝「仅本次」／「**本会话内**允许向此应用注入」／「本会话内允许向**所有**应用注入」；`docs/PLAN.md:2147` 确认卡第三选项＝「本会话内允许 `<工具>` 于 `<路径模式>`」，**会话结束即失效**。
- **我们目前根本没有"长期"这一档**：`docs/PLAN.md:2157` 逐字 `workspace`/`user` 持久档**仍为 RESERVED**。⇒ 移动这一屏给出的对照值：**`Always` 那一枚把"要存成什么规则"印在按钮上**（`PermissionCard.tsx:349-360`），这正是"长期档不敢让人看清自己在存什么"的反面教材。
- **我们的确认卡现在连"允许"都没有**：`docs/specs/SPEC-08-ui-ball-panel.md:195` 逐字"按钮只有「拒绝」与「查看完整参数」——**没有「允许」**（F2 定案）"；`docs/specs/SPEC-06-security-gatekeeping.md:19` 逐字「"允许"只接受**原生侧来源**」⇒ 面板侧（＝移动那一屏的同族位置）结构性给不出"允许一次"。**别家的会话级"坞＋三枚钮＋进度点"这一整套形状，在我们这里被门规砍掉了一半**。
- 我们的答复出口现状（本腿现读）：`internal/agent/approval/approval.go:103-114` 四条"这条通道今天不可用"的实字符串；审批包内**没有** `AllowOnce`／`SessionScope` 这类动作枚举（尺：`grep -rn "AllowOnce\|SessionScope\|session_scope" internal/agent/approval/*.go internal/panel/*.go`＝除测试里的字符串外 **0 命中**，本腿现跑）⇒ **D45 的三档在代码里还没有名册**。

---

## §6 我写错的条目（拿第四版与移动端相关条目逐条对抗）

> 尺：第四版＝`docs/reports/missing-features-2026-09-29-v4.md`；本腿逐条 `sed -n` 回原文核。**只列真错**，不列"还可以更细"。

| # | 第四版原话（行号） | 本轮证据 | 更正后的说法 | 严重度 |
|---|---|---|---|---|
| **6-1** | 「与会话侧栏共用同一判据 `mobile/mobileWidgetSnapshot.ts:18-20,84-98`」（`v4.md:57`，另见 `:298` 同一路径） | 该文件**不在 `packages/mobile/`**：`find packages/mobile -name '*.ts'` 只回 **`packages/mobile/capacitor.config.ts` 一枚**；`mobileWidgetSnapshot.ts` 的真身在 **`packages/ui/src/apps/`**（判据逐字在文件头注释 `:18-20`，本腿已核） | 路径改写为 **`packages/ui/src/apps/mobileWidgetSnapshot.ts:18-20,84-98`**。**这一改不只是修错字**：它说明"移动逻辑"绝大多数**不在移动壳包里，而在共享 UI 包里**（壳包只有原生胶水＋`capacitor.config.ts`）——按原写法去找 `packages/mobile/` 会一无所获 | **路径级错，会导致下位找不到文件** |
| **6-2** | 「点通知直达那一条会话……`mobile/.../deepLinks.ts:23-30,46-53`」（`v4.md:134`） | 同上：`deepLinks.ts`／`deepLinkNavigation.ts` 均在 **`packages/ui/src/apps/`**（`find` 现量），`packages/mobile/` 里没有 | 路径改写为 **`packages/ui/src/apps/deepLinks.ts:23-30,46-53`** ＋ `apps/deepLinkNavigation.ts:15-16,133-135`；意图名册本腿复核＝**7 枚**逐名 `session`／`new-session`／`sessions`／`status`／`settings`／`changes`／`view`（`deepLinks.ts:23-29`），**"7 种"这个数成立**，只有前缀错 | 同上（结论不变） |
| **6-3** | 「openchamber：`Permission required`（`en.ts:619`）＋卡上三枚……」（`v4.md:240`） | 同一个 i18n 目录里**有两枚同字面的键**：`en.ts:619` 逐字 `'sessions.sidebar.session.status.permissionRequired': 'Permission required'` ＝**会话名册行上的状态 chip**；**审批卡的标题**另有其键，逐字 `'chat.permissionCard.title': 'Permission required'` 在 **`en.ts:3354`**（读取点 `PermissionCard.tsx:473`／`PermissionDock.tsx:82`） | 引用**审批卡**时应写 **`en.ts:3354`**；`en.ts:619` 属于名册状态那一族（同族还有 `:615-624` 十枚：`Session active`／`Unread updates`／`Pinned session`／`1 pending question`／`Last turn took {duration}`…）。**字面没错、行号张冠李戴**——两处都叫 "Permission required"，最容易在这种地方混 | 引用精度 |
| **6-4** | 「openchamber **画三枚**（`Allow once`／`Always…`／`Deny`，`PermissionCard.tsx:407-442`）」（`v4.md:85`；`:304` 同） | `:405-450` 是**内联变体**（`variant:'inline'`），注释逐字写着它是给 **BTW 分身层**用的那一张（`PermissionCard.tsx:453`；`PermissionDock.tsx:22-24` 再确认一次）。**聊天屏（含手机）真正渲染的那三枚在同一个文件的另一支：`variant==='dock'` 的 `PermissionCard.tsx:373-393`**，挂载点 `ChatInput.tsx:4161` | 三枚动作这个结论**成立**，但引用应指向 **`PermissionCard.tsx:373-393`（坞变体，聊天屏实际用的那支）**，并注明"内联变体 `:405-450` 只服务 BTW 子会话卡"。**顺带一条只在坞里才有的东西：多枚请求时的进度点阵（`PermissionDock.tsx:111-132`）＋ `{current}/{total}` 步序（`:105-109`）**——引内联那一支就看不见这半套 | 引用到了不渲染的那一支 |
| **6-5** | 第四版把移动壳描述为"壳层只是壳、真 UI 复用 web 构建"这一族说法的延伸（`v4.md:515-519` 只纠了分母；壳内文件构成沿用 09-28 调研件），且第四版自己承认「设置屏 21 项各自的正文未读」「批准卡正文未读」（`v4.md:583`） | 本腿量到的是**三处更硬的事实**：① 移动壳白名单是**写死的 22 页常量**（`MobileApp.tsx:88-111`），不是"整份设置塞小屏"；② **`about` 这一页只有手机上才有**（`lib/settings/metadata.ts:225` 的 `isAvailable` 逐字 `ctx.isMobile && !ctx.isVSCode`），桌面/网页端没有这一页；③ `extensions` 页逐字写着 `!ctx.isVSCode && !ctx.isMobile`（`metadata.ts:233`） | "移动＝桌面那一屏缩小"这一族说法要加限定：**移动是一份自己写的页面白名单＋一枚只在移动上存在（`about`）＋若干枚被点名摘掉（`shortcuts`／`extensions`／`Interface Font Size`／`Enable Spellcheck`／`Send shortcut`／`Session tabs`／`Terminal Quick Keys`）**。设置项总数本腿现量＝注册表 **181 枚**（`lib/settings/search.ts`，尺 `grep -c "^    id: '"`），落在移动 22 页内的 **158 枚** | 结构性低估移动这一屏的独立性 |
| **6-6** | 「**子代理还在跑时父任务的完成通知不许发**」这一族（`v4.md:137`）暗示"子代理状态在别处看得见"；第四版另在 `:72` 只给了 DeepSeek 的族谱树行 | 本腿量到 openchamber **自己承认相反的一半**：工作状态面板注释逐字 *"a permission request raised by a child session **has no representation in the transcript**, so this panel is **the only place** it becomes visible"*（`WorkStatusSubagentsSection.tsx:21-25`）——**而这块面板在手机上整块挂不出来**（`ChatContainer.tsx:959` `workStatusPanelMountable = !isMobile && !isVSCode && …`） | 更正为：**"子会话的阻塞在正文里没有表示"是它自己写下的事实**；移动侧靠"作用域取数含后代子会话"（`ChatContainer.tsx:838-841`）＋坞上的 `From subagent` 徽标（`PermissionDock.tsx:100-104`）补了这一格，**但"有几枚孩子在跑／谁卡住"的汇总视图在手机上缺席**。这对我们票 197／220 那一族的含义变了：**别家也没白送，移动上照样缺一屏** | 结论方向要改 |

> **不在此列的**：第四版 §5-4 对"61 枚"分母的更正（`v4.md:515-519`）本腿复核**成立**（`apps/` 目录 52 枚文件／带 mobile·widget 字样 36 枚，本腿现量），不改。

---

## §7 不建议抄（只给形状级理由）

| 它的形状 | 出处 | ⛔ 为什么不该照抄（形状级，不是"实现成本高"） |
|---|---|---|
| **一份代码里 692 处 `isMobile` 分叉**（108 个文件），外加 36 处平板尺寸判定 | 尺：`grep -rn isMobile packages/ui/src --include='*.ts*' \| grep -v test \| wc -l` ＝ **692**；`useTabletLayout/isTabletLayout` ＝ **36** | 这是**双实现面**：同一功能两端各有一条渲染路径，于是"某端有、另一端静默没有"变成常态而不是事故——本腿现抓两例：桌面版草稿快捷指令条 `DraftPresetChips` 在手机上整块不挂（`ChatInput.tsx:3504-3506` 的 `&& !isMobile`）、子代理汇总面板在手机上整块不挂（`ChatContainer.tsx:959`）。**我们要的是"一屏一套真相"（球＋面板），不该再造一条按端分叉的渲染面** |
| **设置页白名单靠硬编码常量维护**（`MOBILE_SETTINGS_PAGES` 22 枚字面量） | `MobileApp.tsx:88-111` ＋ `:638` | 与注册表（`lib/settings/metadata.ts` 的 `isAvailable(ctx)`）是**两套并行机制**，同一件事（"这一页在这一端有没有"）有两个真相源：`shortcuts` 被 `SettingsView.tsx:254` **再摘一次**、`about` 反过来只在移动存在——**两处都得记得改**。我们有 D36 的"全量 section 树＋三档生效级别"，加端时**应扩展 `isAvailable` 那一侧，不要再叠一份写死的名单** |
| **冷启动/唤醒/切网络三处各自跑一遍"探活＋阶梯重试"** | `MobileApp.tsx:920-956`（冷启动）／`:964-1015`（带持久端点冷启）／`:731-826`（唤醒，含 `retryDelaysMs=[4000,10000]`）／`:840-852`（`online` 事件） | 四段代码各自判"要不要掉回连接屏"、各自持一份 seq 竞态守卫（`nativeResumeValidationSeqRef`）——**同一决策的四个副本**。我们单机常驻、没有"够不着服务器"这一族，**不要为了对称去造一套重试阶梯**；值得抄的只有它写下来的判语（"刚醒的网络会误报不可达，所以快判之后必须再用满预算重判一次"） |
| **把"上次开的是哪条会话"另存一份本地缓存**（带单调戳防同毫秒撞车） | `UI/sync/last-session-cache.ts:59-90`；消费点 `MobileApp.tsx:1037-1098`（含 6s 安全阀 `:1066`） | 这是**第二套真相源**：服务器有权威会话全集、本地还存一份"你上次在哪"。它自己就为这份双源写了三段防护（快照确认存在才恢复 `:1077`、不在册就清掉陈旧指针 `:1081`、超时兜底防卡在闪屏 `:1066`）。**我们的 D35/D43 已经把"权威态"定在 SQLite＋状态机转移表上，再引一份端侧会话指针缓存＝开一类"两边不一致"的新缺陷面** |
| **令牌分端存储：原生进 Keychain/Keystore、网页端进 localStorage** | 09-28 调研件已记（`mobileConnections.ts` 的"native 上令牌绝不进 localStorage"）；本腿只复核到文件在 `apps/mobileConnections.ts`（1,744 行，**正文未逐行读**） | 存储面按端分叉＝**每个端都要各自证明一次"我这一面不泄"**。我们的对应面是 DPAPI（`internal/secret/`）＋ D36 的配置面，**应坚持"一个秘密只有一个存放处"，不要按端开第二处** |
| （对照面，非"不建议"）**移动壳里那些 OS 级件**：锁屏小组件、NSE 离线改快照、GCKeyboard 实体键盘判定、扫码配对 | 09-28 调研件 §Q1 已列，本腿不重复取证 | 这些**不是"不建议抄"，是"落不了原样"**（Win32 无 App Group/WidgetKit 对位物）。登记在此只为防止下一位把它们误当"可迁移形状"——**可迁移的是设计账（快照函数＋事件驱动刷新＋离线推送件），不是 API 面** |

---

## §8 没做完／留给编排者

### 8.1 本腿**故意不碰**的（那是块 B／块 C 的地界，越界即缺陷）

| 留给 | 内容 | 为什么它不属块 A |
|---|---|---|
| **块 B** | 同一仓的 **VS Code 扩展那一屏**：`packages/vscode`（114 枚 `.ts`、bridge 32 种消息、git 24 操作、20 家订阅配额面板、行评论线程）；以及 `VSCodeApp.tsx`／`renderVSCodeApp.tsx` 这两枚移动同名兄弟件 | 它们是"IDE 那一屏"，本腿只在 `MobileApp.tsx` 的 `isVSCode` **分叉点**上引用（为了说明移动这一屏哪里与 IDE 共用判定），未读任何 vscode 包正文 |
| **块 B** | `packages/extensions`（3 枚文件）与 `packages/sdk` 正文 | 扩展生态面 |
| **块 C** | **DSH 那几十枚 `ui-*` 分包** | 另一家的另一种分包形状 |
| **块 C** | 本仓 `packages/ui` 里**非移动命名**的那 46 枚 `apps/` 件（`App.tsx`／`ElectronMiniChatApp.tsx`／`VSCodeApp.tsx` 等）与 `components/session/sidebar/**` 桌面侧栏族 | 桌面/迷你聊天那一屏 |

### 8.2 块 A 内**没读完**的（具名到文件，别把本腿当"移动屏已读穿"）

1. **`UI/apps/MobileSessionsSheet.tsx`＝2,642 行，本腿只按控件点名、未通读**。已量的入口：`SessionSearchInput` 引用 `:41`、视图切换 `:150-151,2386`、重排态 `:1718,1880`。**未读**＝行内所有子件正文、分页/无限滚动、工作树区、归档区实现。
2. `UI/apps/MobileSessionMetadata.tsx`（468 行）＝**只读前 130 行**（环形进度件与三态 `ContextDisplay`）；面板其余部分、`usageDisplayMode` 那一段未读。
3. `UI/apps/MobileChangesSurface.tsx`／`MobileFilesSurface.tsx`／`MobileTimelineList.tsx`／`MobileInstancesSurface.tsx`／`MobileProjectEditSurface.tsx`／`MobileFullscreenSurface.tsx`（270 行，只读了挂载点）＝**逐枚未通读**；§1.F/§1.G 里给的是它们在壳里的接线，不是内部控件清单。
4. `mobileConnections.ts`（1,744 行）＝只读关键行（`connectionDisplayUrl` 用法、探活结果枚举 `:812,1007` 附近）；令牌存储那一段正文未逐行读（§7 那一行据此降级为"沿用 09-28 结论＋文件在位已核"）。
5. `mobileNativeChrome.ts`（462 行）＝**本腿一枚未读**（键盘/状态栏/返回键的原生接线细节都在里面；§1.A2、§1.G 的结论来自 `MobileApp.tsx` 一侧）。
6. **审批与子代理两问已答死，但相邻三处没扩**：`permissionToolPresentation.tsx`（工具→图标/名字的映射表）、`permissionSummary.ts`、`permissionCardPatterns.ts` 正文未读。
7. **设置项 158 枚里本腿只逐枚核了 12 枚**（§2.1）。其余按页登记，标签原文未取。**下一位若要给 owner 一份完整设置对照表，得按 `lib/settings/search.ts` 的 181 枚逐条去 `en.settings.ts` 取字**——本腿没做，这是量级取舍不是遗漏。
8. **服务端一侧本腿一枚未读**（`packages/web/server/**`）：凡涉及"服务端到底回什么字段"的结论（尤其 `permission.save` 里那些模式由谁算、坞上"长期"存到哪里、推送触发的四分类）**一律是客户端反推**〔客户端反推＝不是已证〕。
9. `mobile.changes.placeholder.description`（`en.ts:292`）＝"这里将来会放工作树评审/同步/提交动作"这一占位句，本腿量到字面，**未确认它挂在哪一枚件上、是否出厂就会露出**——别拿它当"已实现的改动页"。

### 8.3 要编排者拍的三件事（本腿不自行填）

1. **§6-4 那条会影响票面措辞**：第四版引的是**不渲染的那一支**（内联卡）。若已按 `:407-442` 写过票面（197／213／219／220 那一族），要不要改引 `:373-393` 并补"坞上的进度点阵"这一格——**涉及票面文字，我不动**。
2. **§4 与 §6-6 合起来是一条我方判断**："子会话的阻塞在正文里没有表示"＝别家自己写下的事实，而移动上补位那屏缺席。⇒ **我们票 197／220 若要给"子代理卡住"做人可见的表示，别家没有可抄的现成答案**（这条与第四版 `v4.md:319` 那句"这一格早晚得我们自己拍"同向，本腿只是把移动侧也钉上了）。**是否据此改判据，请编排者定。**
3. **§5.3 末行**：D45 的三档在我们代码里**还没有动作名册**（`grep AllowOnce|SessionScope|session_scope internal/agent/approval/*.go internal/panel/*.go`＝除测试字符串外 0 命中，本腿现跑），而 `docs/PLAN.md:2157` 逐字写着 `workspace`/`user` 持久档**仍为 RESERVED**。别人移动屏上的第四枚（会话级 `ask|safety|auto`，与卡上三枚**分家**）是**一种我们没有的切法**：它把"这一会话要不要问"做成一枚**常驻可循环的钮**、而不是确认卡上的一档。⇒ **碰 `Q-49`（面板侧不许给 L2"允许"），我不自行填。**

---
