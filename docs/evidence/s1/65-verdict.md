# 65 — 非实现者裁决（玻璃码半格 + [ball].style 配置）

> ★**2026-11-07 追加：A1594 改道后的第二轮裁决见 §7（orb 集成三腿 a2e535674/d687dd910/8fb821bba）。
> §1–§6 是 wave43（D2D wensity/jarvis 路线）的历史裁决，原样保留。**

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

---

## 7. 【2026-11-07 追加】第二轮裁决 —— A1594 改道（集成 LerSent001/orb）后 orb 三腿

**裁决腿:** verdict（A1594 改道后第 65 票第二轮；非实现者，裁决者≠实现者）
**被裁决:** `a2e535674`（orb 集成腿）＋ `d687dd910`（出屏修复腿）＋ `8fb821bba`（黑窗修复腿）
**裁决时 HEAD:** `f5bbc97f7`（已核：`df331ec09..f5bbc97f7` 期间 3 笔 commit 全部只碰 composer/145/190 证据面，`git diff --stat df331ec09..HEAD -- internal/ball/` = 空，被裁内容零影响）
**环境:** Go 1.27.1 windows/amd64 · gofumpt `D:\work\base\gopath\bin\gofumpt.exe` · GOTMPDIR/TMPDIR/TEMP/TMP unset · PATH 前插 build+sherpa-onnx · 双突变在做完即还（`git status --porcelain internal/ball/` 复归空读数）
**隔离:** 探针/基线全在仓外 clone（`/d/tmp/wisp/20261107/verdict-pre-orb` @`a2e535674^`、`/d/tmp/wisp/20261107/verdict-head` @HEAD）；共享工作树只发生过两次已还原的突变，其余零写。

### 7.1 实现腿声明逐条独立复跑

| # | 声明 | 独立复跑读数 | 判 |
|---|---|---|---|
| 1 | 资产=上游件逐字节 | `sha256sum internal/ball/assets/orb-siri.html` = `d93da649d11bc4bc3ab91c006777bbadce1c2eb7d0ffc9d91608ad8d866e3bc6`，与钉死值一致；`git log --all -- <asset>` 仅 a2e535674 一笔；`git diff a2e535674..HEAD -- <asset>` 空；89660 字节 | **PASS** |
| 2 | 像素回归真有牙 | 基线亲跑 `go test -tags winlive ./internal/ball/ -run TestBallLiveOrbPageRendersPixels -count=1 -v` → **PASS 625/1764 = 35.43% 非黑**（实现腿自报 594/1764=33.67%，同量级）；★突变①注释 `browser.Show()` → **FAIL：best non-dark ratio 0.000% / 1764 采样**，25s 内全黑；还原后 diff 空 | **PASS（有牙）** |
| 3 | 出屏回归真有牙 | 基线亲跑 TestBallLiveOrbFollowStaysInsideWorkArea → **PASS，读数 work={0,0,3440,1392} ball=(3340,1292)-(3440,1392) orb=(3196,1260)-(3328,1392)**（全内，翻左+垂直夹紧）；★突变②`orbFollowXY` 首行改回无条件 `return ball.r + gap, ball.t` → **FAIL：orb=(3452,1292)-(3584,1424) 双轴出屏**（水平右缘 3452>3440，垂直底缘 1424>1392——精确复现 A1596 的 bug 形状）；还原后 diff 空 | **PASS（有牙）** |
| 4 | AC#4 行为断言未改 | 三笔 commit 触碰文件清单均无 `dock_test.go`/`live_windows_test.go`；`git diff a2e535674^..HEAD -- internal/ball/dock_test.go internal/ball/live_windows_test.go` = **0 字节**；dock_test.go 现存 12 枚 Test func（62 计 11，只增未删）；11 枚 dock 测试在 winlive 全套里逐枚亲见 PASS | **PASS** |
| 5 | 实现腿门禁自报 | gofumpt 空（亲跑含 cmd/wisp+cmd/balldebug）/ vet rc0 / build rc0 / 默认套件 ok 0.482s / d22scan clean（939+224 Go files，亲跑）——全部复跑一致 | **PASS** |

### 7.2 门禁读数表（本裁决亲跑）

| 门 | 读数 |
|---|---|
| `gofumpt.exe -l internal/ball/ cmd/wisp/ cmd/balldebug/` | 空（rc0） |
| `go vet ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1 -timeout 20m` | ok 0.482s |
| `go test -tags winlive ./internal/ball/ -count=1 -timeout 20m` | **4 FAIL / 5 SKIP / 其余 PASS**（4.3 逐枚归因） |
| `sh scripts/d22scan.sh` | clean（internal/=939、cmd/=224 Go files，零 violation） |

### 7.3 winlive 全套 4 FAIL 逐枚归因（本裁决的核心发现）

基线对照用仓外 clone `a2e535674^`（pre-orb）＋同机同时段复跑：

| 测试 | HEAD | pre-orb 基线 | 归因 |
|---|---|---|---|
| TestBallLiveEdgeDock | FAIL（`edge 1: docked orb boundary -4, want the work edge 0`） | **FAIL（同一读数 -4 vs 0）** | **存量环境失败**（dock 判据 vs 本机任务栏/缩放几何，与 orb 三腿无关） |
| TestLiveConfirmingCancelAndEscReturned | FAIL（Esc 未取消 Confirming） | **FAIL（同形）** | **存量环境失败**（Esc 全局借用受桌面环境干扰，与 orb 无关） |
| TestLiveNeverStealsFocus | FAIL（WFP=第三方窗） | **PASS** | **orb 集成引入的真回归**（见下） |
| TestLiveTransparentCornerFallsThrough | FAIL（WFP=第三方窗） | **PASS** | **orb 集成引入的真回归**（同根，见下） |
| （SKIP 5 枚） | SKIP-LOUD（热键被占/指针被拒的明文跳过） | 同 | 环境条件跳过，非失败 |

