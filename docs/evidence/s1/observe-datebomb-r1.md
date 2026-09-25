# observe 日期定时炸弹：修复 + 兄弟普查（r1）

程：实现程（write-capable）。派单：编排者 `A212 next=①`（"先让日期炸弹那枚落地"，`pending-and-issues.md:5890`）。
本格只动 `internal/observe/**` + 本文件。**生产码一行未动**（见 §2）。

---

## 0. 锚点与时刻

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 08:36 +0800
$ git rev-parse --short HEAD
038eb02
```

⚠ 共享工作树：本程开局 HEAD＝`5540087`，写作期间 HEAD 走到 `42d141d`（编排者登记本炸弹的那一枚）→`98873a5`→`038eb02`。
下面每一处 `file:line` 都在**当前 HEAD** 上重读或重跑过；本节锚点 `038eb02`。

---

## 1. 前提核账（派单给的数，逐条对盘）

派单的四个 `file:line` **全部命中，一字未漂**：

| 派单声称 | 现读 | 结论 |
|---|---|---|
| `logging_test.go:190` `old := …logFilePrefix+"20260101-001"+logFileExt` | 同一行、同一文本 | 对 |
| `logging_test.go:198` `fresh := …logFilePrefix+"20260918-001"+logFileExt` | 同一行、同一文本 | 对 |
| `logging.go:162` `now func() time.Time // wall clock, injectable for tests` | 同一行 | 对 |
| `logging.go:181` `func newRollingWriterClock(dir string, sizeMB, days int, now func() time.Time)` | 同一行 | 对 |

派单里两处**需要更正**的东西，都不是事实错误而是抄写错误，但第二处是会漏检的方法错误：

1. **包路径**：派单印的是 `github.com/github.com/CarlosShao/wisp/internal/observe`（双 `github.com`）。
   现跑真实读数是 `github.com/CarlosShao/wisp/internal/observe`（`go.mod:1 module github.com/CarlosShao/wisp`）。派单自己已打了问号，按现跑数取。

2. **⚠ 派单点名的第一条 grep 是一条恒不匹配的模式，它连本票这枚炸弹都抓不到。**
   ```
   $ grep -rn '"20[0-9][0-1][0-9][0-3][0-9]' --include='*.go' internal cmd tools
   （无输出，exit=1；加 -E 同样为空）
   ```
   原因不是没有炸弹，是**字符类排错位**：该模式把"月份"类 `[0-1][0-9]` 放在第 4–5 位，于是整条只能匹配 `200X-MM-DD` 形状的年份，
   而 `2026-09-18` 的第 4 位是 `6`、不属 `[0-1]` ⇒ **所有 202X 年的紧凑日期一律漏检**，包括 `logging_test.go:190/198` 这两枚本票主角。
   派单第二条 `2026[0-9][0-9][0-9][0-9] --include='*_test.go'` 是好的、命中 4 行（见 §8），但它只扫 `_test.go`、且只扫 `2026` 开头。
   ⇒ **本节普查用的是改正后的模式**（`2[0-9]{3}(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])` 等），枚数见 §8，
   并附一条工具警告：Qoder 的 Grep 工具对"量词后接分组择一"的模式会**静默返回 No matches**（同一模式在 `grep -E` 下同样为空，故此处两者都不可信），
   本节所有计数以 `grep -rEn` 现跑为准。

---

## 2. 修法：把两处字面日期改成从注入时钟派生

**生产文件改动：零。** `newRollingWriterClock` 这一枚缝**本身就够用**，理由（现读 `logging.go`）：
`ensureFileLocked` 取日名用 `w.now()`（`logging.go:226`）、`sweepLocked` 算 cutoff 用 `w.now()`（`:322`）、
`Close` 不读时钟；该用例不建 registry、不起 `flushLoop`，所以注入的 `fixedNow` 覆盖了这个用例摸得到的全部时钟面。
`newRollingWriter`（`:172`，把 `time.Now` 烧进去的那一枚）在本用例里**被换掉而不是被改**，因此 `logging.go` 不必动一字。

改动只在 `internal/observe/logging_test.go` 的 `TestRollingWriterRetentionSweep` 一枚函数体（`+23/−5`）：

```go
fixedNow := time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC)
oldDay   := fixedNow.Add(-12 * 24 * time.Hour) // outside: must be pruned
freshDay := fixedNow.Add(-2 * 24 * time.Hour)  // inside: must survive
…
w, err := newRollingWriterClock(dir, 0, 7, func() time.Time { return fixedNow })
```

