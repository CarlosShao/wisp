# 197-r4 证据件 — 票 197 `AC#4` 的 `blockedOnApproval`：生产路径判据（产码零改动）

> 派单：编排者 2026-09-28 深夜口头派（代号 `197-r4`），范围写死"给票 197 的 `AC#4` 补上生产路径判据，
> **产码零改动、只写测试**"。依据＝`docs/evidence/s1/197-subagent-carrier-r3.md` §⑤1 与 §⑥（前一腿
> 197-r3 那句"这棵树里没有任何办法从外面造出那枚卡"被 `197-r3b` 判为**讲重了**）。
> 生成时刻 `2026-09-28 22:3x +08`。**本文件只写盘上可核的事实**；所有读数出自本腿亲手跑的命令，
> 探针输出全在 `.scratch/wisp/probes/197/r4/`。

---

## ① 起手锚点与绿名册点数（本腿现取，未采信派单里的任何行号）

| 项 | 读数 | 尺 |
|---|---|---|
| 起手 HEAD | **`da1d9289`** | `git rev-parse HEAD` |
| 起手工作树里的本票产码 | **0 枚**（本腿没有动过任何产码文件；`git status --porcelain -- cmd internal` 起手为空） | 同上 |
| 起手包读数：`internal/tools` | `rc=0` **PASS=158 FAIL=0 SKIP=0** | `go test -count=1 -v`（带 sherpa/build PATH），`.scratch/wisp/probes/197/r4/start-v-tools.txt` |
| 起手包读数：`internal/panel` | `rc=1` **PASS=95 FAIL=4 SKIP=0**＝在册那 4 枚（逐名见下） | 同上，`start-v-panel.txt` |
| 起手包读数：`cmd/wisp` | `rc=0` **PASS=96 FAIL=0 SKIP=0**（**整包 `-v`，不是 `-run` 单跑**） | 同上，`start-v-cmdwisp.txt` |
| 起手 `go build ./...` / `go vet ./cmd/wisp ./internal/panel ./internal/tools` | **rc=0 / rc=0**，输出空 | `start-gate.txt` |

