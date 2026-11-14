# 票 449 对抗验收（审计腿，非实现者）— AC#1–AC#5 五条数字声称

**被审**：commit `f9b6d60f7`（parent `ead1ee91f`），票 `.scratch/wisp/issues/449-level-ruler-dc-removal-landing.md`。
**审计腿**：另一条腿，**只核不改**。**钟点** 2026-11-14 22:18–22:4x +08；**树锚** 开工 `be6d50245`、收工 `116a2efff`（审计期间编排者提交了 runbook；`git diff --stat f9b6d60f7 HEAD -- internal/audio` 为空 ⇒ 下面的读数仍对应被审树）。
**方法**：⛔ 未抄 `docs/evidence/s1/291-level-dc.md`，⛔ 未抄 `.scratch/wisp/probes/291/land/*`。自建 overlay 台件：**同一份**读数探针（`.scratch/wisp/probes/449audit/zz_449_audit_test.go`）经 `go test -overlay` 注入 `package audio`，分别对 **HEAD 的 `level.go`**、**`f9b6d60f7^` 的 `level.go`**、**两枚 mutant 尺**跑；产码零改动（审计后 `git status --porcelain -- internal/` 为空）。
**命令约定**（全部在仓根跑）：`NEW` = `-overlay .scratch/wisp/probes/449audit/overlay-add-test.json`（HEAD 尺）；`OLD` = `-overlay .scratch/wisp/probes/449audit/overlay-old-level.json`（`f9b6d60f7^` 尺）。原始读数：同目录 `01-`…`11-*.txt`。

## 五条裁决（一行五列）

| # | 声称（票面原文） | 我的读数（自己现跑） | 判 | 一句依据 | 可核命令 |
|---|---|---|---|---|---|
| AC#1 | 恒定帧 `{0,1,−1,12345,32767,−32768}` → `LevelOfSamples` 恰好 0、`FrameLevel` 恰好 `float32(0)`、**非负零**；改前 1 LSB 帧读 `3.0517578125e-05` | 六形全部：`LevelOfSamples=0`（bits `0000000000000000`）、`FrameLevel=0`（f32bits `00000000`）、`Signbit(f64)=Signbit(f32)=false`；空窗也 +0。改前：`1`/`−1` = `3.0517578125e-05`（bits `3f00000000000000`／f32 `38000000`）、`12345` = `0.376739501953125`、`32767` = `0.999969482421875`、`−32768` = `1` | **满足** | 零是 **+0 的位形**（不是 `−0` 的 `8000…`）；同一探针在旧算术下逐字复现票面的 `3.0517578125e-05` | `go test ./internal/audio -run Test449AC1 -count=1 -v NEW`；对照 `… -run Test449AC1 -count=1 -v OLD`（`01-`/`02-`） |
| AC#2 | 正弦 3000 加 `c∈{−12345,−3000,−1,1,3000,12345,20000}` ⇒ `Float64bits`／`Float32bits` **逐位相等**（⛔ 用 `==` 了事）；改前 `0.06473924235500596 → 0.38226145742487416` | base = `0.064739242355005958`（bits `3fb092c040aebcd1`）；7 个平移**全部** `Float64bits` 相同（同一 `3fb092c040aebcd1`），`FrameLevel` 全部 `3d849602`。旧算术：7 个**全部不同**——`c=−12345` → `0.38226145742487416`（`3fd876f8c28a4ff8`）、`c=±3000` → `0.1121297136009876`、`c=±1` → `0.064739249547880712`（`…5f936312`）、`c=+20000` → `0.61377536554255063` | **满足** | 逐位比较、不是 `==`；连 `±1` 这种 1 LSB 平移都换了低位（`…40aebcd1`→`…5f936312`），旧尺确实把 DC 算进去了 | `go test ./internal/audio -run Test449AC2 -count=1 -v NEW`；对照 `… OLD`（`01-`/`02-`） |
| AC#3 | `DC 12345 + 正弦 6000` 读 `0.1294765192` vs 闭式 `0.1294751186`，差 `1.4e-6`；改前 `0.3983676962` | 读 `0.12947651924196049`（＝`0.1294765192`）；闭式 `0.12947511862546646`；差 `1.401e-06`（`0.0397 × SineLevelTolerance`）；同正弦去/带 DC 位级相同（`3fc092afc3df3fd1`）。改前 `0.39836769618896367`（差 `0.2689`） | **满足** | 差值 `1.401e-06` 是**我实测**，不是抄的 `1.4e-6`；且 `withDC` 与 `withoutDC` 位级相等，说明不是两个邻近值的巧合 | `go test ./internal/audio -run Test449AC3 -count=1 -v NEW`；对照 `… OLD`（`01-`/`02-`） |
| AC#4 | 最热合法帧 `+32767/−32768` 读 `0.9999847412109375`（＝`MaxLevel−LevelLSB/2`），且 `MaxLevel=1.0` **无合法帧可达** | 读 `0.9999847412109375`，与 `(32767+32768)/2/32768` **恰好相等**；`FrameLevel` = `3f7fff00`；`MaxLevel` 未达；极值族 + 20 万随机帧搜索 max 仍是该值（无帧 ≥ `1.0`）；`MaxLevel−LevelLSB/2` 差 `1.52587890625e-05`。旧算术 `0.9999847413273546` | **满足** | 可达性有**解析界**：任一 `s` 与均值都在 `[−32768,32767]` ⇒ `\|s−mean\|≤32767.5` ⇒ `level≤32767.5/32768<1.0`（`32767.5` 可表示、除 `32768` 是 2 的幂精确） | `go test ./internal/audio -run Test449AC4 -count=1 -v NEW`（`01-`/`02-`） |
| AC#5 | `amplitudepair_297` 的 `meanLow<=0` 守门**不是装饰**：`ampLow` 置 0 ⇒ 带守门 **FAIL**、去掉守门 **PASS** 且日志 `mean ratio = +Inf` | `ampLow=0`：**带守门 FAIL** `window A mean = 0: the ruler read silence for a 0-count square, so the ratio has no denominator`；**去掉守门 PASS**，日志 `mean ratio = +Inf, amplitude ratio = +Inf`（window B mean `0.399994`）。前提算术：`\|+Inf−+Inf\| = NaN`（bits `7ff8000000000000`），`NaN > 0.01` = **false** | **满足** | 两条读数用**同一份**测试逐字改 `ampLow`（第二条只多删守门块），PASS/FAIL 的分界就是守门本身；`NaN` 比较为 false 是实测，不是推断 | `go test ./internal/audio -run TestAC297KnownAmplitudePairSeparates -count=1 -v -overlay .scratch/wisp/probes/449audit/overlay-ampLow0-with-guard.json`；`… overlay-ampLow0-without-guard.json`（`03-`/`04-`/`11-`） |

