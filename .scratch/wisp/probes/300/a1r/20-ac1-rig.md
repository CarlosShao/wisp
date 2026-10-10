# 300-a1r — `20` `AC#1` 决定性一发（仓外成对树，⛔ 恒红用例入库）

锚＝`b2933dc9`。夹具布局**不来自本仓任何一枚实现**，来自 `10-ac0-authority.md` 的权威结论（`SubFormat` ＠ **24**）。
本程**零产码改动**：跟踪文件一枚没动（越界尺见 `90-hygiene.md`）。

## 0. 导出树（怎么来的，逐字命令）

```
cd "D:\work\workspace\projects plans\Wisp" && mkdir -p /tmp/wisp300-a1r/tree \
  && git archive b2933dc9 | tar -x -C /tmp/wisp300-a1r/tree
```
- 尺＝管道右端 `tar` 的退码：`rc-pipeline(tar)=0`（⚠ **不是** `git archive` 的退码，那把尺在这条形下取不到；
  我用后续存在性读数代替它：`ls /tmp/wisp300-a1r/tree/internal/audio/wasapi_windows.go` ⇒ 存在，
  `ls $T/internal/audio/*.go | wc -l` ⇒ **19**（尺＝`wc -l`，射程＝导出树的 `internal/audio` 目录），
  `grep -n "unsafe.Add(p, 26)" $T/internal/audio/wasapi_windows.go` ⇒ **196:		sub := *(*windows.GUID)(unsafe.Add(p, 26))**
  ⇒ 锚上那枚错位读法**确实被完整导出**，用例测的就是它）。
- 树路径（MSYS／Windows 两种写法，留着给我核）：`/tmp/wisp300-a1r/tree`
  ＝ `C:\Users\swq\AppData\Local\Temp\wisp300-a1r\tree`（尺＝`cygpath -w`）。
- ★**这棵树留着⛔ 删**（编排者要核）；用例源码就躺在它的 `internal/audio/parse_wave_format_300_windows_test.go`（**6172 字节**，尺＝`wc -c`）。

### 0.1 一条派单示例没覆盖的坑（具名，⛔ 我悄悄绕过）

派单给的 harness 里 `$PWD/third_party/sherpa-onnx` 在**导出树里是空的**：
- 尺＝`git ls-files third_party/sherpa-onnx | wc -l` ⇒ **0**（tracked 枚数＝0）
- 同一条尺的盘上一面＝`ls third_party/sherpa-onnx/*.dll | wc -l` ⇒ **3**（⚠ 未跟踪）
⇒ `git archive` 只带 tracked 件，所以那 3 枚 DLL **进不了导出树**。
处置＝**从仓里 `cp` 进导出树同一路径**（`dll-copied=3`，⛔ 反向写仓、⛔ 改仓里任何东西），派单那条 PATH 便按原样成立。
⚠ 顺带量到：尺＝`grep -rl 'import "C"' <导出树>/internal | head -3` ⇒ **0 命中**（`internal/**` 整棵**无 cgo**），
所以 `internal/audio` 那一发**本来不靠 sherpa DLL 载入**；我仍**逐字带上了那条 PATH harness**（派单写死），
并交回"它到底跑成了没有"的判据＝下面两发都有 `--- FAIL`/`--- PASS` 行、**没有** `exit status 0xc0000135`（＝用例真跑了，⛔ "没跑成"伪装成红）。

## 1. 用例源码（整块，未截断；落点＝导出树，⛔ 仓）

`internal/audio/parse_wave_format_300_windows_test.go`：

