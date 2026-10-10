# 297-a1 · 10 `convertPacket` 可达名册

锚＝`e96da9f4`（锚漂声明见 00 件 §1.1）。本件所有 `.go` 读数的 provenance＝**工作树读**，且同发尺 `git status --porcelain -- internal/audio/` 输出为**空**（现量，见 00 件 R6 后那一发）⇒ 工作树逐字等同 `e96da9f4`。对象层复跑的那几发标 `git show e96da9f4:<path>`。

---

## 1. 票面那句"3 命中全在非测试文件"——我复跑了，**命中数核上了，结论没核上**

### 1.1 尺与射程

| 尺号 | 逐字命令 | 射程目录 | 现量/算的 |
|---|---|---|---|
| R6 | `grep -rn "convertPacket" --include='*.go' internal cmd` | `internal/`＋`cmd/`，含 `_test.go`，含注释行 | 现量 |
| R5 | `git grep -n "convertPacket" e96da9f4` | **全仓**（含 `.scratch/`、`docs/`），对象层 | 现量 |

R6 现量结果＝**3 枚命中，全部在 `internal/audio/wasapi_windows.go`，零枚 `_test.go`**：

```
internal/audio/wasapi_windows.go:385:			out = append(out, convertPacket(data, int(frames), ch, floating, flags&audclntBufferFlagsSilent != 0)...)
internal/audio/wasapi_windows.go:393:// convertPacket turns one interleaved device packet into mono int16.
internal/audio/wasapi_windows.go:394:func convertPacket(data unsafe.Pointer, frames, channels int, floating, silent bool) []int16 {
```

（逐字，整块 3 行，未省行。）

⚠ 一处与票面的**可复现性差别**：票面把尺写作 `--include=*.go`（未加引号）。Git Bash 下未加引号的 `*.go` 会先被 shell 展开。⇒ 我补了一发尺判这两形是否等价：`ls *.go` 现量回 **NONE**（仓根零 `.go`）⇒ **在本仓根目录下两形等价**；换一枚仓根有 `.go` 的工作树就不等价。我跑的是加引号那一形。

R5（全仓、对象层）现量＝**12 枚命中**：上面 3 枚 `.go` ＋ 9 枚**文档/票面/台账里的提及**（票 297 自身 5 枚、`probes/240/c1/census.md:76`、`probes/orchestrator/ci223-1/raw_range_roster.txt:228`、`docs/reports/pending-and-issues.md:14913` 与 `:14984`）。⇒ 票面那条尺的射程**只看 `.go`**，这点票面自己写了（`internal cmd`），我核上；但**要提醒落地腿**：全仓射程里 `convertPacket` 已经被两枚前腿件当作"结论"引用过（`240-c1`、`ci223-1` 名册），那些是** prose，不是仪器**。

### 1.2 我没核上的那一半：票面的推论

票面 `:16` 逐字（从工作树读，内容锚＝该行开头 `★**\`convertPacket\` 的射程＝今天没有任何仪器跑过它**`）：

> ⛔ 零 `_test.go`）。⇒ **它唯一的生产调用点在真设备上**，所以"切错字节"这一形**在仓内不可能被测到**。

我的判语：**前半对，后半今天不成立，而且它正是本票最贵的一处误导。**

- **前半对**：`convertPacket` 的**唯一调用点**是 `wasapiStream.Drain()`（尺 R6，射程＝`internal cmd` 全 `.go`，现量）。而 `Drain()` 的调用名册（尺 R12＝`git grep -n "Drain(" e96da9f4 -- '*.go'`，现量，`internal/llm` 那 20+ 枚命中是同名异物的 `adaptertest.Drain`，已具名排除）只有 4 枚：
  1. `internal/audio/device.go:63` — `deviceStream` 接口声明 `Drain() (samples []int16, err error)`；
  2. `internal/audio/wasapi_windows.go:361` — 真实现（windows tag）；
  3. `internal/audio/wasapimic_windows.go:232` — **唯一生产调用者**：`samples, err := stream.Drain()`；
  4. `internal/audio/hotplug_test.go:108` — `func (s *fakeStream) Drain() ([]int16, error)`，**测试侧替身**。
- **后半不成立**：`convertPacket` 是**包内非方法纯函数**，签名 `func convertPacket(data unsafe.Pointer, frames, channels int, floating, silent bool) []int16`。它的 `data` 参数**不要求来自 COM**——任何 `//go:build windows` 的包内用例都能拿一枚 Go 自己 `make` 的缓冲区取址喂进去（`unsafe.Pointer(&buf[0])`），**不需要麦克风、不需要真设备、不需要 opener 接缝**。⇒ 形③"切错字节"**在仓内完全可测**；它今天零仪器是因为**没人写**，不是因为**测不到**。票面那句"不可能被测到"如果进了落地腿的脑子，那条腿就会去抢真设备窗口（本票排程里最贵的资源），而这是不必要的。

