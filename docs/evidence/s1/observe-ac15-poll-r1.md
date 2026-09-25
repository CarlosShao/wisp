# observe AC#15 - "count without waiting" -> bounded criterion-bearing wait (r1)

agent: `worker-observe-ac15-poll-r1`
ticket face: `.scratch/wisp/issues/136-...md` `AC#15`（`:348-363`，地界句 `:361`）
ruling followed: `docs/reports/pending-and-issues.md` `A212④`（候选 1 ＋ 四条硬约束）
measurements taken from: `docs/evidence/s1/136-ac15-rate-r1.md` §3.2 / §3.3 / §3.4（其中 §3.4 的数按派单要求
**当过期材料处理**，本程自己现量了一版，见 §1.2）
角色: **落地候选 1 并自证它没有把场景洗掉**。本程只动 `internal/observe/**_test.go`；`sampler.go`、
`thresholds.go`、golden、名册、票面勾数一字未动。`AC#15` 那格保持 `[ ]`（翻勾属非实现者终裁）。

临时件：`D:\tmp\observe-ac15-poll-s1\`（只建不删）。

---

## 0. 锚点与仪器

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 09:2x +0800        （本节写于 09:05-09:35 之间，各小节下面各自贴现量时刻）
$ git rev-parse HEAD
ee5a25e689a1f56512fad95f62b58997c3f93ce9     （下面所有 before 读数与 git show 都取这一版）
```

派单简报里"靶子在 `:213/:214`"这一处**在 ee5a25e 上仍成立**（本程现量）：

| 项 | 简报给的 | 本程现量（`ee5a25e`） |
|---|---|---|
| 靶子腿声明 | `coverage:201` | `sampler_settle_coverage_136_test.go:201` ✅ |
| 少判据的那条守卫 | `:213 tree.reads<4` | `:213 if tree.reads < 4 {` ✅ |
| 红句 | `:214 precondition broken: only 2/3 reads taken` | `:214` 逐字 ✅ |
| 影子守卫 | `:216` / `:221` | `:216 if kept < 2 {` / `:221 if lost < 2 {` ✅ |
| `date` 炸弹是否已改过行号 | 简报说 `logging_test.go` 在 `686d7e7` 动过 | 确认 `686d7e7` 只动 `logging_test.go` ＋ 新增 `docs/evidence/s1/observe-datebomb-r1.md`；**`internal/observe` 其余文件行号未漂**（`git diff --numstat ee5a25e 686d7e7^ -- internal/observe/` 只有 logging 一枚） |

仪器与门禁形状（本程自跑，非引用）：

```
$ go test ./internal/observe/ -count=2 -v          （未改动的 HEAD 快照树 D:\tmp\observe-ac15-poll-s1\tree-base）
=== RUN=142  PASS=142  FAIL=0  SKIP=0  ^panic:=0   rc=0     71 枚不同名（声明数 71）
```

⚠ **本程跑测全程在 Windows 上，而 `internal/observe` 在 CI 上只有 `test-core`(ubuntu-latest) 一枚 job 问津**
（`ci.yml:288`；`--scope=windows` 不含本包）。⇒ 承重的分母在 **ubuntu**，本程能给的只有 **windows 本机**的分母，
两味的 OS 不同这一条按派单要求点名写在这里，不当已覆盖。

`scripts/d22scan.sh` 本程**跑不动，且不是本程造成的**：

```
$ sh scripts/d22scan.sh
# github.com/CarlosShao/wisp/tools/d22scan [...].test]
.\scan_test.go:1197:22: not enough arguments in call to ignoredLikeGit
have (string)  want (string, map[string]bool, bool)
FAIL	github.com/CarlosShao/wisp/tools/d22scan [build failed]
d22scan rc=1        （set -eu ⇒ 正控红掉后真扫根本没跑到）
```

`tools/d22scan/**` 是别家正在动的地界（派单已声明），本程不碰、不判它的对错。
**替代读数**（绕过坏掉的正控、只跑真扫那一步）：

