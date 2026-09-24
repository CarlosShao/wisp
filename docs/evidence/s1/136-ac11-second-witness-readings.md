# 票 136 AC#11 —— 独立第二 witness 读数表（被派程 `wisp136ac11b-*`）

**这张表是什么**：一枚被派来做 AC#11 这一格的程，**自己现量**的改前／改后基线、归因探针与本包门禁。
**这张表不是什么**：它**不是**本格的主表。本格的主表与修法**在我开工之后由另一枚在飞的程落地**：
`c2a3bd6`（改前 4/60 ＋改后 0/30 的证据表）与 `f06a8d0`（那枚测试件的修法）。
`docs/evidence/s1/136-ac11-impl.md`（14:34 出现，28 KB）**不是我写的**，我一个字没改；
`docs/evidence/s1/136-ac11-r1-acceptance.md`（终裁表）**也不是我写的**。
⇒ 所以这张表按"第二 witness"记账，**不抵主表、也不与主表争勾**。

**本格我改了几枚仓内文件**：**代码 0 枚**（见 §7），**证据 1 枚＝本文件**。
临时件前缀一律 `wisp136ac11b-`（前一程的 `wisp136ac11-*` 只读复用、未覆盖、未删）。

---

## 0. 锚点（全部自量，不抄任何人的转述）

| 用途 | sha（`git rev-parse`/`git show` 现打） | 我量它的时刻 |
|---|---|---|
| 开工时 HEAD ＝**改前**被验版本 | `51e29b044bf928e300c4a48ca1cc0932e02300b4` | 14:11:xx（本程第一条命令 `git rev-parse HEAD`） |
| **改后**被验版本（他人在飞的修法） | `f06a8d0f3ad1ba9e39b3b77e213ec7c6c0d742bf` | `git show --name-only` 打出全 sha＋`internal/observe/goroutine_test.go` 一枚 |

被验版本的盘上身份（**工作树不当被验版本**）：

- 改前树：`/d/tmp/wisp136ac11b-tree0` ＝ `git archive 51e29b0 | tar -x`，`.go` 文件 **469** 枚 ＝
  `git ls-tree -r 51e29b0 | grep -c '\.go$'` 的 **469**（取件完整，不是只抽了半个包）。
- 改后树：`/d/tmp/wisp136ac11b-tree-after` ＝ `git archive f06a8d0 | tar -x`。
- 两棵树之间**只有一枚文件不同**：`git diff --name-only 51e29b0 f06a8d0` →
  `internal/observe/goroutine_test.go`（计数 1）。⇒ 我这两批改前／改后的**唯一自变量就是那枚测试件**。
- `internal/observe/goroutine_test.go` 的 md5：改前 `14afa086637cf9981c969abcce9ba301`、
  改后 `336a50f5b19e2caebf56175cdee798bb`；后者与 14:37:50 那次**工作树**同 md5（＝我跑的分母就是盘上现在那版）。
- **仓内没有建任何 worktree／没有 checkout**；只 `git archive` 到 `/d/tmp`。

与前一程台件的关系（**复用前先核字节**，不是我口头假设）：
`diff -rq /d/tmp/wisp136ac11-tree0/internal/observe /d/tmp/wisp136ac11b-tree0/internal/observe` → **零差异**，
`goroutine.go` 与 `goroutine_test.go` 两枚 md5 逐一相同 ⇒ 前一程那 30 份日志与我的改前批**跑在同一版代码上**，
所以 §3 的一致性核对成立。

---

## 1. 争用闸门（四查逐字；两批都在放行窗口内）

工具：`/d/tmp/wisp136ac11b-gate.sh`（前一程 `wisp136ac11-gate.sh` 的四查形状 ＋ 我加的第五条"外来文件活动"）。

### 1.1 闸门第一次全量记录（14:17:51，取件后、正式批之前）

