# 136 AC#15 - first real hit-rate measurement (r1)

agent: `worker-ticket136-ac15-rate-r1`
date: 2026-09-24（本文件按阶段渐进写、渐进 commit；每节开头贴本程自量的 `date` 与 `git log -1`）
角色: **量一次命中率 + 按需归因**。本程不改任何产品码、不改任何断言／阈值／golden、不改
`scripts/**` 与 `tools/**`；只新建／追加本文件与 `/d/tmp/wisp136ac15m-*` 前缀的临时件（只建不删）。
票 136 `AC#15` 的门禁保持 `[ ]`（本程不翻勾）。

配方以文件为准（先读再跑），本程按这两张表跑：
`docs/evidence/s1/136-ac15-denominator-census-r1.md`（下称**分母表**）·
`docs/evidence/s1/136-ac15-target-census-r2.md`（下称**靶子表**）。

---

## 0. 阶段 0：锚点、名册按名复核、仪器与门禁形状

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 23:34 +0800
$ git log -1 --format='%h %ad' --date=format:'%H:%M'
8c1d52e 23:34
```

### 0.1 快照锚点与树身份

| 项 | 值 | 怎么量的 |
|---|---|---|
| 快照（跑测用的树）| `beac693ebabf103fa9a539e3012196c608ea08e3` | `git rev-parse HEAD` 于 `git archive` 之前，逐字写进 `/d/tmp/wisp136ac15m-pin.txt` |
| 快照树目录 | `/d/tmp/wisp136ac15m-tree-beac693/` | `git archive --format=tar $PIN \| tar -x -C …`（**仓外**，绝不跑在 work tree 上：分母表 §3.1） |
| 一枚先前未钉号的同名树 | `/d/tmp/wisp136ac15m-tree/` | 第一次 archive 时 HEAD 正在漂（`7175ce0`→`beac693`），故弃用改钉号版；两棵树的 `internal/observe/` `diff -r` **完全一致**，故两版读数可互换；按"只建不删"没删它 |
| 包身份（24 枚 `.go` 的 md5 清单） | `/d/tmp/wisp136ac15m-tree-md5.txt` | `find internal/observe -name '*.go' \| sort \| xargs md5sum` |
| 靶文件 md5 | `sampler_settle_coverage_136_test.go` = `a31968022574b99b7969d0ccddb27eb0` | 与分母表 §3.5 记的 HEAD 值**逐字相同** ⇒ 我这棵树与它的树是同一版靶子 |
| 生产码 md5 | `sampler.go` = `3aa170fc15cf33a121e32345babf0a30` | 同上；本程全程不动它 |
| `go version` | `go1.27.1 windows/amd64` | 与仪器修复表 `136-instr-fixes-r1.md §0` 同版 |

### 0.2 简报给的等式基线（`RUN=71`）在盘上成立，且行号未漂移——按名复核，不背行号

简报提醒"表里的行号是它当时锚点下的行号、盘上可能漂移 ⇒ 用名字定位"。本程**逐名复核**，结论：

**六枚能产出"窗口没读满"这一族的腿（靶子表 §1(c) Group A 的 8 枚守卫塌成 6 枚名字）**——名字→声明行→守卫行→红句行，全部在盘上现量：

| 名字 | 声明于 | 前提腿守卫 | 红句 |
|---|---|---|---|
| `TestCheckSettleHalfTheReadsFailedReportsItsLoss` | `sampler_settle_coverage_136_test.go:201` | `:213 if tree.reads < 4 {` | `:214 precondition broken: only %d reads taken, half-and-half needs a window to lose in` |
| `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | 同文件 `:123` | `:125 if reads < 3 {` | `:126 the window only took %d reads, it must lose some` |
| `TestCheckSettleZeroFootprintDropsAreCountedToo` | 同文件 `:291` | `:293 if reads < 3 {` | `:294 the window only took %d reads` |
| `TestCheckSettleFullyMeasuredWindowReportsNoLoss` | 同文件 `:341` | `:343 if reads < 3 {` | `:344 only %d reads taken` |
| `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` | `sampler_settle_gate_136_test.go:198` | `:201 if kept < 1 \|\| lost < 1 {` | `:202 an alternating tree must both keep and lose reads within the window` |
| `TestSampleStateAllMetricsAndVerdicts` | `sampler_test.go:84` | `:99 if len(rep.Samples) < 3 {` | `:100 expected several samples, got %d`（**不含 `precondition broken` 字样**，故我的族过滤器按句子集合点名它，见 §0.4） |

⇒ **分母表 §4.4 那张名册的行号一枚都没漂**（`coverage:201/123/291/341`、`gate:198/154/236`、`sampler_test.go:84` 逐条对上）。
⇒ 靶子表 §1(a) 写的 `:213-215` 也对上，票面原句 `:189` 确实是旧地址（票面 `>` 块已自行更正过）。
唯一漂移：**仪器缺陷那一处**——靶子表/票面点名的 `sampler_settle_gate_136_test.go:262` 裸下标，现在盘上是
`:271` 带守卫、红句在 `:272`（`136-instr-fixes-r1.md §3` 已交付：守卫 `len(neverRows) != 1`）。
即**本程跑的树里那枚"panic 吞名册"的隐患已经被修掉了**，这与我下面 §0.3 的名册守恒读数相互印证。

八枚见证腿（分母表 §4.4 第二张表）同样按名复核，全部存在：
`TestNoopTaskReturnsToBaseline`(goroutine_test.go:25)、`TestSampleStateZeroSampleWindowFailsClosed`(sampler_zerosample_136_test.go:49)、
`TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`(同文件 :141)、`TestCheckSettleZeroTrustworthySamplesFailsClosed`(sampler_settle_zerosample_136_test.go:37)、
`TestCheckSettleTrustworthyReadsAreRecorded`(同文件 :125)、`TestSamplerGoroutineAccountingFollowsRegistry`(sampler_test.go:319)、
`TestSettleCoverageRowExistsAndPassesWhenFullyMeasured`(gate:154)、`TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured`(gate:236)。

同族口径两枚也现跑对上靶子表 §3：窄口径 **8 站点／3 枚文件**、`grep -c '^func Test'` = **71**（14 枚文件），
`t.Parallel` = **0**，`t.Skip` 只出现在**注释**里两枚（`sampler_test.go:36`、`:335`）⇒ **`SKIP=0` 是结构上必然**，
不是我挑出来的读数。

### 0.3 简报的两枚前提，一枚在盘上复现、一枚要改写措辞

* ✅〔独立复现〕"仪器现态 `RUN=142 PASS=142 FAIL=0 SKIP=0 真 panic=0`、rc=0"：这正是 `136-instr-fixes-r1.md §0/§4.1`
  的 `-count=2 -v` 读数（142 = 71×2、顶层 71 枚各两遍）。与简报给的等式基线 `RUN=71`（单发）**同一棵树、同一口径**，
  差异只是 `-count` 参数。⇒ 简报没有把枚数说错。
* ⚠〔要改写〕简报说"之前那批 30 发命中 0 已被裁为分母太小"：盘上我**没找到**一枚属于 AC#15 名下的 30 发整包批次
  （`find /d/tmp -name '*.v.log' -newermt '2026-09-24 19:00'` ⇒ 0 命中；19:00 之后 `/d/tmp` 新增的是
  `141-acc-*`、`wisp136instr-r1`、`ac14b-r2-136` 等别家的件）。`3/30=10%` 那枚上界在盘上的**原始出处是分母表 §3.4**，
  它讲的是"如果只跑 30 发会怎样"的算术，配的历史批次是 AC#11 那 390 发（内含 30 发一批的形状）。
  ⇒ 结论不受影响（30 发确实不够），但**"那批 30 发"作为一次已经跑过的读数我在盘上找不到**，按本仓口径登记为
  〔简报转述、盘上无件〕，不当我这一程的前提用。我这程的分母是**我自己跑出来的 n**。

### 0.4 本程仪器（两份，都在 `/d/tmp`，都从 AC#11 的归档件抄形再改）

| 件 | 来源与差异 |
|---|---|
| `/d/tmp/wisp136ac15m-batch.sh` | 抄 `/d/tmp/wisp136ac11v-batch.sh`（分母表 §3.1 判"可逐字复用"的那一型：有界争用闸门 60 轮×30s、超带 `exit 8`、`gate.txt` 记树身份、`gate-post.txt` 复扫、逐发 `BATCH-TAG`/`RUN-DATE` 头与 `RC=` 尾）。**三处刻意的差异，逐条报名**：① 闸门第 2 枚 `docker ps` 被我摘掉——本程硬约束禁 docker，代之以 §0.5 那枚一次性 docker 普查（是**换形**不是**漏检**）；② `index.txt` 逐发多记 `bytes=` 与包结果线，让"丢读数"不用开日志就看得见；③ 树身份 md5 点名**五枚**时序族文件而不是枚一枚（分母表 §3.1 末段点名 AC#15 必须补这一条）。严格串行：`while` 单发、无 `&`、无 `nohup`、无 `./...`、只名 `./internal/observe/` 一枚包、无 `-parallel`。 |
| `/d/tmp/wisp136ac15m-summarize.py` | 抄 `/d/tmp/wisp136ac11-summarize.py`（分母表 §3.3 逐条评过的那枚读数器），**修掉它点名的两枚缺陷**：① 顶层判定用严格 `^--- (PASS\|FAIL\|SKIP):`，缩进的子测试行**另记一枚 `indent=` 计数**（AC#11 那版 `^\s*---` 把两者混进同一计数器，今天包内 0 枚 `t.Run` 所以不咬人，但不再继承这枚隐患）；② **命中按句子判、不按测试名判**（分母表 §3.5 末条：按名筛会把 AC#14 那 17 枚故意变异红当目击、把率吹 18 倍）。另外每一枚顶层 FAIL 都连它的消息行原文打印出来，分类可复核。 |

族句集合（我的 `FAMILY` 正则，逐枚对应 §0.2 那六枚前提腿）：`only \d+ reads taken` · `the window only took \d+ reads` ·
`the fixture kept only \d+ of \d+ reads` · `an alternating tree must both keep and lose reads` ·
`this leg needs a window that recorded every read|nothing` · `expected several samples, got \d+`（`sampler_test.go:100`，
它**不带** `precondition broken` 字样——若只按那四个字筛就会漏掉这枚腿，这正是"按句子筛"要写成集合而不是一个串的原因）·
`the seam lost \d+ of \d+ reads`。

### 0.5 开跑前的争用普查（一次性，替代被禁的 docker 检查；其余三枚检查仍逐批发）

```
$ powershell Get-Process | Where ProcessName -match 'Runner.Worker|^go$|compile|^cgo'
(空)
$ gh run list -R CarlosShao/wisp --limit 5 --json databaseId,status,headSha
5 条全 "status":"completed" ⇒ in_progress = 0
$ docker ps        # 只读普查，不跑任何容器动作
4 枚第三方常驻容器（union-proxy 8天 / clipsync 5天 / clipsync-minio 13天 / clipsync-admin-int 2周）
  ⇒ 都是多日常驻、与 wisp 无关、非本窗新增；本机自托管 runner 的 slo-full 不在跑
$ ls -ld /e/work/base/actions-runner/_work
mtime 2026-09-24 21:11:49 +0800   ⇒ 距开跑 2 小时 22 分无 runner 活动
```

同窗口 `/d/tmp` 里**别家正在动**：`wisp-oldrule-check` 23:24、`citation-integrity-2026-09-24` 23:19、
`141-acc-*` 22:5x–23:0x。它们都不是 `go test` 进程（闸门第 1 枚空），但**本程每一批的 gate.txt 都留这一枚读数**，
批末 `gate-post.txt` 复扫；两枚都空才算这批在净窗内。

---

## 1. 试跑（口径 A 的形状标定）与口径 B（同进程连发）

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-24 23:47 +0800
$ git log -1 --format='%h %ad %s' --date=format:'%H:%M'
f1c9471 23:40 docs(A202): 清点程又挖出 6 枚从未被更正过的旧断言 …（本程只取 sha 与时刻，正文与 internal/observe 无关）
$ git diff --name-only beac693..HEAD -- internal/observe/ | wc -l
0                      ⇒ 快照钉 beac693 以来本包字节未动，下面所有读数与 §0 的名册同树
```

### 1.1 试跑 5 发（`/d/tmp/wisp136ac15m-A-pilot/`，只标定形状，不进分母）

`sh /d/tmp/wisp136ac15m-batch.sh <tree> …/wisp136ac15m-A-pilot 5 PILOT` ⇒ 闸门第一轮即 CLEAR。

| run | bytes | RUN | top PASS | top FAIL | top SKIP | indent | family | 真 `^panic:` | fatal | 包结果线 | rc | 包时 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 1-5 | **7547**（五发同一个值） | 71 | 71 | 0 | 0 | 0 | 0 | 0 | 0 | 1 | 0 | 3.851 / 3.966 / 3.882 / 3.993 / 3.996 s |

* 等式基线七项**逐发全过**；`PASS+FAIL+SKIP=71=RUN`；名册 `1 distinct top-level run-name sets over 5 runs; modal size=71 occurs=5`；
  逐发对 golden（`/d/tmp/wisp136ac15m-golden.txt`，71 名，从试跑第 1 发现取）**两向 `comm` 全空**。
* **字节带是单一值 7547**（比靶子表 §4.3 预期的"约 7.3 KB"更紧），⇒ 本程把它当 tripwire 用：任何一发不是 7547 就要打开看，
  任何一发小于 7547 就是丢读数。**分母表的 65 名基线 6717 B 与我的 71 名 7547 B 不冲突**（名册大 6 名）。
* 单发端到端墙钟 **8.1-8.5 s**（包时 3.85-4.0 s ＋ 固定开销 4.0-4.5 s）。那枚固定开销我单独量过一次
  （`go test -count=1 -v -run 'ZZZNosuchTest' ./internal/observe/` ⇒ 包时 0.149 s、端到端 3.98 s，日志
  `/d/tmp/wisp136ac15m-timing-probe.log`；**这一发是仪器定时探针，不是读数，不进任何分母**）。
  ⇒ 分母表 §4.2 按 3.8 s/发估的"A 1200 发约 76 分钟"在本机**偏低约 2.2 倍**，本程实测按 8.3 s/发记。
* 靶子腿试跑五发逐字 `--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)`。

### 1.2 口径 (B)：同进程 `-count=1200 -run '^TestCheckSettleHalfTheReadsFailedReportsItsLoss$'` —— **命中 3 发**

日志 `/d/tmp/wisp136ac15m-B1-target.log`（**完整原始日志留在快照里，未截断、未删**；1200 发 RUN 行全在）。
命令逐字（分母表 §4.1 (B) 的形状；`-run` 的 `$` 锚是承重的，没有它就会把 71 名各跑 1200 遍）：

```
cd /d/tmp/wisp136ac15m-tree-beac693 && go test -count=1200 -v -run '^TestCheckSettleHalfTheReadsFailedReportsItsLoss$' ./internal/observe/
```

等式基线（口径 B 的形状：单名，所以名册只该有这一枚；RUN 与判定必须闭合）：

| 项 | 读数 | 判定 |
|---|---|---|
| `^=== RUN` | 1200 | ✅ 名义发数 |
| `^--- PASS:` 顶层 | 1197 | ✅ |
| `^--- FAIL:` 顶层 | **3** | 见下 |
| `^--- SKIP` | 0 | ✅ |
| 闭合 | 1197＋3＋0 = 1200 = RUN | ✅ **没有丢读数**，⇒ 这 3 枚红不是"读数被吞"的假象 |
| 真 `^panic:` | 0 | ✅ |
| 包结果线 | 恰 1（`FAIL … 125.737s`） | ✅ |
| rc | 1（有红必然为 1，与"丢读数"是两味） | 记 |
| 闸门 | 批前 `/d/tmp/wisp136ac15m-B-gate.txt`：host procs 空、gh `5 completed`＝in_progress 0、runner `_work` mtime 21:11:49 | ✅ 净窗 |

**三枚命中的原文**（逐字从日志取，`--- FAIL` 前一行是它的 `=== RUN`）：

```
    sampler_settle_coverage_136_test.go:214: precondition broken: only 2 reads taken, half-and-half needs a window to lose in
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.11s)
    sampler_settle_coverage_136_test.go:214: precondition broken: only 2 reads taken, half-and-half needs a window to lose in
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.16s)
    sampler_settle_coverage_136_test.go:214: precondition broken: only 3 reads taken, half-and-half needs a window to lose in
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)
```

⇒ **口径 (B) 命中率 = 3 / 1200 = 0.250%**，名中的是 §0.2 第一枚前提腿（靶子腿）的 `:213 if tree.reads < 4` 那条前提守卫，
红在 `:214`。同发的另外两枚影子守卫（`:216 kept<2`、`:221 lost<2`）确实如靶子表 §1(a) 说的**没轮到响**——
三枚红全部印的是 `:214`，与"被 `:213` 遮住的形状"一致。

命中发生在连发的第 **319 / 511 / 522** 遍（不是开头也不是结尾，不聚簇）。单遍耗时分布（PASS 1197 遍 ＋ FAIL 3 遍 = 1200 闭合）：

| 该遍耗时 | 0.10s | 0.11s | 0.12s | 0.13s | 0.14s | 0.15s | 0.16s | 0.17s | 0.88s |
|---|---|---|---|---|---|---|---|---|---|
| PASS | 992 | 168 | 18 | 7 | 6 | 1 | 2 | 2 | 1 |
| FAIL | **1** | **1** | 0 | 0 | 0 | 0 | **1** | 0 | 0 |

* 众数 0.10s（100 ms 窗 ＋ 收尾），⇒ 三枚红里两枚"读数是 2"、一枚"读数是 3"。
* 有一枚 PASS 花了 **0.88s**（众数的 8.8 倍）却仍拿到 ≥4 枚读数 ⇒ 本包确实存在**百毫秒级的调度停顿**，
  但它落在窗外时不影响判定；这与靶子表 §4 C1 从归档那枚目击反推的"一次约 90 ms 的 inter-tick 空档"是同一类物理事件，
  **不是一次形状的重复**：我今天看到的是 `reads=2`（一枚约 90 ms 空档，与归档同形）**与** `reads=3`（约 33 ms 一拍，整体慢拍，归档没有这一形）。
  只登记，不外推。

### 1.3 口径 (B) 第二发：六枚前提腿同进程连发 —— **命中 0**

```
cd /d/tmp/wisp136ac15m-tree-beac693 && go test -count=200 -v -run
  '^(TestCheckSettleHalfTheReadsFailedReportsItsLoss|TestCheckSettleSingleTrustworthyReadReportsItsLoss|
    TestCheckSettleZeroFootprintDropsAreCountedToo|TestCheckSettleFullyMeasuredWindowReportsNoLoss|
    TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed|TestSampleStateAllMetricsAndVerdicts)$' ./internal/observe/
