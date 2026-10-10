# 300-a1r — `10` `AC#0` 权威偏移：`SubFormat` 在 **24**（盘上 SDK 头文件逐字，＋ Learn 交叉核）

锚＝`b2933dc9`（起手闸门读数见 `00-anchor.md`）。本格**零产码改动**、⛔ 动任何跟踪文件。

## 结论（先给答案，再给凭据）

| 问 | 答 | 凭据种类 |
|---|---|---|
| ① `SubFormat` 相对结构起始的偏移 | **24**（⛔ 不是 26） | **盘上头文件** `mmreg.h`（`#include "pshpack1.h"` ⇒ 1 字节 Packing，逐字节可算）＋ Learn 字段序交叉核 |
| ② `cbSize` 文档值／结构总长 | `cbSize = 22`（Learn 写 "**at least 22**"，头文件两处 `/* Format.cbSize = 22 */` 写死 22）；结构总长 **40** ＝ 18 ＋ 22 | 同上 |
| ③ `wValidSamples` 与 `dwChannelMask` 是不是同一个 union | **不是**——两者都是 `WAVEFORMATEXTENSIBLE` 的**独立成员**；`Samples` 那枚 union 里只有**三枚 WORD**（`wValidBitsPerSample` / `wSamplesPerBlock` / `wReserved`），`dwChannelMask` 是它外面的一枚 **DWORD**；标准结构里**根本没有**叫 `wValidSamples`/`nValidSamples` 的字段 | 同上（头文件 typedef ＋ Learn Syntax 块两把一致） |

⇒ **`internal/audio/wasapi_windows.go` 的 `parseWaveFormat` 是错的那一枚；同包 `wavinjector.go`（读 24、守卫 `size < 40`）与权威布局同形。**
★派单的硬门是"**如果权威结论是 26 才对就停手上报**"——权威结论是 **24** ⇒ **不停手**，`AC#1` 的夹具按 **24** 造（见 `20-ac1-rig.md`）。

## 凭据 A（第一选择＝盘上那份，⛔ 拿本仓实现当裁判）

### A.1 文件定位

- 尺＝`find "/c/Program Files (x86)/Windows Kits" "/c/Program Files/Windows Kits" -maxdepth 6 -name 'mmreg.h'`
  ⇒ 命中枚数＝**1**（尺名＝`find` 的 stdout 行数；射程＝两棵 `Program Files*` 下的 `Windows Kits`，`-maxdepth 6`）：

```
/c/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/mmreg.h
```

- ⚠ 派单示例给的是 `Include/*/um/mmreg.h`，**盘上那枚在 `shared/` 不在 `um/`**（尺＝`ls ".../Include/"*/um/mmreg.h` ⇒ rc=2、0 行）。
  这是派单示例路径的过期处，⛔ 影响结论（同一枚 SDK，`shared/mmreg.h` 才是它的真位置）。
- 文件总行数＝**3814**（尺＝`wc -l`）。原始抽取全部逐字落 `logs/ac0-mmreg-extract.txt`（78 行，尺＝`wc -l`）。

### A.2 Packing 是这格的关键（先钉它，再算偏移）

`mmreg.h` 第 31–36 行逐字：

```c
#if !defined( RC_INVOKED ) && defined( _MSC_VER )
#if (_MSC_VER <= 800)
#pragma pack(1)
#else
#include "pshpack1.h"   /* Assume byte packing throughout */
#endif
#endif  /* RC_INVOKED */
```

- 尺＝`grep -n "pragma pack" mmreg.h` ⇒ **4 枚命中**（`grep -n` 数行）：`33`、`3767`、`3782`、`3809`。
  ⇒ 第 33/34 行打开 **1 字节 Packing**，一直有效到 3782 的 `#pragma pack(push, 1)`／3809 的 `pop`；
  我抽的两枚 typedef（2425–2442、2521–2534）与两处 `cbSize = 22` 注释（2540、2550）**全在 Packing=1 的射程里**。
- ★这一条把"编译器会不会插补空"这个问题**当场关掉**：Packing=1 ⇒ 字段**紧挨着排**，偏移就是字段宽度累加，
  没有对齐间隙 ⇒ 下面那张偏移表是**算出来的**（尺名＝算术，输入是头文件里逐字的字段宽度），不是猜的。

### A.3 `WAVEFORMATEX` typedef（`mmreg.h` 第 **2425–2442** 行，逐字整块）

```c
#ifndef _WAVEFORMATEX_
#define _WAVEFORMATEX_
typedef struct tWAVEFORMATEX
{
    WORD    wFormatTag;        /* format type */
    WORD    nChannels;         /* number of channels (i.e. mono, stereo...) */
    DWORD   nSamplesPerSec;    /* sample rate */
    DWORD   nAvgBytesPerSec;   /* for buffer estimation */
    WORD    nBlockAlign;       /* block size of data */
    WORD    wBitsPerSample;    /* Number of bits per sample of mono data */
    WORD    cbSize;            /* The count in bytes of the size of
                                    extra information (after cbSize) */

} WAVEFORMATEX;
typedef WAVEFORMATEX       *PWAVEFORMATEX;
typedef WAVEFORMATEX NEAR *NPWAVEFORMATEX;
typedef WAVEFORMATEX FAR  *LPWAVEFORMATEX;
#endif /* _WAVEFORMATEX_ */
```

