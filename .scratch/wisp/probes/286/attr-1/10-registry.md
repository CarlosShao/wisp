# 票 286 AC#3 — 全仓 "corr 等于 taskID" 旧等值钉名册（只读归因腿 golden-attr-1）

取数时刻：2026-10-09 11:4x +0800（`date` 现量 11:41）。
读数一律 `git show HEAD:<path>` / `git grep ... HEAD --`；HEAD 现量 = `98482adc`（11:40，另一枚腿落的台账件）。
`internal/agent/loop.go` 自 `bd124b2a` 起未再被任何 commit 改；`internal/agent/loop_golden_test.go` 自 `62dda11d` 起未再被改 ⇒ 本名册的行号/形状对 `bd124b2a..HEAD` 全段一致。
尺覆盖（不止一种字样）：`CorrelationID <==|=> <task-token>` / `CorrelationID: <task-token>` / `CorrelationID = <task-token>` / 位置参数 `taskID, taskID` / `orDefault(Corr, Task)` / 注释里的 `TaskID == CorrelationID` 与 `CorrelationID == TaskID` 断言句 / 反形（两枚 corr 比相等）。

三档判语：**该搬**＝请求侧那枚被 242 打红的旧等值钉（写腿按 AC#2 搬）；**该留**＝这一侧今天/设计上就该等于 taskID 或本就是正确新尺；**待人裁**＝非红、非请求侧、且落在本票射程外（注释在别的包）。

## 一、相等断言族（corr 与 task 被比较）——7 枚

| # | HEAD:file:line | 该行原文 | 侧 | 判 | 为什么 |
|---|---|---|---|---|---|
| 1 | `HEAD:internal/agent/loop_golden_test.go:69` | `	if call.Req.TaskID != res.TaskID \|\| call.Req.CorrelationID != res.TaskID {` | 请求侧 | **该搬** | 唯一被真值打红的钉；产码 `loop.go:676 callCorr` 已铸 `taskID#call_e1`，永不等 taskID。AC#2 裁"搬期望、不回退产码"。 |
| 2 | `HEAD:internal/agent/loop_golden_test.go:341` | `	if row.Tool != "echo" \|\| row.CorrelationID != res.TaskID {` | 账本侧/journal | **该留** | 钉的是 DB `tool_call` 行的 corr；该行 corr 由 `loop.go:369 newTaskJournal(..., taskID, taskID)` → `journal.go:81 CorrelationID: t.corrID` 落成 taskID，今天就是绿的。票 282 AC#4(a) 已裁"留"（`A758`）。⛔ 不许为"看起来一致"去动它/`journal.go:81`。 |
| 3 | `HEAD:cmd/wisp/panel_pump_test.go:323` | `	if chunk.CorrelationID != f.taskID() {` | 结果流侧/results | **该留** | `chunk = finalSnap.Results[0]`；结果块的 corr = 流 key，而 run 用 **task id** 追加/收尾（`run.go:1266 c.stream.Append(e.TaskID, e.Text)`、`:1293 c.stream.Close(e.TaskID)`，落 `pump.go:526/549 ResultChunk{CorrelationID: key}`）。文件头 `:20` 就把这条写成"results[].correlationId vs the task id"；同一文件 `:221` 证明**审批卡**那一侧才带 per-call corr（`card.CorrelationID != "pump-corr"`）。⇒ 两个不同面，钉在 task-keyed 那一面，绿。**不是第二枚红。** |
| 4 | `HEAD:internal/agent/ticket285_corr_distinct_rulers_test.go:119` | `			if first.CorrelationID == res.TaskID \|\| second.CorrelationID == res.TaskID {` | 请求侧新尺 | **该留** | 反形正确尺：若任一 corr **塌回** taskID 就 Errorf（还配 `:132 HasPrefix`/`:136 suffix==CallID`）。这是 242 之后"必须不等"的那一族。 |
| 5 | `HEAD:internal/tools/loop_approval_test.go:219` | `	if r.CorrelationID == "" \|\| r.CorrelationID == res.TaskID \|\|` (+ `:220 !strings.HasPrefix(r.CorrelationID, res.TaskID)`) | 请求侧新尺 | **该留** | 就是 `bd124b2a` 把当年 `:213` 那枚等值钉搬成的形状（非空＋前缀＋≠taskID）。已收口。 |
| 6 | `HEAD:internal/tools/ticket283_corr_identity_rulers_test.go:222` | `		if r.CorrelationID == "" \|\| r.CorrelationID == res.TaskID \|\| !strings.HasPrefix(r.CorrelationID, res.TaskID) {` | 行侧新尺 | **该留** | 同 4/5 的正确形状。 |
| 7 | `HEAD:internal/tools/ticket285_corr_rows_rulers_test.go:72` | `		if r.CorrelationID == res.TaskID {`（→ `:73 Fatalf "塌在 task id 上"`）+ `:75 HasPrefix` | 行侧新尺 | **该留** | 要求行 corr **不等于** taskID 且以其为前缀；正确尺。 |

