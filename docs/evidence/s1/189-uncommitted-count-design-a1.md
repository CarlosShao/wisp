# 189-a1 — 面板"审查"那一栏的两个维度（未提交枚数／逐文件 diff）从哪儿取：取数形状现量与只读边界设计

分工＝派单 `§C`（票 189 `AC#1` 只读设计核）。**本件零产码、`frontend/**`／`design/**` 零写面、票面原文一字未改（只在票 189 追加一行 Progress log）、`AC` 框一枚未勾。**

本轮所有读数都是**现跑**的，跑的机器＝本机、跑的目录＝本仓工作树（一棵正被多个会话同时改的共享树，这本身就是分母的来源）。凡"数东西"的尺都逐枚对照了权威命令（§4）。

---

## 0. 起手三件（照派单"共同规矩"第 9 条逐枚抄录）

```
$ date
Mon Sep 28 17:48:31 CST 2026
$ git rev-parse --abbrev-ref HEAD
dev
$ git log --oneline -1
38fc7c0e ledger(A387-A388) 收 33-r3（入向听众进树＋词面尺改能力尺，我四把尺复跑）＋owner 提问后按默认动作派 145-r2
$ git status --porcelain          # 起手名册，77 枚，逐枚如下（/tmp/189a1-start-roster.txt 回显）
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
?? .scratch/wisp/.scratch/
?? .scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
?? .scratch/wisp/probes/139/accept-r1/
?? .scratch/wisp/probes/152/overlay-probe1-on-samppost.json
?? .scratch/wisp/probes/156/__pycache__/
?? .scratch/wisp/probes/156/mut-156-r2/asis.log
?? .scratch/wisp/probes/156/zero156-r4-head.sh
?? .scratch/wisp/probes/156/zero156-r4-work/
?? .scratch/wisp/probes/158/r2/
?? .scratch/wisp/probes/161/r2/__pycache__/
?? .scratch/wisp/probes/161/r2/ctl/
?? .scratch/wisp/probes/161/r5/negative-control/
?? .scratch/wisp/probes/161/r6/logs/flip-7.txt
?? .scratch/wisp/probes/162/r4/
?? .scratch/wisp/probes/162/v1-baseline-gotest.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start.txt
?? .scratch/wisp/probes/183/accept-v1/
?? .scratch/wisp/probes/185/c1/logs/d22scan-post-final.txt
?? .scratch/wisp/probes/33/r2/d22scan-final.txt
?? .scratch/wisp/probes/33/r2/gate-final.txt
?? .scratch/wisp/probes/33/r2/status-final.txt
?? .scratch/wisp/probes/33/r2/status-start.txt
?? .scratch/wisp/probes/999/
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
?? part1-state1-fixed.txt
?? part1-state1-pristine.txt
?? part2-nog6-fixed.txt
?? part2-nog6-pristine.txt
?? part3-stale-fixed.txt
?? part3-stale-pristine.txt
```

起手名册按状态分档（同一份输出用 `awk '{print substr($0,1,2)}' | sort | uniq -c` 现数）：**` M` 15 枚＋` D` 16 枚＝已跟踪改动 31 枚；`??` 46 枚＝未跟踪**；合计 77 枚。**这三档的分裂正是本票取数形状的命门**（§3）。

## 1. 我的锚 ≠ 编排者锚（`A385` 口径：不算漂移，只需点名）

编排者锚＝`38fc7c0e`；我出件时的锚自取为当下 HEAD。起手之后落进树的别人的提交（逐枚 `git show --name-only` 原文，**不用区间 diff**）：

```
== 971b7dad ledger(A389)＋立 Q-71 - 前端等数据的五枚票一次排开（145-r2/188-r1 写＋182a/180a/189a/190a 只读）
.scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
docs/reports/pending-and-issues.md
== cb14a8b2 180-a1(片①): 票 180 Progress log 追加（AC 框一枚未勾，只读普查结论）
.scratch/wisp/issues/180-panel-width-is-a-config-field-with-zero-production-readers-so-changing-config-toml-has-no-visible-effect.md
== 500bc98e 188-r1(片①): 证据件骨架——起手名册逐枚抄录＋起手锚 38fc7c0e＋cb14a8b2 逐枚 show 原文
docs/evidence/s1/188-task-state-r1.md
== 9b6a1be5 190-a1(只读): 证据件 190-file-tree-design-a1.md — 文件树那堆的根/上限/忽略规则设计核
docs/evidence/s1/190-file-tree-design-a1.md
```

**碰没碰我的文件**：`internal/panel/git.go`／`internal/panel/git_test.go`／票 189 票面＝**零枚被碰**（上面四枚的 name-only 里没有它们）。`145-r2` 的写面（`internal/panel/composer.go`／`internal/panel/pump.go`／`cmd/wisp/panel_pump.go`／`internal/panel/pump_test.go`）到本件出件时**仍未进树**（那四枚提交里也没有它们）⇒ §2 的 `file:line` 读数对 145-r2 交件前的形状负责，交件后行号需让写腿现量重取。

