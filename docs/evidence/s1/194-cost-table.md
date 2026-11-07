# `194-cost-table` — 面板方法名册对齐：逐枚代价表（AC#1 只读步骤，2026-11-07 现量）

- **程**：`194-cost-table`（只读，零产码）。写面＝本件＋票 194 Progress log 追加一行。
- **工单**：`.scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md`（**本件不勾它的任何 AC 框**；翻框归非实现者裁决腿）。
- **起手锚**：HEAD＝`a1c3f245b`（`git log --oneline -1` 现读，取数时刻 2026-11-07 16:5x–17:1x +08）。
- **本件不是 `194-c1` 的续写**：c1 普查（09-28，锚 `8ca2291f`）交的是**全量普查**（AC#1 原文四问＋分堆＋待裁清单，226 行）。本件只做派单点名的**代价表＋裁定句**这一步，四问的普查面以 c1 与 `A384` 为准；本件全部读数是**今天这棵树**的现量，与 09-28 快照的差异**逐枚具名**（§0）。

## §0 票面/快照 → 今天的树：现状更正（全部现跑）

| # | 票面/旧读数 | 今天现量 | 尺（可复跑） |
|---|---|---|---|
| 0-1 | 票面 `:12`：代码名册＝`bridge.go:42-45` **4 枚** | **7 枚**：原 4 枚（`:42-45`）＋ 票 248 加的 `config.get`/`config.set`（`:66-67`，无 `panel.` 前缀是有意设计）＋ 票 320 加的 `config.models.list`（`:75`） | `sed -n '41,76p' internal/panel/bridge.go` |
| 0-2 | c1：`ComposerDispatch.Handle` 生产调用者＝**零枚** | **2 枚**：`cmd/wisp/panel_host_windows.go:955` ＋ `cmd/wisp/panel_inbound.go:173`（7b2 节 §② 量过的那把尺，今天同一读数） | `git --no-pager grep -n 'disp\.Handle(' HEAD -- cmd internal \| grep -v _test.go` |
| 0-3 | c1 §2.3：规格 18 枚 ∩ 代码 4 枚＝**零枚** | 交集＝**3 枚**（`config.get`/`config.set`/`models.list`，见 §1 注甲） | 本件 §1 的 A2 扫描 |
| 0-4 | c1：`models.list`＝B 堆「没货」；`config.set`＝C 堆「地基不存在」 | **两枚都落地了**：`config.set` 走 195 号票地基 `Manager.WriteSection`（`internal/config/panelwrite.go:168`——**文件名是 `panelwrite.go`**，§1 注甲的精读更正）；`models.list` 的规格拼写**有意不实现**，今天等价物是 `config.models.list`（票 320） | `grep -nE "^func " internal/config/panelwrite.go` |
| 0-5 | c1：`approval.queue` 同名噪声落点 `internal/statemachine/events.go:43` | 今天 `:43` 是 **`approval.queue-drained`**（子串命中，**非独立字面量**）——噪声仍在，但落点应表述为「子串命中 `:43`」 | `sed -n '43p' internal/statemachine/events.go` |
| 0-6 | c1 §4：出向管子＝`publishPanelSnapshot` 多处调用 | 生产调用者今天＝**1 枚**（`cmd/wisp/resident_present_windows.go:197`）；`bookPanelSnapshot` 也在 `run.go:958`（33-r14 提交信息自述，本件未复量该行） | `git --no-pager grep -n "publishPanelSnapshot(" HEAD -- cmd/wisp \| grep -v _test.go \| grep -v "func "` |

## §1 代价表（18 枚 invoke ＋ 5 枚出向事件，今日现量）

列定义：**定义处**＝方法名或其能力的落码处（注甲）｜**调用者数**＝生产（非 `_test.go`）调用者枚数｜**具名用例**＝字面量命中（含测试侧；本轮 16 枚未落地名字测试侧**全零**）｜**改动面**＝落地最少要动的文件数｜**代价**：低＝追路由＋handler；中＝再加判据＋地基补充；高＝先要功能/契约/禁区批准（写腿自决不了）。

**注甲（同名命中逐枚精读）**：`config.get`/`config.set` 命中＝真实现（248 号票的 settings 门：`bridge.go` 常量＋`config_handlers.go` 的字段白名单＋`composer_dispatch.go:220/225` case＋`panel_inbound.go:311-315` 生产装配）；`models.list` 命中＝`config.models.list` 的**子串**（320 号票有意不实现规格拼写，见 `bridge.go:71-75` 注释「deliberately NOT spelled models.list」）＋两处注释；`approval.decide` 命中＝**防御注释**（`cmd/wisp/approval_reply.go:326`、`internal/agent/approval/replies.go:422` 的 panel-allow 拒绝路径）＋ d22scan 文案（c1 已判噪声）；`grants.list`/`grants.revoke` 命中＝`internal/session/grants.go:40` 的**范围排除注释**（票 224 明写这两枚不归它、revocation 入口＝226 的 `SetAllowedDirs` 路由）；`approval.queue` 命中＝状态机事件 `approval.queue-drained`（子串，c1 已判噪声）。`sectionwrite.go` 是我本轮一次转述笔误，正确文件名 `panelwrite.go`。

