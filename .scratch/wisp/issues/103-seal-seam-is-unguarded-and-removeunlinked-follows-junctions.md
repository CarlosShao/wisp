# 103 — 密封的**两处旁路**：`SetPathResolver` 谁都能装（装了个橡皮图章就静默重写外来 DACL），`RemoveUnlinked` 能沿 junction 删别人真文件且返回 nil（票 94 验收的 R-c / R-b）

**Status:** ready-for-review（2026-09-21 18:4x `agent-ticket103` 修完两处守卫并跑完三门；正文仍为 append-only，numstat 里那 1 行删除就是本行）
**Type:** 安全边界（一个是**接缝无守卫**，一个是**今天够不到的陷阱**——两者不同档，同票不同判据）
**Blocks:** 票 94 挂 `-done` 的条件之一 · **Blocked by:** nothing
**Packages:** `internal/winsec/`（seam 的注册口与 `RemoveUnlinked`）。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、
              `internal/risk/pathresolver*.go`（**入口展开那条归票 102，别在这儿顺手改**）、`rules_gateway.go`、
              `tools/d22scan/**`、`allowlist.txt`。

## 两条读数的原文位置

`docs/evidence/s1/94-adversarial-acceptance.md` 的 **PROBE A**（`SetPathResolver` 无守卫）与 **PROBE F / PROBE W**
（`RemoveUnlinked` 沿 junction 删目标 + `removeStray` 的 `WalkDir` **不下降** ⇒ 今天不可达）。

- **R-c（今天可利用）**：`SetPathResolver` 没有任何守卫 ⇒ 包外可以装一个"什么都说 OK"的解析器，
  之后 winsec **静默重写外来主体的 DACL** 并报告成功。
  ⚠ 这是票 94 为了绕开 `winsec → risk` **传递依赖成环**（`risk → observe → secret → winsec`，`go list -deps` 实测）
  才引入的装配形状 ⇒ **形状本身是必要的，缺的是守卫**。别把 seam 拆掉倒回 `filepath.Abs`（那是把票 94 白做）。
- **R-b（今天够不到，但是陷阱）**：`RemoveUnlinked` 收到一个 junction 时可以删掉**别人树里的真文件**并返回 nil。
  验收代理自己实测了"为什么今天够不到"（`removeStray` 用 `WalkDir` 且**不下降进链接**，PROBE W）
  ⇒ 所以本票的判据是**把它变成"一旦有人接上就会红"的 tripwire**，而不是"现在就有人踩了"。
  ⚠ 修它要动 `internal/memory`（票 18/79 地界）⇒ **AC#2 明写不许顺手扩界**。

## AC（1:1，裁决表 `docs/evidence/s1/103-*.md` 由验收方出）

- [ ] **AC#1（R-c）** seam 只能被**装一次**、且装的必须是**真解析器**：给注册口加守卫
      （幂等/一次性 + 类型上不给伪造留门，或伪造时**响亮失败并审计**）。
      判据用例两条腿：**装第二个 ⇒ 红**；**装一个恒说 OK 的 ⇒ winsec 必须拒**，
      且"什么都没发生"不算绿（要能证明**外来 DACL 没被改**，取 SID 级读数）。
- [ ] **AC#2（R-b）** 给 `RemoveUnlinked` 一条**平台用例**：输入一个 junction/符号链接 ⇒ **必须拒、返回错误且不删目标**；
      并在 `internal/memory/removeStray` 一侧补一条**"下降进链接就会红"的 tripwire**（断言遍历不下降）。
      ⚠ 修 `RemoveUnlinked` 的**语义**若需要动 `internal/memory` ⇒ **停手登记交回**，本票只做守卫与 tripwire。
- [ ] **AC#3** 变异三向：① 去掉注册守卫 ⇒ AC#1 红；② 把"伪造解析器"改成"注册了但被拒" ⇒ **必须仍有一条用例绿**
      （证明它红在守卫而不是红在噪声）；③ 把 tripwire 的"不下降"改成"下降" ⇒ 该用例红。
      锚点=承载行为那一行，**同链 grep 证落地**，`go build` rc=0 先量到（**编译失败不算变异**），
      还原在 `/tmp` 快照里做并证 `diff -q` 干净。
