# 21-verdict — approval gates minimal, segment 2 (card render + depth badge) — NON-implementer verdict

- **Leg:** non-implementer verdict (re-dispatch after the first verdict run died without a verdict; narrow scope), 2026-11-06
- **Verdict agent:** NOT the implementer (implementer = agent=C, wave 34). Rule 6 adversarial pass.
- **Commits examined:** impl `9468fccad` + claim `fb10eae49`, both verified ancestors of dev HEAD (`git merge-base --is-ancestor` rc=0; clone tip during this leg: `48610b495`). Prior context HEAD `218ed2ee9`.
- **Method:** all gates re-run synchronously by THIS agent, in a fresh `git clone --shared` at `/d/tmp/wisp/20261106/21v/tree` (rule 21: main tree's `cmd/wisp` carries ticket 323's untracked `panel_resident_323_windows_test.go` + `panel_resident_assembly_323_test.go`, confirmed on disk). CGO left at default (1); `GOTMPDIR/TMPDIR/TEMP/TMP` unset; PATH = main tree `build` + `third_party/sherpa-onnx`. Zero push; no ticket/production file touched by this leg.

## Gate readings (this agent's own runs)

| Gate | Reading | Note |
|---|---|---|
| `go build ./cmd/wisp ./internal/agent/approval` | **rc=0**, 13.3s | main tree |
| `go test ./internal/agent -run 'Approval\|Card\|Badge\|Depth'` (literal dispatch form) | rc=0, **"no tests to run"**, 3.3s | the `internal/agent` root package has zero tests matching that regex; readings honest |
| same regex over `./internal/agent/...` | **ok** (`approval` 0.373s), rc=0, 2.9s; 23 RUN lines, 0 FAIL | the regex's real target is `internal/agent/approval` |
| `go test ./internal/agent/approval -count=1` (full package, clone) | **ok 0.712s**, rc=0 | covers L1/batch/queue/veto AC machinery |
| `go test ./cmd/wisp -run <43 discovered approval/card/badge test names>` (clone) | first run rc=1, 101s: 2 FAIL (`TestAC246ShippedResidentProcessOwnsItsCancelStep`, `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`) | see triage below |
| — after junctioning gitignored `third_party\sherpa-onnx` (DLLs) into the clone | **`TestAC246Shipped...` PASSES**; only `TestTicket256...` remains red, 94s | DLL glob is clone-relative ⇒ environmental, matches impl receipt's "clone-missing-gitignored-dir" class |
| `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims` at clean pre-impl baseline `090b3f5d9` (second shared clone `21v/base`) | **FAIL identical** (7 of 11 fields vs want 6) | pre-existing drift (ticket 375 added `Options.Answer`, 256 roster stale), NOT this leg — impl receipt claimed exactly this; baseline-proven here |
| gofumpt -l on all 6 touched files | **empty (0 listed)**, rc=0 | |
| Effective focused set | **42/43 cmd/wisp names PASS + full approval package PASS**; the 1 red is baseline-identical | all 5 `TestTicket21S2*` names PASS |

Trust note (per dispatch): full `./...` build / `go vet` / `d22scan` are the orchestrator-accepted 门禁 boxes taken from the implementer's recorded readings (`21-segment2-impl.md`: build rc=0, vet rc=0, d22scan rc=0 clean with positive control PASS=35/FAIL=0). Not re-run here — out of the narrow scope.

**Kill rule:** no single command exceeded 6 minutes (max 101s). Nothing is PARTIAL; the first-run FAILs are triaged above, not hidden.

## Per-AC adjudication (ticket `21-approval-gates-minimal.md`)

