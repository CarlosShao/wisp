# 票 114 AC#1 — 现状表（只读审计，供下一段写码代理作输入）

**审计代理**：`audit-ticket114-ac1`（只读，未参与实现，**本轮零写码、零 git 变更、零测试运行**）
**锚定 sha**：`7ad6eb4`（= 本表所有 `file:line` 的读数时刻；下文行号一律按该 sha 的工作树）
**分支**：`dev`。工作树状态：仅一枚未跟踪文件 `docs/evidence/s1/131-adversarial-acceptance.md`（票 131 的在飞产物，**本代理未 add、未读其内容作为依据**）。
**边界遵守**：`frontend/**` 全程**只读**（读文件 + grep），未写入任何一枚；未碰 `internal/winsec/**`、`cmd/wisp/**` 的写；未碰冻结件（`docs/PLAN.md`、`docs/specs/*`、`internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`）。
**脏工作树校正（重要，本表曾中招并已由本代理改回）**：读数开始时工作树**不干净**——`cmd/wisp/{run.go,models.go,providers.go,doctor.go}` 是他人未提交的在飞修改（`run.go` 净 +14/−1），另有 1 枚已暂存的删除 `.scratch/wisp/issues/121-*.md` 与 1 枚未跟踪 `docs/evidence/s1/131-adversarial-acceptance.md`（**均非本代理所为，本代理一枚未 add**）。
⇒ 本表**所有 `cmd/wisp/**` 行号一律改用 `git show HEAD:<path>` 复核**（唯一被纠的一处：`Modes: rt.modes` 由脏树的 `:379` 改回 HEAD 的 `run.go:366`）；
`internal/`、`frontend/`、`tools/`、`scripts/`、`.github/` 在 `git status --porcelain` 里**零修改** ⇒ 那部分的 `file:line` 工作树 = HEAD，可直接引。
两条"0 生产调用者"的关键断言另按 HEAD 独立复跑一次（`git grep`）：`cmd/` 非测试代码里 `.Set(` = 0、`internal/panel` 的非测试 importer = 1 枚（`cmd/wisp/panel_assets.go`）⇒ **与脏树读数一致**。
**会话期间 HEAD 又前进**：本表提交时父提交是 `7a5dab2`（票 128 的 docs，非本代理所为）。`git diff --numstat 7ad6eb4 HEAD -- internal/panel internal/perm frontend tools scripts .github cmd/wisp/run.go` **输出为空** ⇒ 本表引用的每一枚文件在 `7ad6eb4` 与 `HEAD` 之间**零漂移**，行号两处通用；只有上面那批**未提交**的 `cmd/wisp/*.go` WIP 例外（按 HEAD 引）。
**证据档位图例**：〔独立复现〕= 本代理自己 grep/读原文得出；〔引台账〕= 只引用既有裁决表/账目原文，未独立复算。

命令口径说明（避免仪器假读数）：以下"调用者清单"均由 `grep -rn <符号> --include=*.go .` 在全仓（含 `cmd/`、`internal/`、`tools/`）采得，**逐条人工分类为生产/测试**（分类判据 = 文件名是否以 `_test.go` 结尾 + 是否位于 `cmd`/`internal` 之外）。"0 生产调用者"这类结论下面**逐条列出行号**，不给"应该没有"式断言。

---

## ① 现状表：定义处 / 全部调用者 / 断在哪一行

### ①.1 `panel.ParseComposerRequest`（渲染进程 → 原机的**唯一**入口）

- **定义**：`internal/panel/bridge.go:77` `func ParseComposerRequest(raw string) (ComposerRequest, error)`
- **方法白名单（闭集 switch）**：`internal/panel/bridge.go:97 knownComposerMethod`，被 `bridge.go:83` 调用；四条方法常量在 `internal/panel/bridge.go:35-38`。
- **全部调用者**：

| 调用点 | 分类 |
|---|---|
| `internal/panel/bridge_test.go:66, 83, 92, 102, 112` | 仅测试 |
| — | 生产 = **0** |

- **`knownComposerMethod` 的调用者**：只有 `bridge.go:83`（在 `ParseComposerRequest` 内部）。⇒ 白名单本身也**只在一条没有生产听众的函数里被读**。
- **`RefusedEnvelopeForUser`（拒答形状）的调用者**：`bridge_test.go:119`（仅测试）⇒ 生产 = **0**。
- **断在哪一行**：**`cmd/wisp` 里没有任何一处从 WebView2 的 `WebMessageReceived` 事件取文本并调 `ParseComposerRequest`**。
  具体定位：`internal/panel` 在**整个非测试代码树里只被 `cmd/wisp/panel_assets.go` import**（全仓 `grep -rln "wisp/internal/panel" --include=*.go` 的非测试结果枚数为 1），
  而它只用到 `panel.NewApprovalCardView`（`cmd/wisp/panel_assets.go:47`）、`panel.BuiltinAssets`（`:67`）、`panel.EntryFile`（`:103,114`）——**没有一处 composer 请求的读取**。
  ⇒ 断点**不在 `internal/panel` 内部**，而在"**宿主事件 → `ParseComposerRequest`**"这一跳：它**不存在**，所以无法给出 `file:line`（这是本表最需要报清楚的一条：断的是"缺文件缺函数"，不是"某行没接上"）。

### ①.2 `perm.Store.Set`（档位写入的**唯一**语义点，含 L2 确认腿）

