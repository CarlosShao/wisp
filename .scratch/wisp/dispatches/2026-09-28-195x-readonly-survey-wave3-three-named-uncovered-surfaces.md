# 派单 2026-09-28 19:5x — 只读调研第三轮（三腿）：**上一轮自己点名的三个未覆盖面**，不是再笼统"深挖"

> owner 的原话（09-28，逐字）：「**我需要你反复调研，反复学习，继续深挖到底有没有疏漏之处**」
> 「**取各家之所长，这样才能打磨出属于我们自己独有的 harness，这才是重点**」。
> 上一轮交件自己写了边界（账 `A409`）：**没有逐文件通读任何一家**，且点名三处没碰——
> ① `openchamber` 的 `mobile`／`vscode` 两包；② DeepSeek harness 的 **52 枚 `packages/client/ui-*` 分包**；③ `minimax-code` 的 `packages/agent-modules` 整棵树。
> ⇒ **本轮就按这三处派**，别再交一份"目录级抽样"当覆盖率。

## 0. 三家的位置与我 19:5x 现量的枚数（进去先自己复核一遍再动手）

| 家 | 路径 | 现量 |
|---|---|---|
| OpenChamber | `D:\work\AI\open source\openchamber` | 顶层 `docs electron extensions mobile sdk ui vscode web`；**没有 `.git`**（tarball 解的，`oc.tar.gz` 在同级）⇒ **不许引提交历史，每条结论必须带 `文件:行`** |
| DeepSeek harness | `D:\work\AI\open source\deepseek-harness` | `packages/` 57 枚；**`packages/client/ui-*` 52 枚**（我现列：`ui-agent-preset ui-approval ui-attachment ui-brand-official ui-chat ui-commands ui-conversation ui-deliverables ui-directory-picker-browse ui-directory-picker-native ui-dockkit ui-goal ui-input-trigger ui-jobs ui-layout ui-message-feedback ui-model-selection ui-open-in-app ui-permission-presets ui-plan ui-plugin-manager ui-primitives ui-reference ui-renderer ui-schedule ui-session ui-settings ui-settings-account ui-settings-agent-loop ui-settings-general ui-settings-models ui-settings-plugin-inventory ui-settings-plugins ui-settings-shell ui-settings-subagent ui-settings-web-search ui-shortcuts ui-sidebar ui-sidebar-browser ui-sidebar-documentpreview ui-sidebar-files ui-sidebar-right ui-sidebar-terminal ui-skill ui-slots ui-subagent ui-theme ui-tool ui-trajectory ui-user-questions ui-workflow-run ui-workspace`） |
| MiniMax code | `D:\work\AI\open source\minimax-code` | `packages/` 15 枚；**`packages/agent-modules` 12 个模块、147 枚 `.ts`**：`background-task context-manager conversation-contract cron goal mcp permission plugin-hooks runaway-guard session-report skills system-reminder` |

⚠ 那份名册是我 19:5x 从 `ls` 直接抄的，**抄名可能有手误**（比如 `ui-slots` 这类），进去先自己 `ls` 复核再动手，别把我的抄写当名册。

## 1. 三腿的分工（**各自只写自己那一枚产出件**，别互相覆盖）

- **腿 1（openchamber 未覆盖面）**→ 写 `docs/reports/survey-2026-09-28-oc-mobile-vscode-extensions.md`
  目标：**同一个 harness 换到"手机"和"IDE 里"这两块屏幕上，多了哪些我们在桌面版根本没有的东西**。
  必看：`mobile/`、`vscode/`、`extensions/` 三包逐目录读完（先数枚数再答应做得完）；
  重点问四件：**① 移动端怎么回答"要人批准"那一跳**（推送？回消息？只做只读？）、**② 移动端怎么看待语音/键盘**、
  **③ IDE 版把哪些"我们放球上的信息"搬过去了**（当前目录、当前模型、档位、成本）、**④ 它有没有"离开桌面后继续跑"的续法**。
- **腿 2（DSH 的 52 枚 `ui-*` 分包）**→ 写 `docs/reports/survey-2026-09-28-dsh-ui-packages.md`
  目标：**把"界面"拆成可数的 52 块**，逐块回答"这块对应我们哪个功能、我们有没有"。
  做法：**先全列枚数与每包 README/`index` 的第一段**（别跳读源码就下结论），
  再**按我们的缺口排序**深入这 8 枚（其余点名说"只登记未读"）：
  `ui-approval`、`ui-subagent`、`ui-settings-subagent`、`ui-trajectory`、`ui-user-questions`、`ui-dockkit`、`ui-layout`、`ui-shortcuts`。
  ⚠ 特别回我一个问题：**`ui-trajectory` 与 `ui-subagent` 是不是"点开一枚子代理看它自己那页流式输出"的那个东西**（owner 09-28 18:0x 硬要求：「**子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面**」）。
- **腿 3（MiniMax `agent-modules` 整棵树）**→ 写 `docs/reports/survey-2026-09-28-minimax-agent-modules.md`
  目标：**这是三家唯一把"agent 的护栏"做成独立模块树的一家**，12 枚逐枚读。
  必查：`runaway-guard`（**防跑飞**：怎么判定"它在原地打转"、判定后停谁）、`permission`（**授权梯度与我们 R20 三档怎么对**）、
  `cron`/`background-task`（**定时与后台**：D12 我们只登记了缺口，人家成模块）、`session-report`（**"今天它替我做了哪些决定"那枚日报**，我 `A409` 第九节猜五家都没有——**这条如果被这枚模块推翻，明说推翻我了**）、
  `system-reminder`/`context-manager`/`conversation-contract`（**上下文与提醒的分层**，对我们 D15/D39 有直接参照）、`plugin-hooks`、`skills`、`goal`、`mcp`。

## 2. 共同的硬规矩（只读程也一样）

- **绝不修改 `D:\work\AI\open source\**`**（那是只读参照，不是我们的仓）；**不要在那边建文件**。
- 在本仓**只写你自己那一枚产出件**；**不碰任何源码**、不碰 `docs/reports/pending-and-issues.md`、不碰别人的证据件；**不 commit**（交件由我落账）。
- 不许跑 `frontend/**` 相关的读取与转述（界面那支的地界）；本仓 `PLAN.md`/`docs/specs/**` 只当"我们有什么"的对照读物。
- **每条"他们有"都要带 `文件:行`**；**每条"我们没有"要带我们在仓里的现量**（`grep`/`ls` 读数），量不到就写**〔未量〕**，不许用"应该没有"。
- **枚数带推导式与口径**（"52 枚＝`ls packages/client | grep -c '^ui-'`"）；只报"很多/若干"＝没交件。
- **引别家机制时要写"它在哪一层做"**（进程内/子进程/服务端/OS API），因为我们要判断能不能落到 Win32＋Go 上。
- 结尾必须有两节：**① 我这一腿没读到什么**（具名到目录）、**② 我认为编排者会误读的地方**（我之前的哪一条账可能被这轮推翻）。

## 3. 交件形式

- 一份 md，**大白话标题＋表格**（我会把它并进给 owner 的第三版差集清单，他看的是"人话后果"不是术语）。
- 每条缺口写成一行三段式：**他们怎么做的（带行）／我们现在是什么样（带现量）／落到我们身上是什么后果**。
- 我会**自己抽查 3–5 条**逐名去读原文；**你写歪一条我就整批当未验证处理**。