| AC | Verdict | Readings behind it (all run by this agent unless noted) |
|---|---|---|
| **AC#1** L1: countdown executes on timeout; each working veto channel cancels; non-loaded-KWS message asserted; veto returns failure to LLM | **PASS** | `TestL1WindowTimeoutMeansExecute`, `TestL1WindowVetoedByEachLoadedChannel`, `TestL1WindowUnavailableVoiceChannelIsAnnounced`, `TestVetoTravelsBackToTheLoopAsAFailure`, `TestUnadmittedTaskIsRefusedBeforeQueueing` (full approval pkg run, ok); host-wired veto `TestL1VetoNeedsAChannelTheHostReallyWired` (cmd/wisp focused set PASS); 「语音取消不可用」carried verbatim in the card's channel lines (`TestTicket21S2ChannelsSayUnavailableNotAvailable` PASS) |
| **AC#2** L2: card shows FULL params + rules_hit; native allow executes; simulated panel-source allow REJECTED with forged source; 300s timeout → reject + warning at 270s + replay | **PASS** | card face = this segment's core: `TestTicket21S2CardPrintsFullParamsNotSummary` / `...FallsBackToParamsMap` / `...SaysOutLoudWhenNoParamsArrives` / `...KeepsUnparseableArgsVerbatim` + round-trip `TestTicket21S2ConsoleL2CardShowsFullParamsAndDepth` (real gate → real console card → native allow settles) — all PASS; native allow `TestReplyListenerAllowsAnL2CardFromTheNativeSide`, `TestL2NativeAllowExecutesThroughTheBridge` PASS; forgery `TestNativeHostSeamRefusesAPanelSourcedAllow`, `TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis`, `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` PASS; deadline `TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s`, `TestReplayRedisplaysUnderAFreshGrant`, `TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline` PASS |
| **AC#3** Batch: 10 L1 ops → one confirm; L2 mixed in → no aggregation; 50+ → L2 (R7) | **PASS** | `TestTenOpsInOneToolCallGetOneConfirm`, `TestBatchRefusalsBelowThresholdAndOnSensitiveVerdicts`, `TestR7SizedBatchEscalatesToL2Unaggregated`, `TestR7SizedBatchThatNobodyAnsweredIsRejectedNotRun` (full approval pkg ok); card side now renders it: `TestTicket21S2ConsoleL1BatchCardIsAggregatedView`, `TestTicket21S2CardBatchViewIsD45Shape`, `TestTicket21S2NoBatchMeansNoBatchLines` PASS |
| **AC#4** Cancel-mid-execution → applied-steps report rendered | **PASS** | `TestVetoAfterStartYieldsAppliedStepsReport` (full approval pkg ok, includes the "tool reported no steps" worst case per segment-1 receipt) |
| **AC#5** State machine: Confirming/AwaitingApproval transitions match D43 #17, 21–24 | **PASS** (on this leg's readings) | row-17 mapping `L1→Confirming` / `L2→AwaitingApproval` asserted in the `TestAC246*` block (PASS in focused set); waiting-state name + the two D43 rows `TestReplySeamWaitingStateNamesTheTwoD43Rows`, `TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` (AwaitingApproval while pending) PASS; Confirming veto route `TestLive246ConfirmingCardBorrowsEscVetoesAndReturns` PASS; deadline→reject `TestUnansweredL2CardTimesOutIntoRejectNeverExecution` PASS. Caveat kept honest: `lifecycle_43_windows_test.go` (D43 row-24 void walk) is outside this leg's 43-name focused set and was not re-run by this verdict leg; no contradicting reading anywhere in what ran |

## Rule closeable?

**YES — closeable.** Both implementer commits verified in-tree; every ticket AC has a first-hand PASS reading from a non-implementer run; the only red in the focused set is baseline-identical pre-existing drift (256 roster, ticket 375's follow-up — NOT ticket 21's to fix); gofumpt clean. The AC boxes may be ticked and the ticket marked `-done` with the two residuals below carried forward by name.

## The 2 named residuals (both confirmed real, neither blocking)

1. **`decideRisk` dead-branch cleanup** — `internal/agent/loop.go:916-930`: the declared-L1/L2 branch now passes through whenever `Options.AdmitTask` is bound (every real host does), so the old refusal arm is dead code. Confirmed present at HEAD. Deletion belongs to `internal/agent` (ticket 12/162's write face), outside this leg's scope; `approval/doc.go:42-48` records it as cleanup, not a live fix.
2. **Depth-badge REAL-WINDOW visual signoff = winlive, not fleet** — the badge's headless half is delivered and re-verified here (pure mapping `TestTicket21S2ApprovalDepthBadgeMapping`, AST census `TestTicket21S2BadgeWiringCensus` with the mutation-proven tooth, ordering `TestTicket21S2BadgePostsAreOrderedStateThenDepth` — all PASS). The on-orb visual on a real window belongs to the owner-present (winlive) pass per the 260-r3 precedent; the fleet's headless rulers cannot see pixels and do not pretend to.

## Reproduction

```
cd "/d/work/workspace/projects plans/Wisp"          # main tree: gate 1 build + gofumpt
git clone --shared . /d/tmp/wisp/20261106/21v/tree  # rule-21 clone (323 untracked files stay home)
# junction third_party\sherpa-onnx from the main tree into the clone for DLL-globbing tests
REGEX=$(grep -h "^func Test" cmd/wisp/*approval*_test.go cmd/wisp/*card*_test.go cmd/wisp/*badge*_test.go | sed 's/^func \(Test[A-Za-z0-9_]*\)(.*/\1/' | sort -u | paste -sd'|')
go test ./cmd/wisp -run "$REGEX" -count=1 -timeout 30m          # 42/43 PASS; 256 = baseline drift
go test ./internal/agent/approval -count=1                       # ok
git clone --shared . /d/tmp/wisp/20261106/21v/base && git -C /d/tmp/wisp/20261106/21v/base checkout 090b3f5d9
cd /d/tmp/wisp/20261106/21v/base && go test ./cmd/wisp -run TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims -count=1  # FAILs here too ⇒ pre-existing
```