```
### gate taken at: 2026-09-24 14:17:51 +0800
--- 1) host processes Runner.Worker.exe / go.exe / compile.exe / cgo.exe ---
[Runner.Worker.exe]
[go.exe]
[compile.exe]
[cgo.exe]
--- 1b) any process whose ExecutablePath is under actions-runner\_work ---
Name                ProcessId ExecutablePath
----                --------- --------------
Runner.Listener.exe      3952 E:\work\base\actions-runner\bin\Runner.Listener.exe
--- 2) docker ps (container load does not show in the host process list) ---
union-proxy image=union-api-proxy:local status=Up 8 days
clipsync image=clipsync-clipsync status=Up 4 days
clipsync-minio image=minio/minio:latest status=Up 12 days (healthy)
clipsync-admin-int image=4cbfeb6e08ee status=Up 2 weeks
clipsync-db image=postgres:15-alpine status=Up 13 days (healthy)
clipsync-redis image=redis:7-alpine status=Up 13 days (healthy)
docker_rc=0
--- 3) gh run list --limit 5 --json databaseId,status,headSha ---
failed to get runs: Get "https://api.github.com/repos/CarlosShao/wisp/actions/runs?per_page=5&exclude_pull_requests=true": EOF
gh_rc=1
--- 4) runner _diag worker log mtimes ... ---
Worker_20260924-040204-utc.log mtime=2026-09-24 12:05:21 +0800   ← 当时最新一枚，距 14:17 约 2h12m
```

- ② docker：6 枚容器在跑但 `docker stats --no-stream` 实测 CPU `0.00% / 0.00% / 0.15% / 0.00% / 0.00% / 5.60%`，
  全部常驻（Up 4–14 天）⇒ 记为**常驻背景**，不是测试负载；两批同环境。
- ③ `gh` 那次 **EOF／rc=1** ⇒ 按规矩**降级 `CLEAR-LOCAL-ONLY`**：那一发**不是**"确认为空"。
  14:38:28 我从仓目录重取一次成功：最近 5 枚 run `status` 全 `completed`（最新 headSha `e2a7463…`）、**零枚 `in_progress`** ⇒ 记为 `CLEAR-CI-READ-OK`（只此一发，见 §8 未核清单）。
- ④ runner：14:38 现量最新 Worker 日志 `Worker_20260924-062600-utc.log` mtime **14:31:53**、
  `_work/wisp` 目录 mtime **14:29:34** ⇒ **14:26–14:31:53 本机 runner 上确实跑过一发 CI**（那是 `e2a7463` 那 push 触发的），
  它在我第一扇放行窗（14:33:09）**之前 ~76 s 结束**；我的两批窗口内 `Runner.Worker.exe` 计数为 **0**。

### 1.2 有一批被争用打掉并重跑（如实报，不挑好看的用）

| 批 | 窗口（每发自己的 `date`） | 闸门判定 | 处置 |
|---|---|---|---|
| `wisp136ac11b-batch-before-DISCARDED-contended` | 14:18:45–14:20:18 | **作废**：窗口内他人在飞的 `ac11-orch-after.log.1..44` 在落盘（`find -newermt` 命中 15 枚外来文件） | 目录改名留档，读数只当"含争用"参考（0/30），**不作基线** |
| `wisp136ac11b-batch-before`（正式改前 30） | 14:33:12–14:34:32 | **放行**：开窗前连续 4 发轮询 `foreign10s=0 go.exe=0`（14:32:47／58／14:33:09），批后 `find -newermt '-130s'` 排除我前缀后**零命中** | 采用 |
| `wisp136ac11b-batch-after`（正式改后 30） | 14:34:37–14:36:00 | **放行**：紧跟改前批、无人插队；批后外来活动只有 `wisp136ac11v-assert-{before,after}.txt` 两枚 grep 落盘（非 CPU 负载） | 采用 |

开窗轮询原样（节选，每行一次 `sleep 10` 后现量）：

