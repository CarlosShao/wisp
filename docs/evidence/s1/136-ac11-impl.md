# 票 136 AC#11 —— 实现方证据：`TestNoopTaskReturnsToBaseline` 既有偶发失败（先量复现率 → 归因 → 修 → 复量）

**本格只做 AC#11 那一格。**票面其余各格（AC#8/#9/#10/#12/#13 等）一枚未动。判据原文＝
`.scratch/wisp/issues/136-the-zero-sample-fail-closed-guard-underneath-the-freshness-nail-has-no-nail-54-tests-stay-green-without-it.md:164`
（四条件 ①量复现率并给 n／②归因到"计数窗口缺 happens-before"／③只许把等待变成有判据的等待／④修后 ≥20 发命中 0）。

| 项 | 值 |
|---|---|
| 未修版锚点 | `51e29b044bf928e300c4a48ca1cc0932e02300b4`（本程 `git rev-parse HEAD` 现取，`git cat-file -t`＝`commit`） |
| 修法 commit | `f06a8d0`（只带 `internal/observe/goroutine_test.go` 一枚 pathspec；`git show --name-only` 见 §3 末） |
| 证据 commit | 本文件那一枚（码与证据分两枚，见各节末的原始输出） |
| 改动文件 | **只有** `internal/observe/goroutine_test.go`；生产码 `internal/observe/*.go`（非 `_test.go`）零改动、`internal/risk/**` 未碰 |
| 被验版本 | 全部读数取自 `git archive <sha> \| tar -x` 出的归档树（`/d/tmp/wisp136ac11-tree0`＝51e29b0、`-tree1`＝f06a8d0），**没有一次把仓库工作树当被验版本**；仓内未建 worktree/checkout |
| 工具链 | `go version go1.27.1 windows/amd64`；宿主 Git Bash；`gofumpt v0.12.0 (go1.27.1)`＝`D:\work\base\gopath\bin\gofumpt.exe` |
| 时钟纪律 | 本程全文**没有任何一处用时间戳相减得时长**；每发只报它自己那行的 `date` |

---

## §0 争用闸门（每批发火前，原样）

本机 self-hosted runner 与本机同机（`E:\work\base\actions-runner`，见 `docs/reports/HANDOVER.md:346`）。
闸门脚本：`/d/tmp/wisp136ac11-gate.sh`（v1，四查）与 `/d/tmp/wisp136ac11-gate2.sh`（v2，四查＋进程归因＋`VERDICT=` 一行）。
⚠ `gh` 在本程**间歇性报错**（`EOF`），取不到时按降级口径记 **CLEAR-LOCAL-ONLY**，**不等于"确认为空"**。

**B1＝改前批 1（30 发，`tree0`）── 闸门 `wisp136ac11-gate-before.txt`，14:07:22**

```
--- 1) host processes named Runner.Worker/go/compile/cgo (tasklist) ---        ← 四枚全空（零行输出）
--- 1b) any process whose ExecutablePath is under actions-runner\_work ---
Runner.Listener.exe      3952 E:\work\base\actions-runner\bin\Runner.Listener.exe   ← 常驻监听，不是干活的 worker
--- 2) docker ps ---  union-proxy Up 8 days / clipsync Up 4 days / clipsync-minio Up 12 days (healthy)
                       / clipsync-admin-int Up 2 weeks / clipsync-db Up 13 days (healthy) / clipsync-redis Up 13 days (healthy)
--- 2b) docker stats --no-stream CPU ---  0.00% / 0.32% / 0.00% / 0.00% / 5.53% / 0.46%   ← 无容器在跑 go/编译
--- 3) gh run list ---  failed to get runs: ... EOF  ⇒ 降级 CLEAR-LOCAL-ONLY（未确认远端为空）
--- 4) runner _diag Worker_*.log mtime ---  12:05:21 / 11:04:28 / 10:54:08  ← 距开测 2 小时 02 分无 worker 活动
```
批 1 逐发时刻 14:07:48 → 14:09:37（`wisp136ac11-batch-before/index.txt`）。**反向追认**：§0-B3 那批发 now 才出现的
`Worker_20260924-062600-utc.log`（＝14:26:00 +08 起跑）是本日**唯一**一发 worker，⇒ 批 1 与"CI 抢 CPU"无交集。