- **定义**：`internal/perm/store.go:175` `func (s *Store) Set(ctx context.Context, to risk.Mode, origin, actor string) error`
- **内部已有的 nil-`Confirm` 拒**：`internal/perm/store.go:191-197` —— **只在 `to == risk.ModeAutoApprove` 这一支**上检查 `s.confirm == nil` 并 fail-closed；切到 `ask_every_step` / `ask_high_risk` 两支**不经过任何确认腿**。
- **全部调用者**：

| 调用点 | 分类 |
|---|---|
| `internal/perm/store_test.go:116, 136, 142, 152, 185, 204, 225` | 仅测试 |
| `internal/perm/ticket90_persist_test.go:154, 404` | 仅测试 |
| `cmd/wisp/run_mode101_test.go:231, 234, 520`（`rt.modes.Set(...)`） | 仅测试 |
| — | 生产 = **0** |

- **构造侧的对照**（说明"0"不是 grep 漏）：`cmd/wisp/run.go:337 perm.New(perm.Options{Manager, Confirm, Logf})` 确实**装配了** Store，`cmd/wisp/run.go:347 rt.modes = modeStore`，但生产路径只读它：`run.go:356 modeStore.PermissionMode()`、`run.go:251` 注释、`tools.Options.Modes: rt.modes`（**`run.go:366`**，把 Store 当 `tools.ModeSource` 喂给决策链）。`cmd/wisp` 非测试代码里 **`.Set(` 出现次数 = 0**。
- **L2 确认腿的实现**（AC#2 要钉的那条腿）：`cmd/wisp/run.go:400 (rt *agentRuntime) confirmModeSwitch`，装配点 `cmd/wisp/run.go:333-336`（`confirm := s.modeConfirm; if confirm == nil { confirm = rt.confirmModeSwitch }`）。⇒ 生产注入的**不是 nil**，是"走 C18 gate 的真卡片"；而 `run.go:117-123` 自陈：控制台跑没有原生通道 ⇒ 该卡片**总是超时/被拒**（fail-closed）。
- **断在哪一行**：**没有"composer ModeRequest → `Store.Set`"的处理器**。上游同样断在 ①.1 那一跳；下游 `ModeRequest.Parse()`（`internal/panel/composer.go:139`）的**生产调用者 = 0**（全表：`composer_test.go:192, 201`，皆测试）。⇒ 即使接上 `ParseComposerRequest`，还差 `mode.request` 的 **handler**（该 handler 是 AC#2 的落点）。

### ①.3 面板三块输入（档位 / 附件 / 工作区）逐个

**共同事实**：`panel.Snapshot` / `NewSnapshot` / `NewComposerState` 三者构成"面板显示的那一块真值"，它们的**生产调用者也全为 0**（下表逐条）。
⇒ 票面说的"三块输入"今天**显示腿和请求腿都没接**，不只是请求腿。

| 输入 | 定义（视图 / 请求 / 原生腿） | 全部调用者（区分） | 断在哪一行 |
|---|---|---|---|
| **档位 mode** | 视图 `ModeView` `internal/panel/composer.go:93`、`NewModeView` `:110`、`ModeUnknownView` `:120`；请求 `ModeRequest` `:128`、`ModeRequest.Parse` `:139` | 生产 **0**：`NewModeView` 仅被同文件 `composer.go:202`（`NewComposerState` 内）调；`ModeUnknownView` **零调用者（连测试都没有）**；`ModeRequest.Parse` 仅 `composer_test.go:192,201` | ①`NewComposerState`（`composer.go:197`）生产 0 调用者（仅 `composer_test.go:210`）⇒ 快照从不被构造；②`ModeRequest` 无 handler（①.2）⇒ 两处都断 |
| **附件 attachment** | 载荷 `AttachmentPayload`/`DecodeAttachmentPayload` `internal/panel/attachments.go:381`；执行体 `AttachmentBroker` `:169`、`NewAttachmentBroker` `:179`、`(*AttachmentBroker).Ingest` `:199`；`AcceptedMIMETypes` `:153`、`OutgoingMessage` `composer.go:213`、`ForAgent` `:226` | 生产 **0**：`NewAttachmentBroker` 仅 `attachments_test.go:91,132,229`；`Ingest` 仅测试（`attachments_test.go` 多处）；`DecodeAttachmentPayload` 仅 `attachments_test.go:339,360,373,383`；`ForAgent` 仅 `attachments_test.go:284`；`AcceptedMIMETypes` 的非测试调用者全在 `internal/panel` 自身（`attachments.go:250`、`composer.go:77,205`），而这些函数本身生产 0 | 生产**从不构造 `AttachmentBroker`**（唯一 sink/guard 注入点在测试里）⇒ 断在"`MessageHandler` 里没有 broker 实例可注入"，同样无既有行可指 |
| **工作区 workspace** | 视图 `WorkspaceView` `composer.go:161`；请求 `WorkspaceRequest` `composer.go:180`；原生腿 `RequestWorkspaceSwitch` `internal/panel/workspace.go:76`、接口 `PathScope` `:36`、渲染 `WorkspaceViewFromRoot` `:60`、`UnsetWorkspaceView` `:55` | 生产 **0**：`RequestWorkspaceSwitch` 仅 `workspace_test.go:53,75,103,123,141,161`；`PathScope` 的非测试出现只在 `workspace.go`（`*tools.PathCanonicalizer` 注释声称满足它，但**无静态断言**）；`WorkspaceRequest` **全仓零引用（除定义行 `composer.go:180`）**；`WorkspaceViewFromRoot`/`UnsetWorkspaceView` 非测试调用者仅 `workspace.go:77,82,92,99`（`RequestWorkspaceSwitch` 自己）⇒ 传递性为 0 | 同 mode：缺 `workspace.request` 的 handler；且**`tools.PathCanonicalizer` 是否真满足 `PathScope` 三方法在树级别未被任何编译期断言钉住**（见 `next=` 的前置问题 P3） |

