# 194 — 面板那两份"前端能调后端哪几个动作"的名单互不相认：owner 裁定**以当前规格为准补齐代码侧**，别留债务

- Status: **ready-for-agent（先只读核代价，再分批落地）**。⚠ **功能票＋契约变更已获点名批准**（owner 09-28 16:5x 原话逐字在台账 `A381`）。
- 来路：`33-a1`（只读设计核，表 `docs/evidence/s1/33-inbound-hop-design-a1.md` §Q1，台账 `A374`）现量出两份名单交集为零；`Q-67` 摆给 owner，他当场裁：**「方法名册，这个建议你按照当前规格补一下吧，别搞债务了」**。
- 关联：**票 33**（面板宿主层，入向那一跳）· **票 92／R20**（面板里的档位与工作区是**权限输入口**：只许显示＋发起请求）· **票 146／`Q-49`**（面板侧批准出口那一族刚钉过五枚门）· **票 141／`Q-46`**（`ban #6` 的射程）· **票 186／187／191／192**（每一枚都在问"我这个动作的名字算不算批过"）

## 现量（两份名单，逐字，尺可复制；锚 `8ca2291f`）

| 哪一份 | 内容 | 尺 |
|---|---|---|
| **规格那份** | `docs/specs/SPEC-08-ui-ball-panel.md:163-174` 那张表：`panel.resync`／`tasks.list`／`task.detail`／`history.query`／`transcript.get`／`approval.current`／`approval.queue`／`approval.decide`／`config.get`／`config.set`／`grants.list`／`grants.revoke`／`privacy.purge`／`privacy.export`／`cost.summary`／`models.list`／`models.delete`／`diagnostics.export` ＋ 事件推送那一行 | `sed -n '160,176p' docs/specs/SPEC-08-ui-ball-panel.md` |
| **代码那份** | `internal/panel/bridge.go:42-45` 恰好 4 枚：`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send` | `grep -cE '=\s*"panel\.' internal/panel/bridge.go` ⇒ **4**（⚠ 别用 `grep -c 'panel\.'`，那把尺起手是 5——`bridge.go:35` 是一句含 `"panel.*"` 的注释，`33-r1` 具名报过，见 `A380`） |
| **交集** | **零枚实现**；两处同名噪声（`approval.queue` 在 `internal/statemachine/events.go:43` 是**状态机事件名**；`approval.decide` 在 `tools/d22scan/main.go:23/361/827` 是 **ban #6 的注释与拒绝文案**） | `grep -rlF "<方法名>" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go` 逐枚 |

## 为什么值得做（不做会怎样）

不补，**每一枚新界面功能都要重新问一遍"这个动作的名字算不算批过"**——这两天票 186／187／191／192 全都卡在这同一枚问题上，owner 那句"别搞债务了"点的就是它。⚠ 另一面也要报：**规格那份表是 S5/S7 切片卡的草稿口径**，里面有几枚（`privacy.purge`、`models.delete`、`config.set`）**今天连后端能力都还没有** ⇒ "按规格补代码"不等于"把 18 枚一次做完"，本票的 AC#1 就是逐枚把"有货／没货／撞禁区"分开。

## AC（每格都要答"这一发在未修码上响不响"）

