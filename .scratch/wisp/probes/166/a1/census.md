# 票 166 · a1 只读普查（Go 侧：会话的 list / search / rename / archive / fork）

- 子代理：`166-a1`（只读，无任何写产码动作）
- 锚点（`date -Iseconds` ＋ `git log -1 --format=%H` **自取**，不引别人的号）：
  - 读数时刻 `2026-10-03T10:59:09+08:00` 同刻 `git log -1` = `70a935ce3709d0f1b42f095a7462b2808b7b7ab9`（本件全部行号对齐这一枚之后的工作树）
  - 本程开工时刻同命令 = `1a14f203b8ae42cb61178feae9969b48c6677649`（两枚之间差一枚 commit，期间有写腿落地 ⇒ 锚点仍在动，复量者请自取当刻号）
- 树状态：**不是干净树**（`git status --short` 有他腿在飞的条目，含 `cmd/wisp/panel_host_windows.go`、`cmd/wisp/panel_resident_windows.go`）。本件所有行号取自**工作树现读**，不是取自 `git show <锚点>`；复量者若按锚点取文，可能漂。
- 射程：`cmd internal tools docs scripts` 四枚显式根（+ `.scratch/wisp/issues`、`.scratch/wisp/probes` 用于票面与台账对账）。`frontend/**`、`design/**` **未读、未引**（两层禁令）。
- 本程**未跑任何 Go 命令**（无 `go test`/`build`/`vet`/`list`、无 `wisp`）。所有结论来自读文件与 grep；必须跑了才知道的格子全部进「G. 量不到」并写明复量法，未做任何猜测。

---

## 0. 一句话结论

票面标题级判断（"列不出来、搜不到、改不了名、断不了叉"）**成立**，但**成立的方式比票面写的更硬**：不是"有出口没接线"，而是**盘上根本没有"会话"这一枚实体**——`D35` 那九张表里没有一张存对话、存消息、存标题；会话内容今天只活在 `agent.Loop` 的一块内存切片里，进程退出即蒸发。五件事里**只有 `list` 有一枚真实的 SQL 出口**（`memory.Store.ListTaskLogs`），且它**生产调用者＝0**。`search`／`rename`／`archive`／`fork` 四枚**存储层那一跳都不存在**。

另有一枚**票面没预料的硬冲突**：AC#3 要的"加列迁移"会直接撞 `internal/memory/schema_test.go:173` 那族 schema 契约钉（它把表枚数、逐表列数、DDL 语句枚数、`schema_version` 值四件事全钉死），而 AC#5 又要求 `docs/specs/**` 零字节——`schema.go:10-13` 自陈 DDL 必须与 `SPEC-02 §3` **逐字节相同**。**这两格今天互相拉扯，属"未定义即停"，见 H 节。**

---

## A. 五件事的现量名册

判"接上了没有"的规矩照派单：**数显式调用者枚数答不了能力问题**，逐枚问三处——**谁产出／谁投递／谁落盘**。`go build`/`go vet` rc=0 在本程一次都没跑，也证明不了任何接线（本件不拿它当凭据）。

### A0 先把"会话"在本仓的三个候选实体钉死（不钉这一格，后面每一问都会漂）

| 名字 | 真身 | 盘上有没有 | 键 |
|---|---|---|---|
| 对话/会话（用户意义上的"那次聊天"） | **无实体**。只有 `agent.Loop.history []llm.Message`（`internal/agent/loop.go:187`，注释 `:175` 逐字"single-task agent loop holding one conversation history"） | **无**（不落盘，进程结束即无） | 无键 |
| 任务（一次 `wisp run` / 一枚子代理） | `task_log` 表（`internal/memory/schema.go:53-65`） | 有 | `task_log.id`＝UUIDv4 形状，`newTaskID()` 铸（`internal/agent/loop.go:1121-1123`） |
| 宿主安全会话（审批用的那个"本次会话"） | `internal/session.ID`（`sess_` + 32 hex，`internal/session/session.go:47,94`），落 `approval_grant.session_id` | 有 | `session.ID.String()`，**随进程死**（`:19-20` 逐字） |

⇒ 票 166 要的"会话管理"在盘上**没有主键对象**。任何落地腿第一件事是**选一枚已有键当会话身份，或新造一枚**——这是契约动作，不是接线动作。

### A1 list（列出所有会话）

| 四问 | 读数 | 尺（可复制，均为读文件/grep） |
|---|---|---|
| 存储层有没有那一跳 | **有**（但列的是任务不是会话）：`internal/memory/dao_tasklog.go:93` `ListTaskLogs`，SQL `ORDER BY started_at DESC, id DESC`（`:97`） | 读 `dao_tasklog.go:92-111` |
| Go 有没有导出方法 | **有**，且另有一枚更宽的：`internal/memory/privacy.go:47` `ListPrivacy`（对 `profile/memory/task_log/tool_call/artifacts` 五域统一列） | `grep -n "^func (s \*Store)" internal/memory/*.go` |
| 有没有生产调用者 | `ListTaskLogs` 非测试引用行 **4 枚**（注释 1 + 定义 1 + 调用 2），两处调用**都在 `internal/memory` 自己家里**（`privacy.go:76`、`privacy.go:214`）⇒ **包外生产调用者＝0**。`ListPrivacy` 非测试引用 **2 枚**（注释＋定义），**调用者＝0**；命中只在 `internal/memory/privacy_test.go`、`internal/memory/artifacts_path_invariant_test.go:412` | `grep -rn "ListTaskLogs" cmd internal tools scripts --include='*.go' \| grep -v "_test.go" \| wc -l` ⇒ 4（逐行核为 1 注释 + 1 定义 + 2 调用） |
| 面板/CLI 能不能触发 | **都不能**。CLI：`cmd/wisp/main.go:88-132` 的 switch 共 **10 枚命令**（run/providers/doctor/secret/models/slo/panel-assets/panel-inbound/version/help）＋default 拒绝，**无一枚列历史**。面板：快照 `Snapshot` 只有 `pending/results/composer/generatedAt`（+两枚可选 `instructions/tasks`），`internal/panel/composer.go:57-91`，**没有会话名册字段**；入向六枚方法名册（`internal/panel/bridge.go:42-45,66-67`）也**没有列会话那一枚** | 读 `main.go:88-132`；读 `bridge.go:42-67`；读 `composer.go:57-91` |

三处各是谁：**产出＝`memory.Store.ListTaskLogs`（存在）**；**投递＝无人**（`internal/memory` 包外零调用者）；**落盘（到用户眼前）＝无处**（CLI 不印、快照不带）。
⇒ 票面 AC#1 的答案：**量不到"能列出历史会话"的任何一发** ⇒ **本票成立**，不是编排者判断错。

### A2 search（按内容或标题找）

| 四问 | 读数 |
|---|---|
| 存储层有没有那一跳 | **对会话没有**。四枚显式根里唯一一枚真做检索的 SQL 是 `internal/memory/dao_memory.go:52` `SearchMemory`——它搜的是 **L2 显式记忆表 `memory`（`keywords` 列，`schema.go:40-50`）**，`LIKE … ESCAPE`（转义助手 `models.go:215 escapeLike` 已在，可复用）。`task_log` 侧**零枚** LIKE/全文：`FROM task_log` 只出现在 `dao_tasklog.go:82,97,116,132` 四枚（by-id / list / delete / purge）。**没有 FTS 虚表**（`grep -rn "CREATE VIRTUAL TABLE" cmd internal tools scripts --include='*.go' | wc -l` ⇒ **0**） |
| Go 有没有导出方法 | 有 `SearchMemory`（但那不是会话搜索） |
| 有没有生产调用者 | **0 枚**。`grep -rn "SearchMemory" cmd internal tools scripts --include='*.go' \| grep -v _test` ⇒ **2 行，逐行为 `dao_memory.go:47` 注释与 `:52` 定义**，无一行调用。命中只在 `internal/memory/dao_test.go` |
| 面板/CLI 能不能触发 | 都不能。D34 声明的 `search.content`/`search.files`/`memory.recall` 今天**无实现**（`grep -rn "MemorySearch\|Recall" --include='*.go' \| grep -v _test` ⇒ 0 命中）；台账 `A308`（`docs/reports/pending-and-issues.md:7658`）自陈那 30 枚"声明了但不可及" |

