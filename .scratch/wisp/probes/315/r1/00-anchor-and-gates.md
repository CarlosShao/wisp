# 票 315 写腿 315-r1 —— AC#1 落地：任务开始即开 C25 作用域（选案 ⓐ）

- 腿身份＝实现者（判语归后续非实现者腿）；⛔ 未翻框、⛔ 未动 `Status:`/票名；⛔ 未 push。
- 锚（起手现读）＝`5ed56c2f`（派单写"近 `0137ec08`"；其后 174-r5 两笔已入树，与本腿文件面零相交）。
  读数窗口 2026-10-14 **11:07–11:23 +08**（另 11:24 后补门禁复跑）。
- 环境：`PATH=…/third_party/sherpa-onnx:$PATH`（MSYS 形）；`GOTMPDIR=C:\tmp\wisp315r1gotmp`；`GOFLAGS=`；
  判据全用 `-run` 小射程 `-count=1`；⛔ 未跑整包 `./cmd/wisp/`（314 死区）。
- 结论先行：**① 选案 ⓐ 落地**（`cmd/wisp/run.go` 的 `admitTask` 在任务开始处 `OpenTask(taskID)`，与同钩子 revoke 里那枚
  `CloseTask(taskID)` 同 id、同所有者成对）；**② 修前红→修后绿逐字复到**；**③ 钉子去门控**（a1 的 `WISP_315A1`
  跳过在 cli scope 的严格 runner 下是 `SKIP=1` 致命红，见 §4）；**④ 自然涵盖第二形并活体复到**（project-instructions
  种子 mark，见 §5）；**⑤ 门禁全 rc=0**（§6）。

## 1. 设计读出与选案（先读再动；引文逐字）

**原文指向"组合根在任务开始处开"**：

- `internal/risk/provenance.go:453-454`（OpenScope 文档）："OpenScope registers a task/session scope (**composition
  root calls this at task start**) and HANDS BACK the handle that owns it."
- `internal/risk/provenance.go:58-62`（C25 的 `DEFERRED(C25-loop-wiring)` 项 (1)）："…must (1) OpenScope at task
  start and close what it handed back on the boundary that owns the task — LANDED as a SHAPE by ticket 160…"；
  ⚠ "Still OPEN inside (1)" 的那半是 DisposalScope 形状，**不是**"该不该在任务开始开"。
- `internal/risk/provenance.go:706`（unbound fail-closed 的日志原文）："…fail-closed R4 (**ensure OpenScope at task
  start**)"——判定语义自己就把"开在任务开始"写成接线要求。
- `internal/tools/bridge.go:834-838`（OpenTask 文档）："ticket 19's DEFERRED(C25-loop-wiring) item (1), landed as a
  shape by ticket 160… Idempotent; Execute opens lazily so a caller that forgot cannot get an untainted read."
- `internal/tools/bridge.go:859-862`（CloseTask 文档）："The composition root defers this on the task's own boundary:
  cmd/wisp wires it into agent.Options.AdmitTask's revoke…" ⇒ **开与关在同一份契约文本里都指这枚组合根边界**。
- 票 151 验收件 §4 表行 ①（`docs/evidence/s1/151-task-scope-never-closed-r1.md`）："今天这一枚组合根里**唯一**同时
  满足两事的 per-task 边界就是 D47 那枚 `AdmitTask` hook…**零新生命周期、零新接口**。"
- 票 237 的 harness 侧同形（`internal/tools/subagent_222_test.go:331-333`）：`AdmitTask: func(taskID string) func() {
  h.bridge.OpenTask(taskID); return func() { h.bridge.CloseTask(taskID) } }` ⇒ ⓐ 正是这枚**harness 修法的生产侧对称形**。
- 子任务确实经过这枚钩子：`internal/tools/subagent_197.go:311-320` 的 `opt.AdmitTask` 包着 `baseAdmit`（= 组合根的
  `rt.admitTask`）⇒ 根任务与每枚子代理都在这条边界上，ⓑ 的"子环路自开"在覆盖面上并无增量。

**三案裁决（具名理由）**：