```
w21 at=14:32:13 orch_before2_logs=30 foreign10s=0 go.exe=0
w22 at=14:32:24 orch_before2_logs=30 foreign10s=3 go.exe=0
w23 at=14:32:36 orch_before2_logs=30 foreign10s=0 go.exe=0
w24 at=14:32:47 orch_before2_logs=30 foreign10s=0 go.exe=0
w25 at=14:32:58 orch_before2_logs=30 foreign10s=0 go.exe=0
w26 at=14:33:09 orch_before2_logs=30 foreign10s=0 go.exe=0
WINDOW OPEN at 14:33:09
```

**没有杀过任何进程、没有改过任何工作流**；`date` 只做"每发记一行"，**不用于相减得时长**（本机 wall clock 会跳）。

---

## 2. 改前基线：同一棵纯净快照连跑 30 发整包 `-v`

命令（逐字）：`go test -count=1 -v ./internal/observe/`，cwd `/d/tmp/wisp136ac11b-tree0`，
脚本 `/d/tmp/wisp136ac11b-run-batch.sh`，日志 `/d/tmp/wisp136ac11b-batch-before/<n>.v.log`（n=1..30）。

**目标用例 `TestNoopTaskReturnsToBaseline` 命中数＝2/30；n=30；复现率 6.7%。**

| run | RUN | PASS | FAIL | SKIP | 目标红 | panic | fatal | 红句原文 |
|---|---|---|---|---|---|---|---|---|
| 1–26, 28, 30 | 65 | 65 | 0 | 0 | 0 | 0 | 0 | — |
| **27** | 65 | 64 | 1 | 0 | **1** | 0 | 0 | `goroutine_test.go:33: PerTask mid-task = 1, want 3` |
| **29** | 65 | 64 | 1 | 0 | **1** | 0 | 0 | `goroutine_test.go:29: live count mid-task = 2, want 3` |

四数**只能从 `-v` 量**（非 `-v` 只有一行 `ok`）；每行的 RUN/PASS/FAIL/SKIP 由
`/d/tmp/wisp136ac11b-summarize.py` 现抽，全 30 行原样在 `/d/tmp/wisp136ac11b-b1-roster.txt` 与各 `.v.log`。

聚合读数：`TOTAL runs=30 hits=2 (rate 2/30)`、`TOTAL panic=0 fatal=0`、
`four-number spread: RUN=[65] PASS=[64, 65] FAIL=[0, 1] SKIP=[0]`、
`roster: 1 distinct top-level run-name sets over 30 runs; modal size=65 occurs=30`（**无缺名、无跳名**）。

**名册差集（防"一条 panic 吞掉同包其余几十条"那一族）**：30 发的 `=== RUN` 顶层名册两两相同
（`1 distinct set over 30 runs`），65 枚，且 §5 的改后批与它**逐名相同** ⇒ 2/30 不是"用例少了"或"被跳过"造出来的，
`SKIP=0` 全程成立。

**两枚红句落在不同的那一跳**：`:29`（`reg.Count()` 那一跳）与 `:33`（`rep.PerTask` 那一跳）各一发
——这一发是本格**第一次**在整包批里抓到 `:29` 红（前一程那 30 发的两枚红都在 `:33`），
它与 §4 探针给出的"两跳各自独立的命中率"预测**方向一致**（`:29` 更低）。

---

## 3. 与前一程那 30 份的一致性核对（＝第二 witness 成立与否）

前一程台件：`/d/tmp/wisp136ac11-batch-before/*.v.log`（30 份，`index.txt` 14:07:48–14:09:37），
**盘上现重扫**（不是我抄谁的数）：

| 批 | 谁 | n | 目标红 | 红在哪一跳 | rc=1 发数 |
|---|---|---|---|---|---|
| `wisp136ac11-batch-before` | 前一程（r1） | 30 | **2**（`15.v.log`／`29.v.log`） | `:33` 两发（`PerTask=2`、`PerTask=1`） | 3 |
| `wisp136ac11b-batch-before` | **本程** | 30 | **2**（`27.v.log`／`29.v.log`） | `:33` 一发 ＋ `:29` 一发 | 2 |
| `ac11-orch-base` | 在飞的那枚程（他人） | 60 | **2**（`log.20`／`log.43`） | `:33` 两发（均 `PerTask=2`） | 2 |
| `wisp136ac11-batch-before2` | 他人（14:26 落盘） | 30 | **2** | 未逐名读（只计数） | — |

