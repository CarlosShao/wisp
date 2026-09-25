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

---

## 4. 确定性对照：把那一枚红**按按钮叫出来**，再看改前后各自怎么答

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 09:16 +0800
$ git log -1 --format='%h %ad' --date=format:'%H:%M'      （本程 §0-§3 已落地为 d8390aa；HEAD 此刻 885f50b＝别家台账 commit）
```

§1.3 已经自陈：本程在净窗里 3800 枚窗口一枚没抓到 ⇒ **"率"这一味本程给不出前后对比**。
这一节给的是能给出的那一味：**同一个物理事件、同一份补丁、两棵树，改前红、改后绿**。

仪器 `arm_stall.py`（**只在仓外快照里存在，绝不进树**；`D:\tmp\observe-ac15-poll-s1\arm_stall.py`）：
把 `alternatingTree.ReadTree` 的**前 N 枚读**各睡 60ms。`CheckSettle` 的环是
"tick → 读 → 比 deadline"（`sampler.go:499-525`），deadline = 起点的 +100ms ⇒ **一枚慢读＝一枚没收到的读**，
窗口带着 `reads=2` 回来——与归档四枚目击逐字同形（`only 2 reads taken`）。
靶子只此一枚 fixture：`gateAlternatingTree` 是另一枚类型，别家腿一枚都不受影响（改后整包 142 枚守恒，见 §6）。

| 树 | 补丁 | 读数 |
|---|---|---|
| `tree-pre-stall`（未改的 `ee5a25e` 快照） | `stallArm(1)`：这一枚窗的第一读睡 60ms | **RED** `sampler_settle_coverage_136_test.go:214: precondition broken: only 2 reads taken, half-and-half needs a window to lose in`／`--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.13s)`／包时 0.166s rc=1（日志 `logs/stall-pre.log`） |
| `tree-post-stall`（本程改后的快照，同一份补丁） | 同上 | **GREEN** `--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.43s)`（日志 `logs/stall-post.log`）。0.43s ≈ 3 枚饿死的窗＋1 枚满窗 ⇒ **等待真的重开了窗，而重开后的那一窗仍然一边keep一边lose**（`kept>=2 / lost>=2 / kept+lost==reads / 门行点名 sample_errors` 那几枚断言全过了，否则不会绿） |

⇒ **改前那一枚红，就是归档 4/8200 里的那一枚红**（同站点、同句子、同 `reads=2`）；
⇒ **改后同一事件不再红，且它不是因为"场景被洗掉"才不红**——洗掉的读法会让下面 §5 的变异一起变绿，
而 §5 五发变异没有一枚变绿。

**救援不是无限的**（约束 (i) 的另一半：有界）：`tree-post-stall-all`（每枚读都睡，永不救援）——

```
sampler_settle_coverage_136_test.go:237: precondition broken: 16 windows opened inside the 2s monotonic
  bound (clock.go Timeout) all fell short of 4; counts seen: [2 2 2 2 2 2 2 2] (+8 more, all of them short of 4)
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (2.10s)      rc=1
```

⇒ 2s 之后**它自己收场并印出判据**：16 枚尝试、每枚各拿到几枚，全在红句里。
不是 `t.Skip`、不是挂死、不是"多试几次直到绿"。同一棵树未加消息截短版跑出的那一发是
`16 windows ... counts seen: [2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2]`（16 项全列）；
把列表截到 8 项是**本程随后追加的可读性修**（§5 之后那枚 commit），只动消息文本，不动判据。

---

## 5. 变异：把生产侧的计数/披露打断，改写腿必须还咬（约束 (iii)）

仪器 `mutate.py`（仓外快照；每枚变异各一棵树 `tree-mut-M1..M5`，全留不删）。
**每一发都是 `-count=1 -v` 整包**（不是 `-run` 摘名册），红名册/站点由 `pair_reds.py` 从日志回读
（Go 把用例日志印在 `--- FAIL:` **之前**，第一版配对器按"FAIL 后一行"读会把站点错配给下一枚用例——
这个自纠记在这里，下面表中的配对是修正后重跑日志得到的）。

| 变异 | 打断的是哪一味 | 整包读数 | 靶子腿与站点（本程改写后的行号） |
|---|---|---|---|
| **M1** | `CheckSettle` 的读错误分支不再 `SampleErrors++` | 71 RUN / **68 PASS / 3 FAIL** / 0 SKIP / 0 panic | `TestCheckSettleHalfTheReadsFailedReportsItsLoss` 红在 `coverage:251`：**"the seam lost 5 of 10 reads but the report says sample_errors=0"** ⇒ 半丢的场景**还在**（10 枚读、5 枚丢），红的是"没披露"。另两枚：`coverage:156`（Single 腿"9 unaccounted"）、`gate:223`（kept=10 lost=0） |
| **M2** | 零足迹分支不再计数 | 71 / **69 / 2** / 0 / 0 | `coverage:332` "zero-footprint drops must be counted: sample_errors=0 reads=10 kept=1" ＋ `gate:275` |
| **M3** | 出线键名改掉（披露从线上消失） | 71 / **67 / 4** / 0 / 0 | `coverage:203` "wire report has no `sample_errors`"、`:296`、`:365`、`:395` 四枚全红 |
| **M4** | 门行把丢了一半的窗判成"测满了" | 71 / **67 / 4** / 0 / 0 | `coverage:183`、**`:283`**（靶子腿："this window dropped 5 of 10 reads and its own row says it measured enough"）、`:353`、`gate:230` |
| **M5** | `CheckSettle` 一整窗一枚读都不交付 | 71 / **61 / 10** / 0 / 0 | **9 枚红逐字是新增的那句 bound 红**（`193-195 windows ... fell short of 1/2/3/4`、`counts seen: [0 0 ...]`），站点：`coverage:143/:237/:321/:374`、`gate:166/:217/:264`、`settle_zerosample:43/:140`；第 10 枚是本程**没包**的那枚 `TestCheckSettleVerifiesReleaseCounter`（`sampler_test.go:303` "settle should pass"），它靠自己的 `!rep.Pass` 断言逮到了同一件事 ⇒ 正好是 §2 末登记的那 4 枚"连前提都没有"的暴露腿之一 |

**这一节要证的三件事，各自落在哪：**

1. **改写腿没有变成常绿**——M1 打断的正是"报告说没丢"那一味，靶子腿红在 `coverage:251`（披露断言），
   且红句里 `5 of 10 reads` 说明**窗口的半丢形状仍在**：等待没有把场景洗成"全收到了"。
2. **等待不能被产品侧坏掉满足**——M1/M2/M3/M4 里等待一律**一次过**（判据是接缝计数器，那几发红句里的
   reads=10/20 就是证据），红的是等待**下面**的断言；只有 M5（接缝真的没读到）才走到 bound 那一支。
   两支是分开的，这正是"轮询不吞产品 bug"的形状。
3. **每枚被改写的腿都还咬得动**——M5 一发把 9 枚新增等待全打成 bound 红＋1 枚未包腿自己红，
   没有一枚在"整窗零读数"这种废掉的产品上保持绿。

对照基线：同一批变异打在**未改**的树上会红几枚，本程**没有量**（见 §7"没测什么"）。
本程只声明必要的一面：改后仍咬。

---

## 6. after：同一命令的四数、名册差集、两口径稳定性批次

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 09:37 +0800
$ git rev-parse HEAD                          （本程落地后的树；快照 pin 记在 D:\tmp\observe-ac15-poll-s1\pin-post2.txt）
dcb95ae...    ← 稳定性批次用的 archive 树；本程自己的 commit 序列见 §8
```

