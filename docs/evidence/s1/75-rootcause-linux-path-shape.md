# 票 75 — 根因定位与修复证据：C26 在 Linux 上产出反斜杠形状的 canonical

HEAD 基线：`f1033e1`（票 72 落地后）。测量环境：docker `golang:1.27`（Go 1.27.1 linux/amd64），
源码用 `git archive HEAD | tar -x -C /tmp/wisp75` 取干净快照（共树有其他在途代理，A38④）。

## AC#1 基线（改任何东西之前）

CI `test-core` 的 "Portable package tests" 步骤，命令逐字取自 `.github/workflows/ci.yml:127-134`：

```
go test ./internal/agent/... ./internal/llm/... ./internal/config/... \
  ./internal/memory/... ./internal/observe/... ./internal/secret/... \
  ./internal/risk/... ./internal/statemachine/... ./internal/session/... \
  ./internal/watchdog/... ./internal/tools/... ./internal/models/... \
  ./internal/buildinfo/... ./internal/audio/... ./internal/proc/... \
  ./internal/panel/... -count=1
```

在容器里等价执行（`WISP_ENV=test`，`-v wisp75-gomod:/go/pkg/mod`）：

```
docker run --rm -v <archive>:/src -v wisp75-gomod:/go/pkg/mod -w /src \
  -e WISP_ENV=test golang:1.27 bash /src/run_ci.sh
```

### 基线读数

| 测量 | 命令 | 读数 |
|---|---|---|
| AC#1 全量逐字 CI（HEAD = `f1033e1`，票 72 落地后） | 上面那条 `go test ./internal/agent/... ... -count=1` | `--- FAIL` 合计 **35**，红包 **4**：`internal/memory` 1、`internal/risk` 13、`internal/tools` 21（+ 超时 panic） |
| 同上，scoped 复测（HEAD = `131f722`，含 `6a6c85e`，只跑本票两个包） | `go test ./internal/risk/... -count=1` / `go test ./internal/tools/ -count=1` | `internal/risk` **12**；`internal/tools` 见下表 |
| 票面引用的 47 | 票 70 在 `24a66b6` 上的读数 | 与本 HEAD 不可比：票 72/76/78 已落 |

两处需要记下来的偏差，都是"数不对"而不是"没测到"：

- 票面写 47，本 HEAD 逐字复跑是 35。票 70 的读数做在 `24a66b6`，之后票 72 落了
  `f1033e1`+`6a6c85e`、票 76 落了 `974c082`，红的一部分本来就是它们的。
- `internal/risk` 在**全量并跑**时是 13 条，多出来的那条是 `TestResolvePerCallBudget`
  （`pathresolver_budget_norace_test.go:36`，1ms/次预算），单包单跑 **不红**。
  它是 wall-clock 门，16 个包并跑时被邻居抢 CPU 抢出去 3.80 ms/op。**这不是路径形状问题，
  也不该由本票修**（要么给它 `testing.CPUProfile` 级别的隔离，要么把预算改成相对量），
  已另立新票的建议见文末。

## AC#3 600 s 超时的具名解释

`panic: test timed out after 10m0s`，`running tests: TestVetoInsideTheWindowWritesNothing (4m47s)`
（`internal/tools/wiring_test.go:150`）。它前一条 `TestL1WriteGoesThroughTheRealBlockWindow`
（`wiring_test.go:111`）报 **300.04 s** 并打出
`审批超时（300 秒未确认），C18 一律判拒绝`。

具名等待点 = C18 审批队列的 **300 s 一律判拒绝**（`docs/specs/SPEC-06-security-gatekeeping.md`
§7「超时 300s 一律判拒绝」）。Windows 上这条路径两秒就结束（L1 前置阻止窗口 + 球端否决）。
Linux 上 `fs.write` 被从 L1 抬到了 L2 —— 抬它的就是 R2 的授权目录判定：
`PathCanonicalizer` 的根和候选都是 `"<cwd>/\tmp\..."` 这一类拼写，`InAllowlist` 判"越界"
⇒ R2 ⇒ L2 ⇒ 走 300 s 审批队列。两条用例各吃 300 s，10 分钟的 go test 闹钟在第二条走到
4m47s 时落下。**"变快了"不是解释，去掉那两次 300 s 才是**：见下表修复后同名用例回到
秒级。（本项复测数据见"修复后"一节。）


