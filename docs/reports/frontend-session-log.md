# 前端会话工作记录（`frontend-session`）

> 这枚文件归**前端专责会话**所有。编排者的地界是 `docs/reports/pending-and-issues.md`（真相源台账）
> 与 `docs/reports/HANDOVER.md`（停车点），本文件**不取代、不写入**那两份。
> 开机简报：`docs/reports/frontend-session-brief.md`（2026-09-24 23:45）· 交接页：`docs/reports/frontend-handoff.md`
> 锚点：`dev` @ `f1c9471cecaad7ed24b9df14730b4ab972bab4d4`，本机时刻 2026-09-24 23:48 +0800。

---

## 1. 交付物：`design/doubao/demo` → `frontend/` 十屏差距表

### 1.0 这张表怎么读

- **已有**＝真树里能看见这块界面，且有非实现者能复算的证据。
- **缺**＝真树里没有。`缺(死码)` 表示文件在、但全树无人 import、且装的是上游假数据，**不算已有**。
- **冲突**＝照 demo 做会撞上冻结契约或本仓门禁，条目里点名撞哪一条。
- 判定尺是**盘上现量**，不是 demo 的 `README.md` 自述（该自述有四处与盘上不符，见 §1.3）。

### 1.1 先记三条跨屏结构缺口（它们决定十屏里九屏的排程）

| # | 缺口 | 现量证据 | 后果 |
|---|---|---|---|
| **X1** | **真树没有路由，也没有"屏"的概念** | `frontend/src/main.tsx:6-10` 只 `render(<App/>)` 不传 props；`App.tsx:41-47` 的 `snapshot`/`tab` 全走默认值 ⇒ 线上唯一画面是 `App.tsx:34-39` 的**硬编码空态**。`panel-skeleton.tsx:16` 的 `PANEL_TABS` 只有 4 个值（命令/结果/历史/配置），`:37-48` 渲染的是**无 onClick 的 `<span>`** | demo 是 10 屏＋左导航 10 项（`app.js:18-26` 只列 9 项，`firstrun` 另经 `app.js:229` 注册）。要做九屏，**第一步不是画屏而是先立路由/视图切换这一层**，否则每屏各造一套显隐 |
| **X2** | **入向数据通路两端都没接线** | `internal/panel/composer_handlers.go:36-41` 原文："Nothing calls this handler yet, because the WebView2 `event -> ParseComposerRequest` hop does not exist"；全 Go 树 grep `PostWebMessage|ExecuteScript|AddWebResourceRequested` 只命中文档注释（`internal/panel/doc.go:3`）；`cmd/wisp/run.go:224-225` 同口径 | 前端能画的只有**显示＋发起请求**五枚 method（`src/lib/panel.ts:179/237/246/263/291`）。任何"从 Go 真取到数据"的屏（设置/隐私/成本/安全）在票 114 的 Go 半边落地前**只能画到 fixture 级**，不能报"完成" |
| **X3** | **token 生成器的输入源已被 owner 挪出原位** | `frontend/scripts/gen-tokens.mjs:24` 读 `design/assets/tokens.css`；该路径**盘上已不存在**（`git status` 显示 ` D design/assets/tokens.css`，HEAD 里仍有 12942 字节）。owner 把它**移动**成了未跟踪的 `design/old/assets/tokens.css`，两枚文件 md5 相同：`141c570f063606cf663629c717098f40` | 本机 `npm run tokens:check` 与 `npm run tokens` **必坏**（文件找不到）。CI 不受影响（干净检出时该文件在库里）。⇒ 待拍板项 **P1** |

### 1.2 十屏逐判

