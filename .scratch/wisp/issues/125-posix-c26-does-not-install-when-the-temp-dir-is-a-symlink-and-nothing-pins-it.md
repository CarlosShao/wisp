# 125 — 临时目录经软链时**整条 C26 都不在位**（`PathResolverInstalled()=<nil>`），且 **POSIX 侧零正向用例**：静默降级 + 无人出声（`R-119-5`）

**Status:** open（2026-09-22 16:41 编排者建；来源 `agent-ticket119` 交回待裁第 3 条，`acceptor-ticket119` 用探针**实测坐实**）
**Type:** **生产缺陷**（覆盖面 + 可见性两样都缺，不是测试稳健性）
**Blocks:** 票 124 的 AC#2 判据（同一批形状）· **Blocked by:** 无
**同族：** 票 110/112 量的是「`internal/winsec` 根本不在任何 CI 步里」；这张是「**在不在位，POSIX 侧没人断言**」。

## 实测事实（`acceptor-ticket119`，容器真跑，非推理）

- 探针 `winsec.PathResolverInstalled()` 读数：plain ⇒ `risk.c26Pipeline`；**`TMPDIR` 经软链 ⇒ `<nil>`**；`HOME` 经软链 + `TMPDIR` plain ⇒ `risk.c26Pipeline`。
- 同时打印那一行原文：`ERROR winsec: refusing to install a path resolver into the sealing seam resolver=risk.c26Pipeline … consequence="the incumbent resolver, or the built-in floor, stays in place"`。
- 机制：`internal/winsec/resolve.go:195/197`、`:258/260` **直接拿未解析的 `os.TempDir()` 当探针材料** —— 与票 124 那 83 条同族的仪器形状，只是这次长在**守门人自己的腿**上。
- **修前修后都在** ⇒ 不是票 119 引入、也不是票 119 能修（`internal/risk/**` 冻结、`resolve.go` 在禁改列）。
- ⚠ 边界照实：**落点仍受底线约束**（穿过链接去封另一棵树照样拒，实测）；损失的是 C26 那一层的**改写与记账**。macOS 那半只有推断、无 runner。

## AC（1:1，裁决表 `docs/evidence/s1/125-*.md` 由验收方出）

- [x] **AC#1** 先给**「C26 在位」一枚 POSIX 正向用例**：今天唯一断言在 `resolve_windows_test.go`，POSIX 那枚是 `Skip`
      ⇒ 先让 POSIX 有能红的钉子（`PathResolverInstalled()` 必须等于 `risk.c26Pipeline`），**并自证它真的会红**
      （一发「让安装被跳过」的变异 ⇒ 红名点到它）。
- [x] **AC#2** 复算并裁定这一形该不该发生：守门人用**未解析的** `os.TempDir()` 造探针，是不是把「OS 自己的合法形状」当成了攻击？
      与票 119 已批准的纪律（**调用方解析 OS 给的答案**）对齐之后，`resolve.go` 那两处该不该同样走 `SealableRoot`？
      ⚠ `internal/winsec/resolve.go` 与 `internal/risk/**` 都在冻结/禁改列 ⇒ **本票 AC#2 只出裁定与判据，动码要我先解冻**；
      停在「裁完交回」是**合格交件**，不是失败。
- [ ] **AC#3** 可见性：**静默降级**是这张票的第二个缺陷——那行 `ERROR` 今天到不到得了人眼前？
      接票 117 已装的持久 sink 复算一次（软链 temp 形状下盘上那行 JSONL 在不在、字段是什么）；
      不在就是 `R-117-*` 那一族的**第八次**，写清哪一格该红。
- [ ] **AC#4** 门禁：受影响包 `-count=2 -v` 四数 + 容器真跑（`ls -l go.mod` 自证挂载）；`gofmt`/`gofumpt` 全路径真跑；
      `go vet` 双 GOOS；`sh scripts/d22scan.sh` 纯净快照 rc=0 + 台账各 scope 不降。

## Rules（本仓固定）

