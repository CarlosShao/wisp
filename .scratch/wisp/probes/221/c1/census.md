# 票 221-c1 只读普查：`task.spawn` 说明书对模型许诺了没注册的 `task.cancel`

- 锚点：`git rev-parse --short HEAD` ＝ **beaeaeba**（现跑）。
- 票面：`.scratch/wisp/issues/221-task-spawn-description-promises-task-cancel-that-is-not-registered.md`（标题核对＝「对模型许诺了一枚不存在的工具」，未派错号）。
- 工作树现状（`git status --short`，12:45）：`internal/config/allowdirs.go`／`loader.go`／`permmode.go` 为脏（另一枚写腿在飞）。**本腿零编译类命令：一枚 `go test`／`go build`／`go vet`／`gofumpt`／`d22scan`／`gate-clauses` 都没跑**；下文凡「只能靠跑测试才能证」的判据，一律标〔今天没有读数，由后续腿跑〕，不猜结果。
- 下文所有行号为锚点 beaeaeba 工作树**现读**，票给行号仅作对照；所有路径报到包级。
- 禁区遵守：`frontend/**`、`design/**` 零接触（未读、未转述）；`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 只读；未改任何已跟踪文件。

## ① 现实形状复算（读数时刻 12:45–12:46）

- 行 1（说明文字许诺 `task.cancel`）：**属实，行号错**。票面写 `subagent_197.go:189`，实际句在 **`internal/tools/subagent_197.go:196`**（Description() 整体 :194-200，`+` 号拼接片段，逐字与票面引文一致：「它会在任务名册里留下一行有父子关系与状态的记录，可以用 task.cancel 单独停它；」）。现跑 `grep -rn task.cancel --include=*.go`（排除 .scratch）生产命中＝`internal/tools/subagent_197.go:196` ＋注释行 `internal/tools/task.go:22`、`:23`、`:576`、`:593`、`cmd/wisp/panel_pump.go:125`。
- 行 2（`task.cancel` 没注册）：**属实，行号精确**。`internal/tools/task.go:595-597` `BuiltinTaskEntries` 逐字返回一枚 `Entry{{Tool: taskOutput{d: d}, Decl: TaskOutputDecl()}}`——只有 `task.output`（Name 定义 `internal/tools/task.go:479`）。头部 DEFERRED 标记本句实为 **`:22-:23`**（票面说「`:20` 一带」；整段 :18-27，见 §⑤）。
- 行 3（`TaskRoster.Cancel` 生产零调用者）：**属实，行号精确**。函数在 `internal/tools/task.go:442-459`。现跑 `grep -rn` 扫全部 `.Cancel(` 调用（排除 `_test`、.scratch）＝9 处：`internal/agent/loop.go:300`（RunningTask.Cancel 定义）、`:346`、`:531`，`internal/observe/goroutine.go:324`，`internal/observe/logging.go:116`，`cmd/wisp/run.go:678`（replyRoot），`cmd/wisp/approval_always.go:100`（root），`cmd/balldebug/main.go:339`、`:475`——全是 RunningTask／observe 根句柄，**名册那枚（TaskRoster.Cancel）生产调用者＝0**。测试调用者＝3 枚：`internal/tools/subagent_197_test.go:764`、`:789`、`:792`（都在 `Test197CancelIsPerRowAndNeverCascades`，函数起 :722）。票面「唯一命中＝subagent_197.go:39 的一句注释」**不完**：`internal/tools/task.go:175-176`（"read back by Cancel"）与 `:438` 的函数自述注释也是非测试文本命中（不改「零调用者」结论）。
- 行 4（父取消不级联已做到）：**属实，行号全漂移**。`context.WithCancel(context.WithoutCancel(ctx))` 实为 **`internal/tools/subagent_197.go:334`**（票面 :327）；两处「Parent cancellation must not cascade (ticket 197 §0)」注释在 **:316-319** 与 **:330-333**（票面 :309/:323）；模型可见结果文本含「没有被级联取消」在 **:386-390**（句在 :388；票面 :360-365）。形状侧证：`internal/tools/subagent_197_test.go:760-761` 断言父取消后孩子 State 仍 Thinking、`:785-786` 断言兄弟不动〔两枚是否绿＝〔今天没有读数，由后续腿跑〕〕。

**① 结论**：四行事实全部成立；票面行号 4 处漂移/失准（:189→:196、:327→:334、:309/:323→:316/:330、:360-365→:386-390），当场更正，不迁就。

## ② 乙形射程（读数时刻 12:46–12:55）

**现跑 `grep -rn 单独停它`**：生产命中恰 2 处——`internal/tools/subagent_197.go:196`（Description）与 `:389`（等待失败文本尾句）；另 4 处在 .scratch 证据快照（`probes/197/r1c/pre/`、`probes/197/r3b/pre/` 各 2），非生产面。⚠ 票面数漏的**第三处半许诺**：`:197` 句尾括号「（要停它得单独停）」（现跑 `grep -rn 单独停` 才抓到，全仓生产命中＝:196/:197/:389 三处）——今天没有任何工具能"单独停"，这半句同样是对模型的许诺。⇒ 乙形射程＝**三处**。

**逐处「改后仍然为真」的最小措辞候选**（为真凭据附后）：

1. `internal/tools/subagent_197.go:196`：「……有父子关系与状态的记录，可以用 task.cancel 单独停它；」→ 「……有父子关系与状态的记录；」。当下为真＝`internal/tools/task.go:337-342`（PublishSubagent 写 Kind/ParentTaskID/State）＋`:210-237`（Record 落行）。
2. `internal/tools/subagent_197.go:197`：「注意：停掉父任务不会级联停掉子代理（要停它得单独停）；」→「注意：停掉父任务不会级联停掉子代理；」。当下为真＝`internal/tools/subagent_197.go:334`（WithoutCancel 丢父级联）。
3. `internal/tools/subagent_197.go:388-389`：「……它仍在名册里，可以单独停它：它的流键是 %s。」→「……它仍在名册里，会把这一轮跑完并在名册里留下终态：它的流键是 %s。」当下为真＝finish watcher 无条件跑 finalize（`internal/tools/subagent_197.go:372-380`；注释 :362-366 逐字 "INCLUDING on the path where the parent stopped listening"）。

**会不会打红已有钉子——直接回答：不会，前提是四个被钉串一字不碰。**逐枚判：

- `internal/tools/subagent_197_test.go:418-428`（`Test197SpawnDescriptionNamesTheRealPoolCap`，票面说的 :419 那枚）：只断言 desc 含「上限 4 枚」「深度 1」（数字读自 :198 的 `%d`）——不碰 :198 数字段 ⇒ **不红**。
- `internal/tools/subagent_197_test.go:797`（在 :722 `Test197CancelIsPerRowAndNeverCascades` 内）：断言 desc 含「停掉父任务不会级联」——该子串在 :197 片段头部；**「级联不承诺句」（:197）与「cancel 许诺句」（:196）是 `+` 拼接的两枚不同字符串字面量，不是一句**：候选 1/2 保留 :197 整前缀 ⇒ **不红**。
- 同测试 `:747` 与 `:562`：断言结果文本含「没有被级联取消」（:388 前半）——候选 3 保留该子串 ⇒ **不红**。
- 全仓 `grep -rn task.cancel --include=*_test.go`＝**0 命中**（我复核数＝0，与派单"先前跑＝零命中"一致）；`grep -rn 单独停 --include=*_test.go`＝0 ⇒ :196/:389 文本无任何测试钉。
- `internal/tools/bridge_test.go:169-170`：仅要求各工具 Description() 非空 ⇒ 不红。
- 包外仪器：`tools/gate-clauses.sh`、`scripts/d22scan.sh`、`tools/d22scan/allowlist.txt` 现跑 grep「221／task.cancel／单独停／停掉父任务／没有被级联取消」＝**零命中** ⇒ 无包外措辞钉。
- 文档面含该句者（`docs/evidence/s1/197-subagent-entity-r1.md`、`docs/evidence/s1/222-spawn-permit-release-r1.md` 等，12:59 现跑）全是历史证据快照，不是判据钉；乙形不许也不需要改。

〔乙形改完全包是否真零红＝只能靠跑测试才能证：**这条今天没有读数**（config 写腿在飞），由后续腿按票 AC#5 跑并逐名比红名集合。〕

## ③ 甲形射程（读数时刻 12:47–12:56）

**注册形状（要动哪几面）**：

- 件形：`Entry{Tool, Decl}` 在 `internal/tools/registry.go:71-74`；`Registry.Register` 在 `:95`；Tool 接口四法（Name/Description/Parameters/Execute）在 `internal/tools/tool.go:20-31`（Name 须为 `namespace.action` 形）。
- 新工具落点＝**`internal/tools/task.go`（与 taskOutput 同包同族）**：仿 task.output 四件——Name/Description/Parameters `internal/tools/task.go:479-485`、Decl `:581-590`、产出入 `:595-597`——新增 taskCancel 形状并让 `BuiltinTaskEntries`（:595）返回两枚。装配消费端**不需要动**：`cmd/wisp/run.go:424-429` 是 `for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths})` 的循环注册，条目数天然吸收。对照 `task.spawn` 一族：声明 `internal/tools/subagent_197.go:209-218`（TaskSpawnDecl：Capabilities/Needs 皆 nil、Declared L0、Resident、KindBuiltin）、产出 `:224-226`（BuiltinSubagentEntries）、消费 `cmd/wisp/run.go:579-594`。
- **TaskDeps 不需要新字段**：`internal/tools/task.go:30-52` 已有 `Roster *TaskRoster`（:36），取消句柄就住在 roster 内（`:177` cancel map；`AttachCancel` :413-425、`DetachCancel` :429-436；attach 调用点＝`internal/tools/subagent_197.go:337`）。⚠ **更不许往 `SubagentDeps` 加字段**：它今天恰 5 枚（Roster/BaseOptions/ParentTools/Provenance/Stream，`internal/tools/subagent_197.go:128-161`），字段枚举手守卫在 `internal/tools/subagent_197_test.go:826-833`（reflect 枚举）＋`:871-873`（`Test197SubagentHasNoSelfApprovalOutlet`，函数起 :838）钉「!= 5 即红」——加第 6 枚必打红这枚「子代理永不自批」守卫。甲形不需要碰它。
- 全仓没有任何测试钉 `BuiltinTaskEntries` 的条目数（12:52–12:53 逐名查 `*_test.go`：`internal/tools/task_output_leg_test.go:141-149`、`internal/tools/pointer_183_cli_seam_test.go:111`、`pointer_185_cli_seam_test.go:103`、`ticket175r2_stamp_live_test.go:79` 全是 range 循环消费，无长度断言）。〔"注册第二枚后全仓真零红"仍属〔今天没有读数〕，由后续腿跑。〕

**取消的真能力在不在（逐行读 `internal/tools/task.go:442-459`）**：

- `TaskRoster.Cancel(taskID)` 做三件事：查行（:446-449 同时取 cancel 句柄与名册行）；**无行**→false「查不到任务…v1 的名册只在本进程内」（:450-452）；**有行无句柄**→false「没有在跑…没有可停的句柄」（:453-455）；**有句柄→调它**（:457 `cancel()`，返回 true）。⇒ 真停：句柄＝spawn 时存进去的孩子自己的 `context.CancelFunc`（`internal/tools/subagent_197.go:334` 造、`:337` attach），**不是只置状态**。
- 但它**自己不写状态、不记「被谁停」**：签名 `(taskID string) (bool, string)` 无调用者参数（:442）。终态由既有链路补写：cancel 触发→孩子环路终结→`bg.Wait()` 回 cancelled→finish watcher（`internal/tools/subagent_197.go:372-380`）跑 `finalize`（:401-423）：`:404-405` 把 `agent.StatusCancelled` 映射为 `subagentStateStopped`＝**Muted**（:106-111；:100-105 注释逐字 "D43 has no cancelled; Muted is D43's name for silenced on purpose"），`:413-416` Record 落行、`:417` DetachCancel（`internal/tools/task.go:429-436` **只删句柄不删行**）、`:421` feedDone 关流。⇒ **现读答案：停一名孩子后那行名册不凭空消失，终态＝Muted，孩子自己的流有终态（done）；「被谁停的」今天无处落地——甲形 AC#2 要这一条，必须改 Cancel 签名（带 caller）或在落行时加载体，这条要写死进判据。**
- 若目标是 root 行：root 从未 AttachCancel（attach 只在 spawn 路径），Cancel 走 :453-455 拒停——天然封「停根」。

**不许新造状态名——现读裁决：甲形不缺状态名**：D43 词表恰 20 枚逐名在 `internal/statemachine/states.go:11-31`（`StateMuted` ＝:14），judge＝`Valid`（:39-49）；名册行状态字段＝`TaskOutput.State`（`internal/tools/task.go:125`），写侧全过 `StateAnswer`（:153-163）。没有 "Cancelled"，但 **Muted 已被子代理层用作"宿主停掉"的现成停靠位**（`internal/tools/subagent_197.go:109`）。⇒「甲形缺状态名、这半支要动 D43＝契约面」**不成立**；甲形真正碰的契约面只有 D34 行（`docs/PLAN.md:2564`，task.cancel 标 **L1** 非 L0）、`docs/PLAN.md:1531` DEFERRED 行的到期改判、以及台账 A##（§⑤）。

**权力边界（AC#3 可满足性，现读）**：

- 今天 Go 侧**能**知道"当前这次工具调用属于哪枚任务"：ctx 载具 `cancelHandle{corr, bus}`（`internal/tools/cancel.go:34-51`）由桥注入（`internal/tools/bridge.go:457`，corr＝`orDefault(req.CorrelationID, req.TaskID)`；:246-247 CorrelationID 缺省回落 TaskID），工具侧导出读法 `CorrelationID(ctx)`（`internal/tools/cancel.go:53-58`）；生产派发逐字 `TaskID: taskID, CorrelationID: taskID`（`internal/agent/loop.go:647`）⇒ 生产面上 corr 恒＝调用者 taskID。task.spawn 已经在用它当父身份（`internal/tools/subagent_197.go:250` 取 parentID；`:255` Look 读 Kind 做深度闸）。父子读法＝`internal/tools/task.go:263`（Look→`ParentTaskID` 字段 :128）与 `:300`（Descendants）。⇒ **「调用者 taskID == 目标行 ParentTaskID 才准停」今天机械可判；AC#3 两枚正控（子停兄弟：ParentTaskID≠caller；子停自己：目标是 caller 自身或 Kind=root 无句柄）可满足——不要把 AC#3 判成不可满足。**
- ⚠「或用户」那半今天**没有载体**：root taskID 每轮新铸（`internal/agent/loop.go:321-322`、`:332-333` newTaskID），用户新一轮说话≠派生那轮的 ParentTaskID，等值闸会拒"用户跨轮停子"；工具 ctx 的值载体全集只有三枚——cancelKey（`internal/tools/cancel.go:34`）、hostPathBox（`internal/tools/bridge.go:595-607`）、inFlightSlot（`internal/tools/bridge.go:682-695`），无一携带用户身份。⇒ **甲形判据只能钉"同轮父调用"形；"用户单独有停子的路"要么由 owner 另裁通道（面板/veto：`internal/tools/cancel.go:8-31` 的 CancelBus 是调用级否决、不是名册级停子），要么把许诺措辞改到不需要它。派单必须写死这条，不许写腿默默发明身份载体（那是新权力面）。**

## ④ AC#1 能力尺的落点候选（读数时刻 12:47–12:53）

**放哪**：新测试文件 `internal/tools/description_registry_221_test.go`（包内测试即可读四家 Builtin 构造器与 Registry；测试件可以读，但**今天不许跑**，本节只给形状不给绿）。形状：把四枚生产构造器的产出注册进 `tools.NewRegistry`（装配现场同款调用＝`cmd/wisp/run.go:392`）——`BuiltinFSEntries`（`internal/tools/fs.go:325`）、`BuiltinFSWriteEntries`（`internal/tools/fs_write.go:772`）、`BuiltinTaskEntries`（`internal/tools/task.go:595`）、`BuiltinSubagentEntries`（`internal/tools/subagent_197.go:224`）——得「已注册名全集」；再对每枚 `Entry.Tool.Description()` 跑 `task\.[a-z]+` 抽词根，断言每个词根 ∈ 全集；红时报「工具名＋该 Description 全文＋缺的名」，满足票面"红句能指出是哪一句、缺哪一枚"。**不要走全仓源文本扫描**：那会命中注释与历史证据件（`internal/tools/task.go:22-23` 注释本身含 "task.cancel"），把 DEFERRED 登记面误判成许诺面——票 AC#1 自书"扫能力，不扫词面"。
**今天这把尺红在哪几枚（静态推导，非跑测）**：全集＝{fs.read :124, fs.edit :70, fs.write :234, fs.trash :369, fs.move :440, fs.delete :583, task.output :479, task.spawn :187}（八枚，出处包级路径见上）。逐枚 Description() 现读：task.output 的（`internal/tools/task.go:481-483`）不含 `task.` 词根；六枚 fs 的描述不含；**唯一含词根的＝`internal/tools/subagent_197.go:194-200`，且其中只含 `task.cancel` 一枚、它不在全集 ⇒ 今天这把尺红恰好一枚、红句＝task.spawn 的 :196。除 task.cancel 外无第二枚「许诺未注册」——这把尺今天不会挖出新名字**。
**会不会误伤**：(1) 若将来有 Description 合法提到在册的 `task.output`/`task.spawn`：等值判定放行，不误伤（全集含这八枚）。⚠ 真正的误伤面是**包内测试台的假工具**：`internal/tools/bridge_test.go:36`、`:234`、`internal/tools/ticket90_test.go:68`、`internal/tools/subagent_222_test.go:131` 等自带 Description——它们**不进** `tools.NewRegistry`（尺只读四家 Builtin 构造器的返回集），天然不被扫。⇒ **扫描口径必须写死＝"只扫四家生产构造器返回的 Entry.Description()"，既不扫 `_test.go` 源文本、也不扫注释、也不扫 .scratch。** (2) MCP/本地目录工具今天不进注册（`docs/specs/SPEC-12-roadmap-governance.md:80` MCP＝RESERVED 无实现），尺不涉及。

## ⑤ 契约/账目面（读数时刻 12:48–12:54）

**`internal/tools/task.go` 头部逐字（现读 :18-27）**：

```
// The D34 task family (ticket 164). Only one of the three registered names has
// an implementation here, and that is deliberate:
//
//	task.output  L0  read what a background task printed   (this file)
//	task.list    -   DEFERRED with five fields, PLAN.md §7 :1531
//	task.cancel  -   DEFERRED with five fields, PLAN.md §7 :1531
//
// task.output is AC#3's leg. The two rows above stay unregistered because the
// §7 registry says so, and a half-claim is exactly the "在册 + 无实现 + 无人认领"
// shape ticket 164 AC#1 was written to end.
```

**与 `docs/specs/SPEC-12-roadmap-governance.md` §5 对照**：§5 表在 `:59`（标题「推迟项登记表（M0 要求①：五字段缺一不可）」）、行段 `:66-92`。⚠ **§5 里 `task.list`／`task.cancel` 各零行**（现跑 grep 全文＝零命中）——登记实体不在 SPEC-12，而在 `docs/PLAN.md` §7（标题 `:1507`「七、推迟项与技术债登记表（M0 要求①）」）的 **`:1531`**：逐字「**DEFERRED** | **`task.list` / `task.cancel` 的实现**（D34 名册在册、生产注册表零枚）｜票 164 AC#1 现量（2026-09-26）…**先补名册裁定再谈实现**｜（完成判据）注册表里 `task.list`／`task.cancel` **各有一枚真实现**并过契约测试｜依赖＝票 163＋票 164 AC#1 裁定表｜代价＝Agent 侧查不到"谁被阻塞"，D31 只剩 UI 半条腿」。⇒ 是 **DEFERRED**（非 RESERVED/REJECTED），且票面 `:22-23` 说「五字段」＝这一行的五列。另两枚现读：`docs/PLAN.md:2564`（D34 权威表 §16.5.2，标题 :2516）行「task.list / task.cancel｜查看/取消任务｜**L0 / L1**｜—｜**S7**」；`docs/PLAN.md:3115`（S7 切片验收用例含 `task.list`/`task.cancel`）。
**⚠ 甲形若要「DEFERRED 标记与 §5 登记表双向 1:1」＝在要求一枚今天不存在的仪器，判定成立**：SPEC-12:94-95 逐字把 1:1 仪器的定义写成「代码内每个 `// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点` 必须在表有对应条目」，而 (a) `docs/DEFERRED.md` **不存在**（12:48 `ls`＝No such file；它只以规划身份出现在 `docs/specs/SPEC-12-roadmap-governance.md:99`、`docs/PLAN.md:1431`）；(b) `docs/specs/SPEC-12-roadmap-governance.md` §5 **没有** task.* 行可对接；(c) 票 225 已量（`docs/evidence/s1/225-deferred-registry-audit.md:41`）：「§5 那 12 行 DEFERRED 里只有 2 行在代码里有对应标记」，双向 1:1 今天本来就断。⇒ **派单不许把「1:1 对得上」写成甲形 AC#4 的通过条件**（那是先修仪器票 225 的活）；甲形的正确账＝**摘掉 `internal/tools/task.go:22-23` 的 task.cancel 标记**（它已不是推迟）＋ SPEC-12:94 的"反向亦查"在后续 225 落地前不可证。乙形则 `:22-23` 一字不动（票 AC#4 后半，AGENTS §1.1）。

