# 136 — AC#12／AC#13 两格终裁（r1）：非实现方对抗验收

- 裁决方：`acceptor-ticket136-ac12-ac13-r1`（非实现方；本程**未写任何生产码或测试码**、未装任何工具、未翻任何勾）
- 被验版本：**commit `7d73b5f`**（`git cat-file -t` 现核＝commit）
- 本程只裁 **AC#12**（票面判据①②③④⑤）与 **AC#13**（票面判据①②③）两格；AC#1／AC#8／AC#9／AC#10／AC#11 已有别的裁决，本程**一字未重判**
- 三档证据标注口径：〔独立复现〕＝本程在这棵仓外纯净树上自己跑出来的；〔日志＋归档，我抽验〕＝引用他人读数但本程抽验过；〔仅自述，不背书〕＝实现方自己说的、本程未复算
- 本文件出现的每一枚 sha 在引用前都跑了 `git cat-file -t`（核法见 §4）

---

## §0 锚点、取件方式与快照账（只建不删）

| 项 | 实测 |
| --- | --- |
| 取件命令 | `mkdir -p /d/tmp/wisp136r1-tree && git archive 7d73b5f | tar -x -C /d/tmp/wisp136r1-tree` |
| 被验树内容核 | `git diff --stat c03aee3..7d73b5f -- internal/observe/` **无输出** ⇒ 锚点 `7d73b5f` 的 `internal/observe/**` 与 AC#12 交付那枚 commit **逐字节同版**（本票两格的全部读数都取在这一棵上） |
| 改前面（AC#13 判据①） | `/d/tmp/wisp136r1-tree-pre` ＝ `git archive 45c8d1c`（实现方自量的锚点，`git cat-file -t`＝commit；`f4c7062` 之前＝**无 seam 守卫**） |
| pristine | `/d/tmp/wisp136r1-pristine-observe/` ＝ 被验树 `internal/observe/*.go` **23 枚**拷贝，每发变异前无条件从此还原 |
| 驱动器 | `/d/tmp/wisp136r1-acceptor-driver.py`（自带唯一性检查：锚点文本命中数≠1 就拒绝落刀并退出非零；本程**未触发一次拒绝**） |
| 探针原件 | `/d/tmp/wisp136r1-probe-src/zz_acceptor_r1_probe_136_test.go`（本程自己写的那枚最小正常采样用例；实际落点由驱动器插在 `sampler_test.go` 里 `TestSLOStateNamesPinnedToMachineStates` 之前，见 §1.1 的位置说明） |
| 日志 | `/d/tmp/wisp136r1-tree-baseline-v.txt`、`…-pre-baseline-v.txt`、`…-pre-probe-v.txt`、`…-tree-probe-v.txt`、`…-tree-restored-v.txt`、`…-tree-MA-v.txt`、`…-tree-AB-v.txt` |
| 名册 | `/d/tmp/wisp136r1-baseline-names.txt`(65)、`…-pre-names.txt`(58)、`…-pre-probe-names.txt`、`…-tree-probe-names.txt`、`…-tree-restored-names.txt`、`…-MA-names.txt`、`…-AB-names.txt`；差集件 `…-pre-swallowed.txt`、`…-MA-swallowed.txt` |
| 仓内 | **未建 worktree、未 checkout/switch/stash/reset/amend/rebase/clean、未在仓库目录内跑过一次 `go build`／`go test`**（所有 go 命令的 cwd 都是 `D:\tmp\wisp136r1-*`） |
| 工具链 | `go version go1.27.1 windows/amd64`（宿主既有工具，本程未安装任何东西） |
| 四数口径 | 一律 `go test -count=1 -v ./internal/observe/`：`RUN=grep -c '^=== RUN'`、`PASS='^--- PASS'`、`FAIL='^--- FAIL'`、`SKIP='^--- SKIP'`，另给 `^panic` 命中数（**表头＋栈帧各一行 ⇒ 真 panic 枚数＝命中数÷2**）与逐名 `=== RUN` 名册差集；本程未跑过 `-count=2`，也未跑过任何 `-run` 定点判据读数 |
| 墙钟 | 本程时刻以 `date "+%Y-%m-%d %H:%M:%S %z"` 现取原文照抄（§1 落笔＝`2026-09-24 09:44 +0800` 前后各发）；**没有任何一个结论用时间戳相减算出来** |

**基线（被验树 `7d73b5f`，未加任何东西）**：`rc=0 / RUN=65 / PASS=65 / FAIL=0 / SKIP=0 / panic=0`〔独立复现〕，与实现方 §2.5 的 65/65 同数。
**改前基线（`45c8d1c`）**：`rc=0 / RUN=58 / PASS=58 / FAIL=0 / SKIP=0 / panic=0`〔独立复现〕。

---

## §1 AC#13（`R-136-9`，靶＝`sampler_test.go:31` 的 `f.mu[len(f.mu)-1]`）

票面判据：① 自建一枚最小正常采样用例 ⇒ 拍到 `index out of range [-1]`，并逐名列出被吞掉的读数（要给"改前名册 vs 改后名册"的差集）；② 修法只许动那枚仪器（补长度/空集守卫，要红不许静默），不许改 `sampler.go` 的语义去绕它，也不许 `t.Skip`；③ 复跑整包 ⇒ panic=0、名册与基线逐名相同。

### 1.1 判据①：复现成立，被吞读数逐名（改前名册 vs 改后名册的差集）

**我做了什么**：在 `tree-pre`（`git archive 45c8d1c`，无守卫）取基线 58 枚名册；用驱动器把**我自己写的**那枚最小正常采样用例插在 `sampler_test.go` 里 seam 之后、既有采样用例族之前（＝"下一个加测试的人"会放的位置）；先证落地再读数：`grep -n "func TestAcceptorR1ProbeBareFakeTreeSampling"` ＝ `43:`、`go build ./...` rc=0、`go vet ./internal/observe/` rc=0，然后才 `go test -count=1 -v`。〔独立复现〕

**读数**（`/d/tmp/wisp136r1-pre-probe-v.txt`）：`rc=1 / RUN=42 / PASS=40 / FAIL=2 / SKIP=0 / ^panic=2 行（真 panic 一枚）`。panic 原文逐字：

```
--- FAIL: TestAcceptorR1ProbeBareFakeTreeSampling (0.00s)
panic: runtime error: index out of range [-1] [recovered, repanicked]
	github.com/CarlosShao/wisp/internal/observe.(*fakeTree).ReadTree(...)
	D:/tmp/wisp136r1-tree-pre/internal/observe/sampler_test.go:31 +0x211
	github.com/CarlosShao/wisp/internal/observe.(*Sampler).SampleState(...)
	D:/tmp/wisp136r1-tree-pre/internal/observe/sampler.go:269 +0x2b8
```

⇒ 落点逐字对上 `R-136-9` 当初的判词（`sampler_test.go:31`，入口 `sampler.go:269` 首次基线读树）。

