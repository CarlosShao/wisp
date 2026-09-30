# 241-r1 — 电平生产者半：PCM 帧 → `0..1` 电平数（internal/audio 包内）

票面：`.scratch/wisp/issues/241-no-anything-computes-a-0-to-1-level-from-pcm-the-whole-main-module-has-no-sqrt.md`（43 行／5,944 字节）
锚点 commit：`501c6971`（`ledger(A463 …`，票 241 立票那一发）
实现腿：`241-r1`（本文件作者） · 裁决腿：`241-v1`（必须≠本腿）
⛔ 本腿未碰票面 AC 五枚框（`git log -- .scratch/wisp/issues/241-*.md` 只有 `501c6971` 一发，见 §8）。

**交付三问的直答**：
1. 标度定义在 `internal/audio/level.go:17-48`（口径正文）＋ `:54/:59/:64/:76/:79-82`（具名常量），`0` 与 `1` 两端的物理含义与容差见 §1；
2. `240-c1` 两枚断点里 **环②（`encodeFrame` 未导出）本票补上了**（`EncodeFrame`＋新 `DecodeFrame`＋往返测试），**环①（`NewBoundedFrames` 生产消费者）本票只补到包内测试级**，生产那一环具名留给**票 228 之后的消费腿**（球进干活进程本身是票 07 的地界），见 §4；
3. 我没跑的跨包门禁列在 §7（期望编排者整包复跑六发）。

---

## 0. 开工前自复尺（票面「起跑名册」要求）

| 尺 | 读数 | 时刻 |
|---|---|---|
| `tasklist //FI "IMAGENAME eq balldebug.exe"` | **0 枚**（`INFO: No tasks are running which match the specified criteria.`）⇒ 两枚碰真设备的用例可以放心跑 | 12:10:20 前后各查一次，两次都是 0 |
| `go test ./internal/audio/ -count=1`（**改动前**基线） | `ok github.com/CarlosShao/wisp/internal/audio 15.863s` | `2026-09-30 12:02:02 +0800` |
| 撞钉预检复跑（票面事实 6 那条尺） | 改动前 `grep -rn "func .*RMS\|func .*Level\|Peak\|Sqrt" --include='*.go' internal/audio/` ⇒ 本腿另以 `grep -rn "Level\|RMS\|Peak\|Sqrt\|math\."` 复跑：**命中全在既有 `math.Sin/math.Round/math.Inf/math.Log10` 与 `math.Float32frombits`，`RMS`/`Level`/`Peak` 三个符号零命中** ⇒ 票面事实 6 成立，新符号不撞钉 | 12:03 前后 |
| `grep -rn "SilenceLevelGate" --include='*.go' .` | 生产命中 `internal/ball/liquid.go:30`（定义）与 `:215`（用），测试命中 `internal/ball/liquid_test.go:60`、`internal/ball/tokens_table_test.go:552`；**`internal/audio/` 零命中** ⇒ **票面事实 2 的出处写错了**，见 §9 冲突 (a) | 12:03 |
| `grep -rln "wisp/internal/audio" --include='*.go' .` | **零命中** ⇒ 票面事实 3「采集栈写完零 importer」在本腿开工时仍然成立（本票不接消费者，交件后仍成立） | 12:03 |
| Sherpa PATH | 本票**未引 sherpa**；仍按派单姿势 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 开局，避免 0xc0000135 那种"用例根本没跑"的假读数 | — |

---

## 1. 标度定义落在哪枚文件哪一行

**文件**：`internal/audio/level.go`（本票新增，147 行）

