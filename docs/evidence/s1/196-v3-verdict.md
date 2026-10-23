# 票 196 真修法终裁腿 `196-v3`：`196-r3` 够不够格把「迁移连接在建连层关 FK」做成真修法

**票** 196（真修法终裁腿 `196-v3`，**非实现者**）｜只 commit、⛔ push
**被裁** `196-r3`：`da38c2bf4`（码＋测试＋件，`3 files changed, 413 insertions(+), 8 deletions(-)`）＋`c346a88de`（件 +11）＋`889373f93`（件 2 行改现量）
**锚** 开工 HEAD `889373f93`；落笔时 HEAD `86ff73bdf`（编排者只动台账）⇒ `git diff 9909404db..HEAD -- internal/memory` 只有 r3 那一笔
**环境** Go 1.27.1 windows/amd64｜`GOFLAGS=`｜只 `unset GOTMPDIR`｜⛔ `CGO_ENABLED=0`｜gofumpt `/d/work/base/gopath/bin/gofumpt.exe`
**尺的射程**：每条读数标「哪棵树／哪把尺／量的是代码还是注释」。探针 `.go` **只在仓外** `D:\tmp\wisp\20261023\196v3\**`，⛔ 进仓、**工作树零突变窗口**。树＝`treeHead`（HEAD 导出）／`treeInject`（`dsnPragmas` 注入 `foreign_keys(1)`＝"将来那一改"）／`treeM1`／`treeM2`／`treeM3`／`treeM3b`／`treeM4`／`treeGuardOff`／`treeBase`（`9909404db`）／`treeFmt`（只给 gofumpt 甲尺）

## §0 裁决一句话
★★**够格翻。** 迁移连接确实在**建连层**、`BeginTx` 之前把 FK 关掉：我把 `foreign_keys(1)` 注入 `dsnPragmas()`（＝故意的"将来那一改"）后，生产 `Open()` 仍 `err=<nil>`、`schema_version=4`、`task_log=4`／`tool_call=1`／子行仍指原父、`task_log_v4` 残渣 0。**业务连接一根毫毛没动**（注入树 `s.reader`/`s.writer`＝**1/1**、原样树 **0/0**、`dsn()` 输出**逐字节未变**且零 `foreign_keys`）。边界零 diff、门禁五门 rc=0。
★★**它推翻我的技术建议，成立。** modernc 的排序我**逐行核了源码**（`sqlite.go:439`）也**逐条量了**：`(1)`＋`(0)` **200/200 读 1**；8 拼写里"追加覆盖"**4 种**读完仍是 1。⇒ **"删掉 shared 里每条 `foreign_keys`、再追加恰好一条 OFF" 是对的形状**，⛔ 拼写彩票。
★**但⛔ 够格注销隐患**：第 10 步 `foreign_key_check` 基线化仍欠（派单已单独立票）；且它**留下两处已过时的护栏叙述**（`schema.go:386` 注释＋`schema.go:416-421` 报错串）——**第二处比它自报的第一处更该改**（§8）。

## §1 必判八项，逐项判语（读数在后八节）

| # | 判语 | 关键读数 |
|---|---|---|
| 1 票面前提新鲜度 | **真**；两处自报尺误都**改对**（旧件 `17646`／`342`／`406` → 现 `17647`／`396`／`413`，两项我复量） | 现量 `open.go:154/161/183/199/212/220/364/365/402/461/465`、`boot.Close() 396/413`；改前 `9909404db`：`dsn 146`／`migrate 293`／建连 `294/331`／`BeginTx 394`／`Close 325/342` |
| 2 modernc 先排序 | **成立**（源码＋200 条＋8 拼写三发） | §3 |
| 3 改后 `rc=0` 与行数 | **复现**，含子行指原父、残渣 0 | §4 |
| 4 业务连接没被连累 | **复现**，★我加了"连接级、⛔ 文件级"正控 | §5 |
| 5 反向判据 M1/M2/M3 | **三枚全红、逐名逐行复现**；★`M3` 复现并**加一枚 `M3b`** | §6 |
| 6 边界零 diff | **逐项 0/0** | §7 |
| 7 门禁＋格式甲乙 | **五门 rc=0**；甲尺 0 行、乙尺 3 行**无一属本腿** | §9 |
| 8 `schema.go:386` | **要改**——而且**不止**这一句（`416-421` 更该改） | §8 |

