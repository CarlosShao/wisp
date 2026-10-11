# 304-v1b 20 — my own AC#3 gate re-runs (all four rulers)

Judge leg 304-v1b; anchor `3b8066a2`; all runs on the shared worktree, read-only gates, each log self-lands an rc line.

## Rulers + readings (mine)

| # | ruler (exact cmd) | result | rc line / log |
|---|---|---|---|
| 1 | `sh scripts/d22scan.sh` | `clean - no D22 ban violations`; ban#8 internal/=533 Go files incl. comments+_test | `d22scan_rc=0`; `logs/g-own-d22scan.txt` |
| 2 | `sh scripts/check-path-length-budget.sh --with-self-test` | VERDICT GREEN; denominator `tracked paths=11022 over-budget=57 covered by roster=57 not in roster=0`; worst full path=228 (hat 165) | `pathlen_rc=0`; `logs/g-own-pathlen.txt` |
| 3 | `bash scripts/portable-tests.sh --scope=census` | totals line `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` — byte-equal to the ticket's frozen baseline | `census_rc=0`; `logs/g-own-census.txt` |
| 4 | format, two parallel: `gofumpt -l` on worktree ⇄ on `git show HEAD:<path>` blobs (fresh `/d/tmp` tree) | worktree roster = 0 lines; blob roster = 0 lines; `cmp` worktree vs blob OK x4; range = the four files `e875f00f` touched | `fmt_worktree_rc=0` / `fmt_blob_rc=0` / `cmp_rc=0`; `logs/g-own-fmt.txt` |

Over-reach ruler (per-commit form, shared with AC#1): the four `304-r1` commits stay inside
`probes/304/r1/**` + the four `internal/llm` files; only `27f7b197` touches the ticket (1 add / 0 del, boxes untouched).
Logs: `logs/ac1-commits-detail.txt`, `logs/ac1-bullet.txt`. Tag check: newest tag `wisp-orch-backup-0927-1057` (2026-09-27)
predates ticket 304 — no tag movement (`logs/g-own-zerobyte.txt`).

## The three 0-byte files from `304-r1` (mandated check)

`logs/build-1.txt`=0B, `logs/vet-1.txt`=0B, `logs/vet-reverse.txt`=0B. Same cell DOES carry rc-bearing readings elsewhere:
`g-build.txt` = `build_rc=0`, `g-vet.txt` = `vet_rc=0`, `g-summary.txt` = census before/after `IDENTICAL` + stated baseline.
=> These three are empty shells (superseded first-attempt files), not missing cells. Named as hygiene residue, not a hole.
Inherited from dead leg 304-v1 additionally: `g2-fmt-worktree.txt`=0B, `g3-build-body.txt`=0B, `g3-vet-body.txt`=0B;
its `g3-build.txt`/`g3-vet.txt` carry `rc=0` companions, but its format cell never landed anywhere — covered by my ruler #4 above.

## Verdict

`AC#3` = **HOLDS** on my own re-runs. Named residue: the 0-byte shells stay on disk (this project keeps temp artifacts);
the dead leg's format cell was empty and is supplied by my re-run in this file.