**B1c＝改后批（被本程自己判废，30 发，`tree1`）── 闸门 `wisp136ac11-gate-after.txt`，14:14:38**

```
--- 1) ... (tasklist) ---
go.exe                        3984 Console   1   28,700 K
compile.exe                  19168 Console   1    2,992 K
--- 1b) ExecutablePath like *actions-runner* ---  Runner.Listener.exe 3952 （只有这一枚，bin 路径）
--- 3) gh run list --- [{"databaseId":35958260537,"headSha":"c8967b8b…","status":"completed"}, … 共 5 枚，全 completed，零枚 in_progress]  gh_rc=0
--- 4) Worker_*.log mtime --- 12:05:21 / 11:04:28 / 10:54:08
```
⇒ 闸门 ① **命中**（宿主有 `go.exe`＋`compile.exe`）。这一批**按规矩不作数**（`wisp136ac11-batch-after/`，14:14:46→14:16:34，实测 0/30，本文件只当旁证登记）。
⚠ 该 `go.exe`/`compile.exe` **不在** `actions-runner` 路径下（1b 只有常驻 Listener），且 worker 日志 mtime 直到 14:17 仍是 12:05:21
⇒ 它**不是 CI 那一族**；14:17:27 复量时四枚进程已全部消失（`wisp136ac11-gate2-probe.txt`），本程无仪器可归因到具体来源 ⇒ 记为"未归因的本机 go 活动"。

**B2＝改后正式批（30 发，`tree1`）── 有界轮询等空：`wisp136ac11-waitlog-after2.txt`**

```
attempt=1 at=14:18:20 VERDICT=BUSY go=1 compile=0 cgo=0 runnerworker=0 docker_over_25pct=0
attempt=2 at=14:18:44 VERDICT=BUSY go=1 …        attempt=3 at=14:19:08 VERDICT=BUSY go=2 …
attempt=4..6        VERDICT=BUSY go=2 …        attempt=7 at=14:20:47 VERDICT=BUSY go=1 …
attempt=9 at=14:21:33 VERDICT=CLEAR go=0 compile=0 cgo=0 runnerworker=0 docker_over_25pct=0
CLEAR after attempt 9   （上界＝20 次 × 10s；完整闸门留在 …-waitlog-after2.txt.attempt9）
```
批内逐发时刻 14:21:33 → 14:23:2x（`wisp136ac11-batch-after2/index.txt`）；**批后闸门** `wisp136ac11-gate-post.txt`：
`VERDICT=CLEAR go=0 compile=0 cgo=0 runnerworker=0 docker_over_25pct=0`，`### gate end: 2026-09-24 14:23:25 +0800`。
⇒ B2 是**前后都被 CLEAR 夹住**的一批，且整个区间早于 14:26:00 那发 CI ⇒ ④ 的读数出自这一批。

**B3＝改前旁证批（30 发，`tree0`）──  `wisp136ac11-waitlog-before2.txt`**

```
attempt=1 at=14:25:13 VERDICT=CLEAR go=0 compile=0 cgo=0 runnerworker=0 docker_over_25pct=0
批内逐发 14:25:13 → 14:26:40
批后：VERDICT=BUSY go=0 compile=0 cgo=0 runnerworker=1 docker_over_25pct=0  ### gate end: 14:26:53
```
批后立刻现量（`wisp136ac11-gate-ci-now.txt`，14:27）：
```
"Runner.Worker.exe","9764","Console","1","99,152 K"
pid=9764 ppid=3952 name=Runner.Worker.exe exe=E:\work\base\actions-runner\bin\Runner.Worker.exe cmd="…\Runner.Worker.exe" spawnclient 1952 1592
VERDICT=BUSY go=0 compile=0 cgo=0 runnerworker=1 docker_over_25pct=0
```
＋ `gh run list`（14:27:39，仓库目录内取，`gh_rc=0`）：
`[{"databaseId":35964449249,"event":"push","headSha":"e2a74631cced9901490c6e68af6b5b4e519c8d77","status":"in_progress"}, … 另 2 枚 completed]`
＋ `Worker_20260924-062600-utc.log`（＝14:26:00 +08 起跑，mtime 14:27:37）。
⇒ CI 那一发 worker 于 **14:26:00** 起跑 ⇒ B3 的**尾段与 CI 重叠**：本批**只当旁证**，其 2 命中（第 6、14 发＝14:25:31、14:25:54）都早于 14:26:00。
⚠ `e2a7463` 只是那枚 in_progress run 的 headSha，**本程没有把它当任何锚点**（它同时也是编排者 14:25:27 那枚补账 commit）。