```
$ cd tools/d22scan && go run . -root <repo>
d22scan: examined 225 production Go files ...
d22scan: scope ban #8 internal/  examined 406 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations
scan rc=0
```

⇒ 本程新增/改动的 6 枚文件在 **ban #1-5（含 #4 墙钟超时）＋ ban #8（emoji）** 两面上干净；
正控那一发**未跑通**这一事实本程不洗。

---

## 1. before：本程自己现量的基线（全部在 `ee5a25e` 未改动的仓外快照 `tree-base` 里跑）

### 1.1 四数 ＋ 名册

| 项 | 读数 | 怎么量的 |
|---|---|---|
| `^=== RUN` | **142** | `go test ./internal/observe/ -count=2 -v`，日志 `base-count2.v.log` |
| `^--- PASS`（顶层） | **142** | `grep -c '^--- PASS'` |
| `^--- FAIL` | **0** | 同上 |
| `^--- SKIP` | **0** | 同上；本包 `t.Skip` 只在注释里：改前 `ee5a25e` 全包 `grep -rn 't\.Skip' internal/observe/` = **2 枚**（`sampler_test.go:36`/`:335`，两枚都是注释），改后 **3 枚**（本程新文件注释里一枚）。**结构上产不出 SKIP**，不是本程挑出来的读数。（记一笔自纠：本程第一次跑这枚 grep 时接了 `\| grep -v "^.*://"`，那枚过滤器会把**任何带行号的命中**吃掉（`...:36://` 里就有 `://`），第一次因此只数到 1 枚；上表是摘掉坏过滤器重量的。） |
| `^panic:` | **0** | 单独 `grep -c '^panic:'`，不与 "Panic" 字样混（台账 A182 那条口径） |
| 名册 | **71 枚不同名**，`-count=2` ⇒ 142 RUN = 71×2 闭合 | `grep '^--- PASS' \| awk '{print $3}' \| sort -u` → `roster-golden.txt`（71 行） |
| 声明数 | `grep -c '^func Test'` 14 枚 `_test.go` 合计 **71** | 与名册同值 |
| 机器负载 | **净窗**：`go/compile/cgo/asm/link/gcc/Runner.Worker/wisp/staticcheck` 九枚名单现扫 **0 枚** | 跑测前后各扫一次（`Get-Process` 名单式扫描） |

### 1.2 上界算术的原料（本程现量，派单约束 (ii)）

单遍窗口耗时分布，从 `go test -v` 自己印的 `(0.10s)` 字段取（不引入任何计时仪器）：

| 窗口（腿） | n | p50 | p90 | p99 | max | 红 |
|---|---|---|---|---|---|---|
| `TestCheckSettleHalfTheReadsFailedReportsItsLoss` 单名连发 | **2000** | 100ms | 100ms | **100ms** | **110ms** | **0** |
| 六枚 Group A 腿混跑（同进程交替） | 300×6=**1800** | 100ms | 200ms | **200ms** | **200ms** | **0** |

原始日志 `base-B-target-2000.log`（包时 201.298s）、`base-B-family-300.log`（08:58:18→09:01:56），
读数器 `read_durations.py`（本程自写；按**顶层判定行**取时长、按**句子**计命中，两向都不按名筛）。
负载状态：两批之间与之前之后各扫一次九枚争用名单，**全 0 ⇒ 净窗**。

⇒ 本程自己量的数里没有 §3.4 那枚 0.88s 长尾（1200+2000+1800 = 本程 3800 枚 / 归档 1200 枚里各一枚）。
**上界推导（写进代码注释的就是这段）**：

- 家族里最慢的腿 p99 = 200ms（gate 那两枚 200ms 窗），靶子腿 p99 = 100ms；
- 取 **2s = 最慢腿 p99 的 10 倍 = 靶子腿 p99 的 20 倍**；
- 代价上限：中位尝试 ≈ 20 次；即使每次都被拖到归档见过的那枚 880ms 长尾，也还放得下 **2 次完整窗口**
  ⇒ 一枚慢窗**自己吃不掉**预算，预算只有在"连续多窗都被饿死"时才可能到期；