起手在册红逐名（`grep '^--- FAIL' start-v-panel.txt`）＝**4 枚，与派单给的名册一字不差**：
`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／
`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourwayAgree`。
本腿**一未修、二未 Skip、三未放宽**（终态逐名比对见 §④）。

**与派单给的那三个"终态读数"的出入**：派单说 `cmd/wisp 整包 PASS=96`，本腿起手复跑就是 96（一致）；
`internal/tools PASS=158`（一致）；`internal/panel 红＝4 枚`（一致）。**没有出入，本腿照记。**

---

## ② 落点逐跳行号 ＋ 改前／改后逐字读数

⚠ 行号一律**本腿现读**（`grep -n` 跑在 `da1d9289` 这棵树上），不抄派单、不抄前一腿证据件。
派单给的两枚行号本腿复核为**准**：`internal/panel/composer.go:91`（`Tasks *TaskRosterSection`）、
`cmd/wisp/run.go:577`/`:579`（`AdmitTextTask`＋`PendingApproval`，落在 `confirmModeSwitch` 那枚函数里，
函数本体从 `run.go:572` 起）。

**判据本体**＝`cmd/wisp/subagent_blocked_197_test.go:59 TestRunPacketMarksTheRosterRowACardIsHolding`
（本腿唯一新增的件，251 行含注释：用例一枚＋helper `waitChildRow197` 一枚＋两枚常量；
**产码零改动**）。

### (a) 那一态是怎么造出来的（逐跳）

| 跳 | 落点（现读行号） | 这一跳在链条里是什么 |
|---|---|---|
| 1 | `cmd/wisp/subagent_carrier_197_test.go:151 spawn197.launch`（复用前一腿的 helper，**没改它**） | 孩子经**本 run 自己装配的桥**派生：`rt.bridge.Execute(... task.spawn ...)`，`TaskID == CorrelationID == 根的任务 id`（`internal/agent/loop.go:647` 的配对） |
| 2 | `cmd/wisp/subagent_blocked_197_test.go:107 waitChildRow197` → `cmd/wisp/subagent_carrier_197_test.go:215 childRowAfterJoin197` → `cmd/wisp/panel_pump.go:146 taskRosterState` | 孩子的 task id 从**这枚 run 的名册读者**里取，不是从测试自己手里的变量取；行在孩子的第一次模型调用之前就存在（`internal/tools/task.go:337 PublishSubagent`），所以读到的是**在飞**的那一瞬（`inFlightSlots:1`、`status:Thinking`） |
| 3 | `cmd/wisp/subagent_blocked_197_test.go:118`＝`rt.gate.AdmitTextTask(id)` | **宿主自己的入口**，与 `cmd/wisp/run.go:577`（模式切换的 R20/M4 确认）同一枚；D47 的闸门就认这个登记 |
| 4 | `cmd/wisp/subagent_blocked_197_test.go:120`＝`rt.gate.PendingApproval(ctxCard, tools.Decision{TaskID: id, CorrelationID: id, ...})` | 同一族另一枚（`run.go:579`）。真卡进真队列：`internal/agent/approval/queue.go:141 push` ⇒ `it.Corr = d.CorrelationID`，`internal/agent/approval/gate.go:484 promptFor` → `:486 ui.Prompt` → `cmd/wisp/run.go:1012 consoleApprovalUI.Prompt` |
| 5 | `cmd/wisp/run.go:1030`（`Prompt` 末尾的 `if u.publish != nil { u.publish() }`）→ `cmd/wisp/run.go:494 rt.ui.publish = rt.publishPanelSnapshot` | **这一发是 run 自己发布的**：卡上屏的那一刻泵被触发，本腿在此之前**一次 `publishPanelSnapshot` 都没调用**。发布＝`internal/panel/pump.go:284 Publish`→`:263 TaskRosterSectionFrom(p.src.Tasks(), cards)`→`internal/panel/subagent_roster_197.go:193/210`（`waiting[card.CorrelationID]` 与 `BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`）→ 出口 `cmd/wisp/run.go:492 Out: rt.bookPanelSnapshot`（账本记 sha256） |
| 6 | `cmd/wisp/subagent_blocked_197_test.go:211 cancelCard()` → `internal/agent/approval/gate.go:529`（`ctx.Done` ⇒ `q.abandon`） | 撤卡那一半：队列里的条目真离开，第二次读数由本腿调用 run 自己的发布者取（`publishPanelSnapshot`，具名"这一发是本腿采的样"） |

### (b) 改前 ＝ 今天盘上就有的"没有判据"那一形

`AC#4` 的**这一格**在 `da1d9289` 上的状态不是"码没写"，是**码在、判据不在**：
`internal/panel/composer.go:91` 的 `tasks` 节、`internal/panel/subagent_roster_197.go:127` 的
`blockedOnApproval` 字段、`:210` 的 join、`cmd/wisp/panel_pump.go` 的名册读者**全都在产码里**，
生产路径上唯一读到过它的方式是**断言它为 false**（`cmd/wisp/subagent_carrier_197_test.go:311`），
`true` 那一支只在 `internal/panel` 用一枚手搭的 `NativeVerdict` 钉着（前一腿 §⑤1 自己写的）。
⇒ **"改前"的逐字读数**＝本腿 §③-M1/M2 两发单点变异在同一棵树上打出来的字节：
`"blockedOnApproval":false` 而**同一枚包里 `"pending"` 就带着那枚卡**（那是判据瞎掉时生产路径的真形）。

### (c) 改后（生产路径，**整包 `-v`** 的终态读数＝`final-v-cmdwisp.txt` 第 644-647 行）

判据绿，且两半读数都在同一枚用例里：

```
subagent_blocked_197_test.go:168: tasks wire bytes: {"rows":[{"taskId":"27e42e4d-…","label":"总结一下 rootprompt197r3 aaa…","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"27e42e4d-…","blockedOnApproval":false,…},{"taskId":"9fd9cac4-…","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"27e42e4d-…","status":"Thinking","statusKnown":true,"streamKey":"subagent:9fd9cac4-…","blockedOnApproval":true,…}],"inFlightSlots":1,"poolCap":4,…}
subagent_blocked_197_test.go:222: tasks wire bytes: {…同两行：根 blockedOnApproval:false；孩子 "status":"Settling","statusKnown":true,"blockedOnApproval":false…],"inFlightSlots":0,…}
subagent_blocked_197_test.go:234: blocked row on the run's own packet: task=9fd9cac4-2128-450f-9bf5-61a63b0f1261 status=Thinking streamKey=subagent:9fd9cac4-2128-450f-9bf5-61a63b0f1261 card=ticket197.blocked.probe pending=1 bytes=3091 sha256=7d6bf5b43703bd82 | after the card: blocked=false answer=reject why="任务上下文已结束，审批请求已作废并按拒绝处理"
--- PASS: TestRunPacketMarksTheRosterRowACardIsHolding (2.14s)
```

