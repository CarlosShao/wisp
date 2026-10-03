# 220-r1 实现件：名册 blockedOnApproval 的两形 + L1 只读枚举口

## §0 锚点与自取读数

- 起手时刻：`Sat Oct 3 10:05:45 CST 2026`（`date` 现取）
- 分支 `dev`，起手 HEAD `git log -1 --format=%H` = `0583276ef0dc5e6d77404275da933387966104c5`
- 票面：`.scratch/wisp/issues/220-roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops.md`（完整读过；AC 框一枚未碰）
- 落点裁决照做：AC#2＝**甲形**（给 `Gate` 补只读枚举口，编排者具名落账台账 `A562`，先例 `A487`）；
  AC#1 未新增任何导出方法；AC#4 布尔未改成态、未新造 D43 名。

## §1 起手基线（改前整包名册，现跑）

尺：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 后
`GOFLAGS= go test ./internal/panel ./internal/agent/approval -count=1 -v`

基线四数（top-level `--- PASS` 158 ＋ 子测试 `--- PASS` 75 ＝ 233 枚）：

| 包 | PASS(top-level) | FAIL | SKIP | 包行 |
|---|---|---|---|---|
| `internal/panel` | 111 | 4 | 0 | `FAIL github.com/CarlosShao/wisp/internal/panel 4.676s` |
| `internal/agent/approval` | 47 | 0 | 1 | `ok github.com/CarlosShao/wisp/internal/agent/approval 0.526s` |
| 合计（含 75 枚子测试） | **233** | **4** | **1** | rc=1 |

改前今天绿的**用例名**（完整文件：`.scratch/wisp/probes/220/r1/prechange/baseline-panel-pass.txt`／
`baseline-approval-pass.txt`；下面逐名抄录，排序后）：

`internal/panel`（111 枚）：

