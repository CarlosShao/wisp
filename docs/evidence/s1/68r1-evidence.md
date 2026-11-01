# 68-r1 — 票 68 翻转腿交付：`prototypeVisuals` 默认 true ＋ 命中半径对齐可见外沿 ＋ 7 态签收截图

- 腿：`68-r1`（写腿，票 68 续腿，Claimed by agent-ticket68 的延续）
- 日期：2026-11-01 22:0x–22:5x +08；基线 HEAD `e99bd1a4f`（dev，开工三核过：索引干净、写面无他人脏）
- R15 材料机主批复：registry **A1378**「R15 全按推荐」＝①Settling 合成 ≈0.33 接受 ②停靠露出 ≈82% 接受 ③命中半径 21px→26px 修正
- 票面 AC#1/AC#4 前腿已交付未重开；本腿交付 **AC#2（翻转）＋AC#3 的桌面半（差分复测＋签收材料）**；⛔ 全程零编辑 SPEC-08／零编辑票面

## 1. 翻转本身（AC#2）

`internal/ball/statevisual.go:102`：`var prototypeVisuals bool` → `var prototypeVisuals = true`，
注释块按新事实重写（R15 全按推荐 ⇒ 默认构建＝owner 签收的样子；SPEC-08 §2 Sleeping 行 34.72px 即默认配置；
`EnablePrototypeVisuals(false)` 与 `balldebug -frozen` 保留为逃生门）。①②两项按实现值接受＝零代码改动。

`cmd/balldebug/main.go` 的 `EnablePrototypeVisuals(!*frozen)` 语义不变：无旗标＝新默认（true），`-frozen`＝false。

## 2. 命中半径 21px→26px（R15#5，唯一参数面改动）

**推荐原文**（R15#5 / Q-2，机主已按推荐批）：**对齐可见外沿**，代价「改 `HitTest`，无视觉变化」。
**实现**（`internal/ball/hit.go::HitTest`）：

```go
r := float64(orbPx)*scale/2 + ClickTolerancePx*scale
if edge := math.Ceil(float64(orbPx) * scale * 0.75); edge > r { r = edge }
```

即**可点圆＝max(本体半径＋点击容差， 可见外沿 1.5R)**（向上取整：边缘一像素算球）。

- **为什么不是把 `ClickTolerancePx` 4→9**（派单「唯一参数改动」的字面读法）：`ClickTolerancePx` 被
  `docs/evidence/s1/c21-native-tokens.md:239` 行以「点击容差 4px」钉值，
  `TestC21GeometryRowsMatchCodeConstants` 双向核（表→码要求行内点名的常量存在；值要求码值必须是行内未认领数字）。
  改值或删常量都会打红 C21 门，而该表**不在本腿写面**（写面＝`internal/ball/**`＋`cmd/balldebug/**`＋`68r1-*`）。
  而且单一容差在任何配置下都不可能全态对齐可见外沿（外沿随 D 缩放，容差不随）。owner 批的是「对齐可见外沿」这条原理，
  上式是它的直接实现，并在命名格精确落出派单数字：**Sleeping 静息体 21px → 26px**。
- **数账（96 DPI、默认 56）**：调用点 `ball_windows.go:560` 传 `int(curVisual.SizePx)`＝int(34.72)=34；
  旧 r＝34/2+4＝**21px**；新 r＝ceil(34×0.75)＝**26px**。画出的光晕填充外沿＝`haloExtents` 的 fill 半径
  1.5×17.36＝26.04px（RingMarginPx=22 之后全档位不被窗裁）——命中圆与可见外沿差 0.04px（调用点 int 截断的既有亚像素残差，
  不在本腿扩权修）。其余各态同时从 32px 抬到 42px＝各自可见外沿（Q-2 的全态推广）；**frozen 微点保持 10px 不缩**
  （6+4 > 1.5×6=9，容差下限兜底——R15#5 原文点名「frozen 更糟：10px」的方向不被倒退）。
- `ClickTolerancePx` 注释同步改为「可点圆的下限」；新增行为钉：`TestHitTestAndDPIInjection` 里
  r=25 可点 / r=26（边界，`<=`）可点 / r=27 穿透。

## 3. 测试联动账（票面 AC#2：「逐条迁移＋列出原来断什么、现在断什么、为什么等价」）

