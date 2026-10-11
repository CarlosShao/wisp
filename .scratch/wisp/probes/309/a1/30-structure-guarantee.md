# 309-a1 · 30 `AC#2` 的料：⛔ 依赖未来接线的结构保证，仓里有几种形（⛔ 裁，只交料）

尺与三态同前（HEAD blob @`7dc753e3`；`git status --porcelain -- internal cmd`＝0 行；⛔ 编译面命令 ⇒ 全部静态）。

## 一、`bridge.go:872-875` 那段点名的三类未来候选者（逐字）

尺＝`logs/bridge.go` blob `sed -n '868,883p'`（票面 `:26` 说这段的行号会漂 ⇒ 这里给**内容锚**：`// WHICH LEGS STILL HAVE NO OWNER HERE - ticket 154`）。三类候选者逐字在 `:872-875`：

```go
// So ask one question about your own leg: does the TaskID you dispatch with
// come from that boundary? If it comes from anywhere else - a panel or ball
// host inventing its own id, a retry wrapper, a scheduled task carrying one id
// across turns - nothing in this tree closes it, and no test will tell you so.
```

同一块里还有两枚**行号已漂**的引用（料，属 `AC#3`，⛔ 本腿改）：`:871` 逐字 `mints itself (Run/RunAsync -> newTaskID, internal/agent/loop.go:332/:321)`；现量尺＝`grep -n "^func (l \*Loop) Run\|^func (l \*Loop) RunAsync\|newTaskID" logs/loop.go` ⇒ `RunAsync` 起 `:328`（铸造点 `:329`）／`Run` 起 `:339`（`:340`）⇒ **两枚号各漂 4／19 行**（票面说"两枚号全漂"＝成立）。另一枚：`:861` 逐字 `internal/agent/loop.go:366`（TEXT loop 的 defer）⇒ 本腿⛔ 现读该号（欠读数 U-5）。

## 二、三类候选者在仓里现在有没有半成品（逐类，尺名＋射程）

### 候选者 1「a panel or ball host inventing its own id」——**有 1 枚真半成品，就是本名册的 ②**

- 尺＝`grep -rn "newTaskID" HEAD -- internal cmd` ⇒ 命中只在 `internal/agent/loop.go`（`:329`／`:340`／定义 `:1153`）⇒ **环路之外的 id 全部是自造串**。
- 现量自造串：`cmd/wisp/panel_assets.go:174` 逐字 `const taintSourceScopeID = "panel-assets-l2"`，且在 `:69` 被复用为 `CorrelationID: "panel-assets-l2"`。⇒ 一枚 CLI 宿主"自己命名 id 并当作用域用"的**完整形状已经躺在仓里**（⛔ 过桥、⛔ 落审计），它是候选者 1 的原型机。
- 面板/球宿主本体确实在：`git grep -ln "retry" HEAD -- cmd/wisp` 命中 `cmd/wisp/panel_host_windows.go`·`cmd/wisp/resident_ball_windows.go`；`cmd/wisp/panel_inbound_guards_35r3_test.go` 一族钉的是 `bridge.go:132-141` 那三支**身份守卫**（`inboundPanelSource`＝"bridge.go 接受的唯一 sender 身份"），**没有一枚**给面板递 `ToolRequest`。
- 关键现量：生产 `.Execute(` 调用点（尺＝`git grep -n "\.Execute(" HEAD -- internal cmd` 滤 `_test.go`、滤 `internal/agent/loop.go`、滤 `internal/tools/bridge.go`）＝**1 行**＝`internal/tools/subagent_197.go:581` 逐字 `return p.inner.Execute(ctx, req)` ⇒ 那是 **Tool→Tool 的转发**（`subagentToolProvider` 装饰器，`:573-581`），⛔ 是"宿主派发桥"。

### 候选者 2「a retry wrapper」——**有形状、无 id 语义**

- 现成的装饰器就是上面那枚 `internal/tools/subagent_197.go:573-581`：它**逐字转发同一个 `req`**（⛔ 造、⛔ 洗 `TaskID`）⇒ 一枚重试包装器若照这个形状写，id 会原样穿过＝**这正是"手搭 ToolRequest"最容易顺手做出来的地方**（改天有人加 `req.TaskID = newID` 或 `req = agent.ToolRequest{...}` 重建，本票那一肢就第一次真在生产可达）。
- 其余 `retry` 字样在 `internal/agent/spill.go:165-167`/`:335-336`（工件写重试）与 `internal/agent/approval/approval.go:532`（注释），**均⛔ 碰 ToolRequest**。尺＝上面两条 `git grep -n "Retry\|retry"`。

