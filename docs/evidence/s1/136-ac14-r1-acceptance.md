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

---

## §2 判据②【主靶】独立翻那枚开关 ＋ 交两张名册 —— 〔独立复现〕**成立**

**只在仓外快照里翻**：`git archive aef82f5` → `D:\tmp\wisp136ac14r1\snap-gate`，
`sed -i '556s/false/true/'`。跟踪树**未参与**：翻完之后 `git status --porcelain internal/observe cmd/wisp` ⇒ **0 行**，
本程在跟踪树上对 `internal/observe/**` 一字节未写、也未 commit 任何一支的 Gate。

**先证落地再读数**（本仓口径）：
`grep -n "const settleCoverageRowGates" snap-gate/internal/observe/sampler.go` ⇒ **`556:const settleCoverageRowGates = true`**；
`diff -r snap-gate/internal/observe ../snap-post/internal/observe` ⇒ **只有 556c556 一处**（其余一字节同）；
`go build ./...` ⇒ **rc=0**。三条都在，才读下面的数。

命令：`go test -count=2 -v ./internal/observe/`（同一快照，一发）
⇒ **RUN=142 / PASS=138 / FAIL=4 / SKIP=0**，去重顶层 **71** 枚（红 **2** 枚 ×2 计数 ＝ 4 行 FAIL），耗时 6.762s。
`grep -c '^--- SKIP'` ＝ **0** ⇒ 没有一名用 SKIP 换色；红∩绿 ＝ **0**；名册逐名与 §1 的 71 枚**同一份**
（`comm -3` 只吐出那两枚由 PASS 转 FAIL 的名字，**无人消失**）。

### (A) 红名册"恰好"＝哪几枚（**2 枚**，顶层去重）

