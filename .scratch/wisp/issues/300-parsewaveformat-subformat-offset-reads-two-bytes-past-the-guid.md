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
- [ ] **`AC#2` 落地那一发（⛔ 只有 `AC#1` 红过并经非实现者裁之后才动）**：把偏移改成与权威布局同形的那一枚，**同批**把 `parseWaveFormat` 头顶那三行布局注释改成实话（⚠ 注释现在逐字教读者相信 `SubFormat@26`，这比代码更危险——它会喂下一个改这里的人）；⛔ 顺手格式化别人的文件；⛔ 动电平尺定义（`internal/audio/level.go` 的 `SineLevelTolerance`／`internal/ball/liquid.go:30` 的 `SilenceLevelGate = 0.06` 属 C21 冻结 token 面）；⛔ 动 `internal/observe/thresholds.go`（尺＝208 行、`level` token **0 命中**，我现跑）。
- [ ] **`AC#3` 与票 297 那三形的关系（判语归位，⛔ 不越权结案）**：`AC#1`／`AC#2` 成立⇒ 票 297 三支里的 ③（"切错字节"）拿到一枚**可复现的形状**，但⛔ 本票不许把 297 的"是什么"判成已定案——真机那枚混音格式到底是不是 extensible＋float，归**编排者**那一发 `GetMixFormat` dump（〔仅本机可量〕）。⇒ 两票各自留格，⛔ 合并结案。
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


## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-10 12:32:39 +08] agent=300-a1r did=AC#0 权威偏移读定＝SubFormat 在 24（凭据＝盘上 Windows SDK shared/mmreg.h，第 34 行 pshpack1.h 令 Packing=1 ⇒ 偏移表可算；cbSize=22 与总长 40 由 mmreg.h:2540/2550 与 Learn 那句 "must be at least 22" 从"编排者推的"升为"读到的"；union 那一问＝不是，Samples 与 dwChannelMask 是两枚独立成员，标准结构里没有 wValidSamples；据此 AC#0 结论把 wasapi_windows.go 判成错的那一枚、wavinjector 同形，硬门"若是 26 才停手"未触发）＋AC#1 仓外成对树（git archive b2933dc9 到 /tmp/wisp300-a1r/tree，自造 40 字节面 4 形喂 parseWaveFormat，断解析出的 tag 值）：改前 4/4 FAIL（S1/S2/S3 恒读 0，S4 那枚 Data1 挪到 26 的坏面读到 3 ⇒ 这枚尺两个方向都敏感）、只把树上 unsafe.Add(p, 26) 改 24 后同一份用例 4/4 PASS、随后树已 cmp 验过逐字节还原；恒红用例⛔ 入库（件只落 probes/300/a1r/ 的 00/10/20/90 与 logs/ 七枚 .txt）；三把卫生＝d22scan rc=0、GOFLAGS= go build rc=0、git status --porcelain -- internal cmd 0 行 next=编排者裁 AC#0/AC#1（⛔ 我翻框、⛔ 我动 AC#2）；给 AC#2 两条输入：偏移改 24 时头顶注释里 validBits@18/channelMask@22/SubFormat@26 三个数全是坏的、建议入库用例保留 S4 那枚反向形；AC#3 真机 GetMixFormat 与 AC#5 CI 三跳本程⛔ 跑，读数仍欠
- [2026-10-10 12:38:01 +08] agent=编排者 did=收 300-a1r（四笔 70a6b9d9→efd85155、四枚 .md 逐枚 find -printf 现量 4,612／12,509／17,598／11,103、logs 8 枚其中 0 字节＝0；⚠ 一处过程痕迹＝hygiene-build.txt 在 1c007bf6 那笔是 0 字节入库、下一笔才填 191 字节，按"0 字节的件＝那格没交"那一笔当时不合格，它自己补对了）⇒ **翻 AC#0／AC#1 两格**，⛔ AC#2..AC#5 一枚不碰（AC#2 开工权在非实现者 300-v1 手上，票面那句"⛔ 只有 AC#1 红过并经非实现者裁之后才动"我⛔ 让它静默 dissolve）；★**承重两发我自己重跑，⛔ 引用腿的数**：① 我自己开那枚头文件 `C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`（176,680 字节）逐字读 `sed -n '2521,2534p'` 的 WAVEFORMATEXTENSIBLE ＋ `:31-36` 的 pshpack1 ＋ `grep -n "cbSize = 22"` 命中 :2540／:2550 ⇒ 我自己算的偏移＝WAVEFORMATEX 18、Samples@18、dwChannelMask@20、**SubFormat@24**、总长 40 ⇒ **判 24 对、wasapi_windows.go 那枚 26 错**；② 我在腿留的导出树 /tmp/wisp300-a1r/tree 先核那行逐字仍是 unsafe.Add(p, 26)，再跑两形＝**原样 4/4 FAIL**（红句逐字：S1 tag = 0, want 3／S2 tag = 0, want 1／S3 tag = 0, want 7／S4 tag = 3, want 0 ＋ floating 两行）／**只把树上 26 改成 24 ⇒ 4/4 PASS**，还原证明＝同一条命令内 cp 备份→跑→写回→`cmp` ⇒ restore_cmp=IDENTICAL 且 `grep -c 'unsafe.Add(p, 26)'` 回到 1 ⇒ ★这把尺两形配对成立、⛔ 恒真三形之一；由这两发现在能说的最硬一句＝extensible 混音格式下 tag 恒 0 ⇒ floating=false ⇒ convertPacket 走 int16 那条 unsafe.Slice（⛔ 等于"票 297 那 0.3677 就是它"，那一句仍只有真机那一对能量）；★**夹具⛔ 循环**：我逐字节读了 build()（b0=0xFFFE／b16=22／b18=bits 写进 Samples／b20=mask／Data1 写 data1At，S1-S3＝24、S4＝26）⇒ 与头文件逐格对得上，⛔ 拿仓里任一枚实现当布局裁判；★**接受腿两处顶我、只追加⛔ 改原句**：票面第 13 行那句"nValidSamples 与 dwChannelMask 是同一个 union 的成员"⛔＝union 的名字是 Samples（三枚 WORD：wValidBitsPerSample／wSamplesPerBlock／wReserved），dwChannelMask 是独立成员，头文件里⛔ wValidSamples 这个名字（结论层不变＝偏移仍 24）；另我派单里那句示例路径 Include/*/um/mmreg.h 在盘上⛔ 成立（真身 shared/）⇒ 定式并回：派单出现盘上路径，落笔前先 ls/find 一发（与 A797"声明 HEAD 对象层射程前先 git show"同族）；★**新事实改写我 1 小时前刚立的定式（A802 追加更正 A801）**：`git ls-files third_party/sherpa-onnx | wc -l`＝**0** 而盘上有三枚 DLL ⇒ `git archive` 导出树里那条 harness PATH 指向**不存在的目录**，"成对导出树"法⛔ 自带这处 handicap；300-a1r 是 cp 三枚进树才让尺成立 ⇒ 今后用导出树跑带 harness 的包前先 `ls <树>/third_party/sherpa-onnx/*.dll | wc -l` 并具名；⚠ 连带一枚我⛔ 现在下结论的问题＝293-v1 的件里⛔ 提过 DLL（我 `grep -in "dll" .scratch/wisp/probes/293/v1/20-ac4-gates.md`＝0 命中），它那两组 677／684 要不要打折 ⇒ 交 300-v1 三档裁（无影响／射程受限／要打折）；③ 腿建议"AC#2 入库版保留 S4 那一形"＝裁法交 300-v1，但我先在 AC#2 派单钉⛔ 删 S4；next＝在飞 1 枚 `300-v1`（锚 efd85155，攻★ⓑ／★ⓒ 这套凭据＋裁 AC#2 开工权＋答那枚打折问题）；它交回前⛔ 派 300-r1、⛔ 任何人改 internal/audio；之后 300-r1（偏移改 24＋同批把 validBits@18, channelMask@22, SubFormat@26 那三行注释改成实话——那行注释今天逐字教下一个改这里的人相信 26，比代码危险）→ 300-v2 → AC#3（与票 297 关系归位，⛔ 合并结案）→ AC#5（CI 那三跳）；队列其余 296-r2 → 298-r1 → 295-r1 → 294-r1 → 294-v1；票 299 按在 Q-84（⛔ 不催）。⛔ 零 push、⛔ 本票零产码
