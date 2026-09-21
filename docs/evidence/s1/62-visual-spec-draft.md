# 62 — 悬浮球 20 态视觉映射表（**契约草案，等 owner 签字**）

**Ticket:** `.scratch/wisp/issues/62-liquid-glass-ball-visuals.md` AC#7（正文 `:51`、第 7 框）
**登记项:** `docs/reports/pending-and-issues.md` **A30④** —— 本文件此前被票面与 AC#7 引用却**从未被写出**（`find . -name "62-visual-spec-draft*"` 零命中）。本文件闭的就是这个缺口。
**状态:** **草案 / 未签字 / 不是契约**。SPEC-08 §2、§2.1 与 `docs/PLAN.md` 的 D43 表**在本文件里一个字都没改**（本文件也**不构成**签字：签字后由编排者回填 `docs/specs/SPEC-08-ui-ball-panel.md` §2.1，届时本文件降级为来源附注）。
**Machine（全部数字的口径）:** DESKTOP-LVS7839 · Windows 10 Pro 26100 · 单屏 **3440×1440 @ 96 DPI** · **浅色（白）桌布** · 配置尺寸 **56**（默认档）。⇒ **A28 的阴影在本表每一行都成立**：没有一条像素数字是在 >96 DPI 或深色桌布上量的。
**态名来源:** 逐字取自 D43（`docs/PLAN.md:3053-3094`）与 `internal/statemachine/states.go:11-30`；行序取 SPEC-08 §2.1 的现有行序，便于签字时逐行对调。**20 行，无增无减。**

---

## 0. 读法（三档标记 + 两个渲染模式 + 公共几何量）

### 0.1 三档标记（每行末尾一列，不许混用）

| 标记 | 含义 | 必须给什么 |
|---|---|---|
| **〔实测〕** | 该格数字来自差分化像（活帧 vs 死帧逐像素相减） | 目录/表列名 |
| **〔代码真值，未实测〕** | 只从码里读出，没有任何像素/资源实测 | `file:line` |
| **〔待 owner 定〕** | 契约与码/实测**互相矛盾**，本表**不选值** | 问题编号（§3 的 Q-*）+ registry 编号 |

### 0.2 两个渲染模式（这张表有两条"当前真值"，因为构建里有两条）

`prototypeVisuals`（`SV:102`）决定画哪一套：

| 模式 | 谁开着它 | 画什么 |
|---|---|---|
| **冻结档（frozen）** | **库默认值 = false ⇒ 默认构建就是它**（全仓唯一开启者是 `cmd/balldebug/main.go:104` 的 `EnablePrototypeVisuals(!*frozen)`） | SPEC-08 §2.1 原文那 20 行：纯色玻璃滴 + 状态色环，`Sleeping` = **12px 微点 / opacity 0.35**（`SV:128-137`、`TK:379-380`） |
| **原型档（prototype）** | `balldebug` 不带 `-frozen` 时（= owner 2026-09-20 用肉眼签收的那颗） | 票 62 的玻璃液态球（`SV:272-303` 的 `applyGlassForm` + `REN:510-588` 的 `drawGlass`） |

⇒ **本表每一行的"透明度/颜色"都给两档**。registry **A24-D2** 的原话成立：**"owner 签收的球不在默认构建里"**，所以 SPEC-08 §2 INTERIM 现在描述的是一个**非默认配置**。这是 §2 冲突表的第一条（C-01），也是票 68 AC#2 的靶子。

⚠ 一个容易读错的细节：**原型档不是"替换"§2.1，而是"叠加"**。`drawFrame` 的第 6–13 层（状态环 / Listening 波环 / Thinking 光带 / 进度环 / 图标 / 角标 / 队列点 / 中心文字）在 `REN:426-500`，**在 `v.Glass` 分支之外**，两种模式都照画。所以"各态只换液体的颜色/流速/边框"（票 62 What-to-build 第 4 条）在码里**没有**替换掉描边与图标，只是又加了一层材质。

### 0.3 公共几何量（默认配置 56 / 96 DPI，一次性算清，逐行不再重复推导）

| 量 | 公式 | 真值出处 | @56（非 Sleeping 态） | @56 的 Sleeping（原型档） |
|---|---|---|---|---|
| 本体直径 `D` | 非 Sleeping = 配置尺寸（clamp 44–72）；Sleeping 原型 = `配置 × SleepRestRatio(0.62)`，下限 `SleepingRestMinPx=30`；Sleeping 冻结 = `SleepingDotPx=12` 固定 | `SV:121-139`、`TK:332/379/385` | **56.00px** | **34.72px**（冻结档 12px） |
| 窗口边长 `W` | `WindowEdgePx(配置, dpi)` = `(配置 + 2×RingMarginPx(8)) × dpi/96`；**全 20 态共用同一个窗，换态不重建** | `HIT:48-52`、`BW:229-236` + `BW:356-360` 注释 | **72px** | **72px**（窗不随态缩） |
| 可见外沿（光晕填充椭圆的几何半径，非对比度阈值） | `min(1.5 × D/2, W/2)` | `REN:534`（`fillEllipse(c, R*1.5, R*1.5, …)`） | `1.5×28 = 42` **但被窗裁到 36** | **26.04px**（未被裁） |
| 命中半径 `r` | `int(D) × dpi/96 / 2 + ClickTolerancePx(4) × dpi/96` | `HIT:23/27/35-44`，调用点 `BW:555`（传的是 `int(b.curVisual.SizePx)`，**不是**窗口半径也**不是**可见外沿） | **32px** | **21px** |
| 停靠挤压 | `DockSquash(1) = DockOverlapFrac = 0.42`，**只压液斑/高光/rim 的 X 向**；shell/halo/caustic 仍是整圆；`DockPos` 的 `orbR` 取的是**窗口**半径 | `DK:51-57`、`REN:546-557`、`DK:140-141` | 露出 **≈82%**（不是 42%） | 同；实测 `35×46` 框 |
| 帧率 | `MaxAnimFPS=30` ⇒ 任何动画周期 ≥33ms；`framePeriod = 周期/24` 再夹到 ≥33 | `TK:409-411`、`AN:73-80` | — | — |

⇒ 表里的 `r=32 vs 可见 36/42` 与 `r=21 vs 可见 26` 是**同一个缺口的两个读数**。registry **A29③** 只记了 Sleeping 那一格，并写了"**同形排查应扩展到全球包**"——本表把它扩展到全 20 态（Q-2），**不自行改值**。

### 0.4 差分化像证据索引（本文件所有〔实测〕格子只能来自这里）

