# 票 366 / 腿 `366-r2`（写腿）— 补牙：坏尺用例 ＋ 8 处措辞 ＋ 还原核重量（＋两条小上报）

- 腿＝`366-r2`（写腿）· 派单主体＝`docs/evidence/s1/366-v1-verdict.md`（114 行）· 被补对象＝`366-r1`（`d4ce46f97`／`2abf150b2`／`c2e47f70c`，件 143 行）· 起始 HEAD＝`b0bf991b1`（含终裁 `a803c6c16`）· 2026-10-23 19:1x–2x +08
- 写面＝`internal/memory/**`（`open.go`／`schema.go` **只改注释措辞**；`ForeignKeyAudit` 的 `slices.Clone` 是裁决特批的一处行为）＋测试 ＋本件；写面外 **0 字节**。
- 仓外树：`D:\tmp\wisp\20261023\366r2\{headtree,posttree,restoretree,mut,logs,probes}\**`；**五个变异全用 `go test -overlay`**（`posttree` 自身从未被写）⇒ 工作树无突变窗口。
- **尺的射程（先写清）**：夹具＝生产链建的真 v3 库（`Open(dir, withSchemaTarget(3))`）；孤儿＝生产调用（`PurgeTaskLogs`／`DeleteTaskLog`／raw INSERT／`RunRetentionOnce`）；迁移＝生产 `Open(dir)`；`fkc`＝迁移后另开业务连接数 `PRAGMA foreign_key_check` 行数；计数＝`Store.ForeignKeyAudit()`；**hash＝SHA-256、64 位十六进制**（`Get-FileHash -Algorithm SHA256` 与 `sha256sum` 两把同算法尺；⛔ `git hash-object`＝SHA-1/40 位，不混比）。

## §1 先复跑 `366-v1` 的核心读数（⛔ 抄；我在 `headtree`＝HEAD 导出树上自己跑，`logs/10`）

| 夹具 | 前 fkc | `Open()` | `Checked` | v | `OrphanRows` | `child_tables` | 后 fkc | task_log/tool_call | 再开一次 |
|---|---|---|---|---|---|---|---|---|---|
| `clean-v3`（无孤儿正控） | 0 | `<nil>` | true | 4 | 0 | `[]` | 0 | 2/3 | `<nil>`、`Checked=false` |
| `purge-tasklogs` | 3 | `<nil>` | true | 4 | **3** | `[tool_call]` | 3 | 0/3 | `<nil>`、fkc 仍 3 |
| `delete-tasklog` | 2 | `<nil>` | true | 4 | **2** | `[tool_call]` | 2 | 1/3 | `<nil>`、fkc 仍 2 |
| `manual-insert` | 1 | `<nil>` | true | 4 | **1** | `[tool_call]` | 1 | 2/4 | `<nil>`、fkc 仍 1 |
| `retention` | 1 | `<nil>` | true | 4 | **1** | `[tool_call]` | 1 | 2/4 | `<nil>`、fkc 仍 1 |

⇒ 与 `366-v1` ③④⑤ 同值（3/2/1/1，`clean` 0）；**计数与独立 fkc 逐条相等**；孤儿一行未删；第二次 `Checked=false`＝`version == s.target` 早退（如实）。
**坏尺那一支（HEAD 上）**：我在真库上种一张"FK 指到不存在列"的子表（合法 DDL；SQLite 对全库答 `foreign key mismatch`）⇒ `Open()` 仍 `<nil>`、`Checked=false`、孤儿 3 仍在、日志有 WARN 且**没有** `found no dangling references`。**但终裁对：HEAD 上没有用例钉它**（§2 补）。

## §2 缺 1：坏尺 ⇒ 照常开库，**有新用例**，且与"尺好、有孤儿"**区分得开**

新用例：`TestMigration366AuditNeverTurnsABrokenRulerIntoALock`（`migration_fkcheck_366_test.go:494`，4 子例）。接缝两发（无需新产码接缝）：①`t366FailingTx`（终裁具名的 `txExec` 接口，两个方法各返错）；②真库级坏尺 `t366PlantUnreadableRuler`（让 SQLite 自己答错，`Open()` 走完整生产路径）。

