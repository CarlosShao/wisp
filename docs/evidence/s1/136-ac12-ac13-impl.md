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

（补记，写于 §2 之后、不改上面原文：AC#12 的 MD 变异那一发里同一枚 flake **再次命中同一红点** `goroutine_test.go:29: live count mid-task = 2, want 3`，故本程整包级读数累计命中 **2 发**；AC#12／AC#13 的任何判据都不因此改动——两格的红名都是本程仪器，`SKIP=0`、`panic=0`。）

---

## §2 AC#12：四判据逐判（靶＝`R-136-7`：settle 那侧"这一窗丢了几次读数"在报告里查不到）

### 2.1 先量"改前到底是什么形状"——同一份探针跑在两棵树上（探针文件只在仓外快照树，不进仓库）

探针只经下游真正读的那一面：**`json.Marshal(rep)` → key 清点**，不读任何 Go 字段，所以同一枚文件在两棵树上都编得过（`/d/tmp/wisp136ac1213-ac12probe_test.go`）。这类读数用 `-run TestProbeAC12` 定点取，**与判据用的整包读数分开**。

| 形状（本程自己造的 fixture） | 树 | 探针原文读数 |
| --- | --- | --- |
| 整窗只第 1 枚可信（seam 实际取到 10 次读树） | `tree1`＝`f4c7062`（**无 AC#12**） | `seam took 10 reads \| wire samples=1 \| sample_errors ABSENT \| last_sample_error ABSENT \| dropped_reads ABSENT \| back_within_cap_ms=10 final_bytes=4.194304e+06 pass=true` |
| 同上 | `tree2b`＝`c03aee3`（**有 AC#12**） | `seam took 10 reads \| wire samples=1 \| sample_errors PRESENT=9 \| last_sample_error PRESENT="read: acceptor probe shape: transient tree read failure" \| dropped_reads ABSENT \| … pass=true` |
| 一半读数报错（按读序号交替，取到 10 次） | `tree1`（无 AC#12） | `seam took 10 reads \| wire samples=6 \| sample_errors ABSENT \| last_sample_error ABSENT \| dropped_reads ABSENT \| … pass=true` |
| 同上 | `tree2b`（有 AC#12） | `seam took 10 reads \| wire samples=6 \| sample_errors PRESENT=4 \| last_sample_error PRESENT="read: …" \| dropped_reads ABSENT \| … pass=true` |
| 出线 key 全清点 | 两树 | 改前 11 枚：`[back_within_cap_ms cap_bytes elapsed_ms final_bytes free_os_memory_count free_os_memory_requested pass peak_bytes samples started_at target_state]`；改后 13 枚＝**同样 11 枚 ＋ `sample_errors` ＋ `last_sample_error`**（无改名、无删除） |

⇒ 终裁方 §3.2 那两发探针（`samples=1 back=60 pass=true`／`samples=2 pass=true`）在本程锚点上是**复现成立**的（枚数与毫秒随 ticker 抖动，`pass=true` 与三枚 key 的 ABSENT 是稳定项）；改后同一形状的窗口**丢 9 枚就说丢 9 枚**，`pass` 的颜色不动（本格要的是自陈，不是改判据）。
⚠ 如实两点：①探针里"一半"那一形取到 `samples=6 / sample_errors=4`（探针的 fixture 让第 1 枚恒可信，故非严格 5/5）；**仓库里的钉**用的是严格按读序号交替的 fixture，实测 `kept=5 lost=5`，并断言 `|kept-lost|<=1`；②`dropped_reads` 这一枚 key **在两形里都仍然 ABSENT**——本程按"与 `StateReport` 同形"办，只搬 `sample_errors`＋`last_sample_error` 两枚（终裁方给 R-136-7 的修法方向也正是这两枚），第三枚名字是探针当初顺手试的候选，不是 `StateReport` 有的东西。若终裁认为必须再有一枚独立计数字段，那是本程未做档（§5），不是被漏掉的。
⚠ "可信样本数"这一件按同一族形状 rides 在 `samples` 数组（json 标签 `samples`）上，与 `StateReport` 完全同形（那侧也没有独立 count 字段，计数在 `sampling` 门行的 `"%d valid / %d errors"` 里）；本程把这条解释**明写在此请终裁过目**，不当它已成立。

### 2.2 判据①：`SettleReport` 带上字段且有 json 标签；解冻范围实际用了哪几行（＋一处超出授权，报回）

新字段（`internal/observe/sampler.go`，改后坐标 `:442-455`）：

```go
	Samples               []Sample `json:"samples"`
	// SampleErrors counts the reads this settle window DROPPED: ... （注释写明理由，指向 :164-168 与 R-136-7）
	SampleErrors int `json:"sample_errors"`
	// LastSampleError keeps WHY the most recent read was dropped, so the
	// number above can be read without guessing.
	LastSampleError string `json:"last_sample_error,omitempty"`
	Pass            bool   `json:"pass"`
```

计数落点（`CheckSettle` 循环，改后坐标 `:488-500`），丢弃判据本身逐字不变（可信样本仍只有 `err == nil && PrivateWorkingSetBytes > 0` 那一支能进 `Samples`）：

```go
		m, err := s.tree.ReadTree()
		if err != nil {
			rep.SampleErrors++
			rep.LastSampleError = "read: " + err.Error()
		} else if m.PrivateWorkingSetBytes <= 0 {
			// Zero-footprint reads of a live tree are untrustworthy (see
			// SampleState) and are dropped, never recorded as progress.
			rep.SampleErrors++
			rep.LastSampleError = "read returned a zero private working set for a live tree"
		} else {
			... 原样：append + FinalBytes + BackWithinCapMS ...
		}
```

**实际用了哪几行（`git diff -U0 45c8d1c..c03aee3 -- internal/observe/sampler.go` 的旧侧 hunk 头，锚点坐标）**：

