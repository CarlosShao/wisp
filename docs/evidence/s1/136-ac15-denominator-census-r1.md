# 136 AC#15 - denominator census r1 (READ-ONLY)

agent: `auditor-ticket136-ac15-denom-r1`
date: 2026-09-24
role: census only. Nothing built, run, formatted or edited except this file.
Constraint honoured: a resource-sampling CI run is live on this box, so no `go build|test|vet`,
no `gofmt`, no `wisp`, no `docker` were executed, not even once. Every number below is either a
line on disk (Read/Grep) or a `git` read, or an arithmetic derivation shown next to it.
Where only execution could have answered something, the exact command is written down instead
(see section 5).

Working anchors (all self-measured, none copied from the briefing):

| item | value | how measured |
|---|---|---|
| HEAD at start of this file | `5889559b767586ab73dd859f0a5f00e208612f49` | `git rev-parse HEAD` |
| AC#15 target census r2, carried by | `1bb92ccd76f95486dd06cc1e3a2cede1faf15fdc` (2026-09-24 21:20:15 +0800) | `git log --oneline -3 -- docs/evidence/s1/136-ac15-target-census-r2.md` |
| AC#11 pre-fix anchor | `51e29b044bf928e300c4a48ca1cc0932e02300b4` | `git cat-file -t` = commit |
| AC#11 fix commit | `f06a8d0f3ad1ba9e39b3b77e213ec7c6c0d742bf` | `git cat-file -t` = commit |
| AC#11 post-fix anchor used by the acceptance run | `e2a74631cced9901490c6e68af6b5b4e519c8d77` | `git cat-file -t` = commit |

Three shas named in my briefing do NOT exist in this repository and are not used anywhere below:
`840c9a9` (cited to me as "the AC#15 census r2 commit"; `git cat-file -t 840c9a9` =
`fatal: Not a valid object name`; the real carrier is `1bb92cc`). Also non-resolvable were the two
paths my briefing asserted as fact, see section 1.0. This is the "notifications replay stale
fragments" shape; logged, not obeyed.

---

## Section 1. Are AC#11's archived readings still on disk, and what exactly do they contain?

### 1.0 First: the paths my briefing gave do not exist