| 目录 / 表 | 跑的是什么 | 覆盖了哪些态 |
|---|---|---|
| `62-diff-baseline/diff-table.txt` | **改造前**（冻结档） | Sleeping（px≥8 = **0**、max 4/255、无成像框、privWS 10.60MB、timers=no）、Listening 1471、Speaking 4824 |
| `62-diff-glass/diff-table.txt` | 原型档 | Sleeping **2098 / 40×38**、Listening 4855、Speaking 4839 |
| `62-diff-border/diff-table.txt` | 原型档 + 空闲边框 | Sleeping **2120 / 44×44**、timers=no、cpu 0.000、privWS 11.84MB |
| `62-diff-signoff/diff-table.txt` | 原型档 + 停靠两路（**签收跑**） | Sleeping **2120 / 44×44 / 11.76MB / timers=no**；Sleeping-dock-right **1620 / 35×46 / timers=no**；Listening 4856 与 dock 3230；Speaking 4848 与 dock 3220 |
| `docs/SLO.md` §A.2（表体来自 `docs/evidence/s1/12-ball-walk/diff-table.txt`） | 原型档会话走位 | Sleeping **2103 / 46×46**、Listening 4853、**Thinking 4844**、**Acting 4844**、Speaking 4851、**Warm 4625 / 71×69**、**Settling 4692 / 66×66** |
| `62-ball-visual-prototype.md` | `-tour` 14 步放映的读数表 + 未验证清单 | 同上六个态 + 停靠四边（目测） |
| `ball-states/*.png`（20 张，`1e80700`，2026-09-19） | **冻结档**、且是 `-shots` 的 **Go 侧合成图** | 名义上 20 态，但票 62 AC#1 明文"**不用自绘合成图**"⇒ **不能当 AC#1/#3 的证据**，只能当 §2.1 旧文本的示意图 |
| `internal/ball/*_test.go` | 确定性单元测试（不是像素） | 尺寸真相表 `TT:194-249`、反证 `TT:253-284`、单调性 `LQ` 侧 `liquid_test.go:21`、边框非单帧 `liquid_test.go:135` |

**缩写：** `SV`=`internal/ball/statevisual.go`｜`TK`=`internal/ball/tokens.go`｜`AN`=`internal/ball/anim.go`｜`REN`=`internal/ball/renderer_windows.go`｜`BW`=`internal/ball/ball_windows.go`｜`LQ`=`internal/ball/liquid.go`｜`LW`=`internal/ball/liquid_windows.go`｜`DK`=`internal/ball/dock.go`｜`HIT`=`internal/ball/hit.go`｜`TT`=`internal/ball/tokens_test.go`。行号按**工作树**（HEAD `f342413` 之后）标注。

---

## 1. 20 态映射表（逐态一行；D43 态名逐字照抄）

颜色列写**原型档**（签收的样子），token 值取自 `TK:111-159`（暗色板）/ `TK:162-207`（亮色板）；"aurora 液面" = `BlobA #4B49FF@.95 / BlobB #8B5CF6@.90 / BlobC #22D3EE@.85`、`Rim #141A2E@.72`、`Lip #E8ECFF@.75`、`Glow #6D7CFF@.55`（`TK:240-246`）。`mix = mixLook`：aurora Glow 与态色按 **45%** 混，A 通道保留 aurora 的（`SV:307-315`）。原型档的态身份**只走 `GlowColor`**，`CoreColor/CoreAlpha` 沿用冻结档那一行（`SV:272-303` 是后置覆写）。

