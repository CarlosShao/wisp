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
- 09-28 16:2x 只读普查程 `194-c1` 交代价表（**AC 框一枚未勾，产码／`go.mod`／`docs/PLAN.md`／`docs/specs/**` 零字节不动**）：表＝`docs/evidence/s1/194-method-roster-census-c1.md`，读数＝`.scratch/wisp/probes/194/c1/logs/`。锚 `8ca2291f`（未漂）；两枚 commit `a90908b2`／`72199c38`＋文档终枚，只 commit 未 push。
  - **枚数**：票面「18 枚」**对**（现数 `SPEC-08:163-173` 展开＝18 具名），加 `:174` 事件推送那一行 5 枚＝**23 具名**；代码那份现量 **4 枚**（`bridge.go:42-45`），交集**零枚**，⚠ 反向也零枚（那 4 枚规格未列 ⇒ 名册是并集问题不是替换问题）。
  - **三堆**：A 有货只差名字 **7**｜A/B 有数据面缺读取面 **4**｜B 后端真没货 **4**（`panel.resync`／`transcript.get`／`models.list`／`models.delete`）｜C 撞禁区要先裁 **3**。23 枚里产码同名命中只 3 枚，其中 2 枚是噪声 ⇒ **"有货"全靠按能力找同物异名**。
  - **撞禁区逐名（本程只标不裁）**：`approval.decide` 撞①（`SPEC-08:167` 逐字，allow 侧有结构性 nonce 闸 `queue.go:326/343/354`，reject 侧现成 `:393`）；`config.set` 撞③（`:168`）**且其唯一现成写手 `Manager.SetPermissionMode` `permmode.go:64` 天生撞②**；`grants.revoke` 待裁②（反方逐字证据 `queue.go:378-380`「fail-closed … which is why the panel may do it」）。⇒ 派单里那句"停手上报"在这一枚枚上成立：**没有任何一支能靠"按规格补"顺手动到这三条。**
  - **具名更正两处**：① 票 AC#2／派单 §2 的行号错位（allow 侧＝`167` 非 `169`；`config.set`＝`168` 非 `170`，各偏 2 行；引文本身逐字对得上）；② 票面称 `privacy.purge`「今天连后端能力都还没有」——**存储层已实现**（`PurgePrivacy` `privacy.go:170` 等 7 枚），缺的是面板路由与确认门。
  - **事件推送五枚**：`task.delta`／`approval.request`＝**有内容没这个名字**（整份 `Snapshot` 的两段 `composer.go:57-62`）；`ball.state`＝内容在**另一条原生管**（`Ball.SetState` `ball_windows.go:306`）；`tool.chip`／`cost.tick`＝**真没内容**（数据面在 `dao_misc.go:169`，从未进快照）。管子真在跑（`cmd/wisp/panel_pump.go:264`）。
  - ⚠ **别拿名册补齐冒充按钮能点**：`ComposerDispatch.Handle`（`composer_dispatch.go:120`）生产调用者**现量仍零枚** ⇒ 18 枚**无一例外**要等票 33 的 H2/H3。另有未定稿前提：这张表自己的标题写着「【SPEC 提案，**S5 定稿走契约批准**】」（`SPEC-08:156`）。
  - **门禁四数（`A363`：不充当结案凭据）**：`ban #8 internal/`＝**438**（与基线同值＝零漂移，`cmd/`＝45）；`gate-clauses.sh` 红腿名册＝**只 `G6neg`**（与在册一致，但枚数 1→3，本程不动尺不修腿）；`./internal/panel/`＝**3 枚红逐名照实记**（未当绿未修）；写面闸门＝**终态空集＝起手空集，差集为零**。工具调用终值 **32／硬顶 35**（本程中途把自报累计数报高过，已具名更正）；零删除自证通过（树里 `design/**` 的 ` D` 条目非本程所为、未读未引）。
  - 09-28 16:2x **终枚更正一枚（改的是本程自己刚写的那句闸门，不动编排者原句）**：`16:22` 那次 `git status --porcelain -- internal/ cmd/`＝空集，`16:26` 重取时**多出两枚别人的在飞件**（`internal/risk/pointer_185_test.go`、`internal/tools/pointer_185_cli_seam_test.go`，文件名属 `185-r1`）⇒ 按 `A374` **具名登记为"不是我的"、不动、不提交、不还原**；本程三枚 commit 逐枚 `git show --name-status` 自证＝共 8 条文件条目、**零枚 Go**（`internal/**`／`cmd/**`／`go.mod`／`docs/PLAN.md`／`docs/specs/**` 一字节未动）。
  - ⚠ 另记一枚锚点限定：起手现量 HEAD＝`8ca2291f` 属实（未漂），但本程 `16:14` 落第一枚 commit 时线上已被别家推进两枚（`87e810c4` ledger A381、`423cac66` 181-r2）⇒ **本程首枚的父提交是 `87e810c4`**；⇒ 共树里**别用区间 diff 自证只读程的干净**（`git diff --name-only 8ca2291f..HEAD` 会把别人的 `cmd/wisp/panel_pump.go`、`internal/panel/git.go` 算进来），一律逐枚 show。