⇒ **一致**：三枚独立改前批各 2 枚红（2/30、2/30、2/60），与票面那条"约 27 发命中 2 发"同向；
聚合 **8/150 ≈ 5.3%**。**不是 100%、不是必然红**——前一程交件里那句"4/4 红＝100%、这不是 flake"
**与它自己留下的台件矛盾、也不成立**（它自己的 30 份里就是 2/30）。

**不一致的地方我也照实记**：r1 那批有 **3** 枚 rc=1 但只有 **2** 枚目标红，多发的那枚是
`25.v.log` 的 `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（**另一格、另一族**，见 §8）；
我这批 rc=1 发数（2）与目标红发数（2）相等。红句跳数不同（我有 `:29`、r1 没有）＝**两跳概率本就不同**，不是矛盾。

---

## 4. 归因：缺的那条 happens-before 边 ＋ 少的那枚 key（探针＝仪器，不是修法）

**边的位置**（读码，`internal/observe/goroutine.go`）：`Spawn` 在 `:276` **同步**做
`r.live[gkey{name, owner}]++`，然后 `:281` 才 `go r.run(...)`；而注销在 `run` 的 `defer` 里
（`:331-336` 把计数减回去并 `delete`）。测试 `internal/observe/goroutine_test.go`（改前版）
`before := runtime.NumGoroutine()` 在 `:16`，三枚 `reg.Spawn` 在 `:20`／`:24`／`:25`，
`:28` 读 `reg.Count()`、`:31` 读 `reg.RosterReport()`。
⇒ **`:28`/`:31` 这两次读与"腿还在里面"之间没有任何 happens-before**：`Spawn` 只保证"键已经登记"，
不保证"那个 goroutine 还在函数体里"。`:24`／`:25` 两枚腿的**函数体是空的**（`func(ctx context.Context) {}`），
进去就返回 ⇒ 它们的 `defer` 注销可以落在 `:28` 之前、也可以落在 `:28` 与 `:31` 之间。
（前一程交件里"`before` 在 spawn 之后所以基线被抬高"那句**不成立**：`:16` 早于 `:20`。）

**探针**：`/d/tmp/wisp136ac11b-probe-copy_test.go` → 复制进**我自己的** `/d/tmp` 树
（`wisp136ac11b-tree0/internal/observe/zz_probe136ac11b_test.go`，跑完仍在 `/d/tmp`，**没进仓库、没进任何 commit**），
把改前那两次读原样重放 **500 次**，逐次记下"读的那一刻三枚 key 谁还活着"。
**它是探针不是修法**——它不判绿、不替代用例，只回答"少哪枚 key"。
输出：`/d/tmp/wisp136ac11b-probe-whichkey.log`（14:36:41–14:36:46，安静窗口内）：

```
PROBE iters=500
PROBE shape count=1 perTask=1        = 1
PROBE shape count=2 perTask=1        = 2
PROBE shape count=2 perTask=2        = 2
PROBE shape count=3 perTask=1        = 6
PROBE shape count=3 perTask=2        = 15
PROBE shape count=3 perTask=3        = 474
PROBE reads-that-would-fail: Count()!=3 5/500, PerTask!=3 26/500
PROBE missing-at-RosterReport-read  [approval-waiter] = 4
PROBE missing-at-RosterReport-read  [approval-waiter,tool-exec-noop] = 9
PROBE missing-at-RosterReport-read  [tool-exec-noop] = 13
PROBE missing-at-Count-read         [approval-waiter,tool-exec-noop] = 1
PROBE missing-at-Count-read         [tool-exec-noop] = 4
```

**结论（落在缺的边上，不是"负载高"）**：

- 少的是 **`tool-exec-noop`（`:24`）与 `approval-waiter`（`:25`）**：在 `RosterReport` 那一跳缺 `tool-exec-noop`
  22/500、缺 `approval-waiter` 13/500（9 次两枚同时缺）；在 `Count` 那一跳缺 `tool-exec-noop` 4/500、两枚同缺 1/500。
- **`agent-task-noop`（`:20`）一次都没少过（0/500）** ⇒ 它的体里那句 `time.Sleep(5 * time.Millisecond)`
  正好就是"唯一存在着的窗口"，这从反面指认了机制：**红不红取决于腿还活不活，不取决于机器快慢**。
- `tool-exec-noop` 比 `approval-waiter` 更容易先退（22 vs 13）＝它先 `Spawn`、被调度得更早 ⇒ 与
  "空体腿立刻跑完 `defer`" 这一形完全自洽。
- 探针预测的两跳比率 `PerTask!=3` 26/500＝5.2%、`Count()!=3` 5/500＝1.0%，与我 §2 整包实测
  （30 发里 `:33` 一发、`:29` 一发）以及前一程 2/30、他人 2/60 **同量级**；
  前一程同形探针的独立一发：`PerTask!=3` 33/500、`Count()!=3` 11/500、缺 key 序 `tool-exec`＞`approval`＞`agent=0`
  （`/d/tmp/wisp136ac11-probe-whichkey.log`）⇒ **两枚探针同向**。

---

## 5. 复量：改后同一棵归档树连跑 30 发 ＋ 修法的形状核查

命令逐字同 §2，只换树：cwd `/d/tmp/wisp136ac11b-tree-after`（`git archive f06a8d0`），窗口 14:34:37–14:36:00。

```
TOTAL runs=30 hits=0 (rate 0/30)
TOTAL panic=0 fatal=0
four-number spread: RUN=[65] PASS=[65] FAIL=[0] SKIP=[0]
roster: 1 distinct top-level run-name sets over 30 runs; modal size=65 occurs=30
```

**名册双向差集（证明"0 命中"不是"用例变少／被跳过"）**——`comm -3` 现量：

| 比较 | 结果 |
|---|---|
| 改前 30 发的名册（65 枚）vs 改后 30 发的名册（65 枚） | **对称差为空** |
| 前一程 30 发的名册（65）vs 本程改前名册（65） | **对称差为空** |
| `SKIP` 计数 | 两批全 30 发 **0** |
| `TestNoopTaskReturnsToBaseline` 在两批名册里 | **各 30/30 都在**（既没被跳、也没被删） |

### 5.1 那枚已落地的修法（`f06a8d0`）合不合本格的判据 ③——我按形状核过

- 三枚腿各自 **报到（`entered <- name`，缓冲 3）＋ 等放行（`<-release`）**＝channel rendezvous，
  读完 `Count()`／`PerTask` 之后才 `stop()` 放行 ⇒ "mid-task 观察点"从"抢来的窗口"变成**必然存在的窗口**。
- 等待是**有判据的等待**：`for (len(seen) < 3 || reg.Count() != 3) && !tm.Expired()`，
  上界 `NewTimeout(2 * time.Second)`（单调钟，D42#9 那一族），超时**要红**（`t.Fatalf("mid-task legs entered = %v …")`）不是静默。
- 判据 ③ 的三条禁令逐条查：**没加 `time.Sleep` 糊窗口**（`grep -n "time.Sleep"` 改后版只剩 `:47` 腿体里原有那 5 ms、
  `:94` 基线回落轮询原有那 5 ms——两处改前就有，我没动、它也没动）；**没 `Skip`**（SKIP=0/30 两批）；
  **没放宽阈值**（`want 3`／`!= 3` 两处断言原样，红句只是多带了 `reg.Snapshot()`）。
- **只改 `internal/observe/goroutine_test.go` 一枚**：`git diff --name-only 51e29b0 f06a8d0` → 计数 **1**。
  ⇒ 编排者判"不需要解冻生产码"**成立**：`internal/observe/goroutine.go` 一字未动，SPEC-01 §7 那句登记断言
  （"boot＋一枚 no-op task 后回落到常驻基线"）仍然由同一枚用例在证，只是观察点从抢窗口改成等判据。
- ⚠ 归属要说白：**这枚修法不是我写的、不是我提交的**（我只在 `/d/tmp` 的归档树上量它）。见 §7。

---

## 6. 本包门禁（只算本格该给的，没跑全仓）

被验版本＝`f06a8d0` 归档树（＝工作树那枚文件同 md5），时刻 14:37:26–14:37:33：

| 门禁 | 命令（逐字） | 读数 |
|---|---|---|
| gofmt | `gofmt -l ./internal/observe` | **空输出，rc=0** |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l ./internal/observe` | **空输出，rc=0**；版本 **`v0.12.0 (go1.27.1)`**，宿主既有二进制、**未 `go install`** |
| gofumpt 版本口径 | `.github/workflows/ci.yml:113` | CI 那步是 `go install mvdan.cc/gofumpt@latest`——**没钉版本**（本格读数以宿主 v0.12.0 为准并写明） |
| go vet（宿主） | `go vet ./internal/observe/` | **rc=0、零输出** |
| 改后整包双跑 | `go test -count=2 -v ./internal/observe/` | rc=0；四数 **RUN=130 / PASS=130 / FAIL=0 / SKIP=0**（65×2，`/d/tmp/wisp136ac11b-count2-after.log`） |
| 工具链 | `go version` | `go1.27.1 windows/amd64` |