| 内容 | 行 |
|---|---|
| 口径正文（"THE SCALE"段，含两端定义、量化、线性、DC 说明） | `level.go:8-48`（公式在 `:17-18`） |
| `const LevelFullScale = 32768.0`（分母） | `level.go:54` |
| `const LevelLSB = 1.0 / LevelFullScale`（一枚量化步 = 3.0518e-5） | `level.go:59` |
| `const FullScaleSquareLevel = 32767.0 / LevelFullScale`（对称满幅方波的具名上界） | `level.go:64` |
| `const SineLevelTolerance = 5e-5`（具名容差，推导写在 `:66-75`） | `level.go:76` |
| `const MinLevel = 0.0` / `MaxLevel = 1.0`（闭区间两端） | `level.go:79-82` |
| `func LevelOfSamples(samples []int16) float64`（纯函数，int64 精确累加） | `level.go:96` |
| `func FrameLevel(frame []byte) (float32, error)`（接缝形：一帧进、一个数出） | `level.go:118` |
| `func DecodeFrame(frame []byte) ([]int16, error)` | `level.go:137` |
| `func EncodeFrame(samples []int16) []byte`（原 `encodeFrame` 导出，见 §4） | `wavinjector.go:138` |

**公式**：`level = sqrt(mean(sample^2)) / LevelFullScale`，对一帧 512 枚 int16 样本取**无权 RMS**，除以 int16 能表示的最大幅值 `abs(-32768) = 32768.0`。

**`0` 端的物理含义**：整帧每枚样本都是 0 ⇒ 读**严格 0.0**。这不是地板值、不是近似：口径明写"无噪声地板、无偏移"（`level.go:24-26`），一帧读出 0 只可能是它真的是 0。负零也被测试挡住（`TestLevelSilentFrameIsExactZero` 里 `math.Signbit` 一条）。

**`1` 端的物理含义**：整帧每枚样本都坐在 int16 能带的**最大幅值** `-32768` ⇒ 读**严格 1.0**（这就是 `LevelFullScale` 的定义）。真正对称、不削顶的最响信号是 `±32767` 方波，它读 **`32767/32768 = 0.999969482421875`**（具名常量 `FullScaleSquareLevel`），比 1.0 差**一枚 LSB**——那枚差是 int16 量化本身，不是缺陷；票面 AC#1 写的"1.0 或具名上界"就是这个形状，两枚都被钉（`TestLevelFullScaleSquareEndpoints`）。附带钉住最小非零读数：`±1` LSB 方波读**严格 `1/32768 = 3.0518e-5`**（`LevelLSB`）。

**为什么选 32768 而不是票面/普查提到的 32767 作分母**：分母取 `abs(int16 最小值)` 使**任何合法帧的读数天然落在闭区间 `[0,1]`**，消费者永远不必 clamp、也不存在"悄悄饱和"；若分母取 32767，则 `-32768` 那种样本会读出 `1.00003` 而**越界**，要么引入 clamp（多一条没人定义的裁剪规则）要么让标度说谎。这一点写进 `level.go:50-53` 的常量注释里。

**容差 `5e-5` 与它的推导**（`level.go:66-75`，测量核算在 12:05 用整数级 Python 复算过）：
- 一帧装**整整数个周期**时 `mean(sin^2)` 恰好是 `1/2`，于是残差只剩两笔：(a) int16 舍入，逐样本至多 0.5 LSB，折算到 RMS 约 0.7 LSB ⇒ `2.2e-5`；(b) 满幅对称振幅 `32767` 与分母 `32768` 之间那枚 LSB ⇒ `LevelLSB/sqrt(2) = 2.2e-5`。
- 实算最坏偏差：满幅正弦（amp 32767、每帧 8 周期）读 `0.7070836493`，对 `1/sqrt(2) = 0.7071067812` 偏 **`2.31e-5`**；半幅（amp 16384 = 分母的一半，天然整）读 `0.3535522639`，对 `1/(2*sqrt(2))` 偏 **`1.13e-6`**。
- `5e-5` = 上述上界的约 2.2 倍余量，同时**仍是"八度间距"（`0.1768`）的 1/3536**。
- 容差不许是旋钮：`TestLevelSineToleranceIsNamedAndBounded`（`level_test.go:198`）从两侧钉它——`>= LevelLSB`（不然正弦钉在断言一件不可能存在的舍入）、`<= 1e-3`（不然"在容差内"什么也不说）、且"高八度间距必须 > 100 倍容差"。定向突变 X2（把容差放宽到 `1e9`）实测**只**打红这一枚，正说明它是唯一有牙的那枚（见 §5）。

