# 300-v2 — `20` 头顶注释与用例 `why` 的逐句可证伪性（＋指路对不对）

射程＝`internal/audio/wasapi_windows.go:182-199`（改写后的 18 行）＋`internal/audio/parse_wave_format_300_windows_test.go` 头顶与四枚 `why`。
判据＝**这句话现在还有没有一把尺能证伪**；指路另判**文件名＋行号对不对**。
权威＝盘上 `C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`（我这把现开：`ls -l` ⇒ **176,680 字节**，⛔ 我造成）。
⚠ 派单给过的示例路径 `Include/*/um/mmreg.h` 我这把⛔ 用到（真身 `shared/`，腿与编排者都已具名抓过这一枚）。

## 1. 逐句账（★第一列＝注释里的句子，⛔ 转述，逐字自盘上）

| 句子（逐字片段） | ⛔ 有尺／✅ 有尺 | 我这把的尺与读数 |
|---|---|---|
| `WAVEFORMATEX is 18 bytes with tag@0, channels@2, rate@4, avgBytes@8, blockAlign@12, bits@14, cbSize@16` | ✅ 可证伪 | `sed -n '2426,2438p'` ⇒ 字段序逐字 `WORD wFormatTag / WORD nChannels / DWORD nSamplesPerSec / DWORD nAvgBytesPerSec / WORD nBlockAlign / WORD wBitsPerSample / WORD cbSize` ⇒ 2+2+4+4+2+2+2＝**18**，@0/2/4/8/12/14/16 **逐格对上** |
| `which is included under pshpack1.h (packing = 1, so nothing is padded)` | ✅ 可证伪 | `grep -n pshpack\|poppack` ⇒ **`:35`** `#include "pshpack1.h"   /* Assume byte packing throughout */`、`:3769` `#include "poppack.h"` ⇒ 全文件在那一对之间。⚠ **一处限定我这把量到了、注释⛔ 写着**：`:31` 那层 guard 是 `#if !defined( RC_INVOKED ) && defined( _MSC_VER )`（`_MSC_VER<=800` 时走 `#pragma pack(1)`）⇒ "packing=1" 是**按 MSVC 编译这一枚头文件的消费者**成立；`RC_INVOKED`／非 MSVC 那两条路径⛔ 套 pragma。⚠ **腿在 `10` 件 §3 自己写着"我这把⛔ 重开一遍"**、并把"读到还是推的"标在 `90` 件 ⇒ 那一枚**欠的**，我这把**补开了**（本行就是补的读数）。 |
| `the extension -- documented in that header with cbSize = 22 -- adds Samples@18, dwChannelMask@20 and SubFormat@24, 40 bytes in total` | ✅ 可证伪 | `:2521,2534` 现读（字段序 `Format;union{3×WORD}Samples;DWORD dwChannelMask;GUID SubFormat;`）＋`grep -n "cbSize = 22"` ⇒ **`:2540`／`:2550`** 两处逐字 `/* Format.cbSize = 22 */` ⇒ 18＋2＋4＋16＝**40**、`Data1` 起 **24** |
| `Samples@18 is a two-byte union of three WORDs (wValidBitsPerSample, wSamplesPerBlock, wReserved); "validBits" is not a field name.` | ✅ 可证伪 | 同那段现读：三枚成员逐字都是 `WORD`、union 的名字逐字 `} Samples;`；`grep -c ValidSamples mmreg.h` ⇒ **0** ⇒ 标准结构里⛔ 这枚名字（票面 `:13` 那句已被更正的那一枚） |
| `KSDATAFORMAT_SUBTYPE_PCM is 1 and KSDATAFORMAT_SUBTYPE_IEEE_FLOAT is 3, both with a zero high half word` | ✅ 可证伪 | `:2475` `DEFINE_GUIDSTRUCT("00000001-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_PCM)`／`:2483` `DEFINE_GUIDSTRUCT("00000003-…", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT)` ⇒ `Data1`＝1／3 ⇒ 高半字＝0 ✅ |
| **`so a read two bytes late resolves every extensible mix format to 0`** | ⚠ **这句的全称量词眼下⛔ 任何尺能证伪** | 我这把**反证到它的射程边界**：`scan-26.txt` 里 S4 那一形在偏移 26 **读到的是 3、⛔ 0**（因为那枚面把 `Data1` 写在 26）。⇒ "读晚两字节"⛔ 是"恒 0"的机械性质，它是**"当 subtype 的 `Data1` 高半字为 0（且 `Data2` 前两字节为 0）"时才 0**。注释前半句刚枚举了两枚 subtype（都是 0 高半字），后半句却跳成 **every extensible mix format** ⇒ 中间那一步**靠枚举两枚推全称**。四形里⛔ 一枚写 `Data1 ≥ 0x00010000` 的面 ⇒ 全称那一句**没有尺**（⛔ 恒真，也⛔ 证伪）。⚠ 我这把⛔ 判它"假"——真实混音格式（PCM／FLOAT 两枚 subtype）上它成立；我判的是**这句写到了测量之外**。 |
| `because convertPacket branches on floating alone, that would re-slice 32-bit float frames as int16 while returning the very same frame count` | ✅ 可证伪 | 我这把现读函数体：`func convertPacket(data unsafe.Pointer, frames, channels int, floating, silent bool)` ⇒ `if silent \|\| data == nil` 早退（⚠ "branches on floating **alone**"省了这枚早退，但早退⛔ 决定样本宽度，句子在"走哪条切法"这一义上成立）；`if floating { unsafe.Slice((*float32)(data), frames*channels) }`／`i := unsafe.Slice((*int16)(data), frames*channels)`；`bits` 在函数体内**出现 0 次**（尺＝`awk` 抽函数体 `\| grep -c 'bits'` ⇒ **0**）；两条分支都过 `MonoDownmix(…, channels)` ⇒ 输出枚数同＝`frames` ✅ |
| `The sibling wavinjector.go parses the same layout at body+24 under the same 40-byte total: mmreg.h is the authority for both, not that file.` | ✅ 可证伪 | 现读 `wavinjector.go`：**`:192`** `if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]`／**`:193`** `if size < 40 {`／**`:196`** `fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])` ⇒ "同形"三枚数全对上；**并且这句自己就把裁判权推回 `mmreg.h`** ⇒ 与票面 `AC#0` 那枚"⛔ 拿本仓另一枚实现当裁判"的硬门**同向**（⛔ 违门） |