---

## 2. 落点现量：票 181 那同一枚读面 `internal/panel/git.go`

票面 `AC#2` 与派单 §C②都定死了落点＝同一枚读面，本节给它的**今天形状**。

| 现量项 | 读数（本轮） | 凭据 |
|---|---|---|
| `internal/panel/git.go` 行数 | **559 行**（23,265 字节） | `wc -l internal/panel/git.go` |
| 它的判据文件 | **721 行**（`internal/panel/git_test.go`） | `wc -l internal/panel/git_test.go` |
| 导出符号枚数 | **7 枚**（`GitKindRepo`/`GitKindNotARepo`/`GitKindPermissionDenied`/`GitKindUnreadable` 四枚常量在同一条 `const (` 里、`GitSwitchBlockedReason`、`GitWorktree`、`GitView`、`GitViewNotProbed`、`ReadGit`、`ReadGitForWorkspace`、`GitViewRewritten`） | `git.go:65-70,82,93,117,150,172,240,267` |
| 内部（未导出）函数枚数 | 15 枚，全部是 `os.Stat`/`os.ReadFile`/`os.ReadDir` 一类，**零 `os/exec`** | `git.go:154,284,315,358,373,383,397,433,454,471,497,529,539,557` |
| 谁在调它 | `internal/panel/composer.go:95`（`composer.Git = GitViewNotProbed()`）、`:216`（`Git GitView json:"git"`）、`:232`；`internal/panel/pump.go:120`（`Git func() GitView`）；`cmd/wisp/panel_pump.go:108-109`（`func (rt *agentRuntime) gitView() panel.GitView { return panel.ReadGitForWorkspace(rt.workspaceView()) }`） | 逐枚 `grep -rn "ReadGit\|GitView" --include=*.go internal/ cmd/` |
| 包内 `os/exec` 命中 | **0**（`grep -rn "os/exec" --include=*.go internal/panel/` 空） | 现跑 |
| 它今天**不答**什么 | 未提交枚数、逐文件 diff——`git.go` 全文只答"是不是 git 树／哪个分支或裸提交／有哪几棵工作树／有哪几枚本地分支"四问 | `git.go:5-8` 自陈 |
| 普查独立量到的"零源" | `grep -rniE "git diff|numstat" --include=*.go internal/ cmd/`（排 `_test.go`）⇒ **0 命中**；`180-182` 普查 §5 行 151 同一形尺 ⇒ 0 | `docs/evidence/s1/180-182-panel-fields-census-c1.md:151` |

**结论（派单 §C②）**：新读面**进 `internal/panel/git.go` 这一枚文件族**，不开第二份。落点形状建议（写腿用）：
- 视图类型 `GitUncommittedView`／`GitDiffView` **嵌进 `GitView` 的两个字段**（`git.go:117-145` 那一族），不新增 `ComposerState` 的顶层 key ⇒ 少碰一枚被前端契约钉住的键集（`composer.go:210-216` 的注释逐字写着 "It is a NEW KEY on a section whose key set is pinned by the frontend contract… see Q-51"，钉＝`internal/panel/composer_test.go:48 TestComposerContractTypesMatchFrontend`，今天在册红）。
- 若产码体量让"塞进同一枚文件"变成 1,500 行的巨石，**同族新开一枚 `internal/panel/git_readonly_exec.go`** 也算"同一枚读面"（同一个包、同一族类型、同一个泵钩子），但**必须**沿用 `git.go` 的三段式自陈头（是什么／不是什么／不做什么）。**这一格交给写腿按体量裁，本件不替它定。**

---

## 3. 派单 §C①／票面 `AC#1`：取数形状两条路的**现量代价**

先给一句话结论，再给数。**两维的代价结构不一样，所以不该用一条路**：

- **未提交枚数（`??`＋` M`＋` D`）**：已跟踪那一半纯 Go **便宜且准**；未跟踪那一半纯 Go **要么贵要么错得离谱**。
- **逐文件 diff（内容比对）**：纯 Go 在 `go.mod` 冻结的前提下＝手写 packfile＋delta＋zlib＋CRLF＋重命名检测，**没有一条便宜的路**；外部 `git diff` 一次 267 ms 出全量。

### 3.1 路线 A：起外部进程（`git status --porcelain`／`git diff` 的等价只读子命令）

