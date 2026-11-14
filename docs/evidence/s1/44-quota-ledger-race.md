# 44-quota-ledger-race — `QuotaLedger` 无锁并发：复现、修法、修后读数

- **腿**：票 44 家族补丁腿（写腿），范围＝CI 已在 `Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`
  复现的那枚数据竞争（证据件 `docs/evidence/s1/70-ci-red-triage.md` step 10 (b)）。
- **base**：`48457eafe`（= `git rev-parse HEAD`，动笔时现取）。
- **修前读数**：`go test -race` 39 条 `WARNING: DATA RACE` ＋ 1 条
  `fatal error: concurrent map iteration and map write`，rc=1（§1 逐字）。
- **修后读数**：同一枚并发测试 `-race -count=5` 10/10 PASS；5 个独立进程 5/5 rc=0；
  包级 `go test -race -count=1 ./internal/observe/` = `ok … 33.079s`（§4 逐字）。
- **写面核**：只落 `internal/observe/quota_44.go`（非测试产物）、新增
  `internal/observe/quota_44_race_test.go`（测试）、本件、票 70 Progress log 一行。
  ⛔ 未碰 `scripts/**`、`cmd/wisp/**`、`internal/panel/**`、`internal/tools/**`、`frontend/**`、
  阈值/golden/契约；⛔ 零断言放宽。只 commit、⛔ push。

---

## 0. 缺陷与判据链（引用层 / 判断层分开）

**引用层（逐字，CI 日志）**：`docs/evidence/s1/70-ci-red-triage.md` step 10 (b) 抄的
`job-test-windows.log:7821` 尸检：

```
level=ERROR msg="goroutine panic recovered" goroutine=agent-task-11289555-… owner=agent
  recovered="runtime error: invalid memory address or nil pointer dereference"
  … observe.(*QuotaLedger).roll(0x31b54c228900)  internal/observe/quota_44.go:256
     observe.(*QuotaLedger).Spend(…)              internal/observe/quota_44.go:284
     agent.(*QuotaGate).Decide(…)                 internal/agent/quota_44r4.go:155
```

**判断层（本腿现盘核的，可逐行复对）**：修前 `quota_44.go` 里
`:251-259` 是 `roll()`，`:284` 是 `Spend()` 里那一次 `l.roll()`，`:270` 是 `Post()` 里那一次；
三处**都没有任何 `sync.`**（修前全文件 `Select-String -Pattern 'sync\.'` = 0 命中）。
CI 点的 `:256` 正是 `if d != l.dayKey {`——即 `roll` 里的**日期戳比较**；下一行 `:257`
才是它旁边的写。**⇒ 这条竞争不是"某处忘了加锁"，是"整个类型没有任何锁"**。

---

## 1. 复现（⛔ 先于任何改动，修前 HEAD）

**仪器**：新增 `internal/observe/quota_44_race_test.go`，两枚用例：

| 用例 | 形状 | 为什么是这个形状 |
|---|---|---|
| `TestQuotaLedgerConcurrentSpendAcrossBucketRolls` | **16 goroutine × 200 轮**纯 `Spend()`；注入时钟**每读一次前进一天** | 产品里 `QuotaGate.Decide` 每次派发前调 `Spend()`（`internal/agent/quota_44r4.go:155`），而 `roll()` 一旦跨日就是**写者**（换 stamp ＋ `clear` 三张 map）。时钟每天一变，让每一次 `roll` 都有桶要换 ⇒ 精确复刻 CI 的"写着清 map / 读着 range map"这一对 |
| `TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro` | 8 读 `Spend()` ＋ 8 写 `Post()`，各 200 轮，时钟**固定**在同一天 | 这是产品的稳态形状：`Decide`→`Spend` 与 `Book`→`Post` 本来就同时在跑（根代理＋子代理）。固定时钟让 `roll` 走"无变化"分支，把共享面收敛到三张 map 上；末尾断言**每一微单位都记账到位**（丢一次更新也红） |

**命令与读数（逐字）**：

