# Wisp 外部对标 · 第五版块 B —— VS Code 扩展那一屏（Cline ＋ OpenChamber 扩展形态）

> 本腿只做一件事：**装了扩展之后 VS Code 里出现了哪些界面件**，逐块逐控件。
> ⛔ 不越界：移动端＝块 A（`block-A.md` 已交）、组件分包＝块 C；本文件一处不提它们的屏。
> ⛔ 零读零引本仓 `frontend/**` 与 `design/**`；"我方有没有对应物"这一列只引 `docs/**` 行号或本腿现读的 Go 文件行号。
> ⛔ 本腿不跑任何编译／测试／门禁。
> 三档标法：**〔已证〕**＝读到实现／**〔建了但没接〕**＝有代码无调用者／**〔仅文档〕**＝只有 README 或文档里有。

---

## §0 本轮取数的家与出处

| 家 | 仓库 | 取数方式 | commit sha | 日期 | 扩展本体在哪 | 版本 |
|---|---|---|---|---|---|---|
| **Cline** | `github.com/cline/cline` | `curl` 拉 codeload tarball（`main`）解到 `/tmp/v5b/cline-fresh`，40,530,055 字节 | `c604735fe1ce16d745f9dbf799ade79a4e3dad86`（`api.github.com` 具名） | 提交时间 `2026-09-30T02:57:24Z`（"chore(desktop): release v0.0.38"） | `apps/vscode/`（扩展 id 内部名 `claude-dev`、`displayName` ＝ `Cline`、publisher `saoudrizwan`）`apps/vscode/package.json:2-14` | **4.1.21**（`package.json:5`） |
| **OpenChamber** | `github.com/openchamber/openchamber`（10,942 star，`api.github.com` 查得） | **本腿读的是新鲜 tarball**：`curl` 拉 codeload（`main`）解到 `/tmp/v5b/oc-fresh`，74,398,609 字节、`gzip -t` 通过；同机那份 `D:\work\AI\open source\openchamber`（2026-09-28 落的 **2.0.3** 快照）**只用来对前几轮的行号**，不作为引证面 | `d78dac542d796f20536a3bb94971aa94fbb9e48e`（`api.github.com` 具名） | 提交时间 `2026-09-30T00:43:13Z`（"fix(browser): don't retry a restored tab's dead dev server at launch"） | `packages/vscode/`（`displayName` `OpenChamber`、publisher `fedaykindev`）`packages/vscode/package.json:2-6` | **2.0.4**（`packages/vscode/package.json:5`）；⚠ 比 09-28 那份快照**多一个修订号** ⇒ 前几轮报告里指向 `packages/vscode/...` 的行号**可能与本腿读数差几行**，两处已注明（`missing-features-2026-09-29-v4.md:128,138`、`survey-2026-09-28-oc-mobile-vscode-extensions.md:38`）。本腿所有 OpenChamber 引证一律指这份新鲜快照。 |

**取数代价（必须写明白）**：tarball 快照**没有 `.git`、没有 `node_modules`** ⇒
①**读不到提交历史**，本文件**每一条结论都带 `文件:行`**，一律不引 commit／PR／历史；
②**读不到依赖源码**，凡是"依赖 VS Code 宿主能力"的地方只能读调用点、读不到实现；
③Cline 那份是**monorepo**（`apps/vscode` ＋ `apps/cli` ＋ `sdk/` ＋ `docs/`），本腿只读 `apps/vscode` 这一支与其 webview 子目录 `apps/vscode/webview-ui`，**`sdk/` 与 `apps/cli` 只在需要确认"扩展形态有没有对应界面"时点名，不做 CLI 调研**（那是块 C／别的腿）。

**Cline 与 OpenChamber 的"扩展"不是同一种东西**（先定口径，否则表一会读歪）：
- Cline：扩展宿主进程＋**一枚常驻 webview 视图**，所有对话、设置、审批都画在那枚 webview 里；VS Code 原生设置页 **一项都没有**（`contributes.configuration.properties` 空，实测读数见 §2）。
- OpenChamber：扩展宿主进程里跑一台 web 服务器＋**webview 装的是它自家 web 应用**，所以"扩展形态的界面"＝web 界面＋一串只有扩展才有的宿主件（评论线程、标题栏按钮、状态栏、命令）。本腿只登记**扩展独有的那一层**，web 那一层是块 A／第四版的地界（重复即缺陷，见 §0.1）。

### §0.1 与前几轮的分工（防重）

- 第四版 `docs/reports/missing-features-2026-09-29-v4.md` 里 OpenChamber 扩展只被点了**三枚具名件**：`packages/vscode/src/quotaProviders.ts:15-52`（用量窗口）、`bridge-permission-auto-accept-runtime.ts:1-103`（每会话自动接受）、`bridge-git-*.ts` 一族（转引 `docs/reports/survey-2026-09-28-oc-mobile-vscode-extensions.md:38`）。
- `survey-2026-09-28-oc-mobile-vscode-extensions.md` §1-Q1 给过一份"IDE 多出来的"**散文式清单**（评论线程／右键菜单／会话进标签页／配额盘／bridge 32 种消息）。
- **本块的增量＝把它换成三张表**：`package.json` 的**每一枚贡献点逐枚列**（视图／命令 19 枚／键位／七组菜单／getStarted 走查）、**设置页逐条标签原文**、**状态文案逐句**，并**第一次读 Cline**（前四版六大素材家里没有 Cline，`missing-features-2026-09-29-v4.md:13-16` 逐名列的是 Step-Code／deepseek-harness／minimax-code／openchamber／pi／pi-upstream）。
- 已点过的三枚具名件本块**不重述内容**，只在需要"它长在哪一屏"时行号指回。

---

## §1 表一 · VS Code 扩展逐块逐控件表

> 路径口径：Cline 的引用一律相对 **cline 仓根**（扩展本体在 `apps/vscode/`），OpenChamber 一律相对 **openchamber 仓根**（扩展本体在 `packages/vscode/`）。
> 列顺序＝位置｜控件名（界面原文标签，逐字抄）｜它干什么｜我方有没有对应物｜出处 文件:行。
> 我方那一列只引 `docs/**` 或本腿现读的 Go 文件行号；**〔界面地界·不读〕**＝那一层在前端仓，本腿无权判定，只能给 Go 侧的读数。

### 1.0 先看"装完扩展，VS Code 里多了哪几处位置"（总览，逐处计数）

| 位置 | Cline 4.1.21 | OpenChamber 2.0.4 | 出处 |
|---|---|---|---|
| 活动栏容器（左侧那一竖条图标） | **1 枚**，标题原文 `Cline` | **1 枚**，标题原文 `OpenChamber` | `apps/vscode/package.json:110-118` ／ `packages/vscode/package.json:43-51` |
| 侧栏视图 | **1 枚**，`type:"webview"`，id `claude-dev.SidebarProvider`，`name` 是**空串**（标题栏不写字，只靠图标） | **1 枚**，`type:"webview"`，id `openchamber.chatView`，名字走本地化键 `%views.chat.name%`＝原文 `Chat` | `apps/vscode/package.json:119-128` ／ `packages/vscode/package.json:52-60` ＋ `packages/vscode/package.nls.json`（`views.chat.name`→`Chat`） |
| 视图标题栏（那一横排的图标按钮） | **5 枚** | **4 枚** | `apps/vscode/package.json:255-281` ／ `packages/vscode/package.json:212-233` |
| 编辑器标题栏（tab 右上角） | **0 枚** | **1 枚** `Open New Session in Editor` | `packages/vscode/package.json:206-211` |
| 编辑器右键 | **1 枚** `Add to Cline`（仅在有选区时） | **1 个二级菜单** `OpenChamber`（里面 4 项） | `apps/vscode/package.json:283-288` ／ `packages/vscode/package.json:193-198`＋`:162-167` |
| 资源管理器右键 | 0 | **1 枚** `Attach to OpenChamber Chat` | `packages/vscode/package.json:199-205` |
| 终端右键 | **1 枚** `Add to Cline` | 0 | `apps/vscode/package.json:290-294` |
| 源代码管理（Git 面板）标题栏 | **2 枚**（生成中互相替换） | 0 | `apps/vscode/package.json:296-306` |
| 笔记本工具栏／单元格标题 | **3 枚** | 0 | `apps/vscode/package.json:308-325` |
| 编辑器"灯泡"（代码操作） | **4 项**（其中 1 项只在有问题诊断时出现） | 0 | `apps/vscode/src/extension.ts:249-303` |
| 评论线程（原生批注） | 有实现，**零调用者**〔建了但没接〕 | **接了**（行评论＝待发上下文） | 见 §1.H |
| **状态栏（窗口底部那一行）** | **0 枚**〔已证：无〕 | **0 枚**〔已证：无〕 | 见 §1.M |
| 输出面板频道 | **1 条** `Cline` | **2 条** `OpenChamber` / `OpenChamberManager` | `apps/vscode/src/hosts/vscode/hostbridge/env/debugLog.ts:4` ／ `packages/vscode/src/extension.ts:56`、`packages/vscode/src/opencode.ts:29` |
| 终端 tab | **1 族**，名字固定 `Cline` | 走 bridge 的 exec，不占名 | `apps/vscode/src/hosts/vscode/terminal/VscodeTerminalRegistry.ts:28` |
| 命令面板 | **19 条**（含 2 条 dev-only） | **17 条**（其中 2 条**被显式藏掉**） | 见 §1.E |
| 键位（VS Code 里可搜可改的默认键） | **3 条** | **0 条**（`contributes.keybindings` 是 `null`） | `apps/vscode/package.json:233-253` |
| 首装走查（getStarted） | **5 步** | **无** | `apps/vscode/package.json:58-109` |
| OS 通知（弹到系统那一种） | **无**〔已证：名册零命中〕 | **有**，webview 侧发、带可编辑模板 | 见 §1.I |
| 自定义设置项（VS Code 设置页里） | **0 项**〔已证：`properties` 是空对象〕 | **2 项** | 见 §2 |

⚠ 两家的**最低宿主版本差 16 个大版本**：Cline 要求 `engines.vscode: ^1.101.0`（`apps/vscode/package.json:7-9`），OpenChamber 只要求 `^1.85.0`（`packages/vscode/package.json:14-16`）。这一格直接解释了上一表的多出项：Cline 敢用 `vscode.changes`（多文件"更改"审阅页，见 §1.G）是因为它可以要求新宿主；OpenChamber 要兼容老宿主，所以只能自己用 `vscode.diff`＋虚拟文档搭。**我们的宿主是自己写的窗口，这一档由我们说了算**（对照：`docs/specs/SPEC-08*.md:139` 把面板定为 WebView 低频件）。

⚠ 激活时机也是形状差：Cline `activationEvents: onLanguage / onUri / onStartupFinished`（**开机就起**，`apps/vscode/package.json:42-46`）；OpenChamber `onCommand:openchamber.openSidebar / onCommand:openchamber.attachExplorerToChat / onView:openchamber.chatView`（**点了才起**，`packages/vscode/package.json:37-41`）。