### 候选者 3「a scheduled task carrying one id across turns」——**只有 doc，零实现**

- 尺＝`git ls-tree -r --name-only HEAD -- internal/agent/scheduler` ⇒ **1 个文件**＝`internal/agent/scheduler/doc.go`；`git grep -n "TaskID\|Execute\|Run" HEAD -- internal/agent/scheduler` ⇒ **0 命中**。
- 宿主侧的定时来源存在但⛔ 接桥：`cmd/wisp/resident_task_source_windows.go`（`git grep -l` 命中 cron/schedule 那一族）。⇒ 候选者 3 ＝**纸面**，本票⛔ 为它造路径（票面 `AC#2` 末句）。

## 三、"能不能有一枚⛔ 依赖未来接线的结构保证"——能，仓里已有**四族现成尺形**（逐族给取证尺长什么样）

> 共同判据（票面 `AC#2` 原话）＝"将来任一 caller 手搭 `ToolRequest` 递进空 `TaskID`，要让它在盘上可核地响"。下面四族都⛔ 需要任何新接线就能落地。

### 形 A「枚数钉」（gate-clauses 那族，⛔ 依赖接线、⛔ 依赖 CI）

- 现成范本逐字在 `.scratch/wisp/probes/154/gate-clauses.sh`：`:425-427` 定义 `G1 生产码里构造 ToolRequest（排 internal/agent/ 自身）`，pattern `'ToolRequest\{'`、pathspec `':!*_test.go' ':!internal/agent/*'`、`want_n 0`；`:430-432` `G1b 生产码里任意 .Execute( 调用点（排 loop.go 与 bridge.go 自身）` `want_n 0`；`:444-446` `G3 Q-56 那一支落地：Loop 上出现收 taskID 的导出方法` pattern `'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID'` `want_n 0`；头部 `:81` 名册逐字 `LEGS_EXPECT='G1 G1b G2 G3 G4 G5 G5pos G5neg G6 G6pos G6neg G7 G7pos G7neg'`。
- 取证尺长这样（可逐字重跑，本腿已跑 G1/G1b/G3 的等价物）：`git grep -nE '<pattern>' <锚> -- <pathspec> | grep -c .` ⇄ 文件头登记的 `want_n`；不等＝那枚门响。脚本自述基线逐字（`:421`）＝"本程在锚 f1b99a70 现量：G1=0／G1b=0／G2=2／G3=0／G4=0"。
- 现量对账：本腿在同一锚上跑等价 grep ⇒ **G3＝0**（`Loop` 的 30 枚方法里⛔ 导出方法收 `taskID`）＝与文件登记一致；**G1b 等价尺＝1 行**（`subagent_197.go:581`）而脚本登记 0 ⇒ **两把尺的射程差别在 `internal/tools` 自身是否算"宿主派发"**（脚本排除项只列 loop.go/bridge.go，没列 `:581` 这种 provider 转发 ⇒ 它多半在 f1b99a70 时那行还不存在，或 pattern 更窄）。**具名欠读数 U-4**：现跑 `bash .scratch/wisp/probes/154/gate-clauses.sh` 才知道 G1b 今日色与本票要不要顺手重新登记 `want_n`（本腿⛔ 跑它，怕与 `72-r1` 抢面并产出大读数）。
- 弱点（料）：这族尺**长在 `.scratch/` 里**、靠人跑；票面 `AC#4` 的 CI 门禁面（`go build ./...`／`go vet ./internal/tools/`／`sh scripts/d22scan.sh`）⛔ 包含它。

### 形 B「AST 枚数钉」（**仓里有直接范本，且它就在 `cmd/wisp` 里跑着**）