## 二、产码 corr==task（构造/赋值）——2 枚

| # | HEAD:file:line | 该行原文 | 判 | 为什么 |
|---|---|---|---|---|
| 8 | `HEAD:internal/tools/bridge.go:270` | `		req.CorrelationID = req.TaskID` | **该留** | 关在 `:269 if req.CorrelationID == "" {` 里＝**仅当 corr 空**才回落成 task；不是"要求等于"。票 283 注释 `:7` 也点明这枚回落。 |
| 9 | `HEAD:internal/agent/loop.go:369` | `	j := newTaskJournal(l.opt.Journal, taskID, taskID)` | **该留** | 账本侧残余：第二实参 corr 槽喂 taskID（→`journal.go:81`）。`loop.go:363-368` 注释（`bd124b2a` 改过）逐字承认"journal takes the TASK id for both ids；since 242 request-side is per-call"。票 282 AC#4(a) 裁"留"。 |

## 三、测试 fixture 构造 corr==task（宿主直接建 ToolRequest，不经 loop）——判 **该留**（合法输入、绿）；带过期注释的另列第五节

| # | HEAD:file:line | 该行原文（截） |
|---|---|---|
| 10 | `HEAD:cmd/wisp/subagent_blocked_197_test.go:121` | `				TaskID: id, CorrelationID: id, Tool: blocked197Tool,` |
| 11 | `HEAD:cmd/wisp/subagent_carrier_197_test.go:173` | `			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),` |
| 12 | `HEAD:cmd/wisp/task_scope_close_151_test.go:39` | `		TaskID: task, CorrelationID: task, CallID: "call-" + task,` |
| 13 | `HEAD:internal/tools/bridge_scope_open_ticket158_test.go:39` | `		TaskID: task, CorrelationID: task, CallID: "call-" + task,` |
| 14 | `HEAD:internal/tools/pointer_183_cli_seam_test.go:127` | `		TaskID: taskID, CorrelationID: taskID, CallID: "call-" + taskID + "-" + tool,` |
| 15 | `HEAD:internal/tools/pointer_185_cli_seam_test.go:120` | `		TaskID: taskID, CorrelationID: taskID, CallID: "call-" + taskID + "-" + tool,` |
| 16 | `HEAD:internal/tools/task_cancel_221_legs_test.go:119` | `		TaskID: callerID, CorrelationID: callerID, CallID: "call-221-cancel-" + callerID + "-" + targetID,` |
| 17 | `HEAD:internal/tools/ticket175r2_stamp_live_test.go:95` | `		TaskID: taskID, CorrelationID: taskID, CallID: "call-" + taskID + "-" + tool,` |
| 18 | `HEAD:internal/tools/failclosed_236_teeth_test.go:145` | `		TaskID: "caller-236r2", CorrelationID: "caller-236r2", CallID: "call-236r2-cancel-no-roster",` |
| 19 | `HEAD:internal/panel/pump_test.go:57`（+:58 同形） | `				{CorrelationID: "task-1", Text: "前半", Done: false},` |
| 20 | `HEAD:internal/agent/declared_l0_risk_179_test.go:204`（+:227、:248 同形） | `			j := newTaskJournal(store, "task-179-nil-"+tier, "task-179-nil-"+tier)` |

判语：这些是**测试自己**把某次调用的 corr 塌成 taskID（宿主不走 `agent.Loop`，直接 `bridge.Execute`），prod 允许 corr==task（只要求非空/回落），故**绿、且合法**。⛔ 不是本票射程：搬它们没有真值依据，且 `#14/#15/#16/#17/#18` 的 cancel/roster 尺正是靠"corr==task 塌成同值"来制造盲区，动它们会拆票 221/236/158 那几族的判据。`#19`（`pump_test`）是 `ResultChunk` 直构 fixture，key 用 task id 与产码一致。`#20`（`declared_l0_179`）三发都是**账本侧** `newTaskJournal` 同 `#9` 那一残余 ⇒ 该留。

