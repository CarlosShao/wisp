# 328-a1 只读普查 —— 名册与三形代价（票 328 `AC#0` 的料，⛔ 判语 ⛔ 选形 ⛔ 翻框）

腿＝`328-a1`（只读；⛔ 产码 ⛔ 测试件 ⛔ 改断言 ⛔ push）。锚＝**`7535d83b`**（dev，2026-10-17 10:07:32 +0800，本腿 `git log -1` 现认；`4366c1ca`＝220-v1 那枚也已 `git log -1` 现认，10-16 12:18）。
代码一律 `git show 7535d83b:<路径>`，⛔ 读工作树在飞字节（起手 `git status --porcelain` 见 `cmd/wisp/**`＋`internal/proc/**`＝`228-r3`、`internal/panel/**`＝`333-r1`）。
★与锚的漂移本腿现量：`git diff --stat 7535d83b HEAD -- internal/agent/loop.go internal/panel/approval.go internal/panel/subagent_roster_197.go internal/agent/approval/queue.go`＝**空**（blob `ae8c25c756e6`／`61f6db5e0f73`／`35ce6ffb8603`／`66fec7ae3753`）。

起手四把尺（逐字输出＝`.scratch/wisp/probes/328/a1/logs/rulers-starting-328a1.log`）：票 328 未勾＝**5**（`:31-35`＝`AC#0`…`AC#4`）；`ApprovalCardView` 结构在 `approval.go:39`、首字段 `CorrelationID` 在 `:40`；`grep TaskID`＝**0 命中**（正控＝同尺 `CorrelationID` 命中 `:28/:40/:78`）；名册不对称在 `subagent_roster_197.go:213-214` ⇄ `:216-221`；那把尺＝`:277-288`（唯一读者 `:267`，全仓无测试直呼它）。
台面：本腿⛔ 共享工作树跑 `go test`／`go build`／`go vet`。颜色＝**仓外导出台树** `D:\work\workspace\wisp-328a1-census-tree`（`git archive 7535d83b go.mod go.sum internal cmd tools scripts frontend models docs`＋手工三件 `build/*.dll`／`third_party/sherpa-onnx/*.dll`／`frontend/dist`；8.3 短名＝`D:\work\WORKSP~1\WISP-3~4`，尺＝`cygpath -d`；`.go`＝**750** 证非空；`GOTMPDIR=D:\tmp\328a1-gotmp` 先 `mkdir -p`）。
★本腿唯一执行过的代码＝① 仓外纯 stdlib 探针 `D:\tmp\328a1\corrprobe\main.go`、② 同族第二枚 `D:\tmp\328a1\formprobe\main.go`（两枚皆 `go run` rc=0、stderr 0 字节、⛔ import 本仓包；探针①里四枚函数是**逐字重打**，`diff -B` 对锚逐字节相同＝`logs/probe-transcription-diff.log`）、③ 导出台树里三发 `go test -run`（⛔ 改过任何字节）。
读数口径：**量的**＝尺跑出来的；**算的**＝按代码那一行逐字推的；每条都带尺与射程目录。

## 一、`AC#0` 生产名册（今天真会走到 `callCorr` 的调用源）

尺＝`git grep -n 'callCorr' 7535d83b -- internal cmd`（量的；命中 **19** 行，非测试只有 **5** 行＝`loop.go:367`（注释）`:594`（doc）`:603`（def）`:676`（唯一调用）＋`subagent_roster_197.go:58`（一条注释））。

