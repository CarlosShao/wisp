# 68 — 让默认构建画的就是 owner 签收的那个球（`prototypeVisuals` 默认值 + Sleeping 尺寸三方不一致）

**Status:** review（AC#1 静态三列表 + AC#4 已交付并由编排者复现勾框；**AC#2/AC#3 blocked-on-owner**，见 R15 #3/#4/#5）
**Claimed by:** agent-ticket68（报告已交；**勿重开 AC#1/AC#4**；owner 答完 R15#3/#4/#5 前**不要翻默认值**，
否则会把"我录入错误的尺寸"或"合成到 ≈0.33 的 `Settling`"当成签收样子固化进默认构建）
**Last update:** 2026-09-20
**Blocked by:** —（只碰 `internal/ball/` + `cmd/balldebug/`；票 66 在 `internal/proc`/`internal/observe`/`cmd/wisp`，票 67 在 `internal/llm/adaptertest`/`tools/d22scan`，三包不相交）
**Parallel slots:** ≤1 sub-agent；**AC#2/AC#3 需真桌面** ⇒ 若桌面被占，**只做 AC#1 与 AC#4 的静态半，其余保持未勾并写明**
**Spec refs:** SPEC-08 §2 / §2.1（**冻结，本票不得编辑**；含 2026-09-20 的 INTERIM 标记）、D29、D32（不放宽）
**登记项:** 票 12 AC#7 代理报出的争议 **D1 / D2**（registry 待并入 A24）

## 事实（我自己读码确认，非转述）

`internal/ball/statevisual.go:76-84`：

```go
// prototypeVisuals selects the ticket 62 liquid-glass rendering. It is OFF by
// default: SPEC-08 §2.1 (12px/0.35 Sleeping dot etc.) is frozen until the
// owner signs the new look, and the ball's own tests assert that table.
var prototypeVisuals bool
func EnablePrototypeVisuals(on bool) { prototypeVisuals = on }
```

唯一把它打开的生产侧调用者：`cmd/balldebug/main.go:104` → `ball.EnablePrototypeVisuals(!*frozen)`。
被它门控的代码面：`dock_windows.go:90/114/176/206/229/245`、`liquid.go:313`、`liquid_windows.go:52/154`（票 62/64 的吸附与液面全在这里）。

⇒ **owner 2026-09-20 用肉眼签收的那个 44px 静态玻璃体，只在 `balldebug` 里存在；库默认画的是旧的 12px 微点。**
而我今天把 SPEC-08 §2 的 `Sleeping` 行改成了"44px 静态玻璃体"——**契约文本现在描述的是一个非默认配置**。
这个不一致是**我造成的**，本票就是还这笔账。

⚠ 注释里那句"until the owner signs the new look"的**前提今天部分成立**（R13：临时通过、先完成核心功能），
所以翻转默认值是**与已改契约一致**的动作，不是擅自扩大签收范围：**owner 批的只有"尺寸 / 可见 / 零定时器"三项，
质感仍是"赝品"**（见 SPEC-08 §2 INTERIM 与票 65）。

## 第二处不一致（D1，比默认值更要紧）

同一态下三个数对不上：SPEC-08 §2 INTERIM 写 **44px**；
`statevisual.go::stateSize` 在默认 56 基准下按 `SleepRestRatio = 0.62` 算出 **34.72px**（下限 `SleepingRestMinPx = 30`）；
而 `docs/evidence/s1/c21-native-tokens.md` 的旧行写 **12px**。
但 `docs/SLO.md` 附录 A.2 实测到的差分化像框是 **46×46、2103 像素变化 ≥8/255** ⇒ 那个跑法里 Sleeping **确实是 44px 级**。
**所以要先把"到底是哪个尺寸在什么条件下成立"查清楚**，再决定改码还是改文档——**不许**为了让三方一致就去编辑 SPEC-08。

## 验收标准