```
（命令在日志 `/d/tmp/wisp136ac15m-B2-family.log` 头两行原文里，上面为可读性换了行。）

| 项 | 读数 |
|---|---|
| RUN | 1200 ＝ 6 名 × 200 遍 ✅ |
| 顶层 PASS / FAIL / SKIP | 1200 / **0** / 0 ✅ 闭合 |
| 真 `^panic:` / 包结果线 / rc | 0 / 1 / 0 |
| 包时 | 146.873 s（`ok`） |
| family 句命中 | 0 |

⇒ **族口径 (B)：0 命中 / 1200 枚腿执行**；其中靶子腿自己 **0 / 200 遍**。
**不把 B1 与 B2 相加**：B2 里靶子腿的邻座换了（另外 5 枚腿在同一个进程里交替跑，GC/定时器堆状态与 B1 的"单名独跑"不同），
两发的分母不同形。作为噪声核对：若真率就是 B1 量的 0.25%，则 200 遍里 0 命中的概率是 `e^-0.5 = 61%` ⇒
**B2 与 B1 不互相推翻**，也不足以把 B1 拉平。

---

## 2. 口径 (A)：整包单发 `-count=1`，逐批发数

```
$ date "+%Y-%m-%d %H:%M %z"      （§2 开写时刻）
2026-09-25 06:52 +0800
$ git log -1 --format='%h %ad' --date=format:'%H:%M'
d171b42 00:01
$ git diff --name-only beac693..HEAD -- internal/observe/      （输出为空＝本包字节自快照钉以来未动）
$ git log --oneline beac693..HEAD | wc -l
11                     ⇒ 本程期间进了 11 枚 commit，上面那条路径限定 diff 为空 ⇒ 没有一枚动过 internal/observe/
```

⇒ 本程所有 A/B 读数与 §0 的名册、§1 的靶子文件 md5 同树；被验版本自始至终没换过。
（`d171b42` 是别家的台账 commit，本程没碰；`git log -1` 与我自己 commit 的 sha 不同＝共享工作树的正常现象，
每一枚我自己的 commit 都用 `git log -1 --format=%H -- <本文件>` 单独回读。）

### 2.1 一次必须披露的中断：A02 这批跨了一觉（本机挂起 6 小时 32 分）

A02 的 `index.txt` 逐发时刻原文（节选）：

```
run=29 rc=0 bytes=7547 … 3.921s; at=2026-09-25T00:02:04+0800
run=30 rc=0 bytes=7547 … 3.642s; at=2026-09-25T06:34:38+0800   ← 与上一发之间 6 h 32 min 的空档（本机挂起）
run=31 rc=0 bytes=7547 … 3.692s; at=2026-09-25T06:34:48+0800
run=32 rc=0 bytes=7547 … 6.745s; at=2026-09-25T06:34:59+0800   ← 醒后段里最慢的一发（众数约 3.7s）
```

处理方式（**不藏、也不据此判丢**）：

* 每一发都是**独立的 `go test` 进程**，跑的仍是同一棵固定树；挂起期间没有任何发在飞、没有第二个取样进程，
  所以这 50 发的**逐发等式**该过还是过了（见 §2.2 的逐批四数）。
* 但这批的**批前闸门读数对第 30-50 发已经过期**（闸门是 23:49 拍的、批内前 29 发都在 23:57 之前的净窗里跑完）。
  ⇒ 本程把 A02 记作**两段**：A02a＝第 1-29 发（净窗内，闸门有效）、A02b＝第 30-50 发（醒后段 21 发，
  闸门重扫见 `gate-post.txt` 06:36:53，host procs 空／gh `0 in_progress`／`_work` 未动），并把"醒后段第 3 发被拖慢到
  6.745 s（众数约 3.7 s）"这条协变量点名登记——**慢发正是更容易少一拍的形状**，所以它不是无害噪声，
  本程后面每次醒转都重扫闸门。
* A02b 那 21 发的四数仍然全等（71/71/0/0、panic 0、包线 1），所以按本程规则它们是〔有效段〕，
  **不是**〔丢读数〕；我没有因为"机器睡过"就整批作废，也没有把它并入 A02a 当同一净窗。

### 2.2 逐批读数（每批一条闸门记录、一发一份日志）

| 批 | 名义发数 | 过等式基线的**有效发数** | RUN | 顶层 PASS | 顶层 FAIL | 顶层 SKIP | 缩进判定 | 真 `^panic:` | 包结果线 | index 行 | rc≠0 | 名册 distinct sets | 字节带 | **族句命中** |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 试跑 | 5 | 5 | 355 | 355 | 0 | 0 | 0 | 0 | 5 | 5 | 0 | 1（71×5） | 7547 | 0（§1.1，不进分母） |
| A01 | 50 | **50** | 3550 | 3550 | 0 | 0 | 0 | 0 | 50 | 50 | 0 | 1（71×50） | 7546-7547 | **0** |
| A02a | 29 | **29** | 2059 | 2059 | 0 | 0 | 0 | 0 | 29 | 29 | 0 | 1 | 7546-7547 | **0** |
| A02b | 21 | **21** | 1491 | 1491 | 0 | 0 | 0 | 0 | 21 | 21 | 0 | 1 | 7546-7547 | **0** |
| A03 | 60 | **60** | 4260 | 4260 | 0 | 0 | 0 | 0 | 60 | 60 | 0 | 1（71×60） | 7546-7547 | **0** |
| 小计（A01..A03） | 160 | **160** | 11360 | 11360 | 0 | 0 | 0 | 0 | 160 | 160 | 0 | 每批各 1 枚 | — | **0** |

等式基线七项（靶子表 §4.3／分母表 §4.3 那七条）逐批发落点：

| 检查 | 期望 | A01/A02/A03 实测 |
|---|---|---|
| `grep -c '^=== RUN'` | 71 | 逐发 71 ✅ |
| 顶层判定数 | 71 | 逐发 `PASS+FAIL+SKIP = 71` ✅（`PASS=71`、`FAIL=0`、`SKIP=0`） |
| `^--- PASS` | 71 | 71 ✅ |
| `^--- FAIL` | 0 | 0 ✅ |
| `^--- SKIP` | 0 | 0 ✅ |
| 真 `^panic:` | 0 | 0 ✅（另记：日志里含 `Panic` 的行 4 枚/发＝**标识符味**，是 `TestPanicInWorkerSurvivesAndCancelsRoot` 与 `TestFakeTreeEmptyScriptFailsClosedAndNotPanics` 各 `=== RUN`＋`--- PASS` 两行；票面 `>`④ 那句"7 枚命中全是标识符、其中一枚名是 `TestCheckSettlePanicInSamplerDoesNotStopTick`"**按名现量为 0 枚**——`grep -rn TestCheckSettlePanicInSamplerDoesNotStopTick --include=*.go .` 0 命中，那枚用例不存在（与分母表 §3.3 的判读一致），真值＝**4 行/2 枚名**） |
| 包结果线 `^(ok\|FAIL)[[:space:]]+github` | 恰 1 | 逐发 1 ✅（全是 `ok`） |
| `index.txt` 有本发行 | 有 | 逐批 50/50、50/50、60/60 行，`rc=0` 全一致 ✅ |
| 名册两向差集 | 空 | 对 golden（71 名，试跑第 1 发现取）`golden-not-run=[]`／`union-not-in-golden=[]` ✅ |

**没有一段被判〔丢读数〕，`_void 段数_ = 0`**；名义发数＝有效发数＝160。

字节带更正 §1.1 的说法：绿发的字节不是单一值，而是 **7546 / 7547 两值**，差 1 字节，
成因是 Go 打印用例耗时时会裁尾零（`(0.1s)` 与 `(0.10s)`）。⇒ 本程把它当**带**用（7546-7547），
真正的 tripwire 是"四数＋包线＋index 行"那几枚，字节只作辅助。归档的 65 名单值 6717 B 与我的两值带不冲突。

单发墙钟（同批内）：A01 包时中位 4.137 s／A02 中位 3.643 s；端到端 A01 约 8.2 s/发、A03 约 4.9 s/发（凌晨机器更空，
别家没有取样在飞）。⇒ 分母表 §4.2 的"A 1200 发约 76 分钟"在本机是 **100-170 分钟**的量，取决于编队是否安静。

---

## 3. 命中 ≥1 发 ⇒ 归因（**只归因，不动任何码**）与最小修法候选清单

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 07:00 +0800
$ git log -1 --format='%h %ad' --date=format:'%H:%M'
d171b42 00:01        （本程 §2 之后没再进新的别家 commit；internal/observe/ 自 beac693 起 diff 仍为空）
```