| 测试 | 原来断什么 | 现在断什么 | 为什么等价 |
|---|---|---|---|
| `TestVisualForCoversAllTwentyStates` Sleeping 12px@0.35 | 默认（旧＝frozen）下 §2.1 冻结读数 | 默认（新＝原型）断静息玻璃体 34.72px@1.0＋Glass；**12px@0.35 原样保留**在显式 `EnablePrototypeVisuals(false)` 块 | 旧数字一条未删未放宽，只是挂到它真正所属的模式上；默认档新增的是 owner 签收读数 |
| 同上「Sleeping 任意配置 12px」 | 冻结微点不随配置 | 原型档断「随配置放大（下限 30）」；frozen 块原样保留 12px 断言 | 同上，SPEC-08 §2 INTERIM 的两档读法各归各 |
| 同上 Armed 0.6 / Muted 0.4 | 冻结档透明度 | 默认断 0.92 / 0.85（R15#1 Q-1 实测阶梯＝默认目标）；0.6/0.4 进 frozen 块 | 玻璃材质正是靠高不透明度在浅色桌布成像（62-diff-baseline 的 0 像素教训） |
| 同上 Settling 末端 0.35 | 冻结档落点 | 默认断 `RestSettledOpacity`=0.95；0.35 进 frozen 块；入口 1 两档皆真、原样保留 | R15#3：合成 ≈0.33＝0.95(Visual)×0.35(ULW) 按**实现值接受**，本腿零改 |
| 同上（新增） | — | `PrototypeVisualsEnabled()==true` 钉翻转本身 | 翻转的回归钉 |
| `TestBallLiveAudioLiquidGate` | 第一段「frozen mode」靠**旧默认＝false** 才成立（票面点名的地雷，`:618` 注释自称冻结） | 显式 `EnablePrototypeVisuals(false)` 后再建窗 | 票面要求「翻转与这条修正在同一个 commit」；restore/cleanup 前腿已备好 |
| `TestBallLiveLifecycle`（winlive，翻转后发现） | 状态环游断 Acting `TimersAlive=false`——旧默认 frozen 下无 burst | Acting 特例：等有界液面 burst 自退休（≤2s waitFor）后断无**持久**定时器；`want` 表达式简化为 `Animated(s) ‖ Warm ‖ Settling`（语义不变） | 票 62 burst 是 transitionDriven 的有界过渡（草案表第 6 行「实测 timers=yes ⇒ 是票 62 的 burst」）；默认一翻它就出现，纪律＝「不许可持久定时器」，不是「不许出现」 |
| `TestHitTestAndDPIInjection` | — | 新增 R15#5 三点半径钉（25/26/27） | 新契约数没有旧断言可迁，属新增钉 |

迁移后非 winlive 套件 `-count=2` 全绿；winlive 套件 171/172（见 §6 的环境红）。

## 4. 桌面证据（AC#3 的定量半）

**环境发现（如实报）**：默认布点（屏幕中央）当刻被活窗口占满——Task Manager 进程表持续刷新、QQ Mail 弹窗、
活跃微信会话；alive/dead 差分把背景应用的重绘全算进去 ⇒ 第一次 `-diff`（16 行、teardown 全 clean）的
px/框列**不可与 A.2 家族比较**（例：Sleeping 5450px/190×232、Armed 34256px/mean=18），原始表留
`D:\tmp\wisp\20261101\68r1\diff-flip\diff-table.txt`。**timers/CPU/handles/teardown 列不受影响、全部有效**。

**干净区协议**：勘察后把球放到浏览器静态浅色表格区（`-x 1080 -y 330`，窗口 100px、抓取 340×340 全落在静态背景），
A.2 的白底可比性恢复。`-diff-states Sleeping,Listening,Thinking,Speaking,Warm,Settling,Armed -diff-sample 5s`：

| state | timers | cpu_all% | privWS_MB | px≥8 | 框 | teardown |
|---|---|---|---|---|---|---|
| Sleeping | **no** | **0.000** | 11.98 | **2106** | **40×39** | clean |
| Listening | yes | 0.002 | 15.36 | 5555 | 80×82 | clean |
| Thinking | yes | 0.004 | 14.50 | 5555 | 85×77 | clean |
| Speaking | yes | 0.002 | 13.32 | 5539 | 80×80 | clean |
| Warm | yes | 0.001 | 13.77 | 5225 | 70×69 | clean |
| Settling | yes | 0.000 | 12.54 | 4912 | 68×68 | clean |
| Armed | no | 0.000 | 12.11 | 5524 | 76×77 | clean |

