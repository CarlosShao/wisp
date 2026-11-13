# Audit: rename of tickets 269 and 435 (commit `b3b3db813`)

- **Auditor**: read-only audit leg (this session). No repo file other than this report was written.
- **Date**: 2026-11-13 (+08).
- **Commit under audit**: `b3b3db813d5dfe603894940f658ac151e166ee65` — "shorten the two ticket names I was responsible for breaking, and drop their roster entries ..."
- **HEAD at audit time**: `36890d56ac3b8888462b7df3e9d28c9b8c99c9ae` (A1819 ledger entry; its parent is `b3b3db813`).
- **Scope**: the two renames + the two roster deletions. Readings below are verbatim from the commands as run.

## Verdict summary

| # | Check | Verdict |
|---|---|---|
| 1 | rename recorded as `R`? | **PROBLEM** — `A`, and the old paths are still in the commit's tree |
| 2 | old names absent from the index? | **OK** (index only) — but HEAD still holds both |
| 3 | ticket content preserved? | **OK** — 55=55, 147=147, blob hashes identical |
| 4 | roster entries removed correctly? | **OK** — both old entries gone, new names absent; but the given grep command is invalid |
| 5 | gate really green? | **PROBLEM** — GREEN only on the index; HEAD replays RED (71 over-budget, 2 unrostered) |
| + | "those are historical readings, do not rewrite" | **Correct** — per `issues/README.md` rule 6 |

---

## 1. Did the rename actually happen? — **PROBLEM**

Command: `git show b3b3db813 --name-status --format=''`

```
A	.scratch/wisp/issues/269-gofmt-gate-red-by-own-sample-done.md
A	.scratch/wisp/issues/435-locked-suffix-roster-missed-four-names-done.md
M	scripts/check-path-length-budget.sh
```

`A`, not `R`. There is no `D` for either old name. The commit **added** the two new files and **kept** the two old files.

Corroboration:

- `git ls-tree -r --name-only b3b3db813` → **four** entries: old 269 + new 269 + old 435 + new 435.
- `git ls-tree -r --name-only b3b3db813^` → **two** entries: old 269 + old 435.
- `git show b3b3db813 --stat --format=''` → `3 files changed, 202 insertions(+), 2 deletions(-)`; the two new files are 55 and 147 insertions, the script loses 2 lines. **No deletion of either old file.**
- `git diff --cached --name-status` →

```
D	.scratch/wisp/issues/269-gofmt-step-is-held-red-by-its-own-committed-mutation-sample-and-eats-three-unprotected-steps-done.md
D	.scratch/wisp/issues/435-the-33r11-locked-suffix-roster-missed-four-new-names-and-the-question-is-whether-they-earned-the-suffix-done.md
D	.scratch/wisp/issues/438-path-length-gate-is-red-on-nineteen-unrostered-paths-fourteen-of-which-i-created-today.md
```

The deletions are **staged in the index only**; they are not in the commit. HEAD (`36890d56a`) still contains both old names.

The commit message claims: "Renames used git rename rather than a plain file move, which is the specific mistake that created seventeen duplicate paths earlier today." For this commit that claim is **false**: it is an add-only commit that duplicates both paths. It is the same shape of defect the message says it avoided — the old path is still tracked, so the tree holds two copies of each ticket.

## 2. Old names in the index? — **OK** (index), but HEAD is not clean

Command: `git ls-files '.scratch/wisp/issues/269-*' '.scratch/wisp/issues/435-*'`

```
.scratch/wisp/issues/269-gofmt-gate-red-by-own-sample-done.md
.scratch/wisp/issues/435-locked-suffix-roster-missed-four-names-done.md
```

One entry each: the index is clean. But `git ls-tree -r --name-only HEAD` still lists both old and new for each ticket. The index is ahead of the commit; the commit is the deliverable and it is **not** clean.

## 3. Ticket content preserved? — **OK**

```
old 269 (b3b3db813^): 55
new 269 (disk):       55
old 435 (b3b3db813^): 147
new 435 (disk):       147
```

Blob hashes (`git show <rev>:<path> | git hash-object --stdin`):

- 269 old = `d9d4d73d4545444e0a530575b0cd9fe8ebcfcab5` = 269 new
- 435 old = `f05cab14dd913e14930c8df220dfad00b7a87e33` = 435 new

Content is byte-identical. The rename, where it exists, did not alter the tickets.

## 4. Roster deletion correct? — **OK** (roster), but the given grep command is invalid

The commit's script diff removes exactly the two old-name entries:

```
-    R[".scratch/wisp/issues/269-gofmt-step-is-held-red-by-its-own-committed-mutation-sample-and-eats-three-unprotected-steps-done.md"] = "hist: its name was 99 chars until I appended -done on 2026-11-13, which pushed it to 104; the ticket itself predates README rule 9, so this is debt rather than a new violation"
-    R[".scratch/wisp/issues/435-the-33r11-locked-suffix-roster-missed-four-new-names-and-the-question-is-whether-they-earned-the-suffix-done.md"] = "hist: filed 2026-11-13 already over the 100-char hat, before any -done suffix was appended; created too long today, recorded rather than renamed because the -done suffix is the duplicate-claim key"
```

