# 44-quota-seedfrom-lock — `SeedFrom` 锁外直写三张 map：正控、修法、修后读数

- **腿**：票 44 家族补丁腿（写腿），**只收**上一腿 `6f23ce632` 具名留下的那一枚尾巴
  （`docs/evidence/s1/44-quota-ledger-race.md` 的「具名遗置」：`SeedFrom` 仍在锁外直写三张 map）。
- **base**：`6f23ce632`（= 动笔时 `git log --oneline -3` 的首行之外那一枚；本腿起点）。
- **修前读数**：新仪器在 `-race` 下**4 发全红**（rc=1），每发 6–7 条 `WARNING: DATA RACE`，
  fatal 0；竞争对＝`SeedFrom` 的 `costmodel_44b.go:193/195/198` 写 ↔ `Spend` 的
  `quota_44.go:333/334` 读（§1 逐字）。
- **修后读数**：同一枚仪器 `-race` **5 个独立进程 5/5 rc=0 races=0**；同进程 `-race -count=5`
  **5/5 PASS**；包级 `go test -race -count=1 ./internal/observe/` = `ok … 25.134s`（§3 逐字）。
- **写面核**：只落 `internal/observe/costmodel_44b.go`（非测试产物）、
  `internal/observe/quota_44_race_test.go`（扩既有仪器）、本件、票 70 Progress log 一行。
  ⛔ 未碰 `scripts/**`、`cmd/wisp/**`、`internal/panel/**`、`internal/tools/**`、`frontend/**`、
  阈值/golden/契约；⛔ 零断言放宽；⛔ 行为语义零改动。只 commit、⛔ push。

---

## 0. 缺陷（引用层 / 判断层分开）

**引用层（上一腿自己具名的）**：`internal/observe/quota_44.go:289-296` 的 `roll` 注释逐字：

```
// roll is rollLocked with the lock taken, for the one writer whose call site
// lives outside this file: SeedFrom (internal/observe/costmodel_44b.go) opens
// with l.roll() and then books rows into the maps itself. That seed is
// documented as a boot-time call made before any task books
// (QuotaGate.SeedBaseline), so it is not on the concurrent path this leg
// closes; its own map writes are a NAMED open item in
// docs/evidence/s1/44-quota-ledger-race.md, because protecting them needs one
// line in costmodel_44b.go, which is outside this leg's write face.
```

**判断层（本腿现盘核的）**：修前 `SeedFrom` 的体是

```go
l.roll()                                   // 取锁→rollLocked→放锁
for _, r := range rows {
    ...
    l.modelMon[key] += r.Micros            // :193  ┐
    l.modelDay[key] += r.Micros            // :195  ├ 全程锁外
    l.planMon[r.Provider] += r.Micros      // :198  ┘
}
```

而 `Spend` 的三张 map 拷贝**在 `l.mu` 里**（`quota_44.go:332-338`，读点 `:333/:334`）。
⇒ 锁的存在只给了读者「彼此一致」的保证，**没给读者与 `SeedFrom` 之间任何 happens-before 边**：
一写一读，同一块 map 内存，无同步。这正是本腿要关的那一对。

---

## 1. 正控（⛔ 先于任何改动，修前 HEAD）

**仪器**：扩既有 `internal/observe/quota_44_race_test.go`，新增
`TestQuotaLedgerSeedFromOverlapsSpend`：

| 形状 | 数 | 为什么是这个形状 |
|---|---|---|
| `Spend()` 读者（range 三张 map 的拷贝） | 4 goroutine × 200 轮 | 产品里 `QuotaGate.Decide` 每次派发前调 `Spend()`（`internal/agent/quota_44r4.go:155`） |
| `SeedFrom(rows, currentDay, cards)` 写者 | 1 goroutine × 40 轮 | 产品里 `QuotaGate.SeedBaseline`（`internal/agent/permodel_44b.go:103`）在 boot 重放上月 `cost_daily_model` 行 |
| 注入时钟 | 固定在 `2026-11-03 09:00 UTC` | 让 `roll` 走「无变化」分支，把共享面收敛到三张 map 上（stamp 由上一腿已关） |
| seed 行 | ①`2026-11-03` / plan 计费卡 / 7 micros（day＋month＋provider 信用三桶都写）<br>②`2026-11-02` / 无卡 / 11 micros（只写 month） | 三条写语句 `:193/:195/:198` **各自**都要有对手，而不是只打到一张 map |