```
TestAC1BothSettingsNamesAreWhitelistedAndRouted, TestAC1DispatchTableNamesEveryWhitelistedSettingMethod, TestAC1FieldTableCarriesNoAppr
ovalVocabulary, TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg, TestAC1RosteredSettingWithoutALegRefusesLoudly, TestAC1SelectorAndEmpty
ValueRefusals, TestAC1UnlistedMethodIsRefusedAndTheLegNeverRuns, TestAC2CredentialRouteNeverEchoesTheValueItWasGiven, TestAC2SettingsVie
wCarriesNoCredentialMaterial, TestAC2UnreadableConfigIsToldAsUnreadable, TestAC8ReceiptNamesTheTierAndNeverClaimsLiveSave, TestAModeRequ
estIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete, TestAPumpWithNoReadersSaysSoInsteadOfInventingState, TestAPumpWithoutARosterReade
rSendsFourKeys, TestAPumpWithoutAnInstructionsReaderSendsNoKey, TestAStatusOutsideD43NeverReachesTheWire, TestAStricterOrEqualModeReques
tNeedsNoConfirmLeg, TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached, TestAWideningRequestWithAnAttachedLegReachesTheOneWri
terAndRaisesNoSecondCard, TestAcceptedRequestInventsNoReplyLine, TestAcceptsAndStoresEachSupportedType, TestAnUnreadableModeFromALiveRea
derStillRendersUnknown, TestAnchorOnlyBundleIsNotBuiltAndFailsClosed, TestAnsweredPanelRoutesCarryNoApprovalDecision, TestApprovalCardVi
ewFromRealAssessor, TestApprovalCardViewMarksMissingReason, TestAttachedHandlerIsWhatRunsForItsOwnMethod, TestAttachmentPayloadCarriesTh
eBytes, TestAttributionFieldsReachTheHandlerUnchanged, TestBuiltinBundleIsCompleteOrAbsentNeverHalf, TestCheckAcceptsAConsistentBundle,
TestCheckIgnoresOutsideAndFragmentReferences, TestCheckRejectsEntryNamingAnUnembeddedAsset, TestComposerCurrentModelTravelsOnlyFromItsRe
ader, TestComposerEnvelopeAcceptsItsFourRequests, TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests, TestComposerRenderFixtureTel
lsTheTruth, TestComposerStateSurvivesPanelCloseAndReopen, TestDispatcherSpellsNoRouteLiteralOfItsOwn, TestDroppedStreamsAreNamedOnTheWire
, TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped, TestFrontendComposerRequestsMatchTheEnvelope, TestFrontendHasNoEmoji, TestFron
tendNeverNamesAnApprovalDecision, TestFullInboundMethodRosterIsClosed, TestFullInboundRosterRefusesAnEmptyDeclarationFile, TestGitDimens
ionHasNoModelCallableTool, TestGitSectionTravelsInTheSnapshotPacket, TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes, TestGrantWireS
hapesAreRefusedAtTheDoor, TestHardCeilingDropsWholeKeysAndNamesThem, TestJSONKeyDerivationAgreesWithEncodingJSON, TestMessageCarriesAtt
achmentRefsAndRefusals, TestModeRequestResolvesThroughRisksOwnParser, TestModeViewHasNoWriteSurface, TestNPlusOneSubagentsEachGetTheirOw
nStream, TestNoGitSwitchCapabilityInThePanelSurface, TestNoInboundEnvelopeCanBindAnApprovalVerdict, TestOverflowTruncatesEachStreamsOwnH
eadAndTail, TestPanelFrontendIsStateless, TestPanelHostIsAttachedAndNamesTheWindowHops, TestPanelSnapshotSurvivesWebviewRestart, TestPlan
tedComposerModeWriteGoesRed, TestPlantedGitToolShapesGoRed, TestPlantedGrantWiringGoesRedInASnapshot, TestPlantedRendererDoorShapesGoRed,
 TestPlantedUnregisteredInboundMethodNameGoesRed, TestPublishHandsTheBytesToTheAttachedExit, TestPublishWithoutAnExitReportsInsteadOfPre
tending, TestPumpWithoutGitReaderSaysUnreadable, TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce, TestReadGitDetachedSyntheticTre
e, TestReadGitForWorkspaceConsumesRewriteAccount, TestReadGitNonRepoSaysSoOutLoud, TestReadGitOnThisRepositoryIsSelfConsistent, TestRead
GitRelativePointerFile, TestReadGitSyntheticSymrefRepo, TestReadGitVanishedPointerIsUnreadable, TestReadGitWithoutAWorkspaceDoesNotProbe,
 TestReadGitWorktreePointerFile, TestRealGuardRefusesEveryAssemblableApprovalRouteName, TestRefusesUnsupportedAndMasqueradingInputsLoudl
y, TestRegistrationSilencesTheRosterRuler, TestResolveRefusesPathsOutsideTheBundle, TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll, Te
stRosterMismatchBackstopRefusesInsteadOfAccepting, TestRoutingDoesNotBypassTheNativeModeGate, TestSecondIdenticalAttachmentIsDeduplicated
NotOverwritten, TestSubagentStreamKeyIsTheOneSpelling, TestSubsetOnlyRosterDirectionIsBlindToAPlantedName, TestTamperedSourceIsRefusedBy
TheSameDoor, TestTheInboundHopAddsNoSwitchingCapability, TestTheInstructionStatesDoNotShareOneWireShape, TestTheModeLadderThisGateJudge
sAgainstIsTheDocumentedThree, TestThePumpBuildsThePacketFromWhatTheHostHolds, TestTheRendererHoldsExactlyOneDoorToTheHost, TestTheRoster
ReaderPutsSubagentsOnTheWire, TestTheStreamLogTruncatesInsteadOfMerging, TestTheWriteLegsOwnRefusalIsReturnedUntouched, TestTruncationFac
tsRideThePacket, TestWorkspaceSwitchAuditsAndAppliesACleanPath, TestWorkspaceSwitchDetectsAScopeThatMovedAnyway, TestWorkspaceSwitchKeep
sThePreviousScopeOnFailure, TestWorkspaceSwitchPropagatesC26sReparseRefusal, TestWorkspaceSwitchRefusesARewrittenAccountEvenIfTheResolve
rSaysOK, TestWorkspaceViewRenderedForSnapshots
```