- 那个连串的算术：按归档合并率 0.049% 连失 20 窗 = 1e-56；按它最坏分层 0.25% 连失 20 窗 = **1e-52**；
- 到期是**响亮的一发红**（印尝试次数＋每窗实际读到的枚数），不是 Skip、不是降级 ⇒ 残余风险只是
  "机器连续卡满 2 秒时多一枚可读的红"。

⚠ 本程**没有**推翻"上界取多少才对"的那枚前提，只是把它从归档的数换成本程现量的数：
`100/200ms` 两枚 p99 是本程的，`880ms` 长尾与 `0.049%/0.25%` 两枚率**仍是归档的**（本程 3800 枚没抓到一枚，见 §1.3）。

### 1.3 一条必须自陈的取数失败：本程没有在自己手上复现出这枚 flake

本程在 `ee5a25e` 上取的靶子腿分母是 **2000 + 300 = 2300 枚靶子腿窗口 ＋ 另 1500 枚家族窗口**，命中 **0**。
按归档合并率 0.049% 推，3800 枚里期望 1.85 枚、一枚抓不到的概率 `e^-1.85 = 16%` ⇒ **本程的 0 命中与归档不冲突，
但它本身不构成"率变了"的证据**，更不构成"已修好"的证据。要分辨 0.25%（忙窗分层）与 0.033%（净窗分层）需要
再投几千枚，本程没投（代价见 §6）。

⇒ 因此本程对"修法有效"的证明**不建在率上**，建在两味可判的东西上：
①**确定性对照**（§4：人为造一枚饿死的窗，改前必红、改后必绿）；②**变异**（§5：打断生产计数/披露那一面，
改写腿必须转红并点名站点）。这两味本程都跑出了读数。

---

## 2. 家族普查（本程现量；简报那句"六枚"要更正）

判据：**一条前提在窗口刚关时被读、而它数的东西只可能由 ticker 交付产生**（⇒ 少一枚判据、多一枚 race）。

`ee5a25e` 现量结果：**13 枚腿 / 15 处守卫 / 14 个窗口**，不是六枚。六枚是归档把 Group A 塌成名册的数
（`136-ac15-rate-r1.md` §0.2 那张表），**Group B 的尾巴没数进去**：

