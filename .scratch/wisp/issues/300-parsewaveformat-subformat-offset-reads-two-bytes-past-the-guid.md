# 票 300 — `parseWaveFormat` 把 `WAVEFORMATEXTENSIBLE` 的 `SubFormat` 读在**第 26 字节**，而同包另一枚实现读在**第 24**：两把尺互相矛盾，且其中一枚错的话，**extensible 混音格式的浮点样本会一律被当 16 位整型切**（只读普查 `297-a1` 浮出；⚠ 本票**没有**真机读数，全部是读码＋算术）

**立票**：2026-10-10 10:5x 编排者（来路＝票 297 只读普查腿 `297-a1` 件 `10`／`20`，它具名顶回票 297 现量节那句"切错字节这一形在仓内不可能被测到"）
**性质**：★**这是一枚"嫌疑缺陷＋可测性"票，⛔ 不是已定案的 bug**。今天盘上只有两件事是硬的：① 本仓**同一包内**有两枚实现读同一个字段、读在**不同偏移**；② 其中一边的算术与另一边的守卫对得上（`cbSize=22` ⇒ 总长 40 ⇒ `SubFormat@24`）。⇒ 决定性那一发是 `AC#1`（**包内、零真设备**），⛔ 任何人在它交出读数之前都不许把这枚嫌疑写成结论。
**为什么要紧**：如果 `SubFormat@26` 是错的，那么 `tag` 对任何 extensible 混音格式都会读成 **0**（既不是 1 也不是 3）⇒ `floating=false` ⇒ **32 位浮点样本被当 int16 切**，而 `MonoDownmix` 的 `n = len(in)/chans` 会把"少读了一半字节"这件事**从计数上完全藏起来**（输出枚数与正确形状一模一样）⇒ 只剩幅度是坏的。这**正好是票 297 那个怪形的形状**：六窗 mean 全挤在 0.367–0.383、`sound` 不比 `quiet` 响、甚至 `quiet-B` 比 `sound` 高——一把对**枚数**敏感、对**幅度**才敏感的尺，就是它。

## 现量（每条带尺；⛔ 引用前先重跑，行号一律当快照）

- ★**矛盾本身（我 10:5x 两把尺现读，⛔ 不是转述）**：
  - 读在 **26** 的那枚＝`internal/audio/wasapi_windows.go`，`parseWaveFormat` 体内逐字 `		sub := *(*windows.GUID)(unsafe.Add(p, 26))` ＋下一行 `		f.tag = uint16(sub.Data1) // KSDATAFORMAT_SUBTYPE_* low word: 1=PCM 3=FLOAT`；而**同文件注释逐字写着它自己的布局假设**：`// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].`
  - 读在 **24** 的那枚＝同包 `internal/audio/wavinjector.go`，逐字 `				fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])`，头顶注释逐字 `			if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]`，守卫逐字 `				if size < 40 {`。
  - 尺＝`grep -rn "SubFormat\|ChannelMask\|channelMask" --include=*.go internal cmd` ⇒ **4 枚命中，全在这两枚文件里**（⚠ 没有第三枚实现可仲裁）。
- ⚠**算术这一半是我推的，⛔ 是读数**（写清楚）：`WAVEFORMATEX` 自身 18 字节；`cbSize` 文档值为 **22**；22 ＝ 2 枚补空 ＋ 4 枚 `dwChannelMask` ＋ 16 枚 `GUID` ⇒ `SubFormat` 落在 **24**、结构总长 **40**——这与 `wavinjector` 那枚 `size < 40` 守卫**自洽**，与 `parseWaveFormat` 注释里"validBits 单独占 18–20"那一枚假设**不自洽**（标准结构里没有独立的 validBits 字段，`nValidSamples` 是与 `dwChannelMask` **同一个 union** 的成员）。⇒ 判语只能写到"**24 那一边多一枚内部证据（守卫 40）、26 那一边多一枚假设性注释**"，⛔ 判"26 一定是错的"。
- **`convertPacket` 的下游形状（现读全文）**：`internal/audio/wasapi_windows.go` 逐字 `func convertPacket(data unsafe.Pointer, frames, channels int, floating, silent bool) []int16 {` → `	if floating {` → `		f := unsafe.Slice((*float32)(data), frames*channels)` → `		return MonoDownmix(FloatToPCM16(f), channels)` → `	}` → `	i := unsafe.Slice((*int16)(data), frames*channels)`。**`bits` 在这枚函数体里出现 0 次**（尺＝我现读函数体）⇒ 走哪条分支**只看 `floating`**，而 `floating` 就是那枚 `tag == waveFormatFloat`。
- **仓里没有可用的 extensible 样本**（尺＝`find third_party -name '*.wav'` ⇒ **19 枚**，`xxd -l 64` 抽看前四枚全部是 `fmt ` 块长 `0x10`＝16、tag `0x0001`＝plain PCM）⇒ ⛔ 拿真wav顶 `AC#1`，夹具得自己造字节。⚠ 这也把票 297 `AC#1` 那句"仓里现成 wav"的射程收窄了：**现成的是台件，不是 extensible 夹具**。
- **可测性（`297-a1` 报，我⛔ 复跑）**：`parseWaveFormat` 收的是 `unsafe.Pointer`，包内 windows-tagged 用例可以拿一枚 Go 自己 `[]byte` 的首地址喂进去（⛔ 假设备、⛔ 不占麦）；`convertPacket` 同形。⇒ 票面那句"**在仓内不可能被测到**"（它引 `:16`）**不成立**，我接受更正并写在本票 `AC#1`。
- ⚠**这枚缺陷今天在生产里到底响不响，取决于本机混音格式是不是 extensible＋float**（Windows 共享模式常见如此，但**我没量**）⇒ 归 `AC#3` 那一发真机（〔仅本机可量，⛔ 归编排者跑〕＝只调 `GetMixFormat`、把前 40 字节 hexdump 出来并分别按 24／26 取 `Data1` 低字，⛔ 不 `Initialize`／⛔ 不 `Start`、不占麦）。

## 要建什么（先测再修；次序⛔ 不许换）