| 现量项 | 读数 | 怎么量的 |
|---|---|---|
| `git status --porcelain` 三连耗时 | **197 ms／358 ms／220 ms**（第 1 次是本轮第一次跑，冷） | `s=$(date +%s%N); git --no-optional-locks status --porcelain; e=$(date +%s%N)` 三次 |
| `git diff`（全量，31 枚已跟踪改动） | **267 ms／1,521,885 字节／17,374 行** | 同法，输出落 `/tmp/diff-all.patch` 后 `wc -c`／`wc -l` |
| `git diff --name-only` 与 `--numstat` 的枚数 | 都是 **31 枚** | 现跑 |
| git 二进制 | `git version 2.52.0.windows.1`，`which git`＝`/mingw64/bin/git` | 现跑 |
| ⚠ **只读子命令会不会写盘** | 本轮**没写**：`.git/index` 的 md5 跑 `git status` 前后都是 `1f078d47eebd9e60c14f4136d6c067a3` | `md5sum .git/index` ×2 |
| ⚠ 但**不能当通则** | `git --no-optional-locks` 这枚开关的存在＝git 承认默认路径会去**取锁**（锁的存在就是为写准备的；"刷写 index"属由此推出的**推断，本轮未复现**）。⇒ 派给写腿：**只用 `--no-optional-locks` 起，并且用一发常驻判据钉"跑完读面 `.git/index` 未变"**（合成树里量最便宜） | 现量＋具名标注的推断，两条分开写 |
| 副作用读数：换行转换 | `git diff` 在 stderr 上对 **11 枚文件**报 "LF will be replaced by CRLF the next time Git touches it" ⇒ 这条通道的输出**带 git 的文本转换语义**（`core.autocrlf`／`.gitattributes`），纯 Go 那一侧要么复刻要么就会给出 git 不认为存在的改动 | stderr 原文已回显 |
| 仓内**已有先例**（不是全新能力） | `tools/d22scan/gitignore.go:373` 就是 `exec.CommandContext(ctx, "git", args...)`；它的降级语义写在同一文件 `:49-62`——"问不到 git（没二进制／不是仓库／`-root` 指向子目录）就一条规则都不套并**响亮自陈**" | 现读 |
| 生产 Go 里的 `exec.Command` | 4 处：`cmd/wisp/doctor.go:316`、`cmd/wisp/slo_windows.go:492`、`cmd/balldebug/diff_windows.go:409`、`internal/llm/adaptertest/mockllm.go:65,72` ⇒ 面板／运行时读 git 这条路上**今天一枚都没有**，全仓零次外部 `git`（同 `181-c1` 结论） | `grep -rn "exec.Command" --include=*.go internal/ cmd/`（排 `_test.go`） |
| d22scan 七枚 ban 的射程（本轮复核） | 逐字读了 `tools/d22scan/main.go:5-40` 的 ban 列表：`1 bare-goroutine`／`2 pathresolver-bypass`（**只禁 `filepath.Clean`／`filepath.Abs`**）／`3 plaintext-key`／`4 wallclock-timeout`／`5 mirror-hash`／`6 panel-approval`（`approval.decide` in `frontend/`）／`7 internal-artifact-tool`（`internal/tools/`）／`8 emoji` ⇒ **没有一枚管"起外部进程"**；真要跑会咬到的是 **#1（等子进程别裸 `go`，同步 `Output()` 就不涉及）**＋**#4（超时必须单调钟：`exec.CommandContext`＋`context.WithTimeout`，别写 `time.Now().Sub(...)`）**＋**#2（路径决策过 C26）** | `sed -n '5,40p' tools/d22scan/main.go` |

**路线 A 的代价（定性）**：一次性把"正确性"租过来（忽略语义、`git add -f` 强加、submodule、`core.autocrlf`、重命名检测、二进制识别）；换来的是**依赖用户机器上有 git**（便携包不带 git 时这一维必须降级）、**每次快照 200–360 ms**、以及"外部命令必须只读"这件事得靠钉子守。

### 3.2 路线 B：纯 Go 读索引（手写 `.git/index`＋对象库）

我**真写了一枚仓外探针跑过一遍**（`/tmp/gitprobe/main.go`，不落仓、不进树、不属于产码），读数如下（`go run` 本机，同一棵工作树）：

