# 票 196 真修法腿 `196-r3` — 迁移连接在 `BeginTx` 之前关 FK

写腿交付件。⛔ 不是裁决（裁决归非实现者 `196-v2` 那一侧）。
一句话：**迁移链改成走一条迁移专用 DSN**，它在建连时就把 `PRAGMA foreign_keys` 声明成 `OFF`（SQLite 官方 12 步的第 1 步，且早于 `BeginTx`）；**共用 `dsn()` 与两个业务池一字未动**。

## 1 ★`migrate()` 从哪一行开事务（逐字，现跑）

尺 ＝ `Select-String` 于 `internal/memory/open.go`（改后）。

| 锚 | 行 | 逐字 |
|---|---|---|
| `dsn()` | 154 | `func dsn(path string, forWriter bool) string {` |
| `dsnPragmas()` | 161 | `func dsnPragmas() []string {` |
| `migrationDSN()` | 183 | `func migrationDSN(path string) string {` |
| `migrate()` | 364 | `func (s *Store) migrate() error {` |
| 迁移建连 ① | 365 | `boot, err := sql.Open(DriverName, migrationDSN(s.dbPath))` |
| 迁移建连 ②（每步重开） | 402 | `boot, err = sql.Open(DriverName, migrationDSN(s.dbPath))` |
| `applyMigrationStep()` | 461 | `func applyMigrationStep(db *sql.DB, step migrationStep, logger *slog.Logger) error {` |
| **★迁移事务起点** | **465** | **`tx, err := db.BeginTx(ctx, nil)`** |

★改前（`HEAD` ＝ `9909404db`）：`dsn` 146 ／ `migrate` 293 ／ **`BeginTx` 394**，两处迁移建连在 294 与 331，且都用 `dsn(s.dbPath, true)`。
⇒ **"迁移事务从哪一行开始" ＝ `open.go:465` 的 `tx, err := db.BeginTx(ctx, nil)`（改前 394）**；本腿要把 FK 关在这一行**之前**，所以只能落在**建连**这层（DSN），⛔ 不能落在事务里（`PRAGMA foreign_keys` 在事务内是 no-op，`196-r2` 已量过、写在 `schema.go:361-367`）。

## 2 ★改前：FK 开 ＋ 有 `tool_call` 行 ⇒ 生产 `Open()` 打不开（逐字）

对照树＝**仓外叠加**（`-overlay` 把仓外 `open-before-fkon.go` 顶掉 `internal/memory/open.go`，唯一改动是 `dsn()` 末尾多一行 `"&_pragma=foreign_keys(1)"`；工作树零改动、无突变窗口）。
夹具＝`seedLegacyV3`（生产 `ddlV1/ddlV2/ddlV3`）＋ `addLegacyToolCallRow`（子行 `task_id` 取自 `task_log`）。

| 发 | 命令 | rc | 逐字（截到库自述部分） |
|---|---|---|---|
| A 有护栏（现状） | `-run TestTaskState196LegacyV3` | **1** | `memory: migration 3->4 failed (...): apply: memory: refusing to rebuild task_log with PRAGMA foreign_keys=ON: tool_call.task_id REFERENCES task_log(id), so the DROP TABLE the v4 rebuild needs would abort on SQLite's "FOREIGN KEY constraint failed" and name neither a table nor a row. ... Nothing was modified.` |
| B 护栏被仓外摘掉 | 同上 ＋ `schema.go` 里 `requireForeignKeysOff` 首行 `return nil` | **1** | `... apply: constraint failed: FOREIGN KEY constraint failed (787)` |

★当时状态：**A 的两枚夹具都红**（含零子行那枚 `TestTaskState196LegacyV3DatabaseOpensAndReads`）—— 护栏按设计"连靠运气能过的库也拒"，所以 **FK 一开，`Open()` 一律打不开**；B 证明护栏底下那条 `(787)` 仍然在。
⇒ 这正是 `196-v2` 的裁决：护栏把"(787)"变成"第 1 句可读拒、零语句执行"，但**隐患是被声明并前移、⛔ 被消除**。