```go
//go:build windows

package audio

// Ticket 300 AC#1 decisive rig (probe 300-a1r).
//
// This file exists ONLY inside the out-of-repo export tree produced by
//   git archive b2933dc9 | tar -x -C /tmp/wisp300-a1r/tree
// It is deliberately NOT committed: the assertions below are red at the
// anchor, and an always-red tracked test is banned by repo convention.
//
// Fixture byte layout comes from the AC#0 authority (on-disk Windows SDK
// header C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h,
// which is compiled with 1-byte packing, so there is no padding anywhere):
//
//	WAVEFORMATEX      : tag@0 channels@2 rate@4 avgBytes@8 blockAlign@12
//	                  bits@14 cbSize@16                       -> 18 bytes
//	extension (cbSize = 22, mmreg.h:2540 and :2550):
//	                  Samples@18 (union of three WORDs, 2 bytes)
//	                  dwChannelMask@20 (DWORD, 4 bytes)
//	                  SubFormat@24 (GUID, 16 bytes)
//	total 40 bytes; SubFormat.Data1 therefore starts at byte 24, NOT 26.
//
// Nothing here uses a device, a microphone, or a .wav file: parseWaveFormat
// takes an unsafe.Pointer, so the face is built by hand in a []byte.

import (
	"encoding/binary"
	"runtime"
	"testing"
	"unsafe"
)

// face300 describes one 40-byte WAVEFORMATEXTENSIBLE byte face.
type face300 struct {
	channels uint16
	rate     uint32
	bits     uint16
	data1    uint32 // SubFormat.Data1 value to write
	data1At  int    // byte offset the face writes Data1 at (24 = authoritative)
}

// build returns the face. The slice is created with make and returned, so it
// escapes to the heap by construction (a make result that flows out of a
// function cannot be stack-allocated); the caller pins it with
// runtime.KeepAlive across the parseWaveFormat call, and Go's GC does not move
// live objects, so the address handed to unsafe.Pointer stays valid.
func (f face300) build() []byte {
	b := make([]byte, 40)
	le := binary.LittleEndian
	le.PutUint16(b[0:], 0xFFFE)                        // WAVE_FORMAT_EXTENSIBLE (mmreg.h:2376)
	le.PutUint16(b[2:], f.channels)                    // nChannels
	le.PutUint32(b[4:], f.rate)                        // nSamplesPerSec
	le.PutUint32(b[8:], f.rate*uint32(f.channels)*(uint32(f.bits)/8)) // nAvgBytesPerSec
	le.PutUint16(b[12:], f.channels*(f.bits/8))        // nBlockAlign
	le.PutUint16(b[14:], f.bits)                       // wBitsPerSample
	le.PutUint16(b[16:], 22)                           // cbSize = 22 (mmreg.h:2540/2550)
	le.PutUint16(b[18:], f.bits)                       // Samples.wValidBitsPerSample
	le.PutUint32(b[20:], uint32((1<<f.channels)-1))     // dwChannelMask
	le.PutUint32(b[f.data1At:], f.data1)               // SubFormat.Data1
	return b
}

// TestParseWaveFormatSubFormatOffset300 asserts the value parseWaveFormat
// resolves out of an extensible mix format -- not that it was called.
//
// The four faces are chosen so that the assertion cannot be a tautology:
// S1/S2/S3 put SubFormat.Data1 at the authoritative offset 24 and demand the
// value that lives there; S4 is a deliberately malformed face that puts Data1
// at 26 (the offset the anchor's parseWaveFormat reads) and therefore demands
// 0, because at offset 24 of that face the authoritative layout has the tail of
// dwChannelMask plus zero bytes. A ruler that let both S1 and S4 through would
// be a ruler that distinguishes nothing.
func TestParseWaveFormatSubFormatOffset300(t *testing.T) {
	cases := []struct {
		name    string
		face    face300
		wantTag uint16
		why     string
	}{
		{
			name:    "S1_float_at_24",
			face:    face300{channels: 2, rate: 48000, bits: 32, data1: 3, data1At: 24},
			wantTag: 3,
			why:     "KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24",
		},
		{
			name:    "S2_pcm_at_24",
			face:    face300{channels: 2, rate: 48000, bits: 16, data1: 1, data1At: 24},
			wantTag: 1,
			why:     "KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001) at byte 24",
		},
		{
			name:    "S3_bogus_subtype_7_at_24",
			face:    face300{channels: 2, rate: 48000, bits: 32, data1: 7, data1At: 24},
			wantTag: 7,
			why:     "positive control: an obviously wrong Data1 must NOT be swallowed into PCM/FLOAT; the parser has to report the value it actually read, so this face may not go green on a constant",
		},
		{
			name:    "S4_data1_shifted_to_26",
			face:    face300{channels: 2, rate: 48000, bits: 32, data1: 3, data1At: 26},
			wantTag: 0,
			why:     "positive control for offset sensitivity: malformed face with Data1 two bytes late; per the authoritative layout the value at byte 24 is 0, so an assertion that passes here would read nothing meaningful",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			face := tc.face.build()
			if len(face) != 40 {
				t.Fatalf("face length = %d, want 40", len(face))
			}
			p := unsafe.Pointer(&face[0])
			got, err := parseWaveFormat(p)
			runtime.KeepAlive(face)

			if err != nil {
				t.Fatalf("parseWaveFormat error = %v, want nil; face: %s", err, tc.why)
			}
			if got.tag != tc.wantTag {
				t.Errorf("tag = %d, want %d (%s)", got.tag, tc.wantTag, tc.why)
			}
			// Header sanity: these three offsets are read correctly even at the
			// anchor, so a pass line here proves the face is a well-formed
			// WAVEFORMATEXTENSIBLE and that only the SubFormat position is in play.
			if got.channels != tc.face.channels || got.rate != tc.face.rate || got.bits != tc.face.bits {
				t.Errorf("header = {channels %d rate %d bits %d}, want {channels %d rate %d bits %d}",
					got.channels, got.rate, got.bits,
					tc.face.channels, tc.face.rate, tc.face.bits)
			}
			// Real-world consequence named in ticket 300: convertPacket picks the
			// sample path from floating alone, so a wrong tag silently reparses
			// 32-bit float frames as int16.
			gotFloat := got.tag == waveFormatFloat
			wantFloat := tc.wantTag == waveFormatFloat
			if gotFloat != wantFloat {
				t.Errorf("floating = %v, want %v (convertPacket would take the %s path)",
					gotFloat, wantFloat, map[bool]string{true: "float32", false: "int16"}[wantFloat])
			}
		})
	}
}
```