```
$ go test -race -run 'TestQuotaLedgerConcurrent' -count=1 -v ./internal/observe/
=== RUN   TestQuotaLedgerConcurrentSpendAcrossBucketRolls
==================
WARNING: DATA RACE
Read at 0x00c000111318 by goroutine 15:
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).roll()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:251 +0x76
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).Spend()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:284 +0x3e
  github.com/CarlosShao/wisp/internal/observe.TestQuotaLedgerConcurrentSpendAcrossBucketRolls.func2()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44_race_test.go:58 +0xa9

Previous write at 0x00c000111318 by goroutine 25:
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).roll()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:252 +0xa7
  github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).Spend()
      D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:284 +0x3e
  …
==================
```

（修前 `:251` = `if m != l.monthKey {`，`:252` = `l.monthKey = m`。）

**同一发里的四数与收尾（逐字）**：

```
    testing.go:1865: race detected during execution of test
--- FAIL: TestQuotaLedgerConcurrentSpendAcrossBucketRolls (0.13s)
=== RUN   TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro
…
fatal error: concurrent map iteration and map write

goroutine 44 [running]:
internal/runtime/maps.fatal({0x1404b7a3b?, 0xc00006fd88?})
	D:/work/base/go/src/runtime/panic.go:1195 +0x18
internal/runtime/maps.(*Iter).Next(0xc00079ff00?)
	D:/work/base/go/src/internal/runtime/maps/table.go:882 +0x85
github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).Spend.func1(...)
	D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:287
github.com/CarlosShao/wisp/internal/observe.(*QuotaLedger).Spend(0xc000111340)
	D:/work/workspace/projects plans/Wisp/internal/observe/quota_44.go:295 +0x39a
github.com/CarlosShao/wisp/internal/observe.TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro.func2()
	D:/work/workspace/projects plans/Wisp/internal/observe/quota_44_race_test.go:92 +0x9a
FAIL	github.com/CarlosShao/wisp/internal/observe	0.777s
FAIL
```

**census（数出来的，不是感觉）**：

| 读数 | 值 |
|---|---|
| 进程退出码 | **rc=1** |
| `WARNING: DATA RACE` 条数 | **39**（第一枚用例 13 条，`--- FAIL` 之后第二枚 26 条） |
| `fatal error: concurrent map iteration and map write` | **1**（第二枚用例把测试进程打死） |
| 四数（`-v` 数） | `RUN=2  PASS=0  --- FAIL=1  SKIP=0` |
| 被点到的修前行号分布 | `:251×2 :252×4 :253×5 :254×5 :256×2 :257×4 :258×5 :270×5 :274×11 :275×8 :277×9 :284×22 :286×9 :287×12 :293×7 :294×6 :295×8` |

**与 CI 的对应关系（这是"复现"而不是"另一条形状相近的红"）**：CI 的那一帧是
`quota_44.go:256`，本仪器**同样点到 `:256`**（2 条），并且 `roll` 的六个访问点
（`:251/252/253/256/257/258`）**全覆盖**；CI 的下一帧是 `Spend`（`:284`），本仪器 `:284` 22 条。
差异只有一处并已交代：CI 里是"torn read 后 nil 解引用 + 任务 recover 成 exit 1"，
本机是"race detector 逐条报 ＋ 第二枚用例直接 `fatal error`"——**同一个缺陷的两种表现**。

---

## 2. 根因（可核层 / 推断层分开）

**可核**：`QuotaLedger`（修前 `:221-228`）六个字段全部裸暴露；`roll()` 在**每次**
`Spend()`/`Post()` 上执行（`:270` / `:284`）；跨日/跨月分支里它既写 `dayKey`/`monthKey`
**又** `clear()` 三张 map；`Spend()` 随后 `range` 这三张 map 做快照（`:285-297`）；
`Post()` 直接 `+=` 写这三张 map（`:274-277`）。⇒ 三类同时存在且**无任何互斥**：
① stamp 的读/写（`:251` vs `:252`、`:256` vs `:257`）；② `range` vs `clear`（`:286/287` vs `:253/258`）；
③ map 的 `+=` vs `range`（`:274/275/277` vs `:286/287`）。

**推断层（标明）**：CI 那条 `nil pointer dereference` 是 **torn read 的具体表现**——
两个 goroutine 同时读写同一个 `string` 头（指针字 + 长度字），读者可能拿到
"指针已更新、长度还是旧的"这种组合；本机 `-race` 抓到的正是这两个字
（`:251/252` 这一对）。**我没有证据说 CI 那个 nil 到底是 `l` 还是 `l.now`**（两条构造路径都有 nil
兜底：`NewQuotaLedger` / `NewQuotaGate`），此层保持 triage 的原话，不编排量不到的结论。

