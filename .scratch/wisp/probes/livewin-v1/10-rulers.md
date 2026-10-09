# livewin-v1 — 尺名册（每把都现跑，命令原文附后）

取数时刻：2026-10-09 18:1x–18:3x +08，HEAD 起手 `d432d728`（本腿第 1 笔后为 `12b8ae12`），分支 `dev`。
⛔ 本腿零 go 命令（派单禁 build/vet/test，编排者在同一棵树跑整包复跑）。凡"只有跑才能裁"的，见 30 号件的欠账名册。

## R1 凭据本体是不是整族穷举

命令：`ls -la .scratch/wisp/probes/orch/` + `wc -l` 四枚。
读数：findings 101 行／a1-run1 33 行／a1-run2 33 行／c1c2-resident 80 行。同目录另有 `2026-10-09-live-window-runsheet.md`（35,304 字节，**计划件不是读数件**）、`2026-10-09-cmdwisp-panel-rerun.md`、`2026-10-09-cmdwisp-post-290-verify-1.raw.md`（311,834 字节）、`logs/c1-post-1.txt`、`logs/c1-post-2.txt`。
⇒ 派单点名的"三枚原始件"＝`ls` 里 2026-10-09 那天下 a1/c1 前缀的**全部三枚**，整族穷举、不是抽样。findings 的 `:6`–`:10` 自报的字节数（5,015／5,014／17,378）与 `ls` 逐枚对上。

## R2 凭据入库出处

命令：`git log --oneline -1 -- <每一枚>`；`git cat-file -t 73db600a`。
读数：三枚原始件均 `73db600a`（`git cat-file -t` = `commit`）；findings 当前版落在 `d432d728`＝**HEAD 本身**。
⇒ 派单说的"三枚原始件 commit `73db600a`"成立。⚠ findings 的**具名更正那一节**（`:42`）与三枚原始件**不在同一笔**里进盘——原始件 17:47 落、更正 18:0x 落，所以"结论文件比原始件新"这一条是事实，本腿一律以原始件为准、findings 只当作者的读法。

## R3 票面/凭据引的产码行号会不会漂（派单点名"票面现量节的行号会漂"）

命令：`sed -n '<L>p' cmd/wisp/resident_audio_windows.go` 逐行；再 `grep -rn` 现定位。

| 被引的 `文件:行` | 出处 | 现跑结果 |
|---|---|---|
| `resident_audio_windows.go:213` = `gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(...))` | 票 290 现量第 1 条 | **漂**。现量 `:213` = `}`；那句现活在 **`:273`** |
| `resident_audio_windows.go:238` = `采集腿已装配、门处于静音：设备未打开（…）` | 票 290 现量第 5 条＋票 290 `AC#3` 形① 的判据字面 | **漂**。现量 `:238` = `cfgPath := filepath.Join(...)`；那句现活在 **`:298`**（`ra.verdict = fmt.Sprintf(...)`） |
| `resident_audio_windows.go:320`（现量位置）＝`ra.verdict = "采集腿在跑但设备未交接（gate 未 open）：球不会收到电平"` | 票 290 现量第 5 条写作 `:260` | **漂**：现量 `:260` 是注释行，那句现活在 **`:320`** |
| `internal/config/schema.go:277` = `MicMutedDefault bool ... default:"true"` | 票 290 现量第 1 条 | **没漂**，`:277` 逐字命中 |
| `internal/audio/gate.go:97` = `// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it` | 票 290 现量第 1 条 | **没漂**，`:97` 逐字命中 |
| `resident_audio_247_live_windows_test.go:139` = `dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })` | 票 290 现量第 4 条 | **没漂**，`:139` 逐字命中 |
| `resident_audio_247_live_windows_test.go:194`（2× 判据） | findings `:36` | **没漂**：`:194` 逐字 `if loud.mean() <= quietA.mean()*2 \|\| loud.max() <= quietA.max() {`；红句 `:195` = `t.Fatalf("the loud phase is not distinguishable from silence: %s \| %s \| %s"+` |

⇒ 后果具名：**票 290 `AC#3` 形① 的判据文字里那个 `:238` 今天指不到那句产品文案**（文案在 `:298`）。本格要验的是"球打的事还是那句"，句还在、行号漂了 60 行；这一条**不构成判红**，但**构成"票面指针过期"**，编排者勾框时别把 `:238` 当作尺。