### 1.1 `[]byte` 怎么"钉住"（派单点名要我具名写怎么保证的）

三道，⛔ 只讲道理：

1. **堆分配是编译器自己说的**，不是我推的。尺＝在导出树里
   `go test -c -gcflags='-m' -o ../probe300-test.exe ./internal/audio/` 后
   `grep 'parse_wave_format_300' | grep -i 'escapes to heap'`（命中枚数＝**16**，尺＝`grep -c 'escapes to heap'`，射程＝该件），
   其中第 1 行逐字（原始输出已存 `logs/ac1-escape.txt`）：

```
internal\audio/parse_wave_format_300_windows_test.go:49:11: make([]byte, 40) escapes to heap
```

   ⇒ 第 49 行那枚 40 字节面**在堆上**（`build()` 把 `make` 的结果 return 出去 ⇒ 不可能在栈上）。
2. **活性**：`p := unsafe.Pointer(&face[0])` 之后立刻
   `runtime.KeepAlive(face)` 压在 `parseWaveFormat` 调用**之后**（见源码块）⇒ 从取址到解析返回，那块字节面全程是活的根。
3. **不移动**：Go 现行 GC **不搬移活对象** ⇒ 地址在整个调用期间有效。
   ⚠ 我把这句写成"Go 当前实现的性质"，⛔ 写成语言保证；它由 `go version go1.27.1 windows/amd64`（见 `00-anchor.md`）这一发实测成立。
