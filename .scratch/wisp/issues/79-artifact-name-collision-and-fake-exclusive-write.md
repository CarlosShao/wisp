# 79 — `artifactName` folds distinct tool-call ids onto one disk name, and `writeFileExclusive` isn't exclusive

**Status:** implemented — AC#1..AC#5 ticked, commands + printed numbers in the Progress log;
awaiting adversarial acceptance (file deliberately NOT renamed `-done`: that suffix is the
orchestrator's re-claim key)
**Type:** correctness/data-integrity defect (artifact clobbering) + a storage-hygiene gap
**Blocks:** nothing · **Blocked by:** nothing (`internal/agent`/`internal/memory` are free once ticket 76 landed)
**Packages:** `internal/agent/spill.go`, `internal/memory/artifacts.go` + tests. Do **not** touch
`internal/risk` (ticket 75 in flight), `internal/ball`, `cmd/wisp`, `internal/proc` (ticket 78 in flight).
**Found by:** ticket 76's agent, reported and **deliberately not fixed** (correct call — out of that
ticket's scope). Spec refs: C25 provenance/taint, D31 atomic-rename, D6 single source of truth.

## Defect 1 (the real one): two different artifacts can silently overwrite each other

Measured by ticket 76 and re-read by me at `internal/agent/spill.go:132`:
- `artifactName(callID, seq)` keeps only `[A-Za-z0-9_-]`, so **`p/q`, `p\q` and `pq` all collapse to
  `tool-output-pq.txt`**. Those are *different* model-supplied tool-call ids.
- The doc comment admits only "same call id retried"; the sanitizer's collision is wider than that.
- Worse, the function named **`writeFileExclusive` does not use `O_EXCL`** — it's `os.WriteFile`
  (`O_CREATE|O_TRUNC`), so the name is a lie and the second writer silently wins.
- Consequence: the earlier `Spill.Path` (handed to the model / usable via `fs.read`) can point at
  bytes that were replaced by a different call's output. **That is a provenance break, not cosmetics**:
  a tool result the model believes it can re-read may be another tool's output.

### What to build
- Make the on-disk name **injective** with respect to the logical id (hash or escape rather than strip —
  stripping is what merges them), or detect the collision and fail loudly. Say which you chose and why.
- Make `writeFileExclusive` actually exclusive (`O_CREATE|O_EXCL`), or rename it to what it does.
  ⚠ If you make it truly exclusive, decide and **test** what happens on retry of the same id
  (the current documented behavior is that retry overwrites) — don't let that become an error path
  that breaks the caller.
- ⚠ **Do not "fix" the collision by allowing caller-controlled paths through** (that's the invariant
  ticket 76 just pinned: no caller-controlled destination path; see its AC#1/AC#2 tests).

## Defect 2: nobody enforces that `artifacts/` is flat

`listArtifactsDir` does `if e.IsDir() { continue }`, so a stray subdirectory is **invisible** to
`ListArtifacts`, `PurgeArtifacts`, and the 500 MB LRU quota. Ticket 76's fixture proved it: a `nested\`
canary directory **survived a purge** (existing behavior, recorded honestly, not caused by them).

### What to build
Pick one and write the reason in the commit message: (a) the quota/purge walks recursively so a stray
dir can't hide bytes from the cap; or (b) the store treats a non-conforming entry as an error it reports
rather than silently ignoring. **"Silently ignore" is the status quo and it is the option that must not win
by default** — an entry invisible to the quota is how a 500 MB cap becomes 2 GB.

## AC (1:1 verdict table, one row per box)
- [x] **AC#1** A test proving two *different* ids that previously collided now produce *different* on-disk
  names (or a loud failure), with the collision pair from the report (`p/q` vs `pq`) as named subtests.
- [x] **AC#2** `writeFileExclusive` is genuinely exclusive (or renamed + behavior documented), **plus** a test
  pinning what a same-id retry does now (both outcomes are acceptable; *unspecified* is not).
- [x] **AC#3** Mutation: revert the injectivity change ⇒ AC#1 must go red; **grep-prove the mutation landed
  before running and prove the revert after** (this repo has produced several false greens from
  no-op mutations today — anchor on the line that carries the thing you're changing, not on a symbol name).
- [x] **AC#4** Defect 2 resolved one of the two ways above, with a test where a stray subdirectory either
  counts against the quota or is reported — and that test must go red under the current `continue`.
- [x] **AC#5** Gates: `gofmt -l` on touched pkgs empty, `go vet` scoped rc=0, `go test -count=2` scoped, and
  the RUN-count invariant (`=== RUN` == 2 × distinct names) with **the 2 known pre-existing `SKIP`s in
  `internal/memory/concurrent_test.go:186` named explicitly** rather than filtered out of the report.

## Rules
`git commit -q -F - -- <explicit paths> <<'MSGEOF'` (quoted heredoc, A31); no `git add -A`; check
`git diff --cached --name-only`; no `--amend`/`reset`/`rebase`/`stash` (A34); no whole-repo gates
(shared tree); no worktree inside the repo (A38④); don't touch ticket 76's or ticket 20's faces
(those are mine). First checkpoint commit within 15 tool calls; sync Status + boxes + `next=` every commit.

## Progress log (append-only, newest last)

- 2026-09-21（实现代理，checkpoint 1）：只改 Status，形状三选一定案（注入式命名 / 真 O_EXCL +
  同 id 重走原子替换 / Defect 2 取选项 (a)），并把"票 76 的字面量必须跟着动"这件事写在脸上。
- 2026-09-21（实现代理，checkpoint 2）：**Defect 1 + Defect 2 双双落地**，两包 `go test -count=1` 全绿。
  - **命名形状：escape，不是 hash，不是"检测冲突就报错"**。`artifactName` 现在把 `[a-z0-9_-]`
    以外的每个字节 percent-escape（大写十六进制）。选 escape 的理由：它**信息论上就是 1:1**，
    不依赖摘要抗碰撞；`TestArtifactNameRoundTripsToTheExactID` 里带一个测试自己的解码器，
    26 个 id（含 `\x00`、`%2F`、`tool-output-x.txt`、尾随空格）逐个 decode 回原 id——
    "注入"是被证明的，不是被打印出来的。**大写也 escape** 是票面没写但同一个缺陷的另一顶帽子：
    artifacts 在 NTFS 上，NTFS 折叠大小写，留着 `PQ` 原样进名字就等于 `pq`/`PQ` 仍然同一个文件。
    `.` 也 escape 顺带把票 76 在删除侧修的"尾部点别名"在写入侧一起堵了（名字里不再有字面点）。
    超长的切在 escape 边界外 + 追 sha256(原始 id) 前 8 字节：不加上限的话 3 倍膨胀会把
    Windows 255 字节组件上限撞穿，那是把"静默覆盖"换成"直接写失败"。空 id 的退例 token 从
    `seq<N>` 改成 `seq.<N>`——字面点不在编码器值域里，所以退例名也不会再和某个字面 id 撞。
  - **`writeFileExclusive` 现在是 `O_CREATE|O_EXCL`**，名字不再是谎话；写失败会删掉半截文件。
    同 id 重试的行为**定死并测死**：`Prepare` 收到 `fs.ErrExist` 时按"同一个逻辑 id 的自己人"处理
    （名字已经注入式， Occupied 只剩这一种解释），走 temp+rename 的 D31 原子替换，last-writer-wins，
    不报错、不留 temp、重启后同 id 也算重试（`TestSpillAcrossRestartsKeepsRetrySemantics`）。
  - **Defect 2 取 (a)**：`listArtifactsDir` 递归量游离子目录（WalkDir 用 Lstat，不跟符号链接），
    配额把游离子树和规矩文件放进**同一条 LRU 队列**比 mtime，purge 也回收它。
    移除只用 `os.Remove`（不自上而下递归删的 `RemoveAll`），先文件后目录、子先于父。
    选 (a) 不选 (b) 的判决理由写在下一枚 commit 的正文里。
  - **动了票 76 的字面量（如实报告，共 3 处）**：`spill_path_invariant_test.go` 4 个 `wantName`
    （那 4 个字符串**就是**本票要消掉的折叠结果）+ 那个 seq 退例子测试（改名的同时把它**加强**成
    "5 个 id 必须得 5 个名字"——旧代码里那 4 个 id 全折进 `tool-output-seq1.txt` 互相覆盖，
    测试却在为这个结果鼓掌）；`artifacts_path_invariant_test.go` 的 (e) 段：
    "nested 目录必须活过 purge" 按本票判据反了，改成"必须不活过"并把 removed 集合**逐个名字**钉死，
    containment 从"只能是 artifacts 的直接子项"改成"artifacts 树内任意深度 + 集合相等"。
    守卫拒绝的四形状、AST sink 审计、containment 差分、阳性对照**一条没动**，全绿。
  - **新发现、本票不修**：(1) `os.OpenFile(..., 0o600)` 在 Windows 上是装饰——落盘是
    `-rw-rw-rw-`，权限由目录 ACL 继承，所以 artifact 的"只有我可读"从来没成立过（旧代码同样）；
    (2) 名字里现在可能有 `%`，`cmd.exe` 会展开 `%XX%` 形态——模型拿到的是 `fs.read` 路径不过 shell，
    但任何"把 artifact 路径塞进 shell 命令"的新路由都要先想这条；(3) `tool-output-*.txt` 的 8.3 短名
    在前 6 个字符就分叉，长文件名关闭短名时无关，开着则是同一个折叠类（票 20/18 的 junction 证据区）。
  - AC#1/AC#2 勾上；**AC#3/AC#4 的变异检验和 AC#5 门禁还没跑 ⇒ 那三框仍空**。
    `go test -count=1 ./internal/agent/ ./internal/memory/` = ok 1.683s / ok 12.971s。
    **next=做 AC#3（把 `encodeArtifactID` 换回 strip，grep 证明落地，看 AC#1 变红，还原）和
    AC#4（把 `if e.IsDir()` 的 `continue` 装回去，看我的游离子树用例变红，还原），然后跑 AC#5 三门禁 + RUN 计数不变式。**
- 2026-09-21（实现代理，checkpoint 3 = 变异与门禁，AC#3/AC#4/AC#5 全绿）：
  - **AC#3 变异（红→还原→绿）**。锚点不是符号名，是 `encodeArtifactID` 循环体里**真正承载注入性的那三行**
    （`b.WriteByte('%')` / `upperHex[c>>4]` / `upperHex[c&0xF]`），整段换成 `continue`（= 回到"剥掉"）。
    一条 `&&` 链里先 grep 证明落地再跑：`grep -n "MUTATION-AC3\|b.WriteByte(upperHex" spill.go`
    → 只命中 214 行的 MUTATION-AC3，`WriteByte(upperHex` **0 命中**（改真了，不是空跑）。
    `go test -count=1 -run 'TestArtifactName…|TestSpilledBytes…' ./internal/agent/` 打印：
    `TestArtifactNameDoesNotFoldDistinctIDs` 6 个命名子测试全 FAIL，错误文本就是报告里那对——
    `ids "p/q" and "pq" both name "tool-output-pq.txt"`；`TestSpilledBytesSurviveANameThatUsedToCollide`
    3/3 子测试 FAIL；round-trip 那测 5 条 `does not decode back to a literal id`。
    还原：`grep -c MUTATION-AC3`=0、`grep -c "b.WriteByte(upperHex\[c>>4\])"`=1、
    `git diff --quiet -- internal/agent/spill.go` 干净（与 HEAD 逐字节相同），重跑 `ok 0.189s`。
  - **AC#4 变异**：把 `if e.IsDir() { continue }` 原样装回 `listArtifactsDir`（同样先 grep 证明命中 96 行）。
    `go test -count=1 -run 'TestStray|TestPurgeArtifactsReclaims' ./internal/memory/` → 4/4 FAIL，数字是
    **`artifacts tree still holds 4106 bytes against a 20-byte quota (freed=0, files=0)`** 与
    **`after a purge the artifacts tree still holds 1024 bytes in 3 entries`**——就是"500MB 变 2GB"的缩小版实拍。
    还原后 `grep -c MUTATION-AC4`=0、`git diff --quiet` 干净、重跑 `ok 0.253s`。
  - **AC#5 门禁**（范围只有我的两个包，共树不跑整仓）：
    `gofmt -l internal/agent/ internal/memory/` → **0 行**；
    `go vet ./internal/agent/ ./internal/memory/` → **rc=0**；
    `go test -count=2 -v ./internal/agent/ ./internal/memory/` → rc=0，`ok agent 3.208s` / `ok memory 25.356s`，
    `--- FAIL` 行数 **0**；RUN 计数不变式 **`=== RUN` 274 = 2 × 137 distinct**，逐项对上。
    **SKIP 如实点名**：`--- SKIP` 共 **2 行**，两条都是 `TestSubprocessCrashWriter`
    （`internal/memory/concurrent_test.go:186`，`-count=2` 所以同一个名字出现两次）。
    同文件 `:222` 那条 `t.Skip` 是父测试的"我此刻是孩子"守卫，本轮 `WISP_CRASH_CHILD` 未设，
    `TestCrashRecoveryKillMidWrite` **真跑了没跳**——所以"2 个已知 SKIP"量到的是 2 行 / 1 个名字，
    不是 2 个不同测试；本票没有新增任何 SKIP，也没有新增断言弱化。
  - **动过的票 76 面：0 个文件改内容之外的东西**；三处字面量的理由与加强全部写在
    commit `6901aa9` 的正文里（那四个 `wantName` 字符串本身就是本票要消掉的折叠结果）。
    next=clean-HEAD 复测（`git archive HEAD | tar -x -C /tmp/…`，在仓外），确认这枚 HEAD 单独可绿，
    然后把 run id / 数字交给验收。
- 2026-09-21（实现代理，checkpoint 4 = clean-HEAD 复测）：**HEAD `878d9f8` 单独自洽**。
  `git archive HEAD | tar -x -C /tmp/wisp79-head`（在仓外，不建 worktree，A38④），
  在那份快照里 `go vet ./internal/agent/ ./internal/memory/` rc=0，
  `go test -count=2 ./internal/agent/ ./internal/memory/` = **ok agent 2.904s / ok memory 25.683s**。
  ⇒ 我的两枚 commit 不依赖任何人在飞的未提交改动。
  **本票范围外、没动、要有人接手才知道的三件事**：
  (1) `docs/specs/SPEC-05-agent-core.md:115` 与 `SPEC-02-data-storage.md:180` 还写着
  `artifacts\tool-output-<id>.txt`，磁盘上现在是 `<encoded id>`（票 20 面 `:107` 那句
  "只留 `[A-Za-z0-9_-]`，整串剥光退回 `tool-output-seq<N>.txt`"也已被本票作废——那是编排者的面，我没改）；
  (2) `0o600` 在 Windows 上不落地（实测 `-rw-rw-rw-`），要真"仅 owner 可读"得走 ACL，属另一票；
  (3) 名字里可能有 `%`，任何把 artifact 路径拼进 cmd.exe 的新路由要先过一遍这条。
  next=交验收。判据复跑命令全部在上面三条 Progress 记录里，可直接复制。