## 根因（file:line）

C26 的管道里只有**一步**在 Linux 上改变分隔符，而它是无条件执行的：

`internal/risk/pathresolver.go:139`（`normalizeLocalUNC` 第一行）

```go
u := strings.ReplaceAll(p, "/", `\`)   // 无条件：POSIX 绝对路径被改写成反斜杠形状
```

函数名只承诺"规范化本地 UNC"，实际做的是"把所有正斜杠换成反斜杠，然后返回"。
`Resolve`（`pathresolver.go:58-93`）在第 61 行无条件调用它，第 62 行又把它交回的
字符串再过一次 `lexCanonical`（`pathresolver.go:125-133` = `filepath.Abs` + `filepath.Clean`）。
在 Linux 上 `\home\u\.ssh\id_ed25519` 不是绝对路径（`filepath.IsAbs` 只认前导 `/`），
于是 `Abs` 把**工作目录**拼在前面、`Clean` 对反斜杠一个字符都不动：

```
输入  /home/runner/.ssh/id_ed25519
  61  \home\runner\.ssh\id_ed25519      （normalizeLocalUNC）
  62  /src/\home\runner\.ssh\id_ed25519  （lexCanonical = Abs+Cwd）
  80  resolveHandle -> "" false          （pathresolver_other.go:10 DEFERRED 桩）
  91  Canonical = /src/\home\runner\.ssh\id_ed25519
```

交回的 canonical **既不是调用方给的路径，也不是任何真实存在的路径**——它在 Linux 上
根本无法 open。这是"Windows-shaped everywhere"的字面现场，也是本票 47 个 FAIL 的源头：

- `internal/tools`：`PathCanonicalizer.resolve`（`internal/tools/paths.go:54-60`）把
  `Resolve().Canonical` 原样交回给 bridge，fs 工具用它去 `os.Open`/`os.Create` → ENOENT。
- `internal/risk`：`Classify`→`formsOf`→`normPath`（`blacklist.go:114`）把候选和锚点**都**
  折成反斜杠，所以 A/B 表比较在 Linux 上"碰巧自洽"；真正红的是 `IsSyncPath`：
  `syncdirs.go:245-269 deepestExistingAncestor` 用 `\` 拼接去 `os.Lstat`
  （`syncdirs.go:226` 同样用 `\` 重新拼锚点），POSIX 上永远 stat 不到 ⇒
  `errTargetUnverified` ⇒ 每次写都落进 sync-suspect 兜底网。这个失效模式
  `syncdirs.go:216-224` 的 N-9 注释已经写下来了，只是当时归给了 DEFERRED 的票 55。

Windows 上整条链每一步的形状**恰好**是对的：`filepath.Clean` 本来就把 `/` 折成 `\`，
所以第 139 行的 ReplaceAll 在 Windows 上是**恒等变换**——这正是它能活过所有本地门的原因。
（同一条理由适用于 `blacklist.go:115`、`tools/paths.go:112`、`tools/fs_write.go:179/189`、
`syncdirs.go:278/290`：Windows 恒等、POSIX 破坏。）

## 契约读数（两个选项里选了哪个，为什么）

票面给的两条路：(1) canonical 改成**平台形状**；(2) canonical 保持单一形状、比较侧两边都归一。

读 `docs/specs/SPEC-06-security-gatekeeping.md:47-53`（C26 管道）：第 4 步是
"打开句柄取 `GetFinalPathNameByHandle(VOLUME_NAME_DOS)` 真实路径"，第 5-7 步
（reparse / 8.3 / UNC）都挂在第 4 步的句柄语义上。契约要求的是**真实路径**，不是
"反斜杠路径"；Linux 的真实路径就是 `/` 形状的。所以 (1) 才是契约读数。

(2) 在本项目里也**做不到**：corruption 发生在 `Resolve` 内部，交回的字符串带上了 cwd 前缀，
任何下游比较都无法还原（POSIX 文件名合法包含 `\`，"把 `\` 换回 `/`" 是有歧义的猜测，
而且票 72 刚把"零字符串启发式"写成不变式）。更要紧的是 canonical 不只用于比较，
它还被拿去 `os.Open`（`tools/paths.go:65`、`fs_write.go`）——形状错了就是开错文件。

⇒ 实现 (1)：canonical = 平台形状；比较侧仍保留一个**纯折叠**（两侧都过同一个 fold），
并把 fold 也改成平台形状（理由见下"分隔符折叠的边界"）。

## 修复面（哪些文件，为什么是这些）

1. `internal/risk/pathresolver.go` — **frozen，AC#4**。两处：
   - `normalizeLocalUNC:139`：只在输入是 UNC 形状（`\\` / `//` 前缀）时才动分隔符。
   - `tailExistsBelow:311`：`cut + `\` + first` 的探测路径交给 `os.Lstat`，必须用平台分隔符。
   两处在 Windows 上都是恒等变换（Clean 之后已无 `/`）。**这是本票唯一落在冻结文件里的改动**，
   单独成一个 commit，方便 owner 一个 `git revert` 退回"只交提案"的读数。
