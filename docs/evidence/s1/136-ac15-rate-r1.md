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

（下接 §1 试跑与逐批发数。）
