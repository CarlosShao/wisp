# 200-verdict — 票 200「读 AI 说明书」验收腿：八格判据独立复核（非实现者）

- 本程：`200-v1`（wave34 W200 验收腿，只读复核＋一件证据件，零产码改动）。
- 工单：`.scratch/wisp/issues/200-nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do.md`。
- 前任实现腿证据：`200-project-instructions-r1.md`（commit `84feec44`）、`200-project-instructions-r2.md`（commit `a90677e9`）。
- **本件是验收读数，不是裁决表**：翻框与 `-done` 归编排者（`SPEC-12 §4.3` #1/#3，裁决者≠实现者）。
- 读数原始件：`.scratch/wisp/probes/200/v3/logs/01..11`（本程亲跑归档）。

## 0. 起手锚点与在飞自检

- `git rev-parse HEAD` = **`59ac99f922526b521d70e74132c4cf0109668006`**；时刻 2026-11-06 16:21 +08。
  （简报给的锚 `6ae44b1fb` 是 HEAD 的上一枚 commit，`59ac99f92`＝A1570 收 wave33，与本票无关。）
- 票 200 两枚产码 commit 现量在史：`84feec44`（2026-09-28 19:30）、`a90677e9`（2026-09-28 20:12）。
- `git status --porcelain -- docs/evidence/s1 internal/projctx internal/agent internal/panel internal/risk internal/config cmd/wisp`：
  在飞未跟踪件只有 **`cmd/wisp/panel_resident_323_windows_test.go`／`panel_resident_assembly_323_test.go`**（票 323 的腿，逐名可归票）
  ＋ `docs/evidence/s1/` 三枚别票未提件（195/351/378）。**本票写面（`docs/evidence/s1/200-*.md`）零占用，不停手**；
  本程对 `cmd/wisp` 只跑测试零写，未与 323 撞面。
- 本程八条判据全部**亲跑**（三把尺的第三把：自己复跑），未采纳任何腿的转述读数。

## 1. 八格判据：逐格读数（本程亲跑，HEAD `59ac99f9`）

