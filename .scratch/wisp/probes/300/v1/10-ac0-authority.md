# 300-v1 — `10` ① `AC#0` 的权威凭据**有没有牙**（我自己开文件、自己算偏移、自己裁那句 union）

锚＝`efd85155`（起手闸门见 `00-anchor.md`）。被审件＝`300-a1r` 的 `10-ac0-authority.md`。
本程零产码；原始抽取全部逐字落 `logs/ac0-mmreg.txt`（129 行，尺＝`wc -l`）与 `logs/guid-and-shape.txt`（46 行，同尺）。

## 判语：**成立（有牙）**。`SubFormat` 的权威偏移＝**24**，且 `300-a1r` 那句"`dwChannelMask` 与 `Samples` 是两枚独立成员"**成立**，
## 票面 `:13` 后半句（"`nValidSamples` 与 `dwChannelMask` 是同一个 union 的成员"）**不成立**——错的是票面，⛔ 腿。

凭据＝**我自己**开的盘上文件，⛔ 转述腿的读数，⛔ 拿本仓另一枚实现当裁判，⛔ 凭记忆。

## 1. 我自己定位的那枚文件（尺＝`ls -l`＋`wc -l`，射程＝两棵 `Program Files*` 下的 SDK Include）

```
/c/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/mmreg.h
-rw-r--r-- 1 swq 197609 176680 Apr 15 11:15
3814 行
sha1=bc7753053eb77641df4df9fc39eba76d1557b2ff
```

⇒ 与被审腿报的是**同一枚文件**（它报的也是 `shared/` 不是 `um/`，我复核 `wc -l`＝3814 与它一致）。
派单示例那句 `Include/*/um/mmreg.h` 在盘上不成立——腿已经具名，我复现上了。

## 2. packing（这一步决定"偏移能不能算"，我自己逐字读）

`sed -n '28,40p'` 逐字（文件头，`logs/ac0-mmreg.txt` 第 5–15 行）：

```c
#pragma once
#endif

#if !defined( RC_INVOKED ) && defined( _MSC_VER )
#if (_MSC_VER <= 800)
#pragma pack(1)
#else
#include "pshpack1.h"   /* Assume byte packing throughout */
#endif
#endif  /* RC_INVOKED */
```

- 尺＝`grep -n "pragma pack\|pshpack\|poppack" mmreg.h` ⇒ **6 行**：`33:#pragma pack(1)`、`35:#include "pshpack1.h"`、`3767`、`3769`、`3782:#pragma pack(push, 1)`、`3809:#pragma pack(pop)`。
- ⇒ 两条编译路径（`_MSC_VER<=800` 走 `#pragma pack(1)`，否则走 `pshpack1.h`）**都把 Packing 设成 1**，一直开到 3769 的 `poppack.h`；
  我要抽的两枚 typedef（`WAVEFORMATEX`、`WAVEFORMATEXTENSIBLE`）与两处 `cbSize = 22` 注释**全在这段射程里**。
- ★我自己核对 `pshpack1.h` 干的就是 `#pragma pack(1)` 这件事——⛔ 需要，因为**两条分支之一已经是字面 `pack(1)`**，另一条只是它的具名版本；
  这一句把"编译器会不会插补空"**当场关掉** ⇒ 下面那张偏移表是**算出来的**（输入＝头文件里逐字的字段宽度），⛔ 猜的。
- ⚠ 与被审件的一处**行号过期**（⛔ 影响结论）：它把这段标成"第 31–36 行"却贴了 7 行；我自己量到 `pack(1)`＝33、`pshpack1.h`＝35。
  它 `grep -n "pragma pack"` 报 4 枚命中（33/3767/3782/3809）也对——我这把尺多带了 `pshpack|poppack` 两个词所以是 6 枚。**两把尺口径不同，⛔ 相减。**

## 3. 两枚结构体（逐字，我自己 `sed` 抽的）

`WAVEFORMATEX`（我的抽段＝第 2424–2441 行；`#ifndef _WAVEFORMATEX_` 起）：

