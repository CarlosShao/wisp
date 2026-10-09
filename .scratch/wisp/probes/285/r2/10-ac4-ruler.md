# 285-r2 · AC#4 第五形补尺读数（per-call corr 的"每枚调用各不相同"）

腿：`285-r2`（实现腿，不翻勾、不改判据句、不 push）
起手锚：`ce18b3d6cc3b6cc03066d65f108053ae68fcc1d7`（`2026-10-09 11:04 +0800`，`dev`）
起手第一笔 commit：`28258cd2`（`.scratch/wisp/probes/285/r2/00-anchor.md`，先锚后跑）
交件 commit（两件新测试）：`9a00a890`

## 0. 基线（`-count=1`，⛔ 无缓存重放）

| 时刻 | `internal/agent` | `internal/tools` |
|---|---|---|
| 起手（锚 `ce18b3d6`，未落我的件） | **82 PASS / 1 FAIL / 0 SKIP**，`rc=1` | **208 PASS / 0 FAIL / 0 SKIP**，`rc=0` |
| 收工（同一命令，同两枚包） | **84 PASS / 1 FAIL / 0 SKIP**，`rc=1` | **209 PASS / 0 FAIL / 0 SKIP**，`rc=0` |

起手那枚 FAIL ＝ **`internal/agent/loop_golden_test.go:70`（`TestGoldenSingleToolCall`）**，
`git status --porcelain -- internal cmd` 起手为空 ⇒ 它是**锚上就存在的红**，不是本腿带来的，
本腿也未碰它（见 §5 顶回）。收工的红句枚数与名目与起手逐字同一枚。

```
--- FAIL: TestGoldenSingleToolCall (0.00s)
    loop_golden_test.go:70: call identity = {TaskID:fffd5eef-... CorrelationID:fffd5eef-...#call_e1 CallID:call_e1 Name:echo ...}, want task fffd5eef-d85c-40b7-ae57-e81829aefd0b
```

## 1. 交付形状（两枚包都给，落点理由具名）

新建（⛔ 零产码字节改动、⛔ 未改任何既有测试件）：

- 甲：`internal/agent/ticket285_corr_distinct_rulers_test.go`（144 行／两枚用例）
  - `Test285CallCorrIsDistinctPerCall` —— 包内直接对 `callCorr` 的**两支各自作差**
  - `Test285LoopDispatchesOneCorrPerCall` —— 真 loop 派发（`parallel-tools` 夹具，一枚 turn 两枚调用），
    请求侧两枚 corr 作差 ＋ corr 后缀必须等于它自己那一枚 `CallID`
  - 理由：`callCorr` 是**小写未导出**（`internal/agent/loop.go` 定义），甲形要落就得落 `package agent`；
    并且**回落支（`callID == ""` ⇒ 按轮内位置成键）今天没有任何夹具能走到**——
    现量：`grep -rn '#call-' internal/ --include=*_test.go` ＝ **空输出**（全仓 SSE 夹具都带 call id）。
    所以回落支的形状尺只能是包内直接调用，这一发是甲不可替代的部分（读数见 `20-mutations.md` MU-2）。
- 乙：`internal/tools/ticket285_corr_rows_rulers_test.go`（113 行／一枚用例）
  - `Test285RosterRowsCarryDistinctCorrPerCall` —— 续票 283 那条真链
    （真 `agent.Loop` → 真 bridge → 真工具 → 真 `tool_call` 名册行），
    断言面落在 `ListToolCallsByTask` 读回的**两枚持久行**上。
  - 理由：票面判据写的是 `rows[0].CorrelationID != rows[1].CorrelationID`，
    行名册只有经真链才拿得到（`memory.ToolCall` 不带 `CallID` 列，行的自证只能靠 `Tool` 名与 corr 的配对）；
    复用票 283 的现成夹具 `build221`／`windowGate221`／`parent283Provider`／`fake197Provider`（⛔ 未新造捕获面、
    ⛔ 未改该件一字）。

## 2. 断言句原文（`文件:行`，收工 HEAD）

甲 `internal/agent/ticket285_corr_distinct_rulers_test.go`：

