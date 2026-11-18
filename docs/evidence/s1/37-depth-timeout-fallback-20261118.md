# 票 37 三格收件（写腿 agent-139，2026-11-18）— 深度徽章＋corr 路由 / 270·300 超时梯 / 面板遮蔽原生回退

票面：`.scratch/wisp/issues/37-approval-ui-l2.md` AC#1/AC#2/AC#3（AC#4 已勾、AC#6 人眼签收不在本腿）。
派单：`.scratch/wisp/probes/orch/20261118/txt/37-depth-timeout-fallback-dispatch.txt`（存档 commit `10e8d8dbe`）。
新牙文件：`cmd/wisp/approval_depth_timeout_fallback_37_windows_test.go`（四枚测试）＋
`internal/agent/approval/ticket37_timeout_boundary_37_test.go`（两枚边界读数牙）。

## 格 1 — 深度徽章 1→2＋队头单显＋按 correlationId 路由（pre-47 模拟，既有 queue 测试模式）

- **测试**：`Test37TwoSequentialL2BadgeOneToTwoHeadOnlyRoutedByCorrelation`（单一 task、两发顺序 L2，
  `ask37` 走票 43/48 的 raise-and-wait-for-depth 模式）。
- **读数（载具绿跑，`carrier-related-green.log`）**：
  - 徽章输入 1→2：卡 A `Prompt.Depth==1`、卡 B `Prompt.Depth==2`，`approvalDepthBadge`（`ball.SetBadge`
    的唯一生产者输入）映射逐字 1、2；
  - 队头单显：`Panel().Head()==corr-37-a depth 1`，卡 B 在自己 FIFO 位（`View(corr-37-b).Depth==2`）；
  - **键对账（票 87 R-3 验收项）**：队列签发键（`LiveApprovals[].CorrelationID`）＝卡上 corr＝面板 verdict
    corr，逐枚相等——一次显式对账钉在本测试 (1c)；
  - 路由：先答第二发（`Replies.Reject(corr-37-b)`）⇒ B 的 ask 醒、A 的 ask 不醒（cross-delivery 负检查）；
    A 卡 grant 花在 B 卡 ⇒ `ErrBadGrant`、双卡纹丝不动；排水到 0，`AwaitingHuman()==false`。
- **牙①（两发算一发⇒红）**：变异 `gate.go` `p.Depth = g.q.Depth()` → `p.Depth = 1`。载具跑
  `-test.run=Test37TwoSequential` ⇒ **红**（`card B's stamped depth = 1, want 2` ＋
  `approvalDepthBadge(card B) = 1, want 2`），log `.scratch/wisp/probes/37/r1/logs/t1-two-as-one-red.log`。
  已还原（`git diff internal/agent/approval/gate.go` 空）。
- **牙②（按到达序而非 correlationId 答⇒红）**：变异 `queue.go` `reject` 的查找——无视 corr、命中
  `q.pending[0]`。载具跑 `-test.run=Test37TwoSequential|Test37TwoCardsCoexisting` ⇒ **两枚红**
  （共享名取消命中队头＋交叉投递），log `t2-arrival-order-red.log`。已还原（git diff 空）。
- **共存取消不串（R-3）**：`Test37TwoCardsCoexistingCancelOneNeverHitsTheOther`：
  双卡共存时按共享入键（同 taskID）取消 ⇒ `ErrUnknownCorrelation`、两 ask 均不醒（歧义不猜）；
  按精确签发键取消 B ⇒ 仅 B 结清、A 的 grant 仍持；B 走后同一共享名 ⇒ leniency 落到唯一幸存卡
  （拒绝方向、无歧义才落地，边界两向都钉）。

## 格 2 — 270s 预警 / 300s 自动拒 / 重放可用 / 球徽章清零

- **测试**：`Test37TimeoutWarnsAt270AutoRejectsAt300ReplayWorksBadgeClears`（cmd/wisp，注入
  `fakeClock37` 走 `approval.Clock` 缝，gate 零超时选项 ⇒ 300s/30s 来自 queue 冻结默认）＋
  `Test37TimeoutBoundaryIsReadFromTheContractConstants` / `Test37ExpiryWordingDerivesFromTheInjectedTimeout`
  （approval 包边界读数牙）。
- **读数（载具绿跑）**：卡上 `Deadline==300s`（字面钉死）；270s 弹一枚 `EventWarning`（corr 对、
  `Remaining==30s`、文案「30 秒后自动拒绝」），卡仍可答（徽章输入 1）；300s ⇒ `AnswerTimeout`、
  why 含「300」＋「拒绝」、深度 0、ledger 空、`AwaitingHuman()==false`，且真实走了一遍
  `ballCardUI.Update`（settleOrb＋syncBadge 分支，headless 安全）；`Gate.Replay` ⇒ 新 corr
  `corr-37-t#replay`、新 grant ≠ 死卡 grant、新 grant allow ⇒ 原请求补完、徽章输入 0。
