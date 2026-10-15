# 165-r1 读数与面核（票 165 · plan.present）

- 起手锚（`git rev-parse --short HEAD` 现读，16:15:04）：`b74b87d0`
- 落件时锚（16:25:58 现读）：`9b0f7d55`（漂移：`b74b87d0` → `9b0f7d55`，另一腿的提交）
- 写面：`internal/tools/plan_present.go`（新）+ `internal/tools/plan_present_165_test.go`（新）
- 环境：`PATH=third_party/sherpa-onnx`；`GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\wisp165r1gotmp`

## 五格读数

| 格 | 读数 | 证据 |
|---|---|---|
| AC#1 | 未修码名册 10 枚、**零枚** plan/propose 形：`fs.read fs.list fs.write fs.edit fs.trash fs.move fs.delete task.output task.cancel task.spawn` | `Test165UnmodifiedRosterHasNoPlanTool`（只读三个既有 builder，本票新增的 plan 族进不去）；票面 `PLAN.md:2567` 有行、码侧零枚 ⇒ 票成立 |
| AC#2 | 复用形状成立：`Declared=L2` → R1 给 L2 → 走 `<b>Gate.PendingApproval</b>`（D43 #17 的 AwaitingApproval 支），**未**新增状态/通道/模式；`Gate` 仍**2** 个方法、`GrantSource` 仍 1 个只读方法 | `Test165PlanPresentAsksThroughTheExistingL2Card`（approvals=1 / windows=0 / RulesHit 含 R1 / ModeSilenced=false）；`04-baseline-run.txt` 5 条 PASS |
| AC#3① | 清单批过之后，清单外的 `fs.write`（`not-in-plan.txt`）**仍响**：L1 窗口被问 1 次、会话授权行 0 枚、decision 非 `allow_session_grant` | `Test165PlanApprovalRingsAgainForAnUnlistedWrite`；变异 M3 逐字见下 |
| AC#3② | "批过的范围"那一味**按批大小写死**：文本必须含「只限上面这 3 步」「清单之外」；空清单被判 IsError | `Test165ScopeClauseNamesTheBatch`；变异 M2 逐字见下 |
| AC#4 | 只加发起请求的出口：`Decl` 无 allow/approve/grant 字段、参数 schema 无 allow/approve/grant/permission/scope/auto 名（阳性对照：属性非空）、`GrantSource` 仍只读 1 方法、`AppliedSteps=nil && Truncated=false` | `Test165PlanPresentAddsNoAllowExit` |
| AC#5 | 本票字节面＝`internal/tools/plan_present.go` + `internal/tools/plan_present_165_test.go` + `.scratch/wisp/probes/165/**`；`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`internal/agent/approval/**`、`internal/panel/**`、`frontend/**`、`design/**`、阈值/golden/`allowlist.txt`/`tools/d22scan/**` 本票零字节（提交后 `git show --name-only` 为证） | 见提交 |

## 三发变异（仓外备份 + sha256 还原证）

基线 sha256（`01-cap-baseline.sha256`）：

```
3fc878959be94f4c41bd724cf3433ddc0b90c039d00781cc58be27b93f9a37e1 *internal/tools/plan_present.go
a0c0736c639276b0835be7c0f85df9dff517eef73039ee882842759d63accfdd *internal/tools/bridge.go
```

### M2（AC#3②）摘掉 `b.WriteString(planPresentScopeClause(len(steps)))` ⇒ 唯一红

```
--- FAIL: Test165ScopeClauseNamesTheBatch (0.00s)
    plan_present_165_test.go:263: 方案文本里读不到批准范围的边界（缺 "只限上面这 3 步"）：人类点头时看不到「批到哪儿为止」："这批改动要达成：165 探针\n我要做的事（共 3 步）：\n1. 甲\n2. 乙\n3. 丙\n"
    plan_present_165_test.go:263: 方案文本里读不到批准范围的边界（缺 "清单之外"）：人类点头时看不到「批到哪儿为止」："这批改动要达成：165 探针\n我要做的事（共 3 步）：\n1. 甲\n2. 乙\n3. 丙\n"
FAIL	github.com/CarlosShao/wisp/internal/tools	0.062s
```

还原：`3fc878959be94f4c41bd724cf3433ddc0b90c039d00781cc58be27b93f9a37e1`（与基线逐位相等）

### M3（AC#3①）把 L1 那一支的 `b.gate.PendingWindow(ctx, *dec)` 换成 `AnswerAllow, ""`（窗口不再问）⇒ 唯一红

```
--- FAIL: Test165PlanApprovalRingsAgainForAnUnlistedWrite (0.02s)
    plan_present_165_test.go:236: 方案批过之后，清单外的写没有响（L1 窗口被问 0 次，要 1 次）——批准变成了对后续写的放行
FAIL	github.com/CarlosShao/wisp/internal/tools	0.124s
```

还原：`a0c0736c639276b0835be7c0f85df9dff517eef73039ee882842759d63accfdd`（与基线逐位相等）

### M1（AC#2 附证）`Declared: risk.L2` → `risk.L1` ⇒ 两条红（清单不再走卡、只走 3 秒静默窗口）

```
--- FAIL: Test165PlanPresentAsksThroughTheExistingL2Card (0.00s)
    plan_present_165_test.go:161: Declared=L1，冻结文本（PLAN.md:2567 / SPEC-07 §3:73）写的是 L2
    plan_present_165_test.go:179: 清单走了 L1 静默窗口 1 次：这一档要的是卡（3 秒无人应答不是点头）
    plan_present_165_test.go:182: L2 卡被问了 0 次，要 1 次
--- FAIL: Test165PlanApprovalRingsAgainForAnUnlistedWrite (0.03s)
    plan_present_165_test.go:233: L2 卡被问了 0 次，要 1 次（只有清单那一发）
FAIL	github.com/CarlosShao/wisp/internal/tools	0.077s
```

还原：`3fc8789...` 同上（逐位相等）

## 具名没做（不假装）

- **没接线到 `cmd/wisp`**：本票写面是 `internal/tools/**`，`cmd/wisp` 在飞（135-r1）。`BuiltinPlanEntries()` 已备好、可一行注册，但生产装配点未加 ⇒ 今天生产链路里 `plan.present` 还是不存在的门（票 161 那条"没人喂的门"同样适用于装配半格）。
- **没做界面**（票面"本票不解决"）、没动台账/HANDOVER、没翻 AC 框。
- 方案文本的**条数上限**未定（契约没说），本件不发明一个数字。