两枚夹具名由 `oldDay.Format(logDayLayout)` / `freshDay.Format(logDayLayout)` 现算，`os.Chtimes` 也用 `oldDay`（mtime 与名同日，
因为 `sweepLocked` 先解析名、只在解析失败时才退回 mtime，夹具自相矛盾就两样都钉不住）。

**断言强度不变**：仍然是两条、仍然相反——`os.Stat(old)` 必须 `IsNotExist`（窗外那份**被清**），
`os.Stat(fresh)` 必须成功（窗内那份**不被清**）。没有 Skip、没有把断言换成"没报错就算过"、`days=7` 没动、`retentionSweepPeriod` 没动、
`thresholds.go` 与 golden 没动、`tools/d22scan/allowlist.txt` 没动。

**为什么不是把 `20260918` 换成 `time.Now().Add(-2*24h)`**：那只是把炸弹挪到每天午夜，并让 CI 上的 runner 因时区/起跑时刻不同而漂；
注入常量使这枚用例对真实墙钟**永久无感**。派单此条我认。

---

## 3. 证 1：修之前它**确定性**咬人，且咬的原因是日期而不是别的

（以下在**仓外快照** `D:/tmp/observe-datebomb-sess/snap-head`，`git archive HEAD` 于 `42d141d` 解出；仓内文件从未被这一支改动过——
"临时件只建不删"，快照留在 `D:/tmp` 不删。）

**A. 快照内原码＝红**（与本机现跑同形）：

```
--- FAIL: TestRollingWriterRetentionSweep (0.01s)
    logging_test.go:216: fresh file wrongly pruned: GetFileAttributesEx
    C:\Users\swq\...\Temp\TestRollingWriterRetentionSweep822192421\001\wisp-20260918-001.jsonl:
    The system cannot find the file specified.
FAIL	github.com/CarlosShao/wisp/internal/observe	0.111s
```

**B. 只把夹具名挪进窗内（`20260918`→`20260924`，一行，其余原样）＝绿**：

```
=== RUN   TestRollingWriterRetentionSweep
--- PASS: TestRollingWriterRetentionSweep (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/observe	0.042s
```

A→B 之间唯一的自变量是那枚日期 ⇒ **坐实"失败由日期相对真实墙钟出窗造成"，不是环境问题、不是别的改动**。
派单要求的"临时变体"走的是**仓外快照**这一支，仓内那枚临时改动从未存在过；快照未提交、不删。

---

## 4. 证 2：修之后绿，且 `-count=2` 可复现

```
$ go test ./internal/observe/ -run 'TestRollingWriterRetentionSweep' -count=2 -v
=== RUN   TestRollingWriterRetentionSweep
--- PASS: TestRollingWriterRetentionSweep (0.01s)
=== RUN   TestRollingWriterRetentionSweep
--- PASS: TestRollingWriterRetentionSweep (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/observe	0.039s
```

追加一枚**未来免疫**的直接演示（同一份新代码，快照 `snap-mut` 里只把 `fixedNow` 挪到 2031-03-14）：

```
=== V1) fixedNow -> 2031-03-14
--- PASS: TestRollingWriterRetentionSweep (0.01s)
ok  	github.com/CarlosShao/wisp/internal/observe	0.037s
```

夹具日随注入时刻一起走 ⇒ 这枚用例的形状里**没有"再过 N 年就会翻"的常数**。挪回 2026 后快照复绿（§5 的 M0 控制读数）。

---

## 5. 证 3：变异——把生产 sweeper 打断，新用例仍为**对的因**变红

快照 `D:/tmp/observe-datebomb-sess/snap-mut`（`git archive HEAD` + 我这枚新测试文件）。**四发变异，仓内 `logging.go` 全程未动**；
每发只改快照一行，跑完即回，末了 `cp` 回原文件并复跑整包证明快照回到绿（`ok … 3.453s`）。

