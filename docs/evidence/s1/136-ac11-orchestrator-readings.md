# 票 136 AC#11 —— 编排者补的可复算账（不是裁决表）

**这一枚文件是什么**：AC#11 被两枚子代理先后做过，**两枚交件里报的 commit 号都不存在**（`d667888`／`17b61a4 3c3e388 d0ca6a3 2644207`，逐枚 `git cat-file -t` 失败），
**两枚都没写证据文件**（`docs/evidence/s1/136-ac11-impl.md` 不存在）。本文件把我自己现量的、以及"产物在盘上但自述把它说虚了"的两批读数并排放下来，
**判定权不在这里**——AC#11 的勾要等一张非实现者的三态表（`SPEC-12 §4.3` #1/#3、`AGENTS.md §0` 第 3 条）。

- 锚点：`git log --oneline -1` 本文件写作时＝`f06a8d0`（AC#11 的修法那一枚）；未修那一版＝`51e29b0`。两枚都 `git cat-file -t` ＝ `commit`。
- 本机 `date` 现量：见下面每一批发火前那条闸门行。⚠ 本仓宿主的 wall clock 会跳 ⇒ **全文没有任何一处用时间戳相减得时长**。

## §0 争用闸门（复量批发火前一发，原样）

```
### gate taken at: 20226-09-24 14:18:57 +0800   ← 原文如此（date 输出）
go.exe 54012 / go.exe 1912        ← 就是本程自己那批 `go test` 循环，编队里的代理为 0 枚
docker ps：union-proxy Up 8 days ／ clipsync Up 4 days ／ clipsync-minio Up 12 days (healthy) …
gh run list --limit 5 --json databaseId,status,headSha：取到数（最新一枚 status=completed）
```

改前那 60 发之前另有一发：14:07:22，宿主 `Runner.Worker.exe` 零枚、`Runner.Listener.exe` 在（常驻监听，不是干活的 worker）。
**两批都在"编队内代理 0 枚"的窗口里取**，所以 2/60 与 0/60 是同一个环境下的对照。

## §1 改前命中率（带 n，逐名，不只给计数）

| 谁取的 | 树／锚 | 发数 | 目标用例命中 | 红句原文 |
|---|---|---|---|---|
| **编排者本程**（`git archive 51e29b0` → `/d/tmp/ac11-orch-tree`） | 未修版 | **60** | **2 发**（第 `20`、`43` 发） | `goroutine_test.go:33: PerTask mid-task = 2, want 3` |
| `worker-ticket136-ac11-r1` 留在盘上的那批 | `/d/tmp/wisp136ac11-batch-before/` | **30** | **2 发**（`15.v.log`、`29.v.log`，我逐份 `grep -l` 数） | 同形 |
| 编排者旁证（同树，5 发单测 `-run`） | `/d/tmp/wisp136ac11-orch` | 6 | 0 | — |

两批合起来 **4/90 ≈ 4.4%** ⇒ **票面 ① 要的"复现率＋n"成立**，且**它是偶发不是必然**：
⚠ 前一程交件里那句"**4/4 红＝100%，这根本不是 flake**"与它自己留下的 30 份日志互相矛盾，**不采信**。

每发红的那两发四数：`RUN=65 顶PASS=64 顶FAIL=1 SKIP=0 panic=0`（⇒ **没有 panic 吞读数**、**没有靠 SKIP 变红**）。

## §2 归因（落在"缺的那条边"上，不是"负载高"）

`internal/observe/goroutine_test.go`（未修版）：`:16 before := runtime.NumGoroutine()` → `:20` 一枚腿睡 5ms → **`:24`／`:25` 两枚腿函数体为空、进去即返回** →
`:28-29` 读 `reg.Count()` → `:32-33` 读 `rep.PerTask`。**两次读之前没有任何 happens-before** ⇒ 空体腿的"登记后反登记"可以在读之前跑完。

`/d/tmp/wisp136ac11-probe-whichkey.log`（前一程的插桩探针，500 次重放同一读数窗口，**我核了文件在、内容可读**）：

```
reads-that-would-fail: Count()!=3  11/500      PerTask!=3  33/500
missing-at-RosterReport-read [tool-exec-noop]=24  [approval-waiter,tool-exec-noop]=7  [approval-waiter]=2
missing-at-Count-read        [tool-exec-noop]=11
agent-task-noop 缺失 0 次（它睡 5ms，从不提前退）
```