## 2. 指路逐枚（文件名＋行号对不对）

| 出处在哪 | 引的是什么 | 对不对 |
|---|---|---|
| 产码注释 | **只写相对名**：`shared/mmreg.h`／`pshpack1.h`／`cbSize = 22`，⛔ 行号、⛔ 绝对盘路径 | ✅ **判：该留、且形状对**。行号会随头文件腐烂（本仓 `AGENTS.md` §5 给"行号核正"那把尺就是同一枚毛病），绝对路径还会钉死 SDK 版本；相对名两样⛔ 沾。腿在 `10` 件 §2.2 写明了这一条理由，我这把按同一把尺核它的产面：**兑现**。 |
| 用例头顶 `:18` | `C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h` | ⚠ **判：这一枚⛔ 该留在入库件里**（同一枚规矩腿在产码那面守了、在测试这面破了自己）。三条：① **机器专属**＋**SDK 版本专属**（`10.0.26100.0` 换一枚版本这行就指到不存在的文件）；② 与**同批产码注释**的自定口径**互相打脸**（`10` 件 §2.2 逐字"绝对盘路径…⛔ 产码注释"）；③ **没有换来任何东西**——相对名 `shared/mmreg.h` 一样能定位。⚠ 它**⛔ 是假话**（盘上确有此文件、176,680 字节、字段序与引用一致）⇒ 我判的是**易腐烂的指路**，⛔ 判"错"。 |
| 用例 `:24` | `cbSize = 22, mmreg.h:2540 and :2550` | ✅ 两枚都对（`typedef WAVEFORMATEXTENSIBLE WAVEFORMATPCMEX; /* Format.cbSize = 22 */`／`WAVEFORMATIEEEFLOATEX` 同形） |
| 用例 `:58` | `WAVE_FORMAT_EXTENSIBLE (mmreg.h:2376)` | ✅ 对：`:2376` `#define  WAVE_FORMAT_EXTENSIBLE                 0xFFFE /* Microsoft */` |
| 用例 `:98`（S1 `why`） | `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003)` | ✅ 对（那一行逐字就是 `DEFINE_GUIDSTRUCT("00000003-…", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT)`） |
| 用例 `:104`（S2 `why`） | `KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001)` | ❌ **偏 1 行**（★我这把现读：`:2474` ＝ `    DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_PCM)`，那一行**逐字没写** `KSDATAFORMAT_SUBTYPE_PCM` 这枚符号名；符号名在 `:2473` 的 `#define STATIC_…`／`:2475` 的 `DEFINE_GUIDSTRUCT("00000001-…", KSDATAFORMAT_SUBTYPE_PCM)`／`:2476` 的 `DEFINE_GUIDNAMED`）。⇒ **正解＝`:2475`**。⚠ 后果面：按 `:2474` 去开文件的人看到的是**另一枚符号**（同义、名字⛔ 同）⇒ 引路人会以为注释写错了符号。腿在 `90` 件 §2 已**自报这一枚、并⛔ 动它**（理由＝那两枚字符串是 `a1r`／`v1`／编排者三把尺逐字比对的表面，改＝另一笔）⇒ **我复核它的自报＝逐字准确**（含它给的"2474 是 DEFINE_WAVEFORMATEX_GUID(WAVE_FORMAT_PCM)"那一枚判定，与我这把同值）。裁：这一枚**该改、⛔ 混进 `AC#2` 这一批**（腿的处置我对；改法归编排者另裁）。 |

