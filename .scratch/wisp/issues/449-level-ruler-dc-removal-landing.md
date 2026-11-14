# 票 449 — 电平尺去直流的**落地**（票 291 AC#3 的表施工）

**状态**：落地已交件（`f9b6d60f7`），**差真机验收**
**为什么另开一张**：票 291 的 AC#3 逐字写着「落点表（⛔ **本票不许落地**）」，AC#4 又写「落地程（若有）＝写面只许 `internal/audio/**` 与票 247 那一族台件」。⇒ 落地必须由**另一张票**承接。**编排者 2026-11-14 22:1x 把落地派在了 291 名下，越了 AC#3 这条线** —— 本票是补正，不是新需求。

## 交付物（已交件，非规划）

- 施工图：`docs/evidence/s1/291-ac3-landing-table.md`（30,329 B，票 291 AC#3 的产出）
- 落地：commit **`f9b6d60f7`**，证据 `docs/evidence/s1/291-level-dc.md`（4,343 B）
- 改法（施工图选定，⛔ 不改选）：`LevelOfSamples` **先减本帧均值再算 RMS**（两趟：int64 精确求和 → 第二趟累加 `(s-mean)²` 为 float64；**签名不变、无状态、零分配、无时钟、无 goroutine**）

## 判据（可判事实，全部实测）

- [x] **AC#1 恒定帧读恰好 0**：`{0,1,−1,12345,32767,−32768}` 六形 → `LevelOfSamples` 恰好 0、`FrameLevel` 恰好 `float32(0)`、**非负零**。改前 1 LSB 帧读 `3.0517578125e-05`（红）。
- [x] **AC#2 平移不变位级相等**：正弦 3000 加 `c∈{−12345,−3000,−1,1,3000,12345,20000}` ⇒ `Float64bits`／`Float32bits` **逐位相等**。改前 `0.06473924235500596 → 0.38226145742487416`（红）。
- [x] **AC#3 闭式**：`DC 12345 + 正弦 6000` 读 `0.1294765192` vs 闭式 `0.1294751186`（差 `1.4e-6`，在 `SineLevelTolerance=5e-5` 内）；改前 `0.3983676962`（红，差 0.269）。
- [x] **AC#4 新顶可判**：最热合法帧 `+32767/−32768` 读 `0.9999847412109375`（＝`MaxLevel−LevelLSB/2`）。
- [x] **AC#5 守门有阳性对照**：`amplitudepair_297` 加 `meanLow<=0` 守门。**它在今天不是红的**（`wantLow` 兜得住）⇒ 必须证明它不是装饰：把 `ampLow` 置 0，**带守门＝FAIL**（`window A mean = 0 … the ratio has no denominator`），**去掉守门＝PASS 且日志 `mean ratio = +Inf`**（`Inf−Inf=NaN` ⇒ 比较为 false ⇒ 静默通过）。
- [x] **AC#6 冻结面命中 0**：`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`docs/specs/` **一字节未动**；`SilenceLevelGate = 0.06` 那一行（在 **`internal/ball/liquid.go:30`**，不是 `thresholds.go`）**逐字仍在**、`liquid.go` diff = 0。
- [x] **AC#7 门禁**：`go test ./internal/audio` 四数；`go build ./...`／`go vet` rc=0；`gofumpt -l` 空；`d22scan` clean；path-length `VERDICT GREEN`。
- [ ] **AC#8 ⚠ 真机验收（要机主）**：修后同一枚真麦安静态读数必须**从 0.368 掉到交流底噪量级**（施工图实测 0.0013–0.0031），且 `SilenceLevelGate=0.06` 的判决方向**严格单向更低**（实测最大增量 `+0.000e+00`）。★这一步**没做**，也是本票唯一剩余项。

## 语义变更（必须有人知道，⛔ 数字没动）
- `MinLevel` 含义：从「只有全零帧」→「**任何恒定帧**」
- `MaxLevel = 1.0`：变**闭上界、无合法帧可达**（最高合法帧 `32767.5/32768`）
- `SilenceLevelGate = 0.06`：数字一字节未动，**含义更严格（更低）**
- 代价（实测）：20 Hz 最坏 −3.02 dB／15 Hz −7.98／31.25 Hz 0／85 Hz −0.076／≥300 Hz ≤0.009；每帧 `425.8→978.2 ns`