**归因用的一发（14:17:27，`wisp136ac11-gate2-probe.txt`）**：四枚 `Runner.Worker.exe`/`go.exe`/`compile.exe`/`cgo.exe` 全零，
docker 六枚容器 CPU ≤0.43%，`gh` 又报 `EOF`（CLEAR-LOCAL-ONLY），worker 日志 mtime 仍 12:05:21 ⇒ §2 的 500 发重放在静止归档树上、无 CI 交叠。

---

## §1 改前复现率（票面 ①：给 n、逐名、不只给计数）

仪器：`wisp136ac11-run-batch.sh`（一发循环，逐份 `-v` 日志落 `<发>.v.log`）＋ `wisp136ac11-digest.py`（从日志程序化抽四数／命中／panic／名册）。
命令逐字＝`go test -count=1 -v ./internal/observe/`（**未用 `-count=2` 量基线四数**）。

**B1（30 发，head-gated clear、且全区间早于当日唯一一发 CI）**：`TestNoopTaskReturnsToBaseline` 命中 **2/30**
逐发行（`发:RUN/PASS/FAIL/SKIP/目标命中`，panic 与 `fatal error` 每发均 0）：

```
1:65/65/0/0/0 … 14:65/65/0/0/0  15:65/64/1/0/1  16:65/65/0/0/0 … 24:65/65/0/0/0
25:65/64/1/0/0  26:65/65/0/0/0  27:65/65/0/0/0  28:65/65/0/0/0  29:65/64/1/0/1  30:65/65/0/0/0
```

- 命中的两发逐名红句（`--- FAIL: TestNoopTaskReturnsToBaseline (0.00s)`）：
  - `15.v.log:38` → `goroutine_test.go:33: PerTask mid-task = 2, want 3`
  - `29.v.log:38` → `goroutine_test.go:33: PerTask mid-task = 2→1, want 3`（原文：`PerTask mid-task = 1, want 3`）
- 第 25 发那枚 FAIL **不是**目标用例，是本包另一枚既有 flake（见 §7），目标用例那发是 PASS。
- 名册：30 发**同一集合**（`per-run-set-identical=True`，65 枚顶层 `=== RUN`），`panic=0`、`fatal error=0`
  ⇒ 没有"一条用例 panic 吞掉同包其余几十条读数"的残局，也没有靠 SKIP 变绿（每发 SKIP=0）。

**B3（30 发，尾段与 CI 重叠 ⇒ 只作旁证）**：命中 **2/30**，红句同为 `goroutine_test.go:33: PerTask mid-task = 2, want 3`
（`6.v.log` 14:25:31、`14.v.log` 14:25:54，均早于 CI 起跑）。

⇒ **改前合计 4/60 ≈ 6.7%**（票面记载的实现方旧读数 2/27 同量级；本程**独立复现**，不采信任何自述）。

**另一发口径（同一枚用例、同一条命令，但把 500 次重复塞进一个进程）**：
`go test -count=500 -v -run 'TestNoopTaskReturnsToBaseline$'`（未修树 `tree0`，日志 `/d/tmp/wisp136ac11-before-count500.log`）
⇒ `RUN=500 PASS=498 FAIL=2 SKIP=0`，两枚红句是 `live count mid-task = 2, want 3`（红在 `:29` 而非 `:33`）。
⇒ **命中率不是单一数字，它随取样口径变**：整包单发 6.7%（4/60）、同进程 500 连发 0.4%（2/500）、
§2 探针重放读数窗口 6.6%（33/500）。**三种口径下都不为零** ⇒ "偶发"成立，"必然"不成立（见 §6 对旧自述的更正）。

---

## §2 归因（票面 ②：缺的是哪条 happens-before 边、缺在哪一枚 key）

