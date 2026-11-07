# 239-verdict — 票 239 五格的独立裁决（wave50 裁决腿，票 239 验收）

- **裁决者**：wave50 腿（与实现者编排者、与 239-v1 腿互为他人，`AGENTS.md` §0.3）。
- **裁决时锚点**：主树 HEAD `766bf71f4`（票 170），`git status --porcelain internal/ball`＝**空**（进场时与全程收尾均空）。239 落地 commit＝`0589fd9cc`。
- **关键事实：HEAD 不是 239 时代的 HEAD**。`0589fd9c..766bf71` 之间，`internal/ball` 又进了票 65（jarvis 渲染器／玻璃）、68-r1/r2（prototypeVisuals 默认、**泊靠几何改走态半径**、Settling 0.95）、244/245/344 等。本表全部读数取自**当前 HEAD**，与票面/239-v1 的读数漂移逐一归因到那次后续改动，不把它们错记到 239 头上。
- **突发与处置（两发桌面污染）**：本腿两轮桌差分（11:06–11:08）全框读数不可复现（Sleeping px>=8 两发 9586 vs 18645），判**桌面被其他窗口刷新污染**（239-v1 §1 也具名过同类现象"Speaking 那一发被屏幕上别的窗口刷新洗掉了"）。处置：按**窗口矩形**（margin=120、edge=100 ⇒ 框内 [120,220)²）重算 alive/dead 原始逐像素差分，剥离窗外污染；口径先对全框验证与工具逐字同（27342/18645/7728、6014/5339/3526 全中）再取窗内。**窗内读数完全可复现**（clean 两发逐字同）。
- **本腿一处读数错误自纠**：此前的"winlive 全绿 ok"读数作废——`rc` 取的是管道尾命令的退出码。改用输出落文件取真 rc＝**1**，红名册两枚（见 AC#1）。

## 0. 逐格判语总表

| 格 | 判语 | 一句话理由 |
|---|---|---|
| AC#1 | **部分** | 四把尺全过；红名册≠"只有 CSS 那一枚"（CSS 红已消失，但另有两枚 winlive 基线红，均归后续票）；两形差分做了，Sleeping 全框读数被桌面污染不可复现，窗内口径可复现但与家族对不上 |
| AC#2 | **部分** | 耦合仪器有牙（16/21 双红、红句=耦合成因句，工作树全程零字节）；但点名用例 `TestBallLiveEdgeDock` 的**未变异基线在 HEAD 已经是红的**（归票 68-r2），"复绿"对这枚用例不成立 |
| AC#3 | **PASS** | 恒等路径返回值变异（fill/grad→1.4R/1.5R）⇒ 窗内 p8 **1811→1047**（−42%）、box 52x48→38x36，离族 1051 枚，clean 两发噪声 0——这句话有牙 |
| AC#4 | **PASS** | ①②③ 三半全核过：22px 间隙有机主「22可以」批注在场；0.42 一字节没动（diff 零命中）；分母 115600 本腿现跑复算=340² 坐实 |
| AC#5 | **PASS** | 3 枚文件、46±12 与票面逐字同；表↔码同音本腿独立核（变异下 geometry 红句未复现、保持台件口径）；**归口欠账形状已变**：68-code-verdict.md（2026-11-06）已在，且**坐实了那两枚基线红是 68-r2 未跑 winlive 的未爆雷** |
| 总判 | **5 格维持已勾，不退回** | 239 本票的改动在 HEAD 上站得住；新增两枚基线红是 68/245 地界的债，不翻 239 的框 |

## 1. AC#1 — **部分**（四把尺过；红名册偏离；Sleeping 差分受污染）

**四把尺（全〔本腿现跑〕，11:03–11:06）**：

| 尺 | 读数 |
|---|---|
| `go build ./...` | **rc=1**，唯 `build/mut_rvw.go` 十余枚 `undefined:`——该文件**未追踪、不在版本库**（`git ls-files build/` 空、Nov 4 时戳），是别的腿留在 build/（整目录 `.gitignore:14`）的外来件。**收窄面复跑**（`go list ./...` 42 包剔 1 后全构）＝**rc=0**。判：仓库状态 rc=0，偏离具名 |
| `gofumpt -l`（三枚文件） | 零输出 rc=0 |
| `go vet ./internal/ball/ ./cmd/balldebug/` | 空 rc=0 |
| `sh scripts/d22scan.sh` | **clean rc=0**；阳性对照 45 PASS/0 FAIL；real repo 432 产码文件 |

