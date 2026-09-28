# 176-background-start-port-r1 — 起跑口落地（写手腿 176-r1）

- 交件时刻：09-28 09:3x +08（`date` 现跑，见 §0）
- 派单：`.scratch/wisp/dispatches/2026-09-28-091x-write-176-r1-start-port-fills-the-roster-and-cancel-criteria.md`
- 票面：`.scratch/wisp/issues/176-…-waits-for-one-writer.md`（Progress log 追加一条，**一枚框都没勾**）
- 本程性质：产码＋判据＋端到端台件；**只 commit，未 push**

## 0. 起手四件（现跑）

```
$ date
Mon Sep 28 09:09:12 CST 2026
$ git rev-parse --abbrev-ref HEAD
dev
$ git rev-parse HEAD
8058d407bd9146f050d3cca060fdd219538df11a
$ git log --oneline -3
8058d407 派单 176-r1（起跑口·硬顶 45）：…
386716a6 ledger(A358)：175-r2 让 177 的豁免第一次真生效…
a18cafbc evidence(175-r2): 桥级三发的两向读数＋四枚变异读数…
$ git status --porcelain -- internal/ cmd/          # 空＝起手净
（无输出；命令后接 echo 分隔符证明它确实执行了）
```

起手 `sh .scratch/wisp/probes/154/gate-clauses.sh` 的逐名红腿集合（交件时对比用）：

```
# BAD  腿=G6neg 声明=ring 基线=1枚 实测=2枚 因=新增未成对（票 171 AC#2：实测 > 基线）
# 腿数＝14 声明与实测不符＝1
# 腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0
```

即：**起手就只有 G6neg 一枚 BAD**（票 178 已立）。本枚**没抬它的基线、没往测试里补 `OpenTask` 字面**。

## 1. 三处落点（逐枚：文件／行号／用的是哪枚现成函数）

| # | 文件 | 落点 | 用的现成物 |
|---|---|---|---|
| 起跑 | `cmd/wisp/run.go` `:633-634`（`execute()` 组装点，`:622` 那一发 `loop.Run` 换掉） | `bg := loop.RunAsync(ctx, task)` → `res := bg.Wait()` | `internal/agent` 现成的 `(*Loop).RunAsync`（`loop.go:321`，**没新造同类函数**）；id 取 `res.TaskID`（环路 `newTaskID()` 铸的那枚，组合根不自造） |
| 回填 | 同上 `:635-641` | `tools.TaskBackfill{Roster, Spills}.Backfill(bg.Root().Ctx, res.TaskID, res.Text)`；理由非空 → `rt.auditf` 响亮报出 | `RunningTask.Wait()`（`loop.go:308`，同时 join 结果信号与登记册句柄）·`RunningTask.Root()`·`Loop.Budgets()`（`loop.go:253`）·`agent.NewSpiller`（`spill.go:43`）·`agent.Spiller.Prepare`（`spill.go:84`）·`TaskRoster.Record`（`task.go:123`，**本体一字节未动**） |
| 写侧本体 | `internal/tools/task_backfill.go`（新，121 行） | `Backfill`＝写前看一眼 `Stopped(ctx)`＋只登记真落盘的路径 | `tools.Stopped`（`cancel.go:74`）·`agent.BudgetsFor`（`budgets.go:84`） |

命名不裸送：产物 key＝`taskArtifactPrefix + taskID`，前缀字面 `agent-task-` 抄的是环路自己那枚先例（`loop.go:324` `reg.Spawn("agent-task-"+id, "agent", …)`）。**没改 `Spiller` 的覆盖语义**（`spill.go:117-131` 的 last-writer-wins 原样）。

D38：**没为后台记录员新登记名**（复用现成前缀）。台件里 observe 没打出 `CategoryUnknown` 警告（起跑那两发的日志只有 winsec/MODE-READ/C25 那几行）。

落地后的三个数（命令原文，复算同一条）：

```
grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l   # 1（run.go:633）
grep -rn 'Roster\.Record(' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l # 1（task_backfill.go:117）
grep -rn 'ArtifactPath:' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l   # 0 —— 名册不再靠字面量填充
grep -rn 'rec\.ArtifactPath = ' --include='*.go' internal/ cmd/ | grep -v _test.go    # 1 处，task_backfill.go:112（真落盘才赋）
```

⚠ 票 176 AC#1 那句"若任一枚变成 ≥1，本票作废"——**本枚就是那枚 ≥1**，票面已按派单 §7 只追加、不翻勾，作废判断交给编排者。

## 2. 取消语义那对红／绿（票 176 AC#3）

判据件：`internal/tools/ticket176r1_start_port_test.go`（新，四枚常驻判据 L-1…L-4）。挂点＝`tools` 侧＋现成 `Stopped(ctx)`；**没挂桥、没挂 Loop**。

绿色（落地码，真跑）：

