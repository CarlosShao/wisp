# 票 124 终裁腿 124-v1 裁决 part3（sec 五..八）

> 前接 `docs/evidence/s1/124-v1-verdict.md`（sec 0 与 sec 一）与 `124-v1-verdict-part2.md`（sec 二..四）。

## 五、红名册逐名

**批 1 已清零的 33 名（PRE-LINK 减 PRE-PLAIN，扣 2 枚摇摆名；POST 侧逐名归 0）**
Test313CreateProviderRowCreatesAPresetRowAndItsModel, Test313CreateProviderRowRefusesNonPresetAndExistingRow, Test313RealtimeHalfPairFollowsTheEnabledRule, Test313RoleChatPairWritesBothKeysOrNeither, Test313VoiceCloudChainNeedsTheCatalogRowAndRefusesLocalSherpa, Test320CreateCustomProviderRowLandsAndRefusesByName, Test320CreateModelRowSeedsFromDefaultsAndRefusesTwice, Test320SetProviderBillingWritesTheEnumLeafAndRefusesOthers, TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit, TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure, TestAC11SettingWriteKeepsAForeignHandEditedKey, TestAC11SettingWriteReportsTheKeyPathItChanged, TestAC11WritingTheValueTheFileAlreadyHoldsChangesNoBytes, TestAC1AllowedDirsWriteKeepsAHandEditedKey, TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit, TestAC2CleanWriteReportsOneKeyAtInfoLevel, TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept, TestAC3AdoptionClaimsOnlyWhatThisWriteProduced, TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory, TestAC3SettingWriteStoresARefNeverAValue, TestAC4PermissionModeWriteKeepsAHandEditedKey, TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber, TestDeletedLineRestorationIsNotClaimedAsAWrittenKey, TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile, TestHandDeletedDefaultValuedLineIsRestoredInvisibly, TestHandDeletedMapEntryIsNotRestored, TestHandDeletedProviderBaseUrlComesBackAsThePresetValue, TestHandDeletedScalarLineComesBackAtSchemaDefault, TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory, TestTicket257R1AC1CleanMachineRefusesAllSevenFields, TestTicket257R1AC1HandAddedRowsUnlockAllSevenFields, TestTicket257R1AC2MissingRowIsItsOwnSentence, TestTicket261P1SettingsWriteMaterializesFalse
（33 枚按 8 文件分布：`writeguard_226` 8、`settings_248` 6、`writeguard_unnailed_341` 5、`settings_313` 5、`writeguard_handedit_227` 4、`settings_320` 3、`settings_257` 3、`enabled_261` 1；逐名 file:line 表 `/d/tmp/wisp/20261021/124v1/rosters/config-name-file.tsv`）

**批 1 剩下那 2 名（摇摆，非本腿引入；两形都在）**
TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent, TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed

**批 2 `internal/tools` 38 名（16 文件；扣票 123 三枚 300s 腿后 35 枚可转）**
Test165PlanApprovalRingsAgainForAnUnlistedWrite, Test172DocCopyThatLiesReddensTheRuler, Test172FrozenTierMatchesTheTierTheGateGives, Test20r2DirectionEngineAdjudicatesFileLevelAdditions, Test20r2FirstUseAskFlowGoesOnlyThroughConfig, Test221ParentStopsItsOwnChildRowAndStreamSettle, Test283IdentityChainThroughTheRealLoop, Test285RosterRowsCarryDistinctCorrPerCall, Test287DepthRefusalWinsOverThePoolRefusal, Test287RootCallerWithTheSameShapeStillDerives, Test287SubagentCallerRefusedAtTheRosterDepthCheckThroughTheRealLoop, TestBackfillFilesTheRealSpilledArtifact176r1, TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2, TestFSReadOnlyNeverOpensACard, TestHealthyPointerStaysSilent, TestHealthyReplyStillMatchesThePreFixTemplate, TestHostMintedPointerRereadStaysCleanOnRealBridge175r2, TestL1WriteGoesThroughTheRealBlockWindow, TestLateVetoRendersTheApprovalLayersAppliedStepsReport, TestLeakProbeTokenOnALegalPathStaysSilent174r3, TestLongOutputPointerRecoversEveryByte, TestNonContentResultStaysUnmarked175r2, TestPointer183BackfilledArtifactRereadStaysClean, TestPointer183SiblingArtifactInTheSameDirectoryStillHits, TestPointer185AModelNominatedFileIsNotRostered, TestPointer185AnotherTaskCannotBorrowTheRoster, TestPointer185ForeignFileCarryingTheHostPathStillBlocksTheReread, TestPointer185PagedRereadsKeepWorking, TestPointer185SecondRereadOfTheSameArtifactStaysClean, TestPointer185SiblingHostArtifactIsNotLaundered, TestPointerPast256KiBIsNotFullyReadable, TestPointerToMissingCopyFileSpeaks, TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1, TestTicket224GrantedCallBooksAllowSessionGrantWithItsRow, TestTicket224GrantIDColumnRoundTripsThroughTheRealStore, TestTicket224GrantSourceThatCannotAnswerMeansAsking, TestTicket224LiveGrantStopsTheL1Question, TestTicket224PartialPathCoverageStillAsks

