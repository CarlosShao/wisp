# 票 285 · 腿 `285-r1` · 起手锚

生成时刻 `2026-10-09 10:3x +08`（第 1 笔 commit `97392ab` 落盘 10:32；收工复量 10:52，⛔ 本腿时钟一律取 `date` 与
`git log --date=format:%H:%M` 现量，不抄票面上的名义时刻）。本腿＝落地腿（写测试件），非裁决腿。

## 锚

- HEAD（起手）＝`37f2a0bbfd978098875e53352ba1e5722ae200f2`，分支 `dev`。
- 该 commit 的标题即票 285 的立票 commit（`A757 落账＋票 242 结案…＋新立票 285`）⇒ 本腿起跑时票面三格全未勾。
- 起手时 `git status --porcelain -- internal cmd` ＝**空**。

## 基线四数（本腿自跑，`-count=1`，⛔ 非 `-overlay`、非 `-v` 复用）

命令＝`go test ./internal/agent/approval/ -count=1 -v`
⇒ rc=0／`--- FAIL` **0** 枚／`--- PASS` **90** 枚／`FAIL` 行 **0**／`--- SKIP` **1** 枚／`ok github.com/CarlosShao/wisp/internal/agent/approval 0.404s`

- 红名册＝**空**（0 枚 FAIL）。
- SKIP 名册（逐字）＝`TestDefaultDeadlineWallClockMeasurement`（`--- SKIP: TestDefaultDeadlineWallClockMeasurement (0.00s)`，基线输出第 153 行）。
- 与 `242-v2` 现量基线（`docs/evidence/s1/242-grant-binding-v2.md:8`，HEAD `2bdf385d`，90 PASS／0 FAIL）**枚数等值**；本腿不需顶回差。
- ⚠ `--- PASS` 里有同名重复行（表驱动子用例的顶层名各计一次，例：`TestEveryAnswerRouteForAnUnknownCorrelationFailsFast` 9 次、`TestEveryRefusalRouteOnAnUnknownEntryStillRefuses` 6 次）。**90 是行枚数口径**，不是唯一用例名口径
  （现量唯一名＝**72 枚**，`grep -oE '--- (PASS|SKIP): [A-Za-z0-9_]+' | sort -u | wc -l` 于收工前复量；
  ⛔ 初稿在这里写"61 枚"是本腿未量就下笔的一处错，就地订正，见腿的收工 commit）。分母一律按行枚数记。

## 名册（90 PASS ＋ 1 SKIP，按 `-v` 输出 `--- ` 行原样抽取，含重复）