| 变异 | 改动（快照内 `logging.go`） | 结果 | 红的因 |
|---|---|---|---|
| M0 控制 | 无 | **PASS** `ok 0.034s` | — |
| M1 永远删 | `if t.Before(cutoff) \|\| true` | **FAIL** `logging_test.go:225` | sweeper 去删 writer 自己那枚**正打开着**的当日文件，Windows 拒绝 → sweep 返错，用例在 sweep 的错误检查处响亮红 |
| M2 永远不删 | `if t.Before(cutoff) && false` | **FAIL** `logging_test.go:231: stale file not pruned: <nil>` | **删的那一半承重** |
| M3 窗缩到 1 天 | `cutoff := w.now().UTC().Add(-24*time.Hour)` | **FAIL** `logging_test.go:234: fresh file wrongly pruned: … wisp-20260916-001.jsonl` | **留的那一半承重**，且红点是派生出来的窗内夹具名 `20260916` 本身 |
| M4 sweeper 无视注入时钟 | `cutoff := time.Now().UTC().Add(…)` | **FAIL** `logging_test.go:225`；**整包跑下来只有它红**（`--- FAIL: TestRollingWriterRetentionSweep`，其余 70 枚绿） | 新用例是**全包唯一**钉住"sweep 必须用注入时钟"这条的仪器 |

M1/M4 的红点落在 `:225`（sweep 自己的错误返回）而不是 `:234`，是 **Windows 文件锁的形状**：`newRollingWriterClock` 会先按 `fixedNow` 当日开一枚活文件，
"把所有文件都删"的 sweeper 必然撞上这把活锁，于是先在这一行炸。**两发的结论都是红、不是静默通过**；
要干净地只咬"窗内那份被误删"这一半，看 M3——它刻意让活文件留在窗内，红点就是 `:234` 的那条断言。
（这一支也顺带量到：旧写法在同一变异下的红点位置取决于操作系统，Linux 上 `os.Remove` 对打开文件不报错，红点会挪到 `:234`。）

> 派单那句"改完还能过断言的改写比炸弹更坏"——上面 5 发是本程据此给的**唯一**答卷：任何一发若变绿，这枚改写就不该交。

---

## 6. 证 4：包级无回归，四数 + 名册

同一命令、同一棵树，`-count=2 -v`，`rc` 与五个计数（`^panic:` 单列）：

| 读数 | 修前（`08:28`，HEAD `42d141d`） | 修后（`08:34`，HEAD `98873a5`） |
|---|---|---|
| `rc` | **1** | **0** |
| `=== RUN` | 142 | **142** |
| `--- PASS` | 140 | **142** |
| `--- FAIL` | 2 | **0** |
| `--- SKIP` | 0 | **0** |
| `^panic:` | 0 | **0** |
| 包结果线 | `FAIL … 6.950s` | `ok … 7.110s` |

* **修前那 2 枚 FAIL 是同一枚用例的两发**（`-count=2`）：`uniq -c` 读数是 `2 TestRollingWriterRetentionSweep`，distinct 失败名＝**1 枚**。
* **名册没缩**：distinct `=== RUN` 名 **71 → 71**；`comm -23`（修的没了的）＝空，`comm -13`（新增的）＝空。
* **没有 panic 吞读数**：两趟 `^panic:` 都是 0，所以不存在"某几枚用例产不出读数"的缺口——这条是派单点名的，本包这次不适用，
  且它和 `A212` 里记的那枚"空切片无守卫 ⇒ RUN=49/71、22 枚没读数"的形状**不是同一枚**，本包这趟两趟都跑满了 71 枚。
* 与编排者基线对齐：`pending-and-issues.md:5767` 记的 09-24 22:52 基线是 `RUN=142/PASS=142/FAIL=0/SKIP=0/panic=0`——
  与本程**修后**逐字相同 ⇒ 修后回到基线形状，且 `internal/observe/` 相对那枚基线只差了本票这一枚用例的写法。

---

## 7. 证 5：格式与静态检查

```
$ gofmt -l internal/observe/logging_test.go
（无输出＝干净）
$ go vet ./internal/observe/
（无输出＝干净）
```

`newRollingWriter`（`:172`，把 `time.Now` 烧进去的那一枚）仍被 `TestRollingWriterSizeRoll`（现读 `logging_test.go:133`）与生产 `InitLogWithRegistry` 使用，
换掉本用例这一处调用**不产生未使用符号**（vet 已证）。

---

## 8. 兄弟炸弹普查（派单 `next=①` 点名要的那枚"枚数"）

**本节时刻锚点**：`2026-09-25 08:4x +0800`，HEAD `038eb02`（+ 我这枚未提交的证据文件）。
普查在 §6 的修后树上跑，也在 `HEAD~1`（修前）上重跑过对照，两组的差只可能是本票那 3 枚字面量。