**红名册（〔本腿现跑〕）**：
- `go test ./internal/ball/ -count=1` ⇒ **ok 0.210s，零红**。⚠ 239 时代那枚 `TestC21TableColourRowsMatchTokensCSS` 红**已不在场**：`design/assets/tokens.css` 文件仍不在盘上（`ls` 失败），但测试现读数 `78/78 cited rows agree … :root=129 decls, light=65`——仪器读的路径回来了（归 344/后续票地界，本腿未深挖）。
- `go test ./internal/ball/ -tags winlive -count=1` ⇒ **真 rc=1，两枚红**（输出 `probes/239-verdict-w50/winline-full.log`）：
  1. `TestBallLiveEdgeDock`：`live_windows_test.go:384: edge 1: docked orb boundary -4, want the work edge 0`——**归票 68-r2**（4592c76a5 把 `dockGeometry` 的 `orbR` 从窗口半径改为态半径 `dockOrbRadius(SizePx)=int32(34.72/2)=17`（dock_windows.go:69 现读；Sleeping 静息体 34.72=56×SleepRestRatio），测试公式仍按窗口半径 28（=edgePx/2−margin=50−22）算 `int(28×0.42)=11`，而泊靠律实际按 `int(17×0.42)=7` 摆窗，差恰 4 ⇒ 红句的 `-4`）。68-r2 交付表 §3 具名承认 winlive 全套未跑。该测试的边界断言块自 0589fd9c 以来零改动（diff 该块零命中）＝旧公式对新几何。
  2. `TestLiveConfirmingCancelAndEscReturned`：`interaction_live_test.go:239: the taken-over Esc did not cancel Confirming`，独立复跑两发同句确定红——归 245 地界（该文件最后改动＝245-r1 `e6d6426f8`；这是物理注入用例，桌面状态可能参与，**真机归属留给 245 腿**）。
- **偏离票面判据原文**：票面要求红名册＝"只有 CSS 那一枚"；当前 HEAD 实况＝**CSS 红消失＋winlive 两枚新红**。两枚都不由 239 的三枚文件机制造成（归因链如上），故判**部分**而非 FAIL。

**两形差分（〔本腿现跑〕，自建 `.scratch/wisp/probes/239-verdict-w50/bin/balldebug-w50.exe`＝HEAD `go build ./cmd/balldebug`）**：

| 发 | 口径 | px>=3 | px>=8 | px>=24 | box | 分母 | 可复现性 |
|---|---|---|---|---|---|---|---|
| Sleeping（11:06 run1） | 全框（工具） | 9942 | 9586 | 7717 | (150,0) 190x188 | 115600 | 不可复现 |
| Sleeping（11:08 run2） | 全框（工具） | 27342 | 18645 | 7728 | (0,97) 188x243 | 115600 | 不可复现 |
| Sleeping run1/run2 | **窗内 [120,220)²** | — | **1811/1811** | **953/953** | **(144,148) 52x48** | 10000 | 逐字同 |
| m3 突变（见 AC#3） | 窗内 | — | 1047 | 750 | 38x36 | 10000 | — |
| Listening（11:06） | 全框＝窗内（box 在窗内） | 6014 | **5339** | 3526 | (129,129) 82x82 | 115600 | — |

- **Sleeping 与家族（2098/2103/2120）对不上**：窗内 1811 比家族下沿 2098 低 287（−13.7%），但 noise floor＝0（clean 两发逐字同）——不是噪声，是**几何/渲染漂移**。候因：68-r2 的 drawGlass/停靠改动、65 的玻璃重工（HEAD `02182e1e8` 玻璃 v5 在差分两轮里未编译进我自建二进制？**我判定没有**：我的 balldebug-w50 是 11:06 自建、晚于 02182e1e8，读数即 HEAD 形态）。**判不了**：Sleeping 1811 到底是"239 之后谁把它挪了"还是"家族读数本来就带窗边切割成分"——需要一次"还原旧参数再拍"的对照（拿 62-diff-signoff 时代几何在旧窗口拍），本腿不越权翻旧档；具名留给编排者。
- **Listening 全框读数即窗内读数**（box 82x82 完全在窗内）＝**无方形切割**：box 82 < 窗口 100，外沿干净。`-diff-gate` PASS（min px_delta_ge8=800 铺位过）。**"球被窗口切成方形"这一形在 Listening 上不存在**；Sleeping 的全框污染也不呈现方形（box 触框缘是污染形状，窗内 52x48 居中）。
- gate 铺位源注记照录：`floor source: SPEC-08 §2 (rest body 56*SleepRestRatio=34.72px, 1.5R halo => 2130 px predicted, 2098/2103/2120 recorded; SleepingRestMinPx=30 => 1590 px floor)`——窗内 1811 在铺位 1590 之上。