| hunk（旧侧） | 内容 | 在 `:470-501` 之内？ |
| --- | --- | --- |
| `@@ -442 +442,14 @@` | `SettleReport` 的字段区：旧 `:442` 那一行 `Pass … json:"pass"` 换成 14 行（两枚新字段＋注释；gofmt 重排了同组对齐） | ❌ **不在** |
| `@@ -477 +490,4 @@` | 旧 `:477` 的 `if err == nil && m.PrivateWorkingSetBytes > 0 {` 拆成 `if err != nil { … } else if … <= 0 { … } else {` | ✅ 在 |
| `@@ -479,0 +496,3 @@` | 旧 `:479` 之后插入两枚计数＋一句原因 | ✅ 在 |

`git diff --numstat 45c8d1c..c03aee3 -- internal/observe/sampler.go` ＝ **`21 2`**（21 插 2 删；删的两行就是上面 hunk 里被替换掉的旧 `:442` 与旧 `:477`）；AC#12 全范围 `git diff --numstat f4c7062..c03aee3` 里属于本程的只有 `21 2 internal/observe/sampler.go` ＋ `289 0 internal/observe/sampler_settle_coverage_136_test.go`（同范围里另外四枚 `.md` 是兄弟在飞的票 135/137/台账，不是本程）。

> ⚠ **超出授权范围的那一处，在此报回（不自证合理）**：具名解冻写的是 `sampler.go:470-501`，但判据①要求"报告结构体带上字段"，而 `SettleReport` 的声明在锚点坐标 `:431-443`——**加字段必然动 `:442`**。本程判断这是派单给范围时漏写了 struct 那一截（终裁方 R-136-7 的修法方向本身就写着"给 `SettleReport` 补两枚字段"），于是按判据①做了**最小插入**并把越界行号如实钉在这张表里，供终裁按"接受／要求回退"处理。除此之外 `sampler.go` 零改动：`SampleState` 那侧（`:247-348`）**零 hunk**，`thresholds.go`、任何 golden、`scripts/`、`.github/` 在 `45c8d1c..HEAD` 全范围 **`git diff --name-only` 无输出**，D32 的 CPU≤0.5%／RSS≤25MB 未被任何一发触碰（本程所有变异都在仓外快照树，且只落 `sampler.go` 的 `:490-497` 那几行）。

### 2.3 判据②：两发探针各钉一枚用例（另两腿把两支丢弃路与"恒真"堵住）

`internal/observe/sampler_settle_coverage_136_test.go`（commit `f53ad5c`，`289 增 0 删`；`c03aee3` 只做了一处归因精修，`7 增 2 删`，见 §2.7）。四腿：

