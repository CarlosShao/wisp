# 247 r1 — per-AC readings

Leg `247-r1`, 2026-10-09 12:2x–13:1x +08. Gate rc lines and the red-roster diff are in
`10-gates.md`; this file is the per-box reading. No AC box was ticked by this leg.

## AC#1 — the producer half has a real caller

- ruler (the ticket's own command): non-test importers of `internal/audio` **0 before / 1 after**;
  the one name is `cmd/wisp/resident_audio_windows.go`, i.e. the resident leg, not a debug cmd.
- edge rulers: `GOOS=windows go list -deps ./cmd/wisp` now contains
  `github.com/CarlosShao/wisp/internal/audio`; `go list -deps ./internal/audio` still matches
  `internal/ball|internal/statemachine` **0** times. One new edge, exactly the one AC#0 ruled.

## AC#2 — NOT CLOSED BY THIS LEG (named, with what was obtained)

What the harness produced on this machine (real device, real pinned thread; run 12:4x with
`WISP_LIVE_MIC=1`, log excerpt verbatim):

```
level=INFO msg="audio: capture leg running" path=T frame=32ms
AC#2 出处: capture="麦克风 (Realtek(R) Audio)" render="扬声器 (Realtek(R) Audio)" pinned_thread=14772
AC#2 counters at the seam: levels=375 frames_sent=6 frames_dropped=369 last_err=""
AC#2 reading: quiet-A: samples=93 mean=0.369079 max=0.418274 levels_per_s=31.00
AC#2 reading: sound: samples=189 mean=0.382697 max=0.416027 levels_per_s=31.50
AC#2 reading: quiet-B: samples=93 mean=0.382564 max=0.404375 levels_per_s=31.00
--- FAIL: TestAC247LiveMicrophoneLevelsReachTheBallSeam (12.48s)
```

So the microphone opened, frames came, and 375 scalars were handed to the seam. **The box is
still not satisfied**: the loud window is not distinguishable from the quiet ones
(mean 0.3827 vs 0.3691, and its max is *lower* than quiet-A's). The loud window was
`C:\Windows\Media\Alarm01.wav` played out of the render endpoint
(`扬声器 (Realtek(R) Audio)`) by a `System.Media.SoundPlayer` loop; the same player returned
`PLAYED-OK` when run alone, so the bytes reached the renderer — whether they left the speakers
is what this leg cannot see or control. No synthetic envelope, no wav injector, no skip was
used to make this look like a pass; the test went RED by design.

Two numbers worth the adjudicator's attention, both from that same run:

1. The **floor** is ~0.37 of full scale while nothing is being said. `level.go` states the
   scale is unweighted RMS that does **not** remove DC ("DC counts: RMS does not remove a DC
   offset"), so a capture path carrying a DC offset reads high and barely moves — that is a
   property of this machine's capture endpoint plus the raw scale, not of the plumbing.
2. `frames_sent=6 / frames_dropped=369`: the bounded C8 channel has **no consumer in this
   process** (`internal/speech` is still only `doc.go`, and 唤醒词／ASR／TTS are outside this
   ticket). The levels flow anyway, which is what `emitLevel`'s contract states.

The one reading AC#2 still needs, for the owner (or anybody speaking into the machine):

```sh
cd "D:/work/workspace/projects plans/Wisp"
tasklist //FI "IMAGENAME eq balldebug.exe"    # must answer "No tasks running"
tasklist //FI "IMAGENAME eq wisp.exe"         # must answer "No tasks running"
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_LIVE_MIC=1 \
  go test ./cmd/wisp/ -run TestAC247LiveMicrophoneLevelsReachTheBallSeam -count=1 -v -timeout 180s
```

Then: after the line `AC#2 出处: capture=...` prints, **speak into the microphone for the
middle 6 seconds** (the three windows are quiet 3s / sound 6s / quiet 3s). The box closes when
the `sound` window's mean and max are clearly above both quiet windows, which is exactly the
assertion that failed above.

## AC#3 — only a float32 crosses into the ball

Every new cross-package signature is in `10-gates.md` §4. The only call this leg added that
enters `internal/ball` is `residentBall.setAudioLevel(float32) -> Ball.SetAudioLevel(float32)`
(the signature at `internal/ball/liquid_windows.go:42`, which is the shape the 更正 section
names as the judge). No `[]byte` / `[]int16` / samples parameter crosses that boundary. The
pre-existing C8 seam `AudioSource.Start(ctx, chan<- []byte)` stays inside
cmd/wisp <-> internal/audio and is not a ball crossing; that reading is flagged in the leg
report as an interpretation of AC#3's second sentence.

## AC#4 — the privacy tree did not move, and the default does not open the mic

- Ruler `git diff --name-only 39d17f0b HEAD | grep -c internal/config` => `0`. The three
  default lines are still verbatim what the ticket's 更正 section cites:
  `schema.go:197` wake_word `default:"false"`, `:245` voice `default:"true"`,
  `:277` `mic_muted_default default:"true"`. No read site of those keys was edited either; the
  only new read sites are this leg's (a *new* reader is what P1 甲 demanded, and the ticket-255
  roster in `ada5c563` registers them).
- Device-not-opened reading, `TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice` (PASS):
  with the shipped table the gate is `Muted()==true`, `Open()==false`,
  `mic.Stats().FramesSent == 0`, and `rt.Registry.CountByName("audio-capture") == 0` — the
  capture thread is never spawned, so `wasapi_windows.go`'s synchronous open never runs.
  The boot sentence names the key that decided it (posture contains both `设备未打开` and
  `mic_muted_default`).
- `voice.enabled=false` reading, `TestAC247VoiceDisabledBuildsNoCollector` (PASS): no gate, no
  source, no thread, and **no step-4 registrant** (a leg that owns no thread does not claim the
  step).
- Side effect worth recording: `MicMutedDefault` was listed by the 更正 section as a 哑键 with
  zero production readers. After this leg it has one (`resident_audio_windows.go:196`); the
  ticket's own P1 甲 asked for exactly that ("map it here at boot wiring", gate.go:55-56).

## AC#5 — owner, recover, join, and the frozen ten steps

- The thread is spawned through `audio.SpawnCapture` (audio.go:180, first production caller)
  into `rt.Registry` — no bare `go func(` anywhere in this leg's files; `d22scan` rc=0 with
  `bans #1-5 internal/=229, cmd/=39` examined.
- One registry only, proven live and in tests: `rt.Registry.CountByName("audio-capture")==1`
  while `observe.Default.CountByName(...)==0` (live log + `TestAC247CaptureRegistryReplacesTheDefault`).
- Teardown: `rt.RegisterShutdownHook(proc.StepStopAudio, ra.stop)` — the existing slot, no
  eleventh step. `TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder` (PASS) walks the ten records
  and checks every name in order; the live run shows step 4 in its place between the skipped
  neighbours:

```
step=3 name=cancel-task-roots skipped
msg="audio capture leg stopped: levels_delivered=376 frames_sent=6 frames_dropped=370 reopens=0"
step=5 name=release-speech-sessions skipped
```

  (`mic.Stop()` carries the bounded 2s join, `wasapimic_windows.go:94-115`; the thread id was
  back to 0 after it, asserted in the live test before it failed on the loud phase.)
- Roster inflation: none — the level is computed inside the same `audio-capture` goroutine
  (P8 甲). P8's cost clause reading: `levels_per_s = 31.00 / 31.50` against the nominal
  `1/FrameDuration = 31.25`, i.e. delivery on the pinned thread costs nothing measurable.

## AC#6 — the degraded shapes, two of three (this leg did all three)

`TestAC247DeviceFailureShapesStillBootAndSayTheLoss` — subtests `occupied`,
`permission_denied`, `no_device`, all PASS. Each asserts: the assembly returns (the boot is not
stopped), the verdict names the class `audio_device` **and** the guidance sentence and the
device name, no level reaches the seam, step 4 is still owned, and the trail stays 10 records.
The errors are built by the package's own `audio.DeviceError(0x8889000A / 0x80070005 /
0x88890004, dev, ...)`, so the guidance under test is the text at `device.go:86/:89/:103`.
No state was pushed: `git diff --name-only 39d17f0b HEAD | grep -c internal/statemachine` => `0`
and `internal/statemachine/table.go:97` still reads verbatim
`D43: 14, From: StateListening, Event: EvAudioDeviceLost, To: StateError,`.

Honest limit: these are assembled failures, not a physically unplugged microphone. The verbatim
HRESULT this machine's real device answers when it is genuinely occupied is one of the three
"本格不闭合的三寸" the ticket already carries.

## AC#7 — path roster

`git diff --name-only 39d17f0b HEAD` => the 10 files listed in `10-gates.md` §5; the
forbidden-path grep over that list returns **0**. `frontend/**` and `design/**` were neither
read nor cited (the only `design/**` mention in this file is the *existence* check
`git cat-file -e HEAD:design/assets/tokens.css`, which reads no content).

## AC#8 — see `10-gates.md` §1

## AC#10 — connected is not visible

`TestAC247HandingTheLevelToTheSeamIsNotVisibility` (PASS), log line:

```
AC#10 reading: PrototypeVisualsEnabled()=false levels_reaching_the_ball_seam=1 (a number arriving is not a pixel moving)
```

`ball.EnablePrototypeVisuals` was not called by this leg anywhere:
`grep -rn "EnablePrototypeVisuals" cmd/ internal/ | grep -v statevisual.go` after this leg still
lists only `cmd/balldebug/main.go:122` in production and internal/ball's own tests.

## What this leg added that the ticket did not name (flagged, not hidden)

1. `cmd/wisp/config_readers_255.go` — the ticket-255 roster instrument reddens on any new
   production read site of a config section it adjudicates (`TestTicket255RosterStillMatchesThe
   ActualReadSites`). Adding `[audio]`/`[voice]` reader entries is that instrument's own
   maintenance path; no assertion was weakened and both verdicts stay in the
   `other-process`/no-reader class rather than `consumed`.
2. `captureSource` / `assembleCapture(rt, dataDir, out, newSource)` — an unexported seam inside
   cmd/wisp so AC#6's three device failures are testable at the process level. No exported
   symbol changed semantics; `startResidentAudio(rt, out)` keeps the ruled shape.
3. `config.NewDefaults()` as the fallback when `config.toml` is unreadable at boot, said out
   loud (`cfgSource = "compiled defaults (config.NewDefaults)"` + a Warn + nothing opened). The
   ticket does not cover a missing file; the fallback is the repo's own compiled table (the one
   `firstrun.go:92` writes) and it is the conservative side, since that table's
   `mic_muted_default` is `true`. Reported for ruling, not presented as settled.
4. The C8 frame channel this leg opens has no consumer (see AC#2 point 2), so D38d's drop
   counter ticks and `meter.push` warns at most once per second for as long as the leg runs.
   That is the honest interim shape and the boot line says it; a quieter option would be a
   silent drain, which D38d forbids.
