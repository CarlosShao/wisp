# 69 — C21 token 表的 61 行新内容没有任何机器检查（A24-D4）

**Status:** review（AC#1/#3/#4 已交付；AC#2 的三次变异见
`docs/evidence/s1/69-mutation-tokens.md`。**新登记的 6 条分歧 + 29 条零消费者等编排者判**，
本票没有自行改表、没有自行改任何默认值）
**Claimed by:** agent-ticket69
**Last update:** 2026-09-21
**Blocked by:** 68（它正在改 `internal/ball/`，同包并行=假并行）⇒ **本票在票 68 收尾前不得开工**
**Parallel slots:** ≤1 sub-agent（`internal/ball/` + `docs/evidence/s1/c21-native-tokens.md`）
**Spec refs:** C21、SPEC-08 §2（**冻结，不得编辑**）、D29
**登记项:** A24-D4

## 为什么要这张票

票 12 AC#7 的代理把 `docs/evidence/s1/c21-native-tokens.md` 从 **133 行补到 227 行**
（新增 **25 条几何/动效常量 + 36 个 look 色**），并把 78 条配色行与
`design/assets/tokens.css`、`internal/ball/tokens.go` **逐值对账到零漂移**。这是真交付。

**但它同时暴露了一个空档**：今天只有 **20 条**配色受
`TestTokenGoldenValues` 保护，`TestNoHardcodedColorsInBallPackage` 管的是"不许在绘制码里写字面量"，
**不是**"表与码一致"。⇒ **新加的 61 行没有任何机器检查**：
改了 `tokens.go` 不改表、或改了表不改码，都不会被任何东西抓到。
表头现在也写明了"36 个 look 色在 `tokens.css` 里无出处、只有人工表这一道保证，
**不得当契约级事实引用**"——这句话今天是**真的**，本票的任务就是让它变成不需要这句话的东西。

## 验收标准

- [x] **AC#1 反向检查落地**：一条机器断言，遍历本票负责的两类 token
  （几何/动效常量、look 色），证明**表里每一行都能在码里找到同名同值的对应物，且反之亦然**。
  实现方式自选（testdata 里的表 + 解析，或从 Go 侧生成再 diff），但**判据必须是双向**：
  只查"码里的都在表里"是不够的（那正是漂移的方向）。
  **交付**：`internal/ball/tokens_table_test.go`（新建，未改 `tokens.go` 一个字节）
  解析 markdown 表本体（**没为好测改表结构**），两条用例：
  `TestC21TableColourRowsMatchCode`（40 暗 + 38 亮 + 36 look = **114 配色行**，按 `tokens.go`
  自己的 `hex()/rgba()` 求值后逐通道相等）与
  `TestC21GeometryRowsMatchCodeConstants`（**25 行 / 56 个导出常量**，值必须是该行真写出的数字，
  单位按名字承诺：`*Ms` 收 `1.6s` 这类秒写法、`*Px` 要 px、`*FPS` 要 fps、比率要裸数）。
  反向：`Palette` 40 个 Color 字段、`looks` 4×9 个 look 色、`tokens.go` 56 个导出常量
  各要求**恰好一行**声明；亮表未列的字段额外断言"真的沿用 `:root`"（级联：`BallHalo`、`OnSolid`）。
- [x] **AC#2 变异检验**：故意把 `tokens.go` 里一个几何常量改一个数（或删一行 look 色），
  新用例**必须转红**；然后还原并证明全套件回到绿。⚠ **先 grep 证明变异真的落地**再跑测试，
  跑完立刻 `git checkout --` 还原；原始输出留在 `docs/evidence/s1/69-mutation-*.md`。
  **交付**：三次（码侧 1 + 表侧 2），每次先 `grep -c MUTATION`=1 + `git diff --stat` 证落盘，
  还原后 `grep -c MUTATION`=0 且 `git diff` 空。原始输出与"全套件只有那一条红、
  既有 47 条全绿"的对照在 `docs/evidence/s1/69-mutation-tokens.md`。
- [x] **AC#3 不许弱化**：`TestTokenGoldenValues` 现有 20 条断言**一条都不许删或放宽**；
  只能是超集。若某条 look 色因票 65（玻璃质感返工）注定要变，**登记成显式豁免行**而不是删断言。
  **交付**：`internal/ball/tokens_test.go` **一个字节未动**（`git status` 可证）；新用例是超集
  （20 条金标准里的每一条都被 114 行的表↔码断言覆盖）。今天**没有**需要显式豁免的 look 色，
  显式豁免清单在码侧（`c21OutTableExempt`，15 条带理由）——见下「新登记」。