## 2. AC#2 — **部分**（仪器有牙；但基线 Dock 用例已被 68-r2 打断，"复绿"不成立）

**方法**：不走工作树编辑，走 `go test -overlay`（239-v1 同法）——变异件各与盘上 tokens.go 差**恰好 1 行**（`diff | grep -cE '^[<>]'`＝2，即 1 删 1 增；LF 保留），`RingMarginPx` 一枚未动。全程 `git status --porcelain internal/ball`＝空（插桩前后自证），收尾工作树零字节。**效果等价于"改回 16"且收尾零字节**，还原核验＝`git status --porcelain internal/ball` 空＋撤 overlay 后复跑基线回到 HEAD 真实状态（A→B→A 闭环）。

**A→B→A 读数（〔本腿现跑〕，11:14–11:19）**：

| 发 | overlay | `TestBallLiveEdgeDock` | `TestBallLiveEdgeDockHover` | 红句 |
|---|---|---|---|---|
| 基线 A（HEAD，24） | — | **FAIL**（0.38s） | PASS（0.26s） | `:384 docked orb boundary -4, want the work edge 0` |
| m1（16） | ✓ | FAIL（1.04s） | FAIL（0.09s） | `:354 DebugDock did not commit the dock`＋`:521 would not dock on the right edge` |
| m2（21） | overlay | FAIL | FAIL | 同上两枚耦合成因句 |
| 基线 A′（撤 overlay） | — | FAIL（0.45s） | PASS | `:384 boundary -4`（回到 HEAD 真实状态） |

- **耦合有牙**：16 与 21 两发都让**两枚点名用例红**，红句与 239-v1 记录**逐字同**（`did not commit the dock`／`would not dock`）；21（≠16）也红＝咬的是"trigger<margin"这条律本身，不是旧值回归。Hover 用例基线绿→变异红，**干净的翻转证据**。
- **基线 Dock 红的归因（本腿算术现算）**：68-r2 后 `orbR=dockOrbRadius(SizePx)=int32(34.72/2)=17`（dock_windows.go:69），测试仍按窗口半径 28 算期望 `int(28×0.42)=11`，泊靠律实际摆窗按 `int(17×0.42)=7`，边界差 11−7=4 ⇒ `boundary -4`。**差异 Δ=4 逐字等于红句的 -4**，机制闭合。
- **"测完还原并复绿"这半句的判语**：Hover 用例**成立**（绿→红→绿闭环）；Dock 用例**在 HEAD 上无从复绿**（基线即红），本腿判**部分**——票面预想的"变异红→还原绿"形状对这枚用例今天是**绿→更红→红**（同一用例上叠加了耦合成因句，其余三边几何检查全过——红句是 `:384` 边界句而非 `:354` 提交句，说明 Dock 提交、D32 零定时器、防藏没挂）。
- **给编排者的口信**：两枚基线红都是 68/245 地界的债——**不要把它们记到 239 头上，也不要借 239 的手去修 68/245 的债**。

## 3. AC#3 — **PASS**（"对 Sleeping 恒等"这句话有牙）

**方法**：m3＝盘上 renderer_windows.go **纯插入 3 行**（diff 行数=3，插入点＝原两层 clamp 之后、return 之前）：
`if R*1.9 <= float32(r.w)/2-1 && R*1.9 <= float32(r.h)/2-1 { fill, grad = R*1.4, R*1.5 }`。