### 6.1 同一命令 `go test ./internal/observe/ -count=2 -v` 的前后四数

| | before（`ee5a25e`，未改，仓外快照 `tree-base`） | after（本程落地后，工作树） |
|---|---|---|
| `^=== RUN` | **142** | **142** |
| `^--- PASS` | **142** | **142** |
| `^--- FAIL` | **0** | **0** |
| `^--- SKIP` | **0** | **0** |
| `^panic:` | **0** | **0** |
| rc | 0 | 0 |
| 不同名数 | **71**（声明数 71） | **71**（声明数 71：`grep -c '^func Test'` 于 15 枚 `_test.go`；**新增的那枚文件零枚用例**） |
| 名册两向差集 | — | `comm -23 golden after = 空`、`comm -13 golden after = 空` ⇒ **名册不缩、不增、不换**（`roster-golden.txt` / `roster-after.txt`，各 71 行） |

日志：`base-count2.v.log`（before）、`after-count2.v.log`（after）。

### 6.2 口径 A（整包单发 `-count=1`，一发＝一进程＝靶子腿的一次执行）

仪器 `batch_a.sh`（本程自写；**逐发串行、无 `&`／nohup／`-parallel`、只名一枚包**，逐发扫 16 枚争用名单
＝`scripts/slo-check.ps1` 那 15 枚 ＋ `Runner.Worker`，扫前先静 3s——第一版没静那 3 秒，
把**自己上一发的 `go.exe` 尾巴**当成了别家负载（20 发里 2 发 `foreign=go:1`），这条自纠记在这儿）：

