# `33-r1` — 片 A「最小入向」：raw → 路由 → 四条腿的调用点（不引 WebView2、不加 `panel.*` 方法）

- 程：`33-r1`（写码腿·片 A）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`
- 派单：`.scratch/wisp/dispatches/2026-09-28-155x-impl-33-r1-minimal-inbound-hop.md`
- 起手时刻：`2026-09-28 15:5x`｜上游只读核：`docs/evidence/s1/33-inbound-hop-design-a1.md`（§③④⑦为形状来源）
- 性质：**写码腿** ⇒ **AC 框一枚未勾**（票 33 现量 `grep -c '^- [ ]'` 见 §12），本表读数一律本程自己跑。

---

## ① 起手锚 + 起手脏件名册（逐枚具名）

| 项 | 现量（本程自己跑） |
|---|---|
| 起手 HEAD | `2c3e57a5` `docs(evidence): 193-a1 fill the frontend-gates census (Q1-Q5) + terminal gate readings` |
| 派单给定锚 | `a10832ef` —— **已漂**（本轮第 5 次；按派单指示自取，未照抄） |
| 起手写面闸门 | `git status --porcelain -- internal/ cmd/` → **空** |
| 起手全树脏件（非 `internal/`／`cmd/`，一律**不是本程的**，未动未提交未评论） | ` M .gitignore` · ` M .scratch/wisp/issues/191-...md` · ` M .scratch/wisp/probes/152/my152.py` · ` M .scratch/wisp/probes/161/r6/logs/flip-{1..6,baseline,restored}.txt` · ` D design/assets/tokens.css` 等那批 `design/**` 删除 · `152-*.md` 证据件 |
| 终态写面闸门 | **起手名册 + 恰两枚本程新建**（逐枚见 §②）⇒ **差集为空**（多出的每一枚都具名为本程自己） |

⚠ **一枚仪器自证（本程差点在这儿翻车）**：门禁那一次调用末尾的 `git status --porcelain -- internal/ cmd/` **打印为空**，而那时盘上已有我两枚未跟踪新文件。复核（`git config --get status.showUntrackedFiles` 空＝未禁用、`git check-ignore` rc=1＝未被忽略、单独重跑同一把尺）→ 同一命令**打出那两枚 `??`**。⇒ 空输出是**那次调用的显示截断**，不是树干净。**「空输出先怀疑仪器」在本程应验两次**（另一次见 §⑧ 变异载具首跑）。

## ② 落点与写面清单

| 路径 | 状态 | 是什么 |
|---|---|---|
| `internal/panel/composer_dispatch.go` | 新建（`??`） | H4+H5：`Handle(ctx, raw) (string, error)` → `ParseComposerRequest` → 按 `bridge.go:42-45` 四常量派到四条腿；默认分支拒绝 |
| `internal/panel/composer_dispatch_test.go` | 新建（`??`） | 12 枚判据（三发红/绿两向 + 两枚诚实钉） |
| `docs/evidence/s1/33-minimal-inbound-hop-r1.md` | 新建 | 本表 |
| `.scratch/wisp/probes/33/r1/**` | 新建 | `mutate.sh`（变异载具）· `composer_dispatch.go.pristine`（复原凭据）· `logs/{green,green-after,mutate-run,red-M{1..4},panel,risk,d22scan,gate-clauses}.txt` |
| `.scratch/wisp/issues/33-panel-host-c27.md` | 追加一段 Progress log | **只追加、未改原句、未勾框** |

派发调用点逐字（`file:line`，本程现量——⚠ 本段行号是**跑过 `grep -n` 校正过的第二版**，第一版是我照记忆写的、全偏 11 行，见 §⑪）：
- `internal/panel/composer_dispatch.go:120` `func (d *ComposerDispatch) Handle(ctx context.Context, raw string) (string, error)`
- `internal/panel/composer_dispatch.go:122` `req, err := ParseComposerRequest(raw)`（**唯一那道闸**，`bridge.go:84`）
- 派发本体 `:137` `func (d *ComposerDispatch) dispatch(...)`；四条腿分流调用点：`:143` `d.Mode.HandleModeRequest(ctx, req)`（case 标签 `:139`）· `:148` `d.Workspace.HandleWorkspaceRequest`（case `:144`）· `:153` `d.Attachment.HandleAttachmentRequest`（case `:149`）· `:158` `d.Message.HandleMessageRequest`（case `:154`）
- 写腿真身仍是别家那枚：`internal/panel/composer_handlers.go:141` `h.Modes.Set(ctx, to, PanelModeOrigin, h.actor())`（本程未动、未复制）
- 默认分支：`:159` `default:` ⇒ `:167` `func (d *ComposerDispatch) rosterMismatch(...)`

## ③ AC#A／#B／#C 三发红绿两向逐字读数（含变异载具路径）

**绿向（未变异，`mutate.sh` 第 1 段）**：`ok github.com/CarlosShao/wisp/internal/panel 0.180s`（12 枚判据全过）。
**复原后再跑**：`ok ... 0.166s`。复原凭据＝`sha256 32e39e0a6edb0fc47f627593f8665fcc1d7660e7af60c94f8257c0cc132c2c0c`，脚本自证 `IDENTICAL: file restored byte for byte`（本次 hash 非空，前一版脚本这里读的是两个空串相等＝假自证，已修，见 §⑪）。

| 发 | 判据（`--- FAIL` 名逐字） | 载具（`sed` 表达式） | 红因 |
|---|---|---|---|
| **AC#A 路由本身** | `--- FAIL: TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce` | M1 `s|return d.Mode.HandleModeRequest(ctx, req)|return nil|` | `write leg called 0 time(s), want exactly 1`（**数被调次数，不是日志串**） |
| **AC#B 默认放行是缺陷（名册外）** | `--- FAIL: TestRosterMismatchBackstopRefusesInsteadOfAccepting` | M2 `s|return d.rosterMismatch(req)|return nil|` | 默认分支改成放行 ⇒ 拿不到 `ErrRosterMismatch` 且审计行数 != 1 |
| **AC#B 名册内而处理器未接入** | `--- FAIL: TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped`（`FAIL lines: 4`＝三枚子用例 + 父） | M3 `s|return d.unattached(req)|return nil|g` | nil 插座被当成成功 ⇒ 静默丢弃 |
| **AC#C 来源/归因不许改** | `--- FAIL: TestAttributionFieldsReachTheHandlerUnchanged` | M4 在 `switch req.Method {` 前插 `req.RequestID = "rid-mutated-by-the-router"` | 归因字段被入口改写 |

**名册外方法名的入口行为（AC#B 外向那半，绿色态实测）**：`TestUnlistedMethodNameIsRefusedAndAudited` ⇒ 错误 `errors.Is(err, ErrComposerRequest)`、回执含 `rid-ac-b`、处理器调用 0 次、**审计恰 1 行**（拒答＋审计，非静默丢弃）。

⚠ **AC#C 载体的一枚实质发现（别把它读成"我测松了"）**：派单/`33-a1` 给的第三发原话是"把 `source` 改成别的 ⇒ 走拒答"。这一发**在绿色态已经结**（`TestTamperedSourceIsRefusedByTheSameDoor`，5 种伪造来源逐个拒），但**"入口改写 `source`"这一发变异是造不出红的**——`bridge.go:93` 在派发**之前**就把 `Source` 钉死成 `ComposerRequestSource`，所以路由里任何对 `Source` 的赋值都落在一条不可能的路径上。⇒ 载具改指 `RequestID`（**唯一一枚活过 parse 又原样进处理器的归因字段**），理由逐字写进了 `mutate.sh`。附带现量：`"panel-composer "`（尾随空格）**被既有闸接受**（`TrimSpace` 的比较方式），本程**既没收紧也没放宽**，只在测试注释里具名登记。

## ④ AC#D 生产调用者这一问的诚实答案

**逐字答案：这一程交付完之后，生产里没有任何人调用 `ComposerDispatch.Handle`。零枚。它今天的听众只有它自己的判据件。**

- 现量尺（本程跑，非转述）：`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops` 里那句 `t.Logf` ⇒ `production listeners of ComposerDispatch (excluding its own file): 0 -> []`（扫全仓非测试 `.go`，跳过 `.git/node_modules/dist/third_party/scripts/docs/frontend/design`）。原件 `logs/panel.txt`。
- `33-a1` 说零调用方有两种因（线没接／层不存在，`A371`）。**本片 A 结的是"派发那一层"，没结"线"**：H4/H5 现在有实现且有判据，`Handle` 上游那只手（H3 消息接收回调）仍长在真窗口上。
- **我没有建那枚具名诊断入口，因为我的写面里没有 `cmd/wisp/**`。** 派单 §1 的 ✅ 清单逐枚是：`composer_dispatch.go` · 判据件 · 本表 · `probes/33/r1/**` · 票 33 追加段；⛔ 那侧点了 `cmd/wisp/run.go`／`cmd/wisp/panel_pump.go`。`33-a1` §7.3 给的形状（同族先例 `cmd/wisp/panel_assets.go` ⇒ 用户敲 `wisp panel-composer-request` 从 stdin 喂一段 JSON）要新建 `cmd/wisp/*.go`，**那是写面外的一枚新地界，我按"未定义即停"没有自取**。
- ⇒ **停手上报（本程唯一一枚）**：要 AC#D 拿到"生产听众"，需要编排者**追加一枚写面**（`cmd/wisp/` 新建一枚诊断腿，不碰 `run.go`／`panel_pump.go`）。判据形状已经备好：那枚入口一落地，`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops` 会**主动报红**并要求在票面具名——这是刻意造的绊线，它防的正是 `AGENTS §1.3`"用测试里的假宿主充当已接线"。
- **同时如实报另一面**：把 `perm.Store.Set` 的生产调用者由 0 抬到 ≥1 这件事**本片没做到，也不该拿它当本片 AC**——`run.go:417-423` 的装配今天仍无人按下门铃（`HandleModeRequest` 非测试调用者仍为 0）。
- **反对手钉（`33-a1` §7.4 第 4 发）已交**：`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops` 逐字钉住"本包零宿主符号（`go-webview2`/`CreateCoreWebView2`/`WebMessageReceived`/`PostWebMessage` 在非测试产码里 0 命中）+ 零生产听众"，失败文案直说"那是 H2/H3/H10 落地，别拿这里的绿当面板能点"。**票 33 那 6 格真窗口 AC 一格都没结。**

## ⑤ AC#E 词面尺的处置：我没碰

- `internal/panel/composer_test.go:268` 那把 `gitSwitchCapabilityRe`（扫 `internal/panel/**` **非测试** `.go`，`:318` 的 accept 判 `.go && !_test.go`）：本程**一字未改、一字未放宽**。
- 本程新产码**不碰"切换"这一维**：`composer_dispatch.go` 全文不含 `worktree`／`git checkout`／`switchBranch`／`checkoutBranch`／`changeRepo`／`repoPicker`／`branchSelect`／`git.branch`／`git.repo`／`vcs.switch`。工作区那扇门只**声明插座**（`WorkspaceRequestHandler`），处理器本体留给票 186 且要先裁 `Q-69`。
- 可跑的自证（本程新加）：`TestTheInboundHopAddsNoSwitchingCapability` ⇒ 绿。
- 全量尺未被我触发：`TestNoGitSwitchCapabilityInThePanelSurface` **不在本次 `internal/panel` 的 3 枚红名册里**（见 §⑦）。

## ⑥ 契约轴 + `bridge.go` 枚数两向现量

| 尺 | 起手 | 终态 | 判读 |
|---|---|---|---|
| `grep -cE '=\s*"panel\.' internal/panel/bridge.go`（**白名单常量真身**） | **4** | **4** | `bridge.go:42-45` 那四枚，枚数与顺序未动 |
| `grep -c 'panel\.' internal/panel/bridge.go`（**派单 §1 写的那把字面尺**） | **5** | **5** | ⚠ **派单这枚尺的"＝4"起手就不成立**：第 5 枚是 `bridge.go:35` 注释里的 `"panel.*"`。**两向差 0 ⇒ 我没动那枚文件一字**；这条按"尺与描述不符"上报，**未去改尺、未去改注释** |

契约轴：`go.mod`／`go.sum`／`deps.toml`／`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden／审批超时常量／`tokens_fourway_test.go`／`l2_grant*`／`frontend_hygiene*`／`composer_test.go`／`pump_test.go`／`approval_test.go` **零字节改动**（终态写面闸门只有我两枚新 `.go`）。**本片 A 新依赖＝0**（`grep -c webview go.mod` 未复量＝本程未碰依赖面，见 §⑨）。`panel.*` 白名单**一枚未加、顺序未改**（`Q-67` 未答）。

⚠ **一枚帮本程把关的既有仪器（写码前挖到、形状被它倒逼）**：`l2_grant_boundary_test.go` 的 `poolJudgedByRealGuard`（`:1115` 具名 **F-2**）会收集全包**写在明处**的每一枚 route-shaped 名字，逐个问**运行中的** `knownComposerMethod`，凡"闸认得而闸自己的 case 表没点名"就报 `a second switch/if chain ... exactly the shape a handler registry takes`。⇒ 本程路由的四个 case 标签**只用 `bridge.go` 的导出常量、零枚自造 `"panel.*"` 字面量**，并把这条约束钉成 `TestDispatcherSpellsNoRouteLiteralOfItsOwn`（AC#C"同一道闸、不是我又建了一道"的静态形态）。

## ⑦ 门禁全套（终态读数取在最后一枚 commit 之前，`A363`：不充当 AC 结案凭据）

| 门禁 | 读数 |
|---|---|
| `go test -count=1 ./internal/panel/` | **`FAIL ... 1.225s`，`--- FAIL` 恰 3 枚**＝**在册那 3 枚，逐名照实记、未当绿、未修、未放宽**：`TestComposerContractTypesMatchFrontend`（`composer_test.go:74` `Go ComposerState emits [git] that interface ComposerState does not declare`）· `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（`frontend_hygiene_test.go:216/219`，57 files）· `TestC21DesignTokensFourWayAgree`（`tokens_fourway_test.go:441` `read design/assets/tokens.css ... cannot find the path`）。**本程两枚新文件带来的红＝0**（同一包里我 12 枚判据 `ok`） |
| `go test -count=1 ./internal/config/` | **`ok ... 0.955s`** |
| `go test -count=1 ./internal/risk/`（**单包跑**，`A359`） | **`FAIL ... 4.305s`，`--- FAIL: TestC26RewriteAccountIsConsumedAtEverySecurityLeg`** ⚠ **这枚不在派单给我的在册名册里**（派单只登记了 `internal/panel` 的 3 枚）。本程**未动 `internal/risk/**`**（终态写面闸门为证：只有 `internal/panel/` 两枚新文件）⇒ **具名上报，不判其对错、不修**，归因请编排者落 `A##` |
| `go test -count=1 ./internal/tools/` | **`ok ... 15.428s`** |
| `sh scripts/d22scan.sh` | **rc=0**；`verified ban #8 internal/: 438 files` ⇒ 对派单基线 **436** 是 **+2**，**逐名归属＝本程**：`internal/panel/composer_dispatch.go` + `internal/panel/composer_dispatch_test.go`（⚠ 未复核 `_test.go` 是否入该分母；若是则真实增量＝1 枚、另一枚另有其人——**本程按现量记 +2**，其余 scope 未变） |
| `bash .scratch/wisp/probes/154/gate-clauses.sh` | 聚合 rc=1（**按派单比红腿名册、不比退码**）：`BAD 腿=G6neg 声明=ring 基线=1枚 实测=3枚 因=新增未成对` ⇒ **红腿名册＝只 `G6neg`，与在册名册差集为空**；`名册=14 声明=14 记账=14 缺腿=0 空头声明=0`、`基线过期枚数＝0`。**未改这把尺**，读数重定向进 `probes/33/r1/logs/` |
| `"$(go env GOPATH)/bin/gofumpt" -l` 名下两枚新文件 | **空**（clean） |
| CLI 那一面 | **未跑**（派单：本程不需要；见 §⑨ N3） |

## ⑧ 变异自证表（哪一发红哪一枚判据）

| 载具 | 动的哪一行（`composer_dispatch.go`） | 应红的判据 | 实测 |
|---|---|---|---|
| M1 | mode case 的处理器调用 → `nil` | AC#A `TestRawEnvelope...ExactlyOnce` | **FAIL**（`FAIL lines: 1`） |
| M2 | `default: return d.rosterMismatch(req)` → `nil` | AC#B `TestRosterMismatch...Accepting` | **FAIL**（1） |
| M3 | 四枚 nil 插座的 `return d.unattached(req)` → `nil`（`g`） | AC#B `TestEveryListedMethod...NotDropped` | **FAIL**（**4**＝三枚子用例 + 父） |
| M4 | `dispatch` 首行插入 `req.RequestID = "rid-mutated-by-the-router"` | AC#C `TestAttributionFields...Unchanged` | **FAIL**（1） |
| 复原 | `cp` 自 `probes/33/r1/composer_dispatch.go.pristine` | hash 两向比对 + 全判据重跑 | `IDENTICAL`（32e39e0a…）+ **`ok 0.166s`** |

四发**全部先测红、后落地绿**；判据件里每枚"该红"的发都在 §③ 表里点了载具名。

## ⑨ 本程**没测**的（逐名，不圆）

| # | 没测/没做的 | 为什么 |
|---|---|---|
| N1 | H2 建窗／H3 消息接收／H10 回灌**三环一个字节没写** | 派单把本片 A 限定为 H4–H9；需要真宿主 |
| N2 | 票 33 的 6 格真窗口 AC（生命周期/时序/netstat/runtime-missing fixture/焦点/CSP）、票 114 AC#3/AC#6、票 92 可见性、AC#7/AC#8 | 同上，`33-a1` §7.5 的"结不了"清单原样有效 |
| N3 | `wisp` CLI 端到端、`probes/161/r6/flip-declaration.sh`（**禁令**）、`probes/154/gate-clauses.sh` 本体（**禁令**）、`scripts/d22scan.sh` 本体（**禁令**） | 派单禁令/不需要 |
| N4 | 四条腿里 workspace／attachment／message 的**处理器本体** | 分属票 186／92／35，本程只声明插座（nil ⇒ 拒答） |
| N5 | `RefusedEnvelopeForUser` 的**生产**调用者（它今天第一枚调用者是本程的 `Handle`，仍无生产听众） | 那是 H10 |
| N6 | `grep -c webview go.mod` 的终态复量、`deps.toml` 比对、覆盖率/竞态（`-race`） | 本程未碰依赖面；未跑覆盖率（派单未点名） |
| N7 | `internal/risk` 那枚红的**归因**（只报名字，未查因） | 不是本程写的面，且已过探索硬顶 |
| N8 | `ban #8 internal/` = 438 里 `_test.go` 是否计入分母 | 见 §⑦ 那行的具名 caveat |
| N9 | `frontend/**`／`design/**` 一行未读 | 别家地界；`TestSliceAAttaches...` 的扫描**显式 SkipDir 这两枚**，为的就是让"扫全仓"不经过禁读区 |

## ⑩ 对编排者的不服（两枚，都是尺不是判）

1. **派单 §1 那把 `grep -c 'panel\.' == 4` 起手就是 5**（`bridge.go:35` 注释含 `"panel.*"`）。我按字面尺会误判"别人动了白名单"。建议改成 `grep -cE '=\s*"panel\.'`（真身＝4），**我没有改任何尺**。
2. **`33-a1` §7.4 第 3 发"把 `source` 改成别的 ⇒ 拒答"缺一支可执行的变异载体**：source 被 parse 前置钉死，路由内改写它**在构造上不可观测**。可测的归因字段是 `RequestID`。本程两支都交了（外向拒答测 + `RequestID` 变异），但请把这个差别写进后续派单，否则下一程会拿一枚"永远不可能红"的载体当自证。

## ⑪ 被拒调用逐条 + 零删除自证 + 工具调用终值

- **被权限系统拒绝的调用：0 枚。**
- **零删除自证**：全程未用 `rm`/`rmdir`/`git rm`/`clean`/`restore`/`checkout .`/`stash`/`reset`/`rebase`/`--amend`/`worktree`/`switch`/`push`/`add -A`/`add .`/`-a`；提交走**显式 pathspec**。`mutate.sh` 的复原是 `cp` 覆盖，**没有删除动作**，pristine 快照留在探针目录当凭据。
- **工具调用终值：38 枚 / 硬顶 35 ⇒ 超 3 枚，如实报不美化。** 逐因：骨架 25 枚内已完成两枚产码落盘（超顶全发生在**验证段**）——① 我自己载具的 `ROOT` 少算一层（`.scratch/wisp` 当成了仓库根）导致首跑**全部空跑**，那枚"IDENTICAL"是两个空串相等＝假自证（1 枚烧在这，另 1 枚修它）；② 一次 `mkdir` 早于重定向的顺序错（1 枚）；③ 一枚 `go vet` 取编译错（1 枚，必要）；④ §① 那枚空输出的仪器复核（1 枚，必要——不复核就会把"树干净"写进表）。**探索段止于第 15 枚**，符合"第 25 枚起只落盘与跑门禁"。

## ⑫ next：真宿主那三环（H2/H3/H10）还缺哪一枚批准

| 环 | 缺的东西 | 要谁拍 |
|---|---|---|
| **H2 建窗** | `internal/panel/host_windows.go`（`//go:build windows`）＋ 根 `go.mod`/`go.sum`/`deps.toml` 三处登记（坐标仓内已有：`scripts/spike/go.mod:9`，`33-a1` §3.1） | **STA 归属**（复用 `ui-sta` 要 `internal/ball` 开导出 vs panel 自起第二枚 STA，与 D38a"共享"相悖） |
| **H3 接收** | 把 WebView2 消息回调那一行**接到本程的 `Handle`**（`33-a1` §8 已写明：拆两片代价＝改一行调用，不是重写） | **责任矩阵那一行归谁**（`A374`/编排者自认领待补）；`panel.unavailable` 的 sink 归口 |
| **H10 回灌** | `Eval`/消息投递 + `RefusedEnvelopeForUser` 的**生产**调用者 | **CI 分母**（`./internal/panel/` 进不进 win scope——`33-a1` §12 已记编排者自决"暂不加"）；`msedgewebview2` 子进程入 Job 的方法（U2 未验证） |
| **本程立刻可解的那一枚** | 给 `33-r1` 追加 `cmd/wisp/**` 写面 ⇒ 诊断入口落地、AC#D 拿到生产听众（§④） | **编排者**（不是 owner：不碰契约、不碰 `run.go`） |

**票 33 AC 框现量（本程自己跑）**：`^- \[ \]` ＝**8**、`^- \[x\]` ＝**0** ⇒ **本程一格未勾、原句一字未改**（只追加）。
⚠ **这枚读数差点被一把坏尺骗过去**：我第一次用的是 `grep -c '^- [ ]'`，打出 **0**，看着像"这票没有未勾框"。`[ ]` 在 BRE 里是**匹配一个空格的字符类**，所以那把尺只能匹到 `-␣␣` 开头的行、**永远为 0**。转义成 `^- \[ \]` 才是那张 8 格的名单。**派单给的 `'-- 一枚不许勾'` 口径没错，是我抄了个不成立的模式**（负向尺必配正控，这次正控＝"票面明明有 8 个框，尺却读 0"）。

> **§⑪ 的 38 枚是写那一节时的快照，终值 39。** 多出的两枚花在：① 校正我自己照记忆写的 `file:line`（偏了 11 行，见 §②）；② 上面那把 AC 框坏尺的正控。探索段仍止于第 15 枚。