### 1.A 活动栏／侧栏视图（两家都只有一枚 webview，区别在"名字给不给"）

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| 活动栏 | `Cline`（容器标题）＋图标 `assets/icons/icon.svg` | 点一下开侧栏那屏 | 〔界面地界·不读〕Go 侧有常驻球与面板宿主规划（`docs/specs/SPEC-08*.md:139` 把面板列为 WebView 低频件） | `apps/vscode/package.json:110-118` |
| 侧栏视图 | `name: ""`（**空标题**） | 视图头不显示文字，把整条宽度让给 webview 自己的导航 | 没有这种"故意不写名"的先例；我们的面板名册是 C17 四条方法（`internal/panel/bridge.go:42-45`） | `apps/vscode/package.json:119-127` |
| 侧栏视图 | `Chat`（本地化键 `views.chat.name`） | OpenChamber 给视图写了名字，并且这名字**跟着 VS Code 界面语言走**（另有 fr／tr 两份 nls） | 没有〔已证：仓里没有 `package.nls.*` 这一类宿主文案〕 | `packages/vscode/package.json:52-60`、`packages/vscode/package.nls.json`、`packages/vscode/package.nls.fr.json`、`packages/vscode/package.nls.tr.json` |
| webview 宿主配置 | `retainContextWhenHidden: true` | 侧栏被切走时**不销毁网页**，回来还在原来那一屏 | 我们的面板是"关掉就没状态"：`docs/specs/SPEC-08*.md:151` 逐字要求"前端不得缓存任何跨 show 的业务状态（历史读 SQLite、配置读 TOML）" | `apps/vscode/src/extension.ts:123-125`、`packages/vscode/src/SessionEditorPanelProvider.ts:162-170` |
| 首装走查 | `Meet Cline, your new coding partner` ＋5 步标题（`Start with a Goal, Not Just a Prompt`／`Let Cline Learn Your Codebase`／`Always Use the Best AI Models`／`Extend with Powerful Tools (MCP)`／`You're Always in Control`） | VS Code 的 getStarted 页：装完第一屏，带进度点 | **没有**〔已证：`cmd/wisp` 与 `internal/ball` 无首屏走查〕 | `apps/vscode/package.json:58-109` |

### 1.B 视图标题栏按钮（两家同层不同名，且 Cline 缺一枚 MCP）

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| 视图标题栏 | `New Task`（`$(add)`，位置 `navigation@1`） | 清空当前任务、开新对话 | 没有〔已证：C17 名册 4 条里没有"开新任务"，`internal/panel/bridge.go:42-45`〕 | `apps/vscode/package.json:129-135`、`:256-261`；处理器 `apps/vscode/src/extension.ts:132-139` |
| 视图标题栏 | `Customize`（`$(wrench)`，`navigation@2`） | 打开"技能／MCP／插件"那一屏 | 没有那一屏；配置键有（`internal/config/schema.go:556-573` 的插件开关与白名单） | `apps/vscode/package.json:141-145`、`:262-266` |
| 视图标题栏 | `History`（`$(history)`，`navigation@3`） | 历史任务列表 | 没有〔已证：`internal/session/` 只有 `doc.go` 一枚，无列表出口〕 | `apps/vscode/package.json:146-150`、`:267-271` |
| 视图标题栏 | `Account`（`$(account)`，`navigation@5`） | 账户／额度／订阅 | 没有那一屏 | `apps/vscode/package.json:151-155`、`:272-276` |
| 视图标题栏 | `Settings`（`$(settings-gear)`，`navigation@6`） | 打开扩展自带设置页（不是 VS Code 设置页） | 有配置、无那一屏（`internal/config/schema.go` 全量 section 树） | `apps/vscode/package.json:156-160`、`:277-281` |
| 视图标题栏 | **`MCP Servers` 不在标题栏**〔已证〕 | 命令注册了（`apps/vscode/src/extension.ts:140`）、`contributes.commands` 里有（`:136-140`，标题原文 `MCP Servers`），**但 `menus.view/title` 那一段里没有它**（那一段只有 5 枚，且 `navigation@4` 直接空号） | — | `apps/vscode/package.json:136-140` 对比 `:255-281` |
| 视图标题栏 | `New Session`（`$(add)`，`navigation@1`） | 新会话 | 同上，没有 | `packages/vscode/package.json:212-233`＋nls `command.newSession.title` |
| 视图标题栏 | `Open Session in Editor`（`$(link-external)`，`navigation@2`） | 把当前会话从侧栏**搬到编辑器标签页**（可多开、可并排） | 没有这种"同一会话两个容器"的形态 | `packages/vscode/package.json:212-233`；实现 `packages/vscode/src/SessionEditorPanelProvider.ts:143-158` |
| 视图标题栏 | `Run on Several Models`（`$(repo-forked)`，`navigation@3`） | 开一枚"草稿标签页"，输入区切到并行模式（同一条提示词丢给几枚模型） | 我们拍板的是多子代理（`A399`/`A401` 改判），**不是多模型并列**；界面上两者都没有 | `packages/vscode/package.nls.json` `command.openAgentManager.title`；实现 `packages/vscode/src/extension.ts:233-238`、`src/SessionEditorPanelProvider.ts:117-124` |
| 视图标题栏 | `Settings`（`$(settings-gear)`，`navigation@4`） | 在 webview 里跳到设置视图（不是原生设置页） | 没有那一屏 | `packages/vscode/package.json:212-233` |

### 1.C 编辑器／资源管理器／终端／Git／笔记本的菜单挂载点

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| 编辑器右键 | `Add to Cline`（`when: editorHasSelection`） | 把选区作为上下文塞进输入框 | 〔界面地界·不读〕Go 无"选区/文件引用"通道（`internal/panel/bridge.go:42-45` 四条方法名册里没有） | `apps/vscode/package.json:173-177`、`:283-288` |
| 终端右键 | `Add to Cline` | 把终端最近输出塞进输入框 | 没有 | `apps/vscode/package.json:178-182`、`:290-294`；取值实现 `apps/vscode/src/hosts/vscode/terminal/get-latest-output.ts` |
| Git 面板标题栏 | `Generate Commit Message with Cline`（自定义图标 `$(cline-icon)`）／生成中换成 `Generate Commit Message with Cline - Stop`（`$(debug-stop)`） | **同一位置两枚按钮按状态互斥切换**，靠上下文键 `cline.isGeneratingCommit` 翻 | 没有；我们有 `internal/proc`（进程）但没有 SCM 面板宿主 | `apps/vscode/package.json:188-199`、`:296-306`；翻键 `apps/vscode/src/hosts/vscode/commit-message-generator.ts:195,258,264` |
| 笔记本工具栏 | `Generate Jupyter Cell with Cline`（`$(sparkle)`） | 在 notebook 里让 AI 生一个 cell | 没有 | `apps/vscode/package.json:210-215`、`:308-313` |
| 单元格标题（行内） | `Explain Jupyter Cell with Cline`（`$(question)`）／`Improve Jupyter Cell with Cline`（`$(lightbulb)`） | 每个 cell 标题处两枚行内按钮（`inline@1`／`inline@2`） | 没有 | `apps/vscode/package.json:216-227`、`:315-325` |
| 编辑器右键二级菜单 | `OpenChamber` → `Explain`／`Improve Code`／`Add to Context`／`Add Comment` | 一个二级菜单收纳 4 项 | 没有 | `packages/vscode/package.json:162-167`＋`:234-248`＋`:193-198` |
| 资源管理器右键 | `Attach to OpenChamber Chat`（`when: resourceScheme == file && explorerResourceIsFolder == false`） | 把一个**文件**（不是文件夹）挂到当前会话 | 有受理器、没人把结果发出去：`internal/panel/pump.go:226` 仍传空实参（第四版已记〔建了但没接〕） | `packages/vscode/package.json:199-205`；nls `command.attachExplorerToChat.title` |
| 编辑器标题栏 | `Open New Session In Editor`（图标分深浅两版 `assets/icon.svg`/`icon-titlebar.svg`） | 常驻的一枚"开新标签页"入口 | 没有 | `packages/vscode/package.json:61-161`（命令定义）＋`:206-211`（挂到编辑器标题栏） |

### 1.D 编辑器里的"灯泡"（代码操作）——Cline 独有的第四入口

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| 光标处的 Quick Fix 列表 | `Add to Cline`（Kind=QuickFix） | 任何时候都出现 | 没有 | `apps/vscode/src/extension.ts:268-271` |
| Refactor 列表 | `Explain with Cline`（Kind=RefactorExtract） | 出现于重构建议里 | 没有 | `:276-280` |
| Refactor 列表 | `Improve with Cline`（Kind=RefactorRewrite） | 出现于改写建议里 | 没有 | `:284-288` |
| Quick Fix 列表 | `Fix with Cline`（`isPreferred: true`，**只有当该处有诊断错误时才加**） | 把"这个报错你帮我修"一次点出；`isPreferred` 意味着 VS Code 会在灯泡上直接执行它 | 没有 | `:291-300`（判据逐字 `if (context.diagnostics.length > 0)`）；处理器 `:325-332` |
| 声明的贡献种类 | `providedCodeActionKinds`: QuickFix／RefactorExtract／RefactorRewrite | 提前声明，避免 VS Code 空跑一次提供器 | — | `:306-311` |

⚠ 注释里写明为什么处理器要**自己重新取选区与诊断**（`:260-264`）：VS Code 的 `CommandsConverter` 缓存会在 code action 列表销毁时把命令拆掉，直接执行会报 `Actual command not found`。〔已证〕——这是一条"宿主 API 的坑"，抄形状之前要先读这段。

### 1.E 命令面板能调出什么（逐枚；两家的"藏"法不同）

Cline 19 枚（`apps/vscode/package.json:129-232`，标题逐字）：
`New Task`／`MCP Servers`／`Customize`／`History`／`Account`／`Settings`（6 枚按钮型）、`Create Test Tasks` 与 `Expire MCP OAuth Tokens (for testing)`（**`when: cline.isDevMode`，只开发态可见**，`:162-171`）、`Add to Cline`（编辑器）、`Add to Cline`（终端，同名不同 id）、`Jump to Chat Input`、`Generate Commit Message with Cline`、`Generate Commit Message with Cline - Stop`、`Explain with Cline`、`Improve with Cline`、`Generate Jupyter Cell with Cline`、`Explain Jupyter Cell with Cline`、`Improve Jupyter Cell with Cline`、`Open Walkthrough`。
另外 `commandPalette` 那一段（`:327-336`）只把两枚 Git 命令按 `cline.isGeneratingCommit` 互斥显示——**其余 17 枚没写 `commandPalette` 规则，即默认可见**。

OpenChamber 17 枚（`packages/vscode/package.json:61-161`，标题逐字）：
`Open Sidebar`／`Focus Chat`／`Restart API Connection`／`Show OpenCode Status`／`Run on Several Models`／`Open Active Session in Editor`／`Open New Session in Editor`／`Open Session in Editor`／`Add to Context`／`Add Comment`／`Comment`／`Remove Comment`／`Explain`／`Improve Code`／`New Session`／`Settings`／`Attach to OpenChamber Chat`。
**`commandPalette` 里显式写了 `"when": "false"` 两枚**＝`Comment`（提交行评论）与 `Remove Comment`：命令存在、面板调不出、只能从评论线程里点。〔已证〕`packages/vscode/package.json:183-192`。

诊断类命令 `Show OpenCode Status` 一次列出 server URL／模式／启动重启就绪计数／检测端口／API 前缀／上次退出码／连接时长（上一轮已记，本腿只补一条：它的文案不在 `package.nls.json` 的 28 条里，说明**内容不是本地化的**，出处 `packages/vscode/package.nls.json` 全量 28 键 vs `src/extension.ts:755-795`）。〔已证〕

### 1.F 键位

| 键 | Cline | 说明 | OpenChamber |
|---|---|---|---|
| `ctrl+'`（mac `cmd+'`） | `Add to Cline` **when `editorHasSelection`** | **同一枚键位按"有没有选区"分成两个命令**：有选区→加上下文，没选区→`Jump to Chat Input` | 无键位贡献（`contributes.keybindings`＝`null`） |
| `ctrl+'` | `Jump to Chat Input` **when `!editorHasSelection`** | 同上，第二支 | — |
| 无默认键 | `Generate Commit Message with Cline` | 只挂了 `when: config.git.enabled && scmProvider == git`，键位让用户自己在 VS Code 里设 | — |