## 3 ★选了哪一形 ＋ 为什么不影响业务连接

形：**迁移专用 DSN**（`migrationDSN`），不是全局关，也不是"迁移里关完再恢复"。

```
dsn()            = buildDSN(path, forWriter, dsnPragmas())                  ← 业务池用它，一字未动
migrationDSN()   = buildDSN(path, true, migrationPragmas(dsnPragmas()))     ← 只 migrate() 的两处建连用
migrationPragmas = 把 shared 里每一条 foreign_keys 条目删掉，再追加恰好一条 foreign_keys(OFF)
```

★★**为什么⛔ 是"追加一条覆盖"**（这是我改形的原因，读数为证）：
modernc 会**先排序再逐条执行** `_pragma`（`sqlite.go:433-449`，v1.59.0：`busy_timeout` 排最前、其余按字典序），**不按 DSN 顺序**。实测：

| shared 里的写法 | "追加 OFF 覆盖"读完 | 本腿的形读完 |
|---|---|---|
| `foreign_keys(1)` | 0（侥幸） | 0 |
| `foreign_keys = 1` | 0（侥幸） | 0 |
| `foreign_keys(ON)` | **1** | 0 |
| `foreign_keys(on)` | **1** | 0 |
| `foreign_keys(true)` | **1** | 0 |
| `foreign_keys(yes)` | **1** | 0 |
| `foreign_keys(0)` | 0 | 0 |
| `foreign_keys(OFF)` | 0 | 0 |

⇒ **"追加覆盖"是一场拼写彩票：8 种写法里 4 种读完仍是 1**。另量：`(1)` 与 `(0)` 并存时 **200/200 条新连接全读 1**（`0` 排在 `1` 前 ⇒ 后应用）。**恰好一条**才与排序无关。
★**为什么不影响业务连接**（三层，逐层可核）：
1. `PRAGMA foreign_keys` 是**连接级开关、⛔ 写进文件**（`journal_mode` 才是会写进文件的那种），所以迁移连接说的这句对别的连接**没有记忆**；
2. 迁移连接用完即关（`migrate()` 的 `boot.Close()`：**396** 每步备份前、**413** 收尾、371/375/379/388/407 是各条错误返回；改前分别为 325／342），业务池在它**之后**用 `dsn()` 另开；
3. `dsn()` 输出**逐字节等于改前**（用例钉住），且**一条 `foreign_keys` 都不带** ⇒ 用户在业务连接上设什么就是什么，本包不替他决定。

## 4 ★改后：FK 开 ＋ 有 `tool_call` 行 ⇒ `Open()` 绿

同一发、同一夹具，对照树换成 **"将来的那一改"**：`dsnPragmas()` 里加一条 `foreign_keys(1)`（⚠ 这不是本腿的改动，是**故意注入的未来**，让"共用的 dsn 开了 FK"这件事真的发生）。

| 读数 | 值 |
|---|---|
| `Open()` err | `<nil>`（**rc=0**） |
| `schema_version` | **4** |
| `task_log` / `tool_call` / 子行仍指向原父 | **4 / 1 / 1** |
| `task_log_v4` 残渣 | **0** |
| 迁移连接实际 DSN | `..._pragma=wal_autocheckpoint(1000)&_pragma=foreign_keys(OFF)&_txlock=immediate` |
| 同一棵树上业务连接 DSN | `..._pragma=wal_autocheckpoint(1000)&_pragma=foreign_keys(1)` |

★注意最后两行**同时**出现：注入的 `foreign_keys(1)` 进了业务 DSN，**没进**迁移 DSN —— 这就是分流。

## 5 ★业务连接 FK 仍按用户设的值

