# 票 124 落地腿 124-r1 交件 part2（sec 6：变异自证与逐名红册）

> 前接 `docs/evidence/s1/124-r1-impl.md`（sec 0..5）。

## 6. 变异自证（票面 AC#4）

`MUTATION-124-R1`：把 `sealableTempDir124` 的函数体改回原样交 `t.TempDir()`（6 个调用点不动）。

先证**落地**，再读数：

```
$ grep -c 'proc.SealableRoot(t.TempDir())' internal/config/tempdir_resolved_124_test.go
0                                    <- 解析那一跳命中 0，确已摘掉
$ grep -n 'MUTATION-124-R1' internal/config/tempdir_resolved_124_test.go
33:	return t.TempDir() // MUTATION-124-R1
$ go build ./internal/config/    -> rc=0
$ go vet   ./internal/config/    -> rc=0     （先补 `_ = proc.SealableRoot` 保住 import，
                                              否则 rc=1 是 unused import 的编译错，
                                              不是行为读数——第一发就是这么撞的）
```

读数：`MUT-LINK FAIL=42`。**恢复**后 `RESTORED-LINK FAIL=2`，且那 2 枚是改前普通形也红的**既有**项（见 sec 7）。

红名**逐名**（变异把基线 35 名全数带回，另有 6 名新增，全部是被同一枚 helper 喂根的用例）：

```
回归的基线 33 名（only-in-link 那批，逐名）：Test313CreateProviderRowCreatesAPresetRowAndItsModel,
Test313CreateProviderRowRefusesNonPresetAndExistingRow, Test313RealtimeHalfPairFollowsTheEnabledRule,
Test313RoleChatPairWritesBothKeysOrNeither, Test313VoiceCloudChainNeedsTheCatalogRowAndRefusesLocalSherpa,
Test320CreateCustomProviderRowLandsAndRefusesByName, Test320CreateModelRowSeedsFromDefaultsAndRefusesTwice,
Test320SetProviderBillingWritesTheEnumLeafAndRefusesOthers, TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit,
TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure, TestAC11SettingWriteKeepsAForeignHandEditedKey,
TestAC11SettingWriteReportsTheKeyPathItChanged, TestAC11WritingTheValueTheFileAlreadyHoldsChangesNoBytes,
TestAC1AllowedDirsWriteKeepsAHandEditedKey, TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit,
TestAC2CleanWriteReportsOneKeyAtInfoLevel, TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept,
TestAC3AdoptionClaimsOnlyWhatThisWriteProduced, TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory,
TestAC3SettingWriteStoresARefNeverAValue, TestAC4PermissionModeWriteKeepsAHandEditedKey,
TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber, TestDeletedLineRestorationIsNotClaimedAsAWrittenKey,
TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile, TestHandDeletedDefaultValuedLineIsRestoredInvisibly,
TestHandDeletedMapEntryIsNotRestored, TestHandDeletedProviderBaseUrlComesBackAsThePresetValue,
TestHandDeletedScalarLineComesBackAtSchemaDefault, TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory,
TestTicket257R1AC1CleanMachineRefusesAllSevenFields, TestTicket257R1AC1HandAddedRowsUnlockAllSevenFields,
TestTicket257R1AC2MissingRowIsItsOwnSentence, TestTicket261P1SettingsWriteMaterializesFalse

变异额外翻红 6 名（同一枚 helper 的其它被喂用例，属同族、不是新缺陷）：
TestMigrateNoVersionKeyAssumedV1, TestMigrateV1Fixture, TestResolvedNeverPersists,
TestRoundTripLoadMarshalLoad, TestSaveFileRoundTripsThroughLoad, TestUnwiredGuardLeavesHonestConfigsAlone
```

变异红行**点到用例自己**（样例逐字，两行）：

```
boundary_test.go:191: save: config: config.toml temp file seal: winsec: refusing to seal
  /varlink/w124tmp/TestRoundTripLoadMarshalLoad4294828850/001/.wisp-config-2871620155.tmp: the installed risk.c26Pipeline answered ...
boundary_test.go:239: config: config.toml temp file seal: winsec: refusing to seal
  /varlink/w124tmp/TestResolvedNeverPersists4270695038/001/.wisp-config-2422300932.tmp: ...
```
