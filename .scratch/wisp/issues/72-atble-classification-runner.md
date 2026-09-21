# 72 — A 表（受保护路径）在 CI runner 上退化成 B 表：C26 的纵深防御真破了一格

**Status:** review（修复与用例已入库 `f1033e1` + 插单回执 commit；AC#1/#2/#3/#5 本机可证已勾，**AC#4 的 runner 半边与 AC#6 的改名都不在代理手里**）
**Claimed by:** implementer（票 72 代理，2026-09-21 10:05 接手；根因已由票 70 的 run 35549859581 诊断输出提供）
**Last update:** 2026-09-21 11:00（implementer：执行 10:08 插单，放行侧改窄 + round-3 变异）
**Blocked by:** —（与票 70 共享同一条测试，但**修的是实现不是 CI**；票 70 只负责"这条 CI 变绿"）
**Parallel slots:** ≤1 sub-agent（碰 `internal/risk/`，那里是**冻结契约**，见"硬约束"）
**Spec refs:** C26 PathResolver、SPEC-06 §4 解析管线、D22（契约变更需人批）、票 18 的红队四连

## 症状（不是占位，不是断言写松）

票面/A27 原本的假设是"`test-windows` 那个步骤名自带 placeholder ⇒ 可能是设计上先红的占位"。
**这个假设被 `a04d3e2` 实测否证**：那条命令跑的是 `internal/risk/pathresolver_junction_windows_test.go:88`
的**票 18 真红队用例**（同文件 Case 1..10 覆盖 junction / 8.3 / UNC / `\?\` / 大小写混排 /
豁免只按精确路径 / A 表不可覆盖 / B 表默认拒 / 幂等 / 展开），**一条都不是占位**；
步骤名只是 18/20 落地后没改的**陈旧命名**。

真实失败（GitHub Actions run **35547905707**，commit `5c8f9e4`，`windows-latest`）：

```
pathresolver_junction_windows_test.go:104: canonical "C:\Users\runneradmin\AppData\Local\Temp\
  TestPathResolverJunctionWindows4168261687\001\.ssh\id_testkey" classified B, want ClassA (defense in depth)
