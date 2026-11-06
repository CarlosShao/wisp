# 票 68 · 再派腿 68-r2（Settling 合成 0.95 + 停靠露出 ≈42%）· 实现腿交付表

- **日期**：2026-11-06 19:1x +08（执行程 `agent-ticket68-r2`，派单＝A1583 后的窄口径再派）
- **起手**：主树 HEAD `ab808c359`（派单锚 f7505f122 的当前 HEAD），主树 `git status --porcelain -- internal/ball cmd/balldebug` ＝ **0 行**
- **规则 21 隔离**：`git clone --shared` → `/d/tmp/wisp/20261106/68r2/repo`（branch `dev`，clone 后 `status --porcelain` 空；本腿全部实现与门禁读数出自这棵树；⛔ 主树 gc 前别 prune——alternates 指回主树）
- **写面**：`internal/ball/**` + 本表。SPEC-08 / C21 表 / 票面 / 台账 **零编辑**；`docs/reports/pending-and-issues.md` **零接触**；零 push。
- **环境**：Go 1.27.1 windows/amd64，CGO 未清零（未设 `CGO_ENABLED=0`），`GOTMPDIR/TMPDIR/TEMP/TMP` 已 unset。

## 0. 起手现量核对（哪些 AC 已在 HEAD 在场，避免重开）

HEAD 已带 68-r1（commit `4b9efc889`，2026-11-01）交付的半边，本腿**逐条现读确认、未重开**：

- **(c) 命中半径对齐可见外沿**＝已在场：`hit.go::HitTest` 可点圆＝max(本体半径+ClickTolerancePx, 可见外沿 1.5R=ceil(0.75·orbPx))；调用点 `ball_windows.go:561` 传 `int(b.curVisual.SizePx)`（Sleeping 静息体 int 截断 34 → ceil(34·0.75)=26px）。**本腿零改动**，测试钉 `tokens_test.go::TestHitTestAndDPIInjection` 在场且绿。
- **(d) `prototypeVisuals` 默认 true**＝已在场（`statevisual.go:96`）；`EnablePrototypeVisuals(false)` 与 `cmd/balldebug:128` 的 `-frozen` 逃生门在场。**本腿零改动**（只改其注释，见下）。
- **(e) 旧默认断言迁移 + 地雷**＝已在场：`tokens_test.go` 二十态锚点已逐条迁到新默认读数、冻结数字原样挂显式 `EnablePrototypeVisuals(false)` 块（零删零宽）；`live_windows_test.go::TestBallLiveAudioLiquidGate` 已显式 `EnablePrototypeVisuals(false)`（同 commit 在场）。**本腿零改动**。

## 1. 本腿改动（编号对齐派单 (a)/(b)）

### (a) Settling 合成 → RestSettledOpacity 0.95（修 0.35×0.95 双乘）

- **机制**：Settling 的 260ms 渐隐由 ULW `SourceConstantAlpha` 独驱（框体本身画的是 fadeProgress=0 的静息体），终点却是冻结微点的 `SleepOpacity=0.35`，与视觉侧落点 `RestSettledOpacity=0.95` 叠乘 ⇒ 合成 ≈0.33（R15#3 实测量级）。
- `statevisual.go`：新增**纯函数 `settlingAlpha(p)`**＝淡出终点唯一权威：prototype 终点 `RestSettledOpacity`（0.95）、frozen 终点 `SleepOpacity`（0.35），两端 clamp、单调；`ball_windows.go::onAnimTick` 的 `AnimFade` 分支改读它（驱动代码从 9 行缩到 7 行，定时器策略/stop 时点不变）。`anim.go` 注释同步。
- `statevisual.go:88-96` 的 `prototypeVisuals` 块注释改写：删掉"0.33/82% 已按实现值接受"的过期句（那是 2026-11-01 的中间读法，已被台账后续「R15 全按推荐＝#3→0.95、#4→修 ~50%」取代），改为指向 `settlingAlpha` 与 `dockOrbRadius` 两条现值。
- **合成结果（prototype）**：框体 Opacity 1.0 × 窗口 alpha 终点 0.95 ＝ **0.95**；frozen：1.0 × 0.35 ＝ 0.35（逐字保留 §2.1，逃生门不变）。
- **测试钉**：`TestVisualForCoversAllTwentyStates` 新增 `settlingAlpha` 三向钉（p=0→1、prototype p=1→0.95、10 步单调不升）＋冻结块内 `settlingAlpha(1)→0.35` 钉。

### (b) 停靠露出 → DockOverlapFrac ≈42%（"藏一半"，取代 ~82%）