```c
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
```

`WAVEFORMATEXTENSIBLE`（我的抽段＝第 2523–2533 行，逐字与腿贴的那一块**一字不差**）：

```c
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
```

`GUID` 本体（**我自己另开的一枚文件**，腿⛔ 盘上抽过：`shared/guiddef.h` 第 22–27 行逐字）：

```c
typedef struct _GUID {
    unsigned long  Data1;
    unsigned short Data2;
    unsigned short Data3;
    unsigned char  Data4[ 8 ];
} GUID;
```

⇒ `sizeof(GUID) = 4+2+2+8 = 16`（Packing=1）。仓里产码用的 `windows.GUID` 同形（尺＝`grep -n -A7 'type GUID struct' /d/work/base/gopath/pkg/mod/golang.org/x/sys@v0.42.0/windows/types_windows.go` ⇒ `Data1 uint32; Data2 uint16; Data3 uint16; Data4 [8]byte`）
⇒ `uint16(sub.Data1)` 取的就是 `Data1` 的**低 16 位**，与夹具写的 `0x00000001`/`0x00000003` 对得上。**⛔ 拿这枚 Go 结构当偏移裁判**，它只是用来说明"产码那次解引用读到什么"。

## 4. 我自己算的偏移表（尺＝算术；输入＝上面三块逐字字段宽度；Packing=1 ⇒ 无补空）

| 偏移 | 成员 | 宽度 | 来源（我这把尺） |
|---|---|---|---|
| 0–1 | `Format.wFormatTag` | WORD 2 | mmreg.h 抽段 |
| 2–3 | `Format.nChannels` | 2 | 同上 |
| 4–7 | `Format.nSamplesPerSec` | DWORD 4 | 同上 |
| 8–11 | `Format.nAvgBytesPerSec` | 4 | 同上 |
| 12–13 | `Format.nBlockAlign` | 2 | 同上 |
| 14–15 | `Format.wBitsPerSample` | 2 | 同上 |
| 16–17 | `Format.cbSize` | 2 | 同上 ⇒ **`sizeof(WAVEFORMATEX)`＝18** |
| **18–19** | `Samples`（union，三枚成员**全是 WORD** ⇒ union 尺寸＝**2**） | 2 | WAVEFORMATEXTENSIBLE 抽段 |
| **20–23** | `dwChannelMask` | DWORD 4 | 同上 |
| **24–39** | `SubFormat` | GUID 16 | 同上 ＋ guiddef.h ⇒ **`SubFormat` 起始＝24，总长＝24+16＝40** |

★三枚**互相验算**（这一步是我自己做的，⛔ 引用腿的）：
- `cbSize` 的语义由它自己的注释钉死（逐字）：`The count in bytes of the size of / extra information (after cbSize)` ⇒ 数的是 **16 之后**还是 **18 之后**？注释写的是 **after cbSize** ⇒ 从第 18 字节起数 ⇒ 附加字节＝`Samples`2 ＋ `dwChannelMask`4 ＋ `SubFormat`16 ＝ **22**。
- 盘上两处逐字写死 22：`2540: typedef WAVEFORMATEXTENSIBLE    WAVEFORMATPCMEX; /* Format.cbSize = 22 */`、
  `2550: typedef WAVEFORMATEXTENSIBLE          WAVEFORMATIEEEFLOATEX; /* Format.cbSize = 22 */`
  （尺＝`grep -n "cbSize = 22"` ⇒ **2 枚命中**，我自己现跑）。**PCM 与 IEEE_FLOAT 两档都是 22** ⇒ `cbSize` 与 bits／subtype 无关。