| # | 态（D43） | 尺寸：公式 ⇒ @56 | 透明度 冻结 ⇒ 原型 | 颜色 / 渐变（原型档） | 动效（有无 · 速率 · 周期） | 描边 / 叠加 | 零定时器要求 ⇒ 现状 | 当前实现真值 | 与 SPEC-08 §2.1 现文本 | 档 |
|---|---|---|---|---|---|---|---|---|---|---|
| 1 | `Sleeping` | `56×0.62` ⇒ **34.72px**（下限 30；冻结档 **12px 固定**）；窗 72 | `0.35` ⇒ **`1.0`** | aurora 液面（静止不转）+ halo `mulA(Glow,0.45)`；**无边框** | **无**（静态单帧）；电平在 `LW:51/56-64` 处被直接丢（`liquidDriven` 不含它） | 玻璃 rim 1.2px + lip 1.0px；**无** `BorderRing`（`borderAtRest=0`） | **要求=是（D32/票 07 AC#4）** ⇒ `AN:61-67` AnimNone，**实测 timers=no 四次**（baseline/glass/border/signoff），cpu 0.000/0.000，privWS 11.76–11.91MB | `SV:128-136`＋`SV:155-157`＋`SV:283-287`；差分歧列 `TT:194-249` | **原型档=一致**（§2.1 已按 34.72 更正）；**冻结档=不一致**（默认构建仍画被否证的 12px/0.35） | 〔实测〕成像 +〔待 owner 定〕**Q-1/Q-2/Q-3/Q-4** |
| 2 | `Armed` | `56×1` ⇒ **56.00px** | `0.6` ⇒ **`0.92`** | aurora 液面静止 + halo = **aurora Glow 原样**（`CoreAlpha=0` ⇒ 不进 mix）⇒ **态色缺席** | 无（静态） | rim/lip + **`BorderRing 1.8px` 边框（单帧出现，无过渡）** | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:158-159`、`SV:288-289`、`SV:260-266`＋`LW:139-146`（边框为何是常量） | **两处不一致**：契约"直径 44px" vs 码 56（既有分歧，R15#2）；契约 0.6 vs 原型 0.92；边框是契约没有的新增维度 | 〔待 owner 定〕**Q-5/Q-1** |
| 3 | `Muted` | `56×1` ⇒ **56.00px** | `0.4` ⇒ **`0.85`** | aurora 液面静止 + halo = Glow 原样；斜杠走 `FgSecondary #A6A9B0` | 无 | 斜杠 `SlashStrokePx=1.5px` 对角线（±0.45R，`REN:667-671`）+ 1.8px 边框（常量） | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:160-162`、`SV:290-291`、`TK:414` | 斜杠 **一致**；透明度 **不一致**（0.85 vs 0.4） | 〔待 owner 定〕**Q-1** |
| 4 | `Listening` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | halo = `mix(Glow, **BallHalo** rgba(255,255,255,.10))` —— 该 case 不设 `CoreColor`，所以**混的是默认白雾，不是 accent**；accent 只活在环与图标上 | **有**：AnimFrame，周期 `ListeningWaveMs=2400` ⇒ 周期/24 = **100ms（10fps 有效）**，三层 0.8s 错峰波环（**时基**）；电平另驱动 `halo permille += 150×level`（`REN:524`）与转速 `0.45+4.20×level rad/s`（`LQ:27-28`） | accent 环 `RingStrokePx=1.5px` @R+5；`IconAudioLines` 白描（1.5px）；1.8px 边框（`transitionDriven` 含它 ⇒ 走 ramp，非单帧） | 要求=否（白名单内）⇒ **实测 timers=yes**，cpu 0.023/0.034/0.000，px≥8 **4856 / 4855 / 4853**（三次原型跑；冻结档那条是 1471），框 72×72 | `SV:163-166`、`AN:42-43`、`REN:438-447`、`TK:398/416` | **部分不一致**：契约"外环**随音量**波动"，码里波动环是**时基**、电平只改光晕与转速；且 halo 42px 被 72px 窗裁到 36（C-18） | 〔实测〕 +〔代码真值，未实测〕音量耦合 |
| 5 | `Thinking` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Accent #86C2B9 / CoreAlpha .35` ⇒ halo `mix(Glow, Accent)`；三液斑照转 | **有**：AnimFrame，`ThinkingSweepMs=1200` ⇒ 50ms 夹到 **33ms（30fps）**，底部 2px 光带左→右；液面随电平 | 无状态环；光带 `ThinkingBandPx=2`；1.8px 边框 | 要求=否 ⇒ AnimFrame ✓；**实测 timers=yes、px≥8 4844、框 72×72、cpu 0.016** | `SV:167-172`、`AN:44-45`、`REN:450-458`、`TK:399/418` | **一致**（契约只写光带，未写液面 ⇒ 液面是新增维度 C-15） | 〔实测〕 |
| 6 | `Acting` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Accent / CoreAlpha .95` ⇒ halo `mix(Glow, Accent)`、`GlowColor AccentLine` 被覆盖 | **冻结档无**（`AN:50-53` 明写 AnimNone："允许但不需要"）；**原型档可有** 33ms 涌动 burst（`LQ:78-83` `transitionDriven` 含它，`LQ:299-307` `busy()`）⇒ 一次进态会自己退休 | 无环；1.8px 边框 | 要求=否但**实测 timers=yes**（`12-ball-walk`，cpu 0.003）⇒ 那个 timer 是票 62 的 burst，不是票 07 的动画 timer | `SV:173-176`、`AN:50-53`、`LQ:62-70/78-83` | **轻微不一致**：契约"opacity 1 常亮实心"，码 `CoreAlpha .95` 且原型档**在动**（不"常亮"） | 〔实测〕 +〔代码真值，未实测〕 |
| 7 | `Confirming` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | 无核心色（`CoreAlpha=0`）⇒ halo = aurora Glow 原样；身份全在 danger 环 | **有**：AnimFrame，`ConfirmingPulseMs=2000` ⇒ 83ms；环 alpha = `0.55+0.45·(0.5+0.5cos(2πφ))` —— **每 2s 无限循环**（`REN:428-435`） | danger 环 1.5px @R+5（**脉冲**）；**倒计时**：`BadgeText` 走 `microFmt = FontSizeMicroPx 11px` 且画在**球心**（`REN:497-500`）；`CountdownFontPx=10` **零消费者**（全仓 `TK:430` 一处） | 要求=否 ⇒ AnimFrame ✓；**成像未实测** | `SV:177-180`、`AN:46-47`、`TK:400/430` | **三处不一致**：契约"仅一次或低频" vs 无限循环；"环下"vs 球心；micro 字号与 `CountdownFontPx` 无对应 | 〔代码真值，未实测〕 +〔待 owner 定〕**Q-6** |
| 8 | `AwaitingApproval` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | 无核心色 ⇒ aurora Glow 原样 | **无**（`AN:61-67` 默认分支：静态） | danger 环 **1.5px 常亮**（无 pulse 位）+ 右上**角标 `BadgeDiameterPx=17px`**、2px 描边、白字 `BadgeFontPx=10`；深度默认 **`BadgeCount=1`**，真值要 `SetBadge` | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:181-184`、`TK:421-424`、`BW:310-316`（`SetBadge` **零生产调用者**，仅 `cmd/balldebug/main.go:304`） | **两处不一致**：契约"12px 圆"vs 码 17px；"队列深度"未接线 ⇒ 永远显示 1 或调试器的 3 | 〔代码真值，未实测〕 +〔待 owner 定〕**Q-6** |
| 9 | `Speaking` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Success #7DB896 / .95` ⇒ halo `mix(Glow, Success)`（`SuccessLine` 被覆盖） | 冻结档：AnimFrame `SpeakingBreathMs=1600` ⇒ **66ms** 呼吸。**原型档：呼吸量取自 `v.LiquidAngle`（`REN:521-523`），而 `drawGlass` 不接收 `phase`** ⇒ 无电平、burst 退休后 `LiquidAngle` 停止累积 ⇒ **1.6s 周期在原型档不成立**，66ms 的动画 timer 每帧画出同一张图 | `IconVolume` 白描 1.5px；**无边框**（`borderAtRest(Speaking)=0`，说话时液体接管球体，`SV:262-264`） | 要求=否 ⇒ 实测 timers=yes，cpu 0.016/0.006/0.047；px≥8 **4848/4851/4839**、框 72×72 | `SV:185-189`、`SV:262`、`AN:48-49`、`REN:519-521`、`TK:397` | **不一致**：契约"success 呼吸（1.6s）"在原型档失义（**这条是本表最像缺陷的一条**，码上可判、无人量过） | 〔实测〕成像 +〔代码真值，未实测〕呼吸 +〔待 owner 定〕**Q-7** |
| 10 | `Settling` | `56×1` ⇒ **56.00px** | 冻结 `1→0.35` ⇒ 原型 `Visual 1→**0.95**`（`SV:292-295`，落点 `RestSettledOpacity=0.95`）**但 `BW:471/474` 在两种模式下都把 ULW `constAlpha` 驱到 `SleepOpacity=0.35`** ⇒ 合成 **0.95×0.35 ≈ 0.33** | aurora 液面（会话已退出 ⇒ 静止）+ 边框**淡入**（`transitionDriven` 含它） | **有**：AnimFade 一次性 `SettlingFadeMs=260` ⇒ 33ms，到 1 就 `stopAnimTimer`（`AN:57-60`）⇒ 不是循环 | 1.8px 边框渐入（实测逐帧 `0→0.11→0.32→0.53→0.74→0.96→1`，winlive `TestBallLiveIdleBorderTransition`，**代理自述、本表未复跑**） | 要求=否 ⇒ 实测 **timers=yes、px≥8 4692、框 66×66、cpu 0.000** | `SV:190-192`、`SV:292-295`、`BW:468-476`、`TK:386/402` | **不一致（A29⑤ / R15#3）**：契约"opacity 1→0.35"字面只在冻结档成立；原型档合成 ≈0.33，既非 0.35 也非 0.95 | 〔实测〕 +〔待 owner 定〕**Q-1** |
| 11 | `Warm` | `56×1` ⇒ **56.00px** | `1` ⇒ `1`（Visual 不动）；呼吸走 **ULW constAlpha 0.55↔0.70**（`BW:462-467`） | `CoreColor Warm #D6B184 / .9` ⇒ halo `mix(Glow, Warm)`；`WarmLine` 被覆盖 | **有**：AnimAlpha，`WarmBreathFPS=10` ⇒ **100ms**、周期 `WarmBreathPeriodMs=2400`，余弦 ease-in-out；**不重画像素，只改窗口 alpha** | 1.8px 边框（`transitionDriven` 含 Warm ⇒ 走 ramp） | 要求=否 ⇒ **实测 timers=yes、px≥8 4625、框 71×69、cpu 0.012** | `SV:193-196`、`AN:54-56`、`TK:396/405-406/411`、`BW:462-467` | **一致**（SPEC §2 明文允许"Warm ≥2s opacity 呼吸走合成器"，§2.1 的 2.4s/0.55↔0.7 逐值对上） | 〔实测〕 |
| 12 | `Conversation` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Danger #E07A70 / .95` ⇒ halo `mix(Glow, Danger)` | **无定时器**（`AN:61-67` 默认分支）⇒ 但 `liquidDriven` 含它（`LQ:62-70`）⇒ **液面只由电平回调推进**；无喂送则静止，且 `rampBorder` 只在 `pushLevel` 里跑 ⇒ **静音时边框能否到位取决于采集是否还在发样本**（`LQ:143-150/190-245`） | **danger 常亮环 `ConvRingStrokePx=2.0px`**（> Confirming 的 1.5 ⇒ "不得弱于"成立）+ `IconMic`；1.8px 边框 | 要求=是（"不得渐隐"）⇒ AnimNone ✓；**成像与电平驱动均未实测** | `SV:197-204`、`TK:417`、`LQ:62-70/119-128/143-150` | **一致**（常亮、2px、不渐隐三条都对上）；"边框依赖电平回调"是契约没有的**行为风险** | 〔代码真值，未实测〕 +〔待 owner 定〕**Q-8** |
| 13 | `Downloading` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor BGOverlay rgba(30,36,44,.50) / .7` ⇒ halo `mix(Glow, BGOverlay)` | **无**定时器；`Progress` 默认 `0.0` ⇒ 只画**轨道环**（无弧），有弧要靠 `SetProgress` 逐次重绘 | 进度环：`BorderSoft` 轨道 + `Accent` 弧，2px @R+5（`REN:461-469`）；百分比文字走 `BadgeText`（`SetBadgeText`） | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:205-208`、`BW:318-324`、`BW:326-332`（`SetProgress`/`SetBadgeText` **零生产调用者**，仅 `main.go:307-308`） | **未落地**：契约"环形进度 + 百分比"没有数据源接线（票 26/D26 侧未装配）⇒ 规格成立、能力为零 | 〔代码真值，未实测〕 +〔待 owner 定〕**Q-6** |
| 14 | `Error` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Danger / .95` ⇒ halo `mix(Glow, Danger)` | **无** | `IconX` 白描 1.5px（两条 ±0.35R 对角线）+ 1.8px 边框（常量） | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:209-213`、`REN:672-676` | **一致**（"danger 底 + x 图标"逐字对上；边框是新增维度） | 〔代码真值，未实测〕 |
| 15 | `NoNetwork` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor `**`TintNeutralMid #6A6D75`**` / .8` ⇒ halo `mix(Glow, TintNeutralMid)` | **无** | `IconWifiOff`（两道弧 + 斜线）+ 1.8px 边框 | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:214-217`、`TK:103/68` | **不一致（token 走错门）**：契约"**fg-tertiary** 底"（`#71747C`），码用 `--tint-neutral-mid`（`#6A6D75`）⇒ 两枚不同 token，C21 对账表里也是两行 | 〔代码真值，未实测〕 +〔待 owner 定〕**Q-9** |
| 16 | `Unconfigured` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Warn #D9B26A / .95` ⇒ halo `mix(Glow, Warn)` | **无** | `IconKey`（key-round）+ 1.8px 边框 | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:218-222` | **一致** | 〔代码真值，未实测〕 |
| 17 | `FirstRun` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Accent / .95` ⇒ halo `mix(Glow, Accent)` | **无**：`FirstRunGuidePulseMs=2600` 只是 token，`AN:61-67` 走默认静态分支，注释明写"FirstRun's guide pulse … land with their wiring tickets; S1 keeps them static per the whitelist" | `IconArrowUpRight` + 1.8px 边框 | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:223-227`、`TK:401`、`AN:62-67` | **不一致（脉冲未实现）**：契约"accent 底 + **引导脉冲**"，脉冲这一半在码里不存在也没有票 62 的替代 | 〔代码真值，未实测〕 +〔待 owner 定〕**Q-6** |
| 18 | `WatchdogAlert` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Warn / .95` ⇒ halo `mix(Glow, Warn)` | **无** | `IconAlertTriangle` + 1.8px 边框 | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:228-232` | **一致** | 〔代码真值，未实测〕 |
| 19 | `Queued` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Accent / .95` ⇒ halo `mix(Glow, Accent)`；**没有"主态"** | **无** | 左下 `QueueDotPx=9px` **info** 点（+2px 描边），偏移 `R*0.8`（`REN:487-495`）；1.8px 边框 | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:233-240`（注释自陈："Queued is an overlay on the main state … **the 'main state' compositing lands with the session wiring**"）、`TK:424` | **两处不一致**：①契约"**主态** + 小点叠加"，码是**独立态**（进 `Queued` 就把主态视觉整张换掉，D43 #20 的语义"在等谁"丢了）；②点 **9px vs 契约 6px** | 〔代码真值，未实测〕 +〔待 owner 定〕**Q-10** |
| 20 | `Stuck` | `56×1` ⇒ **56.00px** | `1` ⇒ `1` | `CoreColor Warn / .95` ⇒ halo `mix(Glow, Warn)` | **无** | `IconRefresh`（330° 弧 + 箭头，`REN:722-729`）+ 1.8px 边框 | 要求=是 ⇒ AnimNone ✓；**未实测** | `SV:241-245` | **一致** | 〔代码真值，未实测〕 |

