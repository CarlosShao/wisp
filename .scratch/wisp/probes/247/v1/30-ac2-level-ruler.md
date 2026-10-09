# 247 v1 — question 1: is the AC#2 reading the right shape, and is the 0.37 floor the ruler's fault?

This leg ran the live harness itself (real microphone, real pinned thread) rather than reasoning
about the implementer's log. Prerequisites, re-checked in the same command as the run:
`tasklist //FI "IMAGENAME eq balldebug.exe"` and `… eq wisp.exe` => both
`INFO: No tasks are running which match the specified criteria.`

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_LIVE_MIC=1 \
  go test ./cmd/wisp/ -run TestAC247LiveMicrophoneLevelsReachTheBallSeam -count=1 -v -timeout 180s
rc_live=1
    resident_audio_247_live_windows_test.go:156: AC#2 出处: capture="麦克风 (Realtek(R) Audio)" render="扬声器 (Realtek(R) Audio)" pinned_thread=23820
    resident_audio_247_live_windows_test.go:181: AC#2 counters at the seam: levels=375 frames_sent=6 frames_dropped=369 last_err=""
    resident_audio_247_live_windows_test.go:184: AC#2 reading: quiet-A: samples=93 mean=0.368462 max=0.402494 levels_per_s=31.00
    resident_audio_247_live_windows_test.go:184: AC#2 reading: sound: samples=189 mean=0.381574 max=0.407773 levels_per_s=31.50
    resident_audio_247_live_windows_test.go:184: AC#2 reading: quiet-B: samples=93 mean=0.382262 max=0.404195 levels_per_s=31.00
resident_audio_247_live_windows_test.go:195: the loud phase is not distinguishable from silence: … (speakers off, or routed to an endpoint this microphone cannot hear)
--- FAIL: TestAC247LiveMicrophoneLevelsReachTheBallSeam (12.45s)
    … level=INFO msg="audio capture leg stopped: levels_delivered=376 frames_sent=6 frames_dropped=370 reopens=0"
```

Two independent runs (this leg's at 13:46, `247-r1`'s at 12:4x) agree to the third decimal on the
floor and to the *pattern*: quiet-A ≈ 0.3685/0.3691, sound ≈ 0.3816/0.3827, quiet-B ≈ 0.3823/0.3826.
So the reading is not a one-off and not something this leg can attribute to a bad capture session.

---

## (a) Can these two shapes satisfy AC#2? Is this the shape AC#2 asks for?

**No, and no.**

1. Not the shape. `247…md:25` verbatim: "一发真机读数＝**对麦克风说话**与不说话时 `SetAudioLevel`
   收到的值**不同**". What the harness puts in the room is a WAV out of the render endpoint
   (`resident_audio_247_live_windows_test.go:230-239`, `playAlarmThroughSpeakers`, a
   `System.Media.SoundPlayer` loop over `%WINDIR%\Media\Alarm01.wav`). Loudspeaker-to-microphone
   acoustic loopback is a *different* shape from voice-into-microphone: it depends on room
   geometry, speaker volume and AGC, and the process cannot see whether the speakers were off,
   muted per-app, or routed to another endpoint — the failure text says exactly that
   (`(speakers off, or routed to an endpoint this microphone cannot hear)`).
2. Not a satisfied reading even taken as-is. This leg's own numbers quantify why:

| 比较 | 由 RMS 反推的"新增交流能量" |
|---|---|
| sound − quiet-A | `sqrt(0.381574² − 0.368462²)` = **0.09917** FS (−20.1 dBFS) |
| quiet-B − quiet-A（纯漂移，中间什么都没放） | `sqrt(0.382262² − 0.368462²)` = **0.10178** FS |

   The apparent "signal" of the loud window is **smaller than the drift between the two silent
   windows**. So no, these readings cannot satisfy "说话与不说话时收到的值不同": there is no
   discriminator in them, and the two quiet windows alone would already have produced the same
   size of difference.
3. The bar the harness sets is unreachable while the floor sits where it sits — arithmetic:
   `:194` requires `loud.mean > 2 × quietA.mean = 0.7369`, i.e. an added AC component of
   `sqrt(0.7369² − 0.3685²)` = **0.6382 FS = −3.90 dBFS at the microphone**. That is clipping
   territory, not speech. So the next leg must not "close AC#2 by relaxing `:194`"; it must first
   bring the floor down (device/gain side) or re-shape the discriminator on the consumer side
   (which `level.go:41-48` explicitly assigns to the consumer leg), and only then speak into the mic.

## (b) Is the ~0.37 floor a defect of the level ruler, or of the fixture/device gain?

**Verdict: not an arithmetic defect of the ruler; the ruler is behaving exactly as its documented
scale says, and a DC offset of ≈12,074 int16 units would reproduce the floor almost exactly — but
this repository has no instrument today that can separate "DC offset" from "hot analog gain", and
that missing measurement is what the next leg owes, not a code fix.**

What the ruler computes (`internal/audio/level.go:96-107`, read verbatim by this leg):

```go
	var sumSquares int64
	for _, s := range samples {
		v := int64(s)
		sumSquares += v * v
	}
	rms := math.Sqrt(float64(sumSquares) / float64(len(samples)))
	return rms / LevelFullScale