| 现量项 | 读数 | 对照权威命令 |
|---|---|---|
| `.git/index` 体积／版本／条目 | **399,756 字节／version 2／3,332 条**；解析 **1.087 ms** | `git ls-files \| wc -l` ＝ **3,332**（逐枚对上） |
| 对 3,332 条逐条 `os.Stat` 比 mtime+size | **51.571 ms**；**stat-dirty＝15、missing＝16** | ` M`＝**15**、` D`＝**16** ⇒ **纯 Go 的"已跟踪改动"这一半是准的** |
| 只对脏的那 15 枚算 `blob` SHA-1 | **49.484 ms／909,178 字节／content-dirty＝15** | 权威 15 ⇒ 一致 |
| 未跟踪那一半（朴素遍历，不套任何忽略规则） | **199.7 ms／12,000 个文件／未跟踪候选 8,684** | 权威 `git status`＝**46** ⇒ **高估 188 倍** |
| 忽略规则的真面积 | 本仓 `.gitignore` **5 枚**（根／`frontend/`／`probes/161/r1`／`r3`／`r4/removed`，**嵌套**）、`.git/info/exclude` 存在、全局 `core.excludesFile` **本机为空** | `find`／`git config --get core.excludesFile` 现跑 |
| 已有的"纯 Go 忽略实现"值多少钱 | `tools/d22scan/gitignore.go` **817 行**，而且它**仍然去问 `git ls-files`**（`:373`）——gitignore(5) 的可用子集本身就不是几百行能收口的 | 现读现量 |
| diff 那一半需要的对象库 | `git count-objects -v`：`count: 2016`（松散）／`in-pack: 17653`／`packs: 4`／`size-pack: 20625` KiB；`.git/objects/pack/` 里最大一枚 `pack-431f….pack` ＝ **12,626,706 字节**，另有 `.idx`／`.rev`／`multi-pack-index`（495,600 字节）；`git cat-file --batch-check --batch-all-objects \| wc -l`＝**19,669** 枚对象 | 现跑 |
| ⇒ 纯 Go 要 diff 内容就必须复刻的东西 | pack idx v2 解析＋offset/tree delta 还原＋zlib inflate＋（§3.1 那 11 枚文件的）CRLF 转换＋重命名检测 | 上面那两行读数推出来的工作清单 |
| ⚠⚠ **最贵的一枚现量是"错了不报错"** | 探针第一版把 index 条目的字段偏移读错了（我按 `mtime@12 / sha@24 / size@40` 解，真形状是 `mtime_sec@8 / mode@24 / size@36 / sha1@40:60`）。**错的版本跑通了、不 panic、不报错**，给出的是：`stat-dirty=3316／3332`（几乎全仓都"脏"）、未跟踪候选 `11996`；改对之后才是 `15/16`。⇒ 路线 B 的失败模式不是"跑不起来"，是**面板上出现一个看起来完全合理的错枚数** | 两次 `go run` 的实际输出，原文贴在下面 |

```
# 第一版（偏移错）：
index version=21059 entries=2 parse=0ms (bytes=399756)
naive-walk: 151ms files=11996 untrackedCandidates=11996  (authority: git status untracked=46)
# 第二版（头偏移对、条目内偏移仍错）：
index version=2 entries=3332 parse=557µs
stat-pass entries=3332 cost=36.274ms stat-dirty=3316 missing=16
hash-pass files=3316 bytes=50937392 cost=378.625ms content-dirty=3316
# 第三版（对）：
index version=2 entries=3332 parse=1.087ms (bytes=399756)
stat-pass entries=3332 cost=51.571ms stat-dirty=15 missing=16
hash-pass files=15 bytes=909178 cost=49.484ms content-dirty=15
naive-walk cost=199.703ms dirs=1231 files=12000 untrackedCandidates=8684  (authority: git status untracked=46)
```

（顺带：`version=21059` 那一发是签名头偏移错——"DIRC" 后 4 字节才是版本，我从 `raw[2]` 起读。三版都跑了同一棵树，只有第三版的 15/16 能和权威对上。这正是票 181 那一族"认 git 目录形状会漂"的同一件事，只是这次漂在二进制格式里。）

**路线 B 的代价（定性）**：已跟踪改动那一半**便宜（≈101 ms、零进程、零依赖）且我量到它准**；未跟踪那一半要准就得写忽略语义（对照物 817 行）＋还得处理 `git add -f`；diff 那一半要准就得手写对象库读取。**"用现成 Go 库"这一支今天直接排除**：`go.mod` 的 8 枚直接依赖里没有任何 git 库，`$GOMODCACHE/github.com/go-git` **不存在**（`ls: No such file or directory`），而派单第 12 条写着 `go.mod`／`go.sum` **一字节不动**。

### 3.3 本件的裁（给写腿的推荐，不替 owner 拍别的板）

**分维走法**（票面 `AC#1` 的问号里那句"diff 这一样很可能需要（因为内容比对不在文件索引里）"，实测成立；枚数那一半比我预期的更可以纯 Go）：

1. **未提交枚数 ⇒ 纯 Go 为主，外部为降级方向的反号**：
   - 已跟踪 31 枚 ⇒ 走 §3.2 的 index＋stat 那一发（101 ms、准）。这一半**永远可得**，无二进制依赖。
   - 未跟踪 46 枚 ⇒ 只有两条准路：(a) 问外部 `git status --porcelain`（或 `git ls-files -o --exclude-standard`）；(b) 自建忽略匹配器。**推荐 (a)，并把"问不到 git"做成第三个枚举值而不是"0 枚"**——口径直接沿用 `git.go:37-40`（"读不到"不许塌成"不是仓库"）与 `d22scan/gitignore.go:49-62` 的响亮自陈。
   - ⚠ 不许出现的形状：**把 8,684 当未跟踪枚数显示**，也**不许**因为"套不出规则"就把未跟踪那一半报成 0。