### ①.4 "接通之后哪一条从 0 变成 N"（票 114 AC#1 明问的那格）

按上表，**接一次 composer 线会同时把 8 个符号从 0 抬成 N**，且**抬的顺序决定风险**：

1. `ParseComposerRequest`：0 → ≥1（宿主事件处理器）
2. `knownComposerMethod`：随 #1 被动生效（它只在 #1 内部被读）
3. `RefusedEnvelopeForUser`：0 → ≥1（拒答路径，AC 要求"不许静默丢"）
4. `ModeRequest.Parse`：0 → ≥1
5. **`perm.Store.Set`：0 → ≥1 ← 这条是唯一真正的权限写**
6. `NewAttachmentBroker` + `Ingest` + `DecodeAttachmentPayload`：0 → ≥1
7. `RequestWorkspaceSwitch`（连带 `PathScope` 的生产实现）：0 → ≥1
8. `NewSnapshot`/`NewComposerState`/`Snapshot`：0 → ≥1（显示腿，票 35 的泵）

⇒ **判据落点**：AC#2 的原生侧门必须落在 **#5 的那一枚新 handler 上**（唯一写点），
且**门二（文本门）在 #1 落地那天立刻从"防线之一"变成"防线全部"**（`R-92-2` 的原话），所以 #1/#5 与 AC#2 的门**不能拆成两批发**。

---

## ② 原生侧门该落在哪（AC#2 / `R-92-2` 的正身）

### ②.1 今天已经有的那半条门（先看清，别重复造）

- `internal/perm/store.go:191-197`：`Set` 在 `to == risk.ModeAutoApprove` 且 `s.confirm == nil` 时**已经**fail-closed（写审计 + 返回 error）。
- 装配根已经保证生产里 `Confirm` **不是 nil**：`cmd/wisp/run.go:333-336`（`confirm := s.modeConfirm; if confirm == nil { confirm = rt.confirmModeSwitch }`），真身 `cmd/wisp/run.go:400`，它走 C18 gate 签发一张 L2 卡片；控制台无原生通道 ⇒ 该卡片恒超时/被拒（`run.go:117-123` 自陈）。
- **`internal/panel` 里今天没有任何 `Confirm` 概念**：`grep -rn "Confirm" --include=*.go internal/panel` 的非 `L2Confirm` 命中 = **0**。
- ⇒ **已有的那半条只挡 `auto_approve` 一档**。面板 handler 若直接照抄"`Set` 里已经有门"，那**切到 `ask_high_risk` / `ask_every_step` 两支今天没有任何确认腿**（`store.go:191` 的 if 不覆盖它们）。AC#2 的原文是"**mode 写入处理器**在注入的 `Confirm` 为 nil 时必须拒"——**主语是处理器，不是 `Set`**，所以它不是重复劳动；但也**不是**已经存在的性质。
  ⚠ 这条差异（"handler 级 nil 拒 = 三档全拒"vs"沿用 `Set` 的语义 = 只拒 auto_approve"）**是本表判不出来、必须写码代理先拍的问题**，见文末 `P1`。

### ②.2 最小改动面（按"只写 Go、零 `frontend/`"给出；行数是量级不是承诺）

| # | 文件 | 动作 | 量级 |
|---|---|---|---|
| 1 | **新建** `internal/panel/composer_handlers.go` | 一个 `ComposerHandlers` 依赖注入结构（`ModeWriter`（= `perm.Store.Set` 的签名，`store.go:175`）、`Confirm perm.ConfirmFunc`、`Attachments`、`Paths PathScope`、`Audit AuditFunc`）+ 三个 handler（`mode.request`/`workspace.request`/`attachment.add`）。**mode handler 第一句 = `if h.Confirm == nil { 审计 + 拒 }`**，且**必须在调用 `ModeWriter` 之前** | ~60-90 行 |
| 2 | **新建** `internal/panel/composer_handlers_test.go` | ①nil-`Confirm` ⇒ 拒 **且断言注入的 `ModeWriter` 一次都没被调用**；②非 nil ⇒ 走到 `ModeWriter`；③两条变异判据（见 ②.3） | ~80-120 行 |
| 3 | `cmd/wisp/run.go`（现有装配点附近，锚 `:333-347`） | 把 `rt.modes` 与 `confirm` 交进 handler；handler 挂到哪一行取决于 ②.4 | ~10-20 行 |
| 4 | `internal/panel/bridge.go` | **可能一行都不用改**：白名单 `:97` 已是闭集，路由到 handler 的 switch 属于 #1 | 0-10 行 |

**#1/#2 是本票 AC#2 的正体；#3 是本票 AC#1/AC#3 的正体。二者同批，否则门是空门。**

### ②.3 变异判据（票面 AC#2 明写，写码代理必须自证）

- 删掉 `if Confirm == nil` 那一段 ⇒ #2 的用例①**必须红**。
- 把它改成"只 `audit(...)` 不 return" ⇒ 用例①**也必须红**——**判据是"`ModeWriter` 没被调到"，不是"日志里说了拒"**；只断字符串会把这一发变异放过去（memory：判"放水"看断言有没有被动）。

