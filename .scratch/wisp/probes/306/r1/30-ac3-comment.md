# 306-r1 · 30 `AC#3` 那一格（只改注释，⛔ 与 `AC#1`/`AC#2` 同批＝本腿第 3 笔）

文件＝`internal/audio/wave_format_float_300_windows_test.go`（`//go:build windows`，票 300 `AC#6` 那枚钉）
改动＝**一段注释里的半句假话**，名册：`+16 / -5` 行，全部 `//` 开头（自查尺见 §3）。⛔ 动任何一行断言、⛔ 动任何一枚期望值、
⛔ 删 `waveFormatPCM`、⛔ 顺手改别的注释。

## 1. 那句为什么⛔ 成立（盘上现量，⛔ 抄票面/编排者的数）

尺＝`sed -n`／`grep -n` 对盘上 SDK 头（`C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/`），
件＝`logs/30-sdk-headers.txt`（每把尺自落 `rc=`）：

| 尺 | 逐字读数 |
|---|---|
| `grep -n KSDATAFORMAT_SUBTYPE_IEEE_FLOAT …/shared/mmreg.h` | **命中 4 行**＝`:2480` `#if !defined( STATIC_… )`／`:2481` `#define STATIC_…`／**`:2483` `DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);`**／`:2484` `#define KSDATAFORMAT_SUBTYPE_IEEE_FLOAT DEFINE_GUIDNAMED(…)` ⇒ **`mmreg.h` 里有这份 GUID**，那句 "not in `mmreg.h`" ⛔ 成立 |
| `sed -n '845,860p' …/shared/ksmedia.h` | `ksmedia.h:850` ＝ `#if defined(_INC_MMREG)` 包着 `:851` 的 `#if !defined( STATIC_… )`，`:853` `DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_IEEE_FLOAT)`、**:854** `DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT)` ⇒ `ksmedia.h:850-855` 是**同一枚 GUID 的另一份（替代）定义**，两边⛔ 是"一份在 mmreg／一份在 ksmedia"的两个对象，是**同一枚 GUID 的两处定义** |
| `sed -n '2108,2112p' …/shared/mmreg.h` | `:2110` ＝ `#define  WAVE_FORMAT_IEEE_FLOAT                 0x0003 /* Microsoft Corporation */` ⇒ 那枚钉真正引用的出处＝**格式标签**，⛔ 是 GUID |
| `sed -n '2416,2420p'` / `sed -n '2374,2378p'` `mmreg.h` | `:2418` `#define WAVE_FORMAT_PCM         1`／`:2376` `#define  WAVE_FORMAT_EXTENSIBLE                 0xFFFE` ⇒ 与同文件 `:132-134`／`:161`／`:171`／`:180` 已在用的三枚出处逐字一致（那三处**⛔ 改动**） |

⚠ 与派单措辞的一枚⛔ 同处（具名报回，见 `90-final.md` §冲突）：派单写"`ksmedia.h:851-855` 是另一份替代定义"，盘上那一份的**守卫从 `:850` 的 `#if defined(_INC_MMREG)` 起**、`DEFINE_GUIDSTRUCT` 那行在 **`:854`**（⛔ 是 `:851`）；
`mmreg.h` 侧派单写 `:2480-2484`（✓ 对得上），其 `DEFINE_GUIDSTRUCT` 在 **`:2483`**（✓ 与该文件 `:101` 兄弟件已引的号一致）。
本件注释按盘上写 `ksmedia.h:850-855`，⛔ 抄派单的 `:851-855`。

## 2. 改前／改后**逐字对拉**

改前（`HEAD` blob `:37-43`，尺＝`git diff` 的 `-` 侧，件＝`logs/30-ac3-selfcheck.txt`）：

```go
// Terminology, because AC#6 asked that two things not be merged into one
// citation: the WORD values WAVE_FORMAT_PCM, WAVE_FORMAT_IEEE_FLOAT and
// WAVE_FORMAT_EXTENSIBLE are defined in mmreg.h. The GUID
// KSDATAFORMAT_SUBTYPE_IEEE_FLOAT is a different object and lives in ksmedia.h,
// not in mmreg.h; mmreg.h is also what fixed the field offsets back in AC#0.
// Same header family, different things -- AC#0 pinned the OFFSET, this case
// pins the VALUE.
```

改后（工作树，同一把尺的 `+` 侧）：

