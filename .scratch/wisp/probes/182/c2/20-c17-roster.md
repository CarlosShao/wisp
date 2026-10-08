# 182-c2 · C17 名册真身与"兑现 vs 改契约"逐堆判定 — 锚 HEAD `a7e9b2e7`

## 1. 名册现量（本腿自己跑，⛔ 不按任何票面或台账的记忆数字）

| 面 | 逐字出处 | 现量 |
|---|---|---|
| `panel.*` 四枚 | `internal/panel/bridge.go:42-45`：`MethodModeRequest = "panel.mode.request"`／`MethodWorkspaceRequest = "panel.workspace.request"`／`MethodAttachmentAdd = "panel.attachment.add"`／`MethodMessageSend = "panel.message.send"` | **4** |
| 无前缀两枚 | `bridge.go:66` `MethodConfigGet = "config.get"`、`:67` `MethodConfigSet = "config.set"` | **2** |
| **入向名册真身** | 上两行相加 | **6 枚** |
| 判定函数 | `bridge.go:146 func knownComposerMethod(m string) bool`；拒口 `bridge.go:132 if !knownComposerMethod(r.Method) {` | — |
| 派发表 | `internal/panel/composer_dispatch.go:176 switch req.Method {` ＋ `:177/:182/:187/:192/:197/:202` 六枚 case ＋ `default: d.rosterMismatch(req)` | 6 case |
| 名册尺的钉 | `internal/panel/inbound_roster_253_test.go:59` `wantFullInboundRosterSize253 = 6`（断言体 `:445`/`:787`；DECLARED 信号定义 `:24-27`：`Method` 前缀 **或** 值是点分名，**前缀无关**） | 6→7 必红 |

### ⛔ 顶回一枚票面数字（不是本票的，是别家票面的，本腿不改）
**票 181:17 逐字"C17 白名单今天恰好四枚"——作"入向名册"计数是错的，作 `panel.*` 计数才是对的。**
那把尺（同行写明）＝`grep -n '"panel\.' internal/panel/bridge.go | grep -v _test`，它**结构上看不见**不带 `panel.` 前缀的 `config.get`/`config.set`；而无前缀是 `bridge.go:54-60` 逐字说明的**刻意**选择（"They carry no `panel.` prefix on purpose, and that is load-bearing twice over…"）。
⇒ 引这一句今后必须写两格：**`panel.*` 四枚 ＋ `config.*` 两枚 ＝ 入向真身六枚**。本腿只指认，⛔ 不替 181 改字。

## 2. 名册的**权威规格面**不在 `PLAN.md:1367`（派单那句要补正）

- `PLAN.md:1367` 逐字（本腿整行读，未截断）："**C17** | **`PanelBridge`** | 前端↔Go 双向通道：`invoke(method, args) → result` + Go→前端事件推送（流式结果、审批请求、任务状态）。**回复必须按 correlationId 路由**；前端**必须无状态**（WebView 销毁后一切从 Go 侧重读）| D29, D31" ⇒ **派单那句"PLAN.md:1367 已把 correlationId 路由写进 C17"成立**（这正是要防的"以为要改契约、其实是兑现"那一形的一枚实例）。
- 但**方法名册那张表在 `docs/specs/SPEC-08-ui-ball-panel.md:161-174`**。现量尺（结构尺，不是词频尺）：`awk 'NR>=163 && NR<=173' docs/specs/SPEC-08-ui-ball-panel.md | grep -oE '`[a-z]+\.[a-z]+`' | sort -u | wc -l` ⇒ **18 枚**，逐名：`approval.current` `approval.decide` `approval.queue` `config.get` `config.set` `cost.summary` `diagnostics.export` `grants.list` `grants.revoke` `history.query` `models.delete` `models.list` `panel.resync` `privacy.export` `privacy.purge` `task.detail` `tasks.list` `transcript.get`；另 `:174` 一行事件推送 5 枚（`approval.request` `ball.state` `cost.tick` `task.delta` `tool.chip`）⇒ **表内共 23 枚点分名**。
- `PLAN.md:2968-2973`（D39 C17 补四项）逐字：①方法白名单（未列出→拒绝并记日志）②每方法标注 capability 与是否需原生侧二次授权 ③correlationId 路由＋推送背压与合并 ④**面板每次 `show` 必须发 `panel.resync`，由 Go 侧全量推状态**。
- ⚠ 一条纪律：判"要不要人工批准"要**同时**对三张面——`bridge.go` 的 6 枚真身／`SPEC-08` 的 18+5 枚名册／票 194（"两块名册互不相识，owner 裁：代码向规格对齐"，7 枚框全未勾）。**今天只有 2 枚真名（`config.get`/`config.set`）在两张面上都对得上**。

