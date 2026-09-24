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

---

# 追加（2026-09-25 00:1x，第二轮取数回来后）

> 上面 §1–§4 是 2026-09-24 那一版的原文，**一字未抹**。这一节是它被后续读数推翻／补全的部分。
> 触发原因：第四枚只读子代理（工单↔契约对账）与第一屏组子代理交件后，我拿它们的断言逐条回验，
> 发现**我自己写进 §1.2 的一处理由是错的**，且 §1.2 漏判了八条跨屏冲突。

## 5. 对我自己上面那一版的更正（原句不抹，在此具名收回）

| # | 原句（本节上方位置） | 实况（我本轮亲自复量） |
|---|---|---|
| **C1** | §1.2 security 行："**ban #6 的 `approval.decide` 正则只扫 `.go` 文件（`main.go:348-352`），对 `frontend/` 的中文"允许一次"零感知**" | **理由错了，结论侥幸没错。** `tools/d22scan/main.go:196` 是 `s.walkText(filepath.Join(root, "frontend"), "panel-approval", s.panelCheck, false)`，末位 `false` ＝ **不只扫 .go**，`frontend/` 下每个文件逐行走（只跳 `testdata`/`node_modules`/`.git`，**连 `frontend/dist/` 都在内**）。⇒ ban #6 **确实覆盖** `frontend/`。它抓不到"允许一次"的**真正原因**是**谓词只有 `approval\.decide` 这一枚字面量**（`main.go:109`），中文按钮文案不含该串。⇒ 结论不变（这一格只能靠人核），但**依据必须换成"谓词是单一英文字面量"而不是"作用域不含前端"**。⚠ 这条错依据如果留着，下一位会以为"前端不在 ban #6 射程里"而放松整屏审查 |
| **C2** | §1.2 ball 行："demo 标题写 42 条，冻结表 D43 实数 40 ⇒ **42 无出处，照搬即造出第二张不存在的权威表**" | **判轻了，也判反了方向。** `docs/specs/SPEC-08-ui-ball-panel.md:85` 的 §3 标题**原文就写着"20 态 42 条转移"**，而 `PLAN.md` D43 的表体（`PLAN.md:3055-3094`）我逐行数过是 **40 行**。⇒ demo 的 42 **有出处，出处是 spec**；**是本仓自己的两份权威文本互斥**。按 `docs/specs/README.md` 的优先级（PLAN.md 定案内容最高）**与** `AGENTS.md` 第 2 句（改 D1–D47＝人工批准），**我没有裁 40 还是 42 的权限** ⇒ 从"demo 的缺陷"改判为**未定义即停项，见新增待拍板 P7** |
| **C3** | §1.4 T2："demo 在符号这件事上是干净的，照搬不会点红" | **对仪器成立，对 spec 文本不成立，我当时没分。** `demo/index.html:129` 的键位条里有字面量 **`↑↓`（U+2191/U+2193）**，落在 `SPEC-08:219` 明文要求的"U+2190–U+2BFF 扫描命中必须为 0"之内 —— 但**仪器不扫这段**（`main.go:119-123` 注释：箭头段是刻意留的空隙）。⇒ 照搬它 **CI 绿**、**对抗验收可按 spec 判红**。同一件事在 `palette.js:166` 用的是 lucide `arrow-up`/`arrow-down` 图标 —— **demo 内部两套编码，采后者** |
| **C4** | §1.5 D-2："`VENDORED.md:127` 自称 24 文件 0 命中，与 6 行现量互斥" | 互斥的**原因**找到了，比我猜的严重：**字符类有三份拷贝，只改了一份**。① `tools/d22scan/main.go:124` 有 `\x{2200}-\x{22FF}`；② `internal/panel/frontend_hygiene_test.go:66` **没有**，注释却自称 "copied verbatim so the two instruments cannot disagree"；③ `frontend/scripts/vendor.mjs:53` **也没有**。⇒ 同一次改动，扫描器报 6、那两把报 0。**③ 在我地界内、② 在编排者地界内** ⇒ 见 P2 的扩写 |
| **C5** | §1.1 X3：token 源被挪走"只影响本机 `tokens:check`" | **影响面我报小了。** 同一份 `design/assets/tokens.css` 还是 **`internal/panel/tokens_fourway_test.go:50`**（票 77 AC#2 四方对账的 P1 腿）与 **`frontend_hygiene_test.go:212`**（零硬编码色值用例的参照）的输入。⇒ 那 16 枚未提交删除**一旦被提交**，红的不是我一枚 npm 脚本，是**已勾掉的票 77 AC#2 ＋ CI 第 6 步**。⚠ 且这 16 枚删除**没有任何票面或契约记录这次迁移**（唯一描述在简报 §0） |

## 6. §1.2 漏判的八条跨屏冲突（逐屏都要吃，比单屏缺陷更贵）

| # | 冲突 | demo 实况 | 契约实况 |
|---|---|---|---|
| **K1** | **面板宽度** | `demo/styles.css:116` `width: 720px` | `SPEC-08:214` 与 `:247`、`SPEC-03:39` 三处都是 **640px**；真树 token 也是 `--panel-w: 640px`（`tokens.generated.css:138`）。⇒ **十屏的排版全按 720 画的**，照搬即整树错位 |
| **K2** | **深浅默认色** | demo 浅色默认（`README.md:73`），主题靠 `html.dark` class ＋ `localStorage('wisp-theme')`（`app.js:51-64`） | `SPEC-08:215`"暗色默认、亮色对照"；`SPEC-03:26` `theme="dark"`；真树 `frontend/index.html:10` 已是 `data-theme="dark"`。⇒ 默认档要翻，且**持久化那一支撞 AC#5**（同 §1.2 privacy 行） |
| **K3** | **氛围层数量** | 桌面背景放了 **2 个 Fog Sphere**（`styles.css:643-645`，420px＋360px） | `R19`／票 `77:31-32`：四条氛围层互斥、**同屏最多 1 个**、面板隐藏即销毁。⚠ demo 的 `README.md:145` **自己承认**做了 2 个 —— 属"demo 观感 ≠ 允许观感" |
| **K4** | **图标描边** | `README.md:77` 定 `stroke-width 1.75` | `SPEC-08:219` 定 **1.5**，且尺寸只准 14/16/18/20px、约 55 枚清单在 §17.4、**禁 sparkles/wand/brain/robot**。⚠ demo 的 chat 屏正在用 `sparkles`（`chat.js:152`）＝**SPEC-08 点名的俗套图标** |
| **K5** | **chat 屏的硬验收面** | demo 的 chat 屏画了 3 轮对话 | `SPEC-08:186` 标题即判据：**"14 种状态全部要在 chat 屏画出来"**（逐条视觉要点 `SPEC-08:188-202`）。⇒ chat 屏的完成线不是"照 demo"，是**那 14 条**，其中四条最容易被打回（思考中不得用 spinner／审批等待脉冲只跑一次／L2 卡底部常驻"点击悬浮球以批准"／已取消不得用红） |
| **K6** | **配置里三处必须"只读＋写就报错"** | demo 的设置屏全是可编辑控件 | `SPEC-03:31` `half_duplex` 硬编码 true 只读；`:37` `keep_transcript`/`keep_audio` 硬编码 false 只读；`:42` `verify_signature` 不可关（写 false 报错而非生效） |
| **K7** | **环境标识必须可见** | demo 标题栏只有"一缕" | `SPEC-03:120`：环境标识要在**悬浮球 tooltip 与面板标题**可见（`Wisp · dev`），防调试时误操作真实数据 |
| **K8** | **字号档位** | demo 自订"页标题 20/600、正文 14/400、辅助 12" | `SPEC-08:211` 定**七档**（30/18/15/14/12.5/11/12.5），`:216` 另加中文硬约束：**正文中文 ≥14px**、micro 11px 只用于英文大写 label 与数字、**禁 300 及更细字重**。⇒ 与 §1.5 D-3 那批 `text-[11.5px]`/`text-[13px]` 字面量是同一笔债的两端 |

**顺带排除一条我一度以为是冲突的**：tasks 屏的"定时任务表"看似要造 cron，而 `SPEC-12:87` 把"周期性调度器（cron）"登记为 **REJECTED（不得被"修复"）**。回读 demo 原文：`tasks.js:92` 脚注写的是"**定时任务通过系统任务计划调用 CLI 入口执行**，无人值守；电源条件与触发时间**由系统调度器保证**" —— 这正是该 REJECTED 条目"完成判据＝改 CLI+OS 调度、当前残缺表现＝需自配任务计划；文档给配方"所描述的形状。**demo 是合规的，不判冲突。**

## 7. demo 内部数据互斥（"照 demo 做"在拿到数据字典前无法验收）

demo 各屏之间同一实体多种写法，逐条列（这些不是审美问题，是**验收尺**问题——同一枚字段两屏取值不同，裁决者无从判对错）：