| 场景 | 读数 |
|---|---|
| 上面那发（注入 `foreign_keys(1)`）之后 `s.reader` / `s.writer` | **1 / 1** |
| 原样树（`dsn()` 不设，SQLite 默认）`s.reader` / `s.writer` | **0 / 0** |
| 用例：`dsn(...)+"&_pragma=foreign_keys(1)"`（调用方自设） | **1** |

⇒ **⛔ 被本改动永久关掉**：调用方要 1 就是 1。

## 6 ★新用例（`internal/memory/open_196r3_test.go`，逐名）

| 用例 | 钉什么 | 为什么今天响（⛔ 空转） |
|---|---|---|
| `TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys` | `dsn()` 逐字节等于改前字面量 ＋ 业务 DSN 不带 `foreign_keys` ＋ 调用方设 1 读到 1 | 谁往 `dsn()`/`dsnPragmas()` 加 FK（无论 0 还是 1），字符串比对与控制读数立刻红 ⇒ 硬边界 1 |
| `TestMigrationPragmas196r3ExactlyOneForeignKeysOff` | 5 个输入（含 `foreign_keys(1)`、`(ON)`、`= ON`、只有开关）⇒ 输出**恰好一条**且 `= OFF`，其余条目逐条原序保留 | 摘掉"过滤＋恰好一条"这一跳 ⇒ 5 个子用例全红（见 §7 M1） |
| `TestMigrationDSN196r3SwitchesForeignKeysOffOnTheWire` | 对照连接读 1、迁移连接读 0；再对 **8 种拼写**逐一开真连接读 | ⛔ 只比字符串：它证明驱动真的照做了，且"追加覆盖"的 4 种拼写会在真连接上读到 1 |

★每条都**不依赖** `Open()` 之外的东西；夹具沿用 `seedLegacyV3`／`sealableTempDir124`，⛔ 手写 DDL。

## 7 ★反向判据：三读数

| 读数 | 变异体（仓外 `-overlay`，⛔ 动工作树） | rc | 逐名报红 |
|---|---|---|---|
| **变异红 M1** | `migrationPragmas` 首行 `return shared`（★摘掉新加的那一跳），其余原样 | **1** | `TestMigrationPragmas196r3ExactlyOneForeignKeysOff` 全 5 子例（`open_196r3_test.go:105` 调用数 0≠1、`:110` 末条不是 `OFF`）＋ `TestMigrationDSN196r3SwitchesForeignKeysOffOnTheWire`（`:183`，8 种拼写里 6 种读到 1） |
| **变异红 M2** | M1 ＋ `dsnPragmas()` 注入 `foreign_keys(1)`（未来那一改真的发生） | **1** | `TestTaskState196LegacyV3WithToolCallRowsStillOpens`（`taskstate_196r2_test.go:70`）逐字回到 §2-A 那句可读拒 |
| **变异红 M3** | 只用"追加覆盖"（`return append(shared, OFF)`，保留 shared 那条）＋ 注入 FK | **1** | `TestDsn196r3…`（`:46`／`:49`，注入本身违边界 1）＋ `:105`（2～3 条 `foreign_keys`）＋ `:183`（`ON`/`on`/`true`/`yes` 四种读到 1） |
| **还原绿** | 无 overlay（工作树原样） | **0** | `ok github.com/CarlosShao/wisp/internal/memory 24.5s` |
| **原样绿** | `HEAD` 导出树（`git archive` 到仓外，`9909404db`）跑同一包 | **0** | `ok github.com/CarlosShao/wisp/internal/memory 21.807s`（⇒ HEAD 本来就绿，§2 的红⛔ 基线抖动） |

★**还原核**：先问"谁算的、什么算法、几位" —— 下表是 **SHA-256（64 位十六进制）**，由**本机 PowerShell 7 `Get-FileHash -Algorithm SHA256`** 算；★**⛔ 拿 `git hash-object` 去比**（那是 SHA-1、40 位，`A1149②` 的坑）。