我插的是 `R*1.4, R*1.5`——比 239-v1 的 `R*1.9, R*1.1` 温和，**变异方向一致（恒等路径上的两返回值不再等于原函数返回值）**。语义：fill=1.4R＜原 1.5R、grad=1.5R＜原 1.9R，fill<grad 仍成立、恒等路径条件与原代码一致、Sleeping 的 1.9R=32.98≤fit=49 为真。⚠ 对 239-v1 那种"两值互换"变异：票面字面是"fill 与 grad 两个返回值…强行互换"，严格互换 fill=1.9R>grad=1.5R 违反原函数不变式 `fill ≤ grad`，可能触发下游 clamp 分支差异；**我选择温和变异以隔离"返回值变了"这一个因子**。等价性陈述：任何能改变返回值的变异若读数不动才算"没牙"，我的变异改变了返回值且读数动了 ⇒ 有牙判定成立。

**读数（窗内口径，同 AC#1）**：

| 发 | p8 | p24 | box |
|---|---|---|---|
| clean（两发逐字同） | **1811** | 953 | 52x48 |
| m3（恒等路径返回值改变） | **1047** | 750 | 38x36 |
| Δ | **−764（−42%）** | −203 | 收缩 52x48→38x36 |

离族 1051 枚（=家族下沿 2098−1047，−50.1%）——**"haloExtents 对 Sleeping 是恒等变换"由可观测后果撑着，不是装饰**。比 239-v1 的离族 404（−19.3%）更硬，因为我的口径剥离了污染，noise floor=0。

**还原核验**：m3 只存在于 overlay JSON＋probes 目录，工作树零字节（`git status --porcelain internal/ball` 空）；m3 构建无子进程污染。**PASS**。

## 4. AC#4 — **PASS**（①②③ 全核过）

- **① 22px 贴边间隙**：〔本腿现读码〕`hit.go:31 const RingMarginPx = 22.0`＋`:28` 注释逐字"**22 is the smallest value that fits the waves at every size in 44..72**"在场；22 最小性推导＝波纹外沿 `R+5s+t*16s+RingStrokePx/2·s = R+21.75s`（renderer_windows.go:503 `r.circle(c, R+5*s+float32(t)*16*s, RingStrokePx*s)` 现读）→ 22 是 44..72 全档最小整数边距。机主批注「22可以」已写进票面（2026-10-22 11:0x 编排者翻框注记，原话逐字在场）⇒ **已裁**。
- **② `DockOverlapFrac` 一字节没动**：`git diff d5a59d66..0589fd9c -- internal/ball/` 里 `DockOverlapFrac|ClickTolerancePx` 赋值行 grep＝**0**（本腿现跑）。⚠ **但**：票 68-r2（4592c76a5）之后 `DockSquash(1)` 的**消费面**变了（orbR→态半径 ⇒ 窗口可见伸出量变化），本腿判：239 的"没动"依然逐字真，68-r2 对 0.42 的改法在 68-code-verdict.md 里已被非实现者裁决过 ⇒ ② 的"没顶掉已被接受的数"依然成立。
- **③ 分母**：本腿窗内重算的分母 10000=100²（差分图 340x340、margin 120、窗口 100）⇒ 分母=（窗口+2×margin）²=340²=115600 与 239-v1 结论一致；diff 的 `of 115600` 在工具输出里直接可见。**"分母不变"那句错账在 239-v1 已更正，本腿现跑同向坐实**。
- **残余（票面已登记，本腿现核仍在）**：裸字面量 `16` 走距仍在 `renderer_windows.go:503`（不在 C21 名册）；0.25px 富余（满档 72 时 R+21.75 vs R+22）无测试钉。这两枚**归口本票的残余条款，未清**，但票面已具名未结案、未加 `-done` ⇒ 不阻碍本表判语。

## 4′. 残余两枚的"是否真的还没修"核查（本腿现跑）

- **裸 16 仍在**：`renderer_windows.go:503` 现读 `r.circle(c, R+5*s+float32(t)*16*s, RingStrokePx*s)`——**未收进 C21 名册**（c21-native-tokens.md 无对应行）。
- **0.25px 富余**：满档 72 时波纹外沿 57.75 vs 窗口半 58，余量 0.25px；无差分用例钉"72 档 Listening box < 窗口 edge"这一形（本腿 AC#1 只拍了默认档 Listening，满档未拍——**判不了满档是否仍不被切**，具体留编排者）。
- **cmd/wisp 全量**：票面自述那发 `ok 112.312s` 是**改前读数**。本腿在 HEAD 没跑（耗时 110s+，不属本表五格判据；但 NOTE：`go build` 收窄面 rc=0 已把 cmd/wisp 的编译面盖了）。**判不了**（要编排者下一枚整包门禁）。