| 子例 | 我自己跑出来的读数 |
|---|---|
| `failing transaction` | `Checked=false`、`SchemaVersion=4`、`OrphanRows=0`、`ChildTables=[]`；日志 `level=WARN`＋`could not be read`＋`opening anyway`＋`schema_version=4`；无 `level=ERROR` |
| `broken ruler with orphans` | **`Open()` `err=<nil>`**、`Checked=false`、v=4、`OrphanRows=0`、`ChildTables=[]`；**孤儿仍在 3**（`tool_call WHERE task_id NOT IN task_log`＝3；`tool_call` 行数＝3，一行未删）；`schema_version=4`；日志同上，**无** `found no dangling references` |
| `broken ruler on a clean database` | `Open()` `err=<nil>`、`Checked=false`；**干净库⛔ 打** `found no dangling references`（不"假装检查过"），只打 WARN |
| `good ruler control on the same fixture` | 同一 `purge` 夹具、尺好 ⇒ `Checked=true`、`OrphanRows=3`、`ChildTables=[tool_call]` ⇒ **与坏尺读数不同**（子例内显式比较：两读数同则红） |

- 为什么"区分得开"必须写死：只测"尺坏 ⇒ `Checked=false`"，一个恒 `false` 的实现也会绿；对照发把 `{true,3,[tool_call]}` 也要到。
- 牙的两向都在：①"假装检查过"（`Checked=true`）⇒ M3 红；②"坏尺变锁"（`!Checked` 就拒开）⇒ M5 红（§6）。
- 真库级坏尺的 `Open` 日志逐字（`logs/10`）：`level=WARN msg="migration foreign-key baseline could not be read; opening anyway" component=memory schema_version=4 err="memory: PRAGMA foreign_key_check: SQL logic error: foreign key mismatch - \"t366_child\" referencing \"t366_parent\" (1)"`。

## §3 缺 2：8 处措辞逐字（改前 → 改后；终裁给的措辞**照抄**，行号现跑）

1. `schema.go:204-206`（原 `:197-198`）：`rebuild is SQLite's documented 12-step ALTER procedure, and this is its / step 10.` → `... and this is an / extra check at the point where that procedure may check integrity / (unverified here).`
2. `schema.go:312-313`（原 `:303`）：`Step 10 of that procedure - PRAGMA foreign_key_check - cannot live in this` → `The integrity check (possibly step 10 of that procedure; unverified here) - / PRAGMA foreign_key_check - cannot live in this`
3. `schema.go:520-523`（原 `:509-510`）：`... referential-integrity check: the step 10 / of the documented 12-step ALTER TABLE procedure that the v4 rebuild follows / (ddlV4), and the one step ...` → `... check: an extra / guard for the v4 rebuild (possibly the same place as step 10 of the documented / procedure; unverified here), and the one check this chain did not run until / ticket 366.`
4. `schema.go:529`（原 `:518`）：`ForeignKeyAudit is the migration chain's step-10 record:` → `ForeignKeyAudit is the migration chain's post-rebuild integrity record:`
5. `open.go:332-334`（原 `:331-332`）：`... before the commit (the 12-step / ALTER procedure's step 10; ticket 366).` → `... commit (extra / insurance that may coincide with step 10 of the 12-step ALTER procedure; / unverified here; ticket 366).`
6. `open.go:522-524`（原 `:513-514`）：`SQLite's 12-step procedure checks referential integrity after the rebuild / and before the commit (its step 10).` → `An extra integrity check after the rebuild and before the commit (where the / 12-step procedure may put one; that step is conditional on enforcement / having been ON, and is unverified here).`
7. 用例 `:3`：`Ticket 366 leg r1: the 3->4 rebuild now runs SQLite's 12-step procedure step 10 / (PRAGMA foreign_key_check) ...` → `Ticket 366: the 3->4 rebuild runs an extra referential-integrity check (PRAGMA / foreign_key_check) ... It is extra insurance that may coincide with step 10 of / SQLite's 12-step ALTER procedure; that step is unverified here and is conditional / on enforcement having been ON, which it never is in this package (migrationDSN turns it OFF).`
8. 用例 `:428`（原 `:424`）：`nothing for step 10 to check` → `nothing for the extra check to look at`