- `cbSize` 的**语义**由它自己的注释给出（逐字，第 2435–2436 行）：
  `/* The count in bytes of the size of / extra information (after cbSize) */`
  ⇒ `cbSize` 数的是 **`cbSize` 之后**的字节，⛔ 含 `cbSize` 自己那 2 字节。
- 由字段宽度累加（WORD=2、DWORD=4，Packing=1）：`wFormatTag@0`、`nChannels@2`、`nSamplesPerSec@4`、
  `nAvgBytesPerSec@8`、`nBlockAlign@12`、`wBitsPerSample@14`、`cbSize@16` ⇒ **`sizeof(WAVEFORMATEX) = 18`**。
  ⚠ 顺带一句：`parseWaveFormat` 现读的 `tag@0 / channels@2 / rate@4 / bits@14` **这四位与权威一致**（错位只在 `SubFormat`）。

### A.4 `WAVEFORMATEXTENSIBLE` typedef（`mmreg.h` 第 **2521–2534** 行，逐字整块）

```c
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
```

头顶注释（第 2512–2520 行）逐字包含：`//  WAVEFORMATEXTENSIBLE structure. WAVEFORMATEXTENSIBLE allows you to`、
`//  define a new GUID value for the WAVEFORMATEXTENSIBLE.SubFormat field`、
`//  WAVEFORMATEXTENSIBLE.Format.wFormatTag field.` ⇒ 头文件自己就把 `SubFormat` 与 `Format.wFormatTag` 点名为这枚结构的成员。

### A.5 `cbSize = 22` 由头文件自己写死（两枚，逐字）

- 第 **2540** 行：`typedef WAVEFORMATEXTENSIBLE    WAVEFORMATPCMEX; /* Format.cbSize = 22 */`
- 第 **2550** 行：`typedef WAVEFORMATEXTENSIBLE          WAVEFORMATIEEEFLOATEX; /* Format.cbSize = 22 */`
- 尺＝`grep -n "cbSize = 22" mmreg.h` ⇒ **2 枚命中**（`grep -n` 数行）。
  ★**PCM 与 IEEE_FLOAT 两档都写 22** ⇒ `cbSize` 与 `bits`／与 subtype **无关**，票面推的那枚 22 **被盘上头文件证成**。

### A.6 偏移表（算出来的，输入＝A.3/A.4 逐字字段宽度；Packing=1）

| 字节偏移 | 成员 | 宽度来源 |
|---|---|---|
| 0–1 | `Format.wFormatTag` (WORD) | A.3 |
| 2–3 | `Format.nChannels` (WORD) | A.3 |
| 4–7 | `Format.nSamplesPerSec` (DWORD) | A.3 |
| 8–11 | `Format.nAvgBytesPerSec` (DWORD) | A.3 |
| 12–13 | `Format.nBlockAlign` (WORD) | A.3 |
| 14–15 | `Format.wBitsPerSample` (WORD) | A.3 |
| 16–17 | `Format.cbSize` (WORD) | A.3 |
| **18–19** | `Samples`（union，三成员皆 WORD ⇒ **2 字节**） | A.4 |
| **20–23** | `dwChannelMask` (DWORD) | A.4 |
| **24–39** | `SubFormat` (GUID = `Data1` DWORD 4 ＋ `Data2` WORD 2 ＋ `Data3` WORD 2 ＋ `Data4` 8 字节) | A.4 ＋ `x/sys/windows.GUID` 同形 |

- `SubFormat` 起始＝**24**。结构总长＝24＋16＝**40**。
- ★与 `cbSize` **互相验算**（这一步让两枚独立凭据咬合）：`cbSize` 数的是 18 之后的全部附加字节，
  按上表 18 之后＝ `Samples`(2) ＋ `dwChannelMask`(4) ＋ `SubFormat`(16) ＝ **22** ＝ A.5 那两枚注释写的 22；
  总长 18 ＋ 22 ＝ **40** ＝ `wavinjector.go` 那枚 `if size < 40` 守卫的 40。
  ⇒ **三枚各自独立的东西（字段序、`cbSize` 语义值、总长 40）在 24 这一形上闭合**；
  换成 26 那一形则 `cbSize` 得是 24、总长得是 42，与头文件那两枚 `/* Format.cbSize = 22 */` **直接矛盾**。

### A.7 夹具要用的三枚常量（逐字，同为盘上凭据）