### 3.1 命中的是哪一枚腿的哪一条前提腿

三枚命中**同一枚腿、同一条前提腿**：

| 项 | 值 |
|---|---|
| 腿 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（§0.2 名册第一枚，票 136 AC#12 装的） |
| 文件 | `internal/observe/sampler_settle_coverage_136_test.go`（快照内 md5 `a3196802…`，与分母表 §3.5 记的 HEAD 值同） |
| 响的守卫 | `:213 if tree.reads < 4 {` ⇒ 红句在 `:214`，逐字 `precondition broken: only %d reads taken, half-and-half needs a window to lose in` |
| 被遮的两枚影子守卫 | `:216 kept < 2`、`:221 lost < 2` —— 三枚命中里**一枚都没轮到响**，与靶子表 §1(a) 的"被 `:213` 遮罩"判读一致 |
| 结构上不可能响的那枚 | `:227 kept > lost+1 \|\| lost > kept+1`（`alternatingTree` 按 `reads%2` 交替 ⇒ `\|kept-lost\|<=1` 恒成立），本程零命中，与靶子表一致 |
| 完整原始日志 | `/d/tmp/wisp136ac15m-B1-target.log`（1200 发全量、未截断；三枚红分别在第 319／511／522 遍） |

红因族属：`precondition broken: only N reads taken` 那一族——**窗口没读满**，不是断言失败。

