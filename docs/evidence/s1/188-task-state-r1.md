# 188-task-state-r1 — 票 188 `AC#2`（后台任务状态维）写腿 `188-r1`：现量与停手上报

- 分工出处：`.scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md` §E
- 本程锚点：起手 `git log --oneline -1` = `38fc7c0e`；写这一格时 HEAD 已是 `cb14a8b2`
  （`180-a1` 在我干活期间提交，`git show --name-only cb14a8b2` 原文见下，它只动票 180 那一枚工单，未碰我的文件 ⇒ 按 `A385` 不算漂移，只点名）
- 结论一句话：**`AC#2` 那一格落在本单写的写面之外，撞钉预检也命中了负向钉 ⇒ 按共同规矩「停手具名报回」，本程零产码**。

## 0. 起手名册（`git status --porcelain` 逐枚原文抄录，2026-09-28 17:4x +08）

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
?? .scratch/wisp/.scratch/
?? .scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
?? .scratch/wisp/probes/139/accept-r1/
?? .scratch/wisp/probes/152/overlay-probe1-on-samppost.json
?? .scratch/wisp/probes/156/__pycache__/
?? .scratch/wisp/probes/156/mut-156-r2/asis.log
?? .scratch/wisp/probes/156/zero156-r4-head.sh
?? .scratch/wisp/probes/156/zero156-r4-work/
?? .scratch/wisp/probes/158/r2/
?? .scratch/wisp/probes/161/r2/__pycache__/
?? .scratch/wisp/probes/161/r2/ctl/
?? .scratch/wisp/probes/161/r5/negative-control/
?? .scratch/wisp/probes/161/r6/logs/flip-7.txt
?? .scratch/wisp/probes/162/r4/
?? .scratch/wisp/probes/162/v1-baseline-gotest.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start.txt
?? .scratch/wisp/probes/183/accept-v1/
?? .scratch/wisp/probes/185/c1/logs/d22scan-post-final.txt
?? .scratch/wisp/probes/33/r2/d22scan-final.txt
?? .scratch/wisp/probes/33/r2/gate-final.txt
?? .scratch/wisp/probes/33/r2/status-final.txt
?? .scratch/wisp/probes/33/r2/status-start.txt
?? .scratch/wisp/probes/999/
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
?? part1-state1-fixed.txt
?? part1-state1-pristine.txt
?? part1-state2-fixed.txt
?? part1-state2-pristine.txt
?? part2-nog6-fixed.txt
?? part2-nog6-pristine.txt
?? part3-stale-fixed.txt
?? part3-stale-pristine.txt
```

终态闸门口径按共同规矩那条：**终态等于起手名册（逐枚具名差集为空）**，不是"必须为空"（`A374`）。

## 1. git log / git show 原文（本程每一节都往这里追加）

### 1.1 起手锚 `38fc7c0e`

```
38fc7c0e ledger(A387-A388) 收 33-r3（入向听众进树＋词面尺改能力尺，我四把尺复跑）＋owner 提问后按默认动作派 145-r2
```

### 1.2 写第 0-1 节时树上多出的那一枚（别人家的活，逐枚 `git show --name-only` 现读）

```
$ git show --name-only --format="%H %s" HEAD
cb14a8b237382b32d4d2a442c9eddd63741d83cc 180-a1(片①): 票 180 Progress log 追加（AC 框一枚未勾，只读普查结论）

.scratch/wisp/issues/180-panel-width-is-a-config-field-with-zero-production-readers-so-changing-config-toml-has-no-visible-effect.md
```

## 2. 起手现量（`AC#2` 那一格今天到底缺什么）