（第一行是 `-v` 打出来的**原始字节**，`packetTasks197` 把 `tasks` 那节的 `raw` 整段打进日志；
`…` 与第二行的 `{…}` 是本文件为排版省略的**同一行内**其余字段，逐字节原文在那一发文件的第 644-645 行。
⚠ 每发 `-v` 的 uuid 与 sha256 都不同：本腿另两遍同内容的读数是 `d2d0d49f5a2c109d`
（`run1-single-newcase.txt`，`-run` **单跑**，只算中间读数）与上面这发 `7d6bf5b43703bd82`（**整包**）——
**"绿"只引整包那一发**，理由见 §③ 末的台件口径与派单的仪器纪律。）

三件派单点名的断言，逐条在场：
1. `"blockedOnApproval":true` 出现在**这枚 run 自己发布**的包上（`pending` 里就是那枚卡，且本腿在那之前没发过任何一包）；
2. **同一行的 `status` 是 D43 的名字**：`"Thinking"`（在飞那一瞬）与 `"Settling"`（收口之后），
   判据是 `statemachine.Valid`（`subagent_blocked_197_test.go:176`）——本腿没有换词表、没有把在飞值钉成常量；
3. **`streamKey` 形状不变**：`subagent:9fd9cac4-…`，与 `panel.SubagentStreamKey(id)` 比（`:180`），
   单一铸造点那枚扫描（`TestSubagentStreamKeyHasOneMintSite`）终态照绿（§④）。

**顺带钉住的三枚不该动的东西**：
- **只有被点名的那一行翻**。同一枚包里根行 `blockedOnApproval:false`（`:190`），
  而 `ApprovalCardView`（`internal/panel/approval.go:39-59`）**根本没有 task id 字段**——
  所以"按卡上的任务 id join"这一形在包里无法表达，判据只能是 `correlationId == 行的 taskId`。
- **载体不是批准出口**：撤卡后 `answer=reject`、`why="任务上下文已结束，审批请求已作废并按拒绝处理"`
  （`internal/agent/approval/gate.go:529` 的 abandon 分支 ⇒ `queue.go:427` 返 `AnswerReject`），
  用例把它也断言了（`:230`）；票 197 `AC#5`／`Q-49 丙` 不因本腿多了一枚能答卡的口。
- **顶层键一枚没加**、`pump_test.go:123`/`:291` 那两枚"四枚顶层键"的钉一字未动（终态读数 §④）。

---

## ③ 正控：造得出那一态 ⇒ 判据绿；拔掉读者／把卡撤掉 ⇒ 必须真红

台件：`.scratch/wisp/probes/197/r4/mutate-197r4.sh`（单点、串行、每发跑完立刻
`git cat-file blob HEAD:<path> > <path>` 还原，再 `git status --porcelain -- <path>` 自证）。
逐发输出＝同目录 `m0-*`／`m1-*`／…／`m5-*.txt`，汇总＝`mutation-log.txt`。
⚠ **口径**：这一节的每一发都是 `-run` 过滤跑（为了把变异台件控制在几分钟内），**只算中间读数**；
"绿"的终态声称来自 §④ 的整包 `-v` 点数。前一腿 197-r3 与 197-r3 之前的两腿就是死在把单跑当终态上，本腿照 this rule 写清。

**基线（HEAD＝`237e64f4`，未变异）**：`m0-baseline-cmdwisp rc=0 ok 5.780s`（新用例＋前一腿三枚载体判据一起跑）／
`m0-baseline-panel rc=0 ok 0.050s`（前一腿五枚名册判据）。