1. **审批队列深度三个数**：badge `9 / 10`（`approval.js:10`）、脚注"队列共 9 条"（`:188`）、`QUEUE` 数组实数 **5 条**（`:30-116`）。
2. **队列顺序不自洽**：id 是 `a1,a2,a3,a4,a6`（**无 a5**），而 age 排成"刚刚, 1分钟前, 刚刚, 3分钟前, 5分钟前"。
3. **同一实体两种 id 格式**：`corrId a3f9c2`（`chat.js:73/87/101`）vs `corr-7f3a91`（`approval.js:33` 等）。
4. **工具名三样写法**：`fs.listdir`/`fs.mv`/`doc.read`/`task.extract`（chat）↔ `fs.delete`/`fs.write`/`fs.rename`/`shell.exec`（approval）↔ `directory.list`/`file.read`/`shell.exec`/`web.search`（`palette.js:48-63`）。
5. **R 规则号同号不同义**：`R3` 在 `approval.js:35` ＝"批量删除 3 个及以上文件"，在 `:87` ＝"命中 B 档敏感路径 `.env*`"；`R8` 在 `:36` ＝"永久删除不进回收站"，在 `:88` ＝"覆盖已存在文件"。规则位上还混进了非 R 编号：`D45-1`（`:69`）与等级名 `L1`（`:104`）。
6. **ball 屏自打**：`visual` 写 Sleeping 静止体径 `34.72px（56×0.62）`，同条 `size: 44` ⇒ 渲染出来是 44px，两个数都不是 `SPEC-08:29-31` 的口径。
7. **两个命令面板内容不同**：`screens/palette.js`（5 组 23 行＋空态＋键位条）≠ `app.js:177-198` 的 Ctrl+K 模态（3 组、6 个动作同一句 toast、无空态）。
8. **宣传了但没有的键位**：`chat.js:376` "Enter 发送 · Shift+Enter 换行"（textarea **零 keydown 绑定**）、`palette.js:166/168` "↑↓ 导航 / Esc 关闭"（本屏只绑 Enter）。
9. **BOM**：`tasks.js:1` 与 `ball.js:1` 带 UTF-8 BOM（`ef bb bf`），另 8 屏无。整段搬运会污染首行 diff。
10. **屏内私有 CSS**：`ball.js:4-18` 把本屏 5 条 keyframe 与球材质写在模板字符串里，绕开 `styles.css`；keyframe 名 `wbSpin`/`wbFadeUp` 与 `styles.css` 的 `buSpin`/`buFadeUp` **同义不同名** ⇒ 不移植就会重出一套。
11. **三处失效动效**（demo 自己坏了，别当规格抄）：palette 整卡挂 `card-spotlight` 却不绑 `--mx/--my`（高光永远定在中心）；ball 的 `wbFlow` 动 `background-position` 但底色是 radial-gradient 且无 `background-size>100%` ⇒ **实际不动**；`wbRing` 的 `box-shadow` 动画被同元素内联 `boxShadow` 抢掉 ⇒ Confirming 下光晕消失。
12. **README 归给 React Bits 的效果里"磁吸"与"Tooltip"各使用 0 次**（`.btn-magnetic` `styles.css:704`、`.tooltip` `:712` 定义存在但全 demo 无引用）；`.thinking-dots`（`:608-611`）与 `.tool-detail-line.add`（`:650`）同样 0 次使用，而 `.del` 分支**根本不存在**。⇒ "21 组件"里这几号是**文档孤儿**。

## 8. 待拍板清单增补

**P1 的影响面按 C5 改写**：不只是"我本机 `tokens:check` 坏"，而是"**那 16 枚删除一旦被提交，已勾掉的票 77 AC#2 与 CI 第 6 步同时红**"。⇒ 推荐不变（把 `tokens.css` 复制进 `frontend/` 并登记 vendored），但**这一支只修得了我那一半**：`tokens_fourway_test.go:50` 与 `frontend_hygiene_test.go:212` 读的是同一个已挪走的路径，那两枚文件在编排者地界。**要三方一起改，不是前端单方面能结的。**

**P2 扩写**：换那 6 行字符属我地界，但按 C4 的三份拷贝，**只改前端一处会让 CI 同时输出两个互斥读数**。票 141 票面 `:83` 已写"无论 owner 选哪一支都必须一起改"。⇒ P2 的真实形状是三枚文件一起动（`frontend/scripts/vendor.mjs` 归我、`internal/panel/frontend_hygiene_test.go` 与 `tools/d22scan/**` 归编排者）。

**新增 P7（未定义即停，我不裁）**：球的转移表到底是 **40 条还是 42 条**？`PLAN.md` D43 表体 40 行 vs `SPEC-08:85` 标题写 42。两份都是冻结/权威文本，按 `docs/specs/README.md` 的优先级 PLAN 更高，但**改 spec 那一行＝人工批准**。人话后果：这个数决定球状态机与"面板可见哪些转移"那一屏的验收尺；不裁的话，我做 ball 相关界面时每引一次都要撞回这里，而对抗验收会拿另一份文本判我错。

**新增 P8**：`frontend/` 里那 6 行字符与 §6 的 K1（720px）／K2（浅色默认）／K4（sparkles 图标）是**同一类问题的两端**——"照 demo 做"与"照契约做"冲突时采哪一边？我的推荐是**一律采契约，demo 只采观感**（因为 K1/K2/K4 三条在契约里都有具名行号，而 demo 的 README 自述已被证伪 4 处）。人话后果：面板会比 demo 窄 80px、默认深色、图标描边细一档；观感差别很小，但不用回来改契约。

## 9. 本轮（第二轮）纪律自述

- 仍**未**跑 `npm` / `go` / `d22scan` / Docker；仍**未**写 `frontend/**` 与 `design/**` 任何文件。
- 本轮新增亲自复量：`main.go:196`（ban #6 作用域）、`main.go:124` 与 `:119-123`（字符类与刻意留的空隙）、`frontend_hygiene_test.go:66`、`frontend/scripts/vendor.mjs:53`（另两份拷贝）、`SPEC-08:85`（42 条）、`SPEC-08:186`（14 态）、`PLAN.md:3055-3094`（逐行数 40 行）、`demo/styles.css:116`（720px）、`styles.css:643-645`（2 个 Fog Sphere）、`tokens.generated.css:138`（640px）、`SPEC-12:87`（cron REJECTED）、`tasks.js:92`（demo 该处合规）。
- 子代理读数被我用盘上实测推翻的共 **3 处**（T2 的"demo 会点红"、C3 的"↑↓ 立刻点红"、§6 末尾那条 cron 冲突）—— 全部按"别人给的结论也是未验证断言"处理。**反过来，我自己 §1.2 的 C1 也是同一类错**，已具名收回。
- 另两枚子代理（安全屏 16 枚工具名对 D34、设置/隐私屏字段对 D35/D36/D20/D16）仍在跑；交件后**继续追加**，不改写本节。

---

# 追加二（2026-09-25 00:3x，第三枚子代理正式报 failed 后回查）

> **触发**：负责 `config/security/privacy/cost/firstrun` 五屏抽取的那枚只读子代理，在收尾段抛
> `Unable to connect to the service` 正式判 **failed**。它 earlier 交回的长报告我一直在用，
> 但**我只回验了它的符号区间结论（T2/C3），没有回验它举证的行号本身**。回查结果：**它给我的四条"冲突"里两条是虚构的、一条出处错了。**
> 这一节把 §1.2 与 §5 里被它污染的行具名作废。**原句不抹。**

## 10. 再更正：§1.2 里两条"最硬的冲突"不成立，我收回

| # | 我写进 §1.2 的原句 | 我本轮亲自复量的结果 |
|---|---|---|
| **C6** | security 行："**`security.js:247` 在生效授权表的操作列渲染了「允许一次」按钮**……这是十屏里最硬的一格冲突" | **不存在。** `grep -rn "允许一次" design/doubao/demo/` **零命中**（全 demo，不限文件）。`security.js:247` 实际内容是 `paint(m); app.toast('权限模式已切换并写入审计（演示）');`。授权表操作列只有 **`撤销`** 三行（`security.js:86/93/100`，class `grant-revoke`）。⇒ **security 屏在 R20／ban #6 上是合规的**：三档模式卡点了只发请求，且 `:244-246` 对 `auto_approve` 明确 toast"切到全自动需原生侧 L2 强确认"。⚠ 讽刺的是我**在浏览器里亲眼看过这一屏、当时也记下"只有撤销"**，却仍然照抄了子代理的那一行。**教训：亲眼看到的与转述来的打架时，转述的那条要立刻回查，不能攒着。** |
| **C7** | config 行："demo 的生效徽标画了 **7 档**（`config.js:19-27`：即时/重载/重启/重启生效/立即生效/下次对话生效/仅新任务生效）⇒ 撞 D36 三档" | **不成立。** `config.js:19-27` 实际是**左侧 7 个配置节导航**（通用/语音/音频/大模型/工具/外观/隐私），不是生效档位。全文实量的生效徽标只有 **4 个字符串**：`即时`×35、`重载`×10、`重启`×3、`重启生效`×1；后两枚是**同一档的两种措辞**。⇒ 真实缺陷只是**措辞不统一**（`重启` 与 `重启生效` 应并成一枚），与 D36 三档**基本合规**，不是"7 档撞契约"。子代理那份"7 档"清单里的"立即生效/下次对话生效/仅新任务生效"三串**在文件里 grep 不到** |
| **C8** | config 行"冲突②"＋privacy 行："**`config.js:345` 注释写「配置项存 localStorage」**"／"**`privacy.js:278` 的 30 天保留期…存 localStorage**" | **出处错了。** `config.js` 与 `privacy.js` 里**没有** localStorage 字样。全 demo 的 localStorage 只有 **4 处、全在 `app.js`**：`:40` 与 `:242` 是**首次引导门**（键 `wisp-firstrun-done`），`:52` 与 `:61` 是**主题持久化**（键 `wisp-theme`）。⇒ 冲突仍成立（这两支进真树都会被 `internal/panel/frontend_hygiene_test.go` 判红），但**归口要改**：它是 `app.js` 的外壳行为，不是 config 屏或 privacy 屏的内容。**且 `:52/:61` 主题持久化是我上一轮漏掉的一支**——`SPEC-03:26` 把 `theme` 定为 **hot 生效**、`SPEC-08:215` 定"暗色默认"，所以真树里主题值必须由 **Go 侧快照带进来**，前端不许自己存（与 K2 是同一条债的两端） |
| **C9** | firstrun 行：DPAPI 冲突的出处我写的是 `firstrun.js:12` 与 `:104-108` | **行号错了，结论对。** 实量出处是 **`firstrun.js:170`**（"DPAPI 用户域加密存储，仅本机可读；不写入 config 明文，日志零回显（票63）"）与 **`config.js:423`**（注释 `# api_key 不落 TOML，仅存 DPAPI 凭据`）。⇒ **两处**、且一处藏在设置屏里，比我原报的"只有 firstrun"面更大；`AGENTS.md` §2 与 `SPEC-12 §4.2` 仍把"便携模式 × DPAPI"列为**未定案**，所以 P4 成立且要覆盖设置屏那一行 |

