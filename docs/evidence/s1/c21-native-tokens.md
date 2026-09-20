# C21 原生侧 token 对照表（ticket 07）

`s1/` 证据 · 生成：T07-impl · 2026-09-19 · 全表反向对账补录：票 12 AC#7 · 2026-09-20

权威真相源：`design/assets/tokens.css`（C21 契约参考实现）。原生侧（`internal/ball/tokens.go`）
逐值复制下表；`docs/evidence` 本表供人工核对。规则：球体绘制代码出现任何硬编码色值即违规
（SPEC-08 §2）—— 该约束由 `TestNoHardcodedColorsInBallPackage` 机器执行（仅 `tokens.go` 允许
出现 `rgba(` / 6 位 hex 字面量）。

值格式：`rgba(r,g,b,a)` 为 CSS 通道；Go 列为 `tokens.go` 中的等价构造（`rgba()` 0..255 通道 +
alpha，`hex()` 0xRRGGBB + alpha）。D2D 使用直通 alpha 的 `D2D1_COLOR_F`，premultiplied 仅在
`UpdateLayeredWindow` 位图写出处派生（`Color.Premultiplied()`）。

范围（ticket 12 AC#7 对账时明确，2026-09-20）：本表只覆盖**原生侧已复制的切片** —— 即
`internal/ball/tokens.go` 的 40 个 `Palette` 字段 + 全部导出的几何/动效常量。`tokens.css`
共 129 条声明，其中 80 条是面板 / screens 专用、原生侧无对应字段的 token（`--bg-base`、
`--ambient-a..d`、`--grain`、`--sheen`、`--bg-raised`、`--bg-subtle`、`--glass-blur*`、
`--border-hair`、`--border-strong`、`--fg-disabled`、`-*-hover`/`-*-pressed`/`-*-soft`、
`--accent-fg`、`--ball-core`/`--ball-mid`/`--ball-edge`、`--scrim`、`--shadow-*`、
`--mask-solid`、`--font-mono`、`--t-*`/`--lh-*`/`--fw-*` 排版阶梯、`--ls-micro`、`--r-*`、
`--s-*`、`--ease-*`、`--panel-w`/`--palette-w`/`--chip-h`/`--row-h`/`--input-h`、`--icon-*`、
`--ring`），**按契约不进本表**（面板的原生切片随票 33+ 扩表，不是遗漏）。反查命令见「审计」。

## 明色系（:root，默认）