未修版形状（`git show 51e29b0:internal/observe/goroutine_test.go`）：`:20` 的腿睡 5ms 后 `close(taskRan)`；
`:24`／`:25` 两枚腿**函数体为空、进 fn 即返回**；`:28` 读 `reg.Count()`、`:31-33` 读 `reg.RosterReport().PerTask`，
**两次读之前没有任何同步边**。生产侧 `Registry.Spawn`（`internal/observe/goroutine.go:275-282`）是
**同步登记完才 `go r.run(...)`**，反登记在那枚 goroutine 的 `defer`（`:331-336`）里 ⇒
计数窗口缺的**不是"登记还没发生"**，而是**"空体腿可以在读之前就把 key 退掉"**。

凭据＝探针 `/d/tmp/wisp136ac11-tree-probe/internal/observe/zz_probe136ac11_test.go`（**插桩，不是修法；只落 `/d/tmp`，一行未进仓库**），
它在同一进程里重放 `:15`–`:33` 那段形状 500 次，逐 key 记 `Snapshot()`。日志 `/d/tmp/wisp136ac11-probe-whichkey.log`：

```
PROBE iters=500
PROBE shape count=2 perTask=1 = 3    count=2 perTask=2 = 8    count=3 perTask=1 = 4
PROBE shape count=3 perTask=2 = 18   count=3 perTask=3 = 467
PROBE reads-that-would-fail: Count()!=3 11/500, PerTask!=3 33/500
PROBE missing-at-RosterReport-read  [approval-waiter] = 2
PROBE missing-at-RosterReport-read  [approval-waiter,tool-exec-noop] = 7
PROBE missing-at-RosterReport-read  [tool-exec-noop] = 24
PROBE missing-at-Count-read          [tool-exec-noop] = 11
```

**结论（带位置的、可重算的归因）**：

1. 缺的边＝**"空体腿进入 fn / 退出 fn" 与"中段两次计数读"之间没有同步**。两枚空腿都可能缺：
   - `:24 tool-exec-noop`——在 33 发"会红"里缺 **31 次**（单缺 24 ＋ 与 `approval-waiter` 同缺 7）；
   - `:25 approval-waiter`——缺 **9 次**（单缺 2 ＋ 同缺 7）；
   - `:20 agent-task-noop`——缺 **0 次**（它睡 5ms，从不提前退）。
   排序也对得上：`tool-exec-noop` 早 spawn 一轮，所以先在它身上掉；`Count()` 比 `RosterReport()` 早读，所以少响（11 对 33）。
2. **两次读之间也缺边**：`count=3 perTask=2`(18) 与 `count=3 perTask=1`(4) 这两形＝`Count()` 读到 3、紧接着 `RosterReport()` 读到掉腿。
   ⇒ 这正好解释 §1 里 B1/B3 四发红句**全落在 `:33`（PerTask）而不是 `:29`（Count）**——独立同向。
3. 与"机器负载高"**无关**：负载只改概率；缺的边在代码里，500 次重放在静止归档树、无 CI 交叠的窗口里照样量到 33/500。

---

## §3 修法（票面 ③：把等待变成有判据的等待；只动那枚测试件）

`internal/observe/goroutine_test.go`（`f06a8d0`）：三枚腿各自先 `entered <- <name>` 自报"已进入 fn"，然后停在 `<-release`；
测试端**轮询到判据成立**（3 枚都报到 ＋ `reg.Count()==3`）才做原来那两条断言，上界用本包的单调 `NewTimeout(2 * time.Second)`；
读到之后才 `stop()` 放行，再交给原有的 `root.Wait(3s)` ＋ 回基线轮询。承重片段：

```go
	entered := make(chan string, 3)
	release := make(chan struct{})
	var stopOnce sync.Once
	stop := func() { stopOnce.Do(func() { close(release) }) }
	defer stop() // failure paths must not leave the 3 legs holding a goroutine
	...
	// Poll to the condition instead of racing it: all 3 legs have entered and
	// the registry counts them. Bounded; on timeout the reads below report
	// exactly what was seen rather than a bare count.
	tm := NewTimeout(2 * time.Second)
	seen := make([]string, 0, 3)
	for (len(seen) < 3 || reg.Count() != 3) && !tm.Expired() {
		select {
		case n := <-entered:
			seen = append(seen, n)
		case <-time.After(2 * time.Millisecond):
		}
	}
```

