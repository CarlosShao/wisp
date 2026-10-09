time=2026-10-09T12:01:00.852+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestCompressFoldsOldestKeepsLastThreePreservesIDs
time=2026-10-09T12:01:00.855+08:00 level=INFO msg="agent: history compressed" tokens_before=30090 tokens_after=6586 threshold=12000 msgs_before=45 msgs_after=10 compressed_msgs=36 kept_raw_rounds=3 history_changed=true
--- PASS: TestCompressFoldsOldestKeepsLastThreePreservesIDs (0.00s)
=== RUN   TestCompressPreservesIDsEvenWhenSummarizerDropsThem
time=2026-10-09T12:01:00.856+08:00 level=INFO msg="agent: history compressed" tokens_before=30090 tokens_after=6108 threshold=12000 msgs_before=45 msgs_after=10 compressed_msgs=36 kept_raw_rounds=3 history_changed=true
--- PASS: TestCompressPreservesIDsEvenWhenSummarizerDropsThem (0.00s)
=== RUN   TestCompressTriggerScalesWithWindow
time=2026-10-09T12:01:00.856+08:00 level=INFO msg="agent: history compressed" tokens_before=505 tokens_after=401 threshold=384 msgs_before=10 msgs_after=7 compressed_msgs=4 kept_raw_rounds=3 history_changed=true
--- PASS: TestCompressTriggerScalesWithWindow (0.00s)
=== RUN   TestCompressNoopUnderThreshold
--- PASS: TestCompressNoopUnderThreshold (0.00s)
=== RUN   TestCompressionTraceBooksCountsOnSuccess
--- PASS: TestCompressionTraceBooksCountsOnSuccess (0.00s)
=== RUN   TestCompressionTraceSilentWhenNothingFolded
--- PASS: TestCompressionTraceSilentWhenNothingFolded (0.00s)
=== RUN   TestCompressionTraceSurvivesTheLoopWiring
--- PASS: TestCompressionTraceSurvivesTheLoopWiring (0.01s)
=== RUN   TestCompressionTraceDoesNotAlterTheFold
time=2026-10-09T12:01:00.866+08:00 level=INFO msg="agent: history compressed" tokens_before=1236 tokens_after=780 threshold=384 msgs_before=18 msgs_after=10 compressed_msgs=9 kept_raw_rounds=3 history_changed=true
--- PASS: TestCompressionTraceDoesNotAlterTheFold (0.00s)
=== RUN   TestCompressionTraceSilentWhenNothingFoldableOverThreshold
--- PASS: TestCompressionTraceSilentWhenNothingFoldableOverThreshold (0.00s)
=== RUN   TestCompressionTraceCarriesTheOwningTaskID
--- PASS: TestCompressionTraceCarriesTheOwningTaskID (0.01s)
=== RUN   TestCompressionTraceNeverInventsATaskID
--- PASS: TestCompressionTraceNeverInventsATaskID (0.00s)
=== RUN   TestControlWordsNeverReachTheProvider
--- PASS: TestControlWordsNeverReachTheProvider (0.01s)
=== RUN   TestNearMissUtterancesGoToTheLoop
--- PASS: TestNearMissUtterancesGoToTheLoop (0.00s)
=== RUN   TestDefaultControlCancelsRunningTask
--- PASS: TestDefaultControlCancelsRunningTask (0.00s)
=== RUN   TestControlWithoutAConsumerIsVisible
--- PASS: TestControlWithoutAConsumerIsVisible (0.00s)
=== RUN   TestSteeringInsertsIntoNextRound
--- PASS: TestSteeringInsertsIntoNextRound (0.40s)
=== RUN   TestSteeringDisabled
--- PASS: TestSteeringDisabled (0.00s)
=== RUN   TestRecordSinkToleratesConcurrentPublishers
--- PASS: TestRecordSinkToleratesConcurrentPublishers (0.00s)
=== RUN   TestCorrPerCallTwoAsksSameTask242
--- PASS: TestCorrPerCallTwoAsksSameTask242 (0.00s)
=== RUN   TestDeclaredL0RunsWhenPassThroughIsOff179
--- PASS: TestDeclaredL0RunsWhenPassThroughIsOff179 (0.07s)
=== RUN   TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179
--- PASS: TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179 (0.07s)
=== RUN   TestDeclaredL1L2RoutingUnchanged179
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L1/no-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L1/with-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L2/no-approval-layer
=== RUN   TestDeclaredL1L2RoutingUnchanged179/L2/with-approval-layer
--- PASS: TestDeclaredL1L2RoutingUnchanged179 (0.06s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L1/no-approval-layer (0.01s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L1/with-approval-layer (0.00s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L2/no-approval-layer (0.00s)
    --- PASS: TestDeclaredL1L2RoutingUnchanged179/L2/with-approval-layer (0.00s)
=== RUN   TestFailedTaskBooksOpenCallRowWithDecision
--- PASS: TestFailedTaskBooksOpenCallRowWithDecision (0.06s)
=== RUN   TestCancelledTaskPersistsTerminalRows
--- PASS: TestCancelledTaskPersistsTerminalRows (0.07s)
=== RUN   TestLoopGuardLadderRemindersThenStuck
--- PASS: TestLoopGuardLadderRemindersThenStuck (0.01s)
=== RUN   TestLoopGuardLadderIsConfigurable
--- PASS: TestLoopGuardLadderIsConfigurable (0.00s)
=== RUN   TestLoopGuardResetsOnDifferentCall
--- PASS: TestLoopGuardResetsOnDifferentCall (0.00s)
=== RUN   TestLoopGuardTokenBudgetScalesWithWindow
--- PASS: TestLoopGuardTokenBudgetScalesWithWindow (0.00s)
=== RUN   TestPerToolTimeoutFires
--- PASS: TestPerToolTimeoutFires (0.12s)
=== RUN   TestPerToolTimeoutOfContractHonestToolIsToolClass
--- PASS: TestPerToolTimeoutOfContractHonestToolIsToolClass (0.12s)
=== RUN   TestLoopGuardCanonicalizesEquivalentArgs
--- PASS: TestLoopGuardCanonicalizesEquivalentArgs (0.00s)
=== RUN   TestLoopGuardKeepsDistinctArgsDistinct
--- PASS: TestLoopGuardKeepsDistinctArgsDistinct (0.00s)
=== RUN   TestLoopGuardRepeatsNonJSONArgs
--- PASS: TestLoopGuardRepeatsNonJSONArgs (0.00s)
=== RUN   TestProjectInstructionsReachTheWireRequestBody
--- PASS: TestProjectInstructionsReachTheWireRequestBody (0.01s)
=== RUN   TestAssemblySplitsGuidanceFromUntrustedData
--- PASS: TestAssemblySplitsGuidanceFromUntrustedData (0.01s)
=== RUN   TestHostileInstructionFileMovesNoAuthorityKnob
--- PASS: TestHostileInstructionFileMovesNoAuthorityKnob (0.01s)
=== RUN   TestProjectInstructionsAreNotANewD39Section
--- PASS: TestProjectInstructionsAreNotANewD39Section (0.00s)
=== RUN   TestGoldenTextReply
--- PASS: TestGoldenTextReply (0.00s)
=== RUN   TestGoldenSingleToolCall
--- PASS: TestGoldenSingleToolCall (0.00s)
=== RUN   TestGoldenParallelToolCalls
--- PASS: TestGoldenParallelToolCalls (0.00s)
=== RUN   TestGoldenToolResultFeedsNextTurn
--- PASS: TestGoldenToolResultFeedsNextTurn (0.00s)
=== RUN   TestGoldenConcurrencyCeiling
--- PASS: TestGoldenConcurrencyCeiling (0.13s)
=== RUN   TestToolExecutionRunsFourAcross
--- PASS: TestToolExecutionRunsFourAcross (0.01s)
=== RUN   TestGoldenBudgetExhaustionStuck
--- PASS: TestGoldenBudgetExhaustionStuck (0.01s)
=== RUN   TestGoldenCancellationEndToEnd
--- PASS: TestGoldenCancellationEndToEnd (0.00s)
=== RUN   TestGoroutineBudgetDrains
--- PASS: TestGoroutineBudgetDrains (0.00s)
=== RUN   TestTaskLogAndToolCallRows
time=2026-10-09T12:01:02.102+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTaskLogAndToolCallRows172428908\001\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:02.108+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:02Z duration_ms=5
time=2026-10-09T12:01:02.113+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTaskLogAndToolCallRows172428908\001\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:02.118+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:02Z duration_ms=5
time=2026-10-09T12:01:02.123+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:02.124+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTaskLogAndToolCallRows (0.05s)
=== RUN   TestGoldenProviderFailureIsVisible
--- PASS: TestGoldenProviderFailureIsVisible (0.00s)
=== RUN   TestCachePrefixIsByteStableAcrossTurns
--- PASS: TestCachePrefixIsByteStableAcrossTurns (0.00s)
=== RUN   TestPromptSectionCanonicalOrder
--- PASS: TestPromptSectionCanonicalOrder (0.00s)
=== RUN   TestVolatileSuffixPutsSceneLast
--- PASS: TestVolatileSuffixPutsSceneLast (0.00s)
=== RUN   TestD39SectionBudgetsAndTotalCeiling
--- PASS: TestD39SectionBudgetsAndTotalCeiling (0.00s)
=== RUN   TestEnforceTotalTerminatesAtEveryWindow
--- PASS: TestEnforceTotalTerminatesAtEveryWindow (0.00s)
=== RUN   TestEnforceTotalTerminatesOnABelowBudgetSection
--- PASS: TestEnforceTotalTerminatesOnABelowBudgetSection (0.00s)
=== RUN   TestAC3SpillArtifactLandsPrivate
    spill_acl_windows_test.go:29: icacls tool-output-call_acl.txt -> [NT AUTHORITY\SYSTEM BUILTIN\Administrators DESKTOP-LVS7839\swq]
    spill_acl_windows_test.go:30: icacls artifacts -> [NT AUTHORITY\SYSTEM NT AUTHORITY\SYSTEM BUILTIN\Administrators BUILTIN\Administrators DESKTOP-LVS7839\swq DESKTOP-LVS7839\swq]
    spill_acl_windows_test.go:37: icacls tool-output-call_acl.txt -> [NT AUTHORITY\SYSTEM BUILTIN\Administrators DESKTOP-LVS7839\swq]
--- PASS: TestAC3SpillArtifactLandsPrivate (0.07s)
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/p_slash_vs_pq
    spill_name_injectivity_test.go:90: names: map[tool-output-p%2Fq.txt:p/q tool-output-pq.txt:pq]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/p_backslash_vs_pq
    spill_name_injectivity_test.go:90: names: map[tool-output-p%5Cq.txt:p\q tool-output-pq.txt:pq]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/p_slash_vs_p_backslash
    spill_name_injectivity_test.go:90: names: map[tool-output-p%2Fq.txt:p/q tool-output-p%5Cq.txt:p\q]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/the_whole_report_triple
    spill_name_injectivity_test.go:90: names: map[tool-output-p%2Fq.txt:p/q tool-output-p%5Cq.txt:p\q tool-output-pq.txt:pq]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/dot_and_dotdot_vs_empty
    spill_name_injectivity_test.go:90: names: map[tool-output-%2E%2E.txt:.. tool-output-%2E.txt:. tool-output-seq.1.txt:]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/colon_and_unc_shape
    spill_name_injectivity_test.go:90: names: map[tool-output-%5C%5Ca%5Cb.txt:\\a\b tool-output-a%3Ab.txt:a:b]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/underscore_dash_and_percent
    spill_name_injectivity_test.go:90: names: map[tool-output-a%25b.txt:a%b tool-output-a-b.txt:a-b tool-output-a_b.txt:a_b]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/case_folding
    spill_name_injectivity_test.go:90: names: map[tool-output-%50%51.txt:PQ tool-output-pq.txt:pq]
=== RUN   TestArtifactNameDoesNotFoldDistinctIDs/digits_vs_letters
    spill_name_injectivity_test.go:90: names: map[tool-output-0o.txt:0o tool-output-o0.txt:o0]
--- PASS: TestArtifactNameDoesNotFoldDistinctIDs (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/p_slash_vs_pq (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/p_backslash_vs_pq (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/p_slash_vs_p_backslash (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/the_whole_report_triple (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/dot_and_dotdot_vs_empty (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/colon_and_unc_shape (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/underscore_dash_and_percent (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/case_folding (0.00s)
    --- PASS: TestArtifactNameDoesNotFoldDistinctIDs/digits_vs_letters (0.00s)
=== RUN   TestSpilledBytesSurviveANameThatUsedToCollide
=== RUN   TestSpilledBytesSurviveANameThatUsedToCollide/slash_id_then_bare_id
=== RUN   TestSpilledBytesSurviveANameThatUsedToCollide/bare_id_then_backslash_id
=== RUN   TestSpilledBytesSurviveANameThatUsedToCollide/dotdot_id_then_word_id
--- PASS: TestSpilledBytesSurviveANameThatUsedToCollide (0.03s)
    --- PASS: TestSpilledBytesSurviveANameThatUsedToCollide/slash_id_then_bare_id (0.01s)
    --- PASS: TestSpilledBytesSurviveANameThatUsedToCollide/bare_id_then_backslash_id (0.01s)
    --- PASS: TestSpilledBytesSurviveANameThatUsedToCollide/dotdot_id_then_word_id (0.01s)
=== RUN   TestArtifactNameRoundTripsToTheExactID
--- PASS: TestArtifactNameRoundTripsToTheExactID (0.00s)
=== RUN   TestArtifactNameIsCanonicalAndBounded
--- PASS: TestArtifactNameIsCanonicalAndBounded (0.00s)
=== RUN   TestWriteFileExclusiveIsExclusive
--- PASS: TestWriteFileExclusiveIsExclusive (0.00s)
=== RUN   TestSpillSameIDRetryOverwrites
--- PASS: TestSpillSameIDRetryOverwrites (0.01s)
=== RUN   TestSpillAcrossRestartsKeepsRetrySemantics
--- PASS: TestSpillAcrossRestartsKeepsRetrySemantics (0.01s)
=== RUN   TestSpillCallIDHostileShapesSanitizedToBareNames
=== RUN   TestSpillCallIDHostileShapesSanitizedToBareNames/separator
=== RUN   TestSpillCallIDHostileShapesSanitizedToBareNames/dotdot
=== RUN   TestSpillCallIDHostileShapesSanitizedToBareNames/drive_letter
=== RUN   TestSpillCallIDHostileShapesSanitizedToBareNames/unc
=== RUN   TestSpillCallIDHostileShapesSanitizedToBareNames/encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence
--- PASS: TestSpillCallIDHostileShapesSanitizedToBareNames (0.06s)
    --- PASS: TestSpillCallIDHostileShapesSanitizedToBareNames/separator (0.01s)
    --- PASS: TestSpillCallIDHostileShapesSanitizedToBareNames/dotdot (0.01s)
    --- PASS: TestSpillCallIDHostileShapesSanitizedToBareNames/drive_letter (0.01s)
    --- PASS: TestSpillCallIDHostileShapesSanitizedToBareNames/unc (0.01s)
    --- PASS: TestSpillCallIDHostileShapesSanitizedToBareNames/encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence (0.02s)
=== RUN   TestSpillContainmentByDirectoryListing
    spill_path_invariant_test.go:284: control: escape id "..\\..\\..\\..\\CONTROL-escape" resolves to "C:\\Users\\swq\\AppData\\Local\\Temp\\TestSpillContainmentByDirectoryListing1502504596\\001\\CONTROL-escape.txt" (root is "C:\\Users\\swq\\AppData\\Local\\Temp\\TestSpillContainmentByDirectoryListing1502504596\\001")
    spill_path_invariant_test.go:307: literal-backslash control: id "..\\..\\..\\..\\CONTROL-literal" resolves to "C:\\Users\\swq\\AppData\\Local\\Temp\\TestSpillContainmentByDirectoryListing1502504596\\001\\CONTROL-literal.txt"
--- PASS: TestSpillContainmentByDirectoryListing (0.05s)
=== RUN   TestSpillIntoRealStoreThenDeleteStaysUnderDataDir
--- PASS: TestSpillIntoRealStoreThenDeleteStaysUnderDataDir (0.08s)
=== RUN   TestSpillAPITakesNoCallerControlledDestinationPath
    spill_path_invariant_test.go:545: audited 3 filesystem-mutating call(s) in spill.go
--- PASS: TestSpillAPITakesNoCallerControlledDestinationPath (0.00s)
=== RUN   TestSpillPointerAuthorityAsksWithTheCanonicalAnswer174r4
    spill_pointer_authority_ac3_174r4_test.go:87: arm 1 verbatim: Text="[…输出已落文件：省略 65221 字符，总长 68021 字节 / 约 17005 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestSpillPointerAuthorityAsksWithTheCanonicalAnswer174r43944670815\\001\\tool-output-call_honesty.txt…]"
    spill_pointer_authority_ac3_174r4_test.go:110: arm 2 verbatim: Text="[…输出已落文件：省略 65221 字符，总长 68021 字节 / 约 17005 token；注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestSpillPointerAuthorityAsksWithTheCanonicalAnswer174r43944670815\\002\\tool-output-call_honesty.txt…]"
--- PASS: TestSpillPointerAuthorityAsksWithTheCanonicalAnswer174r4 (0.01s)
=== RUN   TestSpillPointerUnwiredJudgeFailsClosed174
--- PASS: TestSpillPointerUnwiredJudgeFailsClosed174 (0.00s)
=== RUN   TestSpillPointerOutsideAllowlistNamesBothRoadsBack174
--- PASS: TestSpillPointerOutsideAllowlistNamesBothRoadsBack174 (0.00s)
=== RUN   TestSpillPointerCanonicalizeFailureQuotesNoUpstreamText174
--- PASS: TestSpillPointerCanonicalizeFailureQuotesNoUpstreamText174 (0.01s)
=== RUN   TestSpillHealthyPointerStillMatchesThePreFixSentence174
--- PASS: TestSpillHealthyPointerStillMatchesThePreFixSentence174 (0.01s)
=== RUN   TestSpillPointerNoticeAddsNoSecondPath174
=== RUN   TestSpillPointerNoticeAddsNoSecondPath174/outRoot
=== RUN   TestSpillPointerNoticeAddsNoSecondPath174/canonErr
=== RUN   TestSpillPointerNoticeAddsNoSecondPath174/unwired
--- PASS: TestSpillPointerNoticeAddsNoSecondPath174 (0.02s)
    --- PASS: TestSpillPointerNoticeAddsNoSecondPath174/outRoot (0.00s)
    --- PASS: TestSpillPointerNoticeAddsNoSecondPath174/canonErr (0.01s)
    --- PASS: TestSpillPointerNoticeAddsNoSecondPath174/unwired (0.01s)
=== RUN   TestSpillPointerNoticeKeepsHeadTailAndBudget174
--- PASS: TestSpillPointerNoticeKeepsHeadTailAndBudget174 (0.01s)
=== RUN   TestLoopSpillNoticeReachesHistoryUnwired174
--- PASS: TestLoopSpillNoticeReachesHistoryUnwired174 (0.01s)
=== RUN   TestLoopHonorsConfigPointerJudgeQuietly174
--- PASS: TestLoopHonorsConfigPointerJudgeQuietly174 (0.01s)
=== RUN   TestSpillThresholdScalesWithWindow
--- PASS: TestSpillThresholdScalesWithWindow (0.01s)
=== RUN   TestSpillTokenBoundary
--- PASS: TestSpillTokenBoundary (0.00s)
=== RUN   TestSpillArtifactAndStubShape
--- PASS: TestSpillArtifactAndStubShape (0.01s)
=== RUN   TestSpillRawHardCap
--- PASS: TestSpillRawHardCap (0.00s)
=== RUN   TestSpillThroughLoop
--- PASS: TestSpillThroughLoop (0.02s)
=== RUN   TestSpillArtifactRespectsRawCap
--- PASS: TestSpillArtifactRespectsRawCap (0.01s)
=== RUN   Test285CallCorrIsDistinctPerCall
--- PASS: Test285CallCorrIsDistinctPerCall (0.00s)
=== RUN   Test285LoopDispatchesOneCorrPerCall
--- PASS: Test285LoopDispatchesOneCorrPerCall (0.00s)
=== RUN   TestMaxTokensFailsAllToolCallsOfThatMessage
time=2026-10-09T12:01:02.598+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestMaxTokensFailsAllToolCallsOfThatMessage1564849954\001\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:02.604+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:02Z duration_ms=6
time=2026-10-09T12:01:02.611+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestMaxTokensFailsAllToolCallsOfThatMessage1564849954\001\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:02.615+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:02Z duration_ms=4
time=2026-10-09T12:01:02.619+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:02.622+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestMaxTokensFailsAllToolCallsOfThatMessage (0.06s)
=== RUN   TestOpenToolCallAfterFailureIsFailed
--- PASS: TestOpenToolCallAfterFailureIsFailed (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/agent	1.826s
time=2026-10-09T12:01:00.956+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile
=== RUN   TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile/overwrite_branch_target_keeps_every_old_byte
=== RUN   TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile/new_file_branch_target_never_appears
=== RUN   TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile/the_next_write_reclaims_the_staging_residue
    bridge_a18_kill_windows_test.go:230: 残留暂存文件 .wisp-tmp-4d49082e-32660-3054473418：4 字节（子进程被杀时已经写进暂存区的字节数）
--- PASS: TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile (0.45s)
    --- PASS: TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile/overwrite_branch_target_keeps_every_old_byte (0.22s)
    --- PASS: TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile/new_file_branch_target_never_appears (0.20s)
    --- PASS: TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile/the_next_write_reclaims_the_staging_residue (0.02s)
=== RUN   TestBridgeRefusesARealJunctionOnTheReadRoute
=== RUN   TestBridgeRefusesARealJunctionOnTheReadRoute/fs.read_through_a_real_junction
=== RUN   TestBridgeRefusesARealJunctionOnTheReadRoute/fs.list_through_a_real_junction
--- PASS: TestBridgeRefusesARealJunctionOnTheReadRoute (0.07s)
    --- PASS: TestBridgeRefusesARealJunctionOnTheReadRoute/fs.read_through_a_real_junction (0.00s)
    --- PASS: TestBridgeRefusesARealJunctionOnTheReadRoute/fs.list_through_a_real_junction (0.00s)
=== RUN   TestBridgeWritesNothingThroughARealJunction
=== RUN   TestBridgeWritesNothingThroughARealJunction/fs.write_over_the_junctioned_target
=== RUN   TestBridgeWritesNothingThroughARealJunction/fs.write_a_new_file_through_the_junction
=== RUN   TestBridgeWritesNothingThroughARealJunction/fs.trash_through_the_junction
=== RUN   TestBridgeWritesNothingThroughARealJunction/fs.move_through_the_junction
=== RUN   TestBridgeWritesNothingThroughARealJunction/fs.delete_through_the_junction
--- PASS: TestBridgeWritesNothingThroughARealJunction (0.32s)
    --- PASS: TestBridgeWritesNothingThroughARealJunction/fs.write_over_the_junctioned_target (0.07s)
    --- PASS: TestBridgeWritesNothingThroughARealJunction/fs.write_a_new_file_through_the_junction (0.07s)
    --- PASS: TestBridgeWritesNothingThroughARealJunction/fs.trash_through_the_junction (0.07s)
    --- PASS: TestBridgeWritesNothingThroughARealJunction/fs.move_through_the_junction (0.06s)
    --- PASS: TestBridgeWritesNothingThroughARealJunction/fs.delete_through_the_junction (0.06s)
=== RUN   TestJunctionInsideAnAllowedRootCannotReachAnAListFile
--- PASS: TestJunctionInsideAnAllowedRootCannotReachAnAListFile (0.06s)
=== RUN   TestBridgeRefusesTheRealShortNameOfAnAListFile
--- PASS: TestBridgeRefusesTheRealShortNameOfAnAListFile (0.01s)
=== RUN   TestShortNameSpellingGetsTheSameVerdictAsTheLongOne
=== RUN   TestShortNameSpellingGetsTheSameVerdictAsTheLongOne/an_approved_out_of_scope_short_name_reads_the_file_the_card_named
--- PASS: TestShortNameSpellingGetsTheSameVerdictAsTheLongOne (0.01s)
    --- PASS: TestShortNameSpellingGetsTheSameVerdictAsTheLongOne/an_approved_out_of_scope_short_name_reads_the_file_the_card_named (0.00s)
=== RUN   TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing
time=2026-10-09T12:01:01.967+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing13923931\003\data\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:01.974+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:01Z duration_ms=6
time=2026-10-09T12:01:01.981+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing13923931\003\data\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:01.986+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:01Z duration_ms=4
time=2026-10-09T12:01:01.990+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:01.992+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing (0.13s)
=== RUN   TestRealToolCallWritesRewriteAccountIntoAudit
--- PASS: TestRealToolCallWritesRewriteAccountIntoAudit (0.01s)
=== RUN   TestUnusableRootIsVisibleInTheAuditRecord
--- PASS: TestUnusableRootIsVisibleInTheAuditRecord (0.02s)
=== RUN   TestAccountRecordIsWrittenEvenWithNothingToReport
--- PASS: TestAccountRecordIsWrittenEvenWithNothingToReport (0.02s)
=== RUN   TestSensitiveReadAcrossTheBridgeOpensItsTaskScope
--- PASS: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope (0.02s)
=== RUN   TestC1ToolContract
--- PASS: TestC1ToolContract (0.00s)
=== RUN   TestToolConcurrencyCeilingIsFour
--- PASS: TestToolConcurrencyCeilingIsFour (0.05s)
=== RUN   TestPerToolTimeoutHonoredViaContext
--- PASS: TestPerToolTimeoutHonoredViaContext (0.04s)
=== RUN   TestCapabilitySetIsFrozen
--- PASS: TestCapabilitySetIsFrozen (0.00s)
=== RUN   TestUndeclaredCapabilityHardRejects
--- PASS: TestUndeclaredCapabilityHardRejects (0.00s)
=== RUN   TestC4SlotsWithoutAnImplementationAreNotSilent
--- PASS: TestC4SlotsWithoutAnImplementationAreNotSilent (0.00s)
=== RUN   TestC4RegistryRejectsBadDeclarations
--- PASS: TestC4RegistryRejectsBadDeclarations (0.00s)
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch/L0_passes_and_no_gate_is_consulted
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch/L1_uses_the_pending-window_callback
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch/L1_window_timeout_means_execute
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch/L1_veto_returns_a_failure_to_the_model,_not_to_the_user
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch/L2_uses_the_pending-approval_callback
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch/L2_approval_timeout_auto-rejects
=== RUN   TestRiskDecisionRoutesEachLevelToItsBranch/NoGate_fails_closed_above_L0
--- PASS: TestRiskDecisionRoutesEachLevelToItsBranch (0.00s)
    --- PASS: TestRiskDecisionRoutesEachLevelToItsBranch/L0_passes_and_no_gate_is_consulted (0.00s)
    --- PASS: TestRiskDecisionRoutesEachLevelToItsBranch/L1_uses_the_pending-window_callback (0.00s)
    --- PASS: TestRiskDecisionRoutesEachLevelToItsBranch/L1_window_timeout_means_execute (0.00s)
    --- PASS: TestRiskDecisionRoutesEachLevelToItsBranch/L1_veto_returns_a_failure_to_the_model,_not_to_the_user (0.00s)
    --- PASS: TestRiskDecisionRoutesEachLevelToItsBranch/L2_uses_the_pending-approval_callback (0.00s)
    --- PASS: TestRiskDecisionRoutesEachLevelToItsBranch/L2_approval_timeout_auto-rejects (0.00s)
    --- PASS: TestRiskDecisionRoutesEachLevelToItsBranch/NoGate_fails_closed_above_L0 (0.00s)
=== RUN   TestDeclaredRiskIsOnlyAFloor
--- PASS: TestDeclaredRiskIsOnlyAFloor (0.01s)
=== RUN   TestToolCallRowsAreComplete
time=2026-10-09T12:01:02.197+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestToolCallRowsAreComplete1424534539\001\data\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:02.203+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:02Z duration_ms=6
time=2026-10-09T12:01:02.212+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestToolCallRowsAreComplete1424534539\001\data\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:02.218+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:02Z duration_ms=6
time=2026-10-09T12:01:02.222+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:02.224+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestToolCallRowsAreComplete (0.07s)
=== RUN   TestBridgeIsTheLoopToolProvider
--- PASS: TestBridgeIsTheLoopToolProvider (0.00s)
=== RUN   Test236R2TaskCancelRefusesWhenHostGaveNoCallerID
time=2026-10-09T12:01:02.249+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-d5c5c840-db96-4ff3-9b70-c6a0b487f5e8 owner=tools
--- PASS: Test236R2TaskCancelRefusesWhenHostGaveNoCallerID (0.00s)
=== RUN   Test236R2TaskCancelRefusesWhenRosterIsUnwired
--- PASS: Test236R2TaskCancelRefusesWhenRosterIsUnwired (0.00s)
=== RUN   Test236R2TaskSpawnRefusesWhenRosterIsUnwired
--- PASS: Test236R2TaskSpawnRefusesWhenRosterIsUnwired (0.00s)
=== RUN   Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired
--- PASS: Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired (0.00s)
=== RUN   Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired
--- PASS: Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired (0.00s)
=== RUN   Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID
--- PASS: Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID (0.00s)
=== RUN   TestFSEditKeepsABOMItWasNotAskedToTouch
    fs_edit_ac34_test.go:103: AC#3 BOM 正向读数：31 字节 → 31 字节，前三字节仍为 EF BB BF，BOM 出现 1 次，文案="fs.edit 已改写 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditKeepsABOMItWasNotAskedToTouch1628529276\\001\\report.md：1 枚编辑全部生效，31 字节 / 2 行 → 31 字节 / 2 行（临时文件+原子重命名）"
--- PASS: TestFSEditKeepsABOMItWasNotAskedToTouch (0.01s)
=== RUN   TestFSEditOnACRLFFileRefusesAnLFSpelledOld
    fs_edit_ac34_test.go:161: AC#3 CRLF 反向读数：23 字节的纯 CRLF 文件里 LF 写法的 old 被响亮拒绝（文案含文件长度 23，证明读到了整文件），同一文件上不含换行的 old 立即命中并落盘 ⇒ 这是正确拒绝，不是读不到
--- PASS: TestFSEditOnACRLFFileRefusesAnLFSpelledOld (0.01s)
=== RUN   TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF
    fs_edit_ac34_test.go:202: AC#3 CRLF 正向读数：23 字节 → 30 字节，CRLF 对 3→4，裸 CR/裸 LF 均为 0（行尾未被改写），文案="fs.edit 已改写 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF746378507\\001\\crlf-hit.txt：1 枚编辑全部生效，23 字节 / 3 行 → 30 字节 / 4 行（临时文件+原子重命名）"
--- PASS: TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF (0.01s)
=== RUN   TestFSEditDoesNotRewriteAnLFFilesLineEndings
    fs_edit_ac34_test.go:250: AC#3 LF 双向读数：20 字节 → 27 字节，裸 LF 4 条、CR 0 枚（LF 文件的行尾没被改）
--- PASS: TestFSEditDoesNotRewriteAnLFFilesLineEndings (0.02s)
=== RUN   TestFSEditLineEndingForensicsRefusesRealSpellings
=== RUN   TestFSEditLineEndingForensicsRefusesRealSpellings/crlf_file_lf_spelled_old
    fs_edit_ac34_test.go:330: 归一取证：23 字节文件（行尾形 纯 CRLF）被 0 命中拒掉一枚归一后可命中 1 次的 old；落盘前后一字未变
=== RUN   TestFSEditLineEndingForensicsRefusesRealSpellings/cr_only_file_lf_spelled_old
    fs_edit_ac34_test.go:330: 归一取证：20 字节文件（行尾形 纯 CR）被 0 命中拒掉一枚归一后可命中 1 次的 old；落盘前后一字未变
=== RUN   TestFSEditLineEndingForensicsRefusesRealSpellings/mixed_file_old_over_the_crlf_line_lf_spelled
    fs_edit_ac34_test.go:330: 归一取证：22 字节文件（行尾形 混合(CRLF=2 CR=0 LF=1)）被 0 命中拒掉一枚归一后可命中 1 次的 old；落盘前后一字未变
=== RUN   TestFSEditLineEndingForensicsRefusesRealSpellings/mixed_file_old_over_the_lf_line_crlf_spelled
    fs_edit_ac34_test.go:330: 归一取证：22 字节文件（行尾形 混合(CRLF=2 CR=0 LF=1)）被 0 命中拒掉一枚归一后可命中 1 次的 old；落盘前后一字未变
=== RUN   TestFSEditLineEndingForensicsRefusesRealSpellings/lf_file_crlf_spelled_old
    fs_edit_ac34_test.go:330: 归一取证：20 字节文件（行尾形 纯 LF）被 0 命中拒掉一枚归一后可命中 1 次的 old；落盘前后一字未变
--- PASS: TestFSEditLineEndingForensicsRefusesRealSpellings (0.04s)
    --- PASS: TestFSEditLineEndingForensicsRefusesRealSpellings/crlf_file_lf_spelled_old (0.01s)
    --- PASS: TestFSEditLineEndingForensicsRefusesRealSpellings/cr_only_file_lf_spelled_old (0.01s)
    --- PASS: TestFSEditLineEndingForensicsRefusesRealSpellings/mixed_file_old_over_the_crlf_line_lf_spelled (0.01s)
    --- PASS: TestFSEditLineEndingForensicsRefusesRealSpellings/mixed_file_old_over_the_lf_line_crlf_spelled (0.01s)
    --- PASS: TestFSEditLineEndingForensicsRefusesRealSpellings/lf_file_crlf_spelled_old (0.01s)
=== RUN   TestFSEditSilentShapeChangesLandToday
=== RUN   TestFSEditSilentShapeChangesLandToday/an_old_that_carries_the_bom_is_refused_r3
    fs_edit_ac34_test.go:375: AC#3b① 翻转读数（r2 的缺陷读数①在此变红）：14 字节带 BOM 的文件，old 从偏移 0 起并含 BOM ⇒ 现在是硬拒、一字未落，前三字节仍 EF BB BF，IsError=true ErrorClass="tool"
=== RUN   TestFSEditSilentShapeChangesLandToday/an_lf_spelled_new_in_a_crlf_file_is_refused_r3
    fs_edit_ac34_test.go:404: AC#3b② 翻转读数（r2 的缺陷读数②在此变红）：纯 CRLF 文件里一枚 LF 写法的 new ⇒ 硬拒，落盘前后仍 3 组 CRLF / 0 裸 LF，IsError=true ErrorClass="tool"
--- PASS: TestFSEditSilentShapeChangesLandToday (0.02s)
    --- PASS: TestFSEditSilentShapeChangesLandToday/an_old_that_carries_the_bom_is_refused_r3 (0.01s)
    --- PASS: TestFSEditSilentShapeChangesLandToday/an_lf_spelled_new_in_a_crlf_file_is_refused_r3 (0.01s)
=== RUN   TestFSEditKilledMidWriteLeavesTheTargetByteIdentical
    fs_edit_ac34_test.go:495: AC#4(a)(b) 读数：改前 26 字节 → 杀点 write:16（边界在下一块落笔前触发，暂存件里此刻 8 字节）→ 目标仍是 26 字节原字节（字节断言现在跑在边界探针之前）；目录内 .wisp-tmp-* 残件 0 枚，台账 ["在 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditKilledMidWriteLeavesTheTargetByteIdentical2806235523\\001 创建临时文件 .wisp-tmp-4d49082e-15948-74209320" "在步骤「write:16」前停止：模拟进程在此刻被杀" "向临时文件写入 8 字节" "写入未完成：已停止于 write:16：模拟进程在此刻被杀" "删除临时文件 .wisp-tmp-4d49082e-15948-74209320（目标从头到尾未被改动）"]
--- PASS: TestFSEditKilledMidWriteLeavesTheTargetByteIdentical (0.01s)
=== RUN   TestFSEditWithoutTheStagedWriterWouldCorruptTheTarget
    fs_edit_ac34_test.go:598: AC#4(c) 反面实验读数：同一批字节、同一道 write:16 杀点，摘掉 temp+rename 之后目标从 26 字节变成 8 字节（"ALPHA-LO"），既不是原文也不是要落的内容；事后 fs.read 返回 IsError=false ErrorClass="" Truncated=false，把损坏原文照读照回 ⇒ 这发损坏从此看不见。装回 temp+rename（上一枚用例）则目标一字未动。
--- PASS: TestFSEditWithoutTheStagedWriterWouldCorruptTheTarget (0.02s)
=== RUN   TestFSEditRefusesAnOldThatSwallowsTheFileBOM
    fs_edit_ac3b_test.go:134: AC#3b① 改后读数（回执）：IsError=true ErrorClass="tool" AppliedSteps=0
    fs_edit_ac3b_test.go:136: AC#3b① 改后读数（回执文案）："fs.edit 拒绝执行：第 1 枚编辑的 old 从文件第 0 字节起，圈住了文件头那 3 字节 EF BB BF（UTF-8 BOM，文件 14 字节）。BOM 是文件自己的编码标记，不是内容：请把 old 和 new 都从 BOM 之后的第一个字节开始抄；old 必须与文件内容逐字一致，包含全部空白与换行（the old string must match exactly including all whitespace and newlines）。本次调用一个字也没有写，目标文件保持原状。"
    fs_edit_ac3b_test.go:137: AC#3b① 改后读数（审计行，外部可见）：tools: call kind=error task=task-1 corr=corr-1 tool=fs.edit risk=L2 decision=allow outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason="R1: 工具声明为下界（L2）"
    fs_edit_ac3b_test.go:154: AC#3b① 合法对照：14 字节 → 14 字节，前三字节仍是 EF BB BF，文案="fs.edit 已改写 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditRefusesAnOldThatSwallowsTheFileBOM3378615408\\001\\bom-zero.txt：1 枚编辑全部生效，14 字节 / 2 行 → 14 字节 / 2 行（临时文件+原子重命名）"
--- PASS: TestFSEditRefusesAnOldThatSwallowsTheFileBOM (0.02s)
=== RUN   TestFSEditStillLandsAnOldStartingAtByteZeroWithoutABOM
    fs_edit_ac3b_test.go:183: AC#3b① 合法对照（无 BOM 文件、old 从第 0 字节起）：11 字节 → 11 字节，未被拒绝
--- PASS: TestFSEditStillLandsAnOldStartingAtByteZeroWithoutABOM (0.01s)
=== RUN   TestFSEditCannotBeFedAPartialBOMThroughJSON
    fs_edit_ac3b_test.go:215: AC#3b① 射程读数：old 只带 BOM 的第 3 个字节时，落盘的是 9 字节、第 0 字节是 U+FEFF，传进来的已经不是那三个字节："fs.edit 拒绝执行：第 1 枚编辑的 old 在文件里找不到（0 命中，文件 9 字节）。old 必须与文件内容逐字一致，包含全部空白与换行（the old string must match exactly including all whitespace and newlines）；先用 fs.read 读回原文再照抄那一段。本次调用一个字也没有写，目标文件保持原状。"
--- PASS: TestFSEditCannotBeFedAPartialBOMThroughJSON (0.01s)
=== RUN   TestFSEditRefusesALFSpelledNewInAPureCRLFFile
    fs_edit_ac3b_test.go:251: AC#3b② 改后读数（回执）：IsError=true ErrorClass="tool" AppliedSteps=0
    fs_edit_ac3b_test.go:253: AC#3b② 改后读数（回执文案）："fs.edit 拒绝执行：第 1 枚编辑的 new 会改写这份文件的行尾约定：文件原本是纯 CRLF，改完会变成 CRLF 3 组 / 裸 CR 0 条 / 裸 LF 1 条。请把 new 里的换行按文件自己的行尾拼写（纯 CRLF 的文件：换行写成 \\r\\n，回车加换行两字节）；确实要把整份文件换成别的行尾，请改用 fs.write 整文件写回，本工具不猜行尾。old 必须与文件内容逐字一致，包含全部空白与换行（the old string must match exactly including all whitespace and newlines）。本次调用一个字也没有写，目标文件保持原状。"
    fs_edit_ac3b_test.go:254: AC#3b② 改后读数（审计行，外部可见）：tools: call kind=error task=task-1 corr=corr-1 tool=fs.edit risk=L2 decision=allow outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason="R1: 工具声明为下界（L2）"
    fs_edit_ac3b_test.go:272: AC#3b② 合法对照一（同一枚编辑、new 按 CRLF 拼写）：23 字节 → 31 字节，CRLF 3→4 组、裸 CR/裸 LF 均 0
--- PASS: TestFSEditRefusesALFSpelledNewInAPureCRLFFile (0.02s)
=== RUN   TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak
=== RUN   TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak/lf_file_lands_an_lf_spelled_new
    fs_edit_ac3b_test.go:300: AC#3b② 合法对照二（LF 文件 + LF 写法 new）：20 字节 → 27 字节，未被拒绝
=== RUN   TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak/mixed_file_is_not_refused
    fs_edit_ac3b_test.go:321: AC#3b② 合法对照三（已混合的文件不被告知该长成什么样）：22 字节 → 29 字节，行尾形 混合(CRLF=2 CR=0 LF=2)
--- PASS: TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak (0.02s)
    --- PASS: TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak/lf_file_lands_an_lf_spelled_new (0.01s)
    --- PASS: TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak/mixed_file_is_not_refused (0.01s)
=== RUN   TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten
    fs_edit_ac3b_test.go:382: 反向尺读数：22 字节混合文件（混合(CRLF=2 CR=0 LF=1)）改完 22 字节，CRLF 仍 2 组、裸 LF 仍 1 条、裸 CR 0 枚，文案="fs.edit 已改写 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten2695247753\\001\\mixed-reverse-ruler.txt：1 枚编辑全部生效，22 字节 / 3 行 → 22 字节 / 3 行（临时文件+原子重命名）"
--- PASS: TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten (0.01s)
=== RUN   TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue
    fs_edit_ac4b_kill_windows_test.go:200: AC#4b 子进程输出（应只有 taskkill 的强制终止，没有任何 Go 清理）："time=2026-10-09T12:01:02.509+08:00 level=INFO msg=\"winsec: sealing path resolver installed\" resolver=risk.c26Pipeline probes_passed=2\n"
    fs_edit_ac4b_kill_windows_test.go:220: AC#4b(a) 字节级读数：目标 31 字节＝改前原文（"alpha\nbravo\ncharlie\ndelta\necho\n"），要落的 74 字节 "ALPHA-KILLED-MID-WITHOUT-RUNNING-ANY-GO-CLEANUP\n\nbravo\ncharlie\ndelta\necho\n" 一字未出现在盘上目标里
    fs_edit_ac4b_kill_windows_test.go:253: AC#4b(b) 残件点名：.wisp-tmp-4d49082e-15680-3385571916（目录 C:\Users\swq\AppData\Local\Temp\TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResi1019627690\001\edit 内）：4 字节、mtime 12:01:02，内容 "ALPH" ＝要落的新内容的前 4 字节；本程不删（只建不删），它是可回收的垃圾不是损坏
    fs_edit_ac4b_kill_windows_test.go:272: AC#4b(c) 读数：事后 fs.read IsError=false ErrorClass=""，内容 "alpha\nbravo\ncharlie\ndelta\necho\n"（原文完整）；目录里 1 枚残件（[C:\Users\swq\AppData\Local\Temp\TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResi1019627690\001\edit\.wisp-tmp-4d49082e-15680-3385571916]，各 [4] 字节）仍在盘上等人回收
--- PASS: TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue (0.23s)
=== RUN   TestFSEditAndFSWriteShareTheOutOfScopeVerdict
    fs_edit_ac5_gate_r4_test.go:191: 并排读数 fs.write: level=L2 rules_hit=[R1 R2 R8] window=0 approval=1 class="user_rejected" reason="R1: 工具声明为下界（L1）; R2: 目标路径在授权目录之外: C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditAndFSWriteShareTheOutOfScopeVerdict1051713326\\002\\note.txt; R8: 不可逆操作（覆盖已有内容）"
    fs_edit_ac5_gate_r4_test.go:193: 并排读数 fs.write: 回执="L2 审批未通过，已拒绝执行"
    fs_edit_ac5_gate_r4_test.go:191: 并排读数 fs.edit: level=L2 rules_hit=[R1 R2] window=0 approval=2 class="user_rejected" reason="R1: 工具声明为下界（L2）; R2: 目标路径在授权目录之外: C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditAndFSWriteShareTheOutOfScopeVerdict1051713326\\002\\note.txt"
    fs_edit_ac5_gate_r4_test.go:193: 并排读数 fs.edit: 回执="L2 审批未通过，已拒绝执行"
    fs_edit_ac5_gate_r4_test.go:207: 并排审计行 fs.write：tools: call kind=refused task=task-1 corr=corr-1 tool=fs.write risk=L2 decision=reject outcome=error rules_hit=[R1 R2 R8] in_allowlist_scope=false grant_id=0 reason="R1: 工具声明为下界（L1）; R2: 目标路径在授权目录之外: C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditAndFSWriteShareTheOutOfScopeVerdict1051713326\\002\\note.txt; R8: 不可逆操作（覆盖已有内容）"
    fs_edit_ac5_gate_r4_test.go:207: 并排审计行 fs.edit：tools: call kind=refused task=task-1 corr=corr-1 tool=fs.edit risk=L2 decision=reject outcome=error rules_hit=[R1 R2] in_allowlist_scope=false grant_id=0 reason="R1: 工具声明为下界（L2）; R2: 目标路径在授权目录之外: C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditAndFSWriteShareTheOutOfScopeVerdict1051713326\\002\\note.txt"
--- PASS: TestFSEditAndFSWriteShareTheOutOfScopeVerdict (0.01s)
=== RUN   TestFSEditRoutesToApprovalNotTheL1Window
--- PASS: TestFSEditRoutesToApprovalNotTheL1Window (0.02s)
=== RUN   TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord
    fs_edit_ac5_gate_r4_test.go:343: 读数：越界 fs.edit 经 R2 判 L2、卡片 reason="R1: 工具声明为下界（L2）; R2: 目标路径在授权目录之外: C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord1139993222\\002\\quiet.txt"，人批准后落盘
--- PASS: TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord (0.01s)
=== RUN   TestFSEditRefusesAPathC26CannotCanonicalize
    fs_edit_ac5_gate_r4_test.go:384: 读数（未修码即响）：空白 path ⇒ IsError=true class="tool" 回执="缺少 path 参数"
    fs_edit_ac5_gate_r4_test.go:389: 审计行：tools: call kind=error task=task-1 corr=corr-1 tool=fs.edit risk=L2 decision=allow outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason="R1: 工具声明为下界（L2）"
--- PASS: TestFSEditRefusesAPathC26CannotCanonicalize (0.00s)
=== RUN   TestFSEditWriteReclaimsADeadWritersOrphan
    fs_edit_ac5_sweep_r4_windows_test.go:77: 读数：fs.edit 一次成功写盘带走死进程孤儿 .wisp-tmp-4d49082e-30264-4242；台账="清扫 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditWriteReclaimsADeadWritersOrphan2885405200\\001 里本程序留下的历史暂存文件（上次写盘被中断的残口）：清扫 1 个，跳过 0 个；跳过的都不是孤儿：创建者进程仍在运行、正被占用、或不是可归因的普通文件"
--- PASS: TestFSEditWriteReclaimsADeadWritersOrphan (0.06s)
=== RUN   TestFSEditRefusesANonLiteralOld
    fs_edit_test.go:130: 判据 1+2 读数：命中前 50 字节 / 命中失败后 50 字节（一字未变），文案="fs.edit 拒绝执行：第 1 枚编辑的 old 在文件里找不到（0 命中，文件 50 字节）。old 必须与文件内容逐字一致，包含全部空白与换行（the old string must match exactly including all whitespace and newlines）；先用 fs.read 读回原文再照抄那一段。本次调用一个字也没有写，目标文件保持原状。"
--- PASS: TestFSEditRefusesANonLiteralOld (0.01s)
=== RUN   TestFSEditRefusesAnAmbiguousOldIsADistinctError
    fs_edit_test.go:171: 判据 3 读数：old 在 40 字节文件里命中 2 处，先数后拒，落盘前后一字未变
--- PASS: TestFSEditRefusesAnAmbiguousOldIsADistinctError (0.01s)
=== RUN   TestFSEditRefusesAnEmptyOldIsAWholeFileInsert
    fs_edit_test.go:210: 判据 4 空 old 读数：两形（old:"" 与缺 old 字段）都拒，文件保持 29 字节
--- PASS: TestFSEditRefusesAnEmptyOldIsAWholeFileInsert (0.01s)
=== RUN   TestFSEditIsAllOrNothingAcrossABatch
=== RUN   TestFSEditIsAllOrNothingAcrossABatch/no_op_edit_in_the_batch
    fs_edit_test.go:282: 全有或全无读数：本批 2 枚编辑，其中一枚不成立（空批＝整批不成立）⇒ 一个字也不落盘，文件仍是 26 字节
=== RUN   TestFSEditIsAllOrNothingAcrossABatch/overlapping_edits
    fs_edit_test.go:282: 全有或全无读数：本批 2 枚编辑，其中一枚不成立（空批＝整批不成立）⇒ 一个字也不落盘，文件仍是 26 字节
=== RUN   TestFSEditIsAllOrNothingAcrossABatch/second_edit_misses
    fs_edit_test.go:282: 全有或全无读数：本批 2 枚编辑，其中一枚不成立（空批＝整批不成立）⇒ 一个字也不落盘，文件仍是 26 字节
=== RUN   TestFSEditIsAllOrNothingAcrossABatch/empty_batch
    fs_edit_test.go:282: 全有或全无读数：本批 0 枚编辑，其中一枚不成立（空批＝整批不成立）⇒ 一个字也不落盘，文件仍是 26 字节
--- PASS: TestFSEditIsAllOrNothingAcrossABatch (0.03s)
    --- PASS: TestFSEditIsAllOrNothingAcrossABatch/no_op_edit_in_the_batch (0.01s)
    --- PASS: TestFSEditIsAllOrNothingAcrossABatch/overlapping_edits (0.01s)
    --- PASS: TestFSEditIsAllOrNothingAcrossABatch/second_edit_misses (0.01s)
    --- PASS: TestFSEditIsAllOrNothingAcrossABatch/empty_batch (0.01s)
=== RUN   TestFSEditAppliesABatchOnDisk
    fs_edit_test.go:318: 对照读数：26 字节 / 4 行 → fs.edit 已改写 C:\Users\swq\AppData\Local\Temp\TestFSEditAppliesABatchOnDisk2787993068\001\ok.txt：2 枚编辑全部生效，26 字节 / 4 行 → 32 字节 / 5 行（临时文件+原子重命名）
    fs_edit_test.go:319: 对照读数：AppliedSteps=["在 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditAppliesABatchOnDisk2787993068\\001 创建临时文件 .wisp-tmp-4d49082e-15948-1301852070" "向临时文件写入 32 字节" "临时文件已落盘并关闭" "原子重命名 .wisp-tmp-4d49082e-15948-1301852070 → C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSEditAppliesABatchOnDisk2787993068\\001\\ok.txt（目标此刻起为新内容）"]
--- PASS: TestFSEditAppliesABatchOnDisk (0.12s)
=== RUN   TestSweepReclaimsOnlyItsOwnStagingFiles
--- PASS: TestSweepReclaimsOnlyItsOwnStagingFiles (0.06s)
=== RUN   TestSweepNamingSchemeRejectsNearMisses
--- PASS: TestSweepNamingSchemeRejectsNearMisses (0.04s)
=== RUN   TestSweepNeverDeletesThroughARealJunction
--- PASS: TestSweepNeverDeletesThroughARealJunction (0.12s)
=== RUN   TestSweepSparesATempFileHeldByAnotherProcess
--- PASS: TestSweepSparesATempFileHeldByAnotherProcess (0.47s)
=== RUN   TestFSReadReturnsTheFileAndTaintsIt
--- PASS: TestFSReadReturnsTheFileAndTaintsIt (0.01s)
=== RUN   TestFSReadTaintFeedsR4
--- PASS: TestFSReadTaintFeedsR4 (0.01s)
=== RUN   TestFSListSummarizesADirectory
--- PASS: TestFSListSummarizesADirectory (0.01s)
=== RUN   TestFSListHonoursItsCap
--- PASS: TestFSListHonoursItsCap (0.01s)
=== RUN   TestEmptyAllowlistAuthorizesNothing
--- PASS: TestEmptyAllowlistAuthorizesNothing (0.01s)
=== RUN   TestSensitiveFileIsDeniedNotEscalated
--- PASS: TestSensitiveFileIsDeniedNotEscalated (0.01s)
=== RUN   TestFSRegistrationIsTheD34Roster
--- PASS: TestFSRegistrationIsTheD34Roster (0.00s)
=== RUN   TestDeleteEnabledAddsFSDelete
--- PASS: TestDeleteEnabledAddsFSDelete (0.00s)
=== RUN   TestToolsDirectoryShape
--- PASS: TestToolsDirectoryShape (0.00s)
=== RUN   TestD34WriteMatrix
=== RUN   TestD34WriteMatrix/read_inside_allowlist_is_L0
=== RUN   TestD34WriteMatrix/read_out_of_allowlist_is_L2_and_is_a_verdict_not_a_dormant_rule
=== RUN   TestD34WriteMatrix/list_out_of_allowlist_is_L2
=== RUN   TestD34WriteMatrix/write_new_file_is_L1
=== RUN   TestD34WriteMatrix/write_over_existing_is_L2_via_R8
=== RUN   TestD34WriteMatrix/write_out_of_allowlist_is_L2
=== RUN   TestD34WriteMatrix/trash_is_L1_because_the_bin_can_give_it_back
=== RUN   TestD34WriteMatrix/move_same_volume_new_destination_is_L1
=== RUN   TestD34WriteMatrix/move_over_existing_destination_is_L2
=== RUN   TestD34WriteMatrix/move_across_volumes_is_L2
=== RUN   TestD34WriteMatrix/delete_when_enabled_is_L2
--- PASS: TestD34WriteMatrix (0.07s)
    --- PASS: TestD34WriteMatrix/read_inside_allowlist_is_L0 (0.00s)
    --- PASS: TestD34WriteMatrix/read_out_of_allowlist_is_L2_and_is_a_verdict_not_a_dormant_rule (0.01s)
    --- PASS: TestD34WriteMatrix/list_out_of_allowlist_is_L2 (0.00s)
    --- PASS: TestD34WriteMatrix/write_new_file_is_L1 (0.01s)
    --- PASS: TestD34WriteMatrix/write_over_existing_is_L2_via_R8 (0.00s)
    --- PASS: TestD34WriteMatrix/write_out_of_allowlist_is_L2 (0.01s)
    --- PASS: TestD34WriteMatrix/trash_is_L1_because_the_bin_can_give_it_back (0.00s)
    --- PASS: TestD34WriteMatrix/move_same_volume_new_destination_is_L1 (0.01s)
    --- PASS: TestD34WriteMatrix/move_over_existing_destination_is_L2 (0.01s)
    --- PASS: TestD34WriteMatrix/move_across_volumes_is_L2 (0.00s)
    --- PASS: TestD34WriteMatrix/delete_when_enabled_is_L2 (0.00s)
=== RUN   TestCrossVolumeRuleIsR8
--- PASS: TestCrossVolumeRuleIsR8 (0.00s)
=== RUN   TestOverwriteDetectionFollowsTheCanonicalPath
--- PASS: TestOverwriteDetectionFollowsTheCanonicalPath (0.01s)
=== RUN   TestAtomicWriteKillsMidWrite
=== RUN   TestAtomicWriteKillsMidWrite/overwrite_target_keeps_its_old_bytes
=== RUN   TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all
=== RUN   TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips
--- PASS: TestAtomicWriteKillsMidWrite (0.04s)
    --- PASS: TestAtomicWriteKillsMidWrite/overwrite_target_keeps_its_old_bytes (0.01s)
    --- PASS: TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all (0.01s)
    --- PASS: TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips (0.02s)
=== RUN   TestWriteRefusesABodyOverTheCap
--- PASS: TestWriteRefusesABodyOverTheCap (0.01s)
=== RUN   TestFSTrashGoesToTheRecycleBin
--- PASS: TestFSTrashGoesToTheRecycleBin (0.16s)
=== RUN   TestFSTrashNeverFallsBackToAnUnlink
--- PASS: TestFSTrashNeverFallsBackToAnUnlink (0.01s)
=== RUN   TestDeleteIsAbsentFromTheRosterWhenTheFlagIsOff
--- PASS: TestDeleteIsAbsentFromTheRosterWhenTheFlagIsOff (0.00s)
=== RUN   TestDeleteEnabledRegistersItAsL2
=== RUN   TestDeleteEnabledRegistersItAsL2/refused_by_the_gate_it_writes_nothing
=== RUN   TestDeleteEnabledRegistersItAsL2/approved_it_deletes_and_says_it_is_permanent
--- PASS: TestDeleteEnabledRegistersItAsL2 (0.01s)
    --- PASS: TestDeleteEnabledRegistersItAsL2/refused_by_the_gate_it_writes_nothing (0.00s)
    --- PASS: TestDeleteEnabledRegistersItAsL2/approved_it_deletes_and_says_it_is_permanent (0.00s)
=== RUN   TestFSMoveSameVolume
--- PASS: TestFSMoveSameVolume (0.02s)
=== RUN   TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop
--- PASS: TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop (0.01s)
=== RUN   TestVetoAfterWorkStartedProducesAppliedSteps
--- PASS: TestVetoAfterWorkStartedProducesAppliedSteps (0.01s)
=== RUN   TestCancelBusIsNotWiredMeansNoVetoChannel
--- PASS: TestCancelBusIsNotWiredMeansNoVetoChannel (0.00s)
=== RUN   TestFSWriteSilentLossIsNotReported
time=2026-10-09T12:01:04.104+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestFSWriteSilentLossIsNotReported1674845569\001\data\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:04.112+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:04Z duration_ms=8
time=2026-10-09T12:01:04.120+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestFSWriteSilentLossIsNotReported1674845569\001\data\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:04.126+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:04Z duration_ms=5
time=2026-10-09T12:01:04.131+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:04.133+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    fswrite_silentloss_ac1_test.go:124: AC#1(a) 落盘前 672 字节 / 24 行；落盘后 450 字节 / 18 行；少了 6 行 222 字节
    fswrite_silentloss_ac1_test.go:142: AC#1(b) IsError=false Truncated=false RiskLevel=L2 ErrorClass=""
    fswrite_silentloss_ac1_test.go:144: AC#1(b) Result.Text="已写入 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSWriteSilentLossIsNotReported1674845569\\001\\notes.txt（450 字节，临时文件+原子重命名）"
    fswrite_silentloss_ac1_test.go:146: AC#1(b) AppliedSteps[0]="在 C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSWriteSilentLossIsNotReported1674845569\\001 创建临时文件 .wisp-tmp-4d49082e-15948-186863130"
    fswrite_silentloss_ac1_test.go:146: AC#1(b) AppliedSteps[1]="向临时文件写入 450 字节"
    fswrite_silentloss_ac1_test.go:146: AC#1(b) AppliedSteps[2]="临时文件已落盘并关闭"
    fswrite_silentloss_ac1_test.go:146: AC#1(b) AppliedSteps[3]="原子重命名 .wisp-tmp-4d49082e-15948-186863130 → C:\\Users\\swq\\AppData\\Local\\Temp\\TestFSWriteSilentLossIsNotReported1674845569\\001\\notes.txt（目标此刻起为新内容）"
    fswrite_silentloss_ac1_test.go:174: AC#1(c) tool_call 行：tool=fs.write risk=L2 decision=allow outcome=success error_class=""
    fswrite_silentloss_ac1_test.go:176: AC#1(c) tool_call 行的全部列名：[ID TaskID Seq Tool ArgsJSON RiskLevel Decision DecidedAt StartedAt EndedAt Outcome ErrorClass CorrelationID GrantID]
    fswrite_silentloss_ac1_test.go:193: AC#1(c) audit: "tools: call kind=success task=task-1 corr=corr-1 tool=fs.write risk=L2 decision=allow outcome=success rules_hit=[R1 R8] in_allowlist_scope=true grant_id=0 reason=\"R1: 工具声明为下界（L1）; R8: 不可逆操作（覆盖已有内容）\""
    fswrite_silentloss_ac1_test.go:193: AC#1(c) audit: "tools: PATH-ACCOUNT task=task-1 tool=fs.write roots=1 rewritten=[] unusable=[]"
    fswrite_silentloss_ac1_test.go:205: AC#1(c) fs.read 事后复读：450 字节，Truncated=false，文本里没有任何行数/大小字段（Text 就是文件内容本身）
    fswrite_silentloss_ac1_test.go:211: AC#1 读数结论：丢了 6 行 222 字节；Result.Text / AppliedSteps / tool_call 行 / audit 行 / fs.read 事后复读 五处信号全部为零 ⇒ 今天的 fs.write 对"整文件写回时漏抄"没有任何检测
--- PASS: TestFSWriteSilentLossIsNotReported (0.08s)
=== RUN   TestTicket224LiveGrantStopsTheL1Question
--- PASS: TestTicket224LiveGrantStopsTheL1Question (0.01s)
=== RUN   TestTicket224PartialPathCoverageStillAsks
--- PASS: TestTicket224PartialPathCoverageStillAsks (0.03s)
=== RUN   TestTicket224SessionGrantNeverCoversL2
--- PASS: TestTicket224SessionGrantNeverCoversL2 (0.01s)
=== RUN   TestTicket224SessionGrantNeverCoversDeny
--- PASS: TestTicket224SessionGrantNeverCoversDeny (0.01s)
=== RUN   TestTicket224GrantSourceThatCannotAnswerMeansAsking
=== RUN   TestTicket224GrantSourceThatCannotAnswerMeansAsking/reports_that_it_could_not_look
=== RUN   TestTicket224GrantSourceThatCannotAnswerMeansAsking/panics_mid-lookup
--- PASS: TestTicket224GrantSourceThatCannotAnswerMeansAsking (0.01s)
    --- PASS: TestTicket224GrantSourceThatCannotAnswerMeansAsking/reports_that_it_could_not_look (0.00s)
    --- PASS: TestTicket224GrantSourceThatCannotAnswerMeansAsking/panics_mid-lookup (0.01s)
=== RUN   TestTicket224PathlessCallIsNeverGranted
--- PASS: TestTicket224PathlessCallIsNeverGranted (0.00s)
=== RUN   TestTicket224GrantedCallBooksAllowSessionGrantWithItsRow
--- PASS: TestTicket224GrantedCallBooksAllowSessionGrantWithItsRow (0.01s)
=== RUN   TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority
--- PASS: TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority (0.00s)
=== RUN   TestPathCanonicalizerAccountsForRewrittenRoots
--- PASS: TestPathCanonicalizerAccountsForRewrittenRoots (0.01s)
=== RUN   TestTicket252P1AllowlistLegs
    paths_shortname_252_probe_test.go:229: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_probe_test.go:232: long spelling : "D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_probe_test.go:233: short spelling: "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"
    paths_shortname_252_probe_test.go:234: equal?        : false
=== RUN   TestTicket252P1AllowlistLegs/Q1_roots=LONG_asked=SHORT+missing-leaf
    paths_shortname_252_probe_test.go:252: CASE Q1 roots=LONG asked=SHORT+missing-leaf
    paths_shortname_252_probe_test.go:252:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_probe_test.go:252:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   risk.Resolve        : err=<nil> Canonical="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:252:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   NOTE: Canonicalize and risk.Resolve disagree: "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" vs "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:252:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:252:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" ok=true
    paths_shortname_252_probe_test.go:252:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:252:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:252:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252P1AllowlistLegs/Q2_roots=SHORT_asked=LONG+missing-leaf
    paths_shortname_252_probe_test.go:252: CASE Q2 roots=SHORT asked=LONG+missing-leaf
    paths_shortname_252_probe_test.go:252:   roots as configured : ["D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"]
    paths_shortname_252_probe_test.go:252:   asked raw           : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:252:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:252:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:252:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" ok=true
    paths_shortname_252_probe_test.go:252:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:252:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:252:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252P1AllowlistLegs/P1_roots=LONG_asked=LONG+missing-leaf
    paths_shortname_252_probe_test.go:252: CASE P1 roots=LONG asked=LONG+missing-leaf
    paths_shortname_252_probe_test.go:252:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_probe_test.go:252:   asked raw           : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:252:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:252:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:252:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" ok=true
    paths_shortname_252_probe_test.go:252:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:252:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:252:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252P1AllowlistLegs/P2_roots=LONG_asked=SHORT_of_an_EXISTING_dir
    paths_shortname_252_probe_test.go:252: CASE P2 roots=LONG asked=SHORT of an EXISTING dir
    paths_shortname_252_probe_test.go:252:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_probe_test.go:252:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"
    paths_shortname_252_probe_test.go:252:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp" Resolved=true Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:252:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_probe_test.go:252:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:252:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp"
    paths_shortname_252_probe_test.go:252:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:252:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp" ok=true
    paths_shortname_252_probe_test.go:252:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp"
    paths_shortname_252_probe_test.go:252:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:252:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:252:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252P1AllowlistLegs/P3_roots=SHORT_asked=SHORT_of_an_EXISTING_dir
    paths_shortname_252_probe_test.go:252: CASE P3 roots=SHORT asked=SHORT of an EXISTING dir
    paths_shortname_252_probe_test.go:252:   roots as configured : ["D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"]
    paths_shortname_252_probe_test.go:252:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"
    paths_shortname_252_probe_test.go:252:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp" Resolved=true Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:252:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_probe_test.go:252:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:252:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp"
    paths_shortname_252_probe_test.go:252:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:252:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp" ok=true
    paths_shortname_252_probe_test.go:252:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp"
    paths_shortname_252_probe_test.go:252:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:252:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:252:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252P1AllowlistLegs/Q3b_roots=SHORT_asked=SHORT+missing-leaf
    paths_shortname_252_probe_test.go:252: CASE Q3b roots=SHORT asked=SHORT+missing-leaf
    paths_shortname_252_probe_test.go:252:   roots as configured : ["D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"]
    paths_shortname_252_probe_test.go:252:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   risk.Resolve        : err=<nil> Canonical="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:252:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   NOTE: Canonicalize and risk.Resolve disagree: "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" vs "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:252:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:252:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" ok=true
    paths_shortname_252_probe_test.go:252:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
    paths_shortname_252_probe_test.go:252:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:252:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:252:   FIRST FAILING LEG              : none: both legs hold
--- PASS: TestTicket252P1AllowlistLegs (0.01s)
    --- PASS: TestTicket252P1AllowlistLegs/Q1_roots=LONG_asked=SHORT+missing-leaf (0.00s)
    --- PASS: TestTicket252P1AllowlistLegs/Q2_roots=SHORT_asked=LONG+missing-leaf (0.00s)
    --- PASS: TestTicket252P1AllowlistLegs/P1_roots=LONG_asked=LONG+missing-leaf (0.00s)
    --- PASS: TestTicket252P1AllowlistLegs/P2_roots=LONG_asked=SHORT_of_an_EXISTING_dir (0.00s)
    --- PASS: TestTicket252P1AllowlistLegs/P3_roots=SHORT_asked=SHORT_of_an_EXISTING_dir (0.00s)
    --- PASS: TestTicket252P1AllowlistLegs/Q3b_roots=SHORT_asked=SHORT+missing-leaf (0.00s)
=== RUN   TestTicket252P1OwnerMachineReachability
    paths_shortname_252_probe_test.go:260: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_probe_test.go:270: workspace dir="D:\\work\\workspace\\projects plans\\Wisp\\docs" short="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\docs"
    paths_shortname_252_probe_test.go:271: CASE W-long docs
    paths_shortname_252_probe_test.go:271:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_probe_test.go:271:   asked raw           : "D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252p1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:271:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:271:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:271:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252p1-new-note.md" ok=true
    paths_shortname_252_probe_test.go:271:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:271:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:271:   FIRST FAILING LEG              : none: both legs hold
    paths_shortname_252_probe_test.go:272: CASE W-short docs
    paths_shortname_252_probe_test.go:272:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_probe_test.go:272:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   risk.Resolve        : err=<nil> Canonical="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\docs\\q252p1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:272:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   NOTE: Canonicalize and risk.Resolve disagree: "D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252p1-new-note.md" vs "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:272:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:272:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252p1-new-note.md" ok=true
    paths_shortname_252_probe_test.go:272:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\docs\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:272:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:272:   FIRST FAILING LEG              : none: both legs hold
    paths_shortname_252_probe_test.go:270: workspace dir="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" short="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools"
    paths_shortname_252_probe_test.go:271: CASE W-long internal\tools
    paths_shortname_252_probe_test.go:271:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_probe_test.go:271:   asked raw           : "D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252p1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:271:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:271:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:271:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252p1-new-note.md" ok=true
    paths_shortname_252_probe_test.go:271:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:271:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:271:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:271:   FIRST FAILING LEG              : none: both legs hold
    paths_shortname_252_probe_test.go:272: CASE W-short internal\tools
    paths_shortname_252_probe_test.go:272:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_probe_test.go:272:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   risk.Resolve        : err=<nil> Canonical="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools\\q252p1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_probe_test.go:272:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   NOTE: Canonicalize and risk.Resolve disagree: "D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252p1-new-note.md" vs "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_probe_test.go:272:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_probe_test.go:272:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252p1-new-note.md" ok=true
    paths_shortname_252_probe_test.go:272:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\internal\\tools\\q252p1-new-note.md"
    paths_shortname_252_probe_test.go:272:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_probe_test.go:272:   InAllowlist(canonical)         : true
    paths_shortname_252_probe_test.go:272:   FIRST FAILING LEG              : none: both legs hold
--- PASS: TestTicket252P1OwnerMachineReachability (0.01s)
=== RUN   TestTicket252P1EnvSpellingCensus
    paths_shortname_252_probe_test.go:294: env USERPROFILE  : value="C:\\Users\\swq" alias="C:\\Users\\swq" -> same-as-long (no usable alias) lstat_err=<nil>
    paths_shortname_252_probe_test.go:294: env APPDATA      : value="C:\\Users\\swq\\AppData\\Roaming" alias="C:\\Users\\swq\\AppData\\Roaming" -> same-as-long (no usable alias) lstat_err=<nil>
    paths_shortname_252_probe_test.go:294: env LOCALAPPDATA : value="C:\\Users\\swq\\AppData\\Local" alias="C:\\Users\\swq\\AppData\\Local" -> same-as-long (no usable alias) lstat_err=<nil>
    paths_shortname_252_probe_test.go:294: env TEMP         : value="C:\\Users\\swq\\AppData\\Local\\Temp" alias="C:\\Users\\swq\\AppData\\Local\\Temp" -> same-as-long (no usable alias) lstat_err=<nil>
    paths_shortname_252_probe_test.go:294: env TMP          : value="C:\\Users\\swq\\AppData\\Local\\Temp" alias="C:\\Users\\swq\\AppData\\Local\\Temp" -> same-as-long (no usable alias) lstat_err=<nil>
    paths_shortname_252_probe_test.go:289: env CD           : (unset)
    paths_shortname_252_probe_test.go:294: env HOME         : value="C:\\Users\\swq" alias="C:\\Users\\swq" -> same-as-long (no usable alias) lstat_err=<nil>
    paths_shortname_252_probe_test.go:297: process cwd     : value="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" alias="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools" -> DIFFERENT spelling exists
    paths_shortname_252_probe_test.go:299: os.TempDir()    : value="C:\\Users\\swq\\AppData\\Local\\Temp" alias="C:\\Users\\swq\\AppData\\Local\\Temp" -> same-as-long (no usable alias)
    paths_shortname_252_probe_test.go:305: os.UserHomeDir(): value="C:\\Users\\swq" err=<nil> alias="C:\\Users\\swq" -> same-as-long (no usable alias)
    paths_shortname_252_probe_test.go:309: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_probe_test.go:312: Resolve("D:\\work\\workspace\\projects plans\\Wisp") -> Canonical="D:\\work\\workspace\\projects plans\\Wisp" Resolved=true Rewritten=false err=<nil>
    paths_shortname_252_probe_test.go:312: Resolve("D:\\work\\WORKSP~1\\PROJEC~1\\Wisp") -> Canonical="D:\\work\\workspace\\projects plans\\Wisp" Resolved=true Rewritten=false err=<nil>
--- PASS: TestTicket252P1EnvSpellingCensus (0.00s)
=== RUN   TestTicket252R1ShortSpellingOfNewFileIsAuthorized
    paths_shortname_252_r1_test.go:76: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_r1_test.go:79: root long : "D:\\work\\workspace\\projects plans\\Wisp"
    paths_shortname_252_r1_test.go:80: root short: "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"
=== RUN   TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A1_roots=LONG_asked=SHORT+missing-leaf
    paths_shortname_252_r1_test.go:98: CASE A1 roots=LONG asked=SHORT+missing-leaf
    paths_shortname_252_r1_test.go:98:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_r1_test.go:98:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   risk.Resolve        : err=<nil> Canonical="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_r1_test.go:98:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   NOTE: Canonicalize and risk.Resolve disagree: "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" vs "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:98:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_r1_test.go:98:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" ok=true
    paths_shortname_252_r1_test.go:98:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_r1_test.go:98:   InAllowlist(canonical)         : true
    paths_shortname_252_r1_test.go:98:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A2_roots=SHORT_asked=LONG+missing-leaf
    paths_shortname_252_r1_test.go:98: CASE A2 roots=SHORT asked=LONG+missing-leaf
    paths_shortname_252_r1_test.go:98:   roots as configured : ["D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"]
    paths_shortname_252_r1_test.go:98:   asked raw           : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_r1_test.go:98:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:98:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_r1_test.go:98:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" ok=true
    paths_shortname_252_r1_test.go:98:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_r1_test.go:98:   InAllowlist(canonical)         : true
    paths_shortname_252_r1_test.go:98:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A3_roots=SHORT_asked=SHORT+missing-leaf
    paths_shortname_252_r1_test.go:98: CASE A3 roots=SHORT asked=SHORT+missing-leaf
    paths_shortname_252_r1_test.go:98:   roots as configured : ["D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"]
    paths_shortname_252_r1_test.go:98:   asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   risk.Resolve        : err=<nil> Canonical="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_r1_test.go:98:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   NOTE: Canonicalize and risk.Resolve disagree: "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" vs "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:98:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_r1_test.go:98:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" ok=true
    paths_shortname_252_r1_test.go:98:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_r1_test.go:98:   InAllowlist(canonical)         : true
    paths_shortname_252_r1_test.go:98:   FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A4_roots=LONG_asked=LONG+missing-leaf
    paths_shortname_252_r1_test.go:98: CASE A4 roots=LONG asked=LONG+missing-leaf
    paths_shortname_252_r1_test.go:98:   roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
    paths_shortname_252_r1_test.go:98:   asked raw           : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   risk.Resolve        : err=<nil> Canonical="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" Resolved=false Rewritten=false Rewrites=[]
    paths_shortname_252_r1_test.go:98:   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:98:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_r1_test.go:98:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" ok=true
    paths_shortname_252_r1_test.go:98:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
    paths_shortname_252_r1_test.go:98:   LEG2 rootsContain(roots,foldRF): true
    paths_shortname_252_r1_test.go:98:   InAllowlist(canonical)         : true
    paths_shortname_252_r1_test.go:98:   FIRST FAILING LEG              : none: both legs hold
--- PASS: TestTicket252R1ShortSpellingOfNewFileIsAuthorized (0.01s)
    --- PASS: TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A1_roots=LONG_asked=SHORT+missing-leaf (0.00s)
    --- PASS: TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A2_roots=SHORT_asked=LONG+missing-leaf (0.00s)
    --- PASS: TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A3_roots=SHORT_asked=SHORT+missing-leaf (0.00s)
    --- PASS: TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A4_roots=LONG_asked=LONG+missing-leaf (0.00s)
=== RUN   TestTicket252R1BothSpellingsAnswerTheSame
    paths_shortname_252_r1_test.go:107: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
=== RUN   TestTicket252R1BothSpellingsAnswerTheSame/Wisp
    paths_shortname_252_r1_test.go:116: tree long="D:\\work\\workspace\\projects plans\\Wisp" short="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"
    paths_shortname_252_r1_test.go:121: roots=[D:\work\workspace\projects plans\Wisp] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\workspace\projects plans\Wisp] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\workspace\projects plans\Wisp] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\workspace\projects plans\Wisp] legs(SHORT): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:121: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(SHORT): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:121: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(SHORT): leg1=true leg2=true final=true
=== RUN   TestTicket252R1BothSpellingsAnswerTheSame/docs
    paths_shortname_252_r1_test.go:116: tree long="D:\\work\\workspace\\projects plans\\Wisp\\docs" short="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\docs"
    paths_shortname_252_r1_test.go:121: roots=[D:\work\workspace\projects plans\Wisp] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\workspace\projects plans\Wisp] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\workspace\projects plans\Wisp] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\workspace\projects plans\Wisp] legs(SHORT): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:121: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\docs] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\docs] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\docs] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\docs] legs(SHORT): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:121: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\docs\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(SHORT): leg1=true leg2=true final=true
=== RUN   TestTicket252R1BothSpellingsAnswerTheSame/tools
    paths_shortname_252_r1_test.go:116: tree long="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" short="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools"
    paths_shortname_252_r1_test.go:121: roots=[D:\work\workspace\projects plans\Wisp] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\workspace\projects plans\Wisp] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\workspace\projects plans\Wisp] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\workspace\projects plans\Wisp] legs(SHORT): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:121: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\internal\tools] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\internal\tools] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\internal\tools] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp\internal\tools] legs(SHORT): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:121: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(LONG)="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:122: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] canonical(SHORT)="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools\\q252r1-round-trip.md"
    paths_shortname_252_r1_test.go:123: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(LONG): leg1=true leg2=true final=true
    paths_shortname_252_r1_test.go:124: roots=[D:\work\WORKSP~1\PROJEC~1\Wisp] legs(SHORT): leg1=true leg2=true final=true