## 3. 用例里那三枚"扫描面"句子（★对拉我 `10` 件那把实测）

- `:11-13`「`300-v1`'s and the orchestrator's own scans put this suite **green at 24** and **red at each of the other offsets they swept (20, 22, 26, 28)**」
  ⇒ ✅ 与我 `10` 件 §1 表**逐枚同值**（顶层 24 `rc=0`／其余四枚 `rc=1`）。★**"they swept" 那三个限定字是这句的命门**：它⛔ 声称"任何别的偏移都红"，只声称"扫过的那四枚红"⇒ **诚实边界**。⚠ 顺带一把：我这把**⛔ 扫全**（只扫了派单点名的五枚），但按夹具算术，四形同时对上的偏移**只可能是 24**（S1/S2/S3 要同一枚 `N` 分别读出 3／1／7 ⇒ `N`⛔ 唯一）⇒ 那句收紧**⛔ 过头、也未失硬**。
  ⇒ ★**这一处是腿自己抓回来的一次过判**（`30` 件 §5 末条：原句"red at **every other one**"改成"each of the other offsets **they swept**"），它**具名**交给我裁 ⇒ **我裁：改对了，且那一发"注释散文改在跑完全量之后"的处理（`cmp` 工作树↔treeB＋补一发 audio-only）也成立**（我这一把是在**改后的最终字节**上重跑的，`10` 件 §2 那两枚 `IDENTICAL` 就是它的还原尺）。
- `:84-86`「**S4 is the weakest face in the set — on its own it stays green at offsets 22, 24 and 28 — and S1 on its own stays green at 20 and 24. No single face here is the decisive one; the four together only accept 24.」
  ⇒ ✅ **三枚分句逐枚对上我这把实测**：S4 绿＝{22,24,28}、S1 绿＝{20,24}、"⛔ 单枚决定"＝每一枚都至少放行一枚错偏移（S2/S3 只放行 0 枚 ⇒ "no single face is decisive"这一句按"⛔ 一枚单独把范围缩到 {24}"读是**成立**的：S2 把范围缩到 {24}…… ⚠ **等等，这一句我这把判它⛔ 严丝合缝**：S2 单形绿＝{24}（见 `10` 件 §1 表），所以**"no single face is the decisive one"这一句与我这把实测冲突**——S2（以及 S3）单形就把范围钉死在 24。我这把的读数是：**S2／S3 各单枚即只绿 24 ⇒ 它们是决定性的**；"no single face"只对 S1／S4 成立。⇒ 见 §4 判语，这一枚是**用例注释里唯一一句与实测不同向**的（措辞层面，⛔ 判据层面；判据⛔ 因此变弱，四形全绿仍是 24）。
- `:15`「an always-red tracked test is banned here」⇒ ✅ 与票面"编排者补格"那枚定式同文；我这把量到入库那一枚**一进门就绿**（`rc=0`）⇒ 定式⛔ 被破〔盘上现量〕。

## 4. 本件判语（⛔ 我修，只裁）

- **产码那段注释：留、且逐句有尺**，只有两处该记：① 全称那句 `every extensible mix format → 0` 写到了测量之外（§1 倒数第 3 行）；② `branches on floating alone` 省了 `silent‖nil` 早退（§1 最后一行前一格，义⛔ 变）。两枚都＝**散文精度**，⛔ 判据、⛔ 代码。
- **用例头顶那枚绝对路径：判⛔ 该留**（§2 第二行三条），**改＝另一笔**。
- **`mmreg.h:2474`：判该改成 `:2475`**，腿的"⛔ 混进本批"处置我**认可**；改法与那一笔归编排者裁。
- **`:84-86` 里"No single face here is the decisive one"与实测不同向**（S2／S3 单形即决定性）⇒ 我这把的读数支持把它写成"**S2/S3 才是承载判别力的两枚，S1/S4 各放行错偏移**"（这也正是 `300-v1` 自己给的理由，见票面 13:0x 那节与腿 `20` 件 §3 末行——**用例注释把 v1 那句"S4 最弱"抄对了，把"判别力扛在 S2/S3"抄丢了**）〔盘上现量：`10` 件 §1 表〕。
- ★**"引用仓外绝对路径该不该留"这一枚，本仓有先例可引**：腿⛔ 在注释里写现在时计数（会被下一位改动打成假话）＝同一族"注释里写着会腐烂的东西"。绝对盘路径**腐烂得更快**（换 SDK 版本即断）。⇒ 我这把把两枚（`:2474` 行号、绝对路径）**归成同一类：入库件里易腐烂的外部坐标**，⛔ 一枚属"假"、⛔ 一枚属"错判据"，但都该由**编排者另起一笔**处置，⛔ 塞进 `AC#2`。