出处：`apps/vscode/package.json:233-254`。⚠ 我方对照：`internal/config/schema.go:180-183` 有四枚热键键、**没有可改键的那一屏**（第四版 §2 已记〔建了但没接〕，本腿不重述）。

### 1.G diff 编辑器怎么呈现改动（本块最厚的一格）

| 位置 | 控件名／形状 | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| 虚拟文档 scheme | `cline-diff`（`DIFF_VIEW_URI_SCHEME`） | 内容走 **URI 的 base64 query**，不可变；两侧都用它 ⇒ 能开"从未存在过的文件"的 diff | 没有虚拟文档这一层；`internal/panel` 里 `diff` 只出现在审批与 spill 相关处（尺：`grep -rn "diff" internal/panel internal/tools --include=*.go \| grep -v _test \| wc -l` ＝ **30**，全在文本与审批里，无宿主文档 provider） | `apps/vscode/src/hosts/vscode/VscodeDiffContentProvider.ts:4`、注册 `src/extension.ts:149` |
| 第二套 scheme | `cline-edit-preview`（`EDIT_PREVIEW_URI_SCHEME`） | **可变**虚拟文档：`set()` 触发 `onDidChange`，VS Code 就地重渲染同一侧 | 没有 | `apps/vscode/src/hosts/vscode/VscodeEditPreview.ts:6,24-41` |
| diff 标签页 | 一条 `vscode.diff` 调用，`{preview:false, preserveFocus:true}` | 开**只读预览**：真实文件从不被打开或改动；关掉的一定是预览、不会误关源文件 | 〔界面地界·不读〕 | `apps/vscode/src/hosts/vscode/VscodeEditPreview.ts:92-95`，设计注释逐字 `:62-70` |
| 流式动画 | 无文字标签；黄色覆盖层＋高亮当前行 | **模拟打字**：右侧先显示原文，然后新内容一行行打进去；起手停在文件顶（`ANIMATION_START_BEHIND` 注释在 `:8-13`）、到底后 250ms 再跳到第一处改动居中 | 我们的球／面板没有"在编辑器里演一遍"这一层；`internal/streamkey` 只有关键字流 | `apps/vscode/src/hosts/vscode/VscodeEditPreview.ts:107-168`；装饰色逐字 `rgba(255,255,0,0.1)`/`opacity 0.4` 与 `rgba(255,255,0,0.3)`＋1px 边框 `src/hosts/vscode/DecorationController.ts:3-14` |
| 大文件闸门 | 无文字 | **超过 3000 行不演**，直接落最终 diff，只用一次前缀扫描把视口瞄准第一处差异 | 没有这一形 | `apps/vscode/src/hosts/vscode/VscodeEditPreview.ts:16-17,114-120` |
| 视口等待 | 无文字 | `vscode.diff` 返回后编辑器不一定立刻出现在 `visibleTextEditors`，故**最多重试 10 次找那一侧** | — | `apps/vscode/src/hosts/vscode/VscodeEditPreview.ts:170-180` |
| 多文件一次开 | `vscode.changes`（VS Code 原生"更改"多文件审阅页）＋立刻 `workbench.action.closePanel` 收起底部面板腾地方 | 一批改动开成**一个多标签审阅页**，每格仍是虚拟文档对比 | 没有 | `apps/vscode/src/hosts/vscode/hostbridge/diff/openMultiFileDiff.ts:9-31` |
| **VS Code 侧没做的两件事** | `openDiff`／`closeAllDiffs` 在 VS Code 宿主里**直接抛错**：逐字 `openDiff is not supported by the VS Code diff service.` | 说明这两条是给别的宿主（JetBrains／CLI）准备的，扩展形态**不走这条路** | — | `apps/vscode/src/hosts/vscode/hostbridge/diff/openDiff.ts:3-5`、`closeAllDiffs.ts:3-5` |
| webview 内的差异行 | 折叠条：文件图标按动作分色（`Add`＝`FilePlus`／`Delete`＝`FileX`／改＝`FileText`，左边框同色）＋文件路径（悬停下划线、`title="Open file in editor"`）＋`+N · -N` 统计＋一枚跳转图标 | **聊天里先看摘要、点开展开逐行**（限高 80、可横滚、带行号列） | 〔界面地界·不读〕 | `apps/vscode/webview-ui/src/components/chat/DiffEditRow.tsx:27-29,139-166,193-201` |
| OpenChamber 的 diff | 只注册**一枚**虚拟文档 provider（`bridge-system-runtime.ts:135`），开 diff 走 `vscode.diff(originalUri, modifiedUri, title)`（`:367`） | webview 想要"在宿主里开个对比"时代开 | 没有 | `packages/vscode/src/bridge-system-runtime.ts:135,367` |
| OpenChamber 的 git 差异 | `src/gitPathDiff.ts`（按路径算差异、供逐块应用） | 上一轮已记 `apply-hunk`；本腿补位置：**它没有界面**，界面在 web 那一侧 | 我们的 git 出口还停在"连换分支都不通"：`internal/panel/git.go:82` 那句实话 | `packages/vscode/src/gitPathDiff.ts`；对照第四版 `docs/reports/missing-features-2026-09-29-v4.md:380` |

### 1.H 评论线程（"行评论"这一形：一家接了、一家没接）

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| OpenChamber | 控制器 id `openchamber.inlineComments`，线程标题 `Comment on line {0}`／`Comment on lines {0}-{1}`，作者名 `OpenChamber` | 在编辑器**选区或行号槽**起一枚原生评论线程；评论＝**待发上下文卡**，随下一条消息一起发出去 | 没有：我们没有编辑器宿主；且 `internal/panel/bridge.go:42-45` 名册里没有"引用某几行"这一条 | `packages/vscode/src/InlineCommentThreads.ts:102`；文案 `packages/vscode/l10n/bundle.l10n.json`（28 键全量） |
| OpenChamber | 线程未提交时显示 `Not sent yet` | **告诉用户这条评论还没进模型** | 没有 | `bundle.l10n.json` 键 `Not sent yet` |
| OpenChamber | 线程标题按钮 `Remove Comment`（`$(close)`，仅当 `commentThread == openchamberAttached`） | 摘掉这枚待发评论 | 没有 | `packages/vscode/package.json:176-182` |
| OpenChamber | 评论框内联按钮 `Comment`（`group: inline`） | 提交＝把这条塞进会话 | 没有 | `packages/vscode/package.json:169-175` |
| OpenChamber | **失败老实说**：`OpenChamber [Add Comment]: The comment never reached the chat and was discarded` | 宿主把评论送进 webview、**没收到回执**就明说"已丢弃"，不假装挂上 | 同脾气先例只此一处：`internal/panel/git.go:82`；评论这条**没有** | `packages/vscode/src/extension.ts`（该串出现处，尺：`grep -rn "never reached the chat" packages/vscode/src --include=*.ts \| grep -v test`） |
| OpenChamber | 别的守卫句：`[Add Comment]: No active editor`／`File is outside the workspace` | 每种点不着的原因各一句 | 没有 | `packages/vscode/l10n/bundle.l10n.json` |
| Cline | 控制器 id `cline-ai-review`，label `Cline AI Review`，作者名 `Cline`，**流式占位文案 `_Thinking..._`**（逐字，Markdown 斜体） | 设计上是"AI 边想边在行上打字"：先建线程放占位、`appendToStreamingComment` 逐段刷 | 没有 | `apps/vscode/src/hosts/vscode/review/VscodeCommentReviewController.ts:26`、`:63-70`、`:104-115`、`:170-188` |
| Cline | ⚠ **〔建了但没接〕**：这套评论线程**没有任何调用者** | 抽象基类（`addReviewComment`／`startStreamingComment`／`appendToStreamingComment`）＋两枚实现都在，工厂也已挂上 `HostProvider`（`apps/vscode/src/extension.ts:572`），**但全仓没人调这三枚方法** | — | 尺（本机复跑）：`cd /tmp/v5b/cline-fresh && grep -rn "addReviewComment\|startStreamingComment\|appendToStreamingComment\|ReviewComment" apps/vscode/src sdk --include=*.ts` 去掉两处实现与抽象基类后**命中 0 行** |
| Cline | ⚠ 它会**改用户的 VS Code 全局设置**：把 `comments.openView` 写成 `"never"` | 防止加评论时 Comments 面板自己弹出来抢地方 | 我们没有这种"顺手改宿主全局配置"的动作 | `apps/vscode/src/hosts/vscode/review/VscodeCommentReviewController.ts:34-42`（`ConfigurationTarget.Global`） |

### 1.I OS 通知（这一格两家完全相反）

| 位置 | 控件名／原文 | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| OpenChamber | 4 类通知模板，英文默认句：`{agent_name} is ready` ＋ `{model_name} completed the task`／`Tool error` ＋ `{last_message}`／`Input needed` ＋ `{last_message}`／子代理完成另有一套（模板键 `completion`/`subtask`/`error`/`question`） | 标题与正文**用户可编辑**，可用变量共 8 枚：`project_name`／`worktree`／`branch`／`session_name`／`agent_name`／`model_name`／`last_message`／`session_id` | 我们的通知是**托盘气泡**、一句写死的标题＋正文，且只有 `wisp run` 一条生产路径（`cmd/wisp/run.go:161-163` 逐字 `s.notify = postSystemNotification`；`cmd/wisp/notify_windows.go:6-13` 逐字"It is a tray balloon tip, not a toast"）；**模板、变量、按类开关全没有** | `packages/vscode/webview/main.tsx:1929-1939,1947-1965,2009-2017,2027-2033,2044-2050` |
| OpenChamber | 闸门：`nativeNotificationsEnabled`／`notificationMode`（≠`always` 时 `requireHidden`）／`notifyOnCompletion`／`notifyOnSubtasks`／`notifyOnError`／`notifyOnQuestion` | **子代理是否也推**是一枚独立开关；"只看别的端在不在"是另一枚 | 〔已证·同向〕我们的等批准队列长度已现成（`internal/agent/approval/queue.go`）但"数上去／推出去"没有；四枚开关一枚都没有 | 同上 `:1984,1987,2004,2022,2041` |
| OpenChamber | 同会话 **5 秒冷却**（`READY_NOTIFICATION_COOLDOWN_MS = 5000`），完成与错误各一张表 | 防"一件事连响两次" | 没有 | `packages/vscode/webview/main.tsx:1842,2005-2008,2023-2026` |
| OpenChamber | **点通知＝直达那一条会话**：`notification.onclick` 设 `setCurrentSession(sessionId)` 再广播 `openchamber:navigate {view:'chat'}` | 不新建会话、不弹选择 | 没有 | `packages/vscode/webview/main.tsx:1806-1814` |
| OpenChamber | "这张通知该由哪个面板发"要先 **claim**（`claimOpenChamberNotification`），抢不到就不发；窗口焦点位 `window.__OPENCHAMBER_VSCODE_WINDOW_FOCUSED__` 由宿主写（`onDidChangeWindowState`） | 多面板/多标签同开时**不重复弹**；"人就在窗口前"就不弹 | 没有 | `packages/vscode/webview/main.tsx:1783-1786,1800`；宿主侧 `packages/vscode/src/extension.ts:226` |
| Cline | **没有 OS 通知**〔已证〕 | `grep -rn "new Notification\|navigator.requestPermission\|showNotification" apps/vscode/webview-ui/src apps/vscode/src` 去测试后 0 命中 | — | 同上尺 |
| Cline | ⚠ **〔建了但没接〕**：设置里有个 `enableNotifications` 布尔位（注释逐字 `Show notifications for approval and task completion`） | 字段有默认值 `false`、有三条写入路径、有一枚读取钩子，**但名册里没这一项、界面画不出来、宿主也没有发通知的实现** | 与我们的"写了不管用就响亮拒收"正好相反：它**不响亮** | `apps/vscode/src/shared/AutoApprovalSettings.ts:25,43`；写入 `src/core/controller/state/updateAutoApprovalSettings.ts:22`、`src/core/controller/task/newTask.ts:32-33`；读取钩子 `webview-ui/src/hooks/useAutoApproveActions.ts:14-15,29,59-63`；**渲染名册只有 5 项、不含它** `webview-ui/src/components/chat/auto-approve-menu/constants.ts:3-35` |
| Cline | 它只有一条 `showInformationMessage`／`showWarningMessage`／`showErrorMessage` 的**通用转发器**（`window/showMessage.ts:15-21`） | SDK 要弹宿主对话框时走这里；不主动弹 | 我们有一条对应的桥：`docs/specs/SPEC-08*.md` C17 四条方法名册（`internal/panel/bridge.go:42-45`），**不含"弹宿主框"** | `apps/vscode/src/hosts/vscode/hostbridge/window/showMessage.ts:15-21` |