- **唯一铸造点**：`internal/agent/loop.go:676` `TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i)`，两支＝`:605`（`callID==""` ⇒ `taskID#call-<i>`）／`:607`（⇒ `taskID#<provider号>`）。
- 父链（逐枚现跑）：`executeCalls` 只被 `loop.go:504` 调 ⇐ `run`（`:346`）⇐ `RunAsync`（`:328`→`:333`）与 `Run`（`:339`→`:340`）；`taskID`＝`newTaskID()`（`:1157`，UUIDv4 形、**不含 `#`**，尺＝探针①真实实例值 `e2bbeaab-2b94-4a9e-ac59-ec13ce118077`）。
- 生产入口两枚（尺＝`git grep -nE '\.RunAsync\(|\.Run\(ctx' 7535d83b -- internal cmd` 去掉 `_test`、四行 `bridge.Run`（`cmd/balldebug/main.go:252`、`cmd/wisp/models.go:311`、`cmd/wisp/resident_ball_windows.go:362` 等，另 `cmd/balldebug/main.go:662` 是注释）与 `internal/ball/hotkey_reload.go:22/:125` 两行注释）：
  1. `cmd/wisp/run.go:1176` `bg := loop.RunAsync(ctx, task)`；2. `internal/tools/subagent_197.go:372` `bg := child.RunAsync(childCtx, prompt)`。
  ⚠ **票面/前腿那两个行号已漂**（前腿记 `run.go:1164`、`subagent_197.go:367`；本锚实测 `1176`／`372`）。`ToolRequest{` 非测试构造点全仓只有 `internal/agent/loop.go:675` 一枚（尺＝`git grep -n 'ToolRequest{' 7535d83b -- internal cmd`，另一枚命中 `internal/tools/bridge.go:890` 是注释）⇒ 生产 corr 只有 `callCorr` 一个出处。
- 到不到的两支（量的）：陪聊全双工与 realtime 到不了 —— `approval/gate.go:229-241` `admitCheck` 对未登记任务逐字回 `"实时语音（Path C）不具备工具权限；该任务未经文本循环登记（AdmitTextTask），已 fail-closed 拒绝"`，而 `AdmitTextTask` 只由文本循环侧持（`gate.go:182-194`、`loop.go` 的 `opt.AdmitTask`）。

每形 corr 逐字节 ＋ 那把尺剥不剥（量的＝探针①，射程＝`logs/corr-shape-probe-run.log`；`T`＝探针①现生成的真实例值）：

| 形（谁产的） | corr 逐字节 | `splitClashSuffix` 剥？ | base | `waiting[T]` 点亮行？ |
|---|---|---|---|---|
| 支B 模型没给号（`loop.go:605`，index=0/1/2） | `T#call-0`、`T#call-1`、`T#call-2` | **false** | `""` | **false** |
| 支A 模型给号·openai-chat 逐字透传（`toolasm.go:59`；号取仓内真值 `call_e1`／`call_p1`／`call_p2`，尺＝`loop_golden_test.go:66`） | `T#call_e1` 等 | **false** | `""` | **false** |
| 支A·openai-chat 空号回退（`toolasm.go:114` `call_unnamed_%d`） | `T#call_unnamed_0` | **false** | `""` | **false** |
| 支A·responses（`openairesponses/stream.go:446-451` 先 `call_id` 后 `id`） | `T#resp_9fK2`（示例位：仓内**无**钉住该形的前腿件） | **false** | `""` | **false** |
| 支A·anthropic（`anthropic/stream.go:258 id: blk.ID`） | `T#<wire 号>`（★更正前腿：仓内 `toolu_` **零命中**，`toolu_` 不是本仓量到的形状） | **false** | `""` | **false** |
| 队列撞名改写叠在支B上（`queue.go:168`） | `T#call-0#7` | **true** | `T#call-0` ≠ T | **false** |
| 队列撞名改写叠在支A上（`queue.go:168`） | `T#call_e1#7` | **true** | `T#call_e1` ≠ T | **false** |
| 队列 replay（`:614`／`:616`） | `T#call-0#replay`／`T#call-0#9` | false／true | `""`／`T#call-0` | **false** |
| 队列空 corr 回退（`queue.go:165`） | `approval-7`（无 `#`） | **false** | `""` | **false** |
| 宿主发起卡（`run.go:889`＋`approval_always.go:64`＋`config_reload.go:89`；corr 落 `gate.go:281` 的 `orDefaultText(d.CorrelationID, d.TaskID)`） | `host:mode-switch`／`host:fs-allow-widen`／`host:config-reload` | false | `""` | **false**（`panel_pump.go:232-239` 的名册行要 `rt.tasks.Look(id)` 有记录，`host:*` 从不入册） |
| **正控·名册今天认的那一形**（队列对**裸任务号**的撞名改写） | `T#9` | **true** | `T` | **true**（唯一点亮形，而生产不产它） |