- [ ] **AC#1 先交一张逐枚代价表（只读）**：规格那 18 枚（＋事件推送那一行）**逐枚**答四件事——① 今天有没有等价实现（尺要现跑，不许按名字猜）；② 落点在哪（复用现有 handler／新建）；③ **撞不撞禁区**（见 AC#2 那三条）；④ 要不要新的门控与审计。**枚数与名册一律现跑，不许从本票面抄。**
- [ ] **AC#2 三条禁区不许因为"按规格补"被顺手动到**（这三条各出自不同批次的人工批准，owner 这次只批"补名册"，没批动它们）：① **`approval.decide` 的「allow」侧永不允许从面板发起**（`SPEC-08:169` 逐字：「**「allow」拒绝一切面板来源（F2）；仅 `reject` 可面板发起**」，出处 D33／`Q-49` 那一族）；② **面板里的档位与工作区是权限输入口**——只许显示＋发起请求，不许面板直接改权限状态（R20／票 92）；③ **`config.set` 的安全节放宽要走 L2 重新确认**（`SPEC-08:170` 逐字，出处 `SPEC-03 §4.2`）。⇒ 任何一支"必须动这三条才补得齐"＝**停手上报**。
- [ ] **AC#3 白名单扩张要**逐枚**落 `A##`**（文件／行／理由／边界／撤销口令），不许拿"owner 批了名册"当全域通行证一次洗平；⚠ `bridge.go` 那四枚常量的**顺序与命名**别顺手改（有判据按磁盘源码数枚数，`181-r1` 的 `TestGitDimensionHasNoModelCallableTool` 就是这种尺）。
- [ ] **AC#4 每一枚新增方法必须自带三枚常驻判据**：① 白名单外方法名 ⇒ 拒答＋审计（不许静默丢弃）；② **面板侧来源的"批准"不许变成放行**（票 146／`Q-49` 那一族，⚠ 这一形历史上被验收造出过一次完整绕过）；③ 归因字段（`RequestID`／taskID）**不许被入口重写**——⚠ 别照抄票 33 派单里那句"改 `source` 必红"，那发今天走在不可能路径上（`bridge.go:93` 在路由前就钉死 `Source`），**可观测字段是 `RequestID`**（`A380`）。
- [ ] **AC#5 与入向那一跳的依赖要具名**：`33-r1` 已把六环建好（`internal/panel/composer_dispatch.go`，台账 `A380`），但**生产调用者今天仍是零枚**（真窗口那一环 H2/H3 未建，且我派单里"要生产听众"与"禁改 `cmd/wisp`"互斥，矛盾在我）。⇒ 本票**不许拿"名册补齐"冒充"按钮能点"**；结题时要逐枚写明"这枚方法今天有没有真听众"。
- [ ] **AC#6 契约轴**：`docs/PLAN.md` 一字不动（**C17 的权威文字在 `PLAN.md:1367`**，改它＝另一次人工批准）；`docs/specs/**` 那张表**本票也不改**（owner 裁的是"代码向规格对齐"）；`allowlist.txt`、`thresholds.go`／golden、审批超时常量、`go.mod`／`go.sum` 一字节不动；`internal/panel/tokens_fourway_test.go`／`composer_test.go`／`l2_grant_boundary_test.go` 不动。
- [ ] **AC#7 门禁**：逐包 `go test -count=1 ./internal/panel/ ./internal/config/ ./internal/risk/ ./internal/tools/`（⚠ `internal/risk` **单包跑**＝`A359`；⚠ `./internal/panel/` 今天有 **3 枚在册红**，照实记不许当绿）；**收尾必跑一次全仓不变式面**（`A380` 那格账：`sh scripts/d22scan.sh`，基线 `ban #8 internal/`＝**438**；`gate-clauses.sh` **比红腿名册不比退码**，在册只 `G6neg`）；CLI 那一面要跑必须带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，**两形都贴**。

## 本票**不**解决

不做真窗口（H2/H3，归票 33）；不改 `docs/PLAN.md` 里 C17 的权威文字；不把面板做成"能直接改权限／能批准"；不替前端会话画界面（`frontend/**` 归别家）。

## Progress log

- 09-28 16:5x 编排者立票：来路＝`33-a1` 现量＋`Q-67` 由 owner 当场裁「按当前规格补，别搞债务了」（逐字 `A381`）。⚠ 我把他的裁定读成**以 `SPEC-08` 那张表为准、把代码侧补齐**（不是反过来改规格迁就代码），并把它拆成"先交逐枚代价表、再分批落地"两层——因为那 18 枚里有几枚今天连后端能力都没有，一次做完会把 AC#2 那三条禁区卷进来。未派。
