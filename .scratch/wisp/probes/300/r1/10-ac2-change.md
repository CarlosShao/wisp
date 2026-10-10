# 300-r1 — `10` `AC#2` 落地那一发（ⓐ 偏移 ＋ ⓑ 同批注释改写）

锚＝`ec5cdc51`（起手三读见 `00-anchor.md`；我的第一笔＝`ce06cbe6`，只加本程锚件）。
本件射程＝`internal/audio/wasapi_windows.go` 一枚文件的两处：`parseWaveFormat` 的头顶注释块 ＋ 体内那枚偏移数字。
★**ⓐ 与 ⓑ 同一笔 commit**（名册尺见 `30-gates.md` §5，逐字 `git show --stat`）。

## 1. ⓐ 偏移：`26` ⇒ `24`（改动已落地，`landed-as` 尺逐字）

```
$ grep -n "unsafe.Add(p, " internal/audio/wasapi_windows.go
206:		channels: *(*uint16)(unsafe.Add(p, 2)),
207:		rate:     *(*uint32)(unsafe.Add(p, 4)),
208:		bits:     *(*uint16)(unsafe.Add(p, 14)),
211:		sub := *(*windows.GUID)(unsafe.Add(p, 24))
```

- **改的就是第 4 枚读数**（`26` ⇒ `24`）。前三枚（`@2`／`@4`／`@14`）**一字没动**，
  加上 `tag` 那次读取 `*(*uint16)(p)`（`:205`，不在上面这把尺的射程里，因为它不带 `unsafe.Add`）——
  派单写死这几枚**本来就对**，⛔ 动。⛔ 动 `convertPacket`。
- ⚠ **行号从 `:196` 移到 `:211`**（＋15 行＝ⓑ 那段注释改写的代价）。前腿件里逐字引的 `196:` 是本程之前的快照，
  引用者请按**内容**认那一行（尺＝`grep -n "sub := \*(\*windows.GUID)"`）。
  本程⛔ 碰 `cmd/wisp`，票 255 那份按行号引产码的名册**行数没变**（越界尺见 `30-gates.md` §6）。

## 2. ⓑ 同批注释改写（★这比代码更要紧：原文逐字教下一个改这里的人相信 26）

**改前**（`wasapi_windows.go:182-184`，逐字三行）：

```go
// parseWaveFormat reads WAVEFORMATEX (+ WAVEFORMATEXTENSIBLE SubFormat) from
// a CoTaskMem-allocated format pointer. Layout: tag@0, channels@2, rate@4,
// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].
```

**改后**（逐字整块，共 17 行，含一个空注释行分段）：

```go
// parseWaveFormat reads WAVEFORMATEX (+ WAVEFORMATEXTENSIBLE SubFormat) from
// a CoTaskMem-allocated format pointer. The layout is the WAVEFORMATEXTENSIBLE
// declared in the Windows SDK header shared/mmreg.h, which is included under
// pshpack1.h (packing = 1, so nothing is padded): WAVEFORMATEX is 18 bytes
// with tag@0, channels@2, rate@4, avgBytes@8, blockAlign@12, bits@14,
// cbSize@16, and the extension -- documented in that header with cbSize = 22
// -- adds Samples@18, dwChannelMask@20 and SubFormat@24, 40 bytes in total.
// Samples@18 is a two-byte union of three WORDs (wValidBitsPerSample,
// wSamplesPerBlock, wReserved); "validBits" is not a field name.
//
// SubFormat is the member that carries the real subtype behind a 0xFFFE
// wFormatTag, so Data1 is read at byte 24. KSDATAFORMAT_SUBTYPE_PCM is 1 and
// KSDATAFORMAT_SUBTYPE_IEEE_FLOAT is 3, both with a zero high half word, so a
// read two bytes late resolves every extensible mix format to 0; and because
// convertPacket branches on floating alone, that would re-slice 32-bit float
// frames as int16 while returning the very same frame count. The sibling
// wavinjector.go parses the same layout at body+24 under the same 40-byte
// total: mmreg.h is the authority for both, not that file.
```

### 2.1 那六个数的逐枚处置（★按**编排者那把**写，⛔ 按"三个数全坏"那把）

| 原注释里写着的 | 判 | 依据 |
|---|---|---|
| `tag@0` | ✅ 对 | `WAVEFORMATEX.wFormatTag` 是首字段（WORD）；产码 `*(*uint16)(p)` 同形 |
| `channels@2` | ✅ 对 | `nChannels` @2 |
| `rate@4` | ✅ 对 | `nSamplesPerSec` @4 |
| `bits@14` | ✅ 对 | `wBitsPerSample` @14 |
| `cbSize@16` | ✅ 对 | `cbSize` @16 |
| `validBits@18` | ⚠ **偏移对、名字不是字段名** | 真身是名为 `Samples` 的 union（成员 `wValidBitsPerSample`／`wSamplesPerBlock`／`wReserved` 三枚 WORD）；改后注释直接写 `Samples@18` 并把这层说破 |
| `channelMask@22` | ❌ **坏，应为 20** | `dwChannelMask` 是 union **外面**的 DWORD，@20 |
| `SubFormat@26` | ❌ **坏，应为 24** | 见 §3 |

