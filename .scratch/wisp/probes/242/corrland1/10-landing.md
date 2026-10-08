# 242-corrland-1 · 10 落地件（每次工具调用铸一枚独立 corr）

锚：起手 `bbbc18fa`；本腿 commit `133b1bfa`（00-anchor）→ `bd124b2a`（落地，7 件）。⛔ 未 push、零翻框。

## 1. 改动逐处（工作树现量行号，同一发现取）

| 文件:行 | 改了什么 | rc |
|---|---|---|
| `internal/agent/loop.go:602-613` | 新增 `callCorr(taskID, callID, index)`：形状 `taskID + "#" + callID`（空 callID ⇒ `taskID#call-%d` 按回合内位次回退）；注释写明 corr 与 task 的关系＝C18 原文（队列项字段＋答复路由键、与任务标识并列非等值） | 编入下两行 rc=0 |
| `internal/agent/loop.go:676` | 调度那一跳 `CorrelationID: taskID` ⇒ `CorrelationID: callCorr(taskID, p.call.ID, i)`（唯一产码行为改动） | agent 新用例 rc=0；tools 钉子 rc=0 |
| `internal/agent/loop.go:363-369` | 原注释 `// C18: correlation == task id`（引用与原文不符候选）**具名重判**：改为"journal 两个 id 都取 TASK id（它记的是任务级行）；C18 **不**要求 corr==task；请求侧 corr 每次调用一枚、任务级溯源仍走 taskID" | 同上 |
| `internal/tools/loop_approval_test.go:215-222` | 那枚**唯一等值断言**具名重判（原句/新句逐字见 §2） | rc=0（3.08s 实跑，`=== RUN`/`--- PASS` 已眼见） |
| `internal/agent/corr_percall_242_test.go`（新） | 新用例 `TestCorrPerCallTwoAsksSameTask242`（同任务并发两问） | rc=0（`=== RUN`/`--- PASS` 已眼见） |
| `internal/tools/cancel.go:29-34,67-76` | `cancelHandle` 加 `taskID`；新增访问器 `TaskID(ctx)`（连带，见 §4） | tools 全靶向 rc=0 |
| `internal/tools/bridge.go:524-530` | 样章 `cancelHandle{..., taskID: req.TaskID, ...}`（连带） | 同上 |
| `internal/tools/subagent_197.go:259-262` | `parentID := CorrelationID(ctx)` ⇒ `TaskID(ctx)`（连带；行号 262＝现量） | tools 197/221/222 21 枚 rc=0 |
| `internal/tools/task.go:702-708` | `caller := CorrelationID(ctx)` ⇒ `TaskID(ctx)`＋注释同步改（连带） | 同上 |

## 2. 那枚钉子怎么重判的（原句/新句逐字）

原句（`loop_approval_test.go:213-214`，HEAD `bbbc18fa` 现量）：
```go
	if r.CorrelationID == "" || r.CorrelationID != res.TaskID {
		t.Errorf("correlation_id = %q, want the task id %q (C18)", r.CorrelationID, res.TaskID)
	}
```
新句（本腿 `bd124b2a`，现量 `:215-222`，含理由注释）：
```go
	// C18 makes correlationId the key the approval reply is routed by, a
	// queue-item field BESIDE the task id; it never required the two to be
	// equal (docs/PLAN.md:1368). Since ticket 242 the loop mints one
	// correlation id per tool call, traceable back to this task by its task-id
	// prefix, so the row must carry a non-empty per-call id that routes back
	// to this task - not the task id itself, and not an empty string.
	if r.CorrelationID == "" || r.CorrelationID == res.TaskID ||
		!strings.HasPrefix(r.CorrelationID, res.TaskID) {
		t.Errorf("correlation_id = %q, want a non-empty per-call id routing back to task %q (C18)",
			r.CorrelationID, res.TaskID)
	}
```
理由：保留"非空"＋"能路由回本任务"（前缀＝任务级可溯源），并**加一行 `== res.TaskID` 即红**——旧形（corr 塌成 taskID）仍被这枚尺抓住 ⇒ 强度不降反增（旧形、空串、乱串三形皆红）。⛔ 未删、未放宽。

## 3. 新用例（必交）