### 8.1 扫了哪些形状（每形的**现跑枚数**，不是目测加总）

| 形状 | 命令要点 | 现测量 | 处理 |
|---|---|---|---|
| 紧凑 `YYYYMMDD` | `grep -rEno '2[0-9]{3}(0[1-9]|1[0-2])(0[1-9]\|[12][0-9]\|3[01])' --include='*.go' internal cmd tools` | **修前 6 枚 → 修后 5 枚** | 逐枚分类见 8.2；5 枚里有 1 枚是 Go 参考布局 `20060102`（`logging.go:33`，惰性）、1 枚是**我新写的注释里复述死掉的那枚 `20260918`**（`:193`，惰性）⇒ 真夹具 5→3 |
| ISO `YYYY-MM-DD` | 同域同法 | **48 行命中** | 拆：**30 行是注释正文**（记的是裁定/测量日期，惰）＋**18 行是代码行**，18 里再减 **6 行 Go 参考布局**（`prompt.go:262`／`retention.go:209`／`retention_test.go:136,137,138`）⇒ **真 ISO 字面量 12 枚**，逐枚见 8.3 |
| `time.Date(` | `grep -rEn 'time\.Date\('` | **6 枚**（修前后同数；炸弹那枚 `time.Date(2026,1,1,…)` 被换成 `fixedNow` 派生，本文件仍留 2 枚 `time.Date`） | 逐枚见 8.4 |
| `time.Unix(` | `grep -rEn 'time\.Unix\('` | **13 枚** | 8.4 末行：全是注入的假时钟或纯换算 |
| 相对窗 `Add(-N*24h)`（测试里） | `grep -rEn 'Add\(-?[0-9]+ \* 24 \* time\.Hour\)' --include='*_test.go'` | **修前 0 枚 → 修后 2 枚**（都是我的 `oldDay`/`freshDay`，锚在注入时钟上） | 见 8.5 那条"相对形"单独结论 |
| 混合启发式（同文件既有日期字面量又调 `time.Now()`） | 两遍 grep 取交集 | **3 个文件** | 这把"哪枚日期真的挨着活时钟"的合同收缩到 3 份文件里逐一读码，是 8.2–8.4 定性成立的依据 |
| 覆盖面 | `grep -rEl` 三种形状任一 | **38 个文件**含日期形状 | 38 文件全部过一遍上面 5 形，无遗漏 |

**分类结论（本程实际量到的）**：

- **(1) 炸弹：1 枚用例／3 枚字面量** —— 就是本票修的 `TestRollingWriterRetentionSweep`
  （`:190` `20260101`、`:198` `20260918` 两枚夹具名，加 `:194` 那枚配套 mtime `time.Date(2026,1,1,…)`）。
  **除它之外，全仓 `internal`/`cmd`/`tools` 再没有第二枚炸弹。**派单"一枚写死的日期不会是唯一一枚"这一条**没有兑现**：
  这个形状确实成族（见下"同族病"），但成族的是**同一个函数的历史病例**，不是散落的兄弟。
- **(2) 安全：其余全部**（紧凑 3 枚 + ISO 12 枚 + `time.Date` 4 枚 + `time.Unix` 13 枚 = **32 枚站点**，逐条见 8.3/8.4）。
- **(3) 已过期但仍绿：0 枚。** 判据是读码给出的、不是"它现在是绿的"给出的：
  一枚夹具要落进 (3)，必须"它已经出窗，而断言恰好因此变得空洞"。全仓只有 `sweepLocked`（`logging.go:317`）
  与 `memory/retention.go` 两套日期淘汰逻辑，后者**已经**是相对形（8.5），前者只有本票这一枚用例在测——
  ⇒ 没有第三处日期淘汰码可被写空，(3) 天然为 0。
- **⚠ 同族病（已修，但不是"兄弟"）**：`TestRollingWriterDayRoll` 的注释（现读 `logging.go:176-180`）逐字记着
  它当年正是**同一枚病**——"green only on the two dates its fixture hardcoded"，修法是那枚 `newRollingWriterClock` 缝。
  ⇒ 本票的修法不是新发明，是把同一剂药打到这个函数剩下的那一枚用例上。**这也解释了为什么只剩 1 枚**：这一族早被清过一轮。

### 8.2 紧凑日期逐枚（修后剩下的 5 枚）