**停靠（dock）行不单独占态**——它是每一态的第二形态，几何见 §0.3、问题见 Q-1/Q-2。已量的 dock 行只有三条：`Sleeping-dock-right` **1620px / 35×46 / timers=no / cpu 0.000**、`Listening-dock-right` 3230 / 47×72、`Speaking-dock-right` 3220 / 47×72（`62-diff-signoff`）。其余 17 态的 dock 形态**零实测**。

---

## 2. 草案与 SPEC-08 现行文本的逐条冲突点

（"SPEC 位置"指 `docs/specs/SPEC-08-ui-ball-panel.md`；PLAN 的 D43/§2 只在需要时点名。**这一节是本文件最有价值的产出：它不合并、不概括。**）

| # | SPEC 位置与原文要点 | 草案（码 + 实测）说的 | 性质 | 归属 |
|---|---|---|---|---|
| **C-01** | §2 尺寸条 + §2.1 `Sleeping` 行：「静态玻璃体 …… 默认 **34.72px**」 | 默认构建（`prototypeVisuals=false`）画的是 **12px / 0.35**（`SV:128-131`、`TK:379-380`），全仓唯一开启者是 `cmd/balldebug/main.go:104` | **契约描述了一个非默认配置**（A24-D2 原话："默认构建画的不是签收的样子"） | 票 68 AC#2（翻默认值），**本表不翻** |
| **C-02** | §2 动画纪律：「**仅** Listening/Thinking/Acting/Speaking/Confirming 启动动画定时器」 | 实际有定时器的是 **七态 + 一个额外 timer**：Warm（AnimAlpha 100ms）、Settling（AnimFade 33ms）是 §2 自己后面认下的例外；票 62 的 33ms 液态/边框 burst（`LQ:78-83`）在**白名单内额外挂第二个 timer**，且 `Acting` 的实测 `timers=yes`（`12-ball-walk`）来自 burst 而非 §2 说的"动画" | **§2 与 §2.1 内部互斥**（§2.1 要求 Settling 260ms 渐隐与 Warm 呼吸，都需要帧），加上票 62 的新定时器族 | 建议回填时把"五态白名单"改写成"**五态循环动画 + Warm/Settling 两个 SPEC 认可的有界过渡**"，并写明 burst 只在同一集合内（Q-11） |
| **C-03** | §2 视觉一致性：「所有颜色/圆角/阴影/缓动**必须取自 C21**（`design/assets/tokens.css` 为**唯一**样式真相源）」 | 票 62 的 4×9 = **36 个 look 色在 `tokens.css` 里不存在**（原生首创），只有 `TK:240-265` 一个出处；`c21-native-tokens.md` 表头已按 D3 加了「覆盖面限定」，且**无机器检查**（D4） | **契约的"唯一真相源"对票 62 色板不成立** | 票 69（机器检查）+ 签字时把 look 表升进 `tokens.css` 或改 §2 措辞 |
| **C-04** | §2.1 `Armed`：「opacity 0.6，**直径 44px**，静态」 | 码：尺寸 = 配置 **56**（`SV:138`，`Armed` 不进 `stateSize` 的 Sleeping 分支）；opacity 冻结 `0.6` ✓、**原型 `0.92`** ✗ | **既有契约-代码分歧**（编排者已声明不自行动，R15#2）+ 原型档第二处偏离 | **Q-5** |
| **C-05** | §2.1 `Muted`：「opacity **0.4** + 1.5px 斜杠」 | 原型 `0.85`（`SV:290-291`）；斜杠 1.5px ✓ 且两档都画 | 原型档偏离（同 `Armed`/`Sleeping` 一族：玻璃材质靠高不透明度才在浅色桌布成像） | **Q-1**（"原型档各态静止 opacity 到底定几"） |
| **C-06** | §2.1 `Settling`：「opacity **1→0.35** 渐隐，260ms」 | 原型 `Visual.Opacity` 落在 **0.95**、ULW `constAlpha` 落在 **0.35**，**两者相乘 ≈ 0.33**（`BW:471/474` 两档共用 0.35；`SV:295` 用 0.95） | **A29⑤ / R15#3：两个 alpha 打架**，契约两边都不满足 | **Q-1 附属**；修哪一半属渲染变更 ⇒ 票 68 AC#2/AC#3 同批 |
| **C-07** | §2.1 `Listening`：「56px + 外环**随音量**波动 + accent 描边」 | 三层波环由 `animPhase`（时基，2.4s 周期）驱动（`REN:438-447`）；电平只驱动 halo permille（+150×level）与液斑胀缩/转速（`REN:524/543`、`LQ:224-228`） | **文字承诺的耦合方式与码不同**（视觉上有"随音量"的效果，但不是环在随音量） | 回填时改写措辞（Q-12），**不改码** |
| **C-08** | §2.1 `Speaking`：「success 呼吸（**1.6s**）」 | 冻结档 ✓（`SpeakingBreathMs=1600`）。原型档的呼吸幅度 = `cos(LiquidAngle·0.5)`（`REN:521-523`），`drawGlass` **签名里没有 `phase`** ⇒ 无电平且 burst 退休后**不呼吸**；同时 66ms 动画 timer 仍在逐帧重画同一张图 | **原型档 1.6s 不成立**（本表唯一一条"看起来像缺陷"的冲突） | **Q-7**（要么把 `phase` 传进 `drawGlass`，要么承认原型不呼吸并改契约） |
| **C-09** | §2.1 `Confirming`：「danger 脉冲环（2s，**仅一次或低频**）+ **环下** micro 倒计时数字」 | 环 alpha 每 2s **无限循环**（`REN:428-435`，无"仅一次"计数）；文字用 `BadgeText`，画在**球心**、字号 `microFmt=FontSizeMicroPx 11`；`CountdownFontPx=10` **全仓零消费者**；`SetBadgeText` **零生产调用者** | **三处不一致**（循环/位置/字号），外加"数字根本没有数据源" | **Q-6**（接线归票 21 段 2 / 票 12 装配） |
| **C-10** | §2.1 `AwaitingApproval`：「danger 环 + 右上角标（**12px 圆**，白色数字 = **队列深度**）」 | 角标 **17px**（`BadgeDiameterPx`）+ 2px 描边 + 10px 白字；深度 `BadgeCount=1` 常量，真值靠 `SetBadge`（零生产调用者，调试器给 3） | 尺寸不符 + D31 强制要求（"角标必须显示队列深度"）**未落地** | **Q-6** |
| **C-11** | §2.1 `Downloading`：「环形进度 + **百分比**」 | `Progress=0.0` ⇒ 只画轨道；弧与百分比都需 `SetProgress`/`SetBadgeText`（零生产调用者；`main.go:307-308` 是唯一喂数的人） | **未落地**（D26/票 26 侧未装配） | **Q-6** |
| **C-12** | §2.1 `NoNetwork`：「**fg-tertiary** 底 + wifi-off」 | 码用 `--tint-neutral-mid`（`#6A6D75`），契约指 `--fg-tertiary`（`#71747C`）——**两枚不同 token、两行不同的 C21 对账行** | 纯 token 错配，可机判 | **Q-9** |
| **C-13** | §2.1 `FirstRun`：「accent 底 + **引导脉冲**」 | 无脉冲：`FirstRunGuidePulseMs=2600` 只是常量，`AN:61-67` 明写 S1 保持静态、"等接线票" | **未落地**（有 token、无行为、无票面把它列成框） | **Q-6** |
| **C-14** | §2.1 `Queued`：「**主态** + 左下 **6px** info 小点叠加」 | 点是 **9px**；且 `Queued` 在码里是**独立态**（accent 芯 + 点），`SV:233-236` 注释自陈"主态叠加随会话接线落地" ⇒ D43 #20「面板可见在等谁」在球侧无载体 | 尺寸不符 + **叠加语义缺失** | **Q-10** |
| **C-15** | §2.1 **全 20 行**：没有任何一行提到"边框（`BorderRingPx=1.8`）"或"液面" | 原型档给**除 Sleeping/Speaking 外的每一态**加了一圈 1.8px `Lip` 边框 + 1px `Rim`（`REN:581-587`），由 `borderAtRest` 单点决定（`SV:260-266`）；且**静止系状态（Armed/Muted/Error/NoNetwork/Unconfigured/FirstRun/WatchdogAlert/Queued/Stuck/…）的边框是常量、单帧出现**——`frameVisual` 的进入条件 `active‖summon>0‖ownsBorder` 三皆假（`LW:139-146`）⇒ `applyTo` 不被调用，`BorderAlpha=1` 原样落地 | **票 62 引入的两个新维度完全不在 §2.1**；且与票 62 自己提案（`62-ball-visual-prototype.md` §7 第 4 条"禁止单帧跳变"）**自相矛盾** | 回填时新增"边框/液面"列（Q-13）；"静止系要不要也走过渡"需 owner 判 |
| **C-16** | §2.1 `Sleeping` 与票 62 提案 §7 第 5 条：靠边「收缩成贴边标签（**可见比例 0.42**）」 | `DockOverlapFrac=0.42` 只压液斑/高光/rim（`REN:546-557`），shell/halo/caustic 不压；`DockPos` 的 `orbR` 取**窗口**半径 28 而非态半径 ⇒ 落点窗口在屏 **47/72px**，静止球本体 **≈82%** 露出（A29②；实测佐证的 `35×46` 框） | **文档与实测相反**（且**票 62 自己的提案文本 carrying 了错数字** ⇒ 若照 §7 回填就把 0.42 固化进契约） | **Q-4**（R15#4） |
| **C-17** | §2「点击穿透：球体本体可点」 | 命中圆 `r = int(D)/2+4`，可见外沿 `1.5R`（被窗裁到 ≤W/2）⇒ Sleeping **21 vs 26**、其余各态 **32 vs 36/42**（frozen 档同形：`SV:130` 时 r=**10**、窗 72） | **A29③ 的推广**（registry 只记 Sleeping 一格，并写明"同形排查应扩展到全态"）；一圈 4–10px **看得见点不着** | **Q-2**（R15#5） |
| **C-18** | §2.1 `Listening`「外环」+ 各态"环/角标不越界"的隐含要求 | 56 配置下窗 **72**、中心 36：halo 名义半径 42、`Listening` 波环最大 `R+5+16 = 49` ⇒ **两者都被窗口裁掉**；`62-diff-signoff` 非 Sleeping 行的成像框正是 **72×72**（= 窗边，不是内容边界） | **裁切**：§2.1 画的"外环"有一半时间不可见；成像框列因此不可当尺寸读（A29① 已确立"框列浮动 ±6px、`px≥8` 才稳"） | 回填时应写明"W = 配置+16 ⇒ 环预算 = 8px"（Q-14），改窗宽属渲染变更 |
| **C-19** | §2 多显示器：「Per-Monitor V2 DPI 感知」 | `SizePx` 不对称：D2D 目标 `dpiX/dpiY=96` ⇒ 本体是物理像素不缩放，而描边/ring margin 乘 `dpi/96`、窗口走 `WindowEdgePx` ⇒ **144 DPI 下 56 的体装进 108 的窗**（`SV:66-75` 注释；A28） | **本表"@56 的实际像素"列只在 96 DPI 成立**；且 A28 连带命中测试。**零条 >96 DPI 实测** | **Q-3**（R15#6）；修法（纯函数断言 vs 真机跑）待选 |
| **C-20** | §2 尺寸条：「配置范围 44–72（默认 56）」+ §2.1 各行给的数全是 56 | 只有 **56** 档被逐态量过（`62-ball-visual-prototype.md` §6 明说）；**44 / 72 两档零实测、零截屏**。`SleepingRestMinPx=30` 让 44/48 配置下 Sleeping 本体恒为 30（`SV:133-135`、`TT:216-224`），**这条下限行为也没量过** | 契约允许域内 **2/3 的档位没有证据**（`BallSizeSmallPx=44` 还曾被误读成体径，A29①） | **Q-15**；签收前补跑（`-diff … -size 44 / 72`） |
| **C-21** | 非 SPEC，但同表冲突：`docs/SLO.md` §A.2 **节标题**仍写「真实悬浮球在 `Sleeping`（**44px** 静态玻璃体…）」 | A29① 已定量否证 44px，SPEC-08 §2 也已更正为 34.72px；SLO 的**标题**没跟着改（表体数字 2103/46×46 不受影响） | 文档间残留矛盾（`SLO.md` 不是冻结契约，但它在 SLO 门路径上） | 编排者的 `SLO.md`；本文件**不动它** |
| **C-22** | 非 SPEC：SPEC-08 §2 更正块与 A29① 都把推导证据写成「`docs/evidence/s1/68-*`」 | **该路径下零文件**（`find` 只命中 `.scratch/wisp/issues/68-…md`）。34.72 的真实推导只活在 `TT:194-249` + `TT:253-284` + 三个 `62-diff-*` 目录里 | **又一条指向不存在文件的引用**（与 A30④ 同形） | 建议编排者在 registry 记一笔或把引用改指 `TT`；本文件**不改 SPEC** |

