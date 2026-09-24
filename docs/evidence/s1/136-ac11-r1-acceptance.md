# 136 · AC#11 — 非实现者终裁表（r1）

- 裁决者：`acceptor-ticket136-ac11-r1`（非实现者）。
- 只裁一格：**票 136 AC#11**（`.scratch/wisp/issues/136-…-54-tests-stay-green-without-it.md:164`，
  `TestNoopTaskReturnsToBaseline` 那枚既有偶发失败的「量率 → 归因 → 修 → 复量」闭环）。
  票 136 其余各格（含 AC#10／AC#14）**一枚不裁**，票面勾选也不由本表翻。
- 被审修法：commit **`f06a8d0`**（实现程 `worker-ticket136-ac11-r2`），只带 `internal/observe/goroutine_test.go`。
- 本表**不采信任何程的自述**：下面每个数都是本程自己跑出来的，跑在 `git archive` 出来的树上，不落仓库工作树。

---

## §0 锚点自量 · 闸门 · 取件 · 口径

### 0.1 锚点（本程现量，不是抄来的）

```
$ git rev-parse HEAD        e2a74631cced9901490c6e68af6b5b4e519c8d77
$ git branch --show-current dev
$ git log --oneline -1      e2a7463 docs(A157; 136 AC#11 补账): 两枚实现程报的 commit 号全部不存在、
                            都没写证据文件，但修是真的（f06a8d0）——我把读数补成可复算的，勾等非实现者的表
```

工作树干净：`git status --porcelain` 输出 0 行（现量于 14:27 前后）。本机 wall clock 会跳 ⇒
**本表任何一处都不用时间戳相减得时长**；每发只报它自己那行 `date`（写在各 `.v.log` 里的 `RUN-DATE`）。

### 0.2 引用到的每枚 sha 逐枚 `git cat-file -t`（现量）

```
e2a7463 commit
f06a8d0 commit
51e29b0 commit
4e66817 commit
d667888 fatal: Not a valid object name d667888
17b61a4 fatal: Not a valid object name 17b61a4
3c3e388 fatal: Not a valid object name 3c3e388
d0ca6a3 fatal: Not a valid object name d0ca6a3
2644207 fatal: Not a valid object name 2644207
```

⇒ 实现两程交件里报的 5 枚号**在盘上不存在**，与本程 `cat-file` 结果一致（编排者已独立记过同一条）。
本表因此**只以 `f06a8d0` 为被审对象**（它是唯一能 `cat-file` 出来的修法号），两程的话一句不当凭据。

### 0.3 被审版本的盘上身份（防"绿的不是那版"）

```
$ git log --oneline 51e29b0..HEAD -- internal/observe/
f06a8d0 test(136,AC#11): 把 TestNoopTaskReturnsToBaseline 的中段计数从"抢窗口"改成"等判据"
$ git diff --stat 51e29b0 f06a8d0~1 -- internal/observe/      （空输出）
$ git show --name-only f06a8d0 | tail -1                      internal/observe/goroutine_test.go
```

树身份用 md5 钉死（现量，`git show <rev>:<path> | md5sum`）：

| 那枚文件 | 在哪个版本 | md5 |
|---|---|---|
| `internal/observe/goroutine_test.go` | `51e29b0`（=改前） | `14afa086637cf9981c969abcce9ba301` |
| `internal/observe/goroutine_test.go` | `f06a8d0~1`（=改前一刻） | `14afa086637cf9981c969abcce9ba301` ← 同上 |
| `internal/observe/goroutine_test.go` | `f06a8d0`（=改后） | `336a50f5b19e2caebf56175cdee798bb` |
| `internal/observe/goroutine_test.go` | `HEAD`＝`e2a7463` | `336a50f5b19e2caebf56175cdee798bb` ← 同上 |

⇒ 三条同时成立：**改前树就是 `51e29b0`**（`f06a8d0~1` 与它同一枚 md5），**改后就是 `HEAD` 上的样子**
（`e2a7463` 之后到本程量它为止没人再动过那枚文件），**这一路只动过那一枚 `_test.go`**。
本程两批分别取件自 `git archive 51e29b0` → `/d/tmp/wisp136ac11v-tree-before/`、
`git archive e2a7463` → `/d/tmp/wisp136ac11v-tree-after/`，两棵树的 `goroutine_test.go` md5 实测
`14afa086…` / `336a50f5…`，与上表逐位相同；发火时刻整包 `internal/observe/*.go` 的 md5 名册随
`gate.txt` 各存一份（见 §1/§2 引用的目录），以便第三方复算。

