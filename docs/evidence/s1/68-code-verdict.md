# 票 68 · 非实现者代码裁决（W68b/68-r2 码：4592c76a5 + bca0990f8）· 68-code-verdict.md

- **日期**：2026-11-06 19:3x +08（裁决腿＝非实现者；实现腿＝agent-ticket68-r2，交付表 `68-r2-settling-dock.md`）
- **被裁决对象**：commit `4592c76a5`（(a) Settling 合成落 0.95 ＋ (b) 停靠露出修到 ≈42%）＋ `bca0990f8`（交付表）。
  验证锚＝主树 HEAD `86fb145ee`（A1584 台账件，两笔票 68 commit 均为其祖先，`git merge-base --is-ancestor` 现跑通过）。
- **规则 21 隔离**：`git clone --shared` → `/d/tmp/wisp/20261106/68code-verdict/repo`（HEAD＝`86fb145ee`，clone 后 `status --porcelain` 空；
  本文件全部门禁读数出自这棵克隆树的新鲜复跑。⛔ 主树 gc/prune 前别动——alternates 指回主树；克隆树按规则 8 只建不删，留在原处）。
- **环境**：Git Bash（`D:\work\soft\Git\bin\bash.exe`），Go 门禁在克隆树内现跑；`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；
  `GOTMPDIR/TMPDIR/TEMP/TMP` 已 unset；**未设 `CGO_ENABLED`（更未置 0）**。
- **硬禁令遵守**：零 GUI 进程、零 `balldebug.exe`/`wisp.exe` 运行、零 `*winlive*` 测试、零桌面像素测量（前腿 74 分钟挂死教训）。
  本裁决全部证据＝静态读码 ＋ 非桌面新鲜门禁。

## 1. 门禁读数（克隆树内、本腿新鲜复跑，非采信交付表自述）

| 门 | 命令 | 读数 |
|---|---|---|
| gofmt | `gofmt -l internal/ball` | **空** |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/ball` | **空**（v0.12.0） |
| vet | `go vet ./internal/ball ./cmd/balldebug` | **clean** |
| test ball | `go test ./internal/ball -short -count=2` | **ok 0.347s** |
| test balldebug | `go test ./cmd/balldebug -count=2` | **ok 0.062s** |
| d22scan | `sh scripts/d22scan.sh` | **clean，rc=0**（尾部：ban #8 internal/=908、cmd/=217、frontend/=116 文本件，ban #10 internal/risk/=21，"no D22 ban violations"；首跑即绿——本环境 Git Bash 自带 `git.exe`，未复现交付表记的 6 个 FAIL 环境因） |

定向钉（`go test ./internal/ball -count=1 -v -run ...`，全部 PASS）：
`TestDockFollowsTheDrawnBody` · `TestSleepingSizeTruthTable` · `TestVisualForCoversAllTwentyStates` · `TestHitTestAndDPIInjection` · `TestAnimationPolicyZeroTimerInSleeping`。

## 2. 静态确认（逐条现读 `4592c76a5` 的 diff 与 HEAD 上的现文件）

1. **settlingAlpha 合成 0.95（非旧 0.33）＝确认**。`statevisual.go:316-337` 纯函数：prototype 终点 `RestSettledOpacity`（`tokens.go:416`＝**0.95**）、
   frozen 终点 `SleepOpacity`（`tokens.go:410`＝0.35），两端 clamp（p≤0→1，p≥1→end），测试钉了 10 步单调不升。
   合成链现读闭环：`ball_windows.go:347` `applyStateLocked` 里 `curVisual = VisualFor(..., StateSettling, 0)`（fadeProgress=0 ⇒ 框体 Opacity=1.0），
   `:478` `AnimFade` 分支 `b.constAlpha = u8(settlingAlpha(p)*255)` ⇒ **合成＝框体 1.0 × 窗口 0.95 ＝ 0.95**；
   旧码（本 commit 前）`AnimFade` 终点写 `SleepOpacity=0.35` 叠乘视觉 0.95 落点 ⇒ ~0.33，与 R15#3 记载一致，已消除。
   frozen 侧合成＝1.0×0.35＝0.35，§2.1 逐字保留，逃生门不变。