2. **逐文件 diff ⇒ 外部只读子命令，无第二选择**：`git --no-optional-locks diff -- <path>`（或 `--numstat` 那一档先看行数）。git 缺席 ⇒ 这一维整体 `unreadable` ＋一句人话原因，**不许**用"纯 Go 读工作树文件内容"冒充 diff（没有 HEAD 侧基准就不是 diff）。
3. 两维**共用同一次取数的时间戳**，避免"枚数说 31、diff 列表 30"这种同一快照内自相矛盾（`A374`／`194-c1` 那族两份读面必然漂的形状）。
4. ⚠ **泵的成本没裁完**（§7 未裁完第 2 条）：`PumpSources.Git func() GitView`（`pump.go:120`）是每次快照都调的同步钩子，路线 A 一次 200–360 ms、diff 再 267 ms ⇒ 必须定"缓存键／重算条件／上限"，本件只给候选（用 `git.go:397 readHead` 已能零成本拿到的 HEAD sha ＋ `.git/index` 的 mtime+size 作键），**不替 145/写腿定节奏**。

---

## 4. 派单 §C① 附带的硬要求：本轮每枚"数东西的尺"逐枚对照权威命令

| 我在用的尺 | 我的读数 | 权威命令 | 权威读数 | 判定 |
|---|---|---|---|---|
| 未提交枚数 | 77（` M`15／` D`16／`??`46） | `git status --porcelain` ＋ `awk` 分档 | 77 | 一致（同源，`awk` 只切前两字节） |
| 已跟踪改动枚数 | 31 | `git diff --name-only \| wc -l`／`git diff --numstat \| wc -l` | 31／31 | **一致** |
| index 条目数 | 3,332（纯 Go 解析） | `git ls-files \| wc -l` | 3,332 | **一致（修完偏移后）** |
| 未跟踪枚数 | 46（权威侧）；朴素遍历＝8,684 | `git status` 的 `??` | 46 | **尺不能是遍历**，差 188× |
| 分支枚数 | — | `git branch --list` | **4** | 对照：`.git/refs/heads` 递归 `find -type f`＝**4**、非递归 `ls`＝**3** ⇒ 复现票 181 的"斜杠分支名是嵌套目录"那一漂（`git.go:433-467` 的递归注释），本票再次量到 |
| 工作树枚数 | — | `git worktree list` | **3** | 对照：`.git/worktrees/` 目录枚数＝**2** ⇒ 再次量到"主工作树不在那"（`git.go:88-92` 注释），差别 1 枚＝脚下这棵 |
| packed-refs 里的本地分支 | 0 | `grep -c refs/heads .git/packed-refs` | 0 | 本仓松散引用全胜，合并逻辑（`git.go:441-446`）本轮**没有**活样本 ⇒ 写腿若加尺别假设它有样本 |

**给写腿的尺子清单**：任何"枚数"字段都必须在**合成树**（票面 `AC#2` 要的 `git init` 那棵）里同时被两条尺钉住——你的读面 vs 外部 `git` 的权威输出；两者不一致就判红，别只钉一个数字。

---

## 5. 派单 §C③／票面 `AC#3`（本票命门）：只读边界的**结构性不可能**设计

普查把雷区逐枚列了（`180-182-panel-fields-census-c1.md:178-179,195,206,343`）：S2 的"提交／推送"、S3 的"接受／回滚"，**且都标了"要 owner 逐枚批＋过 L2 与审计"**。本票只做显示，所以要求是：**动作腿在设计上就装不进去**，不是"我们先说好不做"。