| CSS 变量 | tokens.css 值 | Go 字段 | Go 值 |
|---|---|---|---|
| `--bg-inset` | `rgba(6, 8, 11, 0.38)` | `Palette.BGInset` | `rgba(6,8,11,0.38)` |
| `--bg-overlay` | `rgba(30, 36, 44, 0.50)` | `Palette.BGOverlay` | `rgba(30,36,44,0.50)` |
| `--glass-ring` | `rgba(255, 255, 255, 0.14)` | `Palette.GlassRing` | `rgba(255,255,255,0.14)` |
| `--glass-hi` | `rgba(255, 255, 255, 0.18)` | `Palette.GlassHi` | `rgba(255,255,255,0.18)` |
| `--glass-hi-soft` | `rgba(255, 255, 255, 0.08)` | `Palette.GlassHiSoft` | `rgba(255,255,255,0.08)` |
| `--border-soft` | `rgba(255, 255, 255, 0.10)` | `Palette.BorderSoft` | `rgba(255,255,255,0.10)` |
| `--fg-primary` | `#F2F3F5` | `Palette.FgPrimary` | `hex(0xF2F3F5,1)` |
| `--fg-secondary` | `#A6A9B0` | `Palette.FgSecondary` | `hex(0xA6A9B0,1)` |
| `--fg-tertiary` | `#71747C` | `Palette.FgTertiary` | `hex(0x71747C,1)` |
| `--accent` | `#86C2B9` | `Palette.Accent` | `hex(0x86C2B9,1)` |
| `--accent-line` | `rgba(134, 194, 185, 0.36)` | `Palette.AccentLine` | `rgba(134,194,185,0.36)` |
| `--danger` | `#E07A70` | `Palette.Danger` | `hex(0xE07A70,1)` |
| `--danger-line` | `rgba(224, 122, 112, 0.38)` | `Palette.DangerLine` | `rgba(224,122,112,0.38)` |
| `--warn` | `#D9B26A` | `Palette.Warn` | `hex(0xD9B26A,1)` |
| `--warn-line` | `rgba(217, 178, 106, 0.38)` | `Palette.WarnLine` | `rgba(217,178,106,0.38)` |
| `--success` | `#7DB896` | `Palette.Success` | `hex(0x7DB896,1)` |
| `--success-line` | `rgba(125, 184, 150, 0.38)` | `Palette.SuccessLine` | `rgba(125,184,150,0.38)` |
| `--info` | `#8AAAD2` | `Palette.Info` | `hex(0x8AAAD2,1)` |
| `--info-line` | `rgba(138, 170, 210, 0.38)` | `Palette.InfoLine` | `rgba(138,170,210,0.38)` |
| `--warm` | `#D6B184` | `Palette.Warm` | `hex(0xD6B184,1)` |
| `--warm-line` | `rgba(214, 177, 132, 0.38)` | `Palette.WarmLine` | `rgba(214,177,132,0.38)` |
| `--orb-hi` | `rgba(255, 255, 255, 0.90)` | `Palette.OrbHi` | `rgba(255,255,255,0.90)` |
| `--orb-body` | `rgba(255, 255, 255, 0.30)` | `Palette.OrbBody` | `rgba(255,255,255,0.30)` |
| `--orb-body-2` | `rgba(255, 255, 255, 0.07)` | `Palette.OrbBody2` | `rgba(255,255,255,0.07)` |
| `--orb-rim` | `rgba(255, 255, 255, 0.30)` | `Palette.OrbRim` | `rgba(255,255,255,0.30)` |
| `--orb-shadow` | `rgba(0, 0, 0, 0.45)` | `Palette.OrbShadow` | `rgba(0,0,0,0.45)` |
| `--ball-halo` | `rgba(255, 255, 255, 0.10)` | `Palette.BallHalo` | `rgba(255,255,255,0.10)` |
| `--tint-accent-hi` | `#C4E8E2` | `Palette.TintAccentHi` | `hex(0xC4E8E2,1)` |
| `--tint-success-hi` | `#BFE3CF` | `Palette.TintSuccessHi` | `hex(0xBFE3CF,1)` |
| `--tint-success-lo` | `#4E8A69` | `Palette.TintSuccessLo` | `hex(0x4E8A69,1)` |
| `--tint-warm-hi` | `#F0DCB8` | `Palette.TintWarmHi` | `hex(0xF0DCB8,1)` |
| `--tint-warm-lo` | `#96794C` | `Palette.TintWarmLo` | `hex(0x96794C,1)` |
| `--tint-warn-hi` | `#F0DCB8` | `Palette.TintWarnHi` | `hex(0xF0DCB8,1)` |
| `--tint-warn-lo` | `#8F7434` | `Palette.TintWarnLo` | `hex(0x8F7434,1)` |
| `--tint-danger-hi` | `#F3C0BA` | `Palette.TintDangerHi` | `hex(0xF3C0BA,1)` |
| `--tint-danger-lo` | `#9C4A43` | `Palette.TintDangerLo` | `hex(0x9C4A43,1)` |
| `--tint-neutral-hi` | `#B9BCC4` | `Palette.TintNeutralHi` | `hex(0xB9BCC4,1)` |
| `--tint-neutral-mid` | `#6A6D75` | `Palette.TintNeutralMid` | `hex(0x6A6D75,1)` |
| `--tint-neutral-lo` | `#3E4148` | `Palette.TintNeutralLo` | `hex(0x3E4148,1)` |
| `--on-solid` | `#FFFFFF` | `Palette.OnSolid` | `hex(0xFFFFFF,1)` |

## 亮色系（`[data-theme="light"]`，D29 双主题）

