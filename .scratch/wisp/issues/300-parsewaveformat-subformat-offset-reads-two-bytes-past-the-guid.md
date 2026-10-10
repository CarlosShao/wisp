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

- [ ] **`AC#0` 权威偏移（⛔ 不许拿本仓另一枚实现当裁判，⛔ 不许凭记忆）**：具名回答"`WAVEFORMATEXTENSIBLE` 里 `SubFormat` 的偏移到底是 24 还是 26"的**外部凭据**——`mmreg.h` 的字段序 ＋ `cbSize` 文档值 ＋ union 的对齐规则，或直接引一枚权威头文件/文档段落（带可核的出处）。★**为什么这格必须挡在 `AC#1` 前面**：`AC#1` 的夹具本身就是按某种布局造的，**如果布局假设错，那枚用例的红什么也证不了（循环）**。⇒ 派单里写死：**`AC#1` 的凭据只有配合 `AC#0` 才成立**；两者不一致时**停手上报**，⛔ 由腿自己改布局再跑一遍。
- [ ] **`AC#1` 决定性一发（⛔ 先红再修，本格只产**测试侧**夹具，⛔ 不改产码语义）**：在 `internal/audio` 新增一枚 windows-tagged 用例，**自己造** 40 字节的 `WAVEFORMATEXTENSIBLE`（`wFormatTag=0xFFFE`、`cbSize=22`、`bits` 各设 16／32、`SubFormat.Data1` 分别置 **1**（PCM）与 **3**（FLOAT）），调 `parseWaveFormat` 并断言返回的 `tag` **等于 1／3**。★判据形状＝**不许**断"被调用过"，断的是**解析出来的那个值**；⚠ 预期今天**两枚都红**（都读出 0）——那就是本票唯一需要的"改前必红"凭据。⛔ 不许为了让它绿而同时改产码（先交红，修归 `AC#2`）。⛔ 不许用真设备。
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

## 编排者补格（2026-10-10 12:2x；⛔ 本节不改写上面任何一行、不勾任何框——只补一条**派单里必须逐字带上**的落地纪律）

- ★**`AC#1` 那枚用例的形状要改一支（我裁，⛔ 让恒红的跟踪测试入库）**：票面写的是"在 `internal/audio` 新增一枚 windows-tagged 用例"，而它**今天必然红**（预期两枚断言都读出 0）。本仓定式＝**恒红测试⛔ 入库**（`A799`/`A801` 那条族：红的跟踪测试会被下一位当成"这棵树本来就红"，还会污染别人的红名册作差）。⇒ 现行射程＝**在仓外成对树里跑那一发**（`git archive` 把当前 `HEAD` 导到仓外，把用例加在**那份导出树**里，`go test` 在那儿跑，红句**逐字**抄进 `.scratch/wisp/probes/300/**` 的 `.md`，用例源码以代码块整块进件），⛔ 往仓里落这枚测试。⇒ 真正的跟踪测试随 **`AC#2`**（修完）同批入库，那时它应当**一进门就绿**，而它的"改前必红"凭据＝本程仓外那一发的逐字读数。⚠ 这不改变 `AC#1` 的判据形状（断的是**解析出来的值**、⛔ 断"被调用过"）。
- ★**`AC#0` 的取法我给一条更便宜的第一选择**：本机若有 Windows SDK 头文件（`mmreg.h`／`mmdeviceapi.h` 之类）就**先读盘上那份**并逐字抄 `WAVEFORMATEXTENSIBLE` 的字段序＋`cbSize` 注释；盘上没有再取外部文档。**两把都要带出处**（文件全路径＋那一行逐字，或 URL＋取回日期），⛔ 拿本仓另一枚实现当裁判，⛔ 凭记忆。
- ⛔ 两格⛔ 由同一枚腿顺手做掉修：`AC#0`＋`AC#1` 交回后**归我裁**，`AC#2` 要**非实现者**裁过才派。


## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