--- PASS: TestTicket252R1BothSpellingsAnswerTheSame (0.02s)
    --- PASS: TestTicket252R1BothSpellingsAnswerTheSame/Wisp (0.01s)
    --- PASS: TestTicket252R1BothSpellingsAnswerTheSame/docs (0.01s)
    --- PASS: TestTicket252R1BothSpellingsAnswerTheSame/tools (0.01s)
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization
    paths_shortname_252_r1_test.go:163: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization/N1_sibling_tree,_long
    paths_shortname_252_r1_test.go:210:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:210:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\q252r1-no-such-tree\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG1 rootsContain(roots,folded): false
    paths_shortname_252_r1_test.go:210:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\q252r1-no-such-tree\\q252r1-out-of-scope.md" ok=true
    paths_shortname_252_r1_test.go:210:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\q252r1-no-such-tree\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG2 rootsContain(roots,foldRF): false
    paths_shortname_252_r1_test.go:210:   InAllowlist(canonical)         : false
    paths_shortname_252_r1_test.go:210:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization/N2_sibling_tree,_short_spelling_of_the_parent
    paths_shortname_252_r1_test.go:210:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:210:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\q252r1-no-such-tree\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG1 rootsContain(roots,folded): false
    paths_shortname_252_r1_test.go:210:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\q252r1-no-such-tree\\q252r1-out-of-scope.md" ok=true
    paths_shortname_252_r1_test.go:210:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\q252r1-no-such-tree\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG2 rootsContain(roots,foldRF): false
    paths_shortname_252_r1_test.go:210:   InAllowlist(canonical)         : false
    paths_shortname_252_r1_test.go:210:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization/N3_the_root's_own_parent
    paths_shortname_252_r1_test.go:210:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:210:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG1 rootsContain(roots,folded): false
    paths_shortname_252_r1_test.go:210:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\q252r1-out-of-scope.md" ok=true
    paths_shortname_252_r1_test.go:210:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG2 rootsContain(roots,foldRF): false
    paths_shortname_252_r1_test.go:210:   InAllowlist(canonical)         : false
    paths_shortname_252_r1_test.go:210:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization/N4_the_root's_parent_spelled_short
    paths_shortname_252_r1_test.go:210:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:210:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG1 rootsContain(roots,folded): false
    paths_shortname_252_r1_test.go:210:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\q252r1-out-of-scope.md" ok=true
    paths_shortname_252_r1_test.go:210:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG2 rootsContain(roots,foldRF): false
    paths_shortname_252_r1_test.go:210:   InAllowlist(canonical)         : false
    paths_shortname_252_r1_test.go:210:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization/N5_one_character_off_the_root's_tail
    paths_shortname_252_r1_test.go:210:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:210:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp-evil\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG1 rootsContain(roots,folded): false
    paths_shortname_252_r1_test.go:210:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp-evil\\q252r1-out-of-scope.md" ok=true
    paths_shortname_252_r1_test.go:210:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp-evil\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG2 rootsContain(roots,foldRF): false
    paths_shortname_252_r1_test.go:210:   InAllowlist(canonical)         : false
    paths_shortname_252_r1_test.go:210:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization/N6_a_name_on_a_volume_that_is_not_there
    paths_shortname_252_r1_test.go:210:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:210:   foldPath(canonical)           : "z:\\q252r1-dead-volume\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG1 rootsContain(roots,folded): false
    paths_shortname_252_r1_test.go:210:   resolvedForm(canonical)        : "Z:\\q252r1-dead-volume\\q252r1-out-of-scope.md" ok=true
    paths_shortname_252_r1_test.go:210:   foldPath(resolvedForm)         : "z:\\q252r1-dead-volume\\q252r1-out-of-scope.md"
    paths_shortname_252_r1_test.go:210:   LEG2 rootsContain(roots,foldRF): false
    paths_shortname_252_r1_test.go:210:   InAllowlist(canonical)         : false
    paths_shortname_252_r1_test.go:210:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
