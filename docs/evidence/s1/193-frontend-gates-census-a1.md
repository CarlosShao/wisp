# 票 193 · `193-a1` 只读普查表 —— "卡前端的门禁"逐枚点名（A1 腿）

> 性质：**只读普查**。本表不勾任何 AC 框、不动 `internal/**`／`cmd/**`／`frontend/**`／`design/**`／`go.mod`／`docs/PLAN.md`／`docs/specs/**` 一字。
> 一切数字均为本程现跑；抄自别家表的读数一律标〔转述，未复量〕；本程**未打开 `frontend/**` 与 `design/**` 的任何一枚文件内容**（只有计数，以及 Go 侧测试自己打印出来的失败句——那些句里的前端行号是仪器打的，不是本程摘的）。

## ① 起手锚＋写面闸门

- 派单锚点：`398e3344`｜本程**起手** HEAD 现量：**`72078779`** ⇒ 已漂；按派单 §3 首行处置（照实际 HEAD 做、登记、不停下等）。
- 共享工作树的实时性（现量）：本程跑到一半，`181-r1` 落了 `caebcbfe`，本程的骨架落在 `420dda62`，终态取数时 HEAD 已是 **`a10832ef`**。**这三枚都不是本程写的产码**，本程只提交了 §⑨ 列出的那两片文档。
- `git status --porcelain -- internal/ cmd/ go.mod` 起手现量（两行，**全是别人的在飞件，本程不动不提交不评论**）：
  - ` M internal/panel/git.go`
  - ` M internal/panel/git_test.go`
  - 中间一次复量＝**空**（`181-r1` 把 `git.go` 提交进 `caebcbfe` 了）；终态复量＝` M internal/panel/git_test.go`（同一枚在飞件的新版本）。
- `git status --porcelain -- design/` 起手现量：**16 删除／4 修改／11 未跟踪**（与派单 §2-Q1 预警同形；本程只报计数，未读其中任何内容）。
- 全仓脏件现量：**74 行**（含 `.scratch/wisp/probes/152/**`、`probes/161/r6/logs/flip-*`、`docs/evidence/s1/152-*.md` 等别家的件）。

## ② Q1 今天真的会点前端红的门禁——逐枚

### ②.1 `scripts/d22scan.sh` 名册（起手现量＝终态现量，两次逐字相同）

```
runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0
d22scan: examined 232 production Go files under internal/ and cmd/
d22scan: scope bans #1-5 internal/      examined 209 production Go files
d22scan: scope bans #1-5 cmd/           examined  23 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  21 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 436 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  45 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations            （退码 0，两次都是）
```

- 分母读数（现跑，非自造）：**`frontend/` 在 ban #6 与 ban #8 各 85 枚**；**`design/` 在 ban #8＝39 枚**（`design/` 的 39 是"当前工作树里还剩几枚文本件"，与 §① 的 16 枚未 staged 删除同时成立——删除件不在树里就不计数）。
- **`ban #8 internal/` 起手＝436，不是基线 433** ⇒ 按派单 §3 具名登记：**不当漂移、不去动它**。旁证（现跑）：`git ls-files --others --exclude-standard -- internal/ cmd/` ＝ **0 行**（无未跟踪 `.go`），`find internal cmd -name '*.go' | wc -l` ＝ **481 ＝ 436＋45** ⇒ 436 是这棵树今天的真实枚数，不是脏件造成的读数污染（与 `33-a1` 交的 436〔转述，未复量〕同值，但我这两把尺是自己跑的）。
- 扫描器自身正向对照：`PASS=34 FAIL=0 SKIP=0`，两次一致。

### ②.2 `tools/d22scan/main.go` 逐枚

| 门禁 | `file:line` | 今天实际扫哪些目录 | 卡不卡前端视觉/动画 |
|---|---|---|---|
| ban #6 面板批准面 scope | 走树 `tools/d22scan/main.go:249`；名册登记 `:395-410`（label `ban #6 frontend/`，live） | 只有 `frontend/`，今天 85 枚文本件 | **不卡动画**。它只匹配一个标识符（`approval.decide`），防的是"面板成为批准/执行入口"（D33/F2） |
| ban #8 零 emoji | 字符类 `tools/d22scan/main.go:164`；射程清单 `:557-564`（`design/`、`frontend/`(everyFile)、`internal/`(goOnly)、`cmd/`(goOnly)） | 今天 `design/` 39、`frontend/` 85、`internal/` 436、`cmd/` 45 | **会点前端红**（写进面板文案/图标的字符类命中即红），但它**不是性能门**，是 owner 自己定的视觉规范 ⇒ 见 §⑥ 硬约束① |
| ban #1-5／#7 | `:388-412` | `internal/`(209)／`cmd/`(23)／`internal/tools/`(21) | **射程里根本没有前端**，本单一枚不碰 |

`tools/d22scan/scan_test.go` 里凡带 `frontend/` 的断言几乎全是**夹具**（`seedFile(t, root, ...)` 写临时根），唯一碰真树的是被 `scripts/d22scan.sh:20-28` 点名的 `TestScannerSelfScanOfRealRepoIsGreen`——今天它 PASS（包含在上面那 34 里）。

### ②.3 全仓"读 `frontend/` 或 `design/` 内容并据此判红"的 `*_test.go`——一把不漏

尺（现跑）：`frontend/|design/` 在 `**/*_test.go` 里命中 **9 枚文件**；逐枚现读后，**真正会读前端/设计内容并判红的测试共 11 枚**，全在 `internal/panel` 与 `internal/ball`：