- **机制**（两处合成才落得到承诺值）：
  1. **`dock_windows.go::dockGeometry` 喂给 `DockPos` 的 `orbR` 由窗口半径（edgePx/2−margin）改为态半径**＝新增纯函数 `dockOrbRadius(stateSizePx)=int32(SizePx/2)`（`dock.go`）。窗口半径在两处皆错：Sleeping 静息体 34.72px 装在 56px 窗里（17.36 vs 28），且 DPI>96 时本体不缩放而窗口缩放（`statevisual.go` 的 DPI 注记）。
  2. **`renderer_windows.go::drawGlass` 让壳体真正吃进 dock 挤压**：dockX/dockY 上提到 step 2 之前，玻璃壳 `fillEllipse(c, R·dockX, R·dockY)`、远侧折射核、rim/lip/边框三圈改走新增 `arcEx`（逐轴半径椭圆 polyline，`arc`/`pointOn` 保留为 `arcEx` 的圆特例；ABI 纪律不变——DrawLine 宽度仍在栈槽浮点位）。**此前壳体是圆的**，与 dock.go 自己写的"tab 就是压扁的球"法则相违，也让 82% 的泄漏无法修到 42%。
- **数字（96 DPI、默认 56、Sleeping 静息体 34.72px）**：`orbR=17` ⇒ `off=int(17×0.42)=7`；压扁本体沿轴 2×17.36×0.42≈14.6px ＝ 直径的 **42%**（=DockOverlapFrac，"文档承诺 42%"；台账记作 ≈50% 的"藏一半"语义）；整窗留在屏上 36+7=43px（72px 窗口口径）。改动前：off=int(28×0.42)=11 ＋ 圆壳 17.36 ＝ 29.1px ≈ **82%**。
- **测试钉**：`TestSleepingSizeTruthTable` ③ 的注释与断言迁移（旧 47px 行保留为注入半径 28 的**纯法则钉**并改写说明；新增 call-site 半径 17 钉 + 43px 钉）；`dock_test.go` 新增 `TestDockFollowsTheDrawnBody`（非 Sleeping 态 96 DPI 仍 28、静息体 17、57/100px 窗口径、42%±1% 露出算术钉）。

### 未改（D32 纪律）

- **零新增定时器**：Settling 的 `AnimationPolicy{AnimFade, framePeriod}` 逐字未动；Sleeping 零定时器钉（`TestAnimationPolicyZeroTimerInSleeping`）未动且绿。dock 仍零 SetTimer（`dock_windows.go` 头注约束未动）。渲染改动全是静态几何，无逐帧 COM 分配新增。

## 2. 门禁读数（全部出自克隆树，改后现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `gofumpt -l internal/ball` | **空**（v0.12.0，`/d/work/base/gopath/bin/gofumpt.exe`） |
| vet | `go vet ./internal/ball`（另 `./cmd/balldebug`） | **均 clean** |
| build | `go build ./cmd/balldebug` | **rc=0** |
| test | `go test ./internal/ball -short -count=2` | **ok 0.330s**（`-v` 口径 `--- PASS` 148 行）；`go test ./cmd/balldebug -count=2` **ok 0.051s** |
| d22scan | `sh scripts/d22scan.sh` | **clean**（0 个 `--- FAIL`；尾部读数：ban #8 internal/=908、cmd/=217、frontend/=116 文本件，ban #10 internal/risk/=21，"no D22 ban violations"） |

**口径说明**：
- `-short` 在本包是**空操作**：包内没有任何 `testing.Short()` 分支；桌面/live 测试是被 `//go:build windows && winlive` 构建标签排除的（本腿未带该标签编译），不是被 `-short` 跳过的。
- d22scan 首跑 6 个 FAIL（`TestWalksSkipGitIgnoredPaths` 等）经查**全部同一环境因**：测试子进程 `%PATH%` 无 `git.exe`（`scan_test.go:1463` 原文）；把 `/d/work/soft/Git/cmd` 加进 PATH 复跑即 0 FAIL——与代码改动零机制关联（二跑日志 `d22-after2.log` 在克隆树外留档）。

## 3. 明确未测（留给真机，本腿零 GUI 进程）

- **AC#3 真桌面复测（整格未做）**：`build/balldebug.exe -diff` 差分化像（Sleeping px≥8 ≥2098、像框 40–46 家族）、`-frozen` 冻结对照（px≥8=0/框 0x0）、`-state Sleeping -hold` 的 `timers=no` 与 `closed handles≤600`、`cpu_all≤0.5%`、WS≤25MB、winlive `-tags winlive` 全套。**本腿未运行任何 `balldebug.exe`/`wisp.exe`，未建窗、未点鼠标、未热键注入、未截屏、未开任何 GUI 进程**（前次 W68 腿因此挂死 74 分钟，本轮禁令遵守）。
- (b) 的 42%/0.95 两个数是**纯几何/纯函数推演＋码读**，未经桌面像素级证实——AC#3 落真机时按票面第 70-85 行命令集复核。
- 提交：克隆树两笔（码＋本表），主树以字节同 patch 落地两笔同文提交；**均未 push**。