| # | 单点改动（落点） | 结果 | 逐字读数 |
|---|---|---|---|
| **M1** 拔掉 join 的读者（泵不再把卡交出去） | `internal/panel/pump.go:263` `TaskRosterSectionFrom(p.src.Tasks(), cards)` → `(..., nil)` | `cmd/wisp` **红 1 枚**（`rc=1`）、`internal/panel` 五枚名册判据里 **红 1 枚**（`TestTheRosterReaderPutsSubagentsOnTheWire`）⇒ 两头发牙，但**只有生产路径那一发看得见"卡就在同一枚包里"** | `--- FAIL: TestRunPacketMarksTheRosterRowACardIsHolding (2.11s)`；`:171: the packet carries a card for "bb825e0c-…" and its roster row still reads blockedOnApproval=false: 「被阻塞的子代理在所有界面上都不可见」 is the incident ticket 197 AC#4 names…`；同一发 `:168` 打出的字节里**那一行是 `"blockedOnApproval":false`、`"inFlightSlots":1`**，而 `pending` 里就带着那枚卡 |
| **M2** 字段恒假（写成常量 false） | `internal/panel/subagent_roster_197.go:210` → `BlockedOnApproval: false,` | `cmd/wisp` **红 1 枚**、`internal/panel` **红 1 枚**；名册其余字段照旧全绿 ⇒ 红在**这一枚字段**上，不是把整节带崩 | `--- FAIL: … (2.07s)`；`:171` 同一句红；`:234` 那行照旧打全了两半读数（`pending=1 … after the card: blocked=false`）——**卡活着那一半红、卡没了那一半绿**＝恒假的签名 |
| **M3** 字段恒真（写成常量 true，＝票 181 AC#7 的"空转"形） | 同一枚行 → `BlockedOnApproval: true,` | `cmd/wisp` **红 1 枚**，红在**另外两枚断言**上（`:191` 与 `:224`），**不是** `:171` ⇒ "挂卡那一发绿"对恒真毫无抵抗力，是本用例**第二半与根行那一半**在替它承重 | `:191: the root row reads blocked although no card names it (the only card is "14fa6c1e-…"): the join is by correlation id, not by anything being in flight` ＋ `:224: the card is gone and the row still reads blocked: the field would then be a constant…`；`:222` 的字节里两枚行都是 `"blockedOnApproval":true`、`"inFlightSlots":0` |
| **M4** 装配根拔掉整枚载体 | `cmd/wisp/run.go:491` `Tasks: rt.taskRosterState,` → `Tasks: nil,` | `cmd/wisp` **红 1 枚**（本用例）；`internal/panel` 那五枚**全绿**（前一腿 M1 同结论：**包内判据对装配回归是瞎的**） | `:168: the packet carries no tasks key at all: {"pending":[{"correlationId":"53cc9506-…","tool":"ticket197.blocked.probe",…,"level":"L2",…}],"results":[…],…,"instructions":{…}}`——`grep` 那串字节里**顶层四枚＋instructions、没有 `tasks`**，而 `pending` 里那枚卡的 `correlationId` **就是孩子的 task id** |
| **M5** 卡还在队列、但**不指这一行**（测试侧单点） | `cmd/wisp/subagent_blocked_197_test.go:121` `CorrelationID: id` → `id + "-elsewhere"` | `cmd/wisp` **红 1 枚**，红在**前置读数**上（轮询 600×10ms 耗尽，8.33s） | `:165: no packet the run published carries the card on its queue (pending [{CorrelationID:a5da9b13-…-elsewhere Tool:ticket197.blocked.probe … Level:L2 …}], wanted one card named "a5da9b13-…"); the run displayed 1 card(s)` |

每发还原后的复跑：`m1-restored-cmdwisp rc=0 ok 1.992s`／`m2-restored rc=0 ok 2.130s`／
`m3-restored rc=0 ok 1.763s`／`m4-restored rc=0 ok 2.437s`／`m5-restored rc=0 ok 1.940s`，
五枚 `RESTORED_CLEAN`，脚本末 `git status --porcelain -- cmd internal` **输出为空**。
两发包内红的逐字同一句（M1 与 M2 各一发）：
`subagent_roster_197_test.go:165: the child has an L2 card on the queue and the packet does not say so:
「被阻塞的子代理在所有界面上都不可见」 is somebody else's recorded incident`（`internal/panel`，`rc=1`）。

**M5 到底证明了什么、没证明什么（不替它吹）**：它证明本用例**不是**"队列里有任何一枚卡就算数"——
前置读数要求那枚卡的 `correlationId` 就是这枚行的 task id。它**没有**证明"产码若写成`有卡即真`会被 M5 抓到"：
那一形在 M5 里会先撞前置红，抓不到 join 那一步。**能抓`有卡即真`的是 M3**（`:191` 那枚根行断言），
所以 M5 记作"前置不空转"的那一发，不记作 join 键的那一发。

**"哪一枚 commit 算未修码"**（派单点名要写清的）：
- 判据**不存在**的那棵树＝**`da1d9289`**（本腿起手 HEAD，也是 197-r3b 交完之后编排者入库第三轮对标件的那一枚）。
  那一棵树上**产码是全的**（`blockedOnApproval` 的 join、`tasks` 节、名册读者都在），
  **缺的只有 cmd/wisp 这一发** ⇒ `da1d9289` 算未修码的是**判据面**、不是产码面。
  本腿**没有**用 checkout/reset 去把那棵树打出来（共享工作树禁），而是用 **M1/M4 在同一棵树上复现**：
  M1 打出"join 的读者拔掉了"的字节（`pending` 带着卡、行读 false），M4 打出"整节没上线"的字节（连 `tasks` 都没有）。