### 3.2 归因（逐条对齐靶子表 §4 的候选，指明我今天的数据支持到哪一步）

* **C1「读数控条完全由 tick 交付决定，且没有前置读」——支持。** 盘上再确认：`CheckSettle` 的读循环
  `sampler.go:491-526`（`deadline := startAt.Add(within)` / `t := time.NewTicker(interval)` / 读只在 `case <-t.C` 之后发生），
  本腿调用点 `coverage:208` 传 `within=100ms`、`interval=10ms` ⇒ 期望约 10 拍、门槛 4 拍。
  今天三枚实测读数 **reads = 2 / 2 / 3**，即交付被砍到期望的 20-30%。
* **C2「窗口与断言之间缺一条'允许被判定'的判据」——支持，且今天第一次有了正例数据。**
  这条是判据② 要的根因本体：**没有任何一处代码在等这一窗读满**（守卫是 one-shot，票面 `>`② 已就此裁过现量方的"重试预算"维——我这程复核的是**更强的一条**：
  `grep -n 'attempts\|retries\|for i := 0\|time.Sleep\|Poll\|Ticker\|Timeout' internal/observe/sampler_settle_coverage_136_test.go internal/observe/sampler_settle_gate_136_test.go`
  ⇒ **0 命中**（两枚文件里既没有重试计数、也没有轮询、也没有 Sleep、也没有自建 Ticker），所以"某枚腿比另一枚多一层等待"在盘上确实没有代码支撑）。
