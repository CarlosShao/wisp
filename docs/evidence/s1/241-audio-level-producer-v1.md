# 241-v1 — 裁决腿验收：`internal/audio` 电平生产者半（票 241 五格 AC）

票面：`.scratch/wisp/issues/241-no-anything-computes-a-0-to-1-level-from-pcm-the-whole-main-module-has-no-sqrt.md`
被验收对象：`internal/audio/level.go`（147 行／6,941 字节）＋ `level_test.go`（398 行／15,247 字节）＋ `wavinjector.go`／`wasapimic_windows.go` 的改名
实现腿：`241-r1`（**本腿不信它的自述**） · 裁决腿：`241-v1`（本文件作者，≠实现者）
本腿起手锚点：`git log --oneline -1` = `aa66ab51`（`ledger(A469 …`）
起手时刻：`2026-09-30 14:02:00 +0800`（`date` 原样输出，非手打）

⛔ 本腿未碰票面五枚 `- [ ]` 框（勾框归编排者）；只碰 `internal/audio`；只 commit 不 push。

---

## §0 逐格 1:1 表（五枚 AC ＋ 编排者要我核的三处改账）

| 格 | 判据原文要点 | 本腿读数凭据（本节只登记"用哪把尺判"，读数随各节补） | 判语 |
|---|---|---|---|
| AC#1 标度是真的 | 静默⇒严格 0；满幅方波⇒1.0 或具名上界；正弦对 `1/√2`（两头都钉）有具名容差；三发缺一不算成立 | §1：逐枚现读用例＋本腿独立复算的端点读数＋突变 M1/M2 | （待填） |
| AC#2 两枚断点逐条判 | `NewBoundedFrames` 零消费者／`encodeFrame` 未导出：要么补＋包内测试，要么具名写给哪枚票；不许默默留空 | §2：现跑全仓 import 名册＋包内消费者测试归因＋欠票号可核性 | （待填） |
| AC#3 攻恒真 | 定向突变（分子改常数 或 容差放宽到任何输入都能过）⇒ 指名用例必须红 | §3：本腿自己的突变台件与前后读数（不复用 r1 的 X1–X4） | （待填） |
| AC#4 越界检查 | `git diff` 出现 `cmd/wisp`／`internal/ball`／`internal/speech`／`PLAN.md`／`docs/specs/**`／`thresholds.go` 任一 ⇒ 退回；另核"生产那一环"有没有假装做完 | §4：按 241 五枚 commit 逐枚 `--name-only` 筛＋全范围差集的归因 | （待填） |
| AC#5 门禁四数 | `go build ./...`／`gofumpt -l`／`go vet ./internal/audio/`／`./tools/d22scan/d22scan.exe` 终态逐名照抄 | §6：本腿自己这一轮的四数（含 gofumpt 的正控） | （待填） |
| 附：改账核三处 | 事实 2 出处／AC#1 半幅正弦口径／AC#5 与禁区自相矛盾那处的裁法 | §5：逐处现读盘上字节＋算术复算 | （待填） |

---

## §1 AC#1 判语：标度是真的还是装饰

**判语：成立**（附两处**文档精度**缺陷登记，都不动摇 AC#1 的三发；见 §1.5/§1.6）。

### 1.1 标度在盘上的形状（本腿逐行现读，非转述 r1）

`level = sqrt(mean(sample²)) / LevelFullScale`，`LevelFullScale = 32768.0`（`level.go:54`），
分子在 `level.go:105` 的 `math.Sqrt` 上落地；累加走 `int64`（`level.go:100-104`）。
接缝形 `FrameLevel([]byte) (float32, error)` 在 `level.go:118`，解码 `DecodeFrame` 在 `level.go:137`。
帧口径三枚常量现读：`FrameSamples = 512`、`FrameBytes = FrameSamples*2`、
`FrameDuration = FrameSamples*time.Second/TargetRate`（`audio.go:47/50/53`），与 `TestLevelFrameIsOneSeamFrame` 钉的一致。

**事实 1 到此翻面**：本腿现跑 `grep -rn "Sqrt" --include=*.go internal cmd tools` ⇒ 三命中，
其中**生产码（排 `_test.go`）只有 `internal/audio/level.go:105` 一枚真运算**（另两枚是同文件 `:11` 的注释与
`level_test.go:133` 的 `math.Sqrt2`）。票面事实 1"全主模块 `Sqrt` 零命中"自这枚改动起**成为历史读数**。

### 1.2 三发的读数（本腿独立复算，Python 整数级 ＋ Go 实跑，两边都对得上）

载具：`.scratch/wisp/probes/241/v1/probe_v1_readings_test.go`，**只用 `go test -overlay` 注入**，
`internal/audio/` 盘上字节一枚未动（终态尺见 §8）。读数存档 `probe-readings.txt`，时刻 14:07:41。

| 发 | 输入 | 本腿实算读数 | 它声称的 | 判 |
|---|---|---|---|---|
| 第一发 静默 | 整帧 512 枚 0 | `LevelOfSamples` = `+0.0`（`Float64bits=0`、`Signbit=false`）；`FrameLevel` = `+0`（`Float32bits=0`、`Signbit=false`）；`LevelOfSamples(nil)` 非 NaN 且 `= MinLevel` | 严格 `0.0`，负零也被钉 | **成立** |
| 第二发 响端 | 整帧 `-32768` | `1` 精确（float64 与 float32 两形都是） | `1.0` 精确 | **成立** |
| 同上 | `±32767` 方波 | `0.999969482421875` ＝ `FullScaleSquareLevel` 逐位等 | `32767/32768`，比顶差一枚 LSB | **成立**（派单第 3 发点名的读数，落在声称值上） |
| 同上 | `1.0 - FullScaleSquareLevel` | `3.0517578125e-05` ＝ `LevelLSB` 逐位等 | 恰差一枚 LSB＝量化不是缺陷 | **成立** |
| 同上 | `±1` 方波 | `3.0517578125e-05` ＝ `LevelLSB` | 最小非零读数（见 §1.6 的口径修正） | **成立** |
| 同上 | `±16384 / ±8192 / ±4096` 方波 | `0.5 / 0.25 / 0.125` 逐位等（分母是 2 的幂） | 半刻度、四分刻度 | **成立** |
| 第三发 正弦 满幅头 | amp 32767、512 样本、8 整周期 | `0.707083649288286`，对 `1/√2 = 0.7071067812` 偏 **`2.313189826153028e-05`** ≤ `5e-5` | 偏差有具名容差 | **成立** |
| 同上 半幅头 | amp 16384 | `0.3535522639041434`，对 `1/(2√2)` 偏 `1.1266891303263193e-06` | 同上 | **成立** |
| 同上 四分幅头 | amp 8192 | `0.176778960510595`，对 `1/(4√2)` 偏 `2.2652139581302855e-06` | 同上 | **成立** |

**"两头都钉"这件事本腿亲自验了**：票面 AC#1 那句被改过的口径（满幅读 `1/√2`、半幅读 `1/(2√2)`）
在 `level_test.go:132-166` 里**确实是两头发各自一条 `t.Fatalf`**，不是只钉一头再把另一头当推论。