```
go test -count=1 -run '176r1' ./internal/tools/
--- PASS: TestBackfillFilesTheRealSpilledArtifact176r1 (0.02s)
--- PASS: TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1 (0.01s)
--- PASS: TestNoRosterWriteAfterCancel176r1 (0.00s)
--- PASS: TestTerminalWriteCarrierStaysOutOfScope176r1 (0.00s)
ok  github.com/CarlosShao/wisp/internal/tools 0.072s
```

L-3 的绿色那发是**真起一个慢任务**：`observe.Registry.Spawn("agent-task-"+id, …)` → 取消 → `<-h.Done()` → `root.Wait()` 归零 → 才 `Backfill(root.Ctx,…)`；断言 `Look(id)` 不命中、`Count()` 不增，并且**同一条腿带对照组**（不取消 → 必须真写进去），所以它不能靠"永远不写"蒙绿。

变异读数（**先证落地再跑**；四枚变异都是 `go test -overlay` 换成 `probes/176/r1/mut/*.go`，跟踪文件一字节没动）：

```
grep -n 落地证明：
M1  112:			rec.ArtifactPath = ""
M2  104:		sp, err := b.Spills.Prepare(taskID, text)
M3   97:	if stopped, why := Stopped(ctx); false && stopped {
M4   97:	if stopped, why := Stopped(ctx); true || stopped {

M1 → --- FAIL: TestBackfillFilesTheRealSpilledArtifact176r1
       :93: 超长输出没有产物路径：这就是 PLAN.md:2564 的「只截不指」
M2 → --- FAIL: TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1
       :145: 产物名 "tool-output-e3f1c0a97b2d4865.txt" 不带 agent-task- 前缀…
       :158: 模型 supplied 的 call id 与任务 id 落在同一个产物名上
M3 → --- FAIL: TestNoRosterWriteAfterCancel176r1
       :224: 任务被取消之后它的输出还是被写进名册了：Look() 命中（why=""）
       :227: 取消之后 Count() = 1, want 0
       :230: 跳过回填却没交代理由
M4 → --- FAIL: TestTerminalWriteCarrierStaysOutOfScope176r1
       :269: 终态载体上回填被拒了：任务已取消，不再写它的输出
```

"我显式排除了终态写入那一支"的凭据＝L-4 本体：`context.WithoutCancel(cancelledRootCtx)`（逐字是 `loop.go:956-958` 那枚载体的形状）→ `Stopped` 必须说"没停"，且回填必须**真写进去**；M4 一发它当场红。`Record` 里**没**加"取消即拒写"（`task.go:121-122` 的 last-writer-wins 一字节未动，`git diff 8058d407 HEAD -- internal/tools/task.go` 为空）。

## 3. 端到端那一发（派单 §4，本票最想要的东西）

台件：`.scratch/wisp/probes/176/r1/zz176r1_e2e_test.go` ＋ `overlay-e2e.json`（`go test -overlay` 把它编进 `cmd/wisp`，跟踪目录没多文件）；读数 `logs/last-run.txt`。合法注入接缝＝**CLI `wisp run`**（`runTextTask` 全入口）＋真 OpenAI-chat SSE 端点（`httptest`，逐帧 `data:` ＋ usage ＋ `[DONE]`），真环路／真桥／真 C26／真审批门／真文件，零 mock 顶替真件。

**前半截拿到了（绿）**：

```
run A（真 CLI 跑一个后台任务）：
  wisp run: 任务 0ac4c512-b05e-4bad-8a54-a5da202970f4 结束（completed，1 轮，0 次工具调用…）
  名册命中同一枚 id；ArtifactPath 指向真落盘的 <data>\artifacts\tool-output-agent-task-0ac4c512-….txt
  产物字节 == 任务打印的 20000 字节（os.ReadFile 对全文比对通过）
  控制台与审计里没有出现「后台任务的输出没有进名册」
```

⇒ **这是这条管子上第一次"真机端到端"的写侧读数**：起跑口落地前，这里连一条记录都不会有。

**后半截没拿到（红，且我没改判据）**：模型在同一进程里调 `task.output`，被环路拒了 7 次后卡进 Stuck：

```
[工具 task.output -> error]  ×7
wisp run: 任务 1f8003a3-… 结束（stuck，8 轮，7 次工具调用…）
名册里回给模型的那一句（台件把 tool 消息原样打了出来）：
  tool||call-176r1-task-output|tool error: 风险未分级且直通开关关闭，已拒绝执行
```

出处：`internal/agent/loop.go:622` 调 `decideRisk`，而 `:773-790` 那个 switch 只把**声明为 L1/L2** 的当作有分级，`default` 一支把**声明 L0** 与"未分级"混在一起，`Config.PassThroughUnclassifiedRisk` 为 false 就拒。`internal/tools/bridge.go:1092 levelString` 对 `risk.L0` 返回 `memory.RiskL0`，所以内置 L0 工具走的正是 default。全仓把该开关设成 true 的地方只有测试夹具（`internal/agent/harness_test.go:116`、`probes/151-accept/accept151_control_test.go`），组合根没设（`cmd/wisp/run.go:606-613` 的 `agent.Config`）。