=== RUN   TestTicket252R1AlignmentAddsNoAuthorization/N7_inside_the_root_but_unresolvable
    paths_shortname_252_r1_test.go:226:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_shortname_252_r1_test.go:226:   foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252r1<invalid>:*?.md"
    paths_shortname_252_r1_test.go:226:   LEG1 rootsContain(roots,folded): true
    paths_shortname_252_r1_test.go:226:   resolvedForm(canonical)        : "" ok=false
    paths_shortname_252_r1_test.go:226:   foldPath(resolvedForm)         : ""
    paths_shortname_252_r1_test.go:226:   LEG2 rootsContain(roots,foldRF): false
    paths_shortname_252_r1_test.go:226:   InAllowlist(canonical)         : false
    paths_shortname_252_r1_test.go:226:   FIRST FAILING LEG              : LEG 2 resolvedForm (paths.go:152-153)
=== NAME  TestTicket252R1AlignmentAddsNoAuthorization
    paths_shortname_252_r1_test.go:244: short ask roots check
--- PASS: TestTicket252R1AlignmentAddsNoAuthorization (0.01s)
    --- PASS: TestTicket252R1AlignmentAddsNoAuthorization/N1_sibling_tree,_long (0.00s)
    --- PASS: TestTicket252R1AlignmentAddsNoAuthorization/N2_sibling_tree,_short_spelling_of_the_parent (0.00s)
    --- PASS: TestTicket252R1AlignmentAddsNoAuthorization/N3_the_root's_own_parent (0.00s)
    --- PASS: TestTicket252R1AlignmentAddsNoAuthorization/N4_the_root's_parent_spelled_short (0.00s)
    --- PASS: TestTicket252R1AlignmentAddsNoAuthorization/N5_one_character_off_the_root's_tail (0.00s)
    --- PASS: TestTicket252R1AlignmentAddsNoAuthorization/N6_a_name_on_a_volume_that_is_not_there (0.00s)
    --- PASS: TestTicket252R1AlignmentAddsNoAuthorization/N7_inside_the_root_but_unresolvable (0.00s)
