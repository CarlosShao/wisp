# 65 — Jarvis 粒子风格渲染器实现腿（W65-jarvis-impl，第 3 次尝试）

**Status:** delivered（代码 + 测试 + 接线全落，门禁 PASS，d22scan 一条预存红非本腿）
**Leg:** W65-jarvis-impl（wave47，jarvis 第 3 次派单；前两次零文件未交付，本次以「先提交、宁瘦勿空」为铁律执行）
**HEAD 锚点:** 派单时 `4656ec715`，实际工作基线 `9dd310cde`（A1590 派单提交，编排者已前移）
**Owner 方向:** A1581（风格可切换）＋ A1575/A1576；本腿只做 `65-style-direction.md` §1.2/§2 的 Style B「jarvis」
**票面:** `.scratch/wisp/issues/65-ball-glass-quality-rework.md`
**Spec refs:** SPEC-08 §2/§2.1 · D32（静帧纪律）· C21（禁色机器）· D36（style 键已由前一腿落 config 面本腿零触碰）

---

## 1. 本腿交付了什么（可核）

| 文件 | 内容 |
|---|---|
| `internal/ball/jarvis.go`（新） | 纯映射层，仿 `wensity.go` 先例，无 Win32/D2D，任何 GOOS 可测：固定种子 xorshift64* PRNG、`jarvisLayout`（确定性粒子场）、`jarvisParticleCount`（44px→75，带内 60–120）、`jarvisMainRingFrac`（态结构：thinking 0.67R / alert 0.39R / 其余 0.50R）、`jarvisRingSpin`（仅 Thinking 推进相位；静止恒 0）、`jarvisListenPulse`/`jarvisScatterAmp`（电平驱动，静止恒等）、cyan/violet 家族选择器 |
| `internal/ball/jarvis_windows.go`（新） | 纯 D2D 粒子球渲染器（镜像 `wensity_windows.go` 结构）：`setJarvis`（state bake：基底辉光 brush + 核心光斑 brush + 固定种子粒子布局）、`drawJarvis`（基底辉光 → ≤120 个小椭圆粒子点（2-3 条椭圆轨道 inner 0.16R / 主环按态 / core 0.07R + 稀疏散点，倾斜率 0.62）→ 核心亮点）、`dropJarvis`。逐帧只 fill 已有 COM 对象（D32 CPU 纪律），无每帧 RNG，无新定时器 |
| `internal/ball/jarvis_windows_test.go`（新） | 6 个纯测试：令牌行齐全（含 cyan/violet 双家族色相验证 + 家族选择器词表）、逐态环半径纯函数（经 `VisualFor` 钉 0.67R/0.39R/0.50R）、44px 粒子数带 60–120、布局确定性 + 有界性、**TestJarvisRestingFrameIsStatic**（镜像 `TestWensityRestingFrameIsStatic`：静止态零频谱/零散布/零脉动/零自旋/怠速环径，全态 `VisualFor` 不设频谱）、动效只出现在活跃态 |
| `internal/ball/tokens.go`（改） | 追加 `JarvisPalette` 令牌集 + `jarvisTokens` 值 + `JarvisTokens()` 访问器（§2）。tokens.go 是包内唯一允许颜色字面量的文件（`TestNoHardcodedColorsInBallPackage` 机器规定） |
| `internal/ball/renderer_windows.go`（改，4 处最小接线） | ① renderer 结构体加 `jarvis jarvisAssets` 字段；② `setAssets` 在 `activeStyle == StyleJarvis` 时 bake；③ `dropAssets` 调 `r.dropJarvis()`；④ `drawFrame` 的 glass 分支按 `activeStyle` 分派：jarvis → `drawJarvis`，glass → `drawGlass`（原路径零改动） |
| `docs/evidence/s1/c21-native-tokens.md`（改，删 3 行） | jarvis 让 `Palette.TintSuccessHi/TintDangerHi/TintDangerLo` 三个令牌重新接线（alert/success 语义色派生），测试判定其「未接线豁免」行陈旧并明确指示 `delete the row`——按指示删除该 3 行（§4） |

## 2. 令牌设计（C21 合规）

- **零硬编码色字面量**：渲染代码只引用 `JarvisPalette`（tokens.go 里的实测值）与 C21 `Palette`。禁色扫描（`TestNoHardcodedColorsInBallPackage`）对非 tokens.go 文件连**注释**都不豁免——首版 jarvis.go 注释里的 `#9F63FF` 被抓，已改为 `0x9F63FF` 记法。
- **cyan 家族**（hue≈187°实测）：Dim `0x1B7C86` / Bright `0x36E7F5` / Core `0xC6FFFF`；thinking 覆盖行 Bright `0x5EFAFF` / Core `0xFFFFFF`；tile 基底辉光 `0x050913`。
- **violet 变体**（hue≈265°实测）：Dim `0x562D98` / Bright `0x9F63FF` / Core `0xFFFFFF`。v1 经 `SetJarvisFamily("cyan"|"violet")` 选择器可达（调试/测试缝）；`[ball] jarvis_palette_hue` 自定义色相按 §7 ⑥ 留 v2，本腿不落 config。
- **alert/success 语义色按 §3.2 从 C21 派生**：alert 粒子取 `Palette.TintDangerLo/Hi`、success 核心取 `Palette.TintSuccessHi`（ember tile 采样按 §1.2 诚实备注不取用）——这正是三条豁免行变陈旧的原因。

