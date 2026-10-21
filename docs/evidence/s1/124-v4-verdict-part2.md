# 票 124 终裁腿 124-v4 裁决 part2（sec 4..7）

> 前接 `docs/evidence/s1/124-v4-verdict.md`（sec 0..3：台面、八项判语、逐格结论、判决）。

## 4. 门禁 rc 与末行逐字

```
$ GOFLAGS= go test -count=2 -v -timeout 40m ./internal/config/   -> rc=0
ok  	github.com/CarlosShao/wisp/internal/config	7.369s
$ GOFLAGS= go test -count=2 -v -timeout 40m ./internal/tools/    -> rc=0
ok  	github.com/CarlosShao/wisp/internal/tools	61.853s
$ go build ./...                                                 -> rc=0（0 行输出）
$ go vet ./internal/config/ ./internal/tools/ ./cmd/wisp/         -> rc=0（0 行输出）
容器 golang:1.27：go build ./... rc=0 / vet config+tools rc=0 / vet cmd/wisp rc=0 / gofmt -l 0 行
$ sh scripts/d22scan.sh                                          -> rc=0
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42,
ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116,
ban #8 internal/=610, ban #8 cmd/=145
$ sh scripts/check-path-length-budget.sh                         -> rc=0
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a
reason, and the roster equals the tree
```

`cmd/wisp` 整包（**按 `A1100` NOT 当判据**）：宿主 `-count=1 -v` rc=1，末行逐字
`FAIL	github.com/CarlosShao/wisp/cmd/wisp	302.250s`；容器 `-count=2 -v` 软链形 rc=1，
末行逐字 `FAIL	github.com/CarlosShao/wisp/cmd/wisp	<各发时长>`；两份红名册双向 **0 行**。

## 5. 红名册逐名（我一手量，NOT 抄任何腿）

**(a) 票 123 那族豁免三枚（`internal/tools`；POST 软链形仅剩这 3 名，`wiring_test.go` 本票零 hunk）**

| 名 | 址（HEAD blob） | 我量到的时延 | 归口 |
|---|---|---|---|
| `TestL1WriteGoesThroughTheRealBlockWindow` | `wiring_test.go:112` | 300.08s | 归票 123 |
| `TestLateVetoRendersTheApprovalLayersAppliedStepsReport` | `wiring_test.go:182` | 300.04s | 归票 123 |
| `TestFSReadOnlyNeverOpensACard` | `wiring_test.go:291` | 300.13s | 归票 123 |

三枚**改前改后都在红**、失败文案 ＝ `审批超时（300 秒未确认），C18 一律判拒绝` 那一族（NOT 本族）
=> 豁免成立、NOT 新增红、NOT 顺手修。

**(b) `internal/config` 两形俱在的已知 POSIX 摇摆（NOT 本批引入；house 账已登）**

`TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent`（`writeguard_handedit_227_test.go:106`）、
`TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed`（`writeguard_unnailed_341_test.go:209`）。
两枚宿文件**均不在批 1 的 4 枚改动文件里**。我这一发：`PRE-P` 只红第一枚、`POST-P` 两枚都红 => 摇摆方向不可预测，
但它**不在任何一批的 went-green 里**（第一枚在 `PRE-L` 也红，第二枚两形都红）。

**(c) `cmd/wisp` 两形都红 68 名（顶层；我的 `POST-P` 名册，与 r3 自报 set-equal）**

