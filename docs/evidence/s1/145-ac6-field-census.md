# 145 AC#6 取证腿（field census）——逐枚字段 → 具名断言用例对照表

- 日期：2026-11-07 13:1x–13:2x +08；性质＝**只读取证腿，零产码、零选形、零翻框**（勾与 Status 一字未动）。
- 起手锚 `acda6330d`，收手锚 `d687dd910`（共享树自前进，本程行号以现跑尺为准）。
- 写面：本件＋票 145 末尾 Progress log 一行；两处显式 pathspec、只 commit 不 push。
- 上游判语：`docs/evidence/s1/145-screen.md` §2（AC#6 定档①能做、只读取证程可派）。
- 措辞纪律：每枚字段要么具名用例、要么具名写「无，归票 N」；全件无「待补」一类含糊词。

## 0. 现跑提取（⛔ 抄票面、⛔ 抄旧证据件）

字段名全部由本程从 `internal/panel/composer.go` 生产码现跑抽出：

| 尺 | 读数 |
|---|---|
| `sed -n '57,92p' internal/panel/composer.go` | `Snapshot`＝**6 枚**直接字段：`Pending`/`Results`/`Composer`/`GeneratedAt`/`Instructions`/`Tasks` |
| `sed -n '271,309p' internal/panel/composer.go` | `ComposerState`＝**12 枚**字段：`Mode`/`Workspace`/`Attachments`/`AttachmentMIMEs`/`MaxAttachmentB`/`AttachmentReason`/`Git`/`Files`/`CurrentModel`/`ModelKnown`/`Credential`/`CredentialKnown` |
| 合同尺自报（本程 §3 那一发） | `Snapshot <-> PanelSnapshot: 6 JSON keys reconciled`＋`ComposerState <-> ComposerState: 12 JSON keys reconciled` —— 6+12 与上两行同读数，双向印证 |

## 1. 对照表（18 枚字段逐枚；「新」＝145 基线四字段之后落地的键，即 AC#6 的对象集）

每行「复跑」列即具名尺：`git grep` 命中处＝用例名＋文件行。

