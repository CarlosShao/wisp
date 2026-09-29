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

- VERDICT: PASS, both shapes. Upgraded from [self-report only].
- Criterion (ticket 110 `:39-41`, read this leg): break a winsec security assertion in a `/tmp` snapshot (the face's own example: "私有集改成按名字比") ⇒ the new step must be rc≠0;
  and prove it is not an empty instrument (package list without winsec ⇒ a "扫描空=红" guard or an explicit failure). Each mutation: same-chain `grep -n` landing proof, `go build` rc=0 first.
- Snapshot: `git archive HEAD | tar -x -C /tmp/gate-rerun-1-snap` (outside the repository, per A38④). Repo worktree never touched; `git status --porcelain -- internal/ cmd/wisp` = empty at the end.
- [0] Positive control on the pristine snapshot: `bash scripts/winsec-tests.sh` rc=**0**, and it printed its own top-level result line
  `ok github.com/CarlosShao/wisp/internal/winsec 7.778s` / `RUN=101 PASS=58 FAIL=0 SKIP=0` (`/tmp/gr1-m3-pristine.log`).
- [A] Mutation = the ticket's own example, reused from M-2 leg ②: `winsec_windows.go:174` `set[ace.trustee]` -> `set[ace.text]` (private set judged by name string instead of SID).
  Landing: `grep -c 'set[ace.text]'` = 1, `set[ace.trustee]` = 0 in the snapshot. `go build ./internal/winsec/` rc=0 **inside the snapshot** (so this is a mutation, not a compile failure).
  `bash scripts/winsec-tests.sh` -> rc=**1**, `FAIL ... internal/winsec 8.325s`, and **zero** `0xc0000135` hits (`dll_load_failures=0`), so the step really executed the binary.
  10 named top-level reds (`m3-mutA-red-names.txt`), the same roster as M-2 leg ② — what is judged red here is the STEP, not just the test file.
- [B] Not-an-empty-instrument proof: `bash scripts/winsec-tests.sh ./internal/config/` -> rc=**2** with the explicit GUARD 1 text
  (`scripts/winsec-tests.sh:60-66`): "the package list for this step does not name ./internal/winsec/; it is [./internal/config/]" (`/tmp/gr1-m3-mutB.log`).
  That is the "显式失败" branch the criterion accepts: a removed package fails the step instead of shortening the run.
- Restore proof: snapshot's `winsec_windows.go` copied back and `diff -q` against the repo file = identical (`snapshot_identical_to_repo=yes`).
- Caveat registered, NOT fixed (zero production code in this ticket): `scripts/winsec-tests.sh:96-100` computes its four numbers with `grep -c '^--- PASS'`,
  which by the AC#1 finding counts top-level results only (58) while its `RUN` (101) includes subtests — the step's printed four numbers mix two granularities.
  It does not affect the step's rc (verdicts are delegated to `scripts/portable-tests.sh` -> `tools/d22scan/runtests.sh`, and guard 2 anchors on the `^ok|FAIL` package line).
- Not re-run here: guard 2's "a build tag empties the whole package" shape (the criterion's empty-instrument branch is satisfied by guard 1).

### M-4 ticket 113 `:54` — AC#4 two shapes

- VERDICT: PASS, both shapes, measured in a REAL linux container (not `GOOS=linux go vet`).
- Criterion (ticket 113 `:54`, read this leg): turn the new leg off ⇒ AC#1's case must be red; then try the "ancestor chain consulted only one level" half-fix shape ⇒ must also be red
  (proving it bites the whole set and not one line). The acceptance table's own four-shape reading is `docs/evidence/s1/113-adversarial-acceptance.md:16` (MUT-A 9 red / MUT-B 1 red / MUT-C 9 red / MUT-D2 2 red).
- Anchor read this leg: `internal/winsec/winsec_other.go` (`//go:build !windows`) — `:155-163 platformVerifyPlacement` walks `pathPieces(path)` and refuses at `:157-160`; `:188-194 ancestorIsLink` is the predicate.
- Where it ran: `git archive HEAD | tar -x -C /tmp/gate-rerun-1-snap` mounted at `/work` in `golang:1.27` (go1.27.1 linux/amd64).
  False-green guard fired FIRST and passed: `/work/go.mod` exists, `go_files_under_internal=460`, `posix_test_files=7`, `winsec_go=32` — so the mount carried the real tree, and the bind mount was NOT the silent-empty shape the dispatch warns about.
  (One attempt did fail, for a different reason: Git Bash rewrote `-w /work` into `C:/Users/.../bin/git/work` and docker exited 125. Fixed with `MSYS2_ARG_CONV_EXCL='*'`; rc=125 was recorded, never read as a green.)
