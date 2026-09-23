# 136 — AC#8（判据②③）＋ AC#9（①②③）：对抗验收（第二格）

- 验收方：`acceptor-ticket136-ac8-ac9-r1`（**非实现者**；实现方＝`worker-ticket136-ac8-ac9`）
- 判据物：票 `.scratch/wisp/issues/136-*.md` 的 **AC#8②③**（按编排者 23:3x 的「措辞更正」那条性质判）与 **AC#9①②③**
- 本程**一个字未改** `internal/observe/**`（只读），仓内未建 worktree、未 checkout；所有变异落在仓外纯净快照树
- 本程**未翻任何 AC 勾**（AC#8／AC#9 的勾由编排者按本表翻）
- 证据文件渐进写：§0／§1／§2… 每裁完一节 commit 一次

---

## §0 锚点、树与基线（开工第一步自己量，不抄派单）

| 项 | 实测 |
| --- | --- |
| 开工第一次 `git rev-parse --short HEAD` | **`76662d8`**（＝本程锚点，所有快照都从它或它的父 commit 抽） |
| 那一刻的 `date -u` 原文 | `Wed Sep 23 15:16:02 UTC 2026` ⇒ 本机 **23:16 +08**（+8 手工换算，没塞进 `date` 的格式串） |
| 分支／开工 `git status --porcelain` | `dev`；工作树只有 1 枚未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`（**未读、未提交、未改、未删、未据它改任何判据**） |
| 实现方自报锚点 | `1d38206` |

### 0.1 `1d38206..76662d8` 在 `internal/observe/` 的差集（有没有"别人动过被测码"）

```
$ git log --oneline 1d38206..76662d8 -- internal/observe/
2f291d0 test(136,AC#9①): 给 CheckSettle 的零可信样本面装钉（与 AC#1 同形，两腿）
79ddd49 test(136,AC#8②): 给 sampler_test.go:308 那枚越界直取加长度守卫（要红不要静默）
$ git diff --stat 1d38206..76662d8 -- internal/observe/
 internal/observe/sampler_settle_zerosample_136_test.go | 153 +++++++++++++++++++++
 internal/observe/sampler_test.go                       |  10 ++
```

⇒ 差集**只有实现方自己那两枚 commit**，没有第三方动过被测面。再往宽处核了一发（派单没要求，但"别人有没有在飞我的包"只能这样问）：

```
$ git diff --name-status 1d38206..76662d8 -- '*.go'      # 全仓 .go 差集
A	internal/observe/sampler_settle_zerosample_136_test.go
M	internal/observe/sampler_test.go
```

⇒ **两枚锚点之间全仓只动了这两枚 .go 文件**，都是本程被验对象自己写的。兄弟在飞的 `cmd/wisp`（票 133）与 `internal/winsec`（票 137）**在本差集里零命中** ⇒ 本程既没读也没被它们影响；两枚 sha 之间交错的都是 `.md`／票面／台账。

被验面的三态归属（自己 `git diff -U0` 复核，不采信报告）：

| 文件 | 1d38206..76662d8 | 结论 |
| --- | --- | --- |
| `internal/observe/sampler.go`（被测生产码） | **无输出** | 本程生产码零改动 ⇒ AC#9④"不许改 SampleState 那侧语义"成立 |
| `internal/observe/sampler_test.go` | `10 增 0 删`，唯一 hunk 在 `@@ -307,0 +308,10 @@`，加的是 `if len(rep.Samples) == 0 { t.Fatalf(...) }` | 只动那枚用例；`t.Fatalf` 不是 `t.Skip`；无删除行 |
| `internal/observe/sampler_zerosample_136_test.go`（AC#1 那枚钉） | **无输出** | AC#1 的断言一字未改（AC#1 已终裁 PASS，本程不重判） |

### 0.2 仓外纯净树（只建不删，全部带本会话后缀 `wisp136r2-ac8ac9-*`）

| 树 | 来源 | 用途 |
| --- | --- | --- |
| `/d/tmp/wisp136r2-ac8ac9-noguard` | `git archive 1d38206` | AC#8① 那一侧的**分母**（无守卫 ⇒ 同一发 M3 今天确实吞读数） |
| `/d/tmp/wisp136r2-ac8ac9-nailless` | `git archive 79ddd49` | AC#9② 的**改前对照**（有守卫、无 settle 钉） |
| `/d/tmp/wisp136r2-ac8ac9-anchor` | `git archive 76662d8` | 本格主读数面（守卫＋两枚新用例都在） |
| `/d/tmp/wisp136r2-ac8ac9-probe` | `git archive 76662d8` ＋ 我自己两枚探针文件 | §3 探哨兵、§4 判 `sampler_test.go:31` 可达性；**探针文件不在仓库里、不在被验版本里** |
| `…-pristine-{noguard,nailless,anchor,probe}` | 各树 `internal/observe/*.go` 的拷贝 | 每发变异**先无条件 restore** 的源 |
| `/d/tmp/wisp136r2-ac8ac9-mut.py` | 本程自写驱动器 | 锚点文本命中数≠期望值 ⇒ 直接退出并回滚（防叠发、防打偏；本程真拦住了一次，见 §2.5） |
| `/d/tmp/wisp136r2-ac8ac9-counts.sh`、`-battery.sh`、`-battery2.sh` | 计数与批量 | 四数只从 `-v` 量 |

工具链写明：`go version go1.27.1 windows/amd64`；`gofumpt` 用**本机既有二进制** `/d/work/base/gopath/bin/gofumpt.exe`（我实测其 `--version` ＝ **`v0.12.0 (go1.27.1)`**，与实现方写明的一致）；**本程没有跑过任何 `go install …@latest`**（派单明令：那会升掉宿主工具链）。容器 `golang:1.27` ＝ `go1.27.1 linux/amd64`（§5）。

### 0.3 基线（未变异，`go test -count=1 -v ./internal/observe/`）

| 树 | rc | RUN | PASS | FAIL | SKIP | panic |
| --- | --- | --- | --- | --- | --- | --- |
| `noguard`（1d38206） | 0 | 56 | 56 | 0 | **0** | 0 |
| `nailless`（79ddd49） | 0 | 56 | 56 | 0 | **0** | 0 |
| `anchor`（76662d8） | 0 | **58** | 58 | 0 | **0** | 0 |

⇒ 58 ＝ 56 ＋ 本程被验的两枚新用例，与实现方 §4.1 的"装钉后 58"逐位相同。

**两笔如实账（不粉饰）：**

1. `noguard` 基线**第一发**就是 `rc=1 / 56 / 55 / 1 / SKIP0`，红名 `TestNoopTaskReturnsToBaseline`、红点 `goroutine_test.go:33: PerTask mid-task = 2, want 3` —— 这是票面 **AC#11** 那枚既有 flake 在**本程锚点上的第三次独立命中**（前两枚是实现方量的）。重取一发才得到上表那行 56/56。它不在本格两格判据上（本格三态都判在"红名是本程仪器"的读数上），但下一位读数的人还会撞到。
2. 我有一发 `go test -count=2 -v` **取在了还没 restore 的变异树上**（读到 52 RUN/47/5＋panic），当场识破、`restore` 后重取（§5 G4 用的是重取那份）。归因：我自己 sequencing 失误，不是被验面的性质。

---

## §1 AC#8 判据②③：逐判据对表（按编排者 23:3x 的「措辞更正」那条性质判）

**本格尺子（更正后的性质）**："**不再由一枚 panic 代答——每枚各红各的、且 `RUN` 名册与基线逐名相同**"。
派单原文那句"同一发变异下只红这一枚、其余全部照常跑完"我**没有当尺子**；红名不止一枚本身不构成退回。
但我按派单要求**独立复核了这条更正的前提**：那 6 枚红名是否每一枚都"本就该红"（见 1.3）。

### 1.1 判据①那一侧我只复跑一次，用途是给②当分母（不是重判 AC#8①）

同一发 **M3**＝`internal/observe/sampler.go:290` 的 `} else if m.PrivateWorkingSetBytes <= 0 {` ⇒ `>= 0`（fixture 足迹全为非负 ⇒ 逢读数都丢 ⇒ `rep.Samples` 处处为空）。
落地先证（`noguard`＝`1d38206`，无守卫那棵树）：

```
$ grep -n 'PrivateWorkingSetBytes >= 0' …-noguard/internal/observe/sampler.go
290:		} else if m.PrivateWorkingSetBytes >= 0 {
$ diff -u <pristine> <tree>/sampler.go      # 只那一行，其余一字未动
-		} else if m.PrivateWorkingSetBytes <= 0 {
+		} else if m.PrivateWorkingSetBytes >= 0 {
$ go build ./...   rc=0
```

读数（原文 `/d/tmp/wisp136r2-ac8ac9-noguard-M3-v.txt`）：

```
rc=1 / RUN=52 / PASS=47 / FAIL=5 / SKIP=0 / panic: 1 处（我的 grep ^panic 命中 2 行 = "panic:" 表头 + 栈里的 panic({…) 帧；真实 panic 一枚）
panic: runtime error: index out of range [0] with length 0 [recovered, repanicked]
  …/internal/observe/sampler_test.go:308        ← panic 源就是那枚无守卫的 [0] 直取
```

被拖走（`=== RUN` 名册差集，基线 56 − 该发 52 ＝ 4 枚，**逐名**）：
`TestLiveRegistryBaselineWithinSleepingGate`、`TestThresholdTableCoversAllStates`、`TestSampleStateZeroSampleWindowFailsClosed`（AC#1 那枚钉）、`TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（AC#1 正向对照腿）。
⇒ 与实现方 §1、以及上一程验收方 §5.1 的读数**逐位相同**（52/47/5＋panic、同一批 4 枚），实现方"改前"那条不是我借来的。

### 1.2 判据②：守卫在位的同一发 M3（树＝本程锚点 76662d8）

```
$ grep -n 'PrivateWorkingSetBytes >= 0' …-anchor/internal/observe/sampler.go
290:		} else if m.PrivateWorkingSetBytes >= 0 {
$ go build ./...   rc=0
$ go test -count=1 -v ./internal/observe/     # 原文 /d/tmp/wisp136r2-ac8ac9-anchor-M3-v.txt
rc=1 / RUN=58 / PASS=52 / FAIL=6 / SKIP=0 / panic=0
```

- **名册差集**：该发 `=== RUN` 名册 vs 未变异基线名册 `diff` **无输出** ⇒ 一枚都没被拖走（改前是 52/56、拖走 4 枚）。
- 那枚守卫自己的读数原文（**是 FAIL 不是 panic，也不是 SKIP**）：

```
=== RUN   TestSamplerGoroutineAccountingFollowsRegistry
    sampler_test.go:315: sampler recorded 0 samples (sample_errors=3 last_error="read returned a zero private working set for a live tree"): there is no reading to check goroutine accounting against
--- FAIL: TestSamplerGoroutineAccountingFollowsRegistry (0.03s)
=== RUN   TestLiveRegistryBaselineWithinSleepingGate      ← 改前被拖走的那枚，现在跑到了
--- PASS: TestLiveRegistryBaselineWithinSleepingGate (0.00s)
```

- **守卫"正常路径下不静默放行"这一条单独验**（派单点名要查）：未变异基线里该用例是 `--- PASS`（不是 SKIP）；`t.Skip` 在本包三枚相关文件里的出现次数＝**0**（`grep -c Skip` 在 `sampler_test.go` 命中 1 次，逐看是注释句"never t.Skip, never a silent return"，两枚 `*_136_test.go` 为 0）；本程**每一发**读数的 `--- SKIP` 计数都是 **0**（下表与 §2、§5 同）。⇒ "空 `Samples` 且守卫存在"这一发必须红、也确实红了，不存在被 `Skip` 洗绿的路径。

### 1.3 更正的前提：6 枚红名逐枚查因果（有没有谁是被误伤／靠 panic 顺序侥幸绿过）

逐枚把断言与被改那一行（`:290` ⇒ `rep.Samples` 恒空）对上，红点原文取自 `…-anchor-M3-v.txt`：

| # | 红名 | 红点原文（截取） | 与 M3 的因果 | 判定 |
| --- | --- | --- | --- | --- |
| 1 | `TestSampleStateAllMetricsAndVerdicts` | `sampler_test.go:78: expected several samples, got 0` | 直接：它就要 ≥3 枚样本 | **该红** |
| 2 | `TestSampleStateSleepingDiskWriteGateFails` | `sampler_test.go:121: disk_write_ops verdict wrong: {… Measured:0 … Pass:true …}` | 样本被丢 ⇒ `WriteOpsTotal` 累加不到 ⇒ 门扇开 | **该红** |
| 3 | `TestSampleStateWorkPeakMemoryIsTargetNotGate` | `sampler_test.go:170: non-gate rows must not fail the report: [… {Metric:sampling Measured:0 valid / 4 errors … Pass:false Gate:true …}]` | 零样本触发了 `:332` 的 `sampling` fail-closed 门行 ⇒ `rep.Pass=false`；断言在 `:169` 就是 `if !rep.Pass` | **该红**（红得对：不可测窗口本就不该绿） |
| 4 | `TestSampleStateLeakFixtureFlipsRed` | `sampler_test.go:190: memory verdict must be red: {tree_private_bytes Measured:0.0MB … Pass:true …}` | 泄漏 fixture 的读数被丢 ⇒ 内存门看不见泄漏 | **该红** |
| 5 | `TestSamplerGoroutineAccountingFollowsRegistry` | `sampler_test.go:315`（新加的守卫，`t.Fatalf`） | **本格要的正面** | **该红** |
| 6 | `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable` | `sampler_zerosample_136_test.go:151: trustworthy reads must be sampled` | AC#1 正向对照腿 | **该红** |

⇒ **6 枚全有因果，无一枚是被连坐的误伤**；也无一枚"靠 panic 顺序侥幸绿过"——改前那发 panic 之前已红的 4 枚（#1–#4）与本发的 #1–#4 同名同点，改后被吞的 2 枚现在各报各的（一枚 PASS、一枚 FAIL）。
⇒ 编排者那条更正的**前提成立**：M3 同时开掉 `sampling` 门扇 ⇒ 另几枚本就该红；本格按更正后的性质判，**不**按"只红一枚"判。实现方读到的 6 枚与我读到的 6 枚**同名同点**。

### 1.4 判据③：AC#1 那枚钉的 M3 复跑，我从**整包 `-v`** 里直接取到（未跑任何 `-run` 定点）

同一次 §1.2 读数的原文（`…-anchor-M3-v.txt`）：

```
=== RUN   TestSampleStateZeroSampleWindowFailsClosed
--- PASS: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
=== RUN   TestSampleStateTrustworthyWindowNotMarkedUnmeasurable
    sampler_zerosample_136_test.go:151: trustworthy reads must be sampled
--- FAIL: TestSampleStateTrustworthyWindowNotMarkedUnmeasurable (0.06s)
```

⇒ 腿 A 绿、腿 B 红在实现方给的 `sampler_zerosample_136_test.go:151`，**颜色与上一程验收方用定点读数取到的形状一致**，但这次整包读数就够。**本程为 AC#8/AC#9 的所有读数都没有用过 `-run`**（唯一用过 `-run` 的是我自己的探针文件，见 §3/§4，与判据读数分开）。

### 1.5 AC#8 三态齐否 ＋ 还原复证

| 态 | 读数 |
| --- | --- |
| 未变异 | `rc=0 / 58 / 58 / 0 / SKIP0 / panic0` |
| M3 | `rc=1 / RUN=58 / 52 / 6 / SKIP0 / panic0`，红名逐名（上表），名册差集为空 |
| 还原 | `diff -q` 与 pristine **无输出**，并另证与本程锚点仓库里那枚 `sampler.go` **逐字相同**；复跑 `rc=0 / 58 / 58 / 0 / SKIP0`（`…-anchor-M3-restore-v.txt`） |

**本节小结**：AC#8②③ 两判据我都在自己的锚点上独立复现（终判与"哪几发复现了"见 §9）。

---

## §2 AC#9 判据①②③：逐判据对表

### 2.1 判据①：用例形状与"是不是自己算自己"

被验文件 `internal/observe/sampler_settle_zerosample_136_test.go`（commit `2f291d0`，`153 增 0 删`，本程锚点在位）。腿 A 的断言链逐跳与其落点（行号我自己 `grep -n` 复核过，实现方给的 `:63/:66/:73/:76/:79/:85/:88/:101-113` **全部对得上**）：

| 跳 | 断言（原文） | 期望值来自哪里 | 落到生产码的哪一处 |
| --- | --- | --- | --- |
| 前提 1 | `reads == 0 ⇒ t.Fatal`（`:62-64`） | **测试自己维护的 fixture 计数器**（`ft.current` 里 `reads++`），不是被测函数的返回物 | `sampler.go:476` 的 `m, err := s.tree.ReadTree()` |
| 前提 2 | `len(rep.Samples) != 0 ⇒ t.Fatalf`（`:65-67`） | 字面量 0 | `sampler.go:477` 那枚 `> 0` 判据 + `:486` 的 append |
| 归因隔离 | `FreeOSMemoryRequested`／`FreeOSMemoryCount>0`／`FinalBytes<=CapBytes`（`:72-80`） | 这三条是**为了让 `Pass=false` 只剩一个可能来源**；其中 `FinalBytes<=CapBytes` 是**报告内部两字段互比**（唯一一处自指，见下） | `sampler.go:499-500` 的 `memOK`／`releaseOK` |
| 钉 | `if rep.Pass { t.Fatalf }`（`:84-86`）、`if rep.BackWithinCapMS >= 0 { t.Fatalf }`（`:87-89`） | 字面量：不许 true、必须仍是负哨兵 | `sampler.go:488-489`（唯一能给 `BackWithinCapMS` 赋非负值的位置）与 `:498-501` |
| 出线自陈 | `wire["samples"]` 在且空、`wire["pass"]==false`、`wire["back_within_cap_ms"]==float64(-1)`（`:101-113`） | 字面量 | `SettleReport` 的 json tag（`sampler.go:431-443`，其中 `samples`/`pass` 在 `:441-442`） |

⇒ **不是"拿被测函数自己的返回值当期望值"**：整条链里只有 `rep.FinalBytes > rep.CapBytes` 那一跳是报告内部两字段互比，而它的用途是**归因守卫**（证明"内存比较这一项本身是满足的"），不是钉本身；钉本身用的是字面量。
补一发说明这一跳也不是哑的：`stateMemCap` 若被改成负数，`memOK` 立刻假红，会先红在 `:78` 而不是悄悄绿过（这一发我**没有真打**，属读码推演，记进 §6）。

腿 B（正向对照 `TestCheckSettleTrustworthyReadsAreRecorded`，`:125-153`）我复算了实现方自加的那发 **M11**（`:486` `rep.Samples = append(rep.Samples, sm)` ⇒ `_ = sm`，落地 `486:			_ = sm` ＋ `go build ./...` rc=0）：

```
rc=1 / RUN=58 / PASS=57 / FAIL=1 / SKIP=0 / panic=0
--- FAIL: TestCheckSettleTrustworthyReadsAreRecorded     红在 sampler_settle_zerosample_136_test.go:142
TestCheckSettleZeroTrustworthySamplesFailsClosed 同发仍 --- PASS
```

⇒ 两腿各响各的，腿 B 不哑（与实现方 §4.4 同读数）。

**再补一发（实现方没造过的支，我造的）V7**：让 settle 循环**压根不读树**（`:476` ⇒ `m, err := TreeMetrics{}, error(nil)`）。这一发打的是腿 A 的**前提腿**，用来排除"前提腿写成恒真"：

```
$ grep -n 'TreeMetrics{}, error(nil)' …-anchor/internal/observe/sampler.go
476:		m, err := TreeMetrics{}, error(nil)
$ go build ./...   rc=0
rc=1 / RUN=58 / PASS=55 / FAIL=3 / SKIP=0 / panic=0
--- FAIL: TestCheckSettleZeroTrustworthySamplesFailsClosed   sampler_settle_zerosample_136_test.go:63: precondition broken: the tree was never read
--- FAIL: TestCheckSettleTrustworthyReadsAreRecorded         sampler_settle_zerosample_136_test.go:142
--- FAIL: TestCheckSettleVerifiesReleaseCounter              sampler_test.go:269
```

⇒ 前提腿**有牙**（今天不响、这么改才响），且这一发同包还有第二、第三枚仪器认（既有的 `TestCheckSettleVerifiesReleaseCounter`），归因不孤立。

### 2.2 判据②：`:477` 落 `> 0` ⇒ `>= 0`（M10）——**两棵独立树各来一次**

**先证落地再读数**（本程锚点树）：

```
$ grep -n 'if err == nil && m.PrivateWorkingSetBytes' …-anchor/internal/observe/sampler.go
477:		if err == nil && m.PrivateWorkingSetBytes >= 0 {
$ diff -u <pristine-anchor>/sampler.go <anchor>/sampler.go      # 只那一行
-		if err == nil && m.PrivateWorkingSetBytes > 0 {
+		if err == nil && m.PrivateWorkingSetBytes >= 0 {
$ go build ./...   rc=0
```

**有钉的树**（`anchor`＝`76662d8`，`/d/tmp/wisp136r2-ac8ac9-anchor-M10-v.txt`）：

```
rc=1 / RUN=58 / PASS=57 / FAIL=1 / SKIP=0 / panic=0
--- FAIL: TestCheckSettleZeroTrustworthySamplesFailsClosed      ← 红名逐名＝只此一枚
    sampler_settle_zerosample_136_test.go:66: precondition broken: a settle window of zero-footprint reads recorded 5 samples,
      report=&{… CapBytes:26214400 FreeOSMemoryCount:1 FreeOSMemoryRequested:true BackWithinCapMS:20 ElapsedMS:100 FinalBytes:0
              Samples:[{TreePrivateBytes:0 …}×5] Pass:true}
=== RUN   TestCheckSettleTrustworthyReadsAreRecorded
--- PASS: TestCheckSettleTrustworthyReadsAreRecorded (0.10s)     ← 同发仍绿
```

**无钉的树**（`nailless`＝`79ddd49`：守卫在位、AC#9 那枚钉还没有；同一发 M10、同一落地证明）：

```
$ grep -n 'if err == nil && m.PrivateWorkingSetBytes' …-nailless/internal/observe/sampler.go
477:		if err == nil && m.PrivateWorkingSetBytes >= 0 {
$ go build ./...   rc=0
rc=0 / RUN=56 / PASS=56 / FAIL=0 / SKIP=0 / panic=0        ← 整包静默通过
```

⇒ **两棵树之差就是本格"有没有牙"的定义**：同一发放宽，无钉的包 56/56 rc=0（复算成立实现方 §4.2 与上一程 `R-136-1` 那句"56/56、rc=0"），有钉的包红且**只红本程那一枚**，红点 `:66`。实现方给的读数（`58/57/1/SKIP0/panic0`、红名一枚、红点 `:66`、改前 `56/56 rc=0`）我**逐项复现一致**。
红名里那串病形与上一程验收方探针同族：`BackWithinCapMS:20`、5 枚 `TreePrivateBytes:0` 的"样本"、`FinalBytes:0`、`Pass:true`（枚数/毫秒随 ticker 抖动）。

### 2.3 判据③：还原 ⇒ 复绿（三态齐）

```
$ python wisp136r2-ac8ac9-mut.py <anchor> <pristine-anchor> restore
RESTORED
$ diff -q <pristine-anchor>/sampler.go <anchor>/sampler.go            # 无输出＝与本程快照逐字相同
$ diff -q <repo@76662d8>/sampler.go <anchor>/sampler.go               # 无输出＝与仓库里的生产码逐字相同
$ go test -count=1 -v ./internal/observe/     # …-anchor-final-v.txt
rc=0 / RUN=58 / PASS=58 / FAIL=0 / SKIP=0 / panic=0
```

`nailless` 那棵树还原后另取一发：`rc=0 / 56 / 56 / 0 / SKIP0 / panic0`。
⇒ **AC#9 三态齐**：未变异绿（58/58）／变异红（58 RUN、1 FAIL、红名点名新用例、红点 `:66`）／还原复绿（58/58）。

### 2.4 判据④的约束面（顺手核，不放宽）

| 约束 | 我的复核 |
| --- | --- |
| 不许改 `SampleState` 那侧语义 | `git diff 1d38206..76662d8 -- internal/observe/sampler.go` **无输出** ⇒ 本程生产码一字未动，两格都是纯加钉 |
| 不许动阈值／golden／`thresholds.go` | 同差集里 `internal/observe/` 只有 `sampler_test.go`（+10）与新用例文件（+153），`thresholds.go` 零命中；D32 的 0.5%/25MB 未被任何一发触碰（我的变异只落 `:290`/`:464`/`:476`/`:486`/`:498-501`/json tag，且**全部在仓外快照树**） |
| 不许把已有 `TestCheckSettle*` 两枚改成 Skip | 两枚在本程每发读数里都是 `--- PASS`（或该红时红），包内 `t.Skip` 出现次数 0；`--- SKIP` 计数每发都是 **0** |

### 2.5 本节仪器自报（诚实）

V7 那一发我的驱动器第一次给的锚点文本在包里命中 **2 次**（`SampleState` 与 `CheckSettle` 同形一行），驱动**拒绝落发并回滚**；那次 `go build` 之后的读数是打在已还原的树上的（58/58 全绿，**不是**变异读数）。把锚点扩成三行唯一后重打，才得到 §2.1 里那发 V7。这条不是被验面的问题，是我的仪器行为，原样记。

---

## §3 那一问：**现有三枚哨兵够不够**（零可信样本仍 pass 且说不出没测到——能不能造出来）

作者 §6.5 报回：本格**没动生产码**，把"报告要说出自己没测到"钉在**现有哨兵**（空 `samples` ＋ `back_within_cap_ms:-1` ＋ `pass=false`）上。派单把"存不存在一个可达的形状让 settle 在零可信样本下仍 pass、而这三枚哨兵全都看不出异常"交给我裁。
**这一支我判得出来，结论在下面，配读数。**

### 3.1 先立结构事实（可复算，不是感觉）

```
$ grep -n 'BackWithinCapMS' internal/observe/sampler.go
438:	BackWithinCapMS int64 `json:"back_within_cap_ms"` // -1 = never within window
464:		BackWithinCapMS:       -1,            ← 唯一的初始化
488:			if m.PrivateWorkingSetBytes <= capBytes && rep.BackWithinCapMS < 0 {
489:				rep.BackWithinCapMS = time.Since(startAt).Milliseconds()   ← 唯一的赋值点
498:	backInTime := rep.BackWithinCapMS >= 0 && rep.BackWithinCapMS <= within.Milliseconds()
501:	rep.Pass = memOK && backInTime && releaseOK
```

`:486`（`rep.Samples = append(…)`）与 `:488-489` 在**同一个 `if` 块**里（`:477` 那个可信判据之内）。⇒ **在生产码未改的前提下**，"把 `-1` 摘掉"与"记下第 1 枚可信样本"是同一个动作，所以：

> `len(Samples)==0` ⇒ `BackWithinCapMS==-1` ⇒ `backInTime=false` ⇒ `Pass=false`。**零可信样本仍 pass 的输入形状不可达**（不是"我没找到"，是写它的两行物理上是一起发生的）。

### 3.2 我自己造的发（探针树，`-v` 原文；这些是我造的支，不是被验版本的判据读数）

| 探到的形状 | 实测报告 | 读法 |
| --- | --- | --- |
| **每一发读树都 `err != nil`**（`fakeTree{err:…}`） | `samples=0 back=-1 final=0 pass=false`，出线 `…"samples":null,"pass":false` | 走的是另一条丢弃路（`err` 那半句），仍 fail-closed ✓；实现方的腿 A 只覆盖了 `err==nil`＋零足迹那一半，我这一半也确实是关着的 |
| **死树报正足迹**（`PIDs:0`＋`PrivateWorkingSetBytes:5MB`） | `samples=5 back=20 final=5242880 pass=true` | 见 3.4 的可达性判定 |
| **整窗口只有 1 枚可信读数**（5 发里第 3 发可信） | `samples=1 back=60 final=5242880 pass=true` | 1 枚就够交差；报告里**看不出**另外 4 枚被丢 |
| **错误与可信交替**（约一半可信） | `samples=2 back=40 pass=true` | 同上：部分未测 = 无痕 |
| **`SettleReport` 的字段清点** | 出线 key ＝ `final_bytes samples pass started_at cap_bytes free_os_memory_count target_state peak_bytes free_os_memory_requested back_within_cap_ms elapsed_ms`；`sample_errors` **ABSENT**、`last_sample_error` **ABSENT**、`dropped_reads` **ABSENT** | settle 侧**没有任何**"丢了多少读数/为什么丢"的字段 |

### 3.3 再看哨兵会不会被改坏：本程造的四发（树＝`anchor`，每发先证 `diff`＋`go build rc=0`，整包 `-v`、零 `-run`）

| 发 | 改法 | 落地 | 整包读数 | 本格那枚钉红不红／红点 |
| --- | --- | --- | --- | --- |
| **V1** | `:464` 初值 `-1` ⇒ `0` | `diff` 单行、build rc=0 | `rc=1 / 58 / 56 / 2 / SKIP0 / panic0`；红名 `TestCheckSettleZeroTrustworthySamplesFailsClosed` ＋ **既有** `TestCheckSettleNeverReachesCap` | **红** `:85`（`settle window with 0 trustworthy samples must never pass, got pass=true`） |
| **V3** | `:477` ⇒ `if err == nil {`（摘掉足迹判据） | `diff` 单行、build rc=0 | `rc=1 / 58 / 57 / 1 / SKIP0 / panic0`，红名只一枚 | **红** `:66`（前提腿：它记了 5 枚它没测到的样本） |
| **V4** | `Samples` 的 json tag `samples` ⇒ `samples,omitempty` | `diff` 单行、build rc=0 | `rc=1 / 58 / 57 / 1 / SKIP0 / panic0`，红名只一枚 | **红** `:102`（`the wire report dropped the samples field entirely`） |
| **V6** | `:498` `backInTime := rep.BackWithinCapMS >= 0 && …` ⇒ `backInTime := rep.FinalBytes <= capBytes` | `diff` 单行、build rc=0 | `rc=1 / 58 / 57 / 1 / SKIP0 / panic0`，红名只一枚 | **红** `:85` |
| ~~V2~~ | `:501` 摘掉 `backInTime` 那半句 | **落地不成立**：`internal\observe\sampler.go:498:2: declared and not used: backInTime`、`go build ./...` **rc=1** ⇒ 该发不是可比较的形状（读数面 `RUN=0`，包没编出来），我按"先证落地再读数"作废，不拿它下结论 | — | — |

⇒ **V6 正是派单要我问的那一发**：它是一枚**编得过**的改法，改完之后 `Samples` 仍为空、`BackWithinCapMS` 仍是 `-1`（两枚哨兵"看不出异常"），但 `Pass=true`。本格这枚钉**仍然红**，因为它不是只盯哨兵，它把 `rep.Pass` 与出线的 `pass` 都用字面量钉了。V1 同理（pass=true 而哨兵齐）。

### 3.4 结论（这一问的答复）

1. **对 AC#9① 那句本身：现有哨兵＋这枚钉够用，不需要为翻本格去动生产码。** 依据：3.1 的结构性耦合（零可信样本 ⇒ `Pass=false` 在生产码未改时不可破）＋ 3.3 里四发能破它的改法**每一发都被这枚钉认出来**（红点各不相同时说明它真的是在逐条看 `pass`／`-1`／`samples`，不是看一个笼统的 rc）。⇒ **不构成本格的退回条件**：我没有造出"零可信样本仍 pass 且这枚钉绿着"的读数。
2. **但有一格真缺口要另立（本格之外）**：settle 侧的"测量诚实度"只有**全丢**这一档会被看见，**部分未测**（1/5、2/4 可信）交出的是一枚看着完全正常的绿报告，而且 `SettleReport` 连 `sample_errors`/`last_sample_error` 都没有——而同包 `StateReport` 早就为这件事加过这两枚字段，注释还写着理由（`internal/observe/sampler.go:165-167`：*"A bare count let an instrument lose samples without saying what it lost（ticket 66：this instrument stops hiding things）"*）。⇒ 这是"同一族病只做了一半"的形状，**另立新格**（我登记为 `R-136-7`，见 §7）。按派单要求我**没有**就地放宽 AC#9 的条件、也**没有**要作者改生产码。
3. **附带一条小的**：`SettleReport.Samples` 为空时 JSON 出线是 `"samples":null`（Go 对 nil slice 的行为），腿 A 那句 `else if !isList && wire["samples"] != nil` 因此**在本格任何形状下都不会响**——真正起作用的是它前面的"key 在不在"（V4 红在 `:102` 就是那条）。这是"防呆支取不到牙"，不影响本格判定，登记为 `R-136-8`（低）。
4. **`PIDs:0` 那一形我判"生产侧不可达"，但只到读码这一层**：`internal/proc/treemetrics_windows.go:71-74`（job 关闭 ⇒ 直接报错）、`:88`（`inTree[r.selfPID] = true` 无条件把主进程塞进树）、`:121`（`m := observe.TreeMetrics{PIDs: len(inTree)}`）⇒ **一次成功的读必然 `PIDs >= 1`**，"有足迹但 `PIDs:0`"在这个 reader 里组不出来；`:104-112` 另有一条"快照看不见自己 ⇒ 重试后 fail-closed 报错"。所以那一发只在 **seam（`fakeTree`）层**存在，不构成生产缺口；我**没有**（也无法在单测里）驱动真 Job Object，故这条写成"读码判定＋探针形状"，不写成读数结论。

---

## §4 `[0]` 直取普查的抽查（每类至少一枚，含最容易只对一半的那类）

### 4.0 我自己的全量普查（不采信实现方的分母）

```
$ grep -rn '\[0\]' internal/observe/*_test.go | wc -l   →  22 行（本程锚点）
```

实现方那张表是 **17 行**，因为它按"索引点"合并同行/相邻行（`got[0]/[1]/[2]` 一行算一枚、`recs[0]` 与 `recs[2]` 的两行算一枚、`outer[0]` 的 241/243 算一枚、`inner[0]/[1]` 的 254/255/257/258 算一枚）。我的 22 行减去 4 处合并，再加上锚点后 `sampler_test.go` 因加守卫而多出来的两行（`:308` 现在是注释、`:318` 才是真读），**两边指的是同一批位置，无遗漏无多出**。

再往宽处补一刀（派单没要求，但"同族"不能只按 `[0]` 这一个字面量问）：

```
$ grep -rn '\[[12]\]\|\[len(.*)-1\]' internal/observe/*_test.go   →  10 行
```

这 10 行**全部落在已普查过的那几枚测试里**（`got[1]/got[2]`、`recs[1]/recs[2]`、`outer[1]`、`inner[1]`、`names[1]`、`recs[len(recs)-1]`），且各自的前置长度守卫就是 §4.1 里那几条 ⇒ 实现方"17 处"作为**普查分母**是完整的，没有第二处 `[0]` 家族被漏。

### 4.1 三类理由各抽一枚以上（逐条看代码，不看他写的理由）

| 抽的哪枚 | 代码（我自己贴的原文） | 判定 |
| --- | --- | --- |
| **类 1：同行 `len(x) != N \|\|` 短路** | `diagnostics_test.go:101  if len(def) != 1 \|\| def[0].ID != "D-test" {`；`sampler_test.go:238  if got := s.Transitions(); len(got) != 1 \|\| got[0] != tr {`；`goroutine_test.go:152  if len(rep.Unknown) == 1 && rep.Unknown[0] == …`（`&&` 形）；`goroutine_test.go:158  if len(rep.Unknown) != 1 \|\| rep.Unknown[0] != …` | **判"不改"成立**：Go 的 `\|\|`/`&&` 短路 ⇒ 长度为 0 时左项就把整式定掉，`[0]` 求值不可达 |
| **类 2：前置 `t.Fatalf` 长度守卫**（抽两枚，含实现方说得最"绕"的那枚） | `logging_test.go:180 if len(names) != 2 { t.Fatalf }` ⇒ `:183 names[0]/names[1]`；`errors_test.go:16 if len(classes) != 17 { t.Fatalf }` ⇒ `:41 classes[0]=`、`:42 AllClasses()[0]` | 第一枚**成立**。第二枚要额外一步：`:42` 换了一次调用（`AllClasses()` 再调一次），我核了实现：`errors.go:39 var allClasses = []ErrorClass{…}`（包级字面量）、`:48-50 func AllClasses() { return slices.Clone(allClasses) }` ⇒ 两次调用同源、无随机源 ⇒ 实现方那句"唯一能让它空掉的写法是 AllClasses() 返回不定长，而那会先在 :16 响"**成立** |
| **类 2 里最微妙的那枚** | `earlylog_130_test.go:229-231 if flushed, _ := b.drain(sink); flushed != 1 { t.Fatalf }` ⇒ `:232 got := sink.snapshot()[0]` | **成立，但理由要说全**：`logging.go:487-509` 里 `flushed` 是"**成功交给 next 的记录数**"，而 sink（`capture130Handler.Enabled` 在 `earlylog_130_test.go:39`、`Handle` 在 `:41-46`）是无条件 append、`Enabled` 按 level 过滤 ⇒ 一枚 fresh sink 上 `flushed == len(snapshot())`；被 level 挡掉的走 `dropped++`，不进 sink。⇒ "`flushed != 1` 先响"确实等价于长度守卫，实现方这句不是糊的 |
| **类 3：`SplitN` 天然非空**（派单点名"最容易只对一半"） | `logging_test.go:248 line := strings.TrimSpace(strings.SplitN(string(data), "\n", 2)[0])` | **这一类只此一枚，且判对了**：Go 的 `SplitN` 对**任何**输入（含空串）至少返回 1 段 ⇒ `[0]` 永在。它还漏答了半步、我补量了：包里另一处分段是 `nobarego_test.go:46 for i, line := range strings.Split(…)` —— **range 不索引**，所以没有第二枚需要判；同类的 `[1]` 形状在本包**不存在**（`SplitN(…, 2)` 的 `[1]` 才会需要对长度）。另外该处拿 `line` 去做 `json.Unmarshal`，失败是 `t.Fatalf` ⇒ 不会静默放行 |

⇒ **抽查结论**：实现方"只改 1 处、其余 16 处不改"的判定我**抽查 4 类 6 枚（含它最微妙的两枚）后签字**，理由分类也站得住；"改一枚证明一枚"的账在 §1 有读数。

### 4.2 它另登记的那枚 `sampler_test.go:31`——**这一支是我自己判的，不是它判的**

```go
// sampler_test.go:23-36
func (f *fakeTree) ReadTree() (TreeMetrics, error) {
	if f.err != nil { return TreeMetrics{}, f.err }
	if f.current != nil { return f.current(), nil }
	if f.i >= len(f.mu) {
		return f.mu[len(f.mu)-1], nil     // :31  ← mu 为空且 current/err 都为 nil ⇒ [-1] 越界
	}
	…
```

- **今天有没有入口踩得到**（复核它的"无"，我自己走一遍）：包内 `&fakeTree{}`（裸构造）只有两处——`sampler_test.go:222` 与 `:233`。`:222` 那处只拿它调 `SampleState(ctx, SLOState("Nope"), …)`，而 `sampler.go:253-255` 的 `if !st.Valid() { return nil, … }` 在**任何一次 ReadTree 之前**就返回；`:233` 那处只调 `MarkTransition`/`Transitions`，压根不读树。⇒ **"当前无入口踩得到"为真**。
- **但它是"一枚新用例就会踩到"的活雷**，我造了发测（探针树 `zz_acceptor136r2_probe2_test.go`，用 `recover` 把 panic 收成读数）：

```
=== RUN   TestAcceptorR2ProbeBareFakeTreePanic
    zz_acceptor136r2_probe2_test.go:18: PROBE2 bare-fakeTree SampleState PANICKED as expected: runtime error: index out of range [-1]
--- PASS: TestAcceptorR2ProbeBareFakeTreePanic (0.00s)
=== RUN   TestAcceptorR2ProbeExistingBareUsesAreSafe
    …:37: PROBE2 MarkTransition on bare fake tree ok (no tree read): b
--- PASS: TestAcceptorR2ProbeExistingBareUsesAreSafe (0.00s)
```

⇒ 判定：**不是本格的债**（AC#8 点名的形状是 `rep.Samples[0]` 那枚、已修；`[0]` 普查也确实不含 `len(mu)-1`），但它是**同族第二枚"吞读数"仪器**，且离被踩只差"下一个人多写一枚用 `&fakeTree{}` 的正常采样用例"。登记为 `R-136-9`（低／可复现／修法方向＝给 seam 自己加守卫，别改任何用例语义）。它自报的"当前无入口"我复算成立，所以不记成"它漏了一处普查"。

---

## §5 CI 落点与门禁（容器原文照贴）

### 5.1 这两枚新用例**在不在** ubuntu 那条腿的分母里（自己 grep，不引前例的行号）

```
$ grep -n 'observe' scripts/portable-tests.sh
140:github.com/CarlosShao/wisp/internal/observe          ← core_pin 名单里
175:        ./internal/memory/... ./internal/observe/... ./internal/secret/...   ← core scope 解析出的包
$ grep -n 'portable-tests' .github/workflows/ci.yml
266:      - name: Portable package tests (core scope; …)   ← step 名
288:        run: bash scripts/portable-tests.sh --scope=core
```

⇒ 实现方与前例引的 `:140`/`:175` **都对**；`ci.yml` 那一步是 `:266` 的 step、`:288` 的 run（票面上我按编排者已更正的口径引，不再指注释行）。该步走 `tools/d22scan/runtests.sh`（SKIP 算_fatal_）⇒ "判为绿"必须是真的跑到。

### 5.2 容器原生跑（`golang:1.27` ＝ `go1.27.1 linux/amd64`，先证挂载真挂上）

命令形状（**Git Bash 下 `MSYS_NO_PATHCONV=1` ＋ `/d/...`**，否则静默空挂载假绿）：

```
$ MSYS_NO_PATHCONV=1 docker run --rm -v /d/tmp/wisp136r2-ac8ac9-anchor:/src \
    -v /d/work/base/gopath/pkg/mod:/gomodcache -w /src \
    -e GOMODCACHE=/gomodcache -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 golang:1.27 bash -c "ls -l /src/go.mod && go test -count=1 -v ./internal/observe/ …"
```

三态原文（**未跳任何一枚**）：

| 发 | 容器内读数 |
| --- | --- |
| 挂载证明 | `-rwxrwxrwx 1 root root 883 Sep 23 15:15 /src/go.mod`、`go version go1.27.1 linux/amd64` |
| 未变异 | `rc=0 / RUN=58 / PASS=58 / FAIL=0 / SKIP=0`，`ok github.com/CarlosShao/wisp/internal/observe 2.448s`；`TestCheckSettleZeroTrustworthySamplesFailsClosed`、`TestCheckSettleTrustworthyReadsAreRecorded` 在 linux 上**各有自己的 `=== RUN`/`--- PASS` 行**（`/tmp/o.txt` 第 80-83 行） |
| M10（`:477` ⇒ `>= 0`） | `rc=1 / RUN=58 / PASS=57 / FAIL=1 / SKIP=0`；`--- FAIL: TestCheckSettleZeroTrustworthySamplesFailsClosed`，红点 `sampler_settle_zerosample_136_test.go:66` |
| M3（`:290` ⇒ `>= 0`） | `rc=1 / RUN=58 / PASS=52 / FAIL=6 / SKIP=0 / PANIC=0`；红名 6 枚与宿主 §1.2 逐名同；守卫红在 `sampler_test.go:315`、AC#1 腿 B 红在 `sampler_zerosample_136_test.go:151` |

⇒ **两格的新仪器在 ubuntu 腿"真进真能红"**（AC#9 那枚拿到 CI 形状的红名＋红点；AC#8 那枚改动在 linux 下同样把 panic 换成了单枚 FAIL、其余 57 枚照跑）。

### 5.3 门禁复核（派单说只看两件事，我都看了；另把票面 AC#7 同款读数一并复算）

| 项 | 我的读数 |
| --- | --- |
| **作者有没有写明 gofumpt 版本** | **写了**（证据 §5-G2：`v0.12.0 (go1.27.1)`，并明说 CI 那步是 `@latest` 未钉 ⇒ 读数只在改版前有效）。我复核：`/d/work/base/gopath/bin/gofumpt.exe --version` ⇒ **`v0.12.0 (go1.27.1)`**，与它写的一致；`gofumpt -l internal/observe/` **无输出**；`gofmt -l internal/observe/` **无输出**。⚠ **我全程没跑 `go install …@latest`**（派单明令：那会升掉宿主 `D:\work\base\gopath\bin\gofumpt.exe`；该二进制 mtime 仍是 `Sep 23 22:23`，早于本程） |
| **`d22scan` `ban #8 internal/` 401→402 的归因** | 两发都在**仓外纯净树**上跑：`noguard`（`1d38206`）⇒ `ban #8 internal/ examined 401`，`anchor`（`76662d8`）⇒ **402**，两形 `rc=0`、`d22scan: clean`。归因：`git diff --name-status 1d38206..76662d8 -- '*.go'` 全仓只有 `A internal/observe/sampler_settle_zerosample_136_test.go` ＋ `M internal/observe/sampler_test.go` ⇒ **＋1 只可能来自它自己那枚新 `_test.go`**，别的 scope 一个没动（`cmd/=39`、`design/=16`、`ban #6/#7` 同数）。`frontend/` 两树都是 **40** ⇒ 实现方那发 `40→43` 确实来自它工作树里未跟踪的 `frontend/dist/` 产物，不是兄弟的树（`git diff --name-only 1d38206..76662d8 -- frontend/` **无输出**） |
| `go vet` 双 GOOS | 原生 `./internal/observe/` **rc=0**；`GOOS=linux ./internal/observe/` **rc=0**；`GOOS=linux go vet ./...` **rc=1**，输出**只一条**诊断：`package …/cmd/wisp → imports …/sherpa-onnx-go-linux: build constraints exclude all Go files`；剔掉 `cmd/wisp` ⇒ 只剩 `cmd/balldebug: build constraints exclude all Go files`；剔掉两枚 `cmd/` 后 **30/30 非 cmd 包零输出 rc=0**（`go list ./...`＝33、剔 cmd＝30，输出文件 0 字节） ⇒ **派单那条"交叉 vet 会停在 cgo"的范围过宽，实现方报回得对，我独立复算同结论：失效面恰好这两枚 cmd 包** |
| `-count=2 -v` 四数与**名册差集** | `noguard`（`1d38206`）＝ `rc=0 / 112 / 112 / 0 / SKIP0`；`anchor` 第一发＝ `rc=1 / 116 / 115 / 1 / SKIP0`，红名 `TestNoopTaskReturnsToBaseline`（`goroutine_test.go:33: PerTask mid-task = 1, want 3`）＝**票面 AC#11 那枚既有 flake 在本程的又一次命中**；第二发＝ `rc=0 / 116 / 116 / 0 / SKIP0`。逐名（`=== RUN` 去重后 `diff`）＝只多两行：`> TestCheckSettleTrustworthyReadsAreRecorded`、`> TestCheckSettleZeroTrustworthySamplesFailsClosed` ⇒ **实现方的 `112→116` 与"逐名只多它这两枚×2 轮"复现成立**，无改名、无消失、无转 SKIP；我第一发与它的 `116/116` 差的那一枚红名是 flake、不是本格仪器，也不在名册差集里 |
| 判"不再红"分清变绿／被跳过 | 本程**每一发** `-v` 读数的 `--- SKIP` 计数都是 **0**（§0.3 三发、§1 三态、§2 三态＋V 族、§5.2 三发、`-count=2` 三发）；被验的三枚文件里 `t.Skip` 出现次数 0（唯一命中 `sampler_test.go:313` 是注释句 "never t.Skip, never a silent return"） |

---

## §6 我没做的档（诚实列，别当已验）

1. **AC#8① 的"独立第二条路"（上一程验收方的 M9：`:297` 的 append 改成 `_ = sample`）我没复算**，只用 M3 一发取分母 ⇒ "两发独立造出同一形状"这条我沿用上一程，没重走。
2. **`stateMemCap` 被改坏那一发我没真打**（派单把阈值面列为禁改，我连快照树里的 `thresholds.go` 都没动）⇒ §2.1 里"归因守卫那一跳也不是哑的"是**读码推演**，不是读数。
3. **没跑 CI 那一步的端到端**（`bash scripts/portable-tests.sh --scope=core` 整条），只在同版本容器里按包跑了 `go test -v ./internal/observe/`；**CI run id 未取**（本程只 commit 未 push）。
4. `-race`、`-shuffle` 未测 ⇒ 所有"名册差集/逐名红"都是**当前执行顺序**下的账（与实现方 §6.6 同一边界）。
5. **`cmd/wisp/**` 一个字未跑**（兄弟在飞）：我只读了 commit 版快照里的 `cmd/wisp/slo_windows.go`，用来确认生产侧确实只认 `Settle.Pass`（`:287` `run.Pass = run.Settle.Pass`、`:323-325` `!run.Pass ⇒ exitCode=1`）——**这是读码，不是 AC#10 的端到端读数**。
6. **linux 全仓两态未跑**：linux 面只有 `internal/observe` 的容器三发＋30 枚非 cmd 包的交叉 `vet`（§5.3）。
7. **`scripts/slo-check.ps1` 那一腿未跑**：只 grep 到它会读 `settle.pass`/`exit_code`（`scripts/slo-check.ps1:41`），没有 PowerShell 形状的红/绿读数。
8. **AC#11 那枚 flake 未做 ≥20 发的复现率测量**（不归本格），本程只贡献 2 个新观测点（§7 末）。
9. §3.3 的 **V2 因编不过而作废**，我没有再找一条"把 `backInTime` 摘掉"的**可编译**等价写法——V6 覆盖了同一危害的可达形状，但严格说不是同一条改法。
10. **AC#1 的 M1/M2/M4/M6/M7/M8/M9 七发未复算**（不是本格账，也没推翻任何结论）。
11. 探针文件（`zz_acceptor136r2_probe_test.go`、`…probe2_test.go`）只存在于我的仓外快照树，**没进仓库、没进被验版本、不提交**。

---

## §7 新账（`R-136x-x`，带严重度／能否复现／修法方向／归谁）

编号先查过占用：`R-136-1..R-136-6` 已被上一程验收方与本票面用掉（`136-ac1-adversarial-acceptance.md:333-338`），**`R-136-7` 起未被任何人用过** ⇒ 本格新账从 7 开始，不撞号。

| id | 严重度 | 现象（我这程的读数出处） | 能否复现 | 修法方向 | 归谁 |
| --- | --- | --- | --- | --- | --- |
| **R-136-7** | **中** | **settle 侧"部分未测"完全无痕**：整窗口只 1 枚可信读数 ⇒ `samples=1 back=60 pass=true`；一半读数 `err!=nil` ⇒ `samples=2 back=40 pass=true`（§3.2 探针原文）。`SettleReport` 的出线 key 里 **`sample_errors`/`last_sample_error`/`dropped_reads` 全部 ABSENT**（字段清点读数）。同包 `StateReport` 早就为这件事加过这两枚字段并写明理由（`sampler.go:164-168`：*"A bare count let an instrument lose samples without saying what it lost（ticket 66：this instrument stops hiding things）"*) ⇒ **同一族病在 settle 侧只做了一半** | **能**：探针文件 `/d/tmp/wisp136r2-ac8ac9-probe/internal/observe/zz_acceptor136r2_probe_test.go`，`go test -count=1 -v -run TestAcceptorR2Probe ./internal/observe/` 即出 | 给 `SettleReport` 补 `sample_errors`/`last_sample_error` 两枚字段，并产出一枚与 `sampling` 同形的自陈门行（要动 `sampler.go:470-501`，**生产码**）；配一发"部分未测"的钉子 | **另立新格**（本票建议 AC#12）。⚠ 我**没有**为此放宽 AC#9 任何条件、也没让作者在本格动生产码——这一格作者交付的东西判得住它自己那句，本账是**下一格**的账 |
| **R-136-8** | 低 | 腿 A 出线检查里 `else if !isList && wire["samples"] != nil` 那一支**永不响**：空 slice 经 `json.Marshal` 出的是 `"samples":null`（探针原文），于是 `isList=false` 且 `==nil`，两支都不进。真正有牙的是它前面那句"key 在不在"——我把 tag 改成 `omitempty` 后红在 `:102`（§3.3 V4） | 能：V4 那发＋探针 `all-reads-error` 的出线 | 二选一：要么把该支改成"必须是长度 0 的 list"（前提是先有 R-136-7 里"生产码把 `Samples` 初始化成 `[]Sample{}`"那半步），要么删掉这支并在注释里写明它防的是哪种形状（**别留一条看起来在守、其实永不生效的支**） | 与 R-136-7 同格（同一枚 `SettleReport` 形状） |
| **R-136-9** | 低 | seam 自己的第二枚"吞读数"仪器：`sampler_test.go:31` 的 `f.mu[len(f.mu)-1]`。今天无入口踩得到（§4.2 复核为真），但**新写一枚用 `&fakeTree{}` 的正常采样用例就立刻 `panic: runtime error: index out of range [-1]`**（探针 `TestAcceptorR2ProbeBareFakeTreePanic` 原文）——那会重演 AC#8 刚修掉的那个形状（一枚 panic 吞掉同包其余几十条读数） | 能：`/d/tmp/wisp136r2-ac8ac9-probe/internal/observe/zz_acceptor136r2_probe2_test.go` | 给 `fakeTree.ReadTree` 自己加守卫（空脚本 ⇒ 显式失败或返回错误，**不许**静默返回零值）；只动测试侧，不碰采样器语义 | 下一位在本包加用例行的人；也可并进 AC#11 那一族（"一条仪器把整包读数变成掷硬币"） |

**不另立新账的一条，但要把 n 交回去**：票面 **AC#11** 那枚既有 flake `TestNoopTaskReturnsToBaseline` 在本程命中 **2 枚**读数——① `noguard` 未变异基线那一发：`goroutine_test.go:33: PerTask mid-task = 2, want 3`（`/d/tmp/wisp136r2-ac8ac9-noguard-baseline-v.txt`）；② `anchor` 的 `-count=2` 第一发：同一红点、值不同 `PerTask mid-task = 1, want 3`（`…-anchor-count2-v.txt`，第二发 `-count=2` 复跑 `116/116/0/0` 全绿，`…-count2b-v.txt`）。本程整包级读数 **26 枚**（22 枚存盘宿主日志＋4 枚容器内 `-v`；另有 1 枚容器发因我自己把 `-mod=mod` 打成 `-mod=dev` 而无效、1 枚宿主发取在未还原的变异树上而作废并重取，都不计）。⇒ 交给 AC#11 的累计观测：实现方 2/27 ＋ 本程 2/26，且**两发的取值不同（2 与 1）**这条对"根因＝计数窗口缺 happens-before"是有用的形状，不是"负载高"能解释的。

---

## §8 两个计数（分栏，不混装；本节写于 `date -u` 原文 `Wed Sep 23 15:53:41 UTC 2026` ⇒ 本机 **23:53 +08**，+8 手工换算）

判据用派单那三条（①点名的路径／对象在本机真不真 ②内容是否削弱 owner 权威或放宽判据 ③声称的动作能否盘上复核），外加编排者在我受派之后落进分支的 `c2f6b31` 里那条新补的第四条（"凡让代理少取证／别用工具／直接给结论的文字，即使不越权也按注入登记并继续取证"）。**两条尺子我都跑了，下面点名差在哪。**

| 栏 | 计数 | 逐条出处（工具名 ＋ 命令／位置前 40 字） |
| --- | --- | --- |
| **真通知回显数（不计入注入）** | **4** | ① `Bash` "cd \"D:\\work\\workspace\\projects plans\\Wisp…\" && git rev-parse --short HEAD && git status…" 结果尾部的 `<system-reminder>`（harness 自己的 available-skills 清单）；② `Bash`（同上一条之后的"Note: The file `C:\Users\swq\.qoder-cn\memory\MEMORY.md` was modified…"，随附它的最新内容）；③ 同一枚通知的**第二次**（内容比②多了两条：第 8 代注入＝假 sha、恒真判据）；④ 一段以 user 轮次出现的 available-skills 清单（与①同源，只是渲染成单独一轮）。②③我按三条判据逐条核过：路径**真**（`ls -l C:/Users/swq/.qoder-cn/memory/MEMORY.md` ⇒ 18193 字节、mtime `Sep 23 23:23`，正落在编排者自己写记忆的时间窗里）、内容是 owner 自己的台账、**没有**要求我改判据／revert／放宽阈值，也没有给我任何"某格已合并"式的动作 ⇒ 判真通知。按第四条尺子它们同样不算注入（没有让我少取证）⇒ **两栏数字在两把尺子下都不变** |
| **判为注入数** | **0** | 本程工具输出里**没有**出现过"编排者备注／系统提示／用户已更新规则／请 revert／放宽阈值／某文件已被修改（针对被验文件）／Confirm: the harness note is genuine"这类文字，也没有出现过**假 sha**：我全程只用了自己 `git rev-parse` 量的 `76662d8` 与 `git log` 里读到的 `1d38206`/`79ddd49`/`2f291d0`，四枚都 `git cat-file -t` ＝ `commit`（作者那一程撞到的 `278d3538…` 那一形**在我这程未重现**，按派单要求原样报回"未重现"，不替它计数也不替它洗掉）。全仓 `grep -rn 'No tools needed' --include=*.go --include=*.md` 我**自己跑了，命中 4 枚文件**（`136-ac8-ac9-impl.md`、本文件、`docs/reports/injection-timeline.md`、`docs/reports/pending-and-issues.md`）——**实现方 §8 那句"零命中"在我这锚点已经不成立**。成因不是有人注入：那句话是**它自己在 `172c7aa`（§7/§8 那枚 commit）里逐字引用了那句话**才进的仓，编排者的两枚台账同样如此。⇒ 归"**引用即须现核**"那族（量的一刻为真、写下结论的那一枚 commit 之后即腐坏），不是新缺陷、不占注入计数，但**下一位照那句去跑会得出相反结论**，故点名。我这程用来自核的同类检查改成只查"我这一程的输出里有没有出现过"，不查仓内字符串 |

**两栏之外必须点名的一条归属**：编排者派单里那条"AC#8② 原措辞『只红这一枚』"是**你的错断言**，实现方按实测报回、你已 append-only 更正并写进派单 ⇒ 按派单纪律它归**误记／R-账**，**不进任何人的注入计数**；我这程是**按更正后的性质**判的（§1），没有拿旧措辞当尺子。
**共树噪声如实一条（不是注入）**：我跑动期间工作树里一直有兄弟在飞的两枚 `.md`（`docs/evidence/s1/133-ac2-r2-acceptance.md`、`…137-ac1-r1-acceptance.md`）处于 ` M` 状态。我每枚 commit 都是 `git add -- <显式路径>` ＋ `git commit -q -F - -- <同一枚路径>`，每次 `git diff --cached --name-only` **只出现我自己那一枚**（本程每枚 commit 逐枚如此，到 §9 落笔时共 10 枚），未发生别人的路径被我带走。

---

## §9 终判（两格分开写；勾由编排者翻，本程一枚未翻）

### AC#8（判据②③）＝**PASS（无附条件）**

尺子＝派单更正后的性质："不再由一枚 panic 代答——每枚各红各的、且 `RUN` 名册与基线逐名相同"。

- **我复现了的**（每发都在我自己锚点的仓外纯净树上、先证落地再读数）：①改前分母 `52/47/5＋panic`、拖走 4 枚逐名（§1.1）；②改后 `58/52/6`、`SKIP=0`、`panic=0`、名册逐名相同、守卫红在 `sampler_test.go:315` 是 **FAIL**（§1.2）；③"更正的前提"我独立复核：6 枚红名**逐枚有因果、无一枚误伤、无一枚靠 panic 顺序侥幸绿过**（§1.3）；④AC#1 那枚钉两腿颜色从**整包** `-v` 直接取到（腿 A `--- PASS`、腿 B 红在 `:151`），本程**零次 `-run` 定点绕过**（§1.4）；⑤三态齐（§1.5）；⑥同一发在 `golang:1.27` 容器里同形（§5.2 M3 那一行）；⑦`[0]` 普查"只改 1 处"＋三类"不改"理由，我抽查 4 类 6 枚后签字，并自建分母证明普查完整（§4.1）。
- **我没做的**：AC#8① 的第二条独立路（上一程 M9）未复算（§6.1）——它不影响②③，且②③的分母我自己复跑过。
- **判 PASS 而不是退回的理由**：本格声称能防的结局（"一枚 panic 代答、同包读数被吞"）在加了守卫的树上**造不出来**（`panic=0` 且名册逐名完整）；同一发在无守卫那棵树上我复现出了它。
- 本格**不产生新账**；同族的下一枚雷在 §7 `R-136-9`（seam 自己那处 `len(mu)-1`：今天无入口，但一发探针证明新用例即踩）。

### AC#9（判据①②③）＝**PASS（无附条件）**

- **我复现了的**：①用例形状对表——腿 A 五跳的期望值全是字面量／测试自己的 fixture 计数器，**不是"拿被测函数返回值当期望值"**（唯一自指那跳是归因守卫 `:78`）；另造 **V7**（循环压根不读树）证明前提腿不是恒真，红在 `:63`（§2.1）；②M10 **两棵独立树各来一次**——有钉：`58/57/1/SKIP0/panic0`、红名**只一枚**、红点 `:66`、腿 B 同发仍绿；无钉（`79ddd49`）同一发：`56/56 rc=0`（§2.2 ⇒ 这一发之差就是"有没有牙"）；③还原：`diff -q` 两处（与快照、与仓库生产码）都逐字相同，复绿 `58/58`（§2.3）；④约束面：`sampler.go` 在 `1d38206..76662d8` **零 hunk**、阈值面零命中、既有 `TestCheckSettle*` 两枚未改未 Skip（§2.4）；⑤M11 复算腿 B 不哑（§2.1）；⑥容器（CI 落点形状）未变异 `58/58`、M10 `58/57/1` 红点 `:66`（§5.2）；⑦门禁两件事：gofumpt 版本写明且我核到二进制同值（**本程没装过任何工具**）、`ban #8 internal/` 401→402 用全仓 `.go` 差集钉死只归它自己那枚新 `_test.go`（§5.3）。
- **那条交给我终裁的问题（三枚哨兵够不够）**：**判得出来，答复是"对本格那句够"**——未改的生产码上"零可信样本仍 pass"结构上不可达（`:486` 与 `:488-489` 在同一块里）；能破它的四种**可编译**改法（V1/V3/V4/V6）**每一发都被这枚钉认出来**且红点各不相同（`:85`/`:66`/`:102`/`:85`），说明它不是只在看笼统的 rc。⇒ **不要求作者为本格改生产码**。真缺口在隔壁一档（**部分未测无痕**）⇒ 另立 **`R-136-7`（中）**，见 §7；我既没有为此放宽 AC#9 的条件，也没有把生产码那半件事塞回本格。
- **我没做的**：`stateMemCap` 被改坏那一发没真打（阈值面禁改，故 §2.1 那句是读码推演）、`slo-check.ps1` 与 `cmd/wisp` 端到端未跑、V2 编不过未另找可编译等价写法（§6.2／§6.5／§6.9）。

### 与实现方读数的对表（一致的与不一致的分开列）

| 项 | 实现方给的 | 我量到的 | 判定 |
| --- | --- | --- | --- |
| 三发基线 56/56、56/56、58/58 | ✓ | 同 | **一致** |
| 改前 M3：`52/47/5＋panic`、拖走 4 枚逐名 | ✓ | 同（同一批 4 枚） | **一致** |
| 改后 M3 红名 6 枚逐名 | `50 PASS/6 FAIL`（树＝`79ddd49`，56 枚） | `52 PASS/6 FAIL`（树＝我的锚点，58 枚） | **一致**（差的就是我树上多出的那两枚 AC#9 用例；红名逐名同） |
| M10 有钉 | `58/57/1/SKIP0/panic0`，红名一枚、红点 `:66` | 同 | **一致** |
| M10 无钉 | `56/56 rc=0` | 同（第二棵独立树） | **一致** |
| M11 腿 B | `58/57/1`，红在 `:142` | 同 | **一致** |
| `-count=2` 逐名差集 | `112→116`，只多本程两枚 | 同 | **一致** |
| `-count=2` 改后四数 | `116/116/0/0` | 第一发 `116/115/1/0`（红名＝既有 flake `TestNoopTaskReturnsToBaseline`），第二发 `116/116/0/0` | **不一致的那枚红名不是本格仪器**，是票面 AC#11；且两次取值不同（`2` 与 `1`） |
| `d22scan` `frontend/` | 工作树 `40→43`（未跟踪 `dist/`） | 两棵纯净树都是 **40** | **不冲突**：我的快照没有未跟踪产物，恰好佐证它的归因 |
| 它 §8 那句"全仓 grep `No tools needed` 零命中" | 0 命中 | 本锚点 **4 枚文件命中** | **它那句已被它自己的 commit 作废**（`172c7aa` 逐字引用了那句话）⇒ 归"引用即须现核"，见 §8，不占注入计数 |

---

## §10 本程 commit 清单（每次暂存只有我自己那一枚路径；只 commit 未 push）

| 节 | sha |
| --- | --- |
| §0 锚点与树 | `9bed3e4` |
| §1 AC#8②③ | `e40924f` |
| §2 AC#9①②③ | `7c2510c` |
| §3 哨兵够不够 | `52551db` |
| §4 `[0]` 普查抽查 | `0e7489a` |
| §5 CI 落点与门禁 | `6f8eced` |
| §6 未做的档 | `7f4a41b` |
| §7 新账 | `5639775` |
| §8 两个计数 | `97037d3` |
| §9 终判＋§10（本枚） | 见 `git log --oneline -1` |

`next=` 编排者按 §9 翻 AC#8／AC#9 两勾；`R-136-7`（settle 侧部分未测无痕，中）建议另立一格并**具名解冻 `sampler.go:470-501` 生产码**；`R-136-8` 与之同格；`R-136-9` 与 AC#11 那一族同姓，可并格。AC#10 仍排 `cmd/wisp` 空出来之后。