## ⑥ 别家怎么做的（外部快照，只读；读数 12:49–12:54；基准目录 `D:/work/AI/open source/`；快照无 .git，以下零提交史引用）

**独立工具名 vs 父任务的一个动作**：

- **deepseek-harness＝独立工具**：模型侧 `job_kill`（`deepseek-harness/packages/jobs/tool-jobs/src/index.ts:372-373`，描述逐字 "Request cancellation of a running background job."；system 引导句 :253 教模型收集后 "job_kill jobs that stopped mattering"）；人侧不走工具、走 Remote `job.kill`（`.agents/notes/implemented/feature/2026-08-26-human-job-kill.md:19`），模型自己停与人停要在"通知是否欠模型"上分账（同 note :5、:16）。
- **minimax（agent-modules）＝独立工具**：`task_stop`（`minimax-code/packages/agent-tools/src/desktop/builtin-defs.ts:699-702`，与 `task_output`（:691-696）成家族；实现 `minimax-code/packages/agent-tools/src/desktop/local-task-control.ts:162-188`）。
- **Step-Code＝父任务的一个动作，不是独立工具**：`agent_send` 的 `action: "reply" | "stop"` 参数（`Step-Code/packages/coding-agent/src/features/step-subagent.ts:614`、schema :427-429，描述逐字 'action:"stop" interrupts the lane and ends its child process'）；注释记录旧独立件 `agent_reply/agent_wait/agent_interrupt/agent_list` 已硬删并入（:646-648）。
- **openchamber＝找不到**：快照顶层无 `src/` 工具目录，grep `TaskStop/task_stop/subagent` 零命中（12:53 现跑）——明写找不到，不编。