**这就是重叠的造法**：boot seed 与第一次派发的 `Decide` 本来就可能同时在跑（`SeedBaseline`
在 `permodel_44b.go` 里是同一次启动路径上的调用）。40 轮 seed × 200 轮读，同一个固定日，
写者与读者都不带任何会让步的 I/O ⇒ 两者在进程里必然交叠。

**修前读数（4 发，逐发现跑，rc 与条数逐发不同——不是抄一次）**

| 发 | 命令 | rc | `WARNING: DATA RACE` | `fatal error` | 命中帧 |
|---|---|---|---|---|---|
| 1 | `go test -race -count=1 -run TestQuotaLedgerSeedFromOverlapsSpend ./internal/observe/ -v` | 1 | 6 | 0 | seed `193/195/198` |
| 2 | 同上（`-v` 省略） | 1 | 6 | 0 | seed `193/195/198` |
| 3 | 同上 | 1 | 7 | 0 | seed `193/195/198` |
| 4 | 同上 | 1 | 6 | 0 | seed `193/195/198` |

**逐字（第 1 发第一条报告，原文未改，只删了 `created at` 里重复的两段栈）**：

```
==================
WARNING: DATA RACE
Write at 0x00c0000210b0 by goroutine 14:
  runtime.mapassign_faststr()
      D:/work/base/go/src/internal/runtime/maps/runtime_faststr.go:261 +0x0
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).SeedFrom()
      D:/work/workspace/projects plans/Wisp/internal/observe/costmodel_44b.go:193 +0x1d1
  github.com/CarlosShao/wisp/internal/observe.TestQuotaLedgerSeedFromOverlapsSpend.func3()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44_race_test.go:172 +0x14d

Previous read at 0x00c0000210b0 by goroutine 12:
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).Spend.func1()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:333 +0x2ba
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).Spend()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:341 +0x28c
  github.com/CarlosShao/wisp/internal/observe.TestQuotaLedgerSeedFromOverlapsSpend.func2()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44_race_test.go:160 +0xcc
==================
```

同类报告的另一半（读在 `:334`、写者在 `:195`）逐字：

```
WARNING: DATA RACE
Read at 0x00c000022d18 by goroutine 11:
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).Spend.func1()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:334 +0x21d
  ...
Previous write at 0x00c000022d18 by goroutine 14:
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).SeedFrom()
      D:/work/workspace/projects plans/Wisp/internal/observe/costmodel_44b.go:195 +0x2b9
```

收尾三行（第 1 发，逐字）：

```
    testing.go:1865: race detected during execution of test
--- FAIL: TestQuotaLedgerSeedFromOverlapsSpend (0.03s)
FAIL	github.com/CarlosShao/wisp/internal/observe	0.153s
```

**⇒ 正控成立：`SeedFrom` ↔ `Spend` 的重叠造得出来，修前 4/4 出 `WARNING: DATA RACE`。**
（另注：本枚仪器在修前**没有**出现上一腿那种
`fatal error: concurrent map iteration and map write`——4 发 fatal 都是 0。
因为读者是 `for k, v := range` + 写 `dst`，写者只有 1 个 goroutine，命中窗口比上一腿的
16 写者窄；红是检测器报的，不是运行时崩的。）

---

## 2. 修法（⛔ 只改 `costmodel_44b.go` 的这一枚函数）

```go
func (l *QuotaLedger) SeedFrom(rows []CostModelDayTotal, currentDay string, cards map[string]QuotaModel) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rollLocked()
	for _, r := range rows {
		...                                  // 原样，一行未动
	}
	return nil
}
```

**为什么外层取同一把 `mu`、并把 `l.roll()` 换成 `l.rollLocked()`**
（这正是派单点名要的形状）：
`roll()` 自己 `l.mu.Lock()`；若外层先 `Lock` 再调 `roll()`，同一个 `sync.Mutex` 不可重入
⇒ **当场死锁**。所以持锁段内只能调 `rollLocked`（`quota_44.go:274`，契约「调用者持锁」）。