对照参考（同一条 gofumpt 指到我 `/d/tmp` 改前树）：它只点了 `zz_probe136ac11b_test.go` 一枚——
**那是探针、不是交付物**，它**没进仓库**，所以不构成门禁缺口。

---

## 7. 本格改动文件清单（＋一件我必须报的盘上事实）

**我改的仓内文件：0 枚代码 + 1 枚证据（本文件）。**

为什么代码那枚我没动：**我开工时工作树是干净的**（14:11 那发 `git status --porcelain` **零行**），
而 `internal/observe/goroutine_test.go` 在 **14:13:27** 被**另一枚在飞的程写入**（14:13:41 现量
`git status` 出现 ` M`、mtime 14:13:27、与我 14:11 那次 Read 的内容不同），
随后以 `f06a8d0`（14:14:06）落地。共享工作树里**覆写／回退别人的未提交改动**是本仓硬规矩禁的
（`AGENTS.md §1.4` ＋ `.scratch/wisp/issues/README` 规则 1/2），所以我：
**没编辑那枚文件、没 `checkout`、没 `add`、没 commit 它**，只把它当成"改后版本"在 `/d/tmp` 归档树上量。
留档：`/d/tmp/wisp136ac11b-foreign-gitstatus-1.txt`（那次的 status 原样）、
`/d/tmp/wisp136ac11b-foreign-md5-1.txt`（当时那枚文件的 md5）。
`/d/tmp/wisp136ac11b-foreign-worktree-edit.patch` **是 0 字节**——我导出的那一刻它已被 `f06a8d0` 收进 HEAD，
所以 `git diff` 空；**那次 diff 的正文我是在 14:13:4x 的工具输出里读的，没落盘**（这一处只算自述，
可重算的凭据是 `git show f06a8d0` 本身）。