| # | 问题 | 现读结论 | 凭据（`file:line`，本程现跑） |
|---|---|---|---|
| 1 | 本单写的写面 `internal/memory/**` 里有没有名为 `Record` 的载体 | **没有**。该包 `Record` 只出现在两个无关位置：方法名 `RecordProviderError`、导出信封字段 `Records` | `grep -rn "Record" --include=*.go internal/memory/` 全部命中：`internal/memory/dao_providerhealth.go:131`、`internal/memory/privacy.go:194` |
| 2 | 票面 现量 点名的 `Record`／`Look`／`Count` 是哪三个符号 | `*tools.TaskRoster` 的三个方法（票 188 现量表写的是 `internal/tools/task.go:110-152`，行号逐字对得上） | `internal/tools/task.go:110`（`type TaskRoster struct`）、`:116`（`NewTaskRoster`）、`:123`（`Record`）、`:139`（`Look`）、`:152`（`Count`） |
| 3 | 那一枚名册的**载体结构体**今天装了什么 | `TaskOutput` **两枚字段，零状态维**：`Text string`、`ArtifactPath string` ⇒ 票面"缺状态字段"这一格在 `internal/tools`，不在 `internal/memory` | `internal/tools/task.go:93`（`type TaskOutput struct`）、`:97`（`Text`）、`:101`（`ArtifactPath`） |
| 4 | 名册的写侧今天谁在写 | `internal/tools/task_backfill.go`（票 176 补的写腿），它自己写明**不改 `Record`**、last-writer-wins | `internal/tools/task_backfill.go:18`（"it does not change Record…"） |
| 5 | `internal/memory` 这一侧的后台任务记录有没有状态维 | **有字段、无枚举、词表非 D43**：`TaskLog.State string`，schema 是 `state TEXT NOT NULL` 且**刻意不建 CHECK**，DAO 只校验非空 | `internal/memory/models.go:58-70`（`TaskLog`，`:62` 是 `State`）、`internal/memory/schema.go:53-65`（`:57` state 列）、`internal/memory/schema.go:14-17`（枚举形状列走 DAO 校验、不走 SQL CHECK 的登记）、`internal/memory/models.go:185-197`（`validateTaskLog` 只要求非空） |
| 6 | 那一枚 `State` 的生产者是谁、写的什么词 | 两处，都在我的写面之外：`internal/agent/loop.go:983-998` `taskLogState`（`done`/`stuck`/`cancelled`/`interrupted`/`error`，`:983` 注释自称 "a free-text column"）、`cmd/wisp/run.go:651,655-670`（复用同一枚映射）；`running`/`succeeded` 是 memory 层与测试自己的词（`internal/memory/dao_tasklog.go:24` 注释、`internal/memory/dao_test.go:253,262,282,307`、`internal/memory/privacy_test.go:28,34`） | 同上 |
| 7 | 面板那一栏今天的来源 | **零**。`internal/panel/` 非测试码里 `task_log`／`TaskLog`／`TaskOutput`／background task 全部零命中，快照里没有任务这一维可显示（与票 188 现量"前端是示意数据"一致） | `grep -rniE "task_log\|tasklog\|TaskOutput\|background.?task" --include=*.go internal/panel/ \| grep -v _test.go` => 空输出 |
| 8 | `internal/panel/` 里 `roster` 这个词的命中是不是任务名册 | **不是**，那些命中全是 C17 方法白名单那枚名册（`composer_dispatch.go:160-178`、`l2_grant_boundary_test.go`），词面撞车，不构成任务来源 | 同上第 7 行的 grep 面 + `grep -rn "Roster" internal/panel/*.go` |

⇒ **本节结论**：票面 `AC#2` 那一格（"名册已在、缺字段"）落在 `internal/tools/task.go:93,123`，
而本单 §E 给的写面是 `internal/memory/**`。两者不重合的那一枚载体正好就是"缺字段"的那一枚（`TaskOutput`）。
`internal/memory` 这一侧不是"缺字段"，而是"字段已在、无枚举、词表与 D43 不一致，且两个生产者都在写面之外"。

## 3. D43 那批状态名：出处与"没新造"的证明

- 那张表在仓里的行号：标题 `docs/PLAN.md:3055`（`### D43 — 状态机权威转移表（**C12 冻结的就是这张表；docs/STATE_MACHINE.md**）`），
  40 条转移正文行＝ `docs/PLAN.md:3061`（`#1`，`FirstRun` -> `Sleeping`）到 `docs/PLAN.md:3100`（`#40`，`Stuck` -> `Acting`/`Settling`）。
  **本程对 `docs/PLAN.md` 零写字节**（见 §6 的门禁与终态名册）。
