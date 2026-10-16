# 165-v1 读数件 01：变异重放（我自己物理改字节，八发＋一枚瞎尺读数＋两枚去牙正控）

台面＝`D:\tmp\165v1\t1`（`git archive HEAD \| tar -x`，1091 枚 `.go`），⛔ `-overlay`、⛔ 共享工作树、⛔ 仓内 worktree。
每发的记账形状＝**基线 sha256 → 突变 sha256 → 红句逐字 → 还原 sha256 →（还原后再跑一遍确认回绿）**。
全部跑完之后我又把这六枚被改过的文件逐枚对 `git cat-file -p HEAD:<f> \| sha256sum` 对账，六枚全 MATCH（见读数件 00 §1）。

## 0. 改前基线（t1，⛔ 用交付方的读数代替）

- `go test ./internal/tools -run Test165 -v -count=1`：**rc=0、RUN=5、PASS_TOP=5、FAIL=0、SKIP=0**
  （五枚＝`Test165UnmodifiedRosterHasNoPlanTool`／`...AsksThroughTheExistingL2Card`／`...RingsAgainForAnUnlistedWrite`／`Test165ScopeClauseNamesTheBatch`／`...AddsNoAllowExit`）。
- `go test ./cmd/wisp -run Test165r2 -v -count=1`（**靶向 `-run`，⛔ 整包**）：**rc=0、RUN=3、PASS_TOP=3、FAIL=0**，
  关键 Logf 逐字：`装配可达读数：registry 10 枚含 plan.present；卡=plan.present；行 risk=L2 decision=allow；授权行 0 枚`、
  `AC#3① 装配读数：方案卡=plan.present（yes）→ 清单外卡=fs.write（no，IsError=true）→ 授权行 0 枚 → 目标文件未落盘`、
  `阳性对照成立：cmd/wisp/run.go:783 range=true register=true`、`接线读数：cmd/wisp/run.go:822 range=true register=true ⇒ 产码调用点 1 枚，形状与先例同形`。

基线哈希：`plan_present.go`=`3fc878959be94f4c41bd724cf3433ddc0b90c039d00781cc58be27b93f9a37e1`（与 r1 的 `01-cap-baseline.sha256` 逐位相等⇒我复放的是**同一份字节**）、
`bridge.go`=`a0c0736c639276b0835be7c0f85df9dff517eef73039ee882842759d63accfdd`（同上相等）、`run.go`=`7f28db4fa8c47147d1a2a2110c1387db1ba7d898376d941ac61f579126ee6eb9`、
`tool.go`=`602d632b9a0a5d52bfb0133d89d8b366c3bef1dc1ba952aeacf0f389fbd03d63`、`plan_present_165_test.go`=`7fa7645adc50bf5f5a310001bda5f4a2ae8bd75f3379924db47144d53c768ebd`、
`plan_assembled_165r2_test.go`=`7b5f4e2291e14c78efef522032d6e43f21e0e00f06e1a18209147be1dd85bfc6`。

## 1. AC#3② 那一发：摘掉「批过的范围」这一味

### M2＝重放 r1 的 M2（摘整句范围条款）

改法＝`sed -i '172d'`（删掉 `b.WriteString(planPresentScopeClause(len(steps)))` 那一行）。
突变=`c7d2a7546cb6b7f1126af930f0a9797456434ccf5c66e12d8dea0c6735bf2e84`；`-run Test165` rc=1、**PASS_TOP=4／FAIL_TOP=1**，唯一红＝`Test165ScopeClauseNamesTheBatch`：

```
    plan_present_165_test.go:263: 方案文本里读不到批准范围的边界（缺 "只限上面这 3 步"）：人类点头时看不到「批到哪儿为止」："这批改动要达成：165 探针\n我要做的事（共 3 步）：\n1. 甲\n2. 乙\n3. 丙\n"
    plan_present_165_test.go:263: 方案文本里读不到批准范围的边界（缺 "清单之外"）：人类点头时看不到「批到哪儿为止」："这批改动要达成：165 探针\n我要做的事（共 3 步）：\n1. 甲\n2. 乙\n3. 丙\n"
FAIL	github.com/CarlosShao/wisp/internal/tools	0.185s
```

