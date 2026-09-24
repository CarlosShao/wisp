# 136 — AC#13（第二枚"吞读数"仪器）＋ AC#12（settle 侧自陈"这一窗丢了几次读数"）：实现方交件

- 实现方：`worker-ticket136-ac12-ac13`（本文件两格都是实现者交付；两格终裁与翻勾归非实现者）
- 做序：按派单，**先 AC#13（不依赖生产码），后 AC#12（动生产码）**
- 证据渐进写：§0／§1… 每节一枚 commit；本文件里出现的每个数字都是本程自己量的，**没有一处抄终裁方或前例的读数当结论**（前例只当形状参照）
- 两格都**未翻任何 AC 勾**；票面只追加 log（append-only）

---

## §0 锚点、树、工具链（开工第一步自己量）

| 项 | 实测 |
| --- | --- |
| 开工第一次 `git rev-parse --short HEAD` | **`45c8d1c`** ＝ 本程唯一权威锚点（那一刻 `date -u` 原文 `Wed Sep 23 16:02:32 UTC 2026` ⇒ 本机 **09-24 00:02:32 +08**，+8 手工换算，没塞进 `date` 格式串） |
| 派单给的约值 | "约 `45c8d1c`"——与本程自量值相同；**但本程只认自己 `rev-parse` 那一发**，不认派单转述 |
| 开工 `git status --short` | 只有 1 枚未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`（**未读、未提交、未改、不据它改判据**）；分支 `dev` |
| 外来 sha 是否核过 | **全部现核**。本程引用到的每一枚外来短哈希都跑了 `git cat-file -t`，逐枚 `commit`：`76662d8`（终裁方锚点）/`1d38206`/`79ddd49`/`2f291d0`（前两格实现方锚点与 commit）/`172c7aa`/`c2f6b31`/`6effb7e`/`e8190bf`。⚠ 没有任何一枚 sha 是从工具输出／别人的报告里抄来当锚点用的；本程所有快照只从 `45c8d1c` 与我自己那枚 commit 抽 |
| 兄弟在飞的推进 | 跑动期间 HEAD 从 `45c8d1c` 前进到 `e2ba9f0`→`7ad0ec6`→`94b7267`→`69414d3`（票 135/137 的 `.md`）。`git log --oneline 45c8d1c..HEAD -- internal/observe/` **只有我自己那一枚 commit**；`git diff --stat 45c8d1c..HEAD -- internal/observe/sampler.go` **无输出** ⇒ 被测面在锚点系上没被第三方动过。我的 commit 落在当时的 tip 上（共树正常形态），`git merge-base --is-ancestor 45c8d1c f4c7062` ＝真 |
| 仓外纯净树（**只建不删**，目录名带本会话后缀 `wisp136ac1213-*`） | `/d/tmp/wisp136ac1213-tree0` ＝ `git archive 45c8d1c`（AC#13① 复现面，**无守卫**）；`/d/tmp/wisp136ac1213-tree1` ＝ `git archive f4c7062`（AC#13②③ 与变异面）；`/d/tmp/wisp136ac1213-pristine-tree1/` ＝ tree1 的 `internal/observe/*.go` 22 枚拷贝（每发变异前**无条件 restore** 的源）；AC#12 的 `tree2`/`pristine-tree2` 见 §2 |
| 仓内 | **未建 worktree、未 checkout、未 stash/reset/amend/rebase**；临时件一律只建不删 |
| 工具链 | `go version go1.27.1 windows/amd64`；`gofumpt` 用**既有二进制** `/d/work/base/gopath/bin/gofumpt.exe`，`--version` 原文 **`v0.12.0 (go1.27.1)`**；⚠ **本程没有跑过任何 `go install …@latest`**（派单明令，那会升掉宿主工具） |
| 四数怎么量 | 一律 `go test -count=1 -v`，`RUN=grep -c '^=== RUN'`、`PASS='^--- PASS'`、`FAIL='^--- FAIL'`、`SKIP='^--- SKIP'`，另给 `^panic` 命中数与**逐名 `=== RUN` 名册差集**；`-count=2` 只用于门禁那一发并单独标注 |

---

## §1 AC#13：三判据逐判（靶＝`sampler_test.go:31` 的 `f.mu[len(f.mu)-1]`，同族于 AC#8 但形状不是 `[0]` 直取）

### 1.1 判据①：复现那一发，并逐名列出被它吞掉的同包读数

**基线（tree0＝`45c8d1c`，未加任何东西，`go test -count=1 -v ./internal/observe/`）**：
`rc=0 / RUN=58 / PASS=58 / FAIL=0 / SKIP=0 / ^panic=0`（原文 `/d/tmp/wisp136ac1213-tree0-baseline-v.txt`）。
`internal/observe/sampler.go` 与终裁方锚点 `76662d8` 逐字相同（`git diff --stat 76662d8..45c8d1c -- internal/observe/` 无输出），所以基线 58 枚与"守卫＋AC#9 钉"那版同集合。

**我自己写的最小"正常采样"用例**（落在 `sampler_test.go` 里、其它采样用例之前——就是"下一个加测试的人"会放的位置；用例只走 `SampleState`，不写任何断言，先让仪器自己表态）：

```go
func TestProbeAC13BareFakeTreeNormalSampling(t *testing.T) {
	s := NewSampler(&fakeTree{}, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 40*time.Millisecond)
	if err != nil {
		t.Fatalf("bare fake tree: err=%v rep=%+v", err, rep)
	}
	t.Logf("bare fake tree: %d samples pass=%v", len(rep.Samples), rep.Pass)
}
```

**读数（原文 `/d/tmp/wisp136ac1213-tree0-probe-v.txt`）**：`rc=1 / RUN=43 / PASS=42 / FAIL=1 / SKIP=0`，panic **一枚**（`grep ^panic` 命中 2 行＝`panic:` 表头＋栈里的 `panic({…})` 帧）。红／panic 原文：

```
=== RUN   TestProbeAC13BareFakeTreeNormalSampling
--- FAIL: TestProbeAC13BareFakeTreeNormalSampling (0.00s)
panic: runtime error: index out of range [-1] [recovered, repanicked]
	github.com/CarlosShao/wisp/internal/observe.(*fakeTree).ReadTree(...)
	.../internal/observe/sampler_test.go:31 +0x211
	github.com/CarlosShao/wisp/internal/observe.(*Sampler).SampleState(...)
	.../internal/observe/sampler.go:269 +0x2b8
```

⇒ panic 源逐字落在派单点名的那一支（`sampler_test.go:31`），入口是 `sampler.go:269` 的首次基线读树。

**被它吞掉的同包读数（基线名册 58 － 该发名册 42 ＝ 16 枚，按基线 `=== RUN` 顺序逐名）**：

| # | 被吞的用例 | # | 被吞的用例 |
| --- | --- | --- | --- |
| 1 | `TestSampleStateAllMetricsAndVerdicts` | 9 | `TestMarkTransitionTimestamps` |
| 2 | `TestSampleStateSleepingDiskWriteGateFails` | 10 | `TestCheckSettleVerifiesReleaseCounter` |
| 3 | `TestSampleStateSleepingTCPGateFails` | 11 | `TestCheckSettleNeverReachesCap` |
| 4 | `TestSampleStateHandleGateUsesRulingLimit` | 12 | `TestSamplerGoroutineAccountingFollowsRegistry`（AC#8 刚修的那枚） |
| 5 | `TestSampleStateWorkPeakMemoryIsTargetNotGate` | 13 | `TestLiveRegistryBaselineWithinSleepingGate` |
| 6 | `TestSampleStateLeakFixtureFlipsRed` | 14 | `TestThresholdTableCoversAllStates` |
| 7 | `TestSampleStateCPUTotalDrivenMean` | 15 | `TestSampleStateZeroSampleWindowFailsClosed`（**AC#1 的钉**） |
| 8 | `TestSampleStateUnknownStateAndReaderError` | 16 | `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（**AC#1 的正向对照腿**） |

⇒ "一挂吞一片"是读数不是感觉：**这 16 枚既不算红也不算绿**，包级只留一句 `FAIL github.com/CarlosShao/wisp/internal/observe`，而本票 AC#1／AC#8 花两格装上的钉就在名单里。
⚠ 如实一句位置依赖：探针插在采样用例族之前，所以吞 16 枚；同一枚探针若放在包尾文件里就吞 0 枚——**危害大小由"谁排在它后面"决定，这正是不能留活雷的理由**，不是反驳。

### 1.2 判据②：修法只动那枚仪器（要红，不静默）；`sampler.go` 语义一字未动

改的只有 seam（`internal/observe/sampler_test.go`，**22 增 0 删**，唯一 hunk 在 `ReadTree` 的"脚本耗尽"那一支）：

```go
	if f.i >= len(f.mu) {
		if len(f.mu) == 0 {
			// Nothing was ever scripted: there is no last reading to repeat.
			return TreeMetrics{}, errFakeTreeNoScript
		}
		return f.mu[len(f.mu)-1], nil
	}
```

`errFakeTreeNoScript` 的文案（同文件 `:38`）：`observe test seam: fakeTree has no scripted reads (bare &fakeTree{} is not a measurement; script mu/current or set err)`。

- **只动仪器**：commit `f4c7062` 的 `git show --numstat` ＝ `111 0 internal/observe/sampler_faketree_guard_136_test.go` ＋ `22 0 internal/observe/sampler_test.go`，**没有第三枚文件**；`git diff 45c8d1c..f4c7062 -- internal/observe/sampler.go` **无输出** ⇒ 采样器语义（含 `SampleState:269-272` 的错误通道、`:290` 的零足迹丢弃、`:332-340` 的 fail-closed 门）未被本程触碰，也没有"为绕开它改语义"。
- **要红不静默**：守卫不返回零值、不 `t.Skip`、不静默 `return`；它走 `TreeReader` 自己的错误通道，误用的那枚用例经 `sampler.go:269-272` 拿到 err 并 `t.Fatal` ⇒ **单枚红、别的照跑**。有脚本时"越界重复末读"的老行为**逐字保留**（`len(f.mu)==0` 才拒），由控制腿钉住。
- **新增钉**（`internal/observe/sampler_faketree_guard_136_test.go`，111 增 0 删，三腿）：
  1. `TestFakeTreeEmptyScriptFailsClosedAndNotPanics`：裸 seam 直接调一次，必须 `errors.Is(err, errFakeTreeNoScript)`、必须带那句"no scripted reads"、必须给零值指标；调用包在 `recoverRead` 的 `recover` 里 ⇒ **若守卫被摘掉，这一腿把 panic 收成本枚自己的红**（而不是让二进制去代答）。
  2. `TestFakeTreeScriptedReadsStillRepeatLastReading`：控制腿，防守卫写宽（逢读就拒 ⇒ 这腿红）。
  3. `TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel`：下一个人的普通采样用例——必须经错误通道红、`rep` 必须是 nil、seam 那句话必须到得了调用方。

### 1.3 判据③：复跑整包 ⇒ panic=0、`RUN` 名册与基线逐名相同

| 发 | 树 | 四数（`-count=1 -v`） | panic | 名册账 |
| --- | --- | --- | --- | --- |
| 基线 | `tree1`＝`git archive f4c7062` | `rc=0 / RUN=61 / PASS=61 / FAIL=0 / SKIP=0` | 0 | 与 tree0 基线 `diff` **只多我三枚**（`TestFakeTreeEmptyScriptFailsClosedAndNotPanics`、`TestFakeTreeScriptedReadsStillRepeatLastReading`、`TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel`），**无改名、无消失、无转 SKIP** ⇒ 新增的绿不是被跳过 |
| 同一发探针落在**修好的**树上 | `tree1` ＋ §1.1 那枚探针（逐字同一份代码、同一插入位置） | `rc=1 / RUN=62 / PASS=61 / FAIL=1 / SKIP=0` | **0** | 名册＝tree1 基线 ＋ 探针一枚，**缺 0 枚**（改前那一发缺 16 枚） |
| 还原 | `tree1`（每发先无条件 restore） | `rc=0 / 61 / 61 / 0 / SKIP=0` | 0 | 22/22 枚 `.go` 与 pristine `diff -q` 逐字相同；并另证与仓库工作树那枚 `sampler_test.go` 逐字相同 |

探针在修好之后的那一枚红，原文（`/d/tmp/wisp136ac1213-tree1-probe-v.txt:92-94`）：

```
=== RUN   TestProbeAC13BareFakeTreeNormalSampling
    sampler_test.go:93: bare fake tree: err=resource: observe: baseline tree read: observe test seam: fakeTree has no scripted reads (bare &fakeTree{} is not a measurement; script mu/current or set err) rep=<nil>
--- FAIL: TestProbeAC13BareFakeTreeNormalSampling (0.00s)
```

⇒ 判据②"要红不许静默"与判据③"panic=0、名册完整"同时成立：**误用的读数由它自己那枚用例报出来，别人不再被代答**。

### 1.4 钉不是哑的：两发变异各响各的腿（都先证落地再读数，每发先 restore）

两发都落在 `tree1`（`internal/observe/sampler_test.go`，测试侧），落地证明＝`diff -u` 只有那一处 hunk ＋ `go build ./...` rc=0 ＋ **`go vet ./internal/observe/` rc=0**（`go build` 不编 `_test.go`，所以对测试文件的落地证明必须带 vet 这一发）。

| 发 | 改法 | 落地证据 | 整包 `-v` 读数 | 红名逐名 ＋ 红点 |
| --- | --- | --- | --- | --- |
| **MA** | 摘掉守卫（回到旧那一支：`if f.i >= len(f.mu) { return f.mu[len(f.mu)-1], nil }`） | `diff -u` 单 hunk（4 删 0 增）、`grep -n "f.mu\[len(f.mu)-1\]"` ＝ `49:`（守卫已不在）、`go build` rc=0、`go vet` rc=0 | `rc=1 / RUN=42 / PASS=40 / FAIL=2 / SKIP=0 / ^panic=2 行（真 panic 一枚）` | ① `TestFakeTreeEmptyScriptFailsClosedAndNotPanics` 红在 `sampler_faketree_guard_136_test.go:46`，消息逐字：`fakeTree.ReadTree panicked instead of failing closed: runtime error: index out of range [-1] (a panic here takes the whole package's readings with it)`；② `TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel` `--- FAIL`（`_test.go:103` 之前就被 panic 打断）；控制腿 `TestFakeTreeScriptedReadsStillRepeatLastReading` 同发仍 `--- PASS`。**panic 复发并再吞 19 枚**（基线 61－该发 42＝19，逐名见 §1.5）⇒ 摘掉守卫就是回到本格的病，红名先到我这两腿 |
| **AB** | 守卫在，但把拒绝换成**静默零值**：`return TreeMetrics{}, nil` | `diff -u` 单 hunk（1 增 1 删＋注释行）、`grep -n "MUTATION AB"` ＝ `50:`、`go build` rc=0、`go vet` rc=0 | `rc=1 / RUN=61 / PASS=59 / FAIL=2 / SKIP=0 / panic=0`，名册与基线逐名相同 | ① `TestFakeTreeEmptyScriptFailsClosedAndNotPanics` 红在 `:57`（`bare &fakeTree{} returned a reading it never measured (metrics={PIDs:0 …全部 0…}): the seam must fail closed`）；② `TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel` 红在 `:103`（把那份"看着正常"的报告整枚打出来：`Samples:[] SampleErrors:4 LastSampleError:read returned a zero private working set for a live tree … Verdicts:[… {Metric:sampling Measured:0 valid / 4 errors Pass:false Gate:true}] Pass:false`）；控制腿仍绿 ⇒ **"不许静默"这一半是有牙的，不是注释** |

**本程没有跑过任何 `-run` 定点绕过**：§1.3／§1.4 的每一个数都是整包 `-count=1 -v` 的读数。

### 1.5 MA 那一发被吞的 19 枚（按 tree1 基线名册顺序，逐名）

`TestCheckSettleZeroTrustworthySamplesFailsClosed`、`TestCheckSettleTrustworthyReadsAreRecorded`、`TestSLOStateNamesPinnedToMachineStates`、`TestSampleStateAllMetricsAndVerdicts`、`TestSampleStateSleepingDiskWriteGateFails`、`TestSampleStateSleepingTCPGateFails`、`TestSampleStateHandleGateUsesRulingLimit`、`TestSampleStateWorkPeakMemoryIsTargetNotGate`、`TestSampleStateLeakFixtureFlipsRed`、`TestSampleStateCPUTotalDrivenMean`、`TestSampleStateUnknownStateAndReaderError`、`TestMarkTransitionTimestamps`、`TestCheckSettleVerifiesReleaseCounter`、`TestCheckSettleNeverReachesCap`、`TestSamplerGoroutineAccountingFollowsRegistry`、`TestLiveRegistryBaselineWithinSleepingGate`、`TestThresholdTableCoversAllStates`、`TestSampleStateZeroSampleWindowFailsClosed`、`TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`。

### 1.6 与 AC#11 分开（派单明令：本格不碰 flake）

本程在 `-count=2 -v` 门禁那一发**命中既有 flake 一枚**：`TestNoopTaskReturnsToBaseline` 红在 `goroutine_test.go:29: live count mid-task = 2, want 3`（第二发 `-count=2` 复跑 `rc=0 / 122 / 122 / 0 / SKIP=0`）。**未修、未 Skip、未调阈值**（那是 AC#11 的账，判据要求 ≥20 发复现率，本程不做）。两处如实登记：①票面 AC#11 记的红点是 `goroutine_test.go:33`（`PerTask mid-task`），本程命中的是同一枚用例里前一跳的 `:29`（`live count mid-task`，`reg.Count()`）——**同函数两跳、缺的都是同一个 happens-before**，不改判、只记差异；②本格所有三态都判在"红名是本程仪器"的读数上，两态复绿那两发 `SKIP=0`、`panic=0`。
