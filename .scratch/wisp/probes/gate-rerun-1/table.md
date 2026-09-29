# gate-rerun-1 — ticket 234: 12 cells that only need a reading (mutation + gate rerun on current HEAD)

- Leg: `gate-rerun-1` (non-implementor rerun). Ticket: `.scratch/wisp/issues/234-the-twelve-cells-that-only-need-a-reading-mutation-and-gate-rerun-on-current-head.md`
- Anchor self-taken at start: `git rev-parse --short HEAD` = `015f6be1`, branch `dev`.
- Blocker check: ticket face Status requires ticket 222's acceptance leg to finish first.
  Evidence it finished: `d39f730e` (222-v1 shou-gong), `f23586fa` (222-v1 table + raw readings),
  `015f6be1` (ledger A447 collects 222-v1). Block released before this leg started.
- Scope: zero production code. `frontend/**` and `design/**` never read, never listed, never quoted.
- Order of work: 5 mutation cells -> 7 gate/regression cells -> AC#5 ledger relocation -> AC#6 bad yardstick.
  One commit per cell so that a mid-run death only loses the remaining cells.

## Baseline roster (AC#1)

- VERDICT: baseline established, all green on `015f6be1`.
- Yardstick (verbatim from the ticket face): `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/winsec/ ./internal/risk/ ./internal/config/ ./internal/agent/approval/ -count=1 -v`.
  This leg added `-p 1` (packages run sequentially) on purpose: `internal/risk TestResolvePerCallBudget` is the registered parallel-load-sensitive case (`R-116-1`), and `-p 1` removes the contention variable instead of having to re-litigate it later.
- READING: `go.exe` count before run = 0. rc=0. Per package (tab-separated `ok` lines): winsec 11.966s / risk 4.025s / config 0.892s / agent/approval 0.346s.
  Four numbers from `-v`: RUN=471 PASS=469 FAIL=0 SKIP=2.
