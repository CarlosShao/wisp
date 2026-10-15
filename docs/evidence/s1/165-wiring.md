# 票 165-r2 · 装配半格（把 plan 家族接进生产装配）—— 落地腿自证

- 腿：165-r2（实现者；判语归非实现者，本件不翻票面任何 AC 框）
- 起手：2026-10-15 22:08 +08，HEAD `818ff733`（派单存档）；写面期间 HEAD 漂到 `6f53583d`（别腿提交，本腿零碰）
- 仓库：`D:\work\workspace\projects plans\Wisp`，分支 `dev`，Shell＝Git Bash
- 派单：`.scratch/wisp/probes/orch/20261015/txt/165r2-dispatch.txt` ＋ 释放记录 `165r2-send.txt`（冲突以 release 为准）
- 读数件：`.scratch/wisp/probes/165/r2/00..11`（每把尺一枚，件里都带 `rc=` 那一行）
- 本件所有"读数"均为本腿现跑（量来的），凡"算来的"都就地标注。

## 0. 这一单为什么存在（"写了没接线"那一族）

- 尺＝`grep -rn "BuiltinPlanEntries" --include=*.go internal cmd | grep -v _test`（射程＝internal＋cmd 的产码 .go）
- 读数（件 00）：命中 2 行，全在定义面自己（`internal/tools/plan_present.go:134` 注释、`:138` func）⇒ **产码调用者 0 枚**；
  测试侧另有 2 枚（`plan_present_165_test.go:76` 测试自己 append、`:286` 取 `[0]`）⇒ 与编排者 21:2x 的量逐字一致。
- 结论：`go build`/`go vet` 的 rc=0 证不了接线，本腿的交付物是"那一行＋能红的钉"。

## 1. AC#1 未修码读数（接线之前量，⛔ 写成"应该做不到"）

判据：在**未接线**的树上，"只读地出一份方案"这一发到不到得了。三把尺同发：

1. 调用点枚数尺（件 02，`go test ./cmd/wisp/ -run 'Test165r2PlanWiringCallSiteCensus' -v -count=1`，rc=1）
   阳性对照先绿：`cmd/wisp/run.go:783 range=true register=true`（已接先例 subagent 那枚）⇒ 尺不瞎；
   plan 那枚：**调用点 0 枚，命中清单＝[]**。
2. 装配 registry 尺（件 03，`go test ./cmd/wisp/ -run 'Test165r2' -v -count=1`，rc=1）
   真 `assembleRuntime` 出来的 registry 名册 **9 枚**＝`[fs.edit fs.list fs.move fs.read fs.trash fs.write
   task.cancel task.output task.spawn]`，零枚 plan 形。
3. 把这一发真打一遍（同一发，c) 用例的判定行）
   经产码 bridge 执行 `plan.present` 的回文＝**"未知工具 plan.present，可用工具见 list_tools"**，且**没升起任何卡**。
- 成立：未修码这一发**到不了**，本票前提由量撑住（不是推断）。

## 2. 接线（产码一处，形状照先例）

- 落点：`cmd/wisp/run.go`，subagent 那块 for-range 之后，同一 registry、同一错误回口：
  `for _, e := range tools.BuiltinPlanEntries() { if err := reg.Register(e); ... return rt, 2 } }`
- 尺＝`git diff --numstat -- cmd/wisp/run.go` ⇒ **29 / 0**（新增 29、删除 0；"29 行＝1 块 for＋26 行注释"是本腿**数注释行算来的**，产码行枚数见件 09 逐字列出）。
- ⛔ 未改 `internal/tools/**` 一字：件 07 记母仓 `internal/tools/plan_present.go` 的 sha256 `3fc87895...`
  与仓外快照树基线**逐位相等**；`git status --porcelain -- internal/tools` 本腿期间只出现过别腿那枚 staged 文件（后由 `287-r1` 自己提交）。
- ⛔ 另造一套装配：调用点枚数尺要求那一处**必须是** `for range 该调用 ＋ 循环体调 Register` 那一形（件 02／04 的 `range=true register=true`）。

## 3. 接线钉（本票最要紧的一格）＋ 正控两发

钉在两枚新件里（都在 `cmd/wisp/`，都是 `_test.go`）：

- a) `plan_wiring_165r2_test.go` = 调用点枚数尺（go/ast，射程＝internal＋cmd 全部 .go 剥 `_test.go` 与 `/testdata/`）
- b) `plan_assembled_165r2_test.go` = 装配可达尺（真装配 registry／真 `bridge.Tools()` 模型面／真确认卡／真 `tool_call` 行／真会话授权表）
- c) 同件里 = 批准只批这一批（装配版，见 §5）

**接线后绿（件 04）**：`go test ./cmd/wisp/ -run 'Test165r2' -v -count=1` → rc=0，三枚全 PASS，读数逐字：