* **C3「定时器粒度」——不支持"单独成因"。** 若只是粒度粗（15.6 ms），期望读数会稳定落在 6 左右而不是 2-3；
  今天两枚 `reads=2` 只能是"至少一次约 50-90 ms 的空档"，与归档那枚目击同形。
* **C4「STW／GC 停顿」——今天的数据第一次把它从"待验"推到"有指向"，但方向与靶子表 §5.2 的预期相反。**
  靶子表 §4 C4 与 §5.2 写的是：同进程连发**看不到**这个症状（归档 0/1004），A 才是承重的率口径。
  我今天量到的是 **B 命中 3 枚、A 前 260 发命中 0 枚**。⇒ 记两条，别混：
  1. **A 的 0/260 不构成"A 更低"的证据**：若真率就是 0.25%，260 发全绿的概率是 `(1-p)^260 = 52%`，一枚硬币。
  2. 但**"B 是 insensitive 那一维"这句前提被实测推翻了**——B 不仅看得见，而且是今天唯一看见的地方。
     机制上的候选解释（只算候选，本程不动码、不再造仪器去证）：B 里这枚腿每秒被重跑约 10 次、
     堆与定时器堆持续 churn，而 A 里它整包只跑 1 次且前面刚跑完 70 枚别的用例。这与靶子表 §4 C4 自己写的
     "`-count=N` in one process makes GC pressure per shot much higher than a single whole-package shot"
     **方向一致**——它当初用这句话来解释"为什么 B 看不见"，我看见的却是相反的一侧：**GC 压力高到一定量，
     就会把一次 STW 塞进那 100 ms 里**。⇒ 这一条留给裁决者：靶子表 §5.2 那句"B 不是能证'没有'的探针"仍成立
     （本程 0/200 与 3/1200 不互相推翻），但它附带的"B 大概率不重现"预测**在当版本上不成立**。