⚠ 顺带一枚对票面有利的更正，防我被读成"票面全错"：**真设备不可达的只有 `Drain()` 这一层**（它要 COM vtable 指针，测试里造不出合法 vtable）。票面把"Drain 不可测"写成了"convertPacket 不可测"，射程宽了一格。

---

## 2. 进入 `convertPacket` 的输入是从哪几枚变量来的（不是自由参数）

`convertPacket` 的 4 枚实参在 `Drain()` 里全部由**流自己**决定，名册如下（逐字块＝`git show e96da9f4:internal/audio/wasapi_windows.go`，函数 `Drain`，整块未省行）：

```go
func (s *wasapiStream) Drain() ([]int16, error) {
	ch := int(s.format.channels)
	floating := s.format.tag == waveFormatFloat
	var out []int16
	for {
		var packet uint32
		if hr, _, _ := comCall(s.capture, 5, uintptr(unsafe.Pointer(&packet))); !hrOK(hr) {
			return out, DeviceError(hr, s.dev, "GetNextPacketSize")
		}
		if packet == 0 {
			return out, nil
		}
		var (
			data           unsafe.Pointer
			frames, flags  uint32
			devPos, qpcPos uint64
		)
		if hr, _, _ := comCall(s.capture, 3,
			uintptr(unsafe.Pointer(&data)), uintptr(unsafe.Pointer(&frames)),
			uintptr(unsafe.Pointer(&flags)),
			uintptr(unsafe.Pointer(&devPos)), uintptr(unsafe.Pointer(&qpcPos))); !hrOK(hr) {
			return out, DeviceError(hr, s.dev, "GetBuffer")
		}
		if frames > 0 {
			out = append(out, convertPacket(data, int(frames), ch, floating, flags&audclntBufferFlagsSilent != 0)...)
		}
		if hr, _, _ := comCall(s.capture, 4, uintptr(frames)); !hrOK(hr) {
			return out, DeviceError(hr, s.dev, "ReleaseBuffer")
		}
	}
}
```

⇒ 四条读法（全部＝**算的**，由上面那块逐字推）：
1. **`frames == 0` 永远进不了 `convertPacket`**（调用点上有 `if frames > 0` 卫）。⇒ 落地腿若为"零帧"写分支，那一条**在本仓不可达**，⛔ 不许算作覆盖了形③。
2. `channels` ＝ `s.format.channels`，**每枚流固定**，来自 `parseWaveFormat`（下节）。
3. `floating` ＝ **只看 `s.format.tag` 一枚数**，`tag == 3`（`waveFormatFloat`）。
4. `silent` ＝ `flags & audclntBufferFlagsSilent != 0`，常量 `audclntBufferFlagsSilent = 0x2`（逐字见 `wasapi_windows.go:103`）。
5. ⚠⚠ **`s.format.bits` 在这一层从头到尾没有被读过**（尺＝`git grep -n "format.bits" e96da9f4 -- '*.go'` 我没跑；但射程内可读的是上面整块 `Drain`＋下节 `convertPacket` 全文，两块里 `bits` 出现 **0 次** ⇒ 这条是**码层直读＋全文可复核**，不是 grep 现量，落地腿要引用请先自己现跑 `grep -c bits` 在那两个函数体内）。

---

## 3. 分支名册：`(tag/bits/channels/frames/silent)` 各组合落到哪一支

`convertPacket` 全文（逐字，整块，从工作树读，函数定义在 `internal/audio/wasapi_windows.go:393-404`，含注释行）：

```go
// convertPacket turns one interleaved device packet into mono int16.
func convertPacket(data unsafe.Pointer, frames, channels int, floating, silent bool) []int16 {
	if silent || data == nil {
		return make([]int16, frames)
	}
	if floating {
		f := unsafe.Slice((*float32)(data), frames*channels)
		return MonoDownmix(FloatToPCM16(f), channels)
	}
	i := unsafe.Slice((*int16)(data), frames*channels)
	return MonoDownmix(i, channels)
}
```

⇒ **恰 3 支出路**（尺＝逐字块数 `return`，射程＝该函数体）：

| 支 | 触发条件 | 输出长度（**算的**） | 归因能力 |
|---|---|---|---|
| A `silent \|\| data == nil` | `flags` 带 `0x2`，**或** COM 回了 nil 指针 | `frames` 枚（全零） | 与 0.3677 **互斥**，见 §5 |
| B `floating` | `tag == 3` | `frames*channels` 枚 float32 全被读；downmix 后 `frames` 枚 | 正常支 |
| C 非浮点 | `tag != 3`（**含 `tag == 0` 这种"没解析出来"的值**） | 按 `frames*channels` 枚 **int16** 读＝只吞掉缓冲区**前一半字节** | 形③的入口 |

