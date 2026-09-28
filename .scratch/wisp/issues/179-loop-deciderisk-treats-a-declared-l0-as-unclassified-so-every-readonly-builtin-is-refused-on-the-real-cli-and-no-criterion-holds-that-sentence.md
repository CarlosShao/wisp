# 179 — **`Loop.decideRisk` 把"声明为 L0"当成"未分级"**：真机 CLI 上每一枚只读内置工具都被自家门拒（逐字文案 `风险未分级且直通开关关闭，已拒绝执行`）——而全仓**零枚**用例钉这句话，测试桥具恰好把 `PassThroughUnclassifiedRisk` 设成 `true`，把这个洞整片盖住

- Status: **ready-for-agent（写码＋判据）**。本票是**缺陷修复**，方向由已定案的 D4 决定，**不是契约变更**（理由见下面"裁定轴"）。
- 来源：`176-r1` 第一次把真机 CLI 端到端腿跑出来（`.scratch/wisp/probes/176/r1/logs/last-run.txt`，27580 字节，**逐字 7 次拒绝**），编排者 09-28 09:5x 在 HEAD `cf527f32` 上**自己现读复算**了整条链（每把尺贴在下面）。台账 `A359`。
- 关联：**D4**（`docs/PLAN.md:122`，那张三级表的 L0 行）· 票 21（`bridge.go:210-213` 那段注释早就点名了这一枚 hazard）· 票 176 AC#3/AC#4/AC#5 · 票 177 AC#3（它的端到端条件格就卡在这一枚上）· 票 175 AC#5 · 票 174 AC#2b（**正交**，见末节）。

## 这是什么（人话）

模型要让宿主读一个文件、列一个目录、或者把超长输出的那份副本读回来——这三件事在方案里是**最低风险档（L0：只读、无副作用，"直接执行，不打断"）**。可是编排 loop 在放行之前自己先问了一题"这枚工具是几级"，它只准备了两个答案（L1、L2），**其余一律当作"没分级"**，然后按"没分级＋直通开关关着"的规则**拒绝执行**。

于是真机上：**最常用、风险最低的那三枚工具，全部被自家的门挡在门外**。上一轮票 176 好不容易把后台任务的起跑口接通、第一次跑到端到端，读回来的就是这个——同一次运行里连续 7 次拒绝，模型反复重试，最后被判"卡住"。

而这件事**在既有测试里完全看不见**：测试用的那套桥具把"未分级就放行"的开关打开了（`harness_test.go:116`），它自己的两枚测试工具（`testtools_test.go:62`、`:67`）又恰好声明成 L0——**所有既有 loop 用例都是靠那枚开关过去的**，没人问过"开关关掉时，声明为 L0 的应当发生什么"。

## 现量证据（每把尺可复制；锚 `cf527f32`，09-28 09:5x 本程现跑）

1. **判定那一支**（尺：`sed -n '773,790p' internal/agent/loop.go`）：
   `switch risk { case RiskL1, RiskL2: … default: if !PassThroughUnclassifiedRisk { reject "风险未分级且直通开关关闭，已拒绝执行" } }`
   ⇒ `default:` 同时吃进**两条完全不同的东西**：真正没分级的 `""`，和 `levelString` 正常产出的 `"L0"`。
2. **`"L0"` 是谁产的**（尺：`sed -n '1092,1101p' internal/tools/bridge.go`）：`levelString` 的 `default: return memory.RiskL0`，而 `internal/risk/assessor.go:40-44` 里 `L0` 是一个**合法档位**（`type Level int` / `L0 Level = iota`），不是"没值"。
   桥把声明档送进 `ToolInfo.RiskLevel` 的那一跳在 `bridge.go:224`。
3. **"未分级"这个词在本仓的逐字定义**（尺：`sed -n '32,38p' internal/agent/tools.go`）：`RiskUnclassified = ""`。⇒ 判定该拿 `== RiskUnclassified` 问，而不是拿 `default:` 兜。
4. **同一枚文件里已经有正确写法**（尺：`sed -n '1084,1089p' internal/agent/loop.go`）：`riskLabel` 逐字 `if r == RiskUnclassified { return "未分级" }`。
   ⇒ **同一枚 `loop.go`、两个函数、对 `"L0"` 是不是"未分级"给出相反答案**：`list_tools` 报给模型的是 `[内置/L0]`，门却按"未分级"拒它。
