# 247 v1 — questions 3 and 4: the self-added seam, and the `config.NewDefaults()` fallback

Both answered independently by this leg; the dispatch says neither may be defaulted by the
orchestrator, so each carries its own verdict line.

---

## Q3 — `captureSource` / `assembleCapture(…, newSource)`: legal seam, or "mock instead of real"?

**Verdict: a legitimate test seam — it does not violate the README's seam rule and it does not fake
any completion. But it is the reason AC#6 can only be judged 部分成立: what the seam buys is a
*control-flow* proof, never a *device* proof, and the device proof is still owed.**

The seam, verbatim:

- `cmd/wisp/resident_audio_windows.go:149-153`
  ```go
  type captureSource interface {
  	audio.AudioSource
  	Stats() audio.Stats
  	Err() error
  }
  ```
- `:156-158` `func newRealCaptureSource(opts ...audio.CaptureOption) captureSource {` /
  `	return audio.NewWASAPIMicrophone(opts...)`
- `:169-171` `func assembleCapture(rt *proc.Runtime, dataDir string, out func(float32),` /
  `	newSource func(opts ...audio.CaptureOption) captureSource,` / `) *residentAudio {`

Rulers this leg ran to establish that **production never takes the injected branch** (a positive
claim needs no ruler; the wiring chain does):

- `cmd/wisp/main.go:66` `		runResident()` → `cmd/wisp/resident_windows.go:278`
  `	raudio := startResidentAudio(rt, rb.setAudioLevel)` → `resident_audio_windows.go:139-141`
  `startResidentAudio` = `buildResidentAudio(rt, rt.Layout.DataDir, out)` → `:165-167`
  `buildResidentAudio` = `assembleCapture(rt, dataDir, out, newRealCaptureSource)`.
  So the only factory reachable from `main` is the real WASAPI one.
- `grep -n "assembleCapture" cmd/wisp/*.go | grep -v _test` => the two production lines above only;
  the injecting call site is exactly one, in a test:
  `cmd/wisp/resident_audio_247_windows_test.go:251`
  `			ra := assembleCapture(rt, dir, tap.call, func(...audio.CaptureOption) captureSource { return fake })`.
- The other five assembly tests pass `newRealCaptureSource` (`:108`, `:144`, `:186`, `:297`, `:360`),
  and the live AC#2 test does too (`resident_audio_247_live_windows_test.go:147`) then type-asserts
  back to the concrete type (`:151` `	mic := ra.mic.(*audio.WASAPIMicrophone)`).
- The seam is unexported, so no exported semantics moved: AC#0 P4's "⛔ 不许改任何既有导出的语义"
  (`247…md:44`) holds; `startResidentAudio`'s signature is the ruled shape.

Against the README's seam list ("Tests inject at seams only: C8 AudioSource (wav) · C5 LlmProvider ·
C17 PanelBridge · CLI `wisp run`", `.scratch/wisp/issues/README.md:217-218`): this is a **fifth**
injection point, and it is inside the package under test rather than one of the four named edges.
Two facts make it the same species rather than a new licence:

1. `cmd/wisp` already injects at exactly this kind of boundary before ticket 247 —
   `proc.Boot(buildinfo.EnvTest, proc.WithRegistry(observe.NewRegistry()))`
   (`resident_audio_247_windows_test.go:52`, and the same shape in tickets 133/246's tests), and
   `internal/audio` itself has long injected `newWASAPIMicrophoneWith(w DeviceWatcher, o streamOpener)`
   (`wasapimic_windows.go:61`, the line this leg read verbatim
   `func newWASAPIMicrophoneWith(w DeviceWatcher, o streamOpener, opts ...CaptureOption) *WASAPIMicrophone {`) —
   the enumeration seam, ticket 13's own words at
   `capturelevel_windows_test.go:14-17`: "the enumeration seam is injected, the audio-data seam never
   is (SPEC-04 sec 9)". Ticket 247's AC#6 asks a question *about this process's boot*, which cannot
   be asked without a factory seam; the alternative was an exported knob, which P4 forbids.
2. Nothing claimed is a device claim. The three subtests assert only: assembly returns
   (`:253-255`), posture does not say "running" (`:256-258`), the verdict names class + guidance +
   device name (`:259-264`), the boot line repeats the loss (`:265-267`), zero levels reach the seam
   (`:271-273`), step 4 is owned (`:274-282`), the trail stays 10 (`:283-285`).

**"Did any red sentence dress an injected failure up as a real device failure?" — no.** This leg
read every user-visible string the branch produces:
`resident_audio_windows.go:253` `	ra.verdict = fmt.Sprintf("麦克风不可用（分类 %s）：%s；球不会收到电平，本进程继续跑", class, lastErr)`
and `:256-258` (Printf of the class, then the sentence naming D43's legal-edge limit and ticket 128's
single refusal). None of them says "真设备", "本机的", or names a real endpoint; the device string
comes from the error object handed in. The evidence file itself states the limit
(`probes/247/r1/20-ac-readings.md`, "Honest limit: these are assembled failures, not a physically
unplugged microphone") — and this leg confirms by ruler rather than by that sentence:
`git ls-tree -r --name-only HEAD internal/speech` still shows only `doc.go`, so no real consumer or
real device test was added either.

One production-path question the seam *cannot* answer, checked separately because it would have been
a silent hole: does `assembleCapture` risk reading `Stats().LastError` before the capture thread has
posted it (landing in the bland `!gate.Open()` branch at `:259-261` instead of the class branch)?
No — `internal/audio/wasapimic_windows.go:87` creates `started := make(chan error, 1)` and `:98`
returns `<-started`, i.e. `mic.Start` blocks until the thread reported the open outcome, and
`gate.go:247-252 openInnerLocked` posts `meter.fail(err)` inside that call before returning. So the
class branch is reachable in production, not only in the fixture. That is the strongest thing the
seam test *does* prove about the real path.

**AC#6's "任取两形"**: mechanically satisfied twice over (three subtests, all PASS in this leg's
run), but as assembled shapes — hence the 部分成立 verdict in `20-ac-verdicts.md` and the two named
missing shots there (real device refusal; ball + task pipeline booted alongside).

---

## Q4 — `config.NewDefaults()` when `config.toml` cannot be read

**Verdict: not in conflict with ticket 198's existing ruling, and not a new "silent default" shape —
it is loud and it is the conservative table. What it does add is a fourth per-boot default-consumer
with no instrument behind "which table decided the posture", i.e. it inherits exactly the toothless
shape ticket 198 already booked. Whether to keep it is the orchestrator's call; this leg only lays
out 现象 / 影响面 / 撤要动哪几枚文件.**

### 现象 (what the code actually does)

`cmd/wisp/resident_audio_windows.go:179-194`, verbatim head:

```go
	c, _, err := config.LoadFile(cfgPath, nil)
	cfgSource := cfgPath
	if err != nil || c == nil {
		...
		c = config.NewDefaults()
		cfgSource = "compiled defaults (config.NewDefaults)"
		slog.Warn("audio: config unreadable at boot; the capture leg is built from the compiled default table",
			"path", cfgPath, "err", configErrText(err, c),
			"mic_muted_default", c.Audio.MicMutedDefault, "voice_enabled", c.Voice.Enabled)
	}
```

Not silent: the source string is carried into every downstream sentence — `:203` (`…（配置来源 " +
cfgSource + "）"`), `:205` Printf, `:238`
`	ra.verdict = fmt.Sprintf("采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 %s）", cfgSource)`,
and `configErrText` (`:337-345`) distinguishes "err non-nil" from "nil config with nil error".
And it cannot open more than the file would have: `config.NewDefaults()` is the table
`firstrun.go:92` writes, whose `[audio]` half is `mic_muted_default = true` — this leg re-read the
tag at `internal/config/schema.go:277` verbatim
`	MicMutedDefault bool `toml:"mic_muted_default" default:"true"`` and the shipped table is asserted
by `TestAC247ShippedDefaultsAreTheOnesThisLegReads` (PASS in this leg's run) at
`resident_audio_247_windows_test.go:87-98`.

### It is not a new shape in this package

Three pre-existing per-boot fallbacks to compiled values, each loud, all read by this leg:

- `cmd/wisp/resident_windows.go:183-188` — `config.LoadFile` error → `slog.Warn("ball: [hotkey]
  source unreadable at construction; falling back to the compiled defaults", …)` + a Printf naming
  `DefaultHotkeys`.
- `cmd/wisp/resident_approval_windows.go:475-481` — error → Warn naming
  `"fallback", "DefaultApprovalTimeout=300s / DefaultL1Window=3s"` and a *typed* provenance
  sentinel `riskProvenanceUnreadable`.
- `cmd/wisp/panel_resident_windows.go:201-206` — error → Warn `"[panel] geometry source unreadable
  at boot; sizing at the host's own default"`.
- The opposite shape exists too and is confined to a CLI subcommand: `cmd/wisp/providers.go:98-101`
  prints `wisp providers: 配置未就绪（Unconfigured）` and `return 2`.

So 247's branch is the house shape for the resident process. The one respect in which it is
*weaker* than the `[risk]` precedent: `resident_approval_windows.go` carries its fallback as a
**sentinel value** (`riskProvenanceUnreadable`) that the receipt/audit layer can assert on, while
247 carries it only as a **string in the prose**. That is the concrete hole to name.

### 与票 198 既判的关系

Ticket 198's own unsolved cell, verbatim from `.scratch/wisp/issues/198-a-fresh-machine-cannot-run-
wisp-because-nothing-creates-config-toml-done.md:55`: "AC#2 未闭合 … **而回执那句"全部取值来自内置
默认表"在这下发的是假话、无人响**" (mutation V3a: a hand-set non-default value still produced the
"all values come from the built-in default table" sentence, and nothing reddened).

Consequences for 247's branch, stated without deciding for the orchestrator:

- No contradiction: 198's ruling is about the **reload receipt's** sentence; 247's branch is a
  per-boot decision by one subsystem, and it says its own source. AC#4's tree is untouched, so the
  branch cannot be read as changing a default.
- But it widens the same blind spot: after this batch, four subsystems in the resident process can
  act on compiled values (`[hotkey]`, `[risk]`, `[panel]` geometry, `[voice]`+`[audio]`), and no
  instrument ties "which table a subsystem used" to what the 198 receipt claims. A test that
  reddens when the audio leg silently picks the default table while the receipt says "read from
  file" does **not** exist. Ruler for that negative: all 6 `assembleCapture` call sites in the 247
  assembly tests take a directory produced by `writeAudioConfig`, which always writes a file
  (`resident_audio_247_windows_test.go:69-79`, `SaveFile` appears exactly once in the file); no test
  points `dataDir` at a missing `config.toml`. So today the branch is reachable only by an untested
  path — and the asymmetry is worth naming, because the *identical* shape for `[hotkey]` IS
  instrumented: `cmd/wisp/resident_hotkey_258_windows_test.go:131` verbatim
  `		t.Errorf("chain answered the compiled defaults for a config that named other values: the [hotkey] section was ignored")`
  and `:175` asserting the sentinel `hotkeyProvenanceDefaults`. The audio leg carries no such
  sentinel and no such test.
- Privacy direction of the branch is safe *by value*, not *by instrument*: `mic_muted_default=true`
  is why it cannot open a device. The missing instrument is what would catch the day somebody
  changes the compiled table.

### 影响面 (blast radius of keeping or dropping it)

- Kept: a broken/absent `config.toml` yields an armed-and-muted leg + a Warn + a boot line, i.e. the
  same posture as a healthy default install. Nothing in the ball/task pipeline changes.
- Dropped (i.e. "no collector when the file cannot be read"): strictly *more* muted, but it would
  also mean a machine whose file has a syntax error loses the whole capture leg silently-ish at the
  feature level, and ticket 128's refusal list must not grow (`247…md:45`,
  "⛔ 也不许把'拒绝启动'扩大"), so "refuse" is not an available answer here.
- Either way AC#2's missing shot is unchanged: with `mic_muted_default=true` the device never opens,
  and (this leg's ruler) nothing in production can open it later — `SetMuted` has zero call sites.

### 撤要动哪几枚文件

One production file, one contiguous hunk: `cmd/wisp/resident_audio_windows.go:181-194` (the
`if err != nil || c == nil` branch and its Warn) plus, if the source should stay visible after
removal, the three format strings at `:203`, `:205`, `:238` that interpolate `cfgSource`. No test
asserts the branch's existence (zero coverage, above), so no test needs touching; the ticket-255
roster rows at `config_readers_255.go:150`/`:199-202` cite lines **196/197**, which sit *after* the
branch, and `sed -n '196p;197p'` still matches after deleting 181-194 only if the remaining lines
shift by ≤14 — they do shift, so a removal must re-cite the two anchors in the same commit or
`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` reddens (which is the guard working, not a
defect). Recorded so nobody discovers that the hard way.