- [x] **AC#1 尺寸真相**：以代码 + 一次实测（若桌面可得）钉死 `Sleeping` 在
  ①`prototypeVisuals=false`、②`=true` 且未靠边、③`=true` 且已吸附（dock ramp 之后）三种情形下的**实际像素尺寸**，
  逐情形给出 `stateSize` 的输入与输出、以及差分化像框。**完成判据是一张三列对照表**，不是一个结论句。
  与 SPEC-08 §2 的 44px 不符的那一格**如实报不符**，由我裁定。
  ⇒ **静态半已交付**（三列表 + 逐格函数/输入/输出/消费者，见 Progress log 与 agent 报告；桌面不可得，未做新实测）。
  **三格全部与 44px 不符，待裁定**：①12px、②34.72px、③34.72px 本体（靠边只改窗口原点与液面/高光的 0.42 挤压，球壳与光晕不挤压）。
  A.2 的 `46×46 / 2103px` 经码算归属于 **②**，且 `44` 那个数是 ≥24/255 的**包围盒**、不是本体直径。
  —— **编排者 23:16 勾框（我自己复现过，非采信自述）**：`go test ./internal/ball/ -run TestSleepingSizeTruthTable`
  → PASS，并读了它的断言体（`tokens_test.go:214` `rest56 = BallSizeDefaultPx × SleepRestRatio // 34.72`、
  `:218-220` 配置 44/48 落到 `SleepingRestMinPx` 下限）。**勾的是"静态三列表"这个交付物**；
  像素级新实测**没有**做，也不得被这张勾解释成做过。
  ⇒ **裁定已下（R15#1）**：SPEC-08 §2 的"44px"是我录入错误，已按 34.72px 更正（原文与推导留在更正块里）。
- [ ] **AC#2 默认值翻转（需桌面复测）**：`prototypeVisuals` 默认改为 `true`，
  使**默认构建 == owner 签收的样子**；`EnablePrototypeVisuals(false)` 与 `balldebug -frozen` 保留为"对照旧冻结规格"的逃生门。
  连带把 `internal/ball` 里**断言旧默认**的测试改到位：⚠ **不得**为了让测试变绿而删除断言或放宽阈值——
  必须**逐条**把断言迁移到"新默认下应有的行为"，并在 log 里列出"这条原来断什么、现在断什么、为什么等价"。
  若某条测试的存在意义就是"钉住旧的冻结规格"，**保留它**并显式 `EnablePrototypeVisuals(false)`，别删。
  - ⚠ **动手前必改的一条地雷（代理登记在 Progress log，这里提到 AC 正文以免漏）**：
    `internal/ball/live_windows_test.go:609` `TestBallLiveAudioLiquidGate` 的第一段
    以"默认即 frozen"为前提（`:618` 注释"模式冻结优先"）**却没有显式 `EnablePrototypeVisuals(false)`**
    ⇒ 默认一翻，它会**自称在测冻结、实际在测原型**。翻转与这条修正在**同一个 commit**，不许分两批。
  - ⛔ **AC#2 被 R15 三项挡住，未答之前不得翻转**（翻了就是照错的数字固化）：
    **#3** `Settling` 的 alpha 今天双乘到 ≈0.33、**#4** 停靠实测露出 ≈82% 而文档承诺 42%、
    **#5** 命中半径 21px < 可见外沿 26px。这三条都改变"被签收的样子"到底是什么，
    owner 未答前我只推进**不依赖它们**的部分（AC#1 已完成、AC#4 已绿）。
