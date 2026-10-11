# 票 304 — 测试二进制被外力杀掉时，它起的子进程**不会跟着退**：机主台面上现量躺着 2 枚 `mockllm.exe`（最老那枚已经活了 **57 小时**），而编队的起手闸门⛔ 看不见这一族

**立票**：2026-10-10 17:4x 编排者（来路＝为销掉台账 `A820` 那笔"24 枚孤儿 webview"欠账而做的父链取证；取证**否证了那笔账**，却抓到一枚真的残留——全过程逐字在 `.scratch/wisp/probes/orch/webview-orphans/`，见 `A821`）
**性质**：★**缺陷票（测试夹具的子进程生命周期）**。⛔ 产码缺陷，⛔ 归口票。
**为什么要紧**：① 编队每程起手都跑 `tasklist` 数 `wisp.exe`／`balldebug.exe` 当"桌面干净"的凭据——**`mockllm.exe` 不在这两名册里**，所以我们可能一直踩着一台"有活的假大脑在听 127.0.0.1"的机器自认为是干净的。② 它**会累积**：今天这两枚分别起于 10-08 08:42 与 11:35＝那两天被外力杀掉的测试二进制留下的；每死一次腿就可能多一枚。③ 机制层的缺口是通用的：`t.Cleanup` 只在测试二进制**自己走到收尾**时才跑，而 Windows **没有"父死子随"**。

## 现量（每条都带尺；⛔ 引用前先重跑，枚数与行号一律当快照）

- ★**残留名册两枚**（尺＝`Get-CimInstance Win32_Process`，件 `forensics-2.txt` §R2 ＋ `forensics-3.txt` 逐枚明细）：
  `pid=20904` 起于 `2026-10-08 08:42:21`（当时已活 **3420 分钟**）／`pid=25660` 起于 `2026-10-08 11:35:29`（**3247 分钟**）；各自 **父 PID 现已不存在**（`11724`、`20036` 两枚 PPID 都查不到进程）。占用＝工作集 6.6／7.2 MB、private 44.5／45.9 MB、累计 CPU **0.2 s**（＝不是 CPU  hog，是**挂着不退**）。各持一枚仍在 `Listen` 的回环套接字：`127.0.0.1:54887`、`127.0.0.1:58918`。
- ★**归属怎么证的**（这是本票最关键的一把尺，因为它同时**否证**了另一笔账）：`mockllm.exe` 的命令行里带着**它自己所在的 `t.TempDir()` 全路径**和**本仓 golden 目录**——路径首段就是**测试名**：
  `…\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution2689719599\001\mockllm.exe -addr 127.0.0.1:0 -print-addr -golden-dir "D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden"`
  `…\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed4174395780\001\mockllm.exe …（同形）`
  ⇒ 一枚进程属不属于本仓，**⛔ 看枚数、⛔ 看时间相邻，只看它命令行里那个具名路径**。
- ★**spawn 点两枚**（尺＝`grep -n 'print-addr'` 全仓 `.go`，命中恰两枚）：`internal/llm/adaptertest/mockllm.go:64-81`（`exe := filepath.Join(t.TempDir(), "mockllm"+exeSuffix())` → `go build -o exe tools/mockllm` → `exec.Command(exe, …)`），回收在 **同一件 `:98-104`** 的 `t.Cleanup`（`Process.Kill()` ＋ `Process.Wait()` ＋ `spawnRoot.Wait(2s)`）；第二枚同形＝`internal/llm/openaichat/mockllm_integ_test.go:73`。⚠ **两枚 spawn 点必须同批改**，只改一枚＝留下第二个泄漏源（这一条要写进落地腿的硬约束）。
- ★**同一次杀盘上的第二种症状**：`%TEMP%` 下名字以 `Test` 开头的目录 **21 枚**（尺＝`Get-ChildItem -Directory | Where Name -like 'Test*'`）。Go 的 `t.TempDir` 本该在用例结束时删；留着＝那批测试**没走到收尾**。⚠ 件里"oldest/newest 五枚"那两栏的**日期栏落盘为空**＝那半把尺没打中，⛔ 引那两栏的时序。
- ★**处置已经发生（立案时就把这两枚停了）**：`kill-mine.ps1` 具名 `Stop-Process -Id`，守卫三条＝**名字必须是 `mockllm.exe`** 且 **命令行必须含本仓 golden 目录** 且 **父 PID 现量不存在**，任一条不满足就 SKIP 不动。改后现量＝`mockllm.exe` **0 枚**，`msedgewebview2.exe` **仍 24 枚**（⛔ 动过，理由见 `A821`）。件 `kill-mine-1.txt`。
- ⚠**既有相近先例，⛔ 混为一谈**：产码那侧的子进程死法**早就裁过了**＝`Q-57`（乙已批：`exited()` 改问 OS；丙＝reaper 协程要**另一次**批准；落地＝票 156）。本票只管**测试夹具**这一侧，⛔ 顺手动 `wisp slo` 那条路。另：CI 日志里那行 `goroutine outside the D38 roster (leak symptom) goroutine=mockllm-stdout-reader owner=test`（如 `.scratch/ci-logs/run-37158259050.log:1718`）是**同一枚 helper 的 goroutine 面**、⛔ 进程面，两件事⛔ 混一把尺。