Test165r2AssembledBridgeReachesPlanPresent, Test165r2UnlistedWriteRingsItsOwnCardAfterAnApprovedPlan,
Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands, TestAC1AC2DispatchHopGate133,
TestAC1AlwaysBranchDoesNotRevertAHandEditedKey, TestAC2CredentialSentinelAppearsInNoArtifact,
TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef, TestAC4EveryLegIsNailedOrRuled,
TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit, TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card,
TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused, TestBookedSnapshotReachesHostBitForBit329,
TestCleanCheckoutBuilds_AC11, TestComposedGateBlocksAWriteForTwoSeconds, TestCompositionRootClosesTheLoopTasksTaintScope,
TestHostDispatchThroughTheAssembledBridge, TestL1VetoNeedsAChannelTheHostReallyWired, TestMissingBlobFailsUnconfiguredNeverSilently,
TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole, TestNativeHostSeamRefusesAPanelSourcedAllow,
TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject, TestReplyListenerAllowsAnL2CardFromTheNativeSide,
TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel, TestRunBooksWithASnapshotOfItsLiveQueue,
TestRunPacketCarriesTheLoadedInstructionFiles, TestRunPacketCarriesTheSubagentItsRosterRowFed,
TestRunPacketLightsTheRosterRowOfItsLiveL1Window, TestRunPacketMarksTheRosterRowACardIsHolding,
TestRunPacketReportsTheStreamLogPastItsBound, TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound,
TestRunTextTaskFailNextIsClassified, TestRunTextTaskKeyResolvesInTheStore, TestRunTextTaskTextPathEndToEnd,
TestSecretEndToEndConfigRefResolvesAtRequestTime, TestSecretFailurePathsLogAndPrintNoPlaintext,
TestSecretFromStdinWritesNoIntermediateFile, TestSecretOverwriteIsAnnounced, TestSecretPortableModeUsesTicket06Seam,
TestSecretSameNameUnderThreeEnvsIsThreeBlobs, TestSecretSetGetListRoundTrip, TestSecretUnsetRefusesWhileReferenced,
TestSnapshotWorkspaceSectionReportsTheNarrowing, TestTaskOutputPointerNotice174r8SpeaksTheRealReason,
TestTaskOutputPointerNotice174r8StaysQuiet, TestThePumpIsDrivenNotJustAssembled, TestTicket101ManualSwitchSurvivesRestart,
TestTicket101ModeSwitchUsesTheRealL2Gate, TestTicket101SessionGrantDoesNotCrossRestart,
TestTicket101UnreadableModeFailsLoudlyAndStrict, TestTicket101UntouchedConfigRestartsAtDefault,
TestTicket223FailureSentencesAreDistinct, TestTicket223HandEditedFsLooseningCostsAnL2Card,
TestTicket223ModeLooseningChangesTheRunningModeAfterAllow, TestTicket223RefusedLooseningKeepsOldValues,
TestTicket223RestartTierSaysItWillNotApply, TestTicket223RunArmsTheReloadTick, TestTicket223TighteningRaisesNoCard,
TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking, TestTicket224ProductionSessionDoesNotSurviveRestart,
TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun, TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere,
TestTicket255ReceiptOmitsPanelFromTheImmediateSentence, TestTicket255ReceiptSentenceAssemblyIsFiltered,
TestTicket255ReceiptStillNamesTheLiveReadSection, TestUnansweredL2CardTimesOutIntoRejectNeverExecution,
TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask, TestZZ315A1RootLoopFirstBridgeCall,
TestZZ315A1SiblingChildFirstBridgeCall

子测试 9 枚：`TestL1VetoNeedsAChannelTheHostReallyWired/{console_posture:_veto_refused,_window_still_executes,host_declares_and_loads_esc,_veto_lands}`、
`TestSecretFailurePathsLogAndPrintNoPlaintext/{bad_dpapi_decrypt_(corrupted_blob,_non-portable),bad_dpapi_decrypt_(portable,_P13),store_write_failure_surfaces_the_ref_only}`、
`TestTicket223FailureSentencesAreDistinct/{schema未知键,声明了版本但坏在后面,缺失,语法错}`。
机制＝DPAPI 在 POSIX 结构性不可用（`internal/secret/protect_other.go:19`）＋票 131/133 的两道门＋
`TestCleanCheckoutBuilds_AC11` 需要真 git 仓（`git rev-parse --show-toplevel: exit status 128`）。
**一枚都不许顺手修**，也不计入任何出口判据。

**(d) 宿主 Windows 档**：三发红名并集 **1 枚**＝`TestAC14GoSideEvalPushReachesThePage`（真窗族；同址同文在未改 HEAD 上照红）。

**(e) 变异的同族连带（NOT 新增红）**：我重打的三发变异里，目标名逐名全回归之外另有连带——
config `34 => 40`（+6）、tools `38 => 62`（+24）、wisp `37 => 41`（+4）。
四枚 wisp 连带名逐字为 `TestProvidersDiscoverListsWhatTheServerServes`、`TestProvidersProbeRecordsMeasuredThinkingFalse`、
`TestProvidersProbeRecordsMeasuredThinkingTrue`、`TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless`——
它们与批 4 那 5 枚同喂 `providers_test.go:40` 那一枚 helper，摘掉解析必然一起红。
=> 判据是"目标名必须全覆盖 + 连带只许增不许减"，NOT 总数相等。