三处：**产出＝只有 L2 那一跳（且它自己也没人调）**；**投递＝无**；**落盘＝无**。

### A3 rename（改标题）

| 四问 | 读数 |
|---|---|
| 存储层有没有那一跳 | **没有"标题"这一枚列**。`task_log` 十二枚列见 `schema.go:53-65` 与 `models.go:58-70`：`id/started_at/ended_at/state/query_text/summary_text/cost_tokens_in/cost_tokens_out/cost_amount_micro/currency/error_class`。**没有 title/name/label**。最接近的两枚都不是标题：`query_text` 是**已脱敏的用户原话**（`schema.go:58` 注释 + §14.4），`summary_text` 是**模型产出的收尾摘要**，由 `FinishTaskLog`（`dao_tasklog.go:49-74`，SQL `:58-61`）在任务结束时**一次性写**，生产写点＝`cmd/wisp/run.go:1135` |
| Go 有没有导出方法 | **没有改名方法**：`internal/memory` 全族 `func (s *Store)` 名册里无任何 `Rename*/SetTitle*/UpdateTitle*`（尺：`grep -rn "^func (s \*Store)" internal/memory --include='*.go' \| grep -v _test`，逐枚看过） |
| 有没有生产调用者 | 不适用（无方法可调用） |
| 面板/CLI 能不能触发 | 都不能 |

三处：**产出＝无**；**投递＝无**；**落盘＝无**。⇒ 这一格是**从零造**，不是接线。⚠ 造的时候有两个后果必须选：(a) 写进 `summary_text` ＝ 把模型摘要与用户标题混进同一列，`FinishTaskLog` 的 `COALESCE(?, summary_text)`（`dao_tasklog.go:59`）会在下一次收尾时**覆盖用户的改名**；(b) 新增一列 ＝ 撞 schema 契约钉（见 D/E/H）。

### A4 archive（归档 / 取消归档，票面还带"撤销"）

| 四问 | 读数 |
|---|---|
| 存储层有没有那一跳 | **没有**。今天唯一像"归档"的两个东西都不是归档：① `state` 列（`schema.go:57`）是一枚**自由文本、无 CHECK**（`validateTaskLog` 只要求非空，`models.go:185-197`；`grep -rn "CHECK(" internal/memory --include='*.go' | wc -l` ⇒ **0**，`PRAGMA foreign_keys` 在 `cmd internal tools scripts` 源码里 ⇒ **0**，只以二进制字符串出现在编译产物 `build/*.exe` 里，那是驱动的自带串、不等于我们开了它）；② **保留期清理**（`internal/memory/retention.go:32 TaskLogTTL = 30*24h`，删行 SQL 在 `:194-204`：`DELETE FROM <table> WHERE <ts> IS NOT NULL AND <ts> < ?`） |
| Go 有没有导出方法 | **没有 Archive*/Unarchive*。⚠ 且"撤销"这一形在真相源只追加的规矩下今天无处可放（`DeleteTaskLog:114`/`PurgeTaskLogs:129` 都是**真删**，不可逆）** | 
| 有没有生产调用者 | `RetentionJob` **生产调用者＝0**：`grep -rn "RetentionJob" cmd internal tools scripts --include='*.go' \| grep -v "^internal/memory/"` ⇒ **1 行，且是注释**（`internal/observe/doc.go:20`）；`NewRetention*`/`StartRetention` 0 命中。⇒ 30 天窗口今天**没有任何东西在跑**（这与"归档会不会让它消失"直接相关，见 B） |
| 面板/CLI 能不能触发 | 都不能。隐私页那一族（`ListPrivacy/DeletePrivacyItem:135/PurgePrivacy:170/ExportPrivacy:198`）是**规格要求的归档邻居**（`docs/specs/SPEC-02:167-168`、`docs/PLAN.md:2715-2716` 规则 2 逐字要"一键清空/导出/逐条删除"），**四枚生产调用者全为 0** |

三处：**产出＝无**；**投递＝无**；**落盘＝无**。

### A5 fork（分叉出新会话）

| 四问 | 读数 |
|---|---|
| 存储层有没有那一跳 | **没有**。`grep -rniE "\bfork\b" cmd internal tools scripts --include='*.go' \| grep -v _test` 的命中**逐枚都不是会话分叉**：`cmd/wisp/doctor.go:226,290-297`、`cmd/wisp/secret.go:88,98,154`（数据根/环境**布局分叉**）、`cmd/wisp/resident_approval_windows.go:270`（注释指票 246-a1 的"分叉"）、`cmd/wisp/testdata/esclistener/main.go:179-183`（测试桩里的一枚 goroutine 起名 fork）、`internal/models/archive.go`（tar 解包，与"归档"同名不同物） |
| Go 有没有导出方法 | 无 |
| 有没有生产调用者 | 不适用 |
| 面板/CLI 能不能触发 | 都不能 |

三处：**产出＝无**；**投递＝无**；**落盘＝无**。⇒ 分叉要复制的是**历史本体**，而历史本体今天不落盘（A0），所以 **fork 排在 transcript 持久化之后，不能先行**；票面 AC#4（"分叉不许吞原件"）在现量上**没有可吞的原件**，这一格的判据要等落地腿才有对象。

### A6 会话内容今天到底散在哪几处（写腿要知道的"现状地板"）

| 载体 | 位置 | 是否落盘 | 边界 |
|---|---|---|---|
| 对话消息（`llm.Message` 序列） | `internal/agent/loop.go:187 history` | **否** | `History()` 只给内存副本（`:256-262`）；`Reset()` 清空（`:264-270`），且 **`Reset()` 全仓零生产调用者**（尺：`grep -rn "\.Reset()" cmd internal tools --include='*.go' \| grep -v _test` ⇒ 3 行，全在 `internal/llm/golden/replay.go:112`、`tools/d22scan/main.go:175,183`、`tools/mockllm/chat.go:431`，**没有一枚是 `Loop.Reset`**）⇒ 连"会话结束"这个动作都没人按 |
| 流式片段（面板看到的正文） | `internal/panel/pump.go:416 StreamLog`（`cmd/wisp/run.go:698` 装配，`DefaultStreamKeys` 上限） | **否**（内存，按 key 数封顶，溢出改名/截断） | 快照 `results` 只带 `{correlationId, text, done}`（`composer.go:95-99`） |
| 任务收尾摘要/成本 | `task_log.summary_text`/`cost_*`（`cmd/wisp/run.go:1124-1141`） | **是** | 每任务一行，非消息级 |
| 工具取证 | `tool_call`（`dao_toolcall.go`） | 是 | `correlation_id == task_id`（`loop.go:356` 逐字"C18: correlation == task id"） |
| 溢出／附件工件 | `internal/memory/artifacts.go:71 PutArtifact`——**生产调用者＝1 枚**：`internal/panel/attachments.go:264`（面板附件落工件，接口声明在 `:57-60`）。⚠ `internal/agent/spill.go`（`:84 Prepare`）**自己写盘、不经 `Store.PutArtifact`** ⇒ 两条落盘路是两套守卫，别混 | 是（文件） | 500MB LRU（`retention.go:37-38`，配额常量在 `ArtifactsQuotaBytes`），**不是会话** |
| 转写文本 | `task_log.query_text`，§14.4 掩码 | 是 | `docs/PLAN.md:1874` 逐字："L3 任务日志 / SQLite 历史：这是产品功能（用户要能看历史），**必须存**" ⇒ **会话历史落盘是规格要求的，不是可选项**；`[privacy] keep_transcript` 硬编码 false（`internal/config/schema.go:504-515`、`internal/config/validate.go:81-91`）管的是"逐字转写全文"那一档，见 H#3 |

