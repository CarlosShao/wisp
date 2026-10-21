# 136-r1 STOP REPORT — the premise of this leg is falsified by measurement

**Leg:** 票 136 落地腿 136-r1 (writer). **Date:** 2026-10-21 09:2x +08. **Anchor HEAD:** `1d138a0b`.
**Verdict: STOPPED before writing any new case.** The premise this leg was dispatched on
("the zero-sample fail-closed guard has no test today") does **not** hold on HEAD.
Thirteen mutations were landed one at a time and **every one of them is caught today**.
Writing another nail would be a tautology (forbidden by the dispatch) and a false green.

Nothing was written to production code. `internal/observe/sampler.go` is byte-identical to HEAD
(sha256 `04A9ACD8EC8B046335DB7E110D4610B00A4FD8D207EACAA89033B368F20AE2F8`, `git diff` empty).

## 1. The two premise items that do not match HEAD

| # | Dispatch / ticket says | Measured on HEAD `1d138a0b` |
| --- | --- | --- |
| P1 | the zero-sample fail-closed guard "has no test today"; removing it leaves the package green (ticket line 19: "54 tests rc=0") | **False.** The guard is nailed in both functions. Removing or weakening it turns the package red today (13 of 13 mutations, section 3) |
| P2 | "the existing `TestCheckSettle*` two cases" | **8 cases**, not 2. `sampler_test.go` has the original two; ticket 136 itself added six more across `sampler_settle_zerosample_136_test.go` and `sampler_settle_coverage_136_test.go` |
| P3 | ticket AC#9 is open work for this leg | **AC#9 is already checked `[x]`** in the ticket, flipped 2026-09-24 00:0x by the orchestrator on acceptor table `136-ac8-ac9-r1-acceptance.md`. The nail it asked for is commit `4fc65dd5` / `d8390aaf`, both present at HEAD |

Ticket name note: the dispatch quotes a slug that does not exist on disk
(`...-has-no-test.md`). The real file is
`136-the-zero-sample-fail-closed-guard-underneath-the-freshness-nail-has-no-nail-54-tests-stay-green-without-it.md`.
The `54` in that real name is itself the stale number: the package is **71** top-level cases today.

## 2. Verbatim, content-anchored (no line numbers trusted across anchors)

StateReport guard, today at `sampler.go:332-340`:

```
	if len(rep.Samples) == 0 {
		// No trustworthy sample in the whole window: the report fails with
		// its own gate verdict rather than silently passing on zeros.
		rep.Verdicts = append(rep.Verdicts, Verdict{
			Metric: "sampling", Measured: fmt.Sprintf("%d valid / %d errors", len(rep.Samples), rep.SampleErrors),
			Limit: ">0 valid samples", Pass: false, Gate: true,
			Note: "fail-closed: unmeasurable windows never pass (D22 mode-6 guard)",
		})
	}
```