★诚实一条（尺）：派单说「`dsn()` 一字未动」——**源码**不是（被重构成 `dsn()=buildDSN(path,forWriter,dsnPragmas())`），**输出**才逐字节未变（§5 另发核）。r3 件里写的是"输出"，**措辞对**；派单那句按字面读会错。

## §2 够不够格翻 ＋ 缺什么
**够格翻**：①真修法落在 `BeginTx` **之前**（DSN 层），⛔ 事务内（`PRAGMA` 在事务内 no-op，r2 已量）②⛔ 全局关（`dsn()` 零 `foreign_keys`）③⛔ 第 10 步硬门槛（全仓 `.go` 零 `foreign_key_check`，`git grep` **rc=1 无命中**）④三条新用例**摘一跳就红**⑤边界零 diff⑥门禁五门 0⑦推翻我的建议**有源码有实测**，方向更对。

**缺什么（按优先级）**
- (a) ★**第 10 步 `foreign_key_check` 基线化**仍欠——本腿按派单"立法不修"，必须**真的立成票**，⛔ 让"隐患已消"盖住它。r3 修的是「重建不被 FK 挡住」，⛔ 「重建后有没有造孤儿」。
- (b) ★**`schema.go` 两处叙述过时**（`386` 注释、`416-421` 报错串），见 §8：留着**会招回真回归**。
- (c) ★**接线没有牙**：把 `open.go:365/402` 的 `migrationDSN(...)` 改回 `dsn(s.dbPath,true)`（树 `treeM4`），**整包 `go test ./internal/memory/` 仍 rc=0**——今天 `dsn()` 不带 FK，两种写法行为相同，没用例能分辨。⚠ 缺口**有界**：那一天 `dsn()` 真带 FK，r2 的生产用例 `TestTaskState196LegacyV3WithToolCallRowsStillOpens` 会在 `:70` 红（M2 实测）。⇒ 可修可不修、⛔ 当缺陷，但票面别把"有用例钉住"说满。
- (d) 件里 §7 **M2 的红名单不全**：只写 `taskstate_196r2_test.go:70`，实际注入本身还红 `open_196r3_test.go:46/49/105/110/169/183`（`treeM2` 复现）。读数不错、名单不全。★极轻一条：`isForeignKeysPragma` 只认 `foreign_keys(` 与 `foreign_keys=` 两形，裸 `foreign_keys`（读语句）会活下来；无害（裸形按前缀先排、OFF 后落）。

## §3 ★「modernc 先排序再执行 `_pragma`」——我自己核的读数
- **源码**（`/d/work/base/gopath/pkg/mod/modernc.org/sqlite@v1.59.0/sqlite.go`）：`sort.Slice` 在**第 439 行**；比较器逐字＝`busy_timeout` 前缀排最前、其余按 `TrimSpace(ToLower(...))` 字典序；执行循环 `"pragma " + v` **按排序后**逐条 `exec`。⇒ **排序行为是真的，⛔ 按 DSN 顺序**。
- **`(1)`＋`(0)` 200 条**（树 `treeHead`；尺＝每条路 `sql.Open`＋读 `PRAGMA foreign_keys` 后立即 `Close` ⇒ **200 条全新物理连接**）：

| DSN 里的顺序 | 读 1 的条数 | 首／末 |
|---|---|---|
| `foreign_keys(1)` 然后 `foreign_keys(0)` | **200/200** | 1／1 |
| `foreign_keys(0)` 然后 `foreign_keys(1)` | **200/200** | 1／1 |

- **8 拼写 ×"追加覆盖"**（尺＝`buildDSN(shared+[S]+[OFF])` 与 `[OFF]+shared+[S]` 两种摆放都量；`shared`＝`dsnPragmas()`）：

| shared 里的写法 | 追加 OFF 读完 | OFF 摆前面读完 | r3 的形（删净＋恰好一条） |
|---|---|---|---|
| `foreign_keys(1)` | 0 | 0 | 0 |
| `foreign_keys = 1` | 0 | 0 | 0 |
| `foreign_keys(ON)` | **1** | **1** | 0 |
| `foreign_keys(on)` | **1** | **1** | 0 |
| `foreign_keys(true)` | **1** | **1** | 0 |
| `foreign_keys(yes)` | **1** | **1** | 0 |
| `foreign_keys(0)` | 0 | 0 | 0 |
| `foreign_keys(OFF)` | 0 | 0 | 0 |