**真回归的机制（探针实证，非推断）**：a2e535674 在 `ball_windows.go:419`（renderFrame）与 `:554`（present）加了 `if b.orbLive() { return }` 守卫——orb 挂上后 D2D 帧停画，**但没有任何代码再往分层球窗 ULW 一帧**。分层窗从进程启动起从未提交过内容 ⇒ 每像素 alpha 恒 0 ⇒ OS 把整窗当全透明：看不见、WindowFromPoint 直接穿透。探针读数（同一停靠位 (3316,1268)-(3416,1368)，同机同分钟）：

- pre-orb clone：`WFP(球心) = WispBallWindow（isBall=true）`
- HEAD clone：`WFP(球心) = "Hermes"(pid 30436) / Tabbit(pid 13996)`——球被 OS 跳过（两枚第三方窗只是恰好压在穿透落点上，**不是**它们挡了球：它们均非 topmost，且 pre-orb 同屏同位下 WFP 照样返回球）

**后果评级**：球不可见+不可点＝产品级视觉回归。像素用例（TestBallLiveOrbPageRendersPixels）之所以仍 PASS，是因为它量的是 **orb 窗**自己的合成像素（WebView2 内容在画），量不到**球窗**已透明这件事——两扇窗，两回事。此发现与三枚实现腿的自证证据不矛盾（它们各自声明范围内属实），是**集成面缺口**：a2e535674 的证据件"已知边界"只写了"D2D 对象仍创建是无害浪费"，没写"球窗从此不画"。

**裁决意见**：orb 窗与球窗的可见/点击关系需要一条收尾腿定形状（方向候选：orb 接管像素的同时把 D2D 静态帧 ULW 一次保底 / 或由球窗继续画一帧再让 orb 覆盖 / 或 orb 窗收编球窗职责）——本裁决只定性，不开方、不动码（裁决者≠实现者）。

### 7.4 票面 6 格 AC 判（A1594 改道后口径）

| # | AC | 判 | 读数/依据 |
|---|---|---|---|
| 1 | 参考图入库＋逐条量化分析 | **沿用 §4 判**（替代件 65-style-direction.md 判 PASS 在案）+ 本轮补充：本路线上游件 LerSent001/orb 的 13 预设导出页本身即"参考实现"，逐字节上舰（§7.1#1） | 沿用 PASS |
| 2 | 属性映射 tokens；零新增硬编码色值 | **部分（改道后口径待收）** | 本路线：页内色值在**上游 HTML 资产内**（集成路线的正当形态，资产零改动是本票的硬要求）；宿主侧 Go 代码（orb_windows.go/orbasset.go 等 6 件）人工复读：色值仅 CSP 字符串与注释，零绘制字面量。但"映射到 C21 DesignTokens"这条 AC 在集成路线下的**对应物**（页面参数如何经 setState/setAudioBands 对齐 C21）尚未有任何腿登记——不改判历史格，改道后的归属留给编排者 |
| 3 | 并排差分产物在档 | **真桌面件，未跑** | 机器侧像素证据已在（像素用例 35.43%）；三张并排（现状/返工后/参考图）仍需机主桌面趟 |
| 4 | 票 62 dock 测试仍绿且行为断言未改 | **部分** | "未改"=**PASS**（§7.1#4，0 字节 diff）；"仍绿"=dock_test.go 12 枚确定性测试全 PASS，但 live 侧 TestBallLiveEdgeDock 存量红（基线同红，非本票引入）且**球窗透明让"dock 的球"这一可点击实体在 orb 场景失义**（§7.3）——判部分，落点=收尾腿 |
| 5 | Sleeping 差分重跑（timers=no / cpu≤0.5% / WS≤25MB / px_delta≥2120） | **未跑——且基线已换渲染器，读数将整表重定义** | D2D 路径的旧差分对 WebView2 渲染器无效（集成腿证据件"已知边界"自报同点）；webview2 常驻 GPU 合成面的空闲功耗必须仪器重测，本裁决不代跑不估读 |
| 6 | **owner 第二次看 `-tour` 并说"这次对了"** | **归机主**（⛔ 本裁决不判） | 集成腿的 tour 入口=balldebug（20 态循环）；判据仍是人眼判决。在此之前 SPEC-08 §2 INTERIM 不摘 |

### 7.5 判语总表

- **a2e535674（集成腿）**：实现腿声明范围内**全部属实**（资产逐字节/映射钉死/回退结构/门禁）；但留下**集成面真回归**——球窗被守卫停画后全透明不可点（§7.3，探针实证）。
- **d687dd910（出屏腿）**：**PASS 全格**。基线 PASS（orb=(3196,1260)-(3328,1392) 全内）＋突变 ②必红（(3452,1292)-(3584,1424) 双轴出屏）＋纯函数表测六形状＋WS_EX 语义零动，复跑与自报一致。
- **8fb821bba（黑窗腿）**：**PASS 全格**。三层根因证据链自洽；像素用例基线 35.43%＋突变 ①0.000% 全黑必红——**证明像素回归有牙**；资产 sha 复测一致；CSP 以 meta 形态保住（文档 API 只紧不松，形态合法）。
- **票 65 整体**：不结案。AC#3/#5 真桌面件未跑；AC#6 归机主；**新增收尾项：球窗透明回归（§7.3）须由实现腿修复后再走桌面趟**。

### 7.6 撤销口令

- 集成整体：`git revert a2e535674 d687dd910 8fb821bba`（球回 D2D wensity 路线）。
- 仅黑窗修：`git revert 8fb821bba`；仅出屏修：`git revert d687dd910`。
