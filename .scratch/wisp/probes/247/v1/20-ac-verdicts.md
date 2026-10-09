# 247 v1 — nine-cell verdicts (all readings taken by this leg)

Leg `247-v1`, non-implementer. Clock of this file's rulers: 2026-10-09 13:2x–13:4x +08.
Tree under review = HEAD `b04ec963` (= the batch's last commit; this leg's own commits are
evidence-only and add zero production bytes).
Implementer = `247-r1`. Nothing below is copied from `probes/247/r1/*` as a reading: every
number was re-run here, and the pre-change baseline was re-measured **from git history**
(`git grep … 39d17f0b`) instead of trusting `00-baseline.md`.

Verdicts, one line each; the evidence for each is under it.

| 格 | 判语 |
|---|---|
| AC#1 | **成立** |
| AC#2 | **不成立**（缺的正是"说话 vs 不说话"那一发，且今天这把尺的门槛在本机不可达） |
| AC#3 | **成立**（"球边界"这个读法我独立判为正确，理由在下面；同时票面第二句自身有歧义，具名报回） |
| AC#4 | **部分成立**（缺"默认档真机开机那一发"：今天所有默认档证据都在装配测试面） |
| AC#5 | **成立** |
| AC#6 | **部分成立**（三形都是造出来的失败，且"球与任务管线仍在跑"这半句没有任何仪器） |
| AC#7 | **成立** |
| AC#8 | **成立**（四把尺我逐条自跑：build rc=0、gofumpt 对名册为空、test rc=1 逐名新增红 0、d22scan rc=0） |
| AC#10 | **成立** |

---

## AC#1 — 生产者半有人真调用：**成立**

Ticket sentence (`247…md:24`) verbatim: "接完之后 `internal/audio` 的非测试 importer **≥1 枚**
（尺＝现量 1 那把，接前接后各一发，读数进表），且那枚 importer 就是**跑着的那条腿**
（⛔ 不许是又一枚调试用 cmd）"。