---

## B. 会话身份

### B1 一枚会话在盘上用什么键认（今天的答案：**没有会话键**）

| 键 | 由谁铸 | 几枚生产写点 | 存在哪 | 生命周期 |
|---|---|---|---|---|
| `task_log.id`（UUIDv4 形状） | `internal/agent/loop.go:1123 newTaskID()`（`crypto/rand`，随机源坏了退回**时钟派生**，`:1125-1131`） | 写：`cmd/wisp/run.go:1128`（`StartTaskLog`）、`internal/agent/journal.go:129` | `task_log.id` PK | **进程内**；行留 30 天（且今天清理无人跑，A4） |
| `tool_call.correlation_id` | **不另铸**，直接等于 task id（`internal/agent/loop.go:356 newTaskJournal(l.opt.Journal, taskID, taskID)`，注释逐字"C18: correlation == task id"） | `journal.go:81` | `tool_call.correlation_id`（`idx_tool_call_corr`，`schema.go:86`） | 同上 |
| `approval_grant.session_id` | `internal/session.Mint()`（`session.go:94`，`crypto/rand` 128bit，`sess_` 前缀，`:47`） | **1 枚生产调用者**：`cmd/wisp/run.go:473`（下一行 `:478 NewLedger`）；测试 0 枚直调 | `approval_grant.session_id`（`idx_grant_session`，`schema.go:99`） | **随进程死**（`session.go:19-20` 逐字："结束点今天＝进程退出"） |
| 面板流式 key | `panel.StreamLog` 按 task id / `subagent:<taskID>` | `cmd/wisp/run.go:1266,1293` | 内存 | 进程内 |

⇒ 分工一句话：**task id 是"一次活"，correlation id 今天就是它的别名，session id 是"一次开机"**。三者**没有一枚能充当跨重启的会话身份**——`session_id` 甚至是**故意**不能跨重启（`session.go:16-30` 整节论证：若可重算，D45 的"本次会话内允许"会变成永久免审通行证）。

### B2 改名会不会动到那个键

**不会动，也不能动**——现量：

1. 键就是 `task_log.id`（TEXT PK），改名要写的列与它无关（A3）。
2. `tool_call.task_id TEXT NOT NULL REFERENCES task_log(id)`（`schema.go:71`）与 `idx_tool_call_task`（`:85`）都**指向它**。⚠ 但这条 FK **不生效**：`PRAGMA foreign_keys` 在源码四枚根里 **0 命中**（尺：`grep -rn "foreign_keys" cmd internal tools scripts docker build | wc -l` ⇒ 4 行，逐行都是 `Binary file build/wisp*.exe matches`，**没有一行是源码**；SQLite 默认不开外键）。⇒ 真去改 id，**不会拦、也不会级联**，`tool_call` 会留下指空的孤儿行。**这条必须写进落地腿的禁区**：动 id ＝ 静默断取证链，且没有数据库层护栏。
3. 成本账不键在 task 上，而是按天聚合：`cost_daily(day)`（`schema.go:102-108`，`dao_misc.go:132 BumpCostDay`）。⇒ **改名碰不到成本账**；碰得到成本账的是"任务的终账"（`FinishTaskLog` 一次写三枚数），不是名字。

### B3 归档会不会让它从名册里消失（票面要求"这两个后果必须写清"）

分三种"归档"形，逐形现量：

| 形 | 会不会从名册消失 | 依据 |
|---|---|---|
| **甲：`state` 列塞一个 `archived`** | **会从"当前任务名册"消失，且是静默的**。`ListTaskLogs`（`dao_tasklog.go:97`）**不带 WHERE**、列**全部**行；而 `RetentionJob` 的删行谓词是 `WHERE started_at < cutoff`（`retention.go:196-198`），**不看 state** ⇒ 归档**挡不住 30 天删除**，也**不会因归档而多活一天**。另撞一枚已登记的缺陷：状态词表今天有两套且互不相认（票 196 现量：真写进库的 `running/succeeded/cancelled/done` 一枚都不在 D43；schema 无 CHECK ⇒ **第五个词加进去门不会响**）。⇒ 甲形＝把"归档"塞进一枚**无人校验的枚举**，正是票 196 警告的形状 |
| **乙：新加 `archived_at INTEGER NULL` 一列** | **不会从名册消失**（`ListTaskLogs` 仍会带它，只要 SQL 不改），但**必红 schema 契约钉**（`schema_test.go:259` 逐表列数＝`contractSpec`；`:195` 表枚数＝9；`:229/:233` 语句枚数＝15/1；`:222` 版本值＝"2"）；且 `schema.go:10-13` 要求 DDL 与 `SPEC-02 §3` **逐字节相同**，而 AC#5 禁改 `docs/specs/**` |
| **丙：真删 / `purgeAllDomains`** | **永久消失且无撤销**（`DeleteTaskLog:114`、`PurgeTaskLogs:129`、`privacy.go:262`），与票面硬约束 3"不许静默删历史"、`issues/README` 真相源只追加**直接相冲** ⇒ 丙形**不能当归档用** |

⇒ 对审批队列/成本账那一族的连带判断：**只要归档选甲/乙且不动 id，`tool_call` 取证行、`approval_grant`（它键在 session id 上、与任务无关）、`cost_daily`（按天）全都不受影响**。真正的雷只有两枚：**(1) 动 id（无 FK 护栏，B2）；(2) 拿"归档"当保留期豁免（RetentionJob 不看它，B3 甲）**。⚠ 另注意 `RetentionJob` 今天**生产零调用者**（A4）⇒ 30 天窗口**当前不会自动兑现**，落地腿若依赖"归档行会到期"，那个前提是空的。

---

## C. 与既有票的分工（编排者给的候选逐格对，指认到格；无人认领的明写）

先报一枚**候选名册本身的缺项**：编排者给的候选（145/167/182/197/213/214/215/216）**漏了票 36 与票 194/196**，而票 36 恰恰是票池里**唯一一枚已经白纸黑字写着"历史任务列表 → 详情 → 逐条删/清空/导出，从 SQLite 经 bridge 取"**的票（`.scratch/wisp/issues/36-result-history-panel.md:28`；`Status: ready-for-agent`、`:4 Claimed by: —`，即**至今未派、未领**）。⇒ 本节按真名册派格，不按候选清单派格。