- `装配可达读数：registry 10 枚含 plan.present；卡=plan.present；行 risk=L2 decision=allow；授权行 0 枚`
- `AC#3① 装配读数：方案卡=plan.present（yes）→ 清单外卡=fs.write（no，IsError=true）→ 授权行 0 枚 → 目标文件未落盘`
- `接线读数：cmd/wisp/run.go:822 range=true register=true ⇒ 产码调用点 1 枚，形状与先例同形`

**正控第一发＝摘掉那一行必须红（件 05）**：物理摘除（仓外快照 `/d/tmp/165r2-snap-a/run.go.wired`，
sha256 `7f28db4f...`；⛔ `-overlay`）后同发三枚全红，rc=1，红句首行逐字：

1. `接线 RED：生产装配的 registry 里没有 plan.present；名册 9 枚：[fs.edit fs.list fs.move fs.read fs.trash fs.write task.cancel task.output task.spawn]（这正是开工前 AC#1 的未修码形状：这份清单在 wisp run 里够不到）`
2. `接线 RED：批方案那一发没升起卡（回文="未知工具 plan.present，可用工具见 list_tools"）⇒ 装配里没有这一发，c) 的后续读不到`
3. `接线 RED：tools.BuiltinPlanEntries 在 internal/ 与 cmd/ 的产码里调用点 0 枚，要 1 枚（cmd/wisp/run.go 的装配那一处）；命中清单＝[]（0 枚＝写了没接线那一族，plan.present 在生产装配里根本够不到）`

同一发里阳性对照仍绿（先例 subagent 仍 1 枚）⇒ 红是那一块造成的，不是尺坏了。

**正控第二发＝还原后绿（件 06）**，还原证明三把尺都本腿自算：

- `sha256sum cmd/wisp/run.go` 与快照 `run.go.wired` ⇒ 两个 `7f28db4f...` 逐位相等；
- `git diff --numstat -- cmd/wisp/run.go` ⇒ 仍是 `29 0`（回到摘除前那一发）；
- `go test ./cmd/wisp/ -run 'Test165r2' -count=1` ⇒ rc=0、`ok ... 2.874s`。

**名册凭据的口径**：`plan_present_165_test.go:76` 那份名册是**测试自己 append** 的，⛔ 当"已接"凭据；
本钉的 b) 那把尺读的是 `assembleRuntime` 产码装配的 `bridge.Registry()`（件 03／04 两色对照就是证据）。

## 4. AC#2 复用形状（⛔ 新增状态）

- 尺＝`git status --porcelain -- <契约轴逐枚点名>`（件 08，射程＝docs/PLAN.md／docs/specs／internal/risk／
  internal/statemachine／internal/panel／internal/agent/approval／internal/agent/statemachine／frontend／
  ci.yml／.github／tools/d22scan/allowlist.txt／internal/llm/testdata／golden／thresholds.go）
  ⇒ 每档 **0 行**。
- 尺＝`grep -cE '^\s+[A-Za-z]+ +State = ' internal/statemachine/states.go` ⇒ **20 枚**（D43 分母，本腿没动那棵树）；
  `git diff` 里 statemachine 命中 **0 枚**、新增 State 常量 **0 枚**。
- 行为面复用尺（现跑，rc=0）：`Test165PlanPresentAsksThroughTheExistingL2Card` PASS——`Gate` 仍 2 个回调、
  走的是现成 L2 卡（D43 #17/#23/#24/#25 那一族）。
- 成立：零新增状态／零新卡型／零新答复通道。**没做完**的部分：无（这一格本腿射程内已量齐）。

## 5. AC#3 两发变异

**① 批完方案后来一发没在方案里的写 ⇒ 必须响**——两形都量：

- 产码／单测尺（rc=0，件 09）：`Test165PlanApprovalRingsAgainForAnUnlistedWrite` PASS（L1 窗口被问 1 次、授权表零行、decision 非 allow_session_grant）。
- **装配尺（新，本腿钉的 c)）**：同一次 `assembleRuntime` 里先 `plan.present`（列 1 步）答 `yes`，再 `fs.write` 到
  夹具里没被 `[fs]` 罩住的那一发 ⇒ 升起**第二张卡**（tool 名＝`fs.write`）、答 `no` 后未执行、文件未落盘、
  `approval_grant` 表 **0 行**。读数行见 §3 的绿句；audit 同发有
  `ANSWER-ALLOW corr=... tool=plan.present route=native decision=allow` 与 `ANSWER-REJECT corr=... tool=fs.write decision=reject`，
  且用例另钉一条 `ANSWER-EXPIRED` 不许出现（钟到点≠人点头）。

