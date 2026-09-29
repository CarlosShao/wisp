# 236-a1 — 票 236 六格的落点／载具／代价普查（只读腿，零产码）

- 起手锚点：**`82110540`**（自取 `git rev-parse --short HEAD`；票面写的起手 `c5d88a7f` 已被后续三推覆盖，本腿一律以 `82110540` 现读）。
- 起手名册：`.scratch/wisp/probes/236/a1/start-status-porcelain.txt`＝**151 行**（工作树本来就有别人的脏改动，终态判据＝等于这份名册，不是"为空"）。
- 本腿做了什么：**只读**。读文件 / `grep` / `git show|log|ls-files|cat-file` / `gh run view|gh api` / 一次 `gofumpt -l`（纯静态、不编译，编排者授权的例外；跑前先存起手名册）。
- ⛔ 未跑 `go test`／`go build`／`go vet`／任何编译（同时段 `235-r1` 独占产码与突变）。凡是"摘掉某行会红几枚"的判据，本腿**只给可复现的静态读数＋形状推理**，⛔ 不代跑突变。
- 禁区遵守：未读未引 `frontend/**`、`design/**`、`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件。未改任何跟踪文件、未 commit、未 push。
- 台件（每份 `wc -c` 非 0，见文末名册）：起手名册、`tracked-go-zlist.bin`、`gofumpt-tracked-stdout.txt`、`gofumpt-tracked-stderr.txt`、本文件。

---

## AC#1 — `taskCancel.Execute` 两支 fail-closed 无尺

### 落点
- 被测码（`internal/tools/task.go`，现读）：
  - `Roster == nil` 支＝**`task.go:699-701`**，字面量「任务名册未接线（fail-closed：拒绝停掉任何任务——名册才是唯一的停法）」。
  - `caller == ""` 支＝**`task.go:707-709`**，字面量「这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）」。`caller` 来自 `task.go:706` 的 `CorrelationID(ctx)`。
  - 同函数里紧邻的另外三支（`查不到`＝`:710-715`、`target == caller`＝`:716-723`、`ParentTaskID == ""`＝`:724-729`、`ParentTaskID != caller`＝`:730-735`）**都有尺**（`task_cancel_221_legs_test.go:395-447`、`:517-535`）。⇒ 缺口精确地就是票面那两支。
- 建议判据落点（具名，⛔ 不新造载具）：**追加两枚用例到既有文件 `internal/tools/task_cancel_221_legs_test.go`**（同包同文件，不动 tracked 分母）。建议名 `Test221TaskCancelWithoutRosterRefuses`／`Test221TaskCancelWithoutHostMintedCallerRefuses`。两枚都必须**钉字面量**（钉上面那两句原文），因为"仍回 IsError"这件事删支之后照样成立（会掉进 `查不到` 或 `不是调用者的孩子` 那一支）——只断 `IsError` 的尺＝没有牙。

### 载具：**两支都有现成载具，零新 infra**（这是本格最要紧的一条）
- `caller == ""` 支：`build221(t, windowGate221(), provider, false)`（`task_cancel_221_legs_test.go:74`）＋现成 helper `x.cancel(t, "", target)`（`:116-127`；它把 callerID 同时写进 `TaskID` 与 `CorrelationID`）。真桥**不会自造 id**：`bridge.go:246-247` 只做 `CorrelationID==""→TaskID` 的回落，`bridge.go:457` 落 `orDefault(req.CorrelationID, req.TaskID)`（`bridge.go:1181-1186` 的 `orDefault` 对空串返回默认值，两枚皆空⇒空），于是 `CorrelationID(ctx)` 按 `cancel.go:55-58` 返回 `""`。L1 窗口用 `windowGate221()`（`:62-69`）答 `AnswerTimeout`，`bridge.go:381-387` 逐字"running out MEANS EXECUTE"⇒ 进到 Execute。`assessorFor("")`（`bridge.go:852-857`）在 `taskID==""` 时返回不带 taint 的 `b.assess`，不会 panic。
  - 目标行现成可用：`x.h.roster.PublishSubagent("kid-221", parent197, …)` 这副写法已在 `:521` 用过。
- `Roster == nil` 支：现成 helper **`mustRegisterTaskEntries(t, TaskDeps{Roster: nil})`（`task_output_leg_test.go:141-150`）** 注册的正是 `BuiltinTaskEntries` 全家（含 `task.cancel`，`task.go:600-605`），配 `New(Options{… Gate: windowGate221()})` 即可。**同形先例已在仓**：`task_output_leg_test.go:113` `TestMissingArgsAndUnwiredRoster` 的 `:126-138` 用完全一样的形状钉了 `taskOutput` 那一支（`task.go:514-516`）⇒ 本格＝把它抄到 cancel 的调用名与字面量上，不是新发明。
- ⛔ 注意 `build221` 复用不上的那一半：`task_cancel_221_legs_test.go:101` 把 roster 写死成 `TaskDeps{Roster: x.h.roster}`，所以 **nil-roster 不能走 `build221`**；走 `mustRegisterTaskEntries` 那条（或给 `build221` 加一个 deps 参数＝改动既有 8 处调用者，代价更大，不推荐）。

### 同包另一副 harness 的可用性（编排者点名要读的 `subagent_222_test.go`，现读）
- `newH222`（`internal/tools/subagent_222_test.go:251-291`）的桥是 **`Gate: NoGate{}`（`:271`）**，且 `wireSpawn`（`:296-322`）只注册 `BuiltinSubagentEntries`（`:317`），**根本不注册 task 家族**、roster 恒非 nil（`:254`）。⇒ **h222 不能当 AC#1 的载具**：`task.cancel` 是 L1，NoGate 的 `PendingWindow` 答 `AnswerReject`，调用进不了 Execute——这正是 `task_cancel_221_legs_test.go:13-19` 那段头部注释逐字警告过的形状（"a refusal from the gate is not a refusal from the authority check"）。
- ⚠ **在飞冲突提示（写腿必读）**：本腿运行期间 `235-r1` 落了 3 枚提交（`82110540`→`a71f0be9`），其中 `internal/tools/subagent_222_test.go` 改 93 行、`internal/tools/subagent_197_test.go` 改 18 行（后者只动 `:398` 那段注释⇒本腿引的 `subagent_197_test.go:201` 未漂移）。222 里新占用到的 helper 名：`awaitChildOnBridge222`（`:391`）、`drainBridgeArrivals`（`:413`）⇒ AC#1 的新用例别与它们重名，也别指望 `h222` 那套桥位计数工具可复用。

### 今天确实零尺（本腿复算的读数）
- `grep -rn "任务名册未接线" --include=*.go .`＝**6 命中，全在源码／注释**：`subagent_197.go:254`、`task.go:515`、`task.go:700`、`task.go:828`（注释）、`task_backfill.go:108`；**`_test.go` 0 命中**。
- `grep -rn "没有宿主给的任务 id" --include=*.go .`＝**2 命中**：`subagent_197.go:261`、`task.go:708`；**`_test.go` 0 命中**。
- `grep -rn "x.cancel(t, \"\"" --include=*_test.go internal/tools`＝**0 命中**（rc=1）；`cancel(t,` 现有 8 处调用者（`:211/:302/:405/:415/:432/:436/:506/:522`）传的全是 `parent197`／`aID`／`bID`／`"kid-221"`，非空。
- nil-roster 注册今天**有一处但从不派发**：`ticket175r2_stamp_live_test.go:353` 用 `BuiltinTaskEntries(TaskDeps{})` 只为枚举名字（该用例 `:349` 起，通体不 Execute）⇒ 分支进不去。
- ⇒ 票面"摘掉后各 0 枚红"与静态事实**一致**（本腿不代跑突变，⛔ 235-r1 在飞）。

### 代价
- 谁会被打红：纯新增用例，不删不改任何现有断言⇒在册红名册不受影响。两枚新用例进 `internal/tools`（今晚该包在 CI 是失败包之一，见 AC#5），写腿交回时必须**逐名证明这两枚为绿**，否则就是"新增红"。
- 门的分母：**追加到既有测试文件⇒tracked `.go` 分母不变**；新开文件⇒分母 +1，且立刻被 `ci.yml:136` 与 `:152` 那两步重看（AC#6 说明那两步今天的颜色）。
- 契约邻接：**不属**。不注册 `task.list`（D34 名册不动）、`Declared` 仍是 `risk.L1`（`task.go:658-667` 不动）、不碰 SLO/golden。
- 争用：与在飞的 `235-r1` 同包（`internal/tools`）——派单时要错开或说明。

### 顺带顶出的一格（超出票面射程，交编排者定）
同族还有 **3 枚同样零尺**的 fail-closed 支，在 `task.spawn` 侧：`internal/tools/subagent_197.go:253-255`（`Roster == nil`）、`:256-258`（`BaseOptions/ParentTools == nil`）、`:259-263`（`parentID == ""`）。尺与读数同上（两条字面量 `_test.go` 0 命中）。要不要与 AC#1 同批收，是**射程决定**，不是本腿能定的事；只把"它们存在且同样无尺"摊平。

---

## AC#2 — 那把 DEFERRED 尺的 (b) 支无牙

### 机制（本腿现读复认，并补一层）
- 尺＝**`internal/tools/task_cancel_221_legs_test.go:223-247`** `Test221DeferredMarkerForCancelLiftedButListStillMarked`：
  - `:224` `os.ReadFile("task.go")`；`:229-239` 逐行要求**同一行同时含** `DEFERRED` 与 `task.list`；`:244-246` 只在 `listMarked == 0` 时报。
- 为什么撑得住：`grep -n "DEFERRED" internal/tools/task.go`＝**3 行**（`:23`、`:32`、`:279`），其中同时含 `task.list` 的是 **2 行**（`:23` 名册抬头、`:279` `Count()` 的叙述注释）。摘掉 `:23`⇒`listMarked` 仍是 1⇒PASS。**票面机制成立，本腿复算通过。**
- ⚠ 本腿补的第二层（比"词面被叙述句撑住"更硬）：`:224` 是**读盘型尺**——`go test -overlay` 只替换编译器看到的字节，测试二进制运行时 `os.ReadFile` 打的是**物理盘**。⇒ 对 `task.go` 做任何 overlay 突变，这把尺**结构上看不见**，今天连"测它的牙"都做不到（不是不想测，是测不出）。

### 改扫能力之后放哪（载具现状，具名）
`BuiltinTaskEntries` 的名册枚举今天被这些测试在用（尺＝`grep -rn "BuiltinTaskEntries" --include=*.go .`＝17 命中，其中测试 8 处）：
1. `internal/tools/task_cancel_221_legs_test.go:171`（`Test221TaskCancelIsRegisteredAtItsFrozenLevel`）——`:183-185` 是 **`len(names) != 2` 计数形**，报错文案里提 `task.list`，但判据不是名集。
2. `internal/tools/task_cancel_221_test.go:42-69` `allBuiltinEntriesHere(t)`（四家 `Builtin*Entries` 的并集）＋ `:130-132` 自证分母非空＋`:134-144` "说明书许诺 vs 注册并集" ⇒ **这已经是能力尺**（摘掉 `task.cancel` 注册它就响）。
3. `internal/tools/ticket175r2_stamp_live_test.go:349` `TestEveryRegisteredToolIsClassifiedForMarking175r2`＋`:353` 用 `BuiltinTaskEntries(TaskDeps{})` 枚举，名册里每枚必须在 `classified175r2`（`:333-347`，含 `"task.cancel"` 那一行在 `:345`）有分类，否则红 ⇒ **新增注册必被它拦**（有人把 `task.list` 注册进来而没答分类表，它响）。
4. 其余为装配用途，不判名集：`task_output_leg_test.go:40/:144`、`task_output_pointer_notice_test.go:53`、`subagent_197_test.go:201`、`pointer_183_cli_seam_test.go:111`、`pointer_185_cli_seam_test.go:103`、生产装配 `cmd/wisp/run.go:441`。
- ⇒ **"task.list 不在名册里"的名集形判据今天零载具**。尺＝`grep -rn "task\.list" --include=*_test.go .`＝**4 命中**，全在 `task_cancel_221_legs_test.go`（`:184` 计数文案、`:222` 注释、`:236` 词面、`:245` 报错文案），没有一枚写 `registered["task.list"]`。落点建议：在 (1) 那枚里把计数换成名集（或新增一枚 `Test221TaskListStaysOutOfRoster`，读 `allBuiltinEntriesHere`／`BuiltinTaskEntries` 的名字集），词面腿可保留但要就地标注它测不到能力。

### 这条判据将来怎么测它的牙（⛔ 不新造门的批准）
- 形状 (i)（能力形，本腿认为唯一能自证的形状，**不是裁定**）：判据读**编译后的注册表** ⇒ overlay `internal/tools/task.go` 的 `BuiltinTaskEntries`（`task.go:600-605`）加/删一枚 `Entry` 就是**可见突变**。牙＝两发：M+（把 `task.list` 塞进名册）⇒ 该枚具名用例红；还原⇒绿。同时它会被 (2)(3) 两把既有能力尺**一起拦**⇒ 单点摘标记不会静默。
- 形状 (ii)（保留词面腿时的唯一可测法）：overlay 只能替换**测试文件自己**——把 `:224` 的读路径改成 `t.TempDir()` 副本、并把副本路径一起 overlay 进去，配正控（副本里删掉 `:23` 那行⇒该发必须红）。代价：突变面不在产品码上，**测的是尺不是码**，很容易被读成"自证"；交回时必须写成"词面腿的牙只能用这种方式量，且它不等价于能力判据"。
- 与票 225 分开算账（现读确认）：`grep -rn "DEFERRED" tools/d22scan/*.go`＝**0 命中**；`scripts/portable-tests.sh:50/:95` 只是注释里的 DEFERRED 字样。⇒ d22scan 里**没有** DEFERRED↔`SPEC-12 §5` 的双向对账尺，那件事仍是票 225 的；本格只管"这把尺自己有没有牙"。

### 代价
- 谁会被打红：把 (1) 的计数换成名集后，**注册任何第三枚 task 工具**不再被误判（计数形会误判合法新增为名集违规——那是它现在的假红风险，本格顺手收窄，不算放宽）。⚠ 但 `len(names) != 2` 那枚若被删＝动既有断言，⛔ 本票禁止；建议**保留计数、另加名集**。
- 门的分母：不动。契约邻接：**不属**（D34 名册内容不变，只是判它的方式从词面改能力）。

---

## AC#3 — "任务 id 必须由宿主铸造"零钉

### 落点
- 铸造处＝`internal/agent/loop.go:1121-1136` `newTaskID()`。
  - **熵源失败退化支＝`loop.go:1125-1131`**：`if _, err := rand.Read(b[:]); err != nil { n := time.Now().UnixNano(); for i := range b { b[i] = byte(n >> (8 * uint(i%8))) } }`。形状后果：16 字节里只有低 8 字节随 ns 变、高 8 字节恒为同一 ns 的重复填充，且随后仍被 `:1132-1133` 盖上 UUIDv4 版本/变体位 ⇒ **输出仍是合法的 36 字符 uuid 形状**，但同一 ns 内两次调用可撞车（`hex`＋`-` 分段不变）。`rand` 是 `loop.go:5` 直引 `crypto/rand`，**没有任何 var 间接层**（现读 `grep -n "crypto/rand\|rand\." internal/agent/loop.go`＝3 命中：`:5` 引包、`:1122` 注释、`:1125` 唯一调用）。
  - 调用者只有 2 枚：`loop.go:322`（`RunAsync`，同一函数 `:324` 用 `"agent-task-"+id` 进 D38 名册）与 `loop.go:333`（`Run`）。
- 名册键的生产来源＝`cmd/wisp/run.go:892`：`…}).Backfill(bg.Root().Ctx, res.TaskID, res.Text)` ⇒ 键是环路铸的；`TaskRoster.Record`（`task.go:215-218`）对键的内容**零假设**（只挡空串）。
- ⚠ 写侧还有一枚同族零尺 fail-closed＝**`internal/tools/task_backfill.go:109-110`**（`taskID == ""` ⇒「这个任务没有宿主铸的 id（id 只由环路生成，组合根不许自造一枚）」）。尺＝`grep -rn "没有宿主铸的 id" --include=*.go .`＝**1 命中（源码自己），`_test.go` 0 命中**。

### 载具（⚠ 本格要顶正票面一格：不是"零钉"，是"钉错了维"）
- **"裸 task id 会静默换字节"那一风险——今天有牙，具名**：`internal/tools/ticket176r1_start_port_test.go:128-168` **`TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1`**，两半都钉 `taskArtifactPrefix`：
  - `:144-146`：产物文件 basename 必须含 `agent-task-<taskID>`，否则 `t.Errorf("产物名 %q 不带 agent-task- 前缀…")`（`:141-143` 那句注释自己点名"this prefix is the only thing separating the two key spaces"）；
  - `:152-167`：用"模型可以重复的裸 task id 当 call id"真的撞一次 `spills.Prepare(taskID, foreign)`（`:153`），要求产物路径不同（`:157-159` `t.Fatalf`）＋任务那份字节没被换（`:160-167` 逐字节比）。
  - ⇒ 摘掉 `task_backfill.go:130` 的 `taskArtifactPrefix+`（票面点名的 `:34-42` 风险的那处落点）**这枚用例必红，且两半都红**。前缀生效性**不是零载具**。
- **"谁铸造"那一维今天确实零钉**（这部分票面对）：`grep -rn newTaskID --include=*_test.go .`＝**5 命中**，逐条读：`internal/agent/compress_trace_test.go:524-527` 是**形状**检查（`len(task) != 36 || strings.Count(task, "-") != 4`，射程是 trace 的 `"task"` 属性那一族；⛔ 别当保护，模型给一枚 36 字符 4 横杠的 id 它一样绿）、`:557-559` 只断"两枚不同"、`internal/tools/pointer_183_cli_seam_test.go:59` 是注释不是断言。⇒ **没有任何尺会在"roster 键／Backfill 的 taskID 换成模型给的值"时变红**（`Record`/`Look`/`Cancel` 对键的内容零假设，来源、熵、可否被参数复写都没人查）。
- 熵源失败那一支的载具：**今天零**。尺：`grep -rn "UnixNano" --include=*_test.go internal/agent | grep -i taskid`＝0；根因＝没有注入缝（见上 `:5`/`:1125` 的直引）。

### 三条修法与各自代价
1. **性质钉**（零产品改动）：同进程连铸 N 枚 ⇒ 两两不同 ＋ 形状 36/4。代价：测不到退化支本身（退化支在多数机器上根本不触发），⛔ 交回时必须写明"覆盖不到熵失败"，别报成全覆盖；否则是"用假绿换一条勾"。
2. **加注入缝再突变**（`var randRead = rand.Read` 之类）：能让 `:1125-1131` 可测。代价：动产品码形状，且**新开一枚"测试可注入失败"的面**——AGENTS §1.3 只允许在既有接缝注入，新接缝本身得先有自己的判据，否则正好是"为测试开后门"那一形。⚠ 属〔形状改动，非本腿可批〕。
3. **来源钉**（本格真正的缺口）：在 `Backfill`／`Record` 边界断"键必须来自 `RunningTask.ID`"——可做的最轻形状＝在 `cmd/wisp` 那侧把 `res.TaskID` 与 `bg.ID` 相等钉住（现读 `run.go:892` 用的是 `res.TaskID`，而 `RunningTask.ID` 在 `loop.go:323` 由 `newTaskID()` 填 ⇒ 两者一致这件事**今天没有尺**）。代价：需要一发跑到真组合根的 CLI 接缝用例（现成载具候选：`internal/tools/pointer_183_cli_seam_test.go`／`pointer_185_cli_seam_test.go` 那两枚 CLI-seam 形状，或 `cmd/wisp` 的 `run` 测试族）；⛔ 不许把 `allowed_dirs` 做成硬边界、不许新造工具面（票面禁区，本腿也没往这个方向摊）。

### 名册内视——三行读法（照编排者要求，⛔ 不加重定性）
1. **现象在哪出现**：`task_backfill.go:34-45` 注释自己写明"普通 spill 路径的 call id 是 model-supplied、裸 task id 是完全可以被够到的 call id"。风险面是**产物文件名的命名空间**，不是权限面；名册键今天只由 `run.go:892` 的 `res.TaskID`（环路铸）填。
2. **有没有本机被入侵的证据**：**没有**。本腿全程只读，未见任何越权、外泄或被利用迹象。这一格是"没人能证明它没坏"，⛔ 不是"它坏了"。
3. **最坏后果是什么形状**：如果前缀那一支哪天被摘（或有人把模型给的 id 直接当 roster 键），最坏＝**同一产物目录上 last-writer-wins 的静默换字节**（注释引 `spill.go:117-131` 不报错）＋ C25 provenance 失真，回执仍"成功"。**不是提权、不是任意文件写**：`task.*` 两枚的 `PathParams` 都是空（`task.go:591`、`task.go:663`），唯一的模型输入是名册键，C26 在这条路上没有可判路径。

---

## AC#4 — `ci.yml:519-524` 那句过期注释

### 现读原文（`.github/workflows/ci.yml`）
- 步名＝**`:513`** `"PathResolver junction placeholder (real cases tickets 18/20)"`；run 行＝**`:527`** `bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows`。
- 那 6 行（`:519-524`）逐字：
  `:519` `# DO NOT READ THIS WIRING AS A FIX OR A MASK: the step FAILS on`
  `:520` `# windows-latest today for its own real reason (run 35551819606 job`
  `:521` `# 106188167785: USERPROFILE's 8.3 form RUNNER~1 makes the A-tier anchor`
  `:522` `# classify B, want ClassA), and runtests.sh propagates a non-zero `go``
  `:523` `# test` exit code unchanged (registry A43② documents the same discipline`
  `:524` `# for the vet failure that hid this gate).`
  （`:514-518` 是同一注释块的前 5 行，讲"这步不是 mask"，与过期无关。）

### 读数（本腿自己现读台件；⛔ 那 18 份日志未入库，本地读得到）
- 今晚一发 `run-36559498617`（test-windows job 正文）：`.scratch/wisp/probes/ci-red/logs/run-36559498617-testwindows-job.txt:4016` 起是**这一步自己的 group**，`:4024` `--- PASS: TestPathResolverJunctionWindows (0.03s)`，`:4025-4026` `PASS` + `ok github.com/CarlosShao/wisp/internal/risk 0.067s`，`:4027` `runtests.sh: OK - packages=[./internal/risk/ -run TestPathResolverJunctionWindows] top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0`。⇒ **这一步今晚为绿**，`:519-520` 那句"the step FAILS … today"对**这一步**为假。
- 今晚另一发 `run-36556778053`：`run-36556778053-log-failed.txt:4044-4047` 同一读数（`--- PASS` ＋ `runtests.sh: OK … PASS=1 FAIL=0`），但**步名列整列是 `UNKNOWN STEP`**（归因只能靠正文，见 AC#5）。
- ⚠ **8.3 短名那一因没有消失，只是不砸在这一步上**：同一发里 `TestPathResolverShortNameAListDenied`（`:3405`）、`TestPathResolverUNCAListDenied`（`:3408`）、`TestPathResolverExtendedLengthPrefixAListDenied`（`:3411`）三枚 FAIL，失败正文正是 `want "C:\\Users\\RUNNER~1\\…"`（`:3404`／`:3407`／`:3410`）；`RUNNER~1` 串出现次数（现量 `grep -c`）：`run-36559498617-testwindows-job.txt`＝**612**、`run-36556778053-log-failed.txt`＝**612**、`run-36559498617-log-failed.txt`＝**212**。

### 判断：**改，但不是"过期"，是"指错对象"**
- 最小有害改法：只动 `:519-524` 的**文字**——把"the step FAILS … today"改成带 run 号＋日期的读数（并注明今晚两发该步绿），并补一句"8.3 那一因今天的落点是同包里那三枚姊妹用例，不是本步"。⛔ 不碰步名、`if:`、`run:`、顺序、不加 `continue-on-error`。
- 风险（为什么这一格仍要小心）：动 `.github/workflows/ci.yml` 哪怕只改注释，会被读成**动门的形状**——本文件自己的纪律先例逐字写在 `ci.yml:146-151`（"no existing step was moved, edited, deleted, given an `if:`, or made continue-on-error"）。注释虽不是步，但它长在门的正文里；且 `ci.yml` 的 diff 会触发评审者对"门的分母/颜色"的整套警觉。票面 AC#4 自己也写了"低利害、可逆、不该拿去找 owner 拍"⇒ 编排者可自裁，⛔ 本腿不代裁。
- 若"留"：代价已经被验证了一次——票面现量表第 4 行就把红的步当成了 tracked 分母那一步（实际红的是 `ci.yml:136`，见 AC#6）。**过期注释的代价不是文案，是把归因方向带偏。**

---

## AC#5 — "逐名比红名集合"要自带隔离复量那一半

### 这条尺今天写在哪儿（现读，答案：只写在停车点散文里）
- 正尺形态写在 **`docs/reports/HANDOVER.md:331`**：「正尺＝`grep -cE '^[[:space:]]*--- (PASS|FAIL)'` ＋拿 `=== RUN` 当分母对拉」（同一行还记着"我那把逐名尺少算 37% 的用例"：`^--- PASS`=295／缩进 `    --- PASS`=174／`=== RUN`=471）。
- 规矩本体（"整包复跑＋逐名比红名集合"）写在 **HANDOVER 的时间戳散文**里：`:103`、`:118`、`:131`、`:260`、`:277`、`:305`、`:357`。
- **派单模板 0 命中**：`grep -cE "红名|隔离复量|逐名" .scratch/wisp/issues/README.md`＝**0**；`scripts/` 与 `.github/` 亦 0 命中（同尺）。⇒ 这条尺**只活在编排者记忆里＋停车点散文**，正是本仓反复栽的那一类（`:305`／`:357` 两次都是我自己写进账的数字被非实现者腿推翻）。
- ⚠⚠ **本腿新量到的一格（会误导所有拿 four numbers 互咬的人）**：`scripts/portable-tests.sh:426-430` 的四把计数就是 `count '^=== RUN'`／`'^--- PASS'`／`'^--- FAIL'`／`'^--- SKIP'`——**行首锚定**，结构上看不到缩进子测试（`:432` 打印那行自称 "all from -v output"）。⇒ four numbers 与"名级 25 行（顶层 23）"**只能对拉顶层口径**；拿它当"子测试层也数过了"就是假绿。`scripts/winsec-tests.sh:134` 同形。

### 现量（名册口径必须先声明，否则名册自己会换位）
- 今晚两发红名册**逐字相同**：`diff <(grep '^Test' names-run-36556778053… ) <(…36559498617…)`＝**IDENTICAL**，各 **25 行名**（其中 2 行是缩进子测试 `TestL1VetoNeedsAChannelTheHostReallyWired/…`⇒顶层 23）。
- four numbers 互咬成立：`names-run-36559498617…` 末 3 行 FAIL＝4／7／12⇒相加 **23**＝顶层数；`winsec-tests.sh` 那发 FAIL=0。
- 与本机基线**差 8 枚**：`roster-baseline-f7478d37.txt`＝**17 枚**，今晚 25 ⇒ 新增 `TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestL1VetoNeedsAChannelTheHostReallyWired`（＋2 枚子测试）、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestRunPacketCarriesTheLoadedInstructionFiles`、`TestTicket223PermissionDeniedSitsInItsOwnSentence`。⇒ 这两份的**跑范围口径不同**（一份是本机包集、一份是 CI 全 job 的 `--log-failed`），不声明就会再撞一次"11 vs 23"。
- 步名归属实测：`gh api repos/CarlosShao/wisp/actions/jobs/109376664767 --jq '.steps[]'` 能给出准确步名（lint job 12 步有名有结论）；而 `--log-failed` 的**日志列**在 `run-36556778053`／`run-36500353477` 上是 `UNKNOWN STEP`（`run-36559498617` 同发却给得出）。⇒ 规矩要写成"步名归因前先确认这一步名给不给得出；给不出就用正文里的自报行（`runtests.sh: OK - packages=…`／`portable-tests.sh: four numbers`）认步"。
- 隔离复量的现成先例（⛔ 本腿不重跑）：`docs/evidence/s1/221-task-cancel-v1.md:165`（`cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues` 整包 5.15s 红，出处 `config_reload_223_test.go:379`）＋ `:171`（单包隔离复量 **3/3 全绿** 5.055／5.019／5.034s，日志 `recheck223-1..3.txt`）＝计时红不记账；`:22` 与 `:272` 就是"在册名册不是稳定集合、两发会换位"的原始读数（那发是 `cmd/wisp` 红、另一发是 `risk/R-116-1` 红）；`:175`／`:179` 是这条更正的原文两句。本机另一形同证据＝`docs/reports/HANDOVER.md:212`、`:224`（`TestResolvePerCallBudget` 整包红、`-count=3` 安静复量转绿，同日第三、四例争用型假红）。

### 应落位置候选 ＋ 各自会误导谁
1. **`.scratch/wisp/issues/README.md`（派单规矩节）**：继承性最好（每枚腿起手必读、AGENTS §1.4 逐字指回它）。⚠ 会误导：它是 **git／派单纪律**的权威文本，把"取数口径"混进去，读的人容易当成"纪律"而不是"怎么采数"，也可能被当成可选注脚（因为纪律条目里夹着工具说明）。
2. **`scripts/portable-tests.sh` 头部注释（`:20` 附近，紧挨它自己那句 "the four numbers printed below"）**：离读数最近，最可能自动生效。⚠ 会误导：写腿看到注释容易以为**脚本已经做了隔离复量**（它只打印 four numbers，且 `:426-430` 是顶层锚定）⇒ 必须同批写明"脚本不复量、不数缩进子测试"。
3. **新开一枚仪器说明（例如 `docs/reports/` 下一份）**：⚠ 会误导：本仓有"索引写了但文件不存在"的前科（AGENTS.md §4 末尾那段 ⛔ 就是为此写的），新开一份最容易被漏读；真要新开，必须同时进 `docs/reports/pending-and-issues.md` 台账并被 AGENTS 索引指向，否则等于第二份没人读的文档。
4. 现状最低成本路线＝**HANDOVER:331 那条正尺升级成两条**（整包名册红 ＋ 单包隔离复量 ＋ 口径声明），并在台账追加一条 A## 指向它。⛔ 落点由执行者／编排者定，本腿只摊"各自误导谁"。

---

## AC#6 — tracked 分母被入库夹具污染

### 全仓现量（本腿自己跑的静态尺，不是转述）
- 尺：`git ls-files -z '*.go' | xargs -0 "$(go env GOPATH)/bin/gofumpt.exe" -l`；分母＝**726 枚 tracked `.go`**（其中 `.scratch/**/*.go`＝**171 枚**）。落盘：`gofumpt-tracked-stdout.txt`／`gofumpt-tracked-stderr.txt`／`tracked-go-zlist.bin`。
- **stdout＝7 枚需要格式化（全部在 `.scratch/wisp/probes/**`）**：
  1. `.scratch/wisp/probes/163/a1/main.go`（票 163 只读核台件的正经探针程序，3,927 字节）
  2. `.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go`（overlay 夹具，5,353 字节）
  3. `.scratch/wisp/probes/183/r2/mut/task-boxset-off.go`（`mut/` 突变样本）
  4. `.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go`（`mut-m1/` 突变样本）
  5. `.scratch/wisp/probes/197/r1c/pre/subagent_197.go`（**产品码改前快照**，19,314 字节）
  6. `.scratch/wisp/probes/197/r1c/pre/subagent_197_test.go`（同上）
  7. `.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go`（`mutations/` 样本）
- **stderr＝1 枚解析错误**：`.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations`（xargs rc=123）。⇒ 现量总数＝**8 枚脏样本（7 可列 ＋ 1 让尺退出 2）**。
- ⚠ 与台账的差：`docs/reports/HANDOVER.md:397` 与 `docs/reports/pending-and-issues.md:8983` 写的是"名单确实只有 **4** 枚 ＋ 第 5 枚错误行"。现量＝7 ＋ 1，**多出的 3 枚是 09-29 之后新入库的**（`197/r1c/pre/` 两枚、`222/v1/mutations/` 一枚）。⇒ 那两句要么标过期，要么按台账规矩**追加**更正（⛔ 不删原句）。引用它们做决策的下一程会低估分母。

### 今晚 CI 到底哪一步红（`gh api` 实测，⛔ 不是推测）
`gh api repos/CarlosShao/wisp/actions/jobs/109376664767`（run `36559498617` 的 lint job）步级结论：
| # | 步名 | 结论 |
|---|---|---|
| 4 | D22 scanner positive control | success |
| 5 | D22 scanner self-test (161 AC#2) | success |
| 6 | D22 seven-ban + emoji scan | success |
| **7** | **gofmt (gofumpt)** | **failure**（正文：`##[error].scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations` 然后 `##[error]Process completed with exit code 2`，即 `run-36556778053-log-failed.txt:7080-7081` 同形） |
| **8** | **gofmt (gofumpt) - the tracked set is the denominator (161 AC#7 form A)** | **skipped**（＝今晚**根本没跑**） |
| **9** | **go vet (module)** | **skipped** |
| **10** | **go vet (tools/d22scan module)** | **skipped** |
| 11 | staticcheck | failure |
| 12 | mockllm module vet | success |

⇒ **票面现量表第 4 行要更正一处**："它进了 tracked 分母那一步 ⇒ `lint::gofmt (tracked set)` 常红"——**红的是 `ci.yml:136-144` 那一步（工作树遍历 `gofumpt -l . tools/d22scan tools/mockllm`）**，`ci.yml:152`/`:172` 那一步（tracked 分母）今晚是 `skipped`，从没给出读数。"吃掉两道 `go vet`"属实且现在**有了 step conclusion 级实据**（`skipped`×2）。这也正是 `ci.yml:153-156`（同一注释块 `:146-171`，票面与台账称"那 20 行注释"）自己承认的"adds NO REACH TODAY… the step above and this one read the same bytes"——两步都被同一枚夹具拦住，但只有先跑的那步留下红。

### 甲（只排 import 顺序、保留符号坏点）
- 被改对象现读：`.scratch/wisp/probes/185/c1/mut/fs_broken.go`＝**11,041 字节／338 行**；形状＝`package tools` / `var probe185c1BrokenMarker = undefinedSymbol185c1` 排在 `import (…)` **之前**（所以 gofumpt／编译器报的是 4:1 的 import 位置错）。甲＝把那三行挪到 import 块之后⇒**符号坏点 `undefinedSymbol185c1` 保留**、解析错误消失、字节总数与行数不变、**内容顺序变⇒文件 md5 变**。
- **票 185 名下引用查没查到：查到了，但不是字节数/md5**（尺：`grep -rn "fs_broken" docs .scratch/wisp/issues`＝**13 命中**；再叠 `grep -i "md5\|字节\|bytes\|wc -c\|行"`＝下面这些**行列／错误文案**引用，**0 枚字节数或 md5 引用**）：
  - `.scratch/wisp/probes/185/c1/logs/overlay-proof-broken.txt:2` ＝ `…fs_broken.go:4:1: syntax error: imports must appear before other declarations`（**具名行列**）。
  - `.scratch/wisp/probes/185/c1/overlay-proof.json:1` ＝ 把 `internal/tools/fs.go` 映射到该枚文件的 overlay 名册（路径引用，无校验和）。
  - `docs/evidence/s1/185-reread-owner-census-c1.md:48` ＝ 「仪器自证（防恒绿）」那一格：**正引这条 rc=1 与那句 syntax error** 来证明 overlay 映射生效（"零枚红是读数不是没跑"）；同文件 `:185`、`:191` 点名该件"保留不删"。
  - `docs/evidence/s1/197-subagent-entity-r1b.md:105`（表行：`4813567e` 09-28 14:42 ＋ 那句 `4:1`）、`:288`；`docs/evidence/s1/200-project-instructions-r2.md:153`；`docs/reports/HANDOVER.md:360`、`:397`；`docs/reports/pending-and-issues.md:8983`、`:9566`；`.scratch/wisp/probes/ci-red/ci-red-1.md:112`。
- ⇒ **甲的真实代价被票面写歪了一点**：不会"洗掉字节数引用"（没人引字节数），会**改掉三处以上被引用的那一行错误文本的语义**（`4:1` 之后变成 `undefined: undefinedSymbol185c1`，`go build -overlay` 的 rc=1 仍在但**理由换句**）。185 那格的"防恒绿自证"正是靠那句 error 文本。动之前必须**在台账追加一条**（例：「自 `<新锚>` 起 `185/c1/mut/fs_broken.go` 的 overlay 生效证明错误行由 `4:1 imports must appear before other declarations` 变为 `undefined: undefinedSymbol185c1`，rc=1 不变，原判据不受影响」），否则后程复算 185 会判"凭据对不上"。⛔ 不许删件（`issues/README` 规则 8／票面禁区）。
- 附带代价：甲修的是**内容**，`ci.yml:136` 那步的 stdout 清单里**还剩 7 枚**（见现量名册），所以甲**不能让 lint 变绿**，只能让它从"尺退出 2、后续步全 skip"变成"尺退出 1、打印 7 行清单"⇒ 后面两道 `go vet` 依然被同一 `bash -e` 的 `exit 1` 吃掉。**⇒ 甲单独做完拿不到 `go vet` 读数**，与票面"完成判据＝任选一支之后必须真拿到一次 `go vet` 读数"直接冲突。这一条是本腿认为最要紧的代价信息：**任何只动 `fs_broken.go` 一枚的修法都完不成该判据**。

### 乙（分母排除 `mut/` 这类刻意夹具目录）——⛔〔契约邻接〕，本腿只摊代价、不推荐执行
- **动的是哪枚文件：既不是 `ci.yml`、也不是 `scripts/`**，而是 **`.scratch/wisp/probes/161/r5/attrib.sh:341`**（`A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1)`；`:337` 先数分母、`:342-345` rc>1 直接 `exit 2`）。`ci.yml:152` 那步只是 `sh` 它（`ci.yml:172`）。
- ⚠ 但**今晚的红不在那一步**（见上表 step 7 vs step 8）⇒ 只改 `attrib.sh` **修不到颜色**；要修 step 7 就得动 `ci.yml:137-144` 的正文（把 `.` 换成排除式或改成喂 tracked 名册）＝**动门的形状**，正是 `ci.yml:146-151` 那段先例禁止"移动/编辑/加 `if:`／`continue-on-error` 既有步"的射程，须人工批准。
- ⚠ **排除会漏掉真的该红的东西吗？会，而且现量给了具名样本**：`mut/` 一条规则只盖得住 8 枚里的 **2 枚**（`183/r2/mut/task-boxset-off.go`、`185/c1/mut/fs_broken.go`）；`mut-m1/`、`mutations/m11-…`、`pre/`（两枚）以及 `163/a1/main.go`、`174/c2/zz…_test.go` **都盖不住**。若放宽到"整片剔 `.scratch/**`"⇒ tracked 分母从 **726 掉到 555**（171 枚免检），其中 `197/r1c/pre/subagent_197.go` 与 `…_test.go` 是**产品码 `internal/tools/subagent_197.go` 的改前快照**、`174/c2/zz…_test.go` 是**会被 overlay 编进 `package tools` 的夹具**⇒ 这两类一旦免检，"格式门"对它们的真实形状从此没人看。逐目录列名单＝**名单本身变成一把要维护的尺**（新夹具目录默认不免检⇒又会红；全免检⇒又漏）。
- 现量代价小结（给编排者派单用）：**乙做完仍不足以满足"拿到一次 `go vet` 读数"**，除非同时处理 step 7 的 7 枚清单（要么逐枚整理台件、要么把 step 7 的分母也换成 tracked 名册）。

### 丙（不动、只写文档）
- ⚠ **现读事实：本仓没有一份"CI 拓扑文档"**。`grep -rln "lint::|步级|step conclusion|job 名" docs/reports docs/BUILD.md docs/PRECHECK.md docs/SLO.md docs/HOTKEYS.md`＝只命中 `docs/reports/HANDOVER.md` 与 `docs/reports/pending-and-issues.md`。⇒ 票面丙那句"写进 CI 拓扑文档"**零载具**。
- `docs/BUILD.md` 有 `## 7. CI 备注（票 08 前不建 workflow）`（`:129-135`，5 行，讲的是 runner 装工具链），且文件头 `:3-5` 自称 **FROZEN（S0 ticket 01 定稿）**，并要求"改其中任何一个必须同步改 `deps.toml`／`build.ps1`／`fetch-deps.ps1` 并重新走干净 clone 验证"⇒ 往这里加一节 CI 拓扑＝顺手动一枚**自称冻结**的交付物，代价与收益不匹配（除非编排者明确豁免同步条款）。
- 最不会被漏读的组合（本腿的代价比较，非裁定）：**`docs/reports/pending-and-issues.md` 追加一条 A##**（AGENTS §1.5 逐字指定它是真相源台账、"只追加不删"，每枚腿被要求读）＋**同一句在 `HANDOVER.md` §4 最新时间戳那节落一条指针**。⚠ 只写 HANDOVER 的代价＝它是停车点、按时间戳追加，下一位只读最新一节就会漏掉根因（今天 `:360`/`:397` 已经写过，仍需要票 236 才浮出来，这就是证据）。
- 丙零改动的代价：**lint 会一直红、两道 `go vet` 会一直被吃**（今晚实测），且每次推送都多花一遍 `go install mvdan.cc/gofumpt@latest`（`ci.yml:138`）才红；写腿会把"CI 没测"当环境噪声——这正是票面"为什么必须先修"第 2 条的形状。

---

## 台件名册（`.scratch/wisp/probes/236/a1/`，只建不删）

| 文件 | 作用 |
|---|---|
| `start-status-porcelain.txt` | 起手名册（151 行，终态比对基准；`.scratch/wisp/probes/236/` 那一条**已在起手名册里**，因为重定向先建了目录） |
| `end-status-porcelain.txt` | 终态名册（175 行）；与起手的差**全部由别的腿造成**，见下一节 |
| `start-anchor.txt` | 起手锚点 |
| `tracked-go-zlist.bin` | `git ls-files -z '*.go'` 原始名册（726 枚） |
| `gofumpt-tracked-stdout.txt` | 全仓 tracked 尺的 `-l` 清单（7 行） |
| `gofumpt-tracked-stderr.txt` | 解析错误行（1 行，`fs_broken.go:4:1`） |
| `census.md` | 本文件 |

---

## 运行期名册漂移（⛔ 不是我造成的，但会影响"终态＝起手"这条判据）

起手 `82110540`／151 行 → 交回时 **`5759d7fb`**／175 行（其间推进 **4 枚提交**）。`diff start end` 里**没有一条与 `236` 有关**⇒ 本腿的建件没有让名册相对起手发生变化（`.scratch/wisp/probes/236/` 那条在起手名册里就已存在）。逐条归属：
- `235-r1` 落了 3 枚：`795ed767`（骨架＋`probes/235/r1/readings.md`）、`1d2ad737`（`subagent_197_test.go` 注释 `:398-410`）、`a71f0be9`（`subagent_222_test.go` ＋93 行，三枚 222 用例加前置读数 `awaitChildOnBridge222`）⇒ 名册里 `?? .scratch/wisp/probes/235/` 从折叠目录展开成 21 条文件。
- 编排者落了 `5759d7fb`（台账 A454＋对照表＋票 230＋两枚 probes 件），并在索引里暂存了 5 条（前缀 `M `/`A `＝**暂存，不是我**——本腿从未 `git add`）。
- 另有两枚新未跟踪目录 `probes/132/`、`probes/v5-survey/`（同样不是我）。
- ⚠ **对本腿读数的影响范围已界定**：`git diff --name-only 82110540..5759d7fb`＝**10 枚文件**，其中只有 2 枚是 Go（`subagent_197_test.go`、`subagent_222_test.go`）。本 census 引到的其它文件（`internal/tools/task.go`、`task_cancel_221*.go`、`task_output_leg_test.go`、`ticket176r1_start_port_test.go`、`ticket175r2_stamp_live_test.go`、`task_backfill.go`、`bridge.go`、`cancel.go`、`subagent_197.go`、`internal/agent/loop.go`、`cmd/wisp/run.go`、`.github/workflows/ci.yml`、`attrib.sh`、`scripts/portable-tests.sh`、`docs/evidence/s1/185-*.md`）**一枚都没被这四推碰过**⇒ 那些 `file:line` 到 `5759d7fb` 仍然成立。
- **本腿的净改动＝只在 `.scratch/wisp/probes/236/a1/` 下建文件**（只建不删）。`git diff --name-only` 列出的跟踪文件全部在起手名册里已是 `M`/`D`（`.gitignore`、`probes/152`、`probes/161/r6/logs/**`、`design/**`、`docs/evidence/s1/152-*.md`）⇒ 跟踪文件零改动、零 commit、零 push。
- ⚠ 因此"终态名册＝起手名册"这一条**在本轮客观上不可能成立**（同时段有产码腿与编排者在写）。上面这份归属就是让编排者能用"逐条对得上别人"代替"名册相同"来验收本腿。

---


1. **AC#1／AC#2／AC#3 的突变读数本腿一枚都没跑**（⛔ 235-r1 独占编译与突变）。上面给的是"落点＋载具＋形状推理＋静态零命中读数"，不是"摘掉后恰 1 枚红"的实测温差。派写腿时这条必须补成两发读数（红→还原绿）＋`git status --porcelain -- internal cmd` 为空。
2. **AC#4 的"过期是否长期成立"我判不了**：我只验到今晚两发（`36556778053`／`36559498617`）该步为绿，且那两发的日志是 `--log-failed`（不含成功步的完整正文）。托管 runner 的 8.3 短名策略由 GitHub 决定，⛔ 我无法预测它哪天变回来——注释里最好写"读数＋日期＋run 号"而不是"今天为真"，但这属于裁法，不属普查。
3. **甲/乙/丙 三支的"任选一支之后必须真拿到一次 `go vet` 读数"这条完成判据，现量给不出任何一支单独能满足**（甲：step 7 仍红于 7 枚清单；乙只改 `attrib.sh`：红的那步根本不是它；丙：不动）。⇒ 要么把判据改成"step 7 绿"（＝必须处理全部 8 枚脏样），要么接受"这一步的颜色不是本轮能翻的"。这是**门的判据**层面的决定，⛔ 本腿不代裁，也不许写腿顺手把 step 7 改成 `continue-on-error`。
4. **未逐一读那 7 枚 `-l` 脏样的正文**（只读了前 4 行做形状归类）；其中 `197/r1c/pre/` 两枚是产品码快照，若编排者决定"顺手把它们格式化"，会改动**别人台件的字节**，与 AC#6 甲同一条风险（要挪到别的票去算）。
5. **票 225 的 DEFERRED↔`SPEC-12 §5` 双向对账尺我未普查**（只确证 d22scan 里 0 命中 DEFERRED）。AC#2 的 (b) 支改完之后，是否要顺带规定"标记行仍必须与 §5 登记表 1:1"属 225 的射程——两票的账别并成一张。
6. **AC#3 的"来源钉"具体落点没定**：`run.go:892` 用 `res.TaskID`、`loop.go:323` 填 `RunningTask.ID`，两者相等这件事该由 `cmd/wisp` 侧还是 `internal/tools` 侧钉，取决于写腿拿到的是 CLI 接缝还是包内 harness；本腿只列出候选载具（`pointer_183/185_cli_seam_test.go` 那两枚形状）。
7. **同族那 3 枚 `subagent_197.go` 的无尺支（AC#1 末节）在不在票 236 射程内**——需编排者一句话；不在就单立一枚，别让它 silently 变成"以后加固"。
8. **本腿所有行号是按起手锚点 `82110540` 现读的，交回时 HEAD 已到 `5759d7fb`（推进 4 枚）**。本腿复核过唯一被撞到的引用（`subagent_197_test.go:201` 未漂移），其余引用文件都不在这四推的 `git diff --name-only`（10 枚、其中 Go 只有 2 枚）里；⚠ 但写腿派单前**仍应按最新 HEAD 重跑一遍本文件每张表里的 `file:line`**（本仓已因"行号偏 1"打回多次，而漂移是共享工作树的常态，不是我的读数错误）。
9. `internal/tools/subagent_222_test.go` 本腿只读了 harness 与新 helper 名段（`:251-322`、`:378-413`），**没逐枚读三枚 222 用例正文**（它在 `235-r1` 手上正在变）⇒ 若写腿想复用 `h222` 那套"桥位计数／到达读数"工具，必须先确认它此刻的形状，且 AC#1 那两枚支路并不需要它（见 AC#1 的 h222 判定：**NoGate 桥进不了 Execute**）。