**为什么不选「先 `l.roll()` 再 `l.mu.Lock()` 只锁写」**（形状上也过 `-race`、还不用碰
`roll` 的死活）：那样 roll 与写分处两个临界区，中间可以跨日——`Spend` 在午夜把 day map
`clear` 掉之后，seed 才把"昨天那一行"写进今天刚清空的 day 桶。这与修前**同一枚窗口**，
但既然本腿就是来关窗口的，就没有理由把窗口留着。一次临界区是**行为上严格更好**的那一边，
且与 `Post`/`Spend` 的既有形状（`Lock` → `rollLocked` → 干活）完全同形。

**行为语义零改动（逐项对）**

| 语义点 | 修前 | 修后 |
|---|---|---|
| 每行入桶的等式 | `modelMon[key] += Micros`；`r.Day == currentDay` 时 `modelDay[key] += Micros`；卡存在且 `BillingPlan` 时 `planMon[r.Provider] += Micros` | 一字未改（同一段循环体） |
| 负行拒绝 | `r.Micros < 0` ⇒ `fmt.Errorf("observe: quota ledger refuses negative stored spend (%d micros for %s)", …)` | 同一句、同一时机（roll 之后、入桶之前） |
| 返回值 | 只有负行返回 error，其余 `nil` | 同 |
| 可见字段 | `Spend()` 的三张 map 拷贝口径不变 | 同 |
| roll 时机 | 函数入口一次 | 函数入口一次（同一把锁内） |

既有的语义钉子照旧全绿（未被本腿改动）：`costmodel_44b_test.go`
`TestQuotaLedgerSeedFromBucketsDayMonthPlan`（day/month/plan 三桶 ＋ 负行拒绝）、
`TestQuotaLedgerSeedRolloverIsModelCorrect`（跨月/跨日 roll 语义）——两枚都在
§3 的包级 `-race` 全量里跑过。

仪器尾部还加了**安全半边**（不靠检测器也能红）：seed 40 轮 × (7, 11) micros，
末尾断言 `ModelDay[raceprov/racemodel]=280`、`ModelMonth=280`、`ProviderPlan[raceprov]=280`、
`ModelMonth[seedprov/seedmodel]=440`、`ModelDay[seedprov/seedmodel]=0`
⇒ **丢任意一次更新也红**，不是"检测器安静就算过"。

---

## 3. 修后读数（全部现跑）

**同一枚仪器，5 个独立进程（`-count=1` 各一发）**

| 发 | rc | `WARNING: DATA RACE` |
|---|---|---|
| 1 | 0 | 0 |
| 2 | 0 | 0 |
| 3 | 0 | 0 |
| 4 | 0 | 0 |
| 5 | 0 | 0 |

单发逐字：`ok  	github.com/CarlosShao/wisp/internal/observe	1.167s`

**同进程 5 连发**：`go test -race -count=5 -run TestQuotaLedgerSeedFromOverlapsSpend ./internal/observe/ -v`

```
--- PASS: TestQuotaLedgerSeedFromOverlapsSpend (0.01s)
--- PASS: TestQuotaLedgerSeedFromOverlapsSpend (0.01s)
--- PASS: TestQuotaLedgerSeedFromOverlapsSpend (0.01s)
--- PASS: TestQuotaLedgerSeedFromOverlapsSpend (0.01s)
--- PASS: TestQuotaLedgerSeedFromOverlapsSpend (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/observe	1.161s
```

**包级**：`go test -race -count=1 ./internal/observe/` ⇒
`ok  	github.com/CarlosShao/wisp/internal/observe	25.134s`，rc=0。

**门禁**

| 门 | 命令 | 读数 |
|---|---|---|
| 构建 | `go build ./internal/... ./cmd/... ./tools/...` | rc=0，无输出 |
| vet | `go vet ./internal/observe` | rc=0，无输出 |
| 格式 | `gofumpt -l internal/observe` | 空（0 行） |
| D22 | `sh scripts/d22scan.sh` | 见下，**第一次红在本人刚写的那一行注释上，建件后重跑即 clean** |

**D22 第一次红（如实记）**：本腿先在 `costmodel_44b.go` 的 `SeedFrom` 注释里引了本件的路径，
而本件当时还没落盘 ⇒ 扫描器自己抓到「注释引用不存在的仓库路径」：

