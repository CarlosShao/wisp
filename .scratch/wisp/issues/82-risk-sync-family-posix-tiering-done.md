# 82 — `internal/risk` 的 sync-root 家族在 ubuntu 上 8 条红：判据没写错，是**POSIX 侧这条探测根本没实现**

**Status:** coded-awaiting-review（六框自评全 PASS，证据与两侧数字在 Progress log；`next=` 见最后一条）
**Type:** 门禁归属裁定（A52④）+ 测试分层落地
**Blocks:** 票 70 的 AC#6（`test-core` 在票 81/82 之前不可能绿）· **Blocked by:** nothing
**Packages:** `internal/risk/` 的**测试文件**（`syncdirs*_test.go`、写门/取证相关用例）。
**禁改**：`internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、`rules_gateway.go`、
`docs/PLAN.md`、`docs/specs/*.md`、`internal/config/**`（票 80 的地界）、
`internal/agent`+`internal/memory` 的测试（票 81 的地界）、`tools/d22scan/**`、`allowlist.txt`。

## 实测证据（CI run `35558750456`，headSha `17efc2c`，`test-core` → `Portable package tests`，ubuntu）

红着的 8 条（`gh run view --log-failed` 逐条点名）：
`TestWriteGatePlainLocalWriteNotFlagged`、`TestWriteGateNotSelectedByPayloadKey`、
`TestWriteGateEveryPathTargetJudged`、`TestWriteGateAllPathsNonSyncStaysExempt`、
`TestSyncNormalNewFileWriteNotFlagged`、`TestSyncFallbackNotDisarmableByWeakRoot`、
`TestSyncEnvConfiguredRoots`、`TestSyncDotDotTailFailsClosed`。

根因不是断言写错：`syncdirs_other.go` 顶部自带 `DEFERRED(P12-macos)`，`registryProbe` 在非 Windows 直接返回 nil
⇒ POSIX 上"sync root"这个概念**尚未实现**（那是**票 55** macOS 移植的活）。
所以这 8 条在 ubuntu 上测的是一个不存在的能力。

## 编排者已裁定（照此执行，不要再问）

**不**把这些包/用例从 portable 清单里搬走（搬走 = 把覆盖面挪到没人看的地方）。
采用票 70-c 为 `internal/secret` 用过的**同一档修法：分层**——
Windows 专属判据进 `//go:build windows`，POSIX 侧留一条**断言当前真实行为**的用例。
机制票 75 已经铺好：`internal/risk/provenance_syncdirs_windows_test.go`（windows 层，7 个子用例）
与 `internal/risk/provenance_syncdirs_other_test.go`（`!windows` 层），
后者里那条由 `SyncDetectionComplete()` 门控的**休眠断言**（A51⑧）就是 POSIX 侧的占位。

## AC（1:1，裁决表 `docs/evidence/s1/82-*.md`）

- [x] **AC#1** 上面 8 条**逐条定性**并列表：哪些属于"只在 Windows 上有意义"（进 windows 层）、
  哪些其实**两侧都该跑**（那就是真缺陷，不许搬走，就地修）。
  ⚠ 不许整族打包贴 tag —— 那正是 A49② 里票 78 代理拒绝过的"整包被静默排除"。
  **判据**：每条给一句"为什么这条在 POSIX 上没有可断言的对象"或"这条两侧都有对象、留下"。
  ⇒ 定性表见 Progress log 2026-09-21 第二条（6 条两侧都有对象、就地修；2 条的"确认等级拆掉兜底网"半边进 windows 层）。
- [x] **AC#1b（2026-09-21 由 A57④ 追加，编排者写的判据）** 碰 `internal/risk` 的**折叠语义**时，
  必须同时检查 `internal/risk/rules_test.go:198` 是否**还在测东西**：它用
  `C:\Program Files\Git\bin\git.exe` 这条字面 Windows 路径做输入，今天两侧同字节同结论（不是问题），
  但它断言的是折叠行为——**一旦本票把折叠收窄成"只在 Windows 上折"，这条会在 ubuntu 上静默变空仪器**
  （跑过、绿、什么也没测到）。判据：给出该用例在**两侧各自的**断言命中数或红/绿读数，
  若它在 Linux 侧变成空转，就地改成平台派生的期望值（照票 81 的做法），**不许留着当装饰**。
  ⇒ **本票没有收窄任何折叠**（落点只有测试文件；`shellAllowlisted` 的 `LastIndexAny(argv0, "/\\")` 一行未动）。
  两侧读数 + 变异自证见 Progress log 2026-09-21 第三条末段：两侧绿→改 `"git.exe"`→`"git"` 后**两侧同一条红**。

- [x] **AC#2b（编排者据 A49② 追加：本票最怕的失败模式）** 分层**必须用构建约束表达，不能用文件名表达**：
  `*_windows_test.go` 这个名字**什么门都不施加**——真正的门是文件里 `//go:build windows` 那一行。
  若有人只改文件名、或让整包退出构建，`go test ./internal/risk/` 会打 **`ok ... [no test files]`**：
  看着全绿，实际**整包判据被静默跳过**（票 78 的代理当初就是按这条拒绝了我那个"省事修法"）。
  ⇒ **机器可检查**：windows 层与 POSIX 层各跑一次 `-count=1 -v`，**两侧都要贴出自己的 `=== RUN` 计数与测试名清单**；
  任一侧出现 `[no test files]` / `no tests to run` / RUN=0 ⇒ **本票不过**。
  ⇒ 两侧计数与名单见 Progress log 第三条。**两处 `no test files` = 0、`no tests to run` = 0、RUN=0 = 无**。
  ⚠ **对本票面的一句事实更正（要编排者拍板，我没改票面正文）**："`*_windows_test.go` 这个名字什么门都不施加"
  **不成立**——Go 的文件名构建约束是先剥掉 `_test` 再看尾部 `_GOOS`，所以 `syncdirs_windows_test.go`
  **光靠文件名就被 gate 到 windows**（实测：把 `//go:build windows` 删掉、文件名不动，ubuntu 上打
  `testing: warning: no tests to run` + `ok … [no tests to run]` + rc=0，见 AC#4 段 [1]）。
  反向不成立：`*_other_test.go` 里 "other" 不是 GOOS，**POSIX 层唯一的门就是那行 `//go:build !windows`**，
  所以我按 AC#4 做了**对称变异**（在 Windows 上剥掉那行 ⇒ POSIX 层 3 条里 2 条真红）。
  含义：**"只改文件名不改 tag"这条路在 windows 侧其实是生效的（靠文件名），正因如此它更隐蔽**——
  审 code 的人看不见 tag 行就以为没分层。票 75/78 的结论（分层要用 tag 表达）我照办了，
  但建议在 A5x 里把这句判据写成"两层都必须有自己的 `//go:build` 行，且**两侧各做一次剥 tag 变异**"。

- [x] **AC#2b 的编排者裁定（2026-09-21 16:2x）＝接受上面那条事实更正**："`*_windows_test.go` 不施加门"这句**是我写的，且它错了**——Go 的文件名构建约束先剥 `_test` 再看 `_GOOS` 尾部，所以 windows 层**光靠文件名就已经被 gate**（`acceptor-ticket82` 独立复现同一读数：只删 tag、保文件名 ⇒ ubuntu 打 `no tests to run` + **rc=0，不红**）。原句**留在上面不删**（本仓规矩：保留原文 + 追加更正）。⇒ 判据最终版：**两层都必须有自己的 `//go:build` 行，且两侧各做一次"剥 tag 变异"**；并且**"只改文件名"在 windows 侧真的生效**这件事本身要写进简报——它让漏写 tag 的人**看不到症状**。

- [x] **AC#2** POSIX 层**不许是空文件、不许全是 `t.Skip`**：它必须断言 POSIX **当前**的真实行为
  （fail-closed / 尚未检出），并在注释里指名**票 55 落地时这一层要怎么转成活的**（A51⑧）。
  ⇒ `internal/risk/syncdirs_other_test.go`：3 条活用例、**0 个 `t.Skip`/`t.Log` 早退**（ubuntu 上 3/3 `--- PASS`），
  头部写死票 55 落地的 1-2-3 步转活程序（搬 windows 层用例 / 删本文件两条 / 解 `provenance_syncdirs_other_test.go` 的门控分支）。
- [x] **AC#3** 双向证据：(i) ubuntu（docker `golang:1.27`，挂 `git archive` 快照）这 8 条**不再出现**在 FAIL 里，
  且 `internal/risk` 的**红数不增加**；(ii) Windows 本机 `go test ./internal/risk/ -count=2` 全绿，
  `=== RUN` 行数是 `-count=2` 的偶数倍，`--- SKIP` 逐条点名。
  ⇒ (i) ubuntu `TOPFAIL=0`（改前 8）、(ii) Windows rc=0 / `=== RUN`=302=2×151 / SKIP 只有 `TestSyncRegistryProbeLive`×2 —— 数字在 Progress log 第三条。
- [x] **AC#4** 变异：把 windows 层的 tag 去掉 ⇒ 在 ubuntu 上它必须**变红**（证明"tag 不是在藏东西，
  而是这条判据在 POSIX 上确实无对象"）。锚点=承载行为的那一行，同链 grep 自证，还原后 `diff -q` 对拍。
  ⇒ 剥 tag 不够（文件名还压着一道门，见 AC#2b 的更正），所以变异体另存成 **无 windows 后缀的文件名**再跑：
  ubuntu 上两条 windows 层用例 **`--- FAIL` ×2、rc=1、失败信息就是那三句等级判据**（不是编译失败）。
  对称方向也做了：Windows 上剥掉 POSIX 层的 `//go:build !windows` ⇒ 3 条里 2 条真红。全程在 `/tmp` 快照里做，工作树 `git diff -- internal/risk/` 干净。
- [x] **AC#5** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <file>` 空、`go vet` rc=0、
  `GOOS=linux go vet` rc=0、`go test -count=2` rc=0。
  ⇒ 见 Progress log 第三条；`GOOS=linux go vet ./internal/risk/`（**按包作用域**，不是 A54③ 禁的 `./...`）rc=0，另加 `GOOS=darwin` rc=0。


## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc；禁 `git add -A`/`.`；commit 前核对 `git diff --cached --name-only`；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；不在仓内建 worktree（A38④）；
docker 挂**快照**不挂工作树；票面 Progress log append-only，**要改的那行先读再替换**；
四种假绿逐跑点名（SKIP 当 ok / `-run` 空匹配 / `-count=N` 没核对倍数 / 步骤被静默跳过）。

## Progress log（append-only）

- 2026-09-21 16:2x（编排者，**独立对抗验收判 PASS ⇒ 本票结案；两条改名前置当场落**）：
  `acceptor-ticket82` 在 `git archive 629ce3c` 的仓外快照 + docker `golang:1.27`（命名卷带会话后缀）里
  把六个框逐条自己跑过：ubuntu `internal/risk` **rc=0 / RUN 124 / PASS 123 / FAIL 0 / SKIP 1**，
  那 8 条**逐条 `--- PASS`**；改前基线 `42ed13f^` 独立复现回**顶层红 8 条、名字一字不差**。
  四组变异（正向删 tag、只删 tag 保文件名、Windows 侧剥 `!windows`、四次生产码行为变异）**全部如预期**，
  还原用 `git archive -- internal/risk` + `diff -r` rc=0 证明。⇒ **判 PASS，可挂 `-done`**。
  - **R-1（已落在本票面 AC#2b 那条裁定里）**：那条"文件名不施加门"的更正**我接受了**，
    它把判据推向更严的版本（见 AC#2b 末句）。
  - **R-2（当场登记并转票 55 当锚点）**：验收代理的 **MUT-c** 实测：把 `gradeConfirmed` 加宽去认 `"default"`
    ⇒ **ubuntu rc=0 全绿**、Windows rc=1 ⇒ 本票表里第 6 行那句"弱等级永不拆网**两侧有对象**"**是 over-claim**
    （Linux 侧那条其实是**空仪器**）。⇒ 已追加进票 55 的锚点段；**本票不改表正文**（保留原文 + 这条更正）。
  - **R-3（本票之外、已立案）**：`TestSyncRegistryProbeLive` 在 ubuntu 是**结构性永久 SKIP**，
    且它的 skip 文案在 Linux 上说了一句 Windows 的谎 ⇒ 归**票 93**（portable 步把 SKIP 记成 ok）。
    本票 portable 步逐字重放读数：**TOPSKIP=7，全部点名 + 归包**。
  - **R-4（归票 85）**：`lint` 里 `staticcheck` 一红就**连带 `mockllm module vet` 被 skipped**
    ⇒ 与 A44① 同族的"步骤次序造成的结构性空窗"，判据是 `if: always()`。
  - **顺带一条大事（对票 70 的账有用）**：验收代理取到真 CI 读数——`test-core` 自 `42ed13f` 起
    **连续 4 个 run success**（`942ab5a` 逐步骤全绿），更早的 `e5e5eb7` 是 failure
    ⇒ **"ubuntu 的 portable 测试全绿"今天第一次成立**，票 70 的 AC#2 可据此重判；
    但 AC#6 判据是"五 job 全 pass"，`lint` 仍红在 `staticcheck` ⇒ 那条**还没**成立。
  next= 挂 `-done`（两条前置已落）。本票**没有留未接的码**；R-2/R-3/R-4 的归宿分别是票 55 / 93 / 85。

- 2026-09-21 写码代理认领本票（HEAD `c923327`，局面基线 = CI run `35562680354` / headSha `e5e5eb7` 的 8 条同族红）。
  本枚 checkpoint 只改票面 Status + 这条 log，先落账再动代码。
  读过的机制件：`internal/risk/syncdirs_other.go:20 registryProbe` 直接 `return nil`（`DEFERRED(P12-macos)` ticket 55）、
  `internal/risk/pathresolver_other.go:10 resolveHandle` 返回 `("", false)` ⇒ `syncSet.add` 的 `canonical` 恒 false
  ⇒ `syncSet.finalize` 的 `s.complete` 恒 false ⇒ `SyncDetectionComplete()` 在 POSIX 恒 false。
  **本票的落点推论（AC#1 定性前先记）**：这 8 条**不是**整族都"在 POSIX 上无对象"——
  `match()` 只在 `!s.complete && isUnder(cand, s.home)` 时兜底，所以**只要被测路径落在 profile home 之外**，
  POSIX 上"是不是 sync 通道"就仍然**只由 root membership 决定**（两侧都有对象）。
  ⇒ 分层方案 = ①把 4 条 `TestWriteGate*` + `TestSyncDotDotTailFailsClosed` + `TestSyncNormalNewFileWriteNotFlagged`
  的**引擎形状**改成两侧可判（root 与对照目录都在 home 之外），把"confirmed 等级本身"这条判据留到分层里；
  ②`"确认等级能否拆掉兜底网"` 这一半在 POSIX 确实无对象，进 windows 层，POSIX 层留一条断言
  "今天任何等级都拆不掉（fail-closed）"的活用例，并指名票 55 落地时怎么转活（AC#2）。

- 2026-09-21 **AC#1 八条逐条定性**（"POSIX 上没有可断言的对象"这句必须能落到具体一行代码上，所以每条都点了那一行）：
  先给共性根因：`syncSet.finalize` 只在**同时**满足 `canonical`（=`Resolve().Resolved`，POSIX 恒 false，见
  `pathresolver_other.go:10`）与 `gradeConfirmed(Source)` 时置 `complete=true`。这 8 条里凡死在
  `SyncDetectionComplete()` 前置条件上的，**死因是"POSIX 造不出 canonical root"，不是"POSIX 判不出 sync"**——
  `match()` 在兜底网之前就先比 root membership，且兜底网**只覆盖 profile home 之下**，所以
  **把 fixture 的两个目录搬到 home 之外，membership 两侧都判得出**。⇒ 六条属"两侧都有对象"，就地修；
  | # | 用例 | 定性 | 处置 |
  |---|------|------|------|
  | 1-4 | `TestWriteGate{NotSelectedByPayloadKey,PlainLocalWriteNotFlagged,EveryPathTargetJudged,AllPathsNonSyncStaysExempt}` | 两侧都有对象：被测的是 writeGate 的"豁免只看调用形状、不看 payload 键名"，平台无关代码；红的唯一来源是 `m7Engine` 那句 `SyncDetectionComplete()` sanity | 就地修：`m7Engine` 换成 home 之外形状（POSIX 可判），并断言 `Root.Source != "suspect-fallback"` 自证判决来自 membership |
  | 5 | `TestSyncNormalNewFileWriteNotFlagged` | 两侧都有对象（B-1 假阳性锤：新文件落在非 root 目录不是通道）；POSIX 无对象的只是 `sanity: registry-grade root is confirmed` 那一行 | 就地修（`membershipEngine`）+ 补正向对照（同一引擎写进 root 必须 hit），防"什么都不 flag" |
  | 8 | `TestSyncDotDotTailFailsClosed` | 两侧都有对象：`hasFoldedDotDot` 是纯词法判据，跑在任何 OS 咨询之前；红的来源仍是那句 confirmed sanity | 就地修：换 home 之外形状后 `filepath.Clean(raw)==folded` 前提两侧同构 |
  | 6 | `TestSyncFallbackNotDisarmableByWeakRoot` | **一半一半**：弱等级永不拆网两侧有对象；"registry/config/env 拆网 + 拆网后 under-profile 判非 sync"在 POSIX 无对象（拆网要 canonical） | 拆：portable 留弱等级半边；`TestSyncConfirmedGradesDisarmFallbackWindows` 进 windows 层；POSIX 层 `TestSyncNoGradeIsConfirmedOnPosix` 断"7 个等级名一个都不拆" |
  | 7 | `TestSyncEnvConfiguredRoots` | **一半一半**：`envConfiguredRoots` 是 `probeEnv` 的纯函数（等级/provider/空白剔除）、写进 env root 判 sync 且归因 `Source=="env"` —— 两侧都有对象；"env ⇒ complete" 无对象 | 拆：portable 留纯函数 + membership 归因半边；complete 半边进 windows 层 |
  ⇒ 没有任何一条被"整族打包贴 tag"：新增的 windows 层只装上面**两条**用例的 confirmed-grade 半边（2 个用例），
  其余 6 条反而**多**在 POSIX 上跑了起来（覆盖面净增，不是搬走）。
  POSIX 层另加 3 条活用例（`syncdirs_other_test.go`，无一条 `t.Skip`）：等级现实、membership 现实、
  `registryProbe` 桩本身（票 55 落地即红的 tripwire）。

- 2026-09-21 **两侧数字（AC#2b / AC#3 / AC#4 / AC#5）**，代码 = commit `42ed13f`。
  **Windows 本机**（`go test ./internal/risk/ -count=2 -v`，工作树）：rc=0；
  `=== RUN` 总（含子用例）**302 = 2×151**、顶层（无 `/`）**180 = 2×90**；`--- PASS` 300、`--- FAIL` **0**；
  `--- SKIP` **2** 条，逐条点名 = `TestSyncRegistryProbeLive` × 2（本机 HKCU Accounts 无 UserFolder，P12 文档里的现实；改前也是这条 SKIP，非本票新造）；
  `[no test files]`=**0**、`no tests to run`=**0**。
  windows 层自己的 RUN 计数：`TestSyncConfirmedGradesDisarmFallbackWindows` 2、`TestWriteGatePlainTargetInsideProfileWindows` 2
  （票 75 的 `TestExfilSyncWriteWindowsSpellingInvariant` 2）；POSIX 层 3 条在 Windows RUN=0（被 tag 挡住，正是分层的含义）。
  `gofmt -l internal/risk/` 空、`gofumpt -l internal/risk/` 空（新文件写完才跑的）、`go vet ./internal/risk/` rc=0、
  `GOOS=linux go vet ./internal/risk/` rc=0、`GOOS=darwin go vet ./internal/risk/` rc=0（**按包作用域**，A54③ 禁的是 `./...` 那条永远 rc=1 的命令）。
  **Linux（docker `golang:1.27.1`，`git archive HEAD` 解到 `/tmp/w82run3`，未挂工作树）**：
  改前基线（同一容器、同一包、HEAD 前一格）`--- FAIL` = **8**（正是票面点名的 8 条，全部死在 `SyncDetectionComplete()` 前置条件那一行）；
  改后 `-count=1 -v`：rc=0、`=== RUN` 总 **123**、顶层 **71**、`--- PASS` 70、`--- FAIL` **0**、`--- SKIP` **1**（`TestSyncRegistryProbeLive`）、
  `[no test files]`=**0**、`no tests to run`=**0**。8 条逐条读数：`TestWriteGate{NotSelectedByPayloadKey,PlainLocalWriteNotFlagged,EveryPathTargetJudged,AllPathsNonSyncStaysExempt}`、
  `TestSync{NormalNewFileWriteNotFlagged,FallbackNotDisarmableByWeakRoot,EnvConfiguredRoots,DotDotTailFailsClosed}` = **8/8 `--- PASS`**；
  POSIX 层 `TestSyncNoGradeIsConfirmedOnPosix` / `TestSyncMembershipDecidesOffProfilePosix` / `TestSyncRegistryProbeIsAStubHere` = **3/3 `--- PASS`**；
  windows 层 2 条不出现在 ubuntu 名单里。**POSIX 侧覆盖是净增**：顶层用例里 `TestWriteGate*`×4 与 `TestSync*`×8 第一次在 ubuntu 上真跑。
  **AC#1b（`rules_test.go:198`）**：两侧绿（该表项的 3 个断言 = level/RulesHit 长度/R6 名，任一不满足就 `t.Fatalf`，无 `t.Run` 子用例可数）；
  变异 = 把 198 行的 `ShellAllowlist: []string{"git.exe"}` 改成 `{"git"}`（锚点就在这一行，`sed -n 198p` 前后对拍）⇒
  **Linux 与 Windows 同一条红、同一句信息**：`rules_test.go:205: clean allowlisted full path -> L1: want L1/[R6], got L2/[R6] (R6: 程序不在 shell 白名单: C:\Program Files\Git\bin\git.exe)`
  ⇒ 它在 ubuntu 上**不是空仪器**（`shellAllowlisted` 用 `LastIndexAny(argv0, "/\\")`，两种分隔符一起找，与 `unifySeparators` 的平台分支无关），
  本票也没收窄任何折叠 ⇒ 不需要照票 81 改成平台派生期望值。真正会被"只在 Windows 折"影响的是 `splitPathComponents`/`hasFoldedDotDot`，
  它们的 ubuntu 读数由 `TestSyncDotDotTailFailsClosed` + `TestResolveKeepsSeparatorsOutOfPosixNames` 两条 PASS 兜着。
  **AC#4（三次变异，全在 `/tmp` 快照里做，仓内 `git diff -- internal/risk/` 空）**：
  [1] ubuntu 只删 `//go:build windows`、文件名不动 ⇒ `no tests to run` + rc=0 ⇒ **文件名自己就是一道门**（AC#2b 正文那句判据要更正，见上）；
  [2] ubuntu 把同一文件另存成无 windows 后缀的 `syncdirs_winmut_test.go` 且无 tag ⇒ **rc=1**、
  `--- FAIL: TestSyncConfirmedGradesDisarmFallbackWindows`（6 句：registry/config/env × "must count as confirmed" + "must lift the blanket suspect net"）、
  `--- FAIL: TestWriteGatePlainTargetInsideProfileWindows`（precondition 那句）⇒ 是真行为红、不是编译失败 ⇒ tag 没在藏东西；
  [3] 对称方向：Windows 上删 `syncdirs_other_test.go` 的 `//go:build !windows` ⇒ 3 条里 2 条真红
  （`Resolve(...)` 在 Windows 报 `Resolved:true`、under-profile 判 `Sync:false`），第三条 `TestSyncRegistryProbeIsAStubHere`
  在本机恰好也绿（这台机器 registry 无 UserFolder ⇒ 桩判据同形），所以 POSIX 层的承重断言是前两条，已在用例名里写明。
  **残留（本票没闭，要报）**：① `TestSyncRegistryProbeLive` 两侧都 `--- SKIP`，而 CI 的 portable 步是裸 `go test` ⇒ SKIP 记成 ok
  （只有 `tools/d22scan/runtests.sh` 拒 SKIP）；这是票 82 之外、与票 71 AC#3 同形状的残留假绿，我没动它（不属那 8 条红，动它=扩界）。
  ② `test-core` 整条 job 是否归零，要编排者在 CI 上复跑确认（本包 8→0 已两侧实测）。
  **next=** 编排者验收（三件待你拍板：① AC#2b 那句"文件名不施加任何门"要不要写进 A5x 更正；
  ② `TestSyncRegistryProbeLive` 的 SKIP-当-ok 要不要另开一张小票；③ 票 55 的 Constraints 里要不要把我那 3 条 POSIX 层的
  "落地即红" tripwire 名字抄进去当验收锚点）。我没动 `docs/PLAN.md`、`docs/evidence/s1/**`（票面写了裁决表路径 `82-*.md`，
  那张表按规矩是验收方出的，不是我自裁）。

- 2026-09-21 第二枚 commit（代码）：**Windows 侧数字已取，Linux 侧同批数字下一条追加，故本枚不勾任何框**。
  Windows（工作树，HEAD 42ed13f + 本枚）：`gofmt -l internal/risk/` 空、`gofumpt -l internal/risk/` 空（**新文件写完才跑**）、
  `go vet ./internal/risk/` rc=0、`go test -count=2 -v` rc=0：ALLRUN=302 / 顶层 TOPRUN=180（=2×90，偶数倍 ✓）/
  ALLPASS=300 / ALLFAIL=0 / ALLSKIP=2 / `[no tests to run]`=0 / `[no test files]`=0；
  两条 SKIP 逐条点名：**都是 `TestSyncRegistryProbeLive`**（×2，因为 `-count=2`；本机 HKCU Accounts 无 UserFolder，
  票面/ PRECHECK §P12 记的机器现实，**不是本票新造的休眠**，也不属那 8 条）。
  本枚另做三件 AC#2 要求的"POSIX 层不许当装饰"的收尾：
  ① `posixSyncEngine`（票 75 的 POSIX 层）现在真的建出 `OneDrive/Notes` 那一级祖先，
     于是同一用例把 ticket 75 的**两条路由都跑到**（祖先链可确认 vs 只能词法回退），名单里多一个子用例；
  ② `syncdirs_test.go` 里两处 `IsSyncPath(` + 字面 `D:\plain\data.txt` + `)` 换成 `outOfProfile(t, home)`：
     那条字面 Windows 路径在 POSIX 上是个**相对路径**（`\` 在 POSIX 不是分隔符），
     "profile 之外不该嫌疑"这条判据在 ubuntu 上从没真正落到一个绝对路径上——两侧现在都是真平台形状；
  ③ 顺手查了 portable 文件里**有没有"两侧同绿但不同因"的字面 Windows 形状输入**：`D:\plain\data.txt` 那两处就是（已按 ② 换掉）；
     换完后 `internal/risk/` 里**再没有** portable 用例吃 `\\?\` / 8.3 / 盘符字面量当输入
     （那些判据本来就在 Windows 层 `syncdirs_redteam_windows_test.go`，`grep -rn "FinalPath\|D:.plain" internal/risk/*_test.go`
     只剩 ② 的那条历史注释 ⇒ **无需再搬**，这是一条负结果，写下来免得下一个人重找一遍）。
  **next=** Linux 侧 AC#2b/AC#3(i)/AC#1b/AC#4 数字（docker `golang:1.27` + `git archive HEAD` 仓外快照），取完回来勾框。