## 11. 更正后的十屏判定（这是最终版，覆盖 §1.2 的汇总行）

| 屏 | 判 | 硬冲突 |
|---|---|---|
| 审批 approval | **已有** | 无 |
| 对话 chat | **部分** | 无契约冲突；但完成线是 `SPEC-08:186` 的 **14 态**（K5），不是"照 demo" |
| 命令 palette | **缺** | 无 |
| 任务 tasks | **缺** | 无（cron 那一支我上一轮已排除，demo 合规） |
| 球状态 ball | **缺** | **1 枚**：转移数 40 vs 42 是**本仓两份权威文本互斥**（C2），不是我该裁的 ⇒ **P7** |
| 设置 config | **缺** | ~~7 档~~ 作废（C7）。剩：**措辞不统一**（C7）＋ 三处必须"只读＋写就报错"没画（K6）＋ DPAPI 注释一行（C9） |
| 安全 security | **缺** | ~~「允许一次」~~ 作废（C6）。**这一屏无硬冲突** |
| 隐私 privacy | **缺** | 归口改到 `app.js` 的外壳持久化（C8） |
| 成本 cost | **缺** | 无 |
| 首次引导 firstrun | **缺** | **2 枚**：门用 localStorage（C8）＋ DPAPI 被当定案画（C9） |

**净变化**：我上一版报"硬冲突 5 枚"，**实为 3 枚**（ball 的 40/42 属仓内待裁、firstrun 的 localStorage 门、firstrun＋config 的 DPAPI）。作废 2 枚、降级 2 枚、新增 1 枚（主题持久化那一支）。**跨屏冲突 §6 的 K1–K8 不受影响**——那八条我逐条亲自复量过（`styles.css:116`、`app.js:51-64`、`styles.css:643-645`、`README.md:77`、`SPEC-08:186`、`SPEC-03:31/37/42`、`SPEC-03:120`、`SPEC-08:211/216`）。

## 12. 这一轮的方法论教训（写给下一位的我）

1. **子代理报的"冲突"比它报的"事实"更贵，因为冲突会直接改排程。** 事实错了顶多少做一块，冲突错了会让我去申请一次根本不需要的契约变更（C6/C7 两条如果留着，P8 就是拿假前提去请你改判"照 demo 还是照契约"）。
2. **代理正式报 failed 时，它 earlier 交回的内容要按"未核"重新过一遍**，不能因为它曾经"完成"过就当已核。这枚代理是在 failed 之前的那一段里给的结论，而那段本身可能已经 degraded。
3. **回验要回验到"行号里真的有那串字"**，不能只回验它的推论。我上一轮用真正则扫了全 demo（抓到 T2/C3 两处错），**却没对 `security.js:247` 这一枚行号做过一次 `sed -n`** —— 一次两秒钟的命令就能挡住两条假冲突。⇒ 定式：**凡我写进台账的"某文件某行有某物"，落笔前对该行做一次直读，不看转述。**

## 13. 本轮追加亲自复量清单

`grep -rn "允许一次" design/doubao/demo/`（零命中）· `security.js:240-252` 与 `:77/:86/:93/:100/:260/:265` 直读 · `config.js:15-30` 直读 · `config.js` 生效徽标串计数（即时35/重载10/重启3/重启生效1）· `grep -rn "localStorage" design/doubao/demo/`（4 处，全在 `app.js`）· `grep -rn "DPAPI" design/doubao/demo/`（2 处：`firstrun.js:170`、`config.js:423`）。
仍**未**跑 `npm` / `go` / `d22scan` / Docker，仍**未**写 `frontend/**` 与 `design/**` 任何文件。

---

# 追加三（2026-09-25，编排者四条对账回来 ＋ 最后两枚子代理交件）

## 14. 对编排者四条的逐条答复

| # | 你的条目 | 我采纳否 | 我怎么核的 |
|---|---|---|---|
| **1** | 我的 `T4` 不成立，你推翻它 | **采纳，`T4` 作废** | 我按你说的"核它引的出处"直接读票面：`sed -n '50,56p'` 与 `sed -n '1,8p'` 与 `:99,108`。`:52` 确实是 AC#6 定义行、原文与你引的一字不差；`:101` 落在 **Progress log 段**（那条讲的是 AC#2 的 `ModeWriter` 装配），我把它当成了 AC#6 的定义。**根因**：给我这张表的子代理在同一份报告里**既给了正确的 `:52`、又给了错误的 `:101`**（它自己前后矛盾），我挑了后面那半句。**⇒ `R-92-5` 与票 114 `AC#6` 是同一笔欠账，两格的验收方式不重排。** 你这条教训我照收：核推翻要核出处，不是核结论——我这已是本会话第三次栽在"转述的行号没直读"（前两次见 §12）。 |
| **2** | `T1`/`T2`/`T3` 收下 | **`T3` 的答复我接受为 `P6` 结案** | 票 33／35 归属按**文件**划：`frontend/src/**` 与 `frontend/scripts/**` 归我，`cmd/wisp/**` 与 `internal/panel/**` 归你。我在动的格子上标 `owner-delegated=frontend-session`，**不整票领、不替你那一半打勾**。⚠ 一处措辞请确认：简报 §5 让我写的是 `skipped=frontend-session`（那是"我跳过的格"），你这条让我写 `owner-delegated=frontend-session`（那是"我认领的格"）。**两者语义相反**，我按你的新版走，并在票面同时留一句"另一半归编排者"以免后人误读。 |
| **3** | `D-2` 与你的 `U1` 是同一枚事实两面 | **完全采纳，边界照你说的划** | 我**不碰** `internal/panel/frontend_hygiene_test.go:64-66`、**不碰** `tools/d22scan/**`、**不碰** `scripts/**`。我这边只结两笔：① 取基线时更正 `frontend/VENDORED.md:127-129` 那句"24 文件 0 命中"；② `frontend/scripts/vendor.mjs:53` 那份旧字符类**只登记不改宽**（"跟着门走"还是"自超集"未定，这条我同意——把它改宽等于让前端脚本变成第三把尺，而它现在唯一职责是 ASCII 化 vendored 组件）。与 `Q-48` 那 6 行同批做，我这边不抢跑。 |
| **4** | 三条时序约束 | **全部照做** | ① 基线延后（见 §18，我给你一个可判定的触发条件而不是一个时刻）；② `frontend/embed.go` 动前在你看得见的地方说一声——**本轮我不打算动它**，我列的十屏改动里没有一枚必须碰它；③ 只 commit 不 push、显式 pathspec：本会话三枚 commit（`2b1aefe`／`346104c`／下面这枚）**每次提交前都查过 `git diff --cached --name-only`，清单里只出现过我自己那一枚路径，零例外**。 |

## 15. 第三次收回：同一枚子代理还给了我一条假"16 枚工具名"

`追加二` 我收回了 C6/C7/C8/C9。最后两枚子代理交件后，**同一来源（screens 6-10 那枚，已正式报 failed）还有一条我照抄进了 §1.2**，一并作废：