- [ ] **AC#3 签收面复测（需桌面）**：翻转后跑一次 `cmd/balldebug` 的差分化像，
  证明 `Sleeping` 的 px≥8/255 与成像框**不低于** A.2 已录的 2103 像素 / 46×46，且 `timers=no` 仍成立（D32 零定时器）。
  - **桌面空出后照此逐条跑（代理 AC#1 报告给的、与 A.2 同口径的命令集，落盘以免随上下文丢失）**：
    1. `export PATH="$PWD/third_party/sherpa-onnx:$PATH"`
    2. 洁净前置：`tasklist | grep -iE "wisp.exe|balldebug.exe"` **必须无输出**；
       并确认显示器缩放是 **100%**（A.2 全部数字都是 96 DPI 下测的，否则不可比——见 A28/A29 的 DPI 不对称）
    3. `go build -o build/balldebug.exe ./cmd/balldebug`
    4. 自由 + 靠边两行：`build/balldebug.exe -diff build/t68-flip -diff-states Sleeping,Listening,Thinking,Acting,Speaking,Warm,Settling -diff-sample 5s -diff-dock right`
       ⇒ 门：`Sleeping` **px≥8 ≥ 2098**、框落在 40–46 家族、`timers=no`、`cpu_all ≤0.5%`、teardown `clean`
    5. **冻结对照（证明逃生门还活着）**：同上一条加 `-frozen` ⇒ 期望 **px≥8 = 0**、框 `0x0`
    6. 零定时器与句柄门：`build/balldebug.exe -state Sleeping -hold`（Ctrl+C；打印 `timers=false` 与 `closed handles=` ≤600）
    7. `go test -tags winlive -count=2 -v ./internal/ball/`
    8. 翻转后重跑 AC#4 四道门
  - ⚠ 第 5 条（`-frozen` 对照）**不是可选项**：逃生门若坏掉而没人知道，票 65 返工时就没有基线可比。
  - 数字不过预算就**报 FAIL 并附全样本**；不许调阈值、不许重测到运气好的那次、不许用平均抹掉坏尾部。
  数字不过就报 FAIL 并附样本，**不许调阈值、不许重测到运气好的那次**。
- [x] **AC#4 无桌面的半格**：`gofmt -l` 空、`go vet ./internal/ball/ ./cmd/balldebug/`、
  `go test -count=2 ./internal/ball/`、`go test -count=2 ./cmd/balldebug/`（非 winlive 部分）。贴原始输出。

## 编排者已裁定

1. **SPEC-08 与一切 `docs/specs/*` 冻结**：本票**零编辑**契约文件。尺寸若与契约不符，**改代码或报不符**，不许改文本。
2. **D32 的零定时器与 ≤0.5% CPU 不因"只是默认值"而放松**。
3. 若桌面被票 66/67 占着：AC#2/AC#3 **保持未勾**并写明"等桌面"，**先交付 AC#1 + AC#4**，不要为了推进而偷跑测量。

## Progress log（append-only）