**取消后名册那行留什么状态**（三家一致：行不删，落**既有**状态词，无一新造词表）：

- DSH：终态 ∈ `'completed' | 'killed' | 'failed'`（`deepseek-harness/packages/jobs/jobs/src/types.ts:17-18`，注释逐字 "cancelled (`killed`)"）；kill 原因并进 `killed` 的 detail 文本、模型与人同见（`packages/jobs/jobs-local/src/index.ts:584`；note :18 逐字 "signal: SIGTERM; cancelled by the user"——**"被谁停的"落在那一行本身**）。
- minimax：`'stopping'` 中间态→`'canceled'` 终态，行留在 store（`minimax-code/packages/local-runtime/src/background-task/task-stop.ts:31-52`；模型可见状态枚举含 stopping/canceled：`packages/agent-tools/src/desktop/builtin-defs.ts:644-649`）；原因走 `lastError.code`（TASK_STOP_REQUESTED/TASK_CANCELED）落在行上（task-stop.ts:33-35、:46-48）。
- Step-Code：lane 状态 ∈ `"running"|"completed"|"failed"|"aborted"`（`Step-Code/packages/coding-agent/src/features/step-subagent.ts:159`），停后行仍在 details，父侧文本渲染成 "interrupted"（`packages/coding-agent/src/features/subagent/execute.ts:58`、`:284`；共享 stop 动词 `src/features/subagent/lane-lifecycle.ts:261`）。