2. **dockOrbRadius 用态本体半径（34.72→17）＝确认**。`dock.go:131-143` `dockOrbRadius(stateSizePx)=int32(SizePx/2)`；
   调用点 `dock_windows.go:66` `orbR = dockOrbRadius(b.curVisual.SizePx)`；Sleeping 静息体＝`stateSize`（`statevisual.go:120-138`）＝
   56×`SleepRestRatio 0.62`＝34.72 ⇒ `int32(34.72/2)`=17（测试钉 17 在 `TestDockFollowsTheDrawnBody` 与 `TestSleepingSizeTruthTable` ③ 双处）。
   非 Sleeping 态仍走 28（56/2），注释与测试口径一致。窗口半径（edgePx/2−margin）旧读法已替换。
3. **drawGlass 吃进 DockSquash（42%）＝确认**。`renderer_windows.go:530-537` dockX/dockY 上提到函数开头（`DockSquash(v.DockProgress)`，
   `tokens.go:382` `DockOverlapFrac=0.42`）；吃进挤压的层＝玻璃壳 `:565` `fillEllipse(c, R*dockX, R*dockY, r.body)`、三个液团 `:579`、
   远侧折射核 `:586`、高光 `:590-591`、rim/lip/边框四圈 `:596-605`（新增 `arcEx` 逐轴椭圆 polyline；`arc`/`pointOn`/`circle` 保留为圆特例，ABI 纪律注释未动）。
   露出算术：off=int(17×0.42)=7，压扁本体沿轴 2×17.36×0.42≈14.6px＝34.72px 直径的 **42%**（旧读法 int(28×0.42)=11＋圆壳 17.36≈29.1px≈82%）。
   钉在 `TestDockFollowsTheDrawnBody`（42%±1% 算术钉）与 `TestSleepingSizeTruthTable` ③（47px 行保留为**注入半径 28** 的纯法则钉＋新 call-site 17/43px 钉）。
4. **命中半径未动、已在场＝确认**。`git show 4592c76a5 --stat -- internal/ball/hit.go` ＝ 空（本 commit 零触碰）；
   HEAD 上 `hit.go::HitTest` 在场：`r = orbPx/2+ClickTolerancePx`，可见外沿 `edge = ceil(orbPx*0.75)`（=1.5R）取大——R15#5 交付的读法原样，
   钉 `TestHitTestAndDPIInjection` 绿。
5. **prototypeVisuals 默认 true ＋ 逃生门完好＝确认**。`statevisual.go:101` `var prototypeVisuals = true`（翻转本体是 68-r1 `4b9efc889` 交付的，
   本腿只改注释：删掉"0.33/82% 已按实现值接受"的过期句——那是 2026-11-01 的中间读法，已被台账"R15 全按推荐＝#3→0.95、#4→修~50%"取代）；`EnablePrototypeVisuals` `:105` 在场；
   生产侧逃生门 `cmd/balldebug/main.go:128` `ball.EnablePrototypeVisuals(!*frozen)` 在场。
6. **地雷修复在场＝确认**。`live_windows_test.go::TestBallLiveAudioLiquidGate`（函数定义现于 `:640`，票面写的 `:609` 行号已漂移，所指对象没变）
   `:648` 显式 `EnablePrototypeVisuals(false)`，注释逐字写明"翻转与这条修正在同一个 commit"——与 AC#2 地雷要求吻合；frozen 段后 `:664` 显式回 true，顺序正确。
7. **D32 纪律＝确认**。`AnimationPolicy{AnimFade, framePeriod}` 逐字未动（Settling 260ms 一次性 burst，30fps 上限内）；
   `TestAnimationPolicyZeroTimerInSleeping` 绿；`dock_windows.go` 零 SetTimer 约束未动；渲染改动全为静态几何/纯函数，无逐帧 COM 分配新增。

## 3. 逐 AC 裁决

