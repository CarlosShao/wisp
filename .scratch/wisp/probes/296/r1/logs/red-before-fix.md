time=2026-10-09T20:22:09.917+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
=== RUN   Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots
time=2026-10-09T20:22:10.199+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T20:22:10.215+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=defaults
wisp: ball hotkey reload bridge armed (provenance=defaults); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T20:22:10.215+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T20:22:10.217+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=summon
time=2026-10-09T20:22:10.217+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=mute
time=2026-10-09T20:22:10.217+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T20:22:10.217+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=panel
time=2026-10-09T20:22:10.217+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon="" mute="" cancel=Esc panel="" live=0
    resident_hotkey_296_windows_test.go:176: AC#1 form-1 (owner's three empty slots) RED: ConfiguredHotkeys() = {Summon: Mute: Cancel:Esc Panel:}, want the merged set {Summon:Ctrl+Alt+Q Mute:Ctrl+Alt+M Cancel:Esc Panel:Ctrl+Alt+P} - the bridge rebound onto the raw empty slots
    resident_hotkey_296_windows_test.go:176: AC#1 form-1 (owner's three empty slots) RED: after the bridge's first Check live=0 [], want live=3 (summon, mute, panel; missing [summon mute panel]), configured={Summon: Mute: Cancel:Esc Panel:}, per-slot=summon=disabled (unset in config)="" mute=disabled (unset in config)="" cancel=not bound while idle (cancel is borrowed only during Confirming)="Esc" panel=disabled (unset in config)="". The reload source handed the bridge un-merged empty slots and the keys were unregistered - ticket 296.
time=2026-10-09T20:22:10.221+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- FAIL: Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots (0.30s)
=== RUN   Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable
time=2026-10-09T20:22:10.246+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Esc
time=2026-10-09T20:22:10.251+08:00 level=INFO msg="ball: hotkey reload bridge armed (polls config, rebinds on change)" poll=1s provenance=defaults
wisp: ball hotkey reload bridge armed (provenance=defaults); a hand edit of [hotkey] rebinds the live keys without a restart
time=2026-10-09T20:22:10.252+08:00 level=INFO msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live" gestures="recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four veto channels stay reduced here to the one the assembly root injected"
time=2026-10-09T20:22:10.252+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=summon
time=2026-10-09T20:22:10.252+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=mute
time=2026-10-09T20:22:10.252+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=cancel
time=2026-10-09T20:22:10.252+08:00 level=INFO msg="hotkey disabled (unset in config)" hotkey=panel
time=2026-10-09T20:22:10.252+08:00 level=INFO msg="ball: hotkeys rebound after config change" summon="" mute="" cancel="" panel="" live=0
    resident_hotkey_296_windows_test.go:217: AC#1 form-2 (config.toml unreadable at tick time) RED: ConfiguredHotkeys() = {Summon: Mute: Cancel: Panel:}, want the merged set {Summon:Ctrl+Alt+Q Mute:Ctrl+Alt+M Cancel:Esc Panel:Ctrl+Alt+P} - the bridge rebound onto the raw empty slots
    resident_hotkey_296_windows_test.go:217: AC#1 form-2 (config.toml unreadable at tick time) RED: after the bridge's first Check live=0 [], want live=3 (summon, mute, panel; missing [summon mute panel]), configured={Summon: Mute: Cancel: Panel:}, per-slot=summon=disabled (unset in config)="" mute=disabled (unset in config)="" cancel=disabled (unset in config)="" panel=disabled (unset in config)="". The reload source handed the bridge un-merged empty slots and the keys were unregistered - ticket 296.
time=2026-10-09T20:22:10.259+08:00 level=INFO msg="ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
--- FAIL: Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable (0.04s)
FAIL
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.436s
FAIL