- [x] **`AC#0` 权威偏移（⛔ 不许拿本仓另一枚实现当裁判，⛔ 不许凭记忆）**：具名回答"`WAVEFORMATEXTENSIBLE` 里 `SubFormat` 的偏移到底是 24 还是 26"的**外部凭据**——`mmreg.h` 的字段序 ＋ `cbSize` 文档值 ＋ union 的对齐规则，或直接引一枚权威头文件/文档段落（带可核的出处）。★**为什么这格必须挡在 `AC#1` 前面**：`AC#1` 的夹具本身就是按某种布局造的，**如果布局假设错，那枚用例的红什么也证不了（循环）**。⇒ 派单里写死：**`AC#1` 的凭据只有配合 `AC#0` 才成立**；两者不一致时**停手上报**，⛔ 由腿自己改布局再跑一遍。
- [x] **`AC#1` 决定性一发（⛔ 先红再修，本格只产**测试侧**夹具，⛔ 不改产码语义）**：在 `internal/audio` 新增一枚 windows-tagged 用例，**自己造** 40 字节的 `WAVEFORMATEXTENSIBLE`（`wFormatTag=0xFFFE`、`cbSize=22`、`bits` 各设 16／32、`SubFormat.Data1` 分别置 **1**（PCM）与 **3**（FLOAT）），调 `parseWaveFormat` 并断言返回的 `tag` **等于 1／3**。★判据形状＝**不许**断"被调用过"，断的是**解析出来的那个值**；⚠ 预期今天**两枚都红**（都读出 0）——那就是本票唯一需要的"改前必红"凭据。⛔ 不许为了让它绿而同时改产码（先交红，修归 `AC#2`）。⛔ 不许用真设备。
- [x] **`AC#2` 落地那一发（⛔ 只有 `AC#1` 红过并经非实现者裁之后才动）**：把偏移改成与权威布局同形的那一枚，**同批**把 `parseWaveFormat` 头顶那三行布局注释改成实话（⚠ 注释现在逐字教读者相信 `SubFormat@26`，这比代码更危险——它会喂下一个改这里的人）；⛔ 顺手格式化别人的文件；⛔ 动电平尺定义（`internal/audio/level.go` 的 `SineLevelTolerance`／`internal/ball/liquid.go:30` 的 `SilenceLevelGate = 0.06` 属 C21 冻结 token 面）；⛔ 动 `internal/observe/thresholds.go`（尺＝208 行、`level` token **0 命中**，我现跑）。
- [x] **`AC#3` 与票 297 那三形的关系（判语归位，⛔ 不越权结案）**：`AC#1`／`AC#2` 成立⇒ 票 297 三支里的 ③（"切错字节"）拿到一枚**可复现的形状**，但⛔ 本票不许把 297 的"是什么"判成已定案——真机那枚混音格式到底是不是 extensible＋float，归**编排者**那一发 `GetMixFormat` dump（〔仅本机可量〕）。⇒ 两票各自留格，⛔ 合并结案。
- [ ] **`AC#4` 门禁＋越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；格式两把尺并排（工作树＋**仓外** blob 各一把）新增 0 枚；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/audio/ ./cmd/wisp/ -count=1 -v` 改前改后各 ≥2 发取**红名交集**＝新增红 0（⚠ 三数必带尺名：`^--- PASS` 顶层 vs 含子测试）；每枚门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；`git show --stat` 名册只含 `internal/audio/**`＋`probes/300/**`；⛔ `frontend/**`／`design/**`／三枚冻结件／`thresholds.go`／`allowlist.txt`／D43 表／`PLAN.md` 零字节；⛔ 零 push、commit 必带显式 pathspec。
- [ ] **`AC#5` CI 可见性（这格会决定落点，⛔ 不许按 API 名字推）**：`297-a1` 报「`./internal/audio/` 只在 core(ubuntu) 的 scope 里、**不在 windows/cli** ⇒ `internal/audio` 的 windows-tagged 用例**整条 CI 无一档真跑**」（我⛔ 复跑；我只现读到 `scripts/portable-tests.sh:241` 那行确实把 `./internal/audio/...` 写进某一枚 scope、`:593` 把 `TestLiveWasapiSmoke` 列为 fixture 例外）。⇒ 落地腿**自己现跑那三跳**（`.github/workflows/ci.yml` → `scripts/wisp-cli-tests.sh` → `portable-tests.sh --scope=<名>`）并具名回答：**本票新造的用例到底进不进 CI 执行面**；不进就写清"凭据只有本机"，⛔ 把用例落点迁就门禁来选（那等于把仪器藏起来）。

## 边界与已知禁区（⛔ 派单要逐字带上）

- ⛔ **本票不许改票 297 的判语**，也不许把 297 的 `AC#2`"三形择一"顺手裁掉——两票各自交格。
- ⛔ 真设备／真麦克风那一发归编排者（会占麦、机主在场才做）；⛔ 派给腿。
- ⚠ **仪器面**：票 255 的名册（`cmd/wisp/config_readers_255.go`）用**行号**引产码行，本票若动 `cmd/wisp` 的**行数**会打红它（今天已有一枚同族红由写腿当场复位，见 `A799`）⇒ 本票射程只写 `internal/audio/**`，⛔ 动 `cmd/wisp`；若真要动，先具名解冻。
- ⚠ **`297-a1` 的这枚嫌疑⛔ 经过非实现者**：它与"票 297 现量节"两句冲突（`:16` 可测性／`:18` 已知幅度尺），编排者只复跑了它读数里**能现读**的部分。⇒ 派 `300-a1`？⛔——`AC#1` 本身就是可跑的裁判，**直接派落地腿前的那一发测试**比再派一枚读腿便宜，且它自己会红。

**Status:** **未开工**（本票 10:5x 立，六格 `AC#0`..`AC#5` 全未勾；`AC#0`→`AC#1` 有**次序硬门**：布局权威没定，用例的红什么也证不了）。⛔ 零翻框、⛔ 零产码、⛔ 零 push。

**Status:** **本票的最新态（2026-10-10 12:4x，编排者；⛔ 上面那条原句一字不改——它证的还是"立票那一刻"，这一条为准）**：`AC#0`／`AC#1` **两格已勾**（凭据＝我自己现跑，见「编排者裁定（2026-10-10 12:4x）」节 ★ⓑ／★ⓒ），盘上剩 **4 枚 `- [ ]`**＝`AC#2`..`AC#5` ⇒ ⛔ 本票⛔ 改 `-done`、⛔ 零产码。**★这票的形状从"嫌疑"升成"已证的偏移错"**：盘上权威头文件（`shared/mmreg.h`，`pshpack1`，`cbSize = 22` 两处）⇒ `SubFormat` 在 **24**、`wasapi_windows.go` 那枚 **26 是错的**；决定性台件两形我自己配对跑过＝原样 **4/4 红**、只把 26 写成 24 ⇒ **4/4 绿**（还原 `cmp` 同值）。⚠ 但**最硬只到这一句**：extensible 混音格式下 `tag` 恒 0 ⇒ `floating=false` ⇒ 浮点样本按 int16 切；⛔ 说到"票 297 那个 0.3677 就是它"——那一支仍只有真机 `GetMixFormat` 那一对量得到（归我那一发，⛔ 归腿、⛔ 需机主在场）。★`AC#2`（改产码）**开工权在 `300-v1`**（票面那条硬门我自己写的，我⛔ 绕）；它同时在飞。

## 编排者补格（2026-10-10 12:2x；⛔ 本节不改写上面任何一行、不勾任何框——只补一条**派单里必须逐字带上**的落地纪律）

- ★**`AC#1` 那枚用例的形状要改一支（我裁，⛔ 让恒红的跟踪测试入库）**：票面写的是"在 `internal/audio` 新增一枚 windows-tagged 用例"，而它**今天必然红**（预期两枚断言都读出 0）。本仓定式＝**恒红测试⛔ 入库**（`A799`/`A801` 那条族：红的跟踪测试会被下一位当成"这棵树本来就红"，还会污染别人的红名册作差）。⇒ 现行射程＝**在仓外成对树里跑那一发**（`git archive` 把当前 `HEAD` 导到仓外，把用例加在**那份导出树**里，`go test` 在那儿跑，红句**逐字**抄进 `.scratch/wisp/probes/300/**` 的 `.md`，用例源码以代码块整块进件），⛔ 往仓里落这枚测试。⇒ 真正的跟踪测试随 **`AC#2`**（修完）同批入库，那时它应当**一进门就绿**，而它的"改前必红"凭据＝本程仓外那一发的逐字读数。⚠ 这不改变 `AC#1` 的判据形状（断的是**解析出来的值**、⛔ 断"被调用过"）。
- ★**`AC#0` 的取法我给一条更便宜的第一选择**：本机若有 Windows SDK 头文件（`mmreg.h`／`mmdeviceapi.h` 之类）就**先读盘上那份**并逐字抄 `WAVEFORMATEXTENSIBLE` 的字段序＋`cbSize` 注释；盘上没有再取外部文档。**两把都要带出处**（文件全路径＋那一行逐字，或 URL＋取回日期），⛔ 拿本仓另一枚实现当裁判，⛔ 凭记忆。
- ⛔ 两格⛔ 由同一枚腿顺手做掉修：`AC#0`＋`AC#1` 交回后**归我裁**，`AC#2` 要**非实现者**裁过才派。


## 编排者裁定（2026-10-10 12:4x，来路＝只读＋台件腿 `300-a1r`，四笔 `70a6b9d9`→`cc878604`→`1c007bf6`→`efd85155`，件 `.scratch/wisp/probes/300/a1r/`＝四枚 `.md` 我逐枚现量（`00` 4,612／`10` 12,509／`20` 17,598／`90` 11,103）＋`logs/` 8 枚（**0 字节＝0 枚**；⚠ 一处过程痕迹：`hygiene-build.txt` 在 `1c007bf6` 那笔里是 **0 字节**入库、下一笔 `efd85155` 才填上 191 字节——按本仓"0 字节的件＝那格没交"的尺，那一笔当时是不合格交付，它自己补对了，⛔ 抹）；★**承重两发我自己重跑**，见 ★ⓑ／★ⓒ）

- ⓐ **翻两格：`AC#0`／`AC#1` ⇒ `[x]`**（两格要的是"外部权威定没定＋那枚红交没交出来"，属我现跑可裁；`AC#2`..`AC#5` ⛔ 一枚不碰——`AC#2` 开工权在**非实现者** `300-v1` 手上，见 ⓔ）。⛔ 本票⛔ 改 `-done`、⛔ 动产码、⛔ push。
- ★ⓑ **`AC#0` 我自己开那枚头文件读了（⛔ 照抄腿的抄件）**：文件＝`C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`（`ls -la` ⇒ 176,680 字节，⛔ 我造成、也⛔ 我改过）。逐字（我 `sed -n '2521,2534p'` 的 stdout）＝
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
  同文件我另两把尺：packing＝`:31-36` 那一段逐字含 `#include "pshpack1.h"   /* Assume byte packing throughout */`（⇒⛔ 补空）；`cbSize` 文档值＝`grep -n "cbSize = 22"` ⇒ 命中 **`:2540`**（`typedef WAVEFORMATEXTENSIBLE WAVEFORMATPCMEX; /* Format.cbSize = 22 */`）与 **`:2550`**（`…WAVEFORMATIEEEFLOATEX; /* Format.cbSize = 22 */`）。⇒ **我自己算的偏移**＝`WAVEFORMATEX` 2+2+4+4+2+2+2＝**18**、`Samples`＝2（**@18**）、`dwChannelMask`＝4（**@20**）、`SubFormat`＝16（**@24**）⇒ 总长 **40**、`Data1` 起始 **24**。⇒ **判：24 对、`parseWaveFormat` 那枚 26 错。**
- ⚠**★追加更正（⛔ 改上面「现量」第 13 行原句，只追加）**：我那句"`nValidSamples` 是与 `dwChannelMask` **同一个 union** 的成员"——**错**。逐字对照上面那段：union 的名字是 `Samples`、成员是 `wValidBitsPerSample`／`wSamplesPerBlock`／`wReserved` 三枚 `WORD`，而 `dwChannelMask` 是**独立成员**；标准结构里⛔ `wValidSamples` 这个名字（那是 `WAVEFORMATEXTENSIBLE` 早期文档里的叫法，⛔ 这枚头文件里的）。腿具名顶我，✅接受。**结论层不变**（偏移仍是 24），错的只是我用来解释"为什么 26 不自洽"的那句机制。**另一枚同时入账**：我给腿的派单里那句示例路径 `Include/*/um/mmreg.h` 在盘上⛔ 成立（真身 `shared/`）⇒ 定式并回：**派单里出现盘上路径，落笔前先 `ls`／`find` 一发**（与 `A797` 那条"声明 HEAD 对象层射程之前先 `git show`"同族）。
- ★ⓒ **`AC#1` 那两发我自己重跑（⛔ 引用腿的数当凭据）**：台件树＝`/tmp/wisp300-a1r/tree`（`git archive b2933dc9` 导出，腿留着⛔ 删）。**先核未突变**：`grep -n "sub := "` ⇒ 逐字 `		sub := *(*windows.GUID)(unsafe.Add(p, 26))`（⇒ 我下面这发是**改前**）。
  - **第一发＝原样（偏移 26）**：`--- FAIL` ×4 ＋包行 `FAIL github.com/CarlosShao/wisp/internal/audio 0.026s`。四形红句逐字（我的 stdout）＝`S1_float_at_24: tag = 0, want 3 (KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24)` ＋ `floating = false, want true (convertPacket would take the float32 path)`／`S2_pcm_at_24: tag = 0, want 1 (KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001) at byte 24)`／`S3_bogus_subtype_7_at_24: tag = 0, want 7 (positive control: an obviously wrong Data1 must NOT be swallowed into PCM/FLOAT…)`／`S4_data1_shifted_to_26: tag = 3, want 0 (positive control for offset sensitivity…)` ＋ `floating = true, want false`。
  - **第二发＝只把树上那枚 `26` 写成 `24`**：`--- PASS` ×4 ＋ `ok github.com/CarlosShao/wisp/internal/audio 0.027s`。**还原证明**＝同一条命令内 `cp` 备份 → 跑 → `cp` 写回 → `cmp` ⇒ `restore_cmp=IDENTICAL`，且尺 `grep -c 'unsafe.Add(p, 26)'` 回到 **1**。⇒ 两形配对成立，这把尺**不是**恒真那三形（同一条断言在 26 全红、在 24 全绿）。
  - ★**由这两发现在能说的最硬一句**：extensible 混音格式下 `tag` 恒读 **0** ⇒ `floating=false` ⇒ `convertPacket` 走 int16 那条 `unsafe.Slice((*int16)(data), frames*channels)`。**⛔** 这一句仍⛔ 等于"票 297 那 0.3677 就是它造成的"——见 ⓓ。
- ⓓ **夹具布局与权威对拉（我自己逐字节读了 `build()`）**：`b[0]=0xFFFE`／`b[2]=channels`／`b[4]=rate`／`b[8]=avgBytes`／`b[12]=blockAlign`／`b[14]=bits`／`b[16]=22`／`b[18]=bits`（写进 `Samples.wValidBitsPerSample`）／`b[20]=dwChannelMask`／`Data1` 写在 `f.data1At`（S1/S2/S3＝**24**，S4＝26）⇒ **与上面那段头文件逐格对得上**，⛔ 拿仓里任一枚实现当布局裁判 ⇒ 我立票时那条反循环硬门（"夹具按哪种布局造本身就预设了答案"）**没有**被踩。
- ⓔ **腿抛来的四枚争议，我的处置**：① 那句 union／成员名——✅接受，见 ★ⓑ 那条追加更正。② 派单路径示例——✅接受，记我。③ **`AC#2` 入库版要不要保留 S4 那一形**——⛔ 我现在裁"保留"，**裁法交 `300-v1`**（它判"入库版少了 S4 会不会变成一把只能单向响的尺"）；但我先在 `AC#2` 的派单里钉死 **⛔ 删 S4**。④ ★**新事实，直接改写我 1 小时前刚立的定式**：`git ls-files third_party/sherpa-onnx | wc -l` ⇒ **0**，而盘上有三枚 DLL（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`）⇒ **`git archive` 的导出树里那条 harness PATH 指向不存在的目录**，"成对导出树"法⛔ 自带这一处 handicap；`300-a1r` 是 `cp` 三枚进树才让尺成立的。⇒ 台账 **`A802` 追加更正 `A801` 那条定式**：用导出树跑带 harness 的包之前先 `ls <导出树>/third_party/sherpa-onnx/*.dll | wc -l` 并具名（0 枚就先 `cp`）。⚠ 连带一枚我⛔ 现在不下结论的问题：`293-v1` 的件里⛔ 提过 DLL 这回事（我 `grep -in "dll" .scratch/wisp/probes/293/v1/20-ac4-gates.md` ⇒ **0 命中**），它那两组 677／684 要不要打折——**交 `300-v1` 的④三档裁**（无影响／射程受限／要打折），我⛔ 先替它判。
- ⓕ **排程（今）**：在飞＝**1 枚** `300-v1`（攻 ★ⓑ／★ⓒ 这套凭据＋裁 `AC#2` 能不能开工＋回答 ⓔ④ 那枚打折问题）。**它交回之前⛔ 派 `300-r1`、⛔ 任何人改 `internal/audio`**。之后序＝`300-r1`（`AC#2`：偏移改 **24** ＋同批把 `// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].` 那三行注释改成实话——★那行注释今天**逐字教下一个改这里的人相信 26**，比代码危险）→ `300-v2`（非实现者裁 `AC#2`）→ `AC#3` 与票 297 的关系归位（⛔ 合并结案）→ `AC#5` CI 可见性那三跳。**队列其余**：`296-r2` → `298-r1` → `295-r1` → `294-r1` → `294-v1`（⛔ 同包写腿⛔ 并发；`cmd/wisp` 被 `300-v1` 的只读射程让位，它⛔ 改）。票 299 按在 `Q-84`（⛔ 不催）。⛔ 零 push。


## 编排者现量（2026-10-10 12:5x，⛔ 翻框、⛔ 不改上面任何原句；★这节是**预跑给落地腿的料**，它仍要自己逐跳复跑并在件里具名回答）

`AC#5` 那三跳我 12:5x 自己走了一遍（全部只读尺，⛔ 编译面 0 发、⛔ 与在飞的 `300-v1` 抢仪器）。行号＝今天的快照，引用前重跑。

- ★**结论先写死：本票新造的 `//go:build windows` 用例进不了 CI 任何一档执行面，而且⛔ 任何一枚门会为此变红。** 四跳链条逐字如下（我每一步都 `sed`/`grep` 现读到那行本身，⛔ 按文件名推）：
  - 跳① `.github/workflows/ci.yml:403` `test-core:` → `:403` `runs-on: ubuntu-latest` → `:470` `run: bash scripts/portable-tests.sh --scope=core`；`scripts/portable-tests.sh:234` `core)` 的清单里 `:241` 确有 `./internal/audio/...`，而 `:168` 的 `core_pin` 里确有 `github.com/CarlosShao/wisp/internal/audio` ⇒ **包在名册上是被"core"认领的**，但 core 那档跑在 **ubuntu**，windows-tagged 文件⛔ 进不了那一档的测试二进制——这⛔ 我的推断，是本仓脚本自己写的：`scripts/portable-tests.sh:565-566` 那两行注释逐字 `# GUARD A: the loud empty denominator. \`.TestGoFiles\`+\`.XTestGoFiles\` are the` / `# files that COMPILE INTO THIS PLATFORM's test binary, so a build tag that`，而 `:364` 那行逐字 `        counts=$(go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}' "$p" 2>/dev/null || echo '?/?')` ⇒ **census 数的也是"编进本平台二进制"的测试文件枚数**。⇒ 用例存在、⛔ 执行。
  - 跳② `ci.yml:516` `test-windows:` → `:517` `runs-on: windows-latest` → `:604` `run: bash scripts/wisp-cli-tests.sh` → `scripts/wisp-cli-tests.sh:113` 逐字 `bash "$portable" --scope=cli` → `portable-tests.sh:256` `cli)` ＋ `:257` 那行注释逐字 `    # cmd/wisp's own tests. Only scripts/wisp-cli-tests.sh may call this, because`（下文 `:258-259` 给的理由＝测试二进制没有 sherpa DLL 会在 **LOAD** 期死）＋ `:260` `scope=(./cmd/wisp/)`。⇒ **cli 档只有 `cmd/wisp`，⛔ audio。**
  - 跳③ 同一个 `test-windows` 作业里还有 `ci.yml:780` `run: bash scripts/portable-tests.sh --scope=windows` → `portable-tests.sh:248-253` 的 `windows)` 清单逐枚＝`./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/ ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/ ./internal/session/ ./internal/projctx/` ⇒ **⛔ `./internal/audio/`**。⚠ 这一枚是 `297-a1` 那句"audio 不在 windows/cli"里我⛔ 复跑过的那一半，现在逐行对上。
  - 跳④（★**这半枚是 `297-a1` 没报的，而我读了才敢说"门拦不住"**）`ci.yml:744` `run: bash scripts/portable-tests.sh --scope=census` ＝ GUARD D 唯一调用点。它的判定是**包级**而⛔ 用例级：`portable-tests.sh:382-397` 只对 `$all` 里每枚**导入路径**问"有没有任一档认领它"（`NO-SCOPE`＋`counts` 非 `0/0` 才记 `unclaimed`），而 `internal/audio` 已被 `core` 认领 ⇒ **它永远不会进 `unclaimed`**。⇒ **本仓那把"CI 只测 33 个包里的 20 个"的尺（`:413` 那句自己的话）能看见"整枚包没人测"，⛔ 看不见"一枚包的 windows 那一半没人测"** ⇒ ⛔ 任何人以为"用例落进 internal/audio 就有 CI 兜着"是**读反了**。
- ⇒ **判语只能写到这一句（⛔ 多说没有尺）**：本票的凭据**只有本机**（编排者那台 Windows 开发机，成对导出树那一形）。落地面⛔ 因此被卡——`parseWaveFormat` 是**包内未导出**函数 ⇒ 判据只能长在 `audio` 包里，⛔ "为了让它进 CI 就把用例搬到 `cmd/wisp`"这条路**物理上不存在**（除非先导出，而那是另一枚射程、⛔ 本票）。
- ★**如果要把 CI 可见性真买回来，最便宜的形我量到这里（⛔ 本票不落地、⛔ 由我另裁）**：把 `./internal/audio/` 拉进 `portable-tests.sh:248` 那枚 `windows)` 清单，并**同一笔 commit** 里更新 `:193` 起的 `win_pin`（★**两枚门各管一半**：`:417-418` 那两行逐字 `"…Pull the package"` / `"…into a named scope and update that tier's pin in the SAME commit, or"` 是 **GUARD D** 的报错正文，而"清单变了而 pin 没跟着变"由 **GUARD C** 点红）。⚠ 代价我⛔ 全量：这一改会把该包里**其余** windows-tagged 用例一起拉进分母（含 `:593` 那行 ledger 里已登记的 `TestLiveWasapiSmoke|./internal/audio/|windows|fixture`——**那行本来就按"audio 会在 windows 档"的形状写着**，我这把⛔ 验证它今天是否 inert），⛔ 拿"改一枚清单"当零成本。**归口＝一枚独立的仪器票（票 255 那一族），⛔ 塞进本票 `AC#2` 那一批。**
- ⚠**排程事实（记在本票防下一位重踩）**：`300-v1` 的锚＝`efd85155`；本票面这节与 `2a74812f`／`3fc002e5` 两笔都⛔ 碰源码（写面＝`docs/reports/**`＋本 `.md`）⇒ 它那四条必答的盘上锚点⛔ 因我漂移。

## 编排者裁定（2026-10-10 13:0x，来路＝非实现者验收腿 `300-v1`，五笔 `25a36ccb`→`c7881f56`→`df41c6a8`→`890cd466`→`3da1a1a1`；件 `.scratch/wisp/probes/300/v1/`＝五枚 `.md` 我逐枚 `find -printf '%s %p\n'` 现量（`00` 5,657／`10` 10,287／`20` 11,716／`30` 14,327／`99` 2,924）＋`logs/` 17 枚、**0 字节＝0 枚**；逐笔 `git diff-tree -r --no-commit-id --name-only` ⇒ **越界 0 枚**，名册只含本票 `probes/300/v1/**` 与本票票面）

⚠ 本节⛔ 翻框（`AC#2`..`AC#5` 仍 4 枚 `- [ ]`；翻框＝落地腿交回并经 `300-v2` 之后）。★`AC#2` 的**开工权我在此放行并已派 `300-r1`**——票面 `AC#2` 我自己写的那句"⛔ 只有 `AC#1` 红过并经非实现者裁之后才动"到这一行**闭合**（红过＝`AC#1` 三把尺；非实现者裁＝本腿"开"）。

- ★**它四条必答里我能坐实的两枚，都用我自己的手复跑过**：
  - **②「`AC#1` 是不是循环论证」＝⛔ 循环，我接受**，而凭据⛔ 它那句"夹具落在我自算的权威表上"（那句我 12:4x 已自己读 `build()` 逐字节核过）。我另外把它**没做的那一发**补跑了＝**五枚偏移扫描**（尺在 `.scratch/wisp/probes/300/orch/2026-10-10-orch-five-offset-scan.txt`，仓外导出树、`sed` 射程只打 `sub := *(*windows.GUID)(unsafe.Add(p, N))` 那一行、跑完 `cmp` 还原＝IDENTICAL、`grep -c 'unsafe.Add(p, 26)'` 回到 1）：**只有 24 四形全绿**（`rc-gotest=0`），26 四形**全红**，20／22／28 各只绿一枚。⇒ "一把尺两种形状都放行"这一形被排除，`AC#1` 的牙**由这枚扫描而硬**。
  - ⚠⚠**我自己在这发里踩了两枚坑，具名留痕⛔ 抹**（都写进上面那件的头两行与末行）：**第一发我把 `sed` 射程打宽到整文件** ⇒ 连 `channels@2`／`rate@4`／`bits@14` 一起改了，五档全红＝**那枚台件无效**（⛔ 任何读数取自它）；**第二发的 `SUBTEST tally` 我把 `grep -c` 打在整份日志而⛔ 只打当节** ⇒ 印出来的是**累计值**、⛔ 是一把尺（有效读数只有逐节 `rc-gotest`，与我事后按节重算的 20→1/3、22→1/3、**24→4/0**、26→0/4、28→1/3 那张表）。⇒ 定式并回：**"逐节计数"这类尺必须把 grep 的射程钉在当节内**（先例＝票 293 那把"三数必带尺名"，这次中尺的是我自己）。
  - **③「S4 那一形留不留」＝留，而我采它的理由、⛔ `300-a1r` 的理由**：我这把现量＝S4 在 **20/22/24/28 全都绿、只在 26 红** ⇒ 它是四形里**判别力最弱**的一枚（`300-a1r` 说它用来"杀读 26 那一形"——按我这把**恰好说反**）；真正扛判别力的是 S1+S2+S3 **同时**只在 24 绿。⇒ **保留它＝⛔ 因为它有牙，而因为一枚"在四个偏移里只红一次"的哨兵⛔ 该被删掉**（删了就等于给下一个改这里的人留一枚可以悄悄放宽的口子），入库用例里它头顶注释要按我这把写、⛔ 照 `a1r` 那句写。
- ⚠**它"无影响、⛔ 打折"那一判我接受结论，但它抛给我的两枚我一枚拆掉、一枚认下来**：
  - ★**"下一位照 `293-v1` 逐字配方重跑得到 313 而⛔ 677＝复现性缺陷"这一句我判它⛔**（尺＝它自己 `30` 件里那行命令逐字 `go test ./internal/panel/ ./internal/audio/ ./internal/ball/ -count=1 -v -timeout 240s` ⇒ **⛔ `cmd/wisp` 那一档**，而 677 那对是**四枚包并集**，两把⛔ 同形）。它同一页下一行其实已经写着"缺的 364 只能由 `cmd/wisp` 那一档补齐"⇒ **这不是复现失败，是它拿一枚没跑的档去比别人的并集**。⇒ ⛔ 记它纪律、⛔ 退回它任何东西；后果＝**"313↔677 复现不上"这句话今后⛔ 许任何人引用**（包括我），要坐实 677 只有一条路＝**我自己照逐字配方建一棵新导出树跑 `cmd/wisp` 那一档**（〔仅本机可量，归编排者，欠我一发；⛔ 与在飞的 `300-r1` 并发跑整包〕）。
  - ★**"那三枚 dll 怎么进 `293-v1` 那棵导出树的，它⛔ 能说明 ⇒ 归我核"这一枚我认**，且它替我把一条线索量出来了：它 `logs/cmdwisp-exporttree.txt` 逐字印着 **"module-cache lib dir: onnxruntime.dll / sherpa-onnx-c-api.dll / sherpa-onnx-cxx-api.dll"** ⇒ **三枚 dll 在模块缓存里本来就躺着** ⇒ "没有 dll 就载入死"那一形并非只有 `cp` 一条解法。⇒ `A802` 那条"导出树先数 dll"的前置尺**仍然对**（数＝最便宜的自证），但**成因我此前写窄了**（我只写了"未跟踪 ⇒ 目录不存在 ⇒ 载入死"，没量"还有第二枚来源"）。⛔ 改 `A802` 原句，本条只做指针。
- ★**它顺手替我量掉的一枚我⛔ 有的坑（⛔ 本票的账）**：派单里我举的例子"`internal/audio` 带 `sherpa` 字样那条"**实测 0 命中**（尺＝`grep -rln 'sherpa\|onnxruntime' internal/audio` ⇒ 0），⇒ **`internal/audio` 与那三枚 DLL 无关、那条 harness PATH 对它 inert**（`300-a1r` 的 `cp` 是多余的保险，⛔ 污染它的读数）。⇒ 并回＝**派单里举"哪枚包依赖哪枚 dll"这种话，落笔前先 `grep -rln` 一发**（与 `A803` 那条"路径先 `ls`"同一枚毛病的第三个面）。
- **排程（今）**：在飞＝**1 枚** `300-r1`（锚＝我 13:0x 的 `HEAD`；写面＝`internal/audio/**`＋`probes/300/r1/**`；三件事＝偏移 26⇒**24** ＋同批把 `// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].` 那三行改成实话（★我给它的是**我这把**：只有 **两枚数**坏＝`channelMask` 22⇒20、`SubFormat` 26⇒24，而"validBits@18"偏移对、名字该是 union `Samples` 的成员——⛔ 采它上面那句"三个数全是坏的"）＋把台件转成**入库用例**（`//go:build windows`、断解析出来的值、一进门就该绿、S4 保留按我这把写注释）。它交回前⛔ 派第二枚走 Go 编译面的腿、⛔ 我跑那发 `cmd/wisp` 复现；⛔ 任何人动 `internal/audio`。之后＝`300-v2`（非实现者裁 `AC#2`）→ `AC#3`（与票 297 关系，＋我那发真机 `GetMixFormat`）→ `AC#5`（我已量完，只等它把凭据写进交付）。队列其余＝`296-r2` → `298-r1` → `295-r1` → `294-r1` → `294-v1`；票 293 只剩 `AC#5`（合并真机窗口）；票 299 按在 `Q-84`（⛔ 不催）。⛔ 零 push。

## 编排者现量（2026-10-10 13:5x，`AC#3` 里那句〔仅本机可量、归编排者〕的 `GetMixFormat` dump；⛔ 翻框、⛔ 不改上面任何一节原句）

件＝`.scratch/wisp/probes/300/orch/2026-10-10-ac3-getmixformat-hexdump.md`（跑＝编排者本人，⛔ 腿；载具**在仓外导出树**、⛔ 入库）。

- **结论**：这台机器上**四枚**默认端点（render/capture × console/multimedia）的混音格式**都是** `WAVE_FORMAT_EXTENSIBLE`(65534)，
  `cbSize=22`、2 声道 48000 Hz 32 bit，`dwChannelMask=3`（前置左右），而 **`SubFormat` 的 16 字节 GUID 逐字节等于**
  `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` `{00000003-0000-0010-8000-00AA00389B71}`。⇒ 票面那句"到底是不是 extensible＋float"＝**是**。
- **那 2 字节之差现在是真事**：偏移 **24** 处是 `03 00 00 00`（＝FLOAT 的 `Data1`），偏移 **26** 处是 `00 00`
  ⇒ 旧写法读出的 `tag=0` ⇒ `floating=false`。**⛔ 止于 `Data1`**：这次连整枚 GUID 都对上了。
- **尺与来路**（逐字在件里）：`go test ./internal/audio/ -run TestOrch300MixFormatHexdump -v -timeout 180s` ⇒ `rc=0`／`--- PASS`；
  跑的是哪一版码**⛔ 靠"应该是新的"**——`cmp` 导出树 `wasapi_windows.go` ↔ `git show HEAD:internal/audio/wasapi_windows.go` ＝ **IDENTICAL**，
  两边 `:211` 都是 `unsafe.Add(p, 24)` ⇒ 这发落在**已改到 24 的那一版**上（`028529fc` 起）。起终两枚 `date` 都在件里（13:53:37→13:53:38）。
- **射程**：**只读字节面**——⛔ `Initialize`、⛔ `Start`、⛔ 打开麦克风、⛔ 出声、零窗口、零弹窗（机主屏幕上⛔ 变化）。
  现量 `Get-Process wisp`＝0／`balldebug`＝0；CPU 37%→39%、MEM 54.7%→56%（跑前跑后各一发）。
- **⛔ 这发证明的**（写死，免得下一程读大）：⛔ 采集链路能跑通／球进干活进程（票 228 那一族）；⛔ CI 现在会跑这一族用例（票 301 `AC#5` 未闭合，`A804`）；
  ⛔ 别的机器／驱动／采样率同形——**只这一台、只这一刻、四枚端点**。
- ★**与 `A805` 那发仓外五枚偏移扫描对拉**＝两把⛔ 同一把尺（一把合成夹具：24 四形全绿／26 四形全红；一把真机字节面：24 处是真 GUID、26 处是 0）
  ⇒ `feedback-evidence-teeth` 第 8 条那一形（一把尺两种形状都放行）被两把独立尺同时排除。
- ⚠**本格⛔ 翻**：`AC#3` 有半枚是"判语归位"（297 的 ③ 拿到形状，⛔ 把 297 的"是什么"判成已定案），那一半归 `300-v2` 之后的非实现者裁；
  本节只把**编排者欠的那半枚读数**交上盘。
- ⚠**载具的两枚边界具名留着**（转成常驻仪器前必须先修，见件第 1 节末尾）：① 它读**固定 40 字节**、⛔ 先看 `cbSize`（本机安全，普通 `WAVEFORMATEX` 端点上就是越界读；生产那一支先判 `tag==extensible && cbSize>=22`，⇒ **这是载具的边界⛔ 生产的**）；
  ② 件第四节是从终端回显转写的，第四枚端点那节我误排过一行，**凭据是那张解码表＋四组 hexdump 逐字节相同**，⛔ 那块代码块冒充"原始件"（要纯原始件就照件里那条命令重跑并 `>` 到文件）。
- ★**我自己踩的三形已具名进件**（⛔ 抹）：slot 判给生产写法（前两形按 slot 3 调 `GetDefaultAudioEndpoint` ⇒ `0x80070057`／`S_OK` 而指针为空，那正是 `EnumAudioEndpoints` 的 `stateMask` 语义；本仓 `mmdevice_windows.go:85` 逐字 `comCall(enum, 4, …)` 且注释写着 `(slot 4)`）；
  写文件工具把载具规范化到 `C:\tmp\...`（那棵树连 `go.mod` 都没有 ⇒ 第二形**从未被编译**，谈不上读数）；打印循环 `b[i:i+16]` 未夹长度会 panic。

## 编排者收 `300-v2`（2026-10-10 14:2x）：`AC#2` 翻勾＋★**我自己那两句假话被这枚腿顶回来**（原句⛔ 删，只追加本节）

- **收**：四笔 `69b3bbb8`→`21c9abe3`→`2328934e`→`2d51a1b2`，件＝`probes/300/v2/**`，逐笔越界 0 枚；卫生＝`rc-build=0`／`rc-d22scan=0`／`gofmt` 名册空／emoji 0 枚／它的件 0 字节＝0 枚。
  双锚照 `A807` 那条规矩：起手 `41329475`（porcelain `-- internal cmd`＝**0 行**）／交回 `2d51a1b2`；期间 HEAD 漂过 `849ce6e9`（`301-a2`）与 `29681406`（我这发 AC#3 读数），它自己加了一把尺 `git diff --name-only 41329475..HEAD -- internal cmd`＝**0 行** ⇒ **产码面⛔ 漂移**。★那把"跨锚证产码没动"的尺是它自加的，记它，并回定式层。
- ★★**它顶回我两句话，两句都写在盘上而证据⛔ 支持**：
  ① 我 13:0x 那一节第 97 行写着"S4 在 **20/22/24/28** 全都绿、只在 26 红"＝**错**。真值＝S4 绿 **22／24／28**、红 **20／26**；20 那一档绿的是 **S1**。
     ⚠ 最难看的一处＝**凭据就在我自己那份按节表里**（`.scratch/wisp/probes/300/orch/2026-10-10-orch-five-offset-scan.txt` 第 20 行逐字 `--- FAIL: .../S4_data1_shifted_to_26`、第 112 行 `offset=20 pass=1 fail=3`），我落笔那句时⛔ 拿那张表对过我自己那句话。
     ⇒ 定式并回：**"我自己写的按节表"与"我自己写的结论句"必须逐枚对拉再过一遍**——枚数对⛔ 等于名字对（这是 `A805` 那条"逐节计数把射程钉在当节内"的**下一层**）。三处各自追加更正（那份件的新节／本节／台账 `A810`），⛔ 删原句。
  ② 我据此喂给落地腿的注释口径（"四形里最弱的一枚哨兵"）方向⛔ 变（S4 放行面确实最宽＝最弱），但同段那句"**没有一枚单独是决定性的**"也⛔ 成立：**S2／S3 各自只绿 24＝单形即已钉死偏移**。
     ★ 而落地那枚用例头顶的注释当时写的是**真值**（`stays green at offsets 22, 24 and 28`）⇒ **腿⛔ 照抄我的假话**，这一点记它。
- **判语＝`AC#2` 成立**（四件事各有尺：偏移配对扫描五档 `24→4/0`、`26→0/4`、`20/22/28→1/3`；重写那三行注释的逐句可证伪性由它现读 `mmreg.h` 证；入库用例一进门就绿；配对基线两发 `432/320`→`437/321`、`comm -13/-23/-3` 三把全空 ⇒ **新增红 0**）⇒ 本格翻勾（凭据＝非实现者判语 ＋ 我自己复跑的按节 tally 与四形 PASS）。
- ★**新归因（它量、我认，⛔ 我复跑——尺是它的，结论进账）**：**"基线 5 枚红"⛔ 全是仓／机既有红**：`TestCleanCheckoutBuilds_AC11` 在真实工作树里 `rc=0`／`PASS 26.10s`，只在**⛔ 有 `.git` 的导出树**里红（`git rev-parse` 退 128＝**装置红**）⇒ 正确写法＝"**4 枚真窗族＋1 枚装置红**"。这是 `A802` 那条"导出树⛔ 自带 handicap"的**第二面**：⛔ 只缺 DLL，**还缺 git 元数据**。
  ⇒ 定式：**跑带 `git` 断言的包之前，先问那句断言在⛔ `.git` 的树里能不能绿。**
- **两枚交付形状缺陷**（具名，⛔ 阻断 `AC#2`）：`probes/300/r1/logs/post-commit-recheck.txt`（817 B）**从未入库**，而 `40` 件 §1／§2 两处写着"逐字＝该件"⇒ 一枚 **blob 层拿不到的悬空引用**（尺＝`git ls-files --error-unmatch` ⇒ `Did you forget to 'git add'?`；读数本身在 `.md` 里有，所以那格⛔ 空）；`30-gates.md:116` 残留未填的 `<!--ROSTER-->`。⇒ 并回＝**文中指名的件与注释里指名的测试同级，一律当待验断言**（`A` 账第五形的那一支）。
- **纪律两枚，都⛔ 豁免、都⛔ 用它退回读数**：① `300-r1` 那笔 `git add` ＋ `git reset -- 自己刚 stage 的路径` ⇒ 记**腿**一枚具名纪律缺陷（`300-v2` 三条理由我全采：禁的是**形状**／这把尺事后⛔ 可造＝⛔ 可证伪／改契约⛔ 归腿），正解它自己验通（`git check-attr text eol -- <未跟踪 .go>` ⇒ `text: set`／`eol: lf`、`rc=0`）；② `300-v2` 自报在仓内造了两枚探针件然后 `rm`（2 发删除）⇒ 记**腿**，理由＝`issues/README` 规则 8 原文"临时件**只建不删**"，⛔ 按"这条的理由在我这儿不成立"自行放宽。**撤销口令「撤 300 两枚纪律缺陷」**。
- ⛔ **翻的两格说清楚为什么不翻**：`AC#3` 还剩"判语归位"那一半（297 三支的 ③ 已有可复现形状＋我这发真机读数，但⛔ 把 297 的"是什么"判成定案）⇒ **另派一枚非实现者**（⛔ 我自己裁，因为那半枚是**判语**）；`AC#5` ⛔ 翻——本格逐字指名"**落地腿**自己现跑那三跳"，而实际把三跳量到的是**验收腿**（`windows)` 清单 audio 0 命中／`cli)`＝只有 `./cmd/wisp/`／core 清单有 audio 而该档 `ubuntu-latest`／用例 `:1` 是 `//go:build windows` ⇒ 凭据确实只有本机）；**指名的人⛔ 跑**＝那一格⛔ 闭合，它的自然闭合点＝票 301 `AC#1` 落地腿的前后名册作差（分母按 **10 枚用例**，见 `A809`）。
- **代笔一笔**（编排者本人，⛔ 动产码语义／⛔ 动产码偏移／⛔ 动断言与夹具字节）：那枚入库用例的注释与 `why:` 里三处说明（"⛔ 单形决定性"那句／`mmreg.h:2474`⇒`:2475`／绝对 SDK 路径⇒`shared/mmreg.h`）——件＝`.scratch/wisp/probes/300/orch/2026-10-10-comment-fix-gates.md`，七把门禁各落 `rc=`：`gofumpt`／`gofmt` 名册空、改动区非 ASCII 命中 0、`go vet` rc=0、定向用例 rc=0（尺＝顶层 `--- PASS` 1 枚＋`=== RUN` 5 枚＝1 父＋4 子，⛔ 一把合计数）、`d22scan` rc=0、`GOFLAGS= go build ./...` rc=0。
  ⚠ **记我**：第一发我看见本机 PATH ⛔ `gofumpt` 就把那格当"跑⛔ 了"，第二发才改用 `$(go env GOPATH)/bin/gofumpt` 取到读数——**命令找⛔ 到⛔ 等于检查通过，先定位再判**。
- ⛔ 零 push（机主从未授权）。**排程（改）**：`301-r1`（audio 进 `windows)` ＋同批补 `win_pin`，分母按 10）**现在可派**——写面 `internal/audio/**` 已⛔ 有人在飞、Go 编译面空 ⇒ 我按完这笔账就派；其后 `301-v1` → `300-v3`（只裁 `AC#3` 判语归位那一半）→ 票 111 `AC#12` 的普查腿。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-10 12:32:39 +08] agent=300-a1r did=AC#0 权威偏移读定＝SubFormat 在 24（凭据＝盘上 Windows SDK shared/mmreg.h，第 34 行 pshpack1.h 令 Packing=1 ⇒ 偏移表可算；cbSize=22 与总长 40 由 mmreg.h:2540/2550 与 Learn 那句 "must be at least 22" 从"编排者推的"升为"读到的"；union 那一问＝不是，Samples 与 dwChannelMask 是两枚独立成员，标准结构里没有 wValidSamples；据此 AC#0 结论把 wasapi_windows.go 判成错的那一枚、wavinjector 同形，硬门"若是 26 才停手"未触发）＋AC#1 仓外成对树（git archive b2933dc9 到 /tmp/wisp300-a1r/tree，自造 40 字节面 4 形喂 parseWaveFormat，断解析出的 tag 值）：改前 4/4 FAIL（S1/S2/S3 恒读 0，S4 那枚 Data1 挪到 26 的坏面读到 3 ⇒ 这枚尺两个方向都敏感）、只把树上 unsafe.Add(p, 26) 改 24 后同一份用例 4/4 PASS、随后树已 cmp 验过逐字节还原；恒红用例⛔ 入库（件只落 probes/300/a1r/ 的 00/10/20/90 与 logs/ 七枚 .txt）；三把卫生＝d22scan rc=0、GOFLAGS= go build rc=0、git status --porcelain -- internal cmd 0 行 next=编排者裁 AC#0/AC#1（⛔ 我翻框、⛔ 我动 AC#2）；给 AC#2 两条输入：偏移改 24 时头顶注释里 validBits@18/channelMask@22/SubFormat@26 三个数全是坏的、建议入库用例保留 S4 那枚反向形；AC#3 真机 GetMixFormat 与 AC#5 CI 三跳本程⛔ 跑，读数仍欠
- [2026-10-10 12:38:01 +08] agent=编排者 did=收 300-a1r（四笔 70a6b9d9→efd85155、四枚 .md 逐枚 find -printf 现量 4,612／12,509／17,598／11,103、logs 8 枚其中 0 字节＝0；⚠ 一处过程痕迹＝hygiene-build.txt 在 1c007bf6 那笔是 0 字节入库、下一笔才填 191 字节，按"0 字节的件＝那格没交"那一笔当时不合格，它自己补对了）⇒ **翻 AC#0／AC#1 两格**，⛔ AC#2..AC#5 一枚不碰（AC#2 开工权在非实现者 300-v1 手上，票面那句"⛔ 只有 AC#1 红过并经非实现者裁之后才动"我⛔ 让它静默 dissolve）；★**承重两发我自己重跑，⛔ 引用腿的数**：① 我自己开那枚头文件 `C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`（176,680 字节）逐字读 `sed -n '2521,2534p'` 的 WAVEFORMATEXTENSIBLE ＋ `:31-36` 的 pshpack1 ＋ `grep -n "cbSize = 22"` 命中 :2540／:2550 ⇒ 我自己算的偏移＝WAVEFORMATEX 18、Samples@18、dwChannelMask@20、**SubFormat@24**、总长 40 ⇒ **判 24 对、wasapi_windows.go 那枚 26 错**；② 我在腿留的导出树 /tmp/wisp300-a1r/tree 先核那行逐字仍是 unsafe.Add(p, 26)，再跑两形＝**原样 4/4 FAIL**（红句逐字：S1 tag = 0, want 3／S2 tag = 0, want 1／S3 tag = 0, want 7／S4 tag = 3, want 0 ＋ floating 两行）／**只把树上 26 改成 24 ⇒ 4/4 PASS**，还原证明＝同一条命令内 cp 备份→跑→写回→`cmp` ⇒ restore_cmp=IDENTICAL 且 `grep -c 'unsafe.Add(p, 26)'` 回到 1 ⇒ ★这把尺两形配对成立、⛔ 恒真三形之一；由这两发现在能说的最硬一句＝extensible 混音格式下 tag 恒 0 ⇒ floating=false ⇒ convertPacket 走 int16 那条 unsafe.Slice（⛔ 等于"票 297 那 0.3677 就是它"，那一句仍只有真机那一对能量）；★**夹具⛔ 循环**：我逐字节读了 build()（b0=0xFFFE／b16=22／b18=bits 写进 Samples／b20=mask／Data1 写 data1At，S1-S3＝24、S4＝26）⇒ 与头文件逐格对得上，⛔ 拿仓里任一枚实现当布局裁判；★**接受腿两处顶我、只追加⛔ 改原句**：票面第 13 行那句"nValidSamples 与 dwChannelMask 是同一个 union 的成员"⛔＝union 的名字是 Samples（三枚 WORD：wValidBitsPerSample／wSamplesPerBlock／wReserved），dwChannelMask 是独立成员，头文件里⛔ wValidSamples 这个名字（结论层不变＝偏移仍 24）；另我派单里那句示例路径 Include/*/um/mmreg.h 在盘上⛔ 成立（真身 shared/）⇒ 定式并回：派单出现盘上路径，落笔前先 ls/find 一发（与 A797"声明 HEAD 对象层射程前先 git show"同族）；★**新事实改写我 1 小时前刚立的定式（A802 追加更正 A801）**：`git ls-files third_party/sherpa-onnx | wc -l`＝**0** 而盘上有三枚 DLL ⇒ `git archive` 导出树里那条 harness PATH 指向**不存在的目录**，"成对导出树"法⛔ 自带这处 handicap；300-a1r 是 cp 三枚进树才让尺成立 ⇒ 今后用导出树跑带 harness 的包前先 `ls <树>/third_party/sherpa-onnx/*.dll | wc -l` 并具名；⚠ 连带一枚我⛔ 现在下结论的问题＝293-v1 的件里⛔ 提过 DLL（我 `grep -in "dll" .scratch/wisp/probes/293/v1/20-ac4-gates.md`＝0 命中），它那两组 677／684 要不要打折 ⇒ 交 300-v1 三档裁（无影响／射程受限／要打折）；③ 腿建议"AC#2 入库版保留 S4 那一形"＝裁法交 300-v1，但我先在 AC#2 派单钉⛔ 删 S4；next＝在飞 1 枚 `300-v1`（锚 efd85155，攻★ⓑ／★ⓒ 这套凭据＋裁 AC#2 开工权＋答那枚打折问题）；它交回前⛔ 派 300-r1、⛔ 任何人改 internal/audio；之后 300-r1（偏移改 24＋同批把 validBits@18, channelMask@22, SubFormat@26 那三行注释改成实话——那行注释今天逐字教下一个改这里的人相信 26，比代码危险）→ 300-v2 → AC#3（与票 297 关系归位，⛔ 合并结案）→ AC#5（CI 那三跳）；队列其余 296-r2 → 298-r1 → 295-r1 → 294-r1 → 294-v1；票 299 按在 Q-84（⛔ 不催）。⛔ 零 push、⛔ 本票零产码
- [2026-10-10 12:55:06 +08] agent=编排者 did=把 AC#5 那三跳预跑了一遍并写进本票"编排者现量"一节（⛔ 翻框、⛔ 不改上面各节原句；全只读、⛔ 编译面 0 发、⛔ 与在飞的 300-v1 抢仪器）：四跳逐跳读到那一行＝ci.yml:403/470 → portable-tests.sh:234 的 core 清单 :241 确有 ./internal/audio/...（core_pin :168 也确有该导入路径）而 core 那档 runs-on ubuntu-latest／ci.yml:604 → wisp-cli-tests.sh:113 → portable-tests.sh:256 的 cli 档 scope 只有 ./cmd/wisp/（:260）／ci.yml:780 的 --scope=windows 清单 :248-253 十枚路径逐枚读过、⛔ audio／★第四跳 297-a1 没报＝ci.yml:744 --scope=census 的判定单位是**导入路径**（:382-397），internal/audio 已被 core 认领 ⇒ 永远进不了 unclaimed ⇒ ★通用新知＝本仓"有没有人跑这枚用例"那把门只到包级、⛔ 到用例级，任何包里 windows-tagged 的判据天生不在 CI 执行面而⛔ 任何一枚门会为此变红；ubuntu 那档⛔ 红的原因也量了＝GUARD A 数的是"编进本平台二进制"的文件（:565-566、:364），尺（限定每枚文件前 3 行找 go:build windows）＝internal/audio 七枚 _test.go 里 2 枚 tagged（capturelevel_windows_test.go、hotplug_test.go★文件名不带 _windows 而内有 tag，按文件名分类会读错）余 5 枚 linux 编得进 ⇒ 三门全绿而这枚新用例从未被执行。⇒ 本票 AC#2 那批的凭据只有本机；⛔ "为进 CI 把用例搬到 cmd/wisp"这条路物理上不存在（parseWaveFormat 包内未导出）；要买回可见性＝把 ./internal/audio/ 加进 windows) 清单并同一笔更新 win_pin（:193），代价⛔ 我全量（该包其余 windows-tagged 用例一起进分母，含 :593 已登记的 TestLiveWasapiSmoke）⇒ 归口一枚独立仪器票（票 255 那一族），⛔ 塞进本票 AC#2。通用账并回台账 A804。next=⛔ 派 300-r1（等 300-v1 四条必答，含"AC#5 我这条预跑成不成"）；我自己那发 GetMixFormat hexdump（AC#3 欠的、〔仅本机可量、归编排者〕）本要现在跑，量完前提后按住＝它要新写一枚 Go 载具＝走 Go 编译面＝撞 A804 那条按住 ⇒ 欠账记我、⛔ 记任何腿，等 300-v1 交完第一件事就跑
- [2026-10-10 12:56:21 +08] agent=300-v1 did=非实现者验收腿三格＋④裁完（锚 efd85155；⛔ 盖章⛔ 修产码⛔ 翻框）。① AC#0 **成立**：偏移＝**24** 由我自己 sed 抽的盘上原文闭合（mmreg.h 的 WAVEFORMATEX 7 字段＝18、Samples union 三成员逐字皆 WORD＝2、dwChannelMask DWORD、SubFormat GUID 16 ⇒ 18+2+4+16＝40；mmreg.h:2540/2550 两处逐字 cbSize = 22；26 那一形会让 cbSize＝24、总长＝42 ⇒ 与原文直接矛盾）；guiddef.h 我自己开（Data1 unsigned long ⇒ 16 字节）。①b **union 那一争＝腿对、票面 :13 后半句错**（Samples 与 dwChannelMask 是两枚独立成员，grep -c ValidSamples mmreg.h ⇒ 0 ⇒ 标准结构⛔ nValidSamples 这枚名字）。② AC#1 台件**⛔ 循环、有牙**：夹具 10 处写入逐处落在我自算的权威表上（⛔ 一处用 @18/@22/@26）；我在看读数**之前**用权威＋夹具手算出锚上四形读 0/0/0/3，我重跑实测逐字命中（4/4 FAIL、rc-gotest=1、无 0xc0000135）；我另把偏移扫成五形 **20/22/24/26/28 ⇒ 全绿只在 24**（26 全红）⇒ 判据⛔ 是"两形都放行"的尺；还原自证 cmp 逐字节同、grep -c 26＝1、24＝0、20/22/28＝0。③ S4 **要保留，但⛔ 因为腿给的那条理由**：S1+S2+S3 子集在 20/22/26/28 已经全红 ⇒ 少了 S4⛔ 变单向尺；我给的两条独立理由＝S4 是唯一让错位读取出"看似合法的 FLOAT(3)"的形＋正确偏移下它是绿（成本零）；同时⚠ **S4 是四枚里最弱的一枚（放行 22/24/28 三形）**，S1 单形在 20 也绿 ⇒ 判别力实际扛在 S2/S3，入库⛔ 只留一枚。**AC#2 ＝ 开**（四条输入进件 §4）。④ 三档结论＝**无影响**：我自己两形对照证明那枚坑是真的且是载入级（无 dll 的 git archive 树里 go test ./cmd/wisp/ -run 空 ⇒ exit status 0xc0000135、零用例名、rc=1；cp 三枚 dll 后同一条命令 ok/rc=0），但 293-v1 的 677/684 **⛔ 打折**——我把它 POST 红名册抽的 8 枚逐名定位 **8/8 全在 cmd/wisp/** ⇒ 它的树当时 DLL 齐、cmd/wisp 真跑了；第二把尺同向（TestAC293 本体在 cmd/wisp 而它交回 PASS=7），第三把尺同向（我 3 包实测 313 vs 它 677，缺的 364 只能由 cmd/wisp 补）。记两笔具名缺陷：293-v1 ⛔ 命名这枚前置条件 ⇒ 照它逐字配方重跑会得到 313 而⛔ 677（复现性缺陷，⛔ 读数错误）；**我⛔ 能说明它那三枚 dll 怎么进树的，那一笔归编排者核**。next=编排者翻 AC#0/AC#1 框＋派 AC#2 落地腿（带 §4 四条）＋自己核 293 那两棵树的 dll 来路＋处置票面 :13 那句；AC#3 真机 GetMixFormat 读数仍欠（"这枚缺陷今天响不响"仍未定，⛔ 本票结案）；三把卫生 rc-d22scan=0／rc-build=0／git status --porcelain -- internal cmd＝0 行；件＝probes/300/v1 的 00/10/20/30 四枚＋logs 十一枚 .txt
- [2026-10-10 13:12:04 +0800] agent=编排者 did=收 300-v1（五笔 25a36ccb→3da1a1a1；五枚 .md 逐枚现量 5,657／10,287／11,716／14,327／2,924＋logs 17 枚、0 字节=0 枚；逐笔 git diff-tree 越界 0 枚）⇒ AC#2 开工权放行并派 300-r1（写面 internal/audio/** ＋ probes/300/r1/**；三件事＝偏移 26→24、同批把那三行布局注释改成实话、把仓外台件转成入库用例且一进门就该绿）；★四条必答里我坐实两枚、拆掉一枚、认下一枚：②AC#1 ⛔ 循环（凭据另有我自己在仓外导出树重跑的五枚偏移扫描＝只有 24 四形全绿、26 四形全红、20/22/28 各只绿一枚，还原 cmp IDENTICAL；件在 probes/300/orch/2026-10-10-orch-five-offset-scan.txt，★里面还留着我踩的两枚坑＝第一发 sed 射程打宽到整文件那枚台件无效、第二发 grep -c 打在整份日志＝累计值⛔ 一把尺）③S4 保留我采 v1 的理由而⛔ 采 a1r 的（我这把现量 S4 在 20/22/24/28 全绿、只在 26 红＝四形里最弱的一枚哨兵，入库注释按我这把写）④打折问题结论我接受（无影响）而它那句"照 293-v1 逐字配方重跑得到 313 而⛔ 677＝复现性缺陷"我判它⛔：它自己跑的命令⛔ cmd/wisp 那一档（go test ./internal/panel/ ./internal/audio/ ./internal/ball/），677 是四枚包并集，两把⛔ 同形；⇒ 今后引用任何人"复现不上"之前先把两条命令逐字对拉射程（本仓这族第四次） next=dll 来源那一枚归我核（它量出模块缓存 lib 目录里本来就躺着三枚 dll ⇒ A802 那条前置尺的成因我写窄了，数 dll 仍要做）；我欠自己两发（照 293-v1 逐字配方含 cmd/wisp 重跑坐实 677、真机 GetMixFormat 前 40 字节 hexdump），都因在飞的 300-r1 按住；⛔ 翻任何框（AC#2..AC#5 仍 4 枚未勾）；★过程留痕一条记我：本节 14 行我 13:0x 写完没立刻提交，而同一时刻 300-r1 正在飞（锚 ce06cbe6）＝正是 A804 我给自己立的那条"编排者＝一枚隐形写腿，改完跟踪文件当场就交"的头一次实弹应用，这一笔按那条补交；此刻工作树里 300-r1 已有未提交改动＝wasapi_windows.go 偏移 26→24 ＋ 那段布局注释整段重写，我逐字读过、⛔ 动、⛔ 判（判归 300-v2），③那枚入库用例此刻还没出现；⛔ 零 push
- [2026-10-10 13:39:31 +08] agent=300-r1 did=AC#2 落地（实现者腿，锚 ec5cdc51，起手三读＋进程前置尺见 probes/300/r1/00-anchor.md）：ⓐ internal/audio/wasapi_windows.go 那枚 sub := *(*windows.GUID)(unsafe.Add(p, 26)) 改成 24（行号随 ⓑ 从 196 移到 211，落地尺 grep -n "unsafe.Add(p, " 四枚读数只第 4 枚变），ⓑ 同批把头顶三行布局注释整块改写（tag@0/channels@2/rate@4/bits@14/cbSize@16 五枚保留、validBits@18 改写成 union Samples@18 并说破"validBits 不是字段名"、channelMask@22⇒20、SubFormat@26⇒24，另写总长 40／cbSize=22／出处 shared/mmreg.h，⛔ 现在时绿测试计数；按编排者那把＝坏掉的只有两枚数，⛔ 采纳"三个数全坏"那把，这一处分歧具名交 300-v2），ⓒ 入库 internal/audio/parse_wave_format_300_windows_test.go（//go:build windows、包内 audio、⛔ 真设备/真麦/Initialize/Start，断的是解析出来的 tag/floating/bits 值⛔ 断"被调用过"，四形 S1..S4 原样沿用 a1r 夹具：剥掉行尾注释与 a1r 那枚台件 diff ⇒ 十处 le.PutUint 写入＋四枚 face300 字面量＋断言逐字节同形，唯一语义差＝S4 的 why 按 v1 §6.2 加"最弱一枚"短语；S4 在 24 上断的是 tag==0 且 floating==false，锚上读 26 它拿 3＝凭据 a1r/v1/orch 三把）⇒ 一进门就绿（工作树两发 rc=0、pass-all=5/pass-top=1、load0xc=0；入库后 treeB 用最终字节再跑一发同形）；门禁七把全交 rc-build=0、rc-d22scan=0、格式两把（工作树＋仓外 HEAD blob＋仓外导出树副本）三处空列表、裸 gofumpt 具名不在 PATH、成对导出树各两发：改前 432/320 红5、改后 437/321 红5，新增红=0、消失红=0、两侧两发红名册对称差皆空（基线那 5 枚 8/8 定位在 cmd/wisp，含本机稳定红的 TestAC4Focus…；抖项 Test197… 四发全绿），三数带尺名；dll 前置尺两棵树各量：cp 前 0（⚠ 那枚目录在导出树里根本不存在的坑：cp 报 is not a directory，第一发基线两发因此 LOAD 死 0xc0000135＝作废，补 mkdir -p 后重跑），改后树用 cmp 证与工作树逐字节同；AC#5 照编排者那把写死"本枚用例进不了 CI 任何一档、凭据只有本机"，⛔ 我重跑四跳、⛔ 为迁就门禁搬落点，该可见性已立成独立仪器票 301；自报一枚命令形状违规＝为取未跟踪件的 ls-files --eol 读数跑了 git add + git reset -- <同一路径>（后果零、⛔ 碰别人名册，正解是 git check-attr），另自报两把我踩坏的尺（grep -c $'' 嵌在命令替换里成空模式＝数了全部行；台件"无设备"那把尺把注释散文算成命中）＋全仓格式存量在派单点名三枚之外还有四枚 i/lf w/crlf 的 internal 生产件（pending_read/tools/provenance/bridge，具名作废⛔ 我动）；⛔ 翻任何框、⛔ 改票面原句、⛔ push、⛔ 动 cmd/wisp 与冻结面。next=300-v2（非实现者裁 AC#2：偏移/注释/入库用例/门禁作差四问；我给它的三条待裁＝① 我与 a1r 关于"注释坏了几枚数"的分歧 ② S4 保留但已标注最弱、以及我⛔ 补第五形（v1 §4.2 的 plain PCM 面）的取舍 ③ 我那一枚 reset 形状违规）；AC#3 真机 GetMixFormat 那发仍欠（归编排者、⛔ 腿），我这把⛔ 重新裁决"该不该修"
- [2026-10-10 13:59:01 +0800] agent=编排者 did=交票面 AC#3 那半枚〔仅本机可量、归编排者〕的 GetMixFormat dump（件 probes/300/orch/2026-10-10-ac3-getmixformat-hexdump.md，载具在仓外导出树、⛔ 入库）＝这台机器四枚默认端点（render/capture × console/multimedia）混音格式**全是** WAVE_FORMAT_EXTENSIBLE(65534)＋cbSize=22＋2ch/48000/32bit＋mask=3，且 SubFormat 那 16 字节 GUID **逐字节**等于 KSDATAFORMAT_SUBTYPE_IEEE_FLOAT {00000003-0000-0010-8000-00AA00389B71}；偏移 24 处是 03 00 00 00、偏移 26 处是 00 00 ⇒ 旧写法读出 tag=0 ⇒ floating=false ⇒ AC#3 那句"到底是不是 extensible＋float"＝**是**（⛔ 止于 Data1，这次连整枚 GUID 都对上）；尺与来路＝go test -run TestOrch300MixFormatHexdump -v ⇒ rc=0／--- PASS，跑的是哪一版码⛔ 靠"应该是新的"（cmp 导出树 wasapi_windows.go ↔ git show HEAD: 同一枚＝IDENTICAL、两边 :211 都是 unsafe.Add(p, 24) ⇒ 落在已改到 24 那一版），起终两枚 date＝13:53:37→13:53:38；★与 A805 那发仓外五枚偏移扫描对拉＝两把⛔ 同一把尺（夹具那把：24 全绿／26 全红；真机这把：24 处是真 GUID／26 处是 0）⇒ feedback-evidence-teeth 第 8 条那一形被两把独立尺同时排除；⛔ 翻框（AC#3 还有"判语归位"那一半，归 300-v2 之后的非实现者）；⛔ 这发证明的四件事写死在票面新节里（⛔ 采集链路跑通、⛔ CI 会跑这一族、⛔ 别的机器同形、⛔ 整条路就对了）；⚠ 载具两枚边界具名留着（读固定 40 字节⛔ 先看 cbSize＝载具的边界⛔ 生产的；第四节是从终端回显转写、第四枚端点我误排过一行 ⇒ 凭据是那张解码表＋四组 hexdump 逐字节相同，⛔ 那块代码块冒充"原始件"）；★我自己踩的三形进件（slot 判给生产 mmdevice_windows.go:85 的 comCall(enum, 4, ...)——前两形按 slot 3 调拿到 0x80070057／S_OK 而指针为空，那正是 EnumAudioEndpoints 的 stateMask 语义；写文件工具把载具规范化到 C:\tmp 而那棵树连 go.mod 都没有 ⇒ 第二形从未被编译、谈不上读数；b[i:i+16] 未夹长度会 panic）；⛔ Initialize／⛔ Start／⛔ 开麦／⛔ 出声／零窗零弹窗，现量 wisp=0 balldebug=0、CPU 37→39%、MEM 54.7→56% next=① 等 300-v2（判 AC#2，含换形也红那把反恒真尺、⛔ reset 族那一笔独立裁、logs/vet.txt 0 字节那格怎么判）；② 票 301：301-a2 已交、待我复跑它顶回的四条再裁 AC#3 落／不落；③ 我还欠自己一发＝照 293-v1 逐字配方（含 cmd/wisp 那档）在仓外导出树重跑坐实 677 ⛔ 零 push
- [2026-10-10 14:22:49 +0800] agent=编排者 did=收 300-v2（四笔 69b3bbb8→2d51a1b2；件 probes/300/v2/**；逐笔越界 0；卫生 rc-build=0／rc-d22scan=0／gofmt 名册空／emoji 0；双锚起手 41329475 porcelain -- internal cmd＝0 行、交回 2d51a1b2，并自加一把"跨锚证产码没动"的尺＝git diff --name-only 41329475..HEAD -- internal cmd 0 行，记它）⇒ ★AC#2 判语成立、本格翻勾（凭据＝非实现者判语＋我自己复跑的按节 tally 24→4/0、26→0/4、20/22/28→1/3 ＋ 定向用例四形 PASS）；★★它顶回我两句话且两句都⛔ 站得住：① 我 13:0x 第 97 行"S4 在 20/22/24/28 全都绿"是⛔ 的（真值＝S4 绿 22/24/28、红 20/26，20 那档绿的是 S1），而凭据就在我自己那份按节件第 20 行 --- FAIL .../S4 与第 112 行 pass=1 fail=3 里 ⇒ 定式并回＝我自己写的"按节表"和我写的"结论句"必须逐枚对拉再过一遍（枚数对⛔ 等于名字对，A805 那条的下一层）；三处各自追加更正⛔ 删原句（orch 件新节＋本节＋A810）② 同段"没有一枚单独是决定性的"也⛔ 成立——S2／S3 单形即只绿 24＝各自有牙；★落地那枚注释当时写的是真值（22/24/28）、⛔ 照抄我的假话，记它；★新归因＝"基线 5 枚红"⛔ 全是既有红：TestCleanCheckoutBuilds_AC11 在真实工作树 rc=0／PASS 26.10s、只在⛔ 有 .git 的导出树红（git rev-parse 退 128＝装置红）⇒ 写法换成"4 真窗族＋1 装置红"，这是 A802"导出树⛔ 自带 handicap"的第二面（⛔ 只缺 DLL 还缺 git 元数据），定式＝跑带 git 断言的包之前先问那句断言在⛔ .git 的树里能不能绿；两枚交付形状缺陷具名（post-commit-recheck.txt 817B 从未入库而 40 件两处写"逐字＝该件"＝blob 层拿不到的悬空引用，尺＝git ls-files --error-unmatch；30-gates.md:116 残留未填 <!--ROSTER-->）⇒ 并回＝文中指名的件与注释里指名的测试同级、一律当待验断言；纪律两枚都⛔ 豁免都⛔ 据此退读数（300-r1 的 git add＋git reset -- 自己刚 stage 的路径＝记腿，v2 三条理由我全采、正解 git check-attr 它自己验通；300-v2 自报在仓内造两件探针然后 rm＝记腿，规则 8 原文"只建不删"⛔ 按理由自行放宽；撤销口令「撤 300 两枚纪律缺陷」）；⛔ 翻 AC#3（还剩"判语归位"那一半，是判语⇒另派非实现者⛔ 我）与 AC#5（本格逐字指名"落地腿自跑那三跳"，实际量到的是验收腿，指名人⛔ 跑＝⛔ 闭合，自然闭合点＝票 301 AC#1 落地腿前后名册作差、分母按 10 枚见 A809）；编排者代笔一笔（⛔ 动产码偏移⛔ 动断言⛔ 动夹具字节：注释里"⛔ 单形决定性"那句／mmreg.h:2474⇒:2475／绝对 SDK 路径⇒shared/mmreg.h），件 probes/300/orch/2026-10-10-comment-fix-gates.md 七把各落 rc：gofumpt 名册空 rc=0、gofmt 名册空 rc=0、改动区非 ASCII 命中 0、go vet rc=0、定向用例 rc=0（顶层 --- PASS 1 枚＋=== RUN 5 枚＝1 父＋4 子，⛔ 一把合计数）、d22scan rc=0、GOFLAGS= go build rc=0；⚠ 记我＝第一发见 PATH ⛔ gofumpt 就当那格跑⛔ 了，第二发才用 $(go env GOPATH)/bin/gofumpt 取到读数——命令找⛔ 到⛔ 等于检查通过；⛔ 零 push next=写面 internal/audio 已空、Go 编译面已空 ⇒ 派 301-r1（audio 进 windows) 档＋同笔补 win_pin，AC#2 前后名册作差分母按 10）

## 编排者收 `300-v3`（2026-10-10 17:3x，两笔 `f76091cd`→`9fc5c410`，件 `.scratch/wisp/probes/300/v3/**`，逐笔名册越界 **0**，⛔ push）⇒ **翻 `AC#3`（成立，带两枚限定）**＋★它量到一枚**盘上⛔ 任一把尺钉得住的常量**⇒ 本票追加 `AC#6`＋认三处顶回（含我自己那句"CI 兜着"已过期）

### `AC#3` 的判语（我复跑过的才落，⛔ 复跑的照口径标注）

- **成立**："票 297 那一支拿到一枚**可复现的形状**"。凭据＝tracked 用例 `internal/audio/parse_wave_format_300_windows_test.go`（尺＝`git ls-files internal/audio | grep 300`，我现跑命中）＋腿那把**仓外**突变；红句逐字 `tag = 0, want 3 (…at byte 24)` **四形齐**，把 24 写回 26 那一发 `rc=1／FAIL=5`。
- **成立但只到本机此刻**："本机那份混音格式确实是 extensible＋float"。凭据＝**我那一发** `GetMixFormat` dump（`probes/300/orch/2026-10-10-ac3-getmixformat-hexdump.md`，四枚端点同一份 40 字节、24 处 `03 00 00 00`、26 处 `00 00`，件末 `rc=0`）——本格把这一枚明写成〔仅本机可量、归编排者〕，⛔ 腿⛔ 我任何一方可重复。
- **⛔ 成立**："票 297 那个 0.3677 的底噪就是这枚偏移造成的"。★我复跑并追认：`grep -c '^- \[ \]'` 票 297＝**5 枚全未勾**、`^- [x]`＝**0**；`ls .scratch/wisp/probes/297/`＝**只有 `a1/`**（⛔ 一枚 `go test` 件）；票 297 的 `AC#0` 那枚"真设备判别仪器"**盘上⛔ 存在**。⇒ 本票⛔ 许把 297 的"是什么"判成定案（原句如此，这一格照办）。

### ★顶回①：支号错了，而且错得会把人带到另一枚落点——归我收

票面 `AC#3` 写"票 297 三支里的 **③（"切错字节"）**"。我 `sed -n '24,30p'` 现读票 297 的权威清单：**②**＝"字节切错（`convertPacket` 那一层的 frames/channels/bits 拼装）"，**③**＝"设备本底就是这样"。⇒ **两处同名不同号**。
- ⚠ 更要紧的是**落点也不同层**：票 297 的 ② 指的是 `internal/audio/wasapi_windows.go:385-403` 那层**拼装**；本票的偏移在 `parseWaveFormat`（读 `SubFormat` 的位置），后果是 `tag` 读错 ⇒ `floating` 判错 ⇒ `convertPacket` 走错路。**"同一族、不同层"**，所以本票只算给"字节切错"那一族拿到可复现形状，⛔ 判 297 的 ② 成立。
- 处置＝票 297 一字不改（⛔ 我改它的清单），**由本节把"③"读成"②那一族的名字"**；今后任何腿引用两票关系一律**按标签引、⛔ 按号引**（同族＝行号锚腐烂的兄弟：编号锚同样会腐烂）。

### ★顶回②：我那句"CI ⛔ 兜着／一进门就绿"已过期（`A817`/票 301 落地之后）

票面 `AC#5` 引的是 `297-a1` 那句"`./internal/audio/` 只在 core(ubuntu)、⛔ 在 windows/cli"。**现状已变**＝票 301 的落地笔 `0a0f62ef` 已把 `./internal/audio/` 写进 `windows)` 档并在同笔补了 `win_pin`（我 15:5x 复跑过那笔名册）。⇒ `AC#5` 那格**⛔ 翻**（它问的是"CI 到底求值过这枚 tagged 用例没有"，那一发色**仍欠**，只有 CI 给），但票面凡"audio 的 windows 半边⛔ 会被任何档求值"这句**从今天起按快照读**。

### ★顶回③＋新格：那枚常量今天**没有任何一把尺钉得住**

`gofmt -l cmd internal scripts tools` 名册 5 枚（含 `cmd/wisp/models.go`）＝⛔ 一枚在腿的射程、⛔ 顺手格式化；`go vet ./internal/audio/` 的 stdout 是 0 字节而 `rc=0`，两件事⛔ 混写。
★恒真攻击的**收获**（这条比翻勾值钱）：三侧反形全红（夹具侧／期望侧／产码语义侧各一发，每发先 `LANDS` 自证、跑完 `cp`＋`cmp` 逐字节还原＝`restore_cmp=IDENTICAL`），**但**腿一发 `N2` 把 `waveFormatFloat` 由 3 改成 1 ⇒ **整包 41 枚顶层全绿 `rc=0`**。
⇒ 我用读码复核了它为什么会绿（⛔ 复跑，⛔ 需要复跑）：`internal/audio/wasapi_windows.go:378` 是 `floating := s.format.tag == waveFormatFloat`，而用例里 `parse_wave_format_300_windows_test.go:152-153` 两侧都拿同一个常量比（`gotFloat := got.tag == waveFormatFloat` / `wantFloat := tc.wantTag == waveFormatFloat`）⇒ **改常量同时挪两侧，恒等式照成立** ⇒ 那行 `floating` 断言⛔ 能独立红。全仓 `grep -rn 'floating'` 只命中"生产那行＋注释＋用例那两行"（我现跑），**零枚**尺钉"这个常量⛔ 是 3"。
⇒ 本票因此追加 **`AC#6`**（不是"顺手修"，是一枚**新判据**，⛔ 由实现腿自勾）：

- [ ] **`AC#6`（2026-10-10 17:3x 编排者追加，来路＝`300-v3` 的恒真攻击 `N2`）给 `waveFormatFloat` 那枚常量装一枚钉**：判据＝**把常量换成任一其它值（3→1 与 3→0 两形）指名用例必须红**，且红句具名指向那枚常量而⛔ 指向 `tag` 的等值比较；凭据形状＝成对两发（改⇒红／`cmp` 逐字节还原⇒绿）。**硬约束三条**：① ⛔ 动 SLO 阈值／golden／`thresholds.go`；② ⛔ 为了让它红而**放宽任何既有的 `tag` 断言**（那条是本票 `AC#1`/`AC#2` 的凭据本体）；③ 新判据⛔ 能写成"常量 == 3"的字面等值再包一层同义反复——判据＝**同一枚换形攻击跑在新尺上必须红**（先例＝`A710` 那枚"模板照抄会造恒绿假钉"，⛔ 匹配形状而⛔ 匹配语义）。落点⛔ 预设：`internal/audio` 的这枚常量该⛔ 由既有 `parse_wave_format_300_windows_test.go` 加一段，还是该钉在别处（"格式常量与权威表一致"那一族），由落地前的只读一发裁。

### 本票状态与 next

**4 勾（`AC#0`／`AC#1`／`AC#2`／`AC#3`）／3 未勾（`AC#4` 门禁＋越界、`AC#5` CI 色欠、`AC#6` 新格）**，⛔ 改 `-done`。
- Progress log：
  - [2026-10-10 17:3x +0800] agent=编排者 did=收 `300-v3`（两笔 `f76091cd`→`9fc5c410`，越界 0）⇒ **翻 `AC#3`**（成立，两枚限定：真机那一枚〔仅本机可量、归编排者〕；"0.3677 就是它"**⛔ 成立**，凭据＝票 297 五格全未勾＋`probes/297/` 只有 `a1/`，我两把尺现跑对上）＋★**追加 `AC#6`**（`waveFormatFloat` 那枚常量全仓零尺钉得住：改 3→1 整包 41 枚照绿，我用读码复核出原因＝用例两侧同用该常量，恒等式⛔ 能独立红）＋认三处顶回（支号 ③ 应为 297 的 ② 那一族、且**落点不同层**；我那句"CI ⛔ 兜着"在 `0a0f62ef` 之后过期；`gofmt` 名册 5 枚⛔ 一枚在腿射程）；⚠ 枚数口径那格按票面⛔ 现在拍、⛔ 在此合并两形；next＝`AC#6` 的落点只读一发 → `AC#4`/`AC#5` 的 CI 色（跟票 303 修复后那批推送一起取，⛔ 为它单推）

## 编排者收 `300-a2`（2026-10-10 17:5x，六笔 `ae5ee86d`→`fbeefd1d`，件 `.scratch/wisp/probes/300/a2/**`；越界＝我自己把六笔名册取并集重算＝**0 枚**；⛔ push）⇒ **裁 `AC#6` 落点＝C2 形**（新枚同包 `windows`-tagged 用例；⛔ C1／⛔ C3 落盘／⛔ C4 当裁判／⛔ C5）＋★补一枚⛔ 在腿射程里的权威读数（本机 `mmreg.h:2110` 逐字＝`0x0003`）＋两把词面尺降级

**先说这发的质量与两处口径**：五格承重事实我逐枚在 **HEAD blob 层**复跑并打中（下面 ① 段），它自报的三把自己的尺错（射程打宽吞进别腿两笔／`echo` 标签混进作差管道／反引号坑复发）全部**只追加⛔ 抹**，正合票 301 `AC#4b` 那三件正向闸门。⚠ 两处枚数口径要具名写：① 它交件说 `logs/` **12** 枚，我数 **11**（尺＝`find .scratch/wisp/probes/300/a2/logs -type f | wc -l`；顶层件 5／`snap/` 4）⇒ **按尺报数**这件事我这也要照自己办；② `F4` 那句"injector 的 `fmtTag==3`／`0xFFFE` 两支今天**零覆盖**"**是词面尺**（`grep -i float|FFFE|extensib`）——我另补一把同族词面尺（`float32|Float32` 在 `wavinjector_test.go` ＝ **0 命中** ⇒ 连造 float32 载荷的写法都没有，比它那把强，⛔ 覆盖尺）。覆盖要 `go test -coverprofile`，那归落地腿，⛔ 现在取（`303-r1` 在飞）。

### ① 我复跑并打中的（⛔ 引它的话）

- **F1 常量射程**：`git grep -n waveFormatFloat HEAD -- '*.go'` ⇒ 全仓 **4 枚**＝`internal/audio/wasapi_windows.go:98`（声明 `= 3`）、`:378`（`floating := s.format.tag == waveFormatFloat`，在 `Drain()` 体内）、`internal/audio/parse_wave_format_300_windows_test.go:152`＋`:153`（恒等式两行）。⇒ **硬后果成立**：任何直接引用该常量的判据⛔ 可能不带 `windows` tag（linux 档连编译面都进不去）。
- **F2**：`waveFormatPCM` 只命中声明行 `:97`（**声明后零使用**）；`waveFormatExt` ＝声明 `:99` ＋唯一使用 `:210`。
- **F3**：`convertPacket` 生产调用点**只有 `:400` 一枚**，测试侧 **0**（另两枚命中是 `:149`／`:155` 的注释与红句文案）⇒ "生产真的这么判"今天除 `Drain()`（要活流）外无法观察。
- **F5 枚数**：逐文件 blob `grep -c '^func Test'` ＝ 1/4/5/8/10/1/7/6 ⇒ **42 枚顶层**，`func TestMain` **0** 命中。⇒ **票面那句"整包 41 枚顶层全绿"与它的 42 是两把尺**（41＝`--- PASS` 那把运行时尺；42＝源码 `^func Test` 尺）⛔ 矛盾，引用必带尺名。
- **N5（对我这格最要紧的一枚）**：`git grep -n -e wasapi_windows -e parse_wave_format HEAD -- cmd/wisp/config_receipt_255_test.go` ⇒ **0 命中** ⇒ 255 那族名册钉⛔ 按行号引 audio 面 ⇒ 落地腿新增用例⛔ 会把它们打红。

### ② ★我补的那把尺在腿的射程之外（件＝`.scratch/wisp/probes/300/orch/2026-10-10-mmreg-authority.txt`，17:58:13 现量，可逐字重跑；⚠ 那把尺读的是**仓外**文件，补位只读腿⛔ 授权射程＝由我代跑并具名）

本机 `C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h` 现量三行逐字＝

- `:2110  #define  WAVE_FORMAT_IEEE_FLOAT                 0x0003 /* Microsoft Corporation */`
- `:2376  #define  WAVE_FORMAT_EXTENSIBLE                 0xFFFE /* Microsoft */`
- `:2418  #define WAVE_FORMAT_PCM         1`