## 要建什么（`AC#0` 之前⛔ 任何产码；三形都必带代价，⛔ 只摆两形）

- [ ] **`AC#0` 只读普查（这一格⛔ 任何产码）**：交回三张表——
  ① **全名册**：本仓所有"测试里起子进程"的地方（尺要写清是 `grep -n 'exec.Command'` 打在 `_test.go` ＋ 测试夹具包（`adaptertest/**`、`internal/**/testdata` 之类）还是打在全仓；**枚数必写"整族还是抽样＋射程目录"**）；逐枚带**文件:行＋它有没有回收＋回收挂在 `t.Cleanup` 还是 `defer` 还是根本没挂**。
  ② **三种死法分开实测**（⛔ 读文档当读数）：(a) 用例正常结束／(b) `go test -timeout` 触发 panic／(c) 外力 `TerminateProcess`（leg 被杀、shell 超时、我按 Ctrl+C）。每一形都要量"**那一发的子进程还在不在**"（尺＝具名 `Get-CimInstance` 名册，⛔ "应该会被回收"）。⚠ 前两形大概率是绿的，那也要把读数交回来——**只有 (b)(c) 的形状能让修法成立**（记忆里那条"改前必红的夹具要跑全部能产坏状态的世界"）。
  ③ **三形代价表**：甲＝**Windows Job Object** ＋ `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`（helper 建 job、把子进程塞进去；测试二进制一死，句柄关，OS 连子进程一起端）；乙＝**子进程自己盯父 PID**（`mockllm` 加一枚 `--parent-pid`，父不在就自退）；丙＝**不动**，把这 N 枚登记成具名已知残留，并把 `mockllm.exe` **加进编队起手闸门的名册**（谁在哪个文件里数进程，要具名到行）。＋**必答两问**：(I) 甲形在**托管 runner** 上会不会把 CI 的别的进程一起端掉（job 的射程边界）；(II) 乙形那种"轮询父 PID"要不要撞 D22 的"⛔ 用墙钟时间差实现超时"（`AGENTS.md` §1.2）——**先答这两问再摆给我裁**。判据＝每条带"哪把尺＋射程目录＋blob 还是工作树"，⛔ 裸数。