| 本票的格 | 归口 | 依据（票面原文行） |
|---|---|---|
| **列（按时间）** 的**面板载体**（快照里带不带会话名册） | **票 145**（`Snapshot` 扩字段；其 AC#2 目前**未勾**，编排者自陈"AC#2 的量出来的答案是落地集＝空"，`145-…:11` 一带） | 票 145 就是"Go 侧载体扩出来"那一票；本票是其**上游数据源**（票 166 面自陈"本票是它的上游"） |
| **列** 的**存储/DAO 那一跳**（会话名册的真 SQL） | **票 36**（`36-result-history-panel.md:28` 要 "task list (30d) → task detail … from SQLite via bridge"） | 票 36 的 Go 侧依赖今天没人做；本票普查给出的答案是"**这一跳今天不存在**"（A1） |
| **详情/transcript（逐条消息可见）** | **票 36** B 槽（`:7-8` "B: history/transcript + purge/export + virtualization"） | 同上；⚠ 这一格是 A6 那枚"历史本体不落盘"的直接下游 |
| **逐条删 / 一键清空 / 导出 JSON** | **票 36**（`:28-29`）＋**票 40**（`40-security-privacy-cost-pages.md`，隐私/成本页）；规格出处 `SPEC-02:167-168`、`PLAN.md:2715-2716` | 现量：`internal/memory/privacy.go` 四枚 API **生产调用者＝0**（A4）⇒ 归口明确、缺的是接线不是设计 |
| **搜（标题与正文）** | **无人认领**（票池 0 枚要求会话搜索；`search.content`/`search.files`/`memory.recall` 只在 D34 声明、`A308`（`pending-and-issues.md:7658`）列为"声明了不可及"） | ⇒ 本票普查判：**要么票 36 扩一格，要么新开票**，编排者不许把它塞进 145/167/182（那三枚射程里没有"检索"） |
| **改名** | **无人认领** | 全仓无 title 列（A3）。尺：`grep -rniE "改名|重命名" .scratch/wisp/issues/*.md` ⇒ 命中**逐枚都不是给会话改名**，三族分别是「票文件名加 `-done` 的**改名权**归验收方」（`113-…:19`、`119-…:5`、`125-…:411`）、「commit／常量**改名**」（`11-…:56`、`110-…:49`）、「Go module 路径**改名成本**」（`docs/PLAN.md:1567,1713,2017` 的 P10 那一族）⇒ **会话改名 0 枚票** ⇒ 必须新开票或并入票 36 的新 AC |
| **归档（含撤销）** | **无人认领** | 尺：`grep -rn "归档" .scratch/wisp/issues/*.md` 命中很多，但**语义逐枚是"证据件／派单／报告归档"**（`08-…:73` 归档的 slo 报告、`125-…:405` 派单归档制度、`131-…:440` 纯净树 `git archive`、`146-…:83` 锚点归档副本）；与"会话／对话"同现者仅 **2 行**，一行是 `12-…:163` 的"本会话是唯一 agent"与"票 08 归档跑"同句错拼，另一行是**票 166 面自己** ⇒ **会话归档 0 枚票**。⚠ `internal/models/archive.go` 是 tar 解包、`cmd/wisp/doctor.go:181 ArchivePath` 是便携模式路径字段，两枚同名不同物，别当已有实现。⚠ 落地时会**借用票 196 的地盘**（状态词表／要不要 CHECK），那枚**未派、待 owner 拍**（`196-…:3`）⇒ 两票互相具名，不许互相"归口到对方已勾的格"（另见 H#4） |
| **分叉** | **无人认领**，且**今天不可派**（要等 transcript 先落盘，A5） | 票 213/214/215 的加号菜单里没有这一类（213 管命令、214 管"往这次对话里加东西"、215 管技能/插件/MCP；三票面逐枚读过，`214-…:3` 那句"往这次对话里加东西"指的是**附件**，不是会话派生） |
| **入向方法名册（若会话动作要走面板通道）** | **票 194**（`194-…:1-3`：两份名册互不相认，owner 已裁"按当前规格补齐代码侧、别搞债务"）＋**票 253**（`253-…` 三枚尺洞，含 `config.*` 落在两枚前缀尺之外） | 具名后果见 D2 与 E |
| **占用/排队/停止/草稿那四枚小出口** | **票 167**（面已列） | 会话名册与那四枚**不重叠**：167 管"这一次"的输入输出，166 管"上一次"的检索。⚠ 唯一相邻格：**票 167 AC#5 草稿**若把草稿落盘，会想用同一张新表——那一格属 167，不属本票，别在这儿顺手做 |
| **任务监控栏那几堆** | **票 182**（`182-…` 只读普查，零产码；它自陈"整块不在票 145 那张十四行表里"） | 本票的"会话列表"与 182 的"任务名册"**是两个对象**：182 数的是**在跑/跑过的任务**（今天有 `tasks` 那一枚可选载体，票 197 落的），本票数的是**跨重启的对话**。⚠ 别并：并了会让人以为 `tasks` 载体＝会话名册已通 |
| **子代理三层（`tasks` 载体本身）** | **票 197**（已落 `internal/panel/subagent_roster_197.go`） | 会话若也走"列表＋点击进页"那一形，**形状可参考 197，数据源不是它**；197 的行是 `TaskRow`，键在 task id 上（`composer.go:91 Tasks`） |
| **界面那一腿（列长什么样、在哪儿点）** | **不归本编队**（票 166 面 `skipped=frontend(owner-delegated)`；`frontend/**`／`design/**` 两层禁令） | 本件全程未读那两层；凡"必须看界面才能定"的格，本件登记为"需前端会话自取"，见 G#6 |

---

## D. 改动面最小形状（只给面，不选形、不写码）

⚠ 前置：这一节**不是方案**，是"要动哪几枚文件、每枚几行、要不要新开导出 API"的**面积测量**。真正选形之前有两格必须先由人拍（H#1 schema 契约 vs AC#5；H#2 会话主键）。

### D-最小可闭集（"列＋搜＋改名"三件，**归档与分叉今天落不了**，理由见 D3）

| # | 文件 | 动哪一处 | 估算行数（只算面积） | 是否新开导出 API |
|---|---|---|---|---|
| 1 | `internal/memory/schema.go` | `ddlV1/ddlV2` 之后加 `ddlV3`＋`migrationChain`（`:168-171`）追加 `{from:2,to:3}`＋`SchemaVersionTarget`（`:129`）改 3 | 约 +25 至 +45 | 是（会话 DAO 的表） |
| 2 | `internal/memory/schema_test.go` | `contractSpec`（`:44` 起）、`contractTableOrder`（`:168-171`）、`:222` 版本断言（"2"→"3"）、`:229/:233` 语句枚数常量、新增一列/表规格 | 约 +30 至 +70，**且这是"改断言"不是"加断言"** ⇒ 必须非实现者复判 | 否 |
| 3 | `internal/memory/dao_session.go`（**新文件**） | 会话行模型 + `List*/Search*/Rename*/Archive*` DAO 原语；可复用 `models.go:215 escapeLike`（LIKE 转义已在） | 约 +120 至 +200 | **是**（这是票 145 等的上游出口） |
| 4 | `internal/memory/dao_session_test.go`（**新文件**） | 逐枚 AC 的断言 + 迁移崩溃原子性（G#4） | 约 +150 至 +250 | 否 |
| 5 | `internal/panel/session_roster.go`（**新文件**，或并入 `composer.go`） | 快照载体：新 section + `PumpSources`（`internal/panel/pump.go:196-210` 一带）新增 reader | 约 +40 至 +90 | 是（面板读面） |
| 6 | `internal/panel/composer.go` | `Snapshot`（`:57-91`）新增**一枚 `omitempty` 指针字段**（照 `:74 Instructions`、`:91 Tasks` 那两枚的先例形） | +6 至 +12 | — |
| 7 | `internal/panel/pump.go` | reader→键的装配（`:175-186,206` 一带注释已定规矩："有 reader 必发键、无 reader 不发"） | +10 至 +20 | — |
| 8 | `cmd/wisp/run.go` | 泵装配处把会话 reader 递进去（`rt.store` 已在 `:1128` 可得，**不必新开 store 装配**） | +6 至 +14 | — |
| 9 | （可选）**`internal/panel/bridge.go`** | 若改名/归档要走**面板入向**：新 const（`:42-45,66-67` 那一族）+ `knownComposerMethod` 的 case 列表（`:146-148`） | +6 至 +10 | **是，且这是 C17 契约面**（见 D2） |
| 10 | （可选）**`internal/panel/composer_dispatch.go`** | 新 handler interface（`:69-112` 一族）+ `dispatch` 的 case（`:175-214`）+ `ComposerDispatch` 字段（`:121-139`） | +25 至 +45 | — |
| 11 | （可选）**`cmd/wisp/panel_inbound.go:271-281`** | 把 `nil` 换成真 handler（今天 `Workspace/Attachment/Message` 三枚是**写死的 nil**，`:275-277` 逐字"refused by name when they arrive"） | +8 至 +20 | — |
| 12 | （可选）**`cmd/wisp/main.go` + `run.go`** | 若走 **CLI 腿**：新增 `case` 与 usage 行（`:88-132`、`:21-56`）。⚠ 这一动**必带** `leg_dispatch_gate_133_test.go` 的名册后果（E#7） | +20 至 +60 | — |
| 13 | **`scripts/portable-tests.sh`** | **若新落点在 `internal/session/`**：该包**不在 core scope 也不在 core_pin**（scope 列表 `:207-217`、pin `:140-166`，两处都逐字没有 `internal/session`；注释 `:96-99` 还写着它"doc.go-only、USED to be here and are gone"——**那句对今天的它已过期**：`ls internal/session` ⇒ `doc.go`、`grants.go`、`session.go` ＋ **两枚**测试文件 `grants_test.go`、`ticket224_pattern_dialect_test.go`）⇒ 落 `internal/session` 必须同批改这两枚列表，否则**新增用例 CI 一根不跑** | +2 | — |

