# 票 368 / 腿 `368-a1` 普查（只读，`AC#0`）— `COALESCE` 遇上三列全 NULL：孤儿造得出、**永生**，且**今天就有生产写入者**

- 腿＝`368-a1`（**只读腿**，只做 `AC#0`）· 开工 anchor HEAD＝`127784149ea3b6c45896f88caeb7b5384f681fdd` · 写件时 HEAD＝`bea959a7273b2b01ad63216d5beef3d7729d7ebf`（期间进了别人 3 笔：`ed7ecc420`／`fd3c6ae09`／`bea959a72`）。★**被本件引用的 10 枚文件在这两个 HEAD 之间逐枚 `git rev-parse` 全 `SAME`** ⇒ 下面每条行号两处同时成立；漂移另报（§1.4）。
- 写面＝本件 ＋ 仓外 `D:\tmp\wisp\20261024\368a1\**`。仓内 `.go` **0 字节**：3 枚探针（memory／agent／tools）全在仓外，经 `-overlay` 注入 `git archive HEAD` **导出树**跑；导出树 8 枚 `git hash-object` ＝ HEAD blob（逐枚 `SAME`）。
- 本腿运行期间**别的腿在工作树里改了 `cmd/wisp/*.go` 5 枚**（开工 09:33 我现跑 `git status --porcelain -- internal/ cmd/ scripts/`＝**空**，收尾时出现 `M`）⇒ 本件**所有代码读数一律取 HEAD blob**（`git show`／`git grep HEAD`），⛔ 工作树读数。

**尺的射程（落笔前五问）**：①`PRAGMA foreign_key_check` 的行数＝当前 DDL 里**声明了的** FK 的悬空行，**不要求** `foreign_keys` 开着（`366-a1` §1 已证）⇒ 全部读数取在**业务形状连接**（`dsn()`，FK 关）。②`RunRetentionOnce` 量的是**生产那一发**（同一个 `runRetention`），`StartRetentionJob` 只是它的定时器外壳 ⇒ "没接线"≠"这一发不删"。③夹具＝生产链 `Open(dir, withSchemaTarget(3))` 建出的**真 v3 库**，迁移＝生产 `Open(dir)`，⛔ 手抄 DDL。④"三列全 NULL" 数的是 `tool_call` 的 `decided_at`/`started_at`/`ended_at` 三列 **SQL `NULL`**，⛔ Go `nil` 指针、⛔"没写 decision/outcome"（§7 恰有一行 decision/outcome 有值而三列全 NULL）。⑤行号用 `git grep -n <pat> HEAD -- internal cmd` 现跑；票面数字一律夹具复现，⛔ 采信。

---

## §1 三条具名断言逐条复核（含行号漂没漂）

**§1.1 `366-a1` 实测 —— ★复现成功、逐数对上**。F4＝**它的形状逐字**（`task-old` 起于 31 天前 ＋ 1 枚 NULL 子行 ＋ 1 枚老时间戳子行 ＋ `task-fresh` 对控）⇒ 真 v3 库 ＋ 生产 `RunRetentionOnce`：

| 读数 | 我的 F4 | `366-a1` 报的 | 判 |
|---|---|---|---|
| `tasklog_deleted`／`toolcall_deleted` | **1**／**1** | 1／1 | ✅ |
| 留下的孤儿 | **1**（`tool_call|rowid=1|parent=task_log|fkid=0`） | 1 | ✅ |
| fkc：retention 前→後，迁移後 | **0→1**，迁移后仍 **1** | 0→1、迁移后 1 | ✅ |
| 今天 `Open()` | **`<nil>`／schema_version=4** | 能开 | ✅ |

⚠**一处口径要改**：`toolcall_deleted=1` **不是** NULL 子行被删，是旁边那枚**老时间戳**子行被删 —— F1 只放 NULL 子行时 `toolcall_deleted=0`、孤儿仍是 **1** ⇒ **1/1 这对数是夹具成分，⛔ 这条路径的不变量**。

**§1.2 `StartRetentionJob` 零生产调用者 —— ★真，且比票面更宽**。`git grep -n "StartRetentionJob\|RunRetentionOnce\|retentionLoop" HEAD -- internal cmd` ⇒ 非测试命中只有 `retention.go:95`（注释）`:98`（定义）`:100`（panic 串）`:104`（自身调用）`:108`（定义）`:233-234`（注释）`:235`（定义）；调用者只有**测试** `retention_test.go:150`／`:198`、`migration_fkcheck_366_test.go:159` ⇒ ★**`RunRetentionOnce` 也零生产调用者**（`:233-234` 注释承诺的 "the doctor command" **盘上不存在**）⇒ **整条 retention 子系统今天一次都不跑**。