**② 摘掉"批过的范围"这一味 ⇒ 是否存在一发越界从此打不红**（件 07，台件＝仓外快照树 `/d/tmp/165r2-snap-b/repo`，
`git archive HEAD` 导出，导出树 `.go` 705 枚、sherpa dll 先数得 0 枚再 cp 3 枚＋build dll 3 枚＝A802 那条坑照做；
母仓 `internal/tools/**` 零字）：

- 基线同树：`internal/tools -run Test165` rc=0（5 枚全 PASS）；`cmd/wisp -run Test165r2` rc=0（3 枚全 PASS）。
- M-B1 摘掉整句（删 `planPresentText` 里那一行 `b.WriteString(planPresentScopeClause(...))`）⇒ rc=1 两包，
  红的只有两枚文字尺：`Test165ScopeClauseNamesTheBatch`（缺 "只限上面这 3 步" ＋ 缺 "清单之外"）与本腿
  `Test165r2AssembledBridgeReachesPlanPresent`（同一对 needles）。
- M-B2 换成**不提条数**的通用句（保留"只限这一批"与"清单之外"）⇒ rc=1 两包，红的仍是同一对文字尺，
  且**只缺一味**（缺 "只限上面这 N 步"，"清单之外"不红）⇒ 那把尺的两味是真分开的。
- **两发之下行为尺都还绿**（`Test165PlanApprovalRingsAgainForAnUnlistedWrite` 与本腿 c) 全 PASS）。
- 具名报回（⛔ 顺手补）：**不存在"一发越界从此打不红"**——越界那一维不由这段文字裁决，它由
  "L2 永不进持久授权"那条结构裁决，而那条有两形尺各自量得到；范围句只决定**人读卡上看没看到边界**，
  那一维红在两枚文字尺上。要补的是"条数之外的第二味"（如"清单之外"被摘而条数还在）这种半摘形状，
  本腿只报形状不补尺。

## 6. AC#4 面板不许代答

- 尺＝`git diff -- cmd/wisp/run.go` 的 `+` 行扫描（件 09）：含 allow/approve/grant 字面的新增行 **2 枚，全在注释散文**
  （"no allow exit"、"grantRead stays the one place"）；非注释产码行 **5 枚**＝`for/range`、`if err := reg.Register(e)`、
  `fmt.Fprintf`、`return rt, 2`、`}`——零枚答复出口。
- 现跑冻结边界钉四枚（`internal/panel`，rc=0）：`TestAnsweredPanelRoutesCarryNoApprovalDecision`／
  `TestRealGuardRefusesEveryAssemblableApprovalRouteName`／`TestNoInboundEnvelopeCanBindAnApprovalVerdict`／
  `TestGrantWireShapesAreRefusedAtTheDoor` 全 PASS；`internal/panel` 本腿零字节（件 08）。
- 成立：加的是"发起请求"的入口（L2 声明走现成 C18 卡），没加"给出允许"的出口。

## 7. AC#5 契约轴零字节

- 尺＝件 08 的逐路径 `git status --porcelain` ⇒ 名单里每一档 **0 行**；`internal/` ban #8 分母
  起手 **568** ＝ 终态 **568**（零移动，件 01b 与件 11 两发同尺），`cmd/` 136 → 138（＝本腿两枚新 `.go`，分母如期 +2）。
- `design/**` 31 行脏 ＝ **机主外部 agent 的地界**（AGENTS §1.3：`frontend/**`、`design/**` 本编队零写面），
  本腿在其中零字节；`go.mod`／`.gitignore` 的 M 同理不是本腿的。

## 8. 门禁各把尺（终态件 11，逐条带 rc）

| 尺 | 射程 | rc |
|---|---|---|
| `go build ./...` | 全模块 | 0 |
| `gofmt -l`（本腿三枚文件） | run.go＋两枚钉 | 0（名册空） |
| `gofumpt -l`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`，本腿三枚） | 同上 | 0（名册空） |
| `gofumpt -l cmd internal` | 两棵树 | 0，**列出 6 枚存量件**（`cmd/wisp/models.go`、`cmd/wisp/panel_transport_35r2_test.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、`internal/tools/task_output_pointer_notice_errtext_174r7_test.go`）＝起手件 01 同样 6 枚，**本腿零新增**、⛔ 顺手修 |
| `go vet -tags winlive ./cmd/wisp/` | 共享编译面 | 0 |
| `sh scripts/d22scan.sh` | 全仓七禁＋正控 | 0（`runtests.sh PASS=35 FAIL=0 SKIP=0`；live `clean - no D22 ban violations`，ban #8 internal/=568、cmd/=138） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 全仓跟踪名册 | 0（positive control PASSED；tracked=13428，over-budget=60＝roster=60，not in roster=0，VERDICT GREEN） |
| `go test ./cmd/wisp/ -run 'Test165r2' -v -count=1` | 定向三枚 | 0 |
| `go test ./internal/tools/ -run 'Test165...(五枚)' -count=1 -v` | 定向五枚 | 0 |
| `go test ./internal/panel/ -run '<四枚边界钉>' -count=1 -v` | 定向四枚 | 0 |
| `go test ./cmd/wisp/ -run 'TestTicket255|TestAC1AC2DispatchHopGate133|TestAC3SelfProof133|TestAC4Distinctness133' -count=1 -v` | 邻座两族 | 0（23 枚 PASS／0 FAIL／0 SKIP，含 `TestTicket255RosterStillMatchesTheActualReadSites` PASS 0.73s） |