=== RUN   TestTicket252R1CanonicalStillNamesOneTree
    paths_shortname_252_r1_test.go:262: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
--- PASS: TestTicket252R1CanonicalStillNamesOneTree (0.00s)
=== RUN   TestTicket107AllowlistJudgmentTwoShapes
    paths_ticket107_portable_test.go:57: shape R root="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107AllowlistJudgmentTwoShapes1289076499\\001\\proj" target="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107AllowlistJudgmentTwoShapes1289076499\\001\\proj\\a.txt" roots=[c:\users\swq\appdata\local\temp\testticket107allowlistjudgmenttwoshapes1289076499\001\proj] InAllowlist=true
    paths_ticket107_portable_test.go:58: shape E root="%WISP107_ROOT%\\proj" target="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107AllowlistJudgmentTwoShapes1289076499\\001\\proj\\a.txt" roots=[c:\users\swq\appdata\local\temp\testticket107allowlistjudgmenttwoshapes1289076499\001\proj] rewritten=["%WISP107_ROOT%\\proj" -> C:\Users\swq\AppData\Local\Temp\TestTicket107AllowlistJudgmentTwoShapes1289076499\001\proj (env)] unusable=[] InAllowlist=true
    paths_ticket107_portable_test.go:60: cross feed (root from shape R, canonical from shape E) InAllowlist=true
