# 票 291 `AC#3` — 落点表（**只读腿，本票未落地**；票面逐字「本票不许落地」）

- 钟点：2026-11-14 21:5x–22:1x +08（`date` 现跑）。
- 树锚：`acc1e9e23`（主仓 HEAD，本腿开工时现读；`git rev-parse --short HEAD`）。
- 写面：**只有本件** ＋ `.scratch/wisp/probes/291/ac3/**`。**零产码改动进主仓**（`git status --porcelain -- internal cmd` 开工/收工均为空，见 §9）。**零 push**、⛔ 真麦（未设 `WISP_LIVE_MIC`）、⛔ 阈值/golden/契约/票面/台账。
- 方法：主仓**只读**；全部改前/改后读数在一棵**仓外导出树**跑（`git clone --shared` 到 `D:\tmp\wisp\20261114\291-ac3-tree`，HEAD 同串 `acc1e9e23`）。导出树里的改动是**探针**，不是落地：文件逐字见 `07-proposed-levelgo-change.patch.txt`／`08-proposed-fixture-rebaseline.patch.txt`／`09-proposed-new-nail.txt`。
- 落点口径：**"最小改法"＝只动 `internal/audio/**`**（票面 `AC#4` 允许的落地写面），**不碰** `cmd/wisp/**`、`internal/ball/**`、任何冻结件。

**一句话**：去直流的最小落点是 **`internal/audio/level.go` 的求和段改成"先去本帧均值、再算 RMS"（两趟、float64、无状态、无分配）**；它**不碰** `SineLevelTolerance`/`MinLevel`/`MaxLevel` 三枚常数、**不碰** `thresholds.go`/golden/`allowlist.txt`/三枚冻结件/C21 token；但它会让 **2 枚既有钉必红**（不是 1 枚：`TestLevelFullScaleSquareEndpoints` ＋ 票 297 那枚 `TestAC297KnownAmplitudePairSeparates`，后者票面此前没点名）、**1 枚夹具静默失去见证力**（`capturelevel_windows_test.go`），并让 **`MaxLevel = 1.0` 变成不可达**、`SilenceLevelGate = 0.06` 的**物理含义搬家（数字一字不动）**。

---

## 0. 结论先行：落点 5 处（5 枚文件）＋ 2 枚只登记不动

| # | 文件:行（`acc1e9e23` 现读） | 动作 | 读到的原文/现状 | 必红？ |
|---|---|---|---|---|
| 1 | `internal/audio/level.go:96-107`（求和段 `:100-104`，返回 `:105-106`） | **改算术**：先 `sum` 求本帧 `mean`，第二趟累加 `(s-mean)²`，`sqrt(...)/LevelFullScale` | `var sumSquares int64` / `sumSquares += v*v` / `rms := math.Sqrt(...)` —— **全波 RMS，无去均值步** | — |
| 2 | `internal/audio/level.go:17-39` 与 `:89-95`（尺的定义面注释） | **改注释**（不是放宽断言）：`:17` 的 `level = sqrt(mean(sample^2))`、`:22-26` 的端点段、`:37-39` 的 `DC counts` 段、`:89-91` 的 int64 精确性段 | `:37-39` 逐字「`DC counts: RMS does not remove a DC offset, so a frame held at -32768 reads 1.0.`」 | 注释，不红；**但它是尺的定义，必须与算术同批改** |
| 3 | `internal/audio/level_test.go:72-88`（`TestLevelFullScaleSquareEndpoints` 的 `-32768` 块） | **改期望＋改断言句**：恒定帧从「恰好 `MaxLevel`」改成「恰好 `MinLevel`」 | `:79-81` `if l := LevelOfSamples(hot); l != MaxLevel`；`:86-88` `FrameLevel(...) != MaxLevel` | **必红**（实测，§6） |
| 4 | `internal/audio/amplitudepair_297_windows_test.go:59-65`（夹具信号）＋`:111-128`（判据） | **重定夹具信号**：纯 DC 窗 → **DC 偏置＋±A 方波**；并在比值前加 `meanLow <= 0` 守门 | `:59-65` 两窗全是**常量**（`samples = append(samples, lowCounts)`），`:122-128` 断言 `meanLow ≈ |counts|/FS` 与 `ratio ≈ 8.0` | **必红**（实测，§6）——**票面/派单此前只点了 `capturelevel`，漏了这枚** |
| 5 | `internal/audio/capturelevel_windows_test.go:21-24` ＋ `:48-54` | **期望值不改**（派生式）；**追加**一条「恒定 DC 帧读恰好 0」的断言，把失去的见证力换成正面见证 | `:23` `samples[i] = 12345 // a steady DC level: RMS is exactly this magnitude`；`:50` `want := float32(LevelOfSamples(...))` | **不红**（派生式，实测 0==0）；真代价＝那句注释变假话、夹具空转 |
| 6 | **新建** `internal/audio/level_dc_291_test.go`（形状见 `09-proposed-new-nail.txt`） | 新增 3 枚钉：恒定帧读 0／**平移不变（位级）**／最高合法帧 `< MaxLevel` | — | 新增；其中 2 枚在**改前是红的**（实测，§6） |
| — | `cmd/wisp/resident_audio_247_live_windows_test.go:194-201` | **只登记、不改**（归票 247，且要非实现者裁） | `:194` `if loud.mean() <= quietA.mean()*2 \|\| loud.max() <= quietA.max()` | 真麦台件，**不在** `go test ./internal/audio/` 名册；去直流后**含义变松**，见 §7.4 |
| — | `internal/ball/liquid.go:30`（C21 冻结 token，机钉 `internal/ball/tokens_table_test.go:552`／`:572`） | **⛔ 一字不动**，具名交编排者摆 owner | `SilenceLevelGate = 0.06 // envelope below this counts as "nobody speaking"`（用在 `liquid.go:272`） | 不红；**语义搬家**（数字不动、单位动） |