- 范本逐字：`cmd/wisp/panel_assets_143_test.go:267` `file, err := parser.ParseFile(fset, "panel_assets.go", nil, 0)` ⇒ 一枚**用 go/parser 读生产源码**的钉已经在本票 ② 的那枚测试文件里（票 143 交的）。另三族同类：`cmd/wisp/leg_dispatch_gate_133_test.go:372/438`、`cmd/wisp/leg_sink_gate_131_test.go:578-586`、`cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go:277`、`cmd/wisp/config_receipt_255_test.go:295`（都是"遍历目录里非 `_test.go` 的 `.go`＋解析/正则数一遍"的形状），还有最直白的文本形 `cmd/wisp/firstrun_257_test.go:430` `source, err := os.ReadFile(filepath.Join(".", "firstrun.go"))`。
- 取证尺会长这样：在 `internal/tools/*_test.go` 里枚一枚 `Test309…`，`os.ReadFile("bridge.go")`＋`go/parser` 走 `CompositeLit`，断言"`agent.ToolRequest{` 的**生产构造点枚数**＝1（现量＝`internal/agent/loop.go:675`）或"`assessorFor` 调用点枚数"＝1；⛔ 依赖任何未来接线，一行新 caller 就红。
- ⚠ 射程要说清的代价：源文件路径依赖 `go test` 的工作目录（`-run ./internal/tools` 时 `.` ＝包目录 ⇒ 能读到 `bridge.go`，`cmd/wisp` 那几枚范本正是这么读的）；跨包构造点（`internal/agent/loop.go`）在 `internal/tools` 包里读⛔ 到，得走 `../../internal/agent/loop.go` 之类的相对路径——**这是"枚数钉"落到 Go 测试里唯一的工程别扭处**，形 A（shell）没有它。
- ⚠ 与本票写面冲突（⛔ 裁，只登记）：票面 `AC#4` 允许写面＝`internal/tools/**`＋`probes/309/**`＋票尾 1 条 ⇒ 若枚数钉要覆盖 `cmd/wisp` 或 `internal/agent` 的构造点，**写面会出界**。

### 形 C「入口断言／拒绝」＝ⓐ 那一形（**今天就被测试可达，⛔ 等新接线**）

- 现量证据链（本腿逐枚数过）：全仓 `ToolRequest{` **60 枚**＝生产 1（`internal/agent/loop.go:675`）＋测试构造 57＋`reflect.TypeOf(agent.ToolRequest{})` 2（`internal/tools/grant_test.go:509`、`internal/tools/ticket90_test.go:454`，它们⛔ 构造请求）；测试里递空 id 的载体**恰好 2 枚**：字面量 1 枚 `internal/tools/failclosed_236_teeth_test.go:288`（`TaskID: "", CorrelationID: ""`）＋参数化 1 枚 `internal/tools/task_cancel_221_legs_test.go:119`（`TaskID: callerID`）由 `failclosed_236_teeth_test.go:118`（`x.cancel(t, "", target)`）喂空。
- ⇒ 结论性事实（料）：**ⓐ 那条新分支今天就有测试可达路径**，所以它作为结构保证⛔ 依赖未来接线；代价见 `20-shape-a-cost.md`（2 枚断言必红、③/④ 文案被遮蔽）。
- 反向事实：**ⓐ 单用会遮蔽 ③/④ 两支**，作为"结构保证"它把保证从"注册/盖戳不许没有父"降级成"不许没有身份"——语义不亏（`subagent_197.go:474-477` 的 `case parentID == ""` 盖戳侧同受保护），教育面亏。

### 形 D「reflect API 钉」（另一族⛔ 依赖接线的现成尺）

- 范本逐字：`internal/tools/grant_test.go:509-518`（遍历 `agent.ToolRequest` 字段名，命中 `session|grant` 即 `t.Errorf`）与 `internal/tools/ticket90_test.go:454-463`（禁 `mode|allow|approve|auto`）。它们证明本仓接受"用 reflect 把 API 形状钉死"这种**零接线**保证。
- 对 `AC#2` 可用的形：一枚"字段名/构造器"钉（例如禁止出现任何让 caller 声明"我不需要任务 id"的字段）——⚠ 它**挡不住**手搭 `ToolRequest` 递空串（那是形 C 的活），只挡"把空 id 合法化"那一条路。**别把它当成 AC#2 的交付物**，它是 AC#2 的补充。

### 形 E（⛔ 现有尺形，但ⓘ 想法）"⛔ 不许存在⛔ 探测器的那一枚"

- 如果 `AC#0` 裁ⓑ，`AC#2` 的"响"只剩一条诚实路径＝**形 B 或形 C**；ⓑ 的注释自身⛔ 构成保证。反过来说：ⓑ＋形 B（枚数钉）＝**行为不变、新增一枚"手搭就红"**，ⓑ＋形 C＝就是ⓐ。⇒ 两形的真实差＝**"改行为换保证"⇄"加一枚保证钉"**。

## 四、本格欠读数（⛔ 实测⛔ 猜，逐枚归编排者）

- U-4 现跑 `bash .scratch/wisp/probes/154/gate-clauses.sh` 取 G1/G1b/G3 今日命中行数与 `want_n` 是否一致（静态等价尺给出 G1b＝1 vs 登记＝0，差别归因于 pattern 细节）。
- U-5 `bridge.go:861` 那句 `internal/agent/loop.go:366` 的行号是否也漂（本腿只核了 `:332/:321` 那两枚）。