## 3. D32 静帧纪律（本腿的硬约束执行）

- Sleeping/静止 = **一帧静态渲染**：非活跃态 `AnimationPolicy` 原封不动（`TestStyleSwitchNeverWakesSleeping` 仍绿），粒子布局 bake 一次、静止帧就是 bake 的恒等几何。
- **零定时器**：jarvis 不新增任何 timer/goroutine/COM 常驻对象；活跃态帧由既有 `AnimFrame` 计时器供给，`phase` 参数进 `drawJarvis`。
- **静止不自旋**：`jarvisRingSpin` 只在 `v.Sweep`（Thinking）非零；静止帧恒 0（测试钉死）。 Listening 脉动/散布严格电平比例（电平 0 = 恒等），避免「静默活跃帧」漂移。
- 布局确定性：固定种子 xorshift64*，同 (state,size) 两次 bake 逐位相等（测试钉死）；无每帧 RNG。

## 4. 结构读法与偏差（诚实声明）

1. **态结构从 Visual 标记读出，不加 State 字段**：包法律是「renderer 是 Visual 的纯函数」（statevisual.go:26 注释），本腿不碰 Visual/statevisual.go——Thinking=`v.Sweep`、alert=`v.RingPulse||v.Icon==IconX`、其余（含 success/Speaking）=idle 0.50R。逐态映射测试经 `VisualFor` 走真映射。
2. **jarvis 不画玻璃壳**：§2.2 规定 jarvis 结构 = 基底辉光 + 3 条粒子环 + 核心光斑；粒子可读性由 `JarvisPalette.BgGlow` 暗基底承担（参考站 tile 底色实测）。dock 挤压/状态环/角标/图标等共享层照旧（在 `drawFrame` 的 jarvis 分支之后统一画）。
3. **c21-native-tokens.md 删 3 行豁免**：这是测试指示的动作（token 接线后豁免必须撤销），不是缩契约——三个令牌的表行本来就在，只是消费者从无到有。
4. **预存红（非本腿）**：`sh scripts/d22scan.sh` 在 `internal/agent/pathmask.go:38` 报 phantom-citation（引用不存在的 `tools/task.go`）。该文件是**共享树里另一在飞腿的未跟踪文件**（`git status --porcelain -- internal/agent/pathmask.go` 输出 `??`；`git show HEAD:internal/agent/pathmask.go` fatal——不在 HEAD）。本腿写面零输出于 internal/agent，finding 数 = 1 且全部落在这个外来文件；本腿全部文件过扫零 finding。
5. `cmd/balldebug/main.go:97` 的 `-style` 帮助文案仍写着「jarvis draws the glass frame until its renderer lands」——现渲染器已落地，该文案过时但属共享编译面最小写面之外，留编排者收口（一行改动，无行为影响）。

## 5. 门禁读数（本腿实跑）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `gofumpt -w <5 文件>` | 已格式化，无 diff 残留 |
| gofmt | `gofmt -l internal/ball` | **空** |
| go vet | `go vet ./internal/ball ./cmd/balldebug` | **PASS** |
| go build | `go build ./cmd/balldebug` | **PASS** |
| go test | `go test ./internal/ball -short -count=2` | **ok 3.927s**（含 6 个新 jarvis 测试 ×2） |
| d22scan | `sh scripts/d22scan.sh` | 1 finding，全落 `internal/agent/pathmask.go:38`（另一腿未跟踪文件，预存/外来，非本腿写面）；本腿文件零 finding |

环境：Git Bash；`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；`unset GOTMPDIR TMPDIR TEMP TMP`；`CGO_ENABLED=1`（未用 0）；未开 GUI、未跑 wisp.exe、未跑任何 *winlive* 测试。

## 6. 写面核实

- git 提交：显式 pathspec 枚举上述 7 个文件，无 `add -A`、无 push、无 amend/reset。
- 未触碰：`frontend/**`、`design/**`、`docs/reports/pending-and-issues.md`、`internal/agent/**`、`internal/config/**`、`cmd/wisp/**`（共享树里 `cmd/wisp/run.go` 的在飞改动属于别的腿，本腿零接触）。
- 可核命令：`git show --stat <本腿hash>`（只含上表 7 文件）；`git status --porcelain -- internal/ball docs/evidence/s1`（提交后只剩其他腿的未跟踪件）。

## 7. 残留（交编排者）

| # | 事项 | 建议 | 不做的代价 |
|---|---|---|---|
| ① | owner 桌面验收：`balldebug -style jarvis -diff -diff-gate` 差分门重跑（timers=no / px_delta_ge8≥2120 / CPU≤0.5% / WS≤25MB） | 编排者按 §2.5 AC#5 跑桌面差分 | 测试绿≠像素验收，票 65 完成判据是人 |
| ② | `balldebug -style` 帮助文案一行更新（§4.5） | 随下一枚 cmd/balldebug 写腿顺带 | 文案误导一次，无行为影响 |
| ③ | `[ball] jarvis_palette_hue` 自定义色相 | v2（§7 ⑥ 已裁决） | 晚一步拿到自定义色，不影响风格切换 |
| ④ | `internal/agent/pathmask.go` phantom-citation（外来在飞腿的文件） | 归该腿或票 212 收口腿处理 | d22scan 维持一条红直到该腿收口 |