★另数一发（AC#0 附带，量的）：派单给的起评尺 `git grep -nE 'CorrelationID:[[:space:]]*"[^"]*"' 7535d83b -- internal cmd`＝**62** 枚，其中尾部是撞名形（`…#<数字>"`）＝**0 枚** —— 硬写撞名形不在这把尺的射程里，它们经夹具 `card197()` 传参。换一把尺才数到（`git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' 7535d83b -- internal cmd`，量的）：**6 枚值＋1 枚注释**，全在 `internal/panel/subagent_blocked_220_test.go`：`:108`／`:112`／`:130`／`:131`／`:132`／`:186` 各一枚 `child-1#9`/`child-1#147`/`other-1#9`，注释 `:15` 提 `task-7#9`；假阳（`"AC#3"` 那一族）已逐枚剔。⇒ 涉及硬写撞名形的顶层用例 **3 枚**：`TestASecondCardForTheSameTaskStillLightsItsRow`（`:108/:112`）、`TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt`（`:130/:131/:132`）、`TestABlockedRowKeepsItsRunningDimension`（`:186`）（＋名册侧 `subagent_roster_197_test.go:108` 用的是**裸任务号形** `CorrelationID: "child-1"`，不经 `callCorr`，今天量到 `--- PASS`，见 `logs/panel-color-keynails-7535d83b.log`）。

## 二、三形代价表（每形同一条尺；⛔ 推荐，只列代价）

★先记一处**票面与三形的对应**（本腿现读票 `AC#1` 逐字）：甲＝卡片也带任务号；**乙＝改 `callCorr` 把可剥的数字放尾巴（票给的例子 `task#<seq>#<callID>`）**；丙＝放宽名册那把尺（按已知任务号前缀匹配）。⚠ 前一枚普查腿（件①②③，提交 `59c46b7a`/`fdea4402`/`135e9ed0`）把乙与丙的内容**互换**了，本件按票面重列，那一处请按"前腿缺陷"处理。