```

- No DC removal, and it is **written down as intent** — `level.go:37-39` verbatim:
  `//	  DC counts: RMS does not remove a DC offset, so a frame held at -32768` /
  `//	  reads 1.0. Consumers that need "voice only" gate it themselves (the ball` /
  `//	  does, with SilenceLevelGate).`
  And `level.go:19-20` verbatim "It is a physical number, not a perceptual one: no gain, no
  compression, no noise floor, no smoothing."
- The ruler is not "always high": `TestLevelSilentFrameIsExactZero` **PASS** in this leg's run, and
  `LevelOfSamples` of an empty window is `MinLevel` (`level.go:97-99`).
- **The minimal reading this leg built itself** (no new repo file, no mutation): the repo's own DC
  fixture. `internal/audio/capturelevel_windows_test.go:21-24`, verbatim:
  ```go
  	samples := make([]int16, frames*FrameSamples)
  	for i := range samples {
  		samples[i] = 12345 // a steady DC level: RMS is exactly this magnitude
  ```
  Its level is `12345/32768 = 0.376740`, and the test that asserts it (PASS in this leg's run) is
  the proof that a **DC-only frame with no sound in it reads 0.377 on this scale**. The machine's
  measured floor 0.3685 corresponds to a DC of `0.368462 × 32768 = 12,074` LSB — within 2% of that
  fixture. So "a DC offset of ~12k LSB fully explains the floor" is a demonstrated possibility, not
  a guess.
- Counter-evidence that it is *not pure* DC: per-window max/mean ratios are 1.092 / 1.069 / 1.057,
  and the two quiet windows drift 3.7% apart. A pure constant DC would make every frame's RMS
  identical (ratio → 1.000) and the drift → 0. So there is real AC content riding on top of whatever
  the offset is, i.e. the endpoint is delivering a genuinely hot floor (gain/boost or ambient), with
  DC as a plausible co-contributor.
- What this leg **ruled out by reading the code** (the tempting "the scale is decoding garbage"
  hypothesis): the capture path does not misread float32 as int16.
  `internal/audio/wasapi_windows.go:363` `	floating := s.format.tag == waveFormatFloat` and `:398-400`
  `	f := unsafe.Slice((*float32)(data), frames*channels)` / `		return MonoDownmix(FloatToPCM16(f), channels)`;
  `FloatToPCM16` (`internal/audio/resample.go:119-128`) scales by 32767/32768 with round-half-away
  and clamps, and `MonoDownmix` (`:102-115`) averages channels. Nothing there can turn silence into
  0.37.
- Downstream consequence worth putting on the record (no action taken by this leg): with a raw floor
  of 0.3685, the ball's own silence gate can never see silence —
  `internal/ball/liquid.go:30` verbatim
  `	SilenceLevelGate   = 0.06 // envelope below this counts as "nobody speaking"`,
  consumed at `internal/ball/liquid.go:215` (`	if raw < SilenceLevelGate {`), i.e. 6× below the floor
  this machine produces. That is a consumer-side decision (`level.go:41-48` says the curve/gain
  belongs there), so it is **not** ticket 247's to fix and this leg changed nothing.

**所以 (b) 的判决**：底噪不是 `level.go` 算错，是"进到尺里的样本本身就有 ~0.37 满刻度的能量（含直流
的可能性已被仓里现成的 DC 夹具演示）"＝**设备/增益/夹具面**；⛔ 修它不许碰 `level.go` 的尺度语义、
不许碰 `:194` 的断言。要判死"DC vs 增益"那一刀，缺的是一枚能把 `mean(samples)` 与 `rms(samples)`
并列报出来的读数件——今天全仓没有（`grep -rn "LevelOfSamples\|FrameLevel" --include=*.go` 的产码
消费者只有 `captureopt.go:105 emitLevel`，它只留 RMS）。

## (c) `frames_sent=6 / frames_dropped=369`: is "接上了" over-claimed?

**No over-claim in any string this leg could find, but "接上了" must be read as "装配可达、电平在测试
里活"，不是"帧有人在吃"，也不是"用户开麦就有数"。**

- The drop is stated in the process's own boot vocabulary, not hidden:
  `cmd/wisp/resident_audio_windows.go:322-323` verbatim
  `	return ra.verdict + "；帧消费侧（ASR/KWS，internal/speech 一块没写）不属于本票，" +` /
  `		"帧按 D38d 计数丢弃并由 meter 限流告警"` — and this leg independently confirms the consumer
  really is absent: `git ls-tree -r --name-only HEAD internal/speech` => `internal/speech/doc.go`.
- The counters are printed at teardown, not just asserted in a test: `:298-299` builds
  `"audio capture leg stopped: levels_delivered=%d frames_sent=%d frames_dropped=%d reopens=%d"`, and
  this leg's own live run printed `levels_delivered=376 frames_sent=6 frames_dropped=370`.
- Levels deliberately flow despite the drop, and that is documented as the contract, not a hack:
  `internal/audio/captureopt.go:96-100` ("The sink runs even when the frame did not reach the
  consumer's channel … a process that has no frame consumer yet … must still get its loudness").
- The stronger "接上了" gap is **not** the PCM side, it is the mute side: at the shipped default no
  number ever arrives, because nothing in production can open the gate —
  `grep -rn "SetMuted" --include=*.go . | grep -v _test` returns only the definition (`gate.go:168`)
  and two comments (`:97`, `:165`), and the batch added no caller
  (`grep -n "^+.*SetMuted" /tmp/247-v1-batch.diff` => empty). Ruler details in
  `20-ac-verdicts.md`, AC#2 section.

## 还欠哪一发（给下一枚腿，一句话＋命令）

机主在场的那一发＝先让门真的能开（今天只有改文件这一条路：`%APPDATA%\wisp\config.toml` 里写
`[audio]` `mic_muted_default = false` 并重启），再跑
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_LIVE_MIC=1 go test ./cmd/wisp/ -run TestAC247LiveMicrophoneLevelsReachTheBallSeam -count=1 -v -timeout 180s`，
在打印 `AC#2 出处: capture=…` 之后**中间 6 秒对麦克风正常说话**；判据＝`:194` 那条断言转绿。
若两枚静默窗仍漂 ~0.014 而说话窗只到 0.4x，那么先做的不是换断言，而是把底噪量出来并归零：
补一枚并列报 `mean(samples)`/`rms(samples)` 的读数件（`audio.DecodeFrame` + 同一帧，写在
`.scratch/wisp/probes/<next>/` 面），据此判"消直流／降增益"该落在设备面还是消费面（`level.go:41-48`
已把曲线判给消费腿）。