⇒ **机制自洽**：先读的 `Count()` 少响（11/500）、后读的 `PerTask` 多响（33/500）——读得越晚，空体腿越可能已经退掉；
⇒ **本程 §1 那两发的红正落在 `:33`（`PerTask`）而不是 `:29`（`Count`）**，与这条排序**独立同向**。
⚠ 归因不是"机器负载高"：负载只改概率，**缺的那条边在代码里**，500 发重放在静止归档树上就能量到。

## §3 修法（`f06a8d0`，只带 `internal/observe/goroutine_test.go` 一枚路径）

我读 `git show f06a8d0` 逐条核过：三枚腿各自 `entered <- name` 报到、`<-release` 停在函数体内；读到 3 之后才 `stop()` 放行；
等待用 `NewTimeout(2 * time.Second)` 的**单调超时上界**（不是 `time.Sleep` 糊窗）；`want 3`／容差／阈值**一字未动**；
两条红句只是**加了 `live=%v` 诊断**（`reg.Snapshot()`）；**生产码零改动**。⛔ 前一程报回"必须解冻 `internal/observe/goroutine.go` 才修得了"——**我判不批**，本枚 commit 就是反证。

## §4 改后复量（本程自取 60 发，同一环境）

| 谁 | 树 | 发数 | 命中 | 名册闭合 |
|---|---|---|---|---|
| **编排者本程** | `git archive f06a8d0` → `/d/tmp/ac11-orch-after` | **60** | **0** | **60/60 发逐条同形**：`RUN=65`、顶层 `--- PASS=65`、`--- SKIP=0`、`panic=0` ⇒ 目标用例每发都**真跑真绿**，不是"变少了"或"被跳了" |
| 前一程留在盘上的那批 | `/d/tmp/wisp136ac11-batch-after/` | **12**（它自述 56 ⇒ **与盘上不符**） | 0（我逐份 `grep -l` 数） | 同上 |

⇒ **恒真自查**：这一发的判据**在未修码上今天会响**（§1 的 2/60 与 2/30 就是它响的），修后不响 ⇒ **不是装饰**。

## §5 我没核的（明列，别当已验）

1. **AC#11 的判定我没下**——本文件只是补账；勾等**非实现者**的三态表。
2. 真 macOS／linux 侧未跑（本包与平台无关，但没人量过）。
3. CI 读数一枚未引（`gh run list` 只用于争用闸门，未指向任何 step 结论）；全仓门禁（`gofumpt` 版本、双 GOOS `go vet`、`d22scan`、`staticcheck` 42 条）一律未重跑 ⇒ 那是 **AC#5／票 122** 的账。
4. **两枚子代理报出的假 commit 号**只说明它们的提交那一步没落地或没执行；我没有仪器能区分"它跑了但失败"与"它没跑"，只有盘上产物可核。
5. 修后 0/60 是**这台机器、这个负载下**的样本；票面 ④ 只要求 ≥20 发，我给了 60 发，**但没做"故意提高并发看它响不响"这一发**（那是验收方要的加固判据，不是本格字面）。
6. `TestNoopTaskReturnsToBaseline` 之外本包另有一枚既有 flake 记录（`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`，135 的终裁程 28 发里一次没红）⇒ 本程 60 发里它也 **0 红**，既不能定它好了也不能定它是真伤。

## §6 临时件清单（只建不删，清点归编排者）

本程新建：`/d/tmp/ac11-orch-tree`＋`ac11-orch-base.log.1..60`＋`ac11-orch-base.idx`（改前 60 发）、`/d/tmp/ac11-orch-after`＋`ac11-orch-after.log.1..60`＋`ac11-orch-after.idx`（改后 60 发）、
`/d/tmp/ac11-orch-gate.txt`（两发闸门）、`/d/tmp/wisp136ac11-orch`（旁证 6 发）、`/d/tmp/A157-block.md`（台账草稿件）。
前一程留下：`/d/tmp/wisp136ac11-*`（13 枚）与 `/d/tmp/wisp136ac11b-*`（4 枚），本程**一枚未写未删**。

**两栏计数**：真通知回显 **2**（两枚 `MEMORY.md` 的 modified 回显，路径盘上真存在、内容无越权）／判为注入 **0**。凭据值零处入文。