票 123 那族三枚 300s 腿**逐名在册**（我一手量到时延）：`TestL1WriteGoesThroughTheRealBlockWindow`(300.03s)、`TestLateVetoRendersTheApprovalLayersAppliedStepsReport`(300.03s)、`TestFSReadOnlyNeverOpensACard`(300.04s)。
文件分布：`pointer_185_cli_seam` 6、`task_output_pointer_notice` 4、`grant` 4、`wiring` 3、`ticket175r2_stamp_live` 3、`subagent_depth_fixture_287` 3、`tier_consistency_172_cases` 2、`ticket176r1_start_port` 2、`task_output_leg` 2、`refusal_20r2` 2、`pointer_183_cli_seam` 2、`ticket285_corr_rows_rulers` 1、`ticket283_corr_identity_rulers` 1、`ticket224_grantid_column` 1、`task_cancel_221_legs` 1、`plan_present_165` 1。
子测试红 3 枚：`Test20r2FirstUseAskFlowGoesOnlyThroughConfig/approved:_the_append_lands_via_Manager.AddAllowedDir_and_the_next_assembly_passes`、`TestTicket224GrantSourceThatCannotAnswerMeansAsking/panics_mid-lookup`、`TestTicket224GrantSourceThatCannotAnswerMeansAsking/reports_that_it_could_not_look`。

**批 3 `cmd/wisp` 37 名（13 文件）**
Test313EmptyConfigToStepfunRowChainAndPair, Test313MultiElementCommaChainKeepsEveryElement, Test320CustomProviderRowThroughTheRealLeg, Test320ModelRowCreateThroughTheRealLeg, Test320ModelRowDeleteThroughTheRealLeg, Test320ModelsListFetchesThroughTheRealLeg, Test320ModelsListNeedsAProvider, Test320ModelsListRefusalsDoNotTouchTheNetwork, Test320ProviderBillingThroughTheRealLeg, Test320ProviderRowDeleteRefusesWhileReferenced, TestAC1InboundLegAnswersSettingsRouteEndToEnd, TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket, TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere, TestAC2SnapshotReportsRefWithoutBlobAsAPartState, TestAC7InvalidSettingsValueIsRefusedBeforeTheFile, TestAC9InboundFlagSurfaceIsNarrow, TestAC9InboundLegFromStdinReachesTheWriteLeg, TestAC9InboundLegRefusesRosterMethodWithNoHandler, TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt, TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables, TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten, TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone, TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences, TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags, TestTicket198R2AC4ReceiptNamesTheRealEntryPoints, TestTicket198R2J1TheAssemblyRootStillCreatesNothing, TestTicket223PanelInboundSaysHotReloadIsDisabled, TestTicket223R2FailureSentenceRouting, TestTicket255RestartTierKeysAreBackedByATest, TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain, TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads, TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields, TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven, TestTicket257R2AC1ReceiptStatesTheNonPresetCondition, TestTicket257R2AC2EffectTimingSaysThreeProcessShapes, TestTicket257R2AC2ThreeRefusalsStayThreeSentences, TestTicket257R2AC3CredentialSurfaceUntouched

文件分布（共 **13** 文件，`UNKNOWN=0`）：`config_reload_223`、`config_sentences_223r2`、`firstrun_198`、`firstrun_198r2`、`firstrun_257_nonpreset`、`firstrun_257`、`panel_config_248`、`panel_config_313`、`panel_config_320`、`panel_delete_320`、`panel_inbound_33`、`panel_models_320`、`restart_tier_keys_255r2`。
票 123 三枚豁免腿在 `cmd/wisp` 名册 **0 命中**，300s 腿 **0 枚** => 本批无豁免项。

## 六、批 2 / 批 3 可行性读数

- **`internal/tools`（38 枚 / 16 文件）逐名清单齐全**：我自己从两形名册 `comm` 量出 **38**，与实现腿自报的 38/16 **逐数对上**；
  16 文件逐枚落到定义文件，`UNKNOWN=0`。**两形都红 = 0**（该包普通形 PASS=212 FAIL=0 SKIP=4）=> 分母干净。
  唯一需要事先钉死的是那三枚 300s 豁免腿在册（会拖长台面：本发 `-timeout 60m` 跑了 **914.7 s**）。
