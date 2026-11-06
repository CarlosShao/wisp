# 65 — 非实现者裁决（玻璃码半格 + [ball].style 配置）

**Status:** code-verdict-complete（**非最终验收**——AC#6 是人判，65 不可由本件结案）
**裁决腿:** verdict65（wave43；非实现者，裁决者≠实现者 D22 双角色）
**被裁决:** `30d8f2e38`（W65-impl 码）＋ `629d6a422`（实现腿交付表）；两 commit 均在 dev
**裁决时 HEAD:** `42a9ae5c9`（A1586 台账笔——**已核**：`git show --stat` 仅 `docs/reports/pending-and-issues.md` +8 行，零产码字节，不影响被裁内容）
**方向依据:** 台账 **A1581**（机主批 D36 `[ball].style` 键，原话在案）＋ A1586 收件记录；本件只追加台账引用，**未写台账**（只追加面归编排者）
**票面:** `.scratch/wisp/issues/65-ball-glass-quality-rework.md`（6 格 AC）
**隔离:** 规则 21 共享克隆 `D:\tmp\wisp\20261106\65verdict\repo`（alternates 指回主树——**主树 gc 前勿 prune**；`build/`、`third_party/sherpa-onnx` 为 junction 复用主树，非副本）
**环境:** Go 1.27.1 windows/amd64 · gofumpt v0.12.0（`D:\work\base\gopath\bin`）· `CGO_ENABLED=1`（默认，未清零——按派单 NEVER CGO_ENABLED=0）· GOTMPDIR/TMPDIR/TEMP/TMP unset · PATH 前插 build+sherpa-onnx
**硬约束遵守:** 零 GUI/窗口、零 exe 运行、零 *winlive* 测试、零桌面像素测量（前腿 74min 挂死教训）；未碰 `docs/reports/pending-and-issues.md`、`frontend/**`、`design/**`、票面、产码

---

## 1. 门禁读数（克隆树内改后现跑，全部 fresh）

| 门 | 命令 | 读数 |
|---|---|---|
| gofmt | `gofmt -l internal/ball` | **空**（rc=0） |
| gofumpt | `gofumpt -l internal/ball internal/config cmd/balldebug` | **空**（rc=0） |
| vet | `go vet ./internal/ball ./cmd/balldebug ./internal/config` | **clean** |
| build | `go build ./cmd/balldebug` | **rc=0** |
| build all | `go build ./...` | **rc=0** |
| test | `go test ./internal/ball -short -count=2` | **ok 0.338s** |
| test | `go test ./cmd/balldebug -count=2` | **ok 0.051s** |
| test | `go test ./internal/config -count=2` | **ok 3.644s** |
| d22scan | `sh scripts/d22scan.sh` | **clean rc=0**：424 产码文件；bans#1-5 internal/=366 cmd/=58；ban#8 internal/=912 cmd/=217；ban#6 frontend/=116（克隆含 tracked frontend/，读数与主树一致） |
| D32 双钉 | `-run "TestWensityRestingFrameIsStatic\|TestStyleSwitchNeverWakesSleeping" -count=2 -v` | 两测各 **PASS ×2** |
| C21 机器 | `-run "TestC21\|TestStyle\|TestWensity\|TestPseudo\|TestNoHardcoded\|TestAnimationPolicy" -v` | 全 PASS：colour rows **verified 126/132**（+6 retired 行钉 ABSENT）· 37 palette 字段 code→table · look colours **45**（5 套）· geometry 表过 · 零消费者判定过 · **CSS 三方 78/78**（tokens.css 侧零新增，无陈旧红）· `TestNoHardcodedColorsInBallPackage` PASS · `TestAnimationPolicyZeroTimerInSleeping` PASS |

`-short` 口径与 68-r2 相同：包内无 `testing.Short()` 分支；桌面/live 由 `//go:build windows && winlive` 排除（克隆内 grep 全包确认，`winlive` 仅 3 个 live 测试文件头），非被 `-short` 跳过。

## 2. 量化规格逐枚静态核实（票 65 硬数字 → 码落点）

