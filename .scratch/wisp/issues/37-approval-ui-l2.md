# 37 — Approval UI: queue view, L2 card panel-grade, native fallback card verification

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 21-approval-gates-minimal, 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: queue/task list UI + depth badge sync; B: panel L2 card +
native fallback polish + batch view)
**Spec refs:** SPEC-08 §5.4 (approval rows), §17.5, C18, D31 visibility mandates, D45-1 batch UI,
C27 fallback

## What to build
The approval experience at panel grade: queue head + depth indicator, per-request full context
(task, tool, params, reason), batch confirmation view, and the verified native fallback card
path when the panel is unavailable — with the ball depth-badge as the always-visible channel.

## Key constraints
- Visibility mandates (D31): queue depth ALWAYS on ball badge (sync <100ms); queued/blocked
  tasks visible in panel task list WITH "waiting on what"; every request self-describes
  (task, tool, FULL params, risk, reason) — never a bare "需要确认".
- Panel L2 card: top 2px danger bar, tool + L2 pill, complete params, rules_hit reasons,
  batch expandable list (first 5 + expand), persistent "点击悬浮球以批准" guide + pulsing ball
  glyph, buttons = reject + view-full-params ONLY (no allow — F2/§15#6; redundant assert here).
- Native fallback (panel unavailable): 320×140 native card (allow + reject + open-panel ghost
  link) — verify C27 rule "L2 capability never disappears"; visuals per §17.5 native-row.
- Timeout UX: 270s warning state on card (30s before 300s auto-reject); post-reject one-click
  replay button (C18).
- Queue FIFO order display; correlationId shown for forensics; decisions (allow/reject/timeout)
  land in tool_call rows with grant linkage hooks.
- Batch aggregated confirmation renders total + affected dirs + expandable detail; L2 batches
  never aggregate (guard re-asserted at UI level).

## Out of scope
- Session-grant options in card (49 adds the third option); multi-task scheduler itself (47).

## Acceptance criteria
- [x] E2E: two sequential L2 requests → depth badge 1→2 (simulated pre-47 via queue test mode),
      head-only display, decision routing by correlationId.
- [x] Timeout: card warning at 270s → auto-reject at 300s → replay works; ball badge clears.
- [x] Fallback: with panel runtime masked, L2 still approvable via native card (33 fixture).
- [x] No-allow invariant: panel card markup contains no allow action wired to decision API
      (DOM + bridge double test).
- [x] Batch view: 10-op batch renders aggregated card; 50+ upgrades to L2 card (R7).
- [ ] Human visual sign-off of card + fallback (screenshots vs design/screens/approval.html).

## Progress log (append-only, newest last)

- 2026-09-21 15:4x（编排者，**接线时必须做的一次键对账 = 票 87 的残留 R-3**）：
  `acceptor-ticket87` 实测（其报告 §6.2 攻击 2）：**卡片 B 的"入键"与卡片 A 的队列签发键撞车时，
  宿主按入键取消 B 会拒到 A**。方向是拒绝 ⇒ **不触"查不到就放行"那条底线**，而且 `reject()` 这条路线
  **先于票 87 存在**；票 87 新增的只是"`Veto` 这条以前什么都不拒的路线现在也会拒错"。
  ⇒ **判据（票 37 接线时算你的验收项）**：`internal/panel/approval.go` 带来的 C17 入键（`CorrelationID`
  取自 `ApprovalSubject.CorrelationID`，见 `:28/:40/:78`）与队列自己的签发键（`queue.go:133-139` 的 `it.Corr`）
  必须做**一次显式对账**，并且要有一条用例：**两张卡共存时，取消其中一张绝不能命中另一张**
  （今天两者相同只因为 `internal/agent/loop.go:644` 把 corr 写成 taskID 且队列没重签发——**那是巧合，不是设计**）。
  相关：票 87 已 `-done`（库层就绪、用户层未接线的判词在 `docs/evidence/s1/87-adversarial-acceptance.md` §12），
  R-1/R-2 在**票 97**。

> **2026-09-26 进度追记（owner 把前端交 ZCode 直管后的可视壳进度，只追加不勾框）**：
> 本票的可视面已按 beautiful-ui 主题收编第四代重画，演示形态在 `?harness=1`
> 应用全观（左栏三段式/会话更改卡/右栏审查·终端·浏览器）与 `?harness=2` 陈列室。
> **本票 AC 的真数据 e2e 格仍等票 145/35/33 的字段与宿主，一格未勾、不由本程
> 勾**；可视壳的 commit 链见票 77 进度追记（`424ac84`…`5e23d99`）与
> `docs/reports/frontend-session-log-zcode.md` §1–§7。

- 2026-11-16 22:5x（票 37 写腿，前端本编队铁律 4；只追加不勾框）：
  队列读面接线（面板侧半张票）。生产行进链已核现量：
  `internal/agent/approval/pending_read.go:162` `Queue.LiveApprovals()`（队头在先 FIFO，
  `Position` 是 1-based 队位）→ `cmd/wisp/panel_pump.go:78` `nativeVerdictsFromGate`
  （"copy, not a judgement"；CallChain 生产恒空、注释点名票 37 证据文件）→
  `internal/panel/pump.go:106` `CardView()` → `Snapshot.pending`
  （`internal/panel/composer.go:58`）→ 页面收包 `frontend/src/panel-root.tsx:37`。
  本次改动（全部前端，未动 Go 一字节）：
  1. `l2-approval-card.tsx`：卡上加 correlationId 原文一行（title 挂全文）——票面
     "correlationId shown for forensics" 的面板卡半格；
  2. `approval-screen.tsx` 重写：FIFO 序号（1..n 按线上顺序）+「队头在先（FIFO），悬浮球
     只答队头」+ 容量句改 `QUEUE_CAPACITY=8`（= `internal/agent/approval/queue.go:130`
     DefaultMaxPending，原来写 10 是没产读者的假数）+ 每行「等待于：」走快照任务名册
     `blockedOnApproval` join（= Go `TaskRosterSectionFrom` 的 indexWaitingKey 配对，
     `subagent_roster_197.go:186-247`）；
  3. 新组件 `approval-queue.tsx`：审批队列条（深度 + 队头 + 每卡 corr/等待于），挂
     `App.tsx` 卡片上方每一视图（D31 队列常显的面板半格）；空队列零渲染（不再摆常驻
     零件训练忽略）；
  4. 新牙 `scripts/render-approval-37.tsx` + `npm run render:approval37`：喂线上形状真值
     ⇒ 断言同一份值上屏（corr/规则/argv/原因/深度 2/队头/等待于/容量 8），并钉三枚
     "具名缺席"：D45-1 批量摘要（`BatchView` 只走原生 Prompt 面 `ui.go:193`，卡线无
     字段）、270s/300s 倒计时（线无 deadline 字段，倒计时在原生卡 `cardrender.go`）、
     每卡深度戳（Prompt.Depth 走球角标，页面不画）——面板 L2 批量聚合过 R7 的
     升级守卫继续由 Go 侧 `gate.go:317-324` 钉住；无 allow 控件断言对两处新面重申。
  未勾格与原因（真数据 e2e 仍缺宿主，与 2026-09-26 追记同因）：
  - AC#1（两卡顺序 1→2、badge 同步 <100ms、correlationId 路由）：球角标走原生
    `approvalDepthBadge`（`resident_approval_windows.go:1266`），非面板；100ms 同步与
    decide-by-corr 的 e2e 要真宿主，本轮以 `render-approval-37` 的深度 2/队头/corr 断言
    覆盖面板半格（`internal/agent/approval/ticket48_c18_full_test.go:135` 是 Go 侧
    FIFO/corr 路由真 e2e）。
  - AC#2（270s 警告→300s 自动拒→重放）：线无 deadline 字段（具名缺席见上）；Go 侧
    时梯已在 `ticket48_c18_full_test.go:310` 钉住；重放 `Gate.Replay`
    （`gate.go:885`）面板无路由，一键重放按钮属未接线（缺：把 replay 暴露进 C17 白名单
    ——改契约须人工批）。
  - AC#3（面板遮蔽时原生卡可批）：C27 原生卡是 `resident_approval_windows.go` 的
    320×140 卡，非面板写面；本腿零涉及。
  - AC#4（面板卡无 allow）：本轮重申 + `render-l2` 原有断言继续钉。
  - AC#6（人工签收）：等编排者排屏（`?harness=1` 可看新队列条与卡 corr 行）。

---

## 编排者收件·第一批（2026-11-17 09:2x，`232a5da9e`，仪器 `render:approval37`）

**尺1** `232a5da9e` 在。**尺3 我自己复跑**：`npm run render:approval37` **rc=0**；build/typecheck/lint 全 rc=0。
- **勾 AC#4（No-allow invariant）**：`frontend/src/components/l2-approval-card.tsx` 里**没有**任何可点的 allow 动作接到 decision API —— grep `allow|decide|decision` 只命中注释与 `判定来自 ${view.decidedBy}` 展示行（`:264/:347`），**零操作接线**。
- 顺带修一处假数：容量句原写「10」无读者，已改 `QUEUE_CAPACITY=8` ＝`internal/agent/approval/queue.go:130 DefaultMaxPending`。
- **留白（具名，均需真宿主或 Go 读面）**：AC#1 深度 1→2 e2e（需宿主）、AC#2 270/300s 倒计时＋球徽标清（球侧）、AC#3 原生回退卡、AC#5 BatchView（`internal/agent/approval/ui.go:193` 只走原生 Prompt 面，卡无字段）、AC#6 人眼签收（机主）。

---

## 编排者追加·BatchView 字段已落地（2026-11-17 10:3x，`b1d2162dd`）

`internal/panel/approval.go:81 ApprovalCardView` 新增五枚平铺 `omitempty` 字段 → `internal/panel/pump.go:126`（NativeVerdict.Batch 映射）→ 数据源 `internal/agent/approval/cardsummary_37.go:39 CardSummary`（display 字段 L1/L2 ≥3 目标都答，R4/R3 污染不出；**L2 永不聚合契约原样、`gate.go` R7 守卫未动**）。
⇒ 前端 AC#5 剩：把聚合卡五字段画出来（10 操作一张卡），50+ 升级 L2 走既有守卫。

---

## Progress log（续，2026-11-17 11:4x，本票前端消费腿 agent-109 · 读面消费波；只追加不勾框）

**本轮交付（frontend/**＋internal/panel/composer_test.go 一处；证据件 `docs/evidence/s1/readface-consume-20261117-r1.md`）：**

1. **AC#5 页半落地**：`frontend/src/lib/panel.ts` 的 `ApprovalCardView` 声明五枚可选 batch 键（batchTotal/batchDirs/batchFirst/batchDeferred/batchSummary，对应 `internal/panel/approval.go:81-85` omitempty 平铺字段；门＝batchTotal>0＝Go HasBatch()，非 batch 卡 wire 逐字节同旧形）；`frontend/src/components/l2-approval-card.tsx` 新 `BatchBlock`：batchSummary 原句逐字＋本批共 N 项（含未展开余数）＋涉及目录（Go 去重排序逐字）＋前 N 项路径预览（调用自身顺序逐字）；**块内零控件**——no-allow invariant 不变（「批准权在原生侧 · 面板不提供允许按钮」与球提示行原样）。
2. **同 commit 删 ledger 登记**：`internal/panel/composer_test.go` composerTSGapLedger 的 `"ApprovalCardView"` 条目删除——ledger 双向规则自点名（page 一声明⇒stale-red，须删条目由普通双向比对接管）；`TestApprovalCardViewJSONKeysMatchFrontendTypes` 复跑 rc=0。
3. **牙**：`render-approval-37.tsx` 新 4b 节——10 操作 ⇒ **一张卡**（「需要你的确认」×1）＋五字段逐字回屏；普通卡零 batch 痕（无「本批共/涉及目录」）；batch 卡按钮名断言无 allow 形；panel.ts 五键声明钉＋Go ledger 条目消失钉（双向都钉）。

**门禁读数**：build/typecheck/lint rc=0；`render:approval37` rc=0；`render:l2`（喂 fixtures/l2-card-fs-delete.json）rc=0；`go test ./internal/panel` rc=0；d22scan rc=0；⛔ 未 push；dist 未入库。
**留白不变**：AC#1 深度 e2e（真宿主）/AC#2 270/300s 倒计时（原生卡）/AC#3 原生回退卡/AC#6 人眼签收（机主）。

## 编排者收件·第二批（2026-11-17 11:4x，`f392f34a2`，件 `docs/evidence/s1/readface-consume-20261117-r1.md`）

**尺1** `f392f34a2` 在（16 件）。**尺3 我自己复跑**：`render:approval37` rc=0；`go test ./internal/agent/approval -run 'TestR7…|TestBatch'` **ok 0.052s**；build/typecheck/lint rc=0。
- **勾 AC#5**：①10 操作 ⇒ **一张卡**（「需要你的确认」×1）＋五字段逐字（batchSummary 原句／本批共 10 项／余 7 项／目录×2／预览×3），普通卡零 batch 痕、batch 块**零 allow 形控件**；②50+ 升级 L2 **不聚合** —— Go 侧 `internal/agent/approval/batch_test.go:87 TestR7SizedBatchEscalatesToL2Unaggregated` 现跑 PASS（断言 level=L2、不得聚合、`rules_hit` 带 R7）。
- 留白：AC#1 深度 1→2 e2e（需宿主）、AC#2 270/300 倒计时＋球徽标（球侧）、AC#3 原生回退卡（33 夹具/宿主）、AC#6 人眼签收（机主）。

## 写腿收件·三格 Go 侧（2026-11-18 21:x，写腿 agent-139；只追加不勾框）

证据件：`docs/evidence/s1/37-depth-timeout-fallback-20261118.md`（三格读数、三枚变异牙红 log、门禁全读数）。
新牙：`cmd/wisp/approval_depth_timeout_fallback_37_windows_test.go`（四枚）＋
`internal/agent/approval/ticket37_timeout_boundary_37_test.go`（两枚边界读数牙）。
- **AC#1（Go/球侧半格）**：两发顺序 L2（票 43/48 的 raise-and-wait-for-depth queue 测试模式）⇒ 徽章输入
  1→2（`Prompt.Depth` 逐字 1、2，`approvalDepthBadge` 映射钉）、队头单显、按 corr 分别路由（先答第二发
  ⇒ B 醒 A 不醒；跨卡花 grant ⇒ ErrBadGrant）；票 87 R-3 键对账（签发键＝卡 corr＝面板 verdict corr）与
  「两卡共存取消其一绝不能命中另一张」（共享名歧义 ⇒ ErrUnknownCorrelation、精确键 ⇒ 仅命中己卡）都钉进
  `Test37TwoCardsCoexistingCancelOneNeverHitsTheOther`。牙：两发算一发（Depth 钉 1）⇒红、按到达序答
  （reject 无视 corr 命中队头）⇒红，均载具现跑现还原。
- **AC#2（Go/球侧半格）**：注入 Clock 走 `approval.Clock` 缝（⛔墙钟判据），gate 零超时选项 ⇒ 300s/30s
  来自 `queue.go` 冻结默认：270s 一枚 Warning（Remaining 30s、「30 秒后自动拒绝」）、300s `AnswerTimeout`
  （why 含 300＋拒绝）、深度 0＋ledger 空＋`AwaitingHuman()==false`＋真实走 `ballCardUI.Update` 的
  settle/syncBadge 分支；`Gate.Replay` 新 corr 新 grant、死 grant 打不开重放卡、新 grant 补完原请求。
  牙：`DefaultApprovalTimeout` 300→299 ⇒ approval 包直跑红（`Queue().Timeout()=4m59s`）＋载具红
  （卡 Deadline 读数 4m59s），证边界是读来的；⚠ 270/300 不在 `internal/observe/thresholds.go`，真源头
  `internal/agent/approval/queue.go:125`（冻结常量，只读未改，变异为临时探针已还原）。
- **AC#3**：复用票 33 AC#5 fixture 本体（`withMaskedRuntime`）⇒ 按名拒绝＋`panel.unavailable`＋
  `Unavailable()==(true,reason)`，被遮世界里两发 L2 经 `Replies.Allow`（原生卡/托盘动词）先后批准、
  徽章 2→0；负向：活 grant 走面板路由报允许 ⇒ `ErrPanelAllow`＋`PANEL-ALLOW-REJECTED`＋nonce 烧毁
  （原生再花 ⇒ ErrBadGrant）；拒绝方向仍通。⛔ 未给面板侧开 allow 口子。
- **门禁**：d22scan 0 findings（末行 OK）；载具相关族 rc=0；`go build ./...` rc=0；前端 build/typecheck/
  lint＋`render:approval37` 全 rc=0；载具整包收口一发读数见证据件。
- **留白不变**：AC#6 人眼签收（机主）；重放面板一键按钮仍缺 C17 白名单 verb（改契约须人工批）；
  Win32 实球 posting 的实窗核验仍归 winlive 档（AST census 牙继续钉接线）。

## 编排者收件·三格（2026-11-18 21:4x，`64dfe1290`，件 `docs/evidence/s1/37-depth-timeout-fallback-20261118.md`）

**尺1** `64dfe1290` 在（4 件）。**尺3 我自己复跑**：载具四枚 PASS（Badge1→2 0.01s／TwoCardsCoexisting 0.01s／Timeout270-300-Replay 0.02s／MaskedNativeCard 0.02s）；`internal/agent/approval ^Test37TimeoutBoundaryIsReadFromTheContractConstants` PASS 0.01s；`render:approval37` rc=0；`d22scan` clean。
- **勾三格**：①徽章 1→2（卡面 Depth 逐字 1,2）、恒队头、**按 correlationId 分别答**（先答 B 只醒 B；跨卡 grant=ErrBadGrant；R-3 精确键只中己卡、歧义名=ErrUnknownCorrelation）②**270s 预警（"30 秒后自动拒绝"）／300s AnswerTimeout 自动拒／重放换新 corr 能补完／球徽章清零**，边界值读自 `queue.go:125` 冻结默认（270/300 **不在 thresholds.go**，该文件零动；变异 300→299 两面红）③面板被遮 ⇒ 原生卡经 `Replies.Allow`（球/托盘生产路径）两发先后批准；面板 allow=ErrPanelAllow+PANEL-ALLOW-REJECTED、nonce 烧毁。
- ⚠ 整包 3 红分诊：33r5／320（台账在册环境红）＋`Test39NativeLooseningCardAcceptedAppliesAndAudits`（**不在名册**；我单跑 3/3 全绿 ⇒ 整包内互扰，非产品红；见下一条处置腿）。