- **en 总句**（代码注释，落 `schema.go:194-201`，逐字照抄）：`An extra referential-integrity guard around the v4 rebuild. SQLite's 12-step ALTER TABLE procedure may run the same check at its step 10, but this repository has not verified that text, and that step is conditional on enforcement having been ON - which it never is here (migrationDSN turns it OFF). Extra insurance that may coincide with step 10; not a step this code claims to follow.`
- **中文句**（本件口径）：`额外保险（可能与官方 12 步的第 10 步重合；官方原文本机取不回、未核；且官方那一步以"原本开着 FK"为前提，本仓迁移连接一直是 OFF）`。
- ⚠ **三处改不动的旧措辞（终裁点名，⛔ 硬改；如实记）**：`schema.go:334`（原 `:324`，**在 `ddlV4` 字符串里**，动它就破硬边界 2）· `open.go:185-186`（原 `:184-185`）· `schema.go:416-417`（原 `:406-407`）与 `schema.go:452`（原 `:442`）—— 后三处属 `196-r3`。
- ★**牙在哪儿（如实）**：以上 8 处**全是注释**，运行时行为 0 变化 ⇒ **没有、也不该假称有**"今天响"的用例；它们的判据是**终裁逐字对照**＋`gofmt/gofumpt`/`go vet`。本腿真正带牙的改动只有两处行为面：坏尺分支（§2，M3）与 `slices.Clone`（§5，M4）。

## §4 缺 3：还原核 4 枚 SHA-256（**64 位全文**、算法写明、量在最后一发变异之后）

| 文件 | SHA-256（64 位十六进制） | 字节 |
|---|---|---|
| `internal/memory/open.go` | `08BDFA15D3A32FF712AB38B288B1410375FC9412F00DCE2DC7082793641D4AD7` | 27170 |
| `internal/memory/schema.go` | `DF2C999885C3A91018F82A0CD9AF8DD79E238A1C05E4C6C391DEC50447DC3CA4` | 28563 |
| `internal/memory/taskstate_196r2_test.go` | `A34A2591026E911AF5041DE22E91BCC4CCA98CACB3DBA1482F5F3C258D4D6303` | 8831 |
| `internal/memory/migration_fkcheck_366_test.go` | `4A9E75DD00587AEDD44F7DC2258AA2F181ED46020C4C7D4ECC8C8E613006D86A` | 26657 |

- 算法＝**SHA-256**（`Get-FileHash -Algorithm SHA256` 与 `sha256sum` 交叉各算一次，同值）；**量在最后一发变异（M5）之后**（M5 后重量：12/12 仍同值，四枚值不变）；**工作树＝`posttree`＝`restoretree` 三处逐字相同**；提交后再与 `git show 43ca90fc8:<path>` 的 blob 比（§10，三枚逐枚相同）。
- **为什么 r1 的 `3d9f2655…` 对不上**（⛔ 猜，可复现）：它逐字等于 r1 那份 `D:\tmp\wisp\20261023\366r1\posttree\internal\memory\open.go` 的 SHA-256 全文＝`3D9F265541DFD34A046A6D92F15F4B04E79321DCB6AC88222715ECC8D3D70F56`（len 26745）；入库 blob 是 `d13000b1…`（len 26781）。差 **36 字节＝3 行注释重排**（`ForeignKeyAudit()` 的 doc：posttree 写 `(see ForeignKeyAudit's doc comment) and never / deletes a row (tool_call is the forensic surface)`；入库写 `(the ForeignKeyAudit type's doc comment says why, with the measurement) and never / deletes a row ...`）⇒ **它量在"最后一次注释改写之前的那份 posttree"上，之后没把 posttree 与交付字节重新对齐**；既不是"量在变异之间"（M1＝`dde55d69…`、M2＝`977b6bdd…`），也不是别的文件的 hash。佐证：r1 同目录另三枚与 HEAD blob **逐字相同**，只有 `open.go` 差。⇒ 这是"读数不可复现"，不是"读数不符"——我能复现 `3d9f2655…`，复现出来的是那份旧副本。