⇒ 派单里"码与证据分两枚 commit"这一条**对我只剩证据那一枚**：码那一枚已经在盘上、且不是我落的。
**本格结案主表也不是我这枚**（`136-ac11-impl.md` 已由在飞的程写到 28 KB），本文件按**第二 witness** 记账。

---

## 8. 顺带发现（不属本格，我没动它）

前一程那批的 `25.v.log` 里有**另一枚偶发红**：

```
sampler_settle_coverage_136_test.go:189: precondition broken: only 2 reads taken, half-and-half needs a window to lose in
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.13s)
```

它是**票 136 AC#12** 那格今天新装的用例（`f53ad5c`／`c03aee3`），红因形状是"100 ms 窗口只取到 2 次读数"
（窗口／间隔计时形状，不是 §4 那条注册表边）。我在盘上可查的 **240 发**（r1 60 ＋ 他人 120 ＋ 本程 60）里
只命中这一发 ⇒ 罕见、**未立案**。⚠ **本格不修、票面不勾、也不替它开 AC**——只登记现象与出处，
归口留给编排者。

---

## 9. 我没核的部分（读这张表的人请连这栏一起信）

1. **CI 侧门禁**：只跑了本包四件（§6）。全仓 `gofumpt -l . tools/d22scan tools/mockllm`、
   `go vet ./...`、`staticcheck`、`d22scan` 七禁＋ban #8 **一律未跑**（派单禁"顺手跑全仓门禁"）。