**本票不决定的那半**：`240-c1` §5 把 **D-1**（"要不要在线性之外加增益/压缩"）登记为**待人裁的观感判定项**。本腿只把**可测的那一半**定成代码：生产者输出的是**物理数**，不含任何曲线（`level.go:36-41` 明写，并有 `TestLevelIsLinearInAmplitude` 作钉——把压缩折进生产者会当场打红）。若 owner 裁"视觉上要增益"，那条曲线属**消费腿**（票 228 之后、球那一侧）叠在这个数之上，`level.go` 保持原始量程的定义权。⚠ 这条判断的算术后果，见 §9 冲突 (c) 下面那段"给 owner 的一句话"。

---

## 2. 新增/改动符号

| 符号 | 位置 | 形状 | 为什么是这个形状 |
|---|---|---|---|
| `LevelFullScale` / `LevelLSB` / `FullScaleSquareLevel` / `SineLevelTolerance` / `MinLevel` / `MaxLevel` | `level.go:54-82` | 具名常量 | AC#1 要"具名上界""具名容差"；无名的数字在裁决时不算定义 |
| `LevelOfSamples([]int16) float64` | `level.go:96` | 纯函数，无状态无时钟无 goroutine | D38b 名册没有"只算电平"的常驻槽（`240-c1` §3.0 P4），消费者在自己线程上调；int64 精确累加保证同输入同位 |
| `FrameLevel([]byte) (float32, error)` | `level.go:118` | 接缝形：`FrameBytes` 一帧进，`float32` 出 | 消费端唯一入口 `Ball.SetAudioLevel(float32)`（`internal/ball/liquid_windows.go:42`）要的就是 float32；帧长不符 ⇒ 报错并把拿到的字节数写进消息，绝不"半帧读成小声" |
| `DecodeFrame([]byte) ([]int16, error)` | `level.go:137` | 导出解码 | 补 `240-c1` 环②；包外不必重抄小端（普查 R21 推荐的就是这一形） |
| `EncodeFrame([]int16) []byte` | `wavinjector.go:138` | `encodeFrame` 改名导出 | 同上；全包只有两处调用（`wavinjector.go:129`、`wasapimic_windows.go:231`），**零测试引用**，改名前后语义逐字相同（含尾零填充） |

未新增：任何 goroutine、任何定时器、任何 `chan`、任何 `import sherpa`、任何消费者接线。

---

## 3. 测试与读数（新增 10 枚名册）

| 用例 | 位置 | 钉的是哪一格 |
|---|---|---|
| `TestLevelSilentFrameIsExactZero` | `level_test.go:48` | AC#1 第一发：静默 ⇒ 严格 0；含 nil 窗口读地板、绝不 NaN |
| `TestLevelFullScaleSquareEndpoints` | `level_test.go:72` | AC#1 第二发：1.0 端、具名上界 `32767/32768`、上界与顶之间恰差一枚 `LevelLSB`、`±1` ⇒ 最小非零读数、半/四分刻度读恰 `0.5/0.25/0.125`、任何合法帧不出 `[0,1]` |
| `TestLevelSineAgainstRootTwo` | `level_test.go:132` | AC#1 第三发：满幅正弦 vs `1/sqrt(2)`、半幅 vs `1/(2*sqrt(2))`、四分幅 vs `1/(4*sqrt(2))`，全用 `SineLevelTolerance`；float32 接缝形同读数 |
| `TestLevelIsLinearInAmplitude` | `level_test.go:172` | 线性（增益/压缩折进生产者即红）＋单调（更响不读更低） |
| `TestLevelSineToleranceIsNamedAndBounded` | `level_test.go:198` | 容差双侧钉住（AC#3 的"放宽到任何输入都能过"唯一克星） |
| `TestLevelSameInputSameBits` | `level_test.go:215` | "同一输入进去同一数出来"钉到**位**（1000 次 `math.Float32bits` 全等）＋换切片承载同样本不变 |
| `TestFrameLevelRejectsNonSeamFrames` | `level_test.go:239` | 七种坏帧长（nil/空/奇数/半帧/差一字节/多一字节/两帧）必须报错、必须返回 `MinLevel`、**错误消息必须含拿到的字节数** |
| `TestSeamCodecRoundTrip` | `level_test.go:276` | AC#2 环②：`EncodeFrame`/`DecodeFrame` 往返含 `±32767/-32768`、尾零填充到 `FrameBytes`、奇数字节不静默解码、**逐字节钉小端** |
| `TestLevelOverBoundedChannelFromWavInjector` | `level_test.go:310` | AC#2 环①（包内侧）：`NewWavInjector` → `NewBoundedFrames()` → 排空 → 每帧 `FrameLevel` **严格 0.5**（`±16384` 方波＝分母的一半）；并断言 `Stats()` `FramesSent=3`/`FramesDropped=0` |
| `TestLevelFrameIsOneSeamFrame` | `level_test.go:369` | 生产者与接缝对"一帧"的口径必须一致（`FrameSamples=512`/`FrameBytes=1024`/`FrameDuration=32ms`） |