## R4 六窗读数与 findings 自报的算术（awk 现跑，非 go）

命令：`awk 'BEGIN{...}'`，输入＝两枚 a1 raw 的 `:18`–`:20`。

| 发/窗 | samples | mean | max | max/mean（本腿算） |
|---|---|---|---|---|
| run1 quiet-A | 93 | 0.367691 | 0.409285 | **1.1131** |
| run1 sound | 189 | 0.382482 | 0.413813 | **1.0819** |
| run1 quiet-B | 94 | 0.383328 | 0.408784 | **1.0664** |
| run2 quiet-A | 93 | 0.369753 | 0.408560 | **1.1050** |
| run2 sound | 189 | 0.382104 | 0.402393 | **1.0531** |
| run2 quiet-B | 93 | 0.381650 | 0.405371 | **1.0622** |

派生比（本腿自算，findings 里没有这三行）：
- run1 `sound/quiet-A` = **1.0402**；run1 `quiet-B/quiet-A` = **1.0425** ⇒ ★**同条件两窗之间的漂移（1.0425）比"响 vs 静"的效应（1.0402）还大**。
- run2 `sound/quiet-A` = **1.0334**；run2 `quiet-B/quiet-A` = **1.0322**。
- 六窗 mean 全距 = 0.367691 … 0.383328（findings/派单说的"挤在 0.367–0.383"**成立**）。
- 六窗 max/mean 全距 = **1.0531 … 1.1131**（findings `:38` 与**票 297 标题**都写作"全挤在 1.08–1.11"⇒ ✗ **三枚窗（1.0664／1.0531／1.0622）在这条带之外**）。

## R5 c1c2 那把 `dropped_frames_total` 尺的语义（产码现读，非推）

命令：`grep -rn "\.push(" internal/audio --include=*.go \| grep -v _test`；`Read internal/audio/audio.go:95-183`；`sed -n '215,250p' internal/audio/wasapimic_windows.go`；`sed -n '90,200p' internal/audio/gate.go`。

1. **谁加的**：`audio.go:126` `m.dropped++`，在 `meter.push` 的 `default:` 支（bounded channel 满 ⇒ 丢帧）；`audio.go:130` `logIt := total == 1 \|\| time.Since(m.lastWarn) >= time.Second` ⇒ 盘上每行是**每秒一发的抽样，打印的是累计值**。
2. **是不是累计量**：是。`dropped` 挂在 `meter`（`audio.go:103`），meter 属于 wasapimic source 对象；source 在装配时建一次（`resident_audio_windows.go:269 newSource(...)`），**gate 开合不重建 meter** ⇒ 盘上 `704 → 722`（关门 20 秒后再开，计数接着走）与代码一致：**跨 reopen 不回零**。
3. **会不会把 boot 期算进去**：不会，且是有条件的不会。`gate.go:99` 起的 `Start` 在关门态**根本不启动 inner source**（`:97` 那句注释；`:108` 才是 `if g.effectiveOpen() {` → `:109 g.open = g.openInnerLocked()`）；帧只在采集循环里产生：`wasapimic_windows.go:232 samples, err := stream.Drain()` → `:239 pending = append(pending, res.Process(samples)...)` → `:240 for len(pending) >= FrameSamples {` → `:241 frame := EncodeFrame(pending[:FrameSamples])` → `:242 m.meter.push(buf, frame)`。⇒ 门没开 ⇒ 没有 `Drain()` ⇒ 没有帧 ⇒ 没有 dropped。**boot 那 19 分钟（17:27:04→17:46:21）盘上零条 drop 行**与此吻合。
4. ★**它到底证到哪一层**：`push` 只在 `stream.Drain()` 交回过样本、且攒够一整帧（32 ms）之后才被调。所以它证的是：**该 WASAPI 采集流已启动、且音频引擎在按设备周期递交包**。它**证不到**"递来的包里有声学内容"——共享模式采集在纯静音房间照样按点交包（内容是近零）。⇒ 派单那一问的答案：**"设备真在采集"这一形今天只证到"设备已交接、包在出"这一层，证不到"麦克风在被有效读"**。而"包里有无人声"恰好就是电平尺那一维，正是票 297 立票在问的那件事。
5. 同一循环里 `wasapimic_windows.go:249 m.cfg.emitLevel(frame)`，其注释 `:245`–`:246` 逐字：`Emitting regardless of the push outcome is deliberate: the drop is a slow-consumer fact, the loudness is a microphone fact` ⇒ `levels=376 / frames_dropped=370 / frames_sent=6`（几乎全丢）**不是缺陷**，且电平确实是从同一串 `Drain()` 出来的帧上算的（不是另一路假数）。

