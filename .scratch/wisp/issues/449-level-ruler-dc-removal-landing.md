# 票 449 — 电平尺去直流的**落地**（票 291 AC#3 的表施工）

**状态**：落地已交件（`f9b6d60f7`），**差真机验收**
**为什么另开一张**：票 291 的 AC#3 逐字写着「落点表（⛔ **本票不许落地**）」，AC#4 又写「落地程（若有）＝写面只许 `internal/audio/**` 与票 247 那一族台件」。⇒ 落地必须由**另一张票**承接。**编排者 2026-11-14 22:1x 把落地派在了 291 名下，越了 AC#3 这条线** —— 本票是补正，不是新需求。

## 交付物（已交件，非规划）

- 施工图：`docs/evidence/s1/291-ac3-landing-table.md`（30,329 B，票 291 AC#3 的产出）
- 落地：commit **`f9b6d60f7`**，证据 `docs/evidence/s1/291-level-dc.md`（4,343 B）
- 改法（施工图选定，⛔ 不改选）：`LevelOfSamples` **先减本帧均值再算 RMS**（两趟：int64 精确求和 → 第二趟累加 `(s-mean)²` 为 float64；**签名不变、无状态、零分配、无时钟、无 goroutine**）

## 判据（可判事实，全部实测）

- [ ] **AC#1 恒定帧读恰好 0**：`{0,1,−1,12345,32767,−32768}` 六形 → `LevelOfSamples` 恰好 0、`FrameLevel` 恰好 `float32(0)`、**非负零**。改前 1 LSB 帧读 `3.0517578125e-05`（红）。
- [ ] **AC#2 平移不变位级相等**：正弦 3000 加 `c∈{−12345,−3000,−1,1,3000,12345,20000}` ⇒ `Float64bits`／`Float32bits` **逐位相等**。改前 `0.06473924235500596 → 0.38226145742487416`（红）。
- [ ] **AC#3 闭式**：`DC 12345 + 正弦 6000` 读 `0.1294765192` vs 闭式 `0.1294751186`（差 `1.4e-6`，在 `SineLevelTolerance=5e-5` 内）；改前 `0.3983676962`（红，差 0.269）。
- [ ] **AC#4 新顶可判**：最热合法帧 `+32767/−32768` 读 `0.9999847412109375`（＝`MaxLevel−LevelLSB/2`）。
- [ ] **AC#5 守门有阳性对照**：`amplitudepair_297` 加 `meanLow<=0` 守门。**它在今天不是红的**（`wantLow` 兜得住）⇒ 必须证明它不是装饰：把 `ampLow` 置 0，**带守门＝FAIL**（`window A mean = 0 … the ratio has no denominator`），**去掉守门＝PASS 且日志 `mean ratio = +Inf`**（`Inf−Inf=NaN` ⇒ 比较为 false ⇒ 静默通过）。
- [ ] **AC#6 冻结面命中 0**：`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`docs/specs/` **一字节未动**；`SilenceLevelGate = 0.06` 那一行（在 **`internal/ball/liquid.go:30`**，不是 `thresholds.go`）**逐字仍在**、`liquid.go` diff = 0。
- [ ] **AC#7 门禁**：`go test ./internal/audio` 四数；`go build ./...`／`go vet` rc=0；`gofumpt -l` 空；`d22scan` clean；path-length `VERDICT GREEN`。
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