---

## 3. 修法：`sync.Mutex` 保护整个类型（含两支 stamp）

**改了什么**（全在 `internal/observe/quota_44.go`，行为语义零改动）：

| 位置 | 改动 |
|---|---|
| struct | 加 `mu sync.Mutex`，doc 里写明"guards every field below"；六个字段名字/类型/顺序不变 |
| `roll` | 原体搬进 `rollLocked()`（**要求持锁**），`roll()` 变成"上锁→`rollLocked()`→解锁"的薄壳——因为 `SeedFrom`（`internal/observe/costmodel_44b.go:187`）**在本腿写面之外**直接调用 `roll`，改名会编译红（§5） |
| `Post` | 负数检查**留在锁外**（不碰状态，错误文案逐字不变）；`mu.Lock/defer Unlock` 之后 `rollLocked()`，其余 `+=` 语句逐字不变 |
| `Spend` | `mu.Lock/defer Unlock` 之后 `rollLocked()`，快照拷贝**在锁内**完成；返回值形状不变 |

**为什么选互斥而不是"不可变快照 ＋ 原子替换"（理由，按主流 harness 的取舍）**：

1. **读路径本来就要拷三张 map**。`Spend()` 的返回值 `QuotaSpend` 是深拷贝（修前 `:285-297` 已有），
   所以 `Spend` 的成本是 O(桶数) **与用哪种发布方式无关**。原子替换只能省掉"取锁"这一下
   （纳秒级），买不到任何真实开销下降——**复杂不加分，买不到价值就不选它**（机主铁律 3）。
2. **原子替换会把成本搬到写路径**。要让读者无锁读，`Post` 每次就得复制三张 map 再整体换指针
   （每派发一次 3 次分配＋GC 压力），而 `Post` 与 `Spend` 同频（一次派发各一次）；
   互斥方案下 `Post` 仍是原地 `+=`，零新增分配。
3. **`roll` 的语义是"两个字段 ＋ 三张 map"必须一起换**（跨月时 `monthKey` 与两张月 map 同时变）。
   互斥让这个"一起"天然成立；原子方案得把整个状态装进一个 struct 做 CAS，
   **同一个不变量换来的是更多机械**，且 `Spend` 的判定会因为快照粒度不同而看到"半换"的中间态
   ——这正是我要消灭的那类形状。
4. **不构成串行瓶颈（这是"不许为省事把 Spend 变瓶颈"的答复）**：临界区里只有内存 map 读写，
   **无 I/O、无时钟等待、无锁外依赖**；锁是**每派发一次**，不是每模型一次；
   而一次派发本身是一次网络往返（毫秒～秒级），锁的持有时长是微秒级拷贝——
   量级差 3～6 个数量级，**它不可能成为吞吐天花板**。真正每模型一次的是 `JudgeModel`，
   它读的是**已经拷好的快照**，仍然零锁（`QuotaParams` 没有被加上任何锁）。

**行为语义零改动（逐条对位，供复核）**：

| 语义 | 修前 | 修后 |
|---|---|---|
| 负数支出拒收 | `Post` 开头 `micros < 0` 返回同一句错文案 | 逐字不变（只是留在锁外） |
| `micros == 0` 行为 | 先 `roll` 再直接返回 nil | 逐字不变（`rollLocked` 仍在返回前） |
| 80/100 判定、Binding、Remaining、Estimated、Notice | `JudgeModel`/`SelectChain`/`quotaStopNotice`/`quotaWarnNotice` | **一字未动**（本腿没碰这些函数） |
| 桶换日/换月即清空 | `roll` 内 `clear` 三张 map | 逐字不变（`rollLocked` 原体搬运） |
| 可见字段/导出符号 | `QuotaModel`/`QuotaProvider`/`QuotaSpend`/`QuotaVerdict`/`QuotaParams` | **一个都没加/删/改名**；新增的 `mu` 是未导出字段 |
| 导出方法集 | `Post`/`Spend`/`SeedFrom` | 签名、返回值、错误分类**全部不变** |

---

## 4. 修后读数

**4.1 同一枚并发测试，`-race` 下 5 发（`-count=5`）**