- 判据**入库**的第一枚 commit＝**`237e64f4`**（本腿，1 枚文件、**251 插入／0 删除**、**产码 0 枚**）。
  它的父＝`da1d9289`。⇒ 若将来谁把 `subagent_blocked_197_test.go` 改回"只断言 false"，
  未修码那一枚就是**改回去的那一枚**，而 M1/M2/M3 三发红在 `237e64f4` 之后的树上仍可复跑（台件已入库）。
- ⚠ 本腿**没有**新增 commit 去改产码来"让正控更方便"：五发变异全部只存在于台件运行的一瞬，
  盘上留下的只有读数文件（`git status --porcelain -- cmd internal` 空、`gofumpt -l internal cmd` 空）。

---

## ④ 门禁逐包整包点数 ＋ panel 红名册逐名比对 ＋ 格式门

| 门禁 | 起手（`da1d9289`，本腿亲手跑） | 终态（同一内容入库为 `237e64f4` 之后复跑） |
|---|---|---|
| `go test -count=1 -v ./internal/panel/` | `rc=1` **PASS=95 FAIL=4 SKIP=0**，`FAIL 1.291s` | `rc=1` **PASS=95 FAIL=4 SKIP=0**，`FAIL 2.201s`（`final-v-panel.txt`）；变异还原后第三遍 `FAIL 1.453s`、枚数照旧（`final2-v-panel.txt`） |
| `go test -count=1 -v ./internal/tools/` | `rc=0` **PASS=158 FAIL=0 SKIP=0**，`ok 15.306s` | `rc=0` **PASS=158 FAIL=0 SKIP=0**，`ok 16.181s`＝**＋0 枚**（本腿没动 tools 包） |
| `go test -count=1 -v ./cmd/wisp/`（**整包，不是 `-run`**） | `rc=0` **PASS=96 FAIL=0 SKIP=0**，`ok 94.520s` | `rc=0` **PASS=97 FAIL=0 SKIP=0**，`ok github.com/CarlosShao/wisp/cmd/wisp 100.159s`＝**96 ＋ 本腿这 1 枚** |
| `go build ./...` | **rc=0**，输出空 | **rc=0** |
| `go vet ./cmd/wisp/ ./internal/panel/ ./internal/tools/` | **rc=0**，输出空 | **rc=0** |
| `sh scripts/d22scan.sh` | — | **rc=0 clean**：`no D22 ban violations`；ban#1-5 `internal/`=216、`cmd/`=24；ban#7 `internal/tools/`=22；ban#8 `internal/`=**453**、`cmd/`=**51**（＝197-r3b 那发的 50 ＋1，多的就是本腿这枚新测试件）、`design/`=39、`frontend/`=85。逐字＝`d22scan-final.txt` |
| gofumpt（`"$GOPATH/bin/gofumpt.exe"`） | `-l internal cmd` **输出为空＝0 枚未净**（新件写完即净，本腿没有跑过 `-w`） | **仍 0 枚**（`final-gate.txt`／`final2-gate.txt`） |
| `git diff --numstat`（跟踪集自证） | — | 本腿唯一入库件＝**251 插入／0 删除**（`git show --numstat 237e64f4`）；**删除列＝0** |

**终态那一列跑了三遍，本腿把每一遍都留下**（因为 §③ 的五发变异动过 `internal/panel/pump.go`、
`internal/panel/subagent_roster_197.go`、`cmd/wisp/run.go` 与 `cmd/wisp/subagent_blocked_197_test.go` 四枚件，
"还原干净"这句话必须是**跑出来的**、不能是 `git status` 一空就算完）：

| 遍 | 树 | 读数 | 文件 |
|---|---|---|---|
| 第一遍 `final-*` | `da1d9289` 的内容＋本腿新件（尚未 commit） | panel `PASS=95 FAIL=4`／tools `PASS=158 FAIL=0`（`ok 16.181s`）／cmd/wisp **`PASS=97 FAIL=0`**（`ok 100.159s`） | `final-gate.txt`、`final-v-{panel,tools,cmdwisp}.txt` |
| 第二遍 `final2-*` | **commit 之后、五发变异还原之后**（`git rev-parse HEAD`＝`237e64f4`，`git status --porcelain -- cmd internal` 为空） | panel `PASS=95 FAIL=4`（`FAIL 1.453s`，逐名同下表）／cmd/wisp **`PASS=97 FAIL=0`**（`ok 98.216s`）；build／vet／gofumpt 三枚 rc=0 且输出空 | `final2-gate.txt`、`final2-v-{panel,cmdwisp}.txt` |
| §③ 台件里的 `m*-restored` | 每一发变异还原后立刻单跑该判据 | 五枚全 `rc=0`（逐字秒数见 §③ 末） | `m1-…-restored-cmdwisp.txt` 等 |