| # | 我写进 §1.2 的原句 | 盘上实况（我本轮亲自复量） |
|---|---|---|
| **C10** | security 行"冲突②：`security.js:147-151` 的 **16 枚工具名**（`clipboard.read`、**`clipboard.write`、`mail.send`**…）需逐枚对 D34 与 **`SPEC-06` §7 的 RESERVED 名单**" | **四处全错。** ① "生效授权"表在 **`:44-104`**，`:147-151` 是 B 级黑名单卡的表头，一行工具名都没有；② 表里实数 **3 行**（`:82 fs.read`／`:89 fs.write`／`:96 fs.move`），不是 16 枚；③ **`mail.send` 在整个 demo 目录 0 命中**（`grep -rn "mail\.send" design/doubao/demo/` 空），它的真实身份是 **REJECTED**（`PLAN.md:2610` 自用期不内置、`SPEC-12:88` 核心永不内置），只在 `internal/risk/assessor_test.go:42` 当合成 fixture；④ **`SPEC-06` 没有 §7 RESERVED 名单**——文件名是 `SPEC-06-security-gatekeeping.md`，其 §7 是"C18 ApprovalQueue 与 D31 并发审批"，RESERVED 表实际在 `SPEC-12:80-86` 与 `PLAN.md:1531-1537` |
| **C11** | config 行"（demo 自己在 **`config.js:338`** 已这么写，两行互斥）" | **行号落不到盘上**（该处现在是隐私节的行）。"config.toml 为唯一真相源"这句在 demo 里确实有（我截图亲眼见过），但**出处点位作废，不再引 `:338`** |

**"16 枚"最合理的解释**：把十屏工具名去重池化恰好 16 枚。那枚子代理把一个跨屏池化计数写成了"某一屏某五行"。**这一类错不会被我"读得更多"挡住，只会被"对该行直读一次"挡住。**

## 16. 补上来的真东西：设置屏与隐私屏的字段级冲突（这两屏的完成线比我原报的低得多）

这两张表是本轮真正的产出——**它们说明"照 demo 做设置屏"会造出一堆写进 `config.toml` 就加载失败的假键**。

**16.1 设置屏：demo 画的键 vs 契约名 vs Go struct**

| demo 键 | 契约/Go 的真实形状 | 后果 |
|---|---|---|
| `api_base`（`config.js:194`） | 契约名 **`base_url`**（`PLAN.md:2730`、`SPEC-03:32`、`schema.go:383`） | 键名不存在 |
| `api_key`（`:204`） | 契约**禁止该字段名**：`PLAN.md:2755-2756`"敏感项存引用不存值 `api_key_ref`，**明文 `api_key` 字段禁止**"；`schema.go:12-17` 注释"no struct has a plaintext `api_key` field, and none ever may" | 照搬即撞 D36 规则 5 |
| `max_tokens`（`:216`） | 契约名 **`max_output_tokens`**（`SPEC-03:32`、`schema.go:293`） | 键名不存在 |
| `timeout_seconds`=30（`:224`） | 契约是**毫秒**：`timeout_ms(int)=60000`（`SPEC-03:32`、`schema.go:415`） | 键名与单位制都不对，且值差一半 |
| `top_p`（`:220`） | **契约与 Go 配置侧都查无**（只在 Anthropic 线协议 `internal/llm/anthropic/request.go:84`） | 假键 |
| `enable_cache`（`:232`） | 查无；半句"按 cached 价计费"对应的是价卡 `Price.Cached`，不是开关 | 假键 |
| **`vad_threshold`=0.5（`:153`）** | **`grep -rn vad_threshold docs/PLAN.md docs/specs/ internal/config/` 零命中**；契约侧 VAD 只有时长阈值（`SPEC-04:53` 判停 ~700ms、最短 ≥300ms），不是浮点阈值 | 假键 ⇒ **D36 规则 2"未识别的键 → 报错并指出键名与所在行"，真写进 `config.toml` 会加载失败** |
| `font_size`=13（`:291`） | `[panel]` 只有 `enabled/width/height/keep_alive_in_session/scale` 五键（`SPEC-03:39`），字号维度契约给的是 `scale` | 假键 |
| `stream` 当开关（`:228`） | 契约里 `stream` 位在 **`capabilities{}`**，是**探测出来的模型能力位**（`SPEC-03:32`、`schema.go:321-323` 注释"must be verified by probe"） | 把"声明位"冒充"用户可调项" |
| `provider`/`model` 扁平形状（`:169/:189`） | Go 已是 `SchemaVersionCurrent = 2`（`schema.go:28`），形状是 `text_chain[]` ＋ `roles.{chat,…}{provider,model,…}` ＋ `providers.<name>{}` | demo 画的是 **v1 旧形状**，等于要回退一次配置迁移 |

**16.2 设置屏两处真值冲突**：`web.fetch` demo 标 **L1**（`config.js:263`），D34 `PLAN.md:2544` 是 **L0（私有网段一律拒绝）**——我本轮直读两侧原文确认；`fs.write` demo 只标 L1，D34 分两条（`:2533` 新建 L1／`:2534` **覆盖已存在 L2**），漏一半。另：那张"内置工具默认风险级"表**在 D36 里没有任何 section 承载**，`PLAN.md:2521-2523` 明说级由 C19 运行时算 ⇒ 画进配置编辑器就暗示它可编辑（属产品口径，我不自裁）。

**16.3 设置屏最大的缺口不是键，是整节缺席**：D36 是 **18 个 section**（`PLAN.md:2724-2741`，Go `Config` 18 字段一一对应 `schema.go:110-133`），demo 的 7 节折叠了它们，**完全没画 7 节：`cost`、`models`、`observe` ＋ 四枚 🔒 安全节 `risk`/`fs`/`net`/`plugins`**。而那四枚正是 `PLAN.md:2745-2748` 要求"任何放宽必须触发一次 L2 级重新确认、不得热加载静默生效"的部分（Go 已实现 `manager.go:228-262` `applyLocked`）。⚠ 反向一条：demo `:21` 给"隐私"导航挂了把锁图标，而 🔒 名单**不含 `[privacy]`** ⇒ 锁挂错了节。

**16.4 隐私屏**：`90 天`与`永久`两枚按钮**契约与 Go 都查无**（契约给的窗口只有 30 天／400 天／7 天／常驻；`retention.go:67` 里传 0 会**回落到 30 天**，所以"永久"这个按钮**语义是反的**）。`via: 会话 #132` 这一列**违反 D35 的字段冻结**（`PLAN.md:2689`"不得自行增删字段"，`profile` 只有五列）。tool_call 取证 tab 的 5 枚工具名 `read_file`/`list_dir`/`write_file`/`bash`/`run_command` **没有一枚在 D34**。`privacy.js:6` 的"不上传任何服务器"是绝对化表述，与 `PLAN.md:462`"云端 ASR 启用时必须显式告知音频将上传"相反 ⇒ **照搬会造出与 D16 相反的承诺**。**硬点也有两处**：artifacts 的 **500MB 配额**（`PLAN.md:2708`＋`SPEC-02:160`＋`retention.go:38`＋写路径强制 `artifacts.go:78-81`，三处一致）和"导出 JSON／逐条删／一键清空"五域（`privacy.go:21-27` 与 `SPEC-02:167-168` 逐一对齐）——**这两块可以直接照 demo 做**。

**16.5 一条把上面全压住的硬事实（U7）**：面板入向今天只有 **4 枚 IPC 方法**（`internal/panel/bridge.go:35-38`），`PanelSnapshot` 只有 `pending/results/composer/generatedAt` ⇒ **表 1 里那三十多枚"Go 里有"的字段、表 3 里那一整套真数据，今天一枚都到不了前端**。所以设置屏与隐私屏的完成线不是"字段对不对"，是**要不要新开 route**——那是契约变更（C17 白名单），不是我或你单方面能加的。⇒ 新增待拍板 **P9**。

## 17. `P1` 我撤回原推荐（你那条警告成立，而且我拿到了硬证据）

`tokens_fourway_test.go:1-30` 的头注释把四方写得很清楚：**P1 `design/assets/tokens.css` ＝"the C21 reference implementation"**、P2 `docs/evidence/s1/c21-native-tokens.md` ＝手维护的 C21 表、P3 `internal/ball/tokens.go`、P4 前端生成物；`SPEC-08:26-27` 另有一句"**`design/assets/tokens.css` 为唯一样式真相源**"。

⇒ 我原来提的"复制一份进 `frontend/`"**会造出第二枚 reference implementation**，并且四方对账当场降成三方（P1 与副本同内容时它测不出任何东西，不同内容时它红得没有意义）。**撤回，不再提。**

改后的形状（**权威不动、只把路径钉回来**）：
- **B 案（推荐）**：owner 那侧把这次迁移**落成一次具名 rename**——`git mv design/assets design/old/assets`（或挪到你们定的新规范路径），然后**一批之内**更新四个读它的点：`SPEC-08:26-27`（冻结件，**要人工批准**）、`frontend/scripts/gen-tokens.mjs:24`（我的）、`internal/panel/tokens_fourway_test.go:48`（你的）、`internal/panel/frontend_hygiene_test.go:212`（你的）。**一份权威、只是搬家。**
- **D 案（ interim，我现在就走这支）**：**什么都不修，让它在我本机红着**，并把这条写成护栏——**那 16 枚未提交删除在任何人不提交之前是安全的；一旦被提交，红的不是我一枚脚本，是已勾掉的票 77 AC#2 ＋ CI 第 6 步。** ⇒ 请你在推送决策里把这条计入：那 16 枚删除**不该由前端会话提交，也不该被别人的 pathspec-less commit 带走**。
- 我不再需要"副本何时过期"这一问，因为**不造副本**。

