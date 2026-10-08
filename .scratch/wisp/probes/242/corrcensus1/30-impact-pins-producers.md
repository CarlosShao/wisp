# 242-corrcensus-1 · ③ 改动影响（机制，不裁定）、钉尺与产码点

锚：HEAD `8dbae00`。本件**只述机制**：不写"好/坏"，不写修法，不动任何产码（本腿 0 Go 命令、0 文件改写，写面仅本目录 `.md`）。
改动假设（他人待裁的那条）：`internal/agent/loop.go:654` 逐字 `TaskID: taskID, CorrelationID: taskID, CallID: p.call.ID,` 里第二项从"本任务 taskID"改成"每次调用唯一值"，其余不动。

---

## 1. 逐点机制（编号对应 ② 表 D 编号）

- D1/D2/D3/D4/D5（tools/bridge 记账链）：corr 全部经 `orDefault(req.CorrelationID, req.TaskID)` 取值。改动后 corr 非空且等于新唯一值 ⇒ `cancelHandle.corr`、`Complete`、`Vetoed`、`cancelText` 的键都从"任务 id"换成"调用 id"；同一条调用内五个点同源（都读同一 request/decision），彼此仍对齐。
- D6（gate.go:281-284 L1 窗口键）：窗口键从 taskID 变调用 id。`openWindow` 的查重语义"同一 correlation_id 已有窗口在跑 ⇒ fail-closed"从"同任务第二条 L1 会被拒"变"同一调用 id 才会撞"。
- D7/D8（Veto 查找）：`g.windows`/`g.running` 两表都是 corr 键；用旧地址（task id）提出的否决不再精确命中这两表，落到 D8 的 `q.reject(v.CorrelationID)` 通道；该通道按队列 byID 解析，而 byID 含别名（见第 2 节第 4 条）。
- D9（gate.go:516-522,552 重戳）：改动只动"入队前的 incoming 值"。`incoming` 与 `it.Corr`（队列地址，撞名时另发 `#seq`）的"两串结构"保留；`markStarted(incoming)` 与 bridge 的 `orDefault(req.CorrelationID, req.TaskID)` 仍同源对齐。
- D10/D11/D12（gate.go:727-748 原生/面板路由）：路由改按新值；`revokeGrants` 与 `reject` 命中的队列项不变（同一次请求同源）。
- D13/D14/D15（queue 发行）：主键来源值从 taskID 变调用 id；`#seq` 撞名条件由"同任务第二卡"变"同调用 id 重复登记"；`otherNames` 仍同时收 `d.CorrelationID` 与 `d.TaskID`（别名集合内容换一半）；重放发行 `it.Corr+"#replay"` 不变。
- D16-D20（replies 台账与路由）：`Record` 非空门、`byID` 键、allow/reject 的 `Request{corr}`、`Veto{corr}` 全部跟随新值；答复面必须携带展示时的值（replies.go:171 非空门是硬性入口）。
- D21（tools/subagent_197.go:259-263）：`CorrelationID(ctx)` 读的是 bridge.go:525 的 handle 值。改动后 `parentID` 从"父任务 id"变"派生调用的 id"；`Roster.Look(parentID)` 与深度判定照此值查找，命中与否取决于名册按什么键登记派生者。
- D22（tools/task.go:706 起）：`caller` 同 D21 换值；`target==caller` 与 `rec.ParentTaskID != caller` 的比较对象一并换成新值；注释（task.go:702-705）称两侧"by construction"同源，即 spawn 侧填 `ParentTaskID` 的来源若同步变化，比较仍成对。
- D23（panel roster join）：`indexWaitingKey` 收 card corr 与 w/task 两键；换值后卡侧键变，task 侧键（w.TaskID）不变，join 仍持有两半。
- D24/D25/D26/D27（resident/run 台账）：`Veto/Reject/Forget` 都以 card/prompt 自身携带的 corr 为键，随请求同源换值，机制上不产生交叉。
- D28（models.go:167）：只要求非空，换值不影响校验分支。
- D29（report.go:95-104）：`running[corr]` 查找随同源换值；报告字段（StartedBeforeVeto/Vetoed/Channel）按命中与否决定。

