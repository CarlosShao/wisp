# 75 — Path canonicalization emits Windows-shaped (backslash) paths on Linux, so `internal/tools` is 19 FAIL + a 600 s timeout in `test-core`

**Status:** in-progress（owner 指派的四项里 1/2/3 已完成并测量；剩 AC#4 由 owner 用交接段落关闭、AC#6 需 owner 推送后的 run 证据）
**Claimed by:** implementer agent（2026-09-21 10:2x）
**Evidence:** `docs/evidence/s1/75-rootcause-linux-path-shape.md`（根因 file:line + AC#1 docker 基线 + AC#4 提案）
**Type:** portability defect (CI-blocking)
**Blocks:** ticket 70 AC#2 (`test-core`)
**Blocked by（原文是"票 73"，那是错的）**：**票 72** 才是真冲突——它此刻正在 `internal/risk/` 里改
`pathresolver*.go`，而本票的正解**极可能就是同一个文件**（C26 的 canonical 在 Linux 产反斜杠）。
同包并行 = 假并行（会把对方的改动吞进我的 diff 或反过来）。票 73 已 `done`，`internal/tools` 现在反而是空的。
⇒ **派发条件 = 票 72 落地之后**。等待期不浪费：我另派了一个**只读**代理去做根因定位
（只查不改，产出"哪一行在 Linux 上产反斜杠"的证据表），这样本票开工时不必从冷搜索开始。
**Packages:** `internal/tools` (+ possibly `internal/risk/pathresolver*` = **frozen, D22 — see below**)
**Evidence:** ticket 70's run forensics (`24a66b6`), docker `golang:1.27` reproduction with the CI command verbatim

## Proven premise
`test-core` is red with **4 packages / 47 `--- FAIL`** (not "about 4" as the ticket-70 face said).
Two of those families are the same root cause:
- `internal/risk` 17 FAIL
- `internal/tools` 19 FAIL **+ a 600 s test-binary timeout panic**

The reported root cause: **C26 canonicalization still produces backslash-shaped paths on Linux.**
On Windows the shape is right, so every local gate is green — which is exactly why this survived.

## What to build
Make the canonical form **platform-shaped**, not "Windows-shaped everywhere", or (if the contract
demands one canonical shape on all platforms) make the *comparison* normalize both sides. Which of
the two is correct is a **contract question**: read C26 in `docs/specs/` first and write down which
reading you implemented and **why, in the commit message**.

## AC (1:1 verdict table required)
- [x] **AC#1** Reproduce in docker with the CI command verbatim (not a local `go test`): record the
  exact command + `47`-ish baseline FAIL count in the evidence file **before** changing anything.
- [x] **AC#2** `internal/tools` goes green in `golang:1.27` **and** stays green on Windows
  (`-count=2`). Both sides, or the box stays unticked — a fix that only moves the failure is not a fix.
- [x] **AC#3** The 600 s timeout is explained by name (which test, which wait), not just "it got faster".
- [ ] **AC#4** ⚠ **D22 gate**: if the correct fix lands in `internal/risk/pathresolver*.go` or
  `assessor.go`, **stop and hand me the one-paragraph diff proposal instead of editing** — those files
  are the frozen security surface, and A38/Q-17 already has an unrelated change queued behind that gate.
- [x] **AC#5** No assertion is weakened, no test is build-tagged away to make `test-core` green.
  (Ticket 70 legitimately used `//go:build windows` for **DPAPI**, because C28 says DPAPI is
  Windows-only and it wrote the coverage cost down. That justification does **not** transfer here:
  path shape is not a platform-API limitation, it is our bug.)
- [ ] **AC#6** Runner-visible proof: after landing, the next `dev` push's `test-core` job is quoted by
  run id + job id, with the FAIL count before/after in one line. **Cancelled runs do not count as evidence**
  (see A40① — three consecutive `dev` push runs had conclusion `cancelled`, which nobody has explained yet).

## Rules (non-negotiable, learned the expensive way today)
- Commit form: `git commit -q -F - -- <explicit paths> <<'MSGEOF'` — **quoted** heredoc (A31: an
  unquoted one let bash execute backticks and created 16 junk files at the repo root).
- No `commit --amend` / `reset` / `rebase` / `stash` / `checkout .` (A34, shared worktree).
- After a rename, verify exactly one face survives (A39: `git ls-tree --name-only HEAD <dir> | grep <n>`
  must return 1 line and `git diff --cached --name-only` must be empty).
- Do not "fix" a red by adding an allowlist entry or a `t.Skip`.
- If a widened gate's coverage now catches existing violations, **fix them in the same batch**
  (A40⑤: the emoji-gate widening left `HEAD` red in `lint` for ~13 min for exactly one `U+26A0`).
