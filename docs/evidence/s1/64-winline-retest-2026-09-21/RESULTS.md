# 票 64 AC#4 / MINOR-1 / MINOR-2 —— winlive 逐跑记录（2026-09-21 08:29–08:35）

## 为什么要专门做这一份（不是走形式）
票 64 的 **MINOR-1** 登记原文：AC#4 的"不被 skip"当时是**结构成立、逐跑记录缺证**——
理由写得很具体：包里 8 处 `t.Skipf` 全带 `SKIP-LOUD`、`requireQuietBallDesktop` 是 `t.Fatalf`，
但代理只对 A1/A1d 两项写了"未走 skip 分支"，而 **"winlive 14 项全 PASS"这句话在 Go 语义里含 SKIP 也算 ok**。
⇒ 缺的不是"再跑一次"，是**把每一条的名字和它真实的 PASS/FAIL/SKIP 落到纸上**。

## 两次独立跑（不是我复述同一次）
| # | 时间（本地 +08:00） | 谁跑的 | 命令 | 真 exit |
|---|---|---|---|---|
| run1 | 08:29:01 → 08:29:26 | `agent-ticket64-winline`（**跑完后在写票面之前被平台闪断杀死**） | `go test -tags winlive -count=2 -v ./internal/ball/` | 归档在 `run1.log`（291 行，末行 `ok ... 23.795s`） |
| run2 | 08:33:xx → 08:33:55 | **编排者本人**（在同一台机、同一棵树、机器安静：`tasklist` 无 `go.exe`/wisp/balldebug） | 同一条，逐字一致 | **`REAL_EXIT=0`** |

⚠ 那个代理死掉时**只剩"把数字写进票面"这一步没做**，而它整包跑测的产物（`run1.log`、
`results-pass{1,2}.txt`、`table-body.txt`）都在。我没有重跑一遍就当它不存在，
也没有拿它的产物直接当自己的结论——**run2 是我自己那一次**。

## 计数（run2，我亲自从日志 grep 出来的）
- `=== RUN` **128** 行 = **58 个不同测试 × 2 次**（`-count=2` 真的生效，不是缓存假象）
- 顶层 `--- PASS` **116** = 58 × 2；另有 12 条缩进的子测试 PASS；**116 + 12 = 128，一条不落**
- `--- SKIP` **0**；日志里字符串 `SKIP` 出现次数 **0**（⇒ 8 处 `SKIP-LOUD` 分支一个都没走）
- `FAIL` **0**；包末行 `ok github.com/CarlosShao/wisp/internal/ball 24.494s`
- run1 的计数形状相同（128 RUN / 116 顶层 PASS / 0 SKIP / 0 FAIL / `ok 23.795s`）

## 票 64 AC#4 依赖的交互项——逐条真名与结果（run1 与 run2 **都是 PASS×2，无 SKIP**）
| 测试 | 它钉的是哪件交互 |
|---|---|
| `TestLiveClickSummonsAndDragDoesNot` | 单击唤起 / 拖拽不误触发 |
| `TestLiveConfirmingCancelAndEscReturned` | Confirming 态 Esc 取消 |
| `TestLiveNeverStealsFocus` | 焦点不抢 |
| `TestLiveTransparentCornerFallsThrough` | 透明区点击穿透 |
| `TestLiveHotkeyRebindEndToEnd` | 热键端到端重绑定 |
| `TestLiveMuteHotkeyEndToEnd` | M 键切 Muted |
| `TestLiveHotkeyOccupiedVsNotAttempted` | 热键被占 vs 未注册可区分 |
| `TestLiveSleepingZeroTimerHandles` | Sleeping 零定时器（**同批把票 07 第 4 框的"藏在 tag 后面"这条也答了**） |
| `TestBallLiveEdgeDock` / `TestBallLiveEdgeDockHover` | 靠边停靠与悬停弹出 |
| `TestBallLiveLifecycle` / `TestBallLivePositionPersistence` / `TestBallLiveIdleBorderTransition` / `TestBallLiveAudioLiquidGate` | 其余 winlive 项 |

## 判定
- **MINOR-1：闭。** 逐条记录已在纸上（本文件 + `run1.log` + `run2-count2.log` + `table-body.txt`），
  且 **0 条 SKIP** ⇒ 那句"全 PASS"不再是含 SKIP 的假概括。
- **MINOR-2：闭**（同一次跑顺带覆盖）。
- **AC#4：保持已勾**，且现在**有逐跑支撑**（不是"结构成立"这种半句话）。
  退回条件（"交互四项任一 SKIP 就退回未勾"）**没有触发**。
- **连带闭掉票 07 的一处含糊**：`TestLiveSleepingZeroTimerHandles` 真跑过两次。
  我 08:25 在票 07 写的措辞是"断言存在、但藏在 `-tags winlive` 后面、今天是否真执行未证"，
  现在有证了 ⇒ 那句话按下面的更正块处理（票 07 的框**仍不勾**，因为那条 AC 还要 CPU≈0 的完整判据链，
  而"跑过了"只是其中一半）。
- ⚠ **本目录不能证明的事**：这些是**桌面安静条件**下的结果（`requireQuietBallDesktop` 会 `Fatalf` 而不是 Skip，
  所以"能跑完"本身就说明桌面是静的）。它们**不等于**多显示器真拖（票 07 第 6 框，今天零证据）、
  也不等于票 68 AC#2/3（那两条阻在 owner 的 R15 #3/#4/#5）。
