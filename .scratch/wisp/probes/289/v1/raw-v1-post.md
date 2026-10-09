time=2026-10-09T14:34:30.356+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestApprovalCardViewFromRealAssessor
    approval_test.go:54: real assessor verdict: level=L2 rules=[R1 R8] reason="R1: 工具声明为下界（L1）; R8: 不可逆操作（永久删除）"
--- PASS: TestApprovalCardViewFromRealAssessor (0.00s)
=== RUN   TestApprovalCardViewMarksMissingReason
--- PASS: TestApprovalCardViewMarksMissingReason (0.00s)
=== RUN   TestApprovalCardViewJSONKeysMatchFrontendTypes
    approval_test.go:134: ApprovalCardView <-> ApprovalCardView: 10 JSON keys reconciled
    approval_test.go:134: ResultChunk <-> ResultChunkView: 3 JSON keys reconciled
    approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
    approval_test.go:134: Snapshot <-> PanelSnapshot: 6 JSON keys reconciled
--- FAIL: TestApprovalCardViewJSONKeysMatchFrontendTypes (0.00s)
=== RUN   TestPanelSnapshotSurvivesWebviewRestart
    approval_test.go:167: snapshot round-tripped through 944 bytes of JSON with 1 pending card(s)
--- PASS: TestPanelSnapshotSurvivesWebviewRestart (0.00s)
=== RUN   TestAnchorOnlyBundleIsNotBuiltAndFailsClosed
--- PASS: TestAnchorOnlyBundleIsNotBuiltAndFailsClosed (0.00s)
=== RUN   TestCheckRejectsEntryNamingAnUnembeddedAsset
--- PASS: TestCheckRejectsEntryNamingAnUnembeddedAsset (0.00s)
=== RUN   TestCheckAcceptsAConsistentBundle
--- PASS: TestCheckAcceptsAConsistentBundle (0.00s)
=== RUN   TestCheckIgnoresOutsideAndFragmentReferences
--- PASS: TestCheckIgnoresOutsideAndFragmentReferences (0.00s)
=== RUN   TestBuiltinBundleIsCompleteOrAbsentNeverHalf
    assets_test.go:144: verdict: bundle complete - entry 1044 bytes, 2 asset refs, 4 embedded files
--- PASS: TestBuiltinBundleIsCompleteOrAbsentNeverHalf (0.00s)
=== RUN   TestResolveRefusesPathsOutsideTheBundle
--- PASS: TestResolveRefusesPathsOutsideTheBundle (0.00s)
=== RUN   TestAcceptsAndStoresEachSupportedType
=== RUN   TestAcceptsAndStoresEachSupportedType/shot.png
    attachments_test.go:125: stored: artifact=attachment-e1f5c234e80b19af.png bytes=16 mime=image/png
=== RUN   TestAcceptsAndStoresEachSupportedType/photo.jpg
    attachments_test.go:125: stored: artifact=attachment-9f6c9d639ad1155a.jpg bytes=12 mime=image/jpeg
=== RUN   TestAcceptsAndStoresEachSupportedType/clip.mp4
    attachments_test.go:125: stored: artifact=attachment-a097c197ed8b9340.mp4 bytes=20 mime=video/mp4
=== RUN   TestAcceptsAndStoresEachSupportedType/anim.gif
    attachments_test.go:125: stored: artifact=attachment-01075758a3f9dded.gif bytes=18 mime=image/gif
=== RUN   TestAcceptsAndStoresEachSupportedType/pic.webp
    attachments_test.go:125: stored: artifact=attachment-583434d1fd31cfcc.webp bytes=32 mime=image/webp