| 批 | 树 | 发数 | RUN 合计 | PASS | FAIL | SKIP | 真 `^panic:` | 族句命中 | 包时 min/中位/max | 负载标签 |
|---|---|---|---|---|---|---|---|---|---|---|
| A-pre | `ee5a25e` 未改 | **20** | 1420 | 1420 | **0** | 0 | 0 | **0** | 3.447 / 3.476 / 3.560 s | **净窗**（2 发见 `go:1`＝本批自己的尾巴；批前批后九枚名单 0 枚；`gh` in_progress 0） |
| A-post | 本程改后 | **20** | 1420 | 1420 | **0** | 0 | 0 | **0** | 3.458 / 3.475 / 3.580 s | **净窗**（同上；批末 09:37:41 复扫 build-family 合计 0） |

⇒ **绿路零代价**：整包中位 3.476 → 3.475 s（**没量到差别**），逐发字节两批同为 **7481**（日志字节数一模一样 ⇒
新增的等待在绿路上不印任何东西）。
⇒ ⚠ **这两批 40 发不构成"率降了"的证据**：改前那批自己也是 0 命中（本程在净窗里根本抓不到这枚 flake，§1.3），
两批的"0/20"是同一枚读数。派单要求的"before/after 四数"这一栏交的是**名册与形状守恒**，不是率的对比。

### 6.3 口径 B（同进程 `-count=N`，归档 4 枚目击**全部**出自这一口径）

| 批 | 命令形状 | n | 红 | 逐遍耗时 p50/p99/max | 判读 |
|---|---|---|---|---|---|
| B-post-target | `-count=2000 -run '^TestCheckSettleHalfTheReadsFailedReportsItsLoss$'` | **2000 枚靶子腿窗** | **0** | 100 / 100 / 110 ms | 与改前基线（§1.2：100/100/110）**逐档同形** ⇒ 等待在绿路上一次都没重开窗 |
| B-post-all13 | `-count=150 -run '^(13 枚被改写的腿)$'` | **1950 枚腿执行** | **0** | 每枚腿的 p50＝max＝它单窗应有的时长（100/120/200/400/30/50/60ms），**没有任何一枚跳到 2× 档** | ⇒ 1950 枚里**一次重开都没发生**（重开必然多烧一整枚窗） |