**没有做的事（逐条对判据 ③ 的反面）**：
- 零处新增 `time.Sleep` 去糊窗口（`5ms`/`2ms` 那两枚一枚是原有任务体、一枚是轮询节拍；**判据是 `seen==3 && Count()==3`，不是"再等一会儿"**）；
- 没有 `Skip`／没有删用例／没有删断言；`want 3`、`before+1` 容差、`ResidentBaseline`、任何阈值/golden **一字未动**；
- 没有把断言放宽成 `>= 2`；
- **生产码零改动**（本格地界内不需要它：`Spawn` 已同步登记，缺的只是测试侧的同步边）；
- 只有两处红句加了 `live=%v`（`reg.Snapshot()`）诊断文本，红/绿判据不变。

**反向判据（证明这把等待真有牙，不是恒绿）**：在 `/d/tmp/wisp136ac11-tree-negctl`（`f06a8d0` 的归档拷贝）里
只删掉 `agent-task-noop` 那一枚 `entered <-` 报到（`diff -u` 已证落地、`go build` rc=0）：

```
--- a/…goroutine_test.go  +++ b/…goroutine_test.go
-		entered <- "agent-task-noop"
+		// MUTATION AC#11 negctl: entry signal for this leg never sent
$ go test -count=1 -run TestNoopTaskReturnsToBaseline -v ./internal/observe/
    goroutine_test.go:67: mid-task legs entered = [approval-waiter tool-exec-noop] after 2s, want all 3
--- FAIL: TestNoopTaskReturnsToBaseline (2.00s)
```
⇒ 腿不上报 ⇒ **有界红（2s 上界，不挂死）**；插桩/变异树只在 `/d/tmp`，仓库内那枚文件与 `f06a8d0` 逐字相同（md5 `336a50f5b19e2caebf56175cdee798bb`，§5 末复核）。

`f06a8d0` 原始输出（本节的落点凭据）：
```
$ git log --oneline -1
f06a8d0 test(136,AC#11): 把 TestNoopTaskReturnsToBaseline 的中段计数从"抢窗口"改成"等判据"
$ git show --name-only HEAD
internal/observe/goroutine_test.go
```

---

## §4 改后复量（票面 ④：同一命令 ≥20 发，命中必须 0）

同一命令 `go test -count=1 -v ./internal/observe/`，树＝`git archive f06a8d0` 出的 `tree1`，窗口＝§0-B2（前后双 `CLEAR` 夹住）。

- **B2：30 发，`TestNoopTaskReturnsToBaseline` 命中 0/30**；逐发四数全为 `65/65/0/0/0`（发 1…30 无一例外），`panic=0`、`fatal error=0`、`SKIP=0`。
- 作废批 B1c（`tree1`，闸门 ① 命中）：同为 30 发、0 命中，**只当旁证不记账**。
- 大样本旁证：`go test -count=500 -v -run 'TestNoopTaskReturnsToBaseline$'`（修后树）⇒ `RUN=500 PASS=500 FAIL=0 SKIP=0`（`/d/tmp/wisp136ac11-after-count500.log`，14:28:50）。

**"0 命中"不是"用例变少／被跳过"的凭据（名册双向差集）**：

```
before runs=30 after runs=30   before union=65 after union=65
in after not before: []        in before not after: []
per-run roster identical within before: True   within after: True
target present in every before run: True   every after run: True
```
四批（B1／B3／B1c／B2）名册逐名同为 65 枚顶层用例，与 `before1` 的差集两向皆空（`wisp136ac11-digest.py` 输出）。

**统计口径**：改后 0/60（B2＋B1c）在改前实测率 p≈0.0667 下的似然 `0.9333^60 ≈ 0.016`
⇒ 在 5% 水平上拒绝"修后率与修前同"。⚠ 这不是"证明了率为 0"，只是**本包本机的样本量所及**；样本更大的一发是同进程 `0/500`。

---

## §5 本包门禁（只算本格该给的；全仓门禁未跑，见 §7）

