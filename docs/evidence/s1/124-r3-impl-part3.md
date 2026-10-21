# 票 124 落地腿 124-r3 交件 part3（sec 7..9：门禁 rc、红名册逐名、提交与正控）

> 前接 `docs/evidence/s1/124-r3-impl.md`（sec 0..4）与 `124-r3-impl-part2.md`（sec 5..6）。

## 7. 门禁（我自己现跑，`unset GOTMPDIR TMPDIR TEMP TMP` + `GOFLAGS=` 每一发）

| 门 | 命令 | rc | 末行逐字 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** | （0 行输出，无末行） |
| vet | `go vet ./cmd/wisp/` | **0** | （0 行输出） |
| test 1 | `go test -count=1 -v -timeout 30m ./cmd/wisp/` | **1** | `FAIL\tgithub.com/CarlosShao/wisp/cmd/wisp\t319.926s` |
| test 2 | `go test -count=1 ./cmd/wisp/` | **1** | `FAIL\tgithub.com/CarlosShao/wisp/cmd/wisp\t298.502s` |
| test 3 | `go test -count=1 -v -timeout 30m ./cmd/wisp/` | **1** | `FAIL\tgithub.com/CarlosShao/wisp/cmd/wisp\t276.382s` |
| d22scan | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; ... bans #1-5 internal/=244, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=610, ban #8 cmd/=145` |
| pathlen | `sh scripts/check-path-length-budget.sh` | **0** | `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |

**d22scan 正向对照逐字**（`runtests.sh` 那一步）：`top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77` => 门是活的。

**台账不降**：同一支 `sh scripts/d22scan.sh` 在**纯净快照**上也跑了一发（`GATE-D22SCAN-SNAP.log`）：
八 scope `244/42/116/25/14/116/610/145` ⇄ 工作树 `244/42/116/25/39/116/610/145`。
=> **唯一变化是 `ban #8 design/` 14 -> 39**，那 25 枚是**未跟踪的 `design/doubao/**`**
（`git ls-files design` = 14、`git status --porcelain -- design` 全是 `??`），
`design/**` 是机主外部 agent 的写面、**与本腿零因果**；其余七数（含 `cmd/` **145**）**一枚不降**。

### 三条 test 的门禁读数与死点（`A1100`：整包 rc=1 今天不能当判据）

按派单**只按名作差 + 逐发登记死点**，不把既有红算本腿头上：

| 发 | 顶层 RUN | PASS | FAIL | SKIP | 子 RUN | 子 PASS | 子 FAIL | 恰缺判决行的 RUN | 死点（逐发） | 时长 |
|---|---|---|---|---|---|---|---|---|---|---|
| (1) `-v -timeout 30m` | 159 | 158 | 0 | 0 | 58 | 58 | 0 | 1 | `TestAC14GoSideEvalPushReachesThePage` | 319.9s |
| (2) （无 `-v`） | — | — | — | — | — | — | — | — | 无名可读（未带 `-v`），只计数 | 298.5s |
| (3) `-v -timeout 30m` | 160 | 158 | 1 | 0 | 58 | 58 | 0 | 1 | `TestPanelThreadIsSTAAndExitsCleanly` | 276.4s |

- 三发时长 **276-320s**、死点**每发落在不同面板用例** => 与 `A1100` 记的
  "每发恰 1 枚 RUN 无判决行、每次落在不同面板用例、死点 277-343s" **逐条同形**。
- 两发带 `-v` 的**红名并集只有 1 枚**：`TestAC14GoSideEvalPushReachesThePage`。

### 那 1 枚红名：归**真窗族**、与本腿零因果（我自己现跑的证词）

```
$ git diff --name-only -- cmd/wisp | grep -c panel_resident_windows_test.go   => 0     (本腿未碰该文件)
$ 该文件依赖本腿 9 枚文件里任何 helper（grep sealableTempDir124|settingsDir248|newInboundDataDir33|
  empty313Dir|t257CleanMachine|t257npCleanMachine|reloadOneAppKey255）        => 0 命中
$ git show HEAD:cmd/wisp/panel_resident_windows_test.go | diff - <工作树那枚>  => 逐字节相同
```

**决定性一条：把本腿那 9 枚文件逐枚换回 HEAD blob，拿同一棵树单跑这一名**（`headcheck/` ＋ 补入被
`.gitignore` 掉的 `third_party/sherpa-onnx/*.dll`，否则宿主 `exit status 0xc0000135`）：