---

## 1. 求和那一段：现在怎么算、DC 从哪进、改哪几行、为什么选这一种

### 1.1 现在怎么算（逐字，`level.go:96-107`）

```go
func LevelOfSamples(samples []int16) float64 {
	if len(samples) == 0 {
		return MinLevel
	}
	var sumSquares int64
	for _, s := range samples {
		v := int64(s)
		sumSquares += v * v
	}
	rms := math.Sqrt(float64(sumSquares) / float64(len(samples)))
	return rms / LevelFullScale
}
```

- **全波 RMS**：`mean(s²)` 里的 `s` 是**带直流**的原始样值；分母 `LevelFullScale = 32768`（`:54`）。
- 一圈 `512` 样／`32 ms`（`audio.go:47-53`），`FrameLevel`（`level.go:118-128`）只是"解码 + 调它"。
- 生产唯一调用点：`captureopt.go:121` `level, err := FrameLevel(frame)`（`emitLevel`，`:117-127`）。
- 因此**任何恒定分量都整份计入**：`-32768` 恒定帧读 `1.0`（`:79` 逐字钉着），`12345` 恒定帧读 `12345/32768 = 0.376740`。

### 1.2 DC 从哪进（四条入口，逐条给读数/出处）

| 入口 | 位置 | 量级（现读/实测） |
|---|---|---|
| ① 设备/ADC 直流偏置 | 真麦 → `wasapi_windows.go:404-414` `convertPacket` | 票面 10-09 安静窗 `0.368462`（≈12073 counts）；`247-v1:92` 同族；本腿**没有**真麦读数，⛔ 不复述成"今天还这样"（`291-mic` 11-08 同机安静 `mean 0.003067`，票面前提当日不可复现） |
| ② **字节切错**（float32 包按 int16 切） | `convertpacket_297_windows_test.go:120-144` 的 S2 形；已由票 300 @24 修 | 本腿实测该合成窗 **均值 = −3998 counts**；基线电平 `0.383045` → 去直流后 `0.363094`（仍 `> 0.1`，该钉绿色） |
| ③ `FloatToPCM16` 的**非对称刻度** | `resample.go:119-128`：`+32767`／`-32768` 两侧系数不同 | 对满幅正弦引入 `≈0.5 LSB` 的负偏（`1.5e-5` 电平，**低于** `SineLevelTolerance`），去直流后归零 |
| ④ 多声道下混 | `resample.go:102-116` `MonoDownmix`，`acc/chans` 向零截断 | 对称截断 ⇒ **不产生**直流（负零两侧同幅），列出来是为了把"不是它"也钉住 |

### 1.3 最小改法（逐行，可直接照抄）

```go
	// 票 291：先去掉本帧自己的直流（本帧均值），再取 RMS。
	// 两趟、float64、无状态、无分配；一圈 512 样 / 32 ms。
	var sum int64
	for _, s := range samples {
		sum += int64(s)
	}
	mean := float64(sum) / float64(len(samples))
	var sumSquares float64
	for _, s := range samples {
		d := float64(s) - mean
		sumSquares += d * d
	}
	rms := math.Sqrt(sumSquares / float64(len(samples)))
	return rms / LevelFullScale
```