| CSS 变量 | tokens.css 值 | Go 字段 | Go 值 |
|---|---|---|---|
| `--bg-inset` | `rgba(20, 24, 28, 0.06)` | `Palette.BGInset` | `rgba(20,24,28,0.06)` |
| `--bg-overlay` | `rgba(255, 255, 255, 0.72)` | `Palette.BGOverlay` | `rgba(255,255,255,0.72)` |
| `--glass-ring` | `rgba(255, 255, 255, 0.66)` | `Palette.GlassRing` | `rgba(255,255,255,0.66)` |
| `--glass-hi` | `rgba(255, 255, 255, 0.90)` | `Palette.GlassHi` | `rgba(255,255,255,0.90)` |
| `--glass-hi-soft` | `rgba(255, 255, 255, 0.60)` | `Palette.GlassHiSoft` | `rgba(255,255,255,0.60)` |
| `--border-soft` | `rgba(16, 20, 24, 0.12)` | `Palette.BorderSoft` | `rgba(16,20,24,0.12)` |
| `--fg-primary` | `#17191C` | `Palette.FgPrimary` | `hex(0x17191C,1)` |
| `--fg-secondary` | `#55585F` | `Palette.FgSecondary` | `hex(0x55585F,1)` |
| `--fg-tertiary` | `#82858C` | `Palette.FgTertiary` | `hex(0x82858C,1)` |
| `--accent` | `#3E837A` | `Palette.Accent` | `hex(0x3E837A,1)` |
| `--accent-line` | `rgba(62, 131, 122, 0.38)` | `Palette.AccentLine` | `rgba(62,131,122,0.38)` |
| `--danger` | `#BC544C` | `Palette.Danger` | `hex(0xBC544C,1)` |
| `--danger-line` | `rgba(188, 84, 76, 0.36)` | `Palette.DangerLine` | `rgba(188,84,76,0.36)` |
| `--warn` | `#96762C` | `Palette.Warn` | `hex(0x96762C,1)` |
| `--warn-line` | `rgba(150, 118, 44, 0.36)` | `Palette.WarnLine` | `rgba(150,118,44,0.36)` |
| `--success` | `#47805F` | `Palette.Success` | `hex(0x47805F,1)` |
| `--success-line` | `rgba(71, 128, 95, 0.36)` | `Palette.SuccessLine` | `rgba(71,128,95,0.36)` |
| `--info` | `#4A6C96` | `Palette.Info` | `hex(0x4A6C96,1)` |
| `--info-line` | `rgba(74, 108, 150, 0.36)` | `Palette.InfoLine` | `rgba(74,108,150,0.36)` |
| `--warm` | `#8F6E3C` | `Palette.Warm` | `hex(0x8F6E3C,1)` |
| `--warm-line` | `rgba(143, 110, 60, 0.36)` | `Palette.WarmLine` | `rgba(143,110,60,0.36)` |
| `--orb-hi` | `rgba(255, 255, 255, 0.96)` | `Palette.OrbHi` | `rgba(255,255,255,0.96)` |
| `--orb-body` | `rgba(255, 255, 255, 0.62)` | `Palette.OrbBody` | `rgba(255,255,255,0.62)` |
| `--orb-body-2` | `rgba(255, 255, 255, 0.24)` | `Palette.OrbBody2` | `rgba(255,255,255,0.24)` |
| `--orb-rim` | `rgba(20, 24, 28, 0.16)` | `Palette.OrbRim` | `rgba(20,24,28,0.16)` |
| `--orb-shadow` | `rgba(20, 24, 28, 0.20)` | `Palette.OrbShadow` | `rgba(20,24,28,0.20)` |
| `--tint-accent-hi` | `#7FB5AC` | `Palette.TintAccentHi` | `hex(0x7FB5AC,1)` |
| `--tint-success-hi` | `#7FAE93` | `Palette.TintSuccessHi` | `hex(0x7FAE93,1)` |
| `--tint-success-lo` | `#35624A` | `Palette.TintSuccessLo` | `hex(0x35624A,1)` |
| `--tint-warm-hi` | `#C2A071` | `Palette.TintWarmHi` | `hex(0xC2A071,1)` |
| `--tint-warm-lo` | `#6E5732` | `Palette.TintWarmLo` | `hex(0x6E5732,1)` |
| `--tint-warn-hi` | `#B99A4E` | `Palette.TintWarnHi` | `hex(0xB99A4E,1)` |
| `--tint-warn-lo` | `#6B5522` | `Palette.TintWarnLo` | `hex(0x6B5522,1)` |
| `--tint-danger-hi` | `#C57E76` | `Palette.TintDangerHi` | `hex(0xC57E76,1)` |
| `--tint-danger-lo` | `#83372F` | `Palette.TintDangerLo` | `hex(0x83372F,1)` |
| `--tint-neutral-hi` | `#9C9FA6` | `Palette.TintNeutralHi` | `hex(0x9C9FA6,1)` |
| `--tint-neutral-mid` | `#71747B` | `Palette.TintNeutralMid` | `hex(0x71747B,1)` |
| `--tint-neutral-lo` | `#4A4D54` | `Palette.TintNeutralLo` | `hex(0x4A4D54,1)` |

