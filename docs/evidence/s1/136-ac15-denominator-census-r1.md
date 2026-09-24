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

---

## Section 2. Is the "1 in 240" claim salvageable from any archived source?

**Short answer: yes - and this contradicts the ruling currently in the tree.** The face's newest
`>` block and the truth-source ledger both now state that the 240 has *no* archival support. That
statement is checkable, and it is wrong on the checkable part. Details, then the caveat that keeps
this from being a clean number.

### 2.1 What the repo was asked to look like versus what it says

Two places currently assert the 240 is unsupported:

| where | text (trimmed) | state |
|---|---|---|
| ticket face `:375` (uncommitted `M` at the time of this section) | its item 3: "my earlier '1 red in 240' is registered as [not reproducible]: nothing in the roster supports that denominator, and it came from a relay" | **not reproducible** ruling |
| `docs/reports/pending-and-issues.md:5698`, bullet b | "'1 red in 240' has no archival support at all (no record in the roster, the source is a relay), so register it as [not reproducible] per the AC#11 precedent" | same ruling |
| ticket face `:351` (the AC#15 cell itself, older) | "在**盘上可查的 240 发整包 `-v` 里红 1 发**（14:3x 第二 witness 报回并留在它的表里；我核过那枚文件与那一段行号真实存在）" | claims the opposite, and names the table |

### 2.2 The record exists, verbatim, in one place - and it is a first-party reading

`docs/evidence/s1/136-ac11-second-witness-readings.md:289-291`:

```
它是**票 136 AC#12** 那格今天新装的用例（`f53ad5c`／`c03aee3`），红因形状是"100 ms 窗口只取到 2 次读数"
（窗口／间隔计时形状，不是 §4 那条注册表边）。我在盘上可查的 **240 发**（r1 60 ＋ 他人 120 ＋ 本程 60）里
只命中这一发 ⇒ 罕见、**未立案**。⚠ **本格不修、票面不勾、也不替它开 AC**——只登记现象与出处，
```

This is not a relay. It is the second witness writing about a batch it selected, counted and
re-scanned off disk ("盘上现重扫（不是我抄谁的数）", same file `:137`). Two further citations of the
same figure exist, both downstream of it and both honest about that:
`136-ac14-impl.md:252` and `136-ac14-r1-acceptance.md:383` (the latter re-locates the face's pointer
from `:266` to the cell's own line). A repo-wide `grep -rn "240" docs/ .scratch/` returns **no other**
candidate source: every other hit is a line number, a sha (`4e5d240`), a byte count, or an unrelated
240-line CI log (`134-ac6-contended-no-conclusion.md:17`).

### 2.3 The 240 reconstructs off disk, file by file, with no gaps

Reading the `r1 60 / 他人 120 / 本程 60` decomposition against section 1.1's sets:

| leg | sets it names | shots |
|---|---|---|
| r1 60 | `wisp136ac11-batch-before` (30) + `wisp136ac11-batch-before2` (30) | 60 |
| 他人 120 | `ac11-orch-base.log.1..60` (60) + `ac11-orch-after.log.1..60` (60) | 120 |
| 本程 60 | `wisp136ac11b-batch-before` (30) + `wisp136ac11b-batch-after` (30) | 60 |
| total | | **240** |

Verified mechanically, not eyeballed: a loop over exactly those 240 files counting `^=== RUN` per
file returned `240-subset: files with RUN=65 = 240 ; others = 0`. All 240 are whole-package
`-count=1` invocations of `./internal/observe/`, i.e. one denominator, the right one. The single hit
is `wisp136ac11-batch-before/25.v.log:89-90`, quoted in section 1.5.

Two things make the 240 *internally admissible* rather than a grab bag:
* The four legs run on **three different trees** (AC#11 pre-fix `51e29b0`, fix `f06a8d0`, post
  `e2a7463`), but the AC#15 target file is **byte-identical on all three**: md5
  `79711ce3a032fdfe3e632cc59912b7eb` at each of `51e29b0`, `f06a8d0`, `e2a7463`
  (`git show <a>:internal/observe/sampler_settle_coverage_136_test.go | md5sum`). AC#11's fix touched
  only `internal/observe/goroutine_test.go` (`git diff --name-only 51e29b0 f06a8d0` = one path). So
  pooling them is legitimate **for the AC#15 symptom**, even though pooling would be illegitimate for
  AC#11's own rate.
* 150 of the 240 are on pre-fix trees and 90 on post-fix trees; the one hit sits in the 150. Stated
  both ways: 1/240 pooled, 1/150 on the pre-fix subset, 0/90 on the post-fix subset.

### 2.4 Where the face's wording genuinely overreaches (four points, all checkable)

1. **"1 in 240" is a chosen subset, not the archive.** 150 more whole-package invocations of the same
   identical target file are sitting on disk and were left out: `wisp136ac11-batch-after` (30),
   `-after2` (30), `wisp136ac11b-batch-before-DISCARDED-contended` (30), `wisp136ac11v-batch-before`
   (30), `wisp136ac11v-batch-after` (30). The full honest on-disk figure is **1 hit in 390**
   invocations (0.256%), all 390 verified at RUN=65. The subset direction is anti-conservative for
   the claim's *authority* (a smaller denominator flatters the rate) and irrelevant to its
   *existence* (0 extra hits either way).
2. **It is a sighting, not a rate.** 1 hit in 240 has a 95% interval of roughly 0.01% to 2.3%, so it
   is compatible with 1/10000 and with 1/44. Anyone reusing it must carry that.
3. **It is denominator (i) only.** Zero of the 240 are same-process repeats. The 1004 same-process
   repeats on disk (`*-count500.log`, two `count2` logs) contain **0** occurrences of the AC#15
   symptom, so the (ii) reading of this event is 0/1004 and no composite number is lawful.
4. **It is stale as a baseline.** `internal/observe/` has moved since `51e29b0`:
   `git diff --name-only 51e29b0 HEAD -- internal/observe/` returns 4 paths, and the target file
   itself is now md5 `a31968022574b99b7969d0ccddb27eb0`, 365 lines against 289. Per this repo's own
   reuse rule (check `git diff <anchor>..HEAD -- <pkg>` is empty before reusing an archived reading),
   the 240 fails the check for *current-tree* use, and census r2 section 2 quantifies why it matters:
   the package's ticker-window budget per shot has roughly doubled (about 1560 ms then, about
   2960 ms now), so 1/240 more likely *understates* today's rate than overstates it.

### 2.5 Verdict for section 2

The repair the ledger asked for is available without running anything: **restore the 240 as a sourced
sighting and correct its pointer, do not delete it.** Concretely, the admissible sentence is
"1 hit in 240 whole-package `-count=1` invocations of `./internal/observe/`, pooled from four
archived batches (`136-ac11-second-witness-readings.md:290`; all 240 logs re-counted on disk at
RUN=65 on 2026-09-24 by this table); **1 hit in 390** across the complete on-disk set;
sighting not rate, 95% interval about 0.01% to 2.3%; denominator (i) only; taken on a package
version whose sampling exposure is about half the current one." The *inadmissible* form is the
bare "1 in 240" that the cell currently prints, and equally inadmissible is the replacement the
ledger now carries ("no archival support"), because the archival support is a named file, a named
line, and 240 files that are all still on disk.

I am not editing the face or the ledger for that - both are outside my one writable path, and the
face's `>` blocks are append-only by rule. This section is the counter-evidence, filed for whoever
owns those two files.

---

## Section 3. Portability: what AC#11 can AC#15 reuse verbatim, and what it cannot

### 3.1 (a) The outer-loop command shape - two variants are archived, only one is reusable

Both drive `go test -count=1 -v ./internal/observe/` against a `git archive` tree in `D:\tmp`
(never the repo work tree), one `.v.log` per shot.

**Shape 1, the plain loop** - `/d/tmp/wisp136ac11b-run-batch.sh` (15 lines, r1 and second witness):

```
for i in $(seq 1 "$N"); do
  ( cd "$TREE" && go test -count=1 -v ./internal/observe/ ) > "$OUT/$i.v.log" 2>&1
  rc=$?
  printf 'run=%s rc=%s at=%s\n' "$i" "$rc" "$(date '+%Y-%m-%d %H:%M:%S %z')" >> "$OUT/index.txt"
done
```

Gate is a separate hand-run script. `date` is written per shot and **never subtracted** (both scripts
carry the comment "host wall clock jumps, so no durations are derived").

**Shape 2, the gated loop** - `/d/tmp/wisp136ac11v-batch.sh` (about 110 lines, the AC#11 acceptance
run). Same core loop plus four things AC#15 needs and shape 1 lacks:
* a bounded contention gate, four checks: (1) `powershell Get-Process` matching
  `Runner.Worker|^go$|compile|^cgo$`; (2) `docker ps` (container load is invisible in the host list);
  (3) `gh run list -R CarlosShao/wisp --limit 5 --json databaseId,status,headSha` counting
  `"status":"in_progress"`, degrading to `CLEAR-LOCAL-ONLY` when `gh` errors, which is explicitly
  **not** an assertion that the queue is empty; (4) `ls -ld` mtime of
  `/e/work/base/actions-runner/_work`;
* poll every 30 s, max 60 rounds (30 min), `exit 8` if it never clears - "never takes numbers inside a
  contention window";
* **tree identity captured into `gate.txt` at fire time**: md5 of the target `_test.go` plus
  `find internal/observe -name '*.go' | sort | xargs md5sum`;
* post-batch gate re-scan into `gate-post.txt`, and per-log `BATCH-TAG`/`RUN-DATE` header + `RC=` trailer.

Verdict: **reuse shape 2 verbatim**, changing only `<tree> <outdir> <runs> <tag>`. Reusing shape 1
would reproduce AC#11's own weakest moment - the second witness had to void a whole 30-shot batch
(`wisp136ac11b-batch-before-DISCARDED-contended`) because it had no in-loop gate and someone else's
44 log files landed during its window.

One thing AC#15 must add and AC#11 did not need: AC#11's target was in `goroutine_test.go`, which the
fix itself rewrote, so pre/post used two different trees by design. AC#15's fix will also change
`internal/observe/**_test.go`, so the md5 block must name **all four** timing-family files, not one,
otherwise a pre-fix batch and a post-fix batch can silently differ in more than the intended hunk.

### 3.2 (b) Shell-loop `-count=1` versus same-process `-count=N`: AC#11 ran BOTH, and the archive settles the census's structural claim

| denominator | did AC#11 run it | archived volume | exact archived command |
|---|---|---|---|
| (i) shell loop of `-count=1` | yes | 390 invocations (section 1.1) | `go test -count=1 -v ./internal/observe/` |
| (ii) same-process `-count=N` | yes, **but `-run`-filtered** | 4 invocations / 1004 repeats | `go test -count=500 -v -run 'TestNoopTaskReturnsToBaseline$' ./internal/observe/` (both `*-count500.log`), and `go test -count=2 -v ./internal/observe/` (two `count2` logs, RUN=130) |

So AC#15 can reuse both shapes, but note the asymmetry: AC#11's same-process volume was taken on
**one test name**, not on the whole package. For AC#15 that is exactly the right precedent, because
`-count=1440` on the whole 71-name package would cost about 80 minutes of ticker windows in one
process and would also drift the process's own heap/timer state across the run.

**Confirming the census claim "one `go test` invocation can never yield 30 `-count=1` runs": CONFIRMED
from the archive, not by argument.** Counting package result lines (`^(ok|FAIL)[[:space:]]+github`) per
file: `wisp136ac11-batch-before/1.v.log` = **1**, `wisp136ac11-before-count500.log` = **1**. A
`-count=500` log has RUN=500 and 500 per-name verdicts but exactly one package outcome, one `rc`, one
process. Conversely, 30 whole-package outcomes occupy 30 separate files. The two readings are
therefore not interconvertible, which is the physical content of AC#15 criterion 1's ban on summing.

The size of the difference is already measured, for the *same* case, in AC#11's own table
(`136-ac11-impl.md:119-120`): whole-package single-shot **4/60 = 6.7%** versus same-process 500
repeats **2/500 = 0.4%**, a 16.7x spread on one identical defect. For AC#15's symptom the archive
gives (i) 1/390 and (ii) 0/1004.

### 3.3 (c) AC#11's roster and panic counting rules, verbatim from the instrument

Source: `/d/tmp/wisp136ac11-summarize.py` and its declared copy `/d/tmp/wisp136ac11b-summarize.py`
(the copy's only stated change is writing the roster dump under a new prefix).

| rule | regex / predicate | portability note for AC#15 |
|---|---|---|
| RUN names | `^=== RUN\s+(Test\w+)` | portable as-is; `\w+` cannot cross `/`, so subtests are excluded from RUN by construction |
| verdicts | `^\s*--- (PASS\|FAIL\|SKIP):\s+(\S+)` | **latent bug**: the leading `\s*` lets indented subtest verdicts into the same counters that RUN excludes them from. Harmless today - `internal/observe/**_test.go` has **0** `t.Run(` occurrences (measured file by file: 14 files, all 0) - but it silently mis-tallies the first time anyone adds a subtest. Fix or note it before reusing |
| hit | `nm == target or nm.startswith(target + "/")` | portable; pass the full name `TestCheckSettleHalfTheReadsFailedReportsItsLoss` |
| panic | `line.startswith('panic')` or (`' [running]:' in line and line.startswith('goroutine ')`) | already the *true-panic* sense. Keep it |
| fatal error | `line.startswith('fatal error')` | keep |
| roster closure | `Counter` of sorted per-run name tuples, then prints `roster: K distinct top-level run-name sets over N runs; modal size=S occurs=M` and per-run `missing=` / `extra=` against the modal set; dumps `<prefix>-<label>-roster.txt` | portable verbatim |
| cross-batch | `comm -3` on flattened per-run name lists, e.g. `/d/tmp/wisp136ac11b-b1-names.txt`, `-a1-names.txt`, `-r1-names.txt` (2350 B, 65 names each) | portable; the golden list must be regenerated at 71 names, see 3.5 |

Two archived outputs to copy the reporting format from:
pre-fix `TOTAL runs=30 hits=2 (rate 2/30)` / `four-number spread: RUN=[65] PASS=[64, 65] FAIL=[0, 1]
SKIP=[0]` / `roster: 1 distinct top-level run-name sets over 30 runs; modal size=65 occurs=30`;
post-fix identical but `hits=0`, `PASS=[65]`, `FAIL=[0]`.

**The panic two-senses rule, and one name the face gets wrong.** The instrument's line-initial rule is
the right one, and every archived log confirms it is quiet: `grep -c '^panic:'` is 0 on all 30
`wisp136ac11-batch-before` logs, and no file among the 90 `before` logs of the `r1` and `v` sets or
the 60 `ac11-orch-base` logs has a single line-initial `panic`. But the identifier sense is a
different quantity and the face's new `>` block, item 4, mis-names it: it says the 7 identifier hits are
"`wantPanic` 与 `TestCheckSettlePanicInSamplerDoesNotStopTick`", yet
`grep -rn "TestCheckSettlePanicInSamplerDoesNotStopTick" --include=*.go .` returns **nothing** - that
test does not exist. The only two test names in the repo containing `Panic` are
`TestPanicInWorkerSurvivesAndCancelsRoot` (`internal/observe/goroutine_test.go`) and
`TestFakeTreeEmptyScriptFailsClosedAndNotPanics` (`internal/observe/sampler_faketree_guard_136_test.go:54`),
and the identifier mass in `internal/observe/**_test.go` is `PanicCount` (3), `PanicEvent` (2),
`panic(` (2 sites) plus those two names. Whoever runs AC#15 should re-measure this and not carry the
7.

### 3.4 (d) What "30 shots with 0 hits" must mean numerically, per denominator

Cost basis, measured not guessed: `/d/tmp/ac14b-r2-136/test-head.txt` is a `-count=2 -v` whole-package
run at 20:32 today - RUN=142, top PASS=142, top FAIL=0, SKIP=0, one package line,
`ok github.com/CarlosShao/wisp/internal/observe 6.626s`, i.e. **3.313 s of package time per
whole-package shot at 71 names** (against 2.67 s/shot observed across the archived 65-name batch
spans). Budget 3.3-4.0 s per invocation including process startup.

**Denominator (i), 30 shell-loop `-count=1` invocations at today's HEAD.** Expected file shape:
30 logs; per log `^=== RUN` = 71, top `--- PASS` = 71, top `--- FAIL` = 0, top `--- SKIP` = 0,
`^panic:` = 0, exactly 1 package result line, index 30 lines all `rc=0`. Batch totals: RUN 2130,
PASS 2130, FAIL 0, SKIP 0; `1 distinct set over 30 runs`, modal size 71 occurs 30; both `comm`
directions empty. Statistical content: with k=0 the 95% one-sided upper bound is `3/30` = **10%**.
That number is *blind* to the effect it is meant to detect: the archived sighting is 1/240 = 0.417%
and the corrected on-disk figure is 1/390 = 0.256%. So 0 hits in 30 whole-package shots is an entry
ticket, not a measurement.

**Denominator (ii), `-count=30 -v -run '<target>$' ./internal/observe/` in ONE invocation.**
Expected: RUN 30, top PASS 30, FAIL 0, SKIP 0, `^panic:` 0, **one** package result line, one `rc`,
about 30 x 0.13 s = 4 s of wall time. Same `3/30` = 10% arithmetic, but the distribution being bounded
is the insensitive one - AC#11 measured that same shape as roughly 16x less likely to show the defect
than (i), and AC#15's symptom has 0 occurrences in the 1004 same-process repeats already on disk.
Reporting this as "30 shots, 0 hits, fixed" is the single most misleading thing this cell could do.

**So criterion 4's ">= 30 shots with 0 hits" is satisfiable literally twice over and means something
different each time.** Write it as: `(A) 0 hits in n_A whole-package -count=1 invocations, giving an
upper bound of 3/n_A; (B) 0 hits in n_B same-process repeats, reported separately as a mechanism
probe, never merged.` AC#11's precedent for that wording is
`136-ac11-orchestrator-readings.md:101-103` ("any 'this cell is X%' phrasing must carry its
denominator, otherwise it becomes the next person's source of error").

### 3.5 What is NOT portable

1. **The roster size.** AC#11's golden list is 65 names; today's package is 71
   (`grep -c '^func Test' internal/observe/*_test.go` = 71 across 14 files, and independently
   confirmed by `ac14b-r2-136/snap*-v.txt` at RUN=71 and `test-head.txt` at 142 for `-count=2`).
   The 6-name delta is AC#14's `sampler_settle_gate_136_test.go`. Every AC#11 `comm -3` baseline must
   be regenerated; reusing `wisp136ac11b-b1-names.txt` as the golden set would report 6 phantom
   "extra" names on every shot.
2. **The tree-identity check.** AC#11's 240/390 shots all sit on `sampler_settle_coverage_136_test.go`
   md5 `79711ce3…`; at HEAD it is `a3196802…` (289 lines to 365), and
   `git diff --name-only 51e29b0 HEAD -- internal/observe/` lists 4 paths. Under this repo's reuse
   rule the archived rate is therefore historical, not a baseline.
3. **Exposure per shot has doubled.** Census r2 section 2 tallies the package's ticker-window budget
   at about 1560 ms then versus about 2960 ms now, of which the AC#15 target itself contributes 100 ms
   in both. Reusing 1/240 as a *prediction* for the current tree would understate the rate, and the
   pre-fix-on-current-tree rate is the only admissible baseline.
4. **The `-run` filter AC#11 used for (ii) is load-bearing**, not decoration: `AC#15`'s target shares
   its binary with 70 other cases, so an unfiltered `-count=1440` would be a 130-minute run that also
   measures 1440 executions of 70 unrelated legs.
5. **Do not count AC#14's mutation reds as sightings.** An exhaustive (uncapped) scan of all of
   `D:\tmp` over `*.log` and `*.txt` finds **18** files containing
   `^--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss`, written to
   `/d/tmp/wisp136ac15denom-r1-scan-failnames.txt` (created, not deleted):
   AC#11's `wisp136ac11-batch-before/25.v.log`, plus 17 mutation-anchored reds from other runs
   (`ac14b-r2-136/snapC-v.txt` `snapD` `snapE` `snapF`, `wisp136ac1213-tree2-MD-v.txt`,
   `-tree2b-MD-v.txt`, `-tree2b-MF-v.txt`, `wisp136ac14/mut/M1.log`, `M2.log`,
   `M3b-gate-true-fold-out/M3b-gate-true-fold-out.log`, `wisp136ac14r1/logs/mut-M1.txt`, `mut-M2.txt`,
   `mut-M3a-t.txt`, `mut-M4.txt`, `wisp136r1-tree-MD-v.txt`, `wisp136r1-tree-MF-v.txt`,
   `wisp141gate/logs/m1-v.txt`). That is 18 exactly: the one flake plus 17 deliberate reds. They fail
   on a *different sentence*, e.g. snapC line 91 at
   `sampler_settle_coverage_136_test.go:251`
   `AC#14: a half-covered settle window must carry its own sampling verdict row, got verdicts=[]`.
   Re-scanning all of `D:\tmp` for AC#15's own flake sentence
   (`precondition broken: only 2 reads taken, half-and-half`) returns **exactly 1 file**, the same
   `25.v.log` (list at `/d/tmp/wisp136ac15denom-r1-scan-sentence.txt`). So across everything ever
   archived on this box, AC#15's symptom has been seen **once**. Filter on the sentence, never on the
   test name, when building the rate; a name-only filter would report 18 sightings and inflate the
   rate 18x with deliberate mutation reds.

---

## Section 4. Cheapest defensible plan for AC#15's first measurement

Everything here is derived from sections 1-3; **nothing in this section was executed.**

### 4.1 The two command lines, one per denominator

**(A) denominator (i), whole-package `-count=1`, one process per shot.** Preferred form is AC#11's
gated driver copied verbatim to a new prefix and widened (section 3.1), invoked as

```
sh /d/tmp/wisp136ac15-batch.sh /d/tmp/ac15-tree-pre /d/tmp/wisp136ac15-A 120 PRE-A
```

The command inside that loop, spelled out so it is checkable without reading the script:

```
( cd /d/tmp/ac15-tree-pre && go test -count=1 -v ./internal/observe/ ) > /d/tmp/wisp136ac15-A/$i.v.log 2>&1
```

Run in 12 chunks of `n = 120` (chunk, not one 1440-line loop) because AC#11's own gate bounds the
wait for a clear window at 60 rounds x 30 s = 30 min, and a 120-shot chunk is about 8 min, so each
chunk gets its own pre-fire gate record and its own post-batch re-scan.

**(B) denominator (ii), same-process repeats of the target alone.**

```
cd /d/tmp/ac15-tree-pre && go test -count=1200 -v -run 'TestCheckSettleHalfTheReadsFailedReportsItsLoss$' ./internal/observe/ > /d/tmp/wisp136ac15-B-target.log 2>&1; echo "RC=$?" >> /d/tmp/wisp136ac15-B-target.log
```

The `-run` anchor with `$` is not cosmetic: without it `-count=1200` replays all 71 names 1200 times.
AC#11's precedent for this exact shape is `go test -count=500 -v -run 'TestNoopTaskReturnsToBaseline$'`
(`136-ac11-impl.md:117`).

### 4.2 The n arithmetic for a "1 in about 240" symptom

The right question is not "how many shots until I see it" but "what does k=0 prove". With zero hits,
the one-sided 95% upper bound is the p solving `(1-p)^n = 0.05`, i.e. `p_up = 1 - 0.05^(1/n)`,
which for these n is the rule-of-three `3/n`.

| n (shots of A) | `p_up = 3/n` | what it rules out |
|---|---|---|
| 30 | 10% | nothing - the archived sighting is 24x smaller |
| 300 | 1.0% | still 2.4x looser than 1/240 |
| 720 | 0.417% | **exactly equal to 1/240, so it proves nothing** (this is why `n = 720` buys nothing) |
| 1167 | 0.257% | exactly equal to the corrected on-disk 1/390 |
| **1440** | **0.208%** | strictly tighter than **both** 1/240 (0.417%) and 1/390 (0.256%); equivalent to "bounded below 1 in 481" |
| 1200 (for *estimating*, not bounding) | expects about 5 hits at 1/240, about 3 at 1/390 | a count you can put an interval on |

The decisive way to see that 30 is not a measurement is detection probability, `1 - (1-p)^n`:

| n | chance of catching it at p = 1/240 | at p = 1/390 |
|---|---|---|
| 30 | 11.8% | 7.4% |
| 100 | 34.2% | 22.6% |
| 240 | 63.2% | 46.0% |
| 500 | 87.7% | 72.4% |
| 1200 | 99.4% | 95.5% |

So a 30-shot clean batch is the **expected** outcome roughly 9 times out of 10 even if the defect is
untouched. Recommendation: **n = 1200 for the pre-fix rate measurement (A), n = 1440 for the post-fix
proof (A)**, each with its own (B) companion at the same n, and the post-fix number reported as
"95% upper bound 0.208%, i.e. below the archived sighting", never as "0 hits therefore fixed".

Cost, using the measured 3.313 s of package time per 71-name shot (section 3.4) plus startup, budget
3.8 s end to end: A at n=1200 is about 76 min, A at n=1440 about 91 min, B at `-count=1200` about
156 s (the target case is 0.10-0.13 s per execution, read off the archived PASS/FAIL timings).
Total for one full pre/post cycle: about 3 hours 20 minutes of gate-checked wall time, and **B is
where the cheap repetitions belong**.

One budget warning the later run must not be surprised by: AC#15 criterion 3's fix *adds* bounded
waiting to up to 6 legs. If each added bound is 200-400 ms and all 6 names run per shot, A's per-shot
cost grows by roughly 1.2-2.4 s, i.e. n=1440 goes from about 91 min to about 2.5-3 h. Recompute the
chunking **after** the fix lands, from the pilot's own measured per-shot span, not from this number.

### 4.3 The baseline file's expected four numbers, written to catch a LOST reading

The point of the baseline is that it must fail loudly when a reading goes missing, which a
"FAIL=0" expectation does not.

Per shot of (A) at HEAD, all seven of these must hold, and any single failure makes the shot **void,
to be re-fired**, never a clean zero:

| check | expected | catches |
|---|---|---|
| `grep -c '^=== RUN'` | 71 | truncated / interrupted log |
| `grep -c '^ *--- \(PASS\|FAIL\|SKIP\):'` | 71 | results that never printed |
| `grep -c '^--- PASS'` | 71 | |
| `grep -c '^--- FAIL'` | 0 | the symptom itself |
| `grep -c '^--- SKIP'` | 0 | skip laundering "untested" as "passed" |
| `grep -c '^panic:'` | 0 | a panic eating the rest of the roster |
| `grep -cE '^(ok\|FAIL)[[:space:]]+github'` | exactly 1 | a run that never reached a package verdict |
| `index.txt` has a `run=N rc=…` line for this N | present | a lost process whose output landed but whose rc did not |

Batch-level expectations for `n = 1200` on a green pre-fix baseline: 1200 logs, **RUN = 85200**,
`PASS + FAIL + SKIP = 85200`, top FAIL = 0, SKIP = 0, `^panic:` = 0, 1200 package lines, 1200 index
lines, and `roster: 1 distinct top-level run-name sets over 1200 runs; modal size=71 occurs=1200`.
A byte-size tripwire is cheap and worth having: in the archive a green 65-name shot is **exactly
6717 B in 264 of 390 shots**, and the three red variants are 6777 / 6780 / 6848 B. A 71-name green
shot should land near 7.3 KB; take the real number from the pilot's first shot and band it, do not
assume it.

The precedent for why this is mandatory rather than tidy: `internal/observe/sampler_test.go:330-335`
records the pre-guard shape measuring `52 RUN / 47 PASS / 5 FAIL + panic, 4 cases swallowed` - a
shrunken denominator that is pixel-identical to a clean batch. AC#11's own zero-hit claims were
defended exactly this way (`136-ac11-r1-acceptance.md:600`: "`RUN` 逐发同数、目标用例逐发都有
`--- PASS` 行").

For (B), the per-file expectation is different and must not be cross-checked against (A)'s numbers:
`-count=1200 -run '<target>$'` gives `^=== RUN` = 1200, `^--- PASS` = 1200, `^--- FAIL` = 0,
`^--- SKIP` = 0, and **exactly one** package line. If RUN is 1200 but PASS+FAIL is less, repeats were
lost mid-process and the whole log is void.

### 4.4 The exact names that must be in the roster diff

14 names. The 6 flake-capable ones (census r2's Group A collapses its 8 guard *rows* onto these 6
*owners*; row count is not name count), then 8 witnesses a "fix" could quietly damage.

Flake-capable set (Group A, must each appear every shot):

| name | file, owning func line | guard that can fire |
|---|---|---|
| `TestCheckSettleHalfTheReadsFailedReportsItsLoss` | `sampler_settle_coverage_136_test.go:201` | `:213 tree.reads < 4` (**the target**), plus shadowed `:216 kept < 2`, `:221 lost < 2` |
| `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | same file `:123` | `:125 reads < 3` |
| `TestCheckSettleZeroFootprintDropsAreCountedToo` | same file `:291` | `:293 reads < 3` |
| `TestCheckSettleFullyMeasuredWindowReportsNoLoss` | same file `:341` | `:343 reads < 3` |
| `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` | `sampler_settle_gate_136_test.go:198` | `:201 kept < 1 \|\| lost < 1` |
| `TestSampleStateAllMetricsAndVerdicts` | `sampler_test.go:84` | `:99 len(rep.Samples) < 3` |

Witness set (must stay present and green, so that a zero is not bought by deleting or skipping
something):

| name | why it is on this list |
|---|---|
| `TestNoopTaskReturnsToBaseline` | AC#11's just-closed leg, same package, same binary. Deleting or skipping it would show as a shrunken roster, not as a red |
| `TestSampleStateZeroSampleWindowFailsClosed` | AC#1's fail-closed nail (`136-ac1` whole cell rests on this name) |
| `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable` | AC#1's leg B, the counterpart that must stay green |
| `TestCheckSettleZeroTrustworthySamplesFailsClosed` | settle-side twin of AC#1's nail, guard `settle_zerosample:62 reads == 0` |
| `TestCheckSettleTrustworthyReadsAreRecorded` | same file, guard `:141 len(rep.Samples) == 0` |
| `TestSamplerGoroutineAccountingFollowsRegistry` | carries AC#8's `len(rep.Samples) == 0` guard at `sampler_test.go:336`, the fix that stopped the 4-case swallow |
| `TestSettleCoverageRowExistsAndPassesWhenFullyMeasured` | gate file `:154`, guard `:156 reads < 1 \|\| len(rep.Samples) != reads` (Group B) |
| `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | gate file `:236`; **`:262` indexes `buildSettleVerdicts(SettleReport{})[0]` with no length guard**, so any fix that makes that builder return an empty slice turns this into a panic that eats every later name (census r2 5.4). Must be in the diff, and a FAIL here voids the batch rather than counting as a hit |

Roster mechanics, both directions, every chunk: golden list regenerated at 71 names from a pilot shot
(not from `wisp136ac11b-b1-names.txt`, which is 65), then `comm -23 run.txt reported.txt` empty
(ran but never reported) and `comm -13 run.txt golden.txt` empty (in golden but never ran), plus
`comm -3` between the pre-fix and post-fix golden lists so nothing appears or disappears.

### 4.5 Which names must run sequentially, and why

**All 6 flake-capable names, plus the whole package, strictly one invocation at a time.** The package
already cannot parallelise itself - `grep -rn 't.Parallel' internal/observe/` returns **0** hits, so
`go test` runs the 71 names in declaration order inside one binary. The parallelism that must be
banned is therefore external, and there are three distinct sources of it:

1. **An outer loop that backgrounds shots** (`… &`, `xargs -P`, `parallel`, or a nested
   `go test … & wait`). AC#11's archived loop is serial `for i in $(seq 1 N)` with no `&`, and it
   must stay that way. Two shots in flight would double the measured rate and make the number
   incomparable with the 390 archived single-shot readings.
2. **Sibling packages via `go test ./...`**: Go runs up to GOMAXPROCS *packages* concurrently, so
   the observe batch must name its one package explicitly, as every archived AC#11 command does.
3. **The self-hosted runner on this box** (the `slo-full` D32 job starts on push and steals CPU).
   This is what AC#11's four-check gate exists for; the gate record is part of the deliverable, and a
   batch fired in an uncleared window is void, not merely noisy.

Two of the 6 have a second, independent reason to stay serial: they read **process-wide** state.
`TestSampleStateAllMetricsAndVerdicts` and the registry legs consume the live
`Registry` counts (`RuntimeGoroutines`, `reg.Count()`), so any sibling holding a goroutine at that
instant changes the answer regardless of CPU load; AC#11's own root cause was exactly this class
(missing happens-before edge on a shared registry, `goroutine.go:276` versus the `defer` at
`:331-336`). And `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` is the panic candidate
above: under parallel invocation, two short logs look like two clean-ish runs, whereas serially its
truncation is unmissable in the roster diff.

Concretely for the batch script: keep `( cd "$TREE" && go test … ) > "$OUT/$i.v.log" 2>&1` exactly as
archived, add no `&`, no `-parallel`, no second package path, and gate around each 120-shot chunk.

---

## Section 5. Things only execution can answer, so they are written down instead

Not run, per the machine constraint. Each line is runnable as-is by the later measurement run.

1. The real per-shot wall time and the real green byte band at 71 names:
   `sh /d/tmp/wisp136ac15-batch.sh /d/tmp/ac15-tree-pre /d/tmp/wisp136ac15-A-pilot 5 PILOT` then
   `wc -c /d/tmp/wisp136ac15-A-pilot/*.v.log` and the `index.txt` span.
2. Whether the target's hit rate on the **current** tree differs from the archived 1/390:
   needs section 4.1 (A) at n=1200; no shortcut exists.
3. `grep -c '^panic:' /d/tmp/wisp136ac15-A/*[!x].v.log` for the true-panic count at HEAD, and the
   identifier-sense count `grep -rin 'panic' internal/observe/*_test.go | wc -l`, to settle the face's
   "7 hits" and its one non-existent name (section 3.3).
4. `git diff <pre-fix-anchor>..<post-fix-HEAD> -- internal/observe/` must list only
   `internal/observe/**_test.go`; anything under `sampler.go` or `thresholds.go` means AC#15's
   boundary was crossed and the cell's stop-hand line is live.
5. The `d22scan` baseline check the fix submission owes: `sh scripts/d22scan.sh` with `internal/`
   not rising above its recorded baseline.

## Section 6. Discipline receipt

* Executed: `Read`, `Grep`/`Glob` equivalents over repo files, `git log|show|cat-file|rev-parse`,
  `ls`, `wc`, `cat`, `stat`, `md5sum`, `grep`, `awk`, `sed -n` on `D:\tmp` log files, and one `head`
  pipe. **Zero** `go build|test|vet`, zero `gofmt`, zero `wisp`, zero `docker`.
* Created (not deleted, per the create-only rule): this file;
  `/d/tmp/wisp136ac15denom-r1-scan-failnames.txt`;
  `/d/tmp/wisp136ac15denom-r1-scan-sentence.txt`. Nothing under `D:\tmp` was moved, renamed,
  overwritten or removed; no worktree or checkout was made inside the repo.
* Symbol hygiene: this file's own prose carries no emoji and no arrow or math symbols. The single
  exception is the block quote in section 2.2, reproduced byte-for-byte from
  `136-ac11-second-witness-readings.md:289-291` because its exact characters are the evidence; every
  other quoted cell here was paraphrased out of those codepoints rather than retyped with them.
* Staging hygiene: every commit used `git add -- <this one path>` then
  `git diff --cached --name-only` before `git commit -q -F - -- <this one path>`. No foreign path
  ever appeared in the staged list. The 16 `design/**` deletions and untracked `design/old/`,
  `design/doubao/` were never staged.
* Shared-worktree hazard actually observed: between two of my commits another agent landed
  `9f75cd2` (face) and `bc67475` (CI read), and the ticket face plus `136-ac15-target-census-r2.md`
  changed **while this file was being written** (census r2 grew from 113 to 617 lines; the face grew
  the `>` block that rules the 240 "not reproducible"). `git log -1` immediately after a commit can
  therefore show someone else's sha; each of my own commits was re-read with
  `git log -1 --format=%H -- <my path>` and verified with `git show --name-only`.
* Two counters, kept separately as required:
  * **genuine-notification echoes: 6** - two `MEMORY.md was modified` notices, one date-change notice,
    the `agents.md` memory injection, and three background-task completion notices. None carried an
    instruction I acted on.
  * **judged injection / not obeyed: 0** - no text in this session claimed authority asking me to
    loosen a criterion, revert, skip, or confirm itself. Separately, three *factual* claims in my own
    briefing were checked and found false (the two `D:\tmp` paths, `840c9a9`, and "the 240 was never
    recorded anywhere"); they are logged as stale briefing fragments in sections 1.0 and 2, not as
    injections, and I did not act on any of them.

next = whoever owns AC#15's measurement should read section 4 as the plan and section 2 as the thing
to correct on the face before the face's "not reproducible" ruling gets cited again.