## 9. 台面故障（另节具名，⛔ 计入判据）

- 一发红来自**台件**不是代码：`台面故障（不是判据）：扫描根 "internal" 读不到：GetFileAttributesEx internal:
  The system cannot find the file specified.`（`plan_wiring_165r2_test.go` 第一次跑）
  ⇒ 成因＝go test 的 CWD 是包目录（`cmd/wisp`），不是仓库根；处置＝钉里加 `repoRoot165r2 = "..",".."` 一跳并用
  `filepath.Rel` 归一名，改后同一把尺拿到 0 枚／1 枚两色（件 02／04）。红的是我自己那行 `t.Fatalf(扫描根…)`，
  不是断言。
- 一发红来自**台件**：`Test165r2UnlistedWrite...` 第一版 41.34s 后红，audit 逐字
  `approval: REPLY REFUSED corr="清单里没有这一步，不批" ... correlation_id 无对应待审批项` ＋
  `approval: ANSWER-EXPIRED corr=t165r2-batch-write ... decision=timeout->reject after=40s`
  ⇒ 成因＝本腿把理由打在了 corr **前面**（`splitReplyLine` 取第二段当 corr）；处置＝答复行改回
  `verb corr reason` 语序，并给用例加一条"不许出现 `ANSWER-EXPIRED`"（到点≠点头）。改后 1.28s 绿。
- `0xc0000135` 那形本腿**没撞上**，因为每发 `cmd/wisp`／`internal/tools` 都前置了
  母仓绝对路径 `third_party/sherpa-onnx`（3 枚 dll）与 `build`（3 枚 dll）；`GOTMPDIR` 先 `mkdir -p`。
  尺＝件 01/04/06/11 各发都同批写了 dll 枚数与 `rc=`。

## 10. 没做完／欠账（具名，⛔ 写成"应该没问题"）

1. **"有没有人真的按门铃"没做完**：缺的料＝模型真会调 `plan.present` 的读数（golden SSE 或真机一轮对话里出现
   `tool=plan.present`）。本腿三把尺只证"够得到装配"，票 161 那一格仍零读数，归口不在本票。
2. **面板那一屏没做完**：缺的料＝`frontend/**` 的渲染面（机主外部 agent 写面），本腿零字。
3. **整包 `go test ./cmd/wisp`（不带 `-run`）没跑**：缺的料＝车道（票 314 整进程死亡区＋`cmd/wisp` 同一时刻只许
   一枚腿的纪律）⇒ 请编排者在自己车道上按 `scripts/wisp-cli-tests.sh` 那形复跑取四数。
4. **`-tags winlive` 那族没跑**（255 真窗两枚在 `*_winlive_test.go`）：本腿按定向口径只跑了 tag 外名册，
   winlive 真机读数仍只有本机、且此刻台面有别腿。
5. **"方案条数上限"没发明数字**：165-r1 票尾具名欠的契约，本腿 c) 只用 1 步与 2 步各一发，没新造阈值。

## 11. 件数与写面对照（派单 §3 白名单）

| 件 | 路径 | 类型 |
|---|---|---|
| 1 | `cmd/wisp/run.go` | 装配那一处（29/0） |
| 2 | `cmd/wisp/plan_wiring_165r2_test.go` | 钉 a（181 行） |
| 3 | `cmd/wisp/plan_assembled_165r2_test.go` | 钉 b＋c（262 行，超一件预算＝拆两枚的唯一原因，尺不变） |
| 4 | `docs/evidence/s1/165-wiring.md` | 本件 |
| 5 | `.scratch/wisp/issues/165-...md` | 票尾 1 条 bullet（⛔ 翻 AC 框／改 Status／动台账） |
| — | `.scratch/wisp/probes/165/r2/*.txt` | 读数件（派单要求的"每把尺一枚"，00／01／01b／02／03／04／05／06／07／08／09／11） |
| — | `/d/tmp/165r2-snap-a/`、`/d/tmp/165r2-snap-b/repo/`、`/d/tmp/165r2-*.log` | 两枚突变快照树与原始日志，只建不删 |