- **`cmd/wisp`（37 枚 / 13 文件）逐名清单齐全**：我自己量出 **37**/13，与自报对上，`UNKNOWN=0`。
  一处我必须上报的分歧：实现腿与台账记的"同包两形都红 **28 枚**"，我这一发量到顶层 **68 枚** 加子测试 **9 枚**。
  逐行看失败串，几乎全是 `secret: DPAPI is only available on Windows; use an env: ref on this platform`
  => 这是 `internal/secret/protect_other.go:19` 的 POSIX 结构性不可用，**不是票 124 那族**，也一枚都不许顺手修。
  68 与 28 的差额我不归因（可能 2b-4 当时只数了"含 not provably resolved 串"的那一部分），但**批 3 的分母必须按 68 立账**，
  否则"两形都红 29 枚"那条豁免名单会低估约 40 枚、把包级 `rc` 读成破口。
- **`cmd/wisp` 能不能与 `tools` 同批**：**技术上能，治理上不许**。
  我现量依赖方向：`internal/tools` **不** import `cmd/wisp`（`go list` 的 Imports 里 `wisp/cmd` 0 命中），
  `cmd/wisp` **import** `internal/tools`（命中）=> 两包无环、编译上互不阻塞。
  但 AGENTS.md 明写 `cmd/wisp/**` 与 `internal/panel/**` 是**共享编译面、每波只许一枚写腿**，
  且 `go test ./cmd/wisp` 重包同一时刻只许跑一个 => **必须单排**，不能与 `tools` 同批。
  另：批 3 那 37 枚全在 `cmd/wisp/**/*_test.go`，与票 133 在动的 `leg_dispatch_gate_133_test.go` **零文件交集**
  （我按文件级核过：13 枚文件里没有它）=> 文件级冲突这一条今天是干净的。

## 七、建议编排者翻哪几格（我自己一格未翻）

1. **`AC#3`：建议翻**（批 1 这一发的证据链我四条尺全过，且 33 枚绿名改前普通形就 PASS 这条是决定性的）。要票级口径就等批 3 合并态，但那时也只会更成立。
2. **`AC#4`：建议翻，但附一条更正**——翻的理由是三态我独立复现（42 / 6 名 / 41）；更正＝它缺 `sha256` 还原证词，翻格时请把"还原靠 `git status` 干净证明"这句写进交件，别让"sha256 已核"这种没发生的动作留在记录里。
3. **`AC#2b` / `AC#2` / `AC#5` / `AC#6`：建议 NOT 翻**（理由见 part2 sec 二表）。
4. **建议补一条账**：把那 2 枚 POSIX-only 摇摆登进 `pending-and-issues.md`（part2 sec 三），并同时登记"整包普通形红数在 0..2 之间非确定"这一事实——它是出口判据"普通形一枚不许多红"的根本限制。
5. **建议更正两处数**：(a) 本腿"6 处递根点" => **9 处 / 9 hunk**（件 sec 3 那张表自己就是 9，派单的 6 同源作废）；
   (b) `124r1-ruling-batch.txt:5` 的"19 处调用 + 1 行注释" => **18 处调用 + 2 行注释**（合 20 不变）。

## 八、提交号、逐枚 `.go`、越界与正控

- **本裁腿提交号**：见本件所在 commit（主题以 `票 124 终裁腿 124-v1` 开头）。只 commit NOT push。
- **被裁两笔逐枚 `.go`**：
  - `d32b459c`：4 枚路径、**4 枚 `.go`**，全是 `internal/config/**/*_test.go`（`enabled_261_test.go` / `manager_test.go` / `settings_257_test.go` / `writeguard_226_test.go`），**越界 0 枚**，产码 0 枚。
  - `abd09226`：3 枚路径、**0 枚 `.go`**，全在 `docs/evidence/s1/124-r1-impl*.md`。
- **正控（同一把尺打真含 `.go` 的提交）**：`fafe2b4a` => 9/9 是 `.go`；`dfa3dc4` => 18/18；`5265c3a` => 3/3；`62dda11` => 8/8。=> 尺是活的、不是恒 0。
- **负控（同上尺打纯文档提交）**：`abd09226` 0/3、`8b4989c0` 0/1、`3d7ff365` 0/1、`4ac4549c` 0/1。
- **本裁腿写面**：只 `docs/evidence/s1/124-v1-verdict.md` 与本三份。零 `.go`、零 `internal/**`、零 `cmd/**`、零 `scripts/**`、零 `.scratch/**`、零 `docs/reports/**`、零 `docs/PLAN.md`、零 `docs/specs/**`。
- 临时件只建不删（`D:\tmp\wisp\20261021\124v1\**` 全部保留）；未建 worktree/checkout；未用 `-overlay`；未建 `.out`。