- M-4a leg off: `if ancestorIsLink(prefix) {` -> `if false && ancestorIsLink(prefix) {`. Landing hits=1. `go build ./internal/winsec/` rc=0. `go test -count=1 -v -p 1` -> rc=1, `FAIL ... 0.077s`, **11 named top-level reds** (`m4-g6-container.log`):
  the three `TestAC1POSIX...` AC#1 carriers (`TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing`, `TestAC1POSIXSealDirThroughASymlinkRefuses`, `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`, `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`, `TestAC1POSIXSealFileThroughABackslashNamedLink`, `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`) plus the 118/119/125 POSIX legs. So AC#1's cases do go red.
- M-4b half-fix (one level only): `for _, prefix := range pathPieces(path) {` -> `... pathPieces(path)[:1] {`. Landing hits=1. `go build` rc=0. -> rc=1, **same 11 named reds**.
  That is the stronger of the two readings: a walk that consults only the first prefix is caught by the全集-shaped cases, not by one line — notably `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`.
- Restore proof: `cp -p` of the pristine copy back, `diff -q` = identical after each shape (`restored=identical` twice), and the container's final check `winsec_other_identical_to_pristine=yes`.
  Repo worktree: `git status --porcelain -- internal/ cmd/wisp` empty (checked in the host runs); nothing in the repo was ever mutated.

### M-5 ticket 115 `:58` — AC#4 two shapes