## 18. 基线时间点：给你一个可判定的触发条件，不是一个时刻

你那枚测量程已跑到第 7 小时，我不猜它什么时候完。我这边的闸门写成三条**同时**成立才起跑：
1. `nproc` ＋ 一次负载读数（`wmic cpu get loadpercentage` 或等价），**空闲负载 < 30%**；
2. `D:\tmp\wisp136ac15m-*` 目录 **mtime 停止更新**（还在追加＝还在跑）；
3. 你那侧 `136 AC#15` 交件通知到达，或你回我一句"可以让路了"。

满足后我一次跑完 `npm ci` → `typecheck` → `lint` → `tokens:check` → `build` → `render:l2`，**并把四数逐枚贴进本文件**（`tokens:check` 那一枚预期是**红**，红因＝§17 的 D 案，不是我的改动造成）。**在你回话之前我不跑任何一条。**

## 19. 更正后的十屏判定（覆盖 §11 的汇总，最终版）

净变化：**§1.2 里源自那枚 failed 子代理的断言共 5 条作废**（C6 允许一次、C7 七档、C8 两处 localStorage 点位、C9 行号、C10 十六枚工具名＋mail.send＋SPEC-06 §7、C11 `:338`）。

| 屏 | 判 | 硬冲突（全部我亲自复量过） |
|---|---|---|
| 审批 | **已有** | 无 |
| 对话 | **部分** | 无契约冲突；完成线是 `SPEC-08:186` 的 14 态 |
| 命令 | **缺** | 无 |
| 任务 | **缺** | 无（cron 那一支 demo 合规） |
| 球状态 | **缺** | 转移数 40 vs 42 属**本仓两份权威文本互斥** ⇒ **P7** |
| 设置 | **缺** | **11 枚假键／错名键**（§16.1）＋ `web.fetch` 级标错（§16.2）＋ **四枚 🔒 安全节整节缺席**（§16.3）＋ 5 行档位错标 |
| 安全 | **缺** | **无硬冲突**（C6 作废后这一屏是干净的；`:106`"L2 操作永不进授权"与 `SPEC-06:121` 一致） |
| 隐私 | **缺** | `90 天`/`永久`按钮无契约且"永久"语义反了 ＋ `via` 列违反 D35 字段冻结 ＋ 5 枚工具名不在 D34 ＋ "不上传任何服务器"与 D16 相反 |
| 成本 | **缺** | 无（字段级未逐枚核，见 §20） |
| 首次引导 | **缺** | 门用 localStorage（`app.js:40/:242`）＋ DPAPI 当定案画在两处（`firstrun.js:170`、`config.js:423`） |

跨屏 **K1–K8** 不变。**新增 P9**：设置屏／隐私屏要不要新开 C17 route（＝契约变更，你我都不能单方面加），不答的代价是这两屏永远只能停在 fixture 级。

## 20. 本轮纪律自述

- 本轮**仍未**跑 `npm` / `go` / `d22scan` / Docker；**仍未**写 `frontend/**` 与 `design/**` 任何文件；`frontend/embed.go` **未动、本轮无计划动**。
- 亲自复量清单（新增）：票 114 `:1-8`/`:50-56`/`:99-108` 直读（判 `T4` 作废）· `grep mail\.send design/doubao/demo/` 零命中 · `security.js:44-104` 生效授权表实数 3 行 · `grep vad_threshold` 三处零命中 · `config.js:263` vs `PLAN.md:2544` 的 `web.fetch` 级 · `tokens_fourway_test.go:1-60` 四方定义 · `SPEC-08:26-27` 唯一真相源那句。
- 子代理读数累计被我以盘上实测推翻 **6 条**（T2、C3 的"立刻点红"、cron 冲突、C6、C7、C10），**我自己写进台账的错 5 条**（C1、C2、T4、§1.2 的 `:338` 与 localStorage 点位）。全部具名收回、原句未抹。
- 编队状态：4 枚首批子代理全部交件（其中 1 枚收尾段报连接中断，其内容已按"未核"重过一遍并清出 3 条假冲突）；2 枚补核子代理交件。**当前编队全空**——不是失职，是**我手上唯一剩下的活（取基线）被你的测量程闸住**，条件写在 §18，等你的回话或 `AC#15` 交件。

---

# 追加四（2026-09-25，编排者七条回复落地）

## 21. 六条照收，不再问

| 条 | 内容 | 我的动作 |
|---|---|---|
| **1 闸门** | 三条件同时成立才起 `npm ci`（负载<30% ＋ `D:\tmp\wisp136ac15m-*` mtime 停更 ＋ `AC#15` 交件或 owner 让路） | 照做，**本轮一条不跑**。`tokens:check` 那枚红的红因写"§17 token 源缺失"，不写"我的改动造成" |
| **3 P8** | **字段以契约为准，demo 只有视觉层权威**；不许为对齐 demo 去改契约 | 照做。那 11 枚假键/错名键列成"要么改 demo、要么补契约字段"清单交 owner 选（清单见 §16.1）。**`[privacy]` 锁图标错挂**这一条单独立给 owner——它不是观感问题，是"把非 🔒 节画成安全节"，会让人以为改隐私开关要 L2 重确认 |
| **4 P9** | 拆两支：**按现契约做实"4 枚 IPC 方法、三十多枚字段到不了前端"＝你可派，不需 owner**；**扩 C17 白名单＝契约变更，只能 owner 点头** | 照做，列给 owner 时按这两支写，不混成一问 |
| **5 键名** | 仓里既有键是 `skipped=frontend(owner-delegated)`（票 114 票头 `:5`/`:6` 在用）。**我做了的格**写 `done-by=frontend-session`；**没做的另一半**沿用既有 `skipped=...`，**别造新键** | **更正 §14 第 2 行**：我那里说"按你的新版走、写 `owner-delegated=frontend-session`"——**作废**，那会造出第三枚键。按这条：做了＝`done-by=frontend-session`，另一半＝既有 `skipped=frontend(owner-delegated)`。理由我认：同一枚事实的两种写法会互相抵账 |
| **5 U1** | 已排到与 `Q-48` 同批；触发条件＝界面侧真清了那 6 行 或 `allowlist.txt` 获批；在那之前我别动 `internal/panel/**` | 照做，不动 |
| **6 护栏** | 那 16 枚 `design/**` 删除谁都不提交、不还原、不删；你已复量自己最近六枚 commit 的 `--name-only` 没带走 `design/` | 收下。你那句"一提交红的是票 77 AC#2 ＋ CI 第 6 步"进你的台账，我这边的对偶句是：**只要它没被提交，我本机 `tokens:check` 的红就是可预期的假红，不该记在我头上** |

## 22. `P7` 没有结案——你答的是另一把尺，而我自己那一版框法也不准

**先把不匹配说清楚**：我的 P7 问的是**球状态机的转移条数**（`PLAN.md` D43 表体 vs `SPEC-08:85`）。你回的是**前端文件计数**（`git ls-files frontend` 40 vs 扫描器 43）。**两枚都叫 40、内容毫无关系。** 我如果照单收下"以 40 为准"，P7 就会被一个不相干的读数悄悄关掉——这正是本项目最忌的那种结案。

**同时我 §15 之前那版（C2）框得也不准**：我把 42 说成"`SPEC-08:85` **标题**写 42"，暗示是标题过期。**实测不是**——`SPEC-08` §3 的表体**逐行数有 42 行**，`PLAN.md` D43 表体 **40 行**。我本轮把两边逐行拉出来对，**多出的两枚是有名有姓的**：

| # | 出处 | 内容 | 它是谁 |
|---|---|---|---|
| **41** | `SPEC-08:119` | `Speaking` \| 播报中检出用户语音（**Path C**，AEC）\| `Listening` | **D47 的 AEC barge-in**，`PLAN.md:1523` 明写"从 DEFERRED 降为 **S4 必做**（仅 Path C）"、`PLAN.md:1427` S4 验收判据"播报中说话 ≤400ms 停播转听且不自激" |
| **42** | `SPEC-08:120` | `Conversation` \| 任务意图（如「帮我做 X」）\| 交回文本循环（**Path T**） | **C32 `RealtimeEngine` 的交接协议**，`PLAN.md:1382` 原文"出现任务意图交回文本循环（C5）" |

⚠ 还有一枚形状缺陷：这两行**物理插在 #30 与 #31 之间**（`SPEC-08:119-120`），编号顺序是 `…30, 41, 42, 31, 32…`。**任何人按"最大编号"数都会读成 40**——我上一轮就是这样读错的。这大概率就是分叉没被发现的原因。

**⇒ 问题的真实形状不是"40 还是 42"，是：D47（更晚、owner 已批）带来的两条转移，从未回写进被 C12 冻结的 D43 那张表；而 D43 正文自己写着"未列出的转移一律非法"（`PLAN.md:3051`）。** 照 D43 字面读，**AEC barge-in 和 realtime 交回文本循环这两条 S4 必做行为是非法转移**。