5. **这句文案全仓只有一处、且零枚用例钉它**（尺：`grep -rn "风险未分级" --include=*.go .` ⇒ **1 行**（`loop.go:785` 产码），`_test.go` 里 **0 命中**）⇒ **无牙**：改它不会有任何现有用例报警，它错了也没人知道。
6. **真机上的分母＝三枚注册工具声明为 L0**（尺：`grep -rn "Declared:     risk.L0" --include=*.go internal/ | grep -v _test.go`）：
   `internal/tools/fs.go:300`（`fs.read`）、`fs.go:313`（`fs.list`）、`internal/tools/task.go:285`（`task.output`）。
   注册发生在 `cmd/wisp/run.go:347-370`。⇒ **CLI 上只读面今天整片不可用。**
7. **端到端那一发（红之前的凭证）**：`probes/176/r1/logs/last-run.txt`；尺 `grep -ac "风险未分级" .scratch/wisp/probes/176/r1/logs/last-run.txt` ⇒ **7**。7 次全是同名的 `task.output` 重试，不是 7 枚工具。
8. **真机没设那枚开关**（尺：`grep -rn "PassThroughUnclassifiedRisk" --include=*.go .` ⇒ 生产码里**零枚**赋值点；赋 `true` 的唯一一处是 `internal/agent/harness_test.go:116`，另有一处显式 `false` 在 `internal/tools/loop_approval_test.go:128`）⇒ 缺陷**只在真机形状上成立**，而真机形状上一轮才第一次被跑出来（票 176）。
   ⚠ **顺带更正一处转述**：`176-r1` 的表（`docs/evidence/s1/176-background-start-port-r1.md` §3 末段）把
   `.scratch/wisp/probes/151-accept/accept151_control_test.go` 也算成"把开关设成 `true` 的夹具"之一——
   尺：`grep -n "PassThroughUnclassifiedRisk" .scratch/wisp/probes/151-accept/accept151_control_test.go` ⇒ **只有 `:14` 一行注释**，没有赋值。
   ⇒ 全仓唯一把该开关设成 `true` 的地方是 `internal/agent/harness_test.go:116`。这一处**不改变任何结论**（两支禁修方向它都碰不到），只是名数要准。

## 裁定轴：这是"修 bug 恢复 D4 一致"，不是"改契约"（写手据此不必再来问我，但也**不许越过下面 AC#5 的禁令**）

- D4 那张表逐字（`docs/PLAN.md:127-129`）：`| **L0** | 只读，无副作用 | 搜索、查询、读文件、读剪贴板 | 直接执行，不打断 |`。
- `loop.go:755-770` 那段注释**自己**只列了两形（"a declared L1/L2 call …"、"an unclassified call still obeys PassThroughUnclassifiedRisk"）；`tools.go:21-25` 也把开关语义写成针对 `""`。**声明 L0 该怎么走，注释里根本没写**——代码把没写的那一支并进了 `default:`。
- ⇒ 修它是**让代码回到已定案的 D4 与既有注释的语义**，不动 D1–D47／C1–C32／R1–R9／D43 任何一格，**不需要新的批准**。
- ⚠ **"修完会把谁顶出去"本票已答**（写手要复核，不要照抄）：`decideRisk` 只是**下界预门**；真正的裁决在 `Bridge.Execute` 里逐发做——C3 能力核查（`bridge.go:255-256`）→ C19 `Assess`（`:272`）→ 权限模式筛（`:286`）→ 门。所以"让声明 L0 在 loop 里通过"**不会**让任何一枚 L1/L2 少问一次；未知工具名那一形也不受影响：`info := byName[c.Name]`（`loop.go:610`）取不到就是零值 `ToolInfo`，`RiskLevel == ""`，仍走 fail-closed。

## AC（每格都要答"这一发在**未修码**上响不响"；先测→再写→再提交）