### D2 ⚠ 新增 `panel.*` 入向方法名＝C17 契约面（具名点出）

- 若落地腿选择**新加 `panel.session.*`** 这一形：它会**直接红**`internal/panel/git_test.go:385`（`whitelistMethodsFromSource` 扫 `internal/panel/bridge.go` 与 `wantMethods`（`:381-383` 四枚）比相等，失败消息在 `:386-388` 逐字"want the four methods this ticket may not extend"），**并且**红同一文件 `:517` 的正控（`if real := …; len(real) != 4`）。⇒ 这是**票 181 那一族"git 维度不许多出入向门"的钉**，不是随手的计数。
- 若落地腿选择**新加 `config.*`/无前缀那一形**（如 `session.list`）：那两枚**前缀尺都看不见它**——本程自跑复核（不是引用别人一句话）：
  - `panelMethodRe`（`internal/panel/git_test.go:394`）逐字 = `"(panel\.[a-z0-9_.-]+)"`，**前缀硬编码**；
  - 它只扫**一枚文件**：`filepath.Join(root, "internal", "panel", "bridge.go")`（`git_test.go:385` 实参），不看别处；
  - ⇒ **同文件里今天真实存在的 `"config.get"`／`"config.set"`（`bridge.go:66-67`）它一枚也不计**：所以 `whitelistMethodsFromSource` 的返回集恰是四枚 `panel.*`。⇒ **另一枚只读腿报的"那把锚看不见已存在的 `config.*` 方法名"这一句，本程复核成立**；
  - 但它报的**行号 `:407` 需要更正**：`:407` 是 `func whitelistMethodsFromSource` 的**声明行**；真正"恰四枚"的**断言在 `:385-389`**（另一枚正控在 `:517`），而"看不见 `config.*`"的**根因在 `:394` 的 regexp**。票面若要引用锚，引 `:385`（断言）与 `:394`（射程），不要引 `:407`。
  - ⚠ 这条盲区的**另一半**：`composer_test.go:393 routeLiteralRe` 同样是 `panel.` 前缀尺（`inbound_roster_253_test.go:437` 的失败消息逐字点名这两枚："The two prefix rulers cannot see a name that carries no \"panel.\" prefix (panelMethodRe at git_test.go:394, routeLiteralRe at composer_test.go:394)"）。⇒ 唯一认全名的尺是 `inbound_roster_253_test.go`：`wantFullInboundRoster253`（`:69-76`，**恰六枚名**：`config.get`/`config.set`/`panel.attachment.add`/`panel.message.send`/`panel.mode.request`/`panel.workspace.request`）、`wantFullInboundRosterSize253 = 6`（`:59`）、`wantNonPanelPrefixedInbound253 = 2`（`:64`）。**任何新增入向名，三处常量必须同批改**，否则 `TestFullInboundMethodRosterIsClosed`（`:419`）按五个分母逐一红。
- ⇒ 结论：**新增 `panel.*` 入向名＝改 C17 白名单＝人工批准**（`AGENTS.md` §1.1：`C1–C32` 变更须人工批准；`bridge.go:56` 的注释也把这枚锚具名回指 `git_test.go`）。本票普查判：**"改名/归档"若非要面板当场触发，Go 侧的读面（快照）＋ CLI 腿可以先闭，避免这一枚契约动作**——但**这是编排者/owner 的选形，不由本程定**。

### D3 为什么"归档带撤销"与"分叉"今天**落不了**（面积之外还缺东西）

1. **撤销要有可逆对象**：今天 `task_log` 没有归档位（A4），而"撤销"要么记旧值要么记事件；`internal/memory` 全族**没有任何 append-only 事件表**（九张表名册见 `schema.go:22-120,134-146`）。⇒ 撤销那一格的最小面积＝**再一张表或一列 + 一条幂等键**，落在与 DDL 契约钉同一枚面上。
2. **分叉要有可复制的历史**：A5——历史本体不落盘（A6 第一行）。⇒ fork 的前置是 transcript 持久化（票 36 B 槽那一格），**不是本票的 DAO**。
3. **归档别与保留期混**：`retention.go:32 TaskLogTTL` 的删行谓词不看 state（B3 甲），而 `RetentionJob` 今天**零生产调用者**（A4）⇒ "归档 30 天后真清"这一形今天**两头都不成立**。

---

## E. 邻居尺名册（词面型／计数型负向尺；文件:行 ＋ 它数什么 ＋ 落地怎么隔离才不误响）