⇒ 结论只敢写到这里：**CLI 这条路上"模型发起的声明-L0 工具调用"从来没被任何既有用例跑过**（risk 级／桥级全绿≠端到端通，本枚第三次验证这句话）。修法两支都不在本枚写面内——改 `internal/agent`（禁区）或在组合根打开 `PassThroughUnclassifiedRisk`（宿主级策略开关，且会把"未分级"与"声明 L0"一并放过＝削弱门）——**两支都没动，等编排者裁**。
`查不到这个任务`／`任务名册未接线`／`路径授权判定者未接线` 三个字符串在这发读数里都＝0 命中（`grep -c` 见 `logs/last-run.txt` 复算），所以拒绝**不是**名册侧的失败：名册有记录、路径是真的、失败发生在**调用进入工具之前**。

## 4. 三条禁区逐条自证

- **(i) 没给 `agent.Loop` 加收 taskID 的导出方法**：`internal/agent/**` 零字节（`git diff --numstat 8058d407 HEAD -- internal/agent/` 空）。🔴 **同时把"这条腿安静≠合规"抄一遍**：G3 按字面 `taskID` 抓，我把参数改名叫 `id` 也能躲过去，所以本枚的凭据不是"G3 没响"，而是"我没新增任何 Loop 方法"这一条 numstat。
- **(ii) 宿主内部 artifacts 写入没做成受门控的 Tool**：回填用的是 `agent.Spiller`（`spill.go:29-31` 逐字 "deliberately not a gated tool"），`internal/tools/registry.go`／C1 面（`tool.go`）一字节未动，没注册新工具名。
- **(iii) 裸 `go func(`**：产码里 **0 枚**（`grep -rn 'go func(' internal/tools/task_backfill.go cmd/wisp/run.go` 空）；台件与判据里的并发全走 `observe.Registry.Spawn`（owner＋recover 在 `internal/observe/goroutine.go`）。`scripts/d22scan.sh` rc=0 是同条 claim 的仪器读数。
- **(iv) `filepath.Clean|Abs` 越界**：新代码里 0 处；路径决策全部交给 C26（`PathCanonicalizer` 在 `pointerNotice` 那一侧，本枚没碰），`risk.PathResolver` 之外的文件系统决策没增加。
- **G2 的雷（派单 §1）避开了**：`run.go` 新写的说明里不含 `OpenTask|CloseTask` 字面（`grep -c -E 'OpenTask|CloseTask' cmd/wisp/run.go` 落地后仍＝起手那两行：`:564` 注释与 `:581` 真实调用），交件终态读数见回禀第 7 节。
- **契约轴**：`docs/PLAN.md`／`docs/specs/**`／`docs/reports/**`／`internal/risk/**`／`internal/panel/**`／`thresholds.go`／golden／`allowlist.txt`／审批超时常量——本枚零字节（`git diff --name-only 8058d407 HEAD` 的清单只有：`cmd/wisp/run.go`、`internal/tools/task_backfill.go`、`internal/tools/ticket176r1_start_port_test.go`、`.scratch/wisp/probes/176/r1/**`、票面与本表）。

## 5. 本程没测什么（别把它读成"都过了"）

- **没**跑全仓 `go test ./...`；跑的是 `-count=1 ./internal/risk/ ./internal/tools/ ./cmd/wisp/`。
- **没**跑 `probes/161/r6/flip-declaration.sh`（派单 §5 明禁）。
- 端到端后半截**没拿到**（原因见 §3），所以"模型读回指针 → `fs.read` 不被自家门拒"这一发在本枚**没有读数**——它现在被 §3 那枚 L0/未分级混判挡住，挡在票 177 的豁免之前，所以**它既不能算通过也不能算被豁免拦下**。
- `TaskRoster` 的跨进程／重启行为**没测**（票 164 定案②：v1 不做）。
- `task.list`／`task.cancel` 的实现本体**没做**（`PLAN.md:1531` DEFERRED 在册）。
- 取消判据的射程是"回填这一支"；**没**覆盖"别的写者绕过 `Backfill` 直接 `Record`"那一形（`Record` 仍是谁都能调的公开方法，这是 last-writer-wins 裁决留下的形状，不是本枚的疏忽修复项）。

## 6. 交件

提交（只 commit，**未 push**）：产码两枚（`internal/tools/task_backfill.go`、`cmd/wisp/run.go`）／判据一枚（`internal/tools/ticket176r1_start_port_test.go`）／台件＋本表一枚。终态门禁**在最后一次提交之后**现跑，读数在回禀第 7 节（本表不抄未跑的数）。