--- PASS: TestTicket107AllowlistJudgmentTwoShapes (0.01s)
=== RUN   TestTicket107AllowlistBoundaryIsComponentWise
--- PASS: TestTicket107AllowlistBoundaryIsComponentWise (0.01s)
=== RUN   TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing
    paths_ticket107b_probes_test.go:108: mklink /J C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\001\proj -> C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\001\outside: Junction created for C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\001\proj <<===>> C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\001\outside
    paths_ticket107b_probes_test.go:109: premise signals for "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\\001\\proj" -> "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\\001\\outside": lstat(err=<nil>, modeLink=false) readlink("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\\001\\outside", err=<nil>) eval(link)="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\\001\\proj",<nil> eval(through)="",The system cannot find the path specified.
    paths_ticket107b_probes_test.go:120: Canonicalize("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\\001\\proj\\marker.txt") refused: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions (judging the lexical spelling instead)
    paths_ticket107b_probes_test.go:121: roots=[] rewritten=[] unusable=["%WISP107B_A_ROOT%\\proj": risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions] judged="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\\001\\proj\\marker.txt" opened="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing913458824\\001\\proj\\marker.txt"
--- PASS: TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing (0.06s)
=== RUN   TestTicket107bProbeCLinkInsideAllowedRootStaysOutside
    paths_ticket107b_probes_test.go:158: mklink /J C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\001\proj\esc -> C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\001\outside: Junction created for C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\001\proj\esc <<===>> C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\001\outside
    paths_ticket107b_probes_test.go:159: premise signals for "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\\001\\proj\\esc" -> "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\\001\\outside": lstat(err=<nil>, modeLink=false) readlink("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\\001\\outside", err=<nil>) eval(link)="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\\001\\proj\\esc",<nil> eval(through)="",The system cannot find the path specified.
    paths_ticket107b_probes_test.go:165: roots=[c:\users\swq\appdata\local\temp\testticket107bprobeclinkinsideallowedrootstaysoutside3991050933\001\proj] rewritten=["%WISP107B_C_ROOT%\\proj" -> C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\001\proj (env)] unusable=[]
    paths_ticket107b_probes_test.go:178: Canonicalize("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeCLinkInsideAllowedRootStaysOutside3991050933\\001\\proj\\esc\\marker.txt") refused: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions (judging the lexical spelling instead)