读数（`WISP_LIVE_MIC` 未设、`balldebug.exe` 进程数 0）：

- 12:10:30 `go test ./internal/audio/ -count=1` ⇒ `ok ... 15.751s`
- 12:11:xx `go test ./internal/audio/ -count=1 -v` ⇒ 顶层 **36 枚**（起跑名册 26 ＋ 本票新增 10），四数：**`=== RUN` 38**（含 `TestOpenOccupiedAndPermissionDenied` 的 2 枚子测试）、顶层 **`--- PASS` 35**、**`--- FAIL` 0**、**`--- SKIP` 1**；整包 `ok ... 15.959s`。
  - SKIP 点名：`TestLiveWasapiSmoke`（`hotplug_test.go:527`，既有跳过条件 `WISP_LIVE_MIC != 1`，注释逐字"真机冒烟待票 16"）⇒ **不是本票引入的**：那枚 SKIP 只看环境变量，本腿全程 `WISP_LIVE_MIC` 未设（尺：`printenv | grep -c WISP_LIVE_MIC` ⇒ 0）；12:02 那发改动前基线是无 `-v` 的 `ok 15.863s`，所以它没打出 SKIP 行——同一条件同一环境，同一条当时也在跳（这是**机制推定**，不是我在 12:02 那发上直接读到的行）。
  - 同一次跑里 `TestPinnedThreadStable10s` 打出 2 行 `audio frames dropped: bounded channel full (slow consumer)`（`source=wasapi-mic`）——既有真设备用例的正常背压日志，非失败。
- 自我修正两枚（都不是为了变绿放宽断言，方向相反）：① `FrameLevel` 起初没锁帧长，`TestFrameLevelRejectsNonSeamFrames` 把 nil/空帧放过 ⇒ 在 `level.go:119-122` 补上"必须恰为 `FrameBytes`"；② 单调性循环用 `int16` 计数在 32767 之上回绕成负幅值 ⇒ 换成 `int` 计数（`level_test.go:183-185` 注释写明这是测试缺陷不是标度缺陷）。两枚都在提交前修掉，提交 `193daefb` 里就是修好的形状。

---

## 4. AC#2：`240-c1` 认定的两处断点逐条判

`240-c1` §1.4 那张五环表把"采到一帧 PCM"到"算出一个 RMS 电平"拆成 ①通道消费者 ②`[]byte`→int16 解码 ③RMS 本体 ④0..1 标度 ⑤投递免费。本票票面点名的两枚断点是 ①②。

