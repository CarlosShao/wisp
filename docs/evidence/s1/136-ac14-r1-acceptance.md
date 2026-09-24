# 136 — AC#14 终裁 r1（非实现者）**acceptance**

> 裁决者 `acceptor-ticket136-ac14-r1`，非实现者（实现程＝`worker-ticket136-ac14`）。
> 本格只裁 **AC#14**（`.scratch/wisp/issues/136-...md:251-271` ＋ `:273-292` 的具名解冻＋四裁定 ＋
> `:294-310` 的 19:2x 交回裁定）。**本程不翻任何勾、不裁"本格成立"以外的动作、不 push。**
> **被验码 = `aef82f5`**（三枚路径全在 `internal/observe/`：`sampler.go`、
> `sampler_settle_gate_136_test.go`（新建）、`sampler_settle_coverage_136_test.go`）。

## §0 锚点、口径、可复算性

| 项 | 读数（本程自己量） |
| --- | --- |
| 开工自量锚点 | `git rev-parse --short HEAD` ⇒ **`6451625`**（第一次量）；**程中漂到 `ddbd3a1`**（连量三发同为 `ddbd3a1`）。派单简报写的 `98665ef` 是**中间某一发**，不是我的锚。 |
| 被验码与 HEAD 的关系 | `git merge-base --is-ancestor aef82f5 HEAD` ⇒ 是祖先；`git diff --stat aef82f5..HEAD -- internal/observe/` ⇒ **空** ⇒ 在 HEAD 上量 `./internal/observe/` ＝ 在 `aef82f5` 上量（逐字节同码）。本程四数即在 HEAD 取，变异/翻布尔全在 `aef82f5` 快照取。 |
| `aef82f5` 写集 | `git show --name-only aef82f5` ⇒ 恰三枚路径，全在 `internal/observe/` ✓（与简报一致） |
| 取数用文件版本 | `internal/observe/sampler.go`、`sampler_settle_coverage_136_test.go`、`sampler_settle_gate_136_test.go` 的 `git log -1` 均为 **`aef82f5`（Thu Sep 24 19:08:12 2026 +0800）**；票面 `136-...md` 的 `git log -1` 为 **`2f22ab3`**，`docs/reports/pending-and-issues.md` 为 **`ddbd3a1`（19:32:50）**。下文每处 file:line 都在这些版本上量。 |
| 工具链 | `go version` 见 §5；`gofumpt` 盘上现量在 §5（本程不抄任何人的版本号） |
| 临时件落点 | 全部在 `D:\tmp\wisp136ac14r1\`（`snap-pre`＝`git archive ca2c55e`、`snap-post`＝`git archive aef82f5`、`snap-gate`／`mut/*`＝翻布尔与逐发变异、`logs/`＝逐发 `-v` 日志）。**仓库目录内除本文件外未新建任何东西。** |
| 工作树纪律 | 只 commit 本文件一枚路径；`internal/observe/**` 与其他禁改面**一字节未写**；owner 的 `design/**` 未提交删除与两枚未跟踪目录未碰、未还原、未提交 |

**快照自证**（防"静默挂空＝假绿"那一形）：`snap-pre/go.mod`、`snap-post/go.mod` 均存在非空，且
`diff -q snap-post/internal/observe/sampler.go <跟踪树>/internal/observe/sampler.go` ⇒ **IDENTICAL**
（量于 HEAD `ddbd3a1`）。

---

## §1 判据① 复现四数与名册 —— 〔独立复现〕**成立**

命令（本程自己跑，两遍，顺序执行以免自相争用）：`go test -count=2 -v ./internal/observe/`，
树＝跟踪树 HEAD（`internal/observe/` 与 `aef82f5` 逐字节同，`git status --porcelain internal/observe` 量前后均 0 行）。

| 口径 | run1（本程） | run2（本程） | 实现件自报 |
| --- | --- | --- | --- |
| `^=== RUN` | **142** | **142** | 142 |
| `^--- PASS` | **142** | **142** | 142 |
| `^--- FAIL` | **0** | **0** | 0 |
| `^--- SKIP` | **0** | **0** | 0 |
| 去重后顶层名数 | **71** | **71** | 71（顶层 71 枚 ×2） |
| 整包耗时 | 6.743s | 6.799s | — |

⇒ **142/142/0/0 ×2 逐格对上**；`-count=2` 的乘子确实存在（**71 × 2 = 142**），报数不带乘子即为误读。

**基线**（`git archive ca2c55e` → `D:\tmp\wisp136ac14r1\snap-pre`，同一条命令 `-count=2 -v`）：
**RUN=130 / PASS=130 / FAIL=0 / SKIP=0，去重 65 枚**，耗时 3.848s ⇒ 实现件 §0.1 的
"基线 65×2=130、净增 6"**逐格对上**。

**名册差集（除四数外的那一维，本程自己算）**

| 比对 | 命令 | 读数 |
| --- | --- | --- |
| 两遍之间 | `grep '^--- ' \| awk '{print $2,$3}' \| sort` → `comm -3` | **0 行** ⇒ 两遍逐名相同 |
| 基线→改后「谁消失了」 | `comm -23 pre.u post.u`（两侧同用 `sort -u`） | **0 枚** ⇒ 无人消失 |
| 基线→改后「谁新增」 | `comm -13 pre.u post.u` | **恰 6 枚**，逐名＝`TestFoldSettlePassOnlyGateRowsVeto`、`TestSettleCoverageRowExistsAndPassesWhenFullyMeasured`、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`、`TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured`、`TestSettleReportPassNeverContradictsItsGateRows`、`TestStateReportVerdictBuilderStaysSinglePurpose` ⇒ 与实现件 §3 那张新增清单**逐字同** |
| 绿转 SKIP | `grep -c '^--- SKIP\|^    --- SKIP'` 三发日志 | **0 / 0 / 0** ⇒ 无任何一名转 SKIP |

⚠ **本程自己制造过一枚假读数并已作废**：第一次算差集时基线侧用了 `sort -u`（65 行）而改后侧用了
不带 `-u` 的 `sort`（142 行），`comm` 两侧归一化不一致 ⇒ 吐出一大堆"新增"（含 78 行既有名）。
**那一发已作废**，上表是两侧都 `sort -u` 之后重跑的。登记在此，免得下一位以为差集天然干净。

**panic 那一维复核**（简报要我复核"8 处全是用例名带 Panic"这句）：
`grep -ci panic` 于两发日志各 **8 行**；`grep -c '^panic:'`（真 panic）＝ **0**。
那 8 行逐行为 `TestPanicInWorkerSurvivesAndCancelsRoot` 与
`TestFakeTreeEmptyScriptFailsClosedAndNotPanics` 各 2 事件（`=== RUN` ＋ `--- PASS`）×2（`-count=2`）＝ 8。
⇒ **实现件 §0.1 那句解释成立**：本包**没有真 panic**，8 是名字里的 `Panic`。
这一条要紧的因果本程独立核过：真 `^panic:` 会把同包其余几十条读数一起吞掉、包级 rc 只会说"这一包失败"，
而本包两遍都是 **142/142 全点名、FAIL=0**，与"零真 panic"自洽。

**AC#15 那枚已知 flake 这一发响没响**：**没响**。`grep 'precondition broken: only'` 于
`head-v-1.txt`/`head-v-2.txt`/`pre-v-1.txt` ⇒ **0 命中**；三发 FAIL 均为 0。
本程额外把它记在样本量上：改后 **142 发**（两遍 `-count=2`）＋ 基线 **130 发** ⇒ 本程这一程共 272 发顶层读数，
**0 发碰到那枚前提腿**。⚠ 这是"没响"，不是"已修"——票面 AC#15 那格（`:325-`）的复现率判据
（要 n、要分口径）**本格不裁、也不许被这句当结论地基**。后续 §2/§4 的每一发变异日志本程同样逐名数过，
那几发的账记在 §2/§4。

**档位**：**〔独立复现〕判据① 成立**（四数 ×2、名册三向差集、panic 两味、flake 未响，全部本程自己跑出来的数）。