- **AC#3 门逐条**：Sleeping px≥8 **2106 ≥ 2098** ✓；框 40×39 ∈ 已录 40–46 家族（旧录 2098/2103/2120，框 40×38/44×44/46×46）✓；
  timers=no ✓；cpu_all 0.000 ≤0.5% ✓；teardown clean ✓。
- Listening/Thinking/Speaking/Armed 的框 76–85px 大于旧录的 72×72 封顶：旧值是 RingMarginPx=8 时代被 72px 窗**裁切**的读数，
  现在 100px 窗不再裁波纹/光晕——与现行几何一致，不是回归。
- 停靠行：干净布点跑不了 dock（`-diff-dock` 钉死右屏边缘＝Task Manager 占用区），第一次跑的 dock 行同样被污染
  （例：Sleeping-dock 2127/181×189）。**dock 行的定量复测保持开放**，等右缘清场的时刻；停靠露出 82% 本身 R15#4 已按实现值接受。

**冻结对照（票面第 5 条，非可选）**：同布点加 `-frozen`，Sleeping：**px≥8=0、框 0×0、max=1**（期望逐字一致）——逃生门健在。

**句柄门（票面第 6 条）**：`-state Sleeping`（2s dwell 自退；不进 `-hold` 是因为无人值守跑没有 Ctrl+C 可发，
同样的两行读数照打印）：`state=Sleeping handles=376 timers=false` → `closed handles=379`，`balldebug: OK`（379 ≤600）。

## 5. 7 态签收截图（＝本文件同目录 `68r1-<State>.png`，7 枚）

干净布点跑的 alive 真桌面帧（BitBlt+CAPTUREBLT 分层合成，非自绘合成图）：Sleeping（静息玻璃体无边框）、
Listening（accent 环＋audio-lines）、Thinking（accent 芯＋底部光带）、Speaking（success 芯＋volume）、
Warm（暖芯）、Settling（≈0.33 合成落点）、Armed（边框＋halo）。逐枚人工目验过。

## 6. winlive（票面第 7 条）与环境红

`go test -tags winlive -count=2 ./internal/ball/`：**171/172 过**（含修复后的 `TestBallLiveLifecycle`、
`TestBallLiveEdgeDockHover`、`TestBallLiveAudioLiquidGate`、`TestLiveTransparentCornerFallsThrough` 等）。

唯一红＝`TestLiveConfirmingCancelAndEscReturned`（count=2 两轮＋单跑一轮，同一断言 `interaction_live_test.go:239`
「注入的物理 Esc 未取消 Confirming」）：测试前半段的**合成** WM_HOTKEY 取消路由、Esc 借走/归还两端全部通过，
只有**物理注入**的 Esc 没有变成 WM_HOTKEY。归因＝桌面环境：机主正在活跃使用（微信输入法组字、前台窗口、可能的低级钩子
在 RegisterHotKey 派发之前吃键）；该测试自己的 SKIP-LOUD 机制就是为「注入不落」这类不稳设计的。本腿三处生产改动
（默认值/命中圆/注释）零键盘路径，机制上无法产生此失败。CI 对 winlive 只跑编译门（`ci.yml` `go vet -tags winlive`），
winlive 测试从未入 CI——这条绿期望本来就只在安静桌面的手动跑里成立。**建议编排者找安静窗口复跑该一条**（命令：
`go test -tags winlive -count=1 -run TestLiveConfirmingCancelAndEscReturned ./internal/ball/`）。

## 7. 四数＋d22scan（最终树，含全部本腿改动）

```
gofumpt -l internal/ball cmd/balldebug   -> （空）
go vet ./internal/ball/ ./cmd/balldebug/ -> clean
go test -count=2 ./internal/ball/        -> ok 0.340s
go test -count=2 ./cmd/balldebug/        -> no test files
d22scan                                  -> clean, no D22 ban violations
```

基线（改动前）同套件亦先跑过一遍全绿留档：`D:\tmp\wisp\20261101\68r1\baseline-ball.txt`。

## 8. 边界自查

- SPEC-08 §2/§2.1 零字节改动；`docs/evidence/s1/c21-native-tokens.md` 零改动；票面零改动；台账零改动（本文件即申报面）。
- 渲染逻辑零改动（`renderer_windows.go`/`liquid*.go`/`dock*.go` 未碰）；命中半径例外落在 `hit.go::HitTest` 一处。
- thresholds/golden/testdata 仓根未碰；`ClickTolerancePx`=4.0、`RingMarginPx`=22.0 数值未动（C21 表仍逐值为真）。
- 临时件全在 `D:\tmp\wisp\20261101\68r1\**`（只建未删）；无 push。
