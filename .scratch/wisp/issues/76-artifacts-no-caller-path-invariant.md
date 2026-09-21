# 76 — Pin the artifacts path as "no caller-controlled path enters it" (closes ticket 20's `:103` box honestly)

**Status:** in progress（实现代理 checkpoint 1：`internal/memory` 侧四层证明已绿 + 顺手坐实一个真缺陷；`internal/agent` 侧未写）
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
- [ ] **AC#1** The four hostile component shapes (separator / `..` / drive letter / UNC) are each
  exercised by a named subtest, and each one is **either rejected or sanitized**, with the resulting
  on-disk name asserted. A test that only checks "no panic" does not count.
- [ ] **AC#2** Containment is proven by **real directory listing**, not string comparison: after the
  writes, list the data dir and assert nothing appeared outside it (and that no file appeared at a
  path the caller named).
- [ ] **AC#3** Mutation: neutralize the sanitize/reject step ⇒ AC#1 and AC#2 both go red. Grep-prove the
  mutation landed before running; grep-prove `git diff --quiet` on the file after restoring.
- [ ] **AC#4** Hand me the replacement sentence for ticket 20 `:103` + the test names that back it.
- [ ] **AC#5** Gates: `gofmt -l` on touched pkgs empty, `go vet ./internal/agent/... ./internal/memory/...`
  rc=0, `go test -count=2` on those two packages, and **`=== RUN` line count == 2 × distinct test names
  with zero `SKIP`** in the log.

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
