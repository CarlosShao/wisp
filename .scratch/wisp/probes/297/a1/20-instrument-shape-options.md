# 297-a1 · 20 仪器形状候选（票面 `AC#0` 那一枚"打喂进电平尺的到底是哪串字节"）

锚＝`e96da9f4`。**本件只列形状与代价，⛔ 不是实现，⛔ 我没写任何测试文件、没跑任何测试。**
每形末尾必答那一问："能不能把 `0.3677` 归因到三支里的某一支"（三支＝①声音没到麦／②设备本底真这么高／③字节切错；票面 `AC#2` 的编号口径里 ①②③ 顺序与票面 `:26-28` 一致——⚠ **票面 `AC#2` 把"字节切错"编成 ②、"设备本底"编成 ③，与票面开头 `现量` 节的叙述顺序不同，我下面一律用**名字**不用号，避免撞号**）。

---

## ⓐ 真设备在程用例（票面 `AC#0` 直读那一形）

### 已有的 skip 形状名册（尺 R8＝`git grep -n "t.Skip" e96da9f4 -- '*.go'`，射程＝全仓 `.go` 含注释行，现量 **119 枚命中行**，按文件计头名如下；尺 R16＝对下表 5 枚件逐枚 `git show e96da9f4:<file> | grep -n -B5 -A5 "t\.Skip"` 现读上下文）

| # | 件（内容锚＝函数名） | 门是什么 | 跳过时打印什么 |
|---|---|---|---|
| 1 | `cmd/wisp/resident_audio_247_live_windows_test.go` 的 `TestAC247LiveMicrophoneLevelsReachTheBallSeam` | `if os.Getenv("WISP_LIVE_MIC") != "1" { t.Skip(...) }` | 逐字理由串：`"AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0"` |
| 2 | 同件的 `playAlarmThroughSpeakers` | `if _, err := os.Stat(media); err != nil { t.Skipf(...) }`（`media` = `%WINDIR%\Media\Alarm01.wav`） | 逐字：`"no alarm sample to play out loud (%v): AC#2 then needs a person speaking into the microphone"` |
| 3 | `internal/audio/hotplug_test.go` 的 `TestLiveWasapiSmoke` | 同 `WISP_LIVE_MIC != "1"` | 逐字：`"live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone (真机冒烟待票 16)"` |
| 4 | `internal/audio/hotplug_test.go` 的 `TestPinnedThreadStable10s` | `if testing.Short()` | `"10s pinned-thread window (acceptance criterion) skipped in -short"` |
| 5 | `internal/ball/live_windows_test.go` 两处（hover 投递失败／`SetCursorPos` 被拒） | 运行期探测不到输入 | **`SKIP-LOUD:` 前缀那一族**，见下 |
| 6 | `internal/models/manifest_real_test.go` 两处 | `WISP_IT_REAL_MIRROR` 未设 | `"real-network spot check; set WISP_IT_REAL_MIRROR=1 to run"` |

（⚠ 这 6 组是本腿**逐枚读了上下文**的；R8 的全仓 119 枚命中分布在 **61 个文件**（尺＝`cut -d: -f2 /tmp/297a1-skip.txt | sort | uniq -c | sort -rn`，现量），我没逐枚读，只点名与音频／真机／"跳过≠通过"有关的这几族。）

### "跳过≠通过"在本仓靠什么保证（这是 ⓐ 形状的生死题，三条都是现量）

1. **`SKIP-LOUD` 文案规矩**（尺 R16，逐字块，`internal/ball/live_windows_test.go`，整块未省行）：
```go
			t.Skipf("SKIP-LOUD: SetCursorPos refused (the pointer is owned elsewhere): asked (%d,%d), at (%d,%d). "+
				"This run therefore did NOT exercise the real WM_MOUSEMOVE/WM_MOUSELEAVE path. "+
				"The pop-back BEHAVIOR is proven by the non-skipping deterministic test "+
				"TestDockHoverPopBackWalksTheRampHome in dock_test.go, which drives the same "+
				"dockRampFrame step law the hover and leave handlers call; this test only proves the "+
				"WIRING delivers the event, and it proved nothing this run.",
```
   ⇒ 定式＝**跳过句里必须自己写清"这一发没证什么"＋点名那枚不跳过的确定性用例谁在扛**。⚠ 自查留痕：这块我第一次落笔时把第 4 行抄成了 `step law the hover and leave calls;`（吞掉一枚 `handlers`），已在同笔 commit 前改回工作树原文；同一函数体内**第二枚** `t.Skipf`（R16 现量，`hover undeliverable this run` 那一枚）写的是 `step law these handlers call; what this test could not show today is the wiring.`——两枚文案不同，⛔ 引这形时别把它们混成一枚。