**§1.3 第五条路径的代码 —— ★在**（射程见 §6）。`loop.go:391` `if err := j.startTask(ctx, input); err != nil {` → `:392` 只 `Warn`、`:393` 收块、**无 `return`**；随后 `:671` `p.rowID = j.startCall(...)`；`journal.go:68` 定义、`:75-82` 插入语句。

**§1.4 行号漂没漂（原引 vs HEAD 现跑）**

| 原引 | HEAD 现跑 | 漂 |
|---|---|---|
| `journal.go:75-82`（票面） | `:75-82` 逐字同一段 | ✅ 未漂 |
| `retention.go:146`／`:154-155`／`:197`、`models.go:160-183`（`366-a1`） | 同左（`validateToolCall` `:160`–`:183`） | ✅ 未漂 |
| `agent/loop.go:670/:909/:941`（`366-a1` 的 `startCall` 三处） | **`:671/:910/:942`** | ⚠ **+1** |
| `agent/loop.go:390-392`（票面"`startTask` 失败只 `Warn`"） | `:391`＝`if`、`:392`＝`Warn`、`:393` 收块（`:390` 是上一句 `l.append`） | ⚠ **起点 −1**（事实未变） |
| `schema.go:82`（`tool_call.task_id REFERENCES`）／`schema.go:261`（`session_turn.session_id REFERENCES`） | **`:83`**／**`:286`** | ⚠ **+1**／**+25** |
| `open.go:374`（票面硬提醒 1："`version==s.target` 早退"） | **`:413`** | ⚠ **+39**（票 366 `r1`/`r2` 大改） |

---

## §2 `startCall` 插的三列：★真全 `NULL`（逐字）

`internal/agent/journal.go:75-82` —— **结构体字面量里根本没有 `DecidedAt`/`StartedAt`/`EndedAt`**：

```go
	rowID, err := t.j.InsertToolCall(ctx, memory.ToolCall{
		TaskID:        t.taskID,
		Seq:           t.seq,
		Tool:          tool,
		ArgsJSON:      string(args),
		RiskLevel:     riskColumn(risk),
		CorrelationID: t.corrID,
	})
```

它落到 `internal/memory/dao_toolcall.go:24-31`（逐字）：`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, decision, decided_at, started_at, ended_at, outcome, error_class, correlation_id, grant_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`，实参逐字 `tc.TaskID, tc.Seq, tc.Tool, tc.ArgsJSON, tc.RiskLevel, nullIfEmpty(tc.Decision), nullableInt64(tc.DecidedAt), nullableInt64(tc.StartedAt), nullableInt64(tc.EndedAt), nullIfEmpty(tc.Outcome), nullIfEmpty(tc.ErrorClass), tc.CorrelationID, nullableInt64(tc.GrantID)` ⇒ nil ⇒ 绑 `NULL`。
**实测**（F1 pre、真库）：`child id=1 decision=<NULL> decided_at=<NULL> started_at=<NULL> ended_at=<NULL>`。
★**附读数（`AC#1` 会用）**：`tool_call.started_at` **生产上没有任何写入者**——除 `dao_toolcall.go:25` 的 INSERT 列名与两处 SELECT/ORDER BY 外，全仓无一条 `UPDATE … started_at`，两个生产构造点（`journal.go:75`、`bridge.go:1210`）**都不设**它 ⇒ 该列**永远 NULL**，`COALESCE` 的**第二臂在生产里是死臂**。

## §3 `COALESCE` 逐字 SQL ＋ 它最后一句老化条件

`retention.go:154-155`（逐字，tsExpr 是**表达式**）：`res.ToolCallDeleted, err = deleteOlderThan(ctx, tx, "tool_call", "COALESCE(ended_at, started_at, decided_at)", c.TaskLogTTL, now)`
`retention.go:195-198`（逐字，`fmt.Sprintf` 拼成**唯一一句**）：`cutoff := now.Add(-ttl).Unix()`；`` q := fmt.Sprintf(`DELETE FROM %s WHERE %s IS NOT NULL AND %s < ?`, table, tsExpr, tsExpr) ``；`res, err := tx.ExecContext(ctx, q, cutoff)`。

⇒ **最后一句老化条件**（代入后逐字）：