--- PASS: TestTicket107bProbeCLinkInsideAllowedRootStaysOutside (0.06s)
=== RUN   TestTicket107bProbeBJudgedTreeIsTheOpenedTree
    paths_ticket107b_probes_test.go:201: mklink /J C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\001\proj -> C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\001\outside: Junction created for C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\001\proj <<===>> C:\Users\swq\AppData\Local\Temp\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\001\outside
    paths_ticket107b_probes_test.go:202: premise signals for "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\proj" -> "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\outside": lstat(err=<nil>, modeLink=false) readlink("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\outside", err=<nil>) eval(link)="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\proj",<nil> eval(through)="",The system cannot find the path specified.
    paths_ticket107b_probes_test.go:209: Canonicalize("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\proj\\marker.txt") refused: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions (judging the lexical spelling instead)
    paths_ticket107b_probes_test.go:211: unexpanded link root: roots=[] unusable=["C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\proj": risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions] judged="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\proj\\marker.txt" opened="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket107bProbeBJudgedTreeIsTheOpenedTree2448249515\\001\\proj\\marker.txt" InAllowlist=false
--- PASS: TestTicket107bProbeBJudgedTreeIsTheOpenedTree (0.05s)
=== RUN   TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses
    paths_twocontainments_252_r2_test.go:105: mklink /J C:\Users\swq\AppData\Local\Temp\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\001\proj\esc -> C:\Users\swq\AppData\Local\Temp\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\001\outside: Junction created for C:\Users\swq\AppData\Local\Temp\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\001\proj\esc <<===>> C:\Users\swq\AppData\Local\Temp\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\001\outside
    paths_twocontainments_252_r2_test.go:106: premise signals for "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\proj\\esc" -> "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\outside": lstat(err=<nil>, modeLink=false) readlink("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\outside", err=<nil>) eval(link)="C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\proj\\esc",<nil> eval(through)="",The system cannot find the path specified.
    paths_twocontainments_252_r2_test.go:139: allowed root        : "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\proj" -> roots=["c:\\users\\swq\\appdata\\local\\temp\\testticket252r2resolvedlegrefuseswhatonlyitrefuses2257602281\\001\\proj"]
    paths_twocontainments_252_r2_test.go:140: escaping ask        : "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\proj\\esc\\marker.txt"
    paths_twocontainments_252_r2_test.go:141: LEG1 lexical        : true
    paths_twocontainments_252_r2_test.go:142: resolvedForm        : "" ok=false
    paths_twocontainments_252_r2_test.go:143: LEG2 resolved       : false (refused because ok=false, or because the answer is outside the root)
    paths_twocontainments_252_r2_test.go:144: InAllowlist         : false
    paths_twocontainments_252_r2_test.go:145: link target written : "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\outside"
    paths_twocontainments_252_r2_test.go:146: what the link gives : "C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses2257602281\\001\\proj\\esc\\marker.txt"
--- PASS: TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses (0.06s)
=== RUN   TestTicket252R2BothContainmentsReachable
    paths_twocontainments_252_r2_test.go:283: S1 scanned 23 non-test .go files of this package; InAllowlist is declared at paths.go:201
    paths_twocontainments_252_r2_test.go:344: S1 carrier: receiver="p" parameter="canonical" decl=paths.go:201
    paths_twocontainments_252_r2_test.go:347: containment paths.go:208 key=f calls=[foldPath] tainted=false refuses=true
    paths_twocontainments_252_r2_test.go:347: containment paths.go:221 key=foldPath(rf) calls=[foldPath resolvedForm] tainted=true refuses=true
    paths_twocontainments_252_r2_test.go:350: S1 re-resolution names: [ok rf]
    paths_twocontainments_252_r2_test.go:351: S1 names an if-with-return-false consults: [f foldPath ok p rf roots rootsContain string workspace]
--- PASS: TestTicket252R2BothContainmentsReachable (0.01s)
=== RUN   TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses
    paths_twocontainments_252_r2_windows_test.go:44: cwd="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" derived repo root="D:\\work\\workspace\\projects plans\\Wisp"
=== RUN   TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses/B1_short-of-existing-tree_+_missing-leaf
    paths_twocontainments_252_r2_windows_test.go:76:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_twocontainments_252_r2_windows_test.go:76:   foldPath(canonical)           : "d:\\work\\worksp~1\\projec~1\\wisp\\q252r2-lexical-leg-only.md"
    paths_twocontainments_252_r2_windows_test.go:76:   LEG1 rootsContain(roots,folded): false
    paths_twocontainments_252_r2_windows_test.go:76:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r2-lexical-leg-only.md" ok=true
    paths_twocontainments_252_r2_windows_test.go:76:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252r2-lexical-leg-only.md"
    paths_twocontainments_252_r2_windows_test.go:76:   LEG2 rootsContain(roots,foldRF): true
    paths_twocontainments_252_r2_windows_test.go:76:   InAllowlist(canonical)         : false
    paths_twocontainments_252_r2_windows_test.go:76:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
    paths_twocontainments_252_r2_windows_test.go:107:   aligned Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_twocontainments_252_r2_windows_test.go:107:   aligned foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252r2-lexical-leg-only.md"
    paths_twocontainments_252_r2_windows_test.go:107:   aligned LEG1 rootsContain(roots,folded): true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r2-lexical-leg-only.md" ok=true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252r2-lexical-leg-only.md"
    paths_twocontainments_252_r2_windows_test.go:107:   aligned LEG2 rootsContain(roots,foldRF): true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned InAllowlist(canonical)         : true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned FIRST FAILING LEG              : none: both legs hold