⇒ **r3 的 §3 表逐格复现**：4 种拼写读完仍是 1，且**与 DSN 顺序无关**（两列同值）。"追加覆盖是拼写彩票"**成立**。

## §4 ★改后那一发：`rc=0` 与行数（我自造夹具）
夹具＝`seedLegacyV3(theFourLegacyWords)`（生产 `ddlV1/ddlV2/ddlV3`）＋`addLegacyToolCallRow`（子行 `task_id` 从 `task_log` `SELECT` 出来）。`treeInject`＝`dsnPragmas()` 只多一条 `"foreign_keys(1)"`（`diff -u` **一行**）。

| 读数 | `treeHead`（原样） | `treeInject`（注入＝将来那一改） |
|---|---|---|
| `Open()` err | `<nil>` | **`<nil>`** |
| `schema_version` | 4 | **4** |
| `task_log` / `tool_call` / 子行仍指原父 | 4 / 1 / 1 | **4 / 1 / 1** |
| `task_log_v4` 残渣（`sqlite_master`） | 0 | **0** |
| 迁移 DSN 实串 | `…wal_autocheckpoint(1000)&_pragma=foreign_keys(OFF)&_txlock=immediate` | **同左**（注入的 `(1)` **没进**迁移 DSN） |
| 业务 DSN 实串 | 无 `foreign_keys` | `…_pragma=foreign_keys(1)&_txlock=immediate` |

★**分流真的发生**：同一棵注入树上 `foreign_keys(1)` 进了业务 DSN、**没进**迁移 DSN。⚠ 该树整包 `rc=1` 是**设计使然**（`TestDsn196r3…` 在 `:46/:49` 红＝注入本身违边界 1），⛔ `Open()` 失败。
★**反证另一半**：`treeM2`（摘新跳＋注入）生产 `Open()` 逐字回到
`memory: migration 3->4 failed (…bak-3-4): apply: memory: refusing to rebuild task_log with PRAGMA foreign_keys=ON: … Nothing was modified.`
再在 `treeGuardOff`（M2 ＋ `requireForeignKeysOff` 首行 `return nil`）摘护栏 ⇒ `apply: constraint failed: FOREIGN KEY constraint failed (787)`。⇒ r3 §2 的 A/B 两行**逐字复现**。

## §5 ★「业务连接没被连累」——这一格最要紧，我独立核（含一发它没做的正控）

| 场景 | 尺 | 读数 |
|---|---|---|
| `treeInject` 上 `s.reader` / `s.writer` | 生产 `Open()` 之后的活池，读 `PRAGMA foreign_keys` | **1 / 1** |
| `treeHead`（原样）同上 | 同上 | **0 / 0** |
| `dsn()` 输出 | 对**改前实现的逐字副本**比串（3 路径 × writer/reader） | **逐字节相等**，且 `strings.Contains(…,"foreign_keys")`＝**false** |
| `dsn(...)+"&_pragma=foreign_keys(1)"`（调用方自设） | 真连接读回 | **1**（r3 用例覆盖，我复跑 PASS） |
| ★**连接级、⛔ 文件级**（r3 没单独量） | `treeInject` 跑完 `Open()`＋`Close()` 后，**同一文件**用**不带任何 pragma** 的 DSN 重开 | **0** |

⇒ ★末行是"⛔ 全局关"的硬证明：迁移连接说过的 `OFF` **没写进文件**，同一个库换条连接就读 0。⇒ 机理三层（连接级／用完即关／`dsn()` 沉默）**逐层成立**：`Open()` 里 `s.migrate()` 在 **261** 行、两个业务池在 **269/276** ⇒ 迁移连接先关（`396`/`413`）、业务池后开。

## §6 ★反向判据：我自己造，逐名逐行报红
尺＝仓外导出树整份改（`diff -u` 证单变量；工作树零突变），跑 `go test ./internal/memory/ -run 'TestDsn196r3|TestMigrationPragmas196r3|TestMigrationDSN196r3|TestTaskState196LegacyV3WithToolCallRowsStillOpens' -count=1 -v`。