* **没被排除、但也无法在本程定的那一维**：本机醒转/挂起。三枚命中全在 23:40-23:42 那段净窗里取的（闸门 §1.2 原文），
  所以**不能**归给"机器睡眠"；A 段里跨睡的那 21 发也零命中。

### 3.3 最小修法候选清单（**交回给你判，本程一枚都没动手**）

按"只许把等待变成有判据的等待／⛔`time.Sleep` 糊窗／⛔Skip／⛔放宽阈值／⛔删断言／地界＝`internal/observe/**_test.go`"这把尺排：

| # | 候选 | 射程 | 代价与副作用（我这程数据能说的） | 我的判读 |
|---|---|---|---|---|
| 1 | **轮询到判据＋单调超时上界**：把"开一窗就直接数 `tree.reads`"改成"`while tree.reads < 4 且未到 observe.Timeout 上界：再开一窗（或再等一拍）`"，到条后再判 kept/lost | 只动 `sampler_settle_coverage_136_test.go`，六枚族腿同形可复用 | 需要一枚**有界的**等待；上界必须由 §3.4 的算术推出来，不是拍脑袋（票面 `>`⑤ 已把这条写进判据） | **唯一与判据③ 字面同形的形状**。代价：单发墙钟增加，A 口径的取样暴露面变大 ⇒ 复跑批次要重算 n |
| 2 | 给 `CheckSettle` 补一枚"前置读"（对齐 `SampleState` 的 `2 + ticks` 形状），使读数控条不再纯靠 tick | **要动 `sampler.go`** | 改变 `CheckSettle` 的出线语义形状 | ⛔ 触发本格的停手线（票面 `>`⑥：那属人工批准面），**不是我能派的修法**，只登记 |
| 3 | 只把 `within` 从 100 ms 抬到比如 300 ms | 只动测试 | 期望拍数 10→30，门槛 4 的余量变大 | ⛔ 靶子表 §4 C3 已点名"抬 interval／加宽窗口只是移均值，不补判据"，且今天的 `reads=2` 是一枚大空档，宽窗照样可能撞上 ⇒ **判"不算修法"** |
| 4 | 把前提腿从"红"改成"重开一次并在报告里自陈重开了几次" | 只动测试 | 需要新增一条自陈行，且要与 AC#14 的门行折叠逻辑不打架（门已 `true`） | 可行但要额外一枚断言钉"重开次数有限且可见"，否则就是 1 的弱化版＋藏信息 ⇒ **次选，须带自陈钉** |
| 5 | `t.Skip`／删这条腿／把 `< 4` 调成 `< 2` | — | — | ⛔ 三枚都是本票口径明禁（跳过＝把"没测"洗成"通过"；调门槛＝放宽断言）。列在这里只为让"没做"看得见 |

