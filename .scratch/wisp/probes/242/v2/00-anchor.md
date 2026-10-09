# 242-v2 起手锚（非实现者对抗验收腿，重裁票 242 三格）

⚠ 本目录**不是空**：`.scratch/wisp/probes/242/v2/` 此刻已有 `anchor.md`／`verdict.md`／`logs/`，
时刻戳＝`2026-10-08 11:52 +0800` 一枚同名前腿（它在仓外用 `go test -overlay` 造错法）。
本派单**明令 ⛔ 不许用 `-overlay`** ⇒ 那枚前腿的读数**不作本腿凭据**，本腿全部错法种在盘上并逐发还原。
⛔ 前腿三件**不覆盖、不删**；本腿只写新文件名（`00-anchor.md` 起）。

- 起手 `date` ⇒ `2026-10-09 10:08 +0800`，rc=0。
- `git rev-parse HEAD` ⇒ `2bdf385d3037b7063d797bd07428dc1d37cb2561`。
- `git log --oneline -3` ⇒
  - `2bdf385d 死腿遗产代提（两枚腿 09:5x 被服务端连接中断掐掉，⛔ 不按额度定式重派）`
  - `6029b743 277-r1 A 道真机读数…`
  - `8ccc9568 280-r1 起手锚…`
- `git status --porcelain -- internal cmd` ⇒ **空**（与票 242「起跑名册」那句一致；工作树别处的在飞改动
  ＝`design/**` 删除批／`M .gitignore`／大量 `probes/**`，本腿零碰零转述）。
- 写面确认：`docs/evidence/s1/242-grant-binding-v2.md` 起手**不存在**（`ls -la docs/evidence/s1/242-*.md`
  只有 `242-grant-binding-v1.md` 29,758 字节，Oct 3 09:26）。

## 本腿基线名册（派单前必做撞钉预检，票 242「排程」那节原话）

尺＝`go test ./internal/agent/approval/ -count=1 -v`（本腿 10:0x 现跑，⛔ 不引编排者 09:4x 那句 `ok 0.369s`）
⇒ rc=0，`ok github.com/CarlosShao/wisp/internal/agent/approval 0.372s`，
`--- PASS` **90** 枚／`--- FAIL` **0** 枚／`=== RUN` 91。

⚠ **顶回派单/票面转述的枚数**：派单与本票 §5 反复写「66 全绿」＝**10-03 的快照**，
今天同名尺件已长到 **90 PASS**。下面每一发的"绿/红"都按本腿现量报，⛔ 不引 66 当基线。

本腿只动的两族尺（现名册，逐枚取自上面的 `-v` 输出，非抄票面）：

- `TestTicket242BindDigestSeparatesItemsBySequenceNumber`
- `TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`
- `TestTicket242PanelItemReadFaceIsFullyDeclared`
- `TestTicket242PanelItemStaysGrantFree`
- `TestTicket242QueuedItemsBindGrantsToTheirOwnDigest`
- `TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce`
- `TestTicket242SpendRequiresTheExactBinding`
- `TestTicket259R1DenialNamesCardNotPending`
- `TestTicket259R1DenialNamesMisbound`
- `TestTicket259R1DenialNamesMissingNonce`
- `TestTicket259R1DenialNamesSpentOrNeverLiveNonce`
- `TestTicket259R1FourDenialLabelsArePairwiseDistinct`
- `TestTicket259R1OutwardAnswerStaysMergedAcrossDenials`
- `TestTicket259R1PanelAPIMethodSetIsClosed`
- `TestTicket259R1PanelCarrierCarriesNoAnswerVerb`
- `TestTicket259R1PanelItemCarriesNoCredentialValue`
- `TestTicket259R1PanelItemFieldCapabilityFence`
- `TestTicket259R1PanelItemProjectionIgnoresGrantState`
- `TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer`

