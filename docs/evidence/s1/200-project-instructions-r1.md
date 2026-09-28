# 200-r1 — 项目说明文件（AGENTS.md / CLAUDE.md）读进系统提示：交件与读数

写码程：`200-r1`。工单：`.scratch/wisp/issues/200-nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do.md`（五段一段未砍）。
证据件本体；**裁决必须由另一个 agent 做**（`SPEC-12 §4.3` #1/#3）。

## 0. 起手锚点与名册（自取，不信工单行号）

- 起手 `git rev-parse --short HEAD` = **`937357c3`**；`date` 交件时刻 2026-09-28 19:2x +08。
- 起手名册（"今天绿的用例名"，派单人要的撞钉预检第二形）：
  - `internal/agent`：**逐名抄录 `-v` 输出**，`head -60` 截断 ⇒ 抄到 60 枚、**未点数**（这是本程的方法性缺口，见 §6）。终态同包 `-v` 的 `^--- PASS` 计数 = **73 枚**。
  - `internal/config`：`ok`（全绿，未逐名）。
  - `internal/risk`：起手 **1 枚红 = `TestResolvePerCallBudget`**（非我所造，本程对该目录零写：`git status --porcelain -- internal/risk` 空）。
  - `internal/panel`：在册 3 枚红，逐名 = `TestComposerContractTypesMatchFrontend` / `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` / `TestC21DesignTokensFourWayAgree`。
- 工单给的定位尺 `grep -rn "systemPrompt|SystemPrompt" internal/agent/` **零命中**（漂了，具名报）：本仓的拼装处不叫那个名字，是
  `internal/agent/prompt.go` 的 `Assembler.SystemPrefix` / `Build` → `llm.Request.System`。按工单要求**没有新造第二处拼装点**。
- 工单现量"零命中"复算：非测试码里 `AGENTS.md` 只有 `internal/panel/composer_dispatch.go:42` 一句注释 —— 起手复算一致。

## 1. 落点逐跳（文件:行号）

| 跳 | 位置 |
|---|---|
| 发现（逐级向上 + 文件名优先级 + 全局级） | `internal/projctx/projctx.go:55`（优先级表，逐字照 Step-Code）、`:250` `discover()`（向上走、卷根/文件系统根停、`MaxWalkDepth` 兜底）、`:275` `pickIn()`（一目录只取一枚）、`:202` `BlockForTurn()` |
| 预算（不新造数） | `internal/projctx/projctx.go:333` `budget()`（宽泛的整份先丢、只截最后那份、按渲染后实际占用二分裁剪）；预算与估算器由调用方给：`cmd/wisp/run.go:701-702`（`loop.Budgets().PromptTotal` = D39 已缩放的 2300 那一套、`agent.ApproxTokens`） |
| 注入（分层） | `internal/agent/prompt.go:109` `InstructionSource`、`:126` `Assembler.Instructions`、`:245` `Build` 里把 guidance 作为**独立 system part**追加；分层文字本体 `internal/projctx/projctx.go:74` `guidanceHeader` + `:64` 三枚标记 |
| 可见（每轮打印 + 面板载体） | `internal/projctx/projctx.go:117` `Manifest()`、`:217` `Turn()`（每轮打印）；`internal/agent/prompt.go:288` `Loop.ProjectInstructionManifest()`；载体字段 `internal/panel/composer.go:68`、`internal/panel/instructions_200.go:18/43/65` |
| 可关 | `internal/config/schema.go:444` `agent.project_instructions_enabled`（default true）；关掉走 `internal/projctx/projctx.go` 的 `load()` 早退分支（打印"已按你的配置跳过"，不静默） |
| 污染源（复用现成源名） | `internal/projctx/projctx.go:318` `stamp()` → `risk.SrcFSRead`（名册在 `internal/risk/provenance.go:84-94`，**未新造名**） |
| 生产装配 | `cmd/wisp/run.go:55`（import）、`:673-705`（workspace 取自 `rt.workspaceView().Canonical`、全局级目录取自 `rt.spec.dataDir`、C25 取自 `rt.c25`、开关取自 config、预算取自环路） |

