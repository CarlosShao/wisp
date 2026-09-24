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

**§3 的 commit 回显**：

```
$ git log --oneline -1
ffe5731 accept(136 AC#14 r1 §3): 复核"三处同形"——改前 if !rep.Pass 原始计数 4（第 4 枚＝健康窗），三枚 disclosure 形归类逐名给；"只红一处"不可满足成立，A182③ 确认并补两数
$ git show --name-only HEAD
docs/evidence/s1/136-ac14-r1-acceptance.md
```

---

## §4 判据④ 承重判定（真摘真跑，全在仓外快照）—— 〔独立复现〕**成立**

**口径**：五发变异，每发都是 `cp -r snap-post mut/<名>`（`snap-post`＝`git archive aef82f5` 的纯净快照）→
`sed` 改那**一行** → `grep -n`/`sed -n` 印出改后那行（下表"落地凭据"）→ `go build ./...` rc=0 →
`go test -count=1 -v ./internal/observe/`。**一律 `-count=1`，四数里的顶层数＝71（不带 ×2 乘子）**。
跟踪树在五发里都未参与（每发跑完 `git status --porcelain internal/observe` 复核 0 行）。

| 发 | 摘/改了什么 | 落地凭据（改后盘上现量） | build | 四数（RUN/PASS/FAIL/SKIP） | 红名（逐名） |
| --- | --- | --- | --- | --- | --- |
| **M1 摘门行本体** | 删 `sampler.go:529` `rep.Verdicts = buildSettleVerdicts(*rep)`（构造点） | 文件 597→**596** 行；`grep -n buildSettleVerdicts` ⇒ **只剩 3 处注释（`:460`/`:540`/`:557`）＋ `:563` 函数定义，构造点 0** | rc=0 | **71/67/4/0**（去重 71） | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`、`TestSettleCoverageRowExistsAndPassesWhenFullyMeasured`、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`、`TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` ⇒ **4 枚，与实现件 §4-M1 逐名同** |
| **M2 让门行不看丢读** | `:565` `covered := rep.SampleErrors == 0 && len(rep.Samples) > 0` → `covered := len(rep.Samples) > 0` | `sed -n '565p'` ⇒ `covered := len(rep.Samples) > 0` | rc=0 | **71/69/2/0** | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` ⇒ **2 枚，与实现件 §4-M2 逐名同** |
| **M3a 单点回退（gate 仍 false）** | `:537` `rep.Pass = foldSettlePass(memOK && backInTime && releaseOK, rep.Verdicts)` → `rep.Pass = memOK && backInTime && releaseOK`；`:556` 保持 `= false` | `sed -n '537p;556p'` ⇒ 前者已回退、后者仍 `false` | rc=0 | **71/71/0/0 ⇒ 全绿** | 无 |
| **M3a′ 两味一起动（gate=true 上重做 M3a）** | 同一处回退 **＋** `:556` = `true` | `sed -n '537p;556p'` ⇒ 回退行 ＋ `const settleCoverageRowGates = true` | rc=0 | **71/68/3/0** | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`、`TestSettleReportPassNeverContradictsItsGateRows` ⇒ **3 枚，与实现件 §4-M3b 逐名同** |
| **M4 摘"红句点名"（本程加的一发，简报未要求）** | `:570` 失败分支 note 串里 `sample_errors=%d` → `dropped=%d` | `sed -n '570p'` ⇒ `note = fmt.Sprintf("fail-closed disclosure: dropped=%d of %d reads taken (…` | rc=0 | **71/69/2/0** | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（站点 `coverage_136_test.go:232`）、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`（站点 `gate_136_test.go:215`） |

**五发的 SKIP 全 0、名册逐名守恒（71 枚一枚不少）** ⇒ 没有任何一发用 SKIP 或"用例变少"换色。
`grep -c "precondition broken: only"` 于五发日志 ⇒ **全 0**（AC#15 那枚已知 flake 在本程 5 发变异 ＋ §1 两发 ＋ §2 一发里都没响；
如实登记：**没响**，不裁它已修）。

### 4.1 承重结论（按"摘掉任意一味，是否存在一发变异从此打不红"逐味判）

| 味（`sampler.go` 落地点） | 判 | 依据读数 |
| --- | --- | --- |
| 门行**构造点** `:529` | **承重** | M1 ⇒ 立刻红 4 枚（含 AC#14 判据② 要的钉子 `TestCheckSettleHalfTheReadsFailedReportsItsLoss`） |
| 门行**看丢读**那半条件 `:565` `SampleErrors == 0` | **承重** | M2 ⇒ 红 2 枚；且红因正是"行说它测够了"那句自陈，方向对 |
| 折叠**调用点** `:537` | **终态承重、当前形状下按定义不承重** | M3a′（gate=true）⇒ 红 3 枚 ⇒ **翻布尔后它立刻变承重且有钉子看着**；M3a（gate=false）⇒ 0 红 |
| 开关常量 `:556` ＋ 行的 `Gate: settleCoverageRowGates` 接线 `:580` | **承重** | §2 只翻 `:556` 一字节 ⇒ 行为改变（红 2 枚），名册不漂 ⇒ 这根线是活的，不是死常量 |
| 红句**点名 `sample_errors`** `:570` | **承重**（本程新增证据） | M4 ⇒ 红 2 枚，站点逐名可查 ⇒ 票面 `:285` 那一维有钉 |
| 声明块注释改写（实现件改动 #2，`:464-467`） | **按定义装饰**，本程**同意实现件的自述** | 它是注释，不参与任何断言；实现件 §4.1 末行已自报"按定义装饰，不参与任何读数"，本程不重复记账 |

**⚠ 本程按票面 `A182⑧` 预先挡住的那条误判，这里明确不犯**：M3a（gate=false 撤折叠）⇒ 71 枚全绿，
**不是**实现件的缺陷。`settleCoverageRowGates=false` ⇒ 门行 `Gate=false` ⇒ 折叠**按定义**永不否决 ⇒
"撤了没反应"是**编排者暂不批 Gate 的后果**。本程**没有**拿它判"装饰"而退回，而是按简报要求
**在 `gate=true` 的快照上重做同一发（M3a′）**：那一发红 3 枚 ⇒ **折叠那行在终态里不是装饰**。

### 4.2 关于"不许要求同一发变异里新加那支先响"

本程在 M1/M2 里没有要求"新加那支先响"作判据：M1 的红 4 枚里既有既有钉（`:226` 那一族）也有新探针，
M2 的 2 枚同理——**承重判据是"摘掉后存在一发变异打不红"，本程五发都取到了红，唯 M3a 那一发的"不红"已由 M3a′ 补掉**。
⇒ 判据④ **成立**，实现件 §4 的六发本程独立复现到**逐名＋逐站点同**：M1（红 4，站点 `coverage:226`／`gate:162`／`gate:207`／`gate:245`）、
M2（红 2，站点 `coverage:229`／`gate:209`）、M3a（全绿）、M3b＝本程 M3a′（红 3，站点 `coverage:236`／`gate:226`／`gate:325`）、
S1＝本程 §2 的翻布尔（红 2，站点 `coverage:144`／`coverage:281`）；
**正向对照**那一发本程未单独重跑，它就是 §1 在跟踪树（＝`aef82f5` 同码）上取的 **142/142/0/0**——
形状同一（未改动的码 ⇒ 全绿），本程以此为准，**不另立一发冒充独立复现**。

**§4 的 commit 回显**：

```
$ git log --oneline -1
24ab73f accept(136 AC#14 r1 §4): 承重判定五发变异（M1 红4/M2 红2/M3a 全绿/M3a′gate=true 红3/M4 摘红句点名 红2），逐站点同，A182⑧ 预挡的误判未犯
$ git show --name-only HEAD
docs/evidence/s1/136-ac14-r1-acceptance.md
```

---

## §5 判据⑤ 门禁与格式 —— 〔独立复现〕**成立**（未放宽、未转 Skip、未动阈值/golden）

工具链本程盘上现量：`go version go1.27.1 windows/amd64`；
`"$(go env GOPATH)/bin/gofumpt.exe" -version` ⇒ **`v0.12.0 (go1.27.1)`**
⇒ **简报写的"盘上现量 v0.12.0"对上；票面/旧报告里的 v0.7.0 确为过期值**（本程不引任何人的转述版本号）。

| 门 | 命令（本程自己跑） | 读数 |
| --- | --- | --- |
| gofmt | `gofmt -l internal/observe` ＋ 三枚显式路径单列 | **空输出，rc=0** |
| gofumpt | `gofumpt -l internal/observe` ＋ 三枚显式路径单列 | **空输出，rc=0** |
| vet | `go vet ./internal/observe/` | **rc=0** |
| 构建 | `go build ./...`（`snap-post`／`snap-gate`／五发变异树各自） | **全 rc=0** |
| d22scan | `git archive` 仓外纯净快照上 `sh scripts/d22scan.sh`：`snap-pre`＝`ca2c55e`、`snap-post`＝`aef82f5` | **两发各 rc=0** |

**八 scope 命中数（本程从 d22scan 自己的 "live scope work" 行逐字取，不抄任何表）**

| scope | `ca2c55e`（改前） | `aef82f5`（改后） | 判 |
| --- | --- | --- | --- |
| bans #1-5 `internal/` | 203 | 203 | 不降 |
| bans #1-5 `cmd/` | 22 | 22 | 不降 |
| ban #6 `frontend/` | 40 | 40 | 不降 |
| ban #7 `internal/tools/` | 18 | 18 | 不降 |
| ban #8 `design/` | 16 | 16 | 不降 |
| ban #8 `frontend/` | 40 | 40 | 不降 |
| ban #8 `internal/` | **404** | **405** | ＋1＝新增那枚 `_test.go`，**不降** |
| ban #8 `cmd/` | 39 | 39 | 不降 |

正对照（同两发日志内）：`runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` ⇒ 两发同值。
台账口径出处本程也现量了：`docs/reports/pending-and-issues.md`（量于 `ddbd3a1`）**`:3922`** 记
"台账八 scope `203/22/40/18/16/40/390/37` 逐格不降"——后两格（390/37）**低于**本程今量值（404→405/39），
即台账那一份是**当天更早的锚**、方向仍是涨不是降。⚠ 因此"不降"这一判**必须连锚一起引**：
本程用的是**同一把尺的两枚相邻快照**（`ca2c55e` → `aef82f5`），不是"跟任意历史值比"。

**"不许为了变绿放宽断言／改成 Skip／动阈值·golden"这一维，本程逐条量过（不是采信自述）**：

| 检查 | 读数 |
| --- | --- |
| 新增 `t.Skip` | `git diff ca2c55e..aef82f5 -- internal/observe/ \| grep -c "^+.*t\.Skip"` ⇒ **0**；三枚路径全文 `t.Skip` ⇒ **0/0/0** |
| 名册里转 SKIP | §1/§2/§4 九发日志 `^--- SKIP` ⇒ **全 0** |
| `thresholds.go` | `git diff --stat ca2c55e..aef82f5 -- internal/observe/thresholds.go` ⇒ **空** |
| 冻结面（`cmd/wisp`、`frontend`、`scripts`、`.github`、`docs/PLAN.md`、`docs/specs`、`internal/risk`、`tools/d22scan`） | 同一条 `--stat` 一次问八枚路径 ⇒ **全空** |
| 被改的那枚断言（`:208-210`）是放宽还是收紧 | `git diff -U2` ⇒ 该文件**只有一枚 hunk `@@ -206,6 +206,34 @@`**；删的是 `if !rep.Pass { t.Fatalf("…this window still passes…") }` 两行，加的是四道**更强**的钉（行必须存在／行必须自己 `Pass==false`／红句必须含 `sample_errors=`／"失败 gate 行与 `pass=true` 并存"必须永不出现在同一枚报告里）＋ 8 段说明注释 ⇒ **收紧，不是放宽** |
| `ca2c55e..aef82f5` 整仓 `--stat` 落在 `internal/observe/` 的路径 | **恰 3 枚**（`git diff --numstat`：`sampler.go` **＋77/−2**、`sampler_settle_coverage_136_test.go` **＋30/−2**、`sampler_settle_gate_136_test.go` **＋349/−0**）⇒ 与 `git show --numstat aef82f5` 逐枚同、与 `git show --name-only aef82f5` 的三枚路径一致，无第 4 枚；也与 `A182①` 的"77 增 2 删／30 增 2 删／349 增 0 删"对上 |

**emoji／卫生那一维的一处自述不成立（不属门禁失败，属自述不准）**：实现件 §3 写"新码与新测试文件的注释**逐字 ASCII**"。
本程现量：`sampler.go` 新增段（`:455-612`）非 ASCII 行 **0** ✓；但新文件
`sampler_settle_gate_136_test.go:152` 注释含中文词 **`同形`**（非 ASCII）。
⇒ **这句自述不成立**。要紧的是它**不破任何门**：`同形` 落在 CJK 段，不在 `AGENTS.md §1.2` 禁的
`U+2190–U+2BFF`／`U+1F300–U+1FAFF`／`U+FE0F` 里，且 `ban #8 internal/` 405 枚文件 rc=0 ＝ **外部读数背书**。
登记为**自述精度缺陷**，不构成退回理由。

**争用与机时**：本程**没有取任何一发 `wisp slo -settle` 真取样读数**（五条判据都不要求），
⇒ `scripts/slo-check.ps1` 那份 15 枚争用名单（本程现量：**块在 `:153-155`，不是简报与 `A182⑦` 写的 `:155-157`**；
枚数 **15 枚**逐名对上：`go`/`gofmt`/`cgo`/`compile`/`asm`/`link` ＋ `gcc`/`g++`/`cc1`/`cc1plus`/`as`/`ld` ＋ `wisp`/`wisp-cli`/`staticcheck`；
`slo-check.ps1` 的 `git log -1` ＝ **`decb7b9`（Wed Sep 23 13:51:05 2026）**，即那处行号偏差**不是漂移、是引错**）
对本程不构成取数约束。**反向账本程如实交**：本程在 `19:3x–20:0x` 于这台 6C12T 上跑了 9 发整包 `go test` ＋ 6 发 `go build ./...`
⇒ **同一时间窗内别的程在这台机取真数，那些数按 `slo-check.ps1:153` 的口径应判"无效样"**（与实现件 §6-4 同一条口径）。

**§5 的 commit 回显**：

```
$ git log --oneline -1
512d33d accept(136 AC#14 r1 §5): 门禁与格式（gofmt/gofumpt v0.12.0/vet/build 全 rc=0、d22scan 两快照 rc=0、八 scope 逐格不降 404->405、零新增 Skip、阈值与八面冻结路径 diff 全空、记一处自述不成立=同形 CJK）
$ git show --name-only HEAD
docs/evidence/s1/136-ac14-r1-acceptance.md
```

---

## §6 五条判据总裁 ＋ 本格状态

| 判据 | 本程档 | 裁 |
| --- | --- | --- |
| ① 复现四数与名册 | **〔独立复现〕** | **成立**：142/142/0/0 ×2、顶层 71（×2 乘子）、基线 130/65、名册三向差集（两遍 0 行／消失 0／新增恰 6 逐名同）、`SKIP=0`、`grep -ci panic`=8 而真 `^panic:`=0 **且那 8 行确为两枚用例名**、AC#15 那枚已知 flake **本程 9 发未响** |
| ②【主靶】独立翻布尔＋交两张名册 | **〔独立复现〕** | **成立**：(A) 红名册**恰好 2 枚**（站点 `:144`／`:281`，逐字红句与本程报告值都在 §2）；(B) 必须仍绿**点名 9 枚**（含 `:307` 健康窗与 `:208` 改写腿）＋总数守恒 69/2/0/71。**两张都交得出 ⇒ 不触发"答不出 (B) 就退回"** |
| ③ "三处同形"复核 | **〔独立复现〕** | **实质成立、计数口径要改**：`ca2c55e` 上 `if !rep.Pass {` **原始 4 枚**，其中 **3 枚**是"部分未测却要求 pass"（`:143`/`:208`/`:252`）、第 4 枚 `:279` 是健康窗＝(B) 那一枚。**"只红一处"的中间形状不存在**（本程给了两条：比例判据必非单调；按丢读原因分流**分不开 `:143` 与 `:208`，二者 `LastSampleError` 串逐字同**）⇒ **票面 `:287`"预先授权改那一处"在盘上不可满足，是编排者的纸的缺陷；`A182③` 独立确认成立**。实现程停手报回**动作正确**，不判放水、不判越界 |
| ④ 承重判定 | **〔独立复现〕** | **成立**：M1 红 4／M2 红 2／M3a(gate=false) 全绿／**M3a′(gate=true) 红 3**／M4（本程加，摘红句点名）红 2。**没有一味是"摘掉后仍然全绿"**；折叠调用点在**终态承重**已用真读数钉死 ⇒ **本程没有、也不该拿 M3a 判它装饰**（`A182⑧` 那条误判已挡） |
| ⑤ 门禁与格式 | **〔独立复现〕** | **成立**：gofmt/gofumpt（盘上 **v0.12.0 (go1.27.1)**）/vet/build 全 rc=0；d22scan 两枚仓外纯净快照各 rc=0，八 scope 逐格不降（唯一变化 `ban #8 internal/` **404→405**＝新增那枚 `_test.go`）；**零新增 `t.Skip`、名册零转 SKIP、`thresholds.go` 与八面冻结路径 diff 全空、被改那枚断言是收紧不是放宽** |

### 6.1 本格（AC#14）裁：**附条件**——条件在编排者手里，不在实现件手里

- **落到的**：出线带与 `sampling` 同形的自陈门行（`Verdicts []Verdict json:"verdicts"` ＋ 行名 `sampling` ＋ `Measured` 逐字沿用 `StateReport` 那行的 `"%d valid / %d errors"` 格式，`sampler.go:577` 对 `:336`）；"部分未测"那一形**由那一行自己说 not-pass**并点名 `sample_errors`；钉子齐（M1/M2/M4 各响各腿）。⇒ 票面 `:262-263` 判据① 的**前半**＋ `:285` 的"红句必须点名"**已成立**。
- **没落到的**：判据① 后半"不许只靠 `pass=false` 一个总布尔代答"里**更强**的那一支——门行今天**不否决**总布尔（`Gate=false`）。⇒ **本格不算结案**，实现件自己 §6-1 也未翻勾、盘上一致（票 136 仍 5 勾／10 未勾，`AC#14` 仍 `[ ]`）。
- **本程不把它记成实现件的债**：那一步的唯一前置是**编排者对 `settleCoverageRowGates` 的裁定**（票面 `:310`／`A182⑤` 白纸黑字"那一支的授权在我手里"）。
- **若批 Gate，代价本程已量成读数**（不是推理）：生产码动 **1 枚布尔字面量**（`sampler.go:556`）＋ 改写 **2 枚既有断言**（`sampler_settle_coverage_136_test.go:143-145` 与 `:280-282`，改法＝本程实测已被 `:208-236` 那一形验证过的"两味取值都成立"写法）；**红名册恰好 2 枚、健康窗必须仍绿**，两张名册都在 §2。
- ⚠ **档三要写清**：批 Gate 那张纸上唯一的**真取样**前置＝实现件 §1.2 那六发 `wisp slo -settle` 的 `sample_errors` **逐发 0**。那一维**本程未复现**（不占机、不在五条判据内）⇒ 按档位记 **〔仅自述，不背书〕**。编排者若据 §2 批 Gate，"不误伤存量合法窗口"这一支**仍立在它自己那六发上**，别当成终裁背过书。

### 6.2 顺带交回给 AC#15 的一枚现量（票面 `:288` 要"到时要现量重划，不许照抄今天的行号"）

前提腿（会偶发红那一形）**今量共 7 枚站点，分布在两枚文件**：

| 站点 | 红句 | AC#14 前后 |
| --- | --- | --- |
| `sampler_settle_coverage_136_test.go:189` | `precondition broken: only %d reads taken, half-and-half needs a window to lose in` | **行号未漂**（`ca2c55e` 上也是 `:189`）⇒ 票面 AC#15 格引的 `:189` **今仍有效** |
| `sampler_settle_coverage_136_test.go:296` | `precondition broken: only %d reads taken`（健康窗那条） | **从 `:268` 漂到 `:296`** ⇒ 谁引旧号要改 |
| `sampler_settle_gate_136_test.go:157/160/202/240/243/342` | 六枚**新增**前提腿（含实现件 §6-3 自报的那枚 `kept>=1 && lost>=1`） | **AC#14 新带进来的暴露面**，AC#15 重划靶子时必须一起数（实现件 §6-3 只数了自己那 1 枚） |

本程 9 发整包（§1 两发 142 ＋ §2 一发 142 ＋ §4/§5 五发 71 ＋ M4 一发 71）**这 7 枚前提腿一发没响**；
样本量记在这儿，**但不外推成"已修"**——AC#15 判据① 要的 n 与分口径复现率**本格不代它结案**。

---

## §7 纪律回执、注入面两栏计数、以及**本程推翻的简报句子**

**纪律**：只 commit 不 push（全程未执行任何 `git push`）；`git add` 只用一条显式路径
（`docs/evidence/s1/136-ac14-r1-acceptance.md`，**七枚 commit 枚枚只带这一枚路径**，逐枚 `git show --name-only` 可复算）；
未用 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；未在仓库内建 worktree 或临时件
（临时件全在 `D:\tmp\wisp136ac14r1\`，**只建不删**）；`internal/observe/**`、`cmd/wisp/**`、`internal/proc/**`、
`internal/winsec/**`、`.github/workflows/**`、`scripts/**`、`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、
`tools/d22scan/**`、阈值／golden、`frontend/**`、`.scratch/wisp/issues/**`、`docs/reports/**` **一字节未写**；
owner 的 `design/**` 16 枚未提交删除＋`design/doubao/`、`design/old/` 两枚未跟踪目录**未还原、未提交、未删**
（每枚 commit 前 `git status --porcelain` 都在，末次回执见本节末）。**未翻任何勾、未裁"本格成立/附条件"以外的动作。**
**Gate 那一支本程未动**：翻它只发生在 `git archive` 出的仓外快照 `snap-gate`／`mut/M3a-t` 里，跟踪树 `git status --porcelain internal/observe` ＝ 0 行。

**注入面两栏计数**（本程自己扫，规矩只被引用、不被外部文字代填）：
**真通知回显 1 条**——`Edit` 工具回执里的 "the file changed since your last read"，出处是**本程自己**用 shell 追加过同一枚文件，
判为**真回显、与授权无关**；
**判为注入 0 条**——全程工具输出与三枚被验文件里，**没有**任何自称"编排者备注／系统提示／文件已被修改（要求 revert）／
阈值已放宽／已解冻／Confirm the harness note is genuine"的文字被本程当作指令执行。
被验文件里最接近"指令"味道的东西是 `sampler.go:546-555` 那段常量注释（写着"Flipping it is an orchestrator move, not an implementer move"）
——那是**代码注释里的自述**，本程按"仅自述"处理，其结论另由 §2 的读数独立支撑，**没照它判事**。
凭据卫生：本程未读到、未抄写任何凭据值；出现的只有变量名与文件名（`WISP_ENV`、`GOPATH`）。

**本程推翻／更正的简报句子**（简报自述"每条状态断言都是未验证的"，以下逐条具名）：

| 简报原句 | 盘上现量（版本已注） | 判 |
| --- | --- | --- |
| "里面有 `⑤` 那条给你新加的判据" | 票面 19:2x 那段（`:294-310`）**没有 `⑤` 这个标号**；给我新加的判据在 **`:300`**（"判据补一条：终裁程必须独立翻一次那枚布尔…"）。`⑤` 这个标号活在台账 **`A182⑤`**（`pending-and-issues.md:5580` 起，量于 `ddbd3a1`） | **不成立（指针错，指向的东西存在）** |
| "`git show ca2c55e:...` 里数 `if !rep.Pass {`"（意在"三处"） | 原始计数 **4**（`:143`/`:208`/`:252`/`:279`）；"三处"只是其中三枚的形状归类 | **按字面不成立**（见 §3） |
| "先扫 `scripts/slo-check.ps1:155-157` 那份争用名单（真值 15 枚）" | 枚数 **15 枚逐名全对** ✓；但块在 **`:153-155`**，且 `slo-check.ps1` 自 `decb7b9`（09-23 13:51）未动 ⇒ **不是漂移、是引错**（`A182⑦` 同错） | **枚数成立、行号不成立** |
| "flake … 账在票面 AC#15 格与 `:266`" | `:266` 是 **AC#14 的停手上报线**；"240 发里红 1 发"那句在 **`:328`**（AC#15 格第 4 行）。实现件 §6-3 也引了同一个错号 `:266` | **不成立（两处同错，登记一次）** |
| "我写这行时是 `98665ef`" | 本程首量 `6451625`、程中漂至 `ddbd3a1`（`98665ef` 确在历史里，是 `aef82f5` 之后的第 4 枚） | **成立（简报已自预告会漂）** |
| "实现件自报 142/142/0/0（71×2、基线 65×2=130、净增 6）"、"前者 8、后者 0"、"恰好红 2 枚"、"M3a 71 枚全绿"、"`:556` 常量／`:537` 折叠"、"被验码 `aef82f5` 三枚路径全在 `internal/observe/`" | 六条**本程全部独立复现到逐名同** | **成立** |

**实现件里另三处自述不成立（不改变任何一格的判，但都会被下一位当事实引，故登记）**：
① §2.1 落点表的行号整体偏早（折叠调用点写 `:540`＝实 **`:537`**；`settleCoverageMetric` 写 `:549-552`＝实 **`:544`**；
两个构造函数写 `:568-600`/`:602-612`＝实 **`:564-583`**/**`:589-597`**；`Verdicts` 字段写 `:455-466`＝实 **`:456-463`**），
而 §1.3 那张三处表与 §4 变异表的坐标（`:556`/`:537`/`:529`/`:565`/`:570`）**逐枚正确**；
② §1.3 说"批准之后会红哪**三**枚"与 §4.1/§4.2 说"代价 **2** 枚"自相矛盾，**盘上值是 2**（§2 的 (A)）；
③ §3 说"新码与新测试文件的注释逐字 ASCII"——`sampler.go` 新增段确实 0 枚非 ASCII，但
`sampler_settle_gate_136_test.go:152` 含中文词 `同形`（非 emoji 段、不破 d22scan，见 §5）。

**本程没做、也不冒充做过的**：未复现 §1.2 那六发真取样（档三，见 6.1）；未跑容器/linux 分母
（`scripts/portable-tests.sh` 的真实容器读数实现件自己也没取，账在它 §6-2 与本程这里各记一笔）；
未裁 AC#10／AC#15／`settleCoverageRowGates` 该不该批；未做票面 §3「每片完成后五件事」里的缺口审计与失败预演
（那两件的授权方另有人，本程只是 `AC#14` 一格的裁决者）。

---

## §8 同形那一维的补量 ＋ 收尾 commit 账

**"与 `sampling` 同形"本程逐字比过两枚格式串**（`internal/observe/sampler.go`，量于 `aef82f5`＝HEAD 同码）：

| 侧 | 站点 | 逐字 |
| --- | --- | --- |
| `StateReport`（冻结侧，票面 `:278`） | `sampler.go:336` | `Metric: "sampling", Measured: fmt.Sprintf("%d valid / %d errors", len(rep.Samples), rep.SampleErrors),` |
| `SettleReport`（本格新造） | `sampler.go:577` | `Measured: fmt.Sprintf("%d valid / %d errors", len(rep.Samples), rep.SampleErrors),` |

⇒ **行名与 `Measured` 格式串逐字同**，且 `buildSettleVerdicts` 另起一枚（`:564`）、`buildVerdicts` 未被做成两用
（`TestStateReportVerdictBuilderStaysSinglePurpose` 在 §1/§2/§4 九发里全绿）。票面 `:262` 那句"与 `sampling` 同形的门行"
这一维本程**判成立**（不是"只加了个字段"）。

**本节（§6-§7）的 commit 回显**：

```
$ git log --oneline -1
87f1e86 accept(136 AC#14 r1 §6-§7): 五条判据总裁（①-⑤ 全成立·独立复现，本格附条件在编排者手里）＋ AC#15 前提腿 7 站点现量 ＋ 纪律回执 ＋ 推翻简报四句
$ git show --name-only HEAD
docs/evidence/s1/136-ac14-r1-acceptance.md
```

**本程 commit 账（逐枚只带 `docs/evidence/s1/136-ac14-r1-acceptance.md` 一枚路径；`git log --oneline aef82f5..HEAD -- docs/evidence/s1/136-ac14-r1-acceptance.md` 可复算）**：
`ca09955` §0-§1 → `bb723f7` §2 → `ffe5731` §3 → `24ab73f` §4 → `512d33d` §5 → `87f1e86` §6-§7 →（本节那一枚在最后，hash 由下一位从 `git log` 读）。
中间穿插的 `ddbd3a1`／`6451625`／`98665ef` 等**都不是本程的 commit**（编排者与票 140 程在飞）；本程七枚（含本节）**枚枚单路径、未 push**。

**收尾回执**（`git status --porcelain` 末次）：`design/**` 16 枚删除 ＋ `design/doubao/`、`design/old/` 两枚未跟踪目录
**仍在原状**（未还原、未提交、未删）；`internal/observe/**` 与其余禁改面 **0 行改动**。