- 动的是 **`level.go:96-107` 之内**（`:100-104` 那 5 行被替换成 10 行）；**函数签名、返回类型、零分配、无状态、无时钟、无 goroutine 全部不变**。
- **两趟而不是一趟**：一趟的代数式 `E[s²] − (E[s])²` 会在大直流＋小交流时**灾难性抵消**（纯 DC 帧会算出 ~1e-6 的假电平，而不是**恰好 0**）；两趟能保证"恒定帧偏差逐项为 0"（实测 `float64 bits = 0000000000000000`，§6）。
- **均值用 `int64` 求和再除**：`sum` 在 512 样上是精确整数（`512×32768 = 2^24`），`mean` 只要真商可表示（恒定帧必然）就**逐位等于该常数** ⇒ 恒定帧残差**恰好 0**。
- **不采用**"整数精确方差 `n·Σs²−(Σs)²`"：`LevelOfSamples` 接受任意长度窗，`n²·32768² > 2^63` 时 `int64` 会**静默回绕**；实测安全上界 **n = 92681 样（16 kHz 下 5.79 s）**（`arith_probe.out.txt`）。要么加长窗就错的隐患，要么多一个分支——**买不到价值**。

### 1.4 为什么选"去本帧均值"，不选"DC 阻断器"或"一阶高通"（按主流音频前端做法，⛔ 按动作小）

三种候选的**实测/算得**对比（同一把 512 样、16 kHz 尺）：

| 候选 | 直流行为 | 低频代价 | 与本仓契约的冲突 | 判决 |
|---|---|---|---|---|
| **A. 去本帧均值（推荐）** | 恒定帧**恰好 0**；任意 `s+c` **位级不变** | 20 Hz 最坏 **−3.02 dB**、15 Hz −7.98 dB、31.25 Hz **0 dB**、85 Hz **−0.076 dB**、≥300 Hz ≤0.009 dB（实测扫频，`arith_probe.out.txt`） | **无**：纯函数、帧内、无状态、无分配；同一枚尺（一条定义） | ✅ |
| B. 有状态一阶 DC 阻断（`y[n]=x[n]−x[n−1]+R·y[n−1]`，主流前端常用 `R∈[0.99,0.999]` ⇒ 16 kHz 下拐点 2.5–25 Hz） | 直流稳态为 0，但**有尾**：一次 DC 跳变/咔哒会拖 2–3 帧（50–100 ms）假响度 | 85 Hz 约 **−0.10 dB**（与 A 基本同量级） | **三处**：① 电平不再只由本帧决定 ⇒ `TestLevelSameInputSameBits`（`level_test.go:215-234`）的"same input same bits"要改成"同一状态机、同一历史"；② 生产者变成"带隐藏历史的尺"，正是 `level.go:41-48` 明写要避免的形状；③ 需要新 API（`LevelMeter` 对象）＋改 `captureopt.go:121` 调用点 ⇒ **同一包两条尺**（`LevelOfSamples` 仍含直流） | ❌ |
| C. 一阶高通/差分（`x[n]−x[n−1]`） | 直流为 0 | **85 Hz ≈ −29.5 dB**（`2sin(ω/2)`），把男声基频整段削掉 | 无状态冲突，但**这不是电平表，是微分器**；响度不再随幅度线性 | ❌ |

**理由（为什么 A 最有价值，不是"动作小"）**：这个数字的**用途**是让球随声音呼吸、并喂 `SilenceLevelGate`——它要的是"这一帧里**变化**了多少"，不是"这一帧的绝对摆幅"。A 用**最小的机制**得到**唯一正确的那件事**（平移不变），同时保住 `level.go:84-95` 写死的"纯函数／无分配／无状态／同输入同比特"四条；B 的低频代价并不比 A 小（85 Hz 都是 ~0.1 dB），却要付"有尾＋两条尺＋改 API"的代价；C 直接削掉语音基频。⇒ **A**。

---

## 2. `SineLevelTolerance` 那把双边尺：**不碰，逐条说清为什么**

`level.go:66-76` 的常数与 `level_test.go:198-212` 的 `TestLevelSineToleranceIsNamedAndBounded` **一字节都不需要动**，而且**实测改后仍绿**（`03-*`／`05-*` 两份日志里逐名 PASS）。逐条：

1. **常数值不动**：`SineLevelTolerance = 5e-5` 的推导（`:66-75`）只有两项误差——int16 舍入（≤`0.7 LSB = 2.2e-5`）与 `32767/32768` 的 1-LSB 缺口（`LevelLSB/√2 = 2.2e-5`）。去直流**不改变**这两项，也**不新增**第三项：整周期正弦的窗均值本来就是 0（它的舍入残差均值是 `O(0.01 LSB)`），减去它带来的电平变化是**二阶小量** `m²/(2·MS) ≈ 1e-12`。
2. **双边钉（下界 `≥ LevelLSB`、上界 `≤ 1e-3`、100× 八度差）全部不变**：`LevelLSB` 不变、八度差（`16384` vs `8192` 整周期正弦，窗均值都 ≈0）不变。实测 `TestLevelSineToleranceIsNamedAndBounded` 改后 **PASS**。
3. **实测的逐枚误差一字未动**（`arith_probe.out.txt`，改前＝改后）：