## 2. 结构性机制（跨点）

1. **键粒度**：以 corr 为键的台账（gate.windows/g.running、bridge.cancel、queue.byID 主键、replies.byID、roster waiting join、resident cards、panel 投影 correlationId）从"每任务一条"变"每调用一条"。
2. **同源对**：bridge.go:525/532/596/778、gate.go:521 的 `orDefault(...)` 与 gate.go:552 的 `markStarted(incoming)` 在同一请求内读同一个值，改动后仍成对；gate.go:517-520 注释明确"the two are not the same string once the queue has re-issued one"——队列重发地址（it.Corr）与本改动无关，两串结构保留。
3. **空值回退家族**：`orDefault`/`orDefaultText` 只在 corr 为空时回退 taskID（bridge.go:269-270、gate.go:281/521、queue.go:164、report.go:95、bridge.go:1077/1120）；corr 恒非空后这些点全部改用新值。
4. **taskID 为锚的旧地址**：queue.go:191 仍把 `d.TaskID` 收进 `names`（别名集合；注释："a card that is on screen is addressable by whatever the display was built from"）；push 后 `q.byID[corr]=it; q.indexLocked(it)`（queue.go:180-181，函数体未展开）是把主键与别名接入寻址的地方。不经队列的精确查找（gate.Veto 的 windows/running）只按精确值命中。
5. **`CorrelationID(ctx)` 语义**：两处（D21 subagent parentID、D22 task caller）今天读作"任务 id"，依据是产码把 corr 派成 taskID（loop.go:654、loop.go:363；注释原话 internal/panel/subagent_roster_197.go:57 "on the production path this reads as the parent's task id"）；改后两者读作"调用 id"。cancel.go:53-54 的访问器文案本身是"the id of the call the running tool belongs to"。
6. **持久行与读模型**：`tool_call.correlation_id` 由 bridge.go:1120（`orDefault`）与 journal.go:81（`t.corrID`，来源 loop.go:363）写入；该列取值随两处来源同步变化；列有索引 idx_tool_call_corr（SPEC-02:89），校验只要求非空（models.go:167）。
7. **面板投影**：panel_pump.go:66、panel/approval.go:78、ResultChunk/json 字段承载新值；答复回传须与展示值一致（replies.go:171 非空门 + byID 键）。

## 3. 钉尺：有没有测试逐字钉住"两枚同值"

- 尺一：`git grep -nF 'CorrelationID: taskID' HEAD -- '*_test.go'` **rc=0**，3 命中（均为测试载具构造 ToolRequest 时同填 taskID）：
  - internal/tools/pointer_183_cli_seam_test.go:127
  - internal/tools/pointer_185_cli_seam_test.go:120
  - internal/tools/ticket175r2_stamp_live_test.go:95
- 尺二：`git grep -nE 'CorrelationID.*==.*TaskID|TaskID.*==.*CorrelationID' HEAD -- '*.go'` **rc=0**，7 命中：
  - 测试断言 1 处（逐字钉住）：internal/tools/loop_approval_test.go:213-214 `if r.CorrelationID == "" || r.CorrelationID != res.TaskID {` `t.Errorf("correlation_id = %q, want the task id %q (C18)", ...)` —— 钉的是**落库行**（经产码 loop 走到 tool_call 行）corr == taskID；
  - 产码注释 1 处：internal/panel/subagent_roster_197.go:57；
  - 测试注释 3 处：cmd/wisp/subagent_carrier_197_test.go:20、internal/panel/subagent_blocked_220_test.go:155、internal/panel/subagent_roster_197_test.go:468；
  - `.scratch` 历史拷贝 2 处（probes/197/r3b，非现役）。