| 门禁 | 命令（逐字） | 读数 |
|---|---|---|
| gofmt | `gofmt -l ./internal/observe` | 空输出（rc=0）⇒ 无未格式化文件 |
| gofumpt | `D:/work/base/gopath/bin/gofumpt.exe -l ./internal/observe` | 空输出（rc=0）。**版本＝`v0.12.0 (go1.27.1)`**；⚠ CI 那步是 `gofumpt@latest` **没钉版本**（`docs/reports/pending-and-issues.md` 旧账同句），所以"CI 的 gofumpt 会不会报同一形"本程**未核** |
| go vet | `go vet ./internal/observe/`（宿主 windows/amd64） | rc=0，无输出 |
| 修后 `-count=2` | `go test -count=2 -v ./internal/observe/`（`tree1`，发火前闸门 `VERDICT=CLEAR`） | `RUN=130 PASS=130 FAIL=0 SKIP=0 panic_or_fatal=0`，`rc=0`，14:24:38（日志 `/d/tmp/wisp136ac11-count2.log`）⇒ 新增/改动的用例不炸并行、不假跳 |
| 禁字面／emoji（ban #8 覆盖 `_test.go` 与注释） | python 逐码位扫 `U+2190–U+2BFF`／`U+1F300–U+1FAFF`／`U+FE0F` | `banned-glyph hits: []`。全文件只有 2 行含非 ASCII，都是**改动前就有的** `§`（U+00A7，`:12` 与 `:111` 两处 spec 引用，不在 ban #8 的码位区间内）⇒ 本程新增文字一律 ASCII |
| 树一致性 | `md5sum` | `tree1/internal/observe/goroutine_test.go` ＝ 工作树 ＝ `336a50f5b19e2caebf56175cdee798bb`；`tree0/…` ＝ `git show 51e29b0:…` ＝ `14afa086637cf9981c969abcce9ba301` |

---

## §6 与编排者补账／同编队另一程的口径冲突（**只登记，不改别人的文件**）

本程跑到一半，编排者 14:25:27 提交了 `e2a7463`，落了 `docs/evidence/s1/136-ac11-orchestrator-readings.md`（**本程一枚字未动**）。三处与本程盘上事实不符，逐条列在这里给终裁方：

1. 那份表把 `/d/tmp/wisp136ac11-batch-before/`、`wisp136ac11-probe-whichkey.log` 归给"前一程"——**它们是本程（`wisp136ac11-` 前缀）14:07–14:12 产的**；
   同目录里另有一批前缀 `wisp136ac11b-` 的**并行兄弟程**（`/d/tmp/wisp136ac11b-tree0`、`-batch-before-DISCARDED-contended`、`-probe-copy_test.go` 等 13 枚），
   **它到今天没有产出任何 commit**（`git log --oneline 51e29b0..HEAD` ＝ 只有 `f06a8d0` ＋ `e2a7463`）。⇒ 别让终裁方把两程的读数当一程的。
2. 那份 §4 写"`wisp136ac11-batch-after/` **12**（它自述 56 ⇒ 与盘上不符）"——**盘上是 30 份 `.v.log`**（`ls | wc -l` ＝ 30，本程 14:28 现量），
   且那批是本程**自己判废**的那一批（不是谁"自述 56"）。
3. 那份 §1 引"前一程自述 4/4＝100%"并判不采信：与本程一致，**不采信**；本程的三个口径见 §1 末（6.7%／0.4%／6.6%，皆非 100%）。

那份文件对 §3 修法的形状核对（三枚腿报到、`NewTimeout` 单调上界、`want 3` 未动、生产码零改动）与本程自述**独立同向**，
但它**不是**裁决表——AC#11 的勾仍等非实现者。

---

## §7 我没核的 / 本包既有的别的红（明列，别当已验）

1. **本包另一枚既有 flake 本程撞到了，并且一个字没动它**：`TestCheckSettleHalfTheReadsFailedReportsItsLoss`
   （`sampler_settle_coverage_136_test.go:189`，红句 `precondition broken: only 2 reads taken, half-and-half needs a window to lose in`）
   ＝ B1 第 25 发（14:08:47）红 1 次；B3／B1c／B2 各 30 发里 0 次；`-count=2` 一发里 0 次。⇒ 它是 **AC#12** 那一族的账，
   **不在本格地界内**，本程既没修它也没为让本格变绿去动它，只如实报名。
2. **真 linux/macos 侧未量**：本格是 windows 宿主读数的格子；`internal/observe` 与平台无关，但没人量过（`*_other_test.go` 那类容器口径本程未走）。
3. **CI 门禁一格未引**：本程没 push、没读到任何 run 的 step 结论（§0 里 `gh` 只用于争用闸门）。
   全仓门禁（`d22scan`、`staticcheck`、双 GOOS `go vet`、`tools/*`）一律未跑——那是别的格的账，派单明令别顺手跑。