### 甲：已落地（4/18）

| 方法 | 定义处 | 调用者数 | 具名用例 | 改动面 | 代价 |
|---|---|---|---|---|---|
| `config.get` | `internal/panel/config_handlers.go`（248 号票，字段白名单只读侧） | 2（`panel_host_windows.go:955`、`panel_inbound.go:173`） | `TestAC1BothSettingsNamesAreWhitelistedAndRouted` 等（`config_route_248_test.go`）＋313 增量（`config_313_test.go`） | 0（已在册） | 已完成 |
| `config.set` | 同上＋地基 `Manager.WriteSection`（`internal/config/panelwrite.go:168`，195 号票） | 同上 | 同上＋`TestAC8ReceiptNamesTheTierAndNeverClaimsLiveSave` 等 | 0 | 已完成（安全节锁定＝195 地基内建） |
| `models.list` | **有意不落地**：320 号票以 `config.models.list` 实现同一能力（`cmd/wisp/panel_models_320.go:124` 起），规格拼写让位给本地已签名模型库命名空间 | （等价物）2 | （等价物）`Test*` 于 `cmd/wisp/panel_models_320_test.go`（本件未逐名列） | 0（或 1：若裁决要规格拼写，加一行别名常量即撞 320 注释——须先裁） | 已完成（拼写分歧待裁） |
| `panel.resync` | 语义落码＝出向 `NewSnapshot`（`internal/panel/composer.go`）＋`publishPanelSnapshot`（`cmd/wisp/panel_pump.go:491`）；**入向「请求重发」面不存在** | 1 | **零** | 2（常量＋case＋handler＋装配——`composer_dispatch.go` 与装配在 2 个文件，常量在 bridge.go，实际 3–4） | 中 |

### 乙：未落地·地基已在（9/18）

| 方法 | 定义处（最接近地基，今日 `file:line`） | 调用者数 | 具名用例 | 改动面 | 代价 |
|---|---|---|---|---|---|
| `tasks.list` | `ListTaskLogs`（`internal/memory/dao_tasklog.go:95`） | 3（crash_recovery.go:145、privacy.go:76、:214） | 零 | 2 | 低 |
| `task.detail` | `ListToolCallsByTask`（`internal/memory/dao_toolcall.go:172`） | **零** | 零 | 2 | 低 |
| `history.query` | 无查询面（c1 判 A/B 之间，今日无新面） | — | 零 | 2＋查询面 | 中 |
| `approval.current` | `Snapshot.Pending`（`internal/panel/composer.go`） | — | 零 | 1 | 低 |
| `approval.queue` | 同上 | — | 零 | 1 | 低 |
| `grants.list` | `ListGrants`（`internal/memory/dao_misc.go:90`）；⚠ 票 224 注释把这两枚挡在射程外（`internal/session/grants.go:39-44`） | **零** | 零 | 2 | **高**（先裁 435 banned-root 表口径） |
| `privacy.purge` | `PurgePrivacy`（`internal/memory/privacy.go:170`）等 7 枚 | **零** | 零 | 2 | 中（确认门＋审计） |
| `privacy.export` | `ExportPrivacy`（`internal/memory/privacy.go:198`） | **零** | 零 | 2 | 中（交付形态：artifacts 写入形状） |
| `cost.summary` | `ListCostDaily`（`internal/memory/dao_misc.go:169`） | **零** | 零 | 2 | 低 |
| `diagnostics.export` | `BuildDiagnosticsBundle`（`internal/observe/diagnostics.go:90`） | **零** | 零 | 2 | 中（交付路径＋脱敏） |

### 丙：未落地·后端没货（5/18）