---

## 3. 〔待 owner 定〕的问题清单（本表**不选值**，每条给"选项 / 不答的代价 / 对应 R15 / 对应 registry"）

| Q | 问题 | 选项 | 不答的代价 | R15 / registry |
|---|---|---|---|---|
| **Q-1** | 原型档的**静止系 opacity 阶梯**（Sleeping 1.0 / Armed 0.92 / Muted 0.85）与契约（0.6 / 0.4）差 0.3–0.45，玻璃材质在浅色桌布上正是靠这个才成像（`62-diff-baseline` 的 0 像素教训）。**以哪一边为准？** | ①按实测阶梯（0.35→1.0 全档上抬）改契约 ②按契约压回（并预期浅色桌布不可辨） ③只定 Sleeping 一档 | 票 65 返工时又要重猜一次；SPEC-08 §2.1 三行长期与默认构建不符 | R15#1 邻近 / A24-D1、A29① |
| **Q-2** | 命中半径与可见外沿（全 20 态的 `r=32 vs 36/42`，Sleeping `21 vs 26`；frozen `10 vs 72`）。 | 对齐可见外沿 / 对齐本体（放弃外沿可点）/ 维持 | 用户反复点边缘没反应，且最难归因（票 64 隐形边框同族） | **R15#5** / A29③ |
| **Q-3** | DPI：本体不放大而窗与描边放大；**全部证据是 96 DPI**。 | ①借高分屏跑一条我给的命令 ②先加纯函数断言锁住 ③推到 S3 | 笔记本用户那里"球小得像噪点" | **R15#6** / A28 |
| **Q-4** | 停靠露出实测 ≈**82%** vs 文档/常量 **42%**。要的是迅雷那种"半个球藏进边里"吗？ | 修码到 ≈50% / 把文档改成 82% | 票 62 的吸附效果与签收时看到的不是同一个东西；**且票 62 提案 §7#5 会把 0.42 固化进契约** | **R15#4** / A29② |
| **Q-5** | `Armed` 的"直径 44px"是既有契约文本，码自 S1 起用配置尺寸（默认 56）且有实测。 | 改契约文本为"配置尺寸" / 改码到 44 / 保持现状 | 每次对照规格都要先解释"为什么差 12px" | **R15#2** / A24-D1 |
| **Q-6** | `Confirming` 倒计时、`AwaitingApproval` 队列深度、`Downloading` 进度/百分比、`FirstRun` 引导脉冲——**四个格子的数据源/行为在生产侧不存在**（`SetBadge`/`SetProgress`/`SetBadgeText` 零非调试器调用者、脉冲无实现、`CountdownFontPx` 零消费者）。回填时写成"契约要求 + 未落地"还是"先降级措辞、落地后再升"? | 保留要求并挂残缺标记（推荐，与 A25 的"零证据判据必须挂在 gate 票上"同法）/ 改弱措辞 | 又一条"人人以为别人做了"的假达标；D31 的强制可视化要求落空 | 新登记（本表首次逐格点名） |
| **Q-7** | `Speaking` 的 **1.6s 呼吸在原型档失义**（呼吸值取自 `LiquidAngle`，`drawGlass` 不接 `phase`）。 | ①把 `phase` 接进 `drawGlass`（渲染变更）②承认原型不呼吸、改契约 ③并入票 65 质感返工 | 契约的 1.6s 数字是空的；且每 66ms 重画同一帧（实测 cpu 只 0.006–0.016%，不成预算问题，但成"看起来在动其实没动"的观感问题） | 新登记（本表发现） |
| **Q-8** | `Conversation` 的边框/液面推进**依赖电平回调而非定时器**：采集停摆 ⇒ 边框停在 0。 | 接受（省一个定时器，D32 友好）/ 给会话态一个有界 burst / 写进契约当"必须持续喂电平" | 开麦但静音时"红色常亮环 + 边框"这一组承诺可能不成立，而这正是 D43 #28 要求"麦克风开着的唯一交代" | 新登记（本表发现）；`LQ:119-128/143-150` |
| **Q-9** | `NoNetwork` 用 `tint-neutral-mid` 而契约指 `fg-tertiary`。 | 改码 / 改契约 | C21 双侧一致性核对（SPEC §8）每轮都要重新解释 | 新登记（本表发现）；关联 A24-D4 |
| **Q-10** | `Queued` 的"主态叠加"未实现（换态即丢主态视觉），点也 9px≠6px。 | 实现叠加（会话接线的票）/ 契约改为"独立态 + 9px 点" | D43 #20「在等谁」在球侧没有载体 | 新登记；归属需编排者指派 |
| **Q-11** | §2 的"五态定时器"条文与 Warm/Settling 例外、票 62 burst 三件事叠在一起。 | 按 C-02 的建议改写条文 / 维持 | 条文与现实不符 ⇒ 后来的代理会自己发明纪律 | 票 62 提案 §7 第 2/3 条已经这样写过 |
| **Q-12** | `Listening`"外环随音量波动"的措辞（环是时基，电平驱动光晕与转速）。 | 改措辞（推荐）/ 改码让环真随电平 | 每次验收都要口头解释 | 新登记 |
| **Q-13** | "边框 + 液面"这两个新维度要不要进 §2.1（以及静止系边框单帧出现算不算违约票 62 自己的"禁止单帧跳变"）。 | 新增两列 / 只写一句总则 / 不写 | 签字后票 65 又按"液体颜色/流速/边框"三变量重做，但没有格子可填 | 票 62 What-to-build 第 4 条 |
| **Q-14** | 环预算：窗 = 配置+16 ⇒ halo(42) 与波环(49) 被裁。 | 加宽窗口（渲染变更，成本进 SLO）/ 缩小环预算并改契约 / 如实写"环按窗裁切" | 成像框列继续被误读成尺寸（A29① 就是被这个坑过一次） | 新登记 |
| **Q-15** | **44 / 72 两档、深色桌布、>96 DPI 三个维度零实测**——签字前要不要补跑。 | 三条都补 / 只补深色桌布（AC#1 明文要求）/ 明知缺口先签 | AC#1 的"浅色**与深色**都清晰可见"至今只在浅色上成立 ⇒ 框不能勾 | 票 62 AC#1、`62-ball-visual-prototype.md` §6 |