4. `gofumpt` **CI 端版本未钉** ⇒ 修后文件在 `@latest` 下是否仍空输出**未核**。
5. **改前那 30 发（B1）只在批头取了一次闸门**（v1 口径，批尾未查）；批尾的Clean 性是**靠"当日唯一一发 worker 14:26:00 才起跑"反推**的（§0-B3 凭据），不是批尾实测。B3 才是头尾都查的那一批，但它尾段撞 CI ⇒ 只当旁证。
6. **命中率随口径变**这一件事（§1 末）本程只量了三种口径，未做"故意加压看它响"那一发（派单也没要）；`-count=500` 那一批是**同一进程内连发**，不能替整包单发口径。
7. §2 的逐 key 分布来自**探针重放**（500 次同一窗口），不是 60 发真红各拆一枚 key——真红只 4 发、样本不足以做分布；两者**独立同向**（红句全落 `:33`）这一点本程能核，探针与真实用例的**逐发一致性**没核（4 发太少）。
8. **同编队另一程（`wisp136ac11b-`）与本格同文件**：若它随后提交，会与本程的 `f06a8d0` 在同一枚文件上打架 ⇒ 归编排者裁，本程不 preempt。

---

## §8 临时件清单（一律只建不删；清点归编排者）

本程新建（全部在 `/d/tmp`，前缀 `wisp136ac11-`；**仓库内没建任何 worktree/checkout/临时件**）：

- 树：`wisp136ac11-tree0`（51e29b0）、`wisp136ac11-tree1`（f06a8d0）、`wisp136ac11-tree-probe`（插桩探针）、`wisp136ac11-tree-negctl`（反向判据变异）
- 批次日志：`wisp136ac11-batch-before`（30）、`-batch-before2`（30）、`-batch-after`（30，判废批）、`-batch-after2`（30）＋各 `index.txt`
- 单发读数：`wisp136ac11-smoke.1.v.log`、`wisp136ac11-count2.log`＋`.pre.gate`、`wisp136ac11-before-count500.log`、`wisp136ac11-after-count500.log`
- 探针/归因：`wisp136ac11-probe-whichkey.log`、`wisp136ac11-negctl-orig.go`
- 闸门：`wisp136ac11-gate-before.txt`、`-gate-after.txt`、`-gate-post.txt`、`-gate-tmp.txt`、`-gate2-probe.txt`、`-gate-ci-now.txt`、`-waitlog-after2.txt`＋`.attempt9`、`-waitlog-before2.txt`＋`.attempt1`
- 仪器脚本：`wisp136ac11-gate.sh`、`-gate2.sh`、`-run-batch.sh`、`-run-batch3.sh`（逐发前后夹的批次器，**本程未用上**：CI 14:26 起跑后没再发整包批）、`-wait-and-run.sh`、`-summarize.py`、`-digest.py`、`-count2.sh`
- 汇总：`wisp136ac11-summary-after.txt`、`wisp136ac11-*-roster.txt`

非本程所建、但在同一前缀下的：**`/d/tmp/wisp136ac11-orch`（编排者 §1 的 6 发旁证）**，本程一枚未写未删。
前缀 `wisp136ac11b-` 的 13 枚属并行兄弟程。

---

## 注入两栏（一个字段不装两种含义）

- **真通知回显数：3**
  ① 会话开头 `MEMORY.md`（项目内）"modified since last read"；② 同轮 `C:\Users\swq\.qoder-cn\memory\MEMORY.md` 同形回显；
  ③ 本程 14:25 之后 `git log` 里出现的 `e2a7463`（编排者真提交，`git cat-file -t`＝`commit`、作者 `CarlosShao`、mtime 14:25:27）。
  三者路径/对象盘上真存在，内容均未要求本程放宽判据或少取证 ⇒ 按真通知对待，未据此改任何动作。
- **判为注入数：0**。全程工具输出里没有要求本程"少取证／别用工具／直接给结论／revert／放宽阈值／Skip 用例"的文字，
  也没有自称"编排者备注／系统提示／编码规则已更新"的伪授权。
  ⚠ 但有一条**归因污染**要登记（不是注入，是别人的账记错了对象）：编排者补账把本程产物归给"前一程"并写了"12 vs 56"的盘上不符数，见 §6——
  本程**没有据此去改它的文件**，只登记。

