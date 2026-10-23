# 票 366 / 终裁腿 `366-v2`（非实现者）— `AC#2` **能翻**：三项缺全补上了（每项我自己另起一把尺复量）

- 腿＝`366-v2`（**非实现者裁腿**）· 裁的对象＝`366-r2`（`43ca90fc8` ＋ 3 笔只改件补录 `5bc0ff51c`／`ec042da5f`／`9e0127c2f`；件 **122 行**）· 裁时 HEAD＝`4bb9f8219`（编排者 A1208）；`git diff --stat 43ca90fc8 HEAD -- internal/memory` **空** ⇒ 我量的是交付字节。
- 被补对象＝`366-v1`（114 行）判的**三项缺**；形仍是〔甲〕基线化（只记录、⛔ 拒开库），本腿**不重开形**。
- 写面＝本件 ＋ 仓外 `D:\tmp\wisp\20261023\366v2\**`（`headtree`＝HEAD 导出树 · `posttree`＝交付字节树 · `restoretree`＝**最后一发变异之后**的新导出 · `pretree`＝`6d47e1ff4`（366 改前）· `pretree-clean`／`r1tree`＝四数对照 · `mut/` `probes/` `logs/`）。**仓内 `.go` 0 字节**；五发变异**全走 `go test -overlay`**（仓外），工作树无突变窗口（⑩）。
- **尺的射程（先写清）**：夹具＝生产链真 v3（`Open(dir, withSchemaTarget(3))`）＋生产孤儿路径 `PurgeTaskLogs`；坏尺①＝`v2BadTx`（实现 `txExec` 两方法各返错）直调 `auditForeignKeys`；坏尺②＝**真库级**（子表 FK 指不存在列 ⇒ SQLite 对全库答 `foreign key mismatch`）⇒ `Open()` 走**完整生产路径**；孤儿读数＝**另开一条业务连接**数 `tool_call WHERE task_id NOT IN (SELECT id FROM task_log)`；hash＝**SHA-256／64 位**（`sha256sum`，⛔ `git hash-object`＝SHA-1/40 位）。

## ① 票面前提新鲜度（对 HEAD blob 复量）

| 票面的话 | HEAD 复量（我自己跑） | 判 |
|---|---|---|
| 重建**自称**走官方 12 步 | `schema.go:308` 逐字仍在（普查件引 `:283-284`、v1 记 `:299-300` ⇒ 又被 r2 推到 `:308`） | ✅ 事实在、行号漂 |
| 第 1 步（关 FK）`196-r3` 已补 | `migrationDSN→buildDSN(...,true,migrationPragmas(...))`；`migrationPragmas` 末尾 `append(out, pragmaForeignKeysOff)`，常量 `"foreign_keys(OFF)"`（`open.go:47`）；业务 `dsn()` **0** 处 `foreign_keys` | ✅ 仍在 |
| ★"第 10 步**从来没跑过**" | **过期**：`foreign_key_check` 命中 schema 10／open 5／用例 19＝**34 行** | ⚠ 已被本票实现改写（正常） |
| 孤儿**三条**路径 | **四条**：`dao_tasklog.go:118`（删父不删子）／`:134`（清父表）／`:130` "own domain" 注释未漂 · raw INSERT（绕 DAO）· `retention.go:197` 只删**非 NULL 时间戳**行 ⇒ NULL 子行留下 | ⚠ 票面"三条"少计一条（同 v1） |
| `196-v2` 实测"迁移前后都＝1" | "1"是**夹具行数**不是常量：我量 purge＝3、clean＝0 | ✅ 方向真、数字非常量 |

## ② ★★AC#2 判语：**能翻**

- `366-v1` 判的三项缺，**逐项都在交付字节上补上了**，且每项我另起一把尺复量过：缺 1＝③（新用例＋两发坏尺＋对照半＋`M3`/`M5` 双红）· 缺 2＝④（8 处逐字改＋总句逐字＋三处改不动真没动）· 缺 3＝⑤（四枚 64 位**逐字相等**＋归因**逐字复现**）。
- 行为面在 r2 里**只多一处**：`ForeignKeyAudit()` 的 `slices.Clone`（`open.go:350-354`），牙＝`TestMigration366AuditChildTablesIsACopy`（`:621`），我造的 M4 红在 `:635`（⑥）。其余产码改动**全是注释**（`git show 43ca90fc8` 的非注释增删行只有 `"slices"` import 与 clone 四行；schema.go 增删**零**非注释行）。
- 边界、门禁、格式、纪律**全绿**（⑦⑧⑨⑩）；⛔ 未翻任何框、⛔ 改票名/`Status:`/`-done`、⛔ push。

