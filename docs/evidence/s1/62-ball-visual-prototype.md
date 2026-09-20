# 62 — 液态玻璃悬浮球视觉原型：owner 签收证据

**Ticket:** .scratch/wisp/issues/62-liquid-glass-ball-visuals.md
**Scope of this file:** 可运行原型 + 实测证据，供 owner 实机签收。**SPEC-08 §2 与 docs/PLAN.md 未改动**
（契约变更须 owner 实机签收后由编排者走 D22 签字流程，措辞见文末第 7 节）。
**Machine:** DESKTOP-LVS7839 · Windows 10 Pro 26100 · i7-8750H（12 逻辑核）· 单屏 3440x1440 ·
浅色（白）桌布 —— 即 owner 实机。
**本文件所有数字都是 2026-09-20 在 dc8cd03 之后的工作副本上实测重跑的**（不是抄前面的记录）。

---

## 0. owner 要跑的那一条命令

在仓库根目录 `D:\work\workspace\projects plans\Wisp` 下，** cmd 窗口里三行**：

```
cd /d "D:\work\workspace\projects plans\Wisp"
set PATH=%CD%\third_party\sherpa-onnx;%PATH%
build\balldebug.exe -tour
```

- 第 2 行是因为这个 exe 要用 `third_party\sherpa-onnx` 里的语音 DLL，不给路径它会直接起不来。
- `-tour` = 自动放映：14 步、每步 5 秒，命令行里每一步都印一行中文提示，
  同时把"这一刻球有没有定时器、有没有吸附、句柄数"打在下一行。**你只要看着球 + 读提示。**
- 中途想停：命令行按 `Ctrl+C`，或右下角托盘图标右键 → Exit。
- 放映结束后**球会停在屏幕右缘变成一个"半隐的标签"**，最后一步是留给你自己动手的（见下表第 14 步）。
- 这一步不会留下任何进程：`Ctrl+C` / 托盘 Exit 都会把球窗口和托盘正常销毁。

### 按顺序看什么（14 步，一步一句话）

| # | 命令行会印的编号 | 看什么 | 已经替你量到的数 |
|---|---|---|---|
| 1 | tour 01 AT REST | 白色桌布上这个静止的玻璃球**看不看得清**、边光够不够。这就是平时挂在桌面上那颗。 | 定时器 no，CPU 0.000%，私有工作集 11.76MB |
| 2 | tour 02 SUMMONED | 按快捷键那一刻的样子：里面的液体开始流动、整球变亮。**没有声音**也在流，这是模型自己的"开场涌动"，几秒后自己停 | timers=yes（会话态才允许有） |
| 3 | tour 03 SPEAKING WITH A VOICE FED IN | 喂进去一条假语音（电平 0.70，每第 3 句压到 1/4）：液体**转起来、跟着起伏**，安静的那一句它会明显变懒 | 单核 0.016% |
| 4 | tour 04 STOPPED TALKING | 不说话了（还在听）：液体收住，**外圈边框"滑"进来**，不是"啪"地出现 —— 这就是你要的那个过渡动画 | 单核 0.023% |
| 5 | tour 05 BACK TO REST | 回到静止：边框和液体都松开，又是**一张不动的帧**。跟第 1 步对比一下，应该一模一样 | timers=no，CPU 0.000% |
| 6 | tour 06 DOCKED ON THE LEFT EDGE | 吸到左边缘：球被压扁成**贴在边上的一条**，只剩一条看得见、点得到 | 定时器 no，CPU 0.000% |
| 7 | tour 07 POPPED BACK OUT OF THE LEFT | 那条又**滑回完整球**，整个都在屏幕内 | 定时器 no |
| 8-9 | tour 08/09 | 右边缘：同上（吸住 / 弹回） | 定时器 no |
| 10-11 | tour 10/11 | 上边缘：同上。**注意它没有钻到任务栏底下**，也没跑到屏幕外 | 定时器 no |
| 12-13 | tour 12/13 | 下边缘：同上 | 定时器 no |
| 14 | tour 14 NOW YOUR TURN | **真鼠标**试：① 指针划过右缘那条 → 弹回完整球；② 指针移开 → 又缩回标签；③ **点一下**标签 → 弹出来并且待着不走；④ 把它拖到屏幕中间松手 → 自由悬浮 | —— 这一步只有你的真指针能演，见第 4 节 |