### 1.3 "同一个数"这一格：本腿把射程切开判

`TestLevelSameInputSameBits`（`level_test.go:215`）把 1000 次调用的 `math.Float32bits` 钉成全等，
再加"换一片切片承载同样本给出同一位形"。
**判读要切开说，否则这枚钉会被当成它没做到的那件事的凭据**：

- 它**真正**防的是**隐藏状态／非确定性**（结果依赖调用次数、依赖哪块内存、依赖时钟或 map 迭代序）。这一格它是有牙的（见 §3 的 M13 定向突变：把累加器提成包级变量，这枚钉当场红）。
- 它对**标度本身对不对完全无感**：一枚恒输出 `0.5` 的生产者在这枚钉下同样全绿。所以 AC#1 的"标度是真的"**不是**这枚钉交付的，交付它的是 §1.2 那些**逐位等**的端点钉与容差钉。
- `int64` 累加的精确性论证本腿复核成立：最坏 `512 × 32768² = 2⁹ × 2³⁰ = 2³⁹`，距 `int64` 上界 `2⁶³` 有 24 位余量 ⇒ 同输入必然同位，与容差不相关的部分不引入浮点噪声。

### 1.4 十枚新增用例的名册（本腿 `-v` 名册现读，非抄 r1）

`TestLevelSilentFrameIsExactZero` `TestLevelFullScaleSquareEndpoints` `TestLevelSineAgainstRootTwo`
`TestLevelIsLinearInAmplitude` `TestLevelSineToleranceIsNamedAndBounded` `TestLevelSameInputSameBits`
`TestFrameLevelRejectsNonSeamFrames` `TestSeamCodecRoundTrip` `TestLevelOverBoundedChannelFromWavInjector`
`TestLevelFrameIsOneSeamFrame` ⇒ 本腿这一轮**十枚全 `--- PASS`**，顶层名册 26＋10＝36 枚，与 r1 自报枚数相符。

### 1.5 容差推导的**分解写歪了**（登记，不动摇结论）

`level.go:69-72` 把 `5e-5` 的来源写成两笔：(a) "int16 舍入…折算到 RMS 约 0.7 LSB = 2.2e-5"、
(b) "32767 与 32768 之间那枚 LSB ⇒ `LevelLSB/√2 = 2.2e-5`"。**这两笔不能相加**，相加得 4.4e-5，
与它自己写的实测 `2.31e-5` 矛盾。本腿把实测定位的两笔拆开重算（Python，整数级）：

- 纯 (b) 项：`FullScaleSquareLevel/√2 = 0.70708520…`，对 `1/√2` 偏 **`2.158e-5`** ⇒ 这是**主项**。
- 剩下的才是 (a) 项：实测 `0.707083649288286` 与上式的差 **`1.55e-6`** ≈ **`0.05 LSB`**，不是它写的 `0.7 LSB`。
- 两笔合起来 `2.158e-5 + 1.55e-6 = 2.313e-5` ＝ §1.2 那条实测值。⇒ **总数对、上界对、`5e-5` 的余量对，但 (a) 那一笔的归因写错了量级（把主项的数抄给了舍入项）**。
- 本腿顺手把容差的射程推得比测试远一点：整周期数 `cycles ∈ {1,3,7,8,9,64,255}` 满幅正弦偏差全部 ≤ `2.313e-5`（最差仍是测试用的 8 周期那枚）。⇒ 注释说的"Measured worst case over the frames in level_test.go"**没有夸大**，实测射程还略宽一点。
- ⚠ 但这条容差**只对整周期帧成立**：注释 `level.go:66-67` 逐字写了 `on a whole-cycle frame`，射程自陈是准的；非整周期窗口 `mean(sin²) ≠ 1/2`，偏差可以远大于 `5e-5`。**任何人把 `5e-5` 拿去当"任意窗口的通用容差"就是误用**——这一条本腿判它注释写对了，但值得进台账（见 §6.2）。

### 1.6 "smallest non-zero reading" 那一行是**窄了的全称句**（登记）

`level.go:32-33` 的表头写 `smallest non-zero reading / a +-1 LSB square wave reads exactly LevelLSB`。
本腿按派单第 2 发喂了"整帧只有**一枚** `±1` 样本、其余全 0"的帧：

- 读数 **`1.3486991523486091e-06`**（`+1` 与 `-1` 两形同值），**不是 `0.0`** ⇒ AC#1 第一发"静默⇒严格 0"**真严格**：只有逐枚全零才读 0，这跟 `level.go:24-26` 那句 "a frame reads 0 only when it IS 0" 完全自洽。
- 但这个数是 `LevelLSB` 的 **1/22.627**（＝ `LevelLSB/√512`），**低于** `3.0518e-5`。⇒ 那行表头作为"标度能读出的最小非零值"是**假的**；真下界是"一枚 LSB 摊在 512 样本窗上"＝ `1.3487e-6`。作为"**恒幅**帧的最小非零读数"它是真的。
- 判语：**不是 AC#1 的失败**（AC#1 没有要求任何"最小非零读数"），但它是一枚**会误导消费者的注释**：谁把 `LevelLSB` 当"分辨率地板"就会算错 22.6 倍。建议进台账、由票 228 之后的消费腿或补票把那一行收窄成 "smallest non-zero reading of a constant-magnitude frame"。
- 与球侧门限的关系：`1.3487e-6 ≪ SilenceLevelGate = 0.06`（`internal/ball/liquid.go:30`，本腿现读逐字节），单枚脉冲尖峰不会把边框点亮 ⇒ 这一格没有隐藏风险。

### 1.7 DC 那一句：说中了，而且把代价一起说中了

`level.go:37-39` 声称 "DC counts: RMS does not remove a DC offset, so a frame held at -32768 reads 1.0"。
本腿实测：整帧恒 `-32768` ⇒ `1`；整帧恒 `20000`（**纯直流、零交流**）⇒ `0.6103515625` ＝ `20000/32768` 精确。
⇒ 声称与实读一致。物理后果也一并登记在这里：**采集中链上任何直流偏置都会读成"响"**，
本票选择"不除 DC、并在注释里点名由消费者自己开门"，这与票面"线性、无增益、无压缩、无噪声地板"是同一句话的另一面，不是缺陷。

### 1.8 `[0,1]` 闭区间那句话

`level.go:51-53` 的分母论证（取 32768 而非 32767 ⇒ 任何合法帧天然落在闭区间、消费者不必 clamp）本腿在枚举行的形状上全部验证：
包括混合帧"前 256 枚 `-32768` ＋ 后 256 枚 `32767`"读 `0.9999847413273546`（仍 ≤ 1）。
反面也验了：**若分母换成 32767**，整帧 `-32768` 会读 `1.000030518509476` ⇒ **越界**，
`TestLevelFullScaleSquareEndpoints` 的"任何合法帧不出 `[0,1]`"那一条正是抓它的（突变 M6 的读数在 §3）。
⇒ 这不是"两种口径都行"，而是只有一头能保证闭区间；本票选了那一头，并且有钉。**判成立。**


---

## §2 AC#2 判语：两枚断点有没有逐条判、有没有默默留空