| 格 | 判据（票面原文） | 本程仪器与读数 | 锚（现读行号） |
|---|---|---|---|
| AC#1 | 加载的字**真出现在发给模型的请求里**，正控＝独特标记逐字命中 | `TestProjectInstructionsReachTheWireRequestBody`（internal/agent）**PASS**：断言的是 `h.requestBodies()` 的 wire JSON 字节含 `WISP200ARRIVALMARKER`＋两枚分层标记；反控＝同一 harness 不 attach ⇒ 零命中且 manifest 报「没有接入」。**票 181 空转教训的对治在位**（尺读 wire 不读结构体） | 仪器 `internal/agent/instructions_test.go:69`；接缝 `internal/agent/prompt.go:280`（`AttachProjectInstructions`→`:126` `Assembler.Instructions`） |
| AC#2 | 逐级向上＋文件名优先级；同一物理文件两次可见只注入一次 | 四枚全 **PASS**（`internal/projctx` 全包 **10 PASS／0 FAIL，ok 0.257s**）：`TestWalksUpFromWorkspaceInPriorityOrder:66`（global→外→内，general-first 逐对递增）·`TestUpwardWalkTerminatesAtTheRoot:104`（根停＋`MaxWalkDepth=64` 兜底）·`TestFileNamePriorityInsideOneDirectory:127`（override>AGENTS.md>… 且一目录只取一枚）·`TestSameFileSeenTwiceIsInjectedOnce:149`（解析路径＋字节哈希双尺去重，`DuplicateOf` 落账、正文 `strings.Count==1`） | 实现 `internal/projctx/projctx.go:55`（优先级表）/`:276`（discover）/`:301`（pickIn）/`:319-333`（去重） |
| AC#3 | 超上限要**说实话**：丢了哪份、截了多少 | 两枚全 **PASS**：`TestOverBudgetDropsBroadFirstAndTruncatesOnlyTheLast:187`（宽泛整份 `Dropped=true` 且正文不入 block、最具体那份 `TruncatedBytes>0`、打印含"被截断/整份丢弃"；正控＝预算放宽两份全保）＋`TestZeroBudgetIsStatedNotSilent:428`（预算 0 也出声）。预算数不新造＝调用方给（`run.go:1719` `loop.Budgets().PromptTotal`，D39 已缩放那套） | 实现 `projctx.go:359`（budget，按渲染后实际占用二分裁）/`:130`（Manifest 逐份报状态） |
| AC#4 | 指令＝指引、文件内容＝不可信数据，**分层在拼装处存在**（能力尺） | `TestAssemblySplitsGuidanceFromUntrustedData` **PASS**（log `07-agent-ac4-v.txt`）：按**位置**钉 `MarkerGuidance < MarkerUntrustedStart < 正文 < MarkerUntrustedEnd`，guidance 头五句弃权（不能改档位／不能改 allowed_dirs／不能替代 L2／面板侧来源不许允许／不可信数据）全在 untrusted 起点**之前**；反控＝原始 fixture 无标记（排除空转）＋空 loader 不加 system part | 实现 `projctx.go:63-77`（三枚标记＋`guidanceHeader`）；拼装 `internal/agent/prompt.go:270` 支 |
| AC#5 | hostile 说明文件**顶不动安全面**，且有能红的正控 | `TestHostileInstructionFileMovesNoAuthorityKnob` **PASS**：hostile 文本（`auto_approve`/`allowed_dirs`）真进 wire（正控 1），带/不带文件两次 Build 的前缀段/Messages/Tools/model 参数 `DeepEqual` 纹丝不动，两台 harness 的 `guard`/九项 Config 权限旋钮/`Budgets()` 逐项等（正控 2＝未 attach 那台不含 hostile 字节）；能力尺 `TestLoaderHasNoImportPathToPermOrConfig:239` **PASS**：`go/parser` 扫包内生产 .go，`internal/perm`/`internal/config`/`internal/approval` **零命中**，正控吃 2 处违规源码报满 2（仪器不瞎） | 仪器 `instructions_test.go:187`；实现层无门＝`projctx.go` 包头安全声明（`:18-31`） |
| AC#6 | 加载内容带**现成** C25 源名盖戳；去掉盖戳能判红 | `TestLoadedContentCarriesTheExistingC25SourceName` **PASS**：`LoadedFile.Source == risk.SrcFSRead`（`"fs.read"`，名册 `internal/risk/provenance.go:84` 现名，未新造）；`CheckText` 在被盖 scope 命中；**变异反控内建**＝从未 Mark 的同文本 scope 不命中（＝删 `stamp()` 后的状态，仪器看得见盖戳的缺席）；`IsSensitiveSource(fs.read)` 为真 | 实现 `projctx.go:344`（stamp，`Prov==nil` 出声不静默） |
| AC#7 | 每轮打印"读了哪几份、各多少字节、截没截"；面板载体拿得到 | 三层全 **PASS**：①每轮打印＝`Turn()`→`Manifest()`（AC#1 用例内 `logged` 非空＋`ProjectInstructionManifest()` 点名文件）；②Go 载体＝`TestTheInstructionStatesDoNotShareOneWireShape:142`／`TestAPumpWithoutAnInstructionsReaderSendsNoKey:235`／`TestWithInstructionsKeepsTheR1CarrierContract:287`／`TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll:60` 四枚 **PASS**（五状态 `loaded/none_found/off/refused/not_run` marshal 互异、reader 接了必发键、改写视图拒读）；③**生产端到端**＝cmd/wisp 两枚 **PASS**（`TestRunPacketCarriesTheLoadedInstructionFiles 2.68s`：真工作区真 `AGENTS.md` ⇒ `status=loaded` `tier=project depth=0 bytes=269 source=fs.read`；`TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound 4.72s`：开关关 ⇒ `status=off`＋逐字跳过句、控制组 ⇒ `loaded`，两 packet 逐字节不同） | 载体 `internal/panel/composer.go:74`（`instructions,omitempty` 一枚对象）＋`instructions_200.go:27-40`（五状态）/`:120`/`:162`；泵读者 `internal/panel/pump.go`（`PumpSources.Instructions`）＋`cmd/wisp/run.go:947`/`:1705-1728`（装配根接线现读在位） |
| AC#8 | 一枚开关默认开；关掉要打印"已按你的配置跳过" | 两枚全 **PASS**：`TestOffSwitchReadsNothingAndSaysWhy:395`（`Enabled=false` ⇒ `Block==""`/`Files==0`/日志含"已按你的配置跳过"；正控＝同接线打开即非空）；端到端 `TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound`（见 AC#7③）。开关键 `internal/config/schema.go:449` `project_instructions_enabled` **default:"true"** 现量在位 | 实现 `projctx.go:255`（`SkipConfigOff` 早退，`:256` 逐字跳过句） |