2. **CI 侧的硬闸**：`scripts/portable-tests.sh` 头部规矩 2 逐字（尺 R18＝`grep -nE "scope|--tags|go test|SKIP|skip|-run" scripts/portable-tests.sh`，现量）：
```
#   2. any top-level `--- SKIP` in what actually ran is fatal;
```
   且 `scripts/portable-tests.sh:703` 一带 `if [ "$skipped" -ne 0 ]; then` 打印"每条 SKIP 的 file:line 与它自己写的理由"。**⇒ 一条会 skip 的新用例如果落进 CI 执行面、又没进那枚 SKIP ledger，CI 直接红**（ledger 见 40 件 §3）。
3. **`WISP_LIVE_MIC=1` 那族**今天**已经在 ledger 里有名**（`internal/audio` 那枚 `TestLiveWasapiSmoke`，逐字行见 40 件 §3）⇒ 说明本仓处理"真机用例在 CI 里不许假绿也不许假红"的既有办法＝**登记进 ledger 让 `-skip` 摘走**，⛔ 不是让它跑到那儿 skip。

### ⓐ 怎么做才不造恒绿（结论，⛔ 我没实现）

- 门的形状照 1／3 名册抄：`os.Getenv("WISP_LIVE_MIC") != "1"` ⇒ skip，**且 skip 文案必须走 `SKIP-LOUD` 定式**（写明"这一发没测到 ⓐⓑⓒⓓ 里的哪几枚"）。
- ⛔ **不许**把票面 `AC#0` 要的 ⓑ（前 16 个 `int16`）写成"可选打印"——票面判据是"四行都有"。
- ⚠ ⓐ 有个**票面没写的第二枚恒绿口**：上面第 2 行那枚 `t.Skipf`（`Alarm01.wav` 不在）。它 skip 的是**外放那一步**，不是整枚用例 ⇒ 用例照跑、照读窗、照不红，只是 `sound` 窗变成"没人说话"。票面 `AC#0` 若复用这台件，**这一枚 skip 会把"有声音"悄悄降级成"没声音"**。⇒ 落地腿要么把它改成 `t.Fatal`，要么在仪器输出里把"外放到底开没开"打成一行。这是本腿对票面的一条**新增缺陷指认**（算的，凭 §1 名册第 2 行的现量内容锚）。
- **归因能力**：⊙ 能分"字节切错"与另外两支（ⓑ 呈大/小交替或量级恒在 12,000–15,500 ⇒ 切错）；⊙ 能分"设备本底"（ⓑ 平滑小信号＋ⓓ RMS 真≈0.37）；⛔ **不能**单独分"声音没到麦"——那一支要靠 ⓓ 的窗间对比＋机主在场，票面自己也是这么排的。**⚠ 而它要独占真设备窗口（票面排程段逐字："297-a1 要真设备 ⇒ 必须独占 `cmd/wisp` 与音频设备面"），这是它唯一但很贵的代价。**

---

## ⓑ 注入接缝＝`C8 AudioSource`／wav（票面 `AC#1` 那一形）

### 仓里现成 wav 夹具名册（尺 R9＝`find . -name '*.wav' -not -path './.git/*'`，射程＝**整个工作树**，现量 **19 枚**）

- `third_party/model-fetch/` — **12 枚**
- `third_party/spike-models/` — **7 枚**
- ⇒ **⛔ 零枚在 `testfixtures/` / `testdata/` / `internal/audio/` 下**（尺＝同一发的目录分布 `cut -d/ -f2-3 | sort | uniq -c`）。全部 19 枚都在 `third_party/`（根 `.gitignore:12` 逐字 `third_party/` ⇒ **未跟踪、由 `scripts/fetch-deps.ps1` 现拉**）。
- ⚠ 后果（算的，凭上面两条现量）：**本仓没有一枚"已知幅度的 wav"是跟踪夹具**。`AC#1` 要的"0.05 与 0.4 同长窗"两枚夹具**要么新写生成器、要么落进 `internal/audio` 的测试里用 `math.Sin` 手搓**。票面 `AC#1` 写的"走注入接缝（`C8 AudioSource`＝wav，仓里现成：`internal/audio/wavinjector.go:52`、`internal/audio/capturelevel_windows_test.go:19 TestAC247RealCaptureLoopEmitsLevels`）"——**"仓里现成"指的是那台件（注入环），不是那两枚 wav**；引用时别把这两件事混成一件事。