- [ ] **AC#1 红之前的读数不许"待跑"**：在**当前 HEAD**（本票开出来这一刻）复跑一次 `176-r1` 那发端到端台件，逐字贴任意一次拒绝行＋尾态（今天已有读数＝上面第 7 把尺）。⚠ 这一格**只登记事实**，修完之后它不翻绿——绿的是 AC#8。
- [ ] **AC#2 正向判据（常驻，落 `internal/agent/`）**：把直通开关**关着**，声明为 `L0` 的一发必须**真的执行**——回执里有工具产出的正文（不是拒绝文案），且 journal 那行的 `decision` 写的是放行。⚠ 断言要钉在"**执行了**"上，不许钉"没报错"：今天它报的是 policy reject（`ClassPermissionDenied`，见 `loop.go:624-625`），所以那条链是**响得出来**的。落地时**逐枚贴两向读数**（未修码红／修后绿），并把尺贴在表里。
- [ ] **AC#3 反向判据（这枚最值钱，与 AC#2 成对，缺一形＝装饰）**：同样关着开关，构造一发 `RiskLevel == ""` 的调用（**目录里没有那个工具名**那一形最省事，因为它天然产零值 `ToolInfo`），要求**仍然被拒**、文案仍是那句。⚠ 任何修法若把 `default:` 改成"一律放行"，这一发就变安静＝**把门洗掉了**，本票判不通过。这一格在未修码上今天也绿（它本来就是对的），所以它是**"不许弄坏"的守卫**，不许拿来充当本票的新增牙。
- [ ] **AC#4 另一枚"不许弄坏"守卫**：L1/L2 两形路由**一字不变**——`AdmitTask == nil` 时仍拒（逐字 `"该操作属于 " + risk + " 级，审批通道尚未接入（ticket 21），已拒绝执行"`，`loop.go:777-779`）；非 nil 时仍返回"放行但不记账"。⚠ 同 AC#3：这两枚今天也绿，只能当守卫。
- [ ] **AC#5 禁用毒修法（这格是给我自己看的闸）**：**不许**通过在组合根（`cmd/wisp/run.go`）把 `PassThroughUnclassifiedRisk` 设成 `true` 来"解决"本票。理由逐字：那枚开关的语义是"给**未分级**调用用的**宿主级政策开关**，不是裁决"（`tools.go:21-25`、`loop.go:767-770` 两处注释），设成 `true` 会把**真未分级**那一发（AC#3）一并放行——那是**放宽门**，与 AGENTS §1.1"不许为了变绿放宽任何断言"同一族，只不过这次放宽的是门不是断言。⇒ 写手若判"该设"，**停手上报**，别自己决定。
- [ ] **AC#6 顺带把两枚过期指针改成实话（只改注释，零产码行为）**：
  ① `cmd/wisp/run.go:357` 逐字仍写 "RunAsync still has zero production call sites (measured, …)"，而**同一枚文件 `:633` 就是那枚调用点**（尺：`grep -rn "\.RunAsync(" --include=*.go cmd/ internal/ | grep -v _test.go` ⇒ **1 行**）。
  ② `bridge.go:212` 引 `internal/agent/loop.go:719-738`，实际 `decideRisk` 现在在 `773-790`。
  ⚠ 只改**指向与枚数**，不改原话的判断内容；这是"注释里编造出一个已不存在的调用者/枚数"那一族（票 97 已定过案）。
