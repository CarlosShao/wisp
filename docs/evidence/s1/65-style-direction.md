# 65 — 球风格方向定案（W65-design 只读设计分析腿）

**Status:** analysis-complete（实现腿 W65-impl 待 W68 收口后另派）
**Leg:** W65-design（wave37，只读分析+证据；本腿写面 = 仅本文件；未改任何生产代码）
**HEAD 锚点:** `fba64a925`（现 HEAD `9c540abf7`，为 A1575 派单提交，分析基于其父锚点内容未变）
**Owner 方向:** 台账 **A1575**（`docs/reports/pending-and-issues.md:31375-31381`，只读引用）＋ A1576 确认「就这么做」
**票面:** `.scratch/wisp/issues/65-ball-glass-quality-rework.md`
**Spec refs:** SPEC-08 §2/§2.1（动画纪律、C21 强制、各态视觉）· D32（PLAN.md §16.3）· D36（PLAN.md:2731-2757）· C21 DesignTokens
**参考图:** `.scratch/wisp/probes/65/refs/`：`wensity-siri-orb.png`(1173×680) · `youtube-audio-visualizer.png`(1072×650) · `jarvis-orb.png`(1341×679)

---

## 0. 一页结论

1. **风格可切换**：默认 **glass**（Wensity 深色玻璃球＋水平彩虹光带＋边缘折射），可选 **jarvis**（青色粒子环球，idle/thinking/success/alert 四态＋自定义配色）。两风格的活跃态（Listening/Speaking/Thinking）动效统一吃音频驱动。
2. **全部量化属性已从三张参考图逐像素采样**（§1），采样脚本在档可复跑（§5）。
3. **实现完全落在 internal/ball 的现有骨架上**：`BallStyle` 枚举＋每风格一份 token 集（走 C21 双胞胎校验机器）＋ renderer 分支；**动画计时器零新增**（`AnimationPolicy` 原封不动），音频输入保持现有标量包络缝（`SetAudioLevel`），v1 不碰 internal/audio。
4. **D32 一条不放宽**：Sleeping/静止＝单帧静态渲染、零定时器；动画只在活跃态；无常驻渲染循环、无 WebGL/RAF（全程 D2D on ui-sta）；验收沿用票 62 差分门（`balldebug -diff -diff-gate`，`px_delta_ge8 ≥ 2120` 基线、`timers=no`、CPU ≤0.5%、私有工作集 ≤25MB、44px 下限不变）。
5. **设置面有一个真契约点**：`[ball] style` 新键属 **D36 契约面**（[ball] 行是键清单＋未知键报错），需按 SPEC-12 §4.1 批准并登记——A1575 备注已预告此事；一行 diff 已备好（§4.2）。C17 白名单**不需要**新方法（`config.get/config.set` 已在册）。
6. 一个必须让机主知道的偏差：**jarvis 的 idle 在参考站里是持续旋转的粒子场，但 D32 禁止**——我们的 jarvis idle（及一切非活跃态）是静态粒子帧。这是 D32 硬约束的直接结果，不是偷工（§2.4）。

---

## 1. 参考图量化属性（逐像素采样，非目测）

**方法与误差声明**：以下全部数字来自对三张 PNG 的程序化采样（Pillow+numpy；脚本与复跑命令见 §5）。坐标为图像像素，颜色为 sRGB hex。浏览器截图含页面 chrome（标签栏/侧栏），采样均已扣除。JPEG/缩放噪声 ⇒ 颜色按 ±8/通道读；实现腿落 token 前应按 §5 命令复核一遍。

### 1.1 风格 A「glass」＝Wensity siri-orb（`wensity-siri-orb.png`）

**关键结构发现（先于数字）**：球的暗玻璃体（#05050E–#1B1B1E）与页面背景（#0A0B0C）**几乎同色**——可辨识性几乎全部来自 ①水平彩虹光带 ②顶部白色弧形高光 ③右下亮唇折射 ④底部微内辉光。这直接决定 token 设计：**体色只是一个近似融入背景的暗基底，带/高光/唇才是主角**（对 Wisp 的推论：浅色壁纸下走 LightPalette 的白玻璃变体，与 D29 现行安排一致）。

**几何**（由行/列扫描钉定）：球心 ≈ (750,305)，半径 R ≈ 125px（bbox x 623–877 / y 143–435，含辉光外扩 ~4%）。