## 6. 提交号、逐枚 `.go`、越界与正控

| 提交 | 件数 | `.go` | 内容 | 越界 |
|---|---|---|---|---|
| `d32b459c`（批 1 码） | 4 | **4** | `internal/config/{enabled_261,manager,settings_257,writeguard_226}_test.go` | 0 |
| `abd09226`（批 1 件） | 3 | 0 | `docs/evidence/s1/124-r1-impl*.md` | 0 |
| `99f1f9c7`（批 2 码） | 15 | **15** | `internal/tools/**/*_test.go` | 0 |
| `d055326e`／`c832ac3f`（批 2 件） | 3／4 | 0 | `124-r2-impl*.md` | 0 |
| `6b712571`（批 3 码） | 9 | **9** | `cmd/wisp/**/*_test.go` | 0 |
| `c83aca3f`（批 3 件） | 4 | 0 | `124-r3-impl*.md` | 0 |

- 合计 **28** 枚 distinct 路径：`*_test.go` **28/28**、生产码 0 枚、`internal/{winsec,proc,risk}` 0 枚、
  禁改列（`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／`tools/d22scan/**`）0 枚、
  新建文件 0 枚（逐枚 `git cat-file -e <commit>^:<path>` 成功）。
- **正控**（同一把尺打真含 `.go` 的提交）：`fafe2b4a` => 9 件 / **9** 枚 `.go`。
  **负控**：`abd09226` => 3 件 / **0** 枚。=> 这把尺是活的，NOT 恒 0。
- 交付件 11 份我全部按 HEAD blob 现量：被禁码位 **0**、末字节 `0a`；行数逐枚
  r1 = 107 / 57 / 52、r2 = 126 / 130 / 154 / 49、r3 = 121 / 128 / 159 / 53 —— **全在 160 行帽内**。
- **本裁腿提交**：主题以 `票 124 终裁腿 124-v4` 开头；`git add` / `commit` 两处显式 pathspec（NOT `add -A`／`.`）；
  只 commit NOT push；写面只有 `docs/evidence/s1/124-v4-verdict.md` 与 `-part2.md` 两份。

## 7. 未验证 / 交回编排者

1. **`internal/panel` 那 1 枚两形都红**我**未复测**（不在本票写面，也没在我这台秤的分母里）；sec 5 的 68+9 只覆盖 `cmd/wisp`。
2. 三发变异我重打时**没有**做"逐发单独还原并核 sha256"的正式台账（变异都跑在仓外 `snap/`，工作树从头到尾未脏）；
   这一点上我不比 r2 强，也不比它弱——`git status --porcelain -- internal cmd` 全程为空。
3. 我**未**跑 `go vet` 的宿主交叉 `GOOS=linux`（本仓必然撞 `sherpa-onnx-go-linux` build constraints）；
   那一侧用**容器原生**读数替代，与 r3 交件的处置同形。
4. 那条"AC#6 移交 359 / 本票按 7 格结"的落账行、以及 sec 2 那两条口径更正，**要编排者写**；
   票面勾选状态、`Status:`、`-done` 后缀、台账、`docs/reports/**` 我**一字未动**。
5. 票面 AC#5 那句"`gofumpt v0.7.0 存在`"与今天手上的二进制（**v0.12.0**）不同；四批各报自己那一枚的处置我沿用（`ci.yml` 是 `@latest` 未钉版本）。
6. ★**收口期间树在动**：票 `359` 在 20:25–20:32 落地（`7835d67f`／`4ca369d5`／`0cea1927`），我 20:4x 收尾时现量到并在它上面补打了两发
   （基线 `75/47/0/0`、两发变异各 `75/46/1/0` 红名同一枚 `TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359`）。
   它只动 `internal/winsec/sameabsoluteness_359_other_test.go` 与 359 自己的件与票面，`resolve.go` blob 五处同一串 `e0b0e22f`；
   `git diff --name-only ef4c9702 0cea1927 -- internal/config internal/tools cmd/wisp` **空** => **本票三批的读数不受它影响**。