`MonoDownmix(in []int16, chans int)`（`internal/audio/resample.go:100-116`，逐字见下）**关键在它的长度算法把 C 支的"少读一半"藏掉了**：

```go
// MonoDownmix averages interleaved multi-channel samples to mono. chans <= 1
// returns the input unchanged.
func MonoDownmix(in []int16, chans int) []int16 {
	if chans <= 1 || len(in) == 0 {
		return in
	}
	n := len(in) / chans
	out := make([]int16, n)
	for i := 0; i < n; i++ {
		var acc int32
		for c := 0; c < chans; c++ {
			acc += int32(in[i*chans+c])
		}
		out[i] = clampI16(acc / int32(chans))
	}
	return out
}
```

⇒ **算的（本件最重要的一条推论）**：真机 48k／立体声／float32 那一形若走进 C 支，`unsafe.Slice((*int16)(data), frames*channels)` 取到的其实是**前 `frames` 枚 float32 的字节**（枚枚 float32 = 2 枚 int16），downmix 时 `chans = 2` ⇒ `n = (frames*2)/2 = frames` ⇒ **输出枚数与正确形状一模一样**。⇒ 后果：**任何"数帧数／数样本数／查长度"的既有断言都对这一支不敏感**，只有**幅度**会变形。这正是票面 `AC#0` 要打印"前 16 个 `int16` 逐枚"的理由；也 ⇒ **本票判形③不许用计数类判据**（要引用这条的落地腿请注意它是**算的**，我⛔ 没跑任何测试）。

---

## 4. ★ 码层浮出一枚**票面没点的更上游嫌疑**：`parseWaveFormat` 的 SubFormat 偏移与同包另一枚解析器**对不上**

`floating` 的唯一来源是 `s.format.tag`，而 `tag` 由 `parseWaveFormat` 定。它自己**零仪器**（尺 R10＝`git grep -n "parseWaveFormat" e96da9f4 -- '*.go'`，现量＝**3 枚命中全在 `wasapi_windows.go`**：注释 `:182`、定义 `:185`、调用 `:263`；⛔ 零 `_test.go`）。全文（逐字，整块，工作树读，`internal/audio/wasapi_windows.go:182-203`）：

```go
// parseWaveFormat reads WAVEFORMATEX (+ WAVEFORMATEXTENSIBLE SubFormat) from
// a CoTaskMem-allocated format pointer. Layout: tag@0, channels@2, rate@4,
// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].
func parseWaveFormat(p unsafe.Pointer) (waveFormat, error) {
	if p == nil {
		return waveFormat{}, observe.New(observe.ClassAudioDevice, "GetMixFormat returned no format")
	}
	f := waveFormat{
		tag:      *(*uint16)(p),
		channels: *(*uint16)(unsafe.Add(p, 2)),
		rate:     *(*uint32)(unsafe.Add(p, 4)),
		bits:     *(*uint16)(unsafe.Add(p, 14)),
	}
	if f.tag == waveFormatExt {
		sub := *(*windows.GUID)(unsafe.Add(p, 26))
		f.tag = uint16(sub.Data1) // KSDATAFORMAT_SUBTYPE_* low word: 1=PCM 3=FLOAT
	}
	if f.rate == 0 || f.channels == 0 {
		return f, observe.New(observe.ClassAudioDevice, "degenerate mix format")
	}
	return f, nil
}
```

同包**另一枚** WAVEFORMATEXTENSIBLE 解析器（WAV 文件侧，尺 R11＝`git grep -n "waveFormatExt\|SubFormat\|0xFFFE\|waveFormatFloat\|waveFormatPCM" e96da9f4 -- '*.go'`，现量 9 枚命中里唯一非 wasapi 的一枚在 `wavinjector.go:192`）逐字（整块，`internal/audio/wavinjector.go:192-197`）：

```go
				if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]
					if size < 40 {
						return nil, 0, 0, fmt.Errorf("extensible fmt chunk too short (%d)", size)
					}
					fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])
				}
```