**判语：成立（附两枚具名条件）**——两枚断点都逐条给了判、没有默默留空、也没有把"生产那一环"假装做完；
但**环② 的闭合形状带一枚本票新开的 panic 面**，而**欠票的那两个名里有一枚是已关闭的票**。两枚都具名在下面。

### 2.1 "有没有假装做完"这一格：本腿现跑尺，一票否决形状先排掉

派单点名要核这一条。尺与读数（本腿 14:1x 现跑）：

- `grep -rln "CarlosShao/wisp/internal/audio" --include=*.go .`（排 `.scratch`）⇒ **零命中**。
  ⇒ 票面事实 3"采集栈写完零 importer"交件时**仍未被本票改变**，`internal/audio` 之外没有任何新增接线。
- `grep -rn "FrameLevel|LevelOfSamples|DecodeFrame|EncodeFrame|LevelFullScale|SineLevelTolerance|FullScaleSquareLevel|LevelLSB|MinLevel|MaxLevel" --include=*.go internal cmd tools scripts`（排 `internal/audio/`）
  ⇒ **与本票符号相关的生产命中：零**。唯一那几条命中是 `cmd/wisp/leg_sink_nail_131_windows_test.go:333`
  与 `resident_sink_nail_127_windows_test.go:140/457` 的结构体字段 `MinLevel string `json:"min_level"`
  ——**同名不同物**（日志档位，不是电平），本腿逐行读过才这么判；它不是本票的消费者，也不会与本票的常量冲突（不同包、不同类型）。
  ⚠ 顺手给编排者一句：`MinLevel`/`MaxLevel` 这对手感很泛的名字将来会被装配根 import，
  到时候与日志/配置侧的同名概念在同一个文件里并存，**是命名债不是今天的缺陷**。
- `grep -rn "Sqrt" --include=*.go internal cmd tools`（排 `_test.go`）⇒ 生产命中只有 `internal/audio/level.go:105` 一枚。

⇒ **结论：它没装。** 交件自述（`241-r1` §4 表格）说的是"补到测试级消费者路径为止"，并把生产那一半具名欠出去；
盘上事实与这句话**逐字相符**，不存在"注释里有个不存在的调用者"那一形（本腿按 09-30 那条教训专门查了这一枚）。

### 2.2 环②（`encodeFrame` 未导出）：方向补齐了，**编码那一半是partial**

做了的（本腿逐行现读）：`encodeFrame` → `EncodeFrame`（`wavinjector.go:138`），
`git show 193daefb` 的 diff 本腿逐字节读过——**两处调用点同批改完**（`wavinjector.go:129`、`wasapimic_windows.go:231`），
函数体逐字未动（含尾零填充），语义确实"改名前后相同"，r1 这句成立。
新增 `DecodeFrame`（`level.go:137`）是全仓此前不存在的反向工具（普查 R21 逐字"全仓无 `PCM16ToFloat`/`decodeFrame`"）。
包内测试 `TestSeamCodecRoundTrip`（`level_test.go:276`）钉往返、尾零、奇数字节拒解、**并单独钉了小端字节序**。

⛔ **但本腿的探针抓到一枚未声明的前置**（存档 `probe-readings.txt`，14:07:41）：

```
V1API EncodeFrame panics on over-length input: runtime error: index out of range [1024] with length 1024
```

- 复现形状：`EncodeFrame` 收到 **513 枚** `int16`（`FrameSamples+1`）。函数体是
  `frame := make([]byte, FrameBytes)` 后 `for i, s := range samples { frame[2*i] = ... }`（`wavinjector.go:139-143`），
  **既没检查也不文档化 `len(samples) <= FrameSamples`**。
- 今天打不到：两处生产调用点都把长度夹在 `FrameSamples` 以内（`wavinjector.go:129` 的 `min(off+FrameSamples, len(...))`、
  `wasapimic_windows.go:231` 的 `pending[:FrameSamples]`），本腿逐行核过。⇒ **不是线上缺陷。**
- 但它是**本票新开的一枚面**：改名之前 `encodeFrame` 是包内私有，调用者只有那两枚、夹得住；
  导出之后它的消费者恰恰是"装配根里那个还没写的循环体"（正是 §2.3 欠出去的那一环），
  而那一步最自然的写法就是"把手上攒的样本一把丢进 `EncodeFrame`"。
  **测试名册对此是空的**：`grep -n "EncodeFrame(" internal/audio/*.go` ⇒ 九处引用全是 ≤512 枚的正常帧
  与一枚 2 样本的小帧，**没有一枚喂过 `FrameSamples+1`**（本腿现跑名册，不是推断）。
- 与另一半还**不对称**：`DecodeFrame` 明确文档化"any whole number of samples"（`level.go:135-136`），
  本腿实喂 2048 字节 ⇒ 解出 1024 样本、`err=nil`（读数 `V1API DecodeFrame(2048 bytes) -> 1024 samples`）。
  **解码是全函数、编码是偏函数，而这份不对称只写在了一半上。**
- 判语分量：**不推翻 AC#2 的字面要求**（它的要求是"要么补上并给出包内测试，要么具名留给哪枚票"，它补了、也测了）。
  记为**附条件①**：导出形带一枚**未声明、未文档、零测试覆盖**的 panic 前置。
  修法有两条且都很小：要么在 `EncodeFrame` 头部夹掉/报错，要么把前置写进注释并补一枚
  `FrameSamples+1` 的用例把边界钉住。**这条不归本腿做**（裁决者不产码），落给编排者派。

### 2.3 环①（`NewBoundedFrames` 零消费者）：包内侧真通了，**欠票的两个名里有一枚是死票**

包内侧这一半**是真的端到端，不是自说自话**：`TestLevelOverBoundedChannelFromWavInjector`（`level_test.go:310`）
真建 `NewBoundedFrames()`（`:326`，并钉 `cap(buf) == BoundedFrameCapacity()`）、
真用 C8 接缝的 `WavInjector` 作生产者（`:322`，素材是 `t.TempDir()` 里现写的 wav，**没有用 mock 代替真源**）、
在消费者线程上逐帧排空、每帧断言**严格 `0.5`（零容差）**、并回头核 `Stats().FramesSent/FramesDropped`。
本腿的突变把这一枚打死过四次（分子归零、分母换 32767、折压缩、字节序颠倒都让它红，见 §3），
⇒ 它是本票最有牙的端到端钉之一。

⛔ **附条件②：欠票具名具了一枚已关闭的票。** r1 §4 把生产那一半写成"地界是**票 07 的 GUI 腿** ＋ 票 228 之后的消费腿"。
本腿现查工单池：

- `.scratch/wisp/issues/07-ball-state-machine-core-done.md` ⇒ 文件名带 `-done`，正文首行逐字
  `# 07 — Ball shell + state machine core (20 states, 40 transitions, hotkeys, tray) (DONE ✅)`、
  `**Status:** done`。按 AGENTS.md §1.5，`-done` 后缀是**防重领的唯一键** ⇒ **票 07 是关着的，债不能停在它身上**。
- 真正活着的那枚是 `.scratch/wisp/issues/228-ball-and-tray-are-not-in-the-process-that-runs-tasks.md`（无 `-done`），
  它的标题说的就是"球和托盘根本不在能干活的那条进程里"，且票面逐字 `⛔ **AC#1..AC#6 一格不勾**（它们是实现格，实现腿还没动）`。
  ⇒ 这一环的地界**只有 228（及其后的消费腿）这一枚名成立**。
- ⚠ **根子不在写腿身上，在代码注释里**：`cmd/wisp/resident_windows.go:81` 现读逐字
  `wisp: empty event loop running; the floating ball arrives in ticket 07 (Ctrl+C exits cleanly)`。
  票 07 已经关闭而这行注释还说球"arrives in ticket 07"，普查 `240-c1` §3.1 第 1 行又照抄了它，
  r1 再照抄普查 ⇒ **一条过期的落点指认在盘上把三程串了一遍，而今天没有任何仪器看得见它**。
  本腿**不碰 `cmd/wisp`**（那是票 224 的地界，硬约束 2），只把这行原文与它的过期性上交。
- 判语分量：**"不许默默留空"这一条它是满足的**（具了名、也说了 owner 那一裁 D-1），
  所以 AC#2 记成立；但**"具名"的质量不够**——两枚名里一枚是死票，
  按这条债今天实际**没有有效归属**。这一枚必须落台账（§6.2）。

### 2.4 另外三环（③RMS／④标度／⑤投递）：本腿逐条看过，判定与票面口径一致

- 环③（`Sqrt` 零命中＝缺一整块）：**已闭合**，见 §2.1 最后那把尺。
- 环④（`0..1` 标度无定义，普查 §5 把它登记成**判定项 D-1**）：本票只把**可测那一半**定成代码，
  **观感那一半（要不要增益/压缩）没有替 owner 裁**——`level.go:41-48` 逐字把这条留给消费腿，
  并有 `TestLevelIsLinearInAmplitude` 作"不许偷偷折曲线"的钉（本腿的 M8 纯增益与 M9 开方压缩两发都把它打红了，见 §3）。
  ⇒ **判它这一格处理得对**：写代码的程没有顺手替人拍板，而这正是编排者 12:22 那句"维持线性、等真声音实测再摆"要的形态。
- 环⑤（投递到球的那根线程）：**未碰，且应当未碰**（在 `internal/ball`，票面禁区）。本腿确认
  `internal/ball` 零改动（§4 的逐枚 pathspec 筛）。


---

## §3 AC#3 判语：定向突变的读数（本腿自己这一轮的数）

**判语：成立。** 17 发突变（15 发改 `level.go` ＋ 1 发改载体 `audio.go` ＋ 1 发复算 r1 自己的 X3）
＋ 1 发**注释级正控**。**十枚新增用例没有一枚是装饰**——每一枚都被至少一发突变打红（§3.3 的覆盖表就是这一句的证明）。

台件：`.scratch/wisp/probes/241/v1/mutation_v1.py`（rev4，读数 `mutation-run-rev4.log` ＋ `mutation-v1-readings.txt`）
与 `mutation_v1_m16.py`／`mutation_v1_m17.py`。时刻 `14:22:20 → 14:27:00`（rev4）、`14:27:58`、`14:28:41 → 14:28:57`。

⛔ **盘上零污染**：全部突变**只用 `go test -overlay` 注入 `.scratch` 里的副本**，
`internal/audio/` 的四个文件**一次都没被写过**，所以本腿**没有"还原"这一步**，也没有"把跟踪文件清成 0 字节"的可能。
终态尺：`git status --porcelain -- internal cmd` ⇒ **空**（＝本腿起手名册）；
`wc -c internal/audio/{level.go,level_test.go,wavinjector.go,wasapimic_windows.go}` ⇒ **6,941／15,247／6,676／8,587**，
与本腿起手头一把尺逐字节相同。

### 3.1 ⚠⚠ 先交本腿这把尺自己的两次故障（不交这两条，下表那 14 发红不可信）

派单说"宁可表里留没判完，也不要交一枚全在脑子里的件"，本腿把这条扩展到**尺本身**：

- **rev1 故障**：M1 的锚点 `sumSquares += v * v → sumSquares += 0` **删掉了循环变量 `v` 的唯一用处** ⇒
  突变体**根本编不过**（`declared and not used: v`）。而 rev1 把它读成 `rc=1 red_count=0`，
  **看起来完全像"测试没抓到这一发"**。这是本腿差点误判 AC#3 的成因，与 09-30 那条
  "`exit status 0xc0000135` 且没有 `--- FAIL` 行＝用例根本没跑，别判成红"是同一味病的另一形。
- **rev2 故障（更严重）**：解析器匹配 `"--- FAIL "`，而 Go 打的是 `"--- FAIL: "`（冒号粘在词上）⇒
  **12 发突变全体报成 `red_count=0`**，也就是"包全绿、测试是装饰"——**一个彻底反过来的读数**。
  抓它的原因是**唯一一次把 baseline 也当读数用**：`BASELINE rc=0 skip=[]`，
  而这个包已知必须 SKIP `TestLiveWasapiSmoke` ⇒ **skip 为空＝解析器瞎了**，不是包干净。
- rev3/rev4 起的四道防护（都在台件里，编排者可复跑）：① baseline 必须 `rc=0` **且**必须复现那枚已知 SKIP，否则 harness 拒跑；
  ② 出现 `build failed` **直接抛异常**，绝不当成"绿"或"没抓到"；③ **C0 注释级正控**（只把注释里 `C8` 改成 `C9`）
  **必须存活**——一枚连注释都能变红的尺不是在量标度；④ 末尾 `INSTRUMENT-SELFCHECK` 行，
  rev4 实读 **`control_survived=True real_mutants_with_red=14/15 => OK`**。
- ⚠ 一枚**假警报**也如实记：14:21 单发跑 M13 时自检查行报过 `BLIND-HARNESS-ALERT`，
  原因是过滤运行里正控 `C0` 压根没执行（`ctl_line` 为空）。rev4 全量跑已给 `OK`；
  台件里那一支现在对过滤跑明写 `NOT-APPLICABLE` 而不喊 BLIND。

### 3.2 十七发突变的前后读数

| 发 | 改法（锚点 → 改成） | rc | 指名打红的用例 |
|---|---|---|---|
| C0 **正控** | 注释里 `C8 seam`→`C9 seam`（语义中性） | **0** | **无（必须存活；它活了）** |
| M1 | 分子恒零：`sumSquares += v * v`→`v * 0` | 1 | Endpoints, OverBoundedChannel, SineAgainstRootTwo, SineToleranceBound（4）|
| M2 | 返回常数：`return rms / LevelFullScale`→`return math.Sqrt(0.25) + rms*0` | 1 | Endpoints, IsLinear, **SilentFrameIsExactZero**, SineAgainstRootTwo, SineToleranceBound（5）|
| M3 | 容差放宽到 `1e9` | 1 | **只有 SineToleranceIsNamedAndBounded**（1）|
| M4 | 容差收窄到 `1e-9` | 1 | IsLinear, SineAgainstRootTwo, SineToleranceBound（3）|
| M5 | 容差**正好摆在钉的上界 `1e-3`** | **0** | **存活**（见 §3.5，这不是洞）|
| M6 | 分母换 `32767.0`（票面自己提的那一档） | 1 | Endpoints, OverBoundedChannel（2）|
| M7 | 整条取负 `return -rms / LevelFullScale` | 1 | Endpoints, IsLinear, OverBoundedChannel, **SilentFrameIsExactZero**, SineAgainstRootTwo, SineToleranceBound（6）|
| M8 | 偷偷加增益 ×2 | 1 | Endpoints, OverBoundedChannel, SineAgainstRootTwo（3）|
| M9 | 偷偷折压缩（开方） | 1 | Endpoints, IsLinear, OverBoundedChannel, SineAgainstRootTwo（4）|
| M10 | 帧长判据 `!=`→`<`（放过两帧） | 1 | **只有 TestFrameLevelRejectsNonSeamFrames**（1）|
| M11 | 均值分母差一：`len(samples)`→`len(samples)-1` | 1 | Endpoints, OverBoundedChannel, SineAgainstRootTwo（3）|
| M12 | `DecodeFrame` 字节序颠倒 | 1 | Endpoints, IsLinear, OverBoundedChannel, SineAgainstRootTwo, **SeamCodecRoundTrip**（5）|
| M13 | **隐藏状态**：累加器提成包级 `v1Accum` ＋ 调用次数进分子 | 1 | Endpoints, IsLinear, OverBoundedChannel, **SameInputSameBits**, SineAgainstRootTwo, SineToleranceBound（6）|
| M14 | 检测器换成**均值绝对值**（RMS 的最像替身） | 1 | **只有 SineAgainstRootTwo**（1）|
| M15 | 均值绝对值 **＋** 容差同时放到被允许的 `1e-3` | 1 | **只有 SineAgainstRootTwo**（1）|
| M16 | 改**载体**：`audio.go` 的 `FrameSamples 512→511` | 1 | **只有 TestLevelFrameIsOneSeamFrame**（1）|
| M17 | 复算 r1 自述的 X3：`LevelFullScale 32768→16384` | 1 | Endpoints, OverBoundedChannel, SineAgainstRootTwo, SineToleranceBound（4）|

（表中用例名省了公共前缀 `TestLevel`；`SineToleranceBound` ＝ `TestLevelSineToleranceIsNamedAndBounded`。）

### 3.3 派单点名的三件事，逐个直答

**① "AC#1 那三发到底有没有牙"** ⇒ **有牙，但牙的分布不均，这一条必须写清**：

- 票面 AC#3 点名的两形都验了：**分子改成常数**＝M1（4 枚红）／M2（5 枚红）；**容差放宽**＝M3（1 枚红）。
  ⇒ **AC#3 的字面要求成立**。
- ⚠ **M1（分子恒零）之下 `TestLevelSilentFrameIsExactZero` 是绿的**——一枚恒零生产者读静音帧也是 0。
  即 AC#1 的"第一发（静默⇒严格 0）"**单独照不到一个死掉的生产者**；真正把 M1 打红的是**第二发（端点）**。
  票面那句"三发缺一发算不成立"到这里被实测坐实成了一句**有内容的话**。
- ⚠⚠ **本腿最结构性的发现**：M14/M15 把 RMS 换成**均值绝对值**之后，
  **十枚用例里只有一枚变红**（`TestLevelSineAgainstRootTwo`）。
  ⇒ **方波那几枚端点钉根本不区分"RMS"与"另一个齐次幅值泛函"**（方波的 RMS 就等于它的均值绝对值）；
  **唯一钉住"开方均值"这个函数形状的是第三发那枚正弦钉**。
  直接后果（已落 §6.2 台账动作）：**谁将来把正弦那枚用例当冗余删掉，这个包就再也区分不了 RMS 与均值绝对值，
  而包级门禁会全绿。** 这条也反过来支撑 §1 的判语：AC#1 的牙主要在正弦发上。
- 十枚用例的**红名覆盖表**（本腿从 rev4＋M16 的原始输出聚合，非手写）：
  `SilentFrameIsExactZero`←M2/M7；`FullScaleSquareEndpoints`←M1/M2/M6/M7/M8/M9/M11/M12/M13/M17；
  `SineAgainstRootTwo`←M1/M2/M4/M7/M8/M9/M11/M12/M13/M14/M15/M17；
  `IsLinearInAmplitude`←M2/M4/M7/M9/M12/M13；`SineToleranceIsNamedAndBounded`←M1/M2/M3/M4/M7/M13/M17；
  `SameInputSameBits`←**只有 M13**；`FrameLevelRejectsNonSeamFrames`←**只有 M10**；
  `SeamCodecRoundTrip`←**只有 M12**；`OverBoundedChannelFromWavInjector`←M1/M2/M6/M7/M8/M9/M11/M12/M13/M17；
  `FrameIsOneSeamFrame`←**只有 M16（载体突变）**。
  ⇒ **10/10 都有牙**，但**其中四枚各只被一发照到**，而那四发**不在 r1 自己的名册里**（它只打了 X1–X4）
  ⇒ **删掉那四枚中任意一枚，就等于关掉唯一照它的那条缝**。这一条是本腿对 AC#3 的增量，写进台账。

**② "严格 0 是不是真严格"** ⇒ **真严格，两路都验**：

- 实测路（§1.2/§1.6）：单枚 `±1` 样本的帧读 `1.3486991523486091e-06`，**不是 `0.0`**；
  只有逐枚全零才读 0；`LevelOfSamples` 与 `FrameLevel` 两形的 `Signbit` 都是 `false`。
- 突变路：M7（整条取负）把 `TestLevelSilentFrameIsExactZero` 打进红名册（6 枚之一）
  ⇒ **`math.Signbit` 那一枚断言不是摆设，它就是"负零也被钉"这句话的实证**。
- 与球侧 `SilenceLevelGate = 0.06` 的关系（本腿只做算术、不碰 `internal/ball`）：
  `internal/ball/liquid.go:215` 拿 `raw` 与门限直接比（`if raw < SilenceLevelGate`），而 `raw` 就是
  `SetAudioLevel` 收到的那个数 ⇒ **`level.go` 定的量程与球侧门限活在同一个数域上，这句话在盘上说不说得通：说得通**。
  换算本腿自己复算：`20*log10(0.06) = -24.436974992327126 dBFS`；`0.06` 作 RMS 对应正弦峰值 `2780.457` 枚 int16；
  本腿实喂 amp 2780 的整周期正弦读 `0.05998977395253116`（**没过**），amp 3000 读 `0.06473793649783069`（**过了**）
  ⇒ **r1 §9(c) 那笔算术在它自己的量程上对得上**。
  **判不了的部分照写**：正常说话的 RMS 落在哪一档，仓里零真声音读数，本腿也零开麦（§7）。

**③ 端点声明** ⇒ **两枚读数都落在声称值上，无一超出**：
满幅方波（整帧 `-32768`）读 **`1`**（float64 与 float32 两形）；`±32767` 方波读 **`0.999969482421875`**
＝具名常量 `FullScaleSquareLevel` 逐位等。⇒ 任一读数没落在声称之外，AC#1 那半句不因端点判不成立。
⚠ 一处**口径精度**要说清（不是失败）：`FrameLevel` 返回 `float32`，打印成 `0.9999695` 是 Go 的
**float32 最短往返表示**、不是精度损失——`32767×2⁻¹⁵` 只需 15 位有效位而 `float32` 有 24 位，
所以 `level_test.go:96` 的 `float64(got) != FullScaleSquareLevel` 是**逐位成立**的（M6/M9/M11 都把它打红过，证明它是活的）。

### 3.4 r1 自述的四发突变：本腿复算**四发全部逐名对上**

r1 §5 那张 X1–X4 表不是空话，但按规矩**复跑才算裁决凭据**：

| r1 声称 | 本腿复算载体 | 枚数与名册是否一致 |
|---|---|---|
| X1 分子改常数 ⇒ 5 枚 | 本腿 M2（同一形状） | **一致**（同 5 枚） |
| X2 容差 `1e9` ⇒ 1 枚 | 本腿 M3 | **一致**（同 1 枚） |
| X3 分母 `16384` ⇒ 4 枚 | 本腿 M17（照它原值复跑） | **一致**（同 4 枚，集合相等） |
| X4 折压缩 ⇒ 4 枚 | 本腿 M9（同一形状） | **一致**（同 4 枚） |

⇒ **r1 §5 的读数表判为诚实**。本腿另加的 13 发是**它没打过的形状**，
其中 M5/M14/M15/M16 是专门为"这枚钉的射程到底到哪里"设计的。

### 3.5 M5 那枚"存活"是本票设计里的**上界本身**，不是洞

`TestLevelSineToleranceIsNamedAndBounded` 的三支是 `>= LevelLSB`、`<= 1e-3`、`gap > 100*tol`；
把容差**正好摆到 `1e-3`** 时三支全为真（`<=` 含等号；`100×1e-3 = 0.1 <` 实测八度间距 `0.176778960510595`）⇒ 整包绿。
这不违反 AC#3，因为 AC#3 要的是"**放宽到任何输入都能过**"，而 `1e-3` 远做不到：
**M15 就是这一问的实答**——容差放到 `1e-3` 的同时换成均值绝对值检测器，`TestLevelSineAgainstRootTwo` **照样红**。
⇒ 判读：**钉的可用区间是 `[3.0517578125e-05, 1e-3]`（两端含），而这个区间的上端仍然窄到照得出"另一种幅值泛函"**。
本腿把它写成台账建议而不是缺陷：`level.go` 的注释可补一句"容差上界是 `1e-3` 本身、含等号，摆到 `1e-3` 仍在钉内"，
免得后人以为 `5e-5` 是被钉死的唯一值。


---

## §4 AC#4 判语：越界检查与"生产那一环"的诚实性

**判语：成立**（并附一枚**尺本身的口径纠正**给编排者，见 §4.2——按票面字面跑这把尺会**误退本票**）。

### 4.1 逐枚 pathspec 筛（本腿现跑）

`241-r1` 名下共 **5 枚 commit**（`bccc6165` 骨架 / `193daefb` 码 / `1f339ad6` 表满 / `73231f46` §8.1 终检 / `bfd4e524` 尺自指修正），
`git show --pretty=format: --name-only` 逐枚抽出来去重，**全集只有 7 枚文件**：

```
.scratch/wisp/probes/241/r1/mutation-readings.txt
.scratch/wisp/probes/241/r1/mutation241.py
docs/evidence/s1/241-audio-level-producer-r1.md
internal/audio/level.go
internal/audio/level_test.go
internal/audio/wasapimic_windows.go
internal/audio/wavinjector.go
```

按 AC#4 那六个禁区路径名筛（`cmd/wisp`／`internal/ball`／`internal/speech`／`PLAN.md`／`docs/specs/`／`thresholds.go`），
再顺手把 `docs/SLO.md`／`allowlist.txt`／golden／三枚冻结件一并筛上 ⇒ **筛为空，rc=1**（"没找到＝好消息"）。

**AC 框那一格本腿换了把更硬的尺**：不是"框数还是 5"，而是**谁碰过这枚工单**。
`git log --oneline -- .scratch/wisp/issues/241-*.md` ⇒ **只有两发**：`501c6971`（编排者立票）与
`541a9b2b`（编排者收件改账），**两枚都是编排者的 commit，写腿从头到尾没碰过工单文件**。
现跑 `- [ ]` = **5 枚**、`- [x]` = **0 枚** ⇒ 五枚框一枚未翻。这比"数一下框数"强：
它排除了"改完又勾回去"这种数框数照不到的形状。

**阈值与冻结件（整段范围尺，不只本票）**：`git diff --name-only 501c6971..HEAD` 里
`thresholds.go`／`docs/SLO.md`／`allowlist.txt`／`golden`／三枚冻结件（`tokens_fourway_test.go`／
`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）**全部零命中（rc=1）**。
⇒ 不止 241 没动，**整段锚点范围里没有任何程动过它们**。

**`frontend/**`／`design/**` 两层禁**：本票 7 枚文件里零命中；本裁决件不含任何来自那两层的结论
（本腿一次都没读它们，`d22scan` 输出里那两行的**枚数**是仪器的被扫计数、不是内容转述，本腿只把它当门禁读数用）。

### 4.2 ⚠ 这把尺按票面字面跑会**误退本票**——具名给编排者

票面 AC#4 逐字写的是"**`git diff` 里出现** `cmd/wisp`／…任一路径 ⇒ 直接退回"。
本腿先照字面跑了整段范围那一形：

```
git diff --name-status 501c6971..HEAD | grep -E "^(M|A)\s+(cmd/wisp|internal/ball|internal/speech|PLAN\.md|docs/specs/|internal/observe/thresholds\.go)"
  ⇒ cmd/wisp/approval_reply.go
    cmd/wisp/run.go
    cmd/wisp/run_mode101_test.go
    cmd/wisp/subagent_selfapproval_197_test.go
    cmd/wisp/ticket224_assembly_test.go
```

**这五枚不是本票的**：它们属于并行在飞的**票 224／票 197**（`git log` 上对应 `1c601fab`／`ba5db093`／`e8ed5fef`／`43d9096f`／`aa66ab51` 那一串，
且 `internal/session/`、`internal/tools/grant*.go`、`internal/agent/approval/*` 同批出现，是 224 的形状）。
⇒ **共享工作树里"git diff 出现某路径"这句必须有归属口径**，否则任何一张票的越界检查都会被邻居的活打红。
本腿采用的口径是：**AC#4 只判本票名下 commit 的 pathspec 全集**（§4.1 那把尺），
并在 §4.2 这里把字面尺的读数一并交出来，让编排者知道两种读法各自给什么。
**建议把这句口径补进票面 AC#4 或禁区段**（一句话即可："`git diff` 指本票 commit 的 pathspec 集，不含并行写腿"）——
本腿**不改票面**，只登记（§6.2）。

### 4.3 "有没有把生产那一环悄悄假装做完"

判**没有**，凭据在 §2.1：全仓对 `internal/audio` 的 import **零命中**、新符号在 `internal/audio` 之外**零生产调用者**，
而 `241-r1` §4 那张表自己写的是"补到测试级消费者路径为止／生产消费者具名欠两枚票"。
**盘上事实与它的自述逐字相符**，不存在"注释里虚构一个调用者"那一形（本腿逐枚核了它引用的符号名与测试名，见 §5.4）。


---

## §5 附：编排者三处就地改账核得对不对

派单说这三处是"编排者自己写错、被实现腿当场顶正并就地改账"，要我顺手核改得对不对。
**三处全部改对了**，逐处给本腿自己的尺。凭据底座：`git diff --numstat 501c6971..HEAD -- .scratch/wisp/issues/241-*.md`
⇒ **`3 3`**（增 3 删 3），与台账 `A466` 那句"删除列 3＝正是这三处"**逐字对上**；
`git log --oneline -- 该文件` ⇒ 只有 `501c6971`（立票）与 `541a9b2b`（收件改账）两发，都是编排者的。

### 5.1 事实 2 的出处（`internal/audio/gate.go` → `internal/ball/liquid.go:30`）：**改对了**

本腿现跑 `grep -rn "SilenceLevelGate" --include=*.go .`（排 `.scratch`）⇒ 五枚命中：

- 定义 `internal/ball/liquid.go:30`（本腿 `sed -n '30p'` 逐字节读回：`SilenceLevelGate   = 0.06 // envelope below this counts as "nobody speaking"`）
- 使用 `internal/ball/liquid.go:215`（`if raw < SilenceLevelGate {`）
- 测试 `internal/ball/liquid_test.go:60`、`internal/ball/tokens_table_test.go:552`
- **`internal/audio/` 里只剩 `level.go` 的两行注释在引它**（`:13` 与 `:39`）——正是票面新写的那句"`internal/audio` 里今天只有 `level.go` 的注释在引它"。
- `grep -c "SilenceLevelGate" internal/audio/gate.go` ⇒ **0 枚（rc=1＝没找到，这就是好消息）**：旧出处那枚文件里根本没有这个符号。

⇒ 改账准确，且**判语没被改动**（"它是判静默的门限、不是量程"这条在真出处上同样成立：`:215` 那支是"要不要开始算静默时长"的比较，不是量程定义）。
`240-c1` §1.4 那行本来就写的是 `liquid.go:30` ⇒ 普查是对的、抄票时抄错，这也与票面新写的括注一致。

### 5.2 AC#1 的"半幅正弦"口径（改成"两头发各自钉"）：**改对了，而且盘上真是两头都钉**

算术这一头本腿用 Python 独立复算（整数级、不走被测代码）：
满幅正弦 `1/√2 = 0.7071067812`、半幅 `1/(2√2) = 0.3535533906`。原句"半幅正弦 ⇒ 与 `1/√2` 的偏差"
把两端口径混了，**顶正是对的**。
盘上那一头本腿逐行读了 `level_test.go:132-166`：里面是**四条独立 `t.Fatalf`**——
满幅对 `1/√2`、半幅对 `rootTwo/2`、`FrameLevel` float32 接缝形的满幅与半幅各一条、再加四分幅对 `rootTwo/4`；
**不是"钉一头、另一头当推论"**。⇒ 票面新口径"判据按两头发各自钉收（本票终态：两头都钉住了）"**与盘上形状一致**。
本腿的实测偏差：满幅 `2.313e-05`、半幅 `1.127e-06`、四分幅 `2.265e-06`，全在具名容差 `5e-5` 内（§1.2）。

### 5.3 AC#5 与禁区那句自相矛盾（裁"四数照跑、被禁的是跨包 `go test`"）：**裁得对，且本腿按同一口径执行了**

票面 `:43` 原文"本腿不跑 `go test` 之外的门禁"与 `:39` 要四数**确实互相拉扯**——
`build`／`gofumpt`／`vet`／`d22scan` 本来一枚都不是 `go test`，按字面执行反而拿不到 AC#5 要求的数。
编排者的裁法（四数照跑、禁的只是"跑别的包的 `go test`"）是**唯一能让两句同时为真**的读法，
不是择一忽视。`241-r1` 的执行形状本腿核过：`gofumpt`/`vet` 收窄到 `internal/audio`、
`build ./...` 与 `d22scan -root .` 天然全模块/全仓（这两发**无法限定范围**，如实交读数并把归因交给编排者）⇒
**读成合规**，与票面现在写的"它按最窄一致解执行……不读成越界"一致。
本腿自己这一轮跑的四数（§6.1）**同一口径**，且**一枚跨包 `go test` 都没跑**（本腿唯一跑过的 `go test` 目标是 `./internal/audio/`）。

⚠ **本腿给这一格补一句该写进票面的话**：`go build ./...` 与 `d22scan -root .` 在**共享工作树**里
天然会把**别人半成品**的状态一起吃进来。本轮这两发都 rc=0，所以没暴露问题；
但**一旦它们红了，红的可能不是本票的码**。建议 AC#5 那行补"红时须按 commit 归因，不得直接算本票不成立"——
这与 §4.2 那枚 AC#4 的口径缺口是**同一味药**。

### 5.4 顺带：本票引用名字的**存在性**本腿逐枚验了（09-30 那条"注释里指名一枚不存在的测试"的教训）

`level.go`/`level_test.go` 的注释与断言里点名的东西，本腿逐枚查了真身：
`FrameSamples`/`FrameBytes`/`FrameDuration`/`TargetRate`/`BoundedFrameCapacity`/`NewBoundedFrames`/
`NewWavInjector`/`Stats.FramesSent`/`FramesDropped` **全部在 `audio.go`/`wavinjector.go` 里存在**；
`level_test.go` 自用的辅助 `decodeOrFail`/`mustLevel`/`squareFrame`/`wholeCycleSine`/`writeWav` 全部有定义
（`writeWav` 在同包既有测试里，本腿确认它不是新造的）。
⇒ **本票没有"拿一枚不存在的名字当凭据"这一形。**


---

## §6 没做完／留给编排者的台账动作（含本腿的门禁四数读数）

### 6.1 门禁四数（本腿自己现跑，不复用前人读数）

⛔ 前人读数（编排者 12:21 的 `ok 15.55s`、`241-r1` §6 那四数）**一律不充当本腿凭据**；下表每一行都是本腿这一轮的数。

| 门禁 | 命令（原样） | 时刻 +08 | 读数 |
|---|---|---|---|
| 1 `go build ./...` | `go build ./...` | 14:09:53 | **rc=0，零输出**（整模块，天然全仓，见 §6.2 归因条款） |
| 2 `gofumpt -l` | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/ cmd/` | 14:10:03 | **空输出、rc=0**。仪器在位凭据：`D:\work\base\gopath/bin/gofumpt.exe`，5,028,352 字节，`ls -la` 现读在档 |
| 2b gofumpt **正控** | `"$(go env GOPATH)/bin/gofumpt.exe" -l .scratch/wisp/probes/241/v1/posctl/` | **14:13:18** | **列得出故意写歪的那枚文件**（`.scratch\wisp\probes\241\v1\posctl\badly_formatted.go`，rc=0）⇒ 上面那行"空输出"是**真的没有待格式化文件**，不是仪器没跑（派单点名的 rc=127 `command not found` 那一形本腿**没有出现**：`ls -la` 已证 exe 在位、正控已证它会列文件）。正控载具是本腿新建的 `.scratch` 内文件，`internal/` 与 `cmd/` 一枚未碰 |
| 2c gofumpt 逐枚 | `gofumpt -l internal/audio/level.go internal/audio/level_test.go internal/audio/wavinjector.go internal/audio/wasapimic_windows.go` | **14:13:18** | **空、rc=0**（本票碰过的四枚文件逐枚点名，不只靠目录级） |
| 3 `go vet ./internal/audio/` | `go vet ./internal/audio/` | 14:09:57 | **rc=0，零输出** |
| 4 d22scan | `./tools/d22scan/d22scan.exe -root .` | 14:10:36 | **`d22scan: clean - no D22 ban violations`，rc=0**；`examined 252 production Go files under internal/ and cmd/`；skipped as git-ignored：1 file under `frontend/dist/assets/`；口径行：`bans #1-5 internal/=223, bans #1-5 cmd/=29, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=473, ban #8 cmd/=65`。⚠ 这些是**被扫文件数、不是违规数**。⚠ 与 `241-r1` 12:13 那发的差别（`internal/` 469→**473**、`cmd/` 64→**65**）**不归本票**：那 5 枚增量来自并行写腿（票 224/197 落在 `internal/session`／`internal/tools`／`cmd/wisp` 的新测试件），本腿只记差别、不做归因结论。⚠ 未跑 `go run ./tools/d22scan`（独立模块，那样必失败）。⚠ exe 非陈旧：`d22scan.exe` 时间戳 Sep 26 23:33 **晚于** `main.go` Sep 26 23:09 |
| 附 包测试 | `go test ./internal/audio/ -count=1`（`-v` 同一轮） | 14:04:44 → 14:05:02 | `ok github.com/CarlosShao/wisp/internal/audio **16.593s**`；四数 **`=== RUN` 38／`--- PASS` 35／`--- FAIL` 0／`--- SKIP` 1**，顶层用例 **36 枚**（起跑名册 26 ＋ 本票新增 10，逐枚名册见 §1.4）。⚠ 38 = 36 顶层 ＋ `TestOpenOccupiedAndPermissionDenied` 的 2 枚子测试。⚠ 跑前 `tasklist //FI "IMAGENAME eq balldebug.exe"` = **0 枚**（`INFO: No tasks are running which match the specified criteria.`，14:0x 现读）；跑全程 `WISP_LIVE_MIC` 未设（`printenv \| grep -c WISP_LIVE_MIC` ⇒ **0 枚**，即"没找到＝环境变量确实没设"）。Sherpa PATH 已按派单姿势导出 |

**四数终态判读**：四发**全绿**，且每发都带了"仪器真的会响"的正控或时间戳凭据。**AC#5 判成立**（判语在 §6.3）。


### 6.2 本腿没做／做不动／该由编排者落的动作

- **票面五枚 AC 框一枚不翻**：本腿只给判语，翻框与加 `-done` 归编排者。
- **AC#5 那四数的"跨包归因"本腿不裁**：`go build ./...` 与 `d22scan -root .` 天然全模块，工作树里另有 `cmd/wisp`／`internal/tools`／`internal/agent/approval`／`internal/session` 的并行写腿半成品；本腿只交读数，红了算谁的由编排者按 commit 归因。
- **`internal/ball` 那侧的消费验证**（`SetAudioLevel` 收到的数与纹波的对应、`live_windows_test.go` 那四枚手工喂 `0.9` 的用例）本腿**一枚未跑**：那是票 224／239／65／68 的地界，且要真桌面。归票 228 之后的消费腿。
- **生产消费者那一环**（装配根里 `for frame := range buf` 的循环体与它的 owner）本腿**没有代做、也不判它"完成"**；判语只在 §2 落在"具名欠票"这句话可不可核。
- **增益/压缩那一半（普查 §5 D-1）** 本腿不裁：编排者 12:22 已裁"现在不改、登记欠账、必须等真声音实测"（`541a9b2b`）。本腿只核这句话与"线性标度"的代码事实是否自洽。
- **`-race`／`-count=2` 复跑**本腿未做：同一包内两枚真设备用例（`TestPinnedThreadStable10s`）各带 10 秒量级，整包 `-race` 会显著拉长并可能与别的在飞腿争 CPU。要不要跑归编排者裁。

---

## §7 我攻不动的地方（判不动的就写判不动）

- **真机声学侧完全攻不动**：本腿零开麦。`FrameLevel` 从真实 WASAPI 采集出来的帧读数是否名副其实、`SilenceLevelGate = 0.06` 折算 `-24.4 dBFS` 会不会真的与正常说话响度贴住（r1 §9(c) 末段那句"接了可能看不出在动"），**仓内没有任何真声音实测读数可以判它对错**。本腿只能判"这句话在算术上说不说得通"，判不了它在真房间里成不成立。
- **"任何合法帧天然落在闭区间 `[0,1]`"这句的全称量**：本腿能钉的是**枚举出来的**那些形状（方波、整周期正弦、±32768、±1）与一条 `LevelOfSamples` 的 int64 精确累加论证；**非整周期正弦、DC 加交流混合、单枚 `-32768` 与 `32767` 混排的帧**这些输入本腿没有穷尽，也没有工具穷尽（全称量要的是证明不是测试）。⇒ 判语只能说"声称未被证伪"，不能说"声称被证明"。
- **容差推导里的"实算最坏 2.31e-5"**：本腿可以在**测试用到的那几枚帧**上复算，但 r1 那句"Measured worst case over the frames in level_test.go"的射程是"这些帧"、**不是"所有帧"**；要把它推广成标度的全局上界需要一条我没做的界证明。⇒ 若编排者日后要在别处（VAD／SLO）引 `5e-5` 当全局容差，这一格得重开。
- **`TestLiveWasapiSmoke` 那一枚 SKIP**：它的跳过条件是既有代码（`hotplug_test.go`，`WISP_LIVE_MIC != 1`），注释自陈"真机冒烟待票 16"。本腿**没有跑过 `WISP_LIVE_MIC=1` 的那一支**（要真麦克风、会抢设备、且不在本票地界）。⇒ 本腿能判"这不是本票藏的红灯"，判不了"真机那支本身是绿的"。
- **`float32` 折算的电平上界**：`FrameLevel` 返回 `float32(LevelOfSamples(...))`。本腿能论证 int16 域内只有有限多个值、`float32` 精度远超；但"消费端拿 float32 会不会在某条后续曲线上丢单调性"属票 228 之后那一侧，本腿不判。
- **D38b goroutine 名册那一句**（生产者不开协程所以不占常驻槽）：本腿能读码证明 `level.go` 无 `go func`/`chan`/`time.`；**判不了**装配根将来必须在哪条线程上调它（P3 那条"采集线程不得跑其他 Go 代码"的裁量属票 07／228，见普查 §3.0 P3）。
