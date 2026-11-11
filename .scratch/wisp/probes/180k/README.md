# 180-census follow-up · 落地腿 K（19 枚孤儿 config 字段）· 反向突变与门禁读数

时间：2026-11-11 12:3x–12:5x +08。起手锚：工作树（HEAD `6c0cbd244` 之上未提交字节）。
本文件与 `logs/*.txt` 只建不删。

## 反向（摘喂线 ⇒ 红，字节级还原 ⇒ 绿；每发命令、红句、还原读数在 logs/）

| 发 | 摘掉哪根喂线 | 落刀 | 哪个用例红 | 读数（logs/） |
|---|---|---|---|---|
| A1 | installLogSink 的 observe.level 喂线（`level = c.Observe.Level` → 硬编码常量） | logsink.go | Test180KObserveLevelGatesTheFile | mut-a-logsink-feed.txt：4 例全红（A1–A4 同发） |
| A2 | privacy.redact_paths 喂线（`RedactPaths: false` 字面量） | logsink.go | Test180KRedactPathsMasksTheFile | 同上 |
| A3 | roll.size_mb 喂线（`RollSizeMB: 10` 字面量） | logsink.go | Test180KRollSizeGatesTheFile | 同上 |
| A4 | roll.days 喂线（`RollDays: 7` 字面量） | logsink.go | Test180KRollDaysSweepsTheFile | mut-a4-rolldays-feed.txt：`roll.days=2 kept wisp-20261106-001.jsonl`；还原后 ok |
| B1 | memory.l1_max 喂线（run.go 摘 `memory.WithProfileCap(cfg.Memory.L1Max)`） | run.go | Test180KMemoryKnobsReachTheAssemblyRoot | mut-b-assembly-feeds.txt：`store.ProfileCap() = 20, want the file's 7` |
| B2 | L3 窗口喂线（run.go `days := 30` 字面量） | run.go | 同上 | 同上：audit 行读 `窗口 30 天` ≠ 文件的 7 |
| B3 | memory.l1_enabled 闸（attachMemoryL1IfEnabled180K 无条件放行） | memory_retention_180k.go | Test180KL1EnabledGateGatesTheAttach | 同上：`l1_enabled=false still reached the attach path` |
| B4 | audio.input_device 喂线（assembleCapture 摘 `audio.WithInputDevice(...)`） | resident_audio_windows.go | Test180KInputDeviceReachesTheCaptureConfig | 同上：`DeviceName "", want "Blue Yeti"` |
| C1 | conversation_mode 喂线（speakDoneFacts180K 硬编码 false） | resident_voice_engines_windows.go | Test180KConversationModeReachesTheSpeakDoneRow | mut-c-voice-feeds.txt：`speakDoneFacts lost the plan value` |
| C2 | aec.enabled 喂线（摘 `rv.bargein.setAECEnabled(plan.AECEnabled)`，两处构造点） | resident_voice_windows.go | Test180KAECSwitchGatesTheMonitor | 同上：arm 拒因变成 `没有 C11 作用域`，不含 `aec.enabled=false` |
| D1 | validate.go 的 audio.sample_rate 守卫 | internal/config/validate.go | Test180KSampleRateOffChainValueIsALoadError | mut-d-validate-guards.txt：`sample_rate = 48000 must be a load error` |
| D2 | validate.go 的 voice.aec.echo_ref 守卫 | internal/config/validate.go | Test180KAECEchoRefVocabulary | 同上：`echo_ref = loopback must be a load error` |

每一发还原后同用例复跑 = ok（正控/终控 GREEN）。⚠ 第一次 D 发（`if false && ref := …`）是非法
Go，编译失败不作为读数——已用语法合法的 `&& false` 重打（mut-d-…txt 为重打后读数）。

## 门禁读数（共享工作树，终态）

- gofumpt（cmd/wisp + internal/config + internal/audio + internal/memory）：空列表 rc0
- go vet ./cmd/wisp ./internal/config ./internal/audio ./internal/memory：rc0
- go build 同四包 + internal/observe：rc0
- go test ./internal/config 全包 2.60s ok / ./internal/audio 18.44s ok / ./internal/memory 30.53s ok（-count=1 -timeout 30m）
- go test ./cmd/wisp -run 'Test180K|Test423|TestTicket255|TestTicket324'：ok 19.86s
- TestRunPacketLightsTheRosterRowOfItsLiveL1Window -count=2：ok（整包首跑的那枚红 = 与前台突变跑并发共享树的干扰，恢复树后两连绿）
- 255 roster（TestTicket255*/TestTicket324* 全族）：ok——audio/privacy/memory/observe 四行按新读口改判（privacy/memory/observe：no-reader → snapshot-only；audio 保持 other-process、补 input_device 半边），sectionReadSites 三枚 nil 位补文件，另按 AC#6 规则把 9 行过期的 :LINE 读数对到 token 现位（token 未变）
- d22scan：clean rc0（bans #1-5 internal/=401 cmd/=68；ban #8 internal/=1002 cmd/=261）
- path-length：VERDICT GREEN rc0（longest=185，roster 62）
- 整包 ./cmd/wisp：终态树复跑见 logs（首跑 411s 的 5 枚红中 3 枚＝外观波 commit 具名的 HEAD 存量红：TestTicket198FirstRunCallerIsTheRunEntryOnly / TestLockedSuffixNeverTakesTheLockItself33r11（stampPublishLocked 未入册）/ Test313WritableUnionEqualsTheTwoRosters；roster 两枚与 RunPacket 一枚＝并发突变期共享树读源干扰，终态树复跑绿）

## 接线一览（19 枚逐枚裁决见 commit message 与完成报告）

接了 11 枚：observe.level、roll.size_mb、roll.days、privacy.redact_paths、memory.l1_max、
memory.l3_retention_days、privacy.retention_days、memory.l1_enabled、voice.conversation_mode、
voice.aec.enabled（+ voice.aec.echo_ref 守卫形）、audio.input_device、audio.sample_rate（守卫形）。
DEFERRED 6 枚（写面外/无生产对象）：app.single_instance（internal/proc.Boot）、
voice.wake_word.veto_words（internal/speech kwsVetoWords）、voice.tts.provider（云 TTS speak 半边未实现）、
voice.tts.speed（同缺口 + 本地速度旋钮在 internal/speech）、privacy.diagnostics_opt_in
（BuildDiagnosticsBundle 无生产调用点）、observe.slo_sample_interval_sec（watchdog 1s tick 是
SPEC-03 §4.3 冻结的双职责共享节拍）。