⇒ 交给你的问题只有一个：**候选 1 的超时上界取多少、由谁算**。我这边能给的算术在 §3.4。

### 3.4 上界算术的原料（本程现量，不抄归档）

| 量 | 读数 | 来源 |
|---|---|---|
| 该腿单遍众数耗时 | **0.10 s**（1200 遍里 0.10 s 占 993 遍、0.11 s 占 169 遍，两档合计 1162/1200） | B1 日志的每遍 `(0.10s)` 字段 |
| 该腿单遍最长 | 0.88 s（一枚 PASS） | 同上 |
| 窗口预算 | `within=100 ms`、`interval=10 ms` ⇒ 期望约 10 拍 | 调用点 `coverage:208` 字面量 |
| 门槛 | 4 拍（`:213`） | 现量 |
| 一次拍间空档的量级 | 命中枚 `reads=2` ⇒ 单档约 50-90 ms；`reads=3` ⇒ 平均约 33 ms | 由 `reads` 反推，**不是测到的 tick 间隔**（本程没加计时仪器） |

⇒ 若按"上界＝期望拍数打满所需时间的若干倍"来给候选 1 定超时，**这几枚数是现成的**；具体倍数与是否复用
`internal/observe/clock.go` 里那枚 `observe.Timeout`（靶子表 §4 C2 指名它）请你拍。

### 2.3 A04 / A05 落账、一枚争用事件的现量、以及仪器中途加强（batch2）

> 编号是 §2 的，**位置在 §3 之后**——因为 A05 这批在写 §3 时还没跑完（07:10 才落账）。渐进写的代价，不是笔误。

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 07:12 +0800
$ git log -1 --format='%h %ad' --date=format:'%H:%M'
58e2e40 07:03      （＝我 §3 那枚 commit；此后别家未再进 commit，`git diff --name-only beac693..HEAD -- internal/observe/` 仍为空）
```

| 批 | 名义发数 | 有效发数 | RUN | 顶层 PASS | FAIL | SKIP | 缩进 | 真 `^panic:` | 包线 | index 行 | rc≠0 | 名册 distinct sets | 字节带 | **族句命中** |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| A04 | 100 | **100** | 7100 | 7100 | 0 | 0 | 0 | 0 | 100 | 100 | 0 | 1（71×100） | 7547-7549 | **0** |
| A05 | 120 | **120** | 8520 | 8520 | 0 | 0 | 0 | 0 | 120 | 120 | 0 | 1（71×120） | 7547-7549 | **0** |
| **累计 A01..A05** | **380** | **380** | **27080** | **27080** | **0** | **0** | **0** | **0** | **380** | **380** | **0** | 每批各 1 枚 | — | **0** |

`_void 段数_ = 0`（没有任何一段因等式不过被判丢读数）。名册两向差集对 golden 逐批为空。

单发包时（`index.txt` 里那条 `ok … NNNs` 现算，n＝该批发数）：

| 批 | min | 中位 | p90 | max | 这段机器在干什么 |
|---|---|---|---|---|---|
| A01 | 3.777 | 4.137 | 4.385 | 4.794 | 23:37-23:49，别家还在动（`141-acc-*` 22:5x-23:0x、`wisp-oldrule-check` 23:24、`citation-integrity` 23:19/23:34 连投 5 枚 commit）⇒ **本批是全程最被拖慢的一段** |
| A02 | 3.506 | 3.643 | 3.899 | **6.745** | 尾段醒转（§2.1） |
| A03 | 3.448 | 3.495 | 3.599 | 3.888 | 凌晨净窗 |
| A04 | 3.446 | 3.482 | 3.556 | 3.613 | 净窗（分布最紧） |
| A05 | 3.439 | 3.473 | 3.520 | 3.646 | 净窗，见下面的争用事件 |

**一枚争用事件，现量到、没洗掉**：07:03:38 我在 A05 跑动期间扫到别家的 `compile.exe`（pid 54408）与一枚 06:55 新建的
`/d/tmp/wisp-clean-43check`。当时的判断是"这批发在争用窗里 ⇒ 按 AC#11 口径该整批判脏"，但**读数不支持这个判语**：
A05 的包时分布是 3.439-3.646 s、中位 3.473，与紧邻的净窗批 A04（中位 3.482）**同带**，那枚 `compile.exe` 只存在了一瞬。
⇒ 处置：**A05 不作废、也不改口径**，只在这里点名；同时把仪器加严（不做"挑掉慢发"这类事后修饰）：

* 新仪器 `/d/tmp/wisp136ac15m-batch2.sh`（保留 `wisp136ac15m-batch.sh` 不动，A01-A05 那 380 发仍由旧版产生，可复算）：
  **每 10 发在两次发射之间重扫一次**别家取样进程（`go.exe`／`compile.exe`／`cgo.exe`／`Runner.Worker.exe`，
  探针对语义是"上一发已返回 ⇒ 扫到的必是他家"），非空就**等**（同 60×30 s 的界，超界 `exit 8`），
  绝不穿过争用窗取样；每发的 `index.txt` 多带一枚 `foreign=go:N,compile:N,cgo:N,runner:N`。
  探针函数单独验过：`go=0,compile=0,cgo=0,runner=0 -> 0`、`go=1,compile=2,cgo=0,runner=0 -> 3`。
* A06 起改用 batch2。

**执行位形的披露（防被误读成"我并发了"）**：A02 与 A04 因为单发×发数超过了工具的 10 分钟上限而被 harness **自动转后台**；
A05 之后我主动用后台任务位承载这条串行循环。两种情况下**循环内部都仍是单发串行**——脚本里没有 `&`、没有 `nohup`、
没有 `xargs -P`/`parallel`、包路径只名一枚 `./internal/observe/`、无 `-parallel`；跑测期间我只做过轻 I/O
（写这份表、`grep`、`git` 读、`tasklist`），**没有并发过第二枚取样命令**，任何时刻在飞的 `go test` 至多一枚。

### 2.4 A06（batch2 首批）与口径 B 的第二发 B3（`-count=3000`）

> 编号归 §2，**位置在 §3 之后**：A06 与 B3 是在写完结论段 §3 之后才落地的读数（渐进写、不回头改已提交的段落）。

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 07:39 +0800
$ git log -1 --format='%h %ad' --date=format:'%H:%M'
7f99963 07:12      （＝我 §2.3 那枚 commit；`git diff --name-only beac693..HEAD -- internal/observe/` 仍为空）
```