- **牙（300 改 299 逻辑等价⇒红，证边界是读来的）**：变异 `queue.go` `DefaultApprovalTimeout`
  300s→299s。两处现红：
  - approval 包直跑：`Test37TimeoutBoundaryIsReadFromTheContractConstants`（`Queue().Timeout() = 4m59s,
    want the frozen 300s`）＋既有 `TestTicket48TimeoutLadderWarn270Reject300Replay`，
    log `t3-299-approval-pkg-red.log`；
  - 载具跑 `-test.run=Test37TimeoutWarns` ⇒ **红**（`the card's deadline reading = 4m59s, want the frozen
    300s C18 boundary`），log `t3-299-carrier-red.log`。
  已还原（git diff 空）。⚠ 270/300 不在 `internal/observe/thresholds.go`（那文件本腿未动一字）；
  真源头是 `internal/agent/approval/queue.go:125` 冻结常量，只读未改，变异为临时探针、跑完即还。
- ⛔ 全程无墙钟判据：Clock 注入 Advance；`waitFor37` 的真实时间只作有界轮询预算（既有 48 模式）。

## 格 3 — 面板运行时被遮（33 fixture）⇒ L2 仍能经原生卡批准；⛔面板侧 allow 开口

- **测试**：`Test37MaskedPanelL2StillApprovableViaNativeCard33fixture`。直接复用票 33 AC#5 的 fixture
  本体（`withMaskedRuntime`，`panel_runtime_missing_33ac5_windows_test.go`——该文件注释里点名的
  消费者就是本卡）。
- **读数（载具绿跑）**：mask 后 `mgr.Show` 按名拒绝（`errPanelRuntimeMissing`，窗口未建）、
  `panel.unavailable` 事件在 slog、`Unavailable()==(true, reason)`；同一被遮世界里两发 L2 经原生卡
  动词 `Replies.Allow`（球/托盘按钮的生产路径）先后批准、两 ask 各自醒、徽章输入 2→0——
  C27「L2 能力不因面板消失」成立；
  **负向**：第三张卡的真实活 grant 走面板路由报允许 ⇒ `ErrPanelAllow`（F2 第三层，路由级拒绝）＋
  `PANEL-ALLOW-REJECTED` 审计线（gate Logf 面）＋该 nonce 被烧（原生侧再花 ⇒ `ErrBadGrant`）；
  面板拒绝方向仍通（`DecideFromPanel Allow:false` 落地）。**未给面板侧开任何 allow 口子。**

## 门禁读数（2026-11-18 现跑）

- `sh scripts/d22scan.sh`：**0 findings**，末行 `d22scan.sh: OK - control rc=0, scan rc=0`（log `d22scan.log`）。
- 载具相关族（Test37×4＋Test37VerdictBatch＋Test48×3＋TestAC5×4＋TestTicket21S2×3＋TestTicket48×3）：
  `ok github.com/CarlosShao/wisp/cmd/wisp`，**rc=0**（`carrier-related-green.log`）。
- `go build ./...`：**rc=0**。
- 前端：`npm run build` rc=0 / `npm run typecheck` rc=0 / `npm run lint` rc=0（50 warnings 0 errors，
  既有水位）/ `npm run render:approval37` rc=0。
- 载具整包收口一发：见下节（后台单发，串行纪律遵守，跑前 `Get-Process wisp.test` 空）。

## 载具整包读数（照实记）

载具单发整包（`carrier-full-close.log`，2026-11-18 21:06–21:39，墙钟 2018s）：
**PASS=755 / FAIL=3 / SKIP=13，rc=1**。按 A1889 三件套口径逐枚归因：

| 红名 | 成因 | 证据 |
|---|---|---|
| `TestAC4FocusReturnToPriorWindowGap33r5` | **台账红名册在册**（前台被占，环境） | `pending-and-issues.md:37334` |
| `Test320FetchLegWritesTheKeyNowhereOnDisk` | **台账红名册在册**（整包内扫 `%TEMP%` 互扰；定向单跑绿） | `pending-and-issues.md:37332/:35852` |
| `Test39NativeLooseningCardAcceptedAppliesAndAudits` | **不在台账红名册**；本腿定向单跑复验 **PASS 0.03s rc=0**（`test39-targeted-recheck.log`）；整包内失败读数＝审计 sink 竞态（dump 里 SECTION-WRITE 仍是 `outcome=queued`，accepted 行未及落，`accepted SECTION-WRITE lines = 0`）——与 320 同族「整包内互扰」形状，**报编排者裁定是否入红名册** | 失败段在 `carrier-full-close.log` 21:09:47 段 |

- 本腿四枚新牙在整包里**全 PASS**（0.01–0.02s）。
- 该整包含 33r5 红＝**既有环境读数，照实记**，不由本腿修复（不碰他人在飞件、不改任何断言）。

## 边界与留白

- 本腿只动：`cmd/wisp/**`（测试）、`internal/agent/approval/**`（测试＋临时变异探针已还原）、
  `docs/evidence/s1/`、票面追记。`internal/panel/**`、`frontend/src/**` 本轮零改动（面板半格与
  页面半格由此前 readface 腿交付，本轮读数继续成立）。
- AC#6（人眼签收）不属写腿；重放的面板一键按钮接线仍缺（C17 白名单加 verb 属改契约，须人工批，
  台账既有记录）。
- 变异牙三枚全部「改→跑→还原→git diff 验净」，底稿在 `.scratch/wisp/probes/37/r1/*.pristine`。