| 屏 | demo 文件 | 判 | 依据与缺口 |
|---|---|---|---|
| **对话 chat** | `screens/chat.js` | **部分** | 已有：结果流 `components/result-stream.tsx:10-25`、输入行 `components/composer.tsx:58-64`（档位下拉/工作区/textarea/附件/发送）。缺：**用户消息气泡与多轮线程**——`PanelSnapshot`（`src/lib/panel.ts:129-137`）只有 `pending/results/composer`，**没有任何"用户说过什么"的字段**，所以历史分组（`chat.js:41`）、会话标题下拉、重播、停止按钮、推理折叠全部无处安放。缺：附件缩略图（demo 有 `1240×860 · 384 KB`＋移除键，真树只有文本行）。缺：干活/陪聊/自动三态切换与陪聊红环（`chat.js:17-21`，真树 `composer.tsx:41-46` 只有三档**权限**名，不是**模式**名） |
| **审批 approval** | `screens/approval.js` | **已有**（十屏里唯一完整的一格） | `components/l2-approval-card.tsx` 全 11 字段真渲染（level / args 逐行 / reason / rulesHit / callChain / C25 污染 / 判定来源），props 走真通路 `snapshot.pending[]`（`App.tsx:54-56`）。demo 的"面板仅拒绝＋查看、无允许键"与真树一致。**但**：生产入向未接（X2），线上恒为 0 张卡；且票 77 `AC#3` 已被独立复算判 **PARTIAL**（`docs/evidence/s1/77-ac3-l2-card-real-risk-fields-r1.md:200-212`）——PARTIAL 的理由正是 X2，不是卡片本身 |
| **命令 palette** | `screens/palette.js` | **缺** | 全树零实现。只有两枚痕迹：token `--palette-w: 680px`（`src/styles/tokens.generated.css:139`）与 tab 文字"命令"（`panel-skeleton.tsx:16`，死 span）。demo 的 5 组结构、搜索过滤＋空态、逐项 `kbd` 快捷键、上下键导航／Enter 执行／Esc 关闭、Ctrl+K 唤起（`app.js:329-330`）全部要从零做。**冲突（弱）**：demo 里 `fs.delete` 标"高危·需确认"，真树侧 `--danger` 色已有，无契约冲突 |
| **任务 tasks** | `screens/tasks.js` | **缺(死码)** | `src/components/ai-native/task-rows.tsx` 在，但**无人 import**、装上游假数据＋自播放定时器，文件头 `:14` 自述 "NOT mounted: demo data intact"；`VENDORED.md:55-59` 记有一条机器守卫 `TestVendoredDemoComponentsAreNotMounted` 故意维持"未挂载"。demo 的六态色标、路径锁文案（"在等 task-6 释放 Desktop"）、cron 中文化＋电源条件、新建任务 Modal 全无对应实现 |
| **球状态 ball** | `screens/ball.js` | **缺 ＋ 冲突** | 真树无球渲染（`frontend/src` grep 球相关 0 命中），只有颜色/尺寸 token（`tokens.generated.css:63-71`、`:140-142`）。**冲突①（契约）**：简报 §3 与 `VENDORED.md:87` 明写桌面球永远 Win32/Direct2D，React Bits 的 `Agentic Ball` 只能当**面板内次级指示器** ⇒ demo 那颗 CSS 预览球（`ball.js:171` 的 `orbPreview`）可以照搬进面板，**不可以**被读成"球进 Web 侧"。**冲突②（数字）**：demo 标题写"20 态 · **42** 条权威转移"（`ball.js:13`），冻结表 `PLAN.md` D43 是 **20 态 · 40 条**（`PLAN.md:3051` 正文，转移行 `:3055-3094` 实数 40）⇒ 42 无出处，照搬即造出第二张不存在的权威表 |
| **设置 config** | `screens/config.js` | **缺 ＋ 冲突** | 真树只有一个"配置"tab 名（死 span），无设置视图、无设置项 struct。**冲突①**：demo 的生效徽标画了 **7 档**（`config.js:19-27`：即时/重载/重启/重启生效/立即生效/下次对话生效/仅新任务生效），冻结契约 D36 是**三档** ⇒ 要压回三档，且这属于"改 demo"不是"改契约"，需 owner 认。**冲突②**：`config.js:345` 注释写"配置项**存 localStorage**，WebView2 启动时读取" ⇒ 直接撞票 77 `AC#5`（前端零本地持久化），守卫在 `internal/panel/frontend_hygiene_test.go:48-56`。**注意**：这条只是注释、不是真代码，但**照搬注释也会被读成口径**，进树前要改写成"config.toml 为唯一真相源"（demo 自己在 `config.js:338` 已这么写，两行互斥） |
| **安全 security** | `screens/security.js` | **缺 ＋ 冲突（最硬的一格）** | 真树 0 命中。**冲突①（R20／ban #6 精神）**：`security.js:247` 在"生效授权"表的操作列渲染了 **`允许一次`** 按钮（与 `撤销` 并排）。契约 R20 与 ban #6 要求面板侧**不得**产生授权决定；demo 的审批屏（`approval.js:227`）反而**明确写了**"面板不提供允许按钮（面板来源的允许在服务端被结构拒绝）"——**同一份 demo 里两屏互斥**，安全屏那一支是错的，不许照搬。⚠ 而且**门禁抓不到它**：ban #6 的 `approval.decide` 正则只扫 `.go` 文件（`tools/d22scan/main.go:348-352`），对 `frontend/` 的中文"允许一次"零感知 ⇒ 这一格只能靠人核，不能靠 CI。**冲突②**：`security.js:147-151` 的 16 枚工具名（`clipboard.read`、`clipboard.write`、`mail.send`…）需逐枚对 D34 内置工具权威表（`PLAN.md:2514` 起）与 SPEC-06 §7 的 RESERVED 名单，**本轮我未逐枚核完**，见 §1.4。**非冲突（已排除）**：`security.js:118` 写"进程隔离：Job Object（已生效）"，我一度以为与裁定 R21 矛盾，实测 `internal/proc/jobscope_windows.go` 存在、`JobObject` 9 处命中 ⇒ 代码层面确有实现，demo 的措辞站得住（是否真挂在生产启动路径上属票 114 Go 半边，不归我判） |
| **隐私 privacy** | `screens/privacy.js` | **缺 ＋ 冲突** | 真树 0 命中。demo 的 5 tab（L1 画像/L2 记忆/L3 日志/tool_call 取证/artifacts）、保留期 30/90/永久、麦克风开关、存储位置（`wisp.db 12.4MB`、`artifacts 86/500MB`）全无实现。**冲突**：同 config 的 localStorage 一条，且这里**不只是注释**——门本身在 `app.js:40` 与 `:242` 真读写 `localStorage`（键 `wisp-firstrun-done`），照搬即被 `frontend_hygiene_test.go:49` 判红。另 `privacy.js:278` 的"30 天保留期…存 localStorage"同样要改写 |
| **成本 cost** | `screens/cost.js` | **缺** | 真树 0 命中（`cost` 只出现在英文注释如 `panel.ts:231`）。demo 的今日/本月/累计三卡、超预算红态、80% 告警线、纯 CSS 柱状图、日/月视图切换、token 分解（IN/OUT/CACHED）、Top5 任务榜、C23 微单位对账说明全部从零。**依赖**：所有数字都要来自 C23 真字段，X2 未解前只能画到 fixture 级 |
| **首次引导 firstrun** | `screens/firstrun.js` | **缺 ＋ 冲突** | 真树无 onboarding。**冲突①**：demo 用 `localStorage` 当门（`app.js:40/:242`），被 AC#5 守卫禁 ⇒ 真树里这个门必须由 Go 侧快照的"是否已完成首次配置"字段驱动。**冲突②（未定案项）**：`firstrun.js:12` 与 `:104-108` 把 **DPAPI** 当定案画（"Windows 凭据 · DPAPI 加密"），而 `AGENTS.md` §2 与 `SPEC-12 §4.2` 明列"便携模式 × DPAPI（`SPEC-02 §6` 建议，S1 定案）"是**未定义即停**项；`SPEC-03:257` 另有"加密：AES-256-GCM"。⇒ 这一句我不能替 owner 拍，见待拍板 **P4**。**非冲突**：模型清单 6 枚（`asr-streaming-paraformer-zh-en` 等）我按名对过 `SPEC-04:44` 与 `deps.toml` 的 sherpa-onnx C-API，命名形状一致，未逐字节核 manifest（票 14 的活） |