亮色块未重定义的变量（`--ball-halo`、`--on-solid`）沿用 `:root` 值 —— `LightPalette()` 以
`DarkPalette()` 为基底覆写，与 CSS 级联语义一致。

## 几何 / 动效 token（CSS 长度与时长，96 DPI 基准，按显示器 DPI 缩放）

| CSS 变量 / 规则 | 值 | Go 常量 | 用途 |
|---|---|---|---|
| `--ball-size` | 56px（可配 44–72） | `BallSizeDefaultPx` / `BallSizeMinPx` / `BallSizeMaxPx` | 球体直径 |
| `--ball-sm` | 44px | `BallSizeSmallPx` | Armed 档参照 |
| `--ball-lg` | 72px | `BallSizeLargePx` | 大尺寸档参照 |
| SPEC-08 §2（Sleeping 微点，**原值，保留不删**） | 12px / opacity 0.35 | `SleepingDotPx` / `SleepOpacity` | 不随用户尺寸放大 —— **生产默认路径仍是此值**：`statevisual.go::stateSize` 在 `prototypeVisuals=false` 时返回 `SleepingDotPx`，`SleepOpacity` 另作 Settling 渐隐终点 |
| SPEC-08 §2 INTERIM（2026-09-20，静态玻璃体） | 休眠体 = 用户尺寸 × 0.62，下限 30px（默认档 56 → 34.72px）；渐隐落点 opacity 0.95 | `SleepRestRatio` / `SleepingRestMinPx` / `RestSettledOpacity` | 仅 `prototypeVisuals=true` 时生效；由 `fd8f838`（票 62）引入 —— **[待裁定] 见本节末「值漂移」** |
| `--dur-fast` | 120ms | `DurFastMs` | 动效时长档 |
| `--dur-base` | 180ms | `DurBaseMs` | 动效时长档 |
| `--dur-slow` | 260ms | `DurSlowMs` | 动效时长档（Settling 渐隐周期 = 260ms） |
| Warm 呼吸周期 | 2.4s（opacity 0.55↔0.7） | `WarmBreathPeriodMs` / `WarmOpacityLow` / `WarmOpacityHigh` | 合成器友好：仅改整窗常量 alpha，10fps |
| Speaking 呼吸 | 1.6s | `SpeakingBreathMs` | success 光晕呼吸 |
| Listening 波环 | 2.4s（3 层 0.8s 错峰） | `ListeningWaveMs` | accent 外环 |
| Thinking 光带 | 1.2s | `ThinkingSweepMs` | 底部 2px accent 光带 |
| Confirming 脉冲 | 2s | `ConfirmingPulseMs` | danger 脉冲环 |
| FirstRun 引导脉冲 | 2.6s | `FirstRunGuidePulseMs` | 引导脉冲周期（S1 静态，接线票启用） |
| Settling 渐隐 | 260ms，1→0.35 | `SettlingFadeMs` | 一次性有界动画（非循环） |
| 帧率上限 | ≤30fps | `MaxAnimFPS` / `MinFrameMs` / `WarmBreathFPS` | 所有动画定时器周期 ≥33ms；Warm 走 10fps 慢拍（`WarmBreathFPS`，此前只出现在用途栏未入 Go 常量列） |
| 描边宽 | 1.5px 斜杠 / 图标；Conversation 环 2px；Thinking 底部光带 2px | `SlashStrokePx` / `IconStrokePx` / `RingStrokePx` / `ConvRingStrokePx` / `ThinkingBandPx` | SPEC-08 §2.1（`ThinkingBandPx` 为本票补录） |
| 角标 / 队列点 | 17px 角标、10px 数字、2px 角标描边、9px info 点 | `BadgeDiameterPx` / `BadgeFontPx` / `BadgeBorderPx` / `QueueDotPx` | ball.html（`BadgeBorderPx` 为本票补录） |
| 字体 | `--font-sans` CJK 头 + `--t-mono` 12.5 / `--t-micro` 11；环下倒计时 10px | `FontFamily` / `FontSizeMonoPx` / `FontSizeMicroPx` / `CountdownFontPx` | DWrite 文本（`CountdownFontPx` 为本票补录） |
| 液态斑位置（票 62，无 CSS 对应） | 半径 0.72 / 0.62 / 0.50 × 球半径；偏心 0.22 / 0.30 / 0.40 | `LiquidRadiusA` / `LiquidRadiusB` / `LiquidRadiusC` / `LiquidOffsetA` / `LiquidOffsetB` / `LiquidOffsetC` | 三枚软场叠加才读成「液体」；渲染器只旋转与胀缩，不重建 brush |
| 玻璃边缘权重（票 62，无 CSS 对应） | 外缘暗环 1.2px / 内亮唇 1.0px / 焦散 0.30 / 「未说话」边框环 1.8px（96 DPI 物理 px，绘制时按 DPI 缩放） | `GlassRimPx` / `GlassLipPx` / `GlassCaustic` / `BorderRingPx` | 浅色桌布上的对比度锚点 |
| 音频包络 → 液体运动（票 62，无 CSS 对应） | 每单位电平：偏心增益 0.55、转速增益 1.0 | `SwimLevelGain` / `SpinLevelGain` | 无分配：只改已有 brush 的几何 |
| 边框与流动时序（票 62/64，无 CSS 对应） | 边框淡入 220ms / 收声汇聚 180ms / 唤起一次性流动 900ms | `BorderOpenMs` / `BorderCloseMs` / `SummonFlowMs` | 180ms 属 `--dur-base` 档；220ms 属 `--dur-slow` 族；900ms 是一次性动效，不受「UI 过渡 ≤300ms」预算约束（SPEC-08 §2 时长预算） |
| 边缘吸附（票 64，无 CSS 对应） | 泊靠动画 160ms / 泊靠后仍露出 0.42 直径 / 距边 16px 内开始挤压（96 DPI，按显示器缩放） | `DockAnimMs` / `DockOverlapFrac` / `DockTriggerPx` | 挤压量按间隙距离读出差分，不需要定时器（D32 空闲零定时器纪律） |
| 缓动 | `cubic-bezier(0.32,0.72,0,1)` | D2D 侧以分段线性近似（S1），周期动画用 ease-in-out 正弦近似 | 见 renderer |