```
scan_test.go:269: repo HEAD violates: internal/observe/costmodel_44b.go:190: [phantom-citation]
comment cites repo path "docs/evidence/s1/44-quota-seedfrom-lock.md" which does not exist on disk
(ticket 212 ban #9): fix the citation or create the file; shorthand forms are a separate
prescription, not a violation
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen
--- FAIL: TestRealRepoLedgerIsHonest
```

⇒ 不是别人的文件、也不是放宽判据，就是这件证据件本身还没写。**建件后重跑**（本件落盘之后现取）：

```
d22scan.sh: positive control - runtests.sh -C tools/d22scan ./...
runtests.sh: OK - packages=[./...] top-level: PASS=45 FAIL=0 SKIP=0, === RUN=93, '[no tests to run]'=0
ok  	github.com/CarlosShao/wisp/tools/d22scan	145.146s
d22scan: examined 481 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: clean - no D22 ban violations; ... ban #8 internal/=1009, ban #8 cmd/=278; ...
```

rc=0，正控 PASS=45 / FAIL=0 / SKIP=0，真扫 481 枚生产 Go 文件 0 发现。

---

## 4. 复跑（两条命令，任何人在干净树上重放）

```sh
go test -race -count=1 -run TestQuotaLedgerSeedFromOverlapsSpend ./internal/observe/ -v
go test -race -count=1 ./internal/observe/
```

判据：前者 `ok`、`races=0`；后者 `ok`。修前同一枚前者的读数见 §1。

---

## 5. 具名后果／未做（不越界，如实登记）

1. **`QuotaLedger.roll()` 现在没有任何生产调用点。** 修前它是 `SeedFrom` 的唯一入口；
   本腿按派单把 `SeedFrom` 改成 `rollLocked` 之后，`internal/observe/quota_44.go:297` 的
   薄壳 `roll` 只剩测试以外无人引用。现跑读数（`staticcheck` 2026.2.2，本机）：

   ```
   internal\observe\quota_44.go:297:23: func (*QuotaLedger).roll is unused (U1000)
   ```

   这是 CI `staticcheck` step **pass 2（`GOOS=windows -checks=all`，含 U1000）**会读到的新增
   一条（pass 1 是 `-checks=all,-U1000`，读不到）。**没有在本腿修**，因为修法要动
   `internal/observe/quota_44.go`（删 `roll` 与其上那段已过期的注释），**不在本腿写面**。
   下一位的一句话修法：删 `quota_44.go:289-301` 那 13 行（`roll` ＋ 注释），并把
   `SeedFrom` 的注释里对 `roll` 的指涉对齐（已在本腿写成「持锁段内只能调 `rollLocked`」，
   不依赖 `roll` 存在）。**注意**：`A1773`「不许删 U1000 报的符号」那条裁定针对的是
   **windows 侧活着、linux 侧假死**的符号；这一枚是**两个 GOOS 都死**（`SeedFrom` 是唯一调用点，
   它在平台中立文件里），删它不改变任何平台的可达代码。
2. **`SeedFrom` 的并发产品形状本腿没有断言**（`SeedBaseline` 与第一次 `Decide` 的先后由
   `internal/agent` 的装配决定，不属 `internal/observe` 的可见面）。本腿只断言**内存安全**
   与**种子值不丢**；不声称"产品里一定发生过跨线程重叠"——这不是本腿能看到的层。

---

## 6. U1000 摘壳（2026-11-14，`44-quota-u1000` 腿）

§5.1 把这一枚具名留给了下一位（"删 `quota_44.go:289-301` 那 13 行"）。本节就是那一步的读数。

### 6.1 做了什么（`internal/observe/quota_44.go`，只注释与删除）

1. **删掉** `func (l *QuotaLedger) roll()` 这枚薄壳（`rollLocked` 加锁版）**及其 doc 注释**。
   删前先数调用点：全仓 `.go` 里 `\.roll()` 命中 **6 处 = 1 处定义（`quota_44.go:297`）
   + 5 处注释**，**生产与测试的调用点均为 0**。所以这不是"某平台假死"，是**两个 GOOS 都死**：
   唯一真调用点（`SeedFrom`）在 `756ebbcfa` 已改成 `rollLocked`，而 `SeedFrom` 位于平台中立文件。
   `A1773`（不许删 U1000 报的符号）针对的是 windows 活着、linux 假死的形状，与这一枚不同类。