⇒ 本腿写"绿"用的是**整包 `-v` 现数**（97＝96＋本腿 1 枚），`-run` 单跑只进 §③ 的变异表、不当终态。
`internal/tools` 第二遍**没重跑**——本腿一枚 tools 文件都没碰；若要三包都按终态数，那是验收程的活，
本腿在此具名而不替它报数（沿用 197-r3b §④.2 的写法）。

**在册红名册逐名比对（`internal/panel`，起手 4 枚 vs 终态 4 枚）**：| 派单给的在册红 | 终态现量（`grep '^--- FAIL' final-v-panel.txt`） | 结论 |
|---|---|---|
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `--- FAIL` 在 | **名册未变**（本腿一未修、二未 Skip、三未放宽） |
| `TestComposerContractTypesMatchFrontend` | `--- FAIL` 在 | **名册未变** |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | `--- FAIL` 在 | **名册未变** |
| `TestC21DesignTokensFourwayAgree` | `--- FAIL` 在 | **名册未变** |

**没有第 5 枚、也没有任何一枚消失**；`cmd/wisp` 与 `internal/tools` 终态各 **0** 枚红。
两枚名册钉的**红因本腿一字未加**（红因仍是前一腿 `tasks` 键那一跳，逐字照旧）：

```
approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
composer_test.go:74:  Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
```

⇒ 界面那支要补的键集**没有因为本腿变长**：`blockedOnApproval` 早在 `TaskRowView` 里（197-r3b §④.3 已逐字段交给界面），
本腿**没有新增任何 JSON 键**，只是把那枚已有键的**判据**补上。两枚"四枚顶层键"的钉终态逐名照绿：
`--- PASS: TestThePumpBuildsThePacketFromWhatTheHostHolds`（钉在 `pump_test.go:123`）、
`--- PASS: TestPublishHandsTheBytesToTheAttachedExit`（钉在 `pump_test.go:291`）、
`--- PASS: TestAPumpWithoutARosterReaderSendsFourKeys`。
`--- PASS: TestSubagentStreamKeyHasOneMintSite` 亦照旧（本腿没碰字面量）。

---

## ⑤ 我没测什么（具名，不留空）

1. **"子代理因为自己的调用被卡住"这一形今天仍然造不出来**——前一腿那半判断本腿现读**认同**：
   mockllm 不回 tool call、`internal/agent` 不发 `tool_choice`，孩子的请求到不了门。
   所以本腿那枚卡是**宿主代签**的（`AdmitTextTask`＋`PendingApproval`，与 `run.go:572 confirmModeSwitch` 同形），
   它钉住的是**join 这一跳在真包上成立**，不是"孩子自己会撞门"。
   ⇒ 要真事件那一形，得先有能让 provider 回 tool call 的仪器（`C5` 的 golden SSE 那支），**那是新范围，本腿没扩**。
2. **L1 确认窗口不会被这枚字段报出来**。现读：`internal/agent/approval/gate.go:218 PendingWindow` 走
   `g.openWindow(w)`／`windows` 那张表，**不进 `q.pending`**；而 `blockedOnApproval` 的输入是
   `pending_read.go:104 LiveApprovals()`，它只遍历 `q.pending`。⇒ **一枚正在 L1 三秒窗口里等孩子的父调用，
   名册行读 blockedOnApproval=false**。这是 AC#4 语汇上的一处**真实缺口**（"被授权卡住"在 D43／B1 里有两形，
   载体只报了 L2 那一形），本腿**具名交回，不改产码**——把它报出来要不要动 `LiveApprovals` 或加一节，
   是契约级判断，归编排者／owner。
3. **多枚卡**：队列里同时两枚卡、其中一枚指这名册行——没测（`DefaultMaxPending`=8，形同 `panel` 那发的单卡）。
   一枚卡**指向根行**时根行必须 blocked——没测（本腿只测了反向：卡指孩子、根行必须不 blocked）。
4. **前置读数没被 M5 覆盖的那一半**（见 §③ 末）：产码若写成"队列里有任意一枚卡 ⇒ 每行都 blocked"，
   抓到它的是 `:190` 那枚根行断言（M3 那一发已经证明它会红），**不是** M5。本腿没有再为那一形单独造一发。