=== RUN   TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses/B2_short-of-existing-tree_itself
    paths_twocontainments_252_r2_windows_test.go:76:   Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_twocontainments_252_r2_windows_test.go:76:   foldPath(canonical)           : "d:\\work\\worksp~1\\projec~1\\wisp"
    paths_twocontainments_252_r2_windows_test.go:76:   LEG1 rootsContain(roots,folded): false
    paths_twocontainments_252_r2_windows_test.go:76:   resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp" ok=true
    paths_twocontainments_252_r2_windows_test.go:76:   foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp"
    paths_twocontainments_252_r2_windows_test.go:76:   LEG2 rootsContain(roots,foldRF): true
    paths_twocontainments_252_r2_windows_test.go:76:   InAllowlist(canonical)         : false
    paths_twocontainments_252_r2_windows_test.go:76:   FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
    paths_twocontainments_252_r2_windows_test.go:107:   aligned Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
    paths_twocontainments_252_r2_windows_test.go:107:   aligned foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp"
    paths_twocontainments_252_r2_windows_test.go:107:   aligned LEG1 rootsContain(roots,folded): true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp" ok=true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp"
    paths_twocontainments_252_r2_windows_test.go:107:   aligned LEG2 rootsContain(roots,foldRF): true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned InAllowlist(canonical)         : true
    paths_twocontainments_252_r2_windows_test.go:107:   aligned FIRST FAILING LEG              : none: both legs hold
--- PASS: TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses (0.00s)
    --- PASS: TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses/B1_short-of-existing-tree_+_missing-leaf (0.00s)
    --- PASS: TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses/B2_short-of-existing-tree_itself (0.00s)
=== RUN   TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing
--- PASS: TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing (0.01s)
=== RUN   TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26
=== RUN   TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26/the_refusal_leg_is_untouched
=== RUN   TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26/a_recorded_account_comes_back_as_written
--- PASS: TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26 (0.00s)
    --- PASS: TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26/the_refusal_leg_is_untouched (0.00s)
    --- PASS: TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26/a_recorded_account_comes_back_as_written (0.00s)
=== RUN   TestSetWorkspaceRootRefusesAnAccountAboutAnotherTree
--- PASS: TestSetWorkspaceRootRefusesAnAccountAboutAnotherTree (0.01s)
=== RUN   TestClearWorkspaceDropsTheAccountWithTheRoot
--- PASS: TestClearWorkspaceDropsTheAccountWithTheRoot (0.00s)
=== RUN   TestWorkspaceSwitchNarrowsWhatTheAssessorJudges
    paths_workspace_test.go:86: AC#3(i): C:\Users\swq\AppData\Local\Temp\TestWorkspaceSwitchNarrowsWhatTheAssessorJudges3020217768\001\beta write L1 -> L2 after switching to C:\Users\swq\AppData\Local\Temp\TestWorkspaceSwitchNarrowsWhatTheAssessorJudges3020217768\001\alpha; inside the workspace still L1
--- PASS: TestWorkspaceSwitchNarrowsWhatTheAssessorJudges (0.01s)
=== RUN   TestWorkspaceSwitchRefusesOutOfScopeAndMissingPaths
--- PASS: TestWorkspaceSwitchRefusesOutOfScopeAndMissingPaths (0.01s)
=== RUN   TestWorkspaceSwitchRefusesAnExpandedSpelling
=== RUN   TestWorkspaceSwitchRefusesAnExpandedSpelling/leading_tilde
    paths_workspace_test.go:155: expanded spelling refused as required: risk: expansion moved this path onto a tree the caller did not name: ~ expands to C:\Users\swq (home); act on the expanded tree only if you asked for it by name
=== RUN   TestWorkspaceSwitchRefusesAnExpandedSpelling/environment_variable
    paths_workspace_test.go:172: expanded spelling refused as required: risk: expansion moved this path onto a tree the caller did not name: %TEMP% expands to C:\Users\swq\AppData\Local\Temp (env); act on the expanded tree only if you asked for it by name
--- PASS: TestWorkspaceSwitchRefusesAnExpandedSpelling (0.00s)
    --- PASS: TestWorkspaceSwitchRefusesAnExpandedSpelling/leading_tilde (0.00s)
    --- PASS: TestWorkspaceSwitchRefusesAnExpandedSpelling/environment_variable (0.00s)
=== RUN   TestWorkspaceSwitchRefusesAJunctionToOutside
    paths_workspace_test.go:220: junction workspace refused with C26's reason: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions
--- PASS: TestWorkspaceSwitchRefusesAJunctionToOutside (0.05s)
=== RUN   TestCanonicalizeReturnsAPathTheOSCanOpen
--- PASS: TestCanonicalizeReturnsAPathTheOSCanOpen (0.00s)
=== RUN   TestCanonicalizeAgreesWithTheOSName
--- PASS: TestCanonicalizeAgreesWithTheOSName (0.00s)
=== RUN   TestFoldPathKeepsPosixBackslashesDistinct
--- PASS: TestFoldPathKeepsPosixBackslashesDistinct (0.00s)
=== RUN   TestDirOfBaseOfCutOnThePlatformSeparator
--- PASS: TestDirOfBaseOfCutOnThePlatformSeparator (0.00s)
=== RUN   TestPointer183BackfilledArtifactRereadStaysClean
--- PASS: TestPointer183BackfilledArtifactRereadStaysClean (0.02s)
=== RUN   TestPointer183SiblingArtifactInTheSameDirectoryStillHits
--- PASS: TestPointer183SiblingArtifactInTheSameDirectoryStillHits (0.02s)
=== RUN   TestPointer185SecondRereadOfTheSameArtifactStaysClean
--- PASS: TestPointer185SecondRereadOfTheSameArtifactStaysClean (0.02s)
=== RUN   TestPointer185PagedRereadsKeepWorking
    pointer_185_cli_seam_test.go:242: 读数：三发分页续读后 scope=185r1-paged 的 mark 枚数=4
--- PASS: TestPointer185PagedRereadsKeepWorking (0.02s)
=== RUN   TestPointer185ForeignFileCarryingTheHostPathStillBlocksTheReread
--- PASS: TestPointer185ForeignFileCarryingTheHostPathStillBlocksTheReread (0.02s)
=== RUN   TestPointer185AModelNominatedFileIsNotRostered
--- PASS: TestPointer185AModelNominatedFileIsNotRostered (0.03s)
=== RUN   TestPointer185AnotherTaskCannotBorrowTheRoster
--- PASS: TestPointer185AnotherTaskCannotBorrowTheRoster (0.03s)
=== RUN   TestPointer185SiblingHostArtifactIsNotLaundered
--- PASS: TestPointer185SiblingHostArtifactIsNotLaundered (0.03s)
=== RUN   TestRecycleBinStructLayoutMatchesWin32
--- PASS: TestRecycleBinStructLayoutMatchesWin32 (0.00s)
=== RUN   TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate
--- PASS: TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate (0.07s)
=== RUN   TestShellTrashRefusesAMissingPathWithoutDeletingAnything
--- PASS: TestShellTrashRefusesAMissingPathWithoutDeletingAnything (0.00s)
=== RUN   Test197SpawnPublishesIdentityRow
time=2026-10-09T12:01:04.960+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-2e04c4e5-f420-4e1e-83e4-67d3a3fb663a owner=tools
--- PASS: Test197SpawnPublishesIdentityRow (0.00s)
=== RUN   Test197RowExistsBeforeFirstChildModelCall
time=2026-10-09T12:01:04.961+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-12ccc3c1-f122-4b06-b9b8-b017b17b941b owner=tools
--- PASS: Test197RowExistsBeforeFirstChildModelCall (0.00s)
=== RUN   Test197SubagentPoolNeverExceedsBridgeCeiling
--- PASS: Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)
=== RUN   Test197SpawnDescriptionNamesTheRealPoolCap
--- PASS: Test197SpawnDescriptionNamesTheRealPoolCap (0.00s)
=== RUN   Test197SubagentPoolCapsAtBridgeCeiling
time=2026-10-09T12:01:04.961+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-965227e5-c900-429f-91dc-9b734c53857e owner=tools
time=2026-10-09T12:01:04.961+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-4e80b47c-3976-4877-a030-12418a267716 owner=tools
time=2026-10-09T12:01:04.962+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-3d75e458-56bf-4aa0-b889-fa187babc0fe owner=tools
time=2026-10-09T12:01:04.962+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-74454224-b08e-48cb-8b5d-8449303e593f owner=tools
--- PASS: Test197SubagentPoolCapsAtBridgeCeiling (0.00s)
=== RUN   Test197FullPoolRefusesNextSpawnWithReadableReason
time=2026-10-09T12:01:04.962+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-16a898fb-d95d-4dc3-8033-98eba71a3f67 owner=tools
time=2026-10-09T12:01:04.962+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-fd370957-02a4-4a54-beb0-aaff94afcab5 owner=tools
time=2026-10-09T12:01:04.962+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-84cee51d-1f98-46f1-a8a3-8bca27005aec owner=tools
time=2026-10-09T12:01:04.963+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-b8076de2-d04b-4356-8295-5b1747c5b905 owner=tools
--- PASS: Test197FullPoolRefusesNextSpawnWithReadableReason (0.00s)
=== RUN   Test197ChildCannotDeriveSubagent
time=2026-10-09T12:01:04.963+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-3906749e-4b73-4289-8a52-3d4c99770fd4 owner=tools
--- PASS: Test197ChildCannotDeriveSubagent (0.00s)
=== RUN   Test197ConclusionCarriesTaskOutputTaint
time=2026-10-09T12:01:04.963+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c62fe2ac-49de-49ec-b4d7-0426d0f8d77a owner=tools
time=2026-10-09T12:01:04.963+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-7edf08f8-1025-46a3-9bbc-8b029482dc08 owner=tools
--- PASS: Test197ConclusionCarriesTaskOutputTaint (0.00s)
=== RUN   Test197CancelIsPerRowAndNeverCascades
time=2026-10-09T12:01:04.964+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-98eb7ead-56a4-4dc2-a3e1-322ea75812d7 owner=tools
--- PASS: Test197CancelIsPerRowAndNeverCascades (0.00s)
=== RUN   Test197SubagentHasNoSelfApprovalOutlet
time=2026-10-09T12:01:04.965+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-9e5363d0-eaf5-4582-b698-132da53b8720 owner=tools
--- PASS: Test197SubagentHasNoSelfApprovalOutlet (0.00s)
=== RUN   Test197StreamKeyShapeIsLiteral
--- PASS: Test197StreamKeyShapeIsLiteral (0.00s)
=== RUN   Test222SpawnConclusionArrivesThroughRealBridgeChildren
time=2026-10-09T12:01:04.965+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-9fd2a270-7cc0-4239-9c22-868440858fdd owner=tools
time=2026-10-09T12:01:04.965+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-0241023c-3b4d-481d-9e8c-e76a5cf76128 owner=tools
time=2026-10-09T12:01:04.965+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-cb875de3-ba04-4166-a63d-6c1bfc4592e7 owner=tools
time=2026-10-09T12:01:04.965+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-828df8ea-bbbb-489b-b707-f06b06546002 owner=tools
--- PASS: Test222SpawnConclusionArrivesThroughRealBridgeChildren (0.00s)
=== RUN   Test222WaitingParentHoldsNoBridgeSlot
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c4e01f82-d023-4724-85b3-5f58a2804934 owner=tools
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-82141e84-b87b-4885-b2ef-f33f97d4a576 owner=tools
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-d100deda-0cf8-4533-a3bb-056d0be429ce owner=tools
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e3d4a8f1-c3a3-4d21-83ee-5dec53d84535 owner=tools
--- PASS: Test222WaitingParentHoldsNoBridgeSlot (0.00s)
=== RUN   Test222CeilingStillCapsExecutedCallsWhileParentsWait
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-69d8c936-e20f-4c46-ae8e-833e49f41e90 owner=tools
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-abd0235c-890f-4fc7-874f-13094775f2d1 owner=tools
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ab17ecfe-af40-48b8-850f-11448f53dd41 owner=tools
time=2026-10-09T12:01:04.966+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f6ffc5b8-b772-40c4-8b9e-10d6de452abc owner=tools
--- PASS: Test222CeilingStillCapsExecutedCallsWhileParentsWait (0.00s)
=== RUN   Test221TaskCancelIsRegisteredAtItsFrozenLevel
time=2026-10-09T12:01:04.967+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-dc2cb7df-4826-4294-ae62-8b5e8686c95b owner=tools
--- PASS: Test221TaskCancelIsRegisteredAtItsFrozenLevel (0.00s)
=== RUN   Test221DeferredMarkerForCancelLiftedButListStillMarked
--- PASS: Test221DeferredMarkerForCancelLiftedButListStillMarked (0.00s)
=== RUN   Test221SpawnDescriptionPromisesOnlyWhatIsTrue
--- PASS: Test221SpawnDescriptionPromisesOnlyWhatIsTrue (0.00s)
=== RUN   Test221ParentStopsItsOwnChildRowAndStreamSettle
time=2026-10-09T12:01:04.997+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test221ParentStopsItsOwnChildRowAndStreamSettle2115422086\001\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:05.008+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:04Z duration_ms=10
time=2026-10-09T12:01:05.016+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test221ParentStopsItsOwnChildRowAndStreamSettle2115422086\001\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:05.022+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:05Z duration_ms=6
time=2026-10-09T12:01:05.028+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:05.030+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T12:01:05.033+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-598c24b6-be5d-489f-86f5-5faab0281d06 owner=tools
--- PASS: Test221ParentStopsItsOwnChildRowAndStreamSettle (0.08s)
=== RUN   Test221SubagentCannotStopSiblingOrItself
time=2026-10-09T12:01:05.046+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-00de2146-0de2-49a3-8618-def032fb24ce owner=tools
time=2026-10-09T12:01:05.046+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-e88e162d-c3e0-4d72-b442-35453c2dd41a owner=tools
--- PASS: Test221SubagentCannotStopSiblingOrItself (0.00s)
=== RUN   Test221ParentCancellationStillDoesNotCascade
time=2026-10-09T12:01:05.047+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-dacc9347-0c8d-4ba7-a9c8-e0b520ec695d owner=tools
--- PASS: Test221ParentCancellationStillDoesNotCascade (0.00s)
=== RUN   Test221TaskCancelUnderNoGateStopsAtTheWindow
--- PASS: Test221TaskCancelUnderNoGateStopsAtTheWindow (0.00s)
=== RUN   Test221EveryPromisedTaskNameIsRegistered
--- PASS: Test221EveryPromisedTaskNameIsRegistered (0.00s)
=== RUN   TestTaskOutputAC2BeforeLegIsUnreachable
    task_output_ac2_before_test.go:52: AC#2 BEFORE verbatim: IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.00s)
=== RUN   TestCanonicalizeFailureFailsClosedInReply174
    task_output_canonicalize_fail_174_test.go:73: canonicalize-error verbatim: IsError=false Truncated=true announce="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，C26 连规范化都没通过（tools: empty path）；注意：这条路径现在读不到，宿主登记的那份副本文件并不存在，全文见    …]"
--- PASS: TestCanonicalizeFailureFailsClosedInReply174 (0.00s)
=== RUN   TestCanonicalizeFailureIsNotTheUnwiredArm174
    task_output_canonicalize_fail_174_test.go:144: unwired control verbatim: announce="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理），全文见    …]"
--- PASS: TestCanonicalizeFailureIsNotTheUnwiredArm174 (0.00s)
=== RUN   TestTaskOutputAC2AfterLegIsReachable
    task_output_leg_test.go:70: AC#2 AFTER verbatim: IsError=false ErrorClass="" Text="后台任务打印的第一行"