## 2. 界面那一跳的现量（只读，`frontend/**` 零写面）

- `frontend/src/lib/panel.ts:252` 现量 **`instructions?: InstructionsSection;`** 已声明（`:187` 节注释、`:173` 元素注释同在）。
- ⇒ 票面那枚"待人一跳"已由界面那支交割：`internal/panel` 全包复跑（log `06-panel-full.txt`，`FAIL 2`）里
  **`TestApprovalCardViewJSONKeysMatchFrontendTypes` 与 `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` 已不在红名册**
  （r2 时代四枚红 ⇒ 今天两枚红）。转绿判据（票面 ：74-75）前半已兑现。
- 今天仍在册的两枚红**逐名归票**：`TestComposerContractTypesMatchFrontend`（红句 `Go ComposerState emits [files]` ⇒ **票 197** 的 `files` 键）与
  `TestC21DesignTokensFourWayAgree`（tokens 四方 ⇒ C21/票 344 一族）。**没有一枚归 200**，名册无新增。
- 面板"真画出五状态"归界面那支的活，本程零读面零证据，不代报。

## 3. 门禁终态（本程亲跑）

| 尺 | 读数 |
|---|---|
| `go test -count=1 ./internal/projctx/` | **ok 0.257s**（10 PASS／0 FAIL，log 01） |
| `go test -count=1 ./internal/agent/`（四枚 200 判据定向） | **4/4 PASS**（logs 02/07） |
| `go test -count=1 -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg ./internal/risk/` | **PASS**（log 03；200-r2 那格修复在 HEAD 仍绿） |
| `go test -count=1 ./internal/panel/` | FAIL 2（逐名归 197/C21，见 §2；定向 5 枚 200 判据 5/5 PASS，log 05） |
| `go test -count=1 -run 'TestRunPacketCarries…|TestRunPacketSaysOff…' ./cmd/wisp/` | **2/2 PASS**（log 04；包编译含 323 在飞件无碍） |
| `go build ./internal/... ./cmd/...` | **rc=0**。⚠ `go build ./...` 报 `build\mut_rvw.go` 11 枚 undefined——`build/` 在 .gitignore 里（`git status -- build/` 空、无历史），是**本机被忽略目录里的变异残渣**，不在仓、不归本票、不归任何在库面，已如实具名不处置 |
| `go vet ./...` | 同上一枚因 `build/` 残渣报错；internal/cmd 面干净 |
| gofumpt（`git ls-files -z '*.go' \| xargs -0 gofumpt -l`） | **0 枚未净**（log 10，rc=0） |
| `sh scripts/d22scan.sh` | **clean**（log 11：bans #1-8 全扫，internal 907／cmd 217 枚含注释与 _test.go） |

## 4. 结案判定

**PASS／可结**：票面八格 AC 在 HEAD `59ac99f9` 上**全部由本程亲跑转绿**，产码两枚 commit 在史、面在盘、行号现读对上。
残量三笔**均不阻结**、均有在册归属：①子代理环不带项目说明（归票 197，r2 证据 §6.3 具名）；②生产 C25 戳 R4 休眠（`run.go` `NoProbe` 一族，票 181/在册老账，本票只证盖戳动作与仪器牙齿）；③`cmd/wisp` 端到端未造 `refused/not_run` 两态（panel 层真加载器用例已钉，r2 证据 §6.2 具名）。**界面显示**那一半 AC#7 归界面那支（声明已落，`panel.ts:252`），Go 侧载体与端到端读数已齐。

翻框与 `-done` 动作（归编排者）：八格 `[ ]` 可全翻 → `git mv` 加 `-done` 后缀。