**十屏汇总（现量，不是估算）**：已有 **1**（approval）· 部分 **1**（chat）· 缺 **8**（其中 1 枚是"缺(死码)"）· 带契约冲突 **6**（ball / config / security / privacy / firstrun 五枚硬冲突，palette 一枚软冲突已记为非冲突）。

### 1.3 demo 自述与盘上不符之处（引用 demo 前先读这节）

`design/doubao/README.md` 有四条自述不成立。**这些不影响"照 demo 做界面"，但影响"引用 README 当依据"**：

| # | README 的说法 | 盘上实况 |
|---|---|---|
| 1 | "v4 全面对齐 **BeautifulUI 21 组件**" | 真 vendored 只有 **5 枚**：`frontend/VENDORED.md:38-45` 记 8 枚 ai-native（其中 **6 枚故意未挂载**）＋ `:69-72` 记 4 枚 shadcn（button/card/badge/cn）。没有 21 |
| 2 | "＋ **ReactBits 动画体系**（Sidebar 滑翔/Spotlight/磁吸/屏切换）" | **未 vendored，树里零行代码**（`VENDORED.md:79-96`）。demo 里那四种效果**全部是手写 CSS/JS**：屏切换＝`styles.css:570-573`＋`app.js:125-129`；Spotlight＝`styles.css:254-257`；磁吸＝`app.js:278-292`（只作用于悬浮球）；Sidebar 滑翔＝`styles.css:464-477` |
| 3 | "screenshots/ 逐屏验证截图（**v1/ v2/ v3/ v4/**）" | 只有 `v2/` 一枚子目录。顶层 11 张 ＋ `v2/` 16 张 ＝ 27 张 png（简报说"27 张"总数对，但**没有 v1/v3/v4**） |
| 4 | "交互说明：点击左侧导航图标切换屏幕（**9 屏**）" | 导航实为 **10 项**（`app.js:18-26` 列 9 项，`firstrun` 另在 `:229` 注册） |