| 属性 | 测量值 | 采样出处 |
|---|---|---|
| **光带垂直位置** | 峰值 y≈289，= 球心上方 **0.13R**；整带亮区 y 276–336 | x=750 垂直亮度剖面（lum 71@y276 → 255@y288-296 → 97@y336） |
| **光带宽度** | FWHM ≈ **49px ≈ 0.39R**（半高截 y282–330，v4 精测）；高斯型软落，无硬边 | 同上 |
| **白核** | band 内 lum>235 且 chroma<40 的像素（扣除顶部弧：y>球心−0.6R）n=2843，质心 (−0.08R, −0.06R)，横向跨度 −0.84R..+0.63R（v4 精测） | 饱和/白核聚类 |
| **光带色谱（左→右，按色相聚类质心）** | 深红 `#C2162F`@(−0.72R,+0.27R) → 品红 `#E649AF`@(−0.50R,+0.15R) → 紫 `#8C4BCB`@(+0.01R,+0.13R) → 蓝 `#246BD1`@(+0.55R,−0.24R) → 青 `#48D7E6`@(+0.65R,−0.11R) → 橙 `#D4994B`@(+0.59R,+0.05R)（青/蓝在中心线**上方**，品红/红在**下方**——带整体向右微翘） | chroma>70 的 10317px 按色相 8 bin；每 bin n=93–3020，maxchroma 162–248 |
| **顶部弧形高光** | 纯白 `#FFFFFF`，贴上缘 −0.95R 处，厚 **4–6px（0.03–0.05R）**，横向覆盖 x 660–840（≈0.9R 宽）；x=800 处已消失（弧随球面弯出采样列） | (728,185) 峰值；x700/x750 剖面 lum>90 行段 |
| **边缘折射** | 右下 45° 方向 0.97R 处**亮唇 `#F7FAFF`**（1–2px）；其余三向边缘为暗缘 `#0B0D11–#27292F`（非对称折射是参考图的特征，不是均匀描边） | 四条 45° 对角径向走样 |
| **底部内辉光** | 0.75R 下方 ≈ `#53565E`（弱、中性偏冷） | v2 采样 |
| **体色** | `#05050E`–`#1B1B1E`（对页背景 #0A0B0C 的对比 <10/通道）；带辉光渗漏区 `#0E35A8`/`#2C96FF`（深蓝渗光） | 体采样 4 点 |
| **径向亮度剖面** | 中心 250 → 0.5R≈76 → 0.92R≈21（谷）→ 0.96R≈58（唇峰）→ 域外 20；谷-峰结构即「暗体+亮缘」 | 2° 步进平均 |
| 饱和度上限 | 全图 maxchroma ≈ 248（蓝/青段），即带内颜色接近全饱和 | 色相聚类 |

**Wensity 采样备注**：绿 bin（n=93，`#E7EA90`）与右缘暗红 bin（n=283，`#8A3C15`）样本过少，判为带末端的过渡混色，不单列 token。

### 1.2 风格 B「jarvis」＝cyber1443 jarvis-ai-orb（`jarvis-orb.png`）

页面为 2×4 tile 演示矩阵；每 tile 取以球心为圆心 r≤135px 的窗口统计（亮粒子 = max(R,G,B)>60；背景 = <25）。