- [x] **AC#4 门禁**：`gofmt -l` 触及包为空、`go vet ./internal/ball/`、
  `go test -count=2 ./internal/ball/`（**非 winlive**；若本票不需要桌面就别碰桌面）。贴原始输出。
  **交付**：本票**不需要桌面**，没碰 winlive。`gofmt -l internal/ball` 空、`go vet` rc=0、
  `go test -count=2 ./internal/ball/` → `ok 0.555s`；`-v` 计数 `PASS=48 FAIL=0`，
  其中本票新增 **4 条**用例（`--- PASS` 四条，命令与输出照抄在证据文件「基线」段）。

## 边界（不要越界）

- **不碰** `docs/specs/*`（SPEC-08 冻结）。若发现表与契约冲突，**报出来**由我裁定，不改文本。
- **不在本票改球面默认值**（那是票 68 AC#2 的活）；两票若都动 `internal/ball` 必须**串行**。
- 若票 65 的质感返工尚未开始，本票**只锁现有值**，不预判新值。

## 交付物与覆盖面（数字，别说"全部"）

新建 `internal/ball/tokens_table_test.go`（**4 条用例**），`internal/ball/tokens.go` 与
`docs/evidence/s1/c21-native-tokens.md` **均未改动一个字节**（三次变异的还原证据在证据文件里）。

| 判据 | 今天真被断言的量 |
|---|---|
| 表 → 码：配色行逐通道相等 | **114 行**（40 暗 + 38 亮 + 36 look） |
| 表 → 码：几何/动效行的数值 | **25 行 / 55 个有数值的常量**（51 个带单位标记 + 4 个只能裸匹配，测试日志显式列出） |
| 码 → 表：`Palette` Color 字段 | **40 / 40** |
| 码 → 表：`looks` 色 | **36 / 36**（4 套 × 9 字段） |
| 码 → 表：`tokens.go` 导出常量 | **56 / 56**（含 `CountdownFontPx`） |
| 级联断言（亮表未列必须等于 `:root`） | **2**（`BallHalo`、`OnSolid`） |
| 表外码侧常量的显式豁免 | **15 条带理由**（`c21OutTableExempt`） |

⇒ 票面说的 61 行：**36 个 look 色全部**进了硬断言（值相等），**25 个几何/动效常量全部**进了
双向名字检查 + 数值检查。既有 20 条金标准一条没动（AC#3）。

## 零消费者清单（报告，不判红 —— AC#3 的理由：这是 owner 未批的规格）

`TestC21TokenConsumerReport` 每次跑都打印，**29 / 105** 条没有任何非测试引用
（18 个绘制文件被 AST 扫过；AST 不吃注释，所以"只在注释里出现"不算消费者）：

- **`tokens.go` 常量 10 条**：`BallSizeSmallPx`、`BallSizeLargePx`、`DurFastMs`、`DurBaseMs`、
  `DurSlowMs`、`FirstRunGuidePulseMs`、`FontSizeMonoPx`、**`CountdownFontPx`**（A33② 点名的那条）、
  `SwimLevelGain`、`SpinLevelGain`。
- **`Palette` 字段 18 条**：`GlassRing`、`GlassHi`、`GlassHiSoft`、`FgPrimary`、`FgTertiary`、
  `InfoLine`、`OrbShadow`、`TintAccentHi`、`TintSuccessHi`、`TintSuccessLo`、`TintWarmHi`、
  `TintWarmLo`、`TintWarnHi`、`TintWarnLo`、`TintDangerHi`、`TintDangerLo`、`TintNeutralHi`、
  `TintNeutralLo`（12 条 tint 里 11 条零消费者，只有 `TintNeutralMid` 被 `statevisual.go:215` 读）。
- **look 字段 1 条**：`looks[*].Deep`（四套 look 的"背面折射色"没有任何绘制路径读它）。

## 新登记给编排者判的分歧（本票一条都没自行"改到迁就另一边"）

- **D7 `internal/ball/hit.go` 两条几何常量不在表里**：`RingMarginPx = 8.0`、`ClickTolerancePx = 4.0`。
  窗口/命中半径都从它们推（`DockPos` 的 ring margin 8 就在票 68 AC#1 的算式里）⇒ 是"码里用了、表里没有"。
  今天它们进 `c21OutTableExempt` 并写明理由，**不等于裁定**：要么进表，要么书面判域外。