## ③ 缺 1（**核心**）：坏尺 ⇒ 照常开库 —— 我自己造两发坏尺，五枚读数逐条复现

**我的五枚读数**（`logs/10-v2-probe.txt`，探针＝`probes/zz366v2_probe_test.go`，只在导出树里编译）：

| 发 | 读数（我自己跑） |
|---|---|
| 坏尺① `v2BadTx` 直调 | `{checked false, version 4, orphans 0, tables []}`；日志 `level=WARN`＋`could not be read`＋`opening anyway`＋`schema_version=4`、**无** `found no dangling references`、**无** `level=ERROR` |
| 坏尺② 真库级＋孤儿（`PurgeTaskLogs`＝3 孤儿） | 迁移前 `orphans=3／tool_call=3／task_log=0`；坏尺实证错＝`SQL logic error: foreign key mismatch - "v2_child" referencing "v2_parent" (1)`；**`Open()` `err=<nil>`**、`checked=false`、`version=4`、`orphans=0`、`tables=[]`；迁移后 **`orphans_still=3`／`tool_call=3`／`schema_version=4`**（一行未删）；日志同上 + **无** clean 行 |
| 坏尺② 干净库 | `Open()` `err=<nil>`、`checked=false`、**无** `found no dangling references`（不假绿）、有 WARN |
| 对照半（**同一 purge 夹具、尺好**） | `checked=true`、`version=4`、`orphans=3`、`tables=[tool_call]` |
| 两读数显式比较 | **DISTINCT=true**（checked `false/true`、orphans `0/3`） |
| 真库级坏尺 WARN 逐字 | `level=WARN msg="migration foreign-key baseline could not be read; opening anyway" component=memory schema_version=4 err="memory: PRAGMA foreign_key_check: SQL logic error: foreign key mismatch - \"v2_child\" referencing \"v2_parent\" (1)"` |

- ★**"它真的区分得开"——判：做到了。** 理由不是"`Checked=false` 出现了"，而是：①**同一夹具**在尺好时给出**必须不同**的 `{true,3,[tool_call]}`（交付用例 `:606` 断言 good 半、`:610` 断言两读数不同）；②**两向故障都咬得住**：坏尺**谎报** `Checked=true` ⇒ 我造的 `M3` 红（1 顶层/4 子例，`:499,:530,:554,:575,:580,:611`）；坏尺**变锁**（`!Checked` 拒开）⇒ 我造的 `M5` 红（1 顶层/3 子例，`:524,:570,:593`）——`M5` 咬的正是**硬边界 #1**（库必须能开）；③把尺换成**恒 `false`** 的常量实现也过不了（good 半 `:607` 要 `{true,3}`）。⇒ `366-v1` 说的"信号饱和"（4 红 1 绿、分不出拒开与开库）**已解除**。
- 交付用例本身：`TestMigration366AuditNeverTurnsABrokenRulerIntoALock`（`:494`，**4 子例**）＋ `t366FailingTx`（`:436`，终裁具名的 `txExec` 接缝）＋ `t366PlantUnreadableRuler`（`:451`，真库级且**先自证尺坏**再迁移）⇒ 与件 §2 报的四子例逐条对得上（我跑绿、变异红）。

## ④ 缺 2：8 处措辞**逐处核**（改前→改后，r1 blob→HEAD blob）