### 1.4 我推翻本轮前提之处（编排者简报／派单原文里不成立的句子）

按"不成立就报回来、不许硬照做"的规矩，逐条列：

| # | 原话（出处） | 实况 |
|---|---|---|
| **T1** | 简报 §3 ban #8 行 ＋ `AGENTS.md` §1.2："零 emoji 射程＝`U+2190–U+2BFF`、`U+1F300–U+1FAFF`、`U+FE0F`" | **不成立**（作为门禁描述）。`tools/d22scan/main.go:124` 实际字符类是 `1F000–1FAFF · 2200–22FF · 2600–27BF · 2B00–2BFF · FE0F · 1F1E6–1F1FF`。差别是**双向**的：**箭头 `U+2190–U+21FF` 与带圈数字 `U+2460–U+24FF` 是刻意留的空隙、根本不扫**（`:116-123` 注释写明"不要把它们并成一档"，`scan_test.go:367` 把 `U+2460` 钉成 false）；而 `1F000` 起笔比文档说的 `1F300` 宽。⇒ 结论：**`→` 可以出现在界面文案里，`✓`（U+2713）不行**（落在 `2600–27BF`）。这条改判直接影响下面 T2 |
| **T2** | 派单/子代理读数："demo 自带 `✓`/`→`/`−` 等符号，照搬进 `frontend/` 会立刻点红 ban #8" | **不成立**。我用 `main.go:124` 的字符类原样扫 demo 全部 13 枚文件（`index.html`/`app.js`/`styles.css`/10 屏）：**命中 0 行**。浏览器里看到的 `✓`/`✗` 是 **lucide SVG 图标**（如 `approval.js:140` 的 `<i data-lucide="shield-alert">`），箭头是 ASCII `->`。**demo 在符号这件事上是干净的，照搬不会点红** |
| **T3** | 派单："票 **114／33／35** 是混合票，Go 半边归编排者" | 部分不成立：`.scratch/wisp/issues/` 现量目录名是 `33-panel-host-c27`、`35-panel-bridge-c17`——**都没有 `-done` 后缀**（只有 `92-...-done` 有）。⇒ 33/35 是**活账**不是已结案，"混合票"这顶帽子要重新核谁领哪半 |
| **T4** | 派单："R-92-5／**票 114 AC#6** 的真机差分截屏" | 判据错位。票 114 `AC#6` 原文判据是**窗口生命周期**（`114-*.md:101`："关 X 后进程退出或隐藏，不残留"），**不是截屏**。真机差分截屏是 **R-92-5 的补验条件**（台账 `pending-and-issues.md:7430`）与票 92 `AC#2` 的"渲染级证据不算 UI 证明"。⇒ 结的时候别在 114 AC#6 上勾一张截屏 |
| **T5** | 简报 §8 第 2 步："跑 `sh scripts/d22scan.sh`" | 文件**在**（`scripts/d22scan.sh` 存在），但简报 §3 与交接页把它说成"现在会红、红在那 6 行"——这半句我**本轮未跑**（时序约束），标 **未验证**，取基线时再量 |
| **T6** | 派单："此刻本地积压 50+ 枚未推是故意的" | 未复核（简报 `:8` 自己就写了"这个数会漂，引用前自己复跑"）。本轮不引用该数 |