⇒ 三个常量**同一份权威文本**都能对上；⚠ 且权威那一行写的是 **`0x0003` 而⛔ `3`**——新用例注释引源时要按那一行逐字写，⛔ 把它"顺手化简"成 `3`（那等于又造一次同源）。
另⚠ 一处术语更正（腿在 `C3` 里混用了）：`KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` 那枚 **GUID** 在 `ksmedia.h`，⛔ 在 `mmreg.h`；`mmreg.h` 管的是**那枚 WORD 值**（`0x0003`）。本票 `AC#0` 当初定 `mmreg.h` 是为**偏移表**（`WAVEFORMATEX` 的字段布局），这格要钉的是**值**——两者同源⛔ 同物。

### ③ 裁语（五段，理由与判据逐段带）

- **`AC#6` 落点＝C2**：新枚同包 `//go:build windows` 的 `_test.go`（例名 `wave_format_float_300_windows_test.go`，真名由腿定且⛔ 与既有 8 枚同名），**⛔ 改任何既有件**。四条硬要求——
  ① **期望侧⛔ 得经该常量**（经了就＝`:152-153` 那枚恒等式的复活，那是本格的病灶）；权威字面按上面 `:2110` 那一行写成 `0x0003` 形，并在注释里**逐字记录出处**（路径＋行号＋那一行原文）。
  ② 红句**具名点名 `waveFormatFloat`**、⛔ 只写 `tag` 的等值比较——★**这一半是我派单漏抄的**，`AC#6` 原文逐字有这句，腿按原文办是对的（顶回①成立）。
  ③ 顺带把 `waveFormatExt`、`waveFormatPCM` 一起钉（同一份权威文本里就有对应两行，边际成本≈0）。⚠ **⛔ 顺手删 `waveFormatPCM`**——"声明后零使用"是词面尺读数，删它是另一件事、属另一格授权。
  ④ 凭据＝**成对两发**（3→1、3→0 各一发：改⇒红，句里点名那枚常量；`cmp` 逐字节还原⇒绿），⛔ 只报枚数。⚠ 新增那枚**顶层用例**会把 audio 多推一枚进 `test-windows` 的红名册作差面（`0a0f62ef` 之后 audio 已在 `windows)` 档，⚠ 那句的世代锚见票面 `AC#5` 那节），落地腿要按 `AC#4` 那把名册尺重跑并**具名说多出来的是哪一枚**。