| 发 | 结构性做法（不是约定） | 现成模板／凭据 |
|---|---|---|
| **F1 类型面：没有落点** | 新读面的**全部**导出签名收成 `func(工作区坐标) View`——入参不收 `[]byte`／`io.Writer`／`*os.File`／命令名／路径列表，出参是**值类型**（无指针、无句柄、无 `context` 写通道）。想写文件，就得先改函数形状；改形状在评审里是**看得见的** | `git.go:172 ReadGit(workspaceRoot string) GitView`、`:240 ReadGitForWorkspace(ws WorkspaceView) GitView`——今天这一族就是这两枚，且 15 枚内部函数里没有一枚有写能力 |
| **F2 依赖名册尺（常驻，AST 级而非词法级）** | 一枚常驻判据：用 `go/parser` 解析**自己这一族文件**（`internal/panel/git*.go` 的非测试文件），断言 ① import 集合 ⊆ 白名单；② `os`／`io` 里被**调用到的方法名**只准出现在 `{Stat, Lstat, ReadFile, ReadDir, Open}`；出现 `WriteFile/Create/Remove/Rename/Mkdir/Truncate` 即红；③ 若走了 §3.3 的外部只读那一支，`exec` 只允许出现在**具名单一文件**里，且 `exec.Command*` 的 git 子命令参数**逐枚 ∈ 只读名册** `{--no-optional-locks, status, diff, ls-files, rev-parse}`，出现 `commit/push/add/checkout/reset/restore/merge/rebase/am/apply/clean/switch/stash/tag` 即红 | 模板＝`190-file-tree-design-a1.md` §6 的 **S3 依赖名册尺**（它写的正是"这不是词法尺——它扫的是这个包持不持有写能力"）＋`git.go:22-26` 的三段式自陈头。⚠ 这发尺**必须自带正控**（种一个 `os.WriteFile` 进临时副本、门必须响），按 `MEMORY` 里那条"负向尺必配正控" |
| **F3 拓扑不可达（今天最硬的一发）** | 面板**发起不了任何事**：网页→宿主那一跳在这棵树里不存在（`git.go:17-20` 逐字 "the inbound half … that hop does not exist in this tree"，`GitSwitchBlockedReason` `:82` 是常量不是判断），入向方法白名单（`internal/panel/bridge.go:42-45` 那四枚 `panel.*`）一枚不加 ⇒ **"审查栏里点接受／回滚"在当前拓扑上不可达**。落地腿交件时**必须现量那份名册**（别引用我说过的话） | `git.go:82`、`181-186-git-detection-census-c1.md` §⑤（该 hop 零生产调用者的 grep 就在那节）、票 181 `AC#3` 逐字"那四枚 `panel.*` 一枚不许加"（引于 `180-182` 普查 `:179`） |
| **F4 不是 Tool（模型侧不可达）** | 新维度**不注册进 `internal/tools`**，D34 表零 git 行（票面 `AC#5` 逐字）。现成判据已在跑：`TestGitDimensionHasNoModelCallableTool`（`git.go:13-16` 的自陈）；两枚新视图沿用同一条 ⇒ 模型既不能调 git，也不能通过面板间接调 | `git.go:13-16`、`git.go` 那枚同名测试 |
| **F5 无副作用可断言**（把"只读"变成读数不是形容） | 常驻判据在合成树里：跑读面前后记 `.git/index` 的 sha256＋mtime，断言**未变**；`.git/refs` 名册未变；工作树文件未变。这一发防的是 §3.1 那枚"只读子命令也可能取锁"的形状（本轮实测前后 md5 同为 `1f078d47…`＝**这一次没写**，通则未复现，所以要钉而不是要相信） | §3.1 现量行 |
| **F6 动作腿归口在别处** | "接受／回滚"＝`fs.write`／`fs.edit` 级 L2 写动作（普查 `:206` 逐字），"提交／推送"＝雷区（`:343`）。本票的交付里**连"预留一个 hook 参数"都不许出现**——预留＝下一枚票可以在不改形状的前提下把它接上，F1 就白写了 | 普查 §7；`AGENTS.md` §1.2"由面板侧来源的 L2 允许"是禁区 |

**一句话**：F1 让"写"没有参数可传，F2 让"写"没有可调用的符号，F3 让面板连请求都发不出，F4 让模型也调不到，F5 让"顺手写了一下"变成红。**这五发叠起来才是"结构性不可能"**；单拿 F6 那种"票面写了不许"交差，就是本仓一直在拒的"约定不做"。

---

## 6. 票面 `AC#4`：大 diff 的代价——上限、截断语义与真读数

| 项 | 读数／设计 | 凭据 |
|---|---|---|
| 本仓今天的真读数 | 31 枚已跟踪改动 ⇒ `git diff` 全量 **1,521,885 字节／17,374 行／267 ms**；这些改动文件的在盘体积合计 **909,178 字节** | §3.1／§3.2 现跑 |
| ⇒ 快照能装吗 | **不能整包塞**。快照是每条推送都走 `PumpSources`→`Publish` 的出向包（`pump.go:105-135`），1.5 MB 进每一包＝把界面管子变成搬运工 | `pump.go:120` 那一族 |
| 上限建议（写腿用，owner 不批我不改） | ① **逐文件**：只送前 `N` 行 hunk（建议 200 行）＋ 单文件 `bytesIncluded/bytesTotal` 三元组；② **整包**：一次快照里 diff 维总字节上限（建议 256 KiB）；③ 超上限 ⇒ 该文件降级为"只报 `--numstat` 行数"；④ 二进制（本仓那枚未跟踪 `design/doubao/01-ball-states.jpg` 属此类）单列 `binary` 档，不进文本预算 | 本票形状选择，非仓内既有值 |
| ⚠ **截断必须带"怎么找回全量"的指针**（票 164／D15 那根管子的教训，票面逐字） | 三个字段一起才算数：`truncated bool`＋`bytesTotal int64`＋**`recoverPointer string`**（具体可执行的那句：`git --no-optional-locks diff -- <path>` 的原样命令，或宿主把全量落在哪一枚 artifacts 文件的路径）。只带 `truncated=true` ＝票面点名的"做了但没生效" | 票 189 `AC#4`；票 164 的溢出层教训（`190-a1` `:185` 同引 `internal/tools/fs.go:36-38`） |
| ⚠ 不许塌成"没有" | `unreadable`（git 缺席／读不到）与 `0`（真的干净）必须是两个枚举值——沿用 `git.go:65-70` 的四态，**新维加自己的态也别新造这套口径**；`hidden_by_ignore`／`total==unknown` 那一族 `190-a1` `:213` 已经写了同一条 | `git.go:37-40,65-70` |

