# 136-instr-fixes-r1 - instrument-side fixes in `internal/observe/sampler_settle_gate_136_test.go`

Worker: `worker-ticket136-instr-r1`. Date `2026-09-24`, tree anchored at HEAD `9048866`
(verified with `git cat-file -t 9048866` -> `commit`).

Batch scope: two defects, both **instrument-side (test code)** in the same file, one batch:

| # | coordinate | defect | disposition |
|---|---|---|---|
| 1 | `sampler_settle_gate_136_test.go:62` | `func gateFailed is unused` (staticcheck U1000 shape) | investigated first, then **used**, not deleted - see section 2 |
| 2 | `sampler_settle_gate_136_test.go:262` | `buildSettleVerdicts(SettleReport{})[0]` indexed with no length guard | length guard added - see section 3 |

Frozen and not touched by this batch: `internal/observe/sampler.go` (production),
`internal/risk/**`, `rules_gateway.go`, `tools/d22scan/**`, `allowlist.txt`, any
`thresholds.go`/golden, `frontend/**`, `design/**`, `docs/PLAN.md`, `docs/specs/**`,
the ticket face under `.scratch/wisp/issues/`, and `docs/reports/**`.

Shared-tree items deliberately left alone: `tools/d22scan/main.go` and
`tools/d22scan/scan_test.go` are ` M` uncommitted edits owned by another worker
(`worker-ticket141-q46c` per `docs/reports/HANDOVER.md:99`); `design/**` has 16 owner
deletions plus untracked `design/old/`, `design/doubao/`. None of these was staged,
committed, reverted or cleaned here.

Instrument note: the roster/four-number rules used below are top-level-line counts
(`^=== RUN  Test`, `^--- PASS: / ^--- FAIL: / ^--- SKIP: `, real panic = `^panic:`,
package result line = `^ok|^FAIL`), the same rules the AC#15 census instrument
`/d/tmp/wisp136ac15-verify/summarize.py` applies. Script kept at
`/d/tmp/wisp136instr-r1/summarize.sh` (temp artefacts live outside the repo; created,
never deleted).

## 0. Baseline, taken by this worker before any edit

`go test -count=2 -v ./internal/observe/` at HEAD `9048866`, log
`/d/tmp/wisp136instr-r1/pre.log`:

```
RUN=142  PASS=142  FAIL=0  SKIP=0     pkg_result_lines=1  real_panic(^panic:)=0
unique_top_level_names=71  names_seen_exactly_2x=71  verdict_lines=142
ok  github.com/CarlosShao/wisp/internal/observe  6.905s
```

That is the documented healthy state (142/142/0/0, 71 unique names each twice, 0 SKIP,
0 real panic). Settle-leg elapsed times at baseline, both passes:

```
TestSettleCoverageRowExistsAndPassesWhenFullyMeasured      0.20s  0.20s
TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed     0.20s  0.20s
TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured  0.40s  0.40s
TestFoldSettlePassOnlyGateRowsVeto                         0.00s  0.00s
TestSettleReportPassNeverContradictsItsGateRows            0.60s  0.60s
TestStateReportVerdictBuilderStaysSinglePurpose            0.00s  0.00s
```

Baseline gates (same tree, before the first edit):