| 字段（json 键） | 新? | 有具名用例? | 用例名＋位置 | 复跑 |
|---|---|---|---|---|
| `Snapshot.Pending`（pending） | 基线 | 有 | `TestThePumpBuildsThePacketFromWhatTheHostHolds`（断言 :67–:90）＋`TestApprovalCardViewFromRealAssessor` | `git grep -n 'len(snap.Pending) != 1' -- internal/panel/pump_test.go` |
| `Snapshot.Results`（results） | 基线 | 有 | `TestThePumpBuildsThePacketFromWhatTheHostHolds`（断言 :102） | `git grep -n 'snap.Results\[1\].Done' -- internal/panel/pump_test.go` |
| `Snapshot.Composer`（composer） | 基线 | 有 | 四键钉 `TestThePumpBuildsThePacketFromWhatTheHostHolds`（:123）与 `TestPublishHandsTheBytesToTheAttachedExit`（:276–:284）；键集双向尺 `TestComposerContractTypesMatchFrontend` | `git grep -n 'composer,generatedAt,pending,results' -- internal/panel/pump_test.go` |
| `Snapshot.GeneratedAt`（generatedAt） | 基线 | 有 | `TestThePumpBuildsThePacketFromWhatTheHostHolds`（断言 :105）＋出口字节 :313 | `git grep -n 'want the injected clock' -- internal/panel/pump_test.go` |
| `Snapshot.Instructions`（instructions） | **新** | **有，归票 200** | `TestTheInstructionStatesDoNotShareOneWireShape`／`TestAPumpWithoutAnInstructionsReaderSendsNoKey`／`TestWithInstructionsKeepsTheR1CarrierContract`（均 instructions_200_test.go）；`TestRunPacketCarriesTheLoadedInstructionFiles`／`TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound`（cmd/wisp/instructions_200r2_test.go） | `git grep -n 'TestAPumpWithoutAnInstructionsReaderSendsNoKey' -- internal/panel cmd/wisp` |
| `Snapshot.Tasks`（tasks） | **新** | **有，归票 197** | `TestTheRosterReaderPutsSubagentsOnTheWire`／`TestAPumpWithoutARosterReaderSendsFourKeys`／`TestAStatusOutsideD43NeverReachesTheWire`／`TestTruncationFactsRideThePacket`／`TestDroppedStreamsAreNamedOnTheWire`（均 subagent_roster_197_test.go）；流侧五枚 `subagent_stream_197_test.go:35/:57/:146/:209/:275` | `git grep -n 'TestAPumpWithoutARosterReaderSendsFourKeys' -- internal/panel` |
| `ComposerState.Mode`（mode） | 基线 | 有 | `TestThePumpBuildsThePacketFromWhatTheHostHolds`（断言 :93）／`TestAPumpWithNoReadersSaysSoInsteadOfInventingState`（:131）／`TestAnUnreadableModeFromALiveReaderStillRendersUnknown`／`TestUnknownModeNeverRendersAsASafeOne`（composer_test.go:250） | `git grep -n 'composer.mode.current' -- internal/panel/pump_test.go` |
| `ComposerState.Workspace`（workspace） | 基线 | 有 | `TestThePumpBuildsThePacketFromWhatTheHostHolds`（断言 :96）／`TestWorkspaceViewRenderedForSnapshots`（workspace_test.go:155）／`TestWorkspaceViewCarriesEveryHalfOfTheAccount`（workspace_account_181r3_test.go:45）／`TestWorkspaceAbsentAndWorkspaceUnreadAreTwoSentences`（display_scrub_216_test.go:266） | `git grep -n 'TestWorkspaceAbsentAndWorkspaceUnreadAreTwoSentences' -- internal/panel` |
| `ComposerState.Attachments`（attachments） | 基线 | **答不出（恒空枚一）** | 见 §2 第 1 行 | （§2） |
| `ComposerState.AttachmentMIMEs`（acceptedAttachmentMimes） | 基线 | 有 | `TestPublishHandsTheBytesToTheAttachedExit`（断言 :299–:303：线上值＝broker 名册，不许 null） | `git grep -n 'acceptedAttachmentMimes' -- internal/panel/pump_test.go` |
| `ComposerState.MaxAttachmentB`（maxAttachmentBytes） | 基线 | 有 | `TestThePumpBuildsThePacketFromWhatTheHostHolds`（断言 :99）／`TestAPumpWithNoReadersSaysSoInsteadOfInventingState`（:140）／`TestPublishHandsTheBytesToTheAttachedExit`（:305） | `git grep -n 'maxAttachmentBytes =' -- internal/panel/pump_test.go` |
| `ComposerState.AttachmentReason`（attachmentError） | 基线 | **答不出（恒空枚二）** | 见 §2 第 2 行 | （§2） |
| `ComposerState.Git`（git） | **新** | **有** | `TestGitSectionTravelsInTheSnapshotPacket`（git_test.go:528，线上逐键）／`TestPumpWithoutGitReaderSaysUnreadable`（:576）／`TestGit189UncommittedAndDiffComeFromTheRealRepo`（git_189_test.go:151，真仓取数）／`TestGit189NoGitAndNoRepoReportAThirdStateNeverZero`（:315）／`TestGit189UncommittedJSONShape`（:403） | `git grep -n 'TestGitSectionTravelsInTheSnapshotPacket' -- internal/panel` |
| `ComposerState.Files`（files） | **新** | **线级无，归票 190**（reader 级有） | 线级：全仓零测试挂 `PumpSources.Files`、零 `"files"` 顶层键断言（本程现跑两把空尺，§2 第 3 行）；reader 级具名在盘：`TestReadGitTree190RealRepo`（tree_reader_190_test.go:34）／`TestTree190RootDoesNotExistUnreadable`（tree_190_ac234_test.go:12）／`TestTree190LinkOutsideRefused`（:25）／`TestTree190ReadSurfaceHasNoWriteCalls`（:44） | `git grep -n 'PumpSources{' -- '*_test.go' -A2 \| grep -i files`（空读数）＋`git grep -n 'TestReadGitTree190RealRepo' -- internal/panel` |
| `ComposerState.CurrentModel`（currentModel） | **新** | **有（本票对象用例）** | `TestComposerCurrentModelTravelsOnlyFromItsReader`（composer_test.go:916，头注 :905 逐字自认 ticket 145 AC#6；三臂＋出口字节钉 `"currentModel":"glm-5"` :947） | `git grep -n 'TestComposerCurrentModelTravelsOnlyFromItsReader' -- internal/panel` |
| `ComposerState.ModelKnown`（modelKnown） | **新** | **有（同一用例）** | `TestComposerCurrentModelTravelsOnlyFromItsReader`（断言 :923／:931／:936：known 只能由读口置位） | `git grep -n 'modelKnown = false for a model a reader named' -- internal/panel` |
| `ComposerState.Credential`（credentialState） | **新** | **有，归票 248** | `TestAC2CredentialSentinelAppearsInNoArtifact`（cmd/wisp/panel_config_248_test.go:154，键在钉 :173）＋`TestAC2SnapshotReportsRefWithoutBlobAsAPartState`（:510，值与 known 双向钉 :536–:548） | `git grep -n 'credentialState' -- cmd/wisp/panel_config_248_test.go` |
| `ComposerState.CredentialKnown`（credentialKnown） | **新** | **有（同一对用例）** | 同上（:539 true 臂／:547–:548 no-reader false 臂） | `git grep -n 'credentialKnown' -- cmd/wisp/panel_config_248_test.go` |