```
$ go test -count=3 -run '^TestAC14GoSideEvalPushReachesThePage$' ./cmd/wisp/   => FAIL（rc=1）
$ go test -count=1 -run '^TestAC14GoSideEvalPushReachesThePage$' ./cmd/wisp/
--- FAIL: TestAC14GoSideEvalPushReachesThePage (1.20s)
    panel_resident_windows_test.go:944: AC#14 nail 2 (Eval push hop), page's own words: title=""
    panel_resident_windows_test.go:946: Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK" ...
```

=> **在未改动的 HEAD 码上同址同文同样红**。台账 `pending-and-issues.md` 也早已具名登过它
（`:20471` A1100 格逐字 "`cmd/wisp` 整包 `rc=1` 有**既有红 1 枚** `TestAC14GoSideEvalPushReachesThePage`（与本票无关）"、
`:20542`、`:18173`、`:17943`）=> 归**真窗族（A875/314 族）**，NOT 本腿。

### gofmt / gofumpt

```
$ gofumpt --version                      => v0.12.0 (go1.27.1)     （绝对路径 /d/work/base/gopath/bin/gofumpt.exe）
$ gofumpt -l cmd/wisp   （工作树）        => 1 行：cmd/wisp/models.go
$ gofmt   -l cmd/wisp   （工作树）        => 1 行：cmd/wisp/models.go
$ gofmt -l /tmp/mh.go  （HEAD blob 副本） => 0 行  （HEAD 版是干净的）
```

NOTE **`cmd/wisp/models.go` 是既有形状，不是本腿**：`git status --porcelain -- cmd/wisp/models.go` **空**、
`git diff --stat` **空**、`git hash-object` 与 `git rev-parse HEAD:` **同一串** `4259a20f…`。
根因＝本机 `core.autocrlf=true`，那枚文件的**工作树副本**是 CRLF（334 个 CRLF、14213 字节），
而 **HEAD blob / 纯净快照都是纯 LF**（13879 字节，`gofmt -l` 0 行）。
=> 提交后 CI / 纯净快照上 gofmt 与 gofumpt **各 0 行**；工作树这一行是
`.gitattributes` 的 `*.go text eol=lf` 与本机 `core.autocrlf` 的既有冲突（**本腿的 9 枚文件全纯 LF、0 CRLF**，见 sec 4）。
`cmd/wisp` 整包 140 枚 `.go` 里**只有 models.go 这一枚带 CRLF**。

## 8. 红名册逐名（新增红 0）

**新增红 = 0 名**（判据＝名册集合）：
`POST-PLAIN` vs `PRE-PLAIN` 双向 `comm` **0 行**；`POST-LINK` vs `PRE-PLAIN` 双向 `comm` **0 行**；
变异发相对 `PRE-LINK` 的增项全部是**同族连带**（sec 5(1) 逐名五枚），不是"新增红"。

### 豁免 / 既有红逐名（各自注明归属）

**(a) 票 124 那族（本腿清零目标）：37 -> 0。豁免项 0 枚。**
`cmd/wisp` 名册里**没有**票 123 那三枚 300 s 腿（`TestL1WriteGoesThroughTheRealBlockWall` 一族在
`internal/tools/wiring_test.go:112/:182/:291`）；八发 `test timed out` 命中 **0**、300 s 量级腿 **0 枚**。

**(b) 两形都红 68 枚（顶层）+ 9 枚（子测试）＝归 DPAPI 族 / 既有门，本腿未修、未计入出口判据。**

顶层 68 逐名（= 裁腿 `124-v1` 那本，本腿一手量与它 set-equal）：

