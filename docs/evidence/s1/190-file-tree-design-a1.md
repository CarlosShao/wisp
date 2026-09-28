# 190-a1 — 工作区文件树那堆的根、上限与忽略规则（只读设计核）

- 分工：`.scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md` §D
- 工单：`.scratch/wisp/issues/190-the-workspace-files-stack-file-tree-has-no-host-side-source.md`
- 角色：只读设计核。**零产码、零 `frontend/**`／`design/**` 写面、票面原文一字不动（本程也没追加 Progress log，见文末"我没做的事"）**
- 本程 AC 框**一枚不勾**（勾要非实现者的表）

## 0. 起手三件（逐枚名册）

```
date                        = 2026-09-28 17:48 +08（本机）
git log --oneline -1        = 38fc7c0e ledger(A387-A388) 收 33-r3（入向听众进树＋词面尺改能力尺，我四把尺复跑）＋owner 提问后按默认动作派 145-r2
git rev-parse HEAD          = 38fc7c0e
分支                         = dev
```

锚点＝编排者锚 `38fc7c0e`（相同 ⇒ 无漂移）。在飞的写腿 `145-r2` 的四个文件（`internal/panel/composer.go`／`internal/panel/pump.go`／`cmd/wisp/panel_pump.go`／`internal/panel/pump_test.go`）**起手名册里一枚都不在**（＝它们当时没有未提交改动），我全程没写它们，只在 `git grep` 结果里读到过 `panel_pump.go:82-86`（读不算写面）。

