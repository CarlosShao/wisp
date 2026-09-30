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

（待填）

---

## §3 AC#3 判语：定向突变的读数（本腿自己这一轮的数）

（待填）

---

## §4 AC#4 判语：越界检查与"生产那一环"的诚实性

（待填）

---

## §5 附：编排者三处就地改账核得对不对

（待填）

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