- [2026-09-20T14:49:27Z] agent=agent-ticket68 did=claim 票 68（Status in-progress / Claimed by）。桌面被票 66（`internal/proc`+`internal/observe`+`cmd/wisp`+`cmd/balldebug/diff_windows.go` 的 proc 迁移，工作区里正有未提交改动）与票 67 占用，按裁定 3 只做 AC#1 静态半 + AC#4；AC#2/AC#3 未跑、未勾、未改默认值。
- [2026-09-20T14:59:41Z] agent=agent-ticket68 did=AC#1 三列尺寸真相表（纯码 + 已录证据，**零新实测**）+ AC#4 四道门全绿 + 码内纠偏（零行为改动）：`statevisual.go` 三处注释把「12px 仍是冻结契约」「SizePx 稍后按 DPI 缩放」两个假前提改成事实——渲染目标 dpiX/dpiY=96，`drawFrame` 直接用 `R=SizePx/2`，只有描边/环边距走 `dpi/96`、窗口边走 `WindowEdgePx`，所以 144 DPI 下 56 配置会画成「108px 窗口里的 56px 本体」。新增两个纯函数测试：`tokens_test.go::TestSleepingSizeTruthTable`（①12 ②34.72 ③34.72；`SleepingRestMinPx` floor 在配置 44..48 生效、49 起失效；`DockSquash(0/1)=1/0.42`、落点窗口只剩 47/72px 在屏）与 `TestRecordedSleepingDiffBoxIsNotA44pxBody`（π(1.5R)² 反证：34.72px 预测 2130 ≈ 已录 2120/2103，44px 预测 3421，故 A.2 的 46×46 是 ② 那格的 ≥24/255 **包围盒**而不是本体直径；此测试**故意**在有人改 `SleepRestRatio` 而不重测时变红）。同时在 `TestVisualForCoversAllTwentyStates` 里点名 AC#2 要逐条迁移的断言（Sleeping 尺寸/opacity、Armed 0.6、Muted 0.4、Settling 末端 0.35），一条未改、一条未删。**未跑**：AC#2 的默认值翻转与复测、AC#3 差分复测、任何 `-tags winlive` 构建、`build/balldebug.exe`/`build/wisp.exe` 的编译或启动、建窗/点击/热键注入/截屏——桌面仍被票 66/67 占着（本轮工作区里能看到 `docs/evidence/s1/66/` 与 `slo-check` 的活跃改动），等桌面。翻转本身只有 1 行，但它会让默认构建成 AC#1 判定为「与 SPEC-08 §2 的 44px 不符」的那一格（34.72px），未拿到裁定就翻等于把不一致写进默认路径，故停在报告里。
- [2026-09-20T15:04:14Z] agent=agent-ticket68 did=把 AC#1 的反证测试从"单次实测"升级成**三次实测**（`TestRecordedSleepingDiffBoxIsNotA44pxBody` 现在遍历 2098/2103/2120，容忍度 2% 未动）：新引 `docs/evidence/s1/62-diff-glass/diff-table.txt`（2098px，像框 40×38）与 `docs/evidence/s1/62-diff-baseline/diff-table.txt`（**列①冻结微点：px≥8=0、px≥3=17、max=4/255、无像框**）。三跑实测的 ≥8 计数稳在 2098–2120（34.72px 本体预测 2130，误差 0.48–1.54%），而 ≥24 的"像框"在**同一帧**上 40×38 / 44×44 / 46×46 抖了 6px ⇒ 44 那个数不是本体直径，是包围盒读数。全部数字在仓库外用 `go run` 独立复算过（未改 `internal/ball/tokens.go` 做变异，因桌面/计时代理仍活着；变异后果由 3421 vs 2098 = 63% ≫ 2% 直接算出：谁把本体改成 44px 而不重测，这条就红）。四道门在 HEAD `caaa783` + 本次改动上重跑仍全绿。另记两条 AC#1 顺带查出的静态事实（未改码，见报告）：`SleepWindowEdgePx` 在生产侧零调用者，`hit.go` 的真实点击圈走 `curVisual.SizePx`（整数截断后 ①r=10px ②r=21px，都小于可见光晕半径 26px）；`ball_windows.go:471/474` 的 ULW `constAlpha` 在两种模式下都用 `SleepOpacity`=0.35 作 Settling 终点，而 `applyGlassForm` 的落点是 `RestSettledOpacity`=0.95 ⇒ 原型模式合成结果是 0.95×0.35≈0.33。AC#2/AC#3 依旧未跑、未勾、等桌面。
- [2026-09-20T15:06:41Z] agent=agent-ticket68 did=清掉 `statevisual.go` 里最后一条同源假前提（`Visual` 的 Glass 字段组注释仍写"契约测试保持权威直到 owner 签收新 look (D22)"，而签收已发生在 2026-09-20），改成"冻结读法现在是逃生门"。纯注释、零行为改动；四道门在 HEAD `1d1179e` + 本次改动上重跑：gofmt 空、vet 空、`go test -count=2 ./internal/ball/` ok 0.079s、`./cmd/balldebug/` no test files。**给 AC#2 留的地雷（未动）**：`live_windows_test.go:609 TestBallLiveAudioLiquidGate` 的第一段以"默认即 frozen"为前提（618 行注释"模式冻结优先"）却没显式 `EnablePrototypeVisuals(false)` ⇒ 翻转后它会**自称在测冻结、实际在测原型**，是 winlive 用例（本地跑不了，所以我只登记不改）。