**值漂移（本票不自行裁定，列给 owner）**

1. **休眠体尺寸三方不一致**：`docs/specs/SPEC-08-ui-ball-panel.md` §2 INTERIM 写「静态玻璃体，直径
   44px」；`internal/ball/statevisual.go::stateSize` 在原型模式下实算 `用户尺寸 × SleepRestRatio
   (0.62)`，默认档 56 → **34.72px**（下限 30px，不是 44px）；本表此前记 **12px**（INTERIM 前的原值，
   已在上面保留标注）。涉及 commit：原值 `8464da6`（2026-09-19，票 07 建表），INTERIM 改动 `fd8f838`
   （2026-09-20，票 62「checkpoint the glass Sleeping body」）。SPEC 属冻结文件、`internal/` 本票不得
   触碰，故两边都只登记不改。
2. **INTERIM 描述的球体在生产路径不生效**：`prototypeVisuals` 声明为 `var prototypeVisuals bool`
   （默认 **false**），只有 `cmd/balldebug/main.go:104` 与 winlive 测试调用
   `EnablePrototypeVisuals(true)`。⇒ 默认构建下 `Sleeping` 仍画 12px 微点，即 SPEC-08 §2 INTERIM
   目前只对调试器成立。登记为发现，非文档笔误。
3. **票 62 的 36 个 look 色值没有 CSS 真相源**：`tokens.go` 的 `looks` 表（aurora / glacier /
   nebula / solar × 9 字段）是原生侧自有字面量，`tokens.css` 中无同名变量 —— `TestNoHardcodedColorsInBallPackage`
   放过它（字面量确在 `tokens.go` 内），但表头「权威真相源 = tokens.css」这条契约对它不成立。
   下表按代码逐值补录，CSS 两列如实标 `无（原生自有）`。