**五条闭合 5／5。**

## 结论：4 枚新钉里，**咬住去直流的是 `TestLevelTranslationInvariant`（第 2 枚，位级平移不变）**

判法不是"哪枚改前红"（改前四枚全红，见下 §三），而是"**哪枚在'DC 没被真正去掉、只是看起来像去掉了'的尺上仍然红**"。我造了两枚 mutant 尺，用**真钉**去咬它们（`05-pins-vs-MUTANT-A/B.txt`）：

| 钉 | HEAD（去均值） | OLD（数 DC） | MUTANT A：恒定帧特判 0，其余仍数 DC | MUTANT B：只去掉 99.99% 的均值（DC 12345 留 1.2345 counts 未去） |
|---|---|---|---|---|
| 1 `TestLevelDCFrameReadsZero` | 绿 | 红（`3.0517578125e-05`） | **绿（被骗过）** | 红（`3.0517578124996788e-09`） |
| 2 `TestLevelTranslationInvariant` | 绿 | 红（`0.0647…→0.3822…`） | **红** `adding a constant DC of -12345 moved the level… the scale still counts DC` | **红**（`0.06473924235500596 → 0.06473925331687663`） |
| 3 `TestLevelDCPlusSineReadsClosedFormRMS` | 绿 | 红（`0.3983… vs 0.1294…`，差 `2.689e-01`） | 红（同左） | 红在**位级等式**那条：读 `0.1294765247` 与闭式差 `1.406e-06`，**容差半边 PASS**；`0.1294765192419605`（无 DC）vs `0.1294765247229794`（带 DC）才红 |
| 4 `TestLevelHottestLegalFrameStaysBelowMaxLevel` | 绿 | 红（`0.9999847413273546`） | 红（旧算术） | **绿（被骗过）** |

三条判定，逐条点名：

1. **第 1 枚不是牙**：它只要求"恒定帧读 0"，MUTANT A 用一个 `if constant {return 0}` 特判就骗过它（绿），而同一条尺在正弦上仍把 DC 算进去。⇒ 必要、不充分。
2. **第 4 枚不是牙（与腿自己的说法一致，且我量出了原因）**：MUTANT B 留着 `1.2345 counts` 的 DC，第 4 枚**照样绿**；它改前的红是 `−0.5` 均值在末位（`0.9999847413273546` vs `0.9999847412109375`），不是"DC 去没去"的判据。
3. **第 3 枚的"闭式"半边也不是牙**：它的断言带 `SineLevelTolerance=5e-5`，折合 `1.64 counts` 的 RMS；MUTANT B 留 `1.2345 counts` 的 DC，闭式半边**PASS**（`1.406e-06 ≤ 5e-05`）。真正咬住 MUTANT B 的是它后半段的**位级等式**——而那正是第 2 枚的命题。⇒ 第 3 枚是旁证，不是牙。