| 环 | 票面认定的形状 | 本票做了什么 | 剩下那一半具名给谁 |
|---|---|---|---|
| ② `encodeFrame` 未导出（普查 R21：包外拿不到帧编解码，只能重抄小端） | "要么导出，要么包外重抄约 10 行" | **补上了，环内闭合**。`encodeFrame` → `EncodeFrame`（`wavinjector.go:138`，两处调用点同批改完、零测试引用），并补上普查 R21 点名"全仓不存在"的反向工具 `DecodeFrame`（`level.go:137`）。包内测试 `TestSeamCodecRoundTrip` 钉往返＋小端字节序＋尾零填充。 | 不再欠谁：包外现在拿得到双向编解码。消费腿若要自己拼帧，不必再抄小端。 |
| ① `NewBoundedFrames` 零消费者（`audio.go:65`；开工时复跑 `grep -rln "wisp/internal/audio"` ⇒ 全仓零 import 者，票面事实 3 在交件时仍成立） | "缺一整块：一个 `for frame := range buf` 的循环体＋它的 owner" | **补到测试级消费者路径为止**：`TestLevelOverBoundedChannelFromWavInjector` 真建 `NewBoundedFrames()`、真用 `WavInjector`（C8 接缝）作生产者、在消费者线程上逐帧排空并算电平，端到端读**严格 0.5**、且 `Stats()` 证明一枚没掉。这条测试同时把 ①②③④ 四环串起来跑通。 | **生产消费者具名欠两枚票**：(1) 装配根里那个 `for frame := range buf` 的循环体与它的 owner —— 票面 ❌ 明写"不接消费者、不改 `cmd/wisp`"（`224-r1` 正在那枚包里），且 `240-c1` §3.3 已量到**跑任务的进程里今天没有球**（`cmd/wisp/resident_windows.go:81` 逐字 `the floating ball arrives in ticket 07`）⇒ 这一环的地界是**票 07 的 GUI 腿 + 票 228 之后的消费腿**；(2) 电平送进球之后的"第六环形状"（`liquidDriven` 不含 `Armed`）票面逐字归**票 228 之后的消费腿**，本腿未碰。 |
| ③ RMS 本体（全主模块 `Sqrt` 零命中） | "缺一整块，约 15 行" | **补上了**：`LevelOfSamples`（`level.go:96`）就是那枚逐帧 RMS；`math.Sqrt` 在主模块生产码里自此**有命中**（本腿尺：`grep -rn "Sqrt" --include='*.go' internal cmd tools` 交件时命中 `internal/audio/level.go:105`，另两枚在同包的 `level.go:11` 注释与 `level_test.go:133`）。 | — |
| ④ `0..1` 标度无定义 | "判定项，不是代码量"（普查 §5 D-1） | **可测那一半定成代码**（§1）：公式、两端、量化、容差、线性全部具名并有钉。**观感那一半（要不要增益/压缩）本腿未裁、也不裁**，保持 D-1 待人裁状态。 | owner 裁 D-1；裁完由消费腿落地。 |
| ⑤ 投递免费（`SetAudioLevel` 自己 `PostTask`） | 属消费侧成本 | 未碰（在 `internal/ball`，票面 ❌ 禁区）。 | 票 228 之后的消费腿。 |

**没有默默留空**：上表两枚票号（07／228 之后的消费腿）与 owner 那一裁（D-1）都具名。

---

## 5. AC#3：定向突变（本腿自演练；⚠ 裁决须由 `241-v1` 独立复跑）

尺子：`.scratch/wisp/probes/241/r1/mutation241.py`（先存 `level.go.pristine`，逐枚改→跑→还原，**临时件只建不删**）。
跑法：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 后 `python .scratch/wisp/probes/241/r1/mutation241.py`。读数存档：`.scratch/wisp/probes/241/r1/mutation-readings.txt`。时刻 `2026-09-30 12:13:39 +0800`。