| 站点 | 值 | 类 | 为什么（读码，非猜测） |
|---|---|---|---|
| `logging.go:33` | `20060102` | (2) | Go 的**参考布局**（2006-01-02 是 Go 规定的那个绝对时刻），不是日期；它定义 `logDayLayout` 本身 |
| `diagnostics_test.go:24` | `20260919` | (2) **但见下方"潜伏"警告** | `BuildDiagnosticsBundle` 抄日志那一段（`diagnostics.go:99-130`）**没有任何日期谓词**：它按前缀/后缀挑文件、只按 `MaxLogBytes` 跳过大文件。`o.Now` 只喂 `manifest.json` 的 `created_at`（`:69-71`、`:88`），而用例不断言那一句 ⇒ 文件名里写哪一年都不承重 |
| `logging_test.go:183`（两枚） | `20260919`/`20260920` | (2) | DayRoll 的**期望渲染值**，两侧都由注入时钟 `day1`/`day2`（`:159-160`、`:172`）生成；窗内两枚文件相对注入时刻同向，且 boot sweep 用的是注入时刻 ⇒ 墙钟再走也不翻 |
| `logging_test.go:193` | `20260918` | (2) | **是注释**——我在 §2 里复述那枚死掉的常量做溯源。惰性的，但**会**被下一轮同样的 grep 命中，所以在这里先自报，别被当成第二枚炸弹 |

> **潜伏项上报（不属本票射程，故不改）**：`diagnostics_test.go:24` 现在是 (2)，
> 但它是**同一枚引信只是还没接上药**：诊断包未来若加"只收最近 N 天日志"这类特性（`BundleOptions.Now` 已经在场，
> 加日期谓词是一行的事），这枚 `20260919` 立刻变成第二枚炸弹，而且**今天这枚的翻脸时刻表它要重走一遍**。
> 我没动它，因为今天那条特性不存在、改了就是给一段没有的行为写守卫（过度工程）。归编排者决定要不要立账。

### 8.3 ISO 日期字面量逐枚（12 枚真代码行）

| 站点 | 值 | 类 | 判据 |
|---|---|---|---|
| `llm/anthropic/adapter.go:90` | `2023-06-01` | (2) | `anthropic-version` 请求头常量，协议规定的绝对值；与墙钟无关 |
| `memory/dao_test.go:413,416,422` | `2026-09-19` ×3 | (2) | `BumpCostDay`/`CostDay` 的 `day` 是**主键参数**：`dao_misc.go:154-166` 是 `WHERE day=?` 的纯查，**不默认今天** ⇒ 任意时刻同结果 |
| `memory/dao_test.go:429` | `2026-01-01` | (2) | 同一枚纯查的**反控制**（断言那行不存在）；它不是"相对窗"，是"没种过的键" |
| `models/downloader_test.go:550` | `tiny-model-2024-01-01` | (2) | 上游产物**名字**的一部分，测的是解包扁平化，不读时钟 |
| `models/manifest_real_test.go:42` | `…-3.3M-2024-01-01` | (2) | 同上，manifest 里的模型键名 |
| `models/minisign_test.go:27` | `timestamp: 2026-09-19` | (2) | 签名 payload 的**注释字段**，签完再验同一段文本；不是时间戳判定 |
| `panel/approval_test.go:148` | `2026-09-21T13:00:00Z` | (2) | 结构体字面量字段，只做 JSON 往返（`Marshal`→`Unmarshal`→`DeepEqual`）。**且 `internal/panel` 的生产码里 `time.Now` 现量 0 处**（grep 全包空），时钟只能从 `NewSnapshot(…, now)` 参数进 ⇒ 永不出窗 |
| `observe/thresholds.go:103` | `2026-09-19` | (2) | `Note:` 字段里记 SLO 裁定日期的**散文**，不是判据。⚠ 而且 `AGENTS.md` §1.1 明写 `thresholds.go` 一字节都不许动 ⇒ **本程对它零改动**；普查扫到它、并**刻意不碰**它，这两件事都要记下来 |
| `tools/d22scan/main.go:370,378` | `2026-09-21`/`2026-09-25` | (2) | 在**字符串里的散文**（"armed by ticket 88 on …"），是 `note:` 文案，不是判据。⚠ `tools/d22scan/**` 此刻由另一枚验收程持有，本程**只读不写**，这里只给它定性 |