- 反证（已拆）：cmd/wisp/subagent_selfapproval_197_test.go:116 逐字 `TaskID: taskID, CorrelationID: corrID,`（两变量，票 259 AC#4 的拆分形状）。
- 分母：`git grep -ln 'CorrelationID' HEAD -- '*_test.go'` rc=0，80 文件（62 仓内 ＋ 18 枚 `.scratch` 拷贝；未逐枚判读，仅登记）。

## 4. 产码点全量：除 loop.go:654 外，还给 CorrelationID 赋值的地方（非测试）

**新造/改变取值（12 组）**
1. internal/agent/loop.go:363 —— `newTaskJournal(l.opt.Journal, taskID, taskID) // C18: correlation == task id`（journal 的 corrID 直取 taskID；落进 tool_call 行的就是它）
2. internal/tools/bridge.go:270 —— 空则 `req.CorrelationID = req.TaskID`（回退默认）
3. internal/tools/bridge.go:1120 —— `tool_call` 行 `CorrelationID: orDefault(req.CorrelationID, req.TaskID)`
4. internal/agent/approval/gate.go:522 —— `d.CorrelationID = corr`（重戳为队列地址；corr 来源 gate.go:516 `it.Corr`）
5. internal/agent/approval/gate.go:281 —— `corr := orDefaultText(d.CorrelationID, d.TaskID)`（L1 窗口键取值）
6. internal/agent/approval/queue.go:165 —— 空 ⇒ `"approval-"+seq`
7. internal/agent/approval/queue.go:168 —— 撞名 ⇒ `原值#seq`
8. internal/agent/approval/queue.go:614-616 —— 重放 ⇒ `it.Corr+"#replay"`（撞名 ⇒ `#seq`）
9. internal/panel/composer_handlers.go:122 —— `ModeRequest{To: req.To, CorrelationID: req.RequestID}`（面板 requestId 充当）
10. internal/panel/pump.go:526,549 —— `ResultChunk{CorrelationID: key, ...}`（key 由 Append/Close 调用方给；SubagentStreamKey 一族）
11. cmd/wisp/panel_assets.go:69 —— 字面常量 `CorrelationID: "panel-assets-l2"`
12. internal/agent/approval/report.go:95-97 —— `orDefaultText` 进 CancellationReport

**誊抄（10 组）**
13. internal/tools/bridge.go:375 —— `Decision.CorrelationID = req.CorrelationID`（誊抄）
14. internal/tools/bridge.go:1077 —— 日志行 `orDefault(...)`（誊抄）
15. internal/agent/approval/gate.go:308/322/341/437/554/561/573 —— Event 誊抄；gate.go:599 —— Prompt 誊抄
16. internal/agent/approval/queue.go:572、pending_read.go:125、window_read.go:73 —— 读模型誊抄 it.Corr/w.corr
17. internal/agent/approval/replies.go:175（+186/194）—— ReplyCard 台账誊抄；replies.go:329/405/438/463 —— 答复面入参誊进 Request/Veto
18. internal/agent/journal.go:81 —— `CorrelationID: t.corrID`（来源第 1 组）
19. internal/panel/pump.go:101 —— NativeVerdict→ApprovalSubject 誊抄
20. cmd/wisp/panel_pump.go:66 —— 队列视图→NativeVerdict 誊抄
21. internal/panel/approval.go:78 —— ApprovalCardView 誊抄

（其余命中为定义/注释/读点，不属赋值点；resident 与 run 侧全部是读 card/prompt 已有值。）

## 5. 未展开/未定（如实登记）

- `q.indexLocked`（queue.go:181）函数体未逐行读：别名如何进 byID 只据 36-41 注释与 186-198 行为推断，未取到断言。
- `task.spawn` 侧 `ParentTaskID` 的填值点未定位（仅据 task.go:702-705 注释称同源）。
- 80 枚含 CorrelationID 的测试文件未逐枚判读其是否依赖"两值同源"；本程只交付两条指定钉尺的全部命中。
- 未碰 D43 状态机转移表、未碰 C17 方法名册（本腿无该面动作；有冲突应停手上报的那条未触发）。