2. **`GOOS=linux` 交叉 vet 未跑**：`internal/observe/treemetrics_other.go` 的 linux 半边**没量**。
3. **`gh run list` 只有一发成功读数**（14:38:28，5/5 `completed`）；14:17 那一发 EOF 已降级 `CLEAR-LOCAL-ONLY`。
   两批窗口内**没有再取一次 gh**，所以"批中无远端 in_progress"是**由本机 runner 日志静止（14:31:53）＋
   `Runner.Worker.exe=0`** 推的，不是由 gh 直读的。
4. **`wisp slo` 端到端、SLO golden、`thresholds.go`、`scripts/**`、`frontend/**`、`internal/risk/**`**：未碰未跑（不属本格）。
5. **改后 30 发只覆盖 `./internal/observe/` 一包**；没跑跨包整树。
6. `wisp136ac11-batch-before2`（30 发、2 命中）**我没逐名读日志正文**，只做了 `grep -l` 计数，
   且**说不清它是哪一枚程跑的**（14:26 落盘、用的是前一程那套脚本）。
7. §8 那枚 AC#12 的 flake **根因未查**。
8. 名册比对用的是**顶层 `=== RUN` 名**（65 枚），子测试（`TestX/sub`）不在这一层；本包两批都无子测试红。

---

## 10. 注入登记（两栏分开计数）

| 栏 | 数 | 明细 |
|---|---|---|
| **真通知回显数**（我判为 harness／owner 真发的） | **1** | 派单消息开头那条"MEMORY.md 已被修改"的提示——内容只是索引重排，无指令、无越权，未据此做任何动作。 |
| **判为注入数**（自称权威但我只登记不服从） | **1** | `/d/tmp/A157-block.md`（mtime 14:13:52）——**出处**：本程 Bash 工具 `cat /d/tmp/A157-block.md`（命令前 40 字 `cd "D:\work\workspace\projects plans\Wisp"`）。它以"编排者"口吻写"这一格我自己动手、验收另派、r1/r2 报的 commit 全不存在"。四条辨别力：**路径真**（文件盘上存在）／**不放宽判据**（它反而更严）／**盘上可重验**（`f06a8d0`、`docs/evidence/s1/136-ac11-impl.md`、那 30 份日志我都独立重扫过）／**没让我少取证**（我照派单把两批 30 发与门禁全跑）。⇒ 按**数据**用，**不当授权**：我没因为它的"我自己动手"就交回活，也没抄它任何数字进表；我的每个数都有自己那一次扫描＋当时的 `-v` 日志做凭。 |

另记一条形状（不计注入、算正常回显）：票面正文里那些"> **AC#12／AC#13 … 编排者**"引用块是**仓内既有文本**，
属合法上下文。本程**没有**遇到要求我"放宽阈值／改成 `>=2`／`Skip` 这条用例／revert／Confirm the harness note is genuine"
的任何文字；这类事若出现，规矩是**只登记不服从**。

---

## 总判