### 1.J 输出面板／终端／状态栏

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| 输出面板 | `Cline`（一条频道） | 调试日志落点 | 我们有 `internal/observe`，但没有"人可翻的那一条频道" | `apps/vscode/src/hosts/vscode/hostbridge/env/debugLog.ts:4` |
| 输出面板 | `OpenChamber`／`OpenChamberManager`（**两条**：一条给扩展、一条给被管的 opencode 进程） | 把"外部进程的 stdout"和"扩展自己的日志"分开两条 | 〔已证·同向〕我们有 `internal/proc`＋`internal/observe`，但没有分两条给人看 | `packages/vscode/src/extension.ts:56`、`src/opencode.ts:29` |
| 终端 | 固定名 `Cline` 的终端，靠 **OSC 633 shell 集成标记**判断命令起止、复用与驱逐 | "在真终端里跑"而不是后台管道；标记丢了另有兜底（注释逐字提到无标记流如 SSH） | 我们用 `internal/proc`，没有 shell 集成标记这一层 | `apps/vscode/src/hosts/vscode/terminal/VscodeTerminalRegistry.ts:25-44`、`VscodeTerminalManager.ts:213,229,237,308,382`、`osc633Parser.ts` |
| 状态栏 | **两家都没有状态栏项**〔已证〕 | — | 我方由**悬浮球**承担这一格（`internal/ball/statevisual.go`），形状不同：我们是常驻窗口不是底部一行 | 尺（本机复跑）：`grep -rni "statusbar" apps/vscode/src`（去掉 `src/test/vscode-mock.ts:63` 的测试假件后 0 命中）／`grep -rn "createStatusBarItem" packages/vscode/src`（0 命中） |

### 1.K webview 那一屏里面的控件（逐件，原文标签）

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| webview 顶部导航 | 5 枚图标按钮，tooltip 逐字 `New Task`／`Customize`／`History`／`Account`／`Settings` | **与视图标题栏那 5 枚同名同去向**（同一屏两个入口） | 没有这一重复；我们只有一条 C17 桥 | `apps/vscode/webview-ui/src/components/menu/Navbar.tsx:15-62`（`data-testid` 为 `tab-<id>`） |
| webview 任务头 | 折叠箭头（aria `Expand task header`／`Collapse task header`）＋任务标题＋工作目录徽章＋`$0.0000` 价签 | 点标题展开/收起细节 | 没有价签这一枚 | `apps/vscode/webview-ui/src/components/chat/task-header/TaskHeader.tsx:137,171-186` |
| 任务头按钮族 | `Start a New Task`（aria `New Task`）／`Copy Text`（aria `Copy`）／`Delete Task`／`Compact Task`／（仅开发态）`Open Disk Conversation History` | 展开时才出现 5 枚小图标 | 〔已证〕我们有取消入口四条并各说实话（`internal/agent/approval/approval.go:106-114`），**"压缩这一条任务"的按钮没有** | `apps/vscode/webview-ui/src/components/chat/task-header/buttons/{NewTaskButton,CopyTaskButton,DeleteTaskButton,CompactTaskButton,OpenDiskConversationHistoryButton}.tsx:12-33`；开发态闸门 `TaskHeader.tsx:156` |
| 上下文窗口条 | 进度条 aria `Context window usage progress`；弹出的确认框两枚按钮标题逐字 `Yes, compact the task`／`No, keep the task as is` | **压缩前二次确认**，且确认框写成人话正反两句 | 〔已证〕有压缩器与预算（`internal/agent/compress.go`、`internal/agent/budgets.go`），**没有那一条进度显示、也没有二次确认**；⚠ 且 `docs/reports/missing-features-2026-09-29-v4.md:238` 记过"20 枚状态名里没有『正在压缩』这一名" | `apps/vscode/webview-ui/src/components/chat/task-header/ContextWindow.tsx:39-52,176` |
| 输入框 | placeholder 两句逐字 `Type a message...`／`Type your task here...`（有任务与无任务各一句） | — | 〔界面地界·不读〕 | `apps/vscode/webview-ui/src/components/chat/ChatView.tsx:373-375` |
| 输入框左下角 | `@` 引用弹层，条目原文：`Problems`／`Terminal`／`Paste URL to fetch contents`／`No results found`，外加文件／文件夹／git 提交三类检索（枚举 `File`/`Folder`/`Problems`/`Terminal`/`URL`/`Git`/`NoResults`） | 六类上下文来源一屏选 | 〔建了但没接〕只有文件/图片受理（`internal/panel/attachments.go`；送出仍传空实参 `internal/panel/pump.go:226`） | `apps/vscode/webview-ui/src/utils/context-mentions.ts:82-90`、`components/chat/ContextMenu.tsx:96-102`、`components/chat/ChatTextArea.tsx:293-296` |
| 输入框右下角 | 模式滑条两档原文 `Plan`／`Act`，悬停解释句逐字 `In Act  mode, Cline will complete the task immediately` / `In Plan  mode, Cline will gather information to architect a plan`，并附 `Toggle w/ <kbd>` 键提示 | **档位旁边永远写一句"这一档会发生什么"＋怎么切** | 〔已证〕我们有档位默认最严（`internal/risk/mode.go:44-56`、`internal/config/schema.go:470`），**没有那句解释、没有那枚切换出口** | `apps/vscode/webview-ui/src/components/chat/ChatTextArea.tsx:1682-1700` |
| 输入框上方 | `Auto-approve:` ＋ 已开项的短名串（`None` 也算一句），点开是逐条开关 | 见 §5（这是"长期"那一档的家） | 〔已证〕没有这一枚条 | `apps/vscode/webview-ui/src/components/chat/auto-approve-menu/AutoApproveBar.tsx:109,44`；`ChatView.tsx:416` |
| 审批行（底部按钮区） | 逐态原文两枚：`Approve`／`Reject`；改文件时换成 `Save`／`Reject`；跑命令时 `Run Command`／`Reject`；浏览器与 MCP 与子代理各是 `Approve`／`Reject` | **按钮名随"要干什么"变，不是永远 Approve** | 〔已证〕面板只能拒、"允许"被结构禁掉（`docs/specs/SPEC-08*.md:156-176`），第四版 `missing-features-2026-09-29-v4.md:85` 已记本格；名册出处 `apps/vscode/webview-ui/src/components/chat/chat-view/shared/buttonConfig.ts:46-105` | 同上 |
| 审批行 | `Proceed While Running`（命令还在跑时）／`Cancel`／`Retry`／`Start New Task`／`Proceed Anyways`／`Resume Task`／`Start New Task with Context`／`Condense Conversation`／`Report GitHub issue` | 每种"等一件事"各有专属动词 | **没有这套动词**〔已证：我方状态名册 20 枚里没有对应动词出口〕；最接近的是四条不可用时的实话句（`internal/agent/approval/approval.go:106-114`） | `apps/vscode/webview-ui/src/components/chat/chat-view/shared/buttonConfig.ts:28-213` |
| 子代理行 | 标题逐字 `Cline wants to use a subagent:`／`Cline wants to use subagents:`；每条一行状态图标＋提示词两行截断＋`Show more`；统计句逐字 `N tools called · N tokens · $X.XX`；`Show output`／`Hide output`；跑着的时候只有一行 `latestToolCall` 文本；没数据时逐字 `Subagent status update unavailable.` | 见 §4（这是必答一的证据） | 〔已证〕我们有子代理名册与流键件（`internal/panel/subagent_roster_197.go`、`internal/streamkey`），**界面上那一行没有** | `apps/vscode/webview-ui/src/components/chat/SubagentStatusRow.tsx:200,186,232,250-270,256-266`；状态枚举 `apps/vscode/src/shared/ExtensionMessage.ts:306-322` |
| 历史页 | 排序与过滤原文：`Newest`／`Oldest`／`Most Expensive`／`Most Tokens`／`Most Relevant`／`Workspace Only`／`Favorites Only`；分组标题 `Today`／`Older`；页头 `History` | 七种排法，含**按钱**与**按相关性** | 没有〔已证：`internal/session/` 只有 `doc.go`〕 | `apps/vscode/webview-ui/src/components/history/HistoryView.tsx:29-35,348-351,384` |
| 自定义页 | 页头 `Customize`；三族原文 `Skills`／`MCP`／`Plugins`（单复数各一句：`MCP server`／`MCP servers`）；两标签 `Installed`／`Marketplace`；来源标签 `Global`／`Workspace`／`Remote`；逐条开关 tooltip `Disable <名>`／`Enable <名>` | 一个"能力商店＋本机已装"双页 | 〔已证〕只有配置键、没有那一屏（`internal/config/schema.go:556-573`） | `apps/vscode/webview-ui/src/components/marketplace/MarketplaceView.tsx:52-125,858-862,1227` |
| 工作树页 | 失败句逐字 `Failed to load worktrees`／`Failed to delete worktree`／`Failed to merge worktree`／`Failed to create task for Cline`；有 `CreateWorktreeModal`／`DeleteWorktreeModal` 两枚弹层 | 在扩展里管理工作树（建新树、删、合并、给某棵树开任务） | 〔已证·同方向但没通〕`internal/panel/git.go:82` 那句"切换分支／切换工作树今天不可用"；⚠ 视图标题栏入口在 Cline 那边**不存在**（见下一行） | `apps/vscode/webview-ui/src/components/worktrees/WorktreesView.tsx:76,147,213,235` |
| ⚠ 工作树按钮 | `cline.worktreesButtonClicked`：处理器已注册（`apps/vscode/src/extension.ts:147`）、订阅通道也在（`src/core/controller/ui/subscribeToWorktreesButtonClicked.ts:16-56`），**但 `contributes.commands` 与 `menus` 里都没有这一枚**（`package.json` 里搜 `worktree` 0 命中） | 〔建了但没接〕（对"扩展宿主入口"这一层而言） | — | 尺（本机复跑）：`grep -n -i worktree apps/vscode/package.json` → 0 命中；`grep -rn "WorktreesButton" apps/vscode/src` → `extension.ts:147`、`registry.ts:19` |
| OpenChamber 侧栏之外 | 会话可开成**编辑器标签页**，标题分别是 `New Session`／`Run on several models`／`Session`（拿不到会话名时的兜底）；标签页带深浅两版图标 | 同一会话在两个容器里都能有一份 | 没有 | `packages/vscode/src/SessionEditorPanelProvider.ts:113-124,143-158,162-176`；标题文案 `l10n/bundle.l10n.json` |
| OpenChamber 兜底句 | `OpenChamber: No folder is open. Open a folder to start a new session.` | 没开文件夹时**不猜**上一个目录、明说要去开 | 同脾气先例：`internal/panel/git.go:82`；这一格没有 | `packages/vscode/src/SessionEditorPanelProvider.ts:125-137`（注释逐字解释"否则会退回共享状态里上一条会话的目录（这就是它修的 bug）"） |
| OpenChamber 首屏 | 加载屏＝**一枚会呼吸的 logo 立方体**（无文字即代表在干活）＋一条状态行；注释逐字 `Only show text when something is wrong — progress states stay silent`；且 `@media (prefers-reduced-motion: reduce)` 时关掉动画 | "进度不写字、错误才写字"＋尊重系统减少动效设置 | 〔已证·可对照〕我们的球本来就是这个思路（20 态视觉），但 `prefers-reduced-motion` 这一条我们没有 | `packages/vscode/src/webviewHtml.ts:116-125,269-286` |
| OpenChamber 首屏文案表 | 7 句 × 3 语言（en/fr/tr）内联在宿主 HTML 里：`Starting OpenCode API…`／`Initializing…`／`Connecting…`／`Connected!`／`Connection error`／`Reconnecting…`／`OpenCode CLI not found. Please install it first.` | 模块脚本还没加载前就得有话可说 | 没有首屏文案表〔已证：`internal/panel/assets.go` 无该层〕 | `packages/vscode/src/webviewHtml.ts:255-299`；另有一份**权威表**在 `packages/ui/src/lib/i18n/bootstrap.ts:232-290`（含 `disconnected`、`startingDevServer`、`waitingDevServer`、`loadingData`，语言更多）⇒ 同一族文案两处各拼一份，见 §7 |