### ⚠⚠ ⓑ/ⓒ 这一族的**射程天花板**（本件对票面最重要的一条不同意）

`capturelevel_windows_test.go` 那枚"真采集环＋假 opener＋真 wav 字节"的台件，**它的假对象是 `fakeStream`，而 `fakeStream` 自己实现了 `Drain()`**（尺 R12 现量：`internal/audio/hotplug_test.go:108: func (s *fakeStream) Drain() ([]int16, error)`）。逐字（整块，工作树读，`internal/audio/capturelevel_windows_test.go:19-38`）：

```go
func TestAC247RealCaptureLoopEmitsLevels(t *testing.T) {
	const frames = 8
	samples := make([]int16, frames*FrameSamples)
	for i := range samples {
		samples[i] = 12345 // a steady DC level: RMS is exactly this magnitude
	}
	dev := DeviceDescriptor{ID: "dev-level", Name: "Level Fixture Mic"}
	watcher := &fakeWatcher{devs: []DeviceDescriptor{dev}}
	stream := &fakeStream{
		dev: dev, rate: TargetRate, period: 10 * time.Millisecond,
		chunks: makeChunks(samples, 160),
	}
	opener := &fakeOpener{streams: []*fakeStream{stream}}

	rec := &levelRecorder{}
	reg := observe.NewRegistry()
	mic := newWASAPIMicrophoneWith(watcher, opener,
		WithCaptureRegistry(reg),
		WithLevelSink(rec.call),
	)
	baseDefault := observe.Default.CountByName("audio-capture")
```

⇒ **算的**：注入面从 `fakeStream.Drain()` **已经交出 `[]int16`** 那一刻起就绕过了 `convertPacket`／`MonoDownmix`／`FloatToPCM16` 全链。所以：
- ⊙ ⓑ 能证的：`emitLevel`→`FrameLevel`→`LevelOfSamples` 这一段**对已知幅度分得开**（票面 `AC#1` 判据"同一把尺在注入面能分开、在真麦面分不开"这条成对读数，**ⓑ 恰好就是它的"注入面"那一半**，这一格 ⓑ 够用、且不需要机主）。
- ⛔ **ⓑ 不能把 0.3677 归因到三支里的任何一支**。它连嫌疑函数都没进。票面 `AC#1` 的末句"拿到之后才允许说'问题在设备/转换那一侧，不在 `LevelOfSamples`'"——**注意它只把 `LevelOfSamples` 摘干净了，"设备"与"转换"两支在 ⓑ 之后仍然缠在一起**。⇒ 谁要只交 ⓑ 就宣称"归因完成"，是对票面判据的读宽。
- ⚠ 另一条既有事实对 ⓑ 不利：上面块里 `samples[i] = 12345 // a steady DC level` ⇒ **这台件今天见证的正是"含直流的 RMS"**（票 291 编排者裁定段 `:35` 已写"那枚夹具从此失去见证力"是真代价、`want` 是派生式所以不会红）。若 `AC#2` 裁成"改尺去直流"，ⓑ 的期望值不红但**它测的东西搬家**。

---

## ⓒ 包内假设备（要把 `convertPacket` 切进去得换哪几枚调用）

尺＝R12＋逐字读 `device.go`/`wasapi_windows.go`/`wasapimic_windows.go` 的接缝声明。要伪造"真设备"到能喂 `convertPacket`，需替换的调用**恰 2 枚**（算的，枚数尺＝`grep -c comCall` 在 `Drain` 函数体内＝**3**：slot 5 `GetNextPacketSize`、slot 3 `GetBuffer`、slot 4 `ReleaseBuffer`，加上 `s.capture` 这枚 COM 指针本身）：

