# 票 179 写腿 179-r3 —— AC#9 尾测补 `decision` 列钉（M2 变异实证）＋ AC#10 二选一归口（选②）

- 腿身份＝实现者（判语归非实现者腿）；⛔ 未翻框、⛔ 未动 `Status:`/票名；⛔ 未 push；⛔ 未跑 AC#8。
- 锚（起手现读）＝`f6fda72d`；派发锚＝`7e8ac3fa` ⇒ **漂移 2 枚**：`8707a7ed`（318-r1 产码，`internal/agent/approval/**`）＋`f6fda72d`（318-r1 交付）——逐报；两者均不触本腿面（本腿面 `internal/agent/loop.go`＋`declared_l0_risk_179_test.go`，起手 `git status --porcelain -- internal/agent/` 只有本腿两枚 ` M`）。产码 commit `0d45d6f4`；读数窗口 2026-10-14 17:1x–17:2x +08。
- 环境：`PATH=…/third_party/sherpa-onnx:$PATH`（MSYS 形）；`GOFLAGS=`；`GOTMPDIR=TMP=TEMP='C:\tmp\wisp179r3gotmp'`（反斜杠形，先 `mkdir`）；判据全 `-run` 小射程 `-count=1`；⛔ 未跑 `./...`、未碰 `internal/risk/`。
- 结论先行：**① AC#9 落地**——尾测补一句读 journal `decision` 列的断言（改前 M2 下照绿→改后 M2 下必红于新断言 `:262`，逐字在 §1）；**② AC#10 选②**——注释改成实话＋归口票 21/162（选型理由与现读证据在 §2、`logs/60`）；**③ 门禁全 rc=0**（§3）；**④ AC#8 未跑**（`cmd/wisp` 本波被 255-r8 独占），下波前置在 §4。

## 1. AC#9：尾测补一句读 `decision` 列的断言（零产码行为改动）

- 改前现量（起手现读复认 `179-v1` 的 M2 少牙形）：尾测 `TestDeclaredL1L2RoutingUnchanged179` 末段（当时 `:244-255`）只断 `ok/unbooked/why` 三值；全文件 `decision` 列断言只长在 `:104`（正向 C1）与 `:220`/`:239`（L1/L2 守卫）三处。
- 落修（commit `0d45d6f4`，diff 逐字）：
  ```go
  	// AC#9: the return values above are only half the promise - the direct L0
  	// arm must also book its allow into the decision column.
  	rows, err := store.ListToolCallsByTask(ctx, j.taskID)
  	if err != nil {
  		t.Fatalf("tool_call lookup: %v", err)
  	}
  	if len(rows) != 1 || rows[0].Decision != DecisionAllow {
  		t.Errorf("decision = %+v, want one row with %q booked by the direct L0 arm", rows, DecisionAllow)
  	}
  ```
  ⛔ 零产码行为改动、未动 `thresholds.go`、未 `t.Skip`、未放宽任何既有断言。
- **改前照绿**（`179-v1` 原判，其日志在册可复算：`.scratch/wisp/probes/179/v1/logs/m2-full-package.txt:6-8`，2026-09-28 12:00:22）——M2 全包跑的红名单只有 `TestDeclaredL0RunsWhenPassThroughIsOff179`（`:104 decision = "", want "allow"`）与 `TestCancelledTaskPersistsTerminalRows`，**尾测不在红名单**＝照绿。本程现读复认该少牙形（上面"改前现量"）。
- **变异落地**（`logs/20`，HEAD=`0d45d6f4` 之上；grep 原文＋build rc=0 逐字在案）：
  - before：`grep -n 'j.decide(ctx, rowID, DecisionAllow)' internal/agent/loop.go` ⇒ `828:`（RiskL0 支，本发靶）与 `837:`（未分级＋开关开着支，**不许碰**）。
  - 靶面锚：`sed -n '822,826p'` ⇒ `… return true, true, ""` / `case RiskL0:`（即 M2 靶＝RiskL0 支那一句）。
  - `sed -i '828d' internal/agent/loop.go`；after：`sed -n '824,836p'` 逐字可见 `case RiskL0:` 下 decide 句已无、`default:` 支的 `j.decide(ctx, rowID, DecisionAllow)` 原样；`grep` ⇒ 只剩 `836:` 一行（＝原 `837` 位移）。
  - `go build ./...` ⇒ **rc=0**、输出 0 字节（`logs/21`，0B 即读数）。