⇒ **r1 报的「M2 摘范围句 ⇒ 唯一红 `Test165ScopeClauseNamesTheBatch`」复现成立**（两行 Errorf 是同一枚顶层用例的两条子判据，⛔ 记成两枚红）。
还原=`3fc878959b…`（MATCH）。

### M2b＝同一味的**另一种坏法**（边界句还在、数不跟批）

改法＝`planPresentScopeClause(len(steps))` → `planPresentScopeClause(1)`（批 3 步却写「只限上面这 1 步」）。
突变=`54b067a6901a6e62516cf047e13bdbdad7e8a5268c32be32921b92d2f094a4ce`；rc=1、FAIL_TOP=1（仍唯一红），逐字（**文本里「清单之外」四个字确实在**）：

```
    plan_present_165_test.go:263: 方案文本里读不到批准范围的边界（缺 "只限上面这 3 步"）：人类点头时看不到「批到哪儿为止」："这批改动要达成：165 探针\n我要做的事（共 3 步）：\n1. 甲\n2. 乙\n3. 丙\n批准范围只限上面这 1 步、只限这一批：清单之外的任何写（包括随后想到的新步骤、以及之后的每一条命令）都仍会各自弹卡，不会因为这份清单被点头而顺带放行。"
```

⇒ 这枚尺**不是**只认"句子在不在"，它把句子钉在批大小上（派单问的"一把尺两种形状是否都放行"：这两形都不放行）。还原=MATCH。

## 2. AC#3① 那一发：批完方案后再来一发清单外的写必须响

### M3＝重放 r1 的 M3（L1 支不再问人）

改法＝`internal/tools/bridge.go:451` 的 `a, why := b.gate.PendingWindow(ctx, *dec)` 换成 `a, why := AnswerAllow, "v1m3 window not asked"`。
突变=`10bf7cccf37f78004758f8e57308d4bce47b3a709c067c6e84f46be14e2bb066`；rc=1、FAIL_TOP=1，唯一红＝`Test165PlanApprovalRingsAgainForAnUnlistedWrite`：

```
    plan_present_165_test.go:236: 方案批过之后，清单外的写没有响（L1 窗口被问 0 次，要 1 次）——批准变成了对后续写的放行
FAIL	github.com/CarlosShao/wisp/internal/tools	0.173s
```

⇒ **r1 报的「M3 ⇒ 唯一红 `Test165PlanApprovalRingsAgainForAnUnlistedWrite`」复现成立**。还原=`a0c0736c639276b0835be7c0f85df9dff517eef73039ee882842759d63accfdd`（MATCH）。

### M9＝那一发的**另一种坏法**（授权行真被铸出来，L1 支被 cover 掉）——这一形**谁都还没造过世界**

我在 t1 的 `plan_present_165_test.go` 里、清单外那一发写之前种了一行假授权：
`led.rows = append(led.rows, "fs.write")`（`plan165Ledger.Covering` 因此返回 `(1,true)` ⇒ 桥走 `agent.DecisionAllowGrant` 静默放行）。
突变=`1865b5ae927d213e9e8e783aa64504daa2278203e544d02dd253d649d148b174`；rc=1、FAIL_TOP=1，两条子判据逐字：

```
    plan_present_165_test.go:237: 方案批过之后，清单外的写没有响（L1 窗口被问 0 次，要 1 次）——批准变成了对后续写的放行
    plan_present_165_test.go:241: 批准一份清单之后，会话授权里多出了 [fs.write]：批准只批那一发，不许铸出授权行
```

⇒ 这一形**尺有牙**（交付件里那两条断言都会红），但**交付方从未造过这个世界**（r1/r2 的读数件里零记录）⇒ 本腿替它补上并记在案。
还原=`7fa7645adc50bf5f5a310001bda5f4a2ae8bd75f3379924db47144d53c768ebd`（MATCH）。

### M8＝装配那一半的同一形（真往库里种一行）

改法＝t1 的 `plan_assembled_165r2_test.go` 里，方案那一发答完 yes 之后、清单外写之前，用真 API 落一行：
`h.openStore().InsertGrant(ctx, memory.ApprovalGrant{Scope: GrantScopeSession, Tool: "fs.write", Pattern: ToSlash(h.outside)+"/*", SessionID: minted, …})`。
突变=`b7d715036ea6c028e04cf8cf280eed14df6f702b6b576d1c51f1feda668b848a`；靶向 rc=1，红位与读数逐字：