| # | 测试 | `file:line` | 它今天断言什么（射程） | 卡的是动画/视觉自由度，还是别的东西 |
|---|---|---|---|---|
| T1 | `TestPanelFrontendIsStateless` | `internal/panel/frontend_hygiene_test.go:169`（清单 `:48-56`） | `frontend/src/**.{ts,tsx,css}` 里出现 `localStorage`/`sessionStorage`/`indexedDB`/`document.cookie`/`caches.*`/`navigator.serviceWorker`/`node:fs` 即红 | **别的东西**（状态持有者唯一性，PLAN.md:1044） |
| T2 | `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | `internal/panel/frontend_hygiene_test.go:196`；字面量尺 `:61`；闭包尺 `:101`（`main.tsx` 的 import 闭包，今天 `57` 枚可达文件由失败句打印） | **闭包内任何非 `tokens.generated.css` 文件写十六进制/`rgb(a)` 字面量即红**，要改色必须回 `design/assets/tokens.css` 的 C21 token | **就是这枚卡视觉自由度**——它把"能用什么颜色/风格"钉在单一色源上。今天它红（仪器点名 `frontend/src/components/harness/right-rail.tsx:89` 一枚色值字面量，值不复述） |
| T3 | `TestVendoredDemoComponentsAreNotMounted` | `internal/panel/frontend_hygiene_test.go:222` | `frontend/src/components/ai-native/` 里**至少一枚** vendored 组件必须被 `main.tsx` 闭包挂载（`:240` 那句红话写的是"this library is the base"的前提） | 卡的是**组件基座前提**（票 77 的 D29 落点），不是性能；今天绿 |
| T4 | `TestFrontendHasNoEmoji` | `internal/panel/frontend_hygiene_test.go:246`（字符类 `:71` 逐字抄扫描器 `:164`） | `frontend/src`＋`frontend/scripts`＋`frontend/index.html` 里出现 emoji 字符类即红 | 卡前端**文案/图标**，非性能；注释 `:64-71` 明写这枚 copy **没有注释豁免、树比 gate 窄，"can be stricter than the gate, never wider"** ⇒ 它是比 ban #8 更严的那一份 |
| T5 | `TestFrontendNeverNamesAnApprovalDecision` | `internal/panel/frontend_hygiene_test.go:281`（尺 `:73`） | ban #6 的**内侧孪生**：整棵 `frontend/`（排 `node_modules/`、`dist/`）不得出现 `approval.decide` | 安全边界，非视觉 |
| T6 | `TestComposerContractTypesMatchFrontend` | `internal/panel/composer_test.go:48`（读 `frontend/src/lib/panel.ts` 于 `:52`） | Go 快照每枚 JSON 键必须在 TS 接口里有声明（双侧同步钉） | **别的东西**（契约同步）。今天红＝`Go ComposerState emits [git] that interface ComposerState does not declare` ⇒ 因是 `181-r1` 刚落的 git 字段，**不是门禁太严** |
| T7 | `TestTheRendererHoldsExactlyOneDoorToTheHost` | `internal/panel/composer_test.go:502`（尺 `:381`，钉死句 `:510`"host call sites in frontend/src = 2, want 2"） | `frontend/src` 里 `.postMessage(` 调用点必须恰好 2 枚且都在 `src/lib/panel.ts` | **结构/安全门**（C17 单门）。⚠ 对动画开关的实际含义：**新字段走快照不算新调用点，不阻塞**；另建一条桥腿才会红 |
| T8 | `TestNoGitSwitchCapabilityInThePanelSurface` | `internal/panel/composer_test.go:276`（红话 `:321`，日志 `:324`） | `frontend/`＋`internal/panel` 生产件里不得出现 git 分支/仓库切换入口（票 92 AC#7，owner 亲自砍掉的产品面） | 产品边界，非视觉 |
| T9 | `TestPlantedComposerModeWriteGoesRed`／`TestModeViewHasNoWriteSurface` | `internal/panel/composer_test.go:95`／`:171` | 面板不得写模式（L2/模式门） | 安全边界 |
| T10 | bridge／approval 的同形双侧钉 | `internal/panel/bridge_test.go:135`、`internal/panel/approval_test.go:109`＋`:126` | 读 `frontend/src/lib/panel.ts` 断言接口没被单边改名 | 契约同步 |
| T11 | `TestC21DesignTokensFourWayAgree` | `internal/panel/tokens_fourway_test.go:10-13`（P1 `design/assets/tokens.css` `:50`、P4 `frontend/src/styles/tokens.generated.css` `:53`；失败点 `:441`） | C21 令牌四向一致（Go 表／CSS／生成主题…） | **owner 09-25 亲自按住的那枚。本单只报它是什么，不提议顺手改**（AC#0） |
| T12 | `TestC21TableColourRowsMatchTokensCSS` | `internal/ball/tokens_table_test.go:1364`（`c21TokensCSSPath = "design/assets/tokens.css"`），失败点 `:1468` | **`internal/ball` 也在运行时读 `design/`** 并据此判红，红话自陈"the CSS leg of this check must never skip" | 与 T11 同族（令牌权威链），**本程同样不提议动它** |

`internal/panel/l2_grant_boundary_test.go` 命中 `frontend/` 的行都在**注释与 `git show 53a1359^:…` 的历史取件**里（`:19`、`:25`、`:1252`、`:2085`），并且 `:2017` 自陈"Nothing here writes into frontend/, design/"——它是否另有运行时读前端内容的断言，本程**未逐枚读完**（见 §⑦）。

### ②.4 "这条门禁存在"与"它今天为什么红"是两件事（派单 §2-Q1 点名要分开报）

| 今天红的测试 | 红的**因**（现量） | 是不是"门禁太严" |
|---|---|---|
| T6 `TestComposerContractTypesMatchFrontend` | `181-r1` 在 `caebcbfe` 给 Go `ComposerState` 加了 `git` 字段，TS 侧尚未声明（该包**起手**跑：`go test -count=1 ./internal/panel/`） | **不是**。是双侧同步钉正常起作用 |
| T2 `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | 仪器点名 `frontend/src/components/harness/right-rail.tsx:89` 的一枚色值字面量 | **这才是射程之争**：它确实把一枚新颜色挡在门外（要它绿要么前端回 token，要么放开这枚门） |
| T11 `TestC21DesignTokensFourWayAgree` | `git status --porcelain -- design/` 现量 16 枚未 staged 删除里有 `design/assets/tokens.css`；失败句是 `open …design\assets\tokens.css: The system cannot find the path specified` | **不是**。文件不见了，与阈值无关 |
| T12（`internal/ball`） | 同一枚缺失文件；`go test -count=1 ./internal/ball/` 起手现量 **1 枚 FAIL**（终态未复跑，见 §⑦） | **不是**，同一因 |

⇒ **今天"会点前端红"的门禁共 2 类 4 枚**：真正因**视觉自由度**而点红的只有 **T2**（加硬约束下的 T4／ban #8 属"文案字符类"，今天没红但会点）；其余 3 枚的红因是**别家会话的文件缺失/单边改动**。把"让它绿"当成"放开它"的理由，本表一律不采纳。

## ③ Q2 "轻量级"在规格里到底是怎么写给前端的（逐字原句）

### ③.1 真的能拿来拒绝一个动画／动效的句子（＝ owner 要收窄的射程）

| 判据句 | 出处 | 它拿什么理由管前端 |
|---|---|---|
| `| 动效 | **motion**（原 Framer Motion） | MIT | ⚠ **只用于面板内**；悬浮球一律原生（Direct2D）。且**禁止无限循环动画**（D32） |` | `docs/PLAN.md:3516` | **D32 的数**当理由，禁掉面板内无限循环动画——今天没有任何仪器执行它，纯靠引用 |
| `Material Icons（与 Raycast 视觉语言冲突）· Lottie（体积与 CPU，违背 D32）。` | `docs/PLAN.md:3451` | **体积与 CPU ⇒ D32**，禁掉一整类动效技术选型 |
| `**审批等待（L2）** … ❌ 无限循环脉冲（视觉噪音 + 违背 D32 的 CPU 约束）❌ 用红色` | `docs/PLAN.md:3486` | 同上：CPU 约束落到一枚具体脉冲 |
| `**所有循环动画只允许出现在非空闲态**（`Thinking`/`Acting`/`Speaking`/`Listening`/审批等待）。` | `docs/PLAN.md:3497` | 直接给面板动效划生命周期 |
| `……（走合成器，不触发重绘）。**面板关闭后不得有任何动画在跑。** 这是原型阶段就要立下的规矩，` | `docs/PLAN.md:3499` | 同一句里前半是渲染成本（合成器/重绘），后半是"关闭后不得跑" |
| `**视觉规范的可执行版本在 §17.3**（颜色/字号/间距/圆角/动效的完整 token）` | `docs/PLAN.md:990` | 把"动效"钉进 C21 token 体系（与 §②.3 T2/T11 同源） |
| `UI 验收：视觉质量纳入切片 done 判据，且**必须人工验收（用户本人）**` | `docs/PLAN.md:1054` | 不是拒绝句，是**owner 本人那一票**——放开与否的最终判据在这条 |
| `- **动画纪律**：`Sleeping` 态完全静态（不启动定时器，CPU≈0）；仅 … 帧率上限 **30fps**，` ／ `motion（仅面板内，禁无限循环动画）` ／ `| 审批等待 | chip 变 warn 描边 + 队列深度角标；脉冲只跑一次（不无限循环） |` | `docs/specs/SPEC-08`（§2 动画纪律 `:19-22`；技术栈 `:181`；动效表 `:190`、`:194`） | 把上面三句 PLAN 句抄进切片规格，其中 `:181`/`:194` 落在**面板**、`:19-22` 落在**悬浮球（原生侧）** |

**可直接转述给 owner 的三句（他要收窄的就是这三句）**：`PLAN.md:3516`「禁止无限循环动画（D32）」、`PLAN.md:3451`「Lottie（体积与 CPU，违背 D32）」、`PLAN.md:3486`「❌ 无限循环脉冲（视觉噪音 + 违背 D32 的 CPU 约束）」。

### ③.2 本来就管不到前端的句子（动了是白改契约）

| 句 | 出处 | 射程 |
|---|---|---|
| `**轻量级 = 运行时资源约束，不是代码行数约束。**` ＋ `约束对象：用户不用时必须真睡眠、不持续占用 PC 资源…不允许动辄数 GB 内存，不允许 CPU 长期十几%~几十%` ＋ `**不**约束：代码总行数、功能覆盖面、错误处理完备度、测试投入、权限模型、签名分发` | `docs/PLAN.md:65-69`（D1） | **D1 通篇没有一句提到动画／过渡／视觉**。"轻量级"作为拒绝动画的口径，在 D1 里**找不到授权句**——这是本票最重要的一条现量结论 |
| `| `Sleeping`（空闲，KWS 关，**无子进程**） | **≤ 25MB / ≤ 40MB** | ≤ 0.5% | …… 悬浮球**静态无动画** · 句柄 <300 · GDI <200 · goroutine ≤6 |` | `docs/PLAN.md:2255`（D32 表） | **悬浮球＝Win32 原生 Direct2D**，不在 `frontend/`；这条管的是常驻内存 |
| `- **动画与 D18 的冲突处理**：`Sleeping` 态**完全静态**……且**帧率上限 30fps**、动画期间 CPU 计入 D18 的工作态预算` | `docs/PLAN.md:1899-1901` | 同上，**球（原生侧）**，不是面板 |
| `首 token → TTS 首帧音频 ≤ 800ms` | `docs/PLAN.md:538`、`:2330` | 语音链路 |
| `| **`[panel]`** | `enabled` `width` `height` `keep_alive_in_session`(true) `scale` | `hot` |` | `docs/PLAN.md:2743` | 生效级别表，**不是约束**；但它是 §④ 那一档的落点 |
| `| **P11** | **WebView2 冷拉起实测**……` ＋ `若实测冷拉起 >2s……会牺牲 D29 的视觉质量` | `docs/PLAN.md:1568`、`:3227`、`:1595` | **启动延迟**，不是动画帧率；本单一枚不碰 |
| `**悬浮球视觉的平台约束**：它叠在**任意桌面壁纸**上 → 必须跟随系统明暗主题或自适应对比度` | `docs/PLAN.md:1041` | 球的可读性约束；与 §②.3 T2 的"单一色源"是同族判据，**不是性能句** |

⇒ **结论一列**：owner 要收窄的只有 **`3516`／`3451`／`3486` 三句（外加 `3497`／`3499` 前半句）**；`D1` 与 `D32` 的表、球侧 30fps、TTS 首帧、WebView 拉起——**本来就没写前端动画，一枚都不用动**，动了纯属白改契约（AC#6 的人工批准也省了）。

## ④ Q3 动画开关要落的字段与生效链（复核＋现量）

### ④.1 上游"18 枚"读数**不成立**；"零读者"成立

尺（现跑，逐枚）：
- `type PanelSection struct` 定义在 `internal/config/schema.go:525-535`，字段共 **5 枚**：`Enabled`(default `true`)、`Width`(default `640`)、`Height`(无 default)、`KeepAliveInSession`(default `true`)、`Scale`(无 default) ⇒ **带 `default` 的是 3 枚，不是 18 枚**。
- 生产读者：`grep -rn "PanelSection" --include=*.go internal/ cmd/ | grep -v _test.go` ＝ **只有 schema.go 自身的 3 行**（`:123` 字段挂载、`:524` 注释、`:525` 定义）。
- `grep -rn "\.Panel\." --include=*.go internal/ cmd/ | grep -v _test.go | wc -l` ＝ **0**。
- `grep -rn "Panel\.Enabled|Width|Height|Scale|KeepAliveInSession" --include=*.go .`（排 `_test.go`）＝ **0 命中**。
- 分母旁证：`internal/config/schema.go` 全文件带 `default:"…"` 的字段共 **69 枚**（`grep -c 'default:"'`）。

⇒ **复核结论**：**"整节零读者"这条成立且比上游更强**（`[panel]` 5 枚字段全零读者，含票 180 说的那批）；**"18 枚带 `default` 的字段"这个枚数我复不出来**（现量 5 枚字段／3 枚带 default）。若上游那 18 枚是别的口径（跨 section 的视觉字段并集？），需要 `180-c1`／`182-c1` 自己更正口径——**本表按现量记 5/3，不照抄**。

### ④.2 这一枚开关应该归 D36 哪一档

- 判据句：`docs/PLAN.md:2726`「**三档生效级别**：`hot`（立即生效）· `reload`（需重载子系统，如换 ASR 模型）· `restart`（需重启进程）」；`docs/PLAN.md:2743` 把 `[panel]` 整节钉成 **`hot`**；`internal/config/schema.go:524` 的注释也写「`PanelSection is [panel]; hot-tier.`」。
- ⇒ **推荐：`hot`**，落点就是 `[panel]`（例：`[panel] animation_tier = "full" | "reduced" | "off"`，与既有 `enabled/width/height/scale` 同节，**不新造 section**——SPEC-12 §薄索引 与 D36 的 section 树是枚举过的，新节名会让 GUI 与 TOML 对不上，`PLAN.md:2723-2724` 写的正是这个病）。
- **归 `hot` 意味着用户改完要做什么才看得见效果**：什么都不用做——**但前提是那条腿真的在跑**：Go 侧读到变更 → 推一次面板快照 → 面板换档。今天 `[panel]` 零读者（§④.1），所以 `hot` 这句现在只是纸面。

### ④.3 最小生效链是哪一条腿（指现有最近邻，不新造抽象）

- 现成最近邻＝**快照双侧钉**：`internal/panel/composer_test.go:48`（＋ `bridge_test.go:135`、`approval_test.go:109/:126`）已经在断言"Go 快照每枚 JSON 键必须在 `frontend/src/lib/panel.ts` 的接口里有声明"。
- ⇒ 最小生效链＝**三格**：`schema.go` 的 `[panel]` 加字段（`default`＋三档枚举校验，照 `CostSection.OverBudget` 那种"enum enforced"写法）→ `internal/panel` 的快照结构带该字段并真的读 `cfg.Panel`（**这条是今天唯一缺的腿**：`\.Panel\.` 零读者）→ `frontend/src/lib/panel.ts` 接口声明它。
- 缺任何一格，**现有仪器自己会红**（少第三格 → `TestComposerContractTypesMatchFrontend` 那形；少第二格 → 字段仍是零读者装饰品，票 180 的病）。**写腿不需要新造门禁**，这是本普查给它的最大便利。

## ⑤ Q4 关闭档"真的不跑动画"要怎么证（Go 侧／测试侧）

**具名结论：今天没有任何机读证据能证明"关闭档下动画真停了"。** 逐条自证：

1. `internal/panel`／`internal/ball` 现有读前端的 12 枚测试（§②.3 全表），断言面是**色值字面量、emoji 字符类、存储 API、host 调用点枚数、契约字段同步、git 入口**——**没有一枚**看 `requestAnimationFrame`、`transition`、`animation`、定时器或 CSS 变量档。
2. 规格里那句"零定时器"纪律（`SPEC-08:19-22`、`PLAN.md:1899-1901`）**射程是悬浮球（原生 Direct2D 侧）**，Go 侧确实可测（球自己的定时器不启动），但它**证明不了面板的动画停了**。
3. `d22scan` 四条射程（§②.1）里没有"动效"这一维 ⇒ 它既不会拦、也不会证。

⇒ **不许拿"设了个 CSS 变量"当证明**这句，本程把它翻成判据形状给写腿（**不落地、不产码**）：

- **判据形状 A（产物差分，最便宜）**：同一份前端产物在 `off` 档与 `full` 档之间必须存在**机读差异**——例：`off` 档下生成的主题/样式里不出现 `@keyframes`/`animation:`/`transition:` 字样（或 `rAF` 注册点数）。尺的现成最近邻＝`frontend_hygiene_test.go:101` 的 `main.tsx` import 闭包走法（它已经能"只扫真的会进 bundle 的文件"），照它的形状加一条字符类即可。⚠ 这条要求**关闭档改变产物或至少改变注入的样式源**，只切一个 `var(--anim: 0)` 不算。
- **判据形状 B（Go 侧可证的那半）**：快照里 `animation_tier` 字段值可断言（`off` 必须由 Go 推出去，且与 `[panel]` 值一致），并且**`off` 档下面板关闭后 Go 侧不启动任何与动效相关的定时器**——这半今天就能测，但**它只证明 Go 没驱动，不证明渲染侧没跑**。
- **判据形状 C（今天不可得）**：真正"不注册 rAF"要运行时观测渲染侧，需要 DevTools/CDP 类探针；仓里没有任何这类基础设施（本程未在 `internal/panel` 见到，前端侧未读）。⇒ **明写：这条路今天不通，别把它写成 AC。**

## ⑥ Q5 逐枚"放开／保留"清单（一行一枚，给 owner 一票一票地批）

> 五字段＝**文件／行／理由／边界／撤销口令**（票 193 AC#2 自陈的那五格）。
> 三条硬约束先钉住并已遵守：**① 零 emoji＝保留**（owner 自己的视觉规范，非性能门）；**② ban #6"面板不得成为批准/执行入口"＝保留**（票 191/192 另议，本单不开门）；**③ `C17` 白名单未加一枚、`thresholds.go`／golden／审批超时／SLO 阈值一字节未动**——本表**没有任何一行**建议放宽资源预算的**数**。

| # | 门禁／条款 | `file:line` | 今天真的会点谁的红 | 放开它买到什么 | 放开它的代价（哪一类本该被拦的东西不再算证据） | 建议 | 若要放开，`A##` 五字段 |
|---|---|---|---|---|---|---|---|
| 1 | 规格引用口径：**面板内禁止无限循环动画，理由 D32** | `docs/PLAN.md:3516` | **没有仪器执行它**——今天不点任何红；它只在人/agent 拒方案时被引用 | 用户可选的呼吸/脉冲类动效不再被"D32"三个字挡回 | 失去"D32 的 CPU 数**从不被面板动效消耗**"这条**纸面**保证；代价真实存在但**无仪器背书**（今天本来也没测） | **收窄（口径级，不改一字）** | 文件＝`docs/PLAN.md`｜行＝`:3516`｜理由＝owner 09-28 点名"轻量级不体现在前端动画上"，此句是三条引用口径之一｜边界＝**只撤"D32 的数"这半**理由；`3499` 的"面板关闭后不得有动画在跑"与 `thresholds.go` 的数**不在本条射程**｜撤销口令＝「193 那条口径收回，无限循环动画重新按 D32 拒」 |
| 2 | 规格引用口径：**Lottie 因"体积与 CPU，违背 D32"被禁** | `docs/PLAN.md:3451` | 同上，无仪器 | 前端可评估重动效库 | 失去"包体/运行时体积不被动效库冲击"这条**引用级**约束；`D29:1012` 的"语料丰富度 > bundle 体积"仍是明文，不会因此失守 | **收窄（口径级）** | 文件＝`docs/PLAN.md`｜行＝`:3451`｜理由＝同 1｜边界＝**只撤 D32 这半**；"Material Icons 与 Raycast 视觉语言冲突"是审美条款，**保留**｜撤销口令＝「Lottie 那条恢复 D32 引用」 |
| 3 | 规格引用口径：**审批等待脉冲"违背 D32 的 CPU 约束"** | `docs/PLAN.md:3486` | 无仪器 | L2 审批态可以有可选脉冲 | **⚠ 这一枚有安全含义**：审批卡的可感知性是 D33 的展示要求（"完整展示待执行命令与参数"，`PLAN.md:1428` S5 done 判据）。放开后"视觉噪音"不再算拒绝理由，但**不得放开到"审批态不可见"** | **收窄（只撤 D32 半句；审批态可见性保留）** | 文件＝`docs/PLAN.md`｜行＝`:3486`｜理由＝同 1，且 owner 原话点名"动画给用户选"｜边界＝**"不得弱化/不可见"这半不动**（`PLAN.md:2203` 那句"红色常亮环，不得渐隐、不得弱化"不在放开之列）｜撤销口令＝「审批脉冲那条恢复」 |
| 4 | 生命周期条款：**循环动画只允许在非空闲态／面板关闭后不得有动画在跑** | `docs/PLAN.md:3497`、`:3499`、`SPEC-08:19-22` | 无仪器 | （若一并放开）面板关闭后仍可能有动画在烧资源 | **这是真性能门，且与"轻量级"的定义直接对得上（`PLAN.md:67` 的"不用时必须真睡眠"）**——面板销毁后残留 rAF/定时器就是"持续占用 PC 资源" | **保留** | 不适用（不放开）。若 owner 要放开，本程标 `needs-approval`：它等于改 D1 的射程 |
| 5 | **T2 单一色源**：色值字面量只能活在生成主题里 | `internal/panel/frontend_hygiene_test.go:196`（尺 `:61`、闭包 `:101`） | **今天红**（仪器点名 `frontend/src/components/harness/right-rail.tsx:89`） | 前端可以就地写颜色/渐变/自定义样式，视觉自由度实质上升，且**不再需要动 `design/assets/tokens.css` 就能出新效果** | **失去"C21 令牌四向一致"的物质基础**：T11/`internal/ball` T12 都靠"唯一色源"才成立；放开后**主题切换/明暗自适应（`PLAN.md:1041`）不再有机读保证**，且**散落的颜色**这一类缺陷从此不算证据 | **收窄（不是全放）**：建议把射程从"任何非生成主题文件的任何色值字面量"收到"**不得新增第二套主题源**（局部一次性色值需在前端登记处出现）"。⚠ **改它＝改 `internal/panel/*_test.go` 产码，只读程不动** | 文件＝`internal/panel/frontend_hygiene_test.go`｜行＝`:196`/`:216`（红话）｜理由＝它是今天**唯一一枚因视觉自由度而点红**的门禁，owner 点名要逐枚决定｜边界＝**T11（`tokens_fourway_test.go`）与 `internal/ball` T12 不在本条射程**（owner 09-25 按过的不动）；emoji／批准面／契约钉不动｜撤销口令＝「色值字面量那条恢复严格」 |
| 6 | **T1 面板无状态**（禁 localStorage/IDB/cookie/SW/fs） | `internal/panel/frontend_hygiene_test.go:169`（清单 `:48-56`） | 今天不红（失败面为空） | 买到"动画偏好可以存前端"——**但那是假买**：偏好该存 `config.toml`（§④） | 失去"状态持有者唯一"这条 D22 级结构保证；`PLAN.md:1044` 明写禁止；面板一旦自存，历史列表分叉第一次就出现在这里 | **保留** | 不适用 |
| 7 | **T3 vendored 基座必须至少一枚被挂载** | `internal/panel/frontend_hygiene_test.go:222`（红话 `:240`） | 今天绿；若前端换掉/删空 `ai-native` 即红 | 前端可自由换组件基座 | 失去"票 77 的 D29 落点确实在用那套基底"这条证据；纯产品决定，不是安全 | **保留（换基座是 owner 的一句话，不是放开门禁）** | 不适用；owner 若要换基座，应写 `A##` 记"前提已改"，不是删这枚测试 |
| 8 | **T4 面板零 emoji（Go 内侧孪生）** | `internal/panel/frontend_hygiene_test.go:246` | 今天不红；写进面板的 emoji 字符类即红 | 面板文案可用 `✓`/`≤`/`→` 类字形 | 失去 owner 亲自定的**视觉规范**（非性能门）；且它比 ban #8 更严（`:64-71` 自陈无注释豁免）——放开它等于改视觉语言，不是改性能 | **保留**（派单硬约束①） | 不适用。**如实登记一处已知未定案**：字符类"规格比仪器宽/窄"的缺口在票 141／台账 `A201②`，本程未再动 |
| 9 | **ban #8 扫描器零 emoji（`frontend/` 85 枚、`design/` 39 枚）** | `tools/d22scan/main.go:164`＋`:557-564` | 今天 clean；命中即整扫描退码非 0 | 同上 | 同上，且它是 CI 自动扫的那把 | **保留**（硬约束①） | 不适用 |
| 10 | **ban #6 面板批准面 scope** | `tools/d22scan/main.go:249`＋`:395-410` | 今天不红（85 枚零命中） | **什么都不买到**——它不挡动画、不挡视觉 | 放开＝"面板侧不得成为批准/执行入口"这条**安全边界**失守（D33/F2） | **保留**（硬约束②；票 191/192 另议，本单不开门） | 不适用 |
| 11 | **T5 ban #6 内侧孪生测试** | `internal/panel/frontend_hygiene_test.go:281` | 今天不红 | 同 10 | 同 10 | **保留** | 不适用 |
| 12 | **T7 渲染侧只留一道门**（`frontend/src` 内 `postMessage` 调用点＝2） | `internal/panel/composer_test.go:502`（钉死句 `:510`，尺 `:381`） | 今天绿；任何新增桥腿即红 | 买到"前端可自辟一条与宿主通信的腿" | 失去 C17 单门这一**安全边界**；**动画开关不需要它**（走快照不算新调用点，见 §④.3） | **保留**，并在表里为写腿记一句：**别为动画开关加桥腿** | 不适用 |
| 13 | **T8 面板无 git 切换入口** | `internal/panel/composer_test.go:276`（红话 `:321`） | 今天绿 | 无（与动画无关） | 失去 owner 在票 92 AC#7 亲手砍掉的产品面 | **保留** | 不适用 |
| 14 | **T6／T10 快照双侧同步钉** | `internal/panel/composer_test.go:48`、`bridge_test.go:135`、`approval_test.go:109`/`:126` | **T6 今天红**（Go 的 `git` 字段未在 TS 声明） | 买到"Go 可以先推、前端以后再跟" | 失去"两边说同一种话"这一机读保证；**这条正是动画开关生效链的护栏**（§④.3），放开它＝把开关变成装饰品的那条路 | **保留**（且是 §④ 的正面依据） | 不适用。今天这枚红的处置＝**等 `181-r1` 自己补齐 TS 侧**，不是修门 |
| 15 | **T11 C21 令牌四向一致** | `internal/panel/tokens_fourway_test.go:10-13`、`:441` | **今天红**（红因＝`design/assets/tokens.css` 被别家会话删除未 staged） | — | — | **本单不动、不提议**（AC#0：只等 owner 一句「动 tokens 那枚」） | 不适用。若 owner 批，`A##` 由他起草 |
| 16 | **T12 `internal/ball` 的 CSS 腿** | `internal/ball/tokens_table_test.go:1364`、`:1468` | **今天红**（同一枚缺失文件；红话自陈"must never skip"） | — | — | **本单不动**（与 T11 同族，同一句口令） | 不适用 |
| 17 | **D1 / D32 的数与表**（内存/CPU/句柄/goroutine/启动） | `docs/PLAN.md:65-71`、`:2219` 起（`:2255` 表行）、`internal/slo` 的 `thresholds.go` | 不点前端红（射程是常驻进程与球） | 无 | 动它＝放宽资源预算本身，**owner 说的是"其他地儿严格没关系"** | **保留，一字节不动**（硬约束③） | 不适用 |
| 18 | **球侧动画纪律 30fps／`Sleeping` 零定时器** | `docs/PLAN.md:1899-1901`、`SPEC-08:19-22` | 射程＝**悬浮球（Win32 原生 Direct2D）**，不在 `frontend/` | 若 owner 也要"球动画可选"，那是一枚**新的原生侧票**，不是本票的放开项 | 动它＝D32 的 Sleeping 行失守 | **保留**，并**具名上报**："球动画可选项"本票未覆盖 ⇒ **未定义即停** | `needs-approval`（若 owner 说要） |
| 19 | **`AwaitingApproval` 红色常亮环"不得渐隐、不得弱化"** | `docs/PLAN.md:2203` | 不点红；但会被"动画可关"顺手动到 | — | 放开＝审批态的不可逆警示强度不再被保证（D33 必修族） | **保留**，且写腿**不得**把它当"动画档"一起关掉 | 若 owner 说要，`needs-approval` |
| 20 | **`C17` 方法白名单／审批超时／golden／SLO 阈值** | 派单 §2-Q5 硬约束③点名 | — | — | — | **一枚未加、一字节未动**（本程零产码） | 不适用 |

**统计**：会点前端红的门禁共 **2 类 4 枚**（`d22scan` ban #8 的 `frontend/` 分支＝T4 孪生；T2 单一色源）；今天真的红的是 **4 枚测试**，其中 **3 枚红因不是门禁**（T6、T11、T12）、**1 枚红因正是射程之争**（T2）。建议放开／收窄的只有 **4 行**（#1、#2、#3＝口径级引用收窄；#5＝射程收窄，需改产码 ⇒ 派写腿），其余**保留**。

## ⑦ 我判错/没测的（逐名）

1. **`internal/panel/l2_grant_boundary_test.go` 未逐枚读完**：它命中 `frontend/` 的 4 处里我看到的是注释与历史取件（`:19`/`:25`/`:1252`/`:2085`），但文件有 2000+ 行，我**没有**证明"它今天不读 `frontend/` 工作树内容"。⇒ 未现量。
2. **`composer_test.go:677` `TestComposerRenderFixtureTellsTheTruth` 未读**：注释提到 `frontend/fixtures/composer-states.html`（`:668`）是渲染产物，它是否读前端内容并据此判红——**未现量**。
3. **`internal/ball` 的枚数只取了起手一次**（1 FAIL／其余 PASS，`go test -count=1 ./internal/ball/`），**终态未复跑**（工具调用顶到硬顶前我优先复跑了派单 §3 点名的那四把）。
4. **`go test ./internal/panel/` 起手与终态各跑一次，两次名册逐字相同（3 FAIL）**——派单预警"至多 4 枚因在飞脏件的红"今天只兑现 1 枚（T6），另外 `composer.go`/`pump.go` 起手已不在脏件清单里（现量只有 `git.go`/`git_test.go`）；**我没有去验证那两枚是否已被 `181-r1` 提交**（那需要读它们的 diff，属于评论别家的活，未做）。
5. **"18 枚零读者"的上游口径我没拿到定义**：我只证明了自己的现量（5 枚字段／3 枚带 default／`\.Panel\.` 零读者），**没有**去读 `docs/evidence/s1/180-182-panel-fields-census-c1.md` 的表体（本程按"不许照抄、自己现跑"做，枚数分歧已如实登记在 §④.1）。
6. **Q4 判据形状 A 的可行性未证**：我**没有**打开 `frontend/` 的任何产物或源文件（写面闸门禁止），所以"`off` 档是否改变生成的样式源、能否机读差分"这一条**只是形状，不是现量结论**；写腿落地前不许当 AC。
7. **`d22scan` 的 `design/ 39` 我没拆枚数构成**（哪些是 `.html`/`.md`），因为拆开要列 `design/` 下的文件名，接近写面边界；只给了总数与"16 枚未 staged 删除同时成立"的计数解释。
8. **`ban #8 internal/`＝436 的"为什么不是 433"我给的是旁证（0 未跟踪＋find=481），不是逐枚归因**：到底是哪几枚票新增了 3 枚 `.go`，我没查 blame（那是别人的历史，且与"要不要放开前端"无关）。
9. **一次工具调用上的自纠**：起手那把 `sh scripts/d22scan.sh` 我用了 `tee`，但当时 `probes/193/a1/` 还没建，`tee` 报了 `No such file or directory`，**读数取自 stdout 未被截断的关键行**，随后我把同一份读数逐字抄进 `probes/193/a1/d22scan-start.txt`（终态复跑逐字一致，已交叉核对）。这是本程唯一一次写面失败，未产生任何文件副作用。
10. **我没有对本票之外的任何东西给出建议**：票 191/192（面板终端与浏览器）在本表里只出现在"另议，不开门"这一句里。

## ⑧ 门禁终态（最后一枚 commit 之后按 §⑨ 说明取数）

| 门禁 | 起手 | 终态 | 处置 |
|---|---|---|---|
| `sh scripts/d22scan.sh` 退码 | **0**（clean） | **0**（clean） | 一致 |
| `ban #8 internal/` examined | **436**（≠ 基线 433，已具名登记为别家已落地件） | **436** | 未动 |
| `ban #6 frontend/` / `ban #8 frontend/` / `ban #8 design/` | 85 / 85 / 39 | 85 / 85 / 39 | 一致 |
| `bans #1-5 internal/`、`cmd/`、`ban #7` | 209 / 23 / 21 | 209 / 23 / 21 | 一致 |
| 正向对照 `runtests.sh` | `PASS=34 FAIL=0 SKIP=0` | `PASS=34 FAIL=0 SKIP=0` | 一致 |
| `gate-clauses.sh`（**比红腿名册不比退码**） | BAD 名册＝**只有 `G6neg`**（声明 1 枚／实测 3 枚），腿数 14／不符 1／缺腿 0 | **只有 `G6neg`**，同值 | 与在册基线一致（无新增红腿） |
| `go test -count=1 ./internal/config/` | `ok` 0.951s | `ok` 1.055s | 绿 |
| `go test -count=1 ./internal/panel/`（**只照实记名册**） | **3 FAIL**：`TestComposerContractTypesMatchFrontend`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestC21DesignTokensFourWayAgree` | **同样 3 FAIL**，逐字同名 | **不许当绿、不许修**；两枚预告红（T2、T11）＋1 枚在飞件红（T6） |
| `go test -count=1 ./internal/ball/` | **1 FAIL**：`TestC21TableColourRowsMatchTokensCSS`（`:1468`） | 未复跑（见 §⑦.3） | 照实记 |
| `git status --porcelain -- internal/ cmd/ go.mod` | 起手 ` M internal/panel/git.go`、` M internal/panel/git_test.go` | ` M internal/panel/git_test.go` | **别人的在飞件，本程一枚未提交、未还原** |
| CLI 那一面 | **未跑**（派单 §3 写"本单不需要"） | 未跑 | 若跑须带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（票 98 的 `0xc0000135`） |

## ⑨ 被拒调用＋零删除自证＋工具调用终值

- **本程提交清单（显式 pathspec，各一枚 commit）**：
  1. `420dda62` — `docs/evidence/s1/193-frontend-gates-census-a1.md` ＋ `.scratch/wisp/probes/193/a1/d22scan-start.txt`
  2. （本表终稿这一枚，见提交信息）＋ 票 193 的 Progress log 一枚**追加** ＋ `probes/193/a1/` 的读数影件。
- **AC 框**：票 193 的 8 枚复选框（AC#0–AC#7）**一枚未勾**；原句一字未改，只在 `## Progress log` 末尾**追加**一行。
- **零删除自证**：本程未执行 `rm`／`rmdir`／`unlink`；未跑 `git clean`／`git restore`／`git checkout .`／`git stash`／`git reset`／`git rebase`／`git commit --amend`／`git worktree`；未 `git push`；未 `git add -A`／`git add .`／`-a`；**未跑 `probes/161/r6/flip-declaration.sh`**（派单 §1 明令）；**未改 `scripts/d22scan.sh` 与 `probes/154/gate-clauses.sh`**——重跑 `gate-clauses.sh` 前先按 §1 那把尺 grep 过它的写出路径（现量：命中里**没有** `WriteFile`/`OpenFile`/追加重定向，只有 `printf`→stdout 与 `git grep` 只读操作），确认原地不洗别人的读数。
- **写面自证**：`internal/**`、`cmd/**`、`frontend/**`、`design/**`、`go.mod`／`go.sum`、`docs/PLAN.md`、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量 **零字节未动**（§⑧ 那两条 `git status` 即凭据：终态只剩别家那一枚 ` M internal/panel/git_test.go`，本程一枚未提交）。
- **被拒调用**：无（本程没有被权限系统挡下的调用，因此没有绕行尝试）。
- **工具调用终值：30 枚／硬顶 35**。第 25 枚之后**未再开新探索**（此后全部用于落表、终态取数与提交）；派单 §4 要求"≤5 枚调用内落骨架"，本程实际在**第 12 枚**落盘、**第 16 枚**提交第一枚 ⇒ **这一条我没做到**，原因是起手把 Q1 的三把尺（扫描名册、测试 grep、脏件清单）连着跑完了才动手，记在这里不辩解。

## ⑩ next

**要 owner 先批（一票一枚，照 §⑥ 的行号点）**：
- `A##` #1/#2/#3——三条**引用口径**收窄（`PLAN.md:3516`、`:3451`、`:3486` 的 D32 半句）。**不改 `PLAN.md` 一字**，只登记"这三句不得再被引用来拒绝一个用户可选的动画档"。这是本票唯一"零产码就能生效"的放开。
- `A##` #5——`frontend_hygiene_test.go:196` 的射程收窄（**要改产码**，本只读程不动）。owner 批了才能派写腿。
- #15/#16（T11、T12）——只等 owner 那句「动 tokens 那枚」；不说就继续当已知常红读。
- #18——"悬浮球动画要不要也做成可选项"是**未定义**（本票只写面板前端），**未定义即停**，请 owner 定。
- #19——确认"审批态常亮环不得渐隐/弱化"**不进入动画档的射程**。

**可以直接派写腿（不需要 owner 先批）**：
- §④.2/§④.3 那三格最小生效链：`[panel] animation_tier`（`hot`，D36 `PLAN.md:2743` 现成档位）→ `internal/panel` 真的读 `cfg.Panel`（今天 `\.Panel\.` **零读者**，这是唯一缺的那条腿）→ `frontend/src/lib/panel.ts` 接口声明。护栏用现成的 `TestComposerContractTypesMatchFrontend` 那形，**不需要新造门禁**。
- §⑤ 判据形状 B（Go 侧：`off` 档由 Go 推出、且 Go 不因面板动效起定时器）——今天就能测，先落这条。
- 写腿提示（来自 §⑥ #12）：**别为动画开关新增一条 `postMessage` 桥腿**，走快照；否则会在 `TestTheRendererHoldsExactlyOneDoorToTheHost` 上吃一发与需求无关的红。

**要转给前端会话的**（本程只给计数，不引内容）：`d22scan` 对 `frontend/` 的射程今天 **85 枚文本件**（ban #6 与 ban #8 各 85），`design/` **39 枚**；其中**只有"色值字面量"那一枚（T2）是视觉自由度之争**，红点已由仪器点名一次（`right-rail.tsx:89`）。零 emoji（ban #8／T4）与批准面（ban #6／T5）**不在放开之列**，请前端会话按现状遵守。