| 幅度 | 理想 `A/LevelFullScale/√2` | 改前误差 | 改后误差 | 容差 |
|---|---|---|---|---|
| 32767 | 0.707085202 | −1.553e-06 | **−1.553e-06（位同）** | 5e-5 |
| 16384 | 0.353553391 | −1.127e-06 | **−1.127e-06（位同）** | 5e-5 |
| 8192 | 0.176776695 | +2.265e-06 | **+2.265e-06（位同）** | 5e-5 |

4. **唯一"碰"到它的是注释**：`level.go:35` 那句「a full-scale sine reads 1/√2 within SineLevelTolerance」仍然**成立**，可以不改；但 `:17` 的公式串要改成"先减均值"。
5. ⚠ **不许借这次把容差放宽/收紧**（票面禁区）。本表**没有**任何数值改动提案。

---

## 3. 两枚 DC 夹具：一枚不会红但会空转，一枚**必红**（票面漏点）

### 3.1 `capturelevel_windows_test.go`（票面点名的那枚）——期望值**不用改**，实测也不会红

- 夹具：`:21-24` 八帧、每样 `12345`（`// a steady DC level: RMS is exactly this magnitude`）。
- 断言：`:50-53` `want := float32(LevelOfSamples(samples[i*FrameSamples:(i+1)*FrameSamples]))` 再与回路吐出的 `l` 比 —— **两边走同一个函数**，派生式，全仓零枚 `0.3767` 字面（本腿现跑 `grep` 复核）。
- 实测：改前 `want = 0.376739502`、`l = 0.376739502`；改后 `want = 0`、`l = 0` ⇒ **绿**。
- **真代价＝见证力**：这枚夹具原本见证"含直流"，改后它见证"**0 等于 0**"（空转）。**处置建议（写进落点 5）**：不改期望值，**追加一条** `l != 0 → t.Fatalf`（恒定 DC 帧必须读恰好 0）＋把 `:23` 那句注释改成"ride a DC offset; the level must ignore it"。这样它在**回路层**（convert→downmix→resample→frame→emitLevel 全链）正面见证去直流，比原来更有牙。

### 3.2 `amplitudepair_297_windows_test.go`——**票面/派单没点名的第二枚必红钉**

- 夹具：`:59-65` 两窗**全是常量**（`lowCounts=1638`、`highCounts=13107`），`:28-32` 注释逐字写着选择理由「**DC makes the level exact** (RMS of a steady value v is |v|)」。
- 断言：`:122-124` `ratio ≈ ampRatio(8.0018)`（容差 0.01）、`:125-128` `meanLow ≈ 1638/32768 = 0.049987793`（容差 `SineLevelTolerance`）。
- 实测（`03-*`）：均值 `0.000000`，`:127` 报 `window A mean = 0.00000000, want 0.04998779` ⇒ **FAIL**。
- ⚠ **附带浮出一枚判据盲点（实测量到，不是推测）**：同一轮日志里 `mean ratio = NaN, amplitude ratio = 8.0018` **却通过了**——因为 `math.Abs(NaN−8) > 0.01` 是 `false`，**0/0 让比值断言静默变绿**。⇒ 落地时**必须**在除法前加 `meanLow <= 0 → Fatalf`（否则"尺读数归零"会被比值断言当成"分离良好"）。
- **改法（保住这枚夹具的用途，不是拆掉它）**：两窗改成 **同一 DC 偏置（`12345`）＋±A 方波**（`±lowCounts`／`±highCounts`）。推导：方波均值恰为 `dcOffset` ⇒ 残差 `±A` ⇒ RMS `= A` ⇒ 电平 `= A/32768`，与改前**同一个期望值**（`0.049987793` / ratio `8.0018`）；同时它**额外**见证"偏置不泄漏"。实测（`05-*`）：`mean=0.049988`／`0.399994`、`ratio = 8.0018`、**PASS**。
- 为什么必须改这个夹具而不是改期望值为 0：改成 0 等于**把"尺能分离已知幅度"这枚仪器拆了**（票 297 `AC#1` 的承重件）。重定信号保住仪器，还多一枚去直流见证。

---

## 4. 冻结面：**逐条 0 字节，且给可复跑尺**