### 1.L "打开侧栏 / 聚焦输入框"这类跨窗口动作（宿主与网页各持一半）

| 位置 | 控件名（原文） | 它干什么 | 我方对应物 | 出处 |
|---|---|---|---|---|
| Cline | `Jump to Chat Input`（`ctrl+'` 的空选分支） | 从编辑器一键把焦点丢进侧栏输入框 | 没有〔已证：C17 四条名册无焦点类方法〕 | `apps/vscode/package.json:183-187`、`src/extension.ts:353` |
| Cline | `cline://...` URI：`onUri` 激活、任务深链路径常量两条、处理失败**只记警告不报错** | 从浏览器/别的程序把一条任务塞进侧栏 | 〔已证〕全仓无深链协议（第四版 `missing-features-2026-09-29-v4.md:134` 同款负向读数） | `apps/vscode/package.json:42-46`、`src/extension.ts:158-179` |
| OpenChamber | `Open Sidebar`／`Focus Chat`／`Restart API Connection` | 三枚"宿主侧救生圈"，尤其第二枚是 webview 卡住时的出口 | 没有；我们的等价物只有球的重启路径 | `packages/vscode/package.nls.json`（`command.openSidebar.title` 等三条） |
| OpenChamber | `Show OpenCode Status` | 把"现在连的是哪、哪个模式、起来几次、上次退出码"一次列全 | 〔已证·无〕 | `packages/vscode/package.nls.json`；实现 `src/extension.ts:755-795`（上一轮已记，本腿只补"文案未本地化"） |

---

## §2 表二 · 设置页逐项表

### 2.0 先定"设置面在哪"这一格（两家都不在 VS Code 设置页）

| 家 | VS Code 原生设置页里的条目数 | 其余设置画在哪 | 出处与尺 |
|---|---|---|---|
| Cline | **0 项**〔已证〕——`contributes.configuration` 只剩一个标题 | webview 自己的 Settings 视图，7 个标签 | `apps/vscode/package.json:338-341`（逐字 `"properties": {}`）；标签名册 `apps/vscode/webview-ui/src/components/settings/SettingsView.tsx:52-99` |
| OpenChamber | **2 项**〔已证〕 | webview（＝它自家 web 应用）里的设置页；扩展独有的只有那两枚＋通知页的一句提示 | `packages/vscode/package.json:249-264`；名册 `packages/ui/src/components/sections/openchamber/OpenChamberVisualSettings.tsx:307` |

尺（本机复跑，两条都跑过）：
- `python -c "import json;d=json.load(open('apps/vscode/package.json'));print(len(d['contributes']['configuration']['properties']))"` ＝ **0**
- `python -c "import json;d=json.load(open('packages/vscode/package.json'));print(len(d['contributes']['configuration']['properties']))"` ＝ **2**
- OpenChamber 外观/行为一页的"可见设置名"名册＝**50 枚**（口径：把 `OpenChamberVisualSettings.tsx:307` 那行 `type VisibleSetting = ...` 的联合类型逐名抽出、去重后计数；**这是"这一页可能被画出来的项目数"，不是"画出来的项目数"**——后者还要各 `shouldShow(...)` 与运行时开关）

### 2.1 Cline 的设置页逐条（标签原文＝界面字面，逐条抄）

页头原文 `Settings`（`apps/vscode/webview-ui/src/components/settings/SettingsView.tsx:253`）；7 个标签的原文名与提示：`API Configuration`／`Features`（提示 `Feature Settings`）／`Terminal`（`Terminal Settings`）／`General`（`General Settings`）／`Remote Config`（`Remotely configured fields`）／`About`（`About Cline`）／`Debug`（`Debug Tools`）——均出自 `SettingsView.tsx:52-99`。

| 设置项标签原文 | 改了会怎样 | 我方对应物 | 出处 |
|---|---|---|---|
| `Allow error and usage reporting` ＋正文 `Help improve Cline by sending usage data and error reports. No code, prompts, or personal information are ever sent.` | 遥测开关；**被组织远程配置管住时打不掉并挂一枚锁图标**（tooltip 逐字 `This setting is managed by your organization's remote configuration`） | **没有这一项**〔已证：`internal/config/schema.go` 搜 `telemetry` 0 命中〕；"被上级管住"这一形我们也没有（远程配置属未定案） | `apps/vscode/webview-ui/src/components/settings/sections/GeneralSettingsSection.tsx:23-46,30-31,37-39` |
| 界面语言（`PreferredLanguageSetting` 一件） | 选对话/命令输出用哪种语言 | **有**：`internal/config/schema.go:138` `app.language` 默认 `zh-CN` | `apps/vscode/webview-ui/src/components/settings/PreferredLanguageSetting.tsx`（由 `GeneralSettingsSection.tsx:19` 引入） |
| `Auto Compact` ／描述 `Automatically compress conversation history.` | 打开后自动压上下文 | **有压缩器、没有这枚开关**〔已证：`internal/agent/compress.go` 在产码里；`schema.go` 搜 `compact`/`condense` 0 命中〕 | `apps/vscode/webview-ui/src/components/settings/sections/FeatureSettingsSection.tsx:33-40` |
| `Auto Compact Strategy`：两档 `Basic`／`Agentic` | 压缩策略选档 | 没有 | `FeatureSettingsSection.tsx:198,208-209` |
| `Web Search` ／描述 `Let the model search the web when the selected provider and model support it. Applies to new tasks.` | 允许模型搜网；**注明"只对新建任务生效"** | 待定案项（`AGENTS.md` §2 逐字列了 `web.search` 实现路径未定）；配置里没有 | `FeatureSettingsSection.tsx:214-215` |
| `Feature Tips` ／`Show rotating tips during the thinking phase to help you discover Cline features.` | 思考阶段轮播功能提示 | 没有〔已证：`schema.go` 无此项〕 | `FeatureSettingsSection.tsx:42-48` |
| `Background Edit` ／`Allow edits without stealing editor focus` | 改文件时不抢编辑器焦点 | 无对应物（我们没有"别人的编辑器"这个概念） | `FeatureSettingsSection.tsx:50-56` |
| `Checkpoints` ／`Save progress at key points for easy rollback` | 打快照、可回滚 | 没有〔已证：`schema.go` 搜 `checkpoint`/`rollback` 0 命中〕 | `FeatureSettingsSection.tsx:58-64` |
| `Worktrees` ／`Enables git worktree management for running parallel Cline tasks.` | 打开工作树管理（并行任务） | **有活、没有这枚开关**〔已证：切工作树在面板里还是那句实话 `internal/panel/git.go:82`；工单 186 在办〕 | `FeatureSettingsSection.tsx:66-72` |
| `Hooks` ／`Enable lifecycle and tool hooks during task execution.` | 生命周期与工具钩子 | 没有〔已证：`schema.go` 搜 `hooks` 0 命中；D46 是外部命令插件、不是钩子〕 | `FeatureSettingsSection.tsx:75-81` |
| 分组小标题原文：`Editor`／`Advanced`（大写小标） | 把上面那些分组 | 我们的配置有 section 树但没有"分组小标"的界面（D36 只在 TOML 层） | `FeatureSettingsSection.tsx:224,245` |
| `MCP Display Mode`：`Plain Text`／`Rich Display`／`Markdown` ＋描述 `Controls how MCP responses are displayed` | 工具返回怎么渲染 | **没有、且是 D13 未定案区**（不自行填） | `FeatureSettingsSection.tsx:262-271` |
| `VS Code Terminal`／`Background Exec`（执行模式下拉） | 命令是在**看得见的终端**里跑，还是后台管道跑 | 只有后台进程（`internal/proc`），**没有这一选**〔已证：`schema.go` 无相关键〕 | `apps/vscode/webview-ui/src/components/settings/sections/TerminalSettingsSection.tsx:105-106`（两枚 `VSCodeOption` 的原文）；判据在 `:25-26` |
| `Shell integration timeout (seconds)` ＋校验句 `Please enter a positive number` | 等待 shell 集成就绪的秒数；**输入非法就地报错、不静默纠正** | 没有；我们有超时但没有这一枚键 | `TerminalSettingsSection.tsx:38-39,118` |
| `Enable aggressive terminal reuse` | 复用已有终端而不是每次都新起 | 没有 | `TerminalSettingsSection.tsx:140-144` |
| `Default Terminal Profile` | 从宿主提供的 profile 名册里挑一枚 | 没有 | `TerminalSettingsSection.tsx:146-148` |
| 帮助句原文 `Having terminal issues?` | 出问题时的入口（指向文档） | 没有 | `TerminalSettingsSection.tsx:176` |
| About 页分组原文 `Development`／`Resources` ＋链接 `GitHub`／`Documentation`／`Discord` | 版本与出口 | 有构建信息但没有那一屏（`internal/buildinfo`） | `apps/vscode/webview-ui/src/components/settings/sections/AboutSection.tsx:39-58` |
| ⚠ `Remote Config`（标签原文 `Remotely configured fields`，提示 `Remotely configured fields`） | **一整页只读**：把"被远端下发的字段"列给人看，界面上不许改 | 没有这一页〔已证：`schema.go` 无 remote 概念〕；这是 owner 级"看得见哪些不是我说了算"的一格 | `SettingsView.tsx:80-82` |
| ⚠ `Debug`（提示 `Debug Tools`） | 扩展自带调试工具页 | 我们有 `cmd/balldebug`，但那是给开发者的可执行件、不是给人点的页 | `SettingsView.tsx:97-99` |
| API 那一页的**提供商条目数** | 逐家提供商一张卡＋模型挑选器 | 我们按 D8 抽象、配置里是"提供方＋模型名"两枚文本键（`internal/config/schema.go:210-218` 同款形状） | 名册尺（本机复跑）：`ls apps/vscode/webview-ui/src/components/settings/providers \| grep -c 'Provider.tsx$'` ＝ **28**（口径：只数文件名以 `Provider.tsx` 结尾的组件件；同目录共 42 个文件，其余是 `GenericProviderSettings.tsx` 这类公共件与 `*.test.tsx` ⇒ 「28＝有专属设置卡的提供商数」，**不是**「支持的提供商总数」） |