**有没有"谁都能停"的防护**：

- DSH：注册表 owner fence 是**唯一**访问规则，unknown/foreign 一律并成 `job/not-found`；子会话自己的 job 从它自己列表可停（note :19）；工具侧调用带归属 `exec.agent?.id`（`tool-jobs/src/index.ts:402`）。
- minimax：会话级可见闸 `isVisibleToOwnerSession`（ownerSessionId == ctx.sessionId，`minimax-code/packages/local-runtime/src/background-task/service.ts:457-462`），`stop()` 先过它（:272-290）——**同会话任意调用都能停**（比"只有父轮能停"宽，是会话粒度不是 turn 粒度）。
- Step-Code：**无防护**——模型可 `agent_send{to:{all:true}, action:"stop"}` 停掉全部 lane（`Step-Code/packages/coding-agent/src/features/step-subagent.ts:625-643`），即票 221 禁区明令不许开的"任何任务都能停任何任务"形，作反面警示件。

## ⑦ 我给下一位的三件（收尾读数 12:56–13:04）

**A. 票面写错/写漏（逐条，均现跑复核）**：

1. 行号漂移四处：`subagent_197.go:189`→**:196**；票 AC#3「:327 那形」→**:334**；现量表「:309/:323 注释」→**:316-319/:330-333**；「:360-365 含没有被级联取消」→**:386-390**。（`internal/tools/task.go:595`、`:442` 两处精确。）
2. 乙形射程少列两处：票正文只圈了「:189 那一句」——现读该改的是 **三处**（:196、:197 尾括号「（要停它得单独停）」、:389「可以单独停它」；13:03 现跑 `grep -rn 单独停` 生产命中恰此三处，测试命中 0）。:197 那半句藏在**被 :797 钉住的同一枚片段里**，派单不点名，写手要么漏改要么把钉拽红。
3. 「非测试唯一命中＝subagent_197.go:39 注释」不完：`internal/tools/task.go:175-176`、`:438` 的函数自述注释也是非测试文本命中（零调用者结论不变）。
4. AC#1 口径自撞：「词根必须在 `BuiltinTaskEntries` 的注册名册里」——`task.spawn` 恰恰不在 `BuiltinTaskEntries`（在 `BuiltinSubagentEntries`，`internal/tools/subagent_197.go:224`）；分母必须写死为**全体生产注册名并集**（§④八枚），否则尺一上线就指错对象。
5. AC#4（甲形半格）要求一枚今天不存在的仪器（§⑤判定：`docs/DEFERRED.md` 不存在、SPEC-12 §5 无 task.* 行、225 已量 12 行仅 2 标记）；且甲形落地那刻 task.cancel 的 DEFERRED 标记是「**应摘**」不是「应对上」——照票面原样写，甲形永远绿不了。
6. 「先例＝task.spawn 作为 D34 新增行」在 `docs/PLAN.md` 现文里**找不到那行**（12:56 现跑 grep 'task.spawn'＝零命中；task.spawn 的批准走台账形）。能引用的真实加行先例是 `docs/PLAN.md:2564` 的 **task.output**（行内逐字「票 164（2026-09-27 owner 批准新增）」）。两路对「要不要动 PLAN.md 表格」含义相反，派单必须让 owner 当场选路。
7. AC#2 担心的「缺状态名」**不成立**：Muted 是现成停靠位（`internal/tools/subagent_197.go:100-109`＋`internal/statemachine/states.go:14`），甲形不必动 D43＝契约面这一半可以划掉；真正剩的契约面只有 §⑤ 那三件。
8. 禁区闸门未开：票 201／票 220 的 issue 文件均无 `-done` 后缀（12:57 现跑 `ls .scratch/wisp/issues/`）——甲形同写面（`internal/tools/task.go`＋`subagent_197.go`＋`cmd/wisp`），**当前脏腿只碰 internal/config 不等于闸门已空**，派甲形前单独确认 201/220 让位。
9. AC#3 没有被判死，反而可满足（§③）；但「用户跨轮停子」没有 ctx 载体（root taskID 每轮新铸，`internal/agent/loop.go:321-322`/`:332-333`）——票面「只有父任务（**或用户**）能停它的孩子」里「或用户」三个字今天没有落点，派单必须要么删要么改判为"另裁通道"。