```
    plan_assembled_165r2_test.go:267: AC#3① 装配读数：方案卡=plan.present（yes）→ 清单外卡=fs.write（no，IsError=true）→ 授权行 1 枚 → 目标文件未落盘
--- FAIL: Test165r2UnlistedWriteRingsItsOwnCardAfterAnApprovedPlan (1.20s)
```

⇒ 装配那半**也会红**，红在「授权行 0 枚」那一枚断言上；⚠ 具名一处形状差：这一发里**卡照样升起**（清单外那发在这套夹具里是 L2 支，`Covering` 只短路 L1 支），
所以装配用例对这个坏世界的凭据是**授权行计数**那一枚、⛔ 是"卡没升"那一枚。还原=`7b5f4e2291e14c78efef522032d6e43f21e0e00f06e1a18209147be1dd85bfc6`（MATCH）。

## 3. AC#2 附证那一发＋装配那一半的去牙

### M1＝重放 r1 的 M1（`Declared: risk.L2` → `risk.L1`，`plan_present.go:127`）

突变=`a5ad06734ad6b6a3c51b4de0fe43ce2fa813da8f7906a5f7bf1b0383d8f8cf73`；rc=1、**FAIL_TOP=2**，逐字：

```
    plan_present_165_test.go:161: Declared=L1，冻结文本（PLAN.md:2567 / SPEC-07 §3:73）写的是 L2
    plan_present_165_test.go:179: 清单走了 L1 静默窗口 1 次：这一档要的是卡（3 秒无人应答不是点头）
    plan_present_165_test.go:182: L2 卡被问了 0 次，要 1 次
--- FAIL: Test165PlanPresentAsksThroughTheExistingL2Card (0.00s)
    plan_present_165_test.go:233: L2 卡被问了 0 次，要 1 次（只有清单那一发）
--- FAIL: Test165PlanApprovalRingsAgainForAnUnlistedWrite (0.02s)
FAIL	github.com/CarlosShao/wisp/internal/tools	0.077s
```

⇒ 读数复现（与 r1 自己的 `05-m1-level-downgrade.txt` 一致）；⚠ 但 r1 的**票尾 bullet** 写的是「M1 Declared L2→L1 ⇒ **AC#2 那条红**（附证）」＝把**两枚顶层红**报成一枚。
这是自陈与盘上读数的**一处不符**，判语里记为过期描述（不影响 AC#2 成立，影响的是那份自陈的可信度刻度）。还原=MATCH。

### MA＝重放 r2 的摘线发（`cmd/wisp/run.go` 删 822–827 那 6 行注册循环）

突变=`a85118e9013bd48c6a96dac168c6cd14c286a878bb96b79ca4de4f5f6fd69aa9`；靶向 `-run Test165r2` **rc=1、RUN=3、PASS_TOP=0、FAIL_TOP=3**，三枚红句逐字：

```
    plan_assembled_165r2_test.go:147: 接线 RED：生产装配的 registry 里没有 plan.present；名册 9 枚：[fs.edit fs.list fs.move fs.read fs.trash fs.write task.cancel task.output task.spawn]（这正是开工前 AC#1 的未修码形状：这份清单在 wisp run 里够不到）
    plan_assembled_165r2_test.go:220: 接线 RED：批方案那一发没升起卡（回文="未知工具 plan.present，可用工具见 list_tools"）⇒ 装配里没有这一发，c) 的后续读不到
    plan_wiring_165r2_test.go:191: 接线 RED：tools.BuiltinPlanEntries 在 internal/ 与 cmd/ 的产码里调用点 0 枚，要 1 枚（cmd/wisp/run.go 的装配那一处）；命中清单＝[]（0 枚＝写了没接线那一族，plan.present 在生产装配里根本够不到）
```

⇒ **r2／编排者报的「摘掉那一行三枚全红」复现成立**（另两枚各红 11s＝真走到装配那一跳才红，⛔ 红在某句文案上）。
还原=`7f28db4f…`（MATCH），还原后靶向重跑 **rc=0**（`logs/ma-restore-green.log`）。

## 4. AC#4 那两枚词面尺的牙（＋一枚**绿的**突变＝射程边界）