Corrected ERE grep:

```
$ grep -cE "269-gofmt-step-is-held-red|435-the-33r11-locked-suffix" scripts/check-path-length-budget.sh
0
$ grep -cE "269-gofmt-gate-red-by-own-sample-done|435-locked-suffix-roster-missed-four-names-done" scripts/check-path-length-budget.sh
0
```

Roster total: **69** entries (66 under `issues/` + 3 under `dispatches/`), matching the gate's own "the roster holds 69 entries".

> ⚠ The command given in the task, `grep -c 'A|B'`, is **not a valid test**. In GNU grep BRE, `|` is a literal pipe, so the pattern matches the literal string `A|B`, which never occurs in any file. It returns 0 unconditionally. The verbatim run returned `0`, but that 0 is meaningless. Use `grep -cE` (or `grep -nE`). With the corrected regex the roster is genuinely clean.

## 5. Gate really green? — **PROBLEM**

Index-based run (the command given):

```
check-path-length-budget.sh: positive control PASSED
check-path-length-budget.sh: denominator read: 19499 tracked paths
check-path-length-budget.sh: denominator: tracked paths=19499  over-budget=69  covered by roster=69  not in roster=0
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
```

The gate reads `git ls-files` (the index). The index has the old files staged for deletion, so this green is real for the current index — and **only** for the current index.

HEAD replay (read-only: a temporary index populated from HEAD via `GIT_INDEX_FILE=<tmp> git read-tree HEAD`, then the same script; the repo's own index and working tree were not touched):

```
check-path-length-budget.sh: denominator read: 19502 tracked paths
check-path-length-budget.sh: denominator: tracked paths=19502  over-budget=71  covered by roster=69  not in roster=2
check-path-length-budget.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/269-gofmt-step-is-held-red-by-its-own-committed-mutation-sample-and-eats-three-unprotected-steps-done.md  (relative 125 chars, name 104 chars, full path on the self-hosted runner 169 chars, budget 165)
check-path-length-budget.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/435-the-33r11-locked-suffix-roster-missed-four-new-names-and-the-question-is-whether-they-earned-the-suffix-done.md  (relative 136 chars, name 115 chars, full path on the self-hosted runner 180 chars, budget 165)
check-path-length-budget.sh: RED - count guard: the roster holds 69 entries, the tree has 71 over-budget tracked paths
check-path-length-budget.sh: VERDICT RED
```

So the debt is **not paid in the commit**. The green is a preview of the uncommitted index. If the staged deletions are committed, the gate will stay green; as committed, HEAD is RED with the two old paths unrostered.

---

## Extra: "those are historical readings, not to be rewritten" — **judgment correct**

`issues/README.md` rule 6 (lines 34–38), verbatim:

> 补裁另有 2 框 FAIL、2 框 PARTIAL。旧报告不追溯改写，但 13 张已 done 票的未决框由票 64 等消化。

`issues/README.md` rule 8 exception (lines 53–54), verbatim:

> ⚠ **例外**：证据文件里被当作"可重跑凭据"引用的快照目录**不许任何人事先删**——引用它的读数一旦失效，那张裁决表就从〔独立复现〕掉回〔仅自述，不背书〕。

**Conclusion: the judgment is correct.** Rule 6 says old reports are not retrospectively rewritten; the ledger is append-only; the probe logs are moment-in-time readings. Rewriting the 19 + 3 references would falsify the record. The correct handling is what was done: record the rename in a **new** ledger entry (A1819, `docs/reports/pending-and-issues.md:36183-36184`), which carries the old→new mapping. The blob hashes are identical, so the referenced content still exists under the new name and the reading remains reproducible through that mapping.

Two corrections to the premise:

1. **The location claim "都在台账和证据件里" is inaccurate.** Of the 19 references to old 269: 16 are in `.scratch/wisp/probes/**`, 1 in `.scratch/wisp/issues/438-...done.md`, 2 in `docs/reports/pending-and-issues.md`, **0** in `docs/evidence/s1/`. Of the 3 references to old 435: 1 in `docs/evidence/s1/audit-435-verdict.md`, 1 in the ledger, 1 in the 438 ticket. Most are probe logs, not ledger/evidence. This does not change the verdict — probe logs are also historical readings — but the claim as stated is wrong.
2. **"改名破坏了 ... 按名引用" is true only in the working tree.** At HEAD the old files still exist, so the references still resolve in the committed state. The breakage is a consequence of the uncommitted staged deletions, not of the commit.

---

## Bottom line

`b3b3db813` does **not** contain the rename. It adds the two short names and removes the two roster entries, but leaves both old paths tracked. The index has the deletions staged; HEAD does not. The gate is green only because it reads the index.

To actually pay the debt, the staged deletions must be committed (the 269 and 435 old paths, and the 438 old path alongside them), or the commit must be redone as a real `git mv`/rename so that `git show --name-status` reports `R`. Until then, the committed tree holds duplicate paths and the committed gate is RED.
