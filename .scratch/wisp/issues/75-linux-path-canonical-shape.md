# 75 — Path canonicalization emits Windows-shaped (backslash) paths on Linux, so `internal/tools` is 19 FAIL + a 600 s timeout in `test-core`

**Status:** ready-for-agent, **but hold dispatch until ticket 73 lands** (same package `internal/tools`)
**Type:** portability defect (CI-blocking)
**Blocks:** ticket 70 AC#2 (`test-core`) · **Blocked by:** ticket 73 (package conflict, not semantics)
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