| 突变 | 改法 | rc | 指名打红的用例 |
|---|---|---|---|
| X1 分子改成常数 | `LevelOfSamples` 末尾 `return rms / LevelFullScale` → `return math.Sqrt(0.25)`（恒 0.5） | 1 | `TestLevelSilentFrameIsExactZero`、`TestLevelFullScaleSquareEndpoints`、`TestLevelSineAgainstRootTwo`、`TestLevelIsLinearInAmplitude`、`TestLevelSineToleranceIsNamedAndBounded`（5 枚） |
| X2 容差放宽到任何输入都能过 | `SineLevelTolerance = 5e-5` → `1e9` | 1 | `TestLevelSineToleranceIsNamedAndBounded`（1 枚） |
| X3 分母换成 32767 那一档的错值 | `LevelFullScale = 32768.0` → `16384.0` | 1 | `TestLevelFullScaleSquareEndpoints`、`TestLevelSineAgainstRootTwo`、`TestLevelSineToleranceIsNamedAndBounded`、`TestLevelOverBoundedChannelFromWavInjector`（4 枚） |
| X4 把压缩曲线折进生产者 | `return rms / LevelFullScale` → `return math.Sqrt(rms / LevelFullScale)` | 1 | `TestLevelFullScaleSquareEndpoints`、`TestLevelSineAgainstRootTwo`、`TestLevelIsLinearInAmplitude`、`TestLevelOverBoundedChannelFromWavInjector`（4 枚） |
| 还原 | 逐字节写回 pristine | 0 | 无（`red=0`；`git status --short -- internal/audio/level.go` 空） |

判读两点，供裁决者复核：
- **X2 只打死一枚**不是弱点，而是这把尺的设计：正弦钉本身对"容差变大"天然无感（放大容差只会让它们更容易过），所以容差**必须**被单独双侧钉住，否则它就是装饰。这一枚就是那条钉。
- **X4 打了 4 枚**证明"标度是真的、不是装饰"的另一半：任何偷偷的增益/压缩都会被 `0.5`/`0.25` 那两枚逐位钉和端到端那枚 `0.5` 抓住。

---

## 6. AC#5：门禁四数终态（逐名照抄）

| 门禁 | 命令（原样） | 时刻 | 读数 |
|---|---|---|---|
| 1 `go build ./...` | `go build ./...` | 12:12:13 +0800 | **rc=0，零输出**（整模块；见 §9 冲突 (b) 的射程说明） |
| 2 `gofumpt -l` | `gofumpt -l internal/audio/` | 12:12:13 | **空输出，rc=0**（gofumpt 在位：`$(go env GOPATH)/bin/gofumpt.exe`） |
| 3 `go vet ./internal/audio/` | `go vet ./internal/audio/` | 12:12:13（另 12:07 一发同果） | **rc=0，零输出** |
| 4 d22scan | `./tools/d22scan/d22scan.exe -root .` | 12:13 前后 | **`d22scan: clean - no D22 ban violations`，rc=0**；`examined 252 production Go files under internal/ and cmd/`；`skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`；live scope work：`bans #1-5 internal/=223, bans #1-5 cmd/=29, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=469, ban #8 cmd/=64`。⚠ 这些数字是**被扫文件数**，不是违规数；⚠ 未跑 `go run ./tools/d22scan`（独立模块，那样必失败）。 |
| 附：包测试 | `go test ./internal/audio/ -count=1` | 12:11 | `ok ... 15.959s`，四数见 §3 |

---

## 7. 我没跑的跨包门禁（期望编排者整包复跑哪几发）

票面禁区原文是"本腿不跑 `go test` 之外的门禁"，括注的理由是"`internal/audio` 之外的包一律不跑，因为有并行写腿的半成品在里面，红了归不清是谁的"；同一张票的 AC#5 又点名要四数。本腿按**最窄一致解**执行：**测试只跑 `internal/audio` 一枚**，其余三发都跑但**范围能收窄就收窄**（`gofumpt -l internal/audio/`、`go vet ./internal/audio/`），`go build ./...` 与 `d22scan -root .` 这两发**天然全模块/全仓**、无法限定范围，故如实交读数并把归因责任交给编排者。冲突记在 §9 (b)。

期望编排者整包复跑／补跑的六发（本腿一枚未跑）：