## 四、FROZEN（out of action scope，只此一行具名，不转述其余）——1 枚

| # | HEAD:file:line | 备注 |
|---|---|---|
| 21 | `HEAD:internal/perm/ticket90_persist_test.go:123` | `TaskID: "t90-persist", CorrelationID: "t90-persist"` —— 三枚冻结件之一，⛔ 一字不动、不在名册动作射程。仅记它存在（git grep 顺带命中，未整读）。 |

## 五、注释里的过期等值 CLAIM（把"loop 派发 corr==task"当产码事实写——自 `bd124b2a` 起为假）——4 枚，判 **待人裁**

| # | HEAD:file:line | 该行原文 | 备注 |
|---|---|---|---|
| 22 | `HEAD:internal/panel/subagent_roster_197.go:57` | `	// loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),` | 双重过期：形状（现为 `callCorr`）＋行锚（647 已漂到 676）。是产码 `.go` 注释，⛔ 不属 AC#4 的"搬 golden:69"那一笔，改它属另一格。 |
| 23 | `HEAD:internal/panel/subagent_roster_197_test.go:468` | `//     CorrelationID == TaskID (internal/agent/loop.go:647); a host that used a` | 同 22，测试注释。 |
| 24 | `HEAD:internal/panel/subagent_blocked_220_test.go:155` | `	// The production pairing (the loop dispatches with CorrelationID == TaskID,` | 配对的 `:157` fixture 用 corr="child-1"/TaskID=""，测试本身仍绿；仅注释句过期。 |
| 25 | `HEAD:cmd/wisp/subagent_carrier_197_test.go:20` | `// bridge with the ids the agent loop itself uses (TaskID == CorrelationID,` | 描述 `:173` 那枚 fixture；把它当"loop 自己用的 id"是过期叙述（loop 现用 callCorr）。 |

判语：这 4 枚都不是红、都不是请求侧断言，且散在 `internal/panel`/`cmd/wisp`。本票 AC#4 只搬 `loop_golden_test.go:69`＋反形自证；要不要同批订正这 4 句注释（属"242 未收口账"的另一面）＝**待你裁**。⛔ 不建议写腿在 AC#4 那一笔里顺手改（越出 AC#5 射程＋跨包）。

## 六、非本族（同行含两 token 但不是 corr==task）——排除说明（顶"只搜一种字样"的假话形）

- `HEAD:internal/agent/approval/pending_read_test.go:67` `	if d.CorrelationID != want.CorrelationID || d.TaskID != want.TaskID {` —— 是"实际↔期望"逐字段比（corr 对 corr、task 对 task），**不是** corr 对 task。排除。
- `HEAD:internal/tools/bridge.go:528/536/600/782/1081/1124`、`internal/agent/approval/gate.go:281/521`、`report.go:95/97`、`queue.go:191`、`internal/tools/fs_write_test.go:752` 等 `orDefault[Text](CorrelationID, TaskID)` —— 是"有 corr 用 corr、没有才落 task"的**路由回落链**，与 C18 一致，非等值钉。排除。
- 各 `CorrelationID: corr`（distinct 变量/字面量，如 `approval_reply_201_test.go`、`run_mode101_test.go:196`、`fakes_test.go:446`、`window_read.go:73` 等）—— corr **≠** task 的正确形状，不属"旧等值钉"。排除。

## 七、.scratch 探针桶（整段 out of action scope）

`git grep CorrelationID..task HEAD -- '.scratch/**/*.go'` = **162 行命中**，跨约 29 个探针/突变件（`probes/151*/bridge.shipped.go`、`probes/151/ac1-probe-source.go:22` 的 `TaskID: task, CorrelationID: task`、`probes/158/**` 与 `probes/222/**` 的 bridge 突变副本、`probes/175/**` 等）。这些是死探针/突变台件，非在跑测试套，⛔ 不在本名册动作射程；只报枚数供你判"是否要清"。

## 汇总（真树）
- 该搬：**1**（`loop_golden_test.go:69`）
- 该留：**19**（`#2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20`）
- 待人裁：**4**（`#22,23,24,25` 过期注释）
- 冻结 out-of-scope：**1**（`#21` ticket90_persist_test.go:123）
- 探针桶（不算动作）：162 行 / ~29 件