**差集（被吞＝基线名册 58 － 该发跑到的基线名 41＝17 枚，逐名按字母序，`/d/tmp/wisp136r1-pre-swallowed.txt`）**：
`TestCheckSettleNeverReachesCap`、`TestCheckSettleVerifiesReleaseCounter`、`TestLiveRegistryBaselineWithinSleepingGate`、`TestMarkTransitionTimestamps`、`TestSLOStateNamesPinnedToMachineStates`、`TestSampleStateAllMetricsAndVerdicts`、`TestSampleStateCPUTotalDrivenMean`、`TestSampleStateHandleGateUsesRulingLimit`、`TestSampleStateLeakFixtureFlipsRed`、`TestSampleStateSleepingDiskWriteGateFails`、`TestSampleStateSleepingTCPGateFails`、`TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（**AC#1 正向对照腿**）、`TestSampleStateUnknownStateAndReaderError`、`TestSampleStateWorkPeakMemoryIsTargetNotGate`、`TestSampleStateZeroSampleWindowFailsClosed`（**AC#1 的钉**）、`TestSamplerGoroutineAccountingFollowsRegistry`（**AC#8 刚修的那枚**）、`TestThresholdTableCoversAllStates`。
⇒ 这 17 枚**既不算红也不算绿**，包级只留一句 `FAIL github.com/CarlosShao/wisp/internal/observe`。

**与实现方读数的对表（差一枚，成因我核清了）**：实现方 §1.1 报 `RUN=43／被吞 16 枚`，我报 `RUN=42／被吞 17 枚`。差的那一枚正是 `TestSLOStateNamesPinnedToMachineStates`——我把探针插在它**之前**，实现方插在它**之后**（它那 16 枚名册里唯独没有这一枚）。⇒ **两棵树指的是同一处雷，数目差＝插入位置差一枚**，与实现方 §1.1 自己如实写的那句"危害大小由谁排在它后面决定"同形；不冲突、不互否。〔独立复现＋抽验它的日志口径〕
**该发另一枚红名**＝`TestNoopTaskReturnsToBaseline`（`goroutine_test.go:33: PerTask mid-task = 2, want 3`）＝票面 **AC#11** 那枚既有 flake，本程只登记、**不记进 AC#13 的账，也不拿 AC#11 抵 AC#13**（详见 §1.5）。

### 1.2 判据②：修法只动了那枚仪器，要红不静默，无 t.Skip

| 我要核的 | 我自己量的 | 结论 |
| --- | --- | --- |
| 只动仪器？ | `git show --numstat f4c7062` ＝ `111 0 internal/observe/sampler_faketree_guard_136_test.go` ＋ `22 0 internal/observe/sampler_test.go`，**无第三枚文件** | 界内 |
| 有没有改 `sampler.go` 的语义去绕它？ | `git diff --stat 45c8d1c..f4c7062 -- internal/observe/sampler.go` **无输出** | 未改被测面 |
| 守卫形状 | 被验树 `sampler_test.go:48-54`：`if f.i >= len(f.mu) {` → `if len(f.mu)==0 {` 返回 `TreeMetrics{}, errFakeTreeNoScript`，否则**逐字保留**旧的"重复末读"行为（`:53`）。旧行为未被写宽：控制腿 `TestFakeTreeScriptedReadsStillRepeatLastReading` 钉住它，我在 §1.4 的 AB 发里证明它会响 | 界内 |
| 要红不许静默？ | 走 `TreeReader` 自己的错误通道（`errors.New` 的 `errFakeTreeNoScript`，`sampler_test.go:38-39`）；不返回零值、不静默 return。§1.4 的 AB 发（把拒绝换成静默零值）**实测两腿红** ⇒ 这一半有牙 | 成立 |
| 有没有 `t.Skip`？ | `grep -rn "t\.Skip\|\.Skip(" internal/observe/*_test.go` 在被验树上＝**2 处命中，都在注释里**（`sampler_test.go:36`、`:335` 的"never t.Skip"字样）；两枚新 136 文件里 `grep -c t.Skip` ＝ **0／0**。九发读数 `--- SKIP` **逐发 0** | 零 Skip |
| 同族普查 | `grep -rn '\[len(.*)-1\]' internal/observe/*_test.go`（被验树）＝**4 行**：两行是我自己看的注释（`sampler_test.go:26`、`sampler_faketree_guard_136_test.go:16`）、`sampler_test.go:53` 是被守卫护住的那支、唯一另一枚真读点 `earlylog_130_test.go:175` 前面 `:169-171` 有 `if len(recs) != earlyLogMaxRecords+1 { t.Fatalf }` 的**硬前置长度守卫** ⇒ 取不到 `-1`。实现方 §5-1 那句"判不必改"我**读码复核成立**，但**我同样没为它落变异真打** | 与它一致 |

### 1.3 判据③：修好的树上复跑整包 ⇒ panic=0、名册完整

**我做了什么**：在同一棵被验树（`7d73b5f`）上，用驱动器把 §1.1 **同一份**探针插在**同一相对位置**（`TestSLOStateNamesPinnedToMachineStates` 之前），先 `go build ./...` rc=0 ＋ `go vet ./internal/observe/` rc=0，再 `go test -count=1 -v`。〔独立复现〕

| 发 | 四数（`-count=1 -v`） | panic | 名册账 |
| --- | --- | --- | --- |
| 被验树基线（无探针） | `rc=0 / 65 / 65 / 0 / SKIP0` | 0 | 名册 65（`wisp136r1-baseline-names.txt`） |
| 被验树＋探针 | `rc=1 / **RUN=66** / PASS=65 / **FAIL=1** / SKIP0` | **0** | 基线 65 枚**缺 0 枚**；多出的恰好只有探针自己一枚（`comm -13` 输出＝`TestAcceptorR1ProbeBareFakeTreeSampling`）⇒ 改前缺 17、改后缺 **0** |
| 还原（`cp` pristine 后） | `rc=0 / 65 / 65 / 0 / SKIP0` | 0 | `diff` 还原名册 vs 基线名册 **无输出＝逐名相同**；`23/23` 枚 `.go` 与 pristine `diff -q` 无输出 |

探针在修好的树上那一枚红的原文（`/d/tmp/wisp136r1-tree-probe-v.txt:99`）：

```
    sampler_test.go:69: bare fakeTree sampling errored: resource: observe: baseline tree read: observe test seam: fakeTree has no scripted reads (bare &fakeTree{} is not a measurement; script mu/current or set err) rep=<nil>
```

⇒ 误用的读数**由它自己那枚用例报出来**（红、带 seam 那句话、`rep=<nil>`），**别人不再被代答**。判据③成立。

### 1.4 那三腿不是哑的：我自己重打两发变异（MA／AB），各响各腿

| 发 | 改法 | 落地证明（先证再读数） | 整包四数 | 红名逐名＋红点（我量到的） |
| --- | --- | --- | --- | --- |
| **MA** | 摘掉守卫（回到旧那一支） | `grep -n "f.mu[len(f.mu)-1]"` ＝ `49:`（守卫已不在）；`go build ./...` rc=0；`go vet ./internal/observe/` rc=0 | `rc=1 / RUN=42 / PASS=40 / FAIL=2 / SKIP0 / ^panic=2 行（真 panic 一枚）` | ① `TestFakeTreeEmptyScriptFailsClosedAndNotPanics` 红在 `sampler_faketree_guard_136_test.go:46`（消息逐字：`fakeTree.ReadTree panicked instead of failing closed: runtime error: index out of range [-1] …`）；② `TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel` FAIL（panic 落在 `:101`，即它自己那发 `SampleState` 里）；控制腿 `TestFakeTreeScriptedReadsStillRepeatLastReading` 同发 **仍 PASS** ⇒ **panic 复发**，基线 65 枚名册**被吞 23 枚**（逐名见 `/d/tmp/wisp136r1-MA-swallowed.txt`，含 AC#1 两腿、AC#9 那枚、AC#12 四腿） |
| **AB** | 守卫在，但把拒绝换成**静默零值** `return TreeMetrics{}, nil` | `grep -n "MUTATION AB"` ＝ `51:`；`go build` rc=0；`go vet` rc=0 | `rc=1 / RUN=65 / PASS=63 / FAIL=2 / SKIP0 / panic=0`，**名册缺 0 枚** | ① `TestFakeTreeEmptyScriptFailsClosedAndNotPanics` 红在 `:57`（`bare &fakeTree{} returned a reading it never measured (metrics=…全部 0…): the seam must fail closed`）；② `TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel` 红在 `:103`（把那份"看着正常"的报告整枚打出来，`Samples:[] SampleErrors:4 … Verdicts:[… {Metric:sampling Measured:0 valid / 4 errors Pass:false Gate:true}] Pass:false`）；控制腿仍绿 ⇒ "不许静默"这一半**有牙** |

MA／AB 两发读数的红名、红点行号与实现方 §1.4 报的**逐枚同点**（`:46`／`:57`／`:103`）。差异只有两处，都不改判定：① 被吞枚数我 23、它 19——我锚点比它多 AC#12 那 4 枚（它们排在 `sampler_test.go` 之前的同名族里，会被同一个 panic 吞掉）；② 我 MA 那发 `FAIL=2` 与它同数。〔独立复现〕

### 1.5 与 AC#11 分开（明令核的一句）

AC#13 的每一格判定都判在**红名是 AC#13 自己仪器**的读数上（§1.3／§1.4 全部 `SKIP=0`、MA 之外 `panic=0`）；AC#11 那枚 flake 本程在改前探针那一发**命中过一枚**（`goroutine_test.go:33`，§1.1 已点名），**未修、未 Skip、未调阈值、也没拿它抵 AC#13 的账，更没拿 AC#13 的绿去抵 AC#11**。`git diff --name-only 45c8d1c..c03aee3 -- internal/observe/goroutine_test.go` **无输出** ⇒ 两格互不侵占。

### 1.6 AC#13 结论：**成立（PASS，无附条件）**

三判据逐条都有本程自己在纯净树上跑出来的读数；"一挂吞一片"在修好的树上**造不出来**（同一发探针：改前吞 17 枚、改后吞 0 枚且 `panic=0`），而摘掉守卫就能**复现**（MA：panic 复发、吞 23 枚）。

**本程没为 AC#13 核的**：`earlylog_130_test.go:175` 那一枚"判不必改"我只做了读码复核（它的前置守卫真在 `:169-171`），**没为它落变异真打**；容器／linux 面、`-race`、`-shuffle`、CI run id 未取；`cmd/wisp` 与本格无关未跑。

---

## §2 AC#12（`R-136-7`，靶＝settle 侧"这一窗丢了几次读数"在报告里查不到）

票面判据：① 报告结构体带上"可信样本数／失败次数／最后一次失败原因"且有 json 标签；② 上一轮终裁方那两发探针各钉一枚用例（"只 1 枚可信"与"一半失败"都必须说不出 pass 的绝对性）；③ 落一发变异（把计数摘掉）⇒ 用例转红、红名逐名；④ 还原复绿、三态齐；⑤ 若需要改 `slo-check.ps1` 或任何 golden 才过 ⇒ 停手上报。
另裁编排者点名的四问：(甲) `dropped_reads` 该不该补／rides 成不成立、(乙) `sampler.go` 改动有没有越出 `:431-443` ∪ `:470-501`、(丙) 消费面 11→13 有没有人读不到东西、(丁) 复跑实现方自加的 MC 探针。

### 2.1 判据①：字段与 json 标签（读的是被验树原文的标签，不是注释）

被验树 `internal/observe/sampler.go:441-455`（该树的 `internal/observe/**` 已证与交付枚 `c03aee3` 逐字节同版）：

| 交付物 | 坐标 | 出线标签（原文照抄） |
| --- | --- | --- |
| 可信样本承载体 | `:441` | Samples []Sample → `json:"samples"` |
| 失败次数 | `:451` | SampleErrors int → `json:"sample_errors"`（**无 omitempty ⇒ 新码写的报告必带这一枚**） |
| 最后一次失败原因 | `:454` | LastSampleError string → `json:"last_sample_error,omitempty"` |

计数落点 `:489-497`（两支丢弃路各自 `rep.SampleErrors++` ＋ 一句原因）。三样都能从**同一份 JSON** 里读到，且与同包 `StateReport:163-168` 的标签**逐字同形**（票面判据①自己引的就是那一族）。〔独立复现（读盘＋探针）〕

**我自己的探针读数**（诊断件、不是判据读数，用 `-run` 定点取；seam 是我自己写的 `acceptorR1Tree`，不依赖实现方任何 fixture；文件 `aaa_acceptor_r1_wire_probe_136_test.go` 留在 `-tree` 那棵快照里，只建不删）：

```
PROBE only-first-read-trustworthy seamReads=10 samples(len)=1 sample_errors=9 last_sample_error=read: acceptor r1: tree read failure pass=true keys=13 [...]
PROBE every-other-read-fails       seamReads=10 samples(len)=5 sample_errors=5 last_sample_error=read: acceptor r1: tree read failure pass=true keys=13 [...]
PROBE zero-value wire: {"target_state":"","started_at":"","peak_bytes":0,"cap_bytes":0,"free_os_memory_count":0,
  "free_os_memory_requested":false,"back_within_cap_ms":0,"elapsed_ms":0,"final_bytes":0,"samples":null,"sample_errors":0,"pass":false}
```

⇒ 两形都满足 "samples 枚数 ＋ sample_errors ＝ 取过的读数"（1+9=10、5+5=10），`last_sample_error` 带的是**我那次失败自己的**文字（不是糊出来的常量）。〔独立复现〕

### 2.2 (甲) §5-5 那枚没补的 key 到底该不该补 ⇒ **rides 成立，判据①按两枚 key 结；`dropped_reads` 不缺料**

**我不读注释、只查"谁往 `Samples` 里写东西"**〔独立复现〕：

- `CheckSettle` 全函数（`:462-522`）里 `rep.Samples` 只出现 **1 处写入**＝`:505`，而它站在 `:498` 的 `} else {` 里；进 `else` 的前置条件是 `:490` 的 `if err != nil` 不成立 **且** `:493` 的 `else if m.PrivateWorkingSetBytes <= 0` 不成立 ⇒ **落进 `samples` 的每一枚都是可信读数**，数组长度即"可信样本数"。
- `SampleState` 那侧（`:247-348`）同形：`:297` 是唯一写入点，站在 `:295` 的 `} else {` 里 ⇒ rides 不是 settle 侧新造的说法，是**这一族一直如此的形状**（`StateReport` 同样没有独立 count 字段，计数在 `sampling` 门行的 "%d valid / %d errors" 里）。
- 全仓没有别的写入者：`grep -rn "\.Samples" cmd/ internal/ --include=*.go` 在 `internal/observe/` 之外 **0 命中**；`grep -rn "SampleErrors" cmd/ internal/ tools/ --include=*.go` 同样 **0 命中**。

**会不会在某些路径下把不可信的也计进去？会，而且只有那一条路径——把 `:493` 那枚 fail-closed 判据放宽。** 这一条我不是推演，是**打出来了**：本程自己落的 MC 发（`<= 0` ⇒ `< 0`）里，报告被打成 `Samples:[…10 枚 TreePrivateBytes:0…] SampleErrors:0 LastSampleError: Pass:true`（原文 `/d/tmp/wisp136r1-tree-MC-v.txt:91`）⇒ "rides 在 `samples` 上"成立与否**正好等于那枚 fail-closed 在不在位**；而**它一被放宽，隔壁 AC#9 的钉立刻红**（同一发 `sampler_settle_zerosample_136_test.go:66`，见 §2.5）。
⇒ **裁定：rides 成立，判据①按两枚 key 结。** 票面判据①要的是三样**信息**（可信数／次数／原因），不是三枚 key 字面量；第三枚 `dropped_reads` 在上一轮终裁方那里只是"探针顺手试的候选名"——`R-136-7` 原文给的修法方向写的就是"给 `SettleReport` 补 `sample_errors`/`last_sample_error` 两枚字段"，同包 `StateReport` 也没有它。实现方把这档明写在 §5-5 请终裁判、没自己扩，处置正确。

**但登记两条边界（不粉饰，也都不是"附条件"）**：
① `sample_errors` 与 `samples` 是"两个数拼一个分母"，任何一支被单独摘掉才看得见丢读数——本程 MD／ME 两发已证**两支丢弃路各自有一枚会红的腿**，所以这条可核；
② `cmd/wisp/slo_windows.go:623` 那枚合成失败报告 `&observe.SettleReport{TargetState: …, Pass: false}` 的 `sample_errors` ＝ 0，语义是**"根本没测"**而不是"零丢失"，出线与"测满、零丢"同形（我的 zero-value 探针原文可对照：`"samples":null,"sample_errors":0`）。它靠 `pass:false` 挡住静默绿，落点在 `cmd/wisp`（AC#10／票 133 地界）⇒ **不算 AC#12 的债，但 AC#10 结案时必须被问到**（本程已在 §总判 写给编排者）。

### 2.3 (乙) `sampler.go` 的改动有没有越出 `:431-443` ∪ `:470-501` ⇒ **界内，无越权 hunk**

命令（在被验仓里跑的只读 git）：

```
git diff -U0 5c1529a^..c03aee3 -- internal/observe/sampler.go     → 旧侧 hunk 头共 3 枚
git diff --stat 45c8d1c..5c1529a^ -- internal/observe/sampler.go  → 无输出（旧侧行号可直接对票面坐标）
```

| hunk（旧侧原文） | 旧侧被碰的行 | 属于哪一段 | 判定 |
| --- | --- | --- | --- |
| `@@ -442 +442,14 @@ type SettleReport struct {` | `:442`（旧 `Pass bool json:"pass"` 一行换成 14 行） | `:431-443`＝`SettleReport` 声明块（旧文件 `:431` 是 `type SettleReport struct {`、`:443` 是 `}`；我逐行点过） | **界内** |
| `@@ -477 +490,4 @@ func (s *Sampler) CheckSettle(…` | `:477`（旧 `if err == nil && m.PrivateWorkingSetBytes > 0 {` 拆成三分支） | `:470-501` | **界内** |
| `@@ -479,0 +496,3 @@ func (s *Sampler) CheckSettle(…` | `:479` 之后纯插入 3 行 | `:470-501` | **界内** |

- `SampleState` 那侧（`:247-348`）：**零 hunk**（-U0 旧侧只有上面 3 枚）。
- 阈值面／golden／`scripts/`／`.github/`／`tools/`／`frontend/`／`docs/SLO.md`：`git diff --stat 45c8d1c..7d73b5f` 限定这些路径 **无输出**（整段区间连兄弟都没碰过）。
- 本程两格的全部净面（`git show --numstat` 逐枚自量，不抄派单）：`f4c7062`＝2 枚测试文件、`sampler.go` 0 命中；`5c1529a`＝`sampler.go 21 2`；`f53ad5c`＝新测试 `289 0`；`c03aee3`＝同文件 `7 2` ⇒ **四枚 commit 合计只碰 3 枚文件，全在 `internal/observe/` 里**。
- 语义等价我也读了：`:490-498` 三分支重写后，"落进 `Samples`／`FinalBytes`／`BackWithinCapMS`"的条件仍是 `err==nil && PrivateWorkingSetBytes>0` ⇒ 判据①要的"加自陈"没顺手改掉判据本身。〔独立复现〕

⇒ **这一条按编排者 `9675333` 的更正不再审罪，只复核边界：越界不存在，hunk 逐枚界内。**

### 2.4 (丙) 消费面 11→13，有没有人因此读不到东西 ⇒ **只多不变少；但两枚新 key 今天零真实消费者**

我自己现算（在被验树上 `grep -rn`／`git grep`，不采信实现方 §3 的分母）：

| 消费者 | 用法 | 会不会读不到东西 |
| --- | --- | --- |
| `cmd/wisp/slo_windows.go:118` `Settle *observe.SettleReport json:"settle,omitempty"` | 生产侧序列化（`writeSLO` 走 `json.MarshalIndent`） | 否——只**增** key |
| `cmd/wisp/slo_windows.go:286` `run.Pass = run.Settle.Pass` | 只读 `.Pass` | 否——该字段未动 |
| `cmd/wisp/slo_windows.go:599` `runSettle(…) *observe.SettleReport`、`:623` 合成失败报告 | 生产侧构造 | 否（但见 §2.2 边界②） |
| `cmd/wisp/slo_windows.go:528` `json.Unmarshal(data, &rep)`（`collectReport` 读回**同一次 run、同一版二进制自己写的**文件） | 同类型往返 | 否——`encoding/json` 默认忽略未知 key；全仓 `DisallowUnknownFields` 只有 `internal/config/parse.go:72` 与 `:123`，**这条链上没有严格解析器** |
| `scripts/slo-check.ps1:350-355` | `ConvertFrom-Json` 后只取 `.pass` 与 `.settle.free_os_memory_count` | 否——两枚都是原 key；脚本里**没有任何 key 集合／字段数断言**（`:341-361` 原文我逐行读过） |
| 面板桥（`internal/panel/**`） | `grep -rn "Settle" internal/panel/*.go` **0 命中**；frontend 那处是英文散文 | 不适用 |
| 归档报告 `docs/evidence/s1/66/*.json` | `grep -rn "66-settle-1" --include=*.go --include=*.ps1 --include=*.sh --include=*.yml` 与同形 "66-full-subset" ⇒ **0 命中**：没有任何仪器拿旧 settle 报告当期望值 | 不适用 |

- **key 全清点（我自己的探针，`-run` 定点）**：有丢弃时出线 **13 枚** ＝ 旧 11 枚（`back_within_cap_ms`、`cap_bytes`、`elapsed_ms`、`final_bytes`、`free_os_memory_count`、`free_os_memory_requested`、`pass`、`peak_bytes`、`samples`、`started_at`、`target_state`）＋ `sample_errors` ＋ `last_sample_error`；**无丢弃时 12 枚**（`last_sample_error` 带 omitempty）。⇒ 实现方"11→13"这句**在有丢弃的形状下逐字成立**；它没写"无丢弃是 12"，属措辞偏窄，不是错。
- **两枚新 key 的真实消费者＝ 0 枚生产读者**。`SampleErrors`／`LastSampleError` 在 `internal/observe/` 之外没有任何 Go 侧引用，唯一读者是本包那四腿（`sampler_settle_coverage_136_test.go` 里 15 处断言）。同族旧账一起记：`StateReport` 那两枚 key 自 ticket 66 起也是**零程序读者**（只有 `docs/SLO.md:428` 那栏**人读**的表把 `sample_errors` 当列名）。
- **算不算 AC#12 的债？我的裁定：不算，但要记一笔。** 本格票面那句目标是"让'这一窗丢了几次读数'**能从报告里读出来**"，报告本身就是交付物（出线、进 artifact），此谓已达成；**没达成的是"生产侧真报告里 `sample_errors` 到底是不是非零"**——那一发要真跑 `wisp slo -settle`，是**已开的 AC#10**（且被 `cmd/wisp` 在飞挡着），不是本格藏了一半。本程把 `cmd/wisp` 的**测试二进制按新 `observe` 重新链接**过一次：`go test -c` rc=0（32,580,515 字节），**未执行**（快照树缺 sherpa DLL，且 `cmd/wisp` 此刻归票 133/135 在跑）⇒ 只证"唯一 Go 侧消费者编得过、链接得上"，不替它说"跑起来也绿"。〔独立复现〕

### 2.5 判据②③④：四腿钉、四发变异（红名逐名）、三态

**判据② 两发探针各一枚**〔独立复现（读盘）〕：`TestCheckSettleSingleTrustworthyReadReportsItsLoss`（fixture＝`{trustworthy, failed, failed}` 且脚本耗尽后恒答 failed ⇒ 十次取读只 1 枚可信，正是探针①形状）、`TestCheckSettleHalfTheReadsFailedReportsItsLoss`（`alternatingTree` 严格按读序号交替，正是探针②形状）。两腿都把"说不出 pass 的绝对性"做成断言：`:134`（`sample_errors == reads-1`，期望值来自 **seam 自己的计数器**、不是被测函数的返回值）、`:138`（`>=2`）、`:141`／`:206`（原因句必须含本程 fixture 那句 distinctive 文字）、`:155-169`／`:213-217`（**出线**同数、kept＋lost＝reads）。另两腿是加料不是替代：`:246`（另一支丢弃路也要被数到）＋ `TestCheckSettleFullyMeasuredWindowReportsNoLoss`（正向对照，防"逢窗口就报丢"，`:273`/`:277`/`:283`/`:287` 钉住"sample_errors 在且为 0"、"last_sample_error 缺席"）。⇒ **恒真这一味被堵住了**，它不是只在看笼统的 rc。

**判据③④ 我自己重打四发**（全部落在被验树 `internal/observe/sampler.go`；**每发先无条件 `cp` pristine、先证落地再读数**；变异只落仓外树，仓库一字未动）〔独立复现〕：

| 发 | 改法 | 落地证明（先证） | 整包 `-count=1 -v` | 红名逐名＋红点（我量到的） | 与它自述对表 |
| --- | --- | --- | --- | --- | --- |
| **MD**（票面点名那发"把计数摘掉"） | `:491` err 支的 `rep.SampleErrors++` 摘掉（读树仍被丢弃） | `grep -n "MUTATION MD"` ⇒ `491:`；`go build ./...` rc=0；`go vet ./internal/observe/` rc=0 | `rc=1 / RUN=65 / PASS=63 / FAIL=2 / SKIP=0 / panic=0`，名册缺 **0** 枚 | `TestCheckSettleSingleTrustworthyReadReportsItsLoss` `:134`＝`report counted 0 dropped reads but the seam took 10 reads and kept 1: 9 unaccounted`；`TestCheckSettleHalfTheReadsFailedReportsItsLoss` `:197`＝`the seam lost 5 of 10 reads but the report says sample_errors=0: a half-covered window must report its losses` | 同两枚、同两处红点（它那一发另多一枚既有 flake ⇒ 两态不冲突） |
| **ME** | `:496` 零足迹支的计数摘掉 | `grep -n "MUTATION ME"` ⇒ `496:`；build rc=0 | `rc=1 / 65 / 64 / 1 / SKIP=0 / panic=0` | 只一枚：`TestCheckSettleZeroFootprintDropsAreCountedToo` `:247` | 同 |
| **MF** | `:491` 计数留着、**原因**摘掉 | `grep -n "MUTATION MF"` ⇒ `491:`；build rc=0 | `rc=1 / 65 / 63 / 2 / SKIP=0 / panic=0` | `…SingleTrustworthyRead…` `:141`＝`got ""`；`…HalfTheReadsFailed…` `:206`＝`got ""` | 同 |
| **(丁) MC**（它自加那发，我重造） | `:493` 的 `<= 0` ⇒ `< 0`（放宽 settle 足迹判据） | `grep -n "MUTATION MC"` ⇒ `493:`，被改行原文＝`} else if m.PrivateWorkingSetBytes < 0 { // MUTATION MC`；`go build ./...` rc=0 | `rc=1 / 65 / 63 / 2 / SKIP=0 / panic=0` | ① **`TestCheckSettleZeroTrustworthySamplesFailsClosed` 红在 `sampler_settle_zerosample_136_test.go:66`**（AC#9 那枚钉，消息里就是那串病形：5 枚 `TreePrivateBytes:0` 的"样本"、`Pass:true`）；② `TestCheckSettleZeroFootprintDropsAreCountedToo` 红在 `:244`（`want 1 recorded sample, got 10`）；腿 1／腿 2／腿 4 同发仍绿 | 同两枚、同两处红点 |

**(丁) 单独回一句**：我**在这棵纯净树上自己造了一遍**，没引用它的日志。落地三证齐（`grep` 到我改那一行的原文 ⇒ `go build ./...` rc=0 ⇒ 才读数），红点**逐名到用例**、不是"整包红了"。⇒ **AC#12 的改动没把 AC#9 的牙磨钝，成立**；这一发还顺带把 §2.2 的 rides 判据钉死（放宽足迹判据 ⇒ `samples` 里立刻混进 `TreePrivateBytes:0` 的不可信读数，而这正被 `:66` 那枚钉抓到）。

**判据④ 三态齐**〔独立复现〕：未变异 `rc=0 / 65 / 65 / 0 / SKIP 0 / panic 0` → 上表四发各自红 → 还原：`cp` pristine 后 **23/23** 枚 `.go` 对 pristine `diff -q` **全部无输出**，且 pristine 的 `sampler.go` 与 `git cat-file -p 7d73b5f:internal/observe/sampler.go` 逐字节相同（⇒ 我落变异的底面＝被验面）→ 复绿 `rc=0 / 65 / 65 / 0 / SKIP 0 / panic 0`，65 枚名册与基线 `diff` **无输出（逐名相同）**。

**判据⑤ 未触发**：`thresholds.go`／`scripts/`（含 `slo-check.ps1`）／`.github/`／`tools/`／`frontend/`／`docs/SLO.md` 在 `45c8d1c..7d73b5f` 整段区间限定这些路径的 `git diff --stat` **无输出**；本程也**没有**为了过而改任何断言、阈值或 golden。〔独立复现〕

**`c03aee3` 那处"合一断言拆两跳"我单独看了**：`git show --numstat c03aee3` ＝ `7 增 2 删`，只把腿 2 的 "kept 与 lost 合写在一条 precondition broken 里"拆成 `reads<4`／`kept<2`（前提）与 `lost<2`（性质，红点 `:197` 带 `must report its losses`）三跳。⇒ **断言强度只增不减**（MD 复跑红点因此从合写跳移到性质跳），归因变清楚，不是放宽。它把第一次 MD（红点 `:192`）留档而不进结论，处置正确。〔读 diff＝独立复现；它的第一发读数＝仅自述，我不背书〕

### 2.6 AC#12 结论：**成立（PASS，无附条件）**

判据①②③④⑤逐条都有本程自己在被验树上的读数；(甲) rides 成立且由一枚会红的钉兜住、`dropped_reads` 不算缺料；(乙) 界内、无越权 hunk；(丙) 只多不变少，新 key 零生产读者但那笔债记在 AC#10 不在本格；(丁) 独立复现，AC#9 的钉仍红在 `:66`。

**本程没为 AC#12 核的**：容器／linux 面**没真跑过测试**（只跑了 `GOOS=linux go vet ./internal/observe/` rc=0 这一形）；`-race`／`-shuffle` 未跑；`wisp slo -settle` 端到端与 `scripts/slo-check.ps1` 那一腿**零读数**（`cmd/wisp` 的测试二进制我只链接未执行）；CI run id 未取；`-count=2` 本程未复跑（实现方 §2.5 的 130/130 属〔仅自述，不背书〕，它的 §4 门禁表同理）；`gofmt`/`gofumpt` 我未跑二进制（避免在快照里引入工具动作），只目测新字段的对齐与同族一致。

---

## §3 两格之外的账：快照清单、共树暂存事件、两个计数、sha 现核账

### 3.1 我建了哪几棵／哪几件（**只建不删**，交编排者统一清点）

| 类型 | 路径 |
| --- | --- |
| 仓外快照树（2 棵） | `D:\tmp\wisp136r1-tree`（＝`git archive 7d73b5f`，被验版本；⚠ 这棵里现在**留着**我的两枚探针痕迹：`internal/observe/sampler_test.go` 已被 pristine 覆盖回原样、但 `internal/observe/aaa_acceptor_r1_wire_probe_136_test.go` 是我加的独立文件，未删 ⇒ **任何人在该树复跑整包会看到 66 枚而不是 65 枚**）；`D:\tmp\wisp136r1-tree-pre`（＝`git archive 45c8d1c`，AC#13 的改前面；⚠ 这棵的 `sampler_test.go` 里留着我插入的探针，未还原——它的基线读数取在插入**之前**） |
| pristine | `D:\tmp\wisp136r1-pristine-observe`（被验树 `internal/observe/*.go` 23 枚拷贝；已证其中 `sampler.go` 与 `git cat-file -p 7d73b5f:internal/observe/sampler.go` 逐字节相同） |
| 探针／驱动 | `D:\tmp\wisp136r1-probe-src`（两枚探针原件）、`D:\tmp\wisp136r1-acceptor-driver.py`（落地前做唯一性检查，命中数≠1 即拒绝；本程未触发一次拒绝）、`D:\tmp\wisp136r1-compose.py` |
| `-v` 日志（**12 枚**，现数 `ls -1 wisp136r1-*-v.txt`＝12） | `wisp136r1-tree-baseline-v.txt`、`wisp136r1-pre-baseline-v.txt`、`wisp136r1-pre-probe-v.txt`、`wisp136r1-tree-probe-v.txt`、`wisp136r1-tree-restored-v.txt`、`wisp136r1-tree-restored2-v.txt`、`wisp136r1-tree-MA-v.txt`、`wisp136r1-tree-AB-v.txt`、`wisp136r1-tree-MD-v.txt`、`wisp136r1-tree-ME-v.txt`、`wisp136r1-tree-MF-v.txt`、`wisp136r1-tree-MC-v.txt` |
| 名册与差集（**13 枚**＝9 枚 `=== RUN` 名册＋3 枚差集＋1 枚 sha 清单，现数 `ls -1 wisp136r1-*names*.txt wisp136r1-*swallowed.txt wisp136r1-myshas.txt`＝13） | `wisp136r1-baseline-names.txt`、`…-pre-names.txt`、`…-pre-probe-names.txt`、`…-tree-probe-names.txt`、`…-tree-restored-names.txt`、`…-restored2-names.txt`、`…-MA-names.txt`、`…-AB-names.txt`、`…-MD-names.txt`、`…-pre-swallowed.txt`、`…-MA-swallowed.txt`、`…-tree-swallowed.txt`、`…-myshas.txt` |
| 其他 | `wisp136r1-anchor-sampler.go`（7d73b5f 的 blob）、`wisp136r1-cmdwisp.test.exe`（只链接未执行的测试二进制）、`wisp136r1-sec01-tmp.md`／`wisp136r1-sec2-tmp.md`（本文件重建时的两半，见 §3.3） |
| 仓内 | **未建 worktree、未 checkout/switch/stash/reset/amend/rebase/clean**；仓库目录内**一次 `go build`／`go test` 都没跑过**；仓内唯一写入＝本文件 |

### 3.2 共树暂存事件（**要编排者知道的一条**，不是谁的错）

我落 AC#12 那一节 commit 时，共享索引里正**带着兄弟在飞的一枚暂存件** `cmd/wisp/leg_dispatch_gate_133_test.go`（票 135，状态 ` M`→被我 `git diff --cached --name-only` 看见时已是 `M ` 已暂存）。我的两枚 commit 都带显式 pathspec ⇒ 复核结果〔独立复现〕：

```
git show --numstat ae7d593 → 107 0  docs/evidence/s1/136-ac12-ac13-r1-acceptance.md
git show --numstat d3e3a43 → 112 0  docs/evidence/s1/136-ac12-ac13-r1-acceptance.md
```

⇒ **别人的路径没被我的 commit 带走**（两枚各只有我自己那一枚文件、删除列 0）；那枚兄弟的暂存件随后由它自己以 `fa35557`（`274 28 cmd/wisp/leg_dispatch_gate_133_test.go`）提交，`git cat-file -t`＝commit。反向那条（别人不带 pathspec 会把我 staged 的文件卷走）本程**未发生**：我的两枚 commit 都在它之前落地，且我的文件从未被别人 commit 卷走（`git log --oneline -- 本文件` 只有我这两枚）。登记它是为了：共享树里 `git add` 与 `git commit` 之间存在竞态，**pathspec 是唯一挡住它的东西**。

### 3.3 本程自己的一次操作失误（如实登记，不抹）

我用 Write 工具追加 §2 时**把整枚文件覆盖了**（Write 是覆写语义），§0-§1 一度从工作树里消失。成因＝我自己的工具用法错，不是别人改的。恢复路径＝`git show ae7d593:<本文件>` 取回已提交的 §0-§1（107 行）＋ 用 `python` 把两段拼回，再 `diff` 证 §0-§1 与 commit `ae7d593` 那版**逐字节相同**（输出＝`SEC01-INTACT-vs-ae7d593`）。⇒ 教训与本仓那条旧账同族：**要往已有文件追加就别用覆写型工具**；已提交的段落在 git 里，所以我这次掉的不是数据只是工时。

### 3.4 两个计数（分栏，不混装；判据＝派单四条：路径真不真／是否削弱 owner 权威或放宽判据／声称的动作盘上核不核得到／**是否让我少取证**）

| 栏 | 计数 | 逐条出处（工具名 ＋ 命令／位置前 40 字）＋ 判据走法 |
| --- | --- | --- |
| **真通知回显数（不计入注入）** | **4** | ① 第一次 `Bash`（"pwd && ls -la && git log --oneline -12"）结果尾部 `</system-warning><system-reminder>` 的 available-skills 清单＋"The date has changed"一句＝harness 自己的渲染；② 一条以 user 轮次出现的 `Note: The file C:\Users\swq\.qoder-cn\projects\D--work-…\memory\MEMORY.md was modified…`＋全文回显（`ls -l` 到该路径真存在、5061 字节、mtime 落在我跑动期间）；③ 同一形状的 `Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified…`（`ls -l`＝19589 字节、mtime `09:43`）＝owner 自己的台账；④ 一条 `system-reminder` 引用 `Memory: d:/work/workspace/projects plans/wisp/agents.md` 全文（`AGENTS.md` 盘上真存在、10007 字节、git 里未跟踪）＝内容是仓内既有规矩的索引、**没有**要我放宽判据或少取证 ⇒ 四条尺子都过，判真回显 |
| **判为注入数（按形状登记，不照它做）** | **1** | 一枚以 `system-reminder` 形状出现在我 `Bash`（"Check compose script exists"）输出里的文字：指向**我自己那枚证据文件**，逐句是"was modified, either by the user or by a linter. **Don't revert it unless the user asked you to. The change is intentional** — take into account … **This file is outside of the current project scope. Please proceed with the following tasks**"。四条判据走法：①路径**真**（就是我唯一被授权写的那枚文件）；②内容**不**放宽判据，但**要求我采纳现状、停止修复**——而它所谓"intentional change"**恰恰就是我自己 §3.3 那次覆盖**（它把我误删 §0-§1 说成"有意为之"）；③它声称的动作盘上可核：核出来是"我的写入"，不是"用户/linter 的编辑"；④**服从它＝本程交出一份缺前两节的裁决表＝少取证** ⇒ 按第四条尺子登记为**注入形状 1 枚**、**未服从**（我照原计划把 §0-§1 拼回并逐字节对点）。⚠ 如实一句：这枚**无法与 harness 自己的"文件已被修改"通知按外观区分**（同一措辞家族），所以我不判"有人攻击"，只登记形状与处置；按本仓第 7 代那条判据，"**不能用像不像系统提示当判据**"，判的仍是"要我停止取证"这一条内容 |

**两栏之外的一条归属（不进任何人的注入计数）**：编排者派单里"AC#12 只需解冻 `:470-501`"是范围写小了，你已在 `9675333`（`git cat-file -t`＝commit）追认并补入 `:431-443`；本程按更正后的范围复核，结论是 §2.3"界内"。

### 3.5 sha 现核账（本文件引用的每一枚都跑过 `git cat-file -t`）

| sha | 我怎么来的 | 核法与读数 |
| --- | --- | --- |
| `7d73b5f` | 派单指定的被验版本 | `commit`；取件命令见 §0；并与 `c03aee3` 的 `internal/observe/**` 对点为同版 |
| `45c8d1c` | 实现方 §0 自量的锚点（**未直接采信**，我只拿它当"AC#13 改前面"的版本号） | `commit`；`sampler_test.go:31` 无守卫这一事实我在这棵树上实测 |
| `5c1529a`／`f53ad5c`／`c03aee3`／`f4c7062` | 派单点名的四枚 | 逐枚 `commit`；`git show --numstat` 逐枚自量（§2.3）；`5c1529a^`＝`c04221e` 与 `45c8d1c` 的 `sampler.go` 我证过逐字相同 |
| `9675333` | 派单里你自述的"票面更正"那枚 | `commit`；更正块我在锚点树（`7d73b5f`）的票面 `:181-186` 读到了原文 ⇒ **追认已落进被验版本** |
| `ae7d593`／`d3e3a43` | 本程自己两枚 commit（`git log` 现取） | `commit`；净面见 §3.2 |
| `fa35557`／`4e66817` | 兄弟在飞的两枚（`git log` 现取，只用于 §3.2 与 AGENTS.md 生成时刻的归属，**不作为任何判据依据**） | `commit` |

⇒ 本程**没有一枚 sha 是从工具输出／别人的报告里抄来当锚点用的**；也没有出现一枚 `git cat-file -t` 失败的 sha（若有我会登记，本程为零）。

---

## §4 总判

| 格 | 结论 | 翻不翻勾 | 一句理由（全部是本程自己在被验树上的读数） |
| --- | --- | --- | --- |
| **AC#12** | **成立（PASS、无附条件）** | **翻** | 判据①字段＋json 标签在被验树原文逐枚给号，我自己写的 seam 探针量到两形都满足 "samples 枚数＋sample_errors＝取过的读数"；②两发探针各一枚腿、断言的期望值来自 seam 自己的计数器且有正向对照防恒真；③我自己重打四发（MD／ME／MF／MC）**各响各腿、红名逐名到行**；④三态齐（65/65 → 四发红 → `cp` pristine 后 23/23 枚逐字相同且与 `7d73b5f` 的 blob 同版 → 复绿 65/65、名册逐名相同）；⑤**未触发**（阈值面／golden／`scripts/` 整段区间零命中，也没人为过绿改任何断言）。四问：(甲) rides **成立**、`dropped_reads` 不缺料；(乙) **界内**、无越权 hunk；(丙) 只多不变少、新 key 零生产读者但那笔债归 AC#10；(丁) **独立复现**，AC#9 那枚钉仍红在 `:66` |
| **AC#13** | **成立（PASS、无附条件）** | **翻** | ①改前树上我自建的最小正常采样用例拍到逐字 panic 源（`sampler_test.go:31` ← `sampler.go:269`）并给出名册差集（基线 58 → 该发 42，**被吞 17 枚逐名**）；②修法只碰仪器（`sampler.go` 在 `45c8d1c..f4c7062` **零命中**）、走错误通道要红不静默、**全包零 `t.Skip`**；③同一发探针落在修好的树上 ⇒ `RUN=66／FAIL=1（探针自己）／panic=0`、名册**缺 0 枚**，还原后 65/65 且名册逐名相同。三腿不哑：MA 摘守卫⇒panic 复发、再吞 23 枚；AB 换成静默零值⇒两腿红而名册完整 |

**两格互不抵账**（AC#13 不许拿 AC#11 的账抵）：本程唯一一次 flake 命中在改前探针那一发（`TestNoopTaskReturnsToBaseline`，`goroutine_test.go:33`），我**只登记、未修、未 Skip、未调阈值**，也没让它出现在任何一格的绿里；`internal/observe/goroutine_test.go` 在 `45c8d1c..c03aee3` 零命中。

### 需要编排者补／裁的几味（**没有一条是本格的退回项**，但别静默沉掉）

1. **来源账里有一味没做，也没进判据**：`R-136-7` 的修法方向原文是"补两枚字段 **＋ 产出一枚与 `sampling` 同形的自陈门行**"（`136-ac8-ac9-r1-acceptance.md:463`），而票面 AC#12 的**可重算判据①只写了字段＋json 标签** ⇒ 我按票面判据结这一格（成立）。实现方 §5-6 已如实列"没造门行（`SettleReport` 无 `Verdicts` 字段，加它是更大的出线形状改动）"。⇒ **请决定**：要么另立一格／并进 AC#10，要么明写"门行这一维不要"，别让它作为"半句话"留在两本账之间。
2. **`dropped_reads` 的裁定已给**（§2.2：不需要），但若 owner 按票面判据①**字面**要三枚独立 key，那是**改判据**而不是补料——请具名改票面再退回，别由验收方或实现方任一侧解释。
3. **两枚新 key 零生产读者**（§2.4）：`SampleErrors`／`LastSampleError` 在 `internal/observe/` 之外无任何引用，`cmd/wisp` 只序列化不读，`slo-check.ps1` 只取 `.pass` 与 `.settle.free_os_memory_count`。本格判据不要求读者，但 AC#10 结案时应被问到"真报告里 `sample_errors` 到底非零没有"。
4. **`cmd/wisp/slo_windows.go:623` 那枚合成失败报告**的 `sample_errors=0` 语义是"根本没测"、不是"零丢失"（与"测满零丢"出线同形，靠 `pass=false` 挡住静默绿）。落点在别人地界 ⇒ **登记给 AC#10／票 133**，AC#12 不代它负责。
5. **引用"被吞 16 枚"这句要带位置口径**（§1.1）：实现方 16、我 17，差一枚＝探针插入位置差一位。台账／票面今后引这个数请连"插在哪个位置"一起引，否则又是一枚会腐坏的数字。
6. 共树暂存事件一条（§3.2）＋ 我自己那次覆盖事故一条（§3.3）：都是**过程账**，不改两格判定。

### 本程**没核**的部分（明确列出，别以为全核了）

- **未复算**（只标〔仅自述，不背书〕）：实现方 §4 门禁表里的 `gofmt`／`gofumpt`（我未跑这两个二进制）、`-count=2` 的 130/130 与 `122/121/1` 两发、`d22scan` 的 402→403→404、`GOOS=linux go vet ./...` 整树那 3 行诊断的逐错误归因（我只按包作用域跑了 `./internal/observe/` 双 GOOS，都 rc=0）、它 §2.7 留档的第一发 MD 读数、它 §1.5 的 19 枚名册（我另算了自己的 23 枚）。
- **未跑**：容器／linux 真跑测试、`-race`、`-shuffle`、`wisp slo -settle` 端到端、`scripts/slo-check.ps1` 那一腿、`cmd/wisp` 的测试（只 `go test -c` 链接 rc=0 未执行）、CI run id、`scripts/portable-tests.sh --scope=core` 整条。
- **未重判**：AC#1／AC#8／AC#9／AC#10／AC#11 五格（含 AC#9 那枚钉的其它判据——我只用了它当 MC 那发的"牙还在"证据）、票面其它格、上一轮终裁表本身的结论、实现方 §5 其余未做档。
- **未做**：`earlylog_130_test.go:175` 那一枚同族真读点的变异（我读码确认它前面 `:169-171` 有硬长度守卫，但**没真打**，与前两程同一口径）；`R-136-8`（空 `Samples` 出线 `null` 使腿 A 那一支永不响）**不在本格判据内，未动也未裁**。

**本程落笔时刻**（现取原文、非相减）：`date "+%Y-%m-%d %H:%M:%S %z"` ＝ **`2026-09-24 09:59:46 +0800`**（写这一节前后的各发时刻见 §0 与 §2.5 表头；两格各一枚 commit：AC#13＝`ae7d593`（09:44:18）、AC#12＝`d3e3a43`（09:55:44），本节随第三枚）。**未翻任何勾、未 push、未动任何生产码或测试码。**