**只有第 2 枚 `TestLevelTranslationInvariant` 在两枚假去直流尺上都红，且零容差、逐位比较、命题本身就是定义**（`level(s) == level(s+c)` 位级成立 ⇔ 窗口自己的均值不参与）。**它咬住去直流。** 一句话给编排者：**咬住去直流的那枚钉叫 `TestLevelTranslationInvariant`（`level_dc_291_test.go:69`）**。

## 另核三件

### 一、冻结面：0 命中，`liquid.go:30` 逐字仍在（`10-frozen-surface.txt`）

| 查什么 | 读数 |
|---|---|
| `git show --name-only f9b6d60f7` 的 13 个文件里含 `thresholds/golden/allowlist/PLAN.md/docs/specs/` 的 | **0** |
| `git diff --numstat f9b6d60f7^ f9b6d60f7 -- internal/observe/thresholds.go internal/ball docs/PLAN.md docs/specs tools/d22scan/allowlist.txt 全部 golden 目录` | **无输出**（0 文件 0 行） |
| 同一冻结集，**当前工作树 vs `f9b6d60f7`** | **无输出**（提交后没人碰过） |
| blob 级（EOL/BOM 无关）`liquid.go` / `thresholds.go` / `allowlist.txt` / `docs/PLAN.md` | parent＝commit＝worktree：`854caf43…`／`e2677b11…`／`6b61fad5…`／`92690666…` |
| `docs/specs` 14 个跟踪文件 | parent 与 commit 的 tree **完全相同** |
| `internal/ball/liquid.go:30` | `SilenceLevelGate   = 0.06 // envelope below this counts as "nobody speaking"` —— **逐字仍在**，`liquid.go` 三态同 blob |
| 工作树 `git status --porcelain -- internal/ docs/specs docs/PLAN.md` | 空（唯一杂物：`tools/d22scan/check_blob_417.sh` 是**未跟踪**文件，非本次改动，照实记） |

### 二、容差没被放宽：`5e-5` 三处读出同值，三枚正弦误差**改前改后逐位相同**

- `SineLevelTolerance` 在 `f9b6d60f7^` 的 `level.go` 与 HEAD 的 `level.go` **都是 `5e-5`**；`git diff … -- internal/audio/level.go` 里这一行只作为**上下文**出现，没有任何 `+`/`-` 行碰到 `SineLevelTolerance`/`5e-5`（`10-frozen-surface.txt` §[7]）。
- 三枚正弦误差（`level_test.go TestLevelSineAgainstRootTwo` 的那三枚），同一探针分别跑 HEAD 尺与旧尺，**误差 bits 完全一致**（`08-`/`09-`）：

| amp | 误差 | bits（HEAD） | bits（OLD） | 位级相同 |
|---|---|---|---|---|
| 32767 | `2.31318982615302815e-05` | `3ef8416bf1c70000` | `3ef8416bf1c70000` | 是 |
| 16384 | `1.12668913038183049e-06` | `3eb2e717ccbc0000` | `3eb2e717ccbc0000` | 是 |
| 8192 | `2.26521395810252990e-06` | `3ec300828b0d0000` | `3ec300828b0d0000` | 是 |

### 三、新钉是不是牙：改前**四枚全红**（含腿自认"第 3 枚改前其实是红的"那枚）

旧尺跑**真钉**（`06-pins-vs-OLD.txt`）：四枚全红，读数与票面/腿的披露逐字对上——

- 钉 1：`constant 1 frame = 3.0517578125e-05, want exactly 0`
- 钉 2：`adding a constant DC of -12345 moved the level from 0.06473924235500596 to 0.38226145742487416`
- 钉 3：`DC 12345 + sine 6000 reads 0.39836769618896367, want 0.12947511862546646 … (off by 2.689e-01)`
- 钉 4：`hottest legal frame = 0.9999847413273546, want exactly 0.9999847412109375` ← **腿自己说"表里写改前改后都绿、实测改前是红的"的那枚，我复核为红，腿的更正成立**；且它确如腿所说**不是去直流的牙**（§结论第 2 条已量出反证）。
- HEAD 尺跑真钉：`07-pins-HEAD.txt` 四枚全绿。

## 边界与未做

- **未改任何产码**：`internal/audio/**` 零改动（`git status --porcelain -- internal/` 空）；票面/台账/阈值/golden/allowlist/契约一字未动；未放宽任何容差；未 `push`。
- **AC#8 真机不在本审计射程**（要机主对着麦克风静一下），本表**不翻**它。
- 我造的两枚 mutant 尺与旧尺**只存在于 `.scratch/wisp/probes/449audit/`**，经 `-overlay` 注入，从未落进 `internal/audio/`；它们的用途是"用真钉咬假尺"，不是候选实现。
- **五条判据全部自己现跑**；原始 stdout 在同目录 `01-`…`11-*.txt`，命令见上表第 5 列。