| 指标 | cyan/idle | aurora/**thinking** | ember/**success** | cyan/**alert** | violet(自定义) | 'listening'(自定义) |
|---|---|---|---|---|---|---|
| 粒子覆盖率（亮像素/tile面积） | **10.0%** | **12.6%** | 10.2% | 10.8% | 14.8% | 14.8% |
| 粒子平均色 | `#1B7C86` | `#298790` | `#227B84` | `#1B7A84` | `#562D98` | `#258D96` |
| 粒子 p90 亮色 | `#36E7F5` | `#5EFAFF` | `#4AE2EE` | `#35E0EF` | `#9F63FF` | `#64FFFF` |
| 核心色（最亮像素） | `#C6FFFF` | `#FFFFFF` | `#F0FFFF` | `#CFFFFF` | `#FFFFFF` | `#FFFFFF` |
| 轨道环（径向亮度峰 @R） | 0.07, 0.16, **0.50** | 0.05, 0.23, **0.67** | 0.16 | 0.07, 0.16, **0.39** | 0.16, 0.25, 0.63, **0.81** | 0.14, 0.23, 0.70 |
| tile 底色 | `#050913` | 同 | 同 | 同 | 同 | 同 |

**从表读出的设计规则**：

- **粒子密度**：亮粒子覆盖 **10–15%**，即 56px 球（r=28px）≈ **60–120 个可见粒子点**；粒径 1–2px，其中混少量白色高亮粒子（p90 色即「亮粒子」token）。
- **环几何**：**2–3 条椭圆轨道环**（径向剖面峰在 0.16R 附近一条内环、0.39–0.81R 一条主环，核心 0.05–0.07R 一条）；环线为细虚线/粒子链，主环半径**按态伸缩**（idle 0.50R → thinking 0.67R → alert 收缩 0.39R）。
- **态间差异是结构差**（环数/环位/覆盖率），**palette 间差异是色相差**（cyan 族 hue≈187°，violet 族 hue≈265°，`#9F63FF` 实测）；自定义配色 = 换一族同明暗结构的色相。
- **采样诚实备注**：标着 **ember/success** 的 tile 采样为青色（cyan%=93、warm%=0）——演示站该 tile 实际渲染的不是琥珀色系。本表如实记录；ember 命名仅作参考站的调色板名，**实现取色以 cyan/violet 两族实测为准**，success 态的暖色版按 C21 语义色（`--success`/`--tint-success-*`）派生，不从该 tile 采样。
- 状态到 Wisp 状态机的映射：idle→Sleeping/rest（**静帧**，见 §2.4）· thinking→Thinking · success→Speaking · alert→Confirming/Error · listening→Listening。

### 1.3 动效参考「#3」＝YouTube 3D audio visualizer（`youtube-audio-visualizer.png`，非第三风格）

| 属性 | 测量值 | 出处 |
|---|---|---|
| 网格团主体色 | 核心 `#A4D7E0` / 中层 `#ACDCE6` / 边缘 `#B2DEE7`（边缘反比核心略亮，半透网面质感） | 环带均值 r<40 / 40-90 / 90-130 |
| 视频背景 | `#324E53`（不是纯黑；团靠相对亮度差成像，差 ≈ +130 lum） | 中位数采样 |
| **形变振幅（三阈值分档，v4/v5 精测）** | **网格实体缘**（bg+120）：mean 125px，峰-峰 p95−p05 = **15px ≈ 0.12R**；中间层（bg+60）：11px ≈ 0.09R；**辉光外缘**（bg+14）：6px ≈ 0.04R，近圆。轮廓为**大尺度平滑起伏**（7px 平滑、prom≥1.5px 下程序计数 0 尖瓣；目测细碎凸包属更小尺度纹理） | 360° 轮廓走样 ×3 阈值；静帧读数＝该时刻下限，视频内形变随音频更大 ⇒ 实现把振幅做成 token（满电平默认 pp≈0.15R） |
| 轮廓瓣数 | 程序计数（平滑窗 7）= 1（大尺度摆动主导）；目测 12–18 个小凸包为更小尺度纹理——实现按「平滑大尺度起伏＋噪声微扰」建模，不按尖瓣正弦 | 局部极大检测 |
| 辉光衰减 | r=100 lum 218 → r=140 lum 100 → r=170 lum 65≈背景：**衰减带约 0.25R 宽**，bloom 较紧 | 径向亮度 |

⇒ 可移植到 D2D 的三个动效不变量：**①形变集中在网格实体缘（pp≈0.12R），辉光外缘近圆（pp≈0.04R）——不是整团乱涨；②起伏是大尺度平滑起伏，非尖瓣脉动；③亮核+紧辉光（衰减带 ~0.25R）**。振幅做成 token，随包络缩放。

---

## 2. 实现方案（internal/ball —— 由 W65-impl 执行；本腿零生产代码）

### 2.0 总体形状：一个枚举 + 两份 token 集 + 一个 renderer 分支

```
config [ball].style ("glass" 默认 | "jarvis")
   → resident 接线（与 size 同缝，hot-tier 热生效）
   → ball.SetStyle(name)                          （tokens.go；仿 SetLook :328 现形）
   → statevisual.go mapState(...)（Visual 结构不变，逐态映射原封不动）
   → renderer_windows.go drawFrame：switch activeStyle {
        StyleGlass  → drawGlass（现有 ticket-62/68 玻璃路径，零改动起步）
        StyleJarvis → drawJarvis（新文件 jarvis_windows.go，纯 D2D）
     }
```

要点：
- **Visual 管线不动**（`statevisual.go:27-64`）：两风格消费同一 `Visual`（SizePx/Opacity/CoreColor/GlowColor/Ring*/Icon/LiquidLevel/LiquidAngle/SummonFlow/Dock*），renderer 仍是 Visual 的纯函数。新增字段只有 §3 的伪谱带（加在 liquid 侧，见 §3.2）。
- **W68 并发避让**：实现腿只在**新文件**（`jarvis_windows.go`、`style.go`）落新代码，`tokens.go`/`statevisual.go` 的增补点在 W68 收口后由编排者统一排；`liquid.go/anim.go/dock*.go` 的**行为路径零改动**（票 65 硬约束：11 条 dock 确定性测试与全部球测试的行为断言不得动）。
- **现有 4 套 -look（aurora/glacier/nebula/solar，`tokens.go:253-278`）保留**为 glass 风格的子 treatment（balldebug `-look` 继续可用）；Wensity 带状款作为新增第 5 套 look **`wensity`** 并设为默认。