### 2.2 OpenChamber 扩展形态的设置逐条

**A. VS Code 原生设置页（全部 2 枚，描述逐字抄）**

| 设置项标签原文 | 改了会怎样 | 我方对应物 | 出处 |
|---|---|---|---|
| `openchamber.apiUrl`（类型 `string`，默认空串）＝描述逐字 `URL of an external OpenCode API server. Leave empty to auto-start a local instance.` | 填了就连别人的服务器；留空＝本机自起一台 | **没有这一枚**〔已证：`internal/config/schema.go` 无 `api_url`/远端实例键；"连到哪台"这层在块 A 报告的移动端有、桌面没有〕 | `packages/vscode/package.json:252-256`＋`package.nls.json` 键 `configuration.apiUrl.description` |
| `openchamber.opencodeBinary`（类型 `string`，默认空串）＝描述逐字 `Optional absolute path to the opencode CLI binary. Useful if PATH lookup fails. **Requires window reload or API restart to apply.**` | 找不到 CLI 时手动指路径；**描述里就把"要重启才生效"写给了用户** | 我们靠 `deps.toml` 解析依赖，但**界面上没有任何一处会说"这一枚要重启"**（第四版 §3 已记本格，出处 `docs/reports/missing-features-2026-09-29-v4.md:248`） | `packages/vscode/package.json:257-261`＋`package.nls.json` 键 `configuration.opencodeBinary.description` |

**B. "装进扩展之后会消失的那几枚设置"——这一格是本轮独有的读法**（同一份设置代码、按宿主裁剪，逐枚带行号；均〔已证〕）：

| 设置项（名册名／界面原文） | 在 VS Code 里的处置 | 为什么这样裁（注释原话或代码判据） | 我方对应 | 出处 |
|---|---|---|---|---|
| `theme` | **整块隐藏**：`hasThemeSettings = shouldShow('theme') && !isVSCode` | 颜色跟宿主走，不给第二套 | 我们有 `app.theme`（`internal/config/schema.go:140`，`dark` 默认），**且必须自己实现深浅两套**（没有宿主给我们供色） | `packages/ui/src/components/sections/openchamber/OpenChamberVisualSettings.tsx:698` |
| `appearance` 分组 | 在 VS Code 里**只剩** localization（`theme`/`timeFormat`/`weekStart` 的判定式 `hasAppearanceSettings`） | 同上：外观归宿主 | 同上 | `:701-703` |
| `terminalShell`／`terminalLoginShell` | 隐藏（`&& !isVSCode`） | 终端由宿主提供 | 我们没有宿主终端；`internal/proc` 自己起进程 ⇒ **这一枚对我们不存在也不该造** | `:705` |
| `sessionTabs` | 隐藏（`&& !isVSCode && !isMobile`） | 会话标签页是浏览器/PWA 的概念 | 我们的面板只有一条会话视图 | `:705` |
| `sessionGoal`（＋`sessionAssist` 段） | 隐藏；注释逐字 `The goal loop runs in the web server — VS Code only renders goal state, so the settings section is hidden there too.` | **后端没有那台 loop 就别给开关** | ⚠ 同脾气应学：我们的"写了不管用就响亮拒收"在 `internal/config/unwired.go`（第四版 `missing-features-2026-09-29-v4.md:222` 已记 6 枚），这一条形是"干脆不给开关"⇒ **两家两种做法，见 §6** | `:707,742-744`、注释 `OpenChamberVisualSettings.tsx:1935-1936` |
| `diffLayout` | 隐藏（`shouldShow('diffLayout') && !isVSCode`） | diff 的排布由宿主编辑器决定 | 我们没有宿主编辑器 ⇒ 这一枚对我们**必须自己定**而不是隐藏 | `:737,753` |
| 权限默认档（`PermissionDefaultModeField`） | **在 VS Code 里整枚不画**：`{isVSCode ? null : <PermissionDefaultModeField .../>}` | 档位的真相源在服务端/宿主，扩展不给第二次设置 | ⚠ 直接对上我们的 D31/D43 与 `internal/risk/mode.go:44-56`；**注意：这是"两处可设→只留一处"的先例** | `packages/ui/src/components/sections/openchamber/DefaultsSettings.tsx:365` |
| 会话工作区那一项（`SessionWorkSettings` 里 `!isVSCode` 分支） | 隐藏 | 工作区＝宿主 workspaceFolders | 我们的对应物是 `internal/projctx` | `packages/ui/src/components/sections/openchamber/SessionWorkSettings.tsx:16,36` |
| `subagentReadOnlyBanner`／界面原文 `Allow Prompting Subagent Sessions` | **画得出来**（不在隐藏名册里）⇒ 扩展里也能改 | 子代理会话默认**只读**，要人显式打开才能往里发话 | 〔已证〕我们有子代理名册（`internal/panel/subagent_roster_197.go`）但**没有这一枚开关，也没有默认只读这一说**；⚠ 直接对上"子级永不自批"那条改判，见 §4 | `:718,743,1924-1931`；文案 `packages/ui/src/lib/i18n/messages/en.settings.ts:2156-2157` |

**C. 通知那一页在扩展形态里多出来的一句话**

| 原文 | 位置 | 我方对应 | 出处 |
|---|---|---|---|
| `When enabled, notifications are delivered through VS Code native notifications.` | 通知页"投递"分节，在 VS Code 里替换掉浏览器那句 | 我们的通知只有托盘气泡那一条路、没有"换投递通道"这一说（`cmd/wisp/notify_windows.go:6-13`） | `packages/ui/src/lib/i18n/messages/en.settings.ts:1858`；渲染分支 `packages/ui/src/components/sections/openchamber/NotificationSettings.tsx:471-473` |
| 同页其余原文（**逐条可对照**）：`Notification Delivery`／`Enable Notifications`／`Notify While App is Focused`／`Send test notification`／`Notification permission denied. Enable it in your browser settings.`／`Permission granted, but notifications are disabled.`／`Notification Templates`／`Variables:`／模板事件名 `completion`／`Subagent Completion`／`error`／`question`／字段名 `Title`／`Message` | 通知页 | **全部没有**〔已证：`schema.go` 搜 `notification` 0 命中；"试验证一下通知"这一枚按钮我们也从来没有〕 | `en.settings.ts:1849-1877`；`NotificationSettings.tsx:459-553` |

---

## §3 表三 · 状态与文案表

> 只收**扩展形态里真会出现的句子**（宿主弹的、webview 首屏画的、审批行画的），逐句原文；第四版 §3 已收的 web/CLI 句不重抄（`docs/reports/missing-features-2026-09-29-v4.md:226-252`）。

### 3.A Cline 的"要人处置"那一族句式（同一模板换宾语，共 14 形）

模板：**`Cline wants to <动作> <宾语>:`**（等批准）／**`Cline is <动作>…:`**（已自动放行、正在做）。逐形原文与行号：

| 原文（逐字） | 何时出现 | 我方对应句 | 出处 |
|---|---|---|---|
| `Cline wants to execute this command:` | 要跑命令 | **没有这一句**〔已证：`internal/agent/approval` 里的文案是四条"不可用"实话，不描述要干什么〕 | `apps/vscode/webview-ui/src/components/chat/ChatRow.tsx:319` |
| `Cline wants to use a tool on the <server>`／`access a resource`（同一行三元式） | 要动 MCP | 没有（MCP 属 D13 未定案，不自行填） | `ChatRow.tsx:330` |
| `Cline wants to edit this file:` | 要改已有文件 | 没有 | `ChatRow.tsx:411` |
| `Cline is creating patches to edit this file:` | 同上但已自动放行 | 没有"已放行"这一支说法 | `ChatRow.tsx:410` |
| `Cline wants to delete this file:` | 要删文件 | 没有 | `ChatRow.tsx:445` |
| `Cline wants to create a new file:` | 要新建 | 没有 | `ChatRow.tsx:463` |
| `Cline wants to read this file:` | 要读 | 没有 | `ChatRow.tsx:486` |
| `Cline wants to view the top level files in this directory:` | 要列目录 | 没有 | `ChatRow.tsx:527` |
| `Cline wants to recursively view all files in this directory:` | 要递归列 | 没有 | `ChatRow.tsx:549` |
| `Cline wants to view source code definition names used in this directory:` | 要看符号名 | 没有 | `ChatRow.tsx:571` |
| `Cline wants to search this directory for <regex>:` | 要搜 | 没有 | `ChatRow.tsx:591` |
| `Cline wants to fetch content from this URL:`／`Cline wants to search the web for:` | 要抓网／搜网 | 没有（`web.search` 未定案） | `ChatRow.tsx:652,681` |
| `Cline wants to start a new task:` | 要开子任务 | 没有；我们有 `internal/panel/subagent_roster_197.go`（名册）但没有那句话 | `ChatRow.tsx:1125` |
| `Cline wants to use a subagent:`／`Cline wants to use subagents:`（单复数两句） | 要起子代理 | 没有 | `apps/vscode/webview-ui/src/components/chat/SubagentStatusRow.tsx:200` |
| `Cline wants to use the browser:` ↔ `Cline is using the browser:` | **同一枚动作的两时态对照**（要不要问 = 有没有自动放行） | 没有；**这一形是本块最值得记的形状**：文案随"权限判定结果"改时态，而不是加一句"已自动允许" | `apps/vscode/webview-ui/src/components/chat/BrowserSessionRow.tsx:362` |
| `Cline is condensing the conversation:` | 正在压上下文 | **我方名册里没有"正在压缩"这一态**（第四版 `missing-features-2026-09-29-v4.md:238` 已记；本轮复量名册＝20 枚，见 3.D） | `ChatRow.tsx:608` |
| `Cline is having trouble...` | 连续出错的兜底 | 没有 | `ChatRow.tsx:314` |

### 3.B Cline 的按钮与状态动词（一屏之内全览；第四版只点了三枚批准钮）

| 原文 | 用在哪 | 我方对应 | 出处（同一文件＝`apps/vscode/webview-ui/src/components/chat/chat-view/shared/buttonConfig.ts`） |
|---|---|---|---|
| `Approve`／`Reject` | 通用工具批准 | 面板只能拒（`docs/specs/SPEC-08*.md:156-176` C17 名册） | `:49-50`（`tool_approve`） |
| `Save`／`Reject` | **改/建/删文件时按钮换成 `Save`** | 没有这一区分 | `:57-58`（`tool_save`），判据 `:251-261`（只有 `editedExistingFile`/`newFileCreated`/`fileDeleted` 三种才走 Save） |
| `Run Command`／`Reject` | 跑命令 | 没有 | `:67-68` |
| `Proceed While Running` | 命令还在跑、允许它先往下走 | 没有 | `:75-76`、`:190` |
| `Retry`／`Start New Task` | 请求失败 | 没有 | `:31-32` |
| `Proceed Anyways`／`Start New Task` | 连错达上限、仍要继续 | 没有 | `:39-40` |
| `Resume Task` | 回来接上 | 没有 | `:135` |
| `Start New Task with Context` | 带着当前上下文开新任务 | 没有 | `:151` |
| `Condense Conversation` | 手动压上下文 | 没有 | `:161` |
| `Report GitHub issue` | 一键报障 | 没有 | `:169` |
| `Cancel` | 流式中间态唯一的按钮 | 有对应物：四条取消出口各说实话（`internal/agent/approval/approval.go:106-114`） | `:180,191,209` |
| `Thinking`（标题）／`Thinking...`（骨架行） | 想的时候那一行 | 我们的球有"在想"这一态（`internal/ball/statevisual.go`），但没有这句话 | `apps/vscode/webview-ui/src/components/chat/ThinkingRow.tsx:24`、`RequestStartRow.tsx:217,238` |
| `Completed`（小标大写） | 完成块的头 | 没有这句 | `apps/vscode/webview-ui/src/components/chat/CompletionOutputRow.tsx:73` |
| `Subagent status update unavailable.` | 子代理状态拿不到时**明说拿不到** | 同脾气先例只有 `internal/panel/git.go:82`；这一格没有 | `apps/vscode/webview-ui/src/components/chat/SubagentStatusRow.tsx:186` |
| `The model used search patterns that don't match anything in the file. Retrying...` | 搜空了：说清**为什么**在重试 | 没有 | `apps/vscode/webview-ui/src/components/chat/ErrorRow.tsx:179` |
| `Daily free model limit reached`／`Total Spent: `／`Request ID: ` | 额度与排障标识 | 有预算与提醒阈值（`internal/config/schema.go:544-548`）**但没有这三句** | `apps/vscode/webview-ui/src/components/chat/ClineFreeModelLimitError.tsx:128`、`CreditLimitError.tsx:63,65`、`ErrorRow.tsx:108` |