## R6 路径之别：a1 那两发到底跑的是谁的装配

命令：`grep -n "assembleCapture\|WISP_LIVE_MIC" cmd/wisp/resident_audio_247_live_windows_test.go`；`sed -n '253,300p' cmd/wisp/resident_audio_windows.go`；`git grep -n "SetMuted(" HEAD -- internal cmd \| grep -v _test`。

- 台件 `:147` = `ra := assembleCapture(rt, dir, router.route, newRealCaptureSource)` ⇒ **同一个生产装配函数**（exe 起的也是它），差别在：`dir` 是 `:139` 现写的**临时配置目录**（`MicMutedDefault = false`），而 exe 走机主那份 `C:\Users\swq\AppData\Roaming\wisp-dev\config.toml`。台件自己的注释 `:137-138` 也这么说：`assembleCapture reads it through config.LoadFile exactly as the resident boot does`。
- 进程入口之别：a1 两发是 **go test 二进制**（`FAIL github.com/CarlosShao/wisp/cmd/wisp`，raw `:32`）；c1c2 那发是 `wisp 0.0.0-dev (a3782ff0, built 2026-10-09T08:59:56Z)`（raw `:3`）＋ `resident runtime booted`（`:7`）＋ `single instance = true`。
- 生产里真有人开门：`SetMuted(` 非测试命中 **1 枚**＝`cmd/wisp/resident_audio_windows.go:162 gate.SetMuted(!gate.Muted())`（托盘/热键共用的手势处理器，`:172`／`:174` 那两句产品文案逐字命中 raw `:63`／`:38`）。
- 球侧收件人：`resident_audio_windows.go:122 ra.levels.Add(1)`、`:139 rb.b.SetAudioLevel(level)` ⇒ 电平进球那一步在生产代码里，**但常驻路径没有任何一处把电平数值打到盘**（findings `:86` 自己认了"常驻路径零打印器"）。

## R7 tasklist 前置条件有没有读数

命令：`grep -rn -i "tasklist" .scratch/wisp/probes/orch/*.md`；`grep -n -i "tasklist\|wisp.exe\|balldebug" .scratch/wisp/probes/orch/logs/*`；`sed -n '40,50p;78,84p' …runsheet.md`。
读数：`runsheet` 只在**计划**里写了这两把尺（`:45`/`:46` 必须 "No tasks are running"）；台件 `:130` 的 skip 文案也把它写成前提；**四枚凭据件与两枚 logs 里没有任何一发 tasklist 的现量读数**。findings `:92` 作者自陈"我在跑 A1 两发之前没确认盘上没有活的 wisp.exe"。
唯一可推的间接证据：c1c2 raw `:7` 成功起住且带 `single instance = true` ⇒ 17:27:04 那一刻没有第二枚活实例；由此**倒推** 17:25/17:26 两发 a1 期间没有常驻实例（若需要更硬的，那要另跑一发，本腿不跑）。

## R8 机主那份配置的现量（票 290 `AC#3` 形① 与票 247 `AC#4` 的"默认档"落点）

命令：`sed -n '/^\[audio\]/,/^\[/p' "$APPDATA/wisp-dev/config.toml"`；`grep -n -A3 "^\[voice\]"`；`grep -n "wake_word"`。
读数：`[audio] mic_muted_default = true`；`[voice] enabled = true`（`:32`）；`[voice.wake_word]` 节存在（`:38`，本腿未逐字取其中的 enabled）。
⇒ "默认档"这一说法在 `mic_muted_default` 与 `voice.enabled` 两枚上**与编译默认一致**（schema `:277` true／`:245` true），这一条支撑票 247 `AC#4`。
⚠ 但 c1c2 raw 里**没有一处打印 `voice_enabled`**（`grep -c "voice_enabled" = 0`）：`resident_audio_windows.go:252` 那句带 `voice_enabled` 的 Warn 只在**配置文件读不到**的兜底支里走。⇒ 票 247 `AC#4` 要的那枚"默树"读数里，`voice.enabled` 那一半是**由配置文件＋代码分支反推**的，不是那发窗口直接打出来的。