> 上表 12 枚逐枚对上 8.1 那笔"18 代码行 − 6 参考布局 = 12"的算术，**枚枚具名、无一项是目测加总**。

### 8.4 `time.Date` 6 枚 / `time.Unix` 13 枚

| 站点 | 类 | 判据 |
|---|---|---|
| `observe/logging_test.go:159` | (2) | DayRoll 的 `day1`，**是注入时钟本身**（`:165` 进 `newRollingWriterClock`），不是"最近的"墙钟替身 |
| `observe/logging_test.go:199` | (2) | 本票新增的 `fixedNow`，同上——它是**锚**，其余两枚夹具由它派生 |
| `agent/prompt_test.go:54,57` | (2) | `TestCachePrefixIsByteStableAcrossTurns` 用两个固定时刻造两 turn，断言**缓存前缀字节相同**（时间必须落在 volatile 尾巴外）。`agent/prompt*.go` 生产码不读 `time.Now`（现量：全 `internal/agent` 只有 `approval.go:38` 的 `SystemClock`、`loop.go:1108` 的 id 发生器、`tools.go:134` 的 EchoProvider 出厂戳，都不在 prompt 组装路径上）⇒ 纯函数固定输入，永不出窗 |
| `panel/composer_test.go:209` | (2) | 传给 `NewSnapshot(…, now)`，`composer.go:83` 只把它渲染成 `GeneratedAt` 字符串；该断言看的是 mode/workspace/附件，不比时间 |
| `perm/store_test.go:104` | (2) | 注入 `Options.Now` 给**审计戳**用：`store.go:283-288` 的 `stamp()` 只在 `now==nil` 才落 `time.Now`。perm 里**没有任何 TTL/过期判定**（`grep 'Now\b'` 非测试文件只命中那三行声明/装配）⇒ 绝对值输入，惰 |
| `time.Unix` 13 枚 | (2) | 三类，全是惰的：①假时钟基值 `time.Unix(0,0)`/`(1700000000,0)`（`approval/fakes_test.go:39`、`prompt_test.go:128,165,228,317`、`llm/matrix_14_2_test.go:67`、`ratelimit_pace_test.go:125,213,252`、`observe/earlylog_130_test.go:97`）——`earlylog` 那枚的用途是"严格递增的可分辨时间戳"（同文件 `:86-88` 注释写明），只断言**顺序**；②生产换算 `memory/dao_providerhealth.go:213,217`（存进库的 unix 秒读回来）；③`secret/ctime_windows.go:24`（NTFS FILETIME→time.Time，纯换算；全 `internal/secret` 无一处把它和窗口比） |

### 8.5 相对形（`now.Add(-N*24h)`）单独一条结论

`memory/retention_test.go:136-138` 是唯一用相对形搭窗内/窗外的族（`-399/-400/-401*day` 对 `CostTTL`），
它做对了一件本票该学的细节：`:117-120` **先播种、再把时钟钉成播种时读到的 `observe.NowWallUTC()`**，
并注入 `RetentionConfig{Now: …}`（`:150`）⇒ 播种与 cutoff 同源自同一枚钉住的时刻，**跟着墙钟走而不是挨着墙钟写死**。
⇒ 类 (2)，**不动**。派单问的"相对形是否跨 UTC 午夜耦合"：这一族不耦合，因为它压根不比真实午夜，它比的是自己钉住的那一枚。
（同文件 `:200-201`、`:238` 用 `time.Now()` 做轮询预算与播种，是另一件事；`:200` 那对
`deadline := time.Now().Add(2*time.Second)` / `for time.Now().Before(deadline)` 是**用墙钟差实现超时**的形状，
与 `AGENTS.md` §1.2 那条禁令同形——**本程未动、也未判定**，因为它在 `internal/memory`，越出派单给我
"stay inside `internal/observe/**`" 的范围；**具名上报，请编排者归位**。它不是日期炸弹，不会在某一天翻，
但它是那一族被 CI 禁的形状。）

---

## 9. 归因：这枚红是**墙钟漂移今天新造的**，不在任何已计数的名册里

派单要的那句判定：**新红，且不是此前被计数的那批失败之一。**三条独立证据链：