## 编排者收表（09-28 16:4x，账 `A384`）：复跑对得上，三枚待裁我裁完，两处更正入账
- **我复跑并认下的读数**：规格那份逐行展开＝`:163` 1 枚＋`:164-173` 17 枚＝**18 具名**（票面那个数对），`:174` 另 5 枚事件；代码那份＝`bridge.go:42-45` **4 枚**；`ComposerDispatch.Handle` 入向生产调用者＝**零枚**（只命中它自己的定义与注释）。A 堆七枚函数逐名现量到 `file:line`（`dao_tasklog.go:93`／`privacy.go:170`、`:198`／`dao_misc.go:78`、`:169`／`diagnostics.go:62`），B 堆那两枚也真缺（`internal/models/` 只有私有 `removeStaging` `downloader.go:739`、无导出列举面）。
- **裁定① `grants.revoke` 由面板发起＝放行，不撞禁区**。依据是冻结文本两处逐字：`docs/specs/SPEC-06-security-gatekeeping.md:117-118`「面板「安全」页实时列出并可**一键撤销**」＋`:19`「「允许」只接受原生侧来源……面板只能「拒绝/查看」」；方向与 `internal/agent/approval/queue.go:378-380` 那句注释（拒绝/撤销＝fail-closed）一致。**边界**：只准撤销**已存在的作用域授权**；经这个口发放或延长授权一律不算。**这一枚可以进堆1 一起派。**
- **裁定② `config.set` 今天不派**：它要的"按节写入＋安全节放宽走 L2 重新确认"（`SPEC-08:168`＋`SPEC-03 §4.2`）**地基不存在**，而全仓唯一现成写手是 `internal/config/permmode.go:64` 的 `SetPermissionMode`——拿它当捷径接出来的正好是"面板直接改权限状态"那一形。**已另立票 195 先建地基**，本票堆3 只留 `approval.decide` 的 reject 一支。
- **裁定③ 派腿顺序＝先片 B（生产听众）后名册**：`Handle` 零调用方期间，补齐 18 枚只会交付"门牌挂满、一枚点不动"＝`panel.resync` 那病的放大版。⇒ `33-r2` 已派，堆1 排在它之后；**堆1 的写面范围照它给的那份**（`bridge.go` 只**追加**常量、不动 42-45 的顺序与命名，避免撞 `composer_test.go:502` 那枚钉）。
- **两处更正（改的是本票票面与我派单，不动 `194-c1` 的原句）**：① AC#2／派单 §2 的 `SPEC-08` 行号各偏 2 行——allow 侧＝**`:167`**、`config.set`＝**`:168`**（引文本身逐字对得上）；② 票面原写"`privacy.purge` 今天连后端能力都没有"**错**，存储层已实现（`internal/memory/privacy.go:170`），缺的只是面板路由与确认门。
- **另两处口径**：① `runtests.sh: PASS=34 FAIL=0` 与 `internal/panel/` 那 3 枚红**不打脸**——那句是 `scripts/d22scan.sh:51` 用 `-C tools/d22scan` 起的正控，`./...` 只覆盖 scanner 自己那一包（本票登记为"矛盾"那条按此更正）；② `ban #8 internal/` 的 **438 是"扫到的文件枚数"、不是违规枚数**，`16:4x` 我复跑＝**441**（多 3 枚＝在飞 `185-r1` 新建的文件）、门仍 **clean** ⇒ 本票与台账今后只报"门 clean/红＋红点逐名"。
- **事件推送五枚的一条硬约束收下并入本票**：不许把它们塞进 `bridge.go` 的入向白名单常量（入向方法名与出向事件名混成一张表＝本票那枚病的镜像形状）。
- **本票状态**：`194-c1` 只是**代价表**，AC 七格**一格不勾**（勾要等非实现者裁"补齐之后代码与规格是否真同源"）。下一步＝`33-r2` 交件后派**堆1 写腿**。

## 编排者派：过期读数更正（腿 7b2，10-06）