--- PASS: TestAcceptsAndStoresEachSupportedType (0.02s)
    --- PASS: TestAcceptsAndStoresEachSupportedType/shot.png (0.00s)
    --- PASS: TestAcceptsAndStoresEachSupportedType/photo.jpg (0.00s)
    --- PASS: TestAcceptsAndStoresEachSupportedType/clip.mp4 (0.00s)
    --- PASS: TestAcceptsAndStoresEachSupportedType/anim.gif (0.00s)
    --- PASS: TestAcceptsAndStoresEachSupportedType/pic.webp (0.00s)
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/empty_file_is_refused,_not_dropped
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/exe_disguised_as_png_is_refused_and_named_as_an_executable
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/unknown_type_without_a_declared_mime_is_refused
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/lying_declared_mime_is_refused_even_when_the_bytes_are_a_png
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/path-injection_file_name_is_refused_and_nothing_is_written
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/2_GB_is_refused_on_its_declared_size_without_reading_a_byte
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/a_source_lying_about_its_size_is_refused_mid-read,_not_truncated
=== RUN   TestRefusesUnsupportedAndMasqueradingInputsLoudly/nil_source_and_nil_sink_are_refused_at_construction,_not_at_send
--- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly (0.01s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/empty_file_is_refused,_not_dropped (0.00s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/exe_disguised_as_png_is_refused_and_named_as_an_executable (0.00s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/unknown_type_without_a_declared_mime_is_refused (0.00s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/lying_declared_mime_is_refused_even_when_the_bytes_are_a_png (0.00s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/path-injection_file_name_is_refused_and_nothing_is_written (0.00s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/2_GB_is_refused_on_its_declared_size_without_reading_a_byte (0.00s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/a_source_lying_about_its_size_is_refused_mid-read,_not_truncated (0.00s)
    --- PASS: TestRefusesUnsupportedAndMasqueradingInputsLoudly/nil_source_and_nil_sink_are_refused_at_construction,_not_at_send (0.00s)
=== RUN   TestSecondIdenticalAttachmentIsDeduplicatedNotOverwritten
--- PASS: TestSecondIdenticalAttachmentIsDeduplicatedNotOverwritten (0.00s)
=== RUN   TestMessageCarriesAttachmentRefsAndRefusals
    attachments_test.go:302: agent-facing text:
        看看这个
        [附件] name=clip.mp4 mime=video/mp4 kind=video bytes=20 path=C:\Users\swq\AppData\Local\Temp\TestMessageCarriesAttachmentRefsAndRefusals2560043113\001\attachment-a097c197ed8b9340.mp4
        [附件未送达] evil.png: "evil.png" 的类型（可执行文件（MS-DOS/PE 头 "MZ"））不受支持；本产品的附件只接受 image/png / image/jpeg / image/gif / image/webp / video/mp4 / video/quicktime
        （这些附件没有进入本轮，请用户确认后重试或改用受支持的类型）
--- PASS: TestMessageCarriesAttachmentRefsAndRefusals (0.00s)
=== RUN   TestAttachmentPayloadCarriesTheBytes
=== RUN   TestAttachmentPayloadCarriesTheBytes/a_real_png_arrives_whole_through_base64
=== RUN   TestAttachmentPayloadCarriesTheBytes/undecodable_base64_is_refused,_not_stored_as_nothing
=== RUN   TestAttachmentPayloadCarriesTheBytes/a_truncated_payload_is_refused,_not_stored_short
=== RUN   TestAttachmentPayloadCarriesTheBytes/an_empty_payload_is_refused
--- PASS: TestAttachmentPayloadCarriesTheBytes (0.00s)
    --- PASS: TestAttachmentPayloadCarriesTheBytes/a_real_png_arrives_whole_through_base64 (0.00s)
    --- PASS: TestAttachmentPayloadCarriesTheBytes/undecodable_base64_is_refused,_not_stored_as_nothing (0.00s)
    --- PASS: TestAttachmentPayloadCarriesTheBytes/a_truncated_payload_is_refused,_not_stored_short (0.00s)
    --- PASS: TestAttachmentPayloadCarriesTheBytes/an_empty_payload_is_refused (0.00s)
=== RUN   TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8
    backpressure_bound_35r8_test.go:278: 35r8 fan-out reading: 1x text=35072 rows=64 names=336 heap+244888B | 4x text=35072 rows=64 names=1536 heap+472104B | cap text=36864 rows=64
--- PASS: TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8 (2.62s)
=== RUN   TestStreamLogFloodBelowKeyBoundIsNotBounded35r8
    backpressure_bound_35r8_test.go:309: flood under a blocked consumer with no fan-out (3 streams, so the key bound never engages): 1x pinned 1200000 text runes over 3 rows (live heap +1202080 bytes), 4x pinned 4800000 runes over 3 rows (live heap +4825288 bytes), against the cap this package can state (36864 runes over 64 rows). Retention grew 4.0x for a 4x flood and equals the flood exactly, so ticket 35 :49's "no unbounded memory (heap cap asserted)" has no counterpart in internal/panel below the key bound. pinned text 4800000 runes > the cap 36864; live heap + 4825288 bytes > the cap 1196032
--- FAIL: TestStreamLogFloodBelowKeyBoundIsNotBounded35r8 (1.03s)
=== RUN   TestStreamLogDroppedNamingLedgerIsNotBounded35r8
    backpressure_bound_35r8_test.go:336: blocked consumer, fan-out past the ceiling: 1x opened 4000 streams and pins 3936 names (50058 runes, live heap +155072 bytes), 4x opened 16000 and pins 15936 names (211994 runes, live heap +563840 bytes) - rows stayed at the ceiling 64 in both, so the memory the flood bought is the naming ledger itself. Its cap is the event count, not the bound. pinned dropped-name runes 211994 > the cap 36864
--- FAIL: TestStreamLogDroppedNamingLedgerIsNotBounded35r8 (0.58s)
=== RUN   TestFloodRulerFiresOnAnUnboundedSink35r8
    backpressure_bound_35r8_test.go:360: 35r8 positive control fired as required: rows 100 > the hard ceiling 64; pinned text 800000 runes > the cap 36864 (heap +825856 bytes)
--- PASS: TestFloodRulerFiresOnAnUnboundedSink35r8 (0.61s)
=== RUN   TestComposerEnvelopeAcceptsItsFourRequests
--- PASS: TestComposerEnvelopeAcceptsItsFourRequests (0.00s)
=== RUN   TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests
=== RUN   TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/a_request_that_names_an_approval_decision_is_not_a_composer_method
=== RUN   TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/a_missing_requestId_is_refused
=== RUN   TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/a_forged_or_absent_source_is_refused
=== RUN   TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/undecodable_and_non-object_payloads_are_refused,_not_defaulted
=== RUN   TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/the_refusal_text_keeps_the_request_identifiable
--- PASS: TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests (0.00s)
    --- PASS: TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/a_request_that_names_an_approval_decision_is_not_a_composer_method (0.00s)
    --- PASS: TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/a_missing_requestId_is_refused (0.00s)
    --- PASS: TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/a_forged_or_absent_source_is_refused (0.00s)
    --- PASS: TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/undecodable_and_non-object_payloads_are_refused,_not_defaulted (0.00s)
    --- PASS: TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests/the_refusal_text_keeps_the_request_identifiable (0.00s)
=== RUN   TestFrontendComposerRequestsMatchTheEnvelope
--- PASS: TestFrontendComposerRequestsMatchTheEnvelope (0.00s)
=== RUN   TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce
--- PASS: TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce (0.00s)
=== RUN   TestAcceptedRequestInventsNoReplyLine
--- PASS: TestAcceptedRequestInventsNoReplyLine (0.00s)
=== RUN   TestRoutingDoesNotBypassTheNativeModeGate
--- PASS: TestRoutingDoesNotBypassTheNativeModeGate (0.00s)
=== RUN   TestUnlistedMethodNameIsRefusedAndAudited
--- PASS: TestUnlistedMethodNameIsRefusedAndAudited (0.00s)
=== RUN   TestRosterMismatchBackstopRefusesInsteadOfAccepting
--- PASS: TestRosterMismatchBackstopRefusesInsteadOfAccepting (0.00s)
=== RUN   TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped
=== RUN   TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped/workspace
=== RUN   TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped/attachment
=== RUN   TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped/message
--- PASS: TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped (0.00s)
    --- PASS: TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped/workspace (0.00s)
    --- PASS: TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped/attachment (0.00s)
    --- PASS: TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped/message (0.00s)
=== RUN   TestAttachedHandlerIsWhatRunsForItsOwnMethod
--- PASS: TestAttachedHandlerIsWhatRunsForItsOwnMethod (0.00s)
=== RUN   TestTamperedSourceIsRefusedByTheSameDoor
--- PASS: TestTamperedSourceIsRefusedByTheSameDoor (0.00s)
=== RUN   TestAttributionFieldsReachTheHandlerUnchanged
--- PASS: TestAttributionFieldsReachTheHandlerUnchanged (0.00s)
=== RUN   TestDispatcherSpellsNoRouteLiteralOfItsOwn
--- PASS: TestDispatcherSpellsNoRouteLiteralOfItsOwn (0.00s)
=== RUN   TestPanelHostIsAttachedAndNamesTheWindowHops
    composer_dispatch_test.go:438: native-host capability hits in D:\work\workspace\projects plans\Wisp: 5 -> [cmd/wisp/panel_host_windows.go:64:imports github.com/jchv/go-webview2 cmd/wisp/panel_resident_windows.go:73:imports github.com/jchv/go-webview2 go.mod:19:dependency "webview" go.sum:10:dependency "webview" go.sum:9:dependency "webview"] (production=true dependency=true)
=== RUN   TestPanelHostIsAttachedAndNamesTheWindowHops/positive-control
    composer_dispatch_test.go:468: carrier A (fake host source) -> 5 hit(s) [internal/panel/host_windows.go:3:identifier CoreWebView2 internal/panel/host_windows.go:5:identifier CoreWebView2 internal/panel/host_windows.go:5:identifier WebMessage internal/panel/host_windows.go:7:identifier CoreWebView2 internal/panel/host_windows.go:7:identifier WebMessage]
    composer_dispatch_test.go:480: carrier B (webview dependency) -> 2 hit(s) [go.mod:3:dependency "webview" go.sum:1:dependency "webview"]
    composer_dispatch_test.go:504: carrier C (CLI seam, no host) -> 0 hits, as required
=== RUN   TestPanelHostIsAttachedAndNamesTheWindowHops/reverse-positive-control
    composer_dispatch_test.go:533: carrier D (host identifiers neutralised) -> 0 hit(s) [] (production=false dependency=false)
--- PASS: TestPanelHostIsAttachedAndNamesTheWindowHops (0.15s)
    --- PASS: TestPanelHostIsAttachedAndNamesTheWindowHops/positive-control (0.01s)
    --- PASS: TestPanelHostIsAttachedAndNamesTheWindowHops/reverse-positive-control (0.00s)
=== RUN   TestTheInboundHopAddsNoSwitchingCapability
--- PASS: TestTheInboundHopAddsNoSwitchingCapability (0.00s)
=== RUN   TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree
--- PASS: TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree (0.00s)
=== RUN   TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached
--- PASS: TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached (0.00s)
=== RUN   TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard
--- PASS: TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard (0.00s)
=== RUN   TestAStricterOrEqualModeRequestNeedsNoConfirmLeg
--- PASS: TestAStricterOrEqualModeRequestNeedsNoConfirmLeg (0.00s)
=== RUN   TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete
=== RUN   TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/no_write_leg_at_all
=== RUN   TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/current档_unreadable
=== RUN   TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/unknown_spelling
=== RUN   TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/empty_target
=== RUN   TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/another_method
--- PASS: TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete (0.00s)
    --- PASS: TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/no_write_leg_at_all (0.00s)
    --- PASS: TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/current档_unreadable (0.00s)
    --- PASS: TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/unknown_spelling (0.00s)
    --- PASS: TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/empty_target (0.00s)
    --- PASS: TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete/another_method (0.00s)
=== RUN   TestTheWriteLegsOwnRefusalIsReturnedUntouched
--- PASS: TestTheWriteLegsOwnRefusalIsReturnedUntouched (0.00s)
=== RUN   TestComposerContractTypesMatchFrontend
    composer_test.go:74: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
    composer_test.go:79: Snapshot <-> PanelSnapshot: 6 JSON keys reconciled
    composer_test.go:74: Go ComposerState emits [git currentModel modelKnown credentialState credentialKnown] that interface ComposerState does not declare
    composer_test.go:79: ComposerState <-> ComposerState: 11 JSON keys reconciled
    composer_test.go:79: ModeView <-> ComposerMode: 3 JSON keys reconciled
    composer_test.go:79: WorkspaceView <-> ComposerWorkspace: 6 JSON keys reconciled
    composer_test.go:79: AttachmentRef <-> ComposerAttachment: 9 JSON keys reconciled
    composer_test.go:79: ResultChunk <-> ResultChunkView: 3 JSON keys reconciled
--- FAIL: TestComposerContractTypesMatchFrontend (0.00s)
=== RUN   TestPlantedComposerModeWriteGoesRed
=== RUN   TestPlantedComposerModeWriteGoesRed/approval.decide_planted_in_composer_trips_ban_#6
    composer_test.go:132: red as required: 1 line(s) matched, first=composer.tsx:4: <button onClick={() => bridge?.postMessage(JSON.stringify({ method: "approval.decide", target: "mode", to: "auto_approve" }))}>
=== RUN   TestPlantedComposerModeWriteGoesRed/panel.mode.set_planted_in_composer_trips_the_ticket_92_gate
    composer_test.go:132: red as required: 2 line(s) matched, first=composer.tsx:1: export function setMode(next: string): void {
=== RUN   TestPlantedComposerModeWriteGoesRed/door_2_ignores_comment_prose_but_not_code_behind_a_comment_opener
--- PASS: TestPlantedComposerModeWriteGoesRed (0.15s)
    --- PASS: TestPlantedComposerModeWriteGoesRed/approval.decide_planted_in_composer_trips_ban_#6 (0.00s)
    --- PASS: TestPlantedComposerModeWriteGoesRed/panel.mode.set_planted_in_composer_trips_the_ticket_92_gate (0.00s)
    --- PASS: TestPlantedComposerModeWriteGoesRed/door_2_ignores_comment_prose_but_not_code_behind_a_comment_opener (0.00s)
=== RUN   TestModeViewHasNoWriteSurface
--- PASS: TestModeViewHasNoWriteSurface (0.00s)
=== RUN   TestModeRequestResolvesThroughRisksOwnParser
--- PASS: TestModeRequestResolvesThroughRisksOwnParser (0.00s)
=== RUN   TestComposerStateSurvivesPanelCloseAndReopen
    composer_test.go:240: composer state round-tripped through 1198 bytes: mode="ask_high_risk" workspace="D:\\work\\notes" attachments=1
--- PASS: TestComposerStateSurvivesPanelCloseAndReopen (0.00s)
=== RUN   TestUnknownModeNeverRendersAsASafeOne
--- PASS: TestUnknownModeNeverRendersAsASafeOne (0.00s)
=== RUN   TestNoGitSwitchCapabilityInThePanelSurface
    composer_test.go:324: AC#7 negative criterion held: no git switch entry point in frontend/ or internal/panel/ (production files)
--- PASS: TestNoGitSwitchCapabilityInThePanelSurface (0.30s)
=== RUN   TestTheRendererHoldsExactlyOneDoorToTheHost
    composer_test.go:524: renderer door held: 64 files scanned, 2 host call sites (all in src/lib/panel.ts), 5 panel.* route literals
--- PASS: TestTheRendererHoldsExactlyOneDoorToTheHost (0.01s)
=== RUN   TestPlantedRendererDoorShapesGoRed
    composer_test.go:582: red as required: 2 outside call sites, 1 assembled routes, 1 unknown literals, door 2 named 1 line(s)
--- PASS: TestPlantedRendererDoorShapesGoRed (0.01s)
=== RUN   TestComposerRenderFixtureTellsTheTruth
    composer_test.go:728: render fixture verified across 3 painted states
--- PASS: TestComposerRenderFixtureTellsTheTruth (0.00s)
=== RUN   TestComposerCurrentModelTravelsOnlyFromItsReader
--- PASS: TestComposerCurrentModelTravelsOnlyFromItsReader (0.00s)
=== RUN   TestAC1BothSettingsNamesAreWhitelistedAndRouted
=== RUN   TestAC1BothSettingsNamesAreWhitelistedAndRouted/config.get_reaches_the_leg_once_and_its_answer_comes_back
=== RUN   TestAC1BothSettingsNamesAreWhitelistedAndRouted/config.set_reaches_the_leg_once_with_the_parsed_selectors
--- PASS: TestAC1BothSettingsNamesAreWhitelistedAndRouted (0.00s)
    --- PASS: TestAC1BothSettingsNamesAreWhitelistedAndRouted/config.get_reaches_the_leg_once_and_its_answer_comes_back (0.00s)
    --- PASS: TestAC1BothSettingsNamesAreWhitelistedAndRouted/config.set_reaches_the_leg_once_with_the_parsed_selectors (0.00s)
=== RUN   TestAC1UnlistedMethodIsRefusedAndTheLegNeverRuns
--- PASS: TestAC1UnlistedMethodIsRefusedAndTheLegNeverRuns (0.00s)
=== RUN   TestAC1RosteredSettingWithoutALegRefusesLoudly
--- PASS: TestAC1RosteredSettingWithoutALegRefusesLoudly (0.00s)
=== RUN   TestAC1DispatchTableNamesEveryWhitelistedSettingMethod
--- PASS: TestAC1DispatchTableNamesEveryWhitelistedSettingMethod (0.00s)
=== RUN   TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg
--- PASS: TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg (0.00s)
=== RUN   TestAC1FieldTableCarriesNoApprovalVocabulary
--- PASS: TestAC1FieldTableCarriesNoApprovalVocabulary (0.00s)
=== RUN   TestAC1SelectorAndEmptyValueRefusals
--- PASS: TestAC1SelectorAndEmptyValueRefusals (0.00s)
=== RUN   TestAC8ReceiptNamesTheTierAndNeverClaimsLiveSave
--- PASS: TestAC8ReceiptNamesTheTierAndNeverClaimsLiveSave (0.00s)
=== RUN   TestAC2CredentialRouteNeverEchoesTheValueItWasGiven
--- PASS: TestAC2CredentialRouteNeverEchoesTheValueItWasGiven (0.00s)
=== RUN   TestAC2UnreadableConfigIsToldAsUnreadable
--- PASS: TestAC2UnreadableConfigIsToldAsUnreadable (0.00s)
=== RUN   TestAC2SettingsViewCarriesNoCredentialMaterial
--- PASS: TestAC2SettingsViewCarriesNoCredentialMaterial (0.00s)
=== RUN   TestPanelFrontendIsStateless
    frontend_hygiene_test.go:193: scanned 64 frontend/src files for 7 persistence APIs: 0 hits
--- PASS: TestPanelFrontendIsStateless (0.14s)
=== RUN   TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme
    frontend_hygiene_test.go:216: a second style source appeared in code the bundle actually ships:
          frontend/src/components/harness/right-rail.tsx:89: diff: [{ sign: "+", text: "  --page: #fafafb;  /* beautiful-ui 主题收编 */" }],
        Add or reference a C21 token in design/assets/tokens.css instead.
    frontend_hygiene_test.go:219: 57 files reachable from main.tsx carry colour literals only in the generated theme
--- FAIL: TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme (0.05s)
=== RUN   TestVendoredDemoComponentsAreNotMounted
    frontend_hygiene_test.go:242: vendored ai-native: 18 mounted (button.tsx chip.tsx code-block.tsx context-cards.tsx diff-table.tsx entity-chip.tsx loading-state.tsx progress-ring.tsx prompt-bar.tsx segmented-control.tsx shimmer.tsx status-pill.tsx streaming-text.tsx switch.tsx text-row.tsx thinking.tsx tool-chips.tsx value-pill.tsx), 3 unmounted upstream demos (approval-card.tsx stream-text.tsx task-rows.tsx)
--- PASS: TestVendoredDemoComponentsAreNotMounted (0.03s)
=== RUN   TestFrontendHasNoEmoji
    frontend_hygiene_test.go:278: ban #8 self-armed: 72 frontend files scanned, 0 emoji-range characters
--- PASS: TestFrontendHasNoEmoji (0.02s)
=== RUN   TestFrontendNeverNamesAnApprovalDecision
    frontend_hygiene_test.go:314: ban #6 twin check: 84 frontend files (node_modules and dist excluded) carry no approval-decision identifier
--- PASS: TestFrontendNeverNamesAnApprovalDecision (0.07s)
=== RUN   TestReadGitSyntheticSymrefRepo
--- PASS: TestReadGitSyntheticSymrefRepo (0.01s)
=== RUN   TestReadGitDetachedSyntheticTree
--- PASS: TestReadGitDetachedSyntheticTree (0.01s)
=== RUN   TestReadGitWorktreePointerFile
--- PASS: TestReadGitWorktreePointerFile (0.01s)
=== RUN   TestReadGitRelativePointerFile
--- PASS: TestReadGitRelativePointerFile (0.01s)
=== RUN   TestReadGitVanishedPointerIsUnreadable
--- PASS: TestReadGitVanishedPointerIsUnreadable (0.00s)
=== RUN   TestReadGitNonRepoSaysSoOutLoud
--- PASS: TestReadGitNonRepoSaysSoOutLoud (0.00s)
=== RUN   TestReadGitWithoutAWorkspaceDoesNotProbe
--- PASS: TestReadGitWithoutAWorkspaceDoesNotProbe (0.00s)
=== RUN   TestGitDimensionHasNoModelCallableTool
--- PASS: TestGitDimensionHasNoModelCallableTool (0.00s)
=== RUN   TestPlantedGitToolShapesGoRed
--- PASS: TestPlantedGitToolShapesGoRed (0.01s)
=== RUN   TestGitSectionTravelsInTheSnapshotPacket
--- PASS: TestGitSectionTravelsInTheSnapshotPacket (0.01s)
=== RUN   TestPumpWithoutGitReaderSaysUnreadable
--- PASS: TestPumpWithoutGitReaderSaysUnreadable (0.00s)
=== RUN   TestReadGitOnThisRepositoryIsSelfConsistent
--- PASS: TestReadGitOnThisRepositoryIsSelfConsistent (0.00s)
=== RUN   TestReadGitForWorkspaceConsumesRewriteAccount
=== RUN   TestReadGitForWorkspaceConsumesRewriteAccount/not_rewritten_is_the_ordinary_road
=== RUN   TestReadGitForWorkspaceConsumesRewriteAccount/rewritten_must_say_so_out_loud
=== RUN   TestReadGitForWorkspaceConsumesRewriteAccount/unset_view_probes_nothing
--- PASS: TestReadGitForWorkspaceConsumesRewriteAccount (0.00s)
    --- PASS: TestReadGitForWorkspaceConsumesRewriteAccount/not_rewritten_is_the_ordinary_road (0.00s)
    --- PASS: TestReadGitForWorkspaceConsumesRewriteAccount/rewritten_must_say_so_out_loud (0.00s)
    --- PASS: TestReadGitForWorkspaceConsumesRewriteAccount/unset_view_probes_nothing (0.00s)
=== RUN   TestInboundRawCredentialShapeReachesNoArtifact35r7
=== RUN   TestInboundRawCredentialShapeReachesNoArtifact35r7/credential-write-accepted
=== RUN   TestInboundRawCredentialShapeReachesNoArtifact35r7/plain-write-with-key-shaped-value
=== RUN   TestInboundRawCredentialShapeReachesNoArtifact35r7/refused-unlisted-method-carrying-content
=== RUN   TestInboundRawCredentialShapeReachesNoArtifact35r7/refused-forged-source-carrying-content
=== RUN   TestInboundRawCredentialShapeReachesNoArtifact35r7/refused-unattached-handler-carrying-content
--- PASS: TestInboundRawCredentialShapeReachesNoArtifact35r7 (0.00s)
    --- PASS: TestInboundRawCredentialShapeReachesNoArtifact35r7/credential-write-accepted (0.00s)
    --- PASS: TestInboundRawCredentialShapeReachesNoArtifact35r7/plain-write-with-key-shaped-value (0.00s)
    --- PASS: TestInboundRawCredentialShapeReachesNoArtifact35r7/refused-unlisted-method-carrying-content (0.00s)
    --- PASS: TestInboundRawCredentialShapeReachesNoArtifact35r7/refused-forged-source-carrying-content (0.00s)
    --- PASS: TestInboundRawCredentialShapeReachesNoArtifact35r7/refused-unattached-handler-carrying-content (0.00s)
=== RUN   TestInboundLeakRulerFiresOnPlantedShapes35r7
    inbound_raw_leak_35r7_test.go:220: fired: planted-audit matches sk-[A-Za-z0-9]{12,}
    inbound_raw_leak_35r7_test.go:220: fired: planted-audit matches <canary-redacted>
    inbound_raw_leak_35r7_test.go:220: fired: planted-receipt matches sk-[A-Za-z0-9]{12,}
--- PASS: TestInboundLeakRulerFiresOnPlantedShapes35r7 (0.00s)
=== RUN   TestLegalPanelTrafficIsStillWritten35r7
--- PASS: TestLegalPanelTrafficIsStillWritten35r7 (0.00s)
=== RUN   TestDecisionRulesMirrorTicket248Ruler35r7
--- PASS: TestDecisionRulesMirrorTicket248Ruler35r7 (0.00s)
=== RUN   TestFullInboundMethodRosterIsClosed
    inbound_roster_253_test.go:519: full inbound roster held: 6 names declared in bridge.go, 6 guard case labels, 6 router case labels, 2 of them carrying no "panel." prefix ("config.get", "config.set")
--- PASS: TestFullInboundMethodRosterIsClosed (0.02s)
=== RUN   TestPlantedUnregisteredInboundMethodNameGoesRed
    inbound_roster_253_test.go:686: AC#2's judgement shape holds under this ruler: an unregistered inbound name is named in all three planted shapes (plain const, Method-named const, guard case literal) while the unplanted tree reads quiet
--- PASS: TestPlantedUnregisteredInboundMethodNameGoesRed (0.01s)
=== RUN   TestRegistrationSilencesTheRosterRuler
    inbound_roster_253_test.go:761: registration read both ways: declared AND answered reads quiet (7 names, no diff); roster-only reads red on the answered side ("config.registered253"), so adding a name to the list is not how this ruler gets satisfied
--- PASS: TestRegistrationSilencesTheRosterRuler (0.01s)
=== RUN   TestFullInboundRosterRefusesAnEmptyDeclarationFile
    inbound_roster_253_test.go:800: empty-source refusal held: 6 roster names reported missing against a file that declares none
--- PASS: TestFullInboundRosterRefusesAnEmptyDeclarationFile (0.00s)
=== RUN   TestSubsetOnlyRosterDirectionIsBlindToAPlantedName
    inbound_roster_253_test.go:832: relaxed form measured blind on purpose: with the plant in the source the subset-only reading finds no problem (missing=(none)) while the shipped equality names "config.blind253". Subset-only is not admissible as AC#2's credential
--- PASS: TestSubsetOnlyRosterDirectionIsBlindToAPlantedName (0.00s)
=== RUN   TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll
--- PASS: TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll (0.01s)
=== RUN   TestTheInstructionStatesDoNotShareOneWireShape
--- PASS: TestTheInstructionStatesDoNotShareOneWireShape (0.00s)
=== RUN   TestAPumpWithoutAnInstructionsReaderSendsNoKey
--- PASS: TestAPumpWithoutAnInstructionsReaderSendsNoKey (0.00s)
=== RUN   TestWithInstructionsKeepsTheR1CarrierContract
--- PASS: TestWithInstructionsKeepsTheR1CarrierContract (0.00s)
=== RUN   TestAnsweredPanelRoutesCarryNoApprovalDecision
    l2_grant_boundary_test.go:1304: answered inbound routes = [config.get config.set panel.attachment.add panel.message.send panel.mode.request panel.workspace.request]; 6 declared Method* constants; 0 approval-shaped names on either side
--- PASS: TestAnsweredPanelRoutesCarryNoApprovalDecision (0.01s)
=== RUN   TestRealGuardRefusesEveryAssemblableApprovalRouteName
    l2_grant_boundary_test.go:1540: behavioural sweep of the running guard: asked=528 assembled route names, answeredByRealGuard=0, hits=[]
--- PASS: TestRealGuardRefusesEveryAssemblableApprovalRouteName (0.00s)
=== RUN   TestNoInboundEnvelopeCanBindAnApprovalVerdict
    l2_grant_boundary_test.go:1579: ComposerRequest binds 24 JSON keys [artifact attachments configField configModel configProvider configValue dataBase64 declaredMime deduplicated id kind method mime name name path reason requestId sizeBytes sizeBytes source stored text to]; 0 verdict-shaped; AST scan of internal/panel: 0 findings
--- PASS: TestNoInboundEnvelopeCanBindAnApprovalVerdict (0.01s)
=== RUN   TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes
=== RUN   TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes/outcome_on_the_envelope
    l2_grant_boundary_test.go:1626: flagged as required: plantedOutcomeEnvelope.Outcome binds "outcome"
=== RUN   TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes/allowOnce_on_the_envelope
    l2_grant_boundary_test.go:1626: flagged as required: plantedAllowFlag.Allow binds "allowOnce"
=== RUN   TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes/a_verdict_smuggled_inside_an_embedded_struct
    l2_grant_boundary_test.go:1626: flagged as required: plantedVerdictIn.Decision.Verdict binds "Verdict"
--- PASS: TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes (0.00s)
    --- PASS: TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes/outcome_on_the_envelope (0.00s)
    --- PASS: TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes/allowOnce_on_the_envelope (0.00s)
    --- PASS: TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes/a_verdict_smuggled_inside_an_embedded_struct (0.00s)
=== RUN   TestJSONKeyDerivationAgreesWithEncodingJSON
=== RUN   TestJSONKeyDerivationAgreesWithEncodingJSON/i_the_AST_key_rule_matches_encoding/json's,_spelling_by_spelling
=== RUN   TestJSONKeyDerivationAgreesWithEncodingJSON/ii_the_AST_list_and_the_reflection_list_are_one_list
    l2_grant_boundary_test.go:1874: ComposerRequest: 11 keys agreed by both instruments [attachments configField configModel configProvider configValue method path requestId source text to]
    l2_grant_boundary_test.go:1874: ModeRequest: 2 keys agreed by both instruments [correlationId to]
    l2_grant_boundary_test.go:1874: AttachmentPayload: 4 keys agreed by both instruments [dataBase64 declaredMime name sizeBytes]
    l2_grant_boundary_test.go:1874: AttachmentRef: 9 keys agreed by both instruments [artifact deduplicated id kind mime name reason sizeBytes stored]
=== RUN   TestJSONKeyDerivationAgreesWithEncodingJSON/iii_the_decoder_is_the_third_vote_and_the_verdict_words_are_not_bindable
    l2_grant_boundary_test.go:1913: AttachmentPayload: 4 listed keys confirmed bindable, 24 verdict spellings confirmed not bindable
    l2_grant_boundary_test.go:1913: AttachmentRef: 9 listed keys confirmed bindable, 24 verdict spellings confirmed not bindable
    l2_grant_boundary_test.go:1913: ComposerRequest: 11 listed keys confirmed bindable, 24 verdict spellings confirmed not bindable
    l2_grant_boundary_test.go:1913: ModeRequest: 2 listed keys confirmed bindable, 24 verdict spellings confirmed not bindable
--- PASS: TestJSONKeyDerivationAgreesWithEncodingJSON (0.01s)
    --- PASS: TestJSONKeyDerivationAgreesWithEncodingJSON/i_the_AST_key_rule_matches_encoding/json's,_spelling_by_spelling (0.00s)
    --- PASS: TestJSONKeyDerivationAgreesWithEncodingJSON/ii_the_AST_list_and_the_reflection_list_are_one_list (0.00s)
    --- PASS: TestJSONKeyDerivationAgreesWithEncodingJSON/iii_the_decoder_is_the_third_vote_and_the_verdict_words_are_not_bindable (0.00s)
=== RUN   TestGrantWireShapesAreRefusedAtTheDoor
=== RUN   TestGrantWireShapesAreRefusedAtTheDoor/the_historical_wire_shape_is_refused_because_Go_does_not_answer_its_route
=== RUN   TestGrantWireShapesAreRefusedAtTheDoor/every_answered_route_drops_a_smuggled_verdict
--- PASS: TestGrantWireShapesAreRefusedAtTheDoor (0.00s)
    --- PASS: TestGrantWireShapesAreRefusedAtTheDoor/the_historical_wire_shape_is_refused_because_Go_does_not_answer_its_route (0.00s)
    --- PASS: TestGrantWireShapesAreRefusedAtTheDoor/every_answered_route_drops_a_smuggled_verdict (0.00s)
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/A_an_envelope_that_grows_an_outcome_field_goes_red
    l2_grant_boundary_test.go:2271: red as required:
          bridge.go:99: inbound envelope ComposerRequest can bind the JSON key "outcome" - a verdict a page can set, which D33/F2 and R20 keep native-side only: Outcome string `json:"outcome,omitempty"`
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/B_Go_answering_panel.approval.request_goes_red
    l2_grant_boundary_test.go:2283: red as required:
          bridge.go:146: Go answers the inbound route "panel.approval.request", whose own name is an approval decision - a panel-side allow door (AGENTS.md §1.2 ban #6, D33/F2)
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/E_a_guard_route_reached_through_a_var_or_a_concatenation_is_named,_not_dropped
    l2_grant_boundary_test.go:2320: loud as required, and the names are not invented:
          bridge.go:148: "panel.review." + acceptE3Tail
          bridge.go:148: acceptE2Route
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/F_a_route_named_outside_the_guard_still_reaches_the_pool
    l2_grant_boundary_test.go:2366: pool sees "panel.review.allow" at [l2-plant-f.go:7] while the snapshot's guard case list does not; the half that asks the RUNNING guard about that name - and about every other name grantRoutePrefixes x grantRouteWords can build - is TestRealGuardRefusesEveryAssemblableApprovalRouteName, which is standing rather than a negative control here (F-R2-3)
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/G_an_inbound_envelope_with_a_renamed_route_key_is_still_seeded_from_its_decode
    l2_grant_boundary_test.go:2412: seeded from its own decode and red as required:
          l2-plant-g.go:8: inbound envelope grantViaCmd can bind the JSON key "outcome" - a verdict a page can set, which D33/F2 and R20 keep native-side only: Outcome string `json:"outcome"`
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/H_a_verdict_tagged_without_a_name_is_still_found_under_its_field_name
    l2_grant_boundary_test.go:2456: red as required, on the spelling the wire uses:
          l2-plant-h.go:8: inbound envelope grantByFieldName can bind the JSON key "Outcome" - a verdict a page can set, which D33/F2 and R20 keep native-side only: Outcome string `json:",omitempty"`
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/C_a_wired_grant_door_goes_red_on_both_halves
    l2_grant_boundary_test.go:2477: red as required:
          bridge.go:146: Go answers the inbound route "panel.approval.request", whose own name is an approval decision - a panel-side allow door (AGENTS.md §1.2 ban #6, D33/F2)
          grant_handler.go:8: inbound envelope grantThrough can bind the JSON key "outcome" - a verdict a page can set, which D33/F2 and R20 keep native-side only: Outcome string `json:"outcome"`
=== RUN   TestPlantedGrantWiringGoesRedInASnapshot/D_the_historical_JSX_button_is_invisible_to_this_instrument
    l2_grant_boundary_test.go:2501: ZERO INSTRUMENT COVERAGE, measured not assumed: a tree holding the verbatim 53a1359^ button (send("grant") -> panel.approval.request carrying outcome:"grant") reads CLEAN to this instrument, because it reads Go and the violation is drawn in JSX.
--- PASS: TestPlantedGrantWiringGoesRedInASnapshot (0.24s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/A_an_envelope_that_grows_an_outcome_field_goes_red (0.02s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/B_Go_answering_panel.approval.request_goes_red (0.02s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/E_a_guard_route_reached_through_a_var_or_a_concatenation_is_named,_not_dropped (0.02s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/F_a_route_named_outside_the_guard_still_reaches_the_pool (0.03s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/G_an_inbound_envelope_with_a_renamed_route_key_is_still_seeded_from_its_decode (0.04s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/H_a_verdict_tagged_without_a_name_is_still_found_under_its_field_name (0.04s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/C_a_wired_grant_door_goes_red_on_both_halves (0.02s)
    --- PASS: TestPlantedGrantWiringGoesRedInASnapshot/D_the_historical_JSX_button_is_invisible_to_this_instrument (0.02s)
=== RUN   TestThePumpBuildsThePacketFromWhatTheHostHolds
--- PASS: TestThePumpBuildsThePacketFromWhatTheHostHolds (0.00s)
=== RUN   TestAPumpWithNoReadersSaysSoInsteadOfInventingState
--- PASS: TestAPumpWithNoReadersSaysSoInsteadOfInventingState (0.00s)
=== RUN   TestAnUnreadableModeFromALiveReaderStillRendersUnknown
--- PASS: TestAnUnreadableModeFromALiveReaderStillRendersUnknown (0.00s)
=== RUN   TestTheStreamLogTruncatesInsteadOfMerging
--- PASS: TestTheStreamLogTruncatesInsteadOfMerging (0.00s)
=== RUN   TestPublishWithoutAnExitReportsInsteadOfPretending
--- PASS: TestPublishWithoutAnExitReportsInsteadOfPretending (0.00s)
=== RUN   TestPublishHandsTheBytesToTheAttachedExit
--- PASS: TestPublishHandsTheBytesToTheAttachedExit (0.00s)
=== RUN   TestASecondCardForTheSameTaskStillLightsItsRow
--- PASS: TestASecondCardForTheSameTaskStillLightsItsRow (0.00s)
=== RUN   TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt
=== RUN   TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/the_queue's_own_rewrite_lights
=== RUN   TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/a_multi-digit_seq_lights
=== RUN   TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/another_task's_rewrite_does_not_light
=== RUN   TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/a_non-numeric_suffix_is_not_the_queue's_form
=== RUN   TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/a_task_id_that_merely_prefixes_does_not_light
=== RUN   TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/an_empty_correlation_id_names_no_row
--- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt (0.00s)
    --- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/the_queue's_own_rewrite_lights (0.00s)
    --- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/a_multi-digit_seq_lights (0.00s)
    --- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/another_task's_rewrite_does_not_light (0.00s)
    --- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/a_non-numeric_suffix_is_not_the_queue's_form (0.00s)
    --- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/a_task_id_that_merely_prefixes_does_not_light (0.00s)
    --- PASS: TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt/an_empty_correlation_id_names_no_row (0.00s)
=== RUN   TestAL1WindowInWaitingLightsItsRow
--- PASS: TestAL1WindowInWaitingLightsItsRow (0.00s)
=== RUN   TestABlockedRowKeepsItsRunningDimension
--- PASS: TestABlockedRowKeepsItsRunningDimension (0.00s)
=== RUN   TestTheRosterReaderPutsSubagentsOnTheWire
--- PASS: TestTheRosterReaderPutsSubagentsOnTheWire (0.00s)
=== RUN   TestAPumpWithoutARosterReaderSendsFourKeys
--- PASS: TestAPumpWithoutARosterReaderSendsFourKeys (0.00s)
=== RUN   TestAStatusOutsideD43NeverReachesTheWire
--- PASS: TestAStatusOutsideD43NeverReachesTheWire (0.00s)
=== RUN   TestTruncationFactsRideThePacket
--- PASS: TestTruncationFactsRideThePacket (0.00s)
=== RUN   TestDroppedStreamsAreNamedOnTheWire
--- PASS: TestDroppedStreamsAreNamedOnTheWire (0.00s)
=== RUN   TestSubagentStreamKeyIsTheOneSpelling
--- PASS: TestSubagentStreamKeyIsTheOneSpelling (0.00s)
=== RUN   TestNPlusOneSubagentsEachGetTheirOwnStream
--- PASS: TestNPlusOneSubagentsEachGetTheirOwnStream (0.00s)
=== RUN   TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther
--- PASS: TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther (0.00s)
=== RUN   TestOverflowTruncatesEachStreamsOwnHeadAndTail
--- PASS: TestOverflowTruncatesEachStreamsOwnHeadAndTail (0.00s)
=== RUN   TestHardCeilingDropsWholeKeysAndNamesThem
--- PASS: TestHardCeilingDropsWholeKeysAndNamesThem (0.00s)
=== RUN   TestC21DesignTokensFourWayAgree
    tokens_fourway_test.go:441: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified.
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
=== RUN   TestWorkspaceViewCarriesEveryHalfOfTheAccount
=== RUN   TestWorkspaceViewCarriesEveryHalfOfTheAccount/clean_answer_is_the_ordinary_road
=== RUN   TestWorkspaceViewCarriesEveryHalfOfTheAccount/rewritten_answer_must_be_rewritten_in_the_view
=== RUN   TestWorkspaceViewCarriesEveryHalfOfTheAccount/reparse_answer_is_read_not_assumed
=== RUN   TestWorkspaceViewCarriesEveryHalfOfTheAccount/no_answer_still_renders_unset_with_a_reason
--- PASS: TestWorkspaceViewCarriesEveryHalfOfTheAccount (0.00s)
    --- PASS: TestWorkspaceViewCarriesEveryHalfOfTheAccount/clean_answer_is_the_ordinary_road (0.00s)
    --- PASS: TestWorkspaceViewCarriesEveryHalfOfTheAccount/rewritten_answer_must_be_rewritten_in_the_view (0.00s)
    --- PASS: TestWorkspaceViewCarriesEveryHalfOfTheAccount/reparse_answer_is_read_not_assumed (0.00s)
    --- PASS: TestWorkspaceViewCarriesEveryHalfOfTheAccount/no_answer_still_renders_unset_with_a_reason (0.00s)
=== RUN   TestWorkspaceAccountReachesEveryConsumer
=== RUN   TestWorkspaceAccountReachesEveryConsumer/git_dimension_says_the_account_out_loud
=== RUN   TestWorkspaceAccountReachesEveryConsumer/project_instruction_request_refuses_both_trees
=== RUN   TestWorkspaceAccountReachesEveryConsumer/the_account_is_on_the_wire
--- PASS: TestWorkspaceAccountReachesEveryConsumer (0.01s)
    --- PASS: TestWorkspaceAccountReachesEveryConsumer/git_dimension_says_the_account_out_loud (0.00s)
    --- PASS: TestWorkspaceAccountReachesEveryConsumer/project_instruction_request_refuses_both_trees (0.00s)
    --- PASS: TestWorkspaceAccountReachesEveryConsumer/the_account_is_on_the_wire (0.00s)
=== RUN   TestScopeThatChangesItsAccountIsSurfaced
--- PASS: TestScopeThatChangesItsAccountIsSurfaced (0.00s)
=== RUN   TestWorkspaceSwitchAuditsAndAppliesACleanPath
    workspace_test.go:71: audit: workspace: SWITCH from="" to="D:\\work\\Wisp\\notes" spelling="D:\\work\\Wisp\\notes" reparse=false rewritten=false rewrites= result=ok
--- PASS: TestWorkspaceSwitchAuditsAndAppliesACleanPath (0.00s)
=== RUN   TestWorkspaceSwitchPropagatesC26sReparseRefusal
--- PASS: TestWorkspaceSwitchPropagatesC26sReparseRefusal (0.00s)
=== RUN   TestWorkspaceSwitchRefusesARewrittenAccountEvenIfTheResolverSaysOK
--- PASS: TestWorkspaceSwitchRefusesARewrittenAccountEvenIfTheResolverSaysOK (0.00s)
=== RUN   TestWorkspaceSwitchKeepsThePreviousScopeOnFailure
--- PASS: TestWorkspaceSwitchKeepsThePreviousScopeOnFailure (0.00s)
=== RUN   TestWorkspaceSwitchDetectsAScopeThatMovedAnyway
--- PASS: TestWorkspaceSwitchDetectsAScopeThatMovedAnyway (0.00s)
=== RUN   TestWorkspaceViewRenderedForSnapshots
--- PASS: TestWorkspaceViewRenderedForSnapshots (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/panel	6.271s
time=2026-10-09T14:34:31.850+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestAC1AlwaysBranchDoesNotRevertAHandEditedKey
time=2026-10-09T14:34:33.052+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:34:33 mockllm: serving on http://127.0.0.1:61852 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:34:33.077+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey240377665\002\logs min_level=info early_records=2 early_dropped=0
time=2026-10-09T14:34:33.108+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey240377665\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:34:33.116+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:34:33Z duration_ms=8
time=2026-10-09T14:34:33.125+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC1AlwaysBranchDoesNotRevertAHandEditedKey240377665\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:34:33.130+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:34:33Z duration_ms=4
time=2026-10-09T14:34:33.134+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:34:33.136+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:34:33.153+08:00 level=WARN msg="config: wrote merged change into config.toml and kept hand edits this process does not hold (ticket 226); the file was NOT claimed as our own write, so the next reload reads it back and a loosening there is denied until it is confirmed (D36 rule 1). Only the keys listed under wrote changed value; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved" key=fs.allowed_dirs wrote=[fs.allowed_dirs] kept_in_file_not_in_memory=[app.theme]
--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (1.37s)
=== RUN   TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card
time=2026-10-09T14:34:34.239+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:34:34 mockllm: serving on http://127.0.0.1:59062 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:34:34.255+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card2661091256\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:34:34.275+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card2661091256\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:34:34.281+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:34:34Z duration_ms=5
time=2026-10-09T14:34:34.288+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card2661091256\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:34:34.293+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:34:34Z duration_ms=4
time=2026-10-09T14:34:34.298+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:34:34.301+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:34:34.313+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=fs.allowed_dirs wrote=[fs.allowed_dirs]
--- PASS: TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (1.11s)
=== RUN   TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused
time=2026-10-09T14:34:35.276+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:34:35 mockllm: serving on http://127.0.0.1:53390 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:34:35.295+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused3961035345\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:34:35.323+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused3961035345\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:34:35.330+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:34:35Z duration_ms=7
time=2026-10-09T14:34:35.337+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused3961035345\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:34:35.342+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:34:35Z duration_ms=4
time=2026-10-09T14:34:35.346+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:34:35.348+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused (1.05s)
=== RUN   TestReplyListenerAllowsAnL2CardFromTheNativeSide
time=2026-10-09T14:34:36.405+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:34:36 mockllm: serving on http://127.0.0.1:54551 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:34:36.425+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide1014616510\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:34:36.451+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide1014616510\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:34:36.458+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:34:36Z duration_ms=6
time=2026-10-09T14:34:36.464+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerAllowsAnL2CardFromTheNativeSide1014616510\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:34:36.470+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:34:36Z duration_ms=5
time=2026-10-09T14:34:36.475+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:34:36.477+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:34:36.530+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerAllowsAnL2CardFromTheNativeSide (1.16s)
=== RUN   TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel
time=2026-10-09T14:34:37.508+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:34:37 mockllm: serving on http://127.0.0.1:54554 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:34:37.525+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel3213686192\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:34:37.546+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel3213686192\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:34:37.554+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:34:37Z duration_ms=8
time=2026-10-09T14:34:37.566+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel3213686192\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:34:37.572+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:34:37Z duration_ms=5
time=2026-10-09T14:34:37.577+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:34:37.578+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:34:37.618+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel (1.08s)
=== RUN   TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject
time=2026-10-09T14:34:38.571+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:34:38 mockllm: serving on http://127.0.0.1:51191 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:34:38.588+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1357433983\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:34:38.609+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1357433983\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:34:38.615+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:34:38Z duration_ms=6
time=2026-10-09T14:34:38.622+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject1357433983\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:34:38.628+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:34:38Z duration_ms=5
time=2026-10-09T14:34:38.636+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:34:38.638+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:34:38.794+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject (1.18s)
=== RUN   TestUnansweredL2CardTimesOutIntoRejectNeverExecution
time=2026-10-09T14:34:39.822+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:34:39 mockllm: serving on http://127.0.0.1:51200 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:34:39.839+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution2953485631\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:34:39.862+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution2953485631\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:34:39.868+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:34:39Z duration_ms=5
time=2026-10-09T14:34:39.874+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestUnansweredL2CardTimesOutIntoRejectNeverExecution2953485631\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:34:39.883+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:34:39Z duration_ms=8
time=2026-10-09T14:34:39.888+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:34:39.889+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:35:10.927+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestUnansweredL2CardTimesOutIntoRejectNeverExecution (32.13s)
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes
time=2026-10-09T14:35:11.891+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:11 mockllm: serving on http://127.0.0.1:54853 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:11.910+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1673063961\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:11.930+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1673063961\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:11.936+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:11Z duration_ms=6
time=2026-10-09T14:35:11.943+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredconsole_posture_veto_1673063961\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:11.948+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:11Z duration_ms=5
time=2026-10-09T14:35:11.952+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:11.954+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands
time=2026-10-09T14:35:14.899+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:14 mockllm: serving on http://127.0.0.1:57205 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:14.917+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load4060764576\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:14.937+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load4060764576\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:14.943+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:14Z duration_ms=5
time=2026-10-09T14:35:14.951+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestL1VetoNeedsAChannelTheHostReallyWiredhost_declares_and_load4060764576\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:14.956+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:14Z duration_ms=4
time=2026-10-09T14:35:14.960+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:14.962+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:35:15.000+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestL1VetoNeedsAChannelTheHostReallyWired (4.07s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/console_posture:_veto_refused,_window_still_executes (3.07s)
    --- PASS: TestL1VetoNeedsAChannelTheHostReallyWired/host_declares_and_loads_esc,_veto_lands (1.00s)
=== RUN   TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole
time=2026-10-09T14:35:15.924+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:15 mockllm: serving on http://127.0.0.1:57208 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:15.941+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole471698452\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:15.961+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole471698452\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:15.967+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:15Z duration_ms=5
time=2026-10-09T14:35:15.974+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole471698452\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:15.978+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:15Z duration_ms=4
time=2026-10-09T14:35:15.983+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:15.985+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole (1.01s)
=== RUN   TestNativeHostSeamRefusesAPanelSourcedAllow
time=2026-10-09T14:35:16.944+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:16 mockllm: serving on http://127.0.0.1:57211 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:16.960+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow2743921258\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:16.978+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow2743921258\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:16.984+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:16Z duration_ms=6
time=2026-10-09T14:35:16.992+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestNativeHostSeamRefusesAPanelSourcedAllow2743921258\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:16.997+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:16Z duration_ms=4
time=2026-10-09T14:35:17.001+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:17.003+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestNativeHostSeamRefusesAPanelSourcedAllow (1.07s)
=== RUN   TestTicket255SplitOnlyClaimsSectionsWithALiveReader
--- PASS: TestTicket255SplitOnlyClaimsSectionsWithALiveReader (0.00s)
=== RUN   TestTicket255HotRowRosterCoversTheRegistry
--- PASS: TestTicket255HotRowRosterCoversTheRegistry (0.00s)
=== RUN   TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim
--- PASS: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)
=== RUN   TestTicket255RosterStillMatchesTheActualReadSites
--- PASS: TestTicket255RosterStillMatchesTheActualReadSites (0.51s)
=== RUN   TestTicket255ReceiptOmitsPanelFromTheImmediateSentence
time=2026-10-09T14:35:18.540+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:18 mockllm: serving on http://127.0.0.1:60472 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:18.557+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence301321590\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:18.576+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence301321590\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:18.581+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:18Z duration_ms=5
time=2026-10-09T14:35:18.588+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsPanelFromTheImmediateSentence301321590\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:18.593+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:18Z duration_ms=4
time=2026-10-09T14:35:18.597+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:18.599+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsPanelFromTheImmediateSentence (2.04s)
=== RUN   TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere
time=2026-10-09T14:35:20.544+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:20 mockllm: serving on http://127.0.0.1:60476 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:20.560+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere256242582\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:20.619+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere256242582\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:20.625+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:20Z duration_ms=6
time=2026-10-09T14:35:20.632+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere256242582\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:20.637+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:20Z duration_ms=4
time=2026-10-09T14:35:20.641+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:20.643+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere (2.03s)
=== RUN   TestTicket255ReceiptStillNamesTheLiveReadSection
time=2026-10-09T14:35:22.567+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:22 mockllm: serving on http://127.0.0.1:65530 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:22.583+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection791466553\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:22.605+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection791466553\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:22.611+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:22Z duration_ms=5
time=2026-10-09T14:35:22.618+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptStillNamesTheLiveReadSection791466553\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:22.622+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:22Z duration_ms=4
time=2026-10-09T14:35:22.626+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:22.629+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket255ReceiptStillNamesTheLiveReadSection (1.99s)
=== RUN   TestTicket255ReceiptSentenceAssemblyIsFiltered
time=2026-10-09T14:35:24.555+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:24 mockllm: serving on http://127.0.0.1:59268 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:24.573+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered355934888\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:24.597+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered355934888\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:24.606+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:24Z duration_ms=8
time=2026-10-09T14:35:24.612+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket255ReceiptSentenceAssemblyIsFiltered355934888\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:24.617+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:24Z duration_ms=5
time=2026-10-09T14:35:24.621+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:24.623+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    config_receipt_255_test.go:597: SENTENCES
          IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm]"
          HONEST    "wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为（票 255 AC#1：这一半不许说成「已立即生效」；逐段的读者判定见 HOT-RELOAD-READER 行）：[ball] [session] [audio] [agent] [privacy] [memory] [panel] [cost] [models] [observe] [hotkey] [app] [voice]"
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=panel tier_row=[panel] claims=["panel":other-process: cmd/wisp/panel_resident_windows.go:207 [cfg.Panel.Width] - the resident panel host's assembly root re-reads [panel] at every window creation and the window is built from that number (cmd/wisp/panel_host_windows.go:262 [Width:  uint(width)], reached from the create at :392); 面板关窗再开即跟上新值, and since 票 255-r1 a re-show of an already-created window posts this host's currently resolved pair through the library's SetSize on Dispatch - a CLIENT-area request, not the OUTER-FRAME one the create makes, so the same width is not the same on-screen pixels; nothing presses it when config.toml is saved, so the size arrives on the next show request, and this `wisp run` process builds no panel host at all]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=llm tier_row=[llm] claims=["llm":consumed: cmd/wisp/panel_config_store.go:96 [cfg.LLM.Roles.Chat.Model] - configStore.ReadSettings calls s.mgr.Config() per call]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=app tier_row=[app.theme] claims=["app.theme":no-reader: 扫描零命中：nothing outside internal/config reads cfg.App - manager.go's planApp compares and copies it, and no component re-skins from it]
    config_receipt_255_test.go:606: LEDGER [audit] config: HOT-RELOAD-READER section=voice tier_row=[voice.punctuation voice.tts.speed voice.wake_word.thresholds voice.wake_word.veto_words] claims=["voice.punctuation":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the punctuation switch has no consumer yet | "voice.tts.speed":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the TTS knobs are not threaded to the voice path yet | "voice.wake_word.thresholds":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the wake-word tunables are not consumed by the KWS path yet | "voice.wake_word.veto_words":no-reader: 扫描零命中：nothing outside internal/config reads this [voice] key (the section's only production reader since ticket 247 is cmd/wisp/resident_audio_windows.go:197 [c.Voice.Enabled], which decides whether the capture leg is built at all) - the veto list is not consumed by the KWS path yet]
--- PASS: TestTicket255ReceiptSentenceAssemblyIsFiltered (0.98s)
=== RUN   TestTicket223RunArmsTheReloadTick
time=2026-10-09T14:35:25.559+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:25 mockllm: serving on http://127.0.0.1:59271 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:25.576+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick1407489562\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:25.596+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick1407489562\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:25.603+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:25Z duration_ms=6
time=2026-10-09T14:35:25.609+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RunArmsTheReloadTick1407489562\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:25.614+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:25Z duration_ms=4
time=2026-10-09T14:35:25.618+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:25.619+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RunArmsTheReloadTick (2.00s)
=== RUN   TestTicket223HandEditedFsLooseningCostsAnL2Card
time=2026-10-09T14:35:27.589+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:27 mockllm: serving on http://127.0.0.1:52994 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:27.609+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card2748775100\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:27.633+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card2748775100\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:27.640+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:27Z duration_ms=6
time=2026-10-09T14:35:27.648+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card2748775100\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:27.653+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:27Z duration_ms=4
time=2026-10-09T14:35:27.658+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:27.660+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:35:28.683+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.09s)
=== RUN   TestTicket223RefusedLooseningKeepsOldValues
time=2026-10-09T14:35:29.661+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:29 mockllm: serving on http://127.0.0.1:63617 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:29.678+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues4107453778\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:29.699+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues4107453778\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:29.706+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:29Z duration_ms=7
time=2026-10-09T14:35:29.712+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RefusedLooseningKeepsOldValues4107453778\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:29.718+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:29Z duration_ms=5
time=2026-10-09T14:35:29.725+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:29.727+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:35:30.730+08:00 level=WARN msg="config: locked loosening rejected; keeping previous values" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223RefusedLooseningKeepsOldValues (5.04s)
=== RUN   TestTicket223TighteningRaisesNoCard
time=2026-10-09T14:35:34.710+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:34 mockllm: serving on http://127.0.0.1:63107 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:34.727+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard1576835302\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:34.748+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard1576835302\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:34.755+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:34Z duration_ms=6
time=2026-10-09T14:35:34.763+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223TighteningRaisesNoCard1576835302\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:34.767+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:34Z duration_ms=4
time=2026-10-09T14:35:34.772+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:34.774+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:35:35.776+08:00 level=INFO msg="config: locked section tightened, hot-applied" section=fs keys=[fs.allowed_dirs]
--- PASS: TestTicket223TighteningRaisesNoCard (2.03s)
=== RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
time=2026-10-09T14:35:36.763+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:36 mockllm: serving on http://127.0.0.1:64966 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:36.781+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow4190673617\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:36.803+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow4190673617\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:36.809+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:36Z duration_ms=6
time=2026-10-09T14:35:36.817+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223ModeLooseningChangesTheRunningModeAfterAllow4190673617\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:36.822+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:36Z duration_ms=4
time=2026-10-09T14:35:36.829+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:36.830+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:35:37.852+08:00 level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.14s)
=== RUN   TestTicket223RestartTierSaysItWillNotApply
time=2026-10-09T14:35:38.916+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:38 mockllm: serving on http://127.0.0.1:59293 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:38.933+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply1286387492\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:38.951+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply1286387492\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:38.958+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:38Z duration_ms=6
time=2026-10-09T14:35:38.969+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223RestartTierSaysItWillNotApply1286387492\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:38.975+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:38Z duration_ms=5
time=2026-10-09T14:35:38.979+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:38.981+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.07s)
=== RUN   TestTicket223FailureSentencesAreDistinct
=== RUN   TestTicket223FailureSentencesAreDistinct/缺失
time=2026-10-09T14:35:40.938+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:40 mockllm: serving on http://127.0.0.1:59304 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:40.954+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2376099514\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:40.976+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2376099514\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:40.982+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:40Z duration_ms=5
time=2026-10-09T14:35:40.989+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct缺失2376099514\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:40.994+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:40Z duration_ms=5
time=2026-10-09T14:35:40.999+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:41.001+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/语法错
time=2026-10-09T14:35:42.945+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:42 mockllm: serving on http://127.0.0.1:49646 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:42.962+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错3602697446\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:42.981+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错3602697446\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:42.987+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:42Z duration_ms=5
time=2026-10-09T14:35:42.994+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct语法错3602697446\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:42.999+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:42Z duration_ms=5
time=2026-10-09T14:35:43.003+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:43.004+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面
time=2026-10-09T14:35:44.963+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:44 mockllm: serving on http://127.0.0.1:49656 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:44.981+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏4161288628\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:45.000+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏4161288628\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:45.006+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:45Z duration_ms=5
time=2026-10-09T14:35:45.012+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinct声明了版本但坏4161288628\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:45.018+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:45Z duration_ms=5
time=2026-10-09T14:35:45.022+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:45.024+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
=== RUN   TestTicket223FailureSentencesAreDistinct/schema未知键
time=2026-10-09T14:35:46.984+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:46 mockllm: serving on http://127.0.0.1:49662 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:47.001+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键194484738\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:47.021+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键194484738\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:47.027+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:47Z duration_ms=5
time=2026-10-09T14:35:47.033+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223FailureSentencesAreDistinctschema未知键194484738\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:47.038+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:47Z duration_ms=5
time=2026-10-09T14:35:47.044+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:47.046+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223FailureSentencesAreDistinct (8.06s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/缺失 (2.02s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/语法错 (1.99s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/声明了版本但坏在后面 (2.03s)
    --- PASS: TestTicket223FailureSentencesAreDistinct/schema未知键 (2.01s)
=== RUN   TestTicket223PanelInboundSaysHotReloadIsDisabled
--- PASS: TestTicket223PanelInboundSaysHotReloadIsDisabled (0.01s)
=== RUN   TestTicket223PermissionDeniedSitsInItsOwnSentence
time=2026-10-09T14:35:49.021+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:35:49 mockllm: serving on http://127.0.0.1:58826 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:35:49.037+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence2433888060\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:35:49.057+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence2433888060\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:35:49.063+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:35:49Z duration_ms=6
time=2026-10-09T14:35:49.070+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket223PermissionDeniedSitsInItsOwnSentence2433888060\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:35:49.075+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:35:49Z duration_ms=5
time=2026-10-09T14:35:49.079+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:35:49.081+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket223PermissionDeniedSitsInItsOwnSentence (2.08s)
=== RUN   TestTicket223R2FailureSentenceRouting
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_注释以方括号开头_C1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_注释以方括号开头_C1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_CRLF_E1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_CRLF_E1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_无空格_G1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_无空格_G1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明当前版_缩进版本行_L1
    config_sentences_223r2_test.go:123: ROUTING "声明当前版_缩进版本行_L1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明未来版_正文语法坏_J1
    config_sentences_223r2_test.go:123: ROUTING "声明未来版_正文语法坏_J1"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/读不出版本_A1_基线不动
    config_sentences_223r2_test.go:123: ROUTING "读不出版本_A1_基线不动"
          raw err = config: config.toml parse: toml: expected character =
          PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
=== RUN   TestTicket223R2FailureSentenceRouting/声明旧版_坏表头_A2_仍归迁移
    config_sentences_223r2_test.go:123: ROUTING "声明旧版_坏表头_A2_仍归迁移"
          raw err = config: config.toml: cannot migrate from schema version 1: not valid TOML: toml: expected character ]; the file was left untouched - fix or restore it manually, it will never be silently reset
          PRODUCTION LINE = cause=migration detail="config.toml 声明了一个这份 Wisp 不会迁移的 schema_version（文件被原样留着，不会被重置）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"
=== RUN   TestTicket223R2FailureSentenceRouting/解析得开_未知键_不抢语法错
    config_sentences_223r2_test.go:123: ROUTING "解析得开_未知键_不抢语法错"
          raw err = config: config.toml: unknown key "this_key_does_not_exist" at line 3
          PRODUCTION LINE = cause=unknown-key detail="config.toml 语法没问题，但里面有这份 schema 不认的键（拼错的键会被这样拒绝，而不是被忽略）。本次运行继续用内存里的旧配置"
=== RUN   TestTicket223R2FailureSentenceRouting/声明未来版_正文解析得开_归更高版本自己那句
    config_sentences_223r2_test.go:123: ROUTING "声明未来版_正文解析得开_归更高版本自己那句"
          raw err = config: config.toml: schema_version 99 was written by a newer build (this build understands 2); upgrade Wisp or restore a backup
          PRODUCTION LINE = cause=newer-build detail="config.toml 是由一个更新的 Wisp 写出来的（它声明的 schema_version 比这份程序懂得的高；这一条不说语法错，也不说值不合法，因为它还没走到校验）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"
--- PASS: TestTicket223R2FailureSentenceRouting (0.60s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_注释以方括号开头_C1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_CRLF_E1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_无空格_G1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明当前版_缩进版本行_L1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明未来版_正文语法坏_J1 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/读不出版本_A1_基线不动 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明旧版_坏表头_A2_仍归迁移 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/解析得开_未知键_不抢语法错 (0.07s)
    --- PASS: TestTicket223R2FailureSentenceRouting/声明未来版_正文解析得开_归更高版本自己那句 (0.07s)
=== RUN   TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128
=== RUN   TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/dev
=== RUN   TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/prod
--- PASS: TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128 (0.01s)
    --- PASS: TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/dev (0.00s)
    --- PASS: TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128/prod (0.00s)
=== RUN   TestAC2TestDataDirBranchStillResolves128
--- PASS: TestAC2TestDataDirBranchStillResolves128 (0.00s)
=== RUN   TestAC2RefusalMarkersAreNotAShortenableList128
--- PASS: TestAC2RefusalMarkersAreNotAShortenableList128 (0.00s)
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor
=== RUN   TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/resolveSecretLayout
--- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (0.13s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders (0.00s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor (0.03s)
    --- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/resolveSecretLayout (0.00s)
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/run
    dataroot_128_windows_test.go:69: AC#2 real process, leg "run": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128run3870956880\001 -> rc=2
        time=2026-10-09T14:35:54.065+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp run: 本机没有可交互控制台，本轮没有人能答复卡片：L2 卡会等到超时后按拒绝处理，L1 窗口没有人能否决（要能当场答复，请在终端里跑）
        wisp run: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list
    dataroot_128_windows_test.go:69: AC#2 real process, leg "secret-list": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128secret-list1353646917\001 -> rc=2
        time=2026-10-09T14:35:54.105+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp secret: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor
    dataroot_128_windows_test.go:69: AC#2 real process, leg "doctor": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128doctor645787144\001 -> rc=1
        time=2026-10-09T14:35:54.140+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp doctor - build chain self-check
        ------------------------------------------------------------------------
        [INFO] wisp build                         version=0.0.0-dev commit=unknown built=unknown WISP_ENV=dev
        [INFO] Go runtime                         go1.27.1 (toolchain pinned by go.mod)
        [INFO] C29 minisign public key            untrusted comment: wisp models signing key (dev)
        RWSjyHlPP9lPxdEQRvWj3zFLMbc1tTEkKMwTDuVXXQDxsWRpA/m5jk9j (placeholder until C29 lands; hardcoded into buildinfo per SPEC-11 §7.3)
        [PASS] gcc (build-time)                   gcc (Rev3, Built by MSYS2 project) 16.2.0
        [FAIL] sherpa-onnx C API                  runtime 1.13.8 does not match build pin unknown
        [FAIL] onnxruntime.dll version            file version 1.28.2.0 does not match build pin unknown
        [INFO] sherpa-onnx built against onnxruntime 1.28.2
        [PASS] DLL colocated: sherpa-onnx-c-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1281913945566\001\sherpa-onnx-c-api.dll
        [PASS] DLL colocated: sherpa-onnx-cxx-api.dll C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData1281913945566\001\sherpa-onnx-cxx-api.dll
        [INFO] deps.toml cross-check              deps.toml not found near the exe (expected for installed copies; build-time pins already verified above)
        [FAIL] data dir resolvable (dev)          数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
        ------------------------------------------------------------------------
        wisp doctor: FAIL
=== RUN   TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident
    dataroot_128_windows_test.go:69: AC#2 real process, leg "resident": APPDATA unset, WISP_ENV=dev, cwd=C:\Users\swq\AppData\Local\Temp\TestAC2RealProcessRefusesOnEveryLegWithoutAppData128resident825477533\001 -> rc=1
        time=2026-10-09T14:35:54.210+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
        wisp 0.0.0-dev (unknown, built unknown)
        WISP_ENV=dev (data dir rules: SPEC-03 §5)
        sherpa-onnx runtime version: 1.13.8
        wisp: boot failed: proc: user config dir: %AppData% is not defined
--- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 (3.31s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/run (0.08s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/secret-list (0.04s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor (0.07s)
    --- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident (0.04s)
=== RUN   TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord
    early_log_nail_130_windows_test.go:189: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord1425107035\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
    early_log_nail_130_windows_test.go:271: EARLY RECORD ON DISK: index 0 of 2, stamp 2026-10-09T14:35:57.2355498+08:00, resolver="risk.c26Pipeline" probes_passed=0x2c0fd9736150; install at index 1 stamp 2026-10-09T14:35:57.2393118+08:00
--- PASS: TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord (3.05s)
=== RUN   TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot
    early_log_nail_130_windows_test.go:283: early-record sink C:\Users\swq\AppData\Local\Temp\TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot2660307520\002\logs: 1 file(s), 2 record(s), msgs=[winsec: sealing path resolver installed wisp: persistent log sink installed]
--- PASS: TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot (3.00s)
=== RUN   TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone
time=2026-10-09T14:36:00.280+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone4000630585\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198_test.go:101: stage 1 receipt (verbatim stderr):
        wisp run: 已在 C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone4000630585\001\config.toml 新建默认配置：全部取值来自内置默认表（schema 的 default 标签），未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码
        wisp run: 缺的两样各有各的入口。key：先跑 wisp secret set <blob 名>（隐藏输入，不进 argv 也不进日志，也可以 --from-stdin 从管道喂），它把明文交给 DPAPI 存储，并回显一行 api_key_ref = "dpapi:<blob 名>"；这份文件里没有写明文 key 的字段，要补的是那个名字。不想用 DPAPI 就把 api_key_ref 写成 env:<环境变量名>，值由系统环境提供。
        wisp run: 模型：在同一份文件的 [llm.providers.<名>] 里补 api_key_ref、base_url 与 models.<id>（名字对上内置预设的，protocol 与 base_url 可以留空），再在 [llm] 的 text_chain 或 roles.chat 里点名 provider/model；wisp providers discover 与 probe 读这份文件去问真实端点，不替你写。注意：models.<id> 条目里要写 enabled = true——缺这枚键的条目视为关闭，点名它会在起动时被拒。
        wisp run: 上面那句模型只是第一样。设置页那七枚字段（服务商的 base_url、api_key_ref、凭据，模型的 context_window、price.in、price.out，还有聊天模型）今天都不建行，只改已有的行；要在这一页配上模型，得在这份文件里手加三样，缺一不可：
        wisp run: 第一样＝一节 [llm.providers.<名>]，就是服务商那一行（名字对上内置预设的，protocol 与 base_url 可以留空；非预设名必须自己写 protocol，否则这份文件加载不过）。第二样＝它的模型行 [llm.providers.<名>.models.<模型 id>]，里面写 enabled = true。第三样＝就地填已有的 [llm.roles.chat] 那一节，把 provider 与 model 两枚一起点上名（别再追加一节同名的，那在 TOML 里是 duplicate table，文件直接加载不过）。界面不会替你建这一行，它只会告诉你去哪一节建；三样齐了这七枚才全部写得进，只加第一样只解锁服务商那三枚。
        wisp run: 写不进去的时候有三种原因，各是一句不同的话，不会合成一句「配置未生效」：
          第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone4000630585\001\config.toml 现在是真的文件；首启之前没有任何旧配置可言。
          第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，加完才写得进；这一条说的不是你的值不对。
          第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，文件一个字节都没动。
        wisp run: 改完什么时候才算用上，按进程形状分三种说法，不是一句「重启就好」：
          在控制台里跑 wisp run——这个进程带着每 1s 重读一次 config.toml 的看门狗，[llm] 属可热加载档，手改的值一秒内就换进这台进程的内存；但模型通路是启动时建一次的，热加载不会替它换脑，真正发请求还是按启动时那一份。
          没有可答卡入口的常驻形状——任务腿过不去控制台那道闸时，这个进程里可能根本没有会重读盘的东西，手改与面板写在两个方向上都只能等下一次启动。
          设置页那一页——它那条腿自己明说不带轮询，写入回执固定说要重启进程；页面上的读数在重启之前也不会跟着你手改的文件走。
        wisp run: 凭据这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，把那份引用再解一次是新建端点时才做的事，所以换过 key 的引用同样要重启才算用上；这一页任何时候都不回显密钥的值，只说已录入还是没录入。
        wisp run: 文本角色未配置（Unconfigured）：config: llm: role "chat" is unset and text_chain is empty (Unconfigured)
    firstrun_198_test.go:102: stage 1 config.toml: 2571 bytes, first line "schema_version = 2"
time=2026-10-09T14:36:00.288+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone4000630585\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:00.293+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone4000630585\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone (0.02s)
=== RUN   TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten
time=2026-10-09T14:36:00.303+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten1896743795\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten (0.01s)
=== RUN   TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables
time=2026-10-09T14:36:00.311+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedT1104476415\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables (0.01s)
=== RUN   TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences
time=2026-10-09T14:36:00.321+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences3195591600\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences (0.01s)
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
--- PASS: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.07s)
=== RUN   TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags
time=2026-10-09T14:36:00.403+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTa4210777554\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_198r2_test.go:237: tag-derived default leaves: 70; file: 2571 bytes
--- PASS: TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags (0.01s)
=== RUN   TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags
--- PASS: TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags (0.00s)
=== RUN   TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile
time=2026-10-09T14:36:00.419+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile2593791870\001\data\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile (0.01s)
=== RUN   TestTicket198R2AC4ReceiptNamesTheRealEntryPoints
time=2026-10-09T14:36:00.429+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints3573565250\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:00.516+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2AC4ReceiptNamesTheRealEntryPoints3573565250\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2AC4ReceiptNamesTheRealEntryPoints (0.09s)
=== RUN   TestTicket198R2J1TheAssemblyRootStillCreatesNothing
time=2026-10-09T14:36:00.525+08:00 level=INFO msg="audit: perm: MODE-READ-FAILED path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing1279560948\\\\001\\\\config.toml\" err=config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket198R2J1TheAssemblyRootStillCreatesNothing1279560948\\001\\config.toml: The system cannot find the file specified. mode=ask_every_step origin=startup result=fail-closed detail=\"档位读不到：本进程不缓存任何上一次的宽松值，决策链不会被装配（退出码 2）\""
time=2026-10-09T14:36:00.526+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198R2J1TheAssemblyRootStillCreatesNothing1279560948\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket198R2J1TheAssemblyRootStillCreatesNothing (0.01s)
=== RUN   TestTicket257R2AC1ReceiptStatesTheNonPresetCondition
time=2026-10-09T14:36:00.537+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptStatesTheNonPresetCondition59287757\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_nonpreset_test.go:88: premise holds: "t257-r2-ghost-gateway" is absent from the 12 built-in presets [anthropic deepseek mimo minimax moonshot ollama openai openrouter qwen siliconflow stepfun zhipu]
--- PASS: TestTicket257R2AC1ReceiptStatesTheNonPresetCondition (0.01s)
=== RUN   TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads
time=2026-10-09T14:36:00.548+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads1465615102\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:00.561+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.base_url wrote=[llm.providers.t257-r2-ghost-gateway.base_url]
time=2026-10-09T14:36:00.565+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.t257-r2-ghost-gateway.api_key_ref wrote=[llm.providers.t257-r2-ghost-gateway.api_key_ref]
--- PASS: TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads (0.02s)
=== RUN   TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain
time=2026-10-09T14:36:00.569+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain2497260847\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_257_test.go:216: AC#1 clean machine: 7 refused writes by reason: map[第 2 种拒因：行不存在:7]
--- PASS: TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain (0.01s)
=== RUN   TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields
time=2026-10-09T14:36:00.583+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields4277683427\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:00.601+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T14:36:00.605+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T14:36:00.609+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
time=2026-10-09T14:36:00.613+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.in wrote=[llm.providers.deepseek.models.deepseek-chat.price.in]
time=2026-10-09T14:36:00.615+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.price.out wrote=[llm.providers.deepseek.models.deepseek-chat.price.out]
time=2026-10-09T14:36:00.618+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.roles.chat.model wrote=[llm.roles.chat.model]
time=2026-10-09T14:36:00.622+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields (0.04s)
=== RUN   TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven
time=2026-10-09T14:36:00.626+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven1656024319\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:00.641+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.base_url wrote=[llm.providers.deepseek.base_url]
time=2026-10-09T14:36:00.644+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
time=2026-10-09T14:36:00.648+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven (0.03s)
=== RUN   TestTicket257R2AC2ThreeRefusalsStayThreeSentences
time=2026-10-09T14:36:00.651+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2ThreeRefusalsStayThreeSentences3877083503\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2ThreeRefusalsStayThreeSentences (0.02s)
=== RUN   TestTicket257R2AC2EffectTimingSaysThreeProcessShapes
time=2026-10-09T14:36:00.667+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC2EffectTimingSaysThreeProcessShapes1607258361\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC2EffectTimingSaysThreeProcessShapes (0.01s)
=== RUN   TestTicket257R2AC3CredentialSurfaceUntouched
time=2026-10-09T14:36:00.678+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket257R2AC3CredentialSurfaceUntouched500700659\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket257R2AC3CredentialSurfaceUntouched (0.01s)
=== RUN   TestTicket198AC3CreatedFileLandsPrivate
time=2026-10-09T14:36:00.689+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate858942193\001\logs min_level=info early_records=0 early_dropped=0
    firstrun_acl_198_windows_test.go:35: AC#3 icacls(cfgPath) verbatim:
        C:\Users\swq\AppData\Local\Temp\TestTicket198AC3CreatedFileLandsPrivate858942193\001\config.toml NT AUTHORITY\SYSTEM:(F)
                                                                                                         BUILTIN\Administrators:(F)
                                                                                                         DESKTOP-LVS7839\swq:(F)
        
        Successfully processed 1 files; Failed processing 0 files
--- PASS: TestTicket198AC3CreatedFileLandsPrivate (0.03s)
=== RUN   TestRunPacketCarriesTheLoadedInstructionFiles
time=2026-10-09T14:36:01.664+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:36:01 mockllm: serving on http://127.0.0.1:51644 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:36:01.682+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles630222674\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:01.705+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles630222674\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:36:01.711+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:36:01Z duration_ms=5
time=2026-10-09T14:36:01.717+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles630222674\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:36:01.726+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:36:01Z duration_ms=8
time=2026-10-09T14:36:01.731+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:36:01.732+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:179: packet instructions section: status=loaded reason= files=[instr-ws/AGENTS.md tier=project depth=0 bytes=269]
--- PASS: TestRunPacketCarriesTheLoadedInstructionFiles (1.04s)
=== RUN   TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound
time=2026-10-09T14:36:02.706+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:36:02 mockllm: serving on http://127.0.0.1:51648 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:36:02.727+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3503584877\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:02.753+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3503584877\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:36:02.759+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:36:02Z duration_ms=6
time=2026-10-09T14:36:02.770+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3503584877\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:36:02.775+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:36:02Z duration_ms=5
time=2026-10-09T14:36:02.779+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:36:02.781+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:36:03.677+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:36:03 mockllm: serving on http://127.0.0.1:65147 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:36:03.696+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3503584877\004\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:03.717+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3503584877\004\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:36:03.723+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:36:03Z duration_ms=5
time=2026-10-09T14:36:03.730+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound3503584877\004\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:36:03.735+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:36:03Z duration_ms=5
time=2026-10-09T14:36:03.740+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:36:03.742+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    instructions_200r2_test.go:251: off packet section: status=off reason=已按你的配置跳过：agent.project_instructions_enabled=false，本轮一份项目说明都没有读（不是没找到，是被配置关掉的）。 files=[]
    instructions_200r2_test.go:252: on  packet section: status=loaded reason= files=[004/AGENTS.md tier=global depth=-1 bytes=274]
--- PASS: TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound (2.01s)
=== RUN   TestAC1AC2DispatchHopGate133
    leg_dispatch_gate_133_test.go:242: dispatch ledger, read out of func main's own branches at run time (12 legs, 4 claims in this gate's registry):
          leg default        main.go:127              installs=false handoff=false covered=ruling main.go:83                                              entries=attachParentConsole aliases=
          leg doctor         main.go:95               installs=false handoff=false covered=test TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 drives cmdDoctor entries=attachParentConsole|cmdDoctor aliases=
          leg help           main.go:124              installs=false handoff=false covered=ruling main.go:79                                              entries=attachParentConsole aliases=--help|-h
          leg models         main.go:102              installs=true  handoff=true  covered=nail TestAC2ModelsLegBooksItsHandOffVerdictOnDisk -> cmdModels entries=attachParentConsole|cmdModels aliases=
          leg no-args        main.go:60               installs=true  handoff=false covered=nail TestAC1ResidentLegInstallsItsLogListenerOnDisk -> runResident entries=attachParentConsole|runResident aliases=
          leg panel-assets   main.go:112              installs=false handoff=false covered=test TestAC1AC2TaintSourceLegProducesAJudgedR4 drives cmdPanelAssets entries=attachParentConsole|cmdPanelAssets aliases=
          leg panel-inbound  main.go:115              installs=false handoff=false covered=test TestAC1InboundLegAnswersSettingsRouteEndToEnd drives cmdPanelInbound entries=attachParentConsole|cmdPanelInbound aliases=
          leg providers      main.go:92               installs=false handoff=false covered=test TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 drives cmdProviders entries=attachParentConsole|cmdProviders aliases=
          leg run            main.go:89               installs=true  handoff=false covered=nail TestAC2SealNoticeLandsInTheRunLegLogFile -> runTextTask   entries=attachParentConsole|cmdRun aliases=
          leg secret         main.go:100              installs=true  handoff=false covered=nail TestAC3SecretLegBooksItsAuditRecordsOnDisk -> cmdSecret   entries=cmdSecret aliases=
          leg slo            main.go:110              installs=false handoff=false covered=test TestSLO144CorruptReportIsJudgedOnTheFirstRead drives contradictSubjectReportErr entries=cmdSLO aliases=
          leg version        main.go:121              installs=false handoff=false covered=ruling main.go:74                                              entries=attachParentConsole|printVersions aliases=--version|-v
    leg_dispatch_gate_133_test.go:244: blindness disclosure: run-roster disclosure: GOOS=windows, 273 startable cases read from this binary itself (`-test.list '.*'`); case names this round's ledger credited through `covered=test`: 4 distinct, 4 of them startable (one outside the roster is red above); this gate's registry: 4 claims, 0 of them outside this round's roster - every registry claim this gate makes is a case this round's binary can start. This round is a whole-package round (no -test.run filter), so the roster is the set this run's === RUN lines are drawn from.
--- PASS: TestAC1AC2DispatchHopGate133 (0.14s)
=== RUN   TestAC4EveryLegIsNailedOrRuled
    leg_sink_gate_131_test.go:328: nail entry claims, checked against the dispatch closure:
          nail "TestAC1ResidentLegInstallsItsLogListenerOnDisk" -> leg "no-args", entry "runResident": declared with the "subprocess:" prefix, so this row is a name-only claim about a real process and NOT evidence that anybody asserts "no-args"'s semantics in-process.
          nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" -> leg "models", entry "cmdModels": checked, the case calls the entry this leg's dispatch reaches
          nail "TestAC2SealNoticeLandsInTheRunLegLogFile" -> leg "run", entry "runTextTask": checked, the case calls the entry this leg's dispatch reaches
          nail "TestAC3SecretLegBooksItsAuditRecordsOnDisk" -> leg "secret", entry "cmdSecret": checked, the case calls the entry this leg's dispatch reaches
    leg_sink_gate_131_test.go:407: leg ledger, enumerated from source at run time:
          leg --help       main.go:124                install=false records=false ruled=false nails=-                                                          -> no records
          leg --version    main.go:121                install=false records=false ruled=false nails=-                                                          -> no records
          leg -h           main.go:124                install=false records=false ruled=false nails=-                                                          -> no records
          leg -v           main.go:121                install=false records=false ruled=false nails=-                                                          -> no records
          leg default      main.go:127                install=false records=false ruled=false nails=-                                                          -> no records
          leg doctor       main.go:95                 install=false records=false ruled=false nails=-                                                          -> no records
          leg help         main.go:124                install=false records=false ruled=false nails=-                                                          -> no records
          leg models       main.go:102                install=true  records=true  ruled=false nails=TestAC2ModelsLegBooksItsHandOffVerdictOnDisk               -> nailed
          leg no-args      main.go:60                 install=true  records=true  ruled=false nails=TestAC1ResidentLegInstallsItsLogListenerOnDisk             -> nailed
          leg panel-assets main.go:112                install=false records=false ruled=false nails=-                                                          -> no records
          leg panel-inbound main.go:115                install=false records=false ruled=false nails=-                                                          -> no records
          leg providers    main.go:92                 install=false records=false ruled=false nails=-                                                          -> no records
          leg run          main.go:89                 install=true  records=true  ruled=false nails=TestAC2SealNoticeLandsInTheRunLegLogFile                   -> nailed
          leg secret       main.go:100                install=true  records=true  ruled=false nails=TestAC3SecretLegBooksItsAuditRecordsOnDisk                 -> nailed
          leg slo          main.go:110                install=false records=true  ruled=true  nails=-                                                          -> ruled
          leg version      main.go:121                install=false records=false ruled=false nails=-                                                          -> no records
--- PASS: TestAC4EveryLegIsNailedOrRuled (0.08s)
=== RUN   TestAC2ModelsLegBooksItsHandOffVerdictOnDisk
    leg_sink_nail_131_windows_test.go:358: leg sink C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk1080766909\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id "absent-in-signed-manifest-131" not in signed manifest]
    leg_sink_nail_131_windows_test.go:396: MODELS LEG RECORDED ON DISK: "models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id \"absent-in-signed-manifest-131\" not in signed manifest" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk1080766909\001\logs, exit=1)
--- PASS: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.01s)
=== RUN   TestAC3SecretLegBooksItsAuditRecordsOnDisk
    leg_sink_nail_131_windows_test.go:427: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk69945905\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob]
    leg_sink_nail_131_windows_test.go:459: leg sink C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk69945905\001\logs: 4 record(s) [INFO/wisp: persistent log sink installed INFO/wisp secret: stored dpapi blob INFO/wisp: persistent log sink installed WARN/wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref]
    leg_sink_nail_131_windows_test.go:479: SECRET LEG RECORDED ON DISK: "wisp secret: stored dpapi blob" and "wisp secret unset: deleted dpapi:nail131 (WISP_ENV=test, portable=false, forced=true, referenced_fields=1) fields=llm.providers.acme.api_key_ref" (install dir=C:\Users\swq\AppData\Local\Temp\TestAC3SecretLegBooksItsAuditRecordsOnDisk69945905\001\logs)
--- PASS: TestAC3SecretLegBooksItsAuditRecordsOnDisk (0.02s)
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models
=== RUN   TestAC2AC3DegradedLegsStillDeliverTheirVerdict/secret
time=2026-10-09T14:36:04.023+08:00 level=INFO msg="wisp secret: stored dpapi blob" ref=dpapi:nail131-degraded env=test portable=false overwrote=false
--- PASS: TestAC2AC3DegradedLegsStillDeliverTheirVerdict (0.01s)
    --- PASS: TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models (0.00s)
    --- PASS: TestAC2AC3DegradedLegsStillDeliverTheirVerdict/secret (0.01s)
=== RUN   TestAC3LogSinkLandsInsideTheEnvDataRoot
--- PASS: TestAC3LogSinkLandsInsideTheEnvDataRoot (0.00s)
=== RUN   TestAC3EmptyDataRootIsARefusalNotAFallback
--- PASS: TestAC3EmptyDataRootIsARefusalNotAFallback (0.00s)
=== RUN   TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole
--- PASS: TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole (0.01s)
=== RUN   TestAC2SealNoticeLandsInTheRunLegLogFile
time=2026-10-09T14:36:04.071+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2333006923\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:04.072+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2333006923\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:167: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2333006923\001\logs: 1 file(s), 682 bytes, 2 record(s)
    logsink_windows_test.go:196: RECORDED ON DISK: level=WARN msg=winsec: seal cleared principals that stood on this object path=C:\Users\swq\AppData\Local\Temp\TestAC2SealNoticeLandsInTheRunLegLogFile2333006923\001\secrets kind=explicit cleared="S-1-1-0(A;OICI;0x1200a9;;;WD)"
--- PASS: TestAC2SealNoticeLandsInTheRunLegLogFile (0.06s)
=== RUN   TestAC2AuditTrailLandsInTheRunLegLogFile
time=2026-10-09T14:36:05.000+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:36:05 mockllm: serving on http://127.0.0.1:65150 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:36:05.017+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1959517630\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:05.035+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1959517630\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:36:05.041+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:36:05Z duration_ms=5
time=2026-10-09T14:36:05.047+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1959517630\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:36:05.052+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:36:05Z duration_ms=5
time=2026-10-09T14:36:05.056+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:36:05.059+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    logsink_windows_test.go:274: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2AuditTrailLandsInTheRunLegLogFile1959517630\002\logs: 1 file(s), 4345 bytes, 19 record(s)
--- PASS: TestAC2AuditTrailLandsInTheRunLegLogFile (0.99s)
=== RUN   TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite
time=2026-10-09T14:36:05.105+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite692027622\001\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:05.106+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite692027622\001\secrets kind=explicit cleared=S-1-1-0(A;OICI;0x1200a9;;;WD) cleared_inherited="" policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
    logsink_windows_test.go:326: sink dir C:\Users\swq\AppData\Local\Temp\TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite692027622\001\logs: 1 file(s), 706 bytes, 2 record(s)
--- PASS: TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite (0.03s)
=== RUN   TestAC1AC2TaintSourceLegProducesAJudgedR4
--- PASS: TestAC1AC2TaintSourceLegProducesAJudgedR4 (0.00s)
=== RUN   TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard
--- PASS: TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard (0.00s)
=== RUN   TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit
=== RUN   TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/declared_but_absent_from_the_outgoing_call
=== RUN   TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/fragment_below_the_contract_floor_cannot_match
--- PASS: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit (0.01s)
    --- PASS: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/declared_but_absent_from_the_outgoing_call (0.00s)
    --- PASS: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit/fragment_below_the_contract_floor_cannot_match (0.00s)
=== RUN   TestAC2TaintSourceIsVisibleInTheUsageBlock
--- PASS: TestAC2TaintSourceIsVisibleInTheUsageBlock (0.00s)
=== RUN   TestAC1MalformedTaintSourceIsRefused
=== RUN   TestAC1MalformedTaintSourceIsRefused/no_separators
=== RUN   TestAC1MalformedTaintSourceIsRefused/two_parts
=== RUN   TestAC1MalformedTaintSourceIsRefused/empty_tool
=== RUN   TestAC1MalformedTaintSourceIsRefused/empty_content
--- PASS: TestAC1MalformedTaintSourceIsRefused (0.01s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/no_separators (0.00s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/two_parts (0.00s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/empty_tool (0.00s)
    --- PASS: TestAC1MalformedTaintSourceIsRefused/empty_content (0.00s)
=== RUN   TestAC1CmdSideEmitsNoVerdictTokens
--- PASS: TestAC1CmdSideEmitsNoVerdictTokens (0.00s)
=== RUN   TestAC2CredentialSentinelAppearsInNoArtifact
--- PASS: TestAC2CredentialSentinelAppearsInNoArtifact (0.01s)
=== RUN   TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere
--- PASS: TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere (0.01s)
=== RUN   TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef
time=2026-10-09T14:36:05.175+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef (0.01s)
=== RUN   TestAC2SharedEnvelopeCannotCarryTheCredentialValue
--- PASS: TestAC2SharedEnvelopeCannotCarryTheCredentialValue (0.00s)
=== RUN   TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope
--- PASS: TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope (0.00s)
=== RUN   TestAC1InboundLegAnswersSettingsRouteEndToEnd
time=2026-10-09T14:36:05.190+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.models.deepseek-chat.context_window wrote=[llm.providers.deepseek.models.deepseek-chat.context_window]
--- PASS: TestAC1InboundLegAnswersSettingsRouteEndToEnd (0.01s)
=== RUN   TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket
--- PASS: TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket (0.01s)
=== RUN   TestAC7InvalidSettingsValueIsRefusedBeforeTheFile
--- PASS: TestAC7InvalidSettingsValueIsRefusedBeforeTheFile (0.01s)
=== RUN   TestAC2SnapshotReportsRefWithoutBlobAsAPartState
time=2026-10-09T14:36:05.217+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=llm.providers.deepseek.api_key_ref wrote=[llm.providers.deepseek.api_key_ref]
--- PASS: TestAC2SnapshotReportsRefWithoutBlobAsAPartState (0.01s)
=== RUN   TestTransportDoorBindingMatchesRoster253r1
    panel_dispatch_binding_roster_253r1_windows_test.go:358: 253-r1 census: door constant "panelDispatchBinding" = "wispDispatch"; installPanelTransport at panel_host_windows.go:801 made 1 Bind call(s) and 1 Init call(s); roster holds 1 name(s) over 35 production file(s)
    panel_dispatch_binding_roster_253r1_windows_test.go:401: 253-r1 lockstep: w.Init(panelPostMessageForwardInit) at panel_host_windows.go:808 derives its window.<door> reference from "panelDispatchBinding", the same constant the Bind uses
--- PASS: TestTransportDoorBindingMatchesRoster253r1 (0.01s)
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/good
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "good": const="wispDispatch" bound=[wispDispatch] initTied="panelPostMessageForwardInit" reds=[]
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/drifted-bind
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "drifted-bind": const="wispDispatch" bound=[wispStaleDoor] initTied="panelPostMessageForwardInit" reds=[rostered but unbound door wispDispatch unrostered bound door wispStaleDoor]
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/empty-bind
    panel_dispatch_binding_roster_253r1_windows_test.go:563: empty-bind: census bound zero readable doors (1 Bind call(s)) - the disk case reddens this as the empty-roster Fatalf
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "empty-bind": const="wispDispatch" bound=[] initTied="panelPostMessageForwardInit" reds=[rostered but unbound door wispDispatch]
=== RUN   TestBindingRosterBitesItsOwnFixtures253r1/untied-init
    panel_dispatch_binding_roster_253r1_windows_test.go:568: fixture "untied-init": const="wispDispatch" bound=[wispDispatch] initTied="" reds=[w.Init script not tied to the door constant]
--- PASS: TestBindingRosterBitesItsOwnFixtures253r1 (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/good (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/drifted-bind (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/empty-bind (0.00s)
    --- PASS: TestBindingRosterBitesItsOwnFixtures253r1/untied-init (0.00s)
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/no_source_at_all_keeps_the_host's_own_constants
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/a_config_width_and_height_both_arrive
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/height_0_is_NOT_auto-height:_the_ruling_is_260,_and_schema.go's_PanelSection.Height_carries_no_default_tag_so_0_is_what_a_config_saying_nothing_answers
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/width_0_falls_back_to_the_host_constant,_not_to_a_guess
=== RUN   TestTicket255WindowOptionsFollowTheConfigSource/a_negative_answer_is_garbage,_and_garbage_is_not_a_second_spelling_of_auto
--- PASS: TestTicket255WindowOptionsFollowTheConfigSource (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/no_source_at_all_keeps_the_host's_own_constants (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/a_config_width_and_height_both_arrive (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/height_0_is_NOT_auto-height:_the_ruling_is_260,_and_schema.go's_PanelSection.Height_carries_no_default_tag_so_0_is_what_a_config_saying_nothing_answers (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/width_0_falls_back_to_the_host_constant,_not_to_a_guess (0.00s)
    --- PASS: TestTicket255WindowOptionsFollowTheConfigSource/a_negative_answer_is_garbage,_and_garbage_is_not_a_second_spelling_of_auto (0.00s)
=== RUN   TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions
time=2026-10-09T14:36:05.237+08:00 level=WARN msg="panel host: [panel] geometry source unreadable, sizing at the host's own default" path=C:\Users\swq\AppData\Local\Temp\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions1541739637\001\no-such-dir\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions1541739637\\001\\no-such-dir\\config.toml: The system cannot find the path specified." default=420x260
--- PASS: TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions (0.00s)
=== RUN   TestTicket255PanelHostBuildsItsWindowOptions
--- PASS: TestTicket255PanelHostBuildsItsWindowOptions (0.01s)
=== RUN   TestTicket255HostStillDoesNotParseConfigItself
--- PASS: TestTicket255HostStillDoesNotParseConfigItself (0.00s)
=== RUN   TestTicket255PanelRosterVerdictIsTheHonestShape
--- PASS: TestTicket255PanelRosterVerdictIsTheHonestShape (0.00s)
=== RUN   TestPanelHostOpensNoListeningSocketL1
--- PASS: TestPanelHostOpensNoListeningSocketL1 (0.00s)
=== RUN   TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12
    panel_host_gate_test.go:109: AC#12 reading (head e7af1796): shape=page-bundle built=true entry-bytes=1044 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil> manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=2
--- PASS: TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 (0.12s)
=== RUN   TestAC1SessionDisposeHasAProductionTrigger_AC1
    panel_host_gate_test.go:179: AC#1 dispose scan: 2 production constructor(s) [panel_resident_windows.go:253 resident_windows.go:151] | 2 manager-teardown site(s) [panel_resident_windows.go:447 panel_resident_windows.go:501] | 2 other .Destroy() call site(s), WebView2 control teardown and NOT a session dispose [panel_host_windows.go:409 panel_host_windows.go:731]
    panel_host_gate_test.go:188: AC#1 dispose reachability: 2 production constructor(s) [panel_resident_windows.go:253, resident_windows.go:151], 2 manager teardown site(s) [panel_resident_windows.go:447, panel_resident_windows.go:501]
--- PASS: TestAC1SessionDisposeHasAProductionTrigger_AC1 (0.02s)
=== RUN   TestCleanCheckoutBuilds_AC11
    panel_host_gate_test.go:414: clean-checkout go build ./... ok in C:\Users\swq\AppData\Local\Temp\TestCleanCheckoutBuilds_AC112489420679\001
--- PASS: TestCleanCheckoutBuilds_AC11 (19.53s)
=== RUN   TestPanelHostRealWindowHopAndLifecycle
    panel_host_windows_test.go:660: cold bring-up measured on this box: -1.000 ms (HEAD e7af1796 at read time, 2026-10-09T14:36:29+08:00)
    panel_host_windows_test.go:662: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window
--- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.05s)
=== RUN   TestAC4FocusReturnToPriorWindowGap33r5
    panel_host_windows_test.go:944: AC#4 focus hop (head e7af1796): foreground before any panel 0x206b6 | the ruler's own editor window 0x1d60f5a | foreground while hidden (prior) 0x206b6 | after Show 0x206b6 | panel hwnd 0x4780d9e | prevFocus recorded at Show 0x206b6 | after Hide 0x206b6 | Hide attempted restore to 0x206b6 (SetForegroundWindow 1, SetFocus 0)
    panel_host_windows_test.go:948: the panel did not take the foreground on Show: foreground 0x206b6, panel hwnd 0x4780d9e (AC#4 says only the panel takes focus when shown)
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (5.10s)
=== RUN   TestPanelHostLatencyPercentilesAC2
    panel_host_windows_test.go:985: no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate
--- SKIP: TestPanelHostLatencyPercentilesAC2 (0.00s)
=== RUN   TestAC3ListeningSocketRulerSeesItsOwnListener
    panel_host_windows_test.go:1054: AC#3 positive control: IPv4 table rows=481, this pid owned 0 before / 1 while 127.0.0.1:60607 is listening
    panel_host_windows_test.go:1076: AC#3 positive control, IPv6 family: listening on [::1]:60608, table rows=33, this pid owns 1 LISTEN rows (baseline 0)
    panel_host_windows_test.go:1102: AC#3 positive control verdict: ruler counts a real listener (0 -> 1) and stops after Close (0)
--- PASS: TestAC3ListeningSocketRulerSeesItsOwnListener (0.01s)
=== RUN   TestAC9InboundLegFromStdinReachesTheWriteLeg
time=2026-10-09T14:36:35.097+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
--- PASS: TestAC9InboundLegFromStdinReachesTheWriteLeg (0.02s)
=== RUN   TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt
--- PASS: TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt (0.01s)
=== RUN   TestAC9InboundLegRefusesRosterMethodWithNoHandler
--- PASS: TestAC9InboundLegRefusesRosterMethodWithNoHandler (0.01s)
=== RUN   TestAC9ComposerDispatchHasAProductionCaller
--- PASS: TestAC9ComposerDispatchHasAProductionCaller (0.02s)
=== RUN   TestAC9InboundFlagSurfaceIsNarrow
--- PASS: TestAC9InboundFlagSurfaceIsNarrow (0.01s)
=== RUN   TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge
    panel_inbound_guards_35r3_test.go:153: (d)-1 roster guard on the page edge: door reply="面板请求被拒绝 [panel.approval.request (无 requestId)]：panel: composer request refused: 方法 \"panel.approval.request\" 不是面板 composer 通路的能力入口"
--- PASS: TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge (0.00s)
=== RUN   TestInboundSourceGuardRefusesForeignSourceOnPageEdge
    panel_inbound_guards_35r3_test.go:161: (d)-2 source guard on the page edge: door reply="面板请求被拒绝 [panel.mode.request pc-35r3-2-9e7a-4c1b]：panel: composer request refused: 来源 \"panel-composer-x\" 不是 \"panel-composer\"，按伪造/串台拒绝（requestId=\"pc-35r3-2-9e7a-4c1b\"）"
--- PASS: TestInboundSourceGuardRefusesForeignSourceOnPageEdge (0.00s)
=== RUN   TestInboundRequestIDGuardRefusesMissingIDOnPageEdge
    panel_inbound_guards_35r3_test.go:170: (d)-3 requestId guard on the page edge: door reply="面板请求被拒绝 [panel.mode.request (无 requestId)]：panel: composer request refused: 缺少 requestId，无法与审计/卡片对齐（method=\"panel.mode.request\"）"
--- PASS: TestInboundRequestIDGuardRefusesMissingIDOnPageEdge (0.00s)
=== RUN   TestLockedSuffixNeverTakesTheLockItself33r11
    panel_locked_naming_33r11_windows_test.go:233: 33-r11 census: 1 Locked-suffixed func(s) across 35 production source file(s) of package main
    panel_locked_naming_33r11_windows_test.go:235:   PanelManager.setPriorFocusLocked in panel_host_windows.go line 625: 0 lock operation(s) in its own body []
--- PASS: TestLockedSuffixNeverTakesTheLockItself33r11 (0.01s)
=== RUN   TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11
    panel_locked_naming_33r11_windows_test.go:305: 33-r11 behaviour: setPriorFocusLocked returned inside the budget while m.mu was held, and it did record the sample
--- PASS: TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11 (0.00s)
=== RUN   TestLockedNamingCensusBitesItsOwnFixture33r11
    panel_locked_naming_33r11_windows_test.go:398: 33-r11 positive control: planted fixtureHost.parkedLocked reported with [f.mu.Lock f.mu.Unlock], honest name clean, unsuffixed self-locker untouched
--- PASS: TestLockedNamingCensusBitesItsOwnFixture33r11 (0.00s)
=== RUN   TestLockedNamingCensusRejectsTheOldName33r11
    panel_locked_naming_33r11_windows_test.go:416: 33-r11 positive control: the pre-rename shape is reported as roundTripHost.firstRoundTripLocked with [m.mu.Lock m.mu.Unlock]
--- PASS: TestLockedNamingCensusRejectsTheOldName33r11 (0.00s)
=== RUN   TestCompletedWithinReportsAParkedCall33r11
    panel_locked_naming_33r11_windows_test.go:435: 33-r11 positive control: a self-locking Locked call parked, was reported as parked, and had touched nothing
--- PASS: TestCompletedWithinReportsAParkedCall33r11 (0.20s)
=== RUN   TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe
    panel_pageover_33r10_windows_test.go:246: AC#13 33-r10 capture: probe document is 136 byte(s), round trip rtMs=0
    panel_pageover_33r10_windows_test.go:248: AC#13 33-r10 world: embed resolves 1044 entry byte(s), 1 element id(s) [root]
    panel_pageover_33r10_windows_test.go:279: AC#13 33-r10 read: last document 1044 byte(s), 1 of 1 entry id(s) present, entry bytes carried=true
--- PASS: TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe (0.00s)
=== RUN   TestRunBooksWithASnapshotOfItsLiveQueue
time=2026-10-09T14:36:36.340+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:36:36 mockllm: serving on http://127.0.0.1:64524 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:36:36.358+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue313222369\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:36:36.379+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue313222369\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:36:36.384+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:36:36Z duration_ms=5
time=2026-10-09T14:36:36.392+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunBooksWithASnapshotOfItsLiveQueue313222369\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:36:36.397+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:36:36Z duration_ms=5
time=2026-10-09T14:36:36.401+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:36:36.404+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:218: booked record: map[at:2026-10-09T06:36:36Z bytes:1652 depth:1 mode:ask_high_risk pending:pump-corr results:0/0 sha256:d6a7018de397a6ce ws:unset]
    panel_pump_test.go:219: packet bytes (1652): {"pending":[{"correlationId":"pump-corr","tool":"fs.read","args":["{\"path\":\"C:/Windows/win.ini\"}"],"level":"L2","rulesHit":["R2"],"reason":"R2: 目标路径在授权目录之外: C:\\Windows\\win.ini","reasonKnown":true,"sessionOverrideBlocked":false,"callChain":[],"decidedBy":"native"}],"results":[],"composer":{"mode":{"current":"ask_high_risk","names":["ask_every_step","ask_high_risk","auto_approve"],"l2ConfirmNames":["auto_approve"]},"workspace":{"set":false,"spelling":"","canonical":"","reparse":false,"rewritten":false,"reason":"未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定"},"attachments":[],"acceptedAttachmentMimes":["image/png","image/jpeg","image/gif","image/webp","video/mp4","video/quicktime"],"maxAttachmentBytes":67108864,"attachmentError":"","git":{"kind":"unreadable","reason":"本轮未选择工作区，git 这一维没有可探测的目录","branch":"","detachedSha":"","isDetached":false,"repoRoot":"","currentWorktree":"","worktrees":[],"branches":[],"switchBlocked":"切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主），面板只有快照这一条出向通道"},"currentModel":"m1","modelKnown":true,"credentialState":"all_recorded","credentialKnown":true},"generatedAt":"2026-10-09T06:36:36Z","instructions":{"status":"not_run","reason":"加载器已经接线，但这一轮还没有跑过加载，所以这里是没有读数，不是没有说明文件。","files":[]},"tasks":{"rows":[],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}}
--- PASS: TestRunBooksWithASnapshotOfItsLiveQueue (32.11s)
=== RUN   TestSnapshotWorkspaceSectionReportsTheNarrowing
time=2026-10-09T14:37:08.396+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:37:08 mockllm: serving on http://127.0.0.1:62344 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:37:08.414+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing3486637003\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:37:08.440+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing3486637003\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:37:08.446+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:37:08Z duration_ms=6
time=2026-10-09T14:37:08.454+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestSnapshotWorkspaceSectionReportsTheNarrowing3486637003\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:37:08.460+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:37:08Z duration_ms=5
time=2026-10-09T14:37:08.465+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:37:08.467+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    panel_pump_test.go:389: fold-key finding did not reproduce on this run: the two spellings came out identical
--- PASS: TestSnapshotWorkspaceSectionReportsTheNarrowing (1.00s)
=== RUN   TestThePumpIsDrivenNotJustAssembled
time=2026-10-09T14:37:09.437+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:37:09 mockllm: serving on http://127.0.0.1:62347 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:37:09.457+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2599296732\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:37:09.477+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2599296732\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:37:09.486+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:37:09Z duration_ms=8
time=2026-10-09T14:37:09.493+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestThePumpIsDrivenNotJustAssembled2599296732\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:37:09.498+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:37:09Z duration_ms=4
time=2026-10-09T14:37:09.502+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:37:09.505+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestThePumpIsDrivenNotJustAssembled (1.04s)
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/the_configured_pair_goes_out_as_it_stands
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/height_0_keeps_the_host's_260,_the_same_rule_the_create_uses
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/width_0_falls_back_to_the_host_constant,_not_to_a_guess
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/an_unreadable_config_answers_the_host_constants
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_negative_answer_resolves_the_way_the_create_resolves_it
=== RUN   TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_host_nobody_sized_sends_nothing_at_all
--- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/the_configured_pair_goes_out_as_it_stands (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/height_0_keeps_the_host's_260,_the_same_rule_the_create_uses (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/width_0_falls_back_to_the_host_constant,_not_to_a_guess (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/an_unreadable_config_answers_the_host_constants (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_negative_answer_resolves_the_way_the_create_resolves_it (0.00s)
    --- PASS: TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch/a_host_nobody_sized_sends_nothing_at_all (0.00s)
=== RUN   TestTicket255r1ShowOnAnExistingWindowSendsTheResize
--- PASS: TestTicket255r1ShowOnAnExistingWindowSendsTheResize (0.00s)
=== RUN   TestTicket255r1SinkTellsDispatchFromABareCall
--- PASS: TestTicket255r1SinkTellsDispatchFromABareCall (0.00s)
=== RUN   TestTicket255r1ReshowWithNoLiveControlSendsNothing
--- PASS: TestTicket255r1ReshowWithNoLiveControlSendsNothing (0.00s)
=== RUN   TestTicket255r1ShowAsksForGeometryOnlyOnTheReshowBranch
--- PASS: TestTicket255r1ShowAsksForGeometryOnlyOnTheReshowBranch (0.00s)
=== RUN   TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne
--- PASS: TestTicket255r1HostResizeCallSiteAsksForHintNoneAndIsTheOnlyOne (0.00s)
=== RUN   TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe
    panel_resident_windows_test.go:315: AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]
time=2026-10-09T14:37:09.538+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:325: no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all). AC#13 cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T14:37:29.544+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:37:29.544+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.01s)
=== RUN   TestAC13BringUpSurvivesAReusedThreadQuit
    panel_resident_windows_test.go:483: AC#13 reused-thread root cause: planted one WM_QUIT on this locked thread, then ran bringUp - panicked=<nil> err=<nil> created=true
    panel_resident_windows_test.go:485: AC#13 reused-thread release: tid=27444 dispatched 3 message(s) before unlocking; windows left on that thread=0 queue head=empty
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (5.03s)
=== RUN   TestAC13BringUpRefusesAThreadWithAQueuedClose
    panel_resident_windows_test.go:569: AC#13 queued-close: tid=27444 first bringUp ok, pumped 0 to idle, Destroy left a queued close (hwnd 0x9F0EFA, plantQueued=true), then bringUp#2 err=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x9F0EFA. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message created=false; the close was still queued after the refusal=true; released with windows=0 queue head=empty
--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (5.03s)
=== RUN   TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever
time=2026-10-09T14:37:39.604+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
    panel_resident_windows_test.go:689: fifth-shape plant on the panel thread: hwnd=0xA00EFA live=true PostMessageW(WM_CLOSE) returned=true lastErr=The operation completed successfully. | check reads queued=true names-the-live-window=true
time=2026-10-09T14:37:39.644+08:00 level=ERROR msg="panel host: show failed on the panel thread" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xA00EFA. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
wisp: panel could not open (panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xA00EFA. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message): 33r9-nail
time=2026-10-09T14:37:39.644+08:00 level=ERROR msg="panel thread: retiring without a panel window after a named refusal" via=33r9-nail err="panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xA00EFA. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message" shows=1
wisp: panel thread will take no further requests: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xA00EFA. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message
time=2026-10-09T14:37:39.644+08:00 level=INFO msg="panel thread ending" why="no window to pump" window_opened=false
    panel_resident_windows_test.go:759: refusal-to-retirement: settle=5.4542ms second RequestShow accepted=false in 0s | thread retired=true isFinished=true | the planted live window 0xA00EFA IsWindow=true | startUpErr=panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xA00EFA. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message | statusLine="panel thread could not start: panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0xA00EFA. The library's create pump would dispatch it before the control exists and end its own GetMessageW loop (pkg/edge chromium.go:96-111, webview.go:242-243 + 381-383); measured 3/3 panic on a planted thread. The thread's owner must dispatch its queue to empty - purging without dispatching leaks the prior window and can swallow this create's own completion message"
time=2026-10-09T14:37:39.649+08:00 level=INFO msg="panel thread exited cleanly" shows=2 toggles=0 disposals=0
--- PASS: TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever (0.04s)
=== RUN   TestAC14AwaitedBindingReplyReachesThePage
time=2026-10-09T14:37:39.649+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T14:37:59.665+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:37:59.665+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.02s)
=== RUN   TestAC14GoSideEvalPushReachesThePage
time=2026-10-09T14:37:59.665+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
    panel_resident_windows_test.go:866: no report "ac14-push" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's push hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
time=2026-10-09T14:38:19.673+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:38:19.673+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- FAIL: TestAC14GoSideEvalPushReachesThePage (20.01s)
=== RUN   TestPanelThreadIsSTAAndExitsCleanly
time=2026-10-09T14:38:19.673+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T14:38:24.680+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:38:24.680+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
    panel_resident_windows_test.go:910: panel thread up and down: hwnd 0x43f08ac, exits observed, shows=1
time=2026-10-09T14:38:24.680+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (5.01s)
=== RUN   TestPanelThreadNameIsNotInResidentRoster
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
=== RUN   TestBallPanelGesturesReachThePanelThread
time=2026-10-09T14:38:24.682+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (panel-hotkey, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T14:38:29.699+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:38:29.699+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=2 disposals=0
--- PASS: TestBallPanelGesturesReachThePanelThread (5.02s)
=== RUN   TestBallGestureWithoutPanelHostStillRecords
time=2026-10-09T14:38:29.700+08:00 level=WARN msg="ball gesture arrived with no executor" gesture=panel-hotkey why="this process has no task pipeline and no microphone, so the gesture has no executor here; the gestures that DO have one are the cancel key the assembly root wired (ticket 246) and the two panel gestures that reach the resident panel thread (ticket 33)"
wisp: ball panel-hotkey: this process has no task pipeline and no microphone, so the gesture has no executor here; the gestures that DO have one are the cancel key the assembly root wired (ticket 246) and the two panel gestures that reach the resident panel thread (ticket 33)
--- PASS: TestBallGestureWithoutPanelHostStillRecords (0.00s)
=== RUN   TestAC4PriorFocusSurvivesARefusedPanelSample
time=2026-10-09T14:38:29.701+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=panel-sta owner="panel host (ticket 33)"
wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)
time=2026-10-09T14:38:34.778+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-10-09T14:38:34.779+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (5.08s)
=== RUN   TestPagePostMessageEnvelopeReachesDispatchRawViaTransport
--- PASS: TestPagePostMessageEnvelopeReachesDispatchRawViaTransport (0.00s)
=== RUN   TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3
    panel_transport_35r2_test.go:1755: shape ③ delivery: nativeExitCalls=1 receiverOK=1 receiverLost=0 ignoredWrites=(none) doorRounds=1 maxPostDepth=3 capTripped=false capDepth=0 unparsable=0 unboundSlots=0 evalThrew=0 resolved=1 rejected=0 scriptThrows=[] | frame="{\"id\":1,\"method\":\"wispDispatch\",\"params\":[\"{\\\"method\\\":\\\"panel.mode.request\\\",\\\"requestId\\\":\\\"pc-1-3f2b1c0d-9e7a-4c1b-8f14-e45fceea469a\\\",\\\"source\\\":\\\"panel-composer\\\",\\\"to\\\":\\\"ask_every_step\\\"}\"]}" | door reply=""
--- PASS: TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3 (0.00s)
=== RUN   TestShapeA3SecondPagePostStillDelivers
--- PASS: TestShapeA3SecondPagePostStillDelivers (0.00s)
=== RUN   TestShapeA3ForwardingHookIsIdempotentInOneDocument
--- PASS: TestShapeA3ForwardingHookIsIdempotentInOneDocument (0.00s)
=== RUN   TestLegacySubShapeOneHookDiesInAReentryLoop
    panel_transport_35r2_test.go:1844: sub-shape ① under the behavioural yard: re-entered 9 levels (cap 8), native exit called 0 times, door fired 0 times, page error="RangeError: maximum call stack size exceeded - chrome.webview.postMessage re-entered 9 levels (cap 8) with the native exit called 0 times; a page in a real browser dies here"
--- PASS: TestLegacySubShapeOneHookDiesInAReentryLoop (0.00s)
=== RUN   TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot
--- PASS: TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot (0.00s)
=== RUN   TestForwardingHookMustNotLoseTheNativeExitReceiver
    panel_transport_35r2_test.go:1968: M-A face has teeth: 1 native-exit call, all with chrome.webview as receiver (receiverLost=0); a bare native(message) trips it
--- PASS: TestForwardingHookMustNotLoseTheNativeExitReceiver (0.00s)
=== RUN   TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit
    panel_transport_35r2_test.go:2018: M-B1 face has teeth: non-writable world refused the hook (doorRounds=0, unboundSlots=1), writable world armed it (doorRounds=1)
--- PASS: TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit (0.00s)
=== RUN   TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent
    panel_transport_35r2_test.go:2071: door absent: the guard's typeof half routed the page's envelope to the native exit byte-for-byte (frames=1, no throw)
--- PASS: TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent (0.00s)
=== RUN   TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable
    panel_transport_35r2_test.go:2140: M-B3 face has teeth: door present as typeof==="string", guard's typeof half routed the page's envelope to the native exit byte-for-byte (frames=1, no throw)
--- PASS: TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable (0.00s)
=== RUN   TestProvidersProbeRecordsMeasuredThinkingFalse
time=2026-10-09T14:38:35.851+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:38:35 mockllm: serving on http://127.0.0.1:57524 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:38:35.882+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse3581396263\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:38:35.889+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:38:35Z duration_ms=6
time=2026-10-09T14:38:35.896+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingFalse3581396263\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:38:35.901+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:38:35Z duration_ms=4
time=2026-10-09T14:38:35.905+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:38:35.906+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:38:35.923+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingFalse (1.22s)
=== RUN   TestProvidersProbeRecordsMeasuredThinkingTrue
time=2026-10-09T14:38:36.905+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:38:36 mockllm: serving on http://127.0.0.1:50944 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:38:36.933+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue1968847840\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:38:36.939+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:38:36Z duration_ms=5
time=2026-10-09T14:38:36.945+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeRecordsMeasuredThinkingTrue1968847840\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:38:36.950+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:38:36Z duration_ms=4
time=2026-10-09T14:38:36.954+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:38:36.956+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:38:36.968+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeRecordsMeasuredThinkingTrue (0.99s)
=== RUN   TestProvidersDiscoverListsWhatTheServerServes
time=2026-10-09T14:38:37.897+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:38:37 mockllm: serving on http://127.0.0.1:59416 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:38:37.925+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes2161520679\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:38:37.931+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:38:37Z duration_ms=5
time=2026-10-09T14:38:37.938+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersDiscoverListsWhatTheServerServes2161520679\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:38:37.943+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:38:37Z duration_ms=4
time=2026-10-09T14:38:37.948+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:38:37.950+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersDiscoverListsWhatTheServerServes (0.97s)
=== RUN   TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless
time=2026-10-09T14:38:38.875+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:38:38 mockllm: serving on http://127.0.0.1:55504 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:38:38.906+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless268269277\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:38:38.913+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:38:38Z duration_ms=6
time=2026-10-09T14:38:38.919+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless268269277\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:38:38.925+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:38:38Z duration_ms=5
time=2026-10-09T14:38:38.929+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:38:38.931+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless (0.98s)
=== RUN   TestAC246CancelGestureUsesTheInjectedExecutor
time=2026-10-09T14:38:38.945+08:00 level=INFO msg="ball: the cancel key was handled by the injected approval gate" outcome="injected executor ran"
wisp: injected executor ran
time=2026-10-09T14:38:38.945+08:00 level=WARN msg="cancel hotkey fired with no executor" why="the assembly root injected no approval gate into this leg, so the borrow cannot be spent"
wisp: ball cancel-hotkey: no approval gate was injected into this process, so the press decided nothing
--- PASS: TestAC246CancelGestureUsesTheInjectedExecutor (0.00s)
=== RUN   TestAC246EscChannelStaysUnloadedWithoutABallWindow
time=2026-10-09T14:38:38.945+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:38.946+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246EscChannelStaysUnloadedWithoutABallWindow (0.00s)
=== RUN   TestAC246CardWithNoWindowFailsClosedThroughTheRealGate
time=2026-10-09T14:38:38.946+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:38.946+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
time=2026-10-09T14:38:38.946+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host:246-no-window tool=resident.confirmation level=L1
time=2026-10-09T14:38:38.946+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host:246-no-window tool=resident.confirmation: 常驻进程没有悬浮球窗口，卡片无处呈现
--- PASS: TestAC246CardWithNoWindowFailsClosedThroughTheRealGate (0.00s)
=== RUN   TestAC246CancelStepHookRunsOnTheRealShutdownSequence
time=2026-10-09T14:38:38.948+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:38.948+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:38.948+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:38.948+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T14:38:38.948+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T14:38:38.948+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:38.948+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:38.949+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC246CancelStepHookRunsOnTheRealShutdownSequence (0.00s)
=== RUN   TestAC246ShippedResidentProcessOwnsItsCancelStep
    resident_approval_246_windows_test.go:244: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentProcessOwnsItsCancelStep1187512521\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentProcessOwnsItsCancelStep (3.10s)
=== RUN   TestAC246ChannelNeedsBothWindowAndExecutor
time=2026-10-09T14:38:42.051+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.285+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:42.297+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance="no host config view" hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:38:42.297+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
wisp: [audit] resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）
time=2026-10-09T14:38:42.301+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
time=2026-10-09T14:38:42.301+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.302+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246ChannelNeedsBothWindowAndExecutor (0.25s)
=== RUN   TestAC246VetoSentenceWithNoCard
time=2026-10-09T14:38:42.302+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestAC246VetoSentenceWithNoCard (0.00s)
=== RUN   TestAC246StatusLineSaysWhatTheLegDoesNot
time=2026-10-09T14:38:42.302+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.303+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
wisp: [audit] resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）
--- PASS: TestAC246StatusLineSaysWhatTheLegDoesNot (0.00s)
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all
time=2026-10-09T14:38:42.303+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing
time=2026-10-09T14:38:42.304+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta110443446\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta110443446\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:38:42.304+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta110443446\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.305+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConsta110443446\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.306+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.306+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/no_host_view_at_all (0.00s)
    --- PASS: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants/config.toml_missing (0.00s)
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant
time=2026-10-09T14:38:42.308+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowtime1302333866\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading
time=2026-10-09T14:38:42.310+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind1927264329\001\config.toml window_sec_read=2 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded
time=2026-10-09T14:38:42.313+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowboth2067179699\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T14:38:42.315+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind1040198765\001\config.toml window_sec_read=99 confirm_timeout_sec_read=300 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive
time=2026-10-09T14:38:42.318+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindowwind2099362038\001\config.toml window_sec_read=1 confirm_timeout_sec_read=300 gate_window=2s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow (0.01s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/timeout_only_-_the_window_then_comes_from_the_schema_tag,_not_from_the_gate's_compiled_constant (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_seeded_to_2s_-_differs_from_the_pre-256_compiled_3s,_so_it_is_a_real_reading (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/both_seeded (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_way_too_large_-_reverse_control,_the_clamp_must_survive (0.00s)
    --- PASS: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow/window_below_the_floor_-_reverse_control,_the_clamp_must_survive (0.00s)
=== RUN   TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly
time=2026-10-09T14:38:42.320+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly3997344733\001\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.321+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly3997344733\001\config.toml window_sec_read=2 confirm_timeout_sec_read=90 gate_window=2s gate_queue_timeout=1m30s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly (0.00s)
=== RUN   TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims
    resident_approval_risk_256_windows_test.go:494: landing site reading: resident leg passes 6 of 10 declared approval.Options fields (was 3 of 10 before ticket 256)
--- PASS: TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims (0.00s)
=== RUN   TestTicket256ResidentBootPassesTheDataDirToTheGate
--- PASS: TestTicket256ResidentBootPassesTheDataDirToTheGate (0.00s)
=== RUN   TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig
time=2026-10-09T14:38:42.325+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3267087739\001\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3267087739\\001\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:38:42.325+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3267087739\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.327+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3267087739\002\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
wisp: resident [risk]: config.toml is present but was refused at load (config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]); the approval gate falls back to the compiled constants (DefaultApprovalTimeout=300s / DefaultL1Window=3s), so the [risk] numbers written in that file are NOT the numbers this process runs on
time=2026-10-09T14:38:42.327+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3267087739\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.328+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance=config config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfi3267087739\002\config.toml window_sec_read=2 confirm_timeout_sec_read=45 gate_window=2s gate_queue_timeout=45s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig (0.00s)
=== RUN   TestTicket268RefusedRiskConfigReachesStdoutOnce
time=2026-10-09T14:38:42.330+08:00 level=WARN msg="resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2166675832\001\config.toml err="config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]" provenance="defaults (config.toml present but refused at load)" fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:38:42.331+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml present but refused at load)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2166675832\001\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:42.331+08:00 level=WARN msg="resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants" path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2166675832\002\config.toml err="config: config.toml read: no file at this path yet。第 1 种拒因：文件没建：这一页与这条链都不新建 config.toml。在控制台运行一次 wisp run，第一次启动会写出全默认的首份配置（它不替你选任何模型，也不替你建任何服务商行）: open C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket268RefusedRiskConfigReachesStdoutOnce2166675832\\002\\config.toml: The system cannot find the file specified." fallback="DefaultApprovalTimeout=300s / DefaultL1Window=3s"
time=2026-10-09T14:38:42.331+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (config.toml unreadable)" config_path=C:\Users\swq\AppData\Local\Temp\TestTicket268RefusedRiskConfigReachesStdoutOnce2166675832\002\config.toml window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket268RefusedRiskConfigReachesStdoutOnce (0.00s)
=== RUN   TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords
--- PASS: TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords (0.00s)
=== RUN   TestAC247LiveMicrophoneLevelsReachTheBallSeam
    resident_audio_247_live_windows_test.go:130: AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0
--- SKIP: TestAC247LiveMicrophoneLevelsReachTheBallSeam (0.00s)
=== RUN   TestAC247ShippedDefaultsAreTheOnesThisLegReads
--- PASS: TestAC247ShippedDefaultsAreTheOnesThisLegReads (0.00s)
=== RUN   TestAC247VoiceDisabledBuildsNoCollector
time=2026-10-09T14:38:42.337+08:00 level=INFO msg="audio: capture leg not built" reason="voice.enabled=false" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector3313696996\001\config.toml
wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 C:\Users\swq\AppData\Local\Temp\TestAC247VoiceDisabledBuildsNoCollector3313696996\001\config.toml）：球不会收到任何电平，本进程其余部分照常
time=2026-10-09T14:38:42.339+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.339+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.339+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.339+08:00 level=INFO msg="shutdown step skipped (module not present)" step=4 name=stop-audio
time=2026-10-09T14:38:42.339+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.339+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.339+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247VoiceDisabledBuildsNoCollector (0.01s)
=== RUN   TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice
time=2026-10-09T14:38:42.354+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice1853808174\001\config.toml path=T
time=2026-10-09T14:38:42.355+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.355+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.355+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.355+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:38:42.355+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.355+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.355+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice (0.01s)
=== RUN   TestAC247UnmuteReachesTheBallSeam
time=2026-10-09T14:38:42.361+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247UnmuteReachesTheBallSeam1185142444\001\config.toml path=T
time=2026-10-09T14:38:42.362+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.362+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.362+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.362+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=2 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:38:42.362+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.362+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.362+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247UnmuteReachesTheBallSeam (0.01s)
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied
time=2026-10-09T14:38:42.367+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
time=2026-10-09T14:38:42.367+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T14:38:42.367+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.368+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.368+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.368+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Busy Mic (0x8889000A): capture device is held in exclusive mode by another application; close that application or pick another device
time=2026-10-09T14:38:42.368+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.368+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.368+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied
time=2026-10-09T14:38:42.373+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
time=2026-10-09T14:38:42.373+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T14:38:42.373+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.373+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.373+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.373+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Built-in Mic (0x80070005): microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone
time=2026-10-09T14:38:42.373+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.373+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.373+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
=== RUN   TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device
time=2026-10-09T14:38:42.380+08:00 level=ERROR msg="audio source failed" source=gate-T error="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
time=2026-10-09T14:38:42.380+08:00 level=ERROR msg="audio: capture device unavailable" class=audio_device err="audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)" posture="boot continues, no state pushed"
wisp: 麦克风不可用（错误分类 audio_device）：audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
time=2026-10-09T14:38:42.380+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.380+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.380+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.380+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0 last_error=audio_device: capture open on device Gone Headset (0x88890004): device invalidated (unplugged or disabled)
time=2026-10-09T14:38:42.380+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.380+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.380+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss (0.03s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/occupied (0.01s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/permission_denied (0.01s)
    --- PASS: TestAC247DeviceFailureShapesStillBootAndSayTheLoss/no_device (0.02s)
=== RUN   TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder
time=2026-10-09T14:38:42.397+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder988287802\001\config.toml path=T
time=2026-10-09T14:38:42.397+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.397+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.397+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.397+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:38:42.397+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.398+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.398+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder (0.01s)
=== RUN   TestAC247HandingTheLevelToTheSeamIsNotVisibility
time=2026-10-09T14:38:42.404+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Local\Temp\TestAC247HandingTheLevelToTheSeamIsNotVisibility1114089406\001\config.toml path=T
    resident_audio_247_windows_test.go:372: AC#10 reading: PrototypeVisualsEnabled()=false levels_reaching_the_ball_seam=1 (a number arriving is not a pixel moving)
time=2026-10-09T14:38:42.405+08:00 level=INFO msg="shutdown step skipped (module not present)" step=1 name=scheduler-close
time=2026-10-09T14:38:42.405+08:00 level=INFO msg="shutdown step skipped (module not present)" step=2 name=stop-hotkey-kws
time=2026-10-09T14:38:42.405+08:00 level=INFO msg="shutdown step skipped (module not present)" step=3 name=cancel-task-roots
time=2026-10-09T14:38:42.405+08:00 level=INFO msg="audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0"
wisp: audio capture leg stopped: levels_delivered=1 frames_sent=0 frames_dropped=0 reopens=0
time=2026-10-09T14:38:42.405+08:00 level=INFO msg="shutdown step skipped (module not present)" step=5 name=release-speech-sessions
time=2026-10-09T14:38:42.405+08:00 level=INFO msg="shutdown step skipped (module not present)" step=6 name=destroy-panel-webview
time=2026-10-09T14:38:42.405+08:00 level=INFO msg="shutdown step skipped (module not present)" step=7 name=flush-logs-close-db
--- PASS: TestAC247HandingTheLevelToTheSeamIsNotVisibility (0.01s)
=== RUN   TestAC228ResidentLegIsTheBallHost
    resident_ball_228_test.go:214: ball host functions (production): resident_ball_windows.go:startResidentBall
--- PASS: TestAC228ResidentLegIsTheBallHost (0.02s)
=== RUN   TestAC228BallHostAnswersEveryGesture
    resident_ball_228_test.go:323: all 10 gesture callbacks answered by the resident host
--- PASS: TestAC228BallHostAnswersEveryGesture (0.02s)
=== RUN   TestAC228ResidentLegReportsAndBooksItsBall
    resident_ball_228_windows_test.go:104: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ResidentLegReportsAndBooksItsBall1339328699\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:168: AC#1 READING: console up=1 absent=0; records created=9 refused=-1 stopped=14 install=1
--- PASS: TestAC228ResidentLegReportsAndBooksItsBall (3.15s)
=== RUN   TestAC228ExitRequestDuringBootStillLeavesThroughD38E
    resident_ball_228_windows_test.go:211: resident sink C:\Users\swq\AppData\Local\Temp\TestAC228ExitRequestDuringBootStillLeavesThroughD38E859248311\002\logs: 1 file(s), 22 record(s)
    resident_ball_228_windows_test.go:221: AC#1 READING: boot-time break exited clean; install=1 shutdown trail starts at 15 of 22 record(s); ball records created=9 stopped=14
--- PASS: TestAC228ExitRequestDuringBootStillLeavesThroughD38E (3.00s)
=== RUN   TestTicket260R4ShippedConstructorInstallsTheReader
time=2026-10-09T14:38:48.596+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.596+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R4ShippedConstructorInstallsTheReader (0.00s)
=== RUN   TestTicket260R4NoBallHostNamesNoKeyOnTheCard
time=2026-10-09T14:38:48.596+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.596+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T14:38:48.597+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.597+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
--- PASS: TestTicket260R4NoBallHostNamesNoKeyOnTheCard (0.00s)
=== RUN   TestTicket260R4SeamIsNotAProductionShortcut
time=2026-10-09T14:38:48.597+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R4SeamIsNotAProductionShortcut (0.00s)
=== RUN   TestTicket260R3DefaultWordingIsTheOldSentence
time=2026-10-09T14:38:48.597+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.597+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.597+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3DefaultWordingIsTheOldSentence (0.00s)
=== RUN   TestTicket260R3SeededKeyReplacesEsc
time=2026-10-09T14:38:48.598+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.598+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.598+08:00 level=INFO msg="audit: resident-approval: 审批门已装配进常驻进程，取消通道 Ctrl+Alt+Q 已加载（本票只落 Ctrl+Alt+Q 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）"
--- PASS: TestTicket260R3SeededKeyReplacesEsc (0.00s)
=== RUN   TestTicket260R3UnloadBranchesNameNoKey
time=2026-10-09T14:38:48.598+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.598+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T14:38:48.598+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.599+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
time=2026-10-09T14:38:48.599+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.599+08:00 level=INFO msg="audit: resident-approval: 审批门已装配，但本进程没有悬浮球窗口，四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"
time=2026-10-09T14:38:48.599+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.599+08:00 level=INFO msg="audit: resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"
--- PASS: TestTicket260R3UnloadBranchesNameNoKey (0.00s)
=== RUN   TestTicket260R3CancelKeyComesFromTheBallChain
time=2026-10-09T14:38:48.600+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R3CancelKeyComesFromTheBallChain (0.00s)
=== RUN   TestTicket260R3ProductionPathHasNoSeam
time=2026-10-09T14:38:48.600+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.600+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:38:48.600+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket260R3ProductionPathHasNoSeam (0.00s)
=== RUN   TestTicket265ResidentGrantHolderUnboundFailsLoudly
--- PASS: TestTicket265ResidentGrantHolderUnboundFailsLoudly (0.00s)
=== RUN   TestTicket265ResidentGrantHolderBoundDelegatesEveryField
--- PASS: TestTicket265ResidentGrantHolderBoundDelegatesEveryField (0.00s)
=== RUN   TestTicket265LedgerErrorReachesTheOperatorThroughTheHolder
--- PASS: TestTicket265LedgerErrorReachesTheOperatorThroughTheHolder (0.00s)
=== RUN   TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow
--- PASS: TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow (0.00s)
=== RUN   TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed
--- PASS: TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed (0.00s)
=== RUN   TestTicket265GateLiteralCarriesTheResidentHolder
--- PASS: TestTicket265GateLiteralCarriesTheResidentHolder (0.00s)
=== RUN   TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask
--- PASS: TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask (0.00s)
=== RUN   TestTicket265SessionLedgerStillHasOneProductionConstructionSite
--- PASS: TestTicket265SessionLedgerStillHasOneProductionConstructionSite (0.01s)
=== RUN   TestTicket265ResidentApprovalConstructsAnUnboundHolder
time=2026-10-09T14:38:48.609+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
--- PASS: TestTicket265ResidentApprovalConstructsAnUnboundHolder (0.00s)
=== RUN   Test258AssemblyRootWiresTheChainAndTheBridge
--- PASS: Test258AssemblyRootWiresTheChainAndTheBridge (0.00s)
=== RUN   Test258SinkRebindRulerReadsBothSpellings
--- PASS: Test258SinkRebindRulerReadsBothSpellings (0.00s)
=== RUN   Test258ConstructionChainTakesConfigValues
--- PASS: Test258ConstructionChainTakesConfigValues (0.01s)
=== RUN   Test258ConstructionChainMissingFileFallsBackAndSaysIt
--- PASS: Test258ConstructionChainMissingFileFallsBackAndSaysIt (0.00s)
=== RUN   Test258SectionMissingTierWordIsDefaults
--- PASS: Test258SectionMissingTierWordIsDefaults (0.00s)
=== RUN   Test258ProvenanceWordsAreThePrintedOnes
--- PASS: Test258ProvenanceWordsAreThePrintedOnes (0.00s)
=== RUN   Test258BridgeRebindsLiveKeysFromConfigEdit
time=2026-10-09T14:38:48.650+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:48.654+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T14:38:48.654+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:38:48.659+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:48.659+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+R mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T14:38:48.664+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeRebindsLiveKeysFromConfigEdit (0.04s)
=== RUN   Test258BridgeMutationNoSrcKeepsOldBinding
time=2026-10-09T14:38:48.687+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:48.691+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:38:49.109+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeMutationNoSrcKeepsOldBinding (0.45s)
=== RUN   Test258OccupiedCombinationNamesTheNewValue
time=2026-10-09T14:38:49.141+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:49.145+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T14:38:49.145+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:38:49.150+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:49.150+08:00 level=ERROR msg="hotkey occupied by another program, not registered" hotkey=panel binding=Ctrl+Alt+V err="Hot key is already registered."
time=2026-10-09T14:38:49.150+08:00 level=ERROR msg="ball: hotkey reload" binding="hotkey panel = \"Ctrl+Alt+V\" is occupied by another program and was NOT registered; pressing it will do nothing until you pick a free combination in [hotkey]"
time=2026-10-09T14:38:49.150+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+Z mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+V live=2
time=2026-10-09T14:38:49.156+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258OccupiedCombinationNamesTheNewValue (0.05s)
=== RUN   Test258V1ProbeSummonEditRebindsLiveBall
time=2026-10-09T14:38:49.180+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:49.184+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T14:38:49.185+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the cancel key: this leg has no task pipeline and no microphone, and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
time=2026-10-09T14:38:49.189+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T14:38:49.189+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+7 mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T14:38:49.194+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (0.04s)
=== RUN   Test258V1ProbeConstructionMissingFileNamesTheFallback
--- PASS: Test258V1ProbeConstructionMissingFileNamesTheFallback (0.01s)
=== RUN   TestAC1ResidentLegInstallsItsLogListenerOnDisk
    resident_sink_nail_127_windows_test.go:433: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegInstallsItsLogListenerOnDisk2054215221\002\logs: 1 file(s), 14 record(s)
--- PASS: TestAC1ResidentLegInstallsItsLogListenerOnDisk (3.54s)
=== RUN   TestAC1ResidentLegOutlivesItsOwnLogFailure
--- PASS: TestAC1ResidentLegOutlivesItsOwnLogFailure (3.28s)
=== RUN   TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink
    resident_sink_nail_127_windows_test.go:594: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink2777802722\002\logs: 1 file(s), 22 record(s)
    resident_sink_nail_127_windows_test.go:640: RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7], first="shutdown step skipped (module not present)" last="shutdown step skipped (module not present)"
    resident_sink_nail_127_windows_test.go:571: resident leg exit: exited 0
--- PASS: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (3.43s)
=== RUN   TestAC246ResidentPipelineAsksThroughTheOneGate
time=2026-10-09T14:39:00.441+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:39:00 mockllm: serving on http://127.0.0.1:59119 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:39:00.456+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:39:00.483+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate2467830387\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:39:00.495+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:39:00Z duration_ms=12
time=2026-10-09T14:39:00.503+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentPipelineAsksThroughTheOneGate2467830387\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:39:00.507+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:39:00Z duration_ms=4
time=2026-10-09T14:39:00.512+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:39:00.513+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:39:00.513+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_205e6b85411ba37d323304dc7a5ac9f6 (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T14:39:00.515+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate2467830387\\\\002\\\\config.toml\""
time=2026-10-09T14:39:00.515+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentPipelineAsksThroughTheOneGate2467830387\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T14:39:00.523+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T14:39:00.523+08:00 level=ERROR msg="approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝" corr=host-corr-1 tool=fs.write level=L1
time=2026-10-09T14:39:00.523+08:00 level=INFO msg="audit: approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现"
wisp: [audit] approval: L1 prompt failed corr=host-corr-1 tool=fs.write: 常驻进程没有悬浮球窗口，卡片无处呈现
time=2026-10-09T14:39:00.524+08:00 level=INFO msg="audit: tools: call kind=refused task=host:246-one-gate corr=host-corr-1 tool=fs.write risk=L1 decision=reject outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T14:39:00.524+08:00 level=INFO msg="audit: tools: PATH-ACCOUNT task=host:246-one-gate tool=fs.write roots=1 rewritten=[] unusable=[]"
time=2026-10-09T14:39:00.527+08:00 level=INFO msg="audit: wisp run: 风险判定 tool=fs.write level=L1 rules=[R1] reason=\"R1: 工具声明为下界（L1）\""
time=2026-10-09T14:39:00.527+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentPipelineAsksThroughTheOneGate (1.08s)
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled
time=2026-10-09T14:39:01.532+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:39:01 mockllm: serving on http://127.0.0.1:59125 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:39:01.548+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger
time=2026-10-09T14:39:01.572+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled2011466727\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:39:01.578+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:39:01Z duration_ms=5
time=2026-10-09T14:39:01.592+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentGateInjectionIsRefusedHalfAssembled2011466727\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:39:01.597+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:39:01Z duration_ms=5
time=2026-10-09T14:39:01.601+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:39:01.604+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:39:01.604+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_b7a38b3587c6b5a401d9c268590fa8b6 (结束点＝本进程退出，A435 第 2 条)"
=== RUN   TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface
time=2026-10-09T14:39:01.618+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:39:01.618+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_5daa510ee2ad43efef21cda3b530318c (结束点＝本进程退出，A435 第 2 条)"
--- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled (1.09s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_ledger (0.06s)
    --- PASS: TestAC246ResidentGateInjectionIsRefusedHalfAssembled/gate_without_surface (0.01s)
=== RUN   TestAC246ResidentTaskRootCancelStopsTheModelCall
time=2026-10-09T14:39:02.647+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:39:02 mockllm: serving on http://127.0.0.1:59126 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:39:02.664+08:00 level=INFO msg="resident gate: [risk] tier taken at construction" provenance="defaults (no host config view)" config_path="(no host config view)" window_sec_read=0 confirm_timeout_sec_read=0 gate_window=3s gate_queue_timeout=5m0s scope="construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)"
time=2026-10-09T14:39:02.688+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall3376463094\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:39:02.701+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:39:02Z duration_ms=12
time=2026-10-09T14:39:02.708+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAC246ResidentTaskRootCancelStopsTheModelCall3376463094\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:39:02.714+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:39:02Z duration_ms=4
time=2026-10-09T14:39:02.717+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:39:02.719+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:39:02.719+08:00 level=INFO msg="audit: wisp run: SESSION-MINT id=sess_44ae75e9c825b03b5b456a0e8078edc4 (结束点＝本进程退出，A435 第 2 条)"
time=2026-10-09T14:39:02.719+08:00 level=INFO msg="audit: perm: MODE-READ origin=startup mode=ask_every_step source=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall3376463094\\\\002\\\\config.toml\""
time=2026-10-09T14:39:02.720+08:00 level=INFO msg="audit: config: HOT-RELOAD state=armed tick=1s path=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall3376463094\\\\002\\\\config.toml\" goroutine=watchdog owner=config d36_confirm=on restart_notice=on detail=\"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效\""
time=2026-10-09T14:39:02.720+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall3376463094\\\\002\" enabled=true budget=2300"
time=2026-10-09T14:39:02.720+08:00 level=INFO msg="audit: projctx: projctx: 没有找到任何项目说明文件（工作区逐级向上与数据目录都查过）"
time=2026-10-09T14:39:02.725+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T06:39:02Z depth=0 pending= mode=ask_every_step ws=unset results=1/1 done bytes=1910 sha256=5cd90806144e8f37"
time=2026-10-09T14:39:02.725+08:00 level=INFO msg="audit: tools: C25 scope closed task=327fffc9-1256-4ae8-b7ed-28a4f439a122 was_open=false dropped=0 open_scopes=0 close_err=<nil>"
    resident_task_source_246_windows_test.go:199: AC#7 positive control: live root -> execute code 0, provider chat requests 1
time=2026-10-09T14:39:02.730+08:00 level=INFO msg="audit: resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留"
wisp: [audit] resident-approval: 退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张、路由失败 0 张；任务根已取消，无等待残留
time=2026-10-09T14:39:02.730+08:00 level=INFO msg="audit: wisp run: 项目说明加载器 workspace=\"\" globalDir=\"C:\\\\Users\\\\swq\\\\AppData\\\\Local\\\\Temp\\\\TestAC246ResidentTaskRootCancelStopsTheModelCall3376463094\\\\002\" enabled=true budget=2300"
time=2026-10-09T14:39:02.731+08:00 level=INFO msg="panel: SNAPSHOT at=2026-10-09T06:39:02Z depth=0 pending= mode=ask_every_step ws=unset results=2/2 done bytes=2358 sha256=d332ea9fc4077436"
time=2026-10-09T14:39:02.731+08:00 level=INFO msg="audit: tools: C25 scope closed task=134ac25a-596b-49f8-8bdf-c963b0ebc8c7 was_open=false dropped=0 open_scopes=0 close_err=<nil>"
time=2026-10-09T14:39:02.731+08:00 level=INFO msg="audit: wisp run: 后台任务的输出没有进名册：这个任务没有留下正文（空正文不进名册，免得「查不到」和「没打印」被读成同一件事）"
    resident_task_source_246_windows_test.go:225: AC#7 READING (ruling 2.3): after step 3's cancel, execute -> exit 1, provider requests 1 -> 1
time=2026-10-09T14:39:02.732+08:00 level=INFO msg="audit: config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=\"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到\""
--- PASS: TestAC246ResidentTaskRootCancelStopsTheModelCall (1.12s)
=== RUN   TestAC246TestTaskInjectionPredicate
=== RUN   TestAC246TestTaskInjectionPredicate/accepted_in_the_harness'_own_root
=== RUN   TestAC246TestTaskInjectionPredicate/refused_outside_test
=== RUN   TestAC246TestTaskInjectionPredicate/refused_when_the_root_is_not_the_harness'_own
=== RUN   TestAC246TestTaskInjectionPredicate/unset_is_silent
=== RUN   TestAC246TestTaskInjectionPredicate/a_refusal_never_carries_the_whole_text
--- PASS: TestAC246TestTaskInjectionPredicate (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/accepted_in_the_harness'_own_root (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/refused_outside_test (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/refused_when_the_root_is_not_the_harness'_own (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/unset_is_silent (0.00s)
    --- PASS: TestAC246TestTaskInjectionPredicate/a_refusal_never_carries_the_whole_text (0.00s)
=== RUN   TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry
    resident_task_source_246_windows_test.go:340: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry1269564032\002\logs: 1 file(s), 22 record(s)
--- PASS: TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry (3.09s)
=== RUN   TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline
    resident_task_source_246_windows_test.go:401: resident sink C:\Users\swq\AppData\Local\Temp\TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeli1025923129\002\logs: 1 file(s), 23 record(s)
--- PASS: TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline (3.37s)
=== RUN   TestAC246DevLegIgnoresTheTestTaskInjection
--- PASS: TestAC246DevLegIgnoresTheTestTaskInjection (3.47s)
=== RUN   TestTicket255RestartTierKeysAreBackedByATest
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.language
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.autostart
=== RUN   TestTicket255RestartTierKeysAreBackedByATest/app.single_instance
--- PASS: TestTicket255RestartTierKeysAreBackedByATest (0.11s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.language (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.autostart (0.03s)
    --- PASS: TestTicket255RestartTierKeysAreBackedByATest/app.single_instance (0.03s)
=== RUN   TestTicket101ManualSwitchSurvivesRestart
time=2026-10-09T14:39:13.753+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:39:13 mockllm: serving on http://127.0.0.1:58702 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:39:13.769+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart1873821577\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:39:13.801+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart1873821577\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:39:13.807+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:39:13Z duration_ms=6
time=2026-10-09T14:39:13.815+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart1873821577\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:39:13.820+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:39:13Z duration_ms=5
time=2026-10-09T14:39:13.824+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:39:13.826+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:39:13.829+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T14:39:13.831+08:00 level=INFO msg="config: wrote merged change into config.toml (ticket 226: only the keys listed changed value, every other key keeps the value the file already held; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved)" key=risk.permission_mode wrote=[risk.permission_mode]
time=2026-10-09T14:39:13.852+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ManualSwitchSurvivesRestart1873821577\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:39:13.865+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ManualSwitchSurvivesRestart (41.14s)
=== RUN   TestTicket101UntouchedConfigRestartsAtDefault
time=2026-10-09T14:39:54.844+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:39:54 mockllm: serving on http://127.0.0.1:58846 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:39:54.862+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1078290788\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:39:54.881+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1078290788\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:39:54.887+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:39:54Z duration_ms=6
time=2026-10-09T14:39:54.895+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1078290788\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:39:54.900+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:39:54Z duration_ms=5
time=2026-10-09T14:39:54.904+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:39:54.906+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:39:56.944+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1078290788\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:39:56.962+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:39:58.994+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UntouchedConfigRestartsAtDefault1078290788\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:39:59.007+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:381: the key is still absent after 3 cold starts, as expected:
        schema_version = 2
        
        [llm]
        text_chain = ["acme/m1"]
        
        [llm.retry]
        max = 1
        backoff_ms = 1
        
        [llm.providers.acme]
        protocol = "openai-chat"
        base_url = "http://127.0.0.1:58846/v1"
        api_key_ref = "dpapi:acme"
        
        [llm.providers.acme.models.m1]
        # 261-r2: explicit enabled, same reason as run_test.go's fixture comment.
        enabled = true
        context_window = 128000
        
        [fs]
        allowed_dirs = ["C:/Users/swq/AppData/Local/Temp/TestTicket101UntouchedConfigRestartsAtDefault1078290788/002"]
        
        [risk]
        l1_window_sec = 1
        confirm_timeout_sec = 40
--- PASS: TestTicket101UntouchedConfigRestartsAtDefault (7.13s)
=== RUN   TestTicket101SessionGrantDoesNotCrossRestart
time=2026-10-09T14:40:01.984+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:40:01 mockllm: serving on http://127.0.0.1:54225 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:40:02.001+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2633722347\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:40:02.029+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2633722347\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:40:02.035+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:40:02Z duration_ms=5
time=2026-10-09T14:40:02.043+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2633722347\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:40:02.049+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:40:02Z duration_ms=5
time=2026-10-09T14:40:02.057+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:40:02.059+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:40:42.148+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101SessionGrantDoesNotCrossRestart2633722347\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:40:42.162+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:22.219+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (81.18s)
=== RUN   TestTicket101ModeSwitchUsesTheRealL2Gate
time=2026-10-09T14:41:23.171+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:23 mockllm: serving on http://127.0.0.1:55176 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:23.189+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate280034808\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:23.210+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate280034808\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:23.216+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:23Z duration_ms=6
time=2026-10-09T14:41:23.223+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101ModeSwitchUsesTheRealL2Gate280034808\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:23.228+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:23Z duration_ms=4
time=2026-10-09T14:41:23.232+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:23.234+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket101ModeSwitchUsesTheRealL2Gate (1.33s)
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict
time=2026-10-09T14:41:24.459+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:24 mockllm: serving on http://127.0.0.1:55182 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:24.475+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict2430042328\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:24.497+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict2430042328\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:24.503+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:24Z duration_ms=6
time=2026-10-09T14:41:24.511+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict2430042328\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:24.516+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:24Z duration_ms=4
time=2026-10-09T14:41:24.520+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:24.522+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
    run_mode101_test.go:604: control log for reference (no fail-closed line expected here):
        wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。手改会按 D36 三档处理：可热加载段立即生效；[risk]/[fs]/[net]/[plugins] 的放宽要先答一张 L2 卡，不答按拒绝保留旧值；重启档的改动本次不生效，会另有一句告诉你为什么不生效。
        echo: ## scene
        当前时间：2026-10-09 06:41 +00:00
        wisp run: 任务 7882f7ad-a12a-46e8-9989-240a63175c99 结束（completed，1 轮，0 次工具调用，成本 0 CNY（未计价））
        wisp run: 回复已完成
        [audit] wisp run: SESSION-MINT id=sess_64a52a606f4f657f9a14a3fc1b09c058 (结束点＝本进程退出，A435 第 2 条)
        [audit] p
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏
time=2026-10-09T14:41:24.542+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict存储损坏1114007260\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识
time=2026-10-09T14:41:24.552+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict版本不认识3386123974\001\logs min_level=info early_records=0 early_dropped=0
=== RUN   TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到
time=2026-10-09T14:41:24.559+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket101UnreadableModeFailsLoudlyAndStrict权限读不到2707302125\001\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict (1.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识 (0.01s)
    --- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict/权限读不到 (0.01s)
=== RUN   TestRunTextTaskTextPathEndToEnd
time=2026-10-09T14:41:25.495+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:25 mockllm: serving on http://127.0.0.1:52742 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:25.513+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd3657609607\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:25.533+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd3657609607\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:25.539+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:25Z duration_ms=6
time=2026-10-09T14:41:25.547+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskTextPathEndToEnd3657609607\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:25.552+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:25Z duration_ms=5
time=2026-10-09T14:41:25.557+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:25.559+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:25.583+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskTextPathEndToEnd (1.02s)
=== RUN   TestRunTextTaskFailNextIsClassified
time=2026-10-09T14:41:26.505+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:26 mockllm: serving on http://127.0.0.1:52749 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:26.523+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified3963762425\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:26.542+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified3963762425\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:26.550+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:26Z duration_ms=7
time=2026-10-09T14:41:26.556+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskFailNextIsClassified3963762425\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:26.561+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:26Z duration_ms=4
time=2026-10-09T14:41:26.566+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:26.568+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:26.596+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskFailNextIsClassified (1.01s)
=== RUN   TestHostDispatchThroughTheAssembledBridge
time=2026-10-09T14:41:27.498+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:27 mockllm: serving on http://127.0.0.1:53692 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:27.515+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge723743137\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:27.534+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge723743137\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:27.540+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:27Z duration_ms=6
time=2026-10-09T14:41:27.546+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestHostDispatchThroughTheAssembledBridge723743137\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:27.552+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:27Z duration_ms=5
time=2026-10-09T14:41:27.556+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:27.557+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:27.635+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestHostDispatchThroughTheAssembledBridge (1.04s)
=== RUN   TestComposedGateBlocksAWriteForTwoSeconds
time=2026-10-09T14:41:28.560+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:28 mockllm: serving on http://127.0.0.1:53695 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:28.577+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2406989541\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:28.599+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2406989541\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:28.606+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:28Z duration_ms=6
time=2026-10-09T14:41:28.613+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestComposedGateBlocksAWriteForTwoSeconds2406989541\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:28.617+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:28Z duration_ms=4
time=2026-10-09T14:41:28.621+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:28.623+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:30.676+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestComposedGateBlocksAWriteForTwoSeconds (3.04s)
=== RUN   TestRunTextTaskKeyResolvesInTheStore
time=2026-10-09T14:41:31.601+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:31 mockllm: serving on http://127.0.0.1:49446 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:31.619+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore806217637\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:31.639+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore806217637\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:31.646+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:31Z duration_ms=6
time=2026-10-09T14:41:31.652+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunTextTaskKeyResolvesInTheStore806217637\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:31.658+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:31Z duration_ms=5
time=2026-10-09T14:41:31.662+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:31.664+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestRunTextTaskKeyResolvesInTheStore (1.00s)
=== RUN   TestMissingBlobFailsUnconfiguredNeverSilently
time=2026-10-09T14:41:32.616+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:32 mockllm: serving on http://127.0.0.1:49450 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:32.637+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestMissingBlobFailsUnconfiguredNeverSilently738187004\002\logs min_level=info early_records=0 early_dropped=0
--- PASS: TestMissingBlobFailsUnconfiguredNeverSilently (0.96s)
=== RUN   TestSecretArgvCarriesNoSecret
=== RUN   TestSecretArgvCarriesNoSecret/from-stdin
=== RUN   TestSecretArgvCarriesNoSecret/interactive-without-console
--- PASS: TestSecretArgvCarriesNoSecret (2.77s)
    --- PASS: TestSecretArgvCarriesNoSecret/from-stdin (0.10s)
    --- PASS: TestSecretArgvCarriesNoSecret/interactive-without-console (0.04s)
=== RUN   TestSecretRealBinaryRefusesValueFlag
--- PASS: TestSecretRealBinaryRefusesValueFlag (2.80s)
=== RUN   TestProcessCommandLineProbeHelperProcess
--- PASS: TestProcessCommandLineProbeHelperProcess (0.00s)
=== RUN   TestProcessCommandLineProbeDetectsAPlantedValue
--- PASS: TestProcessCommandLineProbeDetectsAPlantedValue (0.06s)
=== RUN   TestSecretSetGetListRoundTrip
--- PASS: TestSecretSetGetListRoundTrip (0.01s)
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/traversal
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/dotdot
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/separator
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/space
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/empty
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/empty_stdin
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/multiline_stdin
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/no_name
=== RUN   TestSecretSetRejectsBadNamesAndEmptyInput/two_names
--- PASS: TestSecretSetRejectsBadNamesAndEmptyInput (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/traversal (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/dotdot (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/separator (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/space (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/empty (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/empty_stdin (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/multiline_stdin (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/no_name (0.00s)
    --- PASS: TestSecretSetRejectsBadNamesAndEmptyInput/two_names (0.00s)
=== RUN   TestSecretSetConfirmationMismatchStoresNothing
--- PASS: TestSecretSetConfirmationMismatchStoresNothing (0.00s)
=== RUN   TestSecretSetWithoutConsolePointsAtFromStdin
--- PASS: TestSecretSetWithoutConsolePointsAtFromStdin (0.00s)
=== RUN   TestSecretFlagsAreBoolOnly
--- PASS: TestSecretFlagsAreBoolOnly (0.00s)
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--value
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--secret
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--key
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/-k
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--api-key
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--from-file
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--file
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--stdin
=== RUN   TestSecretValueCarryingFlagsAreRefusedAndUnechoed/positional
--- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--value (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--secret (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--key (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/-k (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--api-key (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--from-file (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--file (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/--stdin (0.00s)
    --- PASS: TestSecretValueCarryingFlagsAreRefusedAndUnechoed/positional (0.00s)
=== RUN   TestSecretFromStdinWritesNoIntermediateFile
--- PASS: TestSecretFromStdinWritesNoIntermediateFile (0.28s)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/unwritable_store_dir
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)
=== RUN   TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only
--- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext (0.03s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable) (0.01s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13) (0.01s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/unwritable_store_dir (0.00s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob) (0.00s)
    --- PASS: TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only (0.01s)
=== RUN   TestSecretUnsetRefusesWhileReferenced
=== RUN   TestSecretUnsetRefusesWhileReferenced/refused_while_referenced
=== RUN   TestSecretUnsetRefusesWhileReferenced/refused_for_a_voice.realtime_reference_too
=== RUN   TestSecretUnsetRefusesWhileReferenced/an_unreadable_config_fails_closed
=== RUN   TestSecretUnsetRefusesWhileReferenced/force_deletes_and_writes_an_audit_line
=== RUN   TestSecretUnsetRefusesWhileReferenced/unreferenced_delete_is_audited_at_info
=== RUN   TestSecretUnsetRefusesWhileReferenced/no_config_at_all_means_no_references
--- PASS: TestSecretUnsetRefusesWhileReferenced (0.02s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/refused_while_referenced (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/refused_for_a_voice.realtime_reference_too (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/an_unreadable_config_fails_closed (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/force_deletes_and_writes_an_audit_line (0.00s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/unreferenced_delete_is_audited_at_info (0.01s)
    --- PASS: TestSecretUnsetRefusesWhileReferenced/no_config_at_all_means_no_references (0.01s)
=== RUN   TestSecretSameNameUnderThreeEnvsIsThreeBlobs
--- PASS: TestSecretSameNameUnderThreeEnvsIsThreeBlobs (0.02s)
=== RUN   TestSecretPortableModeUsesTicket06Seam
--- PASS: TestSecretPortableModeUsesTicket06Seam (0.01s)
=== RUN   TestSecretEndToEndConfigRefResolvesAtRequestTime
--- PASS: TestSecretEndToEndConfigRefResolvesAtRequestTime (0.01s)
=== RUN   TestSecretUsageAndUnknownSubcommand
--- PASS: TestSecretUsageAndUnknownSubcommand (0.00s)
=== RUN   TestSecretOverwriteIsAnnounced
--- PASS: TestSecretOverwriteIsAnnounced (0.01s)
=== RUN   TestSLO156FixtureChildrenDoWhatTheirNamesSay
--- PASS: TestSLO156FixtureChildrenDoWhatTheirNamesSay (0.05s)
=== RUN   TestSLO156ExitedAsksTheOSForAChildNobodyReaped
--- PASS: TestSLO156ExitedAsksTheOSForAChildNobodyReaped (0.03s)
=== RUN   TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees
    slo_exit_os_156_windows_test.go:206: Wait on a killed child reports the kill as an error (exit status 1); that is the fixture, not a failure
--- PASS: TestSLO156LiveSubjectStaysAliveUntilTheOSDisagrees (0.00s)
=== RUN   TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn
--- PASS: TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn (0.03s)
=== RUN   TestSLO156WaitReadyNamesTheDeadSubjectToo
--- PASS: TestSLO156WaitReadyNamesTheDeadSubjectToo (0.03s)
=== RUN   TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt
--- PASS: TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt (0.01s)
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/html-head
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/stray-comma
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/wrong-type
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/trailing-garbage
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/second-document
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/complete-but-no-state-report
=== RUN   TestSLO144ReportsThatContradictThemselvesAreCorruptNow/empty-object-is-a-finished-lie
--- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/html-head (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/stray-comma (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/wrong-type (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/trailing-garbage (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/second-document (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/complete-but-no-state-report (0.00s)
    --- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow/empty-object-is-a-finished-lie (0.00s)
=== RUN   TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport
--- PASS: TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport (0.03s)
=== RUN   TestSLO144LoopGiveUpSentencesOnRealFiles
--- PASS: TestSLO144LoopGiveUpSentencesOnRealFiles (0.05s)
=== RUN   TestSLO144CorruptReportIsJudgedOnTheFirstRead
--- PASS: TestSLO144CorruptReportIsJudgedOnTheFirstRead (0.00s)
=== RUN   TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes
--- PASS: TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes (0.06s)
=== RUN   TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected
--- PASS: TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected (0.00s)
=== RUN   TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt
--- PASS: TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt (0.02s)
=== RUN   TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead
--- PASS: TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead (0.06s)
=== RUN   TestSLO147OffsetSemanticsRenderThreeDifferentSentences
--- PASS: TestSLO147OffsetSemanticsRenderThreeDifferentSentences (0.04s)
=== RUN   TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt
--- PASS: TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt (0.00s)
=== RUN   TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd
--- PASS: TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd (0.00s)
=== RUN   TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne
--- PASS: TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne (0.00s)
=== RUN   TestSLO149ExitedGiveUpSentenceCarriesTheLastReading
--- PASS: TestSLO149ExitedGiveUpSentenceCarriesTheLastReading (0.03s)
=== RUN   TestRunPacketMarksTheRosterRowACardIsHolding
time=2026-10-09T14:41:40.067+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:40 mockllm: serving on http://127.0.0.1:60918 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:40.085+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding208171164\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:40.105+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding208171164\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:40.111+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:40Z duration_ms=5
time=2026-10-09T14:41:40.118+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketMarksTheRosterRowACardIsHolding208171164\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:40.124+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:40Z duration_ms=5
time=2026-10-09T14:41:40.128+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:40.130+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:40.137+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ecf2b05a-b7d0-4d48-ad58-2d5489b4b81e owner=tools
    subagent_blocked_197_test.go:168: tasks wire bytes: {"rows":[{"taskId":"3156441c-5c9c-4ab1-bf69-2d31f56cfde5","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"3156441c-5c9c-4ab1-bf69-2d31f56cfde5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ecf2b05a-b7d0-4d48-ad58-2d5489b4b81e","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"3156441c-5c9c-4ab1-bf69-2d31f56cfde5","status":"Thinking","statusKnown":true,"streamKey":"subagent:ecf2b05a-b7d0-4d48-ad58-2d5489b4b81e","blockedOnApproval":true,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:222: tasks wire bytes: {"rows":[{"taskId":"3156441c-5c9c-4ab1-bf69-2d31f56cfde5","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"3156441c-5c9c-4ab1-bf69-2d31f56cfde5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ecf2b05a-b7d0-4d48-ad58-2d5489b4b81e","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"3156441c-5c9c-4ab1-bf69-2d31f56cfde5","status":"Settling","statusKnown":true,"streamKey":"subagent:ecf2b05a-b7d0-4d48-ad58-2d5489b4b81e","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_blocked_197_test.go:234: blocked row on the run's own packet: task=ecf2b05a-b7d0-4d48-ad58-2d5489b4b81e status=Thinking streamKey=subagent:ecf2b05a-b7d0-4d48-ad58-2d5489b4b81e card=ticket197.blocked.probe pending=1 bytes=3147 sha256=146c9dd0a6fba106 | after the card: blocked=false answer=reject why="任务上下文已结束，审批请求已作废并按拒绝处理"
--- PASS: TestRunPacketMarksTheRosterRowACardIsHolding (1.33s)
=== RUN   TestRunPacketCarriesTheSubagentItsRosterRowFed
time=2026-10-09T14:41:41.403+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:41 mockllm: serving on http://127.0.0.1:59565 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:41.421+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed1516423123\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:41.444+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed1516423123\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:41.451+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:41Z duration_ms=6
time=2026-10-09T14:41:41.459+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketCarriesTheSubagentItsRosterRowFed1516423123\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:41.467+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:41Z duration_ms=7
time=2026-10-09T14:41:41.471+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:41.473+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:41.479+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c0b3f177-f459-4156-8b7f-1091d9892e10 owner=tools
    subagent_carrier_197_test.go:278: tasks wire bytes: {"rows":[{"taskId":"c0b3f177-f459-4156-8b7f-1091d9892e10","label":"载体层正控：把一句话原样说出来","kind":"subagent","parentTaskId":"c5ae6338-b1ee-43b1-8e6b-1ff8294dead4","status":"Thinking","statusKnown":true,"streamKey":"subagent:c0b3f177-f459-4156-8b7f-1091d9892e10","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c5ae6338-b1ee-43b1-8e6b-1ff8294dead4","label":"总结一下 rootprompt197r3 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"c5ae6338-b1ee-43b1-8e6b-1ff8294dead4","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":1,"poolCap":4,"streamTruncated":false,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:388: packet tasks section: rows=2 poolCap=4 child=c0b3f177-f459-4156-8b7f-1091d9892e10 runStatus=Thinking afterStatus=Settling key=subagent:c0b3f177-f459-4156-8b7f-1091d9892e10 bytes=2845
--- PASS: TestRunPacketCarriesTheSubagentItsRosterRowFed (1.35s)
=== RUN   TestSubagentStreamKeyHasOneMintSite
--- PASS: TestSubagentStreamKeyHasOneMintSite (0.11s)
=== RUN   TestRunPacketReportsTheStreamLogPastItsBound
time=2026-10-09T14:41:42.881+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:42 mockllm: serving on http://127.0.0.1:61662 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:42.897+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1949263771\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:42.919+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1949263771\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:42.924+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:42Z duration_ms=5
time=2026-10-09T14:41:42.930+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestRunPacketReportsTheStreamLogPastItsBound1949263771\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:42.936+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:42Z duration_ms=5
time=2026-10-09T14:41:42.944+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:42.946+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:42.952+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-7bd56cc0-68f2-41f0-83f2-cac94e61cac1 owner=tools
time=2026-10-09T14:41:42.957+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-67b954de-8474-487c-a24d-0759b6c5afa2 owner=tools
time=2026-10-09T14:41:42.964+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-bea573f6-1f9c-47f4-a572-03fa537b6c17 owner=tools
time=2026-10-09T14:41:42.968+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-0a1f1a1a-1404-4a41-adcc-f1047923f58f owner=tools
time=2026-10-09T14:41:42.972+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-db7aef4c-db77-497f-8009-e30a6c7e7933 owner=tools
time=2026-10-09T14:41:42.977+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-8ffb8106-cbd7-4304-9520-651b5eda197e owner=tools
time=2026-10-09T14:41:42.981+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c593879a-7a07-4e0a-b92d-c37be49a63fe owner=tools
time=2026-10-09T14:41:42.986+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-889b02d5-13fe-4e49-b854-27c0f6ec1e83 owner=tools
time=2026-10-09T14:41:42.990+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-911fff00-264a-46cd-967d-a12ace998b1f owner=tools
time=2026-10-09T14:41:42.996+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-d62b1fb2-9da8-46ba-bfee-7d3dd73bea5a owner=tools
time=2026-10-09T14:41:43.000+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-193788cf-72d5-4848-84ad-171ea7c78ce2 owner=tools
time=2026-10-09T14:41:43.004+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-ee40bcfa-d2d3-4cd2-9a2f-7a85d1daa227 owner=tools
time=2026-10-09T14:41:43.010+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-db9ef564-d5e1-4496-b1c7-5b1edf46d240 owner=tools
time=2026-10-09T14:41:43.015+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-3721e513-628b-41f8-b23c-2f6beb0dd9ec owner=tools
time=2026-10-09T14:41:43.019+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-1cc23fd9-3fb5-4cff-a46f-ebe84ba68546 owner=tools
time=2026-10-09T14:41:43.023+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f8f3df5b-7cd3-4206-aacd-f19bbf2e2337 owner=tools
time=2026-10-09T14:41:43.027+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-2116e3d0-d7d3-40cc-9cc7-17d3f01b83b9 owner=tools
time=2026-10-09T14:41:43.032+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-adfe642b-2786-4e6e-aba4-d85cf8138120 owner=tools
time=2026-10-09T14:41:43.036+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-3679d1bc-d3da-4593-98db-c7057a21c7ff owner=tools
time=2026-10-09T14:41:43.041+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-8b9de839-2f15-4e0b-bb9a-c6af14b53f36 owner=tools
time=2026-10-09T14:41:43.045+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-064db3ad-904e-4a4a-a505-f341ccecef37 owner=tools
time=2026-10-09T14:41:43.050+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-26c51934-0d6f-49b3-8772-1ec8b1d6fd1c owner=tools
time=2026-10-09T14:41:43.055+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-c8fbc074-66e3-4412-a182-141a5d09ea12 owner=tools
time=2026-10-09T14:41:43.059+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-32f54c68-3c60-440e-aa1d-aa264f402d07 owner=tools
time=2026-10-09T14:41:43.063+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-873170bc-7e43-46ff-8219-7f65b24cfe0f owner=tools
time=2026-10-09T14:41:43.068+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-5f4388e2-ba17-486a-917e-905165c017d5 owner=tools
time=2026-10-09T14:41:43.072+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-9335cbad-f370-4d4f-9353-3ad88a923537 owner=tools
time=2026-10-09T14:41:43.077+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-444a5b72-b380-4870-8422-81c39f5170fb owner=tools
time=2026-10-09T14:41:43.082+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f765c021-ec37-4c78-be5b-ce4e72b82e59 owner=tools
time=2026-10-09T14:41:43.086+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-879579bb-30e0-4b52-ab7d-b318593ca594 owner=tools
time=2026-10-09T14:41:43.091+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-42f467d4-c136-442e-ba93-8c4ec33ca940 owner=tools
time=2026-10-09T14:41:43.095+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-96ad3f8c-b92e-4f21-85e4-0a711294c57a owner=tools
time=2026-10-09T14:41:43.101+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-8bf1dcb2-36a7-4669-a51e-0eda62e6e057 owner=tools
    subagent_carrier_197_test.go:548: tasks wire bytes: {"rows":[{"taskId":"064db3ad-904e-4a4a-a505-f341ccecef37","label":"溢出正控 20 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:064db3ad-904e-4a4a-a505-f341ccecef37","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"0a1f1a1a-1404-4a41-adcc-f1047923f58f","label":"溢出正控 03 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:0a1f1a1a-1404-4a41-adcc-f1047923f58f","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"193788cf-72d5-4848-84ad-171ea7c78ce2","label":"溢出正控 10 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:193788cf-72d5-4848-84ad-171ea7c78ce2","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"1cc23fd9-3fb5-4cff-a46f-ebe84ba68546","label":"溢出正控 14 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:1cc23fd9-3fb5-4cff-a46f-ebe84ba68546","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"2116e3d0-d7d3-40cc-9cc7-17d3f01b83b9","label":"溢出正控 16 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:2116e3d0-d7d3-40cc-9cc7-17d3f01b83b9","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"26c51934-0d6f-49b3-8772-1ec8b1d6fd1c","label":"溢出正控 21 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:26c51934-0d6f-49b3-8772-1ec8b1d6fd1c","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"32f54c68-3c60-440e-aa1d-aa264f402d07","label":"溢出正控 23 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:32f54c68-3c60-440e-aa1d-aa264f402d07","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"3679d1bc-d3da-4593-98db-c7057a21c7ff","label":"溢出正控 18 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:3679d1bc-d3da-4593-98db-c7057a21c7ff","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"3721e513-628b-41f8-b23c-2f6beb0dd9ec","label":"溢出正控 13 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:3721e513-628b-41f8-b23c-2f6beb0dd9ec","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","label":"总结一下 这份笔记","kind":"root","parentTaskId":"","status":"","statusKnown":false,"statusReason":"宿主没有登记这一维（fail-closed：不替任务编一个状态）","streamKey":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"42f467d4-c136-442e-ba93-8c4ec33ca940","label":"溢出正控 30 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:42f467d4-c136-442e-ba93-8c4ec33ca940","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"444a5b72-b380-4870-8422-81c39f5170fb","label":"溢出正控 27 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:444a5b72-b380-4870-8422-81c39f5170fb","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"5f4388e2-ba17-486a-917e-905165c017d5","label":"溢出正控 25 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:5f4388e2-ba17-486a-917e-905165c017d5","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"67b954de-8474-487c-a24d-0759b6c5afa2","label":"溢出正控 01 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:67b954de-8474-487c-a24d-0759b6c5afa2","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"7bd56cc0-68f2-41f0-83f2-cac94e61cac1","label":"溢出正控 00 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:7bd56cc0-68f2-41f0-83f2-cac94e61cac1","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"873170bc-7e43-46ff-8219-7f65b24cfe0f","label":"溢出正控 24 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:873170bc-7e43-46ff-8219-7f65b24cfe0f","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"879579bb-30e0-4b52-ab7d-b318593ca594","label":"溢出正控 29 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:879579bb-30e0-4b52-ab7d-b318593ca594","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"889b02d5-13fe-4e49-b854-27c0f6ec1e83","label":"溢出正控 07 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:889b02d5-13fe-4e49-b854-27c0f6ec1e83","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"8b9de839-2f15-4e0b-bb9a-c6af14b53f36","label":"溢出正控 19 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:8b9de839-2f15-4e0b-bb9a-c6af14b53f36","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"8bf1dcb2-36a7-4669-a51e-0eda62e6e057","label":"溢出正控 32 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:8bf1dcb2-36a7-4669-a51e-0eda62e6e057","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"8ffb8106-cbd7-4304-9520-651b5eda197e","label":"溢出正控 05 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:8ffb8106-cbd7-4304-9520-651b5eda197e","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"911fff00-264a-46cd-967d-a12ace998b1f","label":"溢出正控 08 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:911fff00-264a-46cd-967d-a12ace998b1f","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"9335cbad-f370-4d4f-9353-3ad88a923537","label":"溢出正控 26 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:9335cbad-f370-4d4f-9353-3ad88a923537","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"96ad3f8c-b92e-4f21-85e4-0a711294c57a","label":"溢出正控 31 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:96ad3f8c-b92e-4f21-85e4-0a711294c57a","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"adfe642b-2786-4e6e-aba4-d85cf8138120","label":"溢出正控 17 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:adfe642b-2786-4e6e-aba4-d85cf8138120","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"bea573f6-1f9c-47f4-a572-03fa537b6c17","label":"溢出正控 02 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:bea573f6-1f9c-47f4-a572-03fa537b6c17","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c593879a-7a07-4e0a-b92d-c37be49a63fe","label":"溢出正控 06 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:c593879a-7a07-4e0a-b92d-c37be49a63fe","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"c8fbc074-66e3-4412-a182-141a5d09ea12","label":"溢出正控 22 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:c8fbc074-66e3-4412-a182-141a5d09ea12","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"d62b1fb2-9da8-46ba-bfee-7d3dd73bea5a","label":"溢出正控 09 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:d62b1fb2-9da8-46ba-bfee-7d3dd73bea5a","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"db7aef4c-db77-497f-8009-e30a6c7e7933","label":"溢出正控 04 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:db7aef4c-db77-497f-8009-e30a6c7e7933","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"db9ef564-d5e1-4496-b1c7-5b1edf46d240","label":"溢出正控 12 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:db9ef564-d5e1-4496-b1c7-5b1edf46d240","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"ee40bcfa-d2d3-4cd2-9a2f-7a85d1daa227","label":"溢出正控 11 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:ee40bcfa-d2d3-4cd2-9a2f-7a85d1daa227","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f765c021-ec37-4c78-be5b-ce4e72b82e59","label":"溢出正控 28 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:f765c021-ec37-4c78-be5b-ce4e72b82e59","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false},{"taskId":"f8f3df5b-7cd3-4206-aacd-f19bbf2e2337","label":"溢出正控 15 llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll","kind":"subagent","parentTaskId":"3e0298d6-d691-4a50-97f3-8ee3861ea8ee","status":"Settling","statusKnown":true,"streamKey":"subagent:f8f3df5b-7cd3-4206-aacd-f19bbf2e2337","blockedOnApproval":false,"streamTruncated":false,"streamElidedRunes":0,"streamDropped":false}],"inFlightSlots":0,"poolCap":4,"streamTruncated":true,"streamElidedRunes":0,"droppedStreamKeys":[]}
    subagent_carrier_197_test.go:591: bound crossed: rows=34 streams=34 truncated=true elided=0 dropped=[]
--- PASS: TestRunPacketReportsTheStreamLogPastItsBound (1.18s)
=== RUN   Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands
time=2026-10-09T14:41:44.085+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:44 mockllm: serving on http://127.0.0.1:61737 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:44.102+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands979333811\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:44.121+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands979333811\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:44.127+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:44Z duration_ms=5
time=2026-10-09T14:41:44.133+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands979333811\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:44.139+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:44Z duration_ms=4
time=2026-10-09T14:41:44.142+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:44.144+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:44.151+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=subagent-finish-f6fd4a7d-d4c2-40dd-b057-8a9d213eeffc owner=tools
    subagent_selfapproval_197_test.go:366: child=f6fd4a7d-d4c2-40dd-b057-8a9d213eeffc corr1=f6fd4a7d-d4c2-40dd-b057-8a9d213eeffc-corr-1 corr2=f6fd4a7d-d4c2-40dd-b057-8a9d213eeffc-corr-2 | 空令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 假令牌=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 自称来源=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 借证=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 重放=approval: correlation_id 无对应待审批项 | 面板递证=面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数 | 烧后再试=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） | 宿主允许=<nil> | 落盘=true/false
--- PASS: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.06s)
=== RUN   Test197NoAllowDoorIsReachableFromASubagentsAssembly
--- PASS: Test197NoAllowDoorIsReachableFromASubagentsAssembly (0.00s)
=== RUN   TestSubagentStreamKeyPrefixAgreesAcrossBothPackages
--- PASS: TestSubagentStreamKeyPrefixAgreesAcrossBothPackages (0.00s)
=== RUN   TestSubagentStreamKeyBuildersAgreeAcrossBothPackages
--- PASS: TestSubagentStreamKeyBuildersAgreeAcrossBothPackages (0.00s)
=== RUN   TestCompositionRootClosesTheLoopTasksTaintScope
time=2026-10-09T14:41:55.075+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:55 mockllm: serving on http://127.0.0.1:58690 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:55.092+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope2770213035\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:55.113+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope2770213035\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:55.119+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:55Z duration_ms=5
time=2026-10-09T14:41:55.125+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestCompositionRootClosesTheLoopTasksTaintScope2770213035\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:55.130+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:55Z duration_ms=4
time=2026-10-09T14:41:55.134+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:55.135+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestCompositionRootClosesTheLoopTasksTaintScope (0.98s)
=== RUN   TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit
time=2026-10-09T14:41:56.070+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:56 mockllm: serving on http://127.0.0.1:58693 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:56.088+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit3988081108\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:56.108+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit3988081108\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:56.114+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:56Z duration_ms=5
time=2026-10-09T14:41:56.120+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit3988081108\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:56.126+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:56Z duration_ms=5
time=2026-10-09T14:41:56.129+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:56.133+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.01s)
=== RUN   TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun
time=2026-10-09T14:41:57.076+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:57 mockllm: serving on http://127.0.0.1:58696 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:57.092+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun1244501135\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:57.111+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun1244501135\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:57.117+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:57Z duration_ms=5
time=2026-10-09T14:41:57.123+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun1244501135\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:57.128+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:57Z duration_ms=5
time=2026-10-09T14:41:57.132+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:57.134+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:57.174+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:41:57.183+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (1.03s)
=== RUN   TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking
time=2026-10-09T14:41:58.103+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:41:58 mockllm: serving on http://127.0.0.1:51389 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:41:58.121+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking1540209797\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:41:58.139+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking1540209797\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:41:58.145+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:41:58Z duration_ms=5
time=2026-10-09T14:41:58.152+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking1540209797\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:41:58.158+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:41:58Z duration_ms=6
time=2026-10-09T14:41:58.162+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:41:58.165+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:00.229+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (3.04s)
=== RUN   TestTicket224ProductionSessionDoesNotSurviveRestart
time=2026-10-09T14:42:01.142+08:00 level=WARN msg="goroutine outside the D38 roster (leak symptom)" goroutine=mockllm-stdout-reader owner=test
2026/10/09 14:42:01 mockllm: serving on http://127.0.0.1:50580 (golden dir D:\work\workspace\projects plans\Wisp\internal\llm\testdata\golden)
time=2026-10-09T14:42:01.158+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart4004592580\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:01.178+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart4004592580\002\backup\wisp.db.bak-0-1 from=0 to=1
time=2026-10-09T14:42:01.184+08:00 level=INFO msg="schema migration step applied" component=memory from=0 to=1 started_at=2026-10-09T06:42:01Z duration_ms=5
time=2026-10-09T14:42:01.190+08:00 level=INFO msg="pre-migration backup written" component=memory backup=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart4004592580\002\backup\wisp.db.bak-1-2 from=1 to=2
time=2026-10-09T14:42:01.198+08:00 level=INFO msg="schema migration step applied" component=memory from=1 to=2 started_at=2026-10-09T06:42:01Z duration_ms=7
time=2026-10-09T14:42:01.206+08:00 level=INFO msg="schema migration complete" component=memory schema_version=2
time=2026-10-09T14:42:01.207+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:01.226+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\Users\swq\AppData\Local\Temp\TestTicket224ProductionSessionDoesNotSurviveRestart4004592580\002\logs min_level=info early_records=0 early_dropped=0
time=2026-10-09T14:42:01.239+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
time=2026-10-09T14:42:03.297+08:00 level=INFO msg="startup WAL checkpoint (TRUNCATE) done" component=memory wal_pages_before=0 pages_moved=0
--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (3.07s)
FAIL
FAIL	github.com/CarlosShao/wisp/cmd/wisp	451.539s
FAIL
rc(go-test-full-post)=1
