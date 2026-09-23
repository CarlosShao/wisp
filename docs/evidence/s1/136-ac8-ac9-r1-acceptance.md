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
2. **但有一格真缺口要另立（本格之外）**：settle 侧的"测量诚实度"只有**全丢**这一档会被看见，**部分未测**（1/5、2/4 可信）交出的是一枚看着完全正常的绿报告，而且 `SettleReport` 连 `sample_errors`/`last_sample_error` 都没有——而同包 `StateReport` 早就为这件事加过这两枚字段，注释还写着理由（`internal/observe/sampler.go:165-167`：*"A bare count let an instrument lose samples without saying what it lost（ticket 66：this instrument stops hiding things）"*）。⇒ 这是"同一族病只做了一半"的形状，**另立新格**（我登记为 `R-136-8`，见 §7）。按派单要求我**没有**就地放宽 AC#9 的条件、也**没有**要作者改生产码。
3. **附带一条小的**：`SettleReport.Samples` 为空时 JSON 出线是 `"samples":null`（Go 对 nil slice 的行为），腿 A 那句 `else if !isList && wire["samples"] != nil` 因此**在本格任何形状下都不会响**——真正起作用的是它前面的"key 在不在"（V4 红在 `:102` 就是那条）。这是"防呆支取不到牙"，不影响本格判定，登记为 `R-136-9`（低）。
4. **`PIDs:0` 那一形我判"生产侧不可达"，但只到读码这一层**：`internal/proc/treemetrics_windows.go:71-74`（job 关闭 ⇒ 直接报错）、`:88`（`inTree[r.selfPID] = true` 无条件把主进程塞进树）、`:121`（`m := observe.TreeMetrics{PIDs: len(inTree)}`）⇒ **一次成功的读必然 `PIDs >= 1`**，"有足迹但 `PIDs:0`"在这个 reader 里组不出来；`:104-112` 另有一条"快照看不见自己 ⇒ 重试后 fail-closed 报错"。所以那一发只在 **seam（`fakeTree`）层**存在，不构成生产缺口；我**没有**（也无法在单测里）驱动真 Job Object，故这条写成"读码判定＋探针形状"，不写成读数结论。