想改尺寸（契约允许 44–72，默认 56）加 `-size 44` 或 `-size 72`；想放大看细节加 `-x 1600 -y 500` 决定放哪。

### 顺手对比老规格（A/B，同一台桌面）

```
build\balldebug.exe -frozen -tour
```

`-frozen` 放的是**现在写在 SPEC-08 §2.1 里的那套**（纯色圆点 + 状态色环，6 个状态轮一遍），
每步同样印数字。它**没有**液体也没有吸附（那两样是本票原型专有的），提示里会直说，不会假装演过。
两条命令背对背跑一遍，就是"改之前 vs 改之后"。

想看四种配色的候选（第 5 节）：

```
build\balldebug.exe -tour -look glacier
build\balldebug.exe -tour -look nebula
build\balldebug.exe -tour -look solar
```

想完全自由地玩（真快捷键 / 托盘 / 点击都活着，不放电影）：`build\balldebug.exe -stay`。

---

## 1. 基线（改动之前）：Sleeping 态在白色桌布上不可见 —— 数值证明

测量方法（**不是**自绘合成图）：同一个屏幕区域截两次，一次悬浮球进程活着、一次把它正常关掉，
两帧逐像素相减。工具：`build\balldebug.exe -diff`（票 62 新增，父进程拉起子进程、正常收尾、
再截同区域，保证"死帧"里没有球也没有托盘残留）。

命令：

```
build\balldebug.exe -diff docs\evidence\s1\62-diff-baseline -diff-states Sleeping,Listening,Speaking -x 1600 -y 500 -diff-margin 90 -diff-sample 4s
```

图片：`62-diff-baseline\01-Sleeping-alive.png` / `-dead.png` / `-diff.png`（diff 已放大 6 倍，
所以 2/255 的差别也看得见）。数字表：`62-diff-baseline\diff-table.txt`。

| 状态 | 变化像素数（差值>=8/255） | 最大差值 | 单核 CPU | 私有工作集 | 活动定时器 |
|---|---|---|---|---|---|
| **Sleeping（现状 12px 微点 / opacity 0.35）** | **0**（>=3 的也只有 17 个） | **4/255** | 0.000% | 10.60MB | 无 |
| Listening（对照组，球是可见的） | 1471 | 115/255 | 0.008% | 12.23MB | 有 |
| Speaking（对照组） | 4824 | 102/255 | 0.027% | 12.90MB | 有 |

结论：**现状的 Sleeping 球在 owner 桌面上的可见像素变化为 0**，整个 232x232 区域里最大只有
4/255（1.6% 亮度）的差别 —— 人眼不可能分辨，这就是本票真正的 bug。对照组证明这套差分截屏
方法能看见球（Listening/Speaking 分别有 1471/4824 个像素变化 >=8/255）。

同时基线记录了 D32 现状：Sleeping 单核 CPU 0.000%、私有工作集 10.60MB、句柄 362、
GDI 10 / User 6、**零活动定时器**（timers=no 是从球对象上实测读出来的，不是查策略表）。

> 说明：`commit` 列里的 privateCommit ≈ 100MB 是 Go+cgo 进程的**提交量**（reserve 未落页），
> 与 D32 的"私有工作集"口径不同；本票一律以 private working set（任务管理器同列）为准。

## 2. 签收前重跑：六个状态行的差分（含"静止"与"贴边"两行零定时器证明）

命令（本文件所有第 2/3/4 节数字都来自这一次跑，输出目录已入库）：

```
build\balldebug.exe -diff docs\evidence\s1\62-diff-signoff -diff-states Sleeping,Listening,Speaking -diff-dock right -x 1720 -y 720 -diff-margin 120 -diff-sample 4s
```

每行三张图：`NN-<状态>-alive.png` / `-dead.png` / `-diff.png`（diff 放大 6 倍）；
数字表：`docs\evidence\s1\62-diff-signoff\diff-table.txt`。
带 `-dock-right` 的行是同一条贴边行为的另一路证据：子进程用与拖拽落点完全相同的
`dockCommit` 路径吸附，父进程照旧"活着截一次 / 正常关掉再截一次"。