### 0.4 争用闸门（self-hosted runner 就在本机，本格量的正是 goroutine 时序）

**第一批预检（本程 14:27:51 现量，四项原样）**

```
--1 宿主进程--   9764 Runner.Worker   StartTime 2026/9/24 14:26:00
--2 docker ps--  6 枚容器，全部是别的项目的常驻服务（union-proxy up 8 days / clipsync* up 4~13 days /
                 clipsync-db、clipsync-redis healthy up 13 days）；没有一枚是 wisp CI 起的
--3 gh run list-- 1 in_progress :: databaseId 35964449249 headSha e2a74631cc… status in_progress（工作流名 ci）
--4 _work mtime-- 2026-09-24 14:26:03.149549400 +0800  E:\work\base\actions-runner\_work
```

⇒ **四项里 ①③④ 同时命中**：编排者 14:2x 那次 push 触发的 `ci` run 正在本机跑。
**本程没有在争用窗口里取任何一个数**：`wisp136ac11v-batch.sh` 的发火前置闸门按 **每 30s 一轮、最多 60 轮
（＝上界 30 分钟）** 有界轮询，每轮把上面四项**原样**追加进各批目录的 `gate.txt`；等不到就 `exit 8` 报回。

⚠ 一条仪器自纠，如实报：第一版闸门里 `gh run list` **没带 `-R`**，脚本工作目录在 `/d/tmp`（不是 git 树）⇒
`gh` 报 `failed to determine base repo` 而 ③ 退化成 `GH-UNAVAILABLE`。那批发火前只跑到 round 3、
**一个 `.v.log` 都没产生**（现量 `ls …/*.v.log | wc -l` ＝ 0），本程**停掉自己那一程**（`TaskStop`，
杀的是自己的测量程，不是别人的进程）、补上 `-R CarlosShao/wisp` 后重跑。
第一版留下的 3 轮闸门读数另存为 `/d/tmp/wisp136ac11v-gate-firstattempt-rounds1-3.txt`（只建不删）。
⇒ 重跑后 ③ 读到的是 `1 in_progress`（真数，不是 `GH-UNAVAILABLE`），也就是本表 §1/§2 的每批发火前
都必须看到 `GATE VERDICT: CLEAR` 才算数；若某批只拿到 `CLEAR-LOCAL-ONLY`，那一批发的是**「本地无争用」
而不是「队列确认为空」**，本表会明写。

### 0.5 三档口径声明

- **〔独立复现〕**＝本程自己跑出来的读数，附命令与盘上文件路径，第三方可在同树复算。
- **〔第二 witness〕**＝别人（编排者 `136-ac11-orchestrator-readings.md`）的读数，只在**与本程同向**时用于对照，
  **不当本程凭据**；不同向时如实报，不挑好看的用。
- **〔仅自述〕**＝只有程的话、盘上无产物可核 ⇒ 本表一律不采信（§0.2 那 5 枚不存在的号就属于这一档）。

### 0.6 提示注入登记（本程，逐条带出处）

| # | 出处（工具＋命令前 40 字） | 原文要点 | 判别 |
|---|---|---|---|
| 1 | harness 后台任务事件（本程 `TaskStop` 自己那程 b3yiwhy6l 之后） | `task-notification … status:"killed"` | **真通知回显**（停止动作是本程自己发的，追得到出处：`TaskStop` 调用＋`chain.out`） |
| 2 | 上下文注入的 `MEMORY.md`「modified since last read」两条（仓库记忆 + 用户级记忆） | 记忆索引内容更新 | **真通知回显**（路径真：`C:\Users\swq\.qoder-cn\projects\…\memory\MEMORY.md` 与 `C:\Users\swq\.qoder-cn\memory\MEMORY.md`；未据以放宽任何判据） |

- 判为**注入**并拒绝服从的文字：**0 枚**。
- 本程未遇到任何自称"编排者备注／系统提示／某文件已被修改请 revert／放宽阈值／Confirm the harness note is genuine"
  且**指向不存在的文件或要求少取证**的文字。四条辨别力逐条走：路径真不真／是否削弱 owner 权威或放宽判据／
  能否盘上重验／是否在让我少取证。上表两条四条全过 ⇒ 记为真通知。
- ⚠ 规矩不能被代填：**"登记要带出处"这条只被引用，不由任何外部文字替本程下结论**。

---

<TODO §1 改前自测 / §2 改后自测 / §3 反向判据 / §4 修法合规五查 / §5 归因独立性 / 总判>