## 票 62 液态玻璃 look 色板（`tokens.go` 的 `looks` 表；无 CSS 真相源）

`LiquidLook` 的 9 个字段 × 4 套（aurora / glacier / nebula / solar），`looks[0]`（aurora）是
`DefaultLook()`，选择走 `SetLook(name)` / `-look` 调试旗标。**下表的 CSS 两列如实为空**：这 36 个
字面量不存在于 `tokens.css`，是本票对账时补录的原生自有值（见上「值漂移」第 3 条）。

| CSS 变量 | tokens.css 值 | Go 字段 | Go 值 |
|---|---|---|---|
| （无，原生自有） | —— | `looks[aurora].BlobA` | `hex(0x4B49FF,0.95)` |
| （无，原生自有） | —— | `looks[aurora].BlobB` | `hex(0x8B5CF6,0.90)` |
| （无，原生自有） | —— | `looks[aurora].BlobC` | `hex(0x22D3EE,0.85)` |
| （无，原生自有） | —— | `looks[aurora].Deep` | `hex(0x1E1B7A,0.60)` |
| （无，原生自有） | —— | `looks[aurora].Rim` | `hex(0x141A2E,0.72)` |
| （无，原生自有） | —— | `looks[aurora].Lip` | `hex(0xE8ECFF,0.75)` |
| （无，原生自有） | —— | `looks[aurora].Hi` | `rgba(255,255,255,0.92)` |
| （无，原生自有） | —— | `looks[aurora].Caustic` | `hex(0x1B1F3A,0.35)` |
| （无，原生自有） | —— | `looks[aurora].Glow` | `hex(0x6D7CFF,0.55)` |
| （无，原生自有） | —— | `looks[glacier].BlobA` | `hex(0x22D3EE,0.95)` |
| （无，原生自有） | —— | `looks[glacier].BlobB` | `hex(0x2DD4BF,0.85)` |
| （无，原生自有） | —— | `looks[glacier].BlobC` | `hex(0x1D4ED8,0.90)` |
| （无，原生自有） | —— | `looks[glacier].Deep` | `hex(0x0B3B66,0.60)` |
| （无，原生自有） | —— | `looks[glacier].Rim` | `hex(0x0E2233,0.72)` |
| （无，原生自有） | —— | `looks[glacier].Lip` | `hex(0xE6FBFF,0.70)` |
| （无，原生自有） | —— | `looks[glacier].Hi` | `rgba(255,255,255,0.92)` |
| （无，原生自有） | —— | `looks[glacier].Caustic` | `hex(0x0B2233,0.34)` |
| （无，原生自有） | —— | `looks[glacier].Glow` | `hex(0x38BDF8,0.52)` |
| （无，原生自有） | —— | `looks[nebula].BlobA` | `hex(0xD946EF,0.92)` |
| （无，原生自有） | —— | `looks[nebula].BlobB` | `hex(0x6366F1,0.90)` |
| （无，原生自有） | —— | `looks[nebula].BlobC` | `hex(0xF472B6,0.80)` |
| （无，原生自有） | —— | `looks[nebula].Deep` | `hex(0x3B0764,0.55)` |
| （无，原生自有） | —— | `looks[nebula].Rim` | `hex(0x1A1030,0.74)` |
| （无，原生自有） | —— | `looks[nebula].Lip` | `hex(0xFCE7FF,0.70)` |
| （无，原生自有） | —— | `looks[nebula].Hi` | `rgba(255,255,255,0.92)` |
| （无，原生自有） | —— | `looks[nebula].Caustic` | `hex(0x241033,0.34)` |
| （无，原生自有） | —— | `looks[nebula].Glow` | `hex(0xC026D3,0.50)` |
| （无，原生自有） | —— | `looks[solar].BlobA` | `hex(0xF59E0B,0.92)` |
| （无，原生自有） | —— | `looks[solar].BlobB` | `hex(0xFB7185,0.88)` |
| （无，原生自有） | —— | `looks[solar].BlobC` | `hex(0x7C3AED,0.80)` |
| （无，原生自有） | —— | `looks[solar].Deep` | `hex(0x4A1D3A,0.55)` |
| （无，原生自有） | —— | `looks[solar].Rim` | `hex(0x241019,0.72)` |
| （无，原生自有） | —— | `looks[solar].Lip` | `hex(0xFFEFD6,0.70)` |
| （无，原生自有） | —— | `looks[solar].Hi` | `rgba(255,255,255,0.92)` |
| （无，原生自有） | —— | `looks[solar].Caustic` | `hex(0x2A1206,0.32)` |
| （无，原生自有） | —— | `looks[solar].Glow` | `hex(0xFDBA74,0.50)` |