- First checkpoint commit within your first 15 tool calls; sync Status + boxes + `next=` every commit.

## Progress

- 2026-09-21 10:2x 开工。票 72 已 `f1033e1` 落地 ⇒ 派发条件满足。
  只读代理的 `docs/evidence/s1/75-rootcause-linux-path-shape.md` **不存在**，根因自己定位：
  `internal/risk/pathresolver.go:139` 的 `normalizeLocalUNC` 第一行无条件
  `strings.ReplaceAll(p, "/", "\\")`，配合 `pathresolver.go:62` 的 `lexCanonical`
  在 Linux 上把 `/home/u/x` 变成 `<cwd>/\home\u\x` —— 一个 OS 根本开不了的字符串。
  Windows 上 `filepath.Clean` 本来就把 `/` 折成 `\`，故该行在 Windows 恒等 = 所有本地门绿。
- 基线测量已在 docker `golang:1.27` 用逐字 CI 命令跑起来（`git archive HEAD` 干净快照，
  共树有票 76 在途代理，未跑任何整仓门）。
- next= 等 docker 基线跑完补齐 AC#1 三个读数 → 落修复（`pathresolver.go` 单独一个 commit
  以便 owner 一键退回"只交提案"的读数）→ Windows `-count=2` 半边门 → 勾框。
- 2026-09-21 接续代理（owner 指定四项）开工。根因与已落地面不重做（见上两条 + `27c6fe5`/`ed74595`/`493b4e7`）。
  **第 1 项 Windows 半边门：已真跑，`-count=2`，四个文件全绿，零 SKIP**（本 box 是 Windows 原生，非交叉编译）：
  - `internal/risk/pathresolver_junction_windows_test.go`（票 18 Case 1–10，10 条顶层）
    → `ok github.com/CarlosShao/wisp/internal/risk 1.640s`，**EXIT=0**，`=== RUN` **20**，顶层 `--- PASS` **20**。
  - `internal/tools/bridge_junction_windows_test.go`（票 20，6 条顶层）
    → `ok internal/tools 1.674s`，**EXIT=0**，`=== RUN` **28**（含子项），顶层 `--- PASS` **12**。
  - `internal/tools/fs_staging_windows_test.go`（票 73，4 条顶层）
    → `ok internal/tools 3.341s`，**EXIT=0**，`=== RUN` **8**，顶层 `--- PASS` **8**。
  - `internal/tools/bridge_a18_kill_windows_test.go`（票 73，1 条顶层）
    → `ok internal/tools 2.364s`，**EXIT=0**，`=== RUN` **8**，顶层 `--- PASS` **2**。
  ⇒ 报告"未证/局限"第一条（"Windows 那一半它没测 ⇒ 恒等只是推理"）**关闭**：形状修复在 Windows 上确实是恒等。
- **第 2 项 P4**：`ed74595` 已把 `internal/tools/paths.go` 的 `const sep` 换成 `pathSep = string(filepath.Separator)`
  + `unifySeparators` 平台分支，跨文件放行泄漏（`a\b` 与 `a/b` 折成同一 key）在**白名单根**这一侧已被
  `TestFoldPathKeepsPosixBackslashesDistinct` 钉住。**但同一折叠在 `internal/risk/blacklist.go:134 normPath`
  → `Gate`/`overrideApplies` 的 `bOverrides` 那一侧只被测了"两种 Windows 拼写不改变裁定"（且那条在
  `_windows_test.go` 层，Linux 上不跑）** ⇒ 本代理补一条不带 tag 的正反两用例（见下一个 commit）。
- **第 3 项**："靠这层脏才绿"的用例：`TestFourChannelExfilSuite/fs.write_into_sync_dir` 的 fixture 已由
  `ed74595` 换成 `t.TempDir()` 下的真实平台形状目录（跨平台那条腿已在）。本代理补的是另两半：
  windows 层的拼写不变式（含"根外目标不得被标"的反向钉，只有 Windows 判得动）与 `//go:build !windows`
  层的 POSIX 等价用例（含休眠的"普通写不嫌疑"腿）。原计划的"不带 tag 的反向钉"被实测否掉了 —— Linux 上那条
  断言现在必然红（票 55 的检测缺口），写死它等于替票 55 糊一层，故改为休眠形式；见下面落地段。
- **第 4 项 AC#4 保持不勾**（owner 裁决 R17：修实现去符合既有契约 ≠ D22 改契约）。交接段落已写进本文件
  Progress 末段，供 owner 亲自关闭。

### AC#4 交接段落（D22 gate；本代理不自行勾框，请 owner 用它关闭）

