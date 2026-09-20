# 票 66 对抗验收报告 — SLO 判据仪器返工（A14 解析末项 + A15 树外 CPU 口径）

**验收人：** orchestrator（**不是**实现人 `agent-ticket66` / `agent-ticket66b`）
**时间：** 2026-09-20 23:55 本地（15:55Z）
**被验收物：** `.scratch/wisp/issues/66-slo-instrument-parse-cpu-observer.md`（AC#1–AC#7，7/7 已勾）
**判据来源：** 票面 AC 原文 + `docs/reports/pending-and-issues.md` **A14/A15** + D32 阈值纪律

> ## 为什么这份报告存在（写在最前面，免得后人以为它是走形式）
> 今天我在票 62 上盖过一个**假 `-done`**（八个 AC 框一个都没勾、且 AC#8 要的裁决表根本不存在），
> 已登记为 **A30①** 并撤销（提交 `fc33532`）。票 66 是同一个规则的**正面执行**：它 7/7 框已勾、
> 归档物齐全，但**缺一张与 AC 编号 1:1 的裁决表**（README 完成规则 6）。本表补的就是这张。
>
> **每一行都标注了"谁复现的"**，三种标签含义固定：
> - **〔编排者独立复现〕** = 我自己敲的命令、我自己读到的真实 exit code，不引用代理的任何自述。
> - **〔代理日志＋归档，我抽验〕** = 数字来自代理的 run，但我核对了归档文件确实存在且字段对得上。
> - **〔仅代理自述〕** = 我没能独立复现，且**不会替它背书**。本表有一行落在这一档，见 AC#1。
>
> ⚠ 复现窗口约束：本文写就时**票 64 的 winlive 复测正在占用桌面**，而 SLO 的 CPU 采样对
> 观测者自身负载极敏感（A15 的根因就是"观测者单次读 ≈1.3ms，250ms 间隔下本身 0.52% > 0.5% 门"，
> 这条是我今天自己量出来的）。因此我**刻意没有**在 23:5x 重跑 `-Subset smoke` 或 `wisp slo`——
> 那既会污染票 64 的逐跑，也会让 AC#4 的复现变成"在受扰机器上测仪器"。AC#4 用的是我 **23:36 的独立复跑**。

## 裁决表（与 AC 编号 1:1）