- 逐字对齐**不需要我在 Go 侧抄那张表**：仓里已经有那批名字的冻结拷贝，且它就是 C12 冻结物的实现——
  - `internal/statemachine/doc.go:1-2`："Package statemachine owns the BallState state vocabulary and the D43 transition table (SPEC-01 §3; C12 freezes the table in D43)"
  - `internal/statemachine/states.go:3`："State is one of the 20 BallState values (D43; C12 frozen)"
  - `internal/statemachine/states.go:8`："The 20 states, names exactly as in D43 / SPEC-08"
  - `internal/statemachine/states.go:11-31`＝那 20 枚常量（值逐字为）：`FirstRun` · `Sleeping` · `Armed` · `Muted` · `Listening` · `Thinking` · `Acting` · `Speaking` · `Warm` · `Conversation` · `Confirming` · `AwaitingApproval` · `Settling` · `Downloading` · `Unconfigured` · `NoNetwork` · `Error` · `Queued` · `Stuck` · `WatchdogAlert`
  - `internal/statemachine/states.go:39-49`＝`Valid()`，同一名册的第二处逐字列举（加注释那枚 `Idle` 是 `Sleeping` 的别名，`states.go:33-36`）
- 与 D43 表交叉核对：这 20 枚与 `PLAN.md:3061-3100` 各行 From/To 列出现的状态名**逐字相同**（`Confirming`/`AwaitingApproval` 在表里带 `(L1)`/`(L2)` 后缀，是档位限定词、不是名字的一部分；`Error(audio_device)` 这类带括号的是 D37 载荷写法）。
- ⇒ **本程没有新造任何状态名**（也确实一枚都没落，见 §5）。真要落的时候，"对齐 D43" 的引用对象是 `internal/statemachine` 这一集，
  而不是在 `internal/memory` 或 `internal/tools` 里再造一份 20 枚的名单——后者必然与 C12 那张表漂。

## 4. 撞钉预检（派腿前必跑那一发）

跑的是共同规矩那条命令形状，逐条读了命中断言。

| # | 我要新增／改动的符号 | 现跑命令 | 命中 | 判定 |
|---|---|---|---|---|
| 1 | 任务状态枚举（`TaskState`／`ValidState`／`StateVocab`）与 D43 那批专名 | `grep -rn "TaskState\|ValidState\|StateVocab\|FirstRun\|AwaitingApproval\|NoNetwork\|WatchdogAlert" --include=*_test.go internal/ cmd/` | `internal/ball/liquid_test.go:121-127,307-309`、`internal/ball/live_windows_test.go:87-94`、`internal/ball/tokens_test.go:121-128,300,318-325`、`internal/models/bridge_test.go:17-39`、`internal/models/handoff_window_109_test.go:62`、`internal/statemachine/table_test.go`（D43 在册） | 那批负向钉的射程＝**球/会话那 20 态名册本身**（枚数与名单两处正反都钉），**射程不盖 `internal/memory`**；但它同时钉住"C12 那张表只有 `internal/statemachine` 一份拷贝" ⇒ 我在 memory/tools 里再造一份 20 枚名单就是钉要打的形状。**引用可以、复制不行** |
| 2 | `task_log.state` 走 D43 白名单校验（即"补状态维"在 memory 侧的唯一落法） | `grep -rn "State: \"\|task_log.state" --include=*_test.go internal/memory/ internal/agent/ cmd/` | **`internal/memory/dao_test.go:262,265` 用 `"succeeded"` 且断言写成功**、`:253,282,307` 与 `internal/memory/privacy_test.go:28,34` 用 `running`/`interrupted`；`internal/agent/forensics_test.go:128-132` 钉 `task_log.state` 终值必须是 `cancelled`；`internal/agent/loop_golden_test.go:315-320` 钉必须是 `done` | **命中即红，且是负向钉**：`done`/`cancelled`/`succeeded`/`running` **一枚都不在 D43 那 20 个名字里**。要让状态维"逐字对齐 D43"就必然让这四组既有断言变红，其中两组在**我的写面**（`internal/memory/dao_test.go`）、两枚在**写面外**（`internal/agent/`）。按共同规矩"**不许自己放宽断言**" ⇒ **停手报回**，不自行改钉 |
| 3 | 给 `TaskOutput` 加字段（票面 `AC#2` 的字面落法） | `grep -rn "TaskOutput{\|NumField\|reflect.TypeOf" --include=*_test.go internal/tools/` | `internal/tools/fswrite_silentloss_ac1_test.go:226-232`（`structFieldsOf` 反射列字段名）、`internal/tools/task_output_*_test.go` 里 18 处 `roster.Record(..., TaskOutput{Text:…, ArtifactPath:…})` 字面量构造 | 加字段本身不撞那枚反射钉（它问的是 fs 那族行结构能不能装下字节数）；但 18 处**位置无关的字面量**＋`task.go:123 Record` 与 `task_backfill.go` 的写侧都在这程**写面之外**（本单 §E 只给了 `internal/memory/**`，且 `internal/tools/pointer_185_cli_seam_test.go` 已被点名为别家的腿）⇒ 不可自行落 |