- **C1 ⛔**：它与 `AC#1`/`AC#2` 的凭据**共用同一枚 tracked 用例文件**（那两格逐字引 `:152-153`）⇒ 改它会把两格的 prose 锚与那对成对基线一起推动；且 C1 钉的是"浮点判定⇔解析值的一致性"，⛔ 那枚常量本身的值。
- **C3 ⛔ 落盘读头文件；★但它的"权威"以另一形状买回来**：最硬那条代价是**结构性**的——无头文件时若走 `t.Skip` ⇒ 未登记 SKIP＝红（`scripts/portable-tests.sh` 的 class 定义逐字 `fixture = the machine or OS cannot supply the subject at all`，尺＝blob `:586-600`，我读过），要登记就⛔ 改 `scripts/**` ⇒ **撞 `AC#4` 名册**（本票只许 `internal/audio/**`＋`probes/300/**`）；另两条＝要改掉那枚 "the absolute path … is deliberately not written out here"（`git show HEAD:internal/audio/parse_wave_format_300_windows_test.go | sed -n '17,21p'` 逐字，形式变更⛔ 由腿自取）＋托管镜像带⛔ 带 SDK 我今天答⛔ 了。⇒ **处置＝"值⇔权威"的证据进探针件、⛔ 进 CI 分母**：落地腿交一件 `probes/300/r2/logs/authority-mmreg.txt`（那把 `grep -n -E 'define +WAVE_FORMAT_(IEEE_FLOAT|PCM|EXTENSIBLE)' <头文件>` 自己落一行 `rc=N`）；**我收件时自己再跑一遍那把尺**才算这格闭。⚠ 这形⛔ " weaker 版 C3"——它⛔ 让 CI 变红的能力来自测试自带的权威字面，盘上那一发只是**出处的证据**，两件职责⛔ 混。
- **C4 ⛔ 当本格的裁判**（★正面回答腿顶回来的那枚问题）：`AC#0` 那句"⛔ 拿本仓另一枚实现当裁判"，**我裁＝对"值"同样成立**。理由＝两枚仓内实现可以**一起错**，而这正是当初造出这枚禁令的成因（偏移那形与值那形同因；放开＝用一枚⛔ 独立的事实去钉另一枚）。另两条独立理由也指⛔：它顺带买到的是**新增覆盖面**（`AC#6` 要⛔ 它），且 `parseWav` 的 `size < 40` 守卫会被拉进判据射程。⇒ injector 那两支的覆盖缺口**具名降成残余 R1**，⛔ 混进本格。
- **C5 ⛔**：要动 `wasapi_windows.go` 的**产码语义面**，而本票只授权过一次产码改动＝`AC#2` 的偏移＋注释那一发。`Drain()` 那行唯一消费者的可测性降成残余 R2。

