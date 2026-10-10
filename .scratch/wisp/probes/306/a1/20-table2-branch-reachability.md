# 306-a1 · 20 表②：那两支为什么零执行者（调用链逐跳＋★必答）

世代＝`f994f95`，`git status --porcelain -- internal cmd tools`＝0 ⇒ 下面每一行的行号在**工作树与 HEAD blob 上逐字相同**（blob 尺：`git show HEAD:internal/audio/wavinjector.go | cmp - …` ＝ IDENTICAL）。

## 调用链（从入口到 `:192`／`:218`，逐跳，每跳带内容锚）

1. **入口＝同包测试**：`NewWavInjector(path, opts...)` — 定义 `internal/audio/wavinjector.go:56`：`func NewWavInjector(path string, opts ...InjectorOption) (*WavInjector, error) {`
   全部 12 枚调用点（尺＝`grep -rn --include=*.go "NewWavInjector(" internal | grep -v "func NewWavInjector"`）：
   `captureopt_test.go:71/116/157` · `gate_test.go:16` · `level_test.go:322` · `wavinjector_test.go:83/133/160/223/244/251/271`
   ⇒ **⛔ 一枚生产（非 `_test.go`）调用者**（见下面「响亮报回」）。
2. `wavinjector.go:61`：`	samples, rate, chans, err := parseWav(data)`（`data`＝刚读进来的整个 wav 字节）
3. `wavinjector.go:165`：`func parseWav(data []byte) (samples []int16, rate, chans int, err error) {`
4. `wavinjector.go:188-191`：`fmtTag`/`chans`/`rate`/`bits` 四枚从 `fmt ` 块 body 里逐字节取（`fmtTag = binary.LittleEndian.Uint16(data[body : body+2])`）
5. `wavinjector.go:192`：`			if fmtTag == 0xFFFE {` → 若真则 `:193` 验 `size < 40` → `:196` 把 `fmtTag` 重写成 SubFormat 低字
6. `wavinjector.go:210-227`：`switch {` 三分支 `:211`（1&&16）／`:218`（3&&32）／`:226`（default 报 `unsupported wav format`）

**唯一旁路入口**＝`internal/audio/hotplug_test.go:581`：`	parsed, rate, chans, err := parseWav(raw)`（`raw` 来自 `:576` 的 `writeWav(t, p, samples, TargetRate, 1)`）。

## `fmtTag` 从哪来＝夹具字节构造处（★这是洞的根）

`internal/audio/wavinjector_test.go:18-41` 的 `writeWav` 是全仓唯一的 wav 夹具编码器，它**把标签写死成 PCM16**（内容锚，逐字）：

- `:26` `	buf.Write(u32le(16))` ← `fmt ` 块 size 恒 16（⇒ 即使标签是 0xFFFE 也会撞 `:193` 的 `size < 40`）
- `:27` `	buf.Write(u16le(1)) // PCM` ← **fmtTag 恒 1**
- `:32` `	buf.Write(u16le(16))` ← **bits 恒 16**

`writeWav` 调用点 **11 枚**（尺＝`grep -rn --include=*.go "writeWav(" internal | wc -l`＝11，按文件：`wavinjector_test.go` 5、`captureopt_test.go` 3、`gate_test.go` 1、`hotplug_test.go` 1、`level_test.go` 1）。

唯一的手工第二形＝`TestWavInjectorBadFile`（`wavinjector_test.go:243-274`）里的 8-bit 面：`:259` `	bad = append(bad, u16le(1)...)`（tag 仍 1）＋`:264` `	bad = append(bad, u16le(8)...)`（bits 8）⇒ 只命中 `:226` default。

## ★必答：今天有没有任何夹具真的喂过 `fmtTag==0xFFFE` 且 `bits==32`？

**没有。一枚都没有。** 凭据（⛔ "看起来会走到"，全部是具名用例＋块计数）：

