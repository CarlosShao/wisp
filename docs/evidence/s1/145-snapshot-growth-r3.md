# 145-r3 — 快照长胖第三轮：`run.go` 那两行由派单具名放开之后，本程真落了一维

派单：`.scratch/wisp/dispatches/2026-09-28-181x-wave2-impl-145r3-and-188r2.md` §F（合并票 174 `AC#2b`）
起手锚：`81e77842`（≠ 编排者锚 `8c671afb`＋后续：起手那枚是编排者自己的 `A397`＋票 197 立票，未碰我写面）
终态 commit（本程）：`4db5f3f6`（格②）／`1ba16de0`（格①）／`0e5c0d4a`（格③）／末枚（格④台件）
时间：2026-09-28 18:12 进场 → 18:5x 交件　　写面：`cmd/wisp/run.go` 两行（具名射程）＋`internal/panel/{composer,pump}.go`＋`cmd/wisp/panel_pump.go`＋`internal/panel/composer_test.go` 追加

## 0. 起手名册（`git status --porcelain` 逐枚，18:12）

```
 M .gitignore
 M .scratch/wisp/issues/145-panel-snapshot-has-four-fields-so-eleven-of-the-fourteen-ui-states-have-no-input-to-render-grow-the-go-side-carrier.md
 M .scratch/wisp/issues/182-the-task-monitor-rail-is-outside-ticket-145s-fourteen-row-table-so-nobody-has-counted-which-of-its-stacks-have-a-go-side-source.md
 M .scratch/wisp/issues/188-the-task-monitor-rail-s-subagent-and-background-task-stacks-have-no-go-side-source.md
 M .scratch/wisp/issues/189-the-review-stack-needs-uncommitted-count-and-per-file-diff-but-go-side-has-zero-source.md
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt … flip-6.txt / flip-baseline.txt / flip-restored.txt（8 枚）
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/assets/base.css, icons.js, theme.js, tokens.css（4 枚）
 D design/index.html
 D design/screens/{approval,ball,chat,config,cost,firstrun,palette,privacy,security,states,tasks}.html（11 枚）
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
 M docs/reports/pending-and-issues.md
?? .scratch/wisp/.scratch/
?? .scratch/wisp/dispatches/2026-09-28-181x-wave2-impl-145r3-and-188r2.md
?? .scratch/wisp/issues/196-two-task-state-vocabularies-already-coexist-….md
?? .scratch/wisp/probes/{139/accept-r1, 158/r2, 161/r2/ctl, 161/r2/__pycache__, 161/r5/negative-control,
   161/r6/logs/flip-7.txt, 162/r4, 176/r1/logs/gate-clauses*.txt（4 枚）, 183/accept-v1, 185/c1/logs/d22scan-post-final.txt,
   33/r2/{d22scan-final,gate-final,status-final,status-start}.txt, 999/, 156/__pycache__, 156/mut-156-r2/asis.log,
   156/zero156-r4-head.sh, 156/zero156-r4-work/, 152/overlay-probe1-on-samppost.json, 162/v1-baseline-gotest.txt}
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/{lib/,rb-files.js,rb-plugins.js,rb-review.js,rb-terminal.js,rightbar.js,sidebar.js,screens/home.js,screenshots/}
?? design/old/
?? part1-state1-fixed.txt part1-state1-pristine.txt part1-state2-fixed.txt part1-state2-pristine.txt
?? part2-nog6-fixed.txt part2-nog6-pristine.txt part3-stale-fixed.txt part3-stale-pristine.txt
```

⚠ 起手名册**里没有** `174-*.md`（本程末枚把它追加进 Progress log 后随台件提交，终态差集＝`M → 已提交消失`，逐枚具名见 §6.4）。
⚠ `188/189/182` 那三枚 ticket 与 `docs/**` 的 ` M`／` D`／`??` 全是别人的活，本程一枚未动、一枚未提交（`design/**` 那 16 枚删除是 owner 账，不还原不提交不删）。

## 1. 派单与票面的前提，逐枚复量（不照抄）

