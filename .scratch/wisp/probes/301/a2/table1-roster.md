# 表① rosters @ blob 86d89478 (invisible-5 detail)
--- pkg internal/agent
  FILE internal/agent/spill_acl_windows_test.go
        CASE TestAC3SpillArtifactLandsPrivate
--- pkg internal/audio
  FILE internal/audio/capturelevel_windows_test.go
        CASE TestAC247RealCaptureLoopEmitsLevels
  FILE internal/audio/hotplug_test.go
        CASE TestHotplugReenumerateOnce
        CASE TestHotplugStaleHandleFailsLoudly
        CASE TestHotplugReenumerateFailureNamesDevice
        CASE TestOpenOccupiedAndPermissionDenied
        CASE TestEndpointsPairQueryable
        CASE TestPinnedThreadStable10s
        CASE TestLiveWasapiSmoke
        CASE TestCaptureLoopWavIntegrity
--- pkg internal/memory
  FILE internal/memory/artifacts_junction_tripwire_windows_test.go
        CASE TestAC2ReclaimWalkNeverDescendsIntoAJunction
--- pkg internal/models
  FILE internal/models/acl_sid_121_windows_test.go
  FILE internal/models/handoff_window_109_windows_test.go
        CASE TestAC1TheWriteWindowIsAnInheritedCrossAccountRight
        CASE TestAC4InstallDirWidthIsDeliberateAndHasABreakLine
  FILE internal/models/no_seal_ruling_windows_test.go
        CASE TestAC3ExtractionIsDeliberatelyNotSealed
--- pkg internal/tools
  FILE internal/tools/bridge_a18_kill_windows_test.go
        CASE TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile
  FILE internal/tools/bridge_junction_windows_test.go
        CASE TestBridgeRefusesARealJunctionOnTheReadRoute
        CASE TestBridgeWritesNothingThroughARealJunction
        CASE TestJunctionInsideAnAllowedRootCannotReachAnAListFile
        CASE TestBridgeRefusesTheRealShortNameOfAnAListFile
        CASE TestShortNameSpellingGetsTheSameVerdictAsTheLongOne
        CASE TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing
  FILE internal/tools/fs_edit_ac4b_kill_windows_test.go
        CASE TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue
  FILE internal/tools/fs_edit_ac5_sweep_r4_windows_test.go
        CASE TestFSEditWriteReclaimsADeadWritersOrphan
  FILE internal/tools/fs_staging_windows_test.go
        CASE TestSweepReclaimsOnlyItsOwnStagingFiles
        CASE TestSweepNamingSchemeRejectsNearMisses
        CASE TestSweepNeverDeletesThroughARealJunction
        CASE TestSweepSparesATempFileHeldByAnotherProcess
  FILE internal/tools/paths_shortname_252_probe_test.go
        CASE TestTicket252P1AllowlistLegs
        CASE TestTicket252P1OwnerMachineReachability
        CASE TestTicket252P1EnvSpellingCensus
  FILE internal/tools/paths_shortname_252_r1_test.go
        CASE TestTicket252R1ShortSpellingOfNewFileIsAuthorized
        CASE TestTicket252R1BothSpellingsAnswerTheSame
        CASE TestTicket252R1AlignmentAddsNoAuthorization
        CASE TestTicket252R1CanonicalStillNamesOneTree
  FILE internal/tools/paths_twocontainments_252_r2_windows_test.go
        CASE TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses
  FILE internal/tools/recycle_windows_test.go
        CASE TestRecycleBinStructLayoutMatchesWin32
        CASE TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate
        CASE TestShellTrashRefusesAMissingPathWithoutDeletingAnything
  FILE internal/tools/task_pointer_authority_ac3_174r4_windows_test.go
        CASE TestPointerAuthorityFollowsC26NotTheSpelling174r4
        CASE TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4
        CASE TestPointerCheckWritesNothing174r4