| 文件 | SHA-256（改后，工作树） |
|---|---|
| `internal/memory/open.go` | `eb3cdbdae3ff59cddd1011101d3e38467edc2632b8007076a663c1354c367392` |
| `internal/memory/open_196r3_test.go` | `20be5657e7052b5de9405487fd6eb1570e5224c79b3e24662ff25a2cacbb7cd4` |

★**突变窗口：无**。全部变异体都是**仓外文件 ＋ `-overlay`**，工作树在整场里没被改过（`git status --porcelain -- internal/memory` 全程只有本腿这两个文件），所以上表 hash 就是"还原态"的 hash。

## 8 门禁 rc ＋ 四数

| 门禁 | 命令 | rc |
|---|---|---|
| 建 | `go build ./...` | **0** |
| 静态 | `go vet ./...`（整模块） | **0** |
| 包测 | `go test ./internal/memory -count=1` | **0**（24.5s） |
| 新用例 | `go test ./internal/memory -run 196r3 -count=1` | **0** |
| 格式 | `gofumpt -l`（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`） | **0**（无输出） |
| D22 | `sh scripts/d22scan.sh` | **0**（`clean - no D22 ban violations`） |
| 路径预算 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0**（`VERDICT GREEN`） |

四数（尺 ＝ `git diff --numstat`，★⛔ 用 `grep -cE '^[+-][^+-]'` 数行，`196-v2` 已证会漏空白行）：

| 文件 | 增 | 删 |
|---|---|---|
| `internal/memory/open.go` | **79** | **8** |
| `internal/memory/open_196r3_test.go`（新件） | **187** | **0** |
| `docs/evidence/s1/196-r3-impl.md`（本件） | 计在提交里 | 0 |

## 9 边界自查（逐条）

1. ★⛔ 全局关 FK：走迁移专用 DSN；`dsn()` 逐字节未变，用例钉住 ⇒ 业务连接不受牵连（§5 读数）。
2. ★⛔ 加第 10 步硬门槛：本腿**一行都没加** `foreign_key_check`（`196-v2` 已证老库可能有孤儿行，那会让这类库打不开）；第 ① 件之外的第 ② 件按派单**单独立票**。
3. ★⛔ 动 `ddlV1`–`ddlV4`、⛔ 动 12 枚名册、⛔ 动 `196-r2` 的护栏：`requireForeignKeysOff`／`pragmaForeignKeysQuery` 一字未动（§2-B 的"摘护栏"只存在于**仓外叠加**里，是读数手段，⛔ 是改动）。
4. ★⛔ 动 `docs/PLAN.md`／`docs/specs/**`：未动。
5. `internal/memory/schema.go`：**未动**（连注释都没改；`schema.go:386` 那句"relies on dsn() never setting the pragma"现在读起来略旧 —— 它描述的是 r2 当时的事实，本腿没被授权改它，留给裁决者定是否补一句）。

## 10 提交与四数对账

| 项 | 值 |
|---|---|
| 提交 | `da38c2bf4`（父 `61220f85a`，只动了 `docs/reports/pending-and-issues.md`，未碰本腿写面） |
| `git show --stat` | **3 files changed, 413 insertions(+), 8 deletions(-)** |
| 对账 | 79 ＋ 187 ＋ 147 ＝ **413** ✓；删除 8 ＝ 8 ✓（★两个方向都合得上，⛔ 只靠一个数自洽） |
| 提交后包测 | `go test ./internal/memory -count=1` ⇒ **rc=0** |
| 路径预算（新件入列后再跑） | `rc=0`，分母 17645→**17647**（＋2 件：新用例、本件），`VERDICT GREEN`、roster=61 |
| 提交字节 ＝ 量过的字节 | `git show HEAD:<p>` 经 `sha256sum` 与 `Get-FileHash -Algorithm SHA256` **三件逐一相等**（同为 64 位十六进制）⇒ 读的就是提交的那份 |