## 2. 「答不出」清单（逐枚具名，处置一律归编排者，本腿零选形、零拿掉动作）

1. **`ComposerState.Attachments`（json: attachments）＝恒空枚一。** 现跑尺：`git grep -n 'NewComposerState(mode, workspace, nil, maxAttachment)' -- internal/panel/pump.go` → `pump.go:277`，泵对装配根**硬传 nil**（构造器 `NewComposerState` composer.go:314–:315 把 nil 归一为空数组）。测试面的 Attachments 命中全是**构造器级／洗刷级**用例（`display_scrub_216_test.go:100`/:140、`attachments_test.go:267`），输入是测试自搓的行，**没有任何用例能断言「值来自真来源」**——因为生产通路根本没有第二处给它真值。⇒ AC#6 答句＝**答不出，归编排者裁**（拿掉该字段、或给泵接附件读口，两形都不归取证腿选）。
2. **`ComposerState.AttachmentReason`（json: attachmentError）＝恒空枚二。** 现跑尺：`git grep -rn 'attachmentError\|AttachmentReason' -- '*.go'` → **全仓唯一命中＝定义行 `composer.go:277`**，构造器不填（composer.go:316–:326 的字面量无此键）、泵不填、零测试引用。⇒ AC#6 答句＝**无具名用例，答不出，归编排者裁**。
3. **`ComposerState.Files`（json: files）的两层拆分（表内已列，此处留判语）**：reader 级（`readGitTree190` 的值来自真盘）**有具名用例**；线级（reader→packet 那一跳）**无具名用例**——两把空尺＝`git grep -n 'Files: *func\|\.Files()' -- '*_test.go'`（0 命中）与 `git grep -rn '"files"' -- '*_test.go'`（命中全是 git 维内层与 instructions 段内层，无顶层键）。⇒ AC#6 答句＝**线级无，归票 190**；且它正是 §3 合同尺在册红的那枚键（红名册照实记，处置归编排者／票 190）。

## 3. 门禁读数

- **合同尺一发**（环境照纪律：`unset GOTMPDIR TMPDIR TEMP TMP`＋`PATH` 前缀 `build`/`third_party/sherpa-onnx`）：
  `go test ./internal/panel/ -run 'TestComposerContractTypesMatchFrontend' -count=1`
  → **FAIL**：`composer_test.go:78: Go ComposerState emits [files] that interface ComposerState does not declare`（`Snapshot` 6 键＋`ComposerState` 12 键＋`ModeView` 3／`WorkspaceView` 6／`AttachmentRef` 9／`ResultChunk` 3 键全部 reconciled）。
  ⚠ **在册既有红**（`145-screen.md` 尺 4 同色，红在 `files` 键），**⛔ 本腿引入＝零**（本腿零 Go 改动，读数与筛查腿同键同色）。
- **gofumpt／vet**：本腿**零 Go 改动** ⇒ **无对象可扫**，未跑（跑了也没有被验物）。
- **AC#6 对象用例脸色**（本程未重跑，引用筛查腿尺 4/5 的绿读数并逐字登记：`TestComposerCurrentModelTravelsOnlyFromItsReader` 复跑绿）。

## 4. 本程没做的（如实）

- 未重跑各具名用例本体（只读取证腿的具名尺＝`git grep` 读数＋合同尺那一发，§1 每行可复跑）。
- 未翻任何 AC 勾、未动 Status、未碰 `frontend/**`、`design/**`、`internal/panel/**` 产码、`cmd/wisp/**`、thresholds/golden。
- 恒空两枚的处置（拿掉或补 reader）**没有选形**，按筛查腿判语原样交编排者。