## 边界
写面：`internal/audio/**`（票 291 AC#4 允许）、`docs/evidence/s1/291-*`／`449-*`、`.scratch/wisp/probes/291/**`／`449/**`、两个票尾（只追加）。
⛔ `thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`docs/specs/**`；⛔ `internal/ball/**`（`0.06` 那行在 `liquid.go:30`，**不改**）；⛔ 票 247 那一族台件的门槛；⛔ `cmd/wisp/**`、`internal/panel/**`、`internal/llm/**`、`internal/tools/**`、`tools/d22scan/**`、`frontend/**`、`design/**`；⛔ 放宽任何容差；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。

## 撤销口令
「票 449 撤」＝删本票＋回退 `f9b6d60f7`（新 commit 反改，不 reset/amend）；「291 落地搬回 291」＝把本票内容并回 291 并把 AC#3 的「不许落地」改掉（那是改票面文字，要机主点头）。
---

## 编排者收件（2026-11-14 22:2x，非实现者）：AC#6／AC#7 翻；AC#1–#5 等审计腿；AC#8 等机主

- **尺1**：`f9b6d60f7` 在；落点＝`internal/audio/level.go`、`level_test.go`、`level_dc_291_test.go`（新）、`amplitudepair_297_windows_test.go`、`capturelevel_windows_test.go`。
- **尺2**：`docs/evidence/s1/291-ac3-landing-table.md` **30329 B** ＋ `291-level-dc.md` **4343 B**。
- **尺3 我自己复跑**：`go test ./internal/audio -count=1 -v` ⇒ **RC=0，`RUN=82 PASS=80 FAIL=0 SKIP=2`**，`ok … 18.974s`；两枚 SKIP 是 `TestLiveWasapiSmoke`／`TestAC297LivePacketByteFace`（**真麦/真包，要环境变量**，具名记）。
- **AC#6 我自己核**：`git show --name-only f9b6d60f7 | grep -E 'thresholds|golden|allowlist|PLAN.md|docs/specs/'` ⇒ **命中 0**。★并且我查清了一件票面没说准的事：**`SilenceLevelGate = 0.06` 不在 `thresholds.go`，在 `internal/ball/liquid.go:30`**（`thresholds.go` 实际在 `internal/observe/thresholds.go`）。那一行**逐字仍在**、`liquid.go` 零改动。
- **AC#7 我自己核**：包侧四数与上面同；`d22scan` clean、path-length `VERDICT GREEN` 腿已给，我另行复跑过 `d22scan` 与本仓 path-length。
- **AC#1–#5 我不翻**：那些是**具体数字声称**（恒定帧恰好 0／位级平移不变／闭式 1.4e-6／新顶／守门阳性对照），只有 4 枚新钉在盘上把它们钉住、由**实现者**写的。⇒ 按「翻框须非实现者」，交给审计腿逐条复跑。
- **AC#8 等机主**：真麦安静态读数必须从 **0.368** 掉到 **0.0013–0.0031** 量级 —— 这一步只有机主能提供（对着麦克风静一下）。

**腿自己报的三处与施工图偏差（我认，都往更可判一侧）**：① 表说第 3 枚钉"改前改后都绿"，**实测改前是红的**（`0.9999847413273546` vs `0.9999847412109375`），它把钉头注释改成实测事实并明写"它不是去直流的牙"；② 表列 3 枚钉它**多加第 4 枚**（闭式），因为派单判据点名了"闭式"；③ `amplitudepair` 顶部那句"DC makes the level exact"表没点、它也改了 —— 不改就成假话。
---

## 审计收件（2026-11-14 22:3x，`eb6c56bb8`）：五条闭合 **5／5**，AC#1–#5 翻

**尺**：审计提交在，非 probes 文件**只有 `docs/evidence/s1/449-adversarial.md`（11522 B）**，产码零改动；★**我自己复跑那枚"牙"**：`go test ./internal/audio -run '^TestLevelTranslationInvariant$'` ⇒ **PASS 0.00s**，定义在 `level_dc_291_test.go:69`。

| 格 | 审计读数（它自己现跑，同一份 overlay 分别注入 HEAD 尺与 `f9b6d60f7^` 尺） | 判 |
|---|---|---|
| AC#1 | 六形全 `LevelOfSamples=0`（bits `0000000000000000`）、`FrameLevel=0`（f32bits `00000000`、Signbit false）；旧尺 `1`→`3.0517578125e-05` | 满足 |
| AC#2 | base `0.064739242355005958`（`3fb092c040aebcd1`），7 个平移 `Float64bits`/`Float32bits` 全同；旧尺 `±1` 就换位（`…5f936312`）、`−12345`→`0.38226145742487416` | 满足 |
| AC#3 | `0.12947651924196049` vs 闭式 `0.12947511862546646`，差 `1.401e-06`（容差 `5e-5`）；旧尺 `0.39836769618896367`（差 `0.2689`） | 满足 |
| AC#4 | 最热帧恰好 `0.9999847412109375`＝`MaxLevel−LevelLSB/2`；解析界 `|s−mean|≤32767.5` ⇒ level `<1.0`，**20 万随机帧搜索无帧 ≥1.0** | 满足 |
| AC#5 | `ampLow=0`：带守门 **FAIL**（`the ratio has no denominator`）；去掉守门 **PASS** 且 `mean ratio = +Inf`；`|+Inf−+Inf|=NaN`，`NaN>0.01`=false | 满足 |

### ★本仓第一次把「哪一枚钉真的咬住」做成实验（不是论断）
审计用**两枚假尺**去试四枚钉：
- **MUTANT A**（恒定帧特判返回 0、其余照旧数 DC）⇒ **骗过第 1 枚**
- **MUTANT B**（只去掉 99.99% 的均值，DC 12345 仍留 1.2345 counts）⇒ **骗过第 3 枚的闭式半边**（`1.406e-06 ≤ 5e-05`，只有它后半段的位级等式才红）**并且把第 4 枚整枚骗绿**
- **只有第 2 枚 `TestLevelTranslationInvariant`（`level(s)==level(s+c)`）在两枚假尺上都红** —— 零容差、逐位比较

⇒ 这正是实现腿说「第 3 枚改前其实是红的、但它不是去直流的牙」的**成因**。**此后凡"去直流/去偏置"类改动，牙必须钉在平移不变上，不许只钉恒定帧特判或闭式近似。**

### 另三件（审计现核）
- 冻结面：`f9b6d60f7` 的 13 个文件里 frozen-name 命中 **0**、冻结集 numstat **0 行**；`liquid.go`／`thresholds.go`／`allowlist.txt`／`docs/PLAN.md` **三态同 blob**；`docs/specs` 14 件 tree 相同
- `internal/ball/liquid.go:30` ＝ `SilenceLevelGate   = 0.06 // …` **逐字仍在**
- `SineLevelTolerance` 改前后都是 `5e-5`，三枚正弦误差 bits 改前改后**完全相同**

### 剩下唯一一格
**AC#8（真麦验收）未做** —— 修后安静态必须从 **0.368** 掉到 **0.0013–0.0031** 量级。已并入「一次真机会话跑本」（`w-runbook`），等机主一次 3 分钟。
