# 102r2-r1 — 票 102 洞重开修复腿交件读数（写腿自录）

**性质**：写腿（实现者）交件读数，不是裁决表；对抗验收与裁决表归验收方另出。
**日期**：2026-10-31 · **写面**：`internal/panel/**` ＋ 本证据件（派单书面范围，逐枚申报见末节）。

## 0. 开工三核与共树事实

- 开工 HEAD＝`dev@c73bde4b7`（`git merge-base --is-ancestor c73bde4b7 HEAD` rc=0）；索引净；`git status --porcelain -- internal/panel` 空（42-v1/42-r1 已收，panel 无他人脏）。
- 腿进行中共树落了 3 枚他人提交（`fad2f5e4f` 票 385 落地腿、`57e57d8bf`/`32c8498e7` 编排者 A1304），`git diff --name-only c73bde4b7..HEAD` 只含 `tools/mockllm/**`、`docs/**`、`.scratch/**` ⇒ `internal/{panel,risk}` 两对照面零变化，本腿全部改前/改后读数有效。

## 1. 缺陷与红句逐字（改前，`go test ./internal/risk/ -run TestC26RewriteAccount -count=1 -v`，rc=1）

`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 四条红句逐字（`pathresolver_rewrite_account_test.go:200`）：

```
..\..\internal\panel\composer.go:127 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
..\..\internal\panel\display_scrub_216.go:26 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
..\..\internal\panel\display_scrub_216.go:37 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
..\..\internal\panel\display_scrub_216.go:218 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
```

与派单「composer.go:127／display_scrub_216.go 三处」计数 1:1（composer 1 处＋display_scrub 3 处＝4 hit），**无行号漂移**，停止条件①不触发。

**判据形状**（测试自述，本腿只读不改）：文件级静态判据——凡读 `Result.Canonical` 字段的非测试 .go（仓根 WalkDir，豁免生产者 `pathresolver.go` 与 `tools/**`），该文件必须出现 `Actable(` 或 `.Rewritten`；扫到零 hit 则 `the instrument is broken`（非扫空＝绿）。display_scrub 的 :26/:37 是文件头 roster 注释行、:218 是代码点——同一文件按判据算一条失败，故**真实现场是两个消费文件**：composer.go:127（`NewSnapshot` 的空行兜底网）与 display_scrub_216.go:218（`scrubComposerForDisplay` 对 `WorkspaceView.Canonical` 的洗）。

**risk 公开面核**（只读）：`Result{Canonical, Reparse, Resolved, Spelling, Rewritten, Rewrites}` ＋ `Actable() (string, error)`——rewrite 后路径就是 `Canonical`（展开后坐标），账本＝`Rewritten`/`Rewrites`/`Actable()` ⇒ 停止条件③（无 rewrite 后路径可读）不触发。

## 2. 修法：两处真实消费（非装饰符号）＋一枚行为钉

panel 包内既有消费先例是同形的：`git.go:251-257`（`ReadGitForWorkspace` 读 `ws.Rewritten` 后分支叙述）、`instructions_200.go:120-131`（`ProjectInstructionLoadRequestFor` 读 `ws.Rewritten` 后拒绝带坐标的目录腿）。本腿按同一形状把两处新消费点接上账：

1. **composer.go（`NewSnapshot` 兜底网，红点 ：127）**：网先读账。`composer.Workspace.Rewritten==true` 且整行无叙述 ⇒ 补 `workspaceRewrittenUnnarratedReason`（新常量：改写发生过、显示的坐标是改写后的树、不是使用者拼写的那一串）。「没有读数」分支**一字未动**地保留在其后——账本证明「这一维被读过」（只有 C26 展开步能把 Rewritten 置 true），该行因此永不背「没有读数」这句假话（票 216 AC#4 对句纪律的另一半）。生产者侧（`WorkspaceViewFromRoot`/`RequestWorkspaceSwitch`）对每条改写行本就写叙述，此网兜的是绕过生产者到达边界的行。
2. **display_scrub_216.go（`scrubComposerForDisplay`，红点 ：26/:37/:218）**：读账后分支——`c.Workspace.Rewritten==true` ⇒ `c.Workspace.Reason` 一并 `scrubForDisplay`。语义：改写行的叙述是 roster 上唯一**引用 C26 坐标与原拼写**的文本（census 判 `Canonical` 必洗的正是「不是使用者打的内容」）；不读账，被洗的坐标就从叙述这条旁路未经转义地到达人眼，显示边界被这行自己的解释击穿。未改写行的叙述是本仓常量句，走恒等快路径。文件头 roster 段同步登记第三个条件项（217 census 的 roster 保持封闭：字段没增，是**账本标记了某字段内容类别改变的行状态**）。
3. **`internal/panel/rewrite_account_102r2_test.go`（新文件，行为钉）**：静态判据只挡符号消失、不挡条件被改（票 102 AC#3 M1 的自述结论），故两处消费各钉行为半边——①改写无叙述的行出边界必带账户叙述句、且永不被盖「没有读数」；未改写空行仍得「没有读数」句（网未移位）。②改写行叙述里引的坐标含 U+202E ⇒ 出边界为 `\u202e` 可见转义、包字节无裸 U+202E、`Canonical` 本身的洗仍在。

**改前 → 改后 diff 要点**（`git diff --cached`：3 件 128 insertions / 0 deletions）：

- composer.go：+21 行＝常量 `workspaceRewrittenUnnarratedReason`（含注释）＋网内账户分支（`if composer.Workspace.Rewritten && strings.TrimSpace(composer.Workspace.Reason) == "" { composer.Workspace.Reason = workspaceRewrittenUnnarratedReason }`）及其注释；原「没有读数」条件与句子零改动。
- display_scrub_216.go：+18 行＝roster 注释登记＋`scrubComposerForDisplay` 开头的账户分支（`if c.Workspace.Rewritten { c.Workspace.Reason = scrubForDisplay(c.Workspace.Reason) }`）及其注释；`Canonical` 洗与附件洗零改动。
- 既有测试一字未动；A380 断言一字未动（未触发停止条件②——panel 无任何断言钉「改写行 Reason 为空」，`composer_test.go:214` 的改写行夹具不钉 Reason，`pump_test.go:123` 字节钉只钉四个顶层键）。

## 3. A380 红→绿（必给读数 2）

- 改后真仓：`go test ./internal/risk/ -run TestC26RewriteAccount -count=1` ⇒ **rc=0**（`ok github.com/CarlosShao/wisp/internal/risk`）；腿末复跑仍 rc=0（HEAD 已前进到 `32c8498e7`，见 §0）。

## 4. 反向判据三读数＋还原核（必给读数 3）

| # | 树 | 状态 | 命令 | rc | 红句 |
| --- | --- | --- | --- | --- | --- |
| ① | 真仓 `c73bde4b7`（改前） | 无修复 | `-run TestC26RewriteAccount -count=1 -v` | 1 | 4 hit 逐字见 §1 |
| ② | 真仓（修复在位，暂存态经 `git write-tree`=`78ebcefc3615dcc73484ebea5e1ca66a5c79774e` 导出） | 有修复 | 同上 | **0** | — |
| ③ | 仓外快照 `D:\tmp\wisp\20261031\102r2\mut102r2-tree`（②的 tar 展开），**两件还原 pristine**（`git show HEAD:…` 的改前副本覆盖） | 摘掉修复 | `-run TestC26RewriteAccountIsConsumedAtEverySecurityLeg -count=1 -v` | **1** | 4 hit **逐字与①相同**（composer.go:127／display_scrub_216.go:26/:37/:218） |

- 快照正控②先验绿（rc=0），排除「快照自身红」的假信号；再摘修复得③红。
- 突变全程在 `D:\tmp\wisp\20261031\102r2\**`（仓外），真仓 `git status --porcelain -- internal/panel docs/evidence/s1` 只见本腿 3 件（2 M＋1 A，全部 staged）；快照内的 pristine 覆盖不触碰真仓 ⇒ 还原核成立。
- ⚠ 方法注记：A380 的运行时扫是 WalkDir 读**盘上源码文本**，`-overlay` 只换编译输入不换运行时读 ⇒ 反向判据必须用仓外快照真改文件，不能用 overlay（这正是本腿用 tar 快照的原因）。

## 5. 零回归按名作差（必给读数 4）

`internal/panel` 全套 `-count=1`（改前 → 改后，逐名）：

| 既存红（派单预告 4 条） | 改前 | 改后 |
| --- | --- | --- |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes`（approval） | FAIL | FAIL（不变） |
| `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`（StreamLog×2 之一） | FAIL | FAIL（不变） |
| `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`（StreamLog×2 之二） | FAIL | FAIL（不变） |
| `TestComposerContractTypesMatchFrontend`（composer 契约漂移） | FAIL | FAIL（不变） |

**零新增、零消失**：改前 rc=1、改后 rc=1，FAIL 名集 1:1 相同（4 条均为既存红，非本腿责任）。
`internal/risk` 全套 `-count=1`：改前恰 1 红＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` → 改后 **rc=0 全绿**，A380 转绿之外零变化。
本腿新增行为钉两枚（`TestRewrittenWorkspaceRowLeavesTheBoundaryWithItsAccountSentence`、`TestRewrittenNarrationQuotingTheCoordinateIsWashedAtTheBoundary`）改后随 panel 全套跑绿。

## 6. 两尺／导出树（必给读数 5）

- `gofumpt -l`（`$(go env GOPATH)/bin/gofumpt.exe`）对本腿三件：**空**（rc=0）。
- `go vet ./internal/panel/`：**rc=0**。
- `sh scripts/d22scan.sh` 于导出树（`D:\tmp\wisp\20261031\102r2\mut102r2-tree`，修复在位版）：**rc=0，零 finding**。git 不可问的响亮模式（`gitignore rules NOT APPLIED`，全扫 319 生产 Go 件），8 scope 台账逐行：`bans #1-5 internal/=274`、`#1-5 cmd/=45`、`#6 frontend/=116`、`#7 internal/tools/=25`、`#8 design/=14`、`#8 frontend/=116`、`#8 internal/=707`、`#8 cmd/=162`；对照票 102 时代参照（#6=37、#8=16/37/335/25）只增不减（后续腿进树件），**本腿三件在台账与 finding 里零点名**。d22scan 自带 self-test 套件同跑 35 PASS / 0 FAIL（含正控：seeded violation rc=1、armed ban #6 rc=1、空 scope rc=2——仪器有效，非扫空＝绿）。

## 7. 写面申报（逐枚，必给读数 6）

| # | 件 | 动作 |
| --- | --- | --- |
| 1 | `internal/panel/composer.go` | 改（+21：账户分支＋常量＋注释） |
| 2 | `internal/panel/display_scrub_216.go` | 改（+18：账户分支＋roster 登记＋注释） |
| 3 | `internal/panel/rewrite_account_102r2_test.go` | 新增（行为钉 2 枚） |
| 4 | `docs/evidence/s1/102r2-r1-verdict.md` | 新增（本件） |

⛔ 未动一字：`internal/risk/**`（347 刚落）、`cmd/wisp/**`、`internal/{speech,session,config,memory}/**`、thresholds/golden/配置键/testdata/frontend/design/scripts/.github/票面/台账/specs/PLAN、任何既有测试（含 A380 断言）。⛔ 未 push；commit 单枚、显式 pathspec、主题逐字 `票 102 洞重开修复腿 102r2-r1`。

## 8. 相邻观察（不属本腿，留编排者记账，未越权改）

- `workspace.go:186-195` `currentAfter` 的范围-账户不一致句内嵌 previous/now 的 `Canonical` 与 `accountSummary`（含 Spelling），该行 `Rewritten` 可为 false ⇒ 本腿「按账洗叙述」的闸不覆盖它。要覆盖＝把 Reason 洗成无条件或另立判据，属 216 roster 的下一次裁定，本腿只记账不动手。
- `WorkspaceView` 不携带 `Rewrites []string`（只有 bool），兜底叙述因此不含改写种类清单；生产者叙述（`workspace.go:124-126`）有清单，兜底网是它的下位替代，这个信息差是有意的最小化。