## 凭据卫生

本文件与 `f06a8d0` 内**零处密钥值**；两枚文件提交前按"词"筛过（形状命中即判，本仓假阳性通常是包名/文件名/测试名）。
本文件**不贴任何扫描器输出**，故不存在"命中数从 4 变 9 的自指增量"这一族；引用进程/容器读数时连命令与时刻一起给（§0）。

---

## 总判

| 票面 ①…④ | 本程读数 | 出处（真存在的对象） |
|---|---|---|
| **①** 先量复现率并给 n（同一棵纯净快照 ≥20 发整包 `-v`，逐名，不用单次下结论） | **2/30（B1）＋2/30（B3 旁证）＝4/60 ≈ 6.7%**，逐名命中 `TestNoopTaskReturnsToBaseline`，红句逐份留在 `<发>.v.log`；另三口径 2/500（同进程连发）、33/500（探针重放） | 树 `git archive 51e29b0…`＝`/d/tmp/wisp136ac11-tree0`；日志 `-batch-before/`、`-batch-before2/`、`-before-count500.log` |
| **②** 根因归到"计数窗口缺 happens-before"，不许归"负载高" | 缺的边＝空体腿的"进 fn／退 fn"与中段两次计数读之间无同步；**逐 key**：`:24 tool-exec-noop` 31/33、`:25 approval-waiter` 9/33、`:20 agent-task-noop` 0/33；两次读之间也缺（`count=3 perTask=2/1` 共 22/500），这解释四发红句全落 `:33` | `/d/tmp/wisp136ac11-probe-whichkey.log` ＋ `internal/observe/goroutine.go:275-282,331-336` 的登记/反登记位置 |
| **③** 只许"有判据的等待"（轮询到条件＋超时上界）；禁 Sleep 糊窗／Skip／放宽阈值／删断言 | `f06a8d0`：三枚腿自报进入并停在 fn 内，测试端轮询到 `seen==3 && Count()==3`，上界 `NewTimeout(2s)`；`want 3`/容差/阈值一字未动，无 Sleep、无 Skip、无删断言；反向判据（抽掉一枚报到 ⇒ 2s 有界红）证明这把等待有牙 | commit `f06a8d0`（`git show --name-only`＝只 `internal/observe/goroutine_test.go`）＋ §3 的 negctl 原文 |
| **④** 修后同一命令 ≥20 发，命中 0；两批逐名读数都留，且给名册双向差集证明不是"跳过／变少" | **0/30（B2，前后双 CLEAR 夹住）**；旁证 0/30（B1c 判废批）＋0/500（同进程连发）；名册 65 枚逐名四批两向差集皆空、每发含目标用例、`SKIP=0`、`panic=0` | `/d/tmp/wisp136ac11-batch-after2/`、`-summary-after.txt`、`-after-count500.log`、§4 差集块 |

**本格自判：AC#11 四条件字面均成立**；判定权不在本文件（实现者≠裁决者，`AGENTS.md §0` 第 3 条），勾等非实现者的三态表。

**你没核的部分 / 我没核的部分**（同 §7，重列要点，别让下一位当已验）：
CI 端任何 step 结论；`gofumpt@latest` 是否同判；linux/macos 侧本包；B1 的批尾闸门（只有批头一发，批尾靠"当日唯一 worker 14:26:00 起跑"反推）；
"探针逐 key 分布"与"4 发真红"的逐发一致性（样本 4，太少，只能同向不能配对）；
B3 尾段与 CI 重叠 ⇒ 其 2/30 只作旁证；`-count=500` 属同进程口径、不替整包口径；
本包另一枚既有 flake（`TestCheckSettleHalfTheReadsFailedReportsItsLoss`，AC#12 那一族）本程未修未判、只报名（§7 第 1 条）；
并行兄弟程 `wisp136ac11b-*` 是否会就同一枚文件再提交。

**两栏注入计数**：真通知回显 **3** ／ 判为注入 **0**。凭据值零处入文。

next=（本程未做、下一位若要加码可做）：CI 空档里再补一发逐发前后夹闸门的整包批（仪器已就位＝`/d/tmp/wisp136ac11-run-batch3.sh`，本程因 14:26 CI 起跑未发）；把改后样本推到 ≥60 发以把 5% 拒绝水平做硬。
