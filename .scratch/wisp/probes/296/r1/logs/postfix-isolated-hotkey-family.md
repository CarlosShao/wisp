time=2026-10-09T21:38:42.451+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   Test258AssemblyRootWiresTheChainAndTheBridge
--- PASS: Test258AssemblyRootWiresTheChainAndTheBridge (0.00s)
=== RUN   Test258SinkRebindRulerReadsBothSpellings
--- PASS: Test258SinkRebindRulerReadsBothSpellings (0.00s)
=== RUN   Test258ConstructionChainTakesConfigValues
--- PASS: Test258ConstructionChainTakesConfigValues (0.01s)
=== RUN   Test258ConstructionChainMissingFileFallsBackAndSaysIt
--- PASS: Test258ConstructionChainMissingFileFallsBackAndSaysIt (0.01s)
=== RUN   Test258SectionMissingTierWordIsDefaults
--- PASS: Test258SectionMissingTierWordIsDefaults (0.01s)
=== RUN   Test258ProvenanceWordsAreThePrintedOnes
--- PASS: Test258ProvenanceWordsAreThePrintedOnes (0.00s)
=== RUN   Test258BridgeRebindsLiveKeysFromConfigEdit
time=2026-10-09T21:38:42.860+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:42.877+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T21:38:42.877+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T21:38:42.883+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:42.883+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+R mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T21:38:42.888+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeRebindsLiveKeysFromConfigEdit (0.40s)
=== RUN   Test258BridgeMutationNoSrcKeepsOldBinding
time=2026-10-09T21:38:42.920+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:42.925+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T21:38:43.337+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258BridgeMutationNoSrcKeepsOldBinding (0.45s)
=== RUN   Test258OccupiedCombinationNamesTheNewValue
time=2026-10-09T21:38:43.374+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:43.378+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T21:38:43.378+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T21:38:43.385+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:43.385+08:00 level=ERROR msg="hotkey occupied by another program, not registered" hotkey=panel binding=Ctrl+Alt+V err="Hot key is already registered."
time=2026-10-09T21:38:43.385+08:00 level=ERROR msg="ball: hotkey reload" binding="hotkey panel = \"Ctrl+Alt+V\" is occupied by another program and was NOT registered; pressing it will do nothing until you pick a free combination in [hotkey]"
time=2026-10-09T21:38:43.385+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+Z mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+V live=2
time=2026-10-09T21:38:43.393+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258OccupiedCombinationNamesTheNewValue (0.06s)
=== RUN   Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots
time=2026-10-09T21:38:43.429+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:43.434+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=defaults
wisp: ball hotkey reload bridge armed (provenance=defaults); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T21:38:43.435+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T21:38:43.445+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots (0.05s)
=== RUN   Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable
time=2026-10-09T21:38:43.481+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:43.486+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=defaults
wisp: ball hotkey reload bridge armed (provenance=defaults); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T21:38:43.486+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T21:38:43.496+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable (0.05s)
=== RUN   Test258V1ProbeSummonEditRebindsLiveBall
time=2026-10-09T21:38:43.528+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:43.533+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=config
wisp: ball hotkey reload bridge armed (provenance=config); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T21:38:43.533+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=config hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T21:38:43.538+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T21:38:43.539+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon=Ctrl+Alt+7 mute=Ctrl+Alt+M cancel=Esc panel=Ctrl+Alt+P live=3
time=2026-10-09T21:38:43.543+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (0.05s)
=== RUN   Test258V1ProbeConstructionMissingFileNamesTheFallback
--- PASS: Test258V1ProbeConstructionMissingFileNamesTheFallback (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/cmd/wisp	1.211s
