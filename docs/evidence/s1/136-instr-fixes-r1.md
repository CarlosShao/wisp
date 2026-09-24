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