| 变异体 | 唯一差别 | rc | 逐名／逐行 |
|---|---|---|---|
| **M1** | `migrationPragmas` 体改成 `return shared`（摘掉新加那一跳） | **1** | `TestMigrationPragmas196r3…`（`:105` 调用数 0≠1、`:110` 末条非 OFF，**5 子例全红**）＋`TestMigrationDSN196r3…`（`:183`，8 拼写里 **6 种**读 1：`(1)`／`= 1`／`(ON)`／`(on)`／`(true)`／`(yes)`） |
| **M2** | M1 ＋ `dsnPragmas` 注入 `foreign_keys(1)` | **1** | `taskstate_196r2_test.go:70` **逐字**回到 §4 那句可读拒（⛔ `(787)`）；★另红 `:46/:49/:105/:110/:169/:183`（件未列全，见 §2-d） |
| **M3** | `return append(shared, OFF)`（**只用"追加覆盖"**，保留 shared 那条）＋注入 | **1** | `:46`/`:49`（注入违边界 1）＋`:105`（**2～3 条** `foreign_keys`）＋`:183`（**`ON`/`on`/`true`/`yes` 四种读 1**） |
| **M3b**（★我加） | 同 M3，注入那条换成 **`foreign_keys(ON)`** | **1** | 迁移 DSN＝`…foreign_keys(ON)&_pragma=foreign_keys(OFF)` ⇒ 生产 `Open()` **直接失败**：`refusing to rebuild task_log with PRAGMA foreign_keys=ON … Nothing was modified.` |
| **M4**（★我加） | `open.go:365/402` 改回 `dsn(s.dbPath,true)`（接线回退） | **0** | **整包绿**（无牙，见 §2-c；有界） |
| 还原绿 | 工作树原样 | **0** | `ok … 22.545s` |
| 原样绿 | `treeBase`＝`9909404db` 同包 | **0** | `ok … 25.611s`（r3 的 21.807s 是负载差，**rc 才是判据**） |

★**`M3` 我复现了，而且这正是它推翻我的实证**：`append(shared, OFF)` 在 `(1)` 这类数值拼写下"侥幸"读到 0（M3 的 `Open()` 甚至 `err=<nil>`），换成 `(ON)`（M3b）**整个迁移就打不开**。⇒ 拼写一变行为就翻；"**删净＋恰好一条**"⛔ 风格，是**唯一与排序无关**的形状。

## §7 边界：逐项零 diff
尺＝`git diff --numstat 9909404db..HEAD`（★⛔ `grep '^[+-][^+-]'`，会漏"加号后为空"）。

| 项 | 读数 |
|---|---|
| r3 三笔名册 | `da38c2bf4`＝**3 枚**（`A 件`／`M open.go`／`A open_196r3_test.go`）；`c346a88de`＝**只 M 件**；`889373f93`＝**只 M 件** |
| 全区间（除台账） | `158/0 件`、`79/8 open.go`、`187/0 测试`——**恰 3 枚** |
| `internal/memory/schema.go` | **0/0** ⇒ `ddlV1`–`ddlV4`／`requireForeignKeysOff`／`pragmaForeignKeysQuery` **逐项零 diff**（该件最后一次改动是 r2 的 `0c08eca14`） |
| `internal/memory/taskstate.go`（12 枚名册） | **0/0**（`taskStateRoster` 字面 **12** 条：7＋5） |
| `docs/PLAN.md`／`docs/specs`／`cmd`／`scripts`／`.github`／`.scratch` | **空** |
| ★第 10 步硬门槛 | `git grep -n "foreign_key_check" -- "*.go"` ⇒ **rc=1 无命中**、⛔ 越界 |
| 台账 94 行 | 属编排者 `61220f85a`／`86ff73bdf`，**⛔ 属 r3** |

★四数对账：`79＋187＋147＝413` ✓（`git show --stat da38c2bf4`：`413 insertions(+), 8 deletions(-)`）；删除 `8` ✓。件现 **158 行**＝147＋11（`c346a88de`）＋2−2（`889373f93`）✓。
★SHA-256（**64 位十六进制**；`sha256sum` 与 `Get-FileHash -Algorithm SHA256` **两把尺同值**，且 `git show HEAD:<p>` 同值）：`open.go` `eb3cdbdae3ff59cddd1011101d3e38467edc2632b8007076a663c1354c367392`、`open_196r3_test.go` `20be5657e7052b5de9405487fd6eb1570e5224c79b3e24662ff25a2cacbb7cd4` ⇒ 与 r3 件 §7 **逐位相同**。