- **M2 ⇒ 尾测必红且逐字点名红行**（同命令，`logs/30`）⇒ **rc=1**：
  - `declared_l0_risk_179_test.go:104: decision = "", want "allow"` ⇒ `--- FAIL: TestDeclaredL0RunsWhenPassThroughIsOff179 (0.12s)`（C1 与 `179-v1` 的 M2 一致）；
  - `declared_l0_risk_179_test.go:262: decision = [{ID:5 TaskID:task-179-l0-direct Seq:1 Tool:fs.read ArgsJSON:{} RiskLevel:L0 Decision: DecidedAt:<nil> … GrantID:<nil>}], want one row with "allow" booked by the direct L0 arm` ⇒ `--- FAIL: TestDeclaredL1L2RoutingUnchanged179 (0.09s)`（**新牙咬住**：`Decision:` 空、点名 `:262`）；
  - 同发四枚 L1/L2 子测照 PASS、`FAIL github.com/CarlosShao/wisp/internal/agent 2.804s`。
- **还原 ⇒ 绿**（`logs/40`）：`git cat-file blob HEAD:internal/agent/loop.go > internal/agent/loop.go` ⇒ `git diff -- internal/agent/loop.go`＝0 字节、grep 2 行再现；同命令 ⇒ **rc=0** 全 PASS（`ok … 3.430s`）；`git status --porcelain -- internal/agent/` 空。
- **改后基线绿**（未变异，`logs/10`）：`-run 'TestDeclared'` ⇒ rc=0、2 PASS（`ok … 2.731s`）。

## 2. AC#10：`decideRisk` 第二返回值注释与产码不一致 ⇒ 选 **②**（注释改成实话）

- 现量：唯一调用点＝`internal/agent/loop.go:651` `ok, _, why := l.decideRisk(ctx, j, p.rowID, info)`（第二返回值被 `_` 丢弃）；旧注释逐字 `"The second result reports \"booked nothing\" (gate-owned routing) so the caller can drop its own pending row instead of leaving one open forever."`——承诺的调用方动作在产码里没有对应物。
- **选②，具名理由**（全文＋现读证据在 `logs/60-ac10-choice.txt`）：
  1. **①（调用方真消费 `unbooked`）在真组合里"没有对象"**：`unbooked=true` 只在"声明 L1/L2 且 AdmitTask 已接"那支产生（`loop.go:820/:825`）；而真组合根正是 `Journal: nil` + `AdmitTask: rt.admitTask`（`cmd/wisp/run.go:1027-1040`）⇒ loop 的 taskJournal 是 `newTaskJournal(nil,…)` 的 no-op 形（`journal.go:59-64`；`startCall` `if t.j == nil { return 0 }`＝`:71`、`decide` 直接 return），**调用方手里根本没有"自己的 pending 行"可 drop**（行归桥写；`internal/tools/loop_approval_test.go:117-119` 逐字同述）。
  2. **①的"该做什么"未定义**：剩下的"门那一跳之后 pending 行归谁收"票面明示属票 21/162 地界；真去实现＝自造判据＋挑一个未定案归宿（drop/close），会动 A13"行留 pending、列归门"的既有分工，有让 AC#3/AC#4 守卫变安静的实险，还撞"未定义即停"。
  3. **②是零行为改动**，恰为格面二选一之一；归口（"注释与产码不一致已归口"）登记在票尾 bullet，台账由编排者登记。
- 落地（commit `0d45d6f4`）新注释逐字：
  ```
  // The second result reports "booked nothing" (gate-owned routing): the row is
  // left pending and the gate writes the decision column (ruling A13). The only
  // caller (executeCalls) currently DISCARDS that value - it drops no row, and
  // this comment must not promise that it does; who settles a gate-owned pending
  // row after the gate hop is deferred (tickets 21/162). Ticket 179 AC#10
  // corrected this text, which used to promise a consumer that never existed.
  ```