| # | 尺（文件:行） | 它数什么 | 本跳落地会不会误响 / 怎么隔离 |
|---|---|---|---|
| E1 | `internal/panel/git_test.go:385`（消息 `:386`），射程 `:394`，正控 `:517` | `bridge.go` 里 `"panel.*"` 字面量的**去重集恰＝四枚**（且只扫 `bridge.go` 一枚文件） | **只认前缀不认 `config.*`**（D2 已复核）。新加 `panel.*` ⇒ 必红（＝契约动作，须批准）；新加 `config.*`／无前缀名 ⇒ **这把尺不响，改由 E2 响**。隔离办法：**读面走快照、不新开入向名** |
| E2 | `internal/panel/inbound_roster_253_test.go:419`，名册 `:69-76`，计数 `:59`=6、`:64`=2，五个分母 `:434-479+` | **全名入向名册**（前缀盲区的补尺）：声明集／计数／`knownComposerMethod` 运行守卫／守卫 AST case 标签／路由点，双向比 | 任一新增入向名**必红**（除非同步改三处常量）。⚠ 这把尺自己说"若计数动了，要么别的尺更盲、要么名册长了，diff 必须说清是哪一种"（`:456-457`）⇒ **不许只改名册不改计数** |
| E3 | `internal/panel/composer_test.go:393`（`routeLiteralRe`），`:417`（known 集），`:501`、`:540-571`（把 `frontend/src/lib/panel.ts` 拷进临时树再扫） | 面板 bundle 里每一枚 `"panel.*"` 字面量**必须是 Go 答得出的**（反向：不许有前端调而后端无的门） | 本程**未读 `frontend/**`**（两层禁令），只读这把尺自己。⚠ 若新增 `panel.*` 而前端未跟上，红的是这把尺（不是 Go），后果归**票 145 AC#4③ 那一格**（它自己就问"这把尺今天是否存在"） |
| E4 | 快照**顶层键**钉，共 5 枚：`internal/panel/pump_test.go:123`、`pump_test.go:291`、`instructions_200_test.go:251`（`len(generic) != 4`）、`subagent_roster_197_test.go:214`、`subagent_stream_197_test.go:130`（`len(generic) != 4`） | **无可选 reader 的泵，序列化后顶层 JSON 键恰＝`composer,generatedAt,pending,results`**（枚数与逐字串两种都有） | ⚠ **票面特别要核的那一格，答案是"会红，但只在一种形下红"**：新增**不带 `omitempty` 的顶层字段** ⇒ 5 枚**同时红**；新增**指针 + `omitempty` 且该 host 无 reader 时不发键**（`composer.go:74 Instructions`、`:91 Tasks` 的既有形，规矩写在 `pump.go:175-186` 与 `subagent_roster_197_test.go:189-216`）⇒ **5 枚全绿**，因为这几把尺钉的是"reader-less 的泵"。⇒ 隔离办法＝**会话载体必须照 `Tasks` 那枚的形：指针 + omitempty + 无 reader 不发键**；这也正是票 197/200 两票各自的 AC 为什么把"没这个能力"和"有能力但这轮空"做成两枚可读状态（`subagent_roster_197_test.go:218-237` 逐字） |
| E5 | `internal/memory/schema_test.go:173`，逐点：`:195`（表枚数＝9）、`:259`（逐表**列枚数**）、`:222`（`schema_version` 值＝`"2"`）、`:229`（ddlV1 语句枚数＝15）、`:233`（ddlV2＝1） | D35 数据模型的**契约计数钉**（表/列/索引/语句/版本五类） | ⚠ **任何加列/加表/加索引/加迁移步骤都会红**，且这些是"改断言"而非"加断言"⇒ 必须走**非实现者复判**，并与 AC#5（`docs/specs/**` 零字节）＋`schema.go:10-13`（DDL 与 SPEC-02 逐字节同）三方对质。**这一格本程判为待人拍板（H#1）**，不许落地腿自行放宽 |
| E6 | `internal/memory/migrate_v1seed_test.go`、`schema_test.go:324/:382/:417/:447`、`:458 TestMigrationFailedStepIsAtomic`、`:510/:564` | 迁移链：新库/老库水位、备份文件名规则、失败步骤原子性、更新版本库拒迁 | 加 `{from:2,to:3}` 后：`:324` 的"fresh database creates current version"会跟着动（它断的是**当前目标版本**）；`:458` 那一枚是 AC#3 崩溃原子性的**现成可借形状**（它已经在拿"步骤中途失败"证 watermark 与半列不落盘）。⚠ 数字前缀当幂等键那一形（票面 AC#3 后半句）：新库**不缺任何列**、老库**缺 `v3` 全部新增**，两者幂等键都是 `schema_meta.schema_version` 那一行（`schema.go:151-159` 注释 + `open.go:403` 种子写入）——**这一格只有跑了才知道现量，见 G#4** |
| E7 | `cmd/wisp/leg_dispatch_gate_133_test.go:151 usageCommand133`（`^  wisp ([a-z][a-z0-9-]*)`）、`:154`、`:178`、`:142 minLegs133`、`censusVsUsage133`（`:218` 调用） | **CLI 腿名册**：`func main` 的 switch case 枚数 ≥ 下限，且**与 usage 块逐名对账**；`func main` 里任何跑在已分类分支之外的调用＝红 | 若走 **CLI 腿**（D#12）：**case 名必须同批出现在 usage 块里**（`main.go:21-56`），并按票 133/131 的规矩为那枚腿配自己的钉或写明 `WISP-LEG-COVERAGE-RULING`。隔离办法：**列/搜先用面板读面（快照）＋ DAO，不开新 CLI 命令**；真要开，别分两次 commit |
| E8 | `scripts/portable-tests.sh`：GUARD A/B/C 的定义（`:58-66`）、core scope 列表（`:207-217`）、`core_pin`（`:140-166`）、census 分支（`:256-299`） | **包级覆盖率名册**：声明的包必须编译进至少一枚测试文件、必须自己印 top-level 结果行、pin 与实际解析集双向比 | 若新代码/测试落 `internal/session/`：**它今天不在 scope 也不在 pin**（D#13，注释 `:96-99` 那句"doc.go-only"已过期）⇒ 不改两列表＝**新增用例 CI 完全不跑**（这不是"尺响"，是**尺看不见**，比响更危险）。隔离办法＝**落 `internal/memory`（在 `./internal/memory/...` glob 内，`:211`）**，或同批把 `internal/session/` 加回两枚列表 |
| E9 | `tools/d22scan`（CI 扫描面见 `.github/workflows/ci.yml:58` 注释：`design/ + internal/ + cmd/`；票面另记 `AGENTS.md §1.2` 的实测带） | 禁形：裸 `go func(`、`PathResolver` 之外的 `filepath.Clean/Abs`、明文密钥、墙钟超时、面板侧来源的 L2"允许"、把内部 artifacts 写成受门控 Tool、emoji（字符串不豁免、注释豁免） | 新增会话 DAO/面板字段都在这三枚目录里 ⇒ 会被扫。**已知缺口别踩**：`AGENTS.md §1.2` 自陈仪器带与规格带双向不等（`→` 界面文案抓不到，`✓` U+2713 与 `≤` U+2264 **会响**）⇒ 会话名册的**文案里不要出现 `✓`/`≤` 这类字形**。⚠ 本程**未跑 d22scan**（禁跑），现量读数见 G#5 |
| E10 | `internal/agent/loop.go:356`（`correlation == task id`）＋`cmd/wisp/approval_reply_201_test.go:681`（`strings.TrimSuffix(corr, "-corr")` 之类形状）＋`run_test.go:295,310,392` | **词面型**：测试里把 correlation id 当 task id 的**拼接形状**用 | 会话若引入自己的 id 空间，这些用例的 id 推导链会变；落地腿要逐枚核。**本程只登记形状，不判会不会红**（要跑才知道，G#3） |
| E11 | `internal/observe/thresholds.go`（票面列的零字面区）＋`internal/agent/spill.go` 一族 ACL/路径不变量测试（`spill_acl_windows_test.go` 等） | 阈值与工件路径不变量 | 会话名册**不该新增受门控 Tool**（`AGENTS.md §1.2` 逐条禁"把宿主内部 artifacts 写入实现成受门控的 Tool"）⇒ 会话列出一旦做成**模型侧工具**（例如 `session.list`）就是**那条禁令的近邻**：本票判它属**宿主侧读面**，理由与票 182 AC#3 那句"模型侧工具 vs 宿主侧读面不是一回事"同形。⚠ 这一格若被落地腿读成"可以做成工具"，请直接退回，别顺走 |

---

## F. 推翻清单（票面与派单里每一句都当待验断言复跑；不成立的具名推翻并给真身）

