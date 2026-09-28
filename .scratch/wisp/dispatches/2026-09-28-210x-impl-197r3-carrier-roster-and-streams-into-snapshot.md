# 派单 2026-09-28 21:0x — 写腿 `197-r3`：**载体层**——把"名册里有哪些子代理、每枚跑到哪、它自己那一页在流什么"真送进面板快照

> owner 09-28 18:0x 原话逐字（票 197 顶部也挂着这句）：「**子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面，能明白吗？这是主流 harness 必做的，不要偷懒**」
> 加上他今天那句「**必须特么做完整功能，明白吗？**」⇒ **这一腿不是做样子**：前两腿已经把**实体**（`task.spawn` 派真任务、进 `TaskRoster`）与**流**（`StreamLog` 按 `subagent:<taskID>` 分键）做完了，
> **今天还没有人把它们搬到用户看得见的那条通道上**。这一腿就是那一跳。

## 0. 起手必做（跳了＝没交件）

1. `date "+%Y-%m-%d %H:%M %z"`＋`git log --oneline -10`＋`git status --short` 自取锚点；**本单行号是我 20:5x 读的，会漂，一律现读**。
2. **撞钉预检（按今天新立的规矩：不许只 grep 标识名，要跑包、抄今天的绿用例名）**：
   `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./internal/panel/ ./internal/tools/ ./cmd/wisp/`
   ⚠ 不带那枚 PATH 会 `exit status 0xc0000135`（**一枚用例都不跑**，看着像包坏＝仪器瞎）。
   **起手在册的红（逐名，我不修你也不许修/不许 Skip/不许放宽）**：`internal/panel` 四枚——
   `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／
   `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourwayAgree`；`internal/risk`、`internal/tools`、`cmd/wisp` 今天全绿。
3. 读：票 `.scratch/wisp/issues/197-*.md`、`.scratch/wisp/issues/211-*.md`；证据件 `docs/evidence/s1/197-subagent-entity-r1.md`、`-r1b.md`、`197-subagent-stream-r2.md`；
   以及 `docs/reports/survey-2026-09-28-dsh-ui-packages.md`（**它给了"点进去那一页"的现量与数据形状**）。

## 1. 硬规矩

- **只 commit、绝不 push**；显式 pathspec；禁 `git add -A`/`.`/`-a`；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`restore`/`clean`/`merge`/`worktree`；
  新文件逐枚 `git add -- <路径>`；**仓内绝不删东西**，写坏了只准 `git cat-file blob HEAD:<path> > <path>`。
- **零写面**：`frontend/**`、`design/**`（不写不读不转述）、`internal/panel/tokens_fourway_test.go`、`internal/tools/bridge.go`、`internal/agent/budgets.go`、
  `thresholds.go`、golden、`allowlist.txt`、`PLAN.md`、`docs/specs/**`。
- ⛔ **不许新增 C17 `PanelBridge` 方法名**（那是契约面）。如果你判断"点进去那一页"必须有**入向**方法（例如"把某会话装进某格"），**停手，在证据件里把那一跳写清楚交给我**——
  我落一条 `A##` 并把界面侧的话写进票面，由 owner 自己带给他用的那支界面 agent。**我不问 owner 要"做不做"，只把风险翻成票里的实现约束。**
- 别人的脏文件不碰不提交：`.gitignore`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`design/**` 那 16 枚删除、`.scratch/wisp/.scratch/**`（一枚写坏路径留下的双层目录）。
- 中文长段用 Edit/Write 工具，**别用 shell heredoc**（反引号会被真执行）；写完跟踪文件用 `git diff --numstat` 自证**删除列＝0**。
- `gofumpt` 在 `"$GOPATH/bin/gofumpt.exe"`（本机 v0.12.0），**不在 PATH**；`gofmt -l` 空 ≠ gofumpt 空。今天跟踪集里 `internal/`＋`cmd/` 是 **0 枚未净**，交件时**必须还是 0**。

## 2. 这一腿做什么

### (a) 名册进快照：用户在面板上能列出"现在有哪几枚子代理、每枚什么状态、是谁派出来的"

- 现读落点：`TaskRoster`（`internal/tools/task.go`）已有身份字段 `ParentTaskID`／`Label`／`Kind`／`State`，状态维**只用 D43 那 20 枚名字**（`internal/statemachine`，票 196 那族钉不许绕）；
  快照与泵在 `internal/panel/`（`Snapshot`、`pump.go`）、装配根在 `cmd/wisp/run.go` 的 `PumpSources` 字面量那一带（今天已被票 200 加过一行读者，**照那个形状接**）。
- **判据必须给到"生产路径"**：不要只在包内构造一个 roster 就宣称面板能看到——
  要有一条**真装配**的读数（`cmd/wisp` 那层现成的 harness），证明"派一枚子代理 ⇒ 快照里出现它、带它的 `taskID`／状态／父任务"。
  ⚠ 今天的现量是 `ResultChunk` 只有 `correlationId`／`text`／`done` 三枚键，**"这行是谁的"不在线上**；名册那一侧要能回答这个问题（不必靠给 `ResultChunk` 加键来解决——加键会牵动在册那两枚名册钉，慎）。
- **新增的快照键会打红 `TestApprovalCardViewJSONKeysMatchFrontendTypes`**（那把尺判的就是"Go 发了界面没声明的键"）。
  **这是它该做的工，不许绕**：不改尺、不加豁免、不为了过尺而少发一个键。红因与"界面那支要补哪一行、逐字段形状"写进你的报告，我转写进票面。
  ⚠ **顶层键那一格我 21:0x 现读了才敢写给你**（别信我上一段的话）：`internal/panel/pump_test.go:123` 与 `:291` 两枚断言**逐字**写着
  `"composer,generatedAt,pending,results"`——**只有四枚、且不含 `instructions`**。它俩今天仍绿，是因为票 200 那枚载体是
  `Instructions *InstructionsSection \`json:"instructions,omitempty"\``（`composer.go:74`），**fixture 里没填 ⇒ 键根本不出现**。
  ⇒ **两件事**：① 你要加的东西**默认应该嵌进现有键**（票 145 当年定的就是"不加顶层键"）；
  ② 如果你确实要一枚**无条件出现**的顶层键，那**这两枚钉必红**——**别改它们**，把"为什么必须破这一条"写成一段具名理由交我裁。
  ⚠ 顺带一枚**已经存在的仪器洞**（不是你的活，但你要是撞上就具名报）：这两枚"四枚顶层键"的钉**会因为 `omitempty` 而看不见新加的键**——
  也就是说"加了顶层键但那条路径今天没值"这一形，这两枚钉**抓不到**，只有 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 那种"从真 marshal 的字节里数键"的尺才抓得到。

### (b) 每枚子代理的流进快照：点进去那页有东西可渲染

- `internal/panel/pump.go` 的 `StreamLog.Append/Chunks` 与 `DefaultStreamKeys`（今天＝池 4 同量级）已就绪；
  流键形状两包各有一枚同名常量（`SubagentStreamKeyPrefix`），**panel 侧是真源、`cmd/wisp` 注入**（账 `A406` 已裁），
  这一腿**把两包重复收敛成一枚真源**，同时**保留/新增那枚"两包各自造同一个键＝红"的常驻判据**（`cmd/wisp/subagent_stream_key_197_test.go` 已有两枚钉，别把牙拆了）。
- **溢出语义今天已是"显式截断"**（保留 256 rune、硬上限是其倍数；合并那条已删）⇒ **把截断标记的读者在真装配里接上**，
  并给一条判据：**被截断时快照里必须能看出"这里少了 N 段/多少字节"**，不许静默丢。

### (c) 一票一票把状态词钉住

- 名册里的状态**只允许 D43 的名字**（`statemachine.Valid` 是判据）；**不许新造第二套状态词**（票 196 就是防这个）。
- 取消语义：**取消不级联**（父任务取消不连带掐孩子）——这是 `A401` 的裁定，若你测出来今天的行为与它不符，**具名报我**，别自己改语义。

## 3. 交件（三样）

1. 代码＋测试一次或两次 commit（message 带 `197-r3` 与票号 197）。
2. 证据件 `docs/evidence/s1/197-subagent-carrier-r3.md`，六节：① 起手锚点与**今天的绿用例名册**（逐包点数，不许"只列表不点数"）；
   ② (a)(b)(c) 各自改前→改后读数＋**落点逐跳行号**；③ **至少两发变异/正控**（例：把泵那根线拔掉 ⇒ "快照里看不到子代理"必须真红；把状态词换成一个非 D43 的名字 ⇒ 那枚钉必须真红），
   并写清**"哪一枚 commit 才算未修码"**；④ 门禁四数＋`internal/panel` 红名册逐名比对（**多了哪一枚、少了哪一枚**）＋全仓 gofumpt 名单；
   ⑤ **我没测什么**（具名）；⑥ **对派单的不服**——我给的行号、常量、"当年定的是什么"若有错，当面写出来（今天这条链上已经被写手推翻四次，全是我自己的账）。
3. 报告只报盘上可核的事实：commit 号、`wc -c`、判据逐名读数。我不认"应该已经"。