- [ ] **`AC#1` 落地形（只在 `AC#0` 交出代价、并经编排者裁"甲／乙／丙"之后才开工）**：硬约束五条——① ⛔ 放宽任何断言；② ⛔ 动 `internal/llm/testdata/golden/**` 一字节（那是票 09／票 11 的回放判据面），也⛔ 改 mockllm 的响应语义；③ **两枚 spawn 点同批**（`adaptertest/mockllm.go` ＋ `openaichat/mockllm_integ_test.go`）；④ ⛔ 碰 `Q-57`／票 156 那条产码子进程路，⛔ 碰 `internal/observe` 的 registry 语义（`mockllmRegistry` 那套 D38 名目是票 151 一族钉过的）；⑤ ⛔ 动档／tag／ledger／`ci.yml`（那是票 302／111 的射程）。
- [ ] **`AC#2` 反恒真那一格**：交付含**成对两发**＝(a) 把新加的回收手段撤掉 ⇒ `AC#0` ② 那形 (b)(c) 当场按预期**测到子进程还在**（具名 PID ＋现量名册逐字，⛔ 只报枚数）；(b) `cmp` 逐字节还原 ⇒ 同形测不到残留。⚠ 若采甲形，必须另附一发证明 **job 真把子进程端了**而不是"测试跑完了所以看起来干净"（尺＝在 (c) 形下发一枚 `Get-CimInstance` 名册）。
- [ ] **`AC#3` 门禁与越界**：`sh scripts/d22scan.sh` rc=0；gofumpt／格式名册**两把并排**（写清射程目录＋blob 还是工作树）；`bash scripts/portable-tests.sh --scope=census` 的 totals 行**逐字未变**（现量基线＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`）；`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺）；每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；⛔ 零 push、commit 必带显式 pathspec 且写在 `$( … )` **之外**、只落自家 `probes/304/<腿名>/**`（票 301 `AC#4b` 那三件正向闸门**从本票起对我和对所有腿生效**）。
- [ ] **`AC#4` 起手闸门名册的具名更正**（无论裁哪形都要这一格）：编队派单模板里"跑真机／整包前 `tasklist` 现量"那一句，现在只数 `wisp.exe`／`balldebug.exe`。本格＝**把该数的进程名定下来**（`mockllm.exe` 必进名册；要不要另数 `*.test.exe` 由 `AC#0` ① 的名册决定），并把"数到 ≠ 归我"那句写死：**数到了要先按命令行归属，归属不明的⛔ 停**（这票的起因正是我把别家的 24 枚当成了自己的）。判据＝改后的那行模板文本逐字交回＋一发用它跑出来的现量。
- [ ] **`AC#4b` 终态锚那一侧的同一条尺（★本格由 `303-r1` 的现形触发，⛔ 与 `AC#4` 合并）**：模板里**终态**那句进程计数（`99-final-anchor.md` 那种 `msedgewebview2.exe=N`）同样必带**归属三件**＝谁起的（命令行里 `--webview-exe-name=` 或具名路径）／父 PID 还在不在／本腿认下哪几枚。⚠ 实测事故链：我在 `A821` 更正完"24 枚＝测试遗留孤儿"那句假话之后 **25 分钟**，`303-r1`（`99-final-anchor.md`，18:01:37）就把别家的进程树报成 `msedgewebview2.exe=18` **枚孤儿**（★我这把按**树**归的尺逐字复量＝件 `probes/orch/webview-orphans/forensics-4-trees.txt`，18:18:05：`total=18 trees=3`，逐树 `SearchHost.exe`／`wetype_update.exe`／`clipsync-desktop.exe` **各 `size=6`**、三枚父进程全在、wisp 系 0 ⇒ ⚠ 那枚**数字 18 是对的**、**"孤儿"那个标签是错的**，而错的标签才是会生成"停法"欠账的那一半）。⚠ 同件还记了另一形：62 秒前同一台面上是 **28 枚含 2 棵 `wisp.test.exe` 树**＝编排者自己那发整包在飞 ⇒ **在飞期间"本机枚数"里必然混着腿自己的窗**，⛔ 拿总枚数当残留数。⇒ 判据＝腿交回的名册里每一枚非零计数都能指到一棵树，指不到就写"归属不明"，⛔ 裸枚数进终态锚、⛔ 由裸枚数生成"停法"类欠账。

## 边界（⛔ 塞进本票）

- ⛔ **`msedgewebview2.exe` 那 24 枚**：它们属于 `SearchHost.exe`／`wetype_update.exe`／`clipsync-desktop.exe`／`magpie-windows-amd64.exe` **四棵别人的进程树**（`--webview-exe-name=` 名册逐枚具名，四枚父进程**全部还活着**，`wisp.exe`／`balldebug.exe` 各 0）。"窗退了而渲染进程没退"那一族在本机**当前零实例** ⇒ ⛔ 为本票立那件事、⛔ 动那 24 枚。详见 `A821`。
- ⛔ **`Q-57`／票 156 的产码子进程路**（`wisp slo` 受测主体的 `exited()`），判据已裁乙、丙要另一次批准。
- ⛔ **托管 runner 上的同类残留**：本票只管**本机可证的夹具行为**；CI 那侧要不要搬档属于票 302／票 111 的射程。
- ⛔ **`mockllm-stdout-reader` 那枚 goroutine 泄漏告警**（D38 roster 那一面）——同一件 helper、**另一层**症状，要修就另立。

## 规矩（本票全程）

- 子代理**只 commit、不 push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 在仓内建 worktree；临时件**只建不删**；证据件⛔ 叫 `.out`（根 `.gitignore` 第 8 行是全仓 `*.out`）。
- ⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件；⛔ 为变绿放宽任何断言。
- **杀进程只停本腿自己起的、且必须先做归属取证**（本票的定式：命令行里的具名路径 ＞ 枚数 ＞ 时间相邻）。CPU／MEM ≥70% 停加派。
- 裁决者≠实现者（`SPEC-12 §4.3`）：本票 `AC#1`/`AC#2` 的判语⛔ 由落地腿自勾。
- ⚠ 三形里**丙＝不做**是合法选项，⛔ 把它当失败；它的代价是具名的（机主台面上会继续累积挂着不睡的假大脑，而我们的起手闸门一直看不见它），选它就把那句写进台账并把 `AC#4` 那格做实。

next=`303-r1`／`300-a2`（在飞）交完 → `304-a1`（`AC#0` 只读，⛔ Go 编译面按补位腿规矩自报）→ 编排者裁形 → `304-r1` → `304-v1`