| # | 腿 | 守卫（改前行号） | 判据（改后 await 的 want） | 归档分组 | 本程 |
|---|---|---|---|---|---|
| 1 | `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | `coverage:125 reads<3` | reads≥3 | A | 已包 |
| 2 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss` | `coverage:213 tree.reads<4`（＋影子 `:216/:221`） | reads≥4 | A **靶子** | 已包 |
| 3 | `TestCheckSettleZeroFootprintDropsAreCountedToo` | `coverage:293 reads<3` | reads≥3 | A | 已包 |
| 4 | `TestCheckSettleFullyMeasuredWindowReportsNoLoss` | `coverage:343 reads<3` | reads≥3 | A | 已包 |
| 5 | `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` | `gate:201 kept<1\|\|lost<1` | reads≥2 | A | 已包 |
| 6 | `TestSampleStateAllMetricsAndVerdicts` | `sampler_test:99 samples<3` | reads≥5 | A | 已包 |
| 7 | `TestSettleCoverageRowExistsAndPassesWhenFullyMeasured` | `gate:156 reads<1\|\|…` | reads≥1 | **B** | 已包 |
| 8 | `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | `gate:239 reads<1\|\|…` | 两窗各 reads≥1 | **B** | 已包（含 `:237` 那枚 `_` 丢弃计数的窗） |
| 9 | `TestSamplerGoroutineAccountingFollowsRegistry` | `sampler_test:336 samples==0` | reads≥3 | **B** | 已包 |
| 10 | `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable` | `zerosample:150 samples==0` | reads≥3 | **B** | 已包 |
| 11 | `TestCheckSettleZeroTrustworthySamplesFailsClosed` | `settle_zerosample:62 reads==0` | reads≥1 | **B** | 已包 |
| 12 | `TestCheckSettleTrustworthyReadsAreRecorded` | `settle_zerosample:141 samples==0` | reads≥1 | **B** | 已包 |
| 13 | `TestSampleStateZeroSampleWindowFailsClosed` | `zerosample:63 reads==0` | — | **B** | **没包，理由在下面** |

**第 13 枚为什么不包**（这一条是判断，不是遗漏）：`SampleState` 在 ticker 环之外**必读两枚**
（`sampler.go:269` 开窗读 ＋ `:306` 收窗读），所以 `reads == 0` 这枚守卫在 `SampleState` 家族里
**结构上不可能由调度造成**——给它套一层"等它非零"的等待，等的是一枚永远不会到的救援。
同一条理由不适用在第 11 枚上（`CheckSettle` 环外**零枚**前置读，`sampler.go:493-501`，那枚守卫真能红）。
这一处已在原地留注释说明（`sampler_zerosample_136_test.go` 那枚守卫上方），守卫本身一字未动。

**普查顺带量到、但不在本格射程的另一形（登记，未动）**：另有 **4 枚腿**的断言**同样**会被饿死的窗打翻，
而它们**连一枚前提守卫都没有**，红句会长成产品 bug 的样子而不是 `precondition broken` 的样子：

| 腿 | 暴露的断言 | 饿死时的红 |
|---|---|---|
| `TestCheckSettleVerifiesReleaseCounter`（`sampler_test.go:290`） | `if !rep.Pass` | 0 枚 tick ⇒ `BackWithinCapMS=-1` ⇒ "settle should pass" 红 |
| `TestSampleStateWorkPeakMemoryIsTargetNotGate` | `if !rep.Pass` | 0 枚样本 ⇒ `sampler.go:334-340` 那枚 `sampling` 门行 veto ⇒ 红 |
| `TestSampleStateSleepingDiskWriteGateFails` | `v.Metric=="disk_write_ops" && v.Pass` | 样本为 0 ⇒ `WriteOpsTotal=0` ⇒ 该门行**过** ⇒ 红 |
| `TestSampleStateSleepingTCPGateFails` | 同形（`tcp_connections`） | 同形 |

⇒ 本格修的是"**有前提守卫却少判据**"那一族（AC#15 判据②③ 字面射程）；这 4 枚是"**连前提都没写**"，
给它们加判据＝**新增断言**，超出候选 1 的射程，本程只做登记并交回。归档四枚目击逐字都是 `reads=2/3`
（`136-ac15-rate-r1.md` §3.1），**一枚 0 窗都没出现过**，所以本程没有拿"理论暴露"当动手理由。

---

## 3. 修法：一枚仪器 ＋ 13 处形状相同的等待

新增 `internal/observe/window_wait_136_test.go`（**只有仪器、零枚用例** ⇒ 名册不动）：

- `settleWindowWaitBound = 2 * time.Second`，推导写在它头上（§1.2 的数）；
- `awaitWindow[R any](t, want, open, short)`：`open()` 每次开**一整枚新窗**（新 fixture、新 report、新计数器），
  `short(w)` 是这枚窗自己的枚数；到判据就返回那枚窗，没到就重来；预算用尽 ⇒ `t.Fatalf` 印尝试次数＋每窗枚数；
- `awaitSettleReads` / `awaitStateReads` 两枚薄封装，判据一律是**接缝侧的读枚数**。

复用与新增的交代（派单那两问的第一问之外、`A212` 判放水两问的第二问）：

| 问 | 答 |
|---|---|
| 有没有既有仪器能干这件事？ | **等待**这一味没有：改前两枚文件各跑一次 `git show ee5a25e:<f> \| grep -cE 'attempts\|retries\|Poll\|Timeout\|time.Sleep\|for i := 0'` ⇒ **0 / 0**（归档 §3.2 用另一枚式子量过同一条）。本包里唯一既有的有界等待形状是 `goroutine_test.go:57-66` 的**内联** `NewTimeout` ＋ `for cond && !tm.Expired()`，它既不能被复用（它等的是 channel，不是计数）、也没被抽出来（抽它＝动 AC#11 那枚已结案腿的文件），所以本程**照它的形状**写、不复用它的码 |
| 新造的仪器是不是"为了能过而造"？ | 不是。它**只**做三件事：再开一窗、比较一个判据、到期时把红印得比原来更详细。它不判产品侧任何字段（判据是 fixture 自己的计数器）、不吞任何条件（到期是 `t.Fatalf`）、不含 `Skip`／`Sleep`／墙钟差 |
| `open` 侧复用了什么？ | 三枚既有开窗器一字未改语义：`settleSUT`（现改为委派 `settleTreeSUT`）、`gateSUT`、各腿自建的 `NewSampler(...).SampleState/CheckSettle`。`settleTreeSUT` 是本程新增的**开窗器**（不是判据器）：靶子腿用的是有状态 `alternatingTree` 而 `settleSUT` 只收脚本；`gateSUT` 本来就能收 `TreeReader`，但 AC#14 那枚文件的表头 `:29-31` 明写"AC#15 管 coverage 那枚文件的偶发账，本文件自带 fixture、互不借用"，所以本程没有跨文件拿它用 |

**断言有没有被改？（判放水第一问，机械读数）**

把 5 枚被改文件里所有以 `if <条件>` 开头的判定行改前/改后各抽一份、去掉行号后逐行 diff：

```
sampler_test.go                     before=15  after=15  diff: (空)
sampler_settle_coverage_136_test.go before=30  after=30  diff: 两枚
      9c9   < if tree.reads < 4 {                 > if reads < 4 {
      12c12 < if kept+lost != tree.reads {        > if kept+lost != reads {
sampler_settle_gate_136_test.go     before=20  after=20  diff: (空)
sampler_settle_zerosample_136_test  before=5   after=5   diff: (空)
sampler_zerosample_136_test.go      before=11  after=11  diff: (空)
```

⇒ **零枚条件被改**：阈值 `< 4`／`< 3`／`< 2`／`== 0`／`!= reads` 逐枚原样；两枚差异是**同一个量的名字**
（`tree.reads` → `reads`，值不变，只是"这一窗的读枚数"现在由等待返回的那枚窗带着，见下面的窗口-账对齐）。
`t.Fatal*` 语句枚数全包 **375 → 375 守恒**；coverage 单文件 49→48 ＋ window_wait 0→1，那一枚的去处是
`settleSUT` 与靶子腿各自重复的 `if err != nil { t.Fatal(err) }` 合并进了 `settleTreeSUT`（每条开窗路径仍各查一枚 err）。

**等待会不会把场景洗掉？（判放水真正要防的那一味）**

不会，且这一条不是"我相信它"：判据是 `settleWindow.reads`，它是 **fakeTree/alternatingTree 自己的计数器**；
`rep.SampleErrors`／`rep.Samples`／门行**都不是**判据。⇒ 生产侧计数/披露坏掉时等待**照样一次过**
（因为读数是接缝侧真发生的），随后 `lost < 2`、`kept+lost != reads`、`coverage.Pass`、`note` 点名那几枚断言
自己会红——§5 的变异就是去打这几枚，读数支持这句话。

**窗口-账对齐**（这一处是本轮改动里唯一有语义形状的地方，值得单说）：靶子腿改前是
"建一枚 `tree` → 开一窗 → 读 `tree.reads`"；改后是"每次尝试建**一枚新** `tree` → 开一窗 → 带回那一窗的
`rep` 与 `reads`"。**没有**把多次尝试的计数累加后再去比单窗的 `rep.Samples`——那会让
`kept+lost == reads` 这枚承重等式变成跨窗假账（也就会把本腿的题洗掉）。每一枚被包的腿同此。