After (my own run of the ticket's own command, `247…md:10`):

```
grep -rl "CarlosShao/wisp/internal/audio" --include=*.go . | grep -v "/internal/audio/" | grep -vc _test
1
```

The one name: `./cmd/wisp/resident_audio_windows.go`.

Before — **re-measured by this leg out of git, not read off the implementer's file**:

```
git grep -l "CarlosShao/wisp/internal/audio" 39d17f0b -- '*.go' | grep -v "/internal/audio/" | grep -v "_test.go" | wc -l
0
git cat-file -e 39d17f0b:cmd/wisp/resident_audio_windows.go
fatal: path 'cmd/wisp/resident_audio_windows.go' exists on disk, but not in '39d17f0b'
```

"就是跑着的那条腿" — the importer is reached from the process entry point, ruler and reading:

- `cmd/wisp/main.go:66` verbatim: `		runResident()`
- `cmd/wisp/resident_windows.go:33` verbatim: `func runResident() {`
- `cmd/wisp/resident_windows.go:278` verbatim: `	raudio := startResidentAudio(rt, rb.setAudioLevel)`
- `cmd/wisp/resident_audio_windows.go:139-141` verbatim:
  `func startResidentAudio(rt *proc.Runtime, out func(float32)) *residentAudio {` /
  `	return buildResidentAudio(rt, rt.Layout.DataDir, out)`

Not a debug cmd: the only other non-test importer candidate is `cmd/balldebug`, and the ruler
above names exactly one file, which lives in `cmd/wisp`.

Edge count (AC#0's "恰好一枚新边"), my own:

```
GOOS=windows go list -deps ./cmd/wisp | grep -cE "wisp/internal/audio$"                 => 1
GOOS=windows go list -deps ./internal/audio | grep -cE "internal/(ball|statemachine)"     => 0
```

**One caveat that belongs to AC#2, not to this cell**: "有人真调用" is true of the *code*, but at
the shipped default nothing can ever open the gate, so a real run delivers zero levels
(`SetMuted` has zero production callers — ruler in the AC#2 section). AC#1's words do not ask for
that, so the cell stands.

---

## AC#2 — 球上的数来自声音，不来自命令行：**不成立**

Ticket sentence (`247…md:25`) verbatim: "一发真机读数＝**对麦克风说话**与不说话时 `SetAudioLevel`
收到的值**不同**，并给出两形的采样数与出处；⛔ 不许用 `balldebug` 的合成包络冒充"。

1. The shape that was measured is **not the shape the box names**. What the harness puts into the
   room is a WAV played out of the render endpoint, not a voice into the capture endpoint —
   `cmd/wisp/resident_audio_247_live_windows_test.go:230-239`, verbatim head:
   `func playAlarmThroughSpeakers(t *testing.T) func() {` …
   `	media := filepath.Join(os.Getenv("WINDIR"), "Media", "Alarm01.wav")` …
   `		cmd := exec.Command("powershell.exe", … "New-Object System.Media.SoundPlayer '" + media + "';" …`
   Whether the bytes left the speakers is not observable from the process; the box asks for
   "对麦克风说话", which is acoustically a different (and observable-by-construction) shape.
2. The box's own reading does not exist. The harness went red by design and its assertion is the
   box's criterion, verbatim at `resident_audio_247_live_windows_test.go:194`:
   `	if loud.mean() <= quietA.mean()*2 || loud.max() <= quietA.max() {`
3. `balldebug` is not used to fake anything: my own ruler over the batch's added lines —
   `grep -cE "^\+.*(balldebug\.|wisp/internal/speech\")" /tmp/247-v1-batch.diff` => `2`, and both
   hits are prose/teardown, not a feed: diff line 58 is a comment
   (`+// shape: not cmd/balldebug's synthetic envelope, …`) and diff line 261 is the ticket's own
   §排程 guard `+	for _, name := range []string{"balldebug.exe", "wisp.exe"}`. The synthetic
   envelope was not reused anywhere, and the live test is gated behind `WISP_LIVE_MIC`
   (`:129-131`, `t.Skip`) — the same env-gate precedent as the pre-existing `TestLiveWasapiSmoke`.
   The `tasklist` precondition in that guard is real and this leg re-ran it before its own live run:
   `tasklist //FI "IMAGENAME eq balldebug.exe"` and `… eq wisp.exe` both answered
   `INFO: No tasks are running which match the specified criteria.` (13:2x and again at the live run).
4. **And on this machine the bar as written is not reachable by speaking** — arithmetic on the
   scale this repo defines (see `30-ac2-level-ruler.md` for the full derivation): with the floor
   this machine reports, `loud.mean > 2 × 0.3691 = 0.7382` needs an added AC component of
   `0.6393` of full scale = **-3.89 dBFS at the microphone**. Normal speech cannot do that; only
   clipping-distance shouting or an electronic feed could. So the next leg must fix the *floor*
   (or the shape of the discriminator) before re-aiming at this box, and must not touch the
   assertion to make it green.

**还欠哪一发（一句可照做的命令，机主在场）**：先把 `[audio] mic_muted_default = false` 写进
`%APPDATA%\wisp\config.toml` 并重启（今天生产里没有任何路径能把门打开，见下面 AC#2 附属尺），
然后跑
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_LIVE_MIC=1 go test ./cmd/wisp/ -run TestAC247LiveMicrophoneLevelsReachTheBallSeam -count=1 -v -timeout 180s`
并在打印 `AC#2 出处: capture=…` 之后的**中间 6 秒对着麦克风说话**；判据＝`:194` 那条断言转绿
（sound 窗 mean 超过两枚静默窗 mean 的 2 倍且 max 高于 quiet-A）。若仍红，红的正确解释不是"没说话"
而是"底噪 0.37 满刻度没被消掉"，那一发是设备/增益面的（见 `30-ac2-level-ruler.md`）。

**附带尺（本腿自己跑的否定句）**：这条腿在 shipped 默认档下**永远收不到电平**，因为把门打开的那枚
调用者在产码里不存在：

```
grep -rn "SetMuted" --include=*.go . | grep -v _test
./internal/audio/gate.go:97   (comment: "source is NOT started; SetMuted(false)/SetSpeaking(false) will start it")
./internal/audio/gate.go:165  (comment: "SetMuted is driven by the mute hotkey.")
./internal/audio/gate.go:168  func (g *HalfDuplexGate) SetMuted(muted bool) {
```
=> three hits, **zero call sites** (`gate_test.go` only). And the batch added none:
`grep -n "^+.*SetMuted" /tmp/247-v1-batch.diff` => empty. The mute hotkey that `gate.go:165`
believes drives it only records a gesture — `cmd/wisp/resident_ball_windows.go:281` verbatim:
`			OnMuteHotkey:    func() { recordBallGesture("mute-hotkey") },`
and `raudio` is used for exactly one thing in `runResident`
(`resident_windows.go:295` prints `raudio.posture()`).

This is not a breach of AC#4 (which demands the default *not* open the mic) and nobody wrote the
opposite claim: the boot line says 设备未打开 (`resident_audio_windows.go:238`). But "接上了" today
means "装配可达、电平路径在测试里活"，不是"用户开麦就有数"，这一句应当由编排者写进程序，不由本腿改任何默认值。

---

## AC#3 — 只过一枚 `float32`：**成立**

Ticket sentence (`247…md:26`) verbatim: "任何新增代码里，跨过 `internal/ball` 边界的数据**只许是一个
标量电平**。尺＝能力型（新增的跨包调用签名里出现 `[]byte`／`[]int16`／samples 即判越界），
⛔ 不许做成"扫注释里有没有 samples 这个词"的词面型"。

Every cross-package call the batch added, enumerated by this leg (not copied from `10-gates.md` §4).
Ruler: `git diff --no-renames -U0 39d17f0b..b04ec963` (1542 lines) plus a selector sweep of the new file.

| 新增跨包调用 | 形参类型 | 跨进 `internal/ball`？ |
|---|---|---|
| `cmd/wisp/resident_audio_windows.go:157` `audio.NewWASAPIMicrophone(opts...)` | `...CaptureOption` | 否 |
| `:210` `audio.WithCaptureRegistry(rt.Registry)` | `*observe.Registry` | 否 |
| `:211` `audio.WithLevelSink(ra.levelOut(out))` | `func(float32)` | 否 |
| `:213` `audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(...))` | `captureSource, Path, CaptureOption` | 否 |
| `:216` `audio.NewBoundedFrames()` | — | 否 |
| `:222` `gate.Start(ctx, buf)` | `context.Context, chan []byte`（票 13 既有 C8 契约） | 否 |
| `:232` `mic.Stats()` / `:246` `mic.Err()` / `:296` `gate.Stop()` | 返回值 `audio.Stats` / `error` | 否 |
| `internal/audio/wasapimic_windows.go:88` `SpawnCapture(m.cfg.Registry, func(c context.Context))` | `*observe.Registry, func(ctx)` | 否（observe 早就是 audio 的依赖） |
| `internal/audio/wavinjector.go:93` `SpawnCapture(w.cfg.Registry, …)` | 同上 | 否 |
| `cmd/wisp/resident_audio_windows.go:132` `rb.b.SetAudioLevel(level)` | `float32` | **是，一枚标量** |

The ball-side signature, quoted from disk before judging:
`internal/ball/liquid_windows.go:42` verbatim `func (b *Ball) SetAudioLevel(level float32) {` and
`:40-41` verbatim "This is the only render-side audio input in the project: no samples and no
transcript text cross into the ball (C25 contamination surface)."

`[]byte` appears in exactly one new cross-package call (`gate.Start`), and that call crosses
`cmd/wisp → internal/audio`, never `→ internal/ball`: ruler
`GOOS=windows go list -deps ./internal/audio | grep -cE "internal/(ball|statemachine)"` => `0`.

**My independent judgement of the reading** (the implementer flagged it, so it is mine to answer, not
theirs): the box's *subject* is the ball boundary — sentence 1 names it — and sentence 2 states the
instrument's shape (capability, not word-level). Two further facts make the wide reading untenable:
(a) the orchestrator's own 更正 section (`247…md:88`) says verbatim "**本票 `AC#3`"只过一枚
`float32"`的判据形状正好由这枚签名自己钉着**" and points at `liquid_windows.go:42`; (b) under the
literal wide reading, AC#1 (≥1 non-test importer that is the running leg) and AC#3 cannot both be
satisfied, because the only way for `cmd/wisp` to drive a C8 source is
`AudioSource.Start(ctx, chan<- []byte)` (`internal/audio/audio.go`, the ticket-13 contract) —
the box would forbid the落点 AC#0 ruled at `247…md:77`. So: **the reading is correct**; AC#3 成立.
The residual defect is in the ticket text, not in the wiring — the second sentence should name its
scope, and that is a wording fix for the orchestrator, which this leg does not make.

---

## AC#4 — 隐私默树一格都不许动：**部分成立**

Ticket sentence (`247…md:27`) verbatim: "`voice.enabled`／`wake_word.enabled` 的**默认值与读取处**
一字不改，且真机读数要证明"默认档下麦克风不会被这条腿打开"（或反过来具名说出它打开了、由哪一行决定）。"

The tree did not move — my own rulers:

```
git diff --name-only --no-renames 39d17f0b..b04ec963 -- internal/config        => (empty)
sed -n '197p;245p;277p' internal/config/schema.go
	Enabled bool `toml:"enabled" default:"false"`
	Enabled bool `toml:"enabled" default:"true"`
	MicMutedDefault bool `toml:"mic_muted_default" default:"true"`
```

Read sites: `wake_word.enabled` gained **no** production reader (`grep -n "^+.*c.Voice.WakeWord"
/tmp/247-v1-batch.diff` shows the only two occurrences at diff lines 395-396, both inside
`TestAC247ShippedDefaultsAreTheOnesThisLegReads` asserting the default is false). `voice.enabled`
and `mic_muted_default` gained **new** readers, which is what AC#0 P1 甲 ordered
(`247…md:41`: "必须把 `c.Audio.MicMutedDefault` 传给 … `WithStartMuted`，并把 `voice.enabled=false`
判为根本不构造采集器"). Read-site *edit* count = 0; the pre-existing reload readers under
`internal/config/manager.go` are untouched because the whole package is.

The default does not open the device — `TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice`
**PASS in my own run** (`--- PASS: TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice (0.01s)`),
whose assertions this leg read verbatim (`resident_audio_247_windows_test.go:149-172`): gate
`Muted()` true, `Open()` false, `ra.started` false, `st.FramesSent != 0 || st.BytesSent != 0` →
fatal, `rt.Registry.CountByName("audio-capture") != 0` → fatal, and the boot sentence must contain
both `设备未打开` and `mic_muted_default`.

**缺的那一发**：that is an *assembly* reading on a real source constructor, not a 真机读数. No run of
the shipped default on this machine exists (`mic_muted_default=false` is precisely what the only
live test writes at `resident_audio_247_live_windows_test.go:139`). Missing shot, one sentence:
以 shipped 默认档起一次真常驻进程（`GOFLAGS= go build -o build/wisp.exe ./cmd/wisp` 后跑该 exe 的常驻
入口），把 boot 行 `wisp: 采集腿：…设备未打开…` 与 Windows 设置→隐私→麦克风 的"最近应用访问"指示
一起抄进证据件——今天没有任何一枚腿跑过这一发（机主在场才能看那个指示器，属编排者的欠账，不是腿的）。

---

## AC#5 — 协程有 owner、有 recover：**成立**

Ticket sentence (`247…md:28`) verbatim: "⛔ 裸 `go func(`（ban #1 …），采集协程必须挂在现成的 owner 上
（`observe.Root`／`rt.Registry`），退出序列里被 join——**不许新造 D38(e) 的第 11 步，十步顺序一字不动**。"

- Bare `go func(`: `grep -rn "go func(" cmd/wisp/resident_audio_windows.go internal/audio/captureopt.go
  internal/audio/wasapimic_windows.go internal/audio/wavinjector.go cmd/wisp/resident_windows.go
  cmd/wisp/config_readers_255.go` => **no match** (rc=1). Over the whole added diff:
  `grep -n "^+.*go func(" /tmp/247-v1-batch.diff` => empty. Independently, d22scan (ban #1) rc=0 clean.
- Owner: `internal/audio/audio.go:180-182` verbatim
  `func SpawnCapture(registry *observe.Registry, run func(ctx context.Context)) *observe.Handle {` /
  `	return registry.Spawn("audio-capture", "audio", nil, run)` — the two new call sites are
  `wasapimic_windows.go:88` and `wavinjector.go:93`, both inside the diff's added lines.
- Recover: `internal/observe/goroutine.go:296` verbatim `		if rec := recover(); rec != nil {`, inside
  `func (r *Registry) run(...)` (`:285`), which is the body `Spawn` launches — so the capture
  goroutine inherits owner+recover from the sanctioned entry.
- D38(b) roster zero inflation: `internal/observe` is untouched
  (`git diff --name-only --no-renames 39d17f0b..b04ec963 -- internal/observe` => empty);
  `internal/observe/goroutine.go:41-44` verbatim:
  `// ResidentBaseline is the frozen resident goroutine budget (D38b: ui-sta,` /
  `// audio-capture, hotkey-listener, db-writer, watchdog, log-flusher).` /
  `const ResidentBaseline = 6` / `var ResidentNames = []string{` `	"ui-sta", "audio-capture", …`.
  The level is computed on that same goroutine (`captureopt.go:101-111 emitLevel`, called from
  `wasapimic_windows.go` inside the frame loop and from `wavinjector.go:138` in `pump`), so the only
  name ever spawned by the batch is `audio-capture`, already on the roster.
  Roster-warning path (`goroutine.go:270-273` `slog.Warn("goroutine outside the D38 roster …")`) never
  fires for this batch — no new name exists to fire it.
- One registry, not two: `internal/audio/captureopt.go:78-87` seeds
  `CaptureConfig{Registry: observe.Default}` and an injected registry **replaces** it; asserted live in
  my run — `TestAC247RealCaptureLoopEmitsLevels` PASS and `TestAC247CaptureRegistryReplacesTheDefault` PASS.
- Join, bounded: `internal/audio/wasapimic_windows.go:116` verbatim `	tm := observe.NewTimeout(2 * time.Second)`
  (a `observe.NewTimeout`, not a wall-clock subtraction — ban #4).
- D38(e) ten steps: `internal/proc` untouched (`git diff --name-only … -- internal/proc` => empty);
  `StepStopAudio` is index 3 = step 4 (`internal/proc/shutdown.go:40`, and
  `internal/proc/shutdown_test.go:62` verbatim `	if order[3] != StepStopAudio || order[4] != StepReleaseSpeechSessions {`).
  `TestAC247CaptureLegOwnsStepFourOfTheFrozenOrder` **PASS in my run**; its own assertions walk ten
  names in order and reject an 11th (`resident_audio_247_windows_test.go:302-314`).
- Registration site: `cmd/wisp/resident_audio_windows.go:276` verbatim
  `	if err := rt.RegisterShutdownHook(proc.StepStopAudio, func(context.Context) error {` — and a leg
  that builds nothing does not register (`:202-207` returns before it), pinned by
  `TestAC247VoiceDisabledBuildsNoCollector` (PASS in my run).

---

## AC#6 — 降级照跑：**部分成立**

Ticket sentence (`247…md:29`) verbatim: "拔麦／无权限／设备被占三形里**任取两形**，本进程**仍要把球与
任务管线跑起来**并且把损失说响亮（现量：票 128 定的是"没有数据根才拒绝启动"，别把它扩大）。"

What exists, re-run by this leg: `TestAC247DeviceFailureShapesStillBootAndSayTheLoss` PASS with all
three subtests in **my** run (`occupied` / `permission_denied` / `no_device`, 0.01s each), so 任取两形
is covered twice over. The errors are built with the package's own constructor
(`resident_audio_247_windows_test.go:247` verbatim
`			want := audio.DeviceError(tc.hr, dev, "capture open")`), i.e. **assembled**, not a device that
actually refused; the device names are fixtures (`"Busy Mic"`, `"Built-in Mic"`, `"Gone Headset"`,
`:234-236`).

The loud half is real: the verdict must carry the class and the guidance and the device name
(`:259-264`), the boot posture must repeat the guidance (`:265-267`), and production prints it —
`cmd/wisp/resident_audio_windows.go:253-258`, verbatim
`	ra.verdict = fmt.Sprintf("麦克风不可用（分类 %s）：%s；球不会收到电平，本进程继续跑", class, lastErr)` /
`		fmt.Printf("wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，" +`
No refusal was widened: `grep -c "os.Exit" cmd/wisp/resident_audio_windows.go` => `0`, and
`internal/statemachine` is untouched (`git diff --name-only --no-renames 39d17f0b..b04ec963 --
internal/statemachine` => empty). This leg read the D43 row itself rather than trusting the ticket:
`internal/statemachine/table.go:97` verbatim
`		D43: 14, From: StateListening, Event: EvAudioDeviceLost, To: StateError,`
(and `:55` verbatim `		D43: 5, From: StateSleeping, Event: EvKwsEnabled, To: StateArmed,` — the
`Armed` entry stays as unreachable as AC#0 P2 甲 said, and the batch pushed no state).

**缺的两发**（写成下一枚腿能照做的句子）：
① "三形"today 只在装配面成立，真机一发＝在 Windows 设置→隐私→麦克风 关掉本应用访问（或在另一进程以
独占模式占住 Realtek 输入设备）后起一次真常驻进程，把 boot 行的 HRESULT 原样抄进证据件——票面
`AC#0 编排者裁决` 节自己也把它列为"本格不闭合的三寸"第①条（`247…md:79`），所以这是**已知欠账**，
但格子要闭合就得有人跑那一发。
② "本进程仍要把球与任务管线跑起来"今天**没有任何仪器**：尺＝
`grep -n "任务管线\|startResidentTaskSource\|residentBall" cmd/wisp/resident_audio_247_windows_test.go`
=> 零命中；该族用例只断言 `assembleCapture` 返回、门存在、电平为 0，从不装配球与任务来源。
下一发＝在同一枚用例里（或新增一枚）用同一棵 `rt` 把 `startResidentBall` 与
`startResidentTaskSource` 一起起来，断言设备失败那形里两者仍在。

---

## AC#7 — 越界检查：**成立**

My own ruler, both ranges, `--no-renames` (the trap the orchestrator booked at `A772`):

```
git diff --name-only --no-renames 39d17f0b..b04ec963 | wc -l                     => 13
git diff --name-only --no-renames 39d17f0b..b04ec963 | grep -Ec \
  "^(internal/speech|scripts/spike|frontend/|design/|PLAN\.md|docs/specs/)|thresholds\.go|golden|allowlist\.txt|tokens_fourway_test\.go|l2_grant_boundary_test\.go|ticket90_persist_test\.go"
                                                                                 => 0   (grep rc=1)
same ruler over 2304ca12^..b04ec963 (widest range, includes the two orchestrator commits) => 0
```

The 13 names: the 10 code/test files + the ticket's Progress-log line + `probes/247/r1/*.md` (3).
`internal/speech` is still one file — my own ruler
`git ls-tree -r --name-only HEAD internal/speech` => `internal/speech/doc.go` (and the same at
`39d17f0b`) — so nothing was smuggled in there. The 4 added lines that mention `internal/speech`
are all prose naming the missing consumer (diff lines 738, 1005, 1239, 1364; 1005 is the production
boot sentence `	return ra.verdict + "；帧消费侧（ASR/KWS，internal/speech 一块没写）不属于本票，" +`),
and `grep -c` for an added `wisp/internal/speech` import => `0`. 唤醒词／ASR／TTS stayed out.

---

## AC#8 — 门禁四数：**成立**

See `10-gates.md` for the four commands, each with its own `rc=N`, all run by this leg. Summary:
`GOFLAGS= go build ./...` rc=0 (zero output) · gofumpt over the batch's 10 files rc=0 and **empty**
(dir-scoped it lists 3 files that are not in this batch) · the ticket's test command rc=1 with
`PASS=374 FAIL=6 SKIP=3` in the `-v` form and **逐名新增红 0 枚** against the dispatch baseline ·
`./tools/d22scan/d22scan.exe` rc=0 `clean - no D22 ban violations` (and `sh scripts/d22scan.sh` rc=0).

---

## AC#10 — 接上了不等于看得见：**成立**

Ticket sentence (`247…md:33`) verbatim (judgement half): "本格判据＝票面与本仓文档**不许**把"电平接好"
写成"用户能看见球在动" … ⛔ 本腿不许顺手 `ball.EnablePrototypeVisuals(true)` 来让 AC#2 好看".

- The flag is off by construction: `internal/ball/statevisual.go:102` verbatim
  `var prototypeVisuals bool`; `:109` verbatim `func PrototypeVisualsEnabled() bool { return prototypeVisuals }`.
- The seam really does return early: `internal/ball/liquid_windows.go:52-54` verbatim
  `func (b *Ball) applyLevel(level float32) {` / `	if !prototypeVisuals {` / `		return`.
- Nobody turned it on: my own caller ruler (a negative claim, so a call-site ruler)
  `grep -rn "EnablePrototypeVisuals" --include=*.go . | grep -v "_test.go"` =>
  `cmd/balldebug/main.go:122:  balle.EnablePrototypeVisuals(!*frozen)` (verbatim
  `	ball.EnablePrototypeVisuals(!*frozen)`) + only comments/definition inside `internal/ball`.
  And over the batch's added lines: `grep -n "^+.*EnablePrototypeVisuals(" /tmp/247-v1-batch.diff` => empty.
- The reading exists: `TestAC247HandingTheLevelToTheSeamIsNotVisibility` **PASS in my run**, and my
  log line (not the implementer's) is
  `resident_audio_247_windows_test.go:372: AC#10 reading: PrototypeVisualsEnabled()=false levels_reaching_the_ball_seam=1 (a number arriving is not a pixel moving)`.
- No text in the batch claims visibility: `grep -n "^+.*\(看得见\|呼吸\|breathe\)" /tmp/247-v1-batch.diff`
  => two hits, both negations (`:651` "not claim the orb breathes on screen", `:731` "the user sees
  the ball breathe" preceded by `What this file does NOT claim`). The production boot line at
  `resident_windows.go:295` prints only the posture sentence.
