# C21 原生侧 token 对照表（ticket 07）

`s1/` 证据 · 生成：T07-impl · 2026-09-19

权威真相源：`design/assets/tokens.css`（C21 契约参考实现）。原生侧（`internal/ball/tokens.go`）
逐值复制下表；`docs/evidence` 本表供人工核对。规则：球体绘制代码出现任何硬编码色值即违规
（SPEC-08 §2）—— 该约束由 `TestNoHardcodedColorsInBallPackage` 机器执行（仅 `tokens.go` 允许
出现 `rgba(` / 6 位 hex 字面量）。

值格式：`rgba(r,g,b,a)` 为 CSS 通道；Go 列为 `tokens.go` 中的等价构造（`rgba()` 0..255 通道 +
alpha，`hex()` 0xRRGGBB + alpha）。D2D 使用直通 alpha 的 `D2D1_COLOR_F`，premultiplied 仅在
`UpdateLayeredWindow` 位图写出处派生（`Color.Premultiplied()`）。

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
| SPEC-08 §2（Sleeping 微点） | 12px / opacity 0.35 | `SleepingDotPx` / `SleepOpacity` | 不随用户尺寸放大 |
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
| 帧率上限 | ≤30fps | `MaxAnimFPS` / `MinFrameMs` | 所有动画定时器周期 ≥33ms |
| 描边宽 | 1.5px 斜杠 / 图标；Conversation 环 2px | `SlashStrokePx` / `IconStrokePx` / `RingStrokePx` / `ConvRingStrokePx` | SPEC-08 §2.1 |
| 角标 / 队列点 | 17px 角标、10px 数字、9px info 点 | `BadgeDiameterPx` / `BadgeFontPx` / `QueueDotPx` | ball.html |
| 字体 | `--font-sans` CJK 头 + `--t-mono` 12.5 / `--t-micro` 11 | `FontFamily` / `FontSizeMonoPx` / `FontSizeMicroPx` | DWrite 文本 |
| 缓动 | `cubic-bezier(0.32,0.72,0,1)` | D2D 侧以分段线性近似（S1），周期动画用 ease-in-out 正弦近似 | 见 renderer |

## 审计

- `go test ./internal/ball/ -run TestNoHardcodedColorsInBallPackage`：`internal/ball` 包内除
  `tokens.go` 外任何 `.go` 文件命中 `rgba(` 或 6 位 hex 字面量即 FAIL。
- `go test ./internal/ball/ -run TestTokenGoldenValues`：抽查上表关键值与 CSS 相等（防抄写漂移）。
