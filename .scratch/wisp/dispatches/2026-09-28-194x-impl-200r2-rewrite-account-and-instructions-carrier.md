# 派单 2026-09-28 19:4x — 写腿 `200-r2`：票 200 交件后我复算出的三处账（改写账户被你的新消费者打红／面板那枚键今天没人填／一枚测试格式）

> 上一腿 `200-r1` 交件 `84feec44`（9 枚 pathspec：`internal/projctx/projctx.go` +468、`projctx_test.go` +442、
> `internal/agent/instructions_test.go` +305、`internal/agent/prompt.go` +58/−1、`internal/config/schema.go` +4、
> `internal/panel/composer.go` +7、`internal/panel/instructions_200.go` +72、`cmd/wisp/run.go` +35/−1、证据件 +87）。
> **我逐枚复跑过**：`projctx` ok 0.187s、`config` ok 1.221s、`internal/tools` 未碰（该 commit 里 `internal/tools` 命中 0 枚）——这些都对。
> **但有几处你写成"归因未定案"的，我已经定案了，而且是判红。**

## 0. 起手必做

1. `date`＋`git log --oneline -8`＋`git status --short` 自取锚点；本单行号是我 19:4x 读的，会漂。
2. **撞钉预检**：动手前跑一遍并把**今天的绿用例名**抄进证据件——
   `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./internal/risk/ ./internal/panel/ ./internal/agent/ ./internal/projctx/ ./cmd/wisp/`
   ⚠ `cmd/wisp` 不带那枚 PATH 会 `exit status 0xc0000135`（一枚不跑＝仪器瞎）。
   起手在册的红**名册**（我 19:4x 现量，逐名）：`internal/risk`＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`；
   `internal/panel`＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、
   `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestC21DesignTokensFourWayAgree`（后三枚是在册常红，**不修不 Skip**）。
3. 读票 `.scratch/wisp/issues/200-*.md`＋上一腿证据件 `docs/evidence/s1/200-project-instructions-r1.md`。

## 1. 硬规矩

- **只 commit、绝不 push**；显式 pathspec；禁 `git add -A`/`.`/`-a`；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`restore`/`clean`/`merge`/`worktree`。
- 新文件逐枚 `git add -- <路径>`。**临时件只建不删**；还原被自己写坏的文件只准 `git cat-file blob HEAD:<path> > <path>`。
- 别人的脏文件不碰不提交：`.gitignore`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`design/**` 那 16 枚删除。
- **零写面**：`frontend/**`、`design/**`（不写不读不转述）、`internal/panel/tokens_fourway_test.go`、
  **`internal/tools/**`（另一腿 197-r1b 正在里面改池常量，你碰它就是撞车）**、`PLAN.md`/`docs/specs/**`/`thresholds.go`/golden/`allowlist.txt`。
- 中文长段走 Edit/Write 工具，**别用 shell heredoc**（反引号会被真执行）。

## 2. 这一腿做三格

### (a) 你把票 102 那枚**在册判据打红了**——这条不是"归因未定案"，红句子里就写着你的行号

- 我 19:3x 现量：`go test -count=1 ./internal/risk/` ⇒ **唯一红**＝
  `pathresolver_rewrite_account_test.go:200: ..\..\cmd\wisp\run.go:686 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer`
- 我做的**因果对照**（不是猜）：`git show 84feec44~1:cmd/wisp/run.go | grep -c '\.Canonical'` ＝ **0**；改后 ＝ **1**（就是你加的那行 `instrWorkspace = ws.Canonical`）；
  且 `grep -c 'Actable(\|\.Rewritten' cmd/wisp/run.go` ＝ **0**。⇒ **这枚红是 200-r1 带进来的，不是环境噪声。**
- 尺的语义（自己读 `internal/risk/pathresolver_rewrite_account_test.go:187-203`）：**按文件**判——任何非产码文件读了 `.Canonical` 就必须同文件里出现 `Actable(` 或 `.Rewritten`。
  它要的不是格式，是"每一次用权威路径之前先看一眼改写账户"。
- **修法（要语义不只是过尺）**：`panel.WorkspaceView` 本身就带账户（见 `internal/panel/workspace.go:104-111`，字段 `Rewritten`/`Reason`），
  所以在把 `ws.Canonical` 交给 `projctx` 之前**读 `ws.Rewritten`**：被改写过的视图 ⇒ **一枚说明文件都不读**，并按你们自己那套"响亮而非沉默"的规矩打一行为什么（`rt.auditf` 现成）。
  ⚠ **不许只塞一句 `.Rewritten` 进去骗过尺**（比如算了又不branch）——判据要落在"喂一枚 Rewritten 的视图 ⇒ 加载器不去读那个目录"这种**行为**上（正控）。
- 交件判据：`go test -count=1 ./internal/risk/` 全绿，且你能报出"哪一枚 commit 算未修码"（改前红读数已在上面）。

### (b) 面板那枚 `instructions` 键：**今天有人在填吗**——没人填就把它拿掉，别长成我们骂过别人的那种样子

- 现量（我 19:4x）：`internal/panel/composer.go:68` `Instructions []ProjectInstructionFile \`json:"instructions,omitempty"\``；
  `internal/panel/instructions_200.go:41/65` 有 `ProjectInstructionsFromBundle` 与 `WithInstructions`，
  但 `grep` 在 `internal/panel/pump.go` **零命中** ⇒ 你说"泵填该键的两跳按'宁可缺一跳并具名'留手"，我复算**证实**了：载体在、装配侧没人喂。
- ⚠ **为什么这比缺一跳更坏**：我在账 `A408` 里刚把 OpenChamber 的 6 处"文案在、控件不存在"点成反面教材；**"JSON 键在、值恒缺"是同一枚病，只是换了层**。
  而且 `omitempty` 让缺值时键根本不出现在 JSON 里 ⇒ 界面按"有没有这个键"分支的代码会永远走"没有"那一支，**没人会发现**。
- **二选一，选一支并写理由与读数**：
  - **甲（我推荐，也是"做完整功能"那一支）**：把两跳接上——每轮从 loop 拿 `Bundle`/清单，经 `WithInstructions` 进快照，
    并给一条**端到端**判据：跑一次真装配（`cmd/wisp` 那层现成的 harness），断言快照 JSON 里 `instructions` **带着真文件名**出现；
    再加一条反向：**没加载任何文件时该键为空数组/缺省，且界面能区分"没开"与"开了但目录里没有"**（你们已经打过 `已按你的配置跳过` 那行字，把它传上来）。
    ⚠ **我 19:5x 现量这枚现形**：`instructions_200.go:43-46` 对 `nil`/空 bundle 一律 `return nil`，而 `composer.go:68` 的标签带 `omitempty`
    ⇒ **`Enabled=false` 与"开着但目录里一份说明都没有"今天 marshal 出的是同一枚"键不存在"**，界面**无从区分**、也**无法说真话**（它只能说"没东西"，不能说出**为什么**没东西）。
    ⇒ 走甲就必须把这第三种状态（`off`／`on-but-none-found`／`on-and-loaded`）**带上来**，形状你自己定，但**判据要能正反各跑一发**。
  - **乙**：这一腿**把键从 marshal 里摘掉**（连 `omitempty` 一起），留到泵那天再长回来——**代价是票 200 AC#7 的"面板可见"那一半当场作废**，必须在票里写明"哪一格退回、为什么"。
- ⚠ **同格附带一枚你我都改不了的**：`TestApprovalCardViewJSONKeysMatchFrontendTypes` 现在红，
  红因是 Go 侧发了 `instructions` 而 `frontend/src/lib/panel.ts` 的 `PanelSnapshot` 没声明它。**`frontend/**` 对我和对你都是零写面** ⇒
  **不要试图绕它**（别改尺、别加豁免、别让 Go 少发一个键来"过尺"——那是把两侧同批移动这条纪律拆掉）。选甲就把它当**在册常红**具名报，我把这一跳写进票面交给 owner 带给界面那支；选乙则它自然消失，也要具名说明。

### (c) 两枚小账，一次扫干净

- **格式**：`"$GOPATH/bin/gofumpt.exe" -l` 现在点着 **`internal/agent/instructions_test.go`**（你那一腿新带的；`gofumpt` 不在 PATH 但在 `$GOPATH/bin`，`gofmt -l` 空 ≠ gofumpt 空）。
  ⚠ **只处理你自己那一枚**；`internal/tools/subagent_197.go`／`subagent_197_test.go` 那两枚未格式化件归 197-r1b，**你别动**。
  交件判据：跟踪集里 `internal/`、`cmd/` 下不再有 gofumpt 未净件（`.scratch/wisp/probes/**` 那 4 枚刻意做坏的变异样本允许还在，逐名列出即可）。
- **一枚不复现的红**：你报 `TestResolvePerCallBudget`"起手即红、中途转绿、终态又红"。我在 `84feec44` 之后**单跑三遍全绿**（`go test -count=1 -run TestResolvePerCallBudget ./internal/risk/` ok 1.2s/1.4s/1.4s）。
  ⇒ 你要不然**给出能复现的完整命令＋读数**，要不然**具名写"我不复现，交回编排者"**。**别把它留在票面上当已知红**（我这条规矩：写"今天应该红"必须自带可复跑的尺）。

## 3. 交件

1. 代码一次 commit（message 带 `200-r2` 与票号 200）。
2. 证据件 `docs/evidence/s1/200-project-instructions-r2.md`：① 起手锚点与名册；② (a)(b)(c) 各自改前→改后读数；
   ③ **两发变异/正控**（(a) 喂 Rewritten 视图必不读；(b) 选甲就断"没接泵时快照里这枚键不出现在真装配的 JSON 里"）；
   ④ 门禁终态逐包读数＋在册红名册逐名比对（**少了哪一枚、多了哪一枚**要具名）；⑤ **没测什么**；⑥ 对我派单里任何一处不服——我给的红因、行号、常量有错就当面写。
3. 只报盘上可核的事实：commit 号、`wc -c`、判据逐名读数。我不认"应该已经"。