> 来路＝只读＋追加式更正腿 `expired-premises-7b2`（锚 HEAD＝`bd39f17f`，10-06 19:2x 自取）。⛔ **本节不改动上面任何一行**，尤其不碰 `:42` 那句（它在 `- [ ] **AC#5 …**` 的框文射程内，改它＝改判据，判据归编排者翻）；AC 七格**一格没勾、一格没改**（本节追加前后 `grep -c '^- \[ \]'` 都＝**7**，我量过）。
> 取料与尺的全文在 `.scratch/wisp/probes/expired-premises/7b2/verdict.md`；本节只落"哪两条链推翻了哪一句"。

**① 被推翻的句子（逐字，`:42`，本腿未改它一字）**：
`⚠ **别拿名册补齐冒充按钮能点**：` + 「`ComposerDispatch.Handle`（`composer_dispatch.go:120`）生产调用者**现量仍零枚** ⇒ 18 枚**无一例外**要等票 33 的 H2/H3。」
⇒ 其中"**生产调用者现量仍零枚**"这一层，在今天这棵树上**不成立**；票 186 普查注 `:54` 里那句"入向那一跳**整条不存在**"同属这一层（见 verdict §3，那一枚票面本腿一字未动）。

**② 尺（锚在带括号的调用形状上，⛔ 不锚符号名）**——命令原文照抄，可复跑：
```
git --no-pager grep -n 'disp\.Handle(' HEAD -- cmd internal | grep -v '_test.go'
```
我现量＝**2 处**：`cmd/wisp/panel_host_windows.go:637`＋`cmd/wisp/panel_inbound.go:163`。
⚠ 这把尺的两处陷阱（本腿都现量过）：不剥 `_test.go` ⇒ 变成 9 行、全测试；改扫 `Handle(ctx` 宽形状 ⇒ 变成 9 处，其中 `cmd/wisp/logsink.go:209`、`:211` 与 `internal/observe/logging.go:138`、`:505`、`:519`、`:574`、`:576` 共 **7 处是 slog handler 的同名 `Handle` 方法，与 `ComposerDispatch` 无关**。⇒ **"有人构造了这个对象" ≠ "有人调用了它的方法"**，判"零调用者"与推翻"零调用者"必须用同一把锚在调用形状上、双向都追到 `main` 的尺。

**③ 推翻它的链 A（常驻 GUI 腿，逐跳本腿现量）**：
`cmd/wisp/main.go:66` `runResident()`（注释 `:135` 自证它住在 `resident_windows.go`）
→ `cmd/wisp/resident_windows.go:33` `func runResident()` → `:151` `rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)`
→ `cmd/wisp/panel_resident_windows.go:228` `func newResidentPanelManager(...)` → `:234` `disp, err := newResidentComposerDispatch(dataDir, auditf)` → `:211` 该函数 → `:215` `return newComposerDispatchChain(dataDir, auditf, residentPanelActor)`（`cmd/wisp/panel_inbound.go:228`/`:271` 装配，`:248` `modeWrites := &panel.ModeWriteHandler{` ⇒ 路由器**真接了处理器**，不是空壳）
→ `cmd/wisp/panel_resident_windows.go:253` `return NewPanelManager(disp, assets, dataPath, …)`
→ `cmd/wisp/panel_host_windows.go:213` `func NewPanelManager(disp *panel.ComposerDispatch, …)` ＋ `:152` 字段 `disp *panel.ComposerDispatch`（**归属证明**：`:637` 调的确实是这一枚类型的方法）
→ 到页那一支：`cmd/wisp/resident_windows.go:218` `return panel.RequestToggle(via)`（在 `startResidentBall(... withPanelHost(...))` 里）→ `cmd/wisp/panel_resident_windows.go:427` `RequestToggle` → `:436` `rp.RequestShow(via)` → `:383` `rp.post(func() { rp.showOnThread(via) })` → `:409` `err := rp.mgr.Show(context.Background())` → `cmd/wisp/panel_host_windows.go:488` `if err := m.bringUp(ctx); err != nil`（`:318` 定义）→ **`:405` `bindErr := w.Bind(panelDispatchBinding, func(raw string) string {`** → `:406` `reply, _ := m.dispatchRaw(ctx, raw)` → `:630` `dispatchRaw` → **`:637` `return m.disp.Handle(ctx, raw)`**。

**④ 推翻它的链 B（CLI 诊断腿，逐跳本腿现量）**：
`cmd/wisp/main.go:115` `case "panel-inbound":` → `:120` `os.Exit(cmdPanelInbound(args[1:], panelInboundIO{}))`
→ `cmd/wisp/panel_inbound.go:103` `func cmdPanelInbound(args []string, s panelInboundIO) int` → `:146` `disp, err := newPanelInboundDispatch(dir, auditf)` → `:209` 该函数返回 `*panel.ComposerDispatch` → `:163` `reply, err := disp.Handle(ctx, raw)`。