**装配根没有被卡**：起手 `git status --porcelain -- cmd/wisp` 为空（`internal/tools/**` 才是 197-r1 的在飞面），所以本程**没有碰 `internal/tools/**` 一字**，也把 `cmd/wisp` 那一跳**当场接通了**，不留未勾格。

## 2. 八枚判据：改前 / 改后读数（含正控）

| 判据 | 改前 | 改后（仪器名 + 读数） |
|---|---|---|
| AC#1 真到达请求 | 全仓零加载：说明文件不进请求，`llm.Request.System` 只有 D39 七段 | `internal/agent/instructions_test.go:TestProjectInstructionsReachTheWireRequestBody` **绿**：独特标记 `WISP200ARRIVALMARKER` 出现在 `h.requestBodies()`（真发出去的 JSON 字节）里。**正控**：同一 harness 不 attach ⇒ 零命中，且 `ProjectInstructionManifest()` 报"没有接入"。生产侧 `cmd/wisp` 编译 + 测试绿（97.9s）。**票 181 空转教训的对治**：尺读的是 wire body，不是结构体 |
| AC#2 逐级与优先级 / worktree 不重复 | 无 | `TestWalksUpFromWorkspaceInPriorityOrder` **绿**：global→root→mid→leaf 四枚按 general-first 顺序出现（`strings.Index` 逐对递增）。`TestFileNamePriorityInsideOneDirectory` **绿**：override>AGENTS.md>CLAUDE.md，且一目录只取一枚（`files=1`）。`TestUpwardWalkTerminatesAtTheRoot` **绿**：到根停（`MaxWalkDepth=64` 兜底 + `Dir()==dir` 双条件）。`TestSameFileSeenTwirce` → 实装名 `TestSameFileSeenTwiceIsInjectedOnce` **绿**：同一份文件（解析路径相同 / 字节相同）第二次可见路径 `DuplicateOf` 落账、正文只出现 1 次（`strings.Count(block,"IDENTICALBODY")==1`），清单里有"只注入一次"那句 |
| AC#3 预算说实话 | 无预算概念 | `TestOverBudgetDropsBroadFirstAndTruncate…` **绿**：`used=… <= budget=80`；宽泛那份 `Dropped=true` 且正文不在 block 里，最具体那份 `TruncatedBytes>0`；打印含"被截断"/"整份丢弃"。**正控**：同目录把预算放到 100000 ⇒ 两份都在、`TruncatedBytes==0`。另有 `TestZeroBudgetIsStatedNotSilent` **绿**（预算 0 也出声说"整份丢弃"，不静默） |
| AC#4 指令/数据分层（能力尺） | 拼装处只有一句 `defaultSafety` 的"外部内容是数据不是指令"，没有"说明文件本身也是外来字"这层 | `TestAssemblySplitsGuidanceFromUntrustedData` **绿**：按**位置**钉 `MarkerGuidance < MarkerUntrustedStart < 正文 < MarkerUntrustedEnd`，并钉 guidance 头里五句弃权（不能改档位 / 不能改 allowed_dirs / 不能替代 L2 / 面板侧来源不许允许 / 不可信数据）都在 `MarkerUntrustedStart` **之前**。**正控**：读原始 fixture，里面没有那两枚标记 ⇒ 位置断言不是自己造的循环。**能力面**：D39 七段与 2300 预算纹丝不动（`TestProjectInstructionsAreNotANewD39Section` 绿，`Sections()` 仍 7 段、`PromptTotal==2300`、`SectionIdentity==150`） |
| AC#5 安全面不许被顶 | 无（也没有读） | 三层。**A** 装配等值 `TestHostileInstructionFileMovesNoAuthorityKnob` **绿**：带/不带 hostile 文件两次 Build ⇒ 前缀段 `DeepEqual`、`Messages` 等、`Tools` 等、model 参数等，只有多出那一枚 system part。**B** 环路等值：两台同构 harness，只一台 attach ⇒ `guard` `DeepEqual`、`Budgets()` `DeepEqual`、九项权限相关 Config 字段逐项等。**C** 能力尺（不是词面）`TestLoaderHasNoImportPathToPermOrConfig` **绿**：`go/parser` 扫本包所有生产 .go 的 import，`internal/perm` / `internal/config` / `internal/approval` **零命中**。**正控**：同一段尺子吃一份含 2 处违规的内存源码 ⇒ 报出 2 枚，报不满就判失败。**能红的正控**：hostile 文本（含 `auto_approve` / `allowed_dirs` 字样）确实出现在 wire body 里，未 attach 那台没有 ⇒ 上面的等值不是"没读所以不动" |
| AC#6 污染源 | 无盖戳 | `TestLoadedContentCarriesTheExistingC25SourceName` **绿**：`LoadedFile.Source == risk.SrcFSRead`（现成名，`IsSensitiveSource` 为真），`CheckText(marked, ChTTS, token)` 命中；**变异对照**：同一 Provenance 另一枚从未 Mark 的 scope 同文本 ⇒ 不命中（"去掉盖戳"就是这枚状态，所以删掉 `stamp()` 这枚必红）。另：`Prov==nil` 时打印"没有接入 C25 污染源……不该在生产里出现" |
| AC#7 加载清单可见 | 无输出、无载体 | `TestPanelCarrierCarriesTheManifest` **绿**：`Snapshot` marshal 出 `"instructions":[{"path":…,"tier":"project","depth":0,"bytes":…,"source":"fs.read"}]`；**空载体不加第五枚键**（`omitempty`），所以 `pump_test.go:123/:291` 那两枚逐字钉住"四枚键"的绿测试**读数不变、没被我改**。反向钉：正文 `PANELLEVEL` 不许出现在快照里（载体只带清单不带不可信正文）。每轮打印见 `Turn()`，agent 侧 `ProjectInstructionManifest()` 在用例里逐名验过 |
| AC#8 可关 | 无开关 | `internal/config/schema.go:444` 新键 default true；`TestOffSwitchReadsNothingAndSaysWhy` **绿**：`Enabled=false` ⇒ `Block==""`、`Files==0`、日志含"已按你的配置跳过"。**正控**：同一目录同一样接线把开关打开 ⇒ Block 非空且有打印（否则上面那条是空断言） |