| 冻结面 | 改法会不会碰 | 尺（本腿现跑）与读数 |
|---|---|---|
| `internal/observe/thresholds.go` | **不会**。本改法只在 `internal/audio` 内，且该文件与"电平"无关 | `grep -nE '[Ll]evel|RMS|audio' internal/observe/thresholds.go` ⇒ **零命中**；落地写面 `internal/audio/**` ⇒ 路径不相交 |
| golden（`internal/**/testdata/golden`、`internal/llm/golden`） | **不会**。全仓 golden 全是 **LLM SSE 回放**，没有一枚装电平数 | 现量 **4 个 golden 目录／56 枚文件**＝53 枚 `.sse` 夹具 ＋ `internal/llm/golden` 的 3 枚 go 源码（`golden.go`／`golden_test.go`／`replay.go`）；`grep -rn '0\.37\|12345' <golden 目录>` 零命中 |
| `tools/d22scan/allowlist.txt` | **不会** | `grep -nE 'audio\|level' tools/d22scan/allowlist.txt` ⇒ **零命中** |
| 三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`） | **不会**（路径不相交） | 本表零引用其内容；落地写面不含 `internal/panel|perm` |
| **C21 token `SilenceLevelGate = 0.06`**（`internal/ball/liquid.go:30`；机钉 `internal/ball/tokens_table_test.go:572` ＝ `"SilenceLevelGate": SilenceLevelGate`，在 `c21MotionGolden` 表 `:562-577` 里；出处 `docs/evidence/s1/c21-native-tokens.md:164`） | **数字一字不动**；但**单位/含义搬家**：改前 0.06 是"含直流的全波 RMS"，改后是"纯交流 RMS"。⇒ **这是要摆 owner 的那一格**，⛔ 不是写腿能顺手重定的 | `git diff -- internal/ball` 落地后必须为空；`go test ./internal/ball -run TestTokens` 必须绿。⚠ 票面/`291-a1` 写的机钉行 `:552` 已漂移（`:552` 现读是 `}`），现读 **`:572`** |
| D43 表 / `C1–C32` / `D1–D47` / `docs/PLAN.md` / `docs/specs/**` | **不会**（`AC#1` 已判：去 DC 是**实现细节、可自办**，不是契约变更） | 落地写面 `internal/audio/**` ⇒ 零路径相交 |
| 票 247 的 2× 门槛 `:194` | **不改**（归票 247，须非实现者裁） | 本表只登记，见 §7.4 |

**⇒ 回答派单第 4 问：不会碰冻结面（0 字节）。唯一"搬家"的是 `0.06` 的含义，属具名上交项，不属改法本身。**

---

## 5. 判据怎么钉：改前后逐名红名册（**实测，不是预测**）

尺＝`go test ./internal/audio/ -count=1 -v`，同一棵导出树（`acc1e9e23`）、同一台机器、连续三发：

| 发 | 树的状态 | 名册 | 红 |
|---|---|---|---|
| ① 基线 | `acc1e9e23` 原样 | 66 枚（64 PASS ＋ 2 SKIP：`TestLiveWasapiSmoke`／`TestAC297LivePacketByteFace`，均 `WISP_LIVE_MIC` 门） | **0** |
| ② 只改 `level.go` 求和段 | ＋去均值 | 69 枚（66＋3 枚新钉） | **2**：`TestLevelFullScaleSquareEndpoints`（`level_test.go:80`：`all-minus-32768 frame = 0, want exactly 1`）＋`TestAC297KnownAmplitudePairSeparates`（`amplitudepair_297_windows_test.go:127`：`window A mean = 0.00000000, want 0.04998779`） |
| ③ ②＋两枚夹具重定＋新钉 | 提案全量 | 69 枚 | **0**（`ok github.com/CarlosShao/wisp/internal/audio 18.577s`） |

日志：`01-baseline-audio-v.txt`／`03-after-levelgo-change-audio-v.txt`／`05-after-change-plus-fixture-rebaseline-audio-v.txt`；逐名作差 `06-roster-diff.txt`。
**同一份名册里 64 枚旧钉逐名仍是 PASS**（含 `TestLevelSineAgainstRootTwo`／`TestLevelIsLinearInAmplitude`／`TestLevelSineToleranceIsNamedAndBounded`／`TestLevelSameInputSameBits`／`TestLevelOverBoundedChannelFromWavInjector`／`TestAC247LevelSinkGetsOneScalarPerSeamFrame`／`TestAC247RealCaptureLoopEmitsLevels`／`TestAC297ConvertPacketSyntheticByteFace`）——`291-a1` 当时列为"读码定不了、必须真跑"的 **T1（三枚正弦/线性钉的实际颜色）与 T2（位级可重现性）在本腿量到了**：全绿，且 `TestLevelSameInputSameBits` 1000 次逐位相同。T3/T4 仍欠（T3＝真麦上的 2× 松紧；T4＝两趟遍历对 pinned 线程节拍的影响——本腿给到的是**算术代价** `425.8 → 978.2 ns/窗`（×2.30），`31.25 窗/s` 下**每秒多 17 µs**，远小于 `:190-193` 的 `[0.6,1.6]×31.25/s` 窗宽；真麦那一发仍归 247/AC#2 窗口）。

### 5.1 新增的钉（形状，逐条一句话＋为什么它咬得住）

1. **`TestLevelDCFrameReadsZero`**：`{0, 1, -1, 12345, 32767, -32768}` 六个恒定帧 → `LevelOfSamples` **恰好** `MinLevel`、`FrameLevel` **恰好** `float32(0)`、且不是负零。**改前必红**（实测：`constant 1 frame = 3.0517578125e-05`），因为`1 LSB` 的直流就够它报错。
2. **`TestLevelTranslationInvariant`**：取一枚整周期正弦 `s`，对 `c ∈ {−12345, −3000, −1, 1, 3000, 12345, 20000}` 断言 `LevelOfSamples(s+c)` 与 `LevelOfSamples(s)` **`math.Float64bits` 逐位相等**；再过 `FrameLevel` 用 `Float32bits` 比一次。**改前必红**（实测：`adding a constant DC of -12345 moved the level from 0.06473924235500596 to 0.38226145742487416`）。
3. **`TestLevelHottestLegalFrameStaysBelowMaxLevel`**：反相方波 `+32767/−32768` 读 **恰好** `32767.5/32768 = MaxLevel − LevelLSB/2`，且严格落在 `(FullScaleSquareLevel, MaxLevel)` 里。**改前改后都绿**（它是"新顶"的钉子，不是去直流的牙，别把它当见证）。

### 5.2 "直流被去掉"的**可判事实**（⛔ 不拿"看起来更低"当判据）

三条，任一可单独复跑：

- **F1 恒定性**：恒定帧（任何 `v`）电平 `== 0`，**恰好**（`float64 bits = 0x0000000000000000`；`float32 bits = 0x00000000`）。这是"直流不进平方和"的直接后果。
- **F2 平移不变性**：对任意帧与任意常数 `c`（不越 int16 界），`level(s) == level(s+c)` **位级相等**。这条是**充要**的"直流不参与"——任何"看起来更低""更平"的措辞都替代不了它。
- **F3 闭式**：`DC + 整周期正弦` 读 `A/√2/LevelFullScale`（实测 `DC 12345 + 正弦 6000`：改前 `0.398367696` → 改后 **`0.129476519`**，闭式 `6000/√2/32768 = 0.1294765`；旧读数把 `12345` 与 `6000` 正交叠加，新读数只留交流）。

---

## 6. 代价与风险（按铁律 ⓪⓪⓪：**只写代价与落点**，不据此选小方案）

### 6.1 会不会吃掉真实低频（实测扫频，worst/best phase over 128 相位，512 样 @16 kHz）

| 频率 | 最坏增益 | 最好增益 | 帧率纹波 |
|---|---|---|---|
| 5 Hz | −28.246 dB | −0.006 dB | 28.240 dB |
| 10 Hz | −15.646 dB | −0.005 dB | 15.640 dB |
| 15 Hz | −7.983 dB | −0.001 dB | 7.982 dB |
| 20 Hz | **−3.019 dB** | −0.000 dB | 3.019 dB |
| 25 Hz | −0.629 dB | −0.000 dB | 0.629 dB |
| 31.25 Hz | **0.000 dB** | 0.000 dB | 0.000 dB |
| 40 Hz | −0.294 dB | −0.000 dB | 0.294 dB |
| 50 Hz | −0.343 dB | −0.000 dB | 0.343 dB |
| 85 Hz（男声基频下限） | **−0.076 dB** | −0.000 dB | 0.076 dB |
| 120 Hz | −0.014 dB | −0.000 dB | 0.014 dB |
| 300 Hz | −0.009 dB | −0.000 dB | 0.009 dB |
| 440 Hz / ≥800 Hz | −0.000／−0.001 dB | 0 | ≤0.001 dB |

**读法**：① 语音基频（≥85 Hz）与共振峰（≥300 Hz）**≤0.08 dB**，听不出、球也看不出；② 被削的是 **≤15 Hz** 的"比一帧还慢"的东西（手持噪声、房间隆隆、ADC 漂移）——那正是要削的；③ **20 Hz 附近是最坏带**（最坏 −3 dB，且帧间纹波 3 dB）；④ **31.25 Hz（帧率的整谐波）恰好 0 dB**，所以常见工频干扰（50 Hz）只有 −0.34 dB 最坏值。⑤ **没有单一"拐点"**：这是**帧内去均值**，增益随"窗落在信号哪一段相位"变化——所以本表给的是**最坏/最好两列**，⛔ 不要引用一个"−3 dB 拐点 = 20 Hz"当唯一数字（严格说 20 Hz 是最坏相位下首次到 −3 dB 的频率）。

### 6.2 窗口多长

就是**一个 C8 缝帧**：`512 样 = 32 ms`（`audio.go:47-53`，`level_test.go:369-380` 钉着）。**没有第二个参数可调**——这是选 A 的额外价值：不新增"窗口长度"这种要有人值守的旋钮。若日后要更低的拐点，只能换 B（有状态），那不是本票。

### 6.3 对 `SilenceLevelGate = 0.06` 的影响方向（严格单调）

- **方向：只会更低，绝不更高。** 数学上 `var(s) = E[s²] − (E[s])² ≤ E[s²]`，实测在 5 组夹具上"改后 − 改前"的**最大增量 = +0.000e+00**（`arith_probe.out.txt`）。
- ⇒ 球侧那道门**只会更容易关上**：10-09 那种"安静读 0.368、门永远开着"的形状，去直流后读数落到交流底噪（11-08 同机实测 `0.003067`／`0.001305`），**门会正常关**。
- ⚠ **但 0.06 的物理含义搬家**（"含直流全波 RMS 的 6%" → "纯交流 RMS 的 6%"）。数字⛔ 不许动；这一句要**具名摆 owner**（`AC#2` 那一发同窗口）。
- ⚠ **反方向的新风险**：若某台机器的交流底噪本来就 >0.06（大声房间/高增益），门会**一直开着**——这与改前"永远开着"形状相同但成因不同。⇒ `AC#2` 四形读数必须带**安静窗 AC 读数**，不能只看"门关上了"。

### 6.4 对票 247 那枚 2× 门槛的影响（**只登记，不改**）

- `:194` `loud.mean() <= quietA.mean()*2`：去直流后 `quietA.mean()` 趋近交流底噪（实测 11-08 为 `0.0013–0.0031`，不是 0）⇒ **判据实际变松**（右边从 0.37 掉到 ~0.006，任何正常说话都过）。**这正是票面禁区那句"改尺让 AC#2 自动变绿"的形状**。
- ⇒ 落地票**必须**同时把这一格交给非实现者：要么换成**绝对量级**（如"有声窗 AC ≥ 某 dBFS"），要么明确它在去直流后**不再判事**。⛔ 本票不碰、⛔ 写腿不自定门槛。
- 附带：`:199-201` 的 `quietB.mean() > loud.mean()` 是**反向保护**，去直流后两静默窗差异变小 ⇒ **敏感度上升**（更容易抓到"这是噪声不是读数"），这是好事，也要在落地票里写上。

### 6.5 其它代价（逐条给量）

| 代价 | 量 | 备注 |
|---|---|---|
| `MaxLevel = 1.0` 变**不可达** | 最高合法帧 `= 1 − LevelLSB/2 = 0.9999847412109375`（反相方波 `+32767/−32768`） | `MaxLevel` 常数**不必改**（仍是闭区间上界）；但 `level.go:25-26` 那句「level 1.0 exactly: every sample sits at the largest magnitude an int16 can carry (-32768), which is the definition of LevelFullScale」**必须删/改**，并新增命名上界（建议名 `BipolarSquareLevel = 32767.5/32768`，或在注释里给闭式） |
| `MinLevel = 0` 的含义搬家 | 从"只有全零帧"变成"**任何恒定帧**" | `level.go:23-24` 那句「a frame reads 0 only when it IS 0」**必须改**，否则注释变假话（本仓最忌讳的形状） |
| 注释里的 int64 精确性主张 | `:89-91`「sum of squares accumulates in int64, which is exact」→ 平方和变成 `float64` | 必须改成"均值用 int64 精确求和；残差平方和按固定顺序 float64 累加，故同输入同比特"（实测 `TestLevelSameInputSameBits` 1000 次逐位同） |
| 每帧算术量 | `425.8 → 978.2 ns/窗`（×2.30），31.25 窗/s ⇒ **+17 µs/s** | 落在 pinned `audio-capture` 线程；`TestPinnedThreadStable10s` 实测仍绿、节拍不变 |
| 帧率纹波（<20 Hz 内容） | 20 Hz 处 3.0 dB、50 Hz 0.34 dB、≥85 Hz ≤0.08 dB | 球侧包络已有 attack/release 平滑（`liquid.go:240-250`），31 Hz 纹波看不出来；纯 DSP 层面它是 A 的真实短板，具名不掩盖 |
| 写腿要动的测试面 | 2 枚既有钉必红（§5）＋1 枚空转夹具＋3 枚新钉 | ⇒ 落地票的写面是 `internal/audio/**` 5 枚文件，**不含** `cmd/`、`internal/ball/` |

---

## 7. 落地票该跑的门禁（⛔ 本票不跑，只写清）

```
GOFLAGS= go build ./...
go test ./internal/audio/ -count=1 -v          # 期望：69/69（0 红；2 枚 SKIP 仍 SKIP）
go test ./internal/ball/ -count=1 -run 'TestTokens|TestLiquid' -v   # 0.06 未动 ⇒ 期望绿
sh scripts/d22scan.sh                          # 期望 clean（本改法无禁用形状）
gofmt -l internal/audio/... ; gofumpt -l ...   # 本腿实测：两把都空（10-format.txt）
git diff --stat -- internal/ball cmd/wisp docs/PLAN.md docs/specs tools/d22scan internal/observe   # 期望空
```
（本腿已在导出树实测 `gofmt -l`／`gofumpt -l` 对四枚动过的件**均为空**。）

**不落地**：以上全部只在 `D:\tmp\wisp\20261114\291-ac3-tree` 里验证过，**主仓 `internal/**` 零字节改动**。

---

## 8. 本腿的物证清单（都在 `.scratch/wisp/probes/291/ac3/`）

| 件 | 是什么 |
|---|---|
| `arith_probe.go` | 纯 stdlib 算术探针（基线 vs 去均值），⛔ 不 import 仓内包，⛔ 不进 `go build ./...`（`.scratch` 被 Go 工具链跳过） |
| `arith_probe.out.txt` | 它的全部读数（夹具、扫频、位检查、端点、溢出界、代价、单调性） |
| `01-baseline-audio-v.txt` | 改前 `internal/audio` 全名册（66 枚，0 红） |
| `02-newtooth-on-current-code.txt` | 新钉打**旧**代码：2 红 1 绿（红的是牙，绿的是"新顶"） |
| `03-after-levelgo-change-audio-v.txt` | 只改 `level.go` 后的全名册（**2 红**，逐名） |
| `04-newtooth-after-change.txt` | 新钉打**新**代码：3 绿 |
| `05-after-change-plus-fixture-rebaseline-audio-v.txt` | 提案全量（含两枚夹具重定）：**0 红** |
| `06-roster-diff.txt` | 三发名册逐名作差 |
| `07-proposed-levelgo-change.patch.txt` | 要落的那 5 行→10 行（逐字 diff） |
| `08-proposed-fixture-rebaseline.patch.txt` | 两枚夹具的改法 diff（含 `meanLow <= 0` 守门） |
| `09-proposed-new-nail.txt` | 新钉全文（3 枚，形状可直接照抄） |
| `10-format.txt` | `gofmt -l` / `gofumpt -l` 读数（均空） |

## 9. 写面与纪律自证

- 主仓 `git status --porcelain -- internal cmd docs/specs docs/PLAN.md tools internal/observe` 收工读数：**只有两枚别人的既有未跟踪件**（`?? tools/d22scan/check_blob_417.sh`、`?? tools/jarvis-export/.verify.js`，开工时就在，非本腿）；**`internal/**`／`cmd/**` 零改动**（对照开工时的同一把尺，读数一致）。
- 本件与 `.scratch/wisp/probes/291/ac3/**` 已随 commit **`8d98e4dfb`** 入盘（显式 pathspec；13 枚文件、1654 行插入）；`git show --stat 8d98e4dfb` 的 13 枚名册里 **`internal/**`／`cmd/**` 零命中**。⚠ 树里另有一枚**别人**的已改未跟踪/已改件（`.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go`，` M`）＝非本腿，**未入我的 commit**。
- ⛔ 未改票面／台账／`PLAN.md`／`specs`／`thresholds.go`／golden／`allowlist.txt`／C 表／D43；⛔ 未跑真麦（未设 `WISP_LIVE_MIC`，两枚真机台件仍 SKIP）；⛔ `add -A`／amend／reset／rebase／stash／checkout .；⛔ push。
- 临时件：导出树 `D:\tmp\wisp\20261114\291-ac3-tree`（仓外，只建不删）。它里面躺着**提案的产码改动**（`level.go` 去均值 ＋ 两枚夹具重定 ＋ 新钉），**只在导出树里存在**，主仓零字节。