## 5. 停手上报：三处前提冲突与需要的具名解冻

按 AGENTS.md §0.1（D22 闸门③「未定义即停」）与本单共同规矩"前提与票面/本单不符就停手具名报回，不要自己替我改规格"。本程**零产码、零判据、零字段**。

**冲突 1：写面与票面 现量 不重合，且不重合处正是那一格。**
票 188 的 现量 表第 2 行与 `AC#2` 说的 `Record`＝`*tools.TaskRoster` 的方法（`internal/tools/task.go:123`，票面行号段 `110-152` 逐字对得上），缺字段的是 `TaskOutput`（`task.go:93-102`）。
本单 §E 给的写面是 `internal/memory/**`，而 `internal/memory` 里**没有任何名为 `Record` 的载体**（现量见 §2 第 1 行）。
⇒ 编排者要二选一，且只能由你点名：
 (i) **具名解冻** `internal/tools/task.go`（只到"给 `TaskOutput` 新增状态字段"为止）与 `internal/tools/task_backfill.go`（那一枚 `Record` 的唯一生产写侧）；
 (ii) 或把 `AC#2` 的载体**改判**为 `internal/memory.TaskLog`——但那一格在 memory 侧不是"缺字段"而是"字段已在、无枚举、词表非 D43"（§2 第 5-6 行），票面那句 现量 要跟着改，且同批必须解冻 `internal/agent/loop.go:983-998` 那枚映射与 `cmd/wisp/run.go:651-670` 调用点（后者还在本单的 `cmd/wisp/**` 禁区内）。

**冲突 2：D43 是**会话/球**的状态机，票面要的是**后台任务**的状态；那张"任务态 -> D43 名字"映射表在规格里不存在。**
`PLAN.md:3061-3100` 那 40 行没有一行讲"一个后台任务处于哪一态"（与任务生命周期同形的只有 `Acting`/`Queued`/`Stuck`/`AwaitingApproval`/`Error` 这几枚，且它们是**会话**在干活时的态）。
而 memory 侧今天写的是 `done`/`cancelled`/`interrupted`/`succeeded`/`running` 这套任务完成度词（§2 第 5-6 行）。
⇒ "逐字对齐 D43"落到哪一枚映射，是**规格真空**（新增它＝新需求，同 `AC#1` 那一格被 `A389` 摆给 owner 的性质）。我不自造映射表，也不自造第 21 枚名字。

**冲突 3：`AC#3` 反向判据此刻钉不成活钉。**
现量：`internal/panel/` 非测试码里 `task_log`／`TaskLog`／`TaskOutput`／background task **全部零命中**（§2 第 7 行），今天"面板侧可写"这一维根本不成立，也没有状态维可钉；
而 `145-r2` 正在 `internal/panel/composer.go`／`pump.go` 加快照字段（本单顶部点名的在飞写腿）。
⇒ 这一判据要钉的是"**状态维真实落地的那一枚载体**没有面板写腿"，载体没定 = 只能交一枚扫空气的死钉（无正控 ⇒ 按 `A389`／"负向尺必配正控"那条我不交）。**载体定案后我再落 `AC#3`，且那一发要自带"种一发面板写腿必响"的正控。**

## 6. 门禁读数（原样，本程零产码 ⇒ 读的是"我没弄脏任何东西"）