```
$ go test -race -run 'TestQuotaLedgerConcurrent' -count=5 -v ./internal/observe/     # rc=0
--- PASS: TestQuotaLedgerConcurrentSpendAcrossBucketRolls (0.05s)
--- PASS: TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro (0.03s)
--- PASS: TestQuotaLedgerConcurrentSpendAcrossBucketRolls (0.02s)
--- PASS: TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro (0.02s)
--- PASS: TestQuotaLedgerConcurrentSpendAcrossBucketRolls (0.03s)
--- PASS: TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro (0.02s)
--- PASS: TestQuotaLedgerConcurrentSpendAcrossBucketRolls (0.03s)
--- PASS: TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro (0.02s)
--- PASS: TestQuotaLedgerConcurrentSpendAcrossBucketRolls (0.04s)
--- PASS: TestQuotaLedgerConcurrentSpendAndPostKeepEveryMicro (0.03s)
PASS
ok  	github.com/CarlosShao/wisp/internal/observe	1.480s
```

`WARNING: DATA RACE` = **0**。

**4.2 五个独立进程（更严：每发新建测试二进制）**

```
run1 rc=0 races=0 :: ok  github.com/CarlosShao/wisp/internal/observe 1.183s
run2 rc=0 races=0 :: ok  github.com/CarlosShao/wisp/internal/observe 1.399s
run3 rc=0 races=0 :: ok  github.com/CarlosShao/wisp/internal/observe 1.162s
run4 rc=0 races=0 :: ok  github.com/CarlosShao/wisp/internal/observe 1.503s
run5 rc=0 races=0 :: ok  github.com/CarlosShao/wisp/internal/observe 1.291s
green=5/5
```

**4.3 包级（门禁要求那一发）**

```
$ go test -race -count=1 ./internal/observe/      # rc=0
ok  	github.com/CarlosShao/wisp/internal/observe	33.079s
```

`WARNING: DATA RACE` = **0**。（33s 是包内 sampler/goroutine 那些既有用例的体量，不是本腿新增的。）

**4.4 定向 `Test197`（派单指定的命令，逐字给四数）**

```
$ go test -count=3 -run Test197 -v ./internal/agent/     # rc=0
time=2026-11-14T12:10:36.526+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
testing: warning: no tests to run
PASS
ok  	github.com/CarlosShao/wisp/internal/agent	0.173s [no tests to run]
四数：RUN=0  PASS=0  FAIL=0  SKIP=0
```

**⇒ 这个读数要如实说清楚**：`internal/agent` **没有** `Test197` 家族的用例
（全仓 `func Test197` 只出现在 `internal/tools/subagent_197_test.go` 与
`cmd/wisp/subagent_selfapproval_197_test.go`）；所以这一发是 **RUN=0**（不是"跑了且绿"）。
**补一发同族真读数**（不是派单要求、但同一家族、且不碰 DLL/夹具）：

```
$ go test -count=3 -run Test197 -v ./internal/tools/     # rc=0
四数：RUN=33  PASS=33  FAIL=0  SKIP=0
ok  	github.com/CarlosShao/wisp/internal/tools	0.122s
```

CI 里真正红的那一枚在 `cmd/wisp`（`subagent_selfapproval_197_test.go:201`）——它要
`scripts/wisp-cli-staged-tests.sh` 的 DLL 暂存＋voice 夹具（triage step 4/13 已量过：
夹具之前 `go test -list ./cmd/wisp/` 就 `exit status 0xc0000135`）⇒ **本腿未跑那一枚，不拿它冒充实测**。