| # | 原断言（出处） | 复跑读数 | 判定 |
|---|---|---|---|
| F1 | "会话今天没有管理这一层：列不出来、搜不到、改不了名、断不了叉"（票 166 标题） | 见 A1–A5：五枚里**只有一枚有真 SQL 跳**（`ListTaskLogs`，`dao_tasklog.go:93`），且**包外生产调用者＝0**；其余四枚**存储层那一跳不存在** | **成立，且比原断言更硬**（原话只说"没出口"，真身是"盘上没有会话这个实体"，A0） |
| F2 | "`internal/store` 或 DAO 里的真 SQL"（派单 A 问的措辞） | **`internal/store` 这枚包不存在**：`ls internal` 得 23 枚目录，无 `store`；真身是 **`internal/memory`**（`schema.go`/`dao_*.go`/`privacy.go`/`retention.go`），且 `internal/memory` 在 `docs/specs/SPEC-02` 里就叫 memory 模块 | **推翻（名册级）**：本票及后续引用请写 `internal/memory`，别让人去找 `internal/store` |
| F3 | "那把锚看不见已存在的 `config.*` 方法名，行号 `:407`"（派单 D 转述另一枚只读腿） | 盲区**成立**：`panelMethodRe`（`git_test.go:394`）前缀硬编码 `panel\.`，且只扫 `bridge.go` 一枚文件（`:385` 实参），而 `bridge.go:66-67` 逐字有 `"config.get"`/`"config.set"` ⇒ 那两枚**不在它的计数集里**。**行号不成立**：`:407` 是**辅助函数声明行**（`func whitelistMethodsFromSource`），断言在 `:385-389`，另有正控 `:517`，根因在 `:394` | **半推翻**：射程判断复量成立；**行号按 `:385`（断言）／`:394`（射程）／`:517`（正控）三处引，`:407` 那个号作废** |
| F4 | "生产可达的工具只有 `fs.*` 六枚"（票 166 来源行引 W1 底盘普查、台账 `A308`） | 今天 `cmd/wisp/run.go` 的注册循环共**三处**：`:515`（fs 族）、`:544`（task 族）、`:780`（subagent 族）。fs 族＝`fs.read`/`fs.list`（`internal/tools/fs.go:326-329`）＋ `fs.write`/`fs.edit`/`fs.trash`/`fs.move`＋受 `[fs] delete_enabled` 门的 `fs.delete`（`fs_write.go:772-783`）；task 族＝`task.output`/`task.cancel`（`task.go:600-604`）；subagent 族＝`task.spawn`（`subagent_197.go:218`） | **推翻（过期）**：现量 **恒在 8 枚（fs 六 + task 二）＋ 门控 1 枚（`fs.delete`）＋ 子代理 1 枚（`task.spawn`）**。`A308`（`pending-and-issues.md:7656-7658`）那句"实现 6 枚"锚在 `192ad56`，票 162/164/221/197 落地后已不覆盖 |
| F5 | "候选归口票＝145／167／182／197／213／214／215／216"（派单 C） | 票池实际相关名册里**至少漏三枚**：**票 36**（`36-result-history-panel.md:28` 逐字要 "task list (30d) → task detail … from SQLite via bridge"，`Status: ready-for-agent`、`Claimed by: —`＝未派未领）；**票 194**（入向方法名册对齐，owner 已裁）；**票 196**（状态词表/CHECK，正是"归档"那一格的地盘，且自陈**未派、待 owner 拍**） | **推翻（名册级）**：候选清单必须补这 3 枚，否则"归档/逐条删/导出"会被错并进 145/167/182 |
| F6 | "`D35` 那张表是定过的；加列要走迁移"（票 166 硬约束 1） | 前半句成立（`schema.go:22-146` 九张表、迁移链 `:168-171`、事务内 + 前置备份）。**但"加列走迁移就够了"这一句不成立**：`schema.go:10-13` 逐字要求 DDL 与 `docs/specs/SPEC-02-data-storage.md §3` 的 ```sql 块**逐字节相同**，"No added or removed columns, tables, indexes or constraints"；`schema_test.go:173` 那族钉把表枚数/列枚数/语句枚数/版本值全钉死；而票 166 AC#5 要求 `docs/specs/**` 零字节 | **推翻（前提级）**：**"加列"与"`docs/specs/**` 零字节"在现仓规下互斥**。这一格必须 owner 拍（H#1），**不是落地腿可以自行放宽任何一侧的**（`AGENTS.md` §1.1：SLO/golden/阈值一字节不许动；spec 文字改动＝契约变更） |
| F7 | "minimax 那条迁移不是崩溃原子的教训：两枚改名＋预写日志那一形"（票 166 硬约束 1） | 本仓**读不到对应物**，三把尺：`grep -rln "崩溃原子\|不是崩溃原子" docs .scratch` ⇒ **2 枚文件，一枚是票 166 面自己、一枚是本件**；`grep -rn "两枚改名" docs .scratch/wisp/issues` ⇒ **0**；`grep -rn "迁移" docs/evidence/s1/*.md` 与 `minimax`/`原子` 同现 ⇒ **0 行**（minimax 那族外部对标只出现在审批面，如 `docs/evidence/s1/219-approval-reply-surface-c1.md:406,441`，与迁移无关）。另：本仓的 WAL 是**连接串 pragma**（`internal/memory/open.go:148-151`：`_pragma=journal_mode(WAL)`＋`busy_timeout`＋`wal_autocheckpoint`），**不是迁移手段**（`grep -rn "journal_mode" internal/memory --include='*.go' \| grep -v _test`） | **登记为外部经验、本仓无对应尺**：不推翻票面那句教训，但**不许把它读成"本仓已有防线"**。崩溃原子性的本仓真凭据是现成可借的两枚：`internal/memory/schema_test.go:458 TestMigrationFailedStepIsAtomic` 与 `internal/memory/concurrent_test.go:179`（"Crash recovery with a real subprocess"） |
| F8 | "`session id／task id／correlation id 的分工`"（派单 B 假设三者已分工） | 现量：**correlation id 与 task id 今天不是分工，是别名**——`internal/agent/loop.go:356 newTaskJournal(l.opt.Journal, taskID, taskID)`，注释逐字"C18: correlation == task id"；`session id` 独立（`cmd/wisp/run.go:473` 铸，随进程死） | **推翻（分工级）**：三枚里只有两枚是真分工（task vs session）。任何"按 correlation 归组会话"的想法今天**等价于按 task 归组**，别把它读成"已经有会话维度的关联键" |
| F9 | "票 145 那格'十四态里十一态无输入可画'会一直闭不掉——缺的就是这层数据"（票 166 §为什么值得做） | 票 145 面自陈：`AC#2/AC#6` **未勾**，且 AC#2 的**量出来的答案是"落地集＝空"**（`145-…` Status 追加段，逐字）。十四态对照表的行来源是 `docs/PLAN.md` 的 §17.5 那张表，真身坐标＝`:3481-3494`（见 F10，不是票面引的 `:3473-3488`） | **部分成立**：145 的"未闭"确实与"没有会话/历史这一层数据"同源（快照里确无会话字段，`composer.go:57-91`）；但**"缺的就是这层数据"这句偏大**——票 145 自己列的缺项里还有 token 计数、推理过程、错误分类、重复计数、成本、取消原因、当前屏（`145-…` §现量），**这些与"会话名册"是不同源的两批**。别在派单时把 145 整枚挂在 166 上 |
| F10 | "`docs/PLAN.md:3473-3488` 那张十四行表"（票 145/182 引用、派单 C 借来的坐标） | 现量该区的真身：`:3474` 是 `## 17.5 Agent 工作状态的无 emoji 视觉设计` 标题行，`:3479` 表头「状态｜视觉｜动效｜明确不用」，**十四枚状态行落在 `:3481-3494`**（尺：`awk 'NR>=3474 && NR<=3500 && /^\| \*\*/ {c++} END {print c+0}' docs/PLAN.md` ⇒ **14**；逐行首列＝思考中/SSE 流式/推理过程/工具调用/工具调用（展开）/审批等待/L2 确认卡/L1 阻止窗口/L2 原生降级卡/错误/Stuck/成本/已取消·被打断/注入检出）⇒ **票面写的 `:3473-3488` 少算了 `:3489-3494` 六行，多算了 `:3473` 一行** | **推翻（坐标过期，且是双向偏）**：十四行表今天真身＝`docs/PLAN.md:3481-3494`。⚠ 本票**未改 `docs/PLAN.md` 一字**，只给坐标读数。这一格对票 145 的账有实际后果：按 `:3473-3488` 数只能数到 8 行，会**把后六行（含 `成本`、`注入检出`）读成不存在** |