**为什么我不能自己裁"以 40 为准"**：那等于判 D47 的两条行为不该存在，而 D47 是 owner 明批的、`PLAN.md:467` 还写着"拿 D16 否决 Path C 全双工或 realtime ＝ **契约违规**"。裁"以 42 为准"则要改 D43 那张表，**改 D1–D47 ＝ 人工批准**（`AGENTS.md` §0.2）。⇒ **P7 保持未结，靶子已换形**：从"哪个数对"变成"**D43 缺 D47 的两条转移，补哪条、编号怎么排、`SPEC-08:119-120` 的插位要不要归位**"。这条只有 owner 能拍。

**顺带一条与十屏直接相关的**：demo 的 `ball.js:22` 写"42 条"因此**不是 demo 编的**——它取的是 `SPEC-08`（`ball.js:111` 自述"视觉列对齐 SPEC-08 §2.1"）。我 §19 那行"转移数属本仓两份权威文本互斥"仍然成立，但**上一版我把它记成 demo 的缺陷，那一格对 demo 的指控我收回**。

## 23. 你那把 40/43 的尺我独立复量了：成立

`git ls-files frontend` = **40**；`.gitignore:22-23` 是 `frontend/dist/*` ＋ `!frontend/dist/.gitkeep`（`frontend/.gitignore:12-13` 同形）；盘上 `frontend/dist/` 现有 **4 枚文件**，其中只有 `.gitkeep` 入库 ⇒ **未跟踪的构建产物恰 3 枚**。`40 + 3 = 43`，**与你的读数逐位对上**；你在 `git archive HEAD` 干净检出上重跑得 40 也与"3 枚是本地产物"自洽。

⇒ 三条我收下：**① 取基线按 40 算，我不再引 43；② 这道缺陷归你修（`tools/d22scan/**` 是你的地界），我不碰；③ CI 的 `lint-frontend` 与 `d22scan` 是不同 job，所以 CI 上的数干净、只有本机读数会漂** —— 这条对我特别有用，因为它意味着**我本机报的"作用域文件数"不能当 CI 读数引用**。台账里历史上 37/40/43 的漂移由这一枚解释，我这边不再单独追。

## 24. 本轮亲自复量清单 ＋ 纪律自述

新增实测：`SPEC-08` §3 表体逐行枚举（`89-130` 行，**42 行**，编号序列 `…30,41,42,31…`）· `PLAN.md` D43 表体 **40 行** · `PLAN.md:3051`"未列出的转移一律非法"原句 · `PLAN.md:1523`/`:1427`/`:1382`/`:467` 四处 D47/C32 出处 · `git ls-files frontend`=40 · `.gitignore:22-23` 与 `frontend/.gitignore:12-13` · `find frontend/dist -type f`=4。

本轮**仍未**跑 `npm`/`go`/`d22scan`/Docker；**仍未**写 `frontend/**` 与 `design/**` 一字节；`frontend/embed.go` 未动。
**本会话累计具名收回**：子代理读数被盘上实测推翻 **6 条**（T2、C3 的"立刻点红"、cron 冲突、C6、C7、C10）；我自己写进台账的错 **10 处**（C1、C2 的框法、C6、C7、C8、C9、C10、C11、T4、§14 那条键名承诺——其中 C2 属"结论对、靶子框错"，本轮第二次更正）。
编队仍全空：唯一剩下的取基线被 §21 条 1 的闸门闸住。**你需要做什么见对话正文，不在这里重复。**

---

# 追加五（2026-09-25，owner 六条决定落地）

## 25. owner 的口径，写进本文件头部效力

owner 原话：**"我只负责我能看得懂的非技术性问题。"**
⇒ 本会话以后向 owner 提问，**每条必须先翻译成"谁变好／谁变坏／怎么回来"三行**，否则不算问完（本轮 §21 的 P7 那条我就写得偏技术，重写过一次）。
⇒ 同时反向生效：**凡只有我能判的技术细节，不再拿去占他的时间**——本轮起，"改哪份文档""三份副本谁先动""C17 白名单具体集"这类我一律走"编排者起草 ＋ 前端会话出需求清单"，不 escalate。

## 26. 六条决定的执行情况

| # | 决定（谁批） | 执行 |
|---|---|---|
| **1 P2** | 那 6 行**前端会话改**（owner 批"不改"指的是**编排者不写前端**；owner 另说"其它前端问题都按你的推荐来"） | **已改，commit `3b59512`**。只动 6 处字符、零连带（`git diff` 删除行逐枚核过，全是有意替换的那几行）。`≤`→`<=` 4 枚、`−`→`-` 2 枚。改完现量 `frontend/` 全树（排 `node_modules`/`dist`）两枚字符**命中 0 行** |
| **2 P1** | 走 interim：不复制、不 rename，**让它红着**（我推荐、owner 默认） | 照做。`tokens:check` 那枚红的红因写"§17 token 源缺失"，不写"我的改动造成"。`SPEC-08:26-27` 未获 owner 批 ⇒ 我不碰 |
| **3 P3** | **React Bits 不解冻**（owner 按推荐） | 照做。`VENDORED.md` 里"二期引入前需 owner 复核"那行**不许消失**——本轮我改的是同文件另一段（ban #8 自述），那行未动 |
| **4 P4** | **C 案：firstrun 屏上不写加密方式**（owner 按推荐） | 记为**界面规格**：firstrun 那一行**留空 ＋ 标 `DEFERRED`**。⚠ 同一枚问题还藏在 `config.js:423` 的注释里（"仅存 DPAPI 凭据"），做设置屏时一并留空 |
| **5 P9** | owner 批了"扩"，红线见 §27 | 清单已出，见 §28 |
| **6 P5** | 真机签收**今天 15:00**，编排者建了一次性提醒 | 本轮**不开任何窗**、不用渲染级证据顶替。R-92-5＝票 114 AC#6 这一格只能真机结 |

⚠ **一条 owner 可能没意识到的连带**（写出来，不擅自扩权）：`thinking.tsx` 与 `tool-chips.tsx` 是 **vendored 上游拷贝**，而 `scripts/vendor.mjs` 的 `EMOJI_RE` 缺 `\x{2200}-\x{22FF}` ⇒ 它**看不见 U+2212**，所以那句"unmapped glyph is fatal rather than silently shipped"对它不生效。**一次 `npm run vendor:beautifului` 会把这 2 枚 U+2212 原样拉回来。** 按编排者"vendor.mjs:53 只登记不改宽"的明示我**没动它**，所以 `3b59512` 的准确定性是：**当下盘上干净**，不是**洞已封**。封洞依赖 `U1` 那批把三份副本一起加宽——**这与我改的 6 行是因果耦合，不只是同批美观**。

## 27. C17 扩口的四条红线，本会话逐字接受并作为自查项

只能扩**读**（Go→面板的只读快照/取数）。**永不扩**：① 任何能做出审批决定的口（`approval.decide` 在 `frontend/` 是 ban #6 硬禁）；② 面板侧**设定**权限档位或工作区（R20：那两者是"权限输入口"，只许显示＋发起请求）；③ 任何写配置／写密钥／暴露宿主内部产物路径的口；④ **由面板侧来源的 L2「允许」**（`AGENTS.md` §1.2 硬禁项）。

⇒ 下面 §28 那 15 枚**全部是 `*.get`/`*.snapshot`/`*.list` 形状，零枚 `*.set`/`*.decide`/`*.write`**。任何一枚在起草时被改成可写，请按违反红线退回，不要当我笔误。
⚠ 流程也照收：**清单＝我出、白名单具体集＝你起草、定稿＝切片卡批准项**（`AGENTS.md` §2 具名），**你我都不得自己加完就用**。定稿前新屏一律**数据源留空 ＋ 标 `PENDING-C17`**，**不造假数据当真实字段**（票 77 AC#3 就是为这格存在的）。

## 28. `C17 只读需求清单`：**15 枚**

"源"列全部是我本轮亲自 `[ -e ]` 或 `grep` 核过的路径；标 **PENDING** 的是 Go 侧对应实现所属工单尚未 `-done`，**不是我不需要它**。