| 派单要的数 | 码落点 | 核实方式 |
|---|---|---|
| **带心 0.13R** | `internal/ball/tokens.go:453` `WensityBandCenterFrac = 0.13`；消费于 `wensity_windows.go:43,100,106`；spine 纯函数 `wensity.go:16` `axisY(t)=Center+Tilt·(t−0.5)`，`axisY(0.5)=0.13` 被测试钉 | 常量在 + 消费在 + 测试钉（`TestWensityBandGeometryPinsTheQuantifiedSpec`） |
| **FWHM 0.39R** | `tokens.go:454` `WensityBandWidthFrac = 0.39`；段椭圆 `rx=R·0.39`、`ry=R·0.39·0.62`（`wensity_windows.go:51,97`——0.62 是 FWHM→高斯 σ 软落读法，注释写明） | 常量在；C21 几何表行在；测试钉 token 本值（椭圆纵横比是绘制选择，非规格数） |
| **冷端上翘 0.33R** | `tokens.go:455` `WensityBandTiltFrac = 0.33`；`axisY(1)−axisY(0)=0.33` 被 `style_65_test.go:131` 逐值钉 | 常量在 + 测试钉 |
| **六色 stop 按序** | `tokens.go:435-441` `#C2162F→#E649AF→white→#246BD1→#48D7E6→#D4994B` ＋ `tokens.go:447` `wensityStopTs=[6]float32{0.06,0.22,0.46,0.66,0.84,0.98}`；最近 stop 不插值（`wensity.go:44-60`）；次序由 `style_65_test.go:157-178` 钉死（t 0→1 扫出恰 6 组、逐位比对） | 六色在 + 位置在 + 次序测试钉 |
| **顶弧 0.95R / 0.04R / 57°** | `tokens.go:459-461` 三常量；`wensity_windows.go:120-123` `arcEx(−90±28.5°, width=0.04R)` | 三常量在 + 绘制消费在 |
| **右下亮唇 0.97R** | `tokens.go:462` `WensityRimLipFrac = 0.97`；`wensity_windows.go:117-119` `arcEx(15°..75°)` ＋ `wensity.RimLip #F7FAFF`；其余三向走 `looks[wensity].Rim #1B1D24` 暗缘（renderer step 6）——非对称折射与参考一致 | 常量在 + 非对称结构在 |
| **42% dock** | `tokens.go:508` `DockOverlapFrac = 0.42`（票 64 面值，**本票零触碰**：`git diff 30d8f2e38^ 629d6a422 -- dock*.go hit.go` 空）；drawGlass 经 `DockSquash` 消费（68 已核的半格，码未回退） | 在场 + 零 diff |
| **0.95 settling** | `tokens.go:542` `RestSettledOpacity = 0.95`（票 68 落的值，**本票零触碰**） | 在场 + 零 diff |
| **满幅形变 0.15R** | `tokens.go:456` `WensityBandDeformFrac = 0.15`；`wensityBandDeform` clamp 单调（`wensity.go:30-38`）；`TestWensityBandDeformFullScalePP` 钉 full-scale=token | 在场 + 测试钉 |
| **伪谱 400/40ms·gain 8** | `liquid.go:45-48` `SpecLowTauMs=400 / SpecHighTauMs=40 / SpecHighRateGain=8`；C21 motion 表三行登记 | 在场 + C21 机器过 |
| **config `[ball].style`** | `schema.go` `Style string \`toml:"style" default:"glass"\`` ＋ 词表 `ballStyles`；`validate.go::validateBall` 枚举拒绝带全词表；漂移钉 `TestBallStyleVocabularyMatchesBall`（config↔`ball.StyleNames()` 双向）；`PLAN.md:2741` 一行 diff（`size` **`style`** `position…`）＝A1581 批的原文 | 在场 + 双向测试钉 |
| **`balldebug -style`** | `main.go:97,130-135`（未知值 exit 2 带词表）＋ `diff_windows.go:54,418-420` 子进程透传 | 在场（真机对拍两风格的入口已备） |

**结论：派单点名的每一个量化数都以 token 形态在码里、进 C21 表、被机器检查；渲染新文件零颜色字面量。**

## 3. D32 纪律核实

- **`anim.go` 零 diff**：`git diff 30d8f2e38^ 30d8f2e38 -- internal/ball/anim.go` = **0 行**；`AnimationPolicy` 权威表一字未动。
- **零新定时器**：全部 7 个新增/改动 ball 文件 grep `time.NewTimer|time.After|time.Tick|time.NewTicker` → **零命中**；伪谱三带在 `pushLevel` 同步派生（音频样本即时钟，无新时钟）；`liquid_windows.go` 的 `syncLiquidTimer` 判据未动（diff 仅 1 行 `applyTo` 签名适配）。
- **帧内零 COM 分配**：16 段带 brush ＋ 白核 ＋ 底辉光全部 `setAssets` 一次烘焙（`renderer_windows.go:252-255→setWensity`）、`dropAssets` 随旧资源释放（:328-333）；帧内只有 `setSolid`（=SetColor 复用共享笔刷）、`fillEllipse`（栈槽 D2D1_ELLIPSE）、`arcEx`（DrawLine 折线，几何在寄存器）。
- **静帧确定性**：谱带只能经 `applyTo(v, bandDeformState(curState))` 盖章（`liquid_windows.go:142`），白名单仅 Listening/Thinking/Speaking（`liquid.go:371-383`）；其余 17 态（含吃音频流的 Conversation/Warm/Acting）恒零；`endSession` 归零三带；`VisualFor` 20 态全零谱带被 `TestWensityRestingFrameIsStatic` 钉。
- **`liquid_test.go` 行为断言未动**：仅 3 处调用点机械补 `, false` 实参，断言体逐行对过零改动。
- 新钉 `TestStyleSwitchNeverWakesSleeping`（SetStyle 全态 policy 零位移 + Sleeping 双白名单外）与 `TestWensityRestingFrameIsStatic` 均 count=2 绿。