`internal/agent/approval`（47 枚，另有 1 枚 SKIP＝包内既有）：

```
TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses, TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard, TestAVetoOnAnUnknownN
ameLeavesTheCardAllowableByItsOwnGrant, TestAnAliasCanNeverBuyAnAllow, TestAnAmbiguousVetoNameIsNotGuessedOnTheRefusalSide, TestAnUnansw
eredL2CardResolvesToRejectAtTheQueueDeadline, TestAnUnreachablePromptSurfaceIsRefusedImmediatelyNotWaitedOut, TestBatchAggregatesHomogene
ousL1Ops, TestBatchRefusalsBelowThresholdAndOnSensitiveVerdicts, TestDefaultQueueDeadlineIsFiniteAtThreeHundredSeconds, TestEveryAnswerRo
uteForAnUnknownCorrelationFailsFast, TestEveryRefusalRouteOnAnUnknownEntryStillRefuses, TestGrantIsSingleUseAndBoundToItsItem, TestHostUn
reachableAndFullQueueFailClosed, TestL1WindowTimeoutMeansExecute, TestL1WindowUnavailableVoiceChannelIsAnnounced, TestL1WindowUnreachable
UIFailsClosed, TestL1WindowVetoOnForeignCorrelationCannotCancel, TestL1WindowVetoedByEachLoadedChannel, TestL2NativeAllowExecutesThrough
TheBridge, TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s, TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord, TestLiveApp
rovalsOnNilQueueIsNil, TestLiveApprovalsReportsPendingRowsAndDoesNotConsumeThem, TestLiveApprovalsSharesNoReferenceSlotWithTheQueue, Test
LiveApprovalsSkipsRowsThatAreNotPending, TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis, TestR7SizedBatchEscalatesToL2Unaggregated,
TestR7SizedBatchThatNobodyAnsweredIsRejectedNotRun, TestReplayRedisplaysUnderAFreshGrant, TestReplySeamCarriesNoAllowWithoutAGrant, TestR
eplySeamWaitingStateNamesTheTwoD43Rows, TestTenOpsInOneToolCallGetOneConfirm, TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNam
ed, TestTicket224ForgedSessionAnswerRecordsNothing, TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow, TestTicket224SessionAnswe
rWithoutARecorderSaysTheScopeWasDropped, TestTicket242ForgedBindingCannotSpendAnotherItemsGrant, TestTicket242PanelItemReadFaceIsFullyDec
lared, TestTicket242PanelItemStaysGrantFree, TestTicket242QueuedItemsBindGrantsToTheirOwnDigest, TestTicket242SpendRejectsForgedBindingAn
dConsumesTheNonce, TestTicket242SpendRequiresTheExactBinding, TestUnadmittedTaskIsRefusedBeforeQueueing, TestVetoAfterStartYieldsAppliedS
tepsReport, TestVetoTravelsBackToTheLoopAsAFailure, TestWideningRuleReadsOneLineOffTheCard
```

### 基线那 4 枚红：具名登记＝起手就红，不是本腿造成的