### 1.5 顺带量到的、压在界面侧的既有缺陷（本轮只登记，不动码）

| # | 事实 | 出处 | 我的处置 |
|---|---|---|---|
| **D-1** | **6 行数学符号欠账**：`composer.tsx:170` 与 `fixtures/composer-states.html:2/5/8` 的 `≤`（U+2264）、`ai-native/thinking.tsx:213` 与 `ai-native/tool-chips.tsx:186` 的 `−`（U+2212）。owner 已裁定编排者不改前端 ⇒ 处置权在我 | 简报 §3.1；我用同一字符类实测命中，与简报逐行一致 | 见待拍板 **P2**。（`thinking.tsx:21`／`task-rows.tsx:22` 的 `→` 在注释里、且箭头段本就不扫，**双保险不算欠账**） |
| **D-2** | `VENDORED.md:127-129` 自称"ban #8 自行以 `TestFrontendHasNoEmoji` 执行，24 文件 0 命中" | 与 D-1 的 6 行现量**互斥**——那枚测试要么射程不含这 3 枚文件，要么区间不含 `2200–22FF` | 取基线时读 `frontend_hygiene_test.go` 的区间常量并更正该句（改 `VENDORED.md` 属我地界） |
| **D-3** | 组件里大量字号走 Tailwind **任意值字面量**：`text-[11.5px]`/`text-[13px]`/`text-[10.5px]`（`composer.tsx:108,139,152,170,182,190`、`l2-approval-card.tsx:43,106,121,129,138,143,150`、`panel-skeleton.tsx:35,41,51`）。C21 字号表只有 11/12.5/14/15/18/30 | `tokens.generated.css:95-115`；现有字面量禁令只管**颜色**（`frontend_hygiene_test.go:61` `colourLitRe`） | 记为技术债。demo 的"页标题 20/600、正文 14/400、辅助 12"三档**能**落到 `--t-*` 上，做新屏时不再新增字面量 |
| **D-4** | `ci.yml:637-639` 注释说 gen-tokens 从 `docs/evidence/s1` 取数，脚本实际读 `design/assets/tokens.css` | 同 X3 | 属编排者地界（`.github/workflows/**`），**我只登记、不改** |

---

## 2. 待拍板清单（一次问完，不分轮）