| 案 | 裁决 | 理由（现读） |
|---|---|---|
| ⓐ `admitTask` 补开 | **采用** | 与既有 `CloseTask` 同边界、同 id、同所有者（桥持 `OpenTask` 存下的 handle）；零新接口、零新生命周期；直落上面 5 条引文。 |
| ⓑ 子环路自开 | 否 | 子环路不持桥（`internal/tools/subagent_197.go:296` 的工具面才是父桥）；要开就得从 `internal/agent` 造一条新路由到 C25 引擎（Q-56 一族的新接口），且**盖不住根任务自己**；而子任务早已有 `AdmitTask` 这条现成边界（`:311-320`）。 |
| ⓒ 桥懒开扩到每次 dispatch | 否 | 会在评估器**之前**登记 scope，使 unbound fail-closed 分支对一切走桥调用不可达 = **改 C19/C25 判定语义**（票面禁区）；且会给"宿主自带 id、无人关"那一形每调用开一枚无人关闭的 scope（票 154 残余形反向放大）。 |

⇒ 未触"必须改 C25 语义"的停手条件；未动 `thresholds.go`/规则集/golden/`allowlist.txt`/`PLAN.md`/`specs/**`。

## 2. 落修（产码 1 行 + 注释）

`cmd/wisp/run.go` 的 `admitTask`，在任务开始处：

```go
	rt.noteTask(taskID)
	// The open side of the pair, at the task's own start: from here until the
	// revocation below this task's scope is registered, so its own first bridge
	// call is inspected against its own (still untainted) scope instead of
	// failing closed as an unbound one (ticket 315).
	rt.bridge.OpenTask(taskID)
	revoke := rt.gate.AdmitTextTask(taskID)
	return func() {
		if revoke != nil { revoke() }
		rt.bridge.CloseTask(taskID)
	}
```

同一函数的文档注释补了 WHY-THE-OPEN 段（引原文、记两案排除、记"懒开仍是边界外调用者的兜底"）。总 diff：
`run.go +34/-5`（其中产码 1 行）、`zz315a1probe_test.go +97/-16`（全在钉与注释面）。

## 3. 修前红 / 修后绿 / 突变逐字

- **修前红**（锚 `5ed56c2f`，未修树；命令＝`WISP_315A1=1 GOFLAGS= go test ./cmd/wisp/ -run
  'TestZZ315A1SiblingChildFirstBridgeCall' -count=1 -v`，11:07:29→11:08:09）⇒ **rc=1**，
  `--- FAIL: TestZZ315A1SiblingChildFirstBridgeCall (31.22s)`；卡逐字
  `[确认 L2 fs.read] R4: 包含来自 unbound-scope scope is not open (no OpenScope for this id, or its own handle closed it) 的内容`；
  `[audit] approval: ANSWER-EXPIRED … decision=timeout->reject after=31s`；
  乙 close 行 `was_open=false dropped=0`（对照甲 `was_open=true dropped=1`）。
- **修后绿**（同树＋本修；`go test ./cmd/wisp/ -run 'TestZZ315A1' -count=1 -v`，11:23:29→11:23:37）⇒ **rc=0**：
  A `PASS (0.15s)`＋B `PASS (0.18s)`＋C `PASS (0.12s)`；乙走桥行逐字
  `call kind=success … tool=fs.read risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true`；
  三条 close 行全 `was_open=true`（甲 dropped=1 / 乙 dropped=1 / 根 dropped=2）；`card snapshot (seen=false)`。
- **突变（钉子有牙）**：`-overlay` 把 `run.go` 换成一枚只摘掉 `rt.bridge.OpenTask(taskID)` 那一行的副本
  （工作树零突变；件＝`mutation/run.mut-noopen.go`＋`mutation/overlay-mut-noopen.json`）
  ⇒ B 形 **rc=1**（11:13:53→11:14:34，`--- FAIL … (31.23s)`，红句与修前同族）、C 形 **rc=1**
  （11:22:33→11:23:16，`--- FAIL … (31.15s)`）⇒ 摘这一行，两枚钉必红；不是断言方向造的绿。
- 稳定性：`-race -count=2 -run 'TestZZ315A1'` ⇒ rc=0（ok 5.690s）；`-race -count=3`（三形齐全前）rc=0。

## 4. 钉子与门控：**去门控**（具名理由）

- 钉子＝`cmd/wisp/zz315a1probe_test.go` 的 **Case B**（`TestZZ315A1SiblingChildFirstBridgeCall`，逐字断言"白名单内
  首调用必须 L0/allow/success"）＋新增 **Case C**（§5）；Case A 对照仍常驻。
