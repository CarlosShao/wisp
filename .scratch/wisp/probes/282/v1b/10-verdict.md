# 282-v1b 复核判决（票 282 五格 · 非实现者）

腿 `282-v1b`；起手锚 `6c90fa7f`（见 00-anchor.md）；盘上原文＝`.scratch/wisp/issues/282-corr-landing-over-scope-verdict-bridge-carries-taskid-and-the-new-accessor.md`。
现量审核（`git log -1`＝`f5c63b7c`，`date` 2026-10-08 20:0x CST）；分支 `dev`。
基线（五枚指名用例一发，rc=0）：`Test197SpawnPublishesIdentityRow` PASS 0.00s／`Test197ChildCannotDeriveSubagent` PASS／`Test221ParentStopsItsOwnChildRowAndStreamSettle` PASS 0.06s／`Test221SubagentCannotStopSiblingOrItself` PASS／`TestLoopPassesDeclaredL1WriteThroughTheGate` PASS 3.08s。

跑法：`cd internal/tools` ＋ `PATH="$PWD/../../third_party/sherpa-onnx:$PATH" go test -count=1 -timeout … .`（sherpa-onnx 一枚，层数现数 2）。每发突变四件套（种前 `git hash-object`／改行 `sed -n` 复量／红句带 file:line／还原后 hash 相等＋`git status --porcelain` 空）逐发内联在命令里；备份在仓外 `D:/work/workspace/projects plans/282v1b-tmp/`。

## 五格判语

| 格 | 判语 | 一句话凭据 |
|---|---|---|
| AC#1 种坏验牙 | **成立**（空形红；corr 回退形今天没有仪器，具名见下） | `TaskID` 种坏成 `""` ⇒ `subagent_197_test.go:302` 红句 rc=1；种坏成 `return h.corr` ⇒ 三枚全绿 rc=0＝该形没有仪器 |
| AC#2 没顺手放宽 | **成立** | 新句三形皆红（等值形/无前缀形红在 `loop_approval_test.go:221` 逐字；空形红在 `:204`）；旧句被推翻的等值形在新句下仍红 |
| AC#3 深度判定静默失效 | **判不动**（缺"per-call corr 形状的真 loop→bridge 整链"这把尺） | 改读前（两跳还原 CorrelationID）与改读后同一形状三枚用例双双全绿 rc=0，套件分不出；该支今天没有仪器 |
| AC#4 两条残余 | **成立**（两条均核实为真；均判**留**） | journal 行坑只在非生产组合（`loop.go:369` 收 taskID；生产 `run.go:1003` 是 `Journal: nil`＋`run.go:750` 桥记行）；bridge CRLF 是检出态（blob 全 LF、工作树 CRLF） |
| AC#5 越界检查 | **成立** | 三笔名册 9 件全在 `.scratch/**`＋`internal/{agent,tools}` 内；`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／`frontend/**`／`design/**`／冻结件一件未碰 |

## AC#1 突变逐发行

| 发 | 改哪行 | 读数 | 还原 |
|---|---|---|---|
| S1 空形 | `internal/tools/cancel.go:69` `return h.taskID` ⇒ `_ = h; return ""` | **rc=1 红**：`subagent_197_test.go:302: 派生失败：这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，就没法登记父子关系，也没法把结论盖戳进父任务的 C25 作用域）` | hash `87ac0624c689eca876b30c13a504639ee8b315d0` 逐字回；porcelain 空。同发内 221 腿**不红而等死**（`=== RUN Test221ParentStopsItsOwnChildRowAndStreamSettle` 后无果——子未起 ⇒ `spawnHeld` 等 `t.Context().Done()`，要 `go test` 10m 超时才收；本腿手工止损） |
| S1 附注 | 首试 `return ""`（h 未用）⇒ **build failed rc=1、0 条 `=== RUN`**，不算读数；已还原（hash 同上） | —— | 同上 |
| S2 corr 回退形 | 同 69 行 ⇒ `return h.corr` | **rc=0 全绿**（197 派生／197 深度／221 兄弟自停 三枚 PASS）＝corr 回退形今天没有仪器 | 还原 hash 逐字回；porcelain 空 |

