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