| gate | command | result |
|---|---|---|
| gofmt | `gofmt -l internal/observe/` | empty, rc=0 |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/observe/` | empty, rc=0; `--version` prints `v0.12.0 (go1.27.1)` |
| vet | `go vet ./internal/observe/` | rc=0 |
| staticcheck | `"$(go env GOPATH)/bin/staticcheck.exe" ./internal/observe/` | **rc=1, no findings emitted at all** - see section 1 |
| d22scan | `sh scripts/d22scan.sh` | rc=1, 6 findings, all `frontend/` ban #8 - see section 4 |

`go version` on this box: `go1.27.1 windows/amd64`.

## 1. staticcheck: the binary is here, but it cannot read this toolchain

`ls "$(go env GOPATH)/bin"` -> `gofumpt.exe`, `staticcheck.exe`. So it **is** installed
locally. But:

```
$ staticcheck.exe --version
staticcheck 2025.1.1 (0.6.1)
```

CI pins `staticcheck@2026.2.1` and runs it under `go-version: 1.25.0`
(`docs/reports/pending-and-issues.md:5728`). This box has go1.27.1. Feeding
`./internal/observe/` to the local build produces only:

```
-: internal error in importing "internal/byteorder" (cannot decode "internal/byteorder", export data version 4 is greater than maximum supported version 2)
-: internal error in importing "internal/cpu"     (... same ...)
-: internal error in importing "internal/goarch"  (... same ...)
-: internal error in importing "math/bits"        (... same ...)
-: internal error in importing "unicode/utf8"     (... same ...)
```

rc=1, five `-: ... (compile)` lines, **zero findings**. Consequences, stated plainly:

- This worker **cannot confirm the U1000 finding, and cannot confirm it is gone**, with a
  local staticcheck. The gate that matters is CI's `staticcheck@2026.2.1`; nobody should
  read section 2 as "the lint gate was verified green here".
- What local staticcheck *does* prove is the negative direction only in one sense: it
  never even reaches the analysis phase, so the same 5 lines must appear before and after
  (they do, `sc-pre.txt` vs `sc-post.txt`), i.e. this batch changes nothing observable to
  that instrument.
- The load-bearing evidence for item 1 is therefore `grep` + the introducing commit's own
  enumeration, not a lint reading. That is what section 2 is built on.

Log files: `/d/tmp/wisp136instr-r1/sc-pre.txt`, `/d/tmp/wisp136instr-r1/sc-post.txt`.

## 2. Item 1, `:62 func gateFailed is unused` - answered before it was touched

The question the ledger demands an answer to (`pending-and-issues.md:5719`: "either somebody
uses it, or it was supposed to be used by somebody - whoever judges it has to answer that")
was answered with four readings, none of them inference from the function's name:

| # | reading | command | what it says |
|---|---|---|---|
| R1 | zero callers, ever | `grep -rn "gateFailed" internal/observe/` | 2 hits: `:61` comment + `:62` declaration. Whole-repo grep adds only the three ledger/evidence mentions. |
| R2 | born unused, not orphaned | `git show aef82f5:...gate_136_test.go \| grep -n gateFailed` | same 2 hits at the introducing commit; matches the census two-tree reading (`ci-read-a1fd5bf-r1.md:247`) that no caller was deleted. `aef82f5` verified: `git cat-file -t aef82f5` -> `commit`. |
| R3 | no leg was dropped | `grep -c "^func Test"` in the file = 6; `aef82f5`'s own message enumerates exactly 六枚探针 and its roster arithmetic says 65 -> 71 top-level = +6 | the delivered leg set is the planned leg set. AC#14 judgement 2's nail ("一半读数报错 ⇒ 该门行红、红名逐名") is `:198` via `gateAlternatingTree`, and `:217` proves the error text reaches the row note. So **no ticket judgement is short of a nail**. |
| R4 | where the shape went | `grep -rn "failedRead\|trustworthyRead()\|zeroFootprint()" internal/observe/*.go` | the fixture trio is copied from the sibling AC#15 file (`sampler_settle_coverage_136_test.go:72/:78/:83`), and there the error step **is** used (`:124`). In this file the "reads fail" nail was routed through `gateAlternatingTree` instead, leaving the copied error constructor without a caller. |

R4 however exposed the one thing that is genuinely not nailed anywhere. Enumerating the
window shapes reachable in package `observe`:

| window | kept / dropped | drop channel | nailed by |
|---|---|---|---|
| kept some, lost some | >0 / >0 | error | `gate_:198` (alternating), `coverage_:124` |
| kept some, lost some | >0 / >0 | zero footprint | `coverage_:292` |
| kept nothing | 0 / >0 | **zero footprint** | `gate_:236` (`gateZeroFootprint`), `gate_:336` |
| kept everything | >0 / 0 | - | `gate_:154`, `coverage_:342` |
| **kept nothing** | **0 / >0** | **error** | **nothing** |
| never measured | 0 / 0 | - | `gate_:270` (at the builder, synthetic) |

Coordinates are as **delivered** by this batch (so `gate_:336` is the pre-existing
`all reads untrustworthy` script line, which sat at `:331` before this file grew 12 lines
above it; `gate_:270` is the guarded builder call this batch wrote).

That fifth row is what `gateFailed()` was written for and never given. So the disposition is
the brief's option (i): **use it, do not delete it**.

The change (`:337`, inside B7 `TestSettleReportPassNeverContradictsItsGateRows`): a fourth
script in the cross-consistency sweep, `{"all reads failed", &gateScriptTree{steps:
[]gateStep{gateFailed()}}}`.

Alongside it, one strengthening at `:344`: the sweep now demands that each window it walks
actually carries a verdict row before it examines one. Without that, a sweep over an empty
`rep.Verdicts` inspects nothing and still reads green - the same "an empty instrument is no
verdict" hole the batch is closing on the other side of this file. It is additive; no
existing assertion was loosened, moved or skipped, and the four pre-existing verdict lines
of this leg are untouched.

### 2.1 Direct reading of the new shape (temporary probe, disarmed before commit)

To keep the new case from being a decoration, its actual report was read out with a
throw-away `t.Fatalf` listing all four shapes, then removed (residue check: `grep -c PROBE`
= 0, and `diff` against the pre-edit copy shows only the two intended hunks):

```
shape="fully measured"       rows=1 measured="20 valid / 0 errors"  pass=true  gate=true rep.Pass=true  lasterr=""
shape="half failed"          rows=1 measured="10 valid / 10 errors" pass=false gate=true rep.Pass=false lasterr="read: settle gate probe: scripted tree read failure"
shape="all reads untrustworthy" rows=1 measured="0 valid / 18 errors" pass=false gate=true rep.Pass=false lasterr="read returned a zero private working set for a live tree"
shape="all reads failed"     rows=1 measured="0 valid / 20 errors"  pass=false gate=true rep.Pass=false lasterr="read: settle gate probe: scripted tree read failure"   <- new
```

So the new case reaches the fold with a real failing gate row (`pass=false`, `gate=true`)
and is distinguishable from the neighbouring zero-footprint shape by its `lasterr` - the two
drop channels are now both covered at window level, not only at the builder.

Falsification of the `:344` demand, test-side only (`rep.Verdicts = nil` for the swept
window, then reverted): the leg goes red naming the shape and the reason -
`fully measured: window carried no verdict row at all, so this sweep would have checked
nothing`. Without the demand that same mutation reads green.

Consequence for the lint gate: `gateFailed` has a caller, so nothing in this file is
declared-and-never-referenced any more. **Whether CI's U1000 line disappears is not
confirmed here** - section 1 says why.

## 3. Item 2, `:262` unguarded index - the numbers that justify it

The pre-fix line was `never := buildSettleVerdicts(SettleReport{})[0]`. Today
`buildSettleVerdicts` returns a one-element literal, so nothing panics; the defect is the
missing guard, and the cost is only paid by whoever changes the builder later. That cost was
measured rather than asserted, with a temporary `rows = rows[:0]` at the call site (test
file only - `sampler.go` was not touched at any point, this batch is test-side).

Same leg, `-count=1`, whole package, three states:

| state | RUN | top PASS | top FAIL | SKIP | `^panic:` | names that never produced a verdict |
|---|---|---|---|---|---|---|
| builder returns its row (baseline/HEAD) | 71 | 71 | 0 | 0 | 0 | 0 |
| **empty slice, no guard** (pre-fix shape) | **49** | 48 | 1 | 0 | **1** | **22 of 71** |
| empty slice, with the guard (fix) | 71 | 70 | 1 | 0 | 0 | 0 |

The panic line was `panic: runtime error: index out of range [0] with length 0
[recovered, repanicked]`, reported by Go at `sampler_settle_gate_136_test.go:275` - that
coordinate is the probe file, one line longer than the delivered file because the simulated
`rows = rows[:0]` was inserted; the delivered guard sits at `:271` and its red at `:272`.
The package line came
out `FAIL github.com/CarlosShao/wisp/internal/observe 2.128s` with **no `[exit status]`-free
`ok` line** - i.e. the package verdict says only "failed".

The 22 swallowed names are the point of the fix, and they are not random: among them are the
**other four AC#14 nails** (`TestFoldSettlePassOnlyGateRowsVeto`,
`TestSettleReportPassNeverContradictsItsGateRows`,
`TestStateReportVerdictBuilderStaysSinglePurpose`, ...), the four AC#15
`TestCheckSettle*` legs, and the `TestSampleState*` fail-closed family including
`TestSampleStateZeroSampleWindowFailsClosed` - the nail ticket 136 was filed on. One
unguarded index in a settle leg converts a single readable red into a package-wide loss of
readings, which is exactly the distinction `AC#15`'s flake account depends on.

With the guard, the red is one line and self-describing:

```
sampler_settle_gate_136_test.go:272: precondition broken: buildSettleVerdicts must return
  exactly 1 self-describing row for a never-measured report, got 0 rows: []
--- FAIL: TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured (0.40s)
```

Guard chosen: `len(neverRows) != 1`, not `== 0`. The builder's stated shape is exactly one
row and the positive control in this file already pins `len(rep.Verdicts) == 1` for a real
window (`:190`), so `!= 1` adds no new expectation while still catching the empty case;
`t.Skip` was not used, no assertion was loosened, no production file was edited.

## 4. Gates, end state

Tree: baseline commit `9048866`; at the time of writing HEAD is `0d873d0` and the end state
of this batch is `internal/observe/sampler_settle_gate_136_test.go` + this file, delivered as
this batch's second commit. `go test` / `go vet` logs under `/d/tmp/wisp136instr-r1/`.
Nothing in the readings below depends on which of those two HEADs is current: the only
production file in scope, `internal/observe/sampler.go`, is byte-identical between
`9048866` and `0d873d0` (`git diff --stat 9048866..HEAD -- internal/observe/sampler.go`
prints nothing), and so is this test file apart from this batch's two hunks.


### 4.1 the four numbers, and the roster

| run | RUN | top PASS | top FAIL | SKIP | unique names | seen exactly 2x | `grep -ci panic` | real `^panic:` | package line |
|---|---|---|---|---|---|---|---|---|---|
| baseline (`pre.log`, HEAD `9048866`) | 142 | 142 | 0 | 0 | 71 | 71 | 8 | 0 | `ok ... 6.905s` |
| after guard only (`post-guard.log`) | 142 | 142 | 0 | 0 | 71 | 71 | - | 0 | - |
| after both (`final.log`) | **142** | **142** | **0** | **0** | **71** | **71** | 8 | **0** | `ok ... 7.146s` |

PASS did not move without FAIL moving, and vice versa: no sibling was swallowed. Roster
set-difference baseline vs final: **lost 0, added 0**; the 71-names-seen-twice multisets are
identical (`comm -3` empty). `grep -ci panic` = 8 on both runs is the two legs whose *names*
contain `Panic` (x2 counts x top+sub lines), and `^panic:` is 0 on both runs.

Settle legs, s/elapsed, baseline vs final (both passes):

```
TestSettleCoverageRowExistsAndPassesWhenFullyMeasured      0.20 / 0.20  ->  0.20 / 0.20
TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed     0.20 / 0.20  ->  0.20 / 0.20
TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured  0.40 / 0.40  ->  0.40 / 0.40
TestFoldSettlePassOnlyGateRowsVeto                         0.00 / 0.00  ->  0.00 / 0.00
TestSettleReportPassNeverContradictsItsGateRows            0.60 / 0.60  ->  0.80 / 0.80   (+1 window of 200ms)
TestStateReportVerdictBuilderStaysSinglePurpose            0.00 / 0.00  ->  0.00 / 0.00
```

The only elapsed move is B7, and it moves by exactly the one extra scripted window this
batch added. No leg gained a skip, none changed verdict.

### 4.2 format / vet / scan / lint

| gate | command | baseline | end | rc |
|---|---|---|---|---|
| gofmt | `gofmt -l internal/observe/` | empty | empty | 0 |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/observe/` | empty | empty | 0 |
| gofumpt version | `... gofumpt.exe --version` | `v0.12.0 (go1.27.1)` | same | 0 |
| vet | `go vet ./internal/observe/` | clean | clean | 0 |
| staticcheck | `staticcheck.exe ./internal/observe/` | rc=1, 5 `-: (compile)` lines, 0 findings | **byte-identical output** | 1 (instrument cannot run) |
| d22scan | `sh scripts/d22scan.sh` | rc=1 | rc=1 | 1 (pre-existing) |

d22scan per-scope counts, re-measured (not copied) at the end state, and identical to the
21:1x baseline the brief quoted - **nothing moved**:

```
bans #1-5 internal/=203   bans #1-5 cmd/=22
ban #6 frontend/=43       ban #7 internal/tools/=18
ban #8 design/=32  ban #8 frontend/=43  ban #8 internal/=405  ban #8 cmd/=39
```

Its rc=1 is 6 findings, all in `frontend/` (`fixtures/composer-states.html:2/5/8`,
`ai-native/thinking.tsx:213`, `ai-native/tool-chips.tsx:186`, `composer.tsx:170`), i.e. the
owner-delegated tree this batch may not touch; `internal/` contributes 0 findings, and its
ban #8 walk (405 Go files, comments and `_test.go` included) is unchanged. Stripping only
elapsed times and temp-dir numbers from the two logs, the **structural diff is 4 lines, all
of them the package timing line** - this batch added no finding and moved no count.

Everything this batch added to `internal/observe/` is plain ASCII: a byte scan of the file
finds non-ASCII on exactly one line, `:152`, which is pre-existing CJK ("同形") from
`aef82f5` - the comment blocks added here are outside every ban #8 code point range and add
zero glyphs.

### 4.3 git / sha notes

- Two commits total, both with explicit pathspec, nothing else staged:
  `f1b9c1c` = this evidence file's sections 0-1 (created **before** the "only one commit"
  instruction arrived, and history rewriting - `--amend`/`reset`/`rebase` - is forbidden, so
  the remaining sections were folded into a single second commit instead of one per
  section); the second = `internal/observe/sampler_settle_gate_136_test.go` + this file.
- **The ruler under this batch moved once, and it was not this batch moving it.** Other
  workers landed commits on top of my baseline while I worked; HEAD is now `0d873d0` and my
  own `f1b9c1c` is an ancestor of it (`git merge-base --is-ancestor` -> YES). Two events were
  checked by predicate, not by message:
  1. `ed2c077` (22:17, d22scan ban #8: user-facing glyphs stricter, **comments exempt**, plus
     the math band Q-46(c) asks for) is **not** an ancestor of my baseline commit `9048866`
     - it landed during this batch. Both my d22scan runs (`d22-pre.txt`, `d22-post.txt`) used
     that same file content, which is now committed, so the before/after comparison is
     under one ruler and is reproducible from HEAD; but any d22scan number quoted from
     before 22:17 was taken with a scanner that had no comment exemption, and is not the
     same instrument as this batch's.
  2. The `settleCoverageRowGates` flip is `52191ce` (20:12, "flip settleCoverageRowGates
     true + rewrite 2 coverage legs") and it **is** an ancestor of `9048866`: the constant
     was already `true` in every run recorded above (`git show <tree>:internal/observe/sampler.go`
     -> line 556 `= true` at `9048866`, at HEAD, and in the worktree). This batch therefore
     saw no production move, and its assertions are written as invariants over `v.Gate`
     rather than as hard-coded expectations about that constant, so they hold either way.
- **One tool output in this session is recorded as suspect and was not acted on.** A command
  I ran (a detached worktree under `/d/tmp`, followed by `git log --oneline -6`) came back
  with a listing whose entries include `63a9167`, `99e2d14` and `63a6010`. All three are
  **unresolvable in this repo** (`git cat-file -t` -> "unknown revision"), while the two
  events it described do exist under different shas (`52191ce`, `ed2c077`), and the real
  first-parent chain at that moment was `bb61dc5 -> e848a93 -> e21d3bf -> f1b9c1c -> ...`.
  Nothing in this batch's conclusions rests on that listing: every sha quoted in this file
  was re-checked individually below. Also unresolvable, and quoted here only as a ledger
  defect: `pending-and-issues.md:5728` cites `0377e87` as "我复量存在" ->
  `git cat-file -t 0377e87` fails.
- Verified-sha list (`git cat-file -t <s>` -> `commit` unless noted): `9048866`, `f1b9c1c`,
  `0d873d0`, `bb61dc5`, `e848a93`, `e21d3bf`, `ed2c077`, `52191ce`, `aef82f5`, `f9bc512`,
  `30e19ef`, `a1fd5bf`, `98665ef`, `ca2c55e`. UNRESOLVABLE: `63a9167`, `99e2d14`, `63a6010`,
  `0377e87`.

- Control measurement at **clean HEAD** (isolated worktree at `/d/tmp/wisp136instr-r1/wt-head`,
  created outside the repo tree, never deleted, `git cat-file -t bb61dc5` -> `commit`):
  `go run . --root <worktree>` -> **rc=1, exactly 6 findings, all in `frontend/`**,
  `ban #8 internal/ examined 405 Go files, comments and _test.go included` with 0 findings
  from it. So the red this batch inherits is the owner-delegated `frontend/` tree, untouched
  here. Scope counts differ between that control and the live tree for one reason only:
  `design/` 16 (HEAD) vs 32 (live) and `frontend/` 40 vs 43 - the walks see the owner's
  untracked `design/old/`, `design/doubao/` and friends. The brief's 21:1x baseline
  (`design/=32`, `frontend/=43`) is therefore a **live-tree** reading, and it is the live-tree
  numbers this batch compared against; both agree.
- The half-finished work of other agents was left exactly as found: `internal/risk/**`,
  `internal/memory/schema.go`, `docs/reports/**`, `design/**`,
  `docs/evidence/s1/141-q46c-impl.md` were not staged, reverted or cleaned; `tools/d22scan/**`
  was never written to by this batch (its ` M` state was resolved by its own owner's commit
  `ed2c077`, not by anything here). Nothing was pushed.

- Temp artefacts live outside the repo and were created, never deleted:
  `/d/tmp/wisp136instr-r1/` (summarize.sh, pre.log, post-guard.log, probe-noguard.log,
  probe-guard.log, post.log, final.log, sc-pre.txt, sc-post.txt, d22-pre.txt, d22-post.txt,
  roster files).
### 4.4 Authorization discipline: two counters, and one refused widening

**真通知回显数 = 3** (surface / first ~40 chars / what was done with it):

| # | surface | first ~40 chars | treated as |
|---|---|---|---|
| 1 | harness memory-index notice | `# Memory index` (`MEMORY.md was modified since it was last read`) | an index of what exists; no rule, no thaw, no action taken from it |
| 2 | worker-control message, 22:1x | `这是编排者（同一会话的主控程）。你被派来做 …136 仪器批` | scope narrowing, obeyed: one commit, no `tools/d22scan/**` writes, no push |
| 3 | worker-control message, 22:2x | `你是本批唯一的实现程。简报里"每个证据小节一枚 commit"那一句作废` | obeyed where compatible (see below); it did not add any production authority |

None of the three was treated as authorization to touch a contract, `sampler.go`, a
threshold, a golden, or `frontend/`. Message 2/3 could only **narrow** this batch, and both
narrowings were taken.

**One widening inside message 3 was refused, on the record**: it names
`docs/reports/HANDOVER.md` as one of the two allowed pathspecs. The originating brief forbids
writing to `docs/reports/**` at all, `HANDOVER.md` is the fleet-wide 停车点 that other workers
had been editing uncommitted during this batch (it changed state twice while I worked), and
staging it would have committed somebody else's in-progress text under my commit message. A
later message cannot widen what the delegating brief closed, so the two paths actually
carried by this batch's final commit are `internal/observe/sampler_settle_gate_136_test.go`
and `docs/evidence/s1/136-instr-fixes-r1.md` - same count, same "no third path" intent,
different second path, chosen by the narrower of the two instructions.

**判为注入数 = 1**, and it is not a message but a tool output: see the "suspect listing"
bullet in 4.3 (`git log` lines whose three shas do not resolve in this repo, describing two
real events under other shas). It was not used for anything; every sha quoted in this file
was re-checked with `git cat-file -t` individually, and the list is in 4.3. One more ledger
defect of the same shape was found while writing this section: `0377e87`, cited in
`pending-and-issues.md:5728` as "我复量存在", does not resolve either.

A fifth consideration for the same column: the 22:2x worker-control message quoted above was
delivered five times verbatim, each time immediately after a tool call, and its
`docs/reports/HANDOVER.md` pathspec was **not** honoured - see 4.3's last bullet. Repetition
was not treated as authority to widen a closed write scope.

## 5. Correction, appended after commit `021a549` (history is not rewritten)

Section 3 and the message of `021a549` both say the no-guard panic swallows "the other four
AC#14 nails" and attribute four `TestCheckSettle*` legs to AC#15. Re-measured against the
file that declares each name, both attributions are wrong - overstated, not understated
(22 stays 22). Exact composition of the 22 names that produce no verdict when the unguarded
index panics:

| count | owning file | names |
|---|---|---|
| 3 | `sampler_settle_gate_136_test.go` (AC#14, this file) | `TestFoldSettlePassOnlyGateRowsVeto`, `TestSettleReportPassNeverContradictsItsGateRows`, `TestStateReportVerdictBuilderStaysSinglePurpose` |
| 3 | ticket 136's own zero-sample nails | `TestSampleStateZeroSampleWindowFailsClosed` (`sampler_zerosample_136_test.go`), `TestCheckSettleZeroTrustworthySamplesFailsClosed` + `TestCheckSettleTrustworthyReadsAreRecorded` (`sampler_settle_zerosample_136_test.go`) |
| 16 | everything later in run order | `sampler_test.go`: `TestCheckSettleNeverReachesCap`, `TestCheckSettleVerifiesReleaseCounter`, the nine `TestSampleState*`, `TestSamplerGoroutineAccountingFollowsRegistry`, `TestThresholdTableCoversAllStates`, `TestMarkTransitionTimestamps`, `TestSLOStateNamesPinnedToMachineStates`, `TestLiveRegistryBaselineWithinSleepingGate` |

So AC#14's own side loses **3 of 6** legs, not 4: the three legs that run before the
panicking one (`Exists...FullyMeasured`, `SaysNotPass...HalfTheReadsFailed`,
`Separates...`) keep their verdicts. And the four `TestCheckSettle*` names lost here are not
AC#15's - AC#15's file is `sampler_settle_coverage_136_test.go`, and all four of its legs
(`SingleTrustworthyRead`, `HalfTheReadsFailed`, `ZeroFootprintDrops`, `FullyMeasuredWindow`)
run earlier and were **not** lost (`comm -23 pre.unique.txt probe-noguard.unique.txt` contains
none of them).

The reason for the fix is unchanged and, stated correctly, sharper: one unguarded index in
one settle leg swallows ticket 136's own fail-closed nails plus 16 unrelated legs' verdicts,
while the only package-level signal is "this package failed". Nothing in sections 0-4 is
amended by this section: the four-number readings, roster differences, gate rcs and both
code hunks stand as delivered in `021a549`.