- a1 的 `WISP_315A1=1` 门控已删除。三条理由，逐条有读数：
  1. **带门控落地本身就是一枚无账 SKIP**：cli scope 经 `scripts/wisp-cli-tests.sh` →`portable-tests.sh --scope=cli`
     →`tools/d22scan/runtests.sh`，其"任何顶层 `--- SKIP` 即致命（票 71 AC#3）"我在本锚点实跑复到：
     `sh tools/d22scan/runtests.sh ./cmd/wisp/ -run '^TestZZ315A1SiblingChildFirstBridgeCall$'` ⇒ **rc=1**，
     `1 test(s) SKIPPED and SKIP is not a pass … PASS=0 FAIL=0 SKIP=1`（11:12:07→11:12:15；
     ledger 无此名，`scripts/portable-tests.sh` 里 `grep 315` 命中 0）。
  2. **去门控不需要动 ledger**：ledger 在 `scripts/portable-tests.sh`——是别的腿（174-r5）的写面；去门控后无 SKIP，
     那枚文件一字不碰。
  3. **默认跑传统与成本**：cmd/wisp 的默认跑里本就有无门控的组合根端到端用例（`run_test.go`、`task_scope_close_151_test.go`）；
     修后 B 形"绿时快"（0.18s，无卡即无 31s），A/B/C 三形断言耗时合计 ≈0.45s。
- 判据（dispatch 第 3 项）＝修前红（§3，含未修树与突变两把尺）→修后绿；且 Case A 对照仍绿（0.15s）。
  ⛔ 未把"会红的测试"默认跑进包——默认态三形全绿。

## 5. 自然涵盖：a1 未试的 `project-instructions 作种子 mark` 形 → 已活体复到（Case C）

- Case C＝根任务的第一发走桥调用之前，数据目录放一份 `AGENTS.md`（projctx 的 global 层），装载器在**第一次模型请求
  构造时就**把内容盖戳进 `run:project-instructions`（`internal/projctx/projctx.go` 的 `stamp`，源名 `risk.SrcFSRead`）。
- 修前（同一 overlay 突变）**rc=1 / 31.15s / L2 timeout**；修后 **PASS 0.16s**，根的第一发 `fs.read` 逐字
  `risk=L0 Decision:allow Outcome:success`，前置条件（`AGENTS.md tier=global bytes=89 状态=完整注入`）被断言而非假定。
- 另一未试形（**并发兄弟**）本腿未单独量：结构上与被修条件同一（"该 id 未登记"现在对所有经 admission 的任务为假，
  admission 是按任务而非按调度发生的；修后 run 里两枚孩子各自的 close 行都 `was_open=true`）——**登记为候选，
  若要读数请补一发双孩子在飞的最小形**，本腿不扩修。

## 6. 门禁 rc 表（改后终树；各日志见 §8）

| 门 | 命令（要点） | rc | 读数/注 |
|---|---|---|---|
| 钉子三形 | `go test ./cmd/wisp/ -run 'TestZZ315A1' -count=1 -v` | 0 | 3 PASS（0.15/0.18/0.12s） |
| 钉子 -race | `… -race -count=2 -run 'TestZZ315A1'` | 0 | ok 5.690s |
| 票 237 家族 | `go test ./internal/tools/ -race -count=1 -run 'Test222'` | 0 | 3 PASS，ok 3.607s |
| 151 关侧钉 | `-run 'TestCompositionRootClosesTheLoopTasksTaintScope\|TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit'` | 0 | 两枚 PASS（3.38/3.50s） |
| 宿主直派 | `-run 'TestHostDispatchThroughTheAssembledBridge'` | 0 | PASS 3.67s |
| 197 自批族 | `-run 'Test197SubagentSelfApproval…\|Test197NoAllowDoor…'` | 0 | 两枚 PASS（13.85/0.00s） |
| 158 开侧钉 | `go test ./internal/tools/ -count=1 -run 'TestSensitiveReadAcrossTheBridgeOpensItsTaskScope'` | 0 | PASS 0.01s |
| 构建 | `GOFLAGS= go build ./...` | 0 | 空输出 |
| vet | `go vet ./cmd/wisp/` | 0 | 空输出 |
| 格式 | `gofumpt v0.12.0 -l` 两文件／`gofmt -l` 两文件 | 0 | 均空 |
| D22 扫描 | `sh scripts/d22scan.sh` | 0 | `clean`；cmd/ examined 125（含 _test.go） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | `VERDICT GREEN`，12599 tracked／over=59／roster=59 |
| 严格 runner（SKIP 尺） | `sh tools/d22scan/runtests.sh ./cmd/wisp/ -run '^…SiblingChild…$'`（**带门控旧形**） | **1** | `SKIP=1` 致命 ⇒ 去门控的读数（§4） |