**A06**：120 发全部过等式基线，命中 0。

| 项 | 读数 |
|---|---|
| RUN / 顶层 PASS / FAIL / SKIP | 8520 / 8520 / 0 / 0（逐发 71/71/0/0，`PASS+FAIL+SKIP=RUN`） |
| 缩进判定 / 真 `^panic:` / fatal | 0 / 0 / 0 |
| 包结果线 / index 行 / rc≠0 | 120 / 120 / 0 |
| 名册 | `1 distinct set over 120 runs; modal size=71 occurs=120`，对 golden 两向为空 |
| 中途闸门 | `gate.txt` 里 **0 枚 `BUSY`** ⇒ 12 次中途重扫全清空，这批没有在争用窗里取样 |
| 字节带 | **7556-7558**＝旧批次的 7547-7549 **＋9**。成因已定位不是谜：batch2 的日志头多一串 ` (batch2)`（9 字节）。⇒ 换仪器会把字节带整体平移，**字节只作辅助 tripwire**，判有效与否看四数＋包线＋index |

**B3**（口径 B 第二发，命令与 B1 逐字同形、只把 `-count` 从 1200 抬到 3000；日志
`/d/tmp/wisp136ac15m-B3-target.log`，372408 字节，包时 303.180 s）：

| 项 | 读数 |
|---|---|
| `^=== RUN` | 3000 ✅ 名义 |
| 顶层 PASS / FAIL / SKIP | 2999 / **1** / 0 ⇒ 闭合 `2999+1+0=3000` ✅ **没有丢读数** |
| 真 `^panic:` / 包结果线 / rc | 0 / 恰 1（`FAIL … 303.180s`） / 1 |
| 命中原文 | `sampler_settle_coverage_136_test.go:214: precondition broken: only 2 reads taken, half-and-half needs a window to lose in` ＋ `--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.12s)` |
| 位置 | 连发第 **578** 遍（B1 的三枚在 319／511／522 ⇒ 四枚命中分布在 319-578 这一段，1200-4200 遍的长尾里没有第五枚） |
| 该遍耗时分布 | 0.10 s 2942 遍、0.11 s 45 遍、0.12 s 11 遍（其中 1 遍红）、0.14 s 1、0.16 s 1 |

⇒ **口径 B 的两发要分开报，不能只报好看的那一枚**：

| 发 | 命中 | 率 | 精确（Garwood）95% 区间 |
|---|---|---|---|
| B1 `-count=1200`（23:40，别家还在动） | 3 | **0.250%** | 0.068% - 0.731% |
| B3 `-count=3000`（07:24-07:29，净窗） | 1 | **0.033%** | 0.001% - 0.185% |
| 合并（同形、同树、间隔 8 小时） | 4 | **4/4200 = 0.095%** | 0.026% - 0.218% |

两发相差 7.5 倍，**但这不构成"率变了"的证据**：按合并率 0.095% 推，B1 那 1200 遍里期望 1.14 枚、看到 ≥3 枚的概率
`1 - e^-1.14(1+1.14+1.14²/2) = 10.5%`；B3 那 3000 遍期望 2.86 枚、看到 ≤1 枚的概率 `e^-2.86(1+2.86) = 22%`。
两枚尾概率都在 10-22% 这一档 ⇒ **判"两发相容、样本还小"，不判"率随时间掉了"**。要分辨 0.25% 与 0.033% 需要
再往同形里投几千遍（本程没投，代价见 §4）。

（下接 §4 两口径命中率总账。）