**(甲) 卡片也带任务号**
- (a) 要动：`internal/panel/approval.go:39-59` 加字段＋`:72-99 CardViewFromDecision` 填；`internal/panel/pump.go:98-113 NativeVerdict.CardView` ＋ `NativeVerdict` 本身要带任务号；`cmd/wisp/panel_pump.go:62-80 liveVerdicts`（任务号在 `it.Decision.TaskID` 里现成，尺＝`approval/pending_read.go:113-131`）；`internal/panel/subagent_roster_197.go:213-214` 改成两枚都索引；**`frontend/src/lib/panel.ts:24-47` 的 `ApprovalCardView` 接口**＝机主外部 agent 的写面 ⇒ 跨面。
- (b) 顶红名册：`internal/panel/approval_test.go:105` 的 `{"ApprovalCardView","ApprovalCardView"}` 一对，红句逐字族＝`:129` `"Go %s emits %v that interface %s does not declare"` ⇒ `Go ApprovalCardView emits [taskId] that interface ApprovalCardView does not declare`（算的）。★**这枚用例今天就红**（量的，导出台树 `logs/panel-color-targeted-7535d83b.log`：`--- FAIL`，红句逐字 `approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`＝`Snapshot`↔`PanelSnapshot` 那一对，与卡片无关）⇒ 甲加的是**同枚红件里的第二句**，颜色分不开、只有文案分得开。另一枚同族红＝`composer_test.go:52 TestComposerContractTypesMatchFrontend`（量的：`--- FAIL`，两句红）。⛔红：`cmd/wisp/panel_assets_143_test.go:52 cardView143` 是影子结构＋`:109 decodeCard143` 容忍未知键 ⇒ 不破（读码逐字）；`pump_test.go:123` 钉的是**顶层四键**，卡内键不撞它。`cmd/wisp/**` 本腿⛔ 跑 ⇒ 那一面的红只交引用尺。
- (c) 拆完仍零读数：① 屏幕上那一行亮不亮——`frontend` 里 `blockedOnApproval` 命中文件数＝**0**（量的，射程 `frontend`），且 `panel_snapshot_deliver_329.go:27-33` 自陈页面侧 `window.wispOnPanelSnapshot` 接收器不存在；② 生产包里的 `blockedOnApproval=true` ——`cmd/wisp/subagent_carrier_197_test.go:598-604` 逐字写着 "blockedOnApproval is never asserted on a production packet"；③ 卡片任务号真被生产填上（今天零枚仪器读 `NativeVerdict` 的任务号）。
- (d) 新 JSON 键独立一问（量的，⛔ 推理）：今天钉键的名册只有三处——`internal/panel/approval_test.go:105`（ApprovalCardView/ResultChunk/Snapshot 三对）、`internal/panel/composer_test.go:52`（六对）、`internal/panel/config_313_test.go:283` 钉的是 **settings payload 那三组键**：顶层 11 枚 `chatModel…voice`、provider 行 10 枚 `apiKeyRef…protocol`、voice 5 枚 `cloudASRChain…realtimeProvider`（`:290/:317/:330` 逐字）⇒ **票 313 那族键钉碰不到卡片键**。撞不撞：`git grep -nE 'json:"(taskId|task|taskID)' 7535d83b -- internal frontend`＝4 命中，其中 `subagent_roster_197.go:107` 已用 `json:"taskId"`（名册行那一枚）⇒ 卡内再取 `taskId` 是**同名不同段**（不撞 Go 结构、不撞 TS：`frontend/src/lib/panel.ts` 里 `taskId` **零命中**，尺同射程）；要不要换名＝裁定的事，本腿不裁。C17 面：`docs/specs/SPEC-08-ui-ball-panel.md` 的 `:156`/`:235` 两处提到方法白名单与 `correlationId` 路由（票 328 引的号）⇒ 新键属契约面，要 `A##`。

**(乙) 改 `callCorr`，把可剥的数字放到尾巴**
- ★实测（量的，探针②`logs/form-candidates-probe.log`）：**票面给的那一例根本不点亮** —— `T#0#call_e1` 剥＝false；把号挪到最后 `T#call_e1#0` 剥＝true 但 base＝`T#call_e1`≠T；队列再撞一次 `T#3#7` base＝`T#3`≠T ⇒ **只有 `T#<纯数字>`（丢掉 provider 号）才点亮**，而那一形与队列的撞名改写（`queue.go:168` 用 `q.seq`）**同一命名空间**，"这串号是谁改的"从此不可读。
- (a) 要动：`internal/agent/loop.go:603-608`（两支）＋`:676` 那行；连带 `approval/queue.go:163-169`＋`:614-616` 的撞名/重放形、`:189-200 otherNames`（别名里带 `d.CorrelationID`）；join 侧一枚不动。
- (b) 顶红名册（红句逐字，用例颜色＝量的：五枚今天 `--- PASS`，见 `logs/agent-color-targeted2-7535d83b.log`）：`internal/agent/ticket285_corr_distinct_rulers_test.go:136`（`Test285LoopDispatchesOneCorrPerCall`）红句 `"call %s corr = %q routes to %q, want its own call id"`（`:137`）⇒ 尾形一改就红；`internal/agent/loop_golden_test.go:77-81`（`TestGoldenSingleToolCall`，量到 PASS）红句 `"call identity = %+v, want task %s and correlation id %s (task id prefix + this call's own id, never the bare task id)"`。不红（只钉前缀与不等）：`corr_percall_242_test.go:139/:142`、`tools/ticket283_corr_identity_rulers_test.go:222`、`tools/ticket285_corr_rows_rulers_test.go:75-83`、`tools/subagent_depth_fixture_287_test.go:140`、`tools/loop_approval_test.go:219-222` ⇒ **六枚今天全绿**（量的：tools 那五枚＝`logs/tools-color-targeted-7535d83b.log` 5/5；agent 那枚 `TestCorrPerCallTwoAsksSameTask242` 在 `logs/agent-color-targeted2-7535d83b.log`）。
- (c) 零读数：① "provider 号到底长什么形"仓内无仪器（anthropic/responses 的线值零钉，尺＝仓内 `toolu_` 0 命中）⇒ 乙改完后仍无法观测"真号会不会又混进尾巴"；② 屏幕上那一行（同甲 ①）；③ 生产包 `blockedOnApproval=true`（同甲 ②）。