- `sh scripts/d22scan.sh`（锚 `945a4d92`）：**clean - no D22 ban violations**。
  分母逐名（`ban #8 internal/`＝**文件枚数、不是违规数**）：
  `bans #1-5 internal/=211, #1-5 cmd/=24, ban #6 frontend/=85, ban #7 internal/tools/=21, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=441, ban #8 cmd/=47`；
  另 `examined 235 production Go files under internal/ and cmd/`、`skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`。
  ⚠ 票 188 `AC#6` 写的基线是 `ban #8 internal/` **examined=433**，现读 **441** ⇒ 别家加文件（含 145-r2／33-r3），按 `A384`/`A387` 口径**涨不等于漂移**。
  同一次跑里 d22scan 自带仪器：`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0`。
- `bash .scratch/wisp/probes/154/gate-clauses.sh`：**BAD 名册＝只 `G6neg`**（`声明=ring 基线=1枚 实测=3枚 因=新增未成对（票 171 AC#2）`），与任务书在册口径一致；
  其余 13 腿逐行 `ok`（`G1 0/0`、`G1b 0/0`、`G2 2/2`、`G3 0/0`、`G4 0/0`、`G5 1/1`、`G5pos 0/0`、`G5neg 8/8`、`G6 1/1`、`G6pos 0/0`、`G7 3/3`、`G7pos 0/0`、`G7neg 4/4`）。
  **`probes/161/r6/flip-declaration.sh` 一枚都没跑**（任务书禁跑，且起手名册里它那 9 枚 `flip-*.txt` 日志是别人在飞的脏件）。
- `"$(go env GOPATH)/bin/gofumpt" -l internal/memory/`：**空输出**（本程对 `internal/memory/` 零写，包形状本来就是干净的）。
- `go test -count=1 ./internal/memory/`：`ok github.com/CarlosShao/wisp/internal/memory 12.810s`。

## 7. 没测到什么（本程射程的洞，逐条具名）

1. **没测 `./internal/tools/` 与 `./internal/agent/`**：本程零产码，且那两包里有别家在飞的腿（`internal/tools/pointer_185_cli_seam_test.go` 被点名、`internal/agent` 是本单禁区）⇒ 票面 `AC#6` 那两包的读数**未采**，别当绿。
2. **没测 `./internal/panel/`**：三枚在册红（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）按任务书不修不当绿，且 `145-r2` 正在里面写 ⇒ 本程不进去量。
3. **没跑 `cmd/wisp` 相关**：要带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 才不报 `0xc0000135`（票 98），而 `cmd/wisp/**` 是本程写面之外、无产码可测。
4. **没有变异载体／正控**：本程零产码 ⇒ 没有任何"哪行改动→哪条红"的实发记录；§4 那三组负向钉是**读断言**得到的，不是我种出来的红。载体定案后落 `AC#2`/`AC#3` 时才需要那一栏。
5. **`TaskOutput` 到底该装哪几枚任务态**＝未定案（§5 冲突 2），我不替规格填。
6. 起手名册里的 30 枚 `design/**` 删除与 `docs/evidence/s1/152-...-accept-r1.md` 那处脏改**不是我动的**，本程一律不碰、不还原、不提交。

## 1.3 / 1.4 本程两枚提交的原文（逐枚 `git show --name-only`）

```
$ git log --oneline -1   (片① 之后)
500bc98e 188-r1(片①): 证据件骨架——起手名册逐枚抄录＋起手锚 38fc7c0e＋cb14a8b2 逐枚 show 原文

$ git show --name-only --format="%H %s" HEAD   (片①)
500bc98eb82a19d4debe729b9f3716dce471f65a 188-r1(片①): 证据件骨架——起手名册逐枚抄录＋起手锚 38fc7c0e＋cb14a8b2 逐枚 show 原文

docs/evidence/s1/188-task-state-r1.md

$ git log --oneline -1   (片② 之后)
945a4d92 188-r1(片②): 证据件 §2-§5 落格——现量（Record 在 internal/tools、memory 侧 State 已在非 D43）＋D43 名字出处（statemachine/states.go:11-31）＋撞钉预检三组命中＋停手上报
```