- 尺＝本腿**自己复跑**的覆盖率：`go test ./internal/audio/ -count=1 -coverprofile=.scratch/wisp/probes/306/a1/logs/40-cover.txt`
  → 末行逐字 `ok  	github.com/CarlosShao/wisp/internal/audio	15.960s	coverage: 59.3% of statements`，件内自落 `rc_test=0`（件＝`logs/40-cover-run.txt`）。
- 块级读数（件＝`logs/40-cover.txt`，`wavinjector.go` 行段）：
  - `188.4,192.24 → 1`（走到 `if` 头，只是**比较**，从未为真）
  - `193.5,193.18 → 0` · `194.6,195.1 → 0` · `196.5,196.65 → 0` ⇒ **extensible 支的整个 body 零执行**
  - `219.3,221.26 → 0` · `222.4,223.1 → 0` · `224.3,224.43 → 0` ⇒ **float32 支的整个 body 零执行**
  - `212.3,214.26 → 1` · `217.3,217.35 → 1` · `226.3,226.111 → 1` ⇒ 今天真执行的只有 PCM16 支＋default 支
- 函数级对照（件＝`logs/41-cover-func.txt`）：`wavinjector.go:165: parseWav 70.0%`、`NewWavInjector 100.0%`；同锚点还量到 `wasapi_windows.go:376: Drain 0.0%`、`:409: convertPacket 0.0%`、`parseWaveFormat 81.8%`。
  ⇒ **`parseWav 70.0%`／`Drain 0.0%`／`convertPacket 0.0%` 三数与 `300-v4` 那把尺逐字对上**（派单点名 ⛔ 引它的数当凭据 ⇒ 上一个是本腿独立复跑的那一个）。
- 具名用例清单（都经 `writeWav`，因此全是 tag=1/bits=16）：`TestWavInjectorFrameExact`（`:77`）· `TestWavInjectorRateConversion`（`:123`）· `TestWavInjectorBackpressureDropCounted`（`:151`）· `TestWavInjectorCtxCancel`（`:219`）· `TestWavInjectorBadFile`（`:243`）· `gate_test.go`/`level_test.go:310`（`TestLevelOverBoundedChannelFromWavInjector`）/`captureopt_test.go` 三枚 · `hotplug_test.go:573`（`TestCaptureLoopWavIntegrity`）。
- 反向对照：仓里**确实存在** extensible+float32 的字节面（`parse_wave_format_300_windows_test.go:99` 的 `bits: 32, data1: 3, data1At: 24`），但它喂的是 `wasapi_windows.go` 的 `parseWaveFormat`（`unsafe.Pointer` 侧、`//go:build windows`），**⛔ 同一条字节进过 `parseWav`** ⇒ "两个 parse 同一份 40 字节布局"里只有⛔ 带牙的那一份被钉过。

## 响亮报回（派单点名要这一句）

**`wavinjector.go` 整枚文件今天⛔ 生产调用者**——尺＝`grep -rn --include=*.go "NewWavInjector\|WavInjector" internal cmd tools | grep -v _test.go | grep -v '^internal/audio/wavinjector.go'`，命中只剩两行**注释**：`internal/audio/doc.go:16` 与 `internal/audio/gate.go:29`。
它是 `doc.go:16` 自称的 "the test backbone; ALL pipeline tests"。

⚠ 给编排者的裁定要点（本腿⛔ 裁、⛔ 翻框）：**"没有生产调用者"⛔ 等于"本票要升级成接线票"**——`parseWav` 是包内私有函数、同包测试直接可达（`hotplug_test.go:581` 就是现成的直接调用先例），所以 `AC#1`/`AC#2` 那两枚牙**纯测试侧就能买到**，⛔ 需要任何生产接线。真正需要接线裁定的只有 R2（`Drain`/`convertPacket` 那两枚 0.0%，票面已划出射程）。

## 表②结论（一行）

`:192`／`:218` 两支零执行者的原因＝**唯一的夹具编码器 `writeWav` 把 `fmtTag`/`bits` 写死成 1/16（`wavinjector_test.go:27`/`:32`）、11 枚调用点全走它，而手工第二形只走 default**；今天⛔ 任何夹具喂过 `0xFFFE`＋`bits==32`（凭据＝块计数 0，本腿复跑）。
