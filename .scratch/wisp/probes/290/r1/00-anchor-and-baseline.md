# 票 290 / 290-r1 — 00 起手锚与改前读数（落地腿＝甲-1：只把门拧开）

工作目录 `D:\work\workspace\projects plans\Wisp`，分支 `dev`。
起手钟点 `2026-10-09 14:57:09 +0800`（`date` 的 stdout，未手打）。

## 1. 锚

```
$ git log --oneline -3
6c976aa6 票 289 / 289-v1 裁决腿第 3 笔（判语表＋名册＋AC#3 正文）：四格全判**成立**…
fcf1d449 票 289 / 289-v1 裁决腿正文：AC#1 成立…
e7af1796 票 289 / 289-v1 裁决腿起手：骨架件入库…
```

HEAD＝`6c976aa6`（派单预期「`6c976aa6` 或其后」，命中第一枚）。票 289 的三枚腿笔都在其下，
⇒ 本腿按票面排程「`290-r1` 按在 `289-r1` 交回之后」成立，起手无阻塞。

在飞件（⛔ 不还原、不提交、不据以判绿）：`M .gitignore`、`M .scratch/**` 若干、
`D design/**` 16 枚、`?? .scratch/**` 若干——本腿一律未动。

## 2. AC#2 那把「枚数」尺的改前读数（本腿现跑，尺＝票面 AC#2 逐字那条）

```
$ git grep -n "SetMuted(\|SetSpeaking(" HEAD -- internal cmd | grep -v _test
HEAD:internal/audio/gate.go:97:// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it
HEAD:internal/audio/gate.go:130:func (g *HalfDuplexGate) SetSpeaking(speaking bool) {
HEAD:internal/audio/gate.go:168:func (g *HalfDuplexGate) SetMuted(muted bool) {

$ … | wc -l
3
```

三枚全不是调用者：`:97` 是注释行、`:130`/`:168` 是两枚**定义行** ⇒ **非测试调用者＝0**，
票面前提成立。

⚠ **具名分歧（以原文为准，不改票面一字）**：派单与我读到的票面 AC#2 写的是「由 **2 枚**变 ≥3 枚」，
而这条尺（`SetMuted(` **与** `SetSpeaking(` 合起来）在 `6c976aa6` 上现跑＝**3 枚**。
「2 枚」这个数只在**只量 `SetMuted(`** 那把尺上成立（＝定义 `:168` ＋注释 `:97` 两枚；
`SetSpeaking` 那族自己贡献 `:130` ＋同一句注释 `:97`）。票面现量段 `:12`/`:12` 两行其实也是分开数的
（SetMuted 命中的是「它自己的定义行＋那句注释」）。⇒ 本件把**两把尺都记下**，
改后判据写作：**合尺 3→≥4、且新命中的那枚在 `cmd/wisp` 的生产链路上**。

同一把尺对工作树（不带 `HEAD`）跑一次，读数同为 3 ⇒ 起手时工作树的 Go 面与 HEAD 无差异。

## 3. 三枚默认值的现读（⛔ 本票一字不动；这是"我确认我没动"的凭据，不是结论）

| 键 | 出处 | 现读逐字 |
|---|---|---|
| `voice.enabled` | `internal/config/schema.go:245` | `	Enabled bool \`toml:"enabled" default:"true"\`` |
| `wake_word.enabled` | `internal/config/schema.go:197` | `	Enabled bool \`toml:"enabled" default:"false"\`` |
| `mic_muted_default` | `internal/config/schema.go:277` | `	MicMutedDefault bool \`toml:"mic_muted_default" default:"true"\`` |
| 语义那句（持久化那一问的现产码答案） | `internal/config/schema.go:276` | `	// MicMutedDefault starts every session muted.` |

读取处（本腿不改，两处，改后名册里也不该出现）：`cmd/wisp/resident_audio_windows.go:196`
`ra.mutedAtBoot = c.Audio.MicMutedDefault`、`:213`
`gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`。

## 4. 真相源那句（产码注释，逐字；本腿的主从依据，⛔ 不许倒过来）

`internal/audio/gate.go:20-21`：

```
//   - Muted (both paths, user intent wins): capture closed; Muted/Unmuted
//     events hook the mute hotkey path into the Muted state (ticket 07).
```