四枚名（改前改后**同一集合**，见 §6）：
`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、
`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestC21DesignTokensFourWayAgree`。

现读出处（非本腿写面）：`git status --porcelain -- design frontend` 起手即有 **31 条改动**，其中含
`D design/assets/tokens.css`／`D design/assets/base.css`／`D design/assets/icons.js`／`D design/assets/theme.js`
等**工作树删除**；`TestC21DesignTokensFourWayAgree` 的失败行逐字＝
`tokens_fourway_test.go:441: read design/assets/tokens.css: open D:\...\design\assets\tokens.css: The system cannot find the path specified.`
⇒ 这四枚是**读 `design/**`／`frontend/**` 的两向对照尺**，正被同机别的在飞腿的临时改动挡住。
本腿**未碰**那两层一字（连读都没读），也未据此判自己失败。

## §2 改点（逐枚文件，带包名）

`internal/agent/approval`
- `gate.go`：`window` 结构补一枚 `taskID` 字段（**只加数据、不改判定**——注释里写明没有任何判定分支读写它）；
  `PendingWindow` 里构造 `w := &window{... taskID: d.TaskID ...}` 那一行随之补一个键。
  `gate.go:346-360`（开／关窗口）、`:377`（`Veto` 按名查）**语义一字未改**。
- `window_read.go`（新增）：`type L1Window struct{ CorrelationID, TaskID, Tool string }` ＋
  `func (g *Gate) LiveL1Windows() []L1Window`。头注释把边界写死：只读、只回"哪些 correlation/task 正在 L1 窗口里等"这一维、
  ⛔不带放行/否决能力、⛔不进 C17（并且**已 started 的调用不报**——`markStarted` 之后那一维属于 `g.running`，不是"在等"）。
- `ticket220_l1_window_read_test.go`（新增，`package approval` 白盒）：见 §5。

`internal/panel`
- `subagent_roster_197.go`：`blockedOnApproval` 的查键改成 **corr 与 taskID 双查** + 认 `#序号` 那一形。
  导出签名 `TaskRosterSectionFrom(state, pending)` **一字未动**（AC#1 不新增导出方法）；真正干活的是新增的
  **非导出** `taskRosterSectionFrom(state, pending, l1Windows)`，导出的那一枚以 `nil` 委托过去。
  新增两枚非导出小函数：`indexWaitingKey`（把一枚在等的 id 同时登记成"字面形"与"去掉队列改写后缀的原形"）、
  `splitClashSuffix`（**只**认队列自己写的那一形：最后一个 `#` 之后全是 ASCII 数字，见 `queue.go:152-154`）。
  文件头"not an approval outlet"那段随之改准：现在它读的是同一包的卡 **＋** 宿主枚举口给的窗口，仍**无任何通往答复的路由**。
- `pump.go`：新增 `type L1WindowWait struct{ CorrelationID, TaskID string }`（**无任何方法**）＋
  `PumpSources.L1Windows func() []L1WindowWait`（一枚 **reader**，沿用 `Instructions`／`Tasks` 那条既有规矩：
  nil＝这个宿主看不见这一维，**不等于**"没有窗口在等"）；`Snapshot()` 里把 `waits` 递给 join。
  ⛔ 窗口**没有**被塞进 `cards`／snapshot 的 pending 段——那一格是用户答复的入口（票 219 的按钮），
  把不可答复的窗口列进可答复段＝造出第二扇门的形状，也是票 146 对 `LiveApprovals` 既有判据的射程。
- `subagent_blocked_220_test.go`（新增）：见 §3/§4。

**未动**（冻结件，只做了射程判断、未引用内容）：`internal/panel/tokens_fourway_test.go`、
`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`。
射程判断（读行号得出，非内容引用）：前两枚分别射程＝C21 令牌四向对照 与 L2 授权边界/入站类型面，
本腿改的是**出向快照里一枚既有布尔的查键**与**一枚宿主侧 reader**，不落在那两枚尺的判据对象上；
第三枚射程＝票据持久化，与本腿写面不相交。`go build ./...` 与两包整包复跑对这些尺均未新增失败（§6）。

## §3 AC#3 两形正控：改前/改后**当场跑的**读数

### 形一：同一任务的第二张卡（`原id#序号`）

**改前**（用只依赖既有 API 的探针 `prechange220_probe_test.go`，现已随读数移到
`.scratch/wisp/probes/220/r1/prechange/`，⛔不在包里过夜）：

```
rc=1
--- FAIL: TestPRECHANGESecondCardOfOneTask (0.00s)
    prechange220_probe_test.go:53: PRE-CHANGE READING: the same task's second card is pending as "root-1" and the row still reads blockedOnApproval=false (want true)
```

同一次运行里第二半（`corr == taskID` 那一枚卡，缺省回落还成立的形状）**没有**报错行 ⇒ 缺陷确实只在第二张卡。

**改后**（常驻用例）：

```
--- PASS: TestASecondCardForTheSameTaskStillLightsItsRow (0.00s)
--- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt (0.00s)
    --- PASS: .../the_queue's_own_rewrite_lights  (0.00s)
    --- PASS: .../a_multi-digit_seq_lights  (0.00s)
    --- PASS: .../another_task's_rewrite_does_not_light  (0.00s)
    --- PASS: .../a_non-numeric_suffix_is_not_the_queue's_form  (0.00s)
    --- PASS: .../a_task_id_that_merely_prefixes_does_not_light  (0.00s)
    --- PASS: .../an_empty_correlation_id_names_no_row  (0.00s)
```

（名册那一格：只给 `child-1#9` 一张卡 ⇒ 该格 true；`child-1` 与 `child-1#9` 两张都在 ⇒ 仍 true；`root-1` 仍 false。）

### 形二：正在 L1 短窗口里等

**改前**：这一维**连枚举口都没有**，当场跑出的读数是编译级缺口（不是"应该会红"）：

```
internal\agent\approval\ticket220_l1_window_read_test.go:250:14: g.LiveL1Windows undefined (type *Gate has no field or method LiveL1Windows)
internal\agent\approval\ticket220_l1_window_read_test.go:332:24: undefined: L1Window
internal\panel\subagent_blocked_220_test.go:67:65: undefined: L1WindowWait
internal\panel\subagent_blocked_220_test.go:77:3: unknown field L1Windows in struct literal of type PumpSources
FAIL	github.com/CarlosShao/wisp/internal/panel [build failed]
FAIL	github.com/CarlosShao/wisp/internal/agent/approval [build failed]
```

配套尺（改前现跑）：`grep -rn 'g\.windows\|windows\[' internal/agent/approval --include=*.go | grep -v _test`
⇒ 只剩 `openWindow`/`closeWindow`/`Veto` 那几行自用，**没有任何枚举口**（票面"现量"那一行的形状复现）。

**改后**（两包各自正控）：

```
--- PASS: TestLiveL1WindowsSeesAWindowInTheLoop (0.00s)      # 真闸门 + 真开着的窗口 ⇒ 枚举回 {corr-220, task-220, fs.write}
--- PASS: TestAL1WindowInWaitingLightsItsRow (0.00s)         # 名册那一格：宿主喂 reader ⇒ true；不喂（nil）⇒ false
```

`internal/panel` 那一枚同时钉住**反向**：同一装配、把 `L1Windows` 留 nil ⇒ 那一格必须 false（"没读"不等于"没在等"）。

### 敏感性自查（编排者那句"判据换成反形还绿吗"）——三发突变，逐表现场跑，跑完即还原

| 突变 | 改的字面 | 现跑结果 |
|---|---|---|
| A 退回"只查字面 corr" | `if base, ok := splitClashSuffix(key); false && ok ...` | `--- FAIL: TestASecondCardForTheSameTaskStillLightsItsRow`、`--- FAIL: TestABlockedRowKeepsItsRunningDimension`（读数逐字：`row "child-1" blockedOnApproval=false，期望 true`） |
| B 换成"有任何卡就点亮全部行"（＝只查 taskID 那一侧的粗暴反形） | `BlockedOnApproval: row.TaskID != "" && len(waiting) > 0` | `--- FAIL: TestASecondCard...`、`--- FAIL: TestAL1Window...`、`--- FAIL: TestABlockedRow...`，且子测试三枚**负向**全红：`another_task's_rewrite_does_not_light`／`a_non-numeric_suffix_is_not_the_queue's_form`／`a_task_id_that_merely_prefixes_does_not_light` |
| C 让泵把枚举口丢掉 | `waits = waits[:0:len(waits)]` | `--- FAIL: TestAL1WindowInWaitingLightsItsRow`、`--- FAIL: TestABlockedRowKeepsItsRunningDimension` |

三发都用 `.scratch/wisp/probes/220/r1/prechange/roster.orig.txt`／`pump.orig.txt` 还原，还原后重跑＝§6 的终态。
⇒ 用例**对方向不敏感的形状过不了**：既不是"只加不减"的前向尺，也不是把整列点亮也能绿的钝尺。

## §4 AC#4 与票面"外部对照"那条设计约束

- `blockedOnApproval` **仍是 `bool`**（`subagent_roster_197.go:127` 那一枚字段一字未动，JSON 键 `blockedOnApproval` 未变），
  未新增 D43 名，`PLAN.md`／`docs/specs/**` 未碰一个字节。
- **附加的一维、不是替换**：常驻用例 `TestABlockedRowKeepsItsRunningDimension` 种的是**两形同时成立**
  （`child-1` 有第二张卡、`root-1` 有 L1 窗口），断言两行的 `blockedOnApproval` **都**为真**且** `status=="Thinking"`、
  `statusKnown==true` 原样带出；再加一条 wire 级否定：`tasks` 段里不许冒出 `blockedState`／`waitingOn` 这类新键
  （点亮的是既有布尔，不是新造一维）。这一枚在突变 A 与突变 C 下都会红（§3 表）。

## §5 AC#2 那枚枚举口"有牙但不越权"

`internal/agent/approval/ticket220_l1_window_read_test.go`（白盒，`package approval`）：

1. `TestReadingL1WindowsMovesNoGateState`：
   把枚举调用**夹在前后两次 `probeState` 之间**，逐维比对 `g.windows` 的**键集合与 `*window` 指针身份**、
   `g.running`、`g.admitted`、队列 `q.pending`（Corr 名册＋深度）、`g.grants` 计数（注入一枚计数用
   `GrantRecorder`）、UI 的 `Update` 次数——**八维零变化**才算过（`reflect.DeepEqual`，两枚指针切片也钉）。
   再连读两次断言返回**完全相同且有序**（名册字节要被台账哈希）。
2. 同枚用例末段：读了两遍之后 `g.Veto(...)` **仍然**打到那枚窗口、答复仍是 `AnswerVeto`
   ⇒ 枚举口**没有**消耗、解除或二次答复窗口（不是第二扇放行门）。
3. `TestLiveL1WindowsSeesAWindowInTheLoop`：种一枚**真**开着的 L1 窗口 ⇒ 枚举回
   `{corr-220 / task-220 / fs.write}` 一字不差；同时断言 `Queue().LiveApprovals()` **仍为空**
   ⇒ 枚举口没把窗口伪装成队列卡（票 219 答复段的边界）；窗口关闭后枚举回空。
4. `TestL1WindowIsPlainDataWithNoRouteToAnAnswer`：反射钉 `L1Window` **0 个方法**、字段**只能是 string**、
   字段名**恰好** `{CorrelationID, TaskID, Tool}`；并扫 `*Gate` 的方法面——除 `LiveL1Windows` 外
   任何方法都不许把这枚类型交出去。
5. `TestStartedCallsAreNotReportedAsWaiting`：`markStarted` 之后（＝已开跑）不许被报成"在等"
   ⇒ 反向的那枚谎（"需要批准吃掉了还在干活"的镜像）也被钉住。
- C17／入站面：本腿**没有**新增任何入站方法名；既有那枚闭集尺 `TestFullInboundMethodRosterIsClosed`
  （`inbound_roster_253_test.go`，闭集 6 枚、非 `panel.` 前缀 2 枚）在两包整包复跑里**仍 PASS**（§6 名册未丢名）。

## §6 门禁四数（自己动过的面）

| 尺 | 读数 |
|---|---|
| `GOFLAGS= go build ./...` | rc=0，无输出 |
| `gofumpt.exe -l internal/panel internal/agent/approval` | **空输出**（先 `-d` 看过再 `-w` 落到本腿 4 枚文件） |
| `./tools/d22scan/d22scan.exe` | `d22scan: clean - no D22 ban violations`，rc=0；覆盖面逐条照打：`bans #1-5 internal/=228, cmd/=37, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, frontend/=85, internal/=492, cmd/=89`（含 `_test.go`） |
| `go test ./internal/panel ./internal/agent/approval -count=1 -v`（终态四数） | **PASS=247**（top-level 115 ＋ 子测试 …；`internal/panel` 111→**115**、`internal/agent/approval` 47→**51**）／**FAIL=4**（**与基线同一集合**，§1 那四枚 design/frontend 尺）／**SKIP=1**（基线同枚）／rc=1；分包行：`FAIL internal/panel 3.943s`、`ok internal/agent/approval 0.748s` |
| 名册逐名对比（`comm` 排序后） | **丢名＝0**（两包都是），新增＝本腿那 4＋4 枚（见 §1 名册文件 `final-panel-pass.txt`／`final-approval-pass.txt`） |

未跑（照编排者令）：`./cmd/wisp`（255-r2 在飞）、`./internal/tools`／`./internal/risk`（252-v1 正在突变）；
本机无可用 staticcheck，未跑、未据其报绿；无 `t.Skip` 新增、无放宽任何断言。

## §7 归页面侧那一跳（只写这一节，⛔由 owner 带给他用的那枚 agent）

名册那一列**今天还差什么**（Go 侧射程之外，界面地界，本腿一字未碰、也不转述任何界面内容）：

1. Go 已经能把两维同时供到字节上：一行的 `status`（运行态那一维）**与** `blockedOnApproval`（在等批准那一维）。
   界面侧那一格**必须画成附加的一格**，不许做成"需要批准 > 还在跑"的单值优先级——
   票面"外部对照"记的就是别人家那条链把"还在干活"吃掉的形状。
2. 界面侧还看不见"在等的是**哪一张**卡"：本腿只点亮了布尔那一维，卡的可见性仍只有 L2 pending 段那一来源。
   "点进去看它自己的流"与"能不能在那页里回复它"是两件事，后者本仓从未许诺过（票 220 禁区；票 197 AC#5／Q-49 丙）。
3. ⚠ 上面这两条属界面判断，本腿**不做、不代达**，只落这一节。

## §8 我判不动的地方（具名，留给编排者）

1. **`cmd/wisp` 那一跳没接**（本腿被令⛔不许碰 `cmd/wisp/**`，连一条命令都没跑）：
   要让**生产**那一格真的为 L1 窗口点亮，装配根还欠两件事——
   (a) 给 `PumpSources.L1Windows` 供一枚 reader（形状＝`rt.gate.LiveL1Windows()` 映射成 `[]panel.L1WindowWait`，
   与 `panel_pump.go:62` 那枚 `liveVerdicts` 同形）；
   (b) **发布触发**：今天 packet 是"审批状态移动时"发的，而 L1 窗口开／关（`openWindow`/`closeWindow`）
   不在那些触发里——不补这一枚，窗口只在下一次发布才被看见（2-3s 那一族几乎不会自己赶上）。
   这两件都在 `cmd/wisp` 射程，⛔不是本腿能勾的。AC#5 的整包终态复跑（含 `./cmd/wisp`）归编排者。
2. **corr 与 taskID 之外的第三种配对**判不动：宿主若用一个**与任务无关**的 correlation id 发**卡**（不是窗口），
   `TaskRow` 里今天没有一枚 correlation 字段可供反查，那一格仍会读 false。要闭合就得给 `TaskRow` 加一维＋改 `cmd/wisp` 的映射，
   射程外，按票面"那是关于配对的陈述"如实留着。
3. **基线那 4 枚红**（§1）＝别腿在 `design/**`／`frontend/**` 的工作树删除所致，本腿**不改那两层一字**、也不据此判自己失败。
4. **AC 框一枚未碰**：票面 `- [ ]` 全部留给编排者翻勾，本腿不自述勾选。
5. 撤销口令（供编排者与 owner，不用在本件里执行）：**「220 改形乙」**。