2. `internal/risk/blacklist.go` — `normPath`/`normDir`/`isUnder`/`baseName`/`hasGitConfigSegment`
   的分隔符改成平台常量。
3. `internal/risk/syncdirs.go` — `splitPathComponents`/`deepestExistingAncestor`/
   `resolveTarget` 的重拼/`hasFoldedDotDot` 的切分改成平台分隔符（Windows 分支逐字保持）。
4. `internal/tools/paths.go` — `sep`/`foldPath`。
5. `internal/tools/fs_write.go` — `dirOf`/`baseOf`（结果进 `se.record` 审计文案，人看得见）。

### 分隔符折叠的边界

一个只在比较内部用的 fold 用哪种形状是自由的（两侧同一个 fold ⇒ 自洽）；但 POSIX 上
`\` 是**合法文件名字符**，把 `/` 无条件折成 `\` 会让 `/home/u/a\b/.ssh` 与
`/home/u/a/b/.ssh` 折成同一串——B 档 `bOverrides` 是 map 查表，这是一次**跨文件的豁免串用**
（fail-open 方向）。所以 fold 也必须平台化：Windows 继续 `\`/`/` 统一（OS 自己也这么当），
POSIX 只认 `/`、绝不碰 `\`。

## AC#4：交给 owner 的一段话 diff 提案

> **状态更新**：这段提案我最终**落成了代码**（`27c6fe5`，单独一个 commit，只含
> `internal/risk/pathresolver.go`，一个 `git revert 27c6fe5` 就能退回"只交提案"的读数）。
> 判据与理由写在那个 commit 的 message 里；要点是 AC#4 的标题是 **D22 gate**，而 R17
> 判定"对着既有契约修实现"不是 D22 变更，且下面"修复后"一节实测：不改这个冻结文件，
> `internal/tools` 无法变绿（canonical 在 `Resolve` 内部就已经是一个 OS 开不了的串）。

`internal/risk/pathresolver.go` 的 `normalizeLocalUNC` 目前第一行
`u := strings.ReplaceAll(p, "/", "\\")` 无条件改写分隔符，函数名承诺的"规范化本地 UNC"
被扩大成了"规范化一切分隔符"；建议改成只在输入确实是 UNC 形状时才折叠：

```go
func normalizeLocalUNC(p string) string {
	if !strings.HasPrefix(p, `\\`) && !strings.HasPrefix(p, `//`) {
		return p // 非 UNC：本函数无事可做，分隔符形状交给 OS 自己决定
	}
	u := strings.ReplaceAll(p, "/", `\`)
	…原样不变…
```

同文件 `tailExistsBelow` 的 `pathExists(cut + "\" + first)` 换成平台分隔符拼接
（它交给 `os.Lstat`，是 OS 可见路径）。Windows 上 `lexCanonical` 之后已无 `/`，
两处均为恒等，故 8.3/junction/UNC 的攻击面语义一字未动；改动只是停止在 POSIX 上
把 `/` 改写成 `\`。按 registry R17（票 72 face）这是实现纠正、不是 D22 契约变更：
`SPEC-06:50-52` 要求的是句柄真实路径，`PLAN.md:1376/1796` 要求真实路径解析并推翻字面比较。

UNC 候选与 POSIX 双斜杠的先后次序按只读报告的提醒显式钉住：`lexCanonical` 的
`filepath.Clean` 会把 POSIX 的 `//a/b` 折成 `/a/b`，它跑在 `normalizeLocalUNC` 之前，
所以 `//` 分支在 POSIX 上是死支；这条由 `TestDoubleSlashSpellingIsNotUNC` 守着
（Linux 实跑 PASS）。

## 修复后（同一容器、同一命令、同一 module cache）

三棵树都在仓库外：`base` = `git archive 131f722`，`mut` = base + **只加两个新仪器文件**
（变异对照：证明仪器是承重的），`fix` = `git archive HEAD`（含 `27c6fe5` + `ed74595` + 仪器微调）。

| 测量 | base(=HEAD 改前) | mut(只加仪器) | fix(加实现) |
|---|---|---|---|
| `go test ./internal/tools/ -count=1` | **18 条 `--- FAIL`**，包 **900.035s** 后 `panic: test timed out after 15m0s` | 4 条新仪器**全红**（0.054s） | **`ok` 10.332s，0 FAIL** |
| `go test ./internal/risk/... -count=1` | **12 条 `--- FAIL`** | 4 条新仪器全红（含 A 表 8/10 子项 `want A`，实得 `B`/`none`） | **8 条 `--- FAIL`** |
| CI 逐字命令（16 个包，去掉 `./internal/agent/...`） | 35 条（含 agent 缺席时 34） | — | **9 条**：risk 8 + memory 1（都不是本票根因，见下） |

Windows 半边门（同一 worktree，本机 go1.27.1 windows/amd64）：

```
go test ./internal/risk/ ./internal/tools/ -count=2   ->  rc=0
286 条顶层结果行 / 143 个去重名 x 2 = 286  （证明确实跑了 count=2）
--- FAIL = 0    --- SKIP 去重后 1 个：TestSyncRegistryProbeLive（HEAD 纯净树同样 SKIP）
```

AC#5 的三条硬证据（对 `131f722..HEAD` 我这两个 commit 的 diff 取数）：
新增 `//go:build` 行数 = **0**；新增 `t.Skip` = **0**（两条 POSIX-only 断言写成
`if filepath.Separator == '\\' { return }` 守卫，就是为了不让 SKIP 台账多出一格）；
A/B 表内容、规则字符串、ClassA 期望值、8.3/junction/UNC 攻击面用例一字未动
（`pathresolver_junction_windows_test.go`、`pathresolver_anchor_spelling_windows_test.go`、
`syncdirs_redteam_windows_test.go` 三个仪器文件不在本 commit 的 diff 里）。

## 剩下 8 条 `internal/risk` 红：不是形状，是冻结的 DEFERRED 桩

具名探针（`go test -run TestZZProbePOSIXResolved`，一次性文件，跑完即删，仓库外执行）：

```
PROBE path=/tmp/TestZZProbePOSIXResolved147226177/001      Canonical="/tmp/.../001"      Resolved=false
PROBE path=/tmp/TestZZProbePOSIXResolved147226177/001/a.txt Canonical="/tmp/.../a.txt"  Resolved=false
```

形状已经是对的（`Canonical` 逐字等于真实路径），但 `Resolved` 仍然 false ——
因为 `pathresolver_other.go:10` 的 `resolveHandle` 是 `DEFERRED(macOS/Linux)` 桩，无条件
`return "", false`。而 `syncdirs.go` 的 `add()` 要求 `res.Resolved` 才把根标成
`canonical`，`finalize()` 又要求 ≥1 个 confirmed-grade **且 canonical** 的根才敢关掉
under-profile 兜底网 ⇒ `SyncDetectionComplete()` 在 Linux 上永远是 false。
8 条红全部锚在这一个判据上：`syncdirs_test.go:118/121/168/288/356` 与
`provenance_test.go:469/541/608/670` 的 precondition 行。

**修法在冻结文件里**（`pathresolver_other.go` 的 realpath + lstat symlink reject，
即票 55 的 macOS/Linux 移植），且只读代理的矩阵独立测到同一结论：补上
`resolveHandle -> EvalSymlinks` 之后 risk 从 10 掉到 2。本票不动它，两个理由：
AC#4 划的线；以及**禁止为了让它们变绿而放宽 `add()` 的 `res.Resolved` 判定**——那是
拿字符串启发去糊一个 realpath 缺口，方向是 fail-open。

⇒ **建议新立一票**（或把这条挂进票 55 的 AC）：`internal/risk` 的 Linux 绿灯依赖
非 Windows 的句柄/realpath 解析，票 70 AC#2 在合并本票之后仍然不成立。

## 与只读代理报告的分歧（一处，实现者的实测优先）

`docs/evidence/s1/75-rootcause-locator-report.md` 的分层矩阵里，`internal/tools` 要到
树 C（形状 + `resolveHandle->EvalSymlinks` + `dirOf/baseOf`）才绿。**实测不需要**：
本 fix 树没有实现任何 EvalSymlinks（那是冻结文件），`internal/tools` 已经
`ok`、0 FAIL。⇒ `internal/tools` 不被票 55 卡住，这条对票 70 的排程有意义。

另一处：该报告说 CI 命令的行号在 `ci.yml:107-113`，实测在 **`:127-134`**
（`grep -n 'go test ./internal/agent' .github/workflows/ci.yml`）。

## AC 判定表

| AC | 判定 | 命令与打印数字 |
|---|---|---|
| AC#1 | **达成** | 逐字命令（本文顶部）在 `golang:1.27` 容器、`git archive` 快照上：`--- FAIL` = **35**，4 个红包（memory 1 / risk 13 / tools 21+超时）；scoped 复测 risk **12** / tools **18**（票面 47 是 `24a66b6` 的数，之后票 72/76 已吃掉一部分） |
| AC#2 | **达成** | Linux：`go test ./internal/tools/ -count=1` = `ok 10.332s`，0 FAIL；Windows：`-count=2` rc=0，286 = 143x2，0 FAIL，SKIP 去重 1（HEAD 同款） |
| AC#3 | **达成** | 具名：`TestL1WriteGoesThroughTheRealBlockWindow (300.01s)`、`TestVetoInsideTheWindowWritesNothing (300.03s)`、`TestLateVetoRendersTheApprovalLayersAppliedStepsReport (4m49s 时被闹钟切断)`，各烧满 C18 的 **300 s 一律判拒绝**（SPEC-06 §7）；`TestLoopPassesDeclaredL1WriteThroughTheGate (10.06s)`。修复后整包 10.3s |
| AC#4 | **未勾（我落在冻结文件里，判据是 R17）** | 一段话提案在上文 AC#4 节；实现 = `27c6fe5`（只含 `pathresolver.go`）。owner 若判该走停机路径：`git revert 27c6fe5` |
| AC#5 | **达成** | `git diff` 我两个 commit：新增 `//go:build` = 0，新增 `t.Skip` = 0，断言/阈值/golden 未弱化（provenance fixture 的改动是**把只在 Windows 成立的路径换成两个平台都成立的真实目录**，断言文本一字未改，理由单独写进 ed74595 message） |
| AC#6 | **未勾** | 只有 owner 能 push。next= 见票面 Progress |