| # | 处所（改后行） | 改前（r1）→改后 |
|---|---|---|
| 1 | `schema.go:204-206`（原 `:197-198`） | `…and this is its / step 10.` → `…and this is an / extra check at the point where that procedure may check integrity / (unverified here).` |
| 2 | `schema.go:312-313`（原 `:303`） | `Step 10 of that procedure - PRAGMA foreign_key_check - cannot live in this` → `The integrity check (possibly step 10 of that procedure; unverified here) - / PRAGMA foreign_key_check - cannot live in this` |
| 3 | `schema.go:520-523`（原 `:509-510`） | `the step 10 / of the documented 12-step…that the v4 rebuild follows / (ddlV4), and the one step…` → `an extra / guard for the v4 rebuild (possibly the same place as step 10 of the documented / procedure; unverified here), and the one check this chain did not run until / ticket 366.` |
| 4 | `schema.go:529`（原 `:518`） | `…migration chain's step-10 record:` → `…migration chain's post-rebuild integrity record:` |
| 5 | `open.go:332-334`（原 `:331-332`） | `(the 12-step / ALTER procedure's step 10; ticket 366)` → `(extra / insurance that may coincide with step 10 of the 12-step ALTER procedure; / unverified here; ticket 366)` |
| 6 | `open.go:522-524`（原 `:513-514`） | `SQLite's 12-step procedure checks referential integrity after the rebuild / and before the commit (its step 10).` → `An extra integrity check after the rebuild and before the commit (where the / 12-step procedure may put one; that step is conditional on enforcement / having been ON, and is unverified here).` |
| 7 | 用例 `:3-7`（原 `:3-4`） | `…now runs SQLite's 12-step procedure step 10 / (PRAGMA foreign_key_check)…` → `…runs an extra referential-integrity check (PRAGMA / foreign_key_check)…may coincide with step 10…unverified here…conditional on enforcement having been ON…` |
| 8 | 用例 `:428`（原 `:424`） | `nothing for step 10 to check` → `nothing for the extra check to look at` |

- **en 总句**（`schema.go:194-201`）：把注释头剥掉、空白归一后与我按终裁文字写的 want **MATCH**（程序化比较，非目测）。
- 交付件里**再没有**把"第 10 步"当已核事实的句子：`grep` 只剩 `:195`（总句里的 `may run…at its step 10, but…not verified` 对冲句）。
- ★**三处改不动的旧措辞，真的一字没动**（`pretree` vs `posttree` 逐字比）：`schema.go:334`（在 `ddlV4` 字符串里，原 `:302`）——★**`ddlV4` 运行值 `c17241cd…`、len 1357／5 句，两树逐字相同**（`logs/51`）⇒ 字符串内那句中文没动；`open.go:185-186`（原 `:180-181`）与 `schema.go:416-417`（原 `:385`）／`:452`（原 `:420`）——三条"disables enforcement for…(its step 1)"句子两树比 **IDENTICAL（3 行）**（`logs/61`）。

## ⑤ 缺 3：还原核 4 枚 SHA-256（**我自己算，64 位全文**）＋ 归因

| 文件 | 我算的 SHA-256（64 位） | 字节 | 与 r2 报的 |
|---|---|---|---|
| `internal/memory/open.go` | `08BDFA15D3A32FF712AB38B288B1410375FC9412F00DCE2DC7082793641D4AD7` | 27170 | **逐字相同** |
| `internal/memory/schema.go` | `DF2C999885C3A91018F82A0CD9AF8DD79E238A1C05E4C6C391DEC50447DC3CA4` | 28563 | **逐字相同** |
| `internal/memory/taskstate_196r2_test.go` | `A34A2591026E911AF5041DE22E91BCC4CCA98CACB3DBA1482F5F3C258D4D6303` | 8831 | **逐字相同** |
| `internal/memory/migration_fkcheck_366_test.go` | `4A9E75DD00587AEDD44F7DC2258AA2F181ED46020C4C7D4ECC8C8E613006D86A` | 26657 | **逐字相同** |

- **四源同值 16/16**：工作树＝`posttree`＝`restoretree`＝`git cat-file blob HEAD:` 逐枚相等（`logs/95`），且 `git status --porcelain -- internal/memory` **0 行**。
- ★**归因核（逐条都是可复现的值，不是说法）**：`366r1/posttree/internal/memory/open.go` 的 SHA-256 **＝** `3D9F265541DFD34A046A6D92F15F4B04E79321DCB6AC88222715ECC8D3D70F56`、**len 26745**；r1 入库 blob（`567e4347f`，`c2e47f70c:internal/memory/open.go`）**＝** `d13000B18C2BD5C38D0FEFA8AFE33624AF9DA09F5A985683E61506D0921CE6DF`、**len 26781**；**26781−26745＝36**，而 `diff -u` 恰是 **3 行删／3 行增**、位置＝`ForeignKeyAudit()` 的 doc（旧 `(see ForeignKeyAudit's doc comment) and never…` → 新 `(the ForeignKeyAudit type's doc comment says why, with the measurement) and never…`）；三行字节 **187／223 ⇒ 差 36**（与文件差**逐字节对上**）。
- ⇒ **r2 的归因对**：r1 量在"最后一次注释改写之前的那份 posttree"，之后没与交付字节对齐；**不是**变异产物（v1 的 M1 `dde55d69…`／len 26728、M2 `977b6bdd…`／len 26972，我现算复核，两值都对；**不是**别的文件：r1 同目录另三枚＝r1 自己的入库 blob，只有 `open.go` 差）。
- 时序：四枚量在 **M5 之后**（`restoretree` 是最后一发变异之后才导出，`logs/71`）；**提交后复跑 G3**：`rc=0`、120/66/0/2（`ok 24.414s`，`logs/99`）。

