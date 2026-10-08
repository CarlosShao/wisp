# 00-anchor.md — 33-v4 (non-implementer acceptance leg) start gate

All numbers below are this leg's own measurements, taken with the rulers written in this file.
No number is copied from the writing leg's `impl.md` or from the orchestrator's relay.

## HEAD at start (measured, not "should be")

```
$ git log -1 --format=%H | cat -A | head -3
165b45c4fab20e59c91eff58968f41e7a68708c3$
$ git rev-parse HEAD
165b45c4fab20e59c91eff58968f41e7a68708c3
$ git rev-parse --abbrev-ref HEAD
dev
$ git log -1 --format="%h %ad %s" --date=iso
165b45c4 2026-10-08 10:26:13 +0800 180-c1: AC#1 premise adjudged STALE ...
```

rc=0 for all four.

⚠ The task relay said "起手 HEAD 不早于 `2dcef1a6`". Measured HEAD is `165b45c4`, which **is**
downstream of `2dcef1a6` (`git log --oneline -8` shows `2dcef1a6` = 182-c1 ticket append, parent
of `165b45c4`). So the floor is satisfied; the anchor of this leg is `165b45c4`, **not** `2dcef1a6`.

⚠ Drift is live: parallel legs are committing during this run (`165b45c4` was created 10:26, this
anchor written ~10:3x). HEAD may move again. The audited objects are pinned by their own SHAs
below, so the audit is not HEAD-dependent.

## Tool versions

```
$ git --version ; go version
git version 2.52.0.windows.1
go version go1.27.1 windows/amd64
```

## Working-tree census (porcelain)

```
$ git status --porcelain | wc -l
754
$ git status --porcelain --ignored | wc -l
804
$ git status --porcelain | grep -vc -E 'design/|\.gitignore|probes/161|probes/33/r10/'
710
```

rc=0. `754` is the raw line count. The `design/` + `.gitignore` + `probes/161/**` in-flight
work is present and will not be touched. `?? -` exists (a file literally named `-`) — not mine,
not touched.

⚠ The relay's first-hand number for this repo state is not quoted here because I did not
reproduce it; if a relay said `1476`, my ruler says `754` and the difference is reported in
`verdict.md` Q6 as a relay discrepancy, not as a repo defect.

## The three audited commits

```
$ git log -1 --format="%H | %ad | %s" --date=iso 44d178c9
44d178c946ca8c21ad35211ef5d93789d96edfe5 | 2026-10-08 10:08:42 +0800 | 33-r10 anchor: AC#13 three-cell re-measurement on 128bf600 (probe-first order already in place; box premise expired, content assertion cell still open)
$ ... 70b00885
70b00885683457956c253713641bfbb9723a97f0 | 2026-10-08 10:16:41 +0800 | test(33-r10 AC#13 item 2): give the cold-start page handover a window-free content ruler
$ ... db6ef1b1
db6ef1b1912f2b7d7b4a1518d6f2d98b472f44b9 | 2026-10-08 10:23:18 +0800 | docs(33-r10): file the AC#13 evidence set and the ticket addendum naming the expired premise

$ git merge-base --is-ancestor <each> HEAD ; echo rc=$?
44d178c9 ancestor rc=0
70b00885 ancestor rc=0
db6ef1b1 ancestor rc=0
```

## `70b00885 --stat` (claims to be 41 lines changed in the product file, +30/-11)

```
$ git show --stat --format="%H %s" 70b00885
 cmd/wisp/panel_host_windows.go                |  41 +++-
 cmd/wisp/panel_pageover_33r10_windows_test.go | 293 ++++++++++++++++++++++++++
 2 files changed, 323 insertions(+), 11 deletions(-)
```

rc=0. Stat confirmed as relayed: 41 changed lines in the product file, 293-line new test file.
`git show --stat` counts a changed line once, so `+30/-11` is the per-side split to re-derive
(see `01-product-diff.md` for my own statement-sequence ruler).

## The three audited files: tracked, not ignored, and identical to HEAD

```
$ git check-ignore -v cmd/wisp/panel_host_windows.go cmd/wisp/panel_pageover_33r10_windows_test.go .scratch/wisp/probes/33/r10/impl.md
rc=1   (no output = none of the three is ignored)
$ git ls-tree HEAD --name-only -- <those three paths>
.scratch/wisp/probes/33/r10/impl.md
cmd/wisp/panel_host_windows.go
cmd/wisp/panel_pageover_33r10_windows_test.go
```

## CRLF ruler (start readings) — worktree CR bytes vs HEAD blob CR bytes vs hashes

```
$ for f in cmd/wisp/panel_host_windows.go cmd/wisp/panel_pageover_33r10_windows_test.go .scratch/wisp/probes/33/r10/impl.md; do
    wt=$(tr -cd '\r' < "$f" | wc -c)
    head=$(git show "HEAD:$f" | tr -cd '\r' | wc -c)
    echo "$f | worktreeCR=$wt | headCR=$head | wtHash=$(git hash-object "$f") | headHash=$(git show "HEAD:$f" | git hash-object --stdin)"
  done
cmd/wisp/panel_host_windows.go                | worktreeCR=0 | headCR=0 | wtHash=1f9060dfff33cffab1317e5f653e1a95f9678e97 | headHash=1f9060dfff33cffab1317e5f653e1a95f9678e97
cmd/wisp/panel_pageover_33r10_windows_test.go | worktreeCR=0 | headCR=0 | wtHash=9b9b27478fcae98841f4d202bce0794bc3c97e8a | headHash=9b9b27478fcae98841f4d202bce0794bc3c97e8a
.scratch/wisp/probes/33/r10/impl.md           | worktreeCR=0 | headCR=0 | wtHash=a77866624accf3ddef867b619e873d12e8c7f1f6 | headHash=a77866624accf3ddef867b619e873d12e8c7f1f6
```

rc=0. All three files are **LF-only in both worktree and HEAD blob**, and worktree hash == HEAD
hash. Consequence for Q6: the `gofmt -l` CRLF-phantom caveat does **not** bite these three files;
it may still bite other files of the package, and that is measured separately in `06-hygiene.md`.

## Evidence dir of the audited leg, contents

```
$ find .scratch/wisp/probes/33/r10 -type f ! -name '*.md' ! -name '*.txt' ! -name '*.tsv'
... 01-product-diff.log, gate-*.log (5), mut/mut-*.log (4) ...
```

rc=0 — i.e. the r10 leg filed `.log` files. `.log` is **not** gitignored (root `.gitignore:8` is
`*.out`), and those files are tracked and clean, so this is a note about my own deliverable
convention, not a defect claim against r10. My files: `.md` / `.txt` / `.tsv` only.

## Root `.gitignore:8`

```
$ sed -n '8p' .gitignore
*.out
```

## Identity note (per instruction)

`QODER.md` in this directory records that this agent identifies as **Qoder** and declines to
discuss underlying-model claims. See that file.

## What this leg will not do

- Will not modify any tracked file. Mutations only via `go test -overlay` against copies in
  `/d/tmp/wisp33v4/`.
- Will not open a real window (constraint 3).
- Will not touch any `- [ ]` box on the ticket face (Q4 gives a verdict only).
- Will not fix, adopt, or count the pre-existing lock-shape finding as this leg's or as 33-r10's
  defect (Q1 attribution rule as instructed).
- Never `git add -A`, `--amend`, `reset`, `rebase`, `stash`, `checkout .`, `clean`, `push`.