| 前提 | 现量 | 判 |
|---|---|---|
| 唯一生产装配根＝`cmd/wisp/run.go:442` 的 `panel.NewSnapshotPump(panel.PumpSources{…})` | 起手 `sed -n '428,470p'` 复量：`Verdicts/Mode/Workspace/Git/Results/Out` 六根线，无 `Model`/`AttachmentMax`/`Now` | **对格** |
| 票 174 `AC#2b`＝`:365` 附近 `TaskDeps{}` 补 `Paths` | 起手复量真身在 `:365`（`for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks}) {`），判定者实例 `rt.paths` 造于 `:330`、已喂 `FSDeps`（`:346`）与桥（`:453 tools.Options.Paths`） | **对格** |
| `run.go:651-670` 那批状态词写入者＝票 196 射程 | 本程零字节（`git show --name-only` 三枚产码 commit 里 `run.go` 只有两枚 hunk，见 §3） | **守住** |
| 票 145 点 5「`cmdModels` 是模型清单」 | **错**（`A393` 现量：`cmd/wisp/models.go:205-225` 打的是 `store.manifest.Models`＝签名下载清单，`:206-208` 注释逐字 "not about this boot's [models] section"） | **本程不动那一格原文**，只在我那节 Progress log 具名更正（派单："具名不改写别的原文"） |
| `PumpSources.AttachmentMax` 能不能当那"一根真读口" | 全仓 `MaxAttachmentBytes` 只有 `attachments.go:47` 那枚常量＋`:187` 的默认分支，**无任何配置旋钮**（`grep` 非测试命中 4 行全在 `internal/panel/`） ⇒ 接进去＝把常量搬个家 | **拒**（装饰品） |

## 2. `AC#2`＋`AC#6`：落地集与拒收集

**落地集＝1 维 2 键**（嵌在既有 `composer` 段内，**顶层键集未增**）：

| 键 | Go 字段 | 真源 `file:line` | 抵达路径 | `AC#6` 答句 |
|---|---|---|---|---|
| `composer.currentModel` | `ComposerState.CurrentModel`（`internal/panel/composer.go` `Git` 之后） | `llm.Endpoint.Model`＝`internal/llm/resolver.go:46`；生产赋值 `cmd/wisp/run.go:308 rt.endpoint = ep`（`:314-315` 的 `rt.provs`/`rt.names` 同一枚 `ep`） | 读口 `rt.currentModel`（`cmd/wisp/panel_pump.go`，`gitView()` 之后）→ `PumpSources.Model`（`pump.go` `Git` 之后）→ 装配根 `run.go:442` 字面量 → `pump.go Snapshot()` 填进 composer 段 | `TestComposerCurrentModelTravelsOnlyFromItsReader`（`internal/panel/composer_test.go` 末节）三臂：有读口 ⇒ 逐字带出；无读口 ⇒ 空＋`modelKnown=false`；读到空串 ⇒ 空＋`false`（"未读"与"没有模型"两态不塌）＋出口字节含 `"currentModel":"glm-5"` |
| `composer.modelKnown` | `ComposerState.ModelKnown` | 派生自同一枚读口（`composer.CurrentModel != ""`），**不是第二份拷贝**：`currentModel` 单独存在时界面分不清"没人读"与"读到空" | 同上 | 同上（第三臂专钉这一对区分） |

**形状判据（为什么它俩不算顶层键）**：`Snapshot` 的 JSON 键集仍是 `composer,generatedAt,pending,results` 四枚——`pump_test.go:124`／`:276` 两条键集钉本程**一字未动、仍绿**（`go test -run 'TestSnapshotJSON|TestThePumpExitCarries|TestAPumpWithNoReaders' ./internal/panel/` ⇒ `ok 0.054s`）。派单那句"只有真落第五枚键解冻才算生效" ⇒ **条件未触发**，`A388` 的预批本程不追用（与 `145-r2` 同一结论，只是这次不是"没键"而是"键嵌在段里"）。

**拒收集（本程逐枚判，宁缺毋造优先）**：