## §5 两条小上报

1. **子例数（自己重量，⛔ 抄 r1）**：排除本腿新增两枚用例的口径 —— `M1`＝**4 顶层／8 子例**（r1 报 10 子例，错）；`M2`＝**3 顶层／5 子例**（r1 报 6 子例，错；它把一枚顶层算进子例）。含本腿新增用例的本次总数：`M1` 6 顶层／11 子例、`M2` 5 顶层／6 子例（`logs/40-mutation-m1|m2.txt`）。
2. **`slices.Clone`（裁决：做）**：`open.go:350-354`——
   `func (s *Store) ForeignKeyAudit() ForeignKeyAudit {` / `audit := s.fkAudit` / `audit.ChildTables = slices.Clone(audit.ChildTables)` / `return audit` / `}`；牙＝`TestMigration366AuditChildTablesIsACopy`（`:621`）：调用方写 `first.ChildTables[0]` 后，第二次读数必须仍是 `tool_call`；变异 M4 摘掉 clone ⇒ 红行 `:635`。

## §6 反向判据：三读数 ＋ 还原核

| 变异（`-overlay`，仓外） | 内容 | 红读数（`logs/40-mutation-*`） |
|---|---|---|
| `M1` | 摘掉取审计那一跳（`audit = auditForeignKeys(...)`） | rc=1；**6 顶层／11 子例**（排除新增两枚＝**4／8**）；红行 `:211,:236,:282×4,:337,:384,:387,:407×4`（＋新用例 `:550×3,:584,:607,:630`） |
| `M2` | 把记录改成门槛（`OrphanRows>0 ⇒ Rollback+拒开`，硬边界 #1） | rc=1；**5 顶层／6 子例**（排除新增＝**3／5**）；`Open refused a database with 3/2/1/1 orphan row(s)`；红行 `:275×4,:333,:396`（＋新 `:600,:625`） |
| `M3`（本腿补牙专用） | 坏尺分支**谎报** `Checked=true`（"假装检查过"） | rc=1；1 顶层／4 子例；红行 `:499,:530,:554,:575,:580,:611` |
| `M4`（`slices.Clone` 专用） | 去掉 clone（算出即弃） | rc=1；1 顶层；红行 `:635` |
| `M5`（坏尺变锁专用） | `!audit.Checked ⇒ Rollback+拒开`（把 ⑧ 的坏尺变成永久锁） | rc=1；1 顶层／3 子例（坏尺两发＋对照发）；红行 `:524,:570,:593` |

- ② **还原绿**：`restoretree`（新导 HEAD＋三枚交付字节）`go test -count=1 -v ./internal/memory/` ⇒ **rc=0**、120 RUN／66 PASS／0 FAIL／2 SKIP（`logs/71`）。
- ③ **原样绿**：`posttree`（同一交付字节，变异前）⇒ **rc=0**、120／66／0／2（`logs/31`）。
- 还原核时机＝**最后一发变异（M5）之后**（§4 的 12 次比较，M5 后重量）；⛔ 在变异之间量。
- `M1` 为什么是"最小摘除"而不是整条回滚：新用例在 HEAD 原码上编不过（`ForeignKeyAudit`／`foreignKeyCheck`／`pragmaForeignKeyCheckQuery` 都还不存在）⇒ 只会得到"编译红"。

## §7 四数（`go test -count=1 ./internal/memory/`，改前／改后**分开**、⛔ 相减；`-v` 下计）

| 读数 | 改前（HEAD 工作树，`logs/21`） | 改后（交付字节，`logs/70`） |
|---|---|---|
| rc | **0** | **0** |
| `ok` 行 | `24.605s` | `27.954s` |
| `=== RUN`（含子例） | **114** | **120** |
| `--- PASS`（**顶层**） | **64** | **66** |
| `--- FAIL` | **0** | **0** |
| `--- SKIP`（顶层） | **2** | **2** |