| 行 | 变化像素 >=8/255 | 最大差值 | 变化包围盒 | 单核 CPU | 私有工作集 | 活动定时器 | 收尾 |
|---|---|---|---|---|---|---|---|
| **Sleeping（静止玻璃本体）** | **2120** | 171/255 | 44x44 | **0.000%** | **11.76MB** | **no** | clean |
| Sleeping 贴在右缘（标签态） | 1620 | 195/255 | 35x46 | **0.000%** | 12.22MB | **no** | clean |
| Listening（液体在流） | 4856 | 178/255 | 72x72 | 0.023% | 15.00MB | yes | clean |
| Listening 贴在右缘 | 3230 | 209/255 | 47x72 | 0.000% | 14.95MB | yes | clean |
| Speaking（音频驱动旋转） | 4848 | 181/255 | 72x72 | 0.016% | 13.18MB | yes | clean |
| Speaking 贴在右缘 | 3220 | 211/255 | 47x72 | 0.016% | 13.16MB | yes | clean |

**D32 硬门的复算（这是本票不能退让的一条）**：Sleeping 行 **timers=no、单核与全核 CPU 都是
0.000%**，与已入库的 `62-diff-glass` / `62-diff-border` 两行**逐项一致**（变化像素 2120、
包围盒 44x44、句柄 362、GDI 10 / User 6）；私有工作集 11.76MB，票 02 spike 的预算是 ≤25MB。
**没有回归，也没有拿"贴边行"去替换"静止行"。** 新加的"贴边 Sleeping"这一行同样是
timers=no / 0.000% —— 也就是说靠边吸附**一个定时器都没有引入**，吸住的那条依然是一张不动的帧
（1620 个变化像素 = 那个半隐标签本体）。会话态的行有定时器是设计允许的（Animated 白名单），
CPU 最高 0.023%，离 0.5% 的门槛还有 20 倍。

## 3. Summon 液态流动 + 音频驱动旋转

- 白名单在 `internal/ball/liquid.go` 的 `liquidDriven()`：只有 Listening / Thinking / Acting /
  Speaking / Conversation / Warm 六个会话态能被音频驱动，Sleeping 收到电平**直接丢掉**。
  这条门由 winlive 实测断言（`TestBallLiveAudioLiquidGate`：Sleeping + 0.9 电平 → timers 仍为 no；
  Listening + 电平 → 有限 33ms 涌动定时器起来，并且**自己退休**）。
- 音频侧只有**一个标量电平**跨进渲染层（`Ball.SetAudioLevel`，C25 污染面约束：样本与文本不进渲染层）。
  放映里那条"假语音"是 `-level` 提供的：默认 0.70，每第 12 帧压到 1/4，所以你能看出它在跟随，
  而不是一直匀速打转。想自己喂：`build\balldebug.exe -stay -level 0.9`。
- 差分数字：第 2 节 Listening / Speaking 两行（4856 / 4848 个像素变化 >=8/255，包围盒 72x72 =
  整颗球都在动），对照第 1 节基线里"看不见"的那 0 个。

## 4. 静默边框过渡 + 靠边吸附（以及"悬停弹回"这一路怎么被证明）

**静默边框**：`enterState()` 只挪目标值、不挪当前值，所以边框一定是**逐帧滑进来**的；
能带这个过渡的状态就是 `transitionDriven()` 白名单 = 冻结动画策略本来就允许定时器的那几个
（Animated 五个 + Warm 呼吸 + Settling 淡出），**票 62 没有在票 07 留有空timer的地方新加任何一个**，
Sleeping 永远进不了这个白名单。实测记录在 winlive `TestBallLiveIdleBorderTransition`
（例如 Settling 逐帧 0 → 0.11 → 0.32 → 0.53 → 0.74 → 0.96 → 1，之后涌动定时器自己退掉）。

**靠边吸附**：几何在 `internal/ball/dock.go`（纯数学，167→186 行），窗口侧在 `dock_windows.go`。
四条边都按**目标屏的 work area**（不是整屏）计算，所以不会吸到任务栏底下、也不会跑到副屏外；
DPI 缩放时触发距离与边距一起缩放。整条链路**零 SetTimer**：靠近时挤压来自 WM_MOUSEMOVE 的推送，
落点由 WM_LBUTTONUP 的 commit 决定，弹回由 hover 的 WM_MOUSEMOVE 驱动，缩回由 WM_MOUSELEAVE 驱动。

**"悬停把球弹回来"这一路的证明分两层（这是本轮补的关键）：**