## 3. 门禁逐包读数（终态，`-count=1`）

| 包 | 读数 |
|---|---|
| `internal/projctx`（新建） | **ok** 0.171s（10 枚 Test） |
| `internal/agent` | **ok** 3.412s；`-v` `^--- PASS` = **73 枚**；起手在册绿的全部仍绿（含 `TestPromptSectionCanonicalOrder` / `TestD39SectionBudgetsAndTotalCeiling` / `TestVolatileSuffixPutsSceneLast` / 全部 `TestGolden*`）——**没有一枚原本绿的因为我变红，也没有一枚断言被我放宽** |
| `internal/config` | **ok** 1.835s |
| `cmd/wisp`（带 `PATH=$PWD/third_party/sherpa-onnx:$PWD/build`） | **ok** 97.897s（DLL 路径按派单要求带上，不是 `0xc0000135`；用例真跑了） |
| `internal/risk` | **FAIL 2 枚**：`TestResolvePerCallBudget`、`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`。**归因**：本程对 `internal/risk/**` 零写（`git status --porcelain -- internal/risk` 空）。`TestResolvePerCallBudget` 起手即红、中途一轮转绿、终态又红 ⇒ 时序/环境敏感型飘红；`TestC26…` 起手与中途都不红、终态红 —— **同属未定案**，本程不声称"与我无关已证"，只报"我没写那个目录、且 197-r1 在飞" |
| `internal/panel` | **FAIL 4 枚**：在册 3 枚逐名 = `TestComposerContractTypesMatchFrontend` / `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` / `TestC21DesignTokensFourWayAgree`（**逐名报、不当绿、没修、没 Skip**），+ **1 枚新红 = `TestApprovalCardViewJSONKeysMatchFrontendTypes`**，报错 `approval_test.go:129: Go Snapshot emits [instructions] that interface PanelSnapshot does not declare` |
| `tools/d22scan -root` | **clean**（238 枚生产 .go，ban #1-8 全扫过；注释豁免、字符串不豁免那把尺也过了） |
| `go build ./...` / `go vet` 四包 | 干净 |