新用例 1＋4 子例＋clone 用例 1＝**+6 RUN**；`+2` 顶层 PASS。★尺：`--- PASS` 只数**顶层**（子例缩进，不计）；`=== RUN` 含子例 ⇒ 两列不同尺。
**提交后复跑**（`ec042da5f` 的字节，产码未改）：rc=**0**、120 RUN／66 PASS／0 FAIL／2 SKIP、`ok … 27.846s`（`logs/90`）⇒ 与上表改后列同值。

## §8 门禁 rc（逐门分开、⛔ 走管道取 rc）＋格式两尺

| 门 | 命令 | rc |
|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **0**（0 字节输出） |
| G2 | `GOFLAGS= go vet ./internal/memory/` | **0**（0 字节输出） |
| G3 | `GOFLAGS= go test -count=1 ./internal/memory/` | **0**（§7） |
| G4 | `sh scripts/d22scan.sh` | **0**（`clean - no D22 ban violations`；ban #8 internal/ 639 枚） |
| G5 | `sh scripts/check-path-length-budget.sh` | **0**（`VERDICT GREEN`；roster 61） |

- **格式甲（入库 blob＝判据）**：`git show <commit>:<p>` 导出到仓外再 `gofumpt -l` ⇒ **见 §10**。**乙（工作树＝只作诊断）**：`/d/work/base/gopath/bin/gofumpt.exe -l internal/memory` ⇒ **0 行**（`[\\/]` 过滤后仍 0）。
- **CRLF 尺（⛔ 用 `gofumpt -l`）**：四枚 `tr -cd "\r" | wc -c` ⇒ **全 0**（LF）。
- **硬边界自证**：`ddlV1–ddlV4` 两树各算 SHA-256 逐字相同（`fda736e1…`／`3e9ece70…`／`d46ee18e…`／`c17241cd…`，len 4463/903/1286/1357，`logs/50|51`）；12 名册两树同为 `running,…,Stuck,Error`；新增行里大小写敏感 `DELETE` **0 命中**；`196-r3` 的 `migrationDSN`／`migrationPragmas`／`dsnPragmas` 与 `196-r2` 的 `requireForeignKeysOff` **零 hunk**。

## §9 纪律格

- **只 commit、⛔ push**；`add`/`commit` 两处**显式 pathspec**（逐个文件名字面点名、⛔ glob）；⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ 碰票 196 原句、⛔ `DELETE`。
- **没写的**：三处改不动的旧措辞（§3 末）**一字未动**；`ddlV1`–`ddlV4` 的 SQL（含 `schema.go:334` 那句中文）**一字未动**；任何"删孤儿"动作**没有**（孤儿夹具迁移后行数逐条不变，§2）。

## §10 提交

- 本件提交＝**`43ca90fc8`**（`4 files changed, 371 insertions(+), 21 deletions(-)`：3 枚 `.go` ＋本件）；`add`/`commit` 各自只点四个文件名（字面点名、⛔ glob）；提交前 `git diff --cached --name-status` 只见 **`A` 本件＋`M` 三枚**（无 `D`）。本件其后另有**只改本件**的补录（见 `git log --oneline -- docs/evidence/s1/366-r2-impl.md`：`5bc0ff51c`＝§10、`ec042da5f`＝M5 等），产码字节 0 变化。
- **格式甲复核（入库 blob＝判据）**：`git archive 43ca90fc8` 导出三枚 `.go` 到仓外再 `gofumpt -l` ⇒ **0 行**（rc=0）；三枚 blob 的 SHA-256 逐枚等于 §4 表（`08BDFA15…`／`DF2C99…`／`4A9E75DD…`）⇒ **还原核落在入库字节上**（`taskstate_196r2_test.go` 未改，其 blob 仍＝`A34A2591…`，与 §4 同）。
- ⛔ push；台账／票面／`196` 原句 0 字节改动。