| 候选 | 拒因 |
|---|---|
| `composer.models`／`composer.efforts`（票 145 `AC#2b` 两维） | 派单具名排除＝票 187 射程（档位词表三张全是包内私有 `var`、全仓无导出访问器）；**本程一枚未塞**，`AC#2b` 因此维持未勾 |
| `pending[].taskId`／`.callId`（`LiveApproval.Decision.TaskID`，生产填于 `internal/tools/bridge.go:331`、经 `approval/gate.go:550` 进队列） | 真源**有**、且是逐字拷贝（同 `Level/RulesHit/Reason` 那一族）——但 `ApprovalCardView` 的键集被**绿的**双向尺 `internal/panel/approval_test.go:105 TestApprovalCardViewJSONKeysMatchFrontendTypes` 钉着（它把 `ApprovalCardView`／`ResultChunk`／`Snapshot` 三枚都对 `frontend/src/lib/panel.ts` 逐键对账），加键＝**把一枚今天绿的钉打红** ⇒ 按"撞钉预检"停手：不在本程落，**报回见 §4 N1** |
| `pending[].position`（队列的 `LiveApproval.Position`＝C18 深度徽章） | 与数组下标 `+1` 恒等 ⇒ 与 `145-r2` 被拒的 `approval.depth＝len(pending)` 同形（复述字段当装饰） |
| `pending[].timeoutMs`／`windowMs`（`approval.Queue.Timeout()` `queue.go:126`，生产由 `run.go:380 cfg.Risk.ConfirmTimeoutSec` 喂；L1 窗口 `gate.go:142 Window()`←`:379 cfg.Risk.L1WindowSec`） | 源真且配置驱动，但**家不对**：`gate.go:416` 逐字"an L2 approval never opens an L1 window"，把窗口值盖到 L2 卡上＝造一形假相邻；`approval` 段是顶层键＝派单禁 |
| `results[].*`（reasoning／usage／tool 事件） | 数据进 `StreamLog` 的唯一地点是 `run.go:800-816` 那三支 case＝**本程禁写的记录点**（`:800-806` 注释自己就写着"the reasoning field is ticket 145 row 3, which is a key the snapshot does not have"） |
| `tools[]`／`cost`／`failures[]`／`view` | 沿 `145-r2` §2 的现量未变：记录点缺、单位口径未裁（`internal/agent/cost.go:16-20` micro-USD vs CNY）、Go 侧无九枚屏 id；且全是顶层键 |

⇒ **本程落 1 维**，其余**一枚未落**。这一维不是凑数：它答得出真源行号，也答得出"为什么 `currentModel` 不是被排除的 `models` 清单"（一枚是"这台机器现在真在用哪枚模型"，源＝运行时那枚 endpoint；一枚是"可以挑哪几枚"，源＝目录∩`Enabled` 且今天无导出面）。

## 3. `run.go` 逐行前后（射程＝派单具名的两行）

```
:365  - for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks}) {
:365  + for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}) {

:430  	Workspace: rt.workspaceView,
:431  	Git:       rt.gitView,
:431+ 	Model:     rt.currentModel,
:432  	Results:   rt.stream.Chunks,
```
两枚都**只改字面量本身**：零新增注释、零逻辑、`run.go` 其余行（含 `:651-670` 状态词写入者、`:800-816` 记录点）一字未动。
复算尺：`git show 4db5f3f6 --numstat`＝`1 1 cmd/wisp/run.go`；`git show 1ba16de0 --numstat`＝`run.go 1 0`（其余三枚文件见 §6.4）。

票 174 `AC#2b` 的**未落半句**照实登记：接上之后精确文案仍只给一条回来的路（加 `[fs] allowed_dirs`），缺 `174-v1` 条件②要的"**或批一张 L2 卡**"那半句，而那两句受 `internal/tools/task.go:299`／`:349` 两枚模板冻结钉系着＝`Q-63` 边界，派单具名"我自己先不碰" ⇒ **`AC#2b` 本程只到装配那一行，未结题**。

## 4. 停手报回／欠账（本程一律未擅自办）