| # | 用例 | 断言站点（本程从 `-v` 输出逐字取出） | 红句（原文截断处） | 报告实值 |
| --- | --- | --- | --- | --- |
| A1 | `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | `sampler_settle_coverage_136_test.go:144` | `this leg pins disclosure, not the verdict; a covered-enough window must still pass` | `SampleErrors:9`、行 `Pass:false Gate:true`、`Pass:false` |
| A2 | `TestCheckSettleZeroFootprintDropsAreCountedToo` | `sampler_settle_coverage_136_test.go:281` | `report=&{TargetState:Sleeping …}` | `SampleErrors:9`、行 `Pass:false Gate:true`、`Pass:false` |

⇒ **实现件 §4 变异 S1 自报"恰好红 2 枚"成立**（本程独立翻，枚数与逐名都对上；红因是 `Pass` 被折叠否决，
不是编译破裂、不是前提腿）。⚠ 顺带更正实现件自己的一处口径混用：它 §1.3 写"批准之后会红哪**三**枚"，
§4.1/§4.2 又写代价＝**2** 枚——**2 枚**是盘上值，"三枚"是它在改动之前数既有断言的口径（它 §4.2 自己已改写）。

### (B) 必须仍绿的名册（翻 Gate 后**实测仍绿**，本程逐名点名）

**尺的用法**：(B) 是"折叠过头"的探测器——若下列任何一枚在翻布尔后转红，就说明折叠否决了**不该被否决**的东西，
门行的判据写宽了。本程读数：**下列全部仍绿**（每枚 `--- PASS` 出现 2 次 ＝ `-count=2`，`--- FAIL` 0 次）。

| # | 用例 | 站点 | 为什么它必须仍绿 |
| --- | --- | --- | --- |
| B1 | `TestCheckSettleFullyMeasuredWindowReportsNoLoss` | `sampler_settle_coverage_136_test.go:307`（红句 `a measurable, settled, released window must pass`） | **健康窗口**：`SampleErrors=0`、行自己 pass ⇒ 折叠无东西可否决。**它就是简报与 `A182⑤` 要的那一枚**，实测绿 |
| B2 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss` | 同文件 `:208-236`（AC#14 改写后的那一形） | 改写刻意写成**与 gate 取值无关**（要求"行存在＋行自己说 not-pass＋印 `sample_errors`＋不得出现失败 gate 行与 `pass=true` 并存"）。翻布尔后仍绿 ⇒ 印证实现件 §4.2"这不是漏，是设计"这句**成立** |
| B3 | `TestSettleCoverageRowExistsAndPassesWhenFullyMeasured` | `sampler_settle_gate_136_test.go:154`（其 `:179` 也是 `a measurable, settled, released window must pass`） | 满测窗口的正对照 ＋ 形状钉 |
| B4 | `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` | `_gate_:198`（`:225` 写作 `if row.Gate && rep.Pass`） | 同一形制：两味取值都成立 |
| B5 | `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | `_gate_:236` | 行只判覆盖率，与是否否决无关 |
| B6 | `TestFoldSettlePassOnlyGateRowsVeto` | `_gate_:281`（**7** 个用例，本程逐字数过：含"失败 gate 必须否决"与"失败记录行不得否决"两形） | 钉的是折叠**规则**本身 |
| B7 | `TestSettleReportPassNeverContradictsItsGateRows` | `_gate_:312` | 跨一致性腿 |
| B8 | `TestStateReportVerdictBuilderStaysSinglePurpose` | `_gate_:337` | 冻结面守卫，翻 Gate 不该碰它 |
| B9 | 其余 **60** 枚与 settle/gate 无关的顶层用例 | — | 全绿（71 − 2 红 − 9 点名 ＝ 60） |

⇒ **(A) 2 枚 ＋ (B) 至少 9 枚点名 ＋ 总数守恒（69 绿 / 2 红 / 0 SKIP / 71 名册不变）**。
**判据② 的两张名册本程都交得出 ⇒ 不触发"答不出 (B) 就退回"那一支。**

**翻 Gate 的争用/前提自查**：该发日志内 `grep -c "precondition broken: only"` ⇒ **0**（AC#15 那枚已知 flake
本程四发主读数＋§4 四发变异里**一次没响**）；真 `^panic:` ⇒ 0。

**§2 的 commit 回显**（原样贴，供核"每枚 commit 只带自己那枚路径"）：

```
$ git log --oneline -1
bb723f7 accept(136 AC#14 r1 §2): 主靶——仓外快照独立翻 settleCoverageRowGates，红名册(A)恰 2 枚、必须仍绿名册(B)点名 9 枚，71 名册守恒 0 SKIP
$ git show --name-only HEAD
docs/evidence/s1/136-ac14-r1-acceptance.md
```

---

## §3 判据③ 复核"三处同形"——它推翻过一张落进票面的裁定 —— 〔独立复现〕**结论成立，但简报给的"三处"这个数不成立**

被验版本：`git show ca2c55e:internal/observe/sampler_settle_coverage_136_test.go`（该文件在 `ca2c55e` 的
`git log -1` ＝ 本程即用此版本；全文 289 行）。

**先给本程自己量的原始数**：`grep -c "if !rep.Pass {"` ⇒ **4**，逐处行号 **`:143` / `:208` / `:252` / `:279`**。
⇒ **简报里"数 `if !rep.Pass {`（意在得三处）这句按字面跑出来是 4 枚，不是 3 枚。**
差的第 4 枚（`:279`）不是漏网的同形，它正是 (B) 名册那一枚健康窗口。**"三处同形"这句作为"三枚要求'部分未测'的窗口必须 pass"的实质论断成立；作为 grep 计数不成立。** 下表逐处给改前行号＋红句原文＋归类：

| 改前行号 | 所属用例（func 行） | 红句原文（逐字） | 归类（本程判） |
| --- | --- | --- | --- |
| `:143-145` | `TestCheckSettleSingleTrustworthyReadReportsItsLoss`（`:123`） | `this leg pins disclosure, not the verdict; a covered-enough window must still pass: %+v` | **部分未测却要求 pass**（10 枚读数只 1 枚可信、丢 9）⇒ **该被门抓** |
| `:208-210` | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（`:176`） | `disclosure leg, not a verdict leg: this window still passes, report=%+v` | **部分未测却要求 pass**（丢一半，`kept≈lost`）⇒ **该被门抓；＝票面 `:287` 唯一预先授权改的那一处** |
| `:252-254` | `TestCheckSettleZeroFootprintDropsAreCountedToo`（`:238`） | `report=%+v` | **部分未测却要求 pass**（1 可信 / 2 枚零足迹、丢 9）；红句没写理由，但形状与判据同 ⇒ **该被门抓** |
| `:279-281` | `TestCheckSettleFullyMeasuredWindowReportsNoLoss`（`:265`） | `a measurable, settled, released window must pass, report=%+v` | **健康窗口要求 pass**（`sample_errors=0`、reads≥3 前提腿）⇒ **翻 Gate 后必须仍绿**＝§2 的 B1 |

改后（`aef82f5`）现量：`if !rep.Pass {` ⇒ **3 处**，行号 `:143` / `:280` / `:307`。
⇒ **票面 `:299` 那句"HEAD 上还剩 `:143`、`:280` 两枚要求 pass，另有 `:307` 一枚是健康窗口必须 pass"逐名逐行对上**；
`git diff -U2 ca2c55e..aef82f5` 于该文件**只有一枚 hunk `@@ -206,6 +206,34 @@`**，头注释 `:25-27`
（"a partially covered window still passes today (changing that verdict is not this cell's job)"）**逐字未动**
⇒ 实现件 §2.4 自述的"没动头注释、把它改准属 AC#15 射程"**成立**。

**"任何抓得住'丢一半读数'的门必然把这几形一起抓红"成不成立？——成立（本程既有测量支撑，也给了判据形状）。**

1. 门行的判据是 `sampler.go:565` `covered := rep.SampleErrors == 0 && len(rep.Samples) > 0`（本程盘上现量，行号已核）。
   三枚"部分未测"形的 `SampleErrors` 都 **>0** ⇒ **同一条件同时命中三形**，不存在"只命中一半那一形"的取值。
2. 想只红一处，只有两条路，都不许走：
   - **按丢读比例设阈值**（例：只红在 `lost/reads ≈ 0.5`）：`{1/10, 5/10, 1/10}` 里要红中间那一枚而放过两端的 9/10，
     这是**非单调**判据 ⇒ 新造阈值，`AGENTS.md §1.1` 与票面 `:279` 的禁面。
   - **按丢读原因分流**：本程从 gate=true 的输出逐字取到——`:143` 腿与 `:208` 腿的 `LastSampleError` **同为**
     `read: settle probe: transient tree read failure`，只有 `:252` 腿是 `read returned a zero private working set for a live tree`。
     ⇒ 按原因**最多把"零足迹"那一形分开，永远分不开 `:143` 与 `:208`**（同原因、只差比例，又回到上一条）。
3. **实测印证**：§2 翻布尔后 → `:143`/`:252`（现为 `:280`）**两枚一起红**、健康窗 `:307` **仍绿**，
   一枚不多一枚不少。这正是"没有只红一处的中间形状"的读数形式。

⇒ **判定：票面 `:287`"派单预先授权它改**那一处**断言"这句在盘上不可满足**——任何满足 AC#14 判据①②的门行，
都会同时把另两形抓红，而那两形不在授权面内。实现程选择"落在记录行＋停手报回"是**对的动作**，
不构成"擅自扩面"，也不构成"放水"。
**这一处是编排者的纸的缺陷，账已在本仓口径下成立**：`docs/reports/pending-and-issues.md` **`A182③`**
（量于 `ddbd3a1`，该条目起于 `:5580`）已如实写下"它推翻我那句'一处断言'，成立"。
**本程独立复核后确认 `A182③` 站得住**，并补两条 `A182③` 没量的数：
① 改前 `if !rep.Pass {` 的**原始计数是 4 不是 3**（第 4 枚是健康窗，`A182③` 与票面 `:296` 数的是"三处同形"，
两个口径都真，但**今后谁引这句得连口径一起引**，否则下一位 `grep -c` 会以为少了一枚）；
② "分不开 `:143` 与 `:208`"的**原因串同一**这条证据（上面第 2 点）在 `A182③` 与实现件 §1.3 里都没有，
本程补上——它把"只剩非单调阈值一条路"从**论述**变成了**有串的读数**。