1. **名册里没有它。** 全 `docs/` 搜这枚用例名，只命中两处，且**两处都是今天由编排者本人写的**：
   `docs/evidence/s1/136-ac15-rate-r1.md:567`（§5.1，08:1x）与 `docs/reports/pending-and-issues.md:5890`（`A212③`，08:3x）。
   ⇒ 它**从未进过** `test-windows` 那批具名红（台账 `:5404` 的 17 枚、`:5513` 降到 16 枚、`:5743` 提到的 12 枚 FAIL），
   也不在 `A52③`（`:2639`）那批 `test-core` 归族红里——那几处的名单里都没有 `internal/observe` 的任何用例。
   编排者自己在 `A212③` 写的原话是"**不在任何登记里**，我已 08:22 现场复现"，与本程独立核对一致。

2. **覆盖它的那道门是 `test-core`，不是 `test-windows`**（这一条把 `136-ac15-rate-r1.md` §5.1 末尾**明说没查**的那一半查完）：
   ```
   $ grep -n 'observe' .github/workflows/ci.yml          → 空（ci.yml 从不按名点名 internal/observe）
   $ grep -rn 'internal/observe' scripts/ .github/workflows/
        scripts/portable-tests.sh:140   （包清单里）
        scripts/portable-tests.sh:175   （core scope：./internal/observe/...）
   $ awk 定位 → ci.yml:288  bash scripts/portable-tests.sh --scope=core
                ci.yml:458  bash scripts/portable-tests.sh --scope=windows
   $ sed -n '184,189p' scripts/portable-tests.sh   # windows scope 块：proc/secret/config/risk/ball/perm/plugin/cmd/llmrecord，**不含 observe**
   $ awk '/^windows\)/{f="windows"} /^core\)/{f="core"} f && /internal\/observe/{print f": line "NR}'  scripts/portable-tests.sh
        core: line 175                             ← 全文件里 observe 只出现在 core 这一块
   ```
   （我最初在 §9 落笔时把块号写成 `184,188p`；现按 `sed -n '184,189p'` 重跑核对后改正——块到 `:189` 的 `;;` 为止。
   这一处是本程自己的笔误，不是事实差错：`internal/observe` 确实**只**在 `core` 作用域里，上面那发 `awk` 是全文件扫描给的旁证。）
   另一条 `./...` 型步骤（`ci.yml:81`）是 `runtests.sh -C tools/d22scan ./...`，**那是 d22scan 自己那个 module**，
   摸不到 `internal/observe`。⇒ **全 CI 只有 `test-core` 一步覆盖它**，`slo-*`/`lint*` 都不跑 Go 测试。

3. **`test-core` 是今天从绿翻红的**：`ci.yml:224-225` `test-core: runs-on: ubuntu-latest`，
   而台账 `:5407`（对已推 run `35967768017` 的读数）明记 **`test-core`/`slo-smoke`/`slo-full`/`lint-frontend` 绿**、红的只有 `lint`/`test-windows`。
   ⇒ 后果比"又一枚存量红"更重：**一枚本来绿的 job，因为日期走到 2026-09-25 00:00 UTC，在没有commit、没有改动的前提下翻红**。
   翻脸分钟数与编排者一致：`136-ac15-rate-r1.md` §5.1 量到 A08 批（07:46–07:57）零枚 `--- FAIL`、A08b（08:00 起）发发红，
   两批之间 `git diff beac693..HEAD -- internal/observe/` 全程为空——**这条是归因最硬的一块**，本程不重复取它。

> ⚠ 一处**平台差**值得记给台账，免得将来对不上数：本机（Windows）红点是 `logging_test.go:216`（旧行号）的
> `fresh file wrongly pruned`，而 M1/M4 那两发变异让我量到——**同一段生产逻辑在 Windows 上会先撞上活文件锁、
> 把红点挪到 sweep 的错误返回**（本表 §5）。`test-core` 跑在 ubuntu-latest，Linux 上 `os.Remove` 对打开文件不报错，
> 所以 **CI 上的红点会是那条干净的 `fresh file wrongly pruned` 断言**，与本机形态不同、结论相同。

---

## 10. 总判

**改了什么**：`internal/observe/logging_test.go` 里 `TestRollingWriterRetentionSweep` 一枚函数体（`+23/−5`）。
两枚夹具名从墙钟字面量（`20260101`／`20260918`）改为**由注入常量 `fixedNow` 派生**（`-12d` 窗外／`-2d` 窗内），
构造子从 `newRollingWriter` 换成 `newRollingWriterClock(dir, 0, 7, func(){return fixedNow})`；配套 mtime 同改为派生日。
**生产码零改动**——那枚缝本来就够（`ensureFileLocked`/`sweepLocked`/`Close` 全走 `w.now()`，用例不起 `flushLoop`）。
**断言一条没松**：仍是"窗外必删／窗内必活"两条相反断言；无 Skip；`days=7`、`retentionSweepPeriod`、`thresholds.go`、
golden、`tools/d22scan/allowlist.txt` 一律未动。