- 两枚守卫零触碰自证：AC#3（未分级仍拒）／AC#4（L1/L2 路由一字不变）在改后（`logs/10`）与还原后（`logs/40`）照绿；本程 diff 只含一注释段＋尾测一断言。

## 3. 门禁（取读于 `0d45d6f4`＋r3 台件在盘；交付 commit 后复跑逐字在回执）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| agent 全包 | `GOFLAGS= go test -count=1 ./internal/agent/` | 0 | `ok github.com/CarlosShao/wisp/internal/agent 4.636s`（`logs/50`，54B） |
| 构建 | `go build ./...` | 0 | 空输出（`logs/51`，0B 即读数） |
| 格式 | `gofumpt v0.12.0 -l internal/agent/loop.go internal/agent/declared_l0_risk_179_test.go` | 0 | 空（`logs/52`，0B 即读数） |
| D22 扫描 | `sh scripts/d22scan.sh` | 0 | `clean - no D22 ban violations`；ban #8 internal/=556 Go files（`logs/53`，23758B 留盘） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | `VERDICT GREEN`；tracked=12887 over=60 roster=60 not in roster=0（`logs/54`，1458B） |

## 4. AC#8 下波前置（本波⛔不跑；`cmd/wisp` 现被 255-r8 独占）

- 命令（`cmd/wisp` 空出后直接跑；探针族＝`probes/176/r1` 的续命版，`179-r2` 的台件已修好 r1 的取数写法）：
  ```
  PATH="…/third_party/sherpa-onnx:$PATH" GOFLAGS= GOTMPDIR='C:\tmp\…' \
  go test -count=1 -v -run 'TestRereadHostPointerOnCLISeam179r2' \
    -overlay=.scratch/wisp/probes/179/r2/overlay-e2e.json ./cmd/wisp/
  ```
- overlay 件＝`.scratch/wisp/probes/179/r2/overlay-e2e.json`（把 `cmd/wisp/zz179r2_e2e_test.go` 映射到 `.scratch/wisp/probes/179/r2/zz179r2_e2e_test.go`，跟踪目录零新增）；读数自动落 `probes/179/r2/logs/`（logdir 由 `runtime.Caller` 钉死，不认 CWD——r2 第一跑落错目录的教训）。
- 预期读法：task.output 执行 `risk=L0 decision=allow outcome=success` → 回执带「全文见 …artifacts\tool-output-agent-task-….txt」→ 模型拿它作 `fs.read` 的 path → **该副本逐字节读回**；上一发（`probes/179/r2/logs/e2e-readings.txt`，28496B）阻断者逐字 `level=L2 rules=[R4] reason="R4: 包含来自 task.output 的内容"`（＝票 183 的面）。分岔照票面：仍撞 L2/R4 ⇒ 未结案、指回票 177/162；拿到逐字节读回 ⇒ 本格结案（同时是票 177 AC#3／175 AC#5／176 AC#3、#4 三枚条件格证据）。该跑由做票 183 的那一程顺带取（票面注）。

## 5. 文件清单（产出 7 枚；大日志留盘不追库）

| 件 | 字节 | 说明 |
|---|---|---|
| `00-anchor-and-gates.md`（本件） | — | 交件表 |
| `logs/10-baseline-green.txt` | 1851 | 改后基线绿 |
| `logs/20-m2-mutation-landed.txt` | 1306 | M2 落地：grep 原文/after＋build rc=0 |
| `logs/30-m2-red.txt` | 2275 | M2 红逐字（`:104`＋`:262`） |
| `logs/40-restore-green.txt` | 2055 | 还原（git cat-file）＋绿 rc=0 |
| `logs/54-pathlen.txt` | 1458 | 长度门 GREEN 60=60 |
| `logs/60-ac10-choice.txt` | — | AC#10 选②理由＋注释前后逐字 |
| 留盘不追库 | | `logs/21-m2-build.txt`(0B)、`logs/50-agent-pkg.txt`(54B)、`logs/51-build.txt`(0B)、`logs/52-gofumpt.txt`(0B)、`logs/53-d22scan.txt`(23758B) |