5. **本腿没有钉"未登记的 task 拿不到卡"这一发**。孩子的 id 本腿自己 `AdmitTextTask` 过，
   所以"D47 会拒一枚没登记的卡"这一形在这枚用例里**不可见**（它在 `cmd/wisp/run_test.go:340` 那发
   `TestComposedGateBlocksAWriteForTwoSeconds` 的第 (0) 段钉着，本腿没有重写、也没有绕过）。
6. **"每一枚包都如此"没测**：本腿断言的是"挂卡期间**存在**一枚 run 自己发布的包，它带那枚卡且那枚行 blocked"。
   收口之后的包当然不带卡（第二半读的就是这一形），所以这一条是**存在量词**、不是全称量词。
7. **孩子活着时那枚 `AdmitTextTask` 的副作用没测**：本腿为孩子的 id 第二次登记，`gate.go:160` 的
   `admitted` 是 map、`admitOrd` 追加两次；两枚 revoke（孩子的环路那一枚与本腿 `defer` 那一枚）
   谁先谁后会短暂删掉那枚 key。本腿**没有断言过这一段的任何语义**，也没测到它对名册的影响；
   之所以这样写：名册行与状态不读 `admitted`（读的是 `tools.TaskRoster`），所以 join 不受它影响——
   **但这句推理本腿没有用判据钉住**，只作具名。
8. **点击／渲染那一跳仍然没测**（前一腿 §⑤4 同一条）：所有读数停在**包字节**，传输归票 33。
   没有任何一枚判据证明"页面渲染过那枚 blocked 的行"。
9. **名册的字节规模上限没量**（前一腿 §⑤5 同一条，本腿没动）。本枚包里 `bytes=3091`（两枚行、一枚卡、
   一条 `results`）——多带几枚卡／几枚行会不会撞上 `observe.MaxLoggedString`／`summaryClamp`=440 那两枚数，
   本腿只拿到了"这一发 3091 字节、账本那行的 sha256 记上了"这一枚读数，**没有量边界**。
10. **根行的状态维今天仍没人填**（裁定 `A394` 留给票 196）：本腿的包字节里根行照旧
    `"status":"","statusKnown":false,"statusReason":"宿主没有登记这一维…"`，
    前一腿那枚"空必须带原因"的断言仍在（`subagent_carrier_197_test.go:323`），本腿没有替它改判。

---

## ⑥ 对派单的不服（逐条；能给读数就给读数）

1. **"预计 20-30 行测试"不成立**，实际 251 行（`grep -c ""`＝251，非注释非空 134 行；
   逐段可核：header 注释 `:3-31`、两枚常量 `:44-56`、用例本体 `:59-237`、helper `:242-251`）。
   差别不是本腿罗嗦，是派单那 20-30 行只装了"挂卡 ⇒ true"这一枚断言，而本腿在同一枚用例里还必须装下：
   撤卡 ⇒ false（否则字段是常量）、根行不翻（否则是有卡即真）、卡回 reject（否则载体成了批准出口，`AC#5`／`Q-49 丙`）、
   `statusKnown`＋`statemachine.Valid`＋`streamKey` 三枚"这枚 join 没把别的字段写坏"、sha256 记账那一枚。
   逐枚都有名字：**M2 只被第一枚抓、M3 只被第二/第三枚抓**（§③）——少了哪一枚，本腿的正控就有一发跑不红。
   本腿**没有**加第二枚用例、没有改前一腿任何 helper、没有动产码一枚。
2. **派单给的行号本腿现读为**准**，没有漂**：`internal/panel/composer.go:91`＝`Tasks *TaskRosterSection json:"tasks,omitempty"`；
   `cmd/wisp/run.go:577`＝`revoke := rt.gate.AdmitTextTask(taskID)`、`:579`＝`ans, why := rt.gate.PendingApproval(ctx, tools.Decision{`
   （两枚都在 `confirmModeSwitch` 里，函数从 `:572` 起）。派单那句"行号自己现读，我给的会漂"本腿照办了，
   并顺手把 197-r3b 引用的 `run.go:421-427`（`composer.go:52` 注释里那串装配根行号）现读为**已漂**：
   今天装配根的 `PumpSources` 字面量在 **`run.go:472-493`**。那串过期行号在**产码注释**里
   （`internal/panel/composer.go:52`），本腿**没有改它**（派单写死"产码不许动"），在此具名交回。