```
Test165r2AssembledBridgeReachesPlanPresent, Test165r2UnlistedWriteRingsItsOwnCardAfterAnApprovedPlan,
Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands, TestAC1AC2DispatchHopGate133,
TestAC1AlwaysBranchDoesNotRevertAHandEditedKey, TestAC2CredentialSentinelAppearsInNoArtifact,
TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef, TestAC4EveryLegIsNailedOrRuled,
TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit, TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card,
TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused, TestBookedSnapshotReachesHostBitForBit329,
TestCleanCheckoutBuilds_AC11, TestComposedGateBlocksAWriteForTwoSeconds,
TestCompositionRootClosesTheLoopTasksTaintScope, TestHostDispatchThroughTheAssembledBridge,
TestL1VetoNeedsAChannelTheHostReallyWired, TestMissingBlobFailsUnconfiguredNeverSilently,
TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole, TestNativeHostSeamRefusesAPanelSourcedAllow,
TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject, TestReplyListenerAllowsAnL2CardFromTheNativeSide,
TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel, TestRunBooksWithASnapshotOfItsLiveQueue,
TestRunPacketCarriesTheLoadedInstructionFiles, TestRunPacketCarriesTheSubagentItsRosterRowFed,
TestRunPacketLightsTheRosterRowOfItsLiveL1Window, TestRunPacketMarksTheRosterRowACardIsHolding,
TestRunPacketReportsTheStreamLogPastItsBound, TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound,
TestRunTextTaskFailNextIsClassified, TestRunTextTaskKeyResolvesInTheStore, TestRunTextTaskTextPathEndToEnd,
TestSecretEndToEndConfigRefResolvesAtRequestTime, TestSecretFailurePathsLogAndPrintNoPlaintext,
TestSecretFromStdinWritesNoIntermediateFile, TestSecretOverwriteIsAnnounced,
TestSecretPortableModeUsesTicket06Seam, TestSecretSameNameUnderThreeEnvsIsThreeBlobs,
TestSecretSetGetListRoundTrip, TestSecretUnsetRefusesWhileReferenced,
TestSnapshotWorkspaceSectionReportsTheNarrowing, TestTaskOutputPointerNotice174r8SpeaksTheRealReason,
TestTaskOutputPointerNotice174r8StaysQuiet, TestThePumpIsDrivenNotJustAssembled,
TestTicket101ManualSwitchSurvivesRestart, TestTicket101ModeSwitchUsesTheRealL2Gate,
TestTicket101SessionGrantDoesNotCrossRestart, TestTicket101UnreadableModeFailsLoudlyAndStrict,
TestTicket101UntouchedConfigRestartsAtDefault, TestTicket223FailureSentencesAreDistinct,
TestTicket223HandEditedFsLooseningCostsAnL2Card, TestTicket223ModeLooseningChangesTheRunningModeAfterAllow,
TestTicket223PanelInboundSaysHotReloadIsDisabled, TestTicket223RefusedLooseningKeepsOldValues,
TestTicket223RestartTierSaysItWillNotApply, TestTicket223RunArmsTheReloadTick,
TestTicket223TighteningRaisesNoCard, TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking,
TestTicket224ProductionSessionDoesNotSurviveRestart, TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun,
TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere, TestTicket255ReceiptOmitsPanelFromTheImmediateSentence,
TestTicket255ReceiptSentenceAssemblyIsFiltered, TestTicket255ReceiptStillNamesTheLiveReadSection,
TestTicket255RestartTierKeysAreBackedByATest, TestUnansweredL2CardTimesOutIntoRejectNeverExecution,
TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask, TestZZ315A1RootLoopFirstBridgeCall,
TestZZ315A1SiblingChildFirstBridgeCall
```

机制归类（`harness/classify.py`）：**DPAPI 族 64 枚**（块内含
`secret: DPAPI is only available on Windows`）＋ **既有门/环境 4 枚**：
`TestAC1AC2DispatchHopGate133`（票 133 门）、`TestAC4EveryLegIsNailedOrRuled`（票 131 门）、
`TestCleanCheckoutBuilds_AC11`（`git rev-parse --show-toplevel: exit status 128`）、
`TestSecretOverwriteIsAnnounced`（DPAPI 下游：`first set: 1`）。

子测试 9 枚（两形同形，逐名）：

```
TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes
TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands
TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)
TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)
TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only
TestTicket223FailureSentencesAreDistinct/schema未知键
TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面
TestTicket223FailureSentencesAreDistinct/缺失
TestTicket223FailureSentencesAreDistinct/语法错
```

**(c) 宿主 Windows 档**：三发红名并集 **1 枚**＝`TestAC14GoSideEvalPushReachesThePage`（真窗族，sec 7 已证同址同文在 HEAD 红）。

**(d) `internal/config` 那 2 枚 POSIX-only 摇摆**（`TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent` /
`TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed`）**不在本包名册**、本腿未碰、未计入任何出口判据。

> 本件续 `124-r3-impl-part4.md`（sec 9：提交号、逐枚 `.go`、越界与正控）。