### ②.4 会不会破另一平台（这条判据在 Linux runner 上有没有分母）

- **门与用例落在 `internal/panel` ⇒ 有分母。** 依据：`internal/panel` 全平台中性（`grep "go:build\|_windows\|runtime.GOOS" --include=*.go internal/panel` = **0 命中**），且它在 ubuntu 的 core scope 里：`scripts/portable-tests.sh:179`（`./internal/panel/ ...`）、`scripts/portable-tests.sh:141`（core 清单的 census 钉）。CI 的 ubuntu 步 = `Portable package tests (core scope)`（`.github/workflows/ci.yml:242-243`）。
- **装配根那一跳（#3）落在 `cmd/wisp` ⇒ Linux runner 上没有分母。** `scripts/portable-tests.sh:195` 的 cli scope = `./cmd/wisp/`，只由 `scripts/wisp-cli-tests.sh` 调（`portable-tests.sh:193-195` 注释：没 DLL 会在**加载期**死，ticket 98/`A105②`），CI 里只有 windows 腿跑它（`.github/workflows/ci.yml:378`，`runs-on: windows-latest` 于 `:310-311`）。
  ⇒ **`GOOS=linux go vet ./...` 只编译不执行**，它能证明"cmd/wisp 在 linux 下还能编"，**不能**当"这条门在 linux 上成立"的证据；反之，**若把 nil-Confirm 用例只写在 `cmd/wisp`，Linux 腿永远不会跑它**（只有 windows 腿有分母）。推荐：门的行为用例放 #2（两平台都有分母），`cmd/wisp` 里只放装配断言（"handler 的 Confirm 非 nil"这一条，接受它 windows-only）。
- **不破另一平台的三个前置条件**（写码代理照做即可，任一条破了 ubuntu 就会红）：
  ① handler 文件**不许带 `//go:build windows`**、不许 import `internal/winsec`/`internal/ball` 的 windows-only 符号；
  ② 任何"原生文件夹选择"必须以 func 注入（今天的 `PathScope` 就是这个形状，`internal/panel/workspace.go:36-47`），不许在 `internal/panel` 里直接调 Win32；
  ③ `internal/panel` 现在**没有** import `internal/perm`（`grep` 结果里 panel 侧无 `perm.` 生产引用）⇒ 一旦 handler 用 `perm.ConfirmFunc`/`perm.Switch` 类型，会新增一条 `panel → perm` 依赖。**这是本轮唯一真实的"破另一平台"风险位**（不是编译平台的破，是**依赖方向**的破：`internal/perm` 是 core scope 成员、且被 `cmd/wisp` 使用，方向安全；但请把 `ModeWriter` 写成**本地小接口**而不是直接收 `*perm.Store`，可把这枚新依赖整个省掉）。〔判据：本表仅静态核对，未编译验证——见 `P2`〕
- **WebView2 事件腿（"渲染 → `ParseComposerRequest`"）不需要写 `frontend/`**：页面侧已经在发（`frontend/src/lib/panel.ts:237,246,263,291` 四条字面量 + `bridge.postMessage` 于 `panel.ts`，`composer_test.go:161` 钉住计数=2）。缺的是 Go 侧接收：`grep "WebMessage\|AddHostObjectToScript" --include=*.go` 非测试命中 **0**，且**宿主票 33（`.scratch/wisp/issues/33-panel-host-c27.md`）状态 = `ready-for-agent`（未做）**、票 35（`35-panel-bridge-c17.md`）同 `ready-for-agent` 且 `Blocked by: 33, 34`。
  ⇒ **票 114 的正解与票 33 有地界重叠**：要么 114 只做 handler+门（可测、可绿、但"通路仍不存在"这条性质**不变**），要么 114 顺手把 33 的接收腿做掉（越界风险，且 AC#6 那条真机截屏已被 `A102③` 冻在 owner 派的外部 agent 手里）。**这是 `P0`：先答它再动手。**

---

## ③ 两道现有门的覆盖面（AC#4 / `R-92-1`）

### ③.1 先更正一件事：`R-92-1` 的"三种扩展名"在**当前树**已经不成立

- 现状：门二 = `composerModeWriteRe`（`internal/panel/composer_test.go:88-90`，7 个字面形状）**×** 文件过滤器 `rendererSourceFile`（`internal/panel/composer_test.go:370-377`）= **7 种扩展名** `.ts/.tsx/.js/.mjs/.jsx/.css/.html`。
- 落地 commit：`91b5fc4 style(92,AC#6)+test(92,AC#1): 把门二从字面量面挪到结构面`（`git log --oneline -- internal/panel/composer_test.go` 的顶端）。其注释 `composer_test.go:356-368` 逐字承认"门二过去只开 `.ts/.tsx/.css` 三扇门"。
- 同时新增了三枚结构判据：`dottedJoinRe` `:384`、`computedSendRequestRe` `:389`、`routeLiteralRe` `:392`，由 `TestTheRendererHoldsExactlyOneDoorToTheHost`（`:502`）执行。
  ⇒ **AC#4 的"要么扩扩展名、要么挂结构"这半格，主体已经被 `91b5fc4` 做过一部分**；票 114 剩下的不是重做，而是下面两个**实测仍存在的缝**（③.2 的差集 + ④ 的扫描根）。

### ③.2 差集（门二的文件面 vs d22scan 的两套文件面）

