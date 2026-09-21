# 73 — Sweep orphan `.wisp-tmp-*` staging files left by a real kill

**Status:** in-flight (implementer) — 清扫器落地中，判据框未勾
**Type:** defect-fix (bookkeeping of a proven leak)
**Blocks:** nothing · **Blocked by:** nothing (the characterization test already exists)
**Spec refs:** D31 atomic-rename, SPEC-07 §2–§3, registry **A18**
**Packages:** `internal/tools` only (do not touch `internal/ball`, `internal/observe`, `tools/d22scan` — other agents are in flight there)

## Proven premise (do not re-derive, do not re-litigate)
`761447f` added `internal/tools/bridge_a18_kill_windows_test.go`: a **real child process** + **real
`taskkill /F`** mid-write. It proves three things, all already measured, not guessed:
1. The destination file is byte-for-byte intact or absent ⇒ **D31 holds under a real fault**.
2. **Each kill leaves exactly one `.wisp-tmp-*` behind** (4 bytes each in the test).
3. **Nothing sweeps it** — a later successful write to the same directory does not remove the orphan.

That test deliberately pins the assertion to "the orphan is still there". **This ticket flips that
one assertion to "the orphan is gone"** and must not weaken any other assertion in the file.

## What to build
Recovery of interrupted atomic writes at the place that already owns the target path, so a killed
session does not accumulate litter in the user's allowed dirs:
- On the next write (or on bridge start-up, your choice — say which in the commit message and why),
  reclaim `.wisp-tmp-*` files in the directory being written **that this bridge itself created**.
- The reclaimer must be **conservative**: never delete a `.wisp-tmp-*` it cannot attribute to itself
  (attribution mechanism is your design — record the naming scheme or a marker; explain the choice),
  and never follow a junction/reparse point or a symlink while sweeping (that would turn a cleanup
  into a delete primitive outside the allowed dirs — this repo's C26 exists precisely for that).
- A temp file that is **currently held by another live process** must survive; a retry-then-skip is
  acceptable, a hard failure is not.

## AC (1:1 verdict table required, one row per box)
- [ ] **AC#1** The flipped A18 assertion: after a real `taskkill /F` mid-write, a subsequent write
  through the bridge leaves the target directory with **zero** `.wisp-tmp-*` entries; the destination
  file is still byte-exact. Mutation check required: revert the reclaimer to no-op ⇒ this test goes
  red (grep-prove the mutation landed before running, restore and grep-prove it is gone).
- [ ] **AC#2** Attribution is real, not "delete every `.wisp-tmp-*` in the dir": a `*.wisp-tmp-*`
  file **the bridge did not create** (e.g. a foreign file with the same prefix, created directly by
  the test) survives the sweep, with a test that names the surviving file.
- [ ] **AC#3** Reparse-point safety: an orphan-shaped name reached through a **real NTFS junction**
  pointing outside the allowed dirs is not deleted. `mklink /J` for real, `t.Fatalf` (not `t.Skip`)
  if the fixture cannot be built, and a positive control proving the target file exists first.
- [ ] **AC#4** A live temp file held by another process is not deleted and does not fail the write:
  spawn a child that opens the temp and blocks, then write. `//go:build windows`, zero `t.Skip`.
- [ ] **AC#5** Gates: `gofmt -l internal/tools` empty, `go vet ./internal/tools/...` rc=0,
  `go test ./internal/tools/... -count=2`, `-race`, and the repo's `d22scan` invoked **exactly as
  the CI line invokes it** (`cd tools/d22scan && go run . -root <abs>`), with a seeded-violation
  self-test shown red at least once so a green scan means something.

## Explicit prohibitions for this ticket
- No `git commit --amend`, `git reset`, `git rebase`, `git stash`, `git checkout .` (shared worktree,
  registry A34). Commit with `git commit -F - -- <explicit paths>` and a **quoted** heredoc
  (`<<'MSGEOF'`) — an unquoted one made the shell execute backticks on `9bf05cb` (A31).
- After `git add`, before `git commit`, run `git diff --cached --name-only` and confirm only your
  paths are staged.
- Do not tick an AC you did not reproduce. Do not add `t.Skip`, do not lower a threshold, do not
  touch frozen contracts (`docs/PLAN.md`, `docs/specs/*`, `internal/risk/assessor.go`,
  `internal/risk/pathresolver*.go`, `rules_gateway.go`).
- If an AC cannot be proven, **leave the box empty and say why** in the Progress log.
- Commit your first checkpoint **within your first 15 tool calls** (turn cap is ~150; a ticket that
  dies with nothing committed is a total loss). Sync the ticket face (Status + boxes + one Progress
  log line ending in `next=`) at **every** commit.

## Progress log (append-only, newest last)

- 2026-09-21（实现代理，checkpoint 1）：开工前只做了设计拍板，代码未落，AC 框一个都没勾。
  **清扫时机 = 下一次写盘**（`stageAndRename`/`crossVolume` 在 `os.CreateTemp` 之前对目标目录扫一次），
  不是桥启动时——理由：本仓 `internal/tools` 里没有"桥启动"这个可挂钩的单点（`New()` 只是构造，
  `cmd/wisp` 不许我碰），而证据 §5 已经写明"下一次启动"在这张票的可达面上只能用"新建一套桥 + 一次成功写盘"代表，
  正好就是 A18 用例第三条在做的事。**归属机制 = 暂存文件名自带身份**：
  `.wisp-tmp-<owner8>-<pid>-<rand>`，owner8 = sha256(可执行文件路径 + 用户名 + tempPrefix) 前 8 位十六进制
  （跨进程稳定 ⇒ 新进程能认死进程的文件；不同程序/不同用户 ⇒ 不同 token，前缀相同也不会互删）；
  pid 用来判"创建者还活着吗"（活人持有的暂存文件不是孤儿）。扫到名字不合式、是 reparse point、
  不是普通文件、创建者 pid 还活着、或删不动（重试后仍失败）的一律**跳过不删**。
  **next=落 `internal/tools/fs_staging.go` + 平台 liveness 文件 + 翻 A18 第三条断言，然后按 AC#1 做变异检验**。