## 3. 逐堆判定：兑现（规格早写了）还是改契约（要人工批准）

尺（本腿现跑，⛔ 不是只搜符号名）：`grep -noE '`(shell|terminal|file|fs|browser|commit|push|attachment|plugin|worktree|task)[a-z.]*`' docs/specs/SPEC-08-ui-ball-panel.md` ⇒ 全文件只有 **4 行**命中：`:109 fs.write`（那是别的节）、`:164 tasks.list`/`task.detail`、`:174 task.delta` ⇒ **SPEC-08 那张名册里没有 `shell.*`／`terminal.*`／文件树／浏览器／commit／push／attachment／plugin／worktree 任何一枚**。

| 堆／雷（182 的栏） | 若要做成"面板取数或动作"，需要的名字 | 在 `SPEC-08:163-174`？ | 判定 | 盘上凭据 |
|---|---|---|---|---|
| K1 token／成本 | `cost.summary`（invoke）／`cost.tick`（推送） | **在** | **兑现**，不是改契约 | `SPEC-08:171`/`:174`；载体归票 145 行 12 |
| K2 活动细维 | `history.query`／`transcript.get`／`tool.chip` | **在** | **兑现** | `SPEC-08:165`/`:174` |
| K3 审批**序号**／队列 | `approval.queue`／`approval.current` | **在** | **兑现**（⚠ 这是本腿新指认：票 167 那格"序号"今天缺的是**载体**，不缺**契约名额**） | `SPEC-08:166`；断点现量＝`pump.go:77-91 NativeVerdict` 无 Position 字段、真源在 `pending_read.go:58/:127` |
| K3 批准动作（雷 1） | `approval.decide` | 在 | **名额有、allow 侧被双重禁死** | `SPEC-08:167` 逐字"「allow」拒绝一切面板来源（F2）；仅 `reject` 可面板发起"＋票 194 AC#2①＋`ui.go:133`/`gate.go:746` |
| K8／K9 子代理与后台任务（显示） | 不需要新名（今天走快照） | — | **零名**（票 197 已落：`run.go:724 Tasks: rt.taskRosterState`） | 若改成"面板主动 invoke 取任务详情"＝`tasks.list`/`task.detail` 早列 ⇒ **兑现** |
| K9／K3 **停止** | `task.stop` 之类 | **不在** | **改契约＝人工批准**（本腿现量：SPEC-08 全文件没有一枚 stop/terminate 名） | 归票 167；167-c2 那句"入向 6 枚无一与停止有关"本腿复量成立 |
| K4 文件树（乙形一层列表） | 不需要新名 | — | **零名**＝不触 C17 | `A699` §5 逐字"⛔ 不需要新面板方法名＝不触 C17＝不欠机主那句话"，⚠ **条件**：宿主只给有限深度一层、展开由前端纯 DOM 做 |
| K4 文件树（**懒展开**） | 要一枚"点了才向宿主要下一层"的新名 | **不在** | **改契约＝人工批准**（A699 同一句逐字条件反面） | `A699` §5"若改成'点了才向宿主要下一层'的懒展开，就**必须新增 `panel.*` 方法名＝契约面＝要人工批准**" |
| K6 diff／未提交枚数 | 不需要新名 | — | **零名**＝不触 C17 | `A699` §1／票 189:26（189-a1 二轮现量"3～4 枚文件／约 355～475 行，⛔ 不需要新面板方法名"） |
| K6／K2 面板 commit／push（雷 4） | `git.commit`／`git.push` 之类 | **不在** | **改契约＝人工批准**；⚠ 且票 189 AC#3 逐字禁在显示票里出现动作腿 | `189:14`/`189:20` |
| 切分支／切工作树（票 186） | 名册 +2 枚 | **不在** | **契约变更，但批准已在盘上**＝`Q-64` 改判＝做（owner 原话入 `A368`） | 票 181:46 逐字"不再是拦路条件…收进新立的票 186"⇒ ⛔ 谁都不许再拿"要人工批准"挡 186 |
| K5 插件列表（读面） | 走快照 ⇒ 零名 | — | **零名**（同 git 显示那一形） | 票面 `:13` 那条"D34/宿主读面 vs 模型工具"之分；`internal/plugin` 现量只有 disposal 一族 |
| K5 插件**开关**（雷 2） | 不动名册，要解冻 `plugins.` 键集 | — | **不属 C17 那道题**（属 D36/C28 设置族） | `internal/panel/config_handlers.go:108` 逐字 `lockedFieldFamilies = []string{"risk.", "fs.", "net.", "plugins.", "privacy.", "models.", "audio."}` |
| K7 内置终端（雷 5） | `shell.*`／`terminal.*` | **不在** | **改契约**＝票 191 **乙支未批**；甲支（走同一套 risk 门控）已批 `A379` | 票 191:1（标题）＋`A699`／c1 §4 第 5 枚 |
| K10 浏览器／产物预览（票 192） | `browser.*` 之类 | **不在** | **改契约**；且 192 AC#3 逐字"本票**不做**表单提交、登录、下载、执行页面脚本；若产品要'浏览器能操作页面'，那是**新一级权限**，要单独立票＋owner 批" | 票 192:16 |
| 附件移除（雷 6） | `attachment.remove` | **不在**（SPEC-08 表无 attachment 任何一行） | **改契约＝人工批准**，且全仓零认领 | 尺 `grep -rn 'attachment.remove' --include=*.go --include=*.md internal docs/specs .scratch/wisp/issues` ⇒ **只命中票 182:90 自己**（rc=0） |
| `panel.resync`（AC#7/④ 相关） | 名册已有名 | **在**（`:163`，方向逐字"Go→前端推送"） | **兑现**，⚠ **但有陷阱**：若落地时把它铸成 `bridge.go` 的 `Method*` 常量，就会落进 DECLARED 信号 ⇒ 名册 6→7＝动那枚刻意钉住的数＝要人工批准。出处＝`inbound_roster_253_test.go:24-27`；同一陷阱 `35-a6` §④ 已具名 |