**(丙) 名册那把尺改成按已知任务号前缀匹配**
- (a) 要动：`internal/panel/subagent_roster_197.go:277-288` 一枚尺＋`:272-276` 那段自陈＋`:257-261 indexWaitingKey` 文档；产码面一行不动。
- (b) 顶红名册（`TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt` 今天 1 顶层＋6 子测试全绿，量的：`RUN=11／PASS_TOP=4／FAIL_TOP=1（＝approval_test 那枚，与本案无关）／SKIP=0`）：子用例 `:133 {"a non-numeric suffix is not the queue's form", "child-1#retry", false}` 与 `:134 {"a task id that merely prefixes does not light", "child-1x", false}` 必反 ⇒ 红句逐字族 `assertBlocked` 的 `"row %q blockedOnApproval=%v，期望 %v； packets=%s"`（`:97-98`）。`:132 other-1#9` 仍 false（前缀不是 `child-1`），不红。
- (c) 零读数：`splitClashSuffix` 全仓只有 `:267` 一枚读者、**没有任何测试直呼它**（量的，尺＝`git grep -n splitClashSuffix 7535d83b -- internal cmd`＝def＋doc＋`:267`）⇒ 放宽后这两枚子用例就是唯一仪器；把它们改掉（放宽）之后，"匹配过宽会不会把两枚任务混起来"这一格**没有任何仪器**——同甲②③那一格之外新增的这一格必须具名写。
- ⚠ 票面已具名：丙是**放宽**，`:275-276` 那句 "which is a wider match than any code in this tree produces" 要重新裁（本腿量到它今天为真：名册能收到的 corr 里没有一形落在这把尺上，见第一节那张表）。

## 三、必答三题（⛔ 本腿盖章）

**① 除测试外有没有生产读者握 L2 卡片的 `CorrelationID`？** 有。尺＝`git grep -n '\.CorrelationID' 7535d83b -- internal cmd frontend`（射程＝`internal/**`、`cmd/**`、`frontend/**`；`_test.go` 与构造侧赋值剔除）。逐枚点名——**卡片那一枚**：`internal/panel/subagent_roster_197.go:214`（本票那一格）、`cmd/wisp/panel_pump.go:384`（台账摘要取前三枚 id）、`internal/panel/pump.go:101`→`approval.go:78`（`NativeVerdict`→卡片字段的搬运），前端三枚：`frontend/src/App.tsx:322`（React key）、`components/l2-approval-card.tsx:265`（答复路由 `requestApprovalResolution(view.correlationId,…)`）、`components/approval-screen.tsx:58`（显示）。**同一串号的其它握法（不是卡片类型，但改形状会一起动）**：`approval/queue.go:163-179`（`byID` 路由键）、`approval/gate.go:427/:445/:476`、`approval/replies.go:171-194`、`approval/report.go:97`、`memory/dao_toolcall.go:30`（入库列）、`cmd/wisp/resident_approval_windows.go:613/:729/:744`、`cmd/wisp/run.go:1458`。⇒ **零枚的那一格是"名册之外没有任何读者会因为它改形而点不亮行"**，不是"零枚读者"。