### 3.C Cline 的后端相位名（7 相）——与我方 20 态对照，不是同一层东西

原文（`apps/vscode/webview-ui/src/components/chat/chat-view/shared/buttonConfig.ts:368-397` 的 `switch (turnState.phase)`）：
`idle`／`streaming`／`completed`／`resumable`／`error`／`awaiting_followup`／`awaiting_approval`。〔已证〕

三点必须说清的差别（免得被读成"他们只有 7 态、我们 20 态更细"）：
1. 这 7 相**只管一件事：底部该画哪两枚按钮**（注释逐字 `Authoritative button configuration derived from the backend-owned TurnState`，`:356-362`）；它不是生命周期全貌。
2. 我方 20 态是**球与进程的生命周期**（`internal/statemachine/states.go:11-30`，尺：`grep -c 'State = "' internal/statemachine/states.go` ＝ **20**，本机复跑），包含声音/网络/下载/看门狗这些他们那层根本没有的东西。
3. **可对齐的一格**：他们的 `awaiting_approval` 有配套 UI（按钮名随工具种类变），我们的 `StateAwaitingApproval` 只有状态、没有那两枚动词出口（`docs/specs/SPEC-08*.md:156-176`）。⚠ 这一格第四版已记（`missing-features-2026-09-29-v4.md:85`），此处只补"他们的相位还带 `seq`＋`anchorTs` 两枚数用来防按钮串到下一条请求"这一形（`buttonConfig.ts:61`）。

### 3.D OpenChamber 扩展宿主侧的全部人话（`l10n/bundle.l10n.json` 共 **28** 键，逐句；尺：本机 `python -c "import json;print(len(json.load(open('packages/vscode/l10n/bundle.l10n.json'))))"` ＝ 28）

| 原文 | 何时出现 | 我方对应 | 出处 |
|---|---|---|---|
| `OpenChamber: Chat sidebar is not ready` | 侧栏 webview 还没起来就点按钮 | 没有这一句（我们的面板要么开得出、要么根本没入口） | `packages/vscode/l10n/bundle.l10n.json`；调用 `packages/vscode/src/extension.ts:174,191` |
| `OpenChamber: No active session` | 要点"把当前会话搬进标签页"但没有当前会话 | 没有 | 同上；`src/extension.ts:256` |
| `OpenChamber: No file selected to mention` | 资源管理器里没选中文件 | 没有 | `bundle.l10n.json`；`src/extension.ts:412` 一族 |
| `OpenChamber: Some selected entries were skipped (folders or unsupported resources)` | **多选里有不能用的：说清跳了、并说清为什么** | 我们的受理器只有一条 `path is not a file` 形状的说法？——本轮未读到，写"没有"〔已证：`internal/panel/attachments.go` 无该句〕 | `src/extension.ts:412` |
| `OpenChamber: API connection restarted` | 手动重连连上了 | 没有 | `src/extension.ts:299` |
| `OpenChamber: Failed to open sidebar - {0}`／`Failed to restart API - {0}` | 失败带**原始原因**插值 | 同脾气：我们的实话句会带工单号（`internal/panel/git.go:82`） | `src/extension.ts:168,301` |
| `OpenChamber: No folder is open. Open a folder to start a new session.` | 没开文件夹时**不猜**上一个目录 | 没有 | `src/SessionEditorPanelProvider.ts:125-137` |
| `OpenChamber [Add to Context]: No active editor`／`No text selected` | 每种点不着各一句 | 没有 | `src/extension.ts:310` 一族＋`bundle.l10n.json` |
| `OpenChamber [Add Comment]: No active editor`／`File is outside the workspace`／**`The comment never reached the chat and was discarded`** | 评论没送进去时**明说已丢弃** | 没有；这形对我们最有价值（面板/球那条路一旦没回执，现在无处可说） | `bundle.l10n.json`；`src/InlineCommentThreads.ts` 一族 |
| `Comment on line {0}`／`Comment on lines {0}-{1}`／`Not sent yet` | 线程标题与"还没发"状态 | 没有 | `bundle.l10n.json` |
| `Explain the following Code / Text:`／`Improve the following Code:` | 塞进会话的那句话的开头 | 没有 | `bundle.l10n.json` |
| `OpenCode CLI not found. Install it and ensure it's in PATH.`／`...Please install it and ensure it's in PATH.`（**两版，一字差**） | 找不到被管 CLI | 我们有 `deps.toml` 解析与拒启路径，但没有这两句 | `bundle.l10n.json`（两条并存＝可核对的漂移） |
| `Failed to start OpenCode: {0}`／`More Info` | 起进程失败＋一枚"看详情"按钮 | 没有 | `bundle.l10n.json` |
| `New Session`／`Session`／`Run on several models` | 编辑器标签页的标题（含兜底名） | 没有 | `bundle.l10n.json` |
| `OpenChamber` | 通知与评论作者名（兜底标题） | 我们叫 Wisp（D28） | `bundle.l10n.json` |
| ⚠ 名册里没有中文 | **该 bundle 的英文键自映射**（`'x' => 'x'`），另外只有 `bundle.l10n.fr.json`／`bundle.l10n.tr.json` 两份译文〔已证：`ls packages/vscode/l10n` ＝ 3 个文件〕 | 我方 `app.language` 默认 `zh-CN`（`internal/config/schema.go:138`）⇒ **宿主层文案我们自己必须写中文**，没有现成机制可借 | `packages/vscode/l10n/`（本机 `ls` 读数 3） |

### 3.E OpenChamber 扩展首屏与通知的默认句

| 原文 | 何时出现 | 我方对应 | 出处 |
|---|---|---|---|
| `Starting OpenCode API…`／`Initializing…`／`Connecting…`／`Connected!`／`Connection error`／`Reconnecting…`／`OpenCode CLI not found. Please install it first.` | 宿主 HTML 里的首屏名册（7 句 × en/fr/tr） | 没有首屏文案表 | `packages/vscode/src/webviewHtml.ts:255-299` |
| 上面那张的**权威版**多出 4 枚：`disconnected`／`startingDevServer`／`waitingDevServer`／`loadingData`（后三枚带插值参数） | ui 包那份才是运行时用的 | 没有 | `packages/ui/src/lib/i18n/bootstrap.ts:232-290` |
| `{agent_name} is ready` ＋ `{model_name} completed the task` | 会话完成通知默认句 | 我们的托盘气泡没有这句式（`cmd/wisp/notify_windows.go:6-13`） | `packages/vscode/webview/main.tsx:2009-2014`；同名默认句也在设置页显示：`packages/ui/src/lib/i18n/messages/en.settings.ts:1876-1877` |
| `Tool error` ＋ `{last_message}`／`Input needed` ＋ `{last_message}` | 出错与"要你答"通知 | 没有 | `main.tsx:2027-2033,2044-2050` |
| 兜底句 `An error occurred`／`Agent is waiting for your response`／`Agent is ready` | 模板解析不出来时**仍有整句**，不留空 | 没有 | `main.tsx:2010,2032,2049` |
| `Only show text when something is wrong — progress states stay silent (the animated logo already signals "working")` | 首屏设计说明（注释，不是给用户看的） | 同思路我们已由球承担 | `packages/vscode/src/webviewHtml.ts:271-272` |

### 3.F 我方独有、别家没有的句子（本轮复量仍在产码里）

| 我方原句 | 出处 | 本轮对标读数 |
|---|---|---|
| 四条出口不可用时各说实话：`语音取消不可用`／`面板取消不可用（票 37 未接入）`／`悬浮球取消不可用`／`Esc 取消不可用` | `internal/agent/approval/approval.go:106-114` | 两家的审批行**没有任何一条"这条出口不通"的说明**——它们只做能做的出口。⇒ 我们更诚实，但也更啰嗦；形状级取舍见 §7 |
| `切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主）` | `internal/panel/git.go:82` | OpenChamber 的同格是一枚**看得见的按钮**＋一句"为什么点不着"（3.D 那 4 条守卫句），不是文本自述 |
| `DEFERRED(playback)` 这类代码标记 | `internal/audio/doc.go:27` | 不在界面层，仅登记 |

---

## §4 必答一：点进某个子任务/子代理，能不能看到它各自的流式工作页面？

### 结论一句话

- **Cline：不能。**〔已证〕扩展形态里子代理只有**一枚"状态行"**，没有它自己的流式页面可点进去；能看到的"活"只有一枚工具名的文本与三个计数。
- **OpenChamber：能，但装进 VS Code 之后形态变了——点它是"整屏换过去"，不是"旁边再开一页"。**〔已证〕

### 4.1 Cline：状态行是唯一可见物（逐项拆）

| 能看到什么 | 原文／形状 | 是不是流式 | 出处（`apps/vscode/webview-ui/src/components/chat/SubagentStatusRow.tsx`） |
|---|---|---|---|
| 标题 | `Cline wants to use a subagent:`／复数 `Cline wants to use subagents:` | 否 | `:199-200` |
| 每条子任务一行 | 状态图标四态：转圈（`running`）／勾（`completed`）／叉（`failed`）／斜线（`cancelled`） | 图标随推送换，**内容是整条替换、不是逐字流** | `:41-54`、`:222-240` |
| 提示词 | 两行截断＋`Show more`（aria `Show full subagent prompt`） | 建提示词那一条在流式时**逐字长出来**（仅提示词本身，非工作过程） | `:154-176`、`:229-230` |
| 一行统计 | 逐字模板 `${toolCalls} tools called · ${contextTokens} tokens · ${cost}` | 否，计数刷新 | `:232` |
| **唯一"正在做什么"的线索** | `latestToolCall` 一枚字符串，**等没有详情可展开时才显示**，单行截断、等宽字体 | 半流式：只是"最新那一次工具调用"的名字 | `:233,268-270`；字段定义 `apps/vscode/src/shared/ExtensionMessage.ts:308-322` |
| 结果 | `Show output`／`Hide output`（aria `Show subagent output`）→ 展开是**最终** `result`（Markdown）或 `error` | 否：只在 `status === "completed"`／`"failed"` 时才有 | `:225-227,254-278` |
| 没数据时 | `Subagent status update unavailable.` | — | `:186` |