---

## G. 量不到的格子（禁跑 Go 的硬闸门下，如实登记 + 复量法；不含任何猜测）

| # | 量不到的那一格 | 为什么量不到 | 复量法（等三枚写腿交完、无并发时跑） |
|---|---|---|---|
| G1 | **基线门禁四数**（`go test ./internal/panel/ ./cmd/wisp/`、`gofmt -l`、`gofmt`/`gofumpt -l`、`sh scripts/d22scan.sh` rc） | 有 3 枚写腿在飞（`255-r3`、`252-r2`、`171-r3`），并发跑测试互洗读数；本程一律不跑 | 锚点记下后：`bash scripts/wisp-cli-tests.sh`（它会先 staged sherpa DLL，票 145 AC#5 明写基线要用这支或带 DLL 取），另跑 `bash tools/d22scan/runtests.sh ./internal/panel/`。**四数之外要点名册差集**（票 145 AC#5 原话） |
| G2 | `E5` 那族 schema 钉**当前是绿的还是红的**（我只读到"它会红"的形状，没读到它的 verdict） | 需要 `go test ./internal/memory/` | 写腿清空后单跑 `./internal/memory/`；`TestSchemaContractIntrospection`/`TestMigration*` 四枚逐枚 `-run` |
| G3 | `E10` 那批**把 correlation id 当 task id 拼形状**的用例，新增 id 空间后会不会红 | 需要跑 `cmd/wisp` 与 `internal/agent` 全族 | `bash scripts/wisp-cli-tests.sh` ＋ `./internal/agent/`；红了按用例名逐枚归因，不许整体放宽 |
| G4 | **AC#3 的崩溃原子性现量**（中途被杀之后重开是否"要么全成要么全不成"、新库/老库各缺哪列的实测） | 需要真起进程真 kill（`go test` 之外的 subprocess；本仓已有形状：`internal/memory/concurrent_test.go:179` 那节"Crash recovery with a real subprocess"） | 照 `TestMigrationFailedStepIsAtomic`（`schema_test.go:458`）+ `concurrent_test.go:179` 的两枚现成形状各跑一发，并**分别**用"新库（0→1→2→3）"与"老库（已 2、跳 3）"两种水位，验 `schema_meta.schema_version` 与列存在性同时要么都动要么都不动 |
| G5 | `d22scan` 的**被扫文件计数现量**（`AGENTS.md` 引的 `internal/=223 / cmd/=29 / ban#8 internal/=473 / cmd/=65` 全是旧锚读数） | 需要跑扫描器 | `./tools/d22scan/d22scan.exe -root .`（独立 module，别在根目录 `go vet ./tools/d22scan/`，票 145 AC#5 原话） |
| G6 | `internal/session/` **今天到底被哪一枚 CI scope 认领**（census 的读数） | `bash scripts/portable-tests.sh --scope=census` 要跑 Go；本程**只读到"它在 scope 列表与 pin 里都逐字缺席"**这一半 | 跑 census，看它是否打印 `CLAIMED BY (none)`。⚠ 台账 `pending-and-issues.md:4557` 那句"`GAP-22`（`internal/session/` HEAD 只有 `doc.go`）"**对今天的它已过期**（它有 `session.go:1-101`、`grants.go`、`grants_test.go`、`ticket224_pattern_dialect_test.go`） |
| G7 | 票面 AC#2 要求的"每件都要答'这一发在未修码上响不响'"的**逐发实际响应** | 响应＝行为，要么跑要么变异后跑；本程一律不跑 | 落地腿起手各跑一发未修码基线：(a) 面板：起 `wisp panel-inbound -data <tmp>`，喂一枚 `{"method":"session.list",…}`，应当被 `knownComposerMethod`（`bridge.go:146-148`）按名拒；(b) CLI：`wisp run --resume` 之类形状应落在 `parseRunArgs` 的未知旗标上；(c) 存储：对真库发一条 `UPDATE task_log SET …` 形状看有没有可改的标题列（答：没有，A3）。**这三发的预期来自本件读数，不等于本件已验** |
| G8 | `frontend/**`／`design/**` 那一侧**会话界面今天有没有**（票面"连列出来都没有出口"若界面已有桩，会影响归口） | 两层禁令：本程**不读不引** | 归**owner 委托的前端会话**自取；本件不代答，也不想象界面形状 |

---

## H. 未定义即停上报（`AGENTS.md` §2 那一族的规矩：碰到就停手，不自行假设）

1. **schema 契约 vs AC#5 互斥**（F6）：票 166 AC#3 要加列/加表的迁移，AC#5 要 `docs/specs/**` 零字节，而 `internal/memory/schema.go:10-13` 要求 DDL 与 `SPEC-02 §3` **逐字节相同**且 `schema_test.go:173` 一族把列数/表数/语句数/版本全钉死。**这三条今天不能同时成立**。⇒ 需 owner 拍：是"扩 `SPEC-02 §3`（＝改 spec，人工批准）"，还是"新表不写进 `ddlV1` 而走独立 v3 段并同步 `contractSpec`"，还是"会话不进 SQLite（那 D35/§14.4 的'必须存'又要重读）"。**本程一律不选。**
2. **会话主键未定义**（A0/B1）：跨重启的"会话"在本仓没有键。新造一枚＝新契约面（台账 `9220` 那行已经为票 224 记过一次同形代价："要新增一枚只读 seam ＝ 新契约面"）。⇒ 需人拍：会话＝`task_log.id` 的分组（则今天只有"一次任务＝一个会话"这一枚退化形），还是新 `conversation.id`。
3. **`[privacy] keep_transcript` 硬编码 false 的射程**（A6 末行）：`internal/config/schema.go:504-515` 与 `internal/config/validate.go:81-91` 把它钉成"写 true 就加载失败"，`docs/PLAN.md:1874` 又逐字要求"SQLite 历史必须存"。二者今天靠"存脱敏后的 `query_text`/`summary_text`、不存逐字全文"来分界——**但会话消息级持久化落在哪一侧，规格没说**。⇒ 需人拍（与 `AGENTS.md` §2 列的"便携模式 × DPAPI"同级地属待定案）。
4. **归档 vs 票 196**：`archived` 若进 `state` 列，等于给一枚**无 CHECK、两套词表并存**的列加第五个词（票 196 现量，且那枚票**未派、待 owner 拍甲/乙/丙**）。⇒ 不许在 196 未拍之前先写 166 的归档形。
5. **分叉的前置不在本票**（D3#2）：历史本体不落盘 ⇒ fork 必须先有 transcript 持久化（票 36 B 槽那一格）。⇒ 票 166 AC#4 在现量上**无可验对象**，建议编排者把它读成"条件 AC"，别让人为了勾它去造桩（`issues/README` 硬约束：不许拿 mock 代替真的来假报完成）。

---

## I. 自查

- 空壳自查：本件每一格都有读数或"量不到＋复量法"，无空格、无待补字样（编排者若用占位符字面尺扫本文件，会唯一命中这一句自己，故此处不写出那几个字面 token）。
- 尺的收尾一律 `| wc -l`；零命中处用文字明写"0 命中"，未用 `grep -c` 接 `&&`。
- 未跑任何 Go 命令；未写 `frontend/**`/`design/**`；未产码一字；未碰任何 `^- [ ]` 勾选框；未改 `docs/PLAN.md`/`docs/specs/**`。
- 全文无 API 密钥值，出现的只有变量名/配置键名（`api_key_ref`、`[privacy] keep_transcript` 等）。
- 本文件落在 `.scratch/wisp/probes/166/a1/`，单文件 commit，显式 pathspec。