- **AC#1（尺寸真相三列表）＝ PASS**（维持既有勾框；不重开）。本腿新鲜复跑 `TestSleepingSizeTruthTable` PASS 作旁证；34.72px 口径与 R15#1 裁定一致。
- **AC#4（非桌面四道门）＝ PASS**。§1 表全部出自克隆树新鲜复跑：gofmt/gofumpt 空、vet clean、两包 `-count=2` ok、d22scan clean rc=0。
- **AC#2（默认值翻转）＝ 代码半就现有证据判 PASS；整格仍 blocked-on-owner**。
  - **判据与范围**：翻转本体（`prototypeVisuals=true`）＋ 断言迁移 ＋ 地雷修复在 68-r1 `4b9efc889` 交付、编排者已收；本裁决对象是 68-r2 两笔把
    "被签收的样子"的另外两个偏差（Settling 合成、停靠露出）修到承诺值的**代码半**。逐条静态确认见 §2；1–7 全部成立，
    且每一处行为都钉了测试（本腿新鲜跑绿）。**就静态＋非桌面证据能到的射程，AC#2 的代码是正确的**：
    默认构建 == owner 签收样子所需的全部码内条件（默认 true、0.95 合成、42% 露出、1.5R 命中、逃生门、地雷）现已同场在场。
  - **不能由本裁决背书的**：0.95/42% 两个数是纯几何/纯函数推演＋码读，**未经桌面像素级证实**；AC#2 自身的"需桌面复测"半格与 AC#3 整格
    **headless 不可做**，维持 blocked-on-owner。票面勾框我不动（非实现者零票面编辑）。
- **AC#3（签收面真桌面复测）＝ 未跑、不可跑，维持 blocked-on-owner**。硬禁令（零 GUI/winlive/桌面像素测量）下本腿零运行；
  owner 须照票面 72–85 行命令集在真桌面逐条跑（见 §4）。

## 4. owner 真桌面必跑命令清单（票面 72–85 行逐字抄录，未做任何改写）

前置：桌面空出、显示器缩放 **100%**（A.2 全部数字是 96 DPI 口径）；**第 5 条 `-frozen` 对照不是可选项**。

1. `export PATH="$PWD/third_party/sherpa-onnx:$PATH"`
2. 洁净前置：`tasklist | grep -iE "wisp.exe|balldebug.exe"` **必须无输出**；确认缩放 100%
3. `go build -o build/balldebug.exe ./cmd/balldebug`
4. 自由＋靠边两行：`build/balldebug.exe -diff build/t68-flip -diff-states Sleeping,Listening,Thinking,Acting,Speaking,Warm,Settling -diff-sample 5s -diff-dock right`
   ⇒ 门：`Sleeping` px≥8 **≥ 2098**、像框落在 40–46 家族、`timers=no`、`cpu_all ≤0.5%`、teardown `clean`
5. **冻结对照（逃生门验活）**：同第 4 条加 `-frozen` ⇒ 期望 **px≥8 = 0**、框 `0x0`
6. 零定时器与句柄门：`build/balldebug.exe -state Sleeping -hold`（Ctrl+C 结束；打印 `timers=false` 与 `closed handles=` ≤600）
7. `go test -tags winlive -count=2 -v ./internal/ball/`
8. 翻转后重跑 AC#4 四道门（＝本文件 §1 表那五行，桌面机上再跑一遍）

数字不过预算就报 FAIL 并附全样本；不许调阈值、不许重测到运气好的那次、不许用平均抹掉坏尾部。

## 5. 残留观察（不构成退回，登记备查）

1. **caustic（接触阴影）/halo（外光晕）两层未吃挤压**（`drawGlass` 步骤 0/1 仍按未压扁半径画）：它们是本体之外的阴影/辉光，
   不影响"本体露出 42%"的算术（42% 钉是本体相对的），但靠边 tab 的光晕外沿可能越过屏边被裁——真桌面差分化像若在 40–46 家族之外，
   先看这两层再定性质。
2. **票面行号漂移**：地雷所在行号票面写 `:609`，现函数定义在 `:640`、修复在 `:648`——对象没变，仅提醒后续引用别按旧行号找。
3. 克隆树 `/d/tmp/wisp/20261106/68code-verdict/repo` 留场（规则 8 只建不删）；主树 gc/prune 前勿动它（alternates 指回主树）。

## 6. 纪律声明

零票面编辑、零 SPEC/契约编辑、零 `docs/reports/pending-and-issues.md` 接触、零 frontend/design 接触、零 push；
本文件以显式 pathspec 单独 commit；全部读数可按 §1/§2 的命令与行号现跑复核。