本票在冻结面 `internal/risk/pathresolver.go` 上只改了三处，且都是"把实现改得符合已冻结的契约"而不是"改契约"：
(1) `normalizeLocalUNC` 开头加守卫 `if !strings.HasPrefix(p, "\\") && !strings.HasPrefix(p, "//") { return p }`，
把该函数从"对每个路径无条件把 `/` 折成 `\` 并返回"收窄为"只对真是 UNC 拼写的串动手，其余原样交回"；
(2) 新增包内常量 `sepStr = string(filepath.Separator)` 与辅助 `unifySeparators`（折叠只发生在 Windows 分支），
供同包 `blacklist.go` / `syncdirs.go` 使用；
(3) `tailExistsBelow` 交给 `os.Lstat` 的那次重拼接从字面 `\` 改成 `sepStr`。
**行为差异**：Windows 上 C26 交回的 canonical 字符串逐字不变（`filepath.Clean` 两行之前已把 `/` 折成 `\`，
守卫分支根本取不到），POSIX 上从 `<cwd>/\tmp\x`（一条没有任何 OS 调用能打开、也不再绝对的串）变成 `/tmp/x`；
A/B 表内容、`Class` 取值、函数签名、调用点、`Resolve` 的 fail-closed 方向一字未动。
**授权它的契约线**：`docs/specs/SPEC-06*.md:50-52`（§4 管线：展开 env/~ → 绝对化 → Clean → 句柄真实路径 →
拒 reparse → 展开 8.3 → **规范化 UNC** —— 最后一步就是 `normalizeLocalUNC` 的全部职责，它的输入本就是 UNC 拼写）
+ `docs/PLAN.md:1376`（C26 唯一入口，"真实路径"）、`:1796`（D33/F1 推翻字面比较那一行）、`:2377`
（macOS 侧 `realpath` + `lstat`，即 canonical 在 POSIX 上必须是 `/` 形状）。
**什么情况下这次编辑会变成一个真的 D22 变更**：如果正解被选成反方向那一种 —— 即主张 C26 的 canonical 在
**所有平台上必须是同一个 `\` 形状单一串**，于是契约要新增一条"跨平台单一形状"断言、比较侧改成双侧折叠、
并在任何 OS 调用之前把串展开回去，同时改写 `SPEC-06 §4` 与票 18/72 已冻结的裁定不变式 —— 那需要 owner 拍板。
本票走的不是这条路，故 AC#4 的框留着不勾，由 owner 用本段落关闭。

- next= 本代理剩余三步：① 给 `bOverrides` 侧补一条不带 tag 的 POSIX 折叠不合并用例（P4 的残余漏口）→
  ② sync-write exfil 用例分层（windows 层字面量腿 + 不带 tag 的"外平台字面量不再一律嫌疑"反向钉）→
  ③ 重跑票 72 的六条不变式 + `GOOS=linux go vet` 两侧门，然后**由 owner 推送**后按 run id + job id 读
  `test-core` 补 AC#6（本代理不 push）。AC#6 框在此之前保持不勾。

- 2026-09-21 接续代理第 2/3 项落地完毕（下面三段是测量，不是计划）。
  **P4 状态**：`ed74595` 已经修掉实现（`internal/tools/paths.go` 的 `const sep` →
  `pathSep = string(filepath.Separator)` + `unifySeparators` 平台分支；`internal/risk/blacklist.go:134`
  的 `normPath` 同一处置），tools 侧由 `TestFoldPathKeepsPosixBackslashesDistinct` 钉住**白名单根**那一侧；
  **`bOverrides`（放行侧 map 查表）那一侧当时没有 POSIX 用例**，本代理补
  `internal/risk/pathshape_portable_test.go:TestOverrideKeyKeepsPosixBackslashesDistinct`
  （不带 build tag；`<home>/id_k/x.pem` 与同名但含反斜杠的 `<home>/id_k\x.pem` 两条都判 B、折叠键必须不同、
  给其中一条的 override 不得解锁另一条；Windows 分支断言"同一文件的两种分隔符拼写必须仍折成一个键"，
  所以它在两个平台上都不是空跑）。**变异对照**：(M1) 把 risk 的折叠改成无条件 ⇒ 该测试 Linux 红
  （docker `golang:1.27`，`--- FAIL`，RC=1）；(M1b) 把折叠在 Windows 侧关掉 ⇒ 同一条测试在 Windows 红。
  **exfil 用例分层**：windows 层新增 `internal/risk/provenance_syncdirs_windows_test.go:
  TestExfilSyncWriteWindowsSpellingInvariant`（7 子项：真实同步根内 4 种拼写必须全升 `ChSyncWrite`，
  根外 3 种拼写必须不升 —— 后者就是"不是一张一律捞起的网"那一半，只有在本机检测能走完时才判得动），
  POSIX 层新增 `internal/risk/provenance_syncdirs_other_test.go`（`//go:build !windows`）
  `TestExfilSyncWritePosixSpellingInvariant`（用 `t.TempDir()` 构造，含"文件名里带 `\` 仍在同步根内"一条；
  它的"普通写不该被标"反向腿**显式做成休眠断言**：`SyncDetectionComplete()` 为假时只 `t.Log`，票 55 落地即自动生效，
  因为今天那条腿在 Linux 上属于票 55 那 8 条红，本代理不糊）。
  **变异对照**：(M3) 令 `deepestExistingAncestor` 永远返回 ""（就是在 Windows 上模拟修好之前的 POSIX 形状）
  ⇒ windows 层那条测试以 3 条 `plain_dir` 子项转红，MUTANT_RC=1；恢复后 `grep -c MUTANT` = 0。