- SKIP named (both are documented, neither is this ticket's): `TestSyncRegistryProbeLive` (no registry-grade sync record on this machine) and `TestDefaultDeadlineWallClockMeasurement` (intentionally slow 300s wall-clock, only under `WISP_84_MEASURE=1`).
- Name-set roster: `.scratch/wisp/probes/gate-rerun-1/ac1-baseline-roster.txt` = 297 top-level names (295 PASS + 2 SKIP). Raw `-v` output: `ac1-baseline-v.txt` (268,251 bytes).
- YARDSTICK DEFECT found while measuring (registering, not using it): the four numbers must NOT be counted with `grep -c '^--- PASS'`.
  Go prints every `=== RUN` at column 0 including subtests, but subtest results are indented `    --- PASS:`. So `^--- PASS` = 295 while the true PASS count is 469 (295 top-level + 174 subtests).
  Correct forms: `grep -c '^=== RUN'` and `grep -c '^ *--- PASS'`. RUN(471) = PASS(469) + SKIP(2) + FAIL(0) reconciles; the naive form under-reports PASS by 174 and would make any "four numbers" reading look wrong.
  Also the dispatch's `grep -P '^FAIL\t'` rule was confirmed in practice: `grep -c '^FAIL'` on this tree's red output double-counts each failing package.
- Historical registered reds (`internal/panel` 4 + `internal/ball` 1) are outside AC#1's package scope, so they do not appear in this baseline; they are handled per cell below.

## Mutation cells (5)

### M-1 ticket 92 `:68` — AC#5 direction (iii), verifier marked [self-report only]

- VERDICT: PASS (has a reading now, no longer self-report only). Upgraded from [self-report only].
- Criterion (ticket 92 `:68`, read this leg): "(iii) 把 reparse 那次的拒绝改成放行 ⇒ **既有**安全用例红（不是本票新写的）" and "编译失败不算变异".
- Anchor read this leg: `internal/winsec/placement_windows.go:59-69` — the walk over `pathPieces(path)` whose refusal is at `:65-68`
  (`if isReparsePoint(prefix) { return "", fmt.Errorf("%w: %s traverses a reparse point at %s ...") }`). The file's own doc comment `:23-25` names this the junction case.
- Mutation: `if isReparsePoint(prefix) {` -> `if false && isReparsePoint(prefix) {` (refusal becomes pass-through; the predicate is still called, so nothing fails to compile).
  Landing proof (`mutate.py` stdout, `m1-*`): needle hits in orig = 1, in mutant = 0, repl hits in mutant = 1, line count unchanged 71 -> 71. `go build -overlay m1-overlay.json ./internal/winsec/` rc=0, so this is a mutation and not a compile failure.
- Instrument: `go test -overlay .scratch/wisp/probes/gate-rerun-1/m1-overlay.json ./internal/winsec/ ./internal/risk/ -count=1 -v -p 1` -> rc=1, `FAIL\t...internal/winsec 8.782s`, `internal/risk` still ok. Zero `0xc0000135` hits in the log (checked, so the red is a real assertion and not a load failure).
- RED roster (7 entries = 2 top-level + 5 subtests), file `m1-red-names.txt`:
  `TestAC3JunctionInputIsRefusedNotSealed` (+ subtest `/with_only_the_built-in_verifier`),
  `TestAC3PlacementFloorHoldsForEverySeparatorSpelling` (+ `/all-forward-slash` `/mixed-separators` `/native-backslash` `/trailing-separator`).
- THE PART THE CRITERION ACTUALLY TURNS ON ("not written by this ticket"): both carriers are **pre-existing** cases —
  `internal/winsec/resolve_windows_test.go:61 TestAC3JunctionInputIsRefusedNotSealed` and `placement_symlink_113_other_test.go`-family spelling coverage —
  neither is in ticket 92's own new files. So direction (iii) holds on current HEAD.
- Restore proof: overlay never writes the worktree; `git status --porcelain -- internal/ cmd/wisp` = **empty** (printed in the same run).
  Re-ran the 7 red names with no overlay -> `m1-restored-green.txt`: rc=0, `ok`, PASS=13 FAIL=0, and all 7 names are green in the AC#1 baseline roster too.
- What is still not covered by this cell: (i) and (ii) of the same AC#5 were judged by the acceptance table itself and stay as they were; this leg only supplies (iii). `R-92b-4` is carried under G-1.

### M-2 ticket 104 `:52` — AC#3 both directions

- VERDICT: PASS on all three legs (criterion says "双向变异" and lists three legs ①②③; all three were re-run).
- Criterion (ticket 104 `:52-55`, read this leg): ① detection back to "explicit ACE only" ⇒ AC#1 red; ② private set judged by **name string** instead of **SID** ⇒ some case red; ③ WARN "every time" ⇒ AC#2 red. Plus "`go build` rc=0 先量到（编译失败不算变异）".
- Anchor read this leg: `internal/winsec/winsec_windows.go:164-198 foreignPrincipals` (buckets at `:191-195`, SID skip at `:174`) and `:311-318` (the `noticeNarrowed` call guarded by `len(explicit) > 0 || len(inherited) > 0`). Doc `:163` states "Judgment is by SID throughout".
- Three mutations, all via `-overlay`, worktree never written. Landing proof in each: needle hits 1 -> 0, repl hits 1, 716 lines -> 716.
  - ① `if ace.inherited {` -> `if false && ace.inherited {` at `winsec_windows.go:191`. `go build` rc=0. `go test -overlay ... ./internal/winsec/ -count=1 -v -p 1` rc=1 -> **4 top-level reds**: `TestAC1SealFileReportsTheInheritedGrantItCleared`, `TestAC1DefaultLogSaysInherited` (both are AC#1's own carriers, exactly what leg ① demands), plus `TestAC118KindFieldSaysBoth...`, `TestAC118KindFieldSaysInherited...`.
  - ② `if ace.grant && set[ace.trustee] {` -> `... set[ace.text] {` at `:174` (whitelist compared by the OS's name rendering instead of the SID). `go build` rc=0 -> rc=1, **10 top-level reds** incl. `TestSealNarrowsAndNamesThePrincipalItRemovedBySID`, `TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem`, `TestSealReportsThePrincipalsItCleared`, `TestNoticeAttributionKeepsTwoTreesApart` (`m2b-red-names.txt`). Leg ② is the "silently fails on a localized machine" guard and it does bite.
  - ③ `if beforeErr == nil && (len(explicit) > 0 || len(inherited) > 0) {` -> `if beforeErr == nil {` at `:311` (report on every seal). `go build` rc=0 -> rc=1, **4 top-level reds**: `TestAC2InheritedNoticeHasANoiseBound` (AC#2's carrier, as the leg demands), `TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem`, `TestSealReportsThePrincipalsItCleared`, `TestNoticeAttributionKeepsTwoTreesApart`.
- Restore proof: `git status --porcelain -- internal/ cmd/wisp` = empty after each leg; the union of all red names re-run without overlay -> `m2-restored-green.txt` rc=0, `ok`, PASS=14 FAIL=0.
- Reading files: `m2a-red.txt` `m2b-red.txt` `m2c-red.txt` + the three `*-red-names.txt` + `m2-restored-green.txt`, all in this directory.

### M-3 ticket 110 `:39` — AC#3 the new step must go red by itself

- UNJUDGED

### M-4 ticket 113 `:54` — AC#4 two shapes

- UNJUDGED

### M-5 ticket 115 `:58` — AC#4 two shapes

- UNJUDGED

## Gate / regression cells (7)

### G-1 ticket 92 `:73` — AC#6 three readings valid only against two old tree snapshots

- UNJUDGED

### G-2 ticket 97 `:55` — AC#5 per-package five numbers re-taken on current HEAD

- UNJUDGED

### G-3 ticket 104 `:57` — AC#4 four packages `-count=2 -v`

- UNJUDGED

### G-4 ticket 105 `:54` — AC#5 five readings + per-scope file counts

- UNJUDGED

### G-5 ticket 110 `:47` — AC#5 `bash -n` / `d22scan.sh` per scope not lower / new step ordering

- UNJUDGED

### G-6 ticket 113 `:56` — AC#5 in-container three packages four numbers + per-package vet + gofmt + `d22scan.sh`

- UNJUDGED

### G-7 ticket 115 `:64` — AC#6 per-package four numbers / gofmt / gofumpt.exe -l / go vet / d22scan

- UNJUDGED

## Old debts carried by 5 of the cells

- UNJUDGED

## Final roster comparison (AC#4)

- UNJUDGED

## AC#5 — 10 duplicate denominators relocated to still-open tickets

- UNJUDGED

## AC#6 — retire one bad yardstick (`NewGate(` zero hits)

- UNJUDGED

## Not finished / not judgeable

- (must not be left empty at the end of the leg)