1. **确定性层，永不 skip**：`internal/ball/dock_test.go` 的 `TestDockHoverPopBackWalksTheRampHome`。
   它把 ramp 的步进律 `dockRampFrame`（`dockHoverMove`/`dockLeave` 每条鼠标消息真正调用的那一个函数）
   **反着走一遍** —— 从吸住的 p=1 走到 p=0 —— 在 4 条边 × 2 种显示器布局（含原点是负数的副屏）上断言：
   电平单调回退、恰好落在 0（不 overshoot、不停在 0.0001）、走的帧数落在 160ms 预算里、
   每一帧"画出来的球"仍与 work 边相切、球的可见宽度单调变宽直到恢复整颗、窗口沿吸附轴**离开**边缘、
   落地后整颗球完全在 work area 内（可达）、半路缩回时倒着走回标签、以及两个端点**按住不动**
   （落地之后不再产生任何帧 → 不耗电）。变异检验：把步进律改成只会上不会下，这个测试当场 FAIL
   （`primary/left: the ramp froze at p=1 with target=0`）。
   `go test -count=2 ./internal/ball/` 现在跑 **11 个 dock 测试**，全部**不带 skip**。
2. **真事件层**：`go test -tags winlive ./internal/ball/ -run TestBallLiveEdgeDockHover -v`。
   背景进程常常拿不到真指针，而 `TrackMouseEvent` 对**合成**的 hover 会立刻回一个 WM_MOUSELEAVE，
   所以这一层只能证明"线接通了"，不能当作行为的证明。它保留两处 `t.Skipf`（OS 拒绝指针时 skip
   是诚实的），但 skip 文案现在是**响的**：明说本跑没有走过真实事件路径、行为由上面那个
   确定性测试证明。本轮实跑正是走了 skip 分支（`SKIP-LOUD: SetCursorPos refused ...`），
   同一时间 `TestBallLiveEdgeDock`（吸住 + 落地后仍零定时器）**PASS**。

一句白话：**弹回来这件事的逻辑已经钉死在不会跳过的测试里；放映的第 14 步只是让你的真指针再演一遍。**

## 5. 四个 `-look` 候选与推荐

同一套材质（玻璃壳 + 折射 + 边缘高光 + 内唇光 + 接触阴影），只差颜色，`-look` 选一个：

| 名字 | 液体主色 / 副色 / 点缀 | 边缘光 | 读法 |
|---|---|---|---|
| **aurora（默认）** | 电光靛 #4B49FF / 紫 #8B5CF6 / 青 #22D3EE | #141A2E | 蓝紫青渐变那一族，**就是你给的 vivo 蓝心小V 参考图的色系** |
| glacier | 青 #22D3EE / 蓝绿 #2DD4BF / 深蓝 #1D4ED8 | #0E2233 | 更冷更安静，白桌布上对比最低，"看得出来在动但不会抢戏" |
| nebula | 品红 #D946EF / 靛 #6366F1 / 粉 #F472B6 | #1A1030 | 最张扬的一档，小尺寸（44px）下会糊成一团粉 |
| solar | 琥珀 #F59E0B / 珊瑚 #FB7185 / 紫 #7C3AED | #241019 | 把"冷色玻璃"读成暖色的一档；在白桌布上边缘光最吃亏 |

**推荐：aurora 作为默认**（它是最贴近 owner 给的参考图色系的那一个，且白桌布上的可见像素已经在
第 2 节的数字里成立：静止态 2120 个 >=8/255、最大差值 171/255）。
若 owner 觉得 aurora 太"花"，唯一建议改的是 `-look glacier`；nebula / solar 不建议进默认。

## 6. 未验证的部分（说清楚，别让你替我们猜）

- **配色与你口味的贴合度：未验证。** 你给的参考图（图1 / 图2）**从来没有被任何一个跑这票的
  agent 拿到过** —— 没有任何一张图被存进仓库，也没有任何一条路径被交出来。因此本票里所有关于
  "色系"的判断都只是文字描述（"蓝紫青渐变"）的实现，**不是对图的还原**。第 5 节那四个候选
  就是为此准备的：签收时请只判颜色，材质与机制另说。
- 44–72px 里**只有默认 56px 被逐态量过**（第 2 节的 72x72 包围盒就是 56px 球 + 8px 边距 ×1.0 DPI）；
  44px 与 72px 需要你用眼睛过一遍，命令：`build\balldebug.exe -tour -size 44` / `-size 72`。