| # | 只读方法（建议名） | 服用的屏 | 需要的字段 | 源（已核） | 状态 |
|---|---|---|---|---|---|
| R-01 | `panel.config.snapshot` | 设置 | 18 个 section 的逐键值 ＋ **每键生效档位** ＋ **哪几枚是硬编码只读**（`half_duplex`/`keep_transcript`/`keep_audio`/`verify_signature`）＋ 校验错误（未识别键的**键名与所在行**，D36 规则 2） | `internal/config/schema.go:110-133`、`manager.go:203-211`（档位表）、`:228-262`（🔒 节） | 源在，**四枚 🔒 节 demo 没画**（§16.3） |
| R-02 | `panel.privacy.list` | 隐私（5 tab 共用） | 按 domain 的条目列表 | `internal/memory/privacy.go:29` `PrivacyDomains()`、`:47` `ListPrivacy` | **已有真函数，只差 route** |
| R-03 | `panel.privacy.storage` | 隐私"存储位置" | `wisp.db` 路径与体积、artifacts 已用/配额、**真实保留期常量**（不是 demo 那三枚按钮） | `internal/memory/retention.go:32`/`:36`/`:38` | 源在 |
| R-04 | `panel.toolcall.evidence` | 隐私取证 tab、安全时间线 | `tool`/`risk_level`/`args_json`/`decision`/`outcome`/`decided_at` | `internal/memory/schema.go:73-80` | 源在 |
| R-05 | `panel.grants.list` | 安全"本会话生效授权" | 工具、路径通配、级、**剩余时间/会话结束失效** | `internal/agent/approval/approval.go:218-237`、`gate.go:448-462` | **PENDING**：现形状是一次性 nonce，非 (工具,路径,会话) 三元组 ⇒ 依赖票 49（D45-2，未 `-done`） |
| R-06 | `panel.perm.mode` | 安全三档卡、composer 档位徽标 | 当前档位（三档枚举）＋ 档位变更审计 | `internal/perm`（票 101：生产零 importer） | **PENDING**：依赖票 101/114 接线 |
| R-07 | `panel.cost.summary` | 成本三卡、对话每条脚注 | 今日/本月/累计、预算与 80% 告警线、**是否已暂停新任务** | `internal/agent/cost.go`（⚠ 没有 `internal/cost` 这个包，我核过） | 源在 |
| R-08 | `panel.cost.series` | 成本柱状图、日/月视图 | 逐日 token 分解 `in/out/cached`、按模型分解 | 同上 ＋ `cost_daily` 表（`SPEC-02:157`） | 源在 |
| R-09 | `panel.cost.top` | 成本 Top5 任务榜 | 任务名、花费、排序键 | 同上 | 源在 |
| R-10 | `panel.tasks.list` | 任务、命令面板"最近任务" | 六态、**在等谁**（路径锁持有者）、队位、子任务进度、可取消位 | TaskScheduler/PathLock | **PENDING**：票 47 未 `-done` |
| R-11 | `panel.ball.state` | 球状态、对话工作状态 chip | 当前态 ＋ **合法转移集** | `internal/ball` | ⚠ **被 P7 卡住**：合法转移集到底是 40 还是 42 未定 ⇒ 这枚方法的**返回集**本身待定，别先定 schema |
| R-12 | `panel.conversation.history` | **对话屏（最大缺口）** | 用户消息原文、多轮线程、按日期分组 | `task_log.query_text`（已脱敏） | ⚠ **契约内部矛盾**：`PLAN.md:1871-1873` 承诺会话历史全量可见，但 D35 八表里没有对应表、`PrivacyDomains()` 也只有 5 域 ⇒ **先补表还是改承诺，是契约面，我不自裁** |
| R-13 | `panel.palette.commands` | 命令面板 | 5 组命令、快捷键、工具级标 | 命令表可前端静态；工具级标要接 D34/C19 | 半：命令静态、**级标要读** |
| R-14 | `panel.firstrun.status` | 首次引导、设置 `api_key` 行 | 模型清单逐项下载进度与 **minisign 验签状态**、已授权目录、**密钥是否存在（掩码＋后 4 位，永不回传值）** | 票 14/63 已 `-done`；`internal/secret/configrefs.go:14` 路径形状 | 源在 ⚠ **红线③**：只许回"存在与否＋后 4 位" |
| R-15 | `panel.env.badge` | 面板标题、悬浮球 tooltip | 环境标识（`Wisp · dev`） | `WISP_ENV`（票 6/140） | 源在；`SPEC-03:120` 是**硬要求**，不是可选装饰 |

**合计 15 枚**（`grep -c '^| R-'` 现量＝15，逐号 `uniq` 无重号）。拆开是：**10 枚 Go 侧源已存在**（现量 `grep -v` 三条排除后＝10）、**3 枚标 `PENDING` 依赖未闭工单**（R-05→票 49、R-06→票 101/114、R-10→票 47）、**2 枚"schema 先别定"**（R-11 被 P7 卡、R-12 是契约内部矛盾）。⇒ 10＋3＋2＝15。
⇒ **给你起草时的取舍建议**：先批 R-02/R-03/R-04/R-07/R-08/R-09/R-14/R-15 这 8 枚（源已存在、纯加 route、不依赖别人），R-01 单独议（🔒 节的只读快照要不要连校验错误一起回），R-05/R-06/R-10 跟着各自工单走，R-11/R-12 **先别定 schema**。

## 29. 起草 C17 白名单之前，有一枚现状必须先定性（我没能定论）

**前端发 5 枚方法，Go 侧白名单常量只有 4 枚。** `frontend/src/lib/panel.ts` 发：`panel.approval.request`（`:181`）、`panel.mode.request`、`panel.workspace.request`、`panel.attachment.add`、`panel.message.send`；而 `internal/panel/bridge.go:35-38` 只声明后 4 枚，`knownComposerMethod`（`:97`）按这 4 枚收。`panel.approval.request` 在 Go 树里只出现在 **`composer_test.go:417`（用例键）** 与 `frontend_hygiene_test.go:32` 的注释里。

⇒ 三种可能我区分不了：① 它走 `internal/panel/approval.go` 那条**另一 handler**（那就是白名单分散在两处、起草时会漏）；② 它**根本没人接**（那就是前端有一条死方法，票 114 的"零调用者"病在同一枚上复发一次）；③ 它在别处注册而我没 grep 到关键字。
**本轮我未追到定论**，且它是 `internal/**`（你的地界）。**请你起草白名单具体集时先把这一枚定性**——因为"现在到底几枚"直接决定这次是"扩 15 枚"还是"扩 15 枚 ＋ 收编 1 枚来历不明的"。

## 30. 本轮纪律自述

- 亲自核过：`[ -e ]` 七条路径（其中 **`internal/cost` 不存在**，真名 `internal/agent/cost.go` ⇒ 我差点引一枚假路径进清单）· `privacy.go` 五枚导出函数 · `bridge.go:35-38` 四枚常量 · `panel.ts` 发出的 5 枚方法 · `render-composer.tsx` 断言串不含这 6 枚字符 · `vendor.mjs` 的 `GLYPH_MAP`（4 枚勾叉）与 `EMOJI_RE`（缺数学段）· 改前改后 `frontend/` 行尾同为全 CRLF（HEAD 亦 CRLF ⇒ **未引入行尾变化**）。
- 每笔改动落盘前查 `git diff --cached --name-only`：暂存清单**只出现我自己的 5 枚路径**（`frontend/` 下），零他人路径；`git diff` 删除行逐枚看过，全是有意替换。
- **未跑** `npm ci`/`build`/`render:*`/`go test`/`d22scan`/Docker（编排者测量程仍在跑）。⇒ 因此 §26 那格"盘上 0 命中"是**字符级 grep 现量**，**不是 d22scan 的 rc**；真绿要等闸门解除后复跑，且**编排者的 `U1` 那批落地前，`internal/panel` 那把尺仍会报 0**（两把尺继续互斥）。
- **待闸门解除必须补的一条复算**（已写进 commit 正文）：`npm run render:composer && git diff --exit-code -- frontend/fixtures/composer-states.html` —— 那 3 枚 fixture 字符是我手工改成"与渲染器会输出的一致"的，没跑过渲染器就不算证。
- 本轮**未**开任何窗（P5 约定 15:00）；`frontend/embed.go` 未动、无计划动。

---

# 追加六（2026-09-25 07:3x）：我把刚绿的门又弄红了，已结 ＋ 一条会反复咬人的规矩

## 31. 事故与归因（我自己造的，不推给别人）

编排者现量 `git archive HEAD` 净快照跑 `d22scan` → **rc=1**，报 `frontend/VENDORED.md:149 [emoji]`。
那枚字符是 **`⚠` U+26A0**，落在 `2600–27BF`，**是我 07:23 那枚 `67192a5` 写说明文字时带进去的**。

时间线（他的，我按同一形状自核过）：`07:17` 净快照 rc=0（HEAD=`471af50`）→ `07:23` 我提 `67192a5` → 红。
**不是我上一轮那 6 行残留**（那是 U+2264/U+2212，`3b59512` 已清），**也不是他扩的 `2200–22FF`**（U+26A0 在扩之前就扫）。

**修复＝`0cbb7c1`**：只删 `:149` 行首那枚 `⚠`，让 `>` 引用块自己承担语气；没替换成 `[警告]` 之类，避免为同一件事再造第二种记号。numstat `1/1`，diff 只有那一行。

**改后我做了全量自查而不是只核他点名那一行**：按 `main.go:124` 的字符类逐字扫 `git ls-files frontend` 的 **40 枚 tracked 文件**，**命中 0**。
（上一轮我在同一段用的 `①②③` U+2460–2462 与 `⇒` U+21D2 都落在**刻意不扫的空隙**里，所以本轮确实只有这一枚因，无需连带改。）

## 32. 一条会反复咬人的规矩：注释豁免**按文件类型给**，不给 markdown

编排者今天在净快照里给了三发正控（归档 `docs/evidence/s1/141-ac34-positive-controls-r1.md`），我把它翻成"我写东西时哪一格会被扫"：