**关键负向证据**〔已证〕：`SubagentStatusItem` 的字段就那 12 枚（`index/prompt/status/toolCalls/inputTokens/outputTokens/totalCost/contextTokens/contextWindow/contextUsagePercentage/latestToolCall/result/error`，`ExtensionMessage.ts:308-322`）——**没有任何"子会话 id／子 transcript 入口／子任务页路由"字段**，所以点不进去不是 UI 没做，是**数据结构里没有那条路**。批准那一侧也只有 `use_subagents` 一枚 ask（按钮 `Approve`/`Reject`，`buttonConfig.ts:98-105`），**子代理内部再要批准时，扩展里没有第二个审批面**。

### 4.2 OpenChamber：扩展里点它是"换屏"，且进去默认只读

| 环节 | 形状 | 是不是流式 | 出处 |
|---|---|---|---|
| 父会话侧的子代理面板 | 折叠段 `Subagents`，摘要位是 `busy/total`（如 `2/3`），无子代理时**整段不画**（`if (children.length === 0) return null`） | 计数随事件刷新 | `packages/ui/src/components/chat/work-status/WorkStatusSubagentsSection.tsx:80-90`；标题原文 `packages/ui/src/lib/i18n/messages/en.ts:3263` |
| 每条子任务一行 | 名字（无题时兜底原文 `Subagent`）＋四值之一：**`needs permission`／`asked a question`／`is working`／`Done`**（优先级就按这个顺序）＋（有则）该子树成本 | 值随推送换 | `:93-121`；四句原文 `en.ts:3255-3269`；成本聚合 `:36-40`（注释逐字：嵌套子代理的花费要"滚"到直接子行下，不然会消失） |
| 点一行的行为（**扩展形态专属分支**） | 判据逐字 `if (isEmbeddedSessionChat() || isMobile \|\| isVSCodeRuntime()) { setCurrentSession(childId, directory); return; }` ⇒ **在 VS Code 里直接把整屏切到那条子会话**；只有浏览器/PWA 才在旁边开一枚只读标签页 | — | `:62-76`；aria 原文 `Open {name}`（`en.ts:3282`） |
| 进去之后 | 输入区被一块横条替掉，原文 **`Subagent sessions cannot be prompted.`**；除非设置页把 `Allow Prompting Subagent Sessions` 打开 | 子会话自身的流**看得见**（就是普通会话视图） | 横条组件 `packages/ui/src/components/chat/ChatContainer.tsx:601-614`；原文 `en.ts:2296`；开关 `OpenChamberVisualSettings.tsx:1924-1931` |
| ⚠ 最响的一句（设计账，逐字抄） | 文件头注释：**"Running subagents and, more importantly, their blockers: a permission request raised by a child session has no representation in the transcript, so this panel is the only place it becomes visible."** | — | `WorkStatusSubagentsSection.tsx:20-24` |
| 面板在扩展里挂不挂 | 挂载判据只有 `workStatusPanelMountable = !isMobile` ⇒ **VS Code 里挂得出来** | — | `ChatContainer.tsx:959-962` |
| 子代理发起的批准会不会被认出来源 | **会**：批准卡上挂一枚来源徽标，原文 `From subagent` | — | `packages/ui/src/components/chat/PermissionDock.tsx:102`；原文 `en.ts:2331` |

### 4.3 我方对照（只引 `docs/**` 与本腿现读 Go）

| 这一格 | 我方现状 | 出处 |
|---|---|---|
| "子任务行"这件事 | 有**名册**、没有那一行：`internal/panel/subagent_roster_197.go`（名册件）与 `internal/panel/subagent_stream_197_test.go`（流件）都在，界面上那一行属于〔界面地界·不读〕 | 本机 `ls internal/panel` 读数 |
| "子代理的批准在父任务里看得见" | **两家的形状互不相同**：Cline 数据结构里根本没有那条路〔已证 4.1〕；OpenChamber 明确写"这是唯一能看见的地方"〔已证 4.2〕 ⇒ 我们撞上的是票 197 那一问（台账在册），**本轮证据说明这一格不是可选加固，是"不建就看不见"** | `docs/reports/missing-features-2026-09-29-v4.md` 的块 3；`AGENTS.md` §2 未定案清单不含此格 |
| "点进去会不会把父任务丢了" | OpenChamber 在扩展里就是"整屏换过去"（`setCurrentSession`），**没有任何"回来"的按钮**；这与我们"面板是低频、可关"的前提冲突 ⇒ 见 §7 | 上表 `:62-76` |

---

## §5 必答二：审批卡有没有"本次／本次会话内／长期"三档＋拒绝/允许的理由输入框？

### 5.1 逐家答"三档"

| 家 | 本次 | 本次会话内 | 长期 | 出处与判据 |
|---|---|---|---|---|
| **Cline（扩展形态）** | **有**：底部两枚按钮，随动作换名（`Approve`／`Save`／`Run Command` 各配 `Reject`） | **没有这一档**〔已证：`BUTTON_CONFIGS` 全名册里没有任何"这一会话内"条目〕 | **有**，但家不在卡上：输入框上方那一横条 `Auto-approve:` ＋ 展开的**能力级**开关（`Read files`／`Edit files`／`Execute commands`／`Fetch web content`／`Use MCP servers`，原文标签） | 卡上按钮 `apps/vscode/webview-ui/src/components/chat/chat-view/shared/buttonConfig.ts:46-105`；条与档名册 `apps/vscode/webview-ui/src/components/chat/auto-approve-menu/constants.ts:3-35`、`AutoApproveBar.tsx:109`；说明句逐字 `Let Cline take these actions without asking for approval.` `AutoApproveModal.tsx:71` |
| ⚠ Cline 的"长期"**默认就是开着的** | — | — | 出厂默认读数：`readFiles: true`／`editFiles: true`／`useBrowser: true`／`useMcp: true`／`executeAllCommands: true`，只有 `executeSafeCommands: false` | `apps/vscode/src/shared/AutoApprovalSettings.ts:28-44`（⚠ 与我方"默认最严"正好相反：`internal/risk/mode.go:44-56`、`internal/config/schema.go:470` 默认 `ask_every_step`） |
| **OpenChamber（扩展形态）** | **有**：卡上三枚按钮里最左那枚，原文 `Allow once`（快捷键 `alt+enter`） | **有，但不在卡上**〔已证：它是设置里一枚**每会话布尔位**，持久在宿主 `globalState['permissionAutoAccept']`，带 `revision` 单调号、跨面板广播、写操作串行化〕 | **有**：中间那枚按钮原文 `Always allow`；**能推断出规则时按钮名直接写成 `Always: {patterns}`**（前 2 条＋省略号，全量挂 `title` 悬停） | 三枚按钮 `packages/ui/src/components/chat/PermissionCard.tsx:408-443`（`once`／`always`／`reject`＋三枚 `<kbd>`）；原文 `packages/ui/src/lib/i18n/messages/en.ts:3377-3380`（`Allow once`／`Always: {patterns}`／`Always allow`／`Deny`）；会话档 `packages/vscode/src/bridge-permission-auto-accept-runtime.ts:1-32`（`sessions: Record<string, boolean>`＋`revision`）、`:54-71`（写入并广播）、`:72-103`（`api:permission-auto-accept:get\|set` 两条消息） |
| ⚠ OpenChamber 在扩展里**不给改默认档** | — | — | 判据逐字 `{isVSCode ? null : <PermissionDefaultModeField agentName={defaultAgent} />}`——权限默认档那一枚设置在 VS Code 里**整块不画** | `packages/ui/src/components/sections/openchamber/DefaultsSettings.tsx:365` |

**两家的"会话内"这一档放的地方不一样**（这是本轮新读出来的形状差，第四版没写过）：Cline 只有"这一次"和"这类永远"，中间那一档**是空的**；OpenChamber 把中间那一档**搬出卡片**、做成宿主里持久、跨端广播的一枚开关 ⇒ **卡上永远是三枚、会话档在旁边**。

### 5.2 逐家答"理由输入框"

| 家 | 有没有专用理由框 | 实际怎么带话 | 出处 |
|---|---|---|---|
| Cline | **没有**〔已证：`ActionButtons` 只画两枚 `VSCodeButton`，没有任何输入件〕 | **靠输入框捎带**：批准/拒绝那一刻把草稿整体取走，有内容就连 `text/images/files` 一起发，`responseType` 分别是 `yesButtonClicked`／`noButtonClicked`，发完 `consumeDraftSnapshot` 清草稿；无内容则**只发 responseType** | `apps/vscode/webview-ui/src/components/chat/chat-view/components/layout/ActionButtons.tsx:10-19,134-155`；发送侧 `chat-view/hooks/useMessageHandlers.ts:490-503` |
| Cline（问答卡那一族） | 也**没有**专用框，但**把输入框当成"理由/补充"字段用**：点选项时拼 `option + ": " + inputValue` | 逐字：`text: option + (inputValue ? \`: ${inputValue?.trim()}\` : "")` | `apps/vscode/webview-ui/src/components/chat/OptionsButtons.tsx:79-85` |
| OpenChamber | **没有**〔已证：`PermissionCard.tsx` 里无 `<textarea>`/`<input>`，`onRespond` 的入参类型是 `PermissionReply`（即 `once`/`always`/`reject` 三个码），**函数签名里就没有能放话的位置**〕 | 不带你为什么说"不" | `packages/ui/src/components/chat/PermissionCard.tsx:366-368`（`onRespond: (response: PermissionReply) => void`）、`:376-390`、`:408-443` |
| 两家的"要人答"另有一路 | OpenChamber 把"提问"独立成一族：卡上挂 `From subagent` 徽标、多问题时显示 `{current} of {total}` 步进、并有一句注释逐字 `Forms replaced questions in OpenCode 2.x: one titled request the user answers.` | 这是"答"，不是"批"，**两家都分开了这两件事** | `packages/ui/src/components/chat/PermissionDock.tsx:102,107`；原文 `en.ts:3354`（`{current} of {total}`）、`en.ts:2331`；注释 `packages/vscode/webview/main.tsx:2039` |
| 我方 | 〔已证〕四条应答出口不可用时各说实话，但**"允许"这一侧被结构禁掉**、也没有理由框：`docs/specs/SPEC-08*.md:156-176`；四条实话句在 `internal/agent/approval/approval.go:106-114` | 我方对应：第四版已记这一格碰 `Q-49`（要 owner 一句话），本腿不自行填 | `docs/reports/missing-features-2026-09-29-v4.md:85` |

### 5.3 本轮额外钉出来的三枚形状（都带行号，供 owner 取舍）

1. **"按了但后端没接住"要有专句**：OpenChamber 在评论这条路上写了 `The comment never reached the chat and was discarded`（`packages/vscode/l10n/bundle.l10n.json`）⇒ 同一族"我点了、其实没送到"的诚实句在他们那里是**成体系的**（见 §3.D）。
2. **批准按钮上直接印"这次点下去会存成什么规则"**：`Always: {patterns}` ＋ `title` 里放全量（`PermissionCard.tsx:349-358`、`en.ts:3379`）；⚠ 第四版 `missing-features-2026-09-29-v4.md:86` 已记同一格（当时行号取自 2.0.3 快照，本腿在 2.0.4 复现＝同一函数、行号未漂）。
3. **安全网扣住的那条要显示**：新增原文键 `heldBySafetyNet`（渲染点 `PermissionCard.tsx:304`）⇒ 他们把"这条被兜底规则扣住了"做成卡片上的一句话，而不是让它默默等超时。

---

## §6 我写错的条目（拿第四版清单逐条对抗本腿自己）

（填写中）

---

## §7 不建议抄（只给形状级理由）

（填写中）

---

## §8 没调研完／留给编排者

（填写中）