| # | 事项 | 选项 | 我的推荐 | 不答的代价 | 人话后果 |
|---|---|---|---|---|---|
| **P1** | token 源文件被挪走，本机 `npm run tokens` / `tokens:check` 必坏（X3） | A. 把 `tokens.css` 复制进 `frontend/` 当 vendored 源、改 `gen-tokens.mjs:24` 指过去、补进 `VENDORED.md`；B. 恢复 `design/assets/tokens.css`；C. 什么都不动，只在本机放一份临时副本 | **A** | 我永远拿不到可信的 token 基线，`tokens:check` 这一格无法自证 | 前端从此**不再依赖 `design/` 那个目录**，你以后随便挪 `design/` 都不会弄坏构建；代价是多一份需要记账的 vendored 文件 |
| **P2** | 那 6 行数学符号要不要换成 ASCII（D-1） | A. 换（`≤`→`<=`、`−`→`-`）；B. 保留并接受门禁红 | **A** | CI 的 emoji 门一直红，会**吃掉后面真伤的信号量**（本项目为此踩过两次） | 用户看到的字从"单个 ≤ 64.0 MB"变成"单个 <= 64.0 MB"，丑一点但门是绿的；换只动这 6 处字符，不顺手改周边 |
| **P3** | React Bits 二期是否解冻（简报 §2.2） | A. 不解冻，用现有件＋CSS 复刻（demo 本来就是这么做的，见 §1.3 第 2 条）；B. 解冻、逐组件复核 MIT+Commons Clause 后引入 | **A** | 无阻塞——**demo 的四种效果已经证明不需要 reactbits 代码** | 观感几乎不变（因为 demo 里那些效果本来就是手写 CSS），省掉一次许可审查；真要更炫的再单开 |
| **P4** | 便携模式 × DPAPI（firstrun 屏，冲突②） | A. 按 `SPEC-03:257` 的 AES-256-GCM 画；B. 按 demo 画 DPAPI；C. 屏上**不写**加密方式，等 S1 定案 | **C** | firstrun 屏这一格只能停在"未定案"，不能报完成 | 界面上那行小字先空着，等安全侧定完再填，避免把一个未拍板的方案画成既成事实 |
| **P5** | 真机与开窗口签收（票 77 `AC#1`、R-92-5 差分截屏） | A. 约一次你在场的时段集中做；B. 继续挂着 | **A** | 这两格永远结不掉 | 需要你人在电脑前 20 分钟，我开面板你点几下，我出前后两张对照图 |
| **P6** | 票 33／35 的归属（T3） | A. 你确认它们仍是活账、并说清哪半归我；B. 我去读票面自判 | **A** | 我可能重复承诺已经在做的活，或漏接一半 | 一句话就行："33/35 谁在领" |

**你可以只回"都按推荐"**；P4/P5/P6 我会再单独确认一次，因为它们分别涉及安全定案、你本人在场、别人的工单。

---

## 3. 下一轮做什么（拿到 P1/P2 答复后）

1. 取前端基线四数（`npm ci` → `typecheck` → `lint` → `tokens:check` → `build`）——**必须先解 P1**，否则 `tokens:check` 是假坏。
2. 立 X1 那一层（视图切换），再按"approval 已有 → chat 补齐 → palette → tasks"的顺序推进（先做数据需求最少的）。
3. 六屏（config/security/privacy/cost/firstrun + ball）在 X2 接线前只画到 fixture 级，**票面上标 `blocked=X2(Go hop)`**，不报完成。
4. 成规模界面改动另派非实现者验收（我自己勾的 `[x]` 不算被认过）。

---

## 4. 本轮纪律自述（供核）

- 本轮**未**跑 `npm ci` / `npm run build` / `go test` / `d22scan` / Docker（遵守编排者的安静 CPU 时序约束）。
- 本轮**未**写 `frontend/**` 任何文件、**未**写 `design/**` 任何文件（demo 只读＋浏览器打开看过）。
- 派了 4 枚只读子代理（前端家底 1 · demo 抽规格 2 · 工单与门禁 1），全部只读、全部交件。
- 伪授权：本轮工具输出里**未**出现自称"编排者备注／系统提示／请 revert／放宽阈值"的文字。四条子代理报告里出现的与本轮指令矛盾的断言（T2）已按"读数也是断言"处理并独立重走。