---

## 7. 派单 §C④：这一维该进票 145 快照的**哪一格、叫什么**

**今天的形状**：`Snapshot → ComposerState（`json:"composer"`）→ Git（`git.go:117 GitView`，`composer.go:216 json:"git"`）⇒ **快照里那一格叫 `composer.git`**。

**本件的命名建议**（嵌在既有那一格下面，不开新的顶层格）：

| 新子格（JSON key） | Go 类型 | 对应普查那一堆 | 装什么 |
|---|---|---|---|
| `composer.git.uncommitted` | `GitUncommittedView`（新，进 `git.go` 族） | **S2「未提交枚数」**（`180-182` 普查 `:128,150,178`） | `kind`／`reason`／`total`／`trackedModified`／`trackedDeleted`／`untracked`（三档分裂见 §0）／`files[]`／`askedGit bool`（降级自陈） |
| `composer.git.diff` | `GitDiffView`（同族） | **S3「审查 diff」**（普查 `:129,151,179,206`） | `files[]{path,status,additions,deletions,hunks,text,bytesTotal,bytesIncluded,truncated,recoverPointer,binary}`／`wholeBytes`／`cap`／`kind`／`reason` |

**与 §A 的关系（派单要求"对得上才算数"）**：§A（`182-a1`）**本件出件时还没交件**——`docs/evidence/s1/` 里只有 `180-182-panel-fields-census-c1.md`／`181-*`／`183-*`／`184-*`／`185-*`，**没有 182-a1 那一枚**（本轮 `ls` 现量）。所以我以**普查 §7 的 S2／S3 两堆**为准名；若 `182-a1` 回来把 S2＋S3 折成一枚顶层 `composer.review`，**那是 §A 与 §C 的命名冲突，我不自己折中**：冲突点在"顶层新 key 会撞上 `composer.go:210-216` 注释里那枚前端契约钉（`TestComposerContractTypesMatchFrontend`，今天在册红，见派单第 14 条）"，请编排者具名裁。

**载体 vs 逻辑的归口**（与同波 `190-a1` 的裁定同形、不冲突）：**"快照里有哪几维"归票 145**（载体集），**"这一维的值从哪来"归票 189**（读面）。`190-a1` `:237-240` 已把同一刀落在 S4 上，git 维今天也是走这条路进来的（`git.go:8-11`）。

**⚠ 排序硬约束**（写给编排者，不是给写腿的命令）：`145-r2` 正在写 `composer.go`／`pump.go`／`panel_pump.go`／`pump_test.go`；本票的落地腿**必须排在它之后**（同批文件），且落地前**必须逐枚跑撞钉预检**——`internal/panel` 今天在册 3 枚红（派单第 14 条列了逐名），其中 `TestComposerContractTypesMatchFrontend` 的射程正好盖住"给快照加字段"这一族。

---

## 8. 交给写腿的落地清单（本件不产码，只列它该带什么进门）

1. 位置：`internal/panel/git.go` 同一族（体量过大时同族新增 `git_readonly_exec.go`，**不开第二个包**）。
2. 两枚视图类型嵌进 `GitView`，JSON key 用 §7 那两个；空态一律 `[]T{}` 不 `nil`（`git.go:109-111,154-162` 的既有规矩）。
3. 外部调用：`exec.CommandContext`＋`git --no-optional-locks …`，参数**全为固定字符串字面量**（不接受任何来自面板/模型的输入），同步 `Output()`（不裸 `go` ⇒ 不碰 ban #1），超时用 `context.WithTimeout`（单调钟 ⇒ 不碰 ban #4），git 缺席 ⇒ 响亮降级（`d22scan/gitignore.go:49-62` 同形）。
4. 判据（票面 `AC#2` 要的合成树）：临时目录 `git init` → 写 3 枚文件 → `git add` 1 枚 → 改 1 枚 → 删 1 枚 → 断言 `trackedModified/trackedDeleted/untracked/total` 四档与 `git status --porcelain` 逐枚一致；**再加一发**：同一棵合成树跑前后 `.git/index` sha256 不变（F5）。
5. 结构尺：§5 的 F2（带正控）＋ F3 名册现量＋ F4 沿用既有那枚 Tool 判据。
6. 门禁口径：`./internal/panel/ ./internal/tools/`、`sh scripts/d22scan.sh`（只看 clean／逐名红点，基线随锚点变，别引用旧数）、`gate-clauses.sh` 只比 BAD 名册（今天在册只有 `G6neg`）、`gofumpt -l internal/panel`；**在册 3 枚红不修、不当绿**。