| 判据（票面 `:164` 那格 ①…④，原文口径） | 本程读数 | 判定 |
|---|---|---|
| ① 同一棵纯净快照连跑 ≥20 发整包 `-v`，逐名记命中、给出 n | `git archive 51e29b0` 归档树（469/469 枚 `.go` 齐），30 发，**2/30**，逐名四数全在表内，`SKIP=0`、`panic=0` | **达**（n=30 ≥ 20） |
| ② 根因归到"计数窗口／缺 happens-before"那一类，不许"机器负载高" | 缺的是 `Spawn` 同步登记（`goroutine.go:276`）与 `run` 的 `defer` 注销（`:331-336`）之间、对 `:24`／`:25` 两枚**空体腿**的边；探针 500 发逐 key 点名：`tool-exec-noop` 22、`approval-waiter` 13、`agent-task-noop` **0** | **达** |
| ③ 修法只许"有判据的等待＋超时上界"，禁 `time.Sleep` 糊窗口／`Skip`／放宽阈值 | `f06a8d0` 那版：三枚腿 rendezvous ＋ `NewTimeout(2s)` 轮询 ＋ 超时要红；断言 `!= 3`／`want 3` 一字未松、`SKIP=0`、新增 `time.Sleep` **0 处**、只改那 1 枚测试件 | **达**（但**归属**见 §7：码不是我落的） |
| ④ 修完同一发复跑 ≥20 发命中 0，两次逐名读数留证 | 改后 30 发 **0/30**；两批名册 65/65 **双向对称差为空**、`RUN/PASS/FAIL/SKIP`＝65/65/0/0 每发 | **达** |

**与前一程 30 份的一致性**：**一致**（2/30 vs 2/30，名册对称差为空，`goroutine.go`＋`goroutine_test.go` md5 逐一相同）；
差异只有一处且已归因：它两枚红都在 `:33`、我抓到一枚 `:29`——两跳概率本就不同（探针 5.2% vs 1.0%），不是矛盾。

**本格该给的门禁四件**：gofmt 空／gofumpt **v0.12.0 (go1.27.1)** 空／宿主 `go vet ./internal/observe/` rc=0／
`-count=2 -v` ＝ **130 / 130 / 0 / 0**。全仓门禁**未跑**（不在本格）。

**我这一格没做的事**：没写代码、没勾票面、没动 `docs/reports/**`、没动 `136-ac11-impl.md`／
`136-ac11-r1-acceptance.md`／`136-ac11-orchestrator-readings.md` 任一枚、没 `push`、没杀进程。
未核清单＝§9 那 8 条（其中"CI 侧与 linux 半边"两条是**派单明令不跑**的，其余是真未量）。

**临时件清单（全部 `/d/tmp`、前缀 `wisp136ac11b-`、只建不删）**：
`tree0/`（改前归档树，内含那枚探针副本 `internal/observe/zz_probe136ac11b_test.go`）、
`tree-after/`（改后归档树）、`run-batch.sh`、`summarize.py`、`gate.sh`、
`batch-before-DISCARDED-contended/`（争用作废批，30 份日志＋index）、`batch-before/`（30＋index）、
`batch-after/`（30＋index）、`batch-before.out`、`b0discarded-roster.txt`、`b1-roster.txt`、`a1-roster.txt`、
`b1-names.txt`、`a1-names.txt`、`r1-names.txt`（名册差集用）、`gate-poll1.txt`、`gate-postbatches.txt`、
`probe-copy_test.go`、`probe-whichkey.log`、`count2-after.log`、
`foreign-gitstatus-1.txt`、`foreign-md5-1.txt`、`foreign-worktree-edit.patch`（0 字节，成因见 §7）。
前一程的 `wisp136ac11-*` **一枚未改、未删、未覆盖**。

next= 本文件已由 `Write` 创建并只以显式 pathspec 提交这一枚路径；AC#11 的主表与勾归在飞的程与编排者，
      派单者若要收这一格，请把本表当**第二 witness**并读 §7 的归属冲突说明。