### 2.1 C21 DesignTokens token 集（含双胞胎机器的登记路径）

**Style A — glass（Wensity 款，新增 `wensity` look）**：在 `LiquidLook`（`tokens.go:239-250`）旁新增平行结构（不动现有 9 字段 × 4 look 的表契约，避免把每 look 行数翻倍）：

```
WensityBand {                                    // 全部实测值来自 §1.1
  BodyBase        Color  // #0A0C12 @ a≈0.90  暗玻璃基底
  BandStops       [6]Color // #C2162F → #E649AF → #FFFFFF → #246BD1 → #48D7E6 → #D4994B
  BandCenterFrac  = -0.13 // 带心在球心上方 0.13R
  BandWidthFrac   = 0.39  // FWHM 0.39R，高斯软落
  BandTiltUpRight = 0.33  // 色谱两端质心高度差：暖端(品红/红)(-0.56R,+0.18R)、冷端(青/蓝)(+0.51R,-0.15R)，冷端上翘（v5 精测）
  SpecArcColor    Color  // #FFFFFF；厚 0.04R，贴 -0.95R 上缘，横向 ~0.9R
  RimLipColor     Color  // #F7FAFF @ 右下 0.97R（非对称亮唇，1-2px）
  RimDarkColor    Color  // #1B1D24（其余三向暗缘）
  BottomGlow      Color  // #53565E @ 0.75R 下方
}
```

**Style B — jarvis（新 `JarvisTokens`）**（实测值来自 §1.2）：

```
JarvisTokens {
  BgGlow          Color  // #050913 基底辉光
  Core            Color  // #C6FFFF→#FFFFFF（核心亮点，直径 ~0.10R）
  ParticleDim     Color  // #1B7C86（暗粒子体色）
  ParticleBright  Color  // #36E7F5（p90 亮粒子）
  RingLine        Color  // 环线 = ParticleDim 的半透形式（虚线弧）
  // 态覆盖（只换色相/亮度，结构不变——参考站实测规则）：
  Thinking:  Ring 主环 0.50R→0.67R，ParticleBright→#5EFAFF，覆盖率 +25%
  Speaking(Success): 核心呼吸用 C21 --tint-success-hi/lo 派生（ember tile 采样不可信，§1.2 备注）
  Alert:     主环收缩 0.39R，色相偏 --danger（Confirming 语义）
  Listening: 同 idle 结构 + §3 的音频形变
  // 自定义配色：config 只存一个基色相 hue（0-360）；
  //   Dim/Bright/Core 由固定 HSL 变换派生（violet 族 = hue 265° 实测校准），
  //   代码里零 hex 字面量（用户值在 config，代码只有变换式）——过 TestNoHardcodedColorsInBallPackage。
}
```

**几何/运动常量**（全部进 tokens，渲染零字面量）：环数 3（内 0.16R/主 0.50R/核 0.07R）、粒子数 N=96（56px 档；44px 档 72，72px 档 128，按尺寸分档）、粒子固定**种子化 PRNG** 点集（每 (size,DPI,palette) 预生成一次；静帧与动画帧共用同一几何，动画只推进相位/半径/亮度，不重掷随机数——保证静帧确定性可测）。