```
69:			if b.corr == "" {
72:			if b.corr == task {
75:			if !strings.HasPrefix(b.corr, task) {
78:			if prev, dup := seen[b.corr]; dup {
84:		if len(seen) != len(branch) {
88:		if branch[0].corr == branch[1].corr {
116:	if first.TaskID != res.TaskID || second.TaskID != res.TaskID {
119:	if first.CorrelationID == res.TaskID || second.CorrelationID == res.TaskID {
123:	if first.CorrelationID == second.CorrelationID {
128:		if req.CorrelationID == "" {
132:		if !strings.HasPrefix(req.CorrelationID, res.TaskID) {
136:		if suffix := strings.TrimPrefix(req.CorrelationID, res.TaskID+"#"); suffix != req.CallID {
```

（红句原文里出现的是各 `t.Errorf`/`t.Fatalf` 那一行，比上面的 `if` 晚一到两行，两者都在证据里逐字抄出。）

乙 `internal/tools/ticket285_corr_rows_rulers_test.go`：

```
69:		if r.CorrelationID == "" {
72:		if r.CorrelationID == res.TaskID {
75:		if !strings.HasPrefix(r.CorrelationID, res.TaskID) {
81:	if rows[0].CorrelationID == rows[1].CorrelationID {
89:		if prev, dup := keyed[r.CorrelationID]; dup {
94:	if len(keyed) != 2 {
106:			if !strings.Contains(corr, want) {
```

判据要求的五个面逐条对上：
两枚 corr 互不相等＝甲 `:88`／甲 `:123`／乙 `:81`；两枚非空＝甲 `:69`＋`:128`／乙 `:69`；
两枚都以 taskID 为前缀＝甲 `:75`／甲 `:132`／乙 `:75`；
各自可路由（不是同一个 key）＝甲 `:78`＋`:84`（路由表里一枚调用一格）／乙 `:89`＋`:94`；
各自路由回自己那一枚调用＝甲 `:136`（后缀＝自己的 `CallID`）／乙 `:106`（spawn 的键含 spawn、cancel 的键含 cancel）。

鉴别器（防恒真）：⛔ 不是"非空就算过"——两枚作差那一发在甲的**两支各自**与乙的行对上都是独立断言；
`seen`/`keyed` 两张表的"格数＝调用数"是第二个独立形状；
把作差那一发摘掉，MU-1 下甲的 `:79`/`:85` 与乙的 `:90`/`:95`/`:107` 仍会响（三处冗余面，见 `20-mutations.md`）。

## 3. 正控（⛔ 不是"只报改前绿"）

- **夹具前置逐字适用**：`TaskID != CorrelationID` 在甲乙两侧都是**断言**而非注释
  （甲 `:72`＋`:119`／乙 `:72`），未突变时它们与全件一同 PASS ⇒ 前置为真：
  乙跑出的两枚行 corr 实测形如 `84007b6a-1089-4dff-b7ce-99885242e488#call-283-spawn` 与
  `...#call-283-cancel`，task id 是 `84007b6a-...`（MU-1 的红句里逐字可见）。
- **三轴独立性**（每一发突变只打死它对应那一轴，其余轴静默）：
  MU-1（同键塌形）⇒ 作差轴红、前缀轴与 task id 轴**静默**；
  MU-2（回落支 index 抹常量）⇒ 只有甲的单元作差轴红，**全仓其它尺（含本腿乙、票 283 的尺二）逐枚绿**；
  MU-3（保住互不相等、摘掉 task 前缀）⇒ 前缀轴红、作差轴**静默**。
  ⇒ 这把尺不是靠某个 bug 才能满足的死格，也不是"改前绿"一句话糊过去。
- **同形独立复认**：`internal/agent/corr_percall_242_test.go` 的既有尺在未突变时同样 PASS
  （它今天已经覆盖了请求侧作差，见 §5）。

## 4. 名册读数（本腿收工 `git grep -nE 'callCorr' HEAD -- '*.go'`）

- 产码：`internal/agent/loop.go:367`（注释）／`:594`（文档注释）／`:603`（定义）／`:676`（唯一调用点）＝ **4 枚**
- 测试面：`internal/agent/ticket285_corr_distinct_rulers_test.go` 的注释 3 枚 ＋ **直接调用 5 枚**
  （`:51`／`:52`／`:58`／`:59`／`:60`）＝"测试面只有注释文字"那一形在本腿之后不成立
- `internal/tools/ticket283_corr_identity_rulers_test.go:13` 仍是注释（⛔ 本腿未动该件）