### ④ 残余（具名，⛔ 塞进 `AC#6`）与欠读数

- **R1** injector 的 `fmtTag==3`／`0xFFFE` 两支的**覆盖面**要用覆盖尺定（今天只有两把词面尺：`float|FFFE|extensib` 2 命中全是 `sineI16At` 那枚 `float64` 生成器；`float32|Float32` 0 命中）。
- **R2** `wasapi_windows.go:378` 那行要不要抽成包内谓词（＝C5）——唯一能覆盖"生产真的这么判"的形，代价＝改产码语义面，⛔ 本票授权。
- **R3** 托管镜像带⛔ 带 Windows SDK 头文件＝**本裁语下⛔ 需要**（C3 既裁⛔ 读盘）；若日后要一枚"CI 级外部权威钉"，那要另开一票、且必须先答"未登记 SKIP＝红"那条结构性代价。
- 腿具名的欠读数逐条回答：**①** 成对两发颜色＝落地腿（本腿⛔ 编译面，正确）；**②** CI 那一发色＝编排者，⛔ 为它单推（跟票 303 修复后那批一起取）；**③** SDK 在否＝⛔ 需要（见 R3）；**④** audio 在 `windows` 档当前基线色＝引票 301 那批读数时**必须带那一发的锚点**（`A819`/`A820` 世代）；**⑤** `TestLiveWasapiSmoke` 在托管是 SKIP 还是红＝与②同批；**⑥** 255 族仪器会否看 audio 字节面＝**我已答**（① 段 N5：0 命中）。

