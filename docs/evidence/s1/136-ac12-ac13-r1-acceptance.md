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