- **D8 `internal/ball/liquid.go` 九条动效常量不在表里**：`SpinBaseRadPerS`、`SpinLevelRadPerS`、
  `SpinSummonRadPerS`、`SilenceLevelGate`、`SilenceHoldMs`、`LevelAttackTauMs`、`LevelReleaseTauMs`、
  `MotionEpsilon`、`FrameIntervalMs`。表自己的范围句写"只覆盖 `tokens.go`"（A24-D6），
  所以本票按范围只报告不判红 —— 但票 62 的液体动效规格事实上分成了两处真相源。
- **D9（最该看的一条）「音频包络 → 液体运动」那两格讲的是代码没在做的事**：表里 `SwimLevelGain = 0.55`、
  `SpinLevelGain = 1.0` 零消费者，而真正驱动旋转的是 D8 里的 `SpinBaseRadPerS 0.45` / `SpinLevelRadPerS 4.20`。
  ⇒ 这不是数值漂移而是**语义漂移**：表在描述一件没有实现的映射。裁定要往"改表 / 补实现 / 删规格"哪个方向走，
  与 A33①（`Speaking` 呼吸到底由什么驱动）是同一族问题。
- **D10 `Palette` 的 18 个零消费者字段**：票 62 之后渲染走 `looks`，`Palette` 的 tint/玻璃边缘一半
  是冻结 SPEC-08 §2.1 时代的规格。要不要它们 = A33② 的 Q 列表（owner 判），本票只让它们**可见**。
- **D11 `looks[*].Deep` 有值无绘制路径**（四套 look 各一份，共 4 个色值今天纯装饰）。
- **D12 时长档 `DurFastMs/DurBaseMs/DurSlowMs` 零消费者**，而票 62/64 自建的 `BorderOpenMs 220`、
  `BorderCloseMs 180`、`SettlingFadeMs 260` 才是真被用的数 —— 表里 `--dur-base 180ms` 那行
  与 `BorderCloseMs = 180` 恰好同值，**同值不等于同源**：改 `--dur-base` 不会有任何测试变红。

## 本票没能证明的部分（如实）

1. **表 ↔ `tokens.css` 仍无人机器检查**（表头那句"权威真相源"依旧只靠票 12 的人工对账）。
2. **几何行是"值出现在行内数字里"，不是双射**：同行撞号可互相掩盖（例：`BallSizeMaxPx` 72→44 不会红，
   但 72 会掉进"无人认领数字"日志 —— 那是报告不是门）。
3. **`FontFamily` 的值无断言**（表只引用 `--font-sans`，不写字符串值）；名字双向仍检查。
4. 零消费者按**字段名**匹配 selector，同名字段会让报告偏乐观（今天实测无冲突：`Visual` 用的是
   `GlowColor`/`RingColor` 这类全名）。
5. **"哪一套 look 真的在被渲染"锁不住**（运行时 `SetLook` 决定），本票锁的是四套的全部 36 个值。

## Progress log（append-only；每个 commit 一行 `- [ISO-UTC] agent=... did=...`）

- [2026-09-21T01:26:00Z] agent=ticket69 did=新建 `internal/ball/tokens_table_test.go`（4 条用例，
  114 配色行 + 25 几何行/56 常量的**双向**断言 + 级联断言 + 零消费者报告 + 表外常量的显式豁免清单）；
  `tokens.go`、`c21-native-tokens.md`、`tokens_test.go` 全部一字节未动。勾 AC#1/#3/#4。
  门禁：`gofmt -l internal/ball` 空、`go vet ./internal/ball/` rc=0、
  `go test -count=2 ./internal/ball/` → `ok 0.555s`，`-v` 下 `PASS=48 FAIL=0`（新增 4 条）。
  新登记 D7–D12 给编排者判（表外码侧常量、`SwimLevelGain`/`SpinLevelGain` 的语义漂移、
  18 个零消费者 `Palette` 字段、`looks[*].Deep`、`Dur*Ms` 与 `BorderCloseMs` 的同值不同源）。
  next=AC#2 的三次变异（码侧 `DockTriggerPx` 16→17、表侧 look 色一位、表侧几何行 16px→17px），
  原始输出落 `docs/evidence/s1/69-mutation-tokens.md`。