⇒ **同一枚结构体、同一个包、两枚解析器、偏移差 2 字节**：WAV 侧读 `+24..+26` 并带 `size < 40` 卫（40 正是 `sizeof(WAVEFORMATEXTENSIBLE)`）；WASAPI 侧读 `+26..+28`。
⇒ **算的（码层，⛔ 不是实测）**：若真机 `GetMixFormat` 回的是 **WAVEFORMATEXTENSIBLE**（`tag == 0xFFFE`），从 `+26` 取到的 16 位是 `SubFormat.Data1` 这个 DWORD 的**高半字**——`KSDATAFORMAT_SUBTYPE_PCM`(0x1) 与 `..._IEEE_FLOAT`(0x3) 的高半字**都是 0** ⇒ `f.tag` 变成 **0** ⇒ `floating == false` ⇒ 走进 §3 的 **C 支**，而 `parseWaveFormat` 的退化卫只查 `rate == 0 || channels == 0`，**tag==0 一路放行**（`bits` 又从不被下游读）。这一支的读数形状就是票面记录的"近常数、max/mean≈1.08、sound 与 quiet 分不开"。
⇒ **如果真机回的是普通 WAVEFORMATEX（`tag == 3`、`cbSize == 0`）**，`+26` 那一条根本不走，`floating == true` 走 B 支 ⇒ **形③当场排除**，0.3677 只能归形②。
⇒ **码层到此为止能答的：两支哪一枚在跑，取决于真机 mix format 的 `wFormatTag` 是不是 0xFFFE。这一问本件答不了（具名＝码层答不了，需 AC#0 的 ⓐ 那一行四元组）。** 而这恰好就是票面 `AC#0` 已经点名的打印项，⛔ 我新建需求。
⚠ **我不在本文定罪**。上面是"算的"＋一枚可复核的偏移之争。要把"26 是错的"钉死，成本**极低且不碰真设备**：`parseWaveFormat` 吃的是裸 `unsafe.Pointer`，包内 windows-tag 用例可以手搓一枚 40 字节 WAVEFORMATEXTENSIBLE 喂进去，断言 `tag == 3`。⇒ 这条我给的是**形状**，不是结论；MSVC 真实布局偏移的最终权威在 SDK 头文件（`audioclient.h`/`mmreg.h`），**本仓没有那两枚头文件**（尺＝`find . -name 'audioclient.h' -o -name 'mmreg.h'` 我没跑 ⇒ 这一条**未证成**）。

---

## 5. `silent` / `AUDCLNT_BUFFERFLAGS_SILENT` 那一支在什么条件下才触发

分三问，前两问码层能答、第三问码层答不了：

1. **代码里的触发条件（逐字，整块，`internal/audio/wasapi_windows.go:101-104`）**：
```go
	audclntShareModeShared    = 0
	audclntStreamFlagsEventCB = 0x00040000
	audclntBufferFlagsSilent  = 0x2
	audclntBufferHns          = 1_000_000 // 100ms in 100ns units
```
   判据＝`IAudioCaptureClient::GetBuffer` 回参 `Flags` 带 bit `0x2`（`*(*uint32)` 由 `comCall(s.capture, 3, …)` 的第 3 个出参写出）。另一条进入 A 支的路与 flag 无关：`data == nil`。
2. **A 支与 0.3677 互斥（算的，凭两条现量）**：A 支返回 `make([]int16, frames)` ＝**恒零**；电平尺 `LevelOfSamples` 对恒零回 `0`（`internal/audio/level.go:96-107`，票面已现量并写"全零帧读 0"有常驻判据）。⇒ 若真机那六窗走的是 A 支，读数会是 `0.000000` 而不是 `0.3677`。⇒ **现象不来自 silent 支；silent 支这一格本票可以当场排除，不必再花仪器。**（⚠ 射程声明：这条依赖"电平尺零形有钉"那枚既有判据，我⛔ 没跑它，引的是票面 `:14` 的现量记录＋`level.go` 全文直读。）
3. **什么真实形状会让采集侧回 SILENT（码层答不了）**：仓内没有任何注释、测试或文档写过"我们的流会在 X 时收到 0x2"；`wasapiStream` 用的 stream flags 只有 `audclntStreamFlagsEventCB = 0x00040000`（逐字见 §2 的 `init` 调用 `comCall(s.client, slotIAudioClientInitialize, audclntShareModeShared, audclntStreamFlagsEventCB, …)`），**没有** `AUDCLNT_STREAMFLAGS_SILENT`。至于 WASAPI 在**采集**方向何时置 `AUDCLNT_BUFFERFLAGS_SILENT`，权威在微软文档，⛔ 本仓没有抄进来（尺＝`grep -rn "AUDCLNT_BUFFERFLAGS" --include='*.md' docs` 我没跑）。⇒ 具名：**码层答不了，且本票不需要它**（§5.2 已把 A 支排除在现象之外）。

---

## 6. 本件欠的尺（具名，落地腿⛔ 不许引为"已证"）

- ⛔ 零 `go build`／`go vet`／`go test`／`go test -list` ⇒ 我**没有**任何"某用例今天跑不跑、跑成什么色"的实测；`convertPacket` 的可测性是我**读签名**读的，不是编出来跑的。
- 没跑 `go env`／`go list`（自报见 00 件 §4 末段＋回报）。
- 没跑 `find` 找 SDK 头文件、没跑 `grep -c bits` 的函数体级尺、没跑 d22scan。
- 真机 mix format 四元组**未知** ⇒ §4 那条偏移之争的"今天到底走哪一支"悬着，手柄在票面 `AC#0` 的 ⓐ。