S2 具名（为什么没有仪器）：今天所有经桥的 `task.spawn`／`task.cancel` 派发都用 `corr==taskID` 构造——`subagent_197_test.go:217`、`task_cancel_221_legs_test.go:119`、`subagent_222_test.go:332`、`cmd/wisp/subagent_carrier_197_test.go:173`、`cmd/wisp/subagent_blocked_197_test.go:121` 全是同值；仓内存在 `corr!=taskID` 的派发只有 `bridge_test.go:149`／`grant_test.go:265,279`／`ticket224_grantid_column_test.go:75`／`wiring_test.go:320`，全是 fs/t224 工具的夹具，打不到这两个读点。

## AC#2 三形逐发（新句＝`:219-220` `r.CorrelationID == "" || r.CorrelationID == res.TaskID || !strings.HasPrefix(r.CorrelationID, res.TaskID)`）

| 形 | 改哪行 | 读数 |
|---|---|---|
| 旧形等值 | `internal/agent/loop.go:676` `CorrelationID: callCorr(taskID, p.call.ID, i)` ⇒ `CorrelationID: taskID` | **rc=1 红**：`loop_approval_test.go:221: correlation_id = "4445a6fe-6029-4f9e-bc43-d13a451a3c6f", want a non-empty per-call id routing back to task "4445a6fe-6029-4f9e-bc43-d13a451a3c6f" (C18)`（corr 塌成 taskID，第二支命中） |
| 无前缀 | `loop.go:607` `return taskID + "#" + callID` ⇒ `return "call:" + callID` | **rc=1 红**：`loop_approval_test.go:221: correlation_id = "call:call_w1", want a non-empty per-call id routing back to task "346f7709-…" (C18)`（第三支命中） |
| 空串 | `internal/tools/bridge.go:1124` 落行处 `CorrelationID: orDefault(req.CorrelationID, req.TaskID),` ⇒ `CorrelationID: "",` | **rc=1 红**，但红在**更早一行**：`loop_approval_test.go:204: tool_call rows = 0, want exactly one per call: []`——空 corr 行根本没落库（`internal/memory/schema.go:82: correlation_id TEXT NOT NULL`）；钉子第一支（`==""`）在真组合里够不到，因为上游先挡两道：`bridge.go:269-270` 与 `:1124` 两层 `orDefault` 会把空上抬成 taskID。**形状仍红，但红色不落在那支上**（具名）。 |

逐形对比说明（"没放宽"逐形）：旧句 `=="" || != res.TaskID` 拦的是{空、≠taskID}；新句拦的是{空、==taskID、无前缀}。旧句拦下的"≠taskID 但带前缀"（即 `taskID#call_x`）在新句下放绿——但这一形正是票 242／`A745` 裁定要铸的形状（C18 原文在位 `docs/PLAN.md:1368`，corr 是答复路由键、与任务标识并列非等值），是新口径的**目标形**而非放宽；旧句原本拦的另一个形状（`==taskID`）在新句下仍红（S 表第一发读数）。三形皆红逐发在案 ⇒ 按本票判据"等价或更强"成立，判不退。

## AC#3 两发读数（同一形状）

- 改读前（`subagent_197.go:262` `parentID := TaskID(ctx)` ⇒ `CorrelationID(ctx)` ＋ `task.go:708` `caller := TaskID(ctx)` ⇒ `CorrelationID(ctx)`，同发还原）：`Test197ChildCannotDeriveSubagent` PASS／`Test221ParentStopsItsOwnChildRowAndStreamSettle` PASS／`Test221SubagentCannotStopSiblingOrItself` PASS，rc=0。
- 改读后（基线同三枚）：全绿，rc=0。
- 判决：同形同绿、无 Δ ⇒ 套件**分不出**这两跳的新旧读法；"深度判定静默失效已修"今天既不能让套件种坏就红、也没有一支会因此变绿的反例。缺的尺＝一把**真 loop（per-call corr）→桥→两读点**的整链尺（今天 197/221 全走 `bridge.Execute` 直供 corr==taskID；`cmd/wisp` 侧同值直供）。242 落地件 §4 自己也写明"无现成尺在真 loop→bridge 整链上钉这两处身份"。机制推理上修法是对的（生产形状 `loop.go:676` 铸 `taskID#call_x`、`bridge.go:528` 把 `req.TaskID` 样入 handle ⇒ `TaskID(ctx)` 取真任务 id），但这是推理，不是本票要的读数。⇒ **判不动（该支今天没有仪器）**。