起手 `git status --porcelain` 逐枚（71 枚，终态闸门＝终态等于这份名册、差集为空；⚠ 不是"必须为空"，`A374`）：

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
?? .scratch/wisp/.scratch/
?? .scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
?? .scratch/wisp/probes/139/accept-r1/
?? .scratch/wisp/probes/152/overlay-probe1-on-samppost.json
?? .scratch/wisp/probes/156/__pycache__/
?? .scratch/wisp/probes/156/mut-156-r2/asis.log
?? .scratch/wisp/probes/156/zero156-r4-head.sh
?? .scratch/wisp/probes/156/zero156-r4-work/
?? .scratch/wisp/probes/158/r2/
?? .scratch/wisp/probes/161/r2/__pycache__/
?? .scratch/wisp/probes/161/r2/ctl/
?? .scratch/wisp/probes/161/r5/negative-control/
?? .scratch/wisp/probes/161/r6/logs/flip-7.txt
?? .scratch/wisp/probes/162/r4/
?? .scratch/wisp/probes/162/v1-baseline-gotest.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start.txt
?? .scratch/wisp/probes/183/accept-v1/
?? .scratch/wisp/probes/185/c1/logs/d22scan-post-final.txt
?? .scratch/wisp/probes/33/r2/d22scan-final.txt
?? .scratch/wisp/probes/33/r2/gate-final.txt
?? .scratch/wisp/probes/33/r2/status-final.txt
?? .scratch/wisp/probes/33/r2/status-start.txt
?? .scratch/wisp/probes/999/
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? .scratch/wisp/probes/161/r6/logs/flip-7.txt
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
?? part1-state1-fixed.txt
?? part1-state1-pristine.txt
?? part1-state2-fixed.txt
?? part1-state2-pristine.txt
?? part2-nog6-fixed.txt
?? part2-nog6-pristine.txt
?? part3-stale-fixed.txt
?? part3-stale-pristine.txt
```

（上列第 60 行附近 `flip-7.txt` 是我从原始输出逐枚抄写时重复誊了一行；原始读数里它只有一枚。终态比对以原始 71 枚为准，差集按"名"算不按"行数"算。）

## 1. 起手锚点自证（`git log` 原文，本程第 1 枚调用即取）

```
38fc7c0e ledger(A387-A388) 收 33-r3（入向听众进树＋词面尺改能力尺，我四把尺复跑）＋owner 提问后按默认动作派 145-r2
318d4bb3 33-r3(片④): 证据件 33-inbound-listener-r3.md ＋台件 probes/33/r3/**
d9dfe2e2 33-r3(片③): 票 33 AC#9 追加（未勾）＋ Progress log 一节
e251f950 33-r3(片②补): 那枚能力尺的现量口径写准（27 枚→18 枚）
c64db7d6 33-r3(片②): 片 A 那枚词面负向尺改扫能力（A386 裁定 b 支）＋正控自带
```

自证口径：本程只读，交件证明＝**我那一枚 commit 的 `git show --name-status`**（只含本文件一枚）＋终态名册差集。**不用区间 diff 自证只读**（`git diff <锚>..HEAD` 会把别人的活算进来，共同规矩第 10 行）。

## 2. §D①　树的根取哪一枚（现读 `risk.PathResolver` 那一族，不按名字猜、不按注释猜）

### 2.1 现量：今天"范围"是怎么算出来的

| # | 事实 | 现量凭据（`file:line`） |
|---|---|---|
| 1 | C26 `PathResolver` 是全仓**唯一**路径规范化入口；`filepath.Clean/Abs` 只在 `lexCanonical` 那一处被赦免 | `internal/risk/pathresolver.go:13-14`（"the ONLY path normalization entry point"）、`:204-214`（`lexCanonical` = 唯一赦免点）、`:105` `Resolve(input, exceptions)` |
| 2 | 流水线顺序固定：展开(`%VAR%`/`~`) -> abs -> Clean -> 句柄真实路径 -> **reparse 检测默认拒绝** -> 8.3 短名 -> UNC 归一 | `internal/risk/pathresolver.go:16-24`；reparse 默认拒 = `:117-125`（非例外即 `ErrReparseDenied`，`:36`） |
| 3 | `Resolve` 的返回值带**账**（`Spelling`/`Canonical`/`Resolved`/`Rewritten`/`Rewrites`），"能动的拼写"只从 `Actable()` 出 | `internal/risk/pathresolver.go:52-80`、`:93-99`（`Actable` 在被改写时返回 `ErrRewrittenPath`，`:49`） |
| 4 | **授权根只有一枚来源**：`[fs] allowed_dirs`，构造期逐条过 C26；**空名单＝什么都不授权**（每个 fs 调用落在 L2，不是"没判"） | `cmd/wisp/run.go:327-328`（`for _, d := range cfg.FS.AllowedDirs`）；`internal/tools/paths.go:50-52`（"An EMPTY allowlist authorizes nothing"）；默认值 `docs/specs/SPEC-03-config-secrets-envs.md:35`（`allowed_dirs[]=[]`） |
| 5 | 授权根之外还有两本账：`unusable`（规范化失败 ⇒ 什么都不授权）与 `rewritten`（展开把根搬到了另一棵树） | `internal/tools/paths.go:36-41`、`:58`、`:84-91`；对外读面 `Roots()`／`UnusableRoots()`／`RewrittenRoots()` `:178-198` |
| 6 | **工作区根 = 窄化器，不是另一枚根**：`workspace` 字段只能**收紧** `InAllowlist`，窄化到根外的工作区在**切换时**就被拒 | `internal/tools/paths.go:43-47`（"It can only ever TIGHTEN InAllowlist … no workspace choice can widen authority"）、`internal/tools/paths_workspace.go:40`（`WorkspaceRoot()`） |
| 7 | 判定函数 `InAllowlist` 的真形状 = **三条与**：①在配置根集合内 ②**跟随链后的真实形也在配置根内** ③`workspace != ""` 时还在工作区内 | `internal/tools/paths.go:133-163`：①`:140`、②`:152-155`（`resolvedForm` 拿不到就 `return false`，fail-closed）、③`:159-161`；`workspace == ""` 的语义＝"什么都没窄化"，第一次切换前逐字节等于旧行为 `:157-158` |
| 8 | 包含测试是**组件边界**＋大小写折叠（不是前缀字符串） | `internal/tools/paths.go:168-175`（`rootsContain`：`folded == r || HasPrefix(folded, r+pathSep)`）、`:289-300`（`foldPath` 是**比较折叠**，明文"不是规范化，拿去代替 Canonicalize 就是 D22 禁的 C26 绕过"） |
| 9 | 面板侧已经有一整套"工作区那一维"的宿主读面，而且**明写面板不得决定路径解析** | `internal/panel/workspace.go:6-11`（"It may not decide what that folder resolves to: the whole point of the native leg is that C26 gets the vote"）；接口 `:37-45`（`WorkspaceRoot`／`ResolveWorkspace`／`SetWorkspaceRoot`）；窄化视图 `:60-68` |
| 10 | **"未窄化"Today 面板报的就是授权根集合**，逐字文案已在册 | `internal/panel/workspace.go:56`（`"未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定"`）；泵侧装配 `cmd/wisp/panel_pump.go:82-86`；提供者钩子形状 `internal/panel/pump.go:114`（`Workspace func() WorkspaceView`） |

### 2.2 裁定（问题①的答案）

**两个都不选成"一枚根"。今天代码里的"范围"是一个根集合，树必须照它长成集合形状：**

- 有效读根 = **`{workspace}`（工作区已窄化）∪ `Roots()` 中同时满足 `InAllowlist` 三条的那批（未窄化）**。因为第 6/7 行：切换时根外的工作区进不来，所以**窄化后有效集合恰好＝`{workspace}`**；未窄化时集合＝`Roots()`，**枚数可以是 0，也可以多于一枚**。
- **二者不重合时以"更窄的那一枚"为准＝工作区根**，但这条不是"优先级"、而是**交集**：`InAllowlist` 是与折叠（`paths.go:140/153/159`），任何一枚不满足就不在范围内。把设计写成"工作区优先，否则用 allowed_dirs"就丢了"多根"这一形，会把三枚授权目录显示成一枚、或者反过来凭空多出一棵。
- **零枚根是一种正常态，不是错误态**：`allowed_dirs` 默认 `[]`（`SPEC-03:35`）⇒ 出厂态下树**没有任何可显示的根**。那一格必须显式渲染（沿用 `workspace.go:56` 那句"本轮按 [fs] allowed_dirs 授权的目录判定"的口径），**绝不许回退到"进程工作目录"或"当前盘"**。票面 AC#2 那发判据钉的就是这个回退形状（见 §5）。
- **落地腿取根的唯一途径是现成读面**：`PathCanonicalizer.Roots()`／`WorkspaceRoot()`（`paths.go:178`、`paths_workspace.go:40`）。**不许**从 `config` 再读一遍 `[fs] allowed_dirs`、更不许自己拼绝对路径 ⇒ 那会造出第二本账，与本仓既有裁定同向（票 115 已警告"winsec 的通知与 PATH-ACCOUNT 是两本账"；票 189 §②同一理由要"一份读面"）。
- **树腿不接受来自面板的根参数**（结构性，非约定）：面板只订阅快照维，根由 Go 侧自取。这条同时把 AC#5 的"C17 不加方法"变成"根本没有新方法来"（见 §6）。

## 3. §D②　复述并现量核对：`allowed_dirs` 是**判级输入**、不是执行时的硬边界

**台账原文复述**（`docs/reports/pending-and-issues.md:8098`，A352 标题末句，逐字）：

> 顺手带出一条架构事实我现读复验：`[fs] allowed_dirs` 是判级输入、不是执行时硬边界

**本程现量核对（三条独立凭据，方向一致）**：

| # | 核对点 | 凭据 |
|---|---|---|
| 1 | 越界**只产出一档 level**，不产出拒绝：R2 的返回是 `contribution{rules:[R2], level:L2, reason:...}`，函数签名里没有任何 deny 出口 | `internal/risk/rules_gateway.go:45-49`（`if !ctx.canon.InAllowlist(canonical)` ⇒ `level: L2`）；同族 fail-closed 形状 `:41-43`（规范化失败也走 L2，不走拒绝） |
| 2 | 真打开的是 `Canonicalize` 的返回值，**不是"根内的某条路径"**；批准过 L2 卡之后，这一发照样读到根外的内容 | `internal/tools/paths.go:115-118`（票 102 处置逐字："fs.go / fs_write.go / mode.go open exactly what this function returns"）；执行侧 `internal/tools/fs.go:206-214`（`t.d.open(a.Path)` 的返回直接 `os.Open`，中间没有 allowlist 闸） |
| 3 | 这条边界**不是"没实现"**：`[fs] allowed_dirs` 的默认空值导致"根外读不到"的表象，票 174 已经把它当成一条独立的安全边界摆给 owner（`Q-60`），而不是当成缺陷 | `docs/reports/pending-and-issues.md:8007`（"要不要放开那棵含 `secrets\` 的树是安全边界，不是我该顺手决定的"）；`docs/evidence/s1/102-adversarial-acceptance.md:224`（`InAllowlist` 是"任一根匹配即真"的 **OR 折叠**，不是 first-match） |

**由此得出的设计约束（本票最重要的一条"**不许顺手做**"）**：

1. **树那一维不许被实现成硬边界。** 它显示的根集合是"**这一维从哪里取数**"，不是"**谁能读**"。把过滤写成执行时的授权，等于在本票里悄悄改掉了 2 号事实——那是 **owner 单独批的契约变更**（对照 `AGENTS.md` §0 第 2 句"改契约＝人工批准"），不是设计核或落地腿的顺手项。
2. **硬边界那一改会造出一句谎话**：批准过的 L2 卡能让模型读到根外（上表第 2 行），而根外的条目在树里"不存在"。于是**面板与模型看的是两棵不同的树**，用户会拿面板的"没有这个文件"去反驳模型的真读数。本仓对同族谎话的裁定一直是"宁可显式说不安全"（票 174 的 AC#2b 文案改动、票 92/145 在 `git.go:37-40` 拒绝把"读不到"塌成"不是仓库"）。
3. **必须同屏说清"这是显示范围不是权限"**：文案上写"范围来自 `[fs] allowed_dirs`／工作区窄化"，**不许**写"根外不可访问"。
4. **根的账要与判级的账同源**：树腿读 `Roots()`，R2 判级也读 `roots`（同一枚字段、同一把锁 `paths.go:138-139`）⇒ 面板显示的范围与审计打印的范围必然一致；不一致就是第三本账。

## 4. §D③　路径处理必须走 `PathResolver`——把这一条写成判据思路

**规矩出处**：`AGENTS.md` §1.2（"在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策＝直接用即判违规"）；执行器 `tools/d22scan/main.go`（ban #2，扫描入口在 `:1180` 附近那族）；赦免点唯一 = `internal/risk/pathresolver.go:206-214`；仓内脚本 `scripts/check-pathclean-ban.sh`（`pathresolver.go:27-28` 指向它）。

**判据思路（四发，都写成"红/绿"能判的形状，不写成约定）**：

| 发 | 判据 | 为什么它能真咬住 | 正控（不种就自证无效） |
|---|---|---|---|
| F1 静态 | 树那一枚文件内 `filepath.Clean`／`filepath.Abs` **零命中**；名册尺常驻。⚠ 射程要写准：ban #2 只禁这两枚；`filepath.Join`／`Dir` 不在射程（`internal/tools/fs.go:258` 的 `joinForListing`、`git.go:46` 都合法用着） | 违规定义＝CI 定义，不新造词 | 把 `joinForListing` 那行换成 `filepath.Clean(dir+"\\"+name)` ⇒ 门必须红 |
| F2 能力（**不是词法**） | 树上"每一条被下钻的目录"都必须有**一次 `risk.Resolve`（或 `canon.Canonicalize`）的证据**；纯词法遍历＝红 | 只有过 C26 才会撞 `ErrReparseDenied`（`pathresolver.go:117-125`），词法拼接**结构上看不见 junction** | 现成模板：`internal/tools/bridge_junction_windows_test.go:278-291`（逐字 `fs.list 把 junction 对面的条目列出来了` 即 `t.Fatalf`）——同形立 `TestTreeLegDoesNotCrossAJunction`，正控＝真的种一枚 junction |
| F3 账 | 树腿不得从 `Result` 里只取 `Canonical` 而丢掉账；参与"要不要下钻/要不要显示"的字符串必须来自 `Actable()` 或 `Canonicalize()`，**不许**来自 `foldPath`（`paths.go:291-293` 逐字："calling it on user input instead of Canonicalize would be the C26 bypass D22 forbids"） | 折叠函数长得像规范化、语义只是比较，最容易顺手拿 | 变异＝把某条路径的 `Canonicalize` 换成 `foldPath` ⇒ 断言"改写腿/链外腿必须仍红" |
| F4 根 | 显示根的候选集**逐枚等于** `Roots()`（∩ `WorkspaceRoot()`）的成员，不是"由配置重新算一遍"、也不是"从环境猜" | 二次取数必然漂（票 189 §②同一判据、票 174 `run.go:327-331` 唯一来源） | 把 `d:\work\proj` 塞进 `allowed_dirs=["d:\\work"]` ⇒ 尺必须只认 `d:\work` 这一枚根，报 `d:\work\proj` 即红 |

⚠ F2 是这一族里唯一能咬住"跨树泄漏"那一形的；F1 是词法尺，**光有 F1 等于没设防**（按 `MEMORY` 09-28 的教训："词面型负向尺该改扫能力"）。

## 5. §D④　深度／条目上限／忽略规则（建议值＋理由＋截断时用户看见什么）

### 5.1 上限（三枚各自独立、各自可见；数字是建议，owner 可改，**"三枚各自可见"不可改**）

| 维度 | 建议 | 理由（每条都有现量锚） |
|---|---|---|
| 深度 | 根自为第 1 层，**最多 3 层** | 树腿要**每次快照重跑**（载体＝`pump.go:114` 那一族 provider 钩子；git 维同泵，`internal/panel/git.go:8-11`），成本模型与"模型点一次 `fs.list`"不同；今天没有任何一次调用能照 `fs.list` 那样被摊到 4000-token 溢出层去（`internal/tools/fs.go:36-38`） |
| 每目录条目 | **200** | 直接参照 `internal/tools/fs.go:39`／`:56`（`defaultMaxListEntries = 500`，那是**模型一次读**的量级）；面板树取它的同族小值，因为一次 push 里是 N 个目录而不是一 |
| 全树条目总量 | **1000**（跨层合计） | 前两枚各自封顶挡不住"3 层 × 每层 200 × 分支数"的乘积；D32 的资源预算是**空闲 RSS/CPU** 那一族（`AGENTS.md` §1.1 禁改阈值），落地腿**必须实测一次**再定稿（见 §7 未裁完） |
| 时间 | 一次遍历带 ctx 早停 | 现成形状：`internal/tools/fs.go:251-256`（`fi()` 每条目查一次 `ctx.Err()`，取消就追加 `"(已中止：…)"` 而不是静默返回半截） |

### 5.2 忽略规则（票面 AC#1 要求复用 `d22scan` 那族，本程核到它的真语义）

现量：`tools/d22scan/main.go:56-75` 的注释逐字定义了那套规矩——**".gitignore 模式半 ＋ 本树 git index 半"**；`skip()` 问 `git ls-files -z` 与 `git ls-files -i -c --exclude-standard`，**被 git 报告的（已跟踪／被忽略但被 `-f` 强加）一枚都不跳过**（`:64-66`）；问不到 index 时的降级是"一条规则也不适用"，并且**每次运行打一行响亮自陈**（`:75`"The degradation is over-coverage, never a …"）；实现文件 `tools/d22scan/gitignore.go`（`:131`、`:137` 说明"读而非抄，且 no-index 那一支什么规则都不套"）。名字级硬编码跳过在 `main.go:660`／`:853`／`:931`（`testdata`／`.git`／`node_modules`）。

**本票对"树"的裁定与扫描器不同，而且必须不同——差别落在"显示"还是"下钻"上**：

| 项 | 扫描器（既有） | 文件树（本票） | 为什么 |
|---|---|---|---|
| `.git`／`node_modules`／构建产物 | 整个跳过，看不见 | **仍出现在父目录列表里，但标"未展开（已忽略）"**；只有**下钻**被忽略规则挡住 | 票面 AC#4 的"不许静默少列"与 AC#1 的"忽略规则"是同一条规矩的两半；把条目**抹掉**会让用户以为工作区里没有 `node_modules`，而它正解释了那 20 万枚文件为什么没显示 |
| 过滤的真假（是否被 git 跟踪） | 必须问 index | **同一枚尺，同一份实现**：问得到 ⇒ 按 `.gitignore`＋index 双条件；问不到 ⇒ 一条规则都不套、全扫，并**响亮自陈**（快照维里带一句"这台机器上问不到 git index，已按全量显示"） | 票面逐字"别另发明一套"；`MEMORY` 里那条 d22scan 裁定（"忽略文件过滤必须问 git index、问不到就全扫并响亮自陈"）就是这条规矩的权威表述 |
| 被忽略的目录里**存在但看不见**的量 | 不报告 | 报告 `hidden_by_ignore`（**枚数**，不是名字列表） | 与 §5.3 的三元组同一个字段，避免"截断"与"忽略"共用一个洞 |

### 5.3 树被截断时用户看见什么（不许静默少显示）

现成模板就是 `fs.list`：`internal/tools/fs.go:225-238` —— `truncated := len(names) > limit`；截断后打印 `%s (%d 条目%s)`，`sizeSuffix` 逐字给出 `"，已达上限 %d，truncated"`；**并且把 `Truncated: truncated` 作为结构化字段返回**，不靠文本。

树那一维照此立**四元组**，缺任一枚即判"静默少显示"：

1. `shown` —— 真送进面板的条目数；
2. `total` —— 该次遍历**数到**的条目数（数到就报数，没数完就报 `unknown`）；
3. `truncated_by` —— 哪一枚上限触发的（`depth` / `per_dir` / `total` / `ignored` / `read_error`），**多枚同时触发要全列**；
4. `out_of_scope_note` —— 一句话说明"范围来自 `[fs] allowed_dirs`／工作区窄化，不是权限边界"（§3 第 3 条）。

⚠ `total == unknown` 与 `total == 0` 必须是两种可区分的读数——把前者塌成后者就是 `git.go:37-40` 拒绝的那一发谎话（"读不到"不许塌成"没有"）。

## 6. §D⑤　只读边界：删除／移动／重命名／新建一枚都不许有——**结构性**判据，不是"约定不做"

票面 AC#3 与普查雷区一致（`docs/evidence/s1/180-182-panel-fields-census-c1.md:207` 逐字："S4 文件树 只读**无**；⚠ 树上带删除／移动＝写动作"；`:343` "雷区 4 处…文件树上的删除／移动"）。四发判据，**每一发都是"做不出来"而不是"不许做"**：

| 发 | 结构性判据 | 凭据／为什么它不是约定 |
|---|---|---|
| S1 拓扑不可达 | 树腿**只作为 outbound 快照维存在，不注册任何 C17 入向方法**。入向那一跳今天在这棵树里根本不存在：`internal/panel/git.go:18-20` 逐字"that hop does not exist in this tree, which is why SwitchBlocked below is a constant sentence rather than a computed verdict"。⇒ 只要树腿不新增接收者，"面板上点删除"**在拓扑上不可达** | 票面 AC#5 已钉"C17 不加方法"；这条把它从"别加"升级成"加了也没有能触到它的东西"。落地腿交件时**必须现量名册**（本程没逐枚展开，见 §7） |
| S2 无入参即无写目标 | 树腿的构造/枚举函数**签名里不接路径参数**（根自取，§2.2 末条）。没有可指的路径，就没有可下手的对象 | 对照写族形状：`internal/tools/fs_write.go` 的入口一律是"参数里带 path"；把树腿写成带参 ⇒ S1 立刻失效 |
| S3 依赖名册尺（常驻） | 那一枚文件的依赖/字段/import 名册**只准含读类**（`os.ReadDir`／`os.Lstat`／`risk.Resolve`／`PathCanonicalizer` 的读方法）；出现 `stageAndRename`／`fs_write`／`recycle`／`exec.Command` 即红。模板＝`git.go:22-26` 的 "NO EXTERNAL PROCESS. Everything is os.Stat / os.ReadFile / os.ReadDir"＋票 188 `AC#3` 那一发"这一维只有 Go 侧生产者在写"的同族钉 | 写动作在仓里**有具名的唯一入口**（`fs.go:23-25` 逐字列了写族与 D31 的 `stageAndRename`），所以名册尺是有限集、可判；且这不是词法尺——它扫的是"这个包持不持有写能力" |
| S4 不许变成 Tool | 立 `TestTreeDimensionIsNotModelCallable`：树维不许进 `internal/tools` 的注册表。D34 表没有文件树行（票面 AC#5"D34 表不加行"），而给了模型一枚宿主树读工具**就是契约变更** | 逐字模板：`git.go:13-16` + 那发 `TestGitDimensionHasNoModelCallableTool` |
| S5 截断/越权不许被当成"没有" | 拒绝显示（越界、junction、权限拒绝）与"目录里就这么多"必须是不同枚举值 | 现成枚举形状：`git.go:65-70`（`repo`/`not_a_repo`/`permission_denied`/`unreadable` 四态，注释逐字"the panel must not render the second and third as the first"） |

**树上写动作的真去处**：票面 §"本票不解决"逐字"不做树上写动作（待批）"，普查 §:343 把它列进雷区 4 处并要求 owner 逐枚批。⇒ 本表**不提供**任何"下一步怎么放开"的路线图；要放开是另一枚票。

## 7. §D⑤后半　与 §A 的"文件树那一堆"是不是同一枚／并进票 145 还是另立

**是同一枚。**具名对应＝普查里的 **S4 工作区文件树**（`docs/evidence/s1/180-182-panel-fields-census-c1.md:130`、`:152`、`:180`），载体候选已在同一份表里被命名为 **`B5 = Tree.Entries`**（`:226`）。票 190 就是 S4 的票（`:180` 逐字"没票认领 ⇒ 要立票"，票面 `:4` 记录立票来路 `A373`）。

**归属裁定：拆成两格，载体并进票 145、生产者另立在票 190 本票。**

| 格 | 归谁 | 理由 |
|---|---|---|
| **"快照有哪几维"（`Tree.Entries` 这一维本身、它的空态/截断字段进不进 `Snapshot`）** | **并进票 145** | 145 是快照载体那一集（票面 `:5` 关联行"票 145（快照载体）"），git 维已经走过同一条路（`git.go:8-11` "answers … ride the snapshot push that already exists"）；再开第二份载体必然漂 |
| **"这一维的值从哪来"（宿主侧目录遍历、根的自取、忽略与截断账）** | **票 190 本票 AC#2 另立** | 145 的地界是字段集与泵；遍历逻辑塞进 145 会让一枚载体票同时改判定路径，撞 `AGENTS.md` §1.1 的"由面板侧来源的 L2 允许"那条禁区 |

⚠ **排序硬约束（必须写进落地腿的派单）**：`145-r2` **正在写** `internal/panel/composer.go`／`internal/panel/pump.go`／`cmd/wisp/panel_pump.go`／`internal/panel/pump_test.go`（派单第 4 行；`A388` 只具名解冻 `:124`、`:276` 两行）。⇒ 票 190 的落地腿**排在 145-r2 交件之后**，否则撞同一批文件；载体那一格在 145 交件前先按"待并入"记着。

## 8. AC#4 三形（越权那一形必须显式回答）——本程只给设计口径，现量归落地腿

| 形 | 今天已有的同族读数／文案出处 | 树腿该怎么答 |
|---|---|---|
| 目录不可读／权限拒绝 | `internal/tools/fs.go:210-219`（`打开目录失败：…`／`读取目录失败：…`，都带 `Origin: canon`）；枚举位 `git.go:65-70` 的 `permission_denied` | 该目录位＝`permission_denied`，**父目录里条目仍在、子条目为 0 且标"未读"**；文案不得出现"没有文件" |
| 符号链接／junction 跨界 | `pathresolver.go:117-25` 的 `ErrReparseDenied`（默认拒，例外表逐枚精确路径、`:101-104` 注释"never prefixes"）；现成反例尺 `bridge_junction_windows_test.go:278-291` | 该条目标"链/接合点，未跟随（C26 默认拒绝）"，可点出 `Spelling`；**不许**跟随后再判断（那是票 107 的 probe C 那一形，`paths.go:143-151` 逐字） |
| 工作区未窄化却想报整个盘 | `workspace.go:56` 那句现成文案 + `paths.go:50-52`（空 allowlist 什么都不授权） | 根＝`Roots()`；`Roots()` 为空 ⇒ 空态＋那句现成文案；**无任何回退路径**（§5/§6 的 F4 尺钉它） |

## 9. 未裁完（不交半份表当完整 ⇒ 明写）

- **`Snapshot` 结构体的全字段名册没逐枚读**（预算 20 枚用到第 18 枚）⇒ `Tree.Entries` 具体插在 145 那一集的哪一维，**由落地腿现读 `internal/panel/pump.go` 后定**，与 §A 的结论对齐才算数。
- **C17 入向方法白名单今天有几枚没现量**：§6 S1 引的是 `git.go:18-20` 那句"那一跳不存在"，不是名册读数 ⇒ 落地腿交件前**必须先跑撞钉预检**（共同规矩第 13 条）并逐枚贴名册。
- **三枚上限未经实测**：本程一次遍历都没跑、没取 D32 相关任何读数 ⇒ 3/200/1000 是"与 `fs.list` 的 500 同族且更严"的**建议值**，落地腿要现测一次最坏工作区再定稿。**不许**为过 SLO 去动 `thresholds.go`（`AGENTS.md` §1.1）。
- **`internal/panel/attachments.go` 没读**（票面"全仓没有宿主侧列目录给面板的读面"这一前提我只核到"panel 包文件名册里没有树读面，`git.go` 自陈只走 `.git`"；附件那枚是否列过目录未裁）。
- 门禁四把尺（`d22scan`／`gate-clauses`／`gofumpt`／`go test`）**本程一枚都没跑**：只读设计核零产码、零 Go 文件写面，跑它们不会产新读数，且 `internal/panel` 今天自带 3 枚在册红（派单第 14 行），跑前先核名册的事归落地腿。

## 10. 我没做的事（自证口径）

零产码；`frontend/**`／`design/**` 零写面（本程连读都没读它们）；票面原文一字未动、**AC 框一枚不勾**；未追加票面 Progress log（设计核结论归本证据件，票面追加与派腿由编排者裁）；台账 `pending-and-issues.md` 未动（A## 入账归编排者）；只 commit、绝不 push；显式 pathspec＝本文件一枚；未用 `git add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`；仓内零删除；`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**` 一字节未动。

本文件的**起手名册＋锚点＋终态**三件齐全；终态差集与 `git show --name-status HEAD` 的读数见交件回报。