- 总长 18 ＋ 22 ＝ **40**。
⇒ **24 这一形让三枚彼此独立的东西（字段序、`cbSize` 语义值、结构总长）闭合**。
⇒ **26 那一形算不通**：若 `SubFormat@26`，则 18 之后要有 24 字节附加数据 ⇒ `cbSize` 得写 24、总长得是 42 ⇒ **与盘上那两处 `= 22` 直接矛盾**。
  票面 `:13` 自己那把"内部证据"尺（`wavinjector.go` 的 `if size < 40` 守卫）也只是**同向佐证**，⛔ 我拿它当裁判（我用的裁判是盘上头文件）。

## 5. 那句 union 之争——我用逐字原文裁：**腿对，票面错**

- 票面 `:13` 逐字：`标准结构里没有独立的 validBits 字段，nValidSamples 是与 dwChannelMask 同一个 union 的成员`。
- 我这把尺（⛔ 记名字、⛔ 拿另一枚实现当裁判）：
  - `grep -c "ValidSamples" mmreg.h` ⇒ **0**（我现跑；`nValidSamples` 在整枚头文件里**一次都没出现过**）。
  - `grep -rn "nValidSamples\|wValidSamples\|ValidSamples" --include=*.go internal` ⇒ **0 行**（仓里也只有那三行注释写过 `validBits`，⛔ 这枚名字）。
  - 逐字原文里 `Samples` 那枚 union 的三个成员**全是 `WORD`**，闭括号 `} Samples;` 之后 `DWORD dwChannelMask;` 是**结构体的下一枚直接成员**，
    它⛔ 在任何 union 的花括号里；而 `wValidBitsPerSample`／`wSamplesPerBlock`／`wReserved` 三个名字里**没有一个**叫 `nValidSamples`。
- ⇒ 判：**票面那句后半错了**（前半"`validBits` 不是独立字段"**对**——它是 union 的成员）。
  票面把两件不同的事混成一枚并不存在的字段名。腿的更正**成立**，且腿正确地指出这⛔ 改变偏移结论（24 靠字段序＋`cbSize=22`＋总长 40 三枚闭合，与那句无关）。
- ⇒ 归编排者处置那一句（票面纪律：我只追加 `Progress log` 一行，⛔ 改票面其它一字）。**这条我⛔ 圆场：票面现量节里躺着一条已被盘上原文驳倒的陈述，勾框之前应当改掉或就地打更正标记。**

## 6. 夹具要用的三枚常量（我自己读到的，⛔ 抄腿的）

- `2376:#define  WAVE_FORMAT_EXTENSIBLE                 0xFFFE /* Microsoft */`
- `2110:#define  WAVE_FORMAT_IEEE_FLOAT                 0x0003 /* Microsoft Corporation */`
- `2418:#define WAVE_FORMAT_PCM         1`
- SubFormat GUID 本体（`2475` 逐字）：`DEFINE_GUIDSTRUCT("00000001-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_PCM);`
  ⇒ `Data1` 低字＝**1**；IEEE_FLOAT 那枚在 `2483`（`00000003-...`）⇒ 低字＝**3**。
- ⚠ 一处**我没复现的**：票面与腿都说 `2110` 是 `WAVE_FORMAT_IEEE_FLOAT`，我这把尺同时抽到 `2482: DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_IEEE_FLOAT)`
  ——那是**同一枚值的另一种写法**（把 tag 塞进 GUID），⛔ 两枚互相冲突。本格⛔ 用它做任何判断。

## 7. 本格没做的两件事（写清边界）

- ⛔ **外部文档那一把我没跑**（Microsoft Learn 交叉核）。盘上头文件这一把已经能**独立闭合**（§4 那三枚互验），
  腿那把 Learn 的字段序读数与盘上原文一字不差 ⇒ 我**接受**它为第二把尺、⛔ 重复取回（⛔ 我据它做判语，判语全建在盘上原文上）。
- ⛔ **真机 `GetMixFormat`**（票面 `AC#3`〔仅本机可量，归编排者〕）——我没跑，⛔ 派给我。
  ⇒ "**本机混音格式到底是不是 extensible＋float**"仍**未定**；这一条**不挡 `AC#2`**（改偏移的正确性由 §4 闭合，⛔ 由本机形状决定）。