**4.5 门禁四门**

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./internal/... ./cmd/... ./tools/...` | **rc=0** |
| vet | `go vet ./internal/observe` | **rc=0** |
| gofumpt | `gofumpt -l internal/observe`（本机 v0.12.0，CI pin v0.7.0） | **输出为空、rc=0** |
| d22 | `sh scripts/d22scan.sh` | 见 §4.6 |

**4.6 d22 的两次读数（第一次红在我自己身上，如实留档）**

第一次（证据件还没落盘时）`go run . -root <repo>` 恰好只抓一条，**抓的就是本件**：

```
internal/observe/quota_44.go:223: [phantom-citation] comment cites repo path
"docs/evidence/s1/44-quota-ledger-race.md" which does not exist on disk (ticket 212 ban #9):
fix the citation or create the file; shorthand forms are a separate prescription, not a violation
exit status 1
```

⇒ 这是 ban #9 按设计工作（引用即导航坐标，落盘前不能引用）。同一次 `sh scripts/d22scan.sh`
除这条真 finding 外还红了 **2 枚自测**（正控 `runtests.sh` = `PASS=43 FAIL=2`）：
`TestRealRepoLedgerIsHonest`（它复跑真扫、rc=1，红因就是上面这条引用）与
`TestScannerSelfScanOfRealRepoIsGreen`（死在
`open …\internal\winsec\wisp129-dr-21912: The system cannot find the file specified`——
那是**别人同一个工作树里正在跑的用例留下的瞬时路径**，我复跑时该路径已不存在；
不归因、只记读数）。本件落盘后复跑见 §7：**两枚都转绿、rc=0**。

---

## 5. 遗留（具名，⛔ 本腿没做，也不假装做了）

**`SeedFrom` 仍是一个不持锁的写者**：`internal/observe/costmodel_44b.go:186-201` 先调
`l.roll()`（现在这一发会拿锁）、**然后在锁外**直接 `l.modelMon[key] += …` / `l.modelDay[key] += …`
/ `l.planMon[r.Provider] += …`。它**不在本腿写面**（写面只有 `quota_44.go` ＋ 测试），所以本腿
**不能**把它包进 `l.mu`；能做的都做了：`roll()` 保留为"上锁薄壳"使它不编译红、
把这段直写具名在此、并确认它不在本红路径上。

**为什么它不是本红的原因（可核）**：它的唯一调用者是 `QuotaGate.SeedBaseline`
（`internal/agent/permodel_44b.go:93-110`），注释与票面口径都是"**seed once at boot before any
task books**"——即任务起来之前那一次；CI 的尸检栈里也没有 `SeedFrom`/`SeedBaseline` 帧。

**交给 owner 的一行修法**（需要放行 `costmodel_44b.go`，或由该文件的所有者补）：
在 `SeedFrom` 函数体最外层加 `l.mu.Lock()` / `defer l.mu.Unlock()`，并把首行 `l.roll()` 改成
`l.rollLocked()` 即可；不改任何语义。**本腿不越界去改。**

---

## 6. 复跑命令（谁都能重来）

```sh
# 1) 复现（在修复前的树上）：应看到 WARNING: DATA RACE 与 rc=1
go test -race -run 'TestQuotaLedgerConcurrent' -count=1 -v ./internal/observe/

# 2) 修后：5 发全绿
go test -race -run 'TestQuotaLedgerConcurrent' -count=5 -v ./internal/observe/

# 3) 包级
go test -race -count=1 ./internal/observe/

# 4) 门禁
go build ./internal/... ./cmd/... ./tools/...
go vet ./internal/observe
gofumpt -l internal/observe
sh scripts/d22scan.sh

# 5) 复现仪器的"红方向"自证（在旧实现上非法硬化只有一种方式：把 mu 去掉重跑，
#    本腿不做临时改动，用 §1 的落盘日志复对；日志保留在
#    %TEMP%\wisp44-race-pre-full.log）
```

---

## 7. 最终门禁读数与写面读数（落盘后现取）

- `sh scripts/d22scan.sh`：**rc=0 / clean**。逐字两段：
  - 正控：`runtests.sh: OK - packages=[./...] top-level: PASS=45 FAIL=0 SKIP=0, === RUN=93,
    '[no tests to run]'=0`
  - 真扫：`d22scan: examined 481 production Go files under internal/ and cmd/ of
    D:/work/workspace/projects plans/Wisp` … `d22scan: clean - no D22 ban violations`（0 finding，
    ban #8 emoji 覆盖 `internal/ 1008`＋`cmd/ 278` 个 Go 文件，含 `_test.go` ⇒ 新增的测试文件也在射程内）。
- `gofumpt -l internal/observe`：**空输出、rc=0**（本机 `v0.12.0 (go1.27.1)`；CI pin 的是 v0.7.0，
  ticket 70 的进度日志记过这两个版本在同一树上给过同样的读数）。
- `git status --porcelain -- internal/observe`：只有 `M internal/observe/quota_44.go`
  ＋ `?? internal/observe/quota_44_race_test.go`。
- 未 push；commit 用显式 pathspec（`git commit -- <三枚路径>`，⛔ add -A / amend / reset / rebase /
  stash / checkout .）。