---

## 4. 我在证据里**找不到数字**的格子（直说找不到，不编）

| 缺什么 | 涉及行 | 我找过哪里 | 结果 |
|---|---|---|---|
| **13 个态的成像数字**（`Armed`/`Muted`/`Confirming`/`AwaitingApproval`/`Conversation`/`Downloading`/`Error`/`NoNetwork`/`Unconfigured`/`FirstRun`/`WatchdogAlert`/`Queued`/`Stuck`）的 `px_delta_ge8` / 成像框 / CPU / 私有工作集 | 第 2、3、7、8、12、13、14、15、16、17、18、19、20 行 | `62-diff-{baseline,glass,border,signoff}/diff-table.txt`、`docs/SLO.md` §A.2、`docs/evidence/s1/12-ball-walk/diff-table.txt` | **一条都没有**。这四个目录的 `state` 列合起来只覆盖 7 个态。**没有替这些态估过任何像素数**。名义上存在的 `ball-states/*.png`（20 张）是**冻结档**、且是 `-shots` 的 **Go 侧合成图**，票 62 AC#1 明文不收合成图 |
| **44px / 72px 两档的逐态数字** | 全 20 行的"尺寸"列 | 同上 + `62-ball-visual-prototype.md` §6 | §6 自己写明"44–72px 里**只有默认 56 被逐态量过**"。⇒ 表里 44/72 档我只能给**公式**，给不出实测 |
| **深色桌布**的任何差分 | 全 20 行 | `docs/evidence/s1/ball-contrast/`（`sleeping-over-dark-taskbar.png`、`same-region-no-ball.png`，2 张，非差分表） | **没有深色桌布的差分跑**。AC#1 的"浅色**与深色**"只完成一半 |
| **>96 DPI** 的任何像素/命中数字 | 全 20 行（尤其 C-19） | `docs/SLO.md` §A.2/A.6、四个 `62-diff-*`、`66/` | 全部注明是本机 **3440×1440 @ 96 DPI**。找不到第二条 DPI 口径 |
| **`Settling` 合成 ≈0.33 的像素证据** | 第 10 行 | 同上 | 只有**码可推**（`BW:471/474` × `SV:295`）+ A29⑤ 的文字登记。**没有一次实测直接量过 0.33** ⇒ 该格标〔待 owner 定〕而不是〔实测〕 |
| **音频驱动液面"随电平单调"的像素/wav 注入实测（票 62 AC#4）** | 第 4、5、6、9、11、12 行的动效列 | `62-diff-*`（无 `-diff-level` 行）、`62-ball-visual-prototype.md`（只有 `-level` 目测 + 一句"单核 0.016%"）、`internal/ball/liquid_test.go:21` | **找不到 wav 注入的落盘实测**。现存最强证据是**确定性单元测试** `TestLiquidRotationIsMonotonicInLevel`（数学单调，不是像素单调）+ 票 62 进度 log 里代理的自述。⇒ AC#4 不因本表变为可勾 |
| **`Speaking` 原型档到底呼不呼吸**（C-08/Q-7） | 第 9 行 | 全部差分表（`Speaking` 行的电平 = 0，未喂 `-diff-level`） | 找不到。我只能给**码可判**的结论：`drawGlass` 不接收 `phase`，无电平 ⇒ `LiquidAngle` 不累积 |
| **`Armed`/`Muted`/`Error`… 的 dock 形态** | 全 20 行的停靠维度 | `62-diff-signoff` 的三条 `-dock-right`（Sleeping/Listening/Speaking） | 其余 17 态的 dock 成像数字**不存在** |
| **SPEC-08 §2 更正块引用的 `docs/evidence/s1/68-*`** | C-22 | `find . -name "*68*"` | **零文件**（只有 `.scratch/wisp/issues/68-…md`）。34.72 的推导真实落点在 `TT:194-249`/`TT:253-284` |
| **`BorderRingPx=1.8` / `GlassRimPx=1.2` / `GlassLipPx=1.0` 是否肉眼可辨**（边框与边缘光的"够不够"） | 第 1、2、3、14–20 行的描边列 | 四张差分表的 `max_delta`（171–211/255）只有 Sleeping/Listening/Speaking 三行 | 找不到逐态的边缘光读数。票 65（质感）要的就是这一格，**本表交不出** |