```
PASS: TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses
PASS: TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard
PASS: TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant
PASS: TestAnAliasCanNeverBuyAnAllow
PASS: TestAnAmbiguousVetoNameIsNotGuessedOnTheRefusalSide
PASS: TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline
PASS: TestAnUnreachablePromptSurfaceIsRefusedImmediatelyNotWaitedOut
PASS: TestBatchAggregatesHomogeneousL1Ops
PASS: TestBatchRefusalsBelowThresholdAndOnSensitiveVerdicts
PASS: TestDefaultQueueDeadlineIsFiniteAtThreeHundredSeconds
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
PASS: TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
PASS: TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
PASS: TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
PASS: TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
PASS: TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
PASS: TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
PASS: TestGrantIsSingleUseAndBoundToItsItem
PASS: TestHostUnreachableAndFullQueueFailClosed
PASS: TestHostUnreachableAndFullQueueFailClosed
PASS: TestHostUnreachableAndFullQueueFailClosed
PASS: TestL1WindowIsPlainDataWithNoRouteToAnAnswer
PASS: TestL1WindowTimeoutMeansExecute
PASS: TestL1WindowUnavailableVoiceChannelIsAnnounced
PASS: TestL1WindowUnreachableUIFailsClosed
PASS: TestL1WindowVetoOnForeignCorrelationCannotCancel
PASS: TestL1WindowVetoedByEachLoadedChannel
PASS: TestL1WindowVetoedByEachLoadedChannel
PASS: TestL1WindowVetoedByEachLoadedChannel
PASS: TestL2NativeAllowExecutesThroughTheBridge
PASS: TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s
PASS: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord
PASS: TestLiveApprovalsOnNilQueueIsNil
PASS: TestLiveApprovalsReportsPendingRowsAndDoesNotConsumeThem
PASS: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue
PASS: TestLiveApprovalsSkipsRowsThatAreNotPending
PASS: TestLiveL1WindowsSeesAWindowInTheLoop
PASS: TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis
PASS: TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis
PASS: TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis
PASS: TestR7SizedBatchEscalatesToL2Unaggregated
PASS: TestR7SizedBatchThatNobodyAnsweredIsRejectedNotRun
PASS: TestReadingL1WindowsMovesNoGateState
PASS: TestReplayRedisplaysUnderAFreshGrant
PASS: TestReplySeamCarriesNoAllowWithoutAGrant
PASS: TestReplySeamWaitingStateNamesTheTwoD43Rows
PASS: TestStartedCallsAreNotReportedAsWaiting
PASS: TestTenOpsInOneToolCallGetOneConfirm
PASS: TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNamed
PASS: TestTicket224ForgedSessionAnswerRecordsNothing
PASS: TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow
PASS: TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped
PASS: TestTicket242BindDigestSeparatesItemsBySequenceNumber
PASS: TestTicket242ForgedBindingCannotSpendAnotherItemsGrant
PASS: TestTicket242PanelItemReadFaceIsFullyDeclared
PASS: TestTicket242PanelItemStaysGrantFree
PASS: TestTicket242QueuedItemsBindGrantsToTheirOwnDigest
PASS: TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce
PASS: TestTicket242SpendRequiresTheExactBinding
PASS: TestTicket259R1DenialNamesCardNotPending
PASS: TestTicket259R1DenialNamesMisbound
PASS: TestTicket259R1DenialNamesMissingNonce
PASS: TestTicket259R1DenialNamesSpentOrNeverLiveNonce
PASS: TestTicket259R1FourDenialLabelsArePairwiseDistinct
PASS: TestTicket259R1OutwardAnswerStaysMergedAcrossDenials
PASS: TestTicket259R1PanelAPIMethodSetIsClosed
PASS: TestTicket259R1PanelCarrierCarriesNoAnswerVerb
PASS: TestTicket259R1PanelItemCarriesNoCredentialValue
PASS: TestTicket259R1PanelItemFieldCapabilityFence
PASS: TestTicket259R1PanelItemProjectionIgnoresGrantState
PASS: TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer
PASS: TestTicket260R3EscUnavailableLineNamesNoKey
PASS: TestTicket260R3LoadedCancelLabelIsTheNamedResidual
PASS: TestTicket260R4DefaultLabelIsTheOldSentence
PASS: TestTicket260R4LabelFollowsTheReaderAcrossAReload
PASS: TestTicket260R4SeededKeyReplacesEscOnEveryFace
PASS: TestTicket260R4SpellingSeamCarriesNoAuthority
PASS: TestTicket260R4UnloadLineNamesNoKey
PASS: TestUnadmittedTaskIsRefusedBeforeQueueing
PASS: TestVetoAfterStartYieldsAppliedStepsReport
PASS: TestVetoTravelsBackToTheLoopAsAFailure
PASS: TestWideningRuleReadsOneLineOffTheCard
SKIP: TestDefaultDeadlineWallClockMeasurement
```

## 种前四枚产码件 hash（`git show HEAD:<path> | git hash-object --stdin`）

| 文件 | hash |
|---|---|
| `internal/agent/approval/approval.go` | `67fb146898dc7b235623bb8ad2fac0e0b2465fe1` |
| `internal/agent/approval/gate.go` | `78e2d28a8b676754514bd80898ace8c5a974d3ce` |
| `internal/agent/approval/queue.go` | `66fec7ae375344d7bc741270b42c4ec5050cf41f` |
| `internal/agent/approval/ui.go` | `55bcd1daf251a9ef513ef43b3c63605787252b13` |

⇒ 四枚与 `242-v2` 表（`242-grant-binding-v2.md:25-28`）**逐字等值**＝HEAD 从 `2bdf385d` 走到 `37f2a0b` 没动这四枚。

## 现成夹具（本腿复用，不新造）

`internal/agent/approval/ticket259_denial_rulers_test.go`：
`:36` `type r1Audit struct`（把 `Options.Logf` 接进切片）· `:48 has(want string)` · `:59 dump()` ·
`:71 r1Card(t, taskID, corr, tool)`（返回 `*Gate, *Queue, *qitem, nonce, *r1Audit`）·
`:83 r1CardIn(t, q, taskID, corr, tool)`（**同一本 queue 上再立一枚活卡**，正是两枚活卡需要的形状）。

## 顶回派单的一处

派单写 `Native().Allow(B的corr, A的grant)`（两参）。盘上现量是**三参带 ctx**：
`ticket259_denial_rulers_test.go:118` 逐字 `err := g.Native().Allow(context.Background(), it.Corr, "")`。
以盘上为准。票面 `:16` 自己写的是 `Native().Allow(B的corr, A的活令牌)`（形状描述，非调用字面），不动。