| 仪器 | 定义处 | 文件面 |
|---|---|---|
| 门二 | `composer_test.go:370-377` + walk `:597-616`（root=`frontend/`，跳 `node_modules`/`.git`/**`dist`**） | 7 类扩展名，**不读**其余 |
| d22scan **ban #6**（AC 里说的那道"文本门"） | 声明 `tools/d22scan/main.go:330-341`；实现 `walkText` `main.go:758-787` | **`frontend/` 下每一个文件**，只跳 `testdata`/`node_modules`/`.git`——**没有扩展名白名单，也不跳 `dist`**（与 `R-92-6` 一致） |
| d22scan 的**通用文本类** `isTextFile` | `main.go:877-881` | 15 类：`.md .txt .html .css .js .ts .tsx .jsx .json .yaml .yml .toml .go .svg .vue`（**不含 `.mjs`**） |

- **差集（ban #6 − 门二）= 除那 7 类之外的全部 `frontend/` 文件**。按 tracked 清单实测（`git ls-files frontend` = 40 枚）门二读不到的是 **11 枚**：
  `frontend/embed.go`、`frontend/VENDORED.md`、`frontend/.gitignore`、`frontend/dist/.gitkeep`、
  `frontend/package.json`、`frontend/package-lock.json`、`frontend/tsconfig.json`、`frontend/tsconfig.app.json`、`frontend/tsconfig.node.json`、`frontend/.oxlintrc.json`、`frontend/fixtures/l2-card-fs-delete.json`；
  外加**未跟踪但被 WebView 真正加载的** `frontend/dist/**`（`frontend/embed.go` 的 `//go:embed all:dist` + `internal/panel/assets.go:39-44`）。本机现存一枚陈旧构建 `frontend/dist/assets/index-UL9kYvYl.js`（mtime 09-21 15:10，早于 92 的 composer 落地；`grep -c` 四条 composer 路由字面量 = **0**）⇒ 这枚"正在被 embed 的产物里没有本票刚写的路由"是**独立于扩展名的一条覆盖面账**，登记备查。
- **反向差集（门二 − `isTextFile`）= `.mjs`（1 类）**，即 `frontend/scripts/{gen-tokens.mjs,vendor.mjs,vendor-shadcn.mjs}`。⇒ **如果 AC#4 被字面理解成"扩到与 `isTextFile` 一致"，会把这 3 枚文件从覆盖面里拿掉**（缩范围，票面明令禁止）。要扩就扩到 ban #6 的"每个文件"，不要扩到 `isTextFile`。

### ③.3 "把门二扩到与 d22scan 一致"会不会让某条已绿的用例变红

**答：不会。** 静态依据（本代理实跑 grep，非推断）：

- 门二的 7 个字面形状在**整个 `frontend/` 树**（40 枚 tracked + 本机 `dist/`，`node_modules` 不存在于跟踪集）里**只命中 1 行**：`frontend/src/lib/panel.ts:230`（`* call that actually changes the mode is perm.Store.Set on the Go side, which ...`）。
  它落在 `codeOnly`（`composer_test.go:644-663`）的"纯注释行"分支（以 `*` 起首）⇒ 被清空、今天不红。**把 `.json/.md/.go/.gitignore/...` 加进来不会新增任何命中**（上面 11 枚文件里 0 命中，逐类实测：`grep -rn '"panel\.' frontend --include=*.json --include=*.md --include=*.go` 空、门二 7 形状在 `dist/assets/*.js` = 0、`approval.decide` 在 `dist/assets/index-*.js` 与 `dist/index.html` 均 = 0）。
- 结构钉那三枚同理：`routeLiteralRe` 扫到的 `"panel.*"` 字面量 = 5 枚，全在 `panel.ts`（`:181,237,246,263,291`）⇒ 扩到 json/md/go 后**新增 0 枚未知路由**（`frontend/src/fixtures/` 是空目录）。
- **会变红的只有一种做法**：把门二的 walk 取消 `dist` 跳过（`:606` 的 `case "node_modules", ".git", "dist"`）。⇒ 那等于把判据挂到"`某人在提交前跑了 `npm run build` 没有`"上——`composer_test.go:356-368` 已经为这件事写了显式拒绝理由，**别动它**；要覆盖 dist 的可达性请交给 `AC#11`（把 `npm run render:composer && git diff --exit-code` 做成 CI 步）与 ban #6 本身（它已经读 dist）。

> ⚠ 以上全部为**静态**判定：本代理**未运行任何 `go test`**（任务禁止全仓跑、且同树有他人 WIP）。"某条已绿的用例不会变红"这条结论的**复跑权在写码代理**：`bash scripts/portable-tests.sh ./internal/panel/`（单包 scope）即可，命令与分母都已核过。

---

## ④ `knownComposerMethod` vs `frontend/` 那份方法清单（`R-92b-1`/`R-92b-2` 的输入）

**读 `frontend/` 只做对照，本代理未写其中任何一枚文件**（`A102③`/README 规则 7）。

### ④.1 三份清单的现状与差集

| 清单 | 定义处 | 成员 |
|---|---|---|
| **原生侧闭集**（Go 真正回答的） | 常量 `internal/panel/bridge.go:35-38`；判定 `internal/panel/bridge.go:97-103`（`switch`，`:83` 调用） | 4：`panel.mode.request`、`panel.workspace.request`、`panel.attachment.add`、`panel.message.send` |
| **测试侧允许词表**（结构钉的白名单） | `internal/panel/composer_test.go:409-418` | 5：上面 4 条 + **`"panel.approval.request"`（`:417`）** |
| **`frontend/` 实际字面量** | `frontend/src/lib/panel.ts:181, 237, 246, 263, 291`（全仓唯一命中处） | 5：同上一行 |

- **`frontend − knownComposerMethod` = {`panel.approval.request`}`（1 条）**。它是**有意**的：审批卡片走的是另一条封套（票 77 的地界），Go 侧回答它的代码**不在 `ParseComposerRequest` 里**（全仓非测试代码 grep `panel.approval.request` = **0 命中**，只有两枚测试文件 `composer_test.go:417`、`frontend_hygiene_test.go:32` 提到它）。⇒ 这条差集是"名单不同步"还是"设计上两份名单不同"，**当前树里没有任何断言能区分**（见 `P4`）。
- **`knownComposerMethod − frontend` = ∅**，而且**这一半已经有门**：`internal/panel/bridge_test.go:131-146` 遍历 4 枚 Go 常量、逐条要求在 `frontend/src/lib/panel.ts` 里出现对应字面量 ⇒ **"往 `knownComposerMethod` 加一条而不同时改 `panel.ts` ⇒ 用例红"这条今天已经成立**（AC#8 的一半不是从零开始，别重造）。它的两处不足：(a) 只读**一枚文件**（`bridge_test.go:133`），(b) 词表 `composer_test.go:409` 是**手抄副本**，没有任何断言把它与 `knownComposerMethod` 对齐（全仓 `grep composerRouteLiterals` = 4 命中：`:409` 定义、`:426` 使用、`:525` 打印，**无交叉断言**）。
- **`R-92b-1` 的 F5 形（运行时拼名）今天仍然全绿**，三枚判据各自留缝（`internal/panel/composer_test.go`）：
  `sendRequest("panel"+"."+noun+"."+verb, …)`
  ① `computedSendRequestRe` `:389` = `sendRequest\(\s*(?:[A-Za-z_$][\w$]*|`|\[)` ⇒ 实参以**引号**起首，**不匹配**；
  ② `dottedJoinRe` `:384` = `\.join\(\s*["']\.["']\s*\)` ⇒ 只认 `.join(".")`，**不认 `+` 拼接**；
  ③ `routeLiteralRe` `:392` = `"panel\.[A-Za-z.]+"` ⇒ 要求点号**在同一个字面量里**，`"panel"` 与 `"."` 分家 ⇒ **不匹配**。
  ⇒ **AC#8 剩下的活 = 堵 `+` 拼接这一形**（三选一：把 `computedSendRequestRe` 的交替里加上 `"`，代价是 `panel.ts` 现有 5 处合法字面量调用会被读成 computed、必须同时改判据；或加一条 `stringConcatRouteRe` 专认 `"panel"\s*\+\s*\.`；或把钉挂在**结构**上而非语法上）。**本代理不替你选**（`P5`）。

### ④.2 结构钉的扫描根比 `ban #6` 与"渲染器实际加载的文件集合"**都窄**——窄在哪几个文件

三枚根（都有实测数字）：

| 仪器 | 扫描根（`file:line`） | 实测文件数 |
|---|---|---|
| 结构钉 `TestTheRendererHoldsExactlyOneDoorToTheHost` | `filepath.Join(root, "frontend", "src")`，`internal/panel/composer_test.go:504`（种植变体 `:562`），walk 里再跳 `node_modules`/`.git`/`dist`（`:434`），再过 `rendererSourceFile`（`:370`） | **21**（= `git ls-files frontend/src` 的枚数；与 92b 报告日志原文 `21 files scanned` 逐字对上） |
| `ban #6`（d22scan） | `frontend/` **每个文件**，只跳 `testdata`/`node_modules`/`.git`：声明 `tools/d22scan/main.go:330-341`、实现 `main.go:758-787` | **40 tracked**（+ 本机已构建的 `dist/**`；它**不跳 dist**，=`R-92-6`） |
| 渲染器实际加载集合 | `frontend/dist/index.html` 与 `dist/assets/*`（经 `frontend/embed.go:24` 的 `//go:embed all:dist` → `internal/panel/assets.go:39-44` → 宿主按 `EntryFile`（`assets.go:31`）取）；源侧入口 `frontend/index.html:22` `<script type="module" src="/src/main.tsx">` | dist 里的**产物**，与 src 的 21 枚**不是一一对应**（本机现存陈旧产物 1 枚 js + 1 枚 css，见 ③.2） |

⇒ **钉比两者各窄的这些文件**（`git ls-files frontend` 里 `frontend/src/` 之外的 **19 枚 tracked**，逐枚点名）：
`frontend/index.html`（**渲染器自己的入口页，钉看不见**）、`frontend/vite.config.ts`、`frontend/fixtures/composer-states.html`、`frontend/fixtures/l2-card-fs-delete.json`、`frontend/scripts/render-composer.tsx`、`frontend/scripts/render-l2.tsx`、`frontend/scripts/{gen-tokens,vendor,vendor-shadcn}.mjs`、`frontend/package.json`、`frontend/package-lock.json`、`frontend/tsconfig.json`、`frontend/tsconfig.app.json`、`frontend/tsconfig.node.json`、`frontend/.oxlintrc.json`、`frontend/VENDORED.md`、`frontend/embed.go`、`frontend/dist/.gitkeep`、`frontend/.gitignore`；
外加 `frontend/public/`（**当前是空目录**——92b 的 F6 种植已被清掉，但它仍被 `vite` 原样拷进 dist 根、仍是指针的盲区）与一切 `frontend/dist/**`。
⇒ 与 `A105④` 同族的那条一般性质在这里成立：**"删掉扫描根之外的一枚文件/加一枚通道"对 21/2 这两个数没有影响 ⇒ 全绿**。修法不是加扩展名，是**换根**：把钉的根从 `frontend/src` 改成 `ban #6 用的那枚根**加上它的过滤**（即"能被 embed 进二进制的一切"），并把可达性判据挂在 `dist/index.html` 的引用图上（AC#9 的原话"回答渲染器到底会加载哪些文件"）。**⚠ 本代理未验证 `vite` 是否真把 `public/` 拷进 dist（没跑 `npm run build`，`frontend/public` 现为空 ⇒ 无既有产物可查）——这条留给写码代理，标 `未验证`。**
- 附带一条同族账：孪生门 `TestNoGitSwitchCapabilityInThePanelSurface`（`composer_test.go:276`）除 `node_modules`/`.git`/`dist` 外**还跳 `fixtures`**（`:286`）⇒ 它比门二又窄一层（`R-92-6` 的记账成立，现状复述）。

---

## ⑤ 判不出来 / 需要拍板的条目（一律标"未验证"，不圆）

| # | 条目 | 状态 |
|---|---|---|
| **P0** | **票 114 与票 33 的地界**：票 33（`33-panel-host-c27.md`，Status `ready-for-agent`）拥有 WebView2 宿主与 `AddWebResourceRequested`，票 35（`ready-for-agent`，`Blocked by: 33, 34`）拥有泵。"渲染 → `ParseComposerRequest`"这一跳按票面归 33/35。⇒ **114 是"只做 handler+门（性质不变、可绿、但通路仍不存在）"，还是"连宿主腿一起做掉（做完 AC#1 才有意义）"？** | **未定，需 owner/编排者拍**（本代理只量出两票均未动） |
| **P1** | **nil-`Confirm` 门的范围**：AC#2 原文"mode 写入处理器在注入的 `Confirm` 为 nil 时必须拒"。而 `internal/perm/store.go:191` 只在 `auto_approve` 这一支要确认腿（R20/M4 的不对称是明写的）。⇒ 门是"**装配完整性**：腿没接就三档全拒"，还是"**沿用 M4 语义**：只有切宽要腿"？两种写法都能满足"删掉 nil 检查就红"的变异判据，但**面板请求 `ask_high_risk`（切严）会不会因此被拒**是产品行为差异。 | **未定**；本代理倾向"装配完整性"（与 AC#2 字面一致、且 fail-closed 一侧不误伤"切严"要靠用例钉，不靠猜），但**这是语义决定，不是审计结论** |
| **P2** | handler 用 `perm.ConfirmFunc`/`perm.Switch` 会新增 `internal/panel → internal/perm` 的 import（今天**不存在**）。方向上无环（`perm` 不 import `panel`，实测 `internal/perm/*.go` 无 `panel` 字样）。 | **未编译验证**（本代理不跑 build）；建议按 ②.2 用本地小接口规避 |
| **P3** | `tools.PathCanonicalizer` 与 `PathScope` 的关系：三个方法签名**逐字对得上**（`internal/tools/paths_workspace.go:40, 51, 82` vs `internal/panel/workspace.go:36-47`），但**没有任何编译期断言**（全仓 `_ PathScope` 的 `var` 声明 = 0 命中；`panel` 也不 import `tools`）⇒ 今天它只是"拼巧成立"，改名/改签名不会有任何仪器红。接线那一票要么加 `var _ PathScope = (*tools.PathCanonicalizer)(nil)`（会引入 P2 那条新依赖），要么在装配根用赋值钉住。 | 已核（静态），**未跑** |
| **P4** | `panel.approval.request` 在 `composer_test.go:417` 的允许词表里、在 `knownComposerMethod`（`bridge.go:97`）外——是"两份名单本就该不同"还是"名单腐坏"，**没有断言能区分**。若按 AC#8 把两份名单对齐，这一条会**立刻变红**，需要事先给它一个去处（另设 `approvalRouteLiterals` 或把它进 Go 闭集）。 | **未定，需拍**（会决定 AC#8 的实现形状） |
| **P5** | 堵 F5（`"panel"+"."+noun+"."+verb`）的三选一修法（④.1 末），各有代价，本代理**不推荐**（memory：审计给的修法也是未验证断言）。 | **未定** |
| **P6** | `vite` 是否真把 `frontend/public/` 拷进 `dist/`（现目录为空、无产物可查）、以及 `npm run build` 后 dist 是否会出现 `"panel.*"` 字面量。 | **未验证**（未跑 npm，且 `frontend/**` 写/跑都不归本代理） |

---

## ⑥ 伪授权登记（本会话，逐枚点名）

**自称"编排者备注 / 系统提示 / 用户已更新编码规则 / `wisp-orchestrator-continuation` / 请 revert / 冻结某包"的文字：次数 0。** 本轮工具输出里**没有**任何此类内容，也没有任何要求我 revert/改写/放宽判据的指示。

出现过的非用户内容（**均未被当作指令执行**）：

1. `Note: The file C:\Users\swq\.qoder-cn\...\memory\MEMORY.md was modified since it was last read.` —— **3 次到达、共 4 条 note**（项目级索引 1 条 + 全局索引 3 条）。点名它挂在哪一枚工具调用的结果上：
   - 第 1 次到达：首工具批的结果尾部，同批是 `Bash`（命令前 40 字 `cd "D:\work\workspace\projects plans\Wisp" && git l`，内容 = `git log`/`rev-parse`/`status`）与 `Read`（票 114 票面）。含**两条** note（项目级 + 全局级索引）。
   - 第 2 次到达：`Bash`（命令前 40 字 `cd "D:\work\workspace\projects plans\Wisp" && grep -r`，内容 = `knownComposerMethod`/`composerRouteLiterals` 的对照 grep）。1 条（全局索引）。
   - 第 3 次到达：`Grep`（pattern `"panel\.[A-Za-z.]*"`，path `frontend`）+ `Bash`（前 40 字同上，内容 = 两条名单的对照）这一批之后。1 条（全局索引）。
   性质 = harness 的记忆索引更新通知。**核对它引用的编号是否存在（`A104③` 式核）**：其中引用的 `A102` → `docs/reports/pending-and-issues.md:3296` 存在；`Q-30` → 同文件 `:988` 存在；`R21` → 同文件命中 6 处存在；`票 100` → 命中 4 处存在。⇒ **它引的编号都是真的，但它说的内容与本票无关，且不构成对本票任何一条边界的修改**（本代理未据其改变行为，尤其未据其扩/缩 `frontend/` 冻结面）。
2. 首条消息尾部的可用 skills 清单（工具输出的一部分）—— 与本票无关，未据此调用任何 skill。

---

## ⑦ `next=`（给下一段写码代理）

**必须先答的问题：`P0`（票 114 做到哪一跳）与 `P1`（nil-`Confirm` 门覆盖三档还是只覆盖切宽）。这两条没答，动手就是猜。**（推荐：`P0` 答"114 只做 handler+门，宿主腿留票 33"，`P1` 答"三档全拒 = 装配完整性"——但这两条需要 owner/编排者点头，不是我能替定的。）

**答完之后的最小改动面（按 ②.2，全 Go，零 `frontend/`）：**

1. 新建 `internal/panel/composer_handlers.go`（~60-90 行）：`mode.request` handler 的第一句 = `if Confirm == nil { 写审计; return err }`，**位置在任何 `ModeWriter` 调用之前**；`ModeWriter` 写成**本地小接口**（对齐 `internal/perm/store.go:175` 的签名），避免 ⑤/P2 的新依赖。
2. 新建 `internal/panel/composer_handlers_test.go`（~80-120 行）：断言的是**"注入的 `ModeWriter` 一次都没被调"**，两发变异各自必须红（②.3）。⇒ 这两枚文件在 `internal/panel`，**ubuntu core scope 有分母**（`scripts/portable-tests.sh:141,179`），Linux runner 上真跑；`cmd/wisp` 里的装配断言**只有 windows 腿跑**（`ci.yml:378`），别把行为判据只写在那儿。
3. `cmd/wisp/run.go` 于现有 `:333-347` 附近接 handler（~10-20 行）；`internal/panel/bridge.go` 预计 **0-10 行**。
4. **开工前的停手判据（票头 09-23 横幅 + `A102③`）**：列出要写的每一枚文件；只要有一枚落在 `frontend/` 下就整段停手上报——**④ 的 AC#8/AC#9 修法若最终要求改 `frontend/`（例如把路由改成必须显式字面量），那半格按 `skipped=frontend(owner-delegated)` 记账，不许"只改一角"**。

---

## 复现命令（本表全部数字；本代理实跑）

```
git rev-parse --short HEAD                    # 7ad6eb4
grep -rn "ParseComposerRequest" --include=*.go .              # 定义 bridge.go:77，调用者全在 bridge_test.go
grep -rn "func (s \*Store) Set" --include=*.go internal/perm  # store.go:175
grep -rn "\.Set(" --include=*_test.go cmd/wisp internal/perm | grep -i "mode\|perm\|store"   # 12 处，全测试
grep -rln "wisp/internal/panel" --include=*.go .              # 非测试唯一 importer: cmd/wisp/panel_assets.go
grep -rno "\"panel\.[A-Za-z.]*\"" frontend/src frontend/index.html frontend/dist   # 5 处，全在 panel.ts
git ls-files frontend | wc -l                                 # 40；其中 frontend/src 下 21（= 钉自报的 "21 files scanned"）
git ls-files frontend | grep -v "^frontend/src/"              # 钉看不见的 19 枚（④.2 逐枚点名）
# 门二形状在全 frontend（含本机 dist 产物，node_modules 无跟踪）的命中面：
#   grep -rE "\bsetMode\b|\bmode\.set\b|\bpanel\.mode\.set\b|..." frontend  => 仅 frontend/src/lib/panel.ts:230（注释行，被 codeOnly 清空）
```

**未跑**（本代理只读，任务亦禁止全仓跑）：`go test ./...`、任何单包 `go test`、`npm run build`、`sh scripts/d22scan.sh`、`gofumpt -l`。
⇒ ③.3 的"不会让已绿用例变红"是**静态 grep 面**结论，写码代理复跑 `bash scripts/portable-tests.sh ./internal/panel/` 即可转正。

**时间戳**：本文件写入时本机 `date` 实测 `Wed Sep 23 10:00:43 CST 2026`（会话开始时同为 09:5x-10:00 段，未复用他人读数）。