CheckSettle drop branch (the AC#9 target), today at `sampler.go:505-510`:

```
		} else if m.PrivateWorkingSetBytes <= 0 {
			// Zero-footprint reads of a live tree are untrustworthy (see
			// SampleState) and are dropped, never recorded as progress.
			rep.SampleErrors++
			rep.LastSampleError = "read returned a zero private working set for a live tree"
		} else {
```

## 3. Mutation self-proof: 13 landings, 13 catches, 0 survive

Every landing was proven before reading: content-anchored `grep -n` of the mutated line
plus `go build ./...` rc=0, then `go test -count=1 -v ./internal/observe/`. Each mutation was
reverted alone (never stacked) and the file sha re-checked to the HEAD value before the next one.

| id | mutation | red names (verbatim) | four |
| --- | --- | --- | --- |
| M1 | settle drop branch short-circuited (`else if false &&`) | `TestCheckSettleZeroFootprintDropsAreCountedToo`, `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured`, `TestCheckSettleZeroTrustworthySamplesFailsClosed` | 71/68/3/0 |
| M4 | **AC#9's own named relaxation**: settle `<= 0` to `< 0` | same three | 71/68/3/0 |
| M10 | settle: drop kept, `rep.SampleErrors++` removed | `TestCheckSettleZeroFootprintDropsAreCountedToo`, `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | 71/69/2/0 |
| M11b | settle: count kept, `rep.LastSampleError` attribution removed | `TestCheckSettleZeroFootprintDropsAreCountedToo` | 71/70/1/0 |
| M13 | settle coverage row loses its zero-sample clause | `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | 71/70/1/0 |
| M2 | StateReport zero-sample block removed | `TestSampleStateZeroSampleWindowFailsClosed` | 71/70/1/0 |
| M3 | guard row kept, `Pass: false` to `Pass: true` (tooth removed) | `TestSampleStateZeroSampleWindowFailsClosed` | 71/70/1/0 |
| M5 | SampleState drop branch `<= 0` to `< 0` | `TestSampleStateZeroSampleWindowFailsClosed` | 71/70/1/0 |
| M7 | StateReport Gate fold removed (`Gate && !Pass` no veto) | `TestSampleStateSleepingDiskWriteGateFails`, `TestSampleStateSleepingTCPGateFails`, `TestSampleStateHandleGateUsesRulingLimit`, `TestSampleStateLeakFixtureFlipsRed`, `TestSampleStateZeroSampleWindowFailsClosed` | 71/66/5/0 |
| M8 | zero-sample row `Gate: true` to `Gate: false` | `TestSampleStateZeroSampleWindowFailsClosed` | 71/70/1/0 |
| M9 | guard narrowed to `len(rep.Samples) == 0 && rep.SampleErrors == 0` | `TestSampleStateZeroSampleWindowFailsClosed` | 71/70/1/0 |
| M11a | SampleState: count kept, attribution removed | `TestSampleStateZeroSampleWindowFailsClosed` | 71/70/1/0 |
| M12 | SampleState: drop kept, `rep.SampleErrors++` removed | `TestSampleStateZeroSampleWindowFailsClosed` | 71/70/1/0 |

Four numbers are `RUN` / `PASS` / `FAIL` / `SKIP`, each counted separately from `-v` lines,
never subtracted. `SKIP` is 0 in every row; no case was skipped and none was relaxed.

The core claim the dispatch asked me to fail: **M4 is exactly the mutation ticket AC#9 names
("relax `> 0` to `>= 0`"). On HEAD it reddens three cases by name.** There is no gap left to fill.

## 4. The three readings the dispatch asked for

- **Pristine (original) green**, before any mutation:
  `rc=0`, `RUN=71 PASS=71 FAIL=0 SKIP=0`, verbatim `ok  github.com/CarlosShao/wisp/internal/observe 3.790s`
- **Mutated red** (M4, the AC#9-named relaxation): `rc=1`, `RUN=71 PASS=68 FAIL=3 SKIP=0`,
  red names as in section 3. Same shape for M1.
- **Restored green**, after reverting the last mutation:
  `rc=0`, `RUN=71 PASS=71 FAIL=0 SKIP=0`, verbatim `ok  github.com/CarlosShao/wisp/internal/observe 3.761s`,
  and `sha256(sampler.go)` back to the HEAD value above with `git diff` empty.

## 5. The two original `TestCheckSettle*` cases, by name, today's colour

Run narrow with `-run 'TestCheckSettle'`; all eight matched for completeness:

```
--- PASS: TestCheckSettleSingleTrustworthyReadReportsItsLoss (0.10s)
--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)
--- PASS: TestCheckSettleZeroFootprintDropsAreCountedToo (0.10s)
--- PASS: TestCheckSettleFullyMeasuredWindowReportsNoLoss (0.10s)
--- PASS: TestCheckSettleZeroTrustworthySamplesFailsClosed (0.10s)
--- PASS: TestCheckSettleTrustworthyReadsAreRecorded (0.10s)
--- PASS: TestCheckSettleVerifiesReleaseCounter (0.20s)
--- PASS: TestCheckSettleNeverReachesCap (0.08s)
```

Four numbers for that narrow run: `RUN=8 PASS=8 FAIL=0 SKIP=0`, `rc=0`.

## 6. Gate readings, run by me on the restored tree

| gate | reading |
| --- | --- |
| `GOFLAGS= go build ./...` | `rc=0`, no output |
| `go vet ./internal/observe/` | `rc=0`, no output |
| `go test -count=1 ./internal/observe/` | `rc=0`, `ok  github.com/CarlosShao/wisp/internal/observe 3.796s`; four numbers from the `-v` form above `71/71/0/0` |
| `go test -count=2 -v ./internal/observe/` | `rc=0`, `RUN=142 PASS=142 FAIL=0 SKIP=0 panic=0`, 71 distinct names, no oscillation, no name red in either pass |
| `sh scripts/d22scan.sh` | `rc=0`, `d22scan: clean - no D22 ban violations`; scopes `244/42/116/25/39/116/608/143` |
| `sh scripts/check-path-length-budget.sh` | `rc=1`, `VERDICT RED` — **not mine**, see section 7 |

## 7. One pre-existing red I must not absorb or hide

`sh scripts/check-path-length-budget.sh` returns `rc=1` on HEAD. Both red lines point at
`.scratch/wisp/issues/139-context-compression-...-done.md` (name 109 chars, relative 130),
failing in both directions of the roster pair:

- (A) over budget and not in the roster, and
- (B) roster entry at `scripts/check-path-length-budget.sh:376` that is no longer an over-budget path.

The file is tracked at HEAD and last touched by `718ddb3f`, not by me. My change surface in this
leg is empty, so this red **pre-exists** and belongs to the 139 leg now in flight. I did not touch
the roster, the script, or that ticket — absorbing it would be exactly the "relax the gate to get
green" move the dispatch forbids.

`sh scripts/d22scan.sh` also refuses to run at all when `GOTMPDIR` sits on volume D
(`runtests.sh` ticket 311 AC#1 shape B, `rc=2`). That is a guard working as designed, not a failure;
the `rc=0` above was taken with the temp vars unset so the root lands on volume C.

## 8. What I did not do, and why

- **Did not add a new case.** The behaviour the dispatch wants pinned already is pinned, by name,
  reversibly. A fourteenth case asserting the same thing would pass under every one of the 13
  mutations above only because the older cases also catch them, and would add no discriminating
  power. That is the恒真 shape the dispatch bans.
- **Did not touch production code.** No解冻 was needed or used. `sampler.go:470-501` and
  `:431-443` were read only.
- **Did not touch `ci.yml`,** the threshold surface, `thresholds.go`, any golden, or `SampleState`
  semantics. Did not skip or relax any existing case.
- **Did not commit anything.** There is no implementation to commit. This file is the only artifact;
  leaving it uncommitted keeps the tree honest about that.

## 9. What the orchestrator should decide

1. **Ticket 136 AC#9 needs no more work.** It is already `[x]`, and section 3 re-proves it
   independently on HEAD. If the leg was meant to re-open it, the ticket and the dispatch disagree.
2. **Ticket 136's own title is now stale.** It still says "54 tests stay green without it"; the real
   number is 71 and the guard does not stay green when removed. Retitling needs the human-approved
   path, not an implementer edit — I did not touch it.
3. If a genuinely un-nailed corner of this family exists, it is **outside `internal/observe`**:
   ticket AC#10's `cmd/wisp` end-to-end (`run.Pass` feeding the exit code) is the one seam I did not
   measure, and `cmd/wisp` is outside this leg's write ceiling. Say the word and I will measure it
   read-only and report, without writing.
