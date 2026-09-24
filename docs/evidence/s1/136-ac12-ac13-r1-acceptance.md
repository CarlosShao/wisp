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