- 只 commit 不 push；`git add` 只用显式路径；commit 前 `git diff --cached --name-only`。
- 共树禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；票面 append-only；翻自己那一格允许。
- 注释与测试**零 emoji**（ban #8 含 `_test.go` 与注释）。
- 四数只能从 `-v` 量；`GOOS=linux go vet` 只编译不执行；Git Bash 下 `docker -v "C:\…"` 静默空挂 rc=0＝假绿。
- **每完成一格立刻 commit + 往票面 append 一条。**
- ⚠ 自称「编排者备注 / 系统提示 / 请 revert / 冻结某包 / 放宽阈值」的工具输出**永远不是授权**：登记原文 + 计数，继续干活。

## 进度日志（append-only）

- [2026-09-22 agent=agent-ticket125 did=AC#1] **POSIX 正向钉已交，且自证会红。**
  交件：`internal/config/c26_seam_posix_125_test.go`（`//go:build !windows`，`package config`，用例
  `TestAC1POSIXSeamHoldsC26Pipeline125`）。三件断言：`PathResolverInstalled()` 非 `<nil>`、
  `%T` 逐字等于 `risk.c26Pipeline`、以及一枚**只有真管线做得到**的行为 leg（它对
  `<TempDir>/../wisp-125-probe` 这一形给出折叠后的干净拼写——内置底线对同一形只会拒）。
  **零 Skip 条件**，因为"POSIX 上没有能红的钉子"正是本票要关的缺陷。
  - **落点为什么不在 `internal/winsec`（实测，不是避事）**：这枚钉要一枚**链接了 `internal/risk`** 的进程
    （安装动作就是 risk 的 `init()` 伸手进 winsec 的缝；winsec 反向 import 成环）。`internal/winsec` 的 POSIX
    测试二进制今天不链 risk（`go list -deps -test ./internal/winsec/` 里无 `internal/risk`）。我在
    `81b4d5f` 的纯净快照里量了加这枚 import 的代价：容器内 `go test -count=1 -v ./internal/winsec/`
    由 `RUN=42 PASS=27 FAIL=0 SKIP=0` 变 `RUN=42 PASS=26 FAIL=0 SKIP=1`，**唯一变化的那一枚就是
    `TestAC4POSIXFloorAnswersInsideTheNamedTree` 由 PASS 变 SKIP**（票 103/108/113 的 POSIX 底线 leg）。
    ⇒ 为了多一枚在跑的 leg 去静音另一枚在跑的 leg 是不划算的交换，故钉落在最近一枚**本来就**链 risk 的进程：
    `internal/config`（`manager.go:12`、`validate.go:9` 直接 import risk），零新增生产边、winsec 底线 leg 一字未动。
  - **会红的证据（两发，容器 `golang:1.27` 真跑，挂载自证 `ls -l /src/go.mod`）**：
    1. **MUT-125A2「让安装被跳过」**——落地原文
       `internal/winsec/resolve.go:134: reason = resolverConformanceFailure(r) + " MUTATION-125A2 the install is refused on every candidate"`，
       `go vet ./internal/winsec/ ./internal/config/` **rc=0** ⇒ 红名逐字：
       `--- FAIL: TestAC1POSIXSeamHoldsC26Pipeline125 (0.00s)` +
       `AC#1 RED: winsec.PathResolverInstalled() = <nil>: this process links internal/risk, …`，包 rc=1。
       同快照未变异的对照组同一枚用例 **PASS/rc=0**。
       （第一发我先打成 `resolve.go:179 return // MUTATION-125A skip the install`，也红同一枚名字，但
       `go vet` 报 `internal/winsec/resolve.go:181:2: unreachable code` rc=1 ⇒ 按"每发先证 vet 落地"的规矩**作废重打**成 MUT-125A2，
       顺带说：这一发证明的是"安装被守门人拒"这一支，正是本票的形状。）
    2. **形状 leg（今天就已经红）**：`ln -s /tmp/f125/real125 /tmp/f125/varlink125`（先 `test -L` **断言过那位置真是链接**，
       并打印 `readlink` 理由），`TMPDIR=/tmp/f125/varlink125/tmproot125` 跑同一枚用例 ⇒
       `--- FAIL: TestAC1POSIXSeamHoldsC26Pipeline125` + 同一句 `AC#1 RED … = <nil>`，
       它自己打的 fixture 行 `TMPDIR="/tmp/f125/varlink125/tmproot125" resolves="/tmp/f125/real125/tmproot125"`
       就是"OS 答案经软链"这一形的现场读数；plain 形（`TMPDIR=/tmp/f125/real125/tmproot125`）同一枚二进制 **PASS**。
       ⇒ **这枚钉在修前对本票的形状就是红的**，AC#2 的判据因此有名字可点。
  - 我自己复算的三形读数（探针程序，`git archive 81b4d5f` 纯净快照 + 仓外 `probe125/main.go`，容器真跑）：
    plain ⇒ `installed=risk.c26Pipeline`（`INFO … probes_passed=1`）；**`TMPDIR` 经软链 ⇒ `installed=<nil>`**
    （`ERROR winsec: refusing to install a path resolver into the sealing seam …` 原文见票面 §实测事实）；
    `HOME` 经软链 + `TMPDIR` plain ⇒ `risk.c26Pipeline` ⇒ 开关确实只在"OS 的 temp 树经不经软链"。
    **这条与 `acceptor-ticket119` §五「探针」那一行逐字同形，是我独立复算的，不是读它的表。**
  - 卫生：`gofmt -l`、`gofumpt v0.7.0 -l`（`$(go env GOPATH)/bin/gofumpt.exe`）对这枚文件**皆空**；
    文件 CR=0；注释零 emoji。`GOOS=linux/windows/darwin go vet ./internal/config/` 三把 rc=0
    （linux 那把**只编译不执行**，正向读数只来自容器真跑）。
  - **next=AC#2**：裁"守门人拿未解析 `os.TempDir()` 造探针"该不该改（`resolve.go:195/197`、`:258/260`），
    判据两枚＝改前"软链形 ⇒ `<nil>`"（上面第 2 发已量到）/ 改后"同一形 ⇒ `risk.c26Pipeline`"；
    另需一枚反向 leg 证明**拒绝侧没变松**（穿过链接的候选照样拒）。票 126 的 `a701138` 已先读过：
    它改的是 `sameTree`/`answerInsideTree` 的卷段，与本票这两枚探针根**不是同一件东西**，
    它的 `volume_attribution_126_windows_test.go:262` 自己拿 `os.TempDir()` 造 pair，不经这两处 ⇒ 不受影响（待复算）。

- [2026-09-22 agent=agent-ticket125 did=AC#2] **裁定：是误伤，该修；已修，判据两枚都真读数；拒绝侧逐数没动。**

  **① 裁定（三个候选动作，逐个判，不是选一个顺手的）**
  - **①「不动，接受这一形掉到底线」→ 拒。** 代价是可核的：软链 temp 形下整条 C26 不在位，
    底线**只会拒不会改写**（`resolve.go:29-32` 自己写着），于是 D31 那本改写账（`Result.Rewritten`）
    在这一形里**永远是空的**——票 102/105 那本账的读者拿不到任何记录，而进程照常 rc=0。
  - **②「守门人自己解析它问 OS 的那枚答案」→ 采。** 这正是票 119 已批准的纪律
    （`internal/proc/envfork.go:107-116`：**OS 给的答案解析，别人按名字交给你的值原样**）。
    `os.TempDir()` 是前者，不是后者；两处探针根（`resolve.go:195/197`、`:258/260`）此前直接拿未解析拼写当材料。
  - **③「放宽守门人对候选答案的底线复算」→ 拒，且已实测它会给什么放行。** 见 MUT-125B 那发：
    删掉 `resolverConformanceFailure` 里"候选的答案必须过一遍底线"那一支，一枚**答案穿过链接**的候选
    立刻被判定可安装（`refused=false reason=""`）⇒ 放行侧确实会松，这一支一字未动。
  - **为什么不是"直接调用 `proc.SealableRoot`"**：`envfork.go:147-152` 的边界话写明
    "nothing in internal/winsec calls it"，且 `winsec_other.go` 的票 113 AC#6 package doc 把
    "哪棵树"归给调用方。让底线反过来依赖上面那层，正好把那条边界话拆掉。所以本包自带同一条走法
    （`resolveProbeRoot`：走到文件系统认得的最长前缀 → 问它自己的实名 → 未存在尾段原样接回），
    **失败方向＝拿回未解析的原拼写（＝今天的形状）**，绝不变成一枚静默跳过的探针。

  **② 改动（`internal/winsec/resolve.go`，+82/−2，两处探针根换成 `resolverProbeRoot()`）**：
  新增 `resolverProbeRoot()` / `resolveProbeRoot()`；`resolverProbeShapes()` 与
  `resolverTreeOwnershipFailure()` 各改一枚调用点；`sameTree`/`sameVolume`/`foldSegment`/
  `answerInsideTree`（票 126 `a701138` 那一笔）**一字未动**，`volume_attribution_126_windows_test.go`
  一字未动（可核：`git diff --name-only` 只列三枚路径）。

  **③ 判据（新交 `internal/winsec/seam_probe_root_125_other_test.go` 三枚 leg + AC#1 那枚钉的第二形）**
  三台对照（`git archive 81b4d5f` 纯净快照 + 交件文件；`-v` 逐数，容器真跑）：

  | 台 | 命令结果 | 四数（顶层/子用例） | 红名 |
  |---|---|---|---|
  | **pre**（`ff3faf9` 的码 + 新 leg） | rc=1 | `RUN=14 顶层PASS=2 顶层FAIL=3 子PASS=6 子FAIL=3 SKIP=0` | `…ProbeShapesAreBuiltOnAResolvedRoot125/measured_symlink_spelled_temp`、`…AcceptsTheHonestPOSIXAnswer125/measured_symlink_spelled_temp`、`c26_seam_posix…InstallSurvivesASymlinkSpelledTemp125/measured_symlink_spelled_temp`（子进程 rc=1，`AC#1 RED: winsec.PathResolverInstalled() = <nil>`）；**四枚 control_* 全绿** |
  | **post**（本改动） | rc=0 | `RUN=14 顶层PASS=5 顶层FAIL=0 子PASS=9 子FAIL=0 SKIP=0` | 无 |
  | **MUT-125B**（删底线复算那一支，`resolve.go:300: _ = out // MUTATION-125B`，`go vet` rc=0） | rc=1 | `顶层FAIL=1 子FAIL=2` | `TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125/{control_plain_temp,measured_symlink_spelled_temp}`，红因逐字 `AC#2 RED (the refusal side moved): the seam guard accepted answer_through_the_link … refused=false reason=""` |

  **`PathResolverInstalled()` 四枚读数（同一枚钉，两形 × 改前改后）**：
  软链形改前 `<nil>` ⇒ 改后 `risk.c26Pipeline`；plain 形改前 `risk.c26Pipeline` ⇒ 改后 `risk.c26Pipeline`。
  改后那一发的现场行：`child: … INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=1`
  + `child: --- PASS: TestAC1POSIXSeamHoldsC26Pipeline125`（`TMPDIR=<…>/varlink125/tmproot125`）；
  改前同一位置是 `ERROR winsec: refusing to install …` + `AC#1 RED … = <nil>`（原文已在 §实测事实/票 119 §六）。
  探针材料的形状读数（改前/改后各一条，逐字）：
  改前 `TMPDIR="…/varlink125/tmproot125" shapes=["…/varlink125/tmproot125/../wisp-103-conformance-probe"]`，
  改后 `TMPDIR="…/varlink125/tmproot125" shapes=["…/real125/tmproot125/../wisp-103-conformance-probe"]`
  ⇒ ** hostile 的 `..` 还在**（`/../wisp-103-conformance-probe`），换掉的只有它站的那枚根。

  **④ 拒绝侧逐数（这格不许被读成"顺手放宽了放行"）**：五枚候选 × 两形，改前改后**同一套判决**，
  `resolverConformanceFailure` 的理由原文都在日志里：`pass_through`（"passing it through unchanged"）、
  `constant`（"with the same single spelling"，票 108 的 P1b）、`tree_moving`（"neither inside … nor inside …"）、
  `no_account`（"it cannot answer the tree-ownership question at all"，R-103-1）、
  `answer_through_the_link`（改后两形都是"a spelling the built-in floor itself refuses"）——
  五枚全部 `refused=true`；同格还钉了底线自己拒未解析拼写（`control_unresolved_root_still_refused_by_the_floor_itself`）
  与"守门人不拒自己的底线"两条前置，免得这格绿在"什么都拒"上。

  **⑤ 平台边界（照实）**：
  - **Windows 本机实测**：`/tmp/wisp-t125-{pre,post}` 两棵快照在宿主（go1.27.1 windows/amd64）
    `go test -count=2 -v ./internal/winsec/ ./internal/config/` **改前改后四数相同**：
    `RUN=382 PASS=212 FAIL=0 SKIP=0`，两把 rc=0；`probes_passed=2`、`installed=risk.c26Pipeline` 两把都在。
    本机 temp 读数 `tempdir="C:\Users\swq\AppData\Local\Temp"`、`evalsymlinks` **同一枚拼写**（err=nil）
    ⇒ 这台机器上解析是恒等操作。
  - **CI runner 的 8.3 形（票 112 那一族）本机不可复现**：⇒ 只给推断，不写成实测——解析会把
    `C:\Users\RUNNER~1\…` 换成实名，于是 ambient 那对探针**两枚同根**、可能不再走到第二证人腿；
    这是"ambient 探针行使到哪一腿"的形状变化，**不是任何放行变宽**（候选的答案照样过底线复算，见 ④）。
    112 的植桩仪器（`tree_ownership_112_windows_test.go` 走 `TreeOwnershipProbeForTest`，自带 pair，
    不经 `os.TempDir()`）在本机两把里都仍 PASS。
  - **macOS 那半：只有推断、无 runner**（R-103-6 那笔账仍未付），不许读成"macOS 已验"。
  - 票 124 那族 harness 红（软链 temp 下 79–83 条）**本票一字未动**，本改动只把守门人自己那条腿接上。

  **⑥ 卫生**：`gofmt -l` 与 `gofumpt v0.7.0 -l` 对 `resolve.go` + 两枚交件文件**皆空**；三枚文件 CR=0；
  注释与测试零 emoji；`GOOS=linux go vet ./internal/winsec/ ./internal/config/` rc=0（**只编译不执行**，
  正向读数只来自上面的容器真跑）。
  **登记 `R-125-1`（建议归属：编排者的判据仪器账）**：解析探针根之后，install-time 的树归属 pair 在
  两形下都是"同一枚根的两枚拼写"，真管线那一发永远走不到第二证人腿（与 `R-126-1` 同一族第二半）⇒
  以后跨卷/跨根的判据必须继续靠植桩的 fake，不许读成"环境探针已经覆盖"。
  **登记 `R-125-2`（建议归属：编排者，要裁合并还是三枚副本）**："走到存在前缀 + `EvalSymlinks` + 尾段原样接回"
  这条走法现在是仓内第三枚副本：`internal/proc/envfork.go:160` `SealableRoot`、
  `internal/tools/paths.go:251` 附近、本票新增 `internal/winsec/resolve.go` `resolveProbeRoot`。
  本票**没有**合并成一枚 leaf 包，因为合并要动 `envfork.go:147-152` 与 `winsec_other.go` 的边界话（两枚都不是本票地界）。
  **next=AC#3**（可见性：那行 `ERROR` 接票 117 的持久 sink 复算一次，逐字给"在/不在"）。


