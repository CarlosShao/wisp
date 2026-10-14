# 136-v method — my own rig (not a rerun of probes/136/s/drive.sh)

Everything lives outside the repo under `/d/tmp/wisp136v/` (created, never deleted):
`bin/` (stub, generators, drivers), `trees/`, `worlds/`, `out/` (all raw logs and rc tables).

## trees (built from git, never edited)

| tree | source | role |
|---|---|---|
| `old`  | `git show b38567a9:scripts/slo-freshness.sh` (378 lines) | pre-AC#2/#3/#4 |
| `pre`  | `git show f706dedd:...` | AC#2+#3+#4 in, before the `b6fc818d` seam fix |
| `head` | `git show HEAD:...`; `cmp` == `git show b6fc818d:...` byte-identical | committed delivery |
| `wt`   | copy of the working-tree file | head + orchestrator `:-20` edit, line 138 only (`diff` = that one line) |

Each tree also carries `.github/workflows/ci.yml` (same blob `80b7a393` as the repo).

## worlds (33) — what each models

`real` fresh valid sample; `forge-branch` / `forge-workflow` / `forge-event` name-right provenance-wrong;
`missing-ts` scanned record with no `created_at` key; `missing-source` scanned record with no
`workflow_run`; `mixed` forged-newer + genuine-older on one page; `seam-forged-aged` a forged record
through the documented `SLO_FULL_LAST_SAMPLE` fixture seam; `seam-empty-first` leading-TAB seam record;
`k1-no-token` / `k2-no-repo` / `k3-runs-fail` / `k4-pages-fail` / `k5-empty` / `k6-runs-empty` the six
"cannot get data / no data" shapes; `k4b-total-fail` / `k7-run-fail` the two new seams (total-count
query, run look-up); `p2-live-path` full API path (runs list + jobs + scan, no seam);
`p-green` / `p-notfound` / `cap-zero` the AC#3 partial-scan shapes; `ctrl188` / `spread2` (newest valid
sample on page 2); `bnd-at` (total==cap=100) / `bnd-over` (101) the strict-inequality boundary;
`live349` (total 349, default cap) live-shaped; `d-cap2000` / `d-cap2001` the 20-page boundary;
`future-sample` / `future-job` / `future-scan` future stamps on seam sample / seam job / scanned record;
`stale-sample` / `stale-job` 10-day-old records; `real` doubles as the AC#4 control.

## `gh` stub (own implementation — `bin/gh` + `bin/ghstub.py`)

- endpoint classification by URL shape; fixtures are per-world JSON files (`total.json`,
  `artifacts-pageN.json`, `runs.json`, `run-<id>.json`, `jobs.json`);
- the caller's `--jq` text is matched against MY transcription of the two artifact filters (3-col and
  4-col), the run filter, the runs filter, the jobs filter, and `.total_count`; whitespace-normalised;
  ANY unrecognised filter => exit 1 => the pin lands in its "cannot look" exit 2, never a made-up answer;
- `STUBGH_FAIL` injects failures per kind (`runs` / `jobs` / `run` / `total` / `pages`);
- a page with no fixture answers `{"total_count":T,"artifacts":[]}` — the past-the-end behaviour the pin
  itself documents (re-measured 09-23); modelling it is what lets the 20-page default be exercised with
  one fixture file;
- stdout forced LF-only (a stray CR in the last TSV field would false-reject the branch).

## driver

```
env -i PATH=<stubdir>:/usr/bin:/bin HOME=... sh -c 'set -a; . <world>/env || exit 9; exec sh <script>'
```

Each world's whole environment (seams, repo/token, cap, injected failures) is one `env` file — nothing
inherited. `bin/driver.sh <tree>` runs all 33 worlds; `bin/driver-dash.sh` is the identical driver with
`dash` as the shell. Raw logs: `out/<tree>--<world>.log`; rc tables: `out/rctable.<tree>.txt`
(bash) and `out/rctable-dash.<tree>.txt` (dash).

## reproduce

```
cd /d/tmp/wisp136v
<PY> bin/gen_worlds.py                  # writes worlds/*
sh bin/driver.sh old ; sh bin/driver.sh pre ; sh bin/driver.sh head ; sh bin/driver.sh wt
sh bin/driver-dash.sh head              # (any one tree for the cross-shell check)
<PY> bin/gen_probe_files.py             # assembles this directory's data files
```

live section (real API, no stub): `GH_TOKEN=$(gh auth token) GITHUB_REPOSITORY=CarlosShao/wisp
sh <tree>/scripts/slo-freshness.sh` — token never written to disk.

## instrument boundary (what these readings can and cannot say)

- Every stub world drives the pin through its own documented seams
  (`SLO_FRESH_NOW` / `SLO_FULL_LAST_TRIGGER` / `SLO_FULL_LAST_SAMPLE` / `SLO_FULL_ARTIFACT_PAGES`)
  plus a stub `gh` first on PATH. No world edits the script, its assertions or `date`; the matrix is
  offline.
- The live section is the only network reading; it runs the real script against the real API and was
  used for the default-cap and production questions. Transient `EOF` errors were observed there twice
  (script page 7; my own raw query) — recorded as-is, they are the reason the fail-closed exit 2 is not
  a theoretical branch.
- The three "cannot look" exit-2 classes (no token / no repo / query failed) are exercised with the
  stub; the real token path is exercised live.