- 深色桌布上**本轮没有重测**：票面 AC#1 要求"浅色与深色桌布都清晰可见"，本轮机器是白色桌布。
  换深色壁纸后请重跑第 2 节那条 `-diff` 命令（或 `-tour` 目测），这是签收前唯一还缺的一项视觉证据。
- CPU / 内存的**长时间**表现（例如挂一整天）不在本票口径内；本票的数字是每条 4 秒采样窗口的实测。

## 7. 待签字的 SPEC-08 §2 措辞（owner 签收后才回填契约；本轮**未**改 SPEC-08 / PLAN）

建议把 §2 的"20 态视觉表"从「圆点 + 状态色环」改写成下面这套口径（状态数、事件、转移**都不动**，
只动视觉描述），并把上面第 2 节的数字作为门限写进验收：

> **§2.1 悬浮球视觉 = 玻璃液态球（票 62）**
> 1. 本体是一颗直径 = 配置值（44–72px，默认 56）的**透明玻璃球**：外壳折射 + 边缘高光 +
>    内唇光 + 接触阴影，内部为可流动的液态色团（配色档 `-look` 四选一，默认 `aurora`）。
>    球体在浅色与深色桌布上均须可辨：与"进程不存在"的同区域差分比较，`>=8/255` 的变化像素数
>    在 Sleeping 态不得低于 **1500**（实测 2120，包围盒不含透明边距）。
> 2. **静止态（Sleeping / Armed / Muted / Error 等非会话态）是一帧静态图：零活动定时器、
>    零动画、CPU 采样 0.000%（D32 门不变，≤0.5% 为上限）。** 为可见性所致的任何加强
>    只能走"提高本体不透明度 / 加深色边缘光"，禁止给静止态添加定时器。
> 3. **会话态（Listening / Thinking / Acting / Speaking / Conversation / Warm）**内部液体流动；
>    其中 `liquidDriven()` 白名单内、且采集到语音时，液面**随音频 RMS 包络旋转/起伏**（单调，
>    电平上限 1.0，帧率上限 30fps）。仅允许**有限时**的涌动定时器（33ms，运动停了就自退休）。
> 4. **静音但会话仍在进行**时，液面收敛、**外圈边框以一个过渡动画（≥3 帧）滑入**，禁止单帧跳变；
>    只有会话态可以带这个过渡定时器（= 票 07 的冻结动画策略本来就允许的那些状态）。
> 5. **靠边吸附**：拖近屏幕上/下/左/右任一**工作区**边缘（默认触发距离 16px，按 DPI 缩放）时，
>    球沿该轴收缩成贴边标签（可见比例 0.42），仍保留一条可命中区；**指针悬停或单击 → 弹回完整球**
>    （单击后保持弹出，直到下一次拖拽），指针离开 → 收回标签。整个过程**不得引入任何常驻定时器**
>    （实现只能由鼠标消息 / 窗口位置事件驱动，实测吸住的 Sleeping 行 timers=no）。
>    吸附计算必须使用目标显示器的工作区，禁止吸到任务栏之下或副屏之外。
> 6. 音频数据与任何文本内容不得进入渲染层（C25）；渲染侧只接受一个电平标量。

同时建议登记两条**残缺表现**（推迟项，不许变成无声技术债）：
(i) 深色桌布的差分复测；(ii) 44/72px 两档的逐态截屏。完成判据即第 2 节那条 `-diff` 命令换参数重跑。

---

## 8. 复现这些数字的命令

```
cd /d "D:\work\workspace\projects plans\Wisp"
set PATH=%CD%\third_party\sherpa-onnx;%PATH%
go build -o build\balldebug.exe .\cmd\balldebug
build\balldebug.exe -tour
build\balldebug.exe -diff docs\evidence\s1\62-diff-signoff -diff-states Sleeping,Listening,Speaking -diff-dock right -x 1720 -y 720 -diff-margin 120 -diff-sample 4s
go test -count=2 ./internal/ball/
go test -tags winlive ./internal/ball/ -run TestBallLiveEdgeDock -v
```

（跑 `-diff` 之前桌面上不要留任何一颗球：残留的演示进程会把"死帧"污染成"活帧"，本仓已因此
产生过一次假测量。`-diff` 自己只拉起/回收它自己的子进程，收尾会打 `clean`。）