- **N1 · `pending[]` 那一族的用例之家与钉**：`taskId`/`callId` 真源齐全（`bridge.go:331`→`gate.go:550`→`pending_read.go:46 LiveApproval`），但会被**绿的** `approval_test.go:105` 双向尺打红。要么 owner 让 `frontend/src/lib/panel.ts` 与 Go 同枚 commit（`Q-51`），要么具名解冻那枚钉。**本程没动它，也没为了避开它去改 composer 段的形状。**
- **N2 · `AC#6` 的用例之家仍是裁量**：本程把断言追加进 `internal/panel/composer_test.go`（非三枚冻结件、非 `pump_test.go` 那两行，只**追加**一枚 `TestXxx`、零删除）。口径依据＝`A389` 写面写的是 `internal/panel/**`；若编排者判"`A388` 那三枚文件才是面"，这一枚属**越界自证**，请具名退回而不是放宽断言。
- **N3 · 正控未跑**：派单的正控义务挂在"落了第五枚顶层键"那一支，本程未落 ⇒ 未跑；**新增那一维的变异正控也没跑**（`-overlay` 禁用、仓外副本跑 `internal/panel` 要带 `frontend/dist`，本程预算到顶）。牙齿从代码上自证：`pump.go` 那段 `if p.src.Model != nil` 摘掉 ⇒ 第一臂 `composer.currentModel = ""` 必红。**〔仅自述，未实测〕**

## 5. 前端"应当长这样"逐键清单（零 `frontend/**` 写面；由 owner 自己带）

| TS 键（挂 `PanelSnapshot.composer` 下） | 类型 | 现状 | 含义与显示纪律 |
|---|---|---|---|
| `currentModel` | `string` | **本程已发**（Go 侧真值） | 这一轮真在用哪枚模型；`modelKnown=false` 时**不许**画成"未配置模型" |
| `modelKnown` | `boolean` | **本程已发** | "没人读过" vs "读到空"两态的区分位；缺省 `false` |
| `models` | `ComposerModelOption[]`（`{provider, modelId, display}`） | 仍未发（票 187） | 空数组＝按钮不画，**不许**拿写死模型名顶 |
| `efforts` | `string[]` | 仍未发（词表无导出面，票 187） | 空＝"档位这一维不画"，与"只有 off"必须在文案上分开 |
| `effortVerified` | `boolean` | 仍未发 | 声明位 vs `provider_health` 核实位；**不许**把"未探测"画成"不支持" |
| `pending[].taskId`／`.callId` | `string` | 源已现量、**未发**（N1 的钉） | 卡片归到哪一路任务；前端目前自己按数组下标算徽章 |

⚠ `TestComposerContractTypesMatchFrontend` 因本程两枚键**更红**（它拿 Go 结构体逐键对 `frontend/src/lib/panel.ts`）＝派单与 `A388` 预先判成的"预期后果"，本程不修、不放宽、不动前端。

## 6. 门禁读数（原样；口径＝顶层 `--- FAIL` 行计数，与 `145-r2` 的 `-v` 口径不同，别当同尺复算）

### 6.1 改前（起手 `81e77842`，18:12）
- `./internal/panel/` **rc=1**、`--- FAIL` **3**：`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`
- `./cmd/wisp/` **rc=0**、FAIL 0（DLL 进 `PATH` 直跑）　`./internal/tools/` **rc=0**、FAIL 0

### 6.2 改后（三枚产码 commit 之后同一把尺）
- `./internal/panel/` **rc=1**、`--- FAIL` **3** ⇒ **名册逐名与改前相同**（未新增红、未吞红；那三枚**不修、不当绿**）
- `./cmd/wisp/` **rc=0**　`./internal/tools/` **rc=0**
- 新增用例单跑：`go test -count=1 -run TestComposerCurrentModelTravelsOnlyFromItsReader ./internal/panel/` ⇒ `ok 0.039s`
- `go build ./...` 净、`go vet ./internal/panel/ ./cmd/wisp/` 净、`gofumpt -l internal/panel/ cmd/wisp/` **空**