- **M6（nested，rc=0＝放行）**：把 `"approved":{"type":"boolean"}` 那一味塞在 `steps` 属性**内部**（改的是 schema 字符串第 68 行，`diff` 回读证改动真在字节上）。
  突变=`7c1ff2c6fc6430cb8de9a3ed9fe94349a42b67aaed2ad887554e2fe6145f5b3d`；`-run Test165` **rc=0、五枚全 PASS**。
  ⇒ 具名射程：`Test165PlanPresentAddsNoAllowExit` 的参数名册只扫**顶层 `properties` 的键**，嵌套层的键它读不到。
  ⛔ 把这枚 0 读成"面板代答做得到"：这一形在码侧同样够不到任何出口——`planPresentArgs` 只有 `Summary`／`Steps` 两个字段（`sed -n '/^type planPresentArgs struct/,/^}/p'` 现读），
  且 schema 顶层写死 `"additionalProperties":false`，所以"嵌套里一个叫 approved 的键"既不是模型能设的顶层参数、也没有任何 Go 字段接它。⇒ 记为**尺的射程边界＋被另两味兜住**，不记为漏洞。
- **M6b（真·顶层参数，rc=1）**：把 `"approved":{"type":"boolean"}` 放进顶层 `properties`。
  突变=`9ae3bc41594ce12b126ad108570a2980818d99bc8b6275ae522d9b0daf4308c2`；rc=1、唯一红＝`Test165PlanPresentAddsNoAllowExit`：
  `    plan_present_165_test.go:309: 参数 "approved" 让模型自己说「这批已获准」：本票只许加发起请求的出口`。还原=MATCH。
- **M7（声明面，rc=1）**：给 `internal/tools/tool.go` 的 `type Decl struct` 加一枚 `AutoApprove bool`。
  突变=`da82df7647e0633223826551ab7bc64f7cbc22850e1c5c3a9baf6c3bc63de060`；rc=1、唯一红＝同一枚用例：
  `    plan_present_165_test.go:292: Decl 多了字段 AutoApprove：声明面能装下「已批准」就等于加了一个给出允许的出口`。还原=`602d632b…`（MATCH，重跑一次复现同句）。

## 5. 两枚阳性对照（不是变异，但同属这族尺的去牙）

- **V1-AC1 反形**（t1＋append `BuiltinPlanEntries()`）：`名册 11 枚、plan 形 1 枚＝[plan.present]`，rc=0。
- **V1-AC3 授权读路径有牙**（我在 t1 的 `cmd/wisp/v1ac3_grant_positive_control_test.go`，⛔ 入库）：同一会话身份
  `ListGrantsBySession` 先读 **0 行** → 用真 `InsertGrant` 种 **1 行** → 再读 **1 行**，rc=0、
  `V1-AC3 阳性对照成立：同一会话 ListGrantsBySession 0 行 → 种 1 行 → 读出 1 行（⇒ 交付件那句「授权行 0 枚」有牙）`。
- **V1-M4（在 t2，⛔ 在 t1）**：`Gate` 加第三条回调 ⇒ 唯一红「`Gate 有 3 个方法…`」——全文在读数件 00 §3。

## 6. 台面故障登记（⛔ 计入判据）与日志名册

- 形一：`sed -i '116a…'` 的插入把制表符打成了字面 `t`，那一发 `go test` 回 `[build failed]`（rc=1）⇒ **⛔ 当判据**；我逐行重写后重跑，真读数见读数件 00 §3 末。
- 形二：一次 `python` 吃了 MSYS 形路径（`/d/tmp/…`）报 `FileNotFoundError`，**那一发同一命令里跑的 `go test` 用的是未突变树 ⇒ 五枚全绿＝假绿**（就是 M4 那发 t2 的第一次尝试），具名作废；改正路径形后重打才取到真读数（M4 真读数＝读数件 00 §3 末那条「Gate 有 3 个方法」）。此后 M6b／M7／M8 三发的突变都用原生路径形，且每发**先 `diff`／`sed -n` 回读改动真在字节上、再跑**（M6 nested 那发即用 diff 证过改动落点在第 68 行）。
- 上面两发的日志与重打同文件名，被后一发覆盖；盘上现存 `D:\tmp\165v1\logs\*.log`＝**22 枚**，逐名册 grep `walk`／`build failed`／`no tests to run`／`0xc0000135` 四形＝**命中 0 枚文件**（尺＝`grep -l <词> *.log`，22 枚全喂）。