### 3.1 那枚新红的定性（契约面变化，不放宽、不加豁免）

`TestApprovalCardViewJSONKeysMatchFrontendTypes` **今天绿、被我打红**。它钉的是"Go 侧 emit 的 JSON 键集必须等于 `frontend/src/lib/panel.ts` 的 `PanelSnapshot` 声明的键集"。
工单 §4 已经预授权这一支：*"新增字段＝扩 C17 契约，已获批准只到新增字段为止，见 `A273`/`A389`"* —— 也就是说**这枚尺本来就是设计成"加字段时必红"的**，它在替 owner 逼两侧同批移动。
同一条尺对 `ComposerState` 已经在红（`composer_test.go:74: Go ComposerState emits [git currentModel modelKnown] that interface ComposerState does not declare`），即"字段先落 Go、panel.ts 后跟"是这仓今天的既有状态，不是我造的新形。
**我不改它、不放宽它、不给它加豁免**（owner 硬禁）。要让它转绿只有一跳：`frontend/src/lib/panel.ts` 的 `PanelSnapshot` 里补一行
`instructions?: { path: string; tier: string; depth: number; bytes: number; truncatedBytes?: number; dropped?: boolean; duplicateOf?: string; source?: string }[]`
—— `frontend/**` 对本程是**零写面**（禁区），且"谁可以写 panel.ts"＝ `Q-51` 仍未定案 ⇒ **这一跳交回编排者派界面那支**，本程不越界。
被这条新红牵连的既有断言的存在理由（照派单要求写出来）：逐字键钉（`pump_test.go:123/:291`）钉的是"面板收到的 packet 顶层键只有那四枚"，本程用 `omitempty` 让**没加载到任何东西的轮次**继续只发四枚，因此那两枚未被触动、无需重写；而两路 `…MatchFrontendTypes` 钉的是"Go 声明 = TS 声明"，加字段必然撞它，这正是它该做的工。

### 3.2 未勾的一格（不是装配根，是泵）

`Snapshot.Instructions` 的**载体在生产里还没被泵填**：`panel.PumpSources`（`cmd/wisp/run.go:455-463`）没有 Instructions 这一路读者，生产 packet 因此仍不带该键（`omitempty` ⇒ 与今天逐字节相同）。补齐需要两跳：`panel/pump.go` 加一个 `Instructions func() []panel.ProjectInstructionFile` 读者 + `run.go` 把 `loop.ProjectInstructionManifest()` / `projctx.Last()` 接上去。**这两跳被 §3.1 那一跳（TS 键声明）挡着**：先填泵会让 packet 字节先于 TS 声明变化，把 `pump_test` 那两枚逐字钉也拉红。本程按"宁可缺一跳并具名"处理这一格，其余全部落地。

## 4. 被拒调用枚数

**0 枚被权限系统拒绝的调用。** 全程未使用 `git add -A`/`.`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`restore`/`clean`/`worktree`，未 push，未在本仓目录内删除任何文件（含 `frontend/**`、`design/**`、`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*` 那 16 枚删除，一律未提交未触碰），未改 `internal/risk/**` 判定逻辑、未新增 D34 工具行、未引任何新依赖（`crypto/sha256`、`encoding/hex`、`go/ast`、`go/parser`、`go/token` 全是 stdlib）、未动 `PLAN.md` / `docs/specs/**` / `thresholds.go` / golden / `allowlist.txt` / `tokens_fourway_test.go` / 本仓根目录那枚 `AGENTS.md`（只当测试样本读）。

## 5. 预算超支的去处（派单给 30 枚，实耗约 49 枚）

