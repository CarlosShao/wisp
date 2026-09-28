# 174-c2（只读代价普查·零产码）——"把路径判定者接到生产"这一跳要动哪几层、会把谁顶出去

- 程：`174-c2`｜性质：**只读**（AC 框一枚没勾、产码零字节不动、不落地、不选方案）
- 派单：`.scratch/wisp/dispatches/2026-09-28-145x-readonly-174-c2-wiring-cost-census.md`
- 工单：`.scratch/wisp/issues/174-task-output-points-at-a-file-the-model-cannot-read-by-default-spill-artifacts-are-outside-fs-allowed-dirs.md`
- 骨架落盘时刻：`2026-09-28 14:42+0800`（现跑 `date '+%Y-%m-%d %H:%M%z'`）

> 这份表**只报代价与形状，不替编排者决定接不接、什么时候接**。
> 凡是本程没跑的，一律明写"本程未验"，不拿推理充数。

---

## 1. 起手锚＋写面闸门

- 分支 `dev`。派单给的锚点是 `03c01d71`；**本程起手现量 HEAD＝`e9ef94d0`**（共享工作树，别人在我之前又落了一枚 `docs(evidence): 181-c1 …`）⇒ 起手锚以现量为准，差的那枚是别人已提交的证据件，与本程无关。
- 起手尺（现跑）：`grep -n "TaskDeps{" cmd/wisp/run.go` → **`365: for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks}) {`**（票面上写的 `:362`／`:609` 都已漂号；产物目录那行现在在 `:612`／`:640`）。
- 起手 `git status --porcelain -- internal/ cmd/` → **空**（这一行是闸门自证，非推理）。
- 终态 `git status --porcelain -- internal/ cmd/` → **空**（跑完三把门禁之后再取一次，仍空；见 §8 末行）。
- 写面（只这些，逐枚）：新建 `docs/evidence/s1/174-wiring-cost-census-c2.md`（本件）、新建 `.scratch/wisp/probes/174/c2/{zz174c2_wiring_pair_windows_test.go, overlay.json, logs/*.txt}`、追加票 174 一段 Progress log。**AC 框：0 枚勾。**
- `docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden／审批超时常量：**零字节**（不在我的 pathspec 里）。`frontend/**`／`design/**`：**未读未写未引用**。

## 2. Q1 接线那一跳的完整形状（逐环：谁构造·谁持有·谁调用）

| 环 | 位置（现量行号） | 内容 |
|---|---|---|
| 1 构造 | `cmd/wisp/run.go:327-331` | `allowed` 逐条拷 `cfg.FS.AllowedDirs` ⇒ `rt.paths = tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)`。**判定者今天已经存在，且只造一枚** |
| 2 持有者 A | `cmd/wisp/run.go:345-348` | `tools.BuiltinFSEntries(tools.FSDeps{Paths: rt.paths, …})` ⇒ fs.* 已接 |
| 3 持有者 B | `cmd/wisp/run.go:451-453` | `rt.bridge = tools.New(tools.Options{ … Paths: rt.paths …})` ⇒ 桥（C19 判级）已接，**与 2 同一枚实例** |
| 4 持有者 C＝缺口 | `cmd/wisp/run.go:364-365` | `rt.tasks = tools.NewTaskRoster()`；`tools.TaskDeps{Roster: rt.tasks}` **没有 `Paths`** ⇒ 本票要补的就是这一处赋值 |
| 5 装配进工具 | `internal/tools/task.go:295-297` | `BuiltinTaskEntries(d)` → `taskOutput{d: d}`，`d` 就是那枚 `TaskDeps` |
| 6 调用点 | `internal/tools/task.go:242-243` | `if rec.ArtifactPath != "" { notice := t.d.pointerNotice(rec.ArtifactPath) }` |
| 7 真用它的那处 | `internal/tools/task.go:318-330` | `pointerNotice`：`:319` `d.Paths == nil` → 现在真机上走的就是这支（通用"未接线"那句）；`:323` `d.Paths.Canonicalize(raw)`；`:326` 规范化失败那一支；`:327` `d.Paths.InAllowlist(canon)`；`:331-334` 只读 `os.Stat` 那一腿（接线前后都已在跑） |
| 8 路径从哪来（上游） | `internal/tools/task_backfill.go:104-115` | `b.Spills.Prepare(taskArtifactPrefix+taskID, text)` → `case sp.Path != "": rec.ArtifactPath = sp.Path` → `b.Roster.Record(taskID, rec)`；`Spills` 在 `cmd/wisp/run.go:640` 造的 `agent.NewSpiller(filepath.Join(rt.spec.dataDir,"artifacts"), …)` |

**要动几枚文件**：
- **最少 1 枚**：`cmd/wisp/run.go:365` 给 `TaskDeps` 加 `Paths: rt.paths`（复用现成实例，不新造判定者、不加根）。
- **要真满足 AC#2b 得 2 枚**：那一支的文案今天只说一条回来的路（`task.go:328-329`「要用户先把所属目录加进 `[fs]` allowed_dirs」），而编排者 09-27 22:3x 第 3.2 条要求改后的字必须说出**两条**（加根，**或批一张 L2 卡**）。只改 `run.go` 那一行 ⇒ 精确文案上线、但**仍是半句话**（本程现量，见 §3 表里那两句逐字）。⇒ 第二枚＝`internal/tools/task.go`，并**必须**同时补/改判据。

**每动一处哪几枚既有判据会红（逐枚点名）**：
- 只动 `cmd/wisp/run.go:365` ⇒ **既有 Go 判据零枚会红**。理由与名册（现跑）：`internal/tools/**` 里所有断言都自己构造 `TaskDeps`，没有一个测试读 `cmd/wisp` 的组合根——`grep -rn "TaskDeps" --include=*_test.go internal/` 只命中 `internal/tools/pointer_183_cli_seam_test.go:111`（`TaskDeps{Roster: roster, Paths: paths}`，**它今天就已经带着判定者**，加 `run.go` 一行不影响它）。`cmd/wisp` 名下没有任何测试文件（现量：`glob cmd/wisp/*_test.go` 为空；本程未跑 `go test ./cmd/wisp/`，174-c1 现量该包在本机报 `0xc0000135` 缺 DLL）。
  - 会变的**不是红，是消音**：`.scratch/wisp/probes/176/r1/zz176r1_e2e_test.go:350-353` 那一格今天是 `t.Log` 而不是断言（原文"已登记的形状：…TaskDeps.Paths 在 assembleRuntime 里没接"）⇒ 接线后这条日志不再出现，**它不会红，但它作为"缺口在场"的记录失效**，要重算名册。
  - `internal/tools/task_output_pointer_notice_test.go:232` 的注释写着"what cmd/wisp/run.go:362 builds today" ⇒ 号已漂（今天 `:365`），**注释不是判据，不红**，但接线那枚程该顺手把它的行号说准。
- 再动 `internal/tools/task.go:328-329`（补第二条回来的路）⇒ **会红的名册**（按断言行逐枚点名，**本程未跑变异**，行号来自现跑 `grep -n`）：
  1. `TestPointerOutsideAuthorizedRootSpeaks`（函体起 `:105`，断言在 **`:125`** `strings.Contains(out.Text,"不在你被授权的目录范围内")`）
  2. `TestPointerToMissingCopyFileSpeaks`（起 `:146`）的**反向钉 `:169`**（不许出现那句授权文案）——若把两条路写成含混的一句，可能反被 `:163` 的正钉放过、却被 `:169` 判红
  3. `TestBothPointerDefectsAreReportedSeparately`（起 `:267`，断言 **`:283`／`:286`**）
  4. `TestUnwiredJudgeFailsClosedInReply`（起 `:234`，**`:253`** 正向钉"未接线/按读不到处理"、**`:256`** 反向钉）——这枚是**接线判据的镜像**：接线只动 `cmd/` 时它仍绿；一旦 `pointerNotice` 的 nil 分支文案被顺手改，它红
  5. `TestHealthyPointerStaysSilent`（起 `:203`，**`:218`** 六枚禁字表含 `"未接线"`）——它管"健康回执不许多嘴"，加句不动它
  6. `TestPointerNoticeKeepsTheD153StubShape`（起 `:299`）与 `TestHealthyReplyStillMatchesThePreFixTemplate`（起 `:349`）：后者是**逐字节等于修码前模板**的钉（174-v1 已现算 md5 基线 `8bd433707c1fd9039ac2103da08e7d22`）⇒ 任何往**健康**那一支加字的行为它必红；只改"根外"那一支不动它。
- 门禁面：`scripts/d22scan.sh` 的 ban #7（`internal/tools/`）与 ban #8（emoji/字符串）扫的是形状，加 `Paths:` 赋值与改文案都不引入 `filepath.Clean|Abs`；**本程只现跑了"未动码"那一发的读数（§8），动码后的读数归落地腿自取**。

## 3. Q2 顶出去效应（接上前／后两发对照）

台件：`.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go` ＋ `overlay.json`（`-overlay` 编进 `internal/tools`，**跟踪件零字节**，`git status --porcelain -- internal/ cmd/` 跑完仍空）。读数全在 `.scratch/wisp/probes/174/c2/logs/probe.txt`（`rc=0`）。
两臂只差一处，正是本票那一跳：`before = TaskDeps{Roster}`／`after = TaskDeps{Roster, Paths: judge}`；**fs.\* 与桥在两臂里都用同一枚 `judge`**，所以"回执说的话"和"`fs.read` 实际给的判定"出自同一个 C26 实例。

判定者语义（现读，非推理）：`internal/tools/paths.go:50-52`「An EMPTY allowlist authorizes nothing: every fs call lands at L2 (R2)」；`:133-163 InAllowlist` 还带两道额外的收窄——`:152 resolvedForm`（跟不了 link＝不授权）与 `:159 p.workspace`（票 92 AC#3 的工作区收窄，**只会更紧**）。

| 臂 | 形状 | `fs.read` 接上前 | `fs.read` 接上后 | 回执 接上前（逐字） | 回执 接上后（逐字） |
|---|---|---|---|---|---|
| 默认 `allowed_dirs=[]` | 真 spill | `isError=true level=L2 class="user_rejected" bytes=57` | **同一发** | `注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理）`（`注意：`×1） | `注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来`（×1） |
| 同上 | 名存实不存在 | 同上 L2/user_rejected | 同上 | 未接线那句（×1） | 授权那句 **＋** `注意：这条路径现在读不到，宿主登记的那份副本文件并不存在`（×2） |
| 同上 | 是枚目录 | 同上 | 同上 | 未接线那句（×1） | 授权那句 **＋** `注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状）`（×2） |
| 同上 | 根外真实文件 | 同上 | 同上 | 未接线那句（×1） | 授权那句（×1） |
| `allowed_dirs=[artifacts]` | **真 spill（健康）** | `isError=false level=L0 bytes=30000`（全文读回） | **完全同一发** | 未接线那句（**×1＝明明读得到却说读不到**） | `注意：` **×0（安静）**，指针 `全文见` 仍在 |
| 同上 | 名存实不存在 | `isError=true level=L0 class="tool" bytes=180` | 同上 | 未接线那句 | `…宿主登记的那份副本文件并不存在` |
| 同上 | 是枚目录 | `isError=true level=L0 class="tool" bytes=132` | 同上 | 未接线那句 | `…它存在但不是一般文件（是目录或别的形状）` |
| 同上 | 根外真实文件 | `isError=true level=L2 class="user_rejected" bytes=57` | 同上 | 未接线那句 | 授权那句 |

**收益（逐字对照里那枚唯一的相变）**：`root-is-artifacts-dir` ＋ 健康 spill 一形——接上前回执说"按读不到处理"而**同一枚桥的 `fs.read` 给出 `isError=false level=L0 bytes=30000`**（假阴性，模型会放弃一条真能走通的路）；接上后 `注意：` 归零、指针仍在。这就是 AC#2b 要的那一跳的全部真实收益。

**退让／顶出去（必须一起报的那一支）**：
1. **物理层面：一枚都没有。** 八形 × 两臂，`fs.read` 的 `isError/level/class/bytes` **逐项相同** ⇒ 接线**不动授权**，只动文字。今天能读的接上后仍可读、今天被拒的接上后仍被拒。**没有"被拒的会变放行"**，因为那一行不加根（`allowed_dirs` 默认值一字未动，`NewPathCanonicalizer(nil,…)` 仍 authorize nothing）。
2. **说法层面新造出一枚"说得太肯定"的退让**：接上后那句是"它不在你被授权的目录范围内，**fs.read 会被拒**"，而 174-v1 已实测接一枚答"允许"的门后同一条桥对根外路径 `isError=false level=L2 bytes=4096` 全文读回（`FSDeps.open` 执行腿只 `Canonicalize`、不查根，`fs.go:93-105`）。⇒ 接上后那句话在**没人批**时为真、在**有人批**时为假，而它只给了"加 `allowed_dirs`"这一条出路。本程在同一处补现量：默认臂那一形 `NoGate{}` 下 `isError=true level=L2 class="user_rejected"`——**接上前它也说"读不到"（按读不到处理），所以这不是接线新造的谎言，而是接线把一句含糊的悲观话换成一句精确但只说半条路的话**。⇒ **代价：只接一行 = AC#2b 的第二支（两条回来的路）没兑现。**
3. **未造的形（本程未验，落地腿必须自己补）**：`paths.go:159` 的 workspace 收窄那一支（spill 在根内但在被收窄的工作区外）与 `:152 resolvedForm` 跟不了 link 那一支（跨盘／junction）。两支都会在接上后**新**给出"不在你被授权的目录范围内"，而接上前只说"无法核实"——**这两支是本表唯一没能现跑的潜在顶出去**（要造 junction，且 `AGENTS §1.2` 与票 107 的形状不在本程写面）。

## 4. Q3 与票 183／185 那根管子的关系

**结论：接上 `Paths` 之后，`task.output` 那条产物路径**不会**自动读得回。**（现量＋名册两头都指同一件事）

- 现量（本程台件）：能不能读回**只由授权根与门决定，与 `TaskDeps.Paths` 无关**——同一形在两臂里 `fs.read` 读数逐项相同（§3 表）。默认配置（`allowed_dirs=[]`）下八发全 `L2 / user_rejected`；把 `artifacts` 配进根后健康那一形 `L0 / 30000 字节`。⇒ 接线只让回执**说实话**，不让管子**通**。
- 名册：票 183 文件名带 `-done`（`183-the-host-minted-pointer-exemption-does-not-hold-on-the-real-cli-…-done.md`），它的教训正是"包级全绿≠端到端通"，而且它**已经**在自己的判据里把判定者接进 `TaskDeps`（`internal/tools/pointer_183_cli_seam_test.go:111`）——**那根管子的读回腿另有 C25/R4 盖戳一段，不在本票这一跳上**。票 185 未 `-done`（`185-every-successful-reread-mints-the-blocker-for-the-next-one-…md`）：即使读回一次成功，`fs.read` 的结果又变成一枚未声明的 C25 mark，给下一发造阻塞。⇒ **174 的授权腿通了也还是通不过 185 那一腿。**
- **本程未验（明写）**：① 端到端 `wisp run` 一发没跑（本机 `go test ./cmd/wisp/` 报 `0xc0000135` 缺 DLL，来路＝174-c1 现量；本台件跑在 `-overlay` 的 `go test -run TestProbe174C2…` 里，**接缝级**）。② 本程台件用 `NoGate{}` ＋ `risk.ProvOptions{NoProbe:true}`，**故意绕开门与盖戳**，所以它一个字也测不到 183／185／177 那条 taint 管子的形状。

## 5. AC#2d 那一支今天到底有没有判据钉着

**没有。现量三条尺（不拿"文案已改"充当）：**

1. `grep -rn "连规范化都没通过" --include=*_test.go internal/` → **0 命中**（全仓 `internal/` 的测试里，没有一处出现过 `pointerNotice` 那一支的逐字文案）。
2. `grep -c '规范化' internal/tools/task_output_pointer_notice_test.go` → **`0`**（编排者 22:3x 第 3.1 条那把尺，今天复算仍是 0）。
3. 对照：另外四支各有钉——`不在你被授权的目录范围内`＝`:125`／`:169`(反)／`:256`(反)／`:283`；`并不存在`＝`:163`／`:256`(反)／`:286`；`不是一般文件`＝`:191`；`未接线`＝`:218`(禁字表)／`:253`。`internal/tools/task_output_pointer_notice_test.go` 的八枚函体名（现量）：`:105 Outside…`／`:146 ToMissing…`／`:178 ToNonRegular…`／`:203 Healthy…`／`:234 Unwired…`／`:267 BothDefects…`／`:299 KeepsTheD153StubShape`／`:349 HealthyReplyStillMatchesThePreFixTemplate`——**没有一枚问 `Canonicalize` 返错那一支**。
- 唯一沾"无法规范化"的既有判据都在**另一条腿**上（R2 判级／卡面，不是 `pointerNotice`）：`internal/tools/bridge_junction_windows_test.go:232`／`:306`／`:446-447`／`:459-460`、`internal/tools/fs_edit_ac5_gate_r4_test.go:358`、以及 `rules_gateway.go` 那条 R2 fail-closed 腿。⇒ **"根里有别的测试"绝不等于"这一支有钉"**。
- 本程**没跑变异**（跑变异要动产码，本程禁），所以"摘掉它会全绿"这一发按 174-v1 的现量登记（`MUT-V5` → `rc=0`、顶层 124 全绿），标为**复算自非本程**。
- ⚠ 顺带一条形状提醒给落地腿：那一支是**五支里唯一把 `err.Error()` 原文拼进模型可见文本**的开口（`task.go:326`），补判据时必须同时问"说了没有"＋"拼进去的这段不许带路径原文／根列表／C26 内部状态"。本程现算：`Canonicalize` 的错文里有 `tools: empty path`（`paths.go:120-121`）这种宿主内部串。

## 6. AC#3 三条禁区自证（只读能答的部分）

| 禁区 | 本程没碰的证据 |
|---|---|
| (i) 宿主内部 artifacts 写入＝不许做成受门控的 Tool | 本程**没注册任何工具、没改任何注册表**。台件只调用现成的 `agent.NewSpiller(…).Prepare` 与 `task.output`／`fs.read`（都按今天树上的形状）。`grep -n "BuiltinFSEntries\|BuiltinTaskEntries" probes/174/c2/*.go` 里两枚都是既有构造函数。相关旁证登记在案：`docs/evidence/s1/176-background-start-port-census-a1.md:206` 逐字引 `spill.go:29-31`「HOST-INTERNAL artifact write … deliberately not a gated tool」 |
| (ii) 不许在 `risk.PathResolver`/C26 之外用 `filepath.Clean\|Abs` 做文件系统决策 | 台件里 `filepath` 只用于**造 fixture**（`Join`/`TempDir`/`ToSlash`，把路径喂给被测件）与日志打印；**判定与决策一律走 `NewPathCanonicalizer`→`Canonicalize`→`InAllowlist`**（`zz174c2…go` 无 `Clean`、无 `Abs`：`grep -c "filepath.Clean\|filepath.Abs" probes/174/c2/zz174c2_wiring_pair_windows_test.go`＝0，见 §9 那行）。`internal/tools/paths.go:14-18` 那段"D22/CI-banned patterns appear NOWHERE in this file"一字未动 |
| (iii) 不许把 `artifacts` 目录塞进 `allowed_dirs` 默认值 | 默认臂用的就是 `NewPathCanonicalizer(nil, nil)`；`root-is-artifacts-dir` 那一臂是**测试内构造的对照根**（`t.TempDir()` 下的临时目录），**没有落进任何配置默认值**，`SPEC-03:35` 那行与 `cfg.FS.AllowedDirs` 零字节。⚠ 本程也**没有**把 `allowed_dirs` 当执行时硬边界用：§3 第 2 条退让就是照"`InAllowlist` 只是判级输入、执行腿 `fs.go:93-105` 只 `Canonicalize`"写的，把根做成硬边界属 `Q-60` 的另一支，须 owner 说话 |
| 逐枚数字自证 | `git status --porcelain -- internal/ cmd/` 起手与终态各取一次，**皆空**；`git show --numstat` 删除列对本程每一枚 commit **逐枚 0**（三枚全是新增文件＋票面纯追加）；`git diff --numstat` 对 `internal/ cmd/ docs/PLAN.md docs/specs/**` 为空 |

## 7. 本程没测什么（逐名）

1. 端到端 `wisp run` 一发没跑（本机缺 DLL `0xc0000135`，来路 174-c1；本程未重跑 `go build ./cmd/wisp` 之外的真机腿）。
2. `go test ./cmd/wisp/` 没跑（该包名下现量无 `_test.go`，且受同一 DLL 事实影响）。
3. **没跑变异** ⇒ §2 里"改文案会红哪几枚"是**逐枚点名但未现跑**；只有"只动 `run.go` 一行则零枚红"是现量（构造名册＋glob）。
4. `paths.go:159` workspace 收窄 × spill 路径那一形（§3 退让第 3 条）未造。
5. `paths.go:152 resolvedForm` 跟不了 link／junction 那一形未造。
6. `AC#2d` 那一支的 `err.Error()` 具体会拼出什么串进上下文，未造真发（只现读了 `tools: empty path` 这一枚内部串）。
7. R4／C25 盖戳那一腿**完全未测**（台件 `NoGate{}` ＋ `NoProbe:true`，故意绕开）。
8. `internal/agent/spill.go` 那枚 D15(3) 桩（AC#2c）没看没测，且派单写明等票 177 裁完。
9. 审批超时／`ErrorClass`（`bridge.go:397-404` 那族，票 173 名下）未碰。
10. `probes/176/r1/zz176r1_e2e_test.go` 没重跑，只在源码里读到 `:350-353` 是 `t.Log` 而非断言。
11. `gate-clauses.sh` 的 `G6neg` 那记红（基线 1／实测 3）**没追来路**——另有程正在跑名册差集，本程只报"在册名册未变"。
12. 真机 `dataDir` 里 `config.toml`／`secrets\` 的暴露面（票 174 "为什么值得做"里那支最坏后果）未测，本程所有 fixture 都在 `t.TempDir()`。

## 8. 门禁

现量时刻：本程最后一枚产码无关改动**之后**、本件与台件落盘之前（此后我只新增文档与追加票面一段，`internal/`／`cmd/` 字节不变 ⇒ 读数对最终树成立）。

- `sh scripts/d22scan.sh` → **`rc=0`**；`d22scan: scope ban #8 internal/ examined 433 Go files, comments and _test.go included`（**＝派单基线 433，未漂移**）；`ban #8 cmd/=45`、`ban #7 internal/tools/=21`、`ban #6 frontend/=85`、`ban #8 design/=39`；末行 `clean - no D22 ban violations`。日志 `probes/174/c2/logs/gate-d22scan.txt`。
- `bash .scratch/wisp/probes/154/gate-clauses.sh` → **`rc=1`**；`腿数＝14 声明与实测不符＝1`；红腿名册**只有 `G6neg`**：`BAD 腿=G6neg 声明=ring 基线=1枚 实测=3枚 因=新增未成对（票 171 AC#2：实测 > 基线）`——**＝派单说的"在册只有 G6neg"，本程未加腿**；`基线过期枚数＝0`；`腿数断言＝相符（名册=14 声明=14 记账=14 缺腿=0 空头声明=0）`。**那把尺一字未改**。日志 `logs/gate-clauses.txt`。
- `go test -count=1 ./internal/tools/`（单跑，**不带 `-overlay`**）→ **`rc=0`**、`ok github.com/CarlosShao/wisp/internal/tools 15.812s`。⇒ 本程序台件**不在**这发读数里（它只经 `-overlay` 虚拟进包，跟踪件零字节）。日志 `logs/gate-gotest.txt`。
- 台件那一发：`go test -count=1 -overlay=… -run TestProbe174C2WiringBeforeAfterPair -v ./internal/tools/` → **`rc=0`**（它只打印读数、不断言，所以 rc 只证明编译通过）。日志 `logs/probe.txt`。
- 跑完三把门禁后再取：`git status --porcelain -- internal/ cmd/` → **空**。

## 9. 被拒调用＋零删除自证＋工具调用终值

- 被拒工具调用：**无**（本程每一次调用都被放行，没有需要改道重试的拒绝；一次编译失败是**我自己的台件写错类型**（`Result`→`agent.ToolOutcome`），非权限拒绝，也未触碰跟踪件）。
- 跑过的删除命令：**无**。本程命令全集：`date`／`git log|status|add|commit|show|diff`／`ls`／`mkdir -p`／`grep`／`go test`／`sh scripts/d22scan.sh`／`bash gate-clauses.sh`／`tail`。`rm`／`del`／`Remove-Item`／`clean`／`restore`／`stash`／`--amend`／`rebase`／`worktree`／`switch` **零枚**；`probes/161/r6/flip-declaration.sh` **未跑**（派单禁）。
- 重跑既有台件前的 sink 检查（派单 ⚠ 条）：`grep -n "WriteFile\|OpenFile" .scratch/wisp/probes/154/gate-clauses.sh` → **0 命中**；对 `probes/174/r1`／`probes/174/v1` 的 `.sh`／`.ps1` 同样 **0 命中**（那两家的读数不是我造的）。⇒ 本程所有新读数只落进 `probes/174/c2/logs/`（`probe.txt`／`gate-d22scan.txt`／`gate-clauses.txt`／`gate-gotest.txt`／`readings.txt`），**没有就地覆盖型 sink**，别家票逐行引用的读数未被洗掉（`A367` 那一族）。
- ⚠ 一处如实：本程**复跑**了 `scripts/d22scan.sh` 与 `gate-clauses.sh` 两把既有尺（它们只读、不写跟踪件，上面有 sink 证明），落点全在我的 `logs/`。
- 工具调用终值：**33 枚／硬顶 35**（本件最后这枚更正在第 32-33 枚，只动本文字）；**新探索在第 23 枚之后停止**，其后全在落盘·造台件·跑门禁·提交。骨架第一枚 commit 落在第 **6** 枚调用（比派单的"≤5"晚一枚：我在写之前先花一枚现量"这枚新文件是否已存在"，以免覆盖跟踪件——`ls docs/evidence/s1 | grep 174` 现量只有 `c1`／`r1`／`v1` 三枚）。
- ⚠ 一处逐字自证：本程 commit 的 `git show --numstat` 里**唯一非零的删除列＝`111 16 docs/evidence/s1/174-wiring-cost-census-c2.md`**，那 16 行是本程自己在第 5 枚调用新建的**骨架占位行**（每节的"（待填…）"），删除列对**任何一枚起手即存在的跟踪件＝逐枚 0**（票 174 面 `2 0`＝纯追加，其余全是新建件）。

## 10. next：落地腿派之前还缺什么（含要不要人先批准哪一格）

1. **接线那一行本身不必再等人批**：`Q-60` 已批＝**丙**（不改能力、改回执；`e6e3fb33` 逐字，票面 22:0x 更正行已定），补 `Paths: rt.paths` 正落在丙的射程内，且 §3 现量它**不放宽任何授权**（八形 × 两臂 `fs.read` 逐项相同）。**没有新根、没有新判定者实例。**
2. **但只接一行交不掉 AC#2b**：判据要求"两向都要"，而接上后的精确文案**只给一条回来的路**（加 `[fs] allowed_dirs`），缺"批一张 L2 卡"那半句——编排者 22:3x 第 3.2 条已把这条并进 AC#2b。⇒ 落地腿的写面**必须含 `internal/tools/task.go:328-329`**，否则 AC#2b 只是把含糊的悲观换成精确的半句话。
3. **要人先说话的那一格**：改那两句文案会撞 `TestHealthyReplyStillMatchesThePreFixTemplate`（`:349`）与 `TestPointerNoticeKeepsTheD153StubShape`（`:299`）这两枚"模板/形状冻结"钉，而 `task.go:250-252` 的注释把 stub 的措辞系在 **`PLAN.md:2564`／票 164 已接受 AC#3** 上。⇒ **"只加一句 L2 卡的出路"算不算动到那句冻结文字**，是契约轴（本票 AC#4 名下），**须 owner 或编排者先裁**；未裁之前落地腿不许动 `docs/PLAN.md` 一字。
4. **AC#2d 与接线解耦，可单独先派**：它是纯补判据（零产码语义变更、零契约），名册现成：`task.go:326` 那一支全仓零覆盖（§5 三把尺）。判据要求两半都问（说了没有＋`err.Error()` 不许带什么）。**它不需要任何 owner 批准**，是本表里性价比最高的一枚。
5. **缺的宿主**：端到端那发要一台不缺 DLL 的机器（174-c1 已把这条记为结题条件），否则 AC#2b 的两向里"接上之后的真机读数"永远只能标接缝级。
6. **顺序建议（不是方案，只是代价排序）**：AC#2d 补判据 → 接线那一行（`cmd/wisp/run.go:365`，一枚文件、零枚红）→ 文案那半句（等第 3 条那句批准）→ 票 177 裁完才动 AC#2c（`spill.go`）。⚠ **别把这三步并成一枚程**：并了之后，"零枚红"和"哪几枚红"就分不开，正是本仓咬过的那类归因丢失。
7. **本程留给落地腿的三个未造形**：workspace 收窄（`paths.go:159`）、`resolvedForm` 跟不了 link（`:152`）、`Canonicalize` 返错时 `err.Error()` 的实际串（`task.go:326`）。前两支是**唯一潜在的顶出去**。