3. **"前一腿说重了"这个判断：本腿认同，但 197-r3b §⑤1 自己那半也说轻了**。它写"难点只剩孩子的 id 要在它活着时读到，
   `taskRosterState()` 每 5ms 轮询即可"。现读：孩子的 id 由环路自己铸——铸点在
   `internal/agent/loop.go:322`（`RunAsync` 的第一枚语句 `id := newTaskID()`，函数本体 `loop.go:1123`），
   宿主拿不到回调，只能轮询——这半它对；**但它没提**名册行比 gate 登记**早一瞬**
   （`internal/tools/subagent_197.go:293-301`：`opt.AdmitTask` 先 `PublishSubagent`、再 `baseAdmit`），
   所以**只依赖孩子自己那枚登记**的第一发 `PendingApproval` 有可能撞上 D47 未登记被拒（`gate.go:453` 直接 return reject、
   队列里什么都没有 ⇒ 用例空转红）。本腿因此按 `confirmModeSwitch` 的**原形**自己 `AdmitTextTask(id)`
   ——那不是绕，那**就是宿主现成的入口**（`run.go:576-577` 一枚合成 id 也是这么登记的）。
   代价具名：见 §⑤7，本腿因此没有钉"未登记会被拒"那一发。
4. **当年裁定不需要重开**：`Q-49 丙`／票 197 `AC#5`（子代理永不自批）在本腿之后仍成立，且有读数——
   撤卡那一发自己回 `answer=reject`／`why="任务上下文已结束，审批请求已作废并按拒绝处理"`（`gate.go:529`→`queue.go:427`）。
   本腿**没有**调用 `Queue().allow/reject`、没有碰 `DecideFromNative/DecideFromPanel`、
   没有新增 C17 方法名（`grep -rn "PanelBridge" internal/ cmd/` 非测试码仍是 197-r3b §⑥4 记的那两枚注释命中）。
   ⇒ **没有新格待人拍板**从本腿产生；唯一**本腿认为该摆出来的新问题**是 §⑤2（L1 窗口不被 `blockedOnApproval` 报出），
   那是**载体语汇的射程**问题，不是判据缺失，本腿只具名、不动产码。
5. **派单那句"internal/panel 红＝在册那 4 枚"本腿复跑＝真**（起手 4 枚、终态 4 枚、逐名同一批，§④ 表）。
   派单给的另两枚终态数（tools 158／cmd/wisp 96）本腿起手复跑**逐枚对上**，
   终态 cmd/wisp＝**97**＝96＋本腿这 1 枚，是本腿唯一带来的枚数变化。
6. **"产码零改动"这条本腿守到了，并且没有"非改产码不可"的停手项**：五枚变异只活在台件的一瞬
   （`git status --porcelain -- cmd internal` 终态空、`git show --numstat 237e64f4`＝251/0）。
   唯一一处**本腿想改而不能改**的是 §⑥-2 那串过期行号注释（`internal/panel/composer.go:52` 写 `run.go:421-427`），
   本腿按边界只具名不落地。
7. **本腿自己的一枚操作账**：修 §② 里三枚行号引用时本腿用了 `python - <<'EOF'` 改 md（派单写"中文长段用
   Edit/Write 工具，别用 shell heredoc"）。替换是逐条精确串、跑前断言"找不到就中止"，盘上结果已核；
   但**这是越界动作**，记在这里，不解释成"反正没坏"。此后所有中文段落一律走 Edit/Write。

---

## ⑦ 入库件清单（本腿，逐枚点名）

| commit | 内容 |
|---|---|
| **`237e64f4`**（判据） | 1 枚文件：`cmd/wisp/subagent_blocked_197_test.go`（251 插入／**0 删除**；**产码 0 枚**） |
| 第二枚＝本文件＋台件 | `docs/evidence/s1/197-blocked-carrier-r4.md` ＋ `.scratch/wisp/probes/197/r4/**`（起手与终态整包 `-v` 读数、`start-gate.txt`／`final-gate.txt`／`final2-gate.txt`、`mutate-197r4.sh` 台件、`m0`–`m5` 逐发输出与 `mutation-log.txt`、`d22scan-final.txt`） |

**只 commit、未 push**；写面之外零触碰——`docs/evidence/s1/152-*.md`、`.gitignore`、`design/**`、
`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/**`、`.scratch/wisp/.scratch/**`、
仓根那 8 枚 `part*-fixed/pristine.txt` 一枚都没进本腿的 commit（`git show --stat` 逐枚点名即可复核）。
`frontend/**`、`design/**` 本腿**未读未写**；`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、
`allowlist.txt`、`internal/panel/tokens_fourway_test.go` **一字未动**。