**双胞胎校验登记**（实现腿必做，否则红）：
1. `docs/evidence/s1/c21-native-tokens.md` 增加 `wensity` look 行（9 字段×新字段另列）与 jarvis 段——`TestC21TableColourRowsMatchCode`（`tokens_table_test.go:668`）按表逐行核对代码值，行数下限检查是 `>=`，**加行不需要改测试**，但**代码加字段不加行必红**。
2. `internal/ball/testdata/tokens_c21.css` fixture（`tokens_table_test.go:1507` 的 `TestC21TableColourRowsMatchTokensCSS` 读它）：fixture 是 `design/old/assets/tokens.css` 的逐字节拷贝（头部已写明更新规程）——**design/** 归机主外部 agent**，本仓不写 design/**；流程 = 先在 c21-native-tokens.md 登记 native 侧行（机器校验以此表为准），fixture 与 CSS 真相源的对齐登记为编排者→外部 agent 的移交项（残留决策 ⑤）。
3. `TestNoHardcodedColorsInBallPackage`（`tokens_test.go:88`）：新文件 jarvis_windows.go/style.go 从第一天起只引用 token 值。

### 2.2 静帧与渲染预算（D32 硬约束逐条）

| D32/SPEC-08 §2 约束 | 两风格的满足方式 | 验收 |
|---|---|---|
| Sleeping/静止 = 静态帧、**零定时器** | `AnimationPolicy`（`anim.go:40-69`）**一字不改**：非活跃态 `AnimNone`。风格切换/配置热加载只触发**一次** repaint（重画一帧静帧→ULW 上屏→停） | `balldebug -diff -diff-gate`：`timers=no` 必须保持 |
| 动画只在活跃态、复用现有计时器 | Listening/Thinking/Confirming/Speaking 的 `AnimFrame`（30fps 上限、framePeriod 分步）＋ Warm `AnimAlpha`（10fps 恒 alpha 通道）＋ Settling 一次性 fade——五个授权态之外不出定时器，`transitionDriven`/`liquidDriven` 白名单原样生效 | 现有 `liquid_test.go`/`anim` 相关测试全绿（不改断言） |
| **无常驻 WebGL/RAF** | 无 WebView、无 GPU 渲染循环：jarvis 全部是 D2D 原语（3 条虚线椭圆弧 + ≤128 个小椭圆点 + 1 个核心光斑），按帧画进现有 ULW 位图；几何按 (size,DPI,palette) 预构建，换档才重建 | 代码评审：包内无新线程、无新定时器、无 COM 新对象常驻 |
| CPU ≤0.5%（Sleeping） | Sleeping 不画第二帧（与票 62 现状同构：`cpu=0.000` 实测基线） | 差分表 CPU 列 |
| 私有工作集 ≤25MB | 增量 = 窗口位图（现状已有）+ 预构建几何（KB 量级）+ N 点坐标数组（128×16B≈2KB）——无新进程、无新 DLL、无大纹理 | 差分表 WS 列（基线 11.76MB） |
| 44px | `BallSizeMinPx=44`（`tokens.go:401`）不动；Sleeping 静止体径仍走 `SleepRestRatio 0.62`/`SleepingRestMinPx 30`（SPEC-08 §2 INTERIM 现文）；A1575 的「44px」读作配置下限，不产生新契约 | `TestSleepingSizeTruthTable` 不动仍绿 |
| 句柄 <300 · GDI <200 · goroutine ≤6 | 零新增常驻对象；D2D 几何在设备重建时随旧资源一起释放（现 discard 路径） | D32 附带指标采样（`slo` 侧已挂） |
| **px_delta_ge8 ≥ 2120** | jarvis 静帧含粒子环，可见像素**只会多不会少**；门只许升不许降（`-diff-gate-min-px8` 语义照旧） | Sleeping 行差分重跑（票 62 基线 `docs/evidence/s1/62-diff-glass/diff-table.txt`） |

活跃态帧成本：jarvis 一帧 = 玻璃款同量级（玻璃款每帧画多层径向渐变+DirectWrite 文本；jarvis 每帧 ~130 个小原语，D2D 填充率远低于前者），活跃态 CPU 本就计入工作态预算（D32 表 Conversation ≤5% 等），以差分/采样实测为准，不做估算替代实测。

### 2.3 行为不动承诺（票 65 硬约束回写）

- 靠边吸附/弹回、音频驱动液面（glass 款）、边框滑入：行为与 11 条 dock 确定性测试**零改动**；jarvis 款下 dock 仍是窗口层行为（`dock_windows.go` 不进 diff）。
- 票 65「只改材质不改行为」在 jarvis 款的读法：**交互行为层（点击/拖拽/吸附/热键）两风格完全一致**；差异仅在像素材质。

### 2.4 与参考站的显式偏差（D32 的代价，报机主）

| 参考站行为 | Wisp 实现 | 原因 |
|---|---|---|
| jarvis idle 粒子场持续旋转 | **静帧**（一次渲染后零活动） | D32：非活跃态零定时器零动画（残留决策 ②，若机主要 idle 微动需新豁免） |
| YouTube 团为实时 WebGL 网格 | D2D 多瓣轮廓近似（预生成轮廓多边形 + 半径按包络缩放） | 无常驻 GPU 管线纪律（SPEC-08 §2 全程 D2D） |
| glass 款带内彩虹可能缓慢流动（参考站为 CSS 动画） | 仅活跃态流动；静止 = 定格单帧 | 同 D32 |

### 2.5 测试面（W65-impl 的 AC 草案）

1. token 双胞胎：新 look/新段行登记后 `TestC21TableColourRowsMatchCode` + `TestC21TableColourRowsMatchTokensCSS` 绿；`TestNoHardcodedColorsInBallPackage` 绿。
2. 风格选择器：`SetStyle` 枚举路由测试（未知名拒绝、默认 glass、切换后 repaint 标志置位）。
3. jarvis 静帧确定性：同 (state,size,DPI,palette,seed) 两次渲染位图哈希一致（粒子种子化）。
4. `AnimationPolicy` 相关现有测试不改一字仍绿（计时器零新增的白盒证明）。
5. Sleeping 差分门重跑（两风格各一次）：`timers=no`、`px_delta_ge8 ≥ 2120`、CPU ≤0.5%、WS ≤25MB、px 可见性不回退。
6. 票 65 既有 AC 全部适用（并排差分产物、owner 第二次 `-tour` 亲判「这次对了」——完成判据是人，不是测试绿）。

---

## 3. 音频频谱映射（两风格共用一个驱动）

### 3.1 现状与 v1 决策

现状：采集侧把 **512 采样 RMS 标量包络**经 `Ball.SetAudioLevel`（`liquid_windows.go:42`）喂给 `liquidMotion.pushLevel`（`liquid.go:190-245`：attack 60ms / release 420ms / 静音门 0.06 / 保持 150ms），`liquidDriven` 六态白名单放行。

**v1 决策：不碰 internal/audio**（W68 在邻区、C8 是注入缝契约）——在 `liquidMotion` 内由标量包络派生**伪谱 3 带**（纯确定性函数，无新时钟）：

```
band-low   = 慢积分器（τ≈400ms）——音量「体感」
band-mid   = 现有 level 本身（τ 60/420ms）
band-high  = 包络正变化率的快低通（τ≈40ms）——「嘶/齿音」近似
```

三带只作为 `Visual` 新字段（`SpectrumLow/Mid/High`，0..1）进 renderer，计算仍在 `applyTo`（`liquid.go:312-323`）同一点、同一放行白名单内。**真 5-band FFT 需要动采集侧（C8 邻接）→ 登记 v2/推迟项**（残留决策 ③）。

### 3.2 映射表（全部只消费 Visual 字段；帧由现有 AnimFrame 计时器供给，无新时钟）

| 态 | glass（Wensity 带） | jarvis（粒子环） | 驱动量 |
|---|---|---|---|
| **Listening** | 带高度/亮度随 level 增缩（幅 ±0.06R，借 `LiquidGatherPerLevel` 同型增益）；带内色谱沿 `LiquidAngle` 横向漂移（转速=现有 ω 公式，`liquid.go:224-226`） | 主环半径脉动 0.50R±0.08R×level；粒子散布 ±(0.05+0.15×level)R；轨道角速=现有 ω | level / 三带 |
| **Thinking** | 带心位置以 `ThinkingSweepMs` 相位上下扫（±0.08R）；亮度随 band-high 闪 | 主环 0.50→0.67R 伸展开（参考站 thinking 实测位）；虚线环相位随帧相位推进 | 帧相位 + band-high |
| **Speaking** | 白核亮度按 `SpeakingBreathMs` 呼吸 ×(0.6+0.4·level)；带随 band-mid 起伏（YouTube 形变不变量①②：网格缘 pp≈0.12R、辉光近圆、大尺度平滑起伏；满电平默认 pp≈0.15R，token 可调） | 核心亮度呼吸；粒子抖动幅 ∝ level；环微扩 0.05R×level | band-mid/level |
| **Confirming**（alert） | 现行 danger 脉冲环叠加在带上（环参数照旧） | palette 换 danger 派生 + 主环收缩 0.39R + 单次 `RingPulse` | 无音频（策略态） |
| **Warm** | 恒 alpha 呼吸照旧（`AnimAlpha`，零重绘） | 同左（粒子帧不变，只变窗口常量 alpha） | 无 |
| 静止/其余态 | 静帧 | 静帧 | 无 |

反归纳护栏（写进测试）：非活跃态 `applyTo` 后三带必须为 0（`liquidDriven` 白名单外的态音频**到不了** renderer——`liquid_windows.go:56` 现行丢弃语义不变）。

---

## 4. 设置面（config 字段 + panel 字段 + D36 契约判定）

### 4.1 config.toml 字段

```toml
[ball]
style = "glass"        # "glass" | "jarvis"（hot-tier，与 size 同节）
# jarvis_palette_hue = 187   # v2 预留：自定义基色相 0-360（v1 不落，见残留决策 ⑥）
```

接线：`internal/config/schema.go` `BallSection`（:167-179，全节 hot-tier）加 `Style string \`toml:"style" default:"glass"\``；`validate()` 枚举两值；设置写入走 `writeOneKey` 单键 setter 模式（`settings_320.go` 同形）；resident 读 `[ball]` 处与 `size` 同缝调 `ball.SetStyle` + repaint。`[ball]` 本就是 hot-tier ⇒ 改 `style` 热生效（一帧重绘，无常驻成本）。

### 4.2 D36 契约判定：**是契约面，需登记**（残留决策 ①）

判定依据：D36 的 `[ball]` 行（`docs/PLAN.md:2741`）是**键清单**（`size position{x,y,monitor} opacity_idle click_through hide_on_fullscreen`），且六规则之 2「未识别的键 → 报错」——新增 `style` 键同时改了「GUI 与 TOML 对得上的面」与 D36 表文字。按 SPEC-12 §4.1，**改 D1–D47 文字须机主批准**。A1575 备注已预告：「加『球外观』设置项触及 config/panel 字段，属机主明确要求的配置新增，按批准办（若确属 D36 契约面则登记）」。

⇒ **本腿判定：属 D36 契约面。** 需要的登记动作（编排者）：
1. 台账 A### 追记一行（D36 面：[ball] 新键 style，A1575 预授权引用）；
2. `docs/PLAN.md:2741` 一行 diff 待批：`| [ball] | size **style** position{x,y,monitor} opacity_idle click_through hide_on_fullscreen | hot |`（style 标 hot）；
3. SPEC-03 的 [ball] 节文字随实现腿同步（同一次批准覆盖）。

**不批的代价**：style 无法落 config，只能 `balldebug -style` 调试切换——机主要「可换风格」落空。

### 4.3 panel 字段

- 设置页「球外观」下拉（简洁玻璃/Jarvis 粒子）＋ Jarvis 选中时的配色子项——属**面板设置 UI**，票 39/40 与 S5 范围（票 65 Out of scope 明说不做配置界面）；本票只保证 config 键与 ball 消费端就绪。
- **C17 白名单无需改动**：面板读写走在册的 `config.get`/`config.set`（`internal/panel/bridge.go:66-67`；SPEC-08 §5.2 表），不新增方法 ⇒ 不触 C17 契约面。
- `balldebug` 加 `-style glass|jarvis`（现 `-look` 保留为 glass 子 treatment；`-frozen` 语义不变）——owner `-tour` 验收时一键对拍两风格。

---

## 5. 复核命令与证据基线

```bash
# 参考图逐像素采样（本腿实际执行的四把尺，可复跑；v4/v5 为修正轮）：
"C:\Users\swq\.dsh\dsh-runtimes\dsh-primary-runtime\dependencies\python\python.exe" \
  .scratch/wisp/probes/65/sample_refs.py        # v2：jarvis 六 tile 统计（§1.2 表即其输出）
"C:\Users\swq\.dsh\dsh-runtimes\dsh-primary-runtime\dependencies\python\python.exe" \
  .scratch/wisp/probes/65/sample_refs_v3.py     # v3：wensity 定心采样 + youtube 背景相对阈值
"C:\Users\swq\.dsh\dsh-runtimes\dsh-primary-runtime\dependencies\python\python.exe" \
  .scratch/wisp/probes/65/sample_refs_v4.py     # v4：带 FWHM 精测 + 白核去顶部弧污染 + yt 辉光缘
"C:\Users\swq\.dsh\dsh-runtimes\dsh-primary-runtime\dependencies\python\python.exe" \
  .scratch/wisp/probes/65/sample_refs_v5.py     # v5：色相质心两端高度差 + yt 三阈值形变分档
```

**修正史（防后人拿初版数字）**：v2/v3 初测的三处被 v4/v5 否证——光带 FWHM 0.30R→**0.39R**；白核统计扣除顶部弧污染（n 3063→2843）；YouTube 形变振幅初测 0.30–0.37R 系扫描上限污染，三阈值精测后**网格实体缘 pp≈0.12R、辉光缘近圆**。本文正文已全部采用修正值。

- 采样脚本为临时件，落 `.scratch/wisp/probes/65/`（只建不删，不进 git pathspec）。
- D32 差分基线（不重测不改）：`px_delta_ge8=2120` · `timers=no` · `cpu=0.000` · 私有工作集 11.76MB（票 62，`docs/evidence/s1/62-diff-glass/diff-table.txt`；SPEC-08 §2 INTERIM 实测）。
- 代码事实复核：`SetAudioLevel` 标量缝 `internal/ball/liquid_windows.go:42`；计时器权威 `internal/ball/anim.go:40`；色值禁令 `internal/ball/tokens_test.go:88`；D32 Sleeping 行 `docs/PLAN.md:2261`；D36 [ball] 行 `docs/PLAN.md:2741`。

## 6. 本腿做了什么 / 没做什么（写面核实）

- **写面**：仅本文件（git pathspec 单文件）；临时采样脚本 2 枚落 `.scratch/wisp/probes/65/`（不提交）。
- **没做**：零生产代码；未碰 `internal/ball`（W68 在飞）、`internal/audio`、`internal/config`、`design/**`、`frontend/**`、`docs/reports/pending-and-issues.md`（台账只读引用 A1575/A1576）；无 push。
- 可核命令：`git show --stat HEAD`（本腿 commit 只含本文件）；`git status --porcelain -- internal/ design/ frontend/`（本腿零输出）。

## 7. 残留机主/编排者决策（编号清单；「不做」列必填）

| # | 决策 | 推荐 | 不做的代价 |
|---|---|---|---|
| ① | **D36 契约文字**：`[ball]` 加 `style` 键（§4.2 的一行 diff） | **批**（A1575 已预授权方向，差的是表文字） | style 进不了 config，只能调试旗标切换，「可换风格」对机主不可达 |
| ② | jarvis **idle 持续旋转**是否豁免 D32 | **不做**——守 D32 静帧（Siri/ChatGPT 静止时也不动） | 若豁免：Sleeping 出现常驻定时器，D32 差分门 `timers=no` 必红，等于改契约+改门禁，代价远超收益 |
| ③ | 真 FFT 5-band（动采集侧/C8 邻接）v1 还是 v2 | **v1 用伪谱 3 带**（零 audio 侧改动）；真频谱登记 v2 推迟项 | 动效吃「包络+派生带」而非真频谱——视觉差异小（3 带已覆盖低/中/高频感），但极端音色的频谱细节不动 |
| ④ | 产物名：本文件 vs 票面 AC 写的 `65-glass-reference-analysis.md` | **本文件为上位替代**（含参考分析+实现方案+映射+设置面四合一），票面 AC 文案由编排者对齐 | 收件核票时 AC 指向的文件名对不上，判红一次 |
| ⑤ | C21 fixture 对齐：native 侧行进 c21-native-tokens.md 后，`design/**`（外部 agent 域）的 tokens.css 与 fixture 的同步 | native 表为机器真相先行；CSS 同步开移交项给外部 agent | fixture 暂陈旧 → `TestC21TableColourRowsMatchTokensCSS` 红一次提醒同步（这正是该测试的设计行为），需编排者接受一次「已知红」窗口 |
| ⑥ | jarvis 自定义配色（config 存 hue）v1 就落还是 v2 | **v2**（v1 内置 cyan/violet 两族实测款） | 晚一步拿到自定义色；不影响风格切换主功能 |
| ⑦ | 旧 4 套 -look（aurora/glacier/nebula/solar）留或砍 | **留**为 glass 子 treatment（零运行成本，balldebug 可对拍） | 砍 = 更干净但删表行删代码；留 = 主树多 4 套无人选的配色 |