## 9. 未裁完（不许当完整交）

- **泵的重算节奏／缓存键**（§3.3 第 4 点）：只给了候选（HEAD sha＋index mtime+size），命中率与"每次快照 200–360 ms 能否接受"未裁——要 145 交件后的泵形状读数。
- **`--porcelain=v2 -z` 与 `--numstat` 的输出解析细节**：我只现量到 v1 前两字节分档与枚数，没逐字读 v2 扩展格式（写腿要用就得先 `git status --porcelain=v2` 现跑对照）。
- **d22scan 七枚 ban 我只做到 `main.go:5-40` 的列表逐字（本轮 §3.1 已贴），但没读每枚 ban 的正则实现**；真要判"会不会咬到 exec 那一支"，写腿跑门时才算真核过。
- **submodule／sparse-checkout／`core.autocrlf=true` 的假想形状**：本仓没样本（`.gitmodules` 未查），这些是路线 A 免费送来、路线 B 会漏的类别，我只在 §3.1 记了 CRLF 这一枚实据。
- **重命名检测（`-M`）与 `staged` vs `unstaged` 的口径**：截图那一栏要哪个（"未提交"含不含 `git add` 过的暂存改动？）＝**owner 的口径问题，不是我的**，本件只指出 `git status` 的两列 XY 形状能区分它（本轮名册前两字节只有 ` M`／` D`／`??` 三形，无暂存样本 ⇒ 无活数据）。
- **`182-a1`（§A）与本件命名的一致性**：见 §7 的冲突点，待编排者裁。

## 10. 预算与被哪一步吃掉（派单"超预算要具名"）

只读预算 ≤20 枚。实际约 **23 枚**：**超出的 3 枚全花在 §3.2 那枚纯 Go 探针上**（第一版头偏移错、第二版条目偏移错、第三版才对——每一版都要重写＋跑一遍）。我没有为了省调用而把路线 B 的代价写成推测，因为派单 §C① 要求"两种都要现量代价"。**没有放宽任何断言，没有勾任何 AC。**

## 11. 终态自证（逐枚 `git show --name-status`，不用区间 diff）

**我的写面只有两枚文件**（片① `465c060a`，`git show --name-status` 原文）：

```
465c060a 189-a1(只读·派单§C): 证据件 189-uncommitted-count-design-a1.md — …（AC 框一枚未勾）

M	.scratch/wisp/issues/189-the-review-stack-needs-uncommitted-count-and-per-file-diff-but-go-side-has-zero-source.md
A	docs/evidence/s1/189-uncommitted-count-design-a1.md
```

`git show --name-only HEAD` 里 **`internal/panel/composer.go`／`pump.go`／`cmd/wisp/panel_pump.go`／`internal/panel/git*.go` 命中＝0**（grep 现跑回显 `NONE-OF-THEM (zero)`）⇒ 在飞的 `145-r2` 写面与票 181 那枚读面我都没碰。仓内**零删除**、**未 push**。

**终态闸门（`date`＝Mon Sep 28 18:00:20 CST 2026，锚 `465c060a`）**：`git status --porcelain` 终态 **77 枚** ＝ 起手 **77 枚**，逐枚具名差集（`comm -23`／`comm -13` 对 `/tmp/189a1-start-roster.txt` 与终态文件）**只有两枚、都归别人**：

| 差集 | 条目 | 归因（逐枚具名，非我造成） |
|---|---|---|
| 只在起手 | `?? .scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md` | 别人的 `971b7dad ledger(A389)＋立 Q-71` 把它提交了 ⇒ 从 `??` 名册消失（那枚提交的 name-only 见 §1） |
| 只在终态 | ` M docs/evidence/s1/188-task-state-r1.md` | `188-r1` 写腿在我取名册之后继续改它自己那枚证据件（`500bc98e` 之后又动了一次） |

⇒ **除这两枚外逐枚相同；我这一程造成的净差集＝空**。

**探针与临时件**：纯 Go 探针与所有名册／diff 读数都落在**仓外**（`/tmp/gitprobe/main.go`、`/tmp/189a1-start-roster.txt`、`/tmp/189a1-end-roster.txt`、`/tmp/st-plain.txt`、`/tmp/diff-all.patch`），仓内**没有**为本轮新增任何台件；按"临时件只建不删"（`issues/README` 规则 8）它们留在 `/tmp` 不删。