### 小结（给编排者的三句话）
1. **票 182 自己六枚框：零枚新 C17 面**（它是只读普查票，票面 `:29` 逐字禁 `internal/**`）。
2. **这一栏里"看起来要契约批准"的题，有三分之一其实规格早写了**：`cost.summary`／`history.query`＋`transcript.get`／`approval.queue`＋`approval.current`／`tasks.list`＋`task.detail`／`panel.resync` ⇒ 那五族是**兑现**（出处全部具名在上表），⛔ 不该拿去摆机主。
3. **真欠人工批准的是这四枚**：停止（`task.stop` 形）、commit/push、附件移除、终端乙支／浏览器操作——`SPEC-08` 那张表里**一个名字都没有**；其中"面板 allow"那一枚虽有名额却被 `SPEC-08:167`＋票 194 AC#2① 双重禁死，⛔ 不许当"待拍板"重摆。

## 4. 本件尺的 rc 自陈（⛔ 每把尺前面都没挂管道）
`grep -n 'knownComposerMethod\|MethodModeRequest =\|MethodConfigSet =' internal/panel/bridge.go` rc=0｜`grep -rn 'wantFullInboundRosterSize253' internal/panel/inbound_roster_253_test.go` rc=0｜`grep -n 'ErrPanelAllow' internal/agent/approval/ui.go internal/agent/approval/gate.go` rc=0｜`grep -rn 'plugin' internal/panel/config_handlers.go` rc=0（唯一命中 `:108`）｜`awk 'NR>=163 && NR<=173' … | grep -oE … | sort -u | wc -l` rc=0（**18**）｜`grep -noE '`(shell|terminal|file|fs|browser|commit|push|attachment|plugin|worktree|task)[a-z.]*`' docs/specs/SPEC-08-ui-ball-panel.md` rc=0｜`grep -rn 'attachment.remove' --include=*.go --include=*.md internal docs/specs .scratch/wisp/issues` rc=0（1 命中＝票 182:90）