--- PASS: TestTaskOutputAC2AfterLegIsReachable (0.00s)
=== RUN   TestUnknownTaskIDIsLoudNotEmpty
--- PASS: TestUnknownTaskIDIsLoudNotEmpty (0.00s)
=== RUN   TestMissingArgsAndUnwiredRoster
--- PASS: TestMissingArgsAndUnwiredRoster (0.00s)
=== RUN   TestLongOutputPointerRecoversEveryByte
--- PASS: TestLongOutputPointerRecoversEveryByte (0.01s)
=== RUN   TestTruncationShapeIsTheD15Triple
--- PASS: TestTruncationShapeIsTheD15Triple (0.00s)
=== RUN   TestTruncationBudgetsAreRead
--- PASS: TestTruncationBudgetsAreRead (0.00s)
=== RUN   TestNoCopyFileAnnouncesLoss
--- PASS: TestNoCopyFileAnnouncesLoss (0.00s)
=== RUN   TestPointerPast256KiBIsNotFullyReadable
--- PASS: TestPointerPast256KiBIsNotFullyReadable (0.09s)
=== RUN   TestRosterIsNotAPathSurface
--- PASS: TestRosterIsNotAPathSurface (0.00s)
=== RUN   TestPointerOutsideAuthorizedRootSpeaks
    task_output_pointer_notice_test.go:118: shape (a) verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerOutsideAuthorizedRootSpeaks3737480910\\001\\tool-output-174-outside.txt…]"
--- PASS: TestPointerOutsideAuthorizedRootSpeaks (0.01s)
=== RUN   TestPointerToMissingCopyFileSpeaks
    task_output_pointer_notice_test.go:162: shape (b) missing file verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，宿主登记的那份副本文件并不存在，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerToMissingCopyFileSpeaks3753537920\\001\\tool-output-does-not-exist.txt…]"
--- PASS: TestPointerToMissingCopyFileSpeaks (0.01s)
=== RUN   TestPointerToNonRegularPathSpeaks
    task_output_pointer_notice_test.go:187: shape (b) directory verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状），全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerToNonRegularPathSpeaks3041644503\\001…]"
--- PASS: TestPointerToNonRegularPathSpeaks (0.01s)
=== RUN   TestHealthyPointerStaysSilent
    task_output_pointer_notice_test.go:214: silent control verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestHealthyPointerStaysSilent1918760832\\001\\tool-output-174-healthy.txt…]"
--- PASS: TestHealthyPointerStaysSilent (0.02s)
=== RUN   TestUnwiredJudgeFailsClosedInReply
    task_output_pointer_notice_test.go:253: fail-closed verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理），全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestUnwiredJudgeFailsClosedInReply1663751365\\001\\tool-output-174-nojudge.txt…]"
--- PASS: TestUnwiredJudgeFailsClosedInReply (0.01s)
=== RUN   TestBothPointerDefectsAreReportedSeparately
    task_output_pointer_notice_test.go:283: both shapes verbatim: IsError=false Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来；注意：这条路径现在读不到，宿主登记的那份副本文件并不存在，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestBothPointerDefectsAreReportedSeparately2928675678\\002\\nowhere\\gone.txt…]"
--- PASS: TestBothPointerDefectsAreReportedSeparately (0.01s)
=== RUN   TestPointerNoticeKeepsTheD153StubShape
--- PASS: TestPointerNoticeKeepsTheD153StubShape (0.00s)
=== RUN   TestHealthyReplyStillMatchesThePreFixTemplate
--- PASS: TestHealthyReplyStillMatchesThePreFixTemplate (0.01s)
=== RUN   TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3
    task_output_pointer_notice_test.go:439: leak-criterion verbatim: IsError=false Truncated=true announce="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，C26 连规范化都没通过（risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions）；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状），全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3492044899\\001\\ALLOWLIST-LEAK-PROBE-174r3…]"
    task_output_pointer_notice_test.go:470: notice region under the leak sweep: "注意：这条路径现在读不到，C26 连规范化都没通过（risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions）；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状）"
--- PASS: TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3 (0.06s)
=== RUN   TestLeakProbeTokenOnALegalPathStaysSilent174r3
    task_output_pointer_notice_test.go:507: legal-token control verbatim: IsError=false Truncated=true announce="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestLeakProbeTokenOnALegalPathStaysSilent174r34042844674\\001\\ALLOWLIST-LEAK-PROBE-174r3.txt…]"
--- PASS: TestLeakProbeTokenOnALegalPathStaysSilent174r3 (0.02s)
=== RUN   TestPointerAuthorityFollowsC26NotTheSpelling174r4
    task_pointer_authority_ac3_174r4_windows_test.go:121: two spellings of one file: long="C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerAuthorityFollowsC26NotTheSpelling174r42938359454\\001\\tool-output-174r4-alias.txt" short="C:\\Users\\swq\\AppData\\Local\\Temp\\TESTPO~1\\001\\TOOL-O~1.TXT"
    task_pointer_authority_ac3_174r4_windows_test.go:128: alias arm verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TESTPO~1\\001\\TOOL-O~1.TXT…]"
    task_pointer_authority_ac3_174r4_windows_test.go:157: alias arm, no root verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TESTPO~1\\001\\TOOL-O~1.TXT…]"
--- PASS: TestPointerAuthorityFollowsC26NotTheSpelling174r4 (0.02s)
=== RUN   TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4
    task_pointer_authority_ac3_174r4_windows_test.go:205: workspace-narrowed arm verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4575825906\\001\\artifacts174r4\\tool-output-174r4-narrowed.txt…]"
    task_pointer_authority_ac3_174r4_windows_test.go:224: workspace cleared verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4575825906\\001\\artifacts174r4\\tool-output-174r4-narrowed.txt…]"
--- PASS: TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4 (0.03s)
=== RUN   TestPointerCheckWritesNothing174r4
    task_pointer_authority_ac3_174r4_windows_test.go:265: ghost arm verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，宿主登记的那份副本文件并不存在，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerCheckWritesNothing174r4234859531\\001\\tool-output-174r4-ghost.txt…]"
    task_pointer_authority_ac3_174r4_windows_test.go:278: intact arm verbatim: IsError=false Truncated=true Text="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token，全文见 C:\\Users\\swq\\AppData\\Local\\Temp\\TestPointerCheckWritesNothing174r4234859531\\001\\tool-output-174r4-intact.txt…]"
--- PASS: TestPointerCheckWritesNothing174r4 (0.01s)
=== RUN   TestTaskState188AC2VocabularyIsD43Only
--- PASS: TestTaskState188AC2VocabularyIsD43Only (0.00s)
=== RUN   TestTaskState188AC2ValueComesFromTheRoster
--- PASS: TestTaskState188AC2ValueComesFromTheRoster (0.00s)
=== RUN   TestTaskState188AC2BackfillIsTheGoSideProducer
--- PASS: TestTaskState188AC2BackfillIsTheGoSideProducer (0.00s)
=== RUN   TestTaskState188AC3PanelHasNoWriteLeg
    task_state_188_test.go:364: AC#3 negative leg: 14 panel non-test sources examined, 0 write legs
    task_state_188_test.go:427: positive control roster-record-call red as required: planted.go:6 calls the roster/backfill write method Record
    task_state_188_test.go:427: positive control roster-backfill-call red as required: planted.go:6 calls the roster/backfill write method Backfill
    task_state_188_test.go:427: positive control taskoutput-state-key red as required: planted.go:9 constructs a TaskOutput with a State key
    task_state_188_test.go:427: positive control state-field-assign red as required: planted.go:9 assigns to a State field
    task_state_188_test.go:427: positive control whitelist-method-claim red as required: planted.go:3 claims a panel-side whitelisted method on the task axis: MethodTaskStateWrite
--- PASS: TestTaskState188AC3PanelHasNoWriteLeg (0.03s)
=== RUN   Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment
    tasklist_deferred_236r3_teeth_test.go:178: 现名册：task.output / task.cancel（DEFERRED 的判据读的是这一行，不是注释里的词面）
--- PASS: Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment (0.00s)
=== RUN   Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt
    tasklist_deferred_236r3_teeth_test.go:222: 成对读数：盘上 2 行 / 编译期 2 行（HEAD 上两数相等＝这把尺没在空转）
--- PASS: Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt (0.00s)
=== RUN   TestTaskOutputStampedOnRealBridge175r2
--- PASS: TestTaskOutputStampedOnRealBridge175r2 (0.00s)
=== RUN   TestNonContentResultStaysUnmarked175r2
--- PASS: TestNonContentResultStaysUnmarked175r2 (0.01s)
=== RUN   TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2
--- PASS: TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2 (0.03s)
=== RUN   TestHostMintedPointerRereadStaysCleanOnRealBridge175r2
--- PASS: TestHostMintedPointerRereadStaysCleanOnRealBridge175r2 (0.01s)
=== RUN   TestEveryRegisteredToolIsClassifiedForMarking175r2
    ticket175r2_stamp_live_test.go:374: fs.edit: unstamped: a confirmation of a local edit
    ticket175r2_stamp_live_test.go:374: fs.trash: unstamped: a confirmation of a local delete
    ticket175r2_stamp_live_test.go:374: fs.move: unstamped: a confirmation of a local move
    ticket175r2_stamp_live_test.go:374: fs.delete: unstamped: a confirmation of a local delete (opt-in family)
    ticket175r2_stamp_live_test.go:374: task.output: marked: C25 roster since ticket 175-r2, its answer is what a task printed
    ticket175r2_stamp_live_test.go:374: fs.read: marked: C25 source roster (SPEC-06 §5)
    ticket175r2_stamp_live_test.go:374: fs.list: unstamped: a listing this process composed, no foreign body
    ticket175r2_stamp_live_test.go:374: fs.write: unstamped: a confirmation of a local write
    ticket175r2_stamp_live_test.go:374: task.cancel: unstamped: a receipt naming which roster row this call stopped
--- PASS: TestEveryRegisteredToolIsClassifiedForMarking175r2 (0.00s)
=== RUN   TestBackfillFilesTheRealSpilledArtifact176r1
--- PASS: TestBackfillFilesTheRealSpilledArtifact176r1 (0.01s)
=== RUN   TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1
--- PASS: TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1 (0.01s)
=== RUN   TestNoRosterWriteAfterCancel176r1
--- PASS: TestNoRosterWriteAfterCancel176r1 (0.00s)
=== RUN   TestTerminalWriteCarrierStaysOutOfScope176r1
--- PASS: TestTerminalWriteCarrierStaysOutOfScope176r1 (0.00s)
=== RUN   TestTicket224GrantIDColumnRoundTripsThroughTheRealStore
time=2026-10-09T12:01:05.491+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224GrantIDColumnRoundTripsThroughTheRealStore1414536130\002\data\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:05.500+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:05Z duration_ms=8
time=2026-10-09T12:01:05.507+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224GrantIDColumnRoundTripsThroughTheRealStore1414536130\002\data\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:05.512+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:05Z duration_ms=4
time=2026-10-09T12:01:05.521+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:05.523+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224GrantIDColumnRoundTripsThroughTheRealStore (0.07s)
=== RUN   TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter
--- PASS: TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter (0.00s)
=== RUN   Test283TaskIDAccessorIsTheTaskIDNotTheCorr
--- PASS: Test283TaskIDAccessorIsTheTaskIDNotTheCorr (0.00s)
=== RUN   Test283IdentityChainThroughTheRealLoop
time=2026-10-09T12:01:05.563+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test283IdentityChainThroughTheRealLoop3299523047\001\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:05.571+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:05Z duration_ms=7
time=2026-10-09T12:01:05.584+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test283IdentityChainThroughTheRealLoop3299523047\001\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:05.591+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:05Z duration_ms=7
time=2026-10-09T12:01:05.596+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:05.599+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T12:01:05.608+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f95f940b-ed7d-4bf6-b3d8-2e0332fcdfd6 owner=tools
--- PASS: Test283IdentityChainThroughTheRealLoop (0.07s)
=== RUN   Test285RosterRowsCarryDistinctCorrPerCall
time=2026-10-09T12:01:05.637+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test285RosterRowsCarryDistinctCorrPerCall3630431616\001\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:05.645+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:05Z duration_ms=8
time=2026-10-09T12:01:05.653+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test285RosterRowsCarryDistinctCorrPerCall3630431616\001\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:05.658+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:05Z duration_ms=4
time=2026-10-09T12:01:05.663+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:05.665+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T12:01:05.675+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-88bd18f9-4c7e-448d-93d5-8a9e8d712648 owner=tools
--- PASS: Test285RosterRowsCarryDistinctCorrPerCall (0.07s)
=== RUN   TestTicket90ModeReachesRoutingAsParameter
--- PASS: TestTicket90ModeReachesRoutingAsParameter (0.00s)
=== RUN   TestTicket90TwoBridgesDoNotShareOneMode
--- PASS: TestTicket90TwoBridgesDoNotShareOneMode (0.00s)
=== RUN   TestTicket90UnwiredModeSourceIsTheStrictestMode
--- PASS: TestTicket90UnwiredModeSourceIsTheStrictestMode (0.00s)
=== RUN   TestTicket90IrreversibleStillAsksInEveryMode
--- PASS: TestTicket90IrreversibleStillAsksInEveryMode (0.00s)
=== RUN   TestTicket90TaintEscalationNeverSilenced
--- PASS: TestTicket90TaintEscalationNeverSilenced (0.00s)
=== RUN   TestTicket90TierADenySurvivesEveryMode
--- PASS: TestTicket90TierADenySurvivesEveryMode (0.02s)
=== RUN   TestTicket90ModeCarriesNoAllowAuthority
--- PASS: TestTicket90ModeCarriesNoAllowAuthority (0.00s)
=== RUN   TestTicket90BlacklistGateIsCalledOnLivePath
--- PASS: TestTicket90BlacklistGateIsCalledOnLivePath (0.01s)
=== RUN   TestLoopPassesDeclaredL1WriteThroughTheGate
time=2026-10-09T12:01:05.736+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestLoopPassesDeclaredL1WriteThroughTheGate2153359646\001\db\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:05.743+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:05Z duration_ms=6
time=2026-10-09T12:01:05.751+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestLoopPassesDeclaredL1WriteThroughTheGate2153359646\001\db\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:05.758+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:05Z duration_ms=6
time=2026-10-09T12:01:05.762+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:05.765+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestLoopPassesDeclaredL1WriteThroughTheGate (3.08s)
=== RUN   TestLoopStillRefusesL1WhenNoGateIsRegistered
time=2026-10-09T12:01:08.817+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestLoopStillRefusesL1WhenNoGateIsRegistered2243834630\001\db\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T12:01:08.823+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T04:01:08Z duration_ms=5
time=2026-10-09T12:01:08.830+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestLoopStillRefusesL1WhenNoGateIsRegistered2243834630\001\db\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T12:01:08.836+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T04:01:08Z duration_ms=6
time=2026-10-09T12:01:08.841+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T12:01:08.843+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestLoopStillRefusesL1WhenNoGateIsRegistered (0.06s)
=== RUN   TestL1WriteGoesThroughTheRealBlockWindow
--- PASS: TestL1WriteGoesThroughTheRealBlockWindow (3.02s)
=== RUN   TestVetoInsideTheWindowWritesNothing
--- PASS: TestVetoInsideTheWindowWritesNothing (0.01s)
=== RUN   TestLateVetoRendersTheApprovalLayersAppliedStepsReport
--- PASS: TestLateVetoRendersTheApprovalLayersAppliedStepsReport (3.01s)
=== RUN   TestApprovedL2ThatStopsMidWriteIsNotRenderedAsNeverStarted
--- PASS: TestApprovedL2ThatStopsMidWriteIsNotRenderedAsNeverStarted (0.01s)
=== RUN   TestFSReadOnlyNeverOpensACard
--- PASS: TestFSReadOnlyNeverOpensACard (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/tools	14.006s