- [ ] **AC#4** 回归：`go test -count=2 ./internal/winsec/ ./internal/memory/` rc=0 且逐条点名 SKIP/FAIL
      （⚠ 票 89 的 `TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs` 刚落地，别把它改松）；
      `gofmt`/`gofumpt` 空、`go vet` 与 `GOOS=linux go vet` **按包** rc=0、
      **收尾前必跑 `sh scripts/d22scan.sh`** 纯净树 rc=0（A64②）。
      ⚠ `go test ./cmd/wisp/` 本机需按**票 98** 的注入命令跑；`TestResolvePerCallBudget` 负载下会假红（票 86）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；
**不 push**；不在仓内建 worktree（A38④，快照带会话后缀）；票面 append-only（改行前先读；标题前插段落要重抄标题，删除列必须 0）；
四种假绿逐条点名；数字不达标写 FAIL 附数字；**真机测 junction 只在临时目录造、测完清掉，绝不删真数据**；
15 次工具调用内交回第一枚 checkpoint；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/winsec/` 刚被 `agent-ticket89b` 改过、它的**复验正在跑** ⇒ 你与它可能同文件；
开工前先 `git log --oneline -5 -- internal/winsec/` 看有没有未结的验收，并在票面登记"我等谁让路"。

## Progress log（append-only）

- 2026-09-21 17:4x（编排者）：建票。**两条我合成一张但判据分开**，理由：它们是同一个包、同一个"密封动作"的两处旁路，
  分两张会互相等（都要碰 `internal/winsec/`）；但**档位不同**我写在标题里了——
  **R-c 今天可利用**、**R-b 今天是陷阱**（验收代理自己实测了"够不到"）⇒
  AC#1 是修，AC#2 只是守卫 + tripwire，**不许把两者混成一个"都已修"的读数**。
  另一个我明写的克制点：**入口展开那条（R-a）不在这张票**，它归**票 102**，因为它在冻结的 C26 实现里、
  影响面比 winsec 宽得多——放一张票会让"守卫"级别的活被"契约级"的活拖住（票 89/95 的分票理由同形）。
  next= 排 in **票 102 之后**（fail-open 优先），但可与它并行——只要不撞同一文件。

- 2026-09-21 18:0x（`agent-ticket103`）开工登记 + **修前红**（本枚 commit 只含用例，生产码一行未动）：
  `git log --oneline -5 -- internal/winsec/` 顶端是 `01e7007`（票 89 退回单#3#4），
  `git status --porcelain internal/winsec/ internal/memory/` **空** ⇒ 没有人的未提交改动挡路，
  我也没有覆盖任何东西；本轮只 commit 显式路径（两份新用例 + 本票面）。
  新增用例两枚：`internal/winsec/seam_guard_windows_test.go`（`package winsec_test`，**故意放在包外**，
  因为 AC#1 的威胁就是"包外能装"）与 `internal/memory/artifacts_junction_tripwire_windows_test.go`。
  **实测红名 + 断言原文**（`go test -v`，`=== RUN` 5 条 winsec + 1 条 memory，全在里面）：
  1. `TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone` **FAIL（6 条 RED）**：
     - `AC#1 leg 1 RED: a pass-through rubber stamp was installed into the sealing seam (now winsec_test.rubberStampResolver)`
     - `AC#1 leg 1 RED: refusing the fake left no audit record naming it; log was:` （空日志）
     - `AC#1 leg 2 RED: the refusal is not about the link, so this leg measures nothing: winsec: ...\data\link is not a directory`
       （这条是**偶然红**的正确用法：今天 `PrivateDirAll` 确实报错，但报的是"不是目录"，与链接无关 ⇒
       我把"拒的原因必须提到 reparse/ErrUnresolvedPath"写成断言，避免票 94 AC#3 那种"红在偶然原因"被当成绿）
     - `AC#1 leg 3 RED: SealFile(...\data\link\sub\keep-me.txt) through a junction returned nil`
     - `AC#1 leg 3 RED: S-1-1-0 was stripped from the foreign file, i.e. its DACL was rewritten: [S-1-5-18 S-1-5-32-544 S-1-5-21-...-1001]`
     - `AC#1 leg 3 RED: foreign file DACL changed: before=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001] after=[S-1-5-18 S-1-5-32-544 S-1-5-21-...-1001]`
     ⇒ **AC#1 的"什么都没发生不算绿"这一条已经能用 SID 级读数咬住**（外来文件的 `S-1-1-0` 由有到无）。
  2. `TestAC1SeamIsSingleUse` **FAIL**：`AC#1 RED: the seam was re-installed from winsec_test.narrowOnlyResolverB ... it is not single-use`
     （两枚**都能拒**的合法解析器之间可以随意换 ⇒ 注册口是活的读—改写。）
  3. `TestAC2RemoveUnlinkedRefusesAPathThroughAJunction` **FAIL（2 条 RED，PROBE F 在树内复现）**：
     - `AC#2 RED: RemoveUnlinked("...\data\link\sub\keep-me.txt") through a junction returned nil`
     - `AC#2 RED: the foreign file behind the junction was deleted (GetFileAttributesEx ...ictim\sub\keep-me.txt: The system cannot find the file specified.)`
     （只在临时目录里造 junction、删的是测试自己的临时文件，真数据零风险；`mklink /J` 无需特权。）
  绿的对照组（同一次运行里点名，免得被当成"全红"）：
  `TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink` **PASS**（独立链接照旧只解链、不碰目标），
  `TestAC1RefusedInstallLeavesTheSealWorking` **PASS**（AC#3② 要求的"必须仍有一条绿"的用例），
  `TestAC2ReclaimWalkNeverDescendsIntoAJunction` **PASS**（tripwire 今天钉住的是
  `WalkDir` 对 junction 只报一个条目：`walk over the junction entry reported [...rtifacts\junklink]`，
  下降即红 ⇒ AC#3③ 的靶子）。
  下一条命令：实现守卫（`SetPathResolver` 一次性 + 一致性探针 + 响亮审计；`ResolvePath` 对
  resolver 的**答案**再跑一遍内置底线；`removeUnlinked` 动手前逐级查祖先链是不是链接），
  然后跑 AC#4 的门与 AC#3 的三向变异。