- 名：`TestCorrPerCallTwoAsksSameTask242`（`internal/agent/corr_percall_242_test.go`，包内）。
- 形状：照本包 `harness_test.go` 既有搭法（golden C5 `parallel-tools`＋`withTools` 薄夹具＝一枚 corr 一个放行闸的 ToolProvider，约 60 行）；同一任务一个回合两枚并发调用，**两问同时在场**（两次 `start` 尚在才继续）；断言：两 corr 非空且**互不相等**；按 c1 放行**只**唤醒 call_p1（且任务未结束——第二问仍被扣着）；再按 c2 放行唤醒 call_p2；`HasPrefix(c1/c2, res.TaskID)` 且 `≠ res.TaskID`（各自可路由回本任务、非塌回任务 id）。
- rc=0；`-v` 眼见 `=== RUN` / `--- PASS (0.01s)`。
- 形状今天**造得出**（C4 端口注入＋既有 golden，是本包既有用法），未触发停手条件。

## 4. 连带（越格自报）：tools 侧"任务级身份"改读 `TaskID(ctx)`

per-call corr 落地后，`CorrelationID(ctx)` 不再能当任务 id 用，而全仓只有两处产码把它读作任务身份（普查 D21/D22）：`subagent_197.go` 的 `parentID`（名册按 taskID 记账、`Look(parentID)` 做"子代理不许再派"深度判定、spawn 侧填 `ParentTaskID`）与 `task.go` 的 `caller`（`target==caller`、`rec.ParentTaskID != caller` 的父子比较）。若不动，同任务的 spawn 调用与 stop 调用各持一枚不同 corr ⇒ **深度判定会静默放行、父任务停不了自己的孩子**。故：bridge 在调用上下文样入 `taskID`，新增 `TaskID(ctx)`，两处读点改读它——新旧两制下取值都不变（旧制 corr==taskID；新制直读 req.TaskID），属**行为保持**而非语义变更。⛔ 未动 `CorrelationID(ctx)` 本身（一键仍是"这次调用"的 id，其访问器文案本就这么写）。
- 存量 197/221/222 全 21 枚靶向 rc=0（含 `Test197ChildCannotDeriveSubagent`、`Test221ParentStopsItsOwnChildRowAndStreamSettle`）；cmd/wisp 的 subagent 用例走 `bridge.Execute` 直供 `CorrelationID`（票面已记），不受本跳影响。
- ⚠ 无现成尺在"真 loop→bridge"整链上钉这两处身份（census D21/D22 只述机制）⇒ 本处的必须性今天靠机制推理；验方可对 `TaskID(ctx)` 做定向突变后走真 loop 顶它。

## 5. 三枚测试构造点逐枚重判（⛔ 未批量改）

- `internal/tools/pointer_183_cli_seam_test.go:127` **不改**：`r183Call` 是直供桥的 CLI 接缝夹具，`CorrelationID: taskID` 此刻是"宿主自选 corr"的合法旧形（C18 未要求同值），且它断言的是路径/CLI 行为、与 corr 取值无关。
- `internal/tools/pointer_185_cli_seam_test.go:120` **不改**：同形（`p185Call`），同上理由。
- `internal/tools/ticket175r2_stamp_live_test.go:95` **不改**：同形（`call175r2` 直供桥＋假 gate），该文件钉的是 C25 盖章与拒绝行，corr 不是被测物。

## 6. 没做的 / 残余

- `internal/tools/bridge.go` 工作树全文件 CRLF（`git ls-files --eol`＝`i/lf w/crlf`；blob 为 LF、gofmt 干净）⇒ `gofmt -l` 报它属**既有**行尾态（仓内 13 枚带 CR 同族），本腿未顺手洗行尾（防 churn/防"别人的在飞改动"噪声）；本腿 6 枚其余文件均 LF、gofmt 干净。
- journal 侧（loop 自持 Journal 的非生产组合）两个 id 仍取 TASK id ⇒ 该组合下 tool_call 行 correlation_id 记任务级；生产组合（`cmd/wisp/run.go:1003` `Journal: nil`，桥记行）行里是 per-call corr。若验方要"journal 行也 per-call"，属下一波。
- 未跑 `d22scan`／未跑整包（在册红不在本腿地界）；未跑 `cmd/wisp` 靶向（其 subagent 用例直供 corr，不受影响，已记）。
- 未 push；未加 `-done`；台账/票面一字未动。

## 7. Commit 清单

1. `133b1bfa` probe(242): corrland-1 起手锚
2. `bd124b2a` 242-corrland-1 落地（7 件：loop 铸 corr＋注释重判＋钉子重判＋新用例＋tools 连带 TaskID(ctx)）
3. 本件（10-landing.md）随下一笔 commit。