- ★附带仪器读数：尺＝`go vet ./internal/audio/`（导出树内，PATH 同 harness）⇒ **`rc-vet=0`**、stdout 空
  ⇒ 这枚 `unsafe.Pointer` 用法过 `unsafeptr` 检查（与仓里 `wasapi_windows.go` 头顶那句
  `// unsafe discipline (vet-clean, same as internal/ball)` 同形）。原始输出＝`logs/ac1-vet.txt`。

## 2. 每一形与它的颜色配对（★两发成对读数，⛔ 只交一形）

跑法（两发**一字同形**，只换树上那枚偏移）：

```
cd /tmp/wisp300-a1r/tree && PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test ./internal/audio/ -count=1 -run 'TestParseWaveFormatSubFormatOffset300' -v
```
- `rc` 取法＝`${PIPESTATUS[0]}`（⛔ `$?` after a `| tee`，那把尺量的是 `tee`）。

| 形 | 字节面（`Data1` 位置／值） | 权威布局该读出 | 锚上实读 | 树上码 | 颜色（改前） | 颜色（突变后） |
|---|---|---|---|---|---|---|
| **S1** `float_at_24` | Data1=**3** ＠ **24**，bits 32 | 3 ⇒ `floating=true` | 0 | 26 | **FAIL** | **PASS** |
| **S2** `pcm_at_24` | Data1=**1** ＠ **24**，bits 16 | 1 ⇒ `floating=false` | 0 | 26 | **FAIL** | **PASS** |
| **S3** `bogus_subtype_7_at_24`（正控·值） | Data1=**7** ＠ **24** | 7（既非 PCM 亦非 FLOAT） | 0 | 26 | **FAIL** | **PASS** |
| **S4** `data1_shifted_to_26`（正控·位置） | Data1=**3** ＠ **26**（故意坏面） | **0**（＠24 处是 dwChannelMask 高字节＋零） | **3** | 26 | **FAIL** | **PASS** |

枚数尺口径（两发同一把尺）：`grep -c '^--- \(FAIL\|PASS\)'`（**顶层**）＝改前 **1** 枚顶层 FAIL、改后 **1** 枚顶层 PASS；
含子测试的顶层＋子行两把尺＝改前 **5** 行 FAIL（1 顶层＋4 子）、改后 **5** 行 PASS。★两把差 4 是常态（4 枚子测试），⛔ 拿一把去减另一把。

### 2.1 改前（树上码＝锚上的 `26`，逐字原始输出，`logs/ac1-pristine.txt`）

```
=== RUN   TestParseWaveFormatSubFormatOffset300
=== RUN   TestParseWaveFormatSubFormatOffset300/S1_float_at_24
    parse_wave_format_300_windows_test.go:122: tag = 0, want 3 (KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24)
    parse_wave_format_300_windows_test.go:138: floating = false, want true (convertPacket would take the float32 path)
=== RUN   TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24
    parse_wave_format_300_windows_test.go:122: tag = 0, want 1 (KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001) at byte 24)
=== RUN   TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24
    parse_wave_format_300_windows_test.go:122: tag = 0, want 7 (positive control: an obviously wrong Data1 must NOT be swallowed into PCM/FLOAT; the parser has to report the value it actually read, so this face may not go green on a constant)
=== RUN   TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26
    parse_wave_format_300_windows_test.go:122: tag = 3, want 0 (positive control for offset sensitivity: malformed face with Data1 two bytes late; per the authoritative layout the value at byte 24 is 0, so an assertion that passes here would read nothing meaningful)
    parse_wave_format_300_windows_test.go:138: floating = true, want false (convertPacket would take the int16 path)
--- FAIL: TestParseWaveFormatSubFormatOffset300 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S1_float_at_24 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26 (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/audio	0.023s
FAIL
```

`rc-gotest=1`（⛔ 这条 `rc` 从 `tee` 取）。**没有** `exit status 0xc0000135` ⇒ 用例真跑了。