1. **其它任何包的 `go test`**——一枚都没跑：`cmd/wisp`（并行写腿 `224-r1` 在里面改 `run_mode101_test.go`）、`internal/tools`、`internal/session`、`internal/perm`、`internal/panel`、`internal/ball`、`internal/statemachine`、`internal/agent/approval`。
2. **`gofumpt -l .`（全仓）**——我只跑了 `internal/audio/` 那一档；全仓输出可能含别的写腿的债，归因归编排者。
3. **`sh scripts/d22scan.sh`**——派单/票面点名的只有 `./tools/d22scan/d22scan.exe` 那一形；脚本还多跑一步**正控**（`tools/d22scan` 自己那枚模块的 `go test`，含 seeded-violation），那是"别的包的 `go test`"，本腿按禁区没跑。
4. **`go test ./...`（整仓）／`-count=2` 复跑／`-race`**——一票未跑。
5. **`internal/ball` 那侧的消费验证**：`Ball.SetAudioLevel` 收到的数与球纹波的对应关系（`live_windows_test.go:626/637/645/674` 那四枚手工喂 `0.9` 的用例）——需要真桌面，属票 228 之后的消费腿。
6. **整包 slo/门禁流水线**（`slo-full` 一类）——本腿未启，理由与派单一致（推上去会自启本机抢 CPU）。

跨包证据的**静态侧**本腿倒是复跑过（只读、不跑测试），读数供编排者核：`grep -rln "wisp/internal/audio" --include='*.go' .` 交件时仍**零命中**（票面事实 3 未被本票改变，这是设计如此，不是漏做）。

---

## 8. 越界自查

- 本腿全部 commit（显式 pathspec，`git show --stat` 自证）：
  - `bccc6165` `evidence(241-r1) 骨架先落盘` → 只含 `docs/evidence/s1/241-audio-level-producer-r1.md`
  - `193daefb` `audio(241-r1) 电平生产者半` → 只含 `internal/audio/level.go`、`internal/audio/level_test.go`、`internal/audio/wavinjector.go`、`internal/audio/wasapimic_windows.go`（+552/−4）
  - 本格（证据件填满）另发一发，只含本文件。