## 5. 顶回（盘上现量与票面／派单冲突，逐条具名）

1. **票面 AC#4 的"今天零尺／`callCorr` 零测试引用"在盘上不成立。**
   `internal/agent/corr_percall_242_test.go:99-105` 已经在作差：
   `c1, c2 := p.corr("call_p1"), p.corr("call_p2")` ⇒ `if c1 == c2 { t.Fatalf(...) }`，
   并在 `:109-133` 用各自的 corr 唤醒各自那一枚调用（真路由自证）、`:139`／`:142` 钉前缀与不等于 task id。
   票面尺①（`git grep 'callCorr'`）之所以看不见它，是因为那一枚尺**按名字扫**，
   而 242 的尺经真 loop 派发、不写 `callCorr` 字样；票面尺②的"全仓只有 `ticket87_veto_l2_test.go:173`"
   漏了这一枚（那一枚是 `==` 形状、`grep` 扫的是 `CorrelationID ==`，242 用的是 `c1 == c2`）。
   ⇒ 本腿**没有**把 AC#4 当成"从零建尺"来做，而是把两枚真正缺的洞补齐：
   **(a)** 名册行级的作差（242 只看请求侧内存表，从不读持久行；票面判据写的正是 `rows[...]`）；
   **(b)** 回落支（`callID == ""`）的作差（现量 `#call-` 在测试面零出现，MU-2 下全仓除甲之外逐枚绿）。
   这两发是本腿的可辩护增量；"本格今天零尺"这句判据本身要请按现量复核，⛔ 本腿未翻勾也未改判据句。
2. **票面尺①的"产码 3 枚"现量是 4 枚**：`loop.go:594` 的文档注释里同样写着 `callCorr`，
   那一枚在锚上就在（本腿零产码改动 ⇒ 不可能是我加的）。
3. **`internal/agent` 整包在锚上就是红的**（`TestGoldenSingleToolCall`，`loop_golden_test.go:70`
   期望 `CorrelationID == TaskID`，与票 242 的 per-call 铸形相互矛盾）。
   同一枚文件里 `:341` 那发期望名册行的 corr 等于 task id 却是绿的
   ⇒ 这两枚期望对同一件事给出了不一致的读数，本腿⛔ 未裁定、未放宽、未顺手改（不属 AC#4 射程）。
   `rc` 起手＝收工＝1，本腿没有让它更坏，也没有假装它绿。
4. **派单里"落点自己裁并具名理由（二选一或都给）"我给的是都给**，理由见 §1：
   只给乙就打死 MU-1 打不死 MU-2（乙的夹具走不到回落支），只给甲就没有名册行级读数
   （票面判据的 `rows[...]` 那一形）；两枚都需要，且两枚都不越界。

## 6. 边界自证

- `git diff --name-status ce18b3d..HEAD -- '*.go'` 原文：

```
A	internal/agent/ticket285_corr_distinct_rulers_test.go
A	internal/tools/ticket285_corr_rows_rulers_test.go
```

  ⇒ 只有两枚 **A**（新增测试件），⛔ 零 `M`＝零产码字节改动。
- `internal/agent/loop.go`（三发突变的种刀对象）收工 hash `8eb37e9f59e563b58325d2ebc898b843`
  与起手逐字等值；`git status --porcelain -- internal cmd` 收工为空（原文无输出）。
- ⛔ 未碰 `internal/agent/approval/gate.go`（票 284 随批搭载格不触发）。
- ⛔ 未碰 `internal/tools/loop_approval_test.go`、三枚冻结件、`frontend/**`、`design/**`、
  `docs/PLAN.md`、`docs/specs/**`、`docs/SLO.md`、`internal/observe/thresholds.go`、任何 golden、`allowlist.txt`。
- ⛔ 零 `t.Skip`、⛔ 零断言放宽、⛔ 零 `-overlay`（三发突变全种在盘上并当场还原）。
- 格式与合规：`gofmt -l` 与 `D:\work\base\gopath\bin\gofumpt.exe -l`（只读 `-l`）对两枚新件均**空输出**；
  `tools/d22scan/d22scan.exe -root .` ＝ `clean - no D22 ban violations`（含 `internal/` 521 枚 Go 文件、
  注释与 `_test.go` 全在内），`rc=0`。
- ⛔ 零 push；本腿三笔 commit 全部显式 pathspec。
