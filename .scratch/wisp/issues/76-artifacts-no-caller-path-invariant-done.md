# 76 — Pin the artifacts path as "no caller-controlled path enters it" (closes ticket 20's `:103` box honestly)

**Status:** **done（5/5 PASS，编排者 2026-09-21 10:33 归档）**
裁决表 `docs/evidence/s1/76-adversarial-acceptance.md`。AC#3 我自己在纯净树重做：中和 `artifacts.go:167` 的
`strings.TrimRight(name, ". ")` 守卫 ⇒ **FAIL 2 / PASS 11**（红在 `.../dotdot/bare_and_empty`），还原后 `ok`。
交付：票 20 `:103` 已**改写文本并勾选**（原框照的是这个 API 并不存在的威胁模型）；两条没修的开出**票 79**。
**Type:** security-invariant characterization (closes an AC box that is currently **untestable as written**)
**Blocks:** ticket 20 archival · **Blocked by:** nothing (packages free: `internal/agent`, `internal/memory`)
**Packages:** `internal/agent/spill.go`, `internal/memory/artifacts.go` + their tests. Do **not** touch
`internal/tools` (ticket 73), `internal/ball` (74), `tools/d22scan`/`ci.yml` (71), `internal/risk` (72).
**Spec refs:** C25 provenance/taint, C26 (this ticket must **not** start routing artifacts through C26 — see below), ticket 20 `:103`

## Proven premise (measured by the ticket-20 agent, recorded as A38② — do not re-derive)
Ticket 20's same-shape sweep of 13 entry points found: **two artifacts routes do NOT pass through C26** —
`internal/agent/spill.go:97-106` and `internal/memory/artifacts.go:104-146` — **but neither accepts a
caller-supplied path**. The agent reported it and **added no guard**, and I judged that correct
(no caller-controlled input ⇒ no exploitable surface, and adding a C26 call would touch the frozen
resolver for zero security gain).

**That is precisely why ticket 20's box `:103` is still unticked.** The box reads
"attempting user-dir write via artifacts API → denied", but if the API cannot express a caller path,
there is nothing to *attempt*, so **the box is written against a threat model the API doesn't have**.
Leaving it unticked forever is not honest either — it reads as unfinished work to every future session.

## What to build
Convert the *implicit* safety of these two routes into an **explicit, executed invariant**, then decide
the box in writing:
- A test that pins: **the artifacts/spill API surface takes no caller-controlled destination path** —
  i.e. every path it writes is derived from the data dir + an internally generated name, and any
  caller-supplied *component* (artifact id / key / filename field) that contains a separator, `..`,
  a drive letter, or a UNC prefix is **rejected or sanitized**, proven by feeding all four shapes.