```

**同一条命令本机 windows 跑是 PASS**（`:88` 的 reparse 门与豁免门都过，只有**分类**不同结果）。

## 为什么它是安全洞而不是"环境差异导致的假红"

A 表（`blacklist.go:132` 的 `~/.ssh/**` 之类）语义是**不可放行的禁区**；B 表是
"确认可放行"。分类从 A 掉到 B ⇒ **一次 L2 确认就能读到本该永不可达的路径**。
所以判据不是"让 CI 绿"，是"**A 表在任何 home 落点下都必须赢**"。

## 已经被排除的（不要重新排除一遍）

死前那位代理逐条查过、并写进 `a04d3e2` 的 message：
- `os.UserHomeDir()` 在 windows 就是 `os.Getenv("USERPROFILE")`（读的是 `$(go env GOROOT)/src/os/file.go:605`，
  **不是** known-folder），`t.Setenv` 一定生效；
- `normPath` 两侧都小写、`isUnder` 是 `dir+\` 前缀匹配；
- A 表那条 `~/.ssh/**` 在 `home == temp 根` 时**必然**命中；B 表命中的是 `id_*`（`:155`）。

⇒ 症状**精确等价于**："Resolve 交回的 canonical 与 USERPROFILE 锚点**不在同一棵子树下**"。
它留下的未决问题（也就是它的死因）：**canonicalizer 是 Windows-by-design 会产出这种锚点，还是这里有一个真 bug**——
最可疑的两条路径是 ①runner 上 temp/home 走了 **8.3 短名或大小写混排**（`RUNNER~1`）而锚点没有，
②canonical 化把 `\\?\` 前缀或卷大小写规范化成与 `USERPROFILE` 不同的形状。

## 诊断增量已经有了，别重做

`a04d3e2` **一字未改断言**，只把失败信息改成能一次跑定位的：多打三行 `USERPROFILE` / `HOME` /
`userHomeDir()` 实值、A 表实际去比的前缀锚点、以及 `isUnder(...)` 的布尔结果。
本机复跑仍 PASS、gofumpt 空。⇒ **下一次真 run 会直接把 runner 上的锚点打印出来**，先读那段再动手。

## Acceptance criteria

- [x] **AC#1 根因定位**：读新 run 的诊断输出（**不要**在本地猜完就改码），给出"A 表锚点为何不比 B 表更近"的
      **机制级**解释，并指明具体那一行代码。允许结论是"Windows-by-design 需要归一化两侧"，但必须给出证据。
      **证据 = 票 70 从 run 35549859581 抠出的实值**（不是本地猜）：锚点 `c:\users\runner~1\...\001\.ssh` `isUnder=false`，
      canonical 是 `C:\Users\runneradmin\...`。那一句就是漏点：`internal/risk/blacklist.go:120` 起
      `normDir(userHomeDir())` 把 `os.Getenv("USERPROFILE")` 的**字面拼写**直接当锚点，从没进过 C26 管线；
      `internal/risk/blacklist.go:132`（现 `:170`）的 `isUnder(p, home+`\.ssh`)` 于是拿"已展开"比"未展开"。
      结论：**Windows-by-design 会产两种拼写，归一化两侧是契约要求**（SPEC-06 §4 `展开 8.3 短名` / C26 句柄真实路径）。
- [x] **AC#2 修复方向：归一化两侧，不是加特例**。修法必须让"A 表优先"成为**不变式**
      （同一条路径既命中 A 又命中 B 时永远判 A），**禁止**用"把 `.ssh` 的特例写进 canonicalizer"或
      "调整 case 顺序"来过关。
      落地：`internal/risk/pathresolver.go` 新增 `pathForms`/`formsOf`/`anchorForms`（锚点走**同一条**句柄管线，
      `GetFinalPathNameByHandle` + `\\?\` 剥离 + UNC 归一，零字符串启发式）；`blacklist.go` 的 A 规则改成
      `aEq`/`aUnder` 多形比较 + `uncertainAnchorMiss` 的 **fail-closed**（拼写无法证明已展开 ⇒ 判 A）。
      没有 `.ssh` 特例、没有 case 顺序调整、没有任何 `isCI`/环境变量分支（A19/M-7/C-3 族）。
- [x] **AC#3 双向变异检验**：①把修复退回旧实现 ⇒ 分类用例必须转红；②**构造一条 A 与 B 同时命中的形状**
      （比如把禁目录做成 temp 下的子目录、文件名同时匹配 B 表模式），证明修复后的不变式判 A。
      ①= `pathresolver.go` 的 `formsOf` 退回"只给 raw 拼写" ⇒ `TestClassifyAnchorSpellingIsNotVerdict` /
      `TestClassifySpellingInvariance` / `TestClassifyFailClosedWhenSpellingUnprovable` 三条转红，其余票 18 用例**仍绿**
      （本机 USERPROFILE 拼写恰好等于句柄拼写 ⇒ 只有新用例看得见这个洞，正是 AC#4 需要 runner 的原因）；
      退回后 `git diff --quiet -- internal/risk` rc=0。②= `TestAListWinsWhereBothTablesHit`。
- [ ] **AC#4 本机 + runner 两侧都过**：贴出两侧**逐字同命令**的输出与真实 exit code；
      runner 侧必须来自一次**真 run**，不是本地等价环境。
      **本机侧已量（`go test ./internal/risk/... -count=2` rc=0；`-count=1 -v` = 123 `=== RUN` / 122 PASS / 1 SKIP / 0 FAIL，
      那 1 条 SKIP 是 `TestSyncRegistryProbeLive`，HEAD 纯净树同样 SKIP，非本票引入）**。
      **框不勾**：runner 侧要 push 才有一次真 run，而本票明写"绝不 push" ⇒ 交编排者跑。
- [x] **AC#5 不放宽任何断言**：`ClassA` 的期望值、A/B 表内容、豁免语义（只按精确路径）
      一律不许改；若证据表明**必须**改契约，**停下来上报编排者**（D22）。
      **取证仪器一字未动**，判据是逐字比对而不是"我觉得没改"：
      `git diff --quiet a04d3e2 HEAD -- internal/risk/pathresolver_junction_windows_test.go` **rc=0**
      （和票 18 那版逐字节相同），没加 `t.Skip`、没加 build tag、锚点拼写没"对齐"、期望值仍是 `ClassA`。
      A/B 表内容与规则字符串逐字保留（只把"拿字面锚点比"换成"拿展开锚点比"）；
      `ErrReparseDenied` 与 `reparse_point_exceptions` 的"只按精确路径"语义零改动。
      **两侧规则现在刻意不对称（按编排者插单 10:08 改，理由写在 `overrideApplies` 的注释里）**：
      判定（deny）侧遍历全部比较形（漏一次就是把 A 降成 B，方向是"多问"）；
      **批准（override / 放行）侧只认唯一已解析形**——caller 给的拼写必须**就是**句柄真实拼写，
      且 key 等于它，才放行；带第二种拼写进来的（8.3 短名可能别名到**另一个目录**、junction、`\\?\`、UNC）一律再点一次确认。
      钉这条不对称的用例 = `TestOverrideOnlyAcceptsTheResolvedForm`（四条断言，含对照组防"什么都不放"）；
      它的变异"把放行侧改回遍历全部形式"→ 断言 (1)(2) 转红，报的正是
      `Allow:true ... B-list single-file override`。
      我没有走第 4 条反驳通道：`TestCanonicalInputGainsNoSecondForm` 只证明"**Resolve 交回的**路径恒一种形"，
      而 `Gate`/`Classify` 是导出 API，caller 可以递进任何拼写 ⇒ "恒为 1"不成立，窄规则是对的。
      插单前我在 AC#5 写的那段"豁免查表改成全部拼写形"**作废**，那正是被指出来的放宽。
      **改动的爆炸半径也被钉住了**：`TestCanonicalInputGainsNoSecondForm` 断言
      "已经是句柄拼写的输入（= Resolve 交给每个生产 caller 的东西）只产一种比较形且 `certain()`"，
      也就是"对没有第二种拼写的路径，本票一个字都不改变判定"这句话现在是可跑的，不是口头承诺。
- [ ] **AC#6 顺带把票 70 的 AC#3 闭掉**：本票的修复就是那条 CI 步骤转绿的因；
      修好后在票 70 的 log 里补一行指回本票，并把"步骤名陈旧"这件事改成非误导性命名（改名不算放宽门）。
      前半done：票 70 log 已补指回本票的一行。**后半没做，因为 `.github/workflows/ci.yml` 是票 71 在途的文件**（共树禁改清单）；
      建议的新步骤名（逐字可用）：`PathResolver red-team cases (junction/8.3/UNC/\\?\ - tickets 18/20/72)`，
      交编排者或票 71 落。

## 硬约束

- `internal/risk/assessor.go`、`internal/risk/pathresolver*.go`（**除测试外**）、`rules_gateway.go`、
  `docs/PLAN.md`、`docs/specs/**` 都是**冻结契约**。本票大概率要动 `pathresolver*.go` 的实现体——
  **动之前先把改法报给编排者判**，不要先改了再说。测试文件（`*_test.go`）可以直接动，但**不许改断言来迁就实现**。
- 共树共索引：`git add` 只加显式路径；commit 前必须 `git diff --cached --name-only` 确认只有你的文件；
  **禁止 `git commit --amend` / `git reset` / `git rebase` / `git stash` / `git checkout .` / `git add -A`**
  （registry A31/A34：这两条今天各真出过一次事故）。**绝不 push**。
- commit message 用 `git commit -F -` + **带引号的 heredoc**（`<<'EOF'`）；中文长段落不要用 `printf`。
- **每个 commit 同步票面**（勾框 + 一行 log + `next=`）：票 70 那位做了 5 个 commit、174 次调用，
  票面 Progress log 里**一条都没写**——工作全在、账目全空，下一个代理会把这些活重做一遍。
- 别碰 `cmd/wisp/providers*`、`internal/tools/bridge_*_test.go`、`internal/ball/tokens*`、`tools/d22scan/**`
  （四个在途代理的包）。

## log
- [2026-09-21T10:10Z] agent=agent-ticket72 did=**认领 + 根因定位（AC#1 的证据链已闭合，先记 checkpoint）**。
  起点证据是票 70 从 run **35549859581**（commit `6c1b5e9`，job `test-windows`）里抠出的 `a04d3e2` 诊断实值：
  `USERPROFILE="C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestPathResolverJunctionWindows2581009681\\001"`
  / `A-tier anchor tried="c:\\users\\runner~1\\...\\001\\.ssh" isUnder=false`
  / `canonical "C:\\Users\\runneradmin\\...\\.ssh\\id_testkey" classified B`。
  **机制级结论**：`Resolve` 的产物两侧不对称——候选路径过了 `resolveHandle`（`GetFinalPathNameByHandle`，已展开 8.3），
  而 A 锚点走 `internal/risk/blacklist.go:120` 的 `normDir(userHomeDir())`，也就是 `os.Getenv("USERPROFILE")` 的
  **字面拼写**，从没进过 C26 管线 ⇒ `blacklist.go:132` 的 `isUnder(p, home+`\.ssh`)` 拿长名去比短名，静默 miss，
  掉到 B 表 `id_*`（`blacklist.go:155`）。修法按 R17 不变式做：**两侧都进句柄真实路径**，锚点无法证明已展开时 fail-closed。
  next=落 `pathresolver*.go` 的锚点展开实现 + `blacklist.go` 改为多形比较（不动表内容），然后跑 AC#3 双向变异。
- [2026-09-21T10:25Z] agent=agent-ticket72 did=**AC#1/AC#2/AC#3/AC#5 落地（AC#2 的实现 + AC#3 的双向变异都跑过）**。
  实现：`pathresolver.go` 加 `pathForms`/`formsOf`/`anchorForms`/`tailExistsBelow`/`pathExists`（锚点走同一条 C26 句柄管线；
  祖先走不动时把"存在但打不开"的那一节标成 `partial`）+ `blacklist.go` 的 A 规则换成 `aEq`/`aUnder` 多形比较与
  `uncertainAnchorMiss` 的 fail-closed 分支，B 档豁免查表改成"同一文件的所有拼写形"。
  新增仪器文件：`internal/risk/pathresolver_anchor_spelling_windows_test.go` 四条
  （`TestClassifyAnchorSpellingIsNotVerdict` 用 junction 当 USERPROFILE，**不依赖 8.3、不依赖 CI、不需要管理员**就复刻出同一种不对称；
  `TestClassifySpellingInvariance` 把不变式直接写成"同一文件九种拼写同判"；
  `TestAListWinsWhereBothTablesHit` 是 AC#3② 的 A∩B；
  `TestClassifyFailClosedWhenSpellingUnprovable` 用**悬空 junction** 造出"存在但打不开"，证明 fail-closed 走的是 A 不是 B）。
  变异（AC#3①）：把 `formsOf` 退回"只交 raw 拼写" ⇒ 上面第 1/2/4 条转红，报的正是
  `canonical ... classified B ... want ClassA` 与三条 exact A 路径直接 `Class:none Allow:true`（比 runner 症状更宽）；
  票 18 原有十条用例在本机**仍全绿** ⇒ 本机看不见这个洞不是"洞不在"，是"这台机器的 USERPROFILE 恰好等于句柄拼写"。
  退回变异后 `git diff --quiet -- internal/risk` rc=0。门：`gofmt -l internal/risk` 空、`go vet ./internal/risk/...` rc=0、
  `go test ./internal/risk/... -count=2` rc=0、`GOOS=linux go vet ./internal/risk/` rc=0（跨平台不炸）。
  next=AC#4 等编排者 push 后的一次真 run（`test-windows`）；AC#6 的步骤改名因 `ci.yml` 属票 71 而未做，建议名已写进票 70 log。
- [2026-09-21T10:40Z] agent=agent-ticket72 did=**AC#3 变异检验在"修复已入库"的状态下重跑一遍（round 2），并补上钉爆炸半径的第五条用例**。
  变异（`formsOf` 退回只交 raw 拼写）后同一命令 `-run 'TestClassify…|TestPathResolver|TestAList|TestBList|TestResolve|TestReparse'`：
  **RUN=15 / PASS=12 / FAIL=3 / SKIP=0**，红的正是 `TestClassifyAnchorSpellingIsNotVerdict`、
  `TestClassifySpellingInvariance`、`TestClassifyFailClosedWhenSpellingUnprovable`；
  撤掉变异后 `git diff --quiet -- internal/risk` **rc=0**、`git status --short -- internal/risk` 空，
  重跑 **RUN=15 / PASS=15 / FAIL=0 / SKIP=0**。
  新增 `TestCanonicalInputGainsNoSecondForm`：Resolve 交回的 canonical 只产**一种**比较形且 `certain()==true`
  ⇒ "对本来只有一种拼写的路径，本票不改变任何判定"这句话现在有用例守着。
  顺手做了消费者侧的礼貌检查（不是我的包，只读数）：`go test ./internal/tools/ -run 'Path|Sensit|Classif|Gate|Approv'` ok。
  next=交回编排者：AC#4 的 runner 半边等一次 push 后的 `test-windows`；AC#6 的改名等 `ci.yml`（票 71）落地，建议名已给。
- [2026-09-21T10:52Z] agent=agent-ticket72 did=**把"本票把 C26 的 1ms 预算跑坏了"这句嫌疑排除掉，并登记一条对别人有用的观测**。
  `-count=1 -v` 全量在编队嘈杂时偶发只红 `TestResolvePerCallBudget`（墙钟预算，`!race` 编译）。同一分钟内做 A/B：
  **本票的树 0.803 ms/op 过**，**HEAD 纯净树（`git archive HEAD` 解到 /tmp/wisp-head72，不含本票任何改动）1.012 ms/op 红**。
  反向样本也有：更早一轮 HEAD 0.515 过 / 本票 0.455 过；再早一轮本票 1.562 红 / HEAD 同条件过。
  ⇒ 判定：**这台共享机器上墙钟预算的抖动，不是本票引入的**（本票 diff 没碰 `Resolve` 的代码路径，
  `resolveHandle` 的调用数一字未改，`Resolve` 也不调 `formsOf`）。
  登记给票 70/71 的"门要能自报"族：墙钟 + 共树编队 ⇒ 这道门自带抖动，红的时候要先做纯净树 A/B 再定罪。
  **我不动它**（改判据=放宽门，不是本票的活）。
  最终门（编队相对安静的一轮）：`gofmt -l internal/risk` 空、`go vet ./internal/risk/...` rc=0、
  `GOOS=linux go vet ./internal/risk/...` rc=0、`go test ./internal/risk/... -count=2` rc=0。
  next=只剩 AC#4 的 runner 半边与 AC#6 的改名，两件都在编排者手里（push / `ci.yml` 归票 71）。
- [2026-09-21T11:00Z] agent=agent-ticket72 did=**执行编排者 10:08 插单：把"放行侧"改窄，并给两侧规则的不对称补上用例与变异**。
  我 10:25 那版 `overrideApplies` 确实是"遍历全部拼写形"，插单指出的放宽成立 ⇒ 已改：
  现在只认 `cand.real`，且要求 `cand.raw == cand.real`（caller 递进来的就必须是句柄真实拼写）；
  `real == ""`（拿不出句柄真值的平台/死卷）退回唯一可用形 = 票 72 之前的行为，不是新增放宽。
  新增 `TestOverrideOnlyAcceptsTheResolvedForm`（junction 造出"一文件两拼写"，四条断言含"规范+规范⇒放"的对照组）。
  **round 3 变异**：把 `overrideApplies` 改回遍历全部形式 ⇒ 该用例断言 (1)(2) 转红，
  打印的就是被禁的那句 `{Allow:true NeedL2:false Class:B Reason:B-list single-file override … id_*}`；
  撤变异后 `git diff --quiet -- internal/risk` 只看本次改动前状态（见下条 commit 后的复核）。
  同一轮补的 `TestCanonicalInputGainsNoSecondForm` 从"反驳通道"降级成"爆炸半径证明"（AC#5 已改写）。
  next=最终门跑完后 commit；之后本票代理侧无未做项，只剩 AC#4 runner 半边与 AC#6 改名在编排者手里。
- 2026-09-21 09:10 编排者建票。**为什么单开一张**：票 70 的代理在 `a04d3e2` 里明确写了
  "这条按真缺陷登记、够格单开一张票"，而它自己撞 turn 上限死了；把一个安全分类失效留在
  "让 CI 变绿"那张票里，很容易被下一个代理用"改断言/加 skip"的最短路径解决掉——那正是本项目最贵的一类错。

## 裁定 R17（2026-09-21 09:52，编排者）：**本票不需要 D22 批准**——契约早就写死了

我建票时在票面写了"⚠ 本票要动 `internal/risk/`＝冻结区 ⇒ 改法先报编排者判"。
**那句话我现在自己更正**，理由是我去读了契约原文而不是继续凭印象设闸：

- `docs/specs/SPEC-06-security-gatekeeping.md:52` 把解析管线明写成
  `→ 展开 8.3 短名 → 规范化 UNC`；
- `docs/PLAN.md:1376`（C26 本体）要求"**取句柄真实路径（Win: `GetFinalPathNameByHandle`）**"；
- `docs/PLAN.md:1796` 更是**直接推翻**"字面比较"的黑名单，原话理由就是
  "**junction/8.3/UNC/`\?\` 可全部绕过**"。

⇒ **8.3 短名必须在解析层被展开，这是既有契约要求，不是新政策。** 所以：
runner 上 `RUNNER~1` 与 `runneradmin` 两种拼写导致 A 表不命中，是**实现缺陷（bug）**，
修它**不触发 D22**。**D22 仍然管住的是**：若有人想把契约文本改成"分类容忍拼写差异"，那才是契约变更，须人批。

**不变式（本票真正的验收对象，写死）**：
> **安全分类的结果不得依赖路径的拼写形式。**
> 任何"A 锚点与候选路径比较"的代码，两侧都必须是**已展开的句柄真实路径**；
> 若某一侧无法证明已展开，**只能 fail-closed（判 A / 拒绝），不允许 fail-open**。

**三种"看起来能变绿"的解法，本票明令禁止**（这才是我设闸的原意）：
1. 把测试里的锚点也换成短名、或比较前对两侧做大小写/短名"对齐"的**局部补丁**——那是**在测试里复刻 bug**；
2. 给 runner 加特例（`if isCI` / 按环境变量放宽）——分类逻辑的输入里出现**调用方可控的选择器**，
   与 registry **A19/M-7/C-3 同一族**（见编排者记忆第 6 条）；
3. 用 `t.Skip` 或 `//go:build` 把它挡出 `test-windows`——**票 70 在 `secret` 上用 build tag 是合法的
   （DPAPI 按 C28 真是 Windows-only），在本票不合法**：路径拼写不是平台 API 限制，是我们的 bug（A40③）。

**优先级说明**：本票高于普通票，因为它让 **A 表（不可放行的敏感路径）在一种真实环境下静默降级成 B 表（一次确认可放行）**。
本机不复现不代表不存在：`GetShortPathNameW` 在启用 8.3 的卷上同样会产短名，
票 20 已经证明"**桥会拒 8.3**"，但那是**工具层第二次解析**在挡（A38②），**判定层本身仍会被拼写骗过**。

## 编排者插单（10:08，看见你的 WIP 之后立刻发，不等验收）

`git diff -- internal/risk` 里 `formsOf(canonical)` + `classifyForms(cand)` 的**方向是对的**：
A 表判定**遍历所有拼写形式**，任一命中即判 A ⇒ 这是从严（记忆里那条"不许 first-match"的同族正解）。

⚠ **但 `overrideApplies()` 不能共用同一个"遍历全部形式"的规则。** 现在它让
"人对某个文件点过一次 B 级确认"这件事，**在该路径的每一种拼写下都成立**——
这是把**放行**侧放宽了，方向与 A 表相反。它不是等价的"同一个文件"：
**8.3 短名可以别名到另一个目录**（`C:\Users\RUNNER~1` 与 `C:\Users\runneradmin`
不保证是同一个 profile——短名由卷上先创建者决定，历史上这正是 8.3 aliasing 攻击面）。
真发生别名时，"我批准过 runneradmin 的那个文件"会顺带放过**另一个人的同名文件**。

**要求（按这条改，改完在票面记一句"为什么两侧规则不同"）**：
1. **判定（deny）侧**：遍历全部比较形式，任一命中 → 从严。保持现状。
2. **批准（override / 放行）侧**：**只认唯一已解析形式**（句柄真实路径那一种，即
   `GetFinalPathNameByHandle` 的结果），其余拼写**不授予**放行。
   理由就是上面那句：**放行的依据必须比拒绝的依据更窄**——宁可让人多点一次确认。
3. 补一条用例钉住这个不对称：同一文件的**非规范拼写**带着已批准的规范拼写 override 进来 ⇒
   **仍然不放行**；反向（override 记在非规范形式上）⇒ 同样不放行。
   这条用例的变异是"把 override 改回遍历全部形式"，**必须转红**。
4. 若你判断我这条错了（例如 C26 已保证 `canonical` 只可能是句柄真实路径，因而 override 侧的
   多形式**永远只有一个元素**），**就用证据说服我**：给出那个保证在哪个 file:line、
   以及一条能证明"多形式集合在 canonical 上恒为 1"的用例。**别静默保留。**

**逐条回执（implementer，10:58，未走反驳通道）**：
1. 保持现状 — `classifyForms` 的 A 侧仍是 `aEq`/`aUnder` 多形遍历（`internal/risk/blacklist.go`）。
2. 已改 — `overrideApplies` 现在只认唯一已解析形：`cand.real != "" && cand.raw == cand.real && bOverrides[cand.real]`；
   `real == ""`（该平台/该卷拿不出句柄真值）时退回唯一可用形，即票 72 之前的行为，不是放宽。
   **两侧规则不同的原因**（按你要求写进票面与代码注释）：判定侧漏一次 = A 静默降 B（放行风险由别人承担），
   批准侧多放一次 = 8.3 别名到另一个目录时放过同名文件；所以**拒绝可以宽，放行必须窄**，宁可让人多点一次。
3. 已补 — `TestOverrideOnlyAcceptsTheResolvedForm` 四条断言：非规范拼写 + 规范 override ⇒ 不放；
   非规范拼写 + 记在非规范形上的 override ⇒ 不放；规范路径 + 非规范 override ⇒ 不放；规范 + 规范 ⇒ 放（对照组）。
   变异"改回遍历全部形式"实测转红：`Allow:true … B-list single-file override`（断言 1 与 2 同时红）。
4. 不采用 — 我确实补了 `TestCanonicalInputGainsNoSecondForm`，但它只覆盖"Resolve 交回的路径恒一种形"；
   `Classify`/`Gate` 是导出 API，caller 能递进任何拼写（实测 `grep -rn "\bGate(" --include=*.go internal/ cmd/`
   的非测试命中只有 `blacklist.go:76` 它自己的定义 ⇒ 今天还没有生产 caller 传 override 表；
   `Classify` 侧的生产 caller 是 `internal/tools/paths.go:131`，它传的是 `Resolve` 的 canonical），所以"恒为 1"不是可依赖的前提。