## 5. AC#5 — **PASS**（3 文件/46±12/表码同音/68 欠账形状已变）

- **stat**：`git show 0589fd9c --stat -- internal/ball/`＝`hit.go 16±`、`renderer_windows.go 35±`、`tokens.go 7±`（3 枚、46 增 12 删）——与票面/239-v1 逐字同。
- **表↔码同音（本腿独立核）**：`c21-native-tokens.md:255` 现读 `环边距 22px / 点击容差 4px`、`:257` `距边 24px 内开始挤压`，与盘上码 `RingMarginPx=22.0`／`DockTriggerPx=24` 同音；且 AC#2 基线（24）下非 live 全量 **ok 0.210s 零红**（含 `TestC21GeometryRowsMatchCodeConstants`）⇒ 表↔码目前同音。⚠ 239-v1 §2 那句"变异下 geometry 那枚红句"的读数，本腿**未复现**（本腿 m1 是点名跑、未复跑全量非 live 套件），保持〔读台件，未复跑〕口径。
- **归口欠账形状已变**：`docs/evidence/s1/68-code-verdict.md`（2026-11-06，非实现者 W68b 写）**已在**——239-v1 时代的"68 名下零份非实现者验收表"欠账**已清**。68-r2 的交付表具名承认 winlive 未跑；**68-code-verdict.md 的门禁读数也没跑 winlive**（其 §1 只列 gofumpt/vet/build/test(`-short`)/d22scan 五门，winlive 缺席）——**本腿 AC#1 的两枚基线红就是这笔未爆雷的现形**。建议编排者把两枚红登 68/245 名下，不是 239 的。

## 6. 本腿判不了／留给编排者（⛔ 不许为空）

1. **Sleeping 窗内 1811 vs 家族 2098 的 −287 差**：候选因（68-r2 drawGlass、65 玻璃 v5、家族读数本身带旧几何成分）**判不了**——需要"还原旧参数再拍"的对照。但 noise floor=0、窗内可复现、`-diff-gate` PASS（铺位 1590 之上），**不阻碍形状判语**（Listening box 82<100 无切割）。
2. **TestLiveConfirmingCancelAndEscReturned**（interaction_live_test.go:239）：245 地界的物理注入用例，桌面状态可能参与；两发同句确定红。**本腿判不了根因**（要 245 腿或真机），**判得了"不是 239 的"**（该文件 239 未动、245-r1 是最后改动者）。
3. **满档 72 的 Listening 差分**：没拍。**判不了满档是否仍不被切**——票面残余条款 1 的完成判据之一，留给编排者/下一枚腿。
4. **cmd/wisp 全量**：没跑（110s+，不属五格判据）。**判不了**，留给编排者整包门禁。
5. **票面 AC#4 机主「22可以」批注**：我核的是票面自述（编排者翻框注记 2026-10-22 11:0x＋原话逐字在场）——**我判不了机主本人是否真的说过**（要机主本人或会话记录）。采信＝票面注记，不采信＝要机主重新裁。
6. **AC#2 基线 Dock 红的最终归属**：归因算术闭合（Δ=4），但**"68-r2 改泊靠几何后该测试该不该重写公式"是设计决定**——公式改不改、改成什么样，归编排者/68 票面，本腿不代判。
7. **239-v1 §2 "geometry 红句"读数**：本腿没复跑全量非 live 变异套件，239-v1 那句保持〔读台件，未复跑〕。

## 7. 收尾自证（⛔ 纪律）

- `git status --porcelain internal/ball` 进场/收尾均**空**（11:19/11:2x 现跑）。
- 变异全走 overlay，工作树零字节；差异文件全在 `.scratch/wisp/probes/239-verdict-w50/**`（临时件只建不删）。
- 本腿写面=本表一枚文件。balldebug 进程 0（每发前后 tasklist 自证）。
- 桌面差分两轮受污染的两发全框读数（9586/18645）**作废**，以窗内口径为准。
- ⛔ 本表不勾票面任何复选框（框归编排者）；本票不加 `-done`（残余两枚未清＋基线两枚新红待归口）。