⚠ **改前既有、与本修无关**：`-run '…roster…' -race` 下 `Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`
出 DATA RACE（栈顶 `bytes.(*Buffer).Write()` ⇄ `auditf()` `run.go:933`，即测试夹具 stderr 缓冲被并发写；同族已在
`docs/reports/pending-and-issues.md` 有账）。**决定性对照**＝同一 `-race` 在 overlay 摘掉本修那一行的树上同样红
（6 枚 WARNING、同 FAIL，11:16:14→11:16:41）⇒ **不是本修引入**；无 `-race` 时该族 5 枚全 PASS（rc=0，11:16:52→11:17:24）。
本腿不顺手修（射程外、夹具面）。

## 7. 具名没做

- ⛔ 未动 C19/C25 判定语义、`internal/risk/**`、`thresholds.go`、规则集、golden、`allowlist.txt`、三枚冻结件。
- ⛔ 未动 `scripts/**`（去门控后也无需动 ledger）；⛔ `frontend/**`/`design/**`；⛔ 未 push、未 `add -A`。
- ⛔ 未修 §6 那枚既有 `-race` 夹具红；⛔ 未单独量"并发兄弟"形（登记为候选，§5）。
- ⛔ 未改票面 `Status:`/框；票尾只追加 1 条 bullet（交付笔）。

## 8. 读数名册（`.scratch/wisp/probes/315/r1/`；大日志留盘，点名＋字节数）

| 件 | 字节 | 是什么 |
|---|---|---|
| `logs/10-prefix-caseA.txt` | 5,932 | 修前 Case A 对照 rc=0（11:07:13→11:07:21） |
| `logs/11-prefix-caseB.txt` | 17,615 | 修前 Case B 红 rc=1 31.22s（**入库**） |
| `logs/12-prefix-gated-skip-under-strict-runner.txt` | 738 | 旧门控在严格 runner 下 SKIP=1 rc=1（**入库**） |
| `logs/20-postfix-caseB-gated.txt` | 10,068 | 修后 B（仍带门控）rc=0 0.21s |
| `logs/21-postfix-ungated-both.txt` | 15,813 | 去门控后 A+B rc=0 |
| `logs/22-postfix-ungated-race-count3.txt` | 48 | `-race -count=3` rc=0 |
| `logs/23-postfix-caseC-instructions-seed.txt` | 6,775 | Case C 修后 rc=0 0.16s |
| `logs/24-final-three-cases.txt` | 22,399 | 终树三形 rc=0（**入库**） |
| `logs/25-final-race-count2.txt` | 48 | `-race -count=2` rc=0 |
| `logs/30-mutation-noopen-red.txt` | 17,615 | 突变（摘那一行）B 红 rc=1 31.23s（**入库**） |
| `logs/31-mutation-caseC-red.txt` | 12,383 | 突变 C 红 rc=1 31.15s（**入库**） |
| `logs/40-test222-race.txt` | 569 | Test222 -race rc=0 |
| `logs/41-cmdwisp-task-approval-roster-race.txt` | 58,806 | 名册 `-race` rc=1（既有夹具红，§6） |
| `logs/42-197selfapproval-race-overlay-noopen.txt` | 33,052 | 同红在 overlay 无本修的树上复到 |
| `logs/43-cmdwisp-task-approval-roster-norace.txt` | 12,977 | 名册无 `-race` rc=0（5 PASS） |
| `logs/44-tools-scope-open-family.txt` | 321 | 158 开侧钉 rc=0 |
| `logs/50-build.txt` / `51-vet-cmdwisp.txt` | 0 / 0 | build/vet rc=0 |
| `logs/52-gofumpt.txt` | 78 | gofumpt/gofmt 均空 |
| `logs/53-d22scan.txt` | 22,972 | d22scan rc=0 clean |
| `logs/54-pathlen.txt` | 1,457 | 长度门 rc=0 GREEN |
| `mutation/overlay-mut-noopen.json` | 179 | overlay 映射（**入库**） |
| `mutation/run.mut-noopen.go` | 69,267 | 突变树全文（留盘；重建＝对现读 `run.go` 跑一次 `sed 's/^\trt\.bridge\.OpenTask(taskID)$/\t\/\/ MUTATION…/'`） |