- **AC 框**：票面那五枚 `- [ ]` 一枚未碰。尺：`git diff 501c6971..HEAD -- .scratch/wisp/issues/241-*.md` 为空（工单文件自锚点后本腿零改动）。
- **禁区路径零改动**：`cmd/wisp`／`internal/ball`／`internal/speech`／`internal/panel`／`PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`allowlist.txt`——见上面三条 commit 的 pathspec 清单，全部不在其中；`git status --short` 里别的写腿的脏文件（`.gitignore`、`design/**` 那 16 枚删除、`probes/152/**`、`probes/161/r6/logs/**`、`cmd/wisp/*`、`internal/agent/approval/*`、`internal/tools/bridge.go`、`internal/panel/*`、`internal/perm/*` 等）一律**未 add、未还原、未提交**。
- **三枚冻结件**一字未动（不在任何 commit 的 pathspec 里）。
- **`frontend/**`／`design/**`**：零读、零写、零转述（本文件不含任何来自这两层的结论）。
- **git 纪律**：只 commit、无 push；无 `add -A`／`.`／`-a`；无 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。恢复突变后的源文件用的是**脚本内自己存的 pristine 副本逐字节写回**，不是 `checkout`。
- **SLO 阈值／golden／`thresholds.go`**：一字节未动；无任何断言被放宽（§3 那两枚自我修正方向都是**收紧**）。
- **未新增**：goroutine、定时器、channel、依赖边、`sherpa` 引用。

### 8.1 交件前终检（12:18:03 +0800 复跑，三枚 commit 落盘之后）

| 尺 | 读数 |
|---|---|
| `for c in bccc6165 193daefb 1f339ad6; do git show --pretty=format: --name-only $c; done` | 全部 7 枚文件：`docs/evidence/s1/241-audio-level-producer-r1.md`、`internal/audio/{level.go,level_test.go,wavinjector.go,wasapimic_windows.go}`、`.scratch/wisp/probes/241/r1/{mutation241.py,mutation-readings.txt}` |
| 同一条流水过 AC#4 禁区筛 `grep -E "^(cmd/wisp\|internal/ball\|internal/speech\|internal/panel\|internal/observe/thresholds.go\|PLAN.md\|docs/specs/\|frontend/\|design/)"` | **零命中（rc=1）** ⇒ AC#4 那一路径筛为空 |
| `git status --short -- internal/audio/` | **空** ⇒ 突变演练还原后的源文件与已提交状态逐字节一致 |
| `go test ./internal/audio/ -count=1`（终态） | `ok github.com/CarlosShao/wisp/internal/audio 15.914s` |
| `git diff 501c6971..HEAD -- .scratch/wisp/issues/241-*.md` | **空**；工单文件仍 `5,944` 字节 ⇒ 五枚 AC 框一枚未碰 |
| 占位残留尺 `grep -c "^（待填）$" docs/evidence/s1/241-audio-level-producer-r1.md` | **0** ⇒ 骨架那八处独占一行的占位全部写实（按 09-30 死腿收尾第三把尺的姿势自量）。⚠ 天真形态 `grep -c "待填"` 实测 **1**，那一枚命中就是本行自己引用的模式串（自指），不是残留——把这条写给裁决者，免得它拿"1"当成没写完。 |

本腿到此交件给 `241-v1`（裁决者≠实现者）。⚠ 给裁决者的两把必复尺：§5 那四枚突变请**独立复跑**（`python .scratch/wisp/probes/241/r1/mutation241.py`，它会自存/自还原 pristine），以及 §7 那六发跨包门禁**只有编排者/裁决者能跑**。

---

## 9. 票面与盘上不一致的三处（上报，票面文字本腿未改）
按派单要求："如果和票面冲突，以票面为准，并在交件里指出冲突"。这三处都是**票面现量表**与盘上的差别，不涉及判据本身：

- **(a) 事实 2 的出处写错**。票面写 `SilenceLevelGate = 0.06` 的出处是 `internal/audio/gate.go`（"编排者本轮现读"）。实测：`internal/audio/gate.go` **零命中**这个符号；它在 **`internal/ball/liquid.go:30`**（定义）与 `:215`（使用）。`240-c1` §1.4 那一行写的是 `liquid.go:30`，是对的——**是票面表格抄错了，不是普查错了**。判语不受影响（"它是判静默的门限、不是量程"这一条完全成立）。
- **(b) AC#5 与"禁区"那一句互相拉扯**。票面 :39 要"门禁四数"（含全模块 `go build ./...` 与全仓 d22scan），票面 :43 又写"本腿不跑 `go test` 之外的门禁"。本腿执行口径见 §7 第一段。⚠ 请 `241-v1` 按票面 :39 判四数，不要因为 :43 那一句把本腿交的四数读成"越界跑了的包测试"——**本腿一枚跨包 `go test` 都没跑**。
- **(c) AC#1 "半幅正弦 ⇒ 与 `1/√2` 的偏差"这句里两端口径不一致**。半幅正弦的 RMS 是 `1/(2*sqrt(2)) = 0.35355`，只有**满幅**正弦才读 `1/√2 = 0.70711`。本腿**两读都钉**（`TestLevelSineAgainstRootTwo` 同时钉满幅 vs `1/√2`、半幅 vs `1/(2*sqrt(2))`、四分幅 vs `1/(4*sqrt(2))`），所以无论编排者本意是哪一头，判据都是成立的那一头；票面文字未改。

**顺带给 owner 的那一句（D-1 裁定的算术后果，不是本腿的裁定）**：在这个线性量程上，球侧现有的 `SilenceLevelGate = 0.06` 折算过去是 **`20*log10(0.06) = -24.4 dBFS`**；桌面麦克风正常说话的 RMS 大致就在 `-25 ~ -20 dBFS` 那一档 ⇒ **直接把这个数喂进球，"静默门限"和"正常响度"几乎贴在一起**，很可能出现 `240-c1` D-1 那句"接了但看不出在动"。这一条只由**仓内常量加算术**推出（**没有真声音实测**，本腿零开麦），要定案需要 owner 拿真桌面裁"要不要增益/压缩、裁在哪一侧"。裁决落点建议：消费腿（票 228 之后）在球那侧做曲线，`level.go` 保持原始量程。