## AC#4 两条残余逐条核实

- **journal 非生产组合行仍 taskID**：属实。现量 `internal/agent/loop.go:369` `j := newTaskJournal(l.opt.Journal, taskID, taskID)`（`journal.go:59` 签名 `newTaskJournal(j Journal, taskID, corrID string)`，两 id 都收 taskID）；生产组合 `cmd/wisp/run.go:1003` 循环侧 `Journal: nil`（`run.go:750` 桥侧 `Journal: mem`，桥在 `bridge.go:1124` 以 per-call corr 落行）⇒ 该残余只在"无审批层宿主/测试"组合可见；无断言钉该组合的 corr 值（`loop_approval_test.go:255-258` 只钉 Decision/RiskLevel）。**判留**：非生产组合、且 242 的口径要求是"请求侧 corr per-call"，该组合是循环自己记账的任务级行；若将来要行也 per-call，属新一波（落地件 §6 已自报"属下一波"，仓内无既有票承接，不新造票号）。
- **bridge.go 既有 CRLF 未洗**：属实且为**既有**。工作树 `tr -cd '\r' | wc -c`＝1284；`git show` 三处 blob（`bd124b2a^`／`bd124b2a`／`HEAD`）CR 计数**均为 0**；`git ls-files --eol internal/tools/bridge.go`＝`i/lf w/crlf`；`git config core.autocrlf`＝`true` ⇒ 工作树 CRLF 是检出态、不是本三笔产物（三笔 blob 全 LF）；porcelain 对 bridge.go 空。**判留**：既有行尾态，不顺手洗是对的（防 churn／防吞在飞改动）；全仓同族（台账记 13 枚）。

## AC#5 越界检查（三笔 `--name-only` 名册逐字）

- `133b1bfa`：`.scratch/wisp/probes/242/corrland1/00-anchor.md`（1 件）。
- `bd124b2a`：`internal/agent/corr_percall_242_test.go`／`internal/agent/loop.go`／`internal/tools/bridge.go`／`internal/tools/cancel.go`／`internal/tools/loop_approval_test.go`／`internal/tools/subagent_197.go`／`internal/tools/task.go`（7 件）。
- `0abe217c`：`.scratch/wisp/probes/242/corrland1/10-landing.md`（1 件）。
- 对禁区：`docs/PLAN.md` 0 件；`docs/specs/**` 0 件；`internal/observe/thresholds.go` 0 件；`frontend/**` 0 件；`design/**` 0 件；冻结件（SLO 阈值/golden/`thresholds.go`）0 件。判决：**成立，零越界面**。

## 本腿没试出来的

- S2/S3 两形（corr 回退、改读前还原）在今天的套件里**造不出红**——没有一把 corr!=taskID 形状下的真 loop→bridge 尺；本腿没有新造尺（写面不许）。
- 钉子第一支（`==""`）在真组合里**点不亮**：需要同时挖 `bridge.go:269-270`＋`:1124` 两层归一化＋memory 的 NOT NULL，三级突变没做（做了也是自造态的读数，价值低于成本）。
- S1 的 221 腿是"等死"不是红句（`spawnHeld` 等 `t.Context`），本腿以手工止损记之，未追它到 `go test` 10m 超时那一句。
- `cmd/wisp` 侧未跑（靶向亦然）：已知其 subagent 用例直供 `corr==taskID`（见上），不受本跳影响；在册红不在本腿地界。

## Commit 清单（本腿）

1. `6c90fa7f` 起手锚（00-anchor.md）。
2. 本件（10-verdict.md）随第二笔。