**兄弟普查的数**（现量，非目测）：**炸弹 1 枚（即本票那枚，含 3 枚字面量）；安全 32 枚站点；已过期但仍绿 0 枚。**
扫过的形状与枚数：紧凑 `YYYYMMDD` 6→5 枚（内含 1 枚 Go 布局 + 1 枚我自己的注释复述）、ISO 48 行命中拆出 **12 枚真字面量**、
`time.Date` 6 枚、`time.Unix` 13 枚、测试里的相对 `Add(-N*24h)` 0→2 枚（都是我的）、含日期形状的文件 **38 个**。
派单"一枚写死的日期不会是唯一一枚"这一条**没有兑现**：成族的只有同一函数的历史病例 `TestRollingWriterDayRoll`（已由 `newRollingWriterClock` 治过），
本票是把同一剂药打到该函数剩下那一枚用例上。另上报 **1 枚潜伏**（`diagnostics_test.go:24` 的 `20260919`，今天惰性，
一旦诊断包加"只收最近 N 天"就复燃）+ **1 枚形状越界**（`memory/retention_test.go:200-201` 的墙钟差超时，越出本程 `internal/observe/**` 射程，未动未判）。

**四数（同命令 `-count=2 -v`，两趟都是 71 枚 distinct）**：
修前 `rc=1 / RUN=142 / PASS=140 / FAIL=2 / SKIP=0 / ^panic:=0` →
修后 `rc=0 / RUN=142 / PASS=142 / FAIL=0 / SKIP=0 / ^panic:=0`。distinct 名册 **71→71，未缩、未增**；
两趟 `^panic:` 都是 0，所以本包这次**不存在**"panic 吞掉兄弟读数"的缺口（与 `A212` 那枚空切片形状不是同一件事）；
修后与编排者 09-24 22:52 记录的 142/142/0/0 基线（`:5767`）逐字对齐。
`gofmt -l` 对我碰过的文件为空，`go vet ./internal/observe/` 干净。

**我没测什么（不许读成"已测"）**：
1. **没在真实未来的墙钟上跑过**。未来免疫靠两件事论证——结构（注入常量覆盖该路径全部时钟面）+ 那发把 `fixedNow` 挪到 **2031** 仍绿的实证；
   我**没有**、也无法在本机改系统时钟来复现 2026-09-26 之后。
2. **没跑 CI**。§9 的 job 归属是把 `ci.yml`／`portable-tests.sh` 的射程**读码**读出来的，不是 `gh run view` 量出来的；
   本程按派单**只 commit 不 push**，`test-core` 翻绿要等编排者推。Linux 侧红点形状（§9 那条 ⚠）同理属**推演**，未实测。
3. **没测 `flushLoop` 里那枚每小时的定期 sweep**（`logging.go:301`）——它由 `w.now()` 驱动、需要 registry/计时器才走得起来，
   本票这枚用例从来不调它，我也没为它新增用例（不在派单射程）。**"retention 在长跑进程里是否正确"这条，本包今天仍无覆盖**。
4. **没重跑另外 4 形之外的存储侧日期逻辑**：`memory/retention.go` 我只定性到"相对形＋播种后钉时钟＝安全"，
   **没有**跑 `internal/memory` 的测试来证明（越界，且 `A212⑤` 明说本包与 `AC#15` 那程共享工作树，越跑越野）。
5. **没测 `A212④` 那枚 `AC#15` 抖动**：修后两趟整包 0 FAIL，所以我**没**撞见它，也**不能**据此说它好了或坏了——
   它的 4/8200 命中率不是本包 `-count=2` 量得出的样本（8200 发才是）。按 `A212⑤` 的顺序，它在**本票落地之后**才派。
6. 变异那 5 发是**同一枚 sweeper 的 5 种打断**，不等于"sweepLocked 的所有可能错误都能被这枚用例抓到"；
   我只证了派单点名的那一形（永远删／永远留／窗缩／无视注入时钟），**没做穷举变异**。