| 腿 | 对应探针／性质 | 断言落点（本文件内行号；**加粗**＝本程变异读数里真红过的那一条） |
| --- | --- | --- |
| `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | 探针①"整窗只 1 枚可信"（`samples=1 pass=true` 那一发） | 前提 `:126`（`reads>=3`）／`:130`（`len(Samples)==1`）；性质 **:134**（`sample_errors==reads-1`）、`:138`（`>=2`）、**:141**（`last_sample_error` 带原因）、`:144`（`pass==true`）；出线 `:150`/`:153`（`samples` 是长度 1 的 list）、`:156`/`:158`（`sample_errors` key 在且与报告同数）、`:161`/`:163`（`last_sample_error` 在且同因）、`:167`/`:169`（`留下＋丢掉==取过`） |
| `TestCheckSettleHalfTheReadsFailedReportsItsLoss` | 探针②"一半读数报错"（`samples=2 pass=true` 那一发） | 前提 `:189`（`reads>=4`）／`:192`（`kept>=2`）；性质 **:197**（`lost>=2`，消息即"a half-covered window must report its losses"）、`:200`（`kept+lost==reads`）、`:203`（`\|kept-lost\|<=1`）、**:206**（原因在）、`:209`（`pass==true`）；出线 `:214`（`sample_errors==lost`）、`:217`（`samples` 长度==kept） |
| `TestCheckSettleZeroFootprintDropsAreCountedToo` | **另一支**丢弃路（AC#9 钉住的零足迹 fail-closed）也要被数到 | 前提 `:241`／**:244**（`len(Samples)==1`）；性质 **:247**（`sample_errors==reads-1` 且 `>=2`）、`:250`（原因句**必须与 `StateReport` 那侧逐字同一句**）、`:253`、`:257`（出线同数） |
| `TestCheckSettleFullyMeasuredWindowReportsNoLoss` | 正向对照（防恒真：逢窗口就报丢读数） | `:268`/`:271`（每枚读树都留下）、`:274`（`SampleErrors==0`）、`:277`（`LastSampleError==""`）、`:280`（绿）、`:284`（出线 `sample_errors==0` **在**）、`:287`（`last_sample_error` **缺席**，omitempty 的形状） |

fixture 全部落在冻结 Sleeping cap（25MB）之下，释放腿（`released=true`＋`FreeOSMemoryCount>0`）与内存比较**恒满足**，所以这些腿能红的地方只有"覆盖率的自陈"；seam 自己也照 AC#13 的纪律不留活雷（`coverageScriptTree` 空脚本直接拒，不索引空切片）。

### 2.4 判据③：落变异 ⇒ 用例转红、红名逐名（每发先 restore、先证落地再读数）

四发都落在 `tree2b`（`git archive c03aee3`）的 `internal/observe/sampler.go`，驱动器 `/d/tmp/wisp136ac1213-mut.py`：锚点文本命中数≠1 ⇒ 拒绝落发并回滚（本程未触发拒绝）。每发落地证明＝`diff -u`（只有那一处 hunk，原文见日志）＋`grep -n` 出被改后那一行＋`go build ./...` rc=0，然后才取整包 `-count=1 -v`。

| 发 | 改法 | 落地证明 | 整包四数 | 红名逐名 ＋ 红点 |
| --- | --- | --- | --- | --- |
| **MD**（判据③点名那发："把计数摘掉"） | err 支的 `rep.SampleErrors++` 摘掉（读树仍被丢弃） | `grep -n "MUTATION MD"` ⇒ `491:`；`diff -u` 单 hunk（`-rep.SampleErrors++`）；`go build ./...` rc=0 | `rc=1 / RUN=65 / PASS=62 / FAIL=3 / SKIP=0 / panic=0` | ① `TestCheckSettleSingleTrustworthyReadReportsItsLoss` `sampler_settle_coverage_136_test.go:134`＝`report counted 0 dropped reads but the seam took 10 reads and kept 1: 9 unaccounted`；② `TestCheckSettleHalfTheReadsFailedReportsItsLoss` `:197`＝`the seam lost 5 of 10 reads but the report says sample_errors=0: a half-covered window must report its losses`；③ 既有 flake `TestNoopTaskReturnsToBaseline`（`goroutine_test.go:29`，见 §1.6 补记，非本程仪器、未计入本格）；腿 3／腿 4 同发仍 `--- PASS` |
| **ME** | 零足迹支的 `rep.SampleErrors++` 摘掉 | `grep -n "MUTATION ME"` ⇒ `496:`；单 hunk；`go build` rc=0 | `rc=1 / 65 / 64 / 1 / SKIP0 / panic0` | 只一枚：`TestCheckSettleZeroFootprintDropsAreCountedToo` `:247`＝`zero-footprint drops must be counted: sample_errors=0 reads=10 kept=1` |
| **MF** | err 支的原因赋值摘掉（计数留着） | `grep -n "MUTATION MF"` ⇒ `492:`；单 hunk；`go build` rc=0 | `rc=1 / 65 / 63 / 2 / SKIP0 / panic0` | 两枚：`…SingleTrustworthyReadReportsItsLoss` `:141`＝`the report must carry the reason it dropped reads, got ""`；`…HalfTheReadsFailedReportsItsLoss` `:206`＝`a dropped read must leave its reason behind, got ""` |
| **MC**（派单没要求，本程自加：证 §2.2 的重排没把隔壁 AC#9 那枚钉弄钝） | `} else if m.PrivateWorkingSetBytes <= 0 {` ⇒ `< 0`（放宽 fail-closed 足迹判据） | `grep -n` ⇒ `493:`（`< 0`）；单 hunk；`go build` rc=0 | `rc=1 / 65 / 63 / 2 / SKIP0 / panic0` | ① **`TestCheckSettleZeroTrustworthySamplesFailsClosed` 红在 `sampler_settle_zerosample_136_test.go:66`**（消息就是那串病形：`recorded 5 samples … Pass:true`，红点与终裁方 V3 那一发同一条）；② 本程腿 3 红在 `:244`（`want 1 recorded sample, got 10`）；腿 1／腿 2／腿 4 同发仍绿 |

⇒ **各发只响各的腿**：MD 只碰 err 支 ⇒ 腿 3（零足迹支）不动；ME 只碰零足迹支 ⇒ 腿 1/2 不动；MF 摘原因不摘计数 ⇒ 计数断言仍过、原因断言红；MC 放宽足迹判据 ⇒ AC#9 那枚钉与本程腿 3 同时红，说明本程把 `:477` 拆成三分支之后，**那一族 fail-closed 的牙还在**。

### 2.5 判据④：还原 ⇒ 复绿，三态齐

| 态 | 读数 |
| --- | --- |
| 未变异 | `tree2b`：`rc=0 / RUN=65 / PASS=65 / FAIL=0 / SKIP=0 / panic=0`（名册＝`tree1` 的 61 枚 ＋ 本程 4 枚，逐名 diff 只多出那 4 行，无改名无消失无转 SKIP） |
| 变异 | 上表四发（红名逐名，`SKIP=0`、`panic=0` 每发都是） |
| 还原 | `python …-mut.py restore` ⇒ `RESTORED`，`diff -q` **23/23** 枚 `.go` 与 pristine 逐字无输出；`go test -count=1 -v` ⇒ `rc=0 / 65 / 65 / 0 / SKIP0 / panic0`；另证 `tree2b/internal/observe/sampler.go` 与仓库工作树那枚**逐字相同** |
| 门禁级复跑 | `-count=2 -v`（还原后的 tree2b）＝`rc=0 / RUN=130 / PASS=130 / FAIL=0 / SKIP=0 / panic=0`，去重名册 65 枚与 `-count=1` 逐名相同 |

### 2.6 判据⑤：有没有被迫去动 `scripts/slo-check.ps1`／golden／阈值 ⇒ **没有，一处未动**

- `git diff --name-only 45c8d1c..HEAD -- internal/observe/thresholds.go scripts/ .github/` **无输出**；`slo-check.ps1` 全程只读（读它 `:341-361` 那段：只取 `$settleReport.pass` 与 `$settleReport.settle.free_os_memory_count`，**没有 key 集合断言、没有字段计数**）。
- 全仓没有消费 settle 出线的 golden：`git grep -ln "back_within_cap_ms|free_os_memory_*"` 命中的 tracked 文件只有 `docs/SLO.md`（那是一句**读数**叙述，不是 schema 清单）、`docs/evidence/s1/66/*.json`（票 66 归档的**报告原件**，没有任何仪器读它：`git grep -n "66-settle-1|66-full-subset"` 只命中票面与 `docs/SLO.md` 的重跑命令）、若干 `.md` 证据文件，加 `internal/observe/**` 与本枚新代码。
- 为了让"旧读者读新报告"这件事不靠推断，本程把 `cmd/wisp`（唯一 Go 侧消费者）的**测试二进制按改动后的 `observe` 重新链接**了一次：`go test -c -o /d/tmp/wisp136ac1213-cmdwisp.test.exe ./cmd/wisp/` rc=0（32,584,099 字节）。**只编不跑**：快照树里没有 `third_party/sherpa-onnx`（`0xc0000135` 是缺 DLL 的宿主现象，不是代码现象），而 `cmd/wisp` 的测试此刻归票 133 在飞——跑它等于踩别人的地界。⇒ "cmd/wisp 端到端在改动后仍全绿"记进 §5 未做档。

### 2.7 一处自纠（append-only，不回改已提交的断言）

`f53ad5c` 里腿 2 把"取到几枚"和"丢了几枚要数出来"合写成一条 `if kept < 2 || lost < 2 { t.Fatalf("precondition broken: …") }`。MD 那一发实测把它判红在合写跳上、消息写成 `precondition broken: kept=5 lost=0 of 10 reads`——**归因含糊**（丢读数没被数到是本腿的本题，不是前提）。`c03aee3` 拆成 `reads<4`（计时前提）／`kept<2`（fixture 前提）／`lost<2`（性质，红点带 `must report its losses`）三跳，断言强度只增不减（MD 复跑因此红点从 `:192` 移到 `:197`，见 §2.4）。生产码未动、阈值未动。第一次 MD 的读数（`tree2`，`rc=1 / 65 / 63 / 2`，红名两枚同名同因）留档 `/d/tmp/wisp136ac1213-tree2-MD-v.txt`，不进结论。

---

## §3 谁在消费 `SettleReport`（加字段会改出线形状 ⇒ 这条不默认"无影响"）

全仓 `grep -rn "SettleReport"`（含 `.go`／`.ps1`／`.sh`／`.md`）命中**非证据**文件只有四处，逐处盘：

| 消费者 | 怎么用 | 加字段会不会让它读不到旧字段 | 我的凭据 |
| --- | --- | --- | --- |
| `cmd/wisp/slo_windows.go:118` `Settle *observe.SettleReport json:"settle,omitempty"` | 生产侧**序列化**（`writeSLO` 用 `json.MarshalIndent`） | 否：只**增** key | §2.1 的 key 全清点：改前 11 枚、改后同样 11 枚 ＋ `sample_errors` ＋ `last_sample_error`，**无改名、无删除**（两串原文都在 §2.1 表里） |
| `cmd/wisp/slo_windows.go:287` `run.Pass = run.Settle.Pass` | 只读 `.Pass` | 否 | 同一枚字段未动（本程 `sampler.go` 的 hunk 不含 `Pass` 的语义行；出线里 `"pass"` 仍在，§2.1 探针两形都打到 `pass=true`） |
| `cmd/wisp/slo_windows.go:523-534`（`collectReport`，subject 模式把自己的报告**反序列化**回来） | `json.Unmarshal(data, &rep)` 进同一个 `sloRun` 类型 | 否：同版同类型；且 `encoding/json` 默认忽略未知 key、缺 key 留零值 | `git grep -n DisallowUnknownFields` 全仓只有 `internal/config/parse.go:72/:123`（配置解析面），slo 这条链上**没有**严格解析器 |
| `scripts/slo-check.ps1:348-355` | `ConvertFrom-Json` 后只取 `$settleReport.pass` 与 `$settleReport.settle.free_os_memory_count` | 否：两枚都还是原 key | 上面那段脚本原文（`:341-361`）逐字读过；`git grep -n "\.settle" -- scripts/ frontend/ tools/` 只命中 `:353` 这一行 |

再往外补三刀（派单没点名，但"下游"不能只问 Go 侧）：

- **golden／阈值面**：`git grep -ln "back_within_cap_ms|free_os_memory_requested|free_os_memory_count"` 的 tracked 命中只有 `docs/SLO.md`（那句是**一次 run 的读数叙述**——`settle 行本次实测：peak 126.5MB → 回 cap(25MB) 261ms、free_os_memory_count=2、pass=true`——不是 schema 清单）、`docs/evidence/s1/66/66-settle-1.json` 与 `66-full-subset-slo-report.json`（票 66 归档的**报告原件**；`git grep -n "66-settle-1|66-full-subset"` 显示没有任何仪器读它，命中的是票面叙述与 `docs/SLO.md` 的重跑命令），加本程自己那两枚 `_test.go`。⇒ **没有任何用例或脚本拿旧 settle 报告当期望值比对**，本程因此未动 golden、未动 `slo-check.ps1`（判据⑤）。
- **`scripts/slo-freshness.sh`**：它的 `--jq` 全在 GitHub API 响应上（`.workflow_runs[]`／`.jobs[]`／`.artifacts[]`，`:221/:231/:312`），不解析报告内容 ⇒ 不受影响。
- **前端**：`git grep -n settle -- frontend/` 只命中 `frontend/src/components/ai-native/thinking.tsx` 的英文散文（"then settles"，讲动画），不是数据字段。
- **`cmd/wisp` 的测试**：`grep -rn "Settle\|settle" cmd/wisp/*_test.go` **无输出** ⇒ 没有第二枚仪器对 settle 形状有断言（也意味着：这一族目前没有 cmd 侧的钉，那是 AC#10 的账）。

**必须留的一条副作用（不粉饰）**：`json.Unmarshal` 对**缺 key** 给零值 ⇒ 拿这份改动后的二进制去读**改前**留下的 settle 报告（例如重放 `docs/evidence/s1/66/66-settle-1.json`），`SampleErrors` 会读成 0、`LastSampleError` 空——"看不出丢过"会被读成"没丢过"。今天没有这条路径（`collectReport` 读的是同一次 run 里同一版二进制写的文件；`slo-check.ps1` 不读这两个 key），但它是"字段是后加的"这件事的真实边界，登记为 §5-adjacent 的一条已知限制，不靠"应该没人这么读"糊过去。

**另一条如实观察（本程未动、也不该本程动）**：`cmd/wisp/slo_windows.go:621-623` 在 `CheckSettle` 报错时造一枚合成的失败报告 `&observe.SettleReport{TargetState: …, Pass: false}`——它的 `sample_errors` 是 0，含义是"根本没测"，不是"零丢失"。那一形 `pass=false` 已经挡掉静默绿，且落点在 `cmd/wisp`（票 133/AC#10 地界）⇒ 只登记，不改。

---

## §4 门禁原文（两格各一套；每套都跑在该格自己的纯净快照树上）

`AC#13` 的门禁面＝`tree1`（`git archive f4c7062`）；`AC#12` 的门禁面＝`tree2b`（`git archive c03aee3`）。工具链同一把尺：`go version go1.27.1 windows/amd64`、`/d/work/base/gopath/bin/gofumpt.exe --version` ＝ `v0.12.0 (go1.27.1)`，**全程未 `go install` 任何东西**。

| 门禁 | AC#13（tree1＝`f4c7062`） | AC#12（tree2b＝`c03aee3`） |
| --- | --- | --- |
| `gofmt -l internal/observe/` | 无输出 | 无输出 |
| `gofumpt -l internal/observe/` | 无输出（版本原文 `v0.12.0 (go1.27.1)`；⚠ CI 那一步是 `@latest` 未钉版本 ⇒ 本读数只在该二进制未被改版时有效） | 无输出（同版本、同边界） |
| `go vet ./internal/observe/`（原生） | rc=0 | rc=0 |
| `GOOS=linux go vet ./internal/observe/` | rc=0 | rc=0 |
| `go vet ./...`（原生整树） | rc=0，输出 **0 行** | rc=0，输出 **0 行** |
| `GOOS=linux go vet ./...`（整树交叉） | **rc=1，输出 3 行、一条诊断**：`package github.com/CarlosShao/wisp/cmd/wisp → imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx → imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8` | 同一形（3 行、一条诊断，逐字同） |
| 上面那条的逐错误归因（**不整树归因成"工具链假象"**） | 该诊断不含任何本仓 `file:line`，是**外部模块** `sherpa-onnx-go-linux@v1.13.8` 的 build constraint；把它点名的包剔掉后 `GOOS=linux go vet $(除 cmd/wisp)` **rc=1**、只剩一条 `package …/cmd/balldebug: build constraints exclude all Go files in D:\tmp\wisp136ac1213-tree1\cmd\balldebug`（我方树内文件，`cmd/balldebug` 在无 CGO/无 linux 形下没有可编文件）；再剔掉这两枚 ⇒ `go list ./...`=**33**、剔后 **31/31 枚非 cmd 包零输出 rc=0**（输出文件 0 字节） | 同结构：整树 3 行一条诊断；剔 `cmd/wisp` ⇒ 只剩 `cmd/balldebug` 那一条；剔两枚 cmd ⇒ **31/31 零输出 rc=0**（0 字节）。⚠ 与终裁方 §5.3 那句"30/30"差一枚：两枚锚点 `git ls-tree cmd/` 都是 `balldebug`/`llmrecord`/`wisp` 三枚、`go list ./...` 都是 33，差异只在**剔几枚 cmd**，不在包集合；本程按自己点名的两枚（授权失效面＝`cmd/wisp`、`cmd/balldebug`）报 31/31，不据此判谁错 |
| `go test -count=2 -v ./internal/observe/` 四数（**只从 `-v` 量**） | 第一发 `rc=1 / RUN=122 / PASS=121 / FAIL=1 / SKIP=0`，红名 `TestNoopTaskReturnsToBaseline`（`goroutine_test.go:29: live count mid-task = 2, want 3`，票面 AC#11 那枚既有 flake，见 §1.6）；第二发 `rc=0 / 122 / 122 / 0 / SKIP=0` | `rc=0 / RUN=130 / PASS=130 / FAIL=0 / SKIP=0 / panic=0`（该发未命中 flake） |
| `-count=2` 的**逐名差集**（四数之外必比的那件） | 基线（`45c8d1c`＝tree0）58 枚 → tree1 去重名册 61 枚，`diff` 只多出本程 3 枚（逐名见 §1.3），**无改名、无消失、无转 SKIP**；`-count=2` 去重后与 `-count=1` 名册逐字相同 | tree1 的 61 枚 → tree2b 的 65 枚，`diff` 只多出本程 4 枚（§2.5 逐名），`-count=2` 去重名册 65 枚与 `-count=1` **逐名相同**（`NAMES-IDENTICAL-vs-count1`） |
| 每一发的 `SKIP` 计数 | 本程 AC#13 全部读数（基线／探针／MA／AB／还原／`-count=2`×2 共 **7 发**）`--- SKIP` 计数**逐发为 0**；包内 `t.Skip` 出现次数仍只在 `sampler_test.go:313` 那句注释里（"never t.Skip, never a silent return"） | 本程 AC#12 全部读数（基线／探针×2 树／MD／ME／MF／MC／还原／`-count=2` 共 **9 发**）`--- SKIP` 计数**逐发为 0**、`^panic` 命中 0 |
| `sh scripts/d22scan.sh`（各 scope 不降） | 两形 rc=0；`ban #8 internal/` 由 tree0（`45c8d1c`）的 **402** → tree1 的 **403**＝本程那枚新 `_test.go`；其余 scope 逐字同数（`#1-5 internal/=203`、`cmd/=22`、`#6 frontend/=40`、`#7 internal/tools/=18`、`#8 design/=16`、`#8 frontend/=40`、`#8 cmd/=39`）；`d22scan: clean`，step-1 `runtests.sh: OK … PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` | 两形 rc=0；同一比对继续走到 tree2b 的 **404**（＋本程第二枚新 `_test.go`），**无任何 scope 下降**（其余数字同上，逐字未动） |

两枚 sha 都在写这行之前重新核过（`git cat-file -t f4c7062`／`git cat-file -t c03aee3` 均 `commit`；前者是本程 AC#13 的 commit、后者是 AC#12 的 commit）。本文件**没有**任何未经 `git cat-file -t` 现核的外来 sha；所有快照只从 `45c8d1c` 或我自己那三枚 commit 抽。

---

## §5 本程未做的档（诚实列，别当已验）

**AC#13 侧**

1. **只修了派站点名的 `sampler_test.go:31` 那一支**（"别一把改完，改一枚证明一枚"）。同族普查我自己重走了一遍：`grep -rn '\[len(.*)-1\]' internal/observe/*_test.go` 在本程最终树上＝**4 行**，其中两行是我自己写的注释、一行是被守卫护住的 `:53`，**唯一另一枚真读点是 `earlylog_130_test.go:175`**，它前面 `:168-171` 就有 `len(recs) != earlyLogMaxRecords+1 ⇒ t.Fatalf` 的前置长度守卫 ⇒ 判"不必改"，但**我没为它落变异真打**。`[0]` 直取本程未重做普查（AC#8 已裁：17 处/终裁方 22 行，本程读数 `grep -rn '\[0\]' internal/observe/*_test.go | wc -l`＝**23 行**，比终裁方那 22 行多的正是本程新文件里那 1 行注释）。
2. **没测"守卫被改成 `t.Skip`"那一形**：逻辑上被腿 1 的 `err==nil ⇒ t.Fatal` 覆盖（Skip 掉就不产生 err），但本程没有真打这一发。
3. 探针只在**仓外快照树**里（`tree0`/`tree1`），**不进仓库**；仓库里长期留下的是 §1.2 那三腿钉。⇒ 派单判据①"拍到 panic 并逐名列出被吞读数"是在快照上做的，仓库里没有一枚"逢裸 `&fakeTree{}` 就 panic"的用例（那是应当的：修好之后它不该再 panic）。
4. 未跑 `cmd/wisp`、未跑任何 `-run` 定点**判据**读数（`-run` 只用于 §2.1 的自造探针打印，与判据读数分开）。

**AC#12 侧**

5. **没有加第三枚 key `dropped_reads`**，也没加独立的"可信样本数"整数字段（§2.1 已写明理由与请终裁过目的位置）。若终裁按票面判据①字面要三枚独立字段，这一档就是缺口而不是选择。
6. **没有给 settle 造一枚与 `sampling` 同形的显式门行**（`SettleReport` 无 `Verdicts` 字段，加它是更大的出线形状改动；终裁方 §3.4 已判"AC#9 那句不需要"，本格判据也没要）。
7. **R-136-8 未做**（`Samples` 为空时出线是 `"samples":null`，使 AC#9 腿 A 那一支防呆永不响）——不在 AC#12 的四判据里，本程未动、也未顺手动。
8. 未跑 `wisp slo -settle` 端到端、未跑 `scripts/slo-check.ps1`（PowerShell 那一腿零读数）；`cmd/wisp` 的测试二进制只**编译链接**（§2.6）未执行。⇒ "生产侧真报告里 `sample_errors` 会是多少"今天仍是读码＋包内读数，不是端到端（那是 AC#10 的地界）。
9. **未测 linux 面**（容器原生 `go test -v ./internal/observe/` 一枚未跑）、未测 `-race`／`-shuffle`；所有名册与逐名账都是**当前执行顺序**下的读数。
10. **未取 CI run id**（本程只 commit 未 push）；未复算 CI 那一步（`bash scripts/portable-tests.sh --scope=core`）整条。
11. 未动票面任何 AC 勾；未做 AC#11（flake 复现率 ≥20 发）、AC#2-AC#7、AC#10、AC#1 各腿。

**仪器级异常一条（不是被验面的性质，但会影响别人读我的时间戳）**

12. 本机墙钟在本程两次工具调用之间**跳了 8h29m**：`date -u` 原文 `Wed Sep 23 16:18:49 UTC 2026` 的下一发是 `Thu Sep 24 00:48:03 UTC 2026`；三条交叉量法＝台件 mtime（`/d/tmp/wisp136ac1213-d22-tree1.txt` mtime `00:18:29 +0800`，而 `date` 报 `08:48 +0800`）、`git log` 提交时间戳、跳变前后各一发 `date -u`。⇒ **本文件里任意两个时间戳相减都不许当耗时读**；本程所有读数不依赖墙上时间（计时类断言一枚未跑，settle 的 `within/interval` 是相对时长、不是墙钟）。兄弟一程已把同一现象记进 `docs/reports/pending-and-issues.md`（A 账）与 137 验收 §8，本条不另立新账、只让本文件自带这个边界。

13. **一条自报的纪律偏离（"临时件只建不删"）**：§2.1 那枚探针我先后拷进 `tree1` 与 `tree2` 各一份，取完两形读数后**用 `rm -f` 删掉了树内那两份副本**——成因是"每发变异前先无条件 restore"要求树上不能留下探针（否则名册差集里混着非判据用例）。**原件始终在盘**：`/d/tmp/wisp136ac1213-ac12probe_test.go`（以及 `/d/tmp/wisp136ac1213-mut.py`、全部 `-v` 日志、`tree0`/`tree1`/`tree2`/`tree2b`/两枚 pristine 目录，均未删）。⇒ 快照树本身一枚未删；被删的是我在树内的临时副本，可由原件一次 `cp` 复原。登记在此不按"顺手清理"混过去。

---

## §6 本程 commit 清单与 `next=`

| 节 | sha | 内容（净面） |
| --- | --- | --- |
| AC#13②③ 代码 | `f4c7062` | `internal/observe/sampler_test.go` `22 0`（seam 守卫）＋ `internal/observe/sampler_faketree_guard_136_test.go` `111 0`（三腿钉）；**`sampler.go` 零命中** |
| AC#13 证据 §0-§1 | `a0075ee` | 本文件首两节（锚点＋复现/修法/复跑＋MA/AB 两发变异＋19 枚被吞逐名） |
| AC#12① 生产码 | `5c1529a` | `internal/observe/sampler.go` `21 2`（两枚字段＋两支计数；越界那一行在 §2.2 表里报回） |
| AC#12② 钉 | `f53ad5c` | `internal/observe/sampler_settle_coverage_136_test.go` `289 0`（四腿） |
| AC#12② 精修 | `c03aee3` | 同文件 `7 2`（腿 2 合写断言拆三跳，§2.7 自纠） |
| AC#12 证据 §2 | `0f49259` | 四判据逐判＋解冻范围实用行号＋MD/ME/MF/MC 四发 |
| 两格证据 §3-§4 | `7dc8086` | 消费者核查＋门禁原文两格各一套 |
| 票面 log（AC#13＋AC#12 各一条） | `1f7875d` | append-only 两行（`16 增 0 删`，别人原文一字未改），未翻勾 |

### 6.1 收尾补记（这一段写在上表之后，故把表里来不及登记的最后一枚一并列出）

| 节 | sha | 内容 |
| --- | --- | --- |
| 证据 §5-§7 | `888486c` | 未做的档 12 条 ＋ 共树/暂存账 ＋ 两个计数 ＋ §7.1 sha 核法表 |
| 本节（§6.1 收尾补记） | 见 `git log --oneline -1 -- docs/evidence/s1/136-ac12-ac13-impl.md` | 收尾复跑（下表）＋共树噪声一条 |

收尾复跑（**仓库工作树**，此刻 `HEAD=1f7875d`；`date -u` 原文 `Thu Sep 24 01:19:44 UTC 2026` ⇒ 本机 **09:19:44 +08**）：`gofmt -l internal/observe/` 无输出；`go build ./...` rc=0；`go test -count=1 -v ./internal/observe/` ＝ **`rc=0 / RUN=65 / PASS=65 / FAIL=0 / SKIP=0 / ^panic=0`**（原文 `/d/tmp/wisp136ac1213-final-worktree-v.txt`）；并证工作树那枚 `internal/observe/sampler.go` 与 §2 全部变异读数所依据的 `pristine-tree2b` 拷贝 `diff -q` **逐字相同** ⇒ 本文件的读数面与最终交付面同版。

共树噪声一条（不是注入、也不是谁的错）：本程收尾时工作树里有**兄弟在飞**的三枚 `internal/winsec/*_test.go` 处于 ` M`（`ancestor_separator_108_other_test.go`、`placement_leaf_118_other_test.go`、`placement_symlink_113_other_test.go`，票 137 的地界）。本程九枚 commit 每一枚的 `git diff --cached --name-only` 都只列出我自己的路径，没把它们带走；本程也**没有**读它们的内容来支撑任何结论。

暂存账：本程**每一枚** commit 都跑过 `git add -- <显式路径>` ＋ `git diff --cached --name-only`，暂存清单里只出现过上面那些我自己的路径；`git commit -q -F - -- <同一批路径>` 全部带 pathspec；未 push；未 `--amend`／`reset`／`rebase`／`stash`／`checkout .`；仓内未建 worktree。

共树账：`45c8d1c..7dc8086` 之间兄弟推进 **11 枚** commit（旧→新：`e2ba9f0` `b802432` `286de6c` `8c7ad1f` `7ad0ec6` `94b7267` `69414d3` `99ac876` `c04221e` `097be8a` `b615092`，全是 `.md`／票面／台账；写这一串之前逐枚由 `git log --format=%h %s --reverse 45c8d1c..HEAD` 现取，非抄来的 sha），我的 commit 落在当时 tip 上属共树正常形态。硬核对两发：`git diff --name-status 45c8d1c..HEAD -- internal/observe/` ＝ **只有本程那四枚**（`M sampler.go`、`M sampler_test.go`、`A sampler_faketree_guard_136_test.go`、`A sampler_settle_coverage_136_test.go`）；工作树 `git status --short` 全程除我自己的文件外只剩那枚未跟踪的 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`（未读未动）。

`next=` ①派**非实现者**按票面 AC#13①②③／AC#12①②③④⑤ 终裁并翻勾（本程一枚未翻）；②**先裁 §2.2 那处越界**——具名解冻是 `:470-501`，判据① 必然要动 `SettleReport` 声明区（锚点 `:442`），本程按判据做了最小插入并报回：接受 ⇒ 请把解冻范围补成"`:431-443` ＋ `:470-501`"或等价措辞，撤销 ⇒ 回退 `5c1529a` 的 struct hunk 即可（计数那两发是同一枚 commit 里的另一半，要整体看）；③终裁顺手判一下 §2.1 那两条解释（"可信样本数 rides on `samples`"、"不加 `dropped_reads`"）够不够判据① 的字面；④AC#12 的端到端（生产侧真报告里 `sample_errors` 非零）仍归 AC#10／`cmd/wisp` 空出来之后；⑤R-136-8 与 §5-7 那几条未做档没被本格覆盖，别当已修。

---

## §7 两个计数（分栏，不混装）＋ 本程引用到的每一枚 sha 的现核账

写这一节之后重新 `date -u`：原文 **`Thu Sep 24 01:15:28 UTC 2026`** ⇒ 本机 **09-24 09:15:28 +08**（+8 手工换算，没塞进 `date` 的格式串；本节的落笔时间以这一发为准）。判据用派单那四条（①点名的路径／对象在本机真不真 ②内容是否削弱 owner 权威或放宽判据 ③声称的动作能否盘上复核 ④**是否让我少取证**——即使不越权也按注入登记并继续取证）。"像不像系统提示"**不作判据**。

| 栏 | 计数 | 逐条出处（工具名 ＋ 命令／位置前 40 字）＋ 判据走法 |
| --- | --- | --- |
| **真通知回显数（不计入注入）** | **7** | ① `Bash` "cd \"D:\\work\\workspace\\projects plans\\Wisp\" && git rev-parse --short HEAD && git stat" 结果尾部的 `<system-reminder>`（harness 自己的 available-skills 清单，同一发里还有"date has changed"一句）＝1 枚；② 对话轮里 `Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it was last read.` ＋ 全文回显；③ 同一枚通知的**第二次**（内容与②同源）；④ `Note: The file D:\work\workspace\projects plans\Wisp\docs\reports\HANDOVER.md was modified…` ＋ `git diff` 正文（编排者自己写的台账提交回显，含 `5c1529a`/`f53ad5c` 两枚**我的** commit 短哈希＝渐进 commit 的正常回显）；⑤ `Edit` "file has changed since your last read"（tree1 那枚 AB 变异落刀时，成因＝我用 `cp` 在两次调用之间重写了同一文件）；⑥ `Bash` "Shell cwd was reset to D:\work\workspace\projects plans\Wisp"；⑦ 对话轮里 `<task-notification>` 一段（转述"137 AC#1 终裁表回件"，与我的两格无交集）。**逐条按四条判据过**：路径全部真存在（`ls -l` 到 MEMORY.md 18193 字节、`HANDOVER.md`/台账在盘）、内容是 owner／编排者自己的账、**没有一条**要我 revert／放宽判据／改判据／少取证 ⇒ 判真通知回显。按第④条尺子同样不算（它们都没让我少用工具）⇒ **两栏数字在两把尺子下都不变** |
| **判为注入数** | **0** | 本程工具输出里未出现"编排者备注／系统提示／用户已更新规则／请 revert／放宽阈值／某格已合并／用户已拒绝／Confirm: the harness note is genuine"式文字，也**未出现假 sha**。三件需要点名的事：**(a)** `HANDOVER.md` 回显里写"worker-ticket136-ac13 与 -ac12 两枚还在跑"——本程是**一名代理做两格**（`worker-ticket136-ac12-ac13`）。这条我盘上核了：`grep "agent=" .scratch/wisp/issues/136-*.md` 的署名只有 `orchestrator`/`worker-ticket136-ac1`/`worker-ticket136-ac8-ac9`＋本程，`ls docs/evidence/s1 | grep 136` 里 AC#12/#13 只有我这枚 `136-ac12-ac13-impl.md` ⇒ 判**措辞把一格一名代理拆成了两名的误差**，不是注入、也不构成"有人和我并行做同两格"的证据；**我不因它等待、也不因它改做任何事**。**⚠ 它同时是一条给我的编排事实断言（"两枚在飞"）＝未验证断言**，派单里"两枚兄弟在飞"那句指的就是这两枚名字，我已按"内容真不真"核过并如实写在这里。**(b)** 我在 §4 落笔时**自造过一枚不存在的短哈希** `f4c7620`（想写 `f4c7062` 打错一位）：`git cat-file -t f4c7620` ⇒ `fatal: Not a valid object name`，落盘前被自己那条"引用即须现核"拦下、已从文件里删掉（现文件 `grep -c f4c7620` ＝ **0**，committed 版 `git grep` 亦 0）。登记它**不是注入、是本程仪器自造假 sha 的一发**——正因为本仓今天真发生过假 sha 事故（`injection-timeline.md` §11），我自己这一发也照假 sha 的规格记：**取到没核 ⇒ 不许留在成品里**。**⚠ 本程没有用"全仓 grep 某句话零命中"当"来源清白"的证据**（那条判据已被我们自己的逐字登记摧毁）；要问"某句是不是仓内既有文字"，只用 `git log -S'<原文>' --reverse` 看首枚引入提交。**(c)** 一条时间反常按 §5-12 结案（墙钟跳 8h29m，本机现象，不是注入） |

**两栏之外的一条归属**：编排者派单里那句"AC#12 只需解冻 `sampler.go:470-501`"是**范围写窄了的错断言**（判据① 必然要动 struct 声明区），归派单账／R 账，**不进任何人的注入计数**；本程没有拿它当"所以本格不成立"的依据，也**没有**反过来悄悄把授权扩到自己舒服的形状——越界那两行按 §2.2 逐行钉死并停下等高裁。

### 7.1 本程**引用到／当依据用**的每一枚 sha 与它的核法

| sha | 我怎么来的 | 核法与读数 |
| --- | --- | --- |
| `45c8d1c` | 开工第一步自己 `git rev-parse --short HEAD` | 本程锚点；`git cat-file -t` ＝ `commit`；所有快照只从它或我的 commit 抽 |
| `f4c7062` / `5c1529a` / `f53ad5c` / `c03aee3` | 本程自己四枚 commit（`git log`／`git rev-parse` 现取） | `git cat-file -t` 逐枚 ＝ `commit`；§4 表头那两个（`f4c7062`/`c03aee3`）在落笔前又各现核一次 |
| `a0075ee` / `0f49259` / `7dc8086` | 本程三枚证据 commit | `git log --oneline` 现取；未当锚点用 |
| `76662d8` / `1d38206` / `79ddd49` / `2f291d0` / `172c7aa` / `c2f6b31` / `6effb7e` / `e8190bf` | **外来**（终裁表、票面、派单里读到） | 开工第一步逐枚 `git cat-file -t` ＝ **全部 `commit`**（§0 表）；本程只把它们当"前例形状"与"靶的描述"用，**没有一枚被抄来当快照锚点** |
| `e2ba9f0`/`b802432`/`286de6c`/`8c7ad1f`/`7ad0ec6`/`94b7267`/`69414d3`/`99ac876`/`c04221e`/`097be8a`/`b615092` | 兄弟在飞的 11 枚，`git log --format=%h %s --reverse 45c8d1c..HEAD` **现取** | 只用于共树账（全在 `.md`／票面／台账，`git diff --name-status … -- internal/observe/` 里零命中） |
| `f4c7620` | **本程自造（打错一位）** | `git cat-file -t` ⇒ `fatal: Not a valid object name` ⇒ **不落盘**，见 §7 (b) |