⇒ 本程交出的分母（自己量的，不并进一枚百分比）：**口径 A 40 发（20 改前＋20 改后）＋ 口径 B 3950 枚改后腿执行
＋ §1 的 3800 枚改前腿执行**。**两口径禁止加总**（票面 `AC#15` 判据①）。
⇒ **能说什么**：改后在 3950 枚同进程靶子执行里没有一枚红、没有一枚需要重开；等待对绿路的代价量不出差别。
⇒ **不能说什么**：不能说"率从 0.049% 降到 0"。要把 0 压到归档合并率的**上界以下**需要
`n > 3/0.00049 ≈ 6100` 枚同形窗（本程 3950，不够），要分辨 0.25% 与 0.033% 两层的差更要几千枚；
本程对"改后是否更不 flake"的可辩护证据是 §4 那枚**确定性对照**，不是这批计数。

---

## 7. 纪律回执、变异复算、总判、以及"本程没有测什么"

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 09:45 +0800
$ git log --oneline -3
（本程三枚 commit 全带显式 pathspec、每枚只带本程自己的路径；序列见下）
```

### 7.1 §5 那五发变异在**最终代码**上复算了两发

§5 表里的读数是**截短红句之前**那一版跑的（§4 末说明了这枚先后）。M1 与 M5 用最终树
（`tree-post2` @ `dcb95ae` 的副本 `tree-v3-M1`／`tree-v3-M5`）重跑：

| 变异 | 首跑（§5 表） | 最终代码复算 | 差在哪 |
|---|---|---|---|
| M1 | 71 RUN / 68 PASS / 3 FAIL / 0 SKIP / 0 panic | **同值**，站点同为 `coverage:156`／`coverage:251`／`gate:223` | 无差别 |
| M5 | 71 / 61 / 10 / 0 / 0，9 枚 bound 红 | **同值**，站点同一批；尝试枚数从 193-195 变 191-195（**同量级**），红句尾部变成 `counts seen: [0 0 0 0 0 0 0 0] (+183 more, all of them short of 3)` | 只有消息文本被截短，判据与红名册一字未变 |

### 7.2 硬约束回执（`A212④` 那四条 ＋ 派单第 4 条的禁形）

| # | 约束 | 本程 |
|---|---|---|
| (i) | 超时不许用墙钟差 | `awaitWindow` 只用 `clock.go` 的 `NewTimeout/Expired/Remaining`（单调）；**全包 `grep -n '\.Sub(time\.Now())'` 于本程改动的 6 枚文件＝0 命中**，`time.Since` 也只出现在既有生产码里。形状照 `goroutine_test.go:57-66`（AC#11 已裁定的同族修法）。`d22scan` 真扫 clean（§0：正控那发因别家在动而 build failed，**没跑通就是没跑通**，本程不洗） |
| (ii) | 上界由本程现量推 | 2s，推导与倍数写在 `window_wait_136_test.go:44-69` 头上；本程现量＝靶子腿 n=2000 p99 100ms、六腿 n=1800 p99 200ms（§1.2）；倍数＝**最慢腿 10x p99／靶子腿 20x p99** |
| (iii) | 不许把场景洗掉 | §4 确定性对照（改前红＝归档逐字同形；改后绿）＋ §5 五发变异（M1 的红句里仍写着 `5 of 10 reads`＝半丢仍在）＋ 判据一律取**接缝自己的计数器**（产品坏掉满足不了等待，M5 那一发把 9 枚等待全打成 bound 红） |
| (iv) | 禁 Skip／禁改 `< 4`／禁删断言／禁动 `sampler.go`／禁动 golden·阈值·allowlist | 名册 `t.Skip` 执行位 0 枚（注释枚数见 §1.1 的自纠）；阈值 13 枚逐枚原样（§3 的机械 diff 表）；`git diff --name-only ee5a25e -- internal/observe/ \| grep -v _test.go` **空**；`thresholds.go`、golden、`tools/d22scan/allowlist.txt`、`.scratch/**`（含票面勾数 7/8 未变）、`docs/PLAN.md`、`docs/specs/**`、`cmd/**`、`scripts/**` **一枚都没写过**（`git status --porcelain` 逐径核过为空） |
| 派单机制 | 只 commit 不 push／显式 pathspec／禁形清单 | 遵守。每枚 commit 前 `git diff --cached --name-only` 核过；**本程没有一次 `git add -A`/`.`、没有 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`push`**。共享树里 HEAD 一直在漂（`ee5a25e`→`885f50b`→…，别家台账与 d22scan 程在动），所以本程每次都用 `git log -1 -- <本文件>` 单读自己那一枚 |
| 快照 | 只建不删、仓外 | `D:\tmp\observe-ac15-poll-s1\`：`tree-base`(未改)／`tree-post`／`tree-post2`／`tree-post-v2`／`tree-pre-stall`／`tree-post-stall`／`tree-pre-stall-all`／`tree-post-stall-all`／`tree-mut-M1..M5`／`tree-v3-M1`／`tree-v3-M5` ＋ `logs/` ＋ 4 份仪器；**没有一枚被删**，**没有在仓库目录内建 worktree 或 checkout** |

### 7.3 本程 commit 序列

| 段 | commit | 时刻 |
|---|---|---|
| §0-§3（代码落地＋基线＋普查＋修法） | `d8390aa` | 09:13 |
| §4-§5（确定性对照＋五发变异＋红句截短） | `333dfe3` | 09:21 |
| §6（前后四数＋名册差集＋两口径批次） | 见 `git log -1 -- docs/evidence/s1/observe-ac15-poll-r1.md` 于本节之前那枚 | 09:39 |
| §7（本节） | 本枚 | 09:4x |

### 7.4 总判

**改了什么**：`internal/observe/**_test.go` 六枚文件——新增 `window_wait_136_test.go`（一枚有界单调等待仪器，
零用例），13 枚"数了但不等"的前提腿改成 `awaitSettleReads/awaitStateReads`（判据＝接缝自己的读枚数），
`gate_136_test.go` 表头那句"这些腿进不了这一族"按本程实测改成实话，coverage 文件的开窗器合并出一枚
`settleTreeSUT`。**生产码一字未动**；`sampler.go` 一行都没碰（候选 2 属人工批准面，只登记不落地）。

**家族枚数更正**：派单说六枚，本程现量 **13 枚腿／15 处守卫／14 个窗**（Group B 的尾巴没被数进去）；
其中 12 枚已包、1 枚（`SampleState` 侧 `reads == 0`）**结构上不可能由调度造成**故没包并留了说明。
另登记 4 枚"连前提守卫都没有"的同类暴露腿（`TestCheckSettleVerifiesReleaseCounter` 等，§2 末表），
其中一枚在 M5 那一发里真的自己红了——**那一族属下一格，不属本格**。

**变异读数（最终代码复算）**：M1 68/3、M2 69/2、M3 67/4、M4 67/4、M5 61/10 ——
**没有一枚变异把改写的腿洗绿**；靶子腿在 M1 下红在 `coverage:251`，红句照旧印着 `the seam lost 5 of 10 reads`。

**四数前后（同一命令 `go test ./internal/observe/ -count=2 -v`）**：净窗里 before `142/142/0/0`、
after `142/142/0/0`、`^panic:` 前后皆 0、名册 71 枚两向差集为空；整包单发 20＋20 发逐发同值，
包时中位 3.476→3.475 s（绿路代价量不出差别）。

**`AC#15` 本格状态：`[ ]` 未勾**（票面 7 勾／8 未勾未变，`git status --porcelain -- .scratch/` 为空）。
翻勾属**非实现者终裁**，本程不翻。

### 7.5 本程**没有**测什么（按"会被下一位当事实引用"的尺度逐条列）

1. **没有量出"率降了"**：改前 20 发 A ＋ 3800 枚 B 在净窗里 0 命中，改后 20 发 A ＋ 3950 枚 B 也 0 命中。
   两批 0 对 0 不构成比较。要压过归档合并率（0.049%）的上界需要 `n > 6100` 枚同形窗，本程没跑够；
   要压过 1/390 需要 `n > 1169` 枚**整包单发**，本程只有 20。
2. **没有忙窗读数**：派单要"逐批记负载状态"，本程四批全部落在净窗（争用名单逐批 0 枚），
   所以**归档那枚"忙窗 3/1200 vs 净窗 1/7000"的分层，本程无法在改后树上复现**——改后的忙窗行为是**未量**的。
   （唯一例外是 §4 那枚人造停顿，那是**确定性**的饿窗，不是随机负载。）
3. **没在 ubuntu/CI 上量过任何东西**：本程所有读数出自 Windows 本机；`internal/observe` 在 CI 上只被
   `test-core`（ubuntu-latest，`ci.yml:288` → `scripts/portable-tests.sh --scope=core`，包列在 `:175`）问津，
   `--scope=windows`（`ci.yml:458`）**不含本包**。⇒ **承重的分母在 ubuntu，本程给的分母在 windows，两味 OS 不同**；
   定时器粒度（Windows 常见 15.6ms）在 ubuntu 上不同形，改后率要在这台之外另量。
4. **没做门禁全套**：`gofmt`／`go vet` 跑了（净），**`scripts/d22scan.sh` 的正控那一发没跑通**（别家在动的
   `scan_test.go:1197` build failed），本程只补了真扫那一发；`lint`/`golangci-lint`、`go test ./...` 全包、
   `wisp slo`（本程禁跑）一律未跑。
5. **没量"改前遇到同样变异会红几枚"**：§5 五发变异只打在**改后**的树上（要证的是"仍咬"）。
   改前树打 M1-M4 的红名册**没量** ⇒ "改后比改前咬得多"这句本程只能给方向（M5 在改前树上必然只红靶子族、
   不会有 9 枚 bound 红），不能给差值。
6. **没动那 4 枚"没有前提守卫"的暴露腿**，也**没动** `t.Parallel`／`-race`（本包 `t.Parallel` 现量 0 枚，
   改后仍 0 枚；`-race` 本程一枚都没跑）。
7. **没测 `settleWindowWaitBound` 该不该随 OS 变**：2s 是本机现量推的； ubuntu runner 上同一判据是否够，
   未量（见第 3 条）。

### 7.6 一枚枚数自纠（§7.4 的"13 枚改成 await"说满了一枚，原句不抹，正确形状在这里）

现量：`grep -rn 'awaitSettleReads(t \|awaitStateReads(t ' internal/observe/*_test.go \| grep -v window_wait`
⇒ **13 个调用点**，逐点站点为：

```
sampler_settle_coverage_136_test.go:143  :237  :321  :374      （4 枚腿 / 4 个窗）
sampler_settle_gate_136_test.go:166      :217  :264  :269      （3 枚腿 / 4 个窗：RowSeparates 两窗）
sampler_test.go:87                       :334                  （2 枚腿 / 2 个窗）
sampler_zerosample_136_test.go:158                             （1 枚腿 / 1 个窗）
sampler_settle_zerosample_136_test.go:43 :140                  （2 枚腿 / 2 个窗）
```

⇒ 精确形状是 **家族 13 枚腿 / 15 处守卫 / 14 个窗；本程包了 12 枚腿 / 13 个窗**，
少的那一枚腿与那一枚窗是 §2 表里第 13 行 `TestSampleStateZeroSampleWindowFailsClosed`
（`reads == 0` 那枚守卫在 `SampleState` 侧结构上不可能由调度造成，§2 已给理由，原地留了注释、守卫未动）。
§7.4 把"家族 13 枚"与"改了 12 枚"并写成一句，读起来像 13 枚全改了——**那是本程写错的枚数，不是漏了一枚没改**。

### 7.7 §7.6 那张站点表里有三行是本程凭记忆写的、现量不同（原表不抹，逐行更正）

`grep -rnE 'await(Settle|State)(Reads|Samples)\(t, ' internal/observe/*_test.go \| grep -v window_wait`
于本文件当前版本现量 **13 个调用点**（枚数与 §7.6 相同，三行的**行号**不同）：

| §7.6 写的 | 现量 | 备注 |
|---|---|---|
| `coverage:143 :237 :321 :374` | **同** ✅ | 四枚腿 4 窗 |
| `gate:166 :217 :264 :269` | `gate:166 :217 :264 **267**` | `RowSeparates` 的第二枚窗在 `:267`，本程写成 `:269` |
| `sampler_test:87 :334` | `sampler_test:**93** :**335**` | 两行都写早了 |
| `zerosample:158` | `zerosample:**152**` | 写晚了 |
| `settle_zerosample:43 :140` | **同** ✅ | 两枚腿 2 窗 |

⇒ 形状结论不变：**家族 13 枚腿／15 处守卫／14 个窗，本程包 12 枚腿／13 个窗**；
⇒ 但**行号只有这一张表的算本程读过盘**，§7.6 那张按 §4/§5 时的记忆写，其中三行不成立。
这一处属派单"只引自己读过的版本"的同一条规矩，自逮自报。

---

## 8. 交件态复核（在提交后的 HEAD 上重跑，不是引用 §6 的旧日志）

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 09:5x +0800
$ git log -1 --format='%h'                     ← 本程最后一枚证据 commit 之后
$ git diff --name-only <stability-pin>..HEAD -- internal/observe/
(空)   ⇒ §6 那两批稳定性读数跑的就是交件态的码，本程在其后只追加过本文件
```

| 检查 | 读数 |
|---|---|
| `go vet ./internal/observe/` | rc **0** |
| `gofmt -l internal/observe/` | **空** |
| `go test ./internal/observe/ -count=2 -v`（HEAD，工作树） | **142 / 142 / 0 / 0**，`^panic:` **0**，包时 6.885s（`-count=2` 两遍合计） |
| `d22scan` 真扫（`go run . -root .`） | **clean - no D22 ban violations**；`ban #8 internal/ examined 407 Go files, comments and _test.go included`（406→407 是别家新增的 `internal/panel/l2_grant_boundary_test.go`，与本程无关） |
| `scripts/d22scan.sh` 全形 | **仍未跑通**（正控那发被 `tools/d22scan/scan_test.go:1197` 的 build failure 挡住，`set -eu` 使真扫不执行）。这是别家地界的在飞态，本程不修、不洗，只把两味分开报 |
| 票 136 勾数 | **7 勾／8 未勾**，`AC#15` 在 `:348` 仍 `[ ]`；`git status --porcelain -- .scratch/` 空 |
| 本程写过的路径 | 只有 `docs/evidence/s1/observe-ac15-poll-r1.md` ＋ 6 枚 `internal/observe/**_test.go`（每枚 commit 前 `git diff --cached --name-only` 核过，逐枚 `git show --name-only` 可复算） |

### 8.1 负载标签里那一枚含糊，按含糊报

§6.2 那两批逐发扫的争用名单里，A-pre 有 **2 发**、A-post 有 **1 发**印 `foreign=go:1`。
这**不一定**是别家：`go test` 返回后它派生的 `go.exe`/`link.exe` 会多活一瞬间，而扫描发生在
**上一发返回之后、下一发发出之前**那个缝里。两批都各有一发 `foreign=`（空）夹在旁边、包时同带
（3.447-3.580s），所以这 3 发**没有证据是争用**，也没有证据不是。
⇒ 处置：**不剔发、不重跑、不改标签**，按"净窗内、3 发的 `go:1` 身份未定"报。
⇒ 顺带一条给下一位修仪器的人：要把这一味钉死，扫描得带 pid 与启动时刻并按父进程链排除自己，
本程没有那枚仪器，所以没有那个结论。