### ⑤ 排程（⛔ 现在派落地腿，具名理由）

`300-a2` 交完 ⇒ 本票状态＝**4 勾（`AC#0`/`AC#1`/`AC#2`/`AC#3`）／3 未勾（`AC#4` 门禁＋越界、`AC#5` CI 色、`AC#6` 待落地腿）**，⛔ 改 `-done`。
⛔ 现在就派 `300-r2`：它要 `go test`（`internal/audio`）并新增一枚 `_test.go`，而 **`303-r1` 正在 `cmd/wisp` 写产码＋跑真窗**（既有定式＝一枚 `cmd/wisp` 写腿在飞时整个导入图⛔ 动源码；先例 `A710`/票 111 那节）。next＝`303-r1` 交完并由我复跑过 → **`300-r2`（`AC#6` 按上面四条硬要求落地）** → `300-v4`（非实现者裁 `AC#6`/`AC#4`，⛔ 落地腿自勾）→ `AC#5` 的 CI 色与票 303／票 111 那批**同一枚推送**一起取。

- [2026-10-10 17:5x +0800] agent=编排者 did=收 `300-a2`（六笔 `ae5ee86d`→`fbeefd1d`，越界我自己取并集重算＝0）⇒ **裁 `AC#6` 落点＝C2**（新枚 `windows`-tagged 用例，期望侧⛔ 经该常量、按 `mmreg.h:2110` 那行的 `0x0003` 形写权威字面、红句具名点名常量〔我派单漏抄那半句，腿顶回对〕、顺带钉 `waveFormatExt`/`waveFormatPCM` 但⛔ 删它、成对两发）；**⛔ C1**（与 `AC#1`/`AC#2` 凭据共用同一枚文件且钉⛔ 值）／**⛔ C3 落盘**（未登记 SKIP＝红 ⇒ 必改 `scripts/**`＝撞 `AC#4` 名册；改为"权威证据进探针件、⛔ 进 CI"，收件时我自己再跑那把 grep）／**⛔ C4 当裁判**（★我裁"⛔ 拿仓内另一枚实现当裁判"**对"值"同样成立**，两枚实现可以一起错；它的覆盖缺口降成残余 R1）／**⛔ C5**（改产码语义面）；★补一枚⛔ 在腿射程里的读数（本机 `mmreg.h` 三行逐字 `0x0003`／`0xFFFE`／`1`）＋一处术语更正（`KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` 那枚 GUID 在 `ksmedia.h`、⛔ `mmreg.h`）；两把词面尺降级（"零覆盖"⛔ 覆盖尺；logs 枚数它报 12 我数 **11**）；排程＝`300-r2` **按住等 `303-r1`**（`cmd/wisp` 写腿在飞时整个导入图⛔ 动源码）→ `300-v4` → `AC#5` 的 CI 色与票 303／111 同批推送取。