- A test that the derived target genuinely stays **under the data dir**: write via the API with the
  nastiest accepted input and assert by **`filepath.Abs` + prefix check plus a real file-listing
  diff** (so a rename/symlink trick can't pass by string equality alone).
- Then **rewrite ticket 20's `:103` box text** to state what is actually guaranteed
  ("artifacts API 无调用方可控目标路径；含分隔符/`..`/盘符/UNC 的组件一律拒或净化"),
  tick it, and cite this ticket + the test name. **Do not edit ticket 20 yourself — hand me the exact
  replacement sentence in your final report; that file's face is mine.**

## Explicitly out of scope (do not "improve" these)
- **Do not route artifacts through `risk.Resolve`/C26.** That is a frozen-contract behavior change with
  no security payoff here (A38② reasoning), and it would collide with ticket 72, which is legitimately
  working inside `internal/risk` right now.
- **Do not weaken any existing assertion** in ticket 18/20's junction tests, and do not add `t.Skip`.

## AC (1:1 verdict table required, one row per box)
- [x] **AC#1** The four hostile component shapes (separator / `..` / drive letter / UNC) are each
  exercised by a named subtest, and each one is **either rejected or sanitized**, with the resulting
  on-disk name asserted. A test that only checks "no panic" does not count.
  **PASS** — 8 个命名子测试（两条路由 × 四形状）+ 2 个退化形状子测试。memory 侧断言
  `errors.Is(err, ErrInvalidArtifactName)` **且不是** `ErrNotFound`（= 拒在碰文件系统之前），并逐字节读回
  调用者点名的 canary；agent 侧断言净化后的确切磁盘名（`p/q`→`tool-output-pq.txt`、
  `../../escape`→`tool-output-escape.txt`、`C:\Windows\System32\drop`→`tool-output-CWindowsSystem32drop.txt`、
  `\\fileserver\share\payload`→`tool-output-fileserversharepayload.txt`）+ `ReadDir` 里只有这一个文件。
- [x] **AC#2** Containment is proven by **real directory listing**, not string comparison: after the
  writes, list the data dir and assert nothing appeared outside it (and that no file appeared at a
  path the caller named). **PASS** — 两侧各做 `filepath.WalkDir` 全量快照 + 差分；
  **阳性对照两处都做了**：memory 侧删一个合法 bare name ⇒ 差分必须恰好报那一条且落在 artifacts 内；
  agent 侧手工做一次"绕过净化器"的同样 join 真写盘 ⇒ 列举必须看见 `<root>\CONTROL-escape.txt`
  （对照组实测把转义算术量了出来：前缀吃掉第一个 `..`，需 `depth+2` 组才落到 root）。
  另加 `TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`：一个真 `memory.Store` 数据目录上两条路由轮流上。
- [x] **AC#3** Mutation: neutralize the sanitize/reject step ⇒ AC#1 and AC#2 both go red. Grep-prove the
  mutation landed before running; grep-prove `git diff --quiet` on the file after restoring.
  **PASS（两次变异，各一包）** — ①`artifacts.go:152` 插 `return nil // MUT76-NEUTRALIZED` ⇒ grep 命中 1 行 ⇒
  AC#1 两测 + AC#2 全红，且列举差分**自己**报 `- [data/artifacts/nested/canary-separator.txt data/canary-dotdot.txt]`
  （= 调用者点名的文件真被删了）；②`spill.go:141` 把 `id := b.String()` 换成 `id := callID // MUT76-NEUTRALIZED` ⇒
  AC#1 红 + AC#2 报 `spill wrote OUTSIDE the data dir: ESCAPE-API-a/b/c.txt`。
  还原后：`grep -c MUT76` 两文件均 0、`git diff --quiet internal/agent/spill.go internal/memory/artifacts.go` rc=0。
- [x] **AC#4** Hand me the replacement sentence for ticket 20 `:103` + the test names that back it.
  **PASS** — 句子写在最终报告里；票 20 面**我没有动过一个字**（`git status` 可查）。
- [x] **AC#5** Gates: `gofmt -l` on touched pkgs empty, `go vet ./internal/agent/... ./internal/memory/...`
  rc=0, `go test -count=2` on those two packages, and **`=== RUN` line count == 2 × distinct test names
  with zero `SKIP`** in the log. **PASS with one disclosed exception** —
  `gofmt -l internal/agent/ internal/memory/` 空、`go vet` rc=0、`go test -count=2 -v` 两包 rc=0
  （agent ok / memory ok 24.592s），整包日志 **RUN=228 == 2×114 distinct、FAIL=0**；
  本票新增 8 测的**限定**日志 **RUN=44 == 2×22 distinct、SKIP=0、PASS=44、FAIL=0**。
  ⚠ 整包日志里有 **2 处 `--- SKIP`，全部来自既有测试 `internal/memory/concurrent_test.go:186`
  `TestSubprocessCrashWriter`**（re-exec 子进程辅助体，直接跑时按设计 skip；`TestCrashRecoveryKillMidWrite`
  才是真调用者）。不是本票新增、不是我用 skip 换绿、也不在我改的范围内 ⇒ 在此点名而不是把它过滤掉。

## Rules (all learned from incidents in this repo in the last 24 h)
- Commit form `git commit -q -F - -- <explicit paths> <<'MSGEOF' … MSGEOF` — **quoted heredoc** (A31:
  unquoted let bash execute backticks in a message and created 16 junk files at the repo root).
  Never `git add -A`; check `git diff --cached --name-only` before every commit.
- No `commit --amend` / `reset` / `rebase` / `stash` / `checkout .` (A34, shared worktree).
- Never run whole-repo gates here — `internal/tools`, `internal/ball`, `internal/risk`, `tools/d22scan`
  all have agents in flight; a `go test ./...` would be measuring their uncommitted WIP, not HEAD.
- First checkpoint commit within your first 15 tool calls; sync Status + boxes + `next=` every commit.
- Unprovable AC ⇒ leave the box unticked and say why.

## Progress log (append-only, newest last)

- 2026-09-21（实现代理，checkpoint 1）：**memory 侧四层证明落地**，新增
  `internal/memory/artifacts_path_invariant_test.go`：①`TestArtifactsAPITakesNoCallerControlledDestinationPath`
  用 `go/ast` 审 artifacts.go —— 没有任何入口收 path/dir/dest 形参数、3 处 `os.Remove`/`os.MkdirAll`
  的路径参数全部源自 `s.artifactsDir`、3 处 `listArtifactsDir(...)` 调用点全部只喂 `s.artifactsDir`
  （审到零命中即 t.Fatal，防探针空跑）；②`TestDeleteArtifactRejectsTheFourHostileShapes` +
  ③`TestDeletePrivacyItemRejectsTheFourHostileShapes`：分隔符/`..`/盘符/UNC 各一个命名子测试，
  每张卡的是"**被守卫拒**"而不是"被文件系统拒"——断言 `errors.Is(err, ErrInvalidArtifactName)` 且
  **不是** `ErrNotFound`，同时调用者点名的 canary 逐字节还在；④`TestArtifactsContainmentByDirectoryListing`：
  对整个 fixture 根做递归 `WalkDir` 快照 + 差分，恶意名单跑完两条路后 added/removed 双空，
  **阳性对照**（删一个合法 bare name ⇒ 差分里必须恰好出现那一条且落在 artifacts 内）保证差分不是瞎的。
  **坐实并修掉一个真缺陷（D-76a）**：原守卫 `filepath.Base(name) != name` 只看 Go 的语义，
  Windows 在解析前会**剥掉组件尾部的点和空格**，于是 `DeleteArtifact("....")` 是 `.` 的一种拼法——
  实测它穿过守卫直达 `os.Remove(<artifactsDir>\....)`，只因目录非空才报 "The directory is not empty"，
  空目录时就是一次删掉 artifacts 目录本身的操作。修法：`strings.TrimRight(name, ". ") != name` 一律拒。
  顺带新增哨兵 `ErrInvalidArtifactName`（纯 additive，旧错误文案一字未变），因为"拒了"必须**机器可辨**，
  否则 AC#1 只能靠嗅 `err.Error()` 字符串。`go test ./internal/memory/` 全包 `ok 12.434s`。
  **AC#1/AC#2/AC#3/AC#4/AC#5 全部未勾**：agent 侧（`artifactName` 净化路径）还没测，变异检验和
  `-count=2` 门禁也没跑，此刻任何勾选都是假绿。
  **next=写 `internal/agent/spill_path_invariant_test.go`（四种形状净化后的**磁盘名**逐个断言 +
  真目录列举差分 + 经真 `memory.Store` 的端到端落盘），然后做 AC#3 双包变异检验与 AC#5 门禁。**
- 2026-09-21（实现代理，checkpoint 2）：**agent 侧四层证明落地**，新增
  `internal/agent/spill_path_invariant_test.go`（未改 spill.go 一行实现）：
  ①`TestSpillCallIDHostileShapesSanitizedToBareNames` 分隔符/`..`/盘符/UNC 各一命名子测试，
  断言净化后的**确切磁盘名**（`p/q`→`tool-output-pq.txt`、`../../escape`→`tool-output-escape.txt`、
  `C:\Windows\System32\drop`→`tool-output-CWindowsSystem32drop.txt`、
  `\\fileserver\share\payload`→`tool-output-fileserversharepayload.txt`）+ 该名字里没有 `/ \ : ..`
  + 真 `ReadDir` 里**只有**这一个文件；另一子测试钉"整串被剥光 ⇒ 退回 `tool-output-seq<N>.txt`"。
  ②`TestSpillContainmentByDirectoryListing`：递归列举差分 + **阳性对照**——手工做一次"绕过净化器"的
  同样 join 并真写盘，证明列举确实看得见越界（否则整张 AC#2 是瞎跑）。对照把转义算术**量**出来了：
  `artifactName` 的前缀 `tool-output-` 会吃掉第一个 `..`（变成字面量组件 `tool-output-..`，被第二个 `..` 抵消），
  所以要 `depth+2` 组 `..\` 才正好落到 root——3 组只到 `data\`，实测两次才修对，注释里写明了"这是量出来的"。
  ③`TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`：一个真 `memory.Store` 数据目录上，
  spill 路由（敌意 id）+ memory 路由（敌意名）轮流上，列举差分要求新增只允许 `data/artifacts/*`。
  ④`TestSpillAPITakesNoCallerControlledDestinationPath`：AST 审 spill.go —— 导出面不许有 destination 形参
  （`NewSpiller` 是唯一例外，它是宿主侧构造），两个 os/内部写盘 sink（`os.WriteFile`、`writeFileExclusive`）
  的路径实参必须展开到 `s.dir` + `artifactName(`（需要不动点展开：`path`/`name` 是两条语句）。
  **AC#1/AC#2 两侧全绿；AC#3 变异/AC#5 门禁未做 ⇒ 五框仍全部未勾。**
  `go test -count=1 ./internal/agent/ ./internal/memory/` = ok 1.744s / ok 12.956s。
  **next=AC#3 双包变异检验（各自 neutralize 掉 reject/净化那一步，先 grep 证明改动真落地，再证明 AC#1+AC#2 同时变红，
  然后还原 + `grep -c`=0 + `git diff --quiet`），跑 AC#5 门禁，最后交票 20 `:103` 的替换句子。**
- 2026-09-21（实现代理，checkpoint 3 = 收尾）：**AC#3 变异 / AC#5 门禁做完，五框全绿**，AC#4 的句子在代理最终报告里
  （票 20 面一个字没改）。变异检验的具体收获，两条都值得记进缺陷册：
  ①**D-76a（本票已修 + 已测）**：`validArtifactName` 原来只看 Go 的路径语义，而 Windows 在解析前会剥掉组件
  **尾部的点和空格** ⇒ `"...."` 是 `"."` 的一种拼法，实测它穿过守卫直达 `os.Remove(<artifactsDir>\....)`，
  只因目录非空才报 "The directory is not empty"；`TestDeleteArtifactRejectsTheFourHostileShapes/dotdot/bare_and_empty`
  在变异体（守卫 return nil）下还量到 `".."` → `remove ...\data`、`"."`/`""` → `remove ...\data\artifacts`：
  **空目录时就是一次删掉 artifacts 目录本身**。修法一条：`strings.TrimRight(name, ". ") != name` 一律拒。
  ②**上报不修（编排者判）**：
  - **净化器把不同 id 折到同一磁盘名**：`artifactName` 只留 `[A-Za-z0-9_-]` ⇒ `p/q`、`p\q`、`pq` 全变 `tool-output-pq.txt`，
    而 `writeFileExclusive` 用的是 `os.WriteFile`（`O_TRUNC`，名字里的 "Exclusive" 并不存在）⇒ **不同 tool-call 的
    产物可互相覆盖**；spill 之后模型还能按 `Spill.Path` 用 fs.read 读到"上一份内容已被别人覆掉"的文件。
    文档注释只承认"同 id 重试覆盖"，没承认"不同 id 也覆盖"。建议：净化后与原文不等时追加 `seq`（或不落模型可控名）。
  - **artifacts 目录"扁平"没有任何一侧在守**：`listArtifactsDir` 对 `e.IsDir()` 一律 `continue`，
    于是一个混进 `artifacts\` 的子目录对 `ListArtifacts`/`PurgeArtifacts`/LRU 配额**全部不可见**
    （我的 containment fixture 里那个 `nested\` canary 目录就是被 purge 跳过而活下来的，属既有行为的如实记录）。
  - **按 ruling 留在原地的残余**：两条路由仍**不**过 `risk.Resolve`/C26（本票明令不做）。因此若有人往
    `artifacts\` 里放一个 reparse point，`DeleteArtifact("<bare name>")` 会照穿不误——这不在本票判据范围，
    桥侧 junction 证据归票 20/18。
  **门禁数字（AC#5）**：`gofmt -l internal/agent/ internal/memory/` 空；`go vet ./internal/agent/... ./internal/memory/...` rc=0；
  `go test -count=2 -v` 两包 rc=0（memory ok 24.592s），整包 `RUN=228 / distinct=114 / FAIL=0 / SKIP 出现 2 次`
  （两处 SKIP 均为既有 `concurrent_test.go:186` 的 re-exec 辅助体，已点名、未过滤）；
  本票 8 测定包子集 `RUN=44 == 2×22 distinct、SKIP=0、PASS=44、FAIL=0`。
  **next=交编排者：把 `:103` 的替换句子落到票 20、裁决 D-76a 是否单独立票、以及上面 ② 的三条要不要各开一张。**