- `mmreg.h` 第 **2376** 行：`#define  WAVE_FORMAT_EXTENSIBLE                 0xFFFE /* Microsoft */`
- 第 **2110** 行：`#define  WAVE_FORMAT_IEEE_FLOAT                 0x0003 /* Microsoft Corporation */`
- 第 **2418** 行：`#define WAVE_FORMAT_PCM         1`
- `SubFormat` GUID 本体（第 **2474**／**2483** 行，逐字）：
  `DEFINE_GUIDSTRUCT("00000001-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_PCM);`
  `DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);`
  ⇒ `Data1` 低字分别是 **1**（PCM）与 **3**（FLOAT），与派单要的取值一致；
  `ksmedia.h` 第 846／854 行是同一对 GUID 的第二次定义（尺＝`grep -n` 抽出，⛔ 另算凭据，只是同值重现）。

## 凭据 B（交叉核＝Microsoft Learn，取回日期 **2026-10-10**）

- URL：`https://learn.microsoft.com/en-us/windows/win32/api/mmreg/ns-mmreg-waveformatextensible`
  （页面自报 `canonicalUrl` 同址、`ms.date: 2023-04-26`、`req.header: mmreg.h`）
- Syntax 块逐字（与 A.4 字段序**一字不差**）：

```cpp
typedef struct {
  WAVEFORMATEX Format;
  union {
    WORD wValidBitsPerSample;
    WORD wSamplesPerBlock;
    WORD wReserved;
  } Samples;
  DWORD        dwChannelMask;
  GUID         SubFormat;
} WAVEFORMATEXTENSIBLE, *PWAVEFORMATEXTENSIBLE;
```

- `cbSize` 那句逐字：`The **cbSize** member must be at least 22.`（在 `Format` 成员条目下，同段还有
  `The **wFormatTag** member must be WAVE\_FORMAT\_EXTENSIBLE.`）
- `dwChannelMask` 成员条目逐字：`Bitmask specifying the assignment of channels in the stream to speaker positions.`
  ⇒ Learn 也把 `Samples` 与 `dwChannelMask` 列成**两个成员**，⛔ 任何一处把 `dwChannelMask` 说成 `Samples` union 的成员。
- Remarks 逐字一句（★这解释了为什么"多出来的字节"就是那 22 枚）：
  `**WAVEFORMATEXTENSIBLE** can safely be cast to **WAVEFORMATEX**, because it simply configures the extra bytes specified by **WAVEFORMATEX.cbSize**.`
- ⛔ **`wValidSamples`／`nValidSamples` 在两把凭据里都找不到**：
  尺＝`grep -n "ValidSamples"` 射程＝`shared/mmreg.h` `shared/ksmedia.h` `um/mmdeviceapi.h` `um/audioclient.h`
  ⇒ **stdout 0 行**（⚠ 这条命令的 `grep rc` 因多文件＋`2>&1` 混在一起，**只认 stdout 行数这把尺**；见 `logs/ac0-mmreg-extract.txt` 末段）。
  ⇒ `wValidSamples` 那枚名字的出处**不在我这两把权威里**，本格⛔ 对它做任何断言；
  我只回答被问的那句：**它不是与 `dwChannelMask` 同一个 union 的成员**（两把权威都把 `dwChannelMask` 列为独立 DWORD）。

## 两枚被驳的假设（⛔ 和稀泥，各自具名）

1. **`parseWaveFormat` 头顶那三行布局注释**（`internal/audio/wasapi_windows.go`，内容锚＝逐字
   `// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].`）——**三处与权威不符**：
   ① 没有叫 `validBits` 的独立字段（是 `Samples` 那枚 **2 字节** union 里的 `wValidBitsPerSample`）；
   ② `channelMask` 在 **20** 不在 22；③ `SubFormat` 在 **24** 不在 26。
   ⇒ 票面 `AC#2` 说"注释现在逐字教读者相信 `SubFormat@26`，这比代码更危险"——**本格证实这句成立**，
   注释里的三个数**全都是**坏的（⛔ 只改 `26` 那一个数就留了两枚假话）。
2. **派单／票面 `:13` 那句"标准结构里没有独立的 validBits 字段，`nValidSamples` 是与 `dwChannelMask` 同一个 union 的成员"**——
   **前半对、后半错**：前半对（`validBits` 不是独立字段）；后半**与两把权威都矛盾**——
   `dwChannelMask` 是 `WAVEFORMATEXTENSIBLE` 的独立 DWORD 成员，`Samples` 是另一枚独立的 union 成员（内容只有三枚 WORD），
   **没有一枚 union 同时装着这两样**。⇒ ⚠ 这条**不改本票结论**（偏移 24 靠的是字段序＋`cbSize=22`＋总长 40，
   与这句无关），但它在票面现量节里，**归编排者处置**（票面纪律＝⛔ 我改票面其它一字，只准追加 `Progress log` 一行）。

## 本格没做的两件事（写清边界，⛔ 冒充交满）

- ⛔ **真机 `GetMixFormat` 那一发**（票面 `AC#3`〔仅本机可量，归编排者〕）——本格没跑，⛔ 也不许派给我。
  ⇒ "本机混音格式到底是不是 extensible＋float"仍**未定**。
- ⛔ **改偏移**（票面 `AC#2`）——本格零产码；票面 `:41` 逐字写着"`AC#0`＋`AC#1` 交回后**归我裁**，`AC#2` 要**非实现者**裁过才派"。