## 300-r2 交回（落地腿，2026-10-10 18:5x；⛔ 翻框、⛔ 改上面任何原句；本节只追加）

**本格＝`AC#6`，落点＝编排者收 `300-a2` 裁语 ③ 定的 C2 形。`AC#6` 的判语归非实现者 `300-v4`，本腿⛔ 自勾任何一枚 `- [ ]`。**

起手锚 `05db4bc6`（porcelain `-- internal cmd docs .scratch/wisp/issues`＝**0 行**）；本腿四笔＝`a0339995`→`9a442923`→`cdced117`→（本节那一笔）→（终名册那一笔）。
件＝`.scratch/wisp/probes/300/r2/`（`00-anchor`／`10-shape`／`20-mutations`／`30-gates`／`40-final` ＋ `logs/` **14** 枚，尺＝`find .scratch/wisp/probes/300/r2/logs -type f | wc -l`，⛔ 0 字节＝本腿自己现量）。⛔ push。

- ★**交付形状**：新增一枚 `internal/audio/wave_format_float_300_windows_test.go`（`//go:build windows`）。
  **期望侧三枚字面量按 `mmreg.h` 那一行逐字写：`0x0003`／`0xFFFE`／`1`**（⛔ 化简成 `3`），期望侧⛔ 出现 `waveFormatFloat`／`waveFormatExt`／`waveFormatPCM` 任何一枚；
  观察侧＝那三枚常量本身 ＋ `parseWaveFormat` 解析出来的 `tag`。⇒ `:152-153` 那枚"两侧同枚常量"的恒等式**没有**被复制。
- ★**成对两发突变都红、且红句点名常量**（⛔ 指向 `tag` 的等值比较那一形只出现在钉 `waveFormatExt` 的两枚子测试里）：
  `waveFormatFloat` 3→1 ⇒ `rc=1`、三句 `Errorf`（`waveFormatFloat = 0x0001, want 0x0003; mmreg.h:2110 …`／`waveFormatPCM and waveFormatFloat both = 0x0001 …`／`… mmreg.h reserves 0x0001 for WAVE_FORMAT_PCM …`）；
  3→0 ⇒ `rc=1`、两句（`waveFormatFloat = 0x0000, want 0x0003 …`／`… which is none of WAVE_FORMAT_PCM 1, WAVE_FORMAT_IEEE_FLOAT 0x0003 or WAVE_FORMAT_EXTENSIBLE 0xFFFE …`）；
  两发各自 `cp` 备份→跑→写回→`cmp` ⇒ **IDENTICAL**，声明行回到 `98:	waveFormatFloat = 3`，复跑 ⇒ `rc=0`。逐字读数＝`20-mutations.md` ＋ `logs/mut-m0..m3`。
  ⚠ 台面＝**仓外导出树**（`git archive HEAD`，⛔ clone、⛔ 在共享工作树里做任何 checkout／改产码），四处两两 `cmp` 证与工作树逐字节同；⛔ 拿"一发母仓、一发 clone"作差。
- ★**八把门禁**：① `go vet ./internal/audio/` `rc=0`（输出 0 字节，两件事⛔ 混写）② 定向用例 pristine `rc=0`（顶层 `--- PASS` 1／子测试 6／`=== RUN` 7）
  ③ 整包两发同一枚台面（改前＝本文件还⛔ 在树里那发）⇒ 顶层红名册尺 `grep -E '^--- (FAIL|SKIP): '` 两发**逐字同形**，`comm -13`／`-23`／`-3` **三把全空 ⇒ 新增红 0**；三数带尺名：顶层 `--- PASS` 41→42、含子测试 6→12、`=== RUN` 48→55
  ④ 格式三把并排（工作树／HEAD blob／锚 `05db4bc6` blob，射程目录都写 `internal/audio`）名册各 **0 枚**；全仓残留 **5 枚**（`cmd/wisp/models.go` 等）⛔ 一枚在本腿射程、⛔ 顺手修、只具名
  ⑤ `sh scripts/d22scan.sh` `rc=0`，输出自证 ban #8 实扫 `internal/` **526** 枚 Go 文件且"comments and `_test.go` included" ⇒ 本用例在仪器射程内、⛔ 命中（红句里⛔ U+26D4／U+2713／U+2264／U+1F000–1FAFF／U+2200–22FF／U+2600–27BF／U+2B00–2BFF／`U+FE0F`／U+1F1E6–1F1FF；整枚新文件纯 ASCII，尺＝`grep -cP '[^\x00-\x7F]'` ⇒ 0）
  ⑥ `bash scripts/portable-tests.sh --scope=census` totals 行**逐字未变**＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`；⚠ 同族一枚读数＝`internal/audio` 那行从 `8/0` 变 **`9/0`**（本腿那枚文件确实编进本平台测试二进制，分母⛔ 动）
  ⑦ 名册差集＝**逐笔** `git show --name-only --format=`，三笔全部落在"允许动的写面"名册内 ⇒ 越界 **∅**（明细见 `40-final.md`）
  ⑧ `internal/audio` 顶层用例枚数（尺＝逐文件 blob `grep -c '^func Test'`）：**42 → 43**；⚠ 口径具名＝**新增一枚顶层用例** `TestWaveFormatConstantsMatchMmregAuthority300` ＋ 它内部 6 枚子测试，⛔ 是"只加子测试"。裁语 ④ 要具名的那枚就是它（`0a0f62ef` 之后 audio 在 `windows)` 档，`test-windows` 红名册作差面多的那一枚⛔ 同这个名字）。
- ⚠**`--- SKIP` 一枚的归因写死**：改前那发盘上**已有** `--- SKIP: TestLiveWasapiSmoke`（`scripts/portable-tests.sh` 夹具台账里已登记的 fixture 例外），两发名册逐字同形 ⇒ 既有、已登记、**⛔ 本腿造的、也⛔ 本腿让它消失的**。本腿⛔ 用 `t.Skip`、⛔ 把 `t.Fatalf` 换成 `t.Skip`。
- ⚠**一枚我造过又拆掉的假仪器（具名留痕，⛔ 抹）**：`waveFormatPCM/Float/Ext` 当 **map 字面量的键**时，3→1 那一发⛔ 是红而**是 `duplicate key 1 in map literal` 编译错误**（`rc=1` 而零句 `--- FAIL`）⇒ 我用一枚仓外最小件证了这条语言行为，交付里换成运行期切片＋两两比较。**突变攻击要的读数是一句具名红，⛔ 是一个非零 rc。**
- ⚠**观察侧到不了生产那一行（具名，⛔ 造假日为面）**：`waveFormatFloat` 唯一消费者在 `(*wasapiStream).Drain()` 体内（要活流＝真设备，票 300 票面第 31 行⛔ 归腿），`convertPacket` 收的是 `floating bool` **参数** ⇒ 残余＝裁语里的 **R2**（走 C5 才买得到＝改产码语义面，本票⛔ 授权）。injector 那两支的覆盖面＝**R1** 未动（本件⛔ 调 `parseWav`，⛔ 拿 C4 当本格裁判）。
- ⚠**⛔ 删 `waveFormatPCM`**：`F2` 那句"声明后零使用"是词面尺；本件按裁语 ③ 只钉它的值（它此前**零尺**），删它归另一格。

**★报回三处派单转述与票面原文的差别（一律按票面原文办）**

1. **裁语 C3 指定落地腿交一件 `probes/300/r2/logs/authority-mmreg.txt`**（票面原文："**落地腿交一件**…那把 `grep -n -E 'define +WAVE_FORMAT_(IEEE_FLOAT|PCM|EXTENSIBLE)' <头文件>` 自己落一行 `rc=N`"），而派单只转述了"本用例⛔ 读盘"。
   ⇒ 两者⛔ 矛盾（一件是**出处的证据**，一枚是**测试的运行时依赖**），但派单**没把这件列进我要交的清单**——我按票面原文补交了，`rc=0`，三行读数与编排者存档 `orch/2026-10-10-mmreg-authority.txt` 逐字同。⚠ 我⛔ 授权读仓外文件去⛔ 是这条给的：票面裁语 C3 那句"补位只读腿⛔ 授权射程＝由我代跑并具名"读起来像**编排者代跑**，与本格"落地腿交一件"两句**互相张力**——我按后一句（指名落地腿的那句）办了，这一处分歧交 `300-v4` 裁。
2. **派单硬要求第 4 条（顺带钉 `waveFormatExt`／`waveFormatPCM`）**⛔ 在 `AC#6` 原文（票面第 186 行那句"给 `waveFormatFloat` 那枚常量装一枚钉"）里，它出自**裁语 ③**。⇒ 我按裁语办了，并注明来源，⛔ 把它算进 `AC#6` 原文的要求。
3. **派单说"⛔ 与既有件同名，名册＝`300-a2` 的 `F5` 那 8 枚"**——票面裁语 ③ 原文是"真名由腿定且⛔ 与既有 8 枚同名"，一致；本腿名册尺已现跑（`10-shape.md` 第一节），8 枚既有名⛔ 撞。
4. ⚠**记我一枚过期读数**：`9a442923` 那笔的 commit 消息写着"子测试 PASS **5**"，实测是 **6** 枚（尺＝`grep -c '^    --- PASS: '`）。已推送的历史⛔ 改（⛔ `--amend`），以本节与 `30-gates.md` 门②为准。