- 2026-09-21 18:4x（`agent-ticket103`）**修复落地 + AC#3 三向变异 + AC#4 门**（commit `0717bf2` 生产码，前一枚 `184af22` 是修前红用例）：
  **AC#1（R-c）已修，形状保留**：`SetPathResolver` 现在三道守卫——①一次性（缝上已有解析器时换一个不同的进来会被拒；
  唯一允许的退路是 `nil`＝内置底线，它只能拒不能改写，所以"拆回底线"永远安全，票 94 测试里
  `resolve_windows_test.go` 的拔—装复位照旧走得通）；②一致性探针（`resolverProbeShapes()` 两个纯拼写形状：
  相对路径与 `<TempDir>\..\x`，都是底线必拒的，可接受的答案只有"拒"或"交回一个底线认的改写"——
  恒说 OK 的橡皮图章两头都不沾，且探针不需要特权、不建 fixture、不动文件系统，所以能在 init 期跑）；
  ③拒装写 ERROR 级 slog 审计并保留在位者，**不 panic**（init 顺序会把整个二进制带走，拒装只让它落回底线）。
  **另加结构性的一半**（守卫不可被绕过的原因）：`ResolvePath` 现在把已装 resolver 交回的**答案**再过一遍
  `builtinVerifier`，所以即便有假解析器进了缝，带链接/`..`/空段/`\?\` 的答案也到不了 `os.Mkdir`；
  错误点名"哪一方给的答复"。**没有**倒回 `filepath.Abs`（票 94 的账不清零）。
  绿读数（`-count=2`，逐条点名）：`TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone`、
  `TestAC1SeamIsSingleUse`、`TestAC1RefusedInstallLeavesTheSealWorking`、
  `TestAC2RemoveUnlinkedRefusesAPathThroughAJunction`、`TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink` 全 PASS；
  SID 级判据此刻是"外来文件 `S-1-1-0` 前后都在"：`before=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001] after=同`。
- **AC#2（R-b）只做了守卫 + tripwire，语义修在 winsec 自己包里，一行 `internal/memory` 生产码都没动**：
  `RemoveUnlinked` 动手前跑 `firstLinkAncestor(path)`（祖先链逐级 Lstat，**不含叶子**——解链本身就是这个口的用途），
  Windows 用 `isReparsePoint`、POSIX 用 `ModeSymlink`（junction 在 Windows 的 `Lstat` 里**不报** ModeSymlink，
  实测 attributes=0x410，所以这个谓词必须分平台）；前缀用平台分隔符拼接，不用 `filepath.Join/Clean`
  （在这里做词法归一化，抹掉的正是它要找的形状）。`internal/memory/artifacts_junction_tripwire_windows_test.go`
  钉住回收遍历对 junction 只报一个条目（`walk over the junction entry reported [...\artifacts\junklink]`）+ 结果面
  （目标树文件与目录都在）+ 阳性对照（普通 stray-dir 照样被清，防止"什么都没删"被当成绿）。
- **AC#3 变异三向**（全部在 `/tmp/wisp-103-agent-ticket103`＝`git archive HEAD` 的仓外快照里做，目录名带会话后缀；
  每发都在同一条 `&&` 链里 `grep -n` 打印被改后的整行，且 `go build` rc=0 先量到；`go test -v` 数 `=== RUN`）：
  - **MUT-1 去掉注册守卫**（两行：`for _, probe := range resolverProbeShapes() {`→`range []string{} {`（落地行 resolve.go:152）、
    `if resolver != nil {`→`if resolver != nil && false {`（落地行 resolve.go:115））：`=== RUN` 3 条，
    **2 FAIL / 1 PASS**：`TestAC1SeamRejectsARubberStamp...` 红在 `AC#1 leg 1 RED: a pass-through rubber stamp was installed
    into the sealing seam (now winsec_test.rubberStampResolver)`（**红在守卫**），`TestAC1SeamIsSingleUse` 红在 seam-not-single-use；
    witness `TestAC1RefusedInstallLeavesTheSealWorking` 仍 PASS。
  - **MUT-2 注册了但被拒、只拆审计**（`slog.Error("winsec: refusing to install a path resolver...`→`slog.Debug(`，落地行 resolve.go:107）：
    `=== RUN` 3 条，**1 FAIL / 2 PASS**，唯一红是 `AC#1 leg 1 RED: refusing the fake left no audit record naming it; log was:`（空日志），
    而"缝上仍是原解析器 + 外来 DACL 的 `S-1-1-0` 一字未动 + 密封照常工作"三条读数全绿
    ⇒ **证明 MUT-1 的红在守卫身上，不在噪声身上**，且审计断言本身是活的。
  - **MUT-3 tripwire 的"不下降"改成"下降"**（`removeStray` 的 `filepath.WalkDir(...)`→`descendFollowingLinks(...)`，
    一个用 `os.Stat`+`os.ReadDir` 跟链接的递归；落地行 artifacts.go:155/322，`go build` rc=0）：
    `TestAC2ReclaimWalkNeverDescendsIntoAJunction` **FAIL**，红在读数
    `removeStray(junklink): ... winsec: entry is a link ... the spelling reaches it through the link at ...\artifacts\junklink`
    ⇒ 一旦有人把遍历接成下降的，**立刻红**（并且因为 AC#2 的守卫在位，它是"响亮地红"而不是"别人树里的文件静默消失"）。
    诚实登记两点：该用例 leg 1（断言 `WalkDir` 本身不下降）测的是标准库行为，MUT-3 不会让它红——咬住的是结果面那条腿；
    第一发用 `filepath.EvalSymlinks` 的变异**没有落地成下降**（Windows 上它对 junction 不解引用，读数 `--- PASS`），已换成上面的版本重测。
    `TestStrayRemovalDoesNotFollowLinks`（票 79 的既有用例）在 MUT-3 下仍 PASS ⇒ 它钉不住下降，这条 tripwire 是新增量。
  - **还原**：三发变异之后把快照里被改的文件用 `git show HEAD:<path>` 逐字节写回，
    `diff -q` 对 resolve.go / winsec.go / winsec_windows.go / winsec_other.go / artifacts.go 全部 `CLEAN`，
    并把快照与一份重新 `git archive HEAD` 的副本 `diff -rq` **零差异**（无残留）。仓库目录内**未建 worktree / 未 checkout**。
- **AC#4 门（真实树）**：`go test -count=2 ./internal/winsec/ ./internal/memory/ ./internal/secret/ ./internal/risk/`
  → `ok winsec 12.436s / ok memory 26.597s / ok secret 0.570s / ok risk 7.798s`，FAIL=0；
  `go test -v ./internal/winsec/ ./internal/memory/` `=== RUN` **102** 条，**SKIP 恰 1 条**：
  `TestSubprocessCrashWriter`（`internal/memory` 既有、与本票无关的环境条件用例，点名在此不算静默绿）；
  受我改动辐射的 `./internal/agent/`（1.903s ok）`./internal/tools/`（13.751s ok）同跑。
  `gofmt -l internal/winsec/ internal/memory/` 空、`$(go env GOPATH)/bin/gofumpt -l` 空；
  `go vet ./internal/winsec/` 与 `GOOS=linux`/`GOOS=darwin go vet ./internal/winsec/`、`GOOS=linux go vet ./internal/memory/` 全 rc=0（按包，仓根整树跑是既有坑不重踩）。
  `sh scripts/d22scan.sh`：**clean - no D22 ban violations**（bans #1-5 internal/=197、cmd/=20、#6 frontend/=40、#7 internal/tools/=17、
  #8 design/=16、frontend/=40、internal/=342 Go 文件含注释与 `_test.go`、cmd/=26；runtests.sh 侧 packages=[./...] top-level PASS=21 FAIL=0 SKIP=0）。
  ⚠ 两点环境账，不是我的红：`go build ./...` 在本机当前树上 rc=1，坏在**别人未提交**的
  `internal/config/parse.go:214 undefined: winsec` + `migrate.go "os" imported and not used`（`git status` 里是别人的活，我没碰），
  故按包量；`TestResolvePerCallBudget`（票 86 的负载假红）本轮未出现。
  另：`tools/d22scan` 我一行未改，`allowlist.txt` 未改，冻结清单未碰；未 push。
  **残留交回**：AC 框我不自勾（裁决表 `docs/evidence/s1/103-*.md` 归验收方）。
  next= `acceptor-ticket103` 按 AC#1..#4 逐格裁；重点复核三处：① 守卫是否给了伪造留了别的门
  （`%VAR%`/前导 `~` 的**改写型**伪造仍归票 102 的 `Actable()`，winsec 看不出"善意改写"与"劫持改写"的区别，
  我在 `resolve.go` 的注释与本 log 里都写明了这条边界）；② POSIX 侧守卫我只能编译期验（`GOOS=linux/darwin go build`+`go vet` rc=0），
  **探针集合在 Linux 上是否会把 risk 的纯词法解析器拒掉，本机测不到**——若 CI 跑 Linux，这是第一条要看的眼色；
  ③ AC#2 的祖先检查是 fail-closed：数据根若真的放在别人 symlink 底下，回收会开始报错（方向是拒，不是删），
  这是有意的语义收紧，验收方若判它过界请说一声。

## 独立对抗验收交回（`acceptor-ticket103`，2026-09-21；append-only，前文一行未动）

**总判：PASS WITH CONDITIONS（通过但有条件）** —— 裁决表 `docs/evidence/s1/103-adversarial-acceptance.md`（AC#1..#4 1:1，含档位标注与命令原文）。
AC 框仍不自勾由验收方裁：本代理的裁法是 **AC#3 无条件通过；AC#1/AC#2/AC#4 通过但有条件**。

**三处点名的直答（细节与读数在裁决表）**
- **① 边界守住了，但描述窄了一格**：`%VAR%`/`~` 的改写型伪造确实被 102 的账拦在缝前（`internal/risk/winsec_c26.go:39-44` → `pathresolver.go:93-100` 的 `Actable()` → `winsec/resolve.go:210` 包成密封失败）。但"看不出善意改写/劫持改写"不止适用于 `%VAR%`/`~`：实测**任何**已装 resolver 交回"干净但换了树"的答案，winsec 照封（`resolve.go:222` 只认拼写形状，不认树归属）⇒ 那一半没有票接着，记 `R-103-1` 的 remediation。
- **② POSIX 验到了**，不必只靠编译期：本代理在**真 Linux 容器**（`docker run golang:1.27`，快照 = `git archive HEAD`）量到 `probes_passed=1`、`--- PASS: TestAC103POSIXRiskResolverPassesTheProbe`（**探针没有把 risk 的纯词法解析器拒掉**）、`--- PASS: TestAC103POSIXRemoveUnlinkedAncestorGuard`，四包 `ok winsec/risk/memory/secret` 干净复跑 RC=0；CI 侧对得上 **run `35587986855`（run_number 148，head `b2fa2ed` ⊇ `0717bf2`）** 步级 `Portable package tests` 里 `ok internal/risk 1.130s`、`ok internal/memory 10.748s`（ubuntu）。**"本机全绿 ≠ CI 过了"仍照登**：148 的 job 结论是 failure，红因不属本票（`internal/tools` 的 `TestPathCanonicalizerAccountsForRewrittenRoots`＝102/106 地界；`test-windows` 的 `internal/secret` 红在 `verifyPrivate: ... ([LA LA])`＝runner 账号映射，`verifyPrivate` 不在 `0717bf2` diff 内）；**当前 HEAD 没有跑完过的 run ⇒ 欠编排者 push 后复跑**。
- **③ 判为契约要求的从严**（方向是拒不是删），但"正常路径不受影响"只到间接读数（`TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink` PASS、memory/secret 本机与 Linux 容器全绿），**没有一条用例钉住"数据根真在别人 symlink 底下 ⇒ 报 ErrIsReparsePoint"这一形态**；POSIX 用 `ModeSymlink` 查祖先 ⇒ macOS 的 `/tmp`、`/var` 本身是链接，以 tmp 为根的回收会开始报错。这属语义收紧，不判过界（票面 AC#2 原句就是"必须拒"），但记 `R-103-7` 请编排者裁"从严即可／要豁免通道"。

**本代理独立复现的红与变异（不抄实现方读数）**
- 两条腿 + SID 级读数：`TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone`、`TestAC1SeamIsSingleUse`、`TestAC1RefusedInstallLeavesTheSealWorking`、`TestAC2RemoveUnlinkedRefusesAPathThroughAJunction` 全 PASS（`-count=2` 两样本）；真机 `mklink /J` 在位（`attributes 0x410`），外来文件 `S-1-1-0` **前后都在**。
- 变异三向自做：MUT-1 两行（`resolve.go:115/152`）⇒ 5 RUN、3 PASS/2 FAIL，红名都在守卫上；MUT-2（`resolve.go:107` Error→Debug）⇒ 8 RUN、7 PASS/**恰 1 FAIL**（只红审计断言，witness 绿）；MUT-3（`artifacts.go:155` → `descendFollowingLinks103`）⇒ `TestAC2ReclaimWalkNeverDescendsIntoAJunction` **FAIL**，且票 79 的 `TestStrayRemovalDoesNotFollowLinks` 仍 PASS（新 tripwire 是真增量）。还原：5 文件 `diff -q` 全 CLEAN，快照复跑 `ok winsec 7.127s / ok memory 13.519s`。
- 攻守卫（我自己的探针，只活在 `/tmp` 快照，未进仓）：**P1b 换个注册顺序**（先 `nil` 释放再装恒改写型伪造）⇒ 探针收下、`SealFile` 返回 nil、外来文件 `S-1-1-0` **由有到无**；**P2 分隔符改写**（`/` 或混合）⇒ `firstLinkAncestor` 一个祖先都不查，沿 junction 删掉外来文件且返回 nil；**P3** 纯底线 + `/` 拼写 ⇒ 同一错法在 `placement_windows.go`（**不是 `0717bf2` 改的文件**）上改掉外来 DACL。类型名撞名方向 fail-safe（伪造进不来），已量到 ERROR 读数。
- 门禁复跑：`gofmt`/`gofumpt`（GOPATH/bin 有二进制，跑了）空；按包 `go vet` 原生 + `GOOS=linux` + `GOOS=darwin` 全 rc=0；`sh scripts/d22scan.sh` rc=0 clean、各 scope 文件数不降（`#8 internal/` 342→346 升）；四包 `-count=2 -v` 里 **FAIL 1 条 = `TestResolvePerCallBudget`（隔离复跑 PASS：0.779 ms/op vs 1.000 ms，票 86 负载假红，两样本都报）**；**SKIP 逐条点名 = `TestSubprocessCrashWriter`（memory）+ `TestSyncRegistryProbeLive`（risk）** ⇒ 你那句"SKIP 恰 1 条"只对 winsec+memory 成立，且"102"是每轮数（`-count=2 -v` 原始 204 行）＝`R-103-5`。

**缺陷登记（验收期未修，各该自开票）**：`R-103-1` 一次性守卫可被注册顺序解除 + 缝上答案的树归属无人查（安全）；`R-103-2` `firstLinkAncestor` 分隔符可绕（守卫可绕，陷阱档）；`R-103-3` `platformVerifyPlacement` 同形既有洞（**非本票引入**，票 94 底线的账）；`R-103-4` 恒等重装分支覆盖率 0 ⇒ 幂等语义没钉；`R-103-5` 计数/SKIP 口径；`R-103-6` darwin 无运行期读数 + Linux 探针少一形状；`R-103-7` symlink 数据根的报错形态无用例。
**未碰**：`internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`、`docs/PLAN.md`、`docs/specs/*`、`internal/memory` 生产码；未 push；仓内未建 worktree。