- [ ] **AC#7 门禁**：逐包 `go test -count=1 ./internal/agent/ ./internal/tools/ ./internal/risk/`（**禁全仓 `./...`**）三数＋名册两向 `comm` 差集；`sh scripts/d22scan.sh`；`sh .scratch/wisp/probes/154/gate-clauses.sh`；`gofumpt -l` 对自己名下两枚 `.go`。⚠ **`./cmd/wisp/` 在本机测不到任何东西**（票 98 仍 open；09-28 09:59 我现跑：`go test -count=1 -run XXX_NONE_PKG ./cmd/wisp/` 同样 `exit status 0xc0000135`＝加载期 `sherpa-onnx-c-api.dll not found`，**与代码无关、不是回退**）⇒ CLI 那一面的读数**只能走 `-overlay` 台件**（样板：`176-r1` 的 `probes/176/r1/overlay-e2e.json` 把 `zz176r1_e2e_test.go` 编进 `cmd/wisp`，跟踪目录不多文件），别把"这个包红了"算成本票的账。⚠ **`./internal/risk/` 要单包跑**：本程实测四包并发时 `TestResolvePerCallBudget` 红（`1.214 ms/op`、1417 samples），同腿单跑 **PASS `0.512 ms/op`、2593 samples**——那枚尺自己吃 CPU 争用；`thresholds.go` 一字节不动，**不许当回退**，也不许为它放宽任何断言。⚠ **最终那一次读数必须在最后一枚 commit 之后取**（票 178 的教训：中途读数冒充终态）。⚠ **G3 陷阱**：不许给 `(*Loop)` 加收任何带 `taskID` 的**导出**方法，普查尺会把安静腿打红。⚠ `probes/161/r6/flip-declaration.sh` 别跑（跑一次就脏跟踪日志）。
- [ ] **AC#8 端到端解封——本票存在的唯一目的**：修完之后重跑 `probes/176/r1` 那发端到端，**后半截必须拿到读数**：模型续读宿主自己写下的指针，在真机 CLI 上**不被自家门拒**，并能逐字节读回那份副本。**这一发同时是三枚条件格的证据**：票 177 AC#3 的"端到端条件"、票 175 AC#5、票 176 AC#3/AC#4。
  ⚠ **分岔要报准**：如果新拿到的拒绝是 `…L2 级…已拒绝执行` 或 `已拒绝执行（风险 R4…）` 而不是"风险未分级"，那说明本枚修好了但**撞上另一枚缺陷**（豁免没生效／审批卡没接），**不许写成本票已结案**，要具名报回并指回票 177/162。

## 写面（超出即越权）