---

## 5. 建议回填 SPEC-08 §2.1 的文本（**提案，等签字；未回填**）

⚠ 本节全部内容为**待签字提案**，**不得**被解释为契约已变更。签字前 `docs/specs/SPEC-08-ui-ball-panel.md` 与 `docs/PLAN.md` 保持现状（票 62「契约变更流程」原文：先原型 → owner 实机签收 → 才回填并签字）。

建议把 §2.1 表**从 1 列"视觉"改为 5 列**：`本体尺寸（公式）` / `静止透明度` / `液面身份色（GlowColor 规则）` / `环·图标·角标（叠加层）` / `动效与定时器`，并：

1. 在表前加一句总则：**「§2.1 的每一行在 `prototypeVisuals` 两档下各有一份真值；本表记的是签收档（prototype），frozen 档作为逃生门另表。」**（不这样写，C-01 会在下一次对照时重新变成一份"谁对谁错"的争论。）
2. `Sleeping` 行保留 34.72px 与"零动画零定时器"，**新增**：`实测 px≥8 家族 2098/2103/2120（三次独立跑，成像框 40×38 / 46×46 / 44×44 浮动 ⇒ 框列不作尺寸门）`、`timers=no`、`privWS 11.76–11.91MB`。
3. `Armed`/`Muted`/`Settling` 三行的透明度**在 Q-1/Q-5 落定前留 `[待定]`，不留数字**。
4. `Confirming`/`AwaitingApproval`/`Downloading`/`FirstRun`/`Queued` 五行按 **Q-6** 的裁定写"要求 + 残缺标记"，**不写"已实现"**。
5. 把票 62 提案 §7 第 5 条的**"可见比例 0.42"删掉或改为 Q-4 的答案**——那是 §7 目前唯一一处会把已知错误写进契约的句子（C-16）。
6. 新增两维（C-15）与环预算裁切（C-18）各一条脚注。