- VERDICT: PASS on both shapes. Upgraded from [undecidable: only the implementor's own commit body + zero tables under 115].
- Criterion (ticket 115 `:58-59`, read this leg): revert the chosen fix to its original state ⇒ AC#3's new case must go red; then try the "only case-insensitive comparison (`EqualFold`)" half-fix shape ⇒ must also be red
  (proving it bites the 8.3-alias family rather than one literal string).
- Anchor read this leg: `internal/winsec/winsec_windows.go:106-112 noticeNamesTree` — the fix is `ResolvePath(spelling)` then `sameTree(n.Path, resolved.String())` at `:111`. Doc `:117-119` names the three comparison faces the old shape used (exact `==`, `strings.EqualFold`, `strings.ToLower` map key) and says all three miss an 8.3 alias.
- First attempt at both shapes FAILED TO COMPILE (`.\m5a-mutant.go:107:2: declared and not used: resolved`, `build_rc=1`) and was therefore **discarded, not counted** — the tickets' own rule "编译失败不算变异" applies to me too.
  Re-cut so the local stays used (`_ = resolved`), then: `go build -overlay` rc=**0** for both shapes. Landing proof in both: needle hits 1 -> 0, repl hits 1, 716 -> 717 lines.
- Shape 1 (revert to comparing the caller spelling): `return sameTree(n.Path, resolved.String())` -> `return n.Path == spelling`.
  `go test -overlay m5a-overlay.json ./internal/winsec/ -count=1 -v -p 1` rc=1, `FAIL ... 7.418s`, **2 named top-level reds**:
  `TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree` (AC#3's new positive half — exactly the carrier the criterion names) and `TestNoticeAttributionKeepsTwoTreesApart`.
- Shape 2 (half-fix, `EqualFold` only): `... -> return strings.EqualFold(n.Path, spelling)`. rc=1, `FAIL ... 7.499s`, **same 2 named reds**.
  This is the stronger reading: case-insensitivity does not bridge a long name to its 8.3 short name, so the guard bites the whole alias family, not one string.
- Restore proof: overlay never writes the worktree, `git status --porcelain -- internal/ cmd/wisp` = empty; both names re-run without overlay -> `m5-restored-green.txt` rc=0, PASS=2 FAIL=0 (named-only run,口径 stated).
- Readings: `m5a-red.txt` `m5b-red.txt` `m5a-red-names.txt` `m5b-red-names.txt` `m5-restored-green.txt`, plus the discarded build-failed attempt left on disk as `m5a-mutant.go` / `m5b-mutant.go` (temp files are only ever created, never deleted).
- NOT delivered by this leg, on purpose: ticket 115 still has **zero** tables named `115-*` under `docs/evidence` (`find docs/evidence -name '115*' | wc -l` = 0, measured this leg), and ticket 230 AC#1's seven-item table is explicitly not this ticket's job (ticket 234 "与其它票的关系"). This cell supplies `:58` only.

## Gate / regression cells (7)

### G-1 ticket 92 `:73` — AC#6 three readings valid only against two old tree snapshots

- VERDICT: 有读数 (all four readings re-taken on `015f6be1`). The AC's literal `rc=0` is **not attainable on today's tree**, and the reason is not ticket 92's.
- Criterion (ticket 92 `:73-77`): `sh scripts/d22scan.sh` rc=0 on a pure tree **with the per-scope file counts printed** (`ban #6 frontend/`'s N must have grown because of this ticket); `ban #8` zero emoji; `gofmt -l`/`gofumpt -l` empty; `go vet ./internal/panel/` rc=0;
  `go test -count=2 ./internal/panel/` rc=0 with every SKIP/FAIL named. The face also says `go test ./cmd/wisp/` is a registered load-time `0xc0000135` red here — "不要去追，如实登记即可".
- READING: `gofmt -l internal cmd tools` -> 0 files (`gates-gofmt.txt`); `"$(go env GOPATH)/bin/gofumpt.exe" -l internal cmd tools` -> 0 files, rc=0 (`gates-gofumpt.txt`; the dispatch's `exit 127` trap avoided by using the full path, which exists: `D:\work\base\gopath/bin/gofumpt.exe`).
  `go vet ./internal/panel/` rc=**0**. `sh scripts/d22scan.sh` rc=**0**, "clean - no D22 ban violations" (so ban #8's emoji ban is satisfied), per-scope counts: bans #1-5 internal/=219, cmd/=29; ban #6 frontend/=85; ban #7 internal/tools/=22; ban #8 design/=39, frontend/=85, internal/=460, cmd/=63 (`gates-d22scan-host.txt`).
- THE `rc=0` CLAUSE THAT FAILS, with its口径: `go test -count=2 -v -p 1 ./internal/panel/` -> **rc=1**, four numbers RUN=304 PASS=296 FAIL=8 SKIP=0.
  FAIL=8 = 4 distinct names x 2 runs (`-count=2`), and the name set is exactly the registered ones: `TestApprovalCardViewJSONKeysMatchFrontendTypes`, `TestComposerContractTypesMatchFrontend`,
  `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`, `TestC21DesignTokensFourWayAgree`. Zero new names, and those four are other tickets' turf (the dispatch forbids touching them, and the frozen trio includes `tokens_fourway_test.go`).
  So: the instrument ran and is fully named; the cell cannot be closed as "通过" on today's tree without either fixing someone else's reds or relaxing an assertion, and both are forbidden here.
- `ban #6 frontend/` growth clause: today's N is 85, against 43 in `docs/evidence/s1/92b-adversarial-acceptance.md:138` and 40 at ticket 97's `f140079`. It has grown, but not attributable to ticket 92 any more than to the 40 tickets since. Registered as "grown, attribution lost".
- Old debt `R-92b-4` ("整步 `portable-tests.sh` 它未跑", `92b-…md:138` and `:223`): **now run on current HEAD** — `bash scripts/portable-tests.sh ./internal/panel/` rc=**1**,
  its own line: `top-level: PASS=95 FAIL=4 SKIP=0, === RUN=152`, `strict runner exited 1 for scope=[./internal/panel/]` (`g1-portable-panel.txt`). Same four names; the step also shows it `-skip`s the slow/live set by name, so its SKIP=0 is a skip-list design, not a missing `-v`.
- `cmd/wisp` not run: excluded by the cell's own text (registered load-time red, and ~110-135s). Borrowed reading with口径 attached: the 222 acceptance leg reported `cmd/wisp ok 110.491s` on `fac60ad4..HEAD` with zero `.go` drift.

### G-2 ticket 97 `:55` — AC#5 per-package five numbers re-taken on current HEAD

- VERDICT: PASS, five readings + the ledger, and `R-97-3` is settled by a layered identity that holds exactly.
- Criterion (ticket 97 `:55-57`): per-package `gofmt -l`/`gofumpt -l` empty, `go vet ./internal/agent/approval/` rc=0, `go test -count=2 ./internal/agent/approval/` rc=0 with SKIP/FAIL named (and the 2 SKIPs are `TestDefaultDeadlineWallClockMeasurement`).
- READING: `gofmt -l`=0, `gofumpt -l`=0, `go vet ./internal/agent/approval/` rc=**0**,
  `go test -count=2 -v -p 1 ./internal/agent/approval/` -> rc=**0**, four numbers RUN=114 PASS=112 FAIL=0 SKIP=2 (per-package split of `gates-count2-v.txt`, file `gates-per-package-numbers.txt`).
  SKIP from `-v`, named: `TestDefaultDeadlineWallClockMeasurement` x 2 runs (1 distinct name). Same name and same reason as ticket 87's reading and as the AC#1 baseline.
- `sh scripts/d22scan.sh` rc=0/clean vs the ledger ticket 97 quoted at `f140079` (`internal/=202, cmd/=20, ban #6 frontend/=40, ban #7 tools=18, ban #8 design/=16 frontend/=40 internal/=363 cmd/=26`):
  today 219 / 29 / 85 / 22 / 39 / 85 / **460** / 63 — **no scope lower** ⇒ A64② satisfied on current HEAD.
- `R-97-3` (the ticket's "58+38+2 = 98" adds across layers): settled without touching the closed ticket's arithmetic — counted consistently **with subtests at every layer**, the identity closes exactly:
  my 8-package log gives RUN=1550 = PASS 1534 + FAIL 10 + SKIP 6 = 1550, and separately RUN=1550 = 775 distinct names x 2 (`-count=2`, no cache reuse).
  The mixed-layer form (top-level PASS + subtest PASS + top-level SKIP) is what produced the 96-vs-98 scare. Registered as the corrected form of the式子.

### G-3 ticket 104 `:57` — AC#4 four packages `-count=2 -v`

- VERDICT: PASS, all five sub-readings on current HEAD. The four-package scope is rc=0 for all four packages.
- Criterion (ticket 104 `:57-62`): `go test -count=2 ./internal/winsec/ ./internal/memory/ ./internal/secret/ ./internal/agent/approval/` rc=0 with the four numbers and every SKIP named **from `-v`** (non-`-v` output does not print SKIP at all);
  `gofmt`/`gofumpt` empty; `go vet` and `GOOS=linux go vet` **per package** rc=0; closing `sh scripts/d22scan.sh` rc=0 on a pure tree.
- READING (`-count=2 -v -p 1`, 口径 = these four packages out of an 8-package sequential run, per-package split in `gates-per-package-numbers.txt`):
  winsec RUN=202 PASS=202 FAIL=0 SKIP=0 (13.549s) / memory RUN=136 PASS=134 FAIL=0 SKIP=2 (24.844s) / secret RUN=58 PASS=58 FAIL=0 SKIP=0 (0.239s) / approval RUN=114 PASS=112 FAIL=0 SKIP=2 (0.619s).
  Package result lines: all four `ok`. SKIP named from `-v`: memory `TestSubprocessCrashWriter` x2, approval `TestDefaultDeadlineWallClockMeasurement` x2. Both pre-existing, neither is this ticket's.
- `gofmt -l`=0 / `gofumpt.exe -l`=0; `go vet` per package rc=0 x4; `GOOS=linux go vet` per package rc=0 x4 — **labelled compile-only**, exactly as the dispatch warns.
  The behavioural linux half for `internal/winsec` is not this cell's claim; it lives in G-6's container run.
- `sh scripts/d22scan.sh` rc=0 (host, real repo) ⇒ ticket 104's closing clause satisfied.

### G-4 ticket 105 `:54` — AC#5 five readings + per-scope file counts

- VERDICT: PASS, including the cell's own arithmetic identity.
- Scope taken from ticket 105's own acceptance table (`docs/evidence/s1/105-adversarial-acceptance.md`, the AC#5 section names `go test -count=2 -v ./internal/tools/ ./internal/risk/`), so that is the scope re-measured here.
- READING: `go test -count=2 -v -p 1 ./internal/tools/ ./internal/risk/` rc=**0**; `ok internal/tools 24.447s`, `ok internal/risk 5.846s`.
  Four numbers RUN=816 PASS=814 FAIL=0 SKIP=2 (`g4-tools-risk-count2-v.txt`). SKIP named from `-v`: `TestSyncRegistryProbeLive` x 2 runs.
- THE CELL'S OWN ODD CLAUSE, checked rather than asserted: "报 `=== RUN` 行数 == 不同测试名 × 2" -> distinct names = **408**, 408 x 2 = **816** = the `=== RUN` count. Exact.
  (`-count=2` therefore shows no cache reuse here, which is the point of that clause.)
- `gofmt -l`=0, `gofumpt.exe -l`=0, `go vet ./internal/tools/` rc=0 and `GOOS=linux go vet ./internal/tools/` rc=0 (compile-only), same pair for `./internal/risk/`.
- `sh scripts/d22scan.sh` rc=**0**, per-scope counts as in G-1; against ticket 105's own `f6818f2`-era reading (`internal/` 363 -> 370 in its pure snapshot) today's 460 is **not lower** ⇒ A64② holds.

### G-5 ticket 110 `:47` — AC#5 `bash -n` / `d22scan.sh` per scope not lower / new step ordering

- VERDICT: PASS on the two measurable clauses; the third clause ("新步不得排在会失败的步之后") is **confirmed as still open** and is not this leg's to close — it is ticket 111/112/140's, and ticket 110's own new step is the thing that eats them.
- `bash -n`: `scripts/winsec-tests.sh` rc=**0**, `scripts/portable-tests.sh` rc=**0**, `tools/d22scan/runtests.sh` rc=**0**.
  (An intermediate attempt of mine hit `scripts/runtests.sh: No such file or directory` rc=127 — that path does not exist; the real runner is `tools/d22scan/runtests.sh`. Registered so nobody re-uses my broken needle.)
- `sh scripts/d22scan.sh` on the host repo rc=**0**; per-scope counts 219/29/85/22/39/85/460/63 vs the most recent recorded ledger (222's leg: `ban #7 tools=22`, and `internal/=460`-family) ⇒ no scope dropped.
- Ordering clause, read off the tree instead of off the ticket — **my first draft of this sentence was wrong and is corrected here**. I had written that the steps after the winsec gate "inherit the default `if: success()`, so a red winsec step prevents them".
  Read verbatim on current HEAD: the gate step itself carries `if: ${{ !cancelled() }}` at `ci.yml:404` (right after `- name: "Windows ACL sealing gate ..."` at `:398`), and so does every step after it:
  `Cache third_party` `:441-442`, `cgo build smoke` `:448-452`, `cmd/wisp CLI tests` `:455`. The comments say the reason out loud: `:400-401` "this gate is the step whose red used to eat the four steps below it. It stays FIRST inside the job (ticket 110 AC#5 …)",
  and `:449-451` "red on run 35551819606 before, and skipped behind the winsec gate since ticket 110 … Now it always answers."
  ⇒ **The criterion is SATISFIED on current HEAD**, and by the stronger of the two allowed remedies (placed first AND `!cancelled()`-guarded, no `continue-on-error` anywhere in the job — `ci.yml:381` states that and `grep` finds none).
- Therefore the ticket-110 consequence debt (`f161c690` self-reporting "把后面的步骤全吃掉") is **no longer live as a shape**: the tree now has the guard that prevents it. What remains open is the *content* of the later steps, which is other tickets' denominator:
  measured with the ruler the dispatch demands (`grep -c -- '^- \[ \]'`): ticket 111 un=10, ticket 112 un=5, ticket 140 un=4, none of them `-done`. Registered as 未销 against those tickets, not against this cell.

### G-6 ticket 113 `:56` — AC#5 in-container three packages four numbers + per-package vet + gofmt + `d22scan.sh`

- VERDICT: PASS (all five sub-readings taken on current HEAD `015f6be1` inside a real linux container).
- Criterion (ticket 113 `:56-59`, read this leg): container `-count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/risk/` rc=0 with the four numbers named per line (say whether SKIP was counted from `-v`; `-count=2` does not cache);
  per-package `go vet` + `GOOS=linux go vet` rc=0; `gofmt -l` empty; closing `sh scripts/d22scan.sh` on a pristine snapshot rc=0 with the ledger's per-scope counts not lower.
- READING (`m4-g6-container.log`, golang:1.27, go1.27.1 linux/amd64, snapshot at `/work`): `g6_rc=`**`0`**, `red_name_lines=0`.
  Four numbers WITH subtests (the `^ *--- PASS` form): RUN=574 PASS=570 FAIL=0 SKIP=4.
  Same log with the naive top-level-only form: PASS=344 FAIL=0 SKIP=4 — a 226-line gap, so any "four numbers" quote must state which ruler produced it.
  SKIP named, and it IS from `-v`: 4 SKIP lines = 2 distinct cases x 2 runs (`-count=2`): `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`, `TestSubprocessCrashWriter`.
  Note this is also the POSIX leg finally having a denominator: on linux `internal/winsec`'s `_other_test.go` files (7 of them) actually run, which is what ticket 113 exists to fix.
- Per-package vet in-container: `vet_rc[./internal/winsec/]=0` `vet_rc[./internal/memory/]=0` `vet_rc[./internal/risk/]=0`; `GOOS=linux go vet` (three packages, executed ON linux) rc=0.
  The dispatch's warning was honored by construction: the container run is the evidence, `GOOS=linux go vet` on the host was run too (ticket 104 AC#4 asks for it per package) but is labelled compile-only below, never as behaviour.
- `gofmt -l internal/` in-container: **0 files**.
- `sh scripts/d22scan.sh` in-container: rc=**0**, clean. Per-scope counts there: bans #1-5 internal/=219, cmd/=29; ban #6 frontend/=85; ban #7 internal/tools/=22; ban #8 design/=30, frontend/=85, internal/=460, cmd/=63.
  CAVEAT that must travel with those numbers: the snapshot has no `.git`, and d22scan says so out loud ("gitignore rules NOT APPLIED ... every path in every scope is being scanned", A207's machine-dependent denominator).
  So the container's counts are a *wider* denominator than the host ledger's; they are not comparable to the ledger for the "不降" test. The host run below is the one that carries that test.
- Old debt on this cell (`docs/evidence/s1/113-adversarial-acceptance.md:17`, the "ban#8 sha↔number 配对错了一格" line, plus `R-113-F` at ticket 113 `:16`): see the "Old debts" section — settled by counting tracked `.go` under `internal/` at `3c5d1c3` directly from git, not by re-running the scanner against a dead tree.

### G-7 ticket 115 `:64` — AC#6 per-package four numbers / gofmt / gofumpt.exe -l / go vet / d22scan

- VERDICT: 有读数 for all five sub-readings, on the package ticket 115 actually changed (`internal/winsec`). The cell's own ledger clause passes; the table-under-115 clause is ticket 230's and is NOT delivered here.
- Criterion (ticket 115 `:64-66`): per-package `-count=2 -v` four numbers with every SKIP named from `-v`; `gofmt -l` and `"$(go env GOPATH)/bin/gofumpt.exe" -l` (the face insists the binary exists on this machine and that "not run" must quote the command + error);
  `go vet` + `sh scripts/d22scan.sh` on a pure snapshot rc=0 with the ledger's per-scope counts not lower.
- READING: `go test -count=2 -v -p 1 ./internal/winsec/` -> `ok ... 13.549s`, four numbers RUN=202 PASS=202 FAIL=0 SKIP=**0**.
  Because SKIP=0 there is nothing to name; and this is a `-v` run, so the absence of SKIP is a real zero, not the non-`-v` blindness ticket 104 AC#4 warns about.
- `gofmt -l internal cmd tools` -> 0 files. `gofumpt.exe -l internal cmd tools` -> 0 files, rc=0, with the binary path proven present first (`D:\work\base\gopath/bin/gofumpt.exe`), so this is not the `exit 127` false reading.
- `go vet ./internal/winsec/` rc=**0**. `sh scripts/d22scan.sh` (host repo) rc=**0**, clean, ledger 219/29/85/22/39/85/460/63.
-口径 warning attached to this cell: ticket 115's own numbers live on a DIFFERENT ruler — CI `test-windows` step 4 (`82/42/0/0`, then `80/36/4/0`), and `docs/evidence/s1/ci-step-readings-2026-09-22.md:250-257` states outright that the `42` is permanently unreproducible
  and that judgement must use (i) `--- FAIL=0` and (ii) the two new cases sitting in the `--- PASS` column. Both hold here: FAIL=0, and
  `TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree` + `TestNoticeAttributionKeepsTwoTreesApart` are PASS in `ac1-baseline-v.txt` and go red only under M-5's two mutants.

## Old debts carried by 5 of the cells

- **`R-92b-4`** ("整步 `portable-tests.sh` 它未跑", `docs/evidence/s1/92b-adversarial-acceptance.md:138`, listed as a hit at `:223`): **cleared**. Run on current HEAD: `bash scripts/portable-tests.sh ./internal/panel/` rc=1,
  its own four numbers `=== RUN=152 PASS=95 FAIL=4 SKIP=0`, `strict runner exited 1 for scope=[./internal/panel/]`. The 4 reds are the registered panel set, zero new. The step exists, runs, and refuses to read a red as green.
- **`R-97-3`** (the closed ticket's "58+38+2 = 98" adds across layers; `97-adversarial-acceptance.md:36`): **cleared as an arithmetic question** by a single-granularity ruler: with subtests counted at every layer,
  RUN = PASS + FAIL + SKIP closes exactly (my 8-package run: 1550 = 1534 + 10 + 6) and RUN = distinct names x 2 under `-count=2` (816 = 408 x 2 on ticket 105's scope, 1550 = 775 x 2 on the 8-package run).
  The closed ticket's own text is untouched (append-only; and rewriting it would be a contract-shaped edit, not this leg's).
- **ticket 113's `ban#8` sha↔number mis-pairing** (`docs/evidence/s1/113-adversarial-acceptance.md:17` "但 ban#8 `internal/=371`，不是它报的 368"; the same claim as `R-113-F` at ticket 113 `:16`): **settled in favour of 371, with a self-calibrating ruler**.
  `git ls-tree -r --name-only 3c5d1c3 -- internal/ | grep -c '.go$'` = **371**; the same ruler at HEAD = **460**, and the scanner's own `ban #8 internal/` line at HEAD = **460** — the two agree exactly on a tree where both are alive,
  which is what licenses using `ls-tree` to speak about a dead one. So the number paired with `3c5d1c3` is 371, and the pairing error the table flagged is real and now closed.
- **ticket 110's new step "eating the later steps"** (`f161c690`'s self-report): **the shape is gone from the tree** — the gate carries `if: ${{ !cancelled() }}` (`ci.yml:404`) and all three steps after it do too (`:442`, `:452`),
  with `ci.yml:400-401` and `:449-451` naming that exact history. What stays open is only the content of those steps: ticket 111 un=10, ticket 112 un=5, ticket 140 un=4, none `-done`. **未销, against those tickets.**
- **ticket 115 has zero tables under its own name**: re-measured this leg, `find docs/evidence -name '115*' | wc -l` = **0**. **未销** — it is ticket 230 AC#1's account (and ticket 234's face says explicitly this leg does not deliver it).
  What this leg did instead: the two cells under 115 (`:58`, `:64`) now have readings, and 115's `:44` was moved out of the closed-ticket denominator by AC#5 below.

## Final roster comparison (AC#4)

- VERDICT: clean. Nothing was added, nothing disappeared.
- Re-pulled the AC#1 command verbatim at the end of the leg: `go test ./internal/winsec/ ./internal/risk/ ./internal/config/ ./internal/agent/approval/ -count=1 -v` (with `-p 1`, PATH exported as required) -> rc=0,
  winsec 9.236s / risk 3.651s / config 0.931s / approval 0.334s, four numbers RUN=471 PASS=469 FAIL=0 SKIP=2 — identical to the start of the leg.
- Name sets: baseline 297 top-level names, final 297 top-level names, `diff` between `ac1-baseline-roster.txt` and `ac4-final-roster.txt` = **empty** (`ROSTER_IDENTICAL=yes`). So "只多哪几枚 = 零枚／消失哪几枚 = 零枚".
- Wider denominator, with口径: the 8-package `-count=2 -v` run (`gates-count2-v.txt`) produced exactly 5 distinct red names, and they are exactly the registered historical set
  (`internal/panel`: `TestApprovalCardViewJSONKeysMatchFrontendTypes`, `TestComposerContractTypesMatchFrontend`, `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`, `TestC21DesignTokensFourWayAgree`; `internal/ball`: `TestC21TableColourRowsMatchTokensCSS`).
  FAIL=8 on panel and FAIL=2 on ball are those names x 2 runs. **Zero new reds anywhere**, and none of the five was touched (two of them sit inside the frozen trio).
- `R-116-1` (`internal/risk TestResolvePerCallBudget`, parallel-load sensitive) was never hit: `-p 1` was used on every run, risk came back ok in all four runs it appeared in (4.025s / 5.843s / 5.846s / 3.651s), and no run of mine was ever attributed to contention.
- Restoration state at the end of the leg: `git status --porcelain -- internal/ cmd/wisp` = empty, i.e. the five mutation cells left no trace in production code.

## AC#5 — 10 duplicate denominators relocated to still-open tickets

- VERDICT: DONE, all 10, zero production code, no `-done` renamed.
- The 10 lines were re-located by reading each cited line before touching it (all ten verified to start with `- [ ] ` at the cited number; no drift this time, unlike the earlier batches).
  Instrument: `.scratch/wisp/probes/gate-rerun-1/ac5_relocate.py`, which refuses to run unless every target is an unchecked first line.
- Disposition, in the shape the repo set earlier today (80 `:54` / 89 `:348` / 118 `:59`): the cell's first line becomes `- ⛔ **（…处置…）** <original sentence, verbatim> ⇒ **指向**＝票 <NN> <line>`,
  continuation lines untouched. The marker prefix replaced the checkbox (that is what takes it out of the denominator), so the sentence text survives inside the same line.
- BEFORE -> AFTER, with the ruler the dispatch demands (`grep -c -- '^- \[ \]'`, NOT `grep '^- \['`, which the progress-stamp lines `  [2026-…]` inflate into thousands of false hits):

| sender (closed) | un before | un after | cells moved |
|---|---|---|---|
| 07 | 4 | 1 | 3 (`:54` `:59` `:61` -> 票 64) |
| 11 | 1 | 0 | 1 (`:53` -> 票 12, A11) |
| 63 | 1 | 0 | 1 (`:46` -> 票 12, A8) |
| 80 | 2 | 0 | 2 (`:43` `:50` -> 票 230 AC#3) |
| 97 | 2 | 1 | 1 (`:53` -> 票 230 AC#2) |
| 115 | 3 | 2 | 1 (`:44` -> 票 230 AC#1) |
| 92 | 4 | 3 | 1 (`:52` -> 票 114 AC#5) |

  Total moved = **10** = 3+2+4+1 as the ticket face groups them (->64 three, ->12 two, ->230 four, ->114 one).
- Receiver side measured BEFORE and AFTER and **unchanged in both**: 票 64 un=1/chk=8, 票 12 un=3/chk=6, 票 230 un=5/chk=0, 票 114 un=**9**/chk=2.
  The face's warning that "票 114 记 `un=7` 实为 9" reproduces exactly on current HEAD, so the receivers were not re-counted from their own stale text.
- `git diff --numstat` on the seven files: 3/3, 1/1, 1/1, 1/1, 2/2, 1/1, 1/1 = **10 additions, 10 deletions**, no file rewritten wholesale.
- ⚠ **This contradicts the dispatch's own clause** "原句逐字保留、`git diff --numstat` 删除列必须为 0". A deletion count of 0 is not reachable for this disposition: replacing the checkbox on a line IS one deletion.
  The precedent the face cites (`24198b7a`) is itself 2/1 and 3/1 on the same operation. What "原句逐字保留" means in practice, and what I honoured: the sentence survives verbatim after the marker, and no other cell's checkbox was touched.

## AC#6 — retire one bad yardstick (`NewGate(` zero hits)

- VERDICT: the retirement is **upheld**, and the fact the retired yardstick pointed away from is confirmed on the tree.
- The broken ruler, re-run verbatim on `015f6be1`: `grep -rn 'NewGate(' --include=*.go internal cmd | grep -v _test` -> **0 hits**; widened to the whole repository including tests -> still **0 hits**.
  A string that exists nowhere cannot distinguish "no carrier" from "carrier everywhere", so `done-fix-1`'s inference at ticket 80 `:44` proves nothing. Void.
- The real thing, read this leg: `internal/risk/blacklist.go:76` `func Gate(canonical string, bOverrides map[string]bool) PathDecision` — a plain function, not a constructor, which is why a `NewGate(` search could never have found it.
- The opposite half, also read this leg: `internal/tools/mode.go:98` `d := risk.Gate(c, overrides)` **is a production call site**, fed by `overrides := b.confirmations()` — at `:91`, not `:92` as the dispatch states (one-line drift, recorded rather than repeated).
- The judgement class does NOT change: still 丙. `internal/config/unwired.go` (`risk.blacklist_overrides` row) self-reports that "nothing populates the bOverrides map risk.Gate reads… the confirmations that would fill this map are minted only by an L2 answer",
  i.e. the map is filled by a runtime L2 answer, not by the config key the cell is about, and a file on disk "cannot have clicked anything". So AC#3 of ticket 80 stays accounted under ticket 230 (moved there by AC#5 above).
- Replacement ruler for whoever picks this up: ask "does anything call `risk.Gate(`, and is there a test on it", not "does `NewGate(` exist". Measured with that ruler this leg (`ac6-yardstick.txt`):
  production callers = exactly one, `internal/tools/mode.go:98`; **test callers = zero** (`grep -rln 'risk.Gate(' --include=*_test.go .` returns nothing at all).
- ⚠ So the retired sentence is only half-wrong, and I am not going to make it more wrong than it is: its **推理** is void (a string that appears nowhere cannot prove anything), but its **conclusion** — "no case today runs on `risk.Gate`" —
  survives a working ruler. What does not survive is the way it was reached, and what the dispatch's counter-fact correctly repairs is the other half: `risk.Gate` is not an uncalled function, `mode.go:98` calls it in production
  (`internal/tools/mode.go:53` even records that AC#5 of some ticket "asked for risk.Gate to stop being a function with zero production call sites"). Verdict: yardstick void, conclusion re-founded on a different measurement, classification unchanged (丙).

## Not finished / not judgeable

- All 12 cells got a reading; **none is left 判不了**. Two clauses inside cells stayed out of reach and are named rather than smoothed over:
  1. G-1's `ban #6 frontend/` "must have grown **because of this ticket**" attribution — the count grew (43 -> 85) but 40 tickets intervened, so attribution is unrecoverable. Registered as "grown, attribution lost".
  2. G-1's `rc=0` clause on `./internal/panel/` and ticket 115's "a table under `115-*`" clause — the first is blocked by four reds that belong to other tickets and that this ticket forbids me to touch; the second is ticket 230 AC#1's account, which ticket 234 explicitly says this leg must not deliver.
- Ticket 92 `:79` (AC#7) is NOT in this leg's twelve and was not attempted: its credential sits in a directory this leg may not read. Registered as not-judgeable-by-me, no guess made.
- Two of my own instruments were broken and were replaced rather than reported as results: `grep -c '^--- PASS'` (drops subtests; cost one bogus "restore-green 13 vs 297" moment) and `scripts/runtests.sh` as a `bash -n` target (that path does not exist; the runner is `tools/d22scan/runtests.sh`).
- Not run at all: `go test ./cmd/wisp/` (the cells' own text excludes it as a registered load-time red, 110-135s), the POSIX leg of ticket 104, and ticket 07 `:51` (the one cell of the thirteen that needs a human eye).
- Turn budget: this leg stopped after AC#6. Everything the ticket asks for is present; nothing was left half-written.