`internal/agent/loop.go`（**只 `decideRisk` 那一支＋它上面那段注释**）· 新判据件 `internal/agent/*_test.go`（一枚即可，命名带 `179`）· `cmd/wisp/run.go` 与 `internal/tools/bridge.go` 的**注释行**（仅 AC#6 那两处指针）· 本票面 Progress log（只追加、勾自己新写的格、留 AC#8 待编排者终判）· `docs/evidence/s1/179-declared-l0-refused-*.md`（交件表）· `.scratch/wisp/probes/179/**`（台件，`logdir` 取脚本自身目录）。
**禁改**：`PassThroughUnclassifiedRisk` 的语义/默认值/任何生产赋值点（AC#5）· `internal/risk/**` · `Bridge.Execute` 的裁决链 · `thresholds.go`／golden／审批超时常量／`allowlist.txt` 一字节 · `docs/PLAN.md` · `docs/specs/**` · `frontend/**`、`design/**`（别家归属，不读不写不引）· 别人的票面与证据件 · `probes/**` 既有台件（只读）。
⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）／`docs/evidence/s1/152-*.md` ⇒ **不提交、不还原、不补完、不评论**。
Git：只 commit 不 push；`git add -A`／`git add .`／`commit -a` 一律禁；一步式 `git commit -q -F - -- <显式路径>`；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout <分支>`／`switch`／`merge`／`worktree`／`clean`；临时件只建不删，还原用 `git cat-file blob HEAD:<路径> > <路径>`。

## 本票**不**解决

- 不裁"该不该给 L0 也上一张卡"——D4 已定案"不打断"，本票只让代码回到它。
- 不裁票 174 **AC#2b**：`cmd/wisp/run.go:362` 的 `tools.TaskDeps{Roster: rt.tasks}` 仍然**没接** `Paths` 字段（尺：`grep -n "TaskDeps{" cmd/wisp/run.go`），所以回执仍会带那句 "C26 未接线" 的 fail-closed 文案。这与本票正交：AC#8 那发读回靠的是**桥的豁免＋fs 根内真文件**，不靠那枚接线；两者别混勾。
- 不裁"为什么模型连续重试 7 次才被判定卡住"（那是票 162 那一族）。
- 不翻任何别人的勾，不 push。

## Progress log

- 09-28 09:5x 编排者立票：上面 8 把尺全部本程现跑（锚 `cf527f32`，`git status --porcelain -- internal/ cmd/` 空）。未派。

- 09-28 10:3x 写手程 `179-r1` 落第一／二枚 commit（`ec520750` 修法＋三枚常驻判据、`b6c34d34` AC#6 两枚指针），交件表 `docs/evidence/s1/179-declared-l0-refused-r1.md`。**AC 框一枚没勾**：AC#2/#3/#4/#5/#6 做完且测过，AC#1 没跑成（`176-r1` 台件在本机不带 PATH 就是 `0xc0000135`，我没重跑那发 7 次拒绝），AC#7 部分（三包 PASS＋gofumpt 空，d22scan／runtests／gate-clauses 三枚读数落在最后一枚 commit 之后取，逐字在交件回复），AC#8 未结案（真机 CLI 上 `task.output` 已执行、那句"风险未分级"已消失，但 `fs.read` 续读那一发断在**我自己台件**的 rig 判据形状——我用 content 文本找 call id，而 call id 在 `tool_call_id` 字段里；**没有**冒出 L2／R4／"查不到这个任务"）。两笔要报的不符：起手 HEAD＝`039efb47`（派单写 `cf527f32`，两者间只有文档件）；票面第 13 段"桥具的两枚测试工具恰好声明成 L0"行号属 `blockingProvider`，`newHarness` 默认那枚 `NewEchoProvider()`（`internal/agent/tools.go:123-131`）**根本不写 RiskLevel**，所以既有 loop 用例是靠"未分级＋开关开着"走过去，我的正向判据因此自建声明 L0 目录。工具调用 55/55，到顶停手。

- 09-28 11:2x 写手程 179-r2（台件腿，**零产码改动**）：台件落 .scratch/wisp/probes/179/r2/**（没动 r1 那枚），状态机只认工具行的 tool_call_id 字段，179-r1 断在"拿正文找 call id"的取数写法不复现——run B 第一次真把 fs.read 打在宿主自己写下的指针上。PATH 两形各一发：带前缀 rc=0（ok cmd/wisp [no tests to run]）、不带前缀 rc=1（exit status 0xc0000135），逐字在 logs/path-*.txt。到手：task.output 执行 risk=L0 decision=allow outcome=success，回执里「全文见 …artifacts\tool-output-agent-task-c66c0634….txt」（总长 20000 字节 / 约 5000 token）被模型逐字用作 fs.read 的 path，且与名册登记那条同一条。没到手：**逐字节读回**——那一发被 level=L2 rules=[R4] reason="R4: 包含来自 task.output 的内容" 拒（in_allowlist_scope=true＝不是越界升档），[确认 L2 fs.read] 窗口 0.0s 即拒，run B 退出 1、任务 cancelled，两发 fs.read 的工具行都是 0 字节（没回到上下文）；那两个数＝产物 20000 字节 / 读回 0 字节。四枚分岔三区计数：风险未分级 0/0/0、查不到这个任务 0/0/0、审批通道未接入那一形「L2 级」0/0/0（但升 L2 真实存在，尺不够宽已具名报：level=L2 2 枚、risk=L2 1 枚、确认 L2 1 枚）、R4＝console 8 枚＝**本发阻断者**。⇒ 照票面 AC#8 的分岔规矩**未结案**，具名指回票 177 shape A（Q-61甲 精确豁免）／175-r2（桥级盖章）：task.go:253-255 的 box.set 与 bridge.go:566 的 mark(...,hostPath) 在真 CLI 上没护住宿主自己的指针，同运行 C25 scope closed ... dropped=1 是第一条线索；另指回票 162（窗口 0.0s 即拒→cancelled 那一族）。票 174 AC#2b 那句"路径授权判定者未接线（C26）"照在回执里、正交未阻断。**AC 框一枚没勾**。自报一处形状：第一跑 runtime.Caller 在 -overlay 下报虚拟路径，读数误落 cmd/wisp/logs/，已原样 mv 成 logs/e2e-readings-run1-misplaced.txt（未删未改），logDir179r2 改为只认探针自身目录／上溯 go.mod，复跑两枚读数同为 28496 字节。门禁在 391822a2 之后自跑：三包 PASS（agent 1.887s／tools 15.444s／risk 单跑 4.575s）、d22scan rc=0 且 ban #8 internal/=431 与派单 10:4x 一致、d22scan runtests rc=0 PASS=34 FAIL=0、gate-clauses rc=1 红腿名册＝只有 腿=G6neg（无新增，全文 91056 字节落 logs/gate-clauses.txt）、gofumpt v0.12.0 对名下 .go 空 diff。工具调用 35/35，到顶停手。表 docs/evidence/s1/179-cli-reread-r2.md。
