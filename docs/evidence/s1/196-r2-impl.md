# 196-r2 实现件：`applyV4` 的 `DROP TABLE` 隐患收口（写腿）

**票**：196（隐患收口腿 `196-r2`）｜**写腿**｜**只 commit、⛔ push**
**本件是**：实现件的读数与判决，⛔ 裁决策（裁决归非实现者）。
**收的是**：`196-v1` 判语第 4 条残余（`docs/evidence/s1/196-v1-verdict.md:142`）——
> `applyV4` 的 `DROP TABLE` 今天靠 `foreign_keys` 关着才安全：若将来打开 FK 且库里有 `tool_call` 行，重建会在 DROP 处失败。

**环境**：HEAD `ef3ec49bd`｜`GOFLAGS=`｜只 `unset GOTMPDIR`（⛔ 动 `TMPDIR/TEMP/TMP`）｜⛔ `CGO_ENABLED=0`｜gofumpt 绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`
**开工前核**：`git diff --cached --name-status` ⇒ **空**（索引里无 `D`/`A`）；`git status --porcelain -- internal/memory docs/evidence/s1` ⇒ **空**（写面无别人的在飞改动）。

---

## §1 `applyV4` 现在怎么重建（逐行引，号取自改后 `internal/memory/schema.go`）

`ddlV4` 由 `buildDDLV4()`（`:294`）**生成**（CHECK 的值从 `taskStateRoster` 拼，⛔ 第二份手抄），`:292 var ddlV4 = buildDDLV4()`，**5 句**：

| # | 行 | 语句 |
|---|---|---|
| 1 | `:303` | `CREATE TABLE task_log_v4 (… state TEXT NOT NULL CHECK (state IN (…)) …)` |
| 2 | `:316` | `INSERT INTO task_log_v4 (…) SELECT … FROM task_log`（搬行） |
| 3 | `:321` | `DROP TABLE task_log;` |
| 4 | `:322` | `ALTER TABLE task_log_v4 RENAME TO task_log;` |
| 5 | `:323` | `CREATE INDEX idx_task_log_started ON task_log(started_at);` |

`applyV4`（`:438`）：`:440` 新护栏 → `:443` 名册普查（`taskLogStatesOutsideRoster`，越册则点名拒、什么都不动）→ `:453` 逐句 `splitSQLStatements(ddlV4)` 执行。整步在一个事务里（`open.go:394` `applyMigrationStep` → `BeginTx`）。

★**"别的表对 `task_log` 的引用"当时怎么处理＝一处都没处理**：重建前不删子行、不搬子行、不 `PRAGMA foreign_key_check`、不碰 `legacy_alter_table`。而**引用确实存在**：`ddlV1` 的 `:82` 逐字
`  task_id        TEXT NOT NULL REFERENCES task_log(id),`（表 `tool_call`，`ddlV1:80`）。
⇒ `DROP TABLE task_log` 是**第 3 句**，前两句已经把新表建好、行搬完（同一事务内）。

## §2 ★"开 FK ＋ 有 `tool_call` 行"那发到底炸不炸（逐字，全部我现跑）

探针件 `zz_196r2_probe_test.go` 只存在于**仓外导出树** `/d/tmp/wisp/20261023/196r2/tree1|tree2`；夹具＝生产 `ddlV1+ddlV2+ddlV3` 建的真 v3 库（`schema_version=3`），塞 3 枚 `task_log` 行 ＋ **1 枚 `tool_call` 行**（`task_id` 从 `task_log` 里 `SELECT` 出来，⛔ 手写 id——孤儿行不会让 DROP 违规，那样测的是错的东西）。

| 探针 | 问 | 逐字读数 |
|---|---|---|
| P1 | 生产 `dsn()` 上 `PRAGMA foreign_keys` 是几 | **`0`** |
| P2 | `applyV4` 能不能自己关掉 FK？ | 事务内 `PRAGMA foreign_keys=ON` 后读回 **`0`**；同一连接、**事务外**同一条语句后读回 **`1`**（正控）⇒ ★**事务内是 no-op** |
| P3 | 裸 `DROP TABLE` 在 FK 开时 | 有 1 枚子行 ⇒ `constraint failed: FOREIGN KEY constraint failed (787)`；子行清零后 ⇒ `<nil>` |
| P4a | **生产 `Open()`**，FK 关（今天） | `<nil>`；`schema_version=4 task_log=3 tool_call=1 leftover_v4=0` |
| P4b | **生产 `Open()`**，FK 开（tree2＝`dsn()` 加 `&_pragma=foreign_keys(1)`，唯一差别） | `memory: migration 3->4 failed (database restored from backup is at …\backup\wisp.db.bak-3-4): apply: constraint failed: FOREIGN KEY constraint failed (787)` |

★★**结论：今天不炸**（P4a 绿），判语那句"**将来**"逐字成立 ⇒ 按派单**继续**（⛔ 停手上报）。
★★**但判语只说了"会在 DROP 处失败"，没说的是：那条报错既不点名表、也不点名行、也不说"库没被动过"**（P4b 逐字就是 `(787)`）。这就是本腿要收的东西。

## §3 `PRAGMA foreign_keys` 今天是谁设的、设成什么（逐字行号）

★**没有人设**。全仓唯一挂 pragma 的地方是 `internal/memory/open.go:146 func dsn`：
```
:148  "?_pragma=busy_timeout(" + itoa(pragmaBusyTimeoutMS) + ")" +
:149  "&_pragma=journal_mode(WAL)" +
:150  "&_pragma=synchronous(NORMAL)" +
:151  "&_pragma=wal_autocheckpoint(" + itoa(pragmaWALAutocheckpoint) + ")"
:152  if forWriter {
:153      d += "&_txlock=immediate"
```
四枚 pragma **没有 `foreign_keys`** ⇒ 走 SQLite 默认＝**关**（P1 实测 `0`）。尺：`git grep -n '_pragma=' -- internal/memory` ⇒ 只有上面 4 行。
尺：`git grep -n 'foreign_keys' -- '*.go' '*.sh' '*.yml' '*.ps1' '*.toml'` ⇒ **本腿之前 0 命中**；本腿之后只剩我自己新增的 7 行（`schema.go:361/368/401/407/411/416/421`）。独立旁证两枚：`docs/evidence/s1/326-append-updated.md:147`、`326-v1-verdict.md:78`（都复跑出 0 命中 rc=1）。

**"什么时候会开 FK"**：★**仓里没有任何一张"将来要开 FK"的票**（⛔ 我编不出来）。★但有三条**"FK 声明了不生效"的欠账**，那才是有人会去开它的现实路径——
- `docs/reports/pending-and-issues.md:11196`（`166-a1` 普查）：`foreign_keys` 零命中 ⇒ 改名/删行**不被数据库拦**，会**静默**断取证链；
- `326-append-updated.md:147` ＋ `326-v1-verdict.md:78`：`ddlV3` 的 `REFERENCES session(id)` 今天不生效、孤儿轮次可达；
- 本票残余三处：`196-v1-verdict.md:142`、`196` 票 `:59`、`pending-and-issues.md:23130`。
⇒ ★**"为了 `REFERENCES` 名副其实而去开 FK"是完全合理的下一步**；一旦有人这么做，本腿的护栏就是那道"说话的门"。

## §4 选了哪一形、为什么它最有价值

**选〔甲〕：在 `applyV4` 里显式声明"重建期间 `foreign_keys` 必须是关的"（断言 ＋ 可读理由）**，`schema.go:368`（查询常量）＋ `:398 requireForeignKeysOff` ＋ `:440` 调用。

为什么是它（不是"风险最小"，是"最有价值"）：
1. ★**〔甲〕是唯一能在本腿写面内、在破坏性语句之前开火的形**。P2 已证：`PRAGMA foreign_keys` 在事务内是 **no-op** ⇒ `applyV4` **不可能**自己关掉 FK；真正"关 FK"只能落在 `open.go` 的 `migrate()`/`dsn()`（**开事务之前**），**不属本腿写面**（见 §8 上交）。
2. ★**它正是主流 harness 的做法**：SQLite 官方 12 步 ALTER 的**第 1 步**就是"若 FK 开着，先关"。`ddlV4` 的注释自称走 12 步（`schema.go:283-284` 逐字 "The rebuild follows SQLite's documented 12-step / procedure"），**实际漏了开头的第 1 步（关 FK）与结尾的 `foreign_key_check`** ⇒ 〔甲〕把这条**恰好成立的前提**改成**声明出来的前提**：失败点从第 3 句挪到第 1 句，报错从 `(787)` 换成点名 pragma ＋ 子表 ＋ 修法 ＋ "什么都没动"。
3. ★**〔乙〕不选，三条硬理由**：①它得**枚举子表**（今天只有 `tool_call`，明天多一张引用方就**静默复发**）——比护栏更坏，因为它**看起来**治好了；②它要往 `ddlV4` 加语句（5 句 → N 句），既撞本腿硬边界 2（⛔ 改 `ddlV4` SQL 语义），也撞既有钉测 `TestTaskState196DDLV4KeepsTheShape`（`want 5`）；③它只买到"FK 开着也能重建"，而"FK 为什么会被开"这一层它一点没答。
4. ★**〔丙〕不单独选**：单测只能钉住**已经会发生**的失败；今天钉就是钉 `(787)`，等于把不可读写进契约。〔甲〕把〔丙〕顺带做了——§5 第二枚用例钉的就是〔甲〕的拒绝文本。
5. **代价，写明**：零 `tool_call` 行的库本来能"靠运气"过，现在也**拒**。故意如此——一次成功取决于"子行恰好为零"的重建，不值得跑。
6. **顺序，写明**：FK 门槛在名册普查**之前**（`:440` vs `:443`）。它是连接级的、且是官方第 1 步；两种都错时先听环境问题，再听数据问题。

## §5 新用例逐名 ＋ 为什么它**今天**就响

`internal/memory/taskstate_196r2_test.go`（新，213 行；尺＝`wc -l`）：

1. **`TestTaskState196LegacyV3WithToolCallRowsStillOpens`**（`:63`）——**绊线**。走**生产 `Open()`**，夹具＝v3 库 ＋ 1 枚 `tool_call` 行；断言 `schema_version=4`、子行仍在且仍指向**同一**父、`task_log_v4` 残留 0。
   ★**今天响的理由**：它把 `dsn()` **今天**的状态变成判据——`dsn()` 里加一枚 pragma（那件"为了让 `REFERENCES` 名副其实"的下一步），这枚用例立刻红。§6 的 R2/R3 就是我先把它打红的读数，⛔ 空口径。
2. **`TestTaskState196RefusesRebuildWithForeignKeysOn`**（`:114`）——**护栏的牙**。同一份夹具跑两遍**真的** `applyMigrationStep`（`migrate()` 调的就是它），唯一差别是 pragma：控制半（FK 关）必须绿且 `tool_call=1` 活着；主体半（FK 开）必须**拒**、必须点名（`foreign_keys` / `tool_call.task_id REFERENCES task_log(id)` / `Nothing was modified.` 三条断言，`(787)` 一条都不占），且拒绝后**一个字节没动**（watermark 仍 3、`task_log` 4 行、`tool_call` 1 行、残留 0）。
   ★**今天响的理由**：§6 R3——把护栏那 3 行删掉，这枚用例当天红，报错退回 `(787)`。

## §6 反向判据三读数 ＋ 还原核（全部仓外导出树，工作树零突变）

| 读数 | 树 | 形状 | rc | 逐字 |
|---|---|---|---|---|
| **R1** | `tree3` | 交付形状（HEAD 导出 ＋ 本腿两件），FK 关 | **0** | 两枚全 `PASS`；拒绝文本见 §5.2 断言 |
| **R2** | `tree4` | `tree3` ＋ `dsn()` 加 `&_pragma=foreign_keys(1)`（唯一差别） | **1** | 两枚全 `FAIL`；`Open()` 现在**可读地拒**：`apply: memory: refusing to rebuild task_log with PRAGMA foreign_keys=ON: tool_call.task_id REFERENCES task_log(id), …` ⇒ ★**失败点前移实证** |
| **R3** | `tree5` | `tree4` ＋ 删掉 `schema.go:440-442` 那 3 行（＝改前形状） | **1** | 两枚全 `FAIL`；报错退回 `apply: constraint failed: FOREIGN KEY constraint failed (787)` ⇒ ★**护栏就是那 3 行换来的**（单变量：`diff` 只有这 3 行） |

**还原核**：`sha256sum`（64 hex，`sha256sum` 算的，⛔ `git hash-object` 的 SHA-1）前＝后——`schema.go 6b985a280092becbfc0c72bcb7e932339a0e08f3a62272ea591a8c3f61084094`、`taskstate_196r2_test.go 9ac9ef48a7e590c78560f0864732c331673445529b0e1e12dbc6a499e8bb627f`；`cmp` 仓内 vs 导出树 **rc=0**；`git status --porcelain -- internal/memory` ⇒ ` M schema.go` ＋ `?? taskstate_196r2_test.go`（恰是本腿意图，⛔ 多一枚）。
★诚实一条：**R2 里绊线用例也红了**，因为它的**控制半**假定 `dsn()` 不设 pragma——这本身是第二条护栏（前提变了就喊）。

## §7 四数 ＋ 门禁 rc

**四数**（尺＝`go test -count=1 -v ./internal/memory/`；四数要 `-v` 才数得出；跑的是工作树，HEAD `ef3ec49bd`）：

| | rc | RUN | PASS_TOP | FAIL_TOP | SKIP_TOP |
|---|---|---|---|---|---|
| **改前** | 0 | 91 | 54 | 0 | 2 |
| **改后** | 0 | **93** | **56** | 0 | 2 |

Δ＝**+2 RUN／+2 PASS**，恰是两枚新用例（⛔ 子测试、⛔ 假红），FAIL 前后都是 0。

**门禁**（逐把真跑、rc 直取⛔ 走管道）：`go build ./...` **0**｜`go vet ./internal/memory/` **0**｜`gofumpt -l`（绝对路径）两件 **0 枚列出**，且**正控**一枚故意错格式的文件**被列出 1 枚**（⛔ 把 0 当结论而不做正控）｜`sh scripts/d22scan.sh` **0** `clean`（internal 632／cmd 148 Go 文件）｜`sh scripts/check-path-length-budget.sh` **0** `VERDICT GREEN`。

## §8 边界核对 ＋ 具名上交（⛔ 自填）

**四条硬边界逐条核**：① `ddlV1/ddlV2/ddlV3` 未动（`git diff` 对 `schema.go` **+73／-0**，零删除行）；② `ddlV4` 的 SQL 语义未动（`:292-324` 逐字节同 HEAD，5 句仍是 5 句，`TestTaskState196DDLV4KeepsTheShape` 绿）；③ **12 枚名册未动**（`internal/memory/taskstate.go` **零 diff**，`git status` 里没有它）；④ `docs/PLAN.md`／`docs/specs/**` 未动。
**写面**：`internal/memory/schema.go`（+73/-0，只 `applyV4` 邻域）＋ `internal/memory/taskstate_196r2_test.go`（新）＋ 本件。⛔ 触 `open.go`／`cmd/**`／`scripts/**`／golden／allowlist／三枚冻结件。

**上交两件（不属本腿写面，⛔ 自己改）**：
1. ★**真正的修法在 `open.go`**：按 SQLite 官方第 1 步，在 `migrate()` **开事务之前**把 `PRAGMA foreign_keys=OFF` 钉在迁移连接上（`dsn()` 加一枚 `_pragma=foreign_keys(0)` 即可），重建后按第 12 步 `PRAGMA foreign_key_check`。那两处都在 `open.go`，**归别的腿**；本腿的护栏与它**不冲突**（护栏在 FK 关时是 no-op）。
2. ★**结尾那道 `foreign_key_check` 本腿没加**，理由具名：今天 FK 关着，库里**可能**已有孤儿 `tool_call` 行（`166-a1`／`326` 那两条欠账），把它做成硬门槛会让那类老库**打不开**——正是派单的停手条件③。要么另开票配"孤儿先清后校"，要么等第 1 条落地时一起做。