---

## 6. 本草案能支撑票 62 哪几框、还缺什么（**不动票面 Status 与勾选框**）

| 票 62 AC | 本文件给了什么 | 还缺什么（ ⇒ 不能勾的理由） |
|---|---|---|
| **#1** 浅色+深色桌布都清晰可见 | 浅色侧的数字与出处全在 §0.4/§1 的〔实测〕列（Sleeping 2120/44×44 等） | **深色桌布零差分**（§4 第 3 条）⇒ 框的两半只成立一半 |
| **#2** Sleeping 静态单帧 + 实测零定时器 + CPU≤0.5% | 四次独立 `timers=no` + cpu 0.000/0.000 + privWS，且明写"读活对象不是查策略表"的出处（`BW:445-451` + `main.go:373`） | 数字够；**但真体径 34.72 与 SPEC 的 44 更正仍待 owner 过目（R15#1）**，且 Q-2 命中缺口就在这颗球上 |
| **#3** 每态颜色/速度/描边可辨（逐态截图） | §1 给满 20 行的颜色/动效/描边真值 + 冲突 | **13 个态零实测、零合格截图**（§4 第 1 条）；质感判决本身在票 65（阻在参考图） |
| **#4** 音频驱动液面单调 | 给了单调律的**确定性测试**位置与电平→角速度/胀缩公式（`LQ:224-228`、`REN:543`） | **wav 注入（票 13 的 C8 seam）实测不存在**（§4 第 6 条）⇒ 框不成立 |
| **#5** 靠左/右/上停靠 + 悬停/单击恢复（跨 DPI 与副屏） | dock 几何真值 + 三条 dock 差分行 + 确定性弹回测试（`dock_test.go`，编排者已做单向变异检验） | 已知两处不成立（**A29② 82% vs 42%**、**A29③ 命中 < 可见**）+ **DPI/副屏零证据**（A28）⇒ 框不成立 |
| **#6** 资源复测写进 SLO 且不降阈值 | §0.4/§1 把 SLO §A.2 与四个差分目录对齐，并复述"阈值一动没动" | 归 `docs/SLO.md`（编排者的文件，本文件不改）；另 §4 说明 13 态没有资源行 |
| **#7** **契约草案成文，等签字后才回填** | **本文件即交付物**：20 行、逐行 file:line、三档标记、22 条冲突、15 个待裁问题、10 类找不到数字的格子 | 本框的"成文"半已闭；**"owner 签字"半未发生**，且 §5 明文未回填 ⇒ 编排者可判"草案存在、签字未决" |
| **#8** 对抗验收 + 1:1 裁决表 | §2 的 C-01…C-22 与 §4 的缺口表**可以直接当裁决表的行底稿** | 本文件**不是**对抗验收（作者是本草案的撰写者，非独立验收人）；`62-adversarial-acceptance.md` **仍不存在** |

**本文件刻意不做的事：** 未改 `docs/specs/**`、未改 `docs/PLAN.md`、未改 `.scratch/wisp/issues/**` 的任何 Status/勾选框、未改任何 `.go`、未跑格式化、未动 `docs/reports/*`、未新增 `docs/evidence/s1/68-*`（它引用的推导在 `TT`）、未在 Q-1…Q-15 上替 owner 选任何一个值。