★**票面 `AC#1` 那句"预期今天两枚都红（都读出 0）"被证实且被加强**：S1/S2/S3 都读出 **0**；
而 S4 读出 **3**——**锚上的码在坏面上恰好返回"那枚坏布局本该给的值"**。
⇒ 这枚断言**两个方向都敏感**：把 `Data1` 挪到 26 去喂它，它**⛔ 放行**（要 0 拿到 3，红）。
按本仓定式（先例 `A710`／"恒真的第三形是一把两种形状都放行的尺"），这就是那枚第三形：**它⛔ 是一枚尺**。

### 2.2 改后（同一棵树、同一份用例，只把 `unsafe.Add(p, 26)` 改成 `24`）

突变落地凭据（尺＝`grep -n "sub := \*(\*windows.GUID)"`，逐字）：

```
196:		sub := *(*windows.GUID)(unsafe.Add(p, 24))
```

原始输出（`logs/ac1-mutated-24.txt`）：

```
=== RUN   TestParseWaveFormatSubFormatOffset300
=== RUN   TestParseWaveFormatSubFormatOffset300/S1_float_at_24
=== RUN   TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24
=== RUN   TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24
=== RUN   TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26
--- PASS: TestParseWaveFormatSubFormatOffset300 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S1_float_at_24 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26 (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/audio	0.035s
```

`rc-gotest=0`。**这一发的作用**（⛔ 把它当"修好了"，那是 `AC#2`）：
它把"四枚全红"这一读数**限定**成"只差那枚偏移数字"——
如果红是我把夹具写坏了，突变后它**⛔ 可能**全绿（夹具与偏移无关地坏）。四枚成对翻色 ⇒
1. 夹具本身良构（且 `header`/`face length` 两把 sanity 检查在改前改后**都没有报错行**）；
2. 这枚尺**能绿** ⇒ 它⛔ 是恒红尺；
3. `AC#2` 把偏移改成 24 之后，这枚用例**一进门就该绿** ⇒ 票面要的"改前必红"凭据齐了。

### 2.3 树已还原（⛔ 把突变留在导出树里，那会打红编排者下一发核）

```
rc-restore=0
restored-identical=yes
196:		sub := *(*windows.GUID)(unsafe.Add(p, 26))
```
尺＝`cmp internal/audio/wasapi_windows.go ../wasapi_windows.go.orig`（无输出＝逐字节相同）＋还原后再读一次那一行。
备份件 `../wasapi_windows.go.orig` 留在 `/tmp/wisp300-a1r/` 里⛔ 删（同一棵待核树）。

## 3. 本格交回的判语（⛔ 越到 `AC#2` 的权）

- `AC#1` **成立**：锚上 `parseWaveFormat` 对**任何** extensible 面都把 `tag` 读成 `Data1@26`，
  即对**正确排布**的面恒读 **0**（S1/S2/S3 现证）⇒ `floating` 恒 `false` ⇒ 32 位浮点帧走 `int16` 路径（S1 那句
  `floating = false, want true (convertPacket would take the float32 path)` 就是票面 §「为什么要紧」说的形状）。
- ⇒ 票面 `AC#1` 那句"预期两枚都红"⛔ 是结论了，它是**读数**（改前一发逐字见 2.1）。
- ⛔ 我判 `AC#2` 怎么修（票面 `:41`：`AC#0`＋`AC#1` 交回归编排者裁，`AC#2` 要非实现者裁过才派）。
  我只具名一条**给 `AC#2` 的输入**：按 `10` 件 §「两枚被驳的假设」，头顶注释里 `validBits@18`／`channelMask@22`／`SubFormat@26`
  **三个数都是坏的**，⛔ 只改 `26` 一个数。
- ⛔ 本格没做的：真机 `GetMixFormat`（`AC#3`，仅本机可量，归编排者）、CI 可见性三跳（`AC#5`，落地腿自己现跑）。
  我这枚用例**今天进不进 CI 执行面**＝⛔ 我读数，⛔ 按门禁选落点（票面 `AC#5` 逐字禁这条）。