| 方法 | 定义处 | 调用者数 | 具名用例 | 改动面 | 代价 |
|---|---|---|---|---|---|
| `transcript.get` | 无 transcript 存储概念（今日 grep 仅注释命中） | — | 零 | 3＋ | 高（先裁数据模型） |
| `models.delete` | 无删已装模型方法（`internal/models` 导出面仅 `UnnamedStreams`/`Ensure`/`VerifyInstalled`/`Cancel`/`VerifyDir`） | — | 零 | 3＋ | 高（不可逆＋PathResolver 规矩＋审计） |
| `approval.decide`（reject 支） | 后端真在：`PanelReject`（`cmd/wisp/approval_reply.go:304` 走 `PanelReject` 调用链）＋allow 侧防御注释 `:326` | PanelReject 生产调用者 **1**（`approval_reply.go:304`） | 零（字面量） | 2 | **高**（禁区①：allow 拒一切面板来源，`SPEC-08:167`；reject 支须绕开 146/Q-49 族判据，逐枚 A## 批准） |
| `grants.revoke` | `SetAllowedDirs` 收紧路由（票 226，`cmd/wisp/approval_revoke.go:84` 是其生产调用者） | 1 | 零 | 2 | **高**（禁区②口径已由 A384 裁定放行——但 224 注释与 435 表仍拦路，先对表再派） |
| `config.set` 安全节分支 | 195 地基 `lockedDirectionOf`（`panelwrite.go:343`） | — | — | 0 | 已完成（195 内建） |

### 丁：出向事件 5 枚（`SPEC-08:174`，Go→前端，另一根管）

| 事件 | 定义处 | 调用者数 | 具名用例 | 改动面 | 代价 |
|---|---|---|---|---|---|
| `task.delta` | 内容在 `Snapshot.Results`（整份快照重发，非 delta 帧） | — | 零 | — | 低（**只补出向名册**，⛔ 进 bridge.go 入向常量） |
| `tool.chip` | 真没内容（`Snapshot` 无 tool 段） | — | 零 | 2＋ | 高（先建数据面） |
| `approval.request` | 内容在 `Snapshot.Pending` | — | 零 | — | 低（同 task.delta） |
| `ball.state` | 另一条原生管：`Ball.SetState`（`internal/ball/ball_windows.go`，⚠ wave62 在飞——本件未进内文，只记包名） | — | 零 | — | **不在此票做**（归球/33 簇） |
| `cost.tick` | 真没内容（`ListCostDaily` 在，无 tick 推送） | — | 零 | 2＋ | 中 |

## §2 裁定：落地顺序与集线器

**顺序**：堆乙 9 枚先行（其中 `tasks.list`/`task.detail`/`approval.current`/`cost.summary` 四枚最低价可同批）；堆丙 5 枚逐枚 A## 批准后落（`approval.decide` reject 支、`grants.revoke`、`models.list` 拼写、`transcript.get`、`models.delete`）；丁堆 5 枚**单独一张出向名册**（c1 §10#8 的建议今天维持，理由不变：入向白名单与出向事件名混表＝本票要治的病的镜像）。**「先片 B（生产听众）后名册」的旧裁定（A384 裁定③）已经被盘上事实追平**：`Handle` 生产调用者今天＝2 枚，268/248/320 三票已在名册补齐之前落了真实听众——本票从「收件类」变回「可派写腿」，不再等票 33。

**集线器**：**`internal/panel/bridge.go` 是唯一强集线器**——23 枚入向名字全部过 `knownComposerMethod`（`bridge.go:154-159`），任何一枚落地都要在这里加一行常量；次级集线器＝`composer_dispatch.go` 的 `dispatch` switch（`:199-237`）＋装配函数 `newComposerDispatchChain`（`cmd/wisp/panel_inbound.go:264-347`）——动一枚要过「常量→case→interface/handler→装配」四件套，但全部是**追加式**（248/320 先例），不回改既有行即不撞 `composer_test.go:502` 那枚「四枚顺序与命名」的钉。

## §3 本程自证

- **只读**：零产码改动；`internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`go.mod`/`go.sum`、thresholds/golden、allowlist 一字节未动；AC 框零触碰、票面原句未改（只追加 Progress log 一行）。
- **起手在飞件（具名登记、未动）**：`M internal/ball/{ball_windows,orb_windows,renderer_windows,orb_live_windows_test}.go`（wave62 腿）＋未跟踪 `cmd/wisp/panel_resident_323_windows_test.go`、`cmd/wisp/panel_resident_assembly_323_test.go`（323 腿）。本件读数全部走 `git grep HEAD`（已提交态）或磁盘只读，不依赖也不影响在飞件。
- **门禁**：零 Go 改动 ⇒ gofumpt/vet 无对象（如实登记「零 Go 改动、无对象可扫」）。本程核验＝自己的 `git grep`/`sed` 读数（各尺已在 §0/§1 具名）。
- **未做（宁缺毋造）**：未跑任何 `go test`；未复量 248/320 用例的完整名册（只列头 12 枚具名）；`bookPanelSnapshot` 在 `run.go:958` 的行号出自 33-r14 提交信息（未复量，已标注）；丁堆 5 枚的出向「改一个牵一片」程度未做单独的调用链普查。
- **一次转述笔误（就地更正）**：195 号票地基文件名我第一稿写成 `sectionwrite.go`，精读后正确为 `panelwrite.go`（`Manager.WriteSection` `:168`）——本表按正确名落定。