## §8 ★`schema.go:386`：**要改**——而且**不止**这一句（第二处更该改）
★我判：**要改**。⛔ 但 r3 只找到一处，**同类另一处更该改**：
- `schema.go:386`（注释）：`// gets that for free because dsn() never sets the pragma (open.go), so what the` —— 现在是**反的**：迁移链"免费拿到"第 1 步的原因**不再是 `dsn()` 沉默**，而是 `migrationDSN` 自己**声明**了 `OFF`；`dsn()` 明天真带 FK，迁移也不受影响。
- `schema.go:416-421`（**运行期报错串**，维护者真会看到）：`… this package relies on dsn() never setting the pragma, so something added it - drop the foreign_keys pragma from the migration DSN and retry.` —— 两句都过时：①"relies on dsn() silence"已不成立②"drop the foreign_keys pragma from the migration DSN"**不再是可做的动作**（DSN 由 `buildDSN` 从内部两张表拼，用户改不了）。
- ★**这是真风险而⛔ 洁癖**：留着它，下一腿读到"机制是 `dsn()` 沉默"，最自然的"修复"就是**删掉 `migrationDSN`**——而 **M4 实测整包仍绿**（§6）。⇒ **过时注释正好在招回那个没有牙的回归**。

**具名措辞（由下一腿做，⛔ 本腿改 `.go`）**
1. `schema.go:385-388` 注释末两句替换为：
   `// disables enforcement for the duration of a rebuild (its step 1); the migration chain`
   `// takes that step for itself in migrationDSN (open.go), which drops every foreign_keys`
   `// entry from the shared list and appends exactly one OFF. This guard is therefore a`
   `// backstop rather than the load-bearing precondition: it fires if that declaration is`
   `// ever lost.`
2. `schema.go:419-421` 报错串三段替换为：
   `"SQLite's 12-step ALTER procedure disables enforcement for a rebuild (its step 1), and "`
   `"the migration connection is supposed to declare that itself (migrationDSN in open.go): "`
   `"something removed that declaration. Nothing was modified."`

（★这不改 `C1–C32`／`D1–D47`，但要随票 196 的收口笔一起做，⛔ 顺手混进别的腿。）

## §9 门禁 rc ＋ 提交号
**五门／四数**（工作树＝`889373f93`；落笔前 `git status --porcelain -- internal/memory docs/evidence/s1` **空**）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| 建 | `go build ./...` | **0** | 输出 0 行 |
| 静态 | `go vet ./...`（整模块） | **0** | 输出 0 行 |
| 包测 | `go test -count=1 ./internal/memory/` | **0** | `ok … 22.545s` |
| D22 | `sh scripts/d22scan.sh` | **0** | `clean`；ban #8 `internal/` **633** 件（r2 时 632，＋1＝新用例）／`cmd/` 148 |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `tracked paths=17647  over-budget=61  covered=61  not in roster=0`、`VERDICT GREEN` ⇒ **r3 报的 17647 对** |

**四数分开**（尺＝`go test -v`；★`RUN` 含 subtest，⛔ 与 PASS 同分母）：`rc=0`、`RUN=101`（**顶层 61**）、`PASS_top=59`、`FAIL_top=0`、`SKIP_top=2`。⇒ 相对 r2 时代（`93 / 顶层 58 / 56 / 0 / 2`）：顶层 **+3**（三枚新用例）、`RUN` **+8**（＋5 个 subtest）＝**对得上**。

**格式甲乙两尺**
- **甲（判据）**：`git archive HEAD` 导出仓外 `treeFmt`，`gofumpt -l internal cmd tools scripts frontend` ⇒ **0 行**；单跑 `open.go`／`open_196r3_test.go`／`schema.go` ⇒ **0 行**。
- **乙（只诊断，工作树）**：**3 行**，按 `[\\/]` 归属＝`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go`（**与 r2 判过的同一 3 枚，⛔ 本腿文件**）；`git ls-files --eol` 三件 r3 文件全 `i/lf w/lf`。

**提交号**：本件 `________`（本笔＝裁决正文；回填笔只改本行，⛔ amend）。提交前 `git diff --cached --name-status` 只有 `A docs/evidence/s1/196-v3-verdict.md` 一枚。
**写面**：本件 ＋ 仓外 `D:\tmp\wisp\20261023\196v3\**`；⛔ 改 `internal/**`／`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**`；⛔ 落任何 `.go`。