```sql
DELETE FROM tool_call WHERE COALESCE(ended_at, started_at, decided_at) IS NOT NULL
  AND COALESCE(ended_at, started_at, decided_at) < ?
```

- 边界：`:192-194` 注释明写 **`age == ttl` 保留**、`>30 天` 才删；`?` ＝ `c.Now().Add(-TTL).Unix()`。
- 父那一侧：`retention.go:146` 逐字 `deleteOlderThan(ctx, tx, "task_log", "started_at", c.TaskLogTTL, now)`，而 `task_log.started_at` 是 `NOT NULL`（`schema.go:67`）⇒ 父**必然**进老化集合。
- **合成语句同一性实测**（F3 孪生夹具）：手跑上面那句 ⇒ `rows_affected=4`，与生产那一发 `toolcall_deleted=4` **相等**，幸存行集合同一个 `id=5` ⇒ 我引的这句就是生产跑的那句。

## §4 那一发"31 天父 ＋ NULL 子行"（孤儿／fkc 迁移前後／能否开库）

| 夹具 **F1**（`task-old` 31 天 ＋ **A：三列全 NULL**；`task-fresh` ＋ **F：三列全 NULL、父活着**） | `task_log` | `tool_call` | fkc | 逐字读数 |
|---|---|---|---|---|
| 造好 | 2 | 2 | **0** | `[]` |
| `RunRetentionOnce` #1 | 1 | 2 | **1** | `tasklog_deleted=1`／`toolcall_deleted=0`／`[tool_call|rowid=1|parent=task_log|fkid=0]` |
| retention #2 | 1 | 2 | **1** | `0`／`0` ⇒ ★**孤儿永不被第二轮选中** |
| 生产 `Open()` 迁移 → 再数 | 1 | 2 | **1** | `open_err=<nil>`、`schema_version=4`、`ForeignKeyAudit{Checked:true, OrphanRows:1, ChildTables:[tool_call]}`；孤儿行 `task_id=task-old` 一字未动 |

★**判**：孤儿**造得出**（⛔"造不出"）、**fkc 0→1**、**迁移既不清也不造**、**今天能开**、**第二轮甚至第 N 轮 retention 都删不掉它** ⇒"永生"**成立**（对三列全 NULL 的行）。★另：**第二次 `Open()`** 现跑 `audit.Checked=**false**` ⇒ `open.go:413` 的 `version==s.target` 早退被复认（票面硬提醒 1 引的 `:374` 就是这处）。

## §5 正控（硬闸）：三列都有值 ⇒ 必须和父一起清掉

| 夹具 | 读数 |
|---|---|
| **F2 正控**（`task-old` 31 天 ＋ `B`：**decided_at/started_at/ended_at 三列都＝40 天前**；`task-fresh` ＋ `G` 同形、父活着） | `tasklog_deleted=1`／`toolcall_deleted=**2**` ⇒ `task_log=1`、`tool_call=**0**`、fkc=**0**；迁移 `open_err=<nil>`、`audit{Checked:true, OrphanRows:**0**}`、迁移後 1/0／fkc=0 |
| **F3 三臂＋顺序**（`C` 只 decided_at／`D` 只 started_at／`E` 只 ended_at，均 40 天；`I`＝ended 40 天但 started/decided 29 天；`H`＝ended 29 天但 started/decided 40 天、父活着） | `toolcall_deleted=**4**`（C,D,E,I）、幸存**只有 `id=5`（H）**、fkc=**0** |

★**正控成立**：三列有值 ⇒ **与父同清**（F2 归零）；**三条 COALESCE 臂各自单独有值都能被读到**（C/D/E 全删）；**顺序＝ended_at 优先**（I 删、H 留，两条只差"哪个时间戳在外"）⇒ **尺能读出"正常老化"，孤儿不是尺坏了**。

## §6 ★第五条路径量完整（`startTask` 失败 → 随后 `startCall` 是否照写）

真 loop ＋ 真 store ＋ 真 SQL，三发（R2 在 `Journal` 接缝拒 `StartTaskLog`；R3 用**真 SQLite 错误**：`CREATE TRIGGER … BEFORE INSERT ON task_log … RAISE(ABORT,…)`，Go 侧零注入）：