⇒ 坏的只有两枚数（`22`⇒`20`、`26`⇒`24`）。一枚前腿那句"注释里三个数全是坏的"我这把⛔ 采纳（`18` 那个位置本身是对的，
错的只是把一个 union 说成一枚字段），而编排者在派单里同判（"只有两枚数是坏的"）。这一处**记成我与 `300-a1r` 的分歧**，⛔ 我裁谁错，交 `300-v2`。

### 2.2 改写里我刻意没写的东西

- ⛔ 现在时计数（"今天有 N 枚绿测试"）——派单写死，这类句子会被下一位改动打成假话。
- ⛔ 把 `wavinjector.go` 当权威：那句话只报"同形"，并**当场补了一句 `mmreg.h is the authority for both, not that file`**，
  因为票面 `AC#0` 的硬门就是⛔ 拿本仓另一枚实现当裁判。
- 出处只写**相对名**（`shared/mmreg.h`、`pshpack1.h`、`cbSize = 22`）⛔ 写绝对盘路径 ⇒ 注释里⛔ 机器专属坐标；
  绝对路径与逐字原文在本件 §3（那是证据，⛔ 产码注释）。

## 3. 权威出处（⛔ 我重新裁决"该不该修"，三把尺已独立量过；这一段只是**我落笔处引的东西我自己开文件核过**）

文件＝`C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`。

```
$ sed -n '2521,2534p' "/c/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/mmreg.h"
#ifndef _WAVEFORMATEXTENSIBLE_
#define _WAVEFORMATEXTENSIBLE_
typedef struct {
    WAVEFORMATEX    Format;
    union {
        WORD wValidBitsPerSample;       /* bits of precision  */
        WORD wSamplesPerBlock;          /* valid if wBitsPerSample==0 */
        WORD wReserved;                 /* If neither applies, set to zero. */
    } Samples;
    DWORD           dwChannelMask;      /* which channels are */
                                        /* present in stream  */
    GUID            SubFormat;
} WAVEFORMATEXTENSIBLE, *PWAVEFORMATEXTENSIBLE;
#endif // !_WAVEFORMATEXTENSIBLE_

$ grep -n "cbSize = 22" "/c/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/mmreg.h"
2540:typedef WAVEFORMATEXTENSIBLE    WAVEFORMATPCMEX; /* Format.cbSize = 22 */
2550:typedef WAVEFORMATEXTENSIBLE          WAVEFORMATIEEEFLOATEX; /* Format.cbSize = 22 */
```

⇒ `Format` 18 ＋ `Samples` 2（@18）＋ `dwChannelMask` 4（@20）＋ `SubFormat` 16（@24）＝ **40**，`Data1` 起始 **24**。
`pshpack1.h` 那句（packing＝1）由 `300-a1r` 与编排者各自读过（`mmreg.h:34` 一段），我这把⛔ 重开一遍：
我读到的结构体里字段之间**没有任何补空**（若 `GUID` 前按 8 字节对齐补空，`dwChannelMask`@20 之后会跳 4 ⇒ `SubFormat` 落 24 仍成立、总长 44 才变，
而 `cbSize = 22` 那两行**逐字写死了 22＝2＋4＋16**，与 packing=1 自洽），这一句是我**读**到的还是**推**的，标在 `90-unrun-rulers.md`。

## 4. 同包那枚兄弟读数（⛔ 裁判，只是"同形"这一事实我自己 grep 到了）

```
$ grep -n "body+24\|size < 40\|WAVE_FORMAT_EXTENSIBLE" internal/audio/wavinjector.go
163:// (format 1) and IEEE float32 (format 3, incl. WAVE_FORMAT_EXTENSIBLE with a
192:			if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]
193:				if size < 40 {
196:				fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])
```

⚠ 这一发**只用来支撑 §2 末尾那句注释里的"同形"**，⛔ 用来定偏移（那是 `AC#0` 的硬门）。

## 5. 本件没做的（具名，⛔ 悄悄不做）

- ⛔ 补第 5 形（`300-v1` §4.2 提的"非 extensible 面（plain PCM）没有用例钉着"）。
  理由＝入库那四形是与编排者五枚偏移扫描尺（`probes/300/orch/…-five-offset-scan.txt`）**逐形可对照**的那一套，
  我加一枚就断了对照；归口见 `90-unrun-rulers.md`。
- ⛔ 动 `convertPacket` 的分派形状（`floating` 单条件）⛔ 把 `bits` 拉进判据（票面 `:14` 现读过它体内 `bits` 出现 0 次）。
- ⛔ 动 `level.go` 的 `SineLevelTolerance`、`internal/observe/thresholds.go`、`internal/ball/liquid.go`（C21 冻结 token 面）。