| AC | 判据 | 裁决 | 证据与复现方式 |
|---|---|---|---|
| **#1** `parseSystemProcesses` 顺序修正 + 构造"目标 pid 恰为快照链末项"的回归用例 | **PASS** | 代码侧〔编排者独立核对〕：`internal/proc/treemetrics_windows.go:239-262` 现在**只经由唯一一份解码器** `WalkSystemProcesses` 取条目，注释 `:235-238` 明写"last entry must be delivered before the `NextEntryOffset == 0` test"——即今天丢末项的那个顺序；另**新加了一道自报门** `len(out)==0 → error`（`:257-259`），这是 A30⑤/票 71"门必须自报工作量"纪律在实现里的体现。回归用例存在：〔编排者独立 grep〕`internal/proc/systemprocs_windows_test.go:79 TestParseSystemProcessesKeepsLastSnapshotEntry`，同文件 `:202 TestSystemProcessSnapshotSeesSelf`。变异检验（把顺序退回旧的→测试必须变红）：**〔仅代理自述＋归档〕** `docs/evidence/s1/66-mutation-parse-order.md`（两轮，第二轮由 `agent-ticket66b` 按我的要求换人重做：`grep -n MUTATION-66B` + `git diff --stat` 先证变异真落盘，撤变异后 `grep -c MUTATION`=0 且 `git diff` 空）。**我没在 23:5x 重做变异**，因为那要改 `.go` 而票 70 正在全仓格式化。 |
| **#2** `wisp slo -state Sleeping -seconds 10 -interval-ms 250` 在**零 keeper**（无残留进程）下可用 | **PASS** | 〔代理日志＋归档，我抽验〕`docs/evidence/s1/66/66b-nokeeper-{1..5}.json` **五份齐全**（另存上一轮 `66-nokeeper-{1..5}.json`），我 23:36 前亲自跑过同类命令并核对 `exit=0`。测量前后 `tasklist \| grep -iE "wisp\|balldebug"` 均为空的纪律由代理记录。 |
| **#3** CPU 获得**树外**口径（父读子 pid，复用 `cmd/balldebug/diff_windows.go` 那一份走链实现），树内只作记录不作门 | **PASS** | 〔代理日志＋归档，我抽验〕同一 30s 窗双条件：树外 **0.0173%**（当门）vs 树内 **0.7629%** 带 `observer_cost:true`（不当门），原始 JSON `docs/evidence/s1/66/66-ac3-30s.json`。⚠ **D32 的 `Sleeping` CPU ≤0.5% 阈值一字未动**：两侧都仍按 `<=0.5%` 判并写 pass/fail，改的只是**哪一条有资格当门**——这个区分由 `internal/observe/observer_cost_test.go` 4 条用例钉住（"只许 CPU 行动、limit 字符串两侧必须相等、非 CPU 行的 gate 不许随口径变化"）。〔编排者独立复现〕我 23:36 亲自复跑时读到"树内记录 0.1612–1.0322%"这一族数字，与"树内口径下这条门无定义"的判定一致。 |
| **#4** `scripts/slo-check.ps1 -Subset smoke` 本地**连续 5 次 exit=0**，且泄漏自检仍能翻红 | **PASS** | 〔代理日志＋归档，我抽验〕`docs/evidence/s1/66/66-smoke-run-{1..5}.log` 五份，崩溃原文 `66-smoke-pre-fix-crash.log`。〔编排者独立复现〕**23:36 我亲自跑**：`REAL_EXIT=0`，`state Sleeping exit=0 pass=True`、`state Warm exit=0 pass=True`、`settle exit=0 pass=True`、`all_pass=True`，且**同一次跑里 `leak exit=1 flipped_to_fail=True`**——最后这条是关键：它证明"全绿"不是把门焊死后看到的绿。 |
| **#5** 追溯票 08（在票 08 的 log 追加一行事实，**不改它的 AC 勾选与 Status**） | **PASS** | 〔编排者独立核对〕`.scratch/wisp/issues/08-observability-slo-done.md:73` 就是那一行，且它**只补事实**：逐字段读票 08 自己归档的 `build/slo/slo-report.json`（`generated_at 2026-09-19T23:35:49Z`）→ `Sleeping exit_code=2/report=null`、`Warm exit_code=2/report=null`、`settle exit_code=1/pass=false`、`all_pass=false`。⇒ 那次九项门**一门都没产出样本窗**，红因是仪器不是预算。**权力边界写得很清楚**（"不改本票任何 AC 勾选与 Status 字段——那是编排者的权"），我在票面上也的确一个字没动（今天核对过：票 08 勾框数未变）。 |
| **#6** 文档同步：`docs/SLO.md` 附录 A/B 的数字若因本票改动则补 C，**阈值一律不降** | **PASS** | 〔编排者独立核对〕`docs/SLO.md:405` 起为**附录 C**（C.1 五连零 keeper／C.2 变异复测／C.3 同窗双条件／C.4 smoke 五连＋第三起缺陷／C.5 full 六态与 skeleton 边界／C.6 逐字命令／C.7 与 A、B 的关系）。**A.1/A.2/A.4/B 的原始样本表一字未改**——我在 23:36–23:5x 只往 B 追加过一条**对我自己预测的更正**（B.3 的"约一半概率随机红"被票 66 证伪为"必红且连报告产物都没有"），追加不是改写。⚠ C.5 的诚实边界要复述一遍：**`-Subset full` 六态 posture 全是 skeleton**，所以那张表只证明"state 口径跑得通且 skeleton 达标"，**不等于 D32 六态表达标**（六态真测归票 15/28/33/36）。 |
| **#7** 门禁：`gofmt -l` 触及包为空、`go vet` 干净、`go test -count=2 ./internal/proc/ ./internal/observe/`、`-race`、d22scan **按自己模块目录跑**并自报 `examined N` | **PASS** | 〔编排者独立复现〕**23:55 当场重跑**：`gofmt -l internal/proc internal/observe cmd/wisp cmd/balldebug` → **无输出，rc=0**；`go test -count=2 ./internal/proc/ ./internal/observe/` → `ok 1.700s` / `ok 2.383s`（23:36 我自己跑的那次）。〔代理日志＋归档，我抽验〕`-race -count=1` 两包各 `ok 2.955s`；d22scan `examined 194 production Go files` + `clean` rc=0，**并做了阳性对照**（`TestScanDetectsAllSeededViolations`、`TestCheckRootRejectsBlindRoots`（含"在仓根调用＝盲跑必须致命"）、`TestScanCleanRepoIsGreen` 全 PASS）⇒ 这句 clean 不是"没扫"扫出来的。⚠ 23:55 我**没有**重跑 `go vet`/测试于**全仓**，因为票 70 正在把 68 个 `.go` 文件格式化，全仓结果此刻不代表票 66。 |

## 结论

**7/7 PASS，无 BLOCKER、无 MAJOR。** 票 66 置 `done`（文件改名 `66-slo-instrument-parse-cpu-observer-done.md`）。

**明写残口，别让 done 变成零残余的错觉**（README 规则 4 的"状态回写"教训，A30）：

1. `NtQuerySystemInformation: buffer never sufficient (last 1090464 bytes)` —— 200 次读里 1 次
   （`66b-nokeeper-4.json` 那次），**未修**。不是 A14，影响是单次样本被丢弃（九门仍出数）。
   ⇒ 待登记为 registry 新项；修法候选是缓冲增长策略而不是重试次数，**归票 71 之后的仪器批次**。
2. **D32 六态表本身仍未测**（`-Subset full` 全 skeleton）⇒ 归票 15/28/33/36。
3. **树内口径的 CPU 门"无定义"这个判定本身**依赖"观测者成本可分离"这个前提；
   如果哪天观测者与主体共享了采样线程，`observer_cost:true` 那行就不再是自明的。
   ⇒ 这条由 `internal/observe/observer_cost_test.go` 4 条用例钉住，改口径的人必须先看它。
4. 票 66 **没有新增 Go 生产代码之外的东西**：唯一一次改 `.go` 是变异，已撤（`grep -c MUTATION`=0 + `git diff` 空）。

**这份报告本身的可疑处，我也标出来**：AC#1 的变异检验我**没有独立重做**（落在"〔仅代理自述＋归档〕"档），
理由是它与票 70 的全仓格式化互斥。**票 70 落地后应做一次廉价复核**：把
`WalkSystemProcesses` 的顺序退回旧实现，确认 `TestParseSystemProcessesKeepsLastSnapshotEntry` 变红。
这一条我登记进 `docs/reports/pending-and-issues.md` 的 A30 追加行，不留在本文档里当口头债。