| 读数 | R1（**对照：不注入**） | R2（接缝拒） | R3（真 SQL 拒） |
|---|---|---|---|
| loop 告警 | `warn=0` | `warn=1`：`WARN|agent: task_log open failed|err=probe 368a1: injected startTask refusal` | `warn=1`：`…err=memory: start task_log ad2d923f-…: constraint failed: probe 368a1: task_log insert refused (1811)` |
| 任务是否继续 | `completed`、工具**真跑 1 次**、2 请求 | **同左**（⛔ 停在失败处） | **同左** |
| `task_log`／`tool_call` 行数 | 1／1 | **0／1** | **0／1** |
| 子行有父吗 | `parents_in_task_log=1` | **`=0`（悬空）** | **`=0`（悬空）** |
| 子行三列 | `decided_at`/`ended_at` 有值、`started_at=<NULL>` | 同左（**有值**） | 同左（**有值**） |
| 造好後 fkc | **0** | **1** | **1** |
| 跑 retention（`Now=+31 天`） | — | `toolcall_deleted=**1**` ⇒ 0 行、fkc **0** | — |

★**判**：★**`startCall` 照写**（整链＝`startTask` 失败 → 只 `Warn` → 继续跑 → `executeCalls` `:671` 插行 → `task_id` 无父 → fkc 1；⛔ 只看最后一跳）。★**但**这条路径的孤儿**带 `decided_at`/`ended_at`** ⇒ 30 天后**会被清掉**（实测 `toolcall_deleted=1`）⇒ **"暂时孤儿"，⛔ 永生**。
★**且生产装配下整条链不可达**：`cmd/wisp/run.go:1080` `Journal: nil`，而 `journal.go:71-73`／`:126-128` 对 nil 是 `return 0`／`return nil`。实测（NJ）：`startTask_err=<nil>`、`startCall_rowID=**0**`、`finishTask_err=<nil>` ⇒ **`startTask` 根本不会失败、`startCall` 一行都不插**。

## §7 ★★**第六条路径（具名上报）：`tools.Bridge.book` 才是今天真正的"永生"写入者**

- **代码**：`internal/tools/bridge.go:1171 func (b *Bridge) book(...)`；`:1210-1221` 构造 `memory.ToolCall{TaskID, Seq, Tool, ArgsJSON, RiskLevel, **Decision**, **Outcome**, ErrorClass, CorrelationID, GrantID}` —— **三列时间戳一个都不设**；`:1225` `if _, err := b.j.InsertToolCall(ctx, tc); err != nil {` **把返回的 id 丢掉** ⇒ 这行**再也无法被 update**。
- **为什么它是"生产唯一"**：`cmd/wisp/run.go:764` `tools.New(tools.Options{…})` 的 `:769` 是 `Journal:  mem`；`:1072` `agent.Options{…}` 的 `:1080` 是 `Journal: nil`（注释逐字 "the loop's own tool_call rows are deliberately NOT booked: the bridge already writes the authoritative one"）。全仓 `InsertToolCall` 只有两个调用点（`journal.go:75`、`bridge.go:1225`）；`DecideToolCall`/`FinishToolCall` 的唯一生产调用者是 `journal.go:101`/`:118`，而它们的守卫 `:98-100`/`:112-114` 在 `j == nil` 时直接 `return` ⇒ **bridge 的行没有任何可能的更新者**。
- **实测 B1**（真 `Bridge` ＋ 真 `Store` 当 journal，一发 L0 调用）：行＝`decision=allow outcome=success decided_at=<NULL> started_at=<NULL> ended_at=<NULL>` ⇒ **decision/outcome 有值、三列仍全 NULL**；`RunRetentionOnce` ⇒ `tasklog_deleted=1`／`toolcall_deleted=**0**`、fkc **0→1**；**第二轮** `0`／`0`、fkc 仍 **1**（永生）；**正控**：把该行 `ended_at` 手动写成 40 天前 ⇒ 同一发 retention `toolcall_deleted=**1**`、fkc **0**（⇒ 不是尺坏，是 NULL）。
- ★**射程（⛔ 夸大）**：它造的是**"会变成孤儿"的永生子行**，不是"当场孤儿"——父还在时 fkc＝0。父一旦没（今天：`PurgePrivacy(PrivacyTaskLog)` → `privacy.go:177` → `dao_tasklog.go:134 DELETE FROM task_log`，用户一键；接线后：retention）⇒ 立刻变孤儿且**永不老化**。

## §8 "接线之后会怎样"：判据（⛔ 猜）