**② 名册对 L1 索引两枚键、对卡片只索引一枚——故意还是漏？** 件里逐字（⛔ 转述）：
- `subagent_roster_197.go:56-58`（★同一枚文件、离那把尺 219 行，写着自己认识生产那一形）：`// The entity leg files the deriving call's task id (TaskID(ctx),` / `// internal/tools/subagent_197.go:262), never its correlation id: since` / `// bd124b2a loop.go's callCorr mints one per call (<taskID>#<callID>, :603/:676).`
- `subagent_roster_197.go:217-218`：`// Both halves of one window: an L1 window is keyed by its correlation, and` / `// the row it blocks is named by its task.`
- `subagent_roster_197.go:184-186`：`// A host that pairs cards with a correlation of its own, with no task id in it, still shows those rows as not blocked; that remains a statement` / `// about the pairing, not a guess to paper over.`
- 卡片侧为什么只有一枚，件里唯一解释写在**测试**侧 `subagent_blocked_220_test.go:43-44`：`// correlation id is the only id the card carries (ApprovalCardView has no task` / `// field, and giving it one would widen the frozen card JSON keys).`
⇒ 读数支持"作者知道这一形不亮并写明不糊"（`:184-186`）＋"卡片没带任务号是被键钉住的原因"（`:43-44`）；**"这一枚键该不该索引"在件里没有定性注释**，故意／漏＝留给裁决腿，本腿不裁。

**③ 本票＝票 220 残余那一格（`AC#5`）吗？** 不是。尺一＝`sed -n '29p'` 票 220（逐字）：`- [ ] **AC#5 整包终态读数**：…`go test ./cmd/wisp ./internal/...` 跑到终态、逐名比红名集合（**`-run` 单跑不算终态**）…具名写清"这一格修完之后界面上还差哪一跳"` —— 那一格缺的是**一发台面世界＋终态颜色名册**。尺二＝台账 `A1006`（`docs/reports/pending-and-issues.md:18793`）逐字列了 220-v1 交回的三枚"归编排者"：**① L2 卡片在生产里点不亮行（已立成本票 328）**／② 七件套缺 `models/`／**③〔缺世界〕那一发整包终态名册归我裁** ⇒ 本票＝①，`AC#5`＝③，两枚是不同事项（同一腿同时交回的另两枚）。

## 四、交件自证（三把尺）

- 门禁四数（共享工作树，本腿 Go 产码**零字节**，只交 `.md`／`.log`／`.sh`）：逐枚 `rc=` 见 `logs/gates-328a1.log`。
- 颜色四数（导出台树，`-count=1 -timeout 30m`，`⛔ ./...`、⛔ `cmd/wisp`）：`panel-color-targeted`＝**RUN 11／PASS_TOP 4／FAIL_TOP 1／SKIP 0**（那枚 FAIL 是键钉用例，起手即红、非本腿造成，红句已具名）；`agent-color-targeted2`＝**5／5／0／0**；`tools-color-targeted`＝**5／5／0／0**；`panel-color-keynails`＝**2／1／1／0**。四数未混用、未相减。⛔ 终态读数（`-run` 单跑只算定点，票 328 `AC#4` 要的那发成对终态归落地腿）。
- 活实例 `wisp.exe` PID 21892 未动、未杀；`internal/**`／`cmd/**`／`docs/specs/**`／`PLAN.md`／台账／`HANDOVER`／`frontend/**`／`design/**`／三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／`ci.yml`＝**零字节**（交件侧 `git status --porcelain -- internal cmd docs/specs frontend design` 复跑为空，见本腿最终回报）。
- 件里每一枚数都带"尺＋射程目录＋量的／算的"；⛔ 本件不选形、不裁契约、不翻框。