## ⑥ 两条小上报（我自己数／自己造牙）

- **子例数**（尺：`-run 'TestMigration366'` 下 `--- FAIL` 行，顶层 vs 缩进子例，`logs/40-mutation-*`）：`M1`＝**6 顶层／11 子例**，**排除 r2 新增两枚＝4 顶层／8 子例**（r1 报 10，**错**）；`M2`＝**5／6**，**排除后＝3／5**（r1 报 6，**错**）⇒ 与 r2 更正后的数**逐字相同**。
- **`slices.Clone`**：`open.go:350-354` 实为 `audit := s.fkAudit; audit.ChildTables = slices.Clone(audit.ChildTables); return audit`；牙＝`:621` 用例；我造的 **M4**（`_ = slices.Clone(...)` 算出即弃）⇒ **红在 `:635`**（1 顶层/0 子例），报错串含 `shares its backing array`。

## ⑦ 边界零 diff（逐项，都是我两树各算一次）

- **`ddlV1`–`ddlV4`**：`fda736e1…`（4463／15 句）· `3e9ece70…`（903／1）· `d46ee18e…`（1286／5）· **`c17241cd…`（1357／5）**——`pretree`(`6d47e1ff4`) 与 `posttree`(HEAD) **四条逐字相同**；**12 名册**两树同＝`running,done,succeeded,cancelled,interrupted,stuck,error,Thinking,Settling,Muted,Stuck,Error`。
- **`196-r2`／`196-r3` 零 hunk**：四个边界函数 body 两树逐字同——`migrationDSN 8b40de05…`(3 行)／`migrationPragmas 97da5dd4…`(10)／`dsnPragmas 2e99efb5…`(8)／`requireForeignKeysOff 8050de10…`(25)；`c2e47f70c..HEAD` 四个 hunk 全在注释/clone 区，无函数体。
- **新增行 `DELETE`**：r2 那一笔 **0** 命中；`6d47e1ff4..HEAD` 的 `internal/memory` 新增行 **0** 命中（大小写敏感）。
- **r2 的产码改动面**：`numstat`＝件 `116/0`、用例 `216/3`、`open.go 15/5`、`schema.go 24/13`；非注释增删行只有 `"slices"`＋clone 四行（其余全是注释）⇒ 与件 §3/§9 的自述一致。

## ⑧ 五发变异（**我自己造的**，`-overlay`、仓外；逐名报红行 `logs/40-mutation-*`）

| 变异 | 我造的内容 | rc／红读数 | 红行 |
|---|---|---|---|
| `M1` | 摘掉 `audit = auditForeignKeys(...)` 那一跳 | **rc=1**；6 顶层／11 子例（排除新增＝4／8） | `:211,:236,:282×4,:337,:384,:387,:407×4`（＋新 `:550×3,:584,:607,:630`） |
| `M2` | 记录改成门槛（`OrphanRows>0 ⇒ Rollback+拒开`，硬边界 #1） | **rc=1**；5／6（排除新增＝3／5）；串 `Open refused a database with 3/2/1/1 orphan row(s)` | `:275×4,:333,:396`（＋新 `:600,:625`） |
| `M3` | 坏尺**谎报** `Checked=true` | **rc=1**；1／4 | `:499,:530,:554,:575,:580,:611` |
| `M4` | 去 clone（算出即弃） | **rc=1**；1／0 | `:635` |
| `M5` | 坏尺变锁（仅审计步 `!Checked` ⇒ 拒开） | **rc=1**；1／3 | `:524,:570,:593` |