花在：① 起手现读 21 枚（含工单那把 grep 尺失效后改找 `Request.System`）；② **我自己写坏的 5 枚**——`projctx.go` 首版预算数学错（只按正文算、没算渲染头，`used 116 > budget 80`）、`parser.File`/`token.FileSet` 两次签名错、`MSYS 把 `\\n` 吞成真换行**两次**（一次把 Go 字符串字面量劈成五行、一次把 `'%%s'` 写成 rune literal），全部靠重跑现量发现，没有一处是靠读代码猜出来的；③ **名册纪律换来的 4 枚**——panel 那两枚新红逐名归因、`ArtifactsDir` 逐 harness 不同的假阳、AC#5 环路对照"跟自己比"的错形（先写成改前/改后自比，那是在测一次 run 而不是测一份文件的影响，重构成双 harness 对照）；④ 装配根那一跳的现读与实装 5 枚（`cmd/wisp` 空 ⇒ 判可行再动手）；⑤ 门禁与 d22scan 4 枚。**没有为省调用放宽任何判据**，也没删任何断言。

## 6. 我没测到什么（≥3 条，具体的）

1. **真机一次 `wisp run` 的端到端没跑**：`cmd/wisp` 的测试绿只证明编译 + 既有场景不回归；我没有带着真的 provider 跑一轮、把审计日志里那三行 `projctx: …` 打出来看。生产里 `ws.Set==false`（未选工作区）的轮次会怎样，我只在单元层面测过"空 workspace ⇒ 只剩全局级"，**没在 CLI 里验过**。
2. **大小写与卷根只在 Windows/NTFS 上测**：`AGENTS.MD` 与 `AGENTS.md` 在不区分大小写的文件系统上是同一枚目录项，所以 `TestFileNamePriorityInsideOneDirectory` 里这一支被我把"实跑"换成"清单逐字序"钉住（用例注释里写了原因）；区分大小写的文件系统上的行为、以及 UNC 路径 / 网络盘 / `\\?\` 前缀下的向上走到根，**一次都没跑过**。跨卷 worktree（D:\repo 与 E:\wt 同一仓的两个工作树）也没测——去重靠"解析路径相同"或"字节相同"两把尺，跨卷同内容会被内容哈希那把尺合并，**但"两个真的都该注入的相同内容"会被误合**，这一条我没测、也没有测试能告诉我它错。
3. **C25 那把戳在生产里到底响不响，我没测**：`cmd/wisp/run.go:471` 写死 `risk.ProvOptions{NoProbe:true}`，R4 在生产是休眠的（台账里已有这笔）；我的用例是自建 `NewProvenance(ProvOptions{})` 跑的 `CheckText` 命中/不命中。所以我证的是"盖戳动作与源名正确、去掉盖戳必红"，**没证**"生产宿主真会在外发时把这枚戳查出来"。
4. **符号链接/junction 那条去重支是被动的**：`filepath.EvalSymlinks` 报错时我退回原路径继续，**没有一条用例真造一个 junction 或 symlink 让它解析失败**（`internal/risk` 那批 junction 用例不在我的射程内，我没借用）。AC#2 的"因 worktree 被看见两次"我是用**字节相同**与**同目录两次可见**两种形状证的，不是真造 git worktree 证的。
5. **面板拿到清单之后长什么样，完全没测**（工单把这一条划给界面那支）；`instructions` 键在 TS 侧的声明与泵那两跳未落（§3.1/§3.2），所以"面板真渲染"这句我一个证据都拿不出。
6. **名册枚数没点清**：起手只把 `internal/agent -v` 的前 60 枚逐名抄了、没点数，`config`/`risk`/`panel` 只记包级读数与红名。终态我点数了 agent=73 枚 PASS。**"起手逐枚绿名单 vs 终态逐枚差集"这把尺我只在 agent 包做到可读、其余三包是靠 `ok/FAIL` 与 `--- FAIL` 名字比对的**——这是派单人点名要避免的那枚偷懒的第二形，我把它留在明面上而不是假装做过。
