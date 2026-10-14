# 票 315 验收腿 `315-v1` —— AC#1 的**非实现者**判语（被裁＝`bfc74ae9` 产码＋`b7a15297` 交付；实现者＝`315-r1`）

- 腿身份＝**验收（非实现者）**。⛔ 未翻框、⛔ 未动 `Status:`/票名；判语交回。
- 锚：起手现读 HEAD＝`b7a15297`；读数窗口 **2026-10-14 11:32–11:45 +08**（窗口内 HEAD 漂到 `e3788f08`＝288-r1 停手上报＋编排者三笔文档；
  逐枚 `--name-only` 核过**与本腿判读面零相交**：`cmd/wisp/run.go` sha256 `71093661…` 全窗前后一致、`zz315a1probe_test.go` porcelain 空）。
  上游 AC#0＝`315-a1`（`e0d6d383`＋`810deeb4`）。
- 环境：`PATH=…/third_party/sherpa-onnx:$PATH`；`GOTMPDIR=C:\tmp\wisp315v1gotmp`；`GOFLAGS=`；全用 `-run` 小射程 `-count=1`；⛔ 未跑整包 `./cmd/wisp/`（314 死区）。
- 读数名册（含字节/sha256 前缀/窗口/rc）＝`.scratch/wisp/probes/315/v1/00-anchor.md`；突变/反证件全在仓外 `C:\tmp\wisp315v1overlay\`（工作树零足迹，还原证在锚件）。
- **结论先行：六格全「成立」**；无停手项。**收口建议：AC#0＋AC#1 双成立 ⇒ 票 315 可收**（翻框归编排者）。

## 判语表（每格＝成立/附条件/不成立；凭据均为本腿现读/现跑，不单引实现者的件）

| # | 格 | 判语 | 凭据（逐条） |
|---|---|---|---|
| 1 | 选案正当性：ⓐ`admitTask` 补开，排除 ⓑ子环路自开／ⓒ桥懒开扩到每次 dispatch | **成立** | ① 四处引证逐字核过（我现读）：`internal/risk/provenance.go:453-454`「composition root calls this at task start」（跨 453/454 两行，453 起）；`:706` 日志句「fail-closed R4 (ensure OpenScope at task start)」（恰在 706）；`docs/evidence/s1/151-task-scope-never-closed-r1.md:192` §4 行①「今天这一枚组合根里**唯一**同时满足两事的 per-task 边界就是 D47 那枚 `AdmitTask` hook…零新生命周期、零新接口」；`internal/tools/subagent_222_test.go:331-333` 逐字 `AdmitTask: func(taskID string) func() {` / `h.bridge.OpenTask(taskID)` / `return func() { h.bridge.CloseTask(taskID) }`。② **用既有机制、非改 C25 语义**：`git diff bfc74ae9^ bfc74ae9 -- cmd/wisp/run.go` 非注释产码增行**恰 1 条**＝`+ rt.bridge.OpenTask(taskID)`（其余 33 增/5 删全注释）；`--name-only` 只两枚（run.go＋zz 探针件），未动 `internal/risk/**`/阈值/golden/allowlist/PLAN/specs。`OpenTask` 是既有 API（`internal/tools/bridge.go:847`），`mark()` 懒开兜底仍在（`bridge.go:649`）⇒ 未过 admission 的直派调用仍走 fail-closed，分支可达性只在"已登记任务"上按设计文本收窄。③ ⓑ 排除复核：子环路工具面＝父桥（`subagent_197.go:296`）、子任务 admission 经 `opt.AdmitTask` 包 `baseAdmit`（`:311-320`）、dispatch 带自己 task id（`loop.go:675-676`）⇒ 子环路本就走这枚组合根边界，自开无覆盖面增量还要新路由；且盖不住根任务。④ ⓒ 排除复核：先于评估器登记 ⇒ unbound 分支对一切走桥调用不可达＝C19/C25 判定语义变化（票面禁区）。⑤ 附注（不降判）：151 §4 把"开在任务开始"**刻意留给编排者判**（该件 `:197/:201` 逐字）；本票框（AC#1 明列 ⓐ/ⓑ/ⓒ 三案）即那枚授权 ⇒ 非"单方面改契约"。 |
| 2 | 读数独立复现 | **成立** | ① 默认三形（无门控，我现跑 11:32:58→11:33:06）：**rc=0**，A `PASS 0.15s`／B `PASS 0.20s`／C `PASS 0.14s`，`ok … 3.052s`；stdout `[确认 L2` 命中 **0**；乙走桥行逐字 `call kind=success … tool=fs.read risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true grant_id=0`；五条 close 全 `was_open=true`（乙 dropped=1、根 dropped=2）。② **自造反证**（我自己的形：仓外 sed 摘掉那一行，`run.v1muted.go` sha `d6a9add6…`，overlay 跑，工作树零足迹）⇒ B **rc=1** `--- FAIL: TestZZ315A1SiblingChildFirstBridgeCall (31.22s)`，卡逐字 `[确认 L2 fs.read] R4: 包含来自 unbound-scope scope is not open (no OpenScope for this id, or its own handle closed it) 的内容`，`ANSWER-EXPIRED … decision=timeout->reject after=31s`，乙 close `was_open=false dropped=0`；C **rc=1** `(31.17s)` 同族（根的首发 `L2/timeout/error`）。③ 交叉：复跑他们入库的 `overlay-mut-noopen.json` ⇒ B **rc=1** `(31.21s)` 同红句；其 `run.mut-noopen.go` 我 `diff` 对现读 run.go＝**恰一行替换**（sha `3b632e95…`）。④ 修前红留档：其 log 11（未修树）B rc=1 31.22s、乙 `was_open=false`、同逐字红句；log 30/31（突变 B/C）同族 —— 与我自造反证**双向坐实**（摘一行必红 31s/unbound，带一行必绿 0.14–0.20s）⇒ 钉子有牙，非断言方向造的绿。 |
| 3 | Case C（`AGENTS.md` 种子）正当性 | **成立** | ① 前置断言写法（现读 `zz315a1probe_test.go:685-688`）：数据目录里那枚 `AGENTS.md` 未进加载清单（stderr 无该路径或 `tier=global`）即 `t.Fatalf`「本形量不到，不算绿」＝**断言而非假定**；绿/红两跑 stderr 均含逐字 `…\AGENTS.md tier=global bytes=89 状态=完整注入`（种子确落）。② 读数：修前（overlay 无本修）`rc=1` **31.17s**（他们留档 31.15s）卡与 `unbound-scope` 同红句；修后 `rc=0` **0.14s**（他们 0.12s／23 号件 0.16s），根首发行 `risk=L0 decision=allow outcome=success`、根 close `was_open=true dropped=1`。③ 机制现读：seed 由 projctx 装载器在请求构造时盖 `run:project-instructions`（作用域 `cmd/wisp/run.go:1081`；`internal/projctx/projctx.go:349` `Mark(ScopeID, risk.SrcFSRead, …)`）⇒ 这是 a1 未试、写腿"自然涵盖"的第 2 形，读数与断言一致，非"为绿而绿"。 |
| 4 | 去门控决定（vs 入 ledger vs 保持门控） | **成立** | ① 我现跑 strict runner：`sh tools/d22scan/runtests.sh ./cmd/wisp/ -run '^TestZZ315A1'` ⇒ **rc=0**，逐字 `runtests.sh: OK - … PASS=3 FAIL=0 SKIP=0, === RUN=3`（11:39:26→11:39:35）。② 前提独立复现（我复造门控形：overlay 给 Case B 加回 `t.Skip`）⇒ 同 runner **rc=1**，逐字 `runtests.sh: 1 test(s) SKIPPED and SKIP is not a pass (ticket 71 AC#3) - … PASS=0 FAIL=0 SKIP=1` ⇒ "保持门控"在 cli scope（`scripts/wisp-cli-tests.sh:218`→`portable-tests.sh --scope=cli`→此 runner）就是常红。③ `grep -n 315 scripts/portable-tests.sh` ⇒ **零命中**；该 ledger 是"平台跑不了的 SKIP 具名入账"面（174-r5 刚动过）＝别的腿写面，且修后这枚钉子能跑（0.2s）——把它注册成永久 SKIP＝消声。三选落法：保持门控＝rc=1 禁；入 ledger＝在别人写面消声；**去门控＝钉子默认跑且 strict 绿** ⇒ 落法正确。 |
| 5 | 既有 `-race` 夹具红归因 | **成立（非本修引入）** | ① 我现跑当前树（带本修）：`-race -run '^Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands$' -v` ⇒ **rc=1**，`WARNING: DATA RACE`×7，`--- FAIL … (13.99s)`；栈顶逐字 `bytes.(*Buffer).Write()` ⇄ `cmd/wisp.(*agentRuntime).auditf()` `run.go:933`，一侧来自测试自起的 goroutine（`subagent_selfapproval_197_test.go:115 startChildWrite197.func1`→`Bridge.route`→`Gate.PendingApproval`→`consoleApprovalUI.Update`→`bookWaitingState`）。② 同一命令 **overlay 摘掉本修那一行**（＝修前生产语义）⇒ **rc=1**，同样 7 枚 warning、`--- FAIL … (13.95s)` ⇒ 本修既没引入也没修复它。③ 无 `-race` ⇒ rc=0（两枚 PASS 13.66/0.00s）。④ 台账有账（我现读）：`docs/reports/pending-and-issues.md` `A865` 段（`:17058`）逐字记过同栈顶 `bytes.(*Buffer).Write()`⇄`auditf()` `run.go:933`；同族测试侧缓冲竞争另有归口（`:10484`，归票 223 一枚新判据）⇒ 归因"既有夹具红"成立，本票不动（射程外）。 |
| 6 | 门禁抽核 | **成立** | 我现跑（窗口尾）：`GOFLAGS= go build ./...` **rc=0**（空输出）／`sh scripts/d22scan.sh` **rc=0** `clean`（cmd/ examined 125）／`sh scripts/check-path-length-budget.sh` **rc=0** `VERDICT GREEN`，逐字 `over the hat=59 … roster entries=59`（**59=59**）／`go test ./internal/tools/ -race -count=1 -run 'Test222'` **rc=0**（3 PASS，ok 4.464s）。另：`-race` 三形一发 rc=0（ok 5.395s）。 |

## 门禁 rc 表（本腿现跑；日志见锚件名册）

| 门 | 命令（要点） | rc | 读数 |
|---|---|---|---|
| 默认三形 | `go test ./cmd/wisp/ -run 'TestZZ315A1' -count=1 -v` | 0 | 3 PASS（0.15/0.20/0.14s）；无 `[确认 L2` |
| 自造反证 B/C | 同上 `-overlay=C:/tmp/wisp315v1overlay/v1-overlay.json` | 1 / 1 | 31.22s / 31.17s，`unbound-scope`＋`after=31s` |
| 交叉（他们件） | `-overlay=…/probes/315/r1/mutation/overlay-mut-noopen.json` | 1 | B 31.21s 同红句 |
| strict runner（去门控） | `sh tools/d22scan/runtests.sh ./cmd/wisp/ -run '^TestZZ315A1'` | 0 | PASS=3 FAIL=0 SKIP=0 |
| strict runner（门控复造） | 同上＋`v1-overlay-gated.json` | 1 | `SKIP is not a pass`，SKIP=1 |
| `-race` 三形 | `go test ./cmd/wisp/ -race -count=1 -run 'TestZZ315A1'` | 0 | ok 5.395s |
| 197 族 `-race`（当前树／无本修） | `… -race -run '^Test197SubagentSelfApproval…$'` | 1 / 1 | 各 7 warning；FAIL 13.99s / 13.95s（既有夹具红） |
| 197 族无 `-race` | 同 `-run`，无 `-race` | 0 | 2 PASS |
| 构建 | `GOFLAGS= go build ./...` | 0 | 空输出 |
| D22 扫描 | `sh scripts/d22scan.sh` | 0 | `clean`；cmd/ examined 125 |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | `VERDICT GREEN`；over=59／roster=59 |
| 票 237 家族 | `go test ./internal/tools/ -race -count=1 -run 'Test222'` | 0 | 3 PASS，ok 4.464s |

## 给编排者的收口建议

- **AC#0（`315-a1` 活体复现）＋ AC#1（本表六格全成立）双成立 ⇒ 票 315 可收**（翻框归编排者）。附条件：无。
- 具名遗留（不阻塞收口）：**并发兄弟形**仍未单独量（写腿已登记候选；机制上 admission 按任务而非按调度，修后 run 里两枚孩子 close 行全部 `was_open=true`；若要读数，补一发双孩子在飞的最小形）。
- 具名没做：⛔ 未跑整包 `./cmd/wisp/`；⛔ 未修既有 `-race` 夹具红；⛔ 未动 C19/C25 语义/阈值/golden/allowlist/PLAN/specs；⛔ 未动 `Status:`/框/票名；⛔ 未 push。