| asserted path | on disk? | evidence |
|---|---|---|
| `D:\tmp\wisp136ac11r2\` | ABSENT (no such directory) | `ls -d /d/tmp/wisp136ac11r2` = `No such file or directory` |
| `D:\tmp\wisp136ac11\` | ABSENT as a directory | `ls -d /d/tmp/wisp136ac11` = `No such file or directory`; what exists is a *prefix* |
| `base_observe.txt` | ABSENT | `find /d/tmp -maxdepth 2 -iname "*base_cmd*"` and `-iname "*observe*"` return no such name |
| `base_cmd.txt` | ABSENT | same |
| `full_observe.txt` | ABSENT | same. The only `*observe*` hits under `/d/tmp` at depth 2 are `w124-ac1-logs/L1.__internal_observe_.txt`, `L2.…`, `P1.…`, `wisp136-ac1-m1-observe-test.txt`, and three `observer_cost_test.go` copies. None is an AC#11 per-path log. |

So: **AC#11's readings ARE still on disk, but not where my briefing says, and not under those
filenames.** AC#11 wrote one `-v` log **per shot** (per-run files), not one log per path. The real
on-disk archive is listed below in full. Nothing under `D:\tmp` was moved, renamed or deleted by
this run (create-only rule); the absent paths were not created either, since recreating someone
else's archive under someone else's name is not my remit.

Prefix families that exist, with the process each belongs to (attribution read from the three
AC#11 evidence tables, then re-verified against the files themselves):

| prefix | whose run | which table says so |
|---|---|---|
| `wisp136ac11-*` | the impl run (r1 of AC#11) | `136-ac11-impl.md` §1 and the temp-artifact list at its foot |
| `wisp136ac11b-*` | the second witness | `136-ac11-second-witness-readings.md:11` ("临时件前缀一律 `wisp136ac11b-`") |
| `wisp136ac11v-*` | the r1 acceptance/verdict run | `136-ac11-r1-acceptance.md:134` (driver `/d/tmp/wisp136ac11v-batch.sh`) |
| `ac11-orch-base.*`, `ac11-orch-after.*`, `ac11-orch-tree/`, `ac11-orch-after/` | the orchestrator's own readings | `136-ac11-orchestrator-readings.md:79` |

### 1.1 Whole-package `-count=1` batches (denominator (i))

Each file is exactly one `go test -count=1 -v ./internal/observe/` invocation against one fixed
`git archive` tree. Tallies are from the files themselves: `^=== RUN` for RUN, and the **top-level**
(no leading whitespace) `^--- PASS` / `^--- FAIL` / `^--- SKIP` for the other three.

| set (all under `D:\tmp`) | files | RUN | top PASS | top FAIL | top SKIP | bytes | mtime span (2026-09-24) |
|---|---|---|---|---|---|---|---|
| `wisp136ac11-batch-before/` | 30 | 1950 | 1947 | 3 | 0 | 201761 | 14:07:48 - 14:09:37 |
| `wisp136ac11-batch-before2/` | 30 | 1950 | 1948 | 2 | 0 | 201630 | 14:25:16 - 14:26:40 |
| `wisp136ac11-batch-after/` | 30 | 1950 | 1950 | 0 | 0 | 201510 | 14:14:59 - 14:16:34 |
| `wisp136ac11-batch-after2/` | 30 | 1950 | 1950 | 0 | 0 | 201510 | 14:21:36 - 14:23:07 |
| `wisp136ac11b-batch-before/` | 30 | 1950 | 1948 | 2 | 0 | 201633 | 14:33:11 - 14:34:32 |
| `wisp136ac11b-batch-after/` | 30 | 1950 | 1950 | 0 | 0 | 201510 | 14:34:37 - 14:36:00 |
| `wisp136ac11b-batch-before-DISCARDED-contended/` | 30 | 1950 | 1950 | 0 | 0 | 201510 | 14:18:49 - 14:20:18 |
| `wisp136ac11v-batch-before/` | 30 | 1950 | 1948 | 2 | 0 | 203691 | 14:36:10 - 14:37:37 |
| `wisp136ac11v-batch-after/` | 30 | 1950 | 1950 | 0 | 0 | 203541 | 14:37:51 - 14:39:42 |
| `ac11-orch-base.log.1` … `.60` (+ `ac11-orch-base.idx`, 716 B) | 60 | 3900 | 3898 | 2 | 0 | 403140 | 14:13:28 - 14:16:30 |
| `ac11-orch-after.log.1` … `.60` (+ `ac11-orch-after.idx`, 716 B) | 60 | 3900 | 3900 | 0 | 0 | 403020 | 14:18:19 - 14:21:18 |

Totals: **390 whole-package invocations**, RUN 25350, top PASS 25339, top FAIL 11, top SKIP 0.
Closure check: 25339 + 11 = 25350 = RUN, so nothing in the archive is a lost or double-counted line.
Every `before`-family set also carries an `index.txt` (or `gate.txt` / `gate-post.txt` for the `v`
sets) with per-shot `rc` and the contention-gate record.

Note on `ac11-orch-base`: my first probe of it failed because it is a **flat file prefix** in
`D:\tmp`, not a directory (`ls -d /d/tmp/ac11-orch-base` fails; `ls /d/tmp/ac11-orch-base.log.*`
returns 60 files). It is present, and it is the 60-shot pre-fix batch the second witness's §3 row
attributes to the in-flight orchestrator run, with the two target hits at `.log.20` and `.log.43`
(both `goroutine_test.go:33: PerTask mid-task = 2, want 3`, both files 6777 B against the 6717 B
baseline). Re-verified here, not taken from the table.

### 1.2 How many `-count=1` invocations each set represents (arithmetic shown)

The per-invocation top-level name count is read **from the logs themselves**, not assumed:

* `grep -c '^=== RUN' <log>` = 65 for every one of the 390 files (spot check `25.v.log` = 65;
  exhaustive check on the 240-file subset in section 2.2 returned 240 of 240 at exactly 65).
* `grep -c '^ *--- '` on the same file = 65, and there are **no** indented sub-result lines
  (`    --- PASS:` count is 0), so 65 is simultaneously the RUN count and the top-level
  verdict count: the package `internal/observe` had **65 top-level test names and zero subtests**
  at this anchor.
* Therefore invocations = (set's RUN total) / 65. Each 30-file set: 1950 / 65 = **30**.
  Each 60-file orch set: 3900 / 65 = **60**. Nine 30-sets + two 60-sets = 270 + 120 = **390**.
* Cross-check against a single file's own size: the 65-run baseline log is 6717 B; a 1-name run is
  not observed anywhere, so no file in these sets is a partial or filtered run.

### 1.3 Same-process `-count=N` artefacts (denominator (ii))

| file | bytes | mtime | invocations | RUN | top PASS | top FAIL | SKIP | package verdict lines |
|---|---|---|---|---|---|---|---|---|
| `wisp136ac11-before-count500.log` | 44182 | 14:28 | 1 (`-count=500`) | 500 | 498 | 2 | 0 | 1 |
| `wisp136ac11-after-count500.log` | 44061 | 14:28 | 1 (`-count=500`) | 500 | 500 | 0 | 0 | 1 |
| `wisp136ac11-count2.log` | 13237 | - | 1 (`-count=2`) | 130 | 130 | 0 | 0 | 1 |
| `wisp136ac11-count2.log.pre.gate` | 1981 | - | aborted probe, RUN=0 | 0 | 0 | 0 | 0 | 0 |
| `wisp136ac11b-count2-after.log` | 13237 | 14:3x | 1 (`-count=2`) | 130 | 130 | 0 | 0 | 1 |
| `wisp136ac11v-probe-500.log` | 1325 | 14:4x | 1 (`-run` filtered probe) | 1 | 1 | 0 | 0 | 1 |

Same-process repeats on disk: 500 + 500 + 2 + 2 = **1004**, spread over 4 invocations.
Grand total of shots AC#11 left on disk: 390 + 1004 = **1394**.

### 1.4 The "1016" in my briefing is not a number this repository contains

`grep -rn 1016 docs/ .scratch/wisp/issues/` returns exactly one hit, and it is unrelated:
`docs/evidence/s1/66/66-nokeeper-2.json:47: "tree_private_bytes": 10162176` (a substring of a byte
count). No AC#11 table, ticket cell, or report records a 1016-run reading. AC#11's own tables state
their denominators as `4/60`, `2/30` (x3 independent pre-fix batches), `2/60`, `0/30`, `0/60`,
`2/500`, `0/500`, `33/500`, `12/500`, `6/500`. Anyone needing a single aggregate should use the
per-file census in 1.1/1.3 above (390 whole-package, 1004 same-process) rather than an unsourced
sum; the two must not be added together (see 1.5).

### 1.5 What AC#11's archived readings say about the **AC#15** symptom specifically

This is the load-bearing fact for AC#15. `--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss`
(the AC#15 target) was searched for across all 390 + 1004 archived shots:

| batch | shots | AC#15-target hits | AC#11-target hits | any other top FAIL |
|---|---|---|---|---|
| `wisp136ac11-batch-before` | 30 | **1** (`25.v.log`) | 2 (`15`, `29`) | none |
| `wisp136ac11-batch-before2` | 30 | 0 | 2 | none |
| `wisp136ac11-batch-after` / `-after2` | 60 | 0 | 0 | none |
| `wisp136ac11b-batch-before` | 30 | 0 | 2 | none |
| `wisp136ac11b-batch-after` | 30 | 0 | 0 | none |
| `wisp136ac11b-batch-before-DISCARDED-contended` | 30 | 0 | 0 | none |
| `wisp136ac11v-batch-before` | 30 | 0 | 2 | none |
| `wisp136ac11v-batch-after` | 30 | 0 | 0 | none |
| `ac11-orch-base.log.1..60` | 60 | 0 | 2 (`.20`, `.43`) | none |
| `ac11-orch-after.log.1..60` | 60 | 0 | 0 | none |
| `*-count500.log` + two `count2` logs | 1004 repeats / 4 inv. | **0** | 2 (both in `before-count500`) | none |

The single AC#15 hit, verbatim (`/d/tmp/wisp136ac11-batch-before/25.v.log`, lines 88-90 of 135):

```
=== RUN   TestCheckSettleHalfTheReadsFailedReportsItsLoss
    sampler_settle_coverage_136_test.go:189: precondition broken: only 2 reads taken, half-and-half needs a window to lose in
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.13s)
```

That file's own four numbers: RUN=65, top PASS=64, top FAIL=1, SKIP=0, size 6848 B, mtime 14:09:16.
The other two FAILs in that batch (15 and 29, both 6777 B) are `TestNoopTaskReturnsToBaseline`, a
different case and a different family - the impl table records this at line 108
("第 25 发那枚 FAIL **不是**目标用例，是本包另一枚既有 flake"), and the second witness records the
same split at §3 ("r1 那批有 3 枚 rc=1 但只有 2 枚目标红").

**Derived rate, stated with its denominator:** 1 hit in **390** whole-package `-count=1`
invocations (0.256%), and 0 hits in **1004** same-process `-count=N` repeats. These two are
separate denominators and must never be summed - AC#11 already nailed that rule for its own case
(`136-ac11-orchestrator-readings.md:101-103`): whole-package single-shot 6.7% versus same-process
500-repeat 0.4% for the *same* case is a 16x spread, which is precisely why a merged percentage is
meaningless.