- **判据 1（谁永生）**：一行永生 ⇔ 插入后**没有任何一条 `UPDATE … decided_at/started_at/ended_at` 落到它身上**。三个现跑事实拼出该集合：①`started_at` **无写入者**（§2 附读数）；②`bridge.book` 的行**连 id 都不留**（§7）⇒ **100% 永生**；③`journal.startCall` 的行只在"终局写"没落地时永生（R1 显示正常收尾会写 `decided_at`+`ended_at`；`FinishToolCall` 在 store 关闭后逐字报 `memory: finish tool_call 1: memory: store is closed`，而 `journal.go:118` 是 `_ =` ⇒ **关库竞态**是第 ③ 类入口）。
- **判据 2（接线后会不会累积）**：**会，且单调不减**。永生集合 ∩ 父被删 ＝ 孤儿；retention 第二轮对它们 `tasklog_deleted=0`/`toolcall_deleted=0`（F1、B1）；全仓 `tool_call` 的删除语句只有三条（`dao_toolcall.go:134` 单行、`:149` 全清、`retention.go:197` 这句）⇒ **除"清空取证"外没有任何机制能让孤儿消失**。⇒ **接线后孤儿数 ≈ Σ(每个 ≥30 天任务里 bridge 写过的工具调用数)**，随时间**只增**。
- **判据 3（`AC#1` 的形会被这条改掉）**：〔甲〕改 `COALESCE` 判据／〔乙〕让 `startCall` 写 `started_at` —— **两形都只治 loop 侧**，对 `bridge.book` 的行**零效果**（它不经过 `startCall`）。要覆盖今天真正的写入者，落点只能是 **`InsertToolCall` 本身**（DAO 层给三列兜底）或 **`bridge.book`**。此条是**读数推出的射程**，选形仍归 `AC#1`。

## §9 门禁 rc（逐门分开、⛔ 相减、⛔ 走管道取 rc）

| 门 | 命令 | rc | 备注 |
|---|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **0** | 0 字节输出 |
| G2 | `GOFLAGS= go vet ./internal/memory/ ./internal/agent/ ./internal/tools/` | **0** | 0 字节输出 |
| G3 | `GOFLAGS= go test -count=1 ./internal/memory/` | **0** | `ok … 63.823s`（比 `366-a1` 的 20.5s 慢＝机器负载） |
| G4 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations` |
| G5 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`、两单位同选 61 条 |

门禁跑在**共享工作树**（当时 `cmd/wisp/*.go` 5 枚 `M` 是别人的活）；rc 全用 `cmd > out 2>&1; rc=$?`。探针三发在**导出树** `-overlay` 跑，rc 全 **0**（`out/memory-probe.log`、`out/agent-probe.log`、`out/tools-probe.log`）。

## §10 纪律格：提交号 · 逐枚 `.go`（应为 **0**）

- 提交号：**首笔只含本件**；本件首笔的哈希由紧随的一笔补录在下一行（⛔ amend、⛔ 改写历史）。
- 补录：首笔＝`bac7d5a29a7cf8ec79997d072ca3f6d673888469`（`A docs/evidence/s1/368-a1-census.md`，148 行）。
- **本腿落进仓的 `.go`＝0 枚**。判据三把（⛔ 只看一把）：①本腿开工 `09:33` 现跑 `git status --porcelain -- internal/ cmd/ scripts/`＝**空**，收尾时出现的 5 枚 `M cmd/wisp/*.go` 归**别的腿**（本腿从未对仓内任何路径执行写命令）；②本腿只写 `docs/evidence/s1/368-a1-census.md`；③探针在 `D:\tmp\wisp\20261024\368a1\probes\**`（仓外），经 `-overlay` 生效，导出树里**不落盘**。
- `add`/`commit` 两处显式 pathspec（逐字点名单一文件）；提交前 `git diff --cached --name-status` 核索引。

## §11 请编排者拍（⛔ 本腿自行假设）

1. ★**`AC#1` 的两形（甲/乙）都不覆盖 `bridge.book`**（§8 判据 3）⇒ 若只按票面两形落，接线后孤儿照旧累积。要治"今天的写入者"，形得落在 `InsertToolCall` 或 `book` 上——选形归你。
2. ★**第五条路径要不要留在票面**：它**今天不可达**（`Journal: nil`）、且造出的孤儿**会自愈**；票面把它写成并列路径会让人以为它是活洞。
3. ★**`366-a1` 的 `1/1` 口径**（§1.1 末）与 **`open.go:374`→`:413`**（§1.4）两处建议回写进票面/台账，⛔ 我改票面。
