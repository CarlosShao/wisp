# 65 — glass 实现腿（W65-impl · wensity 带状款 + style 接缝 + 音频伪谱 v1）

**Status:** impl-complete（本腿交付；**不做 jarvis 渲染**——下一腿）
**Leg:** W65-impl（wave42；规则 21 隔离树内实现，主树字节同 patch 落地）
**HEAD 锚:** 主树开工 `199ab3def`（派单写 `bca0990f8+`，实际开工时 dev tip 已收 68 裁决笔）
**方向依据:** 台账 **A1581**（owner 批准 `[ball].style` 进 D36 键清单）+ `docs/evidence/s1/65-style-direction.md`（§1.1 量化值全数采纳；§2.0 接缝；§3.1 v1 伪谱决策）
**票面:** `.scratch/wisp/issues/65-ball-glass-quality-rework.md`（本腿**未编辑票面**——写面授权只有 internal/ball/** + PLAN.md 一行 + 本件）
**写面核实:** `internal/ball/**` · `internal/config/**`（票 65 派单第 1 条授权）· `cmd/balldebug/main.go`+`diff_windows.go`（`-style` 旗标，§4.3 预告项）· `docs/PLAN.md:2741`（仅 [ball] 键清单一行）· `docs/evidence/s1/c21-native-tokens.md`（token 行登记——派单"add token rows"授权；不写则双胞胎机器必红）· 本件
**没碰:** `internal/audio`（C8 缝零接触）· `frontend/**` · `design/**` · `docs/reports/pending-and-issues.md` · 票面 · `docs/specs/**`（SPEC-03 [ball] 节文字同步＝**移交项**，见 §6）· `internal/panel/**` · 零 push

---

## 0. 交付一览（对派单六条）

| 派单项 | 落点 | 状态 |
|---|---|---|
| (1) CONFIG `[ball] style` | `internal/config/schema.go`（`BallSection.Style`，tag default `glass`）+ `validate.go::validateBall`（枚举 `glass\|jarvis`，未知值报错带词表）+ `schema.go` 词表 `ballStyles` | ✅ |
| (1a) PLAN.md 一行 diff | `docs/PLAN.md:2741`＝`` `size` **`style`** `position{x,y,monitor}` ... ``（A1581 批准的原文） | ✅ |
| (2) API `SetStyle` | `internal/ball/tokens.go`：`Style` 枚举（`StyleGlass`/`StyleJarvis`）· `SetStyle`/`ActiveStyle`/`StyleNames`（仿 `SetLook` :328 现形）；`balldebug -style` + `-diff` 透传 | ✅ |
| (3) GLASS/WENSITY look | 新 look `wensity` ＝ `looks[0]`＝**默认**；量化值逐条落 token（§2）；渲染＝新 `wensity_windows.go`（带 + 白核 + 顶弧 + 非对称唇 + 底辉光），旧 4 套 look 保留为 glass 子 treatment | ✅ |
| (4) AUDIO-REACTIVE v1 | `liquid.go`：标量缝内派生伪谱 3 带（low τ400ms / mid＝level / high τ40ms 上升率）；`bandDeformState` 三态白名单（Listening/Thinking/Speaking）；满幅形变 token `WensityBandDeformFrac=0.15`；**internal/audio 零接触** | ✅ |
| (5) C21 双胞胎 | token 行全部入 `c21-native-tokens.md`；`TestC21TableColourRowsMatchCode` 扩展 `wensity.<Field>` 引用形**双向**；`TestC21GeometryRowsMatchCodeConstants` 扩展金标准（11 枚 Wensity* + 3 枚 Spec*）；`TestNoHardcodedColorsInBallPackage` 绿（新文件零字面量） | ✅ |
| (6) D32 HARD | `AnimationPolicy` **一字未动**；Sleeping 静帧＝带子零形变零漂移（`TestWensityRestingFrameIsStatic` + `TestStyleSwitchNeverWakesSleeping` 钉死）；零新定时器/线程/COM 常驻（带 brush 全部在 `setAssets` 一次烘焙，帧内只 fill/stroke） | ✅ |

## 1. 本腿改动逐文件

### 1.1 `internal/ball/tokens.go`
- `looks` 表头插入 **wensity**（第 5 套、`looks[0]`＝`DefaultLook()`）：9 字段全为 §1.1 实测色（BlobA=`#C2162F` 品红侧镜像、Rim=`#1B1D24` 暗缘、Lip=`#F7FAFF`、Caustic=`#53565E`、Glow=`#0E35A8` 深蓝渗光）；blob 三色仅为共用管线满形（wensity 分支不烘焙不绘制 blob，见 1.3）。
- 新 `WensityBand` 结构（9 Color 字段：6 带端色 + SpecArc + RimLip + BottomGlow）+ 全量值 + `WensityBandTokens()`；`wensityStopTs`（六色沿 span 的位置 0.06/0.22/0.46/0.66/0.84/0.98，未导出＝单 treatment 内部布局）。
- 新 11 枚几何 token：`WensityBandCenterFrac 0.13` · `WensityBandWidthFrac 0.39` · `WensityBandTiltFrac 0.33` · `WensityBandDeformFrac 0.15` · `WensityBandSpanLeftFrac 0.84` · `WensityBandSpanRightFrac 0.68` · `WensitySpecArcRadiusFrac 0.95` · `WensitySpecArcThickFrac 0.04` · `WensitySpecArcSpanDeg 57` · `WensityRimLipFrac 0.97` · `WensityBottomGlowFrac 0.75` —— 全部 §1.1 数字，全部进 C21 几何表。
- `Style` 枚举 + `SetStyle`（glass/jarvis 两值收下，未知 false；**jarvis 本腿收值不渲染**——落玻璃帧，注释与测试都写明，等下一腿）。

### 1.2 `internal/ball/wensity.go`（新，无 build tag）
纯函数四件：`wensityBandAxisY`（spine＝CenterFrac + TiltFrac·(t−0.5)，两端高差恰 0.33R）、`wensityBandAxisX`（span 左 0.84 右 0.68）、`wensityBandDeform`（clamp 后 DeformFrac×mid，单调、0 起点）、`wensityBandSegmentColor`（最近 stop，**不插值**——只画表里声明过的纯 token 色）。包内零 import，任意 GOOS 可测（同 liquid.go 姿态）。

### 1.3 `internal/ball/wensity_windows.go`（新，windows）
- `wensityBandSegs=16` 段预烘焙（同 `liquidLadder` 法则：帧内零 COM 分配；`setAssets` 一次性建 brush，`dropAssets` 释放）。
- `drawWensityBand`：底辉光 → 16 段带斑（spine + 形变偏移 `def`、半宽随 mid 胀）→ 白核。色谱沿**既有** `LiquidAngle` 横向漂移（brush 索引旋转；Sleeping angle=0 ⇒ 定格静帧）。
- `drawWensityEdges`：右下 15°–75° 亮唇弧（0.97R，非对称折射——其余三向走 `lk.Rim` 暗缘全圆）+ 顶部白弧高光（0.95R、厚 0.04R、跨 57°）。
- `renderer_windows.go` 钩子：`drawGlass` step 3 分支（wensity 画带、blob look 画场，blob 代码一字未动）、step 6b 分支、`setLiquid` 提前烘焙 caustic 后对 wensity 跳过 blob、struct 加 3 brush 槽。

### 1.4 `internal/ball/liquid.go` + `statevisual.go` + `liquid_windows.go`
- 伪谱三带：`SpecLowTauMs 400` / `SpecHighTauMs 40` / `SpecHighRateGain 8`；`specFollow` 一阶积分纯函数；`pushLevel` 同一步派生（band-mid＝`m.level` 本身，零新状态）；`endSession` 归零；`motionSnapshot/changed` 纳入两带（量化阈值同 `quantChanged`）。
- `Visual` 加 `SpectrumLow/Mid/High`；`applyTo(v, deform bool)`——**只有 `bandDeformState(s)`＝{Listening, Thinking, Speaking} 才盖章谱带**；Conversation/Warm/Acting 虽在 `liquidDriven` 白名单吃音频流，带保持静帧（派单"only Listening/Speaking/Thinking deform"的落点）；反归纳钉 `TestWensityDeformGate`。
- `liquid_test.go` 三处 `applyTo` 调用点机械加 `false` 实参（**断言一字未动**——票 65 行为不动承诺）。

### 1.5 `internal/config/**` + `cmd/balldebug/**`
- `schema.go`：`Style string \`toml:"style" default:"glass"\`` + 词表 `ballStyles`（`BallStyleGlass`/`BallStyleJarvis` 常量；注释写明 config 不 import ball 的理由）。
- `validate.go::validateBall`：枚举校验，报错带 `glass|jarvis` 全词表（同层 enum 规矩：拒绝 ≠ 静默回落）。
- `balldebug`：`-style` 旗标（未知值 exit 2 带词表）+ `diffOpts.style` 子进程透传（对拍两风格时 `-diff -style …` 可用）。

### 1.6 C21 表（`docs/evidence/s1/c21-native-tokens.md`）
- look 表 4→5 套（+9 行 wensity，prose 同步"9 字段 × 5 套、looks[0]=wensity"）。
- 新节「票 65 wensity 带状款 token」9 行（`wensity.<Field>` 引用形）。
- 几何表 +3 行（11 枚 Wensity* + 3 枚 Spec*，数字逐枚单写——同值互相掩护的坑票 74 已封）。
- 审计节 +1 条登记。

## 2. 量化值对照（§1.1 → 本腿落点）

| 量化属性（65-style-direction §1.1） | 落点 |
|---|---|
| 光带垂直位置 0.13R 上方 | `WensityBandCenterFrac` → `wensityBandAxisY(0.5)` |
| FWHM 0.39R 高斯软落 | `WensityBandWidthFrac`（段椭圆 ry=0.39R·0.62×高斯梯度） |
| 冷端上翘 0.33R | `WensityBandTiltFrac` → `axisY(1)−axisY(0)`（测试逐值钉） |
| 色谱 red→magenta→white→blue→cyan→orange | `WensityBand` 六 stop + `wensityStopTs`，次序由 `TestWensityBandGeometryPinsTheQuantifiedSpec` 钉死 |
| 白核 (−0.08R,−0.06R) | `bandCore` brush 中心（spine 上方 t=0.5 处 x−0.08R、y+0.07R） |
| 顶部白弧 0.95R / 厚 0.03–0.05R / 横向 ~0.9R | `WensitySpecArcRadiusFrac/ThickFrac/SpanDeg(57°)` |
| 右下 0.97R 亮唇 `#F7FAFF`，三向暗缘 | `WensityRimLipFrac` + `wensity.RimLip`；暗缘＝`looks[wensity].Rim #1B1D24` |
| 底部内辉光 0.75R `#53565E` | `WensityBottomGlowFrac` + `wensity.BottomGlow` |
| 体色 `#05050E–#1B1B1E` 近融入背景 | **保留既有玻璃壳**（`pal.OrbBody/OrbBody2`，浅色壁纸走 LightPalette 白玻璃——§1.1 对 Wisp 的推论），不引暗体 |
| 满幅形变 pp≈0.15R | `WensityBandDeformFrac`（测试钉 full-scale＝token 值） |

## 3. 门禁读数（全部出自克隆树 `D:\tmp\wisp\20261106\65glass\repo`，改后现跑；环境＝Go 1.27.1 windows/amd64、CGO 未清零、GOTMPDIR/TMPDIR/TEMP/TMP unset、PATH 前插 build+sherpa-onnx junction）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `gofumpt -l internal/ball`（另 internal/config、cmd/balldebug） | **空**（首跑 4 文件报出 → `-w` 修后空；v0.12.0） |
| gofmt | `gofmt -l internal/ball internal/config cmd/balldebug` | **空** |
| vet | `go vet ./internal/ball ./cmd/balldebug ./internal/config` | **clean**（首跑抓出两处：`Style()` 函数与 `Style` 类型同名冲突→getter 改 `ActiveStyle()`；`wensity.go` 无类型常量默认 float64→改常量算式。均已修） |
| build | `go build ./cmd/balldebug` | **rc=0** |
| test | `go test ./internal/ball -short -count=2` | **ok 0.334s**（首跑两红：style 全局残留→两测试加 save/restore；WensityBand 行数我误记 10 实为 9→改判据。均属新测试自身，非行为断言） |
| test | `go test ./cmd/balldebug -count=2` | **ok 0.052s** |
| test | `go test ./internal/config -count=2` | **ok 3.603s**（本腿加的词表/枚举三测试在内） |
| d22scan | `sh scripts/d22scan.sh` | **clean rc=0**（ban#8 internal/=912 · cmd/=217 全量 examined；"no D22 ban violations"） |
| build all | `go build ./...`（克隆树无未跟踪件） | **rc=0** |
| C21 机器 | `go test -run "TestC21\|TestStyle\|TestWensity\|TestPseudo" -v` | 全 PASS：colour rows **verified 126/132**（+9 look 行 +9 带行）、look colours **45**（5 套）、geometry **81 枚/32 行**、零消费者 **24/136 全部被 21 行豁免覆盖**、CSS 三方 **78/78** |

`-short` 口径（同 68-r2）：包内无 `testing.Short()` 分支；桌面/live 测试由 `//go:build windows && winlive` 排除，非被 `-short` 跳过。

## 4. D32 纪律核实（静态，全部本腿可证）

- `anim.go` **零 diff**（`git diff --stat` 可查）：`AnimationPolicy` 权威表一字未动；`TestAnimationPolicyZeroTimerInSleeping` 原样绿。
- 零新定时器/线程/RAF/WebGL：新绘制全是 `setAssets` 烘焙 brush 的 fill/stroke；伪谱三带在 `pushLevel` 同步派生（音频样本即时钟，`liquid_windows.go` 头注的两帧源不变）；`syncLiquidTimer` 的 busy 判据未动。
- Sleeping 静帧确定性：`TestWensityRestingFrameIsStatic`（零谱带 + 零形变 + `VisualFor` 全 20 态零 spectrum）+ `TestStyleSwitchNeverWakesSleeping`（style 切换零 policy 位移 + Sleeping 双白名单外）。
- 44px / CPU ≤0.5% / WS ≤25MB / `px_delta_ge8 ≥ 2120` / `timers=no`：**留真机**（§5）；静态推断——带与弧只增可见像素，门只升不降。

## 5. 明确未测（留给真机，本腿零 GUI 进程）

- **AC：真桌面差分整格未做**：`build/balldebug.exe -diff -diff-gate`（Sleeping 行 `px_delta_ge8≥2120`、像框家族、`timers=no`、`cpu_all≤0.5%`、私有工作集 ≤25MB）；`-style glass` 与 `-style jarvis` 各一遍对拍；owner 第二次 `-tour` 亲判「这次对了」（票 65 完成判据是人）。本腿未运行任何 exe、未建窗、未截屏、未热键注入。
- 带的**观感**（FWHM/白核/弧在 44/56/72px 档的读感）是几何推演，未经像素证实——数字全部 token 化，桌面上偏了改 token 不改代码。
- `dock_windows.go`/`hit.go` 行为零触碰（11 条 dock 确定性测试原样绿）。

## 6. 移交项（编排者）

1. **SPEC-03 [ball] 节文字**：按 §4.2 应随本次批准同步（style 键）——SPEC 属冻结面，本腿未动，移交编排者随台账登记一并办。
2. **台账 A### 追记**：D36 面 [ball] 新键 style 已落实现（PLAN.md:2741 一行 + config 字段 + 校验），引用 A1581 预授权。
3. **票面 Progress log**：本腿未编辑票面（写面授权未含），收件后由编排者补一行。
4. **jarvis 渲染腿**：`StyleJarvis` 已收值落玻璃帧（注释+测试写明）；下一腿落 `jarvis_windows.go` + JarvisTokens + 静帧确定性测试（65-style-direction §1.2/§2.1 Style B 全套数字已备）。
5. `docs/evidence/s1/c21-native-tokens.md` 的 `design/assets/tokens.css` 移交面**无变化**（wensity 全部原生自有，无 CSS 引用可查）。

## 7. 临时件（只建不删，规矩 8）

- `D:\tmp\wisp\20261106\65glass\repo`（规则 21 共享克隆，全部门禁读数出自此树；alternates 指回主树——**主树 gc 前勿 prune**）
- `D:\tmp\wisp\20261106\65glass\repo\build`、`...\third_party\sherpa-onnx`（junction 复用主树，非副本）