2. **改写 `rollLocked` 的 doc**，把两处对 `roll` 的指涉换成现状：
   原来写 "`roll` takes it, Post and Spend take it"，现在写 **`Post`/`Spend`/`SeedFrom` 各自持
   `l.mu`**，并明说**没有取锁变体可调**——`SeedFrom` 自己持锁、先调 `rollLocked`、再在同一临界区
   入桶，所以 boot seed 对 `Spend`/`Post` 是一次原子步骤（引本件）。
   ⛔ 没有为了"让它被用到"而加任何调用：U1000 报的是**死符号**，加调用是把账藏起来而不是修它。

### 6.2 读数（全部现跑，本机 staticcheck 2026.2.2 / go1.27）

| 门禁 | 修前 | 修后 |
|---|---|---|
| `staticcheck -checks=all ./internal/observe/`（`GOOS=windows`） | `quota_44.go:297:23: func (*QuotaLedger).roll is unused (U1000)`，**rc=1** | **rc=0，逐字无输出** |
| `staticcheck -checks=all ./internal/observe/ ./internal/agent/ ./internal/tools/` | 未跑（同口径含 U1000） | **rc=0，逐字无输出** |
| `go test -race -count=1 ./internal/observe/` | `ok`（`756ebbcfa` 上） | **rc=0，`ok … 29.560s`** |
| `go build ./internal/... ./cmd/... ./tools/...` | rc=0 | **rc=0** |
| `go vet ./internal/observe` | rc=0 | **rc=0** |
| `gofumpt -l internal/observe` | 空 | **空（0 行）** |
| `sh scripts/d22scan.sh` | clean | **rc=0 clean**（真扫 `internal/`+`cmd/` 481 枚生产 Go 文件，正控在） |

文件行数 384 → **374**（删 13 行，doc 改写净增 3 行；`rollLocked`/`Post`/`Spend`/`SeedFrom`
的函数体一个字未动）。

### 6.3 `staticcheck -checks=all ./...` 为什么仍是 rc=1（不属本腿）

全仓那一发 rc=1，三处命中**全部在 HEAD 上就已存在、且文件与 `HEAD` 逐字节相同**
（`git status --porcelain -- internal/ball frontend/embed.go build/` 输出为空）：

```
build\mut_rvw.go:255:6: undefined: voiceTtsFiles   … (11 处 undefined，compile)
frontend\embed.go:4:1: ineffectual compiler directive due to extraneous space (SA9009)
internal\ball\orb_windows.go:204:7: const orbHwndTopmost is unused (U1000)
internal\ball\renderer_windows.go:822:20: func (*renderer).overlayGlyphPixels is unused (U1000)
```

- `build/mut_rvw.go` **未被 git 跟踪**（`git ls-files build/mut_rvw.go` 空输出），是树上的变异
  实验遗留件，所以 `build` 包在 `./...` 下编不过；`go build ./internal/... ./cmd/... ./tools/...`
  （CI 实际跑的那一发）不含 `build/`，rc=0 不受影响。
- `internal/ball` 那两条正是 `429387f46`（票 122）**具名留下的 Tier-4 残渣**（只在
  `-tags=winlive` 下有读者），不是新账。

### 6.4 具名遗留（都在本腿写面外，纯注释，零编译/门禁影响）

1. `internal/observe/costmodel_44b.go:191` 的注释仍写着 "Taking the lock here and then calling
   `roll()` would re-enter the same mutex and deadlock"——`roll()` 已不存在，这句现在指空。
   建议下一腿改成 "there is no lock-taking variant to call"（1 行、纯注释）。
2. `internal/observe/quota_44_race_test.go:19/36/71` 三处注释仍以 `roll()` 指代换日/换月的动作
   （实际执行者是 `rollLocked`）。三处都不是调用，测试语义不受影响。

⚠ 本腿**没有**碰这两处：写面只有 `quota_44.go`、本件、票 70 的一行 log。