1. `wasapiOpener.Open`（`internal/audio/wasapi_windows.go:210`）→ 已有接缝：`newWASAPIMicrophoneWith(watcher, opener, …)` 收 `streamOpener` 接口（`device.go:63` 声明 `Drain`）。**但这枚接缝交给测试的是"整个 stream"，不是"stream 里的 capture client"** ⇒ 用它就只能像 `fakeStream` 那样**替掉 `Drain`**，仍然进不了 `convertPacket`。
2. `comCall(s.capture, 3, …)`（`GetBuffer`）这一层若要保真 `Drain`，就得给 `s.capture` 造一枚**能响 vtable slot 3/4/5 的假 COM 对象**（`comCall` 逐字：`vt := *(*unsafe.Pointer)(this); fn := *(*unsafe.Pointer)(unsafe.Add(vt, slot*8))`）⇒ 手工摆一枚 vtable 数组＋三个 `syscall.SyscallN` 兼容的桩。**代价：高，且形状在本仓零先例**（尺＝`git grep -n "comCall(" e96da9f4 -- '*_test.go'` 我没跑；但 R12 的 `Drain(` 名册里唯一的测试侧实现是 `fakeStream` 那种"整层替掉"，**没有任何一枚件造过假 COM vtable** ⇒ 算的）。

⇒ **⊙ 结论（本腿给票面补的那一形，最便宜）**：ⓑ 与 ⓒ 之间其实还有一形 **ⓓ 包内纯函数假字节**——直接在 `internal/audio` 写一枚 `//go:build windows` 用例，`buf := make([]float32, frames*channels)`，`convertPacket(unsafe.Pointer(&buf[0]), frames, channels, false, false)`，再 `LevelOfSamples` 读它。
- **要替换的调用＝0 枚**，要造的假设备＝0 枚，要机主＝0，要独占 `cmd/wisp`＝0，要 third_party DLL＝0（不在 `cmd/wisp` 包）。
- ⊙ **它能做到的**：把"字节切错会不会产出 `mean≈0.3677 / max/mean≈1.08 / 对声音几乎无感"这一**形状**在仓内**正向复现**（如果复现得上，形"切错"就从"不可能被测到"升成"有仪器"）；同发还能把 §4 里 `parseWaveFormat` 的 `+24/+26` 之争钉死（手搓 40 字节 extensible 块断言 `tag`）。
- ⛔ **它做不到的**：**证明真机走的就是那一支**。它只证"这一支能造出这形状"，不证"这形状来自这一支"（别的支也可能造出同形状）。⇒ **ⓓ 与 ⓐ 是互补不是替代**：ⓓ 负责"可达＋判别式"，ⓐ 负责"真机那一发到底是哪一支"。票面 `AC#0` 的 ⓑ／ⓒ 两枚打印项（前 16 个 `int16`／同段字节按 float32 重读）我在 ⓐ 里原样保留，⛔ 不许因为做了 ⓓ 就省掉。
- ⚠ ⓓ 的一条**必须先由编排者裁的规矩问题**：它要新增一枚 `internal/audio` 包内 `//go:build windows` 测试文件。票面 `AC#0` 原文写的是"新增一枚**真设备**在程用例"＋"只产**测试侧**仪器，⛔ 不动产码语义" ⇒ ⓓ 落在"测试侧仪器"里合法，但**落在"AC#0 这一格"里不合法**（票面那格的判据是"四行都有"，指真机四行）。⇒ **本腿不自作主张把 ⓓ 记成 AC#0 的交付**，只把它作为候选交回，请编排者在 `AC#1`/`AC#2` 之间给它一格。

---

## 归因能力总表（本件唯一要交的判语）

| 形 | 分"切错字节" | 分"设备本底" | 分"声音没到麦" | 要机主 | 要独占 cmd/wisp＋麦 | 新夹具 |
|---|---|---|---|---|---|---|
| ⓐ 真设备在程 | 能（ⓑⓒ 两行对照） | 能 | **单发不能** | **要** | **要** | 否 |
| ⓑ wav 注入 | **完全不能**（进不了 `convertPacket`） | 不能 | 不能 | 否 | 否 | **要**（零枚跟踪 wav） |
| ⓒ 包内假 COM | 能（但代价最高） | 不能 | 不能 | 否 | 否 | 否 |
| ⓓ 包内纯函数假字节（本腿新增候选） | 能证"这形造得出这读数"，不能证"就是它" | 不能 | 不能 | 否 | 否 | 否 |

⛔ 本腿不翻框、不派单、不选形；`AC#2` 的"三形择一"按票面原文**由读数裁、不由写腿挑**。