**⑤ 推翻"整条不存在"的字面尺（票 186 `:54` 给的那把 `grep WebMessage|ReceiveMessage|OnMessage`）**：本腿把同一族形状打到**依赖库源码**上（只读，`D:/work/base/gopath/pkg/mod/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808/`）：
- `pkg/edge/chromium.go:112` `e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")`（页面侧 `window.<名字>(...)` 被改写成 `postMessage`）
- `pkg/edge/chromium.go:201` `e.webview.vtbl.AddWebMessageReceived.Call(...)` ＋ `:233` `func (e *Chromium) MessageReceived(...)` → `:240` `e.MessageCallback(...)`
- `webview.go:103` `chromium.MessageCallback = w.msgcb` → `:139` `msgcb` → `:147` `w.callbinding(d)` → `:164` `f, ok := w.bindings[d.Method]` → 命中 `panel_host_windows.go:405` 绑进去的那枚闭包。
⇒ 票 186 `:54` 那把尺**扫的是本仓 Go 文件**，而这条跳的接收器住在依赖库里；"全仓非 test 没有一枚 WebView2 消息接收器"在**这一枚符号上今天不成立**。

**⑥ ★ 最关键的那句边界（本节不许被读成什么）**：
**过期的是"枚数"这一层，不是那句警告的实质**——"页面点下去到不到 Go"是**运行期**问题，本机 `winlive` 未批、真窗那一发今天量不了 ⇒ ⛔ **本节不得被读成"按钮已能点"**。静态可达（链 A／B 每跳都在）与"owner 真按下去收到回执"是两件事；票 33 `:42` AC#14 登记的"Go→页面那一跳今天没有人投递"（模块 `webview.go:443-448` 的 `Dispatch` 只入队＋`PostThreadMessageW`）本腿**未推翻、也不归本腿推翻**。⇒ `:42` 里"18 枚无一例外要等票 33 的 H2/H3"这半句**照旧成立**，本腿只把"生产调用者现量仍零枚"标成过期读数。

**⑦ 处置与两处具名分歧**：
- 票 194 只**追加**本节（本节以上任何一行一字未改，`:25` AC#5 框文里那句"生产调用者今天仍是零枚"也**没动**——翻不翻由编排者定，见 verdict §4-①）；`docs/evidence/s1/**` 一字未动；票 33 `:182`／`:188` 不在本批射程（答的是"具名诊断入口"那一枚），未动。
- **票 186**：普查腿 `pool-validity-4f` §5 第 3 条指认的"零枚"那句，在 186 票面上**不是这个字面**——`:54` 有同族句子（「入向那一跳**整条不存在**」＋「`ParseComposerRequest` 与 `HandleModeRequest` 非 test **零调用方**」），但**没有**"生产调用者零枚"六字逐字。⇒ 4f 的**内容**指认成立、**字面**落点不成立；本腿按派单纪律⛔ **未动票 186 一个字**（缺的是授权范围那一行，见 verdict §4-②）。
- **票 219**：`grep -n 'ComposerDispatch'` 于票 219＝**零命中**（rc=1）⇒ **4f 报的"票 219 表第 4 行"那一处落点在票面上不存在**（`:39` 表第 4 行答的是 `DecideFromPanel`／`DecideFromNative`／`Gate.Veto`，与 `Handle` 不是同一件事）。⛔ 未动票 219 一个字。
- ⚠ 顺带报单（本批射程外，本腿未动）：票 219 `:39`／`:15` 那两枚"答复侧生产零调用者"**今天也过期**（`internal/agent/approval/replies.go:328`、`:408`、`:413`、`:437`、`:463` 五处非 test 调用者，驱动者＝生产文件 `cmd/wisp/approval_reply.go:215`、`:259`、`:302`、`:304`、`:331` 与 `cmd/wisp/resident_approval_windows.go:612`、`:728`；票面引的 `gate.go:622/:610/:371` 今天＝`:736/:724/:421`）。⇒ 是否另派一枚腿由编排者定。
- ⚠ 锚点：本节起手锚＝`bd39f17f`（19:2x），交件前线上已被推进到 `76370fb5`（19:5x）⇒ 本腿把 §② 那把尺与链 A／B 的 20 枚锚**在新 HEAD 上逐枚复跑**，读数不变（详见 verdict §7）。
零删除自证：`git --no-pager diff --numstat -- .scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md` 我现量终值＝**`45  0 <全名>`**（删除列 **0**；行数 `55 → 100`；`grep -c '^- \[ \]'` 追加前后都＝**7**），读数与两次中间值（`40 0`）都抄在 verdict §6 与 §7。