整包 90 枚用例名逐字（机器追加，取法＝`grep -- '--- PASS' <输出> | sed`）：
```
  TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses
  TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard
  TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant
  TestAnAliasCanNeverBuyAnAllow
  TestAnAmbiguousVetoNameIsNotGuessedOnTheRefusalSide
  TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline
  TestAnUnreachablePromptSurfaceIsRefusedImmediatelyNotWaitedOut
  TestBatchAggregatesHomogeneousL1Ops
  TestBatchRefusalsBelowThresholdAndOnSensitiveVerdicts
  TestDefaultQueueDeadlineIsFiniteAtThreeHundredSeconds
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/DecideFromNative.allow
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/DecideFromNative.reject
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/DecideFromPanel.reject
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/Native.Allow
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/Native.Reject
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/Panel.Reject
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/Veto.empty
  TestEveryAnswerRouteForAnUnknownCorrelationFailsFast/Veto.unknown
  TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
  TestEveryRefusalRouteOnAnUnknownEntryStillRefuses/decide_from_native
  TestEveryRefusalRouteOnAnUnknownEntryStillRefuses/decide_from_panel
  TestEveryRefusalRouteOnAnUnknownEntryStillRefuses/native_reject
  TestEveryRefusalRouteOnAnUnknownEntryStillRefuses/panel_reject
  TestEveryRefusalRouteOnAnUnknownEntryStillRefuses/veto
  TestGrantIsSingleUseAndBoundToItsItem
  TestHostUnreachableAndFullQueueFailClosed
  TestHostUnreachableAndFullQueueFailClosed/宿主不可达
  TestHostUnreachableAndFullQueueFailClosed/队列已满
  TestL1WindowIsPlainDataWithNoRouteToAnAnswer
  TestL1WindowTimeoutMeansExecute
  TestL1WindowUnavailableVoiceChannelIsAnnounced
  TestL1WindowUnreachableUIFailsClosed
  TestL1WindowVetoOnForeignCorrelationCannotCancel
  TestL1WindowVetoedByEachLoadedChannel
  TestL1WindowVetoedByEachLoadedChannel/ball
  TestL1WindowVetoedByEachLoadedChannel/esc
  TestL2NativeAllowExecutesThroughTheBridge
  TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s
  TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord
  TestLiveApprovalsOnNilQueueIsNil
  TestLiveApprovalsReportsPendingRowsAndDoesNotConsumeThem
  TestLiveApprovalsSharesNoReferenceSlotWithTheQueue
  TestLiveApprovalsSkipsRowsThatAreNotPending
  TestLiveL1WindowsSeesAWindowInTheLoop
  TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis
  TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis/面板声称原生并携带真实令牌
  TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis/面板声称自己是原生，且不带令牌
  TestR7SizedBatchEscalatesToL2Unaggregated
  TestR7SizedBatchThatNobodyAnsweredIsRejectedNotRun
  TestReadingL1WindowsMovesNoGateState
  TestReplayRedisplaysUnderAFreshGrant
  TestReplySeamCarriesNoAllowWithoutAGrant
  TestReplySeamWaitingStateNamesTheTwoD43Rows
  TestStartedCallsAreNotReportedAsWaiting
  TestTenOpsInOneToolCallGetOneConfirm
  TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNamed
  TestTicket224ForgedSessionAnswerRecordsNothing
  TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow
  TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped
  TestTicket242BindDigestSeparatesItemsBySequenceNumber
  TestTicket242ForgedBindingCannotSpendAnotherItemsGrant
  TestTicket242PanelItemReadFaceIsFullyDeclared
  TestTicket242PanelItemStaysGrantFree
  TestTicket242QueuedItemsBindGrantsToTheirOwnDigest
  TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce
  TestTicket242SpendRequiresTheExactBinding
  TestTicket259R1DenialNamesCardNotPending
  TestTicket259R1DenialNamesMisbound
  TestTicket259R1DenialNamesMissingNonce
  TestTicket259R1DenialNamesSpentOrNeverLiveNonce
  TestTicket259R1FourDenialLabelsArePairwiseDistinct
  TestTicket259R1OutwardAnswerStaysMergedAcrossDenials
  TestTicket259R1PanelAPIMethodSetIsClosed
  TestTicket259R1PanelCarrierCarriesNoAnswerVerb
  TestTicket259R1PanelItemCarriesNoCredentialValue
  TestTicket259R1PanelItemFieldCapabilityFence
  TestTicket259R1PanelItemProjectionIgnoresGrantState
  TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer
  TestTicket260R3EscUnavailableLineNamesNoKey
  TestTicket260R3LoadedCancelLabelIsTheNamedResidual
  TestTicket260R4DefaultLabelIsTheOldSentence
  TestTicket260R4LabelFollowsTheReaderAcrossAReload
  TestTicket260R4SeededKeyReplacesEscOnEveryFace
  TestTicket260R4SpellingSeamCarriesNoAuthority
  TestTicket260R4UnloadLineNamesNoKey
  TestUnadmittedTaskIsRefusedBeforeQueueing
  TestVetoAfterStartYieldsAppliedStepsReport
  TestVetoTravelsBackToTheLoopAsAFailure
  TestWideningRuleReadsOneLineOffTheCard
```