- 还原核时机＝**最后一发变异（M5 窄形）之后**：`restoretree` 现导，全包 **rc=0／120 RUN／66 顶层 PASS／0 FAIL／2 SKIP**（`ok 25.247s`）；工作树 G3 提交前 **120/66/0/2**（`ok 26.086s`）、**提交后复跑同四数**（`ok 24.414s`）；工作树四枚 hash 与 `posttree`／blob 逐枚相同（⑤）。

## ⑨ 门禁 rc（逐门分开、⛔ 走管道取 rc）＋格式两尺

| 门 | 命令 | rc／读数 |
|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **rc=0**（0 字节输出） |
| G2 | `GOFLAGS= go vet ./internal/memory/` | **rc=0**（0 字节输出） |
| G3 | `GOFLAGS= go test -count=1 -v ./internal/memory/` | **rc=0**；四数分开：`ok 26.086s`／`=== RUN` **120**／`--- PASS`（顶层）**66**／`--- FAIL` **0**／`--- SKIP`（顶层）**2**（提交后复跑：rc=0、同四数、`ok 24.414s`） |
| G4 | `sh scripts/d22scan.sh` | **rc=0**（`clean - no D22 ban violations`） |
| G5 | `sh scripts/check-path-length-budget.sh` | **rc=0**（`VERDICT GREEN`；over-the-hat 61） |

- **四数对照（我自己在导出树上各跑一次，⛔ 相减）**：366 改前(`6d47e1ff4`) **101/59/0/2**（`ok 22.393s`）→ r1 交付(`c2e47f70c`) **114/64/0/2**（`ok 26.684s`）→ 本交付 **120/66/0/2** ⇒ r2 报的"改前 114/64 → 改后 120/66"**对得上**；`+6 RUN／+2 顶层 PASS` 也逐数对。
- **格式甲（入库 blob＝判据）**：`43ca90fc8` 的四枚 blob 导出到仓外再 `gofumpt -l` ⇒ **0 行**（rc=0）。**乙（工作树＝只诊断）**：`gofumpt -l internal/memory` 经 `[\\/]` 过滤 ⇒ **0 行**。**CRLF 尺**（⛔ 不是 gofumpt）：四枚 `tr -cd '\r'` ⇒ 全 **0**。

## ⑩ 提交＋逐笔越界

- 本件提交＝**`030352c6e`**（首笔 120 行；**记本行的一笔**只改本件、号见 `git log -1`）；`add`/`commit` 两处**显式 pathspec**、只点 `docs/evidence/s1/366-v2-verdict.md` 一枚；提交前 `git diff --cached --name-status` 只见 **`A`** 这一行（无 `D`/意外 `A`）。⛔ push。
- **越界＝0**：仓内 `.go` **0 枚**；`internal/**`／`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**` **0 字节改动**；五发变异全 `-overlay`，`git status --porcelain -- internal/memory` **空**；临时件全在 `D:\tmp\wisp\20261023\366v2\**`。

---

## 附：上报表（不占判语）

1. **`M5` 的形要写清**：`!audit.Checked` 若**不加** `step.auditForeignKeys` 守卫，会把 0→1/1→2/2→3 全拒（我试过：6 顶层/11 子例红，但红因是"所有迁移被拒"，不是'坏尺变锁'）；我用的是**窄形**（只咬审计步），红行 `:524,:570,:593` 与件一致。⇒ 后来者复现 `M5` 必须带守卫，否则读数不具可比性。
2. **r2 §4 那句"r1 同目录另三枚与 HEAD blob 逐字相同"**只在"HEAD＝r1 交付时"的读法下成立；**今天**那三枚与当前 HEAD 不同（r2 改了注释）。归因本身不受影响（值我逐枚复现），但复量者别拿今天的 HEAD 去对。
3. **残留（⛔ 本腿扩边界）**：`schema.go:308`（原 `:283`）那句 `The rebuild follows SQLite's documented 12-step procedure` **仍把 12 步形状当已核事实**——它**不属** v1 点名的 8 处、也不在 r2 写面；若要一并加对冲，须另立票。
4. **前提"三条路径"实测是四条**（+`retention` 的 NULL 时间戳子行）；我这一腿**只独立重量了 purge 一条**（3 孤儿），另三条由交付用例 4 子例（3/2/1/1）绿＋v1/r2 实测背书——尺的射程如实写在这里。