## 4. 判 AC 1:1（票面 6 格）

| # | 票面 AC | 判 | 依据 |
|---|---|---|---|
| 1 | 参考图入库＋量化分析件逐条注明来源/采样位置 | **PASS（替代件判定）** | `65-style-direction.md` §1＝逐像素量化（§5 四把采样脚本 v2–v5 在盘可复跑，`.scratch/wisp/probes/65/`）；三张参考 PNG 在盘 `probes/65/refs/`（jarvis-orb / wensity-siri-orb / youtube-audio-visualizer）。AC 字面文件名 `65-glass-reference-analysis.md` **不存在**——设计腿 §7④ 自报「本件为上位替代（分析+方案+映射+设置面四合一）」；本裁决认**实质**（逐条量化+出处）满足，文件名错位是编排者对票面文案的簿记项，非量化缺口。参考图落 `.scratch` 临时件（未入库 git）而非票面建议的 `design/refs/`——`design/**` 归外部 agent，临时件只建不删是仓规，不判缺。 |
| 2 | 属性映射 C21 tokens；渲染零新增硬编码色值（d22scan＋人工抽查双过） | **PASS** | 全部色值以 `hex()/rgba()` 落 `tokens.go` token 结构（§2 表）；`wensity.go`/`wensity_windows.go`/style 测试零字面量（`TestNoHardcodedColorsInBallPackage` 绿）；d22scan 克隆树 clean；本裁决逐文件人工复读确认新代码只引 token。 |
| 3 | 并排差分产物在档：当前实现/返工后/参考图（同区域同底色） | **未做——真桌面件（非代码缺口）** | 65 名下零渲染产物（find 全 probes/65 仅 refs 三图）。该产物必须跑 exe/开窗，本腿明令禁跑（74min 挂死前科）。机器侧前置已备齐：`-diff`/`-diff-gate`/`-style` 透传/`-tour`。落点＝机主桌面趟（§6 第 6 步）。 |
| 4 | 票 62 的 11 个 dock 确定性测试与全部球测试仍绿，行为断言未动 | **PASS** | `dock_test.go` 现存 **12** 枚 Test func（票 62 计 11，现 12——只增未删），`dock*.go`/`hit.go`/`anim.go` 对 65 两 commit **零 diff**；`liquid_test.go` 仅机械实参；全包 `-short -count=2` ok 0.338s。 |
| 5 | Sleeping 差分重跑：timers=no、cpu≤0.5%、WS≤25MB、px_delta_ge8≥2120 | **未跑——真桌面件（本腿明令禁跑）** | 静态侧全证：零新定时器（§3 grep）、带/弧/唇只增可见像素而门是**下限门**（`diffgate_windows.go:36` floor 语义，只许升）⇒ 门不会因 65 变红；但四个数字必须真机量，本裁决不代跑、不估读。 |
| 6 | **owner 第二次看 `-tour` 并明确说「这次对了」** | **人判——本裁决不可代行（派单明示）** | 票 65 的完成判据是一个人的桌面观感判决，测试绿不能替；此前 SPEC-08 §2 的 INTERIM 标注不得摘。本裁决能做的是把可预置的入口全部备好（§6 清单）。 |

**半格结论：AC#1/#2/#4 PASS；AC#3/#5＝真桌面未跑（代码侧前置全备、非代码缺口）；AC#6＝人判未发生。⇒ 65 保持 open，无权也不应结案。**

## 5. 零写面自查

- 本腿写面＝仅本件 `docs/evidence/s1/65-verdict.md`（显式 pathspec commit）；票面、产码、台账、specs、frontend、design 零触碰。
- 零 GUI/窗口进程、零 exe 运行、零 winlive、零桌面测量；无 push。
- 可核：`git show --stat HEAD`（本 commit 仅本件）；`git diff 30d8f2e38^ 629d6a422 --stat -- internal/ cmd/`（=实现腿自己的 17 文件，本裁决零追加）。

## 6. 机主真桌面待办（唯一剩余；编排者照此派/递）

前置：桌面清空、缩放 100%、`tasklist` 确认无 wisp/balldebug 残留。
1. `go build -o build/balldebug.exe ./cmd/balldebug`
2. `build\balldebug.exe -diff -diff-gate`（Sleeping 行核 `timers=no`、`cpu≤0.5%`、私有 WS≤25MB、`px_delta_ge8≥2120`——AC#5）
3. `build\balldebug.exe -diff -diff-gate -style jarvis`（同门对拍 jarvis 收值落玻璃帧）
4. **AC#3 三张并排**：同一区域同底色 截 现状（glass 默认 wensity）/ 旧 look 对照（`-look aurora`）/ 参考图 `probes/65/refs/wensity-siri-orb.png`
5. **AC#6 人判**：`build\balldebug.exe -tour -style glass`——owner 亲眼看，说「这次对了」（或哪里不对，改 token 不改码）
6. 通过后编排者方可摘 SPEC-08 §2 INTERIM、结 65；jarvis 渲染腿（wave43 W65-jarvis）另算