- 门（都是 scoped，共树票 79 仍在写 `internal/agent`、`internal/memory`，未跑整仓门）：
  `gofmt -l internal/risk internal/tools` 空；`go vet ./internal/risk/ ./internal/tools/` rc=0；
  `GOOS=linux go vet` 同两包 rc=0；Windows 全包 `go test ./internal/risk/ ./internal/tools/ -count=1`
  → `ok risk 4.853s` / `ok tools 14.097s`；**票 72 六条不变式重跑 `-count=2`：rc=0，`=== RUN` 12，
  顶层 PASS 12，FAIL 0，SKIP 0**（`TestClassifyAnchorSpellingIsNotVerdict`、`TestClassifySpellingInvariance`、
  `TestClassifyFailClosedWhenSpellingUnprovable`、`TestAListWinsWhereBothTablesHit`、
  `TestOverrideOnlyAcceptsTheResolvedForm`、`TestCanonicalInputGainsNoSecondForm`）。
  docker `golang:1.27`（`git archive HEAD` 快照 + 只覆盖本代理改动的三个文件）：
  `internal/tools` **ok 9.345s**（前后一致），`internal/risk` 红数 **8 → 8，零新增红**，
  那 8 条与基线逐字同名（票 55）；本代理两条 POSIX 新测试 `-count=2` 全绿。
- 顺手记一笔待 owner 裁决的新缺陷候选：`risk.Gate` 的 `bOverrides` 形参**今天没有任何生产调用方**
  （`internal/config/schema.go:457 BlacklistOverrides` 会被解析、`manager.go:356` 还会对它做松/紧方向审计，
  但没有人把它喂进 `Gate`）⇒ "单文件 B 档豁免"这条契约目前只存在于测试里。
- next= 交回 owner：① 推送后按 run id + job id 读 `test-core`（AC#6，step 结论而非 job 结论，
  `jobs.total_count=0` 不算样本）；② AC#4 用上面的交接段落关闭；③ 决定要不要把 `bOverrides` 无调用方
  登记成新票。
- 收尾复测（HEAD=`ea01b9f`，即包含本代理全部改动之后**重新跑**的同一批 Windows `-count=2`，与前面
  第一次的数字逐字一致，证明新加的两个测试文件没有动到别人）：
  risk junction `EXIT=0 / === RUN 20 / 顶层 PASS 20 / SKIP 0`；tools bridge junction
  `EXIT=0 / 28 / 12 / 0`；tools staging `EXIT=0 / 8 / 8 / 0`；tools a18 kill `EXIT=0 / 8 / 2 / 0`。
- next= **AC#6 请 owner 读的 run**：把 `dev` 推到 `HEAD=ea01b9f`（或含它的更新 sha）之后，
  `gh run list --branch dev --workflow ci --limit 3` 取 commit 对上的那条 run id，
  再 `gh run view <run-id> --json jobs --jq '.jobs[]|select(.name|test("test-core"))|{name,conclusion,steps:[.steps[]|{name,conclusion}]}'`
  ——**读 `Portable package tests` 这一步的 conclusion，不读 job 状态**（A44③/A49④），
  且该 run 的 `jobs.total_count=0`（cancelled）时**不算样本**（A40①）。要贴的回执一行是：
  `test-core @ <run-id>/<job-id>: --- FAIL 由 4 包 47 条（`24a66b6` 取证，票 70 面上写的"about 4"是错的）→ <N> 条；
  internal/tools 由 19+ 条 + 600s panic → 本代理 docker 实测 ok 9.345s；internal/risk 剩 8 条 = 票 55`。
  本代理不 push，AC#6 框保持不勾直到上面这行被 run id 填实。