```go
// Terminology, because AC#6 asked that two things not be merged into one
// citation: the WORD values WAVE_FORMAT_PCM, WAVE_FORMAT_IEEE_FLOAT and
// WAVE_FORMAT_EXTENSIBLE are defined in mmreg.h (installed 10.0.26100.0 tree:
// :2418, :2110, :2376). The GUID KSDATAFORMAT_SUBTYPE_IEEE_FLOAT is a different
// object, and this paragraph used to say it lives in ksmedia.h and "not in
// mmreg.h". That second half is wrong on disk, and ticket 306 AC#3 corrects
// exactly that sentence and nothing else: mmreg.h DOES define the GUID, at
// mmreg.h:2480-2484, whose DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT)
// sits at :2483; ksmedia.h:850-855 carries an alternative second definition of
// the same GUID, wrapped in a defined(_INC_MMREG) guard and spelled from the same
// DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_IEEE_FLOAT). One definition per header,
// identical value -- what holds is "both headers define it", not "it is not in
// mmreg.h". Neither definition is what this case pins: the authority behind the
// mmregWaveFormatIEEEFloat literal below is the FORMAT TAG at mmreg.h:2110
// (WAVE_FORMAT_IEEE_FLOAT 0x0003), NOT that GUID; the two merely share the value
// 3 because the GUID is built out of the tag. mmreg.h is also what fixed the
// field offsets back in AC#0. Same header family, different things -- AC#0
// pinned the OFFSET, this case pins the VALUE.
```

改后那句写足的三件事＝票面要求的三件事：**两份定义各在何处**（`mmreg.h:2480-2484` 的 `:2483` ＋ `ksmedia.h:850-855` 的 `:854`，
后者由 `defined(_INC_MMREG)` 包着）、**那枚钉引用的出处是格式标签** `mmreg.h:2110`、以及**⛔ 是 GUID**（并给出两者为何共享同一个 `3`：
GUID 由 `DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_IEEE_FLOAT)` 从标签拼出来）。
原来那两句"mmreg.h 定了 field 偏移（`AC#0`）"与"`AC#0` 钉偏移、本格钉值"逐字**保留**。

## 3. "⛔ 动任何断言"的自查尺（四把，各带退码）

| 尺 | 读数 | 件 |
|---|---|---|
| `git diff -- <file> \| grep -E '^[+-]' \| grep -vE '^(\+\+\+\|---)' \| grep -vE '^[+-]//'` | **零命中**（`rc_noncomment_hits=1`＝grep 没找到任何东西）⇒ 增删的每一行都以 `//` 开头，⛔ 一行代码／⛔ 一枚期望值被碰 | `logs/30-ac3-selfcheck.txt` |
| `git diff --numstat -- <file>` | `16 5` ＝ 16 加 5 删，全在上述那一枚射程里 | 同上 |
| `gofmt -l <file>` | **0 行**（`rc_gofmt=0`，⛔ 输出）⇒ 注释改动没把格式尺弄脏（工作树那一把，blob 那一把在 `40-gates.md`） | 同上 |
| `grep -c waveFormatPCM` | `wasapi_windows.go` 仍 1 行（声明 `:97`）、该测试文件仍 7 行 ⇒ **⛔ 删 `waveFormatPCM`**（票 300 `AC#6` 的附带覆盖面还读着它） | 同上 |

## 4. 形状与复绿（票 300 那两枚钉⛔ 被动过形状）

尺＝`GOFLAGS= go test ./internal/audio/ -count=1 -v -run '<指名>'`，件＝`logs/30-ac6-shape.txt`：

- `TestWaveFormatConstantsMatchMmregAuthority300`：顶层 `--- PASS` **1** 枚 ＋ 缩进子测试 **6** 枚
  （`ieee_float_low_word`／`extensible_tag`／`pcm_tag`／`the_three_constants_map_one_to_one_onto_the_authority_table`／
  `parseWaveFormat_expands_a_face_tagged_with_the_authority_extensible_value`／`control_plain_pcm_tag_is_not_expanded`），
  `ok … 0.033s`，**`rc_run_ac6=0`** ⇒ "顶层 1＋子测试 6"那枚形状**逐枚原样**，本格⛔ 为变绿动过它。
- `TestParseWaveFormatSubFormatOffset300`：顶层 1 ＋ 子测试 4（`S1_float_at_24`／`S2_pcm_at_24`／`S3_bogus_subtype_7_at_24`／`S4_data1_shifted_to_26`），**`rc_run_300ac2=0`** ⇒ `parse_wave_format_300_windows_test.go` 的断言⛔ 动（本票只⛔ 让 `:192`/`:196`/`:218` 三枚行内值有牙，⛔ 让票 300 的两枚钉改形状）。