### 6.3 其余两把
- `sh scripts/d22scan.sh` ⇒ **clean、rc=0**（bans #1-5 `internal/=211`＋`cmd/=24`、ban #6 `frontend/=85`、ban #7 `internal/tools/=21`、ban #8 `design/=39`／`frontend/=85`／`internal/=442`／`cmd/=47`；`internal/` 442 与起手 188-r2 同数＝涨的是别人加的文件，不是违规）
- `bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ BAD 名册**仍只有 `G6neg`**（`腿数＝14 声明与实测不符＝1`）；那格 `基线=1枚 实测=3枚` 是 `188-r2` 已具名登记的既有涨法，**非本程所加**。`flip-declaration.sh` 未跑（派单禁）。

### 6.4 本程 commit 集（逐枚 `--name-only` 复算）
| commit | 路径 |
|---|---|
| `4db5f3f6` | `cmd/wisp/run.go`（1 行，格②） |
| `1ba16de0` | `internal/panel/composer.go`／`internal/panel/pump.go`／`cmd/wisp/panel_pump.go`／`cmd/wisp/run.go`（格①） |
| `0e5c0d4a` | `internal/panel/composer_test.go`（格③判据） |
| 末枚 | 本件＋票 145／174 Progress log（格④台件） |

禁面复算：`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／`internal/panel/bridge.go:42-45`／三枚冻结测试件（`tokens_fourway`／`l2_grant`／`frontend_hygiene`）＝**零字节**；`frontend/**`／`design/**`＝**零写面**；仓内零删除命令；只 commit、**未 push**。

## 7. 本程**没**测什么（不假装核过）

1. 没跑端到端 `wisp run` 验票 174 `AC#2b` 的"接上后必须给出精确文案、安静正例仍然安静"——那要真起宿主＋真 spill 文件，本程只到装配行；`174-c2` 已量过"只动那一行 ⇒ 既有 Go 判据零枚会红"，本程**未复算**这条。
2. 没跑变异正控（§4 N3）。
3. 没测 `rt.endpoint.Model` 在未配置角色时的形状（`run.go` 在 `:305` 那支就 `return rt, 2`，泵根本不装配，所以"空字符串＋known=false"那一臂在**生产里不可达**，只有测试可达）。
4. 没测多角色／多端点（今天只有一枚 `ep`，`rt.names` 一枚元素）。
5. 没核 `frontend/src/lib/panel.ts` 那侧现在声明了哪些键（`frontend/**` 零写面且 `A392` 裁"票面为准"）⇒ §5 那张表是**建议形状**、不是对账结果。
6. `./internal/risk/` 未跑（本程零命中那枚包）。

## 8. 纪律面

- **超预算**：派单硬顶 **≤35 枚**，本程实际 **≈39 枚**（具名被谁吃掉：起手名册＋波次派单＋台账三枚 4 枚、票 145/174 与 `run.go`/`pump.go`/`composer.go`/`panel_pump.go` 现读 6 枚、r2 证据件三节 3 枚、落地集选型（`pending_read.go`/`queue.go`/`gate.go`/`bridge.go`/`resolver.go` 真源逐枚核）5 枚、撞钉预检与名册钉 3 枚、改前门禁 1 枚、产码＋判据编辑 6 枚、门禁＋commit 2 枚、台件 4 枚）。**未据此放宽任何断言**。
- 被拒调用：**0 枚**（无一次工具拒绝）。
- 每次 commit 带显式 pathspec，起手 `git diff --cached --name-only`＝空、全程未动别家暂存态；无 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`。
- 票面 `AC` 框**一枚未勾**（`AC#2`／`AC#6`／`AC#2b`／票 174 `AC#2b` 全部维持未勾；勾要非实现者表）。

## 9. 一句人话

面板那条数据通道今天多送了一件真东西：**这一轮回答你的是哪一枚模型**，以及"这件事到底有没有人读过"。上一程想送十三件，被自己那把"答不出来源就不许加"的尺全挡了；这一程来源能一路追到代码行、且不会撞坏任何一条今天还绿着的测试的，只剩这一件，就只加了这一件。可选模型清单和"思考档位"那两样还是给不了——不是忘了，是那两张词表在 Go 里是关在包里的私有条目，谁也没法从外面读（那是另一张票 187 的活）。