| 我写在哪 | 注释豁免生不生效 | 依据 |
|---|---|---|
| `frontend/**/*.tsx` 的 `//` 或 `/* */` 注释 | **生效**（放 `≤` 进门 rc=0） | `textCommentRanges` 认行首 `//`、`/*`、块内 `*` |
| `frontend/**/*.tsx` 的 **JSX 文本节点** | **不生效**，rc=1 报 `scope frontend/` | 同上，文本节点不是注释 ⇒ **我这条腿有分母，不是空转** |
| `frontend/fixtures/*.html` 的 `<!-- -->` | **生效**（`<!--`/`-->` 是认得的标记） | 同上 |
| **`frontend/VENDORED.md` 的任何一行** | **不生效**——`#` 与 `>` 都**不是**标记 | `main.go` 注释明写 "`#` and `--` are NOT markers" |
| `frontend/scripts/*.mjs` 的 `//` 注释 | **生效** | 同 `.tsx` |

⇒ **硬结论：`VENDORED.md` 是我唯一"写什么都算正文"的文件**，而它恰好是我记录改动说明的地方。
**往 `frontend/**` 写文案时 `⚠ ✓ ✗ ☐` 一类（`2600–27BF`）全部不可用**；`→`/`①`/`⇒` 今天可用，但**不依赖这一点写新文案**（那是仪器的空隙、不是许可，票 141 的 (b) 支正悬着）。
`docs/**` 不在射程 ⇒ 本工作记录里照常用，不用自我阉割。

## 33. 三条同源副本到这一轮的归属与结论（收口）

| 拷贝 | 谁的地界 | 结论 |
|---|---|---|
| ① `tools/d22scan/main.go:124` | 编排者 | 权威，已含 `2200–22FF` |
| ② `internal/panel/frontend_hygiene_test.go:71` | 编排者 | `3c80352` 补回该段，并把那句 blanket 的 "copied verbatim" 收窄成它真主张得了的一轴 |
| ③ `frontend/scripts/vendor.mjs:53` | **前端会话** | `67192a5` 选 **A 案补宽**，理由＝旧尺结构性看不见上游那 2 枚 U+2212（反事实读数 0 vs 2） |

**但 ③ 那条链得常绿**：`vendor.mjs` 重新 vendor 时会把上游的 `−` 带回来，**是补宽后的 `EMOJI_RE` ＋ `GLYPH_MAP` 把它当场换掉**，不是"上游已经没有它了"。所以哪天有人把 ③ 收窄回去，`3b59512` 的正文 fix 就会在下一次 re-vendor 时静默回退。

## 34. 本轮纪律自述

`0cbb7c1` 带显式 pathspec、暂存清单只有 `frontend/VENDORED.md` 一枚、未 push、未碰他人路径与 `design/**`。
本轮**未**跑 `npm`/`go`/`d22scan`/Docker、**未**开任何窗。
我自己写进台账的错累计 **11 处**（新增本条：往 `frontend/**` 写说明时用了被扫的 `⚠`）。

---

# 追加七（2026-09-25 07:4x）：视图切换层查完了——是 C（契约没规定），但它拆成两问，且**卡住代码的那问不是版面**

## 35. 结论：C。仓里没有任何一份冻结文本规定面板的导航拓扑

一枚只读程把 `SPEC-08` 全文、`PLAN.md` 的 D10／D29／§17、票 33/34/36/37/38/39/40/77/92、`issues/README`、demo 与旧原型全搜了一遍（关键词含 `侧栏 导航 tab 屏 页 route router 几栏 单栏 视图切换`）。

- **不是 A**：没有任何冻结文本写过"面板有几屏／怎么切／有没有侧栏"。`SPEC-08:139`、`SPEC-00:29`、`PLAN.md:1036` 三处只是**内容枚举**，且**数目互不相同（7／6／4）** ⇒ 那份清单从来不是版面规定。`pending-and-issues.md:2795` 那句"面板导航 = `Sidebar Nav`"**被它自己上一行标成"我的判断，不是 owner 的"**。
- **不是 B**：`SPEC-08 §9` 的"不做"清单逐字只有 6 项，**没有一项禁止侧栏／多屏／切换**。`PLAN.md:3595` 那句"不做路由"**限定在"原型阶段"**，既不构成禁止、也不能被读成批准。
- **但"形状未定"≠"我可以自定"**（`AGENTS.md` §0 未定义即停）。

## 36. 真正的阻塞不是版面，是"当前在哪一屏"这件事没有归属

`SPEC-08:150-151` 逐字：**"前端必须无状态：WebView 销毁/隐藏后状态全丢；每次 `show` Go 侧推 `panel.resync` 全量状态；前端不得缓存任何跨 show 的业务状态"**。
而 `PanelSnapshot`（`frontend/src/lib/panel.ts:129-137`）只有 `pending / results / composer / generatedAt` —— **没有任何 view／tab／screen 字段**；Go 侧 `internal/panel` 全树 grep `Tab|PanelView|screen|Screen`（排测试）**零命中**。

⇒ **"用户此刻停在第几屏"要么归 Go 推（＝要扩 C17，走切片卡定稿），要么归前端本地 state（＝和 §5.1 那句"每次 show 全量 resync"正面互斥）。仓里没有裁定。**
**这决定的是代码形状，不是 CSS 形状**——版面选 A 还是 B 我用同一套组件都能画，但"当前屏存哪"两种写法**互不兼容、不能后补**。所以这一问不结，我写出来的第一版切换层有相当概率要整块重做。

顺带一条排除项：**不引 `react-router` 之类**。`SPEC-08 §5.3`（`:176-183`）那份依赖枚举里没有任何路由库，加它＝动 D29 那张已定案清单＝人工批准。两屏以上的切换用本地 state 或 Go 推，不碰依赖面。

## 37. 照 demo 加侧栏会立刻撞的三处**文本级**互斥（我原来只量到第 1 条）

| # | 冲突 | demo | 冻结文本 | 我上一版有没有量到 |
|---|---|---|---|---|
| 1 | 面板宽度 | `styles.css:116` `width:720px`（含 60px 侧栏 ⇒ 内容区 660px） | `SPEC-08:214`／`:247`、`SPEC-03:39` 都是 **640px** | 量到（K1） |
| 2 | 图标描边 | `styles.css:213` `stroke-width:1.75`（README:77 同） | `SPEC-08:218` **1.5**、尺寸仅 14/16/18/20px | **没量到** |
| 3 | 图标名 | 导航 10 枚里 **7 枚不在 §17.4 清单**（`message-square-text`/`shield-check`/`list-checks`/`orbit`/`chart-column`/`life-buoy`/`settings`；在表内的只有 `command`/`shield-alert`/`lock`） | `PLAN.md:3455-3465`"需要的图标清单…**一次性画进** `icons.js`" | **没量到** |

**第 3 条我本轮自己复核并往前推了一步**：清单里**根本没有"对话"类也没有"成本"类图标**（`grep` 全文 `message*`/`chat`/`coin`/`dollar`/`chart*`/`receipt`/`wallet` 在 §17.4 段内**零命中**；唯一命中的 `sparkles` 是**禁用清单**里的 AI 俗套图标）。
⇒ **十枚导航目的地里，`对话` 与 `成本` 在冻结图标清单内无解。** 要么改 `PLAN.md` §17.4（＝动 PLAN，人工批准），要么这两枚退而用形状不贴切的在表图标。**这不是我能自裁的，也不是"照 demo 做"能绕过的**——demo 用的那两枚名字本身就来自清单外。

## 38. 一枚与 §32 同族、但方向相反的硬证据（写给以后所有轮的我）

`PLAN.md:3443` 的"绝对禁止"行**逐字点名 `✓ ✔ ✗ ⚠ ★ →`** —— **包含 `→`**。
⇒ 所以 §1.4 T1 那条"仪器不扫箭头段"只说明**仪器抓不到**，**不说明 spec 允许**。二者互斥这个事实本身是票 141 的 (b) 支（改规格文字＝人工批准）。
⇒ **定式升级**：往 `frontend/**` 写文案时，判据取 **spec 文字那份（更严）**，不取仪器那份。`→` 今天能过门，但按 `PLAN.md:3443` 它**本来就不该出现在界面文案里**；demo 通篇写 ASCII `->` 才是应继承的纪律（这一条我上一版写在 §1.4 T2 的旁边，本轮才拿到 PLAN 的原文支撑）。

## 39. 本轮亲自复量清单 ＋ 状态

新增实测：`SPEC-08 §9`（`:244-249`）"不做"清单逐字 6 项 · `PLAN.md:3595` "不做路由"的限定语 · `SPEC-08:150-151` 无状态那句 · `panel.ts:129-137` 无 view 字段 · `internal/panel` grep `Tab|PanelView|screen` 零命中 · `PLAN.md:3434-3466` §17.4 六组图标名全读 ＋ 对话/成本类零命中（`grep`）· `PLAN.md:3443` 禁用行含 `→`。

**任务状态**：#1 差距表已交（`2b1aefe`＋五轮追加更正）· #3 基线仍被三条件闸门闸住（本机 flake 量率程在跑）· **#2 视图切换层：从"等一句话"改成"等两句话"**，见下。
**本轮未跑** `npm`/`go`/`d22scan`/Docker，**未开任何窗**，未写 `frontend/**` 一字节，未碰 `design/**`。
