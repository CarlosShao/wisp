# 票 74 对抗验收 —— C21 令牌表必须对齐真正驱动像素的常量

**验收人**：编排者（非实现者；实现方 `agent-ticket74`，5 个 commit）
**验收时间**：2026-09-21 10:16 · **被验收 commit**：`1d9ab9b` `cf8406e` `186ad31` `ed43bc3` `1d94caf`
**全部动作在 `git archive HEAD` 解到 `/tmp` 的纯净树里做**（当时有四个代理在共享树上写码）。
**档位**：〔独立复现〕／〔日志＋归档，我抽验〕／〔仅自述，不背书〕

| AC | 裁决 | 我跑的判据 | 档位 |
|---|---|---|---|
| **AC#1** 零消费者的 `SwimLevelGain`/`SpinLevelGain` 不得"声明了没人接" | **PASS** | `git show --stat 1d9ab9b` ⇒ **删除与替代同在一个 commit**（5 文件：`tokens.go` 删两个常量 + `renderer_windows.go` 改指真驱动 + 表 + 机器检查 +225 + 票面）。我读了渲染器那段 diff：原本是**裸字面量** `1 - 0.12*v.LiquidLevel` / `1 + 0.10*v.SummonFlow`，被提成 `LiquidGatherPerLevel`/`SummonFlowSpread`。⇒ 漂移的**根**不止"声明未接线"，还有"接线了未声明"。 | 〔独立复现〕 |
| **AC#2** 每个改变像素的常量都要么进表、要么进带理由豁免 | **PASS** | `go test ./internal/ball/ -v` 里 `TestC21CodeTokensAreTabledOrExempt` 逐条打印理由，我在日志中读到 4 条例子（`HTClient`/`HTTransparent`/`HTNowable`… 与 `AltSummonSpace`），全套餐 **FAIL 0**；"29 条漂移"的中间态已被 `cf8406e` 改成**必判**。 | 〔独立复现〕 |
| **AC#3** 新断言必须真会红（containment 改单射） | **PASS（我自己重做）** | 我把 `tokens.go:373 DockTriggerPx = 16` 改成 `17`，先 grep 证落盘（`: DockTriggerPx = 17`）⇒ **`--- FAIL: TestC21GeometryRowsMatchCodeConstants`**，报文是 `this row states no **UNCLAIMED** number equal …`。"UNCLAIMED"这个词就是单射本身——**已被别的行占号的数字不能再顶包**。 | 〔独立复现〕 |
| **AC#4** 表↔`tokens.css` 缺口：加检查或带票号移交 | **PASS，未移交（比要求多做了一半）** | 纯净树跑到 `TestC21TableColourRowsMatchTokensCSS`：`three-way colour check: 78/78 cited rows agree …` 且 `--- PASS`。⇒ 票 69 转给我的那条缺口**在这里关掉了**，不需要第三张票。 | 〔独立复现〕 |
| **AC#5** 门禁 | **PASS** | 我的读数：`-count=1` → RUN 55 / 顶格 PASS 49 / FAIL 0 / **SKIP 0**；更早一轮 `-count=2` → RUN 108 / 顶格 PASS 96（=48×2 不变式成立）。代理自报 `-v -count=2` 为 RUN 110=2×55 ⇒ 与我 55 个不同测试名**同一口径**，不是两套数。 | 〔独立复现〕+〔抽验〕 |

## 结论：**PASS，5/5** ⇒ 归档 `-done`

## 三条它关不掉的缺口，我的处置（都带落点）
1. **`renderer_windows.go` 里仍有一批没名字的裸字面量**（`liqRate{1.0,-1.6,0.55}`、`liqPhase{0,2.1,4.2}`、
   渐变停靠 0.42/0.78、`permille 775+225…`）——**任何"枚举常量"的检查在结构上都看不见它们**。
   提出来是**改绘制代码**，需要真机重看视觉 ⇒ 与票 62/65 的真机签收同批做，**登记为 Q-19 的一部分**。
2. **数字单射只做了一半**：一个数字不会再满足两个常量，但"行内每个数字都得有主"需要给散文加数字标注
   （`SPEC-08 §2.1`、`96 DPI`、commit 哈希都算），日志里剩 36 条 ⇒ 登记，不追。
3. **owner 级三条**（look 三色升进 `tokens.css`、`CountdownFontPx` 与 `BadgeFontPx` 10px 重叠、
   `SleepingDotPx`/`SleepOpacity` 三方值漂移）⇒ 并入 **Q-19**，等拍板。

## 一条方法论（代理与我**在同一小时内各自撞到**，值得单独立档）
代理自陈"两次变异一开始没落地（tab 不对 / 行号漂移），表现为**假绿**，被'跑之前先 grep 证落盘'救回"。
我在票 73/74 的复现里**连踩两次同样的假绿**：先拿 `155s/16px→17px` 打表文件（那行已被重做成字体行），
再拿符号名首次命中的行下手（那是描边宽的行）——**两次都得到 `ok`，而 `ok` 是零信息**。
⇒ **锚点必须是"我要改的那个数字所在的那一行"，不能是"符号名出现的地方"或"记忆里的行号"**；
表格类文件的行号会随别人的重做漂移。第三次我把 `grep -n` 与被改后的整行放进同一条 `&&` 链才拿到真红。
