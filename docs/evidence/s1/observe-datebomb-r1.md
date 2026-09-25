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

`newRollingWriter`（`:172`）仍被 `TestRollingWriterSizeRoll`（现读 `logging_test.go:133`）与生产 `InitLogWithRegistry` 使用，
换掉本用例这一处调用**不产生未使用符号**（vet 已证）。
