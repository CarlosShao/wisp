# 136-v — anchor, drift, surface, gates (verification leg 136-v, non-implementer)

## 0. anchor read at hand-off

- `git rev-parse --short HEAD` (start, 2026-10-14 ~19:5x +08) = **`1740fd48`** ("133-r1 交付（甲）").
  The dispatch's landing chain tops at `f8b92e7e`; HEAD has moved on — drift item 1.
- real repo = `D:\work\workspace\projects plans\Wisp` (has `.git`; default branch `dev`, origin/HEAD -> dev).

## 1. drift, one by one (each measured, nothing inferred)

1. **HEAD moved**: `f8b92e7e` -> `1740fd48`. `git diff f8b92e7e..HEAD -- scripts/` = EMPTY, so no
   neighbour touched the 136 surface.
2. **dispatch's "only slo-freshness.sh in scripts/" does not hold over today's range**:
   `git diff --name-only b38567a9..HEAD -- scripts/` returns TWO files —
   - `scripts/slo-freshness.sh` = the 136s four commits;
   - `scripts/check-path-length-budget.sh` = 1+/1-, commit `41071cdc` ("编排者：178 收口") — another
     leg's commit, NOT 136s. The 136s face itself is clean: each of its seven commits is single-purpose
     (four code commits: only `scripts/slo-freshness.sh`; three delivery commits: only
     `.scratch/wisp/probes/136/s/**` + the ticket file).
3. **dispatch's "133+/22-" is the AC#2 commit only**, not the range total. Measured numstats:
   `bb47e9e7` 133/22 (AC#2), `966c2aa9` 41/6 (AC#3), `f706dedd` 25/3 (AC#4), `b6fc818d` 15/2 (seam);
   sequential sum 214/33; **net `git diff --numstat b38567a9..HEAD` = 211+/30-** (later patches rewrite
   lines earlier patches inserted, so only the net is a fact about the pair of trees);
   `-U0` hunks over the whole range = 20 (per commit: 13+6+7+1 = 27 for the same reason).
4. **P1 surface unchanged**: `git diff --name-only b38567a9..HEAD -- .github/` = EMPTY; `ci.yml` blob is
   `80b7a393` in both trees, so all four of my trees reuse one ci.yml.
5. **the `:-20` 收口 is NOT committed anywhere**: it exists only as an unstaged worktree edit of
   `scripts/slo-freshness.sh` line 138. `git show origin/dev:scripts/slo-freshness.sh | :138` still
   reads `:-3`, while `origin/dev` already contains all five 136s commits (`merge-base --is-ancestor`
   per sha, all IN). Consequence measured live — see §3/live-readings.txt: the 12:07:32Z schedule run
   on dev already reds with `exit 2`.
6. **two misdirected temp trees (my own tooling slips, kept per 只建不删, registered here)**:
   - `D:\d\tmp\wisp136v\worlds\*` — first `gen_worlds.py` run resolved POSIX `/d/tmp/...` under the
     current drive root, writing 33 fixture dirs to `D:\d\tmp\wisp136v\`. Real fixtures live in
     `D:\tmp\wisp136v\worlds`.
   - `D:\work\workspace\projects plans Wisp\.scratch\wisp\probes\136\v\` — first `gen_probe_files.py`
     run used a space instead of the slash before `Wisp`; the 5 data files briefly landed in a lookalike
     twin that has NO `.git` and already held other legs' same-slip content (`probes/272`, an empty
     `probes/orch` chain). The real deliverables were then written to the repo path; the twin copies
     stay in place and this line is their register.

## 2. gates (all re-run by this leg, not copied)

| gate | command | rc | reading |
|---|---|---|---|
| `sh -n` | `sh -n scripts/slo-freshness.sh` (both head and wt) | 0 | clean |
| `dash -n` | `dash -n <script>` (both versions) | 0 | clean (CI runs the pin with ubuntu `/bin/sh` = dash) |
| shellcheck | `docker run --rm -i koalaman/shellcheck:v0.11.0 -s sh -` (both versions) | 0 | no output = clean; CI lint step is exactly `shellcheck -s sh scripts/slo-freshness.sh` |
| d22scan | `sh scripts/d22scan.sh` | 0 | "clean - no D22 ban violations" (scripts/ itself is outside its scopes) |
| length gate | `sh scripts/check-path-length-budget.sh --with-self-test` | 0 | VERDICT GREEN; over-budget=60, roster=60, not-in-roster=0, worst path 229 |
| GOTMPDIR | `mkdir C:\tmp\wisp136vgotmp` | - | created, backslash form only (no Go build was run) |
| `git status -- scripts` | before/after | ` M scripts/slo-freshness.sh` | identical at start and end |

## 3. cross-checks of the artifacts themselves

- **136s's own harness re-run** (WS redirected to `/d/tmp/wisp136v/theirws`, not their temp): its
  `drive.sh` reproduces their `logs/gates.txt` rc lists 1:1 — old: control 0, spread 0, cap-zero 1,
  forge-* 0/0/0, missing-ts 0, missing-source 0, partial-green 0, partial-notfound 1, future 0/0,
  stale 1, k1..k6 2/2/2/2/1/1, k4b 0, k7 0; head: cap-zero 2, forge 1/1/1, missing 2/2, partial 2/2,
  future 2/2, stale 1, k's 2/2/2/2/1/1, k4b 2, k7 2. (Raw: `out/s-drive-old.txt`, `out/s-drive-head.txt`.)
- **my own matrix** (`rctable.tsv`, 33 worlds x 4 trees) agrees with every shape above and adds:
  boundaries `total>cap` (100 vs 101 at cap=1), the 2000/2001 boundary under the `:-20` default, a
  scan-path future stamp, a mixed forgery+real page, and a full-API-path world (p2-live-path).
- **cross-shell**: the whole matrix re-run with `dash` (CI's shell) — 132/132 rc identical.