## 审计

- `go test ./internal/ball/ -run TestNoHardcodedColorsInBallPackage`：`internal/ball` 包内除
  `tokens.go` 外任何 `.go` 文件命中 `rgba(` 或 6 位 hex 字面量即 FAIL。
- `go test ./internal/ball/ -run TestTokenGoldenValues`：抽查上表关键值与 CSS 相等（防抄写漂移）。
  边界（本票对账时实测）：该测试的 20 条断言全是 **`Palette` 色值**，几何/动效常量与本 look 表
  **不在金标准覆盖内** —— 抄写漂移只靠本表人工核对。登记为发现，改断言属代码改动，本票未做。
- 反向对账（本票 AC#7 用过，可复跑）：把 `tokens.go` 的 40 个 `Palette` 字段、全部导出几何常量的
  字面量与 `tokens.css` 声明逐条比对，`design/assets/tokens.css` 129 条声明中 80 条无原生对应
  （见「范围」）。核对结论：色值 78 行（明 40 + 亮 38）**与 CSS、`tokens.go` 三方逐字一致，零漂移**；
  几何/动效缺 25 个常量，本票已全部补录。
- 票 64 的 hotkey 面**不产 token**（本票查过，免下一个代理重查）：`internal/ball/hotkey_windows.go`
  三处 `const` 块是 Win32 修饰键/虚拟键码（`modAlt=0x0001` … `vkEscape=0x1B`）与 `HotkeyStatus`
  枚举，既无尺寸也无颜色；`hkNames` 的用户可见文案零 emoji。票 64 进本表的只有边缘吸附三常量
  （`DockAnimMs` / `DockOverlapFrac` / `DockTriggerPx`）。
- **零 emoji 扫描（AC#7 前半）**：仓内唯一机器执行的 emoji 检查是 `tools/d22scan` 的 ban `emoji`
  （`tools/d22scan/main.go:71` 的 `emojiRe` + `walkEmoji`），**scope 只有 `design/` 与 `frontend/`
  两处**（`main.go:137-142`），且 `frontend/` 在当前 HEAD **不存在** ⇒ 该门今天实际只覆盖 `design/`。
  它**不扫** `internal/`、`cmd/` 里的 Go 字符串字面量，所以原生侧 UI 文案（托盘提示、CLI 输出、
  通知正文）落在门外 —— AC#7 那句「over new UI strings」的前提只成立了一半。
  实跑结论与设计缺口记在票 12 的 Progress log；把 scope 扩到 `internal/`+`cmd/` 属代码改动，本票未做。