**B. 派单里必须写死的边界**：

- 乙形：只许动 `internal/tools/subagent_197.go` 的 :196/:197/:389 三处文本；四个被钉串一字不碰——「上限 4 枚」「深度 1」（`internal/tools/subagent_197_test.go:418-428`）、「停掉父任务不会级联」（:797）、「没有被级联取消」（:747、:562）；`internal/tools/task.go` 连注释一字不动（DEFERRED 标记保留＝票 AC#4 后半）；不新增导出名；不改 gate-clauses/d22scan 措辞（现读零命中，无需动）。
- 甲形：`SubagentDeps` 字段数＝5 是钉（`internal/tools/subagent_197_test.go:871-873`），不许为取消往它加字段；`TaskDeps` 现有 Roster 够用、也不加；新工具落 `internal/tools/task.go`＋`BuiltinTaskEntries`；Declared 按 D34 行取 **L1**（`docs/PLAN.md:2564`），不是 task.output 那枚 L0；权限闸＝`CorrelationID(ctx)`==目标行 ParentTaskID（载体见 §③），封死子停兄弟/子停自己/停 root；「被谁停」需要改 `TaskRoster.Cancel` 签名或加载体（今天没有，§③）；终态用 Muted、禁新造态名；PLAN.md／docs/specs／thresholds／golden／allowlist.txt 一字不动，批准只落台账 A##；:334 的 WithoutCancel 那形不许被顺手改坏（票 AC#3 既有）。
- 两支共用：AC#5 整包终态＋逐名比红名集合今天不可跑（config 写腿在飞）；改前/改后读数与 §④ 那把新尺的首跑都由后续腿出；cmd/wisp 测试记得带 sherpa PATH（票 AC#5 既有口径）。

**C. 本腿未证清单（全部〔今天没有读数〕）**：①-④ 各判据的 go-test 红绿、乙形零红断言的实跑证、§④ 能力尺改前一响、AC#5 终态与红名集合、`Test197*` 全家是否绿。原因同一：共享工作树 `internal/config` 有脏写腿，本腿零编译类命令。

## 入库清单

- 本文件：`.scratch/wisp/probes/221/c1/census.md`（本腿唯一写面；目录 `.scratch/wisp/probes/221/c1/` 为其容器）。
- 字节数（`wc -c` 现测，自引用定宽、含本节）：28379
- 本腿零产码、零编译类命令、零已跟踪文件改动；外部快照只读；`frontend`/`design` 零接触。