⇒ **门侧那枚 `muted` 是主**（`gate.muted`，读口 `gate.Muted()` `:202`，写口 `SetMuted` `:168`）。
本腿不新增任何第二枚"谁在静音"的布尔、⛔ 不 dispatch `EvMuteKey`、⛔ 不碰 `internal/statemachine/table.go`。

`internal/audio/gate.go:97`（静音位上设备根本不打开，本票要接的就是这一句的另一半）：

```
// Start implements AudioSource. When the gate is (currently) closed the inner
// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it
// later with the stored ctx/buf, so Start still returns nil.
```

## 5. 装配顺序（甲-1 的"后置 setter"为什么要后置）

- `cmd/wisp/resident_windows.go:217` 建球（`startResidentBall(...)`，事件闭包在这一刻定型）
- `cmd/wisp/resident_windows.go:278` 建采集腿（`raudio := startResidentAudio(rt, rb.setAudioLevel)`，门在这一刻才存在）

⇒ 球的 `OnMuteHotkey`/`OnTrayMute` 两枚闭包**必须**在 fire 时刻才去取那枚门（本进程已有的
`rb` 指针 → 后置注入的执行者），不能在 `ball.New` 时刻取。

## 6. 改前整包名册（两发，`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -v ./cmd/wisp/ ./internal/audio/ -count=1`）

原始件＝同目录 `raw-pre-1.md`／`raw-pre-2.md`（`.md` 后缀，⛔ 非 `.out`：根 `.gitignore:7` 忽略 `*.out`）。
计数尺＝`grep -c -- '--- PASS'` / `'--- FAIL'` / `'--- SKIP'`（含子测试，两发同一支 harness，⛔ 未加 `-skip`）。

| 发次 | PASS | FAIL | SKIP | `go test` rc | 包级行 | `0xc0000135` | `0xc000013a` |
|---|---|---|---|---|---|---|---|
| pre run 1 | 418 | 5 | 3 | rc=1 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 464.624s`／`ok github.com/CarlosShao/wisp/internal/audio 15.945s` | 0 | 0 |
| pre run 2 | 418 | 5 | 3 | rc=1 | `FAIL …/cmd/wisp 465.462s`／`ok …/internal/audio 16.226s` | 0 | 0 |

改前红名册逐名（两发**完全同一组**，本票改前两发就红）：
1. `TestPanelHostRealWindowHopAndLifecycle` (5.06s / 5.06s)
2. `TestAC4FocusReturnToPriorWindowGap33r5` (5.11s / 5.12s)
3. `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` (20.01s / 20.01s)
4. `TestAC14AwaitedBindingReplyReachesThePage` (20.01s / 20.01s)
5. `TestAC14GoSideEvalPushReachesThePage` (20.01s / 20.01s)

⇒ 全 5 枚都是面板/WebView2 那一族（与 289-v1 记的 panel 族红同族），`internal/audio` 两发都 `ok`。

改前三枚 SKIP 逐名（⛔ 不是读数，只是没跑）：
`TestPanelHostLatencyPercentilesAC2`、`TestAC247LiveMicrophoneLevelsReachTheBallSeam`（`WISP_LIVE_MIC=1` 才跑）、
`TestLiveWasapiSmoke`。

⚠ 基线具名交代**用了哪支 harness**＝带 sherpa DLL 前缀那一支（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`），
⇒ 两发 `0xc0000135` 计数都是 0（派单说的环境红没出现）；本票两发也**没有** 289-v1 那枚 `0xc000013a`
（`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 两发都绿）。
⚠ 派单让本件把 `internal/panel` 也入册：本腿没有单跑 `./internal/panel/`，
理由＝票面 AC#5 与派单第 5 段逐字给的门禁尺只有 `./cmd/wisp/ ./internal/audio/` 两包，
⛔ 不自作主张扩包集（要并入 `internal/panel` 请由裁决腿定），故此处具名欠这一枚、未改尺。


## 7. 真机相关的前置量（进程数必须为 0）

```
$ tasklist //FI "IMAGENAME eq balldebug.exe" | grep -c "balldebug.exe"
0        (rc=1，grep 无命中＝0 枚在跑)
$ tasklist //FI "IMAGENAME eq wisp.exe" | grep -c "wisp.exe"
0        (rc=1，同上)
```