**本腿没做到／票面没覆盖的（具名，⛔ 自己填）**

- `AC#4`／`AC#5` 两格⛔ 闭合：本腿⛔ 跑 `GOFLAGS= go build ./...`（⛔ 在我这八把清单里，⛔ 替 `AC#4` 交它），CI 那一发色仍欠（编排者 17:3x：与票 303／111 同批推送一起取）。⇒ **要不要现在补 `go build ./...` 那一发，我⛔ 自己假设，问一句。**
- 本腿⛔ 跑覆盖尺（`-coverprofile`）⇒ R1 那格"injector 两支到底有⛔ 有覆盖"仍只有词面尺。
- 本腿⛔ 动任何产码 ⇒ `waveFormatFloat` 在生产里"真的这么判"那一枚可观察性⛔ 买到（R2）。
- 本工作树**盘上有 834 行既有 dirt**（`design/**` 一批 ` D`、`.gitignore` ` M`、`.scratch/` 一堆未跟踪件，尺＝`git status --porcelain | wc -l`；⛔ 本腿造的，本腿只碰自己的 pathspec）——⚠ `internal`／`cmd` 两枚目录下 dirt＝**0 行**，⇒ 本腿的门③"同一枚台面"成立；但**这一枚事实我⛔ 能从派单里预知**，具名报回给编排者排程用。

## 编排者收 `300-r2`（2026-10-10 19:2x，七笔 `a0339995`→`9cef4589`，件 `probes/300/r2/**`；我自己的复验件＝`probes/300/orch/logs/r1-verify-300r2-20261010-192222.txt`，脚本＝`probes/300/orch/r1-verify-300r2.sh`）⇒ **`AC#6` ⛔ 由我翻**（判语按 17:5x 裁语归 `300-v4`），★并具名**推翻 `297-a1` 传给 `AC#5` 的那句前提**

### 0. 台面（⛔ 与腿同台面＝同一枚 HEAD 的**仓外导出树**，⛔ 在共享工作树里做任何 checkout／改产码）

- 起手尺＝`git status --porcelain -- internal cmd docs scripts .github` 现量 **0 行**；锚＝`HEAD_sha=9cef4589`（＝腿的终笔之后同一枚台面）。
- 导出＝`git archive 9cef4589 | tar -x` 到 `/tmp/300r2-verify-20261010-192222/repo`；`cmp` 该树里那枚用例 ↔ HEAD blob ＝ **IDENTICAL**。
- 本程⛔ 动任何产码；腿的写面我逐笔核过＝七笔并集只含**一枚新 `internal/audio/wave_format_float_300_windows_test.go`** ＋ `probes/300/r2/**` ＋票 300 本节上面那一节（⛔ 改任何既有件）。

### 1. 我自己现跑的读数（⛔ 抄腿的转述；每条带尺）

| 格 | 尺 | 我这发的读数 |
|---|---|---|
| 门① | `GOFLAGS=-mod=mod go vet ./internal/audio/` | `rc_vet=0` |
| 门② | 定向用例 pristine（`-v -run TestWaveFormatConstantsMatchMmregAuthority300`） | `rc=0`，顶层 `--- PASS` **1** ＋ 子测试 `    --- PASS` **6**（逐名：`ieee_float_low_word`／`extensible_tag`／`pcm_tag`／`the_three_constants_map_one_to_one_onto_the_authority_table`／`parseWaveFormat_expands_a_face_tagged_with_the_authority_extensible_value`／`control_plain_pcm_tag_is_not_expanded`）⇒ 与腿改正后的"6 枚"一致（`9a442923` 那笔 commit 消息写的 5 枚已被腿自己追加更正） |
| ★门③ | `GOFLAGS= go build ./...`（腿⛔ 跑、⛔ 在我 17:5x 派单清单里 ⇒ **漏列记我**，本程我销账） | `rc_build=0` |
| 突变一发 | `waveFormatFloat` 3→1（导出树里 `sed`，`cp` 备份→跑→写回→`cmp`） | `rc=1`；`declare_line_mutated=98: waveFormatFloat = 1`；**三句 `Errorf`** 逐字：`:191 waveFormatFloat = 0x0001, want 0x0003; mmreg.h:2110 …`／`:214 waveFormatPCM and waveFormatFloat both = 0x0001 …`／`:229 … mmreg.h reserves 0x0001 for WAVE_FORMAT_PCM (mmreg.h:2418) …`；红体＝顶层 `--- FAIL` ＋ 两枚子测试 FAIL；`restore_after_3to1=IDENTICAL` |
| 突变二发 | 3→0 | `rc=1`；两句 `Errorf`（`:191 … 0x0000, want 0x0003`／`:235 … which is none of WAVE_FORMAT_PCM 1, WAVE_FORMAT_IEEE_FLOAT 0x0003 or WAVE_FORMAT_EXTENSIBLE 0xFFFE …`）；`restore_after_3to0=IDENTICAL` |
| 还原腿 | 同把定向尺再跑一发 | `rc_restored=0` ⇒ **成对两发闭合**（改⇒红／`cmp` 逐字节还原⇒绿） |
| 整包两发 | `go test -count=1 -v ./internal/audio/`，改前＝把该用例从导出树里**挪走**（`mv`，⛔ 删），同一把尺 `grep -E '^--- (FAIL\|SKIP): '` | 改前名册**逐字**＝`--- SKIP: TestLiveWasapiSmoke (0.00s)`；改后名册＝**同一行**；`comm -13`／`comm -23` **双向 0 枚** ⇒ **新增红 0 枚**；两发皆 `ok … internal/audio`（16.058s ↔ 16.033s） |
| 枚数 | 尺＝整包 `grep '^func Test' \| wc -l`（射程＝导出树 `internal/audio/*_test.go`） | **43**（与腿那把"逐文件 blob `grep -c '^func Test'`"的 42→43 对得上） |
| 门④ | `gofmt -l internal/audio`（工作树一把；HEAD blob 那把腿已交） | **0 枚** |
| 门⑤ | `sh scripts/d22scan.sh`（跑在真仓＝工作树，全仓射程） | `rc=0`，自证行逐字含 `ban #8 internal/ examined 526 Go files, comments and _test.go included` ⇒ 新用例在仪器射程内而⛔ 命中 |
| 门⑥ | `bash scripts/portable-tests.sh --scope=census` totals 行 | 逐字未变＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` |

### 2. ★新现量一枚，直接打在 `AC#5` 的前提上（`297-a1` 那句「`internal/audio` 的 windows-tagged 用例整条 CI 无一档真跑」**⛔ 成立**）

三把独立尺（前两把我自己现读，第三把是归档 CI 字节）：
1. `scripts/portable-tests.sh:253` 逐字——`windows` 档的 `scope` 里**就写着** `./internal/audio/`（该行末元素；core 档在 `:242` 写的是 `./internal/audio/...`）；`:195` 的 `win_pin` 名册里逐字有 `github.com/CarlosShao/wisp/internal/audio` ⇒ **windows 那一档声称覆盖它**。⚠ `win_pin`/`core_pin` 是**包名名册**而⛔ 用例数名册 ⇒ 新增一枚用例文件⛔ 需要动 pin（这是"推下去会不会把 CI 顶红"那枚前提，我先量了才推）。
2. `test-windows` 那一步（job `114194107792`，改后发）里逐字出现 `ok  	github.com/CarlosShao/wisp/internal/audio	16.739s`，且同名 `--- PASS` 名册里有 `TestParseWaveFormatSubFormatOffset300`（＝`parse_wave_format_300_windows_test.go` 那枚，`_windows` 文件名层就只在 windows 建）与 `TestMonoDownmixAndFloatConvert` ⇒ **托管 windows runner 上确实求值了 `internal/audio` 的 windows 档用例**（`_windows_test.go` 那层文件名约束在那台机器上成立）。
3. 改前发（job `114187968428`）的 `packages=[…]` 行逐字同含 `./internal/audio/` ⇒ ⛔ 只在改后发偶然出现。
⇒ `AC#5` 的"这格会决定落点"那半，现在的料是**落点 C2 的 CI 可见性成立**（`windows` 档）。⛔ 我翻 `AC#5`：本格原文还欠"推送后那一发的**颜色**"——那枚新用例落在 `9a442923`，**晚于**我上面分析的那两发（`cf46c24a`／`05db4bc6`），所以它⛔ 进过 CI。本程推送之后才有那一发；已登记成待办（件＝票 303 的 `r9` 那两件同源推送）。

### 3. 判语归属：`AC#6` ⛔ 翻（我这一发＝**复验**，⛔ 终裁）

本格原文（`:186`）自己写了"**⛔ 由实现腿自勾**"，而我 17:5x 的裁语写的是"`AC#6` 的判语归非实现者 `300-v4`"。⇒ 上面第 1 节是编排者本人的独立复验（同台面、同判据、逐字红句对拉），**⛔ 替 `300-v4` 盖章**。交 `300-v4` 的必答题（派单会逐条带上）：
- ⓐ 这枚钉⛔ 是 `AC#6` ③ 禁的那个形状（"常量 == 3 的字面等值再包一层同义反复"）？我这发的观察侧确实读的是常量本身、期望侧⛔ 出现那三枚常量（`cmp` 过那三行 `mmregWaveFormat*` 声明）。
- ⓑ 权威凭据的**强度**：期望侧是"**抄自**仓外 `mmreg.h` 的字面量"＋三行逐字引文，测试**运行时⛔ 读那枚头文件**（腿在分歧 1 里把"出处的证据"与"测试的运行时依赖"分开了，我认它分得对；`logs/authority-mmreg.txt` 现量三行逐字我也复看过）。⛔ 读盘这一形在 `300-v4` 够⛔ 不够。
- ⓒ 顺带钉 `waveFormatExt`／`waveFormatPCM` 那两枚**⛔ 在 `AC#6` 原文里**（出自裁语 ③）⇒ 本格闭合的射程含⛔ 含它们、⛔ 删 `waveFormatPCM` 那一格归谁。
- ⓓ 残余 **R2**（观察侧到不了 `(*wasapiStream).Drain` 里 `floating := s.format.tag == waveFormatFloat` 那行）算⛔ 算本格必须闭的；闭它要形 C5＝动产码语义面，本票⛔ 授权。
- ⓔ 残余 **R1**（injector 那两支的覆盖尺）仍只有词面尺——腿⛔ 跑 `-coverprofile`。
- ⓕ **`AC#4` 的两枚具名欠**（腿自己在"本腿没做到"里报了 `go build`，但⛔ 报这两条）：(i) 票面 `AC#4` 原文要的是 `PATH=… sherpa … go test ./internal/audio/ **./cmd/wisp/** -count=1 -v` **改前改后各 ≥2 发**取红名交集——腿那把门③只跑了 `./internal/audio/` 且各**一发**；(ii) "三数必带尺名"那一半腿做到了，我照它的口径对拉过。`cmd/wisp` 那半要不要补、补给谁＝交 `300-v4` 裁（⚠ 那半要真窗台面与 sherpa `PATH`，⛔ 顺手算成腿的欠账——先例＝我 10-08 那条"派单里就禁掉的资源⛔ 算腿的欠账，记编排者"）。

### 4. 现态与 next

- **现态⛔ 变**＝仍 **4 勾／3 未勾**（`AC#0`/`AC#1`/`AC#2`/`AC#3` 勾；`AC#4`/`AC#5`/`AC#6` 未勾），⛔ 改 `-done`。
- next＝**推送（`A813`，fast-forward）** ⇒ 取"那枚新用例进 CI 之后"的颜色（销 `AC#5` 那半）→ **`300-v4`**（非实现者一次裁 `AC#6`＋`AC#4`，必答题＝上面 ⓐ..ⓕ）→ 票 303 `303-v1` 交完之后 `305-a1`。
- Progress log：
  - [2026-10-10 19:2x +0800] agent=编排者 did=收 `300-r2`（七笔 `a0339995`→`9cef4589`，写面并集＝一枚新用例＋自家 probes＋票面一节，越界 0）＋**我自己同台面复跑**（`git archive` 仓外导出树，`porcelain` 0 行；`cmp` 用例↔HEAD blob＝IDENTICAL）：vet `rc=0`／定向 pristine `rc=0`（顶层 1＋子测试 6 逐名）／★`GOFLAGS= go build ./...` `rc=0`（**我那枚派单漏列，本程销账**）／成对两发 3→1 `rc=1`（三句逐字 `:191`／`:214`／`:229`）与 3→0 `rc=1`（两句 `:191`／`:235`），两发各自 `cp`→跑→写回→`cmp` **IDENTICAL**，还原后 `rc=0`／整包两发同一把尺名册**逐字同形**（都只剩 `--- SKIP: TestLiveWasapiSmoke`，`comm` 双向 0）⇒ **新增红 0 枚**／枚数尺 43／`gofmt` 工作树 0 枚／`d22scan` `rc=0`（ban#8 `internal/`=526 实扫）／census totals 逐字未变；★**新现量推翻 `297-a1` 交给 `AC#5` 的前提**：`portable-tests.sh:253` 的 `windows` 档里逐字写着 `./internal/audio/`、`win_pin` 名册含该包，且 `test-windows` 归档字节里 `ok … internal/audio 16.739s`＋`TestParseWaveFormatSubFormatOffset300` 真被求值 ⇒ **"整条 CI 无一档真跑"⛔ 成立**；⚠ 而那枚**新用例**落在 `9a442923`、晚于我分析的 `cf46c24a`／`05db4bc6` 两发 ⇒ 它进 CI 的颜色**仍欠**，本程推送之后才取得；★另量一枚推前前提＝`core_pin`/`win_pin` 是**包名**名册而⛔ 用例数名册 ⇒ 新增用例文件⛔ 需动 pin；**`AC#6` ⛔ 由我翻**（判语归 `300-v4`，ⓐ..ⓕ 六道必答题已列），`AC#4` 另具名两枚欠（腿的门③射程只 `./internal/audio/`、票面要 `./internal/audio/ ./cmd/wisp/` 各 ≥2 发）；现态仍 4 勾／3 未勾，⛔ `-done`；next=推送（`A813`）→ 取新用例的 CI 色 → `300-v4` → `303-v1` → `305-a1`
