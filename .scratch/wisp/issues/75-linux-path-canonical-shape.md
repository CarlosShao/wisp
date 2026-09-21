# 75 — Path canonicalization emits Windows-shaped (backslash) paths on Linux, so `internal/tools` is 19 FAIL + a 600 s timeout in `test-core`

**Status:** in-progress（票 72 已落地 `f1033e1` ⇒ 派发条件满足；开工即测量）
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
- [ ] **AC#1** Reproduce in docker with the CI command verbatim (not a local `go test`): record the
  exact command + `47`-ish baseline FAIL count in the evidence file **before** changing anything.
- [ ] **AC#2** `internal/tools` goes green in `golang:1.27` **and** stays green on Windows
  (`-count=2`). Both sides, or the box stays unticked — a fix that only moves the failure is not a fix.
- [ ] **AC#3** The 600 s timeout is explained by name (which test, which wait), not just "it got faster".
- [ ] **AC#4** ⚠ **D22 gate**: if the correct fix lands in `internal/risk/pathresolver*.go` or
  `assessor.go`, **stop and hand me the one-paragraph diff proposal instead of editing** — those files
  are the frozen security surface, and A38/Q-17 already has an unrelated change queued behind that gate.
- [ ] **AC#5** No assertion is weakened, no test is build-tagged away to make `test-core` green.
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
